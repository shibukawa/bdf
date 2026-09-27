package fontdb

import (
	"os"
	"testing"
	"testing/fstest"

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
	db := New(nil, nil, true)
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
	prog, ok := l.Program([]rune("Hello あいう"), nil, false, false)
	if !ok {
		t.Fatal("not embeddable")
	}
	sf, err := sfnt.Parse(prog)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("subset %d bytes (from %d), %d glyphs", len(prog), l.Size(), sf.NumGlyphs)
	if l.CanSubset(false) && sf.NumGlyphs > 20 {
		t.Fatalf("subset kept %d glyphs", sf.NumGlyphs)
	}
	if g, ok := sf.Cmap['H']; !ok || g == 0 {
		t.Fatal("H missing from subset cmap")
	}
}

func TestFS(t *testing.T) {
	const dir = "../../pptx/testdata/fonts"
	fsys := fstest.MapFS{"README.md": {Data: []byte("not a font")}}
	for _, name := range []string{"MPLUS1p-Regular-subset.ttf", "MPLUS1p-Bold-subset.ttf"} {
		data, err := os.ReadFile(dir + "/" + name)
		if err != nil {
			t.Fatal(err)
		}
		fsys["fonts/"+name] = &fstest.MapFile{Data: data}
	}
	db := New(fsys, nil, false)
	if len(db.Faces) != 2 {
		t.Fatalf("%d faces, want 2", len(db.Faces))
	}
	r := db.Resolve("Meiryo", true, false, true)
	if r.Face == nil || r.Face.Path != "fonts/MPLUS1p-Bold-subset.ttf" || r.SynthBold {
		t.Fatalf("Meiryo bold resolved to %+v", r)
	}
	l, err := r.Face.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !l.Has('あ') {
		t.Error("loaded face has no あ")
	}
	// The faces of the file system come before those of the directories
	// (M PLUS 1p and the formula font STIX Two Math).
	db = New(fsys, []string{dir}, false)
	if len(db.Faces) != 5 || db.Faces[0].Path != "fonts/MPLUS1p-Bold-subset.ttf" {
		t.Errorf("faces: %d, first %s", len(db.Faces), db.Faces[0].Path)
	}
}

// Formula fonts: found by their MATH table, and subset with the glyphs
// they draw by index (the parts of delimiters) mapped from private use
// characters, CFF outlines included.
func TestMathFont(t *testing.T) {
	for _, name := range []string{"../sfnt/testdata/STIXTwoMath-cff-subset.otf", "../../docx/testdata/fonts/STIXTwoMath-subset.ttf"} {
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		db := New(fstest.MapFS{"math.otf": {Data: data}}, nil, false)
		f := db.ResolveMath("Cambria Math")
		if f == nil || !f.Math || f.Family != "STIX Two Math" {
			t.Fatalf("%s: math font %+v", name, f)
		}
		l, err := f.Load()
		if err != nil {
			t.Fatal(err)
		}
		if l.Math() == nil {
			t.Fatalf("%s: no MATH table", name)
		}
		if b, ok := l.Bounds(l.Font.Cmap['(']); !ok || b.YMax < 0.7 || b.YMin > -0.15 {
			t.Errorf("%s: bounds of ( %+v", name, b)
		}
		part := l.Math().Vertical(l.Font.Cmap['(']).Parts[1].Glyph
		prog, ok := l.Program([]rune("x("), map[rune]uint16{0xF0000 + rune(part): part}, false, false)
		if !ok {
			t.Fatal("not embeddable")
		}
		sf, err := sfnt.Parse(prog)
		if err != nil {
			t.Fatal(err)
		}
		if sf.NumGlyphs != 4 || sf.IsCFF != f.CFF {
			t.Errorf("%s: subset of %d glyphs (CFF %v)", name, sf.NumGlyphs, sf.IsCFF)
		}
		if g, ok := sf.Cmap[0xF0000+uint32(part)]; !ok || g == 0 || g == sf.Cmap['('] {
			t.Errorf("%s: part of ( mapped to glyph %d", name, g)
		}
	}
	data, err := os.ReadFile("../../pptx/testdata/fonts/MPLUS1p-Regular-subset.ttf")
	if err != nil {
		t.Fatal(err)
	}
	if New(fstest.MapFS{"mplus.ttf": {Data: data}}, nil, false).ResolveMath("Cambria Math") != nil {
		t.Error("a font without a MATH table resolved as a formula font")
	}
}
