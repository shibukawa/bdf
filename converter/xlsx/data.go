package xlsx

// The values of the sheets as tables the document holds beside its drawing
// (docs/spec.md §4.6): the used range of each worksheet, with the letters
// of its columns for names, and each table (ListObject) of the sheet, with
// the names of its columns. A cell's value is typed: a number is a number
// (an integer when the whole column is integral), or a date or a date and
// time when its format shows one; a string, a boolean; an error is a null.
// A column whose values are of several types holds them as the text the
// cell shows.

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
	parquetenc "github.com/shibukawa/bdf/encoding/parquet"
)

// listTable is a table (ListObject) of a sheet, as a table of values: its
// range, its name, and the names of its columns.
type listTable struct {
	ref    cellRange
	name   string
	cols   []string
	header int
	totals int
}

// listTables reads the tables (ListObjects) of a worksheet for their
// values; a table whose range does not parse is skipped.
func (c *converter) listTables(ws *worksheet) []listTable {
	var out []listTable
	for _, rid := range ws.tables {
		r, ok := c.pkg.Target(ws.part, rid)
		if !ok {
			continue
		}
		n, err := c.pkg.XML(r.Target)
		if err != nil {
			continue
		}
		ref, ok := parseRange(n.AttrStr("ref", ""))
		if !ok {
			continue
		}
		ref.r0, ref.r1 = c.filledRows(ws, ref.r0, ref.r1)
		t := listTable{ref: ref, header: int(n.AttrInt("headerRowCount", 1)), totals: int(n.AttrInt("totalsRowCount", 0))}
		t.name = decodeXString(n.AttrStr("displayName", ""))
		if t.name == "" {
			t.name = decodeXString(n.AttrStr("name", ""))
		}
		for _, tc := range n.Path("tableColumns").Children("tableColumn") {
			t.cols = append(t.cols, decodeXString(tc.AttrStr("name", "")))
		}
		out = append(out, t)
	}
	return out
}

// storeSheetTables stores the values of the n-th sheet (from 1), shown in
// the view: its used range as the table sheet<n>, and its tables
// (ListObjects) as table-<k>, numbered through the workbook.
func (c *converter) storeSheetTables(ws *worksheet, n int, view string) {
	mode := c.opts.Data
	if mode == "none" {
		return
	}
	if ref, ok := usedRange(ws); ok {
		c.storeRange(ws, ref, "sheet"+strconv.Itoa(n), ws.name, nil, view, mode)
	}
	for _, l := range c.listTables(ws) {
		ref := l.ref
		ref.r0 += l.header
		ref.r1 -= l.totals
		if ref.r0 > ref.r1 {
			continue
		}
		c.nTables++
		c.storeRange(ws, ref, "table-"+strconv.Itoa(c.nTables), l.name, l.cols, view, mode)
	}
}

// usedRange is the rectangle of the cells of a sheet that hold a value.
func usedRange(ws *worksheet) (cellRange, bool) {
	ref := cellRange{r0: -1, c0: maxCols, r1: -1, c1: -1}
	for _, r := range ws.rows {
		for _, cl := range r.cells {
			if cl.kind == cellBlank {
				continue
			}
			if ref.r0 < 0 {
				ref.r0 = r.idx
			}
			ref.r1 = r.idx
			ref.c0, ref.c1 = min(ref.c0, cl.col), max(ref.c1, cl.col)
		}
	}
	return ref, ref.r0 >= 0
}

// maxDataCells and maxDataCellsAll bound the cells of a range stored as a
// table, under the default of the data parameter and under all: the table
// is built in memory, cell by cell, whether the file holds the cells or
// not.
const (
	maxDataCells    = 1 << 24
	maxDataCellsAll = 1 << 26
)

// storeRange stores the values of a range of a sheet as a table, within
// what the data parameter's mode allows.
func (c *converter) storeRange(ws *worksheet, ref cellRange, id, title string, names []string, view, mode string) {
	rows := int64(ref.r1-ref.r0) + 1
	cols := int64(ref.c1-ref.c0) + 1
	if mode != "all" && rows > conv.DefaultTableRows {
		c.warnf("sheet %q: the values (%d rows) are not stored as the table %s: more than the default allows (data=all stores them)", ws.name, rows, id)
		return
	}
	if cols > parquetenc.MaxColumns {
		c.warnf("sheet %q: the values are not stored as the table %s: more than %d columns", ws.name, id, parquetenc.MaxColumns)
		return
	}
	if cells, limit := rows*cols, int64(maxDataCells); cells > limit {
		if mode == "all" {
			limit = maxDataCellsAll
		}
		if cells > limit {
			c.warnf("sheet %q: the values (%d rows × %d columns) are not stored as the table %s: more than %d cells", ws.name, rows, cols, id, limit)
			return
		}
	}
	t := c.valuesTable(ws, ref, names)
	data, err := parquetenc.Encode(t)
	if err != nil {
		c.warnf("sheet %q: the values are not stored as the table %s: %v", ws.name, id, err)
		return
	}
	if !conv.TableFits(mode, rows, len(data)) {
		c.warnf("sheet %q: the values (%d rows, %d MiB) are not stored as the table %s: more than the default allows (data=all stores them)", ws.name, rows, len(data)>>20, id)
		return
	}
	entry := &bdf.Table{ID: id, Title: title, Rows: rows, View: view}
	for _, col := range t.Columns {
		entry.Columns = append(entry.Columns, bdf.TableColumn{Name: col.Name, Type: col.Type.String()})
	}
	if err := c.doc.AddTable(entry, data); err != nil {
		c.warnf("sheet %q: the values are not stored as the table %s: %v", ws.name, id, err)
	}
}

// valuesTable reads the values of a range of a sheet as a table whose
// columns are named names, or by their letters when names are fewer.
func (c *converter) valuesTable(ws *worksheet, ref cellRange, names []string) *parquetenc.Table {
	nCols := ref.c1 - ref.c0 + 1
	// the rows of the range, by their index in it
	rows := make([][]cell, ref.r1-ref.r0+1)
	for i := range ws.rows {
		r := &ws.rows[i]
		if r.idx >= ref.r0 && r.idx <= ref.r1 {
			rows[r.idx-ref.r0] = r.cells
		}
	}
	// the type of a column is that of its values, or string when they are
	// of several types
	types := make([]parquetenc.Type, nCols)
	for _, cells := range rows {
		for _, cl := range cells {
			if cl.col < ref.c0 || cl.col > ref.c1 {
				continue
			}
			k := c.cellType(&cl)
			if k == parquetenc.Other {
				continue
			}
			j := cl.col - ref.c0
			switch t := types[j]; {
			case t == parquetenc.Other, t == k:
				types[j] = k
			case t == parquetenc.Integer && k == parquetenc.Number, t == parquetenc.Number && k == parquetenc.Integer:
				types[j] = parquetenc.Number
			case t == parquetenc.Date && k == parquetenc.DateTime, t == parquetenc.DateTime && k == parquetenc.Date:
				types[j] = parquetenc.DateTime
			default:
				types[j] = parquetenc.String
			}
		}
	}
	t := &parquetenc.Table{Columns: make([]parquetenc.Column, nCols), Rows: make([][]any, len(rows))}
	for j := range t.Columns {
		name := ""
		if j < len(names) {
			name = names[j]
		}
		if name == "" {
			name = colName(ref.c0 + j)
		}
		typ := types[j]
		if typ == parquetenc.Other {
			typ = parquetenc.String
		}
		t.Columns[j] = parquetenc.Column{Name: name, Type: typ}
	}
	for i, cells := range rows {
		vals := make([]any, nCols)
		for _, cl := range cells {
			if cl.col >= ref.c0 && cl.col <= ref.c1 {
				vals[cl.col-ref.c0] = c.tableValue(&cl, t.Columns[cl.col-ref.c0].Type)
			}
		}
		t.Rows[i] = vals
	}
	return t
}

// cellType is the type of a cell's value: Other for a blank or an error
// (a null in any column).
func (c *converter) cellType(cl *cell) parquetenc.Type {
	switch cl.kind {
	case cellNum:
		if nf := c.st.numFormat(c.st.xf(cl.style).numFmt); nf.isDate() {
			if nf.hasTime() || cl.num != math.Trunc(cl.num) {
				return parquetenc.DateTime
			}
			return parquetenc.Date
		}
		if cl.num == math.Trunc(cl.num) && math.Abs(cl.num) < 1<<53 {
			return parquetenc.Integer
		}
		return parquetenc.Number
	case cellStr:
		return parquetenc.String
	case cellBool:
		return parquetenc.Boolean
	}
	return parquetenc.Other
}

// tableValue is the value of a cell in a column of a type: nil for a blank,
// an error, or a number no date can be made of; the text the cell shows
// in a column of strings.
func (c *converter) tableValue(cl *cell, typ parquetenc.Type) any {
	switch cl.kind {
	case cellNum:
		switch typ {
		case parquetenc.Integer:
			return int64(cl.num)
		case parquetenc.Number:
			return cl.num
		case parquetenc.Date, parquetenc.DateTime:
			t, _, ok := serialTime(cl.num, c.date1904)
			if !ok {
				return nil
			}
			if typ == parquetenc.Date {
				return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			}
			return t
		}
		return c.st.numFormat(c.st.xf(cl.style).numFmt).formatNumber(cl.num, c.date1904).text()
	case cellStr:
		if cl.text == nil {
			return ""
		}
		return cl.text.plain
	case cellBool:
		if typ == parquetenc.Boolean {
			return cl.num != 0
		}
		if cl.num != 0 {
			return "TRUE"
		}
		return "FALSE"
	}
	return nil
}

// hasTime reports whether a date format shows a time of day (or an
// elapsed time) as well.
func (nf *numFormat) hasTime() bool {
	if len(nf.sections) == 0 {
		return false
	}
	for _, tok := range nf.sections[0].toks {
		switch tok.kind {
		case tElapsed, tSubsec:
			return true
		case tDate:
			s := strings.ToLower(tok.s)
			if strings.HasPrefix(s, "h") || strings.HasPrefix(s, "s") || strings.HasPrefix(s, "a") {
				return true
			}
		}
	}
	return false
}
