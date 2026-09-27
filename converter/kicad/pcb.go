package kicad

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
)

// Opacities of the board editor's default appearance.
const (
	zoneOpacity = 0.6
	pcbHairline = 0.1 // mm, for lines of no width
)

// pcbLayer is a layer of the board with what is drawn on it.
type pcbLayer struct {
	name  string // canonical name (F.Cu)
	title string // user name (F.Silkscreen)
	d     *cad.Drawing
}

// pcb is a board being converted.
type pcb struct {
	c          *conv
	n          *node
	layers     map[string]*pcbLayer
	order      []string     // every layer, in the board's table order
	copper     []string     // copper layers from the front to the back
	holes      *cad.Drawing // the holes of pads
	vias       []viaSpan    // the vias, drawn over the layers they join
	outline    *cad.Drawing // Edge.Cuts, for the layer views
	vars       func(string) (string, bool)
	maskMargin float64
}

// board converts the board of an input into views.
func (c *conv) board(in *input) ([]pageOut, error) {
	n, err := parse(in.pcbData)
	if err != nil {
		return nil, fmt.Errorf("kicad: %s: %w", in.pcb, err)
	}
	if n.name != "kicad_pcb" {
		return nil, fmt.Errorf("kicad: %s: not a KiCad board", in.pcb)
	}
	if v := n.child("version").int(0); v > 0 && v < 20211014 {
		return nil, fmt.Errorf("kicad: %s: a board of KiCad 5 or earlier (version %d); open it in KiCad 6 or later and save it", in.pcb, v)
	}
	c.iu = 1e6 // 1 nm
	c.defaultPen = 0
	b := &pcb{c: c, n: n, layers: map[string]*pcbLayer{}, holes: &cad.Drawing{}}
	b.readLayers()
	b.maskMargin = n.child("setup").numOf("pad_to_mask_clearance", 0)
	tb := n.child("title_block")
	if c.opts.Title == "" && len(c.doc.Meta.DC.Title) == 0 && tb.str("title") != "" {
		c.doc.Meta.DC.Title = bdf.DCValues{tb.str("title")}
	}
	b.vars = func(name string) (string, bool) {
		switch name {
		case "TITLE":
			return tb.str("title"), true
		case "ISSUE_DATE":
			return tb.str("date"), true
		case "REVISION":
			return tb.str("rev"), true
		case "COMPANY":
			return tb.str("company"), true
		case "FILENAME":
			return in.pcb, true
		}
		if strings.HasPrefix(name, "COMMENT") {
			for _, cm := range tb.children("comment") {
				if "COMMENT"+cm.arg(0) == name {
					return cm.arg(1), true
				}
			}
			return "", true
		}
		return "", false
	}
	b.draw()
	return b.views()
}

// readLayers reads the layer table: names, user names and the copper
// stack.
func (b *pcb) readLayers() {
	for _, l := range b.n.child("layers").lists() {
		name := l.arg(0)
		if name == "" {
			continue
		}
		title := l.arg(2)
		if title == "" {
			title = name
		}
		b.layers[name] = &pcbLayer{name: name, title: title, d: &cad.Drawing{}}
		b.order = append(b.order, name)
		if strings.HasSuffix(name, ".Cu") {
			b.copper = append(b.copper, name)
		}
	}
	// front, inner layers in order, back
	slices.SortStableFunc(b.copper, func(x, y string) int { return copperRank(x) - copperRank(y) })
}

func copperRank(name string) int {
	switch name {
	case "F.Cu":
		return 0
	case "B.Cu":
		return 1000
	}
	n, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(name, "In"), ".Cu"))
	return n
}

// layer returns the drawing of a layer, adding one the table lacks.
func (b *pcb) layer(name string) *pcbLayer {
	if l, ok := b.layers[name]; ok {
		return l
	}
	l := &pcbLayer{name: name, title: name, d: &cad.Drawing{}}
	b.layers[name] = l
	b.order = append(b.order, name)
	if strings.HasSuffix(name, ".Cu") {
		b.copper = append(b.copper, name)
		slices.SortStableFunc(b.copper, func(x, y string) int { return copperRank(x) - copperRank(y) })
	}
	return l
}

// expandLayers turns the layer names of a pad or zone (with *.Cu, *.Mask,
// F&B.Cu) into layers of the board.
func (b *pcb) expandLayers(names []string) []string {
	var out []string
	add := func(s string) {
		if !slices.Contains(out, s) {
			out = append(out, s)
		}
	}
	for _, s := range names {
		switch s {
		case "*.Cu":
			for _, cu := range b.copper {
				add(cu)
			}
		case "F&B.Cu":
			add("F.Cu")
			add("B.Cu")
		case "*.Mask":
			add("F.Mask")
			add("B.Mask")
		case "*.Paste":
			add("F.Paste")
			add("B.Paste")
		case "*.SilkS":
			add("F.SilkS")
			add("B.SilkS")
		default:
			if strings.HasPrefix(s, "*.") {
				add("F." + s[2:])
				add("B." + s[2:])
			} else {
				add(s)
			}
		}
	}
	return out
}

func layerNames(n *node) []string {
	var out []string
	if l := n.child("layers"); l != nil {
		for i := 0; l.arg(i) != ""; i++ {
			out = append(out, l.arg(i))
		}
	}
	if l := n.child("layer"); l != nil && l.arg(0) != "" {
		out = append(out, l.arg(0))
	}
	return out
}

// place maps a footprint's local points onto the board.
type fpPlace struct {
	x, y, angle float64
}

func (f fpPlace) pt(x, y float64) cad.Point {
	rx, ry := rotatePt(x, y, f.angle)
	return cad.Point{X: f.x + rx, Y: f.y + ry}
}

// draw collects what every layer draws.
func (b *pcb) draw() {
	board := fpPlace{}
	for _, it := range b.n.lists() {
		switch it.name {
		case "footprint", "module":
			b.footprint(it)
		case "segment":
			b.track(it)
		case "arc":
			b.trackArc(it)
		case "via":
			b.via(it)
		case "zone":
			b.zone(it)
		case "gr_line", "gr_rect", "gr_circle", "gr_arc", "gr_poly", "gr_curve", "gr_bbox":
			b.graphic(it, board, nil)
		case "gr_text":
			b.boardText(it, board, nil, false)
		case "gr_text_box":
			b.textBox(it, board, nil)
		case "dimension":
			b.dimension(it)
		case "target":
			b.target(it)
		case "image":
			b.image(it)
		case "table":
			b.table(it)
		}
	}
}

// strokeWidth reads the width of a board graphic: (stroke (width w)) or,
// before KiCad 7, (width w).
func strokeWidth(it *node) (float64, string) {
	if st := it.child("stroke"); st != nil {
		t := st.str("type")
		if t == "" {
			t = "solid"
		}
		return st.numOf("width", 0), t
	}
	return it.numOf("width", 0), "solid"
}

func filled(it *node) bool {
	f := it.child("fill")
	if f == nil {
		return false
	}
	switch f.arg(0) {
	case "yes", "solid", "true":
		return true
	}
	return f.str("type") == "solid"
}

// graphic draws a line, rectangle, circle, arc, polygon or curve of the
// board or of a footprint (fp not nil).
func (b *pcb) graphic(it *node, pl fpPlace, fp *node) {
	layers := layerNames(it)
	if len(layers) == 0 {
		return
	}
	pt := func(n *node) cad.Point { return pl.pt(n.num(0), n.num(1)) }
	path := &cad.Path{}
	closed := false
	switch strings.TrimPrefix(strings.TrimPrefix(it.name, "gr_"), "fp_") {
	case "line":
		path.Polyline([]cad.Point{pt(it.child("start")), pt(it.child("end"))}, false)
	case "rect", "bbox":
		s, e := it.child("start"), it.child("end")
		pts := []cad.Point{pl.pt(s.num(0), s.num(1)), pl.pt(e.num(0), s.num(1)), pl.pt(e.num(0), e.num(1)), pl.pt(s.num(0), e.num(1))}
		if r := it.numOf("radius", 0); r > 0 {
			path = roundRect(func(q cad.Point) cad.Point { return pl.pt(q.X, q.Y) }, s, e, r)
		} else {
			path.Polyline(pts, true)
		}
		closed = true
	case "circle":
		c := pt(it.child("center"))
		e := it.child("end")
		cc := it.child("center")
		r := math.Hypot(e.num(0)-cc.num(0), e.num(1)-cc.num(1))
		if r <= 0 {
			return
		}
		path.Circle(c, r)
		closed = true
	case "arc":
		s, m, e := it.child("start"), it.child("mid"), it.child("end")
		if m != nil {
			c, r, a0, a1, ok := arcThrough(pt(s), pt(m), pt(e))
			if !ok {
				path.Polyline([]cad.Point{pt(s), pt(e)}, false)
			} else {
				path.Arc(c, r, a0, a1)
			}
		} else if a := it.child("angle"); a != nil {
			// KiCad 6 and earlier: the center, a start point and the sweep
			c, p := pt(s), pt(e)
			r := math.Hypot(p.X-c.X, p.Y-c.Y)
			a0 := math.Atan2(p.Y-c.Y, p.X-c.X)
			path.Arc(c, r, a0, a0+a.num(0)*math.Pi/180)
		} else {
			return
		}
	case "poly":
		var pts []cad.Point
		for _, q := range it.child("pts").lists() {
			if q.name == "xy" {
				pts = append(pts, pt(q))
			} else if q.name == "arc" {
				s, m, e := q.child("start"), q.child("mid"), q.child("end")
				if c, r, a0, a1, ok := arcThrough(pt(s), pt(m), pt(e)); ok {
					for i := 0; i <= 16; i++ {
						a := a0 + (a1-a0)*float64(i)/16
						pts = append(pts, cad.Point{X: c.X + r*math.Cos(a), Y: c.Y + r*math.Sin(a)})
					}
				}
			}
		}
		if len(pts) < 2 {
			return
		}
		path.Polyline(pts, true)
		closed = true
	case "curve":
		var pts []cad.Point
		for _, q := range it.child("pts").children("xy") {
			pts = append(pts, pt(q))
		}
		if len(pts) != 4 {
			return
		}
		path.MoveTo(pts[0].X, pts[0].Y)
		path.CubicTo(pts[1], pts[2], pts[3])
	default:
		return
	}
	w, dash := strokeWidth(it)
	for _, ln := range b.expandLayers(layers) {
		l := b.layer(ln)
		color := pcbLayerColor(ln)
		if closed && filled(it) {
			l.d.Fill(path, cad.Fill{Color: color}, false)
		}
		if w > 0 || !(closed && filled(it)) {
			p := pen(w, dash, color)
			if w <= 0 {
				p.WorldWidth = 0
			}
			l.d.Stroke(path, p)
		}
	}
}

// track draws a track segment.
func (b *pcb) track(it *node) {
	s, e := it.child("start"), it.child("end")
	w := it.numOf("width", 0.25)
	ln := it.str("layer")
	if ln == "" {
		return
	}
	path := (&cad.Path{}).Polyline([]cad.Point{{X: s.num(0), Y: s.num(1)}, {X: e.num(0), Y: e.num(1)}}, false)
	b.layer(ln).d.Stroke(path, cad.Pen{Color: pcbLayerColor(ln), WorldWidth: w, Cap: cad.CapRound, Join: cad.JoinRound})
}

// trackArc draws an arc of track.
func (b *pcb) trackArc(it *node) {
	s, m, e := it.child("start"), it.child("mid"), it.child("end")
	ln := it.str("layer")
	if s == nil || m == nil || e == nil || ln == "" {
		return
	}
	w := it.numOf("width", 0.25)
	p := func(n *node) cad.Point { return cad.Point{X: n.num(0), Y: n.num(1)} }
	path := &cad.Path{}
	if c, r, a0, a1, ok := arcThrough(p(s), p(m), p(e)); ok {
		path.Arc(c, r, a0, a1)
	} else {
		path.Polyline([]cad.Point{p(s), p(e)}, false)
	}
	b.layer(ln).d.Stroke(path, cad.Pen{Color: pcbLayerColor(ln), WorldWidth: w, Cap: cad.CapRound, Join: cad.JoinRound})
}

// viaSpan is a via and the copper layers it joins.
type viaSpan struct {
	c           cad.Point
	size, drill float64
	color       bdf.Color
	from, to    string
}

// drawVia draws a via's ring and hole.
func drawVia(d *cad.Drawing, v viaSpan) {
	d.Fill((&cad.Path{}).Circle(v.c, v.size/2), cad.Fill{Color: v.color}, false)
	if v.drill > 0 {
		d.Fill((&cad.Path{}).Circle(v.c, v.drill/2), cad.Fill{Color: pcbHoleBackground}, false)
	}
}

// holesFor returns what is drawn over a copper layer's view: the pad
// holes and the vias that reach the layer.
func (b *pcb) holesFor(layer string) *cad.Drawing {
	d := &cad.Drawing{}
	d.Items = append(d.Items, b.holes.Items...)
	rank := copperRank(layer)
	for _, v := range b.vias {
		lo, hi := copperRank(v.from), copperRank(v.to)
		if lo > hi {
			lo, hi = hi, lo
		}
		if v.from == "" || (rank >= lo && rank <= hi) {
			drawVia(d, v)
		}
	}
	return d
}

// via records a via: its ring in the via color and its hole.
func (b *pcb) via(it *node) {
	x, y := it.child("at").xy()
	size := it.numOf("size", 0.8)
	drill := it.numOf("drill", 0.4)
	color := pcbViaThrough
	switch {
	case it.hasAtom("blind") || it.child("type").arg(0) == "blind":
		color = pcbViaBlind
	case it.hasAtom("micro") || it.child("type").arg(0) == "micro":
		color = pcbViaMicro
	}
	v := viaSpan{c: cad.Point{X: x, Y: y}, size: size, drill: drill, color: color}
	if ls := it.child("layers"); ls != nil {
		v.from, v.to = ls.arg(0), ls.arg(1)
	}
	b.vias = append(b.vias, v)
}

// zone draws the filled areas of a zone, as KiCad computed and saved them.
func (b *pcb) zone(it *node) {
	for _, fp := range it.children("filled_polygon") {
		ln := fp.str("layer")
		if ln == "" {
			if ls := layerNames(it); len(ls) > 0 {
				ln = ls[0]
			}
		}
		var pts []cad.Point
		for _, q := range fp.child("pts").children("xy") {
			pts = append(pts, cad.Point{X: q.num(0), Y: q.num(1)})
		}
		if len(pts) < 3 || ln == "" {
			continue
		}
		color := pcbLayerColor(ln)
		r, g, bl, a := rgba(color)
		color = bdf.RGBA(r, g, bl, uint8(float64(a)*zoneOpacity))
		b.layer(ln).d.Fill((&cad.Path{}).Polyline(pts, true), cad.Fill{Color: color}, false)
	}
}

// target draws a layer alignment target: a circle and a cross.
func (b *pcb) target(it *node) {
	x, y := it.child("at").xy()
	size := it.numOf("size", 5)
	w := it.numOf("width", 0.1)
	ln := it.str("layer")
	if ln == "" {
		return
	}
	color := pcbLayerColor(ln)
	d := b.layer(ln).d
	pn := pen(w, "solid", color)
	r := size / 3
	if it.hasAtom("x") {
		r = size / 2
		d.Stroke((&cad.Path{}).Polyline([]cad.Point{{X: x - r, Y: y - r}, {X: x + r, Y: y + r}}, false), pn)
		d.Stroke((&cad.Path{}).Polyline([]cad.Point{{X: x + r, Y: y - r}, {X: x - r, Y: y + r}}, false), pn)
		return
	}
	d.Stroke((&cad.Path{}).Circle(cad.Point{X: x, Y: y}, r), pn)
	d.Stroke((&cad.Path{}).Polyline([]cad.Point{{X: x - size/2, Y: y}, {X: x + size/2, Y: y}}, false), pn)
	d.Stroke((&cad.Path{}).Polyline([]cad.Point{{X: x, Y: y - size/2}, {X: x, Y: y + size/2}}, false), pn)
}

// image draws an image placed on a layer.
func (b *pcb) image(it *node) {
	ln := it.str("layer")
	if ln == "" {
		ln = "Dwgs.User"
	}
	p := &schPage{c: b.c, d: b.layer(ln).d}
	p.image(it)
}

var errNothing = fmt.Errorf("kicad: the board draws nothing")
