package wordproc

import (
	"slices"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// maxSpan limits colspan (and the columns a table may have).
const maxSpan = 1000

// table reads a table: its cells take the slots of the grid their row and
// column spans make, and the columns get the widths their content asks for
// (autoColumns). A caption goes above the table, inside its structure.
func (r *htmlReader) table(n *html.Node, css map[string]string, st *hstyle) {
	if st.tables >= maxTables {
		r.block(r.em(blockGap), r.em(blockGap), func() { r.children(n, st.block()) })
		return
	}
	r.flush()
	line := &border{style: "single", width: 0.75, color: ruleColor}
	t := &table{jc: "left", ind: st.indL, rowBand: 1, colBand: 1, frame: &frame{kind: fTable},
		cellMar: [4]float64{r.em(cellPadV), r.em(cellPadH), r.em(cellPadV), r.em(cellPadH)},
		borders: [6]*border{line, line, line, line, line, line},
		before:  r.em(blockGap), after: r.em(blockGap)}
	if a := strings.ToLower(attrStr(n, "align")); a == "center" || a == "right" {
		t.jc = a
	}
	type trRef struct {
		n    *html.Node
		head bool
	}
	var trs []trRef
	var capNode *html.Node
	for k := n.FirstChild; k != nil; k = k.NextSibling {
		if k.Type != html.ElementNode || hidden(k, cssDecls(k)) {
			continue
		}
		switch k.DataAtom {
		case atom.Caption:
			if capNode == nil {
				capNode = k
			}
		case atom.Thead, atom.Tbody, atom.Tfoot:
			for tr := k.FirstChild; tr != nil; tr = tr.NextSibling {
				if tr.Type == html.ElementNode && tr.DataAtom == atom.Tr && !hidden(tr, cssDecls(tr)) {
					trs = append(trs, trRef{tr, k.DataAtom == atom.Thead})
				}
			}
		case atom.Tr:
			trs = append(trs, trRef{k, false})
		}
	}
	if len(trs) == 0 {
		return
	}
	// the grid: the slots cells of rows above take
	type cellSrc struct {
		ce *cell
		n  *html.Node
		ri int
	}
	var refs []cellSrc
	taken := map[[2]int]bool{}
	ncols := 0
	head := true // still in the header rows at the top
	for ri, tr := range trs {
		rw := &row{}
		col := 0
		allTH, any := true, false
		for k := tr.n.FirstChild; k != nil; k = k.NextSibling {
			if k.Type != html.ElementNode || k.DataAtom != atom.Td && k.DataAtom != atom.Th || hidden(k, cssDecls(k)) {
				continue
			}
			for taken[[2]int{ri, col}] {
				col++
			}
			if col >= maxSpan {
				break
			}
			span := min(max(atoiAttr(k, "colspan", 1), 1), maxSpan-col)
			rs := atoiAttr(k, "rowspan", 1)
			if rs <= 0 {
				rs = len(trs) - ri
			}
			rs = min(max(rs, 1), len(trs)-ri)
			ce := &cell{col: col, span: span, rowSpan: rs, vAlign: valign(k), mar: t.cellMar}
			if rs > 1 {
				ce.vmerge = "restart"
			}
			for dr := range rs {
				for dc := range span {
					taken[[2]int{ri + dr, col + dc}] = true
				}
			}
			rw.cells = append(rw.cells, ce)
			refs = append(refs, cellSrc{ce, k, ri})
			any = true
			allTH = allTH && k.DataAtom == atom.Th
			col += span
			ncols = max(ncols, col)
		}
		head = head && (tr.head || any && allTH)
		rw.header = head
		t.rows = append(t.rows, rw)
	}
	if ncols == 0 {
		return
	}
	// the slots of cells spanning rows, in the rows below
	for _, ref := range refs {
		for dr := 1; dr < ref.ce.rowSpan; dr++ {
			rw := t.rows[ref.ri+dr]
			rw.cells = append(rw.cells, &cell{col: ref.ce.col, span: ref.ce.span, rowSpan: 1, vmerge: "continue", mar: t.cellMar})
		}
	}
	for _, rw := range t.rows {
		slices.SortStableFunc(rw.cells, func(a, b *cell) int { return a.col - b.col })
	}
	// content
	s := st.block()
	s.indL, s.listDepth, s.jc = 0, 0, "left"
	s.tables++
	for _, ref := range refs {
		ce, k := ref.ce, ref.n
		rw := t.rows[ref.ri]
		name := cellRef(ref.ri, ce.col)
		if ce.span > 1 || ce.rowSpan > 1 {
			name += ":" + cellRef(ref.ri+ce.rowSpan-1, ce.col+ce.span-1)
		}
		cs := s.child()
		th := k.DataAtom == atom.Th
		switch {
		case rw.header:
			name += " col"
		case th:
			name += " row"
		}
		if th {
			cs = s.chars()
			cs.rp.setToggle(tB, true)
			cs.jc = "center"
		}
		if rw.header {
			ce.shd, ce.hasShd = headerCellBG, true
		}
		if a := align(k, cssDecls(k)); a != "" {
			cs.jc = a
		}
		ce.frame = &frame{kind: fCell, payload: name}
		ce.ctx = st.inFrames(t.frame, ce.frame)
		cs.frames = ce.ctx
		ce.blocks = r.read(k, cs, false)
	}
	t.grid = make([]float64, ncols)
	t.auto = autoLayout(t)
	w := attrStr(n, "width")
	if cw := css["width"]; cw != "" {
		w = cw
	}
	if strings.HasSuffix(w, "%") {
		t.auto.fill = true
	} else if v := r.length(w); v > 0 {
		t.auto.width = v
	}
	if capNode != nil {
		cs := s.chars()
		cs.rp.sz *= 0.875
		cs.rp.color = mutedColor
		cs.jc = "center"
		cs.frames = st.inFrames(t.frame)
		for _, b := range r.read(capNode, cs, false) {
			if p, ok := b.(*para); ok {
				p.pp.keepNext = true
			}
			r.add(b)
		}
	}
	r.add(t)
}

func atoiAttr(n *html.Node, name string, def int) int {
	v, err := strconv.Atoi(attrStr(n, name))
	if err != nil {
		return def
	}
	return v
}

// valign reads the vertical alignment of a cell (middle by default, as in
// browsers).
func valign(n *html.Node) string {
	v := strings.ToLower(attrStr(n, "valign"))
	if css := cssDecls(n); css["vertical-align"] != "" {
		v = css["vertical-align"]
	}
	switch v {
	case "top", "baseline", "text-top":
		return "top"
	case "bottom", "text-bottom":
		return "bottom"
	}
	return "center"
}
