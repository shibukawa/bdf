package kicad

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// Schematic defaults (KiCad's default_values.h), mm.
const (
	schLineWidth      = 0.1524 // 6 mils
	schBusWidth       = 0.3048 // 12 mils
	schNoConnectSize  = 1.2192 // 48 mils
	schJunctionDiam   = 0.9144 // 36 mils
	schPinNameOffset  = 0.508  // 20 mils
	schTargetPinR     = 0.381  // 15 mils
	schPinTextMargin  = 0.1016 // 4 mils
	schTextFudgeY     = -0.25  // SCH_TEXT's offset to match KiCad 6
	schDNPWidth       = 3 * schLineWidth
	schDirectiveLen   = 2.54  // 100 mils
	schDirectiveSym   = 0.508 // 20 mils
	schDanglingSymbol = 0.3048
)

var schDNPMarker = bdf.RGBA(220, 9, 13, 204)

// ptPerMM converts millimeters to points.
const ptPerMM = 72 / 25.4

// schematic converts the schematic of an input into a view whose pages are
// its sheets.
func (c *conv) schematic(in *input) ([]pageOut, error) {
	c.iu = 10000 // 100 nm
	c.defaultPen = c.sch.line
	sheets, err := c.loadSchematic(in.sch, in.schData)
	if err != nil {
		return nil, fmt.Errorf("kicad: %w", err)
	}
	root := sheets[0]
	for _, s := range sheets {
		if s.parent == nil {
			root = s
		}
	}
	tb := root.file.n.child("title_block")
	title := tb.str("title")
	if title == "" {
		title = c.project
	}
	if c.opts.Title == "" && len(c.doc.Meta.DC.Title) == 0 && tb.str("title") != "" {
		c.doc.Meta.DC.Title = bdf.DCValues{tb.str("title")}
	}
	if d := tb.str("date"); d != "" && len(c.doc.Meta.DC.Date) == 0 {
		c.doc.Meta.DC.Date = bdf.DCValues{d}
	}
	if co := tb.str("company"); co != "" && len(c.doc.Meta.DC.Publisher) == 0 {
		c.doc.Meta.DC.Publisher = bdf.DCValues{co}
	}
	if title == "" {
		title = "Schematic"
	}
	c.schSheet = c.readWorksheet(in.sch, c.schSheetFile)
	view := c.doc.NewView("schematic", bdf.ViewFixed, title)
	sel := c.opts.Pages.Numbers(len(sheets))
	if sel == nil {
		for i := range sheets {
			sel = append(sel, i+1)
		}
	}
	// the page of the view each sheet lands on, for links to sheets
	pageOf := map[*sheetInst]int{}
	for i, n := range sel {
		if n >= 1 && n <= len(sheets) {
			pageOf[sheets[n-1]] = i + 1
		}
	}
	var out []pageOut
	for _, n := range sel {
		if n < 1 || n > len(sheets) {
			return nil, fmt.Errorf("kicad: sheet %d out of range (1-%d)", n, len(sheets))
		}
		out = append(out, c.schPageSafe(view, sheets[n-1], sheets, pageOf))
	}
	return out, nil
}

// paperSize returns the size of a sheet's paper in mm (landscape unless the
// file says portrait), as KiCad defines the formats (in mils).
func paperSize(n *node) (w, h float64) {
	name := n.arg(0)
	mils := map[string][2]float64{
		"A5": {8268, 5827}, "A4": {11693, 8268}, "A3": {16535, 11693}, "A2": {23386, 16535}, "A1": {33110, 23386},
		"A0": {46811, 33110}, "A": {11000, 8500}, "B": {17000, 11000}, "C": {22000, 17000}, "D": {34000, 22000},
		"E": {44000, 34000}, "USLetter": {11000, 8500}, "USLegal": {14000, 8500}, "USLedger": {17000, 11000},
		"GERBER": {32000, 32000},
	}
	if name == "User" {
		w, h = n.num(1), n.num(2)
		if w > 0 && h > 0 {
			return w, h
		}
		return 297, 210
	}
	s, ok := mils[name]
	if !ok {
		s = mils["A4"]
	}
	w, h = s[0]*0.0254, s[1]*0.0254
	if n.hasAtom("portrait") {
		w, h = h, w
	}
	return w, h
}

// schPage is a page being drawn.
type schPage struct {
	c      *conv
	s      *sheetInst
	all    []*sheetInst
	pageOf map[*sheetInst]int
	w, h   float64 // paper, mm
	d      *cad.Drawing
	v6refs map[string]*node // KiCad 6: the symbol instances of the root
	syms   []*symInst
	count  int // sheets in the schematic
	links  []pageLink
}

// pageLink is a link from a sheet's box to its page.
type pageLink struct {
	x, y, w, h float64 // mm
	page       int
}

func (c *conv) schPageSafe(view *bdf.View, s *sheetInst, all []*sheetInst, pageOf map[*sheetInst]int) (out pageOut) {
	w, h := paperSize(s.file.n.child("paper"))
	k := ptPerMM
	page := view.AddPage(float32(w*k), float32(h*k))
	out = pageOut{view: view, page: page}
	if c.drawn > maxDrawn {
		c.warnOnce("sch-drawn", "the schematic is too large to draw in full: sheet %s and the pages after it are left blank", s.namePath)
		out.cv = c.cvs.New()
		out.cv.Obj.SetBBox(0, 0, page.W, page.H)
		return out
	}
	defer func() {
		if r := recover(); r != nil {
			c.warnf("sheet %s: internal error: %v", s.namePath, r)
			out.cv = c.cvs.New()
			out.cv.Obj.SetBBox(0, 0, page.W, page.H)
			out.background = nil
		}
	}()
	p := &schPage{c: c, s: s, all: all, pageOf: pageOf, w: w, h: h, d: &cad.Drawing{}, count: len(all)}
	root := s
	for root.parent != nil {
		root = root.parent
	}
	p.v6refs = map[string]*node{}
	for _, path := range root.file.n.child("symbol_instances").children("path") {
		p.v6refs[path.arg(0)] = path
	}
	m := canvas.Scale(k, k)
	view2 := cad.Rect{}.Add(cad.Point{}).Add(cad.Point{X: w * k, Y: h * k})
	// the paper and the frame of the drawing sheet: the same for every
	// page of the size, so stored once
	bg := &cad.Drawing{}
	bg.Fill((&cad.Path{}).Polyline([]cad.Point{{X: 0, Y: 0}, {X: w, Y: 0}, {X: w, Y: h}, {X: 0, Y: h}}, true), cad.Fill{Color: schBackground}, false)
	c.drawWorksheet(c.schSheet, &wsPage{w: w, h: h, first: s.page == "1", color: schWorksheet, minPen: c.sch.line,
		vars: p.titleVars(), bg: bg, fg: p.d})
	out.background = c.cvs.New()
	out.background.Obj.SetBBox(0, 0, page.W, page.H)
	c.plotter.Plot(out.background, bg, m, view2)

	p.draw()
	c.drawn += s.file.size + c.schSheet.size + 100*(len(bg.Items)+len(p.d.Items))
	out.cv = c.cvs.New()
	out.cv.Obj.SetBBox(0, 0, page.W, page.H)
	c.plotter.Plot(out.cv, p.d, m, view2)
	for _, l := range p.links {
		out.cv.Obj.Link(float32(l.x*k), float32(l.y*k), float32(l.w*k), float32(l.h*k), "#page="+strconv.Itoa(l.page))
	}
	return out
}

// titleVars returns the text variables of the page: the title block, the
// page number and the sheet.
func (p *schPage) titleVars() func(string) (string, bool) {
	tb := p.s.file.n.child("title_block")
	comments := map[string]string{}
	for _, cm := range tb.children("comment") {
		comments["COMMENT"+cm.arg(0)] = cm.arg(1)
	}
	return func(name string) (string, bool) {
		switch name {
		case "TITLE":
			return tb.str("title"), true
		case "ISSUE_DATE":
			return tb.str("date"), true
		case "REVISION":
			return tb.str("rev"), true
		case "COMPANY":
			return tb.str("company"), true
		case "#":
			return p.s.page, true
		case "##":
			return strconv.Itoa(p.count), true
		case "SHEETNAME":
			if p.s.parent == nil {
				return "Root", true
			}
			return p.s.name, true
		case "SHEETPATH":
			return p.s.namePath, true
		case "FILENAME":
			return p.s.fileBase(), true
		case "FILEPATH":
			return p.s.file.path, true
		case "PAPER":
			return p.s.file.n.child("paper").arg(0), true
		case "KICAD_VERSION":
			return kicadVersion(p.s.file.n), true
		}
		if v, ok := comments[name]; ok {
			return v, true
		}
		if strings.HasPrefix(name, "COMMENT") {
			if _, err := strconv.Atoi(name[7:]); err == nil {
				return "", true
			}
		}
		return "", false
	}
}

// expand resolves the text variables of a text on the page.
func (p *schPage) expand(s string) string { return p.c.expand(s, p.titleVars()) }

// draw draws the items of the sheet, in layers as KiCad stacks them.
func (p *schPage) draw() {
	n := p.s.file.n
	for _, sym := range n.children("symbol") {
		if si := p.symbolInstance(sym); si != nil {
			p.syms = append(p.syms, si)
		}
	}
	for _, im := range n.children("image") {
		p.image(im)
	}
	for _, sh := range n.children("sheet") {
		p.sheetBox(sh, true)
	}
	for _, si := range p.syms {
		p.symbolBody(si, true)
	}
	for _, it := range n.lists() {
		switch it.name {
		case "polyline", "rectangle", "circle", "arc", "bezier":
			p.shape(it, nil, false)
		case "rule_area":
			if pl := it.child("polyline"); pl != nil {
				p.shape(pl, nil, true)
			}
		case "text_box":
			p.textBox(it, nil)
		case "table":
			p.table(it)
		}
	}
	for _, si := range p.syms {
		p.symbolBody(si, false)
	}
	for _, sh := range n.children("sheet") {
		p.sheetBox(sh, false)
	}
	for _, it := range n.lists() {
		switch it.name {
		case "wire", "bus":
			p.wire(it)
		case "bus_entry":
			p.busEntry(it)
		}
	}
	for _, it := range n.lists() {
		switch it.name {
		case "junction":
			p.junction(it)
		case "no_connect":
			p.noConnect(it)
		}
	}
	for _, it := range n.lists() {
		switch it.name {
		case "label", "global_label", "hierarchical_label":
			p.label(it)
		case "netclass_flag", "directive_label":
			p.directive(it)
		case "text":
			p.note(it)
		}
	}
	for _, si := range p.syms {
		p.symbolFields(si)
	}
	for _, sh := range n.children("sheet") {
		p.sheetFields(sh)
	}
	for _, si := range p.syms {
		if si.dnp {
			p.dnpMark(si)
		}
	}
}

// stroke is the (stroke ...) of an item.
type stroke struct {
	width    float64 // mm; 0: the default
	dash     string  // solid, dash, dot, dash_dot, dash_dot_dot, default
	color    bdf.Color
	hasColor bool
}

func parseStroke(n *node) stroke {
	s := stroke{dash: "default"}
	if n == nil {
		return s
	}
	s.width = n.numOf("width", 0)
	if t := n.str("type"); t != "" {
		s.dash = t
	}
	s.color, s.hasColor = parseColor(n.child("color"))
	return s
}

// pen returns the pen of a stroke of width w (mm) and style.
func pen(w float64, dash string, color bdf.Color) cad.Pen {
	return dashedPen(w, dash, color, defaultSchSettings.dash, defaultSchSettings.gap)
}

// pen is pen with the dashes of the schematic's project.
func (p *schPage) pen(w float64, dash string, color bdf.Color) cad.Pen {
	return dashedPen(w, dash, color, p.c.sch.dash, p.c.sch.gap)
}

// dashedPen returns a pen whose dashes and gaps are the ratios of its
// width, less and more the width its round caps add
// (STROKE_PARAMS::GetDashLength and GetGapLength).
func dashedPen(w float64, dash string, color bdf.Color, dashRatio, gapRatio float64) cad.Pen {
	p := cad.Pen{Color: color, WorldWidth: w, Cap: cad.CapRound, Join: cad.JoinRound}
	if w <= 0 {
		p.WorldWidth = 0
	}
	dashLen := math.Max(dashRatio-1, 1) * w
	gap := math.Max(gapRatio+1, 1) * w
	dot := 0.2 * w
	switch dash {
	case "dash":
		p.Dash = []float64{dashLen, gap}
	case "dot":
		p.Dash = []float64{dot, gap}
	case "dash_dot":
		p.Dash = []float64{dashLen, gap, dot, gap}
	case "dash_dot_dot":
		p.Dash = []float64{dashLen, gap, dot, gap, dot, gap}
	}
	return p
}

// dimIf dims a color when the symbol it belongs to is not placed (DNP).
func dimIf(c bdf.Color, dnp bool) bdf.Color {
	if dnp {
		return dim(c, schBackground)
	}
	return c
}

// arcThrough returns the circle through three points: its center and
// radius, and the angles of the first and last point (radians, y down)
// such that going from a0 to a1 passes the middle one.
func arcThrough(s, m, e cad.Point) (c cad.Point, r, a0, a1 float64, ok bool) {
	ax, ay := s.X, s.Y
	bx, by := m.X, m.Y
	cx, cy := e.X, e.Y
	d := 2 * (ax*(by-cy) + bx*(cy-ay) + cx*(ay-by))
	if math.Abs(d) < 1e-12 {
		return c, 0, 0, 0, false
	}
	ux := ((ax*ax+ay*ay)*(by-cy) + (bx*bx+by*by)*(cy-ay) + (cx*cx+cy*cy)*(ay-by)) / d
	uy := ((ax*ax+ay*ay)*(cx-bx) + (bx*bx+by*by)*(ax-cx) + (cx*cx+cy*cy)*(bx-ax)) / d
	c = cad.Point{X: ux, Y: uy}
	r = math.Hypot(ax-ux, ay-uy)
	a0 = math.Atan2(ay-uy, ax-ux)
	am := math.Atan2(by-uy, bx-ux)
	a1 = math.Atan2(cy-uy, cx-ux)
	norm := func(a float64) float64 {
		for a < 0 {
			a += 2 * math.Pi
		}
		for a >= 2*math.Pi {
			a -= 2 * math.Pi
		}
		return a
	}
	sweep := norm(a1 - a0)
	if norm(am-a0) > sweep {
		sweep -= 2 * math.Pi
	}
	return c, r, a0, a0 + sweep, true
}

// shapePath builds the path of a graphic item (polyline, rectangle,
// circle, arc, bezier), its points mapped by tr (nil: as they are).
func shapePath(it *node, tr func(cad.Point) cad.Point) (*cad.Path, bool) {
	if tr == nil {
		tr = func(p cad.Point) cad.Point { return p }
	}
	xy := func(n *node) cad.Point { return tr(cad.Point{X: n.num(0), Y: n.num(1)}) }
	p := &cad.Path{}
	switch it.name {
	case "polyline":
		var pts []cad.Point
		for _, q := range it.child("pts").children("xy") {
			pts = append(pts, xy(q))
		}
		if len(pts) < 2 {
			return nil, false
		}
		p.Polyline(pts, false)
		return p, pts[0] == pts[len(pts)-1]
	case "rectangle":
		a, b := xy(it.child("start")), xy(it.child("end"))
		if r := it.numOf("radius", 0); r > 0 {
			return roundRect(tr, it.child("start"), it.child("end"), r), true
		}
		// the corners go through the transform one by one
		s, e := it.child("start"), it.child("end")
		c1 := tr(cad.Point{X: e.num(0), Y: s.num(1)})
		c2 := tr(cad.Point{X: s.num(0), Y: e.num(1)})
		p.Polyline([]cad.Point{a, c1, b, c2}, true)
		return p, true
	case "circle":
		c := xy(it.child("center"))
		r := it.numOf("radius", 0)
		if cc := it.child("center"); cc != nil && it.child("end") != nil {
			// KiCad 6 symbols: a circle by its center and a point on it
			e := it.child("end")
			r = math.Hypot(e.num(0)-cc.num(0), e.num(1)-cc.num(1))
		}
		if r <= 0 {
			return nil, false
		}
		p.Circle(c, r)
		return p, true
	case "arc":
		s, m, e := it.child("start"), it.child("mid"), it.child("end")
		if s == nil || e == nil {
			return nil, false
		}
		if m == nil {
			return nil, false
		}
		c, r, a0, a1, ok := arcThrough(xy(s), xy(m), xy(e))
		if !ok {
			p.Polyline([]cad.Point{xy(s), xy(e)}, false)
			return p, false
		}
		p.Arc(c, r, a0, a1)
		return p, false
	case "bezier":
		var pts []cad.Point
		for _, q := range it.child("pts").children("xy") {
			pts = append(pts, xy(q))
		}
		if len(pts) != 4 {
			if len(pts) >= 2 {
				p.Polyline(pts, false)
				return p, false
			}
			return nil, false
		}
		p.MoveTo(pts[0].X, pts[0].Y)
		p.CubicTo(pts[1], pts[2], pts[3])
		return p, false
	}
	return nil, false
}

// roundRect builds a rectangle with rounded corners.
func roundRect(tr func(cad.Point) cad.Point, s, e *node, r float64) *cad.Path {
	x0, y0 := math.Min(s.num(0), e.num(0)), math.Min(s.num(1), e.num(1))
	x1, y1 := math.Max(s.num(0), e.num(0)), math.Max(s.num(1), e.num(1))
	r = math.Min(r, math.Min(x1-x0, y1-y0)/2)
	const k = 0.5523
	pts := func(pp ...cad.Point) []cad.Point {
		for i := range pp {
			pp[i] = tr(pp[i])
		}
		return pp
	}
	p := &cad.Path{}
	q := pts(cad.Point{X: x0 + r, Y: y0})[0]
	p.MoveTo(q.X, q.Y)
	corner := func(line, c1, c2, end cad.Point) {
		l := pts(line)[0]
		p.LineTo(l.X, l.Y)
		cc := pts(c1, c2, end)
		p.CubicTo(cc[0], cc[1], cc[2])
	}
	corner(cad.Point{X: x1 - r, Y: y0}, cad.Point{X: x1 - r + r*k, Y: y0}, cad.Point{X: x1, Y: y0 + r - r*k}, cad.Point{X: x1, Y: y0 + r})
	corner(cad.Point{X: x1, Y: y1 - r}, cad.Point{X: x1, Y: y1 - r + r*k}, cad.Point{X: x1 - r + r*k, Y: y1}, cad.Point{X: x1 - r, Y: y1})
	corner(cad.Point{X: x0 + r, Y: y1}, cad.Point{X: x0 + r - r*k, Y: y1}, cad.Point{X: x0, Y: y1 - r + r*k}, cad.Point{X: x0, Y: y1 - r})
	corner(cad.Point{X: x0, Y: y0 + r}, cad.Point{X: x0, Y: y0 + r - r*k}, cad.Point{X: x0 + r - r*k, Y: y0}, cad.Point{X: x0 + r, Y: y0})
	p.Close()
	return p
}

// fill is the (fill ...) of an item.
type fill struct {
	kind  string // none, outline, background, color, hatch, reverse_hatch, cross_hatch
	color bdf.Color
	has   bool
}

func parseFill(n *node) fill {
	f := fill{kind: "none"}
	if n == nil {
		return f
	}
	if t := n.str("type"); t != "" {
		f.kind = t
	}
	f.color, f.has = parseColor(n.child("color"))
	return f
}

// shape draws a graphic item of the sheet (not of a symbol): its fill and
// its outline. ruleArea draws a rule area's outline.
func (p *schPage) shape(it *node, tr func(cad.Point) cad.Point, ruleArea bool) {
	path, closed := shapePath(it, tr)
	if path == nil {
		return
	}
	if ruleArea && !closed {
		// a rule area is a closed polygon whether its points close or not
		path.Close()
		closed = true
	}
	st := parseStroke(it.child("stroke"))
	f := parseFill(it.child("fill"))
	if closed || it.name == "circle" || it.name == "rectangle" || it.name == "arc" {
		p.fillShape(path, f, schNote, false, st.width)
	}
	color := schNote
	if ruleArea {
		color = schRuleArea
	}
	if st.hasColor {
		color = st.color
	}
	if st.width < 0 || st.dash == "none" {
		return
	}
	w := st.width
	if w <= 0 {
		w = p.c.sch.line
	}
	p.d.Stroke(path, p.pen(w, st.dash, color))
}

// fillShape fills a shape by its fill type; outline is the color of
// "outline" fills and of hatching, strokeW the width of the shape's
// outline (hatch lines are half as wide).
func (p *schPage) fillShape(path *cad.Path, f fill, outline bdf.Color, dnp bool, strokeW float64) {
	var c bdf.Color
	switch f.kind {
	case "outline":
		c = outline
	case "background":
		c = schBody
	case "color":
		if !f.has {
			return
		}
		c = f.color
	case "hatch", "reverse_hatch", "cross_hatch":
		c = outline
		if f.has {
			c = f.color
		}
		p.hatch(path, f.kind, dimIf(c, dnp), strokeW)
		return
	default:
		return
	}
	p.d.Fill(path, cad.Fill{Color: dimIf(c, dnp)}, false)
}

// hatch draws a hatched fill as KiCad 10 draws it: lines y = ±x + a, a
// stepping by the spacing from 0, half as wide as the outline and 40 of
// their widths apart (at most 100 lines across the shape).
func (p *schPage) hatch(path *cad.Path, kind string, color bdf.Color, strokeW float64) {
	if strokeW <= 0 {
		strokeW = p.c.sch.line
	}
	w := strokeW / 2
	spacing := w * 40
	b := path.Bounds()
	if major := math.Max(b.W(), b.H()); major/spacing > 100 {
		spacing = major / 100
	}
	var slopes []float64
	switch kind {
	case "hatch":
		slopes = []float64{-1}
	case "reverse_hatch":
		slopes = []float64{1}
	default:
		slopes = []float64{1, -1}
	}
	var lines []cad.PatternLine
	for _, sl := range slopes {
		lines = append(lines, cad.PatternLine{Angle: math.Atan(sl), Offset: cad.Point{Y: spacing}})
	}
	if !p.d.Hatch(path, lines, cad.Pen{Color: color, WorldWidth: w, Cap: cad.CapRound}) {
		p.d.Fill(path, cad.Fill{Color: color}, false)
	}
}

// wire draws a wire or a bus.
func (p *schPage) wire(it *node) {
	var pts []cad.Point
	for _, q := range it.child("pts").children("xy") {
		pts = append(pts, cad.Point{X: q.num(0), Y: q.num(1)})
	}
	if len(pts) < 2 {
		return
	}
	st := parseStroke(it.child("stroke"))
	w, color := p.c.sch.wire, schWire
	if it.name == "bus" {
		w, color = p.c.sch.bus, schBus
	}
	if st.width > 0 {
		w = st.width
	}
	if st.hasColor {
		color = st.color
	}
	p.d.Stroke((&cad.Path{}).Polyline(pts, false), p.pen(w, st.dash, color))
}

// busEntry draws the diagonal line from a bus to a wire.
func (p *schPage) busEntry(it *node) {
	x, y := it.child("at").xy()
	sz := it.child("size")
	st := parseStroke(it.child("stroke"))
	w, color := p.c.sch.wire, schWire
	if st.width > 0 {
		w = st.width
	}
	if st.hasColor {
		color = st.color
	}
	p.d.Stroke((&cad.Path{}).Polyline([]cad.Point{{X: x, Y: y}, {X: x + sz.num(0), Y: y + sz.num(1)}}, false), p.pen(w, st.dash, color))
}

// junction draws the dot where wires meet.
func (p *schPage) junction(it *node) {
	x, y := it.child("at").xy()
	d := it.numOf("diameter", 0)
	if d <= 0 {
		d = p.c.sch.junction
	}
	color := schJunction
	if c, ok := parseColor(it.child("color")); ok {
		color = c
	}
	p.d.Fill((&cad.Path{}).Circle(cad.Point{X: x, Y: y}, d/2), cad.Fill{Color: color}, false)
}

// noConnect draws the cross of a pin left unconnected on purpose.
func (p *schPage) noConnect(it *node) {
	x, y := it.child("at").xy()
	s := schNoConnectSize / 2
	pn := p.pen(p.c.sch.line, "solid", schNoConnect)
	p.d.Stroke((&cad.Path{}).Polyline([]cad.Point{{X: x - s, Y: y - s}, {X: x + s, Y: y + s}}, false), pn)
	p.d.Stroke((&cad.Path{}).Polyline([]cad.Point{{X: x + s, Y: y - s}, {X: x - s, Y: y + s}}, false), pn)
}

// image draws an image embedded in the sheet: PNG data in base64, at 300
// pixels per inch unless the image says otherwise, times its scale,
// centered on its position.
func (p *schPage) image(it *node) {
	var sb strings.Builder
	for _, d := range it.children("data") {
		for i := 0; ; i++ {
			s := d.arg(i)
			if s == "" {
				break
			}
			sb.WriteString(s)
		}
	}
	data, err := base64.StdEncoding.DecodeString(sb.String())
	if err != nil || len(data) == 0 {
		p.c.warnOnce("image-data", "an image whose data cannot be read is left out")
		return
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || cfg.Width <= 0 || cfg.Height <= 0 {
		p.c.warnOnce("image-format", "an image in a format that cannot be read is left out")
		return
	}
	scale := it.numOf("scale", 1)
	if scale <= 0 {
		scale = 1
	}
	ppi := pngPPI(data)
	if ppi <= 0 {
		ppi = 300
	}
	mmPerPx := 25.4 / ppi * scale
	w, h := float64(cfg.Width)*mmPerPx, float64(cfg.Height)*mmPerPx
	x, y := it.child("at").xy()
	hash := p.c.doc.AddImage(data)
	m := canvas.Translate(x-w/2, y-h/2).Mul(canvas.Scale(mmPerPx, mmPerPx))
	p.d.Image(&cad.Image{Image: hash, W: cfg.Width, H: cfg.Height, M: m, Smooth: true})
}

// pngPPI returns the resolution a PNG's pHYs chunk gives, in pixels per
// inch (0 when it gives none).
func pngPPI(b []byte) float64 {
	if !bytes.HasPrefix(b, []byte("\x89PNG\r\n\x1a\n")) {
		return 0
	}
	i := 8
	for i+12 <= len(b) {
		n := int(uint32(b[i])<<24 | uint32(b[i+1])<<16 | uint32(b[i+2])<<8 | uint32(b[i+3]))
		typ := string(b[i+4 : i+8])
		if n < 0 || i+12+n > len(b) {
			return 0
		}
		if typ == "pHYs" && n == 9 {
			d := b[i+8:]
			ppu := float64(uint32(d[0])<<24 | uint32(d[1])<<16 | uint32(d[2])<<8 | uint32(d[3]))
			if d[8] == 1 && ppu > 0 { // per meter
				return ppu * 0.0254
			}
			return 0
		}
		if typ == "IDAT" {
			return 0
		}
		i += 12 + n
	}
	return 0
}

// note draws a text of the sheet (SCH_TEXT).
func (p *schPage) note(it *node) {
	// KiCad no longer hides schematic text
	e := parseEffects(it.child("effects"))
	at := it.child("at")
	x, y := at.xy()
	color := schNote
	if e.hasColor {
		color = e.color
	}
	s := p.expand(it.arg(0))
	// KiCad 7 and later shift notes to match KiCad 6, and TrueType ones by
	// part of their line's extra height to match fields
	angle := keepUpright(at.num(2))
	dx, dy := 0.0, schTextFudgeY
	if e.face != "" {
		adjust := p.c.round((p.c.outlineHeight(e) - e.sizeY) * 0.4)
		ax, ay := rotatePt(0, -adjust, angle)
		dx, dy = dx+ax, dy+ay
	}
	p.c.text(p.d, placedText{s: s, x: x + dx, y: y + dy, angle: angle, e: e, color: color, brk: cad.BreakBox})
}

// spin is how a label faces: which side of its anchor the text is on.
type spin int

const (
	spinLeft spin = iota
	spinUp
	spinRight
	spinBottom
)

// keepUpright folds an angle of schematic text to 0 or 90 degrees, as
// KiCad does when it reads one (EDA_ANGLE::KeepUpright).
func keepUpright(angle float64) float64 {
	a := math.Mod(math.Round(angle), 360)
	if a < 0 {
		a += 360
	}
	switch {
	case a > 90 && a <= 270:
		a -= 180
	case a > 270:
		a -= 360
	}
	return a
}

// labelSpin returns the spin of a label as KiCad reads it: its angle kept
// upright, and the side of the anchor its justification puts the text on.
func labelSpin(angle float64, h hAlign) spin {
	vertical := math.Abs(keepUpright(angle)) == 90
	switch {
	case vertical && h == alignRight:
		return spinBottom
	case vertical:
		return spinUp
	case h == alignRight:
		return spinLeft
	}
	return spinRight
}

// spinText returns the angle and justification of a label's text.
func spinText(s spin) (angle float64, h hAlign) {
	switch s {
	case spinUp:
		return 90, alignLeft
	case spinLeft:
		return 0, alignRight
	case spinBottom:
		return 90, alignRight
	}
	return 0, alignLeft
}

// rotateSpin turns a label's shape point for its spin (KiCad's
// RotatePoint: 90° is counterclockwise on the screen).
func rotateSpin(pt cad.Point, s spin) cad.Point {
	switch s {
	case spinUp: // -90
		return cad.Point{X: -pt.Y, Y: pt.X}
	case spinRight: // 180
		return cad.Point{X: -pt.X, Y: -pt.Y}
	case spinBottom: // 90
		return cad.Point{X: pt.Y, Y: -pt.X}
	}
	return pt
}

// label draws a local, global or hierarchical label.
func (p *schPage) label(it *node) {
	e := parseEffects(it.child("effects"))
	at := it.child("at")
	x, y := at.xy()
	sp := labelSpin(at.num(2), e.h)
	angle, h := spinText(sp)
	e.h = h
	e.v = alignBottom
	e.mirror = false
	penW := penWidth(e, p.c.sch.line)
	var color bdf.Color
	var off cad.Point
	textH := e.sizeY
	s := p.expand(it.arg(0))
	switch it.name {
	case "label":
		color = schLabelLocal
		// the text's own pen, not the one it is drawn with
		dist := p.c.round(p.c.sch.textOffset*textH) + penWidth(e, 0)
		if sp == spinUp || sp == spinBottom {
			off.X = -dist
		} else {
			off.Y = -dist
		}
	case "global_label":
		color = schLabelGlobal
		e.v = alignMiddle
		shape := it.str("shape")
		horiz := p.c.round(p.c.sch.labelBox * textH)
		vert := p.c.trunc(textH * 0.0715)
		switch shape {
		case "input", "bidirectional", "tri_state":
			horiz += p.c.trunc(textH * 3 / 4)
		}
		switch sp {
		case spinLeft:
			off = cad.Point{X: -horiz, Y: vert}
		case spinUp:
			off = cad.Point{X: vert, Y: -horiz}
		case spinRight:
			off = cad.Point{X: horiz, Y: vert}
		case spinBottom:
			off = cad.Point{X: vert, Y: horiz}
		}
		if e.hasColor {
			color = e.color
		}
		p.globalShape(x, y, s, e, shape, sp, penW, color)
	case "hierarchical_label":
		color = schLabelHier
		e.v = alignMiddle
		dist := p.c.round(p.c.sch.textOffset*textH) + e.sizeX
		switch sp {
		case spinLeft:
			off.X = -dist
		case spinUp:
			off.Y = -dist
		case spinRight:
			off.X = dist
		case spinBottom:
			off.Y = dist
		}
		if e.hasColor {
			color = e.color
		}
		p.hierShape(x, y, textH, it.str("shape"), sp, penW, color)
	}
	if e.hasColor {
		color = e.color
	}
	p.c.text(p.d, placedText{s: s, x: x + off.X, y: y + off.Y, angle: angle, e: e, color: color, thickness: penW, brk: cad.BreakBox})
	if it.name == "global_label" {
		for _, f := range it.children("property") {
			p.labelField(f, schFields, off)
		}
	} else {
		for _, f := range it.children("property") {
			p.labelField(f, schFields, cad.Point{})
		}
	}
}

// round and trunc round to the internal unit of the file being drawn.
func (c *conv) round(v float64) float64 { return math.Round(v*c.iu) / c.iu }
func (c *conv) trunc(v float64) float64 { return math.Trunc(v*c.iu) / c.iu }

// globalShape draws the outline of a global label.
func (p *schPage) globalShape(x, y float64, s string, e effects, shape string, sp spin, penW float64, color bdf.Color) {
	margin := p.c.round(p.c.sch.labelBox * e.sizeY)
	half := p.c.trunc(e.sizeY/2) + margin
	te := e
	te.h, te.v = alignLeft, alignTop
	_, _, tw, _ := p.c.textBox(s, te)
	symLen := tw + 2*margin
	// the outline is sized by the text's own pen (SCH_TEXT::GetPenWidth)
	textPen := penWidth(e, 0)
	gx := symLen + textPen + 3/p.c.iu
	gy := half + textPen + 3/p.c.iu
	pts := []cad.Point{{X: 0, Y: 0}, {X: 0, Y: -gy}, {X: -gx, Y: -gy}, {X: -gx, Y: 0}, {X: -gx, Y: gy}, {X: 0, Y: gy}}
	xoff := 0.0
	switch shape {
	case "input":
		xoff = -half
		pts[0].X += half
	case "output":
		pts[3].X -= half
	case "bidirectional", "tri_state":
		xoff = -half
		pts[0].X += half
		pts[3].X -= half
	}
	out := make([]cad.Point, 0, len(pts)+1)
	for _, q := range pts {
		q.X += xoff
		q = rotateSpin(q, sp)
		out = append(out, cad.Point{X: q.X + x, Y: q.Y + y})
	}
	out = append(out, out[0])
	p.d.Stroke((&cad.Path{}).Polyline(out, false), p.pen(penW, "solid", color))
}

// hierTemplates are the outlines of hierarchical labels and sheet pins by
// shape and spin (left, up, right, bottom), in halves of the text height.
var hierTemplates = map[string][4][][2]float64{
	"input": {
		{{0, 0}, {-1, -1}, {-2, -1}, {-2, 1}, {-1, 1}, {0, 0}},
		{{0, 0}, {1, -1}, {1, -2}, {-1, -2}, {-1, -1}, {0, 0}},
		{{0, 0}, {1, 1}, {2, 1}, {2, -1}, {1, -1}, {0, 0}},
		{{0, 0}, {1, 1}, {1, 2}, {-1, 2}, {-1, 1}, {0, 0}},
	},
	"output": {
		{{-2, 0}, {-1, 1}, {0, 1}, {0, -1}, {-1, -1}, {-2, 0}},
		{{0, -2}, {1, -1}, {1, 0}, {-1, 0}, {-1, -1}, {0, -2}},
		{{2, 0}, {1, -1}, {0, -1}, {0, 1}, {1, 1}, {2, 0}},
		{{0, 2}, {1, 1}, {1, 0}, {-1, 0}, {-1, 1}, {0, 2}},
	},
	"bidirectional": {
		{{0, 0}, {-1, -1}, {-2, 0}, {-1, 1}, {0, 0}},
		{{0, 0}, {-1, -1}, {0, -2}, {1, -1}, {0, 0}},
		{{0, 0}, {1, -1}, {2, 0}, {1, 1}, {0, 0}},
		{{0, 0}, {-1, 1}, {0, 2}, {1, 1}, {0, 0}},
	},
	"passive": {
		{{0, -1}, {-2, -1}, {-2, 1}, {0, 1}, {0, -1}},
		{{1, 0}, {1, -2}, {-1, -2}, {-1, 0}, {1, 0}},
		{{0, -1}, {2, -1}, {2, 1}, {0, 1}, {0, -1}},
		{{1, 0}, {1, 2}, {-1, 2}, {-1, 0}, {1, 0}},
	},
}

// hierShape draws the outline of a hierarchical label or sheet pin,
// filled with the background.
func (p *schPage) hierShape(x, y, textH float64, shape string, sp spin, penW float64, color bdf.Color) {
	if shape == "tri_state" {
		shape = "bidirectional"
	}
	t, ok := hierTemplates[shape]
	if !ok {
		t = hierTemplates["passive"]
	}
	half := p.c.trunc(textH / 2)
	var pts []cad.Point
	for _, q := range t[sp] {
		pts = append(pts, cad.Point{X: x + q[0]*half, Y: y + q[1]*half})
	}
	path := (&cad.Path{}).Polyline(pts, false)
	p.d.Fill(path, cad.Fill{Color: schBackground}, false)
	p.d.Stroke(path, p.pen(penW, "solid", color))
}

// directive draws a directive label (a net class flag).
func (p *schPage) directive(it *node) {
	e := parseEffects(it.child("effects"))
	at := it.child("at")
	x, y := at.xy()
	sp := labelSpin(at.num(2), e.h)
	color := schNetclassFlag
	if e.hasColor {
		color = e.color
	}
	length := it.numOf("length", schDirectiveLen)
	sym := schDirectiveSym
	penW := penWidth(e, p.c.sch.line)
	tr := func(q cad.Point) cad.Point {
		q = rotateSpin(q, sp)
		return cad.Point{X: q.X + x, Y: q.Y + y}
	}
	pn := p.pen(penW, "solid", color)
	switch shape := it.str("shape"); shape {
	case "dot", "round":
		s := sym
		if shape == "dot" {
			s = p.c.round(sym * 0.7)
		}
		a, b, cc := tr(cad.Point{}), tr(cad.Point{Y: length - s}), tr(cad.Point{Y: length})
		p.d.Stroke((&cad.Path{}).Polyline([]cad.Point{a, b}, false), pn)
		circle := (&cad.Path{}).Circle(cc, s)
		if shape == "dot" {
			p.d.Fill(circle, cad.Fill{Color: color}, false)
		}
		p.d.Stroke(circle, pn)
	case "diamond":
		pts := []cad.Point{{X: 0, Y: 0}, {X: 0, Y: length - sym}, {X: -2 * sym, Y: length}, {X: 0, Y: length + sym}, {X: 2 * sym, Y: length}, {X: 0, Y: length - sym}, {X: 0, Y: 0}}
		for i := range pts {
			pts[i] = tr(pts[i])
		}
		p.d.Stroke((&cad.Path{}).Polyline(pts, false), pn)
	case "rectangle":
		s := p.c.round(sym * 0.8)
		pts := []cad.Point{{X: 0, Y: 0}, {X: 0, Y: length - s}, {X: -2 * s, Y: length - s}, {X: -2 * s, Y: length + s}, {X: 2 * s, Y: length + s}, {X: 2 * s, Y: length - s}, {X: 0, Y: length - s}, {X: 0, Y: 0}}
		for i := range pts {
			pts[i] = tr(pts[i])
		}
		p.d.Stroke((&cad.Path{}).Polyline(pts, false), pn)
	}
	for _, f := range it.children("property") {
		p.labelField(f, color, cad.Point{})
	}
}

// labelField draws a field of a label (its intersheet references, a net
// class): at its position, centered on its box, offset like the label's
// text for global labels.
func (p *schPage) labelField(f *node, color bdf.Color, off cad.Point) {
	e := parseEffects(f.child("effects"))
	if e.hide || f.flag("hide") {
		return
	}
	name, value := f.arg(0), f.arg(1)
	if name == "Intersheetrefs" || name == "Intersheet References" {
		// the pages the label appears on; KiCad fills them in
		value = strings.ReplaceAll(value, "${INTERSHEET_REFS}", "")
	}
	s := p.expand(value)
	if strings.TrimSpace(s) == "" {
		return
	}
	if f.flag("show_name") {
		s = name + ": " + s
	}
	if e.hasColor {
		color = e.color
	}
	at := f.child("at")
	x, y := at.xy()
	angle := at.num(2)
	// plotted where it is, as it is justified (SCH_FIELD::Plot)
	e.mirror = false
	p.c.text(p.d, placedText{s: s, x: x + off.X, y: y + off.Y, angle: keepUpright(angle), e: e, color: color,
		thickness: penWidth(e, p.c.sch.line), brk: cad.BreakBox})
}

// rotatePt turns a point about the origin by a KiCad angle (degrees,
// counterclockwise on the screen, y down).
func rotatePt(x, y, angle float64) (float64, float64) {
	switch int(math.Round(angle)) % 360 {
	case 0:
		return x, y
	case 90, -270:
		return y, -x
	case 180, -180:
		return -x, -y
	case 270, -90:
		return -y, x
	}
	s, c := math.Sincos(angle * math.Pi / 180)
	return x*c + y*s, -x*s + y*c
}
