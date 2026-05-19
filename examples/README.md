# Examples for server operators

Copy-paste templates for the GitHub Releases hosting workflow. See
[`../docs/github-hosting.md`](../docs/github-hosting.md) for the full guide.

## Files

| File             | Goes to (in YOUR patches repo)            | Purpose                                    |
|------------------|-------------------------------------------|--------------------------------------------|
| `release.yml`    | `.github/workflows/release.yml`           | Builds + signs manifest on release publish |
| `patches.toml`   | `patches.toml` (repo root)                | Maps release assets to install destinations |

## Setup checklist

1. Generate keypair (run from launcher repo, NOT patches repo):
   ```sh
   go run ./tools/sign-manifest -gen-key
   ```
2. Create a fresh GitHub repo for your patches (separate from the launcher repo).
3. In that repo: Settings -> Secrets and variables -> Actions ->
   New repository secret -> `MANIFEST_PRIV_KEY` = the **private** hex.
4. Copy `release.yml` to `.github/workflows/release.yml` in the patches repo.
   **Edit the `<OWNER>` placeholder** to point at the launcher fork you're using.
5. Copy `patches.toml` to the patches repo root. Edit `server_id`, `realmlist`,
   etc., and the `[[file]]` list to match the assets you'll upload.
6. Commit + push.
7. Set the launcher's `manifest_url` to
   `https://github.com/<your-org>/<patches-repo>/releases/latest/download/manifest.json`
   and the **public** hex to `[security] manifest_pubkey_hex` in `config.toml`.
   Rebuild the launcher binary.
8. Create a GitHub Release in the patches repo: pick a new tag (`v1`, `v2`, ...),
   drag in your MPQs and other asset files (their filenames must match the
   `asset = ` values in `patches.toml`), click Publish.
9. The Action runs in ~30 seconds and attaches a signed `manifest.json` to the
   same release. Players' launchers pick it up on next startup or Sync.
