package idml

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// Coordinates are points, y down. Every page item carries an ItemTransform
// (a b c d tx ty) that maps its own coordinates, in which its path points
// are given, to those of its parent: the spread for an item on a page, the
// group for a member of one. A page too has one, placing it in the spread;
// its GeometricBounds (top left bottom right) are in the page's own
// coordinates.

// parseTransform reads an ItemTransform attribute (the identity when the
// element has none or it is malformed).
func parseTransform(n *ooxml.Node) canvas.Matrix {
	f := nums(n.AttrStr("ItemTransform", ""))
	if len(f) != 6 {
		return canvas.Identity
	}
	m := canvas.Matrix{f[0], f[1], f[2], f[3], f[4], f[5]}
	if m[0]*m[3]-m[1]*m[2] == 0 {
		return canvas.Identity
	}
	return m
}

// invert returns the inverse of an invertible transform.
func invert(m canvas.Matrix) canvas.Matrix {
	det := m[0]*m[3] - m[1]*m[2]
	if det == 0 {
		return canvas.Identity
	}
	return canvas.Matrix{
		m[3] / det, -m[1] / det,
		-m[2] / det, m[0] / det,
		(m[2]*m[5] - m[3]*m[4]) / det, (m[1]*m[4] - m[0]*m[5]) / det,
	}
}

// nums parses a list of numbers separated by spaces.
func nums(s string) []float64 {
	fields := strings.Fields(s)
	out := make([]float64, 0, len(fields))
	for _, f := range fields {
		v, err := strconv.ParseFloat(f, 64)
		if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
			return nil
		}
		out = append(out, v)
	}
	return out
}

// bounds reads a GeometricBounds attribute (top left bottom right) as a
// rectangle.
func bounds(n *ooxml.Node) (r rect, ok bool) {
	f := nums(n.AttrStr("GeometricBounds", ""))
	if len(f) != 4 {
		return rect{}, false
	}
	return rect{x0: f[1], y0: f[0], x1: f[3], y1: f[2]}, true
}

type rect struct{ x0, y0, x1, y1 float64 }

func (r rect) w() float64 { return r.x1 - r.x0 }
func (r rect) h() float64 { return r.y1 - r.y0 }

func (r *rect) add(x, y float64) {
	if r.x0 > r.x1 {
		r.x0, r.y0, r.x1, r.y1 = x, y, x, y
		return
	}
	r.x0, r.y0 = math.Min(r.x0, x), math.Min(r.y0, y)
	r.x1, r.y1 = math.Max(r.x1, x), math.Max(r.y1, y)
}

// emptyRect is a rectangle that add starts from.
var emptyRect = rect{x0: 1, x1: 0}

// pathPoint is a point of a path with its Bézier handles.
type pathPoint struct{ ax, ay, lx, ly, rx, ry float64 }

// subpath is one GeometryPathType.
type subpath struct {
	pts  []pathPoint
	open bool
}

// geometry reads the PathGeometry of an item: its subpaths in the item's
// coordinates.
func geometry(n *ooxml.Node) []subpath {
	var out []subpath
	for _, g := range n.Path("Properties", "PathGeometry").Children("GeometryPathType") {
		sp := subpath{open: g.AttrBool("PathOpen", false)}
		for _, p := range g.Child("PathPointArray").Children("PathPointType") {
			a := nums(p.AttrStr("Anchor", ""))
			if len(a) != 2 {
				continue
			}
			pt := pathPoint{ax: a[0], ay: a[1], lx: a[0], ly: a[1], rx: a[0], ry: a[1]}
			if l := nums(p.AttrStr("LeftDirection", "")); len(l) == 2 {
				pt.lx, pt.ly = l[0], l[1]
			}
			if r := nums(p.AttrStr("RightDirection", "")); len(r) == 2 {
				pt.rx, pt.ry = r[0], r[1]
			}
			sp.pts = append(sp.pts, pt)
		}
		if len(sp.pts) > 0 {
			out = append(out, sp)
		}
	}
	return out
}

// bbox returns the box the anchors and handles of subpaths span.
func bbox(sps []subpath) rect {
	r := emptyRect
	for _, sp := range sps {
		for _, p := range sp.pts {
			r.add(p.ax, p.ay)
			r.add(p.lx, p.ly)
			r.add(p.rx, p.ry)
		}
	}
	if r.x0 > r.x1 {
		return rect{}
	}
	return r
}

// anchorBox returns the box the anchors alone span: the frame of a
// rectangle (whose handles sit on its anchors) or of a text frame.
func anchorBox(sps []subpath) rect {
	r := emptyRect
	for _, sp := range sps {
		for _, p := range sp.pts {
			r.add(p.ax, p.ay)
		}
	}
	if r.x0 > r.x1 {
		return rect{}
	}
	return r
}

// isRect reports whether the subpaths are one closed rectangle with
// straight sides: four anchors whose handles sit on them, two of them on
// each vertical and each horizontal line.
func isRect(sps []subpath) bool {
	if len(sps) != 1 || sps[0].open || len(sps[0].pts) != 4 {
		return false
	}
	for _, p := range sps[0].pts {
		if p.lx != p.ax || p.ly != p.ay || p.rx != p.ax || p.ry != p.ay {
			return false
		}
	}
	p := sps[0].pts
	return (p[0].ax == p[1].ax && p[1].ay == p[2].ay && p[2].ax == p[3].ax && p[3].ay == p[0].ay) ||
		(p[0].ay == p[1].ay && p[1].ax == p[2].ax && p[2].ay == p[3].ay && p[3].ax == p[0].ax)
}

// buildPath builds the drawing path of subpaths; closeAll closes open subpaths
// too (for filling).
func buildPath(sps []subpath, closeAll bool) *bdf.Path {
	p := &bdf.Path{}
	for _, sp := range sps {
		pts := sp.pts
		p.MoveTo(f32(pts[0].ax), f32(pts[0].ay))
		n := len(pts)
		last := n - 1
		if !sp.open {
			last = n
		}
		for i := 1; i <= last; i++ {
			a, b := pts[i-1], pts[i%n]
			if a.rx == a.ax && a.ry == a.ay && b.lx == b.ax && b.ly == b.ay {
				p.LineTo(f32(b.ax), f32(b.ay))
			} else {
				p.CubicTo(f32(a.rx), f32(a.ry), f32(b.lx), f32(b.ly), f32(b.ax), f32(b.ay))
			}
		}
		if !sp.open || closeAll {
			p.Close()
		}
	}
	return p
}

// roundedRect builds a rectangle whose corners are rounded with radius r
// (InDesign's RoundedCorner option, applied to every corner).
func roundedRect(b rect, r float64) *bdf.Path {
	r = math.Min(r, math.Min(b.w(), b.h())/2)
	p := &bdf.Path{}
	p.RoundRect(f32(b.x0), f32(b.y0), f32(b.w()), f32(b.h()), f32(r))
	return p
}

func f32(v float64) float32 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return float32(v)
}
