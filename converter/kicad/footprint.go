package kicad

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
)

// footprint draws a footprint: its graphics, texts and pads. Children are
// in the footprint's own space, turned by its orientation; pads and texts
// give their orientation on the board.
func (b *pcb) footprint(fp *node) {
	at := fp.child("at")
	pl := fpPlace{x: at.num(0), y: at.num(1), angle: at.num(2)}
	fields := map[string]string{}
	for _, p := range fp.children("property") {
		fields[p.arg(0)] = p.arg(1)
	}
	for _, t := range fp.children("fp_text") {
		switch t.arg(0) {
		case "reference":
			fields["Reference"] = t.arg(1)
		case "value":
			fields["Value"] = t.arg(1)
		}
	}
	name := fp.arg(0)
	vars := func(v string) (string, bool) {
		switch strings.ToUpper(v) {
		case "REFERENCE":
			return fields["Reference"], true
		case "VALUE":
			return fields["Value"], true
		case "FOOTPRINT_NAME":
			if _, n, ok := strings.Cut(name, ":"); ok {
				return n, true
			}
			return name, true
		case "FOOTPRINT_LIBRARY":
			l, _, _ := strings.Cut(name, ":")
			return l, true
		}
		if s, ok := fields[v]; ok {
			return s, true
		}
		return b.vars(v)
	}
	margin := fp.numOf("solder_mask_margin", b.maskMargin)
	for _, it := range fp.lists() {
		switch it.name {
		case "fp_line", "fp_rect", "fp_circle", "fp_arc", "fp_poly", "fp_curve":
			b.graphic(it, pl, fp)
		case "fp_text":
			if it.arg(0) == "reference" || it.arg(0) == "value" || it.arg(0) == "user" {
				b.boardText(it, pl, vars, true)
			}
		case "property":
			if it.child("at") != nil && it.child("layer") != nil {
				b.boardText(it, pl, vars, true)
			}
		case "fp_text_box":
			b.textBox(it, pl, vars)
		case "pad":
			b.pad(it, pl, margin)
		case "zone":
			b.zone(it)
		}
	}
}

// boardText draws a text of the board or of a footprint (inFootprint: kept
// upright unless unlocked), or the outlines KiCad cached for text in a
// TrueType font.
func (b *pcb) boardText(it *node, pl fpPlace, vars func(string) (string, bool), inFootprint bool) {
	e := parseEffects(it.child("effects"))
	if e.hide || it.flag("hide") || it.hasAtom("hide") {
		return
	}
	ln := it.child("layer").arg(0)
	if ln == "" {
		return
	}
	s := it.arg(0)
	switch it.name {
	case "fp_text", "property":
		s = it.arg(1)
	}
	lookup := vars
	if lookup == nil {
		lookup = b.vars
	}
	s = b.c.expand(s, func(v string) (string, bool) {
		if v == "LAYER" {
			return b.layer(ln).title, true
		}
		return lookup(v)
	})
	if strings.TrimSpace(s) == "" {
		return
	}
	color := pcbLayerColor(ln)
	d := b.layer(ln).d
	if rc := it.child("render_cache"); rc != nil && len(rc.children("polygon")) > 0 {
		path := &cad.Path{}
		for _, poly := range rc.children("polygon") {
			for _, pts := range poly.children("pts") {
				var q []cad.Point
				for _, xy := range pts.children("xy") {
					q = append(q, cad.Point{X: xy.num(0), Y: xy.num(1)})
				}
				if len(q) >= 3 {
					path.Polyline(q, true)
				}
			}
		}
		d.Fill(path, cad.Fill{Color: color}, false)
		// the text itself, invisible, for search and selection
		e2 := e
		e2.face = ""
		b.textAt(d, it, pl, s, e2, bdf.RGBA(0, 0, 0, 0), inFootprint)
		return
	}
	b.textAt(d, it, pl, s, e, color, inFootprint)
}

// textAt draws a text at its (at x y angle).
func (b *pcb) textAt(d *cad.Drawing, it *node, pl fpPlace, s string, e effects, color bdf.Color, inFootprint bool) {
	at := it.child("at")
	p := pl.pt(at.num(0), at.num(1))
	angle := at.num(2)
	if inFootprint && !at.hasAtom("unlocked") && !it.flag("unlocked") {
		// kept upright: (-90, 90]
		angle = math.Mod(angle, 180)
		if angle > 90 {
			angle -= 180
		} else if angle <= -90 {
			angle += 180
		}
	}
	b.c.text(d, placedText{s: s, x: p.X, y: p.Y, angle: angle, e: e, color: color, thickness: penWidth(e, 0), brk: cad.BreakBox})
}

// textBox draws a text box of the board or of a footprint.
func (b *pcb) textBox(it *node, pl fpPlace, vars func(string) (string, bool)) {
	ln := it.child("layer").arg(0)
	if ln == "" {
		return
	}
	color := pcbLayerColor(ln)
	d := b.layer(ln).d
	var pts []cad.Point
	if s, e := it.child("start"), it.child("end"); s != nil && e != nil {
		pts = []cad.Point{pl.pt(s.num(0), s.num(1)), pl.pt(e.num(0), s.num(1)), pl.pt(e.num(0), e.num(1)), pl.pt(s.num(0), e.num(1))}
	} else {
		for _, q := range it.child("pts").children("xy") {
			pts = append(pts, pl.pt(q.num(0), q.num(1)))
		}
	}
	if len(pts) < 3 {
		return
	}
	var box cad.Rect
	for _, q := range pts {
		box = box.Add(q)
	}
	if it.flagOr("border", true) {
		w, dash := strokeWidth(it)
		if w > 0 {
			d.Stroke((&cad.Path{}).Polyline(pts, true), pen(w, dash, color))
		}
	}
	lookup := vars
	if lookup == nil {
		lookup = b.vars
	}
	e := parseEffects(it.child("effects"))
	e.color, e.hasColor = color, true
	sp := &schPage{c: b.c, d: d}
	s := b.c.expand(it.arg(0), lookup)
	sp.boxText(s, e, box.Min.X, box.Min.Y, box.Max.X, box.Max.Y, it.child("margins"), color)
}

// pad draws a pad on its layers, and its hole.
func (b *pcb) pad(p *node, pl fpPlace, fpMargin float64) {
	typ, shape := p.arg(1), p.arg(2)
	at := p.child("at")
	pos := pl.pt(at.num(0), at.num(1))
	angle := at.num(2)
	sz := p.child("size")
	w, h := sz.num(0), sz.num(1)
	drill := p.child("drill")
	var ox, oy float64
	if off := drill.child("offset"); off != nil {
		ox, oy = off.num(0), off.num(1)
	}
	tr := func(x, y float64) cad.Point {
		rx, ry := rotatePt(x+ox, y+oy, angle)
		return cad.Point{X: pos.X + rx, Y: pos.Y + ry}
	}
	margin := p.numOf("solder_mask_margin", fpMargin)
	for _, ln := range b.expandLayers(layerNames(p)) {
		l := b.layer(ln)
		color := pcbLayerColor(ln)
		grow := 0.0
		switch {
		case strings.HasSuffix(ln, ".Cu"):
			if typ == "np_thru_hole" {
				continue
			}
			if typ == "thru_hole" {
				color = pcbPadThroughHole
			}
		case strings.HasSuffix(ln, ".Mask"):
			grow = margin
		case strings.HasSuffix(ln, ".Paste"):
			if typ != "smd" {
				continue
			}
		}
		b.padShape(l.d, p, shape, w, h, grow, tr, color)
	}
	if drill == nil || (typ != "thru_hole" && typ != "np_thru_hole") {
		return
	}
	dx := drill.num(0)
	dy := dx
	if drill.hasAtom("oval") {
		// (drill oval w h)
		dx, dy = numArgs(drill)
	}
	if dx <= 0 {
		return
	}
	hole := capsule(dx, dy, func(x, y float64) cad.Point {
		rx, ry := rotatePt(x, y, angle)
		return cad.Point{X: pos.X + rx, Y: pos.Y + ry}
	})
	color := pcbHoleBackground
	if typ == "np_thru_hole" {
		color = pcbNPTHole
	}
	b.holes.Fill(hole, cad.Fill{Color: color}, false)
}

// numArgs returns the first two numeric arguments of a list.
func numArgs(n *node) (float64, float64) {
	var v []float64
	for i := 0; n.arg(i) != "" && len(v) < 2; i++ {
		if f, err := parseNum(n.arg(i)); err == nil {
			v = append(v, f)
		}
	}
	switch len(v) {
	case 0:
		return 0, 0
	case 1:
		return v[0], v[0]
	}
	return v[0], v[1]
}

// padShape draws a pad's shape grown by grow (a mask's margin).
func (b *pcb) padShape(d *cad.Drawing, p *node, shape string, w, h, grow float64, tr func(x, y float64) cad.Point, color bdf.Color) {
	w, h = w+2*grow, h+2*grow
	if w <= 0 || h <= 0 {
		return
	}
	fill := func(path *cad.Path) { d.Fill(path, cad.Fill{Color: color}, false) }
	switch shape {
	case "circle":
		fill((&cad.Path{}).Circle(tr(0, 0), w/2))
	case "oval":
		fill(capsule(w, h, tr))
	case "rect":
		if grow > 0 {
			fill(roundedBox(w, h, grow, tr))
		} else {
			fill(roundedBox(w, h, 0, tr))
		}
	case "roundrect":
		r := p.numOf("roundrect_rratio", 0.25)*math.Min(w-2*grow, h-2*grow) + grow
		if cr := p.numOf("chamfer_ratio", 0); cr > 0 && p.child("chamfer") != nil {
			fill(chamferedBox(w, h, cr*math.Min(w, h), r, p.child("chamfer"), tr))
			break
		}
		fill(roundedBox(w, h, r, tr))
	case "trapezoid":
		dl := p.child("rect_delta")
		tx, ty := dl.num(0)/2, dl.num(1)/2
		hx, hy := w/2, h/2
		pts := []cad.Point{tr(-hx-ty, hy+tx), tr(hx+ty, hy-tx), tr(hx-ty, -hy+tx), tr(-hx+ty, -hy-tx)}
		path := (&cad.Path{}).Polyline(pts, true)
		fill(path)
		if grow > 0 {
			d.Stroke(path, cad.Pen{Color: color, WorldWidth: 2 * grow, Cap: cad.CapRound, Join: cad.JoinRound})
		}
	case "custom":
		opts := p.child("options")
		anchor := opts.str("anchor")
		as := math.Min(w-2*grow, h-2*grow)
		if anchor == "rect" {
			fill(roundedBox(as+2*grow, as+2*grow, grow, tr))
		} else {
			fill((&cad.Path{}).Circle(tr(0, 0), as/2+grow))
		}
		for _, prim := range p.child("primitives").lists() {
			b.padPrimitive(d, prim, grow, tr, color)
		}
	default:
		fill(roundedBox(w, h, 0, tr))
	}
}

// padPrimitive draws a primitive of a custom pad, in the pad's space.
func (b *pcb) padPrimitive(d *cad.Drawing, it *node, grow float64, tr func(x, y float64) cad.Point, color bdf.Color) {
	pt := func(n *node) cad.Point { return tr(n.num(0), n.num(1)) }
	w, _ := strokeWidth(it)
	w += 2 * grow
	pn := cad.Pen{Color: color, WorldWidth: w, Cap: cad.CapRound, Join: cad.JoinRound}
	path := &cad.Path{}
	closed := false
	switch it.name {
	case "gr_poly":
		var pts []cad.Point
		for _, q := range it.child("pts").children("xy") {
			pts = append(pts, pt(q))
		}
		if len(pts) < 3 {
			return
		}
		path.Polyline(pts, true)
		closed = true
	case "gr_line":
		path.Polyline([]cad.Point{pt(it.child("start")), pt(it.child("end"))}, false)
	case "gr_rect":
		s, e := it.child("start"), it.child("end")
		path.Polyline([]cad.Point{tr(s.num(0), s.num(1)), tr(e.num(0), s.num(1)), tr(e.num(0), e.num(1)), tr(s.num(0), e.num(1))}, true)
		closed = true
	case "gr_circle":
		c, e := it.child("center"), it.child("end")
		r := math.Hypot(e.num(0)-c.num(0), e.num(1)-c.num(1))
		path.Circle(pt(c), r)
		closed = true
	case "gr_arc":
		if c, r, a0, a1, ok := arcThrough(pt(it.child("start")), pt(it.child("mid")), pt(it.child("end"))); ok {
			path.Arc(c, r, a0, a1)
		}
	case "gr_curve":
		var pts []cad.Point
		for _, q := range it.child("pts").children("xy") {
			pts = append(pts, pt(q))
		}
		if len(pts) == 4 {
			path.MoveTo(pts[0].X, pts[0].Y)
			path.CubicTo(pts[1], pts[2], pts[3])
		}
	default:
		return
	}
	if path.Empty() {
		return
	}
	if closed && (filled(it) || it.name == "gr_poly") {
		d.Fill(path, cad.Fill{Color: color}, false)
	}
	if w > 0 {
		d.Stroke(path, pn)
	}
}

// capsule is an oval: a circle, or a slot of two half circles.
func capsule(w, h float64, tr func(x, y float64) cad.Point) *cad.Path {
	p := &cad.Path{}
	if math.Abs(w-h) < 1e-9 {
		return p.Circle(tr(0, 0), w/2)
	}
	r := math.Min(w, h) / 2
	pts := []cad.Point{}
	arc := func(cx, cy, a0 float64) {
		for i := 0; i <= 16; i++ {
			a := a0 + math.Pi*float64(i)/16
			pts = append(pts, tr(cx+r*math.Cos(a), cy+r*math.Sin(a)))
		}
	}
	if w > h {
		l := (w - h) / 2
		arc(l, 0, -math.Pi/2)
		arc(-l, 0, math.Pi/2)
	} else {
		l := (h - w) / 2
		arc(0, l, 0)
		arc(0, -l, math.Pi)
	}
	return p.Polyline(pts, true)
}

// roundedBox is a w×h rectangle about the origin with corners of radius r.
func roundedBox(w, h, r float64, tr func(x, y float64) cad.Point) *cad.Path {
	hx, hy := w/2, h/2
	r = math.Max(0, math.Min(r, math.Min(hx, hy)))
	if r == 0 {
		return (&cad.Path{}).Polyline([]cad.Point{tr(-hx, -hy), tr(hx, -hy), tr(hx, hy), tr(-hx, hy)}, true)
	}
	var pts []cad.Point
	corner := func(cx, cy, a0 float64) {
		for i := 0; i <= 8; i++ {
			a := a0 + math.Pi/2*float64(i)/8
			pts = append(pts, tr(cx+r*math.Cos(a), cy+r*math.Sin(a)))
		}
	}
	corner(hx-r, -hy+r, -math.Pi/2)
	corner(hx-r, hy-r, 0)
	corner(-hx+r, hy-r, math.Pi/2)
	corner(-hx+r, -hy+r, math.Pi)
	return (&cad.Path{}).Polyline(pts, true)
}

// chamferedBox is a rectangle whose named corners are cut by c and the
// others rounded by r.
func chamferedBox(w, h, c, r float64, corners *node, tr func(x, y float64) cad.Point) *cad.Path {
	has := func(name string) bool { return corners.hasAtom(name) }
	hx, hy := w/2, h/2
	c = math.Min(c, math.Min(hx, hy))
	r = math.Max(0, math.Min(r, math.Min(hx, hy)))
	var pts []cad.Point
	// top left, top right, bottom right, bottom left (y down)
	type corner struct {
		x, y, sx, sy float64
		name         string
	}
	cs := []corner{{-hx, -hy, 1, 1, "top_left"}, {hx, -hy, -1, 1, "top_right"}, {hx, hy, -1, -1, "bottom_right"}, {-hx, hy, 1, -1, "bottom_left"}}
	for i, k := range cs {
		switch {
		case has(k.name):
			// the cut, from the side before to the side after
			if i%2 == 0 {
				pts = append(pts, tr(k.x, k.y+k.sy*c), tr(k.x+k.sx*c, k.y))
			} else {
				pts = append(pts, tr(k.x+k.sx*c, k.y), tr(k.x, k.y+k.sy*c))
			}
		case r > 0:
			cx, cy := k.x+k.sx*r, k.y+k.sy*r
			a0 := math.Atan2(-k.sy, 0)
			if i%2 == 0 {
				a0 = math.Atan2(0, -k.sx)
			}
			for j := 0; j <= 8; j++ {
				a := a0 + math.Pi/2*float64(j)/8
				pts = append(pts, tr(cx+r*math.Cos(a), cy+r*math.Sin(a)))
			}
		default:
			pts = append(pts, tr(k.x, k.y))
		}
	}
	return (&cad.Path{}).Polyline(pts, true)
}
