package gerber

import (
	"math"
	"slices"
)

// vec is a point in the plane of a file, in mm with y up.
type vec struct{ X, Y float64 }

func (a vec) add(b vec) vec       { return vec{a.X + b.X, a.Y + b.Y} }
func (a vec) sub(b vec) vec       { return vec{a.X - b.X, a.Y - b.Y} }
func (a vec) mul(s float64) vec   { return vec{a.X * s, a.Y * s} }
func (a vec) len() float64        { return math.Hypot(a.X, a.Y) }
func (a vec) cross(b vec) float64 { return a.X*b.Y - a.Y*b.X }

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

// box is a bounding box; the zero value is empty.
type box struct {
	Min, Max vec
	ok       bool
}

func (b box) add(p vec) box {
	if !finite(p.X) || !finite(p.Y) {
		return b
	}
	if !b.ok {
		return box{p, p, true}
	}
	b.Min = vec{min(b.Min.X, p.X), min(b.Min.Y, p.Y)}
	b.Max = vec{max(b.Max.X, p.X), max(b.Max.Y, p.Y)}
	return b
}

func (b box) union(c box) box {
	if !c.ok {
		return b
	}
	if !b.ok {
		return c
	}
	return b.add(c.Min).add(c.Max)
}

// grow widens the box by d on every side.
func (b box) grow(d float64) box {
	if !b.ok {
		return b
	}
	return box{b.Min.sub(vec{d, d}), b.Max.add(vec{d, d}), true}
}

func (b box) offset(d vec) box {
	if !b.ok {
		return b
	}
	return box{b.Min.add(d), b.Max.add(d), true}
}

func (b box) w() float64 { return b.Max.X - b.Min.X }
func (b box) h() float64 { return b.Max.Y - b.Min.Y }

// seg is a segment of a contour or of a stroked chain: a line or a
// circular arc to a point, from the end of the segment before it.
type seg struct {
	to vec
	// arc: the arc turns about c by sweep radians (positive:
	// counterclockwise), starting at the point before
	arc   bool
	c     vec
	sweep float64
}

// contour is a chain of segments from start; the contours of filled
// shapes are closed, outer ones counterclockwise and holes clockwise.
type contour struct {
	start vec
	segs  []seg
}

// end returns the point the contour ends at.
func (c *contour) end() vec {
	if n := len(c.segs); n > 0 {
		return c.segs[n-1].to
	}
	return c.start
}

func (c *contour) lineTo(p vec) { c.segs = append(c.segs, seg{to: p}) }

// arcTo adds an arc about center by sweep radians; the arc ends where the
// circle through the current point reaches, which is to when to lies on
// that circle.
func (c *contour) arcTo(to, center vec, sweep float64) {
	c.segs = append(c.segs, seg{to: to, arc: true, c: center, sweep: sweep})
}

// circle returns a closed contour of a circle, counterclockwise.
func circle(center vec, r float64) contour {
	p := center.add(vec{r, 0})
	return contour{start: p, segs: []seg{{to: p, arc: true, c: center, sweep: 2 * math.Pi}}}
}

// polygon returns a closed contour through pts.
func polygon(pts ...vec) contour {
	c := contour{start: pts[0]}
	for _, p := range pts[1:] {
		c.lineTo(p)
	}
	if pts[len(pts)-1] != pts[0] {
		c.lineTo(pts[0])
	}
	return c
}

// rect returns a closed counterclockwise rectangle centred on c.
func rect(c vec, w, h float64) contour {
	x0, y0, x1, y1 := c.X-w/2, c.Y-h/2, c.X+w/2, c.Y+h/2
	return polygon(vec{x0, y0}, vec{x1, y0}, vec{x1, y1}, vec{x0, y1})
}

// flatten appends points along the contour (arcs as polylines with a point
// at least every step radians) to pts.
func (c *contour) flatten(pts []vec, step float64) []vec {
	pts = append(pts, c.start)
	cur := c.start
	for _, s := range c.segs {
		if s.arc {
			r := cur.sub(s.c).len()
			a0 := math.Atan2(cur.Y-s.c.Y, cur.X-s.c.X)
			n := int(math.Ceil(math.Abs(s.sweep) / step))
			n = max(1, min(n, 4096))
			for i := 1; i < n; i++ {
				a := a0 + s.sweep*float64(i)/float64(n)
				pts = append(pts, vec{s.c.X + r*math.Cos(a), s.c.Y + r*math.Sin(a)})
			}
		}
		pts = append(pts, s.to)
		cur = s.to
	}
	return pts
}

// area returns the signed area of the contour, closed by a line from its
// end to its start (positive: counterclockwise).
func (c *contour) area() float64 {
	a := 0.0
	cur := c.start
	for _, s := range c.segs {
		if s.arc {
			// the chord, and the circular segment between the chord and the arc
			r := cur.sub(s.c).len()
			a += cur.cross(s.to) / 2
			sw := s.sweep
			a += r * r * (sw - math.Sin(sw)) / 2
			// the arc ends on its circle; a to off the circle adds its chord
		} else {
			a += cur.cross(s.to) / 2
		}
		cur = s.to
	}
	a += cur.cross(c.start) / 2
	return a
}

// reverse returns the contour traversed the other way.
func (c *contour) reverse() contour {
	n := len(c.segs)
	out := contour{start: c.end(), segs: make([]seg, n)}
	for i := range n {
		s := c.segs[n-1-i]
		from := c.start
		if n-2-i >= 0 {
			from = c.segs[n-2-i].to
		}
		out.segs[i] = seg{to: from, arc: s.arc, c: s.c, sweep: -s.sweep}
	}
	return out
}

// orient returns the contour counterclockwise when ccw is true, clockwise
// otherwise.
func (c *contour) orient(ccw bool) contour {
	if (c.area() >= 0) != ccw {
		return c.reverse()
	}
	return *c
}

// bounds returns the bounding box of the contour.
func (c *contour) bounds() box {
	b := box{}.add(c.start)
	cur := c.start
	for _, s := range c.segs {
		b = b.add(s.to)
		if s.arc {
			b = b.union(arcBounds(cur, s))
		}
		cur = s.to
	}
	return b
}

// arcBounds returns the bounding box of an arc from p.
func arcBounds(p vec, s seg) box {
	r := p.sub(s.c).len()
	a0 := math.Atan2(p.Y-s.c.Y, p.X-s.c.X)
	b := box{}.add(p)
	lo, hi := a0, a0+s.sweep
	if hi < lo {
		lo, hi = hi, lo
	}
	if hi-lo >= 2*math.Pi {
		return b.add(s.c.sub(vec{r, r})).add(s.c.add(vec{r, r}))
	}
	// the axis points within the sweep
	for k := math.Ceil(lo / (math.Pi / 2)); k*math.Pi/2 <= hi; k++ {
		a := k * math.Pi / 2
		b = b.add(vec{s.c.X + r*math.Cos(a), s.c.Y + r*math.Sin(a)})
	}
	return b.add(vec{s.c.X + r*math.Cos(a0+s.sweep), s.c.Y + r*math.Sin(a0+s.sweep)})
}

// affine is a transform x' = a x + c y + e, y' = b x + d y + f.
type affine [6]float64

var identity = affine{1, 0, 0, 1, 0, 0}

func (m affine) apply(p vec) vec {
	return vec{m[0]*p.X + m[2]*p.Y + m[4], m[1]*p.X + m[3]*p.Y + m[5]}
}

// mul returns m after n (m · n).
func (m affine) mul(n affine) affine {
	return affine{
		m[0]*n[0] + m[2]*n[1], m[1]*n[0] + m[3]*n[1],
		m[0]*n[2] + m[2]*n[3], m[1]*n[2] + m[3]*n[3],
		m[0]*n[4] + m[2]*n[5] + m[4], m[1]*n[4] + m[3]*n[5] + m[5],
	}
}

func (m affine) det() float64 { return m[0]*m[3] - m[1]*m[2] }

// scale returns how much m scales lengths (the square root of its
// determinant's magnitude).
func (m affine) scale() float64 { return math.Sqrt(math.Abs(m.det())) }

func (m affine) linear() bool { return m[4] == 0 && m[5] == 0 }

func translate(d vec) affine { return affine{1, 0, 0, 1, d.X, d.Y} }

func rotate(deg float64) affine {
	switch math.Mod(deg, 360) {
	case 0:
		return identity
	case 90, -270:
		return affine{0, 1, -1, 0, 0, 0}
	case 180, -180:
		return affine{-1, 0, 0, -1, 0, 0}
	case 270, -90:
		return affine{0, -1, 1, 0, 0, 0}
	}
	s, c := math.Sincos(deg * math.Pi / 180)
	return affine{c, s, -s, c, 0, 0}
}

// transform returns the contour transformed by m (which keeps circles
// circles: rotations, reflections, uniform scales and translations).
func (c *contour) transform(m affine) contour {
	flip := m.det() < 0
	out := contour{start: m.apply(c.start), segs: make([]seg, len(c.segs))}
	for i, s := range c.segs {
		t := seg{to: m.apply(s.to), arc: s.arc}
		if s.arc {
			t.c = m.apply(s.c)
			t.sweep = s.sweep
			if flip {
				t.sweep = -s.sweep
			}
		}
		out.segs[i] = t
	}
	return out
}

// translated returns the contour moved by d.
func (c *contour) translated(d vec) contour {
	out := contour{start: c.start.add(d), segs: make([]seg, len(c.segs))}
	for i, s := range c.segs {
		s.to = s.to.add(d)
		if s.arc {
			s.c = s.c.add(d)
		}
		out.segs[i] = s
	}
	return out
}

// hull returns the convex hull of pts, counterclockwise.
func hull(pts []vec) []vec {
	pts = slices.Clone(pts)
	slices.SortFunc(pts, func(a, b vec) int {
		if a.X != b.X {
			if a.X < b.X {
				return -1
			}
			return 1
		}
		switch {
		case a.Y < b.Y:
			return -1
		case a.Y > b.Y:
			return 1
		}
		return 0
	})
	pts = slices.Compact(pts)
	if len(pts) < 3 {
		return pts
	}
	var h []vec
	for pass := 0; pass < 2; pass++ {
		start := len(h)
		for _, p := range pts {
			for len(h) >= start+2 && h[len(h)-1].sub(h[len(h)-2]).cross(p.sub(h[len(h)-2])) <= 0 {
				h = h[:len(h)-1]
			}
			h = append(h, p)
		}
		h = h[:len(h)-1]
		slices.Reverse(pts)
	}
	return h
}
