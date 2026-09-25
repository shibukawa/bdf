package fontdb

import (
	"testing"

	"github.com/shibukawa/bdf/converter/internal/sfnt"
)

func TestNormalize(t *testing.T) {
	for in, want := range map[string]string{"ＭＳ Ｐゴシック": "mspゴシック", "Times New Roman": "timesnewroman", "Noto_Sans-CJK": "notosanscjk"} {
		if got := Normalize(in); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestClassify(t *testing.T) {
	for in, want := range map[string]string{"Times New Roman": Serif, "Calibri": Sans, "Courier New": Mono, "游明朝": Serif, "ＭＳ ゴシック": Mono, "Meiryo": Sans} {
		if got := Classify(in); got != want {
			t.Errorf("Classify(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestHeavyName(t *testing.T) {
	for in, want := range map[string]bool{"Arial Black": true, "Segoe UI Semibold": true, "Calibri": false, "Bodoni": false} {
		if got := heavyName(in); got != want {
			t.Errorf("heavyName(%q) = %v", in, got)
		}
	}
}

func TestSystemResolve(t *testing.T) {
	db := New(nil, true)
	if len(db.Faces) == 0 {
		t.Skip("no system fonts")
	}
	r := db.Resolve("Calibri", true, false, false)
	if r.Face == nil {
		t.Fatal("no face for Calibri")
	}
	t.Logf("Calibri bold -> %s %s (synth bold %v)", r.Face.Family, r.Face.Style, r.SynthBold)
	r = db.Resolve("游ゴシック", false, false, true)
	if r.Face != nil {
		t.Logf("游ゴシック -> %s", r.Face.Family)
	}
	l, err := r.Face.Load()
	if err != nil {
		t.Fatal(err)
	}
	prog, ok := l.Program([]rune("Hello あいう"), false)
	if !ok {
		t.Fatal("not embeddable")
	}
	sf, err := sfnt.Parse(prog)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("subset %d bytes (from %d), %d glyphs", len(prog), l.Size(), sf.NumGlyphs)
	if l.CanSubset() && sf.NumGlyphs > 20 {
		t.Fatalf("subset kept %d glyphs", sf.NumGlyphs)
	}
	if g, ok := sf.Cmap['H']; !ok || g == 0 {
		t.Fatal("H missing from subset cmap")
	}
}
