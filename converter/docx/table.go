package docx

import (
	"math"
	"slices"
	"strconv"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// table is a w:tbl.
type table struct {
	style     string
	look      tblLook
	grid      []float64
	jc        string
	ind       float64
	borders   [6]*border // top, left, bottom, right, insideH, insideV
	cellMar   [4]float64 // top, left, bottom, right
	shd       bdf.Color
	hasShd    bool
	rows      []*row
	frame     *frame
	rowBand   int
	colBand   int
	cellSpace float64
}

// table border sides
const (
	sTop = iota
	sLeft
	sBottom
	sRight
	sInsideH
	sInsideV
)

type row struct {
	cells      []*cell
	height     float64
	exact      bool
	cantSplit  bool
	header     bool
	gridBefore int
}

type cell struct {
	blocks  []block
	col     int // first grid column
	span    int
	vmerge  string // "", restart, continue
	rowSpan int
	borders [6]*border
	shd     bdf.Color
	hasShd  bool
	mar     [4]float64
	vAlign  string
	frame   *frame
	ctx     []*frame // the structure of its content: …, the table, the cell
}

func (c *converter) applyTblPr(t *table, n *ooxml.Node) {
	if n == nil {
		return
	}
	if b := n.Child("tblBorders"); b != nil {
		c.readBorders(&t.borders, b)
	}
	if m := n.Child("tblCellMar"); m != nil {
		readMargins(&t.cellMar, m)
	}
	if s := n.Child("shd"); s != nil {
		t.shd, t.hasShd = c.shading(s)
	}
	if j := n.Child("jc"); j != nil {
		t.jc = jcValue(val(j))
	}
	if i := n.Child("tblInd"); i != nil {
		t.ind = twips(i, "w", 0)
	}
	if v := n.Child("tblStyleRowBandSize"); v != nil {
		t.rowBand = int(v.AttrInt("val", 1))
	}
	if v := n.Child("tblStyleColBandSize"); v != nil {
		t.colBand = int(v.AttrInt("val", 1))
	}
	if v := n.Child("tblCellSpacing"); v != nil {
		t.cellSpace = twips(v, "w", 0)
	}
}

func (c *converter) readBorders(dst *[6]*border, n *ooxml.Node) {
	for i, names := range [][]string{{"top"}, {"left", "start"}, {"bottom"}, {"right", "end"}, {"insideH"}, {"insideV"}} {
		for _, name := range names {
			if b := n.Child(name); b != nil {
				dst[i] = c.parseBorder(b)
			}
		}
	}
}

func readMargins(dst *[4]float64, n *ooxml.Node) {
	for i, names := range [][]string{{"top"}, {"left", "start"}, {"bottom"}, {"right", "end"}} {
		for _, name := range names {
			if m := n.Child(name); m != nil {
				dst[i] = twips(m, "w", 0)
			}
		}
	}
}

// table reads a w:tbl.
func (w *walker) table(n *ooxml.Node) *table {
	c := w.c
	tblPr := n.Child("tblPr")
	t := &table{style: val(tblPr.Child("tblStyle")), jc: "left", cellMar: [4]float64{0, 5.4, 0, 5.4},
		frame: &frame{kind: fTable}, rowBand: 1, colBand: 1}
	if t.style == "" || c.st.byID[t.style] == nil {
		t.style = c.st.defTable
	}
	for _, st := range c.st.chain(t.style) {
		c.applyTblPr(t, st.tblPr)
	}
	c.applyTblPr(t, tblPr)
	t.look = parseLook(tblPr.Child("tblLook"))
	for _, g := range n.Child("tblGrid").Children("gridCol") {
		t.grid = append(t.grid, twips(g, "w", 0))
	}
	// rows and cells
	var trs []*ooxml.Node
	var collect func(n *ooxml.Node)
	collect = func(n *ooxml.Node) {
		for _, k := range n.Elements() {
			switch k.Name {
			case "tr":
				trs = append(trs, k)
			case "sdt":
				collect(k.Child("sdtContent"))
			case "customXml", "ins":
				collect(k)
			}
		}
	}
	collect(n)
	type tcRef struct {
		n   *ooxml.Node
		row int
		c   *cell
	}
	var tcs []tcRef
	for ri, tr := range trs {
		trPr := tr.Child("trPr")
		r := &row{}
		for _, st := range c.st.chain(t.style) {
			applyTrPr(r, st.trPr)
		}
		applyTrPr(r, trPr)
		col := r.gridBefore
		var cellNodes []*ooxml.Node
		var collectTc func(n *ooxml.Node)
		collectTc = func(n *ooxml.Node) {
			for _, k := range n.Elements() {
				switch k.Name {
				case "tc":
					cellNodes = append(cellNodes, k)
				case "sdt":
					collectTc(k.Child("sdtContent"))
				case "customXml", "ins":
					collectTc(k)
				}
			}
		}
		collectTc(tr)
		for _, tc := range cellNodes {
			tcPr := tc.Child("tcPr")
			ce := &cell{col: col, span: max(1, int(tcPr.Child("gridSpan").AttrInt("val", 1))), rowSpan: 1, vAlign: "top"}
			if vm := tcPr.Child("vMerge"); vm != nil {
				ce.vmerge = vm.AttrStr("val", "continue")
			}
			if tcPr.Child("hMerge").AttrStr("val", "") == "continue" {
				// legacy horizontal merge: the cell widens the one before it
				if n := len(r.cells); n > 0 {
					r.cells[n-1].span += ce.span
					col += ce.span
					continue
				}
			}
			col += ce.span
			r.cells = append(r.cells, ce)
			tcs = append(tcs, tcRef{tc, ri, ce})
		}
		t.rows = append(t.rows, r)
	}
	if len(t.rows) == 0 {
		return nil
	}
	ncols := 0
	for _, r := range t.rows {
		n := 0
		for _, ce := range r.cells {
			n = max(n, ce.col+ce.span)
		}
		ncols = max(ncols, n)
	}
	for len(t.grid) < ncols {
		t.grid = append(t.grid, 0)
	}
	// vertical merges
	for ri, r := range t.rows {
		for _, ce := range r.cells {
			if ce.vmerge != "restart" {
				continue
			}
			for rj := ri + 1; rj < len(t.rows); rj++ {
				below := t.rows[rj].cellAt(ce.col)
				if below == nil || below.vmerge != "continue" {
					break
				}
				ce.rowSpan++
			}
		}
	}
	// cell properties and content, with the table style's formatting
	for _, ref := range tcs {
		ce, tc := ref.c, ref.n
		r := t.rows[ref.row]
		idx := 0
		for i, x := range r.cells {
			if x == ce {
				idx = i
			}
		}
		tl := c.tableStyleLayers(t.style, t.look, ref.row, idx, len(t.rows), len(r.cells), t.rowBand, t.colBand, r.header)
		ce.mar = t.cellMar
		if tl != nil {
			for _, n := range tl.tcPr {
				c.applyTcPr(ce, n)
			}
		}
		c.applyTcPr(ce, tc.Child("tcPr"))
		name := cellRef(ref.row, ce.col)
		if ce.span > 1 || ce.rowSpan > 1 {
			name += ":" + cellRef(ref.row+ce.rowSpan-1, ce.col+ce.span-1)
		}
		if r.header {
			name += " col"
		}
		ce.frame = &frame{kind: fCell, payload: name}
		ce.ctx = append(append(append([]*frame(nil), w.base...), t.frame), ce.frame)
		if ce.vmerge == "continue" {
			continue
		}
		cw := &walker{c: c, part: w.part, tl: tl, base: ce.ctx,
			fields: w.fields, dark: ce.hasShd && dark(ce.shd) || (!ce.hasShd && t.hasShd && dark(t.shd)), noteNo: w.noteNo}
		ce.blocks = cw.blocks(tc)
		w.fields = cw.fields
	}
	return t
}

func (r *row) cellAt(col int) *cell {
	for _, ce := range r.cells {
		if ce.col == col {
			return ce
		}
	}
	return nil
}

func applyTrPr(r *row, n *ooxml.Node) {
	if n == nil {
		return
	}
	if h := n.Child("trHeight"); h != nil {
		r.height = twips(h, "val", 0)
		r.exact = h.AttrStr("hRule", "atLeast") == "exact"
	}
	if k := n.Child("cantSplit"); k != nil {
		r.cantSplit = onOff(k)
	}
	if k := n.Child("tblHeader"); k != nil {
		r.header = onOff(k)
	}
	if g := n.Child("gridBefore"); g != nil {
		r.gridBefore = int(g.AttrInt("val", 0))
	}
}

func (c *converter) applyTcPr(ce *cell, n *ooxml.Node) {
	if n == nil {
		return
	}
	if b := n.Child("tcBorders"); b != nil {
		c.readBorders(&ce.borders, b)
	}
	if s := n.Child("shd"); s != nil {
		ce.shd, ce.hasShd = c.shading(s)
	}
	if m := n.Child("tcMar"); m != nil {
		readMargins(&ce.mar, m)
	}
	if v := n.Child("vAlign"); v != nil {
		ce.vAlign = val(v)
	}
}

// cellRef names a cell in the A1 style, from the table's top left.
func cellRef(row, col int) string {
	name := ""
	for col++; col > 0; col = (col - 1) / 26 {
		name = string(rune('A'+(col-1)%26)) + name
	}
	return name + strconv.Itoa(row+1)
}

// tableLayout is a table laid out, before it is placed.
type tableLayout struct {
	t     *table
	left  float64   // of the table, in the column
	colX  []float64 // left of each grid column, relative to the table
	rowH  []float64
	cells map[*cell]*cellLayout
}

type cellLayout struct {
	ops  []op
	h    float64 // content height
	x, w float64 // box, relative to the table
}

// layoutTable lays out the cells of a table in a column of width w.
func (c *converter) layoutTable(t *table, w float64, sec *section, vertical bool) *tableLayout {
	lay := &tableLayout{t: t, cells: map[*cell]*cellLayout{}}
	grid := append([]float64(nil), t.grid...)
	total := 0.0
	for _, g := range grid {
		total += g
	}
	if total <= 0 {
		for i := range grid {
			grid[i] = w / float64(len(grid))
		}
		total = w
	}
	lay.colX = make([]float64, len(grid)+1)
	for i, g := range grid {
		lay.colX[i+1] = lay.colX[i] + g
	}
	switch t.jc {
	case "center":
		lay.left = (w - total) / 2
	case "right":
		lay.left = w - total
	default:
		lay.left = t.ind
		if c.compatMode < 15 && len(t.rows) > 0 && len(t.rows[0].cells) > 0 {
			lay.left -= t.rows[0].cells[0].mar[1]
		}
	}
	lay.rowH = make([]float64, len(t.rows))
	for ri, r := range t.rows {
		h := 0.0
		for _, ce := range r.cells {
			x0 := lay.colX[min(ce.col, len(grid))]
			x1 := lay.colX[min(ce.col+ce.span, len(grid))]
			cl := &cellLayout{x: x0, w: x1 - x0}
			lay.cells[ce] = cl
			if ce.vmerge == "continue" {
				continue
			}
			f := c.subflow(math.Max(cl.w-ce.mar[1]-ce.mar[3], 1), sec, vertical)
			f.blocks(ce.blocks, nil)
			cl.ops = f.all()
			cl.h = f.y
			if ce.rowSpan == 1 {
				h = math.Max(h, cl.h+ce.mar[0]+ce.mar[2])
			}
		}
		switch {
		case r.exact && r.height > 0:
			h = r.height
		default:
			h = math.Max(h, r.height)
		}
		lay.rowH[ri] = math.Max(h, 1)
	}
	// a merged cell taller than its rows stretches the last of them
	for ri, r := range t.rows {
		for _, ce := range r.cells {
			if ce.rowSpan <= 1 || ce.vmerge != "restart" {
				continue
			}
			cl := lay.cells[ce]
			sum := 0.0
			last := min(ri+ce.rowSpan, len(t.rows)) - 1
			for k := ri; k <= last; k++ {
				sum += lay.rowH[k]
			}
			if need := cl.h + ce.mar[0] + ce.mar[2]; need > sum && !t.rows[last].exact {
				lay.rowH[last] += need - sum
			}
		}
	}
	return lay
}

// table places a table, row by row. Rows that do not fit break across
// columns between lines of their cells (unless they may not split), and
// the header rows repeat at the top of each column.
func (f *flow) table(t *table) position {
	lay := f.c.layoutTable(t, f.w, f.sec, f.lc.vertical)
	if !f.atTop || !f.paged {
		f.y += f.prevAfter
	}
	f.prevAfter = 0
	f.prevPara = nil
	start := f.pos()
	nh := 0
	for nh < len(t.rows) && t.rows[nh].header {
		nh++
	}
	fresh := false // at the top of a column, after the repeated headers
	for ri := range t.rows {
		cut := 0.0
		for {
			rem := lay.rowH[ri] - cut
			avail := f.limit() - f.y
			if !f.paged || rem <= avail+0.01 {
				f.emitRow(lay, ri, f.y, cut, lay.rowH[ri])
				f.y += rem
				f.atTop, fresh = false, false
				break
			}
			r := t.rows[ri]
			if !r.cantSplit && !r.header || fresh || f.atTop {
				if c := lay.safeCut(ri, cut, cut+avail); c > cut+1 {
					f.emitRow(lay, ri, f.y, cut, c)
					f.y += c - cut
					cut = c
					f.atTop = false
					f.nextColumn(false)
					fresh = f.repeatHeaders(lay, nh, ri)
					continue
				}
			}
			if fresh || f.atTop {
				// taller than the column: let it run over
				f.emitRow(lay, ri, f.y, cut, lay.rowH[ri])
				f.y += rem
				f.atTop, fresh = false, false
				break
			}
			f.nextColumn(false)
			fresh = f.repeatHeaders(lay, nh, ri)
		}
	}
	return start
}

// rowRefs places the footnotes a row references on the page and records
// where its bookmarks went (its lines were laid out in the cells).
func (f *flow) rowRefs(r *row, y float64) {
	for _, ce := range r.cells {
		walkParas(ce.blocks, func(p *para) {
			f.c.recordBookmarks(f, p, y)
			if !f.paged {
				return
			}
			for _, it := range p.items {
				if it.note == nil || it.note.kind != "f" || slices.ContainsFunc(f.pg.notes, func(u *noteUse) bool { return u.n == it.note }) {
					continue
				}
				if len(f.pg.notes) == 0 {
					f.pg.noteH += noteSeparator
				}
				h := f.c.noteHeight(it.note, f.sec)
				f.pg.notes = append(f.pg.notes, &noteUse{n: it.note, h: h})
				f.pg.noteH += h
			}
		})
	}
}

// walkParas calls fn for the paragraphs of blocks, also those in tables.
func walkParas(bs []block, fn func(*para)) {
	for _, b := range bs {
		switch b := b.(type) {
		case *para:
			fn(b)
		case *table:
			for _, r := range b.rows {
				for _, ce := range r.cells {
					walkParas(ce.blocks, fn)
				}
			}
		}
	}
}

func (f *flow) repeatHeaders(lay *tableLayout, nh, ri int) bool {
	if ri < nh {
		return true
	}
	f.repeating = true
	for h := range nh {
		f.emitRow(lay, h, f.y, 0, lay.rowH[h])
		f.y += lay.rowH[h]
	}
	f.repeating = false
	f.atTop = false
	return true
}

// safeCut returns the lowest position in (from, to] of a row, in row
// coordinates, that no line of its cells crosses (from when there is none).
func (lay *tableLayout) safeCut(ri int, from, to float64) float64 {
	y := to
	for range 100 {
		moved := false
		for _, ce := range lay.t.rows[ri].cells {
			cl := lay.cells[ce]
			if cl == nil || ce.vmerge == "continue" {
				continue
			}
			for _, o := range cl.ops {
				if !o.text {
					continue
				}
				top, bottom := o.y0+ce.mar[0], o.y1+ce.mar[0]
				if top < y-0.01 && bottom > y+0.01 {
					y = top
					moved = true
				}
			}
		}
		if !moved {
			break
		}
	}
	if y <= from+0.01 {
		return from
	}
	return y
}

// emitRow draws the part [c0, c1) of a row (row coordinates) with its top
// at y.
func (f *flow) emitRow(lay *tableLayout, ri int, y, c0, c1 float64) {
	t := lay.t
	r := t.rows[ri]
	x := f.x0 + lay.left
	partial := c0 > 0 || c1 < lay.rowH[ri]
	if c0 == 0 && !f.repeating && f.main {
		f.rowRefs(r, y)
	}
	for ci, ce := range r.cells {
		if ce.vmerge == "continue" {
			continue
		}
		cl := lay.cells[ce]
		h := c1 - c0
		if ce.rowSpan > 1 && !partial {
			for k := ri + 1; k < min(ri+ce.rowSpan, len(t.rows)); k++ {
				h += lay.rowH[k]
			}
		}
		bx, bw := x+cl.x, cl.w
		// background
		shd, has := ce.shd, ce.hasShd
		if !has {
			shd, has = t.shd, t.hasShd
		}
		if has {
			col := shd
			f.emit(op{y0: y, y1: y + h, fn: func(e *emitter, dx, dy float64) {
				e.fillRect(col, bx+dx, y+dy, bw, h)
			}})
		}
		// the cell's structure, also when it is empty
		f.emit(op{y0: y, y1: y, text: true, ctx: ce.ctx, fn: func(*emitter, float64, float64) {}})
		// content
		voff := ce.mar[0]
		if !partial {
			switch ce.vAlign {
			case "center":
				voff += (h - ce.mar[0] - ce.mar[2] - cl.h) / 2
			case "bottom":
				voff += h - ce.mar[0] - ce.mar[2] - cl.h
			}
		}
		if partial {
			clipY, clipH := y, h
			f.emit(op{y0: y, y1: y + h, fn: func(e *emitter, dx, dy float64) {
				e.cv.Obj.Save()
				e.cv.Obj.ClipRect(f32(bx+dx), f32(clipY+dy), f32(bw), f32(clipH))
			}})
		}
		for _, o := range cl.ops {
			top, bottom := o.y0+ce.mar[0], o.y1+ce.mar[0]
			if partial {
				if o.text && (top < c0-0.01 || top >= c1-0.01) && !(o.y0 == o.y1 && top >= c0 && top < c1) {
					continue
				}
				if !o.text && (bottom <= c0 || top >= c1) {
					continue
				}
			}
			o.shift(bx+ce.mar[1], y+voff-c0)
			f.emit(o)
		}
		if partial {
			f.emit(op{y0: y, y1: y + h, fn: func(e *emitter, dx, dy float64) {
				e.cv.Obj.Restore()
				e.invalidate()
			}})
		}
		// borders
		var bs [4]*border
		for side := range 4 {
			bs[side] = lay.cellBorder(ri, ci, ce, side)
		}
		if partial && c0 > 0 && bs[sTop] == nil {
			bs[sTop] = t.borders[sInsideH]
		}
		f.emit(op{y0: y, y1: y + h, fn: func(e *emitter, dx, dy float64) {
			e.cellBorders(bs, bx+dx, y+dy, bw, h)
		}})
	}
}

// cellBorder resolves one side of a cell's border: the cell's own, then
// the table's (outer or inside).
func (lay *tableLayout) cellBorder(ri, ci int, ce *cell, side int) *border {
	if b := ce.borders[side]; b != nil {
		if b.none() {
			return nil
		}
		return b
	}
	t := lay.t
	r := t.rows[ri]
	var b *border
	switch side {
	case sTop:
		b = t.borders[sInsideH]
		if ri == 0 {
			b = t.borders[sTop]
		}
	case sBottom:
		b = t.borders[sInsideH]
		if ri+ce.rowSpan >= len(t.rows) {
			b = t.borders[sBottom]
		}
	case sLeft:
		b = t.borders[sInsideV]
		if ci == 0 {
			b = t.borders[sLeft]
		}
	case sRight:
		b = t.borders[sInsideV]
		if ci == len(r.cells)-1 {
			b = t.borders[sRight]
		}
	}
	if b.none() {
		return nil
	}
	return b
}
