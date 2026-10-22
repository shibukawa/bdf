package parquet_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	full "github.com/shibukawa/bdf/converter/parquet"
	"github.com/shibukawa/bdf/encoding/parquet"
)

func sample() *parquet.Table {
	t := &parquet.Table{Columns: []parquet.Column{
		{Name: "name", Type: parquet.String}, {Name: "value", Type: parquet.Number}, {Name: "count", Type: parquet.Integer},
		{Name: "ok", Type: parquet.Boolean}, {Name: "day", Type: parquet.Date}, {Name: "at", Type: parquet.DateTime},
	}}
	day := time.Date(2026, 10, 22, 0, 0, 0, 0, time.UTC)
	at := time.Date(2026, 10, 22, 13, 45, 6, 789_000_000, time.UTC)
	t.Rows = [][]any{
		{"alpha", 1.5, int64(1), true, day, at},
		{"βeta", -2.25, int64(-40), false, day.AddDate(0, 1, 0), at.Add(time.Hour)},
		{nil, nil, nil, nil, nil, nil},
		{"", 0.0, int64(0), true, unixDay(), time.UnixMilli(0).UTC()},
		{"delta", 3e10, int64(1) << 40, false}, // (a short row: nulls at its end)
	}
	return t
}

func unixDay() time.Time { return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC) }

func TestRoundTrip(t *testing.T) {
	want := sample()
	data, err := parquet.Encode(want)
	if err != nil {
		t.Fatal(err)
	}
	cols, rows, err := parquet.Schema(data)
	if err != nil {
		t.Fatal(err)
	}
	if rows != 5 || fmt.Sprint(cols) != fmt.Sprint(want.Columns) {
		t.Fatalf("schema: %d rows, %v", rows, cols)
	}
	got, err := parquet.Read(data, parquet.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Rows) != 5 {
		t.Fatalf("%d rows", len(got.Rows))
	}
	for i, r := range want.Rows {
		for j := range want.Columns {
			var w any
			if j < len(r) {
				w = r[j]
			}
			g := got.Rows[i][j]
			if tw, ok := w.(time.Time); ok {
				if tg, ok := g.(time.Time); !ok || !tg.Equal(tw) {
					t.Errorf("row %d %s = %v, want %v", i+1, want.Columns[j].Name, g, w)
				}
				continue
			}
			if fmt.Sprintf("%T %v", g, g) != fmt.Sprintf("%T %v", w, w) {
				t.Errorf("row %d %s = %T %v, want %T %v", i+1, want.Columns[j].Name, g, g, w, w)
			}
		}
	}
	// the limit
	few, err := parquet.Read(data, parquet.Limits{MaxRows: 2})
	if err != nil || len(few.Rows) != 2 || few.Rows[1][0] != "βeta" {
		t.Fatalf("limited: %v, %v", few, err)
	}
}

// TestPages writes more rows than a page holds: the values span pages.
func TestPages(t *testing.T) {
	tb := &parquet.Table{Columns: []parquet.Column{{Name: "i", Type: parquet.Integer}, {Name: "s", Type: parquet.String}, {Name: "b", Type: parquet.Boolean}}}
	const n = 70_000
	for i := range n {
		var row []any
		if i%7 == 3 {
			row = []any{nil, nil, nil}
		} else {
			row = []any{int64(i), strings.Repeat("x", i%20), i%3 == 0}
		}
		tb.Rows = append(tb.Rows, row)
	}
	data, err := parquet.Encode(tb)
	if err != nil {
		t.Fatal(err)
	}
	got, err := parquet.Read(data, parquet.Limits{MaxRows: n})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Rows) != n {
		t.Fatalf("%d rows", len(got.Rows))
	}
	for i, r := range tb.Rows {
		if fmt.Sprint(got.Rows[i]) != fmt.Sprint(r) {
			t.Fatalf("row %d = %v, want %v", i, got.Rows[i], r)
		}
	}
}

// TestFullReader reads a written file with the converter's reader, which
// was written from the format's specification independently.
func TestFullReader(t *testing.T) {
	data, err := parquet.Encode(sample())
	if err != nil {
		t.Fatal(err)
	}
	res, err := full.Convert(bytes.NewReader(data), int64(len(data)), &full.Options{Name: "t", Rows: -1, NoSystemFonts: true, SystemFonts: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.Rows != 5 || res.Cols != 6 || res.Shown != 5 {
		t.Fatalf("rows %d cols %d shown %d", res.Rows, res.Cols, res.Shown)
	}
	for _, w := range res.Warnings {
		if strings.Contains(w, "column") {
			t.Errorf("warning: %s", w)
		}
	}
	st, err := res.Doc.SearchText()
	if err != nil {
		t.Fatal(err)
	}
	text := st.Views[0].Pages[0].Text
	for _, want := range []string{"name", "βeta", "-2.25", "2026-10-22", "13:45:06", "true", "1099511627776"} {
		if !strings.Contains(text, want) {
			t.Errorf("%q is not in the sheet:\n%s", want, text)
		}
	}
}

// TestBadInput refuses what is not a Parquet file, and what it cannot
// write.
func TestBadInput(t *testing.T) {
	if _, _, err := parquet.Schema([]byte("PAR1 not really PAR1")); !errors.Is(err, parquet.ErrNotParquet) && err == nil {
		t.Error("a short file was read")
	}
	if _, _, err := parquet.Schema([]byte("hello")); !errors.Is(err, parquet.ErrNotParquet) {
		t.Errorf("not parquet: %v", err)
	}
	if _, err := parquet.Encode(&parquet.Table{}); err == nil {
		t.Error("a table without columns was written")
	}
	if _, err := parquet.Encode(&parquet.Table{Columns: []parquet.Column{{Name: "x", Type: parquet.Integer}}, Rows: [][]any{{"one"}}}); err == nil {
		t.Error("a string in an integer column was written")
	}
	if _, err := parquet.Encode(&parquet.Table{Columns: []parquet.Column{{Name: "x", Type: parquet.Other}}}); err == nil {
		t.Error("a column of type other was written")
	}
	if _, err := parquet.Encode(&parquet.Table{Columns: []parquet.Column{{Name: "x", Type: parquet.Number}}, Rows: [][]any{{1, 2}}}); err == nil {
		t.Error("a row longer than the columns was written")
	}
}

// TestFiles reads the files of other writers in converter/parquet's
// testdata: the flat ones, and the schema of the nested one.
func TestFiles(t *testing.T) {
	for _, name := range []string{"basic.parquet", "duckdb.parquet", "types.parquet", "nested.parquet", "codecs.parquet", "encodings.parquet"} {
		data, err := os.ReadFile("../../converter/parquet/testdata/" + name)
		if err != nil {
			t.Skip(err)
		}
		cols, rows, err := parquet.Schema(data)
		if err != nil {
			t.Errorf("%s: schema: %v", name, err)
			continue
		}
		tb, err := parquet.Read(data, parquet.Limits{})
		kinds := map[parquet.Type]int{}
		for _, c := range cols {
			kinds[c.Type]++
		}
		switch {
		case err != nil:
			// the file holds what the package does not read: it says so
			if !strings.Contains(err.Error(), "converter/parquet") && !strings.Contains(err.Error(), "repeated") {
				t.Errorf("%s: %v", name, err)
			}
			t.Logf("%s: %d columns %v, %d rows: not read: %v", name, len(cols), kinds, rows, err)
		case int64(len(tb.Rows)) != rows:
			t.Errorf("%s: %d of %d rows", name, len(tb.Rows), rows)
		default:
			t.Logf("%s: %d columns %v, %d rows: read", name, len(cols), kinds, rows)
		}
	}
}
