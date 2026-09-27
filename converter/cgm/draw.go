package cgm

import (
	"math"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
)

// vdcPerMM returns the VDC units of a millimetre on the page.
func (in *interp) vdcPerMM() float64 {
	if in.pic == nil || in.pic.k == 0 {
		return 1
	}
	return mm / in.pic.k
}

// sizeVDC returns a width or size in VDC; nominal is the nominal size in
// mm of scaled sizes.
func (in *interp) sizeVDC(s size, nominal float64) float64 {
	v := math.Abs(s.v)
	switch s.mode {
	case modeScaled:
		return v * nominal * in.vdcPerMM()
	case modeFractional:
		return v * in.pic.long
	case modeMM:
		return v * in.vdcPerMM()
	}
	return v
}

// setWidth sets the width of a pen: on paper for scaled and metric
// widths, in VDC for absolute and fractional ones. A width of 0 is the
// thinnest line.
func (in *interp) setWidth(pen *cad.Pen, w size) {
	v := math.Abs(w.v)
	switch w.mode {
	case modeScaled:
		pen.Width = v * nominalLine * mm
	case modeMM:
		pen.Width = v * mm
	default:
		pen.WorldWidth = in.sizeVDC(w, nominalLine)
	}
}

// standardDashes are the patterns of the line types 2 to 5 (dash, dot,
// dash-dot, dash-dot-dot) in mm on paper for lines up to 0.35 mm wide.
var standardDashes = [][]float64{
	2: {3, 1.5},
	3: {0, 1},
	4: {3, 1, 0, 1},
	5: {3, 1, 0, 1, 0, 1},
}

// pen returns the pen of a line or an edge.
func (in *interp) pen(typ int, w size, col colour, capStyle, join int, offset float64) cad.Pen {
	pen := cad.Pen{Color: in.colourOf(col), Cap: cad.CapRound, Join: cad.JoinRound}
	in.setWidth(&pen, w)
	switch capStyle {
	case 2, 5:
		pen.Cap = cad.CapButt
	case 4:
		pen.Cap = cad.CapSquare
	}
	switch join {
	case 2:
		pen.Join = cad.JoinMiter
	case 4:
		pen.Join = cad.JoinBevel
	}
	var dash []float64
	switch {
	case typ >= 2 && typ <= 5:
		// patterns grow with lines wider than 0.35 mm
		wmm := in.sizeVDC(w, nominalLine) / in.vdcPerMM()
		f := math.Max(1, wmm/0.35) * in.vdcPerMM()
		for _, v := range standardDashes[typ] {
			dash = append(dash, v*f)
		}
	case typ < 0:
		d, ok := in.st.lineTypes[typ]
		if !ok || len(d.dash) == 0 {
			in.warnOnce("linetype:"+itoa(typ), "line type %d is not defined; it is drawn solid", typ)
			break
		}
		total := 0.0
		for _, v := range d.dash {
			total += v
		}
		repeat := in.sizeVDC(d.repeat, nominalLine)
		if total <= 0 || repeat <= 0 {
			break
		}
		for i, v := range d.dash {
			if i%2 == 1 {
				v = -v // a gap
			}
			dash = append(dash, v*repeat/total)
		}
		dash, pen.DashOffset = cad.DashPattern(dash)
	case typ > 5:
		in.warnOnce("linetype:"+itoa(typ), "line type %d is drawn solid", typ)
	}
	if len(dash) > 0 {
		pen.Dash = dash
		if offset != 0 {
			total := 0.0
			for _, v := range dash {
				total += v
			}
			pen.DashOffset += offset * total
		}
	}
	return pen
}

// linePen returns the pen of lines, with the attributes that come from
// the line bundle.
func (in *interp) linePen() cad.Pen {
	a := &in.st.a
	typ, w, col := a.lineType, a.lineWidth, a.lineColour
	if r, ok := in.st.lineReps[a.lineBundle]; ok {
		if a.asf[asfLineType] {
			typ = r.typ
		}
		if a.asf[asfLineWidth] {
			w = r.width
		}
		if a.asf[asfLineColour] {
			col = r.colour
		}
	}
	return in.pen(typ, w, col, a.lineCap, a.lineJoin, a.lineTypeOffset)
}

// edgePen returns the pen of edges.
func (in *interp) edgePen() cad.Pen {
	a := &in.st.a
	typ, w, col := a.edgeType, a.edgeWidth, a.edgeColour
	if r, ok := in.st.edgeRs[a.edgeBundle]; ok {
		if a.asf[asfEdgeType] {
			typ = r.typ
		}
		if a.asf[asfEdgeWidth] {
			w = r.width
		}
		if a.asf[asfEdgeColour] {
			col = r.colour
		}
	}
	return in.pen(typ, w, col, a.edgeCap, a.edgeJoin, a.edgeTypeOffset)
}

// fillAttrs returns the interior style, fill colour, hatch and pattern
// indexes, with those that come from the fill bundle.
func (in *interp) fillAttrs() (style int, col colour, hatch, pattern int) {
	a := &in.st.a
	style, col, hatch, pattern = a.interior, a.fillColour, a.hatch, a.pattern
	if r, ok := in.st.fillReps[a.fillBundle]; ok {
		if a.asf[asfInteriorStyle] {
			style = r.style
		}
		if a.asf[asfFillColour] {
			col = r.colour
		}
		if a.asf[asfHatchIndex] {
			hatch = r.hatch
		}
		if a.asf[asfPatternIndex] {
			pattern = r.pattern
		}
	}
	return
}

// line draws an open path with the line attributes, or adds it to the
// figure or compound line being built; connect joins it to what comes
// before in those.
func (in *interp) line(path *cad.Path, connect bool) {
	if path.Empty() {
		return
	}
	switch {
	case in.fig != nil:
		if connect {
			in.fig.path.Continue(path)
		} else {
			in.fig.path.Append(path)
		}
	case in.compound != nil:
		if connect {
			in.compound.Continue(path)
		} else {
			in.compound.Append(path)
		}
	default:
		in.out.Stroke(path, in.linePen())
	}
}

// area fills a closed path with the fill attributes and draws its edges
// (the path itself when edges is nil) when edges are visible, or adds it
// to the figure being built.
func (in *interp) area(fill, edges *cad.Path) {
	if fill.Empty() {
		return
	}
	if in.fig != nil {
		in.fig.path.Append(fill)
		return
	}
	if !in.drawing() {
		return
	}
	style, col, hatch, pattern := in.fillAttrs()
	c := in.colourOf(col)
	switch style {
	case styleHollow:
		in.out.Stroke(fill, cad.Pen{Color: c, Cap: cad.CapRound, Join: cad.JoinRound})
	case styleSolid, styleGeometric:
		in.out.Fill(fill, cad.Fill{Color: c}, true)
	case styleHatch:
		in.hatchFill(fill, hatch, c)
	case stylePattern:
		in.patternFill(fill, pattern, c)
	case styleInterpolated:
		if !in.gradientFill(fill) {
			in.out.Fill(fill, cad.Fill{Color: c}, true)
		}
	}
	if in.st.a.edgeVisible {
		if edges == nil {
			edges = fill
		}
		if !edges.Empty() {
			in.out.Stroke(edges, in.edgePen())
		}
	}
}

// fillRef returns the fill reference point.
func (in *interp) fillRef() cad.Point {
	if in.st.a.fillRefSet {
		return in.st.a.fillRef
	}
	return in.extent()[0]
}

// hatchFill draws the hatch lines of a hatch index clipped by the path:
// the standard ones (horizontal, vertical, positive and negative slope,
// and the two crosshatches) 2 mm apart on paper, or those of HATCH STYLE
// DEFINITION.
func (in *interp) hatchFill(path *cad.Path, idx int, c bdf.Color) {
	ref := in.fillRef()
	sp := hatchSpacing * in.vdcPerMM()
	var lines []cad.PatternLine
	add := func(deg float64) {
		s, co := math.Sincos(deg * math.Pi / 180)
		lines = append(lines, cad.PatternLine{Angle: deg * math.Pi / 180, Base: ref, Offset: cad.Point{X: -s * sp, Y: co * sp}})
	}
	switch idx {
	case 1:
		add(0)
	case 2:
		add(90)
	case 3:
		add(45)
	case 4:
		add(135)
	case 5:
		add(0)
		add(90)
	case 6:
		add(45)
		add(135)
	default:
		d, ok := in.st.hatchStyles[idx]
		if !ok {
			in.warnOnce("hatch:"+itoa(idx), "hatch index %d is not defined; it is drawn as horizontal lines", idx)
			add(0)
			break
		}
		lines = in.hatchLines(d, ref)
	}
	if !in.st.transparency {
		in.out.Fill(path, cad.Fill{Color: in.colourOf(in.st.aux)}, true)
	}
	// the lines across the boundary's box, at most
	if b := path.Bounds(); b.Valid() {
		for _, l := range lines {
			if n := l.Offset.Len(); n > 0 {
				in.points += int(min((b.W()+b.H())/n, cad.MaxHatchLines)) * 2
			}
		}
	}
	if len(lines) == 0 || !in.out.Hatch(path, lines, cad.Pen{Color: c}) {
		in.warnOnce("hatch lines", "hatching of too many lines is filled with a lighter colour")
		in.out.Fill(path, cad.Fill{Color: blend(c, in.st.background, 0.5)}, true)
	}
}

// hatchLines returns the lines of a HATCH STYLE DEFINITION: lines along
// the first direction, repeated every duty cycle along the second, at the
// positions the gap widths divide the cycle into (and the same turned for
// a crosshatch).
func (in *interp) hatchLines(d hatchDef, ref cad.Point) []cad.PatternLine {
	cycle := in.sizeVDC(d.cycle, nominalLine)
	if cycle <= 0 || d.dir1.Len() == 0 {
		return nil
	}
	total := 0.0
	for _, g := range d.gaps {
		total += g
	}
	family := func(u, v cad.Point) []cad.PatternLine {
		u = u.Mul(1 / u.Len())
		if v.Len() == 0 || math.Abs(u.X*v.Y-u.Y*v.X) < 1e-9*v.Len() {
			v = cad.Point{X: -u.Y, Y: u.X}
		}
		v = v.Mul(1 / v.Len())
		ang := math.Atan2(u.Y, u.X)
		var out []cad.PatternLine
		pos := 0.0
		n := max(len(d.gaps), 1)
		for i := 0; i < n; i++ {
			pl := cad.PatternLine{Angle: ang, Base: ref.Add(v.Mul(pos)), Offset: v.Mul(cycle)}
			if i < len(d.types) {
				pen := in.pen(d.types[i], size{0, modeAbsolute}, colour{}, 1, 1, 0)
				for j, x := range pen.Dash {
					if j%2 == 1 {
						x = -x
					}
					pl.Dash = append(pl.Dash, x)
				}
			}
			out = append(out, pl)
			if i < len(d.gaps) && total > 0 {
				pos += cycle * d.gaps[i] / total
			}
		}
		return out
	}
	lines := family(d.dir1, d.dir2)
	if d.cross {
		lines = append(lines, family(d.dir2, d.dir1)...)
	}
	return lines
}

// blend mixes c with bg: t = 0 is c.
func blend(c, bg bdf.Color, t float64) bdf.Color {
	mix := func(sh uint) uint8 {
		a, b := float64(uint8(c>>sh)), float64(uint8(bg>>sh))
		return uint8(math.Round(a + (b-a)*t))
	}
	return bdf.RGB(mix(24), mix(16), mix(8))
}

// gradientFill fills a path with an INTERPOLATED INTERIOR: parallel as a
// linear gradient from the fill reference point along its vector,
// elliptical as a radial one around it.
func (in *interp) gradientFill(path *cad.Path) bool {
	ii := in.st.a.interp
	if ii == nil || len(ii.colours) < 2 || len(ii.geom) < 2 {
		if ii != nil && ii.style == 3 {
			in.warnOnce("triangular", "triangular interpolated interiors are filled with the fill colour")
		}
		return false
	}
	stops := make([]bdf.Stop, len(ii.colours))
	lo, hi := 0.0, 1.0
	if len(ii.stages) == len(ii.colours) {
		lo, hi = ii.stages[0], ii.stages[len(ii.stages)-1]
	}
	for i, c := range ii.colours {
		t := float64(i) / float64(len(ii.colours)-1)
		if len(ii.stages) == len(ii.colours) && hi > lo {
			t = (ii.stages[i] - lo) / (hi - lo)
		}
		stops[i] = bdf.Stop{Offset: float32(math.Min(math.Max(t, 0), 1)), Color: in.colourOf(c)}
	}
	ref := in.fillRef()
	v := cad.Point{X: ii.geom[0], Y: ii.geom[1]}
	if v.Len() == 0 {
		return false
	}
	g := &cad.Gradient{P0: ref, P1: ref.Add(v), Stops: stops}
	if ii.style == 2 {
		g = &cad.Gradient{Radial: true, P0: ref, P1: ref, R1: v.Len(), Stops: stops}
	}
	in.out.Fill(path, cad.Fill{Gradient: g}, true)
	return true
}

// readPoints reads the points left.
func (in *interp) readPoints(p params) []cad.Point {
	var pts []cad.Point
	for p.more() {
		pts = append(pts, p.point())
	}
	in.points += len(pts)
	return pts
}

func (in *interp) primitive(c code, p params) {
	switch c {
	case eText, eRestrictedText, eAppendText:
		in.ensureBody()
		in.textElement(c, p)
		return
	case eTile, eBitonalTile:
		in.ensureBody()
		in.tile(c, p)
		return
	}
	in.flushText()
	in.ensureBody()
	if !in.drawing() {
		return
	}
	switch c {
	case ePolyline:
		pts := in.readPoints(p)
		if len(pts) >= 2 {
			in.line((&cad.Path{}).Polyline(pts, false), true)
		}
	case eDisjointPolyline:
		pts := in.readPoints(p)
		path := &cad.Path{}
		for i := 0; i+1 < len(pts); i += 2 {
			path.MoveTo(pts[i].X, pts[i].Y).LineTo(pts[i+1].X, pts[i+1].Y)
		}
		in.line(path, false)
	case ePolymarker:
		in.markers(in.readPoints(p))
	case ePolygon:
		pts := in.readPoints(p)
		if len(pts) >= 2 {
			in.area((&cad.Path{}).Polyline(pts, true), nil)
		}
	case ePolygonSet:
		in.polygonSet(p)
	case eCellArray:
		in.cellArray(p)
	case eGDP:
		id := p.int()
		in.warnOnce("gdp:"+itoa(id), "generalized drawing primitives (%d) are not drawn", id)
	case eRectangle:
		a, b := p.point(), p.point()
		in.area((&cad.Path{}).Polyline([]cad.Point{a, {X: b.X, Y: a.Y}, b, {X: a.X, Y: b.Y}}, true), nil)
	case eCircle:
		c, r := p.point(), math.Abs(p.vdc())
		if r > 0 {
			in.area((&cad.Path{}).Circle(c, r), nil)
		}
	case eArc3Point, eArc3PointClose:
		a, b, e := p.point(), p.point(), p.point()
		path, centre := arc3(a, b, e)
		if c == eArc3PointClose {
			in.closeArc(path, centre, p.enum("PIE CHORD"))
		} else {
			in.line(path, true)
		}
	case eArcCentre, eArcCentreClose, eArcCentreReversed:
		centre := p.point()
		s := cad.Point{X: p.vdc(), Y: p.vdc()}
		e := cad.Point{X: p.vdc(), Y: p.vdc()}
		r := math.Abs(p.vdc())
		a0, a1 := math.Atan2(s.Y, s.X), math.Atan2(e.Y, e.X)
		if c == eArcCentreReversed {
			if a1 >= a0 {
				a1 -= 2 * math.Pi
			}
		} else if a1 <= a0 {
			a1 += 2 * math.Pi
		}
		path := (&cad.Path{}).Arc(centre, r, a0, a1)
		if c == eArcCentreClose {
			in.closeArc(path, centre, p.enum("PIE CHORD"))
		} else {
			in.line(path, true)
		}
	case eEllipse:
		centre, d1, d2 := p.point(), p.point(), p.point()
		path := (&cad.Path{}).EllipseArcAxes(centre, d1.Sub(centre), d2.Sub(centre), 0, 2*math.Pi).Close()
		in.area(path, nil)
	case eEllipticalArc, eEllipticalArcClos:
		centre, d1, d2 := p.point(), p.point(), p.point()
		s := cad.Point{X: p.vdc(), Y: p.vdc()}
		e := cad.Point{X: p.vdc(), Y: p.vdc()}
		a, b := d1.Sub(centre), d2.Sub(centre)
		t0, ok0 := ellipseParam(a, b, s)
		t1, ok1 := ellipseParam(a, b, e)
		if !ok0 || !ok1 {
			return
		}
		if t1 <= t0 {
			t1 += 2 * math.Pi
		}
		path := (&cad.Path{}).EllipseArcAxes(centre, a, b, t0, t1)
		if c == eEllipticalArcClos {
			in.closeArc(path, centre, p.enum("PIE CHORD"))
		} else {
			in.line(path, true)
		}
	case eConnectingEdge:
		// figures join their parts anyway
	case eHyperbolicArc:
		centre, tr, cr := p.point(), p.point(), p.point()
		s := cad.Point{X: p.vdc(), Y: p.vdc()}
		e := cad.Point{X: p.vdc(), Y: p.vdc()}
		in.line(hyperbola(centre, tr.Sub(centre), cr.Sub(centre), s, e), true)
	case eParabolicArc:
		t, a, b := p.point(), p.point(), p.point()
		c1 := a.Add(t.Sub(a).Mul(2.0 / 3))
		c2 := b.Add(t.Sub(b).Mul(2.0 / 3))
		in.line((&cad.Path{}).MoveTo(a.X, a.Y).CubicTo(c1, c2, b), true)
	case eNUBSpline, eNURBSpline:
		in.spline(p, c == eNURBSpline)
	case ePolybezier:
		continuous := p.index() == 2
		pts := in.readPoints(p)
		path := &cad.Path{}
		for i := 0; i+3 < len(pts); {
			if i == 0 || !continuous {
				path.MoveTo(pts[i].X, pts[i].Y)
			}
			path.CubicTo(pts[i+1], pts[i+2], pts[i+3])
			i += 3
			if !continuous {
				i++
			}
		}
		in.line(path, true)
	case ePolysymbol:
		in.warnOnce("symbol", "symbols (POLYSYMBOL) are not drawn")
	}
}

// closeArc closes an arc as a pie (through the centre) or a chord and
// fills it.
func (in *interp) closeArc(path *cad.Path, centre cad.Point, typ int) {
	if path.Empty() {
		return
	}
	if typ == 0 {
		path.LineTo(centre.X, centre.Y)
	}
	in.area(path.Close(), nil)
}

// arc3 returns the circular arc from a through b to c, and its centre; a
// straight line when the points are in a line, a full circle when a and c
// coincide.
func arc3(a, b, c cad.Point) (*cad.Path, cad.Point) {
	path := &cad.Path{}
	if a == c {
		if a == b {
			return path, a
		}
		centre := a.Add(b).Mul(0.5)
		return path.Arc(centre, b.Sub(a).Len()/2, math.Atan2(a.Y-centre.Y, a.X-centre.X), math.Atan2(a.Y-centre.Y, a.X-centre.X)+2*math.Pi), centre
	}
	// relative to a, for precision
	bx, by := b.X-a.X, b.Y-a.Y
	cx, cy := c.X-a.X, c.Y-a.Y
	d := 2 * (bx*cy - by*cx)
	scale := math.Max(math.Hypot(bx, by), math.Hypot(cx, cy))
	if math.Abs(d) <= 1e-12*scale*scale {
		return path.MoveTo(a.X, a.Y).LineTo(c.X, c.Y), a.Add(c).Mul(0.5)
	}
	b2, c2 := bx*bx+by*by, cx*cx+cy*cy
	ux := (cy*b2 - by*c2) / d
	uy := (bx*c2 - cx*b2) / d
	centre := cad.Point{X: a.X + ux, Y: a.Y + uy}
	r := math.Hypot(ux, uy)
	a0 := math.Atan2(a.Y-centre.Y, a.X-centre.X)
	am := math.Atan2(b.Y-centre.Y, b.X-centre.X)
	a1 := math.Atan2(c.Y-centre.Y, c.X-centre.X)
	norm := func(v float64) float64 {
		v = math.Mod(v, 2*math.Pi)
		if v < 0 {
			v += 2 * math.Pi
		}
		return v
	}
	sm, se := norm(am-a0), norm(a1-a0)
	if sm < se {
		return path.Arc(centre, r, a0, a0+se), centre
	}
	return path.Arc(centre, r, a0, a0-(2*math.Pi-se)), centre
}

// ellipseParam returns the parameter t of the point of the ellipse
// c + a·cos t + b·sin t in the direction v from its centre.
func ellipseParam(a, b, v cad.Point) (float64, bool) {
	det := a.X*b.Y - a.Y*b.X
	if det == 0 || !isFinite(det) {
		return 0, false
	}
	alpha := (v.X*b.Y - v.Y*b.X) / det
	beta := (a.X*v.Y - a.Y*v.X) / det
	return math.Atan2(beta, alpha), true
}

// hyperbola returns the arc of the hyperbola c + a·cosh t + b·sinh t (the
// branch through c + a) between the directions s and e from its centre.
func hyperbola(c, a, b, s, e cad.Point) *cad.Path {
	path := &cad.Path{}
	det := a.X*b.Y - a.Y*b.X
	if det == 0 || !isFinite(det) {
		return path
	}
	param := func(v cad.Point) float64 {
		alpha := (v.X*b.Y - v.Y*b.X) / det
		beta := (a.X*v.Y - a.Y*v.X) / det
		r := beta / math.Abs(alpha)
		if alpha == 0 || math.Abs(r) >= 1 {
			r = math.Copysign(0.999999, beta)
		}
		return math.Atanh(r)
	}
	t0, t1 := param(s), param(e)
	const n = 64
	for i := 0; i <= n; i++ {
		t := t0 + (t1-t0)*float64(i)/n
		q := c.Add(a.Mul(math.Cosh(t))).Add(b.Mul(math.Sinh(t)))
		if i == 0 {
			path.MoveTo(q.X, q.Y)
		} else {
			path.LineTo(q.X, q.Y)
		}
	}
	return path
}

// spline draws a NON-UNIFORM B-SPLINE or NON-UNIFORM RATIONAL B-SPLINE:
// its order, control points, knots (weights) and parameter range.
func (in *interp) spline(p params, rational bool) {
	order, n := p.int(), p.int()
	if order < 2 || n < order || n > 1<<20 {
		in.warnOnce("spline", "B-splines with invalid orders or counts are not drawn")
		return
	}
	ctrl := make([]cad.Point, n)
	for i := range ctrl {
		ctrl[i] = p.point()
	}
	knots := make([]float64, n+order)
	for i := range knots {
		knots[i] = p.real()
	}
	var weights []float64
	if rational {
		weights = make([]float64, n)
		for i := range weights {
			weights[i] = p.real()
		}
	}
	t0, t1 := knots[order-1], knots[n]
	if p.more() {
		t0 = p.real()
		t1 = p.real()
	}
	in.points += n
	in.line((&cad.Path{}).BSplineRange(order-1, knots, ctrl, weights, t0, t1), true)
}

// polygonSet draws a POLYGON SET: points with the flags of the edges that
// leave them (invisible, visible, and close invisible or visible, which
// close the polygon and start another).
func (in *interp) polygonSet(p params) {
	var pts []cad.Point
	var flags []int
	for p.more() {
		pts = append(pts, p.point())
		flags = append(flags, p.enum("INVIS VIS CLOSEINVIS CLOSEVIS"))
	}
	in.points += len(pts)
	fill, edges := &cad.Path{}, &cad.Path{}
	first := 0
	for i, q := range pts {
		closing := flags[i] >= 2 || i == len(pts)-1
		to := pts[first]
		if !closing {
			to = pts[i+1]
		}
		if i == first {
			fill.MoveTo(q.X, q.Y)
		} else {
			fill.LineTo(q.X, q.Y)
		}
		if flags[i]&1 == 1 {
			if !edges.Open() || edges.Current() != q {
				edges.MoveTo(q.X, q.Y)
			}
			edges.LineTo(to.X, to.Y)
		}
		if closing {
			fill.Close()
			first = i + 1
		}
	}
	in.area(fill, edges)
}

// markers draws POLYMARKER: types 1 to 5 are a dot, a plus sign, an
// asterisk, a circle and a cross.
func (in *interp) markers(pts []cad.Point) {
	a := &in.st.a
	typ, sz, col := a.markerType, a.markerSize, a.markerColour
	if r, ok := in.st.markerReps[a.markerBundle]; ok {
		if a.asf[asfMarkerType] {
			typ = r.typ
		}
		if a.asf[asfMarkerSize] {
			sz = r.width
		}
		if a.asf[asfMarkerColour] {
			col = r.colour
		}
	}
	s := in.sizeVDC(sz, nominalMarker) / 2
	c := in.colourOf(col)
	pen := cad.Pen{Color: c, Cap: cad.CapButt}
	if typ < 1 || typ > 5 {
		in.warnOnce("marker:"+itoa(typ), "marker type %d is drawn as an asterisk", typ)
		typ = 3
	}
	path := &cad.Path{}
	dots := &cad.Path{}
	d := s * math.Sqrt2 / 2
	for _, q := range pts {
		switch typ {
		case 1:
			dots.Circle(q, 0.25*in.vdcPerMM())
		case 2:
			path.MoveTo(q.X-s, q.Y).LineTo(q.X+s, q.Y).MoveTo(q.X, q.Y-s).LineTo(q.X, q.Y+s)
		case 3:
			path.MoveTo(q.X-s, q.Y).LineTo(q.X+s, q.Y).MoveTo(q.X, q.Y-s).LineTo(q.X, q.Y+s)
			path.MoveTo(q.X-d, q.Y-d).LineTo(q.X+d, q.Y+d).MoveTo(q.X-d, q.Y+d).LineTo(q.X+d, q.Y-d)
		case 4:
			path.Circle(q, s)
		case 5:
			path.MoveTo(q.X-d, q.Y-d).LineTo(q.X+d, q.Y+d).MoveTo(q.X-d, q.Y+d).LineTo(q.X+d, q.Y-d)
		}
	}
	in.out.Fill(dots, cad.Fill{Color: c}, false)
	in.out.Stroke(path, pen)
}
