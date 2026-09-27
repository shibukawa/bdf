package converter

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestFileMap(t *testing.T) {
	m := FileMap{
		"top.kicad_sch":         []byte("top"),
		"sub/power.kicad_sch":   []byte("power"),
		"sub/deep/io.kicad_sch": []byte("io"),
	}
	if err := fstest.TestFS(m, "top.kicad_sch", "sub/power.kicad_sch", "sub/deep/io.kicad_sch"); err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{"missing", "sub/missing", "/top.kicad_sch", "../top.kicad_sch", "sub/"} {
		if _, err := m.Open(bad); err == nil {
			t.Errorf("%q opened", bad)
		}
	}
	// an empty map is an empty directory
	if err := fstest.TestFS(FileMap{}); err != nil {
		t.Fatal(err)
	}
}

func TestReadRef(t *testing.T) {
	o := &Options{Files: FileMap{
		"top.kicad_sch":         []byte("top"),
		"sub/power.kicad_sch":   []byte("power"),
		"sub/deep/io.kicad_sch": []byte("io"),
		"other/io.kicad_sch":    []byte("other io"),
		"lib/sheet.kicad_wks":   []byte("wks"),
	}}
	for _, tc := range []struct {
		from, rel, want, path string
	}{
		{"", "top.kicad_sch", "top", "top.kicad_sch"},
		{"top.kicad_sch", "sub/power.kicad_sch", "power", "sub/power.kicad_sch"},
		{"top.kicad_sch", `sub\power.kicad_sch`, "power", "sub/power.kicad_sch"},
		// relative to the file that refers
		{"sub/power.kicad_sch", "deep/io.kicad_sch", "io", "sub/deep/io.kicad_sch"},
		{"sub/deep/io.kicad_sch", "../power.kicad_sch", "power", "sub/power.kicad_sch"},
		{"sub/power.kicad_sch", "./deep/../power.kicad_sch", "power", "sub/power.kicad_sch"},
		// by base name when no path matches: the first in lexical order
		{"top.kicad_sch", "elsewhere/io.kicad_sch", "other io", "other/io.kicad_sch"},
		{"top.kicad_sch", "../../outside/sheet.kicad_wks", "wks", "lib/sheet.kicad_wks"},
	} {
		b, p, err := o.ReadRef(tc.from, tc.rel)
		if err != nil || string(b) != tc.want || p != tc.path {
			t.Errorf("ReadRef(%q, %q) = %q, %q, %v; want %q, %q", tc.from, tc.rel, b, p, err, tc.want, tc.path)
		}
	}
	for _, rel := range []string{"", "/etc/passwd", `C:\x.kicad_sch`, "c:/x", "missing.kicad_sch", "sub"} {
		if _, _, err := o.ReadRef("top.kicad_sch", rel); err == nil {
			t.Errorf("ReadRef(%q) read a file", rel)
		}
	}
	if _, _, err := o.ReadRef("", "missing.kicad_sch"); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("a missing file: %v", err)
	}
}

func TestReadRefDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, s := range map[string]string{"sub/a.kicad_sch": "a", "b.kicad_sch": "b in dir"} {
		if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(name)), []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// Files first, then Dir
	o := &Options{Dir: dir, Files: FileMap{"b.kicad_sch": []byte("b in files")}}
	for _, tc := range []struct{ from, rel, want, path string }{
		{"", "sub/a.kicad_sch", "a", "sub/a.kicad_sch"},
		{"sub/a.kicad_sch", "../b.kicad_sch", "b in files", "b.kicad_sch"},
		{"", "b.kicad_sch", "b in files", "b.kicad_sch"},
	} {
		b, p, err := o.ReadRef(tc.from, tc.rel)
		if err != nil || string(b) != tc.want || p != tc.path {
			t.Errorf("ReadRef(%q, %q) = %q, %q, %v; want %q, %q", tc.from, tc.rel, b, p, err, tc.want, tc.path)
		}
	}
	o.Files = nil
	if b, _, err := o.ReadRef("sub/a.kicad_sch", "../b.kicad_sch"); err != nil || string(b) != "b in dir" {
		t.Errorf("from Dir: %q, %v", b, err)
	}
	if _, _, err := (&Options{}).ReadRef("", "b.kicad_sch"); err == nil {
		t.Error("a file was read with neither Files nor Dir")
	}
}
