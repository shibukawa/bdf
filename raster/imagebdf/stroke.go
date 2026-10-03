package imagebdf

import (
	"image"
	"math"

	"github.com/shibukawa/bdf"
)

// lineStyle is the stroke state of a context.
type lineStyle struct {
	width   float64
	cap     byte
	join    byte
	miter   float64
	dash    []float64 // even length, or nil
	dashOff float64
}

// miterReach is how far a miter join reaches from its corner, in half
// line widths: the miter limit, or 1 for the other joins.
func (ls *lineStyle) miterReach() float64 {
	if ls.join != bdf.JoinMiter {
		return 1
	}
	if !(ls.miter > 0) {
		return 10
	}
	return ls.miter
}

// setDash applies setLineDash: an odd list repeats, a list with a negative
// or non-finite value is ignored, and an all-zero list draws solid lines.
func (ls *lineStyle) setDash(segs []float32, off float32) {
	var d []float64
	sum := 0.0
	for _, s := range segs {
		v := float64(s)
		if v < 0 || math.IsNaN(v) || math.IsInf(v, 0) {
			return
		}
		d = append(d, v)
		sum += v
	}
	if len(d)%2 == 1 {
		d = append(d, d...)
	}
	if sum == 0 {
		d = nil
	}
	ls.dash = d
	if o := float64(off); !math.IsNaN(o) && !math.IsInf(o, 0) {
		ls.dashOff = o
	}
}

// Limits of the outline of the strokes of one drawing instruction. A
// document decides the width of a line, its dashes and its joins, so a line
// of two points could otherwise ask for any number of polygons.
const (
	// strokePoints is the most points of an outline, besides
	// strokePointsEach for each point of the lines stroked: round joins
	// and caps of lines far wider than the page take up to 256 each.
	strokePoints     = 1 << 20
	strokePointsEach = 16
	// strokeDashes is the most dashes the lines are cut into.
	strokeDashes = 1 << 20
)

// stroker makes the outline of strokes: polygons that all wind the same
// way, for a nonzero fill.
type stroker struct {
	ls      *lineStyle
	hw, tol float64
	out     []polyline

	// m maps the outline to the device and bounds is what shows of it:
	// with cull, polygons that lie outside are left out, which changes
	// nothing of what shows
	m      matrix
	bounds image.Rectangle
	cull   bool

	points, dashes int  // what the outline may still take
	cut            bool // a limit was passed: the rest of the outline is left out
}

// newStroker starts the outline of strokes with a line style. devScale is
// the device pixels per user unit, for the smoothness of round parts.
func newStroker(ls *lineStyle, devScale float64) *stroker {
	return &stroker{ls: ls, hw: ls.width / 2, tol: 0.2 / math.Max(devScale, 1e-9), points: strokePoints, dashes: strokeDashes}
}

// within makes the stroker leave out what lies outside bounds under m.
func (s *stroker) within(m matrix, bounds image.Rectangle) *stroker {
	s.m, s.bounds, s.cull = m, bounds, !plain
	return s
}

// hidden reports whether a box lies outside of what shows.
func (s *stroker) hidden(x0, y0, x1, y1 float64) bool {
	return s.cull && outside([4]float64{x0, y0, x1, y1}, 0, s.m, s.bounds)
}

// take counts n points of the outline, and reports whether it may have them.
func (s *stroker) take(n int) bool {
	if s.cut || n > s.points {
		s.cut = true
		return false
	}
	s.points -= n
	return true
}

// put adds a polygon to the outline.
func (s *stroker) put(pts []point) {
	if !s.take(len(pts)) {
		return
	}
	if signedArea(pts) < 0 {
		reversePoints(pts)
	}
	s.out = append(s.out, polyline{pts: pts, closed: true})
}

// add adds a polygon, unless it lies outside of what shows.
func (s *stroker) add(pts ...point) {
	if s.cull {
		x0, y0, x1, y1 := pts[0].x, pts[0].y, pts[0].x, pts[0].y
		for _, p := range pts[1:] {
			x0, y0, x1, y1 = math.Min(x0, p.x), math.Min(y0, p.y), math.Max(x1, p.x), math.Max(y1, p.y)
		}
		if s.hidden(x0, y0, x1, y1) {
			return
		}
	}
	s.put(pts)
}

// quad adds the four corners of a segment or a square cap.
func (s *stroker) quad(a, b, c, d point) {
	if s.cull && s.hidden(math.Min(math.Min(a.x, b.x), math.Min(c.x, d.x)), math.Min(math.Min(a.y, b.y), math.Min(c.y, d.y)),
		math.Max(math.Max(a.x, b.x), math.Max(c.x, d.x)), math.Max(math.Max(a.y, b.y), math.Max(c.y, d.y))) {
		return
	}
	s.put([]point{a, b, c, d})
}

// round adds a circle of the width of the line: a round cap or join.
func (s *stroker) round(c point) {
	if s.cut || s.hidden(c.x-s.hw, c.y-s.hw, c.x+s.hw, c.y+s.hw) {
		return
	}
	s.put(circle(c, s.hw, s.tol))
}

// square adds a square of the width of the line: the dot of a square cap.
func (s *stroker) square(c point) {
	hw := s.hw
	s.quad(point{c.x - hw, c.y - hw}, point{c.x + hw, c.y - hw}, point{c.x + hw, c.y + hw}, point{c.x - hw, c.y + hw})
}

// strokePolys returns the outline of a stroke of the polylines (user
// space) as polygons that all wind the same way, for a nonzero fill.
// devScale is the device pixels per user unit, for the smoothness of round
// parts.
func strokePolys(lines []polyline, ls *lineStyle, devScale float64) []polyline {
	s := newStroker(ls, devScale)
	s.stroke(lines)
	return s.out
}

// stroke adds the outline of the polylines (user space).
func (s *stroker) stroke(lines []polyline) {
	ls := s.ls
	if !(s.hw > 0) {
		return
	}
	for _, pl := range lines {
		s.points += strokePointsEach * len(pl.pts)
	}
	for _, pl := range lines {
		if s.cut {
			return
		}
		pts := dedupe(pl.pts)
		closed := pl.closed
		if len(pts) == 1 {
			// a zero-length sub-path: a dot under round and square caps
			if len(pl.pts) >= 2 {
				switch ls.cap {
				case bdf.CapRound:
					s.round(pts[0])
				case bdf.CapSquare:
					s.square(pts[0])
				}
			}
			continue
		}
		if len(pts) < 2 {
			continue
		}
		if ls.dash == nil {
			s.piece(pts, closed)
			continue
		}
		if closed {
			pts = append(pts, pts[0])
		}
		s.dash(pts)
	}
}

// piece outlines one polyline: a quadrilateral per segment, a join at each
// corner and caps at the ends of an open one.
func (s *stroker) piece(pts []point, closed bool) {
	ls, hw := s.ls, s.hw
	pts = dedupe(pts)
	if closed && len(pts) > 2 && pts[0] == pts[len(pts)-1] {
		pts = pts[:len(pts)-1]
	}
	n := len(pts)
	if n < 2 {
		if n == 1 && ls.cap != bdf.CapButt {
			if ls.cap == bdf.CapRound {
				s.round(pts[0])
			} else {
				s.square(pts[0])
			}
		}
		return
	}
	segs := n - 1
	if closed {
		segs = n
	}
	for i := 0; i < segs && !s.cut; i++ {
		a, b := pts[i], pts[(i+1)%n]
		nx, ny := normal(a, b, hw)
		s.quad(point{a.x + nx, a.y + ny}, point{b.x + nx, b.y + ny}, point{b.x - nx, b.y - ny}, point{a.x - nx, a.y - ny})
	}
	for i := 1; i < n-1 && !s.cut; i++ {
		s.join(pts[i-1], pts[i], pts[i+1])
	}
	if closed {
		s.join(pts[n-2], pts[n-1], pts[0])
		s.join(pts[n-1], pts[0], pts[1])
		return
	}
	cap := func(end, from point) {
		switch ls.cap {
		case bdf.CapRound:
			s.round(end)
		case bdf.CapSquare:
			dx, dy := end.x-from.x, end.y-from.y
			l := math.Hypot(dx, dy)
			ux, uy := dx/l*hw, dy/l*hw
			nx, ny := -uy, ux
			s.quad(point{end.x + nx, end.y + ny}, point{end.x + nx + ux, end.y + ny + uy}, point{end.x - nx + ux, end.y - ny + uy}, point{end.x - nx, end.y - ny})
		}
	}
	cap(pts[0], pts[1])
	cap(pts[n-1], pts[n-2])
}

// join fills the corner at v between the segments prev→v and v→next.
func (s *stroker) join(prev, v, next point) {
	ls, hw, tol := s.ls, s.hw, s.tol
	ax, ay := v.x-prev.x, v.y-prev.y
	bx, by := next.x-v.x, next.y-v.y
	la, lb := math.Hypot(ax, ay), math.Hypot(bx, by)
	if la == 0 || lb == 0 {
		return
	}
	cross := ax*by - ay*bx
	dot := ax*bx + ay*by
	turn := math.Atan2(math.Abs(cross), dot) // 0: straight on
	if turn < 1e-9 {
		return
	}
	// the outer side: offsets of both segments away from the turn
	n1x, n1y := normal(prev, v, hw)
	n2x, n2y := normal(v, next, hw)
	if cross > 0 {
		n1x, n1y, n2x, n2y = -n1x, -n1y, -n2x, -n2y
	}
	p1 := point{v.x + n1x, v.y + n1y}
	p2 := point{v.x + n2x, v.y + n2y}
	switch ls.join {
	case bdf.JoinRound:
		if turn*hw < 4*tol {
			s.add(v, p1, p2)
			return
		}
		s.round(v)
	case bdf.JoinMiter:
		// miter length over line width is 1/sin(θ/2), θ the angle between the segments
		theta := math.Pi - turn
		limit := ls.miter
		if !(limit > 0) {
			limit = 10
		}
		if sn := math.Sin(theta / 2); sn > 0 && 1/sn <= limit {
			d := hw / math.Cos(turn/2)
			mx, my := (n1x+n2x)/2, (n1y+n2y)/2
			ml := math.Hypot(mx, my)
			if ml > 0 {
				s.add(v, p1, point{v.x + mx/ml*d, v.y + my/ml*d}, p2)
				return
			}
		}
		s.add(v, p1, p2)
	default:
		s.add(v, p1, p2)
	}
}

// normal returns the left normal of a→b scaled to hw.
func normal(a, b point, hw float64) (float64, float64) {
	dx, dy := b.x-a.x, b.y-a.y
	l := math.Hypot(dx, dy)
	if l == 0 {
		return 0, 0
	}
	return -dy / l * hw, dx / l * hw
}

// circle returns a polygon approximating a circle within tol.
func circle(c point, r, tol float64) []point {
	n := 8
	if r > tol {
		step := 2 * math.Acos(1-tol/r)
		// bounded as a float64: a circle far larger than tol gives a
		// number no int holds
		n = int(math.Min(math.Ceil(2*math.Pi/step), 256))
	}
	n = max(8, min(n, 256))
	pts := make([]point, n)
	for i := range pts {
		s, co := math.Sincos(2 * math.Pi * float64(i) / float64(n))
		pts[i] = point{c.x + r*co, c.y + r*s}
	}
	return pts
}

func signedArea(pts []point) float64 {
	a := 0.0
	for i := range pts {
		p, q := pts[i], pts[(i+1)%len(pts)]
		a += p.x*q.y - q.x*p.y
	}
	return a
}

func reversePoints(pts []point) {
	for i, j := 0, len(pts)-1; i < j; i, j = i+1, j-1 {
		pts[i], pts[j] = pts[j], pts[i]
	}
}

// dedupe drops consecutive repeated points.
func dedupe(pts []point) []point {
	if len(pts) < 2 {
		return pts
	}
	out := make([]point, 0, len(pts))
	out = append(out, pts[0])
	for _, p := range pts[1:] {
		if p != out[len(out)-1] {
			out = append(out, p)
		}
	}
	return out
}

// dash strokes the dashes of a polyline: the pieces a pattern (on, off, …)
// cuts it into, starting dashOff into the pattern.
func (s *stroker) dash(pts []point) {
	pattern, offset := s.ls.dash, s.ls.dashOff
	total := 0.0
	for _, v := range pattern {
		total += v
	}
	if total <= 0 {
		s.piece(pts, false)
		return
	}
	offset = math.Mod(offset, total)
	if offset < 0 {
		offset += total
	}
	idx := 0
	left := pattern[0]
	for offset > 0 {
		if offset < left {
			left -= offset
			break
		}
		offset -= left
		idx = (idx + 1) % len(pattern)
		left = pattern[idx]
	}
	var cur []point
	on := idx%2 == 0
	if on {
		cur = []point{pts[0]}
	}
	// hidden reports whether the dashes of a stretch of a segment lie
	// outside of what shows, with their joins and caps
	reach := s.hw * math.Max(math.Sqrt2, s.ls.miterReach())
	hidden := func(a, b point) bool {
		return s.hidden(math.Min(a.x, b.x)-reach, math.Min(a.y, b.y)-reach, math.Max(a.x, b.x)+reach, math.Max(a.y, b.y)+reach)
	}
	for i := 1; i < len(pts); i++ {
		a, b := pts[i-1], pts[i]
		segLen := math.Hypot(b.x-a.x, b.y-a.y)
		pos := 0.0
		skip := hidden(a, b)
		for segLen-pos > left {
			if s.dashes--; s.dashes < 0 || s.cut {
				s.cut = true // a pattern far too fine for the line
				return
			}
			pos += left
			t := pos / segLen
			p := point{a.x + (b.x-a.x)*t, a.y + (b.y-a.y)*t}
			if on {
				if !skip || len(cur) > 1 {
					s.piece(append(cur, p), false)
				}
				cur = nil
			} else {
				cur = []point{p}
			}
			on = !on
			idx = (idx + 1) % len(pattern)
			left = pattern[idx]
		}
		left -= segLen - pos
		if on {
			cur = append(cur, b)
		}
	}
	if on && len(cur) > 1 {
		s.piece(cur, false)
	}
}
