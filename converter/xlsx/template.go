package xlsx

import (
	"fmt"
	"io"
	"io/fs"
	"sort"
	"strings"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/form/tag"
	"github.com/shibukawa/bdf/image/barcode"
)

// A workbook is a template when its cells hold tags (package form/tag):
// Fill converts it with the tags replaced by data, as the pages it prints.
// The workbook is filled as its rows and cells in memory, before it is laid
// out; no workbook is written.
//
// A tag in a cell is replaced by its value. A cell that is one tag and
// nothing else takes the value as it is: a number is a number of the cell's
// number format (and of the formulas that read the cell), a date a date.
//
// A section whose tags are in different cells repeats rows:
//
//   - {{#items}} and {{/items}} in the same row: the row is there once for
//     each item, and what is below it moves down (tag.ModeRows; with tags in
//     different rows and "| rows", the rows between them are).
//   - {{#rows}} and {{/rows}} in different rows: the rows from the one to
//     the other are slots, a row for an item (tag.ModeSlots). Nothing moves;
//     the rows without an item stay as they are drawn, empty. The tags of
//     the first row are those of every row that has none of its own.
//   - {{#pages | pages}}: as ModeRows, with each copy on a page of its own
//     (tag.ModePages).
//
// The cells of the rows of a section see the item; a row inside sections
// sees the items of each. Sections nest: pages of slots, groups of rows.
//
// What refers to rows follows them as it does when Excel inserts rows:
// merged cells, conditional formats, links, print areas and titles, page
// breaks, drawings, and the references of formulas. A formula in a copied
// row refers to the cells of its own copy; a range that holds a repeated
// row holds all its copies. The formulas are then calculated again (those
// the evaluator of conditional formats knows; the others are left empty
// with a warning).

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
}

// Limits of a filled template.
const (
	// maxFilledCells bounds the cells of the sheets of a filled template.
	maxFilledCells = 4 << 20
	// maxCellText bounds the characters of a cell (Excel's own limit).
	maxCellText = 32767
	// maxCopies bounds the ranges that a merged range, a conditional
	// format or a link becomes in the copies of the rows it is in.
	maxCopies = 1 << 16
	// maxPrecedents bounds the cells looked at to order the formulas.
	maxPrecedents = 16 << 20
	// maxSectionDepth bounds the sections of rows one inside the other.
	maxSectionDepth = 32
	// maxMovedRanges bounds the merged ranges, the ranges of conditional
	// formats and the links of a filled sheet, each.
	maxMovedRanges = 1 << 20
)

// Fill converts a workbook as a template filled with data. The result is
// the pages the filled sheets print (ViewsPrint) unless the options select
// other views.
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

// Inspect reads the tags of a workbook without filling it: the fields that
// data for it has, or the first mistake in its tags.
func Inspect(r io.ReaderAt, size int64, opts *Options) ([]tag.Field, error) {
	o := Options{}
	if opts != nil {
		o = *opts
	}
	f := newFiller(&TemplateData{})
	o.fill = f
	c, _, err := openWorkbook(r, size, &o)
	if err != nil {
		return nil, err
	}
	f.c = c
	if err := f.scanAll(); err != nil {
		return nil, err
	}
	return f.fieldList(), nil
}

// cellFormula is the formula of a cell of a template.
type cellFormula struct {
	src  string
	kind string // the t attribute: "" (normal), shared, array, dataTable
	si   int    // the group of a shared formula, from 1 (0: none)
	// node is the parsed formula (nil when the evaluator does not read
	// it), and dr and dc how far the cell is from the one the formula was
	// written for (a shared formula)
	node   fnode
	dr, dc int
}

// tplCell is the tags of a cell.
type tplCell struct {
	toks []tag.Token
	// rowLevel marks the tags that open and close sections of rows
	rowLevel []bool
}

// tplSection is a section of rows: the rows r0 to r1 of the template.
type tplSection struct {
	name     string
	path     []string
	inverted bool
	mode     string
	r0, r1   int
	at       [2]int // the cell of its opening tag
	parent   *tplSection
	children []*tplSection
	// insts are the copies of its rows in the filled sheet
	insts []*itemInst
}

// secInst is the copies of the rows of a section in a copy of the rows
// around it: the rows dst0 to dst1-1 of the filled sheet.
type secInst struct {
	sec        *tplSection
	items      []*itemInst
	dst0, dst1 int
}

// itemInst is a copy of the rows of a section (or the sheet itself): the
// rows dst0 to dst1-1 of the filled sheet.
type itemInst struct {
	parent     *itemInst
	of         *secInst // nil for the sheet
	dst0, dst1 int
	children   []*secInst // of the sections in the rows that repeat rows, in order
}

// outFormula is a formula of a cell of a filled sheet.
type outFormula struct {
	f     *cellFormula
	dr    int       // rows the cell is from the one the formula was read in (a slot filled from the first)
	ctx   *itemInst // the copy of rows the cell is in
	state uint8     // 0 to do, 1 being calculated, 2 done
}

// tplPicture is a picture over a cell of a filled sheet.
type tplPicture struct {
	row, col int
	hash     bdf.Hash
	w, h     float64
	fit      string
	alt      string
	// sym is the QR code or the barcode that is drawn in the place of a
	// picture
	sym *barcode.Symbol
}

// tplSheet is a sheet of a template.
type tplSheet struct {
	ws    *worksheet
	ref   sheetRef
	cells map[[2]int]*tplCell
	root  *tplSection
	// the filled sheet
	inst     *itemInst
	outF     map[[2]int]*outFormula
	pictures []tplPicture
	breaks   []int // rows that start a page: the copies of ModePages
	calc     *sheetCtx
}

type styleKey struct {
	style int
	code  string
}

// filler fills a workbook as a template.
type filler struct {
	c      *converter
	root   *tag.Scope
	strict bool
	files  fs.FS

	sheets []*tplSheet
	byName map[string]*tplSheet

	fields     map[string]*tag.Field
	fieldOrder []string
	styles     map[styleKey]int
	// fmts are the number formats of tags, parsed once each
	fmts     map[string]*numFormat
	nextFmt  int
	cells    int // cells of the filled sheets so far
	pictures map[string]*tplPicture
	looked   int // cells looked at to order the formulas
}

func newFiller(t *TemplateData) *filler {
	data := t.Data
	if data == nil {
		data = map[string]any{}
	}
	return &filler{root: tag.NewScope(tag.Normalize(data)), strict: t.Strict, files: t.Files, byName: map[string]*tplSheet{},
		fields: map[string]*tag.Field{}, styles: map[styleKey]int{}, fmts: map[string]*numFormat{}, nextFmt: 60000, pictures: map[string]*tplPicture{}}
}

// numFmt returns the number format of a tag ({{v | fmt:"…"}}), parsed
// once: a tag in a repeated row is filled once for each item.
func (f *filler) numFmt(code string) *numFormat {
	if nf, ok := f.fmts[code]; ok {
		return nf
	}
	nf := parseNumFormat(code, f.c.st.palette)
	f.fmts[code] = nf
	return nf
}

// errorf makes an error about a cell of a sheet.
func (f *filler) errorf(ts *tplSheet, r, c int, format string, args ...any) error {
	return fmt.Errorf("xlsx: sheet %q cell %s: %s", ts.ws.name, cellRef(r, c), fmt.Sprintf(format, args...))
}

// run fills the workbook: every sheet is read and filled, then what refers
// to its rows is moved with them and the formulas are calculated again.
func (f *filler) run(c *converter) error {
	f.c = c
	if err := f.scanAll(); err != nil {
		return err
	}
	for _, ts := range f.sheets {
		if err := f.expand(ts); err != nil {
			return err
		}
	}
	for _, ts := range f.sheets {
		f.remap(ts)
	}
	f.recalc()
	return nil
}

// scanAll reads the worksheets and their tags.
func (f *filler) scanAll() error {
	for _, ref := range f.c.sheets {
		if ref.chart {
			continue
		}
		ws, err := f.c.sheetOf(ref)
		if err != nil {
			return fmt.Errorf("xlsx: sheet %q: %w", ref.name, err)
		}
		ts := &tplSheet{ws: ws, ref: ref, cells: map[[2]int]*tplCell{}, outF: map[[2]int]*outFormula{}}
		ts.calc = &sheetCtx{c: f.c, ws: ws, formulas: map[string]fnode{}}
		f.sheets = append(f.sheets, ts)
		f.byName[ws.name] = ts
		f.resolveFormulas(ws)
		ws.noValue = 0 // the formulas are calculated when the template is filled
		if err := f.scan(ts); err != nil {
			return err
		}
	}
	return nil
}

// resolveFormulas parses the formulas of a sheet, and gives the cells of a
// shared formula the formula of the cell that holds it.
func (f *filler) resolveFormulas(ws *worksheet) {
	type master struct {
		cf   *cellFormula
		r, c int
	}
	shared := map[int]master{}
	keys := make([][2]int, 0, len(ws.formulas))
	for k := range ws.formulas {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i][0] < keys[j][0] || keys[i][0] == keys[j][0] && keys[i][1] < keys[j][1]
	})
	for _, k := range keys {
		cf := ws.formulas[k]
		if cf.kind == "array" || cf.kind == "dataTable" {
			continue // not calculated
		}
		if cf.src != "" {
			if node, err := parseFormula(cf.src); err == nil {
				cf.node = node
			}
			if cf.si > 0 {
				shared[cf.si] = master{cf, k[0], k[1]}
			}
			continue
		}
		if m, ok := shared[cf.si]; ok && cf.si > 0 {
			cf.src, cf.node, cf.dr, cf.dc = m.cf.src, m.cf.node, k[0]-m.r, k[1]-m.c
		}
	}
}

// scan reads the tags of the cells of a sheet and pairs the tags of its
// sections of rows.
func (f *filler) scan(ts *tplSheet) error {
	ws := ts.ws
	root := &tplSection{r1: maxRows - 1}
	ts.root = root
	var stack []*tplSection
	for i := range ws.rows {
		rw := &ws.rows[i]
		for j := range rw.cells {
			cl := &rw.cells[j]
			if cl.kind != cellStr || cl.text == nil || !tag.Has(cl.text.plain) {
				continue
			}
			toks, err := tag.Parse(cl.text.plain)
			if err != nil {
				return f.errorf(ts, rw.idx, cl.col, "%v", err)
			}
			tc := &tplCell{toks: toks, rowLevel: make([]bool, len(toks))}
			tagged := false
			var local []int
			// (lifted: the sections of the cell that repeat its row, and the
			// tokens that open them)
			var lifted []*tplSection
			var liftedAt []int
			for k, t := range toks {
				tagged = tagged || t.Kind != tag.Text
				switch t.Kind {
				case tag.Open, tag.Inverted:
					local = append(local, k)
				case tag.Close:
					if n := len(local); n > 0 {
						oi := local[n-1]
						o := toks[oi]
						if t.Name != "" && t.Name != o.Name {
							return f.errorf(ts, rw.idx, cl.col, "{{/%s}} closes the section {{#%s}}", t.Name, o.Name)
						}
						local = local[:n-1]
						if o.Mode != "" {
							// a section that says how it repeats repeats the
							// row, not the text between its tags
							sec := &tplSection{name: o.Name, path: o.Path, inverted: o.Kind == tag.Inverted, mode: o.Mode,
								r0: rw.idx, r1: rw.idx, at: [2]int{rw.idx, cl.col}}
							// (those between its tags are in it)
							m := len(lifted)
							for m > 0 && liftedAt[m-1] > oi {
								m--
							}
							for _, in := range lifted[m:] {
								in.parent = sec
								sec.children = append(sec.children, in)
							}
							lifted, liftedAt = append(lifted[:m], sec), append(liftedAt[:m], oi)
							tc.rowLevel[oi], tc.rowLevel[k] = true, true
						}
						continue
					}
					if len(stack) == 0 {
						return f.errorf(ts, rw.idx, cl.col, "{{/%s}} closes no section", t.Name)
					}
					sec := stack[len(stack)-1]
					if t.Name != "" && t.Name != sec.name {
						return f.errorf(ts, rw.idx, cl.col, "{{/%s}} closes the section {{#%s}} of cell %s", t.Name, sec.name, cellRef(sec.at[0], sec.at[1]))
					}
					sec.r1 = rw.idx
					stack = stack[:len(stack)-1]
					tc.rowLevel[k] = true
				}
			}
			for _, k := range local {
				t := toks[k]
				parent := root
				if n := len(stack); n > 0 {
					parent = stack[n-1]
				}
				if len(stack) >= maxSectionDepth {
					return f.errorf(ts, rw.idx, cl.col, "the section {{#%s}} is inside more than %d sections of rows", t.Name, maxSectionDepth)
				}
				sec := &tplSection{name: t.Name, path: t.Path, inverted: t.Kind == tag.Inverted, mode: t.Mode, r0: rw.idx, r1: -1,
					at: [2]int{rw.idx, cl.col}, parent: parent}
				parent.children = append(parent.children, sec)
				stack = append(stack, sec)
				tc.rowLevel[k] = true
			}
			for _, sec := range lifted {
				sec.parent = root
				if n := len(stack); n > 0 {
					sec.parent = stack[n-1]
				}
				if len(stack) >= maxSectionDepth {
					return f.errorf(ts, rw.idx, cl.col, "the section {{#%s}} is inside more than %d sections of rows", sec.name, maxSectionDepth)
				}
				sec.parent.children = append(sec.parent.children, sec)
			}
			// (a cell of text alone had a brace written as text: \{{)
			if tagged || len(toks) != 1 || toks[0].Text != cl.text.plain {
				ts.cells[[2]int{rw.idx, cl.col}] = tc
			}
		}
	}
	if n := len(stack); n > 0 {
		sec := stack[n-1]
		return f.errorf(ts, sec.at[0], sec.at[1], "the section {{#%s}} is not closed with {{/%s}}", sec.name, sec.name)
	}
	if err := f.check(ts, root); err != nil {
		return err
	}
	f.collectFields(ts)
	return nil
}

// check gives the sections of rows their modes and refuses those that
// cannot be filled.
func (f *filler) check(ts *tplSheet, sec *tplSection) error {
	for i, ch := range sec.children {
		switch {
		case ch.inverted:
			if ch.mode == tag.ModeSlots || ch.mode == tag.ModePages {
				return f.errorf(ts, ch.at[0], ch.at[1], "the inverted section {{^%s}} does not repeat: it takes no %q", ch.name, ch.mode)
			}
			ch.mode = tag.ModeRows
		case ch.mode == "" && ch.r0 == ch.r1:
			ch.mode = tag.ModeRows
		case ch.mode == "":
			ch.mode = tag.ModeSlots
		}
		if ch.mode == tag.ModeSlots && len(ch.children) > 0 {
			in := ch.children[0]
			return f.errorf(ts, in.at[0], in.at[1], "the section {{#%s}} is in the rows of {{#%s}}, which are slots (a row for an item); write {{#%s | rows}} to repeat the rows instead",
				in.name, ch.name, ch.name)
		}
		if i > 0 && ch.r0 <= sec.children[i-1].r1 {
			prev := sec.children[i-1]
			return f.errorf(ts, ch.at[0], ch.at[1], "the sections {{#%s}} and {{#%s}} both repeat row %d; a row is repeated by one section, or by sections one inside the other",
				prev.name, ch.name, ch.r0+1)
		}
		if err := f.check(ts, ch); err != nil {
			return err
		}
	}
	return nil
}

// sectionAt returns the innermost section of rows that a row is in (the
// root for none).
func (ts *tplSheet) sectionAt(r int) *tplSection {
	sec := ts.root
	for {
		var in *tplSection
		for _, ch := range sec.children {
			if ch.r0 <= r && r <= ch.r1 {
				in = ch
				break
			}
		}
		if in == nil {
			return sec
		}
		sec = in
	}
}

// collectFields notes the names that the tags of a sheet look up.
func (f *filler) collectFields(ts *tplSheet) {
	keys := make([][2]int, 0, len(ts.cells))
	for k := range ts.cells {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i][0] < keys[j][0] || keys[i][0] == keys[j][0] && keys[i][1] < keys[j][1]
	})
	for _, k := range keys {
		tc := ts.cells[k]
		at := ts.ws.name + "!" + cellRef(k[0], k[1])
		// the sections of rows around the cell, outermost first
		var in []string
		for sec := ts.sectionAt(k[0]); sec != ts.root; sec = sec.parent {
			in = append([]string{sec.name}, in...)
		}
		rows := ts.sectionAt(k[0])
		for i, t := range tc.toks {
			switch t.Kind {
			case tag.Value:
				f.field(tag.Field{Name: t.Name, Kind: tag.FieldValue, In: in, Format: t.Format}, at)
			case tag.Image:
				f.field(tag.Field{Name: t.Name, Kind: tag.FieldImage, In: in, Code: t.Code}, at)
			case tag.Open, tag.Inverted:
				fd := tag.Field{Name: t.Name, Kind: tag.FieldSection, In: in, Inverted: t.Kind == tag.Inverted}
				if tc.rowLevel[i] {
					// the section of rows this tag opens; it is looked up
					// in the sections around it
					sec := rows
					for sec != ts.root && !(sec.at == k && sec.name == t.Name) {
						sec = sec.parent
					}
					if sec != ts.root {
						fd.Mode, fd.In = sec.mode, nil
						if sec.mode == tag.ModeSlots {
							fd.Slots = sec.r1 - sec.r0 + 1
						}
						for p := sec.parent; p != ts.root; p = p.parent {
							fd.In = append([]string{p.name}, fd.In...)
						}
					}
					f.field(fd, at)
					continue
				}
				f.field(fd, at)
				in = append(append([]string(nil), in...), t.Name)
			case tag.Close:
				if !tc.rowLevel[i] && len(in) > 0 {
					in = in[:len(in)-1]
				}
			}
		}
	}
}

// field notes a use of a name of the data (not of the variables, whose
// names start with "_", nor of the item itself).
func (f *filler) field(fd tag.Field, at string) {
	// (the variables of sections are not of the data; "." is the item of
	// the section around it, and tells that its items are values)
	if strings.HasPrefix(fd.Name, "_") || fd.Name == "." && len(fd.In) == 0 {
		return
	}
	key := fd.Kind + "\x00" + strings.Join(fd.In, "\x00") + "\x00\x00" + fd.Name
	old := f.fields[key]
	if old == nil {
		cp := fd
		cp.In = append([]string(nil), fd.In...)
		f.fields[key] = &cp
		f.fieldOrder = append(f.fieldOrder, key)
		old = &cp
	} else {
		if old.Format == "" {
			old.Format = fd.Format
		}
		if old.Mode == "" {
			old.Mode, old.Slots = fd.Mode, fd.Slots
		}
		old.Inverted = old.Inverted && fd.Inverted
	}
	if n := len(old.At); n < 16 && (n == 0 || old.At[n-1] != at) {
		old.At = append(old.At, at)
	}
}

// fieldList returns the fields of the template in the order of the sheets
// and their cells.
func (f *filler) fieldList() []tag.Field {
	out := make([]tag.Field, 0, len(f.fieldOrder))
	for _, k := range f.fieldOrder {
		out = append(out, *f.fields[k])
	}
	return out
}

// styleWith returns the format of a cell with another number format.
func (f *filler) styleWith(style int, code string, id int) int {
	key := styleKey{style, code}
	if code == "" {
		key.code = fmt.Sprintf("\x00%d", id)
	}
	if i, ok := f.styles[key]; ok {
		return i
	}
	st := f.c.st
	if code != "" {
		id = f.nextFmt
		f.nextFmt++
		if st.numFmts == nil {
			st.numFmts = map[int]string{}
		}
		st.numFmts[id] = code
	}
	x := *st.xf(style)
	x.numFmt = id
	st.cellXfs = append(st.cellXfs, x)
	i := len(st.cellXfs) - 1
	f.styles[key] = i
	return i
}

// timeSerial is the serial date of a time (of its clock, whatever its
// zone).
func timeSerial(t time.Time, date1904 bool) float64 {
	base := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	if date1904 {
		base = time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	wall := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
	sec := wall.Unix() - base.Unix()
	d := float64(sec/86400) + (float64(sec%86400)+float64(wall.Nanosecond())/1e9)/86400
	if !date1904 && d < 61 {
		d-- // before the fictitious 1900-02-29
	}
	return d
}

// isDateFormat reports whether a format of cells shows dates.
func (f *filler) isDateFormat(style int) bool {
	st := f.c.st
	return st.numFormat(st.xf(style).numFmt).isDate()
}
