package xlsx

import (
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// Equation shapes are laid out from their Office Math (the a14 choice), not
// drawn from the text of their fallback.
func TestMath(t *testing.T) {
	_, r := convert(t, "math.xlsx", testOptions())
	text := bdf.PlainText(indexRuns(t, r, r.Manifest.Views[0]))
	for _, want := range []string{"x=(−b±√(b^2−4ac))/(2a)", "判別式 D=b^2−4ac が正なら実数解は二つ。"} {
		if !strings.Contains(text, want) {
			t.Errorf("%q not in %q", want, text)
		}
	}
	if strings.Contains(text, "FALLBACK") {
		t.Errorf("fallback text drawn: %q", text)
	}
}
