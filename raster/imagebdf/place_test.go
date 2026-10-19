package imagebdf

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"testing"

	"github.com/shibukawa/bdf"
)

// placedRows draws a page with PlaceRows and puts the bands together.
func placedRows(t *testing.T, r *Renderer, v *bdf.View, pl Place, mono Mono) image.Image {
	t.Helper()
	var out draw.Image
	next := 0
	err := r.PlaceRows(v, 0, pl, mono, func(y0 int, rows image.Image) error {
		b := rows.Bounds()
		if y0 != next || b.Min.Y != y0 || b.Dx() != pl.W {
			t.Fatalf("band at %d, bounds %v, after %d rows", y0, b, next)
		}
		if out == nil {
			switch rows.(type) {
			case *image.Gray:
				out = image.NewGray(image.Rect(0, 0, pl.W, pl.H))
			default:
				out = image.NewRGBA(image.Rect(0, 0, pl.W, pl.H))
			}
		}
		draw.Draw(out, b, rows, b.Min, draw.Src)
		next = b.Max.Y
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if next != pl.H {
		t.Fatalf("the bands end at %d of %d rows", next, pl.H)
	}
	return out
}

// turned returns an image turned clockwise by a quarter turn.
func turned(img image.Image) *image.RGBA {
	b := img.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, b.Dy(), b.Dx()))
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			out.Set(b.Dy()-1-y, x, img.At(b.Min.X+x, b.Min.Y+y))
		}
	}
	return out
}

// nearPixels counts the pixels whose channels differ by more than tol,
// and returns the largest difference.
func nearPixels(a, b image.Image, tol int) (n, most int) {
	ab, bb := a.Bounds(), b.Bounds()
	if ab.Size() != bb.Size() {
		return -1, 0
	}
	for y := 0; y < ab.Dy(); y++ {
		for x := 0; x < ab.Dx(); x++ {
			r1, g1, b1, a1 := a.At(ab.Min.X+x, ab.Min.Y+y).RGBA()
			r2, g2, b2, a2 := b.At(bb.Min.X+x, bb.Min.Y+y).RGBA()
			d := 0
			for _, p := range [][2]uint32{{r1, r2}, {g1, g2}, {b1, b2}, {a1, a2}} {
				d = max(d, int(math.Abs(float64(p[0]>>8)-float64(p[1]>>8))))
			}
			most = max(most, d)
			if d > tol {
				n++
			}
		}
	}
	return n, most
}

// placeDoc is a page of 100 × 160 units with text, a stroke, an image of
// a gradient, a rotated rectangle and a shadow offset to the right and
// down.
func placeDoc(t *testing.T) (*bdf.Document, *bdf.View) {
	t.Helper()
	font, err := os.ReadFile("../../converter/pptx/testdata/fonts/MPLUS1p-Bold-subset.ttf")
	if err != nil {
		t.Fatal(err)
	}
	d := bdf.NewDocument()
	o := bdf.NewObject()
	fontRef := o.AddFont(bdf.EmbeddedFont(d.AddFont(font), 700, bdf.StyleNormal))
	o.FillColor(0xf8f0e0ff).FillRect(0, 0, 100, 160)
	o.Save().Shadow(0x00000080, 0, 12, 4).FillColor(0x2040c0ff).FillRect(10, 10, 30, 20).Restore()
	o.Font(fontRef, 14).FillColor(0x000000ff).FillText("Place", 10, 60, 0)
	o.StrokeColor(0xc00000ff).Line(2, bdf.CapRound, bdf.JoinRound, 10)
	o.StrokePath(o.AddPath((&bdf.Path{}).MoveTo(10, 80).LineTo(90, 100).LineTo(30, 150)))
	o.Save().Transform(0.8, 0.6, -0.6, 0.8, 70, 110).FillColor(0x00a000ff).FillRect(-8, -8, 16, 16).Restore()
	h, _ := d.AddObject(o)
	v := d.NewView("v", bdf.ViewFixed, "")
	v.AddPage(100, 160, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	return d, v
}

// TestPlaceAsPage checks that a page placed at its own size, unturned
// at the origin, is the page that PageRows draws.
func TestPlaceAsPage(t *testing.T) {
	d, v := placeDoc(t)
	for _, mode := range []AntialiasMode{AntialiasOn, AntialiasOff} {
		r := New(d, &Options{FontDirs: testFonts, NoSystemFonts: true, Antialias: mode, BandRows: 32, ImageCacheBytes: -1})
		for _, scale := range []float64{1, 2.5} {
			whole, err := r.Page(v, 0, scale)
			if err != nil {
				t.Fatal(err)
			}
			b := whole.Bounds()
			got := placedRows(t, r, v, Place{W: b.Dx(), H: b.Dy(), Scale: scale}, 0)
			if n, at := samePixels(whole, got); n != 0 {
				t.Errorf("antialias %d, scale %g: %d pixels differ, the first at %v", mode, scale, n, at)
			}
		}
	}
}

// TestPlaceTurned checks that a page turned a quarter turn is the page
// drawn upright and turned: the pixels of geometry, text and the shadow
// (whose offset turns with the page) land where the turned image has
// them.
func TestPlaceTurned(t *testing.T) {
	d, v := placeDoc(t)
	r := New(d, &Options{FontDirs: testFonts, NoSystemFonts: true, BandRows: 40, ImageCacheBytes: -1})
	upright, err := r.Page(v, 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	want := turned(upright) // 320 × 200
	got := placedRows(t, r, v, Place{W: 320, H: 200, Scale: 2, Turns: 1}, 0)
	// Rasterizing turned geometry is not bit for bit the turned raster
	// of the upright geometry (the edges are sampled across instead of
	// along), so a channel may differ by a few levels at an edge.
	n, most := nearPixels(want, got, 12)
	if n > 0 {
		t.Errorf("a quarter turn: %d pixels differ by more than 12 (at most %d)", n, most)
	}
	t.Logf("a quarter turn: at most %d levels apart", most)
	for turns := 2; turns <= 4; turns++ {
		want = turned(want)
		if turns == 4 {
			break
		}
		w, h := want.Bounds().Dx(), want.Bounds().Dy()
		got := placedRows(t, r, v, Place{W: w, H: h, Scale: 2, Turns: turns}, 0)
		if n, most := nearPixels(want, got, 12); n > 0 {
			t.Errorf("%d quarter turns: %d pixels differ by more than 12 (at most %d)", turns, n, most)
		}
	}
}

// TestPlaceOnPaper checks a page placed on a larger sheet: the margins
// are the background, and the page is where X and Y put it.
func TestPlaceOnPaper(t *testing.T) {
	d, v := placeDoc(t)
	r := New(d, &Options{FontDirs: testFonts, NoSystemFonts: true, BandRows: 16, ImageCacheBytes: -1})
	upright, err := r.PageMono(v, 0, 1, Gray)
	if err != nil {
		t.Fatal(err)
	}
	got := placedRows(t, r, v, Place{W: 140, H: 220, Scale: 1, X: 20, Y: 30}, Gray)
	white := color.Gray{Y: 255}
	for _, p := range []image.Point{{0, 0}, {19, 100}, {139, 219}, {60, 29}, {60, 190}, {120, 100}} {
		if c := got.At(p.X, p.Y); c != white {
			t.Errorf("the margin at %v is %v", p, c)
		}
	}
	if n, _ := samePixels(upright, got.(*image.Gray).SubImage(image.Rect(20, 30, 120, 190))); n != 0 {
		t.Errorf("the page on the sheet: %d pixels differ", n)
	}
}

// mirrored returns an image mirrored left to right, or top to bottom.
func mirrored(img image.Image, x bool) *image.RGBA {
	b := img.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for y := 0; y < b.Dy(); y++ {
		for i := 0; i < b.Dx(); i++ {
			if x {
				out.Set(b.Dx()-1-i, y, img.At(b.Min.X+i, b.Min.Y+y))
			} else {
				out.Set(i, b.Dy()-1-y, img.At(b.Min.X+i, b.Min.Y+y))
			}
		}
	}
	return out
}

// TestPlaceMirrored checks a page mirrored as the back of a sheet:
// left to right, top to bottom, and both (a half turn).
func TestPlaceMirrored(t *testing.T) {
	d, v := placeDoc(t)
	r := New(d, &Options{FontDirs: testFonts, NoSystemFonts: true, BandRows: 40, ImageCacheBytes: -1})
	upright, err := r.Page(v, 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	pl := Place{W: 200, H: 320, Scale: 2}
	for _, c := range []struct {
		x, y bool
		want image.Image
	}{{true, false, mirrored(upright, true)}, {false, true, mirrored(upright, false)}, {true, true, turned(turned(upright))}} {
		pl.MirrorX, pl.MirrorY = c.x, c.y
		got := placedRows(t, r, v, pl, 0)
		if n, most := nearPixels(c.want, got, 12); n > 0 {
			t.Errorf("mirror x %v y %v: %d pixels differ by more than 12 (at most %d)", c.x, c.y, n, most)
		}
	}
}

func TestPlaceErrors(t *testing.T) {
	d, v := placeDoc(t)
	r := New(d, nil)
	nop := func(int, image.Image) error { return nil }
	for name, pl := range map[string]Place{
		"no width":  {W: 0, H: 10, Scale: 1},
		"no scale":  {W: 10, H: 10},
		"NaN":       {W: 10, H: 10, Scale: 1, X: math.NaN()},
		"turns":     {W: 10, H: 10, Scale: 1, Turns: 4},
		"too large": {W: 1 << 15, H: 1 << 15, Scale: 1},
	} {
		if err := r.PlaceRows(v, 0, pl, 0, nop); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
	if err := r.PlaceRows(v, 3, Place{W: 10, H: 10, Scale: 1}, 0, nop); err == nil {
		t.Error("no such page: no error")
	}
}
