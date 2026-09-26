package xlsx

import (
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// fillTexts returns the x of the text an object (and the objects it uses)
// draws, by text.
func fillTexts(t *testing.T, r *bdf.Reader, h bdf.Hash) map[string]float32 {
	t.Helper()
	out := map[string]float32{}
	var walk func(h bdf.Hash)
	walk = func(h bdf.Hash) {
		o, err := r.Object(h)
		if err != nil {
			t.Fatal(err)
		}
		o.Walk(func(in bdf.Instr) {
			switch in.Op {
			case bdf.OpFillText:
				out[in.Args[0].(string)] = in.Args[1].(float32)
			case bdf.OpUse, bdf.OpUseAt:
				walk(o.Objects[int(in.Args[0].(uint64))])
			}
		})
	}
	walk(h)
	return out
}

func sizes(runs []bdf.Run) []float32 {
	var out []float32
	for _, run := range runs {
		for i := 0; i < int(run[0]); i++ {
			out = append(out, run[1])
		}
	}
	return out
}

func testGrid() *Grid {
	return &Grid{Name: "Data", HeaderRows: 1, Rows: [][]GridCell{
		{{Text: "name"}, {Text: "amount"}, {Text: "note"}},
		{{Text: "Ann"}, {Text: "1,234.50", Number: true}, {Text: "two\r\nlines"}},
		{{Text: "Bob"}, {Text: "7", Number: true}, {Text: strings.Repeat("a long note ", 20)}},
		{},
		{{Text: "Cy"}, {Text: "007"}},
	}}
}

func TestConvertGrid(t *testing.T) {
	res, err := ConvertGrid(testGrid(), testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) != 0 || res.Sheets != 1 || res.EmbeddedFonts != 2 {
		t.Errorf("result: %+v", res)
	}
	r := reopen(t, res)
	v := r.Manifest.Views[0]
	if v.Title != "Data" || v.Kind != bdf.ViewSheet || !v.Gridlines || v.Freeze == nil || v.Freeze.Rows != 1 || v.Freeze.Cols != 0 {
		t.Errorf("view = %q %s grid %v freeze %+v", v.Title, v.Kind, v.Gridlines, v.Freeze)
	}
	// A and B keep the default width; C fits its text, at most 50 digits
	// (of 9 pixels in the test font) and the padding
	w := sizes(v.Cols)
	if len(w) < 12 || w[0] != w[3] || w[1] != w[0] || w[2] != (50*9+5)*pxPt {
		t.Errorf("column widths = %v", w)
	}
	h := sizes(v.Rows)
	if h[1] <= h[2] || h[2] != h[0] {
		t.Errorf("row heights = %v (row 2 wraps two lines)", h[:5])
	}
	ms, _ := marks(t, r, tile(t, v, "0,0"))
	for _, want := range []string{"C A1 col", "C C1 col", "C B2", "C C3", "C B5"} {
		if !contains(ms, want) {
			t.Errorf("no mark %q in %q", want, ms)
		}
	}
	// numbers on the right of their cells, text on the left
	xs := fillTexts(t, r, tile(t, v, "0,0"))
	if x := xs["7"]; x < w[0]+w[1]/2 || x > w[0]+w[1] {
		t.Errorf("7 at x %g (column B is %g–%g)", x, w[0], w[0]+w[1])
	}
	if x := xs["007"]; x > w[0]+w[1]/2 {
		t.Errorf("007 at x %g", x)
	}
	text := bdf.PlainText(indexRuns(t, r, v))
	for _, want := range []string{"name\namount\nnote", "Ann\n1,234.50\ntwo lines", "Cy\n007"} {
		if !strings.Contains(text, want) {
			t.Errorf("text index lacks %q:\n%s", want, text)
		}
	}
}

func TestGridTableStyle(t *testing.T) {
	fills := func(g *Grid) int {
		res, err := ConvertGrid(g, testOptions())
		if err != nil {
			t.Fatal(err)
		}
		r := reopen(t, res)
		_, counts := marks(t, r, tile(t, r.Manifest.Views[0], "0,0"))
		return counts[bdf.OpFillRect] + counts[bdf.OpFillPath]
	}
	plain := fills(testGrid())
	g := testGrid()
	g.TableStyle = "medium2"
	if styled := fills(g); styled <= plain+3 {
		t.Errorf("the table style draws %d fills, the plain grid %d", styled, plain)
	}
	g.TableStyle = "TableStyleNope"
	if _, err := ConvertGrid(g, testOptions()); err == nil {
		t.Error("an unknown table style was accepted")
	}
}

func TestGridLimits(t *testing.T) {
	row := make([]GridCell, maxCols+10)
	row[0].Text, row[maxCols-1].Text, row[maxCols+5].Text = "first", "last", "beyond"
	res, err := ConvertGrid(&Grid{Rows: [][]GridCell{row}, Lang: "ja"}, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "16384 columns") {
		t.Errorf("warnings = %q", res.Warnings)
	}
	if v := res.Doc.Views[0]; v.Title != "Sheet1" {
		t.Errorf("title = %q", v.Title)
	}
}
