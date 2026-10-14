package pdfme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/formbox"
	"github.com/shibukawa/bdf/form/tag"
)

// The tests lay text out with the fonts of the tests alone (M PLUS 1p):
// the places they expect are those that the generator of pdfme 6.2.4 gives
// with the same font (test/pdfme/gen.py makes the templates and the data).
func testOptions() *Options {
	return &Options{FontDirs: []string{"../docx/testdata/fonts"}, NoSystemFonts: true}
}

func read(t testing.TB, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func data(t testing.TB, name string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(read(t, "../../test/pdfme/data/"+name), &m); err != nil {
		t.Fatal(err)
	}
	return m
}

// pageText returns the text of the pages of a document.
func pageText(t *testing.T, doc *bdf.Document) []string {
	t.Helper()
	st, err := doc.SearchText()
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, pg := range st.Views[0].Pages {
		out = append(out, pg.Text)
	}
	return out
}

func fill(t *testing.T, name string, td *TemplateData) *Result {
	t.Helper()
	b := read(t, "testdata/"+name)
	res, err := Fill(bytes.NewReader(b), int64(len(b)), td, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	return res
}

// TestInspect reads the fields of a template as the names that data has.
func TestInspect(t *testing.T) {
	b := read(t, "testdata/invoice.json")
	fields, err := Inspect(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range fields {
		s := strings.Join(append(append([]string{}, f.In...), f.Name), ">") + " " + f.Kind
		for _, p := range []string{f.Type, f.Mode, f.Code} {
			if p != "" {
				s += " " + p
			}
		}
		if len(f.Options) > 0 {
			s += fmt.Sprintf(" %v", f.Options)
		}
		if f.Default != nil {
			s += fmt.Sprintf(" default=%v", f.Default)
		}
		if f.Multiline {
			s += " multiline"
		}
		got = append(got, s)
	}
	want := []string{
		"number value string default=INV-0001",
		"issued value date",
		"customer.company value string",
		"customer.person value string",
		"note value string default=いつもありがとうございます。 multiline",
		"paid value boolean default=false",
		"method value string [カード ふりこみ] default=カード",
		"express value boolean default=false",
		"items section rows",
		"items>項目 value",
		"items>数量 value",
		"items>単価 value",
		"items>金額 value",
		"qr image qr default=https://example.com/",
		"code image code128 default=INV0001",
		"terms section",
		"terms>. value",
		"confirmed value boolean default=false",
		"total value", // (what a text of the template shows: "総額 {total}")
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("fields:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	for _, f := range fields {
		if len(f.At) != 1 || f.At[0] != "page 1" {
			t.Errorf("%s is at %v", f.Name, f.At)
		}
	}
}

// TestDefaults reads what the designer left in the fields as the sample
// of the data.
func TestDefaults(t *testing.T) {
	b := read(t, "testdata/invoice.json")
	got, err := Defaults(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		t.Fatal(err)
	}
	out, _ := json.Marshal(got)
	want := `{"code":"INV0001","confirmed":false,"customer":{"company":"カメラのタナカ","person":"タナカ ハナコ"},"express":false,"issued":"2026/10/14",` +
		`"items":[{"単価":"300","数量":"2","金額":"600","項目":"ノート"},{"単価":"120","数量":"10","金額":"1,200","項目":"ペン"}],"method":"カード",` +
		`"note":"いつもありがとうございます。","number":"INV-0001","paid":false,"qr":"https://example.com/","terms":["30 日まで","\tカード","\tふりこみ","返すのは 7 日まで"],"total":""}`
	if string(out) != want {
		t.Errorf("defaults:\n%s\nwant:\n%s", out, want)
	}
	// the template fills with them, without a warning
	res := fill(t, "invoice.json", &TemplateData{Data: got, Strict: true})
	if len(res.Warnings) != 0 || res.Pages != 1 {
		t.Errorf("%d pages, warnings %q", res.Pages, res.Warnings)
	}
}

// TestFill fills a template on a blank base: the table grows with its
// rows over two pages, with the names of its columns again on the second,
// and what is under it moves there.
func TestFill(t *testing.T) {
	res := fill(t, "invoice.json", &TemplateData{Data: data(t, "invoice.json"), Strict: true})
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	if res.Pages != 2 || len(res.Doc.Views[0].Pages) != 2 {
		t.Fatalf("%d pages", res.Pages)
	}
	pg := res.Doc.Views[0].Pages[0]
	if math.Abs(float64(pg.W)-210*mm) > 0.01 || math.Abs(float64(pg.H)-297*mm) > 0.01 {
		t.Errorf("page of %g × %g", pg.W, pg.H)
	}
	text := pageText(t, res.Doc)
	for i, want := range [][]string{
		{"ACME Trading", "インボイス", "INV-2026-1014", "2026.10.14", "カメラのタナカ", "タナカ ハナコ さま", "ふりこみ", "項目", "ストラップ X", "1 / 2"},
		{"ACME Trading", "項目", "数量", "単価", "金額", "ノート Y", "ペン H", "総額 46,200", "1.", "返すのは 7 日まで", "COPY", "2 / 2"},
	} {
		for _, w := range want {
			if !strings.Contains(text[i], w) {
				t.Errorf("page %d has no %q:\n%s", i+1, w, text[i])
			}
		}
	}
	for _, not := range []string{"ノート Y", "総額", "COPY"} {
		if strings.Contains(text[0], not) {
			t.Errorf("page 1 has %q", not)
		}
	}
	if strings.Contains(text[1], "インボイス") || strings.Contains(text[1], "ストラップ X") {
		t.Errorf("page 2 has what is on page 1:\n%s", text[1])
	}
	// the note is longer than its field at its size: the size that fits
	if !strings.Contains(text[0], "このたびのごちゅうもんを表にまとめました。") {
		t.Errorf("the note is not whole:\n%s", text[0])
	}
}

// TestLayout checks the heights and the pages of what grows against those
// of pdfme.
func TestLayout(t *testing.T) {
	b := read(t, "testdata/invoice.json")
	tpl, err := parse(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		t.Fatal(err)
	}
	f := &filler{t: tpl, td: &TemplateData{}, warned: map[string]bool{}}
	f.data, _ = tag.Normalize(data(t, "invoice.json")).(map[string]any)
	f.r = formbox.New(bdf.NewDocument(), formbox.Options{FontDirs: []string{"../docx/testdata/fonts"}, NoSystemFonts: true})
	f.expressions()
	var items *schema
	for _, s := range tpl.pages[0] {
		if s.name == "items" {
			items = s
		}
	}
	tl := f.tableOf(items, f.content(items, false), 0, -1, false)
	// a row of a line: 10 points and the room of 2 mm and a border of 0.1 mm
	// over and under; the fourth row has two lines
	const one = 10*pt2mm + 4.2
	if len(tl.body) != 34 || math.Abs(tl.head.h-one) > 1e-9 || math.Abs(tl.body[0].h-one) > 1e-9 || math.Abs(tl.body[3].h-(one+10*pt2mm)) > 1e-9 {
		t.Errorf("%d rows, the head %g, rows %g and %g; want 34, %g, %g and %g", len(tl.body), tl.head.h, tl.body[0].h, tl.body[3].h, one, one, one+10*pt2mm)
	}
	if got := fmt.Sprint(tl.cols); got != "[87.4 26.6 38 38]" {
		t.Errorf("columns %s", got)
	}
	pages := f.layout()
	if len(pages) != 2 {
		t.Fatalf("%d pages", len(pages))
	}
	at := func(page int, name string) *schema {
		for _, s := range pages[page] {
			if s.name == name {
				return s
			}
		}
		return nil
	}
	// the table: 24 rows on the first page, 10 on the second, from its top
	p1, p2 := at(0, "items"), at(1, "items")
	if p1 == nil || p2 == nil || p1.from != 0 || p1.to != 24 || p1.split || p2.from != 24 || p2.to != 34 || !p2.split || p2.y != 20 {
		t.Fatalf("the table: %+v and %+v", p1, p2)
	}
	if want := 25*one + 10*pt2mm; math.Abs(p1.h-want) > 1e-9 {
		t.Errorf("the first part is %g mm tall, want %g", p1.h, want)
	}
	if want := 11 * one; math.Abs(p2.h-want) > 1e-9 {
		t.Errorf("the second part is %g mm tall, want %g", p2.h, want)
	}
	// what was under the table is as far under its end
	total := at(1, "totalText")
	if total == nil || at(0, "totalText") != nil || math.Abs(total.y-(20+11*one+(113-108))) > 1e-9 {
		t.Errorf("the total: %+v", total)
	}
	// what is over the table stays
	if s := at(0, "title"); s == nil || s.y != 20 || at(1, "title") != nil {
		t.Errorf("the title: %+v", s)
	}
}

// TestBasePDF fills a template over a base PDF: its pages are those of
// the PDF, with the values over them.
func TestBasePDF(t *testing.T) {
	res := fill(t, "certificate.json", &TemplateData{Data: data(t, "certificate.json"), Strict: true})
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	if res.Pages != 1 {
		t.Fatalf("%d pages", res.Pages)
	}
	pg := res.Doc.Views[0].Pages[0]
	if math.Abs(float64(pg.W)-595.28) > 0.01 || math.Abs(float64(pg.H)-419.53) > 0.01 || len(pg.Layers) != 2 {
		t.Errorf("page of %g × %g with %d layers", pg.W, pg.H, len(pg.Layers))
	}
	text := pageText(t, res.Doc)[0]
	for _, want := range []string{"Certificate of Completion", "Hanako Tanaka", "Document Formats 101", "Score 92 / 100 (with distinction)", "October 14, 2026", "With honours"} {
		if !strings.Contains(text, want) {
			t.Errorf("the page has no %q:\n%s", want, text)
		}
	}
	// without the data: the fields are empty, with warnings, or an error
	res = fill(t, "certificate.json", nil)
	if len(res.Warnings) != 8 || !strings.Contains(res.Warnings[0], `the data has no "name"`) ||
		!strings.Contains(res.Warnings[6], `the data has no "points", which a text of the template shows`) {
		t.Errorf("warnings: %q", res.Warnings)
	}
	if text := pageText(t, res.Doc)[0]; strings.Contains(text, "Jane Doe") || !strings.Contains(text, "Certificate of Completion") {
		t.Errorf("the page without data:\n%s", text)
	}
	b := read(t, "testdata/certificate.json")
	if _, err := Fill(bytes.NewReader(b), int64(len(b)), &TemplateData{Strict: true}, testOptions()); err == nil || !strings.Contains(err.Error(), `the data has no "name"`) {
		t.Errorf("strict: %v", err)
	}
}

// blank makes a template of a blank page with fields.
func blank(fields ...string) []byte {
	return []byte(`{"basePdf":{"width":100,"height":50,"padding":[5,5,5,5]},"schemas":[[` + strings.Join(fields, ",") + `]]}`)
}

// TestTemplates reads the forms of templates, and tells what is none.
func TestTemplates(t *testing.T) {
	// the older form: the fields of a page by their names
	old := []byte(`{"basePdf":{"width":100,"height":50},"schemas":[{"b":{"type":"text","position":{"x":1,"y":2},"width":30,"height":8},
		"a":{"type":"text","position":{"x":1,"y":12},"width":30,"height":8}}]}`)
	fields, err := Inspect(bytes.NewReader(old), int64(len(old)))
	if err != nil || len(fields) != 2 || fields[0].Name != "b" || fields[1].Name != "a" {
		t.Errorf("the older form: %v, %v", fields, err)
	}
	res, err := Fill(bytes.NewReader(old), int64(len(old)), &TemplateData{Data: map[string]any{"a": "first", "b": 12.5}}, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if text := pageText(t, res.Doc)[0]; !strings.Contains(text, "12.5") || !strings.Contains(text, "first") {
		t.Errorf("the older form:\n%s", text)
	}
	for _, c := range []struct{ src, want string }{
		{`[1]`, "not JSON"},
		{`{"a":1}`, "no schemas"},
		{`{"schemas":[[]]}`, "no base PDF"},
		{`{"schemas":[[]],"basePdf":"https://example.com/base.pdf"}`, "not in the template"},
		{`{"schemas":[[]],"basePdf":"data:application/pdf;base64,!!"}`, "not base64"},
		{`{"schemas":[[]],"basePdf":"data:application/pdf;base64,aGVsbG8="}`, "the base PDF"},
		{`{"schemas":[[]],"basePdf":{"width":0,"height":10}}`, "0 × 10 mm"},
		{`{"schemas":[5],"basePdf":{"width":10,"height":10}}`, "page 1 of the template"},
	} {
		_, err := Fill(strings.NewReader(c.src), int64(len(c.src)), nil, testOptions())
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v, want %q", c.src, err, c.want)
		}
	}
	if !IsTemplate([]byte(" \n{\"basePdf\":\"data:")) || !IsTemplate([]byte("\xef\xbb\xbf{\"schemas\":[")) || IsTemplate([]byte(`{"items":[]}`)) || IsTemplate([]byte("%PDF-1.7")) {
		t.Error("IsTemplate")
	}
}

// TestValues fills fields with the values of data of each kind.
func TestValues(t *testing.T) {
	tpl := blank(
		`{"name":"n","type":"text","position":{"x":5,"y":5},"width":40,"height":6}`,
		`{"name":"when","type":"dateTime","format":"yyyy-MM-dd HH:mm","position":{"x":5,"y":12},"width":40,"height":6}`,
		`{"name":"at","type":"time","position":{"x":50,"y":12},"width":20,"height":6}`,
		`{"name":"mv","type":"multiVariableText","text":"Dear {who}, {n} items","variables":["who","n"],"position":{"x":5,"y":19},"width":60,"height":6}`,
		`{"name":"ro","type":"text","readOnly":true,"content":"{n * 2} and {missing.x} on {date}","position":{"x":5,"y":26},"width":80,"height":6}`,
		`{"name":"pic","type":"image","position":{"x":70,"y":5},"width":10,"height":10}`,
		`{"name":"pdf","type":"pdf417","position":{"x":82,"y":5},"width":10,"height":10}`,
		`{"name":"bad","type":"ean13","position":{"x":82,"y":20},"width":10,"height":10}`,
		`{"name":"odd","type":"hologram","position":{"x":5,"y":33},"width":10,"height":6}`,
	)
	res, err := Fill(bytes.NewReader(tpl), int64(len(tpl)), &TemplateData{Now: "2026/10/14 09:30", Data: map[string]any{
		"n": 21, "when": "2026-10-14T18:05:00", "at": "07:45", "who": "Hanako", "pic": "missing.png", "pdf": "x", "bad": "12", "odd": "x",
	}}, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	text := pageText(t, res.Doc)[0]
	for _, want := range []string{"21", "2026-10-14 18:05", "07:45", "Dear Hanako, 21 items", "42 and {missing.x} on 2026/10/14"} {
		if !strings.Contains(text, want) {
			t.Errorf("the page has no %q:\n%s", want, text)
		}
	}
	got := strings.Join(res.Warnings, "\n")
	for _, want := range []string{`the data has no "missing"`, `an expression of "ro" is not evaluated`, `the picture of "pic" is left out`, `the pdf417 code of "pdf" is not drawn`,
		`the code of "bad" is left out`, `a field of the type "hologram" is not drawn`} {
		if !strings.Contains(got, want) {
			t.Errorf("no warning %q in:\n%s", want, got)
		}
	}
	if len(res.Warnings) != 6 {
		t.Errorf("%d warnings:\n%s", len(res.Warnings), got)
	}
	// without the clock, the date that a text shows is a value of the data
	fields, _ := Inspect(bytes.NewReader(tpl), int64(len(tpl)))
	var names []string
	for _, f := range fields {
		names = append(names, f.Name)
	}
	if got := strings.Join(names, " "); got != "n when at mv.who mv.n pic pdf bad missing date" {
		t.Errorf("fields: %s", got)
	}
	res, err = Fill(bytes.NewReader(tpl), int64(len(tpl)), &TemplateData{Data: map[string]any{"n": "3", "date": "2026/01/02", "missing": map[string]any{"x": "there"}}}, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if text := pageText(t, res.Doc)[0]; !strings.Contains(text, "6 and there on 2026/01/02") {
		t.Errorf("the page:\n%s", text)
	}
}

// TestFormatDate writes dates in the formats of date-fns.
func TestFormatDate(t *testing.T) {
	d := time.Date(2026, 10, 4, 15, 7, 9, 0, time.UTC)
	for _, c := range []struct{ format, locale, want string }{
		{"yyyy/MM/dd", "en", "2026/10/04"},
		{"yy-M-d", "en", "26-10-4"},
		{"MMMM do, yyyy", "en", "October 4th, 2026"},
		{"EEE, MMM d", "en", "Sun, Oct 4"},
		{"EEEE", "ja", "日曜日"},
		{"yyyy年M月d日(E)", "ja", "2026年10月4日(日)"},
		{"hh:mm:ss a", "en", "03:07:09 PM"},
		{"H:m a", "ja", "15:7 午後"},
		{"'at' HH 'o''clock'", "en", "at 15 o'clock"},
	} {
		if got := formatDate(d, c.format, c.locale); got != c.want {
			t.Errorf("%s (%s): %q, want %q", c.format, c.locale, got, c.want)
		}
	}
}

// TestColors reads the colors of templates.
func TestColors(t *testing.T) {
	for s, want := range map[string]bdf.Color{"#fff": bdf.RGB(255, 255, 255), "#1234": bdf.RGBA(0x11, 0x22, 0x33, 0x44), "#2980ba": bdf.RGB(0x29, 0x80, 0xba),
		"00000080": bdf.RGBA(0, 0, 0, 0x80)} {
		if got, ok := colorOf(s); !ok || got != want {
			t.Errorf("%s: %08x %v, want %08x", s, uint32(got), ok, uint32(want))
		}
	}
	for _, s := range []string{"", "#12", "red", "#ggg", "#00000000"} {
		if _, ok := colorOf(s); ok {
			t.Errorf("%q is a color", s)
		}
	}
}

// FuzzFill fills what may be a template: it does not crash, and ends.
func FuzzFill(f *testing.F) {
	for _, name := range []string{"invoice.json", "certificate.json"} {
		f.Add(read(f, "testdata/"+name), read(f, "../../test/pdfme/data/"+name))
	}
	f.Add(blank(`{"name":"t","type":"table","position":{"x":5,"y":40},"width":90,"height":10,"head":["a","b"],"repeatHead":true,"showHead":true}`,
		`{"name":"x","type":"text","overflow":"expand","dynamicFontSize":{"min":1,"max":900},"position":{"x":5,"y":5},"width":2,"height":2,"rotate":33}`,
		`{"name":"l","type":"list","position":{"x":5,"y":30},"width":50,"height":5}`),
		[]byte(`{"t":[[1,2],["a\nb","c"],{"a":true}],"x":"some words that wrap","l":["a","\tb"]}`))
	opts := &Options{NoSystemFonts: true, FontDirs: []string{"../docx/testdata/fonts"}, NoTextIndex: true}
	f.Fuzz(func(t *testing.T, tpl, js []byte) {
		if len(tpl) > 1<<16 || len(js) > 1<<14 {
			return
		}
		var m map[string]any
		json.Unmarshal(js, &m)
		Inspect(bytes.NewReader(tpl), int64(len(tpl)))
		Defaults(bytes.NewReader(tpl), int64(len(tpl)))
		res, err := Fill(bytes.NewReader(tpl), int64(len(tpl)), &TemplateData{Data: m}, opts)
		if err == nil && res.Doc == nil {
			t.Fatal("no document")
		}
	})
}

// countingFS counts the files read of it.
type countingFS struct {
	fstest.MapFS
	reads int
}

func (c *countingFS) ReadFile(name string) ([]byte, error) {
	c.reads++
	return c.MapFS.ReadFile(name)
}

// TestPicturesOnce reads a picture that the data names once, however many
// pages show it: a picture of the base page is on every page of a table
// that goes on over pages.
func TestPicturesOnce(t *testing.T) {
	var rows []string
	for i := 0; i < 60; i++ {
		rows = append(rows, fmt.Sprintf(`["item %d","%d"]`, i+1, (i+1)*100))
	}
	tpl := []byte(`{"basePdf":{"width":100,"height":50,"padding":[5,5,5,5],"staticSchema":[
		{"name":"logo","type":"image","readOnly":true,"content":"logo.png","position":{"x":80,"y":2},"width":15,"height":3}]},
		"schemas":[[{"name":"items","type":"table","position":{"x":5,"y":6},"width":90,"height":10,"showHead":false,"head":["name","price"],
		"headWidthPercentages":[70,30],"bodyStyles":{"fontSize":8,"padding":1},"content":"[]"},
		{"name":"photo","type":"image","position":{"x":5,"y":30},"width":10,"height":10}]]}`)
	files := &countingFS{MapFS: fstest.MapFS{"logo.png": {Data: tplPNG(30, 6)}, "photo.png": {Data: tplPNG(10, 10)}}}
	res, err := Fill(bytes.NewReader(tpl), int64(len(tpl)), &TemplateData{Files: files, Data: map[string]any{"items": "[" + strings.Join(rows, ",") + "]", "photo": "photo.png"}}, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if res.Pages < 3 || files.reads != 2 || len(res.Warnings) != 0 {
		t.Errorf("%d pages, %d files read, warnings %q", res.Pages, files.reads, res.Warnings)
	}
	images := 0
	for _, part := range res.Doc.Parts() {
		if part.Type == bdf.PartImage {
			images++
		}
	}
	if images != 2 {
		t.Errorf("%d images in the document, want 2", images)
	}
	// a picture that is not read is told once, for all the pages
	files.reads = 0
	res, err = Fill(bytes.NewReader(tpl), int64(len(tpl)), &TemplateData{Files: files, Data: map[string]any{"items": "[" + strings.Join(rows, ",") + "]", "photo": "missing.png"}}, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if files.reads != 1 || len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], `the picture of "photo" is left out`) {
		t.Errorf("%d files read, warnings %q", files.reads, res.Warnings)
	}
}

// tplPNG makes a picture of a size.
func tplPNG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < w*h; i++ {
		img.Set(i%w, i/w, color.RGBA{200, uint8(i), 30, 255})
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}
