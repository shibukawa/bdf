package textform

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf/form/tag"
	"github.com/shibukawa/bdf/form/textpart"
)

// cell is a character of a filled line: a character and the characters of
// no width that follow it.
type cell struct {
	text string
	// w is the columns of the character (textpart.Width), and attrs the
	// flags of the sections of attributes it is in.
	w     int
	attrs uint8
}

// cols is the columns the character takes as printed.
func (c cell) cols() int {
	if c.attrs&attrWide != 0 {
		return 2 * c.w
	}
	return c.w
}

func colsOf(cells []cell) int {
	n := 0
	for _, c := range cells {
		n += c.cols()
	}
	return n
}

// outLine is a line of the form, or a page break.
type outLine struct {
	cells []cell
	brk   bool
	// table is the table of the line (the expansion of a section of rows
	// or slots; 0 for none), and no the line of the template.
	table int
	no    int
}

// tall reports whether a line takes two lines of the page.
func (o *outLine) tall() bool {
	for _, c := range o.cells {
		if c.attrs&attrTall != 0 {
			return true
		}
	}
	return false
}

// cont is the rest of a value that goes on to the lines under its own: a
// piece for each line, in the columns from col.
type cont struct {
	col    int
	pieces [][]cell
}

// lineBuilder is a line being filled.
type lineBuilder struct {
	cells []cell
	cols  int
	conts []cont
}

func (b *lineBuilder) add(c cell) {
	b.cells = append(b.cells, c)
	b.cols += c.cols()
}

// Limits of a filled form.
const (
	maxOutLines = 1 << 20
	maxLineCols = 1 << 16
	// maxSteps bounds the repetitions of sections, which may print
	// nothing
	maxSteps = 8 << 20
)

// filler fills a template with data.
type filler struct {
	t        *template
	td       *TemplateData
	format   tag.Formatter
	out      []outLine
	table    int
	tables   int
	warnings []string
	warned   map[string]bool
	// raw lays the lines of the template out as they are, with their tags.
	raw bool
	// steps counts the repetitions of sections
	steps int
}

// step counts a repetition of a section.
func (f *filler) step(at string) error {
	if f.steps++; f.steps > maxSteps {
		return fmt.Errorf("text: %s: the filled template needs more than %d repetitions of sections", at, maxSteps)
	}
	return nil
}

func (f *filler) warnf(format string, args ...any) {
	f.warnings = append(f.warnings, fmt.Sprintf(format, args...))
}

func (f *filler) warnOnce(key, format string, args ...any) {
	if f.warned[key] {
		return
	}
	f.warned[key] = true
	f.warnf(format, args...)
}

// lookup finds the value of the name of a tag; a name the data does not
// have is told.
func (f *filler) lookup(name string, path []string, sc *tag.Scope, at string) (any, error) {
	v, ok := sc.Lookup(path)
	if ok || sc.IsBlank() {
		return v, nil
	}
	if f.td.Strict {
		return nil, fmt.Errorf("text: %s: the data has no %q", at, name)
	}
	f.warnOnce("field:"+name, "%s: the data has no %q; it is left empty", at, name)
	return nil, nil
}

// lines fills the lines from to to (indexes of the lines of the template)
// of a section, with the sections of lines in them.
func (f *filler) lines(sec *section, from, to int, sc *tag.Scope) error {
	ci := 0
	for i := from; i <= to; {
		for ci < len(sec.children) && sec.children[ci].a < i {
			ci++
		}
		if ci < len(sec.children) && sec.children[ci].a == i {
			ch := sec.children[ci]
			if err := f.section(ch, sc); err != nil {
				return err
			}
			i = ch.b + 1
			continue
		}
		l := f.t.lines[i]
		if err := f.emit(l, l.attrs, sc); err != nil {
			return err
		}
		i++
	}
	return nil
}

// section fills a section of lines.
func (f *filler) section(s *section, sc *tag.Scope) error {
	v, err := f.lookup(s.name, s.path, sc, s.at)
	if err != nil {
		return err
	}
	if !s.inverted && (s.mode == tag.ModeRows || s.mode == tag.ModeSlots) && f.table == 0 {
		f.tables++
		f.table = f.tables
		defer func() { f.table = 0 }()
	}
	if s.mode == tag.ModeSlots {
		items, slots := tag.Items(v), s.slots(f.t)
		if len(items) > len(slots) {
			return fmt.Errorf("text: %s: %d items for the %d lines of {{#%s}}: the lines from this one to that of {{/%s}} are slots, a line for an item (write {{#%s | rows}} to repeat the lines instead)",
				s.at, len(items), len(slots), s.name, s.name, s.name)
		}
		first := f.t.lines[slots[0]]
		for k, li := range slots {
			if err := f.step(s.at); err != nil {
				return err
			}
			isc := sc.Blank(map[string]any{"_blank": true, "_blank_first": k == len(items), "_index": float64(k + 1), "_count": float64(len(items))})
			if k < len(items) {
				vars := tag.Vars(k, len(items), false)
				vars["_blank"], vars["_blank_first"] = false, false
				isc = sc.Item(items[k], vars)
			}
			// (a line of the slots without tags of its own is the first)
			l, tpl := f.t.lines[li], f.t.lines[li]
			if !tpl.tagged {
				tpl = first
			}
			if err := f.emit(tpl, l.attrs, isc); err != nil {
				return err
			}
		}
		return nil
	}
	var items []any
	switch {
	case s.inverted && !tag.Truthy(v):
		items = []any{nil}
	case !s.inverted:
		items = tag.Items(v)
	}
	for k, it := range items {
		if err := f.step(s.at); err != nil {
			return err
		}
		isc := sc
		if !s.inverted {
			isc = sc.Item(it, tag.Vars(k, len(items), s.mode == tag.ModePages))
		}
		if s.mode == tag.ModePages && k > 0 {
			f.out = append(f.out, outLine{brk: true})
		}
		if err := f.lines(s, s.a, s.b, isc); err != nil {
			return err
		}
	}
	return nil
}

// emit fills a line of the template (tpl, whose tags are those of the line)
// and adds it to the form, with the lines that the rest of its values go
// on to.
func (f *filler) emit(tpl *line, attrs uint8, sc *tag.Scope) error {
	if tpl.brk {
		f.out = append(f.out, outLine{brk: true})
		return nil
	}
	if tpl.tagOnly && !f.raw {
		return nil
	}
	b := &lineBuilder{}
	if f.raw {
		f.text(b, tpl.src, 0)
	} else if err := f.tokens(tpl, 0, len(tpl.toks), sc, attrs, b); err != nil {
		return err
	}
	f.out = append(f.out, outLine{cells: trim(b.cells), table: f.table, no: tpl.no})
	more := 0
	for _, c := range b.conts {
		more = max(more, len(c.pieces))
	}
	for k := range more {
		var cl lineBuilder
		for _, c := range b.conts {
			if k >= len(c.pieces) {
				continue
			}
			for cl.cols < c.col {
				cl.add(cell{text: " ", w: 1})
			}
			for _, x := range c.pieces[k] {
				cl.add(x)
			}
		}
		f.out = append(f.out, outLine{cells: trim(cl.cells), table: f.table, no: tpl.no})
	}
	if len(f.out) > maxOutLines {
		return fmt.Errorf("text: the filled form has more than %d lines", maxOutLines)
	}
	if b.cols > f.t.grid.width() {
		f.warnOnce(fmt.Sprintf("wide:%d", tpl.no), "line %d: the filled line is %d columns; the paper has %d after the left margin", tpl.no, b.cols, f.t.grid.width())
	}
	return nil
}

// trim drops the spaces at the end of a line.
func trim(cells []cell) []cell {
	for len(cells) > 0 && cells[len(cells)-1].text == " " {
		cells = cells[:len(cells)-1]
	}
	return cells
}

// tokens fills the tokens from to to of a line.
func (f *filler) tokens(l *line, from, to int, sc *tag.Scope, attrs uint8, b *lineBuilder) error {
	for i := from; i < to; i++ {
		tk := l.toks[i]
		if l.lineLevel[i] {
			continue
		}
		switch tk.Kind {
		case tag.Text:
			f.text(b, tk.Text, attrs)
		case tag.Value:
			v, err := f.lookup(tk.Name, tk.Path, sc, l.at())
			if err != nil {
				return err
			}
			f.value(b, tk, tag.ValueText(v, tk.Format, f.format), attrs, l)
		case tag.Open, tag.Inverted:
			end := l.match[i]
			if a := attrOf(tk.Name); a != 0 {
				if err := f.tokens(l, i+1, end, sc, attrs|a, b); err != nil {
					return err
				}
				i = end
				continue
			}
			v, err := f.lookup(tk.Name, tk.Path, sc, l.at())
			if err != nil {
				return err
			}
			if tk.Kind == tag.Inverted {
				if !tag.Truthy(v) {
					if err := f.tokens(l, i+1, end, sc, attrs, b); err != nil {
						return err
					}
				}
				i = end
				continue
			}
			items := tag.Items(v)
			for k, it := range items {
				if err := f.step(l.at()); err != nil {
					return err
				}
				if err := f.tokens(l, i+1, end, sc.Item(it, tag.Vars(k, len(items), false)), attrs, b); err != nil {
					return err
				}
			}
			i = end
		}
		if b.cols > maxLineCols {
			return fmt.Errorf("text: %s: the filled line is more than %d columns", l.at(), maxLineCols)
		}
	}
	return nil
}

// cells makes the cells of a text: tabs to the next stop of 8 columns (from
// col), control characters dropped, the characters of no width with the
// one before them.
func (f *filler) cells(s string, attrs uint8, col int) []cell {
	var out []cell
	wide := f.t.grid.ambiguousWide
	for _, r := range s {
		switch {
		case r == '\t':
			n := 8 - col%8
			for range n {
				out = append(out, cell{text: " ", w: 1, attrs: attrs &^ attrWide})
			}
			col += n
			continue
		case r < 0x20 || r == 0x7f || r >= 0x80 && r < 0xa0:
			continue
		}
		w := textpart.Width(r, wide)
		switch {
		case unicode.Is(unicode.Co, r):
			f.warnOnce("pua", "the character U+%04X is of the private use area (a character of a font of its own): the text mode of printers does not print it", r)
		case r >= 0x1f000 && r < 0x1fb00 || r >= 0x2600 && r < 0x2800 && w == 2:
			f.warnOnce("emoji", "the character %c (U+%04X) is a pictograph (an emoji): the text mode of printers does not print it", r, r)
		}
		if w == 0 && len(out) > 0 {
			out[len(out)-1].text += string(r)
			continue
		}
		c := cell{text: string(r), w: w, attrs: attrs}
		out = append(out, c)
		col += c.cols()
	}
	return out
}

// text adds text to a line.
func (f *filler) text(b *lineBuilder, s string, attrs uint8) {
	for _, c := range f.cells(s, attrs, b.cols) {
		if c.w == 0 && len(b.cells) > 0 {
			b.cells[len(b.cells)-1].text += c.text
			continue
		}
		b.add(c)
	}
}

// value adds a value to a line: in the columns of its tag, when it has
// them, and with the rest of it on the lines under it when it wraps or has
// line breaks.
func (f *filler) value(b *lineBuilder, tk tag.Token, s string, attrs uint8, l *line) {
	s = strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\r", "\n")
	paras := strings.Split(s, "\n")
	if tk.Width > 0 && tk.Overflow != tag.OverflowWrap && len(paras) > 1 {
		// (a value cut or kept whole is one line)
		paras = []string{strings.Join(paras, " ")}
	}
	col := b.cols
	var pieces [][]cell
	for _, p := range paras {
		cs := f.cells(p, attrs, col)
		switch {
		case tk.Width == 0:
			pieces = append(pieces, cs)
		case tk.Overflow == tag.OverflowWrap:
			for len(cs) > 0 {
				n, w := 0, 0
				for n < len(cs) && (n == 0 || w+cs[n].cols() <= tk.Width) {
					w += cs[n].cols()
					n++
				}
				pieces = append(pieces, pad(cs[:n], tk, attrs))
				cs = cs[n:]
			}
			if len(cs) == 0 && (len(pieces) == 0 || len(p) == 0) {
				pieces = append(pieces, pad(nil, tk, attrs))
			}
		case tk.Overflow == tag.OverflowCut:
			n, w := 0, 0
			for n < len(cs) && w+cs[n].cols() <= tk.Width {
				w += cs[n].cols()
				n++
			}
			pieces = append(pieces, pad(cs[:n], tk, attrs))
		default:
			if w := colsOf(cs); w > tk.Width {
				f.warnOnce("over:"+tk.Name, "%s: the value of {{%s}} is %d columns, more than its %d; it is written whole (| cut cuts it, | wrap goes on to the next line)", l.at(), tk.Name, w, tk.Width)
			}
			pieces = append(pieces, pad(cs, tk, attrs))
		}
	}
	if len(pieces) == 0 {
		pieces = [][]cell{pad(nil, tk, attrs)}
	}
	for _, c := range pieces[0] {
		b.add(c)
	}
	if len(pieces) > 1 {
		b.conts = append(b.conts, cont{col: col, pieces: pieces[1:]})
	}
}

// pad puts cells in the columns of a tag (none: as they are). The spaces
// it adds are of single width, so that they fill an odd number of columns.
func pad(cs []cell, tk tag.Token, attrs uint8) []cell {
	if tk.Width == 0 {
		return cs
	}
	room := tk.Width - colsOf(cs)
	if room <= 0 {
		return cs
	}
	sp := func(n int) []cell {
		out := make([]cell, n)
		for i := range out {
			out[i] = cell{text: " ", w: 1, attrs: attrs &^ attrWide}
		}
		return out
	}
	left := 0
	switch tk.Align {
	case tag.AlignRight:
		left = room
	case tag.AlignCenter:
		left = room / 2
	}
	out := append(sp(left), cs...)
	return append(out, sp(room-left)...)
}

// page is a page of the form: its lines, and the tables they are in.
type page struct {
	lines  []outLine
	tables []int
}

// paginate puts the lines of the form on pages: a page breaks at a form
// feed, at a copy of a section of pages, and, when the form grows, where it
// is full; the lines of the title start each page after the first.
func (f *filler) paginate() ([]page, error) {
	g := &f.t.grid
	rows := g.rows()
	var pages [][]outLine
	var cur, title []outLine
	titleOpen := g.title > 0
	newPage := func() {
		pages = append(pages, cur)
		cur = append([]outLine(nil), title...)
	}
	used := func(ls []outLine) int {
		n := 0
		for i := range ls {
			n++
			if ls[i].tall() {
				n++
			}
		}
		return n
	}
	for _, ol := range f.out {
		if ol.brk {
			titleOpen = false
			newPage()
			continue
		}
		need := 1
		if ol.tall() {
			need = 2
		}
		if rows > 0 && used(cur)+need > rows {
			switch {
			case len(ol.cells) == 0 && !f.t.grows && !f.raw:
				// (empty lines past the end of a page of a fixed form are
				// none of its lines)
				continue
			case !f.t.grows && !f.raw:
				return nil, fmt.Errorf("text: page %d: line %d of the template is past the %d lines of a page (%d lines of paper less the top margin %d and the skip %d); only a form with a section of | rows goes on to another page by itself",
					len(pages)+1, ol.no, rows, g.lines, g.top, g.skip)
			}
			titleOpen = false
			newPage()
			if used(cur)+need > rows {
				return nil, fmt.Errorf("text: page %d: the title leaves no room for line %d of the template", len(pages)+1, ol.no)
			}
		}
		cur = append(cur, ol)
		if titleOpen {
			title = append(title, ol)
			if used(title) >= g.title {
				titleOpen = false
			}
		}
	}
	// (a break at the end of the form makes no empty page)
	if len(pages) == 0 || len(cur) > len(title) || len(f.out) > 0 && !f.out[len(f.out)-1].brk {
		pages = append(pages, cur)
	}
	out := make([]page, len(pages))
	for i, ls := range pages {
		for _, ol := range ls {
			out[i].lines = append(out[i].lines, ol)
			out[i].tables = append(out[i].tables, ol.table)
		}
	}
	return out, nil
}

// image makes the print image of the pages.
func (f *filler) image(pages []page) *textpart.Image {
	g := &f.t.grid
	im := &textpart.Image{Version: textpart.Version, CPI: g.cpi, LPI: g.lpi, Columns: g.columns, Lines: g.lines,
		Margin: textpart.Margin{Left: g.left, Top: g.top}, Skip: g.skip, AmbiguousWide: g.ambiguousWide, Pages: []textpart.Page{}}
	for _, p := range pages {
		tp := textpart.Page{Lines: []textpart.Line{}}
		for _, ol := range p.lines {
			tp.Lines = append(tp.Lines, partLine(ol.cells))
			if ol.tall() {
				tp.Lines = append(tp.Lines, textpart.Line{})
			}
		}
		im.Pages = append(im.Pages, tp)
	}
	return im
}

// partLine makes a line of the print image of the cells of a line.
func partLine(cells []cell) textpart.Line {
	var l textpart.Line
	var b strings.Builder
	col := 0
	var paper, size *textpart.Range
	for _, c := range cells {
		b.WriteString(c.text)
		n := c.cols()
		if c.attrs&attrPaper != 0 {
			if paper != nil && paper.Col+paper.Cols == col {
				paper.Cols += n
			} else if n > 0 {
				l.Paper = append(l.Paper, textpart.Range{Col: col, Cols: n})
				paper = &l.Paper[len(l.Paper)-1]
			}
		}
		if s := sizeOf(c.attrs); s != "" {
			if size != nil && size.Size == s && size.Col+size.Cols == col {
				size.Cols += n
			} else if n > 0 {
				l.Size = append(l.Size, textpart.Range{Col: col, Cols: n, Size: s})
				size = &l.Size[len(l.Size)-1]
			}
		}
		col += n
	}
	l.Text = b.String()
	return l
}

// sizeOf is the size of characters of flags.
func sizeOf(attrs uint8) string {
	switch attrs & (attrWide | attrTall) {
	case attrWide:
		return textpart.SizeWide
	case attrTall:
		return textpart.SizeTall
	case attrWide | attrTall:
		return textpart.SizeBig
	}
	return ""
}
