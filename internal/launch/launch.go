package launch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Run spawns Wow.exe from the given profile root and detaches.
// cwd is set to the profile root because the client resolves Data/ relatively.
// Returns once the process is started; does not wait for exit.
func Run(profileRoot string) (*os.Process, error) {
	wowExe := filepath.Join(profileRoot, "Wow.exe")
	if _, err := os.Stat(wowExe); err != nil {
		return nil, fmt.Errorf("Wow.exe not found at %s: %w", wowExe, err)
	}
	cmd := exec.Command(wowExe)
	cmd.Dir = profileRoot
	// On Windows, exec.Command already detaches stdio sensibly for GUI apps.
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start Wow.exe: %w", err)
	}
	return cmd.Process, nil
}
