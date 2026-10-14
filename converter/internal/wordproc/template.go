package wordproc

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/form/tag"
)

// A Word document is a template when its text holds tags (package
// form/tag): FillDOCX converts it with the tags replaced by data. The
// document is filled as its tree of elements in memory, before it is read
// and laid out: what the data makes longer flows on as text does, onto
// further pages. No document is written.
//
// A tag in a text is replaced by its value, in the format of the run the
// tag starts in (Word splits a tag into runs where it checked the spelling
// or the text was edited: the runs of a paragraph are read as one text). A
// section repeats the smallest of these that holds both its tags:
//
//   - a part of the text of a paragraph (both tags in the paragraph);
//   - the paragraphs and tables from the one tag to the other, for tags in
//     different paragraphs of the body, of a cell or of a text box;
//   - rows of a table: the row, for tags in different cells of it (the
//     table grows), and for tags in different rows the rows from the one to
//     the other, which are slots, a row for an item (tag.ModeSlots;
//     "| rows" repeats them instead).
//
// How a section repeats, when its tag says it, is what it repeats: "| rows"
// the row of the table it is in or, outside tables, its paragraphs, also
// when both tags are in the same paragraph; "| pages" its paragraphs and
// tables with each item on a new page.
//
// A paragraph that holds nothing but the tags of such sections is left out.
// A picture whose alternative text has a picture tag ({{%photo}}) is
// replaced by the picture the data names, fitted in its frame; the rest of
// the alternative text is that of the new picture.
//
// The body, the headers and footers and the text boxes are filled; the
// tags of a section are in the same one of them. Headers, footers and the
// text boxes outside sections see the values of the data, not the items of
// sections. Footnotes and endnotes are not filled.

// TemplateData is what a template is filled with.
type TemplateData struct {
	// Data is the data: the values of JSON (nil, bool, float64, string,
	// []any, map[string]any), and time.Time for dates.
	Data map[string]any
	// Strict makes a tag whose name the data does not have an error;
	// without it the tag is replaced by nothing, with a warning.
	Strict bool
	// Files holds the pictures that the data names by their paths.
	Files fs.FS
	// Format writes values in the number formats of tags
	// ({{total | fmt:"#,##0"}}); converter/xlsx has one (Formatter).
	// Without it the formats of tags are not applied.
	Format tag.Formatter
}

// Limits of a filled template.
const (
	// maxFilledNodes bounds the elements of a filled template.
	maxFilledNodes = 8 << 20
	// maxParaText bounds the characters of a filled paragraph.
	maxParaText = 1 << 16
	// maxSectionDepth bounds the sections one inside the other.
	maxSectionDepth = 32
	// maxExpansionSteps bounds repetitions that produce no nodes or text.
	maxExpansionSteps = 8 << 20
)

// FillDOCX converts a Word document as a template filled with data. It
// makes the page view alone unless the options ask for more: a form is its
// pages.
func FillDOCX(r io.ReaderAt, size int64, t *TemplateData, opts *Options) (*Result, []tag.Field, error) {
	o := Options{}
	if opts != nil {
		o = *opts
	}
	if o.Views == "" {
		o.Views = ViewsPages
	}
	if t == nil {
		t = &TemplateData{}
	}
	f := newFiller(t)
	res, err := convertDOCX(r, size, &o, f)
	if err != nil {
		return nil, nil, err
	}
	return res, f.fieldList(), nil
}

// InspectDOCX reads the tags of a Word document without filling it: the
// fields that data for it has, or the first mistake in its tags.
func InspectDOCX(r io.ReaderAt, size int64) (fields []tag.Field, err error) {
	defer func() {
		if e := recover(); e != nil {
			fields, err = nil, fmt.Errorf("docx: internal error: %v", e)
		}
	}()
	p, err := ooxml.Open(r, size)
	if err != nil {
		return nil, fmt.Errorf("docx: %w", err)
	}
	p.Supported = func(prefix string) bool { return supported[prefix] }
	main := "word/document.xml"
	if rel, ok := p.RelOfType("", "/officeDocument"); ok {
		main = rel.Target
	}
	root, err := p.XML(main)
	if err != nil {
		return nil, fmt.Errorf("docx: %w", err)
	}
	f := newFiller(&TemplateData{})
	if err := f.scan(p, main, root); err != nil {
		return nil, err
	}
	return f.fieldList(), nil
}

// step is a place in an element: the kid of it at an index. The steps
// from a story down to a paragraph are those of its parents.
type step struct {
	parent *step
	cont   *ooxml.Node
	idx    int
	depth  int
}

func (s *step) down(cont *ooxml.Node, idx int) *step {
	d := 0
	if s != nil {
		d = s.depth + 1
	}
	return &step{parent: s, cont: cont, idx: idx, depth: d}
}

// tplLoc is where a paragraph is.
type tplLoc struct {
	at   *step  // the paragraph in what holds it
	name string // "table 2, row 3, cell 1"
}

// tplSec is a section that repeats the kids a to b of an element:
// paragraphs and tables, or the rows of a table.
type tplSec struct {
	name     string
	path     []string
	inverted bool
	mode     string
	a, b     int
	at       string
	children []*tplSec
}

// holder is the sections of an element.
type holder struct{ secs []*tplSec }

// add puts a section among those of the holder: the ones inside it become
// its children.
func (h *holder) add(s *tplSec) {
	var rest []*tplSec
	for _, o := range h.secs {
		if o.a >= s.a && o.b <= s.b {
			s.children = append(s.children, o)
		} else {
			rest = append(rest, o)
		}
	}
	h.secs = append(rest, s)
	sort.SliceStable(h.secs, func(i, j int) bool { return h.secs[i].a < h.secs[j].a })
}

// tplItem is a piece of the content of a paragraph: text or a tag, a child
// of a run that is not text (a tab, a drawing), or another child of the
// paragraph (a hyperlink, a bookmark).
type tplItem struct {
	tok   tag.Token
	run   *ooxml.Node // the run whose format the text has, or whose child kid is
	kid   *ooxml.Node
	node  *ooxml.Node
	wrap  *ooxml.Node // the element in node that holds runs, for a hyperlink and the like
	mark  bool        // node is a mark of a place (a bookmark), not content
	outer bool        // a tag of a section beyond the paragraph
}

// inlineTpl is the pieces of a paragraph with tags, or of an element of a
// paragraph that holds runs (a hyperlink).
type inlineTpl struct {
	items []tplItem
	at    string
	// only is set when a paragraph holds nothing but tags of sections
	// beyond it: it is left out
	only bool
}

// openSec is a section whose closing tag has not been read.
type openSec struct {
	tok tag.Token
	loc tplLoc
	use int // its use among those of the fields (-1: none)
}

// fieldUse is a tag that takes a name from the data.
type fieldUse struct {
	fd    tag.Field
	at    tplLoc
	inner []string // the sections of its paragraph around it
	sec   *tplSec  // of a section that reaches beyond its paragraph
}

// picTag is the tag of a picture, and the rest of its alternative text.
type picTag struct {
	tok tag.Token
	alt string
	at  string
}

// story counts the paragraphs and tables of a story, for messages.
type story struct {
	label         string
	paras, tables int
}

type filler struct {
	c      *converter // nil when the template is only read
	root   *tag.Scope
	strict bool
	files  fs.FS
	format tag.Formatter

	inl     map[*ooxml.Node]*inlineTpl
	holders map[*ooxml.Node]*holder
	pics    map[*ooxml.Node]picTag
	uses    []fieldUse
	boxes   int
	nodes   int
	steps   int // repeated items and inline tokens expanded
	// filled are the headers and footers, filled, by their parts
	filled map[string]*ooxml.Node
}

func newFiller(t *TemplateData) *filler {
	data := t.Data
	if data == nil {
		data = map[string]any{}
	}
	return &filler{root: tag.NewScope(tag.Normalize(data)), strict: t.Strict, files: t.Files, format: t.Format,
		inl: map[*ooxml.Node]*inlineTpl{}, holders: map[*ooxml.Node]*holder{}, pics: map[*ooxml.Node]picTag{}, filled: map[string]*ooxml.Node{}}
}

func tplErrorf(at, format string, args ...any) error {
	return fmt.Errorf("docx: %s: %s", at, fmt.Sprintf(format, args...))
}

func (f *filler) expansionStep(at string) error {
	if f.steps >= maxExpansionSteps {
		return tplErrorf(at, "the filled template needs more than %d expansion steps", maxExpansionSteps)
	}
	f.steps++
	return nil
}

// hfParts lists the headers and footers of a document: the parts, and
// what they are called in messages.
func hfParts(p *ooxml.Package, main string) (parts, labels []string) {
	rels := p.Rels(main)
	ids := make([]string, 0, len(rels))
	for id := range rels {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		if len(ids[i]) != len(ids[j]) {
			return len(ids[i]) < len(ids[j])
		}
		return ids[i] < ids[j]
	})
	for _, id := range ids {
		r := rels[id]
		if r.External || !strings.HasSuffix(r.Type, "/header") && !strings.HasSuffix(r.Type, "/footer") {
			continue
		}
		parts = append(parts, r.Target)
		labels = append(labels, strings.TrimSuffix(path.Base(r.Target), ".xml"))
	}
	return parts, labels
}

// scan reads the tags of the document: of its body, then of its headers
// and footers.
func (f *filler) scan(p *ooxml.Package, main string, root *ooxml.Node) error {
	if err := f.scanStory(root.Child("body"), "", nil); err != nil {
		return err
	}
	parts, labels := hfParts(p, main)
	for i, part := range parts {
		n, err := p.XML(part)
		if err != nil {
			continue // told when the part is drawn
		}
		if err := f.scanStory(n, labels[i], nil); err != nil {
			return err
		}
	}
	for _, h := range f.holders {
		if err := f.check(h.secs, nil); err != nil {
			return err
		}
	}
	return nil
}

// fillAll fills the body and the headers and footers of the document, and
// returns the body.
func (f *filler) fillAll(p *ooxml.Package, main string, body *ooxml.Node) (*ooxml.Node, error) {
	out, err := f.clone(body, f.root)
	if err != nil {
		return nil, err
	}
	parts, _ := hfParts(p, main)
	for _, part := range parts {
		n, err := p.XML(part)
		if err != nil {
			continue
		}
		if f.filled[part], err = f.clone(n, f.root); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// scanStory reads the tags of a story (the body, a header, a footer, a
// text box): its sections are closed in it. in is the paragraph that a
// text box is in.
func (f *filler) scanStory(root *ooxml.Node, label string, in *step) error {
	if root == nil {
		return nil
	}
	var stack []openSec
	if err := f.scanBlocks(root, in, &story{label: label}, &stack, ""); err != nil {
		return err
	}
	if n := len(stack); n > 0 {
		o := stack[n-1]
		return tplErrorf(o.loc.name, "the section {{#%s}} is not closed with {{/%s}}", o.tok.Name, o.tok.Name)
	}
	return nil
}

// scanBlocks reads the paragraphs and the tables of an element; cell
// names the cell of a table they are in.
func (f *filler) scanBlocks(cont *ooxml.Node, at *step, st *story, stack *[]openSec, cell string) error {
	for i, k := range cont.Kids {
		here := at.down(cont, i)
		switch k.Name {
		case "p":
			name := cell
			if name == "" {
				st.paras++
				name = "paragraph " + strconv.Itoa(st.paras)
			}
			if st.label != "" {
				name = st.label + ", " + name
			}
			if err := f.scanInline(k, tplLoc{at: here, name: name}, stack, nil); err != nil {
				return err
			}
		case "tbl":
			tname := cell
			if cell == "" {
				st.tables++
				tname = "table " + strconv.Itoa(st.tables)
			}
			row := 0
			for j, tr := range k.Kids {
				if tr.Name != "tr" {
					continue
				}
				row++
				col := 0
				inRow := here.down(k, j)
				for c, tc := range tr.Kids {
					if tc.Name != "tc" {
						continue
					}
					col++
					cname := cell // (a table in a cell: the cell of the table around it)
					if cell == "" {
						cname = tname + ", row " + strconv.Itoa(row) + ", cell " + strconv.Itoa(col)
					}
					if err := f.scanBlocks(tc, inRow.down(tr, c), st, stack, cname); err != nil {
						return err
					}
				}
			}
		case "sdt", "sdtContent", "customXml", "ins", "moveTo":
			if err := f.scanBlocks(k, here, st, stack, cell); err != nil {
				return err
			}
		}
	}
	return nil
}

// scanDeep finds what is filled in a child of a paragraph that is not
// text: the pictures with tags, and the text boxes, which are stories of
// their own.
func (f *filler) scanDeep(n *ooxml.Node, loc tplLoc) error {
	at := loc.name
	switch n.Name {
	case "txbxContent":
		f.boxes++
		return f.scanStory(n, "text box "+strconv.Itoa(f.boxes), loc.at)
	case "drawing":
		for _, box := range n.Kids {
			descr, ok := box.Child("docPr").Attr("descr")
			if !ok || !tag.Has(descr) {
				continue
			}
			toks, err := tag.Parse(descr)
			if err != nil {
				return tplErrorf(at, "the alternative text of a picture: %v", err)
			}
			pt := picTag{at: at}
			found := 0
			for _, t := range toks {
				switch t.Kind {
				case tag.Text:
					pt.alt += t.Text
				case tag.Image:
					pt.tok = t
					found++
				default:
					found = 2
				}
			}
			if found != 1 {
				return tplErrorf(at, "the alternative text of a picture takes one picture tag ({{%%name}}) and no other tag")
			}
			pt.alt = strings.TrimSpace(pt.alt)
			f.pics[n] = pt
			f.use(tag.Field{Name: pt.tok.Name, Kind: tag.FieldImage, Code: pt.tok.Code}, loc, nil)
			break
		}
	}
	for _, k := range n.Kids {
		if err := f.scanDeep(k, loc); err != nil {
			return err
		}
	}
	return nil
}

// seqText is the text of a run in a row of runs.
type seqText struct {
	run  *ooxml.Node
	text string
}

// scanInline reads the tags of the content of a paragraph, or of an
// element of a paragraph that holds runs (stack is then nil: the sections
// of its tags are closed in it, and base names the sections of the
// paragraph around it). The texts of the runs that follow each other are
// one text.
func (f *filler) scanInline(n *ooxml.Node, at tplLoc, stack *[]openSec, base []string) error {
	return f.scanInlineAtDepth(n, at, stack, base, 0)
}

func (f *filler) scanInlineAtDepth(n *ooxml.Node, at tplLoc, stack *[]openSec, base []string, outerDepth int) error {
	var items []tplItem
	tagged := false
	var seq []seqText
	var marks []tplItem
	flush := func() error {
		items = append(items, marks...)
		marks = marks[:0]
		if len(seq) == 0 {
			return nil
		}
		var text strings.Builder
		starts := make([]int, len(seq))
		for i, s := range seq {
			starts[i] = text.Len()
			text.WriteString(s.text)
		}
		src := text.String()
		runs := seq
		seq = nil
		if !tag.Has(src) {
			for _, s := range runs {
				items = append(items, tplItem{tok: tag.Token{Kind: tag.Text, Text: s.text}, run: s.run})
			}
			return nil
		}
		toks, err := tag.Parse(src)
		if err != nil {
			return tplErrorf(at.name, "%v", err)
		}
		runAt := func(pos int) *ooxml.Node { return runs[max(sort.SearchInts(starts, pos+1)-1, 0)].run }
		for _, t := range toks {
			if t.Kind != tag.Text || t.Text != src[t.Pos:t.End] {
				// a tag, or text with a brace written as text: in the
				// format of its first run
				tagged = true
				items = append(items, tplItem{tok: t, run: runAt(t.Pos)})
				continue
			}
			// text keeps the formats of its runs
			for pos := t.Pos; pos < t.End; {
				k := max(sort.SearchInts(starts, pos+1)-1, 0)
				end := t.End
				if k+1 < len(starts) && starts[k+1] < end {
					end = starts[k+1]
				}
				items = append(items, tplItem{tok: tag.Token{Kind: tag.Text, Text: src[pos:end], Pos: pos, End: end}, run: runs[k].run})
				pos = end
			}
		}
		return nil
	}
	for _, k := range n.Kids {
		switch k.Name {
		case "pPr", "proofErr":
		case "r":
			for _, rk := range k.Kids {
				switch rk.Name {
				case "rPr", "lastRenderedPageBreak":
				case "t":
					seq = append(seq, seqText{k, rk.Content()})
				default:
					if err := flush(); err != nil {
						return err
					}
					items = append(items, tplItem{run: k, kid: rk})
				}
			}
		case "bookmarkStart", "bookmarkEnd", "commentRangeStart", "commentRangeEnd", "permStart", "permEnd":
			// (does not part the text around it: it goes before that text)
			marks = append(marks, tplItem{node: k, mark: true})
		case "hyperlink", "smartTag", "fldSimple", "ins", "customXml", "dir", "bdo", "sdt":
			if err := flush(); err != nil {
				return err
			}
			in := k
			if k.Name == "sdt" {
				in = k.Child("sdtContent")
			}
			items = append(items, tplItem{node: k, wrap: in})
		default:
			if err := flush(); err != nil {
				return err
			}
			items = append(items, tplItem{node: k})
		}
	}
	if err := flush(); err != nil {
		return err
	}
	// what is in the kids that are not text: in the order of the document
	nested := func(it *tplItem, inner []string) error {
		switch {
		case it.mark:
		case it.wrap != nil:
			depth := outerDepth
			if stack != nil {
				depth = len(*stack)
			}
			return f.scanInlineAtDepth(it.wrap, at, nil, inner, depth)
		case it.node != nil:
			return f.scanDeep(it.node, at)
		case it.kid != nil:
			return f.scanDeep(it.kid, at)
		}
		return nil
	}
	if !tagged {
		for i := range items {
			if err := nested(&items[i], base); err != nil {
				return err
			}
		}
		return nil
	}
	tpl := &inlineTpl{items: items, at: at.name}
	f.inl[n] = tpl
	// the sections of the paragraph, and the tags of those beyond it
	// (localFrom: the uses of fields from those in a section on)
	var local, localUse, localFrom []int
	inner := append([]string(nil), base...)
	// outside removes section k of the paragraph from those around the uses
	// in it: it is a section beyond the paragraph
	outside := func(k int) {
		k += len(base)
		for i := localFrom[k-len(base)]; i < len(f.uses); i++ {
			u := &f.uses[i]
			u.inner = append(u.inner[:k:k], u.inner[k+1:]...)
		}
	}
	const beyond = "the tags of a section are both in or both outside a hyperlink, a field or a content control"
	for i := range items {
		it := &items[i]
		if it.node != nil || it.kid != nil {
			if err := nested(it, inner); err != nil {
				return err
			}
			continue
		}
		switch t := it.tok; t.Kind {
		case tag.Value:
			f.use(tag.Field{Name: t.Name, Kind: tag.FieldValue, Format: t.Format}, at, inner)
		case tag.Image:
			return tplErrorf(at.name, "the picture tag {{%%%s}} is in a text: in a Word document it is the alternative text of a picture, which the picture of the data replaces", t.Name)
		case tag.Open, tag.Inverted:
			depth := outerDepth + len(inner)
			if stack != nil {
				depth = len(*stack) + len(inner)
			}
			if depth >= maxSectionDepth {
				return tplErrorf(at.name, "the section {{#%s}} is inside more than %d sections", t.Name, maxSectionDepth)
			}
			localUse = append(localUse, f.use(tag.Field{Name: t.Name, Kind: tag.FieldSection, Inverted: t.Kind == tag.Inverted}, at, inner))
			localFrom = append(localFrom, len(f.uses))
			local = append(local, i)
			inner = append(inner, t.Name)
		case tag.Close:
			if n := len(local); n > 0 {
				oi, ou := local[n-1], localUse[n-1]
				o := items[oi].tok
				if t.Name != "" && t.Name != o.Name {
					return tplErrorf(at.name, "{{/%s}} closes the section {{#%s}}", t.Name, o.Name)
				}
				lift := o.Mode != "" && stack != nil
				if lift {
					outside(n - 1)
				}
				local, localUse, localFrom, inner = local[:n-1], localUse[:n-1], localFrom[:n-1], inner[:len(inner)-1]
				if o.Mode != "" && stack == nil {
					return tplErrorf(at.name, "the section {{#%s}} takes no %q in a hyperlink, a field or a content control", o.Name, o.Mode)
				}
				if !lift {
					continue
				}
				// a section that says how it repeats: it repeats more than
				// the text between its tags
				items[oi].outer, it.outer = true, true
				if err := f.section(openSec{tok: o, loc: at, use: ou}, at); err != nil {
					return err
				}
				continue
			}
			if stack == nil {
				return tplErrorf(at.name, "{{/%s}}: %s", t.Name, beyond)
			}
			n := len(*stack)
			if n == 0 {
				return tplErrorf(at.name, "{{/%s}} closes no section", t.Name)
			}
			o := (*stack)[n-1]
			if t.Name != "" && t.Name != o.tok.Name {
				return tplErrorf(at.name, "{{/%s}} closes the section {{#%s}} of %s", t.Name, o.tok.Name, o.loc.name)
			}
			*stack = (*stack)[:n-1]
			it.outer = true
			if err := f.section(o, at); err != nil {
				return err
			}
		}
	}
	if len(local) > 0 && stack == nil {
		return tplErrorf(at.name, "{{#%s}}: %s", items[local[0]].tok.Name, beyond)
	}
	for k := len(local) - 1; k >= 0; k-- {
		outside(k)
	}
	for k, i := range local {
		if len(*stack) >= maxSectionDepth {
			return tplErrorf(at.name, "the section {{#%s}} is inside more than %d sections", items[i].tok.Name, maxSectionDepth)
		}
		items[i].outer = true
		*stack = append(*stack, openSec{tok: items[i].tok, loc: at, use: localUse[k]})
	}
	tpl.only = stack != nil
	for _, it := range items {
		switch {
		case it.mark || it.outer:
		case it.node != nil || it.kid != nil || it.tok.Kind != tag.Text || strings.TrimSpace(it.tok.Text) != "":
			tpl.only = false
		}
	}
	return nil
}

// section notes a section that repeats more than text of a paragraph, from
// its opening tag to its closing tag: the kids of the innermost element
// that holds both, from the one with the first tag to the one with the
// other (the row, for the cells of a row), or what its tag says.
func (f *filler) section(o openSec, end tplLoc) error {
	t := o.tok
	s := &tplSec{name: t.Name, path: t.Path, inverted: t.Kind == tag.Inverted, mode: t.Mode, at: o.loc.name}
	a, b := o.loc.at, end.at
	for a.depth > b.depth {
		a = a.parent
	}
	for b.depth > a.depth {
		b = b.parent
	}
	for a.cont != b.cont {
		a, b = a.parent, b.parent
	}
	if a.cont.Name == "tr" {
		// in cells of a row: the row
		a, b = a.parent, b.parent
	}
	inTable := false
	for p := a; p != nil; p = p.parent {
		inTable = inTable || p.cont.Name == "tbl"
	}
	switch t.Mode {
	case tag.ModeRows:
		// the rows of the table it is in
		for a.cont.Name != "tbl" && inTable {
			a, b = a.parent, b.parent
		}
	case tag.ModeSlots:
		if a.cont.Name != "tbl" || a.idx == b.idx {
			return tplErrorf(s.at, "the section {{#%s}} takes no %q: slots are the rows of a table from that of {{#%s}} to that of {{/%s}}", t.Name, t.Mode, t.Name, t.Name)
		}
	case tag.ModePages:
		if inTable {
			return tplErrorf(s.at, "the section {{#%s}} is in a table: it takes no %q", t.Name, t.Mode)
		}
		if s.inverted {
			return tplErrorf(s.at, "the inverted section {{^%s}} does not repeat: it takes no %q", t.Name, t.Mode)
		}
	}
	s.a, s.b = a.idx, b.idx
	switch {
	case a.cont.Name != "tbl":
		if s.mode != tag.ModePages {
			s.mode = tag.ModeRows
		}
	case s.inverted || t.Mode == "" && s.a == s.b:
		s.mode = tag.ModeRows
	case t.Mode == "":
		s.mode = tag.ModeSlots
	}
	h := f.holders[a.cont]
	if h == nil {
		h = &holder{}
		f.holders[a.cont] = h
	}
	h.add(s)
	if o.use >= 0 {
		f.uses[o.use].sec = s
	}
	return nil
}

// check refuses the sections that cannot be filled: two that repeat the
// same thing, and sections of rows in slots.
func (f *filler) check(secs []*tplSec, parent *tplSec) error {
	for i, s := range secs {
		if parent != nil && parent.mode == tag.ModeSlots {
			return tplErrorf(s.at, "the section {{#%s}} is in the rows of {{#%s}}, which are slots (a row for an item); write {{#%s | rows}} to repeat the rows instead", s.name, parent.name, parent.name)
		}
		if i > 0 && s.a <= secs[i-1].b {
			return tplErrorf(s.at, "the sections {{#%s}} and {{#%s}} both repeat the same paragraph or row; one section repeats it, or sections one inside the other", secs[i-1].name, s.name)
		}
		if err := f.check(s.children, s); err != nil {
			return err
		}
	}
	return nil
}

// inOf names the sections around a place, outermost first: those of the
// elements that hold it, then those of its paragraph. For a section
// (stop), those around the section.
func (f *filler) inOf(at tplLoc, inner []string, stop *tplSec) []string {
	var levels [][]string
	stopped := -1
	for s := at.at; s != nil; s = s.parent {
		var names []string
		if h := f.holders[s.cont]; h != nil {
			secs := h.secs
		walk:
			for {
				for _, sec := range secs {
					if sec == stop {
						stopped = len(levels)
						break walk
					}
					if sec.a <= s.idx && s.idx <= sec.b {
						names = append(names, sec.name)
						secs = sec.children
						continue walk
					}
				}
				break
			}
		}
		levels = append(levels, names)
	}
	var out []string
	for i := len(levels) - 1; i >= 0 && i >= stopped; i-- {
		out = append(out, levels[i]...)
	}
	if stopped >= 0 {
		return out
	}
	return append(out, inner...)
}

// use notes a tag that takes a name of the data, and returns its place
// among the uses (-1 for a variable, which the data does not have).
func (f *filler) use(fd tag.Field, at tplLoc, inner []string) int {
	// (the variables of sections are not of the data; "." is kept: the
	// item of the section around it, which tells that its items are values)
	if strings.HasPrefix(fd.Name, "_") {
		return -1
	}
	f.uses = append(f.uses, fieldUse{fd: fd, at: at, inner: append([]string(nil), inner...)})
	return len(f.uses) - 1
}

// fieldList returns the fields of the template: one for each name in each
// set of sections, in the order of the document.
func (f *filler) fieldList() []tag.Field {
	var out []tag.Field
	index := map[string]int{}
	add := func(fd tag.Field, at string) {
		key := fd.Kind + "\x00" + strings.Join(fd.In, "\x00") + "\x00\x00" + fd.Name
		if i, ok := index[key]; ok {
			o := &out[i]
			if o.Format == "" {
				o.Format = fd.Format
			}
			if o.Mode == "" {
				o.Mode, o.Slots = fd.Mode, fd.Slots
			}
			o.Inverted = o.Inverted && fd.Inverted
			if n := len(o.At); n < 16 && o.At[n-1] != at {
				o.At = append(o.At, at)
			}
			return
		}
		index[key] = len(out)
		fd.At = []string{at}
		out = append(out, fd)
	}
	for _, u := range f.uses {
		fd := u.fd
		fd.In = f.inOf(u.at, u.inner, u.sec)
		if fd.Name == "." && len(fd.In) == 0 {
			continue // the data itself: not a name
		}
		if s := u.sec; s != nil {
			if fd.Mode = s.mode; s.mode == tag.ModeSlots {
				fd.Slots = s.b - s.a + 1
			}
		}
		add(fd, u.at.name)
	}
	return out
}

// lookup finds the value of a tag; a name the data does not have is told.
func (f *filler) lookup(t tag.Token, sc *tag.Scope, at string) (any, error) {
	v, ok := sc.Lookup(t.Path)
	if ok || sc.IsBlank() {
		return v, nil
	}
	if f.strict {
		return nil, tplErrorf(at, "the data has no %q", t.Name)
	}
	f.c.warnOnce("field:"+t.Name, "%s: the data has no %q; it is left empty", at, t.Name)
	return nil, nil
}

// clone copies an element with its tags replaced; nil for one that is left
// out. What has no tags in it is copied as it is.
func (f *filler) clone(n *ooxml.Node, sc *tag.Scope) (*ooxml.Node, error) {
	if f.nodes++; f.nodes > maxFilledNodes {
		return nil, fmt.Errorf("docx: the filled template has more than %d elements", maxFilledNodes)
	}
	if len(n.Kids) == 0 {
		return n, nil
	}
	if n.Name == "drawing" {
		if pt, ok := f.pics[n]; ok {
			return f.fillPicture(n, pt, sc)
		}
	}
	cp := *n
	if tpl := f.inl[n]; tpl != nil {
		if tpl.only {
			if n.Path("pPr", "sectPr") == nil {
				return nil, nil
			}
			// the last paragraph of a section of the document stays
			cp.Kids = []*ooxml.Node{n.Child("pPr")}
			return &cp, nil
		}
		kids, err := f.fillInline(n, tpl, sc)
		if err != nil {
			return nil, err
		}
		cp.Kids = kids
		return &cp, nil
	}
	cp.Kids = make([]*ooxml.Node, 0, len(n.Kids))
	var secs []*tplSec
	if h := f.holders[n]; h != nil {
		secs = h.secs
	}
	if err := f.fillKids(n, secs, 0, len(n.Kids)-1, sc, &cp.Kids); err != nil {
		return nil, err
	}
	if n.Name == "tc" || n.Name == "txbxContent" {
		// a cell ends with a paragraph, and a text box has one
		if k := len(cp.Kids); k == 0 || cp.Kids[k-1].Name != "p" {
			cp.Kids = append(cp.Kids, emptyPara(n))
		}
	}
	return &cp, nil
}

// emptyPara makes an empty paragraph in the format of the last paragraph
// of an element.
func emptyPara(cont *ooxml.Node) *ooxml.Node {
	space := cont.Space
	var pPr *ooxml.Node
	for _, k := range cont.Kids {
		if k.Name == "p" {
			space, pPr = k.Space, k.Child("pPr")
		}
	}
	p := &ooxml.Node{Space: space, Name: "p"}
	if pPr != nil && pPr.Child("sectPr") == nil {
		p.Kids = []*ooxml.Node{pPr}
	}
	return p
}

// fillKids fills the kids a to b of an element: those of the sections are
// repeated for their items, or filled as slots.
func (f *filler) fillKids(n *ooxml.Node, secs []*tplSec, a, b int, sc *tag.Scope, out *[]*ooxml.Node) error {
	i := a
	one := func(k *ooxml.Node, sc *tag.Scope) error {
		c, err := f.clone(k, sc)
		if c != nil {
			*out = append(*out, c)
		}
		return err
	}
	for _, s := range secs {
		for ; i < s.a; i++ {
			if err := one(n.Kids[i], sc); err != nil {
				return err
			}
		}
		i = s.b + 1
		v, err := f.lookup(tag.Token{Name: s.name, Path: s.path}, sc, s.at)
		if err != nil {
			return err
		}
		if s.mode == tag.ModeSlots {
			items, slots := tag.Items(v), s.b-s.a+1
			if len(items) > slots {
				return tplErrorf(s.at, "%d items for the %d rows of {{#%s}}: the rows from this one to that of {{/%s}} are slots, a row for an item (write {{#%s | rows}} to repeat the rows instead)",
					len(items), slots, s.name, s.name, s.name)
			}
			for k := 0; k < slots; k++ {
				if err := f.expansionStep(s.at); err != nil {
					return err
				}
				isc := sc.Blank(map[string]any{"_blank": true, "_blank_first": k == len(items), "_index": float64(k + 1), "_count": float64(len(items))})
				if k < len(items) {
					v := tag.Vars(k, len(items), false)
					v["_blank"], v["_blank_first"] = false, false
					isc = sc.Item(items[k], v)
				}
				tr, err := f.slotRow(n.Kids[s.a], n.Kids[s.a+k], k, isc)
				if err != nil {
					return err
				}
				*out = append(*out, tr)
			}
			continue
		}
		var items []any
		switch {
		case s.inverted && !tag.Truthy(v):
			items = []any{nil}
		case !s.inverted:
			items = tag.Items(v)
		}
		for k, it := range items {
			if err := f.expansionStep(s.at); err != nil {
				return err
			}
			isc := sc
			if !s.inverted {
				isc = sc.Item(it, tag.Vars(k, len(items), s.mode == tag.ModePages))
			}
			start := len(*out)
			if err := f.fillKids(n, s.children, s.a, s.b, isc, out); err != nil {
				return err
			}
			if s.mode == tag.ModePages && k > 0 && len(*out) > start {
				// each item on a new page
				rest := append([]*ooxml.Node(nil), (*out)[start+1:]...)
				*out = append(append((*out)[:start], breakBefore((*out)[start])...), rest...)
			}
		}
	}
	for ; i <= b; i++ {
		if err := one(n.Kids[i], sc); err != nil {
			return err
		}
	}
	return nil
}

// breakBefore puts a paragraph or a table on a new page: the paragraph
// with a page break before it, the table after a paragraph that has one
// (and is a point high).
func breakBefore(n *ooxml.Node) []*ooxml.Node {
	el := func(name string, attrs ...string) *ooxml.Node {
		e := &ooxml.Node{Space: n.Space, Name: name}
		for i := 0; i+1 < len(attrs); i += 2 {
			e.Attrs = append(e.Attrs, ooxml.Attr{Name: ooxml.Name{Space: n.Space, Local: attrs[i]}, Value: attrs[i+1]})
		}
		return e
	}
	if n.Name != "p" {
		pPr := el("pPr")
		rPr := el("rPr")
		rPr.Kids = []*ooxml.Node{el("sz", "val", "2")}
		pPr.Kids = []*ooxml.Node{el("pageBreakBefore"), el("spacing", "before", "0", "after", "0", "line", "20", "lineRule", "exact"), rPr}
		p := el("p")
		p.Kids = []*ooxml.Node{pPr}
		return []*ooxml.Node{p, n}
	}
	cp := *n
	pPr := el("pPr")
	if old := n.Child("pPr"); old != nil {
		pPr.Attrs = old.Attrs
		for _, k := range old.Kids {
			if k.Name != "pageBreakBefore" {
				pPr.Kids = append(pPr.Kids, k)
			}
		}
	}
	pPr.Kids = append([]*ooxml.Node{el("pageBreakBefore")}, pPr.Kids...)
	cp.Kids = []*ooxml.Node{pPr}
	for _, k := range n.Kids {
		if k.Name != "pPr" {
			cp.Kids = append(cp.Kids, k)
		}
	}
	return []*ooxml.Node{&cp}
}

// ownTags reports whether an element has tags of its own: not only those
// that close the section of its slots.
func (f *filler) ownTags(n *ooxml.Node) bool {
	if tpl := f.inl[n]; tpl != nil && !tpl.only {
		return true
	}
	if _, ok := f.pics[n]; ok {
		return true
	}
	for _, k := range n.Kids {
		if f.ownTags(k) {
			return true
		}
	}
	return false
}

// content reports whether an element has text or drawings that are not
// tags of sections beyond their paragraphs.
func (f *filler) content(n *ooxml.Node) bool {
	if tpl := f.inl[n]; tpl != nil && tpl.only {
		return false
	}
	switch n.Name {
	case "t":
		return n.Text != ""
	case "drawing", "pict", "object":
		return true
	}
	for _, k := range n.Kids {
		if f.content(k) {
			return true
		}
	}
	return false
}

// slotRow fills slot k of a section of slots: the row at that place, with
// the content of the cells of the first row in its cells that have none
// when the row has no tags of its own.
func (f *filler) slotRow(pattern, tr *ooxml.Node, k int, sc *tag.Scope) (*ooxml.Node, error) {
	if k == 0 || f.ownTags(tr) {
		return f.clone(tr, sc)
	}
	cp := *tr
	cp.Kids = nil
	first := pattern.Children("tc")
	c := 0
	for _, kid := range tr.Kids {
		src := kid
		if kid.Name == "tc" {
			if c < len(first) && f.ownTags(first[c]) && !f.content(kid) {
				// the content of the first row in the format of this cell
				// (its kids at the places of those of the first row: the
				// sections of that cell are those of this one)
				merged := *kid
				merged.Kids = make([]*ooxml.Node, len(first[c].Kids))
				for i, x := range first[c].Kids {
					if merged.Kids[i] = x; x.Name == "tcPr" {
						if own := kid.Child("tcPr"); own != nil {
							merged.Kids[i] = own
						}
					}
				}
				if h := f.holders[first[c]]; h != nil {
					f.holders[&merged] = h
					defer delete(f.holders, &merged)
				}
				src = &merged
			}
			c++
		}
		n, err := f.clone(src, sc)
		if err != nil {
			return nil, err
		}
		if n != nil {
			cp.Kids = append(cp.Kids, n)
		}
	}
	return &cp, nil
}

// piece is a piece of a filled paragraph: text in the format of a run, a
// child of a run that is not text, or another child of the paragraph.
type piece struct {
	text []string
	run  *ooxml.Node
	kid  *ooxml.Node
	node *ooxml.Node
}

// fillInline fills the content of a paragraph (or of an element of a
// paragraph that holds runs) and returns its new kids.
func (f *filler) fillInline(n *ooxml.Node, tpl *inlineTpl, sc *tag.Scope) ([]*ooxml.Node, error) {
	items, at := tpl.items, tpl.at
	var pieces []piece
	total := 0
	emit := func(s string, run *ooxml.Node) {
		if s == "" || total >= maxParaText {
			return
		}
		text, n := textPrefix(s, maxParaText-total)
		total += n
		if len(text) < len(s) {
			f.c.warnOnce("paratext:"+at, "%s: the text of a paragraph is cut at %d characters", at, maxParaText)
		}
		if k := len(pieces); k > 0 && pieces[k-1].node == nil && pieces[k-1].kid == nil && pieces[k-1].run == run {
			pieces[k-1].text = append(pieces[k-1].text, text)
			return
		}
		pieces = append(pieces, piece{text: []string{text}, run: run})
	}
	var walk func(lo, hi int, sc *tag.Scope) error
	walk = func(lo, hi int, sc *tag.Scope) error {
		for i := lo; i < hi; i++ {
			if err := f.expansionStep(at); err != nil {
				return err
			}
			it := items[i]
			switch {
			case it.node != nil:
				c, err := f.clone(it.node, sc)
				if err != nil {
					return err
				}
				if c != nil {
					pieces = append(pieces, piece{node: c})
				}
				continue
			case it.kid != nil:
				c, err := f.clone(it.kid, sc)
				if err != nil {
					return err
				}
				if c != nil {
					pieces = append(pieces, piece{run: it.run, kid: c})
				}
				continue
			case it.outer:
				continue
			}
			switch t := it.tok; t.Kind {
			case tag.Text:
				emit(t.Text, it.run)
			case tag.Value:
				v, err := f.lookup(t, sc, at)
				if err != nil {
					return err
				}
				switch v.(type) {
				case []any, map[string]any:
					f.c.warnOnce("fieldkind:"+t.Name, "%s: %q is a list or an object, not a value; it is left empty", at, t.Name)
				}
				if t.HasFormat && t.Format != "" && f.format == nil {
					f.c.warnOnce("format", "%s: the format of {{%s}} is not applied: the template is filled without a formatter of number formats", at, t.Name)
				}
				emit(tag.ValueText(v, t.Format, f.format), it.run)
			case tag.Open, tag.Inverted:
				j, depth := i+1, 1
				for ; j < hi; j++ {
					if items[j].node != nil || items[j].kid != nil || items[j].outer {
						continue
					}
					if k := items[j].tok.Kind; k == tag.Open || k == tag.Inverted {
						depth++
					} else if k == tag.Close {
						if depth--; depth == 0 {
							break
						}
					}
				}
				v, err := f.lookup(t, sc, at)
				if err != nil {
					return err
				}
				if t.Kind == tag.Inverted {
					if !tag.Truthy(v) {
						if err := walk(i+1, j, sc); err != nil {
							return err
						}
					}
				} else {
					list := tag.Items(v)
					for k, e := range list {
						if total >= maxParaText || f.nodes > maxFilledNodes {
							break
						}
						if err := f.expansionStep(at); err != nil {
							return err
						}
						if err := walk(i+1, j, sc.Item(e, tag.Vars(k, len(list), false))); err != nil {
							return err
						}
					}
				}
				i = j
			}
		}
		return nil
	}
	if err := walk(0, len(items), sc); err != nil {
		return nil, err
	}
	var out []*ooxml.Node
	if pPr := n.Child("pPr"); pPr != nil {
		out = append(out, pPr)
	}
	// the pieces of a run that follow each other are a run
	var cur, curSrc *ooxml.Node
	runOf := func(src *ooxml.Node) *ooxml.Node {
		if cur != nil && curSrc == src {
			return cur
		}
		cur, curSrc = &ooxml.Node{Space: src.Space, Name: "r", Attrs: src.Attrs}, src
		if rPr := src.Child("rPr"); rPr != nil {
			cur.Kids = append(cur.Kids, rPr)
		}
		out = append(out, cur)
		return cur
	}
	for _, pc := range pieces {
		switch {
		case pc.node != nil:
			cur = nil
			out = append(out, pc.node)
		case pc.kid != nil:
			r := runOf(pc.run)
			r.Kids = append(r.Kids, pc.kid)
		default:
			r := runOf(pc.run)
			// the line breaks and the tabs of a value are those of the line
			text := strings.ReplaceAll(strings.Join(pc.text, ""), "\r\n", "\n")
			for len(text) > 0 {
				i := strings.IndexAny(text, "\n\t")
				if i < 0 {
					i = len(text)
				}
				if i > 0 {
					r.Kids = append(r.Kids, &ooxml.Node{Space: pc.run.Space, Name: "t", Text: text[:i]})
				}
				if i == len(text) {
					break
				}
				name := "br"
				if text[i] == '\t' {
					name = "tab"
				}
				r.Kids = append(r.Kids, &ooxml.Node{Space: pc.run.Space, Name: name})
				text = text[i+1:]
			}
		}
	}
	f.nodes += len(out)
	return out, nil
}

// textPrefix reads at most n+1 runes without copying an oversized value.
func textPrefix(s string, n int) (string, int) {
	count := 0
	for i := range s {
		if count == n {
			return s[:i], count
		}
		count++
	}
	return s, count
}

// copyTree copies an element with what is in it.
func copyTree(n *ooxml.Node) *ooxml.Node {
	cp := *n
	cp.Attrs = append([]ooxml.Attr(nil), n.Attrs...)
	cp.Kids = make([]*ooxml.Node, len(n.Kids))
	for i, k := range n.Kids {
		cp.Kids[i] = copyTree(k)
	}
	return &cp
}

// setAttr sets an attribute of an element by its local name.
func setAttr(n *ooxml.Node, name, value string) {
	for i := range n.Attrs {
		if n.Attrs[i].Name.Local == name && !strings.HasSuffix(n.Attrs[i].Name.Space, "relationships") {
			n.Attrs[i].Value = value
			return
		}
	}
	n.Attrs = append(n.Attrs, ooxml.Attr{Name: ooxml.Name{Local: name}, Value: value})
}

// fillPicture replaces the picture of a drawing by the picture that its
// tag names, fitted in its frame; nil when the data names none.
func (f *filler) fillPicture(drawing *ooxml.Node, pt picTag, sc *tag.Scope) (*ooxml.Node, error) {
	t, at := pt.tok, pt.at
	v, err := f.lookup(t, sc, at)
	if err != nil {
		return nil, err
	}
	fail := func(msg string) (*ooxml.Node, error) {
		if f.strict {
			return nil, tplErrorf(at, "the picture of %q %s", t.Name, msg)
		}
		f.c.warnf("%s: the picture of %q %s; it is left out", at, t.Name, msg)
		return nil, nil
	}
	var rid string
	var w, h float64
	alt, fit := pt.alt, t.Fit
	if t.Code != "" {
		// a QR code or a barcode of the value, drawn in the place of the picture
		if v == nil || v == "" {
			return nil, nil
		}
		sym, err := t.Symbol(v)
		if err != nil {
			if f.strict {
				return nil, tplErrorf(at, "%v", err)
			}
			f.c.warnf("%s: the code of %q is left out: %v", at, t.Name, err)
			return nil, nil
		}
		rid = f.c.r.InjectDrawing(func(o *bdf.Object, x, y, w, h float32) {
			o.FillColor(bdf.RGB(0, 0, 0)).FillPath(o.AddPath(sym.Path(x, y, w, h)), bdf.NonZero)
		})
		w, h = float64(sym.W), float64(sym.H)
		if alt == "" {
			alt = sym.Text // what the symbol says
		}
		if fit = ""; sym.Linear {
			fit = tag.FitStretch // the bars are as wide and as tall as their place
		}
	} else {
		name, _ := v.(string)
		if name == "" {
			return nil, nil
		}
		data, err := tag.PictureData(f.files, name)
		if err != nil {
			return fail("is not read: " + err.Error())
		}
		var ok bool
		if rid, w, h, ok = f.c.r.InjectImage(data); !ok {
			return fail("is not a PNG, JPEG, GIF, BMP or WebP image")
		}
	}
	cp := copyTree(drawing)
	f.nodes += 32
	for _, box := range cp.Kids {
		docPr := box.Child("docPr")
		if d, ok := docPr.Attr("descr"); !ok || !tag.Has(d) {
			continue
		}
		setAttr(docPr, "descr", alt)
		pic := box.Path("graphic", "graphicData", "pic")
		bf := pic.Child("blipFill")
		blip := bf.Child("blip")
		set := false
		if blip != nil {
			for i := range blip.Attrs {
				if a := &blip.Attrs[i]; a.Name.Local == "embed" || a.Name.Local == "link" {
					a.Name.Local, a.Value, set = "embed", rid, true
				}
			}
		}
		if !set {
			return fail("has no place in the template's picture")
		}
		blip.Kids = nil // the effects of the picture it replaces are not its
		// without the crop of the picture it replaces
		kids := bf.Kids[:0]
		for _, k := range bf.Kids {
			if k.Name != "srcRect" {
				kids = append(kids, k)
			}
		}
		bf.Kids = kids
		ext := box.Child("extent")
		bw, bh := float64(ext.AttrInt("cx", 0)), float64(ext.AttrInt("cy", 0))
		if ext == nil || bw <= 0 || bh <= 0 {
			break
		}
		emu := func(v float64) string { return strconv.FormatInt(int64(math.Round(v)), 10) }
		switch fit {
		case tag.FitStretch:
		case tag.FitCover:
			// the middle of the picture, cut to the frame
			crop := &ooxml.Node{Space: blip.Space, Name: "srcRect"}
			pct := func(v float64) string { return strconv.Itoa(int(math.Round(v * 100000))) }
			if w/h > bw/bh {
				cut := (1 - (bw/bh)/(w/h)) / 2
				setAttr(crop, "l", pct(cut))
				setAttr(crop, "r", pct(cut))
			} else {
				cut := (1 - (w/h)/(bw/bh)) / 2
				setAttr(crop, "t", pct(cut))
				setAttr(crop, "b", pct(cut))
			}
			at := 0
			for i, k := range bf.Kids {
				if k.Name == "blip" {
					at = i + 1
				}
			}
			bf.Kids = append(bf.Kids[:at], append([]*ooxml.Node{crop}, bf.Kids[at:]...)...)
		default:
			// the whole picture: a picture in a line becomes smaller, one
			// placed on the page stays in the middle of its frame
			k := math.Min(bw/w, bh/h)
			nw, nh := w*k, h*k
			setAttr(ext, "cx", emu(nw))
			setAttr(ext, "cy", emu(nh))
			if e := pic.Path("spPr", "xfrm", "ext"); e != nil {
				setAttr(e, "cx", emu(nw))
				setAttr(e, "cy", emu(nh))
			}
			if box.Name == "anchor" {
				for _, mv := range []struct {
					pos string
					by  float64
				}{{"positionH", (bw - nw) / 2}, {"positionV", (bh - nh) / 2}} {
					if off := box.Path(mv.pos, "posOffset"); off != nil {
						if v, err := strconv.ParseFloat(strings.TrimSpace(off.Text), 64); err == nil {
							off.Text = emu(v + mv.by)
						}
					}
				}
			}
		}
		break
	}
	return cp, nil
}
