package pptx

import (
	"bytes"
	"math"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// testOptions restricts fonts to the test font directory so that output
// does not depend on the machine.
func testOptions() *Options {
	return &Options{FontDirs: []string{"testdata/fonts"}, NoSystemFonts: true}
}

func convert(t *testing.T, name string, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
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

// opCounts counts the instructions of an object and the objects it uses.
func opCounts(t *testing.T, r *bdf.Reader, h bdf.Hash) map[byte]int {
	t.Helper()
	counts := map[byte]int{}
	var walk func(h bdf.Hash)
	walk = func(h bdf.Hash) {
		o, err := r.Object(h)
		if err != nil {
			t.Fatal(err)
		}
		if err := o.Walk(func(in bdf.Instr) { counts[in.Op]++ }); err != nil {
			t.Fatal(err)
		}
		for _, c := range o.Objects {
			walk(c)
		}
	}
	walk(h)
	return counts
}

func markCounts(t *testing.T, r *bdf.Reader, h bdf.Hash) map[byte]int {
	t.Helper()
	counts := map[byte]int{}
	o, err := r.Object(h)
	if err != nil {
		t.Fatal(err)
	}
	o.Walk(func(in bdf.Instr) {
		if in.Op == bdf.OpMark {
			counts[byte(in.Args[0].(uint64))]++
		}
	})
	return counts
}

func body(p *bdf.Page) bdf.Hash {
	for _, l := range p.Layers {
		if l.Role == bdf.RoleBody {
			return l.Obj
		}
	}
	return bdf.Hash{}
}

func TestConvertBasic(t *testing.T) {
	res, r := convert(t, "basic.pptx", testOptions())
	if res.Slides != 5 {
		t.Fatalf("slides = %d (the hidden slide must be skipped)", res.Slides)
	}
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	if res.EmbeddedFonts != 2 {
		t.Errorf("embedded fonts = %d, want the regular and bold test fonts", res.EmbeddedFonts)
	}
	v := r.Manifest.Views[0]
	if v.Kind != bdf.ViewFixed || v.Pages[0].W != 720 || v.Pages[0].H != 540 {
		t.Fatalf("view = %s %gx%g", v.Kind, v.Pages[0].W, v.Pages[0].H)
	}
	if r.Manifest.Meta.Title != "PowerPoint test deck" || r.Manifest.Meta.Source != "pptx" {
		t.Errorf("meta = %+v", r.Manifest.Meta)
	}
	text := plainText(t, r)
	for _, want := range []string{
		"BDF from PowerPoint\nSlides rendered directly from DrawingML",
		"• First level bullet with enough words to wrap onto a second line inside the placeholder",
		"– Second level", "1. numbered one\n2. numbered two\n3. numbered three",
		"• A link to example.com",
		"Aligned justify: the quick brown fox jumps over the lazy dog and keeps running.",
		// wrapped East Asian lines join without a space
		"日本語の文章は、単語の区切りに空白を使わないため、文字と文字の間で改行します。句読点「、」や「。」は行頭に来ません。",
		"混在 mixed テキスト x2 underline", "縦書きのテキスト",
		"Header 0", "merged cell with longer text that wraps",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("text lacks %q", want)
		}
	}
	if strings.Contains(text, "Hidden slide") {
		t.Error("hidden slide converted")
	}
	if c := opCounts(t, r, body(v.Pages[1])); c[bdf.OpLink] != 1 {
		t.Errorf("slide 2 links = %d", c[bdf.OpLink])
	}
	if m := markCounts(t, r, body(v.Pages[3])); m[bdf.MarkWrap] == 0 || m[bdf.MarkLine] == 0 || m[bdf.MarkBox] == 0 {
		t.Errorf("slide 4 marks = %v", m)
	}
	// the picture and a table
	if c := opCounts(t, r, body(v.Pages[4])); c[bdf.OpImage]+c[bdf.OpImageSub] != 1 || c[bdf.OpStrokePath] < 10 {
		t.Errorf("slide 5 ops = %v", c)
	}
}

func TestSharedLayers(t *testing.T) {
	res, r := convert(t, "features.pptx", testOptions())
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	v := r.Manifest.Views[0]
	if v.Pages[0].W != 960 || v.Pages[0].H != 540 {
		t.Fatalf("page = %gx%g", v.Pages[0].W, v.Pages[0].H)
	}
	a, b := v.Pages[0].Layers, v.Pages[1].Layers
	roles := func(ls []bdf.Layer) string {
		var s []string
		for _, l := range ls {
			s = append(s, l.Role)
		}
		return strings.Join(s, ",")
	}
	if roles(a) != "background,master,master,body" || roles(b) != roles(a) {
		t.Fatalf("layers = %s / %s", roles(a), roles(b))
	}
	for i := 0; i < 3; i++ {
		if a[i].Obj != b[i].Obj {
			t.Errorf("layer %d (%s) is not shared", i, a[i].Role)
		}
	}
	if a[3].Obj == b[3].Obj {
		t.Error("bodies should differ")
	}
	// slide 5 (Title Only layout) shares the master layer but not the layout one
	c := v.Pages[4].Layers
	if c[1].Obj != a[1].Obj || c[2].Obj == a[2].Obj {
		t.Errorf("master/layout sharing across layouts is wrong")
	}
	text := plainText(t, r)
	for _, want := range []string{"BDF features deck", "Share", "Q1", "East", "45%", "日本語のスライド",
		"箇条書きの一行目です。長い文章は枠の幅で折り返され、句読点は行頭に来ません。", "縦書きで、読みます。", "プリセット形状", "Callout"} {
		if !strings.Contains(text, want) {
			t.Errorf("text lacks %q", want)
		}
	}
}

func TestOptions(t *testing.T) {
	opts := testOptions()
	opts.Hidden = true
	res, r := convert(t, "basic.pptx", opts)
	if res.Slides != 6 || !strings.Contains(plainText(t, r), "Hidden slide") {
		t.Errorf("Hidden: %d slides", res.Slides)
	}
	opts = testOptions()
	opts.Slides = []int{2, 4}
	res, _ = convert(t, "basic.pptx", opts)
	if res.Slides != 2 {
		t.Errorf("Slides: %d", res.Slides)
	}
	opts.Slides = []int{9}
	if _, err := ConvertFile("testdata/basic.pptx", opts); err == nil {
		t.Error("out of range slide accepted")
	}
	opts = testOptions()
	opts.SystemFonts = true
	res, r = convert(t, "basic.pptx", opts)
	if res.EmbeddedFonts != 0 {
		t.Errorf("SystemFonts embedded %d fonts", res.EmbeddedFonts)
	}
	for _, p := range r.Manifest.Parts {
		if p.T == bdf.PartFont {
			t.Fatal("font part with SystemFonts")
		}
	}
	o, err := r.Object(body(r.Manifest.Views[0].Pages[0]))
	if err != nil {
		t.Fatal(err)
	}
	if len(o.Fonts) == 0 || o.Fonts[0].Kind != bdf.FontSystem || !strings.Contains(o.Fonts[0].Family, `"Calibri"`) {
		t.Errorf("fonts = %+v", o.Fonts)
	}
}

func TestDeterministic(t *testing.T) {
	var outs [2][]byte
	for i := range outs {
		res, err := ConvertFile("testdata/features.pptx", testOptions())
		if err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		if err := res.Doc.WriteSingle(&buf); err != nil {
			t.Fatal(err)
		}
		outs[i] = buf.Bytes()
	}
	if !bytes.Equal(outs[0], outs[1]) {
		t.Fatal("conversion is not deterministic")
	}
}

func TestLineBreaking(t *testing.T) {
	for _, c := range []struct {
		a, b rune
		want bool
	}{
		{' ', 'a', true}, {'a', ' ', false}, {'a', 'b', false}, {'-', 'b', true},
		{'日', '本', true}, {'本', '。', false}, {'「', '本', false}, {'本', 'A', true}, {'ー', 'ト', true}, {'テ', 'ー', false},
	} {
		if got := canBreak(c.a, c.b); got != c.want {
			t.Errorf("canBreak(%q, %q) = %v", c.a, c.b, got)
		}
	}
	items := func(s string) []item {
		var out []item
		for _, r := range s {
			out = append(out, item{r: r, w: 10, kind: itemChar})
		}
		return out
	}
	lines := func(pa *para, width float64) []string {
		var out []string
		for _, ln := range layoutParagraph(pa, width, true) {
			var b strings.Builder
			for _, it := range ln.items {
				b.WriteRune(it.r)
			}
			out = append(out, b.String())
		}
		return out
	}
	pa := &para{items: items("aaa bbb ccc"), defTab: 72}
	if got := strings.Join(lines(pa, 75), "|"); got != "aaa bbb |ccc" {
		t.Errorf("latin = %q", got)
	}
	pa = &para{items: items("あいうえ。お"), defTab: 72}
	if got := strings.Join(lines(pa, 40), "|"); got != "あいう|え。お" {
		t.Errorf("kinsoku = %q", got)
	}
	// a word longer than the line is cut
	pa = &para{items: items("abcdefgh"), defTab: 72}
	if got := strings.Join(lines(pa, 35), "|"); got != "abc|def|gh" {
		t.Errorf("long word = %q", got)
	}
	// hanging indent with a bullet
	pa = &para{items: items("ab"), defTab: 72, marL: 20, indent: -20, bullet: &bullet{w: 8}}
	if ls := layoutParagraph(pa, 100, true); ls[0].start != 20 || ls[0].bulletX != 0 {
		t.Errorf("bullet line start %g, bullet at %g", ls[0].start, ls[0].bulletX)
	}
}

func TestAutoNumber(t *testing.T) {
	for _, c := range []struct {
		scheme string
		n      int
		want   string
	}{
		{"arabicPeriod", 3, "3."}, {"arabicParenR", 2, "2)"}, {"romanUcPeriod", 14, "XIV."}, {"romanLcParenBoth", 4, "(iv)"},
		{"alphaLcParenR", 28, "bb)"}, {"alphaUcPeriod", 1, "A."}, {"circleNumDbPlain", 3, "③"}, {"arabicDbPeriod", 12, "１２．"},
		{"ea1JpnChsDbPeriod", 11, "十一．"},
	} {
		if got := autoNumber(c.scheme, c.n); got != c.want {
			t.Errorf("autoNumber(%s, %d) = %q, want %q", c.scheme, c.n, got, c.want)
		}
	}
}

func TestColors(t *testing.T) {
	cc := &colorCtx{scheme: defaultTheme.colors, clrMap: defaultClrMap}
	parse := func(x string) rgba {
		n, err := parseXML([]byte(`<solidFill xmlns:a="a">` + x + `</solidFill>`))
		if err != nil {
			t.Fatal(err)
		}
		c, ok := cc.color(n)
		if !ok {
			t.Fatalf("no color in %s", x)
		}
		return c
	}
	near := func(c rgba, r, g, b uint8) bool {
		d := func(v float64, w uint8) bool { return math.Abs(v*255-float64(w)) <= 1.5 }
		return d(c.R, r) && d(c.G, g) && d(c.B, b)
	}
	// "Blue, Accent 1, Lighter 40%" and "Darker 50%" as Office shows them
	if c := parse(`<schemeClr val="accent1"><lumMod val="60000"/><lumOff val="40000"/></schemeClr>`); !near(c, 0x8f, 0xaa, 0xdc) {
		t.Errorf("lumMod/lumOff = %v", c.bdf().CSS())
	}
	if c := parse(`<schemeClr val="accent1"><lumMod val="50000"/></schemeClr>`); !near(c, 0x1f, 0x38, 0x64) {
		t.Errorf("lumMod = %v", c.bdf().CSS())
	}
	if c := parse(`<schemeClr val="tx1"/>`); !near(c, 0, 0, 0) {
		t.Errorf("tx1 = %v", c.bdf().CSS())
	}
	if c := parse(`<srgbClr val="FF0000"><alpha val="50000"/></srgbClr>`); c.A != 0.5 {
		t.Errorf("alpha = %v", c.A)
	}
	if c := parse(`<prstClr val="dkBlue"/>`); !near(c, 0, 0, 0x8b) {
		t.Errorf("prstClr = %v", c.bdf().CSS())
	}
	if c := parse(`<srgbClr val="4472C4"><shade val="50000"/></srgbClr>`); c.R >= 0x44/255.0 || c.B >= 0xc4/255.0 {
		t.Errorf("shade = %v", c.bdf().CSS())
	}
}

func TestCharts(t *testing.T) {
	for _, c := range []struct {
		v    float64
		code string
		want string
	}{
		{8.200000000000001, "General", "8.2"}, {0.45, "0%", "45%"}, {1234567, "#,##0", "1,234,567"}, {3.14159, "0.00", "3.14"},
		{-5, "General", "-5"}, {12, `"$"#,##0`, "$12"},
	} {
		if got := formatNumber(c.v, c.code); got != c.want {
			t.Errorf("formatNumber(%v, %q) = %q, want %q", c.v, c.code, got, c.want)
		}
	}
	a := niceAxis(0, 8.2, nil, 10)
	if a.lo != 0 || a.hi != 9 || a.step != 1 {
		t.Errorf("axis 0..8.2 = %v..%v by %v", a.lo, a.hi, a.step)
	}
	a = niceAxis(20.4, 90, nil, 6)
	if a.lo != 0 || a.hi != 100 || a.step != 20 {
		t.Errorf("axis 20.4..90 = %v..%v by %v", a.lo, a.hi, a.step)
	}
	a = niceAxis(95, 100, nil, 10)
	if a.lo == 0 {
		t.Errorf("axis 95..100 starts at zero")
	}
}
