package pptx

import (
	"math"
	"strings"
)

// Tables: the grid comes from a:tblGrid and the rows; cell fills, borders
// and text properties come from the cell's own a:tcPr and from the table
// style parts that apply to it (whole table, banded rows and columns,
// first/last row and column, corner cells), later parts winning. Rows grow
// to fit their text.

func (c *converter) loadTableStyles() {
	c.tableStyles = map[string]*node{}
	if r, ok := c.pkg.relOfType(c.presPart, "/tableStyles"); ok {
		if n, err := c.pkg.xml(r.Target); err == nil {
			for _, st := range n.children("tblStyle") {
				c.tableStyles[strings.ToUpper(st.attrStr("styleId", ""))] = st
			}
		}
	}
}

// builtinTableStyle is PowerPoint's default table style, "Medium Style 2 -
// Accent 1", used when a table names a style the file does not define.
const builtinTableStyle = `<tblStyle styleId="{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}">
<wholeTbl><tcTxStyle><fontRef idx="minor"><prstClr val="black"/></fontRef><schemeClr val="dk1"/></tcTxStyle>
<tcStyle><tcBdr>
<left><ln w="12700"><solidFill><schemeClr val="lt1"/></solidFill></ln></left>
<right><ln w="12700"><solidFill><schemeClr val="lt1"/></solidFill></ln></right>
<top><ln w="12700"><solidFill><schemeClr val="lt1"/></solidFill></ln></top>
<bottom><ln w="12700"><solidFill><schemeClr val="lt1"/></solidFill></ln></bottom>
<insideH><ln w="12700"><solidFill><schemeClr val="lt1"/></solidFill></ln></insideH>
<insideV><ln w="12700"><solidFill><schemeClr val="lt1"/></solidFill></ln></insideV>
</tcBdr><fill><solidFill><schemeClr val="accent1"><tint val="20000"/></schemeClr></solidFill></fill></tcStyle></wholeTbl>
<band1H><tcStyle><tcBdr/><fill><solidFill><schemeClr val="accent1"><tint val="40000"/></schemeClr></solidFill></fill></tcStyle></band1H>
<band2H><tcStyle><tcBdr/></tcStyle></band2H>
<band1V><tcStyle><tcBdr/><fill><solidFill><schemeClr val="accent1"><tint val="40000"/></schemeClr></solidFill></fill></tcStyle></band1V>
<band2V><tcStyle><tcBdr/></tcStyle></band2V>
<lastCol><tcTxStyle b="on"><fontRef idx="minor"><prstClr val="black"/></fontRef><schemeClr val="lt1"/></tcTxStyle><tcStyle><tcBdr/><fill><solidFill><schemeClr val="accent1"/></solidFill></fill></tcStyle></lastCol>
<firstCol><tcTxStyle b="on"><fontRef idx="minor"><prstClr val="black"/></fontRef><schemeClr val="lt1"/></tcTxStyle><tcStyle><tcBdr/><fill><solidFill><schemeClr val="accent1"/></solidFill></fill></tcStyle></firstCol>
<lastRow><tcTxStyle b="on"><fontRef idx="minor"><prstClr val="black"/></fontRef><schemeClr val="lt1"/></tcTxStyle><tcStyle><tcBdr><top><ln w="38100"><solidFill><schemeClr val="lt1"/></solidFill></ln></top></tcBdr><fill><solidFill><schemeClr val="accent1"/></solidFill></fill></tcStyle></lastRow>
<firstRow><tcTxStyle b="on"><fontRef idx="minor"><prstClr val="black"/></fontRef><schemeClr val="lt1"/></tcTxStyle><tcStyle><tcBdr><bottom><ln w="38100"><solidFill><schemeClr val="lt1"/></solidFill></ln></bottom></tcBdr><fill><solidFill><schemeClr val="accent1"/></solidFill></fill></tcStyle></firstRow>
</tblStyle>`

func (c *converter) tableStyle(tblPr *node) *node {
	if st := tblPr.child("tableStyle"); st != nil {
		return st
	}
	id := strings.ToUpper(strings.TrimSpace(tblPr.child("tableStyleId").text()))
	if id == "" {
		// no style: plain table
		return nil
	}
	if st, ok := c.tableStyles[id]; ok {
		return st
	}
	n, _ := parseXML([]byte(builtinTableStyle))
	if id != "{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}" {
		c.warnOnce("tblstyle:"+id, "table style %s is not defined in the file; using the default table style", id)
	}
	return n
}

type tcell struct {
	tc         *node
	r, c       int
	rows, cols int
	merged     bool // covered by another cell's span
	block      *textBlock
}

// region is the rectangle of cells a table style part covers for a cell.
type region struct{ r0, c0, r1, c1 int }

type stylePart struct {
	n   *node
	reg region
}

func (s *slideCtx) drawTable(cv *canvas, sh *shape, xf xform, tbl *node) {
	var colW []float64
	for _, gc := range tbl.path("tblGrid").children("gridCol") {
		colW = append(colW, gc.emuAttr("w", 0))
	}
	trs := tbl.children("tr")
	nr, nc := len(trs), len(colW)
	if nr == 0 || nc == 0 {
		return
	}
	rowH := make([]float64, nr)
	cells := make([][]*tcell, nr)
	for r, tr := range trs {
		rowH[r] = tr.emuAttr("h", 0)
		cells[r] = make([]*tcell, nc)
		for c, tc := range tr.children("tc") {
			if c >= nc {
				break
			}
			cells[r][c] = &tcell{tc: tc, r: r, c: c,
				rows: max(1, int(tc.attrInt("rowSpan", 1))), cols: max(1, int(tc.attrInt("gridSpan", 1))),
				merged: tc.attrBool("hMerge", false) || tc.attrBool("vMerge", false)}
		}
	}
	tblPr := tbl.child("tblPr")
	style := s.c.tableStyle(tblPr)
	flags := map[string]bool{}
	for _, f := range []string{"firstRow", "lastRow", "firstCol", "lastCol", "bandRow", "bandCol"} {
		flags[f] = tblPr.attrBool(f, false)
	}
	// style parts that apply to a cell, in application order
	parts := func(r, c int) []stylePart {
		if style == nil {
			return nil
		}
		whole := region{0, 0, nr - 1, nc - 1}
		out := []stylePart{{style.child("wholeTbl"), whole}}
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
				out = append(out, stylePart{style.child(name), region{0, c, nr - 1, c}})
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
				out = append(out, stylePart{style.child(name), region{r, 0, r, nc - 1}})
			}
		}
		if flags["lastCol"] && c == nc-1 {
			out = append(out, stylePart{style.child("lastCol"), region{0, c, nr - 1, c}})
		}
		if flags["firstCol"] && c == 0 {
			out = append(out, stylePart{style.child("firstCol"), region{0, 0, nr - 1, 0}})
		}
		if flags["lastRow"] && r == nr-1 {
			out = append(out, stylePart{style.child("lastRow"), region{r, 0, r, nc - 1}})
		}
		if flags["firstRow"] && r == 0 {
			out = append(out, stylePart{style.child("firstRow"), region{0, 0, 0, nc - 1}})
		}
		cell := region{r, c, r, c}
		if flags["lastRow"] && r == nr-1 && flags["firstCol"] && c == 0 {
			out = append(out, stylePart{style.child("swCell"), cell})
		}
		if flags["lastRow"] && r == nr-1 && flags["lastCol"] && c == nc-1 {
			out = append(out, stylePart{style.child("seCell"), cell})
		}
		if flags["firstRow"] && r == 0 && flags["firstCol"] && c == 0 {
			out = append(out, stylePart{style.child("nwCell"), cell})
		}
		if flags["firstRow"] && r == 0 && flags["lastCol"] && c == nc-1 {
			out = append(out, stylePart{style.child("neCell"), cell})
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
			cl.block = s.layoutText(tf, 0, 0, w, rowH[r], identity)
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
	m := translate(xf.X, xf.Y)
	cv.obj.Save()
	cv.transform(m)
	// table background
	if bg := fillElem(tblPr); bg != nil {
		f := s.resolveFill(bg, sh.part, s.cc, nil)
		if f.kind != fillNone && cv.setFill(f, colX[nc], rowY[nr]) {
			cv.obj.FillRect(0, 0, f32(colX[nc]), f32(rowY[nr]))
			cv.drawn = true
		}
	} else if style != nil {
		if bg := style.child("tblBg"); bg != nil {
			f := s.resolveFill(fillElem(bg), "", s.cc, nil)
			if ref := bg.child("fillRef"); ref != nil {
				f = s.styleFill(ref, s.cc)
			}
			if f.kind != fillNone && cv.setFill(f, colX[nc], rowY[nr]) {
				cv.obj.FillRect(0, 0, f32(colX[nc]), f32(rowY[nr]))
				cv.drawn = true
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
				s.drawBlip(cv, f.blip, f.part, x0, y0, x1-x0, y1-y0, false)
			} else if f.kind != fillNone && cv.setFill(f, x1-x0, y1-y0) {
				cv.obj.Save()
				cv.obj.Translate(f32(x0), f32(y0))
				cv.obj.FillRect(0, 0, f32(x1-x0), f32(y1-y0))
				cv.obj.Restore()
				cv.drawn = true
			}
		}
	}
	cv.obj.Restore()
	// text
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
				b = s.layoutText(tf, 0, 0, x1-x0, y1-y0, identity)
				if b == nil {
					continue
				}
			} else {
				b.h = y1 - y0
			}
			b.fm = m.mul(translate(x0, y0)).mul(b.fm)
			s.drawTextBlock(cv, b)
		}
	}
	// borders: each edge once, the cell's own line winning over the style
	cv.obj.Save()
	cv.transform(m)
	edge := func(r, c int, side string) *line {
		cl := cells[r][c]
		if cl == nil {
			return nil
		}
		own := map[string]string{"left": "lnL", "right": "lnR", "top": "lnT", "bottom": "lnB"}[side]
		if ln := cl.tc.path("tcPr", own); ln != nil {
			return s.resolveLine([]*node{ln}, []string{sh.part}, nil, s.cc)
		}
		var found *node
		for _, p := range parts(r, c) {
			bdr := p.n.path("tcStyle", "tcBdr")
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
			if b := bdr.child(name); b != nil {
				found = b
			}
		}
		if found == nil {
			return nil
		}
		if ln := found.child("ln"); ln != nil {
			return s.resolveLine([]*node{ln}, []string{""}, nil, s.cc)
		}
		if ref := found.child("lnRef"); ref != nil {
			return s.resolveLine(nil, nil, ref, s.cc)
		}
		return nil
	}
	stroke := func(l *line, x0, y0, x1, y1 float64) {
		if l == nil {
			return
		}
		cv.obj.Save()
		cv.setLine(l, math.Abs(x1-x0), math.Abs(y1-y0))
		p := newPath().moveTo(x0, y0).lineTo(x1, y1)
		cv.obj.StrokePath(cv.obj.AddPath(p.p))
		cv.obj.Restore()
		cv.drawn = true
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
						if l2 := edge(above.r, above.c, "bottom"); l == nil || (l2 != nil && above.tc.path("tcPr", "lnB") != nil && o.tc.path("tcPr", "lnT") == nil) {
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
						if l2 := edge(left.r, left.c, "right"); l == nil || (l2 != nil && left.tc.path("tcPr", "lnR") != nil && o.tc.path("tcPr", "lnL") == nil) {
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
	cv.obj.Restore()
}

// cellText builds the text frame of a cell: its body, margins as insets,
// and the text properties of the table style parts.
func (s *slideCtx) cellText(tc *node, parts []stylePart, part string) *textFrame {
	tb := tc.child("txBody")
	if tb == nil {
		return nil
	}
	tcPr := tc.child("tcPr")
	bp := &node{Name: "bodyPr"}
	for _, kv := range [][2]string{{"marL", "lIns"}, {"marR", "rIns"}, {"marT", "tIns"}, {"marB", "bIns"}} {
		if v, ok := tcPr.attr(kv[0]); ok {
			bp.Attr = append(bp.Attr, mkNode("", kv[1], v).Attr...)
		}
	}
	if v, ok := tcPr.attr("anchor"); ok {
		bp.Attr = append(bp.Attr, mkNode("", "anchor", v).Attr...)
	}
	if v, ok := tcPr.attr("vert"); ok {
		bp.Attr = append(bp.Attr, mkNode("", "vert", v).Attr...)
	}
	tf := &textFrame{body: tb, part: part, cc: s.cc, own: tb.child("lstStyle"),
		bodyPrs: chain{bp, tb.child("bodyPr")}, lists: []*node{s.c.defTextStyle}}
	for i := len(parts) - 1; i >= 0; i-- {
		if ts := parts[i].n.child("tcTxStyle"); ts != nil {
			tf.extra = append(tf.extra, tcTxProps(ts))
		}
	}
	return tf
}

// tcTxProps turns a table style's a:tcTxStyle into run properties.
func tcTxProps(ts *node) *node {
	n := &node{Name: "defRPr"}
	for _, a := range []string{"b", "i"} {
		switch ts.attrStr(a, "def") {
		case "on":
			n.Attr = append(n.Attr, mkNode("", a, "1").Attr...)
		case "off":
			n.Attr = append(n.Attr, mkNode("", a, "0").Attr...)
		}
	}
	if fr := ts.child("fontRef"); fr != nil {
		ref := fontRefProps(fr)
		for _, k := range ref.Kids {
			if k.Name != "solidFill" {
				n.Kids = append(n.Kids, k)
			}
		}
	}
	if f := ts.child("font"); f != nil {
		n.Kids = append(n.Kids, f.Kids...)
	}
	for _, k := range ts.Kids {
		if isColor(k) {
			n.Kids = append(n.Kids, &node{Name: "solidFill", Kids: []*node{k}})
		}
	}
	return n
}

func (s *slideCtx) cellFill(tc *node, parts []stylePart, part string) fill {
	if e := fillElem(tc.child("tcPr")); e != nil {
		return s.resolveFill(e, part, s.cc, nil)
	}
	f := fill{}
	for _, p := range parts {
		ts := p.n.child("tcStyle")
		if ts == nil {
			continue
		}
		if e := fillElem(ts.child("fill")); e != nil {
			f = s.resolveFill(e, "", s.cc, nil)
		} else if ref := ts.child("fillRef"); ref != nil {
			f = s.styleFill(ref, s.cc)
		}
	}
	return f
}
