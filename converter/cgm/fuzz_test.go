package cgm

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
)

// FuzzInterp feeds damaged metafiles to both readers and the interpreter,
// whatever their first element: they may draw nothing, but must not panic
// or hang. The seeds are the test files.
func FuzzInterp(f *testing.F) {
	files, _ := filepath.Glob("testdata/*.cgm")
	for _, name := range files {
		if b, err := os.ReadFile(name); err == nil {
			f.Add(b)
		}
	}
	fonts := &cad.Fonts{Set: fontset.New(fontdb.New(nil, []string{"../pptx/testdata/fonts"}, false), func(string) {})}
	// smaller limits: an input of a few bytes may ask for the most
	maxPoints, maxTotalCells, maxCells, maxPictures = 1<<20, 1<<18, 1<<16, 100
	f.Fuzz(func(t *testing.T, b []byte) {
		for _, text := range []bool{false, true} {
			var elems []element
			var next func() (element, bool)
			if text {
				next = (&textReader{s: string(b)}).next
			} else {
				next = (&binReader{data: b}).next
			}
			for {
				e, ok := next()
				if !ok {
					break
				}
				elems = append(elems, e)
			}
			c := &converter{opts: &Options{}, doc: bdf.NewDocument(), warned: map[string]bool{}}
			in := newInterp(c, text)
			in.fonts = fonts
			in.run(iterate(elems))
		}
	})
}
