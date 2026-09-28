package parquet

import "testing"

// TestRowLimit checks the cell cap on the grid: the default and rows=all keep
// their behaviour for a few columns, but a table of many columns is held to
// maxSheetCells so the reader never allocates rows×cols cells from a row
// count the file only states.
func TestRowLimit(t *testing.T) {
	for _, c := range []struct {
		rows, cols int
		want       int64
	}{
		{0, 8, DefaultRows},                    // the default, few columns
		{0, 16384, maxCells / 16384},           // the default caps cells for many columns
		{-1, 8, MaxRows},                       // all rows, few columns: unchanged
		{-1, 16384, maxSheetCells / 16384},     // all rows, many columns: held to the cell cap
		{100000, 16384, maxSheetCells / 16384}, // an explicit count is capped too
		{500, 4, 500},                          // a modest explicit count passes through
		{0, 0, DefaultRows},                    // no columns: no division
	} {
		if got := rowLimit(c.rows, c.cols); got != c.want {
			t.Errorf("rowLimit(%d, %d) = %d, want %d", c.rows, c.cols, got, c.want)
		}
	}
	// the cap never lets rows×cols pass maxSheetCells
	for _, cols := range []int{1, 40, 1000, 16384} {
		if n := rowLimit(-1, cols) * int64(cols); n > maxSheetCells {
			t.Errorf("%d columns allow %d cells (cap %d)", cols, n, maxSheetCells)
		}
	}
}
