package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// State persists per-server manifest metadata between launches so the launcher
// can do conditional fetches (If-None-Match) and detect content changes without
// re-downloading the full manifest each time.
//
// Stored at <UserConfigDir>/WowLauncher/manifest-state.json. Best-effort:
// corrupt or missing state is treated as "no prior knowledge" and the next
// fetch becomes unconditional.
type State struct {
	Servers map[string]ServerState `json:"servers"`
}

type ServerState struct {
	ETag         string `json:"etag"`          // value from previous Response, sent as If-None-Match
	ManifestHash string `json:"manifestHash"`  // sha256 of canonical inner manifest JSON
	LastChecked  int64  `json:"lastChecked"`   // unix seconds
}

// StateStore wraps an on-disk State with a mutex so concurrent fetches
// (multiple servers in parallel) don't corrupt the file.
type StateStore struct {
	path string
	mu   sync.Mutex
	data State
}

func OpenStateStore(dir string) (*StateStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir state dir: %w", err)
	}
	s := &StateStore{path: filepath.Join(dir, "manifest-state.json")}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *StateStore) load() error {
	raw, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			s.data = State{Servers: map[string]ServerState{}}
			return nil
		}
		return fmt.Errorf("read state: %w", err)
	}
	if err := json.Unmarshal(raw, &s.data); err != nil {
		// Corrupt file — start fresh rather than refuse to launch.
		s.data = State{Servers: map[string]ServerState{}}
		return nil
	}
	if s.data.Servers == nil {
		s.data.Servers = map[string]ServerState{}
	}
	return nil
}

func (s *StateStore) Get(serverID string) ServerState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.data.Servers[serverID]
}

func (s *StateStore) Set(serverID string, ss ServerState) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Servers[serverID] = ss
	raw, err := json.MarshalIndent(&s.data, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
