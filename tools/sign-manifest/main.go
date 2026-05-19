// sign-manifest is a tiny CLI for private-server owners.
//
// Usage:
//   sign-manifest -gen-key                    generate keypair, print to stdout
//   sign-manifest -key priv.hex in.json out.json  sign in.json -> out.json (Signed envelope)
//
// The unsigned input is the inner Manifest struct (see internal/manifest/schema.go).
// The output is a Signed envelope ready to host at the manifest_url in config.toml.
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"wow-launcher/internal/manifest"
)

func main() {
	gen := flag.Bool("gen-key", false, "generate Ed25519 keypair and exit")
	keyFile := flag.String("key", "", "path to hex-encoded private key (64 bytes)")
	flag.Parse()

	if *gen {
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			die("gen: %v", err)
		}
		fmt.Printf("public:  %s\n", hex.EncodeToString(pub))
		fmt.Printf("private: %s\n", hex.EncodeToString(priv))
		fmt.Println("\nKeep private SECRET. Put public in config.toml manifest_pubkey_hex.")
		return
	}

	if *keyFile == "" || flag.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "usage: sign-manifest -key priv.hex in.json out.json")
		os.Exit(2)
	}

	privHex, err := os.ReadFile(*keyFile)
	if err != nil {
		die("read key: %v", err)
	}
	priv, err := hex.DecodeString(string(trimSpace(privHex)))
	if err != nil || len(priv) != ed25519.PrivateKeySize {
		die("bad private key hex (need %d bytes)", ed25519.PrivateKeySize)
	}

	inPath, outPath := flag.Arg(0), flag.Arg(1)
	raw, err := os.ReadFile(inPath)
	if err != nil {
		die("read input: %v", err)
	}
	var m manifest.Manifest
	if err := json.Unmarshal(raw, &m); err != nil {
		die("parse input: %v", err)
	}

	canonical, err := json.Marshal(&m)
	if err != nil {
		die("canonicalize: %v", err)
	}
	sig := ed25519.Sign(priv, canonical)
	pub := ed25519.PrivateKey(priv).Public().(ed25519.PublicKey)

	signed := manifest.Signed{
		Manifest:     m,
		SignatureHex: hex.EncodeToString(sig),
		PubkeyHex:    hex.EncodeToString(pub),
	}
	out, err := json.MarshalIndent(&signed, "", "  ")
	if err != nil {
		die("marshal signed: %v", err)
	}
	if err := os.WriteFile(outPath, out, 0o644); err != nil {
		die("write output: %v", err)
	}
	fmt.Printf("signed %s -> %s\n", inPath, outPath)
}

func trimSpace(b []byte) []byte {
	for len(b) > 0 && (b[len(b)-1] == '\n' || b[len(b)-1] == '\r' || b[len(b)-1] == ' ' || b[len(b)-1] == '\t') {
		b = b[:len(b)-1]
	}
	for len(b) > 0 && (b[0] == '\n' || b[0] == '\r' || b[0] == ' ' || b[0] == '\t') {
		b = b[1:]
	}
	return b
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
