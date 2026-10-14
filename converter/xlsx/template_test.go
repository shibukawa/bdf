package xlsx

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shibukawa/bdf/image/barcode"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/form/tag"
)

// fillFile fills a template of testdata with JSON data and returns the
// result with the text of the pages of its views.
func fillFile(t *testing.T, name, data string, td *TemplateData) (*Result, map[string][]string) {
	t.Helper()
	if td == nil {
		td = &TemplateData{}
	}
	if err := json.Unmarshal([]byte(data), &td.Data); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	st, _ := f.Stat()
	res, err := Fill(f, st.Size(), td, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	return res, pageText(t, res.Doc)
}

func pageText(t *testing.T, doc *bdf.Document) map[string][]string {
	t.Helper()
	st, err := doc.SearchText()
	if err != nil {
		t.Fatal(err)
	}
	out := map[string][]string{}
	for _, v := range doc.Views {
		out[v.Title] = make([]string, len(v.Pages))
	}
	for _, v := range st.Views {
		for _, pg := range v.Pages {
			out[v.Title][pg.Page-1] = pg.Text
		}
	}
	return out
}

func testPNG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 40, 20))
	for y := 0; y < 20; y++ {
		for x := 0; x < 40; x++ {
			img.Set(x, y, color.RGBA{200, uint8(x * 6), uint8(y * 12), 255})
		}
	}
	var b bytes.Buffer
	png.Encode(&b, img)
	return b.Bytes()
}

// lines reports whether the text of a page has the lines one after
// another.
func lines(page string, want ...string) bool {
	return strings.Contains("\n"+page+"\n", "\n"+strings.Join(want, "\n")+"\n")
}

const fixedData = `{"no":1024,"date":"2026-10-13","customer":{"name":"ACME"},"amount":12000,"total":13200,"paid":true,
	"memo":"Thank you.\nSecond line","photo":"stamp.png","url":"https://example.com/c/1024","jan":"4901234567894","tags":["a","b","c"]}`

// TestFillFixed fills a form of one page: values of every kind.
func TestFillFixed(t *testing.T) {
	res, text := fillFile(t, "form-fixed.xlsx", fixedData, &TemplateData{Files: fstest.MapFS{"stamp.png": {Data: testPNG()}}})
	pages := text["Form"]
	if len(pages) != 1 {
		t.Fatalf("%d pages", len(pages))
	}
	p := pages[0]
	for _, want := range [][]string{
		{"明細書 Certificate"},
		{"No", "1024"},                            // a number
		{"Date", "2026/10/13", "2026-10"},         // a date in the format of the cell, and in that of the tag
		{"Name", "ACME さま"},                       // a value in text
		{"Amount", "12,000", "13,200.0", "large"}, // the formulas that read the value
		{"Total", "13,200.00"},                    // the format of the tag
		{"Paid", "■ paid"},                        // a section and its inverse
		{"Tags: a, b, c"},                         // a list in text
		{"1 / 1"},
	} {
		if !lines(p, want...) {
			t.Errorf("the page lacks %q:\n%s", want, p)
		}
	}
	if strings.Contains(p, "missing") || strings.Contains(p, "□") || strings.Contains(p, "}}\n") {
		t.Errorf("tags are left in the page:\n%s", p)
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], `"missing"`) || !strings.Contains(res.Warnings[0], "A15") {
		t.Errorf("warnings: %q", res.Warnings)
	}
	// the picture is in the document, in a figure named after its tag
	imgs := 0
	for _, part := range res.Doc.Parts() {
		if part.Type == bdf.PartImage {
			imgs++
		}
	}
	if imgs != 1 {
		t.Errorf("%d images", imgs)
	}
	// and so are the QR code and the barcode, which are drawn
	if got := symbols(t, res.Doc); len(got) != 5 || got[0] != "figure photo" || got[1] != "figure https://example.com/c/1024" || got[3] != "figure 4901234567894" {
		t.Errorf("figures: %q", got)
	} else if !strings.HasSuffix(got[4], " contours below bars") {
		t.Errorf("JAN lacks human-readable vector digits: %q", got[4])
	}
	// the fields of the template
	var names []string
	for _, f := range res.Fields {
		names = append(names, f.Kind+" "+strings.Join(f.In, ">")+" "+f.Name)
	}
	want := []string{"value  no", "value  date", "value  customer.name", "value  amount", "value  total", "section  paid", "value  memo",
		"image  photo", "image  url", "section  tags", "value tags .", "value  missing", "image  jan"}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("fields = %q", names)
	}
	// strictly, a name the data lacks is an error
	f, _ := os.Open("testdata/form-fixed.xlsx")
	defer f.Close()
	st, _ := f.Stat()
	var data map[string]any
	json.Unmarshal([]byte(fixedData), &data)
	if _, err := Fill(f, st.Size(), &TemplateData{Data: data, Strict: true}, testOptions()); err == nil || !strings.Contains(err.Error(), `no "missing"`) && !strings.Contains(err.Error(), "picture") {
		t.Errorf("strict fill: %v", err)
	}
	// a paid of false shows the other box
	_, text = fillFile(t, "form-fixed.xlsx", `{"paid":false,"tags":[]}`, nil)
	if p := text["Form"][0]; !lines(p, "Paid", "□ paid") || !lines(p, "Tags:") {
		t.Errorf("the empty form:\n%s", p)
	}
}

func itemsJSON(n int) string {
	var items []string
	for i := 1; i <= n; i++ {
		items = append(items, fmt.Sprintf(`{"name":"Item %02d","qty":%d,"price":%d}`, i, i%5+1, 100*i))
	}
	return `[` + strings.Join(items, ",") + `]`
}

// TestFillRows fills a table that grows: a row for each item, the
// formulas of the rows and those under them, the pages.
func TestFillRows(t *testing.T) {
	res, text := fillFile(t, "form-rows.xlsx", `{"no":"A-7","customer":"ACME","note":"Thank you.","items":`+itemsJSON(60)+`}`, nil)
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "TODAY()") {
		t.Errorf("warnings: %q", res.Warnings)
	}
	pages := text["Invoice"]
	if len(pages) != 3 {
		t.Fatalf("%d pages", len(pages))
	}
	all := strings.Join(pages, "\n")
	total := 0
	for i := 1; i <= 60; i++ {
		qty, price := i%5+1, 100*i
		total += qty * price
		if !lines(all, itoa2(i), fmt.Sprintf("Item %02d", i), itoa2(qty), comma(price), comma(qty*price)) {
			t.Fatalf("the row of item %d is not in the pages:\n%s", i, all)
		}
	}
	tax := (total + 5) / 10
	for _, want := range [][]string{{"Total", comma(total)}, {"Tax", comma(tax)}, {"60", "Sum", comma(total + tax)}, {"Thank you."}} {
		if !lines(all, want...) {
			t.Errorf("the pages lack %q", want)
		}
	}
	if strings.Contains(all, "no items") || strings.Contains(all, "{{") {
		t.Errorf("the pages show what they should not:\n%s", all)
	}
	// the title rows are on every page, and the break before the note is kept
	for i, p := range pages {
		if !strings.Contains(p, "Invoice A-7") || !lines(p, "No", "項目", "数", "Price", "Amount") || !strings.Contains(p, fmt.Sprintf("%d / 3", i+1)) {
			t.Errorf("page %d lacks the titles or its number:\n%s", i+1, p)
		}
	}
	if !strings.Contains(pages[2], "Thank you.") || strings.Contains(pages[2], "Item") || strings.Contains(pages[1], "Thank you.") {
		t.Errorf("the note is not alone on the last page:\n%s", pages[2])
	}
	// the other sheet reads the filled one
	if p := text["Summary"][0]; !lines(p, "Items", "60", "Sum", itoa2(total+tax), "First", "Item 01", "Unknown") {
		t.Errorf("the summary:\n%s", p)
	}

	// no items: the row is gone, the sum of nothing is 0, the inverse shows
	_, text = fillFile(t, "form-rows.xlsx", `{"no":1,"customer":"x","note":"","items":[]}`, nil)
	if p := strings.Join(text["Invoice"], "\n"); !lines(p, "Amount", "no items", "Total", "0", "Tax", "0", "0", "Sum", "0") {
		t.Errorf("the table without items:\n%s", p)
	}
	if p := text["Summary"][0]; !lines(p, "Items", "0", "Sum", "0", "First", "#REF!") {
		t.Errorf("the summary without items:\n%s", p)
	}
}

func itoa2(n int) string { return fmt.Sprint(n) }

// comma writes a number with separators of thousands.
func comma(n int) string {
	s := fmt.Sprint(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

// TestFillPages fills a slip of ten rows a page: the data gives the rows
// of each page.
func TestFillPages(t *testing.T) {
	page := func(first, n, in, out int) string {
		var rows []string
		for i := first; i < first+n; i++ {
			rows = append(rows, fmt.Sprintf(`{"no":%d,"name":"Part %02d","qty":%d,"amount":%d}`, i, i, i, 1000*i))
		}
		return fmt.Sprintf(`{"carry_in":%d,"carry_out":%d,"rows":[%s]}`, in, out, strings.Join(rows, ","))
	}
	data := `{"slip":"D-1","customer":"ACME","total":300000,"pages":[` + page(1, 10, 0, 55000) + `,` + page(11, 10, 55000, 210000) + `,` + page(21, 4, 210000, 300000) + `]}`
	res, text := fillFile(t, "form-pages.xlsx", data, nil)
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	pages := text["Slip"]
	if len(pages) != 3 {
		t.Fatalf("%d pages", len(pages))
	}
	for i, want := range [][]string{
		{"Delivery D-1", "1 / 3", "To: ACME", "From", "0"},
		{"Delivery D-1", "2 / 3", "To: ACME", "From", "55,000"},
		{"Delivery D-1", "3 / 3", "To: ACME", "From", "210,000"},
	} {
		if !lines(pages[i], want...) || !lines(pages[i], "No", "項目", "数", "Amount") {
			t.Errorf("page %d lacks its head %q:\n%s", i+1, want, pages[i])
		}
	}
	// the rows of each page, the sum of the page and what is carried over
	if !lines(pages[0], "10", "Part 10", "10", "10,000", "Page", "55,000", "To", "55,000") || strings.Contains(pages[0], "Part 11") {
		t.Errorf("page 1:\n%s", pages[0])
	}
	if !lines(pages[1], "11", "Part 11", "11", "11,000") || !lines(pages[1], "Page", "155,000", "To", "210,000") {
		t.Errorf("page 2:\n%s", pages[1])
	}
	// the last page: four rows, the mark of the end in the first empty
	// one, the total
	if !lines(pages[2], "24", "Part 24", "24", "24,000", "-- end --", "Page", "90,000", "Total 300,000", "To", "300,000") {
		t.Errorf("page 3:\n%s", pages[2])
	}
	if strings.Count(pages[2], "-- end --") != 1 || strings.Contains(pages[0]+pages[1], "Total") || strings.Contains(pages[0]+pages[1], "-- end --") {
		t.Errorf("the marks of the last page are elsewhere:\n%s", strings.Join(pages, "\n--\n"))
	}
	// the slots: ten rows
	var slots tag.Field
	for _, f := range res.Fields {
		if f.Name == "rows" {
			slots = f
		}
	}
	if slots.Mode != tag.ModeSlots || slots.Slots != 10 || !reflect.DeepEqual(slots.In, []string{"pages"}) {
		t.Errorf("the field of the slots: %+v", slots)
	}

	// more rows than a page has: an error that says so
	f, _ := os.Open("testdata/form-pages.xlsx")
	defer f.Close()
	st, _ := f.Stat()
	var over map[string]any
	json.Unmarshal([]byte(`{"pages":[`+page(1, 11, 0, 0)+`]}`), &over)
	if _, err := Fill(f, st.Size(), &TemplateData{Data: over}, testOptions()); err == nil || !strings.Contains(err.Error(), "11 items for the 10 rows") || !strings.Contains(err.Error(), "A4") {
		t.Errorf("too many rows: %v", err)
	}
}

// book makes a workbook of sheets given as the content of their sheetData
// (and what follows it), with defined names.
func book(sheets []string, names string) []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	add := func(name, data string) {
		w, _ := zw.Create(name)
		w.Write([]byte(data))
	}
	wb := `<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets>`
	rels := `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`
	for i, s := range sheets {
		n := i + 1
		wb += fmt.Sprintf(`<sheet name="S%d" sheetId="%d" r:id="rId%d"/>`, n, n, n)
		rels += fmt.Sprintf(`<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet%d.xml"/>`, n, n)
		data, rest, _ := strings.Cut(s, "</sheetData>")
		add(fmt.Sprintf("xl/worksheets/sheet%d.xml", n), `<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`+data+`</sheetData>`+rest+`</worksheet>`)
	}
	add("xl/workbook.xml", wb+`</sheets>`+names+`</workbook>`)
	add("xl/_rels/workbook.xml.rels", rels+`</Relationships>`)
	zw.Close()
	return buf.Bytes()
}

// rowXML writes a row of cells: a text is an inline string, "=…" a
// formula, a number a number, "" no cell.
func rowXML(r int, cells ...string) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<row r="%d">`, r)
	for i, c := range cells {
		ref := cellRef(r-1, i)
		switch {
		case c == "":
		case strings.HasPrefix(c, "="):
			fmt.Fprintf(&b, `<c r="%s"><f>%s</f><v>0</v></c>`, ref, c[1:])
		case strings.Trim(c, "0123456789.-") == "":
			fmt.Fprintf(&b, `<c r="%s"><v>%s</v></c>`, ref, c)
		default:
			c = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(c)
			fmt.Fprintf(&b, `<c r="%s" t="inlineStr"><is><t>%s</t></is></c>`, ref, c)
		}
	}
	b.WriteString(`</row>`)
	return b.String()
}

// filled fills a workbook and returns its sheets as they are before they
// are drawn: the values of the cells as text ("" for an empty one).
func filled(t *testing.T, data []byte, jsonData string) (*converter, map[string][][]string, error) {
	t.Helper()
	td := &TemplateData{}
	if err := json.Unmarshal([]byte(jsonData), &td.Data); err != nil {
		t.Fatal(err)
	}
	f := newFiller(td)
	opts := testOptions()
	opts.fill = f
	c, wb, err := openWorkbook(bytes.NewReader(data), int64(len(data)), opts)
	if err != nil {
		t.Fatal(err)
	}
	c.loadPrintNames(wb)
	if err := f.run(c); err != nil {
		return c, nil, err
	}
	out := map[string][][]string{}
	for _, ts := range f.sheets {
		var grid [][]string
		for _, rw := range ts.ws.rows {
			for len(grid) <= rw.idx {
				grid = append(grid, nil)
			}
			for _, cl := range rw.cells {
				for len(grid[rw.idx]) <= cl.col {
					grid[rw.idx] = append(grid[rw.idx], "")
				}
				grid[rw.idx][cl.col] = toStr(cellFval(&cl))
			}
		}
		// without the empty cells at the ends of the rows
		for i, r := range grid {
			for len(r) > 0 && r[len(r)-1] == "" {
				r = r[:len(r)-1]
			}
			if grid[i] = r; len(r) == 0 {
				grid[i] = nil
			}
		}
		out[ts.ws.name] = grid
	}
	return c, out, nil
}

func grid(rows ...string) [][]string {
	var out [][]string
	for _, r := range rows {
		if r == "" {
			out = append(out, nil)
			continue
		}
		out = append(out, strings.Split(r, "|"))
	}
	return out
}

// TestFillGroups fills rows inside rows: a group for each item with rows
// for its items, and the formulas of each.
func TestFillGroups(t *testing.T) {
	sheet := rowXML(1, "Report", "=SUM(B3:B5)") +
		rowXML(2, "{{#groups | rows}}{{name}}", "") +
		rowXML(3, "{{#items}}{{label}}", "{{n}}{{/items}}", "=B3*2") +
		rowXML(4, "sub", "=SUM(B3:B3)", "=B4+B$1") +
		rowXML(5, "{{/groups}}") +
		rowXML(6, "all", "=SUM(B2:B5)", "=COUNT(B3:B3)", "=B3")
	data := `{"groups":[{"name":"G1","items":[{"label":"a","n":1},{"label":"b","n":2}]},{"name":"G2","items":[]},{"name":"G3","items":[{"label":"c","n":10}]}]}`
	c, got, err := filled(t, book([]string{sheet}, ""), data)
	if err != nil {
		t.Fatal(err)
	}
	// a formula in a copy reads its own copy (the sum of a group); one
	// outside reads from the first copy to the last (what is between them
	// too: the range is one range, as when rows are inserted in Excel)
	want := grid(
		"Report|26",
		"G1",
		"a|1|2",
		"b|2|4",
		"sub|3|29",
		"",
		"G2",
		"sub|0|26",
		"",
		"G3",
		"c|10|20",
		"sub|10|36",
		"",
		"all|26|5|1",
	)
	if !reflect.DeepEqual(got["S1"], want) {
		t.Errorf("the filled sheet:\n%s\nwant:\n%s", show(got["S1"]), show(want))
	}
	if len(c.warnings) != 0 {
		t.Errorf("warnings: %q", c.warnings)
	}
}

func show(g [][]string) string {
	var b strings.Builder
	for i, r := range g {
		fmt.Fprintf(&b, "%2d %s\n", i+1, strings.Join(r, "|"))
	}
	return b.String()
}

// TestFillMoves checks that what refers to rows follows them: merged
// ranges, conditional formats, links, the print area and titles, page
// breaks, and the references of other sheets.
func TestFillMoves(t *testing.T) {
	sheet := rowXML(1, "head") +
		rowXML(2, "{{#items}}{{v}}", "x{{/items}}") +
		rowXML(3, "foot", "=A2") +
		`</sheetData><mergeCells><mergeCell ref="A1:B1"/><mergeCell ref="B2:C2"/><mergeCell ref="C1:C3"/><mergeCell ref="A3:A4"/></mergeCells>` +
		`<conditionalFormatting sqref="A2 A3:A4"><cfRule type="cellIs" operator="equal" priority="1"><formula>1</formula></cfRule></conditionalFormatting>` +
		`<hyperlinks><hyperlink ref="B2" location="x"/><hyperlink ref="A3"/></hyperlinks>` +
		`<rowBreaks count="2"><brk id="1" man="1"/><brk id="2" man="1"/></rowBreaks>`
	other := rowXML(1, "=S1!A2", "=SUM(S1!A2:A2)", "=S1!A3", "=COUNTA(S1!A1:A3)")
	names := `<definedNames><definedName name="_xlnm.Print_Area" localSheetId="0">S1!$A$2:$B$3</definedName>` +
		`<definedName name="_xlnm.Print_Titles" localSheetId="0">S1!$1:$2</definedName>` +
		`<definedName name="_bdf.Print_TitleScope" localSheetId="0">S1!$A$1:$B$2</definedName>` +
		`<definedName name="_bdf.Print_KeepTogether" localSheetId="0">S1!$B$2:$C$2,S1!$A$3:$A$4</definedName></definedNames>`
	c, got, err := filled(t, book([]string{sheet, other}, names), `{"items":[{"v":5},{"v":6},{"v":7}]}`)
	if err != nil {
		t.Fatal(err)
	}
	if want := grid("head", "5|x", "6|x", "7|x", "foot|5"); !reflect.DeepEqual(got["S1"], want) {
		t.Errorf("the filled sheet:\n%s", show(got["S1"]))
	}
	if want := grid("5|18|foot|5"); !reflect.DeepEqual(got["S2"], want) {
		t.Errorf("the other sheet:\n%s", show(got["S2"]))
	}
	ws := c.parsed["S1"]
	var merges []string
	for _, m := range ws.merges {
		merges = append(merges, cellRef(m.r0, m.c0)+":"+cellRef(m.r1, m.c1))
	}
	if want := []string{"A1:B1", "B2:C2", "B3:C3", "B4:C4", "C1:C5", "A5:A6"}; !reflect.DeepEqual(merges, want) {
		t.Errorf("merges = %q, want %q", merges, want)
	}
	if got, want := ws.cfs[0].AttrStr("sqref", ""), "A2:A2 A3:A3 A4:A4 A5:A6"; got != want {
		t.Errorf("sqref = %q, want %q", got, want)
	}
	var links []string
	for _, l := range ws.links {
		links = append(links, l.AttrStr("ref", ""))
	}
	if want := []string{"B2:B2", "B3:B3", "B4:B4", "A5:A5"}; !reflect.DeepEqual(links, want) {
		t.Errorf("links = %q, want %q", links, want)
	}
	// a break before the repeated row is before each copy; the one after it has moved
	if want := []int{1, 2, 3, 4}; !reflect.DeepEqual(ws.print.rowBreaks, want) {
		t.Errorf("breaks = %v, want %v", ws.print.rowBreaks, want)
	}
	pn := c.printNames[0]
	if a := pn.areas[0]; a != (cellRange{1, 0, 4, 1}) || *pn.titleRows != [2]int{0, 3} {
		t.Errorf("print area %v, titles %v", a, *pn.titleRows)
	}
	if got, want := *pn.titleScope, (cellRange{0, 0, 3, 1}); got != want {
		t.Errorf("title scope = %v, want %v", got, want)
	}
	if got, want := pn.keepRows, []cellRange{{1, 1, 1, 2}, {2, 1, 2, 2}, {3, 1, 3, 2}, {4, 0, 5, 0}}; !reflect.DeepEqual(got, want) {
		t.Errorf("keep-together ranges = %v, want %v", got, want)
	}
}

// TestFillShared checks the formulas that Excel shares between the cells
// of a column, and those the evaluator does not know.
func TestFillShared(t *testing.T) {
	sheet := `<row r="1"><c r="A1"><v>1</v></c><c r="B1"><f t="shared" ref="B1:B3" si="0">A1*10</f><v>0</v></c></row>` +
		`<row r="2"><c r="A2" t="inlineStr"><is><t>{{n}}</t></is></c><c r="B2"><f t="shared" si="0"/><v>0</v></c></row>` +
		`<row r="3"><c r="A3"><v>3</v></c><c r="B3"><f t="shared" si="0"/><v>0</v></c><c r="C3"><f>B3+B2+B1</f><v>9</v></c></row>` +
		`<row r="4"><c r="A4"><f>NOSUCH(A1)</f><v>7</v></c><c r="B4"><f t="array" ref="B4">A1:A3*2</f><v>8</v></c><c r="C4"><f>C4+1</f><v>1</v></c><c r="D4"><f>1/0</f><v>0</v></c></row>`
	c, got, err := filled(t, book([]string{sheet}, ""), `{"n":2}`)
	if err != nil {
		t.Fatal(err)
	}
	// (a formula that reads its own cell reads it empty)
	if want := grid("1|10", "2|20", "3|30|60", "||1|#DIV/0!"); !reflect.DeepEqual(got["S1"], want) {
		t.Errorf("the filled sheet:\n%s", show(got["S1"]))
	}
	if len(c.warnings) != 2 {
		t.Errorf("warnings: %q", c.warnings)
	}
}

// TestFillValues checks the kinds of values in cells that are one tag,
// and in text.
func TestFillValues(t *testing.T) {
	sheet := rowXML(1, "{{n}}", "{{s}}", "{{b}}", "{{d}}", "{{z}}", "{{l}}", "{{o.k}}") +
		rowXML(2, "n={{n}} s={{s}} b={{b}} d={{d}} z={{z}}.", `{{n | fmt:"0.00"}}`, `{{d | fmt:"yyyy/mm"}}`, `<{{n | fmt:"0.0"}};{{ds | fmt:"m/d"}}>`, "{{big}}") +
		rowXML(3, " {{ n }} ", "{{#o}}{{k}}{{/o}}", "{{#l}}[{{.}}:{{_index}}/{{_count}}]{{/l}}", "{{^z}}none{{/z}}", "{{#l}}{{#_first}}F{{/_first}}{{#_last}}L{{/_last}}{{/l}}")
	c, got, err := filled(t, book([]string{sheet}, ""), `{"n":1234.5,"s":"text","b":true,"d":"2026-10-13","z":null,"l":["x","y"],"o":{"k":"v"},"ds":"2026-01-02T03:04:05","big":1e15}`)
	if err != nil {
		t.Fatal(err)
	}
	want := grid(
		"1234.5|text|TRUE|2026-10-13|||v",
		"n=1234.5 s=text b=TRUE d=2026-10-13 z=.|1234.5|46308|<1234.5;1/2>|1E+15",
		"1234.5|v|[x:1/2][y:2/2]|none|FL",
	)
	if !reflect.DeepEqual(got["S1"], want) {
		t.Errorf("the filled sheet:\n%s\nwant:\n%s", show(got["S1"]), show(want))
	}
	ws := c.parsed["S1"]
	// the kinds: a number, a text, a boolean; a date where one is asked for
	kinds := func(r int) []cellKind {
		var out []cellKind
		for _, cl := range ws.rows[r].cells {
			out = append(out, cl.kind)
		}
		return out
	}
	if got, want := kinds(0), []cellKind{cellNum, cellStr, cellBool, cellStr, cellBlank, cellBlank, cellStr}; !reflect.DeepEqual(got, want) {
		t.Errorf("kinds of row 1 = %v", got)
	}
	if got, want := kinds(1), []cellKind{cellStr, cellNum, cellNum, cellStr, cellNum}; !reflect.DeepEqual(got, want) {
		t.Errorf("kinds of row 2 = %v", got)
	}
	// the formats of the tags are the formats of the cells
	format := func(r, col int) string {
		cl := ws.cellAt(r, col)
		return c.st.numFormat(c.st.xf(cl.style).numFmt).formatNumber(cl.num, false).text()
	}
	if a, b := format(1, 1), format(1, 2); a != "1234.50" || b != "2026/10" {
		t.Errorf("formatted: %q %q", a, b)
	}
	if len(c.warnings) != 1 || !strings.Contains(c.warnings[0], `"l" is a list`) {
		t.Errorf("warnings: %q", c.warnings)
	}
}

// TestFillDates checks the dates of the data: a time is a serial date in
// the format Excel gives a date typed into a cell.
func TestFillDates(t *testing.T) {
	td := &TemplateData{Data: map[string]any{
		"day":  time.Date(2026, 10, 13, 0, 0, 0, 0, time.UTC),
		"time": time.Date(2026, 10, 13, 9, 30, 0, 0, time.FixedZone("JST", 9*3600)),
		"old":  time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
		"far":  time.Date(2400, 1, 1, 0, 0, 0, 0, time.UTC),
	}}
	f := newFiller(td)
	opts := testOptions()
	opts.fill = f
	data := book([]string{rowXML(1, "{{day}}", "{{time}}", "on {{day}} at {{time}}", "{{old}}", "{{far}}")}, "")
	c, _, err := openWorkbook(bytes.NewReader(data), int64(len(data)), opts)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.run(c); err != nil {
		t.Fatal(err)
	}
	ws := c.parsed["S1"]
	shown := func(col int) string {
		cl := ws.cellAt(0, col)
		if cl.kind == cellStr {
			return cl.text.plain
		}
		return c.st.numFormat(c.st.xf(cl.style).numFmt).formatNumber(cl.num, false).text()
	}
	for col, want := range []string{"10/13/2026", "10/13/2026 9:30", "on 10/13/2026 at 10/13/2026 9:30", "1/1/1900", "1/1/2400"} {
		if got := shown(col); got != want {
			t.Errorf("column %d shows %q, want %q", col, got, want)
		}
	}
	if n := ws.cellAt(0, 0).num; n != 46308 {
		t.Errorf("the serial date of the day is %v", n)
	}
	if n := ws.cellAt(0, 3).num; n != 1 {
		t.Errorf("the serial date of 1900-01-01 is %v", n)
	}
}

// TestFillErrors checks the mistakes of templates, told with their cells.
func TestFillErrors(t *testing.T) {
	for name, c := range map[string]struct{ sheet, want string }{
		"unclosed tag":     {rowXML(1, "", "{{name"), "cell B1: the tag {{name is not closed"},
		"unclosed section": {rowXML(1, "{{#a}}") + rowXML(2, "x"), "cell A1: the section {{#a}} is not closed"},
		"stray close":      {rowXML(1, "{{/a}}"), "cell A1: {{/a}} closes no section"},
		"crossed":          {rowXML(1, "{{#a}}", "{{#b}}") + rowXML(2, "{{/a}}", "{{/b}}"), "cell A2: {{/a}} closes the section {{#b}} of cell B1"},
		"crossed in cell":  {rowXML(1, "{{#a}}{{#b}}{{/a}}{{/b}}"), "cell A1: {{/a}} closes the section {{#b}}"},
		"two in a row":     {rowXML(1, "{{#a}}", "{{/a}}", "{{#b}}", "{{/b}}"), "cell C1: the sections {{#a}} and {{#b}} both repeat row 1"},
		"section in slots": {rowXML(1, "{{#a}}") + rowXML(2, "{{#b}}{{x}}", "{{/b}}") + rowXML(3, "{{/a}}"), "cell A2: the section {{#b}} is in the rows of {{#a}}, which are slots"},
		"bad option":       {rowXML(1, "{{a | upper}}"), "cell A1: the tag {{a}} has an option it does not take: upper"},
		"inverted pages":   {rowXML(1, "{{^a | pages}}", "{{/a}}"), `cell A1: the inverted section {{^a}} does not repeat: it takes no "pages"`},
	} {
		_, _, err := filled(t, book([]string{c.sheet}, ""), `{"a":[],"b":[]}`)
		if err == nil || !strings.Contains(err.Error(), c.want) || !strings.Contains(err.Error(), `sheet "S1"`) {
			t.Errorf("%s: %v\nwant an error with %q", name, err, c.want)
		}
	}
	// Inspect tells the same without data
	data := book([]string{rowXML(1, "{{#a}}")}, "")
	if _, err := Inspect(bytes.NewReader(data), int64(len(data)), testOptions()); err == nil || !strings.Contains(err.Error(), "is not closed") {
		t.Errorf("Inspect: %v", err)
	}
	data = book([]string{rowXML(1, "{{#a}}{{x}}", "{{/a}}") + rowXML(2, "{{y}}", "{{%p}}")}, "")
	fields, err := Inspect(bytes.NewReader(data), int64(len(data)), testOptions())
	if err != nil || len(fields) != 4 || fields[0].Mode != tag.ModeRows || fields[1].In[0] != "a" || fields[3].Kind != tag.FieldImage {
		t.Errorf("Inspect: %+v, %v", fields, err)
	}
}

// TestFillModes checks the ways the rows between the tags of a section
// repeat.
func TestFillModes(t *testing.T) {
	data := `{"a":[{"v":1},{"v":2}],"one":{"v":9},"no":false}`
	for name, c := range map[string]struct {
		sheet string
		want  [][]string
	}{
		// rows: the rows between the tags, for each item
		"rows": {rowXML(1, "{{#a | rows}}{{v}}") + rowXML(2, "-{{/a}}") + rowXML(3, "end"), grid("1", "-", "2", "-", "end")},
		// slots by default: the rows that are there
		"slots": {rowXML(1, "{{#a}}{{v}}", "s") + rowXML(2, "", "t") + rowXML(3, "{{/a}}", "u") + rowXML(4, "end"), grid("1|s", "2|t", "|u", "end")},
		// a slot with tags of its own keeps them
		"own slot": {rowXML(1, "{{#a}}{{v}}") + rowXML(2, "<{{v}}>{{/a}}"), grid("1", "<2>")},
		// an object is one item; a false value none
		"object": {rowXML(1, "{{#one}}{{v}}", "{{/one}}") + rowXML(2, "{{#no}}x", "y{{/no}}") + rowXML(3, "end"), grid("9", "end")},
		// an inverted section of rows: there when the value is false
		"inverted": {rowXML(1, "{{^no}}shown", "{{/no}}") + rowXML(2, "{{^a}}hidden", "{{/a}}") + rowXML(3, "end"), grid("shown", "end")},
		// the rows around the tags see the item too
		"whole row": {rowXML(1, "{{v}}", "{{#a}}", "{{/a}}", "{{v}}"), grid("1|||1", "2|||2")},
		// a section that says that it repeats rows, with its tags in a cell:
		// the row, not the text of the cell
		"one cell":         {rowXML(1, "{{#a | rows}}{{v}}{{/a}}", "x") + rowXML(2, "end"), grid("1|x", "2|x", "end")},
		"one cell, nested": {rowXML(1, "{{#a | rows}}{{#one | rows}}{{v}}{{/one}}{{/a}}", "{{v}}"), grid("9|9", "9|9")},
		"one cell, not":    {rowXML(1, "{{^no | rows}}shown{{/no}}") + rowXML(2, "{{^a | rows}}hidden{{/a}}") + rowXML(3, "end"), grid("shown", "end")},
		// the variables of an empty slot
		"blank": {rowXML(1, "{{#a}}{{v}}{{#_blank}}-{{/_blank}}{{#_blank_first}}!{{/_blank_first}}") + rowXML(2, "") + rowXML(3, "") + rowXML(4, "{{/a}}"), grid("1", "2", "-!", "-")},
	} {
		_, got, err := filled(t, book([]string{c.sheet}, ""), data)
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		if !reflect.DeepEqual(got["S1"], c.want) {
			t.Errorf("%s:\n%s\nwant:\n%s", name, show(got["S1"]), show(c.want))
		}
	}
}

// TestFillLimits checks that a template cannot be filled past what a
// sheet holds.
func TestFillLimits(t *testing.T) {
	// a section in a section: two thousand times a thousand rows
	list := func(inner string, n int) string {
		return "[" + strings.TrimSuffix(strings.Repeat(inner+",", n), ",") + "]"
	}
	data := `{"a":` + list(`{"b":`+list(`{"c":[1,2,3]}`, 1000)+`}`, 2000) + `}`
	sheet := rowXML(1, "{{#a | rows}}{{#b | rows}}", "{{#c}}{{.}}{{/c}}", "{{/b}}{{/a}}")
	start := time.Now()
	_, _, err := filled(t, book([]string{sheet}, ""), data)
	if err == nil || !strings.Contains(err.Error(), "cells") && !strings.Contains(err.Error(), "rows of a sheet") {
		t.Errorf("a template of two million rows: %v", err)
	}
	if d := time.Since(start); d > 60*time.Second {
		t.Errorf("it took %v", d)
	}
	// a long text in a cell is cut
	c, got, err := filled(t, book([]string{rowXML(1, "{{#l}}{{s}}{{/l}}")}, ""), `{"l":`+list(`{"s":"`+strings.Repeat("x", 100)+`"}`, 1000)+`}`)
	if err != nil || len([]rune(got["S1"][0][0])) != maxCellText || len(c.warnings) != 1 {
		t.Errorf("a long text: %d characters, %q, %v", len(got["S1"][0][0]), c.warnings, err)
	}
}

// TestValues reads a workbook as data.
func TestValues(t *testing.T) {
	sheet := rowXML(1, "name", "qty", "day") + rowXML(3, "a", "2") +
		`<row r="4"><c r="A4" t="b"><v>1</v></c><c r="C4" s="1"><v>46308</v></c><c r="E4"/></row>`
	data := book([]string{sheet, rowXML(2, "", "x")}, "")
	// a date format for the cell of style 1
	data = withPart(t, data, "xl/styles.xml", `<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><cellXfs count="2"><xf numFmtId="0"/><xf numFmtId="14"/></cellXfs></styleSheet>`,
		`<Relationship Id="rId9" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>`)
	got, err := Values(bytes.NewReader(data), int64(len(data)), nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []SheetValues{
		{Name: "S1", Rows: [][]any{{"name", "qty", "day"}, nil, {"a", 2.0}, {true, nil, time.Date(2026, 10, 13, 0, 0, 0, 0, time.UTC)}}},
		{Name: "S2", Rows: [][]any{nil, {nil, "x"}}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Values = %#v", got)
	}
	// values far from the others are not sheets of empty cells
	last := `<row r="1048576"><c r="XFD1048576"><v>1</v></c></row>`
	far := book([]string{last, last, last, last, last}, "")
	if _, err := Values(bytes.NewReader(far), int64(len(far)), nil); err == nil || !strings.Contains(err.Error(), "cells") {
		t.Errorf("values in the last cells of five sheets: %v", err)
	}
}

// withPart adds a part and its relationship to a workbook made by book.
func withPart(t *testing.T, data []byte, name, content, rel string) []byte {
	t.Helper()
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, f := range zr.File {
		rc, _ := f.Open()
		var b bytes.Buffer
		b.ReadFrom(rc)
		rc.Close()
		s := b.String()
		if f.Name == "xl/_rels/workbook.xml.rels" {
			s = strings.Replace(s, "</Relationships>", rel+"</Relationships>", 1)
		}
		w, _ := zw.Create(f.Name)
		w.Write([]byte(s))
	}
	w, _ := zw.Create(name)
	w.Write([]byte(content))
	zw.Close()
	return buf.Bytes()
}

// TestFillDepth checks that sections of rows are not nested without end.
func TestFillDepth(t *testing.T) {
	var cells []string
	for i := 0; i <= maxSectionDepth; i++ {
		cells = append(cells, "{{#a}}")
	}
	_, _, err := filled(t, book([]string{rowXML(1, cells...)}, ""), `{}`)
	if err == nil || !strings.Contains(err.Error(), "inside more than") {
		t.Errorf("deep sections: %v", err)
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

// symbols lists the figures and their modules/bars, digit contours and bounds.
func symbols(t *testing.T, doc *bdf.Document) []string {
	t.Helper()
	var out []string
	for _, part := range doc.Parts() {
		if part.Type != bdf.PartObject {
			continue
		}
		o, err := bdf.DecodeObject(part.Data)
		if err != nil {
			t.Fatal(err)
		}
		o.Walk(func(in bdf.Instr) {
			switch in.Op {
			case bdf.OpMark:
				if byte(in.Args[0].(uint64)) == bdf.MarkFigure {
					out = append(out, "figure "+in.Args[1].(string))
				}
			case bdf.OpFillPath:
				p := o.Paths[int(in.Args[0].(uint64))].Inline
				rects, contours, box := codeBounds(p)
				if rects == 0 {
					return
				}
				label := ""
				if contours > 0 {
					label = fmt.Sprintf(" + %d contours below bars", contours)
				}
				out = append(out, fmt.Sprintf("%d rects in %.1f × %.1f%s", rects, box[2]-box[0], box[3]-box[1], label))
			}
		})
	}
	return out
}

// TestFillCodes fills picture tags that ask for a QR code or a barcode:
// the symbol is drawn in the cell, with JAN digits below its bars.
func TestFillCodes(t *testing.T) {
	sheet := `<row r="1" ht="60" customHeight="1">` + strings.TrimPrefix(rowXML(1, "{{%url | qr}}", "{{%jan | barcode:jan}}", "{{%url | qr:H | quiet:0}}", "{{%no | barcode:itf | quiet:0}}"), `<row r="1">`) +
		`</sheetData><cols><col min="1" max="4" width="20" customWidth="1"/></cols>`
	wb := book([]string{sheet}, "")
	fill := func(js string, strict bool) (*Result, error) {
		td := &TemplateData{Strict: strict}
		if err := json.Unmarshal([]byte(js), &td.Data); err != nil {
			t.Fatal(err)
		}
		return Fill(bytes.NewReader(wb), int64(len(wb)), td, testOptions())
	}
	const data = `{"url":"https://example.com/invoice/1024","jan":"4901234567894","no":1024}`
	res, err := fill(data, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	count := func(kind, text string, o *barcode.Options) int {
		s, err := barcode.Encode(kind, text, o)
		if err != nil {
			t.Fatal(err)
		}
		return len(s.Rects())
	}
	// the cells are 60 pt tall, and wider: a QR code is a square in the
	// middle, and JAN reserves room for its digits inside the cell.
	got := symbols(t, res.Doc)
	want := []string{
		"figure https://example.com/invoice/1024", fmt.Sprintf("%d rects in ", count(barcode.QR, "https://example.com/invoice/1024", nil)),
		"figure 4901234567894", fmt.Sprintf("%d rects in ", count(barcode.EAN13, "4901234567894", nil)),
		"figure https://example.com/invoice/1024", fmt.Sprintf("%d rects in 58.5 × 58.5", count(barcode.QR, "https://example.com/invoice/1024", &barcode.Options{Level: "H", Quiet: barcode.NoQuiet})),
		"figure 1024", fmt.Sprintf("%d rects in ", count(barcode.ITF, "1024", &barcode.Options{Quiet: barcode.NoQuiet})),
	}
	if len(got) != len(want) {
		t.Fatalf("got %q, want %q", got, want)
	}
	for i := range want {
		if !strings.HasPrefix(got[i], want[i]) {
			t.Errorf("%d: got %q, want %q…", i, got[i], want[i])
		}
	}
	if !strings.HasSuffix(got[3], " contours below bars") {
		t.Errorf("JAN lacks human-readable vector digits: %q", got[3])
	}
	// Like the ITF cell, JAN has a 133.5 × 58.5 pt drawing frame.
	// Include the digit contours when checking that it stays inside it.
	var rects int
	var width, height float64
	if n, err := fmt.Sscanf(got[3], "%d rects in %f × %f", &rects, &width, &height); err != nil || n != 3 || width > 133.5 || height != 58.5 {
		t.Errorf("JAN and its digits exceed the cell: %q", got[3])
	}
	if !strings.HasSuffix(got[7], " × 58.5") || strings.HasSuffix(got[7], "in 58.5 × 58.5") {
		t.Errorf("the bars are not as large as the cell: %q", got[7])
	}
	for _, part := range res.Doc.Parts() {
		if part.Type == bdf.PartImage {
			t.Error("a code is a picture of pixels")
		}
	}
	// the fields tell which code a picture is
	var fields []string
	for _, f := range res.Fields {
		fields = append(fields, f.Kind+" "+f.Name+" "+f.Code)
	}
	if want := []string{"image url qr", "image jan ean13", "image no itf"}; !reflect.DeepEqual(fields, want) {
		t.Errorf("fields = %q, want %q", fields, want)
	}
	// a value that the code can not hold is left out with a warning, or is
	// an error; without a value there is no code
	res, err = fill(`{"url":"","jan":"49012345678","no":null}`, false)
	if err != nil {
		t.Fatal(err)
	}
	if got := symbols(t, res.Doc); len(got) != 0 {
		t.Errorf("codes of no value: %q", got)
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], `cell B1: the code of "jan" is left out`) {
		t.Errorf("warnings: %q", res.Warnings)
	}
	if _, err = fill(`{"url":"","jan":"49012345678","no":null}`, true); err == nil || !strings.Contains(err.Error(), "B1") {
		t.Errorf("strict: %v", err)
	}
}
