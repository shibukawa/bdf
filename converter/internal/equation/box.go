package equation

import (
	"math"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/fontset"
)

// Box is a laid out formula or part of one: its width, its height above
// and depth below the baseline (points), and what it draws, relative to
// the left end of its baseline (y grows downwards).
type Box struct {
	W, H, D float64
	// Italic is the italic correction of a box that ends in a slanted
	// glyph: where a superscript goes past its width.
	Italic float64

	items []item
	class Class
	// what attaching scripts and accents looks at
	single    bool    // one glyph, not grown
	accentX   float64 // where an accent attaches above it (x); NaN for the middle
	largeOp   bool    // an (embellished) large operator
	movable   bool    // … whose limits go to the side outside display style
	stretchH  bool    // an operator grown horizontally (a brace, an arrow)
	accentLow bool    // an accent glyph that sits at the height made for x-height bases
}

// item kinds
const (
	iGlyph = iota
	iRule
	iStroke
)

// item is something a box draws.
type item struct {
	kind  uint8
	x, y  float64 // glyph: the start of its baseline; rule: top left; stroke: origin
	fc    *fontset.Choice
	size  float64
	text  string
	adv   float64
	sx    float64 // horizontal and vertical scale of a grown glyph drawn without a formula font
	sy    float64
	w, h  float64 // rule size; stroke line width in w
	path  *bdf.Path
	color bdf.Color
	// alt is the text of a glyph drawn with a private use character
	alt string
	// ink is the extent of what the item draws relative to (x, y): left,
	// top, right, bottom
	ink [4]float64
}

func newBox() *Box { return &Box{accentX: math.NaN(), class: Ord} }

// add draws c into b with its origin at (x, y).
func (b *Box) add(c *Box, x, y float64) {
	for _, it := range c.items {
		it.x += x
		it.y += y
		b.items = append(b.items, it)
	}
}

// grow extends b's height and depth to hold c placed with its baseline at
// y (positive: lower).
func (b *Box) grow(c *Box, y float64) {
	b.H = math.Max(b.H, c.H-y)
	b.D = math.Max(b.D, c.D+y)
}

// rule adds a filled rectangle.
func (b *Box) rule(x, y, w, h float64, color bdf.Color) {
	if w <= 0 || h <= 0 {
		return
	}
	b.items = append(b.items, item{kind: iRule, x: x, y: y, w: w, h: h, color: color, ink: [4]float64{0, 0, w, h}})
}

// stroke adds a stroked path (drawn relative to x, y) whose points lie in
// ink (left, top, right, bottom).
func (b *Box) stroke(p *bdf.Path, x, y, width float64, color bdf.Color, ink [4]float64) {
	ink = [4]float64{ink[0] - width, ink[1] - width, ink[2] + width, ink[3] + width}
	b.items = append(b.items, item{kind: iStroke, x: x, y: y, w: width, path: p, color: color, ink: ink})
}

// accentAt returns where an accent attaches above b.
func (b *Box) accentAt() float64 {
	if math.IsNaN(b.accentX) {
		return b.W / 2
	}
	return b.accentX
}

// Width, Height and Depth are the box's dimensions in points.
func (b *Box) Width() float64  { return b.W }
func (b *Box) Height() float64 { return b.H }
func (b *Box) Depth() float64  { return b.D }

// Ink returns the extent of what b draws relative to its origin (y grows
// downwards), which may reach past its box (italic overhangs, accents).
func (b *Box) Ink() bdf.Rect {
	x0, y0, x1, y1 := 0.0, -b.H, b.W, b.D
	for _, it := range b.items {
		x0, y0 = math.Min(x0, it.x+it.ink[0]), math.Min(y0, it.y+it.ink[1])
		x1, y1 = math.Max(x1, it.x+it.ink[2]), math.Max(y1, it.y+it.ink[3])
	}
	return bdf.Rect{X: float32(x0), Y: float32(y0), W: float32(x1 - x0), H: float32(y1 - y0)}
}
