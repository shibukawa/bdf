package xlsx

import (
	"bytes"
	"encoding/xml"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf/converter/internal/ooxml"
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
}

// readWorksheet parses a worksheet part.
func (c *converter) readWorksheet(part string) (*worksheet, error) {
	data, err := c.pkg.Read(part)
	if err != nil {
		return nil, err
	}
	ws := &worksheet{part: part, baseColW: 8, showGrid: true, showZero: true}
	d := xml.NewDecoder(bytes.NewReader(data))
	d.Strict = false
	depth := 0
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			depth++
			if depth != 2 {
				continue
			}
			if t.Name.Local == "sheetData" {
				if err := c.readSheetData(d, ws); err != nil {
					return nil, err
				}
				depth--
				continue
			}
			n, err := ooxml.ReadElement(d, t)
			if err != nil {
				return nil, err
			}
			depth--
			ws.element(n)
		case xml.EndElement:
			depth--
		}
	}
	return ws, nil
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
		ws.defColW = n.AttrFloat("defaultColWidth", 0)
		ws.baseColW = n.AttrFloat("baseColWidth", 8)
		ws.defRowH = n.AttrFloat("defaultRowHeight", 0)
		ws.zeroH = n.AttrBool("zeroHeight", false)
	case "cols":
		for _, k := range n.Children("col") {
			cd := colDef{min: int(k.AttrInt("min", 1)) - 1, max: int(k.AttrInt("max", 1)) - 1, width: -1,
				hidden: k.AttrBool("hidden", false), style: int(k.AttrInt("style", -1))}
			if _, ok := k.Attr("width"); ok {
				cd.width = k.AttrFloat("width", 0)
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
)

// readSheetData streams the rows and cells of sheetData.
func (c *converter) readSheetData(d *xml.Decoder, ws *worksheet) error {
	nextRow := 0
	var cur *row
	var cl *cell
	nextCol := 0
	var cellType string
	var inV, hasF, hasV bool
	var vbuf strings.Builder
	var isNode *ooxml.Node
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "row":
				r := row{idx: nextRow, ht: -1, style: -1}
				for _, a := range t.Attr {
					switch a.Name.Local {
					case "r":
						if n, err := strconv.Atoi(a.Value); err == nil && n >= 1 {
							r.idx = n - 1
						}
					case "ht":
						if f, err := strconv.ParseFloat(a.Value, 64); err == nil && f >= 0 {
							r.ht = f
						}
					case "customHeight":
						r.custom = xmlBool(a.Value)
					case "hidden":
						r.hidden = xmlBool(a.Value)
					case "s":
						if n, err := strconv.Atoi(a.Value); err == nil {
							r.style = n
						}
					}
				}
				// the row's format applies to its empty cells only with customFormat
				if !hasTrue(t.Attr, "customFormat") {
					r.style = -1
				}
				if r.idx >= maxRows {
					r.idx = maxRows - 1
				}
				ws.rows = append(ws.rows, r)
				cur = &ws.rows[len(ws.rows)-1]
				nextRow = r.idx + 1
				nextCol = 0
			case "c":
				if cur == nil {
					ws.rows = append(ws.rows, row{idx: nextRow, ht: -1, style: -1})
					cur = &ws.rows[len(ws.rows)-1]
					nextRow++
				}
				cc := cell{col: nextCol, style: 0}
				cellType = "n"
				for _, a := range t.Attr {
					switch a.Name.Local {
					case "r":
						if col, _, ok := parseRef(a.Value); ok {
							cc.col = col
						}
					case "s":
						if n, err := strconv.Atoi(a.Value); err == nil {
							cc.style = n
						}
					case "t":
						cellType = a.Value
					}
				}
				nextCol = cc.col + 1
				cur.cells = append(cur.cells, cc)
				cl = &cur.cells[len(cur.cells)-1]
				vbuf.Reset()
				hasF, hasV = false, false
			case "v":
				inV, hasV = true, true
				vbuf.Reset()
			case "is":
				n, err := ooxml.ReadElement(d, t)
				if err != nil {
					return err
				}
				isNode = n
				hasV = true
			default:
				if t.Name.Local == "f" {
					hasF = true
				}
				// formulas and extensions: skip their content
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.CharData:
			if inV {
				vbuf.Write(t)
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "v":
				inV = false
			case "c":
				if cl != nil {
					c.cellValue(cl, cellType, vbuf.String(), isNode)
					if hasF && (!hasV || strings.TrimSpace(vbuf.String()) == "" && isNode == nil) {
						ws.noValue++
					}
				}
				cl, isNode = nil, nil
			case "row":
				cur = nil
			case "sheetData":
				return nil
			}
		}
	}
}

func hasTrue(attrs []xml.Attr, name string) bool {
	for _, a := range attrs {
		if a.Name.Local == name {
			return xmlBool(a.Value)
		}
	}
	return false
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
func parseRef(s string) (col, row int, ok bool) {
	s = strings.TrimSpace(strings.ReplaceAll(s, "$", ""))
	i := 0
	for i < len(s) && (s[i] >= 'A' && s[i] <= 'Z' || s[i] >= 'a' && s[i] <= 'z') {
		c := s[i]
		if c >= 'a' {
			c -= 'a' - 'A'
		}
		col = col*26 + int(c-'A'+1)
		i++
		if col > maxCols {
			return 0, 0, false
		}
	}
	if i == 0 || i == len(s) {
		return 0, 0, false
	}
	r, err := strconv.Atoi(s[i:])
	if err != nil || r < 1 || r > maxRows {
		return 0, 0, false
	}
	return col - 1, r - 1, true
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

// colStyle returns the format of a column's empty cells (-1 for none).
func (ws *worksheet) colStyle(c int) int {
	for i := len(ws.cols) - 1; i >= 0; i-- {
		if cd := ws.cols[i]; c >= cd.min && c <= cd.max {
			return cd.style
		}
	}
	return -1
}
