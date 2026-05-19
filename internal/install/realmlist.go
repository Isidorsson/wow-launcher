package install

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrEmptyRealmlistAddr is returned when caller passes an empty address.
var ErrEmptyRealmlistAddr = errors.New("realmlist address is empty")

// WriteRealmlist atomically writes realmlist.wtf in the install's locale dir.
//
// Format (CRLF, single line, no trailing whitespace):
//
//   set realmlist <addr>
//
// We deliberately do NOT write `set patchlist` — most 3.3.5 private servers
// don't operate a Blizzard-style patch server, and writing a stale/wrong
// patchlist makes the client try to contact it on every launch (slow + can
// hang the login screen). If a downstream server needs patchlist, extend the
// manifest with a separate field and add a Write* helper for it.
//
// Atomicity: write to realmlist.wtf.tmp then os.Rename to final path. On
// Windows os.Rename is atomic when source/dest are on the same volume, which
// they always are here. Crash mid-write leaves only the .tmp file behind.
func WriteRealmlist(inst *Install, addr string) error {
	if addr == "" {
		return ErrEmptyRealmlistAddr
	}
	if inst == nil || inst.Root == "" || inst.Locale == "" {
		return fmt.Errorf("install missing Root or Locale")
	}

	target := filepath.Join(inst.Root, "Data", inst.Locale, "realmlist.wtf")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("mkdir locale dir: %w", err)
	}

	body := []byte("set realmlist " + addr + "\r\n")
	tmp := target + ".tmp"
	if err := os.WriteFile(tmp, body, 0o644); err != nil {
		return fmt.Errorf("write tmp realmlist: %w", err)
	}
	if err := os.Rename(tmp, target); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename realmlist: %w", err)
	}
	return nil
}
