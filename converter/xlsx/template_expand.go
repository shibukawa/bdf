package xlsx

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif" // the pictures that data names
	_ "image/jpeg"
	_ "image/png"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/shibukawa/bdf/form/tag"
	"github.com/shibukawa/bdf/image/imgconv"
)

// maxCachedFormats bounds the number formats a Formatter keeps parsed.
const maxCachedFormats = 4096

// expander makes the rows of a filled sheet from those of its template.
type expander struct {
	f   *filler
	ts  *tplSheet
	src []row // the rows of the template
	out []row
}

// srcCell is a cell of the template that a cell of the filled sheet is
// made from.
type srcCell struct {
	cl cell
	tc *tplCell
	cf *cellFormula
	dr int // rows from the row the cell was read in (a slot filled from the first)
	// row and col are where the cell is in the template
	row, col int
}

// expand fills the rows of a sheet.
func (f *filler) expand(ts *tplSheet) error {
	x := &expander{f: f, ts: ts, src: ts.ws.rows}
	ts.inst = &itemInst{}
	last := 0
	if n := len(x.src); n > 0 {
		last = x.src[n-1].idx
	}
	if _, err := x.unit(ts.root.children, 0, last, f.root, ts.inst, 0); err != nil {
		return err
	}
	ts.inst.dst1 = maxRows
	ts.ws.rows = x.out
	return nil
}

// rowAt returns a row of the template, or nil.
func (x *expander) rowAt(idx int) *row {
	i := sort.Search(len(x.src), func(i int) bool { return x.src[i].idx >= idx })
	if i < len(x.src) && x.src[i].idx == idx {
		return &x.src[i]
	}
	return nil
}

// vars are the variables of an item of a section.
func vars(k, n int, pages bool) map[string]any { return tag.Vars(k, n, pages) }

// Formatter returns a formatter of values in the number formats of Excel,
// for the templates that are not workbooks (see tag.Formatter): numbers and
// dates in a format code, text in the text section of one. A text that is
// a date is read as one for a date format. The formats are parsed once
// each; the formatter may be called from several goroutines.
func Formatter() tag.Formatter {
	var mu sync.Mutex
	parsed := map[string]*numFormat{}
	parse := func(format string) *numFormat {
		mu.Lock()
		defer mu.Unlock()
		if nf, ok := parsed[format]; ok {
			return nf
		}
		nf := parseNumFormat(format, nil)
		if len(parsed) < maxCachedFormats {
			parsed[format] = nf
		}
		return nf
	}
	return func(v any, format string) string {
		if format == "" {
			if x, ok := v.(float64); ok && !math.IsNaN(x) && !math.IsInf(x, 0) {
				return formatGeneral(x, 11)
			}
			return tag.ValueText(v, "", nil)
		}
		nf := parse(format)
		if s, ok := v.(string); ok && nf.isDate() {
			if d, ok := tag.ParseDate(s); ok {
				v = d
			}
		}
		switch x := v.(type) {
		case float64:
			if math.IsNaN(x) || math.IsInf(x, 0) {
				return ""
			}
			return nf.formatNumber(x, false).text()
		case time.Time:
			return nf.formatNumber(timeSerial(x, false), false).text()
		case string:
			if nf.hasTextSection() {
				return nf.formatText(x).text()
			}
			return x
		}
		return tag.ValueText(v, "", nil)
	}
}

// unit fills the rows r0 to r1 of the template, a copy of the rows of a
// section (or the sheet), from the row dst of the filled sheet: the rows of
// the sections in it are repeated for their items, the others are filled
// as they are. It returns the row after the last.
func (x *expander) unit(children []*tplSection, r0, r1 int, sc *tag.Scope, inst *itemInst, dst int) (int, error) {
	inst.dst0 = dst
	t := r0
	for _, ch := range children {
		var err error
		if dst, err = x.static(t, ch.r0-1, sc, inst, dst); err != nil {
			return 0, err
		}
		t = ch.r1 + 1
		v, found := sc.Lookup(ch.path)
		if !found {
			if err := x.f.missing(x.ts, ch.at[0], ch.at[1], ch.name, sc); err != nil {
				return 0, err
			}
		}
		if ch.mode == tag.ModeSlots {
			items, n := tag.Items(v), ch.r1-ch.r0+1
			if len(items) > n {
				return 0, x.f.errorf(x.ts, ch.at[0], ch.at[1], "%d items for the %d rows of {{#%s}}: the rows from this one to that of {{/%s}} are slots, a row for an item (write {{#%s | rows}} to repeat the rows instead)",
					len(items), n, ch.name, ch.name, ch.name)
			}
			pattern := x.cellsOf(x.rowAt(ch.r0), ch.r0)
			for k := 0; k < n; k++ {
				if err := x.f.expansionStep(x.ts, ch.at[0], ch.at[1]); err != nil {
					return 0, err
				}
				isc := sc.Blank(map[string]any{"_blank": true, "_blank_first": k == len(items), "_index": float64(k + 1), "_count": float64(len(items))})
				if k < len(items) {
					v := vars(k, len(items), false)
					v["_blank"], v["_blank_first"] = false, false
					isc = sc.Item(items[k], v)
				}
				if err := x.slot(ch, k, pattern, isc, inst, dst+k); err != nil {
					return 0, err
				}
			}
			dst += n
			continue
		}
		si := &secInst{sec: ch, dst0: dst}
		inst.children = append(inst.children, si)
		var items []any
		switch {
		case ch.inverted && !tag.Truthy(v):
			items = []any{nil} // once, in the scope around it
		case !ch.inverted:
			items = tag.Items(v)
		}
		for k, it := range items {
			if err := x.f.expansionStep(x.ts, ch.at[0], ch.at[1]); err != nil {
				return 0, err
			}
			ii := &itemInst{parent: inst, of: si}
			si.items = append(si.items, ii)
			ch.insts = append(ch.insts, ii)
			isc := sc
			if !ch.inverted {
				isc = sc.Item(it, vars(k, len(items), ch.mode == tag.ModePages))
			}
			if ch.mode == tag.ModePages && k > 0 {
				x.ts.breaks = append(x.ts.breaks, dst)
			}
			if dst, err = x.unit(ch.children, ch.r0, ch.r1, isc, ii, dst); err != nil {
				return 0, err
			}
			ii.dst1 = dst
		}
		si.dst1 = dst
	}
	dst, err := x.static(t, r1, sc, inst, dst)
	if err != nil {
		return 0, err
	}
	inst.dst1 = dst
	return dst, nil
}

// static fills the rows t0 to t1 of the template, which no section in the
// unit repeats, from the row dst.
func (x *expander) static(t0, t1 int, sc *tag.Scope, inst *itemInst, dst int) (int, error) {
	if t1 < t0 {
		return dst, nil
	}
	i := sort.Search(len(x.src), func(i int) bool { return x.src[i].idx >= t0 })
	for ; i < len(x.src) && x.src[i].idx <= t1; i++ {
		rw := &x.src[i]
		if err := x.emit(rw, x.cellsOf(rw, rw.idx), dst+rw.idx-t0, sc, inst); err != nil {
			return 0, err
		}
	}
	return dst + (t1 - t0 + 1), nil
}

// cellsOf returns the cells of a row of the template with their tags and
// formulas.
func (x *expander) cellsOf(rw *row, idx int) []srcCell {
	if rw == nil {
		return nil
	}
	out := make([]srcCell, len(rw.cells))
	for i, cl := range rw.cells {
		k := [2]int{idx, cl.col}
		out[i] = srcCell{cl: cl, tc: x.ts.cells[k], cf: x.ts.ws.formulas[k], row: idx, col: cl.col}
	}
	return out
}

// slot fills the slot k of a section of slots: the row of the template at
// that place, with the tags and formulas of the first row of the section
// where it has no tags of its own.
func (x *expander) slot(sec *tplSection, k int, pattern []srcCell, sc *tag.Scope, inst *itemInst, dst int) error {
	idx := sec.r0 + k
	rw := x.rowAt(idx)
	cells := x.cellsOf(rw, idx)
	// a cell with tags of its own: not one that only closes the section
	tagged := func(c srcCell) bool {
		if c.tc == nil {
			return false
		}
		for i, t := range c.tc.toks {
			if !c.tc.rowLevel[i] && (t.Kind != tag.Text || strings.TrimSpace(t.Text) != "") {
				return true
			}
		}
		return false
	}
	own := false
	for _, c := range cells {
		own = own || tagged(c)
	}
	if k > 0 && !own {
		// the values and formulas of the first row, in the formats of this one
		byCol := map[int]int{}
		for i, c := range cells {
			byCol[c.col] = i
		}
		for _, p := range pattern {
			if p.tc == nil && p.cf == nil {
				continue
			}
			p.dr = k
			if i, ok := byCol[p.col]; ok {
				if (cells[i].cl.kind != cellBlank && cells[i].tc == nil) || cells[i].cf != nil {
					continue // the row says what is here
				}
				p.cl.style = cells[i].cl.style
				cells[i] = p
				continue
			}
			cells = append(cells, p)
		}
		sort.SliceStable(cells, func(i, j int) bool { return cells[i].col < cells[j].col })
	}
	if rw == nil {
		if len(cells) == 0 {
			return nil
		}
		rw = &row{ht: -1, style: -1}
	}
	return x.emit(rw, cells, dst, sc, inst)
}

// emit adds a row of the filled sheet: a row of the template with its cells
// filled.
func (x *expander) emit(src *row, cells []srcCell, dst int, sc *tag.Scope, inst *itemInst) error {
	if dst >= maxRows {
		return fmt.Errorf("xlsx: sheet %q: the filled sheet has more than the %d rows of a sheet", x.ts.ws.name, maxRows)
	}
	if x.f.cells += len(cells) + 1; x.f.cells > maxFilledCells {
		return fmt.Errorf("xlsx: sheet %q: the filled template has more than %d cells", x.ts.ws.name, maxFilledCells)
	}
	out := *src
	out.idx = dst
	out.cells = make([]cell, 0, len(cells))
	for _, c := range cells {
		cl, err := x.cell(c, dst, sc)
		if err != nil {
			return err
		}
		out.cells = append(out.cells, cl)
		if c.cf != nil && c.tc == nil {
			x.ts.outF[[2]int{dst, c.col}] = &outFormula{f: c.cf, dr: c.dr, ctx: inst}
		}
	}
	x.out = append(x.out, out)
	return nil
}

// missing tells that the data has no value for a name: a warning, or an
// error when the template is filled strictly. A place that no item fills
// has no values.
func (f *filler) missing(ts *tplSheet, r, c int, name string, sc *tag.Scope) error {
	if sc.IsBlank() {
		return nil
	}
	if f.strict {
		return f.errorf(ts, r, c, "the data has no %q", name)
	}
	f.c.warnOnce("field:"+name, "sheet %q cell %s: the data has no %q; it is left empty", ts.ws.name, cellRef(r, c), name)
	return nil
}

// cell fills a cell: a cell without tags is as it is; a cell that is one
// tag takes the value; in other cells the tags are replaced in the text.
func (x *expander) cell(c srcCell, dst int, sc *tag.Scope) (cell, error) {
	cl := c.cl
	if c.tc == nil {
		return cl, nil
	}
	toks, skip := c.tc.toks, c.tc.rowLevel
	// the one tag of the cell, when it is nothing else
	one, n := -1, 0
	for i, t := range toks {
		if skip[i] || t.Kind == tag.Text && strings.TrimSpace(t.Text) == "" {
			continue
		}
		one = i
		n++
	}
	switch {
	case n == 0:
		cl.kind, cl.text, cl.num = cellBlank, nil, 0
		return cl, nil
	case n == 1 && toks[one].Kind == tag.Value:
		t := toks[one]
		v, ok := sc.Lookup(t.Path)
		if !ok {
			if err := x.f.missing(x.ts, c.row, c.col, t.Name, sc); err != nil {
				return cl, err
			}
		}
		x.f.setValue(&cl, v, t, x.ts, c)
		return cl, nil
	}
	// text with tags: the pieces, each with the place in the text of what it stands for
	type piece struct {
		s   string
		pos int
	}
	var pieces []piece
	total := 0
	emit := func(s string, pos int) {
		if s == "" || total >= maxCellText {
			return
		}
		text, n := textPrefix(s, maxCellText-total)
		total += n
		if len(text) < len(s) {
			x.f.c.warnOnce("celltext:"+x.ts.ws.name, "sheet %q cell %s: the text is cut at the %d characters of a cell", x.ts.ws.name, cellRef(c.row, c.col), maxCellText)
		}
		pieces = append(pieces, piece{text, pos})
	}
	var walk func(lo, hi int, sc *tag.Scope) error
	walk = func(lo, hi int, sc *tag.Scope) error {
		for i := lo; i < hi; i++ {
			if err := x.f.expansionStep(x.ts, c.row, c.col); err != nil {
				return err
			}
			if skip[i] {
				continue
			}
			t := toks[i]
			switch t.Kind {
			case tag.Text:
				emit(t.Text, t.Pos)
			case tag.Value:
				v, ok := sc.Lookup(t.Path)
				if !ok {
					if err := x.f.missing(x.ts, c.row, c.col, t.Name, sc); err != nil {
						return err
					}
				}
				emit(x.f.valueText(v, t, cl.style), t.Pos)
			case tag.Image:
				if err := x.picture(t, sc, c, dst); err != nil {
					return err
				}
			case tag.Open, tag.Inverted:
				// its closing tag, in the cell
				j, depth := i+1, 1
				for ; j < hi; j++ {
					if skip[j] {
						continue
					}
					if k := toks[j].Kind; k == tag.Open || k == tag.Inverted {
						depth++
					} else if k == tag.Close {
						if depth--; depth == 0 {
							break
						}
					}
				}
				v, ok := sc.Lookup(t.Path)
				if !ok {
					if err := x.f.missing(x.ts, c.row, c.col, t.Name, sc); err != nil {
						return err
					}
				}
				if t.Kind == tag.Inverted {
					if !tag.Truthy(v) {
						if err := walk(i+1, j, sc); err != nil {
							return err
						}
					}
				} else {
					items := tag.Items(v)
					for k, it := range items {
						if total >= maxCellText {
							break
						}
						if err := x.f.expansionStep(x.ts, c.row, c.col); err != nil {
							return err
						}
						if err := walk(i+1, j, sc.Item(it, vars(k, len(items), false))); err != nil {
							return err
						}
					}
				}
				i = j
			}
		}
		return nil
	}
	if err := walk(0, len(toks), sc); err != nil {
		return cl, err
	}
	var b strings.Builder
	for _, p := range pieces {
		b.WriteString(p.s)
	}
	if b.Len() == 0 {
		cl.kind, cl.text, cl.num = cellBlank, nil, 0
		return cl, nil
	}
	rt := &richText{plain: b.String()}
	if old := cl.text; old != nil && len(old.runs) > 0 {
		// each piece in the font of the run it stands in
		ends := make([]int, len(old.runs))
		end := 0
		for i, run := range old.runs {
			end += len(run.text)
			ends[i] = end
		}
		fontAt := func(pos int) *xfont {
			return old.runs[min(sort.SearchInts(ends, pos+1), len(old.runs)-1)].font
		}
		var text strings.Builder
		var font *xfont
		for _, p := range pieces {
			next := fontAt(p.pos)
			if text.Len() > 0 && next != font {
				rt.runs = append(rt.runs, textRun{text: text.String(), font: font})
				text = strings.Builder{}
			}
			font = next
			text.WriteString(p.s)
		}
		if text.Len() > 0 {
			rt.runs = append(rt.runs, textRun{text: text.String(), font: font})
		}
	}
	cl.kind, cl.text, cl.num = cellStr, rt, 0
	return cl, nil
}

// setValue makes a value the value of a cell that is one tag.
func (f *filler) setValue(cl *cell, v any, t tag.Token, ts *tplSheet, c srcCell) {
	cl.kind, cl.text, cl.num = cellBlank, nil, 0
	// a text that is a date, where a date is asked for
	if s, ok := v.(string); ok && (f.isDateFormat(cl.style) || t.HasFormat && f.numFmt(t.Format).isDate()) {
		if d, ok := tag.ParseDate(s); ok {
			v = d
		}
	}
	switch x := v.(type) {
	case nil:
	case float64:
		if math.IsNaN(x) || math.IsInf(x, 0) {
			return
		}
		cl.kind, cl.num = cellNum, x
	case bool:
		cl.kind = cellBool
		if x {
			cl.num = 1
		}
	case string:
		if x != "" {
			cl.kind, cl.text = cellStr, &richText{plain: clipText(x)}
		}
	case time.Time:
		cl.kind, cl.num = cellNum, timeSerial(x, f.c.date1904)
		if !t.HasFormat && !f.isDateFormat(cl.style) {
			// as Excel formats a date typed into a cell
			id := 14
			if x.Hour() != 0 || x.Minute() != 0 || x.Second() != 0 {
				id = 22
			}
			cl.style = f.styleWith(cl.style, "", id)
		}
	default:
		f.c.warnOnce("fieldkind:"+t.Name, "sheet %q cell %s: %q is a list or an object, not a value; it is left empty", ts.ws.name, cellRef(c.row, c.col), t.Name)
		return
	}
	if t.HasFormat && t.Format != "" && cl.kind != cellBlank {
		cl.style = f.styleWith(cl.style, t.Format, 0)
	}
}

// clipText cuts a text at the characters of a cell.
func clipText(s string) string {
	if len(s) <= maxCellText {
		return s
	}
	text, _ := textPrefix(s, maxCellText)
	return text
}

// textPrefix reads at most n+1 runes, retaining a slice of the source
// rather than allocating for all of an oversized data value.
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

// valueText formats a value for the text of a cell: in the format of the
// tag, or as Excel shows it in a cell without a format.
func (f *filler) valueText(v any, t tag.Token, style int) string {
	st := f.c.st
	var nf *numFormat
	if t.HasFormat && t.Format != "" {
		nf = f.numFmt(t.Format)
		if s, ok := v.(string); ok && nf.isDate() {
			if d, ok := tag.ParseDate(s); ok {
				v = d
			}
		}
	}
	switch x := v.(type) {
	case float64:
		if math.IsNaN(x) || math.IsInf(x, 0) {
			return ""
		}
		if nf != nil {
			return nf.formatNumber(x, f.c.date1904).text()
		}
		return formatGeneral(x, 11)
	case bool:
		if x {
			return "TRUE"
		}
		return "FALSE"
	case string:
		if nf != nil && nf.hasTextSection() {
			return nf.formatText(x).text()
		}
		return x
	case time.Time:
		if nf == nil {
			id := 14
			if x.Hour() != 0 || x.Minute() != 0 || x.Second() != 0 {
				id = 22
			}
			nf = st.numFormat(id)
		}
		return nf.formatNumber(timeSerial(x, f.c.date1904), f.c.date1904).text()
	}
	return ""
}

// picture puts the picture that a tag names over the cell: the path of a
// file of the template's files, or a data: URL.
func (x *expander) picture(t tag.Token, sc *tag.Scope, c srcCell, dst int) error {
	f := x.f
	v, ok := sc.Lookup(t.Path)
	if !ok {
		return f.missing(x.ts, c.row, c.col, t.Name, sc)
	}
	if t.Code != "" {
		// a QR code or a barcode of the value
		if v == nil || v == "" {
			return nil
		}
		sym, err := t.Symbol(v)
		if err != nil {
			if f.strict {
				return f.errorf(x.ts, c.row, c.col, "%v", err)
			}
			f.c.warnf("sheet %q cell %s: the code of %q is left out: %v", x.ts.ws.name, cellRef(c.row, c.col), t.Name, err)
			return nil
		}
		x.ts.pictures = append(x.ts.pictures, tplPicture{row: dst, col: c.col, fit: t.Fit, alt: sym.Text, sym: sym, w: float64(sym.W), h: float64(sym.H)})
		return nil
	}
	name, _ := v.(string)
	if name == "" {
		return nil
	}
	pic := f.pictures[name]
	if pic == nil {
		pic = &tplPicture{}
		f.pictures[name] = pic
		data, err := f.pictureData(name)
		if err != nil {
			if f.strict {
				return f.errorf(x.ts, c.row, c.col, "the picture of %q: %v", t.Name, err)
			}
			f.c.warnf("sheet %q cell %s: the picture of %q is left out: %v", x.ts.ws.name, cellRef(c.row, c.col), t.Name, err)
			return nil
		}
		cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
		if format := imgconv.Sniff(data); err != nil || format == "" || format == "svg" || cfg.Width <= 0 || cfg.Height <= 0 {
			if f.strict {
				return f.errorf(x.ts, c.row, c.col, "the picture of %q is not a PNG, JPEG or GIF image", t.Name)
			}
			f.c.warnf("sheet %q cell %s: the picture of %q is not a PNG, JPEG or GIF image; it is left out", x.ts.ws.name, cellRef(c.row, c.col), t.Name)
			return nil
		}
		if res, err := imgconv.Optimize(data, f.c.opts.Images); err == nil || res.Data != nil {
			data = res.Data
		}
		pic.hash, pic.w, pic.h = f.c.doc.AddImage(data), float64(cfg.Width), float64(cfg.Height)
	}
	if pic.w <= 0 {
		return nil // left out, as told before
	}
	p := *pic
	p.row, p.col, p.fit, p.alt = dst, c.col, t.Fit, t.Name
	x.ts.pictures = append(x.ts.pictures, p)
	return nil
}

// pictureData reads a picture that the data names.
func (f *filler) pictureData(name string) ([]byte, error) { return tag.PictureData(f.files, name) }

// mapRow returns the row of the filled sheet that a row of the template
// became, as a formula in the copy ctx of the rows of a section sees it:
// the row of the same copy; for a row in a section the formula is not in,
// its first copy, or its last when the row ends a range. gone is set when
// the row is in a section that no item fills: the row is then the one
// after where it would be (before, at the end of a range).
func (ts *tplSheet) mapRow(t int, ctx *itemInst, end bool) (dst int, gone bool) {
	inst, r0 := ts.inst, 0
	if inst == nil {
		return t, false
	}
	for {
		grow := 0
		var into *secInst
		for _, si := range inst.children {
			if si.sec.r1 < t {
				grow += si.dst1 - si.dst0 - (si.sec.r1 - si.sec.r0 + 1)
				continue
			}
			if si.sec.r0 <= t {
				into = si
			}
			break
		}
		if into == nil {
			return inst.dst0 + t - r0 + grow, false
		}
		if len(into.items) == 0 {
			if end {
				return into.dst0 - 1, true
			}
			return into.dst0, true
		}
		pick := into.items[0]
		if end {
			pick = into.items[len(into.items)-1]
		}
		for c := ctx; c != nil; c = c.parent {
			if c.of == into {
				pick = c
				break
			}
		}
		inst, r0 = pick, into.sec.r0
	}
}

// eachCopy calls fn with the rows of the filled sheet that the rows r0 to
// r1 of the template became: once for each copy of the rows of the
// innermost section that holds them all, with the copies of the sections
// between them held whole.
func (ts *tplSheet) eachCopy(r0, r1 int, fn func(d0, d1 int)) {
	if ts.inst == nil {
		fn(r0, r1)
		return
	}
	// the innermost section that repeats rows and holds r0 to r1
	sec := ts.root
	for {
		var in *tplSection
		for _, ch := range sec.children {
			if ch.mode != tag.ModeSlots && ch.r0 <= r0 && r1 <= ch.r1 {
				in = ch
				break
			}
		}
		if in == nil {
			break
		}
		sec = in
	}
	one := func(ctx *itemInst) {
		d0, g0 := ts.mapRow(r0, ctx, false)
		d1, g1 := ts.mapRow(r1, ctx, true)
		if g0 && g1 || d1 < d0 || d0 >= maxRows {
			return
		}
		fn(d0, min(d1, maxRows-1))
	}
	if sec == ts.root {
		one(nil)
		return
	}
	for i, ii := range sec.insts {
		if i >= maxCopies {
			break
		}
		one(ii)
	}
}
