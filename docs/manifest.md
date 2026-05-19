# Manifest format

The manifest is the source of truth for what patches the launcher pulls.
You host it at a stable HTTPS URL; the launcher fetches it on every Sync,
and (with ETag-aware conditional GET) on startup to detect changes.

> **Most operators should follow [github-hosting.md](github-hosting.md)** —
> drag-drop patches to a GitHub Release, an Action signs the manifest. This
> page is the underlying schema + the manual workflow for custom CDNs.

## Inner manifest (what you sign)

```json
{
  "schema_version": 1,
  "server_id": "demo",
  "client_version": "3.3.5a-12340",
  "realmlist": "logon.example.com",
  "locale": "enUS",
  "min_launcher_version": "0.1.0",
  "files": [
    {
      "path": "Data/patch-4.MPQ",
      "url": "https://cdn.example.com/demo/patch-4-v3.MPQ",
      "size": 524288000,
      "sha256": "<hex>",
      "required": true
    },
    {
      "path": "Data/patch-9.MPQ",
      "url": "https://cdn.example.com/demo/hd-textures-v1.MPQ",
      "size": 2147483648,
      "sha256": "<hex>",
      "required": false,
      "label": "HD textures (2 GB, optional)"
    }
  ]
}
```

## Signed envelope (what you publish)

```json
{
  "manifest": { /* the inner manifest, exactly as above */ },
  "signature": "<hex Ed25519 sig over canonical JSON of inner manifest>",
  "pubkey": "<hex Ed25519 pubkey, must match config.toml>"
}
```

## Workflow

```sh
# One-time: generate a keypair. Keep private SECRET (server/CI only).
go run ./tools/sign-manifest -gen-key

# Per release — choice of two flows:

# A) Hand-written manifest (full control, custom CDN URLs):
go run ./tools/sign-manifest -key priv.hex manifest-unsigned.json manifest.json

# B) Build from a patches.toml + assets directory (hashes + signs in one shot):
go run ./tools/sign-manifest build \
  --patches-toml patches.toml \
  --assets-dir ./assets \
  --release-tag v3 \
  --repo owner/repo \
  --key priv.hex \
  --out manifest.json

# Upload manifest.json to your manifest_url; upload patches to CDN.
```

## Patch slot conventions

The WoW client loads patches in lexicographic order. Use slots Blizzard didn't:

| Slot                       | Use for                                           |
|----------------------------|---------------------------------------------------|
| `patch-4.MPQ` ... `patch-9.MPQ` | required server data (DBCs, custom maps)     |
| `patch-A.MPQ` ... `patch-Z.MPQ` | overrides / cosmetic / HD packs              |
| `patch-{locale}-4.MPQ` ... `patch-{locale}-Z.MPQ` | locale-specific equivalents |

**Rule**: higher letter wins. If two patches contain the same internal file,
the one with the later slot overrides. Reserve specific slots in your project
so cosmetic packs don't shadow critical fixes. Common convention:
- 4–6: required game fixes
- 7–9: required cosmetics
- A–F: optional player-installable packs
- W–Z: dev/test overrides

## Versioned CDN filenames

Always rename a patch's URL when its content changes (e.g. `patch-4-v3.MPQ`).
The launcher matches by SHA256, not filename, but immutable URLs let you set
infinite cache TTL on the CDN. The manifest itself stays small and is the
only resource you ever invalidate.
