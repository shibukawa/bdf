package xlsx

import (
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// A worksheet can also be laid out as Excel prints it (the view "print"):
// its print areas, or what it uses, cut into pages by the paper, the
// margins and the scaling of its page setup, with the print titles repeated
// on the pages after them, the manual page breaks kept, and the header and
// the footer on every page.
//
// The cells are drawn as they are for the tiles of the sheet view, into
// regions instead of tiles: a region is the rows and columns of a page (or
// the title rows over a page's columns, or the title columns beside its
// rows), an object in sheet coordinates from the region's corner. A page's
// body layer uses its regions, clipped, scaled and moved to their place on
// the paper; regions with the same content (the titles of the pages of a
// column of pages) are one object.
//
// Excel rounds row heights and column widths to the pixels of the printer
// when it prints, so where it breaks a page by itself can differ by a row
// from where this does; a manual page break is where it is.

// View selections (Options.Views).
const (
	ViewsSheet = "sheet"
	ViewsPrint = "print"
	ViewsBoth  = "both"
)

const (
	// maxPrintPages bounds the pages of the print view of a sheet: a print
	// area is the file's to state, and one of a million empty rows is
	// twenty thousand pages.
	maxPrintPages = 1 << 14
	// maxPrintAreas bounds the print areas of a sheet.
	maxPrintAreas = 256
	// Excel's limits of the scale of a printed sheet.
	minPrintScale = 0.10
	maxPrintScale = 4.00
)

// printSetup is what a worksheet says about printing it.
type printSetup struct {
	paper        int     // paperSize; 0 when the sheet names none
	paperW       float64 // paperWidth and paperHeight in points; 0 when not given
	paperH       float64
	landscape    bool
	scale        float64 // percent; 0 when not given
	fitToPage    bool
	fitW, fitH   int // pages across and down to fit on; 0 leaves the direction alone
	firstPage    int
	useFirstPage bool
	overThenDown bool
	margins      *[6]float64 // left, right, top, bottom, header, footer in points
	gridLines    bool
	headings     bool
	hCenter      bool
	vCenter      bool
	rowBreaks    []int // rows that start a page (0-based), ascending
	colBreaks    []int
	hf           *headerFooter
	set          bool // the sheet has a page setup
}

// element reads one of the elements of a worksheet about printing.
func (ps *printSetup) element(n *ooxml.Node) {
	switch n.Name {
	case "sheetPr":
		ps.fitToPage = n.Child("pageSetUpPr").AttrBool("fitToPage", false)
	case "pageMargins":
		in := func(name string, def float64) float64 {
			v := n.AttrFloat(name, def)
			if math.IsNaN(v) || math.IsInf(v, 0) || v < 0 {
				v = def
			}
			return math.Min(v, 100) * 72
		}
		ps.margins = &[6]float64{in("left", 0.7), in("right", 0.7), in("top", 0.75), in("bottom", 0.75), in("header", 0.3), in("footer", 0.3)}
	case "pageSetup":
		ps.paper = int(n.AttrInt("paperSize", 0))
		ps.paperW, ps.paperH = measurePt(n.AttrStr("paperWidth", "")), measurePt(n.AttrStr("paperHeight", ""))
		ps.landscape = n.AttrStr("orientation", "") == "landscape"
		if v := n.AttrFloat("scale", 0); v > 0 && !math.IsInf(v, 0) {
			ps.scale = v
		}
		ps.fitW, ps.fitH = 1, 1
		if _, ok := n.Attr("fitToWidth"); ok {
			ps.fitW = int(min(max(n.AttrInt("fitToWidth", 1), 0), maxPrintPages))
		}
		if _, ok := n.Attr("fitToHeight"); ok {
			ps.fitH = int(min(max(n.AttrInt("fitToHeight", 1), 0), maxPrintPages))
		}
		ps.set = true
		ps.firstPage = int(min(max(n.AttrInt("firstPageNumber", 1), -1<<20), 1<<20))
		ps.useFirstPage = n.AttrBool("useFirstPageNumber", false)
		ps.overThenDown = n.AttrStr("pageOrder", "downThenOver") == "overThenDown"
	case "printOptions":
		ps.gridLines = n.AttrBool("gridLines", false) && n.AttrBool("gridLinesSet", true)
		ps.headings = n.AttrBool("headings", false)
		ps.hCenter = n.AttrBool("horizontalCentered", false)
		ps.vCenter = n.AttrBool("verticalCentered", false)
	case "rowBreaks", "colBreaks":
		limit := maxRows
		if n.Name == "colBreaks" {
			limit = maxCols
		}
		var out []int
		for _, b := range n.Children("brk") {
			// a break is after the row or column of its id, counted from 1:
			// the one of that index, counted from 0, starts the page
			if id := b.AttrInt("id", 0); id > 0 && id < int64(limit) && b.AttrBool("man", false) {
				out = append(out, int(id))
			}
		}
		sort.Ints(out)
		if n.Name == "rowBreaks" {
			ps.rowBreaks = out
		} else {
			ps.colBreaks = out
		}
	case "headerFooter":
		ps.hf = readHeaderFooter(n)
	}
}

// measurePt reads a universal measure ("210mm", "8.5in") in points; 0 for
// none.
func measurePt(v string) float64 {
	v = strings.TrimSpace(v)
	for _, u := range []struct {
		suffix string
		pt     float64
	}{{"mm", 72 / 25.4}, {"cm", 72 / 2.54}, {"in", 72}, {"pt", 1}, {"pc", 12}, {"pi", 12}} {
		if strings.HasSuffix(v, u.suffix) {
			f, err := strconv.ParseFloat(strings.TrimSuffix(v, u.suffix), 64)
			if err != nil || math.IsNaN(f) || math.IsInf(f, 0) || f <= 0 {
				return 0
			}
			return f * u.pt
		}
	}
	return 0
}

const mm = 72 / 25.4

// paperSizes are the paper sizes of ST_PaperSize that printers have, in
// points, upright.
var paperSizes = map[int][2]float64{
	1: {612, 792}, 2: {612, 792}, 3: {792, 1224}, 4: {1224, 792}, 5: {612, 1008}, 6: {396, 612}, 7: {522, 756},
	8: {297 * mm, 420 * mm}, 9: {210 * mm, 297 * mm}, 10: {210 * mm, 297 * mm}, 11: {148 * mm, 210 * mm},
	12: {257 * mm, 364 * mm}, 13: {182 * mm, 257 * mm}, 14: {612, 936}, 15: {215 * mm, 275 * mm}, 16: {720, 1008},
	17: {792, 1224}, 18: {612, 792}, 19: {279, 639}, 20: {297, 684}, 21: {324, 747}, 22: {342, 792}, 23: {360, 828},
	27: {110 * mm, 220 * mm}, 28: {162 * mm, 229 * mm}, 29: {324 * mm, 458 * mm}, 30: {229 * mm, 324 * mm},
	31: {114 * mm, 162 * mm}, 32: {114 * mm, 229 * mm}, 33: {250 * mm, 353 * mm}, 34: {176 * mm, 250 * mm},
	35: {176 * mm, 125 * mm}, 36: {110 * mm, 230 * mm}, 37: {279, 540}, 38: {261, 468}, 39: {1071, 792},
	40: {612, 864}, 41: {612, 936}, 42: {250 * mm, 353 * mm}, 43: {100 * mm, 148 * mm}, 44: {648, 792},
	45: {720, 792}, 46: {1080, 792}, 47: {220 * mm, 220 * mm},
	66: {420 * mm, 594 * mm}, 69: {200 * mm, 148 * mm}, 70: {105 * mm, 148 * mm},
	75: {792, 612}, 76: {420 * mm, 297 * mm}, 77: {297 * mm, 210 * mm}, 78: {210 * mm, 148 * mm},
	79: {364 * mm, 257 * mm}, 80: {257 * mm, 182 * mm}, 81: {148 * mm, 100 * mm}, 82: {148 * mm, 200 * mm},
	83: {148 * mm, 105 * mm}, 88: {128 * mm, 182 * mm}, 89: {182 * mm, 128 * mm},
}

// paperNames are the papers Options.Paper names, by their paperSize.
var paperNames = map[string]int{"letter": 1, "tabloid": 3, "legal": 5, "a3": 8, "a4": 9, "a5": 11, "b4": 12, "b5": 13, "a6": 70, "b6": 88, "postcard": 43}

// PaperNames returns the papers that Options.Paper names.
func PaperNames() []string {
	out := make([]string, 0, len(paperNames))
	for n := range paperNames {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// printNames is what the defined names of a workbook say about printing a
// sheet: its print areas and its print titles.
type printNames struct {
	areas     []cellRange
	titleRows *[2]int // first and last, 0-based
	titleCols *[2]int
	// Optional BDF pagination names; ordinary Excel print titles are unchanged.
	titleScope *cellRange
	keepRows   []cellRange
}

// loadPrintNames reads the print areas, titles and optional BDF pagination
// settings from workbook defined names, by the index of the sheet.
func (c *converter) loadPrintNames(wb *ooxml.Node) {
	c.printNames = map[int]*printNames{}
	for _, dn := range wb.Child("definedNames").Children("definedName") {
		name := dn.AttrStr("name", "")
		area := strings.EqualFold(name, "_xlnm.Print_Area")
		scope := strings.EqualFold(name, "_bdf.Print_TitleScope")
		keep := strings.EqualFold(name, "_bdf.Print_KeepTogether")
		if !area && !scope && !keep && !strings.EqualFold(name, "_xlnm.Print_Titles") {
			continue
		}
		id, ok := dn.Attr("localSheetId")
		if !ok {
			continue
		}
		idx, err := strconv.Atoi(strings.TrimSpace(id))
		if err != nil || idx < 0 {
			continue
		}
		pn := c.printNames[idx]
		if pn == nil {
			pn = &printNames{}
			c.printNames[idx] = pn
		}
		for _, ref := range splitRefs(dn.Content()) {
			rg, ok := parseRange(ref)
			if !ok {
				continue
			}
			switch {
			case scope:
				cp := rg
				pn.titleScope = &cp
			case keep:
				if len(pn.keepRows) < maxPrintAreas {
					pn.keepRows = append(pn.keepRows, rg)
				}
			case area:
				if len(pn.areas) < maxPrintAreas {
					pn.areas = append(pn.areas, rg)
				}
			case rg.c0 == 0 && rg.c1 == maxCols-1:
				pn.titleRows = &[2]int{rg.r0, rg.r1}
			case rg.r0 == 0 && rg.r1 == maxRows-1:
				pn.titleCols = &[2]int{rg.c0, rg.c1}
			}
		}
	}
}

// splitRefs splits the value of a defined name ("Sheet1!$A$1:$C$5,'a,b'!$E:$F")
// into its references without their sheets.
func splitRefs(s string) []string {
	var out []string
	quoted, start := false, 0
	cut := func(end int) {
		ref := s[start:end]
		// the sheet name ends at the last ! outside quotes
		q, bang := false, -1
		for i := 0; i < len(ref); i++ {
			switch ref[i] {
			case '\'':
				q = !q
			case '!':
				if !q {
					bang = i
				}
			}
		}
		if ref = strings.TrimSpace(ref[bang+1:]); ref != "" {
			out = append(out, ref)
		}
	}
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\'':
			quoted = !quoted
		case ',':
			if !quoted {
				cut(i)
				start = i + 1
			}
		}
	}
	cut(len(s))
	return out
}

// printBand is the rows or the columns of a page: the entries i0 to i1-1.
type printBand struct{ i0, i1 int }

// printRegion is cells that pages show: an object in sheet coordinates from
// the corner of its box, in three layers (what is under the borders, the
// borders, and the drawings over them).
type printRegion struct {
	b      box
	layers [3]*tileCv
}

// printArea is a print area cut into pages: page (i, j) shows the rows of
// the band i and the columns of the band j.
type printArea struct {
	rg           cellRange
	rows         []printBand
	cols         []printBand
	rowPos       []float64 // where each band of rows starts on the sheet, and where the last ends
	colPos       []float64
	body         map[[2]int]*printRegion
	top          map[int]*printRegion // the title rows over the columns of band j
	left         map[int]*printRegion // the title columns beside the rows of band i
	hasRowTitles bool
}

// printPage is a page: the bands i and j of an area.
type printPage struct {
	a    *printArea
	i, j int
}

// printPlan is a sheet laid out for printing.
type printPlan struct {
	s      *sheetCtx
	ps     *printSetup
	pw, ph float64 // the paper
	// margins: left, right, top, bottom, and of the header and the footer
	ml, mr, mt, mb, mh, mf float64
	scale                  float64
	titleRows, titleCols   *[2]int
	titleScope             *cellRange
	keepRows               [][2]int // inclusive, sorted and overlapping groups joined
	// the title columns and rows on the sheet
	tx0, tx1, ty0, ty1 float64
	// the size of the row and column headings (0 when they are not printed)
	hdrW, hdrH float64
	areas      []*printArea
	corner     *printRegion // the title rows of the title columns
	pages      []printPage
	// layer is the layer of the regions that is drawn into
	layer int
	all   []*tileCv
	// usedR and usedC are the last row and column the sheet uses (-1 for
	// an empty sheet)
	usedR, usedC int
}

// usedRange finds the last row and column of what the sheet prints without
// a print area: its values, the cells with a fill or a border, the merged
// ranges, the tables, the drawings, and the text that flows out of the last
// cell of a row.
func (s *sheetCtx) usedRange() (r1, c1 int) {
	r1, c1 = -1, -1
	use := func(r, c int) { r1, c1 = max(r1, r), max(c1, c) }
	for i := range s.ws.rows {
		rw := &s.ws.rows[i]
		if rw.idx >= s.nRows {
			break
		}
		last := -1
		for j := range rw.cells {
			cl := &rw.cells[j]
			if cl.col >= s.nCols {
				break
			}
			if cl.kind == cellBlank {
				if f := s.fmtOf(cl.style); f.fill.empty() && !hasBorder(f) {
					continue
				}
			} else {
				last = j
			}
			use(rw.idx, cl.col)
		}
		if last < 0 || rw.cells[last].kind != cellStr {
			continue
		}
		cl := &rw.cells[last]
		if _, merged := s.mergeAt(rw.idx, cl.col); merged {
			continue
		}
		b := s.cellBox(rw.idx, cl.col)
		if b.w <= 0 || b.h <= 0 {
			continue
		}
		f := s.formatAt(rw.idx, cl.col, rw, cl)
		if f == nil {
			f = s.fmtOf(0)
		}
		if lay := s.layoutCell(rw.idx, cl.col, cl, f, b, false); lay != nil && len(lay.lines) > 0 {
			if right := lay.bounds.x + lay.bounds.w; right > b.x+b.w+0.01 {
				use(rw.idx, min(s.cols.index(right-0.01), maxCols-1))
			}
		}
	}
	for _, m := range s.merges {
		use(m.r1, m.c1)
	}
	for _, t := range s.tables {
		use(min(t.ref.r1, s.nRows-1), min(t.ref.c1, s.nCols-1))
	}
	for _, d := range s.drawings {
		if b := s.anchorBox(d); !d.out && b.w > 0 && b.h > 0 {
			use(min(s.rows.index(b.y+b.h-0.01), s.nRows-1), min(s.cols.index(b.x+b.w-0.01), s.nCols-1))
		}
	}
	for _, p := range s.pictures() {
		if p.row < s.nRows && p.col < s.nCols {
			use(p.row, p.col)
		}
	}
	return r1, c1
}

// planPrint lays the sheet out for printing: the paper, the scale and the
// pages. idx is the index of the sheet in the workbook.
func (s *sheetCtx) planPrint(idx int) *printPlan {
	c, ps := s.c, &s.ws.print
	p := &printPlan{s: s, ps: ps, scale: 1}
	p.usedR, p.usedC = s.usedRange()

	// the paper
	size, ok := paperSizes[ps.paper]
	if !ok {
		if ps.paper != 0 {
			c.warnOnce("paper:"+strconv.Itoa(ps.paper), "sheet %q: paper size %d is not known; the default paper is used", s.ws.name, ps.paper)
		}
		def := 1 // Letter, the default of the format
		if c.locale == "ja" {
			def = 9 // A4
		}
		if n, ok := paperNames[strings.ToLower(c.opts.Paper)]; ok {
			def = n
		}
		size = paperSizes[def]
	}
	p.pw, p.ph = size[0], size[1]
	if ps.paperW > 0 && ps.paperH > 0 {
		p.pw, p.ph = math.Min(ps.paperW, 14400), math.Min(ps.paperH, 14400)
	}
	if ps.landscape {
		p.pw, p.ph = p.ph, p.pw
	}
	m := [6]float64{0.7 * 72, 0.7 * 72, 0.75 * 72, 0.75 * 72, 0.3 * 72, 0.3 * 72}
	if ps.margins != nil {
		m = *ps.margins
	}
	// margins that leave nothing of the paper are not those of a page
	if m[0]+m[1] > p.pw-36 {
		m[0], m[1] = 0, 0
	}
	if m[2]+m[3] > p.ph-36 {
		m[2], m[3] = 0, 0
	}
	p.ml, p.mr, p.mt, p.mb, p.mh, p.mf = m[0], m[1], m[2], m[3], m[4], m[5]

	// the titles and the areas
	var names *printNames
	if pn := c.printNames[idx]; pn != nil {
		names = pn
	} else {
		names = &printNames{}
	}
	if t := names.titleRows; t != nil && t[0] < maxRows && t[0] <= t[1] {
		p.titleRows = &[2]int{t[0], min(t[1], maxRows-1)}
		p.ty0, p.ty1 = s.rows.pos(p.titleRows[0]), s.rows.pos(p.titleRows[1]+1)
	}
	if t := names.titleCols; t != nil && t[0] < maxCols && t[0] <= t[1] {
		p.titleCols = &[2]int{t[0], min(t[1], maxCols-1)}
		p.tx0, p.tx1 = s.cols.pos(p.titleCols[0]), s.cols.pos(p.titleCols[1]+1)
	}
	p.titleScope = names.titleScope
	for _, rg := range names.keepRows {
		p.keepRows = append(p.keepRows, [2]int{rg.r0, rg.r1})
	}
	// A scoped table heading belongs with its first detail row, even when
	// it falls at the bottom of the first page.
	if p.titleRows != nil && p.titleScope != nil && p.titleRows[1] < p.titleScope.r1 {
		p.keepRows = append(p.keepRows, [2]int{p.titleRows[0], p.titleRows[1] + 1})
	}
	p.keepRows = joinKeepRows(p.keepRows)
	if ps.headings {
		p.hdrH = s.rows.def
		if p.hdrH <= 0 {
			p.hdrH = 15
		}
		digits := len(strconv.Itoa(max(p.usedR, 0) + 1))
		p.hdrW = (float64(max(digits, 3))*c.mdw + 12) * pxPt
	}
	var ranges []cellRange
	for _, rg := range names.areas {
		// whole columns or rows are those of what the sheet uses
		if rg.r0 == 0 && rg.r1 == maxRows-1 {
			rg.r1 = p.usedR
		}
		if rg.c0 == 0 && rg.c1 == maxCols-1 {
			rg.c1 = p.usedC
		}
		if rg.r0 <= rg.r1 && rg.c0 <= rg.c1 {
			ranges = append(ranges, rg)
		}
	}
	if len(names.areas) == 0 && p.usedR >= 0 && p.usedC >= 0 {
		ranges = []cellRange{{0, 0, p.usedR, p.usedC}}
	}

	// the scale
	switch {
	case ps.fitToPage:
		p.scale = p.fitScale(ranges)
	case ps.scale > 0:
		p.scale = math.Min(math.Max(ps.scale/100, minPrintScale), maxPrintScale)
	}

	// the pages
	room := maxPrintPages
	for _, rg := range ranges {
		if room <= 0 {
			break
		}
		rows, _ := p.cut(s.rows, rg.r0, rg.r1, true, 0)
		cols, _ := p.cut(s.cols, rg.c0, rg.c1, false, 0)
		if len(rows) == 0 || len(cols) == 0 {
			continue
		}
		if n := len(rows) * len(cols); n > room {
			c.warnOnce("pages:"+s.ws.name, "sheet %q: its print areas take more than %d pages; the rest are left out", s.ws.name, maxPrintPages)
			if room < len(cols) {
				break
			}
			rows = rows[:room/len(cols)]
		}
		room -= len(rows) * len(cols)
		a := &printArea{rg: rg, rows: rows, cols: cols, body: map[[2]int]*printRegion{}, top: map[int]*printRegion{}, left: map[int]*printRegion{}}
		for _, band := range rows {
			a.hasRowTitles = a.hasRowTitles || p.rowTitles(band.i0)
		}
		a.rowPos, a.colPos = bandPos(s.rows, rows), bandPos(s.cols, cols)
		p.areas = append(p.areas, a)
		if ps.overThenDown {
			for i := range rows {
				for j := range cols {
					p.pages = append(p.pages, printPage{a, i, j})
				}
			}
		} else {
			for j := range cols {
				for i := range rows {
					p.pages = append(p.pages, printPage{a, i, j})
				}
			}
		}
	}
	return p
}

// bandPos returns where the bands start on the sheet, and where the last
// ends.
func bandPos(ax *axis, bands []printBand) []float64 {
	out := make([]float64, len(bands)+1)
	for k, b := range bands {
		out[k] = ax.pos(b.i0)
	}
	out[len(bands)] = ax.pos(bands[len(bands)-1].i1)
	return out
}

// availW and availH are the paper between the margins.
func (p *printPlan) availW() float64 { return p.pw - p.ml - p.mr }
func (p *printPlan) availH() float64 { return p.ph - p.mt - p.mb }

// rowTitles reports whether a page starting at first needs the title rows:
// after the original titles, and within the table when a scope is set.
func (p *printPlan) rowTitles(first int) bool {
	return p.titleRows != nil && first > p.titleRows[1] &&
		(p.titleScope == nil || first >= p.titleScope.r0 && first <= p.titleScope.r1)
}

// joinKeepRows joins overlapping groups, but leaves adjacent sections
// independent so each can move to the next page on its own.
func joinKeepRows(groups [][2]int) [][2]int {
	sort.Slice(groups, func(i, j int) bool { return groups[i][0] < groups[j][0] })
	out := groups[:0]
	for _, g := range groups {
		if g[0] < 0 || g[1] < g[0] || g[1] >= maxRows {
			continue
		}
		if n := len(out); n > 0 && g[0] <= out[n-1][1] {
			out[n-1][1] = max(out[n-1][1], g[1])
		} else {
			out = append(out, g)
		}
	}
	return out
}

// colTitles is rowTitles for the title columns.
func (p *printPlan) colTitles(first int) bool { return p.titleCols != nil && first > p.titleCols[1] }

// cut cuts the entries first to last of the rows (or the columns) into the
// bands of pages at the plan's scale: a band ends where the next entry
// does not fit on the page, and before an entry with a manual break. With
// a limit, it stops when the bands are more (ok is false).
func (p *printPlan) cut(ax *axis, first, last int, rows bool, limit int) (bands []printBand, ok bool) {
	avail, head, breaks := p.availW()/p.scale-p.hdrW, p.tx1-p.tx0, p.ps.colBreaks
	titles := p.colTitles
	if rows {
		avail, head, breaks, titles = p.availH()/p.scale-p.hdrH, p.ty1-p.ty0, p.ps.rowBreaks, p.rowTitles
	}
	room := func(start int) float64 {
		r := avail
		if titles(start) {
			r -= head
		}
		return r
	}
	bi := sort.SearchInts(breaks, first+1)
	ki := sort.Search(len(p.keepRows), func(k int) bool { return p.keepRows[k][0] >= first })
	start, used, limitAt := first, 0.0, room(first)
	for i := first; i <= last; i++ {
		size := ax.at(i)
		brk := false
		for bi < len(breaks) && breaks[bi] <= i {
			brk = brk || breaks[bi] == i
			bi++
		}
		if rows && ki < len(p.keepRows) && p.keepRows[ki][0] == i {
			end := p.keepRows[ki][1]
			ki++
			if end <= last {
				need := ax.pos(end+1) - ax.pos(i)
				// A group larger than a fresh page must be allowed to split.
				if need <= room(i)+0.01 && used+need > limitAt+0.01 {
					brk = true
				}
			}
		}
		if i > start && (brk || used+size > limitAt+0.01) {
			bands = append(bands, printBand{start, i})
			if limit > 0 && len(bands) >= limit {
				return bands, false
			}
			if len(bands) > maxPrintPages {
				return bands, true
			}
			start, used, limitAt = i, 0, room(i)
		}
		used += size
	}
	return append(bands, printBand{start, last + 1}), true
}

// fitScale returns the largest scale, in whole percents up to 100, at which
// every area is as many pages wide and tall as the page setup fits it to.
func (p *printPlan) fitScale(ranges []cellRange) float64 {
	// a page setup that does not say fits on a page each way
	fitW, fitH := 1, 1
	if p.ps.set {
		fitW, fitH = p.ps.fitW, p.ps.fitH
	}
	// the areas are cut at each scale that is tried: those of a hostile
	// file (hundreds of areas of a million rows) are not all looked at
	const maxFitEntries = 4 << 20
	for i, n := 0, 0; i < len(ranges); i++ {
		if n += ranges[i].r1 - ranges[i].r0 + ranges[i].c1 - ranges[i].c0 + 2; n > maxFitEntries {
			ranges = ranges[:i+1]
		}
	}
	fits := func(pct int) bool {
		p.scale = float64(pct) / 100
		for _, rg := range ranges {
			if fitH > 0 {
				if _, ok := p.cut(p.s.rows, rg.r0, rg.r1, true, fitH); !ok {
					return false
				}
			}
			if fitW > 0 {
				if _, ok := p.cut(p.s.cols, rg.c0, rg.c1, false, fitW); !ok {
					return false
				}
			}
		}
		return true
	}
	lo, hi := int(minPrintScale*100), 100
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if fits(mid) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return float64(lo) / 100
}

// bandsIn returns the bands k0 to k1-1 that the positions a to b touch.
func bandsIn(pos []float64, a, b float64) (k0, k1 int) {
	n := len(pos) - 1
	k0 = sort.Search(n, func(k int) bool { return pos[k+1] > a })
	k1 = sort.Search(n, func(k int) bool { return pos[k] >= b })
	return k0, max(k1, k0)
}

// target returns the object of a region's layer that is drawn into,
// starting it.
func (p *printPlan) target(r *printRegion) *tileCv {
	t := r.layers[p.layer]
	if t == nil {
		t = &tileCv{cv: p.s.c.cvs.New(), ox: r.b.x, oy: r.b.y, bounds: &r.b}
		t.cv.Obj.SetBBox(0, 0, f32(r.b.w), f32(r.b.h))
		r.layers[p.layer] = t
		p.all = append(p.all, t)
	}
	return t
}

// regionsIn returns the regions a rectangle of the sheet touches, as what
// is drawn into for them.
func (p *printPlan) regionsIn(x0, y0, x1, y1 float64) []*tileCv {
	if !(x0 < x1 && y0 < y1) {
		return nil
	}
	var out []*tileCv
	inRows := p.titleRows != nil && y0 < p.ty1 && y1 > p.ty0
	inCols := p.titleCols != nil && x0 < p.tx1 && x1 > p.tx0
	corner := false
	for _, a := range p.areas {
		rowT := a.hasRowTitles
		colT := p.colTitles(a.cols[len(a.cols)-1].i0)
		corner = corner || rowT && colT
		// most rectangles are outside most areas of a sheet with many
		inX := x0 < a.colPos[len(a.cols)] && x1 > a.colPos[0]
		inY := y0 < a.rowPos[len(a.rows)] && y1 > a.rowPos[0]
		if !(inX && inY) && !(inX && inRows && rowT) && !(inY && inCols && colT) {
			continue
		}
		i0, i1 := bandsIn(a.rowPos, y0, y1)
		j0, j1 := bandsIn(a.colPos, x0, x1)
		for i := i0; i < i1; i++ {
			for j := j0; j < j1; j++ {
				k := [2]int{i, j}
				r := a.body[k]
				if r == nil {
					r = &printRegion{b: box{a.colPos[j], a.rowPos[i], a.colPos[j+1] - a.colPos[j], a.rowPos[i+1] - a.rowPos[i]}}
					a.body[k] = r
				}
				if r.b.w > 0 && r.b.h > 0 {
					out = append(out, p.target(r))
				}
			}
		}
		if inRows && rowT {
			for j := j0; j < j1; j++ {
				r := a.top[j]
				if r == nil {
					r = &printRegion{b: box{a.colPos[j], p.ty0, a.colPos[j+1] - a.colPos[j], p.ty1 - p.ty0}}
					a.top[j] = r
				}
				if r.b.w > 0 && r.b.h > 0 {
					out = append(out, p.target(r))
				}
			}
		}
		if inCols && colT {
			for i := i0; i < i1; i++ {
				r := a.left[i]
				if r == nil {
					r = &printRegion{b: box{p.tx0, a.rowPos[i], p.tx1 - p.tx0, a.rowPos[i+1] - a.rowPos[i]}}
					a.left[i] = r
				}
				if r.b.w > 0 && r.b.h > 0 {
					out = append(out, p.target(r))
				}
			}
		}
	}
	if inRows && inCols && corner {
		if p.corner == nil {
			p.corner = &printRegion{b: box{p.tx0, p.ty0, p.tx1 - p.tx0, p.ty1 - p.ty0}}
		}
		if r := p.corner; r.b.w > 0 && r.b.h > 0 {
			out = append(out, p.target(r))
		}
	}
	return out
}

// clipSegment cuts a line to a box, edges included; ok is false when no
// part of it is in the box.
func clipSegment(sg segment, b *box) (segment, bool) {
	t0, t1 := 0.0, 1.0
	dx, dy := sg.x1-sg.x0, sg.y1-sg.y0
	for _, e := range [4][2]float64{{-dx, sg.x0 - b.x}, {dx, b.x + b.w - sg.x0}, {-dy, sg.y0 - b.y}, {dy, b.y + b.h - sg.y0}} {
		pk, qk := e[0], e[1]
		switch {
		case pk == 0:
			if qk < 0 {
				return sg, false
			}
		case pk < 0:
			t0 = math.Max(t0, qk/pk)
		default:
			t1 = math.Min(t1, qk/pk)
		}
	}
	if t0 > t1 {
		return sg, false
	}
	return segment{sg.x0 + t0*dx, sg.y0 + t0*dy, sg.x0 + t1*dx, sg.y0 + t1*dy}, true
}

// paintGridlines draws the lines between the cells the sheet uses, as
// Excel prints them when the page setup asks for them: under the fills, and
// not through merged ranges.
func (s *sheetCtx) paintGridlines() {
	p := s.pr
	if p.usedR < 0 || p.usedC < 0 {
		return
	}
	paths := map[*tileCv]*bdf.Path{}
	var order []*tileCv
	const w = pxPt / 2
	line := func(x0, y0, x1, y1 float64) {
		if x0 == x1 && y0 == y1 {
			return
		}
		for _, t := range s.tilesIn(math.Min(x0, x1)-w, math.Min(y0, y1)-w, math.Max(x0, x1)+w, math.Max(y0, y1)+w) {
			sg, ok := clipSegment(segment{x0, y0, x1, y1}, t.bounds)
			if !ok {
				continue
			}
			pa := paths[t]
			if pa == nil {
				pa = &bdf.Path{}
				paths[t] = pa
				order = append(order, t)
			}
			pa.MoveTo(f32(sg.x0-t.ox), f32(sg.y0-t.oy)).LineTo(f32(sg.x1-t.ox), f32(sg.y1-t.oy))
		}
	}
	// gaps cuts the entries a to b of a line into those outside the spans
	gaps := func(a, b int, spans [][2]int, fn func(a, b int)) {
		sort.Slice(spans, func(i, j int) bool { return spans[i][0] < spans[j][0] })
		for _, sp := range spans {
			if sp[0] > a {
				fn(a, min(sp[0], b+1)-1)
			}
			a = max(a, sp[1]+1)
			if a > b {
				return
			}
		}
		fn(a, b)
	}
	// the lines under the rows: not across the merged ranges that go on below
	for r := -1; r <= p.usedR; r++ {
		var spans [][2]int
		if r >= 0 {
			if s.rows.at(r) == 0 {
				continue
			}
			for _, i := range s.mergeRows[r] {
				if m := s.merges[i]; m.r1 > r {
					spans = append(spans, [2]int{m.c0, m.c1})
				}
			}
		}
		y := s.rows.pos(r + 1)
		gaps(0, p.usedC, spans, func(a, b int) { line(s.cols.pos(a), y, s.cols.pos(b+1), y) })
	}
	// the lines beside the columns: not across the merged ranges that go on
	// to the right
	byCol := map[int][][2]int{}
	for _, m := range s.merges {
		for c := m.c0; c < m.c1 && c <= p.usedC; c++ {
			byCol[c] = append(byCol[c], [2]int{m.r0, m.r1})
		}
	}
	for c := -1; c <= p.usedC; c++ {
		if c >= 0 && s.cols.at(c) == 0 {
			continue
		}
		x := s.cols.pos(c + 1)
		gaps(0, p.usedR, byCol[c], func(a, b int) { line(x, s.rows.pos(a), x, s.rows.pos(b+1)) })
	}
	for _, t := range order {
		o := t.cv.Obj
		o.Line(f32(w), bdf.CapButt, bdf.JoinMiter, 10)
		o.StrokeColor(bdf.RGB(0, 0, 0))
		o.StrokePath(o.AddPath(paths[t]))
		t.cv.Drawn = true
	}
}

// place draws a region at a place on the page, at the plan's scale.
func (p *printPlan) place(cv *canvas.Canvas, r *printRegion, x, y float64) {
	if r == nil {
		return
	}
	for layer, t := range r.layers {
		if t == nil || !t.cv.Drawn {
			continue
		}
		// each object starts in the document's language
		cv.SetLang("")
		o := cv.Obj
		o.Save()
		if layer != 1 {
			// the borders are cut to the region line by line: the lines
			// on its edges are drawn whole
			o.ClipRect(f32(x), f32(y), f32(r.b.w*p.scale), f32(r.b.h*p.scale))
		}
		cv.Transform(canvas.Matrix{p.scale, 0, 0, p.scale, x, y})
		bbox := bdf.Rect{W: f32(r.b.w), H: f32(r.b.h)}
		o.Use(cv.Share(t.cv, bbox))
		cv.Used(t.cv)
		o.Restore()
		cv.Drawn = true
	}
}

// view makes the view of the plan's pages: on each, the header, the footer
// and the cells. It returns the objects of the layers of each page.
func (p *printPlan) view(id, title string) (*bdf.View, [][]*canvas.Canvas) {
	c := p.s.c
	v := c.doc.NewView(id, bdf.ViewFlow, title)
	v.Continuous = &bdf.Continuous{Gap: 24}
	body := &bdf.RectDef{X: f32(p.ml), Y: f32(p.mt), W: f32(p.availW()), H: f32(p.availH())}
	var out [][]*canvas.Canvas
	first := 1
	if p.ps.useFirstPage {
		first = p.ps.firstPage
	}
	count := len(p.pages)
	add := func(n int, cells *canvas.Canvas) {
		bp := v.AddPage(f32(p.pw), f32(p.ph))
		bp.Body = body
		var cvs []*canvas.Canvas
		layer := func(role string, cv *canvas.Canvas, always bool) {
			if cv != nil && (cv.Drawn || always) {
				bp.Layers = append(bp.Layers, bdf.Layer{Role: role})
				cvs = append(cvs, cv)
			}
		}
		head, foot := p.ps.hf.codes(n, first)
		layer(bdf.RoleHeader, p.hfLayer(head, true, n, first+count-1), false)
		layer(bdf.RoleFooter, p.hfLayer(foot, false, n, first+count-1), false)
		layer(bdf.RoleBody, cells, true)
		out = append(out, cvs)
	}
	if count == 0 {
		// a sheet with nothing to print: Excel prints nothing, a view has a page
		count = 1
		add(first, c.cvs.New())
		return v, out
	}
	for k, pg := range p.pages {
		a := pg.a
		cv := c.cvs.New()
		rowT, colT := p.rowTitles(a.rows[pg.i].i0), p.colTitles(a.cols[pg.j].i0)
		w, h := a.colPos[pg.j+1]-a.colPos[pg.j]+p.hdrW, a.rowPos[pg.i+1]-a.rowPos[pg.i]+p.hdrH
		if colT {
			w += p.tx1 - p.tx0
		}
		if rowT {
			h += p.ty1 - p.ty0
		}
		x0, y0 := p.ml, p.mt
		if p.ps.hCenter {
			x0 += math.Max(p.availW()-w*p.scale, 0) / 2
		}
		if p.ps.vCenter {
			y0 += math.Max(p.availH()-h*p.scale, 0) / 2
		}
		hx, hy := x0+p.hdrW*p.scale, y0+p.hdrH*p.scale
		bx, by := hx, hy
		if colT {
			bx += (p.tx1 - p.tx0) * p.scale
		}
		if rowT {
			by += (p.ty1 - p.ty0) * p.scale
		}
		if rowT && colT {
			p.place(cv, p.corner, hx, hy)
		}
		if rowT {
			p.place(cv, a.top[pg.j], bx, hy)
		}
		if colT {
			p.place(cv, a.left[pg.i], hx, by)
		}
		p.place(cv, a.body[[2]int{pg.i, pg.j}], bx, by)
		if p.ps.headings {
			p.headings(cv, pg, x0, y0, rowT, colT)
		}
		add(first+k, cv)
	}
	return v, out
}

// headings draws the letters of the columns and the numbers of the rows of
// a page, along its top and its left.
func (p *printPlan) headings(cv *canvas.Canvas, pg printPage, x0, y0 float64, rowT, colT bool) {
	s := p.s
	st := s.textStyle(s.c.st.font(0), nil, nil)
	st.color = bdf.RGB(0, 0, 0)
	t := &tileCv{cv: cv}
	o := cv.Obj
	o.Save()
	cv.Transform(canvas.Matrix{p.scale, 0, 0, p.scale, x0, y0})
	cv.SetLang("")
	lines := &bdf.Path{}
	label := func(text string, x, y, w, h float64) {
		ln := &tline{}
		for _, r := range text {
			ln.items = append(ln.items, s.item(st, r))
		}
		lineMetrics(ln, st, s)
		ln.x, ln.base = x+(w-ln.width)/2, y+(h+ln.asc-ln.desc)/2
		o.Mark(bdf.MarkParagraph, "")
		s.emitLine(t, ln, 0, 0)
		lines.MoveTo(f32(x), f32(y)).LineTo(f32(x+w), f32(y)).LineTo(f32(x+w), f32(y+h)).LineTo(f32(x), f32(y+h)).Close()
	}
	x := p.hdrW
	cols := func(c0, c1 int) {
		for c := c0; c < c1; c++ {
			if w := s.cols.at(c); w > 0 {
				label(colName(c), x, 0, w, p.hdrH)
				x += w
			}
		}
	}
	if colT {
		cols(p.titleCols[0], p.titleCols[1]+1)
	}
	cols(pg.a.cols[pg.j].i0, pg.a.cols[pg.j].i1)
	y := p.hdrH
	rows := func(r0, r1 int) {
		for r := r0; r < r1; r++ {
			if h := s.rows.at(r); h > 0 {
				label(strconv.Itoa(r+1), 0, y, p.hdrW, h)
				y += h
			}
		}
	}
	if rowT {
		rows(p.titleRows[0], p.titleRows[1]+1)
	}
	rows(pg.a.rows[pg.i].i0, pg.a.rows[pg.i].i1)
	o.Line(pxPt/2, bdf.CapButt, bdf.JoinMiter, 10)
	o.StrokeColor(bdf.RGB(0, 0, 0))
	o.StrokePath(o.AddPath(lines))
	o.Restore()
	cv.Drawn = true
}
