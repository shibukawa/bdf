package pptx

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
	"github.com/shibukawa/bdf/form/tag"
)

// A presentation is a template when the text of its slides holds tags
// (package form/tag): Fill converts it with the tags replaced by data. The
// slides are filled as their trees of elements in memory, before they are
// laid out; no presentation is written.
//
// A tag in a text is replaced by its value, in the format of the run the
// tag starts in (PowerPoint splits a tag into runs where the language or
// the spelling check changes: the runs of a paragraph are read as one
// text). A section repeats the smallest of these that holds both its tags:
//
//   - a part of the text of a paragraph (both tags in the paragraph);
//   - paragraphs of a text (the tags in different paragraphs of a shape or
//     of a cell): those from the one to the other;
//   - rows of a table: the row, for tags in the same row (the table grows),
//     and for tags in different rows the rows from the one to the other,
//     which are slots, a row for an item (tag.ModeSlots; "| rows" repeats
//     them instead);
//   - slides: the slide, for tags in different shapes of it, and the slides
//     from the one to the other for tags on different slides.
//
// How a section repeats, when its tag says it, is what it repeats:
// "| pages" the slide, and "| rows" the row of a table or, in another text,
// the paragraph, also when both tags are in the same paragraph:
// "{{#items | rows}}{{name}}{{/items}}" in a paragraph of a list is a
// paragraph for each item.
//
// A paragraph that holds nothing but the tags of such sections is left out.
// A picture whose alternative text has a picture tag ({{%photo}}) is
// replaced by the picture the data names, fitted in its frame; the rest of
// the alternative text is that of the new picture.
//
// Text that has changed is fitted again as PowerPoint fits it: a shape that
// shrinks its text on overflow (a:normAutofit) gets the scale at which the
// new text fits, and a shape that resizes to fit its text (a:spAutoFit) its
// new height. Shapes do not push each other.
//
// {{_slide}} and {{_slides}} are the number of the slide and that of the
// slides of the filled presentation. A table whose rows were
// repeated and that reaches past the foot of its slide goes on on copies of
// the slide, with the rows before the repeated ones at the top of each.

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

// The variables of the slides: {{_slide}} is the number of the slide a text
// is on and {{_slides}} that of the slides, as they are when the template
// has been filled (with the slides that sections repeated and those that
// tables go on on).
const (
	varSlide  = "_slide"
	varSlides = "_slides"
	// fieldSlides is the type of the field (a:fld) that shows the number
	// of slides.
	fieldSlides = "bdf-slides"
)

// Limits of a filled template.
const (
	// maxFilledSlides bounds the slides of a filled template.
	maxFilledSlides = 10000
	// maxFilledNodes bounds the elements of the slides of a filled template.
	maxFilledNodes = 8 << 20
	// maxParaText bounds the characters of a filled paragraph.
	maxParaText = 1 << 16
	// maxSectionDepth bounds the sections one inside the other.
	maxSectionDepth = 32
	// maxExpansionSteps bounds repetitions that produce no nodes or text.
	maxExpansionSteps = 8 << 20
)

// Fill converts a presentation as a template filled with data.
func Fill(r io.ReaderAt, size int64, t *TemplateData, opts *Options) (*Result, error) {
	o := Options{}
	if opts != nil {
		o = *opts
	}
	if t == nil {
		t = &TemplateData{}
	}
	f := newFiller(t)
	o.fill = f
	res, err := Convert(r, size, &o)
	if err != nil {
		return nil, err
	}
	res.Fields = f.fieldList()
	return res, nil
}

// Inspect reads the tags of a presentation without filling it: the fields
// that data for it has, or the first mistake in its tags.
func Inspect(r io.ReaderAt, size int64, opts *Options) ([]tag.Field, error) {
	o := Options{NoSystemFonts: true, SystemFonts: true}
	f := newFiller(&TemplateData{})
	f.inspect = true
	o.fill = f
	if _, err := Convert(r, size, &o); err != nil {
		return nil, err
	}
	return f.fieldList(), nil
}

// slideJob is a slide to draw: a part of the package, or a filled copy of
// one.
type slideJob struct {
	part string
	node *ooxml.Node // nil: the part as it is
	num  int         // the slide number that its fields show
}

// loc is where a paragraph is.
type loc struct {
	slide int
	tbl   *ooxml.Node // the table it is in, and the row of it
	row   int
	body  *ooxml.Node // its text body, and its place among the paragraphs
	para  int
	name  string // "slide 2: Title 1"
}

// Kinds of sections: what they repeat.
const (
	secParas = iota
	secRows
	secSlides
)

// sec is a section that repeats paragraphs, rows or slides: the units a to
// b of what holds them.
type sec struct {
	name     string
	path     []string
	inverted bool
	mode     string
	kind     int
	a, b     int
	at       string
	children []*sec
}

// holder is the sections of a text body, a table or the presentation.
type holder struct{ secs []*sec }

// add puts a section among those of the holder: the ones inside it become
// its children.
func (h *holder) add(s *sec) {
	var rest []*sec
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

// pItem is a piece of a paragraph: text or a tag, or a child element that
// is not a run.
type pItem struct {
	tok   tag.Token
	run   *ooxml.Node // the run whose format the text has
	node  *ooxml.Node // a child that is kept as it is (a:br, a:fld, a14:m)
	outer bool        // a tag of a section beyond the paragraph
}

// paraTpl is the pieces of a paragraph with tags.
type paraTpl struct {
	items []pItem
	// only is set when the paragraph holds nothing but tags of sections
	// beyond it: it is left out
	only bool
}

// open is a section whose closing tag has not been read.
type open struct {
	tok tag.Token
	loc loc
	use int // its use among those of the fields (-1: none)
}

// fieldUse is a tag that takes a name from the data.
type fieldUse struct {
	fd    tag.Field
	at    loc
	inner []string // the sections of its paragraph around it
	sec   *sec     // of a section that reaches beyond its paragraph
}

// picTag is the tag of a picture, and the rest of its alternative text.
type picTag struct {
	tok tag.Token
	alt string
}

type filler struct {
	c       *converter
	root    *tag.Scope
	strict  bool
	files   fs.FS
	format  tag.Formatter
	inspect bool

	parts  []string
	slides []*ooxml.Node
	paras  map[*ooxml.Node]*paraTpl
	bodies map[*ooxml.Node]*holder
	tables map[*ooxml.Node]*holder
	top    holder
	pics   map[*ooxml.Node]picTag

	uses  []fieldUse
	nodes int
	steps int // repeated items and inline tokens expanded
	// grown are the tables of the filled slides whose rows were repeated:
	// the rows before and after the repeated ones
	grown map[*ooxml.Node][2]int
	// loaded are the pictures that the data names, by their names, as
	// they were stored
	loaded map[string]injected
}

// injected is a picture of the data as it was stored in the document (see
// drawingml.Renderer.InjectImage), or why it is left out.
type injected struct {
	rid  string
	w, h float64
	ok   bool
	why  string
}

func newFiller(t *TemplateData) *filler {
	data := t.Data
	if data == nil {
		data = map[string]any{}
	}
	return &filler{root: tag.NewScope(tag.Normalize(data)), strict: t.Strict, files: t.Files, format: t.Format,
		paras: map[*ooxml.Node]*paraTpl{}, bodies: map[*ooxml.Node]*holder{}, tables: map[*ooxml.Node]*holder{},
		pics: map[*ooxml.Node]picTag{}, grown: map[*ooxml.Node][2]int{}, loaded: map[string]injected{}}
}

// run fills the slides of the parts and returns the slides to draw.
func (f *filler) run(c *converter, parts []string) ([]slideJob, error) {
	f.c, f.parts = c, parts
	var stack []open
	for i, part := range parts {
		n, err := c.pkg.XML(part)
		if err != nil {
			return nil, fmt.Errorf("pptx: slide %d: %w", i+1, err)
		}
		f.slides = append(f.slides, n)
		if err := f.scan(n.Path("cSld", "spTree"), loc{slide: i, name: "slide " + strconv.Itoa(i+1)}, &stack); err != nil {
			return nil, err
		}
	}
	if n := len(stack); n > 0 {
		o := stack[n-1]
		return nil, fmt.Errorf("pptx: %s: the section {{#%s}} is not closed with {{/%s}}", o.loc.name, o.tok.Name, o.tok.Name)
	}
	if err := f.check(f.top.secs, nil); err != nil {
		return nil, err
	}
	for _, h := range f.tables {
		if err := f.check(h.secs, nil); err != nil {
			return nil, err
		}
	}
	for _, h := range f.bodies {
		if err := f.check(h.secs, nil); err != nil {
			return nil, err
		}
	}
	if f.inspect {
		return nil, nil
	}
	var out []slideJob
	if err := f.fillSlides(f.top.secs, 0, len(parts)-1, f.root, &out); err != nil {
		return nil, err
	}
	// tables that reach past the foot of their slides go on on copies
	var paged []slideJob
	for _, j := range out {
		more, err := f.paginate(j)
		if err != nil {
			return nil, err
		}
		paged = append(paged, more...)
		if len(paged) > maxFilledSlides {
			return nil, fmt.Errorf("pptx: the filled template has more than %d slides", maxFilledSlides)
		}
	}
	for i := range paged {
		paged[i].num = c.firstSlide + i
	}
	return paged, nil
}

// errorf makes an error about a place of the template.
func errorAt(at, format string, args ...any) error {
	return fmt.Errorf("pptx: %s: %s", at, fmt.Sprintf(format, args...))
}

func (f *filler) expansionStep(at string) error {
	if f.steps >= maxExpansionSteps {
		return errorAt(at, "the filled template needs more than %d expansion steps", maxExpansionSteps)
	}
	f.steps++
	return nil
}

// shapeName returns the name of a shape of a slide, for messages.
func shapeName(n *ooxml.Node) string {
	for _, k := range n.Elements() {
		if strings.HasPrefix(k.Name, "nv") {
			if name := k.Child("cNvPr").AttrStr("name", ""); name != "" {
				return name
			}
		}
	}
	return ""
}

// scan reads the tags of the texts under an element and pairs the tags of
// the sections that reach beyond paragraphs.
func (f *filler) scan(n *ooxml.Node, at loc, stack *[]open) error {
	if n == nil {
		return nil
	}
	switch n.Name {
	case "sp", "graphicFrame", "pic", "cxnSp", "grpSp":
		if name := shapeName(n); name != "" && at.tbl == nil {
			at.name = "slide " + strconv.Itoa(at.slide+1) + ": " + name
		}
	}
	switch n.Name {
	case "tbl":
		at.tbl = n
		for i, tr := range n.Children("tr") {
			at.row = i
			for _, tc := range tr.Children("tc") {
				if err := f.scan(tc, at, stack); err != nil {
					return err
				}
			}
		}
		return nil
	case "txBody":
		at.body = n
		for i, p := range n.Children("p") {
			at.para = i
			if err := f.scanPara(p, at, stack); err != nil {
				return err
			}
		}
		return nil
	case "pic":
		if descr, ok := cNvPrOf(n).Attr("descr"); ok && tag.Has(descr) {
			toks, err := tag.Parse(descr)
			if err != nil {
				return errorAt(at.name, "the alternative text: %v", err)
			}
			var pt picTag
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
				return errorAt(at.name, "the alternative text of a picture takes one picture tag ({{%%name}}) and no other tag")
			}
			pt.alt = strings.TrimSpace(pt.alt)
			f.pics[n] = pt
			f.use(tag.Field{Name: pt.tok.Name, Kind: tag.FieldImage, Code: pt.tok.Code}, at, nil)
		}
	}
	for _, k := range n.Elements() {
		if err := f.scan(k, at, stack); err != nil {
			return err
		}
	}
	return nil
}

// cNvPrOf returns the non-visual properties of a shape (its name and its
// alternative text).
func cNvPrOf(n *ooxml.Node) *ooxml.Node {
	for _, k := range n.Elements() {
		if strings.HasPrefix(k.Name, "nv") {
			return k.Child("cNvPr")
		}
	}
	return nil
}

// scanPara reads the tags of a paragraph: the texts of its runs, one after
// the other, are one text.
func (f *filler) scanPara(p *ooxml.Node, at loc, stack *[]open) error {
	var items []pItem
	tagged := false
	kids := p.Kids
	for i := 0; i < len(kids); {
		k := kids[i]
		switch k.Name {
		case "pPr", "endParaRPr":
			i++
			continue
		case "r":
		default:
			items = append(items, pItem{node: k})
			i++
			continue
		}
		// the runs that follow each other
		var text strings.Builder
		var starts []int
		j := i
		for ; j < len(kids) && kids[j].Name == "r"; j++ {
			starts = append(starts, text.Len())
			text.WriteString(kids[j].Child("t").Content())
		}
		src := text.String()
		runAt := func(pos int) *ooxml.Node {
			k := sort.SearchInts(starts, pos+1) - 1
			return kids[i+max(k, 0)]
		}
		if !tag.Has(src) && !strings.Contains(src, `\`) {
			for k := i; k < j; k++ {
				items = append(items, pItem{tok: tag.Token{Kind: tag.Text, Text: kids[k].Child("t").Content()}, run: kids[k]})
			}
			i = j
			continue
		}
		toks, err := tag.Parse(src)
		if err != nil {
			return errorAt(at.name, "%v", err)
		}
		for _, t := range toks {
			if t.Kind != tag.Text {
				tagged = true
				items = append(items, pItem{tok: t, run: runAt(t.Pos)})
				continue
			}
			if t.Text != src[t.Pos:t.End] {
				// with a brace written as text: in the format of its first run
				tagged = true
				items = append(items, pItem{tok: t, run: runAt(t.Pos)})
				continue
			}
			// text keeps the formats of its runs
			for pos := t.Pos; pos < t.End; {
				k := sort.SearchInts(starts, pos+1) - 1
				end := t.End
				if k+1 < len(starts) && starts[k+1] < end {
					end = starts[k+1]
				}
				items = append(items, pItem{tok: tag.Token{Kind: tag.Text, Text: src[pos:end], Pos: pos, End: end}, run: kids[i+k]})
				pos = end
			}
		}
		i = j
	}
	if !tagged {
		return nil
	}
	tpl := &paraTpl{items: items}
	f.paras[p] = tpl
	// the sections of the paragraph, and the tags of those beyond it
	// (localFrom: the uses of fields from those in a section on)
	var local, localUse, localFrom []int
	var inner []string
	// outside removes section k of the paragraph from those around the uses
	// in it: it is a section beyond the paragraph
	outside := func(k int) {
		for i := localFrom[k]; i < len(f.uses); i++ {
			u := &f.uses[i]
			u.inner = append(u.inner[:k:k], u.inner[k+1:]...)
		}
	}
	for i := range items {
		it := &items[i]
		if it.node != nil {
			continue
		}
		switch t := it.tok; t.Kind {
		case tag.Value:
			f.use(tag.Field{Name: t.Name, Kind: tag.FieldValue, Format: t.Format}, at, inner)
		case tag.Image:
			return errorAt(at.name, "the picture tag {{%%%s}} is in a text: in a presentation it is the alternative text of a picture, which the picture of the data replaces", t.Name)
		case tag.Open, tag.Inverted:
			if len(*stack)+len(local) >= maxSectionDepth {
				return errorAt(at.name, "the section {{#%s}} is inside more than %d sections", t.Name, maxSectionDepth)
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
					return errorAt(at.name, "{{/%s}} closes the section {{#%s}}", t.Name, o.Name)
				}
				if o.Mode != "" {
					outside(n - 1)
				}
				local, localUse, localFrom, inner = local[:n-1], localUse[:n-1], localFrom[:n-1], inner[:n-1]
				if o.Mode == "" {
					continue
				}
				// a section that says how it repeats: it repeats more than
				// the text between its tags
				items[oi].outer, it.outer = true, true
				if err := f.section(open{tok: o, loc: at, use: ou}, at); err != nil {
					return err
				}
				continue
			}
			n := len(*stack)
			if n == 0 {
				return errorAt(at.name, "{{/%s}} closes no section", t.Name)
			}
			o := (*stack)[n-1]
			if t.Name != "" && t.Name != o.tok.Name {
				return errorAt(at.name, "{{/%s}} closes the section {{#%s}} of %s", t.Name, o.tok.Name, o.loc.name)
			}
			*stack = (*stack)[:n-1]
			it.outer = true
			if err := f.section(o, at); err != nil {
				return err
			}
		}
	}
	for k := len(local) - 1; k >= 0; k-- {
		outside(k)
	}
	for k, i := range local {
		if len(*stack) >= maxSectionDepth {
			return errorAt(at.name, "the section {{#%s}} is inside more than %d sections", items[i].tok.Name, maxSectionDepth)
		}
		items[i].outer = true
		*stack = append(*stack, open{tok: items[i].tok, loc: at, use: localUse[k]})
	}
	tpl.only = true
	for _, it := range items {
		if it.node != nil || !it.outer && (it.tok.Kind != tag.Text || strings.TrimSpace(it.tok.Text) != "") {
			tpl.only = false
		}
	}
	return nil
}

// section notes a section that repeats more than text of a paragraph, from
// its opening tag to its closing tag. What it repeats is the smallest of a
// text, a table and the slides that holds both tags, or what its tag says:
// "pages" is slides, and "rows" the rows of a table or, in another text,
// paragraphs.
func (f *filler) section(o open, end loc) error {
	t := o.tok
	s := &sec{name: t.Name, path: t.Path, inverted: t.Kind == tag.Inverted, mode: t.Mode, at: o.loc.name}
	kind := secSlides
	switch {
	case o.loc.slide != end.slide:
	case o.loc.body == end.body:
		kind = secParas
	case o.loc.tbl != nil && o.loc.tbl == end.tbl:
		kind = secRows
	}
	table := o.loc.slide == end.slide && o.loc.tbl != nil && o.loc.tbl == end.tbl
	switch t.Mode {
	case tag.ModePages:
		kind = secSlides
	case tag.ModeRows:
		if kind == secSlides {
			return errorAt(s.at, "the section {{#%s}} repeats slides (its tags are in different shapes): it takes no %q", t.Name, t.Mode)
		}
		if table {
			kind = secRows
		}
	case tag.ModeSlots:
		if !table || o.loc.row == end.row {
			return errorAt(s.at, "the section {{#%s}} takes no %q: slots are the rows of a table from that of {{#%s}} to that of {{/%s}}", t.Name, t.Mode, t.Name, t.Name)
		}
		kind = secRows
	}
	var h *holder
	switch s.kind = kind; kind {
	case secParas:
		s.a, s.b, s.mode = o.loc.para, end.para, tag.ModeRows
		if h = f.bodies[end.body]; h == nil {
			h = &holder{}
			f.bodies[end.body] = h
		}
	case secRows:
		s.a, s.b = o.loc.row, end.row
		if s.inverted || t.Mode == "" && s.a == s.b {
			s.mode = tag.ModeRows
		} else if t.Mode == "" {
			s.mode = tag.ModeSlots
		}
		if h = f.tables[end.tbl]; h == nil {
			h = &holder{}
			f.tables[end.tbl] = h
		}
	default:
		s.a, s.b, s.mode = o.loc.slide, end.slide, tag.ModePages
		h = &f.top
	}
	h.add(s)
	if o.use >= 0 {
		f.uses[o.use].sec = s
	}
	return nil
}

// check refuses the sections that cannot be filled: two that repeat the
// same unit, and sections of rows in slots.
func (f *filler) check(secs []*sec, parent *sec) error {
	for i, s := range secs {
		if parent != nil && parent.mode == tag.ModeSlots && s.kind == secRows {
			return errorAt(s.at, "the section {{#%s}} is in the rows of {{#%s}}, which are slots (a row for an item); write {{#%s | rows}} to repeat the rows instead", s.name, parent.name, parent.name)
		}
		if i > 0 && s.a <= secs[i-1].b {
			what := []string{"paragraph", "row", "slide"}[s.kind]
			return errorAt(s.at, "the sections {{#%s}} and {{#%s}} both repeat %s %d; one section repeats it, or sections one inside the other", secs[i-1].name, s.name, what, s.a+1)
		}
		if err := f.check(s.children, s); err != nil {
			return err
		}
	}
	return nil
}

// inOf names the sections around a place, outermost first: those of the
// slides, of its table and of its text that hold it, then those of its
// paragraph. For a section (stop), those around the section.
func (f *filler) inOf(at loc, inner []string, stop *sec) []string {
	var out []string
	done := false
	var walk func(secs []*sec, unit int)
	walk = func(secs []*sec, unit int) {
		for _, s := range secs {
			if s == stop {
				done = true
				return
			}
			if s.a <= unit && unit <= s.b {
				out = append(out, s.name)
				walk(s.children, unit)
				return
			}
		}
	}
	walk(f.top.secs, at.slide)
	if h := f.tables[at.tbl]; h != nil && at.tbl != nil && !done {
		walk(h.secs, at.row)
	}
	if h := f.bodies[at.body]; h != nil && at.body != nil && !done {
		walk(h.secs, at.para)
	}
	if done {
		return out
	}
	return append(out, inner...)
}

// use notes a tag that takes a name of the data, and returns its place
// among the uses (-1 for a variable, which the data does not have).
func (f *filler) use(fd tag.Field, at loc, inner []string) int {
	// (the variables of sections are not of the data; "." is kept: the
	// item of the section around it, which tells that its items are values)
	if strings.HasPrefix(fd.Name, "_") {
		return -1
	}
	f.uses = append(f.uses, fieldUse{fd: fd, at: at, inner: append([]string(nil), inner...)})
	return len(f.uses) - 1
}

// fieldList returns the fields of the template: one for each name in each
// set of sections, in the order of the slides.
func (f *filler) fieldList() []tag.Field {
	var out []tag.Field
	index := map[string]int{}
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
			if n := len(o.At); n < 16 && o.At[n-1] != u.at.name {
				o.At = append(o.At, u.at.name)
			}
			continue
		}
		index[key] = len(out)
		fd.At = []string{u.at.name}
		out = append(out, fd)
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
		return nil, errorAt(at, "the data has no %q", t.Name)
	}
	f.c.warnOnce("field:"+t.Name, "%s: the data has no %q; it is left empty", at, t.Name)
	return nil, nil
}

// fillSlides fills the slides s0 to s1 of the template: those of the
// sections are repeated for their items.
func (f *filler) fillSlides(secs []*sec, s0, s1 int, sc *tag.Scope, out *[]slideJob) error {
	i := s0
	for _, s := range secs {
		for ; i < s.a; i++ {
			if err := f.fillSlide(i, sc, out); err != nil {
				return err
			}
		}
		i = s.b + 1
		v, err := f.lookup(tag.Token{Name: s.name, Path: s.path}, sc, s.at)
		if err != nil {
			return err
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
				isc = sc.Item(it, tag.Vars(k, len(items), true))
			}
			if err := f.fillSlides(s.children, s.a, s.b, isc, out); err != nil {
				return err
			}
		}
	}
	for ; i <= s1; i++ {
		if err := f.fillSlide(i, sc, out); err != nil {
			return err
		}
	}
	return nil
}

func (f *filler) fillSlide(i int, sc *tag.Scope, out *[]slideJob) error {
	if len(*out) >= maxFilledSlides {
		return fmt.Errorf("pptx: the filled template has more than %d slides", maxFilledSlides)
	}
	n, err := f.clone(f.slides[i], sc, "slide "+strconv.Itoa(i+1))
	if err != nil {
		return err
	}
	*out = append(*out, slideJob{part: f.parts[i], node: n})
	return nil
}

// clone copies an element of a slide with its tags replaced; nil for one
// that is left out.
func (f *filler) clone(n *ooxml.Node, sc *tag.Scope, at string) (*ooxml.Node, error) {
	if f.nodes++; f.nodes > maxFilledNodes {
		return nil, fmt.Errorf("pptx: the filled template has more than %d elements", maxFilledNodes)
	}
	switch n.Name {
	case "sp", "graphicFrame", "pic", "cxnSp", "grpSp":
		if name := shapeName(n); name != "" {
			if head, _, ok := strings.Cut(at, ":"); ok {
				at = head
			}
			at += ": " + name
		}
	}
	switch n.Name {
	case "tbl":
		return f.fillTable(n, sc, at)
	case "txBody":
		return f.fillBody(n, sc, at)
	case "pic":
		if pt, ok := f.pics[n]; ok {
			return f.fillPicture(n, pt, sc, at)
		}
	}
	cp := *n
	cp.Kids = make([]*ooxml.Node, 0, len(n.Kids))
	for _, k := range n.Kids {
		c, err := f.clone(k, sc, at)
		if err != nil {
			return nil, err
		}
		if c != nil {
			cp.Kids = append(cp.Kids, c)
		}
	}
	return &cp, nil
}

// copyTree copies an element with what is in it, as it is.
func copyTree(n *ooxml.Node) *ooxml.Node {
	cp := *n
	cp.Attrs = append([]ooxml.Attr(nil), n.Attrs...)
	cp.Kids = make([]*ooxml.Node, len(n.Kids))
	for i, k := range n.Kids {
		cp.Kids[i] = copyTree(k)
	}
	return &cp
}

// setAttr sets an attribute of an element.
func setAttr(n *ooxml.Node, name, value string) {
	for i := range n.Attrs {
		if n.Attrs[i].Name.Local == name && n.Attrs[i].Name.Space == "" {
			n.Attrs[i].Value = value
			return
		}
	}
	n.Attrs = append(append([]ooxml.Attr(nil), n.Attrs...), ooxml.Attr{Name: ooxml.Name{Local: name}, Value: value})
}

// fillBody fills a text body: the paragraphs of its sections are repeated
// for their items. Its text is then fitted again when it is drawn.
func (f *filler) fillBody(body *ooxml.Node, sc *tag.Scope, at string) (*ooxml.Node, error) {
	cp := *body
	cp.Kids = nil
	var paras []*ooxml.Node
	for _, k := range body.Kids {
		if k.Name == "p" {
			paras = append(paras, k)
		}
	}
	var secs []*sec
	if h := f.bodies[body]; h != nil {
		secs = h.secs
	}
	changed := len(secs) > 0
	for _, p := range paras {
		changed = changed || f.paras[p] != nil
	}
	var out []*ooxml.Node
	if err := f.fillParas(paras, secs, 0, len(paras)-1, sc, at, &out); err != nil {
		return nil, err
	}
	if len(out) == 0 && len(paras) > 0 {
		// a text has a paragraph: an empty one in the format of the first
		p := *paras[0]
		p.Kids = nil
		for _, k := range paras[0].Kids {
			if k.Name == "pPr" || k.Name == "endParaRPr" {
				p.Kids = append(p.Kids, copyTree(k))
			}
		}
		out = []*ooxml.Node{&p}
	}
	done := false
	for _, k := range body.Kids {
		switch {
		case k.Name == "p":
			if !done {
				cp.Kids = append(cp.Kids, out...)
				done = true
			}
		case k.Name == "bodyPr" && changed:
			bp := copyTree(k)
			setAttr(bp, drawingml.RefitAttr, "1")
			cp.Kids = append(cp.Kids, bp)
		default:
			cp.Kids = append(cp.Kids, copyTree(k))
		}
	}
	return &cp, nil
}

// fillParas fills the paragraphs p0 to p1 of a text.
func (f *filler) fillParas(paras []*ooxml.Node, secs []*sec, p0, p1 int, sc *tag.Scope, at string, out *[]*ooxml.Node) error {
	i := p0
	one := func(p *ooxml.Node, sc *tag.Scope) error {
		n, err := f.fillPara(p, sc, at)
		if n != nil {
			*out = append(*out, n)
		}
		return err
	}
	for _, s := range secs {
		for ; i < s.a; i++ {
			if err := one(paras[i], sc); err != nil {
				return err
			}
		}
		i = s.b + 1
		v, err := f.lookup(tag.Token{Name: s.name, Path: s.path}, sc, at)
		if err != nil {
			return err
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
				isc = sc.Item(it, tag.Vars(k, len(items), false))
			}
			if err := f.fillParas(paras, s.children, s.a, s.b, isc, at, out); err != nil {
				return err
			}
		}
	}
	for ; i <= p1; i++ {
		if err := one(paras[i], sc); err != nil {
			return err
		}
	}
	return nil
}

// piece is a piece of a filled paragraph: text in the format of a run, or
// a child element kept as it is.
type piece struct {
	text []string
	run  *ooxml.Node
	node *ooxml.Node
}

// fillPara fills a paragraph; nil for one that held nothing but the tags
// of sections beyond it.
func (f *filler) fillPara(p *ooxml.Node, sc *tag.Scope, at string) (*ooxml.Node, error) {
	tpl := f.paras[p]
	if tpl == nil {
		if f.nodes += len(p.Kids); f.nodes > maxFilledNodes {
			return nil, fmt.Errorf("pptx: the filled template has more than %d elements", maxFilledNodes)
		}
		return copyTree(p), nil
	}
	if tpl.only {
		return nil, nil
	}
	items := tpl.items
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
		if n := len(pieces); n > 0 && pieces[n-1].node == nil && pieces[n-1].run == run {
			pieces[n-1].text = append(pieces[n-1].text, text)
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
			if it.node != nil {
				pieces = append(pieces, piece{node: it.node})
				continue
			}
			if it.outer {
				continue
			}
			switch t := it.tok; t.Kind {
			case tag.Text:
				emit(t.Text, it.run)
			case tag.Value:
				if _, ok := sc.Lookup(t.Path); !ok && (t.Name == varSlide || t.Name == varSlides) {
					// the number of the slide, and that of the slides: known
					// when the slides are drawn
					typ := "slidenum"
					if t.Name == varSlides {
						typ = fieldSlides
					}
					fld := &ooxml.Node{Space: it.run.Space, Name: "fld", Attrs: []ooxml.Attr{{Name: ooxml.Name{Local: "type"}, Value: typ}}}
					if rPr := it.run.Child("rPr"); rPr != nil {
						fld.Kids = append(fld.Kids, rPr)
					}
					fld.Kids = append(fld.Kids, &ooxml.Node{Space: it.run.Space, Name: "t", Text: "0"})
					pieces = append(pieces, piece{node: fld})
					continue
				}
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
					if items[j].node != nil || items[j].outer {
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
						if total >= maxParaText {
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
	cp := *p
	cp.Kids = nil
	var tail []*ooxml.Node
	for _, k := range p.Kids {
		switch k.Name {
		case "pPr":
			cp.Kids = append(cp.Kids, copyTree(k))
		case "endParaRPr":
			tail = append(tail, copyTree(k))
		}
	}
	for _, pc := range pieces {
		if pc.node != nil {
			cp.Kids = append(cp.Kids, copyTree(pc.node))
			continue
		}
		// a line break in a value is a break of the paragraph's line
		for i, line := range strings.Split(strings.ReplaceAll(strings.Join(pc.text, ""), "\r\n", "\n"), "\n") {
			if i > 0 {
				br := &ooxml.Node{Space: pc.run.Space, Name: "br"}
				if rPr := pc.run.Child("rPr"); rPr != nil {
					br.Kids = []*ooxml.Node{copyTree(rPr)}
				}
				cp.Kids = append(cp.Kids, br)
			}
			if line == "" {
				continue
			}
			r := *pc.run
			r.Kids = nil
			for _, k := range pc.run.Kids {
				if k.Name == "t" {
					r.Kids = append(r.Kids, &ooxml.Node{Space: k.Space, Name: "t", Attrs: k.Attrs, Text: line})
				} else {
					r.Kids = append(r.Kids, copyTree(k))
				}
			}
			if pc.run.Child("t") == nil {
				r.Kids = append(r.Kids, &ooxml.Node{Space: pc.run.Space, Name: "t", Text: line})
			}
			cp.Kids = append(cp.Kids, &r)
		}
	}
	cp.Kids = append(cp.Kids, tail...)
	f.nodes += len(cp.Kids)
	return &cp, nil
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

// fillTable fills a table: the rows of its sections are repeated for their
// items, or filled as slots.
func (f *filler) fillTable(tbl *ooxml.Node, sc *tag.Scope, at string) (*ooxml.Node, error) {
	cp := *tbl
	cp.Kids = nil
	var rows []*ooxml.Node
	for _, k := range tbl.Kids {
		if k.Name == "tr" {
			rows = append(rows, k)
		}
	}
	var secs []*sec
	if h := f.tables[tbl]; h != nil {
		secs = h.secs
	}
	var out []*ooxml.Node
	if err := f.fillRows(rows, secs, 0, len(rows)-1, sc, at, &out); err != nil {
		return nil, err
	}
	done := false
	for _, k := range tbl.Kids {
		if k.Name != "tr" {
			cp.Kids = append(cp.Kids, copyTree(k))
		} else if !done {
			cp.Kids = append(cp.Kids, out...)
			done = true
		}
	}
	// a table whose rows were repeated may go on on copies of the slide:
	// with the rows before the first repeated one, and before those after
	// the last
	head, tail, grows := -1, 0, false
	for _, s := range secs {
		if s.mode == tag.ModeRows && !s.inverted {
			if !grows {
				head = s.a
			}
			tail, grows = len(rows)-1-s.b, true
		}
	}
	if grows && head+tail <= len(out) {
		f.grown[&cp] = [2]int{head, tail}
	}
	return &cp, nil
}

// fillRows fills the rows r0 to r1 of a table.
func (f *filler) fillRows(rows []*ooxml.Node, secs []*sec, r0, r1 int, sc *tag.Scope, at string, out *[]*ooxml.Node) error {
	i := r0
	one := func(tr *ooxml.Node, sc *tag.Scope) error {
		n, err := f.clone(tr, sc, at)
		if n != nil {
			*out = append(*out, n)
		}
		return err
	}
	for _, s := range secs {
		for ; i < s.a; i++ {
			if err := one(rows[i], sc); err != nil {
				return err
			}
		}
		i = s.b + 1
		v, err := f.lookup(tag.Token{Name: s.name, Path: s.path}, sc, s.at)
		if err != nil {
			return err
		}
		if s.mode == tag.ModeSlots {
			items, n := tag.Items(v), s.b-s.a+1
			if len(items) > n {
				return errorAt(s.at, "%d items for the %d rows of {{#%s}}: the rows from this one to that of {{/%s}} are slots, a row for an item (write {{#%s | rows}} to repeat the rows instead)",
					len(items), n, s.name, s.name, s.name)
			}
			for k := 0; k < n; k++ {
				if err := f.expansionStep(s.at); err != nil {
					return err
				}
				isc := sc.Blank(map[string]any{"_blank": true, "_blank_first": k == len(items), "_index": float64(k + 1), "_count": float64(len(items))})
				if k < len(items) {
					v := tag.Vars(k, len(items), false)
					v["_blank"], v["_blank_first"] = false, false
					isc = sc.Item(items[k], v)
				}
				tr, err := f.slotRow(rows[s.a], rows[s.a+k], k, isc, at)
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
				isc = sc.Item(it, tag.Vars(k, len(items), false))
			}
			if err := f.fillRows(rows, s.children, s.a, s.b, isc, at, out); err != nil {
				return err
			}
		}
	}
	for ; i <= r1; i++ {
		if err := one(rows[i], sc); err != nil {
			return err
		}
	}
	return nil
}

// ownTags reports whether a row has tags of its own: not only those that
// close the section of its slots.
func (f *filler) ownTags(n *ooxml.Node) bool {
	if n.Name == "p" {
		tpl := f.paras[n]
		return tpl != nil && !tpl.only
	}
	for _, k := range n.Kids {
		if f.ownTags(k) {
			return true
		}
	}
	return false
}

// hasText reports whether a cell has text.
func hasText(n *ooxml.Node) bool {
	if n.Name == "t" {
		return n.Text != ""
	}
	for _, k := range n.Kids {
		if hasText(k) {
			return true
		}
	}
	return false
}

// slotRow fills slot k of a section of slots: the row at that place, with
// the texts of the cells of the first row in its cells that have no text
// when the row has no tags of its own.
func (f *filler) slotRow(pattern, tr *ooxml.Node, k int, sc *tag.Scope, at string) (*ooxml.Node, error) {
	if k == 0 || f.ownTags(tr) {
		return f.clone(tr, sc, at)
	}
	cp := *tr
	cp.Kids = nil
	first := pattern.Children("tc")
	c := 0
	for _, kid := range tr.Kids {
		if kid.Name != "tc" {
			cp.Kids = append(cp.Kids, copyTree(kid))
			continue
		}
		src := kid
		// (a cell that holds nothing but the tag that closes the section
		// has no text of its own)
		if c < len(first) && f.ownTags(first[c]) && (!hasText(kid) || !f.ownTags(kid) && f.tagOnly(kid)) {
			// the text of the first row in the format of this cell
			merged := *kid
			merged.Kids = nil
			for _, x := range kid.Kids {
				if x.Name == "txBody" {
					merged.Kids = append(merged.Kids, first[c].Child("txBody"))
				} else {
					merged.Kids = append(merged.Kids, x)
				}
			}
			src = &merged
		}
		n, err := f.clone(src, sc, at)
		if err != nil {
			return nil, err
		}
		cp.Kids = append(cp.Kids, n)
		c++
	}
	return &cp, nil
}

// tagOnly reports whether the paragraphs of an element all hold nothing
// but tags of sections beyond them.
func (f *filler) tagOnly(n *ooxml.Node) bool {
	if n.Name == "p" {
		tpl := f.paras[n]
		return tpl != nil && tpl.only || !hasText(n)
	}
	for _, k := range n.Kids {
		if !f.tagOnly(k) {
			return false
		}
	}
	return true
}

// fillPicture replaces a picture by the picture that its tag names, fitted
// in its frame; nil when the data names none.
func (f *filler) fillPicture(pic *ooxml.Node, pt picTag, sc *tag.Scope, at string) (*ooxml.Node, error) {
	t := pt.tok
	v, err := f.lookup(t, sc, at)
	if err != nil {
		return nil, err
	}
	fail := func(msg string) (*ooxml.Node, error) {
		if f.strict {
			return nil, errorAt(at, "the picture of %q %s", t.Name, msg)
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
				return nil, errorAt(at, "%v", err)
			}
			f.c.warnf("%s: the code of %q is left out: %v", at, t.Name, err)
			return nil, nil
		}
		rid = f.c.r.InjectDrawing(func(o *bdf.Object, x, y, w, h float32) {
			o.FillColor(bdf.RGB(0, 0, 0))
			sym.Fill(o, x, y, w, h)
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
		in, ok := f.loaded[name]
		if !ok {
			// (read and stored once: the same picture may be in every item)
			if data, err := tag.PictureData(f.files, name); err != nil {
				in.why = "is not read: " + err.Error()
			} else if in.rid, in.w, in.h, in.ok = f.c.r.InjectImage(data); !in.ok {
				in.why = "is not a PNG, JPEG, GIF, BMP or WebP image"
			}
			f.loaded[name] = in
		}
		if !in.ok {
			return fail(in.why)
		}
		rid, w, h = in.rid, in.w, in.h
	}
	cp := copyTree(pic)
	if nv := cNvPrOf(cp); nv != nil {
		setAttr(nv, "descr", alt)
	}
	bf := cp.Child("blipFill")
	blip := bf.Child("blip")
	if blip == nil {
		return fail("has no place in the template's picture")
	}
	set := false
	for i := range blip.Attrs {
		if a := &blip.Attrs[i]; a.Name.Local == "embed" || a.Name.Local == "link" {
			a.Name.Local, a.Value, set = "embed", rid, true
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
	xfrm := cp.Path("spPr", "xfrm")
	off, ext := xfrm.Child("off"), xfrm.Child("ext")
	if off == nil || ext == nil {
		return cp, nil
	}
	bw, bh := float64(ext.AttrInt("cx", 0)), float64(ext.AttrInt("cy", 0))
	if bw <= 0 || bh <= 0 {
		return cp, nil
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
		// the whole picture in the middle of the frame
		k := math.Min(bw/w, bh/h)
		nw, nh := w*k, h*k
		setAttr(off, "x", emu(float64(off.AttrInt("x", 0))+(bw-nw)/2))
		setAttr(off, "y", emu(float64(off.AttrInt("y", 0))+(bh-nh)/2))
		setAttr(ext, "cx", emu(nw))
		setAttr(ext, "cy", emu(nh))
	}
	return cp, nil
}

// frameOf finds the frame (p:graphicFrame) of a table of a slide.
func frameOf(n, tbl *ooxml.Node) *ooxml.Node {
	if n.Name == "graphicFrame" && n.Path("graphic", "graphicData", "tbl") == tbl {
		return n
	}
	for _, k := range n.Kids {
		if f := frameOf(k, tbl); f != nil {
			return f
		}
	}
	return nil
}

// grownTables lists the tables of a filled slide whose rows were repeated.
func (f *filler) grownTables(n *ooxml.Node, out *[]*ooxml.Node) {
	if n.Name == "tbl" {
		if _, ok := f.grown[n]; ok {
			*out = append(*out, n)
		}
		return
	}
	for _, k := range n.Kids {
		f.grownTables(k, out)
	}
}

// paginate cuts a filled slide whose table reaches past the foot of the
// slide into copies of the slide, each with the rows that fit: the rows
// before the repeated ones at the top of each copy, those after them on the
// last.
func (f *filler) paginate(j slideJob) ([]slideJob, error) {
	var tables []*ooxml.Node
	f.grownTables(j.node, &tables)
	if len(tables) == 0 {
		return []slideJob{j}, nil
	}
	tbl := tables[0]
	tree := j.node.Path("cSld", "spTree")
	frame := frameOf(tree, tbl)
	if frame == nil {
		return []slideJob{j}, nil
	}
	s, err := f.c.newSlideCtx(j.part, j.node, 0)
	if err != nil {
		return nil, err
	}
	hs := s.d.TableRowHeights(tbl, j.part)
	rows := tbl.Children("tr")
	if len(hs) != len(rows) {
		return []slideJob{j}, nil
	}
	x0, y0 := frame.Path("xfrm", "off").AttrEMU("x", 0), frame.Path("xfrm", "off").AttrEMU("y", 0)
	foot := y0 + frame.Path("xfrm", "ext").AttrEMU("cy", 0) // of the table as the template has it
	// the table may reach down to the shape under it, or to the foot of the
	// slide less the margin it has at its side
	limit := f.c.slideH - math.Min(math.Max(x0, 0), f.c.slideH*0.1)
	for _, k := range tree.Kids {
		if k == frame {
			continue
		}
		if off := k.Path("spPr", "xfrm", "off"); off != nil {
			if y := off.AttrEMU("y", 0); y >= foot-1 && y < limit {
				limit = y
			}
		}
	}
	total := 0.0
	for _, h := range hs {
		total += h
	}
	if y0+total <= limit+0.5 {
		return []slideJob{j}, nil
	}
	ht := f.grown[tbl]
	head, tail := ht[0], ht[1]
	headH, tailH := 0.0, 0.0
	for i := 0; i < head; i++ {
		headH += hs[i]
	}
	for i := len(rows) - tail; i < len(rows); i++ {
		tailH += hs[i]
	}
	room := limit - y0 - headH
	// the rows of each copy
	var chunks [][2]int
	lo, used := head, 0.0
	end := len(rows) - tail
	for i := head; i < end; i++ {
		if i > lo && used+hs[i] > room+0.5 {
			chunks = append(chunks, [2]int{lo, i})
			lo, used = i, 0
		}
		used += hs[i]
	}
	chunks = append(chunks, [2]int{lo, end})
	// the rows after the repeated ones: on the last copy when they fit
	lastOwn := tail > 0 && used+tailH > room+0.5 && end > lo
	var out []slideJob
	for ci, ch := range chunks {
		keep := append([]*ooxml.Node(nil), rows[:head]...)
		keep = append(keep, rows[ch[0]:ch[1]]...)
		if ci == len(chunks)-1 && !lastOwn {
			keep = append(keep, rows[end:]...)
		}
		out = append(out, slideJob{part: j.part, node: replaceRows(j.node, tbl, keep)})
	}
	if lastOwn {
		keep := append(append([]*ooxml.Node(nil), rows[:head]...), rows[end:]...)
		out = append(out, slideJob{part: j.part, node: replaceRows(j.node, tbl, keep)})
	}
	return out, nil
}

// replaceRows copies the elements from a slide down to a table, whose rows
// become the rows given; what is beside them is shared.
func replaceRows(n, tbl *ooxml.Node, rows []*ooxml.Node) *ooxml.Node {
	if n == tbl {
		cp := *n
		cp.Kids = nil
		done := false
		for _, k := range n.Kids {
			if k.Name != "tr" {
				cp.Kids = append(cp.Kids, k)
			} else if !done {
				cp.Kids = append(cp.Kids, rows...)
				done = true
			}
		}
		return &cp
	}
	for i, k := range n.Kids {
		if r := replaceRows(k, tbl, rows); r != k {
			cp := *n
			cp.Kids = append([]*ooxml.Node(nil), n.Kids...)
			cp.Kids[i] = r
			return &cp
		}
	}
	return n
}
