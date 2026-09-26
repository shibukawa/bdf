package html

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
	"golang.org/x/text/encoding/japanese"
)

// testOptions lays text out with the Word converter's test fonts only, so
// that the output does not depend on the machine, and fetches nothing.
func testOptions() *Options {
	return &Options{FontDirs: []string{"../docx/testdata/fonts"}, NoSystemFonts: true, NoRemote: true}
}

func convertHTML(t *testing.T, src string, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
	res, err := ConvertBytes([]byte(src), opts)
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

// content is what the objects of a view hold: the text with breaks as
// newlines, the structure marks, the links and the images.
type content struct {
	text    string
	marks   map[byte][]string
	links   []string
	images  [][4]float32 // x, y, w, h
	fonts   []bdf.Font
	strips  int
	widthPt float32
}

func viewContent(t *testing.T, r *bdf.Reader, view int) *content {
	t.Helper()
	c := &content{marks: map[byte][]string{}}
	v := r.Manifest.Views[view]
	c.strips = len(v.Pages)
	if len(v.Pages) > 0 {
		c.widthPt = v.Pages[0].W
	}
	var b strings.Builder
	for _, pg := range v.Pages {
		for _, l := range pg.Layers {
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
			o.Walk(func(in bdf.Instr) {
				switch in.Op {
				case bdf.OpMark:
					kind := byte(in.Args[0].(uint64))
					c.marks[kind] = append(c.marks[kind], in.Args[1].(string))
				case bdf.OpLink:
					c.links = append(c.links, in.Args[4].(string))
				case bdf.OpImage:
					c.images = append(c.images, [4]float32{in.Args[1].(float32), in.Args[2].(float32), in.Args[3].(float32), in.Args[4].(float32)})
				}
			})
			c.fonts = append(c.fonts, o.Fonts...)
		}
	}
	c.text = b.String()
	return c
}

func TestArticle(t *testing.T) {
	opts := testOptions()
	res, err := ConvertFile("testdata/article.html", opts)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Extracted || res.Images != 1 {
		t.Errorf("extracted %v, %d images", res.Extracted, res.Images)
	}
	r := reopen(t, res)
	dc := r.Manifest.Meta.DC
	if dc.Title.First() != "How Canvas Documents Load Quickly" || dc.Creator.First() != "Alex Example" || dc.Language.First() != "en" ||
		dc.Publisher.First() != "The Example Journal" || !strings.HasPrefix(dc.Created.First(), "2026-09-20") ||
		dc.Identifier.First() != "https://journal.example/canvas-documents" || !slices.Equal(dc.Subject, bdf.DCValues{"canvas", "documents", "preview"}) {
		t.Errorf("metadata %+v", dc)
	}
	if r.Manifest.Meta.Source != "html" || len(r.Manifest.Views) != 1 || r.Manifest.Views[0].Kind != bdf.ViewScroll {
		t.Errorf("source %q, views %+v", r.Manifest.Meta.Source, r.Manifest.Views)
	}
	c := viewContent(t, r, 0)
	for _, want := range []string{"Previews of office documents", "What goes into a document", "Figure 1. Parts of a sample document",
		"The fastest code is the code", "bdf generate report.docx report.bdf\nbdf split report.bdf out/"} {
		if !strings.Contains(c.text, want) {
			t.Errorf("text lacks %q:\n%s", want, c.text)
		}
	}
	// the clutter around the article is left out
	for _, junk := range []string{"Advertisement", "Popular", "Privacy", "Subscribe", "tracking"} {
		if strings.Contains(c.text, junk) {
			t.Errorf("text holds %q", junk)
		}
	}
	// the reader view's title, then the article's headings
	if !slices.Equal(c.marks[bdf.MarkHeading], []string{"1", "2", "2"}) {
		t.Errorf("headings %v", c.marks[bdf.MarkHeading])
	}
	// a header cell spanning rows, one spanning columns, then the body
	cells := []string{"A1:A2 col", "B1:C1 col", "B2 col", "C2 col", "A3", "B3", "C3", "A4", "B4", "C4", "A5", "B5", "C5"}
	if len(c.marks[bdf.MarkTable]) != 1 || !slices.Equal(c.marks[bdf.MarkCell], cells) {
		t.Errorf("tables %d, cells %v", len(c.marks[bdf.MarkTable]), c.marks[bdf.MarkCell])
	}
	if len(c.marks[bdf.MarkList]) != 1 || len(c.marks[bdf.MarkListItem]) != 3 {
		t.Errorf("%d lists, %d items", len(c.marks[bdf.MarkList]), len(c.marks[bdf.MarkListItem]))
	}
	// the link to the heading goes to its strip; relative links without a
	// base URL are left out
	if !slices.Equal(c.links, []string{"https://example.com/spec", "#page=1"}) {
		t.Errorf("links %v", c.links)
	}
	if len(c.images) != 1 {
		t.Errorf("images %v", c.images)
	}
}

func TestExtract(t *testing.T) {
	opts := testOptions()
	opts.Extract = ExtractNone
	res, err := ConvertFile("testdata/article.html", opts)
	if err != nil {
		t.Fatal(err)
	}
	c := viewContent(t, reopen(t, res), 0)
	if res.Extracted || !strings.Contains(c.text, "Advertisement") || !strings.Contains(c.text, "Privacy") {
		t.Errorf("extracted %v:\n%s", res.Extracted, c.text)
	}
	// a short document is kept whole
	res, _ = convertHTML(t, `<!DOCTYPE html><html><head><title>Note</title></head><body><nav>Home</nav><p>Hello.</p></body></html>`, testOptions())
	if res.Extracted {
		t.Error("a note was taken for an article")
	}
	opts = testOptions()
	opts.Extract = "bogus"
	if _, err := ConvertBytes([]byte("<p>x"), opts); err == nil {
		t.Error("bad extraction accepted")
	}
}

func TestWhiteSpace(t *testing.T) {
	_, r := convertHTML(t, "<!DOCTYPE html><p>  日本語の\n文章と English\n words  and <b>bold</b>  text. </p>"+
		"<pre>  two\n\tspaces\n</pre><p>a<br>  b</p>", testOptions())
	c := viewContent(t, r, 0)
	want := "日本語の文章と English words and bold text.\n  two\n    spaces\na\nb"
	if c.text != want {
		t.Errorf("text %q, want %q", c.text, want)
	}
}

func TestLists(t *testing.T) {
	_, r := convertHTML(t, `<!DOCTYPE html><ol start="3" type="a"><li>c<li>d</ol><ol reversed><li>two<li>one</ol>`+
		`<ul><li>a<ul><li>b</ul></ul><ul><li><input type="checkbox" checked> done</ul><li>stray</li>`, testOptions())
	c := viewContent(t, r, 0)
	// labels and their items' text are on one line
	want := "c. c\nd. d\n2. two\n1. one\n• a\n◦ b\ndone\n• stray"
	if c.text != want {
		t.Errorf("text %q, want %q", c.text, want)
	}
	if len(c.marks[bdf.MarkList]) != 6 || len(c.marks[bdf.MarkListItem]) != 8 {
		t.Errorf("%d lists, %d items", len(c.marks[bdf.MarkList]), len(c.marks[bdf.MarkListItem]))
	}
}

func TestStructure(t *testing.T) {
	png := base64.StdEncoding.EncodeToString(onePixel())
	_, r := convertHTML(t, `<!DOCTYPE html><html lang="ja"><h2 id="top">見出し</h2><p lang="en">English</p>`+
		`<table><tr><th>Name<td>1</tr><tr><th>Other<td>2</tr></table>`+
		`<figure><img src="data:image/png;base64,`+png+`"><figcaption>図の説明</figcaption></figure>`+
		`<p><img alt="装飾なし" src="data:image/png;base64,`+png+`"> <a href="#top">上へ</a> <a href="javascript:alert(1)">js</a></p>`, testOptions())
	c := viewContent(t, r, 0)
	if !slices.Equal(c.marks[bdf.MarkHeading], []string{"2"}) || !slices.Equal(c.marks[bdf.MarkLang], []string{"en", ""}) {
		t.Errorf("headings %v, languages %v", c.marks[bdf.MarkHeading], c.marks[bdf.MarkLang])
	}
	if !slices.Equal(c.marks[bdf.MarkCell], []string{"A1 row", "B1", "A2 row", "B2"}) {
		t.Errorf("cells %v", c.marks[bdf.MarkCell])
	}
	// a picture without alt takes the figure's caption
	if !slices.Equal(c.marks[bdf.MarkFigure], []string{"図の説明", "装飾なし"}) {
		t.Errorf("figures %v", c.marks[bdf.MarkFigure])
	}
	if !slices.Equal(c.links, []string{"#page=1"}) {
		t.Errorf("links %v", c.links)
	}
	if r.Manifest.Meta.DC.Language.First() != "ja" {
		t.Errorf("language %v", r.Manifest.Meta.DC.Language)
	}
}

// onePixel returns a PNG of one pixel.
func onePixel() []byte {
	b, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAADElEQVR4nGP4z8AAAAMBAQDJ/pLvAAAAAElFTkSuQmCC")
	return b
}

func TestImages(t *testing.T) {
	var fetched []string
	opts := testOptions()
	opts.NoRemote = false
	opts.Dir = "testdata"
	opts.Fetch = func(u string) ([]byte, error) {
		fetched = append(fetched, u)
		if strings.HasSuffix(u, "missing.png") {
			return nil, errors.New("HTTP 404 Not Found")
		}
		if strings.HasSuffix(u, ".svg") {
			return []byte(`<svg xmlns="http://www.w3.org/2000/svg"/>`), nil
		}
		return os.ReadFile("testdata/images/chart.png")
	}
	res, r := convertHTML(t, `<!DOCTYPE html><p><img src="https://example.com/a.png" width="100"> <img src="https://example.com/a.png">`+
		`<img src="https://example.com/missing.png" alt="[missing]"> <img src="https://example.com/badge.svg" alt="build passing">`+
		`<img src="images/icon.png"> <img src="/etc/hosts" alt="[absolute]"></p>`, opts)
	c := viewContent(t, r, 0)
	// fetched once each, the SVG left out with its alternative text
	slices.Sort(fetched)
	if !slices.Equal(fetched, []string{"https://example.com/a.png", "https://example.com/badge.svg", "https://example.com/missing.png"}) {
		t.Errorf("fetched %v", fetched)
	}
	if len(c.images) != 3 || res.Images != 3 {
		t.Fatalf("images %v (%d)", c.images, res.Images)
	}
	// width="100" is 100 CSS pixels; the 400 × 240 px chart is 300 × 180 pt
	if im := c.images[0]; im[2] != 75 || im[3] != 45 {
		t.Errorf("sized image %v", im)
	}
	if im := c.images[1]; im[2] != 300 || im[3] != 180 {
		t.Errorf("natural image %v", im)
	}
	for _, alt := range []string{"[missing]", "build passing", "[absolute]"} {
		if !strings.Contains(c.text, alt) {
			t.Errorf("text %q lacks %q", c.text, alt)
		}
	}
	// NoRemote: only the local picture
	opts.NoRemote = true
	fetched = nil
	res, _ = convertHTML(t, `<p><img src="https://example.com/a.png" alt="remote"><img src="images/icon.png">`, opts)
	if res.Images != 1 || len(fetched) != 0 {
		t.Errorf("NoRemote: %d images, fetched %v", res.Images, fetched)
	}
	// a floating picture goes to the right of the column, the text beside it
	_, r = convertHTML(t, `<p><img src="images/icon.png" align="right" width="64">Text beside the picture.`, opts)
	c = viewContent(t, r, 0)
	if len(c.images) != 1 || c.images[0][0] != 36+432-48 || c.images[0][2] != 48 {
		t.Errorf("floating image %v", c.images)
	}
	// pictures wider than the column shrink to it
	opts.Width = 200
	_, r = convertHTML(t, `<p><img src="images/icon.png" width="1000">`, opts)
	if im := viewContent(t, r, 0).images[0]; im[2] != 200 || im[3] != 200 {
		t.Errorf("fitted image %v", im)
	}
}

func TestLinks(t *testing.T) {
	opts := testOptions()
	opts.BaseURL = "https://example.com/docs/page.html"
	_, r := convertHTML(t, `<!DOCTYPE html><p><a href="other.html">rel</a> <a href="/root">root</a> <a href="mailto:a@example.com">mail</a> <a href="ftp://x">ftp</a></p>`, opts)
	c := viewContent(t, r, 0)
	if !slices.Equal(c.links, []string{"https://example.com/docs/other.html", "https://example.com/root", "mailto:a@example.com"}) {
		t.Errorf("links %v", c.links)
	}
	// a base element takes precedence
	_, r = convertHTML(t, `<!DOCTYPE html><head><base href="https://base.example/x/"></head><p><a href="y">y</a>`, testOptions())
	if c := viewContent(t, r, 0); !slices.Equal(c.links, []string{"https://base.example/x/y"}) {
		t.Errorf("links %v", c.links)
	}
}

func TestHidden(t *testing.T) {
	_, r := convertHTML(t, `<!DOCTYPE html><head><style>p{}</style><script>var x</script></head>`+
		`<p>shown</p><p hidden>attr</p><p style="display: none">css</p><template><p>tpl</p></template><noscript>ns</noscript>`, testOptions())
	if c := viewContent(t, r, 0); c.text != "shown" {
		t.Errorf("text %q", c.text)
	}
}

func TestCharset(t *testing.T) {
	sjis, err := japanese.ShiftJIS.NewEncoder().Bytes([]byte(`<!DOCTYPE html><html><head><meta charset="Shift_JIS"><title>題名</title></head><body><p>日本語の本文</p></body></html>`))
	if err != nil {
		t.Fatal(err)
	}
	res, err := ConvertBytes(sjis, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	r := reopen(t, res)
	if c := viewContent(t, r, 0); c.text != "日本語の本文" || r.Manifest.Meta.DC.Title.First() != "題名" {
		t.Errorf("text %q, title %q", c.text, r.Manifest.Meta.DC.Title.First())
	}
}

func TestMHTML(t *testing.T) {
	png := base64.StdEncoding.EncodeToString(onePixel())
	archive := strings.Join([]string{
		"From: <Saved by Blink>",
		"Snapshot-Content-Location: https://example.com/post/",
		"Subject: A post",
		"MIME-Version: 1.0",
		`Content-Type: multipart/related;`,
		`	type="text/html";`,
		`	boundary="----B"`,
		"",
		"",
		"------B",
		"Content-Type: text/html",
		"Content-ID: <frame-1@mhtml.blink>",
		"Content-Transfer-Encoding: quoted-printable",
		"Content-Location: https://example.com/post/",
		"",
		`<!DOCTYPE html><html><head><meta http-equiv=3D"Content-Type" content=3D"text/html; charset=3DUTF-8"><title>A post</title></head>=`,
		`<body><p>Caf=C3=A9</p><p><img src=3D"pic.png"></p><p><a href=3D"../about">about</a></p></body></html>`,
		"------B",
		"Content-Type: image/png",
		"Content-Transfer-Encoding: base64",
		"Content-Location: https://example.com/post/pic.png",
		"",
		png[:20],
		png[20:],
		"------B--",
		"",
	}, "\r\n")
	if !IsMHTML([]byte(archive)) || IsHTML([]byte(archive)) {
		t.Fatal("not detected as MHTML")
	}
	opts := testOptions()
	opts.NoRemote = false
	opts.Fetch = func(u string) ([]byte, error) { return nil, fmt.Errorf("fetched %s", u) }
	res, r := convertHTML(t, archive, opts)
	c := viewContent(t, r, 0)
	if c.text != "Café\nabout" || res.Images != 1 || !slices.Equal(c.links, []string{"https://example.com/about"}) {
		t.Errorf("text %q, %d images, links %v, warnings %v", c.text, res.Images, c.links, res.Warnings)
	}
}

func TestViewsAndFonts(t *testing.T) {
	src := `<!DOCTYPE html><h1>Title</h1><p>Text <code>code</code></p>`
	opts := testOptions()
	opts.Views = ViewsBoth
	res, r := convertHTML(t, src, opts)
	if len(r.Manifest.Views) != 2 || r.Manifest.Views[0].Kind != bdf.ViewFlow || res.Pages != 1 || res.Strips != 1 {
		t.Fatalf("views %+v", r.Manifest.Views)
	}
	if p := r.Manifest.Views[0].Pages[0]; p.W < 595 || p.W > 596 || p.H < 841 || p.H > 842 {
		t.Errorf("page %v × %v", p.W, p.H)
	}
	// fonts are referred to by name unless embedded
	c := viewContent(t, r, 1)
	if res.EmbeddedFonts != 0 || len(c.fonts) == 0 {
		t.Fatalf("%d embedded, %d fonts", res.EmbeddedFonts, len(c.fonts))
	}
	for _, f := range c.fonts {
		if f.Kind != bdf.FontSystem || !strings.HasSuffix(f.Family, "sans-serif") && !strings.HasSuffix(f.Family, "monospace") {
			t.Errorf("font %+v", f)
		}
	}
	opts = testOptions()
	opts.EmbedFonts = true
	if res, _ := convertHTML(t, src, opts); res.EmbeddedFonts == 0 {
		t.Error("no fonts embedded")
	}
	opts.Views = "bogus"
	if _, err := ConvertBytes([]byte(src), opts); err == nil {
		t.Error("bad views accepted")
	}
}

func TestDeterministic(t *testing.T) {
	var outs [2][]byte
	for i := range outs {
		res, err := ConvertFile("testdata/article.html", testOptions())
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
		t.Error("output differs between runs")
	}
}

func TestRegistered(t *testing.T) {
	f := conv.Lookup("html")
	if f == nil || !slices.Contains(f.Extensions, ".mhtml") {
		t.Fatal("html is not registered")
	}
	// the registry resolves images beside the file
	dir := t.TempDir()
	path := filepath.Join(dir, "page.html")
	os.WriteFile(filepath.Join(dir, "a.png"), onePixel(), 0o644)
	os.WriteFile(path, []byte(`<!DOCTYPE html><p><img src="a.png">`), 0o644)
	res, err := conv.ConvertFile(path, "", &conv.Options{FontDirs: []string{"../docx/testdata/fonts"}, NoSystemFonts: true,
		Params: map[string]string{"remote": "false", "width": "300"}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Summary, "1 image(s)") {
		t.Errorf("summary %q", res.Summary)
	}
	if _, err := conv.ConvertFile(path, "", &conv.Options{Params: map[string]string{"width": "wide"}}); err == nil {
		t.Error("bad width accepted")
	}
}
