package kicad

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

func parseNum(s string) (float64, error) { return strconv.ParseFloat(s, 64) }

// magnification returns how many times larger than the board its pages
// are: the largest whole number that keeps the long side within that of
// A3 (420 mm), as the Gerber converter draws boards; larger boards are
// drawn smaller.
func magnification(long float64) float64 {
	if long <= 0 {
		return 1
	}
	if long > 420 {
		return 420 / long
	}
	return math.Min(math.Floor(420/long), 1000)
}

// stack returns the layers of the board view of a side, from the bottom
// up: the other side's layers, the inner copper, the side's own layers
// with its copper on top, then the drawing layers.
func (b *pcb) stack(front bool) []string {
	side := func(p string) []string {
		return []string{p + "Fab", p + "CrtYd", p + "Paste", p + "Adhes", p + "SilkS", p + "Mask", p + "Cu"}
	}
	var out []string
	var inner []string
	for _, cu := range b.copper {
		if cu != "F.Cu" && cu != "B.Cu" {
			inner = append(inner, cu)
		}
	}
	if front {
		out = append(out, side("B.")...)
		for i := len(inner) - 1; i >= 0; i-- {
			out = append(out, inner[i])
		}
		out = append(out, side("F.")...)
	} else {
		out = append(out, side("F.")...)
		out = append(out, inner...)
		out = append(out, side("B.")...)
	}
	return out
}

// drawingLayers are the layers drawn over the board in its views.
func (b *pcb) drawingLayers() []string {
	var out []string
	for _, n := range b.order {
		if strings.HasSuffix(n, ".Cu") || strings.HasPrefix(n, "F.") || strings.HasPrefix(n, "B.") {
			continue
		}
		if n == "Edge.Cuts" {
			continue
		}
		out = append(out, n)
	}
	return append(out, "Edge.Cuts")
}

// views makes the views of the board: its front and back as the board
// editor shows them, and a view of each layer with something on it.
func (b *pcb) views() ([]pageOut, error) {
	c := b.c
	// the page holds the board and what is drawn around it (dimensions,
	// notes), unless things lie far from the board: then the board and a
	// margin
	var edge, all cad.Rect
	if l := b.layers["Edge.Cuts"]; l != nil {
		edge = l.d.Bounds()
	}
	for _, l := range b.layers {
		all = all.Union(l.d.Bounds())
	}
	all = all.Union(b.holes.Bounds())
	for _, v := range b.vias {
		all = all.Add(v.c)
	}
	extent := all
	if edge.Valid() && all.W()*all.H() > 4*edge.W()*edge.H() {
		gx, gy := edge.W()*0.2, edge.H()*0.2
		extent = edge.Add(cad.Point{X: edge.Min.X - gx, Y: edge.Min.Y - gy}).Add(cad.Point{X: edge.Max.X + gx, Y: edge.Max.Y + gy})
	}
	if !extent.Valid() {
		return nil, errNothing
	}
	long := math.Max(extent.W(), extent.H())
	margin := math.Max(1, 0.03*long)
	mag := magnification(long + 2*margin)
	k := mag * ptPerMM
	x0, y0 := extent.Min.X-margin, extent.Min.Y-margin
	w, h := (extent.W()+2*margin)*k, (extent.H()+2*margin)*k
	m := canvas.Scale(k, k).Mul(canvas.Translate(-x0, -y0))
	pageRect := cad.Rect{}.Add(cad.Point{}).Add(cad.Point{X: w, Y: h})
	c.plotter.Thin = pcbHairline * k
	// the board views show every via over the holes of pads
	allHoles := &cad.Drawing{}
	allHoles.Items = append(allHoles.Items, b.holes.Items...)
	for _, v := range b.vias {
		drawVia(allHoles, v)
	}
	// each layer is drawn once, as an object the views share
	objs := map[string]*canvas.Canvas{}
	layerObj := func(name string, d *cad.Drawing) *canvas.Canvas {
		if cv, ok := objs[name]; ok {
			return cv
		}
		if d == nil || d.Empty() {
			return nil
		}
		cv := c.cvs.New()
		cv.Obj.SetBBox(0, 0, float32(w), float32(h))
		c.plotter.Plot(cv, d, m, pageRect)
		objs[name] = cv
		return cv
	}
	var out []pageOut
	layerMode := strings.ToLower(c.opts.Layers)
	page := func(id, title string, names []string, mirror bool, holes bool) {
		view := c.doc.NewView(id, bdf.ViewFixed, title)
		pg := view.AddPage(float32(w), float32(h))
		cv := c.cvs.New()
		cv.Obj.SetBBox(0, 0, float32(w), float32(h))
		cv.Obj.FillColor(pcbBackground)
		cv.Obj.FillRect(0, 0, float32(w), float32(h))
		if mirror {
			cv.Obj.Transform(-1, 0, 0, 1, float32(w), 0)
		}
		bbox := bdf.Rect{X: 0, Y: 0, W: float32(w), H: float32(h)}
		for _, n := range names {
			var d *cad.Drawing
			if l := b.layers[n]; l != nil {
				d = l.d
			}
			switch {
			case n == ":holes":
				d = allHoles
			case strings.HasPrefix(n, ":holes:"):
				d = b.holesFor(n[len(":holes:"):])
			}
			if ch := layerObj(n, d); ch != nil {
				ref := cv.Share(ch, bbox)
				cv.Obj.Use(ref)
				cv.Used(ch)
			}
		}
		out = append(out, pageOut{view: view, page: pg, cv: cv})
	}
	if layerMode != "layers" {
		front := append(b.stack(true), ":holes")
		front = append(front, b.drawingLayers()...)
		page("board-front", "Front", front, false, true)
		back := append(b.stack(false), ":holes")
		back = append(back, b.drawingLayers()...)
		page("board-back", "Back (seen from below)", back, true, true)
	}
	if layerMode != "board" {
		for _, n := range b.order {
			l := b.layers[n]
			if l == nil || l.d.Empty() || n == "Edge.Cuts" {
				continue
			}
			names := []string{n}
			if strings.HasSuffix(n, ".Cu") {
				names = append(names, ":holes:"+n)
			}
			names = append(names, "Edge.Cuts")
			page("layer-"+n, l.title, names, strings.HasPrefix(n, "B."), false)
		}
		if l := b.layers["Edge.Cuts"]; l != nil && !l.d.Empty() && len(out) == 0 {
			page("layer-Edge.Cuts", l.title, []string{"Edge.Cuts"}, false, false)
		}
	}
	return out, nil
}

// dimension draws a dimension: its lines and arrows, and its text.
func (b *pcb) dimension(it *node) {
	ln := it.str("layer")
	if ln == "" {
		return
	}
	d := b.layer(ln).d
	color := pcbLayerColor(ln)
	st := it.child("style")
	thick := st.numOf("thickness", 0.15)
	arrow := st.numOf("arrow_length", 1.27)
	extOff := st.numOf("extension_offset", 0.5)
	extH := st.numOf("extension_height", 0.58642)
	pn := pen(thick, "solid", color)
	var pts []cad.Point
	for _, q := range it.child("pts").children("xy") {
		pts = append(pts, cad.Point{X: q.num(0), Y: q.num(1)})
	}
	line := func(a, c cad.Point) { d.Stroke((&cad.Path{}).Polyline([]cad.Point{a, c}, false), pn) }
	// arrowhead at a pointing away from c (outward) or towards it
	head := func(a, c cad.Point, inward bool) {
		dx, dy := a.X-c.X, a.Y-c.Y
		l := math.Hypot(dx, dy)
		if l == 0 {
			return
		}
		dx, dy = dx/l, dy/l
		if inward {
			dx, dy = -dx, -dy
		}
		for _, s := range []float64{1, -1} {
			ang := s * 27.5 * math.Pi / 180
			ex := -(dx*math.Cos(ang) - dy*math.Sin(ang)) * arrow
			ey := -(dx*math.Sin(ang) + dy*math.Cos(ang)) * arrow
			line(a, cad.Point{X: a.X + ex, Y: a.Y + ey})
		}
	}
	inward := st.str("arrow_direction") == "inward"
	typ := it.child("type").arg(0)
	txt := it.child("gr_text")
	switch typ {
	case "aligned", "orthogonal":
		if len(pts) < 2 {
			break
		}
		p1, p2 := pts[0], pts[1]
		hgt := it.numOf("height", 0)
		var n cad.Point
		var a, c cad.Point
		if typ == "orthogonal" {
			if it.numOf("orientation", 0) == 0 { // horizontal
				y := p1.Y + hgt
				a, c = cad.Point{X: p1.X, Y: y}, cad.Point{X: p2.X, Y: y}
				n = cad.Point{Y: 1}
			} else {
				x := p1.X + hgt
				a, c = cad.Point{X: x, Y: p1.Y}, cad.Point{X: x, Y: p2.Y}
				n = cad.Point{X: 1}
			}
		} else {
			dx, dy := p2.X-p1.X, p2.Y-p1.Y
			l := math.Hypot(dx, dy)
			if l == 0 {
				break
			}
			// the crossbar lies height away, on the side the direction turns to
			// by 90° on the page
			n = cad.Point{X: -dy / l, Y: dx / l}
			a = cad.Point{X: p1.X + n.X*hgt, Y: p1.Y + n.Y*hgt}
			c = cad.Point{X: p2.X + n.X*hgt, Y: p2.Y + n.Y*hgt}
		}
		sgn := 1.0
		if hgt < 0 || (typ == "orthogonal" && ((n.Y == 1 && a.Y < p1.Y) || (n.X == 1 && a.X < p1.X))) {
			sgn = -1
		}
		ext := func(p, q cad.Point) {
			dx, dy := q.X-p.X, q.Y-p.Y
			l := math.Hypot(dx, dy)
			if l < 1e-9 {
				return
			}
			ux, uy := dx/l, dy/l
			start := cad.Point{X: p.X + ux*extOff, Y: p.Y + uy*extOff}
			end := cad.Point{X: q.X + ux*extH, Y: q.Y + uy*extH}
			if l > extOff {
				line(start, end)
			}
		}
		_ = sgn
		ext(p1, a)
		ext(p2, c)
		line(a, c)
		head(a, c, inward)
		head(c, a, inward)
	case "leader":
		if len(pts) < 2 {
			break
		}
		line(pts[0], pts[1])
		head(pts[0], pts[1], true)
		if txt != nil {
			tx, ty := txt.child("at").xy()
			line(pts[1], cad.Point{X: tx, Y: ty})
		}
	case "center":
		if len(pts) < 2 {
			break
		}
		cx, cy := pts[0].X, pts[0].Y
		dx, dy := pts[1].X-cx, pts[1].Y-cy
		line(cad.Point{X: cx - dx, Y: cy - dy}, cad.Point{X: cx + dx, Y: cy + dy})
		line(cad.Point{X: cx + dy, Y: cy - dx}, cad.Point{X: cx - dy, Y: cy + dx})
	case "radial":
		if len(pts) < 2 {
			break
		}
		line(pts[0], pts[1])
		head(pts[1], pts[0], false)
		if ll := it.numOf("leader_length", 0); ll > 0 && txt != nil {
			tx, ty := txt.child("at").xy()
			line(pts[1], cad.Point{X: tx, Y: ty})
		}
	}
	if txt != nil {
		if txt.child("layer") == nil {
			txt.items = append(txt.items, item{list: &node{name: "layer", items: []item{{s: ln}}}})
		}
		b.boardText(txt, fpPlace{}, nil, false)
	}
}

// table draws a table of the board.
func (b *pcb) table(t *node) {
	ln := t.str("layer")
	if ln == "" {
		return
	}
	color := pcbLayerColor(ln)
	d := b.layer(ln).d
	sp := &schPage{c: b.c, d: d}
	var all cad.Rect
	var cells []cad.Rect
	for _, cell := range t.child("cells").children("table_cell") {
		s, e := cell.child("start"), cell.child("end")
		if s == nil || e == nil {
			continue
		}
		r := cad.Rect{}.Add(cad.Point{X: s.num(0), Y: s.num(1)}).Add(cad.Point{X: e.num(0), Y: e.num(1)})
		if r.W() <= 0 && r.H() <= 0 {
			continue
		}
		all = all.Union(r)
		cells = append(cells, r)
		ef := parseEffects(cell.child("effects"))
		sp.boxText(b.c.expand(cell.arg(0), b.vars), ef, r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, cell.child("margins"), color)
	}
	if !all.Valid() {
		return
	}
	bw, _ := strokeWidth(t.child("border"))
	if bw <= 0 {
		bw = 0.1
	}
	sw, _ := strokeWidth(t.child("separators"))
	if sw <= 0 {
		sw = 0.1
	}
	for _, r := range cells {
		if r.Min.Y > all.Min.Y+1e-6 {
			d.Stroke((&cad.Path{}).Polyline([]cad.Point{r.Min, {X: r.Max.X, Y: r.Min.Y}}, false), pen(sw, "solid", color))
		}
		if r.Min.X > all.Min.X+1e-6 {
			d.Stroke((&cad.Path{}).Polyline([]cad.Point{r.Min, {X: r.Min.X, Y: r.Max.Y}}, false), pen(sw, "solid", color))
		}
	}
	d.Stroke((&cad.Path{}).Polyline([]cad.Point{all.Min, {X: all.Max.X, Y: all.Min.Y}, all.Max, {X: all.Min.X, Y: all.Max.Y}}, true), pen(bw, "solid", color))
}
