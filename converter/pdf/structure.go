package pdf

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/shibukawa/bdf"
)

// Tagged PDF keeps the reading structure in a tree of structure elements
// (catalog /StructTreeRoot) whose leaves are marked-content sequences of the
// page content, identified by MCID. The converter maps what BDF can express
// (docs/spec.md §7.8: headings, paragraphs, lists, tables with cell refs,
// figures with alternative text, languages) to MARKs.
//
// A marked-content sequence resolves to a target: the containers enclosing
// it (LIST with its current LIST_ITEM, TABLE with its current CELL, FIGURE),
// the element whose paragraph its runs belong to, and its language. MARKs
// are emitted lazily by comparing the target with the structure the emitted
// instructions leave open (structState): before a text run, and before a
// drawing only to open the figures it belongs to or close those it does not.
// So sequences that draw no text (backgrounds, borders) and content outside
// the structure (artifacts, no MCID) leave the structure as it is, except
// that they never count as part of a figure. A form XObject continues its
// caller's walk; its MCIDs resolve through its own /StructParents or the
// page's.

const structMaxDepth = 64

type structTree struct {
	p       *pdf
	root    types.Dict
	roleMap types.Dict
	docLang string
	elems   map[string]*structElem
	targets map[*structElem]*mcTarget
	cells   map[*structElem]map[*structElem]string // table → cell → CELL payload
	tables  map[int]*mcidTable                     // by parent tree key
	byPage  map[int]map[int]*structElem            // page object → MCID → element, from walking the tree
}

type structElem struct {
	typ      string // standard structure type after role mapping
	dict     types.Dict
	parent   *structElem
	alt      string // /Alt, or /ActualText when there is no /Alt
	hasAlt   bool
	lang     string
	hasLang  bool
	building bool
}

func newStructTree(p *pdf, catalog types.Dict, docLang string) *structTree {
	root := p.dict(catalog["StructTreeRoot"])
	if root == nil {
		return nil
	}
	return &structTree{p: p, root: root, roleMap: p.dict(root["RoleMap"]), docLang: docLang,
		elems: map[string]*structElem{}, targets: map[*structElem]*mcTarget{}, cells: map[*structElem]map[*structElem]string{}, tables: map[int]*mcidTable{}}
}

func structKey(o types.Object) string {
	switch v := o.(type) {
	case types.IndirectRef:
		return "r" + strconv.Itoa(v.ObjectNumber.Value())
	case *types.IndirectRef:
		if v != nil {
			return "r" + strconv.Itoa(v.ObjectNumber.Value())
		}
	case types.Dict:
		return fmt.Sprintf("d%p", v)
	}
	return ""
}

// elem returns the structure element o refers to (nil if it is none).
func (t *structTree) elem(o types.Object) *structElem { return t.elemAt(o, 0) }

func (t *structTree) elemAt(o types.Object, depth int) *structElem {
	key := structKey(o)
	if key == "" {
		return nil
	}
	if e, ok := t.elems[key]; ok {
		if e != nil && e.building {
			return nil // /P cycle
		}
		return e
	}
	d := t.p.dict(o)
	if d == nil || d["S"] == nil || t.p.name(d["Type"]) == "StructTreeRoot" {
		t.elems[key] = nil
		return nil
	}
	e := &structElem{typ: t.stdType(t.p.name(d["S"])), dict: d, building: true}
	t.elems[key] = e
	if v, ok := d["Alt"]; ok {
		e.alt, e.hasAlt = t.p.text(v), true
	} else if v, ok := d["ActualText"]; ok {
		e.alt = t.p.text(v)
	}
	if v, ok := d["Lang"]; ok && t.p.deref(v) != nil {
		e.lang, e.hasLang = t.p.text(v), true
	}
	if depth < structMaxDepth {
		e.parent = t.elemAt(d["P"], depth+1)
	}
	e.building = false
	return e
}

var stdStructTypes = map[string]bool{}

func init() {
	for _, s := range strings.Fields(`Document DocumentFragment Part Art Sect Div BlockQuote Caption TOC TOCI Index
		NonStruct Private Aside Title FENote Sub P H L LI Lbl LBody Table TR TH TD THead TBody TFoot Span Quote Note
		Reference BibEntry Code Link Annot Ruby RB RT RP Warichu WT WP Figure Formula Form Em Strong Artifact`) {
		stdStructTypes[s] = true
	}
}

// headingLevel returns n for Hn (clamped to 6), or 0.
func headingLevel(typ string) int {
	if len(typ) < 2 || typ[0] != 'H' {
		return 0
	}
	n, err := strconv.Atoi(typ[1:])
	if err != nil || n < 1 {
		return 0
	}
	return min(n, 6)
}

// stdType maps a structure type through the role map to a standard type.
func (t *structTree) stdType(s string) string {
	for i := 0; i < 16; i++ {
		if stdStructTypes[s] || headingLevel(s) > 0 {
			return s
		}
		m := t.p.name(t.roleMap[s])
		if m == "" || m == s {
			return s
		}
		s = m
	}
	return s
}

// Structure element classes.
const (
	clsInline  = iota // transparent: Span, Link, NonStruct, Lbl, LBody, TR, THead, unknown types …
	clsGroup          // Document, Sect, Div …: the paragraph of content nothing finer encloses
	clsPara           // P, Caption, BlockQuote …
	clsHeading        // H, H1-H6, Title
	clsList           // L
	clsItem           // LI
	clsTable          // Table
	clsCell           // TH, TD
	clsFigure         // Figure, Formula with /Alt
)

func (e *structElem) class() int {
	switch e.typ {
	case "Document", "DocumentFragment", "Part", "Art", "Sect", "Div", "TOC", "Index", "Aside":
		return clsGroup
	case "P", "Caption", "BlockQuote", "Note", "FENote", "TOCI":
		return clsPara
	case "H", "Title":
		return clsHeading
	case "L":
		return clsList
	case "LI":
		return clsItem
	case "Table":
		return clsTable
	case "TH", "TD":
		return clsCell
	case "Figure":
		return clsFigure
	case "Formula":
		if e.hasAlt {
			return clsFigure
		}
		return clsPara
	}
	if headingLevel(e.typ) > 0 {
		return clsHeading
	}
	return clsInline
}

// structFrame is an open LIST, TABLE or FIGURE and its current item or cell.
type structFrame struct {
	kind        byte // bdf.MarkList, bdf.MarkTable or bdf.MarkFigure
	id          *structElem
	payload     string
	item        *structElem // LI or TH/TD
	itemPayload string
}

// mcTarget is the structure a marked-content sequence belongs to. A nil
// target is content outside the structure in the document language.
type mcTarget struct {
	tagged      bool // resolved to a structure element
	frames      []structFrame
	leaf        *structElem // element whose paragraph the runs belong to
	leafKind    byte
	leafPayload string
	leafID      *structElem // leaf, or the innermost frame's item or container
	hasFigure   bool
	lang        string // "" = document default
	langSet     bool   // lang comes from the element or the sequence
}

func (t *mcTarget) language() string {
	if t == nil {
		return ""
	}
	return t.lang
}

func sameTarget(a, b *mcTarget) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		x := a
		if x == nil {
			x = b
		}
		return !x.tagged && x.lang == ""
	}
	if a.tagged != b.tagged || a.lang != b.lang || a.leaf != b.leaf || a.leafID != b.leafID || len(a.frames) != len(b.frames) {
		return false
	}
	for i := range a.frames {
		if a.frames[i].id != b.frames[i].id || a.frames[i].item != b.frames[i].item {
			return false
		}
	}
	return true
}

func positionID(frames []structFrame) *structElem {
	if n := len(frames); n > 0 {
		if frames[n-1].item != nil {
			return frames[n-1].item
		}
		return frames[n-1].id
	}
	return nil
}

// target returns the target of content whose innermost element is e.
func (t *structTree) target(e *structElem) *mcTarget {
	if tg, ok := t.targets[e]; ok {
		return tg
	}
	var chain []*structElem
	for x := e; x != nil && len(chain) < structMaxDepth; x = x.parent {
		chain = append(chain, x)
	}
	tg := &mcTarget{tagged: true}
	var leaf, group *structElem
	for i := len(chain) - 1; i >= 0; i-- {
		x := chain[i]
		top := len(tg.frames) - 1
		switch x.class() {
		case clsList:
			tg.frames = append(tg.frames, structFrame{kind: bdf.MarkList, id: x})
			leaf, group = nil, nil
		case clsTable:
			tg.frames = append(tg.frames, structFrame{kind: bdf.MarkTable, id: x})
			leaf, group = nil, nil
		case clsFigure:
			tg.frames = append(tg.frames, structFrame{kind: bdf.MarkFigure, id: x, payload: x.alt})
			tg.hasFigure = true
			leaf, group = nil, nil
		case clsItem:
			if top >= 0 && tg.frames[top].kind == bdf.MarkList && tg.frames[top].item == nil {
				tg.frames[top].item = x
				leaf, group = nil, nil
			} else {
				group = x
			}
		case clsCell:
			if top >= 0 && tg.frames[top].kind == bdf.MarkTable && tg.frames[top].item == nil {
				if ref, ok := t.cellRef(tg.frames[top].id, x); ok {
					tg.frames[top].item, tg.frames[top].itemPayload = x, ref
					leaf, group = nil, nil
					break
				}
			}
			group = x
		case clsPara, clsHeading:
			leaf = x
		case clsGroup:
			group = x
		}
		if x.hasLang {
			tg.lang, tg.langSet = t.normLang(x.lang), true
		}
	}
	if leaf == nil {
		leaf = group
	}
	if leaf != nil {
		tg.leaf, tg.leafID, tg.leafKind, tg.leafPayload = leaf, leaf, bdf.MarkParagraph, leaf.typ
		if leaf.class() == clsHeading {
			tg.leafKind, tg.leafPayload = bdf.MarkHeading, strconv.Itoa(max(headingLevel(leaf.typ), 1))
		}
	} else {
		tg.leafID = positionID(tg.frames)
	}
	t.targets[e] = tg
	return tg
}

func (t *structTree) normLang(s string) string { return normLang(s, t.docLang) }

// normLang returns a language tag, "" when it is the document default.
func normLang(s, docLang string) string {
	s = strings.TrimSpace(s)
	if strings.EqualFold(s, docLang) {
		return ""
	}
	return s
}

// kids returns the /K entries of a structure element or the tree root.
func (t *structTree) kids(d types.Dict) types.Array {
	if a := t.p.array(d["K"]); a != nil {
		return a
	}
	if d["K"] == nil {
		return nil
	}
	return types.Array{d["K"]}
}

// cellRef returns the CELL payload of a TH/TD of table.
func (t *structTree) cellRef(table, cell *structElem) (string, bool) {
	refs, ok := t.cells[table]
	if !ok {
		refs = t.tableGrid(table)
		t.cells[table] = refs
	}
	r, ok := refs[cell]
	return r, ok
}

// tableGrid lays out the cells of a table in structure order (rows of TR,
// possibly grouped in THead/TBody/TFoot) with their row and column spans.
func (t *structTree) tableGrid(table *structElem) map[*structElem]string {
	var rows [][]*structElem
	var collect func(e *structElem, depth int)
	collect = func(e *structElem, depth int) {
		for _, k := range t.kids(e.dict) {
			x := t.elem(k)
			if x == nil {
				continue
			}
			switch x.typ {
			case "THead", "TBody", "TFoot":
				if depth < 4 {
					collect(x, depth+1)
				}
			case "TR":
				var cells []*structElem
				for _, ck := range t.kids(x.dict) {
					if c := t.elem(ck); c != nil && c.class() == clsCell {
						cells = append(cells, c)
					}
				}
				rows = append(rows, cells)
			}
		}
	}
	collect(table, 0)
	refs := map[*structElem]string{}
	used := map[[2]int]bool{}
	budget := 1 << 20
	for r, cells := range rows {
		c := 0
		for _, cell := range cells {
			for used[[2]int{r, c}] {
				c++
			}
			rs, cs, scope := 1, 1, ""
			for _, a := range t.attrs(cell) {
				rs = max(rs, t.p.intOr(a["RowSpan"], 1))
				cs = max(cs, t.p.intOr(a["ColSpan"], 1))
				if s := t.p.name(a["Scope"]); s != "" {
					scope = s
				}
			}
			rs, cs = min(rs, len(rows)-r), min(cs, 1024)
			if budget -= rs * cs; budget < 0 {
				rs, cs = 1, 1
			}
			for dr := 0; dr < rs; dr++ {
				for dc := 0; dc < cs; dc++ {
					used[[2]int{r + dr, c + dc}] = true
				}
			}
			ref := cellName(c, r)
			if rs > 1 || cs > 1 {
				ref += ":" + cellName(c+cs-1, r+rs-1)
			}
			if cell.typ == "TH" {
				if scope == "Row" {
					ref += " row"
				} else {
					ref += " col"
				}
			}
			refs[cell] = ref
			c += cs
		}
	}
	return refs
}

// cellName returns the A1-style name of a zero-based column and row.
func cellName(col, row int) string {
	var b []byte
	for col++; col > 0; col = (col - 1) / 26 {
		b = append([]byte{byte('A' + (col-1)%26)}, b...)
	}
	return string(b) + strconv.Itoa(row+1)
}

// attrs returns the attribute dictionaries of an element (/A: a dictionary
// or an array of dictionaries and revision numbers).
func (t *structTree) attrs(e *structElem) []types.Dict {
	if d := t.p.dict(e.dict["A"]); d != nil {
		return []types.Dict{d}
	}
	var out []types.Dict
	for _, a := range t.p.array(e.dict["A"]) {
		if d := t.p.dict(a); d != nil {
			out = append(out, d)
		}
	}
	return out
}

// mcidTable maps the MCIDs of a content stream to structure elements.
type mcidTable struct {
	t   *structTree
	arr types.Array
	m   map[int]*structElem
}

func (m *mcidTable) elem(mcid int) *structElem {
	if m == nil {
		return nil
	}
	if m.arr != nil {
		if mcid >= 0 && mcid < len(m.arr) {
			return m.t.elem(m.arr[mcid])
		}
		return nil
	}
	return m.m[mcid]
}

// parentTreeValue looks key up in the parent tree.
func (t *structTree) parentTreeValue(key int) types.Object {
	return t.p.numberTree(t.root["ParentTree"], key)
}

// contentTable returns the MCID table of a content stream with /StructParents
// key (hasKey), falling back for pages to the MCIDs the tree assigns to the
// page object.
func (t *structTree) contentTable(key int, hasKey bool, pageObj int) *mcidTable {
	if hasKey {
		if m, ok := t.tables[key]; ok {
			return m
		}
		if arr := t.p.array(t.parentTreeValue(key)); arr != nil {
			m := &mcidTable{t: t, arr: arr}
			t.tables[key] = m
			return m
		}
	}
	if pageObj > 0 {
		if t.byPage == nil {
			t.byPage = map[int]map[int]*structElem{}
			t.walk(t.root, nil, 0, map[*structElem]bool{}, 0)
		}
		if m := t.byPage[pageObj]; len(m) > 0 {
			return &mcidTable{t: t, m: m}
		}
	}
	return nil
}

// walk collects the MCIDs of the elements under d by page.
func (t *structTree) walk(d types.Dict, e *structElem, pg int, seen map[*structElem]bool, depth int) {
	if depth > structMaxDepth {
		return
	}
	if r, ok := d["Pg"].(types.IndirectRef); ok {
		pg = r.ObjectNumber.Value()
	}
	add := func(pg, mcid int) {
		if e == nil || pg == 0 {
			return
		}
		if t.byPage[pg] == nil {
			t.byPage[pg] = map[int]*structElem{}
		}
		t.byPage[pg][mcid] = e
	}
	for _, k := range t.kids(d) {
		switch v := t.p.deref(k).(type) {
		case types.Integer:
			add(pg, int(v))
		case types.Dict:
			switch t.p.name(v["Type"]) {
			case "MCR":
				if v["Stm"] == nil {
					mpg := pg
					if r, ok := v["Pg"].(types.IndirectRef); ok {
						mpg = r.ObjectNumber.Value()
					}
					add(mpg, t.p.intOr(v["MCID"], -1))
				}
			case "OBJR":
			default:
				if x := t.elem(k); x != nil && !seen[x] {
					seen[x] = true
					t.walk(v, x, pg, seen, depth+1)
				}
			}
		}
	}
}

// structState is the structure that the instructions emitted so far leave
// open: what a reader walking them has at this point.
type structState struct {
	frames []structFrame
	leaf   *structElem
	lang   string
}

func (s structState) clone() structState {
	s.frames = append([]structFrame(nil), s.frames...)
	return s
}

func (s structState) equal(o structState) bool {
	if s.leaf != o.leaf || s.lang != o.lang || len(s.frames) != len(o.frames) {
		return false
	}
	for i := range s.frames {
		if s.frames[i].id != o.frames[i].id || s.frames[i].item != o.frames[i].item {
			return false
		}
	}
	return true
}

func (s structState) hasFigure() bool {
	for _, f := range s.frames {
		if f.kind == bdf.MarkFigure {
			return true
		}
	}
	return false
}

// formStruct records the structure a form XObject was converted under: its
// MARKs are only right for a caller in the same state.
type formStruct struct {
	mcids   *mcidTable
	inherit *mcTarget
	entry   structState
	exit    structState
	used    bool // the content depends on the entry state (it has text or MARKs)
}

func (f *formStruct) reusableFor(g *formStruct) bool {
	return g != nil && f.mcids == g.mcids && sameTarget(f.inherit, g.inherit) && f.entry.equal(g.entry)
}

// --- interpreter side ---

type mcEntry struct {
	actual *actualText
	tgt    *mcTarget
	set    bool // tgt replaces the enclosing target
	hidden bool // optional content the default configuration hides
}

// curTarget returns the target of the content being interpreted.
func (in *interp) curTarget() *mcTarget {
	for i := len(in.mcStack) - 1; i >= 0; i-- {
		if in.mcStack[i].set {
			return in.mcStack[i].tgt
		}
	}
	return in.inherit
}

// markedTarget returns the target of a marked-content sequence, or false
// when it does not change the enclosing one.
func (in *interp) markedTarget(tag string, d types.Dict, text func(types.Object) string) (*mcTarget, bool) {
	outer := in.curTarget()
	var t *mcTarget
	set := false
	if tag == "Artifact" {
		t, set = &mcTarget{}, true
	} else if tree := in.c.tree; tree != nil {
		if mcid, ok := in.c.pdf.num(d["MCID"]); ok {
			if e := in.mcids.elem(int(mcid)); e != nil {
				t, set = tree.target(e), true
			}
		}
	}
	if v, ok := d["Lang"]; ok {
		nt := mcTarget{}
		if set {
			nt = *t
		} else if outer != nil {
			nt = *outer
		}
		nt.lang, nt.langSet = normLang(text(v), in.c.docLang), true
		return &nt, true
	}
	return withLang(t, outer), set
}

// withLang gives t the language of outer unless its elements set one.
func withLang(t, outer *mcTarget) *mcTarget {
	if t == nil || t.langSet || t.lang == outer.language() {
		return t
	}
	nt := *t
	nt.lang = outer.language()
	return &nt
}

// targetOf returns the target of an XObject: the element it belongs to as a
// whole (/StructParent), or the enclosing marked content.
func (in *interp) targetOf(d types.Dict) *mcTarget {
	cur := in.curTarget()
	if in.c.tree == nil {
		return cur
	}
	sp, ok := in.c.pdf.num(d["StructParent"])
	if !ok {
		return cur
	}
	e := in.c.tree.elem(in.c.tree.parentTreeValue(int(sp)))
	if e == nil {
		return cur
	}
	return withLang(in.c.tree.target(e), cur)
}

func (in *interp) mark(kind byte, payload string) {
	in.obj.Mark(kind, payload)
	in.structUsed = true
}

// closeFrames emits END for the open frames beyond n.
func (in *interp) closeFrames(n int) {
	st := in.st
	if len(st.frames) <= n {
		return
	}
	for len(st.frames) > n {
		in.mark(bdf.MarkEnd, "")
		st.frames = st.frames[:len(st.frames)-1]
	}
	st.leaf = positionID(st.frames)
}

// commonFrames returns how many open frames the target keeps as they are.
func (in *interp) commonFrames(tf []structFrame) int {
	st := in.st
	i := 0
	for i < len(st.frames) && i < len(tf) && st.frames[i].id == tf[i].id && (tf[i].item == nil || st.frames[i].item == tf[i].item) {
		i++
	}
	return i
}

// syncFrames closes and opens frames so that the open ones are tf.
func (in *interp) syncFrames(tf []structFrame) {
	st := in.st
	i := in.commonFrames(tf)
	if i == len(st.frames) && i == len(tf) {
		return
	}
	if i < len(st.frames) && i < len(tf) && st.frames[i].id == tf[i].id {
		// Same list or table, next item or cell.
		in.closeFrames(i + 1)
		in.mark(itemKind(tf[i].kind), tf[i].itemPayload)
		st.frames[i].item, st.frames[i].itemPayload = tf[i].item, tf[i].itemPayload
		i++
	} else {
		in.closeFrames(i)
	}
	for ; i < len(tf); i++ {
		f := tf[i]
		in.mark(f.kind, f.payload)
		if f.item != nil {
			in.mark(itemKind(f.kind), f.itemPayload)
		}
		st.frames = append(st.frames, f)
	}
	st.leaf = positionID(st.frames)
}

func itemKind(kind byte) byte {
	if kind == bdf.MarkTable {
		return bdf.MarkCell
	}
	return bdf.MarkListItem
}

// closeFigures closes the open figures (and what they enclose) that the
// target frames tf do not keep.
func (in *interp) closeFigures(tf []structFrame) {
	for j := in.commonFrames(tf); j < len(in.st.frames); j++ {
		if in.st.frames[j].kind == bdf.MarkFigure {
			in.closeFrames(j)
			return
		}
	}
}

// syncText emits the MARKs that put the next text run in structure t.
func (in *interp) syncText(t *mcTarget) {
	st := in.st
	if st == nil {
		return
	}
	in.structUsed = true
	if t != nil && t.tagged {
		in.syncFrames(t.frames)
		if t.leafID != st.leaf {
			switch {
			case t.leaf == nil:
				in.mark(bdf.MarkParagraph, "")
			case t.leafKind == bdf.MarkParagraph && st.leaf != nil && st.leaf == positionID(st.frames):
				// The item or cell has just begun (or only its label was
				// read): its first paragraph continues it.
			default:
				in.mark(t.leafKind, t.leafPayload)
			}
			st.leaf = t.leafID
		}
	} else {
		in.closeFigures(nil)
	}
	if lang := t.language(); lang != st.lang {
		in.mark(bdf.MarkLang, lang)
		st.lang = lang
	}
}

// syncDraw emits the MARKs a drawing needs: it opens the figures of its
// target and closes the figures it is not part of. Other structure waits for
// the next text run.
func (in *interp) syncDraw(t *mcTarget) {
	st := in.st
	if st == nil {
		return
	}
	if t != nil && t.hasFigure {
		in.structUsed = true
		in.syncFrames(t.frames)
		return
	}
	if st.hasFigure() {
		in.structUsed = true
		var tf []structFrame
		if t != nil {
			tf = t.frames
		}
		in.closeFigures(tf)
	}
}

// closeStructure ends the open structure at the end of a page.
func (in *interp) closeStructure() {
	if in.st == nil {
		return
	}
	in.closeFrames(0)
	*in.st = structState{}
}

// decodeText converts the bytes of a PDF text string (PDFDocEncoding, or
// UTF-16BE or UTF-8 with a byte order mark) to UTF-8.
func decodeText(b []byte) string {
	if len(b) >= 2 && b[0] == 0xfe && b[1] == 0xff {
		var out []rune
		for i := 2; i+1 < len(b); i += 2 {
			u := rune(b[i])<<8 | rune(b[i+1])
			if u >= 0xd800 && u < 0xdc00 && i+3 < len(b) {
				lo := rune(b[i+2])<<8 | rune(b[i+3])
				u = 0x10000 + (u-0xd800)<<10 + (lo - 0xdc00)
				i += 2
			}
			out = append(out, u)
		}
		return string(out)
	}
	if len(b) >= 3 && b[0] == 0xef && b[1] == 0xbb && b[2] == 0xbf {
		return string(b[3:])
	}
	return pdfDocText(b)
}
