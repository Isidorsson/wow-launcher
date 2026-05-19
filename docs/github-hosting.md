# Hosting patches on GitHub

This guide shows server owners how to host MPQs (and any other patch files) on
GitHub Releases, with a GitHub Action that builds and signs the manifest
automatically every time you publish a release.

You upload files. The Action does the rest. The launcher picks up changes the
next time a player clicks Sync (or on next launch, if you've enabled the
startup check).

## Why GitHub Releases

- **Free**, 2 GB per file, unlimited bandwidth
- Served from GitHub's CDN (`objects.githubusercontent.com`) — fast worldwide
- Immutable asset URLs per release — no cache-busting needed
- `releases/latest/download/<name>` always redirects to your newest release
- Zero infrastructure: no S3, no Cloudflare account, no domain

GitHub repo blobs (committing the MPQ into git) won't work — 100 MB hard limit,
throttled raw URL. Use Releases.

## How it fits together

```
┌──────────────────────────────────────────────────────────────────────┐
│  You (server owner)                                                  │
│    1. Edit patches.toml in your repo                                 │
│    2. Create a new GitHub Release, drag in your MPQs / other files   │
│    3. Click "Publish release"                                        │
└──────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌──────────────────────────────────────────────────────────────────────┐
│  GitHub Action (fires on `release: published`)                       │
│    1. Downloads release assets                                       │
│    2. Hashes each file (SHA256)                                      │
│    3. Builds manifest.json from patches.toml + hashes                │
│    4. Signs it with your Ed25519 private key (stored as a secret)    │
│    5. Uploads manifest.json back to the same release                 │
└──────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌──────────────────────────────────────────────────────────────────────┐
│  Launcher (on player's machine)                                      │
│    - manifest_url in config.toml points at                           │
│      https://github.com/<you>/<repo>/releases/latest/download/       │
│      manifest.json                                                   │
│    - Verifies signature against embedded public key                  │
│    - Downloads only files whose local SHA256 doesn't match           │
└──────────────────────────────────────────────────────────────────────┘
```

## One-time setup

### 1. Generate a signing keypair

On your dev machine, from the launcher repo:

```sh
go run ./tools/sign-manifest -gen-key
```

Output:

```
public:  6f3c...
private: 2b91...
```

**The private key is a secret.** Anyone with it can publish manifests that your
players' launchers will trust. Store it once, then close the terminal.

### 2. Bake the public key into the launcher

Edit `config.toml`:

```toml
[security]
manifest_pubkey_hex = "6f3c..."   # paste the PUBLIC hex here

[[servers]]
id            = "main"
name          = "My Realm"
manifest_url  = "https://github.com/<you>/<repo>/releases/latest/download/manifest.json"
website       = "https://myrealm.example"
news_feed_url = "https://myrealm.example/news.json"
```

Rebuild the launcher binary. The config is embedded at compile time.

### 3. Create your patches repository

This is **separate** from the launcher source repo. It's a small repo that
exists only to host releases. Suggested layout:

```
my-realm-patches/
├── .github/
│   └── workflows/
│       └── release.yml          # the Action below
└── patches.toml                 # list of files + destinations
```

No MPQs in git — they go in the Release, not the repo.

### 4. Add your private key as a repo secret

In the patches repo on GitHub:

- Settings → Secrets and variables → Actions → New repository secret
- Name: `MANIFEST_PRIV_KEY`
- Value: the **private** hex from step 1

### 5. Write `patches.toml`

This file maps GitHub release asset names to install destinations. Edit it
whenever you add, remove, or rename a patch.

```toml
# Top-level fields baked into every manifest the Action builds.
server_id            = "main"
client_version       = "3.3.5a-12340"
realmlist            = "logon.myrealm.example"
locale               = "enUS"
min_launcher_version = "0.1.0"

# One [[file]] per patch. `asset` = filename you'll upload to the Release.
# `path` = where it lands inside the WoW install.

[[file]]
asset    = "patch-4.MPQ"
path     = "Data/patch-4.MPQ"
required = true

[[file]]
asset    = "patch-5.MPQ"
path     = "Data/patch-5.MPQ"
required = true

[[file]]
asset    = "patch-enUS-4.MPQ"
path     = "Data/enUS/patch-enUS-4.MPQ"
required = true

[[file]]
asset    = "hd-textures.MPQ"
path     = "Data/patch-9.MPQ"
required = false
label    = "HD textures (2 GB, optional)"

[[file]]
asset    = "CustomUI.toc"
path     = "Interface/AddOns/CustomUI/CustomUI.toc"
required = true
```

**Notes:**

- `path` uses forward slashes and is relative to the WoW install root
- Non-MPQ files work — addons, WTF configs, anything the launcher should drop
  into the install
- Use [patch slot conventions](manifest.md#patch-slot-conventions) to avoid
  collisions (slots 4–9 = required, A–Z = optional)
- The launcher ships SHA256 mismatches automatically; renaming an asset on
  GitHub re-downloads it for every player on next sync

### 6. Add the workflow

Save this as `.github/workflows/release.yml` in the patches repo:

```yaml
name: Build manifest

on:
  release:
    types: [published]

jobs:
  manifest:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Clone launcher (for sign-manifest tool)
        run: git clone --depth 1 https://github.com/<you>/wow-launcher.git launcher

      - name: Download release assets
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          mkdir -p assets
          gh release download "${{ github.event.release.tag_name }}" \
            --repo "${{ github.repository }}" \
            --dir ./assets

      - name: Write private key
        run: echo "${{ secrets.MANIFEST_PRIV_KEY }}" > priv.hex

      - name: Build + sign manifest
        run: |
          cd launcher
          go run ./tools/sign-manifest build \
            --patches-toml ../patches.toml \
            --assets-dir ../assets \
            --release-tag "${{ github.event.release.tag_name }}" \
            --repo "${{ github.repository }}" \
            --key ../priv.hex \
            --out ../manifest.json

      - name: Upload manifest to release
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          gh release upload "${{ github.event.release.tag_name }}" manifest.json \
            --repo "${{ github.repository }}" \
            --clobber
```

> The `build` subcommand of `sign-manifest` is part of the launcher repo. The
> workflow clones it on the fly so you don't have to vendor the tool. If you
> prefer, copy `tools/sign-manifest` into your patches repo and skip the clone
> step.

## Per-release workflow

Every time you ship a patch:

1. **Edit `patches.toml`** if you've added a new file, renamed an asset, or
   changed `client_version` / `realmlist`. Commit and push.
2. **Create a new release on GitHub:**
   - Releases → Draft a new release
   - Tag: `v2`, `v3`, etc. (each release must be a new tag)
   - Drag-drop your MPQs and any other files into the asset uploader
   - Click **Publish release**
3. **Wait ~30 seconds.** The Action runs, builds `manifest.json`, attaches it
   to the same release.
4. **Done.** Players' launchers hit `releases/latest/download/manifest.json`
   on next sync and download anything whose SHA256 changed.

You don't touch the launcher binary. You don't redistribute it. The manifest
URL is stable, the public key is stable, the launcher just rechecks.

## When does the launcher check for updates?

| Trigger        | What it does                                                              |
|----------------|---------------------------------------------------------------------------|
| On startup     | Background conditional fetch (`If-None-Match`). Shows banner if changed.  |
| Sync button    | Fetches manifest, downloads any file whose SHA256 mismatches.             |
| Play button    | If `min_launcher_version` exceeds installed launcher, refuses.            |

The startup check stores each server's manifest content hash + ETag under
`<UserConfigDir>/WowLauncher/state/manifest-state.json`. On next launch it
sends the saved ETag as `If-None-Match`. If GitHub returns `304 Not Modified`
no bandwidth is spent; if the content hash actually differs from last seen,
the launcher emits `update:available` and the UI shows a banner pointing the
player at the Sync button.

The downloader is SHA256-driven: a file is considered up-to-date if its hash
matches the manifest. Renaming a file, changing its content, or bumping a
patch version all just change the hash and trigger a re-download — there is
no separate "version" field per file.

## Troubleshooting

**"signature invalid" on player launch**
The public key in `config.toml` doesn't match the private key the Action signed
with. Regenerate, rebuild launcher.

**"pubkey mismatch"**
The manifest's embedded `pubkey` field doesn't match the one in `config.toml`.
This is a separate check from signature verification — same fix.

**Action fails: "asset X not found in release"**
You listed an asset in `patches.toml` that wasn't uploaded to the release.
Either upload it or remove the entry.

**Launcher downloads the whole thing every time**
The manifest's recorded SHA256 doesn't match the file on the CDN. Usually
means you uploaded a file, hashed it, then re-uploaded a different one without
re-running the Action. Re-publish the release.

**Player on slow connection times out**
The downloader resumes — partial files are kept on disk and continue next sync.
No special config needed.

**Player reports stuck "update available" banner**
The launcher caches the last-seen manifest content hash + ETag at
`<UserConfigDir>/WowLauncher/state/manifest-state.json`
(`%LOCALAPPDATA%\WowLauncher\state\manifest-state.json` on Windows). Deleting
that file forces a clean re-check on next launch. The cache is best-effort —
a missing or corrupt file is treated as "no prior knowledge", never blocks
launching.

## What this guide does NOT cover

- **News feed JSON** — separate file at `news_feed_url`. See
  [news.go](../internal/news/news.go) for the schema. Host alongside the
  manifest if you like, or anywhere else.
- **Custom CDN** (Cloudflare R2, Backblaze B2, etc.) — use `manifest.md` and
  sign manually with `sign-manifest -key priv.hex in.json out.json`.
- **Private patch distribution** — GitHub Releases are public. For paid
  realms or beta tests, use a private CDN with signed URLs.

## See also

- [manifest.md](manifest.md) — full manifest schema reference
- [README.md](../README.md) — launcher overview and build instructions
