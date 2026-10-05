package imagebdf

import (
	"testing"

	"github.com/shibukawa/bdf"
)

func BenchmarkRenderPage(b *testing.B) {
	for _, name := range []string{"pdf/chrome-doc", "pptx/features", "drawio/showcase"} {
		doc := openDoc(b, "../../testdata/"+name+".bdf")
		view := doc.Views[0]
		if view.Kind == bdf.ViewSheet || len(view.Pages) == 0 {
			b.Fatalf("%s has no fixed page", name)
		}
		for _, scale := range []struct {
			name string
			x    float64
		}{{"thumbnail", 0.25}, {"page", 1}} {
			b.Run(name+"/"+scale.name, func(b *testing.B) {
				r := New(doc, &Options{FontDirs: testFonts, NoSystemFonts: true})
				if _, err := r.Page(view, 0, scale.x); err != nil {
					b.Fatal(err)
				}
				b.ReportAllocs()
				for b.Loop() {
					if _, err := r.Page(view, 0, scale.x); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}
