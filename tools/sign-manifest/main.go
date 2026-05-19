// sign-manifest is a tiny CLI for private-server owners.
//
// Modes:
//
//	sign-manifest -gen-key
//	    Generate an Ed25519 keypair and print hex to stdout.
//
//	sign-manifest -key priv.hex in.json out.json
//	    Sign a hand-written unsigned inner Manifest JSON into a Signed envelope.
//
//	sign-manifest build --patches-toml patches.toml --assets-dir ./assets \
//	    --release-tag v3 --repo owner/repo --key priv.hex --out manifest.json
//	    Walk assets-dir, hash each file referenced in patches.toml, build URLs
//	    from --release-tag + --repo, sign, write signed envelope.
//
// The unsigned input/output schema is internal/manifest.
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"wow-launcher/internal/manifest"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "build" {
		runBuild(os.Args[2:])
		return
	}

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
		fmt.Fprintln(os.Stderr, "   or: sign-manifest -gen-key")
		fmt.Fprintln(os.Stderr, "   or: sign-manifest build --help")
		os.Exit(2)
	}

	priv := loadPrivKey(*keyFile)
	inPath, outPath := flag.Arg(0), flag.Arg(1)
	raw, err := os.ReadFile(inPath)
	if err != nil {
		die("read input: %v", err)
	}
	var m manifest.Manifest
	if err := json.Unmarshal(raw, &m); err != nil {
		die("parse input: %v", err)
	}
	writeSigned(outPath, &m, priv)
	fmt.Printf("signed %s -> %s\n", inPath, outPath)
}

// patchesTOML is the source-of-truth file in a server owner's patches repo.
// Top-level fields populate Manifest's scalar fields; [[file]] entries are
// matched against release assets by Asset name to produce Manifest.Files.
type patchesTOML struct {
	ServerID           string    `toml:"server_id"`
	ClientVersion      string    `toml:"client_version"`
	Realmlist          string    `toml:"realmlist"`
	Locale             string    `toml:"locale"`
	MinLauncherVersion string    `toml:"min_launcher_version"`
	File               []tomlEnt `toml:"file"`
}

type tomlEnt struct {
	Asset    string `toml:"asset"`    // release asset filename (in --assets-dir)
	Path     string `toml:"path"`     // destination inside WoW install
	Required bool   `toml:"required"`
	Label    string `toml:"label"`
}

func runBuild(args []string) {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	patchesPath := fs.String("patches-toml", "patches.toml", "path to patches.toml")
	assetsDir := fs.String("assets-dir", "./assets", "directory containing release assets")
	releaseTag := fs.String("release-tag", "", "GitHub release tag (e.g. v3) — required")
	repo := fs.String("repo", "", "GitHub repo as owner/name — required")
	keyFile := fs.String("key", "", "path to hex-encoded private key — required")
	outPath := fs.String("out", "manifest.json", "output path for signed manifest")
	urlBase := fs.String("url-base", "", "override URL base (default: https://github.com/<repo>/releases/download/<tag>)")
	_ = fs.Parse(args)

	if *releaseTag == "" || *repo == "" || *keyFile == "" {
		fmt.Fprintln(os.Stderr, "build: --release-tag, --repo, --key are required")
		fs.Usage()
		os.Exit(2)
	}

	var pt patchesTOML
	rawTOML, err := os.ReadFile(*patchesPath)
	if err != nil {
		die("read %s: %v", *patchesPath, err)
	}
	if err := toml.Unmarshal(rawTOML, &pt); err != nil {
		die("parse %s: %v", *patchesPath, err)
	}
	if len(pt.File) == 0 {
		die("%s: no [[file]] entries", *patchesPath)
	}

	base := *urlBase
	if base == "" {
		base = fmt.Sprintf("https://github.com/%s/releases/download/%s", *repo, *releaseTag)
	}

	files := make([]manifest.File, 0, len(pt.File))
	for _, ent := range pt.File {
		if ent.Asset == "" || ent.Path == "" {
			die("entry missing asset or path: %+v", ent)
		}
		localPath := filepath.Join(*assetsDir, ent.Asset)
		size, sum, err := hashFile(localPath)
		if err != nil {
			die("hash %s: %v", localPath, err)
		}
		files = append(files, manifest.File{
			Path:     ent.Path,
			URL:      base + "/" + ent.Asset,
			Size:     size,
			SHA256:   sum,
			Required: ent.Required,
			Label:    ent.Label,
		})
		fmt.Printf("  + %s (%d bytes, sha256 %s)\n", ent.Asset, size, sum[:12])
	}

	m := manifest.Manifest{
		SchemaVersion:      1,
		ServerID:           pt.ServerID,
		ClientVersion:      pt.ClientVersion,
		Realmlist:          pt.Realmlist,
		Locale:             pt.Locale,
		MinLauncherVersion: pt.MinLauncherVersion,
		Files:              files,
	}

	priv := loadPrivKey(*keyFile)
	writeSigned(*outPath, &m, priv)
	fmt.Printf("built + signed %d files -> %s\n", len(files), *outPath)
}

func hashFile(path string) (int64, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, "", err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return 0, "", err
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return 0, "", err
	}
	return info.Size(), hex.EncodeToString(h.Sum(nil)), nil
}

func loadPrivKey(path string) ed25519.PrivateKey {
	raw, err := os.ReadFile(path)
	if err != nil {
		die("read key: %v", err)
	}
	priv, err := hex.DecodeString(string(trimSpace(raw)))
	if err != nil || len(priv) != ed25519.PrivateKeySize {
		die("bad private key hex (need %d bytes)", ed25519.PrivateKeySize)
	}
	return ed25519.PrivateKey(priv)
}

func writeSigned(outPath string, m *manifest.Manifest, priv ed25519.PrivateKey) {
	canonical, err := json.Marshal(m)
	if err != nil {
		die("canonicalize: %v", err)
	}
	sig := ed25519.Sign(priv, canonical)
	pub := priv.Public().(ed25519.PublicKey)

	signed := manifest.Signed{
		Manifest:     *m,
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
