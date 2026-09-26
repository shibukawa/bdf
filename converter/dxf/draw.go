package dxf

import (
	"math"
	"slices"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// props are the resolved properties of an entity, and what BYBLOCK means
// for the entities of a block it inserts.
type props struct {
	color bdf.Color
	ltype string
	lw    int // lineweight in hundredths of a mm; negative: the default
}

// ctx is where entities are being drawn.
type ctx struct {
	m   canvas.Matrix // the coordinates of the entity (WCS of its block) to the drawing
	out *cad.Drawing
	// dark selects the palette of a dark background
	dark bool
	// plot leaves out the layers that are not plotted (layouts)
	plot bool
	// insLayer is the layer of the INSERT whose block is drawn: entities
	// on layer 0 take its properties
	insLayer *layer
	byBlock  props
	depth    int
	// frozen are the layers frozen in the viewport (upper-case names)
	frozen map[string]bool
	// ltFactor multiplies line type lengths (1/the viewport scale when
	// $PSLTSCALE sets them in paper units)
	ltFactor float64
	// deferred are drawn once the extent of the drawing is known
	// (construction lines, points sized relative to the view)
	deferred *[]func(view cad.Rect)
	// blocks are the blocks being drawn, outermost first (upper-case
	// names): a block that inserts itself is not drawn again
	blocks []string
}

const maxDepth = 24

// maxEntities bounds the entities drawn for one view, blocks included.
const maxEntities = 5_000_000

// ccwSpan reduces the angles of an arc turning counter-clockwise from a0
// to a1 (radians) to a span of at most a turn; ok is false when an angle
// is not a finite number.
func ccwSpan(a0, a1 float64) (float64, float64, bool) {
	if math.IsNaN(a0) || math.IsInf(a0, 0) || math.IsNaN(a1) || math.IsInf(a1, 0) {
		return 0, 0, false
	}
	a0, a1 = math.Mod(a0, 2*math.Pi), math.Mod(a1, 2*math.Pi)
	if a1 <= a0 {
		a1 += 2 * math.Pi
	}
	return a0, a1, true
}

var defaultLayer = &layer{name: "0", color: 7, trueColor: -1, ltype: "CONTINUOUS", lineweight: -3, alpha: 1}

func (c *converter) layer(name string) *layer {
	if name == "" {
		name = "0"
	}
	if l := c.d.layers[key(name)]; l != nil {
		return l
	}
	return defaultLayer
}

// visible returns the layer whose properties an entity takes (the layer of
// the INSERT for entities on layer 0 of a block), and whether it is shown.
func (c *converter) visible(e *entity, x *ctx) (*layer, bool) {
	if e.int(60, 0) == 1 {
		return nil, false
	}
	// (a layer missing from the table has the properties of layer 0 but
	// keeps its name)
	name := key(e.str(8))
	if name == "" {
		name = "0"
	}
	own := c.layer(name)
	if own.frozen || x.frozen[name] {
		return nil, false
	}
	eff := own
	if name == "0" && x.insLayer != nil {
		eff = x.insLayer
	}
	if eff.color < 0 && e.typ != "INSERT" {
		// off; the entities of an inserted block on other layers stay
		return nil, false
	}
	if name == "DEFPOINTS" || x.plot && eff.noPlot {
		// the definition points of dimensions are never shown
		return nil, false
	}
	return eff, true
}

func (c *converter) aci(i int, dark bool) bdf.Color {
	if i < 0 {
		i = -i
	}
	if i < 1 || i > 255 {
		i = 7
	}
	v := aciLight[i]
	if dark {
		v = aciDark[i]
	}
	return bdf.RGB(uint8(v>>16), uint8(v>>8), uint8(v))
}

func rgb(v int) bdf.Color { return bdf.RGB(uint8(v>>16), uint8(v>>8), uint8(v)) }

func withAlpha(c bdf.Color, a float64) bdf.Color {
	if a >= 1 {
		return c
	}
	return c&^0xff | bdf.Color(uint8(math.Round(math.Max(a, 0)*255)))
}

// props resolves the color, line type and lineweight of an entity.
func (c *converter) props(e *entity, x *ctx, l *layer) props {
	var p props
	ci := e.int(62, 256)
	switch {
	case e.has(420) && (!e.has(62) || ci != 0 && ci != 256):
		p.color = rgb(e.int(420, 0))
	case ci == 0:
		p.color = x.byBlock.color
	case ci == 256 || ci == 257:
		if l.trueColor >= 0 {
			p.color = rgb(int(l.trueColor))
		} else {
			p.color = c.aci(l.color, x.dark)
		}
	default:
		p.color = c.aci(ci, x.dark)
	}
	alpha := l.alpha
	if e.has(440) {
		v := e.int(440, 0)
		switch {
		case v&0x01000000 != 0:
			alpha = float64(uint8(x.byBlock.color)) / 255
		case v&0x02000000 != 0:
			alpha = transparency(v)
		}
	}
	p.color = withAlpha(p.color, alpha)
	switch lt := strings.ToUpper(e.str(6)); lt {
	case "", "BYLAYER":
		p.ltype = l.ltype
	case "BYBLOCK":
		p.ltype = x.byBlock.ltype
	default:
		p.ltype = lt
	}
	switch lw := e.int(370, -1); lw {
	case -1:
		p.lw = l.lineweight
		if p.lw == -2 {
			p.lw = x.byBlock.lw
		}
	case -2:
		p.lw = x.byBlock.lw
	default:
		p.lw = lw
	}
	return p
}

// pen returns the pen of an entity drawn through m.
func (c *converter) pen(e *entity, x *ctx, p props, m canvas.Matrix) cad.Pen {
	lw := p.lw
	if lw < 0 {
		lw = 25 // AutoCAD's default lineweight, 0.25 mm
	}
	pen := cad.Pen{Color: p.color, Width: float64(lw) / 100 * 72 / 25.4, Cap: cad.CapRound, Join: cad.JoinRound}
	if lt := c.d.ltypes[key(p.ltype)]; lt != nil && len(lt.elems) > 0 {
		s := e.num(48, 1) * c.ltscale * cad.Scale(m) * x.ltFactor
		elems := make([]float64, len(lt.elems))
		for i, v := range lt.elems {
			elems[i] = v * s
		}
		pen.Dash, pen.DashOffset = cad.DashPattern(elems)
	}
	return pen
}

// ocs returns the transform of an entity's object coordinate system into
// the XY plane of the world (the arbitrary axis algorithm, then a view
// from the top); elevation is the entity's z in the OCS.
func ocs(e *entity, elevation float64) canvas.Matrix {
	n := e.vec3(210, [3]float64{0, 0, 1})
	return ocsMatrix(n, elevation)
}

func ocsMatrix(n [3]float64, elevation float64) canvas.Matrix {
	l := math.Sqrt(n[0]*n[0] + n[1]*n[1] + n[2]*n[2])
	if l == 0 || (n[0] == 0 && n[1] == 0 && n[2] > 0) {
		return canvas.Identity
	}
	n = [3]float64{n[0] / l, n[1] / l, n[2] / l}
	var ax [3]float64
	if math.Abs(n[0]) < 1.0/64 && math.Abs(n[1]) < 1.0/64 {
		ax = cross([3]float64{0, 1, 0}, n)
	} else {
		ax = cross([3]float64{0, 0, 1}, n)
	}
	ax = norm(ax)
	ay := norm(cross(n, ax))
	return canvas.Matrix{ax[0], ax[1], ay[0], ay[1], n[0] * elevation, n[1] * elevation}
}

func cross(a, b [3]float64) [3]float64 {
	return [3]float64{a[1]*b[2] - a[2]*b[1], a[2]*b[0] - a[0]*b[2], a[0]*b[1] - a[1]*b[0]}
}

func norm(v [3]float64) [3]float64 {
	l := math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
	if l == 0 {
		return v
	}
	return [3]float64{v[0] / l, v[1] / l, v[2] / l}
}

func pt2(v [3]float64) cad.Point { return cad.Point{X: v[0], Y: v[1]} }

// entity draws an entity.
func (c *converter) entity(e *entity, x *ctx) {
	if x.depth > maxDepth {
		c.warnOnce("depth", "blocks nested deeper than %d levels are not drawn", maxDepth)
		return
	}
	if c.count++; c.count > maxEntities {
		c.warnOnce("budget", "only the first %d entities of a view are drawn", maxEntities)
		return
	}
	l, ok := c.visible(e, x)
	if !ok {
		return
	}
	p := c.props(e, x, l)
	switch e.typ {
	case "LINE":
		a, b := e.vec3(10, [3]float64{}), e.vec3(11, [3]float64{})
		path := (&cad.Path{}).MoveTo(a[0], a[1]).LineTo(b[0], b[1])
		x.out.Stroke(path.Transform(x.m), c.pen(e, x, p, x.m))
	case "XLINE", "RAY":
		c.construction(e, x, p)
	case "POINT":
		c.point(e, x, p)
	case "CIRCLE", "ARC":
		ctr := e.vec3(10, [3]float64{})
		m := x.m.Mul(ocs(e, ctr[2]))
		r := e.num(40, 0)
		path := &cad.Path{}
		if e.typ == "CIRCLE" {
			path.Circle(pt2(ctr), r)
		} else {
			a0, a1, ok := ccwSpan(e.num(50, 0)*math.Pi/180, e.num(51, 360)*math.Pi/180)
			if !ok {
				return
			}
			path.Arc(pt2(ctr), r, a0, a1)
		}
		x.out.Stroke(path.Transform(m), c.pen(e, x, p, m))
	case "ELLIPSE":
		c.ellipse(e, x, p)
	case "LWPOLYLINE":
		vs, closed, _ := lwVertices(e)
		m := x.m.Mul(ocs(e, e.num(38, 0)))
		c.polyline(e, x, p, m, vs, closed)
	case "POLYLINE":
		c.heavyPolyline(e, x, p)
	case "SPLINE", "HELIX":
		c.spline(e, x, p)
	case "SOLID", "TRACE":
		pts := []cad.Point{e.pt(10), e.pt(11), e.pt(13), e.pt(12)}
		if !e.has(13) {
			pts = pts[:2]
			pts = append(pts, e.pt(12))
		}
		m := x.m.Mul(ocs(e, e.num(30, 0)))
		path := (&cad.Path{}).Polyline(pts, true).Transform(m)
		x.out.Fill(path, cad.Fill{Color: p.color}, false)
	case "3DFACE":
		c.face(e, x, p)
	case "TEXT", "ATTRIB":
		c.text(e, x, p)
	case "ATTDEF":
		// outside a block the definition shows its tag; in the inserts of
		// its block only a constant attribute shows (its value)
		if f := e.int(70, 0); x.depth == 0 && len(x.blocks) == 0 || f&2 != 0 && f&1 == 0 {
			c.text(e, x, p)
		}
	case "MTEXT":
		c.mtext(e, x, p)
	case "INSERT":
		c.insert(e, x, p, l)
	case "DIMENSION", "ARC_DIMENSION", "LARGE_RADIAL_DIMENSION":
		c.dimension(e, x, p, l)
	case "ACAD_TABLE":
		c.table(e, x, p, l)
	case "LEADER":
		c.leader(e, x, p)
	case "MULTILEADER", "MLEADER":
		c.multileader(e, x, p)
	case "HATCH", "MPOLYGON":
		c.hatch(e, x, p)
	case "WIPEOUT":
		c.wipeout(e, x)
	case "MLINE":
		c.mline(e, x, p)
	case "MESH":
		c.mesh(e, x, p)
	case "VIEWPORT":
		// the border; the layout draws what it shows
		if e.int(69, 0) != 1 && x.depth == 0 {
			ctr := e.pt(10)
			w, h := e.num(40, 0)/2, e.num(41, 0)/2
			if w > 0 && h > 0 {
				path := (&cad.Path{}).Polyline([]cad.Point{{X: ctr.X - w, Y: ctr.Y - h}, {X: ctr.X + w, Y: ctr.Y - h},
					{X: ctr.X + w, Y: ctr.Y + h}, {X: ctr.X - w, Y: ctr.Y + h}}, true)
				x.out.Stroke(path.Transform(x.m), c.pen(e, x, p, x.m))
			}
		}
	case "IMAGE":
		c.warnOnce("IMAGE", "raster images (IMAGE) refer to files outside the drawing and are not drawn")
	case "OLE2FRAME", "OLEFRAME":
		c.warnOnce(e.typ, "embedded OLE objects are not drawn")
	case "3DSOLID", "BODY", "REGION", "SURFACE", "EXTRUDEDSURFACE", "LOFTEDSURFACE", "REVOLVEDSURFACE", "SWEPTSURFACE", "PLANESURFACE", "NURBSSURFACE":
		c.warnOnce(e.typ, "ACIS solids, regions and surfaces are not drawn")
	default:
		c.warnOnce("type:"+e.typ, "entities of type %s are not drawn", e.typ)
	}
}

func (c *converter) construction(e *entity, x *ctx, p props) {
	if x.deferred == nil {
		return
	}
	a, d := e.vec3(10, [3]float64{}), e.vec3(11, [3]float64{1, 0, 0})
	pen := c.pen(e, x, p, x.m)
	m, out := x.m, x.out
	ray := e.typ == "RAY"
	*x.deferred = append(*x.deferred, func(view cad.Rect) {
		if !view.Valid() {
			return
		}
		// clip the line to the view (in drawing coordinates)
		pa := cad.Apply(m, pt2(a))
		pd := cad.Apply(m, pt2(a).Add(pt2(d))).Sub(pa)
		if pd.Len() == 0 {
			return
		}
		t0, t1 := math.Inf(-1), math.Inf(1)
		if ray {
			t0 = 0
		}
		for _, ax := range [2][3]float64{{pa.X, pd.X, 0}, {pa.Y, pd.Y, 1}} {
			lo, hi := view.Min.X, view.Max.X
			if ax[2] == 1 {
				lo, hi = view.Min.Y, view.Max.Y
			}
			if ax[1] == 0 {
				if ax[0] < lo || ax[0] > hi {
					return
				}
				continue
			}
			ta, tb := (lo-ax[0])/ax[1], (hi-ax[0])/ax[1]
			if ta > tb {
				ta, tb = tb, ta
			}
			t0, t1 = max(t0, ta), min(t1, tb)
		}
		if t0 >= t1 {
			return
		}
		q0, q1 := pa.Add(pd.Mul(t0)), pa.Add(pd.Mul(t1))
		out.Stroke((&cad.Path{}).MoveTo(q0.X, q0.Y).LineTo(q1.X, q1.Y), pen)
	})
}

func (c *converter) point(e *entity, x *ctx, p props) {
	mode := int(c.d.hnum("$PDMODE", 0))
	if mode&31 == 1 && mode&96 == 0 {
		return
	}
	loc := e.vec3(10, [3]float64{})
	pen := c.pen(e, x, p, x.m)
	pen.Dash = nil
	size := c.d.hnum("$PDSIZE", 0)
	m, out := x.m, x.out
	draw := func(sz float64) {
		path := &cad.Path{}
		cx, cy := loc[0], loc[1]
		h := sz / 2
		switch mode & 31 {
		case 0:
			path.MoveTo(cx, cy).LineTo(cx, cy)
		case 2:
			path.MoveTo(cx-h, cy).LineTo(cx+h, cy).MoveTo(cx, cy-h).LineTo(cx, cy+h)
		case 3:
			path.MoveTo(cx-h, cy-h).LineTo(cx+h, cy+h).MoveTo(cx-h, cy+h).LineTo(cx+h, cy-h)
		case 4:
			path.MoveTo(cx, cy).LineTo(cx, cy+h)
		}
		if mode&32 != 0 {
			path.Circle(cad.Point{X: cx, Y: cy}, h)
		}
		if mode&64 != 0 {
			path.Polyline([]cad.Point{{X: cx - h, Y: cy - h}, {X: cx + h, Y: cy - h}, {X: cx + h, Y: cy + h}, {X: cx - h, Y: cy + h}}, true)
		}
		out.Stroke(path.Transform(m), pen)
	}
	if size > 0 || mode&31 == 0 && mode&96 == 0 {
		draw(size)
		return
	}
	if x.deferred == nil {
		return
	}
	*x.deferred = append(*x.deferred, func(view cad.Rect) {
		// a size relative to the view: a percentage of its height
		pct := 5.0
		if size < 0 {
			pct = -size
		}
		s := cad.Scale(m)
		if s == 0 || !view.Valid() {
			return
		}
		draw(view.H() * pct / 100 / s)
	})
}

func (c *converter) ellipse(e *entity, x *ctx, p props) {
	ctr := e.vec3(10, [3]float64{})
	maj := e.vec3(11, [3]float64{1, 0, 0})
	n := norm(e.vec3(210, [3]float64{0, 0, 1}))
	ratio := e.num(40, 1)
	mn := cross(n, maj)
	for i := range mn {
		mn[i] *= ratio
	}
	t0, t1, ok := ccwSpan(e.num(41, 0), e.num(42, 2*math.Pi))
	if !ok {
		return
	}
	path := &cad.Path{}
	path.EllipseArcAxes(pt2(ctr), pt2(maj), pt2(mn), t0, t1)
	if math.Abs(t1-t0-2*math.Pi) < 1e-9 {
		path.Close()
	}
	x.out.Stroke(path.Transform(x.m), c.pen(e, x, p, x.m))
}

func (c *converter) spline(e *entity, x *ctx, p props) {
	var knots, weights []float64
	var ctrl, fit []cad.Point
	var start, end *cad.Point
	for _, t := range e.tags {
		switch t.code {
		case 40:
			knots = append(knots, t.f)
		case 41:
			weights = append(weights, t.f)
		case 10:
			ctrl = append(ctrl, cad.Point{X: t.f})
		case 20:
			if n := len(ctrl); n > 0 {
				ctrl[n-1].Y = t.f
			}
		case 11:
			fit = append(fit, cad.Point{X: t.f})
		case 21:
			if n := len(fit); n > 0 {
				fit[n-1].Y = t.f
			}
		case 12:
			start = &cad.Point{X: t.f}
		case 22:
			if start != nil {
				start.Y = t.f
			}
		case 13:
			end = &cad.Point{X: t.f}
		case 23:
			if end != nil {
				end.Y = t.f
			}
		}
	}
	flags := e.int(70, 0)
	path := &cad.Path{}
	switch {
	case len(ctrl) >= 2:
		if flags&4 == 0 {
			weights = nil
		}
		path.BSpline(e.int(71, 3), knots, ctrl, weights)
	case len(fit) >= 2:
		if flags&1 != 0 && fit[0] != fit[len(fit)-1] {
			fit = append(fit, fit[0])
		}
		path.Interpolate(fit, start, end)
	default:
		return
	}
	x.out.Stroke(path.Transform(x.m), c.pen(e, x, p, x.m))
}

func (c *converter) face(e *entity, x *ctx, p props) {
	pts := []cad.Point{e.pt(10), e.pt(11), e.pt(12), e.pt(13)}
	hidden := e.int(70, 0)
	path := &cad.Path{}
	for i := range 4 {
		if hidden&(1<<i) != 0 {
			continue
		}
		a, b := pts[i], pts[(i+1)%4]
		if a == b {
			continue
		}
		path.MoveTo(a.X, a.Y).LineTo(b.X, b.Y)
	}
	x.out.Stroke(path.Transform(x.m), c.pen(e, x, p, x.m))
}

// insert draws the block an INSERT refers to, as an array for MINSERT,
// and its attributes.
func (c *converter) insert(e *entity, x *ctx, p props, l *layer) {
	name := e.str(2)
	b := c.d.blocks[key(name)]
	switch {
	case b == nil:
		c.warnOnce("block:"+key(name), "block %q is missing", name)
		return
	case b.flags&4 != 0 || b.xref != "" && len(b.entities) == 0:
		c.warnOnce("xref:"+key(name), "external reference %q (%s) is not included in the drawing", name, b.xref)
		return
	}
	ins := e.vec3(10, [3]float64{})
	sx, sy := e.num(41, 1), e.num(42, 1)
	rot := e.num(50, 0)
	cols, rows := min(max(e.int(70, 1), 1), 32767), min(max(e.int(71, 1), 1), 32767)
	cs, rs := e.num(44, 0), e.num(45, 0)
	if cols*rows > 10000 {
		c.warnOnce("minsert", "block arrays of more than 10000 inserts are drawn once")
		cols, rows = 1, 1
	}
	base := x.m.Mul(ocs(e, ins[2])).Mul(canvas.Translate(ins[0], ins[1])).Mul(canvas.Rotate(rot))
	for r := range rows {
		for col := range cols {
			m := base.Mul(canvas.Translate(float64(col)*cs, float64(r)*rs)).Mul(canvas.Scale(sx, sy)).Mul(canvas.Translate(-b.base.X, -b.base.Y))
			c.block(b, x, m, p, l)
		}
	}
	if len(e.kids) > 0 {
		ax := *x
		ax.byBlock = p
		for _, a := range e.kids {
			c.entity(a, &ax)
		}
	}
}

// block draws the entities of a block through m, as inserted by an entity
// with the properties p on the layer l.
func (c *converter) block(b *block, x *ctx, m canvas.Matrix, p props, l *layer) {
	k := key(b.name)
	if slices.Contains(x.blocks, k) {
		c.warnOnce("cycle:"+k, "block %q inserts itself; the inner inserts are not drawn", b.name)
		return
	}
	sub := *x
	sub.blocks = append(slices.Clone(x.blocks), k)
	sub.m = m
	sub.insLayer = l
	sub.byBlock = p
	sub.depth = x.depth + 1
	for _, e := range b.entities {
		c.entity(e, &sub)
	}
}

func (c *converter) dimension(e *entity, x *ctx, p props, l *layer) {
	name := e.str(2)
	b := c.d.blocks[key(name)]
	if b == nil {
		c.warnOnce("dimblock", "dimensions without their geometry block are not drawn")
		return
	}
	elev := e.vec3(11, [3]float64{})[2]
	o := ocs(e, elev)
	m := x.m
	if e.has(12) {
		ins := e.vec3(12, [3]float64{})
		w := cad.Apply(ocs(e, ins[2]), pt2(ins))
		m = m.Mul(canvas.Translate(w.X, w.Y))
	}
	m = m.Mul(canvas.Matrix{o[0], o[1], o[2], o[3], 0, 0}).Mul(canvas.Translate(-b.base.X, -b.base.Y))
	c.block(b, x, m, p, l)
}

func (c *converter) table(e *entity, x *ctx, p props, l *layer) {
	b := c.d.blocks[key(e.str(2))]
	if b == nil {
		c.warnOnce("tableblock", "tables without their geometry block are not drawn")
		return
	}
	ins := e.vec3(10, [3]float64{})
	dir := e.vec3(11, [3]float64{1, 0, 0})
	m := x.m.Mul(canvas.Translate(ins[0], ins[1])).Mul(canvas.Rotate(math.Atan2(dir[1], dir[0]) * 180 / math.Pi)).
		Mul(canvas.Translate(-b.base.X, -b.base.Y))
	c.block(b, x, m, p, l)
}

func (c *converter) wipeout(e *entity, x *ctx) {
	pts := imageBoundary(e)
	if len(pts) < 3 {
		return
	}
	bg := bdf.RGB(255, 255, 255)
	if x.dark {
		bg = darkBackground
	}
	x.out.Fill((&cad.Path{}).Polyline(pts, true).Transform(x.m), cad.Fill{Color: bg}, false)
}

// imageBoundary returns the clip boundary of an IMAGE or WIPEOUT in the
// world.
func imageBoundary(e *entity) []cad.Point {
	ins := pt2(e.vec3(10, [3]float64{}))
	u := pt2(e.vec3(11, [3]float64{1, 0, 0}))
	v := pt2(e.vec3(12, [3]float64{0, 1, 0}))
	size := e.pt(13)
	var bp []cad.Point
	for _, t := range e.tags {
		switch t.code {
		case 14:
			bp = append(bp, cad.Point{X: t.f})
		case 24:
			if n := len(bp); n > 0 {
				bp[n-1].Y = t.f
			}
		}
	}
	if len(bp) == 2 {
		// a rectangle by two corners
		a, b := bp[0], bp[1]
		bp = []cad.Point{a, {X: b.X, Y: a.Y}, b, {X: a.X, Y: b.Y}}
	}
	// the boundary's origin is the top left corner of the image, y down
	origin := ins.Add(u.Mul(0.5)).Sub(v.Mul(0.5))
	out := make([]cad.Point, len(bp))
	for i, q := range bp {
		out[i] = origin.Add(u.Mul(q.X)).Add(v.Mul(size.Y - q.Y))
	}
	return out
}

func (c *converter) mesh(e *entity, x *ctx, p props) {
	var verts []cad.Point
	i := 0
	tags := e.tags
	// vertices: 92 count, then 10/20/30 each; faces: 93 size, then 90 values
	for ; i < len(tags); i++ {
		if tags[i].code == 92 {
			break
		}
	}
	for i++; i < len(tags) && tags[i].code != 93; i++ {
		switch tags[i].code {
		case 10:
			verts = append(verts, cad.Point{X: tags[i].f})
		case 20:
			if n := len(verts); n > 0 {
				verts[n-1].Y = tags[i].f
			}
		}
	}
	var list []int
	for i++; i < len(tags) && tags[i].code == 90; i++ {
		list = append(list, int(tags[i].f))
	}
	path := &cad.Path{}
	for j := 0; j < len(list); {
		n := list[j]
		j++
		if n <= 0 || j+n > len(list) {
			break
		}
		face := list[j : j+n]
		j += n
		for k := range face {
			a, b := face[k], face[(k+1)%n]
			if a < 0 || b < 0 || a >= len(verts) || b >= len(verts) {
				continue
			}
			path.MoveTo(verts[a].X, verts[a].Y).LineTo(verts[b].X, verts[b].Y)
		}
	}
	x.out.Stroke(path.Transform(x.m), c.pen(e, x, p, x.m))
}
