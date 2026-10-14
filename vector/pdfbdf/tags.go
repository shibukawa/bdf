package pdfbdf

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
)

// A tagged PDF says what its content is: the structure tree of the file
// (headings, paragraphs, lists, tables, figures with their alternative
// texts) in the order it is read, with each piece of the content of the
// pages marked as that of an element of the tree, or as an artifact (what
// is not content: backgrounds, the masters of slides, running heads, the
// lines of tables). Readers of the file that do not show it (screen
// readers, tools that take the text out or flow it again) go by the tree.
//
// The tree is made of the MARKs of the objects, by the rules of the
// structure of documents (docs/spec.md §7.8), which the viewer follows for
// its text layer: a leaf MARK (PARAGRAPH, BOX, HEADING, a LIST_ITEM outside
// a list, a CELL outside a table) says what the next run of text is of,
// and a structure MARK (LIST, TABLE, FIGURE, END, an item in a list, a
// cell in a table) opens or closes a container. Text is the content of
// the leaf it is in; what is drawn in a figure is the content of the
// figure; everything else is an artifact.

// structElem is an element of the structure tree.
type structElem struct {
	tag    string
	parent *structElem
	kids   []structKid
	ref    ref
	// page is the page of the first content of the element.
	page  *pageOut
	lang  string
	alt   string
	attrs string // the entries of its attribute object, of the owner Table
	// grid is where a cell of a table is, when its MARK told.
	grid *cellRef
	// kept is set for an element with content, or with an element that
	// has some.
	kept bool
}

// structKid is what an element holds: another element, a piece of marked
// content of a stream (of a page, or of a form drawn on it), or an
// annotation.
type structKid struct {
	elem  *structElem
	mcid  int
	s     *stream
	page  *pageOut
	annot *ref
}

// Kinds of the containers that MARKs open.
const (
	nodeList = iota
	nodeItem
	nodeTable
	nodeCell
	nodeFigure
)

// tagNode is an open container.
type tagNode struct {
	kind int
	// elem is the element of the container, and into the element that
	// what is in it goes (the body of an item of a list). A figure
	// without an alternative text is decoration: it has no element.
	elem, into *structElem
	// the row of a table that cells are added to
	row      *structElem
	rowIndex int
}

// leafKind is what a leaf MARK said the next run of text is of.
type leafKind struct {
	tag  string
	cell *cellRef // a cell of a sheet
}

// cellRef is the payload of a CELL: a reference or a range, of a header
// or not.
type cellRef struct {
	row, col, rows, cols int
	scope                string
}

// parseCellRef reads the payload of a CELL ("B2", "A4:B4 col").
func parseCellRef(payload string) *cellRef {
	s := strings.TrimSpace(payload)
	scope := ""
	if i := strings.LastIndexAny(s, " \t"); i >= 0 {
		switch tail := s[i+1:]; tail {
		case "col", "row":
			scope, s = tail, strings.TrimSpace(s[:i])
		default:
			return nil
		}
	}
	one := func(s string) (row, col int, ok bool) {
		i := 0
		for i < len(s) && (s[i] >= 'A' && s[i] <= 'Z' || s[i] >= 'a' && s[i] <= 'z') {
			col = col*26 + int(s[i]|0x20-'a') + 1
			if i++; i > 4 {
				return 0, 0, false
			}
		}
		if i == 0 || i == len(s) || len(s)-i > 8 {
			return 0, 0, false
		}
		n, err := strconv.Atoi(s[i:])
		if err != nil || n < 1 {
			return 0, 0, false
		}
		return n - 1, col - 1, true
	}
	a, b, isRange := strings.Cut(s, ":")
	r0, c0, ok := one(a)
	if !ok {
		return nil
	}
	out := &cellRef{row: r0, col: c0, rows: 1, cols: 1, scope: scope}
	if isRange {
		r1, c1, ok := one(b)
		if !ok {
			return nil
		}
		out.rows, out.cols = max(1, r1-r0+1), max(1, c1-c0+1)
	}
	return out
}

// tagger makes the structure tree of a file, and marks its content.
type tagger struct {
	ex   *exporter
	root *structElem
	lang string // the language of the document

	// the structure of the top-level object being written
	open     []*tagNode
	leaf     *structElem
	reserved *leafKind
	// said is the language of the runs that follow, when a LANG told one
	said string
	// sheet is the table of the cells of a sheet that the last leaf was a
	// cell of: the cell that follows goes on with it
	sheet *tagNode

	// held is the element that everything drawn is the content of: the
	// leaf of a drawing that stands for text
	held *structElem
	// off is set where content is not marked (in soft masks, and in the
	// forms of the layers that pages share), and artifact where all of it
	// is an artifact (in the layers that are not content)
	off, artifact int

	page *pageOut
	// the streams with marked content, by their keys of the parent tree,
	// and the annotations with theirs
	streams []*stream
	annots  []annotTag
	next    int
	elems   int
	// tables are the tables of the tree: their rows are made whole when
	// the tree is written
	tables []*structElem
}

// annotTag is an annotation of the tree: the key of the parent tree that
// leads from it to its element.
type annotTag struct {
	key  int
	elem *structElem
}

// maxStructElems bounds the elements of the tree: what MARKs come after
// them is content of the last elements.
const maxStructElems = 1 << 22

func newTagger(ex *exporter) *tagger {
	t := &tagger{ex: ex, lang: ex.doc.Meta.DC.Language.First()}
	t.root = &structElem{tag: "Document"}
	return t
}

func (t *tagger) elem(tag string, parent *structElem) *structElem {
	e := &structElem{tag: tag, parent: parent}
	parent.kids = append(parent.kids, structKid{elem: e})
	t.elems++
	return e
}

func (t *tagger) top() *tagNode {
	if n := len(t.open); n > 0 {
		return t.open[n-1]
	}
	return nil
}

func (t *tagger) topKind() int {
	if n := t.top(); n != nil {
		return n.kind
	}
	return -1
}

// parent is the element that what comes next is in.
func (t *tagger) parent() *structElem {
	for i := len(t.open) - 1; i >= 0; i-- {
		if e := t.open[i].into; e != nil {
			return e
		}
	}
	return t.root
}

// figure returns the innermost open figure.
func (t *tagger) figure() *tagNode {
	for i := len(t.open) - 1; i >= 0; i-- {
		if t.open[i].kind == nodeFigure {
			return t.open[i]
		}
	}
	return nil
}

func (t *tagger) push(n *tagNode) {
	t.open = append(t.open, n)
	t.leaf, t.reserved, t.sheet = nil, nil, nil
}

func (t *tagger) pop() {
	if n := len(t.open); n > 0 {
		t.open = t.open[:n-1]
	}
	t.leaf, t.reserved, t.sheet = nil, nil, nil
}

// begin starts the structure of a top-level object: no container is
// open, and the language is that of the document.
func (t *tagger) begin() {
	t.open, t.leaf, t.reserved, t.sheet, t.said = nil, nil, nil, nil, ""
}

// item opens an item of the innermost list.
func (t *tagger) item() {
	li := t.elem("LI", t.top().elem)
	t.push(&tagNode{kind: nodeItem, elem: li, into: t.elem("LBody", li)})
}

// cellAttrs writes the attributes of a cell of a table.
func cellAttrs(c *cellRef) string {
	if c == nil {
		return ""
	}
	var b strings.Builder
	if c.cols > 1 {
		fmt.Fprintf(&b, " /ColSpan %d", c.cols)
	}
	if c.rows > 1 {
		fmt.Fprintf(&b, " /RowSpan %d", c.rows)
	}
	switch c.scope {
	case "col":
		b.WriteString(" /Scope /Column")
	case "row":
		b.WriteString(" /Scope /Row")
	}
	return b.String()
}

// cell adds a cell to a table: in the row of its reference, a new row
// when the reference is of another one.
func (t *tagger) cell(tbl *tagNode, c *cellRef) *structElem {
	if tbl.row == nil || c != nil && c.row != tbl.rowIndex {
		tbl.row = t.elem("TR", tbl.elem)
		if c != nil {
			tbl.rowIndex = c.row
		}
	}
	tag := "TD"
	if c != nil && c.scope != "" {
		tag = "TH"
	}
	e := t.elem(tag, tbl.row)
	e.attrs, e.grid = cellAttrs(c), c
	return e
}

// table opens a table.
func (t *tagger) table() *structElem {
	e := t.elem("Table", t.parent())
	t.tables = append(t.tables, e)
	return e
}

// regular makes the rows of a table whole: every row with a cell for each
// column, in the order of the columns. The MARKs of a table tell the
// cells that have text, by their places; a table of a file has every
// cell, and the cells that span rows and columns say so. Rows are in the
// order of their places, and a cell spans the rows that the table has of
// those it spans.
func (t *tagger) regular(tbl *structElem) {
	type slot struct {
		e *structElem
		c cellRef
	}
	var cells []slot
	var rest []structKid // what is not a row: the caption
	var trs []*structElem
	for _, k := range tbl.kids {
		if k.elem == nil || k.elem.tag != "TR" {
			rest = append(rest, k)
			continue
		}
		trs = append(trs, k.elem)
		for _, ck := range k.elem.kids {
			if ck.elem == nil || ck.elem.grid == nil {
				return // a cell without a place: the rows are as they came
			}
			cells = append(cells, slot{ck.elem, *ck.elem.grid})
		}
	}
	if len(cells) == 0 {
		return
	}
	const maxSide = 1 << 12 // (of a table that is made whole: a larger one is left as it came)
	minC, maxC := cells[0].c.col, 0
	rowSet := map[int]bool{}
	for _, s := range cells {
		minC, maxC = min(minC, s.c.col), max(maxC, s.c.col+s.c.cols-1)
		rowSet[s.c.row] = true
	}
	if maxC-minC >= maxSide || len(rowSet) > maxSide {
		return
	}
	rows := make([]int, 0, len(rowSet))
	for r := range rowSet {
		rows = append(rows, r)
	}
	sort.Ints(rows)
	start := map[[2]int]*slot{}
	taken := map[[2]int]bool{}
	for i := range cells {
		s := &cells[i]
		if taken[[2]int{s.c.row, s.c.col}] {
			continue // under a cell that came before
		}
		start[[2]int{s.c.row, s.c.col}] = s
		for _, r := range rows {
			if r < s.c.row || r >= s.c.row+s.c.rows {
				continue
			}
			for c := s.c.col; c < s.c.col+s.c.cols; c++ {
				taken[[2]int{r, c}] = true
			}
		}
	}
	tbl.kids = rest
	for i, r := range rows {
		var tr *structElem
		if i < len(trs) {
			tr = trs[i]
			tr.kids = nil
		} else {
			tr = &structElem{tag: "TR"}
		}
		tr.parent, tr.kept = tbl, true
		tbl.kids = append(tbl.kids, structKid{elem: tr})
		for c := minC; c <= maxC; c++ {
			s := start[[2]int{r, c}]
			if s == nil {
				if !taken[[2]int{r, c}] {
					// a cell with nothing in it
					tr.kids = append(tr.kids, structKid{elem: &structElem{tag: "TD", parent: tr, kept: true}})
				}
				continue
			}
			// the rows it spans, of those the table has
			span := s.c
			span.rows = 0
			for _, rr := range rows {
				if rr >= s.c.row && rr < s.c.row+s.c.rows {
					span.rows++
				}
			}
			s.e.parent, s.e.kept, s.e.attrs = tr, true, cellAttrs(&span)
			tr.kids = append(tr.kids, structKid{elem: s.e})
			c += s.c.cols - 1
		}
	}
}

// mark takes a MARK: a leaf MARK says what the next run is of, a
// structure MARK opens or closes a container (docs/spec.md §7.8).
func (t *tagger) mark(kind byte, payload string) {
	if t.off > 0 || t.artifact > 0 {
		return
	}
	if t.elems >= maxStructElems {
		return
	}
	switch kind {
	case bdf.MarkParagraph, bdf.MarkBox:
		t.reserved = &leafKind{tag: "P"}
	case bdf.MarkHeading:
		level := 1
		if n, err := strconv.Atoi(strings.TrimSpace(payload)); err == nil && n >= 1 && n <= 6 {
			level = n
		}
		t.reserved = &leafKind{tag: "H" + strconv.Itoa(level)}
	case bdf.MarkList:
		t.push(&tagNode{kind: nodeList, elem: t.elem("L", t.parent())})
	case bdf.MarkListItem:
		if t.topKind() == nodeItem {
			t.pop()
		}
		if t.topKind() == nodeList {
			t.item()
		} else {
			t.reserved = &leafKind{tag: "P"}
		}
	case bdf.MarkTable:
		t.push(&tagNode{kind: nodeTable, elem: t.table()})
	case bdf.MarkCell:
		c := parseCellRef(payload)
		if t.topKind() == nodeCell {
			t.pop()
		}
		if tbl := t.top(); tbl != nil && tbl.kind == nodeTable {
			td := t.cell(tbl, c)
			t.push(&tagNode{kind: nodeCell, elem: td, into: td})
		} else if c != nil {
			t.reserved = &leafKind{cell: c}
		} else {
			t.reserved = &leafKind{tag: "P"}
		}
	case bdf.MarkFigure:
		n := &tagNode{kind: nodeFigure}
		if payload != "" {
			n.elem = t.elem("Figure", t.parent())
			n.elem.alt = payload
			t.setLang(n.elem)
		}
		t.push(n)
	case bdf.MarkEnd:
		if k := t.topKind(); k == nodeItem || k == nodeCell {
			t.pop()
		}
		if len(t.open) > 0 {
			t.pop()
		}
	case bdf.MarkLang:
		if t.said != payload {
			// the runs that follow are of another element
			t.said = payload
			if t.leaf != nil && t.reserved == nil {
				t.reserved = &leafKind{tag: "Span"}
			}
		}
	}
}

// leafElem is the element of the run of text that comes next: the one
// that the last leaf MARK told, or a paragraph.
func (t *tagger) leafElem() *structElem {
	if t.leaf != nil && t.reserved == nil {
		return t.leaf
	}
	r := t.reserved
	t.reserved = nil
	if t.elems >= maxStructElems && t.leaf != nil {
		return t.leaf
	}
	// text right in a list is an item of it, right in a table its caption
	if t.topKind() == nodeList {
		t.item()
	}
	switch top := t.top(); {
	case top != nil && top.kind == nodeTable:
		t.leaf = t.elem("Caption", top.elem)
	case r != nil && r.cell != nil:
		// a cell of a sheet: of the table that the cells before it made
		if t.sheet == nil {
			t.sheet = &tagNode{kind: nodeTable, elem: t.table()}
		}
		t.leaf = t.cell(t.sheet, r.cell)
	case r != nil && r.tag == "Span" && t.leaf != nil:
		// the same paragraph, in another language
		p := t.leaf
		for p.tag == "Span" && p.parent != nil {
			p = p.parent
		}
		t.leaf = t.elem("Span", p)
	default:
		tag := "P"
		if r != nil && r.tag != "" && r.tag != "Span" {
			tag = r.tag
		}
		t.leaf, t.sheet = t.elem(tag, t.parent()), nil
	}
	t.setLang(t.leaf)
	return t.leaf
}

// setLang gives an element the language of the runs that follow, when it
// is not the one the element has from those around it.
func (t *tagger) setLang(e *structElem) {
	want := t.said
	if want == "" {
		want = t.lang
	}
	have := t.lang
	for p := e.parent; p != nil; p = p.parent {
		if p.lang != "" {
			have = p.lang
			break
		}
	}
	if want != "" && want != have {
		e.lang = want
	}
}

// link makes the element of a LINK: in the leaf it is in, or in the
// container around it.
func (t *tagger) link() *structElem {
	if t.off > 0 || t.artifact > 0 || t.elems >= maxStructElems {
		return nil
	}
	if f := t.figure(); f != nil && f.elem == nil {
		return nil
	}
	parent := t.parent()
	if t.leaf != nil && t.reserved == nil {
		parent = t.leaf
		for parent.tag == "Span" && parent.parent != nil {
			parent = parent.parent
		}
	} else if k := t.topKind(); k == nodeList || k == nodeTable {
		// (a list and a table hold items and rows alone)
		return nil
	}
	return t.elem("Link", parent)
}

// --- marked content ---

// mcState is the marked content that is open in a stream.
type mcState struct {
	open  bool
	elem  *structElem // nil: an artifact
	depth int         // the q of the stream that it was opened in
}

// q and Q write the operators that save and restore the graphics state.
// Marked content does not go across them: content that was opened inside
// the state that Q ends is closed first.
func (s *stream) q() {
	s.op("q")
	s.qd++
}

func (s *stream) Q() {
	if s.mc.open && s.qd <= s.mc.depth {
		s.endMC()
	}
	s.op("Q")
	s.qd--
}

func (s *stream) endMC() {
	if s.mc.open {
		s.op("EMC")
		s.mc = mcState{}
	}
}

// tag marks what is drawn next: as the content of the element it is of,
// or as an artifact. text says that it is a run of text.
func (ex *exporter) tag(text bool) {
	t := ex.tags
	if t == nil || t.off > 0 {
		return
	}
	var want *structElem
	switch {
	case t.artifact > 0:
	case t.held != nil:
		want = t.held
	default:
		if f := t.figure(); f != nil {
			want = f.elem
		} else if text {
			want = t.leafElem()
		}
	}
	ex.s.setMC(t, want)
}

// tagArtifact marks what is drawn next as an artifact.
func (ex *exporter) tagArtifact() {
	if t := ex.tags; t != nil && t.off == 0 {
		ex.s.setMC(t, nil)
	}
}

// setMC makes the open marked content of the stream that of an element
// (nil: an artifact).
func (s *stream) setMC(t *tagger, want *structElem) {
	if s.mc.open && s.mc.elem == want {
		return
	}
	s.endMC()
	s.mc = mcState{open: true, elem: want, depth: s.qd}
	if want == nil {
		s.op("/Artifact BMC")
		return
	}
	if s.key < 0 {
		s.key = t.next
		t.next++
		t.streams = append(t.streams, s)
	}
	mcid := len(s.mcids)
	s.mcids = append(s.mcids, want)
	want.kids = append(want.kids, structKid{mcid: mcid, s: s, page: t.page})
	if want.page == nil {
		want.page = t.page
	}
	s.op("/%s << /MCID %d >> BDC", want.tag, mcid)
}

// --- the tree ---

// written reports whether the content that a kid stands for is in the
// file: that of a page, of a form that was written (a group that shows
// nothing is not), or an annotation.
func (k *structKid) written() bool {
	return k.annot != nil || k.s != nil && (k.s.page || k.s.ref != 0)
}

// keep marks the elements that hold content, and returns whether the
// element does.
func (e *structElem) keep() bool {
	for i := range e.kids {
		k := &e.kids[i]
		if k.elem == nil && k.written() || k.elem != nil && k.elem.keep() {
			e.kept = true
		}
	}
	return e.kept
}

// write writes the structure tree and returns the entries of the catalog
// for it.
func (t *tagger) write() string {
	pw := t.ex.pw
	t.root.keep()
	t.root.kept = true
	for _, tbl := range t.tables {
		if tbl.kept {
			t.regular(tbl)
		}
	}
	var alloc func(e *structElem)
	alloc = func(e *structElem) {
		e.ref = pw.alloc()
		for _, k := range e.kids {
			if k.elem != nil && k.elem.kept {
				alloc(k.elem)
			}
		}
	}
	alloc(t.root)
	treeRef := pw.alloc()
	var put func(e *structElem, parent ref)
	put = func(e *structElem, parent ref) {
		var b strings.Builder
		fmt.Fprintf(&b, "<< /Type /StructElem /S /%s /P %s", e.tag, parent)
		if e.page != nil {
			fmt.Fprintf(&b, " /Pg %s", e.page.ref)
		}
		if e.lang != "" {
			b.WriteString(" /Lang " + text(e.lang))
		}
		if e.alt != "" {
			b.WriteString(" /Alt " + text(e.alt))
		}
		if e.attrs != "" {
			b.WriteString(" /A << /O /Table" + e.attrs + " >>")
		}
		b.WriteString(" /K [")
		for i := range e.kids {
			k := &e.kids[i]
			switch {
			case k.elem != nil:
				if k.elem.kept {
					b.WriteString(k.elem.ref.String() + " ")
				}
			case !k.written():
			case k.annot != nil:
				fmt.Fprintf(&b, "<< /Type /OBJR /Obj %s /Pg %s >> ", *k.annot, k.page.ref)
			case !k.s.page:
				// (content of a form that is drawn on the page)
				fmt.Fprintf(&b, "<< /Type /MCR /Pg %s /Stm %s /MCID %d >> ", k.page.ref, k.s.ref, k.mcid)
			case k.page == e.page:
				b.WriteString(strconv.Itoa(k.mcid) + " ")
			default:
				fmt.Fprintf(&b, "<< /Type /MCR /Pg %s /MCID %d >> ", k.page.ref, k.mcid)
			}
		}
		b.WriteString("] >>")
		pw.put(e.ref, b.String())
		for _, k := range e.kids {
			if k.elem != nil && k.elem.kept {
				put(k.elem, e.ref)
			}
		}
	}
	put(t.root, treeRef)
	// the parent tree: from the marked content of each stream, and from
	// each annotation, to its element
	var nums strings.Builder
	for _, s := range t.streams {
		if !s.page && s.ref == 0 {
			continue
		}
		var arr strings.Builder
		for _, e := range s.mcids {
			arr.WriteString(e.ref.String() + " ")
		}
		fmt.Fprintf(&nums, "%d %s ", s.key, pw.add("["+arr.String()+"]"))
	}
	for _, a := range t.annots {
		fmt.Fprintf(&nums, "%d %s ", a.key, a.elem.ref)
	}
	parents := pw.add("<< /Nums [" + nums.String() + "] >>")
	pw.put(treeRef, fmt.Sprintf("<< /Type /StructTreeRoot /K [%s] /ParentTree %s /ParentTreeNextKey %d >>", t.root.ref, parents, t.next))
	return " /StructTreeRoot " + treeRef.String() + " /MarkInfo << /Marked true >>"
}
