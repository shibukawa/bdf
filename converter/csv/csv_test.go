package csv

import (
	"bytes"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/xlsx"
)

// testOptions lays text out with the test fonts of the PowerPoint tests
// only, so that the output does not depend on the machine.
func testOptions() *Options {
	return &Options{FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true}
}

func reopen(t *testing.T, doc *bdf.Document) *bdf.Reader {
	t.Helper()
	var buf bytes.Buffer
	if err := doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func sheetText(t *testing.T, r *bdf.Reader) string {
	t.Helper()
	v := r.Manifest.Views[0]
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
	return bdf.PlainText(idx)
}

func TestConvertFiles(t *testing.T) {
	for _, c := range []struct {
		file    string
		charset string
		bom     bool
		dialect Dialect
		rows    int
		cols    int
		text    []string
	}{
		{"basic.csv", "UTF-8", true, Dialect{',', '"', 0}, 11, 9, []string{
			"ID\nProduct\nCategory\nPrice\nQuantity\nDiscount\nReleased\nCode\nNote",
			"Plain, A5", `Nib sizes "F" and "M"`, "Two lines: warm and cool", "$1,249.99", "(10%)", "0107"}},
		{"japanese.tsv", "Shift_JIS", false, Dialect{'\t', '"', 0}, 7, 5, []string{
			"項目\n文字数\n日\nタイム\n説明", "縦書きの文章", "R6.1.2", "一行目 二行目", "English と日本語の混在"}},
	} {
		res, err := ConvertFile("testdata/"+c.file, testOptions())
		if err != nil {
			t.Fatalf("%s: %v", c.file, err)
		}
		if res.Charset != c.charset || res.BOM != c.bom || res.Dialect != c.dialect || !res.Header || res.Rows != c.rows || res.Cols != c.cols {
			t.Errorf("%s: %s bom %v %q header %v %d×%d", c.file, res.Charset, res.BOM, res.Dialect, res.Header, res.Rows, res.Cols)
		}
		if len(res.Warnings) != 0 || res.EmbeddedFonts != 2 {
			t.Errorf("%s: warnings %q, %d font(s)", c.file, res.Warnings, res.EmbeddedFonts)
		}
		r := reopen(t, res.Doc)
		m := r.Manifest
		name := strings.TrimSuffix(c.file, c.file[strings.LastIndexByte(c.file, '.'):])
		if m.Meta.Source != "csv" || len(m.Views) != 1 || m.Views[0].Title != name || m.Views[0].Freeze == nil || m.Views[0].Freeze.Rows != 1 {
			t.Errorf("%s: meta %+v, views %+v", c.file, m.Meta, m.Views)
		}
		text := sheetText(t, r)
		for _, want := range c.text {
			if !strings.Contains(text, want) {
				t.Errorf("%s: text index lacks %q:\n%s", c.file, want, text)
			}
		}
	}
}

func TestCodeColumns(t *testing.T) {
	recs := [][]string{{"Code", "N", "Mixed"}, {"0012", "12", "0012"}, {"1003", "1003", "abc"}, {"", "7", "1003"}}
	build := func(header bool) [][]xlsx.GridCell {
		var rows [][]xlsx.GridCell
		first, rest := make([]digitStat, 3), make([]digitStat, 3)
		for i, rec := range recs {
			stats := rest
			if i == 0 {
				stats = first
			}
			var cells []xlsx.GridCell
			for j, v := range rec {
				cells = append(cells, xlsx.GridCell{Text: v, Number: classify(v).numeric()})
				stats[j].add(v)
			}
			rows = append(rows, cells)
		}
		codeColumns(rows, first, rest, header)
		return rows
	}
	rows := build(true)
	// codes: 1003 becomes text like 0012; plain numbers and mixed columns stay
	if rows[2][0].Number || !rows[2][1].Number || !rows[3][2].Number {
		t.Errorf("cells = %+v", rows)
	}
	// without a header row, "Code" is a value other than digits
	if rows = build(false); !rows[2][0].Number {
		t.Errorf("cells = %+v", rows)
	}
}

func TestParams(t *testing.T) {
	f := conv.Lookup("csv")
	if f == nil {
		t.Fatal("csv is not registered")
	}
	data := []byte("a;b,c\n1;2,5\n3;4,5\n")
	run := func(params map[string]string) (*conv.Result, error) {
		o := &conv.Options{Params: params, FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true, FileName: "data.csv"}
		return f.Convert(bytes.NewReader(data), int64(len(data)), o)
	}
	for _, c := range []struct {
		params  map[string]string
		summary string
	}{
		{nil, "3 row(s) × 2 column(s) (UTF-8, semicolon-separated, header row)"},
		{map[string]string{"delimiter": "comma"}, "3 row(s) × 2 column(s) (UTF-8, comma-separated, header row)"},
		{map[string]string{"delimiter": "tab", "header": "false"}, "3 row(s) × 1 column(s) (UTF-8, tab-separated, no header row)"},
		{map[string]string{"quote": "none", "header": "true", "charset": "windows-1252"}, "(WINDOWS-1252, semicolon-separated, no quotes, header row)"},
		{map[string]string{"table": "TableStyleLight9"}, "semicolon-separated"},
		{map[string]string{"delimiter": `"`}, `'"'-separated, no quotes`},
	} {
		res, err := run(c.params)
		if err != nil {
			t.Errorf("%v: %v", c.params, err)
			continue
		}
		if !strings.Contains(res.Summary, c.summary) {
			t.Errorf("%v: summary %q, want %q", c.params, res.Summary, c.summary)
		}
		if v := res.Doc.Views[0]; v.Title != "data" {
			t.Errorf("%v: sheet %q", c.params, v.Title)
		}
	}
	for _, params := range []map[string]string{{"delimiter": "ab"}, {"delimiter": "é"}, {"quote": "back"}, {"header": "maybe"},
		{"charset": "klingon"}, {"table": "nope"}, {"delimiter": "\n"}, {"delimiter": `"`, "quote": "double"}} {
		if _, err := run(params); err == nil {
			t.Errorf("%v was accepted", params)
		}
	}
}

func TestDetect(t *testing.T) {
	for _, c := range []struct {
		name string
		data string
		want bool
	}{
		{"csv", "a,b,c\n1,2,3\n", true},
		{"tsv", "a\tb\n1\t2\n3\t4\n", true},
		{"utf-16", "\xff\xfea\x00,\x00b\x00\n\x001\x00,\x002\x00\n\x00", true},
		{"shift_jis", "\x96\xbc\x91\x4f,\x94\x4e\x97\xee\n\x8eR\x93c,30\n", true},
		{"one line", "a,b,c\n", false},
		{"one column", "a\nb\nc\n", false},
		{"prose", "Hello, world. This is text,\nwith commas, here and there, and there.\nAnd a line without.\n", false},
		{"binary", "a,b\n\x00\x01\x02,\x03\n", false},
		{"pdf", "%PDF-1.7\n1,2\n3,4\n", false},
		{"empty", "", false},
	} {
		if got := detect([]byte(c.data), bytes.NewReader([]byte(c.data)), int64(len(c.data))); got != c.want {
			t.Errorf("%s: detected %v, want %v", c.name, got, c.want)
		}
	}
	// the first KiB alone does not hold two records: more is read
	long := strings.Repeat("x", 1500) + ",y\n1,2\n3,4\n"
	if !detect([]byte(long[:1024]), strings.NewReader(long), int64(len(long))) {
		t.Error("long first record: not detected")
	}
}
