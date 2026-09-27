package drawio

import (
	"strings"
	"testing"
)

// mathLines lays out a label on a page with math="1" (or not) and returns
// its lines, formulas written as {linear notation}.
func mathLines(t *testing.T, math bool, cell string) []string {
	t.Helper()
	c, l := func() (*converter, *labelBox) {
		root, err := parseXML([]byte(`<mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/>` + cell + `</root></mxGraphModel>`))
		if err != nil {
			t.Fatal(err)
		}
		c := newConverter(testOptions())
		c.m = parseModel(root)
		c.math = math
		c.page = &pageInfo{name: "P", number: 1, count: 1}
		v := newView(c.m, nil, nil)
		st := v.state(c.m.cells["a"])
		return c, c.layoutLabel(st, c.newShape(st))
	}()
	if !c.mathEngine().HasMathFont() {
		t.Fatal("no formula font among the test fonts")
	}
	var out []string
	for _, ln := range l.lo.lines {
		var b strings.Builder
		for _, it := range ln.items {
			if it.eq != nil {
				b.WriteString("{" + it.eq.text + "}")
				continue
			}
			b.WriteRune(it.r)
		}
		out = append(out, b.String())
	}
	return out
}

func TestMathLabels(t *testing.T) {
	geo := `<mxGeometry x="0" y="0" width="200" height="120" as="geometry"/>`
	// a display formula takes a line of its own; inline ones stay in theirs
	cell := `<mxCell id="a" value="Roots of \(ax^2+bx+c=0\):$$x=\frac{-b\pm\sqrt{b^2-4ac}}{2a}$$when \(a \ne 0\)" style="whiteSpace=wrap;html=1;" vertex="1" parent="1">` + geo + `</mxCell>`
	want := []string{"Roots of {ax^2+bx+c=0}:", "{x=(−b±√(b^2−4ac))/(2a)}", "when {a≠0}"}
	if got := mathLines(t, true, cell); strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("lines %q, want %q", got, want)
	}
	// without math="1" the source is text
	if got := mathLines(t, false, cell); len(got) == 0 || !strings.Contains(strings.Join(got, " "), `\frac`) {
		t.Errorf("lines without math %q", got)
	}
	// formulas wrap with the words around them
	cell = `<mxCell id="a" value="Inline \(e^{i\pi}+1=0\) and \(\sum_{k=1}^n k\) wrap with the words around them" style="whiteSpace=wrap;html=1;" vertex="1" parent="1"><mxGeometry x="0" y="0" width="120" height="120" as="geometry"/></mxCell>`
	got := mathLines(t, true, cell)
	if len(got) < 3 || !strings.Contains(strings.Join(got, " "), "{e^(iπ)+1=0}") {
		t.Errorf("wrapped lines %q", got)
	}
	// in plain labels, formulas stay in the lines they are in; AsciiMath
	// between backquotes is text
	cell = "<mxCell id=\"a\" value=\"Plain \\(x^2\\) and `x^2`&#xa;$$y$$ too\" style=\"rounded=0;\" vertex=\"1\" parent=\"1\">" + geo + `</mxCell>`
	want = []string{"Plain {x^2} and `x^2`", "{y} too"}
	if got := mathLines(t, true, cell); strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("plain lines %q, want %q", got, want)
	}
}

func TestNextFormula(t *testing.T) {
	for _, tc := range []struct {
		s, tex  string
		display bool
	}{
		{`a $$x$$ b`, "x", true},
		{`a \(y\) b $$x$$`, "y", false},
		{`\[z\]`, "z", true},
		{`no closing \(x`, "", false},
	} {
		start, _, display, tex := nextFormula(tc.s)
		if tc.tex == "" {
			if start >= 0 {
				t.Errorf("%q: formula %q", tc.s, tex)
			}
			continue
		}
		if tex != tc.tex || display != tc.display {
			t.Errorf("%q: %q display %v", tc.s, tex, display)
		}
	}
}

func TestConvertMath(t *testing.T) {
	res := convertTest(t, "math.drawio")
	for _, w := range res.Warnings {
		if strings.Contains(w, "math") || strings.Contains(w, "MATH") {
			t.Errorf("warning %q", w)
		}
	}
	text := viewText(t, res.Doc, res.Doc.Views[0])
	for _, want := range []string{
		"x=(−b±√(b^2−4ac))/(2a)",
		"e^(iπ)+1=0",
		"二次方程式の判別式 D=b^2−4ac が正のとき",
		"A=(1, 2; 3, 4)",
		"lim_(h→0) (f(x+h)−f(x))/h",
		"AsciiMath stays text: `sqrt(x^2+1)`",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("%q not in %q", want, text)
		}
	}
	if strings.Contains(text, `\frac`) {
		t.Errorf("LaTeX source drawn: %q", text)
	}
}
