package manifest

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
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
	m, _, _, err := fetchConditional(ctx, url, pubkeyHex, "")
	return m, err
}

// FetchResult bundles the outputs of a conditional fetch.
//
// Changed=false + Manifest=nil means the server returned 304 Not Modified.
// Changed=true means we either got a fresh 200 or the manifest content hash
// differs from PrevContentHash (passed in by the caller).
type FetchResult struct {
	Manifest    *Manifest
	ETag        string // value from response, store and send next time
	ContentHash string // sha256 of canonical inner manifest JSON
	Changed     bool   // true if content differs from prevContentHash or no prev known
}

// FetchWithETag performs a conditional GET using prevETag (sent as If-None-Match).
//
// On 304, returns Manifest=nil, Changed=false, ETag=prevETag (unchanged).
// On 200, parses + verifies, then compares ContentHash to prevContentHash to
// decide if the manifest's actual content changed (a server might serve the
// same bytes with a new ETag header — we want to know about content drift).
func FetchWithETag(ctx context.Context, url, pubkeyHex, prevETag, prevContentHash string) (*FetchResult, error) {
	m, etag, hash, err := fetchConditional(ctx, url, pubkeyHex, prevETag)
	if err != nil {
		return nil, err
	}
	r := &FetchResult{ETag: etag, ContentHash: hash}
	if m == nil {
		// 304 — content known unchanged. Preserve prev ETag if server omitted.
		if r.ETag == "" {
			r.ETag = prevETag
		}
		r.Changed = false
		return r, nil
	}
	r.Manifest = m
	r.Changed = hash != prevContentHash
	return r, nil
}

func fetchConditional(ctx context.Context, url, pubkeyHex, ifNoneMatch string) (*Manifest, string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, fetchTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", "", fmt.Errorf("build request: %w", err)
	}
	if ifNoneMatch != "" {
		req.Header.Set("If-None-Match", ifNoneMatch)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", "", fmt.Errorf("fetch manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		return nil, resp.Header.Get("ETag"), "", nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, "", "", fmt.Errorf("manifest http %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20)) // 4MB cap
	if err != nil {
		return nil, "", "", fmt.Errorf("read manifest: %w", err)
	}

	var signed Signed
	if err := json.Unmarshal(raw, &signed); err != nil {
		return nil, "", "", fmt.Errorf("parse manifest: %w", err)
	}

	if pubkeyHex != "" {
		if err := verifySignature(&signed, pubkeyHex); err != nil {
			return nil, "", "", fmt.Errorf("verify manifest: %w", err)
		}
	}

	canonical, err := canonicalJSON(&signed.Manifest)
	if err != nil {
		return nil, "", "", fmt.Errorf("canonicalize: %w", err)
	}
	sum := sha256.Sum256(canonical)
	return &signed.Manifest, resp.Header.Get("ETag"), hex.EncodeToString(sum[:]), nil
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
