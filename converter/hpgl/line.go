package hpgl

import (
	"math"
	"slices"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
)

// --- pens (the palette) ---

// penState is the palette: the color and width of each pen, and the pen
// selected.
type penState struct {
	n        int // palette size (NP)
	colors   []bdf.Color
	widths   []float64 // mm (WU 0) or percentages of the P1–P2 diagonal (WU 1)
	relative bool      // WU 1
	cr       [6]float64
	selected int
	// chosen records an SP: a plot that never selects a pen is drawn with
	// pen 1, as its author must have meant (plotters draw nothing without
	// a pen, and cutters cut without one)
	chosen bool
	valid  bool // the cached pen is up to date
	cached cad.Pen
}

// defaultColors are the first eight pens of the HP-GL/2 palette.
var defaultColors = [8]bdf.Color{
	bdf.RGB(255, 255, 255), bdf.RGB(0, 0, 0), bdf.RGB(255, 0, 0), bdf.RGB(0, 255, 0),
	bdf.RGB(255, 255, 0), bdf.RGB(0, 0, 255), bdf.RGB(255, 0, 255), bdf.RGB(0, 255, 255),
}

const defaultWidth = 0.35 // mm

func newPens() penState {
	p := penState{}
	p.defaults()
	return p
}

// defaults is the palette of IN: eight pens in their default colors, 0.35
// mm wide.
func (p *penState) defaults() {
	p.n = 8
	p.colors = make([]bdf.Color, 8)
	copy(p.colors, defaultColors[:])
	p.widths = []float64{defaultWidth, defaultWidth, defaultWidth, defaultWidth, defaultWidth, defaultWidth, defaultWidth, defaultWidth}
	p.relative = false
	p.defaultColorRange()
	p.valid = false
}

func (p *penState) defaultColorRange() { p.cr = [6]float64{0, 255, 0, 255, 0, 255} }

// defaultColors is the palette of DF: eight pens in their default colors,
// keeping the widths of the pens that remain.
func (p *penState) defaultColors() {
	p.resize(8)
	p.colors = p.colors[:8]
	p.widths = p.widths[:8]
	for i := range p.colors {
		p.colors[i] = defaultColor(i)
	}
	p.defaultColorRange()
	p.valid = false
}

func (p *penState) invalidate() { p.valid = false }

// index maps a pen number into the palette (raster devices take pens past
// its end modulo its size less pen 0).
func (p *penState) index(n int) int {
	if n < 0 {
		return 0
	}
	if n >= p.n {
		if p.n <= 1 {
			return 0
		}
		n = (n-1)%(p.n-1) + 1
	}
	return n
}

func (p *penState) current() int {
	if !p.chosen && p.selected == 0 {
		return p.index(1)
	}
	return p.index(p.selected)
}

// defaultColor is the color a pen takes by default.
func defaultColor(n int) bdf.Color {
	if n < 8 {
		return defaultColors[n]
	}
	return bdf.RGB(0, 0, 0)
}

// resize sets the palette size, keeping the pens that remain.
func (p *penState) resize(n int) {
	for len(p.colors) < n {
		p.colors = append(p.colors, defaultColor(len(p.colors)))
		w := defaultWidth
		if p.relative {
			w = 0.1
		}
		p.widths = append(p.widths, w)
	}
	p.n = n
	p.valid = false
}

func (g *gl) selectPen(a []float64) {
	g.flush()
	n := 0
	if len(a) > 0 && valid(a[0]) {
		n = int(math.Round(math.Max(-1e6, math.Min(1e6, a[0]))))
	}
	g.pens.selected = n
	g.pens.chosen = true
	g.pens.valid = false
}

func (g *gl) penWidth(a []float64) {
	g.flush()
	p := &g.pens
	w := defaultWidth
	if p.relative {
		w = 0.1
	}
	if len(a) > 0 {
		if !valid(a[0]) || a[0] < 0 {
			return
		}
		w = a[0]
	}
	if len(a) > 1 {
		n := int(a[1])
		if n < 0 || n >= p.n {
			return
		}
		p.widths[n] = w
	} else {
		for i := range p.widths {
			p.widths[i] = w
		}
	}
	p.valid = false
}

func (g *gl) widthUnit(a []float64) {
	g.flush()
	p := &g.pens
	rel := len(a) > 0 && a[0] == 1
	if len(a) > 0 && a[0] != 0 && a[0] != 1 {
		return
	}
	p.relative = rel
	w := defaultWidth
	if rel {
		w = 0.1
	}
	for i := range p.widths {
		p.widths[i] = w
	}
	p.valid = false
}

func (g *gl) numberOfPens(a []float64) {
	g.flush()
	n := 8
	if len(a) > 0 {
		if a[0] < 2 || !valid(a[0]) {
			return
		}
		// the next power of two, at most 256
		n = 2
		for n < int(a[0]) && n < 256 {
			n *= 2
		}
	}
	g.pens.resize(n)
}

func (g *gl) penColor(a []float64) {
	g.flush()
	p := &g.pens
	if len(a) == 0 {
		for i := range p.colors {
			p.colors[i] = defaultColor(i)
		}
		p.valid = false
		return
	}
	n := int(a[0])
	if n < 0 || n >= p.n || len(a) == 2 || len(a) == 3 || !valid(a...) {
		return
	}
	if len(a) == 1 {
		p.colors[n] = defaultColor(n)
	} else {
		var rgb [3]uint8
		for i := 0; i < 3; i++ {
			lo, hi := p.cr[2*i], p.cr[2*i+1]
			v := a[1+i]
			// clamp to the range, which may run either way
			if lo < hi {
				v = math.Max(lo, math.Min(hi, v))
			} else {
				v = math.Max(hi, math.Min(lo, v))
			}
			rgb[i] = uint8(math.Round((v - lo) / (hi - lo) * 255))
		}
		p.colors[n] = bdf.RGB(rgb[0], rgb[1], rgb[2])
	}
	p.valid = false
}

func (g *gl) colorRange(a []float64) {
	if len(a) == 0 {
		g.pens.defaultColorRange()
		return
	}
	if len(a) < 6 || !valid(a...) || a[0] == a[1] || a[2] == a[3] || a[4] == a[5] {
		return
	}
	copy(g.pens.cr[:], a[:6])
}

// penColorOf returns the color a pen draws in.
func (g *gl) penColorOf(n int) bdf.Color {
	col := g.pens.colors[g.pens.index(n)]
	if g.c.opts.Colors == Mono || g.grayscale {
		if g.pens.index(n) == 0 {
			return bdf.RGB(255, 255, 255)
		}
		if g.c.opts.Colors == Mono {
			return bdf.RGB(0, 0, 0)
		}
		r, gg, b := float64(uint8(col>>24)), float64(uint8(col>>16)), float64(uint8(col>>8))
		y := uint8(math.Round(0.299*r + 0.587*gg + 0.114*b))
		return bdf.RGB(y, y, y)
	}
	return col
}

// widthOf returns the width of a pen in drawing units (0: the thinnest).
func (g *gl) widthOf(n int) float64 {
	w := g.pens.widths[g.pens.index(n)]
	if g.pens.relative {
		return w / 100 * g.diagonal() * g.scaleOf()
	}
	return w * 40 * g.scaleOf()
}

// --- line types and attributes ---

// lineType is an LT setting.
type lineType struct {
	typ     int // 0 dots at the end points, ±1..8; solid when !set
	set     bool
	length  float64 // pattern length: percentage of the P1–P2 diagonal, or mm
	absolut bool
}

type lineState struct {
	lt, prev  lineType
	user      [9][]float64 // UL patterns
	cap, join byte
	miter     float64
}

// defaultPatterns are the fixed line types 1 to 8: pen-down and pen-up
// stretches in percentages of the pattern length.
var defaultPatterns = [9][]float64{
	nil,
	{0, 100},
	{50, 50},
	{70, 30},
	{80, 10, 0, 10},
	{70, 10, 10, 10},
	{50, 10, 10, 10, 10, 10},
	{70, 10, 0, 10, 0, 10},
	{50, 10, 0, 10, 10, 10, 0, 10},
}

func (l *lineState) defaults() {
	l.lt = lineType{length: 4}
	l.prev = lineType{}
	for i := range l.user {
		l.user[i] = nil
	}
	l.cap, l.join, l.miter = cad.CapButt, cad.JoinMiter, 5
}

func (l *lineState) pattern(typ int) []float64 {
	if typ < 0 {
		typ = -typ
	}
	if typ < 1 || typ > 8 {
		return nil
	}
	if p := l.user[typ]; p != nil {
		return p
	}
	return defaultPatterns[typ]
}

func (g *gl) lineType(a []float64) {
	g.flush()
	l := &g.line
	if len(a) == 0 {
		if l.lt.set {
			l.prev = l.lt
		}
		l.lt.set = false
		return
	}
	t := int(math.Round(a[0]))
	if t == 99 {
		if !l.lt.set && l.prev.set {
			l.lt = l.prev
		}
		return
	}
	if t < -8 || t > 8 || !valid(a...) {
		return
	}
	lt := l.lt
	lt.typ, lt.set = t, true
	if len(a) > 1 {
		if a[1] <= 0 {
			return
		}
		lt.length = a[1]
	}
	if len(a) > 2 {
		switch a[2] {
		case 0:
			lt.absolut = false
		case 1:
			lt.absolut = true
		default:
			return
		}
	}
	l.lt = lt
	l.prev = lineType{}
}

func (g *gl) userLineType(a []float64) {
	g.flush()
	l := &g.line
	if len(a) == 0 {
		for i := range l.user {
			l.user[i] = nil
		}
		return
	}
	n := int(math.Abs(a[0]))
	if n < 1 || n > 8 || !valid(a...) {
		return
	}
	if len(a) == 1 {
		l.user[n] = nil
		return
	}
	gaps := slices.Clone(a[1:min(len(a), 21)])
	sum := 0.0
	for i, v := range gaps {
		if v < 0 {
			gaps[i] = 0
		}
		sum += gaps[i]
	}
	if sum <= 0 {
		return
	}
	for i := range gaps {
		gaps[i] = gaps[i] / sum * 100
	}
	l.user[n] = gaps
}

func (g *gl) lineAttributes(a []float64) {
	g.flush()
	l := &g.line
	if len(a) == 0 {
		l.cap, l.join, l.miter = cad.CapButt, cad.JoinMiter, 5
		return
	}
	for i := 0; i+1 < len(a); i += 2 {
		v := a[i+1]
		switch int(a[i]) {
		case 1:
			switch int(v) {
			case 1:
				l.cap = cad.CapButt
			case 2:
				l.cap = cad.CapSquare
			case 3, 4: // triangular ends are drawn round
				l.cap = cad.CapRound
			}
		case 2:
			switch int(v) {
			case 1, 2:
				l.join = cad.JoinMiter
			case 3, 4: // triangular joins are drawn round
				l.join = cad.JoinRound
			case 5, 6:
				l.join = cad.JoinBevel
			}
		case 3:
			if valid(v) {
				l.miter = math.Max(1, v)
			}
		}
	}
}

// patternLength returns the length of the line type's pattern in drawing
// units.
func (g *gl) patternLength(lt lineType) float64 {
	if lt.absolut {
		return lt.length * 40 * g.scaleOf()
	}
	return lt.length / 100 * g.diagonal() * g.scaleOf()
}

// pen returns the pen lines are drawn with now, without the dashes of an
// adaptive line type (strokePath adapts them to each vector).
func (g *gl) pen() cad.Pen {
	p := &g.pens
	if p.valid {
		return p.cached
	}
	n := p.current()
	pen := cad.Pen{Color: g.penColorOf(n), Cap: g.line.cap, Join: g.line.join, Miter: g.line.miter}
	if w := g.widthOf(n); w > 0 {
		pen.WorldWidth = w
	}
	pen.Color = g.fill.screen(pen.Color, g.fill.transparent)
	if elems := g.dashElements(); elems != nil {
		pen.Dash, pen.DashOffset = cad.DashPattern(elems)
	}
	p.cached, p.valid = pen, true
	return pen
}

// dashElements returns the pattern of a fixed line type in drawing units,
// as CAD line types give them (a dash positive, a gap negative, a dot 0),
// or nil for a solid line.
func (g *gl) dashElements() []float64 {
	lt := g.line.lt
	if !lt.set || lt.typ <= 0 {
		return nil
	}
	L := g.patternLength(lt)
	pat := g.line.pattern(lt.typ)
	if pat == nil || L <= 0 {
		return nil
	}
	elems := make([]float64, len(pat))
	for i, v := range pat {
		v = v / 100 * L
		if i%2 == 1 {
			v = -v
		}
		elems[i] = v
	}
	return elems
}

// strokeState is the path being drawn with one pen: consecutive vectors
// go into one path, drawn when the pen changes or something else is drawn.
type strokeState struct {
	cur        *cad.Path
	pen        cad.Pen
	pendingDot bool // the pen went down and has not moved
}

// strokePath draws a path with the pen: added to the path being drawn, or
// by itself with line type 0 (dots at the vertices) and adaptive line
// types.
func (g *gl) strokePath(p *cad.Path, vertices []cad.Point) {
	lt := g.line.lt
	switch {
	case lt.set && lt.typ == 0:
		for _, v := range vertices {
			g.dot(v)
		}
		return
	case lt.set && lt.typ < 0:
		g.flush()
		g.adaptive(p, lt)
		return
	}
	pen := g.pen()
	s := &g.stroke
	if s.cur != nil && !samePen(s.pen, pen) {
		g.flush()
	}
	if s.cur == nil {
		s.cur, s.pen = &cad.Path{}, pen
	}
	if s.cur.Open() && s.cur.Current() == p.First() {
		s.cur.Continue(p)
	} else {
		s.cur.Append(p)
	}
	if s.cur.Size() > 20000 {
		g.flush()
	}
}

func samePen(a, b cad.Pen) bool {
	return a.Color == b.Color && a.Width == b.Width && a.WorldWidth == b.WorldWidth && a.Cap == b.Cap && a.Join == b.Join &&
		a.Miter == b.Miter && a.DashOffset == b.DashOffset && slices.Equal(a.Dash, b.Dash)
}

// adaptive draws a path with an adaptive line type: each piece (a line or
// a curve) holds a whole number of patterns, starting and ending in the
// middle of a pen-down stretch.
func (g *gl) adaptive(p *cad.Path, lt lineType) {
	pen := g.pen()
	pat := g.line.pattern(lt.typ)
	L := g.patternLength(lt)
	pts := cad.Flatten(p, 16)
	if len(pts) < 2 || pat == nil || L <= 0 {
		g.drawStroke(p, pen)
		return
	}
	length := 0.0
	for i := 1; i < len(pts); i++ {
		length += pts[i].Sub(pts[i-1]).Len()
	}
	if length == 0 {
		return
	}
	n := math.Max(1, math.Round(length/L))
	k := length / n / 100
	// a pattern ending with a pen-up stretch loses half its first dash to
	// its end; one ending with a pen-down stretch shares both out
	var dash []float64
	first := pat[0]
	switch {
	case len(pat) == 1:
		// a pen-down stretch only
		g.drawStroke(p, pen)
		return
	case len(pat)%2 == 0:
		for _, v := range pat {
			dash = append(dash, v*k)
		}
	default:
		first = pat[0] + pat[len(pat)-1]
		dash = append(dash, first*k)
		for _, v := range pat[1 : len(pat)-1] {
			dash = append(dash, v*k)
		}
	}
	pen.Dash = dash
	pen.DashOffset = first * k / 2
	g.drawStroke(p, pen)
}

// flush draws the path being drawn.
func (g *gl) flush() {
	s := &g.stroke
	if s.cur != nil && !s.cur.Empty() {
		g.drawStroke(s.cur, s.pen)
	}
	s.cur = nil
}

func (g *gl) drawStroke(p *cad.Path, pen cad.Pen) {
	if !g.c.budget() {
		return
	}
	g.c.cur.d.Stroke(p, pen)
	g.c.cur.marked = true
}

// dot draws a dot: a round pen touching the paper.
func (g *gl) dot(q cad.Point) {
	g.flush()
	pen := g.pen()
	pen.Dash, pen.Cap = nil, cad.CapRound
	d := 0.01 * g.scaleOf()
	g.drawStroke((&cad.Path{}).MoveTo(q.X, q.Y).LineTo(q.X+d, q.Y), pen)
}
