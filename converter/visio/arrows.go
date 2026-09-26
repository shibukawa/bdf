package visio

import (
	"math"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// Arrowheads (BeginArrow and EndArrow 1-45) are drawn from the shapes
// [MS-VSDX] pictures for them. They are laid out in units of the
// arrowhead's length from its tip, with x running back along the line and
// y across it; the line is cut back under the parts that cover its end.

// arrowPen draws one arrowhead.
type arrowPen struct {
	cv         *canvas.Canvas
	tip        pt
	back, perp pt // unit vectors
	size       float64
	l          *lineStyle
}

func (a *arrowPen) at(x, y float64) pt {
	return a.tip.add(a.back.mul(x * a.size)).add(a.perp.mul(y * a.size))
}

func (a *arrowPen) path(closed bool, xy ...float64) *bdf.Path {
	p := &bdf.Path{}
	for i := 0; i+1 < len(xy); i += 2 {
		q := a.at(xy[i], xy[i+1])
		if i == 0 {
			p.MoveTo(f32(q.x), f32(q.y))
		} else {
			p.LineTo(f32(q.x), f32(q.y))
		}
	}
	if closed {
		p.Close()
	}
	return p
}

// stroke draws an open polyline with the line's pen.
func (a *arrowPen) stroke(xy ...float64) {
	a.cv.Obj.StrokePath(a.cv.Obj.AddPath(a.path(false, xy...)))
}

func (a *arrowPen) fill(xy ...float64) {
	a.cv.Obj.FillPath(a.cv.Obj.AddPath(a.path(true, xy...)), 0)
}

func (a *arrowPen) outline(xy ...float64) {
	a.cv.Obj.StrokePath(a.cv.Obj.AddPath(a.path(true, xy...)))
}

func (a *arrowPen) tick(x float64) { a.stroke(x, -0.45, x, 0.45) }

func (a *arrowPen) circle(cx, r float64, filled bool) {
	c := a.at(cx, 0)
	p := &bdf.Path{}
	p.Ellipse(f32(c.x), f32(c.y), f32(r*a.size), f32(r*a.size), 0, 0, 2*math.Pi, false)
	ref := a.cv.Obj.AddPath(p)
	if filled {
		a.cv.Obj.FillPath(ref, 0)
	} else {
		a.cv.Obj.StrokePath(ref)
	}
}

func (a *arrowPen) square(x0, half float64, filled bool) {
	xy := []float64{x0, -half, x0 + 2*half, -half, x0 + 2*half, half, x0, half}
	if filled {
		a.fill(xy...)
	} else {
		a.outline(xy...)
	}
}

// curve approximates a quadratic curve from (x0,y0) through the control
// point to (x1,y1) as polyline points.
func curve(x0, y0, cx, cy, x1, y1 float64) []float64 {
	var out []float64
	for i := 0; i <= 8; i++ {
		t := float64(i) / 8
		u := 1 - t
		out = append(out, u*u*x0+2*u*t*cx+t*t*x1, u*u*y0+2*u*t*cy+t*t*y1)
	}
	return out
}

func join(parts ...[]float64) []float64 {
	var out []float64
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

// arrowDef is an arrowhead: how far the line is cut back, and its drawing.
type arrowDef struct {
	trim float64
	draw func(a *arrowPen)
}

var (
	bulgedSides = join(curve(0, 0, 0.45, 0.4, 1, 0.45), []float64{1, -0.45}, curve(1, -0.45, 0.45, -0.4, 0, 0))
	sweptSides  = join(curve(0, 0, 0.4, 0.1, 0.7, 0.42), curve(0.7, 0.42, 0.5, 0, 0.7, -0.42), curve(0.7, -0.42, 0.4, -0.1, 0, 0))
	diamond     = []float64{0.3, 0, 0.8, 0.3, 1.3, 0, 0.8, -0.3}
)

var arrowDefs = [...]arrowDef{
	1: {0, func(a *arrowPen) { a.stroke(0.6, -0.4, 0, 0, 0.6, 0.4) }},
	2: {0.35, func(a *arrowPen) { a.fill(0, 0, 0.45, 0.22, 0.38, 0, 0.45, -0.22) }},
	3: {0, func(a *arrowPen) { a.stroke(0.85, -0.38, 0, 0, 0.85, 0.38) }},
	4: {0.9, func(a *arrowPen) { a.fill(0, 0, 1, 0.375, 1, -0.375) }},
	5: {0.8, func(a *arrowPen) { a.fill(0, 0, 1, 0.3, 0.85, 0, 1, -0.3) }},
	6: {0.9, func(a *arrowPen) { a.fill(bulgedSides...) }},
	7: {0, func(a *arrowPen) {
		a.stroke(curve(0.6, -0.45, 0.35, -0.08, 0, 0)...)
		a.stroke(curve(0.6, 0.45, 0.35, 0.08, 0, 0)...)
	}},
	8:  {0.5, func(a *arrowPen) { a.fill(0, 0, 0.65, 0.25, 0.5, 0, 0.65, -0.25) }},
	9:  {0, func(a *arrowPen) { a.stroke(-0.25, 0.5, 0.25, -0.5) }},
	10: {0.3, func(a *arrowPen) { a.circle(0.3, 0.3, true) }},
	11: {0.3, func(a *arrowPen) { a.square(0, 0.3, true) }},
	12: {0, func(a *arrowPen) { a.stroke(1.1, -0.35, 0, 0, 1.1, 0.35) }},
	13: {1.2, func(a *arrowPen) { a.fill(0, 0, 1.3, 0.4, 1.3, -0.4) }},
	14: {1, func(a *arrowPen) { a.outline(0, 0, 1, 0.4, 1, -0.4) }},
	15: {0.45, func(a *arrowPen) { a.outline(0, 0, 0.45, 0.3, 0.45, -0.3) }},
	16: {0.75, func(a *arrowPen) { a.outline(0, 0, 0.75, 0.4, 0.75, -0.4) }},
	17: {0.6, func(a *arrowPen) { a.outline(0, 0, 0.8, 0.4, 0.6, 0, 0.8, -0.4) }},
	18: {1, func(a *arrowPen) { a.outline(bulgedSides...) }},
	19: {0.5, func(a *arrowPen) { a.outline(sweptSides...) }},
	20: {0.6, func(a *arrowPen) { a.circle(0.3, 0.3, false) }},
	21: {0.6, func(a *arrowPen) { a.square(0, 0.3, false) }},
	22: {1.2, func(a *arrowPen) { a.outline(0, 0, 0.6, 0.35, 1.2, 0, 0.6, -0.35) }},
	23: {0, func(a *arrowPen) { a.stroke(0.1, -0.5, 0.6, 0.5) }},
	24: {0, func(a *arrowPen) { a.tick(0.35) }},
	25: {0, func(a *arrowPen) { a.tick(0.35); a.tick(0.55) }},
	26: {0, func(a *arrowPen) { a.tick(0.35); a.tick(0.55); a.tick(0.75) }},
	27: {0, func(a *arrowPen) { a.stroke(0, -0.4, 0.7, 0, 0, 0.4) }},
	28: {0, func(a *arrowPen) { a.stroke(0, -0.4, 0.7, 0, 0, 0.4); a.tick(0.95) }},
	29: {1.3, func(a *arrowPen) {
		a.stroke(0, -0.4, 0.6, 0, 0, 0.4)
		a.stroke(0, 0, 0.7, 0)
		a.circle(1, 0.3, false)
	}},
	30: {0.9, func(a *arrowPen) { a.tick(0.15); a.stroke(0, 0, 0.3, 0); a.circle(0.6, 0.3, false) }},
	31: {0.6, func(a *arrowPen) { a.circle(0.3, 0.3, false); a.tick(0.85) }},
	32: {0.6, func(a *arrowPen) { a.circle(0.3, 0.3, false); a.tick(0.85); a.tick(1.05) }},
	33: {0.6, func(a *arrowPen) { a.circle(0.3, 0.3, false); a.tick(0.85); a.tick(1.05); a.tick(1.25) }},
	34: {1.3, func(a *arrowPen) { a.circle(0.15, 0.15, false); a.outline(diamond...) }},
	35: {0.3, func(a *arrowPen) { a.circle(0.3, 0.3, true); a.tick(0.85) }},
	36: {0.3, func(a *arrowPen) { a.circle(0.3, 0.3, true); a.tick(0.85); a.tick(1.05) }},
	37: {0.3, func(a *arrowPen) { a.circle(0.3, 0.3, true); a.tick(0.85); a.tick(1.05); a.tick(1.25) }},
	38: {1.3, func(a *arrowPen) { a.circle(0.15, 0.15, true); a.outline(diamond...) }},
	39: {1.1, func(a *arrowPen) { a.fill(0, 0, 0.6, 0.35, 0.6, -0.35); a.fill(0.6, 0, 1.2, 0.35, 1.2, -0.35) }},
	40: {1.2, func(a *arrowPen) { a.outline(0, 0, 0.6, 0.35, 0.6, -0.35); a.outline(0.6, 0, 1.2, 0.35, 1.2, -0.35) }},
	41: {0.7, func(a *arrowPen) { a.circle(0.35, 0.35, false) }},
	42: {0.3, func(a *arrowPen) { a.circle(0.3, 0.3, true) }},
	43: {0, func(a *arrowPen) { a.stroke(0.8, -0.4, 0, 0, 0.8, 0.4); a.stroke(1.3, -0.4, 0.5, 0, 1.3, 0.4) }},
	44: {0, func(a *arrowPen) { a.stroke(0.8, -0.4, 0, 0, 0.8, 0.4); a.tick(0.6) }},
	45: {0, func(a *arrowPen) {
		a.stroke(0.8, -0.4, 0, 0, 0.8, 0.4)
		a.stroke(1.3, -0.4, 0.5, 0, 1.3, 0.4)
		a.tick(1.1)
	}},
}

// arrowSizes scale arrowheads by their size (very small … colossal).
var arrowSizes = [...]float64{0.6, 0.8, 1, 1.4, 1.8, 2.6, 3.6}

// arrowLength is the length of an arrowhead in points.
func arrowLength(size int, width float64) float64 {
	f := 1.0
	if size >= 0 && size < len(arrowSizes) {
		f = arrowSizes[size]
	}
	return f * (4.5 + 2.5*math.Max(width, 0.25))
}

func arrowDefOf(kind int) (arrowDef, bool) {
	if kind <= 0 || kind >= len(arrowDefs) || arrowDefs[kind].draw == nil {
		return arrowDef{}, false
	}
	return arrowDefs[kind], true
}

// drawArrow draws an arrowhead with its tip at tip, where the line
// arrives going in direction dir.
func drawArrow(cv *canvas.Canvas, l *lineStyle, kind, size int, tip, dir pt) {
	def, ok := arrowDefOf(kind)
	if !ok {
		return
	}
	n := dir.len()
	if n == 0 {
		return
	}
	back := dir.mul(-1 / n)
	a := &arrowPen{cv: cv, tip: tip, back: back, perp: pt{-back.y, back.x}, size: arrowLength(size, l.width), l: l}
	cv.Obj.Save()
	cv.Obj.Dash(nil, 0)
	cv.Obj.Line(f32(l.width), l.cap, 0, 10)
	cv.Obj.StrokeColor(l.color)
	cv.Obj.FillColor(l.color)
	def.draw(a)
	cv.Obj.Restore()
}

// trimStart cuts a figure back by d points at its start, along its first
// segment.
func trimStart(sp *subpath, d float64) {
	if d <= 0 || len(sp.segs) == 0 {
		return
	}
	s := &sp.segs[0]
	dir := sp.startDir()
	n := dir.len()
	if n == 0 {
		return
	}
	limit := s.p.dist(sp.start) * 0.9
	d = math.Min(d, limit)
	shift := dir.mul(d / n)
	sp.start = sp.start.add(shift)
	if s.cubic {
		s.c1 = s.c1.add(shift)
	}
}

// trimEnd cuts a figure back by d points at its end.
func trimEnd(sp *subpath, d float64) {
	if d <= 0 || len(sp.segs) == 0 {
		return
	}
	i := len(sp.segs) - 1
	s := &sp.segs[i]
	prev := sp.start
	if i > 0 {
		prev = sp.segs[i-1].p
	}
	dir := sp.endDir()
	n := dir.len()
	if n == 0 {
		return
	}
	d = math.Min(d, s.p.dist(prev)*0.9)
	shift := dir.mul(-d / n)
	s.p = s.p.add(shift)
	if s.cubic {
		s.c2 = s.c2.add(shift)
	}
}
