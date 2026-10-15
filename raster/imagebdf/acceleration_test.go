package imagebdf

import (
	stdcontext "context"
	"image"
	"testing"
)

// Separates completed-image reuse from text-mask reuse. The benchmark of
// repeated default Page calls now intentionally measures completed-image hits.
func BenchmarkRenderCaches(b *testing.B) {
	for _, name := range []string{"pdf/chrome-doc", "pptx/features", "drawio/showcase"} {
		doc := openDoc(b, "../../testdata/"+name+".bdf")
		for _, mode := range []struct {
			name         string
			images, text int
		}{{"disabled", -1, -1}, {"text-only", -1, 0}, {"completed", 0, 0}} {
			b.Run(name+"/"+mode.name, func(b *testing.B) {
				r := New(doc, &Options{NoSystemFonts: true, FontDirs: testFonts, ImageCacheBytes: mode.images, TextCacheBytes: mode.text})
				if _, e := r.PageDPI(doc.Views[0], 0, 300); e != nil {
					b.Fatal(e)
				}
				b.ReportAllocs()
				for b.Loop() {
					if _, e := r.PageDPI(doc.Views[0], 0, 300); e != nil {
						b.Fatal(e)
					}
				}
			})
		}
	}
}
func BenchmarkPages(b *testing.B) {
	doc := openDoc(b, "../../testdata/pdf/chrome-doc.bdf")
	v := doc.Views[0]
	for _, workers := range []struct {
		name string
		n    int
	}{{"serial", 1}, {"parallel", 4}} {
		b.Run(workers.name, func(b *testing.B) {
			r := New(doc, &Options{NoSystemFonts: true, FontDirs: testFonts, ImageCacheBytes: -1})
			yield := func(int, *image.RGBA) error { return nil }
			pages := []int{0, 0, 0, 0}
			if e := r.PagesDPI(stdcontext.Background(), v, pages, 300, workers.n, yield); e != nil {
				b.Fatal(e)
			}
			b.ReportAllocs()
			for b.Loop() {
				if e := r.PagesDPI(stdcontext.Background(), v, pages, 300, workers.n, yield); e != nil {
					b.Fatal(e)
				}
			}
		})
	}
}
