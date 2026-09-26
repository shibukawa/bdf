package visio

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// Geometry sections describe a shape's outline in its local coordinates
// (inches, y up, origin at the bottom left of its box). Arcs, elliptical
// arcs and ellipses become cubic Béziers; NURBS and B-splines are sampled.
// The points are mapped to the page as they are added.

type pt struct{ x, y float64 }

func (a pt) add(b pt) pt             { return pt{a.x + b.x, a.y + b.y} }
func (a pt) sub(b pt) pt             { return pt{a.x - b.x, a.y - b.y} }
func (a pt) mul(k float64) pt        { return pt{a.x * k, a.y * k} }
func (a pt) len() float64            { return math.Hypot(a.x, a.y) }
func (a pt) dist(b pt) float64       { return a.sub(b).len() }
func apply(m canvas.Matrix, p pt) pt { x, y := m.Apply(p.x, p.y); return pt{x, y} }

// seg is a line or a cubic Bézier to p.
type seg struct {
	cubic  bool
	c1, c2 pt
	p      pt
}

// subpath is a figure of connected segments.
type subpath struct {
	start pt
	segs  []seg
}

func (sp *subpath) end() pt {
	if len(sp.segs) == 0 {
		return sp.start
	}
	return sp.segs[len(sp.segs)-1].p
}

// closed reports whether the figure ends where it starts: Visio fills only
// closed figures.
func (sp *subpath) closed() bool {
	return len(sp.segs) > 0 && sp.end().dist(sp.start) < 1e-3
}

// startDir and endDir are the directions of travel at the ends.
func (sp *subpath) startDir() pt {
	for _, s := range sp.segs {
		for _, q := range s.points() {
			if d := q.sub(sp.start); d.len() > 1e-6 {
				return d
			}
		}
	}
	return pt{}
}

func (sp *subpath) endDir() pt {
	e := sp.end()
	for i := len(sp.segs) - 1; i >= 0; i-- {
		prev := sp.start
		if i > 0 {
			prev = sp.segs[i-1].p
		}
		s := sp.segs[i]
		cands := []pt{prev}
		if s.cubic {
			cands = []pt{s.c2, s.c1, prev}
		}
		for _, q := range cands {
			if d := e.sub(q); d.len() > 1e-6 {
				return d
			}
		}
	}
	return pt{}
}

func (s seg) points() []pt {
	if s.cubic {
		return []pt{s.c1, s.c2, s.p}
	}
	return []pt{s.p}
}

// pathBuilder collects the figures of a geometry section in page space.
type pathBuilder struct {
	m    canvas.Matrix // local → page
	subs []*subpath
	cur  pt // current point, local
	w, h float64
}

func (b *pathBuilder) moveTo(p pt) {
	b.subs = append(b.subs, &subpath{start: apply(b.m, p)})
	b.cur = p
}

func (b *pathBuilder) last() *subpath {
	if len(b.subs) == 0 {
		b.moveTo(b.cur)
	}
	return b.subs[len(b.subs)-1]
}

func (b *pathBuilder) lineTo(p pt) {
	sp := b.last()
	sp.segs = append(sp.segs, seg{p: apply(b.m, p)})
	b.cur = p
}

func (b *pathBuilder) cubicTo(c1, c2, p pt) {
	sp := b.last()
	sp.segs = append(sp.segs, seg{cubic: true, c1: apply(b.m, c1), c2: apply(b.m, c2), p: apply(b.m, p)})
	b.cur = p
}

// arcCenter is an arc of a circle about c from angle a0 sweeping by sweep
// (radians, counterclockwise when positive), as Béziers mapped by t.
func (b *pathBuilder) arc(t canvas.Matrix, c pt, r, a0, sweep float64) {
	n := int(math.Ceil(math.Abs(sweep) / (math.Pi / 2)))
	if n < 1 {
		n = 1
	}
	d := sweep / float64(n)
	k := 4.0 / 3 * math.Tan(d/4)
	for i := 0; i < n; i++ {
		s0, s1 := a0+float64(i)*d, a0+float64(i+1)*d
		p0 := pt{c.x + r*math.Cos(s0), c.y + r*math.Sin(s0)}
		p1 := pt{c.x + r*math.Cos(s1), c.y + r*math.Sin(s1)}
		c1 := pt{p0.x - k*r*math.Sin(s0), p0.y + k*r*math.Cos(s0)}
		c2 := pt{p1.x + k*r*math.Sin(s1), p1.y - k*r*math.Cos(s1)}
		b.cubicTo(apply(t, c1), apply(t, c2), apply(t, p1))
	}
}

// arcThrough draws the arc of a circle from the current point through mid
// to end; points in the space t maps to local coordinates (identity for
// circular arcs, a stretch and a rotation for elliptical ones).
func (b *pathBuilder) arcThrough(t canvas.Matrix, p0, mid, p1 pt) {
	// circumcenter
	ax, ay := p0.x, p0.y
	bx, by := mid.x, mid.y
	cx, cy := p1.x, p1.y
	d := 2 * (ax*(by-cy) + bx*(cy-ay) + cx*(ay-by))
	if math.Abs(d) < 1e-12 {
		b.lineTo(apply(t, p1))
		return
	}
	ux := ((ax*ax+ay*ay)*(by-cy) + (bx*bx+by*by)*(cy-ay) + (cx*cx+cy*cy)*(ay-by)) / d
	uy := ((ax*ax+ay*ay)*(cx-bx) + (bx*bx+by*by)*(ax-cx) + (cx*cx+cy*cy)*(bx-ax)) / d
	c := pt{ux, uy}
	r := c.dist(p0)
	a0 := math.Atan2(ay-uy, ax-ux)
	am := math.Atan2(by-uy, bx-ux)
	a1 := math.Atan2(cy-uy, cx-ux)
	norm := func(a float64) float64 {
		for a < 0 {
			a += 2 * math.Pi
		}
		for a >= 2*math.Pi {
			a -= 2 * math.Pi
		}
		return a
	}
	// counterclockwise sweep from a0 to a1; clockwise when mid is not on it
	ccw := norm(a1 - a0)
	if norm(am-a0) > ccw {
		ccw -= 2 * math.Pi
	}
	b.arc(t, c, r, a0, ccw)
}

// ellipse adds a closed ellipse given its center and the ends of two
// conjugate half axes.
func (b *pathBuilder) ellipse(c, a, e pt) {
	u, v := a.sub(c), e.sub(c)
	const k = 0.5522847498
	p := func(s, t float64) pt { return c.add(u.mul(s)).add(v.mul(t)) }
	b.moveTo(p(1, 0))
	b.cubicTo(p(1, k), p(k, 1), p(0, 1))
	b.cubicTo(p(-k, 1), p(-1, k), p(-1, 0))
	b.cubicTo(p(-1, -k), p(-k, -1), p(0, -1))
	b.cubicTo(p(k, -1), p(1, -k), p(1, 0))
}

// ellipticalArc draws an arc of an ellipse from the current point through
// ctrl to end; angle is the ellipse's major axis angle and ratio the major
// to minor axis ratio.
func (b *pathBuilder) ellipticalArc(end, ctrl pt, angle, ratio float64) {
	if ratio <= 1e-9 || math.IsNaN(ratio) {
		b.lineTo(end)
		return
	}
	// into a space where the ellipse is a circle
	toCircle := canvas.Scale(1/ratio, 1).Mul(canvas.Rotate(-angle * 180 / math.Pi))
	back := canvas.Rotate(angle * 180 / math.Pi).Mul(canvas.Scale(ratio, 1))
	b.arcThrough(back, apply(toCircle, b.cur), apply(toCircle, ctrl), apply(toCircle, end))
}

// bowArc draws a circular arc to end whose middle is bow away from the
// chord: positive bows lie to the right of the direction of travel.
func (b *pathBuilder) bowArc(end pt, bow float64) {
	d := end.sub(b.cur)
	l := d.len()
	if math.Abs(bow) < 1e-9 || l < 1e-12 {
		b.lineTo(end)
		return
	}
	n := pt{d.y / l, -d.x / l}
	mid := b.cur.add(end).mul(0.5).add(n.mul(bow))
	b.arcThrough(canvas.Identity, b.cur, mid, end)
}

// splinePoints samples a rational B-spline (weights nil: non-rational).
func splinePoints(ctrl []pt, weights, knots []float64, degree int) []pt {
	n := len(ctrl)
	if degree < 1 || n <= degree || len(knots) < n+degree+1 {
		return ctrl
	}
	lo, hi := knots[degree], knots[n]
	if !(hi > lo) {
		return ctrl
	}
	w := func(i int) float64 {
		if weights == nil || weights[i] == 0 {
			return 1
		}
		return weights[i]
	}
	eval := func(u float64) pt {
		// the span
		k := degree
		for k < n-1 && u >= knots[k+1] {
			k++
		}
		type hp struct{ x, y, w float64 }
		dd := make([]hp, degree+1)
		for j := 0; j <= degree; j++ {
			i := k - degree + j
			wi := w(i)
			dd[j] = hp{ctrl[i].x * wi, ctrl[i].y * wi, wi}
		}
		for r := 1; r <= degree; r++ {
			for j := degree; j >= r; j-- {
				i := k - degree + j
				den := knots[i+degree-r+1] - knots[i]
				a := 0.0
				if den != 0 {
					a = (u - knots[i]) / den
				}
				dd[j] = hp{(1-a)*dd[j-1].x + a*dd[j].x, (1-a)*dd[j-1].y + a*dd[j].y, (1-a)*dd[j-1].w + a*dd[j].w}
			}
		}
		if dd[degree].w == 0 {
			return pt{dd[degree].x, dd[degree].y}
		}
		return pt{dd[degree].x / dd[degree].w, dd[degree].y / dd[degree].w}
	}
	steps := 16 * (n - degree)
	if steps > 512 {
		steps = 512
	}
	out := make([]pt, 0, steps+1)
	for i := 1; i <= steps; i++ {
		u := lo + (hi-lo)*float64(i)/float64(steps)
		if i == steps {
			u = hi - 1e-12*(hi-lo)
		}
		out = append(out, eval(u))
	}
	return out
}

// funcArgs parses the arguments of a POLYLINE(…) or NURBS(…) value.
func funcArgs(v, name string) []float64 {
	v = strings.TrimSpace(v)
	i := strings.Index(strings.ToUpper(v), name+"(")
	if i < 0 {
		return nil
	}
	v = v[i+len(name)+1:]
	if j := strings.LastIndexByte(v, ')'); j >= 0 {
		v = v[:j]
	}
	var out []float64
	for _, f := range strings.Split(v, ",") {
		out = append(out, num(f))
	}
	return out
}

// buildGeometry turns a geometry section into figures in page space.
func buildGeometry(g *geoSection, m canvas.Matrix, w, h float64) []*subpath {
	b := &pathBuilder{m: m, w: w, h: h}
	for _, r := range g.rows {
		r.w, r.h = w, h
	}
	rel := func(r *geoRow, xn, yn string) pt { return pt{r.num(xn) * w, r.num(yn) * h} }
	abs := func(r *geoRow, xn, yn string) pt { return pt{r.num(xn), r.num(yn)} }
	// A section that does not start with a move starts at its last point.
	if len(g.rows) > 0 {
		switch g.rows[0].t {
		case "MoveTo", "RelMoveTo", "Ellipse", "InfiniteLine", "SplineStart":
		default:
			last := g.rows[len(g.rows)-1]
			if strings.HasPrefix(last.t, "Rel") {
				b.moveTo(rel(last, "X", "Y"))
			} else {
				b.moveTo(abs(last, "X", "Y"))
			}
		}
	}
	for i := 0; i < len(g.rows); i++ {
		r := g.rows[i]
		switch r.t {
		case "MoveTo":
			b.moveTo(abs(r, "X", "Y"))
		case "RelMoveTo":
			b.moveTo(rel(r, "X", "Y"))
		case "LineTo":
			b.lineTo(abs(r, "X", "Y"))
		case "RelLineTo":
			b.lineTo(rel(r, "X", "Y"))
		case "ArcTo":
			b.bowArc(abs(r, "X", "Y"), r.num("A"))
		case "EllipticalArcTo":
			b.ellipticalArc(abs(r, "X", "Y"), abs(r, "A", "B"), r.num("C"), r.num("D"))
		case "RelEllipticalArcTo":
			b.ellipticalArc(rel(r, "X", "Y"), rel(r, "A", "B"), r.num("C"), r.num("D"))
		case "Ellipse":
			b.ellipse(abs(r, "X", "Y"), abs(r, "A", "B"), abs(r, "C", "D"))
		case "RelCubBezTo":
			b.cubicTo(rel(r, "A", "B"), rel(r, "C", "D"), rel(r, "X", "Y"))
		case "RelQuadBezTo":
			q, e := rel(r, "A", "B"), rel(r, "X", "Y")
			p0 := b.cur
			b.cubicTo(p0.add(q.sub(p0).mul(2.0/3)), e.add(q.sub(e).mul(2.0/3)), e)
		case "PolylineTo":
			args := funcArgs(r.str("A"), "POLYLINE")
			if len(args) >= 2 {
				xRel, yRel := args[0] == 0, args[1] == 0
				for j := 2; j+1 < len(args); j += 2 {
					x, y := args[j], args[j+1]
					if xRel {
						x *= w
					}
					if yRel {
						y *= h
					}
					b.lineTo(pt{x, y})
				}
			}
			b.lineTo(abs(r, "X", "Y"))
		case "NURBSTo":
			b.nurbs(r, w, h)
		case "SplineStart":
			// a B-spline through the knot rows that follow
			ctrl := []pt{b.cur, abs(r, "X", "Y")}
			knots := []float64{r.num("B"), r.num("A")}
			last, degree := r.num("C"), int(r.num("D"))
			for i+1 < len(g.rows) && g.rows[i+1].t == "SplineKnot" {
				i++
				ctrl = append(ctrl, abs(g.rows[i], "X", "Y"))
				knots = append(knots, g.rows[i].num("A"))
			}
			knots = clampKnots(append(knots, last), len(ctrl), degree)
			for _, q := range splinePoints(ctrl, nil, knots, degree) {
				b.lineTo(q)
			}
		case "SplineKnot":
			b.lineTo(abs(r, "X", "Y"))
		}
	}
	return b.subs
}

// nurbs draws a NURBSTo row: its E cell holds the knot of the end, the
// degree, how its coordinates are given and the inner control points.
func (b *pathBuilder) nurbs(r *geoRow, w, h float64) {
	end := pt{r.num("X"), r.num("Y")}
	args := funcArgs(r.str("E"), "NURBS")
	if len(args) < 4 {
		b.lineTo(end)
		return
	}
	lastKnot, degree := args[0], int(args[1])
	xRel, yRel := args[2] == 0, args[3] == 0
	ctrl := []pt{b.cur}
	weights := []float64{r.num("D")}
	knots := []float64{r.num("C")}
	for j := 4; j+3 < len(args); j += 4 {
		x, y := args[j], args[j+1]
		if xRel {
			x *= w
		}
		if yRel {
			y *= h
		}
		ctrl = append(ctrl, pt{x, y})
		knots = append(knots, args[j+2])
		weights = append(weights, args[j+3])
	}
	ctrl = append(ctrl, end)
	knots = clampKnots(append(knots, r.num("A"), lastKnot), len(ctrl), degree)
	weights = append(weights, r.num("B"))
	for _, q := range splinePoints(ctrl, weights, knots, degree) {
		b.lineTo(q)
	}
	b.lineTo(end)
}

// clampKnots completes the knots Visio stores for a B-spline (one per
// control point and a last one) into a clamped knot vector for n control
// points: the first knot starts it degree+1 times and the last ends it, so
// that the curve runs from the first control point to the last.
func clampKnots(k []float64, n, degree int) []float64 {
	if len(k) == 0 || degree < 1 {
		return k
	}
	lead := 1
	for lead < len(k) && k[lead] == k[0] {
		lead++
	}
	for ; lead < degree+1; lead++ {
		k = append([]float64{k[0]}, k...)
	}
	for len(k) < n+degree+1 {
		k = append(k, k[len(k)-1])
	}
	return k[:n+degree+1]
}
