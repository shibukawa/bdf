package equation

import (
	"runtime"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/internal/fontdb"
	"github.com/shibukawa/bdf/internal/mathlayout"
)

// testEngine returns an engine with the test fonts (STIX Two Math and M
// PLUS 1p).
func testEngine(t *testing.T) (*Engine, *fontset.Set) {
	t.Helper()
	set := fontset.New(fontdb.New(nil, []string{"../../docx/testdata/fonts"}, false), nil)
	text := func(family string, r rune, bold, italic bool) *fontset.Choice {
		return set.FaceFor(family, "", bold, italic, r)
	}
	return New(Fonts{Set: set, Math: "Cambria Math", Text: text}), set
}

// The glyphs drawn by index (the larger variants of delimiters) are drawn
// with private use characters of the formula font, which is embedded even
// when the other fonts are referred to by name.
func TestGlyphsByIndex(t *testing.T) {
	e, set := testEngine(t)
	if !e.HasMathFont() {
		t.Fatal("the test font STIX Two Math is not found")
	}
	tall := e.Layout(ParseTeX(`\left( \frac{\frac{a}{b}}{\frac{c}{d}} \right)`), Style{Size: 10, Display: true})
	pua := false
	for _, it := range tall.Items() {
		for _, r := range it.Text {
			if r >= 0xF0000 {
				pua = true
				if it.Alt != "(" && it.Alt != ")" {
					t.Errorf("private use glyph for %q", it.Alt)
				}
			}
		}
	}
	if !pua {
		t.Error("the parentheses around a tall fraction did not grow")
	}
	if n := set.Embed(bdf.NewDocument(), fontset.EmbedOptions{PinnedOnly: true}); n != 1 {
		t.Errorf("embedded %d fonts with PinnedOnly, want the formula font", n)
	}
}

func TestPlace(t *testing.T) {
	e, set := testEngine(t)
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

// The glyphs of a run are joined once when the run is drawn: joining them
// one by one copied the run for each glyph (30000 digits: 450 MB).
func TestDrawLongRun(t *testing.T) {
	e, set := testEngine(t)
	digits := strings.Repeat("1234567890", 3000)
	b := e.Layout(&Row{Kids: []Node{&Atom{Kind: mathlayout.Number, Text: digits}, &Atom{Kind: mathlayout.Ident, Text: "x"}}}, Style{Size: 10})
	cv := canvas.NewBuilder(bdf.NewDocument(), set).New()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	Draw(cv, b, 0, 0)
	runtime.ReadMemStats(&after)
	if n := after.TotalAlloc - before.TotalAlloc; n > 16<<20 {
		t.Errorf("drawing %d digits allocated %d MiB", len(digits), n>>20)
	}
	o, err := bdf.DecodeObject(cv.Obj.Encode())
	if err != nil {
		t.Fatal(err)
	}
	var texts []string
	o.Walk(func(in bdf.Instr) {
		if in.Op == bdf.OpFillText {
			texts = append(texts, in.Args[0].(string))
		}
	})
	// the digits and the letter after them, in the formula font, are a run
	if len(texts) != 1 || texts[0] != digits+"𝑥" {
		t.Errorf("%d runs drawn: %.40q", len(texts), texts)
	}
}
