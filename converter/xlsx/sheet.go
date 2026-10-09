package xlsx

import (
	"bytes"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/tinygodriver/encoding/xmlro"
)

// A worksheet (ECMA-376 Part 1 §18.3) is read as a stream: sheetData, which
// holds the cells and can be very large, into compact rows; the other
// elements (columns, merged cells, conditional formats, hyperlinks, views)
// into node trees.

type cellKind byte

const (
	cellBlank cellKind = iota // formatted but empty
	cellNum
	cellStr
	cellBool
	cellErr
)

type cell struct {
	col   int
	style int
	kind  cellKind
	num   float64
	text  *richText // strings and error values
}

// richText is a string, with runs of their own fonts when rich.
type richText struct {
	plain string
	runs  []textRun // nil for plain text
}

type textRun struct {
	text string
	font *xfont // nil: the cell's font
}

type row struct {
	idx    int
	ht     float64 // points; < 0 when not set
	custom bool
	hidden bool
	style  int // row format (customFormat), -1 for none
	cells  []cell
}

type colDef struct {
	min, max int     // 0-based, inclusive
	width    float64 // characters; < 0 when not set
	hidden   bool
	style    int
}

// cellRange is a rectangle of cells, 0-based and inclusive.
type cellRange struct{ r0, c0, r1, c1 int }

func (r cellRange) contains(row, col int) bool {
	return row >= r.r0 && row <= r.r1 && col >= r.c0 && col <= r.c1
}

type worksheet struct {
	name     string
	part     string
	rows     []row
	cols     []colDef
	defColW  float64 // characters; 0 when not set
	baseColW float64
	defRowH  float64 // points; 0 when not set
	zeroH    bool    // rows are hidden unless they say otherwise
	merges   []cellRange
	showGrid bool
	showZero bool
	rtl      bool
	freezeC  int
	freezeR  int
	cfs      []*ooxml.Node // conditionalFormatting
	x14cfs   []*ooxml.Node // x14:conditionalFormatting in the extension list
	links    []*ooxml.Node
	drawing  string // relationship ids
	legacy   string
	tables   []string
	// autoFilter is the range of the sheet's AutoFilter (tables have their own)
	autoFilter *cellRange
	noValue    int // formula cells saved without a value

	// grids (ConvertGrid)
	fitCols    bool // columns are as wide as their text
	headerRows int  // rows at the top that head the columns
	gridTables []gridTable

	// the col elements by column (indexCols)
	colsIndexed      bool
	colLast, colBase []int32
}

// readWorksheet parses a worksheet part: sheetData as a stream (readSheetData),
// the other elements as trees.
func (c *converter) readWorksheet(part string) (*worksheet, error) {
	rc, err := c.pkg.OpenPart(part)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	ws := &worksheet{part: part, baseColW: 8, showGrid: true, showZero: true}
	r := ooxml.NewReader(rc)
	root, err := rootElement(r)
	if err != nil {
		return nil, err
	}
	for {
		ok, err := r.NextChild(root)
		if err != nil {
			return nil, err
		}
		if !ok {
			return ws, nil
		}
		if xmlro.Equal(r.LocalName(), "sheetData") {
			if err := c.readSheetData(r, ws); err != nil {
				return nil, err
			}
			continue
		}
		n, err := c.pkg.ReadFrom(r)
		if err != nil {
			return nil, err
		}
		ws.element(n)
	}
}

// rootElement advances a reader to the root element of a part.
func rootElement(r *xmlro.Reader) (xmlro.Element, error) {
	for {
		k, err := r.Next()
		if err != nil {
			return xmlro.Element{}, err
		}
		switch k {
		case xmlro.StartElement:
			return r.Element(), nil
		case xmlro.EOF:
			return xmlro.Element{}, io.ErrUnexpectedEOF
		}
	}
}

// element takes what the renderer needs from a top-level element.
func (ws *worksheet) element(n *ooxml.Node) {
	switch n.Name {
	case "AlternateContent":
		for _, k := range n.Elements() {
			ws.element(k)
		}
	case "sheetViews":
		v := n.Child("sheetView")
		ws.showGrid = v.AttrBool("showGridLines", true)
		ws.showZero = v.AttrBool("showZeros", true)
		ws.rtl = v.AttrBool("rightToLeft", false)
		if p := v.Child("pane"); p != nil && strings.HasPrefix(p.AttrStr("state", "split"), "frozen") {
			ws.freezeC = int(p.AttrFloat("xSplit", 0))
			ws.freezeR = int(p.AttrFloat("ySplit", 0))
		}
	case "sheetFormatPr":
		ws.defColW = clampColWidth(n.AttrFloat("defaultColWidth", 0))
		ws.baseColW = clampColWidth(n.AttrFloat("baseColWidth", 8))
		if h := clampRowHt(n.AttrFloat("defaultRowHeight", 0)); h >= 0 {
			ws.defRowH = h
		}
		ws.zeroH = n.AttrBool("zeroHeight", false)
	case "cols":
		for _, k := range n.Children("col") {
			cd := colDef{min: int(k.AttrInt("min", 1)) - 1, max: int(k.AttrInt("max", 1)) - 1, width: -1,
				hidden: k.AttrBool("hidden", false), style: int(k.AttrInt("style", -1))}
			if _, ok := k.Attr("width"); ok {
				cd.width = clampColWidth(k.AttrFloat("width", 0))
			}
			if cd.min < 0 || cd.max < cd.min {
				continue
			}
			cd.max = min(cd.max, maxCols-1)
			ws.cols = append(ws.cols, cd)
		}
	case "mergeCells":
		for _, k := range n.Children("mergeCell") {
			if r, ok := parseRange(k.AttrStr("ref", "")); ok && (r.r1 > r.r0 || r.c1 > r.c0) {
				ws.merges = append(ws.merges, r)
			}
		}
	case "conditionalFormatting":
		ws.cfs = append(ws.cfs, n)
	case "hyperlinks":
		ws.links = append(ws.links, n.Children("hyperlink")...)
	case "drawing":
		ws.drawing = n.RelID("id")
	case "legacyDrawing":
		ws.legacy = n.RelID("id")
	case "autoFilter":
		if r, ok := parseRange(n.AttrStr("ref", "")); ok {
			ws.autoFilter = &r
		}
	case "tableParts":
		for _, k := range n.Children("tablePart") {
			ws.tables = append(ws.tables, k.RelID("id"))
		}
	case "extLst":
		for _, ext := range n.Children("ext") {
			if cfs := ext.Child("conditionalFormattings"); cfs != nil {
				ws.x14cfs = append(ws.x14cfs, cfs.Children("conditionalFormatting")...)
			}
		}
	}
}

const (
	maxRows = 1048576
	maxCols = 16384
	// Excel's own maxima for a row's height (points) and a column's width
	// (characters, its widest digit). Heights and widths past these, or not
	// finite (a "ht" or "width" of "Inf"), are clamped, so a hostile file of
	// a few bytes cannot make the sheet's grid span an unbounded distance.
	maxRowPt   = 409.5
	maxColChar = 255
)

// clampRowHt limits a row height in points: a negative or non-finite value
// is dropped (returns -1, "not set"), the rest is capped at maxRowPt.
func clampRowHt(f float64) float64 {
	if math.IsNaN(f) || math.IsInf(f, 0) || f < 0 {
		return -1
	}
	return math.Min(f, maxRowPt)
}

// clampColWidth limits a column width in characters: a non-finite value
// becomes the maximum, a negative one zero, and the rest is capped.
func clampColWidth(f float64) float64 {
	switch {
	case math.IsNaN(f) || f < 0:
		return 0
	case math.IsInf(f, 0):
		return maxColChar
	}
	return math.Min(f, maxColChar)
}

// readSheetData reads the rows and cells of sheetData, which the reader is
// on, into compact rows: a token at a time, keeping as trees only the
// inline strings. The reader is left on the end of sheetData.
func (c *converter) readSheetData(r *xmlro.Reader, ws *worksheet) error {
	nextRow := 0
	sheetData := r.Element()
	for {
		ok, err := r.NextChild(sheetData)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		if !xmlro.Equal(r.LocalName(), "row") {
			continue
		}
		rw := row{idx: nextRow, ht: -1, style: -1}
		if v, ok := r.Attr("r"); ok {
			if n, err := v.Int(); err == nil && n >= 1 {
				rw.idx = int(min(n, maxRows)) - 1
			}
		}
		if v, ok := r.Attr("ht"); ok {
			if f, err := v.Float(); err == nil {
				rw.ht = clampRowHt(f)
			}
		}
		rw.custom = boolAttr(r, "customHeight")
		rw.hidden = boolAttr(r, "hidden")
		if v, ok := r.Attr("s"); ok {
			if n, err := v.Int(); err == nil && n >= math.MinInt32 && n <= math.MaxInt32 {
				rw.style = int(n)
			}
		}
		// the row's format applies to its empty cells only with customFormat
		if !boolAttr(r, "customFormat") {
			rw.style = -1
		}
		if rw.idx >= maxRows {
			rw.idx = maxRows - 1
		}
		ws.rows = append(ws.rows, rw)
		cur := &ws.rows[len(ws.rows)-1]
		nextRow = rw.idx + 1
		nextCol := 0
		rowEl := r.Element()
		for {
			ok, err := r.NextChild(rowEl)
			if err != nil {
				return err
			}
			if !ok {
				break
			}
			if !xmlro.Equal(r.LocalName(), "c") {
				continue
			}
			cc := cell{col: nextCol}
			if v, ok := r.Attr("r"); ok {
				if col, _, ok := parseRef(v); ok {
					cc.col = col
				}
			}
			if v, ok := r.Attr("s"); ok {
				if n, err := v.Int(); err == nil && n >= math.MinInt32 && n <= math.MaxInt32 {
					cc.style = int(n)
				}
			}
			typ := cellTypeOf(r)
			nextCol = cc.col + 1
			cur.cells = append(cur.cells, cc)
			cl := &cur.cells[len(cur.cells)-1]
			var value string
			var is *ooxml.Node
			hasF, hasV := false, false
			cellEl := r.Element()
			for {
				ok, err := r.NextChild(cellEl)
				if err != nil {
					return err
				}
				if !ok {
					break
				}
				switch {
				case xmlro.Equal(r.LocalName(), "v"):
					v, err := r.ElementText()
					if err != nil {
						return err
					}
					value, hasV = v.String(), true
				case xmlro.Equal(r.LocalName(), "is"):
					if is, err = c.pkg.ReadFrom(r); err != nil {
						return err
					}
					hasV = true
				case xmlro.Equal(r.LocalName(), "f"):
					hasF = true // its content, and extensions, are passed over
				}
			}
			c.cellValue(cl, typ, value, is)
			if hasF && (!hasV || strings.TrimSpace(value) == "" && is == nil) {
				ws.noValue++
			}
		}
	}
}

// cellTypes are the values of a cell's t attribute, each one string.
var cellTypes = []string{"n", "s", "str", "inlineStr", "b", "e", "d"}

// cellTypeOf returns the type of the cell the reader is on ("n" without
// a t attribute), without making a string of a known one.
func cellTypeOf(r *xmlro.Reader) string {
	v, ok := r.Attr("t")
	if !ok {
		return "n"
	}
	for _, t := range cellTypes {
		if v.Equal(t) {
			return t
		}
	}
	return v.String()
}

// boolAttr reads a boolean attribute of the element the reader is on, as
// xmlBool reads one; false when there is none.
func boolAttr(r *xmlro.Reader, name string) bool {
	v, ok := r.Attr(name)
	if !ok {
		return false
	}
	if v.HasEntities() {
		return xmlBool(v.String())
	}
	b := bytes.TrimSpace(v)
	return xmlro.Equal(b, "1") || xmlro.Equal(b, "true") || xmlro.Equal(b, "on")
}

func xmlBool(v string) bool {
	switch strings.TrimSpace(v) {
	case "1", "true", "on":
		return true
	}
	return false
}

// cellValue sets a cell's value from its type and its v or is element.
func (c *converter) cellValue(cl *cell, typ, v string, is *ooxml.Node) {
	switch typ {
	case "s":
		i, err := strconv.Atoi(strings.TrimSpace(v))
		if err == nil && i >= 0 && i < len(c.sst) {
			cl.kind, cl.text = cellStr, c.sst[i]
		}
	case "inlineStr":
		if is != nil {
			cl.kind, cl.text = cellStr, c.readRich(is)
		} else if v != "" {
			cl.kind, cl.text = cellStr, &richText{plain: v}
		}
	case "str":
		cl.kind, cl.text = cellStr, &richText{plain: decodeXString(v)}
	case "b":
		cl.kind = cellBool
		if strings.TrimSpace(v) == "1" || strings.EqualFold(strings.TrimSpace(v), "true") {
			cl.num = 1
		}
	case "e":
		cl.kind, cl.text = cellErr, &richText{plain: strings.TrimSpace(v)}
	case "d":
		if f, ok := isoSerial(strings.TrimSpace(v), c.date1904); ok {
			cl.kind, cl.num = cellNum, f
		}
	default:
		if s := strings.TrimSpace(v); s != "" {
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				cl.kind, cl.num = cellNum, f
			}
		}
	}
}

// isoSerial converts an ISO 8601 date (t="d" cells) to a serial date.
func isoSerial(s string, date1904 bool) (float64, bool) {
	for _, layout := range []string{"2006-01-02T15:04:05.999999999", "2006-01-02T15:04:05Z07:00", "2006-01-02", "15:04:05"} {
		t, err := time.Parse(layout, s)
		if err != nil {
			continue
		}
		if layout == "15:04:05" {
			return float64(t.Hour()*3600+t.Minute()*60+t.Second()) / 86400, true
		}
		base := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
		if date1904 {
			base = time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
		}
		d := t.Sub(base).Hours() / 24
		if !date1904 && d < 61 {
			d-- // before the fictitious 1900-02-29
		}
		return d, true
	}
	return 0, false
}

// readRich reads a string item (si, is): plain text (t) or runs (r).
func (c *converter) readRich(n *ooxml.Node) *richText {
	rt := &richText{}
	var b strings.Builder
	for _, k := range n.Kids {
		switch k.Name {
		case "t":
			b.WriteString(decodeXString(k.Text))
		case "r":
			var f *xfont
			if rp := k.Child("rPr"); rp != nil {
				ff := parseFont(rp, nil)
				ff.name = ""
				ff.size = 0
				if v := rp.Child("rFont").AttrStr("val", ""); v != "" {
					ff.name = v
				}
				if v := rp.Child("sz").AttrFloat("val", 0); v > 0 {
					ff.size = v
				}
				f = &ff
			}
			text := decodeXString(k.Child("t").Text)
			rt.runs = append(rt.runs, textRun{text: text, font: f})
			b.WriteString(text)
		}
	}
	rt.plain = b.String()
	if len(rt.runs) == 1 && rt.runs[0].font == nil {
		rt.runs = nil
	}
	return rt
}

// decodeXString decodes the _xHHHH_ escapes of ST_Xstring and normalizes
// line breaks.
func decodeXString(s string) string {
	if strings.Contains(s, "_x") {
		var b strings.Builder
		for i := 0; i < len(s); i++ {
			if s[i] == '_' && i+6 < len(s) && s[i+1] == 'x' && s[i+6] == '_' {
				if v, err := strconv.ParseUint(s[i+2:i+6], 16, 32); err == nil {
					b.WriteRune(rune(v))
					i += 6
					continue
				}
			}
			b.WriteByte(s[i])
		}
		s = b.String()
	}
	if strings.Contains(s, "\r") {
		s = strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\r", "\n")
	}
	return s
}

// parseRef parses an A1 reference (with optional $) into a 0-based column
// and row.
func parseRef[S ~string | ~[]byte](ref S) (col, row int, ok bool) {
	// spaces around, and the $ of absolute references, are passed over
	i, end := 0, len(ref)
	for i < end && isSpace(ref[i]) {
		i++
	}
	for end > i && isSpace(ref[end-1]) {
		end--
	}
	for ; i < end; i++ {
		c := ref[i]
		if c == '$' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		if c < 'A' || c > 'Z' {
			break
		}
		col = col*26 + int(c-'A'+1)
		if col > maxCols {
			return 0, 0, false
		}
	}
	if col == 0 || i == end {
		return 0, 0, false
	}
	// the row number, as strconv.Atoi reads it
	neg, digits := false, 0
	if ref[i] == '+' || ref[i] == '-' {
		neg = ref[i] == '-'
		i++
	}
	for ; i < end; i++ {
		c := ref[i]
		if c == '$' {
			continue
		}
		if c < '0' || c > '9' {
			return 0, 0, false
		}
		if row = row*10 + int(c-'0'); row > maxRows {
			return 0, 0, false
		}
		digits++
	}
	if digits == 0 || neg || row < 1 {
		return 0, 0, false
	}
	return col - 1, row - 1, true
}

// isSpace reports the white space strings.TrimSpace trims from a reference.
func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f'
}

// parseRange parses "A1:C3" (or a single reference, or whole rows or
// columns "A:C", "3:5").
func parseRange(s string) (cellRange, bool) {
	s = strings.TrimSpace(s)
	a, b, found := strings.Cut(s, ":")
	if !found {
		b = a
	}
	if c0, r0, ok := parseRef(a); ok {
		if c1, r1, ok := parseRef(b); ok {
			return cellRange{min(r0, r1), min(c0, c1), max(r0, r1), max(c0, c1)}, true
		}
		return cellRange{}, false
	}
	// whole columns or rows
	if c0, ok := parseCol(a); ok {
		if c1, ok := parseCol(b); ok {
			return cellRange{0, min(c0, c1), maxRows - 1, max(c0, c1)}, true
		}
	}
	r0, e0 := strconv.Atoi(strings.ReplaceAll(a, "$", ""))
	r1, e1 := strconv.Atoi(strings.ReplaceAll(b, "$", ""))
	if e0 == nil && e1 == nil && r0 >= 1 && r1 >= 1 {
		return cellRange{min(r0, r1) - 1, 0, max(r0, r1) - 1, maxCols - 1}, true
	}
	return cellRange{}, false
}

func parseCol(s string) (int, bool) {
	s = strings.ReplaceAll(s, "$", "")
	col := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		if c < 'A' || c > 'Z' {
			return 0, false
		}
		col = col*26 + int(c-'A'+1)
	}
	return col - 1, col > 0 && col <= maxCols
}

// parseSqref parses a space-separated list of ranges.
func parseSqref(s string) []cellRange {
	var out []cellRange
	for _, f := range strings.Fields(s) {
		if r, ok := parseRange(f); ok {
			out = append(out, r)
		}
	}
	return out
}

// colName returns the letters of a 0-based column.
func colName(c int) string {
	var b []byte
	for c++; c > 0; c = (c - 1) / 26 {
		b = append([]byte{byte('A' + (c-1)%26)}, b...)
	}
	return string(b)
}

func cellRef(row, col int) string { return colName(col) + strconv.Itoa(row+1) }

// sortRows orders rows and their cells, merging duplicates (the last wins).
func (ws *worksheet) sortRows() {
	sort.SliceStable(ws.rows, func(i, j int) bool { return ws.rows[i].idx < ws.rows[j].idx })
	out := ws.rows[:0]
	for _, r := range ws.rows {
		if n := len(out); n > 0 && out[n-1].idx == r.idx {
			out[n-1].cells = append(out[n-1].cells, r.cells...)
			if r.ht >= 0 {
				out[n-1].ht = r.ht
			}
			continue
		}
		out = append(out, r)
	}
	// The backing slice can still point at duplicate rows that were removed.
	clear(ws.rows[len(out):])
	ws.rows = out
	for i := range ws.rows {
		cs := ws.rows[i].cells
		sort.SliceStable(cs, func(a, b int) bool { return cs[a].col < cs[b].col })
		dedup := cs[:0]
		for _, cl := range cs {
			if n := len(dedup); n > 0 && dedup[n-1].col == cl.col {
				dedup[n-1] = cl
				continue
			}
			if cl.col < maxCols {
				dedup = append(dedup, cl)
			}
		}
		clear(cs[len(dedup):])
		ws.rows[i].cells = dedup
	}
}

// rowAt returns the row with an index, or nil.
func (ws *worksheet) rowAt(idx int) *row {
	i := sort.Search(len(ws.rows), func(i int) bool { return ws.rows[i].idx >= idx })
	if i < len(ws.rows) && ws.rows[i].idx == idx {
		return &ws.rows[i]
	}
	return nil
}

// cellAt returns the cell at a position, or nil.
func (ws *worksheet) cellAt(r, c int) *cell {
	rw := ws.rowAt(r)
	if rw == nil {
		return nil
	}
	i := sort.Search(len(rw.cells), func(i int) bool { return rw.cells[i].col >= c })
	if i < len(rw.cells) && rw.cells[i].col == c {
		return &rw.cells[i]
	}
	return nil
}

// cellsIn returns the number of cells of a range that are in the file.
func (ws *worksheet) cellsIn(rg cellRange) int {
	n := 0
	i := sort.Search(len(ws.rows), func(i int) bool { return ws.rows[i].idx >= rg.r0 })
	for ; i < len(ws.rows) && ws.rows[i].idx <= rg.r1; i++ {
		cells := ws.rows[i].cells
		c0 := sort.Search(len(cells), func(k int) bool { return cells[k].col >= rg.c0 })
		c1 := sort.Search(len(cells), func(k int) bool { return cells[k].col > rg.c1 })
		n += c1 - c0
	}
	return n
}

// colStyle returns the format of a column's empty cells (-1 for none).
func (ws *worksheet) colStyle(c int) int {
	ws.indexCols()
	if c < 0 || c >= len(ws.colLast) || ws.colLast[c] < 0 {
		return -1
	}
	return ws.cols[ws.colLast[c]].style
}

// indexCols resolves the col elements by column, once: the last element
// over a column (colLast), whose width and style the column has, and the
// last one that has a style (colBase), which is the base format of the
// column's cells. The elements of a sheet do not overlap; those of a
// malformed one may, and then the last one over a column wins. Looking
// through the elements for each column takes the product of their numbers.
func (ws *worksheet) indexCols() {
	if ws.colsIndexed {
		return
	}
	ws.colsIndexed = true
	if len(ws.cols) == 0 {
		return
	}
	ws.colLast, _ = lastOver(ws.cols, func(colDef) bool { return true })
	ws.colBase, _ = lastOver(ws.cols, func(cd colDef) bool { return cd.style >= 0 })
}

// lastOver returns, for each column of a sheet, the index of the last of the
// col elements that pick accepts over it (-1 for none). It goes through the
// elements from the last one and gives each the columns that none has taken
// yet, so a column is looked at once however many elements are over it;
// steps is the number of elements and columns it looked at.
func lastOver(cols []colDef, pick func(colDef) bool) (out []int32, steps int) {
	out = make([]int32, maxCols)
	// free[c] is a column at or after c that may be free, and is free when
	// it is its own
	free := make([]int32, maxCols+1)
	for c := range free {
		free[c] = int32(c)
		if c < maxCols {
			out[c] = -1
		}
	}
	find := func(c int32) int32 {
		for free[c] != c {
			free[c] = free[free[c]]
			c = free[c]
		}
		return c
	}
	for k := len(cols) - 1; k >= 0; k-- {
		cd := cols[k]
		steps++
		if !pick(cd) || cd.min < 0 || cd.min >= maxCols {
			continue
		}
		for c := find(int32(cd.min)); int(c) <= cd.max && c < maxCols; c = find(c + 1) {
			out[c] = int32(k)
			free[c] = c + 1
			steps++
		}
	}
	return out, steps
}
