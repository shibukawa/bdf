// Package shapes reads an object that draws paths only as a list of filled
// and stroked shapes in the units of the object: the transforms applied,
// every path made of moves, lines, quadratic and cubic curves and closes
// (rectangles, ellipses, arcs and rounded rectangles drawn as those), and
// every shape with its color. It is what a renderer that fills and strokes
// paths and does nothing else (raster/ebitenginebdf) draws an object from.
package shapes

import (
	"fmt"
	"math"

	"github.com/shibukawa/bdf"
)

// Shape is a path that is filled or stroked in one color.
type Shape struct {
	// Path holds VerbMove, VerbLine, VerbQuad, VerbCubic and VerbClose
	// only.
	Path   bdf.Path
	Stroke bool
	// Rule is the fill rule of a filled shape (bdf.NonZero, bdf.EvenOdd).
	Rule byte
	// Color is the color with the alpha of the state applied.
	Color bdf.Color
	// the line of a stroked shape, the width in the units of the object
	Width     float64
	Cap, Join byte
	Miter     float64
}

// UnsupportedError is the error of an object that draws more than paths in
// plain colors.
type UnsupportedError struct {
	Op byte // the first instruction that cannot be drawn
}

func (e *UnsupportedError) Error() string {
	return fmt.Sprintf("the object draws more than paths in plain colors: %s", bdf.OpName(e.Op))
}

// maxStates bounds the states saved at a time.
const maxStates = 1 << 16

type matrix [6]float64

func (m matrix) mul(n matrix) matrix {
	return matrix{
		m[0]*n[0] + m[2]*n[1], m[1]*n[0] + m[3]*n[1],
		m[0]*n[2] + m[2]*n[3], m[1]*n[2] + m[3]*n[3],
		m[0]*n[4] + m[2]*n[5] + m[4], m[1]*n[4] + m[3]*n[5] + m[5],
	}
}

func (m matrix) apply(x, y float64) (float64, float64) {
	return m[0]*x + m[2]*y + m[4], m[1]*x + m[3]*y + m[5]
}

type state struct {
	m            matrix
	fill, stroke bdf.Color
	alpha        float64
	width, miter float64
	cap, join    byte
}

// Of returns the shapes an object draws, in drawing order. An object that
// clips, draws text, images, gradients, patterns, groups, masks, shadows,
// dashed lines or other objects, or blends other than source-over, is an
// *UnsupportedError: those are not paths in plain colors.
func Of(o *bdf.ObjectPart) ([]Shape, error) {
	st := state{m: matrix{1, 0, 0, 1, 0, 0}, fill: 0x000000ff, stroke: 0x000000ff, alpha: 1, width: 1, miter: 10}
	var stack []state
	var out []Shape
	var unsupported *UnsupportedError
	f := func(in bdf.Instr, i int) float64 {
		if i < len(in.Args) {
			if v, ok := in.Args[i].(float32); ok {
				return float64(v)
			}
		}
		return 0
	}
	u := func(in bdf.Instr, i int) uint64 {
		if i < len(in.Args) {
			if v, ok := in.Args[i].(uint64); ok {
				return v
			}
		}
		return 0
	}
	path := func(i uint64) *bdf.Path {
		if i >= uint64(len(o.Paths)) {
			return nil
		}
		return o.Paths[i].Inline // nil for a path of a collection, which Deps reports
	}
	fill := func(rule byte, add func(b *builder)) {
		b := &builder{}
		add(b)
		if len(b.out.Verbs) > 0 {
			out = append(out, Shape{Path: b.out, Rule: rule, Color: withAlpha(st.fill, st.alpha)})
		}
	}
	stroke := func(add func(b *builder)) {
		b := &builder{}
		add(b)
		// a line is as wide as the transform makes it where it scales
		// evenly; where it does not, the width is that of the mean scale
		k := math.Sqrt(math.Abs(st.m[0]*st.m[3] - st.m[1]*st.m[2]))
		if len(b.out.Verbs) > 0 && k > 0 {
			out = append(out, Shape{Path: b.out, Stroke: true, Color: withAlpha(st.stroke, st.alpha),
				Width: st.width * k, Cap: st.cap, Join: st.join, Miter: st.miter})
		}
	}
	err := o.Walk(func(in bdf.Instr) {
		if unsupported != nil {
			return
		}
		switch in.Op {
		case bdf.OpSave:
			if len(stack) < maxStates {
				stack = append(stack, st)
			}
		case bdf.OpRestore:
			if n := len(stack); n > 0 {
				st, stack = stack[n-1], stack[:n-1]
			}
		case bdf.OpTransform:
			st.m = st.m.mul(matrix{f(in, 0), f(in, 1), f(in, 2), f(in, 3), f(in, 4), f(in, 5)})
		case bdf.OpTranslate:
			st.m = st.m.mul(matrix{1, 0, 0, 1, f(in, 0), f(in, 1)})
		case bdf.OpScale:
			st.m = st.m.mul(matrix{f(in, 0), 0, 0, f(in, 1), 0, 0})
		case bdf.OpFillColor:
			st.fill = bdf.Color(u(in, 0))
		case bdf.OpStrokeColor:
			st.stroke = bdf.Color(u(in, 0))
		case bdf.OpLine:
			if w := f(in, 0); w > 0 && !math.IsInf(w, 0) {
				st.width = w
			}
			st.cap, st.join = byte(u(in, 1)), byte(u(in, 2))
			if m := f(in, 3); m > 0 && !math.IsInf(m, 0) {
				st.miter = m
			}
		case bdf.OpAlpha:
			if a := f(in, 0); a >= 0 && a <= 1 {
				st.alpha = a
			}
		case bdf.OpFillRect:
			fill(bdf.NonZero, func(b *builder) { b.m = st.m; b.rect(f(in, 0), f(in, 1), f(in, 2), f(in, 3)) })
		case bdf.OpStrokeRect:
			stroke(func(b *builder) { b.m = st.m; b.rect(f(in, 0), f(in, 1), f(in, 2), f(in, 3)) })
		case bdf.OpFillPath:
			if p := path(u(in, 0)); p != nil {
				fill(byte(u(in, 1)), func(b *builder) { b.m = st.m; b.add(p) })
			}
		case bdf.OpStrokePath:
			if p := path(u(in, 0)); p != nil {
				stroke(func(b *builder) { b.m = st.m; b.add(p) })
			}
		case bdf.OpFillPathAt:
			if p := path(u(in, 0)); p != nil {
				fill(byte(u(in, 1)), func(b *builder) {
					b.m = st.m.mul(matrix{1, 0, 0, 1, f(in, 2), f(in, 3)})
					b.add(p)
				})
			}
		case bdf.OpFillPathRun:
			if len(in.Args) < 2 {
				return
			}
			glyphs, _ := in.Args[1].([]bdf.Glyph)
			fill(byte(u(in, 0)), func(b *builder) {
				for _, g := range glyphs {
					if p := path(uint64(g.Path)); p != nil {
						b.m = st.m.mul(matrix{1, 0, 0, 1, float64(g.X), float64(g.Y)})
						b.start()
						b.add(p)
					}
				}
			})
		case bdf.OpDash:
			if segs, _ := in.Args[0].([]float32); len(segs) > 0 {
				unsupported = &UnsupportedError{Op: in.Op}
			}
		case bdf.OpBlend:
			if u(in, 0) != uint64(bdf.BlendSourceOver) {
				unsupported = &UnsupportedError{Op: in.Op}
			}
		case bdf.OpShadow:
			if u(in, 0)&0xff != 0 {
				unsupported = &UnsupportedError{Op: in.Op}
			}
		case bdf.OpFilter:
			if s, _ := in.Args[0].(string); s != "" && s != "none" {
				unsupported = &UnsupportedError{Op: in.Op}
			}
		case bdf.OpFont, bdf.OpTextStyle, bdf.OpSmoothing, bdf.OpLink, bdf.OpMark, bdf.OpExt:
			// what sets the state of text and images, and what draws nothing
		default:
			unsupported = &UnsupportedError{Op: in.Op}
		}
	})
	if err != nil {
		return nil, err
	}
	if unsupported != nil {
		return nil, unsupported
	}
	return out, nil
}

// withAlpha returns c with its alpha scaled.
func withAlpha(c bdf.Color, alpha float64) bdf.Color {
	if alpha >= 1 {
		return c
	}
	a := math.Round(float64(uint8(c)) * alpha)
	return c&^0xff | bdf.Color(uint8(a))
}

// builder makes a path of moves, lines, curves and closes of the paths
// added to it, under a transform, the way Canvas builds a Path2D.
type builder struct {
	m   matrix
	out bdf.Path
	// the current point and the start of the sub-path, before the
	// transform; has is false before the first point
	cx, cy, sx, sy float64
	has            bool
	// pending is set when the current point starts a sub-path that has
	// not been written yet: after a close, a rectangle or a rounded
	// rectangle, which leave the current point where nothing may follow
	pending bool
}

// start forgets the current point: the next path starts one of its own.
func (b *builder) start() { b.has = false }

func (b *builder) moveTo(x, y float64) {
	px, py := b.m.apply(x, y)
	b.out.MoveTo(float32(px), float32(py))
	b.cx, b.cy, b.sx, b.sy, b.has, b.pending = x, y, x, y, true, false
}

// moveLater makes (x, y) the current point, where a sub-path starts if
// anything is drawn from it.
func (b *builder) moveLater(x, y float64) {
	b.cx, b.cy, b.sx, b.sy, b.has, b.pending = x, y, x, y, true, true
}

// begin writes the start of a sub-path that was left for later.
func (b *builder) begin() {
	if b.pending {
		b.moveTo(b.cx, b.cy)
	}
}

func (b *builder) lineTo(x, y float64) {
	if !b.has {
		b.moveTo(x, y)
		return
	}
	b.begin()
	px, py := b.m.apply(x, y)
	b.out.LineTo(float32(px), float32(py))
	b.cx, b.cy = x, y
}

func (b *builder) quadTo(cx, cy, x, y float64) {
	if !b.has {
		b.moveTo(cx, cy)
	}
	b.begin()
	qx, qy := b.m.apply(cx, cy)
	px, py := b.m.apply(x, y)
	b.out.QuadTo(float32(qx), float32(qy), float32(px), float32(py))
	b.cx, b.cy = x, y
}

func (b *builder) cubicTo(c1x, c1y, c2x, c2y, x, y float64) {
	if !b.has {
		b.moveTo(c1x, c1y)
	}
	b.begin()
	ax, ay := b.m.apply(c1x, c1y)
	bx, by := b.m.apply(c2x, c2y)
	px, py := b.m.apply(x, y)
	b.out.CubicTo(float32(ax), float32(ay), float32(bx), float32(by), float32(px), float32(py))
	b.cx, b.cy = x, y
}

// close closes the sub-path; the next one starts where it started.
func (b *builder) close() {
	if b.has && !b.pending {
		b.out.Close()
		b.moveLater(b.sx, b.sy)
	}
}

func (b *builder) rect(x, y, w, h float64) {
	b.moveTo(x, y)
	b.lineTo(x+w, y)
	b.lineTo(x+w, y+h)
	b.lineTo(x, y+h)
	b.close()
	b.moveLater(x, y)
}

// ellipse adds an elliptical arc with Canvas semantics: a line from the
// current point to the start of the arc, angles measured clockwise from
// the x axis (y down) and rotated by rot.
func (b *builder) ellipse(cx, cy, rx, ry, rot, a0, a1 float64, ccw bool) {
	if rx < 0 || ry < 0 {
		return
	}
	const tau = 2 * math.Pi
	var sweep float64
	if !ccw {
		if a1-a0 >= tau {
			sweep = tau
		} else {
			sweep = math.Mod(a1-a0, tau)
			if sweep < 0 {
				sweep += tau
			}
		}
	} else {
		if a0-a1 >= tau {
			sweep = -tau
		} else {
			sweep = -math.Mod(a0-a1, tau)
			if sweep > 0 {
				sweep -= tau
			}
		}
	}
	sin, cos := math.Sincos(rot)
	tr := func(x, y float64) (float64, float64) {
		x, y = rx*x, ry*y
		return cx + x*cos - y*sin, cy + x*sin + y*cos
	}
	s0, c0 := math.Sincos(a0)
	b.lineTo(tr(c0, s0))
	if sweep == 0 || rx == 0 || ry == 0 {
		s1, c1 := math.Sincos(a0 + sweep)
		b.lineTo(tr(c1, s1))
		return
	}
	n := int(math.Ceil(math.Abs(sweep) / (math.Pi / 2)))
	step := sweep / float64(n)
	k := 4.0 / 3 * math.Tan(step/4)
	a := a0
	for i := 0; i < n; i++ {
		e := a + step
		sa, ca := math.Sincos(a)
		se, ce := math.Sincos(e)
		// control points of the unit circle arc, then scaled and rotated
		c1x, c1y := tr(ca-k*sa, sa+k*ca)
		c2x, c2y := tr(ce+k*se, se-k*ce)
		ex, ey := tr(ce, se)
		b.cubicTo(c1x, c1y, c2x, c2y, ex, ey)
		a = e
	}
}

// arcTo adds a line and an arc of radius r tangent to the lines from the
// current point to (x1, y1) and from there to (x2, y2).
func (b *builder) arcTo(x1, y1, x2, y2, r float64) {
	if r < 0 {
		return
	}
	if !b.has {
		b.moveTo(x1, y1)
	}
	b.begin()
	x0, y0 := b.cx, b.cy
	if (x0 == x1 && y0 == y1) || (x1 == x2 && y1 == y2) || r == 0 {
		b.lineTo(x1, y1)
		return
	}
	ax, ay := x0-x1, y0-y1
	bx, by := x2-x1, y2-y1
	la, lb := math.Hypot(ax, ay), math.Hypot(bx, by)
	cross := ax*by - ay*bx
	if math.Abs(cross) < 1e-12*la*lb {
		b.lineTo(x1, y1) // collinear
		return
	}
	theta := math.Acos(math.Max(-1, math.Min(1, (ax*bx+ay*by)/(la*lb))))
	d := r / math.Tan(theta/2)
	t0x, t0y := x1+ax/la*d, y1+ay/la*d
	t1x, t1y := x1+bx/lb*d, y1+by/lb*d
	// the center is on the bisector
	bisx, bisy := ax/la+bx/lb, ay/la+by/lb
	bl := math.Hypot(bisx, bisy)
	h := r / math.Sin(theta/2)
	cx, cy := x1+bisx/bl*h, y1+bisy/bl*h
	b.lineTo(t0x, t0y)
	b.ellipse(cx, cy, r, r, 0, math.Atan2(t0y-cy, t0x-cx), math.Atan2(t1y-cy, t1x-cx), cross > 0)
}

// roundRect adds a rectangle with corners rounded by r (Canvas roundRect
// with one radius).
func (b *builder) roundRect(x, y, w, h, r float64) {
	if w < 0 {
		x, w = x+w, -w
	}
	if h < 0 {
		y, h = y+h, -h
	}
	r = math.Max(0, math.Min(r, math.Min(w, h)/2))
	if r == 0 {
		b.rect(x, y, w, h)
		return
	}
	b.moveTo(x+r, y)
	b.lineTo(x+w-r, y)
	b.ellipse(x+w-r, y+r, r, r, 0, -math.Pi/2, 0, false)
	b.lineTo(x+w, y+h-r)
	b.ellipse(x+w-r, y+h-r, r, r, 0, 0, math.Pi/2, false)
	b.lineTo(x+r, y+h)
	b.ellipse(x+r, y+h-r, r, r, 0, math.Pi/2, math.Pi, false)
	b.lineTo(x, y+r)
	b.ellipse(x+r, y+r, r, r, 0, math.Pi, 3*math.Pi/2, false)
	b.close()
	b.moveLater(x, y)
}

// add adds a path.
func (b *builder) add(p *bdf.Path) {
	a := p.Args
	i := 0
	f := func(k int) float64 {
		if i+k >= len(a) {
			return 0
		}
		return float64(a[i+k])
	}
	for _, v := range p.Verbs {
		switch v {
		case bdf.VerbMove:
			b.moveTo(f(0), f(1))
			i += 2
		case bdf.VerbLine:
			b.lineTo(f(0), f(1))
			i += 2
		case bdf.VerbQuad:
			b.quadTo(f(0), f(1), f(2), f(3))
			i += 4
		case bdf.VerbCubic:
			b.cubicTo(f(0), f(1), f(2), f(3), f(4), f(5))
			i += 6
		case bdf.VerbClose:
			b.close()
		case bdf.VerbRect:
			b.rect(f(0), f(1), f(2), f(3))
			i += 4
		case bdf.VerbEllipse:
			b.ellipse(f(0), f(1), f(2), f(3), f(4), f(5), f(6), f(7) != 0)
			i += 8
		case bdf.VerbArcTo:
			b.arcTo(f(0), f(1), f(2), f(3), f(4))
			i += 5
		case bdf.VerbRoundRect:
			b.roundRect(f(0), f(1), f(2), f(3), f(4))
			i += 5
		}
	}
}
