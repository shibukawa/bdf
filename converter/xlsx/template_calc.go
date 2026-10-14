package xlsx

import (
	"sort"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/form/tag"
)

// remap moves what refers to the rows of a sheet of a filled template to
// where the rows are now: a range in rows that were copied becomes a range
// in each copy.
func (f *filler) remap(ts *tplSheet) {
	ws := ts.ws
	ranges := func(rg cellRange) []cellRange {
		var out []cellRange
		ts.eachCopy(rg.r0, rg.r1, func(d0, d1 int) {
			out = append(out, cellRange{d0, rg.c0, d1, rg.c1})
		})
		return out
	}
	// what a few bytes of a template become in the copies of its rows is
	// bounded: the rest is left out, with a warning
	full := func(n int, what string) bool {
		if n < maxMovedRanges {
			return false
		}
		f.c.warnOnce("moved:"+ws.name+what, "sheet %q: the copies of its rows have more than %d %s; the rest are left out", ws.name, maxMovedRanges, what)
		return true
	}
	// merged ranges
	var merges []cellRange
	for _, m := range ws.merges {
		if full(len(merges), "merged ranges") {
			break
		}
		merges = append(merges, ranges(m)...)
	}
	ws.merges = merges
	// conditional formats: their ranges
	formats := 0
	sqref := func(s string) string {
		var out []string
		for _, rg := range parseSqref(s) {
			if full(formats, "ranges of conditional formats") {
				break
			}
			for _, r := range ranges(rg) {
				out = append(out, cellRef(r.r0, r.c0)+":"+cellRef(r.r1, r.c1))
				formats++
			}
		}
		return strings.Join(out, " ")
	}
	for _, n := range ws.cfs {
		setAttr(n, "sqref", sqref(n.AttrStr("sqref", "")))
	}
	for _, n := range ws.x14cfs {
		if sq := n.Child("sqref"); sq != nil {
			sq.Text = sqref(sq.Content())
		}
	}
	// links
	var links []*ooxml.Node
	for _, n := range ws.links {
		rg, ok := parseRange(n.AttrStr("ref", ""))
		if !ok {
			continue
		}
		if full(len(links), "links") {
			break
		}
		for _, r := range ranges(rg) {
			cp := *n
			cp.Attrs = append([]ooxml.Attr(nil), n.Attrs...)
			setAttr(&cp, "ref", cellRef(r.r0, r.c0)+":"+cellRef(r.r1, r.c1))
			links = append(links, &cp)
		}
	}
	ws.links = links
	if af := ws.autoFilter; af != nil {
		if r := ts.span(*af); r != nil {
			ws.autoFilter = r
		} else {
			ws.autoFilter = nil
		}
	}
	// the print areas and titles hold every copy of their rows
	if pn := f.c.printNames[ts.ref.idx]; pn != nil {
		var areas []cellRange
		for _, a := range pn.areas {
			if a.r0 == 0 && a.r1 == maxRows-1 {
				areas = append(areas, a) // whole columns
			} else if r := ts.span(a); r != nil {
				areas = append(areas, *r)
			}
		}
		pn.areas = areas
		if t := pn.titleRows; t != nil {
			if r := ts.span(cellRange{t[0], 0, t[1], 0}); r != nil {
				pn.titleRows = &[2]int{r.r0, r.r1}
			} else {
				pn.titleRows = nil
			}
		}
	}
	// page breaks: before each copy of the rows they are before, and
	// before the copies that are pages
	breaks := append([]int(nil), ts.breaks...)
	for _, b := range ws.print.rowBreaks {
		ts.eachCopy(b, b, func(d0, _ int) { breaks = append(breaks, d0) })
	}
	sort.Ints(breaks)
	ws.print.rowBreaks = ws.print.rowBreaks[:0]
	for i, b := range breaks {
		if b > 0 && b < maxRows && (i == 0 || b != breaks[i-1]) {
			ws.print.rowBreaks = append(ws.print.rowBreaks, b)
		}
	}
	ws.formulas = nil
}

// span returns the rows of the filled sheet from the first copy of the
// first row of a range to the last copy of its last; nil when the rows are
// gone.
func (ts *tplSheet) span(rg cellRange) *cellRange {
	d0, g0 := ts.mapRow(rg.r0, nil, false)
	d1, g1 := ts.mapRow(rg.r1, nil, true)
	if g0 && g1 || d1 < d0 || d0 >= maxRows {
		return nil
	}
	return &cellRange{d0, rg.c0, min(d1, maxRows-1), rg.c1}
}

// setAttr sets an attribute of an element.
func setAttr(n *ooxml.Node, name, value string) {
	for i := range n.Attrs {
		if n.Attrs[i].Name.Local == name {
			n.Attrs[i].Value = value
			return
		}
	}
	n.Attrs = append(n.Attrs, ooxml.Attr{Name: ooxml.Name{Local: name}, Value: value})
}

// filledRows moves a range of rows of the template of a sheet to the rows
// of the filled sheet (for what is read after the sheet is filled: tables
// and drawings).
func (c *converter) filledRows(ws *worksheet, r0, r1 int) (int, int) {
	if c.opts == nil || c.opts.fill == nil {
		return r0, r1
	}
	ts := c.opts.fill.byName[ws.name]
	if ts == nil {
		return r0, r1
	}
	if r := ts.span(cellRange{r0: r0, r1: r1}); r != nil {
		return r.r0, r.r1
	}
	return r0, r0
}

// recalc calculates the formulas of the filled sheets again, each after
// the formulas it reads.
func (f *filler) recalc() {
	type at struct {
		ts  *tplSheet
		key [2]int
	}
	for _, ts := range f.sheets {
		keys := make([][2]int, 0, len(ts.outF))
		for k := range ts.outF {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i][0] < keys[j][0] || keys[i][0] == keys[j][0] && keys[i][1] < keys[j][1]
		})
		for _, k := range keys {
			if ts.outF[k].state != 0 {
				continue
			}
			// depth first, without recursion: a chain of formulas may be
			// as long as the sheet
			stack := []at{{ts, k}}
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				of := top.ts.outF[top.key]
				switch of.state {
				case 0:
					of.state = 1
					f.precedents(top.ts, top.key, of, func(ts2 *tplSheet, k2 [2]int) {
						if ts2.outF[k2].state == 0 {
							stack = append(stack, at{ts2, k2})
						}
					})
				case 1:
					f.calcOne(top.ts, top.key, of)
					of.state = 2
					stack = stack[:len(stack)-1]
				default:
					stack = stack[:len(stack)-1]
				}
			}
		}
	}
}

// refRange returns the cells of the filled sheets that a reference of a
// formula of a cell reads; gone when their rows are gone.
func (f *filler) refRange(ts *tplSheet, of *outFormula, n fRef) (ts2 *tplSheet, rg cellRange, gone bool) {
	ts2 = ts
	if n.sheet != "" && n.sheet != ts.ws.name {
		if ts2 = f.byName[n.sheet]; ts2 == nil {
			return nil, rg, true
		}
	}
	dr, dc := of.f.dr+of.dr, of.f.dc
	move := func(v int, abs bool, d int) int {
		if abs {
			return v
		}
		return v + d
	}
	switch {
	case n.wholeCols:
		return ts2, cellRange{0, move(n.c0, n.absC0, dc), maxRows - 1, move(n.c1, n.absC1, dc)}, false
	case n.wholeRows:
		rg = cellRange{move(n.r0, n.absR0, dr), 0, move(n.r1, n.absR1, dr), maxCols - 1}
	default:
		rg = cellRange{move(n.r0, n.absR0, dr), move(n.c0, n.absC0, dc), move(n.r0, n.absR0, dr), move(n.c0, n.absC0, dc)}
		if n.isRange {
			rg.r1, rg.c1 = move(n.r1, n.absR1, dr), move(n.c1, n.absC1, dc)
		}
	}
	if rg.r0 > rg.r1 {
		rg.r0, rg.r1 = rg.r1, rg.r0
	}
	if rg.c0 > rg.c1 {
		rg.c0, rg.c1 = rg.c1, rg.c0
	}
	if rg.r0 < 0 || rg.c0 < 0 || rg.r1 >= maxRows || rg.c1 >= maxCols {
		return ts2, rg, true
	}
	var ctx *itemInst
	if ts2 == ts {
		ctx = of.ctx
	}
	d0, g0 := ts2.mapRow(rg.r0, ctx, false)
	d1, g1 := ts2.mapRow(rg.r1, ctx, true)
	rg.r0, rg.r1 = d0, min(d1, maxRows-1)
	return ts2, rg, g0 && g1 || d1 < d0 || d0 >= maxRows
}

// precedents calls fn for the cells with formulas that a formula reads.
func (f *filler) precedents(ts *tplSheet, key [2]int, of *outFormula, fn func(*tplSheet, [2]int)) {
	var walk func(n fnode)
	walk = func(n fnode) {
		switch x := n.(type) {
		case fRef:
			ts2, rg, gone := f.refRange(ts, of, x)
			if gone || ts2 == nil || len(ts2.outF) == 0 {
				return
			}
			rows := ts2.ws.rows
			i := sort.Search(len(rows), func(i int) bool { return rows[i].idx >= rg.r0 })
			for ; i < len(rows) && rows[i].idx <= rg.r1; i++ {
				cells := rows[i].cells
				j := sort.Search(len(cells), func(j int) bool { return cells[j].col >= rg.c0 })
				for ; j < len(cells) && cells[j].col <= rg.c1; j++ {
					if f.looked++; f.looked > maxPrecedents {
						return // the formulas are calculated in the order of the cells
					}
					k := [2]int{rows[i].idx, cells[j].col}
					if _, ok := ts2.outF[k]; ok && !(ts2 == ts && k == key) {
						fn(ts2, k)
					}
				}
			}
		case fUnary:
			walk(x.x)
		case fBinary:
			walk(x.a)
			walk(x.b)
		case fCall:
			for _, a := range x.args {
				walk(a)
			}
		}
	}
	if of.f.node != nil {
		walk(of.f.node)
	}
}

// calcOne calculates the formula of a cell and makes the result the value
// of the cell; a formula that is not calculated leaves the cell empty.
func (f *filler) calcOne(ts *tplSheet, key [2]int, of *outFormula) {
	cl := ts.ws.cellAt(key[0], key[1])
	if cl == nil {
		return
	}
	cl.kind, cl.text, cl.num = cellBlank, nil, 0
	unsupported := func() {
		f.c.warnOnce("formula:"+ts.ws.name+"\x00"+of.f.src, "sheet %q cell %s: the formula %q is not calculated; the cell is left empty", ts.ws.name, cellRef(key[0], key[1]), of.f.src)
	}
	if of.f.node == nil {
		unsupported()
		return
	}
	env := &fenv{s: ts.calc, dr: of.f.dr + of.dr, dc: of.f.dc, row: key[0], col: key[1]}
	env.remap = func(n fRef) (fval, bool) {
		ts2, rg, gone := f.refRange(ts, of, n)
		switch {
		case ts2 == nil:
			return ferr("#REF!"), true
		case gone && !n.isRange:
			return ferr("#REF!"), true
		case gone:
			// the rows of a range are gone: nothing to add up
			return fval{kind: fRange, rng: &frange{ws: ts2.ws, cellRange: cellRange{rg.r0, rg.c0, rg.r0 - 1, rg.c1}}}, true
		case !n.isRange:
			return cellFval(ts2.ws.cellAt(rg.r0, rg.c0)), true
		}
		return fval{kind: fRange, rng: &frange{ws: ts2.ws, cellRange: rg}}, true
	}
	v, err := env.eval(of.f.node)
	if err != nil {
		unsupported()
		return
	}
	switch v = env.scalar(v); v.kind {
	case fNum:
		cl.kind, cl.num = cellNum, v.num
	case fBool:
		cl.kind, cl.num = cellBool, v.num
	case fStr:
		if v.str != "" {
			cl.kind, cl.text = cellStr, &richText{plain: clipText(v.str)}
		}
	case fErr:
		cl.kind, cl.text = cellErr, &richText{plain: v.str}
	}
}

// pictures returns the pictures that the tags of a filled template put
// over the cells of the sheet. A cell with nothing but a picture is a
// cell that is used.
func (s *sheetCtx) pictures() []tplPicture {
	if s.c == nil || s.c.opts == nil || s.c.opts.fill == nil {
		return nil
	}
	if ts := s.c.opts.fill.byName[s.ws.name]; ts != nil {
		return ts.pictures
	}
	return nil
}

// paintPictures draws the pictures that the tags of a filled template put
// over cells: each in its cell, or the merged range of the cell.
func (s *sheetCtx) paintPictures() {
	for _, p := range s.pictures() {
		if p.row >= s.nRows || p.col >= s.nCols {
			continue
		}
		b := s.cellBox(p.row, p.col)
		if m, ok := s.mergeAt(p.row, p.col); ok {
			b = s.rangeBox(m)
		}
		// inside the borders of the cell
		b = box{b.x + pxPt, b.y + pxPt, b.w - 2*pxPt, b.h - 2*pxPt}
		if b.w <= 0 || b.h <= 0 {
			continue
		}
		w, h := b.w, b.h
		switch {
		case p.fit == tag.FitStretch, p.sym != nil && p.sym.Linear:
			// (the bars of a barcode are as wide and as tall as their place)
		case p.fit == tag.FitCover && p.sym == nil:
			k := max(b.w/p.w, b.h/p.h)
			w, h = p.w*k, p.h*k
		default:
			k := min(b.w/p.w, b.h/p.h)
			w, h = p.w*k, p.h*k
		}
		x, y := b.x+(b.w-w)/2, b.y+(b.h-h)/2
		for _, t := range s.tilesIn(b.x, b.y, b.x+b.w, b.y+b.h) {
			o := t.cv.Obj
			o.Mark(bdf.MarkFigure, p.alt)
			t.save()
			o.ClipRect(f32(b.x-t.ox), f32(b.y-t.oy), f32(b.w), f32(b.h))
			if p.sym != nil {
				t.fillColor(bdf.RGB(0, 0, 0))
				o.FillPath(o.AddPath(p.sym.Path(f32(x-t.ox), f32(y-t.oy), f32(w), f32(h))), bdf.NonZero)
			} else {
				o.Image(t.cv.Image(p.hash), f32(x-t.ox), f32(y-t.oy), f32(w), f32(h))
			}
			t.restore()
			o.Mark(bdf.MarkEnd, "")
			t.cv.Drawn = true
		}
	}
}
