package equation

import (
	"math"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/internal/fontdb"
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
		root, err := ooxml.Parse([]byte(`<root ` + mns + `>` + tc.xml + `</root>`))
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
// PLUS 1p), or with M PLUS 1p alone.
func testEngine(t *testing.T, withMath bool) (*Engine, *fontset.Set) {
	t.Helper()
	dir := "../../docx/testdata/fonts"
	var fsys []string
	if withMath {
		fsys = []string{dir}
	}
	db := fontdb.New(nil, fsys, false)
	if !withMath {
		db = fontdb.New(nil, nil, false)
		for _, f := range fontdb.New(nil, []string{dir}, false).Faces {
			if !f.Math {
				db.Faces = append(db.Faces, f)
			}
		}
	}
	set := fontset.New(db, nil)
	text := func(family string, r rune, bold, italic bool) *fontset.Choice {
		return set.FaceFor(family, "", bold, italic, r)
	}
	return New(Fonts{Set: set, Math: "Cambria Math", Text: text}), set
}

func layout(e *Engine, tex string, display bool) *Box {
	return e.Layout(ParseTeX(tex), Style{Size: 10, Display: display})
}

func TestLayout(t *testing.T) {
	e, set := testEngine(t, true)
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
	// the glyphs drawn by index are embedded
	doc := bdf.NewDocument()
	if n := set.Embed(doc, fontset.EmbedOptions{PinnedOnly: true}); n != 1 {
		t.Errorf("embedded %d fonts with PinnedOnly, want the formula font", n)
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
	e, _ := testEngine(t, false)
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

func TestPlace(t *testing.T) {
	e, set := testEngine(t, true)
	doc := bdf.NewDocument()
	cvs := canvas.NewBuilder(doc, set)
	cv := cvs.New()
	n := ParseTeX(`E = mc^2`)
	Place(cv, e.Layout(n, Style{Size: 12}), 10, 20, Linear(n))
	set.Embed(doc, fontset.EmbedOptions{})
	cvs.Encode()
	if !cv.Drawn {
		t.Error("nothing drawn")
	}
	v := doc.NewView("v", bdf.ViewFixed, "")
	v.AddPage(100, 50, bdf.Layer{Role: bdf.RoleBody, Obj: cv.Hash()})
	idx, err := doc.BuildTextIndex(v)
	if err != nil {
		t.Fatal(err)
	}
	_ = idx
}
