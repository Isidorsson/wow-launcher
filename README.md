# wow-launcher

Wails-based launcher for WoW 3.3.5a private servers.
Open source — fork it, edit `config.toml`, host a signed manifest, ship.

## Features

- Multi-server: one launcher, dropdown to switch realms.
- Per-server profiles: each realm gets its own isolated install (hardlinked
  base files, downloaded patches) so MPQs from server A don't conflict with B.
- Resumable, hash-verified, concurrent-chunk downloads (`Range` requests).
- Ed25519-signed manifests so a CDN compromise can't push malware patches.
- Startup update check with ETag caching — 304 on unchanged manifests, banner
  when content drifts.
- Auto-detect existing 3.3.5 install + folder picker + drag-drop install path.
- News feed per server.

## Architecture

```
Wails binary (Go + Svelte) ──HTTPS──> manifest server (small JSON, signed)
                          ──HTTPS──> CDN (patch MPQs, immutable URLs)
```

Patches live in slots `patch-4.MPQ`–`patch-Z.MPQ` (+ locale variants).
See `docs/manifest.md` for the schema.

## Quickstart (server operator)

**Recommended path — GitHub Releases + Action:** see
[docs/github-hosting.md](docs/github-hosting.md). Zero CDN setup, drag-drop
MPQ uploads, Action signs the manifest on `release: published`. The startup
update check (ETag + content-hash cache) is automatic — no per-player config.

**Manual path — your own CDN:**

1. Fork this repo.
2. Edit `config.toml`:
   - `branding.launcher_name`
   - `servers[0].manifest_url` (where you'll host the manifest)
   - `security.manifest_pubkey_hex` (run `go run ./tools/sign-manifest -gen-key`)
3. Build patches into MPQs (use [MPQEditor] or [StormLib]).
4. Upload patches to your CDN (Cloudflare R2 + bunny.net recommended).
5. Either:
   - Author `manifest-unsigned.json` by hand (see `docs/manifest.md`) and sign:
     `go run ./tools/sign-manifest -key priv.hex manifest-unsigned.json manifest.json`
   - Or use `sign-manifest build --patches-toml ... --assets-dir ... --key ...`
     to hash + sign in one shot.
6. Upload `manifest.json` to `manifest_url`.
7. Build launcher: `wails build -platform windows/amd64`.
8. Ship `build/bin/wow-launcher.exe` to users.

[MPQEditor]: http://www.zezula.net/en/mpq/download.html
[StormLib]: https://github.com/ladislav-zezula/StormLib

## Dev

```sh
wails dev          # hot-reload dev server
wails build        # release build
go build ./...     # backend-only compile check
go test ./...      # backend tests
```

## Status

Backend end-to-end ready. Frontend wired. `go test ./...` green.

Not yet implemented (good next-session work):
- Self-update for the launcher binary itself
- Profile delete / repair UI
- Cancel button during sync (wire ctx.CancelFunc)
- Play button gate on `min_launcher_version`

## License

MIT — pick a license file to include before publishing.
