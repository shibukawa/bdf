package otf

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"testing/fstest"
)

// sample is a font collection of two fonts.
const sample = "../../converter/font/testdata/bold.ttc"

// canMap reports whether the platform maps files.
func canMap() bool {
	switch runtime.GOOS {
	case "linux", "darwin", "windows", "freebsd", "netbsd", "openbsd":
		return true
	}
	return false
}

func TestOpen(t *testing.T) {
	want, err := os.ReadFile(sample)
	if err != nil {
		t.Fatal(err)
	}
	f, err := Open(sample)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if !bytes.Equal(f.Bytes(), want) {
		t.Fatal("the contents differ from ReadFile's")
	}
	if f.Size() != len(want) {
		t.Fatalf("Size = %d, want %d", f.Size(), len(want))
	}
	if string(f.Bytes()[:4]) != "ttcf" {
		t.Fatalf("tag %q, want ttcf", f.Bytes()[:4])
	}
	if canMap() && !f.Mapped() {
		t.Error("the file is read, not mapped")
	}
}

func TestOpenFS(t *testing.T) {
	want, err := os.ReadFile(sample)
	if err != nil {
		t.Fatal(err)
	}
	dir, name := filepath.Split(sample)
	f, err := OpenFS(os.DirFS(dir), name)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if !bytes.Equal(f.Bytes(), want) {
		t.Fatal("DirFS: the contents differ from ReadFile's")
	}
	if canMap() && !f.Mapped() {
		t.Error("DirFS: the file is read, not mapped")
	}
	m, err := OpenFS(fstest.MapFS{"x.ttc": {Data: want}}, "x.ttc")
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()
	if !bytes.Equal(m.Bytes(), want) {
		t.Fatal("MapFS: the contents differ")
	}
	if m.Mapped() {
		t.Error("MapFS: a file that is not an *os.File is mapped")
	}
	if _, err := OpenFS(fstest.MapFS{}, "missing.ttf"); err == nil {
		t.Error("a missing file opens")
	}
}

func TestClose(t *testing.T) {
	f, err := Open(sample)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if f.Bytes() != nil || f.Mapped() || f.Size() != 0 {
		t.Error("a closed file still has contents")
	}
	if err := f.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}
}

func TestEmptyAndMissing(t *testing.T) {
	dir := t.TempDir()
	empty := filepath.Join(dir, "empty.ttf")
	if err := os.WriteFile(empty, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	f, err := Open(empty)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if f.Size() != 0 || f.Mapped() {
		t.Errorf("empty file: size %d, mapped %v", f.Size(), f.Mapped())
	}
	if _, err := Open(filepath.Join(dir, "missing.ttf")); err == nil {
		t.Error("a missing file opens")
	}
	if _, err := Open(dir); err == nil {
		t.Error("a directory opens")
	}
}
