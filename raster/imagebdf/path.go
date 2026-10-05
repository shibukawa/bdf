package imagebdf

import (
	"image"
	"math"

	"github.com/shibukawa/bdf"
)

// matrix is an affine transform as Canvas 2D keeps it: x' = a x + c y + e,
// y' = b x + d y + f.
type matrix [6]float64

var identity = matrix{1, 0, 0, 1, 0, 0}

// mul returns m × n: n applies first (ctx.transform(n) on a context at m).
func (m matrix) mul(n matrix) matrix {
	return matrix{
		m[0]*n[0] + m[2]*n[1],
		m[1]*n[0] + m[3]*n[1],
		m[0]*n[2] + m[2]*n[3],
		m[1]*n[2] + m[3]*n[3],
		m[0]*n[4] + m[2]*n[5] + m[4],
		m[1]*n[4] + m[3]*n[5] + m[5],
	}
}

func (m matrix) apply(x, y float64) (float64, float64) {
	return m[0]*x + m[2]*y + m[4], m[1]*x + m[3]*y + m[5]
}

func (m matrix) translate(x, y float64) matrix { return m.mul(matrix{1, 0, 0, 1, x, y}) }
func (m matrix) scale(x, y float64) matrix     { return m.mul(matrix{x, 0, 0, y, 0, 0}) }

func (m matrix) det() float64 { return m[0]*m[3] - m[1]*m[2] }

func (m matrix) invert() (matrix, bool) {
	d := m.det()
	if d == 0 || math.IsNaN(d) || math.IsInf(d, 0) {
		return matrix{}, false
	}
	return matrix{
		m[3] / d, -m[1] / d, -m[2] / d, m[0] / d,
		(m[2]*m[5] - m[3]*m[4]) / d, (m[1]*m[4] - m[0]*m[5]) / d,
	}, true
}

// maxScale is the largest factor the transform stretches a length by.
func (m matrix) maxScale() float64 {
	return math.Max(math.Hypot(m[0], m[1]), math.Hypot(m[2], m[3]))
}

// areaScale is the square root of the area factor (the zoom of a uniform
// transform).
func (m matrix) areaScale() float64 { return math.Sqrt(math.Abs(m.det())) }

func (m matrix) finite() bool {
	for _, v := range m {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return false
		}
	}
	return true
}

type point struct{ x, y float64 }

// Segment kinds of a path.
const (
	segLine = iota
	segQuad
	segCubic
)

type segment struct {
	kind int
	p    [6]float64 // control points and end point
}

// subpath is a run of segments from a start point.
type subpath struct {
	start  point
	segs   []segment
	closed bool
}

// path is a Path2D: sub-paths in user space.
type path struct {
	subs []subpath
	// box is the bounding box of the points, once extent has found it (a
	// path is not added to after it is drawn)
	box   [4]float64
	boxed bool
}

func (p *path) cur() *subpath {
	if len(p.subs) == 0 {
		return nil
	}
	return &p.subs[len(p.subs)-1]
}

// current returns the current point, and whether there is one.
func (p *path) current() (point, bool) {
	s := p.cur()
	if s == nil {
		return point{}, false
	}
	if n := len(s.segs); n > 0 {
		g := s.segs[n-1]
		switch g.kind {
		case segLine:
			return point{g.p[0], g.p[1]}, true
		case segQuad:
			return point{g.p[2], g.p[3]}, true
		default:
			return point{g.p[4], g.p[5]}, true
		}
	}
	return s.start, true
}

func (p *path) moveTo(x, y float64) {
	p.boxed = false
	p.subs = append(p.subs, subpath{start: point{x, y}})
}

// ensure starts a sub-path at (x, y) when there is none (Canvas: "ensure
// there is a subpath").
func (p *path) ensure(x, y float64) {
	if len(p.subs) == 0 {
		p.moveTo(x, y)
	}
}

func (p *path) lineTo(x, y float64) {
	if len(p.subs) == 0 {
		p.moveTo(x, y)
		return
	}
	p.boxed = false
	s := p.cur()
	s.segs = append(s.segs, segment{kind: segLine, p: [6]float64{x, y}})
}

func (p *path) quadTo(cx, cy, x, y float64) {
	p.boxed = false
	p.ensure(cx, cy)
	s := p.cur()
	s.segs = append(s.segs, segment{kind: segQuad, p: [6]float64{cx, cy, x, y}})
}

func (p *path) cubicTo(c1x, c1y, c2x, c2y, x, y float64) {
	p.boxed = false
	p.ensure(c1x, c1y)
	s := p.cur()
	s.segs = append(s.segs, segment{kind: segCubic, p: [6]float64{c1x, c1y, c2x, c2y, x, y}})
}

// close closes the current sub-path; the next one starts where it started.
func (p *path) close() {
	if s := p.cur(); s != nil {
		s.closed = true
		p.subs = append(p.subs, subpath{start: s.start})
	}
}

func (p *path) rect(x, y, w, h float64) {
	p.moveTo(x, y)
	p.lineTo(x+w, y)
	p.lineTo(x+w, y+h)
	p.lineTo(x, y+h)
	p.close()
	p.moveTo(x, y)
}

// ellipse adds an elliptical arc with Canvas semantics: a line from the
// current point to the start of the arc, angles measured clockwise from
// the x axis (y down) and rotated by rot.
func (p *path) ellipse(cx, cy, rx, ry, rot, a0, a1 float64, ccw bool) {
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
	at := func(a float64) (float64, float64) {
		s, c := math.Sincos(a)
		x, y := rx*c, ry*s
		return cx + x*cos - y*sin, cy + x*sin + y*cos
	}
	sx, sy := at(a0)
	if _, ok := p.current(); ok {
		p.lineTo(sx, sy)
	} else {
		p.moveTo(sx, sy)
	}
	if sweep == 0 || rx == 0 || ry == 0 {
		ex, ey := at(a0 + sweep)
		p.lineTo(ex, ey)
		return
	}
	n := int(math.Ceil(math.Abs(sweep) / (math.Pi / 2)))
	step := sweep / float64(n)
	k := 4.0 / 3 * math.Tan(step/4)
	a := a0
	for i := 0; i < n; i++ {
		b := a + step
		sa, ca := math.Sincos(a)
		sb, cb := math.Sincos(b)
		// control points of the unit circle arc, then scaled and rotated
		p1x, p1y := ca-k*sa, sa+k*ca
		p2x, p2y := cb+k*sb, sb-k*cb
		tr := func(x, y float64) (float64, float64) {
			x, y = rx*x, ry*y
			return cx + x*cos - y*sin, cy + x*sin + y*cos
		}
		c1x, c1y := tr(p1x, p1y)
		c2x, c2y := tr(p2x, p2y)
		ex, ey := tr(cb, sb)
		p.cubicTo(c1x, c1y, c2x, c2y, ex, ey)
		a = b
	}
}

// arcTo adds a line and an arc of radius r tangent to the lines from the
// current point to (x1, y1) and from there to (x2, y2).
func (p *path) arcTo(x1, y1, x2, y2, r float64) {
	if r < 0 {
		return
	}
	p0, ok := p.current()
	if !ok {
		p.moveTo(x1, y1)
		p0 = point{x1, y1}
	}
	if (p0.x == x1 && p0.y == y1) || (x1 == x2 && y1 == y2) || r == 0 {
		p.lineTo(x1, y1)
		return
	}
	ax, ay := p0.x-x1, p0.y-y1
	bx, by := x2-x1, y2-y1
	la, lb := math.Hypot(ax, ay), math.Hypot(bx, by)
	cross := ax*by - ay*bx
	if math.Abs(cross) < 1e-12*la*lb {
		p.lineTo(x1, y1) // collinear
		return
	}
	cosT := (ax*bx + ay*by) / (la * lb)
	theta := math.Acos(math.Max(-1, math.Min(1, cosT)))
	d := r / math.Tan(theta/2)
	t0x, t0y := x1+ax/la*d, y1+ay/la*d
	t1x, t1y := x1+bx/lb*d, y1+by/lb*d
	// centre: along the bisector
	bisx, bisy := ax/la+bx/lb, ay/la+by/lb
	bl := math.Hypot(bisx, bisy)
	h := r / math.Sin(theta/2)
	cx, cy := x1+bisx/bl*h, y1+bisy/bl*h
	a0 := math.Atan2(t0y-cy, t0x-cx)
	a1 := math.Atan2(t1y-cy, t1x-cx)
	p.lineTo(t0x, t0y)
	p.ellipse(cx, cy, r, r, 0, a0, a1, cross > 0)
}

// roundRect adds a rectangle with corners rounded by r (Canvas roundRect
// with one radius).
func (p *path) roundRect(x, y, w, h, r float64) {
	if w < 0 {
		x, w = x+w, -w
	}
	if h < 0 {
		y, h = y+h, -h
	}
	r = math.Max(0, math.Min(r, math.Min(w, h)/2))
	if r == 0 {
		p.rect(x, y, w, h)
		return
	}
	p.moveTo(x+r, y)
	p.lineTo(x+w-r, y)
	p.ellipse(x+w-r, y+r, r, r, 0, -math.Pi/2, 0, false)
	p.lineTo(x+w, y+h-r)
	p.ellipse(x+w-r, y+h-r, r, r, 0, 0, math.Pi/2, false)
	p.lineTo(x+r, y+h)
	p.ellipse(x+r, y+h-r, r, r, 0, math.Pi/2, math.Pi, false)
	p.lineTo(x, y+r)
	p.ellipse(x+r, y+r, r, r, 0, math.Pi, 3*math.Pi/2, false)
	p.close()
	p.moveTo(x, y)
}

// buildPath converts a bdf path to a path.
func buildPath(bp *bdf.Path) *path {
	p := &path{}
	a := bp.Args
	i := 0
	f := func(k int) float64 {
		if i+k >= len(a) {
			return 0
		}
		return float64(a[i+k])
	}
	for _, v := range bp.Verbs {
		switch v {
		case bdf.VerbMove:
			p.moveTo(f(0), f(1))
			i += 2
		case bdf.VerbLine:
			p.lineTo(f(0), f(1))
			i += 2
		case bdf.VerbQuad:
			p.quadTo(f(0), f(1), f(2), f(3))
			i += 4
		case bdf.VerbCubic:
			p.cubicTo(f(0), f(1), f(2), f(3), f(4), f(5))
			i += 6
		case bdf.VerbClose:
			p.close()
		case bdf.VerbRect:
			p.rect(f(0), f(1), f(2), f(3))
			i += 4
		case bdf.VerbEllipse:
			p.ellipse(f(0), f(1), f(2), f(3), f(4), f(5), f(6), f(7) != 0)
			i += 8
		case bdf.VerbArcTo:
			p.arcTo(f(0), f(1), f(2), f(3), f(4))
			i += 5
		case bdf.VerbRoundRect:
			p.roundRect(f(0), f(1), f(2), f(3), f(4))
			i += 5
		}
	}
	return p
}

// polyline is a flattened sub-path.
type polyline struct {
	pts    []point
	closed bool
}

// extent returns the bounding box of the points of the path (control points
// included, sub-paths without segments left out, which draw nothing), and
// whether it has one: whether the path has points and all are finite.
func (p *path) extent() (box [4]float64, ok bool) {
	if !p.boxed {
		x0, y0, x1, y1 := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
		add := func(x, y float64) {
			x0, y0, x1, y1 = math.Min(x0, x), math.Min(y0, y), math.Max(x1, x), math.Max(y1, y)
		}
		for i := range p.subs {
			s := &p.subs[i]
			if len(s.segs) == 0 {
				continue
			}
			add(s.start.x, s.start.y)
			for j := range s.segs {
				g := &s.segs[j]
				add(g.p[0], g.p[1])
				if g.kind != segLine {
					add(g.p[2], g.p[3])
				}
				if g.kind == segCubic {
					add(g.p[4], g.p[5])
				}
			}
		}
		p.box = [4]float64{x0, y0, x1, y1}
		p.boxed = true
	}
	b := p.box
	ok = b[0] <= b[2] && b[1] <= b[3] && !math.IsInf(b[0], 0) && !math.IsInf(b[1], 0) && !math.IsInf(b[2], 0) && !math.IsInf(b[3], 0)
	return b, ok
}

// flatten turns the path into polylines through m, with curves split until
// they deviate from their chords by less than tol (in the output space).
func (p *path) flatten(m matrix, tol float64) []polyline {
	return p.flattenTo(&scratch{}, make([]polyline, 0, len(p.subs)), m, tol)
}

// flattenTo is flatten with the memory of a renderer: the polylines are
// added to out and their points to those sc holds, where they last until
// sc is reset for the next drawing.
func (p *path) flattenTo(sc *scratch, out []polyline, m matrix, tol float64) []polyline {
	if sc.limited {
		return out
	}
	buf := sc.pts
	defer func() { sc.pts = buf }()
	for _, s := range p.subs {
		if len(buf) >= maxPathPoints {
			sc.stop()
			return out
		}
		start := len(buf)
		x0, y0 := m.apply(s.start.x, s.start.y)
		buf = append(buf, point{x0, y0})
		cx, cy := x0, y0
		for _, g := range s.segs {
			switch g.kind {
			case segLine:
				x, y := m.apply(g.p[0], g.p[1])
				buf = append(buf, point{x, y})
				cx, cy = x, y
			case segQuad:
				qx, qy := m.apply(g.p[0], g.p[1])
				x, y := m.apply(g.p[2], g.p[3])
				buf = flattenQuad(buf, cx, cy, qx, qy, x, y, tol)
				cx, cy = x, y
			case segCubic:
				ax, ay := m.apply(g.p[0], g.p[1])
				bx, by := m.apply(g.p[2], g.p[3])
				x, y := m.apply(g.p[4], g.p[5])
				buf = flattenCubic(buf, cx, cy, ax, ay, bx, by, x, y, tol)
				cx, cy = x, y
			}
			if len(buf) > maxPathPoints {
				buf = buf[:maxPathPoints]
				sc.stop()
				return out
			}
		}
		// the polyline keeps the array its points are in now: points added
		// later go after them, or into a larger array
		out = append(out, polyline{pts: buf[start:len(buf):len(buf)], closed: s.closed})
	}
	return out
}

// reset lets go of the polylines of the drawing before.
func (sc *scratch) reset() {
	sc.polys, sc.pts = keep(sc.polys), keep(sc.pts)
	sc.limited, sc.work = false, 0
}

// outside reports whether a box in user space (x0, y0, x1, y1), with reach
// more on every side, lies outside bounds under m by more than a pixel:
// nothing drawn in it shows. A box whose place cannot be told (a corner
// that is no number) is not outside.
func outside(box [4]float64, reach float64, m matrix, bounds image.Rectangle) bool {
	if bounds.Empty() {
		return true
	}
	x0, y0, x1, y1 := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	for _, c := range [4][2]float64{{box[0] - reach, box[1] - reach}, {box[2] + reach, box[1] - reach}, {box[0] - reach, box[3] + reach}, {box[2] + reach, box[3] + reach}} {
		px, py := m.apply(c[0], c[1])
		if math.IsNaN(px) || math.IsNaN(py) {
			return false
		}
		x0, y0, x1, y1 = math.Min(x0, px), math.Min(y0, py), math.Max(x1, px), math.Max(y1, py)
	}
	return x1 < float64(bounds.Min.X)-1 || x0 > float64(bounds.Max.X)+1 || y1 < float64(bounds.Min.Y)-1 || y0 > float64(bounds.Max.Y)+1
}

// boxUnder returns the bounding box of a box under m.
func boxUnder(box [4]float64, m matrix) [4]float64 {
	x0, y0, x1, y1 := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	for _, c := range [4][2]float64{{box[0], box[1]}, {box[2], box[1]}, {box[0], box[3]}, {box[2], box[3]}} {
		px, py := m.apply(c[0], c[1])
		if math.IsNaN(px) || math.IsNaN(py) {
			return [4]float64{px, py, px, py}
		}
		x0, y0, x1, y1 = math.Min(x0, px), math.Min(y0, py), math.Max(x1, px), math.Max(y1, py)
	}
	return [4]float64{x0, y0, x1, y1}
}

// hidden reports whether a path drawn under m lies outside bounds, with
// what reaches up to reach around it (in its units): the width of a line.
func hidden(p *path, reach float64, m matrix, bounds image.Rectangle) bool {
	if plain {
		return false
	}
	box, ok := p.extent()
	return ok && outside(box, reach, m, bounds)
}

// segmentsFor returns how many straight pieces keep a curve whose control
// polygon bends by dd within tol.
func segmentsFor(dd, tol float64) int {
	n := int(math.Ceil(math.Sqrt(dd / (8 * tol))))
	switch {
	case n < 1 || math.IsNaN(dd):
		return 1
	case n > 256:
		return 256
	}
	return n
}

func flattenQuad(out []point, x0, y0, x1, y1, x2, y2, tol float64) []point {
	dd := math.Hypot(x0-2*x1+x2, y0-2*y1+y2) * 2
	n := segmentsFor(dd, tol)
	for i := 1; i <= n; i++ {
		t := float64(i) / float64(n)
		u := 1 - t
		out = append(out, point{u*u*x0 + 2*u*t*x1 + t*t*x2, u*u*y0 + 2*u*t*y1 + t*t*y2})
	}
	return out
}

func flattenCubic(out []point, x0, y0, x1, y1, x2, y2, x3, y3, tol float64) []point {
	dd := 6 * math.Max(math.Hypot(x0-2*x1+x2, y0-2*y1+y2), math.Hypot(x1-2*x2+x3, y1-2*y2+y3))
	n := segmentsFor(dd, tol)
	for i := 1; i <= n; i++ {
		t := float64(i) / float64(n)
		u := 1 - t
		a, b, c, d := u*u*u, 3*u*u*t, 3*u*t*t, t*t*t
		out = append(out, point{a*x0 + b*x1 + c*x2 + d*x3, a*y0 + b*y1 + c*y2 + d*y3})
	}
	return out
}

// isRect reports whether the path is one axis-aligned rectangle (as rect
// makes it), and returns its corners.
func (p *path) isRect() (x0, y0, x1, y1 float64, ok bool) {
	var rs *subpath
	for i := range p.subs {
		s := &p.subs[i]
		if len(s.segs) == 0 {
			continue
		}
		if rs != nil {
			return 0, 0, 0, 0, false
		}
		rs = s
	}
	if rs == nil || len(rs.segs) != 3 && !(len(rs.segs) == 4 && rs.segs[3].p[0] == rs.start.x && rs.segs[3].p[1] == rs.start.y) {
		return 0, 0, 0, 0, false
	}
	pts := []point{rs.start}
	for _, g := range rs.segs[:3] {
		if g.kind != segLine {
			return 0, 0, 0, 0, false
		}
		pts = append(pts, point{g.p[0], g.p[1]})
	}
	if rs.segs[len(rs.segs)-1].kind != segLine {
		return 0, 0, 0, 0, false
	}
	a, b, c, d := pts[0], pts[1], pts[2], pts[3]
	horizontalFirst := a.y == b.y && b.x == c.x && c.y == d.y && d.x == a.x
	verticalFirst := a.x == b.x && b.y == c.y && c.x == d.x && d.y == a.y
	if !horizontalFirst && !verticalFirst {
		return 0, 0, 0, 0, false
	}
	return math.Min(a.x, c.x), math.Min(a.y, c.y), math.Max(a.x, c.x), math.Max(a.y, c.y), true
}
