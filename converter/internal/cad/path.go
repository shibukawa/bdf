package cad

import (
	"math"

	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// Point is a point or a vector in drawing coordinates.
type Point struct{ X, Y float64 }

// Add returns p + q.
func (p Point) Add(q Point) Point { return Point{p.X + q.X, p.Y + q.Y} }

// Sub returns p − q.
func (p Point) Sub(q Point) Point { return Point{p.X - q.X, p.Y - q.Y} }

// Mul returns p scaled by s.
func (p Point) Mul(s float64) Point { return Point{p.X * s, p.Y * s} }

// Len returns the length of p.
func (p Point) Len() float64 { return math.Hypot(p.X, p.Y) }

// Apply transforms a point.
func Apply(m canvas.Matrix, p Point) Point {
	x, y := m.Apply(p.X, p.Y)
	return Point{x, y}
}

// Path verbs.
const (
	moveTo byte = iota
	lineTo
	cubicTo
	closePath
)

// Path is a path in drawing coordinates. Curves are cubic Bézier curves:
// arcs and ellipses are converted when they are added, so that a path
// keeps its shape under any affine transform.
type Path struct {
	verbs []byte
	pts   []Point
	cur   Point
	start Point
	open  bool // a subpath has been started
}

// Empty reports whether the path has no segments.
func (p *Path) Empty() bool { return p == nil || len(p.verbs) == 0 }

// Current returns the current point.
func (p *Path) Current() Point { return p.cur }

// MoveTo starts a subpath.
func (p *Path) MoveTo(x, y float64) *Path {
	p.verbs = append(p.verbs, moveTo)
	p.pts = append(p.pts, Point{x, y})
	p.cur, p.start, p.open = Point{x, y}, Point{x, y}, true
	return p
}

// LineTo adds a line segment (a move when no subpath has been started).
func (p *Path) LineTo(x, y float64) *Path {
	if !p.open {
		return p.MoveTo(x, y)
	}
	p.verbs = append(p.verbs, lineTo)
	p.pts = append(p.pts, Point{x, y})
	p.cur = Point{x, y}
	return p
}

// CubicTo adds a cubic Bézier curve.
func (p *Path) CubicTo(c1, c2, to Point) *Path {
	if !p.open {
		p.MoveTo(c1.X, c1.Y)
	}
	p.verbs = append(p.verbs, cubicTo)
	p.pts = append(p.pts, c1, c2, to)
	p.cur = to
	return p
}

// Close closes the subpath.
func (p *Path) Close() *Path {
	if p.open {
		p.verbs = append(p.verbs, closePath)
		p.cur = p.start
		p.open = false
	}
	return p
}

// Polyline adds the points as line segments, starting a subpath at the
// first; closed closes it.
func (p *Path) Polyline(pts []Point, closed bool) *Path {
	for i, q := range pts {
		if i == 0 {
			p.MoveTo(q.X, q.Y)
		} else {
			p.LineTo(q.X, q.Y)
		}
	}
	if closed && len(pts) > 0 {
		p.Close()
	}
	return p
}

// Append adds the subpaths of q.
func (p *Path) Append(q *Path) *Path {
	if q == nil {
		return p
	}
	p.verbs = append(p.verbs, q.verbs...)
	p.pts = append(p.pts, q.pts...)
	p.cur, p.start, p.open = q.cur, q.start, q.open
	return p
}

// Open reports whether a subpath has been started and not closed.
func (p *Path) Open() bool { return p.open }

// First returns the first point of the path.
func (p *Path) First() Point {
	if len(p.pts) == 0 {
		return Point{}
	}
	return p.pts[0]
}

// Continue adds the segments of q to the open subpath: q's first move
// becomes a line from the current point (dropped when they meet).
func (p *Path) Continue(q *Path) *Path {
	if q.Empty() {
		return p
	}
	if !p.open {
		return p.Append(q)
	}
	i := 0
	q.walk(func(v byte, pts []Point) {
		if i == 0 && v == moveTo {
			if pts[0] != p.cur {
				p.LineTo(pts[0].X, pts[0].Y)
			}
		} else {
			switch v {
			case moveTo:
				p.MoveTo(pts[0].X, pts[0].Y)
			case lineTo:
				p.LineTo(pts[0].X, pts[0].Y)
			case cubicTo:
				p.CubicTo(pts[0], pts[1], pts[2])
			case closePath:
				p.Close()
			}
		}
		i++
	})
	return p
}

// Transform returns the path transformed by m.
func (p *Path) Transform(m canvas.Matrix) *Path {
	q := &Path{verbs: p.verbs, pts: make([]Point, len(p.pts)), open: p.open}
	for i, pt := range p.pts {
		q.pts[i] = Apply(m, pt)
	}
	q.cur, q.start = Apply(m, p.cur), Apply(m, p.start)
	return q
}

// Walk calls f for each segment: the verb's points (1 for a move or a
// line, 3 for a curve, none for a close).
func (p *Path) walk(f func(verb byte, pts []Point)) {
	i := 0
	for _, v := range p.verbs {
		n := 0
		switch v {
		case moveTo, lineTo:
			n = 1
		case cubicTo:
			n = 3
		}
		f(v, p.pts[i:i+n])
		i += n
	}
}

// Bounds returns the bounding box of the path's points (control points
// included).
func (p *Path) Bounds() Rect {
	var r Rect
	for _, pt := range p.pts {
		r = r.Add(pt)
	}
	return r
}

// Arc adds a circular arc around c from angle a0 to a1 (radians,
// counter-clockwise from the x axis in a y-up space when a1 > a0): a line
// to its start when a subpath is open, a move otherwise.
func (p *Path) Arc(c Point, r, a0, a1 float64) *Path {
	return p.EllipseArc(c, Point{r, 0}, 1, a0, a1)
}

// Circle adds a full circle as a closed subpath.
func (p *Path) Circle(c Point, r float64) *Path {
	p.open = false
	p.Arc(c, r, 0, 2*math.Pi)
	return p.Close()
}

// EllipseArc adds the arc of an ellipse centred at c with the major axis
// vector major and the minor to major ratio from parameter t0 to t1
// (radians; the minor axis is the major one turned by 90° counter-
// clockwise): a line to its start when a subpath is open, a move
// otherwise.
func (p *Path) EllipseArc(c, major Point, ratio, t0, t1 float64) *Path {
	return p.EllipseArcAxes(c, major, Point{-major.Y * ratio, major.X * ratio}, t0, t1)
}

// EllipseArcAxes adds the arc c + major·cos t + minor·sin t from t0 to t1:
// a line to its start when a subpath is open, a move otherwise.
func (p *Path) EllipseArcAxes(c, major, minor Point, t0, t1 float64) *Path {
	at := func(t float64) Point {
		s, co := math.Sincos(t)
		return Point{c.X + major.X*co + minor.X*s, c.Y + major.Y*co + minor.Y*s}
	}
	// derivative
	d := func(t float64) Point {
		s, co := math.Sincos(t)
		return Point{-major.X*s + minor.X*co, -major.Y*s + minor.Y*co}
	}
	start := at(t0)
	if p.open {
		p.LineTo(start.X, start.Y)
	} else {
		p.MoveTo(start.X, start.Y)
	}
	sweep := t1 - t0
	n := int(math.Ceil(math.Abs(sweep) / (math.Pi / 2)))
	if n < 1 {
		n = 1
	}
	if n > 64 {
		n = 64
	}
	h := sweep / float64(n)
	k := 4.0 / 3 * math.Tan(h/4)
	for i := 0; i < n; i++ {
		ta, tb := t0+h*float64(i), t0+h*float64(i+1)
		a, b := at(ta), at(tb)
		da, db := d(ta), d(tb)
		p.CubicTo(a.Add(da.Mul(k)), b.Sub(db.Mul(k)), b)
	}
	return p
}

// Ellipse adds a full ellipse as a closed subpath.
func (p *Path) Ellipse(c, major Point, ratio float64) *Path {
	p.open = false
	p.EllipseArc(c, major, ratio, 0, 2*math.Pi)
	return p.Close()
}

// BulgeTo adds a segment from the current point to q that is an arc when
// bulge (the tangent of a quarter of the included angle, positive counter-
// clockwise) is not zero: the segments of DXF polylines.
func (p *Path) BulgeTo(q Point, bulge float64) *Path {
	a := p.cur
	if bulge == 0 || a == q {
		return p.LineTo(q.X, q.Y)
	}
	c, r, a0, a1 := BulgeArc(a, q, bulge)
	if r == 0 || math.IsInf(r, 0) || math.IsNaN(r) {
		return p.LineTo(q.X, q.Y)
	}
	return p.Arc(c, r, a0, a1)
}

// BulgeArc returns the circle of a polyline segment from a to b with a
// bulge, and the start and end angles (a1 < a0 for a clockwise arc).
func BulgeArc(a, b Point, bulge float64) (c Point, r, a0, a1 float64) {
	theta := 4 * math.Atan(bulge) // included angle, signed
	chord := b.Sub(a)
	l := chord.Len()
	r = l / (2 * math.Sin(theta/2))
	// distance from the chord midpoint to the centre, towards the left of
	// the chord for a counter-clockwise arc
	h := r * math.Cos(theta/2)
	mid := a.Add(chord.Mul(0.5))
	n := Point{-chord.Y / l, chord.X / l}
	c = mid.Add(n.Mul(h))
	r = math.Abs(r)
	a0 = math.Atan2(a.Y-c.Y, a.X-c.X)
	a1 = a0 + theta
	return
}

// Rect is an axis-aligned rectangle; the zero value is empty.
type Rect struct {
	Min, Max Point
	ok       bool
}

// Valid reports whether the rectangle holds any point.
func (r Rect) Valid() bool { return r.ok }

// Add returns r grown to hold p.
func (r Rect) Add(p Point) Rect {
	if math.IsNaN(p.X) || math.IsNaN(p.Y) || math.IsInf(p.X, 0) || math.IsInf(p.Y, 0) {
		return r
	}
	if !r.ok {
		return Rect{p, p, true}
	}
	r.Min.X, r.Min.Y = min(r.Min.X, p.X), min(r.Min.Y, p.Y)
	r.Max.X, r.Max.Y = max(r.Max.X, p.X), max(r.Max.Y, p.Y)
	return r
}

// Union returns the rectangle holding r and s.
func (r Rect) Union(s Rect) Rect {
	if !s.ok {
		return r
	}
	return r.Add(s.Min).Add(s.Max)
}

// Transform returns the bounding box of r transformed by m.
func (r Rect) Transform(m canvas.Matrix) Rect {
	if !r.ok {
		return r
	}
	var out Rect
	for _, p := range []Point{r.Min, {r.Max.X, r.Min.Y}, r.Max, {r.Min.X, r.Max.Y}} {
		out = out.Add(Apply(m, p))
	}
	return out
}

// W returns the width.
func (r Rect) W() float64 { return r.Max.X - r.Min.X }

// H returns the height.
func (r Rect) H() float64 { return r.Max.Y - r.Min.Y }

// Scale returns the mean scale factor of a transform (the square root of
// its determinant's magnitude): how lengths change on average.
func Scale(m canvas.Matrix) float64 {
	return math.Sqrt(math.Abs(m[0]*m[3] - m[1]*m[2]))
}
