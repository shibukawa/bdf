package tiff

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// FuzzDecode feeds damaged files to the reader: it may fail, but must not
// panic. The seeds are the test files.
func FuzzDecode(f *testing.F) {
	files, _ := filepath.Glob("testdata/*.tif")
	more, _ := filepath.Glob("../../tiff/testdata/*.tif")
	for _, name := range append(files, more...) {
		if b, err := os.ReadFile(name); err == nil && len(b) < 64<<10 {
			f.Add(b)
		}
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		old := MaxDecodeBytes
		MaxDecodeBytes = 1 << 24
		defer func() { MaxDecodeBytes = old }()
		file, err := Open(bytes.NewReader(b), int64(len(b)))
		if err != nil {
			return
		}
		for _, d := range file.Pages() {
			d.Decode()
			d.JPEG()
			d.Strings(TagDocumentName)
			d.Resolution()
		}
	})
}
