package cgm

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strconv"
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

// The test metafiles are made by test/cgm/gen.py: shapes.cgm (and the
// same metafile in clear text, shapes-text.cgm), illustration.cgm and
// sjis.cgm. Text is laid out with the test fonts of the PowerPoint
// converter.

func testOptions() *Options {
	return &Options{FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true}
}

func convertBytes(t *testing.T, data []byte, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
	if opts == nil {
		opts = testOptions()
	}
	res, err := Convert(bytes.NewReader(data), int64(len(data)), opts)
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

func read(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func convert(t *testing.T, name string, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
	return convertBytes(t, read(t, name), opts)
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

// pictures interprets a metafile and returns its pictures, and the
// warnings.
func pictures(t *testing.T, data []byte) ([]*picture, []string) {
	t.Helper()
	elems, text, _, err := readElements(data)
	if err != nil {
		t.Fatal(err)
	}
	c := &converter{opts: &Options{}, doc: bdf.NewDocument(), warned: map[string]bool{}}
	in := newInterp(c, text)
	set := fontset.New(fontdb.New(nil, []string{"../pptx/testdata/fonts"}, false), func(string) {})
	in.fonts = &cad.Fonts{Set: set}
	in.run(iterate(elems))
	return in.pictures, c.warnings
}

func near(a, b, eps float64) bool { return math.Abs(a-b) <= eps }

func TestShapes(t *testing.T) {
	res, r := convert(t, "shapes.cgm", nil)
	if len(res.Warnings) > 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	m := r.Manifest
	if len(m.Views) != 1 || len(m.Views[0].Pages) != 1 || m.Meta.Source != "cgm" || res.Encoding != "binary" {
		t.Fatalf("manifest: %+v", m)
	}
	// A3 at the metric scale of the picture
	p := m.Views[0].Pages[0]
	if !near(float64(p.W), 420*mm, 0.01) || !near(float64(p.H), 297*mm, 0.01) {
		t.Errorf("page %v × %v pt", p.W, p.H)
	}
	if got := m.Meta.DC.Title.First(); got != "shapes" {
		t.Errorf("title %q", got)
	}
	// JIS X 0208 is in the character set list
	if got := m.Meta.DC.Language.First(); got != "ja" {
		t.Errorf("language %q", got)
	}
	text := plainText(t, r)
	for _, want := range []string{"Turned 45", "Red Blue", "Boxed text", "JIS: 日本語の文字", "縦書き", "CGM test drawing",
		"© converter/cgm", "図枠 A3"} {
		if !strings.Contains(text, want) {
			t.Errorf("text %q lacks %q", text, want)
		}
	}
	if res.EmbeddedFonts == 0 {
		t.Error("no fonts embedded")
	}
}

// The clear text encoding of the same metafile draws the same.
func TestClearTextMatchesBinary(t *testing.T) {
	bin, w1 := pictures(t, read(t, "shapes.cgm"))
	txt, w2 := pictures(t, read(t, "shapes-text.cgm"))
	if len(w1)+len(w2) > 0 {
		t.Errorf("warnings: %q %q", w1, w2)
	}
	if len(bin) != 1 || len(txt) != 1 {
		t.Fatalf("%d and %d pictures", len(bin), len(txt))
	}
	a, b := bin[0].d.Describe(0.5), txt[0].d.Describe(0.5)
	if len(a) < 100 {
		t.Fatalf("only %d items", len(a))
	}
	if !slices.Equal(a, b) {
		for i := range min(len(a), len(b)) {
			if a[i] != b[i] {
				t.Fatalf("item %d:\nbinary %s\ntext   %s", i, a[i], b[i])
			}
		}
		t.Fatalf("%d and %d items", len(a), len(b))
	}
}

func TestIllustration(t *testing.T) {
	res, r := convert(t, "illustration.cgm", nil)
	if len(res.Warnings) > 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	pages := r.Manifest.Views[0].Pages
	if len(pages) != 2 || res.Pictures != 2 {
		t.Fatalf("%d pages", len(pages))
	}
	// abstract pictures are fitted to 297 mm
	if !near(float64(pages[0].W), 297*mm, 0.01) || !near(float64(pages[0].H), 297*mm*600/800, 0.01) {
		t.Errorf("page %v × %v pt", pages[0].W, pages[0].H)
	}
	text := plainText(t, r)
	for _, want := range []string{"Übersicht — ×2 °C", "A1", "Page 2"} {
		if !strings.Contains(text, want) {
			t.Errorf("text %q lacks %q", text, want)
		}
	}
	if strings.Contains(text, "hidden") {
		t.Errorf("the hidden layer is drawn: %q", text)
	}
	images := 0
	for _, p := range r.Manifest.Parts {
		if p.T == bdf.PartImage {
			images++
		}
	}
	if images != 2 {
		t.Errorf("%d images, want the 2 tiles", images)
	}
	// the y axis of the VDC is down, the character up vector (0, -1): the
	// text is upright on the page
	pics, _ := pictures(t, read(t, "illustration.cgm"))
	found := false
	for _, it := range pics[0].d.Items {
		if it.Text() == "Übersicht — ×2 °C" {
			found = true
		}
	}
	if !found {
		t.Fatal("no title")
	}
	desc := pics[0].d.Describe(0.001)
	for _, l := range desc {
		if strings.Contains(l, "Übersicht") {
			// m = [x scale, 0, 0, y scale (negative: up is -y in the VDC), ...]
			var m [6]float64
			f := strings.Fields(strings.TrimSuffix(l[strings.Index(l, "m=[")+3:strings.Index(l, "]")], "]"))
			for i := range min(len(f), 6) {
				m[i], _ = strconv.ParseFloat(f[i], 64)
			}
			if len(f) != 6 || m[0] <= 0 || m[3] >= 0 {
				t.Errorf("title transform %v", m)
			}
			pm := pics[0].m
			if pm[3]*m[3] >= 0 {
				t.Errorf("the title is upside down on the page: %v · %v", pm, m)
			}
		}
	}
}

func TestShiftJIS(t *testing.T) {
	res, r := convert(t, "sjis.cgm", nil)
	if len(res.Warnings) > 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	m := r.Manifest
	if got := m.Meta.DC.Title.First(); got != "テスト図" {
		t.Errorf("title %q", got)
	}
	if got := m.Meta.DC.Language.First(); got != "ja" {
		t.Errorf("language %q", got)
	}
	p := m.Views[0].Pages[0]
	if !near(float64(p.W), 210*mm, 0.01) || !near(float64(p.H), 297*mm, 0.01) {
		t.Errorf("page %v × %v pt", p.W, p.H)
	}
	text := plainText(t, r)
	for _, want := range []string{"図枠", "説明文字 テスト", "日本語の文字"} {
		if !strings.Contains(text, want) {
			t.Errorf("text %q lacks %q", text, want)
		}
	}
	elems, text2, _, _ := readElements(read(t, "sjis.cgm"))
	in := newInterp(&converter{opts: &Options{}, warned: map[string]bool{}}, text2)
	in.scan = true
	in.run(iterate(elems))
	if got := in.fontNames; !slices.Equal(got, []string{"ＭＳ ゴシック", "ＭＳ 明朝"}) {
		t.Errorf("fonts %q", got)
	}
	if in.fontSpecs[0].EastAsian != "MS Gothic" || in.fontSpecs[1].EastAsian != "MS Mincho" {
		t.Errorf("font specs %+v", in.fontSpecs)
	}
}

func TestPages(t *testing.T) {
	opts := testOptions()
	opts.Pages = conv.PageList(2)
	res, r := convert(t, "illustration.cgm", opts)
	if res.Pages != 1 || res.Pictures != 2 {
		t.Fatalf("%d pages of %d pictures", res.Pages, res.Pictures)
	}
	if text := plainText(t, r); !strings.Contains(text, "Page 2") || strings.Contains(text, "Übersicht") {
		t.Errorf("text %q", text)
	}
	opts.Pages = conv.PageList(3)
	data := read(t, "illustration.cgm")
	if _, err := Convert(bytes.NewReader(data), int64(len(data)), opts); err == nil || !strings.Contains(err.Error(), "out of range") {
		t.Errorf("page 3: %v", err)
	}
}

func gzipped(b []byte) []byte {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Write(b)
	zw.Close()
	return buf.Bytes()
}

func TestGzip(t *testing.T) {
	data := gzipped(read(t, "sjis.cgm"))
	if !Detect(data[:min(len(data), 1024)], bytes.NewReader(data), int64(len(data))) {
		t.Fatal("not detected")
	}
	res, _ := convertBytes(t, data, nil)
	if res.Pages != 1 {
		t.Errorf("%d pages", res.Pages)
	}
}

func TestDetect(t *testing.T) {
	for _, c := range []struct {
		name string
		data []byte
		want bool
	}{
		{"binary", read(t, "shapes.cgm"), true},
		{"clear text", read(t, "shapes-text.cgm"), true},
		{"clear text with a comment", []byte("% made by hand %\n  beg_mf 'x';\nmfversion 1;\n"), true},
		{"gzip", gzipped(read(t, "shapes.cgm")), true},
		{"gzip of text", gzipped([]byte("hello, world")), false},
		{"begin metafile alone", []byte{0x00, 0x22, 0x01, 'a', 0x10, 0x22, 0x00, 0x01}, true},
		{"csv", []byte("0,1,2\n3,4,5\n"), false},
		{"word", []byte("BEGMFX 'x';"), false},
		{"empty", nil, false},
		{"pdf", []byte("%PDF-1.7\n"), false},
	} {
		head := c.data[:min(len(c.data), 1024)]
		if got := Detect(head, bytes.NewReader(c.data), int64(len(c.data))); got != c.want {
			t.Errorf("%s: %v", c.name, got)
		}
	}
}

// elem encodes a binary element.
func elem(class, id int, data []byte) []byte {
	var out []byte
	h := uint16(class<<12 | id<<5)
	if len(data) < 31 {
		out = binary.BigEndian.AppendUint16(out, h|uint16(len(data)))
		out = append(out, data...)
	} else {
		out = binary.BigEndian.AppendUint16(out, h|31)
		// two partitions
		half := len(data) / 2 &^ 1
		out = binary.BigEndian.AppendUint16(out, 0x8000|uint16(half))
		out = append(out, data[:half]...)
		out = binary.BigEndian.AppendUint16(out, uint16(len(data)-half))
		out = append(out, data[half:]...)
	}
	if len(out)%2 == 1 {
		out = append(out, 0)
	}
	return out
}

func i16(vs ...int) []byte {
	var b []byte
	for _, v := range vs {
		b = binary.BigEndian.AppendUint16(b, uint16(int16(v)))
	}
	return b
}

// metafile wraps elements of a picture body in a metafile.
func metafile(body ...[]byte) []byte {
	out := slices.Concat(elem(0, 1, []byte{1, 'm'}), elem(1, 1, i16(1)), elem(0, 3, []byte{0}),
		elem(2, 6, i16(0, 0, 100, 100)), elem(0, 4, nil))
	for _, b := range body {
		out = append(out, b...)
	}
	return append(append(out, elem(0, 5, nil)...), elem(0, 2, nil)...)
}

func TestCellArrayRunLength(t *testing.T) {
	// 4 × 2 cells, run-length encoded with 8-bit colour indexes: two runs
	// in the first row, one in the second, each row on a word boundary
	data := slices.Concat(i16(0, 0, 40, 20, 40, 0), i16(4, 2, 8, 0),
		i16(3), []byte{2}, i16(1), []byte{3}, // row 1: 48 bits
		i16(4), []byte{4}, []byte{0}) // row 2: 24 bits, padded to a word
	pics, warns := pictures(t, metafile(elem(4, 9, data)))
	if len(warns) > 0 {
		t.Errorf("warnings %q", warns)
	}
	d := pics[0].d.Describe(0.01)
	if len(d) != 1 || !strings.HasPrefix(d[0], "image 0 0 40 20 4x2") {
		t.Fatalf("%q", d)
	}
	elems, _, _, _ := readElements(metafile(elem(4, 9, data)))
	in := newInterp(&converter{opts: &Options{}, warned: map[string]bool{}}, false)
	for _, e := range elems {
		if e.code == eCellArray {
			p := e.params(&in.pr)
			for range 6 {
				p.vdc()
			}
			p.int()
			p.int()
			p.int()
			mode := p.enum("")
			var got []int
			in.cells(p, 4, 2, 8, mode, maxCells, func(_, _ int, c colour) { got = append(got, c.index) })
			if !slices.Equal(got, []int{2, 2, 2, 3, 4, 4, 4, 4}) {
				t.Errorf("cells %v", got)
			}
		}
	}
}

func TestPartitions(t *testing.T) {
	// a polyline of 20 points (80 bytes) in two partitions
	var pts []int
	for i := range 20 {
		pts = append(pts, i*5, (i%2)*50)
	}
	pics, _ := pictures(t, metafile(elem(4, 1, i16(pts...))))
	d := pics[0].d.Describe(0.01)
	if len(d) != 1 || !strings.HasPrefix(d[0], "stroke 0 0 95 50") {
		t.Errorf("%q", d)
	}
}

func TestArc3(t *testing.T) {
	// counter-clockwise through the top, and clockwise through the bottom
	for _, c := range []struct {
		mid  cad.Point
		want cad.Rect
	}{
		{cad.Point{X: 0, Y: 10}, cad.Rect{}.Add(cad.Point{X: -10, Y: 0}).Add(cad.Point{X: 10, Y: 10})},
		{cad.Point{X: 0, Y: -10}, cad.Rect{}.Add(cad.Point{X: -10, Y: -10}).Add(cad.Point{X: 10, Y: 0})},
	} {
		path, centre := arc3(cad.Point{X: 10}, c.mid, cad.Point{X: -10})
		if centre != (cad.Point{}) {
			t.Errorf("centre %v", centre)
		}
		var b cad.Rect
		for _, q := range cad.Flatten(path, 16) {
			b = b.Add(q)
		}
		if !near(b.Min.X, c.want.Min.X, 0.01) || !near(b.Min.Y, c.want.Min.Y, 0.01) || !near(b.Max.X, c.want.Max.X, 0.01) || !near(b.Max.Y, c.want.Max.Y, 0.01) {
			t.Errorf("through %v: bounds %v, want %v", c.mid, b, c.want)
		}
	}
	// in a line: a straight line
	path, _ := arc3(cad.Point{}, cad.Point{X: 1}, cad.Point{X: 2})
	if pts := cad.Flatten(path, 4); len(pts) != 2 {
		t.Errorf("collinear: %v", pts)
	}
}

func TestDecodeText(t *testing.T) {
	jis := func(s string) []byte {
		// the EUC-JP bytes of JIS X 0208 without their high bits
		b := []byte{}
		for _, c := range mustEncodeEUC(s) {
			b = append(b, c&0x7f)
		}
		return b
	}
	for _, c := range []struct {
		name     string
		charsets []charset
		g0, g1   int
		in       []byte
		want     string
	}{
		{"ascii", nil, 1, 1, []byte("abc"), "abc"},
		{"latin-1 by default", nil, 1, 1, []byte("caf\xe9"), "café"},
		{"shift out to JIS X 0208", []charset{{set94, "B"}, {set94Multi, "B"}}, 1, 2,
			slices.Concat([]byte("A"), []byte{0x0e}, jis("図形"), []byte{0x0f}, []byte("B")), "A図形B"},
		{"8-bit JIS X 0208", []charset{{set94, "B"}, {set94Multi, "$B"}}, 1, 2, mustEncodeEUC("文字"), "文字"},
		{"katakana of JIS X 0201", []charset{{set94, "I"}}, 1, 1, []byte{0x31, 0x32}, "ｱｲ"},
		{"utf-16", []charset{{setComplete, "/L"}}, 1, 1, []byte{0x30, 0xc6, 0x30, 0xb9, 0x30, 0xc8}, "テスト"},
		{"utf-8 with its designation escape", []charset{{setComplete, "\x1b%/I"}}, 1, 1, []byte("日本"), "日本"},
	} {
		in := newInterp(&converter{opts: &Options{}, warned: map[string]bool{}}, false)
		in.charsets = c.charsets
		in.st = newState()
		in.st.a.charSet, in.st.a.altCharSet = c.g0, c.g1
		if got := in.decodeText(c.in); got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
}

func mustEncodeEUC(s string) []byte {
	b, err := japanese.EUCJP.NewEncoder().Bytes([]byte(s))
	if err != nil {
		panic(err)
	}
	return b
}

func TestScaleFactor(t *testing.T) {
	pr := defaultPrecisions()
	fl := binary.BigEndian.AppendUint32(nil, math.Float32bits(0.25))
	fx := slices.Concat(i16(0), []byte{0x40, 0})
	for _, c := range []struct {
		b    []byte
		want float64
	}{{fl, 0.25}, {fx, 0.25}} {
		p := &binParams{b: c.b, pr: &pr}
		if got := p.scaleFactor(); got != c.want {
			t.Errorf("% x: %v", c.b, got)
		}
	}
}

// Truncated and damaged metafiles convert, or fail, without hanging.
func TestDamaged(t *testing.T) {
	for _, name := range []string{"shapes.cgm", "shapes-text.cgm", "illustration.cgm"} {
		data := read(t, name)
		start := time.Now()
		for n := 0; n < len(data); n += max(len(data)/97, 1) {
			Convert(bytes.NewReader(data[:n]), int64(n), &Options{NoSystemFonts: true, NoTextIndex: true, Warn: func(string) {}})
		}
		b := slices.Clone(data)
		for i := 7; i < len(b); i += max(len(b)/61, 1) {
			b[i] ^= 0xa5
			Convert(bytes.NewReader(b), int64(len(b)), &Options{NoSystemFonts: true, NoTextIndex: true, Warn: func(string) {}})
		}
		if d := time.Since(start); d > 30*time.Second {
			t.Errorf("%s: %v", name, d)
		}
	}
}
