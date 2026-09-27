package raster

import (
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

// strokePolys returns the outline of a stroke of the polylines (user
// space) as polygons that all wind the same way, for a nonzero fill.
// devScale is the device pixels per user unit, for the smoothness of round
// parts.
func strokePolys(lines []polyline, ls *lineStyle, devScale float64) []polyline {
	hw := ls.width / 2
	if !(hw > 0) {
		return nil
	}
	tol := 0.2 / math.Max(devScale, 1e-9)
	var out []polyline
	add := func(pts ...point) {
		if signedArea(pts) < 0 {
			reversePoints(pts)
		}
		out = append(out, polyline{pts: pts, closed: true})
	}
	for _, pl := range lines {
		pts := dedupe(pl.pts)
		closed := pl.closed
		if len(pts) == 1 {
			// a zero-length sub-path: a dot under round and square caps
			if len(pl.pts) >= 2 {
				c := pts[0]
				switch ls.cap {
				case bdf.CapRound:
					add(circle(c, hw, tol)...)
				case bdf.CapSquare:
					add(point{c.x - hw, c.y - hw}, point{c.x + hw, c.y - hw}, point{c.x + hw, c.y + hw}, point{c.x - hw, c.y + hw})
				}
			}
			continue
		}
		if len(pts) < 2 {
			continue
		}
		pieces := [][]point{pts}
		pieceClosed := closed
		if ls.dash != nil {
			if closed {
				pts = append(pts, pts[0])
			}
			pieces = dashPolyline(pts, ls.dash, ls.dashOff)
			pieceClosed = false
		}
		for _, p := range pieces {
			strokePiece(p, pieceClosed, ls, hw, tol, add)
		}
	}
	return out
}

// strokePiece outlines one polyline: a quadrilateral per segment, a join at
// each corner and caps at the ends of an open one.
func strokePiece(pts []point, closed bool, ls *lineStyle, hw, tol float64, add func(...point)) {
	pts = dedupe(pts)
	if closed && len(pts) > 2 && pts[0] == pts[len(pts)-1] {
		pts = pts[:len(pts)-1]
	}
	n := len(pts)
	if n < 2 {
		if n == 1 && ls.cap != bdf.CapButt {
			c := pts[0]
			if ls.cap == bdf.CapRound {
				add(circle(c, hw, tol)...)
			} else {
				add(point{c.x - hw, c.y - hw}, point{c.x + hw, c.y - hw}, point{c.x + hw, c.y + hw}, point{c.x - hw, c.y + hw})
			}
		}
		return
	}
	segs := n - 1
	if closed {
		segs = n
	}
	for i := 0; i < segs; i++ {
		a, b := pts[i], pts[(i+1)%n]
		nx, ny := normal(a, b, hw)
		add(point{a.x + nx, a.y + ny}, point{b.x + nx, b.y + ny}, point{b.x - nx, b.y - ny}, point{a.x - nx, a.y - ny})
	}
	joinAt := func(prev, v, next point) {
		join(prev, v, next, ls, hw, tol, add)
	}
	for i := 1; i < n-1; i++ {
		joinAt(pts[i-1], pts[i], pts[i+1])
	}
	if closed {
		joinAt(pts[n-2], pts[n-1], pts[0])
		joinAt(pts[n-1], pts[0], pts[1])
		return
	}
	cap := func(end, from point) {
		switch ls.cap {
		case bdf.CapRound:
			add(circle(end, hw, tol)...)
		case bdf.CapSquare:
			dx, dy := end.x-from.x, end.y-from.y
			l := math.Hypot(dx, dy)
			ux, uy := dx/l*hw, dy/l*hw
			nx, ny := -uy, ux
			add(point{end.x + nx, end.y + ny}, point{end.x + nx + ux, end.y + ny + uy}, point{end.x - nx + ux, end.y - ny + uy}, point{end.x - nx, end.y - ny})
		}
	}
	cap(pts[0], pts[1])
	cap(pts[n-1], pts[n-2])
}

// join fills the corner at v between the segments prev→v and v→next.
func join(prev, v, next point, ls *lineStyle, hw, tol float64, add func(...point)) {
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
			add(v, p1, p2)
			return
		}
		add(circle(v, hw, tol)...)
	case bdf.JoinMiter:
		// miter length over line width is 1/sin(θ/2), θ the angle between the segments
		theta := math.Pi - turn
		limit := ls.miter
		if !(limit > 0) {
			limit = 10
		}
		if s := math.Sin(theta / 2); s > 0 && 1/s <= limit {
			d := hw / math.Cos(turn/2)
			mx, my := (n1x+n2x)/2, (n1y+n2y)/2
			ml := math.Hypot(mx, my)
			if ml > 0 {
				add(v, p1, point{v.x + mx/ml*d, v.y + my/ml*d}, p2)
				return
			}
		}
		add(v, p1, p2)
	default:
		add(v, p1, p2)
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
		n = int(math.Ceil(2 * math.Pi / step))
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

// dashPolyline cuts a polyline into the dashes of a pattern (on, off, …)
// starting offset into it.
func dashPolyline(pts []point, pattern []float64, offset float64) [][]point {
	total := 0.0
	for _, v := range pattern {
		total += v
	}
	if total <= 0 {
		return [][]point{pts}
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
	var out [][]point
	var cur []point
	on := idx%2 == 0
	if on {
		cur = []point{pts[0]}
	}
	guard := 0
	for i := 1; i < len(pts); i++ {
		a, b := pts[i-1], pts[i]
		segLen := math.Hypot(b.x-a.x, b.y-a.y)
		pos := 0.0
		for segLen-pos > left {
			guard++
			if guard > 1<<20 {
				return out // a pattern far too fine for the line
			}
			pos += left
			t := pos / segLen
			p := point{a.x + (b.x-a.x)*t, a.y + (b.y-a.y)*t}
			if on {
				cur = append(cur, p)
				out = append(out, cur)
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
		out = append(out, cur)
	}
	return out
}
