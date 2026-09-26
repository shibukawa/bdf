package xlsx

import (
	"fmt"
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/fontset"
)

// Grid is a sheet of values without formats, such as the records of a CSV
// file. ConvertGrid shows it the way Excel shows such a sheet in a new
// workbook: the default font (游ゴシック for Japanese), gridlines, text on
// the left and numbers on the right, text that flows into empty cells,
// and values with line breaks wrapped in taller rows. The columns are as
// wide as their text, up to about 50 digits.
type Grid struct {
	// Name is the sheet name ("Sheet1" when empty).
	Name string
	// Rows holds the values row by row; rows may differ in length, and
	// empty strings are empty cells.
	Rows [][]GridCell
	// HeaderRows is the number of rows at the top that head the columns:
	// they are bold, stay in view (a frozen pane) and are marked as column
	// headers for screen readers.
	HeaderRows int
	// Lang is the language of the text ("ja", "ko", "zh-CN", "zh-TW" or
	// others; "" when unknown). It picks the font, as Excel's default font
	// depends on the language of Office, and tells the language of East
	// Asian text without kana or hangul.
	Lang string
	// TableStyle formats the grid as an Excel table with a built-in table
	// style ("TableStyleMedium2", or "Medium2" for short): banded rows and,
	// under header rows, filter buttons. "" leaves the cells plain.
	TableStyle string
}

// GridCell is a value of a Grid.
type GridCell struct {
	// Text is the value as it is shown.
	Text string
	// Number aligns the value to the right, as Excel aligns numbers, dates
	// and times; like those it does not flow into the next cell, and turns
	// into "###" when its column is too narrow. The text is shown as it is
	// (Excel would reformat it: 1.50 as 1.5, 007 as 7).
	Number bool
}

// cell formats of grids
const (
	gridPlain = iota
	gridHeader
	gridWrap
	gridHeaderWrap
)

// fitMaxDigits caps the width of the columns fitted to their text.
const fitMaxDigits = 50

// ConvertGrid converts a grid into a document with one sheet view. Of the
// options, Sheets and Hidden do not apply.
func ConvertGrid(g *Grid, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	tableStyle := ""
	if g.TableStyle != "" {
		if tableStyle = builtinTableStyleName(g.TableStyle); tableStyle == "" {
			return nil, fmt.Errorf("xlsx: unknown table style %q", g.TableStyle)
		}
	}
	c := newConverter(nil, opts)
	if opts.Title != "" {
		c.doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	c.eaScript = "Jpan" // as for workbooks of unknown language
	if s := fontset.Script(g.Lang); s != "" {
		c.eaScript, c.eaLang = s, g.Lang
		if s == "Jpan" {
			c.locale = "ja"
		}
	}
	font := xfont{name: gridFont(g.Lang), size: 11, family: 2}
	bold := font
	bold.bold = true
	c.st = loadStyles(c, nil, defaultThemeColor)
	c.st.fonts = []xfont{font, bold}
	plain := xalign{h: "general", v: "bottom"}
	wrap := plain
	wrap.wrap = true
	c.st.cellXfs = []xf{gridPlain: {align: plain}, gridHeader: {font: 1, align: plain},
		gridWrap: {align: wrap}, gridHeaderWrap: {font: 1, align: wrap}}
	c.mdw = c.maxDigitWidth(c.st.font(0))

	name := g.Name
	if name == "" {
		name = "Sheet1"
	}
	ws := &worksheet{name: name, baseColW: 8, showGrid: true, showZero: true, fitCols: true}
	lastCol := -1
	for i, cells := range g.Rows {
		if i >= maxRows {
			c.warnf("only the first %d rows are shown (the most a sheet has)", maxRows)
			break
		}
		if len(cells) > maxCols {
			c.warnOnce("gridcols", "only the first %d columns are shown (the most a sheet has)", maxCols)
			cells = cells[:maxCols]
		}
		rw := row{idx: i, ht: -1, style: -1}
		for j, v := range cells {
			if v.Text == "" {
				continue
			}
			text := v.Text
			if strings.Contains(text, "\r") {
				text = strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", "\n"), "\r", "\n")
			}
			style := gridPlain
			if i < g.HeaderRows {
				style = gridHeader
			}
			if strings.Contains(text, "\n") {
				style += gridWrap
			}
			cl := cell{col: j, style: style, kind: cellStr, text: &richText{plain: text}}
			if v.Number {
				cl.kind = cellNum
			}
			rw.cells = append(rw.cells, cl)
			lastCol = max(lastCol, j)
		}
		if len(rw.cells) > 0 {
			ws.rows = append(ws.rows, rw)
		}
	}
	nRows := min(len(g.Rows), maxRows)
	header := min(max(g.HeaderRows, 0), nRows)
	ws.headerRows, ws.freezeR = header, header
	if tableStyle != "" && lastCol >= 0 {
		ws.gridTables = []gridTable{{ref: cellRange{0, 0, nRows - 1, lastCol}, header: header, style: tableStyle}}
	}
	c.parsed[name] = ws
	v, tiles, err := c.worksheetSafe(sheetRef{name: name}, "sheet1")
	if err != nil {
		return nil, fmt.Errorf("xlsx: %w", err)
	}
	c.finish([]pendingView{{view: v, tiles: tiles}})
	return &Result{Doc: c.doc, Warnings: c.warnings, Sheets: 1, EmbeddedFonts: c.embeddedFonts}, nil
}

// gridFont is the default font of a new workbook in the Office of a
// language: that of the default theme's minor font for the language's
// script, or Calibri.
func gridFont(lang string) string {
	switch fontset.Script(lang) {
	case "Jpan":
		return "游ゴシック"
	case "Hang":
		return "맑은 고딕"
	case "Hans":
		return "等线"
	case "Hant":
		return "新細明體"
	}
	return "Calibri"
}

// defaultThemeColor is the color scheme of Office's default theme (2013 to
// 2022), which the built-in table styles of grids refer to.
func defaultThemeColor(name string) (rgb, bool) {
	v, ok := map[string]uint32{"dk1": 0x000000, "lt1": 0xFFFFFF, "dk2": 0x44546A, "lt2": 0xE7E6E6,
		"accent1": 0x4472C4, "accent2": 0xED7D31, "accent3": 0xA5A5A5, "accent4": 0xFFC000,
		"accent5": 0x5B9BD5, "accent6": 0x70AD47, "hlink": 0x0563C1, "folHlink": 0x954F72}[name]
	return hexRGB(v), ok
}

// builtinTableStyleName returns the name of the built-in table style that
// name stands for (case aside, and without "TableStyle"), or "".
func builtinTableStyleName(name string) string {
	builtinTableStyle("")
	for _, n := range []string{name, "TableStyle" + name} {
		for k := range builtinStyles {
			if strings.EqualFold(k, n) {
				return k
			}
		}
	}
	return ""
}

// gridTable is a table that a grid is formatted as.
type gridTable struct {
	ref    cellRange
	header int
	style  string
}

// fitWidths returns the widths that the columns of a sheet whose columns
// fit their text need: the widest value plus the cell padding and a pixel
// on each side (and a filter button in table headers), at least the
// default width, and for text at most fitMaxDigits digits.
func (s *sheetCtx) fitWidths(def float64) map[int]float64 {
	styles := map[int]*tstyle{}
	widest := map[int]float64{}
	numeric := map[int]float64{}
	for i := range s.ws.rows {
		rw := &s.ws.rows[i]
		for j := range rw.cells {
			cl := &rw.cells[j]
			if cl.text == nil {
				continue
			}
			st, ok := styles[cl.style]
			if !ok {
				st = s.textStyle(&s.fmtOf(cl.style).font, nil, nil)
				styles[cl.style] = st
			}
			w := s.textWidth(st, cl.text.plain)
			if t := s.tableAt(rw.idx, cl.col); t != nil && t.filter && t.isHeader(rw.idx) {
				w += 16 * pxPt
			}
			if cl.kind == cellNum {
				numeric[cl.col] = max(numeric[cl.col], w)
			} else {
				widest[cl.col] = max(widest[cl.col], w)
			}
		}
	}
	limit := (fitMaxDigits*s.c.mdw + 5) * pxPt
	fit := func(w float64) float64 { return (math.Ceil(w/pxPt-0.01) + 5 + 2) * pxPt }
	out := map[int]float64{}
	for col, w := range widest {
		out[col] = max(def, min(fit(w), limit))
	}
	for col, w := range numeric {
		out[col] = max(out[col], def, fit(w))
	}
	return out
}

// textWidth measures the widest line of a value in a style.
func (s *sheetCtx) textWidth(st *tstyle, text string) float64 {
	best, w, n := 0.0, 0.0, 0
	for _, r := range text {
		if r == '\n' {
			best, w = max(best, w), 0
			continue
		}
		if r == '\t' {
			r = ' '
		}
		if r < 0x20 {
			continue
		}
		w += s.item(st, r).w
		if n++; n >= maxShown {
			break
		}
	}
	return max(best, w)
}

// loadGridTables adds the tables a grid is formatted as.
func (s *sheetCtx) loadGridTables() {
	for _, gt := range s.ws.gridTables {
		t := &table{ref: gt.ref, header: gt.header, elems: map[string]styleElem{}, filterButtons: map[int]bool{},
			filter: gt.header > 0, rowStripes: true}
		s.tableStyle(t, gt.style)
		s.tables = append(s.tables, t)
	}
}
