package xlsx

import (
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	parquetenc "github.com/shibukawa/bdf/encoding/parquet"
)

// readTable reads a table of a converted document.
func readTable(t *testing.T, r *bdf.Reader, id string) (*bdf.Table, *parquetenc.Table) {
	t.Helper()
	entry, data, err := r.Table(id)
	if err != nil {
		t.Fatal(err)
	}
	tb, err := parquetenc.Read(data, parquetenc.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	return entry, tb
}

// TestDataTables checks the tables of values of a workbook: the used range
// of a sheet, its table (ListObject) with typed columns, and the charts,
// placed in their views.
func TestDataTables(t *testing.T) {
	_, r := convert(t, "features.xlsx", testOptions())
	var ids []string
	for _, tb := range r.Manifest.Tables {
		ids = append(ids, tb.ID)
	}
	if len(ids) != 4 || ids[0] != "sheet1" || ids[1] != "table-1" || ids[2] != "chart-1" || ids[3] != "chart-2" {
		t.Fatalf("tables = %v", ids)
	}
	// the sheet: its used range, with the letters of the columns; the
	// values are the text shown, the columns holding headings too
	entry, tb := readTable(t, r, "sheet1")
	if entry.Title != "Features" || entry.View != "sheet1" || entry.Page != 0 || entry.Rows != 21 || len(entry.Columns) != 7 || entry.Columns[0].Name != "A" || entry.Columns[6].Name != "G" {
		t.Errorf("sheet1 = %+v", entry)
	}
	if v := tb.Rows[2][0]; v != "3" {
		t.Errorf("A3 = %#v", v)
	}
	// the table: the names of its columns and their types; the header row
	// is not a row of it
	entry, tb = readTable(t, r, "table-1")
	if entry.Title != "Supplies" || entry.View != "sheet1" || entry.Rows != 4 {
		t.Errorf("table-1 = %+v", entry)
	}
	cols := []string{}
	for _, c := range entry.Columns {
		cols = append(cols, c.Name+":"+c.Type)
	}
	if len(cols) != 3 || cols[0] != "Item:string" || cols[1] != "Units:integer" || cols[2] != "Price:number" {
		t.Errorf("columns = %v", cols)
	}
	if row := tb.Rows[0]; row[0] != "Pens" || row[1] != int64(120) || row[2] != 1.5 {
		t.Errorf("row 1 = %#v", row)
	}
	// the charts: the one in the sheet is placed in it, the chart sheet's
	// on its page
	if c1 := r.Manifest.Table("chart-1"); c1.View != "sheet1" || c1.Page != 0 || c1.Title != "Units" {
		t.Errorf("chart-1 = %+v", c1)
	}
	if c2 := r.Manifest.Table("chart-2"); c2.View != "sheet2" || c2.Page != 1 {
		t.Errorf("chart-2 = %+v", c2)
	}
	// data=none stores none
	opts := testOptions()
	opts.Data = "none"
	_, r = convert(t, "features.xlsx", opts)
	if len(r.Manifest.Tables) != 0 {
		t.Errorf("with Data none, tables = %d", len(r.Manifest.Tables))
	}
}

// TestValuesTable checks how the cells of a range are typed: numbers by
// their format (dates, times), booleans, errors as nulls, and columns of
// mixed values as the text shown.
func TestValuesTable(t *testing.T) {
	c := newConverter(nil, testOptions())
	c.st = loadStyles(c, nil, defaultThemeColor)
	c.st.cellXfs = []xf{{}, {numFmt: 14}, {numFmt: 22}, {numFmt: 9}}
	str := func(s string) *richText { return &richText{plain: s} }
	ws := &worksheet{name: "T", rows: []row{
		{idx: 1, cells: []cell{{col: 1, kind: cellNum, num: 1}, {col: 2, kind: cellNum, num: 45000, style: 1}, {col: 3, kind: cellNum, num: 45000.5, style: 2}, {col: 4, kind: cellBool, num: 1}, {col: 5, kind: cellNum, num: 0.25, style: 3}, {col: 6, kind: cellNum, num: 2.5}}},
		{idx: 3, cells: []cell{{col: 1, kind: cellNum, num: 2.5}, {col: 2, kind: cellNum, num: 45001, style: 1}, {col: 3, kind: cellErr, text: str("#N/A")}, {col: 4, kind: cellBlank}, {col: 5, kind: cellStr, text: str("n/a")}, {col: 6, kind: cellNum, num: 3}}},
	}}
	ref, ok := usedRange(ws)
	if !ok || ref != (cellRange{r0: 1, c0: 1, r1: 3, c1: 6}) {
		t.Fatalf("used range = %+v %v", ref, ok)
	}
	tb := c.valuesTable(ws, ref, []string{"n"})
	cols := []string{}
	for _, col := range tb.Columns {
		cols = append(cols, col.Name+":"+col.Type.String())
	}
	want := []string{"n:number", "C:date", "D:datetime", "E:boolean", "F:string", "G:number"}
	if len(cols) != len(want) {
		t.Fatalf("columns = %v", cols)
	}
	for i := range want {
		if cols[i] != want[i] {
			t.Errorf("column %d = %s, want %s", i, cols[i], want[i])
		}
	}
	if len(tb.Rows) != 3 {
		t.Fatalf("rows = %d", len(tb.Rows))
	}
	r1, r2, r3 := tb.Rows[0], tb.Rows[1], tb.Rows[2]
	if r1[0] != 1.0 || r3[0] != 2.5 {
		t.Errorf("numbers = %#v %#v", r1[0], r3[0])
	}
	if d, _ := r1[1].(time.Time); !d.Equal(time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("date = %v", r1[1])
	}
	if d, _ := r1[2].(time.Time); !d.Equal(time.Date(2023, 3, 15, 12, 0, 0, 0, time.UTC)) {
		t.Errorf("datetime = %v", r1[2])
	}
	if r1[3] != true || r3[3] != nil {
		t.Errorf("booleans = %#v %#v", r1[3], r3[3])
	}
	if r1[4] != "25%" || r3[4] != "n/a" {
		t.Errorf("mixed column = %#v %#v", r1[4], r3[4])
	}
	if r3[2] != nil {
		t.Errorf("error = %#v", r3[2])
	}
	for i, v := range r2 {
		if v != nil {
			t.Errorf("empty row: %d = %#v", i, v)
		}
	}
}
