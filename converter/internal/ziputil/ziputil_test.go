package ziputil

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// TestSameAsArchiveZip reads every entry of the zip files of the test
// documents both ways and compares the bytes, with the pooled inflaters
// reused from one entry to the next.
func TestSameAsArchiveZip(t *testing.T) {
	var files []string
	for _, pattern := range []string{"../../gerber/testdata/*.zip", "../../kicad/testdata/*.zip", "../../pptx/testdata/*.pptx", "../../docx/testdata/*.docx", "../../xlsx/testdata/*.xlsx", "../../epub/testdata/*.epub", "../../musicxml/testdata/*.mxl", "../../sxf/testdata/*.p2z"} {
		m, _ := filepath.Glob(pattern)
		files = append(files, m...)
	}
	if len(files) == 0 {
		t.Fatal("no zip files")
	}
	entries := 0
	for _, name := range files {
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		want, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			continue // not a zip (an encrypted document)
		}
		got, err := NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		for i, f := range want.File {
			a := read(t, f)
			b := read(t, got.File[i])
			if !bytes.Equal(a, b) {
				t.Errorf("%s %s: %d bytes vs %d", name, f.Name, len(a), len(b))
			}
			entries++
		}
	}
	t.Logf("%d entries of %d files compared", entries, len(files))
}

func read(t *testing.T, f *zip.File) []byte {
	t.Helper()
	rc, err := f.Open()
	if err != nil {
		t.Fatalf("%s: %v", f.Name, err)
	}
	defer rc.Close()
	b, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("%s: %v", f.Name, err)
	}
	return b
}
