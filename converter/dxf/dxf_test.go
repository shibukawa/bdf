package dxf

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// The test drawings are made by test/dxf/gen.py (ezdxf); text is laid out
// with the test fonts of the PowerPoint converter.

func testOptions() *Options {
	return &Options{FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true}
}

func convertData(t *testing.T, data []byte, opts *Options) (*Result, *bdf.Reader) {
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

func convert(t *testing.T, name string) (*Result, *bdf.Reader) {
	t.Helper()
	res, err := ConvertFile(filepath.Join("testdata", name), testOptions())
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

func plainText(t *testing.T, r *bdf.Reader, view int) string {
	t.Helper()
	h, err := bdf.ParseHash(r.Manifest.Views[view].TextIndex)
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

func TestShapes(t *testing.T) {
	res, r := convert(t, "shapes.dxf")
	if len(res.Warnings) > 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	m := r.Manifest
	if len(m.Views) != 1 || m.Views[0].ID != "model" || m.Views[0].Title != "Model" || len(m.Views[0].Pages) != 1 {
		t.Fatalf("views: %+v", m.Views)
	}
	if m.Meta.Source != "dxf" {
		t.Errorf("source %q", m.Meta.Source)
	}
	text := plainText(t, r, 0)
	for _, want := range []string{"LEFT 字", "CENTER", "FIT テキスト", "ALIGNED", "Ø50 ±0.1 45° UNDER", "カイテン 30°",
		"図の文字は", "おおきく、アンダーライン、Bold", "長い文字は", "CENTERED", "MTEXT", "60", "R12", "MULTI", "LEADER", "A1", "A3"} {
		if !strings.Contains(text, want) {
			t.Errorf("text lacks %q:\n%s", want, text)
		}
	}
	for _, hidden := range []string{"OFF LAYER", "FROZEN LAYER", "DEFPOINTS"} {
		if strings.Contains(text, hidden) {
			t.Errorf("text of a hidden layer %q is drawn", hidden)
		}
	}
	// the page is A3 wide and fits the drawing (300 × 200 mm and margins)
	p := m.Views[0].Pages[0]
	if p.W < 1190 || p.W > 1191 || p.H < 760 || p.H > 830 {
		t.Errorf("page %v × %v", p.W, p.H)
	}
}

// TestBinary checks that a binary DXF file converts to the same objects as
// the same drawing in text.
func TestBinary(t *testing.T) {
	_, a := convert(t, "shapes.dxf")
	_, b := convert(t, "shapes-bin.dxf")
	la, lb := a.Manifest.Views[0].Pages[0].Layers, b.Manifest.Views[0].Pages[0].Layers
	if len(la) != 1 || len(lb) != 1 || la[0].Obj != lb[0].Obj {
		t.Errorf("binary drawing differs: %v vs %v", la, lb)
	}
}

func TestLayout(t *testing.T) {
	res, r := convert(t, "layout.dxf")
	if len(res.Warnings) > 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	m := r.Manifest
	if len(m.Views) != 2 || m.Views[1].Title != "A3 図面" || m.Views[1].ID != "layout1" {
		t.Fatalf("views: %+v", m.Views)
	}
	// A3 landscape
	p := m.Views[1].Pages[0]
	if abs(float64(p.W)-420*72/25.4) > 0.01 || abs(float64(p.H)-297*72/25.4) > 0.01 {
		t.Errorf("paper %v × %v", p.W, p.H)
	}
	text := plainText(t, r, 1)
	if !strings.Contains(text, "テスト図 1:100 / 1:200") {
		t.Errorf("title block missing:\n%s", text)
	}
	// both viewports show the plan; the second one freezes a layer
	if n := strings.Count(text, "区間 A"); n != 2 {
		t.Errorf("%d viewports show the plan, want 2:\n%s", n, text)
	}
	if n := strings.Count(text, "FROZEN IN VIEWPORT 2"); n != 1 {
		t.Errorf("the frozen layer shows %d times, want 1:\n%s", n, text)
	}
	// only one view each when asked for
	for _, c := range []struct {
		views string
		want  []string
	}{{"model", []string{"model"}}, {"layouts", []string{"layout1"}}} {
		opts := testOptions()
		opts.Views = c.views
		res, err := ConvertFile("testdata/layout.dxf", opts)
		if err != nil {
			t.Fatal(err)
		}
		var ids []string
		for _, v := range res.Doc.Views {
			ids = append(ids, v.ID)
		}
		if strings.Join(ids, ",") != strings.Join(c.want, ",") {
			t.Errorf("views=%s: %v", c.views, ids)
		}
	}
}

func TestShiftJIS(t *testing.T) {
	_, r := convert(t, "r12-sjis.dxf")
	if text := plainText(t, r, 0); !strings.Contains(text, "シフトJIS の文字") {
		t.Errorf("text: %q", text)
	}
	if got := r.Manifest.Meta.DC.Language.First(); got != "ja" {
		t.Errorf("language %q", got)
	}
}

// minimal returns a DXF file in text with a header and one TEXT entity.
func minimal(codepage string, text []byte) []byte {
	var b bytes.Buffer
	b.WriteString("0\nSECTION\n2\nHEADER\n9\n$ACADVER\n1\nAC1009\n")
	if codepage != "" {
		b.WriteString("9\n$DWGCODEPAGE\n3\n" + codepage + "\n")
	}
	b.WriteString("0\nENDSEC\n0\nSECTION\n2\nENTITIES\n0\nTEXT\n8\n0\n10\n0\n20\n0\n40\n5\n1\n")
	b.Write(text)
	b.WriteString("\n0\nENDSEC\n0\nEOF\n")
	return b.Bytes()
}

func TestEncodings(t *testing.T) {
	sjis := []byte{0x83, 0x65, 0x83, 0x58, 0x83, 0x67, 0x82, 0xcc, 0x95, 0xb6, 0x8e, 0x9a} // テストの文字
	for _, c := range []struct {
		name     string
		codepage string
		text     []byte
		want     string
	}{
		{"shift_jis", "ANSI_932", sjis, "テストの文字"},
		{"shift_jis labelled as the default", "ANSI_1252", sjis, "テストの文字"},
		{"shift_jis without a code page", "", sjis, "テストの文字"},
		{"latin", "ANSI_1252", []byte("Gr\xf6\xdfe"), "Größe"},
		{"utf-8", "ANSI_1252", []byte("Größe 字"), "Größe 字"},
		{"unicode escape", "ANSI_1252", []byte(`\U+5B57\U+00B0`), "字°"},
		{"multibyte escape", "ANSI_1252", []byte(`\M+18E9A`), "字"},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, r := convertData(t, minimal(c.codepage, c.text), nil)
			if got := strings.TrimSpace(plainText(t, r, 0)); got != c.want {
				t.Errorf("text %q, want %q", got, c.want)
			}
		})
	}
}

func TestDetect(t *testing.T) {
	for _, c := range []struct {
		head string
		want bool
	}{
		{"  0\r\nSECTION\r\n  2\r\nHEADER\r\n", true},
		{"999\ncomment\n0\nSECTION\n", true},
		{"\xef\xbb\xbf0\nSECTION\n", true},
		{"AutoCAD Binary DXF\r\n\x1a\x00\x00\x00", true},
		{"0\nLINE\n", false},
		{"hello\nworld\n", false},
	} {
		if got := Detect([]byte(c.head)); got != c.want {
			t.Errorf("Detect(%q) = %v", c.head, got)
		}
	}
}

func TestTextCodes(t *testing.T) {
	segs := textCodes("%%c10 %%p1%%d %%uA%%u%%%B %%065")
	var b strings.Builder
	under := ""
	for _, s := range segs {
		b.WriteString(s.s)
		if s.underline {
			under += s.s
		}
	}
	if got := b.String(); got != "Ø10 ±1° A%B A" {
		t.Errorf("text %q", got)
	}
	if under != "A" {
		t.Errorf("underlined %q", under)
	}
}

func TestSplitStack(t *testing.T) {
	for _, c := range []struct{ in, num, den string }{
		{"1/2", "1", "2"}, {"+0.1^ -0.2", "+0.1", "-0.2"}, {"a#b", "a", "b"}, {`1\/2/3`, "1/2", "3"},
	} {
		num, den, _ := splitStack(c.in)
		if num != c.num || den != c.den {
			t.Errorf("splitStack(%q) = %q, %q", c.in, num, den)
		}
	}
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
