package docx

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/linebreak"
)

func lbAllowed(a, b rune) bool { return linebreak.Allowed(a, b) }
func lbSpace(r rune) bool      { return linebreak.IsSpace(r) }

// line is a laid-out line of a paragraph. Item positions are relative to
// the left edge of the area the paragraph is in (a column, a cell, a text
// box); baselines to the top of the line.
type line struct {
	p         *para
	items     []item
	label     []item // the list number, first line only
	first     bool
	last      bool    // last line of the paragraph
	endsBreak bool    // ends with a line break
	after     uint8   // kPage or kColumn when a break follows the line
	sep       byte    // how the line joins the previous one: 'L' space, 'W' none, 'P' paragraph break
	asc, desc float64 // of the tallest content
	height    float64
	baseline  float64
	left      float64 // the line's area
	right     float64
	used      float64 // right end of the visible content
}

// lineCtx is what breaking a paragraph's lines needs from where it is.
type lineCtx struct {
	defTab    float64
	grid      float64 // line pitch of the document grid (0: none)
	charPitch float64 // cell width of the character grid (0: none)
	allChars  bool    // every character takes whole cells, not only East Asian ones
	expand    bool    // justify lines that end with a line break
	vertical  bool    // East Asian vertical text: inline objects stand upright in the turned line
}

// indents returns a paragraph's left, right and first-line indents in
// points (the first line's relative to the left one).
func (p *para) indents() (l, r, first float64) {
	pp := p.pp
	l, r, first = pp.indL, pp.indR, pp.indFirst
	ch := p.mark.size // a character of the paragraph's font
	if pp.indLCh != 0 {
		l = pp.indLCh / 100 * ch
	}
	if pp.indRCh != 0 {
		r = pp.indRCh / 100 * ch
	}
	if pp.indFirstCh != 0 {
		first = pp.indFirstCh / 100 * ch
	}
	return
}

// nextTab returns the tab stop after x: a custom stop, the left indent
// when a hanging first line has not reached it, or a default stop.
func (lc *lineCtx) nextTab(p *para, x, indL float64, first bool) tabStop {
	const eps = 0.01
	best := tabStop{pos: math.Inf(1), align: "left", leader: "none"}
	for _, t := range p.pp.tabs {
		if t.align == "bar" {
			continue
		}
		if t.pos > x+eps {
			best = t
			break
		}
	}
	if first && indL > x+eps && indL < best.pos {
		best = tabStop{pos: indL, align: "left", leader: "none"}
	}
	if math.IsInf(best.pos, 1) {
		dt := lc.defTab
		if dt <= 0 {
			dt = 36
		}
		best = tabStop{pos: (math.Floor(x/dt+eps) + 1) * dt, align: "left", leader: "none"}
	}
	return best
}

// breakLine lays out the line of p that starts at item start, between the
// area positions left and right.
func (lc *lineCtx) breakLine(p *para, start int, first bool, left, right float64) *line {
	indL, _, indFirst := p.indents()
	ln := &line{p: p, first: first, left: left, right: right}
	x0 := left
	if first {
		x0 += indFirst
		if p.label != nil {
			ln.label = append([]item(nil), p.label.items...)
		}
	}
	// find the end: fill until the first item that does not fit (placing
	// the items up to well past the right edge is enough to find it)
	items := p.items[start:]
	end := len(items)
	pos := lc.place(p, ln, items, x0, indL, first, right+2*math.Max(right-x0, 100))
	if len(pos) < len(items) {
		end = len(pos)
	}
	lastBrk := -1
	for i := range pos {
		it := &pos[i]
		if it.kind == kBreak || it.kind == kPage || it.kind == kColumn {
			end = i + 1
			break
		}
		overflow := it.x+it.w+it.pad > right+0.01 && i > 0 && !(it.kind == kChar && lbSpace(it.r))
		if overflow && p.pp.overflowPunct && it.kind == kChar && hangingPunct(it.r) && it.x <= right+0.01 {
			overflow = false // hanging punctuation
		}
		if overflow {
			if lastBrk >= 0 {
				end = lastBrk + 1
			} else {
				end = i
			}
			break
		}
		if it.brk {
			lastBrk = i
		}
	}
	if end == 0 && len(items) > 0 {
		end = 1
	}
	// keep what follows a break opportunity that is only spaces on this line
	for end < len(items) && items[end].kind == kChar && lbSpace(items[end].r) {
		end++
	}
	ln.items = lc.place(p, ln, items[:end], x0, indL, first, math.Inf(1))
	if n := len(ln.items); n > 0 {
		switch last := ln.items[n-1]; last.kind {
		case kBreak:
			ln.endsBreak = true
		case kPage, kColumn:
			ln.after = last.kind
		}
	}
	ln.last = start+end >= len(p.items)
	ln.used = x0
	for i := len(ln.items) - 1; i >= 0; i-- {
		it := ln.items[i]
		if (it.kind == kChar && lbSpace(it.r)) || it.kind == kBreak || it.kind == kPage || it.kind == kColumn {
			continue
		}
		ln.used = it.x + it.w + it.pad
		break
	}
	if len(ln.label) > 0 {
		ln.used = math.Max(ln.used, ln.label[len(ln.label)-1].x+ln.label[len(ln.label)-1].w)
	}
	ln.align(lc)
	ln.metrics(lc)
	return ln
}

// hangingPunct are the characters that may hang past the right margin.
func hangingPunct(r rune) bool { return strings.ContainsRune("、。，．,.", r) }

// place positions items from x0, resolving tabs; the list number of a
// first line goes first. It returns positioned copies, up to the first item
// that starts past limit.
func (lc *lineCtx) place(p *para, ln *line, items []item, x0, indL float64, first bool, limit float64) []item {
	out := make([]item, len(items))
	copy(out, items)
	x := x0
	if len(ln.label) > 0 {
		lb := p.label
		lx := x
		switch lb.jc {
		case "center":
			lx -= lb.w / 2
		case "right":
			lx -= lb.w
		}
		for i := range ln.label {
			ln.label[i].x = lx
			lx += ln.label[i].w
		}
		x = math.Max(x, lx)
		switch lb.suffix {
		case "space":
			x += 0.25 * p.mark.size
		case "nothing":
		default:
			t := lc.nextTab(p, x, indL, true)
			x = t.pos
		}
	}
	pend := -1 // pending right, center or decimal tab
	var pendStop tabStop
	var pendStart float64
	finish := func(k int) {
		if pend < 0 {
			return
		}
		segW := x - pendStart
		var shift float64
		switch pendStop.align {
		case "right":
			shift = pendStop.pos - pendStart - segW
		case "center":
			shift = pendStop.pos - pendStart - segW/2
		case "decimal":
			before := segW
			for j := pend + 1; j < k; j++ {
				if out[j].kind == kChar && out[j].r == '.' {
					before = out[j].x - pendStart
					break
				}
			}
			shift = pendStop.pos - pendStart - before
		}
		if shift > 0 {
			out[pend].w = shift
			for j := pend + 1; j < k; j++ {
				out[j].x += shift
			}
			x += shift
		}
		pend = -1
	}
	for i := range out {
		it := &out[i]
		if x > limit && pend < 0 {
			return out[:i]
		}
		if it.kind == kObject && lc.vertical {
			// upright: its height goes along the line
			it.w = it.obj.h + it.obj.ext[1] + it.obj.ext[3]
		}
		inGroup := lc.vertical && it.tcy != 0 && i > 0 && out[i-1].tcy == it.tcy
		if lc.vertical && it.tcy != 0 {
			// horizontal in vertical text: the group takes one square
			it.w = 0
			if !inGroup {
				it.w = it.size
			}
		}
		if it.kind == kTab {
			finish(i)
			t := lc.nextTab(p, x, indL, first)
			it.x = x
			setLeader(it, t.leader)
			if t.align == "left" {
				it.w = t.pos - x
				x = t.pos
			} else {
				it.w = 0
				pend, pendStop, pendStart = i, t, x
			}
			continue
		}
		it.x = x
		if lc.charPitch > 0 && p.pp.snapToGrid && it.kind == kChar && !inGroup && (lc.allChars || fontdb.IsCJK(it.r) || lc.vertical && it.tcy != 0) {
			// the character takes whole cells of the character grid
			cells := math.Max(1, math.Ceil(it.w/lc.charPitch-0.01))
			it.pad = math.Max(0, cells*lc.charPitch-it.w)
		}
		x += it.w + it.pad + it.gap
	}
	finish(len(out))
	return out
}

// setLeader keeps the leader of a tab in its rune: '\t' for none, else the
// character the leader repeats.
func setLeader(it *item, leader string) {
	switch leader {
	case "dot":
		it.r = '.'
	case "hyphen":
		it.r = '-'
	case "underscore", "heavy":
		it.r = '_'
	case "middleDot":
		it.r = '·'
	default:
		it.r = '\t'
	}
}

// align moves a line for centered and right alignment and spreads
// justified lines.
func (ln *line) align(lc *lineCtx) {
	jc := ln.p.pp.jc
	if ln.p.pp.bidi {
		switch jc {
		case "left":
			jc = "right"
		case "right":
			jc = "left"
		}
	}
	extra := ln.right - ln.used
	var shift float64
	switch jc {
	case "center":
		shift = extra / 2
	case "right":
		shift = extra
	case "both", "distribute":
		if extra <= 0 {
			break
		}
		if jc == "both" && (ln.last || ln.after != 0 || (ln.endsBreak && !lc.expand)) {
			break
		}
		ln.justify(extra, jc == "distribute")
	}
	if shift > 0 {
		for i := range ln.items {
			ln.items[i].x += shift
		}
		for i := range ln.label {
			ln.label[i].x += shift
		}
		ln.used += shift
	}
}

// justify spreads extra space over the line after its last tab: over the
// spaces between words and the gaps next to East Asian characters (all
// gaps for distributed alignment).
func (ln *line) justify(extra float64, all bool) {
	items := ln.items
	from := 0
	for i, it := range items {
		if it.kind == kTab {
			from = i + 1
		}
	}
	n := len(items)
	for n > from && (items[n-1].kind != kChar && items[n-1].kind != kObject && items[n-1].kind != kField || items[n-1].kind == kChar && lbSpace(items[n-1].r)) {
		n--
	}
	if n-from < 1 {
		return
	}
	var spaces, gaps int
	for i := from; i < n; i++ {
		it := items[i]
		if it.kind == kChar && lbSpace(it.r) {
			spaces++
			continue
		}
		if i+1 < n && it.kind == kChar && items[i+1].kind == kChar && !lbSpace(items[i+1].r) &&
			(all || fontdb.IsCJK(it.r) || fontdb.IsCJK(items[i+1].r)) {
			gaps++
		}
	}
	if spaces+gaps == 0 {
		return
	}
	add := extra / float64(spaces+gaps)
	shift := 0.0
	for i := from; i < len(items); i++ {
		it := &items[i]
		it.x += shift
		if i >= n {
			continue
		}
		if it.kind == kChar && lbSpace(it.r) {
			it.w += add
			it.ls = -1
			shift += add
			continue
		}
		if i+1 < n && it.kind == kChar && items[i+1].kind == kChar && !lbSpace(items[i+1].r) &&
			(all || fontdb.IsCJK(it.r) || fontdb.IsCJK(items[i+1].r)) {
			if fontdb.IsCJK(it.r) || all {
				it.ls = add
			} else {
				// a Latin letter before East Asian text: the space goes
				// after it, and its word stays one run
				it.gap += add
			}
			shift += add
		}
	}
	ln.used += extra
}

// metrics sets the line's height and baseline by the paragraph's line
// spacing rule, snapped to the document grid.
func (ln *line) metrics(lc *lineCtx) {
	p := ln.p
	visible := false
	grow := func(it *item) {
		switch it.kind {
		case kChar, kField, kTab:
			if it.fc == nil {
				return
			}
			ln.asc = math.Max(ln.asc, it.fc.Asc*it.size+it.shift)
			ln.desc = math.Max(ln.desc, it.fc.Desc*it.size-it.shift)
			if it.kind != kTab {
				visible = true
			}
		case kObject:
			if lc.vertical {
				ln.asc = math.Max(ln.asc, it.obj.w+it.obj.ext[0]+it.obj.ext[2])
			} else {
				ln.asc = math.Max(ln.asc, it.obj.h+it.obj.ext[1]+it.obj.ext[3])
			}
			visible = true
		}
	}
	for i := range ln.items {
		grow(&ln.items[i])
	}
	for i := range ln.label {
		grow(&ln.label[i])
	}
	if !visible {
		fc := p.markFace
		if fc != nil {
			ln.asc = math.Max(ln.asc, fc.Asc*p.mark.size)
			ln.desc = math.Max(ln.desc, fc.Desc*p.mark.size)
		}
	}
	natural := ln.asc + ln.desc
	pp := p.pp
	grid := 0.0
	if pp.snapToGrid && lc.grid > 0 {
		grid = lc.grid
	}
	switch pp.lineRule {
	case "exact":
		ln.height = pp.line / 20
		if ln.height < natural && natural > 0 {
			ln.baseline = ln.asc * ln.height / natural
		} else {
			ln.baseline = ln.height - ln.desc
		}
	case "atLeast":
		min := pp.line / 20
		h := natural
		pad := 0.0
		if grid > 0 {
			g := math.Max(1, math.Ceil(natural/grid-0.01)) * grid
			pad = (g - natural) / 2
			h = g
		}
		ln.height = math.Max(min, h)
		ln.baseline = ln.height - ln.desc - pad
	default:
		mult := pp.line / 240
		if mult <= 0 {
			mult = 1
		}
		if grid > 0 {
			g := math.Max(1, math.Ceil(natural/grid-0.01)) * grid
			ln.height = g * mult
			ln.baseline = ln.height - (g-natural)/2 - ln.desc
		} else {
			ln.height = natural * mult
			ln.baseline = ln.height - ln.desc
		}
	}
}
