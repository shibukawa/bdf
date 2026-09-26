package dxf

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// cursor reads tags in order.
type cursor struct {
	tags []tag
	i    int
}

func (c *cursor) more() bool { return c.i < len(c.tags) }

// want consumes the next tag when it has the code.
func (c *cursor) want(code int) (float64, bool) {
	if c.i < len(c.tags) && c.tags[c.i].code == code {
		c.i++
		return c.tags[c.i-1].f, true
	}
	return 0, false
}

// seek skips to the next tag with the code and consumes it.
func (c *cursor) seek(code int) (tag, bool) {
	for c.i < len(c.tags) {
		t := c.tags[c.i]
		c.i++
		if t.code == code {
			return t, true
		}
	}
	return tag{}, false
}

func (c *cursor) point(code int) (cad.Point, bool) {
	x, ok := c.want(code)
	if !ok {
		return cad.Point{}, false
	}
	y, _ := c.want(code + 10)
	c.want(code + 20)
	return cad.Point{X: x, Y: y}, true
}

type boundary struct {
	flags int
	path  *cad.Path
}

// hatchData is what a HATCH holds.
type hatchData struct {
	paths    []boundary
	solid    bool
	style    int
	lines    []cad.PatternLine
	gradient bool
	gname    string
	gangle   float64
	gshift   float64
	oneColor bool
	tint     float64
	colors   []bdf.Color
}

func (c *converter) readHatch(e *entity, dark bool) *hatchData {
	tags := e.sub("AcDbHatch").tags
	if len(tags) == 0 {
		tags = e.tags
	}
	cur := &cursor{tags: tags}
	h := &hatchData{solid: e.int(70, 0) == 1}
	if _, ok := cur.seek(91); !ok {
		return h
	}
	n := int(cur.tags[cur.i-1].f)
	for range n {
		if _, ok := cur.seek(92); !ok {
			break
		}
		flags := int(cur.tags[cur.i-1].f)
		path := &cad.Path{}
		if flags&2 != 0 {
			hasBulge, _ := cur.want(72)
			cur.want(73)
			nv, _ := cur.want(93)
			var vs []vertex
			for range int(nv) {
				p, ok := cur.point(10)
				if !ok {
					break
				}
				v := vertex{p: p}
				if hasBulge != 0 {
					v.bulge, _ = cur.want(42)
				}
				vs = append(vs, v)
			}
			for i, v := range vs {
				if i == 0 {
					path.MoveTo(v.p.X, v.p.Y)
				}
				path.BulgeTo(vs[(i+1)%len(vs)].p, v.bulge)
			}
			path.Close()
		} else {
			ne, _ := cur.want(93)
			for range int(ne) {
				typ, ok := cur.want(72)
				if !ok {
					break
				}
				readEdge(cur, int(typ), path)
			}
			path.Close()
		}
		if ns, ok := cur.want(97); ok {
			for range int(ns) {
				cur.want(330)
			}
		}
		h.paths = append(h.paths, boundary{flags: flags, path: path})
	}
	if v, ok := cur.want(75); ok {
		h.style = int(v)
	}
	cur.want(76)
	if !h.solid {
		cur.want(52)
		cur.want(41)
		cur.want(77)
		if nl, ok := cur.want(78); ok {
			for range int(nl) {
				var pl cad.PatternLine
				a, _ := cur.want(53)
				pl.Angle = a * math.Pi / 180
				pl.Base.X, _ = cur.want(43)
				pl.Base.Y, _ = cur.want(44)
				pl.Offset.X, _ = cur.want(45)
				pl.Offset.Y, _ = cur.want(46)
				nd, _ := cur.want(79)
				for range int(nd) {
					d, _ := cur.want(49)
					pl.Dash = append(pl.Dash, d)
				}
				h.lines = append(h.lines, pl)
			}
		}
	}
	// the gradient
	inColors := false
	for ; cur.more(); cur.i++ {
		t := cur.tags[cur.i]
		switch t.code {
		case 450:
			h.gradient = t.f != 0
		case 452:
			h.oneColor = t.f != 0
		case 453:
			inColors = true
		case 63:
			if inColors {
				h.colors = append(h.colors, c.aci(int(t.f), dark))
			}
		case 421:
			if inColors {
				col := rgb(int(t.f))
				if n := len(h.colors); n > 0 && cur.i > 0 && cur.tags[cur.i-1].code == 63 {
					h.colors[n-1] = col
				} else {
					h.colors = append(h.colors, col)
				}
			}
		case 460:
			h.gangle = t.f
			inColors = false
		case 461:
			h.gshift = t.f
		case 462:
			h.tint = t.f
		case 470:
			h.gname = strings.ToUpper(t.s)
		}
	}
	return h
}

// readEdge adds an edge of a boundary path.
func readEdge(cur *cursor, typ int, path *cad.Path) {
	join := func(p cad.Point) {
		if !path.Open() {
			path.MoveTo(p.X, p.Y)
		} else if path.Current() != p {
			path.LineTo(p.X, p.Y)
		}
	}
	switch typ {
	case 1:
		a, _ := cur.point(10)
		b, _ := cur.point(11)
		join(a)
		path.LineTo(b.X, b.Y)
	case 2:
		ctr, _ := cur.point(10)
		r, _ := cur.want(40)
		s, _ := cur.want(50)
		e, _ := cur.want(51)
		ccw, _ := cur.want(73)
		a0, a1 := arcAngles(s, e, ccw != 0)
		start := cad.Point{X: ctr.X + r*math.Cos(a0), Y: ctr.Y + r*math.Sin(a0)}
		join(start)
		path.Arc(ctr, r, a0, a1)
	case 3:
		ctr, _ := cur.point(10)
		maj, _ := cur.point(11)
		ratio, _ := cur.want(40)
		s, _ := cur.want(50)
		e, _ := cur.want(51)
		ccw, _ := cur.want(73)
		a0, a1 := arcAngles(s, e, ccw != 0)
		// angles to the parameters of the ellipse
		t0, t1 := angleToParam(ratio, a0), angleToParam(ratio, a1)
		if ccw != 0 {
			for t1 <= t0 {
				t1 += 2 * math.Pi
			}
		} else {
			for t1 >= t0 {
				t1 -= 2 * math.Pi
			}
		}
		minor := cad.Point{X: -maj.Y * ratio, Y: maj.X * ratio}
		st := ctr.Add(maj.Mul(math.Cos(t0))).Add(minor.Mul(math.Sin(t0)))
		join(st)
		path.EllipseArcAxes(ctr, maj, minor, t0, t1)
	case 4:
		deg, _ := cur.want(94)
		rational, _ := cur.want(73)
		cur.want(74)
		nk, _ := cur.want(95)
		nc, _ := cur.want(96)
		var knots, weights []float64
		var ctrl []cad.Point
		for range int(nk) {
			k, _ := cur.want(40)
			knots = append(knots, k)
		}
		for range int(nc) {
			p, ok := cur.point(10)
			if !ok {
				break
			}
			ctrl = append(ctrl, p)
			if w, ok := cur.want(42); ok {
				weights = append(weights, w)
			}
		}
		var fit []cad.Point
		if nf, ok := cur.want(97); ok {
			for range int(nf) {
				p, ok := cur.point(11)
				if !ok {
					break
				}
				fit = append(fit, p)
			}
		}
		var st, en *cad.Point
		if p, ok := cur.point(12); ok {
			st = &p
		}
		if p, ok := cur.point(13); ok {
			en = &p
		}
		seg := &cad.Path{}
		if len(ctrl) >= 2 {
			if rational == 0 {
				weights = nil
			}
			seg.BSpline(int(deg), knots, ctrl, weights)
		} else if len(fit) >= 2 {
			seg.Interpolate(fit, st, en)
		} else {
			return
		}
		path.Continue(seg)
	}
}

// arcAngles returns the start and end angles in radians of a hatch edge
// arc, in its direction: clockwise arcs store their angles negated.
func arcAngles(s, e float64, ccw bool) (float64, float64) {
	a0, a1 := s*math.Pi/180, e*math.Pi/180
	if ccw {
		for a1 <= a0 {
			a1 += 2 * math.Pi
		}
		if a1-a0 > 2*math.Pi {
			a1 -= 2 * math.Pi
		}
		return a0, a1
	}
	a0, a1 = -a0, -a1
	for a1 >= a0 {
		a1 -= 2 * math.Pi
	}
	if a0-a1 > 2*math.Pi {
		a1 += 2 * math.Pi
	}
	return a0, a1
}

// angleToParam converts an angle on an ellipse into its parameter.
func angleToParam(ratio, a float64) float64 {
	if ratio == 0 {
		return a
	}
	t := math.Atan2(math.Sin(a)/ratio, math.Cos(a))
	// the parameter nearest to the angle (the same number of turns)
	return t + 2*math.Pi*math.Round((a-t)/(2*math.Pi))
}

// hatch draws a HATCH: a solid or gradient fill, or the lines of its
// pattern, inside its boundary paths.
func (c *converter) hatch(e *entity, x *ctx, p props) {
	if e.typ == "MPOLYGON" {
		c.warnOnce("MPOLYGON", "MPOLYGON entities are not drawn")
		return
	}
	h := c.readHatch(e, x.dark)
	elev := e.vec3(10, [3]float64{})[2]
	m := x.m.Mul(ocs(e, elev))
	bound := &cad.Path{}
	for i, b := range h.paths {
		switch h.style {
		case 2:
			// ignore islands: the outer boundary only
			if b.flags&1 == 0 && i > 0 {
				continue
			}
		case 1:
			if b.flags&(1|16) == 0 && i > 0 {
				continue
			}
		}
		bound.Append(b.path)
	}
	if bound.Empty() {
		return
	}
	wb := bound.Transform(m)
	switch {
	case h.gradient:
		x.out.Fill(wb, cad.Fill{Gradient: c.gradient(h, bound, m, p)}, true)
	case h.solid || len(h.lines) == 0:
		x.out.Fill(wb, cad.Fill{Color: p.color}, true)
	default:
		pen := c.pen(e, x, p, m)
		pen.Dash, pen.DashOffset = nil, 0
		lines := make([]cad.PatternLine, len(h.lines))
		for i, l := range h.lines {
			dir := cad.Apply(m, cad.Point{X: math.Cos(l.Angle), Y: math.Sin(l.Angle)}).Sub(cad.Apply(m, cad.Point{}))
			off := cad.Apply(m, l.Offset).Sub(cad.Apply(m, cad.Point{}))
			s := dir.Len()
			dash := make([]float64, len(l.Dash))
			for j, d := range l.Dash {
				dash[j] = d * s
			}
			lines[i] = cad.PatternLine{Angle: math.Atan2(dir.Y, dir.X), Base: cad.Apply(m, l.Base), Offset: off, Dash: dash}
		}
		if !x.out.Hatch(wb, lines, pen) {
			c.warnOnce("densehatch", "hatch patterns too dense to draw are filled with a tint of their color")
			x.out.Fill(wb, cad.Fill{Color: withAlpha(p.color, 0.3)}, true)
		}
	}
}

// gradient approximates the gradient fills of AutoCAD with linear and
// radial gradients across the boundary.
func (c *converter) gradient(h *hatchData, bound *cad.Path, m canvas.Matrix, p props) *cad.Gradient {
	c1, c2 := p.color, bdf.RGB(255, 255, 255)
	if len(h.colors) > 0 {
		c1 = h.colors[0]
	}
	if len(h.colors) > 1 && !h.oneColor {
		c2 = h.colors[1]
	} else {
		// one color: towards black or white by the tint
		t := h.tint
		mix := func(a uint8, to float64) uint8 { return uint8(float64(a)*(1-to) + 255*to) }
		if t < 0.5 {
			k := 1 - t*2
			c2 = bdf.RGB(uint8(float64(uint8(c1>>24))*(1-k)), uint8(float64(uint8(c1>>16))*(1-k)), uint8(float64(uint8(c1>>8))*(1-k)))
		} else {
			k := (t - 0.5) * 2
			c2 = bdf.RGB(mix(uint8(c1>>24), k), mix(uint8(c1>>16), k), mix(uint8(c1>>8), k))
		}
	}
	b := bound.Bounds()
	ctr := cad.Point{X: (b.Min.X + b.Max.X) / 2, Y: (b.Min.Y + b.Max.Y) / 2}
	half := math.Hypot(b.W(), b.H()) / 2
	dir := cad.Point{X: math.Cos(h.gangle), Y: math.Sin(h.gangle)}
	// the extent of the boundary along the direction
	lo, hi := math.Inf(1), math.Inf(-1)
	for _, q := range []cad.Point{b.Min, {X: b.Max.X, Y: b.Min.Y}, b.Max, {X: b.Min.X, Y: b.Max.Y}} {
		v := q.Sub(ctr)
		d := v.X*dir.X + v.Y*dir.Y
		lo, hi = min(lo, d), max(hi, d)
	}
	p0, p1 := ctr.Add(dir.Mul(lo)), ctr.Add(dir.Mul(hi))
	shift := h.gshift
	name := strings.TrimPrefix(h.gname, "INV")
	inv := strings.HasPrefix(h.gname, "INV")
	a, bcol := c1, c2
	if inv {
		a, bcol = c2, c1
	}
	g := &cad.Gradient{P0: cad.Apply(m, p0), P1: cad.Apply(m, p1)}
	s := cad.Scale(m)
	switch name {
	case "CYLINDER":
		g.Stops = []bdf.Stop{{Offset: 0, Color: a}, {Offset: float32(0.5 + shift/2), Color: bcol}, {Offset: 1, Color: a}}
	case "SPHERICAL", "HEMISPHERICAL", "CURVED":
		center := ctr.Add(dir.Mul(shift * half))
		if name == "HEMISPHERICAL" {
			center = ctr.Sub(cad.Point{X: -dir.Y, Y: dir.X}.Mul(half / 2))
		}
		g.Radial = true
		g.P0, g.P1 = cad.Apply(m, center), cad.Apply(m, center)
		g.R0, g.R1 = 0, half*s*1.2
		g.Stops = []bdf.Stop{{Offset: 0, Color: bcol}, {Offset: 1, Color: a}}
	default: // LINEAR
		g.Stops = []bdf.Stop{{Offset: 0, Color: a}, {Offset: 1, Color: bcol}}
	}
	return g
}
