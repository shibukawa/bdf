package drawio

import "math"

// Geometry in diagram units (draw.io pixels, 1/96 inch). The layer objects
// scale them to points once, at their top.

// point is a position in diagram units.
type point struct{ x, y float64 }

// rect is an axis-aligned rectangle in diagram units.
type rect struct{ x, y, w, h float64 }

func (r rect) cx() float64 { return r.x + r.w/2 }
func (r rect) cy() float64 { return r.y + r.h/2 }

func (r rect) center() point { return point{r.cx(), r.cy()} }

// grow returns r enlarged by d on every side.
func (r rect) grow(d float64) rect { return rect{r.x - d, r.y - d, r.w + 2*d, r.h + 2*d} }

// union returns the smallest rectangle holding r and o; an empty (zero)
// rectangle is the identity.
func (r rect) union(o rect) rect {
	if r == (rect{}) {
		return o
	}
	if o == (rect{}) {
		return r
	}
	x0, y0 := math.Min(r.x, o.x), math.Min(r.y, o.y)
	x1, y1 := math.Max(r.x+r.w, o.x+o.w), math.Max(r.y+r.h, o.y+o.h)
	return rect{x0, y0, x1 - x0, y1 - y0}
}

// addPoint returns r grown to hold p.
func (r rect) addPoint(p point) rect {
	if r == (rect{}) {
		return rect{p.x, p.y, 0, 0}
	}
	return r.union(rect{p.x, p.y, 0, 0})
}

// rotate90 turns r by 90° about its center (mxRectangle.rotate90).
func (r rect) rotate90() rect {
	t := (r.w - r.h) / 2
	return rect{r.x + t, r.y - t, r.h, r.w}
}

// contains reports whether p lies in r (edges included), as mxUtils.contains.
func (r rect) contains(x, y float64) bool {
	return r.x <= x && r.x+r.w >= x && r.y <= y && r.y+r.h >= y
}

// rotatePoint turns p about c by the angle whose cosine and sine are given
// (mxUtils.getRotatedPoint).
func rotatePoint(p point, cos, sin float64, c point) point {
	x, y := p.x-c.x, p.y-c.y
	return point{x*cos - y*sin + c.x, y*cos + x*sin + c.y}
}

func toRadians(deg float64) float64 { return deg * math.Pi / 180 }

// matrix is an affine transform [a b c d e f] (Canvas order).
type matrix [6]float64

var identity = matrix{1, 0, 0, 1, 0, 0}

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

func translateM(x, y float64) matrix { return matrix{1, 0, 0, 1, x, y} }
func scaleM(x, y float64) matrix     { return matrix{x, 0, 0, y, 0, 0} }
func rotateM(deg float64) matrix {
	s, c := math.Sincos(toRadians(deg))
	return matrix{c, s, -s, c, 0, 0}
}

// rotateAbout rotates by deg about (cx, cy).
func rotateAbout(deg, cx, cy float64) matrix {
	return translateM(cx, cy).mul(rotateM(deg)).mul(translateM(-cx, -cy))
}

// mod is a modulo whose result has the sign of n (mxUtils.mod).
func mod(a, n int) int { return ((a % n) + n) % n }

// ptSegDistSq is the squared distance of (px, py) from the segment
// (x1, y1)-(x2, y2) (mxUtils.ptSegDistSq).
func ptSegDistSq(x1, y1, x2, y2, px, py float64) float64 {
	x2 -= x1
	y2 -= y1
	px -= x1
	py -= y1
	dot := px*x2 + py*y2
	var proj float64
	if dot <= 0 {
		proj = 0
	} else {
		px = x2 - px
		py = y2 - py
		dot = px*x2 + py*y2
		if dot <= 0 {
			proj = 0
		} else {
			proj = dot * dot / (x2*x2 + y2*y2)
		}
	}
	l := px*px + py*py - proj
	if l < 0 {
		l = 0
	}
	return l
}

// intersection returns the intersection of the lines through (x0,y0)-(x1,y1)
// and (x2,y2)-(x3,y3) when it lies on both segments, within 1e-6 of their
// ends against rounding errors (mxUtils.intersection).
func intersection(x0, y0, x1, y1, x2, y2, x3, y3 float64) (point, bool) {
	denom := (y3-y2)*(x1-x0) - (x3-x2)*(y1-y0)
	nume_a := (x3-x2)*(y0-y2) - (y3-y2)*(x0-x2)
	nume_b := (x1-x0)*(y0-y2) - (y1-y0)*(x0-x2)
	ua := nume_a / denom
	ub := nume_b / denom
	const eps = 0.000001
	if ua >= -eps && ua <= 1+eps && ub >= -eps && ub <= 1+eps {
		return point{x0 + ua*(x1-x0), y0 + ua*(y1-y0)}, true
	}
	return point{}, false
}

func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
