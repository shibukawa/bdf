package parquet

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/xlsx"
)

// The files of testdata are written by test/parquet/gen.py; their values
// were checked cell by cell against what pyarrow reads.

func openTestFile(t *testing.T, name string) *os.File {
	t.Helper()
	f, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { f.Close() })
	return f
}

// readTable reads a test file into a grid, failing on warnings.
func readTable(t *testing.T, name string, opts *Options) (*xlsx.Grid, *Result) {
	t.Helper()
	f := openTestFile(t, name)
	st, _ := f.Stat()
	if opts == nil {
		opts = &Options{}
	}
	g, res, err := table(f, st.Size(), opts, func(msg string) { t.Errorf("%s: warning: %s", name, msg) })
	if err != nil {
		t.Fatal(err)
	}
	return g, res
}

// texts returns the texts of a grid's rows, cells separated by " | ", and
// numbers marked with ">".
func texts(g *xlsx.Grid) []string {
	var out []string
	for _, row := range g.Rows {
		var cells []string
		for _, c := range row {
			if c.Number {
				cells = append(cells, ">"+c.Text)
			} else {
				cells = append(cells, c.Text)
			}
		}
		out = append(out, strings.Join(cells, " | "))
	}
	return out
}

func checkRows(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%d rows, want %d", len(got), len(want))
	}
	for i := range min(len(got), len(want)) {
		if got[i] != want[i] {
			t.Errorf("row %d:\n got  %s\n want %s", i, got[i], want[i])
		}
	}
}

func TestBasic(t *testing.T) {
	g, res := readTable(t, "basic.parquet", nil)
	checkRows(t, texts(g), []string{
		"id | name | 説明 | price | qty | ratio | active | day | updated | tags | size",
		"int64 | string | string | decimal(10, 2) | int32 | double | boolean | date | timestamp[ms, UTC] | list<string> | struct<w: int32, h: int32>",
		`>1 | Ann | 本文の説明 | >1200.00 | >3 | >0.125 | true | >2026-01-05 | >2026-09-27 09:30:00.000 | ["red", "blue"] | {"w": 1920, "h": 1080}`,
		`>2 | Bob | 図表 | >35.50 | >12 | >3.14159 | false | >2026-02-14 | >2026-09-27 10:15:30.250 | [] | {"w": 1280, "h": 720}`,
		">3 | さとう | 英語と日本語 | >-8.25 | >0 | >-2.5 |  | >2026-03-01 | >2026-09-27 11:00:00.500 |  | ",
		`>4 | すずき | 一行目` + "\n" + `二行目 | >0.99 | >7 | >1e-07 | true | >2026-04-30 |  | ["green"] | {"w": 800, "h": 600}`,
		">5 |  | 空白 |  |  |  |  |  |  |  | ",
		`>6 | カタカナ | Mixed テキスト | >99999.99 | >42 | >2.5e+21 | false | >2025-12-31 | >2026-09-28 00:00:00.001 | ["a", "b", "c"] | {"w": 3, "h": 4}`,
		`>7 | Émile | Latin-1 à é ü | >10.00 | >-5 | >100 | true | >2024-02-29 | >2026-01-01 00:00:00.000 | ["x"] | {"w": 1, "h": 1}`,
		`>8 | Zoë |  | >0.00 | >1000000 | >0 | false | >2000-01-01 | >1999-12-31 23:59:59.999 | [null, "y"] | {"w": 0, "h": 0}`,
	})
	if g.HeaderRows != 2 || !g.SubHeader || g.Lang != "ja" {
		t.Errorf("grid: %d header rows, subheader %v, lang %q", g.HeaderRows, g.SubHeader, g.Lang)
	}
	if res.Rows != 8 || res.Shown != 8 || res.Cols != 11 || res.RowGroups != 1 || !slices.Equal(res.Codecs, []string{"Snappy"}) ||
		!strings.HasPrefix(res.CreatedBy, "parquet-cpp-arrow") {
		t.Errorf("result: %+v", res)
	}
}

func TestTypes(t *testing.T) {
	g, _ := readTable(t, "types.parquet", nil)
	checkRows(t, texts(g), []string{
		"i8 | i16 | i32 | i64 | u8 | u16 | u32 | u64 | f16 | f32 | f64 | f64b | dec5 | dec15 | dec30 | str | lstr | bin | fixed | uuid | json | date | time_ms | time_us | time_ns | ts_s | ts_ms | ts_us | ts_ns | dict | null",
		"int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float16 | float | double | double | decimal(5, 2) | decimal(15, 3) | decimal(30, 4) | string | string | binary | fixed[4] | uuid | json | date | time[ms] | time[us] | time[ns] | timestamp[ms] | timestamp[ms] | timestamp[us, UTC] | timestamp[ns, UTC] | string | null",
		">-128 | >-32768 | >-2147483648 | >-9223372036854775808 | >0 | >0 | >0 | >0 | >0.1 | >0.1 | >NaN | >1e-07 | >-1.23 | >-123456789012.345 | >-12345678901234567890123456.7890 | text | large | plain text | ABCD | 12345678-1234-5678-1234-567812345678 | {\"a\": 1} | >1970-01-01 | >00:00:00.000 | >00:00:00.000000 | >00:00:00.000000000 | >1970-01-01 00:00:00 | >1970-01-01 00:00:00.000 | >1970-01-01 00:00:00.000000 | >1970-01-01 00:00:00.000000000 | red | ",
		">0 | >1 | >2 | >3 | >1 | >1 | >1 | >1 | >-65500 | >-3.4028235e+38 | >Infinity | >123456789.125 | >0.05 | >0.001 | >0.0001 | かな | string | 0x000102ff | 0x00000001 | 00000000-0000-0000-0000-000000000000 | [1, 2] | >1969-12-31 | >12:34:56.789 | >12:34:56.789012 | >12:34:56.789012345 | >1969-12-31 23:59:59 | >1969-12-31 23:59:59.999 | >1969-12-31 23:59:59.999999 | >1969-12-31 23:59:59.999999999 | green | ",
		">127 | >32767 | >2147483647 | >9223372036854775807 | >255 | >65535 | >4294967295 | >18446744073709551615 | >6e-08 | >1.5 | >-Infinity | >1e+21 | >999.99 | >1.000 | >42.0000 |  | テキスト | 0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f… (40 bytes) | 0xdeadbeef | f24f9b64-81fa-49d1-b74e-8c09a6e31c56 | \"s\" | >9999-12-31 | >23:59:59.999 | >23:59:59.999999 | >23:59:59.999999999 | >2026-09-27 10:30:00 | >2026-09-27 10:30:00.123 | >2026-09-27 10:30:00.123456 | >2026-09-27 10:30:00.123456789 | red | ",
		strings.Repeat(" | ", 30),
	})
}

func TestNested(t *testing.T) {
	g, _ := readTable(t, "nested.parquet", nil)
	checkRows(t, texts(g), []string{
		"list | list_list | list_struct | struct_list | map | map_struct | fixed_list | large_list | deep",
		"list<int32> | list<list<string>> | list<struct<k: string, v: double>> | struct<name: string, scores: list<int16>> | map<string, int64> | map<int32, struct<ok: boolean>> | list<float> | list<int64> | list<list<list<int8>>>",
		`[1, 2] | [["a"], ["b", "c"]] | [{"k": "x", "v": 1.5}] | {"name": "p", "scores": [1, 2]} | {"a": 1, "b": null} | {1: {"ok": true}} | [1, 2, 3] | [10] | [[[1, 2], [3]], [[4]]]`,
		`[] | [[], null] | [{"k": null, "v": null}, null] | {"name": null, "scores": []} | {} | {2: null} |  | [20, 30] | [[[]]]`,
		` |  | [] |  |  |  | [0.5, -0.5, 0] | [] | [[null]]`,
		`[null, 3] | [[null]] |  | {"name": "q", "scores": null} | {"c": 3} | {} | [1e-07, 1e+21, 2] |  | `,
	})
}

// TestEncodings reads data pages v2 of every encoding, each column beside
// the same values in PLAIN.
func TestEncodings(t *testing.T) {
	g, res := readTable(t, "encodings.parquet", nil)
	if res.Rows != 300 || len(g.Rows) != 302 {
		t.Fatalf("%d rows, %d shown", res.Rows, len(g.Rows)-2)
	}
	names := g.Rows[0]
	for j, h := range names {
		if strings.HasPrefix(h.Text, "plain_") {
			continue
		}
		p := slices.IndexFunc(names, func(c xlsx.GridCell) bool { return c.Text == "plain_"+h.Text })
		if p < 0 {
			t.Fatalf("no plain column for %s", h.Text)
		}
		empty := 0
		for i, row := range g.Rows[2:] {
			if row[j] != row[p] {
				t.Errorf("%s row %d: %q, plain %q", h.Text, i, row[j].Text, row[p].Text)
				break
			}
			if row[j].Text == "" {
				empty++
			}
		}
		if empty > 30 {
			t.Errorf("%s: %d nulls", h.Text, empty)
		}
	}
	if got := g.Rows[2+299]; got[0].Text != "-4718" || got[2].Text != "1839506156249" || got[6].Text != "prefix/0099/c" || got[8].Text != "64.75" {
		t.Errorf("last row: %v", got)
	}
}

// TestCodecs reads a column compressed with each codec.
func TestCodecs(t *testing.T) {
	g, res := readTable(t, "codecs.parquet", nil)
	for i, row := range g.Rows[2:] {
		for j := range row {
			if row[j] != row[0] {
				t.Fatalf("row %d: %s %q, uncompressed %q", i, g.Rows[0][j].Text, row[j].Text, row[0].Text)
			}
		}
	}
	if g.Rows[3][0].Text != "value 1 repeated text " {
		t.Errorf("row 1: %q", g.Rows[3][0].Text)
	}
	want := []string{"uncompressed", "Snappy", "gzip", "Brotli", "Zstandard", "LZ4"}
	if !slices.Equal(res.Codecs, want) {
		t.Errorf("codecs %v, want %v", res.Codecs, want)
	}
}

func TestLegacy(t *testing.T) {
	g, _ := readTable(t, "legacy.parquet", nil)
	checkRows(t, texts(g), []string{
		"ts | list",
		"timestamp[ns] | list<int32>",
		">2024-01-01 20:34:56.123456 | [1, 2]",
		">1900-01-01 00:00:00.000000 | []",
		" | ",
	})
}

// TestDuckDB reads DuckDB's output: shredded Variants, UUIDs, intervals and
// native geometries, and the PLAIN_DICTIONARY encoding.
func TestDuckDB(t *testing.T) {
	g, _ := readTable(t, "duckdb.parquet", nil)
	checkRows(t, texts(g), []string{
		"id | v | u | iv | g | color",
		"int32 | variant | uuid | interval | geometry | string",
		">1 | 42 | b5f3c2de-0a41-4d4b-9f3e-1c2d3e4f5a6b | >P3D | POINT (1 2) | red",
		`>2 | "text" |  | >P1Y2M3DT4H5M6.5S | LINESTRING (0 0, 1 1) | blue`,
		`>3 | {"a": 1, "b": [true, null]} | 00000000-0000-0000-0000-000000000000 | >PT0S |  | red`,
		">4 | null |  |  | POLYGON ((0 0, 4 0, 4 4, 0 0)) | ",
		`>5 | [1, "two", 3.5] |  | >PT1H30M | POINT EMPTY | blue`,
	})
}

// TestVariants reads Variants shredded into typed objects and lists.
func TestVariants(t *testing.T) {
	g, _ := readTable(t, "variant.parquet", nil)
	checkRows(t, texts(g), []string{
		"id | obj | arr",
		"int32 | variant | variant",
		`>1 | {"a": 1, "b": "x"} | [1, 2]`,
		`>2 | {"a": 2, "b": "y", "c": true} | [3]`,
		`>3 | {"a": null, "b": "z"} | []`,
		`>4 | 5 | "not a list"`,
		`>5 | {"a": "text", "b": {"deep": [1.5]}} | [4, null]`,
		">6 | null | null",
	})
}

// TestGeo reads the WKB column that GeoParquet metadata names.
func TestGeo(t *testing.T) {
	g, _ := readTable(t, "geo.parquet", nil)
	checkRows(t, texts(g), []string{
		"name | geometry",
		"string | geometry",
		"point | POINT (139.767 35.681)",
		"line | LINESTRING (0 0, 1.5 2, 3 0)",
		"polygon | POLYGON ((0 0, 10 0, 10 10, 0 10, 0 0), (2 2, 3 2, 3 3, 2 2))",
		"multipoint | MULTIPOINT ((1 2), (3 4))",
		"multipolygon | MULTIPOLYGON (((0 0, 1 0, 0 1, 0 0)))",
		"collection | GEOMETRYCOLLECTION (POINT (5 6), LINESTRING (7 8, 9 10))",
		"point z | POINT Z (1 2 3)",
		"empty | POINT EMPTY",
		"null | ",
	})
}

func TestRows(t *testing.T) {
	f := openTestFile(t, "encodings.parquet")
	st, _ := f.Stat()
	var warnings []string
	g, res, err := table(f, st.Size(), &Options{Rows: 5, NoTypes: true}, func(msg string) { warnings = append(warnings, msg) })
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Rows) != 6 || g.HeaderRows != 1 || g.SubHeader || res.Shown != 5 || res.Rows != 300 {
		t.Errorf("%d rows, %d header rows, result %+v", len(g.Rows), g.HeaderRows, res)
	}
	if want := []string{"only the first 5 of 300 rows are shown (the rows parameter shows more)"}; !slices.Equal(warnings, want) {
		t.Errorf("warnings %q", warnings)
	}
	warnings = nil
	if g, _, _ = table(f, st.Size(), &Options{Rows: -1}, func(msg string) { warnings = append(warnings, msg) }); len(g.Rows) != 302 || warnings != nil {
		t.Errorf("all rows: %d rows, warnings %q", len(g.Rows), warnings)
	}
}

func TestConvert(t *testing.T) {
	res, err := ConvertFile(filepath.Join("testdata", "basic.parquet"), &Options{FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) != 0 || res.EmbeddedFonts != 2 {
		t.Errorf("warnings %q, %d embedded font(s)", res.Warnings, res.EmbeddedFonts)
	}
	var buf bytes.Buffer
	if err := res.Doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	m := r.Manifest
	if m.Meta.Source != "parquet" || len(m.Views) != 1 || m.Views[0].Title != "basic" || m.Views[0].Freeze == nil || m.Views[0].Freeze.Rows != 2 {
		t.Fatalf("meta %+v, views %+v", m.Meta, m.Views)
	}
	h, err := bdf.ParseHash(m.Views[0].TextIndex)
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
	text := bdf.PlainText(idx)
	for _, want := range []string{"説明", "decimal(10, 2)", "本文の説明", "2026-09-27 10:15:30.250", `{"w": 1920, "h": 1080}`} {
		if !strings.Contains(text, want) {
			t.Errorf("text index lacks %q:\n%s", want, text)
		}
	}
}

func TestParams(t *testing.T) {
	f := conv.Lookup("parquet")
	if f == nil {
		t.Fatal("parquet is not registered")
	}
	data, err := os.ReadFile(filepath.Join("testdata", "encodings.parquet"))
	if err != nil {
		t.Fatal(err)
	}
	if !f.Detect(data[:1024], bytes.NewReader(data), int64(len(data))) {
		t.Error("not detected")
	}
	run := func(params map[string]string) (*conv.Result, error) {
		o := &conv.Options{Params: params, FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true, FileName: "data.parquet"}
		return f.Convert(bytes.NewReader(data), int64(len(data)), o)
	}
	for _, c := range []struct {
		params  map[string]string
		summary string
	}{
		{nil, "300 row(s) × 16 column(s) (1 row group(s); gzip; written by parquet-cpp-arrow version 25.0.1), 2 embedded font(s)"},
		{map[string]string{"rows": "10", "types": "false"}, "10 of 300 row(s) × 16 column(s)"},
		{map[string]string{"rows": "all", "table": "TableStyleMedium2"}, "300 row(s) × 16 column(s)"},
	} {
		res, err := run(c.params)
		if err != nil {
			t.Fatalf("%v: %v", c.params, err)
		}
		if !strings.Contains(res.Summary, c.summary) {
			t.Errorf("%v: summary %q", c.params, res.Summary)
		}
		if res.Doc.Views[0].Title != "data" {
			t.Errorf("%v: sheet %q", c.params, res.Doc.Views[0].Title)
		}
	}
	for _, params := range []map[string]string{{"rows": "0"}, {"rows": "many"}, {"types": "maybe"}, {"table": "NoSuchStyle"}} {
		if _, err := run(params); err == nil {
			t.Errorf("%v: no error", params)
		}
	}
}

func TestNotParquet(t *testing.T) {
	for name, b := range map[string][]byte{
		"empty":     nil,
		"text":      []byte("id,name\n1,Ann\n"),
		"no footer": []byte("PAR1\x00\x00\x00\x00PAR1"),
		"long":      []byte("PAR1\xff\xff\x00\x00PAR1"),
		"garbage":   append([]byte("PAR1"), append(bytes.Repeat([]byte{0xff}, 20), "\x10\x00\x00\x00PAR1"...)...),
	} {
		if _, _, err := table(bytes.NewReader(b), int64(len(b)), &Options{}, func(string) {}); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
	b := []byte("PARE\x00\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00PARE")
	if _, _, err := table(bytes.NewReader(b), int64(len(b)), &Options{}, func(string) {}); err != ErrEncrypted {
		t.Errorf("encrypted footer: %v", err)
	}
}

// TestDamaged reads the test files with each byte changed and cut short at
// each length: errors and warnings are expected, panics are not.
func TestDamaged(t *testing.T) {
	for _, name := range []string{"basic.parquet", "nested.parquet", "duckdb.parquet"} {
		data, err := os.ReadFile(filepath.Join("testdata", name))
		if err != nil {
			t.Fatal(err)
		}
		b := make([]byte, len(data))
		for i := range data {
			for _, x := range []byte{0xff, 0x01} {
				copy(b, data)
				b[i] ^= x
				table(bytes.NewReader(b), int64(len(b)), &Options{}, func(string) {})
			}
			cut := data[:i]
			table(bytes.NewReader(cut), int64(len(cut)), &Options{}, func(string) {})
		}
	}
}

func FuzzTable(f *testing.F) {
	for _, name := range []string{"basic.parquet", "nested.parquet", "duckdb.parquet", "legacy.parquet"} {
		data, err := os.ReadFile(filepath.Join("testdata", name))
		if err != nil {
			f.Fatal(err)
		}
		f.Add(data)
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		start := time.Now()
		table(bytes.NewReader(b), int64(len(b)), &Options{}, func(string) {})
		if d := time.Since(start); d > time.Second {
			t.Errorf("%v to read %d bytes", d, len(b))
		}
	})
}
