package drawingml

import (
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// plainHost is a host with nothing to inherit.
type plainHost struct{}

func (plainHost) Placeholder(*ooxml.Node, string) ([]Inherited, bool) { return nil, true }
func (plainHost) TextStyle(string) *ooxml.Node                        { return nil }
func (plainHost) Field(string) (string, bool)                         { return "", false }
func (plainHost) Link(string, ooxml.Rel) string                       { return "" }

func TestLayoutText(t *testing.T) {
	body, err := ooxml.Parse([]byte(`<txBody><bodyPr lIns="91440" tIns="0" rIns="91440" bIns="0" anchor="ctr" wrap="square"/>` +
		`<p><pPr algn="ctr"/><r><rPr sz="1200"/><t>centered text</t></r></p></txBody>`))
	if err != nil {
		t.Fatal(err)
	}
	doc := bdf.NewDocument()
	fonts := fontset.New(fontdb.New([]string{"../../../pptx/testdata/fonts"}, false), func(msg string) { t.Error(msg) })
	// no package: drawings without one lay out text too
	d := New(Config{Doc: doc, Fonts: fonts}).NewDrawing("", nil, plainHost{})
	tb := d.LayoutText(body, "", 0, 0, 200, 100, canvas.Translate(10, 20))
	if tb == nil {
		t.Fatal("no text laid out")
	}
	x0, y0, x1, y1 := tb.Bounds()
	// centered both ways, with the 7.2 pt insets on the sides
	if !(x0 > 20 && x1 < 180 && x0 < 100 && x1 > 100) || !(y0 > 30 && y1 < 70) {
		t.Errorf("bounds = %g %g %g %g", x0, y0, x1, y1)
	}
	if mid := (x0 + x1) / 2; mid < 99 || mid > 101 {
		t.Errorf("bounds not centered: %g", mid)
	}
	cv := canvas.NewBuilder(doc, fonts).New()
	tb.Draw(cv)
	if !cv.Drawn {
		t.Error("nothing drawn")
	}
	empty, _ := ooxml.Parse([]byte(`<txBody><bodyPr/><p/></txBody>`))
	if d.LayoutText(empty, "", 0, 0, 100, 100, canvas.Identity) != nil {
		t.Error("an empty body laid out")
	}
}

func TestResolveColor(t *testing.T) {
	n, err := ooxml.Parse([]byte(`<solidFill><schemeClr val="phClr"><lumMod val="50000"/></schemeClr></solidFill>`))
	if err != nil {
		t.Fatal(err)
	}
	if c, ok := ResolveColor(n, nil, 0xFF0000FF); !ok || c != 0x800000FF {
		t.Errorf("phClr lumMod 50%% = %08x, %v", uint32(c), ok)
	}
	n, _ = ooxml.Parse([]byte(`<fill><schemeClr val="accent1"><alpha val="50000"/></schemeClr></fill>`))
	if c, ok := ResolveColor(n, map[string]bdf.Color{"accent1": 0x2E75B6FF}, 0); !ok || c != 0x2E75B680 {
		t.Errorf("accent1 alpha 50%% = %08x, %v", uint32(c), ok)
	}
	n, _ = ooxml.Parse([]byte(`<fill><noFill/></fill>`))
	if _, ok := ResolveColor(n, nil, 0); ok {
		t.Error("a color out of no color")
	}
}
