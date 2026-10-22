package parquet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"time"
)

// pageRows is the most rows of a column in one data page, and pageBytes
// the most bytes of their values.
const (
	pageRows  = 1 << 16
	pageBytes = 1 << 20
)

// Encode writes a table as a Parquet file in memory.
func Encode(t *Table) ([]byte, error) {
	var b bytes.Buffer
	if err := Write(&b, t); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Write writes a table as a Parquet file: its columns as optional fields
// of a flat schema, their values PLAIN, their definition levels RLE, in
// one row group, uncompressed. A value that is not of its column's type
// (see Table) is an error; a row shorter than the columns is padded with
// nulls, and a column of type Other cannot be written.
func Write(w io.Writer, t *Table) error {
	if len(t.Columns) == 0 {
		return fmt.Errorf("parquet: a table needs a column")
	}
	if len(t.Columns) > MaxColumns {
		return fmt.Errorf("parquet: %d columns, more than %d", len(t.Columns), MaxColumns)
	}
	for i, c := range t.Columns {
		if c.Name == "" {
			return fmt.Errorf("parquet: column %d has no name", i+1)
		}
		if c.Type == Other || int(c.Type) >= len(typeNames) {
			return fmt.Errorf("parquet: column %q: a column of type %s cannot be written", c.Name, c.Type)
		}
	}
	for i, r := range t.Rows {
		if len(r) > len(t.Columns) {
			return fmt.Errorf("parquet: row %d has %d values for %d columns", i+1, len(r), len(t.Columns))
		}
	}
	file := []byte("PAR1")
	type chunk struct {
		offset, size int64
		values       int64
	}
	chunks := make([]chunk, len(t.Columns))
	for j, c := range t.Columns {
		ch := chunk{offset: int64(len(file))}
		for start := 0; start < len(t.Rows) || start == 0; start += pageRows {
			end := min(start+pageRows, len(t.Rows))
			page, n, err := encodePage(c, t.Rows[start:end], j)
			if err != nil {
				return err
			}
			file = append(file, page...)
			ch.values += int64(n)
			if end == len(t.Rows) {
				break
			}
		}
		ch.size = int64(len(file)) - ch.offset
		chunks[j] = ch
	}
	// the footer
	e := newTenc()
	e.i32(1, 1) // version
	e.list(2, tStruct, len(t.Columns)+1)
	e.elem()
	e.bin(4, "schema")
	e.i32(5, int32(len(t.Columns)))
	e.end()
	for _, c := range t.Columns {
		e.elem()
		e.i32(1, physOf(c.Type))
		e.i32(3, repOptional)
		e.bin(4, c.Name)
		switch c.Type {
		case String:
			e.i32(6, convUTF8)
			e.begin(10)
			e.begin(logString)
			e.end()
			e.end()
		case Date:
			e.i32(6, convDate)
			e.begin(10)
			e.begin(logDate)
			e.end()
			e.end()
		case DateTime:
			e.i32(6, convTSMillis)
			e.begin(10)
			e.begin(logTimestamp)
			e.boolean(1, true) // isAdjustedToUTC
			e.begin(2)         // unit
			e.begin(unitMillis)
			e.end()
			e.end()
			e.end()
			e.end()
		}
		e.end()
	}
	e.i64(3, int64(len(t.Rows)))
	e.list(4, tStruct, 1) // row groups
	e.elem()
	e.list(1, tStruct, len(t.Columns))
	var total int64
	for j, c := range t.Columns {
		ch := chunks[j]
		total += ch.size
		e.elem()
		e.i64(2, ch.offset)
		e.begin(3)
		e.i32(1, physOf(c.Type))
		e.list(2, tI32, 2)
		e.b = binary.AppendVarint(e.b, encPlain)
		e.b = binary.AppendVarint(e.b, encRLE)
		e.list(3, tBinary, 1)
		e.b = binary.AppendUvarint(e.b, uint64(len(c.Name)))
		e.b = append(e.b, c.Name...)
		e.i32(4, codecNone)
		e.i64(5, ch.values)
		e.i64(6, ch.size)
		e.i64(7, ch.size)
		e.i64(9, ch.offset)
		e.end()
		e.end()
	}
	e.i64(2, total)
	e.i64(3, int64(len(t.Rows)))
	e.end()
	e.bin(6, "bdf")
	footer := e.bytes()
	file = append(file, footer...)
	file = binary.LittleEndian.AppendUint32(file, uint32(len(footer)))
	file = append(file, "PAR1"...)
	_, err := w.Write(file)
	return err
}

// physOf is the physical type of a column.
func physOf(t Type) int32 {
	switch t {
	case String:
		return typeBinary
	case Number:
		return typeDouble
	case Integer, DateTime:
		return typeInt64
	case Boolean:
		return typeBoolean
	case Date:
		return typeInt32
	}
	return typeBinary
}

// encodePage writes a data page (version 1) of the values of column j of
// rows: their definition levels (1 for a value, 0 for a null), RLE, then
// the values, PLAIN. It returns the page and the count of its entries.
func encodePage(c Column, rows [][]any, j int) ([]byte, int, error) {
	var vals []byte
	var defs []byte // the definition levels, a run each change
	run, level := 0, byte(1)
	flush := func() {
		if run > 0 {
			defs = binary.AppendUvarint(defs, uint64(run)<<1)
			defs = append(defs, level)
		}
	}
	var bits byte
	nbits := 0
	for i, r := range rows {
		var v any
		if j < len(r) {
			v = r[j]
		}
		l := byte(1)
		if v == nil {
			l = 0
		}
		if l != level {
			flush()
			run, level = 0, l
		}
		run++
		if v == nil {
			continue
		}
		var err error
		switch c.Type {
		case String:
			var s string
			switch x := v.(type) {
			case string:
				s = x
			case []byte:
				s = string(x)
			case fmt.Stringer:
				s = x.String()
			default:
				s = fmt.Sprint(x)
			}
			vals = binary.LittleEndian.AppendUint32(vals, uint32(len(s)))
			vals = append(vals, s...)
		case Number:
			f, ok := toFloat(v)
			if !ok {
				err = fmt.Errorf("%T is not a number", v)
			}
			vals = binary.LittleEndian.AppendUint64(vals, math.Float64bits(f))
		case Integer:
			n, ok := toInt(v)
			if !ok {
				err = fmt.Errorf("%T is not an integer", v)
			}
			vals = binary.LittleEndian.AppendUint64(vals, uint64(n))
		case Boolean:
			b, ok := v.(bool)
			if !ok {
				err = fmt.Errorf("%T is not a boolean", v)
			}
			if b {
				bits |= 1 << nbits
			}
			if nbits++; nbits == 8 {
				vals = append(vals, bits)
				bits, nbits = 0, 0
			}
		case Date, DateTime:
			tm, ok := v.(time.Time)
			if !ok {
				err = fmt.Errorf("%T is not a time", v)
			}
			if c.Type == Date {
				y, m, d := tm.UTC().Date()
				days := time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Sub(unixDay) / (24 * time.Hour)
				vals = binary.LittleEndian.AppendUint32(vals, uint32(int32(days)))
			} else {
				vals = binary.LittleEndian.AppendUint64(vals, uint64(tm.UnixMilli()))
			}
		}
		if err != nil {
			return nil, 0, fmt.Errorf("parquet: column %q, row %d: %v", c.Name, i+1, err)
		}
		if len(vals) >= pageBytes && i+1 < len(rows) {
			// (the page is large enough: the rest go in the next pages)
			rows = rows[:i+1]
			break
		}
	}
	if nbits > 0 {
		vals = append(vals, bits)
	}
	flush()
	body := binary.LittleEndian.AppendUint32(nil, uint32(len(defs)))
	body = append(body, defs...)
	body = append(body, vals...)
	h := newTenc()
	h.i32(1, pageData)
	h.i32(2, int32(len(body)))
	h.i32(3, int32(len(body)))
	h.begin(5)
	h.i32(1, int32(len(rows)))
	h.i32(2, encPlain)
	h.i32(3, encRLE)
	h.i32(4, encRLE)
	h.end()
	return append(h.bytes(), body...), len(rows), nil
}

func toFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	case int32:
		return float64(x), true
	case uint:
		return float64(x), true
	case uint64:
		return float64(x), true
	case uint32:
		return float64(x), true
	}
	return 0, false
}

func toInt(v any) (int64, bool) {
	switch x := v.(type) {
	case int64:
		return x, true
	case int:
		return int64(x), true
	case int32:
		return int64(x), true
	case int16:
		return int64(x), true
	case int8:
		return int64(x), true
	case uint:
		return int64(x), true
	case uint32:
		return int64(x), true
	case uint64:
		return int64(x), x <= math.MaxInt64
	case float64:
		return int64(x), x == math.Trunc(x) && math.Abs(x) < 1<<53
	}
	return 0, false
}
