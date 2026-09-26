package sxf

import (
	"math"
	"slices"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// sfcRenderer draws the features of an SFC file.
type sfcRenderer struct {
	*renderer
	m     *model
	stack []string // the compound figures being drawn
}

func (r *sfcRenderer) color(code int) bdf.Color {
	if c, ok := r.m.color(code); ok {
		return c
	}
	return bdf.RGB(255, 255, 255)
}

// style reads the color, line type and width codes at i, i+1, i+2.
func (r *sfcRenderer) style(f *feature, i int) lineStyle {
	w := r.m.width(f.int(i + 2))
	return lineStyle{color: r.color(f.int(i)), width: w, dash: r.m.pattern(f.int(i+1), w)}
}

func pts(xs, ys []float64) []cad.Point {
	n := min(len(xs), len(ys))
	out := make([]cad.Point, n)
	for i := range n {
		out[i] = cad.Point{X: xs[i], Y: ys[i]}
	}
	return out
}

func deg(a float64) float64 { return a * math.Pi / 180 }

// arcPath returns the arc of an ellipse (radii rx, ry turned by rot)
// between two angles in degrees, counter-clockwise or clockwise.
func arcPath(ctr cad.Point, rx, ry, rot, a0, a1 float64, cw bool) *cad.Path {
	major := cad.Point{X: rx * math.Cos(deg(rot)), Y: rx * math.Sin(deg(rot))}
	ratio := 1.0
	if rx != 0 {
		ratio = ry / rx
	}
	t0, t1 := deg(a0), deg(a1)
	if math.IsNaN(t0) || math.IsInf(t0, 0) || math.IsNaN(t1) || math.IsInf(t1, 0) {
		return &cad.Path{}
	}
	t0, t1 = arcEnds(t0, t1, cw)
	return (&cad.Path{}).EllipseArc(ctr, major, ratio, t0, t1)
}

// arcEnds returns the parameters an arc from t0 to t1 (radians) runs
// between, counter-clockwise or clockwise: less than a turn, or a full one
// when the ends meet or are a turn apart (360°, or 2π written a hair over).
func arcEnds(t0, t1 float64, cw bool) (float64, float64) {
	const eps = 1e-9
	if math.Abs(t1-t0) >= 2*math.Pi-eps {
		if cw {
			return t0, t0 - 2*math.Pi
		}
		return t0, t0 + 2*math.Pi
	}
	t0, t1 = math.Mod(t0, 2*math.Pi), math.Mod(t1, 2*math.Pi)
	if cw {
		if t1 >= t0-eps {
			t1 -= 2 * math.Pi
		}
	} else if t1 <= t0+eps {
		t1 += 2 * math.Pi
	}
	return t0, t1
}

// items draws the elements of the sheet or of a compound figure through m.
func (r *sfcRenderer) items(items []*feature, m canvas.Matrix) {
	for _, f := range items {
		r.feature(f, m)
	}
}

func (r *sfcRenderer) feature(f *feature, m canvas.Matrix) {
	if !r.visit() {
		return
	}
	// every element but a hatch starts with its layer
	if !r.m.layerVisible(f.int(0)) {
		return
	}
	switch f.name {
	case "point_marker_feature":
		r.marker(f.int(4), cad.Point{X: f.num(2), Y: f.num(3)}, f.num(5), f.num(6), r.paint(r.color(f.int(1))), m)
	case "line_feature":
		p := (&cad.Path{}).MoveTo(f.num(4), f.num(5)).LineTo(f.num(6), f.num(7))
		r.stroke(p, r.style(f, 1), m)
	case "polyline_feature":
		r.stroke((&cad.Path{}).Polyline(pts(f.list(5), f.list(6)), false), r.style(f, 1), m)
	case "circle_feature":
		r.stroke((&cad.Path{}).Circle(cad.Point{X: f.num(4), Y: f.num(5)}, f.num(6)), r.style(f, 1), m)
	case "arc_feature":
		r.stroke(arcPath(cad.Point{X: f.num(4), Y: f.num(5)}, f.num(6), f.num(6), 0, f.num(8), f.num(9), f.int(7) == 1), r.style(f, 1), m)
	case "ellipse_feature":
		ctr := cad.Point{X: f.num(4), Y: f.num(5)}
		rx, ry, rot := f.num(6), f.num(7), f.num(8)
		major := cad.Point{X: rx * math.Cos(deg(rot)), Y: rx * math.Sin(deg(rot))}
		ratio := 1.0
		if rx != 0 {
			ratio = ry / rx
		}
		r.stroke((&cad.Path{}).Ellipse(ctr, major, ratio), r.style(f, 1), m)
	case "ellipse_arc_feature":
		r.stroke(arcPath(cad.Point{X: f.num(4), Y: f.num(5)}, f.num(6), f.num(7), f.num(9), f.num(10), f.num(11), f.int(8) == 1), r.style(f, 1), m)
	case "text_string_feature":
		r.text(r.textAt(f, 1, 2), m)
	case "spline_feature":
		r.stroke(bezierPath(pts(f.list(6), f.list(7))), r.style(f, 1), m)
	case "clothoid_feature":
		c := clothoid(f.num(6), f.num(9), f.num(10), f.int(7) == 1)
		cm := m.Mul(canvas.Translate(f.num(4), f.num(5))).Mul(canvas.Rotate(f.num(8)))
		r.stroke((&cad.Path{}).Polyline(c, false), r.style(f, 1), cm)
	case "sfig_locate_feature":
		r.locate(f, m)
	case "externally_defined_symbol_feature":
		r.warn("symbol", "predefined symbols (externally_defined_symbol_feature) are not drawn")
	case "linear_dim_feature":
		r.linearDim(f, m)
	case "curve_dim_feature", "angular_dim_feature":
		r.angularDim(f, m)
	case "radius_dim_feature":
		s := r.style(f, 1)
		a, b := cad.Point{X: f.num(4), Y: f.num(5)}, cad.Point{X: f.num(6), Y: f.num(7)}
		r.stroke((&cad.Path{}).MoveTo(a.X, a.Y).LineTo(b.X, b.Y), s, m)
		r.dimArrow(f, 8, a, b, s, m)
		if f.int(13) == 1 {
			r.text(r.textAt(f, 1, 14), m)
		}
	case "diameter_dim_feature":
		s := r.style(f, 1)
		a, b := cad.Point{X: f.num(4), Y: f.num(5)}, cad.Point{X: f.num(6), Y: f.num(7)}
		r.stroke((&cad.Path{}).MoveTo(a.X, a.Y).LineTo(b.X, b.Y), s, m)
		r.dimArrow(f, 8, a, b, s, m)
		r.dimArrow(f, 13, a, b, s, m)
		if f.int(18) == 1 {
			r.text(r.textAt(f, 1, 19), m)
		}
	case "label_feature", "balloon_feature":
		s := r.style(f, 1)
		v := pts(f.list(5), f.list(6))
		r.stroke((&cad.Path{}).Polyline(v, false), s, m)
		i := 7
		if f.name == "balloon_feature" {
			r.stroke((&cad.Path{}).Circle(cad.Point{X: f.num(7), Y: f.num(8)}, f.num(9)), s, m)
			i = 10
		}
		if len(v) >= 2 {
			r.arrow(f.int(i), v[0], v[0].Sub(v[1]), f.num(i+1), s, m)
		}
		if f.int(i+2) == 1 {
			r.text(r.textAt(f, 1, i+3), m)
		}
	case "externally_defined_hatch_feature", "fill_area_style_colour_feature", "fill_area_style_hatching_feature", "fill_area_style_tiles_feature":
		r.hatchFeature(f, m)
	default:
		r.m.unknown[f.name]++
	}
}

// textAt reads a text: its color code at ci and font, string, position,
// height, width, spacing, angle, slant, anchor and direction from i on.
func (r *sfcRenderer) textAt(f *feature, ci, i int) textBox {
	return textBox{color: r.color(f.int(ci)), font: r.m.font(f.int(i)), s: f.str(i + 1),
		x: f.num(i + 2), y: f.num(i + 3), h: f.num(i + 4), w: f.num(i + 5), spacing: f.num(i + 6),
		angle: f.num(i + 7), slant: f.num(i + 8), anchor: max(f.int(i+9), 1), vertical: f.int(i+10) == 2}
}

// bezierPath joins cubic Bézier segments: control points 0-3, 3-6, ...
func bezierPath(p []cad.Point) *cad.Path {
	path := &cad.Path{}
	if len(p) < 2 {
		return path
	}
	path.MoveTo(p[0].X, p[0].Y)
	i := 0
	for ; i+3 < len(p); i += 3 {
		path.CubicTo(p[i+1], p[i+2], p[i+3])
	}
	for i++; i < len(p); i++ {
		path.LineTo(p[i].X, p[i].Y)
	}
	return path
}

// locate draws a compound figure where it is placed.
func (r *sfcRenderer) locate(f *feature, m canvas.Matrix) {
	name := f.str(1)
	s := r.m.sfigs[name]
	if s == nil {
		r.warn("sfig:"+name, "compound figure %q is not defined", name)
		return
	}
	if _, ok := backgroundAttribute(name); ok {
		// a group that only carries the background color attribute
		return
	}
	if slices.Contains(r.stack, name) || len(r.stack) > 32 {
		r.warn("cycle", "compound figures that place themselves are drawn once")
		return
	}
	sm := m.Mul(canvas.Translate(f.num(2), f.num(3))).Mul(canvas.Rotate(f.num(4))).Mul(canvas.Scale(f.num(5), f.num(6)))
	if s.flag == 2 {
		// geodetic coordinates: x points north, y east
		sm = sm.Mul(canvas.Matrix{0, 1, 1, 0, 0, 0})
	}
	r.stack = append(r.stack, name)
	r.items(s.items, sm)
	r.stack = r.stack[:len(r.stack)-1]
}

// extension draws an extension line: flag at i, then base, start and end.
func (r *sfcRenderer) extension(f *feature, i int, s lineStyle, m canvas.Matrix) {
	if f.int(i) != 1 {
		return
	}
	p := (&cad.Path{}).MoveTo(f.num(i+3), f.num(i+4)).LineTo(f.num(i+5), f.num(i+6))
	r.stroke(p, s, m)
}

// dimArrow draws the arrow of a dimension whose code, direction code,
// position and scale start at i, on the line from a to b: pointing out of
// the line (direction 1) or into it (2).
func (r *sfcRenderer) dimArrow(f *feature, i int, a, b cad.Point, s lineStyle, m canvas.Matrix) {
	code, dir := f.int(i), f.int(i+1)
	if dir == 0 || code == 0 {
		return
	}
	tip := cad.Point{X: f.num(i + 2), Y: f.num(i + 3)}
	// away from the farther end of the line
	other := a
	if tip.Sub(a).Len() < tip.Sub(b).Len() {
		other = b
	}
	d := tip.Sub(other)
	if dir == 2 {
		d = d.Mul(-1)
	}
	r.arrow(code, tip, d, f.num(i+4), s, m)
}

func (r *sfcRenderer) linearDim(f *feature, m canvas.Matrix) {
	s := r.style(f, 1)
	a, b := cad.Point{X: f.num(4), Y: f.num(5)}, cad.Point{X: f.num(6), Y: f.num(7)}
	r.stroke((&cad.Path{}).MoveTo(a.X, a.Y).LineTo(b.X, b.Y), s, m)
	r.extension(f, 8, s, m)
	r.extension(f, 15, s, m)
	r.dimArrow(f, 22, a, b, s, m)
	r.dimArrow(f, 27, a, b, s, m)
	if f.int(32) == 1 {
		r.text(r.textAt(f, 1, 33), m)
	}
}

// angularDim draws an angular or arc length dimension: an arc around its
// origin, extension lines, arrows along the arc and the value.
func (r *sfcRenderer) angularDim(f *feature, m canvas.Matrix) {
	s := r.style(f, 1)
	ctr := cad.Point{X: f.num(4), Y: f.num(5)}
	rad := f.num(6)
	r.stroke(arcPath(ctr, rad, rad, 0, f.num(7), f.num(8), false), s, m)
	r.extension(f, 9, s, m)
	r.extension(f, 16, s, m)
	for k, i := range []int{23, 28} {
		code, dir := f.int(i), f.int(i+1)
		if code == 0 || dir == 0 {
			continue
		}
		tip := cad.Point{X: f.num(i + 2), Y: f.num(i + 3)}
		a := math.Atan2(tip.Y-ctr.Y, tip.X-ctr.X)
		// the tangent out of the arc: backwards at its start, forwards at
		// its end
		d := cad.Point{X: math.Sin(a), Y: -math.Cos(a)}
		if k == 1 {
			d = d.Mul(-1)
		}
		if dir == 2 {
			d = d.Mul(-1)
		}
		r.arrow(code, tip, d, f.num(i+4), s, m)
	}
	if f.int(33) == 1 {
		r.text(r.textAt(f, 1, 34), m)
	}
}

// boundary returns the closed outline of a composite curve (by code) in
// the coordinates of its members.
func (r *sfcRenderer) boundary(code int) *cad.Path {
	if code < 1 || code > len(r.m.curves) {
		return nil
	}
	c := r.m.curves[code-1]
	var poly []cad.Point
	add := func(p []cad.Point) {
		if len(p) == 0 {
			return
		}
		// continue from the end nearest to the outline so far
		if n := len(poly); n > 0 {
			last := poly[n-1]
			if p[len(p)-1].Sub(last).Len() < p[0].Sub(last).Len() {
				p = slices.Clone(p)
				slices.Reverse(p)
			}
		}
		poly = append(poly, p...)
	}
	for _, mf := range c.members {
		switch mf.name {
		case "polyline_feature":
			add(pts(mf.list(5), mf.list(6)))
		case "line_feature":
			add([]cad.Point{{X: mf.num(4), Y: mf.num(5)}, {X: mf.num(6), Y: mf.num(7)}})
		case "arc_feature":
			add(flatten(arcPath(cad.Point{X: mf.num(4), Y: mf.num(5)}, mf.num(6), mf.num(6), 0, mf.num(8), mf.num(9), mf.int(7) == 1)))
		case "ellipse_arc_feature":
			add(flatten(arcPath(cad.Point{X: mf.num(4), Y: mf.num(5)}, mf.num(6), mf.num(7), mf.num(9), mf.num(10), mf.num(11), mf.int(8) == 1)))
		case "circle_feature":
			add(flatten((&cad.Path{}).Circle(cad.Point{X: mf.num(4), Y: mf.num(5)}, mf.num(6))))
		case "spline_feature":
			add(flatten(bezierPath(pts(mf.list(6), mf.list(7)))))
		}
	}
	if len(poly) < 3 {
		return nil
	}
	return (&cad.Path{}).Polyline(poly, true)
}

// flatten samples a path into points.
func flatten(p *cad.Path) []cad.Point { return cad.Flatten(p, 32) }

// outline draws the outline of a composite curve that shows it.
func (r *sfcRenderer) outline(code int, path *cad.Path, m canvas.Matrix) {
	if code < 1 || code > len(r.m.curves) || path == nil {
		return
	}
	c := r.m.curves[code-1]
	if !c.visible {
		return
	}
	w := r.m.width(c.width)
	r.stroke(path, lineStyle{color: r.color(c.color), width: w, dash: r.m.pattern(c.typ, w)}, m)
}

func (r *sfcRenderer) hatchFeature(f *feature, m canvas.Matrix) {
	// the outline and holes are the last three arguments: out_id, the
	// number of holes and their list
	n := len(f.args)
	if n < 3 {
		return
	}
	out := int(f.num(n - 3))
	holes := capped(f.list(n-1), r.renderer)
	bound := r.boundary(out)
	if bound == nil {
		r.warn("hatchboundary", "hatches whose outline is missing are not drawn")
		return
	}
	area := (&cad.Path{}).Append(bound)
	for _, h := range holes {
		if hp := r.boundary(int(h)); hp != nil {
			area.Append(hp)
		}
	}
	switch f.name {
	case "externally_defined_hatch_feature":
		if name := f.str(1); name == "Area_control" {
			// a blank area: filled with the background
			if r.budget(area.Size()) {
				r.out.Fill(area.Transform(m), cad.Fill{Color: r.bg}, true)
			}
		} else {
			r.warn("hatch:"+name, "the predefined hatch %q is not drawn", name)
		}
	case "fill_area_style_colour_feature":
		if r.budget(area.Size()) {
			r.out.Fill(area.Transform(m), cad.Fill{Color: r.paint(r.color(f.int(1)))}, true)
		}
	case "fill_area_style_hatching_feature":
		var lines []hatchLine
		for i := range min(f.int(1), 4) {
			v := f.list(2 + i)
			if len(v) < 7 {
				continue
			}
			w := r.m.width(int(v[2]))
			lines = append(lines, hatchLine{style: lineStyle{color: r.color(int(v[0])), width: w, dash: r.m.pattern(int(v[1]), w)},
				x: v[3], y: v[4], spacing: v[5], angle: v[6]})
		}
		r.hatch(area, lines, m)
	case "fill_area_style_tiles_feature":
		r.warn("tiles", "hatches of tiled symbols are not drawn")
	}
	r.outline(out, bound, m)
	for _, h := range holes {
		r.outline(int(h), r.boundary(int(h)), m)
	}
}
