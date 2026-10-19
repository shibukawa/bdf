package raster

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

func TestEncodeLine(t *testing.T) {
	for _, c := range []struct {
		line []byte
		bpp  int
		want []byte
	}{
		{[]byte{7}, 1, []byte{0, 0, 7}},
		{[]byte{7, 7, 7}, 1, []byte{0, 2, 7}},
		{[]byte{1, 2, 3}, 1, []byte{0, 254, 1, 2, 3}},
		{[]byte{1, 2, 2, 2, 3}, 1, []byte{0, 0, 1, 2, 2, 0, 3}},
		{[]byte{1, 2, 3, 3}, 1, []byte{0, 255, 1, 2, 1, 3}},
		{[]byte{1, 2, 3, 1, 2, 3, 9, 9, 9}, 3, []byte{0, 1, 1, 2, 3, 0, 9, 9, 9}},
	} {
		got := encodeLine(nil, 0, c.line, c.bpp)
		if !bytes.Equal(got, c.want) {
			t.Errorf("encodeLine(%v, %d) = %v, want %v", c.line, c.bpp, got, c.want)
		}
	}
	// runs longer than 128 are split
	long := bytes.Repeat([]byte{5}, 300)
	if got := encodeLine(nil, 0, long, 1); !bytes.Equal(got, []byte{0, 127, 5, 127, 5, 43, 5}) {
		t.Errorf("a long run: %v", got)
	}
	var lit []byte
	for i := 0; i < 200; i++ {
		lit = append(lit, byte(i))
	}
	got := encodeLine(nil, 0, lit, 1)
	if got[1] != 129 || got[130] != byte(257-72) || len(got) != 1+1+128+1+72 {
		t.Errorf("a long literal: %d bytes, %d %d", len(got), got[1], got[130])
	}
}

// doc is a page of 72 × 108 points (1 × 1.5 inch) with colour, gray,
// text-like strokes and a gradient-like row of bars.
func doc(t *testing.T) (*bdf.Document, *bdf.View) {
	t.Helper()
	d := bdf.NewDocument()
	o := bdf.NewObject()
	o.FillColor(0xffffffff).FillRect(0, 0, 72, 108)
	o.FillColor(0xd02020ff).FillRect(6, 6, 30, 20)
	o.FillColor(0x2060c0ff).FillRect(40, 10, 24, 40)
	for i := 0; i < 8; i++ {
		o.FillColor(bdf.Color(uint32(i*32)<<24|uint32(i*32)<<16|uint32(i*32)<<8|0xff)).FillRect(float32(4+i*8), 60, 8, 12)
	}
	o.StrokeColor(0x000000ff).Line(1.5, bdf.CapRound, bdf.JoinRound, 10)
	o.StrokePath(o.AddPath((&bdf.Path{}).MoveTo(6, 80).LineTo(66, 100).LineTo(10, 104)))
	h, _ := d.AddObject(o)
	v := d.NewView("v", bdf.ViewFixed, "")
	v.AddPage(72, 108, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	v.AddPage(72, 108, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	return d, v
}

// write draws the pages of a view placed on a sheet and writes them as
// PWG Raster.
func write(t *testing.T, r *imagebdf.Renderer, v *bdf.View, pl imagebdf.Place, typ Type, mono imagebdf.Mono) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := NewWriter(&buf)
	for page := range v.Pages {
		if err := w.StartPage(Page{Width: pl.W, Height: pl.H, XDPI: 144, YDPI: 144, Type: typ, PageWidth: 72, PageHeight: 108,
			PageSizeName: "custom_25.4x38.1mm", TotalPages: len(v.Pages)}); err != nil {
			t.Fatal(err)
		}
		if err := r.PlaceRows(v, page, pl, mono, func(_ int, rows image.Image) error { return w.WriteRows(rows) }); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// TestRoundTrip writes pages in each type and reads them back: the
// pixels are those the renderer draws.
func TestRoundTrip(t *testing.T) {
	d, v := doc(t)
	r := imagebdf.New(d, &imagebdf.Options{BandRows: 20, ImageCacheBytes: -1})
	rOff := imagebdf.New(d, &imagebdf.Options{BandRows: 20, ImageCacheBytes: -1, Antialias: imagebdf.AntialiasOff})
	pl := imagebdf.Place{W: 144, H: 216, Scale: 2}
	for _, c := range []struct {
		typ  Type
		r    *imagebdf.Renderer
		mono imagebdf.Mono
	}{{SGray8, r, imagebdf.Gray}, {SRGB8, r, 0}, {Black1, rOff, imagebdf.Bilevel}, {SGray1, rOff, imagebdf.Bilevel}, {SGray8, r, 0}} {
		b := write(t, c.r, v, pl, c.typ, c.mono)
		pages, err := Read(bytes.NewReader(b), 0)
		if err != nil {
			t.Fatalf("%s: %v", c.typ, err)
		}
		if len(pages) != 2 {
			t.Fatalf("%s: %d pages", c.typ, len(pages))
		}
		p := pages[1]
		if p.Type != c.typ || p.Width != 144 || p.Height != 216 || p.XDPI != 144 || p.PageWidth != 72 || p.PageSizeName != "custom_25.4x38.1mm" ||
			p.TotalPages != 2 || p.CrossFeed != 1 || p.Feed != 1 || p.Copies != 1 {
			t.Fatalf("%s: header %+v", c.typ, p.Page)
		}
		if p.Lines >= p.Height {
			t.Errorf("%s: %d lines for %d rows: no line repeated", c.typ, p.Lines, p.Height)
		}
		mono := c.mono
		if c.typ == SGray8 && mono == 0 {
			mono = imagebdf.Gray // written from RGBA, compared with the renderer's gray
		}
		want, err := c.r.PageMono(v, 1, 2, max(mono, imagebdf.Gray))
		if mono == 0 {
			want, err = c.r.Page(v, 1, 2)
		}
		if err != nil {
			t.Fatal(err)
		}
		if n, most := differ(want, p.Image); n > 0 {
			t.Errorf("%s from mono %d: %d pixels differ (by %d at most)", c.typ, c.mono, n, most)
		}
	}
}

// differ counts the pixels whose gray or colour differs by more than 1
// level (the rounding of RGBA to gray), and the largest difference.
func differ(a, b image.Image) (n, most int) {
	ab := a.Bounds()
	if ab.Size() != b.Bounds().Size() {
		return -1, 0
	}
	for y := 0; y < ab.Dy(); y++ {
		for x := 0; x < ab.Dx(); x++ {
			r1, g1, b1, _ := a.At(ab.Min.X+x, ab.Min.Y+y).RGBA()
			r2, g2, b2, _ := b.At(x, y).RGBA()
			if p, ok := a.(*image.Paletted); ok {
				// black first: the reader gives 0 for ink
				if p.ColorIndexAt(ab.Min.X+x, ab.Min.Y+y) == 0 {
					r1, g1, b1 = 0, 0, 0
				} else {
					r1, g1, b1 = 0xffff, 0xffff, 0xffff
				}
			}
			d := 0
			for _, c := range [][2]uint32{{r1, r2}, {g1, g2}, {b1, b2}} {
				d = max(d, int(math.Abs(float64(c[0]>>8)-float64(c[1]>>8))))
			}
			most = max(most, d)
			if d > 1 {
				n++
			}
		}
	}
	return n, most
}

func TestWriterErrors(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	if err := w.WriteRows(image.NewGray(image.Rect(0, 0, 4, 1))); err == nil {
		t.Error("rows before a page")
	}
	if err := w.StartPage(Page{Width: 4, Height: 2, XDPI: 300, YDPI: 300, Type: "cmyk_8"}); err == nil {
		t.Error("an unknown type")
	}
	if err := w.StartPage(Page{Width: 4, Height: 2, XDPI: 300, YDPI: 300, Type: SGray8}); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteRows(image.NewGray(image.Rect(0, 0, 5, 1))); err == nil {
		t.Error("rows of another width")
	}
	if err := w.WriteRows(image.NewGray(image.Rect(0, 0, 4, 1))); err != nil {
		t.Fatal(err)
	}
	if err := w.StartPage(Page{Width: 4, Height: 2, XDPI: 300, YDPI: 300, Type: SGray8}); !errors.Is(err, ErrRows) {
		t.Errorf("a page short of rows: %v", err)
	}
	if err := w.WriteRows(image.NewGray(image.Rect(0, 0, 4, 2))); !errors.Is(err, ErrRows) {
		t.Errorf("rows past the page: %v", err)
	}
	if err := w.Close(); !errors.Is(err, ErrRows) {
		t.Errorf("closing a page short of rows: %v", err)
	}
}

func TestOverWhite(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.StartPage(Page{Width: 2, Height: 1, XDPI: 300, YDPI: 300, Type: SRGB8})
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.RGBA{0, 0, 0, 0})     // transparent: white
	img.Set(1, 0, color.RGBA{128, 0, 0, 128}) // half red over white
	w.WriteRows(img)
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	pages, err := Read(&buf, 0)
	if err != nil {
		t.Fatal(err)
	}
	got := pages[0].Image.(*image.RGBA)
	if c := got.RGBAAt(0, 0); c != (color.RGBA{255, 255, 255, 255}) {
		t.Errorf("transparent: %v", c)
	}
	if c := got.RGBAAt(1, 0); c != (color.RGBA{255, 127, 127, 255}) {
		t.Errorf("half red: %v", c)
	}
}

func TestFit(t *testing.T) {
	a4 := Sheet{Width: 595.276, Height: 841.89, DPI: 300}
	if w, h := a4.Pixels(); w != 2480 || h != 3507 {
		t.Fatalf("A4 at 300 dpi: %d × %d", w, h)
	}
	const dpi = 300.0 / 72
	centre := func(sheet int, page, scale float64) float64 { return math.Round((float64(sheet) - page*scale) / 2) }
	fitA5 := math.Min(595.276/419.528, 841.89/595.276)
	fitLand := math.Min(595.276/841.89, 841.89/595.276)
	for _, c := range []struct {
		name    string
		w, h    float64
		scaling Scaling
		o       Orientation
		turns   int
		f       float64 // the scale of the page, as a factor of its size
		tw, th  float64 // the page's size once turned
	}{
		{"A4 on A4", 595.276, 841.89, ScaleAutoFit, OrientAuto, 0, 1, 595.276, 841.89},
		{"A5 on A4, not enlarged", 419.528, 595.276, ScaleAutoFit, OrientAuto, 0, 1, 419.528, 595.276},
		{"A5 on A4, fit", 419.528, 595.276, ScaleFit, OrientAuto, 0, fitA5, 419.528, 595.276},
		{"A5 on A4, fill", 419.528, 595.276, ScaleFill, OrientAuto, 0, math.Max(595.276/419.528, 841.89/595.276), 419.528, 595.276},
		{"A4 landscape on A4", 841.89, 595.276, ScaleAutoFit, OrientAuto, 3, 1, 595.276, 841.89},
		{"A4 landscape as portrait", 841.89, 595.276, ScaleAutoFit, Portrait, 0, fitLand, 841.89, 595.276},
		{"A3 on A4, shrunk", 841.89, 1190.55, ScaleAuto, OrientAuto, 0, 595.276 / 841.89, 841.89, 1190.55},
		{"reverse landscape", 841.89, 595.276, ScaleNone, ReverseLandscape, 1, 1, 595.276, 841.89},
		{"a half turn", 595.276, 841.89, ScaleNone, ReversePortrait, 2, 1, 595.276, 841.89},
	} {
		pl, err := Fit(c.w, c.h, a4, c.scaling, c.o)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		scale := dpi * c.f
		x, y := centre(2480, c.tw, scale), centre(3507, c.th, scale)
		if pl.W != 2480 || pl.H != 3507 || pl.Turns != c.turns || math.Abs(pl.Scale-scale) > 1e-6 || pl.X != x || pl.Y != y {
			t.Errorf("%s: %+v, want turns %d scale %g at %g, %g", c.name, pl, c.turns, scale, x, y)
		}
	}
	if pl, _ := Fit(595.276, 841.89, a4, ScaleAutoFit, OrientAuto); pl.X != 0 || pl.Y != 0 {
		t.Errorf("A4 on A4 is not at the origin: %+v", pl)
	}
	if _, err := Fit(0, 1, a4, ScaleFit, OrientAuto); err == nil {
		t.Error("an empty page")
	}
	if _, err := Fit(1, 1, a4, "stretch", OrientAuto); err == nil {
		t.Error("an unknown scaling")
	}
}

func TestBack(t *testing.T) {
	for _, c := range []struct {
		back       string
		tumble     bool
		cross, fwd int
	}{
		{"normal", false, 1, 1}, {"normal", true, 1, 1},
		{"flipped", false, 1, -1}, {"flipped", true, -1, 1},
		{"rotated", false, -1, -1}, {"rotated", true, 1, 1},
		{"manual-tumble", false, 1, 1}, {"manual-tumble", true, -1, -1},
	} {
		cf, f, err := Back(c.back, c.tumble)
		if err != nil || cf != c.cross || f != c.fwd {
			t.Errorf("Back(%q, %v) = %d, %d, %v", c.back, c.tumble, cf, f, err)
		}
	}
	if _, _, err := Back("sideways", false); err == nil {
		t.Error("an unknown sheet-back")
	}
	pl := Mirror(imagebdf.Place{}, -1, -1)
	if !pl.MirrorX || !pl.MirrorY {
		t.Errorf("Mirror: %+v", pl)
	}
}

func TestReadRefuses(t *testing.T) {
	if _, err := Read(bytes.NewReader([]byte("RaS3")), 0); err == nil {
		t.Error("another sync word")
	}
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.StartPage(Page{Width: 1000, Height: 1000, XDPI: 300, YDPI: 300, Type: SGray8})
	w.WriteRows(image.NewGray(image.Rect(0, 0, 1000, 1000)))
	w.Close()
	if _, err := Read(bytes.NewReader(buf.Bytes()), 1000); err == nil {
		t.Error("a page past the pixels")
	}
	b := buf.Bytes()
	if _, err := Read(bytes.NewReader(b[:len(b)-1]), 0); err == nil {
		t.Error("a truncated stream")
	}
}

// TestDocuments writes the first page of documents of testdata at 100
// dpi in gray, placed on a sheet of its own size, reads it back and
// compares it with the page the renderer draws; writing it again gives
// the same bytes.
func TestDocuments(t *testing.T) {
	fonts := []string{"../../converter/pptx/testdata/fonts", "../../converter/docx/testdata/fonts", "../../fixture/testdata/fonts"}
	for _, f := range []string{"../../testdata/demo.bdf", "../../testdata/form/pptx-fixed.bdf", "../../testdata/form/rows.bdf", "../../testdata/docx/basic.bdf"} {
		rd, err := bdf.OpenSingleFile(f)
		if err != nil {
			t.Fatal(err)
		}
		d, err := rd.ToDocument()
		rd.Close()
		if err != nil {
			t.Fatal(err)
		}
		var v *bdf.View
		for _, view := range d.Views {
			if len(view.Pages) > 0 && view.Kind != bdf.ViewSheet {
				v = view
				break
			}
		}
		if v == nil {
			t.Fatalf("%s: no view of pages", f)
		}
		p := v.Pages[0]
		r := imagebdf.New(d, &imagebdf.Options{FontDirs: fonts, NoSystemFonts: true, ImageCacheBytes: -1})
		pl, err := Fit(float64(p.W), float64(p.H), Sheet{Width: float64(p.W), Height: float64(p.H), DPI: 100}, ScaleNone, Portrait)
		if err != nil {
			t.Fatal(err)
		}
		once := func() []byte {
			var buf bytes.Buffer
			w := NewWriter(&buf)
			if err := w.StartPage(Page{Width: pl.W, Height: pl.H, XDPI: 100, YDPI: 100, Type: SGray8}); err != nil {
				t.Fatal(err)
			}
			if err := r.PlaceRows(v, 0, pl, imagebdf.Gray, func(_ int, rows image.Image) error { return w.WriteRows(rows) }); err != nil {
				t.Fatal(f, err)
			}
			if err := w.Close(); err != nil {
				t.Fatal(err)
			}
			return buf.Bytes()
		}
		b := once()
		if !bytes.Equal(b, once()) {
			t.Errorf("%s: two writings differ", f)
		}
		pages, err := Read(bytes.NewReader(b), 0)
		if err != nil {
			t.Fatal(f, err)
		}
		// the region of the page that the sheet's pixels show, at the
		// same scale (PageMono would stretch the page to whole pixels)
		region := bdf.Rect{W: float32(float64(pl.W) / pl.Scale), H: float32(float64(pl.H) / pl.Scale)}
		want, err := r.RegionMono(v, 0, region, pl.W, pl.H, imagebdf.Gray)
		if err != nil {
			t.Fatal(err)
		}
		got := pages[0].Image
		if n, most := differ(want, got); n > 0 {
			t.Errorf("%s: %d pixels differ (by %d at most)", f, n, most)
		}
		t.Logf("%s: %d × %d pixels in %d bytes (%d lines)", f, pl.W, pl.H, len(b), pages[0].Lines)
	}
}
