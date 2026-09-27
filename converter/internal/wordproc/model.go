package wordproc

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// The document is read into blocks (paragraphs and tables) whose text is
// already measured: items are characters, tabs, breaks and inline objects
// with their fonts and advances. Laying the blocks out (flow.go) only
// breaks lines and places them.

type block interface{ isBlock() }

func (*para) isBlock()  {}
func (*table) isBlock() {}
func (*rule) isBlock()  {}

// rule is a horizontal line across the column (HTML's hr), with CSS
// margins.
type rule struct {
	before, after float64
	ind           float64 // left indent
	width         float64 // thickness
	color         bdf.Color
}

// para is a paragraph.
type para struct {
	part      string
	pp        *pprops
	style     string
	mark      *runStyle // the paragraph mark's style (empty lines, list numbers)
	markFace  *fontset.Choice
	items     []item
	label     *label
	anchors   []*floatObj
	bookmarks []string
	heading   int      // 1-6 for headings
	ctx       []*frame // structure it is in (tables, lists, text boxes)
	sect      *section // the section this paragraph ends, if any
	hasText   bool
	breakOnly bool // holds nothing but a page or column break
}

// label is the number or bullet of a list paragraph.
type label struct {
	items  []item
	w      float64
	jc     string
	suffix string
}

// item kinds
const (
	kChar = iota
	kTab
	kBreak  // line break (w:br, w:cr)
	kPage   // page break
	kColumn // column break
	kObject // inline drawing
	kField  // PAGE, NUMPAGES … drawn when the page is known
)

// item is one character, tab, break or inline object with its style and
// face, measured.
type item struct {
	kind  uint8
	r     rune
	st    *runStyle
	fc    *fontset.Choice
	size  float64 // font size drawn (smaller for super/subscripts)
	shift float64 // baseline raise
	w     float64 // advance in points, letter spacing included
	gap   float64 // space added after it (between East Asian and Latin text)
	pad   float64 // space added after it to fill its cells of the character grid (set by layout)
	brk   bool    // a line may break after it
	x     float64 // position in the line (set by layout)
	ls    float64 // extra letter spacing from justification (-1: stretched space, not drawn)
	obj   *inlineObj
	fld   *field
	note  *note // the footnote or endnote a reference stands for
	clear bool  // a break that clears floating objects
	tcy   int   // characters set horizontally in one cell of vertical text share a group number
}

// runStyle is what the characters of a run are drawn with.
type runStyle struct {
	rp           *rprops
	fonts        [4]string // ascii, hAnsi, East Asian, complex script families
	size         float64
	bold, italic bool
	color        bdf.Color
	link         string // URL, or "#bm:" + a bookmark in the document
	lang, eaLang string
	key          string
}

// field is text that depends on the page it is drawn on.
type field struct {
	typ    string // PAGE, NUMPAGES, SECTIONPAGES, PAGEREF
	arg    string // PAGEREF: the bookmark
	format string // \* switch: roman, ROMAN, alphabetic …
	cached string // the result the document was saved with
}

// frame is a structure (spec §7.8) that text is in: a table, a cell, a
// list, a list item, a text box. Its identity is the pointer.
type frame struct {
	kind    byte
	payload string
}

const (
	fTable = iota
	fCell
	fList
	fItem
	fBox
	fFigure
)

// section is a part of the document with its own page setup.
type section struct {
	pgW, pgH                      float64
	top, bottom, left, right      float64
	header, footer                float64
	topExact, bottomExact         bool
	cols                          []column
	colSep                        bool
	titlePg                       bool
	hdr, ftr                      map[string]string // default, first, even → part
	typ                           string            // nextPage, continuous, evenPage, oddPage, nextColumn
	pgStart                       int               // -1: continue numbering
	pgFmt                         string
	gridType                      string
	linePitch                     float64
	charSpace                     float64 // character grid: pitch added to the default font size, points
	vertical                      bool    // East Asian vertical text: lines top to bottom, right to left
	blocks                        []block
	borders                       [4]*border
	bordersOffsetFromText         bool
	bordersFirst, bordersNotFirst bool // page borders on the first page only, on all but the first
}

// column is a text column, relative to the left margin.
type column struct{ x, w float64 }

func (s *section) textWidth() float64 { return math.Max(s.pgW-s.left-s.right, 36) }

// lineLength is the length of the lines between the margins: the text
// width, or for vertical text the text height.
func (s *section) lineLength() float64 {
	if s.vertical {
		return math.Max(s.pgH-s.top-s.bottom, 36)
	}
	return s.textWidth()
}

// column returns the start and length of a column of the section's lines
// in the flow's coordinates (see page.rotation for vertical text) on a page
// whose lines are lineLen long.
func (s *section) column(i int, lineLen float64) (x0, w float64) {
	c := s.cols[min(i, len(s.cols)-1)]
	switch {
	case s.vertical && len(s.cols) == 1:
		return 0, lineLen
	case s.vertical:
		return c.x, c.w
	}
	return s.left + c.x, c.w
}

// grid returns the line pitch that lines snap to (0 for none).
func (s *section) grid() float64 {
	if s.gridType == "lines" || s.gridType == "linesAndChars" || s.gridType == "snapToChars" {
		return s.linePitch
	}
	return 0
}

func (c *converter) parseSection(n *ooxml.Node, prev *section) *section {
	s := &section{pgW: 612, pgH: 792, top: 72, bottom: 72, left: 72, right: 72, header: 36, footer: 36,
		typ: "nextPage", pgStart: -1, pgFmt: "decimal", hdr: map[string]string{}, ftr: map[string]string{}}
	if prev != nil {
		s.pgW, s.pgH = prev.pgW, prev.pgH
	}
	if sz := n.Child("pgSz"); sz != nil {
		s.pgW, s.pgH = twips(sz, "w", s.pgW), twips(sz, "h", s.pgH)
	}
	if m := n.Child("pgMar"); m != nil {
		s.top, s.bottom = twips(m, "top", 72), twips(m, "bottom", 72)
		s.left, s.right = twips(m, "left", 72), twips(m, "right", 72)
		s.header, s.footer = twips(m, "header", 36), twips(m, "footer", 36)
		s.left += twips(m, "gutter", 0)
	}
	if s.top < 0 {
		s.top, s.topExact = -s.top, true
	}
	if s.bottom < 0 {
		s.bottom, s.bottomExact = -s.bottom, true
	}
	s.pgW = math.Max(s.pgW, 72)
	s.pgH = math.Max(s.pgH, 72)
	if t := n.Child("type"); t != nil {
		s.typ = t.AttrStr("val", "nextPage")
	}
	s.titlePg = onOff(n.Child("titlePg"))
	if p := n.Child("pgNumType"); p != nil {
		if _, ok := p.Attr("start"); ok {
			s.pgStart = int(p.AttrInt("start", 1))
		}
		s.pgFmt = p.AttrStr("fmt", "decimal")
	}
	if g := n.Child("docGrid"); g != nil {
		s.gridType = g.AttrStr("type", "default")
		s.linePitch = twips(g, "linePitch", 0)
		s.charSpace = g.AttrFloat("charSpace", 0) / 4096
	}
	if td := n.Child("textDirection"); td != nil {
		switch val(td) {
		case "tbRl", "tbRlV", "tbLrV":
			s.vertical = true
		}
	}
	// columns divide the lines: for vertical text they are bands from top
	// to bottom
	tw := s.lineLength()
	s.cols = []column{{0, tw}}
	if cn := n.Child("cols"); cn != nil {
		num := int(cn.AttrInt("num", 1))
		space := twips(cn, "space", 36)
		s.colSep = cn.AttrBool("sep", false)
		if defs := cn.Children("col"); len(defs) > 0 && !cn.AttrBool("equalWidth", true) {
			s.cols = nil
			x := 0.0
			for _, d := range defs {
				w := twips(d, "w", tw/float64(len(defs)))
				s.cols = append(s.cols, column{x, w})
				x += w + twips(d, "space", 0)
			}
		} else if num > 1 && num <= 45 {
			w := (tw - float64(num-1)*space) / float64(num)
			if w > 18 {
				s.cols = nil
				for i := range num {
					s.cols = append(s.cols, column{float64(i) * (w + space), w})
				}
			}
		}
	}
	for _, ref := range n.Children("headerReference") {
		if r, ok := c.pkg.Target(c.main, ref.RelID("id")); ok {
			s.hdr[ref.AttrStr("type", "default")] = r.Target
		}
	}
	for _, ref := range n.Children("footerReference") {
		if r, ok := c.pkg.Target(c.main, ref.RelID("id")); ok {
			s.ftr[ref.AttrStr("type", "default")] = r.Target
		}
	}
	if pb := n.Child("pgBorders"); pb != nil {
		for i, name := range []string{"top", "left", "bottom", "right"} {
			s.borders[i] = c.parseBorder(pb.Child(name))
		}
		s.bordersOffsetFromText = pb.AttrStr("offsetFrom", "text") == "text"
		switch pb.AttrStr("display", "allPages") {
		case "firstPage":
			s.bordersFirst = true
		case "notFirstPage":
			s.bordersNotFirst = true
		}
	}
	if n.Child("lnNumType") != nil {
		c.warnOnce("lnnum", "line numbers are not drawn")
	}
	return s
}

// walker reads the blocks of one story (the body, a header, a note, a text
// box, a cell) of a part.
type walker struct {
	c      *converter
	part   string
	tl     *tableLayers // table style formatting of the cell being read
	base   []*frame     // structure the story is in
	lists  []openList
	fields []*fieldState
	dark   bool // the story is on a dark background (auto text is white)
	noteNo string
}

type openList struct {
	lvl   int
	numID string
	list  *frame
	item  *frame
}

type fieldState struct {
	instr   strings.Builder
	result  bool // past the separator
	typ     string
	dyn     *field
	emitted bool
	link    string
}

// hidden reports whether text is in a field instruction.
func (w *walker) hidden() bool {
	for _, f := range w.fields {
		if !f.result {
			return true
		}
	}
	return false
}

// blocks reads the block-level content of a container (w:body,
// w:tc, w:txbxContent, w:hdr …).
func (w *walker) blocks(n *ooxml.Node) []block {
	var out []block
	for _, k := range n.Elements() {
		switch k.Name {
		case "p":
			out = append(out, w.paragraph(k))
		case "tbl":
			w.closeLists(-1)
			if t := w.table(k); t != nil {
				out = append(out, t)
			}
		case "sdt":
			out = append(out, w.blocks(k.Child("sdtContent"))...)
		case "customXml", "ins", "moveTo":
			out = append(out, w.blocks(k)...)
		case "altChunk":
			w.c.warnOnce("altchunk", "embedded documents (altChunk) are not drawn")
		}
	}
	return out
}

func (w *walker) closeLists(lvl int) {
	for len(w.lists) > 0 && w.lists[len(w.lists)-1].lvl > lvl {
		w.lists = w.lists[:len(w.lists)-1]
	}
}

// listCtx returns the structure of a paragraph: the story's, then the
// lists it is in. A numbered paragraph starts an item of the list for its
// level (lists of deeper levels nest in the items above them); other
// paragraphs with text close the lists, empty ones continue the item.
func (w *walker) listCtx(p *para, listed bool) []*frame {
	ctx := append([]*frame(nil), w.base...)
	switch {
	case listed:
		lvl := p.pp.ilvl
		w.closeLists(lvl)
		if n := len(w.lists); n > 0 && w.lists[n-1].lvl == lvl && w.lists[n-1].numID != p.pp.numID {
			// another list
			w.lists = w.lists[:n-1]
		}
		if n := len(w.lists); n == 0 || w.lists[n-1].lvl < lvl {
			w.lists = append(w.lists, openList{lvl: lvl, numID: p.pp.numID, list: &frame{kind: fList}})
		}
		w.lists[len(w.lists)-1].item = &frame{kind: fItem}
	case p.hasText:
		w.lists = nil
	}
	for _, l := range w.lists {
		ctx = append(ctx, l.list)
		if l.item != nil {
			ctx = append(ctx, l.item)
		}
	}
	return ctx
}

// paragraph reads a w:p.
func (w *walker) paragraph(n *ooxml.Node) *para {
	c := w.c
	pPr := n.Child("pPr")
	styleID := c.st.paraStyle(val(pPr.Child("pStyle")))
	p := &para{part: w.part, style: styleID}
	p.pp = c.paraProps(styleID, w.tl, pPr)
	p.mark = c.runStyle(c.runProps(styleID, "", w.tl, pPr.Child("rPr")), "", w.dark)
	p.markFace = c.face(p.mark, 'x')
	if p.pp.outline >= 0 {
		p.heading = min(p.pp.outline+1, 6)
	}
	if sp := pPr.Child("sectPr"); sp != nil {
		p.sect = c.parseSection(sp, c.lastSection)
		c.lastSection = p.sect
	}
	w.inlines(p, n, "")
	w.finishParagraph(p)
	listed := false
	if p.pp.numID != "" && p.pp.numID != "0" {
		text, lvl := c.num.next(p.pp.numID, p.pp.ilvl)
		if lvl != nil && (text != "" || lvl.format != "none") {
			p.label = c.makeLabel(p, text, lvl)
			listed = p.heading == 0
		}
	}
	p.ctx = w.listCtx(p, listed)
	return p
}

// finishParagraph sets the break opportunities and the spacing between
// East Asian and Latin text.
func (w *walker) finishParagraph(p *para) {
	items := p.items
	for i := range items {
		it := &items[i]
		if it.kind == kField {
			// measured with the result the document was saved with
			text := strings.TrimSpace(it.fld.cached)
			if text == "" {
				text = "0"
			}
			it.w = 0
			for _, r := range text {
				it.w += w.c.fonts.Advance(w.c.face(it.st, r), r) * it.size
			}
		}
		if it.kind == kChar && !unicode.IsSpace(it.r) || it.kind == kField {
			p.hasText = true
		}
	}
	if !p.hasText {
		only := len(items) > 0
		for _, it := range items {
			if it.kind != kPage && it.kind != kColumn {
				only = false
			}
		}
		p.breakOnly = only
	}
	for i := range items {
		it := &items[i]
		var next *item
		if i+1 < len(items) {
			next = &items[i+1]
		}
		switch {
		case it.kind == kObject:
			it.brk = true
		case it.kind != kChar:
			it.brk = it.kind == kTab
		case next == nil || next.kind == kObject || next.kind == kTab || next.kind == kBreak:
			it.brk = true
		case next.kind != kChar && next.kind != kField:
			it.brk = true
		case next.kind == kField, it.tcy != 0 && next.tcy == it.tcy:
			it.brk = false
		default:
			if p.pp.kinsoku {
				it.brk = allowed(it.r, next.r)
			} else {
				it.brk = allowedNoKinsoku(it.r, next.r)
			}
		}
		// spacing between East Asian and Latin text or numbers
		if next != nil && it.kind == kChar && next.kind == kChar {
			a, b := it.r, next.r
			if (p.pp.autoSpaceDE && (eaLetter(a) && latinLetter(b) || latinLetter(a) && eaLetter(b))) ||
				(p.pp.autoSpaceDN && (eaLetter(a) && unicode.IsDigit(b) && b < 0x80 || a < 0x80 && unicode.IsDigit(a) && eaLetter(b))) {
				it.gap = 0.25 * math.Min(it.size, next.size)
			}
		}
	}
}

func eaLetter(r rune) bool {
	return unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana)
}

func latinLetter(r rune) bool {
	return r < 0x250 && unicode.IsLetter(r)
}

// inlines reads the run-level content of a paragraph (or of a hyperlink,
// field, content control … inside it); link is the hyperlink it is in.
func (w *walker) inlines(p *para, n *ooxml.Node, link string) {
	c := w.c
	for _, k := range n.Elements() {
		switch k.Name {
		case "r":
			w.run(p, k, link)
		case "hyperlink":
			l := link
			if a := k.AttrStr("anchor", ""); a != "" {
				l = "#bm:" + a
			} else if r, ok := c.pkg.Target(w.part, k.RelID("id")); ok && r.External {
				l = safeURL(r.Target)
				if a := k.AttrStr("anchor", ""); a != "" && l != "" {
					l += "#" + a
				}
			}
			w.inlines(p, k, l)
		case "fldSimple":
			f := w.beginField(k.AttrStr("instr", ""))
			f.result = true
			w.inlines(p, k, w.linkOf(link))
			w.endField(p, k.Child("r"))
		case "smartTag", "customXml", "ins", "moveTo", "dir", "bdo":
			w.inlines(p, k, link)
		case "sdt":
			w.inlines(p, k.Child("sdtContent"), link)
		case "bookmarkStart":
			if name := k.AttrStr("name", ""); name != "" {
				p.bookmarks = append(p.bookmarks, name)
			}
		case "oMath", "oMathPara":
			w.math(p, k, link)
		}
	}
}

// linkOf returns the link of the innermost HYPERLINK field in effect.
func (w *walker) linkOf(link string) string {
	for i := len(w.fields) - 1; i >= 0; i-- {
		if f := w.fields[i]; f.result && f.link != "" {
			return f.link
		}
	}
	return link
}

// math draws the text of an Office Math object linearly.
func (w *walker) math(p *para, n *ooxml.Node, link string) {
	w.c.warnOnce("math", "equations are drawn as plain text")
	var walk func(n *ooxml.Node)
	walk = func(n *ooxml.Node) {
		for _, k := range n.Elements() {
			if k.Name == "r" {
				w.run(p, k, link)
				continue
			}
			walk(k)
		}
	}
	walk(n)
}

func safeURL(s string) string {
	u, err := url.Parse(s)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "mailto") {
		return ""
	}
	return s
}

func (w *walker) beginField(instr string) *fieldState {
	f := &fieldState{}
	f.instr.WriteString(instr)
	w.fields = append(w.fields, f)
	return f
}

// separate ends a field's instruction: what the field is decides how its
// result is drawn.
func (w *walker) separate(f *fieldState) {
	f.result = true
	fs := strings.Fields(f.instr.String())
	if len(fs) == 0 {
		return
	}
	f.typ = strings.ToUpper(fs[0])
	format := ""
	for i, s := range fs {
		if s == `\*` && i+1 < len(fs) {
			format = fs[i+1]
		}
	}
	switch f.typ {
	case "PAGE", "NUMPAGES", "SECTIONPAGES":
		f.dyn = &field{typ: f.typ, format: format}
	case "PAGEREF":
		if len(fs) > 1 {
			f.dyn = &field{typ: f.typ, arg: fs[1], format: format}
		}
	case "HYPERLINK":
		args := quotedArgs(f.instr.String())
		for i := 0; i < len(args); i++ {
			switch args[i] {
			case `\l`:
				if i+1 < len(args) {
					f.link = "#bm:" + args[i+1]
					i++
				}
			case `\o`, `\t`, `\m`, `\n`:
				i++
			default:
				if f.link == "" && !strings.HasPrefix(args[i], `\`) {
					f.link = safeURL(args[i])
				}
			}
		}
	}
}

func (w *walker) endField(p *para, rn *ooxml.Node) {
	n := len(w.fields)
	if n == 0 {
		return
	}
	f := w.fields[n-1]
	if !f.result {
		w.separate(f)
	}
	w.fields = w.fields[:n-1]
	if f.dyn != nil && !f.emitted && !w.hidden() {
		// a field without a result: draw it in the style of its end
		w.fieldItem(p, f, w.c.runStyle(w.runRProps(p, rn), w.linkOf(""), w.dark))
	}
}

// quotedArgs splits a field instruction into its arguments, honoring
// double quotes.
func quotedArgs(s string) []string {
	var out []string
	var b strings.Builder
	in := false
	for _, r := range s {
		switch {
		case r == '"':
			in = !in
		case unicode.IsSpace(r) && !in:
			if b.Len() > 0 {
				out = append(out, b.String())
				b.Reset()
			}
		default:
			b.WriteRune(r)
		}
	}
	if b.Len() > 0 {
		out = append(out, b.String())
	}
	if len(out) > 0 {
		out = out[1:]
	}
	return out
}

func (w *walker) runRProps(p *para, r *ooxml.Node) *rprops {
	rPr := r.Child("rPr")
	return w.c.runProps(p.style, val(rPr.Child("rStyle")), w.tl, rPr)
}

// run reads a w:r.
func (w *walker) run(p *para, r *ooxml.Node, link string) {
	c := w.c
	rp := w.runRProps(p, r)
	if rp.tcy {
		// horizontal in vertical text: the run's characters share a cell
		start := len(p.items)
		c.tcyGroups++
		defer func() {
			for i := start; i < len(p.items); i++ {
				if p.items[i].kind == kChar {
					p.items[i].tcy = c.tcyGroups
				}
			}
		}()
	}
	var st *runStyle
	style := func() *runStyle {
		if st == nil {
			st = c.runStyle(rp, w.linkOf(link), w.dark)
		}
		return st
	}
	for _, k := range r.Elements() {
		switch k.Name {
		case "fldChar":
			switch k.AttrStr("fldCharType", "") {
			case "begin":
				w.beginField("")
			case "separate":
				if n := len(w.fields); n > 0 {
					w.separate(w.fields[n-1])
				}
			case "end":
				w.endField(p, r)
			}
			st = nil // links may have changed
			continue
		case "instrText":
			if n := len(w.fields); n > 0 && !w.fields[n-1].result {
				w.fields[n-1].instr.WriteString(k.Content())
			}
			continue
		}
		if w.hidden() || rp.has(tVanish) {
			continue
		}
		if f := w.dynField(); f != nil {
			// the result of a page field: drawn when the page is known
			if !f.emitted {
				if t := k.Content(); k.Name == "t" {
					f.dyn.cached += t
				}
				w.fieldItem(p, f, style())
				f.emitted = true
			} else if k.Name == "t" {
				f.dyn.cached += k.Content()
			}
			continue
		}
		switch k.Name {
		case "t":
			c.addText(p, style(), k.Content())
		case "tab":
			p.items = append(p.items, item{kind: kTab, r: '\t', st: style(), fc: c.face(style(), ' '), size: style().size})
		case "ptab":
			p.items = append(p.items, item{kind: kTab, r: '\t', st: style(), fc: c.face(style(), ' '), size: style().size})
		case "br":
			switch k.AttrStr("type", "textWrapping") {
			case "page":
				p.items = append(p.items, item{kind: kPage, st: style()})
			case "column":
				p.items = append(p.items, item{kind: kColumn, st: style()})
			default:
				p.items = append(p.items, item{kind: kBreak, r: '\n', st: style(), clear: k.AttrStr("clear", "none") != "none"})
			}
		case "cr":
			p.items = append(p.items, item{kind: kBreak, r: '\n', st: style()})
		case "noBreakHyphen":
			c.addText(p, style(), "\u2011")
		case "sym":
			ch, err := strconv.ParseUint(k.AttrStr("char", ""), 16, 32)
			if err != nil {
				continue
			}
			font := k.AttrStr("font", "")
			text := fontset.MapSymbol(font, string(rune(ch)))
			sst := style()
			if !fontset.IsSymbol(font) && font != "" {
				cp := *rp
				cp.fonts = [4]string{font, font, font, font}
				cp.themes = [4]string{}
				sst = c.runStyle(&cp, w.linkOf(link), w.dark)
			}
			c.addText(p, sst, text)
		case "footnoteReference", "endnoteReference":
			kind := "f"
			if k.Name == "endnoteReference" {
				kind = "e"
			}
			num := c.noteReference(kind, k.AttrStr("id", ""))
			if num != "" {
				n := len(p.items)
				c.addText(p, style(), num)
				if n < len(p.items) {
					p.items[n].note = c.notes[kind+k.AttrStr("id", "")]
				}
			}
		case "footnoteRef", "endnoteRef":
			c.addText(p, style(), w.noteNo)
		case "separator", "continuationSeparator":
			// drawn by the page layout
		case "drawing":
			w.drawing(p, k, style())
		case "pict", "object":
			w.vml(p, k, style())
		case "ruby":
			// the base text; the reading above it is left out
			for _, b := range k.Child("rubyBase").Children("r") {
				w.run(p, b, link)
			}
		case "softHyphen", "lastRenderedPageBreak", "annotationRef", "commentReference", "delText", "delInstrText":
		}
	}
}

func (w *walker) dynField() *fieldState {
	for i := len(w.fields) - 1; i >= 0; i-- {
		if f := w.fields[i]; f.result && f.dyn != nil {
			return f
		}
	}
	return nil
}

func (w *walker) fieldItem(p *para, f *fieldState, st *runStyle) {
	f.emitted = true
	p.items = append(p.items, item{kind: kField, fld: f.dyn, st: st, fc: w.c.face(st, '0'), size: st.size})
}

// runStyle derives the drawing style of a run.
func (c *converter) runStyle(rp *rprops, link string, darkBG bool) *runStyle {
	st := &runStyle{rp: rp, size: rp.sz, bold: rp.has(tB), italic: rp.has(tI), color: rp.color, link: link}
	for i := range 4 {
		f := rp.fonts[i]
		if t := rp.themes[i]; t != "" {
			if th := themeFontRef(t); th != "" {
				if tf := c.r.ThemeFont(c.main, th, c.eaScript); tf != "" {
					f = tf
				}
			}
		}
		st.fonts[i] = f
	}
	if st.fonts[2] == "" {
		st.fonts[2] = st.fonts[0]
	}
	if rp.autoColor {
		st.color = bdf.RGB(0, 0, 0)
		if darkBG || (rp.hasShd && dark(rp.shd)) {
			st.color = bdf.RGB(255, 255, 255)
		}
	}
	st.lang, st.eaLang = rp.lang[0], rp.lang[1]
	st.key = fmt.Sprintf("%v|%v|%g|%v|%v|%x|%s|%s|%s|%s|%g|%g|%g|%s|%x|%v|%x|%v|%x|%s", st.fonts, rp.toggles, st.size, st.bold, st.italic,
		st.color, link, rp.u, rp.vertAlign, rp.highlight, rp.position, rp.spacing, rp.scale, rp.em, rp.shd, rp.hasShd, rp.uColor, rp.hasUColor, rp.hint, st.lang+"/"+st.eaLang)
	return st
}

func themeFontRef(t string) string {
	switch t {
	case "majorAscii", "majorHAnsi":
		return "+mj-lt"
	case "majorEastAsia":
		return "+mj-ea"
	case "majorBidi":
		return "+mj-cs"
	case "minorAscii", "minorHAnsi":
		return "+mn-lt"
	case "minorEastAsia":
		return "+mn-ea"
	case "minorBidi":
		return "+mn-cs"
	}
	return ""
}

// face picks the face that draws r in a run: the East Asian font for East
// Asian characters (and, with the eastAsia hint, for the symbols both
// kinds of text use), the ASCII font for ASCII and the high ANSI font for
// other characters.
func (c *converter) face(st *runStyle, r rune) *fontset.Choice {
	bold, italic := st.bold, st.italic
	if fontdb.IsCJK(r) {
		return c.fonts.FaceFor(st.fonts[1], st.fonts[2], bold, italic, r)
	}
	if st.rp.hint == "eastAsia" && ambiguousEA(r) {
		return c.fonts.FaceFor(st.fonts[2], st.fonts[1], bold, italic, r)
	}
	if r < 0x80 {
		return c.fonts.FaceFor(st.fonts[0], st.fonts[2], bold, italic, r)
	}
	return c.fonts.FaceFor(st.fonts[1], st.fonts[2], bold, italic, r)
}

// ambiguousEA reports whether r is a character that East Asian and Latin
// text share (quotes, symbols, box drawing), which the eastAsia font hint
// sends to the East Asian font.
func ambiguousEA(r rune) bool {
	switch {
	case r >= 0x2000 && r <= 0x27BF, r >= 0x2E80 && r <= 0x2EFF, r >= 0xFF00 && r <= 0xFFEF:
		return true
	case r == 0xA7 || r == 0xA8 || r == 0xB0 || r == 0xB1 || r == 0xB4 || r == 0xB6 || r == 0xD7 || r == 0xF7:
		return true
	}
	return false
}

// addText appends the characters of text in a run's style, with the
// run's capitals and small capitals.
func (c *converter) addText(p *para, st *runStyle, text string) {
	rp := st.rp
	switch {
	case rp.has(tCaps):
		text = strings.ToUpper(text)
	case rp.has(tSmallCaps):
		small := *st
		small.size = st.size * 0.8
		small.key += "|sc"
		for _, r := range text {
			if unicode.IsLower(r) {
				c.addChar(p, &small, unicode.ToUpper(r))
			} else {
				c.addChar(p, st, r)
			}
		}
		return
	}
	for _, r := range text {
		c.addChar(p, st, r)
	}
}

func (c *converter) addChar(p *para, st *runStyle, r rune) {
	switch r {
	case '\r', '\n':
		p.items = append(p.items, item{kind: kBreak, r: '\n', st: st})
		return
	case '\t':
		p.items = append(p.items, item{kind: kTab, r: '\t', st: st, fc: c.face(st, ' '), size: st.size})
		return
	case '\u00ad':
		// a soft hyphen shows only where Word breaks a word at it
		return
	}
	size, shift := st.size, st.rp.position
	switch st.rp.vertAlign {
	case "superscript":
		shift += 0.33 * size
		size *= 0.65
	case "subscript":
		shift -= 0.08 * size
		size *= 0.65
	}
	switch {
	case fontdb.IsCJK(r):
		c.cjk++
	case unicode.IsLetter(r):
		c.latin++
	}
	fc := c.face(st, r)
	w := c.fonts.Advance(fc, r)*size*st.rp.scale + st.rp.spacing
	p.items = append(p.items, item{kind: kChar, r: r, st: st, fc: fc, size: size, shift: shift, w: w})
}

// makeLabel builds the number or bullet of a list paragraph: its text in
// the paragraph mark's style with the level's run properties.
func (c *converter) makeLabel(p *para, text string, lvl *numLevel) *label {
	// the paragraph mark's own formatting, then the level's
	rp := p.mark.rp.clone()
	c.applyRPr(rp, lvl.rPr)
	rp.u, rp.hasShd, rp.highlight = "none", false, ""
	st := c.runStyle(rp, "", false)
	font := ""
	if lvl.rPr != nil {
		font = lvl.rPr.Child("rFonts").AttrStr("ascii", "")
	}
	if font != "" && fontset.IsSymbol(font) {
		text = fontset.MapSymbol(font, text)
		cp := *rp
		cp.fonts = p.mark.rp.fonts
		cp.themes = p.mark.rp.themes
		st = c.runStyle(&cp, "", false)
	}
	lb := &label{jc: lvl.jc, suffix: lvl.suffix}
	tmp := &para{}
	c.addText(tmp, st, text)
	lb.items = tmp.items
	for _, it := range lb.items {
		lb.w += it.w
	}
	return lb
}

// allowed reports whether a line may break between a and b.
func allowed(a, b rune) bool { return lbAllowed(a, b) }

// allowedNoKinsoku is allowed without the East Asian line breaking rules.
func allowedNoKinsoku(a, b rune) bool {
	if lbSpace(a) {
		return !lbSpace(b)
	}
	if lbSpace(b) {
		return false
	}
	return fontdb.IsCJK(a) || fontdb.IsCJK(b)
}
