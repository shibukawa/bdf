package hpgl

import (
	"math"
	"strings"
	"unicode/utf8"

	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// The HP-GL/2 state. Coordinates in the drawing are "drawing units": for a
// plotter, the plotter units of the default coordinate system (x along the
// longer side of the plot size, y up, the origin at the lower left corner
// of the hard-clip limits); for a PCL printer, plotter units on the
// physical page with y down (base places the picture frame there). RO
// turns the plotter's coordinate system, and SC puts user units on it.

// defaultHardW and defaultHardH are the hard-clip limits of a plot that
// sets no plot size: ISO A4 landscape.
const (
	defaultHardW = 11880
	defaultHardH = 8400
)

// scaling is an SC instruction.
type scaling struct {
	typ                    int
	xmin, xmax, ymin, ymax float64 // type 2: xmin, xfactor, ymin, yfactor
	left, bottom           float64 // type 1: percentages of the space left
}

type gl struct {
	c *converter

	// framed is set when HP-GL/2 is placed in the picture frame of a PCL
	// printer, standaloneSet when a plotter took the job (PJL entered
	// HP-GL/2 or HP RTL, or ESC %-1B)
	framed, standaloneSet bool
	frameAt               frameInfo
	grayscale             bool

	base         canvas.Matrix // unrotated plotter units -> drawing
	hardW, hardH float64       // hard-clip limits (unrotated plotter units)
	psSet        bool
	rot          int
	p1, p2       cad.Point // scaling points (current plotter units)
	scale        *scaling
	dev, um, inv canvas.Matrix // plotter units -> drawing; user units -> drawing; drawing -> user units

	pos      cad.Point // pen location (drawing)
	down     bool
	relative bool
	crPoint  cad.Point // carriage-return point of labels (drawing)

	pens     penState
	line     lineState
	fill     fillState
	window   *window
	chord    int // CT: 0 chord angles, 1 deviation distances
	symbol   []byte
	polygon  polygon
	inPoly   bool
	label    labelState
	stroke   strokeState
	anchor   cad.Point // AC (drawing)
	anchored bool
	merge    bool
	ticks    [2]float64 // TL
}

func newGL(c *converter) *gl {
	g := &gl{c: c, base: canvas.Identity, hardW: defaultHardW, hardH: defaultHardH}
	g.reset()
	return g
}

// reset returns everything to its power-on state (ESC E).
func (g *gl) reset() {
	g.flush()
	g.pens = newPens()
	if !g.framed {
		g.hardW, g.hardH, g.psSet = defaultHardW, defaultHardH, false
	}
	g.initialize()
}

// standalone makes HP-GL/2 that of a plotter (PJL entered HP-GL/2, or
// ESC %-1B).
func (g *gl) standalone() {
	g.standaloneSet = true
	if g.framed {
		g.flush()
		g.framed = false
		g.base = canvas.Identity
		g.hardW, g.hardH, g.psSet = defaultHardW, defaultHardH, false
		g.defaultP1P2()
		g.update()
	}
}

// enterPrinter makes the job that of a PCL printer.
func (g *gl) enterPrinter() {
	g.c.pclSetup()
}

// enter is ESC % # B: HP-GL/2 is entered from PCL or HP RTL; 1 and 3 move
// the pen to the cursor.
func (g *gl) enter(v int) {
	c := g.c
	switch {
	case v < 0:
		g.standalone()
	case c.printer():
		// a PCL printer: HP-GL/2 is drawn in the picture frame
		if f := c.pcl.frame(); !g.framed || f != g.frameAt {
			g.framed = true
			g.placeFrame(f)
		}
	}
	if v == 1 || v == 3 {
		if q, ok := c.rtl.capPoint(); ok {
			g.pos = q
			g.crPoint = q
		}
	}
}

// placeFrame puts the coordinate system in the picture frame of the PCL
// page.
func (g *gl) placeFrame(f frameInfo) {
	g.flush()
	g.frameAt = f
	g.base = f.base
	g.hardW, g.hardH = f.w, f.h
	g.psSet = true
	g.defaultP1P2()
	g.update()
}

// initialize is IN: DF, and the pen, plot rotation, P1 and P2, pen widths.
func (g *gl) initialize() {
	g.defaults()
	g.rot = 0
	g.down = false
	g.pens.defaults()
	g.defaultP1P2()
	g.update()
	g.pos = apply(g.dev, cad.Point{})
	g.crPoint = g.pos
	g.merge = false
	g.ticks = [2]float64{0.5, 0.5}
}

// defaults is DF.
func (g *gl) defaults() {
	g.flush()
	g.closeWindow()
	g.window = nil
	g.anchored = false
	g.relative = false
	g.scale = nil
	g.chord = 0
	g.symbol = nil
	g.inPoly = false
	g.polygon = polygon{}
	g.line.defaults()
	g.fill.defaults()
	g.label.defaults()
	g.pens.defaultColors()
	g.update()
	g.crPoint = g.pos
}

// pageAdvanced is the start of a new page: the pen goes to the origin.
func (g *gl) pageAdvanced() {
	g.pos = apply(g.dev, cad.Point{})
	g.crPoint = g.pos
	g.down = false
	g.stroke = strokeState{}
}

// hardRect returns the hard-clip limits in the drawing.
func (g *gl) hardRect() cad.Rect {
	return cad.Rect{}.Add(cad.Point{}).Add(cad.Point{X: g.hardW, Y: g.hardH}).Transform(g.base)
}

// size returns the hard-clip limits seen in the current rotation.
func (g *gl) size() (float64, float64) {
	if g.rot == 90 || g.rot == 270 {
		return g.hardH, g.hardW
	}
	return g.hardW, g.hardH
}

func (g *gl) defaultP1P2() {
	w, h := g.size()
	g.p1, g.p2 = cad.Point{}, cad.Point{X: w, Y: h}
}

// rotation returns the transform of the current plotter units into the
// unrotated ones.
func (g *gl) rotation() canvas.Matrix {
	w, h := g.hardW, g.hardH
	switch g.rot {
	case 90:
		return canvas.Matrix{0, 1, -1, 0, w, 0}
	case 180:
		return canvas.Matrix{-1, 0, 0, -1, w, h}
	case 270:
		return canvas.Matrix{0, -1, 1, 0, 0, h}
	}
	return canvas.Identity
}

// update recomputes the transforms after the coordinate system changed.
func (g *gl) update() {
	g.dev = g.base.Mul(g.rotation())
	s := canvas.Identity
	if g.scale != nil {
		if m, ok := g.scale.matrix(g.p1, g.p2); ok {
			s = m
		}
	}
	g.um = g.dev.Mul(s)
	g.inv = invert(g.um)
	g.pens.invalidate()
}

// matrix returns the transform of user units into plotter units.
func (s *scaling) matrix(p1, p2 cad.Point) (canvas.Matrix, bool) {
	switch s.typ {
	case 2:
		if s.xmax == 0 || s.ymax == 0 {
			return canvas.Matrix{}, false
		}
		return canvas.Matrix{s.xmax, 0, 0, s.ymax, p1.X - s.xmin*s.xmax, p1.Y - s.ymin*s.ymax}, true
	case 1:
		uw, uh := s.xmax-s.xmin, s.ymax-s.ymin
		w, h := p2.X-p1.X, p2.Y-p1.Y
		if uw == 0 || uh == 0 || w == 0 || h == 0 {
			return canvas.Matrix{}, false
		}
		k := math.Min(math.Abs(w/uw), math.Abs(h/uh))
		sx := k * sign(w) * sign(uw)
		sy := k * sign(h) * sign(uh)
		// the isotropic area in the P1/P2 rectangle, the space left
		// shared out by left and bottom
		ex := math.Abs(w) - k*math.Abs(uw)
		ey := math.Abs(h) - k*math.Abs(uh)
		lx := math.Min(p1.X, p2.X) + ex*s.left/100
		ly := math.Min(p1.Y, p2.Y) + ey*s.bottom/100
		// xmin goes to the side of P1
		tx := lx - math.Min(0, sx*uw) - sx*s.xmin
		ty := ly - math.Min(0, sy*uh) - sy*s.ymin
		return canvas.Matrix{sx, 0, 0, sy, tx, ty}, true
	default:
		uw, uh := s.xmax-s.xmin, s.ymax-s.ymin
		if uw == 0 || uh == 0 {
			return canvas.Matrix{}, false
		}
		sx, sy := (p2.X-p1.X)/uw, (p2.Y-p1.Y)/uh
		return canvas.Matrix{sx, 0, 0, sy, p1.X - s.xmin*sx, p1.Y - s.ymin*sy}, true
	}
}

func sign(v float64) float64 {
	if v < 0 {
		return -1
	}
	return 1
}

// invert returns the inverse of an affine transform (the identity for a
// degenerate one).
func invert(m canvas.Matrix) canvas.Matrix {
	det := m[0]*m[3] - m[1]*m[2]
	if det == 0 || math.IsNaN(det) || math.IsInf(det, 0) {
		return canvas.Identity
	}
	a, b, c, d := m[3]/det, -m[1]/det, -m[2]/det, m[0]/det
	return canvas.Matrix{a, b, c, d, -(a*m[4] + c*m[5]), -(b*m[4] + d*m[5])}
}

// user returns the pen location in user units.
func (g *gl) user() cad.Point { return apply(g.inv, g.pos) }

func apply(m canvas.Matrix, p cad.Point) cad.Point {
	x, y := m.Apply(p.X, p.Y)
	return cad.Point{X: x, Y: y}
}

// linear applies the linear part of m (a displacement).
func linear(m canvas.Matrix, d cad.Point) cad.Point {
	return cad.Point{X: m[0]*d.X + m[2]*d.Y, Y: m[1]*d.X + m[3]*d.Y}
}

// scaleOf returns how lengths in the drawing compare with plotter units.
func (g *gl) scaleOf() float64 { return cad.Scale(g.base) }

// diagonal returns the distance from P1 to P2 in plotter units.
func (g *gl) diagonal() float64 { return g.p2.Sub(g.p1).Len() }

// valid reports whether a parameter is a usable number.
func valid(vs ...float64) bool {
	for _, v := range vs {
		if math.IsNaN(v) || math.IsInf(v, 0) || math.Abs(v) > 1e12 {
			return false
		}
	}
	return true
}

// instruction reads the parameters of an instruction and carries it out.
func (g *gl) instruction(m string) {
	c := g.c
	switch m {
	case "LB":
		g.labelInstruction(c.labelText())
		return
	case "WD", "BL":
		c.labelText()
		return
	case "DT":
		g.defineTerminator()
		return
	case "SM":
		g.symbolMode()
		return
	case "CO", "MG":
		c.quoted()
		return
	case "BP":
		g.beginPlot()
		return
	case "PE":
		g.polylineEncoded()
		return
	case "PB":
		c.params()
		return
	}
	a := c.params()
	switch m {
	case "PA", "PR", "PU", "PD", "AA", "AR", "AT", "RT", "CI", "BZ", "BR":
	default:
		// the pen may draw differently after anything else
		g.pens.valid = false
	}
	switch m {
	// configuration and status
	case "IN":
		g.initialize()
	case "DF":
		g.defaults()
	case "IP":
		g.inputP1P2(a, false)
	case "IR":
		g.inputP1P2(a, true)
	case "IW":
		g.inputWindow(a)
	case "SC":
		g.scaleInstruction(a)
	case "RO":
		g.rotate(a)
	case "PS":
		g.plotSize(a)
	case "PG", "AF", "AH", "FR":
		g.advance(a)
	case "RP", "QL", "MT", "EC", "NR", "VS", "VA", "VN", "ST", "PP", "OE", "OH", "OI", "OP", "OS", "OD", "OA", "OC", "OF", "OG",
		"OK", "OL", "OO", "OT", "OW", "DC", "DP", "IM", "AP", "AS", "BF", "GC", "GM", "GP", "KY", "SG", "DL", "UC", "CM", "CV", "DS", "IV", "EL", "LI":
	// vectors
	case "PA":
		g.relative = false
		g.plot(a)
	case "PR":
		g.relative = true
		g.plot(a)
	case "PU":
		g.penUp()
		g.plot(a)
	case "PD":
		g.penDown()
		g.plot(a)
	case "AA", "AR":
		g.arcInstruction(a, m == "AR")
	case "AT", "RT":
		g.arcThree(a, m == "RT")
	case "CI":
		g.circle(a)
	case "BZ", "BR":
		g.bezier(a, m == "BR")
	case "CT":
		if len(a) == 0 || a[0] == 0 || a[0] == 1 {
			g.chord = 0
			if len(a) > 0 {
				g.chord = int(a[0])
			}
		}
	// polygons
	case "PM":
		g.polygonMode(a)
	case "EP":
		g.edgePolygon()
	case "FP":
		g.fillPolygon(a)
	case "EA", "ER", "RA", "RR":
		g.rectangle(a, m == "ER" || m == "RR", m[0] == 'R')
	case "EW", "WG":
		g.wedge(a, m == "WG")
	// line and fill attributes
	case "SP":
		g.selectPen(a)
	case "PW":
		g.penWidth(a)
	case "WU":
		g.widthUnit(a)
	case "NP":
		g.numberOfPens(a)
	case "PC":
		g.penColor(a)
	case "CR":
		g.colorRange(a)
	case "LT":
		g.lineType(a)
	case "UL":
		g.userLineType(a)
	case "LA":
		g.lineAttributes(a)
	case "FT":
		g.fillType(a)
	case "PT":
		// pen thickness of solid fills on pen plotters
	case "RF":
		g.rasterFill(a)
	case "AC":
		g.anchorCorner(a)
	case "SV":
		g.screenedVectors(a)
	case "TR":
		g.flush()
		g.fill.transparent = len(a) == 0 || a[0] != 0
	case "MC":
		g.merge = len(a) > 0 && a[0] != 0
		if g.merge && (len(a) < 2 || a[1] != 252) {
			c.warnOnce("MC", "merge control (MC) is not applied: later marks replace earlier ones")
		}
	// characters
	case "SI", "SR", "SD", "AD", "SS", "SA", "DI", "DR", "DV", "LO", "ES", "SL", "CP", "TD", "LM", "CF", "FI", "FN", "SB", "CS", "CA":
		g.characterInstruction(m, a)
	case "XT", "YT":
		g.tick(m == "XT")
	case "TL":
		g.tickLength(a)
	default:
		c.warnOnce("unknown "+m, "the instruction %s is not supported", m)
	}
}

// --- configuration ---

// beginPlot is BP: a new plot starts on a new page, initialized; the
// picture name (kind 1) titles the document.
func (g *gl) beginPlot() {
	c := g.c
	b := c.b
	for c.i < len(b) {
		ch := b[c.i]
		switch {
		case ch == ';':
			c.i++
		case ch == ' ' || ch == ',' || ch == '\t' || ch == '\r' || ch == '\n':
			c.i++
			continue
		case ch >= '0' && ch <= '9' || ch == '+' || ch == '-' || ch == '.':
			v, n := number(b[c.i:])
			c.i += max(n, 1)
			if int(v) == 1 {
				// a quoted picture name follows
				for c.i < len(b) && (b[c.i] == ' ' || b[c.i] == ',') {
					c.i++
				}
				if c.i < len(b) && b[c.i] == '"' {
					if s := trimTitle(c.quoted()); s != "" && c.title == "" {
						c.title = s
					}
				}
			}
			continue
		case ch == '"':
			c.quoted()
			continue
		}
		break
	}
	g.flush()
	if !c.cur.d.Empty() && !c.printer() {
		// a new plot (PCL maps BP to IN: its pages end with a form feed)
		c.rtl.end()
		c.endPage()
	}
	g.initialize()
}

// advance is PG: the page ends.
func (g *gl) advance(a []float64) {
	c := g.c
	if c.printer() {
		return // PCL ejects pages
	}
	c.rtl.end()
	c.endPage()
}

// plotSize is PS: the hard-clip limits, the length along the frame advance
// and the width across it; the x axis is the longer side.
func (g *gl) plotSize(a []float64) {
	if g.c.printer() {
		return // the picture frame is the plot size in PCL
	}
	w, h := float64(defaultHardW), float64(defaultHardH)
	if len(a) >= 1 {
		l, wd := a[0], a[0]
		if len(a) >= 2 {
			wd = a[1]
		} else {
			wd = 0
		}
		if l <= 0 || !valid(l, wd) {
			return
		}
		if wd <= 0 {
			// the width of the media, which a file does not tell: keep the
			// length and the ratio of the roll's default (1.5 : 1)
			wd = l / 1.5
		}
		w, h = math.Max(l, wd), math.Min(l, wd)
		g.c.rtl.plotSize(l, wd)
	}
	g.flush()
	g.closeWindow()
	g.hardW, g.hardH = w, h
	g.psSet = true
	g.defaultP1P2()
	g.window = nil
	g.anchored = false
	g.polygon = polygon{}
	g.inPoly = false
	g.update()
	g.pos = apply(g.dev, cad.Point{})
	g.crPoint = g.pos
}

// inputP1P2 is IP (plotter units) and IR (percentages of the hard-clip
// limits).
func (g *gl) inputP1P2(a []float64, rel bool) {
	g.flush()
	if len(a) == 0 {
		g.defaultP1P2()
		g.update()
		return
	}
	if len(a) != 2 && len(a) < 4 || !valid(a...) {
		return
	}
	w, h := g.size()
	pt := func(x, y float64) cad.Point {
		if rel {
			return cad.Point{X: x / 100 * w, Y: y / 100 * h}
		}
		return cad.Point{X: math.Round(x), Y: math.Round(y)}
	}
	p1 := pt(a[0], a[1])
	p2 := p1.Add(g.p2.Sub(g.p1))
	if len(a) >= 4 {
		p2 = pt(a[2], a[3])
	}
	if p2.X == p1.X {
		p2.X++
	}
	if p2.Y == p1.Y {
		p2.Y++
	}
	g.p1, g.p2 = p1, p2
	g.update()
}

// scaleInstruction is SC.
func (g *gl) scaleInstruction(a []float64) {
	g.flush()
	if len(a) == 0 {
		// scaling off: the window, anchor corner and fill spacing keep their
		// place in plotter units
		g.scale = nil
		g.update()
		return
	}
	if len(a) < 4 || !valid(a...) {
		return
	}
	s := &scaling{xmin: a[0], xmax: a[1], ymin: a[2], ymax: a[3], left: 50, bottom: 50}
	if len(a) >= 5 {
		s.typ = int(a[4])
	}
	switch s.typ {
	case 0:
		if a[0] == a[1] || a[2] == a[3] {
			return
		}
	case 1:
		if a[0] == a[1] || a[2] == a[3] {
			return
		}
		if len(a) >= 7 {
			s.left, s.bottom = math.Max(0, math.Min(100, a[5])), math.Max(0, math.Min(100, a[6]))
		}
	case 2:
		if a[1] == 0 || a[3] == 0 {
			return
		}
	default:
		return
	}
	g.scale = s
	g.update()
}

// rotate is RO: the coordinate system turns; the pen keeps its place on
// the paper, P1 and P2 their coordinates.
func (g *gl) rotate(a []float64) {
	r := 0
	if len(a) > 0 {
		r = int(math.Round(a[0]))
	}
	switch r {
	case 0, 90, 180, 270:
	default:
		return
	}
	g.flush()
	old := g.dev
	g.rot = r
	g.update()
	g.crPoint = g.pos
	// the window and the polygon buffer keep their coordinates: they turn
	// on the paper
	turn := g.dev.Mul(invert(old))
	g.polygon.transform(turn)
	if g.window != nil {
		open := g.window.open
		g.closeWindow()
		g.window.path = g.window.path.Transform(turn)
		if open {
			g.openWindow()
		}
	}
}

// --- vectors ---

// penUp and penDown change the pen state; a pen lowered and raised again
// without moving leaves a dot.
func (g *gl) penUp() {
	if g.down && g.stroke.pendingDot && !g.inPoly {
		g.dot(g.pos)
	}
	g.down = false
	g.stroke.pendingDot = false
}

func (g *gl) penDown() {
	if !g.down {
		g.stroke.pendingDot = true
	}
	g.down = true
}

// plot moves the pen through the coordinate pairs of PA, PR, PU and PD.
func (g *gl) plot(a []float64) {
	for i := 0; i+1 < len(a); i += 2 {
		if !valid(a[i], a[i+1]) {
			continue
		}
		var q cad.Point
		if g.relative {
			q = g.pos.Add(linear(g.um, cad.Point{X: a[i], Y: a[i+1]}))
		} else {
			q = apply(g.um, cad.Point{X: a[i], Y: a[i+1]})
		}
		g.lineTo(q)
		g.drawSymbol()
	}
	if len(a) >= 2 {
		g.crPoint = g.pos
	}
}

// lineTo moves the pen to q, drawing when it is down; a vector of no
// length makes a dot.
func (g *gl) lineTo(q cad.Point) {
	if q == g.pos && g.down && !g.inPoly {
		g.stroke.pendingDot = false
		g.dot(q)
		return
	}
	p := (&cad.Path{}).MoveTo(g.pos.X, g.pos.Y).LineTo(q.X, q.Y)
	g.addPath(p, g.down, []cad.Point{q})
	g.pos = q
}

// addPath adds a piece of path that starts at the pen location, drawn when
// down; vertices are its end points (for line type 0 and adaptive line
// types).
func (g *gl) addPath(p *cad.Path, down bool, vertices []cad.Point) {
	if g.inPoly {
		g.polygon.add(p, down)
		return
	}
	if !down {
		return
	}
	g.stroke.pendingDot = false
	g.strokePath(p, vertices)
}

// arcInstruction is AA and AR: an arc around a centre from the pen
// location.
func (g *gl) arcInstruction(a []float64, rel bool) {
	if len(a) < 3 || !valid(a...) {
		return
	}
	u := g.user()
	ctr := cad.Point{X: a[0], Y: a[1]}
	if rel {
		ctr = u.Add(ctr)
	}
	sweep := math.Mod(a[2], 360)
	if a[2] != 0 && sweep == 0 {
		sweep = 360 * sign(a[2])
	}
	chord, explicit := g.chordParam(a, 3)
	r := u.Sub(ctr).Len()
	a0 := math.Atan2(u.Y-ctr.Y, u.X-ctr.X)
	if r == 0 || sweep == 0 {
		if g.down && !g.inPoly {
			g.dot(g.pos)
		}
		g.crPoint = g.pos
		return
	}
	p := g.arcPath(ctr, r, a0, a0+sweep*math.Pi/180, chord, explicit)
	g.addUserPath(p)
	g.crPoint = g.pos
}

// chordParam returns an arc's chord angle (degrees) or deviation, and
// whether the instruction gave it.
func (g *gl) chordParam(a []float64, i int) (float64, bool) {
	if i < len(a) {
		return math.Abs(a[i]), true
	}
	return 5, false
}

// arcPath returns an arc in user units, from angle a0 to a1 (radians),
// starting with a move. A chord angle coarser than the default draws the
// chords a plotter draws; finer ones, and the default, draw the arc.
func (g *gl) arcPath(ctr cad.Point, r, a0, a1, chord float64, explicit bool) *cad.Path {
	p := &cad.Path{}
	n := 0
	if explicit && g.chord == 0 {
		ca := math.Mod(chord, 360)
		if ca > 180 {
			ca = 360 - ca
		}
		if ca > 5+1e-9 {
			n = int(math.Ceil(math.Abs(a1-a0)*180/math.Pi/ca - 1e-9))
		}
	}
	if n >= 1 && n <= 3600 {
		for i := 0; i <= n; i++ {
			t := a0 + (a1-a0)*float64(i)/float64(n)
			q := cad.Point{X: ctr.X + r*math.Cos(t), Y: ctr.Y + r*math.Sin(t)}
			if i == 0 {
				p.MoveTo(q.X, q.Y)
			} else {
				p.LineTo(q.X, q.Y)
			}
		}
		return p
	}
	return p.Arc(ctr, r, a0, a1)
}

// addUserPath adds a path in user units that starts at the pen location;
// the pen moves to its end.
func (g *gl) addUserPath(p *cad.Path) {
	d := p.Transform(g.um)
	end := d.Current()
	g.addPath(d, g.down, []cad.Point{end})
	g.pos = end
}

// arcThree is AT and RT: an arc through an intermediate point to an end
// point.
func (g *gl) arcThree(a []float64, rel bool) {
	if len(a) < 4 || !valid(a...) {
		return
	}
	u := g.user()
	m := cad.Point{X: a[0], Y: a[1]}
	e := cad.Point{X: a[2], Y: a[3]}
	if rel {
		m, e = u.Add(m), u.Add(e)
	}
	chord, explicit := g.chordParam(a, 4)
	defer func() { g.crPoint = g.pos }()
	switch {
	case m == u && e == u:
		if g.down && !g.inPoly {
			g.dot(g.pos)
		}
		return
	case m == u || m == e:
		g.addUserPath((&cad.Path{}).MoveTo(u.X, u.Y).LineTo(e.X, e.Y))
		return
	case e == u:
		// a circle with the diameter from the pen to the intermediate point
		ctr := u.Add(m).Mul(0.5)
		a0 := math.Atan2(u.Y-ctr.Y, u.X-ctr.X)
		g.addUserPath(g.arcPath(ctr, u.Sub(ctr).Len(), a0, a0+2*math.Pi, chord, explicit))
		return
	}
	ctr, ok := circumcentre(u, m, e)
	if !ok {
		// collinear: a straight line
		g.addUserPath((&cad.Path{}).MoveTo(u.X, u.Y).LineTo(e.X, e.Y))
		return
	}
	r := u.Sub(ctr).Len()
	a0 := math.Atan2(u.Y-ctr.Y, u.X-ctr.X)
	am := math.Atan2(m.Y-ctr.Y, m.X-ctr.X)
	a1 := math.Atan2(e.Y-ctr.Y, e.X-ctr.X)
	// counter-clockwise when the intermediate point comes before the end
	ccw := norm(am-a0) < norm(a1-a0)
	sweep := norm(a1 - a0)
	if !ccw {
		sweep -= 2 * math.Pi
	}
	g.addUserPath(g.arcPath(ctr, r, a0, a0+sweep, chord, explicit))
}

// norm returns an angle in [0, 2π).
func norm(a float64) float64 {
	a = math.Mod(a, 2*math.Pi)
	if a < 0 {
		a += 2 * math.Pi
	}
	return a
}

// circumcentre returns the centre of the circle through three points.
func circumcentre(a, b, c cad.Point) (cad.Point, bool) {
	d := 2 * (a.X*(b.Y-c.Y) + b.X*(c.Y-a.Y) + c.X*(a.Y-b.Y))
	scale := math.Max(math.Max(b.Sub(a).Len(), c.Sub(a).Len()), 1e-300)
	if math.Abs(d) < 1e-12*scale*scale {
		return cad.Point{}, false
	}
	a2, b2, c2 := a.X*a.X+a.Y*a.Y, b.X*b.X+b.Y*b.Y, c.X*c.X+c.Y*c.Y
	return cad.Point{
		X: (a2*(b.Y-c.Y) + b2*(c.Y-a.Y) + c2*(a.Y-b.Y)) / d,
		Y: (a2*(c.X-b.X) + b2*(a.X-c.X) + c2*(b.X-a.X)) / d,
	}, true
}

// circle is CI: a circle around the pen location, drawn with the pen down
// whatever its state; the pen stays at the centre.
func (g *gl) circle(a []float64) {
	if len(a) < 1 || !valid(a...) {
		return
	}
	u := g.user()
	r := a[0]
	chord, explicit := g.chordParam(a, 1)
	start := 0.0
	if r < 0 {
		start, r = math.Pi, -r
	}
	if g.inPoly {
		g.polygon.closeSub(g.down)
	}
	if r == 0 {
		if !g.inPoly {
			g.dot(g.pos)
		}
		return
	}
	p := g.arcPath(u, r, start, start+2*math.Pi, chord, explicit).Close()
	d := p.Transform(g.um)
	if g.inPoly {
		g.polygon.addClosed(d, g.down)
		return
	}
	first := d.First()
	g.strokePath(d, []cad.Point{first})
}

// bezier is BZ and BR: cubic Bézier curves from the pen location.
func (g *gl) bezier(a []float64, rel bool) {
	for i := 0; i+5 < len(a); i += 6 {
		if !valid(a[i : i+6]...) {
			continue
		}
		u := g.user()
		pts := [3]cad.Point{{X: a[i], Y: a[i+1]}, {X: a[i+2], Y: a[i+3]}, {X: a[i+4], Y: a[i+5]}}
		if rel {
			for j := range pts {
				pts[j] = u.Add(pts[j])
			}
		}
		p := (&cad.Path{}).MoveTo(u.X, u.Y).CubicTo(pts[0], pts[1], pts[2])
		g.addUserPath(p)
	}
	if len(a) >= 6 {
		g.crPoint = g.pos
	}
}

// polylineEncoded is PE: coordinates encoded in base 64 (or base 32 after
// the 7 flag), relative unless flagged absolute, drawn with the pen down
// unless flagged up.
func (g *gl) polylineEncoded() {
	c := g.c
	b := c.b
	base := 64
	fraction := 0.0 // fractional binary bits
	upNext, absNext := false, false
	wasRelative := g.relative
	moved := false
	readNumber := func() (int64, bool) {
		var v int64
		mul := int64(1)
		for c.i < len(b) {
			ch := b[c.i]
			if ch == ';' || ch == 0x1b {
				return 0, false
			}
			c.i++
			if base == 64 {
				ch &= 0xff
			} else {
				ch &= 0x7f
			}
			switch {
			case base == 64 && ch >= 63 && ch <= 126:
				v += int64(ch-63) * mul
			case base == 64 && ch >= 191 && ch <= 254:
				v += int64(ch-191) * mul
				return v, true
			case base == 32 && ch >= 63 && ch <= 94:
				v += int64(ch-63) * mul
			case base == 32 && ch >= 95 && ch <= 126:
				v += int64(ch-95) * mul
				return v, true
			default:
				continue // ignored characters
			}
			if mul > 1<<40 {
				return 0, false
			}
			mul *= int64(base)
		}
		return 0, false
	}
	signed := func(v int64) float64 {
		if v&1 != 0 {
			return -float64(v >> 1)
		}
		return float64(v >> 1)
	}
	for c.i < len(b) {
		ch := b[c.i]
		if ch == ';' {
			c.i++
			break
		}
		if ch == 0x1b {
			break
		}
		flag := ch & 0x7f
		switch flag {
		case ':':
			// values are encoded as coordinates are, with a sign
			c.i++
			if v, ok := readNumber(); ok && !g.inPoly {
				g.selectPen([]float64{signed(v)})
			}
			continue
		case '<':
			c.i++
			upNext = true
			continue
		case '>':
			c.i++
			if v, ok := readNumber(); ok {
				if f := signed(v); f >= 0 && f < 64 {
					fraction = f
				}
			}
			continue
		case '=':
			c.i++
			absNext = true
			continue
		case '7':
			c.i++
			base = 32
			continue
		}
		if flag < 63 {
			c.i++ // spaces and control characters are ignored
			continue
		}
		xv, ok := readNumber()
		if !ok {
			break
		}
		yv, ok := readNumber()
		if !ok {
			break
		}
		k := math.Pow(2, -fraction)
		x, y := signed(xv)*k, signed(yv)*k
		var q cad.Point
		if absNext {
			q = apply(g.um, cad.Point{X: x, Y: y})
		} else {
			q = g.pos.Add(linear(g.um, cad.Point{X: x, Y: y}))
		}
		if upNext {
			g.penUp()
		} else {
			g.penDown()
		}
		g.lineTo(q)
		g.drawSymbol()
		upNext, absNext = false, false
		moved = true
	}
	g.relative = wasRelative
	if moved {
		g.crPoint = g.pos
	}
}

// --- tick marks (HP-GL) ---

// tickLength is TL: the lengths of tick marks on the positive and negative
// sides, in percentages of P2 − P1.
func (g *gl) tickLength(a []float64) {
	g.ticks = [2]float64{0.5, 0.5}
	if len(a) > 0 && valid(a[0]) {
		g.ticks[0], g.ticks[1] = a[0], 0
	}
	if len(a) > 1 && valid(a[1]) {
		g.ticks[1] = a[1]
	}
}

// tick is XT and YT: a tick mark across the pen location, perpendicular to
// the x axis (XT) or the y axis (YT).
func (g *gl) tick(x bool) {
	d := g.p2.Sub(g.p1)
	var a, b cad.Point
	if x {
		a, b = cad.Point{Y: g.ticks[0] / 100 * d.Y}, cad.Point{Y: -g.ticks[1] / 100 * d.Y}
	} else {
		a, b = cad.Point{X: g.ticks[0] / 100 * d.X}, cad.Point{X: -g.ticks[1] / 100 * d.X}
	}
	p0 := g.pos.Add(linear(g.dev, b))
	p1 := g.pos.Add(linear(g.dev, a))
	g.strokePath((&cad.Path{}).MoveTo(p0.X, p0.Y).LineTo(p1.X, p1.Y), nil)
}

// trimTitle cleans up a picture name for a title.
func trimTitle(s string) string {
	s = strings.TrimSpace(strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, s))
	if !utf8.ValidString(s) {
		s = strings.ToValidUTF8(s, "")
	}
	return s
}
