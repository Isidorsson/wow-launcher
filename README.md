# wow-launcher

Wails-based launcher for WoW 3.3.5a private servers.
Open source — fork it, edit `config.toml`, host a signed manifest, ship.

## Features

- Multi-server: one launcher, dropdown to switch realms.
- Per-server profiles: each realm gets its own isolated install (hardlinked
  base files, downloaded patches) so MPQs from server A don't conflict with B.
- Resumable, hash-verified, concurrent-chunk downloads (`Range` requests).
- Ed25519-signed manifests so a CDN compromise can't push malware patches.
- Auto-detect existing 3.3.5 install + manual path override.

## Architecture

```
Wails binary (Go + Svelte) ──HTTPS──> manifest server (small JSON, signed)
                          ──HTTPS──> CDN (patch MPQs, immutable URLs)
```

Patches live in slots `patch-4.MPQ`–`patch-Z.MPQ` (+ locale variants).
See `docs/manifest.md` for the schema.

## Quickstart (server operator)

1. Fork this repo.
2. Edit `config.toml`:
   - `branding.launcher_name`
   - `servers[0].manifest_url` (where you'll host the manifest)
   - `security.manifest_pubkey_hex` (run `go run ./tools/sign-manifest -gen-key`)
3. Build patches into MPQs (use [MPQEditor] or [StormLib]).
4. Upload patches to your CDN (Cloudflare R2 + bunny.net recommended).
5. Author `manifest-unsigned.json` (see docs).
6. Sign: `go run ./tools/sign-manifest -key priv.hex manifest-unsigned.json manifest.json`.
7. Upload `manifest.json` to `manifest_url`.
8. Build launcher: `wails build -platform windows/amd64`.
9. Ship `build/bin/wow-launcher.exe` to users.

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
- Folder picker dialog (runtime.OpenDirectoryDialog) instead of typed path
- Self-update for the launcher binary itself
- News feed fetch + render
- Profile delete / repair UI
- Cancel button during sync (wire ctx.CancelFunc)

## License

MIT — pick a license file to include before publishing.
