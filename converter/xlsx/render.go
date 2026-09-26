package xlsx

import (
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
)

// A worksheet is drawn into square tiles (spec §4.1): per tile, the cell
// fills, conditional data bars and icons, the cell text, the borders, and
// last the drawings. Cells whose painting crosses a tile edge (overflowing
// or merged text, a border) are drawn in every tile they touch with that
// tile's origin; the viewer clips each tile, and text extraction keeps a
// run only in the tile its anchor lies in. Pictures, shapes and charts are
// objects of their own that the tiles they cover use.

const tileSize = 2048

// view extent beyond the last used cell: some empty grid, as in Excel
const (
	marginRows = 5
	marginCols = 2
	minRows    = 30
	minCols    = 12
)

type sheetCtx struct {
	c      *converter
	ws     *worksheet
	cols   *axis
	rows   *axis
	nRows  int
	nCols  int
	tiles  map[[2]int]*tileCv
	merges []cellRange
	// mergeRows lists the merges that cover a row (within the view)
	mergeRows map[int][]int
	tables    []*table
	cf        *cfResults
	comments  map[[2]int]bool
	drawings  []*anchored
	base      map[int]*cellFmt // formats by xf index
	textCols  int              // columns that text reaches (overflowing text)
	formulas  map[string]fnode // parsed conditional format formulas (nil: unsupported)
}

// tileCv is a tile's canvas and the drawing state emitted into it.
type tileCv struct {
	cv     *canvas.Canvas
	ox, oy float64
	st     emitState
	stack  []emitState
}

type emitState struct {
	fill    bdf.Color
	hasFill bool
	font    bdf.FontRef
	size    float64
	hasFont bool
	ls      float64 // letter spacing
	dashed  bool    // a line dash is set
}

func (t *tileCv) save() {
	t.cv.Obj.Save()
	t.stack = append(t.stack, t.st)
}

func (t *tileCv) restore() {
	t.cv.Obj.Restore()
	if n := len(t.stack); n > 0 {
		t.st = t.stack[n-1]
		t.stack = t.stack[:n-1]
	}
}

func (t *tileCv) fillColor(c bdf.Color) {
	if !t.st.hasFill || t.st.fill != c {
		t.cv.Obj.FillColor(c)
		t.st.fill, t.st.hasFill = c, true
	}
}

func (t *tileCv) setFont(fc *fontset.Choice, size float64) {
	ref := t.cv.Font(fc.Use)
	if !t.st.hasFont || t.st.font != ref || t.st.size != size {
		t.cv.Obj.Font(ref, f32(size))
		t.st.font, t.st.size, t.st.hasFont = ref, size, true
	}
}

func (t *tileCv) letterSpacing(ls float64) {
	if ls != t.st.ls {
		t.cv.Obj.TextStyle(bdf.AlignLeft, bdf.BaselineAlphabetic, bdf.DirInherit, f32(ls))
		t.st.ls = ls
	}
}

func (t *tileCv) rect(x, y, w, h float64) {
	t.cv.Obj.FillRect(f32(x-t.ox), f32(y-t.oy), f32(w), f32(h))
	t.cv.Drawn = true
}

func f32(v float64) float32 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return float32(v)
}

// tilesIn returns the tiles a rectangle (sheet coordinates) touches,
// creating them.
func (s *sheetCtx) tilesIn(x0, y0, x1, y1 float64) []*tileCv {
	if x1 <= x0 || y1 <= y0 {
		return nil
	}
	w, h := s.cols.total(), s.rows.total()
	x0, y0 = math.Max(x0, 0), math.Max(y0, 0)
	x1, y1 = math.Min(x1, w), math.Min(y1, h)
	if x1 <= x0 || y1 <= y0 {
		return nil
	}
	var out []*tileCv
	for ty := int(y0 / tileSize); float64(ty)*tileSize < y1; ty++ {
		for tx := int(x0 / tileSize); float64(tx)*tileSize < x1; tx++ {
			out = append(out, s.tileAt(tx, ty))
		}
	}
	return out
}

func (s *sheetCtx) tileAt(tx, ty int) *tileCv {
	k := [2]int{tx, ty}
	if t, ok := s.tiles[k]; ok {
		return t
	}
	t := &tileCv{cv: s.c.cvs.New(), ox: float64(tx) * tileSize, oy: float64(ty) * tileSize}
	t.cv.Obj.SetBBox(0, 0, tileSize, tileSize)
	s.tiles[k] = t
	return t
}

// worksheetSafe renders a worksheet; a malformed one that trips the
// renderer becomes an empty sheet and a warning.
func (c *converter) worksheetSafe(ref sheetRef, id string) (v *bdf.View, tiles map[string]*canvas.Canvas, err error) {
	defer func() {
		if r := recover(); r != nil {
			c.warnf("sheet %q: internal error: %v", ref.name, r)
			v = c.doc.NewView(id, bdf.ViewSheet, ref.name)
			v.Tile = tileSize
			v.Cols = []bdf.Run{{minCols, 48}}
			v.Rows = []bdf.Run{{minRows, 15}}
			v.Gridlines = true
			v.Tiles = map[string]string{}
			tiles, err = nil, nil
		}
	}()
	return c.worksheet(ref, id)
}

func (c *converter) worksheet(ref sheetRef, id string) (*bdf.View, map[string]*canvas.Canvas, error) {
	ws := c.parsed[ref.name]
	if ws == nil {
		var err error
		if ws, err = c.readWorksheet(ref.part); err != nil {
			return nil, nil, err
		}
		ws.name = ref.name
		ws.sortRows()
		c.parsed[ref.name] = ws
	}
	if ws.noValue > 0 {
		c.warnf("sheet %q: %d formula cell(s) have no calculated value and are left empty (the workbook was saved without calculating them)", ws.name, ws.noValue)
	}
	if ws.rtl {
		c.warnOnce("rtl", "sheet %q: right-to-left sheets are drawn left to right", ws.name)
	}
	s := &sheetCtx{c: c, ws: ws, tiles: map[[2]int]*tileCv{}, base: map[int]*cellFmt{}, formulas: map[string]fnode{}}
	s.loadTables()
	s.loadDrawings()
	s.loadComments()
	s.extent()
	s.indexMerges()
	s.geometry()
	s.cf = s.evalConditionalFormats()

	v := c.doc.NewView(id, bdf.ViewSheet, ws.name)
	v.Tile = tileSize
	v.Gridlines = ws.showGrid
	if ws.freezeC > 0 || ws.freezeR > 0 {
		v.Freeze = &bdf.Freeze{Cols: min(ws.freezeC, s.nCols), Rows: min(ws.freezeR, s.nRows)}
	}
	v.Tiles = map[string]string{}

	s.paintFills()
	s.paintConditional()
	s.paintText()
	s.paintBorders()
	s.paintMarks()
	s.paintLinks()
	s.paintDrawings()
	// the view ends after the text that flows past its last column
	s.nCols = min(max(s.nCols, s.textCols+marginCols), maxCols)
	s.cols.n = s.nCols
	v.Cols = s.cols.runs()
	v.Rows = s.rows.runs()

	out := map[string]*canvas.Canvas{}
	for k, t := range s.tiles {
		if t.cv.Drawn {
			out[strconv.Itoa(k[0])+","+strconv.Itoa(k[1])] = t.cv
		}
	}
	return v, out, nil
}

// extent sets the rows and columns of the view: what the sheet uses, and
// some empty grid around it.
func (s *sheetCtx) extent() {
	ws := s.ws
	maxR, maxC := -1, -1
	use := func(r, c int) {
		maxR, maxC = max(maxR, r), max(maxC, c)
	}
	for _, rw := range ws.rows {
		for _, cl := range rw.cells {
			if cl.kind != cellBlank || cl.style != 0 {
				use(rw.idx, cl.col)
			}
		}
		if rw.ht >= 0 || rw.style >= 0 {
			maxR = max(maxR, rw.idx)
		}
	}
	for _, m := range ws.merges {
		use(min(m.r1, maxRows-1), min(m.c1, maxCols-1))
	}
	for _, t := range s.tables {
		use(t.ref.r1, t.ref.c1)
	}
	for _, d := range s.drawings {
		use(d.toRow, d.toCol)
	}
	for k := range s.comments {
		use(k[0], k[1])
	}
	s.nRows = min(max(maxR+1+marginRows, minRows), maxRows)
	s.nCols = min(max(maxC+1+marginCols, minCols), maxCols)
	if ws.freezeR > 0 {
		s.nRows = min(max(s.nRows, ws.freezeR+marginRows), maxRows)
	}
	if ws.freezeC > 0 {
		s.nCols = min(max(s.nCols, ws.freezeC+marginCols), maxCols)
	}
}

// geometry sizes the columns, then the rows (whose automatic heights
// depend on the text laid out in the columns).
func (s *sheetCtx) geometry() {
	ws, c := s.ws, s.c
	defCol := ws.defaultColPt(c.mdw)
	cw := map[int]float64{}
	// the columns of the whole sheet: text may flow past the view (which then grows)
	for _, cd := range ws.cols {
		for i := cd.min; i <= cd.max && i < maxCols; i++ {
			switch {
			case cd.hidden:
				cw[i] = 0
			case cd.width >= 0:
				cw[i] = colWidthPt(cd.width, c.mdw)
			default:
				cw[i] = defCol
			}
		}
	}
	s.cols = newAxis(maxCols, defCol, cw)

	defRow := ws.defRowH
	if defRow <= 0 {
		defRow = autoRowHeight(c.fontLineHeight(c.st.font(0)))
	}
	visibleDef := defRow
	if ws.zeroH {
		visibleDef = 0
	}
	defLine := c.fontLineHeight(c.st.font(0))
	rh := map[int]float64{}
	for i := range ws.rows {
		rw := &ws.rows[i]
		if rw.idx >= s.nRows {
			break
		}
		switch {
		case rw.hidden:
			rh[rw.idx] = 0
		case rw.ht >= 0:
			rh[rw.idx] = rw.ht
		default:
			h := defRow
			if fit := s.autoHeight(rw, defLine); fit > h {
				h = fit
			}
			rh[rw.idx] = h
		}
	}
	s.rows = newAxis(s.nRows, visibleDef, rh)
}

// autoHeight is the height a row without one of its own needs for its
// text (merged cells do not count, as in Excel). Only cells that may be
// taller than a line of the default font are laid out.
func (s *sheetCtx) autoHeight(rw *row, defLine float64) float64 {
	best := 0.0
	for i := range rw.cells {
		cl := &rw.cells[i]
		if cl.kind == cellBlank {
			continue
		}
		if _, merged := s.mergeAt(rw.idx, cl.col); merged {
			continue
		}
		f := s.fmtOf(cl.style)
		if !f.align.wrap && f.align.rotation == 0 && (cl.text == nil || cl.text.runs == nil) && s.c.fontLineHeight(&f.font) <= defLine+0.01 {
			continue
		}
		lay := s.layoutCell(rw.idx, cl.col, cl, f, box{0, 0, s.cols.at(cl.col), 1e6}, true)
		if lay != nil && lay.textH > best {
			best = lay.textH
		}
	}
	if best == 0 {
		return 0
	}
	return autoRowHeight(best)
}

func (s *sheetCtx) indexMerges() {
	s.mergeRows = map[int][]int{}
	for _, m := range s.ws.merges {
		if m.r0 >= s.nRows || m.c0 >= s.nCols {
			continue
		}
		m.r1, m.c1 = min(m.r1, s.nRows-1), min(m.c1, s.nCols-1)
		i := len(s.merges)
		s.merges = append(s.merges, m)
		for r := m.r0; r <= m.r1; r++ {
			s.mergeRows[r] = append(s.mergeRows[r], i)
		}
	}
}

// mergeAt returns the merged range a cell is in.
func (s *sheetCtx) mergeAt(r, c int) (cellRange, bool) {
	for _, i := range s.mergeRows[r] {
		if m := s.merges[i]; c >= m.c0 && c <= m.c1 {
			return m, true
		}
	}
	return cellRange{}, false
}

// box is a rectangle in sheet coordinates.
type box struct{ x, y, w, h float64 }

func (s *sheetCtx) cellBox(r, c int) box {
	return box{s.cols.pos(c), s.rows.pos(r), s.cols.at(c), s.rows.at(r)}
}

func (s *sheetCtx) rangeBox(m cellRange) box {
	x0, y0 := s.cols.pos(m.c0), s.rows.pos(m.r0)
	return box{x0, y0, s.cols.pos(m.c1+1) - x0, s.rows.pos(m.r1+1) - y0}
}

// cellFmt is the resolved format of a cell.
type cellFmt struct {
	font   xfont
	fill   xfill
	border xborder
	align  xalign
	numFmt *numFormat
	fontID int
	// hideText draws the value invisible (still selectable and searchable):
	// data bars and icon sets that do not show the value
	hideText bool
}

// fmtOf returns the format of an xf index (shared; do not modify).
func (s *sheetCtx) fmtOf(i int) *cellFmt {
	if f, ok := s.base[i]; ok {
		return f
	}
	st := s.c.st
	x := st.xf(i)
	f := &cellFmt{font: *st.font(x.font), fill: st.fill(x.fill), border: *st.border(x.border), align: x.align,
		numFmt: st.numFormat(x.numFmt), fontID: x.font}
	s.base[i] = f
	return f
}

// styleAt returns the xf index of a position: its cell's, else the row's,
// else the column's; -1 for none.
func (s *sheetCtx) styleAt(rw *row, cl *cell, c int) int {
	if cl != nil {
		return cl.style
	}
	if rw != nil && rw.style >= 0 {
		return rw.style
	}
	return s.ws.colStyle(c)
}

// formatAt returns the format of a position with the table style and the
// conditional formats over it (nil when nothing formats it).
func (s *sheetCtx) formatAt(r, c int, rw *row, cl *cell) *cellFmt {
	idx := s.styleAt(rw, cl, c)
	t := s.tableAt(r, c)
	cfr := s.cf.at(r, c)
	if idx < 0 && t == nil && cfr == nil {
		return nil
	}
	f := s.fmtOf(max(idx, 0))
	if t == nil && (cfr == nil || !cfr.hasDxf()) {
		return f
	}
	cp := *f
	if t != nil {
		t.apply(s, &cp, r, c)
	}
	if cfr != nil {
		cfr.apply(&cp)
	}
	return &cp
}

// paintFills draws the cell backgrounds: runs of equal fills along each
// row, joined down the rows while they repeat, and merged ranges whole.
func (s *sheetCtx) paintFills() {
	type open struct {
		x0, x1, y0, y1 float64
		key            string
		f              xfill
	}
	var prev []open
	flush := func(o open) {
		s.fillBox(box{o.x0, o.y0, o.x1 - o.x0, o.y1 - o.y0}, o.f)
	}
	for r := 0; r < s.nRows; r++ {
		h := s.rows.at(r)
		if h == 0 {
			continue
		}
		y0 := s.rows.pos(r)
		var cur []open
		s.eachFormatted(r, func(c0, c1 int, f *cellFmt) {
			if f == nil || f.fill.empty() {
				return
			}
			x0, x1 := s.cols.pos(c0), s.cols.pos(c1+1)
			if x1 <= x0 {
				return
			}
			key := fillKey(f.fill)
			if key == "" { // gradients are per cell
				for c := c0; c <= c1; c++ {
					if w := s.cols.at(c); w > 0 {
						s.fillBox(box{s.cols.pos(c), y0, w, h}, f.fill)
					}
				}
				return
			}
			if n := len(cur); n > 0 && cur[n-1].key == key && cur[n-1].x1 == x0 {
				cur[n-1].x1 = x1
				return
			}
			cur = append(cur, open{x0: x0, x1: x1, y0: y0, y1: y0 + h, key: key, f: f.fill})
		})
		// continue the rectangles of the row above that match exactly
		var next []open
		j := 0
		for _, o := range cur {
			for j < len(prev) && prev[j].x0 < o.x0 {
				flush(prev[j])
				j++
			}
			if j < len(prev) && prev[j].x0 == o.x0 && prev[j].x1 == o.x1 && prev[j].key == o.key && prev[j].y1 == o.y0 {
				o.y0 = prev[j].y0
				j++
			}
			next = append(next, o)
		}
		for ; j < len(prev); j++ {
			flush(prev[j])
		}
		prev = next
	}
	for _, o := range prev {
		flush(o)
	}
	// merged ranges take the fill of their top-left cell
	for _, m := range s.merges {
		rw := s.ws.rowAt(m.r0)
		f := s.formatAt(m.r0, m.c0, rw, s.ws.cellAt(m.r0, m.c0))
		if f != nil && !f.fill.empty() {
			s.fillBox(s.rangeBox(m), f.fill)
		}
	}
}

// eachFormatted calls fn for the formatted positions of a row outside
// merged ranges, as runs of columns with the same format (a nil format for
// columns that nothing formats is skipped).
func (s *sheetCtx) eachFormatted(r int, fn func(c0, c1 int, f *cellFmt)) {
	rw := s.ws.rowAt(r)
	// columns with a format of their own: cells, tables, conditional formats
	special := map[int]bool{}
	if rw != nil {
		for _, cl := range rw.cells {
			if cl.col < s.nCols {
				special[cl.col] = true
			}
		}
	}
	for _, t := range s.tables {
		if r >= t.ref.r0 && r <= t.ref.r1 {
			for c := t.ref.c0; c <= t.ref.c1 && c < s.nCols; c++ {
				special[c] = true
			}
		}
	}
	for _, c := range s.cf.colsIn(r) {
		if c < s.nCols {
			special[c] = true
		}
	}
	cols := make([]int, 0, len(special))
	for c := range special {
		cols = append(cols, c)
	}
	sort.Ints(cols)
	merged := func(c int) bool {
		_, ok := s.mergeAt(r, c)
		return ok
	}
	// base formats: the row's, or the columns'
	type iv struct{ c0, c1, style int }
	var base []iv
	if rw != nil && rw.style >= 0 {
		base = []iv{{0, s.nCols - 1, rw.style}}
	} else {
		for _, cd := range s.ws.cols {
			if cd.style >= 0 && cd.min < s.nCols {
				base = append(base, iv{cd.min, min(cd.max, s.nCols-1), cd.style})
			}
		}
	}
	runStart, runEnd := -1, -1
	var runFmt *cellFmt
	emit := func() {
		if runStart >= 0 && runFmt != nil {
			fn(runStart, runEnd, runFmt)
		}
		runStart, runFmt = -1, nil
	}
	put := func(c int, f *cellFmt) {
		if f == nil || merged(c) || s.cols.at(c) == 0 {
			emit()
			return
		}
		if runStart >= 0 && runEnd == c-1 && runFmt == f {
			runEnd = c
			return
		}
		emit()
		runStart, runEnd, runFmt = c, c, f
	}
	k := 0
	baseAt := func(c int) int {
		style := -1
		for _, b := range base {
			if c >= b.c0 && c <= b.c1 {
				style = b.style
			}
		}
		return style
	}
	// walk the columns that are formatted: special ones and base intervals
	c := 0
	for c < s.nCols {
		if k < len(cols) && cols[k] == c {
			var cl *cell
			if rw != nil {
				cl = s.ws.cellAt(r, c)
			}
			put(c, s.formatAt(r, c, rw, cl))
			k++
			c++
			continue
		}
		next := s.nCols
		if k < len(cols) {
			next = cols[k]
		}
		// the stretch c..next-1 has only base formats
		for c < next {
			style := baseAt(c)
			if style < 0 {
				emit()
				// jump to the next base interval start
				jump := next
				for _, b := range base {
					if b.c0 > c && b.c0 < jump {
						jump = b.c0
					}
				}
				c = jump
				continue
			}
			put(c, s.fmtOf(style))
			c++
		}
	}
	emit()
}

// fillKey identifies a fill for joining cells; "" for gradients.
func fillKey(f xfill) string {
	if f.grad != nil {
		return ""
	}
	return fmt.Sprintf("%s|%v|%v", f.pattern, f.fg, f.bg)
}

// fillBox paints a fill over a box.
func (s *sheetCtx) fillBox(b box, f xfill) {
	st := s.c.st
	for _, t := range s.tilesIn(b.x, b.y, b.x+b.w, b.y+b.h) {
		switch {
		case f.grad != nil:
			p := s.gradientPaint(f.grad, b)
			t.cv.Obj.FillPaint(t.cv.Obj.AddPaint(p))
			t.st.hasFill = false
			// gradient coordinates are in sheet space: move them to the tile
			t.save()
			t.cv.Obj.Translate(f32(-t.ox), f32(-t.oy))
			t.cv.Obj.FillRect(f32(b.x), f32(b.y), f32(b.w), f32(b.h))
			t.restore()
			t.cv.Drawn = true
		case f.pattern == "solid":
			t.fillColor(st.color(f.fg, black).bdf())
			t.rect(b.x, b.y, b.w, b.h)
		default:
			fg := st.color(f.fg, black)
			bg := st.color(f.bg, white)
			img := s.c.patternImage(f.pattern, fg, bg)
			ref := t.cv.Image(img)
			p := bdf.Pattern(ref, bdf.RepeatBoth)
			// one pattern pixel per screen pixel, aligned to the sheet
			p.Matrix = [6]float32{pxPt, 0, 0, pxPt, f32(-t.ox), f32(-t.oy)}
			t.cv.Obj.FillPaint(t.cv.Obj.AddPaint(p))
			t.st.hasFill = false
			t.rect(b.x, b.y, b.w, b.h)
		}
	}
}

// gradientPaint makes the paint of a gradient fill over a box: linear at
// an angle, or a path gradient from a rectangle given by fractions.
func (s *sheetCtx) gradientPaint(g *gradFill, b box) bdf.Paint {
	st := s.c.st
	var stops []bdf.Stop
	for _, gs := range g.stops {
		stops = append(stops, bdf.Stop{Offset: f32(clamp01(gs.pos)), Color: st.color(gs.color, white).bdf()})
	}
	if g.path {
		cx := b.x + b.w*(g.left+(1-g.right))/2
		cy := b.y + b.h*(g.top+(1-g.bottom))/2
		r := math.Hypot(b.w, b.h) / 2
		return bdf.RadialGradient(f32(cx), f32(cy), 0, f32(cx), f32(cy), f32(r), stops...)
	}
	a := g.degree * math.Pi / 180
	cx, cy := b.x+b.w/2, b.y+b.h/2
	// the gradient line spans the box in its direction
	dx, dy := math.Cos(a), math.Sin(a)
	half := (math.Abs(dx)*b.w + math.Abs(dy)*b.h) / 2
	return bdf.LinearGradient(f32(cx-dx*half), f32(cy-dy*half), f32(cx+dx*half), f32(cy+dy*half), stops...)
}
