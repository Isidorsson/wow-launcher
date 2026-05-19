package install

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteRealmlist(t *testing.T) {
	tests := []struct {
		name    string
		inst    *Install
		addr    string
		wantErr error
		check   func(t *testing.T, root string)
	}{
		{
			name: "happy path writes CRLF single line",
			inst: nil, // set per-test below
			addr: "logon.example.com",
			check: func(t *testing.T, root string) {
				got, err := os.ReadFile(filepath.Join(root, "Data", "enUS", "realmlist.wtf"))
				if err != nil {
					t.Fatalf("read: %v", err)
				}
				want := "set realmlist logon.example.com\r\n"
				if string(got) != want {
					t.Errorf("body = %q, want %q", got, want)
				}
			},
		},
		{
			name: "replaces existing realmlist",
			addr: "new.example.com",
			check: func(t *testing.T, root string) {
				p := filepath.Join(root, "Data", "enUS", "realmlist.wtf")
				got, _ := os.ReadFile(p)
				if !strings.Contains(string(got), "new.example.com") {
					t.Errorf("expected new addr, got %q", got)
				}
				if strings.Contains(string(got), "OLD") {
					t.Errorf("old content survived: %q", got)
				}
			},
		},
		{
			name:    "empty addr errors before touching disk",
			addr:    "",
			wantErr: ErrEmptyRealmlistAddr,
		},
		{
			name:    "nil install errors",
			addr:    "x.example.com",
			wantErr: errors.New("install missing"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			if tt.inst == nil && tt.name != "nil install errors" {
				tt.inst = &Install{Root: root, Locale: "enUS"}
			}
			if tt.name == "replaces existing realmlist" {
				// Pre-seed the file before the call.
				dir := filepath.Join(root, "Data", "enUS")
				_ = os.MkdirAll(dir, 0o755)
				_ = os.WriteFile(filepath.Join(dir, "realmlist.wtf"), []byte("set realmlist OLD\r\n"), 0o644)
			}

			err := WriteRealmlist(tt.inst, tt.addr)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if errors.Is(tt.wantErr, ErrEmptyRealmlistAddr) && !errors.Is(err, ErrEmptyRealmlistAddr) {
					t.Errorf("err = %v, want ErrEmptyRealmlistAddr", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.check != nil {
				tt.check(t, root)
			}
		})
	}
}

func TestWriteRealmlist_NoTmpLeftover(t *testing.T) {
	root := t.TempDir()
	inst := &Install{Root: root, Locale: "enUS"}
	if err := WriteRealmlist(inst, "logon.example.com"); err != nil {
		t.Fatalf("write: %v", err)
	}
	entries, _ := os.ReadDir(filepath.Join(root, "Data", "enUS"))
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".tmp") {
			t.Errorf("found leftover tmp file: %s", e.Name())
		}
	}
}
