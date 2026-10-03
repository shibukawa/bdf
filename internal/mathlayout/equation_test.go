package mathlayout

import (
	"math"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/internal/fontdb"
	"github.com/shibukawa/bdf/internal/xmltree"
	"golang.org/x/net/html"
)

func TestTeXLinear(t *testing.T) {
	for _, tc := range []struct{ tex, want string }{
		{`x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}`, "x=(−b±√(b^2−4ac))/(2a)"},
		{`\sum_{k=1}^{n} k^2`, "∑_(k=1)^n k^2"},
		{`\sin\theta + \cos^2 x`, "sin θ+cos^2 x"},
		{`\lim_{h \to 0} f(h)`, "lim_(h→0) f(h)"},
		{`\left( \frac{a}{b} \right)`, "(a/b)"},
		{`\begin{pmatrix} a & b \\ c & d \end{pmatrix}`, "(a, b; c, d)"},
		{`\hat{x} + \vec{v}`, "x̂+v⃗"},
		{`\mathbb{R}^n`, "ℝ^n"},
		{`\text{if } x > 0`, "if x>0"},
		{`\sqrt[3]{x}`, "√(3&x)"},
		{`a \not= b \ne c`, "a≠b≠c"},
		{`f'(x)`, "f^′(x)"},
		{`\binom{n}{k}`, "((n¦k))"},
		{`\int_0^\infty e^{-x}\,dx`, "∫_0^∞ e^(−x) dx"},
		{`\begin{aligned} a &= b \\ &= c \end{aligned}`, "a=b\n=c"},
		{`\unknowncommand{x}`, `\unknowncommandx`},
	} {
		if got := Linear(ParseTeX(tc.tex)); got != tc.want {
			t.Errorf("%s: linear %q, want %q", tc.tex, got, tc.want)
		}
	}
}

func TestMathMLLinear(t *testing.T) {
	for _, tc := range []struct{ src, want string }{
		{`<math><mfrac><mrow><mi>a</mi><mo>+</mo><mi>b</mi></mrow><mn>2</mn></mfrac></math>`, "(a+b)/2"},
		{`<math display="block"><msubsup><mo>∫</mo><mn>0</mn><mn>1</mn></msubsup><mi>x</mi><mi>d</mi><mi>x</mi></math>`, "∫_0^1 xdx"},
		{`<math><msqrt><mi>x</mi></msqrt><mo>+</mo><mroot><mi>y</mi><mn>3</mn></mroot></math>`, "√x+√(3&y)"},
		{`<math><mover accent="true"><mi>x</mi><mo>^</mo></mover></math>`, "x̂"},
		{`<math><mfenced><mi>a</mi><mi>b</mi></mfenced></math>`, "(a,b)"},
		{`<math><semantics><mrow><mi>x</mi></mrow><annotation encoding="application/x-tex">x</annotation></semantics></math>`, "x"},
		{`<math><semantics><annotation encoding="application/x-tex">\frac{1}{2}</annotation></semantics></math>`, "1/2"},
		{`<math><mtable><mtr><mtd><mi>a</mi></mtd><mtd><mi>b</mi></mtd></mtr></mtable></math>`, "a, b"},
		{`<math><mi>sin</mi><mo>&#x2061;</mo><mi>x</mi></math>`, "sinx"},
	} {
		doc, err := html.Parse(strings.NewReader(tc.src))
		if err != nil {
			t.Fatal(err)
		}
		m := findElement(doc, "math")
		n, display := ParseMathML(m)
		if got := Linear(n); got != tc.want {
			t.Errorf("%s: linear %q, want %q", tc.src, got, tc.want)
		}
		if want := strings.Contains(tc.src, `display="block"`); display != want {
			t.Errorf("%s: display %v", tc.src, display)
		}
	}
}

func findElement(n *html.Node, name string) *html.Node {
	if n.Type == html.ElementNode && n.Data == name {
		return n
	}
	for k := n.FirstChild; k != nil; k = k.NextSibling {
		if f := findElement(k, name); f != nil {
			return f
		}
	}
	return nil
}

const mns = `xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math" xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"`

func TestOMMLLinear(t *testing.T) {
	for _, tc := range []struct{ xml, want string }{
		{`<m:oMath><m:r><m:t>x=</m:t></m:r><m:f><m:num><m:r><m:t>a+b</m:t></m:r></m:num><m:den><m:r><m:t>2</m:t></m:r></m:den></m:f></m:oMath>`, "x=(a+b)/2"},
		{`<m:oMath><m:nary><m:naryPr><m:chr m:val="∑"/></m:naryPr><m:sub><m:r><m:t>i=1</m:t></m:r></m:sub><m:sup><m:r><m:t>n</m:t></m:r></m:sup><m:e><m:sSub><m:e><m:r><m:t>x</m:t></m:r></m:e><m:sub><m:r><m:t>i</m:t></m:r></m:sub></m:sSub></m:e></m:nary></m:oMath>`, "∑_(i=1)^n x_i"},
		{`<m:oMath><m:nary><m:e><m:r><m:t>f</m:t></m:r></m:e></m:nary></m:oMath>`, "∫f"},
		{`<m:oMath><m:rad><m:radPr><m:degHide m:val="1"/></m:radPr><m:deg/><m:e><m:r><m:t>x+1</m:t></m:r></m:e></m:rad></m:oMath>`, "√(x+1)"},
		{`<m:oMath><m:d><m:dPr><m:begChr m:val="["/><m:endChr m:val="]"/></m:dPr><m:e><m:r><m:t>a</m:t></m:r></m:e><m:e><m:r><m:t>b</m:t></m:r></m:e></m:d></m:oMath>`, "[a|b]"},
		{`<m:oMath><m:func><m:fName><m:r><m:rPr><m:sty m:val="p"/></m:rPr><m:t>sin</m:t></m:r></m:fName><m:e><m:r><m:t>x</m:t></m:r></m:e></m:func></m:oMath>`, "sin x"},
		{`<m:oMath><m:acc><m:e><m:r><m:t>x</m:t></m:r></m:e></m:acc></m:oMath>`, "x̂"},
		{`<m:oMath><m:eqArr><m:e><m:r><m:t>a&amp;=b</m:t></m:r></m:e><m:e><m:r><m:t>&amp;=c</m:t></m:r></m:e></m:eqArr></m:oMath>`, "a=b\n=c"},
		{`<m:oMath><m:r><m:rPr><m:nor/></m:rPr><m:t>速さ</m:t></m:r></m:oMath>`, "速さ"},
	} {
		root, err := xmltree.Parse([]byte(`<root ` + mns + `>` + tc.xml + `</root>`))
		if err != nil {
			t.Fatal(err)
		}
		n := (&OMML{}).ParseOMML(root.Child("oMath"))
		if got := Linear(n); got != tc.want {
			t.Errorf("%s: linear %q, want %q", tc.xml, got, tc.want)
		}
	}
}

// testEngine returns an engine with the test fonts (STIX Two Math and M
// PLUS 1p), or with M PLUS 1p alone. Like the converters, it draws the
// glyphs no character maps to with private use characters.
func testEngine(t *testing.T, withMath bool) *Engine {
	t.Helper()
	db := fontdb.New(nil, []string{"../../converter/docx/testdata/fonts"}, false)
	faces := map[*fontdb.Face]*Face{}
	face := func(f *fontdb.Face) *Face {
		if f == nil {
			return nil
		}
		if fa := faces[f]; fa != nil {
			return fa
		}
		l, err := f.Load()
		if err != nil {
			t.Fatal(err)
		}
		faces[f] = &Face{Loaded: l}
		return faces[f]
	}
	fonts := Fonts{
		Text: func(family string, r rune, bold, italic bool) *Face {
			return face(db.Resolve("M PLUS 1p", bold, italic, false).Face)
		},
		GlyphRune: func(f *Face, g uint16) (rune, bool) { return rune(0xF0000 + int(g)), true },
	}
	if withMath {
		fonts.Math = face(db.ResolveMath("Cambria Math"))
	}
	return New(fonts)
}

func layout(e *Engine, tex string, display bool) *Box {
	return e.Layout(ParseTeX(tex), Style{Size: 10, Display: display})
}

func TestLayout(t *testing.T) {
	e := testEngine(t, true)
	if !e.HasMathFont() {
		t.Fatal("the test font STIX Two Math is not found")
	}
	x := layout(e, `x`, false)
	if x.W <= 0 || x.H <= 0 || x.H > 6 {
		t.Errorf("x: %v × %v", x.W, x.H)
	}
	frac := layout(e, `\frac{a}{b}`, true)
	if frac.H <= x.H || frac.D <= 0 {
		t.Errorf("fraction: height %v depth %v", frac.H, frac.D)
	}
	sup := layout(e, `x^2`, false)
	var two *item
	for i := range sup.items {
		if sup.items[i].text == "2" {
			two = &sup.items[i]
		}
	}
	if two == nil || two.y >= -2 || two.size >= 10 {
		t.Errorf("superscript: %+v", two)
	}
	// large operators grow in display style
	if in, disp := layout(e, `\sum`, false), layout(e, `\sum`, true); disp.H+disp.D <= (in.H+in.D)*1.1 {
		t.Errorf("display ∑ %v is not larger than inline %v", disp.H+disp.D, in.H+in.D)
	}
	// delimiters grow to what they enclose, with the font's variants
	tall := layout(e, `\left( \frac{\frac{a}{b}}{\frac{c}{d}} \right)`, true)
	pua := false
	for _, it := range tall.items {
		for _, r := range it.text {
			if r >= 0xF0000 {
				pua = true
				if it.alt != "(" && it.alt != ")" {
					t.Errorf("private use glyph for %q", it.alt)
				}
			}
		}
	}
	if !pua {
		t.Error("the parentheses around a tall fraction did not grow")
	}
	// a radical is as tall as its radicand and a gap
	r := layout(e, `\sqrt{\frac{a}{b}}`, true)
	f := layout(e, `\frac{a}{b}`, true)
	if r.H <= f.H || r.W <= f.W {
		t.Errorf("radical %v×%v around %v×%v", r.W, r.H, f.W, f.H)
	}
	// every node lays out, whatever it holds
	for _, s := range sampleTeX {
		b := layout(e, s, true)
		if b.W <= 0 || math.IsNaN(b.W+b.H+b.D) {
			t.Errorf("%s: %v × %v + %v", s, b.W, b.H, b.D)
		}
	}
}

func TestLayoutWithoutMathFont(t *testing.T) {
	e := testEngine(t, false)
	if e.HasMathFont() {
		t.Fatal("a formula font was found")
	}
	for _, s := range sampleTeX {
		b := layout(e, s, true)
		if b.W <= 0 || math.IsNaN(b.W+b.H+b.D) {
			t.Errorf("%s: %v × %v + %v", s, b.W, b.H, b.D)
		}
	}
	r := layout(e, `\sqrt{x}`, false)
	strokes := 0
	for _, it := range r.items {
		if it.kind == iStroke {
			strokes++
		}
	}
	if strokes != 1 {
		t.Errorf("radical without a formula font: %d strokes", strokes)
	}
}

// sampleTeX are formulas of every kind of node (the document of samples of
// converter/internal/equation draws the same).
var sampleTeX = []string{
	`x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}`,
	`e^{i\pi} + 1 = 0,\quad f'(x) = \lim_{h \to 0} \frac{f(x+h) - f(x)}{h}`,
	`\sum_{k=1}^{n} k^2 = \frac{n(n+1)(2n+1)}{6}`,
	`\int_0^\infty e^{-x^2}\,dx = \frac{\sqrt{\pi}}{2}`,
	`\left( \frac{a}{b} \right)^2 \left[ x_1, x_2 \right] \left\{ \frac{1}{2} \right\} \left| \frac{x}{y} \right|`,
	`A = \begin{pmatrix} a_{11} & a_{12} \\ a_{21} & a_{22} \end{pmatrix},\quad \det A = \begin{vmatrix} a & b \\ c & d \end{vmatrix}`,
	`f(x) = \begin{cases} x^2 & x \ge 0 \\ -x & x < 0 \end{cases}`,
	`\hat{x} + \tilde{y} + \bar{z} + \vec{v} + \dot{a} + \ddot{b} + \widehat{xyz} + \overline{AB} + \underline{cd}`,
	`\sqrt[3]{x+1} + \sqrt{\frac{a}{b}} + \sqrt{\sqrt{x}}`,
	`\overbrace{a+b+c}^{3} + \underbrace{x+y}_{2} \xrightarrow{f} \mathbb{R}^n \mathcal{L} \mathfrak{g} \mathbf{v}`,
	`\begin{aligned} (a+b)^2 &= a^2 + 2ab + b^2 \\ &\le 2(a^2 + b^2) \end{aligned}`,
	`\Gamma(z) = \int_0^\infty t^{z-1} e^{-t}\, dt,\quad \binom{n}{k} = \frac{n!}{k!(n-k)!}`,
	`\bigl( \Bigl( \biggl( \Biggl( x \Biggr) \biggr) \Bigr) \bigr) \quad a \bmod b \pmod{n}`,
	`\prod_{i=1}^{n} x_i \bigcup_{\alpha \in A} U_\alpha \oint_C \mathbf{F} \cdot d\mathbf{r}`,
	`\text{if } x \in \mathbb{Z}\text{ then 日本語も}\ \color{red}{x^2} \boxed{E = mc^2} \cancel{y}`,
	`{}_{n}C_{r} \quad \nabla \times \mathbf{B} = \mu_0 \mathbf{J} \quad \alpha\beta\gamma\Delta\Omega`,
}
