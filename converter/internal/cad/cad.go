// Package cad draws the 2D drawings of CAD formats (DXF, JWW, SXF). A
// reader puts what a drawing shows into a Drawing, in the drawing's own
// coordinates (y up, float64, so that survey coordinates keep their
// precision), and Plotter writes the Drawing onto a page the way a plotter
// would: line widths are widths on paper, whatever the scale.
//
// Curves are cubic Bézier curves (see Path), text is measured and drawn
// with the fonts of the document (fontset), and fills can be clipped
// groups, which is how hatch patterns are drawn (see Hatch).
package cad

import (
	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// Line caps and joins (bdf values).
const (
	CapButt   byte = 0
	CapRound  byte = 1
	CapSquare byte = 2
	JoinMiter byte = 0
	JoinRound byte = 1
	JoinBevel byte = 2
)

// Pen is how a path is stroked.
type Pen struct {
	Color bdf.Color
	// Width is the width of the line on paper, in pt; 0 is the thinnest
	// line (Plotter.Thin).
	Width float64
	// WorldWidth, when positive, is the width in drawing units instead (the
	// wide polylines of DXF).
	WorldWidth float64
	// Dash is the dash pattern in drawing units: dash, gap, dash, gap, ...
	// A dash of 0 is a dot. nil draws a continuous line.
	Dash []float64
	// DashOffset shifts the pattern along the path (drawing units).
	DashOffset float64
	Cap, Join  byte
}

// Fill is a solid color or a gradient.
type Fill struct {
	Color    bdf.Color
	Gradient *Gradient
}

// Gradient is a gradient in drawing coordinates.
type Gradient struct {
	Radial bool
	// linear: from P0 to P1; radial: circles around P0 (radius R0) and P1
	// (radius R1)
	P0, P1 Point
	R0, R1 float64
	Stops  []bdf.Stop
}

type kind uint8

const (
	kStroke kind = iota
	kFill
	kText
	kGroup
)

// Item is one thing a Drawing draws.
type Item struct {
	kind    kind
	path    *Path
	pen     Pen
	fill    Fill
	evenOdd bool
	text    *Text
	clip    *Path // group: the clip (nil: none)
	items   []Item
	bounds  Rect
}

// Text returns the string of a text item ("" for other items).
func (it Item) Text() string {
	if it.kind == kText && it.text != nil {
		return it.text.S
	}
	return ""
}

// Drawing is a display list in drawing coordinates.
type Drawing struct {
	Items []Item
	open  []*Item // groups being filled
}

func (d *Drawing) add(it Item) {
	if n := len(d.open); n > 0 {
		g := d.open[n-1]
		g.items = append(g.items, it)
		return
	}
	d.Items = append(d.Items, it)
}

// Stroke draws a path with a pen.
func (d *Drawing) Stroke(p *Path, pen Pen) {
	if p.Empty() {
		return
	}
	d.add(Item{kind: kStroke, path: p, pen: pen, bounds: p.Bounds()})
}

// Fill fills a path; evenOdd selects the even-odd rule (nonzero
// otherwise).
func (d *Drawing) Fill(p *Path, f Fill, evenOdd bool) {
	if p.Empty() {
		return
	}
	d.add(Item{kind: kFill, path: p, fill: f, evenOdd: evenOdd, bounds: p.Bounds()})
}

// Text draws a line of text.
func (d *Drawing) Text(t *Text) {
	if t == nil || t.S == "" {
		return
	}
	d.add(Item{kind: kText, text: t, bounds: t.Bounds()})
}

// Begin opens a group clipped by clip (even-odd; nil: no clip) that holds
// what is drawn until End.
func (d *Drawing) Begin(clip *Path) {
	it := &Item{kind: kGroup, clip: clip}
	d.open = append(d.open, it)
}

// End closes the group opened last; an empty group is dropped.
func (d *Drawing) End() {
	n := len(d.open)
	if n == 0 {
		return
	}
	g := d.open[n-1]
	d.open = d.open[:n-1]
	if len(g.items) == 0 {
		return
	}
	for _, it := range g.items {
		g.bounds = g.bounds.Union(it.bounds)
	}
	if g.clip != nil {
		g.bounds = intersect(g.bounds, g.clip.Bounds())
	}
	d.add(*g)
}

func intersect(a, b Rect) Rect {
	if !a.ok || !b.ok {
		return Rect{}
	}
	r := Rect{Point{max(a.Min.X, b.Min.X), max(a.Min.Y, b.Min.Y)}, Point{min(a.Max.X, b.Max.X), min(a.Max.Y, b.Max.Y)}, true}
	if r.Min.X > r.Max.X || r.Min.Y > r.Max.Y {
		return Rect{}
	}
	return r
}

// Bounds returns the bounding box of what the drawing draws (the widths
// of lines left out).
func (d *Drawing) Bounds() Rect {
	var r Rect
	for _, it := range d.Items {
		r = r.Union(it.bounds)
	}
	return r
}

// Empty reports whether the drawing draws nothing.
func (d *Drawing) Empty() bool { return len(d.Items) == 0 }

// Append adds the items of another drawing transformed by m, clipped by
// clip (in this drawing's coordinates; nil: no clip): the view of a
// viewport.
func (d *Drawing) Append(src *Drawing, m canvas.Matrix, clip *Path) {
	d.Begin(clip)
	s := Scale(m)
	for _, it := range src.Items {
		d.add(transformItem(it, m, s))
	}
	d.End()
}

func transformItem(it Item, m canvas.Matrix, s float64) Item {
	out := it
	if it.path != nil {
		out.path = it.path.Transform(m)
	}
	if it.clip != nil {
		out.clip = it.clip.Transform(m)
	}
	out.bounds = it.bounds.Transform(m)
	switch it.kind {
	case kStroke:
		out.pen.WorldWidth *= s
		out.pen.DashOffset *= s
		if it.pen.Dash != nil {
			out.pen.Dash = make([]float64, len(it.pen.Dash))
			for i, v := range it.pen.Dash {
				out.pen.Dash[i] = v * s
			}
		}
	case kFill:
		if g := it.fill.Gradient; g != nil {
			ng := *g
			ng.P0, ng.P1 = Apply(m, g.P0), Apply(m, g.P1)
			ng.R0, ng.R1 = g.R0*s, g.R1*s
			out.fill.Gradient = &ng
		}
	case kText:
		t := *it.text
		t.M = m.Mul(t.M)
		out.text = &t
	case kGroup:
		out.items = make([]Item, len(it.items))
		for i, c := range it.items {
			out.items[i] = transformItem(c, m, s)
		}
	}
	return out
}
