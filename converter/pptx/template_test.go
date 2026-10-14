package pptx

import (
	"archive/zip"
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
	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/form/tag"
	"github.com/shibukawa/bdf/image/barcode"
)

const tplNS = ` xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"` +
	` xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"` +
	` xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"`

// tplDeck makes a presentation (720 × 540 pt, without layouts) of the
// slides whose shapes are given. The pictures of the slides are
// ppt/media/image1.png, which rId1 of every slide names.
func tplDeck(t testing.TB, slides ...string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	add := func(name, data string) {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(data))
	}
	const relNS = `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`
	const relType = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/"
	add("[Content_Types].xml", `<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">`+
		`<Default Extension="xml" ContentType="application/xml"/><Default Extension="png" ContentType="image/png"/>`+
		`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/></Types>`)
	add("_rels/.rels", relNS+`<Relationship Id="rId1" Type="`+relType+`officeDocument" Target="ppt/presentation.xml"/></Relationships>`)
	var ids, rels strings.Builder
	for i, s := range slides {
		n := strconv.Itoa(i + 1)
		ids.WriteString(`<p:sldId id="` + strconv.Itoa(256+i) + `" r:id="rId` + n + `"/>`)
		rels.WriteString(`<Relationship Id="rId` + n + `" Type="` + relType + `slide" Target="slides/slide` + n + `.xml"/>`)
		add("ppt/slides/slide"+n+".xml", `<p:sld`+tplNS+`><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>`+
			s+`</p:spTree></p:cSld></p:sld>`)
		add("ppt/slides/_rels/slide"+n+".xml.rels", relNS+`<Relationship Id="rId1" Type="`+relType+`image" Target="../media/image1.png"/></Relationships>`)
	}
	add("ppt/presentation.xml", `<p:presentation`+tplNS+`><p:sldIdLst>`+ids.String()+`</p:sldIdLst><p:sldSz cx="9144000" cy="6858000"/></p:presentation>`)
	add("ppt/_rels/presentation.xml.rels", relNS+rels.String()+`</Relationships>`)
	add("ppt/media/image1.png", string(tplPNG(4, 4)))
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
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

func emu(pt float64) string { return strconv.FormatInt(int64(pt*12700), 10) }

func tplXfrm(x, y, w, h float64) string {
	return `<a:off x="` + emu(x) + `" y="` + emu(y) + `"/><a:ext cx="` + emu(w) + `" cy="` + emu(h) + `"/>`
}

// tp makes a paragraph of runs; a run that starts with "<" is markup.
func tp(runs ...string) string {
	var b strings.Builder
	b.WriteString("<a:p>")
	for _, r := range runs {
		if strings.HasPrefix(r, "<") {
			b.WriteString(r)
			continue
		}
		b.WriteString(`<a:r><a:rPr lang="en-US" sz="1800"/><a:t>` + html.EscapeString(r) + `</a:t></a:r>`)
	}
	b.WriteString("</a:p>")
	return b.String()
}

// tsp makes a text box; fit is what its a:bodyPr holds (its autofit).
func tsp(name string, x, y, w, h float64, fit string, paras ...string) string {
	return `<p:sp><p:nvSpPr><p:cNvPr id="2" name="` + name + `"/><p:cNvSpPr txBox="1"/><p:nvPr/></p:nvSpPr><p:spPr><a:xfrm>` + tplXfrm(x, y, w, h) +
		`</a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></p:spPr><p:txBody><a:bodyPr wrap="square" lIns="0" tIns="0" rIns="0" bIns="0">` + fit +
		`</a:bodyPr><a:lstStyle/>` + strings.Join(paras, "") + `</p:txBody></p:sp>`
}

// ttbl makes a table of rows of cells (each the text of a paragraph), with
// rows of a height and columns of the same width.
func ttbl(name string, x, y, w, rowH float64, rows ...[]string) string {
	var b strings.Builder
	b.WriteString(`<p:graphicFrame><p:nvGraphicFramePr><p:cNvPr id="3" name="` + name + `"/><p:cNvGraphicFramePr/><p:nvPr/></p:nvGraphicFramePr><p:xfrm>` +
		tplXfrm(x, y, w, rowH*float64(len(rows))) + `</p:xfrm><a:graphic><a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/table"><a:tbl><a:tblPr/><a:tblGrid>`)
	for range rows[0] {
		b.WriteString(`<a:gridCol w="` + emu(w/float64(len(rows[0]))) + `"/>`)
	}
	b.WriteString(`</a:tblGrid>`)
	for _, row := range rows {
		b.WriteString(`<a:tr h="` + emu(rowH) + `">`)
		for _, cell := range row {
			b.WriteString(`<a:tc><a:txBody><a:bodyPr/><a:lstStyle/>`)
			for _, line := range strings.Split(cell, "\n") {
				b.WriteString(`<a:p><a:r><a:rPr lang="en-US" sz="1000"/><a:t>` + html.EscapeString(line) + `</a:t></a:r></a:p>`)
			}
			b.WriteString(`</a:txBody><a:tcPr marL="0" marR="0" marT="0" marB="0"/></a:tc>`)
		}
		b.WriteString(`</a:tr>`)
	}
	b.WriteString(`</a:tbl></a:graphicData></a:graphic></p:graphicFrame>`)
	return b.String()
}

// tpic makes a picture with an alternative text.
func tpic(name, descr string, x, y, w, h float64) string {
	return `<p:pic><p:nvPicPr><p:cNvPr id="4" name="` + name + `" descr="` + html.EscapeString(descr) + `"/><p:cNvPicPr/><p:nvPr/></p:nvPicPr>` +
		`<p:blipFill><a:blip r:embed="rId1"/><a:stretch><a:fillRect/></a:stretch></p:blipFill><p:spPr><a:xfrm>` + tplXfrm(x, y, w, h) +
		`</a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></p:spPr></p:pic>`
}

// tplFormat stands for the formatter of number formats.
func tplFormat(v any, format string) string {
	if format == "" {
		return fmt.Sprint(v)
	}
	return fmt.Sprintf("[%v as %s]", v, format)
}

func fillDeck(t *testing.T, deck []byte, data map[string]any, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
	if opts == nil {
		opts = testOptions()
	}
	res, err := Fill(bytes.NewReader(deck), int64(len(deck)), &TemplateData{Data: data, Format: tplFormat}, opts)
	if err != nil {
		t.Fatal(err)
	}
	return res, reopen(t, res)
}

func fillErr(t *testing.T, deck []byte, data map[string]any, strict bool) error {
	t.Helper()
	_, err := Fill(bytes.NewReader(deck), int64(len(deck)), &TemplateData{Data: data, Strict: strict}, testOptions())
	if err == nil {
		t.Fatal("no error")
	}
	return err
}

// pageTexts returns the texts of the pages: the lines of each joined by
// " / " (a paragraph may be more than one of them).
func pageTexts(t *testing.T, r *bdf.Reader) []string {
	t.Helper()
	var out []string
	for i := range r.Manifest.Views[0].Pages {
		p := r.Manifest.Views[0].Pages[i]
		var lines []string
		var walk func(h bdf.Hash)
		walk = func(h bdf.Hash) {
			o, err := r.Object(h)
			if err != nil {
				t.Fatal(err)
			}
			o.Walk(func(in bdf.Instr) {
				switch in.Op {
				case bdf.OpFillText:
					lines = append(lines, in.Args[0].(string))
				case bdf.OpUse, bdf.OpUseAt:
					walk(o.Objects[int(in.Args[0].(uint64))])
				}
			})
		}
		if h := body(p); h != (bdf.Hash{}) {
			walk(h)
		}
		out = append(out, strings.Join(lines, " / "))
	}
	return out
}

func wantPages(t *testing.T, r *bdf.Reader, want ...string) {
	t.Helper()
	got := pageTexts(t, r)
	if len(got) != len(want) {
		t.Fatalf("%d pages, want %d:\n%s", len(got), len(want), strings.Join(got, "\n"))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("page %d:\n got %s\nwant %s", i+1, got[i], want[i])
		}
	}
}

func noWarnings(t *testing.T, res *Result) {
	t.Helper()
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
}

func TestFillValues(t *testing.T) {
	deck := tplDeck(t,
		tsp("Title", 36, 36, 600, 40, "", tp("Invoice {{no}} for ", "{{custo", "mer.na", "me}}", "!"))+
			tsp("Body", 36, 100, 600, 200, "",
				tp(`Total {{total | fmt:"#,##0"}}`, " and ｛｛ｎｏ｝｝"),
				tp(`\{{literal}} {{missing}}{{flag}} {{n}}`),
				tp("{{note}}")))
	res, r := fillDeck(t, deck, map[string]any{
		"no": "A-1", "customer": map[string]any{"name": "ACME"}, "total": 1234.5, "ｎｏ": "wide", "flag": true, "n": 3.0,
		"note": "one\ntwo",
	}, nil)
	wantPages(t, r, "Invoice A-1 for ACME! / Total [1234.5 as #,##0] and wide / {{literal}} TRUE 3 / one / two")
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], `slide 1: Body: the data has no "missing"`) {
		t.Errorf("warnings: %v", res.Warnings)
	}
	if err := fillErr(t, deck, map[string]any{}, true); !strings.Contains(err.Error(), `the data has no`) {
		t.Errorf("strict: %v", err)
	}
	// a presentation without tags is filled as it is
	plain := tplDeck(t, tsp("Title", 36, 36, 600, 40, "", tp("Nothing to fill")))
	res, r = fillDeck(t, plain, nil, nil)
	noWarnings(t, res)
	wantPages(t, r, "Nothing to fill")
}

func TestFillParagraphs(t *testing.T) {
	deck := tplDeck(t, tsp("List", 36, 36, 600, 300, "",
		tp("Items:{{#items}} {{name}}{{^_last}},{{/}}{{/items}}."),
		tp("{{#items}}"),
		tp("{{_index}}/{{_count}} {{name}}"),
		tp("{{#tags}}", "- {{.}}"),
		tp("{{/tags}}"),
		tp("{{/items}}"),
		tp("{{^items}}"), tp("nothing"), tp("{{/items}}"),
		tp("end")))
	res, r := fillDeck(t, deck, map[string]any{"items": []any{
		map[string]any{"name": "a", "tags": []any{"x", "y"}},
		map[string]any{"name": "b", "tags": []any{}},
	}}, nil)
	noWarnings(t, res)
	wantPages(t, r, "Items: a, b. / 1/2 a / - x / - y / 2/2 b / end")
	res, r = fillDeck(t, deck, map[string]any{"items": []any{}}, nil)
	noWarnings(t, res)
	wantPages(t, r, "Items:. / nothing / end")

	// a section that says that it repeats: the paragraph, not the text of it
	deck = tplDeck(t, tsp("List", 36, 36, 600, 300, "", tp("{{#items | rows}}- {{name}}{{#tags}} {{.}}{{/tags}}{{/items}}"), tp("end")))
	_, r = fillDeck(t, deck, map[string]any{"items": []any{
		map[string]any{"name": "a", "tags": []any{"x", "y"}}, map[string]any{"name": "b"}}}, nil)
	wantPages(t, r, "- a x y / - b / end")

	// a text whose paragraphs are all left out keeps an empty one
	deck = tplDeck(t, tsp("List", 36, 36, 600, 300, "", tp("{{#items}}"), tp("{{name}}"), tp("{{/items}}"))+tsp("Other", 36, 400, 600, 40, "", tp("x")))
	_, r = fillDeck(t, deck, map[string]any{"items": []any{}}, nil)
	wantPages(t, r, "x")
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
	deck := tplDeck(t, ttbl("Table", 36, 36, 400, 20,
		[]string{"Name", "Qty"},
		[]string{"{{#items}}{{name}}", "{{qty}}{{/items}}"},
		[]string{"Total", "{{total}}"}))
	res, r := fillDeck(t, deck, map[string]any{"items": items(3), "total": 6.0}, nil)
	noWarnings(t, res)
	wantPages(t, r, "Name / Qty / item1 / 1 / item2 / 2 / item3 / 3 / Total / 6")
	_, r = fillDeck(t, deck, map[string]any{"items": items(0), "total": 0.0}, nil)
	wantPages(t, r, "Name / Qty / Total / 0")

	// the tags in different rows: the rows are slots
	deck = tplDeck(t, ttbl("Table", 36, 36, 400, 20,
		[]string{"Name", "Qty"},
		[]string{"{{#items}}{{name}}", "{{qty}}"},
		[]string{"", ""},
		[]string{"", "{{/items}}"},
		[]string{"Total", "{{total}}"}))
	res, r = fillDeck(t, deck, map[string]any{"items": items(2), "total": 3.0}, nil)
	noWarnings(t, res)
	wantPages(t, r, "Name / Qty / item1 / 1 / item2 / 2 / Total / 3")
	_, r = fillDeck(t, deck, map[string]any{"items": items(3), "total": 6.0}, nil)
	wantPages(t, r, "Name / Qty / item1 / 1 / item2 / 2 / item3 / 3 / Total / 6")
	if err := fillErr(t, deck, map[string]any{"items": items(4)}, false); !strings.Contains(err.Error(), "4 items for the 3 rows of {{#items}}") {
		t.Errorf("slots: %v", err)
	}
	fields, err := Inspect(bytes.NewReader(deck), int64(len(deck)), nil)
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

	// "rows" repeats the rows between the tags
	deck = tplDeck(t, ttbl("Table", 36, 36, 400, 20,
		[]string{"{{#items | rows}}{{name}}", ""},
		[]string{"", "{{qty}}{{/items}}"}))
	_, r = fillDeck(t, deck, map[string]any{"items": items(2)}, nil)
	wantPages(t, r, "item1 / 1 / item2 / 2")

	// a row of slots with tags of its own, and the marks of empty slots
	deck = tplDeck(t, ttbl("Table", 36, 36, 400, 20,
		[]string{"{{#items}}{{name}}", "{{^_blank}}x{{/}}"},
		[]string{"second {{name}}{{#_blank_first}}(none below){{/}}", "{{/items}}"}))
	_, r = fillDeck(t, deck, map[string]any{"items": items(1)}, nil)
	wantPages(t, r, "item1 / x / second (none below)")
}

func TestFillSlides(t *testing.T) {
	// tags in different shapes of a slide: the slide is repeated
	deck := tplDeck(t,
		tsp("Cover", 36, 36, 600, 40, "", tp("{{title}}")),
		tsp("Head", 36, 36, 600, 40, "", tp("{{#pages}}{{name}} {{_page}}/{{_pages}}"))+
			tsp("Foot", 36, 400, 600, 40, "", tp("{{carry}}{{/pages}}")),
		tsp("End", 36, 36, 600, 40, "", tp("end")))
	data := map[string]any{"title": "T", "pages": []any{
		map[string]any{"name": "a", "carry": 1.0}, map[string]any{"name": "b", "carry": 2.0}, map[string]any{"name": "c", "carry": 3.0}}}
	res, r := fillDeck(t, deck, data, nil)
	noWarnings(t, res)
	wantPages(t, r, "T", "a 1/3 / 1", "b 2/3 / 2", "c 3/3 / 3", "end")
	if res.Slides != 5 {
		t.Errorf("slides = %d", res.Slides)
	}
	// tags on different slides: the slides from the one to the other
	deck = tplDeck(t,
		tsp("A", 36, 36, 600, 40, "", tp("{{#groups}}A {{name}}")),
		tsp("B", 36, 36, 600, 40, "", tp("{{#members}}"), tp("{{.}}"), tp("{{/members}}")),
		tsp("C", 36, 36, 600, 40, "", tp("C {{name}}{{/groups}}")),
		tsp("None", 36, 36, 600, 40, "", tp("{{^groups | pages}}none{{/groups}}")))
	_, r = fillDeck(t, deck, map[string]any{"groups": []any{
		map[string]any{"name": "g1", "members": []any{"m1", "m2"}}, map[string]any{"name": "g2", "members": []any{}}}}, nil)
	wantPages(t, r, "A g1", "m1 / m2", "C g1", "A g2", "", "C g2")
	_, r = fillDeck(t, deck, map[string]any{"groups": []any{}}, nil)
	wantPages(t, r, "none")
	// the slides of a selection are those of the template
	_, r = fillDeck(t, deck, map[string]any{"groups": []any{}}, &Options{FontDirs: []string{"testdata/fonts"}, NoSystemFonts: true, Slides: conv.PageList(4)})
	wantPages(t, r, "none")
}

func TestFillPicture(t *testing.T) {
	photo := "data:image/png;base64," + base64.StdEncoding.EncodeToString(tplPNG(40, 10))
	deck := tplDeck(t,
		tpic("Contain", "{{%photo}} the photo", 100, 100, 100, 100)+
			tpic("Cover", "{{%photo | fit:cover}}", 300, 100, 100, 100)+
			tpic("Stretch", "{{%photo | fit:stretch}}", 500, 100, 100, 100)+
			tpic("Kept", "a picture", 100, 300, 100, 100))
	res, r := fillDeck(t, deck, map[string]any{"photo": photo}, nil)
	noWarnings(t, res)
	p := r.Manifest.Views[0].Pages[0]
	o, err := r.Object(body(p))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	o.Walk(func(in bdf.Instr) {
		switch in.Op {
		case bdf.OpTranslate:
			got = append(got, fmt.Sprintf("at %v %v", in.Args[0], in.Args[1]))
		case bdf.OpImage:
			got = append(got, fmt.Sprintf("image %v %v %v %v", in.Args[1], in.Args[2], in.Args[3], in.Args[4]))
		case bdf.OpImageSub:
			got = append(got, fmt.Sprintf("sub %v %v %v %v > %v %v %v %v", in.Args[1], in.Args[2], in.Args[3], in.Args[4], in.Args[5], in.Args[6], in.Args[7], in.Args[8]))
		case bdf.OpMark:
			if byte(in.Args[0].(uint64)) == bdf.MarkFigure {
				got = append(got, "figure "+in.Args[1].(string))
			}
		}
	})
	want := []string{
		"figure the photo", "at 100 137.5", "image 0 0 100 25", // the whole picture in the middle of the frame
		"at 300 100", "sub 15 0 10 10 > 0 0 100 100", // the middle of the picture
		"at 500 100", "image 0 0 100 100",
		"figure a picture", "at 100 300", "image 0 0 100 100",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("got\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	images := 0
	for _, e := range r.Manifest.Parts {
		if e.T == bdf.PartImage {
			images++
		}
	}
	if images != 2 {
		t.Errorf("%d images, want that of the data and that of the template", images)
	}
	// without a picture the picture is left out
	res, r = fillDeck(t, deck, map[string]any{"photo": ""}, nil)
	noWarnings(t, res)
	if c := opCounts(t, r, body(r.Manifest.Views[0].Pages[0])); c[bdf.OpImage] != 1 || c[bdf.OpImageSub] != 0 {
		t.Errorf("ops = %v", c)
	}
	res, _ = fillDeck(t, deck, map[string]any{"photo": "missing.png"}, nil)
	if len(res.Warnings) == 0 || !strings.Contains(res.Warnings[0], `slide 1: Contain: the picture of "photo" is not read`) {
		t.Errorf("warnings: %v", res.Warnings)
	}
	if err := fillErr(t, deck, map[string]any{"photo": "missing.png"}, true); !strings.Contains(err.Error(), "is not read") {
		t.Errorf("strict: %v", err)
	}
}

// fontSizes lists the sizes of the fonts that the body of a page sets.
func fontSizes(t *testing.T, r *bdf.Reader, page int) []float64 {
	t.Helper()
	var out []float64
	o, err := r.Object(body(r.Manifest.Views[0].Pages[page]))
	if err != nil {
		t.Fatal(err)
	}
	o.Walk(func(in bdf.Instr) {
		if in.Op == bdf.OpFont {
			out = append(out, float64(in.Args[1].(float32)))
		}
	})
	return out
}

func TestFillRefit(t *testing.T) {
	long := strings.Repeat("The quick brown fox jumps over the lazy dog. ", 12)
	// a text that shrinks on overflow: the scale of the template is that of
	// the text of the template
	deck := tplDeck(t,
		tsp("Shrink", 36, 36, 300, 100, `<a:normAutofit fontScale="90000"/>`, tp("{{text}}"))+
			tsp("Kept", 36, 300, 300, 100, `<a:normAutofit fontScale="50000"/>`, tp("short")))
	_, r := fillDeck(t, deck, map[string]any{"text": "short"}, nil)
	if got := fontSizes(t, r, 0); len(got) != 2 || got[0] != 18 || got[1] != 9 {
		t.Errorf("short: sizes = %v, want 18 (fitted again) and 9 (as the file says)", got)
	}
	_, r = fillDeck(t, deck, map[string]any{"text": long}, nil)
	if got := fontSizes(t, r, 0); len(got) != 2 || got[0] >= 18 || got[0] < 4 || got[1] != 9 {
		t.Errorf("long: sizes = %v, want a smaller size and 9", got)
	}
	// a shape that resizes to its text: the shape under it is not moved
	deck = tplDeck(t, `<p:sp><p:nvSpPr><p:cNvPr id="2" name="Grow"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr><p:spPr><a:xfrm>`+tplXfrm(36, 36, 300, 30)+
		`</a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom><a:solidFill><a:srgbClr val="FFCC00"/></a:solidFill></p:spPr>`+
		`<p:txBody><a:bodyPr wrap="square" lIns="0" tIns="0" rIns="0" bIns="0" anchor="t"><a:spAutoFit/></a:bodyPr><a:lstStyle/>`+tp("{{text}}")+`</p:txBody></p:sp>`)
	height := func(text string) float64 {
		_, r := fillDeck(t, deck, map[string]any{"text": text}, nil)
		o, err := r.Object(body(r.Manifest.Views[0].Pages[0]))
		if err != nil {
			t.Fatal(err)
		}
		return float64(o.BBox.H)
	}
	if a, b := height("short"), height(long); a > 31 || b < 100 {
		t.Errorf("heights = %v (one line) and %v (grown to the text)", a, b)
	}
}

func TestFillTableContinues(t *testing.T) {
	var items []any
	for i := 1; i <= 50; i++ {
		items = append(items, map[string]any{"name": "item" + strconv.Itoa(i)})
	}
	deck := tplDeck(t,
		tsp("Title", 36, 20, 600, 30, "", tp("List {{_slide}}/{{_slides}}"))+
			ttbl("Table", 36, 60, 400, 20,
				[]string{"Name"},
				[]string{"{{#items | rows}}{{name}}{{/items}}"},
				[]string{"Total {{total}}"}),
		tsp("After", 36, 36, 600, 40, "", tp("after")))
	res, r := fillDeck(t, deck, map[string]any{"items": items, "total": 50.0}, nil)
	pages := pageTexts(t, r)
	if len(pages) != 4 {
		t.Fatalf("%d pages:\n%s", len(pages), strings.Join(pages, "\n"))
	}
	// 504 pt less the top of the table at 60: 22 rows of 20 pt, one of them the head
	for i, p := range pages[:3] {
		if !strings.HasPrefix(p, "List "+strconv.Itoa(i+1)+"/4 / Name / item"+strconv.Itoa(1+i*21)+" / ") {
			t.Errorf("page %d: %s", i+1, p)
		}
	}
	if !strings.HasSuffix(pages[2], "item50 / Total 50") || strings.Contains(pages[0], "Total") || pages[3] != "after" {
		t.Errorf("pages:\n%s", strings.Join(pages, "\n"))
	}
	noWarnings(t, res)
	// a table that fits is left on its slide
	res, r = fillDeck(t, deck, map[string]any{"items": items[:5], "total": 5.0}, nil)
	noWarnings(t, res)
	if n := len(r.Manifest.Views[0].Pages); n != 2 {
		t.Errorf("%d pages", n)
	}
}

func TestFillErrors(t *testing.T) {
	box := func(paras ...string) string { return tsp("Box", 36, 36, 600, 100, "", paras...) }
	for _, tc := range []struct {
		name   string
		slides []string
		want   string
	}{
		{"unclosed", []string{box(tp("{{#items}}"), tp("x"))}, "slide 1: Box: the section {{#items}} is not closed"},
		{"stray", []string{box(tp("x{{/items}}"))}, "{{/items}} closes no section"},
		{"mismatch", []string{box(tp("{{#a}}"), tp("{{/b}}"))}, "{{/b}} closes the section {{#a}}"},
		{"inline mismatch", []string{box(tp("{{#a}}x{{/b}}"))}, "{{/b}} closes the section {{#a}}"},
		{"syntax", []string{box(tp("{{#}}"))}, "slide 1: Box:"},
		{"picture in text", []string{box(tp("{{%photo}}"))}, "the picture tag {{%photo}} is in a text"},
		{"picture alt", []string{tpic("Pic", "{{photo}}", 0, 0, 10, 10)}, "takes one picture tag"},
		{"same unit", []string{box(tp("{{#a}}"), tp("x{{/a}}{{#b}}"), tp("{{/b}}"))}, "both repeat paragraph 2"},
		{"slots of paragraphs", []string{box(tp("{{#a | slots}}"), tp("{{/a}}"))}, `takes no "slots"`},
		{"slots of a row", []string{ttbl("T", 0, 0, 100, 20, []string{"{{#a | slots}}", "{{/a}}"})}, `takes no "slots"`},
		{"rows of slides", []string{box(tp("{{#a | rows}}")) + box(tp("{{/a}}"))}, `takes no "rows"`},
		{"rows in slots", []string{ttbl("T", 0, 0, 100, 20, []string{"{{#a}}", ""}, []string{"{{#b}}x", "{{/b}}"}, []string{"", "{{/a}}"})}, "which are slots"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			deck := tplDeck(t, tc.slides...)
			err := fillErr(t, deck, map[string]any{}, false)
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("%v, want %q", err, tc.want)
			}
			if _, ierr := Inspect(bytes.NewReader(deck), int64(len(deck)), nil); ierr == nil || ierr.Error() != err.Error() {
				t.Errorf("Inspect: %v, Fill: %v", ierr, err)
			}
		})
	}
}

func TestInspectFields(t *testing.T) {
	deck := tplDeck(t,
		tsp("Title", 36, 36, 600, 40, "", tp("{{title}} {{_page}}"))+tpic("Logo", "{{%logo}}", 0, 0, 10, 10),
		tsp("Head", 36, 36, 600, 40, "", tp("{{#pages}}{{name}}"))+
			ttbl("Table", 36, 100, 400, 20, []string{"{{#rows}}{{name}}", `{{amount | fmt:"#,##0"}}{{/rows}}`})+
			tsp("Foot", 36, 400, 600, 40, "", tp("{{title}}{{/pages}}")))
	fields, err := Inspect(bytes.NewReader(deck), int64(len(deck)), nil)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range fields {
		got = append(got, fmt.Sprintf("%s %s in %v mode=%s fmt=%s at %v", f.Kind, f.Name, f.In, f.Mode, f.Format, f.At))
	}
	want := []string{
		"value title in [] mode= fmt= at [slide 1: Title]",
		"image logo in [] mode= fmt= at [slide 1: Logo]",
		"section pages in [] mode=pages fmt= at [slide 2: Head]",
		"value name in [pages] mode= fmt= at [slide 2: Head]",
		"section rows in [pages] mode=rows fmt= at [slide 2: Table]",
		"value name in [pages rows] mode= fmt= at [slide 2: Table]",
		"value amount in [pages rows] mode= fmt=#,##0 at [slide 2: Table]",
		"value title in [pages] mode= fmt= at [slide 2: Foot]",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("got\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestFillLimits(t *testing.T) {
	deck := tplDeck(t, tsp("A", 36, 36, 600, 40, "", tp("{{#pages}}x")), tsp("B", 36, 36, 600, 40, "", tp("{{/pages}}")))
	if err := fillErr(t, deck, map[string]any{"pages": make([]any, maxFilledSlides)}, false); !strings.Contains(err.Error(), "more than 10000 slides") {
		t.Errorf("slides: %v", err)
	}
	// a paragraph is cut
	deck = tplDeck(t, tsp("A", 36, 36, 600, 40, "", tp("{{#items}}0123456789{{/items}}")))
	res, _ := fillDeck(t, deck, map[string]any{"items": make([]any, maxParaText/10+5)}, nil)
	if len(res.Warnings) == 0 || !strings.Contains(strings.Join(res.Warnings, "\n"), "is cut at 65536 characters") {
		t.Errorf("warnings: %v", res.Warnings)
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
	return fmt.Sprintf("%d rects at %.2f %.2f of %.2f × %.2f%s", rects, box[0], box[1], box[2]-box[0], box[3]-box[1], label)
}

// TestFillCodes fills picture tags that ask for a QR code or a barcode:
// the symbol of the value is drawn in the frame of the picture.
func TestFillCodes(t *testing.T) {
	deck := tplDeck(t,
		tpic("QR", "{{%url | qr}}", 100, 100, 100, 100)+
			tpic("Wide", "{{%url | qr:H | quiet:0}} the invoice", 300, 100, 200, 100)+
			tpic("JAN", "{{%jan | barcode:jan}}", 100, 300, 200, 60)+
			tpic("Kept", "a picture", 400, 300, 100, 100))
	const url, jan = "https://example.com/invoice/1024", "4901234567894"
	data := map[string]any{"url": url, "jan": jan}
	res, r := fillDeck(t, deck, data, nil)
	noWarnings(t, res)
	o, err := r.Object(body(r.Manifest.Views[0].Pages[0]))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	o.Walk(func(in bdf.Instr) {
		switch in.Op {
		case bdf.OpTranslate:
			got = append(got, fmt.Sprintf("at %v %v", in.Args[0], in.Args[1]))
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
	sym := func(kind, text string, o *barcode.Options, x, y, w, h float32) string {
		s, err := barcode.Encode(kind, text, o)
		if err != nil {
			t.Fatal(err)
		}
		return rectBox(s.Path(x, y, w, h))
	}
	want := []string{
		"figure " + url, "at 100 100", sym(barcode.QR, url, nil, 0, 0, 100, 100), // what the code says is its text
		"figure the invoice", "at 350 100", sym(barcode.QR, url, &barcode.Options{Level: "H", Quiet: barcode.NoQuiet}, 0, 0, 100, 100), // a square in the middle
		"figure " + jan, "at 100 300", sym(barcode.EAN13, jan, nil, 0, 0, 200, 60), // bars and digits inside the frame
		"figure a picture", "at 400 300", "image",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("got\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if len(got) < 9 || !strings.HasSuffix(got[8], " contours below bars") {
		t.Errorf("JAN lacks human-readable vector digits: %q", got)
	}
	if !strings.Contains(want[5], " at 0.00 0.00 of 100.00 × 100.00") {
		t.Errorf("a code without a quiet zone fills its place: %s", want[5])
	}
	images := 0
	for _, e := range r.Manifest.Parts {
		if e.T == bdf.PartImage {
			images++
		}
	}
	if images != 1 {
		t.Errorf("%d images, want that of the template alone", images)
	}
	// the fields tell which code a picture is
	fields, err := Inspect(bytes.NewReader(deck), int64(len(deck)), nil)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, f := range fields {
		names = append(names, f.Kind+" "+f.Name+" "+f.Code)
	}
	if strings.Join(names, ", ") != "image url qr, image jan ean13" {
		t.Errorf("fields = %q", names)
	}
	// a value that the code can not hold is left out with a warning, or is
	// an error; without a value there is no code
	data["url"], data["jan"] = "", "49012345678"
	res, r = fillDeck(t, deck, data, nil)
	if c := opCounts(t, r, body(r.Manifest.Views[0].Pages[0])); c[bdf.OpFillPath] != 0 || c[bdf.OpImage] != 1 {
		t.Errorf("ops = %v", c)
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], `slide 1: JAN: the code of "jan" is left out`) {
		t.Errorf("warnings: %v", res.Warnings)
	}
	if err := fillErr(t, deck, data, true); !strings.Contains(err.Error(), "slide 1: JAN") {
		t.Errorf("strict: %v", err)
	}
}
