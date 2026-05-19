package profile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// baseFiles is the set of files copied from the user's source 3.3.5 install
// into each per-server profile. Everything else (custom patches) is downloaded
// fresh per profile from the server's manifest.
//
// Paths use forward slashes and are joined locale-side at runtime.
var baseFiles = []string{
	"Wow.exe",
	"Data/common.MPQ",
	"Data/common-2.MPQ",
	"Data/expansion.MPQ",
	"Data/lichking.MPQ",
	"Data/patch.MPQ",
	"Data/patch-2.MPQ",
	"Data/patch-3.MPQ",
}

// localeBaseFiles are joined with the locale string at runtime.
var localeBaseFiles = []string{
	"Data/{locale}/locale-{locale}.MPQ",
	"Data/{locale}/expansion-locale-{locale}.MPQ",
	"Data/{locale}/lichking-locale-{locale}.MPQ",
	"Data/{locale}/patch-{locale}.MPQ",
	"Data/{locale}/patch-{locale}-2.MPQ",
	"Data/{locale}/patch-{locale}-3.MPQ",
}

// MaterializeBase populates a profile directory with base client files from
// the user's source install.
//
// ★ LEARNING GAP — implement this yourself.
//
// Why this matters: each profile needs the base ~5GB of client files. Copying
// costs disk space (5GB × num servers). Hardlinking shares one inode across
// profiles = near-zero cost. But hardlinks only work within the same NTFS
// volume; cross-volume falls back to copy. Symlinks need admin on Windows,
// avoid them.
//
// Implementation outline:
//
//   1. Build the full file list: baseFiles + localeBaseFiles with {locale}
//      replaced by the locale arg.
//   2. For each entry:
//        src = filepath.Join(baseInstall, entry)
//        dst = filepath.Join(profileRoot, entry)
//      Ensure dst's parent dir exists with os.MkdirAll.
//   3. Stat src — if missing, skip (Blizzlike installs vary; e.g. older
//      versions don't have patch-3.MPQ). Skip = continue, not error.
//   4. If dst already exists, skip (idempotent — caller may re-run).
//   5. Try os.Link(src, dst) first. If it fails with a cross-device or
//      not-supported error (check errors.Is patterns: syscall.EXDEV on Unix,
//      err.Error() containing "different drive" or "not supported" works on
//      Windows), fall back to a streaming copy:
//        - open src, create dst, io.Copy, close both.
//        - On any error during copy, os.Remove the partial dst.
//   6. Return the first hard error. Don't aggregate — fail fast so the user
//      sees the actual cause.
//
// Tests you should write:
//   - happy path with a fake source tree
//   - missing optional file (patch-3) does not error
//   - hardlink failure path falls through to copy
//   - copy is atomic-ish: failed copy doesn't leave a half file
//
// Hint: filepath.WalkDir is NOT what you want here — we have a known list.
func MaterializeBase(baseInstall, profileRoot, locale string) error {
	// TODO: implement
	_ = baseFiles
	_ = localeBaseFiles
	_ = io.Copy
	_ = os.Link
	_ = filepath.Join
	return fmt.Errorf("MaterializeBase not implemented yet — see comment in materialize.go")
}
