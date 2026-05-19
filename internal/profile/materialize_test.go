package profile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func seedBaseInstall(t *testing.T, present []string) string {
	t.Helper()
	root := t.TempDir()
	for _, rel := range present {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("payload:"+rel), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestMaterializeBase_HappyPath(t *testing.T) {
	files := []string{
		"Wow.exe",
		"Data/common.MPQ",
		"Data/common-2.MPQ",
		"Data/expansion.MPQ",
		"Data/lichking.MPQ",
		"Data/patch.MPQ",
		"Data/patch-2.MPQ",
		"Data/patch-3.MPQ",
		"Data/enUS/locale-enUS.MPQ",
		"Data/enUS/expansion-locale-enUS.MPQ",
		"Data/enUS/lichking-locale-enUS.MPQ",
		"Data/enUS/patch-enUS.MPQ",
		"Data/enUS/patch-enUS-2.MPQ",
		"Data/enUS/patch-enUS-3.MPQ",
	}
	base := seedBaseInstall(t, files)
	dst := t.TempDir()

	if err := MaterializeBase(nil, base, dst, "enUS"); err != nil {
		t.Fatalf("materialize: %v", err)
	}

	for _, rel := range files {
		p := filepath.Join(dst, filepath.FromSlash(rel))
		got, err := os.ReadFile(p)
		if err != nil {
			t.Errorf("missing %s: %v", rel, err)
			continue
		}
		if !strings.Contains(string(got), rel) {
			t.Errorf("%s content mismatch: %q", rel, got)
		}
	}
}

func TestMaterializeBase_MissingOptionalSkipped(t *testing.T) {
	// patch-3 and patch-enUS-3 are common to be absent on older installs.
	files := []string{
		"Wow.exe",
		"Data/common.MPQ",
		"Data/common-2.MPQ",
		"Data/expansion.MPQ",
		"Data/lichking.MPQ",
		"Data/patch.MPQ",
		"Data/patch-2.MPQ",
		// patch-3.MPQ deliberately missing
		"Data/enUS/locale-enUS.MPQ",
		"Data/enUS/expansion-locale-enUS.MPQ",
		"Data/enUS/lichking-locale-enUS.MPQ",
		"Data/enUS/patch-enUS.MPQ",
		"Data/enUS/patch-enUS-2.MPQ",
		// patch-enUS-3.MPQ deliberately missing
	}
	base := seedBaseInstall(t, files)
	dst := t.TempDir()

	if err := MaterializeBase(nil, base, dst, "enUS"); err != nil {
		t.Fatalf("materialize: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dst, "Data", "patch-3.MPQ")); !os.IsNotExist(err) {
		t.Errorf("patch-3.MPQ should be absent in dst, err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(dst, "Data", "common.MPQ")); err != nil {
		t.Errorf("required common.MPQ missing in dst: %v", err)
	}
}

func TestMaterializeBase_Idempotent(t *testing.T) {
	files := []string{"Wow.exe", "Data/common.MPQ", "Data/enUS/locale-enUS.MPQ"}
	base := seedBaseInstall(t, files)
	dst := t.TempDir()

	if err := MaterializeBase(nil, base, dst, "enUS"); err != nil {
		t.Fatalf("first: %v", err)
	}
	// Mutate dst to detect if second run clobbers.
	sentinel := filepath.Join(dst, "Wow.exe")
	if err := os.WriteFile(sentinel, []byte("user-modified"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := MaterializeBase(nil, base, dst, "enUS"); err != nil {
		t.Fatalf("second: %v", err)
	}
	got, _ := os.ReadFile(sentinel)
	if string(got) != "user-modified" {
		t.Errorf("idempotent re-run clobbered existing file: %q", got)
	}
}

func TestMaterializeBase_EmptyArgsError(t *testing.T) {
	cases := [][3]string{
		{"", "/tmp/dst", "enUS"},
		{"/tmp/src", "", "enUS"},
		{"/tmp/src", "/tmp/dst", ""},
	}
	for _, c := range cases {
		if err := MaterializeBase(nil, c[0], c[1], c[2]); err == nil {
			t.Errorf("expected error for %v", c)
		}
	}
}
