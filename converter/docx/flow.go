package docx

import (
	"math"

	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/linebreak"
)

// op is one thing a layout draws, positioned in the coordinates of the
// area it was laid out in (a page, the scroll view's column, a cell, a
// text box). Ops are drawn in order; the structure of the text (spec
// §7.8) is derived from their contexts, so that a table or list cut by a
// page break is reopened on the next.
type op struct {
	y0, y1 float64  // vertical extent
	text   bool     // draws lines of text: never cut through
	ctx    []*frame // structure the text is in (text ops only)
	top    bool     // drawn outside any structure (floating objects)
	leaf   byte     // mark before the text: 'P' paragraph, 'H' heading, 'L' line, 'W' wrap
	level  int      // heading level
	dx, dy float64
	fn     func(e *emitter, dx, dy float64)
}

func (o *op) shift(dx, dy float64) {
	o.dx += dx
	o.dy += dy
	o.y0 += dy
	o.y1 += dy
}

// page is a laid-out page of the page view.
type page struct {
	sec           *section
	secNo         int
	num           int  // the number printed on it
	first         bool // first page of its section
	blank         bool // inserted to start a section on an odd or even page
	body          []op
	behind, front []op // floating objects
	notes         []*noteUse
	noteH         float64
	top, bottom   float64 // the body area
}

// A page of vertical text is laid out as a horizontal page turned by 90°
// clockwise: the flow's x runs down the page from the top of the body (the
// lines), its y runs leftwards from the body's right edge (the order of the
// lines). Headers, footers and floating objects stay upright, in page
// coordinates.

// lineLen is the length of the page's lines.
func (pg *page) lineLen() float64 {
	if pg.sec.vertical {
		return pg.bottom - pg.top
	}
	return pg.sec.textWidth()
}

// block returns where the lines of the body begin and end in the flow's
// coordinates: the top and bottom of the body, or for vertical text its
// right and left edges.
func (pg *page) block() (top, bottom float64) {
	if pg.sec.vertical {
		return 0, pg.sec.textWidth()
	}
	return pg.top, pg.bottom
}

// rotation maps the flow's coordinates of a page of vertical text onto the
// page.
func (pg *page) rotation() canvas.Matrix {
	return canvas.Matrix{0, 1, -1, 0, pg.sec.pgW - pg.sec.right, pg.top}
}

// excl is the area a floating object keeps text out of.
type excl struct {
	x0, y0, x1, y1 float64
	full           bool   // no text beside it (top and bottom wrapping)
	side           string // bothSides, left, right, largest
}

// flow places blocks: into pages and columns (the page view), or into one
// area of unlimited height (the scroll view, cells, text boxes, headers).
type flow struct {
	c      *converter
	paged  bool
	main   bool        // the page or scroll view's flow (not a cell, box or header)
	origin *[2]float64 // where the area's origin is on the page, when it is on one
	lc     lineCtx

	pages []*page
	pg    *page
	sec   *section
	secNo int
	col   int

	x0, w          float64 // the current column: left edge and width
	y, top, bottom float64 // position, and the column's top and bottom
	colTop         float64 // where the columns of the section begin on the page
	colBottom      float64 // the lowest column end of the section on the page

	ops, behind, front []op // unpaged output

	floats    []excl
	prevPara  *para
	prevAfter float64 // space after the previous paragraph, still to add
	atTop     bool    // nothing placed in the column yet
	hard      bool    // the column began after a page, column or section break
	inChain   bool
	repeating bool // drawing the header rows of a table again

	pageStart int  // pending page number of a new section (-1: none)
	secFirst  bool // the next page is its section's first
}

// subflow returns an unpaged flow of the given width, for text that is
// vertical or not.
func (c *converter) subflow(width float64, sec *section, vertical bool) *flow {
	f := &flow{c: c, w: width, sec: sec, atTop: true, pageStart: -1}
	f.lc = c.lineCtx(sec, vertical)
	f.bottom = math.Inf(1)
	return f
}

func (c *converter) lineCtx(sec *section, vertical bool) lineCtx {
	lc := lineCtx{defTab: c.defTab, expand: !c.noExpandShiftReturn, vertical: vertical}
	if sec != nil {
		lc.grid = sec.grid()
		switch sec.gridType {
		case "linesAndChars", "snapToChars":
			lc.charPitch = c.baseSize + sec.charSpace
			lc.allChars = sec.gridType == "snapToChars"
		}
	}
	return lc
}

// emit adds a body op.
func (f *flow) emit(o op) {
	if f.paged {
		f.pg.body = append(f.pg.body, o)
	} else {
		f.ops = append(f.ops, o)
	}
}

func (f *flow) out() *[]op {
	if f.paged {
		return &f.pg.body
	}
	return &f.ops
}

// all returns an unpaged flow's ops: floating objects behind the text,
// the text, floating objects in front.
func (f *flow) all() []op {
	out := make([]op, 0, len(f.behind)+len(f.ops)+len(f.front))
	out = append(out, f.behind...)
	out = append(out, f.ops...)
	return append(out, f.front...)
}

// position is where something was placed: a column of a page.
type position struct{ page, col int }

func (f *flow) pos() position { return position{len(f.pages), f.col} }

// flowMark is a snapshot of a paged flow, to lay content out again.
type flowMark struct {
	npages                  int
	nbody, nbehind, nfront  int
	nnotes                  int
	noteH                   float64
	y, top, bottom, colTop  float64
	colBottom               float64
	x0, w                   float64
	col                     int
	nfloats                 int
	prevPara                *para
	prevAfter               float64
	atTop, hard             bool
	sec                     *section
	secNo, pageStart        int
	secFirst                bool
	nops, nbehindU, nfrontU int
}

func (f *flow) mark() flowMark {
	m := flowMark{npages: len(f.pages), y: f.y, top: f.top, bottom: f.bottom, colTop: f.colTop, colBottom: f.colBottom, x0: f.x0, w: f.w,
		col: f.col, nfloats: len(f.floats), prevPara: f.prevPara, prevAfter: f.prevAfter, atTop: f.atTop, hard: f.hard,
		sec: f.sec, secNo: f.secNo, pageStart: f.pageStart, secFirst: f.secFirst,
		nops: len(f.ops), nbehindU: len(f.behind), nfrontU: len(f.front)}
	if f.pg != nil {
		m.nbody, m.nbehind, m.nfront, m.nnotes, m.noteH = len(f.pg.body), len(f.pg.behind), len(f.pg.front), len(f.pg.notes), f.pg.noteH
	}
	return m
}

func (f *flow) restore(m flowMark) {
	f.pages = f.pages[:m.npages]
	if m.npages > 0 {
		f.pg = f.pages[m.npages-1]
		f.pg.body, f.pg.behind, f.pg.front = f.pg.body[:m.nbody], f.pg.behind[:m.nbehind], f.pg.front[:m.nfront]
		f.pg.notes, f.pg.noteH = f.pg.notes[:m.nnotes], m.noteH
	} else {
		f.pg = nil
	}
	f.y, f.top, f.bottom, f.colTop, f.x0, f.w, f.col = m.y, m.top, m.bottom, m.colTop, m.x0, m.w, m.col
	f.colBottom = m.colBottom
	f.floats = f.floats[:m.nfloats]
	f.prevPara, f.prevAfter, f.atTop, f.hard = m.prevPara, m.prevAfter, m.atTop, m.hard
	f.sec, f.secNo, f.pageStart, f.secFirst = m.sec, m.secNo, m.pageStart, m.secFirst
	f.ops, f.behind, f.front = f.ops[:m.nops], f.behind[:m.nbehindU], f.front[:m.nfrontU]
}

// newPage starts a page of the current section.
func (f *flow) newPage(hard bool) {
	sec := f.sec
	num := 1
	if n := len(f.pages); n > 0 {
		num = f.pages[n-1].num + 1
	}
	if f.pageStart >= 0 {
		num = f.pageStart
		f.pageStart = -1
	}
	pg := &page{sec: sec, secNo: f.secNo, num: num, first: f.secFirst}
	f.secFirst = false
	pg.top, pg.bottom = f.c.bodyArea(pg)
	f.pages = append(f.pages, pg)
	f.pg = pg
	f.col = 0
	top, bottom := pg.block()
	f.colTop, f.top, f.bottom = top, top, bottom
	f.colBottom = 0
	f.x0, f.w = sec.column(0, pg.lineLen())
	f.y = top
	f.atTop, f.hard = true, hard
	f.floats = nil
	f.prevAfter = 0
}

// nextColumn moves to the next column, or to a new page after the last.
func (f *flow) nextColumn(hard bool) {
	if !f.paged {
		return
	}
	if f.col+1 < len(f.sec.cols) {
		f.colBottom = math.Max(f.colBottom, f.y)
		f.col++
		f.x0, f.w = f.sec.column(f.col, f.pg.lineLen())
		f.y, f.top = f.colTop, f.colTop
		f.atTop, f.hard = true, hard
		f.prevAfter = 0
		return
	}
	f.newPage(hard)
}

// limit is the bottom of the space left in the column.
func (f *flow) limit() float64 {
	if !f.paged {
		return math.Inf(1)
	}
	return f.bottom - f.pg.noteH
}

// sections lays out a document's sections into pages.
func (f *flow) sections(secs []*section) {
	for i, s := range secs {
		prev := f.sec
		f.sec, f.secNo = s, i
		f.lc = f.c.lineCtx(s, s.vertical)
		if s.pgStart >= 0 {
			f.pageStart = s.pgStart
		}
		switch {
		case i == 0:
			f.secFirst = true
			f.newPage(true)
		case (s.typ == "continuous" || s.typ == "nextColumn") && sameGeometry(prev, s):
			// on the same page, below what is there (in the next column
			// for a column break)
			if s.typ == "nextColumn" && f.col+1 < len(s.cols) {
				f.col++
			} else {
				f.col = 0
			}
			f.x0, f.w = s.column(f.col, f.pg.lineLen())
			if f.col == 0 {
				f.colTop = math.Max(f.y, f.colBottom)
				f.colBottom = 0
			}
			f.y, f.top = f.colTop, f.colTop
		default:
			f.secFirst = true
			f.newPage(true)
			if (s.typ == "oddPage" && f.pg.num%2 == 0) || (s.typ == "evenPage" && f.pg.num%2 == 1) {
				f.pg.blank = true
				f.secFirst = true
				f.newPage(true)
			}
		}
		f.blocks(s.blocks, nil)
	}
}

func sameGeometry(a, b *section) bool {
	return a != nil && a.pgW == b.pgW && a.pgH == b.pgH && a.left == b.left && a.right == b.right && a.vertical == b.vertical
}

// blocks places blocks one after another. Paragraphs kept with the next
// block move to the next column with it when the pair would be split.
func (f *flow) blocks(bs []block, next block) {
	for i := 0; i < len(bs); i++ {
		after := next
		if i+1 < len(bs) {
			after = bs[i+1]
		}
		if f.paged && !f.inChain && keepNext(bs[i]) && i+1 < len(bs) {
			j := i
			for j < len(bs) && keepNext(bs[j]) {
				j++
			}
			if j < len(bs) {
				m := f.mark()
				f.inChain = true
				first := f.placeBlock(bs[i], bs[i+1])
				var last position
				for k := i + 1; k <= j; k++ {
					nx := next
					if k+1 < len(bs) {
						nx = bs[k+1]
					}
					last = f.placeBlock(bs[k], nx)
				}
				if first != last && !m.atTop {
					f.restore(m)
					f.nextColumn(false)
					for k := i; k <= j; k++ {
						nx := next
						if k+1 < len(bs) {
							nx = bs[k+1]
						}
						f.placeBlock(bs[k], nx)
					}
				}
				f.inChain = false
				i = j
				continue
			}
		}
		f.placeBlock(bs[i], after)
	}
}

func keepNext(b block) bool {
	p, ok := b.(*para)
	return ok && p.pp.keepNext && p.sect == nil
}

// placeBlock places a block and returns where its first line went.
func (f *flow) placeBlock(b block, next block) position {
	switch b := b.(type) {
	case *para:
		return f.paragraph(b, next)
	case *table:
		return f.table(b)
	}
	return f.pos()
}

// spacing returns the space before and after a paragraph.
func (f *flow) spacing(p *para) (before, after float64) {
	pp := p.pp
	lineUnit := 12.0
	if f.lc.grid > 0 {
		lineUnit = f.lc.grid
	}
	before, after = pp.before, pp.after
	if pp.beforeLn != 0 {
		before = pp.beforeLn / 100 * lineUnit
	}
	if pp.afterLn != 0 {
		after = pp.afterLn / 100 * lineUnit
	}
	if pp.beforeAuto {
		before = 14
	}
	if pp.afterAuto {
		after = 14
	}
	return math.Max(before, 0), math.Max(after, 0)
}

// placed tells where the lines of a paragraph went.
type placed struct {
	first, last position
	segs        []int // lines per column
	broken      bool  // split by a page or column break of its own
}

// paragraph places a paragraph with Word's rules for keeping its lines
// together (keepLines) and for widows and orphans.
func (f *flow) paragraph(p *para, next block) position {
	if !f.paged {
		return f.placePara(p, next, nil).first
	}
	m := f.mark()
	res := f.placePara(p, next, nil)
	split := len(res.segs) > 1
	switch {
	case !split || res.broken:
	case p.pp.keepLines && !m.atTop:
		f.restore(m)
		f.nextColumn(false)
		res = f.placePara(p, next, nil)
	case p.pp.widowCtl:
		n := len(res.segs)
		switch {
		case res.segs[0] == 1 && !m.atTop:
			// an orphan: the first line alone at the bottom
			f.restore(m)
			f.nextColumn(false)
			res = f.placePara(p, next, nil)
		case res.segs[n-1] == 1 && res.segs[n-2] >= 3:
			// a widow: the last line alone at the top; take one more
			total := 0
			for _, s := range res.segs[:n-1] {
				total += s
			}
			f.restore(m)
			res = f.placePara(p, next, map[int]bool{total - 1: true})
		case res.segs[n-1] == 1 && res.segs[n-2] == 2 && n == 2 && !m.atTop:
			f.restore(m)
			f.nextColumn(false)
			res = f.placePara(p, next, nil)
		}
	}
	return res.first
}

// placePara places the lines of a paragraph, breaking the column before
// the lines in force.
func (f *flow) placePara(p *para, next block, force map[int]bool) placed {
	var res placed
	if p.pp.pageBreakBefore && f.paged && !(f.atTop && f.col == 0) {
		f.newPage(true)
	}
	before, after := f.spacing(p)
	if p.pp.contextual && f.prevPara != nil && f.prevPara.style == p.style {
		before = 0
		f.prevAfter = 0
	}
	if f.c.suppressSpBfAfterPgBrk && f.atTop && f.hard {
		before = 0
	}
	gap := f.prevAfter + before
	if f.atTop && f.paged {
		gap = before
	}
	bstart := len(*f.out())
	bdrTop, bdrBottom := f.c.paraBorderSides(p, f.prevPara, next)
	boxTop := 0.0 // room of the top border above the first line
	if bdrTop != nil {
		boxTop = bdrTop.space + bdrTop.width
	}
	start := 0
	first := true
	lineNo := 0
	var segStart float64
	var segs []paraSeg
	y := f.y + gap + boxTop
	cur := position{-1, -1}
	var prev *line
	for first || start < len(p.items) {
		if force[lineNo] && !f.atTop {
			segs = f.closeSeg(segs, segStart)
			f.nextColumn(false)
			y = f.y
		}
		var am anchorMark
		if first {
			am = f.anchorMark()
			f.anchors(p, y)
		}
		ln, ly := f.fitLine(p, start, first, y)
		if f.paged && ly+ln.height > f.limit()-f.lineNotes(ln) && !f.atTop && !ln.breakOnly() {
			segs = f.closeSeg(segs, segStart)
			if first {
				// the paragraph's floating objects go with it
				f.undoAnchors(am)
			}
			f.nextColumn(false)
			y = f.y
			if first {
				y += before + boxTop
			}
			continue
		}
		if pos := f.pos(); pos != cur {
			if cur.page >= 0 {
				segs = f.closeSeg(segs, segStart)
			}
			cur = pos
			segStart = ly
			if first {
				segStart -= boxTop
			}
			res.segs = append(res.segs, 0)
			if first {
				res.first = pos
			}
		}
		res.segs[len(res.segs)-1]++
		if prev != nil {
			ln.sep = joinMark(prev, ln)
		}
		prev = ln
		f.placeNotes(ln)
		f.emitLine(ln, ly, first)
		if first {
			f.c.recordBookmarks(f, p, ly)
		}
		f.y = ly + ln.height
		f.atTop = false
		y = f.y
		start += len(ln.items)
		first = false
		lineNo++
		if ln.after != 0 {
			res.broken = true
		}
		switch ln.after {
		case kPage:
			segs = f.closeSeg(segs, segStart)
			if f.paged {
				f.newPage(true)
			}
			y = f.y
			cur = position{-1, -1}
		case kColumn:
			segs = f.closeSeg(segs, segStart)
			f.nextColumn(true)
			y = f.y
			cur = position{-1, -1}
		}
	}
	if bdrBottom != nil {
		f.y += bdrBottom.space + bdrBottom.width
	}
	segs = f.closeSeg(segs, segStart)
	f.paraBackground(p, segs, bstart, bdrTop, bdrBottom)
	res.last = f.pos()
	f.prevPara = p
	f.prevAfter = after
	if p.breakOnly {
		f.prevAfter = 0
	}
	return res
}

// breakOnly reports whether a line holds nothing but a page or column
// break: it stays at the bottom of the column even when it does not fit, so
// that the break does not leave an empty page.
func (ln *line) breakOnly() bool {
	for _, it := range ln.items {
		if it.kind != kPage && it.kind != kColumn {
			return false
		}
	}
	return ln.after != 0
}

// joinMark returns how a line joins the one before it in the text: a
// paragraph break after a line break, nothing between East Asian
// characters, else a space.
func joinMark(prev, ln *line) byte {
	if prev.endsBreak {
		return 'P'
	}
	var a, b rune
	for i := len(prev.items) - 1; i >= 0; i-- {
		if it := prev.items[i]; it.kind == kChar {
			a = it.r
			break
		}
	}
	for _, it := range ln.items {
		if it.kind == kChar {
			b = it.r
			break
		}
	}
	if a != 0 && b != 0 && linebreak.Joins(a, b) {
		return 'W'
	}
	return 'L'
}

// paraSeg is the part of a paragraph in one column, for its shading and
// borders.
type paraSeg struct {
	x0, w, y0, y1 float64
	page          *page
}

func (f *flow) closeSeg(segs []paraSeg, y0 float64) []paraSeg {
	s := paraSeg{x0: f.x0, w: f.w, y0: y0, y1: f.y, page: f.pg}
	if n := len(segs); n > 0 && segs[n-1].page == f.pg && segs[n-1].x0 == f.x0 {
		segs[n-1].y1 = f.y
		return segs
	}
	if s.y1 <= s.y0 {
		return segs
	}
	return append(segs, s)
}

// fitLine breaks the line of p that starts at item start at or below y,
// beside the floating objects on the page. It returns the line and its
// top.
func (f *flow) fitLine(p *para, start int, first bool, y float64) (*line, float64) {
	indL, indR, _ := p.indents()
	left0, right0 := indL, f.w-indR
	if right0 < left0+1 {
		right0 = left0 + 1
	}
	guess := 1.2 * p.mark.size
	for range 50 {
		l, r, until := f.avail(y, guess, left0, right0)
		if r-l < math.Min(36, right0-left0) && until > y {
			y = until
			continue
		}
		ln := f.lc.breakLine(p, start, first, l, r)
		if ln.height > guess+0.01 {
			l2, r2, until2 := f.avail(y, ln.height, left0, right0)
			if r2-l2 < math.Min(36, right0-left0) && until2 > y {
				y = until2
				continue
			}
			if l2 != l || r2 != r {
				ln = f.lc.breakLine(p, start, first, l2, r2)
			}
		}
		if clears(ln) {
			if u := f.floatBottom(y); u > y {
				y = u
			}
		}
		return ln, y
	}
	return f.lc.breakLine(p, start, first, left0, right0), y
}

// clears reports whether a line ends with a break that moves the next
// line below the floating objects.
func clears(ln *line) bool {
	for _, it := range ln.items {
		if it.kind == kBreak && it.clear {
			return true
		}
	}
	return false
}

// avail narrows the line area [left, right] (column coordinates) to the
// side of the floating objects beside a line at y of height h. until is
// the bottom of the lowest object that took the space.
func (f *flow) avail(y, h, left, right float64) (l, r, until float64) {
	l, r = left, right
	for _, e := range f.floats {
		if e.y1 <= y || e.y0 >= y+h {
			continue
		}
		until = math.Max(until, e.y1)
		if e.full {
			return l, l, until
		}
		x0, x1 := e.x0-f.x0, e.x1-f.x0
		if x1 <= l || x0 >= r {
			continue
		}
		leftW, rightW := x0-l, r-x1
		switch e.side {
		case "left":
			r = math.Min(r, x0)
		case "right":
			l = math.Max(l, x1)
		default:
			if leftW >= rightW {
				r = math.Min(r, x0)
			} else {
				l = math.Max(l, x1)
			}
		}
	}
	return l, r, until
}

func (f *flow) floatBottom(y float64) float64 {
	b := y
	for _, e := range f.floats {
		if e.y0 <= y+0.01 && e.y1 > b {
			b = e.y1
		}
	}
	return b
}

// emitLine adds the op that draws a line with its top at y.
func (f *flow) emitLine(ln *line, y float64, first bool) {
	p := ln.p
	x := f.x0
	o := op{y0: y, y1: y + ln.height, text: true, ctx: p.ctx}
	switch {
	case first && p.heading > 0:
		o.leaf, o.level = 'H', p.heading
	case first:
		o.leaf = 'P'
	default:
		o.leaf = ln.sep
	}
	o.fn = func(e *emitter, dx, dy float64) { e.line(ln, x+dx, y+dy) }
	f.emit(o)
}

// anchors places the floating objects of a paragraph whose first line is
// at y.
func (f *flow) anchors(p *para, y float64) {
	if f.paged && f.sec.vertical {
		f.verticalAnchors(p)
		return
	}
	for _, o := range p.anchors {
		x, top := f.floatPos(o, p, y)
		box := [4]float64{x - o.ext[0], top - o.ext[1], x + o.w + o.ext[2], top + o.h + o.ext[3]}
		obj := o
		dop := op{y0: top - o.ext[1], y1: top + o.h + o.ext[3], top: true}
		dop.fn = func(e *emitter, dx, dy float64) {
			e.drawObject(&obj.inlineObj, boxAt(x+dx, top+dy, obj.w, obj.h))
		}
		switch {
		case o.behind:
			if f.paged {
				f.pg.behind = append(f.pg.behind, dop)
			} else {
				f.behind = append(f.behind, dop)
			}
		default:
			if f.paged {
				f.pg.front = append(f.pg.front, dop)
			} else {
				f.front = append(f.front, dop)
			}
		}
		switch o.wrap {
		case "square":
			f.floats = append(f.floats, excl{x0: box[0] - o.dist[2], y0: box[1] - o.dist[0], x1: box[2] + o.dist[3], y1: box[3] + o.dist[1], side: o.side})
		case "topAndBottom":
			f.floats = append(f.floats, excl{x0: box[0], y0: box[1] - o.dist[0], x1: box[2], y1: box[3] + o.dist[1], full: true})
		}
	}
}

// verticalAnchors places the floating objects of a paragraph on a page of
// vertical text: upright, in page coordinates, relative to the page or
// the margins (positions relative to the paragraph or column are taken as
// relative to the margins). The area they keep text out of is turned into
// the flow's coordinates.
func (f *flow) verticalAnchors(p *para) {
	pg, s := f.pg, f.sec
	for _, o := range p.anchors {
		bx, bw := s.left, s.pgW-s.left-s.right
		switch o.hRel {
		case "page":
			bx, bw = 0, s.pgW
		case "leftMargin", "insideMargin":
			bx, bw = 0, s.left
		case "rightMargin", "outsideMargin":
			bx, bw = s.pgW-s.right, s.right
		}
		by, bh := pg.top, pg.bottom-pg.top
		switch o.vRel {
		case "page":
			by, bh = 0, s.pgH
		case "topMargin", "insideMargin":
			by, bh = 0, pg.top
		case "bottomMargin", "outsideMargin":
			by, bh = pg.bottom, s.pgH-pg.bottom
		}
		x, y := alignIn(bx, bw, o.w, o.hOff, o.hAlign), alignIn(by, bh, o.h, o.vOff, o.vAlign)
		obj := o
		dop := op{top: true, fn: func(e *emitter, dx, dy float64) {
			e.drawObject(&obj.inlineObj, boxAt(x+dx, y+dy, obj.w, obj.h))
		}}
		if o.behind {
			pg.behind = append(pg.behind, dop)
		} else {
			pg.front = append(pg.front, dop)
		}
		// page → flow: x runs down from the top of the body, y leftwards
		// from its right edge
		right := s.pgW - s.right
		x0, x1 := x-o.ext[0]-o.dist[2], x+o.w+o.ext[2]+o.dist[3]
		y0, y1 := y-o.ext[1]-o.dist[0], y+o.h+o.ext[3]+o.dist[1]
		e := excl{x0: y0 - pg.top, x1: y1 - pg.top, y0: right - x1, y1: right - x0, side: o.side}
		switch o.wrap {
		case "square":
			f.floats = append(f.floats, e)
		case "topAndBottom":
			e.full = true
			f.floats = append(f.floats, e)
		}
	}
}

// alignIn positions an extent of size in [base, base+span]: by an
// alignment, or at an offset from base.
func alignIn(base, span, size, off float64, align string) float64 {
	switch align {
	case "left", "top", "inside":
		return base
	case "center":
		return base + (span-size)/2
	case "right", "bottom", "outside":
		return base + span - size
	}
	return base + off
}

// anchorMark records the floating objects placed so far, to take back
// those of a paragraph that moves on.
type anchorMark struct {
	pg                                      *page
	nfloats, nbehind, nfront, nbehU, nfronU int
}

func (f *flow) anchorMark() anchorMark {
	m := anchorMark{pg: f.pg, nfloats: len(f.floats), nbehU: len(f.behind), nfronU: len(f.front)}
	if f.pg != nil {
		m.nbehind, m.nfront = len(f.pg.behind), len(f.pg.front)
	}
	return m
}

func (f *flow) undoAnchors(m anchorMark) {
	if m.pg != nil {
		m.pg.behind, m.pg.front = m.pg.behind[:m.nbehind], m.pg.front[:m.nfront]
	}
	if m.pg == f.pg {
		f.floats = f.floats[:m.nfloats]
	}
	f.behind, f.front = f.behind[:m.nbehU], f.front[:m.nfronU]
}

// floatPos returns the top left corner of a floating object anchored in
// a paragraph whose first line is at y. Where there is no page (the scroll
// view, cells, text boxes), positions relative to the page are relative to
// the column and the paragraph, and the object stays in the column.
func (f *flow) floatPos(o *floatObj, p *para, y float64) (float64, float64) {
	sec := f.sec
	onPage := sec != nil && (f.paged || f.origin != nil)
	var pageX, pageY, pageW, marL, marR float64
	switch {
	case onPage:
		pageW, marL, marR = sec.pgW, sec.left, sec.right
		if f.origin != nil {
			pageX, pageY = -f.origin[0], -f.origin[1]
		}
	case sec != nil:
		pageX, pageW, marL, marR = f.x0-sec.left, sec.pgW, sec.left, sec.right
	default:
		pageX, pageW = f.x0, f.w
	}
	bx, bw := f.x0, f.w
	switch o.hRel {
	case "page":
		bx, bw = pageX, pageW
	case "margin":
		bx, bw = pageX+marL, pageW-marL-marR
	case "leftMargin", "insideMargin":
		bx, bw = pageX, marL
	case "rightMargin", "outsideMargin":
		bx, bw = pageX+pageW-marR, marR
	case "character":
		indL, _, _ := p.indents()
		bx, bw = f.x0+indL, f.w-indL
	}
	x := bx + o.hOff
	switch o.hAlign {
	case "left", "inside":
		x = bx
	case "center":
		x = bx + (bw-o.w)/2
	case "right", "outside":
		x = bx + bw - o.w
	}
	if !onPage {
		x = math.Max(f.x0, math.Min(x, f.x0+f.w-o.w))
	}
	by, bh := y, 0.0
	if onPage {
		top, bottom := sec.top, sec.pgH-sec.bottom
		if f.pg != nil {
			top, bottom = f.pg.top, f.pg.bottom
		}
		switch o.vRel {
		case "page":
			by, bh = pageY, sec.pgH
		case "margin":
			by, bh = pageY+top, bottom-top
		case "topMargin", "insideMargin":
			by, bh = pageY, top
		case "bottomMargin", "outsideMargin":
			by, bh = pageY+bottom, sec.pgH-bottom
		}
	}
	top := by + o.vOff
	switch o.vAlign {
	case "top", "inside":
		top = by
	case "center":
		top = by + (bh-o.h)/2
	case "bottom", "outside":
		top = by + bh - o.h
	}
	return x, top
}

// lineNotes returns the height the footnotes referenced by a line need on
// the page that are not there yet.
func (f *flow) lineNotes(ln *line) float64 {
	if !f.paged {
		return 0
	}
	h := 0.0
	for _, it := range ln.items {
		if it.note == nil || it.note.kind != "f" {
			continue
		}
		if len(f.pg.notes) == 0 && h == 0 {
			h += noteSeparator
		}
		h += f.c.noteHeight(it.note, f.sec)
	}
	return h
}

// noteSeparator is the room of the line above the footnotes.
const noteSeparator = 12

func (f *flow) placeNotes(ln *line) {
	for _, it := range ln.items {
		if it.note == nil || it.note.kind != "f" {
			continue
		}
		if !f.paged {
			continue // the scroll view draws them at the end
		}
		if len(f.pg.notes) == 0 {
			f.pg.noteH += noteSeparator
		}
		h := f.c.noteHeight(it.note, f.sec)
		f.pg.notes = append(f.pg.notes, &noteUse{n: it.note, h: h})
		f.pg.noteH += h
	}
}
