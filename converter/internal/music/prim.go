package music

import "math"

// primKind is the kind of a drawing primitive.
type primKind uint8

const (
	primGlyph primKind = iota
	primRect
	primFill // a filled outline
	primLine // a stroked line
	primText
	primMark
)

// prim is a drawing primitive of the layout, in staff spaces relative to
// its staff (or system).
type prim struct {
	kind  primKind
	g     *glyph
	scale float64 // glyphs: the size relative to the staff (1, or graceScale)
	sx    float64 // glyphs: an extra horizontal scale (braces are stretched)
	sy    float64 // glyphs: an extra vertical scale
	x, y  float64
	w, h  float64  // rects: size; lines: the end point is (w, h)
	width float64  // lines: thickness
	dash  float64  // lines: dash length (0 solid)
	out   *outline // fills
	text  *textLine
	align int // text: 0 left, 1 centre, 2 right
	mark  byte
	alt   string
}

// skyline is the outline of what has been drawn on one side of a staff:
// for stretches of x the furthest y reached, to place marks clear of it.
type skyline struct {
	above bool
	segs  []skySeg
}

type skySeg struct{ x0, x1, y float64 }

// add records a box [x0, x1] reaching y.
func (s *skyline) add(x0, x1, y float64) {
	if x1 < x0 {
		x0, x1 = x1, x0
	}
	s.segs = append(s.segs, skySeg{x0, x1, y})
}

// at returns the furthest y over [x0, x1], or base when nothing is there.
func (s *skyline) at(x0, x1, base float64) float64 {
	y := base
	for _, g := range s.segs {
		if g.x1 < x0 || g.x0 > x1 {
			continue
		}
		if s.above {
			y = min(y, g.y)
		} else {
			y = max(y, g.y)
		}
	}
	return y
}

// drawing collects the primitives of a staff with their extents.
type drawing struct {
	prims []prim
	top   *skyline
	bot   *skyline
}

func newDrawing() *drawing {
	return &drawing{top: &skyline{above: true}, bot: &skyline{}}
}

// glyph draws g with its origin at (x, y) and records its box.
func (c *drawing) glyph(g *glyph, x, y, scale float64) {
	c.prims = append(c.prims, prim{kind: primGlyph, g: g, x: x, y: y, scale: scale, sx: 1, sy: 1})
	c.box(x+g.x0*scale, y+g.y0*scale, x+g.x1*scale, y+g.y1*scale)
}

// rect fills a rectangle.
func (c *drawing) rect(x, y, w, h float64) {
	c.prims = append(c.prims, prim{kind: primRect, x: x, y: y, w: w, h: h})
	c.box(x, y, x+w, y+h)
}

// fill fills an outline.
func (c *drawing) fill(o *outline) {
	c.prims = append(c.prims, prim{kind: primFill, out: o})
	x0, y0, x1, y1 := o.bounds()
	c.box(x0, y0, x1, y1)
}

// line strokes a line of thickness w.
func (c *drawing) line(x0, y0, x1, y1, w, dash float64) {
	c.prims = append(c.prims, prim{kind: primLine, x: x0, y: y0, w: x1, h: y1, width: w, dash: dash})
	c.box(min(x0, x1), min(y0, y1)-w/2, max(x0, x1), max(y0, y1)+w/2)
}

// text draws a line of text with its baseline at y; align places x at its
// left, centre or right. sp converts the text's pt to staff spaces.
func (c *drawing) text(t *textLine, x, y float64, align int, sp float64) {
	if t == nil || len(t.runs) == 0 {
		return
	}
	c.prims = append(c.prims, prim{kind: primText, text: t, x: x, y: y, align: align})
	w := t.width / sp
	x0 := x - w*float64(align)/2
	c.box(x0, y-t.asc/sp, x0+w, y+t.desc/sp)
}

// mark adds a structure mark.
func (c *drawing) markText(kind byte, payload string) {
	c.prims = append(c.prims, prim{kind: primMark, mark: kind, alt: payload})
}

// box records a drawn box in the skylines.
func (c *drawing) box(x0, y0, x1, y1 float64) {
	c.top.add(x0, x1, y0)
	c.bot.add(x0, x1, y1)
}

// bounds returns the bounding box of an outline's points.
func (o *outline) bounds() (x0, y0, x1, y1 float64) {
	x0, y0, x1, y1 = 1e9, 1e9, -1e9, -1e9
	for i := 0; i+1 < len(o.args); i += 2 {
		x, y := o.args[i], o.args[i+1]
		x0, x1 = min(x0, x), max(x1, x)
		y0, y1 = min(y0, y), max(y1, y)
	}
	return
}

// curve is a tie or slur: a crescent from (x0, y0) to (x1, y1) whose
// middle line bulges by h (negative: up), end thick at the ends and mid at
// the middle.
func curve(x0, y0, x1, y1, h, end, mid float64) *outline {
	s := math.Copysign(1, h)
	eo, ei := s*end/2, -s*end/2
	// a cubic whose control points stand at a quarter and three quarters
	// of the way reaches e/4 + 3c/4 at its middle
	co := (h + s*mid/2 - eo/4) * 4 / 3
	ci := (h - s*mid/2 - ei/4) * 4 / 3
	dx, dy := x1-x0, y1-y0
	return &outline{verbs: []byte{0, 3, 1, 3, 4}, args: []float64{
		x0, y0 + eo,
		x0 + dx*0.25, y0 + dy*0.25 + co, x0 + dx*0.75, y0 + dy*0.75 + co, x1, y1 + eo,
		x1, y1 + ei,
		x0 + dx*0.75, y0 + dy*0.75 + ci, x0 + dx*0.25, y0 + dy*0.25 + ci, x0, y0 + ei}}
}

// quad is a filled quadrilateral (a beam).
func quad(x0, y0, x1, y1, x2, y2, x3, y3 float64) *outline {
	return &outline{verbs: []byte{0, 1, 1, 1, 4}, args: []float64{x0, y0, x1, y1, x2, y2, x3, y3}}
}
