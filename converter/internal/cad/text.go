package cad

import (
	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
)

// FontSpec names the font of a text.
type FontSpec struct {
	// Family is the font of Latin text ("" for the default sans-serif).
	Family string
	// EastAsian is the font of CJK characters ("" for Family).
	EastAsian    string
	Bold, Italic bool
}

func (f FontSpec) ea() string {
	if f.EastAsian != "" {
		return f.EastAsian
	}
	return f.Family
}

// Break is how a text joins the text drawn before it, for text extraction
// (the MARK kinds of the spec, §7.8).
type Break byte

const (
	// BreakBox starts a text of its own (the default).
	BreakBox Break = iota
	// BreakParagraph starts a paragraph of the same text.
	BreakParagraph
	// BreakLine starts a line of the same paragraph (a space in between).
	BreakLine
	// BreakWrap starts a line wrapped without a space.
	BreakWrap
	// BreakNone continues the line.
	BreakNone
)

// Text is a line of text in one font. Its em space has the origin at the
// start of the baseline, x along the text and y up, in ems; M maps it into
// the drawing (size, width factor, obliquing, rotation, mirroring).
type Text struct {
	S     string
	Font  FontSpec
	Color bdf.Color
	M     canvas.Matrix
	// Advance is the advance of S in ems: its natural advance, or the sum
	// of Cells and Spacing.
	Advance float64
	// Cells, when not nil, are the advances of the characters of S in ems
	// (fixed-pitch text): each run of characters is stretched to the width
	// of its cells.
	Cells []float64
	// Spacing is added after each character (ems).
	Spacing float64
	// Asc, Desc and Cap are the ascent, descent and cap height of the font
	// in ems.
	Asc, Desc, Cap float64

	Underline, Overline, Strike bool
	Break                       Break
	// Vertical sets the characters upright along the text: x runs down the
	// column, and each character is turned back a quarter turn, centred
	// on the x axis in its cell.
	Vertical bool
}

// Bounds returns the bounding box of the text in the drawing.
func (t *Text) Bounds() Rect {
	r := Rect{Point{0, -t.Desc}, Point{t.Advance, t.Asc}, true}
	if t.Vertical {
		r = Rect{Point{0, -0.5}, Point{t.Advance, 0.5}, true}
	}
	return r.Transform(t.M)
}

// Fonts measures text with the fonts of a document.
type Fonts struct {
	Set *fontset.Set
}

func (f *Fonts) face(spec FontSpec, r rune) *fontset.Choice {
	return f.Set.FaceFor(spec.Family, spec.ea(), spec.Bold, spec.Italic, r)
}

// Measure returns the advance of s in ems.
func (f *Fonts) Measure(spec FontSpec, s string) float64 {
	w := 0.0
	for _, r := range s {
		w += f.RuneAdvance(spec, r)
	}
	return w
}

// RuneAdvance returns the advance of a character in ems.
func (f *Fonts) RuneAdvance(spec FontSpec, r rune) float64 {
	if r == '\t' {
		r = ' '
	}
	return f.Set.Advance(f.face(spec, r), r)
}

// Metrics returns the ascent, descent and cap height in ems of the font
// that draws Latin text.
func (f *Fonts) Metrics(spec FontSpec) (asc, desc, capHeight float64) {
	fc := f.face(spec, 'H')
	return fc.Asc, fc.Desc, fc.CapHeight
}

// NewText measures a text; m maps its em space into the drawing.
func (f *Fonts) NewText(spec FontSpec, s string, color bdf.Color, m canvas.Matrix) *Text {
	asc, desc, capH := f.Metrics(spec)
	return &Text{S: s, Font: spec, Color: color, M: m, Advance: f.Measure(spec, s), Asc: asc, Desc: desc, Cap: capH}
}

// NewCellText measures a text of fixed-pitch characters: cells are the
// advances of its characters and spacing the space after each, in ems.
func (f *Fonts) NewCellText(spec FontSpec, s string, cells []float64, spacing float64, color bdf.Color, m canvas.Matrix) *Text {
	t := f.NewText(spec, s, color, m)
	t.Cells, t.Spacing = cells, spacing
	t.Advance = 0
	for _, c := range cells {
		t.Advance += c + spacing
	}
	return t
}
