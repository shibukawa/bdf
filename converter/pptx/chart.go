package pptx

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
)

// Charts are drawn from the values cached in the chart part (c:numCache,
// c:strCache): bar and column charts (clustered, stacked, 100% stacked),
// line, area, pie and doughnut, and scatter charts, with axes, gridlines,
// tick labels, the legend, the title and simple data labels. 3-D variants
// are drawn flat.

type chartSeries struct {
	idx    int
	name   string
	cats   []string
	vals   []float64 // NaN = no value
	xs     []float64 // scatter X values
	spPr   *node
	dPts   map[int]*node
	marker *node
	dLbls  *node
	format string
}

type chartPlot struct {
	kind       string // bar, line, area, pie, doughnut, scatter
	horizontal bool   // bar charts with barDir="bar"
	grouping   string
	vary       bool
	series     []*chartSeries
	gap        float64
	overlap    float64
	hole       float64
	firstAng   float64
	dLbls      *node
	markers    bool
	node       *node
}

type chartCtx struct {
	s          *slideCtx
	cv         *canvas
	part       string
	root       *node
	x, y, w, h float64
	base       *node // chart-wide default run properties
	nseries    int
}

func (s *slideCtx) drawChart(cv *canvas, sh *shape, xf xform, rid string) {
	r, ok := s.c.pkg.target(sh.part, rid)
	if !ok {
		return
	}
	root, err := s.c.pkg.xml(r.Target)
	if err != nil {
		s.c.warnf("chart: %v", err)
		return
	}
	ch := &chartCtx{s: s, cv: cv, part: r.Target, root: root, x: xf.X, y: xf.Y, w: xf.W, h: xf.H}
	ch.base = root.path("txPr", "p", "pPr", "defRPr")
	ch.draw()
	cv.drawn = true
}

func (ch *chartCtx) draw() {
	s, cv := ch.s, ch.cv
	cv.obj.Save()
	cv.obj.Translate(f32(ch.x), f32(ch.y))
	// chart area
	s.boxFill(cv, ch.root.child("spPr"), ch.part, 0, 0, ch.w, ch.h)
	chart := ch.root.child("chart")
	pa := chart.child("plotArea")
	var plots []*chartPlot
	for _, k := range pa.Kids {
		if p := ch.parsePlot(k); p != nil {
			plots = append(plots, p)
		}
	}
	// title
	top, bottom, left, right := 7.0, ch.h-7, 7.0, ch.w-7
	autoDeleted := chart.child("autoTitleDeleted").attrStr("val", "0") == "1"
	if t := chart.child("title"); t != nil && !autoDeleted {
		top += ch.title(t, plots, top)
	} else if t == nil && !autoDeleted && len(plots) == 1 && len(plots[0].series) == 1 {
		// Office titles a single-series chart with the series name
		top += ch.title(&node{Name: "title"}, plots, top)
	}
	// legend
	if lg := chart.child("legend"); lg != nil && len(plots) > 0 {
		top, bottom, left, right = ch.legend(lg, plots, top, bottom, left, right)
	}
	if len(plots) == 0 {
		cv.obj.Restore()
		return
	}
	// plot area
	px, py, pw, ph := left, top, right-left, bottom-top
	manual := false
	if ml := pa.path("layout", "manualLayout"); ml != nil && ml.child("w") != nil {
		px, py = ml.child("x").attrFloat("val", 0)*ch.w, ml.child("y").attrFloat("val", 0)*ch.h
		pw, ph = ml.child("w").attrFloat("val", 1)*ch.w, ml.child("h").attrFloat("val", 1)*ch.h
		manual = ml.child("layoutTarget").attrStr("val", "outer") == "inner"
	}
	switch plots[0].kind {
	case "pie", "doughnut":
		s.boxFill(cv, pa.child("spPr"), ch.part, px, py, pw, ph)
		ch.pie(plots[0], px, py, pw, ph)
	default:
		ch.axesChart(pa, plots, px, py, pw, ph, manual)
	}
	cv.obj.Restore()
}

// boxFill fills and outlines a rectangle with an spPr.
func (s *slideCtx) boxFill(cv *canvas, spPr *node, part string, x, y, w, h float64) {
	if spPr == nil {
		return
	}
	f := s.resolveFill(fillElem(spPr), part, s.cc, nil)
	var ln *line
	if l := spPr.child("ln"); l != nil {
		ln = s.resolveLine([]*node{l}, []string{part}, nil, s.cc)
	}
	if f.kind == fillNone && ln == nil {
		return
	}
	cv.obj.Save()
	cv.obj.Translate(f32(x), f32(y))
	s.paintGeometry(cv, rectGeometry(w, h), f, ln, nil, w, h)
	cv.obj.Restore()
}

func (ch *chartCtx) parsePlot(k *node) *chartPlot {
	kind := strings.TrimSuffix(k.Name, "Chart")
	switch kind {
	case "bar", "bar3D":
		kind = "bar"
	case "line", "line3D", "stock":
		kind = "line"
	case "area", "area3D":
		kind = "area"
	case "pie", "pie3D", "ofPie":
		kind = "pie"
	case "doughnut", "scatter":
	case "radar", "bubble", "surface", "surface3D":
		ch.s.c.warnOnce("chart:"+kind, "%s charts are not supported", kind)
		return nil
	default:
		return nil
	}
	p := &chartPlot{kind: kind, node: k, grouping: k.child("grouping").attrStr("val", "clustered"),
		vary: k.child("varyColors").attrStr("val", "0") == "1", gap: 150, dLbls: k.child("dLbls")}
	if kind == "bar" {
		p.horizontal = k.child("barDir").attrStr("val", "col") == "bar"
		if g := k.child("gapWidth"); g != nil {
			p.gap = g.attrFloat("val", 150)
		}
		if o := k.child("overlap"); o != nil {
			p.overlap = o.attrFloat("val", 0)
		} else if p.grouping == "stacked" || p.grouping == "percentStacked" {
			p.overlap = 100
		}
	}
	if kind == "line" && k.child("grouping") == nil {
		p.grouping = "standard"
	}
	if kind == "doughnut" {
		p.hole = k.child("holeSize").attrFloat("val", 50) / 100
	}
	if kind == "pie" || kind == "doughnut" {
		p.firstAng = k.child("firstSliceAng").attrFloat("val", 0)
	}
	p.markers = kind == "line" && k.child("marker").attrStr("val", "1") != "0" || kind == "scatter"
	for _, sn := range k.children("ser") {
		se := &chartSeries{idx: ch.nseries, spPr: sn.child("spPr"), dPts: map[int]*node{}, marker: sn.child("marker"), dLbls: sn.child("dLbls")}
		if v := sn.child("idx"); v != nil {
			se.idx = int(v.attrInt("val", int64(ch.nseries)))
		}
		ch.nseries++
		se.name = strings.Join(cacheStrings(sn.child("tx")), " ")
		se.cats = cacheStrings(sn.child("cat"))
		var vals *node
		if kind == "scatter" {
			se.xs = cacheNumbers(sn.child("xVal"), nil)
			numeric := false
			for _, x := range se.xs {
				if !math.IsNaN(x) {
					numeric = true
				}
			}
			if !numeric {
				// text X values: the points are numbered 1, 2, 3 …
				se.xs = nil
			}
			se.cats = cacheStrings(sn.child("xVal"))
			vals = sn.child("yVal")
		} else {
			vals = sn.child("val")
		}
		se.vals = cacheNumbers(vals, &se.format)
		for _, d := range sn.children("dPt") {
			se.dPts[int(d.child("idx").attrInt("val", -1))] = d.child("spPr")
		}
		p.series = append(p.series, se)
	}
	return p
}

// cacheStrings returns the cached strings of a c:tx / c:cat / c:xVal.
func cacheStrings(n *node) []string {
	if n == nil {
		return nil
	}
	if v := n.child("v"); v != nil {
		return []string{v.text()}
	}
	var cache *node
	for _, ref := range []string{"strRef", "numRef", "multiLvlStrRef", "strLit", "numLit"} {
		if r := n.child(ref); r != nil {
			cache = r
			if c := r.child("strCache"); c != nil {
				cache = c
			} else if c := r.child("numCache"); c != nil {
				cache = c
			} else if c := r.child("multiLvlStrCache"); c != nil {
				cache = c.child("lvl")
			}
			break
		}
	}
	if cache == nil {
		return nil
	}
	count := int(cache.child("ptCount").attrInt("val", 0))
	pts := cache.children("pt")
	for _, pt := range pts {
		count = max(count, int(pt.attrInt("idx", 0))+1)
	}
	out := make([]string, count)
	for _, pt := range pts {
		if i := int(pt.attrInt("idx", -1)); i >= 0 && i < count {
			out[i] = pt.child("v").text()
		}
	}
	return out
}

// cacheNumbers returns cached numbers (NaN for missing points) and the format code.
func cacheNumbers(n *node, format *string) []float64 {
	if n == nil {
		return nil
	}
	var cache *node
	if r := n.child("numRef"); r != nil {
		cache = r.child("numCache")
	} else if r := n.child("numLit"); r != nil {
		cache = r
	}
	if cache == nil {
		return nil
	}
	if format != nil {
		*format = cache.child("formatCode").text()
	}
	count := int(cache.child("ptCount").attrInt("val", 0))
	pts := cache.children("pt")
	for _, pt := range pts {
		count = max(count, int(pt.attrInt("idx", 0))+1)
	}
	out := make([]float64, count)
	for i := range out {
		out[i] = math.NaN()
	}
	for _, pt := range pts {
		if i := int(pt.attrInt("idx", -1)); i >= 0 && i < count {
			if v, err := strconv.ParseFloat(strings.TrimSpace(pt.child("v").text()), 64); err == nil {
				out[i] = v
			}
		}
	}
	return out
}

// autoColor is the automatic series (or point) color: the accents, then
// darker and lighter variants of them.
func (ch *chartCtx) autoColor(i int) rgba {
	c, _ := ch.s.cc.scheme1("accent" + itoa(i%6+1))
	switch (i / 6) % 6 {
	case 1:
		c = lumMod(c, 0.6, 0)
	case 2:
		c = lumMod(c, 0.8, 0.2)
	case 3:
		c = lumMod(c, 0.8, 0)
	case 4:
		c = lumMod(c, 0.6, 0.4)
	case 5:
		c = lumMod(c, 0.5, 0)
	}
	return c
}

func lumMod(c rgba, mod, off float64) rgba {
	h, s, l := rgbToHSL(c)
	out := hslToRGB(h, s, clamp01(l*mod+off))
	out.A = c.A
	return out
}

// seriesFill is the fill of a series (or of one point of it).
func (ch *chartCtx) seriesFill(p *chartPlot, se *chartSeries, pt int) fill {
	if sp := se.dPts[pt]; sp != nil {
		if e := fillElem(sp); e != nil {
			return ch.s.resolveFill(e, ch.part, ch.s.cc, nil)
		}
	}
	if e := fillElem(se.spPr); e != nil {
		return ch.s.resolveFill(e, ch.part, ch.s.cc, nil)
	}
	if p.vary && len(p.series) == 1 {
		return fill{kind: fillSolid, color: ch.autoColor(pt)}
	}
	return fill{kind: fillSolid, color: ch.autoColor(se.idx)}
}

// seriesLine is the outline of a series; lines and scatter series get
// their automatic color when they set none.
func (ch *chartCtx) seriesLine(p *chartPlot, se *chartSeries, pt int) *line {
	var lns []*node
	if sp := se.dPts[pt]; sp != nil && sp.child("ln") != nil {
		lns = append(lns, sp.child("ln"))
	}
	if ln := se.spPr.child("ln"); ln != nil {
		lns = append(lns, ln)
	}
	parts := make([]string, len(lns))
	for i := range parts {
		parts[i] = ch.part
	}
	l := ch.s.resolveLine(lns, parts, nil, ch.s.cc)
	if p.kind == "line" || p.kind == "scatter" {
		explicitNone := false
		for _, ln := range lns {
			if e := fillElem(ln); e != nil {
				explicitNone = e.Name == "noFill"
				break
			}
		}
		if l == nil && !explicitNone {
			l = &line{width: 2.25, fill: fill{kind: fillSolid, color: ch.autoColor(se.idx)}, cap: 1, join: 1, miter: 8}
			for _, ln := range lns {
				if _, ok := ln.attr("w"); ok {
					l.width = ln.emuAttr("w", 2.25)
					break
				}
			}
		}
	}
	return l
}

// textStyle builds the run style of a chart text element from its txPr
// (c:txPr or the default run properties of c:rich) and the chart default.
func (ch *chartCtx) textStyle(txPr *node, size float64, bold bool) *runStyle {
	// size is the element's default, relative to the chart's base size
	def := mkNode("defRPr", "sz", itoa(int(math.Round(size*100))))
	if bold {
		def.Attr = append(def.Attr, mkNode("", "b", "1").Attr...)
	}
	rc := chain{txPr.path("p", "pPr", "defRPr"), def, ch.base}
	tf := &textFrame{part: ch.part, cc: ch.s.cc}
	return ch.s.runStyle(tf, rc, 1, nil)
}

// measure returns the width of a single-line label.
func (ch *chartCtx) measure(st *runStyle, text string) float64 {
	w := 0.0
	for _, r := range text {
		w += ch.s.c.advance(ch.s.c.faceFor(st, r), r)*st.size + st.spacing
	}
	return w
}

// label draws a single-line label; ha is 'l', 'c' or 'r' and va 't', 'm'
// or 'b' for the side of the text box that (x, y) is on.
func (ch *chartCtx) label(st *runStyle, text string, x, y float64, ha, va byte) {
	if text == "" {
		return
	}
	w := ch.measure(st, text)
	fc := ch.s.c.faceFor(st, 'x')
	switch ha {
	case 'c':
		x -= w / 2
	case 'r':
		x -= w
	}
	asc, desc := fc.asc*st.size, fc.desc*st.size
	switch va {
	case 't':
		y += asc
	case 'm':
		y += (asc - desc) / 2
	case 'b':
		y -= desc
	}
	pa := &para{}
	ch.s.addChars(pa, st, text)
	xx := 0.0
	for i := range pa.items {
		pa.items[i].x = xx
		xx += pa.items[i].w
	}
	em := &textEmitter{cv: ch.cv}
	ch.cv.obj.Save()
	ch.cv.obj.Mark(bdf.MarkBox, "")
	em.emitItems(&textLine{items: pa.items}, x, y)
	ch.cv.obj.Restore()
}

func (ch *chartCtx) lineHeight(st *runStyle) float64 {
	fc := ch.s.c.faceFor(st, 'x')
	return (fc.asc + fc.desc) * st.size
}

// title draws the chart title and returns the height it takes.
func (ch *chartCtx) title(t *node, plots []*chartPlot, top float64) float64 {
	if rich := t.path("tx", "rich"); rich != nil {
		tf := &textFrame{body: rich, part: ch.part, cc: ch.s.cc, own: rich.child("lstStyle"), bodyPrs: chain{rich.child("bodyPr"), mkNode("bodyPr", "wrap", "square")}}
		def := mkNode("defRPr", "sz", itoa(int(math.Round(ch.baseSize()*120))), "b", "1")
		tf.extra = []*node{def, ch.base}
		// center the paragraphs
		ctr := mkNode("defPPr", "algn", "ctr")
		tf.lists = []*node{&node{Name: "lstStyle", Kids: []*node{ctr}}}
		b := ch.s.layoutText(tf, ch.w*0.1, top, ch.w*0.8, ch.h, identity)
		if b == nil {
			return 0
		}
		b.h = b.height()
		ch.s.drawTextBlock(ch.cv, b)
		return b.height()
	}
	text := ""
	if v := t.path("tx", "strRef"); v != nil {
		text = strings.Join(cacheStrings(t.child("tx")), " ")
	} else if len(plots) == 1 && len(plots[0].series) == 1 {
		text = plots[0].series[0].name
	} else {
		// Office shows its default title text
		text = "Chart Title"
	}
	if text == "" {
		return 0
	}
	st := ch.textStyle(t.child("txPr"), ch.baseSize()*1.2, true)
	ch.label(st, text, ch.w/2, top, 'c', 't')
	return ch.lineHeight(st) + 4
}

func (ch *chartCtx) baseSize() float64 {
	if v, ok := ch.base.attr("sz"); ok {
		return atof(v, 1000) / 100
	}
	return 10
}

type legendEntry struct {
	text string
	f    fill
	l    *line
	mark bool
}

// legend draws the legend and returns the plot area left for the chart.
func (ch *chartCtx) legend(lg *node, plots []*chartPlot, top, bottom, left, right float64) (float64, float64, float64, float64) {
	var entries []legendEntry
	for _, p := range plots {
		if (p.kind == "pie" || p.kind == "doughnut") && len(p.series) > 0 {
			se := p.series[0]
			for i := range se.vals {
				name := itoa(i + 1)
				if i < len(se.cats) {
					name = se.cats[i]
				}
				entries = append(entries, legendEntry{text: name, f: ch.seriesFill(&chartPlot{vary: true, series: p.series[:1]}, se, i)})
			}
			continue
		}
		for _, se := range p.series {
			e := legendEntry{text: se.name, f: ch.seriesFill(p, se, -1)}
			if p.kind == "line" || p.kind == "scatter" {
				e.l = ch.seriesLine(p, se, -1)
				e.mark = p.markers && se.marker.child("symbol").attrStr("val", "auto") != "none"
				e.f = ch.markerFill(p, se)
			}
			entries = append(entries, e)
		}
	}
	for _, del := range lg.children("legendEntry") {
		if del.child("delete").attrStr("val", "0") == "1" {
			i := int(del.child("idx").attrInt("val", -1))
			if i >= 0 && i < len(entries) {
				entries[i].text = "\x00"
			}
		}
	}
	kept := entries[:0]
	for _, e := range entries {
		if e.text != "\x00" {
			kept = append(kept, e)
		}
	}
	entries = kept
	if len(entries) == 0 {
		return top, bottom, left, right
	}
	st := ch.textStyle(lg.child("txPr"), ch.baseSize(), false)
	lh := ch.lineHeight(st) * 1.2
	key := st.size * 0.7
	pos := lg.child("legendPos").attrStr("val", "r")
	overlay := lg.child("overlay").attrStr("val", "0") == "1"
	draw := func(e legendEntry, x, y float64) {
		cy := y + lh/2
		if e.l != nil || e.mark {
			if e.l != nil {
				ch.cv.obj.Save()
				ch.cv.setLine(e.l, key*2, 0)
				ch.cv.obj.StrokePath(ch.cv.obj.AddPath(newPath().moveTo(x, cy).lineTo(x+key*2, cy).p))
				ch.cv.obj.Restore()
			}
			if e.mark {
				ch.marker(nil, e.f, x+key, cy, "square", key*0.7)
			}
			ch.label(st, e.text, x+key*2+4, cy, 'l', 'm')
			return
		}
		ch.cv.obj.Save()
		if ch.cv.setFill(e.f, key, key) {
			ch.cv.obj.FillRect(f32(x), f32(cy-key/2), f32(key), f32(key))
		}
		ch.cv.obj.Restore()
		ch.label(st, e.text, x+key+4, cy, 'l', 'm')
	}
	keyW := key + 4
	for _, e := range entries {
		if e.l != nil || e.mark {
			keyW = key*2 + 4
		}
	}
	switch pos {
	case "b", "t":
		widths := make([]float64, len(entries))
		total := 0.0
		for i, e := range entries {
			widths[i] = keyW + ch.measure(st, e.text) + 12
			total += widths[i]
		}
		x := (ch.w - total) / 2
		y := bottom - lh
		if pos == "t" {
			y = top
		}
		for i, e := range entries {
			draw(e, x, y)
			x += widths[i]
		}
		if !overlay {
			if pos == "t" {
				top += lh + 4
			} else {
				bottom -= lh + 4
			}
		}
	default: // r, l, tr
		maxW := 0.0
		for _, e := range entries {
			maxW = math.Max(maxW, keyW+ch.measure(st, e.text))
		}
		total := lh * float64(len(entries))
		y := top + (bottom-top-total)/2
		x := right - maxW
		if pos == "l" {
			x = left
		}
		if pos == "tr" {
			y = top
		}
		for _, e := range entries {
			draw(e, x, y)
			y += lh
		}
		if !overlay {
			if pos == "l" {
				left += maxW + 8
			} else {
				right -= maxW + 8
			}
		}
	}
	return top, bottom, left, right
}

func (ch *chartCtx) markerFill(p *chartPlot, se *chartSeries) fill {
	if f := se.marker.path("spPr"); fillElem(f) != nil {
		return ch.s.resolveFill(fillElem(f), ch.part, ch.s.cc, nil)
	}
	if l := ch.seriesLine(p, se, -1); l != nil {
		return l.fill
	}
	return fill{kind: fillSolid, color: ch.autoColor(se.idx)}
}

var autoMarkers = []string{"diamond", "square", "triangle", "x", "star", "circle", "plus", "dot", "dash"}

// marker draws a data point marker centered at (x, y).
func (ch *chartCtx) marker(se *chartSeries, f fill, x, y float64, symbol string, size float64) {
	if se != nil {
		symbol = se.marker.child("symbol").attrStr("val", "auto")
		if symbol == "auto" {
			symbol = autoMarkers[se.idx%len(autoMarkers)]
		}
		size = se.marker.child("size").attrFloat("val", 5) * 1.4
	}
	if symbol == "none" {
		return
	}
	r := size / 2
	p := newPath()
	stroke := false
	switch symbol {
	case "circle", "dot":
		p.p.Circle(f32(x), f32(y), f32(r))
	case "diamond":
		p.moveTo(x, y-r).lineTo(x+r, y).lineTo(x, y+r).lineTo(x-r, y).close()
	case "triangle":
		p.moveTo(x, y-r).lineTo(x+r, y+r).lineTo(x-r, y+r).close()
	case "x", "star", "plus":
		stroke = true
		if symbol != "plus" {
			p.moveTo(x-r, y-r).lineTo(x+r, y+r).moveTo(x-r, y+r).lineTo(x+r, y-r)
		}
		if symbol != "x" {
			p.moveTo(x, y-r).lineTo(x, y+r).moveTo(x-r, y).lineTo(x+r, y)
		}
	case "dash":
		p.p.Rect(f32(x-r), f32(y-r/4), f32(2*r), f32(r/2))
	default: // square
		p.p.Rect(f32(x-r), f32(y-r), f32(2*r), f32(2*r))
	}
	ref := ch.cv.obj.AddPath(p.p)
	ch.cv.obj.Save()
	if stroke {
		c, _ := f.average()
		ch.cv.obj.Line(1, 0, 0, 10)
		ch.cv.obj.StrokeColor(c.bdf())
		ch.cv.obj.StrokePath(ref)
	} else if ch.cv.setFill(f, size, size) {
		ch.cv.obj.FillPath(ref, 0)
	}
	ch.cv.obj.Restore()
}

// axis is a value axis scale.
type axis struct {
	lo, hi, step float64
	node         *node
	format       string
	deleted      bool
	reverse      bool
}

func (a *axis) pos(v float64) float64 {
	if a.hi == a.lo {
		return 0
	}
	t := (v - a.lo) / (a.hi - a.lo)
	if a.reverse {
		t = 1 - t
	}
	return t
}

// niceAxis picks the axis range and major unit like Office does: zero is
// included unless the data are far from it, the ends get 5% of headroom,
// and the unit is the smallest 1, 2 or 5 times a power of ten that needs at
// most maxSteps intervals.
func niceAxis(lo, hi float64, n *node, maxSteps int) *axis {
	a := &axis{node: n}
	if math.IsInf(lo, 0) || math.IsInf(hi, 0) {
		lo, hi = 0, 1
	}
	if lo > 0 && (hi-lo) > hi/6 {
		lo = 0
	}
	if hi < 0 && (hi-lo) > -lo/6 {
		hi = 0
	}
	if lo == hi {
		if lo == 0 {
			hi = 1
		} else {
			lo, hi = math.Min(0, lo), math.Max(0, hi)
		}
	}
	sc := n.child("scaling")
	fixedMin, fixedMax := sc.child("min") != nil, sc.child("max") != nil
	if fixedMin {
		lo = sc.child("min").attrFloat("val", lo)
	}
	if fixedMax {
		hi = sc.child("max").attrFloat("val", hi)
	}
	span := hi - lo
	if !fixedMax && hi > 0 {
		hi += span / 20
	}
	if !fixedMin && lo < 0 {
		lo -= span / 20
	}
	step := n.child("majorUnit").attrFloat("val", 0)
	if step <= 0 {
		maxSteps = min(10, max(maxSteps, 2))
		mag := math.Pow(10, math.Floor(math.Log10((hi-lo)/float64(maxSteps))))
		for _, m := range []float64{1, 2, 5, 10, 20, 50} {
			step = m * mag
			if math.Ceil(hi/step-1e-9)-math.Floor(lo/step+1e-9) <= float64(maxSteps) {
				break
			}
		}
	}
	a.step = step
	a.lo, a.hi = lo, hi
	if !fixedMin {
		a.lo = math.Floor(lo/step+1e-9) * step
	}
	if !fixedMax {
		a.hi = math.Ceil(hi/step-1e-9) * step
	}
	if a.hi <= a.lo {
		a.hi = a.lo + step
	}
	a.deleted = n.child("delete").attrStr("val", "0") == "1"
	a.reverse = sc.child("orientation").attrStr("val", "minMax") == "maxMin"
	a.format = n.child("numFmt").attrStr("formatCode", "General")
	return a
}

// formatNumber formats a value with a (simplified) number format code.
func formatNumber(v float64, code string) string {
	if code == "" || strings.EqualFold(code, "General") {
		s := strconv.FormatFloat(v, 'f', -1, 64)
		if len(s) > 10 {
			s = strconv.FormatFloat(v, 'g', 6, 64)
		}
		return s
	}
	if i := strings.IndexByte(code, ';'); i >= 0 {
		code = code[:i]
	}
	percent := strings.Contains(code, "%")
	if percent {
		v *= 100
	}
	dec := 0
	if i := strings.IndexByte(code, '.'); i >= 0 {
		for _, r := range code[i+1:] {
			if r == '0' || r == '#' {
				dec++
			} else {
				break
			}
		}
	}
	s := strconv.FormatFloat(math.Abs(v), 'f', dec, 64)
	if strings.Contains(code, ",") {
		intPart, frac, _ := strings.Cut(s, ".")
		var b strings.Builder
		for i, r := range intPart {
			if i > 0 && (len(intPart)-i)%3 == 0 {
				b.WriteByte(',')
			}
			b.WriteRune(r)
		}
		s = b.String()
		if frac != "" {
			s += "." + frac
		}
	}
	if v < 0 {
		s = "-" + s
	}
	prefix := ""
	for _, cur := range []string{"$", "¥", "€", "£", "￥"} {
		if strings.Contains(code, cur) {
			prefix = cur
		}
	}
	s = prefix + s
	if percent {
		s += "%"
	}
	return s
}

// axesChart draws bar, line, area and scatter plots on category and value axes.
func (ch *chartCtx) axesChart(pa *node, plots []*chartPlot, px, py, pw, ph float64, inner bool) {
	s, cv := ch.s, ch.cv
	p0 := plots[0]
	horizontal := p0.kind == "bar" && p0.horizontal
	scatter := p0.kind == "scatter"
	// data range
	ncat := 0
	lo, hi := math.Inf(1), math.Inf(-1)
	xlo, xhi := math.Inf(1), math.Inf(-1)
	var cats []string
	for _, p := range plots {
		stacked := p.grouping == "stacked" || p.grouping == "percentStacked"
		pos, neg := map[int]float64{}, map[int]float64{}
		for _, se := range p.series {
			ncat = max(ncat, len(se.vals))
			if len(se.cats) > len(cats) {
				cats = se.cats
			}
			for i, v := range se.vals {
				if math.IsNaN(v) {
					continue
				}
				if p.grouping == "percentStacked" {
					continue
				}
				if stacked {
					if v >= 0 {
						pos[i] += v
						v = pos[i]
					} else {
						neg[i] += v
						v = neg[i]
					}
				}
				lo, hi = math.Min(lo, v), math.Max(hi, v)
				if stacked {
					lo = math.Min(lo, 0)
				}
			}
			if p.kind == "scatter" {
				for i := range se.vals {
					x := float64(i + 1)
					if i < len(se.xs) {
						x = se.xs[i]
					}
					if !math.IsNaN(x) {
						xlo, xhi = math.Min(xlo, x), math.Max(xhi, x)
					}
				}
			}
		}
		if p.grouping == "percentStacked" {
			lo, hi = math.Min(lo, 0), math.Max(hi, 1)
		}
	}
	var valAxNode, catAxNode *node
	for _, k := range pa.Kids {
		switch k.Name {
		case "valAx":
			if scatter {
				// X is the horizontal value axis
				if pos := k.child("axPos").attrStr("val", "l"); (pos == "b" || pos == "t") && catAxNode == nil {
					catAxNode = k
					continue
				}
			}
			if valAxNode == nil {
				valAxNode = k
			}
		case "catAx", "dateAx":
			if catAxNode == nil {
				catAxNode = k
			}
		}
	}
	// fewer ticks on small charts: vertical axes keep 1.5 label heights per
	// interval, horizontal ones 1.5 label widths
	budget := func(length float64, st *runStyle) int { return int(length / (1.5 * ch.lineHeight(st))) }
	hbudget := func(length float64, st *runStyle, lo, hi float64) int {
		w := math.Max(ch.measure(st, formatNumber(hi, "General")), ch.measure(st, formatNumber(lo, "General"))) + st.size/2
		return int(length / (1.5 * w))
	}
	valStyle := ch.textStyle(valAxNode.child("txPr"), ch.baseSize(), false)
	catStyle := ch.textStyle(catAxNode.child("txPr"), ch.baseSize(), false)
	vb := budget(ph, valStyle)
	if horizontal {
		vb = hbudget(pw, valStyle, lo, hi)
	}
	va := niceAxis(lo, hi, valAxNode, vb)
	if p0.grouping == "percentStacked" {
		if valAxNode.path("numFmt") == nil {
			va.format = "0%"
		}
	}
	if va.format == "General" && len(p0.series) > 0 && valAxNode.child("numFmt").attrStr("sourceLinked", "1") != "0" {
		va.format = p0.series[0].format
	}
	var xa *axis
	if scatter {
		xa = niceAxis(xlo, xhi, catAxNode, hbudget(pw, catStyle, xlo, xhi))
	}
	catReverse := catAxNode.path("scaling", "orientation").attrStr("val", "minMax") == "maxMin"
	catDeleted := catAxNode.child("delete").attrStr("val", "0") == "1"
	showVal := !va.deleted && valAxNode.child("tickLblPos").attrStr("val", "nextTo") != "none"
	showCat := !catDeleted && catAxNode.child("tickLblPos").attrStr("val", "nextTo") != "none"
	// tick labels
	var valLabels []string
	var valTicks []float64
	for v := va.lo; v <= va.hi+va.step*1e-6; v += va.step {
		valTicks = append(valTicks, v)
		valLabels = append(valLabels, formatNumber(v, va.format))
	}
	var catLabels []string
	var catTicks []float64
	if scatter {
		for v := xa.lo; v <= xa.hi+xa.step*1e-6; v += xa.step {
			catTicks = append(catTicks, v)
			catLabels = append(catLabels, formatNumber(v, xa.format))
		}
	} else {
		for i := 0; i < ncat; i++ {
			if i < len(cats) && cats[i] != "" {
				catLabels = append(catLabels, catLabel(cats[i], catAxNode))
			} else {
				catLabels = append(catLabels, itoa(i+1))
			}
		}
	}
	// reserve room for the labels unless the layout is the inner plot area
	if !inner {
		vw, vh := 0.0, ch.lineHeight(valStyle)
		for _, l := range valLabels {
			vw = math.Max(vw, ch.measure(valStyle, l))
		}
		cw, chh := 0.0, ch.lineHeight(catStyle)
		for _, l := range catLabels {
			cw = math.Max(cw, ch.measure(catStyle, l))
		}
		if horizontal {
			if showCat {
				px += cw + 6
				pw -= cw + 6
			}
			if showVal {
				ph -= vh + 4
			}
		} else {
			if showVal {
				px += vw + 6
				pw -= vw + 6
			}
			if showCat {
				ph -= chh + 4
			}
		}
	}
	if pw <= 0 || ph <= 0 {
		return
	}
	s.boxFill(cv, pa.child("spPr"), ch.part, px, py, pw, ph)
	// value → coordinate
	vpos := func(v float64) float64 {
		t := va.pos(v)
		if horizontal {
			return px + t*pw
		}
		return py + ph - t*ph
	}
	// gridlines
	gridLine := func(ax *node, name string) *line {
		g := ax.child(name)
		if g == nil {
			return nil
		}
		if ln := g.path("spPr", "ln"); ln != nil {
			return s.resolveLine([]*node{ln}, []string{ch.part}, nil, s.cc)
		}
		return &line{width: 0.75, fill: fill{kind: fillSolid, color: rgba{0.85, 0.85, 0.85, 1}}, miter: 8}
	}
	strokeLine := func(l *line, x0, y0, x1, y1 float64) {
		if l == nil {
			return
		}
		cv.obj.Save()
		cv.setLine(l, math.Abs(x1-x0), math.Abs(y1-y0))
		cv.obj.StrokePath(cv.obj.AddPath(newPath().moveTo(x0, y0).lineTo(x1, y1).p))
		cv.obj.Restore()
	}
	if gl := gridLine(valAxNode, "majorGridlines"); gl != nil {
		for _, v := range valTicks {
			if horizontal {
				strokeLine(gl, vpos(v), py, vpos(v), py+ph)
			} else {
				strokeLine(gl, px, vpos(v), px+pw, vpos(v))
			}
		}
	}
	if gl := gridLine(catAxNode, "majorGridlines"); gl != nil && scatter {
		for _, v := range catTicks {
			x := px + xa.pos(v)*pw
			strokeLine(gl, x, py, x, py+ph)
		}
	}
	// category geometry
	band := 0.0
	if ncat > 0 {
		if horizontal {
			band = ph / float64(ncat)
		} else {
			band = pw / float64(ncat)
		}
	}
	between := valAxNode.child("crossBetween").attrStr("val", "between") == "between" || p0.kind == "bar"
	catCenter := func(i int) float64 {
		if catReverse != horizontal {
			// bar charts list categories from the bottom up
			i = ncat - 1 - i
		}
		off := float64(i) * band
		if between {
			off += band / 2
		} else if ncat > 1 {
			off = float64(i) * (func() float64 {
				if horizontal {
					return ph
				}
				return pw
			}()) / float64(ncat-1)
		}
		if horizontal {
			return py + off
		}
		return px + off
	}
	// plots
	for _, p := range plots {
		switch p.kind {
		case "bar":
			ch.bars(p, ncat, band, catCenter, vpos, va)
		case "line", "area":
			ch.lines(p, catCenter, vpos, va)
		case "scatter":
			ch.scatterPlot(p, func(x float64) float64 { return px + xa.pos(x)*pw }, vpos)
		}
	}
	// axis lines
	axisLine := func(ax *node) *line {
		if ln := ax.path("spPr", "ln"); ln != nil {
			return s.resolveLine([]*node{ln}, []string{ch.part}, nil, s.cc)
		}
		return &line{width: 0.75, fill: fill{kind: fillSolid, color: rgba{0.53, 0.53, 0.53, 1}}, miter: 8}
	}
	zero := vpos(math.Max(va.lo, math.Min(0, va.hi)))
	if !catDeleted {
		if horizontal {
			strokeLine(axisLine(catAxNode), zero, py, zero, py+ph)
		} else if scatter {
			strokeLine(axisLine(catAxNode), px, zero, px+pw, zero)
		} else {
			strokeLine(axisLine(catAxNode), px, zero, px+pw, zero)
		}
	}
	if !va.deleted {
		if horizontal {
			strokeLine(axisLine(valAxNode), px, py+ph, px+pw, py+ph)
		} else {
			x := px
			if scatter {
				x = px + xa.pos(math.Max(xa.lo, math.Min(0, xa.hi)))*pw
			}
			strokeLine(axisLine(valAxNode), x, py, x, py+ph)
		}
	}
	// tick labels
	if showVal {
		for i, v := range valTicks {
			if horizontal {
				ch.label(valStyle, valLabels[i], vpos(v), py+ph+3, 'c', 't')
			} else {
				ch.label(valStyle, valLabels[i], px-4, vpos(v), 'r', 'm')
			}
		}
	}
	if showCat {
		if scatter {
			for i, v := range catTicks {
				ch.label(catStyle, catLabels[i], px+xa.pos(v)*pw, py+ph+3, 'c', 't')
			}
		} else {
			for i, l := range catLabels {
				if horizontal {
					ch.label(catStyle, l, px-4, catCenter(i), 'r', 'm')
				} else {
					ch.label(catStyle, l, catCenter(i), py+ph+3, 'c', 't')
				}
			}
		}
	}
}

func catLabel(s string, ax *node) string {
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		code := ax.child("numFmt").attrStr("formatCode", "General")
		return formatNumber(v, code)
	}
	return s
}

// bars draws a bar or column plot.
func (ch *chartCtx) bars(p *chartPlot, ncat int, band float64, center func(int) float64, vpos func(float64) float64, va *axis) {
	cv := ch.cv
	m := float64(len(p.series))
	stacked := p.grouping == "stacked" || p.grouping == "percentStacked"
	slots := m
	if stacked {
		slots = 1
	}
	ov := p.overlap / 100
	barW := band / (slots - (slots-1)*ov + p.gap/100)
	group := barW * (slots - (slots-1)*ov)
	totals := make([]float64, ncat)
	if p.grouping == "percentStacked" {
		for _, se := range p.series {
			for i, v := range se.vals {
				if !math.IsNaN(v) {
					totals[i] += math.Abs(v)
				}
			}
		}
	}
	pos, neg := make([]float64, ncat), make([]float64, ncat)
	base := math.Max(va.lo, math.Min(0, va.hi))
	for si, se := range p.series {
		for i, v := range se.vals {
			if math.IsNaN(v) || i >= ncat {
				continue
			}
			if p.grouping == "percentStacked" && totals[i] != 0 {
				v /= totals[i]
			}
			v0 := base
			if stacked {
				if v >= 0 {
					v0 = pos[i]
					pos[i] += v
					v = pos[i]
				} else {
					v0 = neg[i]
					neg[i] += v
					v = neg[i]
				}
			}
			slot := float64(si)
			if stacked {
				slot = 0
			}
			c0 := center(i) - group/2 + slot*barW*(1-ov)
			a, b := vpos(v0), vpos(v)
			x, y, w, h := c0, math.Min(a, b), barW, math.Abs(b-a)
			if p.horizontal {
				// categories run down the band; the slot order is reversed so
				// that the first series is at the bottom like in Office
				c0 = center(i) + group/2 - slot*barW*(1-ov) - barW
				x, y, w, h = math.Min(a, b), c0, math.Abs(b-a), barW
			}
			f := ch.seriesFill(p, se, i)
			l := ch.seriesLine(p, se, i)
			cv.obj.Save()
			cv.obj.Translate(f32(x), f32(y))
			ch.s.paintGeometry(cv, rectGeometry(w, h), f, l, nil, w, h)
			cv.obj.Restore()
			ch.dataLabel(p, se, i, v, x+w/2, y, x+w, y+h/2)
		}
	}
}

// lines draws line and area plots.
func (ch *chartCtx) lines(p *chartPlot, center func(int) float64, vpos func(float64) float64, va *axis) {
	cv := ch.cv
	stacked := p.grouping == "stacked" || p.grouping == "percentStacked"
	n := 0
	for _, se := range p.series {
		n = max(n, len(se.vals))
	}
	totals := make([]float64, n)
	for _, se := range p.series {
		for i, v := range se.vals {
			if !math.IsNaN(v) {
				totals[i] += math.Abs(v)
			}
		}
	}
	acc := make([]float64, n)
	base := vpos(math.Max(va.lo, math.Min(0, va.hi)))
	type pt struct{ x, y, v float64 }
	var prevTop []pt
	for _, se := range p.series {
		var pts []pt
		for i, v := range se.vals {
			if math.IsNaN(v) {
				continue
			}
			if p.grouping == "percentStacked" && totals[i] != 0 {
				v /= totals[i]
			}
			if stacked {
				acc[i] += v
				v = acc[i]
			}
			pts = append(pts, pt{center(i), vpos(v), v})
		}
		if len(pts) == 0 {
			continue
		}
		if p.kind == "area" {
			path := newPath()
			path.moveTo(pts[0].x, pts[0].y)
			for _, q := range pts[1:] {
				path.lineTo(q.x, q.y)
			}
			if stacked && prevTop != nil {
				for i := len(prevTop) - 1; i >= 0; i-- {
					path.lineTo(prevTop[i].x, prevTop[i].y)
				}
			} else {
				path.lineTo(pts[len(pts)-1].x, base).lineTo(pts[0].x, base)
			}
			path.close()
			f := ch.seriesFill(p, se, -1)
			cv.obj.Save()
			if cv.setFill(f, 1, 1) {
				cv.obj.FillPath(cv.obj.AddPath(path.p), 0)
			}
			if l := ch.seriesLine(p, se, -1); l != nil {
				cv.setLine(l, 1, 1)
				cv.obj.StrokePath(cv.obj.AddPath(path.p))
			}
			cv.obj.Restore()
			prevTop = pts
			continue
		}
		l := ch.seriesLine(p, se, -1)
		if l != nil && len(pts) > 1 {
			path := newPath().moveTo(pts[0].x, pts[0].y)
			for _, q := range pts[1:] {
				path.lineTo(q.x, q.y)
			}
			cv.obj.Save()
			cv.setLine(l, 1, 1)
			cv.obj.StrokePath(cv.obj.AddPath(path.p))
			cv.obj.Restore()
		}
		if p.markers {
			mf := ch.markerFill(p, se)
			for _, q := range pts {
				ch.marker(se, mf, q.x, q.y, "", 0)
			}
		}
		for i, q := range pts {
			ch.dataLabel(p, se, i, q.v, q.x, q.y-4, q.x+4, q.y)
		}
	}
}

func (ch *chartCtx) scatterPlot(p *chartPlot, xpos, ypos func(float64) float64) {
	cv := ch.cv
	style := p.node.child("scatterStyle").attrStr("val", "lineMarker")
	for _, se := range p.series {
		var xs, ys []float64
		for i, v := range se.vals {
			x := float64(i + 1)
			if i < len(se.xs) {
				x = se.xs[i]
			}
			if math.IsNaN(v) || math.IsNaN(x) {
				continue
			}
			xs, ys = append(xs, xpos(x)), append(ys, ypos(v))
		}
		if len(xs) == 0 {
			continue
		}
		l := ch.seriesLine(p, se, -1)
		if style == "marker" || style == "none" && se.spPr.child("ln") == nil {
			l = nil
		}
		if l != nil && len(xs) > 1 {
			path := newPath().moveTo(xs[0], ys[0])
			for i := 1; i < len(xs); i++ {
				path.lineTo(xs[i], ys[i])
			}
			cv.obj.Save()
			cv.setLine(l, 1, 1)
			cv.obj.StrokePath(cv.obj.AddPath(path.p))
			cv.obj.Restore()
		}
		if se.marker.child("symbol").attrStr("val", "auto") != "none" {
			mf := ch.markerFill(p, se)
			for i := range xs {
				ch.marker(se, mf, xs[i], ys[i], "", 0)
			}
		}
	}
}

// pie draws pie and doughnut plots.
func (ch *chartCtx) pie(p *chartPlot, x, y, w, h float64) {
	if len(p.series) == 0 {
		return
	}
	cv := ch.cv
	se := p.series[0]
	total := 0.0
	for _, v := range se.vals {
		if !math.IsNaN(v) && v > 0 {
			total += v
		}
	}
	if total == 0 {
		return
	}
	r := math.Min(w, h) / 2 * 0.95
	cx, cy := x+w/2, y+h/2
	ang := (p.firstAng - 90) * math.Pi / 180
	pv := &chartPlot{vary: true, series: []*chartSeries{se}, dLbls: p.dLbls}
	for i, v := range se.vals {
		if math.IsNaN(v) || v <= 0 {
			continue
		}
		sweep := v / total * 2 * math.Pi
		path := &bdf.Path{}
		if p.hole > 0 {
			path.Ellipse(f32(cx), f32(cy), f32(r), f32(r), 0, f32(ang), f32(ang+sweep), false)
			path.Ellipse(f32(cx), f32(cy), f32(r*p.hole), f32(r*p.hole), 0, f32(ang+sweep), f32(ang), true)
			path.Close()
		} else {
			path.MoveTo(f32(cx), f32(cy))
			path.Ellipse(f32(cx), f32(cy), f32(r), f32(r), 0, f32(ang), f32(ang+sweep), false)
			path.Close()
		}
		f := ch.seriesFill(pv, se, i)
		l := ch.seriesLine(pv, se, i)
		if l == nil && fillElem(se.spPr) == nil {
			// Office separates slices with a thin line of the background color
			bg, _ := ch.s.cc.scheme1("bg1")
			l = &line{width: 0.75, fill: fill{kind: fillSolid, color: bg}, join: 1, miter: 8}
		}
		ref := cv.obj.AddPath(path)
		cv.obj.Save()
		if cv.setFill(f, 2*r, 2*r) {
			cv.obj.FillPath(ref, 0)
		}
		if l != nil {
			cv.setLine(l, 2*r, 2*r)
			cv.obj.StrokePath(ref)
		}
		cv.obj.Restore()
		mid := ang + sweep/2
		lr := r * 0.65
		if p.hole > 0 {
			lr = r * (1 + p.hole) / 2
		}
		lx, ly := cx+lr*math.Cos(mid), cy+lr*math.Sin(mid)
		ch.pieLabel(pv, se, i, v, v/total, lx, ly)
		ang += sweep
	}
}

// dLblsFor returns the data label settings that apply to a series.
func dLblsFor(p *chartPlot, se *chartSeries) *node {
	if se.dLbls != nil {
		return se.dLbls
	}
	return p.dLbls
}

func showFlag(d *node, name string) bool {
	return d.child(name).attrStr("val", "0") == "1"
}

// dataLabel draws a value label for bars and line points (above the point,
// or right of horizontal bars).
func (ch *chartCtx) dataLabel(p *chartPlot, se *chartSeries, i int, v, x, y, hx, hy float64) {
	d := dLblsFor(p, se)
	if d == nil || d.child("delete").attrStr("val", "0") == "1" {
		return
	}
	var parts []string
	if showFlag(d, "showSerName") {
		parts = append(parts, se.name)
	}
	if showFlag(d, "showCatName") && i < len(se.cats) {
		parts = append(parts, se.cats[i])
	}
	if showFlag(d, "showVal") {
		code := d.child("numFmt").attrStr("formatCode", se.format)
		parts = append(parts, formatNumber(v, code))
	}
	if len(parts) == 0 {
		return
	}
	st := ch.textStyle(d.child("txPr"), ch.baseSize(), false)
	text := strings.Join(parts, ", ")
	if p.kind == "bar" && p.horizontal {
		ch.label(st, text, hx+3, hy, 'l', 'm')
		return
	}
	ch.label(st, text, x, y-2, 'c', 'b')
}

func (ch *chartCtx) pieLabel(p *chartPlot, se *chartSeries, i int, v, frac, x, y float64) {
	d := dLblsFor(p, se)
	if d == nil || d.child("delete").attrStr("val", "0") == "1" {
		return
	}
	var parts []string
	if showFlag(d, "showCatName") && i < len(se.cats) {
		parts = append(parts, se.cats[i])
	}
	if showFlag(d, "showVal") {
		parts = append(parts, formatNumber(v, d.child("numFmt").attrStr("formatCode", se.format)))
	}
	if showFlag(d, "showPercent") {
		parts = append(parts, formatNumber(frac, "0%"))
	}
	if len(parts) == 0 {
		return
	}
	st := ch.textStyle(d.child("txPr"), ch.baseSize(), false)
	ch.label(st, strings.Join(parts, "\u00a0"), x, y, 'c', 'm')
}
