package pdf

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/imgconv"
)

// pageImages decodes the images the first page draws, in drawing order.
func pageImages(t *testing.T, r *bdf.Reader) []image.Image {
	t.Helper()
	page, err := r.Object(r.Manifest.Views[0].Pages[0].Layers[0].Obj)
	if err != nil {
		t.Fatal(err)
	}
	var out []image.Image
	if err := page.Walk(func(in bdf.Instr) {
		if in.Op != bdf.OpImage {
			return
		}
		b, err := r.Part(page.Images[in.Args[0].(uint64)])
		if err != nil {
			t.Fatal(err)
		}
		img, err := imgconv.Decode(b)
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, img)
	}); err != nil {
		t.Fatal(err)
	}
	return out
}

func nrgba(img image.Image, x, y int) color.NRGBA {
	b := img.Bounds()
	return color.NRGBAModel.Convert(img.At(b.Min.X+x, b.Min.Y+y)).(color.NRGBA)
}

func TestConvertImages(t *testing.T) {
	res, r := convert(t, "images-jpx-jbig2.pdf")
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	imgs := pageImages(t, r)
	if len(imgs) != 5 {
		t.Fatalf("%d images", len(imgs))
	}
	for i, im := range imgs[:3] {
		if b := im.Bounds(); b.Dx() != 180 || b.Dy() != 120 {
			t.Errorf("image %d is %v", i, b)
		}
	}
	// The lossless grayscale one keeps the checker corner exactly.
	for _, p := range [][3]int{{2, 2, 18}, {10, 2, 255}, {10, 10, 18}} { // (20//3)·3: the fixture averages by thirds
		if g := nrgba(imgs[1], p[0], p[1]); g.R != uint8(p[2]) || g.G != g.R || g.A != 255 {
			t.Errorf("gray (%d,%d) = %v, want %d", p[0], p[1], g, p[2])
		}
	}
	// The lossy RGB one: red grows to the right, green downwards.
	if a, b := nrgba(imgs[0], 60, 60), nrgba(imgs[0], 170, 110); !(int(b.R) > int(a.R)+60 && int(b.G) > int(a.G)+40) {
		t.Errorf("rgb %v %v", a, b)
	}
	// The alpha of the RGBA one is its soft mask: opaque near the rings'
	// centre, transparent in the far corner.
	if c, e := nrgba(imgs[2], 108, 54), nrgba(imgs[2], 0, 119); c.A < 200 || e.A > 40 {
		t.Errorf("alpha centre %v corner %v", c, e)
	}

	// JBIG2: the halftone image is black where the bitmap is 1; the text is
	// a stencil mask painted with the fill colour where it is 1.
	for i, base := range []string{"halftone_arith", "globals"} {
		want := readPBM(t, filepath.Join("..", "internal", "jbig2", "testdata", base+".pbm"))
		img := imgs[3+i]
		if b := img.Bounds(); b.Dx() != want.w || b.Dy() != want.h {
			t.Fatalf("%s is %v, want %dx%d", base, b, want.w, want.h)
		}
		bad := 0
		for y := 0; y < want.h; y++ {
			for x := 0; x < want.w; x++ {
				c := nrgba(img, x, y)
				black := want.bits[y*((want.w+7)/8)+x/8]>>(7-x%8)&1 == 1
				var ok bool
				if i == 0 {
					ok = black == (c.R == 0 && c.A == 255) && (black || c.R == 255)
				} else {
					ok = black == (c.A == 255) && (!black || c.R > 100 && c.G < 60)
				}
				if !ok {
					bad++
				}
			}
		}
		if bad != 0 {
			t.Errorf("%s: %d pixels differ", base, bad)
		}
	}
}

type pbm struct {
	w, h int
	bits []byte
}

func readPBM(t *testing.T, path string) pbm {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var p pbm
	lines := bytes.SplitN(b, []byte("\n"), 3)
	if len(lines) != 3 || string(lines[0]) != "P4" {
		t.Fatalf("%s: not a PBM", path)
	}
	if _, err := fmt.Sscanf(string(lines[1]), "%d %d", &p.w, &p.h); err != nil {
		t.Fatal(err)
	}
	p.bits = lines[2]
	return p
}
