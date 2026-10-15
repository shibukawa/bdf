package imagebdf

import (
	stdcontext "context"
	"fmt"
	"image"
	"testing"
)

// BenchmarkRenderHiDPI draws pages at print resolutions: 600 dpi without
// antialiasing, and 360 dpi with and without it.
func BenchmarkRenderHiDPI(b *testing.B) {
	for _, name := range []string{"pdf/chrome-doc", "docx/basic", "pptx/features"} {
		doc := openDoc(b, "../../testdata/"+name+".bdf")
		view := doc.Views[0]
		for _, c := range []struct {
			dpi  float64
			mode AntialiasMode
			tag  string
			mono Mono
		}{{600, AntialiasOff, "600dpi/off", 0}, {600, AntialiasOff, "600dpi/off/gray", Gray}, {600, AntialiasOff, "600dpi/off/bilevel", Bilevel}, {360, AntialiasOn, "360dpi/on", 0}, {360, AntialiasOff, "360dpi/off", 0}} {
			b.Run(fmt.Sprintf("%s/%s", name, c.tag), func(b *testing.B) {
				r := New(doc, &Options{ImageCacheBytes: -1, TextCacheBytes: -1, Antialias: c.mode, FontDirs: testFonts, NoSystemFonts: true})
				if _, err := r.PageMono(view, 0, c.dpi/72, c.mono); err != nil {
					b.Fatal(err)
				}
				b.ReportAllocs()
				for b.Loop() {
					if _, err := r.PageMono(view, 0, c.dpi/72, c.mono); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

// BenchmarkPagesHiDPI draws a batch of pages at 600 dpi without
// antialiasing, as a print job does: in colour, in gray and in black and
// white.
func BenchmarkPagesHiDPI(b *testing.B) {
	doc := openDoc(b, "../../testdata/pdf/chrome-doc.bdf")
	v := doc.Views[0]
	for _, c := range []struct {
		name string
		mono Mono
	}{{"rgba", 0}, {"gray", Gray}, {"bilevel", Bilevel}} {
		b.Run(c.name, func(b *testing.B) {
			r := New(doc, &Options{ImageCacheBytes: -1, Antialias: AntialiasOff, FontDirs: testFonts, NoSystemFonts: true})
			yield := func(int, image.Image) error { return nil }
			pages := []int{0, 0, 0, 0}
			if err := r.PagesMono(stdcontext.Background(), v, pages, 600.0/72, 1, c.mono, yield); err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			for b.Loop() {
				if err := r.PagesMono(stdcontext.Background(), v, pages, 600.0/72, 1, c.mono, yield); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
