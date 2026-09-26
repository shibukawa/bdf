package xlsx

import (
	"math"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/linebreak"
)

// Cell text is laid out the way Excel lays it out: in the cell's box less a
// padding of 2 pixels on the left and 3 on the right (the gridline), on one
// line unless the cell wraps (at the break opportunities of DrawingML text,
// with kinsoku), aligned horizontally and vertically, indented, shrunk to
// fit or repeated to fill. Text that does not fit flows into the empty
// cells next to it; numbers never do: General numbers lose digits and
// other numbers turn into "#".

// maxShown is the number of characters of a value that a cell shows.
const maxShown = 1024

const (
	padL = 2 * pxPt
	padR = 3 * pxPt
	padV = 1 * pxPt
)

// tstyle is the resolved style of a run of cell text.
type tstyle struct {
	latin, ea string
	size      float64
	bold      bool
	italic    bool
	underline string
	strike    bool
	vert      string
	color     bdf.Color
}

type titem struct {
	r    rune
	st   *tstyle
	fc   *fontset.Choice
	w    float64
	kind byte // itemChar, itemGap, itemBreak, itemFill
	x    float64
	ls   float64 // letter spacing added after the character
}

const (
	itemChar = iota
	itemGap
	itemBreak
	itemFill
)

type tline struct {
	items    []titem
	x, base  float64
	width    float64
	asc      float64
	desc     float64
	wrapped  bool // continues the previous line after a wrap
	cjkWrap  bool // wrapped between characters that join without a space
	explicit bool // follows a line break in the text
}

// cellLayout is the text of a cell ready to draw.
type cellLayout struct {
	lines     []*tline
	textH     float64        // height the text needs (automatic row heights)
	clip      *box           // nil when the text stays inside its box
	m         *canvas.Matrix // rotated text: line coordinates → sheet
	bounds    box            // what the text may paint
	lang      string
	langKnown bool
	stacked   bool
}

// indentStep is the width of one indent level: about a character and a
// half of the default font.
func (s *sheetCtx) indentStep() float64 { return (s.c.mdw + 3) * pxPt }

// layoutCell lays out the value of a cell in a box (the cell or its merged
// range). measure skips what only drawing needs (overflow).
func (s *sheetCtx) layoutCell(r, c int, cl *cell, f *cellFmt, b box, measure bool) *cellLayout {
	c0 := s.c
	type seg struct {
		s    string
		font *xfont
		kind byte
	}
	var segs []seg
	numeric, general, boolErr := false, false, false
	var fmtColor *rgb
	var fm formatted
	switch cl.kind {
	case cellNum:
		if cl.num == 0 && !s.ws.showZero {
			return nil
		}
		fm = f.numFmt.formatNumber(cl.num, c0.date1904)
		numeric, general, fmtColor = true, fm.general, fm.color
		for _, p := range fm.pieces {
			segs = append(segs, seg{s: p.s, kind: p.kind})
		}
	case cellStr:
		switch {
		case cl.text.runs != nil:
			for _, run := range cl.text.runs {
				segs = append(segs, seg{s: run.text, font: run.font})
			}
		case f.numFmt.hasTextSection():
			tf := f.numFmt.formatText(cl.text.plain)
			fmtColor = tf.color
			for _, p := range tf.pieces {
				segs = append(segs, seg{s: p.s, kind: p.kind})
			}
		default:
			segs = []seg{{s: cl.text.plain}}
		}
	case cellBool:
		boolErr = true
		if cl.num != 0 {
			segs = []seg{{s: "TRUE"}}
		} else {
			segs = []seg{{s: "FALSE"}}
		}
	case cellErr:
		boolErr = true
		segs = []seg{{s: cl.text.plain}}
	default:
		return nil
	}
	al := f.align
	h := al.h
	if h == "general" || h == "" {
		switch {
		case numeric:
			h = "right"
		case boolErr:
			h = "center"
		default:
			h = "left"
		}
	}
	rot := al.rotation
	wrap := (al.wrap || h == "justify" || h == "distributed" || al.v == "justify" || al.v == "distributed") && rot == 0
	base := s.textStyle(&f.font, nil, fmtColor)
	if f.hideText {
		base.color = bdf.RGBA(0, 0, 0, 0)
	}
	styles := map[*xfont]*tstyle{}
	styleOf := func(x *xfont) *tstyle {
		if x == nil {
			return base
		}
		if st, ok := styles[x]; ok {
			return st
		}
		st := s.textStyle(&f.font, x, fmtColor)
		if f.hideText {
			st.color = bdf.RGBA(0, 0, 0, 0)
		}
		styles[x] = st
		return st
	}
	var items []titem
	fillAt, fillRune := -1, rune(0)
	shown := 0 // Excel shows the first 1,024 characters of a cell
	for _, sg := range segs {
		if shown >= maxShown {
			break
		}
		st := styleOf(sg.font)
		switch sg.kind {
		case pieceGap:
			for _, r := range sg.s {
				it := s.item(st, r)
				it.kind = itemGap
				items = append(items, it)
			}
			continue
		case pieceFill:
			if fillAt < 0 {
				fillAt = len(items)
				fillRune = []rune(sg.s)[0]
			}
			continue
		}
		for _, r := range sg.s {
			if r == '\n' {
				if wrap {
					items = append(items, titem{r: '\n', st: st, fc: s.c.fonts.FaceFor(st.latin, st.ea, st.bold, st.italic, ' '), kind: itemBreak})
				}
				continue
			}
			if r == '\t' {
				r = ' '
			}
			if r < 0x20 {
				continue
			}
			items = append(items, s.item(st, r))
			if shown++; shown >= maxShown {
				break
			}
		}
	}
	lang, known := s.cellLang(items)
	lay := &cellLayout{lang: lang, langKnown: known}
	avail := b.w - padL - padR
	indent := 0.0
	if al.indent > 0 && (h == "left" || h == "right" || h == "distributed") {
		indent = float64(al.indent) * s.indentStep()
		avail -= indent
		if h == "distributed" {
			avail -= indent
		}
	}
	avail = math.Max(avail, 0)
	width := func(its []titem) float64 {
		w := 0.0
		for _, it := range its {
			if it.kind != itemBreak {
				w += it.w
			}
		}
		return w
	}
	if rot == 0 && !wrap {
		w := width(items)
		switch {
		case numeric && w > avail+0.01:
			// numbers do not overflow
			items = s.fitNumber(cl.num, general, base, avail)
		case al.shrink && w > avail+0.01 && w > 0:
			items = s.shrink(items, avail/w)
		}
		if fillAt >= 0 {
			w := width(items)
			fit := s.item(base, fillRune)
			if unicode.IsSpace(fillRune) {
				fit.kind = itemGap // spacing, not text
			}
			if n := int((avail - w) / math.Max(fit.w, 0.01)); n > 0 && fillAt <= len(items) {
				fill := make([]titem, n)
				for i := range fill {
					fill[i] = fit
				}
				items = append(items[:fillAt:fillAt], append(fill, items[fillAt:]...)...)
			}
		}
		if h == "fill" {
			if w := width(items); w > 0 && w < avail {
				n := int(avail / w)
				one := items
				for i := 1; i < n; i++ {
					items = append(items, one...)
				}
			}
		}
	}
	if rot == 255 {
		return s.layoutStacked(lay, items, al, h, b)
	}
	var lines []*tline
	if wrap {
		lines = wrapItems(items, avail)
	} else {
		lines = []*tline{{items: items}}
	}
	for _, ln := range lines {
		lineMetrics(ln, base, s)
	}
	total := 0.0
	for _, ln := range lines {
		total += ln.asc + ln.desc
	}
	if rot != 0 {
		return s.layoutRotated(lay, lines[0], al, h, b, rot)
	}
	lay.textH = total + 2*padV - 2*pxPt
	if measure {
		return lay
	}
	lay.lines = lines
	// horizontal
	for i, ln := range lines {
		x := b.x + padL
		extra := avail - ln.width
		last := i == len(lines)-1 || lines[i+1].explicit
		switch h {
		case "right":
			x = b.x + b.w - padR - indent - ln.width
		case "center", "centerContinuous":
			x = b.x + (b.w-ln.width)/2
		case "justify":
			x += indent
			if !last {
				spread(ln, extra)
			}
		case "distributed":
			x += indent
			spread(ln, extra)
		default:
			x += indent
		}
		ln.x = x
	}
	// vertical
	y := b.y + b.h - padV - total
	switch al.v {
	case "top":
		y = b.y + padV
	case "center":
		y = b.y + (b.h-total)/2
	case "justify", "distributed":
		if len(lines) > 1 && total < b.h-2*padV {
			gap := (b.h - 2*padV - total) / float64(len(lines)-1)
			yy := b.y + padV
			for _, ln := range lines {
				ln.base = yy + ln.asc
				yy += ln.asc + ln.desc + gap
			}
			y = math.NaN()
		} else if al.v == "distributed" {
			y = b.y + (b.h-total)/2
		} else {
			y = b.y + padV
		}
	}
	if !math.IsNaN(y) {
		for _, ln := range lines {
			ln.base = y + ln.asc
			y += ln.asc + ln.desc
		}
	}
	// overflow into empty neighbors, or clip
	x0, x1 := b.x, b.x+b.w
	tx0, tx1 := math.Inf(1), math.Inf(-1)
	for _, ln := range lines {
		tx0, tx1 = math.Min(tx0, ln.x), math.Max(tx1, ln.x+ln.width)
	}
	if len(lines) == 0 {
		tx0, tx1 = b.x, b.x
	}
	if !wrap && !numeric && (tx0 < x0 || tx1 > x1) {
		x0, x1 = s.overflow(r, c, h, b, tx0, tx1)
	}
	ty0, ty1 := lines[0].base-lines[0].asc, lines[len(lines)-1].base+lines[len(lines)-1].desc
	if tx0 < x0-0.01 || tx1 > x1+0.01 || ty0 < b.y-0.01 || ty1 > b.y+b.h+0.01 {
		lay.clip = &box{x0, b.y, x1 - x0, b.h}
		lay.bounds = *lay.clip
	} else {
		lay.bounds = box{math.Min(tx0, b.x), b.y, math.Max(tx1, b.x+b.w) - math.Min(tx0, b.x), b.h}
	}
	return lay
}

// hasTextSection reports whether a format formats text (a fourth section
// or a section with @).
func (nf *numFormat) hasTextSection() bool {
	for i, sc := range nf.sections {
		if sc.text || i == 3 {
			return true
		}
	}
	return false
}

// textStyle resolves the style of cell text: the cell's font, changed by a
// rich text run's font, colored by the number format.
func (s *sheetCtx) textStyle(cellFont, run *xfont, fmtColor *rgb) *tstyle {
	f := *cellFont
	if run != nil {
		if run.name != "" {
			f.name, f.scheme = run.name, run.scheme
		}
		if run.size > 0 {
			f.size = run.size
		}
		f.bold, f.italic, f.underline, f.strike, f.vert = run.bold, run.italic, run.underline, run.strike, run.vert
		if run.color.kind != colorNone {
			f.color = run.color
		}
	}
	col := s.c.st.color(f.color, black)
	if fmtColor != nil {
		col = *fmtColor
	}
	return &tstyle{latin: s.c.fontName(&f), ea: s.c.eaFontName(&f), size: f.size, bold: f.bold, italic: f.italic,
		underline: f.underline, strike: f.strike, vert: f.vert, color: col.bdf()}
}

// item measures a character.
func (s *sheetCtx) item(st *tstyle, r rune) titem {
	fc := s.c.fonts.FaceFor(st.latin, st.ea, st.bold, st.italic, r)
	size := st.size
	if st.vert != "" {
		size *= 2.0 / 3
	}
	return titem{r: r, st: st, fc: fc, w: s.c.fonts.Advance(fc, r) * size}
}

// fitNumber shortens a number that does not fit: General shows fewer
// digits; otherwise, or when nothing fits, the cell fills with "#".
func (s *sheetCtx) fitNumber(v float64, general bool, st *tstyle, avail float64) []titem {
	if general {
		for n := 10; n >= 1; n-- {
			str := formatGeneral(v, n)
			if str == "" {
				continue
			}
			var its []titem
			w := 0.0
			for _, r := range str {
				it := s.item(st, r)
				its = append(its, it)
				w += it.w
			}
			if w <= avail+0.01 {
				return its
			}
		}
	}
	hash := s.item(st, '#')
	n := int(avail / math.Max(hash.w, 0.01))
	its := make([]titem, n)
	for i := range its {
		its[i] = hash
	}
	return its
}

// shrink scales text down to fit (shrink to fit).
func (s *sheetCtx) shrink(items []titem, k float64) []titem {
	scaled := map[*tstyle]*tstyle{}
	out := make([]titem, len(items))
	for i, it := range items {
		st, ok := scaled[it.st]
		if !ok {
			cp := *it.st
			cp.size *= k
			st = &cp
			scaled[it.st] = st
		}
		it.st = st
		it.w *= k
		out[i] = it
	}
	return out
}

// wrapItems breaks text into lines of at most width.
func wrapItems(items []titem, width float64) []*tline {
	var lines []*tline
	explicit := false
	wrapped, cjk := false, false
	for len(items) > 0 || len(lines) == 0 {
		x := 0.0
		end, lastBrk, brokeAt := len(items), -1, -1
		for i := 0; i < len(items); i++ {
			it := items[i]
			if it.kind == itemBreak {
				end, brokeAt = i, i
				break
			}
			if x+it.w > width+1e-6 && i > 0 && !(it.kind == itemChar && linebreak.IsSpace(it.r)) {
				if lastBrk >= 0 {
					end = lastBrk + 1
				} else {
					end = i
				}
				break
			}
			x += it.w
			if it.kind == itemChar && (i+1 >= len(items) || items[i+1].kind != itemChar || linebreak.Allowed(it.r, items[i+1].r)) {
				lastBrk = i
			}
		}
		ln := &tline{items: items[:end], explicit: explicit, wrapped: wrapped, cjkWrap: cjk}
		lines = append(lines, ln)
		if brokeAt >= 0 {
			items = items[brokeAt+1:]
			explicit, wrapped, cjk = true, false, false
			if len(items) == 0 {
				// a trailing line break leaves an empty line
				lines = append(lines, &tline{explicit: true})
				break
			}
			continue
		}
		if end >= len(items) {
			break
		}
		prev, next := items[end-1].r, items[end].r
		cjk = linebreak.Joins(prev, next)
		explicit, wrapped = false, true
		items = items[end:]
		// spaces at the start of a wrapped line are dropped
		for len(items) > 0 && items[0].kind == itemChar && linebreak.IsSpace(items[0].r) {
			items = items[1:]
		}
		if len(items) == 0 {
			break
		}
	}
	return lines
}

// lineMetrics positions a line's items and measures it.
func lineMetrics(ln *tline, base *tstyle, s *sheetCtx) {
	x := 0.0
	for i := range ln.items {
		ln.items[i].x = x
		x += ln.items[i].w
	}
	// the width ends at the last visible item (not trailing spaces)
	ln.width = 0
	for i := len(ln.items) - 1; i >= 0; i-- {
		it := ln.items[i]
		if it.kind == itemChar && linebreak.IsSpace(it.r) {
			continue
		}
		ln.width = it.x + it.w
		break
	}
	for _, it := range ln.items {
		size := it.st.size
		ln.asc = math.Max(ln.asc, it.fc.Asc*size)
		ln.desc = math.Max(ln.desc, it.fc.Desc*size)
	}
	if ln.asc == 0 && ln.desc == 0 {
		fc := s.c.fonts.FaceFor(base.latin, base.ea, base.bold, base.italic, 'x')
		ln.asc, ln.desc = fc.Asc*base.size, fc.Desc*base.size
	}
}

// spread justifies a line: extra space goes between words, or between the
// characters of text without spaces.
func spread(ln *tline, extra float64) {
	if extra <= 0 {
		return
	}
	n := len(ln.items)
	for n > 0 && ln.items[n-1].kind == itemChar && linebreak.IsSpace(ln.items[n-1].r) {
		n--
	}
	spaces := 0
	for i := 0; i < n; i++ {
		if ln.items[i].kind == itemChar && linebreak.IsSpace(ln.items[i].r) {
			spaces++
		}
	}
	shift := 0.0
	switch {
	case spaces > 0:
		add := extra / float64(spaces)
		for i := range ln.items {
			ln.items[i].x += shift
			if i < n && ln.items[i].kind == itemChar && linebreak.IsSpace(ln.items[i].r) {
				// a widened space is not drawn: the words keep their advances
				ln.items[i].w += add
				ln.items[i].kind = itemGap
				shift += add
			}
		}
	case n > 1:
		add := extra / float64(n-1)
		for i := range ln.items {
			ln.items[i].x += shift
			if i < n-1 {
				ln.items[i].ls = add
				shift += add
			}
		}
	case n == 1:
		// one character: centered
		ln.items[0].x += extra / 2
	}
	ln.width += extra
}

// overflow extends a line of text past its box into the empty cells beside
// it, in the direction its alignment grows. It returns the horizontal
// range the text may paint.
func (s *sheetCtx) overflow(r, c int, h string, b box, tx0, tx1 float64) (float64, float64) {
	x0, x1 := b.x, b.x+b.w
	last := c
	if m, ok := s.mergeAt(r, c); ok {
		last = m.c1
		if m.r1 > m.r0 {
			return x0, x1 // multi-row merged cells clip
		}
	}
	right := h == "left" || h == "center" || h == "centerContinuous" || h == "justify" || h == "distributed"
	left := h == "right" || h == "center" || h == "centerContinuous"
	if right {
		for cc := last + 1; cc < maxCols && x1 < tx1; cc++ {
			if !s.emptyAt(r, cc) {
				break
			}
			x1 = s.cols.pos(cc + 1)
		}
	}
	if left {
		first := c
		if m, ok := s.mergeAt(r, c); ok {
			first = m.c0
		}
		for cc := first - 1; cc >= 0 && x0 > tx0; cc-- {
			if !s.emptyAt(r, cc) {
				break
			}
			x0 = s.cols.pos(cc)
		}
	}
	return x0, x1
}

// emptyAt reports whether text may flow over a cell: it has no value and
// is not merged.
func (s *sheetCtx) emptyAt(r, c int) bool {
	if cl := s.ws.cellAt(r, c); cl != nil && cl.kind != cellBlank {
		if cl.kind == cellStr && cl.text.plain == "" {
			return true
		}
		return false
	}
	_, merged := s.mergeAt(r, c)
	return !merged
}

// layoutRotated lays a line of text out turned by the cell's text rotation
// (1–90 counterclockwise, 91–180 clockwise by the value less 90), placed in
// the box by the alignment.
func (s *sheetCtx) layoutRotated(lay *cellLayout, ln *tline, al xalign, h string, b box, rot int) *cellLayout {
	deg := float64(rot)
	if rot > 90 {
		deg = -float64(rot - 90)
	}
	m := canvas.Rotate(-deg)
	bx0, by0, bx1, by1 := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	for _, p := range [][2]float64{{0, -ln.asc}, {ln.width, -ln.asc}, {0, ln.desc}, {ln.width, ln.desc}} {
		x, y := m.Apply(p[0], p[1])
		bx0, by0, bx1, by1 = math.Min(bx0, x), math.Min(by0, y), math.Max(bx1, x), math.Max(by1, y)
	}
	var tx, ty float64
	switch h {
	case "right":
		tx = b.x + b.w - padR - bx1
	case "left", "fill", "justify", "distributed":
		tx = b.x + padL - bx0
	default:
		tx = b.x + (b.w-(bx1-bx0))/2 - bx0
	}
	switch al.v {
	case "top":
		ty = b.y + padV - by0
	case "center", "justify", "distributed":
		ty = b.y + (b.h-(by1-by0))/2 - by0
	default:
		ty = b.y + b.h - padV - by1
	}
	mm := canvas.Translate(tx, ty).Mul(m)
	lay.m = &mm
	ln.x, ln.base = 0, 0
	lay.lines = []*tline{ln}
	lay.textH = by1 - by0 + 2*padV - 2*pxPt
	if bx1-bx0 > b.w+0.01 || by1-by0 > b.h+0.01 {
		lay.clip = &b
	}
	lay.bounds = b
	return lay
}

// layoutStacked lays text out as a column of upright characters (text
// rotation 255).
func (s *sheetCtx) layoutStacked(lay *cellLayout, items []titem, al xalign, h string, b box) *cellLayout {
	lay.stacked = true
	total := 0.0
	var lines []*tline
	for _, it := range items {
		if it.kind != itemChar {
			continue
		}
		ln := &tline{items: []titem{it}, width: it.w, asc: it.fc.Asc * it.st.size, desc: it.fc.Desc * it.st.size}
		lines = append(lines, ln)
		total += ln.asc + ln.desc
	}
	y := b.y + b.h - padV - total
	switch al.v {
	case "top":
		y = b.y + padV
	case "center", "justify", "distributed":
		y = b.y + (b.h-total)/2
	}
	for _, ln := range lines {
		switch h {
		case "left":
			ln.x = b.x + padL
		case "right":
			ln.x = b.x + b.w - padR - ln.width
		default:
			ln.x = b.x + (b.w-ln.width)/2
		}
		ln.items[0].x = 0
		ln.base = y + ln.asc
		y += ln.asc + ln.desc
	}
	lay.lines = lines
	lay.textH = total + 2*padV - 2*pxPt
	if total > b.h {
		lay.clip = &b
	}
	lay.bounds = b
	return lay
}

// cellLang returns the language of a cell's text: that of its kana or
// hangul, the workbook's East Asian language for other East Asian text,
// and the document's for the rest; false for East Asian text of unknown
// language.
func (s *sheetCtx) cellLang(items []titem) (string, bool) {
	cjk := false
	for _, it := range items {
		if l := fontset.RuneLang(it.r); l != "" {
			return l, true
		}
		if fontdb.IsCJK(it.r) && unicode.Is(unicode.Han, it.r) {
			cjk = true
		}
	}
	if !cjk {
		return "", true
	}
	if s.c.eaLang != "" {
		return s.c.eaLang, true
	}
	return "", false
}

// emitText draws a cell's laid-out text into a tile, after a CELL mark
// with ref.
func (s *sheetCtx) emitText(t *tileCv, lay *cellLayout, ref string) {
	t.cv.Obj.Mark(bdf.MarkCell, ref)
	if lay.langKnown {
		t.cv.SetLang(lay.lang)
	} else if fontset.Script(t.cv.Lang()) == "" {
		t.cv.SetLang("")
	}
	depth := len(t.stack)
	if lay.clip != nil {
		t.save()
		cb := lay.clip
		t.cv.Obj.ClipRect(f32(cb.x-t.ox), f32(cb.y-t.oy), f32(cb.w), f32(cb.h))
	}
	ox, oy := t.ox, t.oy
	if lay.m != nil {
		t.save()
		t.cv.Transform(canvas.Translate(-t.ox, -t.oy).Mul(*lay.m))
		ox, oy = 0, 0
	}
	for i, ln := range lay.lines {
		if i > 0 {
			switch {
			case lay.stacked || ln.cjkWrap:
				t.cv.Obj.Mark(bdf.MarkWrap, "")
			default:
				t.cv.Obj.Mark(bdf.MarkLine, "")
			}
		}
		s.emitLine(t, ln, ox, oy)
	}
	for len(t.stack) > depth {
		t.restore()
	}
	t.letterSpacing(0)
}

// emitLine draws a line as runs of characters of one style and face.
func (s *sheetCtx) emitLine(t *tileCv, ln *tline, ox, oy float64) {
	items := ln.items
	n := len(items)
	for n > 0 && items[n-1].kind == itemChar && linebreak.IsSpace(items[n-1].r) {
		n--
	}
	for i := 0; i < n; {
		it := items[i]
		if it.kind != itemChar {
			i++
			continue
		}
		j := i + 1
		for j < n {
			nx := items[j]
			if nx.kind != itemChar || nx.st != it.st || nx.fc != it.fc || nx.ls != it.ls {
				break
			}
			j++
		}
		s.emitRun(t, items[i:j], ln, ox, oy)
		i = j
	}
}

func (s *sheetCtx) emitRun(t *tileCv, run []titem, ln *tline, ox, oy float64) {
	it := run[0]
	st := it.st
	var b strings.Builder
	adv := 0.0
	for _, r := range run {
		b.WriteRune(r.r)
		adv += r.w + r.ls
	}
	size := st.size
	shift := 0.0
	switch st.vert {
	case "superscript":
		size *= 2.0 / 3
		shift = st.size * 0.33
	case "subscript":
		size *= 2.0 / 3
		shift = -st.size * 0.15
	}
	x := ln.x + it.x - ox
	y := ln.base - shift - oy
	t.setFont(it.fc, size)
	t.fillColor(st.color)
	t.letterSpacing(it.ls)
	t.cv.Obj.FillText(b.String(), f32(x), f32(y), f32(adv))
	t.cv.Drawn = true
	th := math.Max(it.fc.ULThick*size, pxPt)
	w := adv - it.ls
	if st.underline != "" {
		uy := y - it.fc.ULPos*size
		t.cv.Obj.FillRect(f32(x), f32(uy), f32(w), f32(th))
		if strings.HasPrefix(st.underline, "double") {
			t.cv.Obj.FillRect(f32(x), f32(uy+2*th), f32(w), f32(th))
		}
	}
	if st.strike {
		t.cv.Obj.FillRect(f32(x), f32(y-0.3*size), f32(w), f32(th))
	}
}

// paintText lays out and draws the values of the cells.
func (s *sheetCtx) paintText() {
	for i := range s.ws.rows {
		rw := &s.ws.rows[i]
		r := rw.idx
		if r >= s.nRows {
			break
		}
		for j := range rw.cells {
			cl := &rw.cells[j]
			if cl.kind == cellBlank || cl.col >= s.nCols {
				continue
			}
			b := s.cellBox(r, cl.col)
			ref := cellRef(r, cl.col)
			if m, ok := s.mergeAt(r, cl.col); ok {
				if m.r0 != r || m.c0 != cl.col {
					continue // hidden under the merged range
				}
				b = s.rangeBox(m)
				ref += ":" + cellRef(m.r1, m.c1)
			}
			if b.w <= 0 || b.h <= 0 {
				continue
			}
			f := s.formatAt(r, cl.col, rw, cl)
			if f == nil {
				f = s.fmtOf(0)
			}
			if f.align.h == "centerContinuous" {
				b = s.centerAcross(r, cl.col, b)
			}
			lay := s.layoutCell(r, cl.col, cl, f, b, false)
			if lay == nil || len(lay.lines) == 0 {
				continue
			}
			if t := s.tableAt(r, cl.col); t != nil && t.isHeader(r) {
				ref += " col"
			}
			bb := lay.bounds
			if right := bb.x + bb.w; right > s.cols.pos(s.nCols) {
				s.textCols = max(s.textCols, s.cols.index(right-0.01)+1)
			}
			for _, t := range s.tilesIn(bb.x, bb.y, bb.x+bb.w, bb.y+bb.h) {
				s.emitText(t, lay, ref)
			}
		}
	}
}

// centerAcross widens the box of a "center across selection" cell over the
// empty cells after it that center across too.
func (s *sheetCtx) centerAcross(r, c int, b box) box {
	for cc := c + 1; cc < s.nCols; cc++ {
		cl := s.ws.cellAt(r, cc)
		if cl == nil || cl.kind != cellBlank {
			break
		}
		if f := s.fmtOf(cl.style); f.align.h != "centerContinuous" {
			break
		}
		b.w = s.cols.pos(cc+1) - b.x
	}
	return b
}
