package install

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ErrNotWowInstall is returned by Validate when path is not a WoW install.
var ErrNotWowInstall = errors.New("not a valid WoW 3.3.5 install")

// Install describes a verified WoW client install on disk.
type Install struct {
	Root    string // absolute path to install root
	Locale  string // e.g. "enUS" (detected from Data/<locale>/)
	WowExe  string // absolute path to Wow.exe
	Version string // best-effort, may be empty
}

// Validate checks whether path looks like a WoW 3.3.5 install and returns
// metadata. The check is structural: presence of Wow.exe + Data/common.MPQ +
// a known locale subfolder. It does NOT enforce that the client is exactly
// 3.3.5a — that's the manifest's job (client_version field).
func Validate(path string) (*Install, error) {
	wowExe := filepath.Join(path, "Wow.exe")
	if _, err := os.Stat(wowExe); err != nil {
		return nil, fmt.Errorf("%w: missing Wow.exe", ErrNotWowInstall)
	}
	dataDir := filepath.Join(path, "Data")
	if _, err := os.Stat(filepath.Join(dataDir, "common.MPQ")); err != nil {
		return nil, fmt.Errorf("%w: missing Data/common.MPQ", ErrNotWowInstall)
	}
	locale, err := detectLocale(dataDir)
	if err != nil {
		return nil, err
	}
	return &Install{Root: path, Locale: locale, WowExe: wowExe}, nil
}

func detectLocale(dataDir string) (string, error) {
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return "", fmt.Errorf("read Data dir: %w", err)
	}
	knownLocales := map[string]bool{
		"enUS": true, "enGB": true, "deDE": true, "frFR": true,
		"esES": true, "esMX": true, "ruRU": true, "koKR": true,
		"zhCN": true, "zhTW": true, "ptBR": true, "itIT": true,
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if knownLocales[e.Name()] {
			// Confirm by presence of locale-{loc}.MPQ
			if _, err := os.Stat(filepath.Join(dataDir, e.Name(), "locale-"+e.Name()+".MPQ")); err == nil {
				return e.Name(), nil
			}
		}
	}
	return "", fmt.Errorf("%w: no locale folder found in Data/", ErrNotWowInstall)
}

// AutoDetect scans common install locations for a WoW 3.3.5 client.
// Returns all hits — caller picks or lets user pick.
func AutoDetect() []*Install {
	candidates := commonPaths()
	var found []*Install
	seen := map[string]bool{}
	for _, p := range candidates {
		abs, err := filepath.Abs(p)
		if err != nil || seen[strings.ToLower(abs)] {
			continue
		}
		seen[strings.ToLower(abs)] = true
		if inst, err := Validate(abs); err == nil {
			found = append(found, inst)
		}
	}
	return found
}

func commonPaths() []string {
	var paths []string
	if runtime.GOOS == "windows" {
		drives := []string{"C:", "D:", "E:", "F:"}
		stems := []string{
			`\Program Files\World of Warcraft`,
			`\Program Files (x86)\World of Warcraft`,
			`\Games\World of Warcraft`,
			`\World of Warcraft`,
			`\WoW`,
			`\WoW 3.3.5a`,
			`\World of Warcraft 3.3.5a`,
		}
		for _, d := range drives {
			for _, s := range stems {
				paths = append(paths, d+s)
			}
		}
		if user := os.Getenv("USERPROFILE"); user != "" {
			paths = append(paths,
				filepath.Join(user, "Desktop", "World of Warcraft"),
				filepath.Join(user, "Documents", "World of Warcraft"),
			)
		}
	}
	return paths
}
