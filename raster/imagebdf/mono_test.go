package imagebdf

import (
	stdcontext "context"
	"image"
	"image/color"
	"testing"

	"github.com/shibukawa/bdf"
)

// TestMonoImages checks the gray and bilevel images against the RGBA one.
func TestMonoImages(t *testing.T) {
	for _, name := range []string{"pptx/features", "pdf/chrome-doc", "demo"} {
		doc := openDoc(t, "../../testdata/"+name+".bdf")
		for _, mode := range []AntialiasMode{AntialiasOn, AntialiasOff} {
			r := New(doc, &Options{FontDirs: testFonts, NoSystemFonts: true, Antialias: mode})
			for _, v := range doc.Views {
				var rgba *image.RGBA
				var gray, bilevel image.Image
				var err error
				if v.Kind == bdf.ViewSheet || v.Kind == bdf.ViewScroll {
					region := bdf.Rect{X: 10.5, Y: 20.25, W: 300, H: 200}
					rgba, err = r.Region(v, 0, region, 300, 200)
					if err == nil {
						gray, err = r.RegionMono(v, 0, region, 300, 200, Gray)
					}
					if err == nil {
						bilevel, err = r.RegionMono(v, 0, region, 300, 200, Bilevel)
					}
				} else {
					rgba, err = r.Page(v, 0, 1)
					if err == nil {
						gray, err = r.PageMono(v, 0, 1, Gray)
					}
					if err == nil {
						bilevel, err = r.PageMono(v, 0, 1, Bilevel)
					}
				}
				if err != nil {
					t.Fatal(name, v.ID, err)
				}
				g, ok := gray.(*image.Gray)
				if !ok || g.Rect != rgba.Rect {
					t.Fatalf("%s %s: gray is %T %v", name, v.ID, gray, gray.Bounds())
				}
				b, ok := bilevel.(*image.Paletted)
				if !ok || b.Rect != rgba.Rect || len(b.Palette) != 2 || b.Palette[0] != (color.Gray{0}) || b.Palette[1] != (color.Gray{255}) {
					t.Fatalf("%s %s: bilevel is %T %v", name, v.ID, bilevel, bilevel.Bounds())
				}
				black, white := 0, 0
				for i := 0; i < len(g.Pix); i++ {
					p := rgba.Pix[4*i : 4*i+4]
					// the luminance over white, from the bytes of the RGBA image
					// (which are rounded once more, so within a step)
					y := (0.299*float64(p[0]) + 0.587*float64(p[1]) + 0.114*float64(p[2]) + float64(255-p[3]))
					if d := float64(g.Pix[i]) - y; d > 1.01 || d < -1.01 {
						t.Fatalf("%s %s: gray pixel %d is %d, the RGBA image gives %g", name, v.ID, i, g.Pix[i], y)
					}
					want := uint8(0)
					if g.Pix[i] >= 128 {
						want = 1
					}
					if b.Pix[i] != want {
						t.Fatalf("%s %s: bilevel pixel %d is %d for gray %d", name, v.ID, i, b.Pix[i], g.Pix[i])
					}
					if b.Pix[i] == 0 {
						black++
					} else {
						white++
					}
				}
				if black == 0 || white == 0 {
					t.Errorf("%s %s: %d black and %d white pixels", name, v.ID, black, white)
				}
			}
		}
	}
}

// TestMonoCache checks that the images of each kind are kept apart and
// returned as copies.
func TestMonoCache(t *testing.T) {
	doc := openDoc(t, "../../testdata/pptx/features.bdf")
	v := doc.Views[0]
	r := New(doc, &Options{FontDirs: testFonts, NoSystemFonts: true})
	rgba, err := r.Page(v, 0, 0.5)
	if err != nil {
		t.Fatal(err)
	}
	gray, err := r.PageMono(v, 0, 0.5, Gray)
	if err != nil {
		t.Fatal(err)
	}
	bilevel, err := r.PageMono(v, 0, 0.5, Bilevel)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.images.items) != 3 {
		t.Fatalf("%d images kept, want one of each kind", len(r.images.items))
	}
	gray.(*image.Gray).Pix[0] ^= 0xff
	again, err := r.PageMono(v, 0, 0.5, Gray)
	if err != nil {
		t.Fatal(err)
	}
	if again.(*image.Gray).Pix[0] == gray.(*image.Gray).Pix[0] {
		t.Error("the gray image kept is not a copy")
	}
	if _, err := r.PageMono(v, 0, 0.5, Bilevel+1); err == nil {
		t.Error("an unknown monochrome format is accepted")
	}
	if err := r.PagesMono(stdcontext.Background(), v, []int{0}, 0.5, 1, Bilevel+1, func(int, image.Image) error { return nil }); err == nil {
		t.Error("an unknown monochrome format is accepted by PagesMono")
	}
	// the batch gives the images kept, of the kind asked for
	n := 0
	err = r.PagesMono(stdcontext.Background(), v, []int{0, 0}, 0.5, 2, Bilevel, func(page int, img image.Image) error {
		n++
		b, ok := img.(*image.Paletted)
		if !ok || b.Rect != rgba.Rect {
			t.Errorf("page %d: %T %v", page, img, img.Bounds())
		}
		if string(b.Pix) != string(bilevel.(*image.Paletted).Pix) {
			t.Errorf("page %d differs from the single drawing", page)
		}
		return nil
	})
	if err != nil || n != 2 {
		t.Fatal(n, err)
	}
	if _, ok := r.images.get(imageKey{}); ok {
		t.Fatal("an empty key is kept")
	}
}
