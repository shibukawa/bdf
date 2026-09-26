package tiff

import (
	"bytes"
	"encoding/binary"
	"image"
	"math"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/tiff"
	"github.com/shibukawa/bdf/imgconv"
)

// The files in testdata come from test/tiff/gen.sh.

func convert(t *testing.T, name string, opts *Options) *Result {
	t.Helper()
	res, err := ConvertFile("testdata/"+name, opts)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	if len(res.Warnings) > 0 {
		t.Errorf("%s: warnings %q", name, res.Warnings)
	}
	return res
}

// pageImage returns the image part a page draws.
func pageImage(t *testing.T, doc *bdf.Document, p *bdf.Page) []byte {
	t.Helper()
	if len(p.Layers) != 1 {
		t.Fatalf("page with %d layers", len(p.Layers))
	}
	o, err := bdf.DecodeObject(doc.Part(p.Layers[0].Obj).Data)
	if err != nil || len(o.Images) != 1 {
		t.Fatalf("page object: %v, %d images", err, len(o.Images))
	}
	return doc.Part(o.Images[0]).Data
}

func near(a float32, b float64) bool { return math.Abs(float64(a)-b) < 0.01 }

func TestScan(t *testing.T) {
	res := convert(t, "scan.tif", nil)
	doc := res.Doc
	pages := doc.Views[0].Pages
	if len(pages) != 3 || res.Pages != 3 || res.Scaled != 1 {
		t.Fatalf("%d pages, %d scaled; want 3 and 1", len(pages), res.Scaled)
	}
	for i, want := range [][2]float64{
		{1000.0 / 300 * 72, 1300.0 / 300 * 72}, // 300 dpi
		{300.0 / 150 * 72, 390.0 / 150 * 72},   // 150 dpi
		{200.0 / 96 * 72, 150.0 / 96 * 72},     // no resolution: DefaultDPI
	} {
		if p := pages[i]; !near(p.W, want[0]) || !near(p.H, want[1]) {
			t.Errorf("page %d is %v × %v pt, want %.2f × %.2f", i+1, p.W, p.H, want[0], want[1])
		}
	}
	// Page 1, a 300 dpi bilevel scan, is scaled down to 192 dpi and stays
	// bilevel: a 1-bit PNG under Keep.
	img := pageImage(t, doc, pages[0])
	if imgconv.Sniff(img) != "png" || img[24] != 1 {
		t.Errorf("page 1: %s of depth %d, want a 1-bit PNG", imgconv.Sniff(img), img[24])
	}
	if w, h := binary.BigEndian.Uint32(img[16:]), binary.BigEndian.Uint32(img[20:]); w != 640 || h != 832 {
		t.Errorf("page 1 stored at %d×%d, want 640×832 (192 dpi)", w, h)
	}
	// Page 2, 150 dpi JPEG strips, is stored as the JPEG stream the strips
	// join into.
	f, err := os.ReadFile("testdata/scan.tif")
	if err != nil {
		t.Fatal(err)
	}
	tf, err := tiff.Open(bytes.NewReader(f), int64(len(f)))
	if err != nil {
		t.Fatal(err)
	}
	joined, ok := tf.Pages()[1].JPEG()
	if !ok || !bytes.Equal(pageImage(t, doc, pages[1]), joined) {
		t.Error("page 2 is not stored as its joined JPEG strips")
	}
	// The descriptive tags of the first page.
	dc := doc.Meta.DC
	if doc.Meta.Source != "tiff" || dc.Title.First() != "Scan test" || dc.Description.First() != "Three pages of a scanner" ||
		!slices.Equal(dc.Creator, bdf.DCValues{"Alice", "Bob"}) || dc.Rights.First() != "(c) 2026 Example" ||
		dc.Modified.First() != "2026-09-26T18:30:00" || doc.Views[0].Title != "Scan test" {
		t.Errorf("metadata %+v, source %q", dc, doc.Meta.Source)
	}
}

func TestFax(t *testing.T) {
	res := convert(t, "fax.tif", nil)
	pages := res.Doc.Views[0].Pages
	if len(pages) != 2 {
		t.Fatalf("%d pages", len(pages))
	}
	// 1728 × 280 pixels at 204 × 98 dpi: the page keeps the fax's
	// proportions, and only the 204 dpi axis is scaled down.
	for i, p := range pages {
		if !near(p.W, 1728.0/204*72) || !near(p.H, 280.0/98*72) {
			t.Errorf("page %d is %v × %v pt", i+1, p.W, p.H)
		}
		img := pageImage(t, res.Doc, p)
		if w, h := binary.BigEndian.Uint32(img[16:]), binary.BigEndian.Uint32(img[20:]); w != 1626 || h != 280 {
			t.Errorf("page %d stored at %d×%d, want 1626×280", i+1, w, h)
		}
	}
	if bytes.Equal(pageImage(t, res.Doc, pages[0]), pageImage(t, res.Doc, pages[1])) {
		t.Error("the two fax pages are the same image")
	}
}

// TestOrientation puts the stored image of each page on its page with the
// page's matrix, as the viewer does, and compares it with page 1 (which
// has the default orientation): each is the same picture stored turned or
// mirrored, with the Orientation tag that undoes it.
func TestOrientation(t *testing.T) {
	b, err := os.ReadFile("testdata/orientation.tif")
	if err != nil {
		t.Fatal(err)
	}
	f, err := tiff.Open(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		t.Fatal(err)
	}
	pages := f.Pages()
	if len(pages) != 8 {
		t.Fatalf("%d pages", len(pages))
	}
	var want *image.RGBA
	for i, d := range pages {
		p := newPage(d, DefaultDPI)
		if p.orientation != i+1 {
			t.Fatalf("page %d has orientation %d", i+1, p.orientation)
		}
		img, _, err := d.Decode()
		if err != nil {
			t.Fatal(err)
		}
		m := p.transform()
		// The inverse of the matrix maps page points to image points.
		det := m[0]*m[3] - m[1]*m[2]
		pw, ph := int(p.pageW), int(p.pageH)
		out := image.NewRGBA(image.Rect(0, 0, pw, ph))
		sx, sy := float32(img.Bounds().Dx())/p.imgW, float32(img.Bounds().Dy())/p.imgH
		for y := range ph {
			for x := range pw {
				px, py := float32(x)+0.5-m[4], float32(y)+0.5-m[5]
				ix, iy := (m[3]*px-m[2]*py)/det, (-m[1]*px+m[0]*py)/det
				out.Set(x, y, img.At(int(ix*sx), int(iy*sy)))
			}
		}
		if i == 0 {
			want = out
			continue
		}
		if out.Rect != want.Rect {
			t.Errorf("orientation %d: page %v, want %v", i+1, out.Rect, want.Rect)
			continue
		}
		bad := 0
		for k := range out.Pix {
			if d := int(out.Pix[k]) - int(want.Pix[k]); d > 8 || d < -8 {
				bad++
			}
		}
		if bad > 0 {
			t.Errorf("orientation %d: %d channel values differ from page 1", i+1, bad)
		}
	}
}

func TestPageSelection(t *testing.T) {
	res := convert(t, "scan.tif", &Options{Pages: converter.PageList(3, 1, 3)})
	pages := res.Doc.Views[0].Pages
	if len(pages) != 2 || !near(pages[0].W, 150) || !near(pages[1].W, 240) {
		t.Errorf("pages 3, 1: got %d pages", len(pages))
	}
	if _, err := ConvertFile("testdata/scan.tif", &Options{Pages: converter.PageList(4)}); err == nil {
		t.Error("no error for page 4 of 3")
	}
}

func TestCaps(t *testing.T) {
	// Without caps nothing is scaled.
	res := convert(t, "scan.tif", &Options{Images: imgconv.Options{MaxDPI: -1, MaxPixels: -1}})
	if res.Scaled != 0 {
		t.Errorf("%d pages scaled without caps", res.Scaled)
	}
	// At 100 dpi the JPEG page is scaled too, and stored as JPEG again.
	res = convert(t, "scan.tif", &Options{Images: imgconv.Options{MaxDPI: 100}})
	if res.Scaled != 2 {
		t.Errorf("%d pages scaled to 100 dpi, want 2", res.Scaled)
	}
	img := pageImage(t, res.Doc, res.Doc.Views[0].Pages[1])
	if imgconv.Sniff(img) != "jpeg" {
		t.Errorf("scaled JPEG page stored as %s", imgconv.Sniff(img))
	}
	// A pixel cap alone.
	res = convert(t, "fax.tif", &Options{Images: imgconv.Options{MaxDPI: -1, MaxPixels: 100_000}})
	img = pageImage(t, res.Doc, res.Doc.Views[0].Pages[0])
	if w, h := binary.BigEndian.Uint32(img[16:]), binary.BigEndian.Uint32(img[20:]); w*h > 100_000 || w*h < 99_000 {
		t.Errorf("fax page stored at %d×%d for a cap of 100000 pixels", w, h)
	}
}

func TestConvertMode(t *testing.T) {
	if !imgconv.Available() {
		t.Skip("codecs not compiled in")
	}
	res := convert(t, "scan.tif", &Options{Images: imgconv.Options{Mode: imgconv.Convert}})
	if f := imgconv.Sniff(pageImage(t, res.Doc, res.Doc.Views[0].Pages[0])); f != "webp" {
		t.Errorf("bilevel page under Convert stored as %s", f)
	}
}

func TestRegistered(t *testing.T) {
	b, err := os.ReadFile("testdata/scan.tif")
	if err != nil {
		t.Fatal(err)
	}
	if f := converter.Detect(bytes.NewReader(b), int64(len(b))); f == nil || f.Name != "tiff" {
		t.Fatalf("scan.tif detected as %v", f)
	}
	res, err := converter.Convert(bytes.NewReader(b), int64(len(b)), "", &converter.Options{Params: map[string]string{"dpi": "200"}})
	if err != nil {
		t.Fatal(err)
	}
	if p := res.Doc.Views[0].Pages[2]; !near(p.W, 72) {
		t.Errorf("page without a resolution at dpi=200 is %v pt wide, want 72", p.W)
	}
	if !strings.Contains(res.Summary, "3 page(s), 2 scaled down") { // page 3 at 200 dpi is over the cap too
		t.Errorf("summary %q", res.Summary)
	}
	if _, err := converter.Convert(bytes.NewReader(b), int64(len(b)), "", &converter.Options{Params: map[string]string{"dpi": "x"}}); err == nil {
		t.Error("no error for dpi=x")
	}
	cr2 := append([]byte("II*\x00\x10\x00\x00\x00CR\x02\x00"), make([]byte, 32)...)
	if f := converter.Detect(bytes.NewReader(cr2), int64(len(cr2))); f != nil {
		t.Errorf("Canon raw file detected as %s", f.Name)
	}
}

func TestDNG(t *testing.T) {
	// A 1×1 grey image with a DNGVersion tag.
	var b bytes.Buffer
	b.WriteString("II*\x00\x08\x00\x00\x00")
	le := binary.LittleEndian
	binary.Write(&b, le, uint16(6))
	for _, e := range [][4]uint32{{256, 3, 1, 1}, {257, 3, 1, 1}, {258, 3, 1, 8}, {262, 3, 1, 1}, {273, 4, 1, 0}, {50706, 1, 4, 0x00000401}} {
		binary.Write(&b, le, uint16(e[0]))
		binary.Write(&b, le, uint16(e[1]))
		binary.Write(&b, le, e[2])
		binary.Write(&b, le, e[3])
	}
	binary.Write(&b, le, uint32(0))
	_, err := Convert(bytes.NewReader(b.Bytes()), int64(b.Len()), nil)
	if err == nil || !strings.Contains(err.Error(), "DNG") {
		t.Errorf("DNG file: %v", err)
	}
}

// TestDamaged checks that a page that does not decode is left blank with
// a warning, and the others are converted.
func TestDamaged(t *testing.T) {
	b, err := os.ReadFile("testdata/scan.tif")
	if err != nil {
		t.Fatal(err)
	}
	tf, _ := tiff.Open(bytes.NewReader(b), int64(len(b)))
	// Give page 3 an unknown compression scheme.
	d := tf.Pages()[2]
	b = slices.Clone(b)
	for i := 0; i+12 <= len(b); i += 2 {
		if int64(i) > d.Offset && int64(i) < d.Offset+2+20*12 && binary.LittleEndian.Uint16(b[i:]) == tiff.TagCompression &&
			binary.LittleEndian.Uint16(b[i+2:]) == 3 {
			binary.LittleEndian.PutUint16(b[i+8:], 40000)
			break
		}
	}
	res, err := Convert(bytes.NewReader(b), int64(len(b)), nil)
	if err != nil {
		t.Fatal(err)
	}
	pages := res.Doc.Views[0].Pages
	if len(pages) != 3 || len(pages[2].Layers) != 0 || len(pages[0].Layers) != 1 {
		t.Fatalf("pages %d, layers of page 3: %d", len(pages), len(pages[2].Layers))
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "page 3") {
		t.Errorf("warnings %q", res.Warnings)
	}
}
