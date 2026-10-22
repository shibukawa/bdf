//go:build js && wasm && previewonly

package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/encoding/parquet"
)

func init() {
	methods["tables"] = tablesCall
	methods["table"] = tableCall
}

// tablesCall lists the tables a single-file bdf holds beside its drawing
// (docs/spec.md §4.6), as bdf tables does.
func tablesCall(_ js.Value, args []js.Value) any {
	data, err := bdfInput(args)
	if err != nil {
		return reject(fmt.Errorf("tables: %w", err))
	}
	password := str(args, "password")
	return promise(func() (any, error) {
		r, err := bdf.OpenSingle(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return nil, err
		}
		if r.Locked() && password != "" {
			if err := r.Unlock(password); err != nil {
				return nil, err
			}
		}
		if r.Locked() {
			return nil, bdf.ErrLocked
		}
		out := make([]any, 0, len(r.Manifest.Tables))
		for _, t := range r.Manifest.Tables {
			out = append(out, tableValue(t))
		}
		return out, nil
	})
}

func tableValue(t *bdf.Table) map[string]any {
	cols := make([]any, len(t.Columns))
	for i, c := range t.Columns {
		cols[i] = map[string]any{"name": c.Name, "type": c.Type}
	}
	return map[string]any{"id": t.ID, "title": t.Title, "part": t.Part, "rows": t.Rows, "columns": cols, "view": t.View, "page": t.Page}
}

// tableCall writes a table of a single-file bdf as CSV, JSON or the
// Parquet file the document holds, as bdf table does.
func tableCall(_ js.Value, args []js.Value) any {
	data, err := bdfInput(args)
	if err != nil {
		return reject(fmt.Errorf("table: %w", err))
	}
	password, id, format := str(args, "password"), str(args, "id"), str(args, "format")
	if format == "" {
		format = "csv"
	}
	switch format {
	case "csv", "json", "parquet":
	default:
		return reject(fmt.Errorf("table: format must be csv, json or parquet, not %q", format))
	}
	rows := 0
	if v := option(args, "rows"); v.Type() == js.TypeNumber {
		rows = v.Int()
	}
	return promise(func() (any, error) {
		doc, err := readDocument(data, password)
		if err != nil {
			return nil, err
		}
		t := doc.Table(id)
		if t == nil {
			return nil, fmt.Errorf("no table %q", id)
		}
		h, err := bdf.ParseHash(t.Part)
		if err != nil {
			return nil, err
		}
		part := doc.Part(h)
		if part == nil {
			return nil, fmt.Errorf("table %s: its part is missing", id)
		}
		out := map[string]any{"table": tableValue(t), "format": format}
		if format == "parquet" {
			out["data"] = bytesToJS(part.Data)
			return out, nil
		}
		tb, err := parquet.Read(part.Data, parquet.Limits{MaxRows: rows})
		if err != nil {
			return nil, err
		}
		var b bytes.Buffer
		if format == "csv" {
			err = parquet.WriteCSV(&b, tb)
		} else {
			err = parquet.WriteJSON(&b, tb)
		}
		if err != nil {
			return nil, err
		}
		out["data"] = bytesToJS(b.Bytes())
		out["rowsRead"] = len(tb.Rows)
		return out, nil
	})
}
