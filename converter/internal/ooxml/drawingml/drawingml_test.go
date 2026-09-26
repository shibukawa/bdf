package drawingml

import (
	"math"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

func TestRunLang(t *testing.T) {
	st := func(lang, alt string) *runStyle { return &runStyle{lang: lang, altLang: alt} }
	run := func(s string, st *runStyle) []item {
		var out []item
		for _, r := range s {
			out = append(out, item{r: r, st: st, kind: itemChar})
		}
		return out
	}
	kana := &para{eaLang: "ja"}
	for _, c := range []struct {
		name  string
		run   []item
		pa    *para
		want  string
		known bool
	}{
		{"latin", run("Hello", st("en-US", "ja-JP")), nil, "en-US", true},
		{"latin without lang", run("Hello", st("", "")), nil, "", true},
		{"kana", run("かな", st("en-US", "")), nil, "ja", true},
		{"hangul", run("한국어", st("en-US", "")), nil, "ko", true},
		{"East Asian lang", run("汉字", st("zh-CN", "en-US")), nil, "zh-CN", true},
		{"East Asian altLang", run("漢字", st("en-US", "ja-JP")), nil, "ja-JP", true},
		{"altLang before script", run("かな", st("en-US", "zh-TW")), nil, "zh-TW", true},
		{"Han in a kana paragraph", run("漢字", st("en-US", "")), kana, "ja", true},
		{"Han alone", run("漢字", st("en-US", "")), &para{}, "", false},
		{"mixed", append(run("ABC", st("en-US", "")), run("です", st("en-US", ""))...), nil, "ja", true},
	} {
		if got, known := runLang(c.run, c.pa); got != c.want || known != c.known {
			t.Errorf("%s: runLang = %q, %v; want %q, %v", c.name, got, known, c.want, c.known)
		}
	}
	for _, c := range []struct {
		r, c int
		want string
	}{{0, 0, "A1"}, {3, 1, "B4"}, {0, 25, "Z1"}, {0, 26, "AA1"}, {9, 701, "ZZ10"}, {0, 702, "AAA1"}} {
		if got := cellRef(c.r, c.c); got != c.want {
			t.Errorf("cellRef(%d, %d) = %s, want %s", c.r, c.c, got, c.want)
		}
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
		n, err := ooxml.Parse([]byte(`<solidFill xmlns:a="a">` + x + `</solidFill>`))
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
