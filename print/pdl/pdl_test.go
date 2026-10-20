package pdl

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/image/barcode"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

// testDoc is a view of pages of a size: text, rules, a photograph (a
// gradient of grays), a Code 128 barcode and a QR code, each after its
// mark, on each page.
func testDoc(t *testing.T, w, h float32, pages int) (*bdf.Document, *bdf.View) {
	t.Helper()
	font, err := os.ReadFile("../../converter/pptx/testdata/fonts/MPLUS1p-Bold-subset.ttf")
	if err != nil {
		t.Fatal(err)
	}
	photo := image.NewGray(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			photo.Pix[y*32+x] = uint8(x * 8)
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, photo)
	d := bdf.NewDocument()
	img := d.AddImage(buf.Bytes())
	bar, err := barcode.Encode(barcode.Code128, "BDF-0042", nil)
	if err != nil {
		t.Fatal(err)
	}
	qr, err := barcode.Encode(barcode.QR, "https://example.com/r/42", nil)
	if err != nil {
		t.Fatal(err)
	}
	v := d.NewView("v", bdf.ViewFixed, "")
	for i := range pages {
		o := bdf.NewObject()
		f := o.AddFont(bdf.EmbeddedFont(d.AddFont(font), 700, bdf.StyleNormal))
		o.Font(f, 14).FillColor(0x000000ff).FillText("Label "+string(rune('A'+i)), 8, 20, 0)
		o.FillRect(8, 26, w-16, 1.5)
		o.StrokeColor(0x000000ff).Line(1, bdf.CapButt, bdf.JoinMiter, 10).StrokeRect(4, 4, w-8, h-8)
		o.Image(o.AddImage(img), 8, 32, 40, 30)
		bar.Fill(o, 54.3, 32.6, w-62, 30)
		qr.Fill(o, 8, 68, 50, 50)
		hb, _ := d.AddObject(o)
		v.AddPage(w, h, bdf.Layer{Role: bdf.RoleBody, Obj: hb})
	}
	return d, v
}

// wantDots draws a page as Write draws it, without the symbols of the
// kinds the printer prints (none when drop is nil).
func wantDots(t *testing.T, d *bdf.Document, v *bdf.View, page int, l label, drop []string) *Bitmap {
	t.Helper()
	r := imagebdf.New(d, &imagebdf.Options{NoSystemFonts: true, ImageCacheBytes: -1, Antialias: imagebdf.AntialiasOff,
		DitherImages: true, AlignBarcodes: true, Barcode: func(b imagebdf.Barcode) bool { return slices.Contains(drop, b.Mark.Kind) }})
	bm := NewBitmap(l.W, l.H)
	if err := r.PlaceRows(v, page, l.place, imagebdf.Bilevel, func(y0 int, rows image.Image) error {
		bm.pack(y0, rows)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return bm
}

// dump writes a bitmap as a PNG file in $BDF_PDL_OUT, to look at.
func dump(t *testing.T, name string, b *Bitmap) {
	dir := os.Getenv("BDF_PDL_OUT")
	if dir == "" {
		return
	}
	img := image.NewGray(image.Rect(0, 0, b.W, b.H))
	for y := 0; y < b.H; y++ {
		for x := 0; x < b.W; x++ {
			if !b.Black(x, y) {
				img.Pix[y*img.Stride+x] = 255
			}
		}
	}
	f, err := os.Create(filepath.Join(dir, name+".png"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	png.Encode(f, img)
}

// sameDots compares bitmaps, the ones of a printer that may be wider (the
// dots a command rounds up to bytes) being white past the other's width.
func sameDots(t *testing.T, what string, got, want *Bitmap) {
	t.Helper()
	dump(t, strings.NewReplacer(" ", "-", "/", "-").Replace(what), got)
	if got.H != want.H || got.W < want.W {
		t.Fatalf("%s: %d × %d dots, want %d × %d", what, got.W, got.H, want.W, want.H)
	}
	for y := 0; y < got.H; y++ {
		for x := 0; x < got.W; x++ {
			if g, w := got.Black(x, y), want.Black(x, y); g != w {
				t.Fatalf("%s: the dot at %d, %d is black %v, want %v", what, x, y, g, w)
			}
		}
	}
}

func TestProfiles(t *testing.T) {
	for _, p := range Profiles() {
		if err := p.Check(); err != nil {
			t.Errorf("%s: %v", p.Name, err)
		}
		if q, ok := LookupProfile(p.Name); !ok || q != p {
			t.Errorf("%s: not found by its name", p.Name)
		}
	}
	for _, lang := range Languages() {
		p, err := ParseProfile("", lang, "", "")
		if err != nil || p.Language != lang || p.Barcodes != IsLabel(lang) {
			t.Errorf("%s: %+v %v", lang, p, err)
		}
	}
	p, err := ParseProfile("zd421-203", "", "300", "4in")
	if err != nil || p.DPI != 300 || p.Width != 1200 || p.Language != ZPL {
		t.Errorf("zd421 at 300 dpi, 4 inches: %+v %v", p, err)
	}
	p, err = ParseProfile("", "esc/pos", "", "58mm")
	if err != nil || p.Width != 464 || p.Language != ESCPOS {
		t.Errorf("58 mm of ESC/POS: %+v %v", p, err)
	}
	for _, bad := range [][4]string{{"nope", "", "", ""}, {"", "pcl", "", ""}, {"", "", "", ""}, {"", "zpl", "x", ""}, {"", "zpl", "", "-3mm"}, {"", "zpl", "20", ""}} {
		if p, err := ParseProfile(bad[0], bad[1], bad[2], bad[3]); err == nil {
			t.Errorf("%v: %+v", bad, p)
		}
	}
}

func TestSettingsCheck(t *testing.T) {
	for _, s := range []Settings{{Darkness: 101}, {Speed: -1}, {Tracking: "holes"}, {Cut: "now"}, {Drawer: 3}, {Mode: "text"}, {Mode: "vector"}, {Gap: 1000}} {
		if err := s.Check(); err == nil {
			t.Errorf("%+v: no error", s)
		}
	}
	if err := (&Settings{Darkness: -100, Speed: 5080, Tracking: "gap", Cut: CutJob, Drawer: 2, Mode: "raster"}).Check(); err != nil {
		t.Error(err)
	}
}

func TestSequence(t *testing.T) {
	rd := &render{o: &Options{Pages: []int{0, 1}, Copies: 2, Collate: true}}
	if p, n := rd.sequence(true); len(p) != 4 || p[2] != 0 || n != 1 {
		t.Errorf("collated: %v %d", p, n)
	}
	rd.o.Collate = false
	if p, n := rd.sequence(true); len(p) != 2 || n != 2 {
		t.Errorf("uncollated, the printer's copies: %v %d", p, n)
	}
	if p, n := rd.sequence(false); len(p) != 4 || p[1] != 0 || n != 1 {
		t.Errorf("uncollated: %v %d", p, n)
	}
	rd.o.Pages, rd.o.Collate = []int{3}, true
	if p, n := rd.sequence(true); len(p) != 1 || n != 2 {
		t.Errorf("one page, collated: %v %d", p, n)
	}
}
