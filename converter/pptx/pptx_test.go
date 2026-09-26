package pptx

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
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

// convertEdited converts a test deck with some of its parts rewritten.
func convertEdited(t *testing.T, name string, edits map[string]func(string) string, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	zw := zip.NewWriter(&out)
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		b, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Fatal(err)
		}
		if edit, ok := edits[f.Name]; ok {
			s := edit(string(b))
			if s == string(b) {
				t.Fatalf("edit of %s changed nothing", f.Name)
			}
			b = []byte(s)
		}
		w, err := zw.Create(f.Name)
		if err != nil {
			t.Fatal(err)
		}
		w.Write(b)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	res, err := Convert(bytes.NewReader(out.Bytes()), int64(out.Len()), opts)
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

func plainText(t *testing.T, r *bdf.Reader) string {
	t.Helper()
	h, err := bdf.ParseHash(r.Manifest.Views[0].TextIndex)
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
	return bdf.PlainText(idx)
}

// opCounts counts the instructions of an object and the objects it uses.
func opCounts(t *testing.T, r *bdf.Reader, h bdf.Hash) map[byte]int {
	t.Helper()
	counts := map[byte]int{}
	var walk func(h bdf.Hash)
	walk = func(h bdf.Hash) {
		o, err := r.Object(h)
		if err != nil {
			t.Fatal(err)
		}
		if err := o.Walk(func(in bdf.Instr) { counts[in.Op]++ }); err != nil {
			t.Fatal(err)
		}
		for _, c := range o.Objects {
			walk(c)
		}
	}
	walk(h)
	return counts
}

func markCounts(t *testing.T, r *bdf.Reader, h bdf.Hash) map[byte]int {
	t.Helper()
	counts := map[byte]int{}
	o, err := r.Object(h)
	if err != nil {
		t.Fatal(err)
	}
	o.Walk(func(in bdf.Instr) {
		if in.Op == bdf.OpMark {
			counts[byte(in.Args[0].(uint64))]++
		}
	})
	return counts
}

// structure lists the structure MARKs of an object and the text between
// them, separated by " | ": runs joined as text extraction joins them
// (LINE is a space, the other MARKs listed), ALT_TEXT standing for the
// drawing op after it, IMAGE for images, and marks with their payload in
// parentheses (always for LANG).
func structure(t *testing.T, r *bdf.Reader, h bdf.Hash) string {
	t.Helper()
	names := map[byte]string{bdf.MarkParagraph: "P", bdf.MarkCell: "CELL", bdf.MarkBox: "BOX", bdf.MarkHeading: "H",
		bdf.MarkList: "LIST", bdf.MarkListItem: "LI", bdf.MarkTable: "TABLE", bdf.MarkFigure: "FIG", bdf.MarkEnd: "END", bdf.MarkLang: "LANG"}
	var out []string
	inText, space := false, false
	var alt *string
	add := func(s string) {
		out = append(out, s)
		inText, space = false, false
	}
	addText := func(s string) {
		switch {
		case !inText:
			out = append(out, s)
		case space:
			out[len(out)-1] += " " + s
		default:
			out[len(out)-1] += s
		}
		inText, space = true, false
	}
	var walk func(h bdf.Hash)
	walk = func(h bdf.Hash) {
		o, err := r.Object(h)
		if err != nil {
			t.Fatal(err)
		}
		err = o.Walk(func(in bdf.Instr) {
			switch in.Op {
			case bdf.OpMark:
				kind, payload := byte(in.Args[0].(uint64)), in.Args[1].(string)
				switch kind {
				case bdf.MarkAltText:
					alt = &payload
				case bdf.MarkLine:
					space = true
				case bdf.MarkWrap:
				default:
					if payload != "" || kind == bdf.MarkLang {
						add(names[kind] + "(" + payload + ")")
					} else {
						add(names[kind])
					}
				}
			case bdf.OpFillText:
				if alt != nil {
					addText(*alt)
					alt = nil
					return
				}
				addText(in.Args[0].(string))
			case bdf.OpImage, bdf.OpImageSub:
				add("IMAGE")
			case bdf.OpUse, bdf.OpUseAt:
				if alt != nil {
					addText(*alt)
					alt = nil
					return
				}
				walk(o.Objects[int(in.Args[0].(uint64))])
			}
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	walk(h)
	for i := range out {
		out[i] = strings.TrimSpace(out[i])
	}
	return strings.Join(out, " | ")
}

func body(p *bdf.Page) bdf.Hash {
	for _, l := range p.Layers {
		if l.Role == bdf.RoleBody {
			return l.Obj
		}
	}
	return bdf.Hash{}
}

func TestConvertBasic(t *testing.T) {
	res, r := convert(t, "basic.pptx", testOptions())
	if res.Slides != 5 {
		t.Fatalf("slides = %d (the hidden slide must be skipped)", res.Slides)
	}
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	if res.EmbeddedFonts != 2 {
		t.Errorf("embedded fonts = %d, want the regular and bold test fonts", res.EmbeddedFonts)
	}
	v := r.Manifest.Views[0]
	if v.Kind != bdf.ViewFixed || v.Pages[0].W != 720 || v.Pages[0].H != 540 {
		t.Fatalf("view = %s %gx%g", v.Kind, v.Pages[0].W, v.Pages[0].H)
	}
	// dc.language: the default text style's, as the core properties name none
	if dc := r.Manifest.Meta.DC; dc.Title.First() != "PowerPoint test deck" || r.Manifest.Meta.Source != "pptx" ||
		dc.Description.First() != "generated using python-pptx" || len(dc.Creator) != 0 ||
		dc.Created.First() != "2013-01-27T09:14:16Z" || dc.Modified.First() != "2013-01-27T09:15:58Z" || dc.Language.First() != "en-US" {
		t.Errorf("meta = %+v", r.Manifest.Meta)
	}
	if v.Title != "PowerPoint test deck" {
		t.Errorf("view title = %q", v.Title)
	}
	text := plainText(t, r)
	for _, want := range []string{
		"BDF from PowerPoint\nSlides rendered directly from DrawingML",
		"• First level bullet with enough words to wrap onto a second line inside the placeholder",
		"– Second level", "1. numbered one\n2. numbered two\n3. numbered three",
		"• A link to example.com",
		"Aligned justify: the quick brown fox jumps over the lazy dog and keeps running.",
		// wrapped East Asian lines join without a space
		"日本語の文章は、単語の区切りに空白を使わないため、文字と文字の間で改行します。句読点「、」や「。」は行頭に来ません。",
		"混在 mixed テキスト x2 underline", "縦書きのテキスト",
		"Header 0", "merged cell with longer text that wraps",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("text lacks %q", want)
		}
	}
	if strings.Contains(text, "Hidden slide") {
		t.Error("hidden slide converted")
	}
	if c := opCounts(t, r, body(v.Pages[1])); c[bdf.OpLink] != 1 {
		t.Errorf("slide 2 links = %d", c[bdf.OpLink])
	}
	if m := markCounts(t, r, body(v.Pages[3])); m[bdf.MarkWrap] == 0 || m[bdf.MarkLine] == 0 || m[bdf.MarkBox] == 0 {
		t.Errorf("slide 4 marks = %v", m)
	}
	// the picture and a table
	if c := opCounts(t, r, body(v.Pages[4])); c[bdf.OpImage]+c[bdf.OpImageSub] != 1 || c[bdf.OpStrokePath] < 10 {
		t.Errorf("slide 5 ops = %v", c)
	}
}

func TestStructure(t *testing.T) {
	_, r := convert(t, "basic.pptx", testOptions())
	pages := r.Manifest.Views[0].Pages
	for _, c := range []struct {
		page int
		want string
	}{
		// the title placeholders are headings, after their BOX
		{1, "BOX | H(1) | BDF from PowerPoint | BOX | Slides rendered directly from DrawingML"},
		// bullets are list items, nested by level; numbered ones too
		{2, "BOX | H(1) | Bullets and levels | BOX | LIST | LI | • First level bullet with enough words to wrap onto a second line inside the placeholder" +
			" | LIST | LI | – Second level | LIST | LI | • Third level | END | END | LI | • A link to example.com" +
			" | LIST | LI | 1. numbered one | LI | 2. numbered two | LI | 3. numbered three | END | END"},
		// Japanese text (kana) is marked; the Han-only run of a Japanese
		// paragraph and the vertical line (in the parent, before its
		// ALT_TEXT) too; Latin runs return to the document's language
		{4, "BOX | LANG(ja) | 日本語の文章は、単語の区切りに空白を使わないため、文字と文字の間で改行します。句読点「、」や「。」は行頭に来ません。" +
			" | P | 混在 | LANG() | mixed | LANG(ja) | テキスト | LANG() | x2 underline | BOX | LANG(ja) | 縦書きのテキスト | BOX | LANG() | Middle anchored"},
		// the picture is a figure with its description; the table has
		// column headers from firstRow and a merged cell as a range
		{5, "BOX | H(1) | Picture and table | FIG(image.png) | IMAGE | END | TABLE" +
			" | CELL(A1 col) | BOX | Header 0 | CELL(B1 col) | BOX | Header 1 | CELL(C1 col) | BOX | Header 2" +
			" | CELL(A2) | BOX | r1c0 | CELL(B2) | BOX | r1c1 | CELL(C2) | BOX | r1c2" +
			" | CELL(A3) | BOX | r2c0 | CELL(B3) | BOX | r2c1 | CELL(C3) | BOX | r2c2" +
			" | CELL(A4:B4) | BOX | merged cell with longer text that wraps | CELL(C4) | BOX | r3c2 | END"},
	} {
		got := structure(t, r, body(pages[c.page-1]))
		if !strings.Contains(got, c.want) {
			t.Errorf("slide %d:\n got %s\nwant %s", c.page, got, c.want)
		}
	}

	_, r = convert(t, "features.pptx", testOptions())
	pages = r.Manifest.Views[0].Pages
	for _, c := range []struct {
		page int
		want string
	}{
		{1, "BOX | H(1) | Master and layout | BOX | LIST | LI | • The band at the bottom comes from the slide master | LI | • The orange bar comes from the layout | END"},
		// Han-only text keeps the Japanese in effect (the table cells)
		{4, "BOX | H(1) | LANG(ja) | 日本語のスライド | BOX | LIST | LI | • 箇条書きの一行目です。長い文章は枠の幅で折り返され、句読点は行頭に来ません。" +
			" | LIST | LI | – 二段目の項目 | END | LI | • 英数字 | LANG() | ABC | LANG(ja) | と | LANG() | 123 | LANG(ja) | の混在 | END" +
			" | BOX | 縦書きで、読みます。 | TABLE | CELL(A1 col) | BOX | 項目 | CELL(B1 col) | BOX | 説明" +
			" | CELL(A2) | BOX | 図形 | CELL(B2) | BOX | プリセット形状 | CELL(A3) | BOX | 表 | CELL(B3) | BOX | 表のスタイル | END"},
	} {
		got := structure(t, r, body(pages[c.page-1]))
		if got != c.want {
			t.Errorf("features slide %d:\n got %s\nwant %s", c.page, got, c.want)
		}
	}
	// text on masters is no heading
	if got := structure(t, r, pages[0].Layers[1].Obj); got != "BOX | BDF features deck" {
		t.Errorf("features master: %s", got)
	}
}

func TestFigures(t *testing.T) {
	const decorative = `<a:extLst><a:ext uri="{C183D7F6-B498-43B3-948B-1728B52AA6E4}">` +
		`<adec:decorative xmlns:adec="http://schemas.microsoft.com/office/drawing/2017/decorative" val="1"/></a:ext></a:extLst>`
	_, r := convertEdited(t, "basic.pptx", map[string]func(string) string{
		"ppt/slides/slide3.xml": func(s string) string {
			return strings.NewReplacer(
				`<p:cNvPr id="3" name="Rectangle 2"/>`, `<p:cNvPr id="3" name="Rectangle 2" descr="Blue box"/>`,
				`<p:cNvPr id="4" name="Rounded Rectangle 3"/>`, `<p:cNvPr id="4" name="Rounded Rectangle 3" descr="Ornament">`+decorative+`</p:cNvPr>`,
				`<p:cNvPr id="5" name="Oval 4"/>`, `<p:cNvPr id="5" name="Oval 4" title="Green oval"/>`,
				`<p:cNvPr id="12" name="Group 11"/>`, `<p:cNvPr id="12" name="Group 11" descr="Two shapes"/>`,
			).Replace(s)
		},
		"docProps/core.xml": func(s string) string {
			return strings.Replace(s, "<dc:subject/>", "<dc:subject/><dc:language>ja-JP</dc:language>", 1)
		},
	}, testOptions())
	if l := r.Manifest.Meta.DC.Language.First(); l != "ja-JP" {
		t.Errorf("dc.language = %q, want the core properties' dc:language", l)
	}
	pages := r.Manifest.Views[0].Pages
	// the figure covers the shape's drawing, not its text; decorative
	// drawings are no figures; a group's figure covers its members. (The
	// title has no lang in its style chain, the other shapes take en-US
	// from the default text style.)
	got := structure(t, r, body(pages[2]))
	for _, want := range []string{
		"BOX | H(1) | Shapes | FIG(Blue box) | END | BOX | LANG(en-US) | Text 0 | FIG(Green oval) | END | BOX | Text 4",
		"FIG(Two shapes) | BOX | grp | END",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("slide 3:\n got %s\nwant %s", got, want)
		}
	}
	if strings.Contains(got, "Ornament") || strings.Count(got, "FIG") != 3 {
		t.Errorf("slide 3 figures: %s", got)
	}
	// "ja" in a "ja-JP" document is the document's language
	got = structure(t, r, body(pages[3]))
	if want := "BOX | LANG() | 日本語の文章は"; !strings.Contains(got, want) || strings.Contains(got, "LANG(ja)") {
		t.Errorf("slide 4:\n got %s\nwant %s…", got, want)
	}
	// the text itself does not change
	_, plain := convert(t, "basic.pptx", testOptions())
	if a, b := plainText(t, r), plainText(t, plain); a != b {
		t.Errorf("text changed:\n%s\n---\n%s", a, b)
	}
}

func TestSharedLayers(t *testing.T) {
	res, r := convert(t, "features.pptx", testOptions())
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	v := r.Manifest.Views[0]
	if v.Pages[0].W != 960 || v.Pages[0].H != 540 {
		t.Fatalf("page = %gx%g", v.Pages[0].W, v.Pages[0].H)
	}
	a, b := v.Pages[0].Layers, v.Pages[1].Layers
	roles := func(ls []bdf.Layer) string {
		var s []string
		for _, l := range ls {
			s = append(s, l.Role)
		}
		return strings.Join(s, ",")
	}
	if roles(a) != "background,master,master,body" || roles(b) != roles(a) {
		t.Fatalf("layers = %s / %s", roles(a), roles(b))
	}
	for i := 0; i < 3; i++ {
		if a[i].Obj != b[i].Obj {
			t.Errorf("layer %d (%s) is not shared", i, a[i].Role)
		}
	}
	if a[3].Obj == b[3].Obj {
		t.Error("bodies should differ")
	}
	// slide 5 (Title Only layout) shares the master layer but not the layout one
	c := v.Pages[4].Layers
	if c[1].Obj != a[1].Obj || c[2].Obj == a[2].Obj {
		t.Errorf("master/layout sharing across layouts is wrong")
	}
	text := plainText(t, r)
	for _, want := range []string{"BDF features deck", "Share", "Q1", "East", "45%", "日本語のスライド",
		"箇条書きの一行目です。長い文章は枠の幅で折り返され、句読点は行頭に来ません。", "縦書きで、読みます。", "プリセット形状", "Callout"} {
		if !strings.Contains(text, want) {
			t.Errorf("text lacks %q", want)
		}
	}
}

func TestOptions(t *testing.T) {
	opts := testOptions()
	opts.Hidden = true
	res, r := convert(t, "basic.pptx", opts)
	if res.Slides != 6 || !strings.Contains(plainText(t, r), "Hidden slide") {
		t.Errorf("Hidden: %d slides", res.Slides)
	}
	opts = testOptions()
	opts.Slides = conv.PageList(2, 4)
	res, _ = convert(t, "basic.pptx", opts)
	if res.Slides != 2 {
		t.Errorf("Slides: %d", res.Slides)
	}
	opts.Slides = conv.PageList(9)
	if _, err := ConvertFile("testdata/basic.pptx", opts); err == nil {
		t.Error("out of range slide accepted")
	}
	opts = testOptions()
	opts.SystemFonts = true
	res, r = convert(t, "basic.pptx", opts)
	if res.EmbeddedFonts != 0 {
		t.Errorf("SystemFonts embedded %d fonts", res.EmbeddedFonts)
	}
	for _, p := range r.Manifest.Parts {
		if p.T == bdf.PartFont {
			t.Fatal("font part with SystemFonts")
		}
	}
	o, err := r.Object(body(r.Manifest.Views[0].Pages[0]))
	if err != nil {
		t.Fatal(err)
	}
	if len(o.Fonts) == 0 || o.Fonts[0].Kind != bdf.FontSystem || !strings.Contains(o.Fonts[0].Family, `"Calibri"`) {
		t.Errorf("fonts = %+v", o.Fonts)
	}
}

func TestDeterministic(t *testing.T) {
	var outs [2][]byte
	for i := range outs {
		res, err := ConvertFile("testdata/features.pptx", testOptions())
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
