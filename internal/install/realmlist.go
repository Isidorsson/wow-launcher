package install

import (
	"fmt"
	"os"
	"path/filepath"
)

// realmlist.wtf format:
//   set realmlist logon.example.com
//   set patchlist logon.example.com   (optional)
//   set portal "us"                   (optional, region)
//
// The file lives at <Data>/<Locale>/realmlist.wtf. The launcher MUST overwrite
// it before launching Wow.exe so the client points at the private server's
// auth/realm gateway instead of Blizzard's.

// WriteRealmlist sets the realm address for an install.
//
// ★ LEARNING GAP — implement this yourself.
//
// Why this matters: realmlist.wtf is what redirects the client to your
// private server. Get it wrong and the client either silently connects to
// the wrong endpoint or crashes on bnet handshake. The function should:
//
//   1. Compute the target path: <inst.Root>/Data/<inst.Locale>/realmlist.wtf
//   2. Build file body. Conventional format (each on its own line):
//        set realmlist <addr>
//      Some servers also want:
//        set patchlist <addr>
//      Decide whether you write patchlist too (most private servers don't
//      need it; including it when not configured = silent breakage).
//   3. Write atomically: write to realmlist.wtf.tmp then rename. Reason:
//      if the launcher is killed mid-write, you don't want a half-written
//      file that crashes the client.
//   4. CRLF line endings on Windows builds (the WoW client tolerates both
//      but be consistent — easier debugging).
//   5. Return descriptive errors with %w so callers can branch on os.IsPermission.
//
// Tests you should write afterwards:
//   - file is written with correct content
//   - existing realmlist.wtf is replaced (not appended)
//   - permission errors are wrapped, not swallowed
//   - empty addr returns error before touching disk
func WriteRealmlist(inst *Install, addr string) error {
	// TODO: implement
	_ = fmt.Sprintf
	_ = os.Rename
	_ = filepath.Join
	return fmt.Errorf("WriteRealmlist not implemented yet — see comment in realmlist.go")
}
