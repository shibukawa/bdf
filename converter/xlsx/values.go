package xlsx

import (
	"fmt"
	"io"
	"time"
)

// SheetValues is the values of the cells of a worksheet.
type SheetValues struct {
	// Name is the name of the sheet.
	Name string
	// Rows holds the values row by row from the first row of the sheet,
	// each from column A to its last value: nil for an empty cell, else a
	// float64, a string, a bool, or a time.Time for a number in a date
	// format. Formulas give the values the workbook was saved with.
	Rows [][]any
}

// maxValueCells bounds the cells, empty ones included, that Values returns
// for a workbook.
const maxValueCells = 4 << 20

// Values reads the values of the worksheets of a workbook (the hidden ones
// too), without converting it: a workbook read as data. Of the options,
// only Warn applies.
func Values(r io.ReaderAt, size int64, opts *Options) ([]SheetValues, error) {
	o := Options{NoSystemFonts: true, SystemFonts: true}
	if opts != nil {
		o.Warn = opts.Warn
	}
	c, _, err := openWorkbook(r, size, &o)
	if err != nil {
		return nil, err
	}
	var out []SheetValues
	room := maxValueCells
	for _, ref := range c.sheets {
		if ref.chart {
			continue
		}
		ws, err := c.readWorksheet(ref.part)
		if err != nil {
			return nil, fmt.Errorf("xlsx: sheet %q: %w", ref.name, err)
		}
		ws.sortRows()
		sv := SheetValues{Name: ref.name}
		for _, rw := range ws.rows {
			n := 0
			for _, cl := range rw.cells {
				if cl.kind != cellBlank {
					n = cl.col + 1
				}
			}
			if n == 0 {
				continue
			}
			if room -= rw.idx + 1 - len(sv.Rows) + n; room < 0 {
				return nil, fmt.Errorf("xlsx: sheet %q: its values reach over more than %d cells", ref.name, maxValueCells)
			}
			for len(sv.Rows) <= rw.idx {
				sv.Rows = append(sv.Rows, nil)
			}
			vals := make([]any, n)
			for _, cl := range rw.cells {
				if cl.col >= n {
					break
				}
				switch cl.kind {
				case cellNum:
					vals[cl.col] = cl.num
					if c.st.numFormat(c.st.xf(cl.style).numFmt).isDate() {
						if t, _, ok := serialTime(cl.num, c.date1904); ok {
							vals[cl.col] = time.Time(t)
						}
					}
				case cellStr:
					vals[cl.col] = cl.text.plain
				case cellBool:
					vals[cl.col] = cl.num != 0
				}
			}
			sv.Rows[rw.idx] = vals
		}
		out = append(out, sv)
	}
	return out, nil
}
