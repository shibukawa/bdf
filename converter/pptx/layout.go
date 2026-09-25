package pptx

import (
	"math"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
)

// Line breaking follows the simple rules Office applies: breaks after runs
// of spaces, between East Asian characters (except before closing
// punctuation and small kana, and after opening brackets: kinsoku), between
// East Asian and other text, and after hyphens inside words.

const noStart = "、。，．,.:;?!)]}」』】〕〉》）］｝〙〗〟”’ゝゞヽヾーぁぃぅぇぉっゃゅょゎゕゖァィゥェォッャュョヮヵヶ々〻‐゠–〜？！：；・…‥%％°℃"
const noEnd = "([{「『【〔〈《（［｛〘〖〝“‘¥$£＄￥"

func isBreakSpace(r rune) bool {
	return r == ' ' || r == '　' || r == ' ' || r == ' ' || r == ' '
}

func canBreak(a, b rune) bool {
	if isBreakSpace(a) {
		return !isBreakSpace(b)
	}
	if isBreakSpace(b) || a == '\u00a0' || b == '\u00a0' {
		return false
	}
	if strings.ContainsRune(noStart, b) || strings.ContainsRune(noEnd, a) {
		return false
	}
	if fontdb.IsCJK(a) || fontdb.IsCJK(b) {
		return true
	}
	if (a == '-' || a == '‐') && unicode.IsLetter(b) {
		return true
	}
	return false
}

// markBreaks sets the break opportunities of a paragraph's items.
func markBreaks(items []item) {
	for i := range items {
		if items[i].kind != itemChar {
			continue
		}
		if i+1 >= len(items) || items[i+1].kind != itemChar {
			items[i].brk = true
			continue
		}
		items[i].brk = canBreak(items[i].r, items[i+1].r)
	}
}

// textLine is a laid-out line.
type textLine struct {
	pa        *para
	items     []item
	first     bool // first line of its paragraph
	afterBr   bool // follows an explicit line break
	last      bool // last line of its paragraph (or before a break)
	cjkWrap   bool // wrapped between characters that join without a space
	start     float64
	width     float64 // up to the end of the last visible item
	asc, desc float64
	height    float64
	baseline  float64
	offset    float64 // alignment shift
	bulletX   float64
}

// layoutParagraph breaks a paragraph into lines of at most width (no limit
// when wrap is false).
func layoutParagraph(pa *para, width float64, wrap bool) []*textLine {
	markBreaks(pa.items)
	firstStart := pa.marL + pa.indent
	bulletX := firstStart
	if b := pa.bullet; b != nil {
		switch {
		case pa.indent >= 0:
			firstStart = pa.marL + pa.indent + b.w
		case bulletX+b.w > pa.marL:
			firstStart = bulletX + b.w
		default:
			firstStart = pa.marL
		}
	}
	firstStart = math.Max(firstStart, 0)
	var lines []*textLine
	items := pa.items
	first, afterBr, cjkWrap := true, false, false
	for {
		start := pa.marL
		if first {
			start = firstStart
		}
		ln := &textLine{pa: pa, first: first, afterBr: afterBr, cjkWrap: cjkWrap, start: start, bulletX: bulletX}
		x := start
		end := len(items)
		lastBrk := -1
		brokeAt := -1
		for i := 0; i < len(items); i++ {
			it := &items[i]
			if it.kind == itemBreak {
				end = i
				brokeAt = i
				break
			}
			if it.kind == itemTab {
				it.w = nextTab(pa, x) - x
			}
			if wrap && x+it.w > width+1e-6 && !(it.kind == itemChar && isBreakSpace(it.r)) && i > 0 {
				if lastBrk >= 0 {
					end = lastBrk + 1
				} else {
					end = i
				}
				break
			}
			it.x = x
			x += it.w
			if it.brk {
				lastBrk = i
			}
		}
		ln.items = items[:end]
		// re-place (an earlier break point may have been chosen)
		x = start
		for i := range ln.items {
			it := &ln.items[i]
			if it.kind == itemTab {
				it.w = nextTab(pa, x) - x
			}
			it.x = x
			x += it.w
		}
		ln.width = start
		for i := len(ln.items) - 1; i >= 0; i-- {
			it := ln.items[i]
			if it.kind == itemChar && isBreakSpace(it.r) {
				continue
			}
			ln.width = it.x + it.w
			break
		}
		ln.width -= start
		lines = append(lines, ln)
		if brokeAt >= 0 {
			ln.last = true
			items = items[brokeAt+1:]
			first, afterBr, cjkWrap = false, true, false
			continue
		}
		if end >= len(items) {
			ln.last = true
			break
		}
		prev, next := items[end-1].r, items[end].r
		cjkWrap = !isBreakSpace(prev) && (fontdb.IsCJK(prev) || fontdb.IsCJK(next))
		items = items[end:]
		first, afterBr = false, false
	}
	return lines
}

func nextTab(pa *para, x float64) float64 {
	for _, t := range pa.tabs {
		if t > x+0.01 {
			return t
		}
	}
	return (math.Floor(x/pa.defTab+1e-9) + 1) * pa.defTab
}

// metrics sets the line's ascent, descent and height.
func (ln *textLine) metrics() {
	pa := ln.pa
	ln.asc, ln.desc = 0, 0
	visible := false
	for _, it := range ln.items {
		if it.kind == itemBreak {
			continue
		}
		size := it.st.size
		shift := 0.0
		if it.st.baseline != 0 {
			shift = it.st.baseline * size
			size *= 2.0 / 3
		}
		ln.asc = math.Max(ln.asc, it.fc.asc*size+shift)
		ln.desc = math.Max(ln.desc, it.fc.desc*size-shift)
		visible = true
	}
	if b := pa.bullet; b != nil && ln.first {
		ln.asc = math.Max(ln.asc, b.fc.asc*b.st.size)
		ln.desc = math.Max(ln.desc, b.fc.desc*b.st.size)
	}
	if !visible {
		ln.asc = math.Max(ln.asc, pa.endFace.asc*pa.end.size)
		ln.desc = math.Max(ln.desc, pa.endFace.desc*pa.end.size)
	}
	natural := ln.asc + ln.desc
	if pa.lnPts > 0 {
		ln.height = pa.lnPts
	} else {
		ln.height = natural * pa.lnPct
	}
}

// align positions a line horizontally in a box of the given width, and
// spreads justified lines.
func (ln *textLine) align(width float64, wrap bool) {
	pa := ln.pa
	avail := width - ln.start
	extra := avail - ln.width
	switch pa.algn {
	case "r":
		ln.offset = extra
	case "ctr":
		ln.offset = extra / 2
	case "just", "dist", "thaiDist", "justLow":
		if !wrap || extra <= 0 || (ln.last && pa.algn != "dist") {
			break
		}
		// Spread over the spaces between words, or over the characters
		// when there are none (East Asian text).
		n := len(ln.items)
		for n > 0 && ln.items[n-1].kind == itemChar && isBreakSpace(ln.items[n-1].r) {
			n--
		}
		spaces := 0
		for i := 0; i < n; i++ {
			if ln.items[i].kind == itemChar && isBreakSpace(ln.items[i].r) {
				spaces++
			}
		}
		if spaces > 0 {
			add := extra / float64(spaces)
			shift := 0.0
			for i := 0; i < len(ln.items); i++ {
				ln.items[i].x += shift
				if i < n && ln.items[i].kind == itemChar && isBreakSpace(ln.items[i].r) {
					ln.items[i].w += add
					ln.items[i].ls = -1 // stretched space: not drawn
					shift += add
				}
			}
		} else if n > 1 {
			add := extra / float64(n-1)
			shift := 0.0
			for i := 0; i < len(ln.items); i++ {
				ln.items[i].x += shift
				if i < n-1 && ln.items[i].kind == itemChar {
					ln.items[i].ls = add
					shift += add
				}
			}
		}
		ln.width = avail
	}
}

// laidOut is a text body after layout.
type laidOut struct {
	lines  []*textLine
	height float64
	width  float64 // widest line
}

// layoutBody lays out paragraphs in a column of the given width and returns
// the lines with baselines relative to the top of the text.
func layoutBody(paras []*para, width float64, wrap bool) *laidOut {
	lo := &laidOut{}
	y := 0.0
	for pi, pa := range paras {
		lines := layoutParagraph(pa, width, wrap)
		for li, ln := range lines {
			ln.metrics()
			if li == 0 && pi > 0 {
				y += spacing(pa.befPct, pa.befPts, ln.asc+ln.desc)
			}
			ln.align(width, wrap)
			ln.baseline = y + ln.height - ln.desc
			y += ln.height
			lo.width = math.Max(lo.width, ln.start+ln.offset+ln.width)
			lo.lines = append(lo.lines, ln)
		}
		if pi < len(paras)-1 && len(lines) > 0 {
			y += spacing(pa.aftPct, pa.aftPts, lines[len(lines)-1].asc+lines[len(lines)-1].desc)
		}
	}
	lo.height = y
	return lo
}

func spacing(p, pts, natural float64) float64 {
	if pts > 0 {
		return pts
	}
	return p * natural
}

// textEmitter writes text instructions and avoids repeating state.
type textEmitter struct {
	cv       *canvas
	font     bdf.FontRef
	fontSize float64
	hasFont  bool
	color    bdf.Color
	hasColor bool
	ls       float64
	links    []linkRect
	m        matrix // text space → slide, for link rectangles
	upright  bool   // East Asian vertical text: lines go to child objects
	vchars   bool   // inside such a line: CJK characters stand upright
}

type linkRect struct {
	x0, y0, x1, y1 float64
	url            string
}

func (e *textEmitter) setFont(fc *faceChoice, size float64) {
	ref := e.cv.font(fc.use)
	if !e.hasFont || ref != e.font || size != e.fontSize {
		e.cv.obj.Font(ref, f32(size))
		e.font, e.fontSize, e.hasFont = ref, size, true
	}
}

func (e *textEmitter) setColor(c bdf.Color) {
	if !e.hasColor || c != e.color {
		e.cv.obj.FillColor(c)
		e.color, e.hasColor = c, true
	}
}

func (e *textEmitter) setLetterSpacing(ls float64) {
	if ls != e.ls {
		e.cv.obj.TextStyle(bdf.AlignLeft, bdf.BaselineAlphabetic, bdf.DirInherit, f32(ls))
		e.ls = ls
	}
}

// emitLines draws laid-out lines whose baselines are offset by dy.
func (e *textEmitter) emitLines(lines []*textLine, dx, dy float64) {
	for i, ln := range lines {
		switch {
		case i == 0:
		case ln.first || ln.afterBr:
			e.cv.obj.Mark(bdf.MarkParagraph, "")
		case ln.cjkWrap:
			e.cv.obj.Mark(bdf.MarkWrap, "")
		default:
			e.cv.obj.Mark(bdf.MarkLine, "")
		}
		base := ln.baseline + dy
		if b := ln.pa.bullet; b != nil && ln.first {
			bx := ln.bulletX
			if ln.pa.algn == "ctr" || ln.pa.algn == "r" {
				bx = ln.start + ln.offset - b.w
			}
			if c, ok := b.st.color(); ok {
				e.setLetterSpacing(0)
				e.setFont(b.fc, b.st.size)
				e.setColor(c.bdf())
				e.cv.obj.FillText(b.text, f32(bx+dx), f32(base), f32(b.w))
				e.cv.drawn = true
				// the bullet and the text are separate words
				e.cv.obj.Mark(bdf.MarkLine, "")
			}
		}
		if e.upright {
			e.emitVerticalLine(ln, dx+ln.offset, base)
		} else {
			e.emitItems(ln, dx+ln.offset, base)
		}
	}
}

// emitVerticalLine draws a line of East Asian vertical text in a child
// object placed with USE_AT at the start of its baseline, where upright
// characters are turned back one by one; an ALT_TEXT before the USE_AT
// carries the line's text, so that extraction sees one run per line.
func (e *textEmitter) emitVerticalLine(ln *textLine, dx, base float64) {
	var b strings.Builder
	first := -1
	for i, it := range ln.items {
		if it.kind != itemChar {
			if it.kind == itemTab {
				b.WriteRune('\t')
			}
			continue
		}
		if first < 0 {
			first = i
		}
		b.WriteRune(it.r)
	}
	text := strings.TrimRight(b.String(), " \u3000\t")
	if first < 0 || text == "" {
		return
	}
	x0 := dx + ln.items[first].x
	size := ln.asc + ln.desc
	ch, ref := e.cv.child(bdf.Rect{X: 0, Y: f32(-ln.asc), W: f32(ln.width), H: f32(size)})
	ch.obj.SetBBox(-f32(size), f32(-ln.asc-size), f32(ln.width+2*size), f32(3*size))
	che := &textEmitter{cv: ch, m: e.m.mul(translate(x0, base)), vchars: true}
	che.emitItems(ln, dx-x0, 0)
	e.links = append(e.links, che.links...)
	e.cv.obj.Mark(bdf.MarkAltText, text)
	e.cv.obj.UseAt(ref, f32(x0), f32(base))
}

// emitItems draws a line's characters as runs of equal style.
func (e *textEmitter) emitItems(ln *textLine, dx, base float64) {
	items := ln.items
	for i := 0; i < len(items); {
		it := items[i]
		if it.kind != itemChar || it.ls < 0 {
			i++
			continue
		}
		j := i + 1
		for j < len(items) {
			n := items[j]
			if n.kind != itemChar || n.ls < 0 || (n.st != it.st && n.st.key != it.st.key) || n.fc != it.fc || n.ls != it.ls {
				break
			}
			if e.vchars && uprightInVertical(n.r) != uprightInVertical(it.r) {
				break
			}
			j++
		}
		// trailing spaces at the end of the line are not drawn
		k := j
		if j == len(items) {
			for k > i && isBreakSpace(items[k-1].r) {
				k--
			}
		}
		if k > i {
			e.emitRun(items[i:k], dx, base)
		}
		i = j
	}
}

func (e *textEmitter) emitRun(run []item, dx, base float64) {
	it := run[0]
	st := it.st
	var b strings.Builder
	adv := 0.0
	for _, r := range run {
		b.WriteRune(r.r)
		adv += r.w + r.ls
	}
	text := b.String()
	size := st.size
	shift := 0.0
	if st.baseline != 0 {
		shift = st.baseline * size
		size *= 2.0 / 3
	}
	x := it.x + dx
	y := base - shift
	if st.highlight != nil {
		e.cv.obj.FillColor(st.highlight.bdf())
		e.hasColor = false
		e.cv.obj.FillRect(f32(x), f32(y-it.fc.asc*size), f32(adv), f32((it.fc.asc+it.fc.desc)*size))
	}
	c, visible := st.color()
	if !visible {
		c = rgba{} // invisible, but still searchable
	}
	e.setFont(it.fc, size)
	e.setColor(c.bdf())
	if e.vchars && uprightInVertical(it.r) {
		e.emitUpright(run, x, y, size)
	} else {
		e.setLetterSpacing(st.spacing + it.ls)
		e.cv.obj.FillText(text, f32(x), f32(y), f32(adv))
	}
	e.cv.drawn = true
	if st.outline != nil {
		e.cv.obj.Save()
		e.cv.setLine(st.outline, adv, size)
		e.cv.obj.StrokeText(text, f32(x), f32(y), f32(adv))
		e.cv.obj.Restore()
	}
	lineColor := c
	if st.uColor != nil {
		lineColor = *st.uColor
	}
	th := math.Max(it.fc.ulTh*size, 0.5)
	if st.underline != "none" && st.underline != "" {
		e.setColor(lineColor.bdf())
		uy := y - it.fc.ulPos*size
		e.cv.obj.FillRect(f32(x), f32(uy), f32(adv-st.spacing-it.ls), f32(th))
		if strings.HasPrefix(st.underline, "dbl") {
			e.cv.obj.FillRect(f32(x), f32(uy+2*th), f32(adv-st.spacing-it.ls), f32(th))
		}
	}
	if st.strike == "sngStrike" || st.strike == "dblStrike" {
		e.setColor(c.bdf())
		sy := y - 0.3*size
		e.cv.obj.FillRect(f32(x), f32(sy), f32(adv), f32(th))
		if st.strike == "dblStrike" {
			e.cv.obj.FillRect(f32(x), f32(sy-2*th), f32(adv), f32(th))
		}
	}
	if st.link != "" {
		e.link(x, y-it.fc.asc*size, x+adv, y+it.fc.desc*size, st.link)
	}
}

// uprightInVertical reports whether a character stands upright in East
// Asian vertical text; Latin text, long vowel marks, dashes, brackets and
// ellipses are turned with the line instead.
func uprightInVertical(r rune) bool {
	if !fontdb.IsCJK(r) {
		return false
	}
	return !strings.ContainsRune("ー－―‐〜～…‥（）「」『』【】〔〕［］｛｝〈〉《》〘〙〖〗＝｜＿", r)
}

// verticalForms are the vertical presentation forms of East Asian
// punctuation; without them the horizontal glyph moves to the top right of
// its square, where vertical text puts it.
var verticalForms = map[rune]rune{'、': '︑', '。': '︒', '，': '︐', '．': '︒', '：': '︓', '；': '︔', '！': '︕', '？': '︖'}

// emitUpright draws East Asian characters of vertical text standing
// upright: each character's em box is turned back by 90° around its center.
func (e *textEmitter) emitUpright(run []item, x, y, size float64) {
	e.setLetterSpacing(0)
	fc := run[0].fc
	half := (fc.asc - fc.desc) * size / 2
	for _, r := range run {
		cx := r.x - run[0].x + x + r.w/2
		glyph, gfc := r.r, r.fc
		dx, dy := 0.0, 0.0
		if v, ok := verticalForms[r.r]; ok {
			if vf := e.cv.c.faceFor(r.st, v); vf.l != nil && vf.l.Has(v) {
				glyph, gfc = v, vf
				e.cv.c.advance(vf, v)
			} else {
				dx, dy = 0.55*size, -0.55*size
			}
		}
		e.setFont(gfc, size)
		e.cv.obj.Save()
		e.cv.obj.Transform(0, -1, 1, 0, f32(cx), f32(y-half))
		e.cv.obj.FillText(string(glyph), f32(-r.w/2+dx), f32(half+dy), f32(r.w-r.st.spacing))
		e.cv.obj.Restore()
	}
}

func (e *textEmitter) link(x0, y0, x1, y1 float64, url string) {
	xs, ys := make([]float64, 0, 4), make([]float64, 0, 4)
	for _, p := range [][2]float64{{x0, y0}, {x1, y0}, {x0, y1}, {x1, y1}} {
		px, py := e.m.apply(p[0], p[1])
		xs, ys = append(xs, px), append(ys, py)
	}
	r := linkRect{xs[0], ys[0], xs[0], ys[0], url}
	for i := range xs {
		r.x0, r.y0 = math.Min(r.x0, xs[i]), math.Min(r.y0, ys[i])
		r.x1, r.y1 = math.Max(r.x1, xs[i]), math.Max(r.y1, ys[i])
	}
	e.links = append(e.links, r)
}
