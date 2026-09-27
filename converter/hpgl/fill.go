package hpgl

import (
	"bytes"
	"fmt"
	"image"
	imgcolor "image/color"
	"image/png"
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// --- the polygon buffer ---

// polygon is the polygon buffer: subpolygons to fill (every vertex) and
// the pieces of their outlines drawn with the pen down (to edge).
type polygon struct {
	// the pen location and state before PM0, which PM2 restores
	savedPos  cad.Point
	savedDown bool
	fill      *cad.Path // closed subpolygons
	sub       *cad.Path // the subpolygon being defined
	edges     *cad.Path // pen-down pieces
}

// start begins a buffer at the pen location.
func (p *polygon) start(at cad.Point) {
	*p = polygon{fill: &cad.Path{}, edges: &cad.Path{}}
	p.sub = (&cad.Path{}).MoveTo(at.X, at.Y)
}

// add adds a piece of path from the last vertex.
func (p *polygon) add(q *cad.Path, down bool) {
	if p.fill == nil {
		p.start(q.First())
	}
	move := p.sub == nil && q.Size() == 2
	switch {
	case move:
		// the point after PM1 starts the next subpolygon: the move to it
		// is not a side
		end := q.Current()
		p.sub = (&cad.Path{}).MoveTo(end.X, end.Y)
	case p.sub == nil:
		// an arc starts at the pen, which is the first point
		p.sub = (&cad.Path{}).Append(q)
	default:
		p.sub.Continue(q)
	}
	if down && !move {
		if p.edges.Open() && p.edges.Current() == q.First() {
			p.edges.Continue(q)
		} else {
			p.edges.Append(q)
		}
	}
}

// transform moves the buffer by m.
func (p *polygon) transform(m canvas.Matrix) {
	for _, q := range []**cad.Path{&p.fill, &p.sub, &p.edges} {
		if *q != nil {
			*q = (*q).Transform(m)
		}
	}
}

// addClosed adds a closed subpolygon (a circle).
func (p *polygon) addClosed(q *cad.Path, down bool) {
	if p.fill == nil {
		p.start(q.First())
		p.sub = nil
	}
	p.closeSub(down)
	p.fill.Append(q)
	p.edges.Append(q)
}

// closeSub closes the subpolygon being defined (PM1, PM2, CI), adding
// its closing side, which is edged when the pen is down.
func (p *polygon) closeSub(down bool) {
	if p.fill == nil {
		return
	}
	if p.sub != nil && !p.sub.Empty() && p.sub.Size() >= 2 {
		first, last := p.sub.First(), p.sub.Current()
		if first != last && down {
			p.edges.Continue((&cad.Path{}).MoveTo(last.X, last.Y).LineTo(first.X, first.Y))
		}
		p.fill.Append(p.sub.Close())
	}
	p.sub = nil
}

func (g *gl) polygonMode(a []float64) {
	n := 0
	if len(a) > 0 {
		n = int(a[0])
	}
	switch n {
	case 0:
		g.flush()
		g.polygon.start(g.pos)
		g.polygon.savedPos, g.polygon.savedDown = g.pos, g.down
		g.inPoly = true
	case 1:
		if g.inPoly {
			g.polygon.closeSub(g.down)
		}
	case 2:
		if g.inPoly {
			g.polygon.closeSub(g.down)
			g.inPoly = false
			g.pos, g.down = g.polygon.savedPos, g.polygon.savedDown
			g.stroke.pendingDot = false
		}
	}
}

// edgePolygon is EP: the pen-down sides of the buffer are drawn.
func (g *gl) edgePolygon() {
	if g.inPoly || g.polygon.edges == nil || g.polygon.edges.Empty() {
		return
	}
	g.flush()
	pen := g.pen()
	lt := g.line.lt
	if lt.set && lt.typ < 0 {
		g.adaptive(g.polygon.edges, lt)
		return
	}
	if lt.set && lt.typ == 0 {
		return
	}
	g.drawStroke(g.polygon.edges, pen)
}

// fillPolygon is FP: the buffer is filled with the fill type, even-odd
// (0) or non-zero (1).
func (g *gl) fillPolygon(a []float64) {
	if g.inPoly || g.polygon.fill == nil || g.polygon.fill.Empty() {
		return
	}
	g.flush()
	g.fillPath(g.polygon.fill, !(len(a) > 0 && a[0] == 1))
}

// rectangle is EA, ER, RA and RR: a rectangle from the pen location to a
// corner, into the buffer, then edged or filled.
func (g *gl) rectangle(a []float64, rel, fill bool) {
	if len(a) < 2 || !valid(a...) || g.inPoly {
		return
	}
	u := g.user()
	q := cad.Point{X: a[0], Y: a[1]}
	if rel {
		q = u.Add(q)
	}
	p := (&cad.Path{}).Polyline([]cad.Point{u, {X: q.X, Y: u.Y}, q, {X: u.X, Y: q.Y}}, true).Transform(g.um)
	g.flush()
	g.polygon.start(g.pos)
	g.polygon.fill, g.polygon.edges, g.polygon.sub = p, p, nil
	if fill {
		g.fillPath(p, true)
	} else {
		g.edgePolygon()
	}
}

// wedge is EW and WG: a wedge of a circle around the pen location from a
// start angle through a sweep, into the buffer, then edged or filled.
func (g *gl) wedge(a []float64, fill bool) {
	if len(a) < 3 || !valid(a...) || g.inPoly {
		return
	}
	u := g.user()
	r := a[0]
	start := a[1] * math.Pi / 180
	sweep := math.Max(-360, math.Min(360, a[2])) * math.Pi / 180
	if r < 0 {
		r = -r
		start += math.Pi
	}
	chord, explicit := g.chordParam(a, 3)
	var p *cad.Path
	if math.Abs(sweep) >= 2*math.Pi-1e-12 {
		p = g.arcPath(u, r, start, start+sweep, chord, explicit).Close()
	} else {
		p = (&cad.Path{}).MoveTo(u.X, u.Y)
		arc := g.arcPath(u, r, start, start+sweep, chord, explicit)
		p.Continue(arc)
		p.Close()
	}
	d := p.Transform(g.um)
	g.flush()
	g.polygon.start(g.pos)
	g.polygon.fill, g.polygon.edges, g.polygon.sub = d, d, nil
	if r == 0 {
		g.dot(g.pos)
		return
	}
	if fill {
		g.fillPath(d, true)
	} else {
		g.edgePolygon()
	}
}

// --- fill types ---

type fillState struct {
	typ         int
	spacing     [5]float64 // FT 3, 4: spacing (plotter units; 0: 1% of the P1–P2 diagonal)
	angle       [5]float64
	shade       float64 // FT 10
	rfIndex     int     // FT 11
	rfPen       bool
	pclPattern  int // FT 21
	transparent bool
	sv          int // SV: 0 solid, 1 shaded, 2 raster fill, 21 PCL pattern
	svShade     float64
	patterns    map[int]*rasterPattern
	images      map[patternKey]bdf.Hash // the images of the patterns drawn
}

// patternKey identifies the image of a raster pattern in the colors it is
// drawn in.
type patternKey struct {
	pat         *rasterPattern
	colors      string // the colors of its pens
	transparent bool
}

// rasterPattern is an RF pattern: pens by pixel, left to right and top to
// bottom.
type rasterPattern struct {
	w, h int
	pens []int
}

func (f *fillState) defaults() {
	*f = fillState{typ: 1, shade: 50, rfIndex: 1, transparent: true, svShade: 50}
}

func (g *gl) fillType(a []float64) {
	g.flush()
	f := &g.fill
	t := 1
	if len(a) > 0 {
		t = int(a[0])
	}
	if !valid(a...) {
		return
	}
	switch t {
	case 1, 2:
	case 3, 4:
		if len(a) > 1 {
			if a[1] < 0 {
				return
			}
			// user units along the x axis, kept in plotter units
			f.spacing[t] = a[1] * math.Abs(g.userScaleX())
		}
		if len(a) > 2 {
			f.angle[t] = a[2]
		}
	case 10:
		if len(a) > 1 {
			f.shade = math.Max(0, math.Min(100, a[1]))
		}
	case 11:
		if len(a) > 1 {
			f.rfIndex = int(a[1])
		}
		f.rfPen = len(a) > 2 && a[2] == 1
	case 21:
		if len(a) > 1 {
			f.pclPattern = int(a[1])
		}
	case 22:
		g.c.warnOnce("FT22", "user-defined PCL fill patterns (FT 22) are filled solid")
	default:
		return
	}
	f.typ = t
}

// userScaleX returns the plotter units of a user unit along the x axis.
func (g *gl) userScaleX() float64 {
	if g.scale == nil {
		return 1
	}
	m, ok := g.scale.matrix(g.p1, g.p2)
	if !ok {
		return 1
	}
	return m[0]
}

func (g *gl) rasterFill(a []float64) {
	g.flush()
	f := &g.fill
	if f.patterns == nil {
		f.patterns = map[int]*rasterPattern{}
	}
	if len(a) == 0 {
		clear(f.patterns)
		return
	}
	n := int(a[0])
	if len(a) < 3 {
		delete(f.patterns, n)
		return
	}
	w, h := int(a[1]), int(a[2])
	if w < 1 || h < 1 || w > 256 || h > 256 {
		return
	}
	pat := &rasterPattern{w: w, h: h, pens: make([]int, w*h)}
	for i := 0; i < w*h && 3+i < len(a); i++ {
		pat.pens[i] = max(0, int(a[3+i]))
	}
	f.patterns[n] = pat
}

func (g *gl) anchorCorner(a []float64) {
	g.flush()
	if len(a) < 2 || !valid(a...) {
		g.anchored = false
		return
	}
	g.anchor = apply(g.um, cad.Point{X: a[0], Y: a[1]})
	g.anchored = true
}

// anchorPoint returns the anchor corner of fill patterns (drawing).
func (g *gl) anchorPoint() cad.Point {
	if g.anchored {
		return g.anchor
	}
	return apply(g.dev, cad.Point{})
}

func (g *gl) screenedVectors(a []float64) {
	g.flush()
	f := &g.fill
	f.sv = 0
	if len(a) == 0 {
		return
	}
	switch t := int(a[0]); t {
	case 0:
	case 1:
		f.sv = 1
		if len(a) > 1 && valid(a[1]) {
			f.svShade = math.Max(0, math.Min(100, a[1]))
		}
	case 2, 21, 22:
		f.sv = 1
		f.svShade = 50
		g.c.warnOnce("SV", "vectors screened with patterns are drawn half shaded")
	}
	g.pens.valid = false
}

// screen applies the screening of vectors to a pen color.
func (f *fillState) screen(c bdf.Color, transparent bool) bdf.Color {
	if f.sv == 0 {
		return c
	}
	return shaded(c, f.svShade, transparent)
}

// shaded returns a color at a shading level (percent): transparent where
// the white of the pattern shows through, or mixed with it.
func shaded(c bdf.Color, level float64, transparent bool) bdf.Color {
	k := math.Max(0, math.Min(1, level/100))
	if transparent {
		return bdf.Color(uint32(c)&0xffffff00 | uint32(math.Round(float64(uint8(c))*k)))
	}
	mix := func(v uint8) uint8 { return uint8(math.Round(float64(v)*k + 255*(1-k))) }
	return bdf.RGB(mix(uint8(c>>24)), mix(uint8(c>>16)), mix(uint8(c>>8)))
}

// fillPath fills a path (drawing units) with the fill type.
func (g *gl) fillPath(p *cad.Path, evenOdd bool) {
	if p.Empty() || !g.c.budget() {
		return
	}
	c := g.c
	n := g.pens.current()
	col := g.penColorOf(n)
	f := &g.fill
	d := &c.cur.d
	c.cur.marked = true
	switch f.typ {
	case 3, 4:
		g.hatch(p, f.typ)
		return
	case 10:
		d.Fill(p, cad.Fill{Color: shaded(col, f.shade, f.transparent)}, evenOdd)
		return
	case 11:
		if img, ok := g.patternImage(f.rfIndex, f.rfPen); ok {
			d.Fill(p, cad.Fill{Pattern: img}, evenOdd)
			return
		}
	case 21:
		g.pclHatch(p, f.pclPattern)
		return
	}
	d.Fill(p, cad.Fill{Color: col}, evenOdd)
}

// hatch fills with parallel lines (FT 3) or cross-hatching (FT 4), drawn
// with the pen and line type from the anchor corner.
func (g *gl) hatch(p *cad.Path, typ int) {
	f := &g.fill
	spacing := f.spacing[typ] * g.scaleOf()
	if spacing <= 0 {
		spacing = 0.01 * g.diagonal() * g.scaleOf()
	}
	if spacing <= 0 {
		return
	}
	angles := []float64{f.angle[typ]}
	if typ == 4 {
		angles = append(angles, f.angle[typ]+90)
	}
	pen := g.pen()
	pen.Dash, pen.DashOffset = nil, 0
	// the dashes of a fixed line type start on a line through the anchor
	// corner, every other line staggered by half the pattern (adaptive
	// types are drawn solid in fills)
	elems := g.dashElements()
	plen := 0.0
	for _, v := range elems {
		plen += math.Abs(v)
	}
	anchor := g.anchorPoint()
	var lines []cad.PatternLine
	for _, a := range angles {
		// the angle is measured from the plotter x axis (turned by RO)
		s, co := math.Sincos(a * math.Pi / 180)
		dir := linear(g.dev, cad.Point{X: co, Y: s})
		t := math.Atan2(dir.Y, dir.X)
		u := cad.Point{X: math.Cos(t), Y: math.Sin(t)}
		nrm := cad.Point{X: -u.Y, Y: u.X}
		if elems != nil {
			lines = append(lines,
				cad.PatternLine{Angle: t, Base: anchor, Offset: nrm.Mul(2 * spacing), Dash: elems},
				cad.PatternLine{Angle: t, Base: anchor.Add(nrm.Mul(spacing)).Add(u.Mul(plen / 2)), Offset: nrm.Mul(2 * spacing), Dash: elems})
		} else {
			lines = append(lines, cad.PatternLine{Angle: t, Base: anchor, Offset: nrm.Mul(spacing)})
		}
	}
	hp := pen
	if !g.c.cur.d.Hatch(p, lines, hp) {
		g.c.warnOnce("hatch", "hatching too fine to draw is filled in a lighter color")
		g.c.cur.d.Fill(p, cad.Fill{Color: shaded(pen.Color, 30, true)}, true)
	}
}

// pclHatch fills with a predefined PCL cross-hatch pattern (FT 21): lines
// 1/300 inch wide, 16/300 inch apart.
func (g *gl) pclHatch(p *cad.Path, n int) {
	var angles []float64
	switch n {
	case 1:
		angles = []float64{0}
	case 2:
		angles = []float64{90}
	case 3:
		angles = []float64{45}
	case 4:
		angles = []float64{135}
	case 5:
		angles = []float64{0, 90}
	case 6:
		angles = []float64{45, 135}
	default:
		g.c.cur.d.Fill(p, cad.Fill{Color: g.penColorOf(g.pens.current())}, true)
		return
	}
	k := g.scaleOf()
	spacing := 16.0 / 300 * 1016 * k
	pen := cad.Pen{Color: g.penColorOf(g.pens.current()), WorldWidth: 1016.0 / 300 * k}
	var lines []cad.PatternLine
	for _, a := range angles {
		s, co := math.Sincos(a * math.Pi / 180)
		dir := linear(g.base, cad.Point{X: co, Y: s})
		t := math.Atan2(dir.Y, dir.X)
		u := cad.Point{X: math.Cos(t), Y: math.Sin(t)}
		lines = append(lines, cad.PatternLine{Angle: t, Base: g.anchorPoint(), Offset: cad.Point{X: -u.Y, Y: u.X}.Mul(spacing)})
	}
	if !g.c.cur.d.Hatch(p, lines, pen) {
		g.c.cur.d.Fill(p, cad.Fill{Color: shaded(pen.Color, 10, true)}, true)
	}
}

// patternImage returns the image pattern of an RF raster fill: a pixel is
// a dot of a 300 dpi device, anchored at the anchor corner.
func (g *gl) patternImage(index int, penColor bool) (*cad.Pattern, bool) {
	pat := g.fill.patterns[index]
	if pat == nil {
		return nil, false
	}
	twoColor := true
	for _, p := range pat.pens {
		if p > 1 {
			twoColor = false
			break
		}
	}
	// the color of each pen the pattern uses (the palette may have changed)
	color := func(pn int) bdf.Color {
		switch {
		case pn == 0:
			return bdf.RGB(255, 255, 255)
		case twoColor && penColor:
			return g.penColorOf(g.pens.current())
		}
		return g.penColorOf(pn)
	}
	var key strings.Builder
	var used [256]bool
	for _, pn := range pat.pens {
		if pn < 256 && !used[pn] {
			used[pn] = true
			fmt.Fprintf(&key, "%d:%08x,", pn, uint32(color(pn)))
		}
	}
	k := patternKey{pat: pat, colors: key.String(), transparent: g.fill.transparent}
	h, ok := g.fill.images[k]
	if !ok {
		img := image.NewNRGBA(image.Rect(0, 0, pat.w, pat.h))
		for i, pn := range pat.pens {
			if pn == 0 && g.fill.transparent {
				continue
			}
			c := color(pn)
			img.SetNRGBA(i%pat.w, i/pat.w, imgcolor.NRGBA{uint8(c >> 24), uint8(c >> 16), uint8(c >> 8), 255})
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return nil, false
		}
		h = g.c.doc.AddImage(buf.Bytes())
		if g.fill.images == nil {
			g.fill.images = map[patternKey]bdf.Hash{}
		}
		g.fill.images[k] = h
	}
	// rows run along the plotter x axis, the first row at the top
	dot := 1016.0 / 300
	a := g.anchorPoint()
	ex := linear(g.dev, cad.Point{X: dot})
	ey := linear(g.dev, cad.Point{Y: -dot})
	return &cad.Pattern{Image: h, M: canvas.Matrix{ex.X, ex.Y, ey.X, ey.Y, a.X, a.Y}, Smooth: true}, true
}

// --- the window (IW) ---

// window is a soft-clip window, a rectangle of the drawing.
type window struct {
	path *cad.Path
	open bool
}

func (g *gl) inputWindow(a []float64) {
	g.flush()
	g.closeWindow()
	g.window = nil
	if len(a) >= 4 && valid(a...) {
		x0, x1 := math.Min(a[0], a[2]), math.Max(a[0], a[2])
		y0, y1 := math.Min(a[1], a[3]), math.Max(a[1], a[3])
		p := (&cad.Path{}).Polyline([]cad.Point{{X: x0, Y: y0}, {X: x1, Y: y0}, {X: x1, Y: y1}, {X: x0, Y: y1}}, true)
		g.window = &window{path: p.Transform(g.um)}
	}
	g.openWindow()
}

// openWindow starts drawing clipped by the window.
func (g *gl) openWindow() {
	if g.window != nil && !g.window.open {
		g.c.cur.d.Begin(g.window.path)
		g.window.open = true
	}
}

// closeWindow ends the drawing clipped by the window.
func (g *gl) closeWindow() {
	if g.window != nil && g.window.open {
		g.flush()
		g.c.cur.d.End()
		g.window.open = false
	}
}
