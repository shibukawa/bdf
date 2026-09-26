package drawingml

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// Tables: the grid comes from a:tblGrid and the rows; cell fills, borders
// and text properties come from the cell's own a:tcPr and from the table
// style parts that apply to it (whole table, banded rows and columns,
// first/last row and column, corner cells), later parts winning. Rows grow
// to fit their text.

func (c *Renderer) tableStyle(tblPr *ooxml.Node) *ooxml.Node {
	if st := tblPr.Child("tableStyle"); st != nil {
		return st
	}
	id := strings.ToUpper(strings.TrimSpace(tblPr.Child("tableStyleId").Content()))
	if id == "" {
		// no style: plain table
		return nil
	}
	if st, ok := c.tableStyles[id]; ok {
		return st
	}
	if st := builtinTableStyle(id); st != nil {
		c.tableStyles[id] = st
		return st
	}
	c.warnOnce("tblstyle:"+id, "unknown table style %s; using the default table style", id)
	st := builtinTableStyle("{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}")
	c.tableStyles[id] = st
	return st
}

type tcell struct {
	tc         *ooxml.Node
	r, c       int
	rows, cols int
	merged     bool // covered by another cell's span
	block      *textBlock
}

// region is the rectangle of cells a table style part covers for a cell.
type region struct{ r0, c0, r1, c1 int }

type stylePart struct {
	n   *ooxml.Node
	reg region
}

func (s *Drawing) drawTable(cv *canvas.Canvas, sh *shape, xf xform, tbl *ooxml.Node) {
	var colW []float64
	for _, gc := range tbl.Path("tblGrid").Children("gridCol") {
		colW = append(colW, gc.AttrEMU("w", 0))
	}
	trs := tbl.Children("tr")
	nr, nc := len(trs), len(colW)
	if nr == 0 || nc == 0 {
		return
	}
	rowH := make([]float64, nr)
	cells := make([][]*tcell, nr)
	for r, tr := range trs {
		rowH[r] = tr.AttrEMU("h", 0)
		cells[r] = make([]*tcell, nc)
		for c, tc := range tr.Children("tc") {
			if c >= nc {
				break
			}
			cells[r][c] = &tcell{tc: tc, r: r, c: c,
				rows: max(1, int(tc.AttrInt("rowSpan", 1))), cols: max(1, int(tc.AttrInt("gridSpan", 1))),
				merged: tc.AttrBool("hMerge", false) || tc.AttrBool("vMerge", false)}
		}
	}
	tblPr := tbl.Child("tblPr")
	style := s.c.tableStyle(tblPr)
	flags := map[string]bool{}
	for _, f := range []string{"firstRow", "lastRow", "firstCol", "lastCol", "bandRow", "bandCol"} {
		flags[f] = tblPr.AttrBool(f, false)
	}
	// style parts that apply to a cell, in application order
	parts := func(r, c int) []stylePart {
		if style == nil {
			return nil
		}
		whole := region{0, 0, nr - 1, nc - 1}
		out := []stylePart{{style.Child("wholeTbl"), whole}}
		if flags["bandCol"] {
			first := 0
			if flags["firstCol"] {
				first = 1
			}
			if c >= first && !(flags["lastCol"] && c == nc-1) {
				name := "band1V"
				if (c-first)%2 == 1 {
					name = "band2V"
				}
				out = append(out, stylePart{style.Child(name), region{0, c, nr - 1, c}})
			}
		}
		if flags["bandRow"] {
			first := 0
			if flags["firstRow"] {
				first = 1
			}
			if r >= first && !(flags["lastRow"] && r == nr-1) {
				name := "band1H"
				if (r-first)%2 == 1 {
					name = "band2H"
				}
				out = append(out, stylePart{style.Child(name), region{r, 0, r, nc - 1}})
			}
		}
		if flags["lastCol"] && c == nc-1 {
			out = append(out, stylePart{style.Child("lastCol"), region{0, c, nr - 1, c}})
		}
		if flags["firstCol"] && c == 0 {
			out = append(out, stylePart{style.Child("firstCol"), region{0, 0, nr - 1, 0}})
		}
		if flags["lastRow"] && r == nr-1 {
			out = append(out, stylePart{style.Child("lastRow"), region{r, 0, r, nc - 1}})
		}
		if flags["firstRow"] && r == 0 {
			out = append(out, stylePart{style.Child("firstRow"), region{0, 0, 0, nc - 1}})
		}
		cell := region{r, c, r, c}
		if flags["lastRow"] && r == nr-1 && flags["firstCol"] && c == 0 {
			out = append(out, stylePart{style.Child("swCell"), cell})
		}
		if flags["lastRow"] && r == nr-1 && flags["lastCol"] && c == nc-1 {
			out = append(out, stylePart{style.Child("seCell"), cell})
		}
		if flags["firstRow"] && r == 0 && flags["firstCol"] && c == 0 {
			out = append(out, stylePart{style.Child("nwCell"), cell})
		}
		if flags["firstRow"] && r == 0 && flags["lastCol"] && c == nc-1 {
			out = append(out, stylePart{style.Child("neCell"), cell})
		}
		return out
	}
	colX := make([]float64, nc+1)
	for c := 0; c < nc; c++ {
		colX[c+1] = colX[c] + colW[c]
	}
	// lay out text and grow rows
	for r := 0; r < nr; r++ {
		for c := 0; c < nc; c++ {
			cl := cells[r][c]
			if cl == nil || cl.merged {
				continue
			}
			tf := s.cellText(cl.tc, parts(r, c), sh.part)
			if tf == nil {
				continue
			}
			c1 := min(nc, c+cl.cols)
			w := colX[c1] - colX[c]
			cl.block = s.layoutText(tf, 0, 0, w, rowH[r], canvas.Identity)
			if cl.block == nil || cl.block.vertical {
				continue
			}
			need := cl.block.height()
			r1 := min(nr, r+cl.rows)
			have := 0.0
			for k := r; k < r1; k++ {
				have += rowH[k]
			}
			if need > have {
				rowH[r1-1] += need - have
			}
		}
	}
	rowY := make([]float64, nr+1)
	for r := 0; r < nr; r++ {
		rowY[r+1] = rowY[r] + rowH[r]
	}
	m := canvas.Translate(xf.X, xf.Y)
	cv.Obj.Save()
	cv.Transform(m)
	// table background
	if bg := fillElem(tblPr); bg != nil {
		f := s.resolveFill(bg, sh.part, s.cc, nil)
		if f.kind != fillNone && s.setFill(cv, f, colX[nc], rowY[nr]) {
			cv.Obj.FillRect(0, 0, f32(colX[nc]), f32(rowY[nr]))
			cv.Drawn = true
		}
	} else if style != nil {
		if bg := style.Child("tblBg"); bg != nil {
			f := s.resolveFill(fillElem(bg), "", s.cc, nil)
			if ref := bg.Child("fillRef"); ref != nil {
				f = s.styleFill(ref, s.cc)
			}
			if f.kind != fillNone && s.setFill(cv, f, colX[nc], rowY[nr]) {
				cv.Obj.FillRect(0, 0, f32(colX[nc]), f32(rowY[nr]))
				cv.Drawn = true
			}
		}
	}
	// cell fills
	for r := 0; r < nr; r++ {
		for c := 0; c < nc; c++ {
			cl := cells[r][c]
			if cl == nil || cl.merged {
				continue
			}
			x0, y0 := colX[c], rowY[r]
			x1, y1 := colX[min(nc, c+cl.cols)], rowY[min(nr, r+cl.rows)]
			f := s.cellFill(cl.tc, parts(r, c), sh.part)
			if f.kind == fillBlip {
				s.drawBlip(cv, f.blip, f.part, f.cc, x0, y0, x1-x0, y1-y0, false)
			} else if f.kind != fillNone && s.setFill(cv, f, x1-x0, y1-y0) {
				cv.Obj.Save()
				cv.Obj.Translate(f32(x0), f32(y0))
				cv.Obj.FillRect(0, 0, f32(x1-x0), f32(y1-y0))
				cv.Obj.Restore()
				cv.Drawn = true
			}
		}
	}
	cv.Obj.Restore()
	// text, in a TABLE with a CELL for each cell that has text (in row-major
	// order of their top left corners)
	inTable := false
	for r := 0; r < nr; r++ {
		for c := 0; c < nc; c++ {
			cl := cells[r][c]
			if cl == nil || cl.merged || cl.block == nil {
				continue
			}
			x0, y0 := colX[c], rowY[r]
			x1, y1 := colX[min(nc, c+cl.cols)], rowY[min(nr, r+cl.rows)]
			b := cl.block
			if b.vertical {
				tf := s.cellText(cl.tc, parts(r, c), sh.part)
				b = s.layoutText(tf, 0, 0, x1-x0, y1-y0, canvas.Identity)
				if b == nil {
					continue
				}
			} else {
				b.h = y1 - y0
			}
			b.fm = m.Mul(canvas.Translate(x0, y0)).Mul(b.fm)
			if !inTable {
				cv.Obj.Mark(bdf.MarkTable, "")
				inTable = true
			}
			ref := cellRef(r, c)
			if r1, c1 := min(nr, r+cl.rows)-1, min(nc, c+cl.cols)-1; r1 > r || c1 > c {
				ref += ":" + cellRef(r1, c1)
			}
			switch {
			case flags["firstRow"] && r == 0:
				ref += " col"
			case flags["firstCol"] && c == 0:
				ref += " row"
			}
			cv.Obj.Mark(bdf.MarkCell, ref)
			s.drawTextBlock(cv, b)
		}
	}
	if inTable {
		cv.Obj.Mark(bdf.MarkEnd, "")
	}
	// borders: each edge once, the cell's own line winning over the style
	cv.Obj.Save()
	cv.Transform(m)
	edge := func(r, c int, side string) *line {
		cl := cells[r][c]
		if cl == nil {
			return nil
		}
		own := map[string]string{"left": "lnL", "right": "lnR", "top": "lnT", "bottom": "lnB"}[side]
		if ln := cl.tc.Path("tcPr", own); ln != nil {
			return s.resolveLine([]*ooxml.Node{ln}, []string{sh.part}, nil, s.cc)
		}
		var found *ooxml.Node
		for _, p := range parts(r, c) {
			bdr := p.n.Path("tcStyle", "tcBdr")
			if bdr == nil {
				continue
			}
			name := side
			switch side {
			case "left":
				if c > p.reg.c0 {
					name = "insideV"
				}
			case "right":
				if c+cl.cols-1 < p.reg.c1 {
					name = "insideV"
				}
			case "top":
				if r > p.reg.r0 {
					name = "insideH"
				}
			case "bottom":
				if r+cl.rows-1 < p.reg.r1 {
					name = "insideH"
				}
			}
			if b := bdr.Child(name); b != nil {
				found = b
			}
		}
		if found == nil {
			return nil
		}
		if ln := found.Child("ln"); ln != nil {
			return s.resolveLine([]*ooxml.Node{ln}, []string{""}, nil, s.cc)
		}
		if ref := found.Child("lnRef"); ref != nil {
			return s.resolveLine(nil, nil, ref, s.cc)
		}
		return nil
	}
	stroke := func(l *line, x0, y0, x1, y1 float64) {
		if l == nil {
			return
		}
		cv.Obj.Save()
		setLine(cv, l, math.Abs(x1-x0), math.Abs(y1-y0))
		p := newPath().moveTo(x0, y0).lineTo(x1, y1)
		cv.Obj.StrokePath(cv.Obj.AddPath(p.p))
		cv.Obj.Restore()
		cv.Drawn = true
	}
	owner := func(r, c int) *tcell {
		// the origin cell covering (r, c)
		for rr := r; rr >= 0; rr-- {
			for cc := c; cc >= 0; cc-- {
				if cl := cells[rr][cc]; cl != nil && !cl.merged && rr+cl.rows > r && cc+cl.cols > c {
					return cl
				}
			}
		}
		return nil
	}
	for r := 0; r < nr; r++ {
		for c := 0; c < nc; c++ {
			o := owner(r, c)
			if o == nil {
				continue
			}
			// top edge of (r, c): between the cell above and this one
			if r == 0 || owner(r-1, c) != o {
				l := edge(o.r, o.c, "top")
				if r > 0 {
					if above := owner(r-1, c); above != nil {
						if l2 := edge(above.r, above.c, "bottom"); l == nil || (l2 != nil && above.tc.Path("tcPr", "lnB") != nil && o.tc.Path("tcPr", "lnT") == nil) {
							l = l2
						}
					}
				}
				stroke(l, colX[c], rowY[r], colX[c+1], rowY[r])
			}
			if c == 0 || owner(r, c-1) != o {
				l := edge(o.r, o.c, "left")
				if c > 0 {
					if left := owner(r, c-1); left != nil {
						if l2 := edge(left.r, left.c, "right"); l == nil || (l2 != nil && left.tc.Path("tcPr", "lnR") != nil && o.tc.Path("tcPr", "lnL") == nil) {
							l = l2
						}
					}
				}
				stroke(l, colX[c], rowY[r], colX[c], rowY[r+1])
			}
			if r == nr-1 {
				stroke(edge(o.r, o.c, "bottom"), colX[c], rowY[nr], colX[c+1], rowY[nr])
			}
			if c == nc-1 {
				stroke(edge(o.r, o.c, "right"), colX[nc], rowY[r], colX[nc], rowY[r+1])
			}
		}
	}
	cv.Obj.Restore()
}

// cellText builds the text frame of a cell: its body, margins as insets,
// and the text properties of the table style parts.
func (s *Drawing) cellText(tc *ooxml.Node, parts []stylePart, part string) *textFrame {
	tb := tc.Child("txBody")
	if tb == nil {
		return nil
	}
	tcPr := tc.Child("tcPr")
	bp := &ooxml.Node{Name: "bodyPr"}
	for _, kv := range [][2]string{{"marL", "lIns"}, {"marR", "rIns"}, {"marT", "tIns"}, {"marB", "bIns"}} {
		if v, ok := tcPr.Attr(kv[0]); ok {
			bp.Attrs = append(bp.Attrs, mkNode("", kv[1], v).Attrs...)
		}
	}
	if v, ok := tcPr.Attr("anchor"); ok {
		bp.Attrs = append(bp.Attrs, mkNode("", "anchor", v).Attrs...)
	}
	if v, ok := tcPr.Attr("vert"); ok {
		bp.Attrs = append(bp.Attrs, mkNode("", "vert", v).Attrs...)
	}
	tf := &textFrame{body: tb, part: part, cc: s.cc, own: tb.Child("lstStyle"),
		bodyPrs: chain{bp, tb.Child("bodyPr")}, lists: []*ooxml.Node{s.host.TextStyle("")}}
	for i := len(parts) - 1; i >= 0; i-- {
		if ts := parts[i].n.Child("tcTxStyle"); ts != nil {
			tf.extra = append(tf.extra, tcTxProps(ts))
		}
	}
	return tf
}

// tcTxProps turns a table style's a:tcTxStyle into run properties.
func tcTxProps(ts *ooxml.Node) *ooxml.Node {
	n := &ooxml.Node{Name: "defRPr"}
	for _, a := range []string{"b", "i"} {
		switch ts.AttrStr(a, "def") {
		case "on":
			n.Attrs = append(n.Attrs, mkNode("", a, "1").Attrs...)
		case "off":
			n.Attrs = append(n.Attrs, mkNode("", a, "0").Attrs...)
		}
	}
	if fr := ts.Child("fontRef"); fr != nil {
		ref := fontRefProps(fr)
		for _, k := range ref.Kids {
			if k.Name != "solidFill" {
				n.Kids = append(n.Kids, k)
			}
		}
	}
	if f := ts.Child("font"); f != nil {
		n.Kids = append(n.Kids, f.Kids...)
	}
	for _, k := range ts.Kids {
		if isColor(k) {
			n.Kids = append(n.Kids, &ooxml.Node{Name: "solidFill", Kids: []*ooxml.Node{k}})
		}
	}
	return n
}

func (s *Drawing) cellFill(tc *ooxml.Node, parts []stylePart, part string) fill {
	if e := fillElem(tc.Child("tcPr")); e != nil {
		return s.resolveFill(e, part, s.cc, nil)
	}
	f := fill{}
	for _, p := range parts {
		ts := p.n.Child("tcStyle")
		if ts == nil {
			continue
		}
		if e := fillElem(ts.Child("fill")); e != nil {
			f = s.resolveFill(e, "", s.cc, nil)
		} else if ref := ts.Child("fillRef"); ref != nil {
			f = s.styleFill(ref, s.cc)
		}
	}
	return f
}

// cellRef returns the A1-style reference of a cell (0-based row and column).
func cellRef(r, c int) string {
	var col []byte
	for c++; c > 0; c = (c - 1) / 26 {
		col = append([]byte{byte('A' + (c-1)%26)}, col...)
	}
	return string(col) + itoa(r+1)
}
