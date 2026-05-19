package profile

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"wow-launcher/internal/events"
)

// baseFiles is the set of files copied from the user's source 3.3.5 install
// into each per-server profile. Everything else (custom patches) is downloaded
// fresh per profile from the server's manifest.
//
// Paths use forward slashes; converted with filepath.FromSlash at use site.
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

// MaterializeBase populates profileRoot with base client files from baseInstall.
//
// Strategy: hardlink each file. If the link syscall fails because src and dst
// are on different volumes (Windows: ERROR_NOT_SAME_DEVICE = syscall.Errno(17);
// POSIX: EXDEV), fall back to a streaming copy. This keeps disk cost ≈ patches
// only when everything is on one NTFS volume, while still working across drives.
//
// Idempotent: missing optional source files are skipped (Blizzlike installs
// vary in which Blizz patches they have applied), and existing dst files are
// left alone so callers can re-run after a partial failure.
//
// Fail-fast: returns on the first hard error so the user sees the real cause
// rather than a wall of secondary failures.
func MaterializeBase(ctx context.Context, baseInstall, profileRoot, locale string) error {
	if baseInstall == "" || profileRoot == "" || locale == "" {
		return fmt.Errorf("baseInstall, profileRoot, locale all required")
	}

	all := make([]string, 0, len(baseFiles)+len(localeBaseFiles))
	all = append(all, baseFiles...)
	for _, p := range localeBaseFiles {
		all = append(all, strings.ReplaceAll(p, "{locale}", locale))
	}

	crossDrive := sameVolume(baseInstall, profileRoot) == false
	if crossDrive {
		events.Emit(ctx, events.StatusMessage,
			"Base install on different drive — copying ~5GB. This is a one-time step and may take several minutes.")
	}

	for i, rel := range all {
		src := filepath.Join(baseInstall, filepath.FromSlash(rel))
		dst := filepath.Join(profileRoot, filepath.FromSlash(rel))

		srcInfo, err := os.Stat(src)
		if err != nil {
			if os.IsNotExist(err) {
				// Optional file absent in source install — fine, skip.
				continue
			}
			return fmt.Errorf("stat %s: %w", src, err)
		}
		if srcInfo.IsDir() {
			return fmt.Errorf("%s is a directory, expected file", src)
		}

		if _, err := os.Stat(dst); err == nil {
			// Already materialized — skip (idempotent re-run).
			continue
		}

		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(dst), err)
		}

		events.Emit(ctx, events.StatusMessage,
			fmt.Sprintf("Preparing base files %d/%d — %s", i+1, len(all), filepath.Base(rel)))

		if err := os.Link(src, dst); err == nil {
			continue
		} else if !isCrossDeviceErr(err) {
			// Real failure (permission denied, etc.) — surface it.
			// Note: we fall through to copy on cross-device. For other
			// link failures we try copy too because some filesystems
			// (FAT32, exFAT) silently disallow hardlinks; copy is a safe
			// universal fallback.
			if !isLinkUnsupportedErr(err) {
				return fmt.Errorf("hardlink %s -> %s: %w", src, dst, err)
			}
		}

		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("copy %s -> %s: %w", src, dst, err)
		}
	}
	return nil
}

// sameVolume reports whether two paths live on the same Windows volume by
// comparing their drive letters (case-insensitive). Falls back to true on
// non-Windows or unparseable paths so we don't emit spurious warnings.
func sameVolume(a, b string) bool {
	va := filepath.VolumeName(a)
	vb := filepath.VolumeName(b)
	if va == "" || vb == "" {
		return true
	}
	return strings.EqualFold(va, vb)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		_ = os.Remove(dst)
		return err
	}
	if err := out.Close(); err != nil {
		_ = os.Remove(dst)
		return err
	}
	return nil
}

// isCrossDeviceErr identifies the "source and dest on different volumes"
// link error across platforms.
func isCrossDeviceErr(err error) bool {
	if errors.Is(err, syscall.EXDEV) {
		return true
	}
	// Windows ERROR_NOT_SAME_DEVICE is 17. Go's syscall.EXDEV on Windows is
	// 18, so errors.Is above misses it — match the errno directly too.
	if errors.Is(err, syscall.Errno(17)) {
		return true
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "different drive") ||
		strings.Contains(s, "different disk drive") ||
		strings.Contains(s, "not same device") ||
		strings.Contains(s, "cross-device")
}

// isLinkUnsupportedErr identifies filesystems that reject hardlinks entirely
// (FAT32, exFAT, some network mounts). Falls back to copy for these.
func isLinkUnsupportedErr(err error) bool {
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "not supported") ||
		strings.Contains(s, "operation not permitted") ||
		strings.Contains(s, "incorrect function")
}
