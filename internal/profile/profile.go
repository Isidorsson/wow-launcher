package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Profile is one per-server install slot.
//
// Strategy: each server gets its own Root directory. The base client files
// (Wow.exe, common.MPQ, expansion.MPQ, lichking.MPQ, locale base) are either
// copied or hard-linked from the user's detected base install on first use;
// custom server patches go into Data/ and Data/<locale>/ on top.
//
// Hard-linking saves ~5GB per profile. Not all filesystems support it
// (FAT32 doesn't, NTFS does). We try hardlink first, fall back to copy.
type Profile struct {
	ServerID    string `json:"server_id"`
	Root        string `json:"root"`         // absolute path to this profile's install
	Locale      string `json:"locale"`       // copied from source install
	BaseInstall string `json:"base_install"` // path to user's base 3.3.5 install (source of base MPQs)
}

// Manager owns profile.json on disk and provides CRUD over profiles.
type Manager struct {
	root      string // directory containing profile.json + <serverID>/ subdirs
	profiles  map[string]*Profile
}

func NewManager(root string) (*Manager, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir profiles root: %w", err)
	}
	m := &Manager{root: root, profiles: map[string]*Profile{}}
	if err := m.load(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Manager) indexPath() string { return filepath.Join(m.root, "profiles.json") }

func (m *Manager) load() error {
	data, err := os.ReadFile(m.indexPath())
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	var list []*Profile
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("parse profiles.json: %w", err)
	}
	for _, p := range list {
		m.profiles[p.ServerID] = p
	}
	return nil
}

func (m *Manager) save() error {
	var list []*Profile
	for _, p := range m.profiles {
		list = append(list, p)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	tmp := m.indexPath() + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, m.indexPath())
}

func (m *Manager) Get(serverID string) (*Profile, bool) {
	p, ok := m.profiles[serverID]
	return p, ok
}

func (m *Manager) List() []*Profile {
	out := make([]*Profile, 0, len(m.profiles))
	for _, p := range m.profiles {
		out = append(out, p)
	}
	return out
}

// Create initializes a new profile by materializing base client files from
// baseInstall into <profiles_root>/<serverID>/. See MaterializeBase for the
// hardlink/copy logic — that's the learning gap.
func (m *Manager) Create(serverID, baseInstall, locale string) (*Profile, error) {
	if _, exists := m.profiles[serverID]; exists {
		return nil, fmt.Errorf("profile %s already exists", serverID)
	}
	profileRoot := filepath.Join(m.root, serverID)
	if err := os.MkdirAll(profileRoot, 0o755); err != nil {
		return nil, err
	}
	p := &Profile{
		ServerID:    serverID,
		Root:        profileRoot,
		Locale:      locale,
		BaseInstall: baseInstall,
	}
	if err := MaterializeBase(baseInstall, profileRoot, locale); err != nil {
		return nil, fmt.Errorf("materialize base: %w", err)
	}
	m.profiles[serverID] = p
	if err := m.save(); err != nil {
		return nil, err
	}
	return p, nil
}

func (m *Manager) Delete(serverID string) error {
	p, ok := m.profiles[serverID]
	if !ok {
		return nil
	}
	if err := os.RemoveAll(p.Root); err != nil {
		return err
	}
	delete(m.profiles, serverID)
	return m.save()
}
