package docx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"math"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
)

// testOptions restricts fonts to the test font directory so that output
// does not depend on the machine.
func testOptions() *Options {
	return &Options{FontDirs: []string{"testdata/fonts"}, NoSystemFonts: true}
}

func convert(t *testing.T, name string, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
	res, err := ConvertFile(filepath.Join("testdata", name), opts)
	if err != nil {
		t.Fatal(err)
	}
	return res, reopen(t, res)
}

// reopen writes a converted document and reads it back.
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

const ns = `xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" ` +
	`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"`

// makeDocx builds a document from the content of its w:body and more parts
// (by name; styles and numbering are related to the document).
func makeDocx(t *testing.T, body string, parts map[string]string) (*Result, *bdf.Reader) {
	t.Helper()
	var buf bytes.Buffer
	z := zip.NewWriter(&buf)
	add := func(name, data string) {
		w, err := z.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(data))
	}
	add("[Content_Types].xml", `<?xml version="1.0"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">`+
		`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>`+
		`<Default Extension="xml" ContentType="application/xml"/>`+
		`<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/></Types>`)
	add("_rels/.rels", `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`+
		`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/></Relationships>`)
	rels := `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`
	for name, data := range parts {
		add(name, data)
		typ := strings.TrimSuffix(strings.TrimPrefix(name, "word/"), ".xml")
		rels += fmt.Sprintf(`<Relationship Id="rId%s" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/%s" Target="%s"/>`,
			typ, typ, strings.TrimPrefix(name, "word/"))
	}
	add("word/_rels/document.xml.rels", rels+`</Relationships>`)
	add("word/document.xml", `<?xml version="1.0"?><w:document `+ns+`><w:body>`+body+`</w:body></w:document>`)
	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
	res, err := Convert(bytes.NewReader(buf.Bytes()), int64(buf.Len()), testOptions())
	if err != nil {
		t.Fatal(err)
	}
	return res, reopen(t, res)
}

func p(text string, pPr ...string) string {
	return `<w:p><w:pPr>` + strings.Join(pPr, "") + `</w:pPr><w:r><w:t xml:space="preserve">` + text + `</w:t></w:r></w:p>`
}

// a small A5-like page: 300 × 400 pt with 50 pt margins
const smallPage = `<w:sectPr><w:pgSz w:w="6000" w:h="8000"/><w:pgMar w:top="1000" w:bottom="1000" w:left="1000" w:right="1000" w:header="500" w:footer="500"/></w:sectPr>`

// runs returns the text runs of a page of a view (every layer, or those
// with the given role), in page coordinates.
func runs(t *testing.T, r *bdf.Reader, view, page int, role string) []bdf.TextRun {
	t.Helper()
	var out []bdf.TextRun
	for _, l := range r.Manifest.Views[view].Pages[page].Layers {
		if role != "" && l.Role != role {
			continue
		}
		o, err := r.Object(l.Obj)
		if err != nil {
			t.Fatal(err)
		}
		rs, err := bdf.ExtractText(o, func(h bdf.Hash) *bdf.ObjectPart {
			c, _ := r.Object(h)
			return c
		})
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, rs...)
	}
	return out
}

// text joins runs as extraction does: breaks become newlines.
func text(rs []bdf.TextRun) string {
	var b strings.Builder
	for i, r := range rs {
		if i > 0 {
			switch r.Sep {
			case bdf.SepSpace:
				b.WriteByte(' ')
			case bdf.SepBreak:
				b.WriteByte('\n')
			}
		}
		b.WriteString(r.Text)
	}
	return b.String()
}

func pageText(t *testing.T, r *bdf.Reader, view, page int, role string) string {
	return text(runs(t, r, view, page, role))
}

// find returns the first run holding s.
func find(rs []bdf.TextRun, s string) (bdf.TextRun, bool) {
	for _, r := range rs {
		if strings.Contains(r.Text, s) {
			return r, true
		}
	}
	return bdf.TextRun{}, false
}

func TestConvertBasic(t *testing.T) {
	res, r := convert(t, "basic.docx", testOptions())
	if len(res.Warnings) > 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	m := r.Manifest
	if len(m.Views) != 2 || m.Views[0].Kind != bdf.ViewFlow || m.Views[1].Kind != bdf.ViewScroll {
		t.Fatalf("views = %+v", m.Views)
	}
	if res.Pages != 3 || len(m.Views[0].Pages) != 3 {
		t.Fatalf("%d pages", res.Pages)
	}
	if m.Meta.Source != "docx" || m.Meta.DC.Title.First() != "BDF Word test" || m.Meta.DC.Language.First() != "en-US" {
		t.Errorf("meta = %+v", m.Meta)
	}
	pg := m.Views[0].Pages[0]
	if pg.W != 595.3 || pg.H != 841.9 || pg.Body == nil || pg.Body.X != 72 || pg.Body.W != 451.3 {
		t.Errorf("page = %v×%v body %+v", pg.W, pg.H, pg.Body)
	}
	var roles []string
	for _, l := range pg.Layers {
		roles = append(roles, l.Role)
	}
	if !slices.Equal(roles, []string{"header", "footer", "body"}) {
		t.Errorf("roles = %v", roles)
	}
	// page fields in headers and footers
	for i := range 2 {
		want := fmt.Sprintf("Page %d of 3", i+1)
		if h := pageText(t, r, 0, i, "header"); !strings.Contains(h, want) {
			t.Errorf("header of page %d = %q", i+1, h)
		}
		if f := pageText(t, r, 0, i, "footer"); strings.TrimSpace(f) != fmt.Sprint(i+1) {
			t.Errorf("footer of page %d = %q", i+1, f)
		}
	}
	// the landscape section's first page has a header of its own and no
	// footer (it defines none for its first page)
	p3 := m.Views[0].Pages[2]
	if p3.W != 841.9 || p3.H != 595.3 || pageText(t, r, 0, 2, "header") != "The first page of the landscape section" ||
		slices.ContainsFunc(p3.Layers, func(l bdf.Layer) bool { return l.Role == bdf.RoleFooter }) {
		t.Errorf("page 3: %v×%v %+v", p3.W, p3.H, p3.Layers)
	}
	body := pageText(t, r, 0, 0, "body")
	for _, want := range []string{
		"1. Text and paragraphs 1", "3. Pages and columns 2", // PAGEREF computed, not the cached 9
		"Runs can be bold, italic,", "第１条 目的を定める。", "Fresh from the orchard", "Footnotes go at the bottom of the page.",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("page 1 lacks %q in\n%s", want, body)
		}
	}
	// the footnote is below the text, the header above it
	rs := runs(t, r, 0, 0, "body")
	note, _ := find(rs, "Footnotes go")
	inline, _ := find(rs, "sits on the baseline")
	if !(note.Y > inline.Y && note.Y < 841.9-72) {
		t.Errorf("footnote at %v, last line at %v", note.Y, inline.Y)
	}
	// the scroll view has everything, the notes at the end
	var all []string
	for i := range m.Views[1].Pages {
		all = append(all, pageText(t, r, 1, i, ""))
	}
	scroll := strings.Join(all, "\n")
	for _, want := range []string{"Text and paragraphs", "After a column break.", "Back to one column", "Footnotes go at the bottom"} {
		if !strings.Contains(scroll, want) {
			t.Errorf("scroll view lacks %q", want)
		}
	}
	if strings.Index(scroll, "Footnotes go") < strings.Index(scroll, "Back to one column") {
		t.Error("footnote before the end in the scroll view")
	}
	if strings.Contains(scroll, "Page 1 of") {
		t.Error("header in the scroll view")
	}
}

func TestStructure(t *testing.T) {
	_, r := convert(t, "basic.docx", testOptions())
	pg := r.Manifest.Views[0].Pages[0]
	body := pg.Layers[len(pg.Layers)-1].Obj
	o, err := r.Object(body)
	if err != nil {
		t.Fatal(err)
	}
	counts := map[byte]int{}
	var cells, headings, figures []string
	o.Walk(func(in bdf.Instr) {
		if in.Op != bdf.OpMark {
			return
		}
		kind, payload := byte(in.Args[0].(uint64)), in.Args[1].(string)
		counts[kind]++
		switch kind {
		case bdf.MarkCell:
			cells = append(cells, payload)
		case bdf.MarkHeading:
			headings = append(headings, payload)
		case bdf.MarkFigure:
			figures = append(figures, payload)
		}
	})
	if !slices.Equal(headings, []string{"1", "1"}) { // on the first page
		t.Errorf("headings = %v", headings)
	}
	// bullets (2 levels), numbers (2 levels), articles: 5 lists, 10 items
	if counts[bdf.MarkList] != 5 || counts[bdf.MarkListItem] != 10 {
		t.Errorf("%d lists, %d items", counts[bdf.MarkList], counts[bdf.MarkListItem])
	}
	if counts[bdf.MarkTable] != 1 || !slices.Equal(cells, []string{"A1 col", "B1 col", "C1 col", "A2", "B2", "C2:C3", "A3", "B3", "A4:B4", "C4"}) {
		t.Errorf("tables %d, cells %v", counts[bdf.MarkTable], cells)
	}
	if !slices.Equal(figures, []string{"A gradient photo"}) {
		t.Errorf("figures = %v", figures)
	}
	// every structure is closed
	if counts[bdf.MarkEnd] != counts[bdf.MarkList]+counts[bdf.MarkTable]+counts[bdf.MarkFigure] {
		t.Errorf("%d ENDs", counts[bdf.MarkEnd])
	}
	// links: the contents to the pages of the headings, the external one
	var links []string
	o.Walk(func(in bdf.Instr) {
		if in.Op == bdf.OpLink {
			links = append(links, in.Args[4].(string))
		}
	})
	if !slices.Equal(links, []string{"#page=1", "#page=1", "#page=2", "https://example.com/"}) {
		t.Errorf("links = %v", links)
	}
}

// TestScrollStrips checks that the strips of the scroll view stack into
// the whole document without cutting lines.
func TestScrollStrips(t *testing.T) {
	res, r := convert(t, "grid.docx", testOptions())
	v := r.Manifest.Views[1]
	if res.Strips < 2 || len(v.Pages) != res.Strips {
		t.Fatalf("%d strips", res.Strips)
	}
	w := v.Pages[0].W
	for i, p := range v.Pages {
		if p.W != w || p.Body != nil || len(p.Layers) != 1 || p.Layers[0].Role != bdf.RoleBody {
			t.Errorf("strip %d: %+v", i, p)
		}
		for _, run := range runs(t, r, 1, i, "") {
			if run.Y-run.Size < -1 || run.Y > p.H+1 {
				t.Errorf("strip %d (height %v): run %q at %v", i, p.H, run.Text, run.Y)
			}
		}
	}
	// one line of text per run: the same text as the pages
	var pages, strips []string
	for i := range r.Manifest.Views[0].Pages {
		pages = append(pages, pageText(t, r, 0, i, "body"))
	}
	for i := range v.Pages {
		strips = append(strips, pageText(t, r, 1, i, ""))
	}
	norm := func(s []string) string { return strings.Join(strings.Fields(strings.Join(s, " ")), "") }
	// the page view repeats the table's header row on each page
	if n, m := strings.Count(norm(pages), "項目説明"), strings.Count(norm(strips), "項目説明"); n < 2 || m != 1 {
		t.Errorf("header rows: %d on pages, %d in the scroll view", n, m)
	}
}

// TestGrid checks the line pitch on a document grid: single lines take one
// grid line, 1.5 lines take one and a half, exact spacing ignores the grid.
func TestGrid(t *testing.T) {
	_, r := convert(t, "grid.docx", testOptions())
	rs := runs(t, r, 0, 0, "body")
	pitches := func(prefix string, skip int) []float64 {
		var ys []float64
		for _, run := range rs {
			if len(ys) == 0 && !strings.HasPrefix(run.Text, prefix) {
				continue
			}
			if len(ys) == 0 && skip > 0 {
				skip--
				continue
			}
			if len(ys) == 0 || run.Y != float32(ys[len(ys)-1]) {
				ys = append(ys, float64(run.Y))
			}
			if len(ys) == 3 {
				break
			}
		}
		return []float64{ys[1] - ys[0], ys[2] - ys[1]}
	}
	check := func(name string, got []float64, want float64) {
		for _, g := range got {
			if math.Abs(g-want) > 0.05 {
				t.Errorf("%s: line pitch %v, want %v", name, got, want)
				return
			}
		}
	}
	check("single", pitches("文書の行は", 0), 18)
	check("1.5 lines", pitches("文書の行は", 1), 27)
	check("exact", pitches("文書の行は", 2), 15)
}

func TestKeepWithNext(t *testing.T) {
	var b strings.Builder
	// 300 pt of body per page: 15 lines of 18 pt, then the heading as the
	// 16th, and no room for the paragraph after it
	for i := range 15 {
		b.WriteString(p(fmt.Sprintf("Filler paragraph %d", i+1), `<w:spacing w:after="0" w:line="360" w:lineRule="exact"/>`))
	}
	b.WriteString(p("Heading kept with the next", `<w:keepNext/><w:spacing w:after="0" w:line="360" w:lineRule="exact"/>`))
	b.WriteString(p("The paragraph after the heading", `<w:spacing w:after="0" w:line="360" w:lineRule="exact"/>`))
	_, r := makeDocx(t, b.String()+smallPage, nil)
	if n := len(r.Manifest.Views[0].Pages); n != 2 {
		t.Fatalf("%d pages", n)
	}
	p2 := pageText(t, r, 0, 1, "body")
	if !strings.HasPrefix(p2, "Heading kept with the next") {
		t.Errorf("page 2 = %q", p2)
	}
	// without keepNext the heading stays on the first page
	_, r = makeDocx(t, strings.ReplaceAll(b.String(), "<w:keepNext/>", "")+smallPage, nil)
	if p2 := pageText(t, r, 0, 1, "body"); !strings.HasPrefix(p2, "The paragraph after") {
		t.Errorf("page 2 without keepNext = %q", p2)
	}
}

func TestWidowControl(t *testing.T) {
	exact := `<w:spacing w:after="0" w:line="360" w:lineRule="exact"/>`
	// 300 pt of body per page, 16 lines of 18 pt; lines of 20 characters of 10 pt
	lines := func(fillers, chars int, pPr string) int {
		var b strings.Builder
		for i := range fillers {
			b.WriteString(p(fmt.Sprintf("Line %d", i+1), exact))
		}
		b.WriteString(p(strings.Repeat("あ", chars), exact, pPr))
		_, r := makeDocx(t, b.String()+smallPage, nil)
		if len(r.Manifest.Views[0].Pages) < 2 {
			return 0
		}
		ys := map[float32]bool{}
		for _, run := range runs(t, r, 0, 1, "body") {
			ys[run.Y] = true
		}
		return len(ys)
	}
	for _, c := range []struct {
		name            string
		fillers, chars  int
		without, withWC int
	}{
		{"widow", 13, 70, 1, 2},  // 3 + 1 lines: one more line goes along
		{"orphan", 15, 50, 2, 3}, // 1 + 2 lines: the first line goes along
		{"three", 14, 50, 1, 3},  // 2 + 1 lines: the paragraph moves
	} {
		if got := lines(c.fillers, c.chars, ""); got != c.without {
			t.Errorf("%s: %d lines on page 2 without widow control, want %d", c.name, got, c.without)
		}
		if got := lines(c.fillers, c.chars, "<w:widowControl/>"); got != c.withWC {
			t.Errorf("%s: %d lines on page 2 with widow control, want %d", c.name, got, c.withWC)
		}
	}
}

func TestFieldsAndSections(t *testing.T) {
	fld := func(instr, cached string) string {
		return `<w:r><w:fldChar w:fldCharType="begin"/></w:r><w:r><w:instrText> ` + instr + ` </w:instrText></w:r>` +
			`<w:r><w:fldChar w:fldCharType="separate"/></w:r><w:r><w:t>` + cached + `</w:t></w:r><w:r><w:fldChar w:fldCharType="end"/></w:r>`
	}
	footer := `<?xml version="1.0"?><w:ftr ` + ns + `><w:p><w:r><w:t xml:space="preserve">p. </w:t></w:r>` + fld(`PAGE \* ROMAN`, "I") +
		`<w:r><w:t xml:space="preserve"> / </w:t></w:r>` + fld("SECTIONPAGES", "9") + `</w:p></w:ftr>`
	sect1 := `<w:sectPr><w:footerReference w:type="default" r:id="rIdfooter1"/><w:pgSz w:w="6000" w:h="8000"/>` +
		`<w:pgMar w:top="1000" w:bottom="1000" w:left="1000" w:right="1000" w:header="500" w:footer="500"/><w:pgNumType w:fmt="lowerRoman"/></w:sectPr>`
	sect2 := `<w:sectPr><w:footerReference w:type="default" r:id="rIdfooter1"/><w:type w:val="evenPage"/><w:pgSz w:w="6000" w:h="8000"/>` +
		`<w:pgMar w:top="1000" w:bottom="1000" w:left="1000" w:right="1000" w:header="500" w:footer="500"/></w:sectPr>`
	body := p("Front matter") + `<w:p><w:r><w:br w:type="page"/></w:r><w:r><w:t>Second front page</w:t></w:r></w:p>` +
		`<w:p><w:pPr>` + sect1 + `</w:pPr></w:p>` + p("Chapter one") + sect2
	_, r := makeDocx(t, body, map[string]string{"word/footer1.xml": footer})
	v := r.Manifest.Views[0]
	if len(v.Pages) != 4 {
		t.Fatalf("%d pages", len(v.Pages))
	}
	var footers []string
	for i := range v.Pages {
		footers = append(footers, pageText(t, r, 0, i, "footer"))
	}
	// the section starts on an even page: a blank page is inserted (with no
	// footer); the numbering continues
	want := []string{"p. I / 2", "p. II / 2", "", "p. IV / 1"}
	if !slices.Equal(footers, want) {
		t.Errorf("footers = %q, want %q", footers, want)
	}
}

func TestOptions(t *testing.T) {
	opts := testOptions()
	opts.Views = ViewsPages
	res, r := convert(t, "basic.docx", opts)
	if len(r.Manifest.Views) != 1 || res.Strips != 0 {
		t.Errorf("pages only: %d views", len(r.Manifest.Views))
	}
	opts = testOptions()
	opts.Views = ViewsScroll
	_, r = convert(t, "basic.docx", opts)
	if len(r.Manifest.Views) != 1 || r.Manifest.Views[0].Kind != bdf.ViewScroll {
		t.Errorf("scroll only: %+v", r.Manifest.Views)
	}
	opts = testOptions()
	opts.Pages = []int{2}
	_, r = convert(t, "basic.docx", opts)
	if n := len(r.Manifest.Views[0].Pages); n != 1 || !strings.Contains(pageText(t, r, 0, 0, "header"), "Page 2 of 3") {
		t.Errorf("page 2: %d pages", n)
	}
	opts = testOptions()
	opts.Views = "none"
	if _, err := ConvertFile("testdata/basic.docx", opts); err == nil {
		t.Error("bad views accepted")
	}
	opts = testOptions()
	opts.SystemFonts = true
	res, r = convert(t, "basic.docx", opts)
	if res.EmbeddedFonts != 0 {
		t.Errorf("SystemFonts embedded %d fonts", res.EmbeddedFonts)
	}
	for _, p := range r.Manifest.Parts {
		if p.T == bdf.PartFont {
			t.Fatal("font part with SystemFonts")
		}
	}
}

func TestDeterministic(t *testing.T) {
	var outs [2][]byte
	for i := range outs {
		res, err := ConvertFile("testdata/basic.docx", testOptions())
		if err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		if err := res.Doc.WriteSingle(&buf); err != nil {
			t.Fatal(err)
		}
		outs[i] = buf.Bytes()
	}
	if !bytes.Equal(outs[0], outs[1]) {
		t.Fatal("conversion is not deterministic")
	}
}

func TestRegistered(t *testing.T) {
	f, err := conv.DetectFile("testdata/basic.docx")
	if err != nil || f == nil || f.Name != "docx" {
		t.Fatalf("detected %v, %v", f, err)
	}
	for _, name := range []string{"../pptx/testdata/basic.pptx", "../pdf/testdata/cairo-cff.pdf"} {
		if f, _ := conv.DetectFile(name); f != nil && f.Name == "docx" {
			t.Errorf("%s detected as docx", name)
		}
	}
}

func TestLineBreaking(t *testing.T) {
	// hanging punctuation: a full stop past the right edge stays on its line
	body := p("あいうえおかきくけこ。次の行", `<w:jc w:val="left"/>`) +
		`<w:sectPr><w:pgSz w:w="4400" w:h="8000"/><w:pgMar w:top="1000" w:bottom="1000" w:left="1000" w:right="1000" w:header="500" w:footer="500"/></w:sectPr>`
	_, r := makeDocx(t, body, map[string]string{"word/styles.xml": `<?xml version="1.0"?><w:styles ` + ns +
		`><w:docDefaults><w:rPrDefault><w:rPr><w:sz w:val="24"/></w:rPr></w:rPrDefault></w:docDefaults></w:styles>`})
	// 120 pt of text width holds 10 characters of 12 pt, and 。 hangs
	got := pageText(t, r, 0, 0, "body")
	if got != "あいうえおかきくけこ。次の行" {
		t.Errorf("lines = %q", got)
	}
	rs := runs(t, r, 0, 0, "body")
	if len(rs) != 2 || rs[0].Y == rs[1].Y || !strings.HasSuffix(rs[0].Text, "こ。") || rs[1].Sep != bdf.SepNone {
		t.Errorf("runs = %+v", rs)
	}
}

// TestCharacterGrid checks the character grid (docGrid linesAndChars): East
// Asian characters take whole cells of the default font size plus the
// grid's extra pitch, which sets how many fit on a line.
func TestCharacterGrid(t *testing.T) {
	styles := map[string]string{"word/styles.xml": `<?xml version="1.0"?><w:styles ` + ns +
		`><w:docDefaults><w:rPrDefault><w:rPr><w:sz w:val="20"/></w:rPr></w:rPrDefault></w:docDefaults></w:styles>`}
	firstLine := func(grid string) string {
		body := p(strings.Repeat("あ", 30), `<w:jc w:val="left"/>`) + `<w:sectPr><w:pgSz w:w="6000" w:h="8000"/>` +
			`<w:pgMar w:top="1000" w:bottom="1000" w:left="1000" w:right="1000" w:header="500" w:footer="500"/>` + grid + `</w:sectPr>`
		_, r := makeDocx(t, body, styles)
		rs := runs(t, r, 0, 0, "body")
		var b strings.Builder
		for _, run := range rs {
			if run.Y != rs[0].Y {
				break
			}
			b.WriteString(run.Text)
		}
		return b.String()
	}
	// 200 pt of text width: 20 characters of 10 pt, 19 cells of 10.5 pt
	if n := len([]rune(firstLine(""))); n != 20 {
		t.Errorf("without a grid: %d characters on the first line", n)
	}
	if n := len([]rune(firstLine(`<w:docGrid w:type="linesAndChars" w:linePitch="360" w:charSpace="2048"/>`))); n != 19 {
		t.Errorf("with a character grid: %d characters on the first line", n)
	}
	// a line grid alone leaves the characters alone
	if n := len([]rune(firstLine(`<w:docGrid w:type="lines" w:linePitch="360" w:charSpace="2048"/>`))); n != 20 {
		t.Errorf("with a line grid: %d characters on the first line", n)
	}
}

// TestVertical checks a section of East Asian vertical text: lines run down
// the page from right to left, one ALT_TEXT run per line in extraction;
// headers stay horizontal, and the scroll view lays the text out
// horizontally.
func TestVertical(t *testing.T) {
	res, r := convert(t, "vertical.docx", testOptions())
	if len(res.Warnings) > 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	body := runs(t, r, 0, 0, "body")
	var lines []bdf.TextRun
	for _, run := range body {
		if run.AltText && strings.Contains(run.Text, "縦書きでは") {
			lines = append(lines, run)
		}
	}
	if len(lines) < 3 {
		t.Fatalf("%d vertical lines: %+v", len(lines), body)
	}
	for i, l := range lines {
		// turned by 90°: the baseline runs down the page
		if math.Abs(float64(l.Matrix[0])) > 1e-3 || math.Abs(float64(l.Matrix[1]-1)) > 1e-3 {
			t.Errorf("line %d matrix %v", i, l.Matrix)
		}
		if i > 0 && l.X >= lines[i-1].X {
			t.Errorf("line %d at x %v, right of the one before at %v", i, l.X, lines[i-1].X)
		}
	}
	text := text(body)
	for _, want := range []string{"縦書きの文書", "縦書きでは、行は上から下へ進み、右から左へ並びます。", "英字のWordや数字の2026は、横に倒して組みます。",
		"第１条 目的を定める。", "2026年10月26日、100個の例。", "表も縦に組まれます。", "脚注も縦書きです。"} {
		if !strings.Contains(text, want) {
			t.Errorf("page 1 lacks %q in\n%s", want, text)
		}
	}
	// the footnote is at the left end of the body, the header horizontal
	note, _ := find(body, "脚注も縦書き")
	if note.X > 150 {
		t.Errorf("footnote at x %v", note.X)
	}
	head := runs(t, r, 0, 0, "header")
	if len(head) == 0 || head[0].Matrix[0] != 1 || head[0].AltText {
		t.Errorf("header runs %+v", head)
	}
	// structure: heading, list, table, picture
	o, err := r.Object(r.Manifest.Views[0].Pages[0].Layers[2].Obj)
	if err != nil {
		t.Fatal(err)
	}
	counts := map[byte]int{}
	var figures []string
	o.Walk(func(in bdf.Instr) {
		if in.Op == bdf.OpMark {
			kind := byte(in.Args[0].(uint64))
			counts[kind]++
			if kind == bdf.MarkFigure {
				figures = append(figures, in.Args[1].(string))
			}
		}
	})
	if counts[bdf.MarkHeading] != 1 || counts[bdf.MarkList] != 1 || counts[bdf.MarkTable] != 1 || counts[bdf.MarkCell] != 4 ||
		!slices.Equal(figures, []string{"A small picture"}) {
		t.Errorf("marks %v, figures %v", counts, figures)
	}
	// the scroll view is horizontal
	for _, run := range runs(t, r, 1, 0, "") {
		if run.AltText || run.Matrix[0] != 1 {
			t.Fatalf("scroll view run %+v", run)
		}
	}
}
