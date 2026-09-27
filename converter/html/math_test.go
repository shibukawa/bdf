package html

import (
	"strings"
	"testing"
)

// Formulas are laid out whatever wrote them, once each: the renderings
// that KaTeX, MathJax and Wikipedia put beside the MathML are dropped.
func TestMath(t *testing.T) {
	src := `<!DOCTYPE html><body>
<p>MathML <math><msup><mi>x</mi><mn>2</mn></msup></math> inline.</p>
<p>KaTeX <span class="katex"><span class="katex-mathml"><math><semantics><mrow><mi>a</mi><mo>+</mo><mi>b</mi></mrow>` +
		`<annotation encoding="application/x-tex">a+b</annotation></semantics></math></span>` +
		`<span class="katex-html" aria-hidden="true"><span class="base">GARBLED</span></span></span> end.</p>
<p>MathJax <mjx-container class="MathJax" jax="CHTML" display="true"><mjx-math><mjx-mi>GLYPHS</mjx-mi></mjx-math>` +
		`<mjx-assistive-mml><math display="block"><mfrac><mn>1</mn><mn>2</mn></mfrac></math></mjx-assistive-mml></mjx-container></p>
<p>Wikipedia <span class="mwe-math-element"><span class="mwe-math-mathml-inline" style="display: none;">` +
		`<math><mi>y</mi></math></span><img src="y.svg" class="mwe-math-fallback-image-inline" alt="{\displaystyle y}"></span> end.</p>
<p>MathJax 2 <span class="MathJax_Preview">PREVIEW</span><script type="math/tex">\sqrt{z}</script> end.</p>
<p>GitHub <math-renderer class="js-inline-math">$\alpha_1$</math-renderer> end.</p>
</body>`
	opts := testOptions()
	opts.Extract = ExtractNone
	res, r := convertHTML(t, src, opts)
	c := viewContent(t, r, 0)
	for _, want := range []string{"x^2", "a+b", "1/2", "y", "√z", "α_1"} {
		if !strings.Contains(c.text, want) {
			t.Errorf("%q not in %q", want, c.text)
		}
	}
	for _, bad := range []string{"GARBLED", "GLYPHS", "PREVIEW", "displaystyle"} {
		if strings.Contains(c.text, bad) {
			t.Errorf("%q in %q", bad, c.text)
		}
	}
	if len(c.images) != 0 {
		t.Errorf("fallback images drawn: %v", c.images)
	}
	// the text refers to fonts by name, the formula font is embedded
	if res.EmbeddedFonts != 1 {
		t.Errorf("%d fonts embedded, want the formula font", res.EmbeddedFonts)
	}
	for _, w := range res.Warnings {
		if strings.Contains(w, "MATH") {
			t.Errorf("warning %q", w)
		}
	}
}

// Articles keep their formulas.
func TestMathInArticle(t *testing.T) {
	p := "<p>" + strings.Repeat("A paragraph of the article, long enough to be picked out as its content. ", 8) + "</p>"
	src := `<!DOCTYPE html><html><head><title>Math</title></head><body><nav>menu</nav><article><h1>Math</h1>` + p +
		`<p>The identity <span class="katex"><span class="katex-mathml"><math><mrow><msup><mi>e</mi><mrow><mi>i</mi><mi>π</mi></mrow></msup>` +
		`<mo>+</mo><mn>1</mn><mo>=</mo><mn>0</mn></mrow></math></span><span class="katex-html">GARBLED</span></span> holds.</p>` + p +
		`</article></body></html>`
	res, r := convertHTML(t, src, testOptions())
	if !res.Extracted {
		t.Fatal("no article extracted")
	}
	c := viewContent(t, r, 0)
	if !strings.Contains(c.text, "e^(iπ)+1=0") || strings.Contains(c.text, "GARBLED") {
		t.Errorf("text %q", c.text)
	}
}
