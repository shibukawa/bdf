package docx

import (
	"math"
	"strings"
	"testing"
)

// Office Math is laid out by the formula engine: each formula is one run
// of text (its linear notation), inline among the text around it or
// centered on a line of its own.
func TestMath(t *testing.T) {
	res, r := convert(t, "math.docx", testOptions())
	for _, w := range res.Warnings {
		if strings.Contains(w, "MATH") || strings.Contains(w, "equations") {
			t.Errorf("warning %q", w)
		}
	}
	rs := runs(t, r, 0, 0, "body")
	all := text(rs)
	for _, want := range []string{
		"x=(−b±√(b^2−4ac))/(2a)",
		"∑_(k=1)^n k^2=(n(n+1)(2n+1))/6",
		"lim_(h→0) (f(x+h)−f(x))/h",
		"A=(a_11, a_12; a_21, a_22)",
		"f(x)={x^2,x≥0\n−x,x<0",
		"v=(距離)/(時間)",
		"で、表面積は S=4πr^2 です",
	} {
		if !strings.Contains(all, want) {
			t.Errorf("%q not in %q", want, all)
		}
	}
	pg := r.Manifest.Views[0].Pages[0]
	// a display formula is centered in the text column (72 pt margins)
	q, ok := find(rs, "x=(−b±")
	if !ok || !q.AltText {
		t.Fatalf("quadratic formula run %+v", q)
	}
	if mid := q.X + q.Advance/2; math.Abs(float64(mid-pg.W/2)) > 2 {
		t.Errorf("display formula centered at %v, page %v wide", mid, pg.W)
	}
	// an inline formula sits on the line of the text around it
	f, _ := find(rs, "ax^2+bx+c=0")
	txt, _ := find(rs, "二次方程式")
	if math.Abs(float64(f.Y-txt.Y)) > 0.5 || f.X <= txt.X {
		t.Errorf("inline formula at %v,%v; text at %v,%v", f.X, f.Y, txt.X, txt.Y)
	}
	if res.EmbeddedFonts != 3 {
		t.Errorf("%d fonts embedded, want M PLUS 1p (regular and bold) and STIX Two Math", res.EmbeddedFonts)
	}
}
