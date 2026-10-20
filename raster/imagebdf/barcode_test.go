package imagebdf

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"sync"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/image/barcode"
)

// symbolDoc is a page of 200 × 120 points with a Code 128 barcode and a
// QR code in boxes that fall between the pixels, each after its mark.
func symbolDoc(t *testing.T) (*bdf.Document, *bdf.View, []*barcode.Symbol) {
	t.Helper()
	bar, err := barcode.Encode(barcode.Code128, "BDF-12345", nil)
	if err != nil {
		t.Fatal(err)
	}
	qr, err := barcode.Encode(barcode.QR, "https://example.com/bdf", &barcode.Options{Level: "Q"})
	if err != nil {
		t.Fatal(err)
	}
	d := bdf.NewDocument()
	o := bdf.NewObject()
	o.FillColor(0x000000ff)
	bar.Fill(o, 10.3, 10.7, 120.9, 30.2)
	qr.Fill(o, 140.45, 20.15, 50.3, 50.3)
	o.FillRect(10, 100, 50, 2) // not a symbol
	h, _ := d.AddObject(o)
	v := d.NewView("v", bdf.ViewFixed, "")
	v.AddPage(200, 120, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	return d, v, []*barcode.Symbol{bar, qr}
}

// runs returns the lengths of the runs of black and white of a row of a
// bilevel image, from x0 to x1.
func runs(img *image.Paletted, y, x0, x1 int) []int {
	var out []int
	n, last := 0, -1
	for x := x0; x < x1; x++ {
		v := int(img.ColorIndexAt(x, y))
		if v != last && n > 0 {
			out = append(out, n)
			n = 0
		}
		last = v
		n++
	}
	return append(out, n)
}

// TestAlignBarcodes checks that the modules of the symbols that marks tell
// are drawn on whole pixels at 203 dpi, as the symbols' grids say.
func TestAlignBarcodes(t *testing.T) {
	d, v, syms := symbolDoc(t)
	scale := 203.0 / 72
	for _, rows := range []int{7, 1000} {
		r := New(d, &Options{Antialias: AntialiasOff, AlignBarcodes: true, BandRows: rows, ImageCacheBytes: -1})
		img, err := r.PageMono(v, 0, scale, Bilevel)
		if err != nil {
			t.Fatal(err)
		}
		bl := img.(*image.Paletted)
		if w := r.Warnings(); len(w) > 0 {
			t.Fatal(w)
		}
		// the barcode: every run across its middle a whole number of
		// modules of the same width
		bar := syms[0]
		x0, x1 := int(10.3*scale), int((10.3+120.9)*scale)+1
		rs := runs(bl, int(25*scale), x0, x1)
		mod := 0
		for _, n := range rs[1 : len(rs)-1] {
			if mod == 0 || n < mod {
				mod = n
			}
		}
		if want := int(120.9 * scale / float64(bar.W)); mod != want {
			t.Fatalf("rows %d: the narrowest module is %d pixels, want %d (%v)", rows, mod, want, rs)
		}
		for _, n := range rs[1 : len(rs)-1] {
			if n%mod != 0 {
				t.Fatalf("rows %d: a run of %d pixels is not of modules of %d", rows, n, mod)
			}
		}
		// the dark modules of the row, as the grid has them
		bx0, bx1 := bar.Bars()
		start := rs[0]
		for x := bx0; x < bx1; x++ {
			px := x0 + start + (x-bx0)*mod + mod/2
			if got := bl.ColorIndexAt(px, int(25*scale)) == 0; got != bar.At(x, 0) {
				t.Fatalf("rows %d: module %d is dark %v, the grid says %v", rows, x, got, bar.At(x, 0))
			}
		}
		// the QR code: each module a square of whole pixels, as the grid
		qr := syms[1]
		qmod := int(50.3 * scale / float64(qr.W))
		ox := 140.45*scale + (50.3*scale-float64(qmod*qr.W))/2
		oy := 20.15*scale + (50.3*scale-float64(qmod*qr.H))/2
		for y := 0; y < qr.H; y++ {
			for x := 0; x < qr.W; x++ {
				for _, c := range [][2]int{{0, 0}, {qmod - 1, qmod - 1}} {
					px, py := int(ox)+x*qmod+c[0], int(oy)+y*qmod+c[1]
					// the corner of the frame is rounded, so allow its pixel either way
					if got := bl.ColorIndexAt(px, py) == 0; got != qr.At(x, y) {
						px2, py2 := int(ox+0.5)+x*qmod+c[0], int(oy+0.5)+y*qmod+c[1]
						if got2 := bl.ColorIndexAt(px2, py2) == 0; got2 != qr.At(x, y) {
							t.Fatalf("rows %d: QR module %d,%d is dark %v, the grid says %v", rows, x, y, got, qr.At(x, y))
						}
					}
				}
			}
		}
	}
}

// TestBarcodeCallback checks that the symbols are told with their boxes on
// the image and their turns, and left out when the callback takes them.
func TestBarcodeCallback(t *testing.T) {
	d, v, _ := symbolDoc(t)
	for _, turns := range []int{0, 1, 2, 3} {
		var mu sync.Mutex
		seen := map[string]Barcode{}
		r := New(d, &Options{Antialias: AntialiasOff, BandRows: 16, ImageCacheBytes: -1, Barcode: func(b Barcode) bool {
			mu.Lock()
			defer mu.Unlock()
			if old, ok := seen[b.Mark.Text]; ok && (old.X != b.X || old.Y != b.Y || old.W != b.W || old.H != b.H || old.Turns != b.Turns) {
				t.Errorf("%q told as %+v and as %+v", b.Mark.Text, old, b)
			}
			seen[b.Mark.Text] = b
			return true
		}})
		pl := Place{W: 200, H: 120, Scale: 1, Turns: turns}
		if turns%2 == 1 {
			pl.W, pl.H = 120, 200
		}
		img := placedRows(t, r, v, pl, Bilevel)
		if len(seen) != 2 {
			t.Fatalf("turns %d: told %v", turns, seen)
		}
		for text, b := range seen {
			if b.Turns != turns {
				t.Errorf("turns %d: %q is told turned %d", turns, text, b.Turns)
			}
		}
		if b := seen["BDF-12345"]; b.W < 120 || b.W > 122 || b.H < 30 || b.H > 31 {
			t.Errorf("turns %d: the barcode is %+v", turns, b)
		}
		if turns == 0 {
			if b := seen["BDF-12345"]; b.X < 10 || b.X > 11 || b.Y < 10 || b.Y > 11 {
				t.Errorf("the barcode is at %g, %g", b.X, b.Y)
			}
		}
		// what is left is the rule
		dark := 0
		bounds := img.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				if c := color.GrayModel.Convert(img.At(x, y)).(color.Gray); c.Y < 128 {
					dark++
				}
			}
		}
		if dark != 100 {
			t.Errorf("turns %d: %d dark pixels, want those of the rule (100)", turns, dark)
		}
	}
}

// TestDitherImages checks that an image of a gray is drawn as dots of
// about its darkness, the same in bands as whole, and that the text and
// rules around it are not dithered.
func TestDitherImages(t *testing.T) {
	gray := image.NewGray(image.Rect(0, 0, 16, 16))
	for i := range gray.Pix {
		gray.Pix[i] = 0x60 // 37.6% luminance
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, gray); err != nil {
		t.Fatal(err)
	}
	d := bdf.NewDocument()
	o := bdf.NewObject()
	o.Image(o.AddImage(d.AddImage(buf.Bytes())), 0, 0, 64, 64)
	o.FillColor(0x808080ff).FillRect(64, 0, 64, 64) // a gray fill: cut at the half
	h, _ := d.AddObject(o)
	v := d.NewView("v", bdf.ViewFixed, "")
	v.AddPage(128, 64, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	var whole *image.RGBA
	for _, rows := range []int{5, 1000} {
		r := New(d, &Options{Antialias: AntialiasOff, DitherImages: true, BandRows: rows, ImageCacheBytes: -1})
		p := placedRows(t, r, v, Place{W: 128, H: 64, Scale: 1}, Bilevel).(*image.RGBA)
		isDark := func(x, y int) bool { return p.Pix[p.PixOffset(x, y)] < 128 }
		dark, fill := 0, 0
		for y := 0; y < 64; y++ {
			for x := 0; x < 64; x++ {
				if isDark(x, y) {
					dark++
				}
				if isDark(64+x, y) {
					fill++
				}
			}
		}
		if want := 64 * 64 * 0.624; float64(dark) < want-64 || float64(dark) > want+64 {
			t.Errorf("rows %d: %d of 4096 pixels of the image dark, want about %g", rows, dark, want)
		}
		if fill != 0 {
			t.Errorf("rows %d: %d pixels of the fill are dark (a gray of 50.2%% is white)", rows, fill)
		}
		if whole == nil {
			whole = p
		} else if !bytes.Equal(whole.Pix, p.Pix) {
			t.Errorf("the dots of bands of %d rows differ from those of the page", rows)
		}
	}
}
