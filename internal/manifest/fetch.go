package manifest

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const fetchTimeout = 15 * time.Second

// Fetch downloads and verifies a signed manifest from url.
//
// If pubkeyHex is empty, signature check is skipped (dev mode).
// Otherwise the manifest's signature MUST verify against pubkeyHex AND the
// manifest's embedded pubkey field must match (defense in depth).
func Fetch(ctx context.Context, url, pubkeyHex string) (*Manifest, error) {
	ctx, cancel := context.WithTimeout(ctx, fetchTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("manifest http %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20)) // 4MB cap
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}

	var signed Signed
	if err := json.Unmarshal(raw, &signed); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}

	if pubkeyHex != "" {
		if err := verifySignature(&signed, pubkeyHex); err != nil {
			return nil, fmt.Errorf("verify manifest: %w", err)
		}
	}

	return &signed.Manifest, nil
}

func verifySignature(s *Signed, expectedPubkeyHex string) error {
	if s.PubkeyHex != expectedPubkeyHex {
		return fmt.Errorf("pubkey mismatch: manifest=%s expected=%s", s.PubkeyHex, expectedPubkeyHex)
	}
	pub, err := hex.DecodeString(expectedPubkeyHex)
	if err != nil || len(pub) != ed25519.PublicKeySize {
		return fmt.Errorf("invalid pubkey hex")
	}
	sig, err := hex.DecodeString(s.SignatureHex)
	if err != nil || len(sig) != ed25519.SignatureSize {
		return fmt.Errorf("invalid signature hex")
	}
	// Canonical bytes = JSON of inner manifest with sorted keys.
	canonical, err := canonicalJSON(&s.Manifest)
	if err != nil {
		return fmt.Errorf("canonicalize: %w", err)
	}
	if !ed25519.Verify(pub, canonical, sig) {
		return fmt.Errorf("signature invalid")
	}
	return nil
}

// canonicalJSON produces a stable byte representation for signing/verifying.
// json.Marshal already emits map keys in sorted order, and struct fields in
// declaration order, so this is sufficient for our schema.
func canonicalJSON(m *Manifest) ([]byte, error) {
	return json.Marshal(m)
}
