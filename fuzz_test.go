package bdf_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/raster/imagebdf"
	"github.com/shibukawa/bdf/thumbnail"
)

// FuzzReader reads a document as a server does: the container, every part,
// the text of every view, the first pages drawn, and the document written
// again. Malformed input is an error, never a panic, a hang or memory out
// of proportion with the input.
func FuzzReader(f *testing.F) {
	for _, pattern := range []string{"testdata/*.bdf", "testdata/*/*.bdf"} {
		files, _ := filepath.Glob(pattern)
		for _, name := range files {
			if b, err := os.ReadFile(name); err == nil {
				f.Add(b)
			}
		}
	}
	raster := imagebdf.Options{FontDirs: []string{"converter/pptx/testdata/fonts", "fixture/testdata/fonts"}, NoSystemFonts: true}
	f.Fuzz(func(t *testing.T, data []byte) {
		r, err := bdf.OpenSingle(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return
		}
		if r.Locked() {
			_ = r.Unlock("demo-パスワード")
			if r.Locked() {
				_ = r.WriteSingle(io.Discard) // the stored parts, copied
				return
			}
		}
		doc, err := r.ToDocument()
		if err != nil {
			return
		}
		if _, err := doc.SearchText(); err != nil {
			t.Logf("search text: %v", err)
		}
		if _, err := thumbnail.Make(doc, &thumbnail.Options{Size: 48, Raster: raster}); err != nil {
			t.Logf("thumbnail: %v", err)
		}
		ren := imagebdf.New(doc, &raster)
		for _, v := range doc.Views {
			if v.Kind == bdf.ViewSheet {
				continue
			}
			for i := range min(len(v.Pages), 2) {
				if _, err := ren.Page(v, i, 0.1); err != nil {
					t.Logf("page %d: %v", i, err)
				}
			}
			if len(v.Pages) > 0 && v.Kind != bdf.ViewSheet && v.Play == nil {
				_ = r.WriteSegment(io.Discard, bdf.SegmentOptions{Segment: bdf.SegmentAt(v, 0, 3)})
			}
		}
		if err := doc.WriteSingle(io.Discard); err != nil {
			t.Fatalf("write: %v", err)
		}
	})
}
