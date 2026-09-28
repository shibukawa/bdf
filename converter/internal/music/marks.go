package music

import (
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
)

// xAt returns the x of an offset of a measure: its slot, or between the
// slots around it.
func (ctx *staffCtx) xAt(mi, off int) float64 {
	col := ctx.e.cols[mi]
	sl := col.slots
	if len(sl) == 0 || off <= sl[0].offset {
		if len(sl) == 0 {
			return col.x + col.bodyStart + col.lead
		}
		return col.x + sl[0].x
	}
	for i, s := range sl {
		if s.offset == off {
			return col.x + s.x
		}
		if s.offset > off {
			p := sl[i-1]
			f := float64(off-p.offset) / float64(s.offset-p.offset)
			return col.x + p.x + f*(s.x-p.x)
		}
	}
	last := sl[len(sl)-1]
	end := col.width - barlineWidth(ctx.e.rightBar(mi))
	f := float64(off-last.offset) / float64(max(ctx.e.s.Measures[mi].Length-last.offset, 1))
	return col.x + last.x + min(f, 1)*(end-last.x)
}

// dotsAndMarks draws augmentation dots and articulations.
func (ctx *staffCtx) dotsAndMarks() {
	c := ctx.c
	for _, l := range ctx.evs {
		if l.rest || len(l.heads) == 0 {
			continue
		}
		s := l.scale
		ev := l.ev
		if ev.Dots > 0 {
			x := l.right + 0.3*s
			if l.up && l.flags > 0 {
				x = max(x, l.stemX+0.9*s)
			}
			seen := map[int]bool{}
			for _, h := range l.heads {
				p := h.pos
				if p%2 == 0 {
					if l.up || ev.Voice%2 == 1 {
						p--
					} else {
						p++
					}
				}
				if seen[p] {
					continue
				}
				seen[p] = true
				for i := 0; i < ev.Dots; i++ {
					c.glyph(sym("augmentationDot"), x+0.5*s*float64(i), float64(p)/2, s)
				}
			}
		}
		ctx.articulations(l)
	}
}

// articulations draws the articulations, fermatas and ornaments of an
// event: staccato, tenuto, accents on the side of the noteheads away from
// the stem; fermatas and ornaments above the staff.
func (ctx *staffCtx) articulations(l *evLayout) {
	c := ctx.c
	ev := l.ev
	if len(ev.Artics) == 0 {
		return
	}
	s := l.scale
	multi := l.sm.voices > 1
	xc := l.x + l.heads[0].g.advance*s/2
	upY := l.top - 0.4*s
	if l.stem && l.up {
		upY = min(upY, l.stemY1-0.4*s)
	}
	downY := l.bottom + 0.4*s
	if l.stem && !l.up {
		downY = max(downY, l.stemY1+0.4*s)
	}
	for _, a := range ev.Artics {
		name := a
		above := true
		outside := false
		switch {
		case a == "fermata":
			above = !(multi && ev.Voice%2 == 0)
			outside = true
		case strings.HasPrefix(a, "ornament"):
			outside = true
		case strings.HasPrefix(a, "artic"):
			if multi {
				above = ev.Voice%2 == 1
			} else if l.stem {
				above = !l.up
			} else {
				above = true
			}
			outside = a == "articAccent" || a == "articMarcato"
		default:
			continue
		}
		if !strings.HasPrefix(a, "ornament") {
			if above {
				name += "Above"
			} else {
				name += "Below"
			}
		}
		g := sym(name)
		if g == nil {
			continue
		}
		gx := xc - (g.x0+g.x1)*s/2
		if above {
			y := upY
			if outside {
				y = min(y, -0.6)
			}
			oy := y - g.y1*s
			if !outside && oy+(g.y0+g.y1)*s/2 > -0.2 && isLine((oy+(g.y0+g.y1)*s/2)*2) {
				oy -= 0.3
			}
			c.glyph(g, gx, oy, s)
			upY = oy + g.y0*s - 0.25
		} else {
			y := downY
			if outside {
				y = max(y, 4.6)
			}
			oy := y - g.y0*s
			if !outside && oy+(g.y0+g.y1)*s/2 < 4.2 && isLine((oy+(g.y0+g.y1)*s/2)*2) {
				oy += 0.3
			}
			c.glyph(g, gx, oy, s)
			downY = oy + g.y1*s + 0.25
		}
	}
}

// isLine reports whether a staff position (half spaces) is on a line of
// the staff.
func isLine(pos float64) bool {
	p := math.Round(pos)
	return math.Abs(pos-p) < 0.3 && int(p)%2 == 0 && p >= 0 && p <= 8
}

// tieSide returns -1 to tie a note above, 1 below.
func tieSide(l *evLayout, i int) float64 {
	if l.sm.voices > 1 {
		if l.ev.Voice%2 == 1 {
			return -1
		}
		return 1
	}
	n := len(l.heads)
	if n > 1 {
		// heads are bottom first: the lower half ties below
		if 2*i+1 < n {
			return 1
		}
		if 2*i+1 > n {
			return -1
		}
	}
	if l.up {
		return 1
	}
	return -1
}

// ties draws the ties of the staff, and the halves of ties that cross to
// the systems before and after.
func (ctx *staffCtx) ties() {
	c, def := ctx.c, ctx.def
	ended := map[*Note]bool{}
	sameNote := func(a, b *Note) bool {
		return a.Pitch == b.Pitch
	}
	for i, l := range ctx.evs {
		if l.rest {
			continue
		}
		for hi, h := range l.heads {
			if !h.note.TieStart {
				continue
			}
			side := tieSide(l, hi)
			s := l.scale
			x0 := h.x + h.g.advance*s + 0.12
			y0 := float64(h.pos)/2 + side*0.5
			found := false
			for _, m := range ctx.evs[i+1:] {
				if m.rest || m.ev.Grace && !l.ev.Grace {
					continue
				}
				for _, h2 := range m.heads {
					if h2.note.TieStop && sameNote(h.note, h2.note) && !ended[h2.note] {
						ended[h2.note] = true
						x1 := h2.x - 0.12
						if x1-x0 < 0.6 {
							x1 = x0 + 0.6
						}
						c.fill(curve(x0, y0, x1, y0, side*min(0.2+0.06*(x1-x0), 0.9), def.tieEnd, def.tieMid))
						found = true
						break
					}
				}
				if found || m.ev.Voice == l.ev.Voice && m.ev.Offset > l.ev.Offset && m.mi > l.mi+1 {
					break
				}
			}
			if !found && l.mi == ctx.sys.last {
				x1 := ctx.sys.width - 0.3
				if x1 > x0+0.3 {
					c.fill(curve(x0, y0, x1, y0, side*min(0.2+0.06*(x1-x0), 0.9), def.tieEnd, def.tieMid))
				}
			}
		}
	}
	// ties from the system before
	for _, l := range ctx.evs {
		if l.rest || l.mi != ctx.sys.first {
			continue
		}
		for hi, h := range l.heads {
			if h.note.TieStop && !ended[h.note] && ctx.sys.first > 0 {
				x0 := ctx.sys.header - 0.5
				x1 := h.x - 0.12
				side := tieSide(l, hi)
				y := float64(h.pos)/2 + side*0.5
				if x1 > x0+0.3 {
					c.fill(curve(x0, y, x1, y, side*min(0.2+0.06*(x1-x0), 0.9), def.tieEnd, def.tieMid))
				}
				ended[h.note] = true
			}
		}
	}
}

// slurEnd is the point a slur starts or ends at on an event.
func slurEnd(l *evLayout, above bool) (x, y float64) {
	s := l.scale
	if l.rest || len(l.heads) == 0 {
		x = (l.x + l.right) / 2
		if above {
			return x, l.top - 0.5
		}
		return x, l.bottom + 0.5
	}
	x = l.x + l.heads[0].g.advance*s/2
	if above {
		y = l.top - 0.6*s
		if l.stem && l.up {
			x = l.stemX
			y = min(y, l.stemY1-0.4*s)
		}
	} else {
		y = l.bottom + 0.6*s
		if l.stem && !l.up {
			x = l.stemX
			y = max(y, l.stemY1+0.4*s)
		}
	}
	return x, y
}

// slurs draws the slurs of the staff, open at the ends of the system when
// they cross to the systems before or after.
func (ctx *staffCtx) slurs() {
	type open struct {
		l     *evLayout
		above *bool
		i     int
	}
	starts := map[int]open{}
	for i, l := range ctx.evs {
		for _, m := range l.ev.Slurs {
			if m.Start {
				continue
			}
			if st, ok := starts[m.Number]; ok {
				ctx.slur(st.l, l, st.i, i, st.above)
				delete(starts, m.Number)
			} else {
				ctx.slur(nil, l, -1, i, m.Above)
			}
		}
		for _, m := range l.ev.Slurs {
			if m.Start {
				starts[m.Number] = open{l, m.Above, i}
			}
		}
	}
	// the slurs that go on in the next system, in the order of their numbers
	for _, n := range slices.Sorted(maps.Keys(starts)) {
		st := starts[n]
		ctx.slur(st.l, nil, st.i, len(ctx.evs), st.above)
	}
}

// slur draws a slur from a to b (nil: the edge of the system), bulging
// enough to clear the events between.
func (ctx *staffCtx) slur(a, b *evLayout, ia, ib int, above *bool) {
	ref := a
	if ref == nil {
		ref = b
	}
	if ref == nil {
		return
	}
	up := true
	switch {
	case above != nil:
		up = *above
	case ref.sm.voices > 1:
		up = ref.ev.Voice%2 == 1
	default:
		// on the noteheads' side when every stem points the same way
		allUp := true
		for k := max(ia, 0); k < min(ib+1, len(ctx.evs)); k++ {
			l := ctx.evs[k]
			if l.rest || !l.stem || l.ev.Voice != ref.ev.Voice {
				continue
			}
			if !l.up {
				allUp = false
			}
		}
		up = !allUp
	}
	var x0, y0, x1, y1 float64
	if a != nil {
		x0, y0 = slurEnd(a, up)
	}
	if b != nil {
		x1, y1 = slurEnd(b, up)
	}
	if a == nil {
		x0, y0 = ctx.sys.header-0.5, y1
	}
	if b == nil {
		x1, y1 = ctx.sys.width-0.3, y0
	}
	if x1-x0 < 1 {
		x1 = x0 + 1
	}
	sign := 1.0
	if up {
		sign = -1
	}
	h := min(0.5+0.1*(x1-x0), 2.5)
	for k := max(ia+1, 0); k < min(ib, len(ctx.evs)); k++ {
		l := ctx.evs[k]
		if l.ev.Voice != ref.ev.Voice {
			continue
		}
		xm, ym := slurEnd(l, up)
		t := (xm - x0) / (x1 - x0)
		if t <= 0.05 || t >= 0.95 {
			continue
		}
		line := y0 + (y1-y0)*t
		need := (line - ym) * -sign // how far the curve must reach past the line
		if need > 0 {
			h = max(h, need/(4*t*(1-t))+0.2)
		}
	}
	h = min(h, 6)
	ctx.c.fill(curve(x0, y0, x1, y1, sign*h, ctx.def.slurEnd, ctx.def.slurMid))
}

// tuplets draws the numbers and brackets of tuplets.
func (ctx *staffCtx) tuplets() {
	c := ctx.c
	groups := map[*Tuplet][]*evLayout{}
	var order []*Tuplet
	for _, l := range ctx.evs {
		if t := l.ev.Tuplet; t != nil && !l.ev.Grace {
			if _, ok := groups[t]; !ok {
				order = append(order, t)
			}
			groups[t] = append(groups[t], l)
		}
	}
	for _, t := range order {
		evs := groups[t]
		if !t.ShowNumber && (t.Bracket == nil || !*t.Bracket) {
			continue
		}
		// the stem side when the stems agree, else above
		up, down := 0, 0
		var beam *beamGroup
		oneBeam := true
		for i, l := range evs {
			if l.rest {
				oneBeam = false
				continue
			}
			if l.stem {
				if l.up {
					up++
				} else {
					down++
				}
			}
			if i == 0 || l.beam != beam {
				if i > 0 {
					oneBeam = false
				}
				beam = l.beam
			}
		}
		if beam == nil {
			oneBeam = false
		}
		above := down == 0 || up > 0 && up >= down
		if evs[0].sm.voices > 1 {
			above = evs[0].ev.Voice%2 == 1
		}
		bracket := !oneBeam
		if t.Bracket != nil {
			bracket = *t.Bracket
		}
		x0 := evs[0].x
		last := evs[len(evs)-1]
		x1 := last.right
		if !last.rest && last.stem && last.up {
			x1 = max(x1, last.stemX+ctx.def.stem)
		}
		y := 0.0
		if above {
			y = 0.0
			for _, l := range evs {
				y = min(y, l.top)
				if l.stem && l.up {
					y = min(y, l.stemY1-ctx.def.beam)
				}
			}
			y -= 1.0
		} else {
			y = 4.0
			for _, l := range evs {
				y = max(y, l.bottom)
				if l.stem && !l.up {
					y = max(y, l.stemY1+ctx.def.beam)
				}
			}
			y += 1.0
		}
		digits := strconv.Itoa(t.Actual)
		nw := 0.0
		for _, d := range digits {
			nw += sym("tuplet" + string(d)).advance
		}
		xc := (x0 + x1) / 2
		if t.ShowNumber {
			nx := xc - nw/2
			for _, d := range digits {
				g := sym("tuplet" + string(d))
				c.glyph(g, nx, y+0.7, 1)
				nx += g.advance
			}
		}
		if bracket {
			hook := 0.6
			if !above {
				hook = -0.6
			}
			w := ctx.def.tupletBracket
			gapL, gapR := xc, xc
			if t.ShowNumber {
				gapL, gapR = xc-nw/2-0.3, xc+nw/2+0.3
			}
			c.line(x0, y+hook, x0, y, w, 0)
			c.line(x0, y, gapL, y, w, 0)
			c.line(gapR, y, x1, y, w, 0)
			c.line(x1, y, x1, y+hook, w, 0)
		}
	}
}

// directions draws dynamics, hairpins, words, chord symbols, pedal marks,
// octave lines, segni and codas.
func (ctx *staffCtx) directions() {
	e, c := ctx.e, ctx.c
	sp := e.sp
	type pending struct {
		x    float64
		kind DirectionKind
		size int
	}
	hairpins := map[int]pending{}
	octaves := map[int]pending{}
	pedal := -1.0
	isAbove := func(d *Direction, def bool) bool {
		if d.Above != nil {
			return *d.Above
		}
		return def
	}
	drawHairpin := func(p pending, x1 float64) {
		x0 := p.x
		if x1-x0 < 1 {
			x1 = x0 + 1
		}
		y := c.bot.at(x0, x1, 4) + 1.2
		o := 0.6
		w := ctx.def.hairpin
		if p.kind == DirCrescendo {
			c.line(x0, y, x1, y-o, w, 0)
			c.line(x0, y, x1, y+o, w, 0)
		} else {
			c.line(x0, y-o, x1, y, w, 0)
			c.line(x0, y+o, x1, y, w, 0)
		}
	}
	drawOctave := func(p pending, x1 float64) {
		name, below := "ottavaAlta", false
		switch {
		case p.size >= 15:
			name = "quindicesimaAlta"
		case p.size < 0:
			name, below = "ottavaBassaVb", true
		}
		g := sym(name)
		x0 := p.x
		if below {
			y := c.bot.at(x0, x1, 4) + 0.6 - g.y0
			c.glyph(g, x0, y, 1)
			c.line(x0+g.advance+0.3, y-0.4, x1, y-0.4, ctx.def.octaveLine, 0.6)
			c.line(x1, y-0.4, x1, y-1.2, ctx.def.octaveLine, 0)
		} else {
			y := c.top.at(x0, x1, 0) - 0.6 - g.y1
			c.glyph(g, x0, y, 1)
			c.line(x0+g.advance+0.3, y-0.6, x1, y-0.6, ctx.def.octaveLine, 0.6)
			c.line(x1, y-0.6, x1, y+0.4, ctx.def.octaveLine, 0)
		}
	}
	for mi := ctx.sys.first; mi <= ctx.sys.last; mi++ {
		for _, d := range e.sm[mi][ctx.si].directions {
			x := ctx.xAt(mi, d.Offset)
			switch d.Kind {
			case DirDynamic:
				if isAbove(d, false) {
					ctx.dynamic(d.Text, x, true)
				} else {
					ctx.dynamic(d.Text, x, false)
				}
			case DirCrescendo, DirDiminuendo:
				if d.Stop {
					if p, ok := hairpins[d.Number]; ok {
						drawHairpin(p, x)
						delete(hairpins, d.Number)
					} else {
						drawHairpin(pending{x: ctx.sys.header, kind: d.Kind}, x)
					}
					continue
				}
				hairpins[d.Number] = pending{x: x + 0.3, kind: d.Kind}
			case DirWords, DirChord:
				if strings.TrimSpace(d.Text) == "" {
					continue
				}
				t := e.text.line(d.Text, e.wordSize(), false, d.Italic)
				w := t.width / sp
				c.markText(bdf.MarkBox, "")
				if isAbove(d, true) {
					y := c.top.at(x, x+w, 0) - 0.4 - t.desc/sp
					c.text(t, x, y, 0, sp)
				} else {
					y := c.bot.at(x, x+w, 4) + 0.4 + t.asc/sp
					c.text(t, x, y, 0, sp)
				}
			case DirPedal:
				if d.Stop {
					g := sym("keyboardPedalUp")
					y := c.bot.at(x, x+g.advance, 4) + 0.6 - g.y0
					c.glyph(g, x, y, 1)
					pedal = -1
				} else {
					g := sym("keyboardPedalPed")
					y := c.bot.at(x, x+g.advance, 4) + 0.6 - g.y0
					c.glyph(g, x, y, 1)
					pedal = x
				}
			case DirOctave:
				if d.Stop {
					if p, ok := octaves[d.Number]; ok {
						drawOctave(p, x)
						delete(octaves, d.Number)
					}
					continue
				}
				octaves[d.Number] = pending{x: x, size: d.Size}
			case DirSegno, DirCoda:
				name := "segno"
				if d.Kind == DirCoda {
					name = "coda"
				}
				g := sym(name)
				y := c.top.at(x, x+g.advance, 0) - 0.5 - g.y1
				c.glyph(g, x, y, 1)
			}
		}
	}
	// those that go on in the next system, in the order of their numbers
	for _, n := range slices.Sorted(maps.Keys(hairpins)) {
		drawHairpin(hairpins[n], ctx.sys.width-0.5)
	}
	for _, n := range slices.Sorted(maps.Keys(octaves)) {
		drawOctave(octaves[n], ctx.sys.width-0.5)
	}
	_ = pedal
}

// dynamicGlyphs are the SMuFL letters of dynamics.
var dynamicGlyphs = map[rune]string{'p': "dynamicPiano", 'm': "dynamicMezzo", 'f': "dynamicForte",
	'r': "dynamicRinforzando", 's': "dynamicSforzando", 'z': "dynamicZ", 'n': "dynamicNiente"}

// dynamic draws a dynamic mark: in the SMuFL letters when it is made of
// them, else in italics.
func (ctx *staffCtx) dynamic(text string, x float64, above bool) {
	c := ctx.c
	var gs []*glyph
	for _, r := range text {
		name, ok := dynamicGlyphs[r]
		if !ok {
			gs = nil
			break
		}
		gs = append(gs, sym(name))
	}
	if gs == nil {
		t := ctx.e.text.line(text, ctx.e.wordSize(), false, true)
		w := t.width / ctx.e.sp
		if above {
			c.text(t, x, c.top.at(x, x+w, 0)-0.4-t.desc/ctx.e.sp, 0, ctx.e.sp)
		} else {
			c.text(t, x, c.bot.at(x, x+w, 4)+0.4+t.asc/ctx.e.sp, 0, ctx.e.sp)
		}
		return
	}
	w := 0.0
	for _, g := range gs {
		w += g.advance
	}
	y := 0.0
	if above {
		y = c.top.at(x, x+w, 0) - 1.0
	} else {
		y = c.bot.at(x, x+w, 4) + 2.2
	}
	for _, g := range gs {
		c.glyph(g, x, y, 1)
		x += g.advance
	}
}

// lyric is a syllable placed on the staff.
type lyric struct {
	l      *evLayout
	syl    *LyricSyllable
	t      *textLine
	x0, x1 float64
}

// lyrics draws the lyrics of each verse on one line below the staff,
// with hyphens between the syllables of a word and extender lines after
// a word's last syllable held over several notes.
func (ctx *staffCtx) lyrics() {
	e, c := ctx.e, ctx.c
	sp := e.sp
	verses := map[int][]*lyric{}
	maxVerse := 0
	for _, l := range ctx.evs {
		if l.rest {
			continue
		}
		for _, syl := range l.ev.Lyrics {
			if syl.Text == "" && !syl.Extend {
				continue
			}
			v := max(syl.Verse, 1)
			t := e.text.line(syl.Text, e.lyricSize(), false, false)
			xc := l.x + l.heads[0].g.advance*l.scale/2
			w := t.width / sp
			verses[v] = append(verses[v], &lyric{l: l, syl: syl, t: t, x0: xc - w/2, x1: xc + w/2})
			maxVerse = max(maxVerse, v)
		}
	}
	if maxVerse == 0 {
		return
	}
	size := e.lyricSize() / sp
	base := c.bot.at(-1e9, 1e9, 4) + 0.5 + 0.85*size
	for v := 1; v <= maxVerse; v++ {
		ls := verses[v]
		y := base + float64(v-1)*1.3*size
		c.markText(bdf.MarkParagraph, "")
		for i, ly := range ls {
			if i > 0 {
				prev := ls[i-1]
				if prev.syl.Syllabic == SyllableBegin || prev.syl.Syllabic == SyllableMiddle {
					c.markText(bdf.MarkWrap, "")
				} else {
					c.markText(bdf.MarkLine, "")
				}
			}
			c.text(ly.t, (ly.x0+ly.x1)/2, y, 1, sp)
		}
		hy := y - 0.3*size
		for i, ly := range ls {
			if ly.syl.Syllabic == SyllableBegin || ly.syl.Syllabic == SyllableMiddle {
				x0 := ly.x1
				x1 := x0 + 1.4
				if i+1 < len(ls) {
					x1 = ls[i+1].x0
				}
				if gap := x1 - x0; gap > 0.5 {
					w := min(0.9, gap*0.5)
					c.rect((x0+x1)/2-w/2, hy-0.06, w, 0.12)
				}
			}
			if ly.syl.Extend {
				// to the end of the last note before the next syllable
				end := ly.x1
				next := -1
				if i+1 < len(ls) {
					next = indexOf(ctx.evs, ls[i+1].l)
				}
				start := indexOf(ctx.evs, ly.l)
				for k := start + 1; k < len(ctx.evs) && (next < 0 || k < next); k++ {
					m := ctx.evs[k]
					if m.rest || m.ev.Voice != ly.l.ev.Voice {
						continue
					}
					end = max(end, m.right)
				}
				if next < 0 && ly.l.mi == ctx.sys.last {
					end = max(end, ctx.sys.width-0.5)
				}
				if end > ly.x1+0.6 {
					c.rect(ly.x1+0.2, y+0.1, end-ly.x1-0.2, ctx.def.staffLine)
				}
			}
		}
	}
}

func indexOf(evs []*evLayout, l *evLayout) int {
	for i, m := range evs {
		if m == l {
			return i
		}
	}
	return -1
}

// measureNumber returns the number shown for a measure: its own, or the
// one it has by counting from the measures before.
func (e *engraver) measureNumber(mi int) string {
	if n := e.s.Measures[mi].Number; n != "" {
		return n
	}
	return strconv.Itoa(e.counts[mi])
}

// systemMarks draws what belongs above the whole system on its top staff:
// the measure number, tempo marks, rehearsal marks and volta brackets.
func (ctx *staffCtx) systemMarks() {
	e, c := ctx.e, ctx.c
	sp := e.sp
	if ctx.sys.first > 0 {
		t := e.text.line(e.measureNumber(ctx.sys.first), 1.6*sp, false, true)
		y := c.top.at(0, t.width/sp, 0) - 0.6 - t.desc/sp
		c.markText(bdf.MarkBox, "")
		c.text(t, 0, y, 0, sp)
	}
	for mi := ctx.sys.first; mi <= ctx.sys.last; mi++ {
		m := e.s.Measures[mi]
		col := e.cols[mi]
		if m.Rehearsal != "" {
			t := e.text.line(m.Rehearsal, 2.4*sp, true, false)
			x := col.x + 0.3
			if mi == ctx.sys.first {
				x = 0
			}
			w := t.width/sp + 0.8
			h := (t.asc + t.desc) / sp
			y := c.top.at(x, x+w, 0) - 0.8
			c.markText(bdf.MarkBox, "")
			c.text(t, x+0.4, y-0.4-t.desc/sp, 0, sp)
			lw := ctx.def.staffLine * 1.5
			c.line(x, y, x+w, y, lw, 0)
			c.line(x, y-h-0.8, x+w, y-h-0.8, lw, 0)
			c.line(x, y, x, y-h-0.8, lw, 0)
			c.line(x+w, y, x+w, y-h-0.8, lw, 0)
		}
		for _, tm := range m.Tempo {
			x := ctx.xAt(mi, tm.Offset)
			if tm.Offset == 0 {
				x = col.x + col.bodyStart
				if mi == ctx.sys.first {
					x = ctx.sys.header
				}
			}
			ctx.tempoMark(tm, x)
		}
		if m.Ending != nil && m.Ending.Start {
			ctx.volta(mi)
		}
	}
	// a volta that started in the system before
	if f := ctx.sys.first; f > 0 {
		if m := e.s.Measures[f]; m.Ending != nil && !m.Ending.Start {
			ctx.volta(f)
		}
	}
}

// tempoMark draws a tempo indication: its words in bold and the
// metronome mark.
func (ctx *staffCtx) tempoMark(tm TempoMark, x float64) {
	e, c := ctx.e, ctx.c
	sp := e.sp
	size := 2.4 * sp
	var words *textLine
	w := 0.0
	if tm.Text != "" {
		words = e.text.line(tm.Text, size, true, false)
		w = words.width/sp + 0.8
	}
	var met *glyph
	var num *textLine
	if tm.PerMinute > 0 {
		names := map[NoteType]string{Whole: "metNoteWhole", Half: "metNoteHalfUp", Quarter: "metNoteQuarterUp",
			Eighth: "metNote8thUp", N16th: "metNote16thUp"}
		if n, ok := names[tm.Beat]; ok {
			met = sym(n)
		} else {
			met = sym("metNoteQuarterUp")
		}
		num = e.text.line(" = "+strconv.FormatFloat(tm.PerMinute, 'f', -1, 64), size, false, false)
		w += met.advance + float64(tm.Dots)*0.6 + num.width/sp
	}
	y := c.top.at(x, x+w, 0) - 1.0
	c.markText(bdf.MarkBox, "")
	if words != nil {
		c.text(words, x, y, 0, sp)
		x += words.width/sp + 0.8
	}
	if met != nil {
		c.glyph(met, x, y-0.2, 0.8)
		x += met.advance * 0.8
		for i := 0; i < tm.Dots; i++ {
			c.glyph(sym("metAugmentationDot"), x+0.2, y-0.2, 0.8)
			x += 0.6
		}
		c.markText(bdf.MarkAltText, "♩")
		c.text(num, x, y, 0, sp)
	}
}

// volta draws the bracket of an ending from measure mi to its last
// measure in the system.
func (ctx *staffCtx) volta(mi int) {
	e, c := ctx.e, ctx.c
	sp := e.sp
	m := e.s.Measures[mi]
	end := mi
	for end < ctx.sys.last && !e.s.Measures[end].Ending.Stop {
		if next := e.s.Measures[end+1].Ending; next == nil || next.Start {
			break
		}
		end++
	}
	x0 := e.cols[mi].x + 0.2
	if mi == ctx.sys.first {
		x0 = ctx.sys.header
	}
	x1 := e.cols[end].x + e.cols[end].width - 0.3
	var t *textLine
	th := 0.0
	if m.Ending.Start && m.Ending.Text != "" {
		t = e.text.line(m.Ending.Text, 1.9*sp, false, false)
		th = (t.asc + t.desc) / sp
	}
	// the bracket and its number clear what is below
	y := c.top.at(x0, x1, 0) - 0.8 - th
	lw := ctx.def.staffLine * 1.4
	if m.Ending.Start {
		c.line(x0, y+1.6, x0, y, lw, 0)
	}
	c.line(x0, y, x1, y, lw, 0)
	if last := e.s.Measures[end].Ending; last != nil && last.Stop && !last.Open {
		c.line(x1, y, x1, y+1.6, lw, 0)
	}
	if t != nil {
		c.markText(bdf.MarkBox, "")
		c.text(t, x0+0.4, y+0.2+t.asc/sp, 0, sp)
	}
}
