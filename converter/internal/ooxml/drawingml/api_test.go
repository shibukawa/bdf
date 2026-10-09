package drawingml

import (
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/internal/fontdb"
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
	fonts := fontset.New(fontdb.New(nil, []string{"../../../pptx/testdata/fonts"}, false), func(msg string) { t.Error(msg) })
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

func TestTextFlow(t *testing.T) {
	// three paragraphs, the second justified with a first-line indent
	body, err := ooxml.Parse([]byte(`<txBody><bodyPr lIns="0" tIns="0" rIns="0" bIns="0"/>` +
		`<p><r><rPr sz="1200"/><t>one two three four five six seven eight nine ten</t></r></p>` +
		`<p><pPr algn="just" indent="180000"/><r><rPr sz="1200"/><t>alpha beta gamma delta epsilon zeta eta theta iota kappa lambda</t></r></p>` +
		`<p><r><rPr sz="1200"/><t>end</t></r></p></txBody>`))
	if err != nil {
		t.Fatal(err)
	}
	doc := bdf.NewDocument()
	fonts := fontset.New(fontdb.New(nil, []string{"../../../pptx/testdata/fonts"}, false), func(msg string) { t.Error(msg) })
	d := New(Config{Doc: doc, Fonts: fonts}).NewDrawing("", nil, plainHost{})
	bodyPr, _ := ooxml.Parse([]byte(`<bodyPr lIns="0" tIns="0" rIns="0" bIns="0"/>`))
	whole := d.LayoutText(body, "", 0, 0, 120, 1000, canvas.Identity)
	_, _, _, wholeBottom := whole.Bounds()
	total := len(whole.b.lo.lines)

	f := d.NewTextFlow(body, "")
	var got []string
	boxes := 0
	for !f.Done() {
		if boxes > 20 {
			t.Fatal("the flow never ends")
		}
		boxes++
		tb := f.Fill(bodyPr, 0, 0, 120, 40, canvas.Identity)
		if tb == nil {
			t.Fatal("an empty box while text is left")
		}
		_, _, _, y1 := tb.Bounds()
		if y1 > 40+1e-6 {
			t.Errorf("box %d: text reaches %g, past the box", boxes, y1)
		}
		for _, ln := range tb.b.lo.lines {
			var s []rune
			for _, it := range ln.items {
				s = append(s, it.r)
			}
			got = append(got, string(s))
		}
	}
	if boxes < 2 {
		t.Errorf("%d box(es) for %g pt of text", boxes, wholeBottom)
	}
	if len(got) != total {
		t.Errorf("%d lines through the flow, %d laid out at once: %q", len(got), total, got)
	}
	if got[len(got)-1] != "end" {
		t.Errorf("last line %q", got[len(got)-1])
	}
	// a continued paragraph keeps the whole text
	var text []rune
	for _, p := range body.Children("p") {
		text = append(text, []rune(p.Child("r").Child("t").Content())...)
	}
	var flowed []rune
	for _, l := range got {
		flowed = append(flowed, []rune(l)...)
	}
	if cleaned := func(r []rune) string {
		var b []rune
		for _, c := range r {
			if c != ' ' {
				b = append(b, c)
			}
		}
		return string(b)
	}; cleaned(text) != cleaned(flowed) {
		t.Errorf("flowed text %q", string(flowed))
	}
	// boxes of every height, with space after the paragraphs: the text
	// comes through whole every time
	spaced, _ := ooxml.Parse([]byte(`<txBody><bodyPr lIns="0" tIns="0" rIns="0" bIns="0"/>` +
		`<p><pPr><spcAft><spcPts val="3000"/></spcAft></pPr><r><rPr sz="1200"/><t>one two three four five six seven</t></r></p>` +
		`<p><pPr><spcAft><spcPts val="3000"/></spcAft></pPr><r><rPr sz="1200"/><t>alpha beta gamma delta epsilon</t></r></p>` +
		`<p><r><rPr sz="1200"/><t>end</t></r></p></txBody>`))
	spacedLines := len(d.LayoutText(spaced, "", 0, 0, 120, 1000, canvas.Identity).b.lo.lines)
	for h := 5.0; h <= 80; h += 5 {
		f := d.NewTextFlow(spaced, "")
		lines, boxes := 0, 0
		for !f.Done() && boxes < 50 {
			boxes++
			if tb := f.Fill(bodyPr, 0, 0, 120, h, canvas.Identity); tb != nil {
				lines += len(tb.b.lo.lines)
			}
		}
		if !f.Done() || lines != spacedLines {
			t.Errorf("height %g: %d lines in %d boxes, done %v", h, lines, boxes, f.Done())
		}
	}
	// a multi-column box takes more lines
	g := d.NewTextFlow(body, "")
	cols, _ := ooxml.Parse([]byte(`<bodyPr lIns="0" tIns="0" rIns="0" bIns="0" numCol="2" spcCol="0"/>`))
	tb := g.Fill(cols, 0, 0, 240, 40, canvas.Identity)
	if tb == nil || len(tb.b.lo.lines) <= len(whole.b.lo.lines)/4 {
		t.Errorf("two columns of 40 pt hold %d lines", len(tb.b.lo.lines))
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
