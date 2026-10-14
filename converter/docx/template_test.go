package docx

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html"
	"image"
	"image/color"
	"image/png"
	"strconv"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/form/tag"
	"github.com/shibukawa/bdf/image/barcode"
)

// tr makes a run of text; a run that starts with "<" is markup.
func tr(text string) string {
	if strings.HasPrefix(text, "<") {
		return text
	}
	return `<w:r><w:rPr><w:sz w:val="20"/></w:rPr><w:t xml:space="preserve">` + html.EscapeString(text) + `</w:t></w:r>`
}

// tp makes a paragraph of runs.
func tp(runs ...string) string {
	var b strings.Builder
	b.WriteString(`<w:p><w:pPr><w:spacing w:before="0" w:after="0"/></w:pPr>`)
	for _, r := range runs {
		b.WriteString(tr(r))
	}
	return b.String() + `</w:p>`
}

// ttbl makes a table of rows of cells; a cell is the text of a paragraph
// (of several, parted by "\n").
func ttbl(rows ...[]string) string {
	var b strings.Builder
	b.WriteString(`<w:tbl><w:tblPr><w:tblW w:w="4000" w:type="dxa"/></w:tblPr><w:tblGrid>`)
	for range rows[0] {
		b.WriteString(`<w:gridCol w:w="` + strconv.Itoa(4000/len(rows[0])) + `"/>`)
	}
	b.WriteString(`</w:tblGrid>`)
	for _, row := range rows {
		b.WriteString(`<w:tr>`)
		for _, cell := range row {
			b.WriteString(`<w:tc><w:tcPr><w:tcW w:w="` + strconv.Itoa(4000/len(row)) + `" w:type="dxa"/></w:tcPr>`)
			for _, line := range strings.Split(cell, "\n") {
				b.WriteString(tp(line))
			}
			b.WriteString(`</w:tc>`)
		}
		b.WriteString(`</w:tr>`)
	}
	return b.String() + `</w:tbl>`
}

// tplFormat stands for the formatter of number formats.
func tplFormat(v any, format string) string {
	if format == "" {
		return fmt.Sprint(v)
	}
	return fmt.Sprintf("[%v as %s]", v, format)
}

func fillDocx(t *testing.T, body string, parts map[string]string, data map[string]any) (*Filled, *bdf.Reader) {
	t.Helper()
	file := docxBytes(t, body, parts)
	res, err := Fill(bytes.NewReader(file), int64(len(file)), &TemplateData{Data: data, Format: tplFormat}, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	return res, reopen(t, &res.Result)
}

func fillDocxErr(t *testing.T, body string, data map[string]any, strict bool) error {
	t.Helper()
	file := docxBytes(t, body, nil)
	_, err := Fill(bytes.NewReader(file), int64(len(file)), &TemplateData{Data: data, Strict: strict}, testOptions())
	if err == nil {
		t.Fatal("no error")
	}
	return err
}

// bodyLines returns the lines of the body of the pages, joined by " / ",
// a string for each page.
func bodyLines(t *testing.T, r *bdf.Reader) []string {
	t.Helper()
	var out []string
	for i := range r.Manifest.Views[0].Pages {
		out = append(out, strings.ReplaceAll(pageText(t, r, 0, i, bdf.RoleBody), "\n", " / "))
	}
	return out
}

func wantBody(t *testing.T, r *bdf.Reader, want ...string) {
	t.Helper()
	got := bodyLines(t, r)
	if len(got) != len(want) {
		t.Fatalf("%d pages, want %d:\n%s", len(got), len(want), strings.Join(got, "\n"))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("page %d:\n got %s\nwant %s", i+1, got[i], want[i])
		}
	}
}

func noWarnings(t *testing.T, res *Filled) {
	t.Helper()
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
}

func TestFillValues(t *testing.T) {
	body := tp("Invoice {{no}} for ", "{{custo", `<w:proofErr w:type="spellStart"/>`, "mer.na", `<w:bookmarkStart w:id="0" w:name="_GoBack"/><w:bookmarkEnd w:id="0"/>`, "me}}", "!") +
		tp(`Total {{total | fmt:"#,##0"}}`, " and ｛｛ｎｏ｝｝") +
		tp(`\{{literal}} {{missing}}{{flag}} {{n}}`) +
		tp("{{note}}") +
		tp(`<w:hyperlink w:anchor="top"><w:r><w:t>see {{no}}</w:t></w:r></w:hyperlink>`, `<w:r><w:tab/><w:t>{{n}}</w:t><w:br/><w:t>{{n}}</w:t></w:r>`)
	data := map[string]any{"no": "A-1", "customer": map[string]any{"name": "ACME"}, "total": 1234.5, "ｎｏ": "wide", "flag": true, "n": 3.0, "note": "one\ntwo"}
	res, r := fillDocx(t, body, nil, data)
	wantBody(t, r, "Invoice A-1 for ACME! / Total [1234.5 as #,##0] and wide / {{literal}} TRUE 3 / one / two / see A-1 3 / 3")
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], `paragraph 3: the data has no "missing"`) {
		t.Errorf("warnings: %v", res.Warnings)
	}
	if len(r.Manifest.Views) != 1 || r.Manifest.Views[0].Kind != bdf.ViewFlow {
		t.Errorf("views = %d (%s), want the pages alone", len(r.Manifest.Views), r.Manifest.Views[0].Kind)
	}
	if err := fillDocxErr(t, body, map[string]any{}, true); !strings.Contains(err.Error(), `the data has no`) {
		t.Errorf("strict: %v", err)
	}
	// a document without tags is filled as it is
	res, r = fillDocx(t, tp("Nothing to fill"), nil, nil)
	noWarnings(t, res)
	wantBody(t, r, "Nothing to fill")
}

func TestFillBlocks(t *testing.T) {
	body := tp("Items:{{#items}} {{name}}{{^_last}},{{/}}{{/items}}.") +
		tp("{{#items}}") +
		tp("{{_index}}/{{_count}} {{name}}") +
		tp("{{#tags}}", "- {{.}}") +
		tp("{{/tags}}") +
		ttbl([]string{"in {{name}}"}) +
		tp("{{/items}}") +
		tp("{{^items}}") + tp("nothing") + tp("{{/items}}") +
		tp("end")
	res, r := fillDocx(t, body, nil, map[string]any{"items": []any{
		map[string]any{"name": "a", "tags": []any{"x", "y"}},
		map[string]any{"name": "b", "tags": []any{}},
	}})
	noWarnings(t, res)
	wantBody(t, r, "Items: a, b. / 1/2 a / - x / - y / in a / 2/2 b / in b / end")
	res, r = fillDocx(t, body, nil, map[string]any{"items": []any{}})
	noWarnings(t, res)
	wantBody(t, r, "Items:. / nothing / end")

	// a section that says that it repeats: the paragraph, not the text of it
	body = tp("{{#items | rows}}- {{name}}{{#tags}} {{.}}{{/tags}}{{/items}}") + tp("{{^items | rows}}none{{/items}}") + tp("end")
	_, r = fillDocx(t, body, nil, map[string]any{"items": []any{
		map[string]any{"name": "a", "tags": []any{"x", "y"}}, map[string]any{"name": "b"}}})
	wantBody(t, r, "- a x y / - b / end")
	_, r = fillDocx(t, body, nil, map[string]any{"items": []any{}})
	wantBody(t, r, "none / end")
}

func TestFillTableRows(t *testing.T) {
	items := func(n int) []any {
		var out []any
		for i := 1; i <= n; i++ {
			out = append(out, map[string]any{"name": "item" + strconv.Itoa(i), "qty": float64(i)})
		}
		return out
	}
	// the tags in a row: the row is repeated
	body := ttbl(
		[]string{"Name", "Qty"},
		[]string{"{{#items}}{{name}}", "{{qty}}{{/items}}"},
		[]string{"Total", "{{total}}"})
	res, r := fillDocx(t, body, nil, map[string]any{"items": items(3), "total": 6.0})
	noWarnings(t, res)
	wantBody(t, r, "Name / Qty / item1 / 1 / item2 / 2 / item3 / 3 / Total / 6")
	_, r = fillDocx(t, body, nil, map[string]any{"items": items(0), "total": 0.0})
	wantBody(t, r, "Name / Qty / Total / 0")

	// the tags in different rows: the rows are slots
	body = ttbl(
		[]string{"Name", "Qty"},
		[]string{"{{#items}}{{name}}", "{{qty}}"},
		[]string{"", ""},
		[]string{"", "{{/items}}"},
		[]string{"Total", "{{total}}"})
	res, r = fillDocx(t, body, nil, map[string]any{"items": items(2), "total": 3.0})
	noWarnings(t, res)
	wantBody(t, r, "Name / Qty / item1 / 1 / item2 / 2 / Total / 3")
	_, r = fillDocx(t, body, nil, map[string]any{"items": items(3), "total": 6.0})
	wantBody(t, r, "Name / Qty / item1 / 1 / item2 / 2 / item3 / 3 / Total / 6")
	if err := fillDocxErr(t, body, map[string]any{"items": items(4)}, false); !strings.Contains(err.Error(), "table 1, row 2, cell 1: 4 items for the 3 rows of {{#items}}") {
		t.Errorf("slots: %v", err)
	}
	file := docxBytes(t, body, nil)
	fields, err := Inspect(bytes.NewReader(file), int64(len(file)))
	if err != nil {
		t.Fatal(err)
	}
	var sec *tag.Field
	for i := range fields {
		if fields[i].Kind == tag.FieldSection {
			sec = &fields[i]
		}
	}
	if sec == nil || sec.Name != "items" || sec.Mode != tag.ModeSlots || sec.Slots != 3 {
		t.Errorf("fields = %+v", fields)
	}

	// "rows" repeats the rows between the tags, and the row of a cell
	body = ttbl(
		[]string{"{{#items | rows}}{{name}}", ""},
		[]string{"", "{{qty}}{{/items}}"}) +
		ttbl([]string{"{{#items | rows}}{{name}}{{/items}}"}) +
		ttbl([]string{"{{#items | rows}}\n{{name}}\n{{/items}}", "x"})
	_, r = fillDocx(t, body, nil, map[string]any{"items": items(2)})
	wantBody(t, r, "item1 / 1 / item2 / 2 / item1 / item2 / item1 / x / item2 / x")

	// paragraphs of a cell, and a cell that keeps a paragraph
	body = ttbl([]string{"{{#items}}\n{{name}}\n{{/items}}", "{{#none}}\n{{/none}}"}) + tp("end")
	_, r = fillDocx(t, body, nil, map[string]any{"items": items(2), "none": []any{}})
	wantBody(t, r, "item1 / item2 / end")

	// a row of slots with tags of its own, and the marks of empty slots
	body = ttbl(
		[]string{"{{#items}}{{name}}", "{{^_blank}}x{{/}}"},
		[]string{"second {{name}}{{#_blank_first}}(none below){{/}}", "{{/items}}"})
	_, r = fillDocx(t, body, nil, map[string]any{"items": items(1)})
	wantBody(t, r, "item1 / x / second (none below)")
}

func TestFillPages(t *testing.T) {
	people := []any{map[string]any{"name": "a"}, map[string]any{"name": "b"}, map[string]any{"name": "c"}}
	body := tp("{{#people | pages}}") + tp("Certificate {{_page}}/{{_pages}}") + tp("for {{name}}") + tp("{{/people}}") + smallPage
	res, r := fillDocx(t, body, nil, map[string]any{"people": people})
	noWarnings(t, res)
	wantBody(t, r, "Certificate 1/3 / for a", "Certificate 2/3 / for b", "Certificate 3/3 / for c")
	// a table at the top of what is repeated
	body = ttbl([]string{"{{#people | pages}}{{name}}"}) + tp("end{{/people}}") + smallPage
	_, r = fillDocx(t, body, nil, map[string]any{"people": people[:2]})
	wantBody(t, r, "a / end", "b / end")
	// one paragraph for each page
	body = tp("{{#people | pages}}{{name}}{{/people}}") + smallPage
	_, r = fillDocx(t, body, nil, map[string]any{"people": people})
	wantBody(t, r, "a", "b", "c")
}

func TestFillHeaderAndTextBox(t *testing.T) {
	const wNS = `xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"`
	box := `<w:r><w:drawing><wp:inline xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing"><wp:extent cx="1270000" cy="635000"/><wp:docPr id="1" name="Box"/>` +
		`<a:graphic xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><a:graphicData uri="http://schemas.microsoft.com/office/word/2010/wordprocessingShape">` +
		`<wps:wsp xmlns:wps="http://schemas.microsoft.com/office/word/2010/wordprocessingShape"><wps:cNvSpPr txBox="1"/><wps:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="1270000" cy="635000"/></a:xfrm>` +
		`<a:prstGeom prst="rect"><a:avLst/></a:prstGeom></wps:spPr><wps:txbx><w:txbxContent>` + tp("{{#lines}}") + tp("box {{.}} of {{name}}") + tp("{{/lines}}") +
		`</w:txbxContent></wps:txbx><wps:bodyPr/></wps:wsp></a:graphicData></a:graphic></wp:inline></w:drawing></w:r>`
	body := tp("{{#people}}") + tp("for {{name}}", box) + tp("{{/people}}") +
		`<w:sectPr><w:headerReference w:type="default" r:id="rIdheader"/><w:footerReference w:type="default" r:id="rIdfooter"/>` +
		`<w:pgSz w:w="6000" w:h="8000"/><w:pgMar w:top="1000" w:bottom="1000" w:left="1000" w:right="1000" w:header="500" w:footer="500"/></w:sectPr>`
	parts := map[string]string{
		"word/header.xml": `<w:hdr ` + wNS + `>` + tp("{{company}} header") + `</w:hdr>`,
		"word/footer.xml": `<w:ftr ` + wNS + `>` + tp("{{#show}}footer of {{company}}{{/show}}") + `</w:ftr>`,
	}
	res, r := fillDocx(t, body, parts, map[string]any{"company": "ACME", "show": true, "lines": []any{"1", "2"},
		"people": []any{map[string]any{"name": "a"}, map[string]any{"name": "b"}}})
	noWarnings(t, res)
	if got := pageText(t, r, 0, 0, bdf.RoleHeader); got != "ACME header" {
		t.Errorf("header = %q", got)
	}
	if got := pageText(t, r, 0, 0, bdf.RoleFooter); got != "footer of ACME" {
		t.Errorf("footer = %q", got)
	}
	got := strings.ReplaceAll(pageText(t, r, 0, 0, bdf.RoleBody), "\n", " / ")
	for _, want := range []string{"for a", "for b", "box 1 of a / box 2 of a", "box 1 of b / box 2 of b"} {
		if !strings.Contains(got, want) {
			t.Errorf("body = %q, want %q in it", got, want)
		}
	}
	var names []string
	for _, f := range res.Fields {
		names = append(names, fmt.Sprintf("%s %s in %v at %s", f.Kind, f.Name, f.In, strings.Join(f.At, ";")))
	}
	want := []string{"section people in [] at paragraph 1", "value name in [people] at paragraph 2",
		"section lines in [people] at text box 1, paragraph 1", "value . in [people lines] at text box 1, paragraph 2", "value name in [people lines] at text box 1, paragraph 2",
		"section show in [] at footer, paragraph 1", "value company in [show] at footer, paragraph 1", "value company in [] at header, paragraph 1"}
	if strings.Join(names, "\n") != strings.Join(want, "\n") {
		t.Errorf("fields:\n%s\nwant:\n%s", strings.Join(names, "\n"), strings.Join(want, "\n"))
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

// tpic makes a run with a picture in the line: 100 × 100 pt, with an
// alternative text.
func tpic(descr string) string {
	return `<w:r><w:drawing><wp:inline xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing"><wp:extent cx="1270000" cy="1270000"/>` +
		`<wp:docPr id="1" name="Picture" descr="` + html.EscapeString(descr) + `"/><a:graphic xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
		`<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture"><pic:pic xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">` +
		`<pic:nvPicPr><pic:cNvPr id="0" name="p"/><pic:cNvPicPr/></pic:nvPicPr><pic:blipFill><a:blip r:embed="rIdmedia/image1.png"/><a:stretch><a:fillRect/></a:stretch></pic:blipFill>` +
		`<pic:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="1270000" cy="1270000"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></pic:spPr></pic:pic></a:graphicData></a:graphic></wp:inline></w:drawing></w:r>`
}

func TestFillPicture(t *testing.T) {
	photo := "data:image/png;base64," + base64.StdEncoding.EncodeToString(tplPNG(40, 10))
	parts := map[string]string{"word/media/image1.png": string(tplPNG(4, 4))}
	body := tp(tpic("{{%photo}} the photo")) + tp(tpic("{{%photo | fit:cover}}")) + tp(tpic("{{%photo | fit:stretch}}")) + tp(tpic("a picture"))
	images := func(r *bdf.Reader) []string {
		var got []string
		for _, l := range r.Manifest.Views[0].Pages[0].Layers {
			o, err := r.Object(l.Obj)
			if err != nil {
				t.Fatal(err)
			}
			o.Walk(func(in bdf.Instr) {
				switch in.Op {
				case bdf.OpImage:
					got = append(got, fmt.Sprintf("image %v %v", in.Args[3], in.Args[4]))
				case bdf.OpImageSub:
					got = append(got, fmt.Sprintf("sub %v %v %v %v > %v %v", in.Args[1], in.Args[2], in.Args[3], in.Args[4], in.Args[7], in.Args[8]))
				case bdf.OpMark:
					if byte(in.Args[0].(uint64)) == bdf.MarkFigure {
						got = append(got, "figure "+in.Args[1].(string))
					}
				}
			})
		}
		return got
	}
	res, r := fillDocx(t, body, parts, map[string]any{"photo": photo})
	noWarnings(t, res)
	want := []string{"figure the photo", "image 100 25", "sub 15 0 10 10 > 100 100", "image 100 100", "figure a picture", "image 100 100"}
	if got := images(r); strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("got\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	// without a picture the picture is left out
	res, r = fillDocx(t, body, parts, map[string]any{"photo": ""})
	noWarnings(t, res)
	if got := images(r); strings.Join(got, "\n") != "figure a picture\nimage 100 100" {
		t.Errorf("got %v", got)
	}
	res, _ = fillDocx(t, body, parts, map[string]any{"photo": "missing.png"})
	if len(res.Warnings) == 0 || !strings.Contains(res.Warnings[0], `paragraph 1: the picture of "photo" is not read`) {
		t.Errorf("warnings: %v", res.Warnings)
	}
}

func TestFillErrors(t *testing.T) {
	for _, tc := range []struct{ name, body, want string }{
		{"unclosed", tp("{{#items}}") + tp("x"), "paragraph 1: the section {{#items}} is not closed"},
		{"stray", tp("x{{/items}}"), "{{/items}} closes no section"},
		{"mismatch", tp("{{#a}}") + tp("{{/b}}"), "{{/b}} closes the section {{#a}} of paragraph 1"},
		{"inline mismatch", tp("{{#a}}x{{/b}}"), "{{/b}} closes the section {{#a}}"},
		{"syntax", tp("{{#}}"), "docx: paragraph 1:"},
		{"picture in text", tp("{{%photo}}"), "the picture tag {{%photo}} is in a text"},
		{"picture alt", tp(tpic("{{photo}}")), "takes one picture tag"},
		{"same unit", tp("{{#a}}") + tp("x{{/a}}{{#b}}") + tp("{{/b}}"), "both repeat the same paragraph or row"},
		{"slots of paragraphs", tp("{{#a | slots}}") + tp("{{/a}}"), `takes no "slots"`},
		{"slots of a row", ttbl([]string{"{{#a | slots}}", "{{/a}}"}), `takes no "slots"`},
		{"pages in a table", ttbl([]string{"{{#a | pages}}"}, []string{"{{/a}}"}), `is in a table: it takes no "pages"`},
		{"rows in slots", ttbl([]string{"{{#a}}", ""}, []string{"{{#b}}x", "{{/b}}"}, []string{"", "{{/a}}"}), "which are slots"},
		{"out of a hyperlink", tp(`<w:hyperlink w:anchor="top"><w:r><w:t>{{#a}}x</w:t></w:r></w:hyperlink>`, "{{/a}}"), "both in or both outside a hyperlink"},
		{"across text boxes", tp("{{#a}}") + ttbl([]string{"{{/b}}"}), "closes the section {{#a}}"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := fillDocxErr(t, tc.body, map[string]any{}, false)
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("%v, want %q", err, tc.want)
			}
			file := docxBytes(t, tc.body, nil)
			if _, ierr := Inspect(bytes.NewReader(file), int64(len(file))); ierr == nil || ierr.Error() != err.Error() {
				t.Errorf("Inspect: %v, Fill: %v", ierr, err)
			}
		})
	}
}

func TestInspectFields(t *testing.T) {
	body := tp("{{title}} {{_page}}", tpic("{{%logo}}")) +
		tp("{{#groups | pages}}") + tp("{{name}}") +
		ttbl([]string{"{{#rows}}{{name}}", `{{amount | fmt:"#,##0"}}{{/rows}}`}) +
		tp("{{title}}{{/groups}}")
	file := docxBytes(t, body, nil)
	fields, err := Inspect(bytes.NewReader(file), int64(len(file)))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range fields {
		got = append(got, fmt.Sprintf("%s %s in %v mode=%s fmt=%s at %v", f.Kind, f.Name, f.In, f.Mode, f.Format, f.At))
	}
	want := []string{
		"value title in [] mode= fmt= at [paragraph 1]",
		"image logo in [] mode= fmt= at [paragraph 1]",
		"section groups in [] mode=pages fmt= at [paragraph 2]",
		"value name in [groups] mode= fmt= at [paragraph 3]",
		"section rows in [groups] mode=rows fmt= at [table 1, row 1, cell 1]",
		"value name in [groups rows] mode= fmt= at [table 1, row 1, cell 1]",
		"value amount in [groups rows] mode= fmt=#,##0 at [table 1, row 1, cell 2]",
		"value title in [groups] mode= fmt= at [paragraph 4]",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("got\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

// codeBounds recognizes bars/modules followed optionally by closed digit
// outlines. Digits must use curves and stay below the bars; the returned box
// includes both, so frame bounds are checked by the caller's expected symbol.
func codeBounds(p *bdf.Path) (rects, contours int, box [4]float32) {
	if p == nil || len(p.Verbs) < 8 || p.Verbs[0] != bdf.VerbRect {
		return
	}
	box = [4]float32{p.Args[0], p.Args[1], p.Args[0], p.Args[1]}
	grow := func(x, y float32) {
		box[0], box[1] = min(box[0], x), min(box[1], y)
		box[2], box[3] = max(box[2], x), max(box[3], y)
	}
	at := 0
	for rects < len(p.Verbs) && p.Verbs[rects] == bdf.VerbRect {
		grow(p.Args[at], p.Args[at+1])
		grow(p.Args[at]+p.Args[at+2], p.Args[at+1]+p.Args[at+3])
		at += 4
		rects++
	}
	if rects < 8 {
		return 0, 0, [4]float32{}
	}
	bottom, curves, closed := box[3], 0, 0
	for _, verb := range p.Verbs[rects:] {
		points := 0
		switch verb {
		case bdf.VerbMove:
			contours++
			points = 1
		case bdf.VerbLine:
			points = 1
		case bdf.VerbQuad:
			curves++
			points = 2
		case bdf.VerbClose:
			closed++
		default:
			return 0, 0, [4]float32{}
		}
		for ; points > 0; points-- {
			x, y := p.Args[at], p.Args[at+1]
			if y <= bottom {
				return 0, 0, [4]float32{}
			}
			grow(x, y)
			at += 2
		}
	}
	if contours > 0 && (curves == 0 || closed != contours) {
		return 0, 0, [4]float32{}
	}
	return
}

// rectBox summarizes code modules/bars, digit contours and their full bounds.
func rectBox(p *bdf.Path) string {
	rects, contours, box := codeBounds(p)
	if rects == 0 {
		return ""
	}
	label := ""
	if contours > 0 {
		label = fmt.Sprintf(" + %d contours below bars", contours)
	}
	return fmt.Sprintf("%d rects of %.2f × %.2f%s", rects, box[2]-box[0], box[3]-box[1], label)
}

// TestFillCodes fills picture tags that ask for a QR code or a barcode:
// the symbol of the value is drawn in the place of the picture, and in
// each row of a table that is filled.
func TestFillCodes(t *testing.T) {
	parts := map[string]string{"word/media/image1.png": string(tplPNG(4, 4))}
	body := tp(tpic("{{%url | qr}}")) + tp(tpic("{{%jan | barcode:jan}} the article")) +
		ttbl([]string{"{{#items}}{{name}}", tpic("{{%code | barcode:code128 | quiet:0}}") + tr("{{/items}}")}) + tp(tpic("a picture"))
	const url, jan = "https://example.com/invoice/1024", "4901234567894"
	data := map[string]any{"url": url, "jan": jan, "items": []any{
		map[string]any{"name": "A", "code": "A-001"}, map[string]any{"name": "B", "code": "B-002"}}}
	codes := func(r *bdf.Reader) []string {
		var got []string
		for _, l := range r.Manifest.Views[0].Pages[0].Layers {
			o, err := r.Object(l.Obj)
			if err != nil {
				t.Fatal(err)
			}
			o.Walk(func(in bdf.Instr) {
				switch in.Op {
				case bdf.OpImage:
					got = append(got, "image")
				case bdf.OpFillPath:
					if s := rectBox(o.Paths[int(in.Args[0].(uint64))].Inline); s != "" {
						got = append(got, s)
					}
				case bdf.OpMark:
					if byte(in.Args[0].(uint64)) == bdf.MarkFigure {
						got = append(got, "figure "+in.Args[1].(string))
					}
				}
			})
		}
		return got
	}
	sym := func(kind, text string, o *barcode.Options, w, h float32) string {
		s, err := barcode.Encode(kind, text, o)
		if err != nil {
			t.Fatal(err)
		}
		return rectBox(s.Path(0, 0, w, h))
	}
	res, r := fillDocx(t, body, parts, data)
	noWarnings(t, res)
	want := []string{
		"figure " + url, sym(barcode.QR, url, nil, 100, 100),
		"figure the article", sym(barcode.EAN13, jan, nil, 100, 100),
		"figure A-001", sym(barcode.Code128, "A-001", &barcode.Options{Quiet: barcode.NoQuiet}, 100, 100),
		"figure B-002", sym(barcode.Code128, "B-002", &barcode.Options{Quiet: barcode.NoQuiet}, 100, 100),
		"figure a picture", "image",
	}
	got := codes(r)
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("got\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if len(got) < 4 || !strings.HasSuffix(got[3], " contours below bars") {
		t.Errorf("JAN lacks human-readable vector digits: %q", got)
	}
	// the fields tell which code a picture is
	file := docxBytes(t, body, parts)
	fields, err := Inspect(bytes.NewReader(file), int64(len(file)))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, f := range fields {
		names = append(names, strings.TrimSpace(f.Kind+" "+f.Name+" "+f.Code))
	}
	if strings.Join(names, ", ") != "image url qr, image jan ean13, section items, value name, image code code128" {
		t.Errorf("fields = %q", names)
	}
	// a value that the code can not hold is left out with a warning, or is
	// an error; without a value there is no code
	data = map[string]any{"url": "", "jan": "49012345678", "items": []any{}}
	res, r = fillDocx(t, body, parts, data)
	if got := codes(r); strings.Join(got, "\n") != "figure a picture\nimage" {
		t.Errorf("got %v", got)
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], `paragraph 2: the code of "jan" is left out`) {
		t.Errorf("warnings: %v", res.Warnings)
	}
	if _, err := Fill(bytes.NewReader(file), int64(len(file)), &TemplateData{Data: data, Strict: true}, testOptions()); err == nil || !strings.Contains(err.Error(), "paragraph 2") {
		t.Errorf("strict: %v", err)
	}
}
