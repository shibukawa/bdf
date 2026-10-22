package bdf_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/encoding/parquet"
)

// TestTables adds tables to a document, writes it and reads them back.
func TestTables(t *testing.T) {
	d := bdf.NewDocument()
	v := d.NewView("slides", bdf.ViewFixed, "Deck")
	v.AddPage(100, 100, bdf.Layer{Role: bdf.RoleBody, Obj: d.AddPart(bdf.PartObject, bdf.NewObject().Encode())})
	rows := &parquet.Table{Columns: []parquet.Column{{Name: "Quarter", Type: parquet.String}, {Name: "Sales", Type: parquet.Number}, {Name: "When", Type: parquet.Date}},
		Rows: [][]any{{"Q1", 20.5, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}, {"Q2", nil, nil}}}
	tb, err := d.AddTableRows("chart-1", "Sales", rows, "slides", 1)
	if err != nil {
		t.Fatal(err)
	}
	if tb.Rows != 2 || len(tb.Columns) != 3 || tb.Columns[1].Type != "number" || tb.Part == "" {
		t.Fatalf("table = %+v", tb)
	}
	if _, err := d.AddTableRows("chart-1", "", rows, "", 0); err == nil || !strings.Contains(err.Error(), "already") {
		t.Fatalf("a second chart-1: %v", err)
	}
	if err := d.AddTable(&bdf.Table{ID: "x"}, []byte("nope")); err == nil {
		t.Fatal("a table that is not Parquet was added")
	}
	var b bytes.Buffer
	if err := d.WriteSingle(&b); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(b.Bytes()), int64(b.Len()))
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Manifest.Tables) != 1 || r.Manifest.Table("chart-1") == nil || r.Manifest.Table("nope") != nil {
		t.Fatalf("tables = %+v", r.Manifest.Tables)
	}
	entry, data, err := r.Table("chart-1")
	if err != nil {
		t.Fatal(err)
	}
	if entry.View != "slides" || entry.Page != 1 || entry.Title != "Sales" {
		t.Errorf("entry = %+v", entry)
	}
	got, err := parquet.Read(data, parquet.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Rows) != 2 || got.Rows[0][0] != "Q1" || got.Rows[0][1] != 20.5 || got.Rows[1][1] != nil {
		t.Errorf("rows = %v", got.Rows)
	}
	if _, _, err := r.Table("nope"); err == nil {
		t.Error("a table that is not there was read")
	}
	// the document read back carries the tables
	d2, err := r.ToDocument()
	if err != nil || len(d2.Tables) != 1 || d2.Table("chart-1") == nil {
		t.Fatalf("ToDocument: %v, %v", d2.Tables, err)
	}
	// a manifest with two tables of one id is refused
	m := r.Manifest
	m.Tables = append(m.Tables, &bdf.Table{ID: "chart-1", Part: m.Tables[0].Part})
	var b2 bytes.Buffer
	if err := d2.WriteSingle(&b2); err != nil {
		t.Fatal(err)
	}
	d2.Tables = append(d2.Tables, &bdf.Table{ID: "chart-1", Part: m.Tables[0].Part})
	b2.Reset()
	if err := d2.WriteSingle(&b2); err != nil {
		t.Fatal(err)
	}
	if _, err := bdf.OpenSingle(bytes.NewReader(b2.Bytes()), int64(b2.Len())); err == nil {
		t.Error("two tables of one id were accepted")
	}
}
