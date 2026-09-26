package xlsx

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// testOptions restricts fonts to the test fonts (shared with the PowerPoint
// tests) so that output does not depend on the machine.
func testOptions() *Options {
	return &Options{FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true}
}

func convert(t *testing.T, name string, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
	res, err := ConvertFile(filepath.Join("testdata", name), opts)
	if err != nil {
		t.Fatal(err)
	}
	return res, reopen(t, res)
}

// convertEdited converts a test workbook with some of its parts rewritten.
func convertEdited(t *testing.T, name string, edits map[string]func(string) string, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	zw := zip.NewWriter(&out)
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		b, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Fatal(err)
		}
		if edit, ok := edits[f.Name]; ok {
			s := edit(string(b))
			if s == string(b) {
				t.Fatalf("edit of %s changed nothing", f.Name)
			}
			b = []byte(s)
		}
		w, _ := zw.Create(f.Name)
		w.Write(b)
	}
	zw.Close()
	res, err := Convert(bytes.NewReader(out.Bytes()), int64(out.Len()), opts)
	if err != nil {
		t.Fatal(err)
	}
	return res, reopen(t, res)
}

func reopen(t *testing.T, res *Result) *bdf.Reader {
	t.Helper()
	var buf bytes.Buffer
	if err := res.Doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func indexRuns(t *testing.T, r *bdf.Reader, v *bdf.View) []bdf.IndexRun {
	t.Helper()
	h, err := bdf.ParseHash(v.TextIndex)
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.Part(h)
	if err != nil {
		t.Fatal(err)
	}
	idx, err := bdf.DecodeTextIndex(b)
	if err != nil {
		t.Fatal(err)
	}
	return idx
}

// marks lists the MARKs of an object (and the objects it uses) with their
// payloads, and counts its instructions.
func marks(t *testing.T, r *bdf.Reader, h bdf.Hash) ([]string, map[byte]int) {
	t.Helper()
	var out []string
	counts := map[byte]int{}
	var walk func(h bdf.Hash)
	walk = func(h bdf.Hash) {
		o, err := r.Object(h)
		if err != nil {
			t.Fatal(err)
		}
		o.Walk(func(in bdf.Instr) {
			counts[in.Op]++
			switch in.Op {
			case bdf.OpMark:
				out = append(out, strings.TrimSpace(string("PLCBAWHlLTFEG"[in.Args[0].(uint64)])+" "+in.Args[1].(string)))
			case bdf.OpUse, bdf.OpUseAt:
				walk(o.Objects[int(in.Args[0].(uint64))])
			}
		})
	}
	walk(h)
	return out, counts
}

func tile(t *testing.T, v *bdf.View, key string) bdf.Hash {
	t.Helper()
	h, err := bdf.ParseHash(v.Tiles[key])
	if err != nil {
		t.Fatalf("tile %s: %v", key, err)
	}
	return h
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

func TestConvertBasic(t *testing.T) {
	res, r := convert(t, "basic.xlsx", testOptions())
	if res.Sheets != 2 {
		t.Fatalf("sheets = %d (the hidden sheet must be skipped)", res.Sheets)
	}
	if res.EmbeddedFonts != 2 {
		t.Errorf("embedded fonts = %d", res.EmbeddedFonts)
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "2 formula cell(s)") {
		t.Errorf("warnings: %q", res.Warnings)
	}
	m := r.Manifest
	if dc := m.Meta.DC; dc.Title.First() != "Excel test workbook" || dc.Creator.First() != "bdf test" ||
		dc.Created.First() != "2026-09-01T09:00:00Z" || m.Meta.Source != "xlsx" {
		t.Errorf("meta = %+v", m.Meta)
	}
	if len(m.Views) != 2 || m.Views[0].Title != "Sales" || m.Views[1].Title != "Big" || m.Views[1].ID != "sheet3" {
		t.Fatalf("views = %+v", m.Views)
	}
	v := m.Views[0]
	if v.Kind != bdf.ViewSheet || v.Tile != 2048 || !v.Gridlines || v.Freeze == nil || v.Freeze.Rows != 3 || v.Freeze.Cols != 0 {
		t.Errorf("sheet view = %s tile %g grid %v freeze %+v", v.Kind, v.Tile, v.Gridlines, v.Freeze)
	}
	// columns: A 14 chars, then B and C 12, D 13, E 18, F 8, the hidden H
	var widths []float32
	for _, run := range v.Cols {
		for i := 0; i < int(run[0]); i++ {
			widths = append(widths, run[1])
		}
	}
	if len(widths) < 12 || widths[7] != 0 || widths[1] != widths[2] || !(widths[4] > widths[3] && widths[3] > widths[1] && widths[1] > widths[5]) {
		t.Errorf("column widths = %v", widths)
	}
	var heights []float32
	for _, run := range v.Rows {
		for i := 0; i < int(run[0]); i++ {
			heights = append(heights, run[1])
		}
	}
	if heights[0] != 28 || heights[12] != 62 || heights[20] != 0 || heights[19] <= heights[18] {
		t.Errorf("row heights = %v (row 20 holds 20 pt text, row 21 is hidden)", heights[:22])
	}
	text := bdf.PlainText(indexRuns(t, r, v))
	for _, want := range []string{
		"Quarterly sales 数字の表", "Region\nQ1\nQ2\nGrowth\nUpdated\nFlag",
		"North ノース\n125,000\n131,250.50\n5.0%\n2026-03-31\nTRUE",
		"South サウス\n98,000\n-1,234.50\n-7.0%", "1,234,568\n1,300,000.00", "2026/12/31 18:30",
		"This text is long and flows over the empty cells to its right.",
		// wrapped Japanese lines join without a space, Latin ones with one
		"日本語の文章は、単語の区切りに空白を使わないため、文字と文字の間で改行します。句読点「、」や「。」は行頭に来ません。",
		"Wrapped Latin text breaks after spaces between words.",
		"indent 2", "shrink to fit this text", "centered across", "rotated 45°", "up 90°", "縦書き",
		"######", "3.14159", "Rich bold red and italic", "Strike", "x2",
		"1-May-26\nFriday", "3/8\n1.23E+04", "$ 1,234.50\n$ (1,234.50)\n#N/A",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("text lacks %q", want)
		}
	}
	if strings.Contains(text, "hidden row") || strings.Contains(text, "Hidden sheet") {
		t.Error("hidden content converted")
	}
	ms, counts := marks(t, r, tile(t, v, "0,0"))
	for _, want := range []string{"C A1:F1", "C A13:C13", "C B4", "G ja", "G"} {
		if !contains(ms, want) {
			t.Errorf("marks lack %q: %v", want, ms[:min(len(ms), 20)])
		}
	}
	if counts[bdf.OpLink] != 1 {
		t.Errorf("links = %d", counts[bdf.OpLink])
	}
	if counts[bdf.OpStrokePath] == 0 || counts[bdf.OpFillPaint] < 3 {
		t.Errorf("borders %d, paints %d (pattern and gradient fills)", counts[bdf.OpStrokePath], counts[bdf.OpFillPaint])
	}
}

func TestTiles(t *testing.T) {
	_, r := convert(t, "basic.xlsx", testOptions())
	v := r.Manifest.Views[1]
	for _, k := range []string{"0,0", "0,1", "0,2", "1,0"} {
		if v.Tiles[k] == "" {
			t.Errorf("no tile %s: %v", k, v.Tiles)
		}
	}
	// Text over a tile edge is drawn in both tiles but indexed once.
	runs := indexRuns(t, r, v)
	count := func(s string) int {
		n := 0
		for _, run := range runs {
			if strings.Contains(run.Text, s) {
				n++
			}
		}
		return n
	}
	if n := count("straddles the tile boundary"); n != 1 {
		t.Errorf("straddling text indexed %d times", n)
	}
	if n := count("flows across the tile boundary"); n != 1 {
		t.Errorf("overflowing text indexed %d times", n)
	}
	for _, key := range []string{"0,0", "1,0"} {
		o, err := r.Object(tile(t, v, key))
		if err != nil {
			t.Fatal(err)
		}
		found := false
		o.Walk(func(in bdf.Instr) {
			if in.Op == bdf.OpFillText && strings.Contains(in.Args[0].(string), "flows across") {
				found = true
			}
		})
		if !found {
			t.Errorf("tile %s does not draw the overflowing text", key)
		}
	}
	if text := bdf.PlainText(runs); !strings.Contains(text, "row 300\n90000") || !strings.Contains(text, "row 137\n18769\nstraddles") {
		t.Errorf("big sheet text: %q …", text[:200])
	}
}

func TestOptions(t *testing.T) {
	opts := testOptions()
	opts.Hidden = true
	res, r := convert(t, "basic.xlsx", opts)
	if res.Sheets != 3 || r.Manifest.Views[1].Title != "Hidden" {
		t.Errorf("Hidden: %d sheets", res.Sheets)
	}
	opts = testOptions()
	opts.Sheets = []int{3}
	if res, _ := convert(t, "basic.xlsx", opts); res.Sheets != 1 || res.Doc.Views[0].Title != "Big" {
		t.Errorf("Sheets: %d", res.Sheets)
	}
	opts.Sheets = []int{9}
	if _, err := ConvertFile("testdata/basic.xlsx", opts); err == nil {
		t.Error("sheet 9 of 3 converted")
	}
	opts = testOptions()
	opts.SystemFonts = true
	if res, _ := convert(t, "basic.xlsx", opts); res.EmbeddedFonts != 0 {
		t.Errorf("SystemFonts embedded %d fonts", res.EmbeddedFonts)
	}
}

func TestDeterministic(t *testing.T) {
	var out [][]byte
	for i := 0; i < 2; i++ {
		for _, name := range []string{"basic.xlsx", "features.xlsx"} {
			res, err := ConvertFile(filepath.Join("testdata", name), testOptions())
			if err != nil {
				t.Fatal(err)
			}
			var buf bytes.Buffer
			if err := res.Doc.WriteSingle(&buf); err != nil {
				t.Fatal(err)
			}
			out = append(out, buf.Bytes())
		}
	}
	if !bytes.Equal(out[0], out[2]) || !bytes.Equal(out[1], out[3]) {
		t.Error("conversion is not deterministic")
	}
}

func TestFeatures(t *testing.T) {
	res, r := convert(t, "features.xlsx", testOptions())
	if res.Sheets != 2 {
		t.Fatalf("sheets = %d", res.Sheets)
	}
	v := r.Manifest.Views[0]
	text := bdf.PlainText(indexRuns(t, r, v))
	for _, want := range []string{"Conditional formats", "18\nbanana", "Item\nUnits\nPrice\nPens\n120\n1.5", "図形のテキスト shape text", "Units"} {
		if !strings.Contains(text, want) {
			t.Errorf("text lacks %q", want)
		}
	}
	ms, counts := marks(t, r, tile(t, v, "0,0"))
	// table header cells are column headers; the picture and the shape
	// are figures with their descriptions
	for _, want := range []string{"C A16 col", "C C16 col", "C A17", "F Picture", "F A rounded note box"} {
		if !contains(ms, want) {
			t.Errorf("marks lack %q: %v", want, ms)
		}
	}
	if counts[bdf.OpImage] != 1 || counts[bdf.OpUseAt] < 3 {
		t.Errorf("images %d, uses %d", counts[bdf.OpImage], counts[bdf.OpUseAt])
	}
	// the chart has the values of the cells it refers to
	chart := r.Manifest.Views[1]
	if chart.Kind != bdf.ViewFixed || chart.Title != "Chart" || len(chart.Pages) != 1 || chart.Pages[0].W != 720 {
		t.Fatalf("chart sheet = %+v", chart)
	}
	if text := bdf.PlainText(indexRuns(t, r, chart)); !strings.Contains(text, "Values") || !strings.Contains(text, "scale") || !strings.Contains(text, "40") {
		t.Errorf("chart sheet text %q", text)
	}
}

func TestConditionalFormats(t *testing.T) {
	c := &converter{}
	s := &sheetCtx{c: c, ws: &worksheet{}, nRows: 20, nCols: 5}
	vals := []cellVal{}
	for i, v := range []float64{3, 18, 7, 42, 25} {
		vals = append(vals, cellVal{r: i, c: 0, kind: cellNum, num: v})
	}
	nums := numbers(vals)
	for _, tc := range []struct {
		op   string
		args []float64
		v    float64
		want bool
	}{
		{"greaterThan", []float64{20}, 25, true},
		{"between", []float64{10, 5}, 7, true},
		{"notBetween", []float64{5, 10}, 7, false},
		{"lessThanOrEqual", []float64{3}, 3, true},
	} {
		var args []cellVal
		for _, a := range tc.args {
			args = append(args, cellVal{kind: cellNum, num: a})
		}
		if got := compare(tc.op, cellVal{kind: cellNum, num: tc.v}, args); got != tc.want {
			t.Errorf("%s %v %v = %v", tc.op, tc.v, tc.args, got)
		}
	}
	if p := percentile(nums, 50); p != 18 {
		t.Errorf("median = %v", p)
	}
	if m, sd := meanStd(nums); m != 19 || sd < 15 || sd > 16 {
		t.Errorf("mean %v sd %v", m, sd)
	}
	// operands: numbers, strings, relative and absolute references
	s.ws.rows = []row{{idx: 0, cells: []cell{{col: 1, kind: cellNum, num: 5}}}, {idx: 1, cells: []cell{{col: 1, kind: cellNum, num: 9}}}}
	for _, tc := range []struct {
		f    string
		r    int
		want float64
	}{{"B1", 0, 5}, {"B1", 1, 9}, {"$B$1", 1, 5}, {"-B$1", 1, -5}, {"12.5", 0, 12.5}} {
		v, ok := s.evalOperand(tc.f, [2]int{0, 0}, tc.r, 0)
		if !ok || v.num != tc.want {
			t.Errorf("operand %s at row %d = %v %v", tc.f, tc.r, v.num, ok)
		}
	}
	if v, ok := s.evalOperand(`"a""b"`, [2]int{}, 0, 0); !ok || v.text != `a"b` {
		t.Errorf("string operand = %q", v.text)
	}
	if _, ok := s.evalOperand("SUM(A1:A3)", [2]int{}, 0, 0); ok {
		t.Error("a formula was evaluated")
	}
}

func TestGeometry(t *testing.T) {
	a := newAxis(10, 15, map[int]float64{2: 30, 5: 0})
	if a.pos(0) != 0 || a.pos(3) != 60 || a.pos(6) != 90 || a.total() != 150 || a.at(5) != 0 {
		t.Errorf("axis pos %v %v %v total %v", a.pos(0), a.pos(3), a.pos(6), a.total())
	}
	if a.index(59) != 2 || a.index(60) != 3 || a.index(1000) != 9 {
		t.Errorf("axis index %v %v %v", a.index(59), a.index(60), a.index(1000))
	}
	if got := a.runs(); len(got) != 5 || got[0] != (bdf.Run{2, 15}) || got[1] != (bdf.Run{1, 30}) || got[3] != (bdf.Run{1, 0}) {
		t.Errorf("runs %v", got)
	}
	// Calibri 11 (7 px digits): the standard width of 8.43 characters is 64 px
	if w := colWidthPt(9.140625, 7); w != 48 {
		t.Errorf("column width = %v", w)
	}
	ws := &worksheet{baseColW: 8}
	if w := ws.defaultColPt(7); w != 48 {
		t.Errorf("default column = %v", w)
	}
	if w := ws.defaultColPt(8); w != 54 {
		t.Errorf("default column (8 px digits) = %v", w)
	}
	if h := autoRowHeight(13.43); h != 15 {
		t.Errorf("row height for Calibri 11 = %v", h)
	}
	for _, c := range []struct {
		ref      string
		col, row int
	}{{"A1", 0, 0}, {"Z10", 25, 9}, {"AA1", 26, 0}, {"$XFD$1048576", 16383, 1048575}} {
		col, row, ok := parseRef(c.ref)
		if !ok || col != c.col || row != c.row {
			t.Errorf("parseRef(%s) = %d %d %v", c.ref, col, row, ok)
		}
		if got := colName(c.col); !strings.Contains(c.ref, got) {
			t.Errorf("colName(%d) = %s", c.col, got)
		}
	}
	if r, ok := parseRange("B2:A1"); !ok || r != (cellRange{0, 0, 1, 1}) {
		t.Errorf("range %v", r)
	}
	if r, ok := parseRange("C:D"); !ok || r.c0 != 2 || r.c1 != 3 || r.r1 != maxRows-1 {
		t.Errorf("column range %v", r)
	}
	if got := decodeXString("a_x000D_\r\nb_x0041_"); got != "a\r\nbA" && got != "a\n\nbA" {
		t.Errorf("xstring %q", got)
	}
}

func TestEdits(t *testing.T) {
	// right-to-left sheets warn; a sheet view without gridlines has none
	res, r := convertEdited(t, "basic.xlsx", map[string]func(string) string{
		"xl/worksheets/sheet1.xml": func(s string) string {
			return strings.Replace(s, "<sheetView ", `<sheetView rightToLeft="1" showGridLines="0" `, 1)
		},
	}, testOptions())
	if r.Manifest.Views[0].Gridlines {
		t.Error("gridlines shown")
	}
	found := false
	for _, w := range res.Warnings {
		found = found || strings.Contains(w, "right-to-left")
	}
	if !found {
		t.Errorf("warnings %q", res.Warnings)
	}
}

func TestFormulas(t *testing.T) {
	var warnings []string
	c := &converter{parsed: map[string]*worksheet{}, warned: map[string]bool{}, opts: &Options{Warn: func(w string) { warnings = append(warnings, w) }}}
	ws := &worksheet{name: "S", rows: []row{
		{idx: 0, cells: []cell{{col: 0, kind: cellNum, num: 4}, {col: 1, kind: cellStr, text: &richText{plain: "done"}}, {col: 2, kind: cellNum, num: 10}}},
		{idx: 1, cells: []cell{{col: 0, kind: cellNum, num: 7}, {col: 1, kind: cellStr, text: &richText{plain: "Open"}}, {col: 2, kind: cellNum, num: 4}}},
		{idx: 2, cells: []cell{{col: 0, kind: cellNum, num: 4}, {col: 2, kind: cellErr, text: &richText{plain: "#N/A"}}}},
	}}
	s := &sheetCtx{c: c, ws: ws, formulas: map[string]fnode{}}
	for _, tc := range []struct {
		f    string
		r, c int
		want string
	}{
		{"MOD(ROW(),2)=0", 1, 0, "TRUE"},
		{"ISEVEN(ROW())", 0, 0, "FALSE"},
		{"$B1=\"DONE\"", 0, 0, "TRUE"},
		{"$B1=\"DONE\"", 1, 0, "FALSE"},
		{"A1>AVERAGE($A$1:$A$3)", 1, 0, "TRUE"},
		{"COUNTIF($A:$A,A1)>1", 2, 0, "TRUE"},
		{"COUNTIF($A:$A,A1)>1", 1, 0, "FALSE"},
		{"AND(A1>3,OR(B1=\"x\",C1>=10))", 0, 0, "TRUE"},
		{"ISNA(C1)", 2, 0, "TRUE"},
		{"IFERROR(C1*2,-1)", 2, 0, "-1"},
		{"LEFT(B1,2)&UPPER(\"k\")", 0, 0, "doK"},
		{"-A1^2+SUM(C1:C2)%", 0, 0, "16.14"}, // unary minus binds before ^,
		{"SEARCH(\"PE\",B1)", 1, 0, "2"},
		{"COUNTIF(A1:A3,\"<5\")", 0, 0, "2"},
		{"COUNTIF(B1:B3,\"d*\")", 0, 0, "1"},
		{"MAX(A1:A3)-MIN(A1:A3)", 0, 0, "3"},
	} {
		v, ok := s.evalFormula(tc.f, [2]int{0, 0}, tc.r, tc.c)
		if got := toStr(v); !ok || got != tc.want {
			t.Errorf("%s at %s = %q %v, want %q", tc.f, cellRef(tc.r, tc.c), got, ok, tc.want)
		}
	}
	if len(warnings) != 0 {
		t.Errorf("warnings %q", warnings)
	}
	for _, f := range []string{"TODAY()>A1", "A1+", "INDIRECT(\"A1\")"} {
		if _, ok := s.evalFormula(f, [2]int{}, 0, 0); ok {
			t.Errorf("%s evaluated", f)
		}
	}
}
