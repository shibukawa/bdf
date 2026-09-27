package hpgl

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"golang.org/x/text/encoding/japanese"
)

// The test plots are made by test/hpgl/gen.py; text is laid out with the
// test fonts of the PowerPoint converter.

func testOptions() *Options {
	return &Options{FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true}
}

func convert(t *testing.T, name string, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
	if opts == nil {
		opts = testOptions()
	}
	res, err := ConvertFile(filepath.Join("testdata", name), opts)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := res.Doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	return res, r
}

func plainText(t *testing.T, r *bdf.Reader) string {
	t.Helper()
	h, err := bdf.ParseHash(r.Manifest.Views[0].TextIndex)
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.Part(h)
	if err != nil {
		t.Fatal(err)
	}
	idx, err := bdf.DecodeTextIndex(b)
	if err != nil {
		t.Fatal(err)
	}
	return bdf.PlainText(idx)
}

// plotted is a plot read with the converter's state left to inspect.
type plotted struct {
	*converter
	ended bool
}

// plot reads a plot: its pages are ended by the first describe, the pen
// is where the plot left it until then.
func plot(t *testing.T, data string) *plotted {
	t.Helper()
	c := newConverter(testOptions())
	db := fontdb.New(nil, []string{"../pptx/testdata/fonts"}, false)
	c.set = fontset.New(db, nil)
	c.fonts = &cad.Fonts{Set: c.set}
	c.doc = bdf.NewDocument()
	c.run([]byte(data))
	return &plotted{converter: c}
}

// describe returns the items of a page's drawing, coordinates rounded to
// whole plotter units.
func describe(c *plotted, page int) []string {
	if !c.ended {
		c.endPage()
		c.ended = true
	}
	if page >= len(c.pages) {
		return nil
	}
	return c.pages[page].d.Describe(1)
}

func near(a, b, tol float64) bool { return math.Abs(a-b) <= tol }

func TestShapes(t *testing.T) {
	res, r := convert(t, "shapes.plt", nil)
	if len(res.Warnings) > 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	m := r.Manifest
	if len(m.Views) != 1 || len(m.Views[0].Pages) != 1 || m.Meta.Source != "hpgl" {
		t.Fatalf("manifest: %+v", m)
	}
	// the plot size is A3
	if s := res.Sizes[0]; !near(s[0], 420, 0.1) || !near(s[1], 297, 0.1) {
		t.Errorf("page %v mm", s)
	}
	if got := m.Meta.DC.Title.First(); got != "HP-GL/2 test plot" {
		t.Errorf("title %q", got)
	}
	text := plainText(t, r)
	for _, want := range []string{"Stick font 0.25 x 0.35 cm", "Large", "Default size: 9 characters per inch", "Turned 45", "Slanted",
		"Line one", "Line two", "LO1", "LO19", "DOWN", "Univers Bold 14 pt", "CG Times Italic 12 pt", "文字 テスト", "テスト図"} {
		if !strings.Contains(text, want) {
			t.Errorf("text lacks %q:\n%s", want, text)
		}
	}
	if res.EmbeddedFonts == 0 {
		t.Error("no fonts embedded")
	}
}

// TestJob reads a job for a large-format plotter: PJL, HP-GL/2 turned to a
// portrait sheet, raster images.
func TestJob(t *testing.T) {
	res, r := convert(t, "job.plt", nil)
	if len(res.Warnings) > 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	if res.Pages != 2 || res.Images != 2 {
		t.Fatalf("%d pages, %d images", res.Pages, res.Images)
	}
	// the first page is seen as RO 90 turned it: portrait A4
	if s := res.Sizes[0]; !near(s[0], 210, 0.1) || !near(s[1], 297, 0.1) {
		t.Errorf("page 1: %v mm", s)
	}
	// the second has a raster image only, 400 × 320 pixels at 300 dpi,
	// and a margin of 5 mm
	if s := res.Sizes[1]; !near(s[0], 400.0/300*25.4+10, 0.1) || !near(s[1], 320.0/300*25.4+10, 0.1) {
		t.Errorf("page 2: %v mm", s)
	}
	if got := r.Manifest.Meta.DC.Title.First(); got != "Plotter job page" {
		t.Errorf("title %q", got)
	}
	if text := plainText(t, r); !strings.Contains(text, "Portrait page") || !strings.Contains(text, "Raster image above") {
		t.Errorf("text %q", text)
	}
	// the image of the first page is upright below the pen: its top left
	// corner at (1000, 6000) of the turned system, 240 × 180 pixels at
	// 150 dpi
	c := plot(t, string(mustRead(t, "testdata/job.plt")))
	var img string
	for _, s := range describe(c, 0) {
		if strings.HasPrefix(s, "image") {
			img = s
		}
	}
	// in the drawing (unturned): x = 11880 − y′, y = x′
	w, h := 240*1016/150.0, 180*1016/150.0
	want := fmt.Sprintf("image %g %g %g %g 240x180", 11880-6000.0, 1000.0, math.Round(11880-6000+h), math.Round(1000+w))
	if img != want {
		t.Errorf("%s, want %s", img, want)
	}
}

func TestHPGL1(t *testing.T) {
	res, r := convert(t, "hpgl1.plt", nil)
	if len(res.Warnings) > 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	text := plainText(t, r)
	for _, want := range []string{"HP-GL PEN PLOTTER", "VERTICAL"} {
		if !strings.Contains(text, want) {
			t.Errorf("text lacks %q: %q", want, text)
		}
	}
	// no plot size: the page is the drawing and a margin of 5 mm
	if s := res.Sizes[0]; !near(s[0], (10250-250)/40.0+10, 0.5) {
		t.Errorf("page %v mm", s)
	}
}

// TestPrinter puts HP-GL/2 in the picture frame of a PCL page: A4
// landscape, the frame between the top and bottom margins of half an inch.
func TestPrinter(t *testing.T) {
	c := plot(t, "\x1bE\x1b&l26A\x1b&l1O\x1b%0BIN;SP1;PA0,0;PD1016,0;PU;\x1b%0A\x1bE")
	d := describe(c, 0)
	if len(c.pages) != 1 || c.pages[0].pcl == nil {
		t.Fatalf("%d pages", len(c.pages))
	}
	if p := c.pages[0].pcl; !near(p.w, 297*40, 1) || !near(p.h, 210*40, 1) {
		t.Errorf("page %v × %v", p.w, p.h)
	}
	// from the left edge of the logical page (0.2 in) on the bottom margin
	y := 210*40 - 1016/2.0
	want := fmt.Sprintf("stroke %g %g %g %g color=000000ff width=0 dash=[]", math.Round(0.2*1016), y, math.Round(1.2*1016), y)
	if len(d) != 1 || d[0] != want {
		t.Errorf("%q, want %q", d, want)
	}
}

func TestPages(t *testing.T) {
	c := plot(t, "IN;SP1;PD100,0;PG;PG;SP2;PD0,100;PG;")
	describe(c, 0)
	if len(c.pages) != 2 {
		t.Fatalf("%d pages (PG on an empty page makes none)", len(c.pages))
	}
	// BP starts a new plot
	c = plot(t, `IN;SP1;PD100,0;BP1,"two";SP1;PD0,100;`)
	describe(c, 0)
	if len(c.pages) != 2 || c.title != "two" {
		t.Errorf("%d pages, title %q", len(c.pages), c.title)
	}
	// a page selection
	opts := testOptions()
	opts.Pages = conv.PageList(2)
	data := "IN;SP1;PD100,0;PG;SP1;PD0,100;PG;"
	res, err := Convert(bytes.NewReader([]byte(data)), int64(len(data)), opts)
	if err != nil {
		t.Fatal(err)
	}
	if res.Pages != 1 {
		t.Errorf("%d pages", res.Pages)
	}
}

// TestScaling checks the user units of SC: anisotropic, isotropic (the
// space left shared by left and bottom) and point factor.
func TestScaling(t *testing.T) {
	p1, p2 := cad.Point{X: 0, Y: 0}, cad.Point{X: 2000, Y: 1000}
	for _, c := range []struct {
		s    scaling
		u    cad.Point
		want cad.Point
	}{
		{scaling{xmin: 0, xmax: 100, ymin: 0, ymax: 10}, cad.Point{X: 50, Y: 5}, cad.Point{X: 1000, Y: 500}},
		// isotropic: 10 units a side fit 1000, centred in 2000
		{scaling{typ: 1, xmin: 0, xmax: 10, ymin: 0, ymax: 10, left: 50, bottom: 50}, cad.Point{X: 0, Y: 0}, cad.Point{X: 500, Y: 0}},
		{scaling{typ: 1, xmin: 0, xmax: 10, ymin: 0, ymax: 10, left: 0, bottom: 50}, cad.Point{X: 10, Y: 10}, cad.Point{X: 1000, Y: 1000}},
		// mirrored: xmin on the side of P1
		{scaling{typ: 1, xmin: 10, xmax: 0, ymin: 0, ymax: 10, left: 50, bottom: 50}, cad.Point{X: 10, Y: 0}, cad.Point{X: 500, Y: 0}},
		{scaling{typ: 2, xmin: 5, xmax: 40, ymin: 0, ymax: 40}, cad.Point{X: 6, Y: 1}, cad.Point{X: 40, Y: 40}},
	} {
		m, ok := c.s.matrix(p1, p2)
		if got := apply(m, c.u); !ok || !near(got.X, c.want.X, 1e-9) || !near(got.Y, c.want.Y, 1e-9) {
			t.Errorf("%+v: %v → %v, want %v", c.s, c.u, got, c.want)
		}
	}
}

// TestRotate turns the coordinate system: the pen keeps its place.
func TestRotate(t *testing.T) {
	c := plot(t, "IN;PS10000,5000;SP1;PU1000,500;RO90;PD1000,500;PU;")
	// (1000, 500) turned 90°: x = 10000 − 500, y = 1000
	want := "stroke 1000 500 9500 1000 color=000000ff width=0 dash=[]"
	if got := describe(c, 0); len(got) != 1 || got[0] != want {
		t.Errorf("%q, want %q", got, want)
	}
}

// TestPolylineEncoded decodes PE: absolute and relative pairs, pen-up
// flags, fractional bits and a pen, in base 64 and base 32.
func TestPolylineEncoded(t *testing.T) {
	enc := func(v int, base int) string {
		u := 2 * v
		if v < 0 {
			u = -2*v + 1
		}
		var b []byte
		for u >= base {
			b = append(b, byte(63+u%base))
			u /= base
		}
		if base == 64 {
			b = append(b, byte(191+u))
		} else {
			b = append(b, byte(95+u))
		}
		return string(b)
	}
	pe := "PE:" + enc(2, 64) + ">" + enc(1, 64) + "<=" + enc(200, 64) + enc(400, 64) + enc(200, 64) + enc(0, 64) + "7" + enc(-100, 32) + enc(100, 32) + ";"
	c := plot(t, "IN;SP1;NP4;PC2,0,0,255;"+pe)
	if g := c.gl; g.pos != (cad.Point{X: 150, Y: 250}) || !g.down || g.relative {
		t.Errorf("after PE: %v down %v relative %v", g.pos, g.down, g.relative)
	}
	// pen 2, one fractional bit: (100, 200) up, then to (200, 200) and
	// (150, 250)
	want := "stroke 100 200 200 250 color=0000ffff width=0 dash=[]"
	if got := describe(c, 0); len(got) != 1 || got[0] != want {
		t.Errorf("%q, want %q", got, want)
	}
}

// TestPolygon fills a subpolygon with a hole: the move after PM1 is not a
// side, and PM2 restores the pen.
func TestPolygon(t *testing.T) {
	c := plot(t, "IN;SP1;PU0,0;PM0;PD100,0,100,100,0,100;PM1;PU20,20;PD80,20,80,80,20,80;PM2;FP;EP;PA200,0;")
	if g := c.gl; g.down || g.pos != (cad.Point{X: 200, Y: 0}) {
		t.Errorf("pen %v down %v", g.pos, g.down)
	}
	d := describe(c, 0)
	if len(d) != 2 || !strings.HasPrefix(d[0], "fill 0 0 100 100") || !strings.HasPrefix(d[1], "stroke 0 0 100 100") {
		t.Fatalf("%q", d)
	}
	g := c.gl
	// two closed subpolygons of four vertices each, the second from (20, 20)
	if v := cad.Flatten(g.polygon.fill, 1); len(v) != 8 || v[4] != (cad.Point{X: 20, Y: 20}) {
		t.Errorf("fill vertices %v", v)
	}
}

// TestArcs draws arcs through three points and with coarse chords.
func TestArcs(t *testing.T) {
	c := plot(t, "IN;SP1;PU0,0;PD;AT100,100,200,0;PU;")
	if d := describe(c, 0); len(d) != 1 || !strings.HasPrefix(d[0], "stroke 0 0 200 1") {
		// the control points of a half circle over (0,0)–(200,0) reach a
		// little above its top, 100
		t.Errorf("%q", d)
	}
	// CI with a chord angle of 90 is a square
	c = plot(t, "IN;SP1;PA0,0;CI100,90;")
	d := describe(c, 0)
	if len(d) != 1 || d[0] != "stroke -100 -100 100 100 color=000000ff width=0 dash=[]" {
		t.Errorf("%q", d)
	}
}

// TestLineTypes checks the dash patterns of fixed and adaptive line types.
func TestLineTypes(t *testing.T) {
	// a pattern of 10 mm (400 plotter units): line type 2 is half dash
	c := plot(t, "IN;SP1;PW0;LT2,10,1;PA0,0;PD1000,0;PU;")
	if d := describe(c, 0); len(d) != 1 || !strings.HasSuffix(d[0], "dash=[200 200]") {
		t.Errorf("%q", d)
	}
	// adaptive: 1000 is 2.5 patterns, drawn as 3 of 333⅓ starting in the
	// middle of a dash
	c = plot(t, "IN;SP1;PW0;LT-2,10,1;PA0,0;PD1000,0;PU;")
	if d := describe(c, 0); len(d) != 1 || !strings.HasSuffix(d[0], "dash=[167 167]") {
		t.Errorf("%q", d)
	}
	// LT;LT99 restores the line type
	c = plot(t, "IN;SP1;LT3,10,1;LT;LT99;PA0,0;PD1000,0;")
	if d := describe(c, 0); len(d) != 1 || !strings.HasSuffix(d[0], "dash=[280 120]") {
		t.Errorf("%q", d)
	}
}

// TestLabel places labels by their origin: the cap height and the pitch
// of the stick font.
func TestLabel(t *testing.T) {
	c := plot(t, "IN;SP1;SI1,1.5;PA1000,1000;LO5;LBAB\x03")
	// the pen does not move after a centred label, but does after LO1
	if g := c.gl; g.pos != (cad.Point{X: 1000, Y: 1000}) {
		t.Errorf("pen %v", g.pos)
	}
	d := describe(c, 0)
	if len(d) != 1 || !strings.HasPrefix(d[0], `text "AB"`) {
		t.Fatalf("%q", d)
	}
	// two cells of 1.5 × 1 cm centred on the pen, and the cap height of
	// 1.5 cm: the text starts at (1000 − 600, 1000 − 300)
	if !strings.HasSuffix(d[0], " 400 700] vertical=false") {
		t.Errorf("%q", d)
	}
	c = plot(t, "IN;SP1;SI1,1.5;PA1000,1000;LBAB\x03")
	if g := c.gl; !near(g.pos.X, 1000+2*600, 1e-6) || g.pos.Y != 1000 {
		t.Errorf("pen %v", g.pos)
	}
}

func TestDecode(t *testing.T) {
	g := newGL(&converter{warned: map[string]bool{}})
	text := func(b []byte) string {
		var sb strings.Builder
		for _, ch := range g.decode(b) {
			if !ch.ctrl {
				sb.WriteRune(ch.r)
			}
		}
		return sb.String()
	}
	// Roman-8
	if got := text([]byte("caf\xc5 \xb3C")); got != "café °C" {
		t.Errorf("Roman-8: %q", got)
	}
	// Shift_JIS by its kana and kanji
	sjis, _ := japanese.ShiftJIS.NewEncoder().String("図面 テスト")
	if got := text([]byte(sjis)); got != "図面 テスト" {
		t.Errorf("Shift_JIS: %q", got)
	}
	// accented Roman-8 letters followed by ASCII letters are not kanji
	if got := text([]byte("\xe7N \xdaZ")); got != "ÓN ÖZ" {
		t.Errorf("Roman-8 pairs: %q", got)
	}
	// 16-bit JIS
	g.label.mode16 = true
	if got := text([]byte("\x25\x46\x25\x39\x25\x48\x00A")); got != "テストA" {
		t.Errorf("JIS: %q", got)
	}
}

func TestDecompress(t *testing.T) {
	if got := packbits([]byte{0xfd, 'U', 0, 'A', 0xff, 'T'}, 100); string(got) != "UUUUATT" {
		t.Errorf("packbits %q", got)
	}
	// the delta rows of the reference guide
	r1 := deltaRow(nil, []byte{0x01, 0xff}, 100)
	r2 := deltaRow(r1, []byte{0x02, 0xf0}, 100)
	r3 := deltaRow(r2, []byte{0x00, 0x0f, 0x22, 0xaa, 0xaa}, 100)
	if !slices.Equal(r3, []byte{0x0f, 0xff, 0xf0, 0xaa, 0xaa}) {
		t.Errorf("delta rows %x %x %x", r1, r2, r3)
	}
	// replacement delta row: a literal byte at offset 2, then a run of
	// four at offset 1
	seed := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	got := replacementDelta(seed, []byte{2 << 3, 0xaa, 0x80 | 1<<5 | 2, 0xbb}, 100)
	if !slices.Equal(got, []byte{1, 2, 0xaa, 4, 0xbb, 0xbb, 0xbb, 0xbb, 9}) {
		t.Errorf("method 9: %x", got)
	}
	// an offset of 15 or more takes extra bytes
	got = replacementDelta(nil, []byte{15<<3 | 0, 5, 0xcc}, 100)
	if len(got) != 21 || got[20] != 0xcc {
		t.Errorf("method 9 offset: %x", got)
	}
	// rows stop at the width they may have
	if got := packbits([]byte{0x81, 1, 0x81, 1, 0x81, 1}, 200); len(got) > 200+128 {
		t.Errorf("packbits: %d bytes", len(got))
	}
	if got := replacementDelta(nil, []byte{0xff, 255, 255, 255, 0, 1}, 50); len(got) > 50 {
		t.Errorf("method 9: %d bytes", len(got))
	}
}

// TestReviewCases are plots that went wrong: they must not panic, and must
// draw what they draw.
func TestReviewCases(t *testing.T) {
	// a direct-by-plane encoding of too many planes is left out
	plot(t, "\x1bE\x1b*v6W\x00\x02\x01\x0a\x0a\x0a\x1b*r1A\x1b*b2W\xff\xff\x1b*rC")
	// PS in polygon mode ends it; CI then draws
	if d := describe(plot(t, "IN;SP1;PM0;PS;CI100;"), 0); len(d) != 1 {
		t.Errorf("CI after PS in polygon mode: %q", d)
	}
	// a user line type of one stretch, adaptive
	if d := describe(plot(t, "IN;SP1;UL1,100;LT-1;PA0,0;PD100,0;"), 0); len(d) != 1 {
		t.Errorf("UL1,100: %q", d)
	}
	// PS closes the window
	if d := describe(plot(t, "IN;IW0,0,5000,5000;PS20000,10000;SP1;PA100,100;PD2000,2000;"), 0); len(d) != 1 {
		t.Errorf("IW then PS: %q", d)
	}
	// BP in a PCL job is IN: one page
	c := plot(t, "\x1bE\x1b&l26A\x1b%0BIN;SP1;PA0,0;PD1000,0;\x1b%0A\x1b%0BBP;IN;SP1;PA0,1000;PD1000,1000;\x1bE")
	if describe(c, 0); len(c.pages) != 1 {
		t.Errorf("BP in PCL: %d pages", len(c.pages))
	}
	// a raster image not ended at the end of the file
	data := "\x1bE\x1b*r1A\x1b*b1W\xff\x1b*b1W\xff"
	if res, err := Convert(bytes.NewReader([]byte(data)), int64(len(data)), testOptions()); err != nil || res.Images != 1 {
		t.Errorf("open raster image: %v, %d images", err, res.Images)
	}
	// a vector of no length is a dot
	if d := describe(plot(t, "IN;SP1;PU100,100;PD100,100;PU;"), 0); len(d) != 1 || !strings.Contains(d[0], "stroke 100 100 100 100") {
		t.Errorf("zero-length vector: %q", d)
	}
	// RO turns the window to (9500–10000, 0–1000), which clips the stroke
	// from (10000, 0) to (9100, 400)
	d := describe(plot(t, "IN;PS10000,5000;SP1;IW0,0,1000,500;RO90;PA0,0;PD400,900;"), 0)
	if len(d) != 2 || d[0] != "group 9500 0 10000 400" {
		t.Errorf("RO and IW: %q", d)
	}
	// SO and SI draw in the size of the font for the rest of the label; SI
	// holds for the next one
	c = plot(t, "IN;SP1;SI2,3;PA0,0;LBA\x0eB\x0fC\x03PA0,5000;LBD\x03")
	d = describe(c, 0)
	if len(d) != 3 {
		t.Fatalf("%q", d)
	}
	scale := func(s string) float64 { var v float64; fmt.Sscan(s[strings.Index(s, "m=[")+3:], &v); return v }
	if a, bc, dd := scale(d[0]), scale(d[1]), scale(d[2]); a != dd || bc >= a {
		t.Errorf("sizes %v %v %v", a, bc, dd)
	}
}

// TestRaster draws a black and white raster image of a PCL page: the
// white pixels are transparent, and the image is 8 × 2 pixels at 75 dpi.
func TestRaster(t *testing.T) {
	c := plot(t, "\x1bE\x1b&l2A\x1b*t75R\x1b*r1A\x1b*b1W\xf0\x1b*b1W\x0f\x1b*rC\x1bE")
	if len(c.pages) != 1 || c.images != 1 {
		t.Fatalf("%d pages, %d images", len(c.pages), c.images)
	}
	d := describe(c, 0)
	// at the top margin, from the left edge of the logical page
	want := fmt.Sprintf("image %g %g %g %g", 0.25*1016, 1016/2.0, 0.25*1016+8*1016/75.0, 1016/2.0+2*1016/75.0)
	if len(d) != 1 || !near(parseBox(d[0])[2], parseBox(want)[2], 1) || !near(parseBox(d[0])[0], parseBox(want)[0], 1) {
		t.Errorf("%q, want %q", d, want)
	}
}

func parseBox(s string) [4]float64 {
	var b [4]float64
	var kind string
	fmt.Sscan(s, &kind, &b[0], &b[1], &b[2], &b[3])
	return b
}

// TestHostile checks that broken and hostile input ends quickly.
func TestHostile(t *testing.T) {
	for _, data := range []string{
		"IN;LBno terminator",
		"IN;PM0;" + strings.Repeat("PD1,1,2,2;", 1000),
		"IN;SP1;CI1e300;AA0,0,1e30;SC0,0,0,0;IP0,0,0,0;PS-1;",
		"IN;SP1;FT3,0.0001,45;RA100000,100000;",
		"IN;RF1,256,256;FT11,1;RA10,10;",
		"\x1bE\x1b*r65535S\x1b*r65535T\x1b*r1A\x1b*b9M\x1b*b3W\xff\xff\xff\x1b*rC",
		"\x1b%-12345X@PJL ENTER LANGUAGE=HPGL2\r\nIN;SP1;PE" + strings.Repeat("\xff", 100) + ";",
		"IN;DT;;SMx;PAx;LM1;LBxx\x00\x03",
	} {
		done := make(chan struct{})
		go func() {
			defer close(done)
			Convert(bytes.NewReader([]byte(data)), int64(len(data)), testOptions())
		}()
		select {
		case <-done:
		case <-time.After(20 * time.Second):
			t.Fatalf("%q did not finish", data)
		}
	}
}

func TestNotHPGL(t *testing.T) {
	for _, s := range []string{"hello", "set terminal png\nset output 'a.png'\nplot sin(x) title 'sine'\n"} {
		if _, err := Convert(bytes.NewReader([]byte(s)), int64(len(s)), nil); err == nil {
			t.Errorf("no error for %q", s)
		}
	}
	// a plot that draws nothing is a blank page
	if res, err := Convert(bytes.NewReader([]byte("IN;SP1;PG;")), 10, testOptions()); err != nil || res.Pages != 1 {
		t.Errorf("empty plot: %v", err)
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// TestAccumulator scales a 600 dpi image down as its rows arrive: by 3 to
// stay under 192 dpi, a two-color image staying two-colored, the ink kept
// where it covers two fifths of a block.
func TestAccumulator(t *testing.T) {
	r := newRTL(&converter{opts: &Options{}, warned: map[string]bool{}})
	img := &raster{r: r, pw: 1016 / 600.0, ph: 1016 / 600.0, width: 30, planes: 1}
	img.acc.init(r, img)
	if img.acc.factor != 3 || !img.acc.bilevel {
		t.Fatalf("factor %d, bilevel %v", img.acc.factor, img.acc.bilevel)
	}
	ink := color.NRGBA{0, 0, 0, 255}
	for y := 0; y < 30; y++ {
		row := make([]color.NRGBA, 30)
		for x := range row {
			// a column of blocks with one pixel of ink in three, and one with
			// two in three
			if x < 15 && x%3 == 0 || x >= 15 && x%3 != 0 {
				row[x] = ink
			}
		}
		img.addRow(row)
	}
	pix, ok := img.acc.image().(*image.Paletted)
	if !ok || pix.Bounds().Dx() != 10 || pix.Bounds().Dy() != 10 {
		t.Fatalf("%T %v", pix, pix.Bounds())
	}
	if pix.ColorIndexAt(0, 0) != 0 || pix.ColorIndexAt(9, 9) != 1 {
		t.Errorf("blocks %d %d", pix.ColorIndexAt(0, 0), pix.ColorIndexAt(9, 9))
	}
}
