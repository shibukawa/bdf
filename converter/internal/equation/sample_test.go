package equation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
)

// TestSamples writes a document of formulas to look at when
// EQUATION_SAMPLES names a directory (render it with test/render.mjs).
func TestSamples(t *testing.T) {
	dir := os.Getenv("EQUATION_SAMPLES")
	if dir == "" {
		t.Skip("EQUATION_SAMPLES not set")
	}
	for _, withMath := range []bool{true, false} {
		fontDir := "../../docx/testdata/fonts"
		if !withMath {
			fontDir = t.TempDir()
			for _, f := range []string{"MPLUS1p-Regular-subset.ttf", "MPLUS1p-Bold-subset.ttf"} {
				data, _ := os.ReadFile(filepath.Join("../../docx/testdata/fonts", f))
				os.WriteFile(filepath.Join(fontDir, f), data, 0o644)
			}
		}
		db := fontdb.New(nil, []string{fontDir}, false)
		set := fontset.New(db, func(m string) { t.Log(m) })
		doc := bdf.NewDocument()
		cvs := canvas.NewBuilder(doc, set)
		cv := cvs.New()
		text := func(family string, r rune, bold, italic bool) *fontset.Choice {
			return set.FaceFor(family, "", bold, italic, r)
		}
		e := New(Fonts{Set: set, Math: "Cambria Math", Text: text, Warn: func(m string) { t.Log(m) }})
		label := set.Choose("M PLUS 1p", false, false, false)
		y := 40.0
		for _, s := range sampleTeX {
			for _, r := range s {
				set.Advance(label, r)
			}
			cv.Obj.Font(cv.Font(label.Use), 9)
			cv.Obj.FillColor(bdf.RGBA(120, 120, 120, 255))
			cv.Obj.FillText(s, 20, float32(y), 0)
			n := ParseTeX(s)
			b := e.Layout(n, Style{Size: 16, Display: true})
			y += b.H + 14
			Place(cv, b, 40, y, Linear(n))
			y += b.D + 24
		}
		set.Embed(doc, fontset.EmbedOptions{})
		cvs.Encode()
		v := doc.NewView("p", bdf.ViewFixed, "samples")
		v.AddPage(720, float32(y+20), bdf.Layer{Role: bdf.RoleBody, Obj: cv.Hash()})
		name := "samples.bdf"
		if !withMath {
			name = "samples-textfont.bdf"
		}
		f, err := os.Create(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		if err := doc.WriteSingle(f); err != nil {
			t.Fatal(err)
		}
		f.Close()
	}
}

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
