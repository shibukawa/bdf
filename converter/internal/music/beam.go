package music

// stemsAndBeams groups the events into beams, settles the stem
// directions, and draws the noteheads, stems, flags and beams.
func (ctx *staffCtx) stemsAndBeams() {
	e := ctx.e
	type key struct{ mi, voice int }
	lists := map[key][]*evLayout{}
	var order []key
	for _, l := range ctx.evs {
		if l.rest || l.ev.Grace {
			continue
		}
		k := key{l.mi, l.ev.Voice}
		if _, ok := lists[k]; !ok {
			order = append(order, k)
		}
		lists[k] = append(lists[k], l)
	}
	for _, k := range order {
		for _, g := range ctx.beamGroups(k.mi, lists[k]) {
			ctx.beamDirection(g)
			for _, l := range g.evs {
				l.beam = g
				l.up = g.up
			}
		}
	}
	for _, l := range ctx.evs {
		if l.rest {
			continue
		}
		ctx.arrange(l)
		ctx.stemBase(l)
	}
	drawn := map[*beamGroup]bool{}
	for _, l := range ctx.evs {
		if l.rest {
			continue
		}
		if l.beam != nil {
			if !drawn[l.beam] {
				drawn[l.beam] = true
				ctx.drawBeam(l.beam)
			}
			continue
		}
		ctx.drawStem(l)
	}
	_ = e
}

// beamable reports whether an event can join a beam.
func beamable(ev *Event) bool {
	return !ev.Rest && !ev.Grace && len(ev.Notes) > 0 && ev.Type >= Eighth
}

// beamGroups divides the events of one voice of a measure into beams: as
// the score says, or by the beats of the time signature (a tuplet is
// beamed whole; in 4/4 eighths alone are beamed by half measures).
func (ctx *staffCtx) beamGroups(mi int, evs []*evLayout) []*beamGroup {
	var out []*beamGroup
	var cur *beamGroup
	flush := func() {
		if cur != nil && len(cur.evs) > 1 {
			out = append(out, cur)
		}
		cur = nil
	}
	t := ctx.e.timeAt(mi)
	beat, compound := t.beat()
	halves := [2]bool{}
	if !compound && t.Beats == 4 && t.BeatType == 4 {
		halves = [2]bool{true, true}
		for _, l := range evs {
			h := min(l.ev.Offset/(2*PPQ), 1)
			if l.ev.Type != Eighth || l.ev.Dots > 0 || l.ev.Tuplet != nil {
				halves[h] = false
			}
		}
	}
	groupOf := func(ev *Event) any {
		if ev.Tuplet != nil {
			return ev.Tuplet
		}
		if h := min(ev.Offset/(2*PPQ), 1); halves[h] {
			return 100 + h
		}
		return ev.Offset / max(beat, 1)
	}
	var curKey any
	for _, l := range evs {
		ev := l.ev
		if ev.Beams != nil {
			// the score's beams
			if len(ev.Beams) == 0 || !beamable(ev) {
				flush()
				continue
			}
			switch ev.Beams[0] {
			case BeamBegin:
				flush()
				cur = &beamGroup{evs: []*evLayout{l}}
			case BeamContinue:
				if cur == nil {
					cur = &beamGroup{}
				}
				cur.evs = append(cur.evs, l)
			case BeamEnd:
				if cur == nil {
					cur = &beamGroup{}
				}
				cur.evs = append(cur.evs, l)
				flush()
			default:
				flush()
			}
			continue
		}
		if !beamable(ev) {
			flush()
			continue
		}
		k := groupOf(ev)
		if cur == nil || k != curKey || cur.evs[len(cur.evs)-1].ev.Offset+cur.evs[len(cur.evs)-1].ev.Duration != ev.Offset {
			flush()
			cur = &beamGroup{}
			curKey = k
		}
		cur.evs = append(cur.evs, l)
	}
	flush()
	return out
}

// beamDirection settles the stems of a beam: as the score or the voice
// says, else away from the middle line for the note furthest from it.
func (ctx *staffCtx) beamDirection(g *beamGroup) {
	for _, l := range g.evs {
		switch l.ev.Stem {
		case StemUp:
			g.up = true
			return
		case StemDown:
			g.up = false
			return
		}
	}
	if l := g.evs[0]; l.sm.voices > 1 {
		g.up = l.ev.Voice%2 == 1
		return
	}
	hi, lo := 1000, -1000
	for _, l := range g.evs {
		for _, h := range l.heads {
			hi, lo = min(hi, h.pos), max(lo, h.pos)
		}
	}
	g.up = lo-4 > 4-hi
}

// stemBase works out where the stem of an event leaves its noteheads and
// its tip without a beam.
func (ctx *staffCtx) stemBase(l *evLayout) {
	if len(l.heads) == 0 {
		return
	}
	s := l.scale
	g := l.heads[0].g
	t := ctx.def.stem * s
	if l.up {
		ax, ay := g.anchor("stemUpSE")
		l.stemX = l.x + ax*s - t
		l.stemY0 = float64(l.heads[0].pos)/2 + ay*s
		l.stemY1 = float64(l.heads[len(l.heads)-1].pos)/2 - 3.5*s
		if s == 1 {
			l.stemY1 = min(l.stemY1, 2)
		}
	} else {
		ax, ay := g.anchor("stemDownNW")
		l.stemX = l.x + ax*s
		l.stemY0 = float64(l.heads[len(l.heads)-1].pos)/2 + ay*s
		l.stemY1 = float64(l.heads[0].pos)/2 + 3.5*s
		if s == 1 {
			l.stemY1 = max(l.stemY1, 2)
		}
	}
	if l.beam == nil {
		l.flags = numberOfBeams(l.ev.Type)
		if l.flags > 2 {
			ext := float64(l.flags-2) * 0.5 * s
			if l.up {
				l.stemY1 -= ext
			} else {
				l.stemY1 += ext
			}
		}
	}
}

var flagNames = []string{"8th", "16th", "32nd", "64th", "128th"}

// drawStem draws the stem, flag, grace slash and tremolo of an event
// without a beam.
func (ctx *staffCtx) drawStem(l *evLayout) {
	c := ctx.c
	if !l.stem || len(l.heads) == 0 {
		return
	}
	s := l.scale
	t := ctx.def.stem * s
	c.rect(l.stemX, min(l.stemY0, l.stemY1), t, abs64(l.stemY1-l.stemY0))
	if l.flags > 0 {
		name := "flag" + flagNames[min(l.flags, len(flagNames))-1]
		if l.up {
			g := sym(name + "Up")
			ax, ay := g.anchor("stemUpNW")
			c.glyph(g, l.stemX-ax*s, l.stemY1-ay*s, s)
		} else {
			g := sym(name + "Down")
			ax, ay := g.anchor("stemDownSW")
			c.glyph(g, l.stemX-ax*s, l.stemY1-ay*s, s)
		}
	}
	if l.ev.Grace && l.ev.Slash {
		d := 1.0
		if !l.up {
			d = -1
		}
		y := l.stemY1 + d*1.6*s
		c.line(l.stemX-0.7*s, y+d*0.5*s, l.stemX+1.2*s, y-d*0.7*s, 0.12*s, 0)
	}
	ctx.tremolo(l)
}

// tremolo draws the strokes of a tremolo across the middle of the stem.
func (ctx *staffCtx) tremolo(l *evLayout) {
	if l.ev.Tremolo <= 0 || !l.stem {
		return
	}
	g := sym("tremolo" + string(rune('0'+min(l.ev.Tremolo, 3))))
	ctx.c.glyph(g, l.stemX+ctx.def.stem*l.scale/2, (l.stemY0+l.stemY1)/2, l.scale)
}

// drawBeam draws the beams of a group and the stems up to them.
func (ctx *staffCtx) drawBeam(g *beamGroup) {
	c := ctx.c
	evs := g.evs
	s := evs[0].scale
	t := ctx.def.stem * s
	th := ctx.def.beam * s
	gap := (ctx.def.beam + ctx.def.beamSpacing) * s
	first, last := evs[0], evs[len(evs)-1]
	xc := func(l *evLayout) float64 { return l.stemX + t/2 }
	// the ends the beam follows: the notes nearest it
	end := func(l *evLayout) float64 {
		if g.up {
			return float64(l.heads[len(l.heads)-1].pos) / 2
		}
		return float64(l.heads[0].pos) / 2
	}
	x0, x1 := xc(first), xc(last)
	d := end(last) - end(first)
	concave := false
	for _, l := range evs[1 : len(evs)-1] {
		y := end(l)
		if g.up && y < min(end(first), end(last)) || !g.up && y > max(end(first), end(last)) {
			concave = true
		}
	}
	if concave {
		d = 0
	}
	d = max(min(d, 1), -1)
	slope := 0.0
	if x1 > x0 {
		slope = d / (x1 - x0)
	}
	beams := 1
	for _, l := range evs {
		beams = max(beams, numberOfBeams(l.ev.Type))
	}
	minLen := (3.25 + float64(beams-1)*0.75) * s
	y0 := 0.0
	if g.up {
		y0 = 1e9
		for _, l := range evs {
			y0 = min(y0, end(l)-minLen-slope*(xc(l)-x0))
			if s == 1 {
				y0 = min(y0, 2-slope*(xc(l)-x0))
			}
		}
	} else {
		y0 = -1e9
		for _, l := range evs {
			y0 = max(y0, end(l)+minLen-slope*(xc(l)-x0))
			if s == 1 {
				y0 = max(y0, 2-slope*(xc(l)-x0))
			}
		}
	}
	yAt := func(x float64) float64 { return y0 + slope*(x-x0) }
	// stems to the outer edge of the beam
	for _, l := range evs {
		l.stemY1 = yAt(xc(l))
		if l.stem {
			c.rect(l.stemX, min(l.stemY0, l.stemY1), t, abs64(l.stemY1-l.stemY0))
		}
		ctx.tremolo(l)
	}
	left := func(l *evLayout) float64 { return l.stemX }
	right := func(l *evLayout) float64 { return l.stemX + t }
	seg := func(level int, xa, xb float64) {
		off := float64(level) * gap
		ya, yb := yAt(xa), yAt(xb)
		if g.up {
			c.fill(quad(xa, ya+off, xb, yb+off, xb, yb+off+th, xa, ya+off+th))
		} else {
			c.fill(quad(xa, ya-off, xb, yb-off, xb, yb-off-th, xa, ya-off-th))
		}
	}
	seg(0, left(first), right(last))
	hook := 1.1 * s
	for level := 1; level < beams; level++ {
		explicit := evs[0].ev.Beams != nil
		i := 0
		for i < len(evs) {
			if numberOfBeams(evs[i].ev.Type) <= level {
				i++
				continue
			}
			j := i
			for j+1 < len(evs) && numberOfBeams(evs[j+1].ev.Type) > level {
				if explicit && level < len(evs[j+1].ev.Beams) && evs[j+1].ev.Beams[level] == BeamBegin {
					break
				}
				j++
			}
			if j > i {
				seg(level, left(evs[i]), right(evs[j]))
			} else {
				// a single note gets a hook toward its neighbour
				forward := i == 0
				if explicit && level < len(evs[i].ev.Beams) {
					forward = evs[i].ev.Beams[level] == BeamForwardHook || evs[i].ev.Beams[level] == BeamBegin
				}
				if forward {
					seg(level, left(evs[i]), left(evs[i])+hook)
				} else {
					seg(level, right(evs[i])-hook, right(evs[i]))
				}
			}
			i = j + 1
		}
	}
}

func abs64(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
