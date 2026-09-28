package csv

import (
	"bytes"
	"strings"
	"testing"
)

// TestNewlineField is a smoke test that a quoted field of many line breaks
// converts quickly: the layout bounds the lines it makes for such a cell (see
// the xlsx package), so the whole file is still two records and conversion
// does not hang.
func TestNewlineField(t *testing.T) {
	data := []byte("a,b\n1,\"" + strings.Repeat("\n", 20000) + "\"\n")
	res, err := Convert(bytes.NewReader(data), int64(len(data)), testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if res.Rows != 2 || res.Cols != 2 {
		t.Errorf("rows %d cols %d, want 2×2", res.Rows, res.Cols)
	}
}
