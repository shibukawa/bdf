package textform

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/xlsx"
	"github.com/shibukawa/bdf/form/tag"
	"github.com/shibukawa/bdf/form/textpart"
	"golang.org/x/text/encoding/japanese"
)

var update = flag.Bool("update", false, "write the plain text of the test forms (testdata/want)")

func testOptions() *Options {
	return &Options{FontDirs: []string{"testdata/fonts"}, NoSystemFonts: true}
}

func fillString(tpl string, data map[string]any, strict bool) (*Result, error) {
	if data != nil {
		tag.Normalize(data)
	}
	return Fill(strings.NewReader(tpl), int64(len(tpl)), &TemplateData{Data: data, Strict: strict, Format: xlsx.Formatter()}, testOptions())
}

// plain fills a template and returns the text of its print image, without
// the margins.
func plain(t *testing.T, tpl string, data map[string]any) (string, *Result) {
	t.Helper()
	res, err := fillString(tpl, data, false)
	if err != nil {
		t.Fatal(err)
	}
	return res.Image.Text(textpart.TextOptions{}), res
}

func TestFrontMatter(t *testing.T) {
	var warnings []string
	warn := func(format string, args ...any) { warnings = append(warnings, fmt.Sprintf(format, args...)) }
	g, err := parseFrontMatter(nil, warn)
	if err != nil || g != defaultGrid() {
		t.Errorf("no front matter: %+v, %v", g, err)
	}
	g, err = parseFrontMatter(strings.Split(`# a slip
cpi: 12            # columns to the inch
lpi: "8"
columns: 136
lines: 66
margin: {left: 2, top: 1}
skip: 3
title: 4
font: 'BIZ UDゴシック'
ambiguous: wide
colour: red`, "\n"), warn)
	want := grid{cpi: 12, lpi: 8, columns: 136, lines: 66, left: 2, top: 1, skip: 3, title: 4, font: "BIZ UDゴシック", ambiguousWide: true}
	if err != nil || g != want {
		t.Errorf("front matter: %+v, %v", g, err)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "colour") {
		t.Errorf("warnings: %q", warnings)
	}
	g, err = parseFrontMatter([]string{"margin:", "  left: 5", "  top: 2 # lines", "lines: 0", "cpi: 13"}, warn)
	if err != nil || g.left != 5 || g.top != 2 || g.lines != 0 || g.cpi != 13 {
		t.Errorf("a block of the margin: %+v, %v", g, err)
	}
	if !strings.Contains(warnings[len(warnings)-1], "13 CPI") {
		t.Errorf("a pitch of no printer: %q", warnings)
	}
	for _, bad := range []string{"cpi: 0", "cpi: ten", "lpi: 31", "columns: 0", "lines: -1", "margin: {left: 80}", "margin: {top: 60}\nskip: 6",
		"title: 66", "ambiguous: half", "just text", "  left: 1", "margin: {left: 1"} {
		if _, err := parseFrontMatter(strings.Split(bad, "\n"), warn); err == nil {
			t.Errorf("parseFrontMatter(%q) takes it", bad)
		}
	}
	front, body, first, has := splitFrontMatter("---\ncpi: 12\n---\nbody\n")
	if !has || len(front) != 1 || body != "body\n" || first != 4 {
		t.Errorf("splitFrontMatter: %q %q %d %v", front, body, first, has)
	}
	if _, body, first, has := splitFrontMatter("---\nno end\n"); has || body != "---\nno end\n" || first != 1 {
		t.Error("a rule without the end of front matter is body")
	}
}

func TestIsTemplate(t *testing.T) {
	for s, want := range map[string]bool{
		"---\ncpi: 10\n---\nx": true,
		"Dear {{name}}":        true,
		"\xef\xbb\xbf｛｛名前｝｝":   true,
		"plain text":           false,
		"a\x00{{b}}":           false,
		"\xff\xfe{{b}}":        false,
		"{{a}} " + "\xe3\x81":  true, // (a character cut at the end)
	} {
		if got := IsTemplate([]byte(s)); got != want {
			t.Errorf("IsTemplate(%q) = %v", s, got)
		}
	}
}

func TestValues(t *testing.T) {
	tpl := `{{a}}|{{b | width:6}}|{{c | right:6}}|{{d | center:7}}|{{n | fmt:"#,##0.0" | right:9}}|
{{j | width:6}}|{{j | right:7}}|{{j | center:8}}|`
	got, res := plain(t, tpl, map[string]any{"a": "x", "b": "ab", "c": "cd", "d": "e", "n": 1234.5, "j": "漢字"})
	want := "x|ab    |    cd|   e   |  1,234.5|\n漢字  |   漢字|  漢字  |\n"
	if got != want {
		t.Errorf("values:\n%s\nwant\n%s", got, want)
	}
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	// a value wider than its columns: kept (with a warning), cut, wrapped
	tpl = "[{{v | width:5}}]\n[{{v | width:5 | cut}}]\n[{{v | right:5 | cut}}]\n<{{w | width:5 | wrap}}> <{{x | right:3 | wrap}}> end\n[{{z | width:3 | wrap}}]"
	got, res = plain(t, tpl, map[string]any{"v": "あいうえ", "w": "東京都千代田区丸の内", "x": "12345", "z": ""})
	want = "[あいうえ]\n[あい ]\n[ あい]\n<東京 > <123> end\n 都千     45\n 代田\n 区丸\n の内\n[   ]\n"
	if got != want {
		t.Errorf("overflow:\n%s\nwant\n%s", got, want)
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "{{v}} is 8 columns, more than its 5") {
		t.Errorf("warnings: %q", res.Warnings)
	}
	// line breaks in a value without columns go on under it; tabs are
	// spaces to the stops of 8 columns
	got, _ = plain(t, "Addr: {{a}} .\nx\ty", map[string]any{"a": "1-2\nTokyo"})
	if want := "Addr: 1-2 .\n      Tokyo\nx       y\n"; got != want {
		t.Errorf("line breaks:\n%q\nwant\n%q", got, want)
	}
	// a line wider than the paper is told
	_, res = plain(t, "---\ncolumns: 10\nmargin: {left: 2}\n---\n{{a}}", map[string]any{"a": "123456789"})
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "9 columns; the paper has 8") {
		t.Errorf("warnings: %q", res.Warnings)
	}
}

func TestMissing(t *testing.T) {
	got, res := plain(t, "a {{x}} b", map[string]any{})
	if got != "a  b\n" || len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], `line 1: the data has no "x"`) {
		t.Errorf("a missing name: %q %q", got, res.Warnings)
	}
	if _, err := fillString("---\nlpi: 6\n---\n\na {{x}} b", map[string]any{}, true); err == nil || !strings.Contains(err.Error(), `line 5: the data has no "x"`) {
		t.Errorf("strict: %v", err)
	}
}

func TestSections(t *testing.T) {
	data := map[string]any{"tags": []any{"a", "b", "c"}, "none": []any{}, "ok": true,
		"items": []any{map[string]any{"n": "one", "q": 1}, map[string]any{"n": "two", "q": 2}}}
	// characters, inverted sections, lines of tags only
	tpl := "{{#tags}}{{.}}{{^_last}},{{/_last}}{{/tags}} {{^none}}none{{/none}}\n{{#ok}}\nshown\n{{/ok}}\n{{#none | rows}}\nhidden\n{{/none}}\nend"
	got, _ := plain(t, tpl, data)
	if want := "a,b,c none\nshown\nend\n"; got != want {
		t.Errorf("sections:\n%q\nwant\n%q", got, want)
	}
	// slots: a line for an item, the lines without tags are the first
	tpl = "{{#items}}\n| {{_index}} {{n | width:5}}|{{q | right:3}}|\n|        |   |\n| blank{{#_blank_first}}!{{/_blank_first}} |   |\n{{/items}}\nend"
	got, res := plain(t, tpl, data)
	if want := "| 1 one  |  1|\n| 2 two  |  2|\n| blank! |   |\nend\n"; got != want {
		t.Errorf("slots:\n%s\nwant\n%s", got, want)
	}
	if len(res.Fields) != 3 || res.Fields[0].Mode != tag.ModeSlots || res.Fields[0].Slots != 3 || res.Fields[1].Name != "n" || !reflect.DeepEqual(res.Fields[1].In, []string{"items"}) {
		t.Errorf("fields: %+v", res.Fields)
	}
	if _, err := fillString(tpl, map[string]any{"items": []any{1, 2, 3, 4}}, false); err == nil || !strings.Contains(err.Error(), "4 items for the 3 lines of {{#items}}") {
		t.Errorf("more items than slots: %v", err)
	}
	// rows, and rows in rows
	tpl = "{{#groups | rows}}\n# {{g}}\n{{#items | rows}}\n  - {{.}}\n{{/items}}\n{{/groups}}"
	got, _ = plain(t, tpl, map[string]any{"groups": []any{map[string]any{"g": "A", "items": []any{"x", "y"}}, map[string]any{"g": "B", "items": []any{"z"}}}})
	if want := "# A\n  - x\n  - y\n# B\n  - z\n"; got != want {
		t.Errorf("rows:\n%s\nwant\n%s", got, want)
	}
	// a line of the same line repeated with | rows
	got, _ = plain(t, "{{#tags | rows}}- {{.}}{{/tags}}", data)
	if got != "- a\n- b\n- c\n" {
		t.Errorf("rows of a line: %q", got)
	}
	for tpl, msg := range map[string]string{
		"{{#a}}\n{{#b | rows}}\n{{/b}}\n{{/a}}":    "is in the lines of {{#a}}, which are slots",
		"{{#a}}x":                                  "line 1: the section {{#a}} is not closed",
		"x\n{{/a}}":                                "line 2: {{/a}} closes no section",
		"{{#a}}{{#b}}\n{{/a}}{{/b}}":               "{{/a}} closes {{#b}}",
		"{{^a | slots}}\n{{/a}}":                   "does not repeat",
		"{{#a | rows}}x{{/a}}{{#b | rows}}y{{/b}}": "both repeat line 1",
		"{{%photo}}":                               "no pictures nor codes",
		"{{#_wide | rows}}x{{/_wide}}":             "marks how characters are printed",
		"{{a | width:0}}":                          "takes columns from 1",
		"{{a | cut}}":                              "takes cut with its columns only",
		"ok\n  {{a | fmt}}":                        "line 2, column 3",
		"{{#a}}\n{{/a}}":                           "has no line to fill",
	} {
		if _, err := fillString(tpl, map[string]any{}, false); err == nil || !strings.Contains(err.Error(), msg) {
			t.Errorf("%q: %v, want %q", tpl, err, msg)
		}
	}
}

func TestPages(t *testing.T) {
	// a list over pages of 5 lines (8 less the margin and the skip), its
	// two lines of title on each
	tpl := "---\nlines: 8\nmargin: {top: 1}\nskip: 2\ntitle: 2\n---\nTitle\n-----\n{{#items | rows}}\n{{.}}\n{{/items}}\nend"
	items := []any{"1", "2", "3", "4", "5", "6", "7"}
	got, res := plain(t, tpl, map[string]any{"items": items})
	want := "Title\n-----\n1\n2\n3\n\fTitle\n-----\n4\n5\n6\n\fTitle\n-----\n7\nend\n"
	if got != want {
		t.Errorf("pages:\n%q\nwant\n%q", got, want)
	}
	if n := len(res.Doc.Views[0].Pages); n != 3 {
		t.Errorf("%d pages", n)
	}
	// a form feed
	got, _ = plain(t, "a\n\fb\nc\fd\n\f", nil)
	if want := "a\n\fb\nc\n\fd\n"; got != want {
		t.Errorf("form feeds: %q, want %q", got, want)
	}
	// a page for each item
	tpl = "---\nlines: 6\n---\n{{#pages | pages}}\n{{_page}}/{{_pages}} {{name}}\n{{/pages}}"
	got, _ = plain(t, tpl, map[string]any{"pages": []any{map[string]any{"name": "a"}, map[string]any{"name": "b"}}})
	if want := "1/2 a\n\f2/2 b\n"; got != want {
		t.Errorf("a page for each: %q, want %q", got, want)
	}
	// a form that does not grow fits its pages; empty lines past them are
	// none of its lines
	if _, err := fillString("---\nlines: 3\n---\na\nb\n{{v | width:2 | wrap}}\nd", map[string]any{"v": "1234"}, false); err == nil || !strings.Contains(err.Error(), "line 6 of the template is past the 3 lines of a page") {
		t.Errorf("a fixed form longer than its page: %v", err)
	}
	got, _ = plain(t, "---\nlines: 3\n---\na\nb\nc\n\n\n", nil)
	if got != "a\nb\nc\n" {
		t.Errorf("empty lines past the page: %q", got)
	}
	// a roll: a page as long as its lines
	res, err := fillString("---\nlines: 0\nmargin: {top: 2}\nlpi: 8\n---\n{{#items | rows}}\n{{.}}\n{{/items}}", map[string]any{"items": items}, false)
	if err != nil {
		t.Fatal(err)
	}
	if p := res.Doc.Views[0].Pages; len(p) != 1 || p[0].H != 81 || p[0].W != 576 {
		t.Errorf("the page of a roll: %+v", p[0])
	}
	// the pages to keep
	o := testOptions()
	o.Pages, _ = conv.ParsePages("2")
	s := "a\fb\fc"
	res, err = Fill(strings.NewReader(s), int64(len(s)), nil, o)
	if err != nil || res.Image.Text(textpart.TextOptions{}) != "b\n" {
		t.Errorf("-pages: %v", err)
	}
}

func TestAttributes(t *testing.T) {
	tpl := "{{#_big}}納品書{{/_big}} {{#_wide}}{{n | right:6}}{{/_wide}}\n{{#preprinted}}\n┌──┐\n{{/preprinted}}\n{{#preprinted}}No.{{/preprinted}} {{no}}\n{{#_tall}}\nTotal {{t}}\n{{/_tall}}\nend"
	got, res := plain(t, tpl, map[string]any{"n": 12, "no": 7, "t": 100})
	want := "納  品  書     1 2\n\n┌──┐\nNo. 7\nTotal 100\n\nend\n"
	if got != want {
		t.Errorf("attributes:\n%q\nwant\n%q", got, want)
	}
	lines := res.Image.Pages[0].Lines
	wantLines := []textpart.Line{
		{Text: "納品書   12", Size: []textpart.Range{{Col: 0, Cols: 12, Size: "big"}, {Col: 15, Cols: 4, Size: "wide"}}},
		{},
		{Text: "┌──┐", Paper: []textpart.Range{{Col: 0, Cols: 4}}},
		{Text: "No. 7", Paper: []textpart.Range{{Col: 0, Cols: 3}}},
		{Text: "Total 100", Size: []textpart.Range{{Col: 0, Cols: 9, Size: "tall"}}},
		{},
		{Text: "end"},
	}
	if !reflect.DeepEqual(lines, wantLines) {
		t.Errorf("the lines of the print image:\n%+v\nwant\n%+v", lines, wantLines)
	}
	if got := res.Image.Text(textpart.TextOptions{NoPaper: true}); !strings.Contains(got, "\n\n    7\n") {
		t.Errorf("without the paper: %q", got)
	}
	if len(res.Fields) != 3 {
		t.Errorf("the sections of attributes are not fields: %+v", res.Fields)
	}
	// the paper is drawn on a layer of the role background
	ls := res.Doc.Views[0].Pages[0].Layers
	if len(ls) != 2 || ls[0].Role != bdf.RoleBackground || ls[1].Role != bdf.RoleBody {
		t.Errorf("layers: %+v", ls)
	}
}

// runs returns the runs a page draws: their text, x, base line, advance
// and font size, and the marks between them.
func runs(t *testing.T, doc *bdf.Document, layer int) []string {
	t.Helper()
	p := doc.Part(doc.Views[0].Pages[0].Layers[layer].Obj)
	o, err := bdf.DecodeObject(p.Data)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	size := float32(0)
	err = o.Walk(func(in bdf.Instr) {
		switch in.Op {
		case bdf.OpFont:
			size = in.Args[1].(float32)
		case bdf.OpFillText:
			out = append(out, fmt.Sprintf("%s@%g,%g+%g/%g", in.Args[0], in.Args[1], in.Args[2], in.Args[3], size))
		case bdf.OpMark:
			out = append(out, fmt.Sprintf("M%d%s", in.Args[0], in.Args[1]))
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestLayout(t *testing.T) {
	// 10 CPI, 6 LPI: a column is 7.2 points, a line 12
	tpl := "---\ncolumns: 40\nlines: 10\nmargin: {left: 1, top: 1}\n---\nAB 漢字 x\n{{#items}}\n{{.}}│\n\n{{/items}}\n{{#_big}}大{{/_big}}"
	res, err := fillString(tpl, map[string]any{"items": []any{"1", "2"}}, false)
	if err != nil {
		t.Fatal(err)
	}
	pg := res.Doc.Views[0].Pages[0]
	if math.Abs(float64(pg.W)-288) > 1e-3 || math.Abs(float64(pg.H)-120) > 1e-3 {
		t.Errorf("the page is %g × %g", pg.W, pg.H)
	}
	got := strings.Join(runs(t, res.Doc, 0), " ")
	want := "M0 AB @7.2,22.56+21.6/12 漢字 @28.8,22.56+36/12 x@64.8,22.56+7.2/12 " +
		"M9 M2A1 1@7.2,34.56+7.2/12 │@14.4,34.56+7.2/12 M2A2 2@7.2,46.56+7.2/12 │@14.4,46.56+7.2/12 M11 " +
		"M0 大@7.2,70.56+28.8/24"
	if got != want {
		t.Errorf("runs:\n%s\nwant\n%s", got, want)
	}
	if res.Image.Font != "BIZ UDGothic" {
		t.Errorf("the font is %q", res.Image.Font)
	}
	// the print image is a part of the document
	im, err := textpart.FromDocument(res.Doc)
	if err != nil || !reflect.DeepEqual(im, res.Image) {
		t.Errorf("the part: %v", err)
	}
}

func TestFontFile(t *testing.T) {
	font, err := os.ReadFile("testdata/fonts/BIZUDGothic-Regular-subset.ttf")
	if err != nil {
		t.Fatal(err)
	}
	tpl := "---\nfont: fonts/slip.ttf\n---\n伝票"
	files := fstest.MapFS{"fonts/slip.ttf": {Data: font}}
	res, err := Fill(strings.NewReader(tpl), int64(len(tpl)), &TemplateData{Files: files}, &Options{NoSystemFonts: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.Image.Font != "BIZ UDGothic" || len(res.Warnings) != 0 {
		t.Errorf("the font of the file: %q, %q", res.Image.Font, res.Warnings)
	}
	for name, msg := range map[string]string{"fonts/none.ttf": "does not exist", "../x.ttf": "not a path in the directory"} {
		tpl := "---\nfont: " + name + "\n---\nx"
		if _, err := Fill(strings.NewReader(tpl), int64(len(tpl)), &TemplateData{Files: files}, &Options{NoSystemFonts: true}); err == nil || !strings.Contains(err.Error(), msg) {
			t.Errorf("font: %s: %v", name, err)
		}
	}
}

func TestShow(t *testing.T) {
	tpl := "---\nlines: 4\n---\n{{#items}}\n{{name | width:4}}|\n{{/items}}\n\f{{#items}}x{{/items}}\na\nb\nc\nd"
	res, err := Show(strings.NewReader(tpl), int64(len(tpl)), nil, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	got := res.Image.Text(textpart.TextOptions{})
	if want := "{{#items}}\n{{name | width:4}}|\n{{/items}}\n\f{{#items}}x{{/items}}\na\nb\nc\n\fd\n"; got != want {
		t.Errorf("Show:\n%q\nwant\n%q", got, want)
	}
	if len(res.Fields) != 2 || !reflect.DeepEqual(res.Fields[0].At, []string{"line 4", "line 7"}) {
		t.Errorf("fields: %+v", res.Fields)
	}
}

func TestPlain(t *testing.T) {
	_, res := plain(t, "---\nmargin: {left: 1, top: 1}\n---\n{{#preprinted}}名前{{/preprinted}} \U00020bb7田 ①\n\fp2", nil)
	b, warnings, err := Plain(res.Image, map[string]string{"encoding": "shift_jis", "newline": "crlf", "paper": "preprinted"})
	if err != nil {
		t.Fatal(err)
	}
	// (the first character is not in Shift_JIS; ① is, as Windows-31J has it)
	want, _ := japanese.ShiftJIS.NewEncoder().String("\r\n      〓田 ①\r\n\f\r\n p2\r\n")
	if string(b) != want || len(warnings) != 1 || !strings.Contains(warnings[0], "U+20BB7") {
		t.Errorf("Plain = % x, %q", b, warnings)
	}
	b, _, err = Plain(res.Image, nil)
	if err != nil || string(b) != "\n 名前 \U00020bb7田 ①\n\f\n p2\n" {
		t.Errorf("Plain in UTF-8: %q, %v", b, err)
	}
	for _, p := range []map[string]string{{"encoding": "klingon"}, {"newline": "cr"}, {"paper": "none"}} {
		if _, _, err := Plain(res.Image, p); err == nil {
			t.Errorf("Plain takes %v", p)
		}
	}
}

// TestForms fills the templates of testdata with the data of test/text/data
// and compares their print images, as plain text, with testdata/want.
func TestForms(t *testing.T) {
	for _, name := range []string{"slip", "rows", "pages", "receipt"} {
		t.Run(name, func(t *testing.T) {
			tpl, err := os.ReadFile("testdata/" + name + ".txt")
			if err != nil {
				t.Fatal(err)
			}
			raw, err := os.ReadFile("../../test/text/data/" + name + ".json")
			if err != nil {
				t.Fatal(err)
			}
			var data map[string]any
			if err := json.Unmarshal(raw, &data); err != nil {
				t.Fatal(err)
			}
			res, err := Fill(bytes.NewReader(tpl), int64(len(tpl)), &TemplateData{Data: data, Strict: true, Format: xlsx.Formatter()}, testOptions())
			if err != nil {
				t.Fatal(err)
			}
			if len(res.Warnings) != 0 {
				t.Errorf("warnings: %q", res.Warnings)
			}
			got := res.Image.Text(textpart.TextOptions{Margins: true})
			file := "testdata/want/" + name + ".txt"
			if *update {
				if err := os.WriteFile(file, []byte(got), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			want, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			if got != string(want) {
				t.Errorf("the print image differs from %s (go test -update writes it):\n%s", file, got)
			}
			// the pages and the part agree
			if len(res.Image.Pages) != len(res.Doc.Views[0].Pages) {
				t.Errorf("%d pages in the part, %d in the view", len(res.Image.Pages), len(res.Doc.Views[0].Pages))
			}
			st, err := res.Doc.SearchText()
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(st.Views[0].Pages[0].Text, strings.Fields(res.Image.Pages[0].Lines[0].Text)[0]) {
				t.Errorf("the text of the first page: %q", st.Views[0].Pages[0].Text)
			}
		})
	}
}

func TestInspect(t *testing.T) {
	tpl, err := os.ReadFile("testdata/pages.txt")
	if err != nil {
		t.Fatal(err)
	}
	fields, err := Inspect(bytes.NewReader(tpl), int64(len(tpl)))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range fields {
		got = append(got, fmt.Sprintf("%s %s %s %s %d %v", f.Kind, strings.Join(f.In, "/"), f.Name, f.Mode, f.Slots, f.At))
	}
	want := []string{
		"section  pages pages 0 [line 9]",
		"value pages customer  0 [line 11]",
		"value pages carry_in  0 [line 12]",
		"section pages rows slots 10 [line 13]",
		"value pages/rows date  0 [line 14]",
		"value pages/rows name  0 [line 14]",
		"value pages/rows amount  0 [line 14]",
		// (in a section of a variable, as the other templates have it)
		"value pages/_last carry_out  0 [line 25]",
		"value pages/_last total  0 [line 25]",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("fields:\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestLimits(t *testing.T) {
	// sections in sections that print nothing end
	list := make([]any, 300)
	for i := range list {
		list[i] = i
	}
	tpl := strings.Repeat("{{#a}}", 4) + strings.Repeat("{{/a}}", 4)
	if _, err := fillString(tpl, map[string]any{"a": list}, false); err == nil || !strings.Contains(err.Error(), "repetitions of sections") {
		t.Errorf("sections of characters: %v", err)
	}
	tpl = strings.Repeat("{{#a | rows}}\n", 4) + strings.Repeat("{{/a}}\n", 4)
	if _, err := fillString(tpl, map[string]any{"a": list}, false); err == nil || !strings.Contains(err.Error(), "repetitions of sections") {
		t.Errorf("sections of lines: %v", err)
	}
	// a line that grows past its bound
	tpl = "{{#a}}{{#a}}0123456789{{/a}}{{/a}}"
	if _, err := fillString(tpl, map[string]any{"a": list}, false); err == nil || !strings.Contains(err.Error(), "more than 65536 columns") {
		t.Errorf("a long line: %v", err)
	}
}
