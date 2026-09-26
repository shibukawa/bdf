package markdown

import (
	"bytes"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
	"golang.org/x/text/encoding/japanese"
)

// testOptions lays text out with the Word converter's test fonts only, so
// that the output does not depend on the machine.
func testOptions() *Options {
	return &Options{FontDirs: []string{"../docx/testdata/fonts"}, NoSystemFonts: true, NoRemote: true}
}

// content is what the objects of the scroll view hold.
type content struct {
	text  string
	marks map[byte][]string
	links []string
	image int
}

func convert(t *testing.T, res *Result, err error) (*bdf.Reader, *content) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := res.Doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	c := &content{marks: map[byte][]string{}}
	var b strings.Builder
	for _, pg := range r.Manifest.Views[0].Pages {
		o, err := r.Object(pg.Layers[0].Obj)
		if err != nil {
			t.Fatal(err)
		}
		rs, err := bdf.ExtractText(o, func(h bdf.Hash) *bdf.ObjectPart { ch, _ := r.Object(h); return ch })
		if err != nil {
			t.Fatal(err)
		}
		for _, run := range rs {
			if b.Len() > 0 {
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
				c.image++
			}
		})
	}
	c.text = b.String()
	return r, c
}

func convertSrc(t *testing.T, src []byte) (*bdf.Reader, *content) {
	t.Helper()
	res, err := ConvertBytes(src, testOptions())
	return convert(t, res, err)
}

func TestBasic(t *testing.T) {
	res, err := ConvertFile("testdata/basic.md", testOptions())
	r, c := convert(t, res, err)
	m := r.Manifest.Meta
	if m.Source != "markdown" || m.DC.Title.First() != "Markdown のテスト文書" || m.DC.Creator.First() != "テスト太郎" ||
		m.DC.Created.First() != "2026-09-26" || m.DC.Language.First() != "ja" || !slices.Equal(m.DC.Subject, bdf.DCValues{"bdf", "markdown"}) {
		t.Errorf("metadata %+v", m)
	}
	if !slices.Equal(c.marks[bdf.MarkHeading], []string{"1", "2", "2", "2", "2"}) {
		t.Errorf("headings %v", c.marks[bdf.MarkHeading])
	}
	for _, want := range []string{
		// line breaks of the source between Japanese characters are no spaces
		"ソースの改行は和文の間では空白になりません。",
		// tabs of code blocks move to the next multiple of four columns
		"func main() {\n    fmt.Println(\"こんにちは、世界\") // タブで字下げ",
		"用語\n定義リストの説明。", "詳細\n折りたたまれた内容も表示します。", "1. 脚注の本文です。",
		"完了したタスク\n未完了のタスク",
	} {
		if !strings.Contains(c.text, want) {
			t.Errorf("text lacks %q:\n%s", want, c.text)
		}
	}
	// the header row of the table, then its body
	if cells := c.marks[bdf.MarkCell]; len(cells) != 12 || cells[0] != "A1 col" || cells[3] != "A2" {
		t.Errorf("cells %v", cells)
	}
	// the picture of the Markdown and the one of the raw HTML
	if c.image != 2 || !slices.Contains(c.marks[bdf.MarkFigure], "グラデーションの画像") || !slices.Contains(c.marks[bdf.MarkFigure], "アイコン") {
		t.Errorf("%d images, figures %v", c.image, c.marks[bdf.MarkFigure])
	}
	// the external link, the link to a heading (GitHub's id of "リスト"),
	// the footnote reference and its way back
	if !slices.Contains(c.links, "https://example.com/") || len(c.links) != 4 {
		t.Errorf("links %v", c.links)
	}
	for _, l := range c.links[1:] {
		if !strings.HasPrefix(l, "#page=") {
			t.Errorf("links %v", c.links)
		}
	}
}

func TestFrontMatter(t *testing.T) {
	for _, c := range []struct {
		src  string
		want map[string][]string
		body string
	}{
		{"---\ntitle: \"A: title\"\nauthors:\n  - Ann\n  - name: Bob\ntags: [x, 'y']\ndescription: >\n  folded\n  text\nlayout: post # not metadata\n---\nbody",
			map[string][]string{"title": {"A: title"}, "creator": {"Ann", "Bob"}, "subject": {"x", "y"}, "description": {"folded text"}}, "body"},
		{"---\nauthor:\n  name: Carol\nlang: en-US\n...\nbody", map[string][]string{"creator": {"Carol"}, "language": {"en-US"}}, "body"},
		{"+++\ntitle = \"T\"\ntags = [\"a\", \"b\"]\n[params]\nauthor = \"no\"\n+++\nbody", map[string][]string{"title": {"T"}, "subject": {"a", "b"}}, "body"},
		// not closed: a thematic break
		{"---\ntext", map[string][]string{}, "---\ntext"},
	} {
		fm, body := splitFrontMatter([]byte(c.src))
		if string(body) != c.body {
			t.Errorf("%q: body %q", c.src, body)
		}
		if len(fm.dc) != len(c.want) {
			t.Errorf("%q: %v", c.src, fm.dc)
		}
		for k, v := range c.want {
			if !slices.Equal(fm.dc[k], v) {
				t.Errorf("%q: %s = %q, want %q", c.src, k, fm.dc[k], v)
			}
		}
	}
}

func TestHeadingIDs(t *testing.T) {
	doc, err := ToHTML([]byte("# Hello, World!\n# Hello, World!\n## 日本語の見出し（1）\n## C++ & Go_lang"))
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{`id="hello-world"`, `id="hello-world-1"`, `id="日本語の見出し1"`, `id="c--go_lang"`} {
		if !bytes.Contains(doc, []byte(id)) {
			t.Errorf("%s missing:\n%s", id, doc)
		}
	}
}

func TestEncodings(t *testing.T) {
	sjis, _ := japanese.ShiftJIS.NewEncoder().Bytes([]byte("# 見出し\n\n本文"))
	_, c := convertSrc(t, sjis)
	if c.text != "見出し\n本文" {
		t.Errorf("Shift_JIS: %q", c.text)
	}
	_, c = convertSrc(t, []byte("\xef\xbb\xbfplain"))
	if c.text != "plain" {
		t.Errorf("BOM: %q", c.text)
	}
	if IsText([]byte("\x00binary")) || !IsText([]byte("text\n")) || IsText(nil) {
		t.Error("IsText")
	}
}

func TestRawHTML(t *testing.T) {
	// a README's centered logo and a collapsed section
	_, c := convertSrc(t, []byte("<p align=\"center\"><b>Logo</b></p>\n\n<details><summary>More</summary>\n\nHidden *text*.\n\n</details>\n"))
	if c.text != "Logo\nMore\nHidden text." {
		t.Errorf("text %q", c.text)
	}
}

func TestRegistered(t *testing.T) {
	f := conv.Lookup("markdown")
	if f == nil || !f.Fallback || !slices.Contains(f.Extensions, ".md") {
		t.Fatal("markdown is not registered as the fallback format")
	}
	res, err := conv.ConvertFile("testdata/basic.md", "", &conv.Options{FontDirs: []string{"../docx/testdata/fonts"}, NoSystemFonts: true,
		Params: map[string]string{"remote": "false", "views": "both"}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Summary, "page(s)") || !strings.Contains(res.Summary, "2 image(s)") || len(res.Doc.Views) != 2 {
		t.Errorf("summary %q", res.Summary)
	}
}
