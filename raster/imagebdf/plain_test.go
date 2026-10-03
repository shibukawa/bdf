package imagebdf_test

import (
	"bytes"
	"image"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/raster/imagebdf"
	"github.com/shibukawa/bdf/thumbnail"
)

// TestShortcutsChangeNoPixel draws the documents of testdata, a thumbnail
// and pages of each, with the shortcuts of the renderer and without them
// (what lies outside the canvas drawn like the rest, the crossings of a
// scanline sorted by insertion, every pixel asked for its coverage, the
// colours of a gradient worked out for each fill, no memory kept from one
// drawing to the next), and checks that the images are the same.
func TestShortcutsChangeNoPixel(t *testing.T) {
	if testing.Short() {
		t.Skip("slow")
	}
	files, _ := filepath.Glob("../../testdata/*/*.bdf")
	if len(files) < 50 {
		t.Fatalf("%d documents in testdata", len(files))
	}
	files = append(files, "../../testdata/demo.bdf")
	defer imagebdf.SetPlain(false)
	fonts := imagebdf.Options{FontDirs: []string{"../../converter/pptx/testdata/fonts", "../../converter/docx/testdata/fonts", "../../fixture/testdata/fonts"}, NoSystemFonts: true}
	for _, f := range files {
		rd, err := bdf.OpenSingleFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if rd.Encrypted() {
			// the segments of testdata/segments, sealed for a key the tests of @bdfkit/core hold
			rd.Close()
			continue
		}
		d, err := rd.ToDocument()
		rd.Close()
		if err != nil {
			t.Fatal(f, err)
		}
		// draw returns the images of the document and its warnings
		draw := func(plain bool) (images []*image.RGBA, warnings []string) {
			imagebdf.SetPlain(plain)
			add := func(img *image.RGBA, err error) {
				if err != nil {
					t.Fatal(f, err)
				}
				images = append(images, img)
			}
			res, err := thumbnail.Make(d, &thumbnail.Options{Raster: fonts})
			if err != nil {
				t.Fatal(f, err)
			}
			add(res.Image, nil)
			warnings = res.Warnings
			for _, v := range d.Views {
				r := imagebdf.New(d, &fonts)
				switch v.Kind {
				case bdf.ViewSheet:
					add(r.Region(v, 0, bdf.Rect{X: 300.5, Y: 200.25, W: 500, H: 400}, 375, 300))
				case bdf.ViewScroll:
					add(r.Region(v, 0, bdf.Rect{X: 100, Y: 300, W: 400, H: 400}, 300, 300))
				default:
					if len(v.Pages) == 0 {
						continue
					}
					add(r.Page(v, 0, 1))
					// a part of the page, the rest of which lies outside
					p := v.Pages[len(v.Pages)-1]
					add(r.Region(v, len(v.Pages)-1, bdf.Rect{X: p.W / 3, Y: p.H / 4, W: p.W / 2, H: p.H / 2}, 200, 200))
				}
				warnings = append(warnings, r.Warnings()...)
			}
			return images, warnings
		}
		fast, fw := draw(false)
		plain, pw := draw(true)
		for i := range fast {
			if fast[i].Rect != plain[i].Rect || !bytes.Equal(fast[i].Pix, plain[i].Pix) {
				n := 0
				for k := range fast[i].Pix {
					if k < len(plain[i].Pix) && fast[i].Pix[k] != plain[i].Pix[k] {
						n++
					}
				}
				t.Errorf("%s: image %d (%v) differs in %d bytes", f, i, fast[i].Rect.Size(), n)
			}
		}
		// the warnings are the same but for those of images that show
		// nowhere, which are not decoded
		for _, w := range fw {
			found := false
			for _, p := range pw {
				found = found || p == w
			}
			if !found {
				t.Errorf("%s: warning %q only with the shortcuts", f, w)
			}
		}
	}
}
