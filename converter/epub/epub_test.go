package epub

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
	_ "github.com/shibukawa/bdf/converter/docx"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/japanese"
)

// testOptions lays text out with the Word converter's test fonts only, so
// that the output does not depend on the machine.
func testOptions() *Options {
	return &Options{FontDirs: []string{"../docx/testdata/fonts"}, NoSystemFonts: true}
}

func convertFile(t *testing.T, name string, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
	if opts == nil {
		opts = testOptions()
	}
	res, err := ConvertFile("testdata/"+name, opts)
	if err != nil {
		t.Fatal(err)
	}
	return res, reopen(t, res)
}

func convertBytes(t *testing.T, data []byte, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
	if opts == nil {
		opts = testOptions()
	}
	res, err := Convert(bytes.NewReader(data), int64(len(data)), opts)
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

// page is what a page of a view holds.
type page struct {
	text   string
	roles  []string
	links  []string
	images [][4]float32 // x, y, w, h (in the coordinates they are drawn in)
	marks  map[byte][]string
	w, h   float32
}

func viewPages(t *testing.T, r *bdf.Reader, view string) []page {
	t.Helper()
	var v *bdf.View
	for _, x := range r.Manifest.Views {
		if x.ID == view {
			v = x
		}
	}
	if v == nil {
		t.Fatalf("no view %q", view)
	}
	var out []page
	for _, pg := range v.Pages {
		p := page{marks: map[byte][]string{}, w: pg.W, h: pg.H}
		var b strings.Builder
		for _, l := range pg.Layers {
			p.roles = append(p.roles, l.Role)
			o, err := r.Object(l.Obj)
			if err != nil {
				t.Fatal(err)
			}
			rs, err := bdf.ExtractText(o, func(h bdf.Hash) *bdf.ObjectPart {
				ch, _ := r.Object(h)
				return ch
			})
			if err != nil {
				t.Fatal(err)
			}
			for i, run := range rs {
				if i > 0 || b.Len() > 0 {
					switch run.Sep {
					case bdf.SepSpace:
						b.WriteByte(' ')
					case bdf.SepBreak:
						b.WriteByte('\n')
					}
				}
				b.WriteString(run.Text)
			}
			var walk func(o *bdf.ObjectPart)
			walk = func(o *bdf.ObjectPart) {
				o.Walk(func(in bdf.Instr) {
					switch in.Op {
					case bdf.OpMark:
						kind := byte(in.Args[0].(uint64))
						p.marks[kind] = append(p.marks[kind], in.Args[1].(string))
					case bdf.OpLink:
						p.links = append(p.links, in.Args[4].(string))
					case bdf.OpImage:
						p.images = append(p.images, [4]float32{in.Args[1].(float32), in.Args[2].(float32), in.Args[3].(float32), in.Args[4].(float32)})
					}
				})
				for _, h := range o.Objects {
					if ch, err := r.Object(h); err == nil {
						walk(ch)
					}
				}
			}
			walk(o)
		}
		p.text = b.String()
		out = append(out, p)
	}
	return out
}

// pageOf returns the 1-based number of the first page whose text holds s.
func pageOf(pages []page, s string) int {
	for i, p := range pages {
		if strings.Contains(p.text, s) {
			return i + 1
		}
	}
	return 0
}

func TestBasic(t *testing.T) {
	res, r := convertFile(t, "basic.epub", nil)
	if len(res.Warnings) > 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	m := r.Manifest
	dc := m.Meta.DC
	if m.Meta.Source != "epub" || !slices.Equal(dc.Title, bdf.DCValues{"A Small Book", "Tests for the EPUB converter"}) ||
		dc.Creator.First() != "Test Author" || dc.Language.First() != "en" || dc.Publisher.First() != "bdf" ||
		dc.Identifier.First() != "urn:uuid:5b0f3c1e-8f5d-4b8e-9a51-6c1d0a7e2f10" || dc.Date.First() != "2026-09-27" ||
		dc.Modified.First() != "2026-09-27T00:00:00Z" || !slices.Equal(dc.Subject, bdf.DCValues{"EPUB", "Testing"}) ||
		dc.Description.First() != "A book that tests the EPUB converter." {
		t.Errorf("metadata %+v", dc)
	}
	if len(m.Views) != 1 || m.Views[0].ID != "pages" || m.Views[0].Kind != bdf.ViewFlow || m.Views[0].Direction != "" {
		t.Fatalf("views %+v", m.Views[0])
	}
	pages := viewPages(t, r, "pages")
	if res.Pages != len(pages) || res.Chapters != 5 || res.Images != 2 || res.Vertical || res.FixedLayout {
		t.Errorf("result %+v", res)
	}
	if p := pages[0]; p.w != float32(A5.Width) || p.h != float32(A5.Height) {
		t.Errorf("page size %v × %v", p.w, p.h)
	}
	// the cover: its picture fills the page, without a page number
	cover := pages[0]
	if len(cover.images) != 1 || cover.images[0][2] != float32(A5.Width) || slices.Contains(cover.roles, bdf.RoleFooter) {
		t.Errorf("cover %+v", cover)
	}
	// reading order: the spine's linear items, then the others
	order := []int{pageOf(pages, "Contents"), pageOf(pages, "This book tests"), pageOf(pages, "A centered line"), pageOf(pages, "A note outside")}
	if !slices.IsSorted(order) || order[0] != 2 || slices.Contains(order, 0) {
		t.Errorf("chapters start on pages %v", order)
	}
	// chapters start pages: the last line of one is not on the page of the next
	if pageOf(pages, "A centered line") == pageOf(pages, "A footnote") {
		t.Errorf("chapter two does not start a page")
	}
	// the page number of a text page, as a footer
	if p := pages[1]; !slices.Contains(p.roles, bdf.RoleFooter) || !strings.HasPrefix(p.text, "2\n") {
		t.Errorf("folio of page 2: %q %v", p.text, p.roles)
	}
	all := ""
	var links []string
	for _, p := range pages {
		all += p.text + "\n"
		links = append(links, p.links...)
	}
	for _, s := range []string{"A page marker written as an empty element does not swallow the rest of the paragraph",
		"Lists keep their bullets.", "Figure 1. A picture of the publication.", "fmt.Println(\"hello, epub\")"} {
		if !strings.Contains(all, s) {
			t.Errorf("missing %q in\n%s", s, all)
		}
	}
	if strings.Contains(all, "hidden by a class") {
		t.Errorf("the hidden paragraph is drawn")
	}
	// links to chapters and to elements of other chapters go to their pages
	want := []string{
		fmt.Sprintf("#page=%d", pageOf(pages, "This book tests")), // the table of contents
		fmt.Sprintf("#page=%d", pageOf(pages, "A centered line")), // ch2.xhtml#sec2
		fmt.Sprintf("#page=%d", pageOf(pages, "A note outside")),  // notes.xhtml#n1 (not linear)
		"https://example.com/",
	}
	for _, w := range want {
		if !slices.Contains(links, w) {
			t.Errorf("no link %s in %v", w, links)
		}
	}
	// the figure's alternative text and the heading structure
	var figures, headings []string
	for _, p := range pages {
		figures = append(figures, p.marks[bdf.MarkFigure]...)
		headings = append(headings, p.marks[bdf.MarkHeading]...)
	}
	if !slices.Contains(figures, "Three bars rising from left to right") || !slices.Contains(figures, "A Small Book") {
		t.Errorf("figures %q", figures)
	}
	if !slices.Contains(headings, "1") || !slices.Contains(headings, "2") {
		t.Errorf("headings %q", headings)
	}
}

func TestVertical(t *testing.T) {
	res, r := convertFile(t, "vertical.epub", nil)
	if len(res.Warnings) > 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	v := r.Manifest.Views[0]
	if len(r.Manifest.Views) != 1 || v.Direction != bdf.DirectionRTL || !res.Vertical {
		t.Fatalf("views %+v, result %+v", v, res)
	}
	pages := viewPages(t, r, "pages")
	p1 := pages[pageOf(pages, "縦書きの本を読む")-1]
	// the lines are drawn one by one as child objects with their text as
	// ALT_TEXT (vertical text), tate-chu-yoko within the line's text
	alts := strings.Join(p1.marks[bdf.MarkAltText], "\n")
	for _, s := range []string{"縦書きの本を読むとき、行は上から下へ、右から左へ進みます。", "12月31日", "傍点を付けた"} {
		if !strings.Contains(alts, s) {
			t.Errorf("no vertical line holding %q in\n%s", s, alts)
		}
	}
	if strings.Contains(p1.text, "かんじ") || !strings.Contains(p1.text, "漢字のように") {
		t.Errorf("ruby: %q", p1.text)
	}
	// the colophon is horizontal (class hltr): no lines as child objects
	col := pages[pageOf(pages, "書名")-1]
	if len(col.marks[bdf.MarkAltText]) != 0 {
		t.Errorf("the colophon is vertical: %q", col.marks[bdf.MarkAltText])
	}
	// the gaiji picture is one em of the text
	var gaiji [4]float32
	for _, im := range p1.images {
		if im[2] < 20 {
			gaiji = im
		}
	}
	if gaiji[2] != 12 || gaiji[3] != 12 {
		t.Errorf("gaiji picture %v (images %v)", gaiji, p1.images)
	}
	// both views: the scroll view is horizontal
	res, r = convertFile(t, "vertical.epub", &Options{Views: ViewsBoth, FontDirs: testOptions().FontDirs, NoSystemFonts: true})
	if len(r.Manifest.Views) != 2 || r.Manifest.Views[1].Kind != bdf.ViewScroll || r.Manifest.Views[1].Direction != "" || res.Strips == 0 {
		t.Fatalf("views %+v", r.Manifest.Views)
	}
	scroll := viewPages(t, r, "scroll")
	if s := scroll[0].text; !strings.Contains(s, "縦書きの本を読むとき") {
		t.Errorf("scroll view %q", s)
	}
}

func TestFixed(t *testing.T) {
	res, r := convertFile(t, "fixed.epub", nil)
	if !res.FixedLayout || res.Pages != 4 || res.Images != 4 {
		t.Errorf("result %+v", res)
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "outside the reading order") {
		t.Errorf("warnings %q", res.Warnings)
	}
	v := r.Manifest.Views[0]
	if len(r.Manifest.Views) != 1 || v.Kind != bdf.ViewFixed || v.Direction != bdf.DirectionRTL || r.Manifest.Meta.DC.Title.First() != "固定レイアウトのテスト" {
		t.Fatalf("view %+v", v)
	}
	for i, p := range viewPages(t, r, "pages") {
		if p.w != 450 || p.h != 600 || len(p.images) != 1 || p.images[0] != [4]float32{0, 0, 450, 600} {
			t.Errorf("page %d: %v × %v, images %v", i+1, p.w, p.h, p.images)
		}
	}
	// the pages of the same picture share it; selected pages
	res, _ = convertFile(t, "fixed.epub", &Options{Pages: conv.Pages{{From: 2, To: 3}}})
	if res.Pages != 2 {
		t.Errorf("selected %d pages", res.Pages)
	}
	if _, err := ConvertFile("testdata/fixed.epub", &Options{Pages: conv.Pages{{From: 9, To: 9}}}); err == nil {
		t.Errorf("page 9 converted")
	}
}

// file is a file of a publication built by a test.
type file struct{ name, body string }

// makeEPUB builds a publication with a package document at OEBPS/content.opf.
func makeEPUB(t *testing.T, opf string, files ...file) []byte {
	t.Helper()
	var buf bytes.Buffer
	z := zip.NewWriter(&buf)
	w, _ := z.CreateHeader(&zip.FileHeader{Name: "mimetype", Method: zip.Store})
	w.Write([]byte("application/epub+zip"))
	all := append([]file{
		{"META-INF/container.xml", `<?xml version="1.0"?><container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
<rootfiles><rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/></rootfiles></container>`},
		{"OEBPS/content.opf", opf},
	}, files...)
	for _, f := range all {
		w, err := z.Create(f.name)
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(f.body))
	}
	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// opf3 writes a package document with the chapters in the spine.
func opf3(meta, spineAttrs string, chapters ...string) string {
	var man, spine strings.Builder
	for i, c := range chapters {
		fmt.Fprintf(&man, `<item id="c%d" href="%s" media-type="application/xhtml+xml"/>`, i, c)
		fmt.Fprintf(&spine, `<itemref idref="c%d"/>`, i)
	}
	return `<?xml version="1.0" encoding="UTF-8"?>
<package xmlns="http://www.idpf.org/2007/opf" version="3.0" unique-identifier="id">
<metadata xmlns:dc="http://purl.org/dc/elements/1.1/"><dc:identifier id="id">x</dc:identifier><dc:title>T</dc:title>` + meta + `</metadata>
<manifest>` + man.String() + `</manifest><spine` + spineAttrs + `>` + spine.String() + `</spine></package>`
}

func chapterDoc(lang, htmlAttrs, head, body string) string {
	return `<?xml version="1.0" encoding="UTF-8"?><html xmlns="http://www.w3.org/1999/xhtml" xml:lang="` + lang + `"` + htmlAttrs +
		`><head><title>x</title>` + head + `</head><body>` + body + `</body></html>`
}

// Ids repeat across chapters: links go to the element of the chapter they
// name.
func TestLinks(t *testing.T) {
	filler := strings.Repeat("<p>Some text to fill the page with lines of words.</p>", 40)
	data := makeEPUB(t, opf3("<dc:language>en</dc:language>", "", "a.xhtml", "sub/b.xhtml"),
		file{"OEBPS/a.xhtml", chapterDoc("en", "", "<style>#gone { display: none }</style>", `<p id="gone">Gone.</p><p id="x">First x.</p>`+filler+
			`<p><a href="sub/b.xhtml#x">to b's x</a> <a href="#x">to a's x</a> <a href="sub/b%20c.xhtml">nowhere</a></p>`)},
		file{"OEBPS/sub/b.xhtml", chapterDoc("en", "", "", filler+`<p id="x">Second x.</p><p><a href="../a.xhtml">to a</a></p>`)},
	)
	_, r := convertBytes(t, data, nil)
	pages := viewPages(t, r, "pages")
	var links []string
	for _, p := range pages {
		links = append(links, p.links...)
	}
	want := []string{fmt.Sprintf("#page=%d", pageOf(pages, "Second x.")), "#page=1", "#page=1"}
	if !slices.Equal(links, want) {
		t.Errorf("links %v, want %v", links, want)
	}
	// id selectors match the ids as written
	if pageOf(pages, "Gone.") != 0 {
		t.Errorf("an element hidden by its id is drawn")
	}
}

func TestDRM(t *testing.T) {
	enc := func(alg string) file {
		return file{"META-INF/encryption.xml", `<?xml version="1.0"?><encryption xmlns="urn:oasis:names:tc:opendocument:xmlns:container"
xmlns:enc="http://www.w3.org/2001/04/xmlenc#"><enc:EncryptedData><enc:EncryptionMethod Algorithm="` + alg + `"/>
<enc:CipherData><enc:CipherReference URI="OEBPS/a.xhtml"/></enc:CipherData></enc:EncryptedData></encryption>`}
	}
	chapter := file{"OEBPS/a.xhtml", chapterDoc("en", "", "", "<p>text</p>")}
	data := makeEPUB(t, opf3("", "", "a.xhtml"), chapter, enc("http://www.w3.org/2001/04/xmlenc#aes128-cbc"))
	if _, err := Convert(bytes.NewReader(data), int64(len(data)), testOptions()); !errors.Is(err, ErrDRM) {
		t.Errorf("encrypted chapter: %v", err)
	}
	// fonts obfuscated as the specification says are not DRM
	data = makeEPUB(t, opf3("", "", "a.xhtml"), chapter, enc("http://www.idpf.org/2008/embedding"))
	convertBytes(t, data, nil)
}

// An EPUB 2 publication: opf: attributes, an XHTML 1.1 chapter with
// entities, a chapter that is not well-formed, one in Shift_JIS, one that
// is a picture.
func TestEPUB2(t *testing.T) {
	sjis, _ := japanese.ShiftJIS.NewEncoder().String(`<?xml version="1.0" encoding="Shift_JIS"?><html xmlns="http://www.w3.org/1999/xhtml"><head><title>s</title></head><body><p>日本語の文書</p></body></html>`)
	opf := `<?xml version="1.0"?><package xmlns="http://www.idpf.org/2007/opf" version="2.0" unique-identifier="BookId">
<metadata xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:opf="http://www.idpf.org/2007/opf">
<dc:title>Old Book</dc:title><dc:creator opf:role="aut" opf:file-as="Writer, A">A Writer</dc:creator>
<dc:identifier id="other">isbn-0</dc:identifier><dc:identifier id="BookId" opf:scheme="UUID">urn:uuid:1</dc:identifier>
<dc:date opf:event="publication">2001</dc:date><dc:date opf:event="modification">2002-03-04</dc:date>
<dc:language>en</dc:language><meta name="cover" content="img"/></metadata>
<manifest><item id="ncx" href="toc.ncx" media-type="application/x-dtbncx+xml"/>
<item id="a" href="a.html" media-type="application/xhtml+xml"/><item id="b" href="b.html" media-type="application/xhtml+xml"/>
<item id="s" href="s.html" media-type="application/xhtml+xml"/><item id="img" href="p.png" media-type="image/png"/></manifest>
<spine toc="ncx"><itemref idref="a"/><itemref idref="b"/><itemref idref="s"/><itemref idref="img"/></spine></package>`
	data := makeEPUB(t, opf,
		file{"OEBPS/a.html", `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.w3.org/1999/xhtml"><head><title>a</title></head><body><p>A&nbsp;&mdash;&nbsp;B</p></body></html>`},
		file{"OEBPS/b.html", `<html><head><title>b</title></head><body><p>Broken<br>markup</p></body></html>`},
		file{"OEBPS/s.html", sjis},
		file{"OEBPS/p.png", string(tinyPNG)},
	)
	res, r := convertBytes(t, data, nil)
	dc := r.Manifest.Meta.DC
	if dc.Title.First() != "Old Book" || dc.Creator.First() != "A Writer" || !slices.Equal(dc.Identifier, bdf.DCValues{"urn:uuid:1", "isbn-0"}) ||
		dc.Date.First() != "2001" || dc.Modified.First() != "2002-03-04" {
		t.Errorf("metadata %+v", dc)
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "b.html: not well-formed XML") {
		t.Errorf("warnings %q", res.Warnings)
	}
	pages := viewPages(t, r, "pages")
	var all string
	for _, p := range pages {
		all += p.text + "\n"
	}
	for _, s := range []string{"A — B", "Broken\nmarkup", "日本語の文書"} {
		if !strings.Contains(all, s) {
			t.Errorf("missing %q in %q", s, all)
		}
	}
	if last := pages[len(pages)-1]; len(last.images) != 1 || res.Images != 1 {
		t.Errorf("the picture in the spine: %+v", last)
	}
}

// tinyPNG is a white 4 × 4 pixel PNG.
var tinyPNG = func() []byte {
	img := image.NewGray(image.Rect(0, 0, 4, 4))
	for i := range img.Pix {
		img.Pix[i] = 255
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}()

// The writing mode of books whose style sheets do not say: vertical for
// East Asian books bound on the right, as Kindle's primary-writing-mode
// says.
func TestWritingModeDefault(t *testing.T) {
	ja := file{"OEBPS/a.xhtml", chapterDoc("ja", "", "", "<p>縦書きの本</p>")}
	for _, c := range []struct {
		meta, spine string
		vertical    bool
	}{
		{"<dc:language>ja</dc:language>", ` page-progression-direction="rtl"`, true},
		{"<dc:language>ja</dc:language>", "", false},
		{"<dc:language>en</dc:language>", ` page-progression-direction="rtl"`, false},
		{`<dc:language>ja</dc:language><meta name="primary-writing-mode" content="vertical-rl"/>`, "", true},
	} {
		res, _ := convertBytes(t, makeEPUB(t, opf3(c.meta, c.spine, "a.xhtml"), ja), nil)
		if res.Vertical != c.vertical {
			t.Errorf("%s %s: vertical %v", c.meta, c.spine, res.Vertical)
		}
	}
	// a style sheet that says horizontal-tb anywhere keeps unmarked chapters horizontal
	css := file{"OEBPS/s.css", ".h { writing-mode: horizontal-tb }"}
	link := `<link rel="stylesheet" href="s.css"/>`
	res, _ := convertBytes(t, makeEPUB(t, opf3("<dc:language>ja</dc:language>", ` page-progression-direction="rtl"`, "a.xhtml"),
		file{"OEBPS/a.xhtml", chapterDoc("ja", "", link, "<p>横書き</p>")}, css), nil)
	if res.Vertical {
		t.Errorf("vertical with a style sheet that says writing modes")
	}
	// inline style on the body
	res, _ = convertBytes(t, makeEPUB(t, opf3("<dc:language>ja</dc:language>", "", "a.xhtml"),
		file{"OEBPS/a.xhtml", strings.Replace(chapterDoc("ja", "", "", "<p>縦</p>"), "<body>", `<body style="writing-mode: vertical-rl">`, 1)}), nil)
	if !res.Vertical {
		t.Errorf("horizontal with an inline writing mode")
	}
}

func TestCSS(t *testing.T) {
	sheets := map[string]string{
		"a.css":     `@charset "UTF-8"; @import url("sub/b.css"); /* c */ .x { text-align: center; color: red }`,
		"sub/b.css": `@import 'c.css'; html.v { -epub-writing-mode: vertical-rl } span.t, .u { text-combine-upright: all }`,
		"sub/c.css": `@font-face { font-family: f; src: url(f.otf) } @media amzn-kf8 { img.g { width: 1em; height: 1em } }
			p .x { text-align: right } .x:first-child { text-align: right } [lang] { display: none } #i { display: none !important }
			.y { display: block } .z { text-align: left } .z.w { text-align: right } .z { text-align: justify }`,
	}
	var rules []cssRule
	parseCSS(sheets["a.css"], func(u string) string { return sheets[u] }, 0, &rules)
	var got []string
	for _, r := range rules {
		got = append(got, fmt.Sprintf("%s|%s|%s|%d %v", r.sel.tag, r.sel.id, strings.Join(r.sel.classes, "."), r.spec, r.decls))
	}
	want := []string{
		"img||g|101 [[width 1em] [height 1em]]",
		"|i||10000 [[display none]]",
		"||y|100 [[display block]]",
		"||z|100 [[text-align left]]",
		"||z.w|200 [[text-align right]]",
		"||z|100 [[text-align justify]]",
		"html||v|101 [[-epub-writing-mode vertical-rl]]",
		"span||t|101 [[text-combine-upright all]]",
		"||u|100 [[text-combine-upright all]]",
		"||x|100 [[text-align center]]",
	}
	if !slices.Equal(got, want) {
		t.Errorf("rules\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	el := func(tag, class, style string) *htmlNode { return newNode(tag, class, style) }
	st := newStyler(rules)
	for _, c := range []struct {
		n    *htmlNode
		want string
	}{
		{el("p", "z w", ""), "text-align: right; "},                                     // specificity
		{el("p", "z", ""), "text-align: justify; "},                                     // order
		{el("p", "z", "text-align: center"), "text-align: justify; text-align: center"}, // own style wins
		{el("img", "g", ""), "width: 1em; height: 1em; "},
		{el("span", "g", ""), ""}, // sizes are for pictures
		{el("div", "y", ""), ""},  // only display: none
		{el("html", "v", ""), "-epub-writing-mode: vertical-rl; "},
	} {
		st.apply(c.n)
		if got := attrVal(c.n, "style"); got != strings.TrimSpace(c.want) && got != c.want {
			t.Errorf("<%s class=%q>: style %q, want %q", c.n.Data, attrVal(c.n, "class"), got, c.want)
		}
	}
	if writingMode(el("html", "", "writing-mode: vertical-rl; -webkit-writing-mode: horizontal-tb")) != modeVertical {
		t.Errorf("the standard property does not win")
	}
}

func TestDetect(t *testing.T) {
	for _, name := range []string{"basic.epub", "vertical.epub", "fixed.epub"} {
		f, err := conv.DetectFile("testdata/" + name)
		if err != nil || f == nil || f.Name != "epub" {
			t.Errorf("%s detected as %v (%v)", name, f, err)
		}
	}
	// a mimetype file that is not first, compressed
	var buf bytes.Buffer
	z := zip.NewWriter(&buf)
	w, _ := z.Create("META-INF/container.xml")
	w.Write([]byte(`<container><rootfiles><rootfile full-path="a.opf" media-type="application/oebps-package+xml"/></rootfiles></container>`))
	w, _ = z.Create("mimetype")
	w.Write([]byte("application/epub+zip"))
	z.Close()
	head := buf.Bytes()[:min(1024, buf.Len())]
	if !Detect(head, bytes.NewReader(buf.Bytes()), int64(buf.Len())) {
		t.Errorf("careless EPUB not detected")
	}
	if f, _ := conv.DetectFile("../docx/testdata/basic.docx"); f == nil || f.Name != "docx" {
		t.Errorf("a Word document is detected as %v", f)
	}
	docx, _ := os.ReadFile("../docx/testdata/basic.docx")
	if Detect(docx[:1024], bytes.NewReader(docx), int64(len(docx))) {
		t.Errorf("a Word document is an EPUB")
	}
}

func TestRegistered(t *testing.T) {
	f := conv.Lookup("epub")
	if f == nil || !slices.Contains(f.Extensions, ".epub") {
		t.Fatalf("format %+v", f)
	}
	res, err := conv.ConvertFile("testdata/basic.epub", "", &conv.Options{FontDirs: []string{"../docx/testdata/fonts"}, NoSystemFonts: true,
		Params: map[string]string{"paper": "b6", "views": "both", "size": "10"}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(res.Summary, "") || len(res.Doc.Views) != 2 || res.Doc.Views[0].Pages[0].W != float32(B6.Width) {
		t.Errorf("summary %q, views %d", res.Summary, len(res.Doc.Views))
	}
	for _, p := range []map[string]string{{"paper": "b99"}, {"size": "-1"}, {"views": "all"}} {
		if _, err := conv.ConvertFile("testdata/basic.epub", "", &conv.Options{Params: p}); err == nil {
			t.Errorf("params %v accepted", p)
		}
	}
}

func TestPaper(t *testing.T) {
	for in, want := range map[string]Paper{"A5": A5, "b6": B6, "letter": Letter, "100x200": {283.46, 566.93}, "128 x 188 mm": {362.83, 532.91}} {
		if got, err := ParsePaper(in); err != nil || math.Abs(got.Width-want.Width) > 0.01 || math.Abs(got.Height-want.Height) > 0.01 {
			t.Errorf("ParsePaper(%q) = %v, %v", in, got, err)
		}
	}
	for _, in := range []string{"", "x", "10x10", "a3"} {
		if _, err := ParsePaper(in); err == nil {
			t.Errorf("ParsePaper(%q) accepted", in)
		}
	}
}

func TestDeterministic(t *testing.T) {
	var outs [2][]byte
	for i := range outs {
		res, err := ConvertFile("testdata/vertical.epub", testOptions())
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
		t.Errorf("two conversions differ")
	}
}

type htmlNode = html.Node

func newNode(tag, class, style string) *html.Node {
	n := &html.Node{Type: html.ElementNode, Data: tag}
	if class != "" {
		n.Attr = append(n.Attr, html.Attribute{Key: "class", Val: class})
	}
	if style != "" {
		n.Attr = append(n.Attr, html.Attribute{Key: "style", Val: style})
	}
	return n
}
