package form

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/form/tag"
)

func testConvert() *converter.Options {
	return &converter.Options{FontDirs: []string{"../converter/pptx/testdata/fonts"}, NoSystemFonts: true}
}

func template(t *testing.T, name string) *bytes.Reader {
	t.Helper()
	b, err := os.ReadFile("../converter/xlsx/testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(b)
}

func text(t *testing.T, doc *bdf.Document) string {
	t.Helper()
	st, err := doc.SearchText()
	if err != nil {
		t.Fatal(err)
	}
	var b strings.Builder
	for _, v := range st.Views {
		for _, p := range v.Pages {
			b.WriteString(p.Text + "\n")
		}
	}
	return b.String()
}

func TestLoadData(t *testing.T) {
	csv := "name,qty,code,ok,note\nA,2,007,true,\n\"B, the second\",1.5,12,FALSE,x\n"
	sjis := []byte{0x96, 0xbc, 0x91, 0x4f, ',', 'n', '\n', 0x8e, 0x52, 0x93, 0x63, ',', '3', '\n'} // 名前,n / 山田,3 in Shift_JIS
	data, err := LoadData(
		Source{FileName: "base.json", Data: []byte(`{"no": 12345678901234567890, "customer": {"name": "ACME"}, "old": 1, "big": 1e3}`)},
		Source{FileName: "more.yaml", Data: []byte("old: 2\ndate: 2026-10-13\ntags: [a, b]\nnested:\n  deep: true\n")},
		Source{Name: "items", FileName: "lines.csv", Data: []byte(csv)},
		Source{FileName: "dir/people.csv", Data: sjis},
		Source{Name: "customer.zip", FileName: "zip.json", Data: []byte(`"100-0001"`)},
	)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"no": 12345678901234567890.0, "old": 2.0, "big": 1000.0,
		"customer": map[string]any{"name": "ACME", "zip": "100-0001"},
		"date":     time.Date(2026, 10, 13, 0, 0, 0, 0, time.UTC),
		"tags":     []any{"a", "b"},
		"nested":   map[string]any{"deep": true},
		"items": []any{
			map[string]any{"name": "A", "qty": 2.0, "code": "007", "ok": true, "note": nil},
			map[string]any{"name": "B, the second", "qty": 1.5, "code": 12.0, "ok": false, "note": "x"},
		},
		"people": []any{map[string]any{"名前": "山田", "n": 3.0}},
	}
	if !reflect.DeepEqual(data, want) {
		t.Errorf("data = %#v", data)
	}
	for _, c := range []struct {
		src  Source
		want string
	}{
		{Source{FileName: "a.json", Data: []byte(`[1]`)}, "is an object"},
		{Source{FileName: "a.json", Data: []byte(`{"a":`)}, "a.json"},
		{Source{FileName: "a.json", Data: []byte(`{} {}`)}, "more than one"},
		{Source{FileName: "a.yaml", Data: []byte("a: 1\n---\nb: 2\n")}, "more than one"},
		{Source{FileName: "a.toml", Data: []byte(`a = 1`)}, "data is JSON"},
		{Source{Name: "x..y", FileName: "a.json", Data: []byte(`1`)}, "empty part"},
		{Source{Name: "l", FileName: "a.csv", Data: []byte("a,a\n1,2\n")}, `two columns are named "a"`},
	} {
		if _, err := LoadData(c.src); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("LoadData(%s): %v, want an error with %q", c.src.FileName, err, c.want)
		}
	}
}

// TestLoadWorkbook reads a workbook as data: a list for each sheet.
func TestLoadWorkbook(t *testing.T) {
	b, err := os.ReadFile("../converter/xlsx/testdata/print.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	data, err := LoadData(Source{FileName: "print.xlsx", Data: b})
	if err != nil {
		t.Fatal(err)
	}
	wide, ok := data["Wide"].([]any)
	if !ok || len(wide) != 40 {
		t.Fatalf("the sheet Wide: %T of %d", data["Wide"], len(wide))
	}
	if row := wide[1].(map[string]any); row["行"] != "Row 2" || row["C1"] != 201.0 || row["C13"] != 213.0 {
		t.Errorf("the second row: %v", row)
	}
	// with a name, the sheets are under it
	data, err = LoadData(Source{Name: "book", FileName: "print.xlsx", Data: b})
	if _, ok := data["book"].(map[string]any)["Fit"].([]any); err != nil || !ok {
		t.Errorf("the sheets under a name: %v, %v", data, err)
	}
}

func TestFill(t *testing.T) {
	data, err := LoadData(
		Source{FileName: "order.yaml", Data: []byte("no: A-7\ncustomer: ACME\nnote: Thank you.\n")},
		Source{Name: "items", FileName: "lines.csv", Data: []byte("name,qty,price\nPen,2,120\nInk,1,800\n")},
	)
	if err != nil {
		t.Fatal(err)
	}
	r := template(t, "form-rows.xlsx")
	res, err := Fill(r, r.Size(), &Options{Data: data, Convert: testConvert()})
	if err != nil {
		t.Fatal(err)
	}
	got := text(t, res.Doc)
	for _, want := range []string{"Invoice A-7", "To: ACME", "Pen\n2\n120\n240", "Ink\n1\n800\n800", "Total\n1,040", "Tax\n104", "Sum\n1,144"} {
		if !strings.Contains(got, want) {
			t.Errorf("the form lacks %q:\n%s", want, got)
		}
	}
	if res.Format != FormatXLSX || len(res.Fields) == 0 {
		t.Errorf("format %q, %d fields", res.Format, len(res.Fields))
	}
	// what is not a template
	for name, want := range map[string]string{
		"form_test.go":         "not an Excel workbook",
		"../testdata/demo.bdf": "not a compiled template",
	} {
		b, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := Fill(bytes.NewReader(b), int64(len(b)), nil); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("Fill(%s): %v, want an error with %q", name, err, want)
		}
	}
}

// TestFillOffice fills the Word and PowerPoint templates with the data of
// the Excel ones: the same data, in the forms of each.
func TestFillOffice(t *testing.T) {
	load := func(name string) map[string]any {
		b, err := os.ReadFile("../test/xlsx/data/" + name + ".json")
		if err != nil {
			t.Fatal(err)
		}
		data, err := LoadData(Source{FileName: name + ".json", Data: b})
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	for _, c := range []struct {
		file, format string
		pages        int
		want         []string
	}{
		{"docx/testdata/form-fixed.docx", FormatDOCX, 1, []string{"明細書 Certificate", "No. 1024", "Date", "2026/10/13 2026-10", "ACME さま", "12,000", "13,200.00", "■ paid", "Thank you.\nSecond line", "Tags: a, b, c", "- tag 3: c"}},
		{"docx/testdata/form-rows.docx", FormatDOCX, 2, []string{"Invoice A-7", "To: ACME", "1\nItem 01\n2\n100", "60\nItem 60\n1\n6,000", "Thank you."}},
		{"docx/testdata/form-pages.docx", FormatDOCX, 3, []string{"Delivery D-1", "3 / 3", "From 210,000", "24\nPart 24\n24\n24,000\n-- end --", "Total 300,000"}},
		{"pptx/testdata/form-fixed.pptx", FormatPPTX, 1, []string{"明細書 Certificate", "No. 1024", "Date 2026/10/13 2026-10", "Name ACME さま", "Amount 12,000   Total 13,200.00", "■ paid", "Memo: Thank you.\nSecond line", "Tags: a, b, c", "- tag 3: c"}},
		{"pptx/testdata/form-rows.pptx", FormatPPTX, 4, []string{"Invoice A-7", "4 / 4", "1\nItem 01\n2\n100", "60\nItem 60\n1\n6,000", "Thank you."}},
		{"pptx/testdata/form-pages.pptx", FormatPPTX, 3, []string{"Delivery D-1", "3 / 3", "From 210,000", "24\nPart 24\n24\n24,000\n-- end --", "Total 300,000"}},
	} {
		tpl, err := os.ReadFile("../converter/" + c.file)
		if err != nil {
			t.Fatal(err)
		}
		name := c.file[strings.LastIndex(c.file, "/")+1 : strings.LastIndex(c.file, ".")]
		res, err := Fill(bytes.NewReader(tpl), int64(len(tpl)), &Options{Data: load(name), Files: os.DirFS("../test/xlsx/data"), Convert: testConvert()})
		if err != nil {
			t.Errorf("%s: %v", c.file, err)
			continue
		}
		got := text(t, res.Doc)
		for _, want := range c.want {
			if !strings.Contains(got, want) {
				t.Errorf("%s: the form lacks %q:\n%s", c.file, want, got)
			}
		}
		if n := len(res.Doc.Views[0].Pages); res.Format != c.format || len(res.Doc.Views) != 1 || n != c.pages || len(res.Fields) == 0 {
			t.Errorf("%s: format %q, %d views, %d pages (want %d), %d fields", c.file, res.Format, len(res.Doc.Views), n, c.pages, len(res.Fields))
		}
		for _, w := range res.Warnings {
			t.Errorf("%s: warning: %s", c.file, w)
		}
		// compiled, it fills to the same form, and shows its tags without a sample
		cmp, err := Compile(bytes.NewReader(tpl), int64(len(tpl)), name, &Options{Convert: testConvert()})
		if err != nil {
			t.Errorf("%s: compile: %v", c.file, err)
			continue
		}
		if got := text(t, cmp.Doc); !strings.Contains(got, "{{") || len(cmp.Doc.Views) != 1 || cmp.Template.Format != c.format {
			t.Errorf("%s: the pages of the compiled template (%d views):\n%s", c.file, len(cmp.Doc.Views), got)
		}
		if !reflect.DeepEqual(cmp.Template.Fields, res.Fields) {
			t.Errorf("%s: the fields of the compiled template: %+v\nof the filled one: %+v", c.file, cmp.Template.Fields, res.Fields)
		}
		var bdft, direct, viaBDFT bytes.Buffer
		if err := cmp.Doc.WriteSingle(&bdft); err != nil {
			t.Fatal(err)
		}
		again, err := Fill(bytes.NewReader(bdft.Bytes()), int64(bdft.Len()), &Options{Data: load(name), Files: os.DirFS("../test/xlsx/data"), Convert: testConvert()})
		if err != nil {
			t.Errorf("%s: filling the compiled template: %v", c.file, err)
			continue
		}
		res.Doc.WriteSingle(&direct)
		again.Doc.WriteSingle(&viaBDFT)
		if !bytes.Equal(direct.Bytes(), viaBDFT.Bytes()) {
			t.Errorf("%s: the compiled template fills to another form than the template", c.file)
		}
	}
}

// TestCompile compiles a template and fills the compiled template: the
// same form as the template gives.
func TestCompile(t *testing.T) {
	sample := map[string]any{"no": "SAMPLE", "customer": "Sample Co.", "note": "", "items": []any{map[string]any{"name": "Sample item", "qty": 1.0, "price": 100.0}}}
	r := template(t, "form-rows.xlsx")
	c, err := Compile(r, r.Size(), "form-rows.xlsx", &Options{Data: sample, Convert: testConvert()})
	if err != nil {
		t.Fatal(err)
	}
	if got := text(t, c.Doc); !strings.Contains(got, "Invoice SAMPLE") || !strings.Contains(got, "Sample item") {
		t.Errorf("the pages of the compiled template:\n%s", got)
	}
	var bdft bytes.Buffer
	if err := c.Doc.WriteSingle(&bdft); err != nil {
		t.Fatal(err)
	}
	// a reader of documents opens it as a document
	rd, err := bdf.OpenSingle(bytes.NewReader(bdft.Bytes()), int64(bdft.Len()))
	if err != nil {
		t.Fatal(err)
	}
	if doc, err := rd.ToDocument(); err != nil || len(doc.Views) == 0 || len(doc.Views[0].Pages) == 0 {
		t.Errorf("the compiled template as a document: %v", err)
	}
	// its description
	def, src, err := ReadTemplate(bytes.NewReader(bdft.Bytes()), int64(bdft.Len()))
	if err != nil {
		t.Fatal(err)
	}
	orig, _ := os.ReadFile("../converter/xlsx/testdata/form-rows.xlsx")
	if def.Format != FormatXLSX || def.Name != "form-rows.xlsx" || !bytes.Equal(src, orig) || !reflect.DeepEqual(def.Sample, sample) {
		t.Errorf("the description: %+v", def)
	}
	fields, err := Fields(bytes.NewReader(bdft.Bytes()), int64(bdft.Len()))
	direct, _ := Fields(bytes.NewReader(orig), int64(len(orig)))
	if err != nil || len(fields) == 0 || !reflect.DeepEqual(fields, direct) {
		t.Errorf("the fields of the compiled template: %+v\nof the template: %+v", fields, direct)
	}
	// it holds the JSON Schema of its data: that of the fields and the sample
	want, err := Schema(direct, sample)
	if err != nil {
		t.Fatal(err)
	}
	var held, made any
	if json.Unmarshal(def.Schema, &held) != nil || json.Unmarshal(want, &made) != nil || !reflect.DeepEqual(held, made) || !bytes.Contains(def.Schema, []byte(`"items"`)) {
		t.Errorf("the schema of the compiled template:\n%s\nwant\n%s", def.Schema, want)
	}
	if got, err := SchemaOf(bytes.NewReader(bdft.Bytes()), int64(bdft.Len()), nil); err != nil || !bytes.Equal(got, want) {
		t.Errorf("SchemaOf a compiled template: %v\n%s\nwant\n%s", err, got, want)
	}
	// (of a template that is not compiled, and with another sample, it is made of the tags)
	bare, _ := Schema(direct, nil)
	if got, err := SchemaOf(bytes.NewReader(orig), int64(len(orig)), nil); err != nil || !bytes.Equal(got, bare) || bytes.Equal(bare, want) {
		t.Errorf("SchemaOf a template: %v\n%s", err, got)
	}
	other := map[string]any{"no": 7.0, "items": []any{map[string]any{"name": 1.0}}}
	theirs, _ := Schema(direct, other)
	if got, err := SchemaOf(bytes.NewReader(bdft.Bytes()), int64(bdft.Len()), other); err != nil || !bytes.Equal(got, theirs) || bytes.Equal(theirs, want) {
		t.Errorf("SchemaOf a compiled template with a sample: %v\n%s", err, got)
	}
	var items tag.Field
	for _, f := range fields {
		if f.Name == "items" {
			items = f
		}
	}
	if items.Kind != tag.FieldSection || items.Mode != tag.ModeRows {
		t.Errorf("the field of the items: %+v", items)
	}
	// filled, it is the form of the template
	data := map[string]any{"no": 1.0, "customer": "ACME", "note": "n", "items": []any{map[string]any{"name": "Pen", "qty": 3.0, "price": 50.0}}}
	a, err := Fill(bytes.NewReader(bdft.Bytes()), int64(bdft.Len()), &Options{Data: data, Convert: testConvert()})
	if err != nil {
		t.Fatal(err)
	}
	b, err := Fill(bytes.NewReader(orig), int64(len(orig)), &Options{Data: data, Convert: testConvert()})
	if err != nil {
		t.Fatal(err)
	}
	var wa, wb bytes.Buffer
	a.Doc.WriteSingle(&wa)
	b.Doc.WriteSingle(&wb)
	if !bytes.Equal(wa.Bytes(), wb.Bytes()) || !strings.Contains(text(t, a.Doc), "Pen\n3\n50\n150") {
		t.Errorf("the compiled template fills to another form than the template (%d and %d bytes)", wa.Len(), wb.Len())
	}
	// a compiled template is not compiled again, and a mistake is told when it is compiled
	if _, err := Compile(bytes.NewReader(bdft.Bytes()), int64(bdft.Len()), "", nil); err == nil || !strings.Contains(err.Error(), "already") {
		t.Errorf("compiling a compiled template: %v", err)
	}
	// without a sample the pages show the tags
	r = template(t, "form-fixed.xlsx")
	c, err = Compile(r, r.Size(), "", &Options{Convert: testConvert()})
	if err != nil {
		t.Fatal(err)
	}
	if got := text(t, c.Doc); !strings.Contains(got, "{{no}}") || !strings.Contains(got, "{{customer.name}} さま") || len(c.Template.Fields) == 0 {
		t.Errorf("the pages without a sample:\n%s", got)
	}
}

// TestSchema makes the JSON Schema of the data of templates: what their
// tags suggest, and what a sample of the data settles.
func TestSchema(t *testing.T) {
	props := func(t *testing.T, schema []byte) map[string]any {
		t.Helper()
		var doc map[string]any
		if err := json.Unmarshal(schema, &doc); err != nil {
			t.Fatalf("%v\n%s", err, schema)
		}
		if doc["$schema"] != SchemaDialect || doc["type"] != "object" {
			t.Errorf("the schema starts with %v, %v", doc["$schema"], doc["type"])
		}
		return doc["properties"].(map[string]any)
	}
	// what describes a property, in a line: its type, the format, and the properties of an object or of the items of a list
	var line func(v any) string
	line = func(v any) string {
		m := v.(map[string]any)
		s := fmt.Sprint(m["type"])
		if f, ok := m["format"]; ok {
			s += "(" + f.(string) + ")"
		}
		if n, ok := m["maxItems"]; ok {
			s += fmt.Sprintf("[max %v]", n)
		}
		inner := m
		if items, ok := m["items"].(map[string]any); ok {
			s += " of " + fmt.Sprint(items["type"])
			inner = items
		}
		if p, ok := inner["properties"].(map[string]any); ok {
			names := make([]string, 0, len(p))
			for k := range p {
				names = append(names, k)
			}
			sort.Strings(names)
			for i, k := range names {
				names[i] = k + ": " + line(p[k])
			}
			s += " {" + strings.Join(names, ", ") + "}"
		}
		return s
	}
	check := func(t *testing.T, schema []byte, want map[string]string) {
		t.Helper()
		p := props(t, schema)
		for name, w := range want {
			v, ok := p[name]
			if !ok {
				t.Errorf("no property %q in\n%s", name, schema)
				continue
			}
			if got := line(v); got != w {
				t.Errorf("%s is %q, want %q", name, got, w)
			}
		}
		if len(p) != len(want) {
			t.Errorf("%d properties, want %d:\n%s", len(p), len(want), schema)
		}
	}
	fields := func(file string) []tag.Field {
		r := template(t, file)
		f, err := Fields(r, r.Size())
		if err != nil {
			t.Fatal(err)
		}
		return f
	}
	sample := func(name string) map[string]any {
		b, err := os.ReadFile("../test/xlsx/data/" + name + ".json")
		if err != nil {
			t.Fatal(err)
		}
		data, err := LoadData(Source{FileName: name + ".json", Data: b})
		if err != nil {
			t.Fatal(err)
		}
		return data
	}

	// a form of one page: values, a date by the format of its tag, an object
	// of a name with a dot, a condition, a list of values, a picture, codes
	fixed, err := Schema(fields("form-fixed.xlsx"), nil)
	if err != nil {
		t.Fatal(err)
	}
	check(t, fixed, map[string]string{
		"no": "[string number]", "date": "string(date)", "customer": "object {name: [string number]}", "amount": "[string number]",
		"total": "number", "paid": "boolean", "memo": "[string number]", "photo": "string(uri-reference)", "url": "string",
		"tags": "array of [string number]", "missing": "[string number]", "jan": "string",
	})
	// the properties are in the order of the template, and say what the template says of them
	if i, j, k := bytes.Index(fixed, []byte(`"no"`)), bytes.Index(fixed, []byte(`"customer"`)), bytes.Index(fixed, []byte(`"jan"`)); i < 0 || i > j || j > k {
		t.Errorf("the properties are not in the order of the template:\n%s", fixed)
	}
	p := props(t, fixed)
	if info := p["jan"].(map[string]any)[SchemaKeyword].(map[string]any); info["kind"] != "image" || info["code"] != "ean13" || fmt.Sprint(info["at"]) != "[Form!B16]" {
		t.Errorf("what the schema says of jan: %v", info)
	}
	if info := p["total"].(map[string]any)[SchemaKeyword].(map[string]any); info["kind"] != "value" || info["format"] != "#,##0.00" {
		t.Errorf("what the schema says of total: %v", info)
	}
	// with the sample, the values have the types of their values
	fixed, err = Schema(fields("form-fixed.xlsx"), sample("form-fixed"))
	if err != nil {
		t.Fatal(err)
	}
	check(t, fixed, map[string]string{
		"no": "number", "date": "string(date)", "customer": "object {name: string}", "amount": "number", "total": "number", "paid": "boolean",
		"memo": "string", "photo": "string(uri-reference)", "url": "string", "tags": "array of string", "missing": "string", "jan": "string",
	})

	// a table that grows: a list of objects
	rows, err := Schema(fields("form-rows.xlsx"), nil)
	if err != nil {
		t.Fatal(err)
	}
	check(t, rows, map[string]string{
		"no": "[string number]", "customer": "[string number]", "note": "[string number]",
		"items": "array of object {name: [string number], price: [string number], qty: [string number]}",
	})
	if info := props(t, rows)["items"].(map[string]any)[SchemaKeyword].(map[string]any); info["kind"] != "section" || info["mode"] != "rows" {
		t.Errorf("what the schema says of items: %v", info)
	}

	// a slip of pages: without a sample the names in the section of the
	// pages are those of a page; the sample has some of them around it
	pages, err := Schema(fields("form-pages.xlsx"), nil)
	if err != nil {
		t.Fatal(err)
	}
	const slots = "rows: array[max 10] of object {amount: [string number], name: [string number], no: [string number], qty: [string number]}"
	check(t, pages, map[string]string{
		"pages": "array of object {carry_in: [string number], carry_out: [string number], customer: [string number], " + slots + ", slip: [string number], total: number}",
	})
	pages, err = Schema(fields("form-pages.xlsx"), sample("form-pages"))
	if err != nil {
		t.Fatal(err)
	}
	check(t, pages, map[string]string{
		"slip": "string", "customer": "string", "total": "number",
		"pages": "array of object {carry_in: number, carry_out: number, rows: array[max 10] of object {amount: number, name: string, no: number, qty: number}}",
	})

	// names that a template looks up around a section too are not those of its items;
	// a section whose value the template also writes is that value; a section of a sample's object is an object
	mixed := []tag.Field{
		{Name: "people", Kind: tag.FieldSection, Mode: tag.ModeRows},
		{Name: "name", Kind: tag.FieldValue, In: []string{"people"}},
		{Name: "company", Kind: tag.FieldValue, In: []string{"people"}},
		{Name: "company", Kind: tag.FieldValue},
		{Name: "memo", Kind: tag.FieldSection},
		{Name: "memo", Kind: tag.FieldValue, In: []string{"memo"}},
		{Name: "none", Kind: tag.FieldSection, Inverted: true},
		{Name: "owner", Kind: tag.FieldSection},
		{Name: "phone", Kind: tag.FieldValue, In: []string{"owner"}},
		{Name: "_index", Kind: tag.FieldValue, In: []string{"people"}},
	}
	got, err := Schema(mixed, nil)
	if err != nil {
		t.Fatal(err)
	}
	check(t, got, map[string]string{
		"people": "array of object {name: [string number]}", "company": "[string number]", "memo": "[string number]", "none": "boolean",
		"owner": "array of object {phone: [string number]}",
	})
	if d := props(t, got)["none"].(map[string]any)["description"]; !strings.Contains(fmt.Sprint(d), "false") {
		t.Errorf("an inverted section: %v", d)
	}
	got, err = Schema(mixed, map[string]any{"owner": map[string]any{"phone": "03"}, "memo": "x", "people": []any{}})
	if err != nil {
		t.Fatal(err)
	}
	check(t, got, map[string]string{
		"people": "array of object {name: [string number]}", "company": "[string number]", "memo": "string", "none": "boolean",
		"owner": "object {phone: string}",
	})

	// the formats of numbers and of dates
	for code, want := range map[string]string{
		"#,##0": "number", "0.0%": "number", "0.00E+00": "number", `"¥"#,##0;[Red]-#,##0`: "number", "yyyy/mm/dd": "date", "yyyy年m月d日": "date",
		"ggge年": "date", "m/d": "date", "yyyy-mm-dd hh:mm": "date-time", "h:mm": "", "[h]:mm:ss": "", "General": "", "@": "", "": "", `"no" 000`: "number",
		"[$-411]ggge年m月d日": "date", "mmm": "date", "aaa": "date",
	} {
		if got := formatKind(code); got != want {
			t.Errorf("the format %q is %q, want %q", code, got, want)
		}
	}
}

// TestFillPlaced fills the templates whose fields have places: a PDF with
// the fields of a form, and a template of pdfme.
func TestFillPlaced(t *testing.T) {
	load := func(name string) map[string]any {
		b, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		data, err := LoadData(Source{FileName: "data.json", Data: b})
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	// (the fonts of the Word tests have the characters of the templates)
	convert := &converter.Options{FontDirs: []string{"../converter/docx/testdata/fonts"}, NoSystemFonts: true}
	for _, c := range []struct {
		file, format, data string
		pages              int
		want               []string
		// the pages of the compiled template without a sample
		bare, schema []string
	}{
		{"pdfform/testdata/application.pdf", FormatPDF, "../test/pdf/data/application.json", 2,
			[]string{"Application", "Taro Yamada", "YAMADA TARO", "12800"},
			[]string{"Application"}, []string{`"applicant"`, `"maxLength": 7`, `"enum"`, `"United States"`, `"type": "boolean"`}},
		{"pdfme/testdata/invoice.json", FormatPDFME, "../test/pdfme/data/invoice.json", 2,
			[]string{"インボイス", "INV-2026-1014", "2026.10.14", "ストラップ X", "数量", "単価", "金額", "総額 46,200", "2 / 2"},
			[]string{"インボイス", "INV-0001", "ノート\n2\n300\n600", "1 / 1"}, []string{`"customer"`, `"format": "date"`, `"mode": "rows"`, `"code": "qr"`, `"total"`}},
		{"pdfme/testdata/certificate.json", FormatPDFME, "../test/pdfme/data/certificate.json", 1,
			[]string{"Certificate of Completion", "Hanako Tanaka", "Score 92 / 100 (with distinction)", "October 14, 2026"},
			[]string{"Certificate of Completion", "Jane Doe"}, []string{`"honours"`, `"points"`, `"kind": "image"`}},
	} {
		tpl, err := os.ReadFile("../converter/" + c.file)
		if err != nil {
			t.Fatal(err)
		}
		data := load(c.data)
		res, err := Fill(bytes.NewReader(tpl), int64(len(tpl)), &Options{Data: data, Strict: true, Convert: convert})
		if err != nil {
			t.Errorf("%s: %v", c.file, err)
			continue
		}
		if res.Format != c.format || len(res.Warnings) != 0 || len(res.Doc.Views[0].Pages) != c.pages || len(res.Fields) == 0 {
			t.Errorf("%s: format %s, %d pages, %d fields, warnings %q", c.file, res.Format, len(res.Doc.Views[0].Pages), len(res.Fields), res.Warnings)
		}
		got := text(t, res.Doc)
		for _, w := range c.want {
			if !strings.Contains(got, w) {
				t.Errorf("%s: no %q in\n%s", c.file, w, got)
			}
		}
		fields, err := Fields(bytes.NewReader(tpl), int64(len(tpl)))
		if err != nil || !reflect.DeepEqual(fields, res.Fields) {
			t.Errorf("%s: the fields: %v\n%+v\nof the filled template:\n%+v", c.file, err, fields, res.Fields)
		}
		schema, err := SchemaOf(bytes.NewReader(tpl), int64(len(tpl)), nil)
		if err != nil {
			t.Errorf("%s: %v", c.file, err)
		}
		for _, w := range c.schema {
			if !bytes.Contains(schema, []byte(w)) {
				t.Errorf("%s: no %s in the schema\n%s", c.file, w, schema)
			}
		}
		// compiled with the data as its sample, it fills to the same form
		comp, err := Compile(bytes.NewReader(tpl), int64(len(tpl)), "template", &Options{Data: data, Convert: convert})
		if err != nil {
			t.Errorf("%s: %v", c.file, err)
			continue
		}
		var bdft, direct, again bytes.Buffer
		if err := comp.Doc.WriteSingle(&bdft); err != nil {
			t.Fatal(err)
		}
		def, src, err := ReadTemplate(bytes.NewReader(bdft.Bytes()), int64(bdft.Len()))
		if err != nil || def.Format != c.format || !bytes.Equal(src, tpl) || !reflect.DeepEqual(def.Fields, fields) {
			t.Errorf("%s: the compiled template: %v, %+v", c.file, err, def)
		}
		filled, err := Fill(bytes.NewReader(bdft.Bytes()), int64(bdft.Len()), &Options{Data: data, Convert: convert})
		if err != nil {
			t.Errorf("%s: %v", c.file, err)
			continue
		}
		res.Doc.WriteSingle(&direct)
		filled.Doc.WriteSingle(&again)
		if !bytes.Equal(direct.Bytes(), again.Bytes()) {
			t.Errorf("%s: the compiled template fills to another form (%d and %d bytes)", c.file, direct.Len(), again.Len())
		}
		// compiled without one, its pages are the template as it is: a
		// template of pdfme with what its designer left in its fields
		comp, err = Compile(bytes.NewReader(tpl), int64(len(tpl)), "", &Options{Convert: convert})
		if err != nil {
			t.Errorf("%s: %v", c.file, err)
			continue
		}
		got = text(t, comp.Doc)
		for _, w := range c.bare {
			if !strings.Contains(got, w) {
				t.Errorf("%s: no %q in the pages without a sample\n%s", c.file, w, got)
			}
		}
		if (comp.Template.Sample != nil) != (c.format == FormatPDFME) {
			t.Errorf("%s: the sample of the compiled template: %v", c.file, comp.Template.Sample)
		}
	}
	// what is no template
	for _, src := range []string{`{"items":[]}`, "plain text", "%PDF-1.7\n%%EOF"} {
		if _, err := Fill(strings.NewReader(src), int64(len(src)), nil); err == nil {
			t.Errorf("%q is filled", src)
		}
	}
}
