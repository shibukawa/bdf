package parquet

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

// Text writes a value as text: a date as 2006-01-02, a date and time in
// RFC 3339, a number as Go prints it, a null as "".
func Text(c Column, v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case float64:
		return strconv.FormatFloat(x, 'g', -1, 64)
	case int64:
		return strconv.FormatInt(x, 10)
	case bool:
		return strconv.FormatBool(x)
	case time.Time:
		if c.Type == Date {
			return x.UTC().Format("2006-01-02")
		}
		return x.UTC().Format(time.RFC3339Nano)
	case []byte:
		return fmt.Sprintf("%x", x)
	}
	return fmt.Sprint(v)
}

// WriteCSV writes a table as CSV: the names of the columns, then the rows,
// their values as Text writes them.
func WriteCSV(w io.Writer, t *Table) error {
	cw := csv.NewWriter(w)
	names := make([]string, len(t.Columns))
	for i, c := range t.Columns {
		names[i] = c.Name
	}
	if err := cw.Write(names); err != nil {
		return err
	}
	row := make([]string, len(t.Columns))
	for _, r := range t.Rows {
		for i, c := range t.Columns {
			var v any
			if i < len(r) {
				v = r[i]
			}
			row[i] = Text(c, v)
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

// WriteJSON writes a table as JSON: an array of objects, a value under
// the name of its column (dates and times as strings, as Text writes
// them; nulls as null; bytes as hexadecimal).
func WriteJSON(w io.Writer, t *Table) error {
	out := make([]map[string]any, len(t.Rows))
	for i, r := range t.Rows {
		o := make(map[string]any, len(t.Columns))
		for j, c := range t.Columns {
			var v any
			if j < len(r) {
				v = r[j]
			}
			switch x := v.(type) {
			case time.Time, []byte:
				v = Text(c, x)
			}
			o[c.Name] = v
		}
		out[i] = o
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(out)
}
