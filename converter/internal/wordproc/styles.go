package wordproc

import (
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// style is a style definition (w:style).
type style struct {
	id, typ, name string
	basedOn       string
	pPr, rPr      *ooxml.Node
	tblPr, trPr   *ooxml.Node
	tcPr          *ooxml.Node
	cond          map[string]*condStyle // table style conditional formatting by type
}

// condStyle is the formatting of one part of a table (first row, banded
// columns …).
type condStyle struct {
	pPr, rPr, tblPr, trPr, tcPr *ooxml.Node
}

// styles is the document's style sheet.
type styles struct {
	byID             map[string]*style
	defPara, defChar string
	defTable         string
	docPPr, docRPr   *ooxml.Node
}

func (c *converter) loadStyles() {
	c.st = &styles{byID: map[string]*style{}}
	r, ok := c.pkg.RelOfType(c.main, "/styles")
	if !ok {
		return
	}
	root, err := c.pkg.XML(r.Target)
	if err != nil {
		c.warnf("styles: %v", err)
		return
	}
	dd := root.Child("docDefaults")
	c.st.docPPr = dd.Path("pPrDefault", "pPr")
	c.st.docRPr = dd.Path("rPrDefault", "rPr")
	for _, s := range root.Children("style") {
		st := &style{id: s.AttrStr("styleId", ""), typ: s.AttrStr("type", "paragraph"), name: val(s.Child("name")),
			basedOn: val(s.Child("basedOn")), pPr: s.Child("pPr"), rPr: s.Child("rPr"), tblPr: s.Child("tblPr"),
			trPr: s.Child("trPr"), tcPr: s.Child("tcPr")}
		for _, t := range s.Children("tblStylePr") {
			if st.cond == nil {
				st.cond = map[string]*condStyle{}
			}
			st.cond[t.AttrStr("type", "")] = &condStyle{pPr: t.Child("pPr"), rPr: t.Child("rPr"), tblPr: t.Child("tblPr"),
				trPr: t.Child("trPr"), tcPr: t.Child("tcPr")}
		}
		if _, dup := c.st.byID[st.id]; !dup {
			c.st.byID[st.id] = st
		}
		if s.AttrBool("default", false) {
			switch st.typ {
			case "paragraph":
				c.st.defPara = st.id
			case "character":
				c.st.defChar = st.id
			case "table":
				c.st.defTable = st.id
			}
		}
	}
}

// chain returns a style and the styles it is based on, base first.
func (s *styles) chain(id string) []*style {
	var out []*style
	seen := map[string]bool{}
	for id != "" && !seen[id] && len(out) < 32 {
		seen[id] = true
		st := s.byID[id]
		if st == nil {
			break
		}
		out = append(out, st)
		id = st.basedOn
	}
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}

// paraStyle returns the paragraph style a paragraph uses (the default
// paragraph style when it names none or an unknown one).
func (s *styles) paraStyle(id string) string {
	if st := s.byID[id]; st != nil && st.typ == "paragraph" {
		return id
	}
	return s.defPara
}

// tableLayers are the formatting layers a table style gives the content of
// one cell: the whole table's, then the conditional parts that apply, in
// the order they override each other.
type tableLayers struct {
	pPr, rPr, tcPr []*ooxml.Node
}

// paraProps resolves the properties of a paragraph: the document
// defaults, the table style's, the paragraph style chain's, the numbering
// level's, and its own.
func (c *converter) paraProps(styleID string, tl *tableLayers, direct *ooxml.Node) *pprops {
	p := defaultPProps()
	c.applyPPr(p, c.st.docPPr)
	if tl != nil {
		for _, n := range tl.pPr {
			c.applyPPr(p, n)
		}
	}
	for _, st := range c.st.chain(c.st.paraStyle(styleID)) {
		c.applyPPr(p, st.pPr)
	}
	// numbering: from the paragraph's own numPr, else its style's
	numID, ilvl := p.numID, p.ilvl
	if n := direct.Child("numPr"); n != nil {
		if id := n.Child("numId"); id != nil {
			numID = val(id)
		}
		if l := n.Child("ilvl"); l != nil {
			ilvl = int(l.AttrInt("val", 0))
		}
	}
	if lvl := c.num.level(numID, ilvl); lvl != nil {
		c.applyPPr(p, lvl.pPr)
	}
	c.applyPPr(p, direct)
	p.numID, p.ilvl = numID, ilvl
	return p
}

// runProps resolves the properties of a run: the document defaults, the
// table style's, the paragraph style chain's, then the character style
// chain's (whose toggle properties flip those set before), then its own.
func (c *converter) runProps(pStyle, rStyle string, tl *tableLayers, direct ...*ooxml.Node) *rprops {
	r := defaultRProps()
	c.applyRPr(r, c.st.docRPr)
	if tl != nil {
		for _, n := range tl.rPr {
			c.applyRPr(r, n)
		}
	}
	for _, st := range c.st.chain(c.st.paraStyle(pStyle)) {
		c.applyRPr(r, st.rPr)
	}
	if rStyle == "" {
		rStyle = c.st.defChar
	}
	if cs := c.st.chain(rStyle); len(cs) > 0 && cs[len(cs)-1].typ == "character" {
		ch := r.clone()
		ch.set = 0
		for _, st := range cs {
			c.applyRPr(ch, st.rPr)
		}
		// toggles set by the character style flip the inherited value
		tog := r.toggles ^ (ch.toggles & ch.set)
		ch.toggles = tog
		r = ch
	}
	for _, d := range direct {
		c.applyRPr(r, d)
	}
	r.set = 0
	return r
}

// tableStyleLayers returns the table style formatting of the cell at row,
// col (0-based) of a table with rows×cols cells, for the table look.
func (c *converter) tableStyleLayers(styleID string, look tblLook, row, col, rows, cols, rowBand, colBand int, header bool) *tableLayers {
	chain := c.st.chain(styleID)
	if len(chain) == 0 {
		return nil
	}
	tl := &tableLayers{}
	add := func(cs *condStyle) {
		if cs == nil {
			return
		}
		if cs.pPr != nil {
			tl.pPr = append(tl.pPr, cs.pPr)
		}
		if cs.rPr != nil {
			tl.rPr = append(tl.rPr, cs.rPr)
		}
		if cs.tcPr != nil {
			tl.tcPr = append(tl.tcPr, cs.tcPr)
		}
	}
	for _, st := range chain {
		add(&condStyle{pPr: st.pPr, rPr: st.rPr, tcPr: st.tcPr})
	}
	types := condTypes(look, row, col, rows, cols, rowBand, colBand)
	for _, t := range types {
		for _, st := range chain {
			add(st.cond[t])
		}
	}
	return tl
}

// tblLook says which conditional formats of a table style apply.
type tblLook struct {
	firstRow, lastRow, firstCol, lastCol, noHBand, noVBand bool
}

func parseLook(n *ooxml.Node) tblLook {
	l := tblLook{firstRow: true, firstCol: true, noVBand: true}
	if n == nil {
		return l
	}
	if v, ok := n.Attr("val"); ok {
		var x int64
		for _, ch := range v {
			x *= 16
			switch {
			case ch >= '0' && ch <= '9':
				x += int64(ch - '0')
			case ch >= 'a' && ch <= 'f':
				x += int64(ch-'a') + 10
			case ch >= 'A' && ch <= 'F':
				x += int64(ch-'A') + 10
			}
		}
		l = tblLook{firstRow: x&0x20 != 0, lastRow: x&0x40 != 0, firstCol: x&0x80 != 0, lastCol: x&0x100 != 0,
			noHBand: x&0x200 != 0, noVBand: x&0x400 != 0}
	}
	for _, a := range []struct {
		dst  *bool
		name string
	}{{&l.firstRow, "firstRow"}, {&l.lastRow, "lastRow"}, {&l.firstCol, "firstColumn"}, {&l.lastCol, "lastColumn"},
		{&l.noHBand, "noHBand"}, {&l.noVBand, "noVBand"}} {
		if _, ok := n.Attr(a.name); ok {
			*a.dst = n.AttrBool(a.name, false)
		}
	}
	return l
}

// condTypes lists the conditional formats of a cell in the order they
// apply: banded columns, banded rows, first and last column, first and
// last row, then the corner cells.
func condTypes(look tblLook, row, col, rows, cols, rowBand, colBand int) []string {
	var out []string
	if rowBand < 1 {
		rowBand = 1
	}
	if colBand < 1 {
		colBand = 1
	}
	isFirstRow := look.firstRow && row == 0
	isLastRow := look.lastRow && row == rows-1
	isFirstCol := look.firstCol && col == 0
	isLastCol := look.lastCol && col == cols-1
	if !look.noVBand && !isFirstCol && !isLastCol {
		c := col
		if look.firstCol {
			c--
		}
		if (c/colBand)%2 == 0 {
			out = append(out, "band1Vert")
		} else {
			out = append(out, "band2Vert")
		}
	}
	if !look.noHBand && !isFirstRow && !isLastRow {
		r := row
		if look.firstRow {
			r--
		}
		if (r/rowBand)%2 == 0 {
			out = append(out, "band1Horz")
		} else {
			out = append(out, "band2Horz")
		}
	}
	if isFirstCol {
		out = append(out, "firstCol")
	}
	if isLastCol {
		out = append(out, "lastCol")
	}
	if isFirstRow {
		out = append(out, "firstRow")
	}
	if isLastRow {
		out = append(out, "lastRow")
	}
	switch {
	case isFirstRow && isFirstCol:
		out = append(out, "nwCell")
	case isFirstRow && isLastCol:
		out = append(out, "neCell")
	case isLastRow && isFirstCol:
		out = append(out, "swCell")
	case isLastRow && isLastCol:
		out = append(out, "seCell")
	}
	return out
}
