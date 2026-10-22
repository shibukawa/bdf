package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTables checks bdf tables and bdf table on the tables of the demo
// document and of a converted presentation.
func TestTables(t *testing.T) {
	demo := filepath.Join("..", "..", "testdata", "demo.bdf")
	code, out, errOut := runOut(t, "tables", demo)
	if code != 0 || !strings.Contains(out, `table-1: 3 rows × 3 columns "Section 1" (doc, page 1)`) || !strings.Contains(out, "R1C1:string, R1C2:integer, R1C3:number") {
		t.Errorf("tables: exit %d\n%s%s", code, out, errOut)
	}
	code, out, _ = runOut(t, "tables", "-json", demo)
	if code != 0 || !strings.Contains(out, `"id": "table-1"`) || !strings.Contains(out, `"type": "integer"`) {
		t.Errorf("tables -json: exit %d\n%s", code, out)
	}
	// the rows as CSV (the default), as JSON, and as the Parquet part
	code, out, _ = runOut(t, "table", demo, "table-1")
	if code != 0 || out != "R1C1,R1C2,R1C3\nR2C1,2,1\nR3C1,3,1.5\nR4C1,4,2\n" {
		t.Errorf("table: exit %d\n%q", code, out)
	}
	code, out, _ = runOut(t, "table", "-format", "json", "-rows", "1", demo, "table-1")
	if code != 0 || !strings.Contains(out, `{"R1C1":"R2C1","R1C2":2,"R1C3":1}`) || strings.Contains(out, "R3C1") {
		t.Errorf("table -format json -rows 1: exit %d\n%s", code, out)
	}
	dir := t.TempDir()
	pq := filepath.Join(dir, "t.parquet")
	if code, _, errOut = runOut(t, "table", "-format", "parquet", demo, "table-1", pq); code != 0 {
		t.Fatalf("table -format parquet: exit %d\n%s", code, errOut)
	}
	if data, err := os.ReadFile(pq); err != nil || !strings.HasPrefix(string(data), "PAR1") || !strings.HasSuffix(string(data), "PAR1") {
		t.Errorf("parquet file: %v, %d bytes", err, len(data))
	}
	// an unknown table, and a document without tables
	if code, _, errOut = runOut(t, "table", demo, "nope"); code == 0 || !strings.Contains(errOut, `no table "nope"`) {
		t.Errorf("unknown table: exit %d\n%s", code, errOut)
	}
	code, out, _ = runOut(t, "tables", filepath.Join("..", "..", "testdata", "demo-encrypted.bdf"))
	if code == 0 {
		t.Errorf("encrypted without a password: exit %d\n%s", code, out)
	}
}
