package music

import (
	"math"
	"slices"
	"sort"
	"strconv"
)

// evLayout is an event placed in a system.
type evLayout struct {
	ev     *Event
	sm     *staffMeasure
	mi     int
	x      float64 // left of the noteheads on the normal side of the stem
	scale  float64
	heads  []headLayout
	up     bool
	stem   bool
	stemX  float64 // left edge of the stem
	stemY0 float64 // where the stem leaves the noteheads
	stemY1 float64 // the stem's tip (or the beam)
	top    float64 // the top of the noteheads
	bottom float64 // the bottom of the noteheads
	right  float64 // the right edge of the noteheads
	beam   *beamGroup
	flags  int
	rest   bool // a rest: x, right, top and bottom are its box
}

// headLayout is one notehead of a chord.
type headLayout struct {
	note *Note
	pos  int
	x    float64
	g    *glyph
}

// beamGroup is events joined by beams.
type beamGroup struct {
	evs []*evLayout
	up  bool
}

// staffCtx lays out one staff of one system.
type staffCtx struct {
	e    *engraver
	sys  *system
	si   int
	c    *drawing
	evs  []*evLayout // in time order
	def  engravingDefaults
	part *Part
}

// numberOfBeams is the number of beams (or flags) of a note value.
func numberOfBeams(t NoteType) int {
	if t < Eighth {
		return 0
	}
	return int(t - Quarter)
}

// layoutStaff draws one staff of a system.
func (e *engraver) layoutStaff(sys *system, si int) *staffLayout {
	ctx := &staffCtx{e: e, sys: sys, si: si, c: newDrawing(), def: e.def, part: e.s.Parts[e.staves[si].part]}
	c := ctx.c
	// the five lines
	for k := 0; k < 5; k++ {
		c.rect(0, float64(k)-e.def.staffLine/2, sys.width, e.def.staffLine)
	}
	ctx.header()
	for mi := sys.first; mi <= sys.last; mi++ {
		ctx.measure(mi)
	}
	ctx.stemsAndBeams()
	ctx.dotsAndMarks()
	ctx.ties()
	ctx.slurs()
	ctx.tuplets()
	ctx.directions()
	ctx.lyrics()
	if si == 0 {
		ctx.systemMarks()
	}
	l := &staffLayout{prims: c.prims}
	l.above = min(c.top.at(-1e9, 1e9, 0), 0)
	l.below = max(c.bot.at(-1e9, 1e9, 4), 4)
	return l
}

// header draws the clef, key signature and time signature at the start of
// the system.
func (ctx *staffCtx) header() {
	e, c := ctx.e, ctx.c
	sm := e.sm[ctx.sys.first][ctx.si]
	x := 0.8
	if g := clefGlyph(sm.clef, false); g != nil {
		pos, _ := clefRef(sm.clef)
		if sm.clef.Sign == ClefPercussion || sm.clef.Sign == ClefTab {
			pos = 4
		}
		c.glyph(g, x, float64(pos)/2, 1)
		x += g.advance + 0.8
	}
	x = ctx.keySig(x, sm.clef, KeySig{}, sm.key)
	if t := e.timeAt(ctx.sys.first); ctx.sys.first == 0 || e.s.Measures[ctx.sys.first].Time != nil {
		ctx.timeSig(x, t)
	}
}

// keySig draws a key signature (after the naturals that cancel prev) and
// returns where it ends.
func (ctx *staffCtx) keySig(x float64, clef Clef, prev, k KeySig) float64 {
	c := ctx.c
	if clef.Sign == ClefPercussion || clef.Sign == ClefTab {
		return x
	}
	if n := cancelled(prev, k); n > 0 {
		nat := sym("accidentalNatural")
		for _, p := range keyPositions(clef, prev) {
			c.glyph(nat, x, float64(p)/2, 1)
			x += nat.advance + 0.1
		}
		x += 0.4
	}
	g := sym("accidentalSharp")
	if k.Fifths < 0 {
		g = sym("accidentalFlat")
	}
	ps := keyPositions(clef, k)
	for _, p := range ps {
		c.glyph(g, x, float64(p)/2, 1)
		x += g.advance + 0.1
	}
	if len(ps) > 0 {
		x += 0.6
	}
	return x
}

// timeSig draws a time signature.
func (ctx *staffCtx) timeSig(x float64, t *TimeSig) {
	c := ctx.c
	if t.Symbol == "common" || t.Symbol == "cut" {
		name := "timeSigCommon"
		if t.Symbol == "cut" {
			name = "timeSigCutCommon"
		}
		c.glyph(sym(name), x, 2, 1)
		return
	}
	num, den := strconv.Itoa(t.Beats), strconv.Itoa(t.BeatType)
	width := func(s string) float64 {
		w := 0.0
		for _, r := range s {
			w += sym("timeSig" + string(r)).advance
		}
		return w
	}
	wn, wd := width(num), width(den)
	w := max(wn, wd)
	for i, s := range []string{num, den} {
		xx := x + (w-[]float64{wn, wd}[i])/2
		for _, r := range s {
			g := sym("timeSig" + string(r))
			c.glyph(g, xx, float64(1+2*i), 1)
			xx += g.advance
		}
	}
}

// measure places the events of one measure.
func (ctx *staffCtx) measure(mi int) {
	e, c := ctx.e, ctx.c
	col := e.cols[mi]
	sm := e.sm[mi][ctx.si]
	mx := col.x
	if mi != ctx.sys.first {
		x := mx + barlineWidth(e.leftBar(mi)) + 0.4
		if sm.clefChange {
			if g := clefGlyph(sm.clef, true); g != nil {
				pos, _ := clefRef(sm.clef)
				if sm.clef.Sign == ClefPercussion || sm.clef.Sign == ClefTab {
					pos = 4
				}
				c.glyph(g, x, float64(pos)/2, 1)
				x += g.advance + 0.6
			}
		}
		if sm.keyChange {
			x = ctx.keySig(x, sm.clef, sm.prevKey, sm.key)
		}
		if t := e.s.Measures[mi].Time; t != nil {
			ctx.timeSig(x, t)
		}
	}
	slotX := map[int]float64{}
	for _, sl := range col.slots {
		slotX[sl.offset] = mx + sl.x
	}
	clef := sm.clef
	changes := sm.clefChanges
	// grace notes lead to the next event
	var graces []*Event
	for _, ev := range sm.events {
		for len(changes) > 0 && changes[0].Offset <= ev.Offset && !ev.Grace {
			ch := changes[0]
			changes = changes[1:]
			clef = ch.Clef
			g := clefGlyph(clef, true)
			if g != nil {
				pos, _ := clefRef(clef)
				gx := slotX[ch.Offset] - ctx.e.cols[mi].slots[0].pre
				if sx, ok := slotX[ch.Offset]; ok {
					gx = sx - g.advance - 0.4 - ctx.preOf(mi, ch.Offset, g.advance)
				}
				c.glyph(g, gx, float64(pos)/2, 1)
			}
		}
		if ev.Hidden {
			continue
		}
		if ev.Grace {
			graces = append(graces, ev)
			continue
		}
		x, ok := slotX[ev.Offset]
		if !ok {
			continue
		}
		if ev.Rest {
			ctx.rest(mi, sm, ev, x)
			graces = nil
			continue
		}
		l := ctx.place(mi, sm, clef, ev, x, 1)
		// grace notes before it, right to left
		gx := x - ctx.accWidth(sm, ev) - 0.3
		for i := len(graces) - 1; i >= 0; i-- {
			gev := graces[i]
			w := noteheadGlyph(gev, nil).advance * graceScale
			gx -= w + ctx.accWidth(sm, gev)*graceScale + 0.3
			ctx.place(mi, sm, clef, gev, gx+ctx.accWidth(sm, gev)*graceScale, graceScale)
		}
		graces = nil
		_ = l
	}
}

// preOf is the room the accidentals of the events at an offset take.
func (ctx *staffCtx) preOf(mi, off int, _ float64) float64 {
	for _, sl := range ctx.e.cols[mi].slots {
		if sl.offset == off {
			return max(sl.pre-ctx.clefChangeWidth(mi, off), 0)
		}
	}
	return 0
}

func (ctx *staffCtx) clefChangeWidth(mi, off int) float64 {
	w := 0.0
	for _, sm := range ctx.e.sm[mi] {
		for _, c := range sm.clefChanges {
			if g := clefGlyph(c.Clef, true); c.Offset == off && g != nil {
				w = max(w, g.advance+0.5)
			}
		}
	}
	return w
}

// accWidth is the room the accidentals of an event take.
func (ctx *staffCtx) accWidth(sm *staffMeasure, ev *Event) float64 {
	pa := placeAccidentals(sm, ev)
	if len(pa) == 0 {
		return 0
	}
	cols := map[int]float64{}
	for _, a := range pa {
		cols[a.col] = max(cols[a.col], accidentalGlyph(a.acc).advance+0.2)
	}
	w := 0.0
	for k := range len(cols) { // in the order of the columns: the sum is the same every time
		w += cols[k]
	}
	return w
}

// rest draws a rest.
func (ctx *staffCtx) rest(mi int, sm *staffMeasure, ev *Event, x float64) {
	e, c := ctx.e, ctx.c
	t := ev.Type
	if ev.MeasureRest {
		t = Whole
		if e.s.Measures[mi].Length >= Breve.Ticks(0) {
			t = Breve
		}
		col := e.cols[mi]
		g := restGlyph(t)
		x0 := col.x + col.bodyStart
		x1 := col.x + col.width - barlineWidth(e.rightBar(mi))
		x = (x0+x1)/2 - g.advance/2
	}
	g := restGlyph(t)
	y := 2.0
	if t == Whole {
		y = 1
	}
	if sm.voices > 1 {
		d := 2.0
		if t <= Half {
			d = 1
		}
		if ev.Voice%2 == 1 {
			y -= d
		} else {
			y += d
		}
	}
	c.glyph(g, x, y, 1)
	ctx.evs = append(ctx.evs, &evLayout{ev: ev, sm: sm, mi: mi, x: x, scale: 1, rest: true,
		right: x + g.advance, top: y + g.y0, bottom: y + g.y1})
	for i := 0; i < ev.Dots; i++ {
		c.glyph(sym("augmentationDot"), x+g.advance+0.3+0.5*float64(i), 1.5, 1)
	}
	if t <= Half && (y < 0 || y > 4) {
		// a whole or half rest outside the staff stands on a leger line
		c.rect(x-0.4, math.Round(y)-e.def.leger/2, g.advance+0.8, e.def.leger)
	}
}

// place lays out the noteheads, accidentals and leger lines of a note or
// chord, and remembers it for the stems and beams.
func (ctx *staffCtx) place(mi int, sm *staffMeasure, clef Clef, ev *Event, x, scale float64) *evLayout {
	e := ctx.e
	e.notes += len(ev.Notes)
	l := &evLayout{ev: ev, sm: sm, mi: mi, x: x, scale: scale}
	for _, n := range ev.Notes {
		l.heads = append(l.heads, headLayout{note: n, pos: staffPos(clef, n.Pitch), x: x, g: noteheadGlyph(ev, n)})
	}
	sort.SliceStable(l.heads, func(i, j int) bool { return l.heads[i].pos > l.heads[j].pos }) // bottom first
	l.up = ctx.stemUp(sm, ev, l.heads)
	l.stem = ev.Type >= Half && ev.Stem != StemNone
	ctx.evs = append(ctx.evs, l)
	return l
}

// columnLefts returns the left edges of the columns of accidentals of the
// widths w (column 0 first), leftwards from x.
func columnLefts(x float64, w map[int]float64) []float64 {
	left := make([]float64, len(w))
	for k := range left {
		x -= w[k]
		left[k] = x
	}
	return left
}

// stemUp decides the direction of the stem of an event alone: by the voice
// when the staff has several, else away from the middle line.
func (ctx *staffCtx) stemUp(sm *staffMeasure, ev *Event, heads []headLayout) bool {
	switch ev.Stem {
	case StemUp:
		return true
	case StemDown:
		return false
	}
	if sm.voices > 1 {
		return ev.Voice%2 == 1
	}
	if len(heads) == 0 {
		return true
	}
	lo, hi := heads[0].pos, heads[len(heads)-1].pos // lowest note (largest pos), highest
	return lo-4 > 4-hi
}

// arrange puts the noteheads of chords on the sides of the stem (a second
// apart they cannot stand in one column) and draws them with their
// accidentals and leger lines.
func (ctx *staffCtx) arrange(l *evLayout) {
	e, c := ctx.e, ctx.c
	s := l.scale
	if len(l.heads) == 0 {
		return
	}
	w := l.heads[0].g.advance * s
	// from the end the stem leaves: bottom up for an up stem
	order := make([]int, len(l.heads))
	for i := range order {
		order[i] = i
	}
	if !l.up {
		slices.Reverse(order)
	}
	prevDisplaced := false
	prevPos := math.MinInt
	for _, i := range order {
		h := &l.heads[i]
		// the second of two notes a second apart goes to the other side
		displaced := prevPos != math.MinInt && abs(h.pos-prevPos) == 1 && !prevDisplaced
		if displaced {
			if l.up {
				h.x = l.x + w - e.def.stem*s
			} else {
				h.x = l.x - w + e.def.stem*s
			}
		}
		prevDisplaced = displaced
		prevPos = h.pos
	}
	l.top, l.bottom = 1e9, -1e9
	l.right = -1e9
	minX := 1e9
	for _, h := range l.heads {
		y := float64(h.pos) / 2
		gw := h.g.advance * s
		c.glyph(h.g, h.x, y, s)
		l.top = min(l.top, y-0.5*s)
		l.bottom = max(l.bottom, y+0.5*s)
		l.right = max(l.right, h.x+gw)
		minX = min(minX, h.x)
	}
	// leger lines
	ext := e.def.legerExt * s
	hi, lo := l.heads[len(l.heads)-1].pos, l.heads[0].pos
	for p := -2; p >= hi; p -= 2 {
		x0, x1 := 1e9, -1e9
		for _, h := range l.heads {
			if h.pos <= p+1 {
				x0, x1 = min(x0, h.x), max(x1, h.x+h.g.advance*s)
			}
		}
		c.rect(x0-ext, float64(p)/2-e.def.leger/2, x1-x0+2*ext, e.def.leger)
	}
	for p := 10; p <= lo; p += 2 {
		x0, x1 := 1e9, -1e9
		for _, h := range l.heads {
			if h.pos >= p-1 {
				x0, x1 = min(x0, h.x), max(x1, h.x+h.g.advance*s)
			}
		}
		c.rect(x0-ext, float64(p)/2-e.def.leger/2, x1-x0+2*ext, e.def.leger)
	}
	// accidentals, in columns leftwards from the noteheads
	pa := placeAccidentals(l.sm, l.ev)
	colW := map[int]float64{}
	for _, a := range pa {
		colW[a.col] = max(colW[a.col], (accidentalGlyph(a.acc).advance+0.2)*s)
	}
	left := columnLefts(minX, colW)
	for _, a := range pa {
		g := accidentalGlyph(a.acc)
		x := left[a.col]
		y := float64(a.pos) / 2
		if a.note.Cautionary {
			pl, pr := sym("accidentalParensLeft"), sym("accidentalParensRight")
			c.glyph(pl, x-pl.advance*s, y, s)
			c.glyph(pr, x+g.advance*s, y, s)
		}
		c.glyph(g, x, y, s)
	}
}
