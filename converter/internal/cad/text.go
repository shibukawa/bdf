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
	// Advance is the natural advance of S in ems.
	Advance float64
	// Pos, when not nil, places each character of S: its x in ems
	// (fixed-pitch text).
	Pos []float64
	// Asc, Desc and Cap are the ascent, descent and cap height of the font
	// in ems.
	Asc, Desc, Cap float64

	Underline, Overline, Strike bool
	Break                       Break
}

// Bounds returns the bounding box of the text in the drawing.
func (t *Text) Bounds() Rect {
	w := t.Advance
	if n := len(t.Pos); n > 0 {
		w = max(w, t.Pos[n-1]+1)
	}
	r := Rect{Point{0, -t.Desc}, Point{w, t.Asc}, true}
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
