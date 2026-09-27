package pptx

import (
	"strings"
	"testing"
)

// Equations are laid out from their Office Math (the a14 choice), not drawn
// from the picture of their fallback; other a14 choices keep their
// fallbacks.
func TestMath(t *testing.T) {
	res, r := convert(t, "math.pptx", testOptions())
	for _, w := range res.Warnings {
		if strings.Contains(w, "MATH") {
			t.Errorf("warning %q", w)
		}
	}
	text := plainText(t, r)
	for _, want := range []string{
		"x=(−b±√(b^2−4ac))/(2a)",
		"∑_(k=1)^n k^2=(n(n+1)(2n+1))/6",
		"Euler's identity e^(iπ)+1=0 holds for every text line.",
		"二次方程式の判別式は D=b^2−4ac で、行の中に置かれます。",
		"A=(1, 2; 3, 4)",
		"Other a14 content keeps its fallback",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("%q not in %q", want, text)
		}
	}
	for _, bad := range []string{"FALLBACK PICTURE", "CHOICE WITHOUT MATH"} {
		if strings.Contains(text, bad) {
			t.Errorf("%q in %q", bad, text)
		}
	}
	if res.EmbeddedFonts != 2 {
		t.Errorf("%d fonts embedded, want M PLUS 1p and STIX Two Math", res.EmbeddedFonts)
	}
	// with system fonts, the formula font is embedded still
	opts := testOptions()
	opts.SystemFonts = true
	if res, _ := convert(t, "math.pptx", opts); res.EmbeddedFonts != 1 {
		t.Errorf("%d fonts embedded with system fonts, want the formula font", res.EmbeddedFonts)
	}
}
