package markdownbdf

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// open reads a test document, skipping the test when it is missing.
func open(t *testing.T, name string) *bdf.Document {
	t.Helper()
	path := filepath.Join("../../testdata", name)
	if _, err := os.Stat(path); err != nil {
		t.Skipf("no %s", path)
	}
	r, err := bdf.OpenSingleFile(path)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	d, err := r.ToDocument()
	if err != nil {
		t.Fatal(err)
	}
	return d
}

// write writes a document, collecting its images.
func write(t *testing.T, d *bdf.Document, opts *Options) (string, map[string][]byte, *Result) {
	t.Helper()
	images := map[string][]byte{}
	if opts == nil {
		opts = &Options{}
	}
	opts.Image = func(name string, data []byte) (string, error) {
		images[name] = data
		return "images/" + name, nil
	}
	var b bytes.Buffer
	res, err := Write(&b, d, opts)
	if err != nil {
		t.Fatal(err)
	}
	return b.String(), images, res
}

func wantLines(t *testing.T, md string, lines ...string) {
	t.Helper()
	for _, l := range lines {
		if !strings.Contains(md, l) {
			t.Errorf("missing %q in:\n%s", l, md)
		}
	}
}

// TestMarkdownRoundTrip writes the Markdown fixture back as Markdown: the
// blocks of the source come back.
func TestMarkdownRoundTrip(t *testing.T) {
	md, images, res := write(t, open(t, "markdown/basic.bdf"), nil)
	wantLines(t, md,
		"---\ntitle: Markdown のテスト文書\ncreator: テスト太郎\n",
		"subject: [bdf, markdown]\n",
		"\n# Markdown のテスト文書\n",
		"段落の本文です。**太字**、*斜体*、取り消し線、`インラインのコード`、[リンク](https://example.com/)を含みます。",
		"\n## リスト\n",
		"- 項目 1\n- 項目 2\n  - 入れ子の項目 A\n  - 入れ子の項目 B\n    1. 番号付き\n    2. 番号付き\n- 項目 3\n",
		"1. 手順 1\n2. 手順 2\n\n   続きの段落。\n",
		"```\nfunc main() {\n    fmt.Println(\"こんにちは、世界\") // タブで字下げ\n",
		"| 左寄せ | 中央 | 右寄せ |\n|---|---|---|\n| りんご | 赤 | 120 |\n",
		"![グラデーションの画像](images/",
		"見出しへのリンク: リスト、脚注<sup>1</sup>。", // links into the document drop their targets; a footnote mark is a superscript
	)
	if strings.Contains(md, "#page=") {
		t.Error("an internal link was kept")
	}
	if len(images) != 2 || res.Images != 2 {
		t.Errorf("%d images written, %d counted", len(images), res.Images)
	}
	for name, data := range images {
		if !strings.HasSuffix(name, ".webp") || !bytes.HasPrefix(data, []byte("RIFF")) {
			t.Errorf("image %s: %d bytes, %q", name, len(data), data[:min(4, len(data))])
		}
	}
	// the scroll view is written: the flow view's pages would repeat it
	if res.Views != 1 {
		t.Errorf("%d views", res.Views)
	}
}

// TestInternalLinks keeps the links into the document when asked.
func TestInternalLinks(t *testing.T) {
	md, _, _ := write(t, open(t, "markdown/basic.bdf"), &Options{InternalLinks: true})
	wantLines(t, md, "[リスト](#page=1)")
}

// TestWord writes a Word document: its headings, the bold of an embedded
// face, lists, a table with a merged cell, and no running heads.
func TestWord(t *testing.T) {
	md, _, _ := write(t, open(t, "docx/basic.bdf"), nil)
	wantLines(t, md,
		"\n# Text and paragraphs\n",
		"Runs can be **bold**, *italic*, underlined,",
		"SMALL CAPS; x<sup>2</sup> and H<sub>2</sub>O. A link to [example.com](https://example.com/) and a note.<sup>1</sup>",
		"- A bulleted item\n  - A nested item\n  - Another nested item\n- Back at the first level\n",
		"1. First step\n   1. a detail\n   2. another detail\n2. Second step,",
		"| Item | Quantity | Note |\n|---|---|---|\n| Apples | 12 | Fresh from the orchard |\n",
		"| Total of the fruit above, in one merged cell |  | 合計は 19 個 |",
		"![A gradient photo](images/",
		"Column paragraph 1.", "Column paragraph 2.", "Column paragraph 3.",
	)
	if strings.Count(md, "Column paragraph 1.") != 1 {
		t.Error("the flow view's pages were written besides the scroll view")
	}
}

// TestSlides writes a deck: a heading a slide, the levels of its bullets,
// a link, a table and a picture, without the backgrounds.
func TestSlides(t *testing.T) {
	md, _, _ := write(t, open(t, "pptx/basic.bdf"), nil)
	wantLines(t, md,
		"\n# BDF from PowerPoint\n\nSlides rendered directly from DrawingML\n",
		"\n# Bullets and levels\n\n- First level bullet with enough words to wrap onto a second line inside the placeholder\n  - Second level\n    - Third level\n- [A link to example.com](https://example.com/)\n  1. numbered one\n  2. numbered two\n  3. numbered three\n",
		"| Header 0 | Header 1 | Header 2 |\n|---|---|---|\n| r1c0 | r1c1 | r1c2 |",
		"![image.png](images/",
	)
	if strings.Contains(md, "<!--") {
		t.Error("pages are marked")
	}
}

// TestSheets writes a workbook: a table a sheet, under the sheet's name.
func TestSheets(t *testing.T) {
	md, _, res := write(t, open(t, "xlsx/basic.bdf"), nil)
	wantLines(t, md,
		"\n# Sales\n",
		"| **Region** | **Q1** | **Q2** | **Growth** | **Updated** | **Flag** |", // bold cells of an embedded face
		"| North ノース | 125,000 | 131,250.50 | 5.0% | 2026-03-31 | TRUE |",
		"[example.com](https://example.com/)",
		"\n# Big\n",
		"| row 2 | 4 |",
	)
	if res.Views != 2 {
		t.Errorf("%d views", res.Views)
	}
	// a cell that straddles tiles is written once
	if n := strings.Count(md, "flows across the tile boundary"); n != 1 {
		t.Errorf("the straddling cell is written %d times", n)
	}
}

// TestTaggedPDF writes a tagged PDF: its structure tree became MARKs.
func TestTaggedPDF(t *testing.T) {
	md, _, _ := write(t, open(t, "pdf/chrome-doc.bdf"), nil)
	wantLines(t, md,
		"\n# BDF conversion fixture\n",
		"It contains emphasis, **strong text**, `monospace code` and a [hyperlink](https://example.com/doc).",
		"- First bullet item\n- Second bullet item",
		"| Name | Value | Note |\n|---|---|---|\n| alpha | 1.00 | first |",
		"\n## Second page\n",
		"followed by a justified paragraph.", // the ligature's ALT_TEXT joins its word
	)
}

// TestUntaggedPDF writes a PDF without structure: its lines, split into
// paragraphs by their spacing, with its bold lines as headings.
func TestUntaggedPDF(t *testing.T) {
	md, _, _ := write(t, open(t, "pdf/reportlab-mixed.bdf"), nil)
	wantLines(t, md,
		"Helvetica: The quick brown fox jumps over the lazy dog",
		"[link annotation over this text](https://example.com/rl)",
		"![](images/",
	)
}

// TestViews picks views by id, and refuses an unknown one.
func TestViews(t *testing.T) {
	d := open(t, "docx/basic.bdf")
	md, _, res := write(t, d, &Options{Views: []string{"pages"}})
	if res.Views != 1 || !strings.Contains(md, "# Text and paragraphs") {
		t.Errorf("views = %d:\n%s", res.Views, md[:min(len(md), 400)])
	}
	if _, err := Write(&bytes.Buffer{}, d, &Options{Views: []string{"nope"}}); err == nil {
		t.Error("an unknown view was accepted")
	}
	// the running heads come back when their roles are kept
	md, _, _ = write(t, d, &Options{Views: []string{"pages"}, SkipRoles: []string{}})
	if !strings.Contains(md, "Footnotes go at the bottom") {
		t.Error("the body of the pages is missing")
	}
}

// TestWriteDir writes the file and its images into a directory.
func TestWriteDir(t *testing.T) {
	dir := t.TempDir()
	res, err := WriteDir(filepath.Join(dir, "out"), "basic", open(t, "markdown/basic.bdf"), nil)
	if err != nil {
		t.Fatal(err)
	}
	md, err := os.ReadFile(filepath.Join(dir, "out", "basic.md"))
	if err != nil {
		t.Fatal(err)
	}
	re := regexp.MustCompile(`!\[[^\]]*\]\(([^)]+)\)`)
	refs := re.FindAllStringSubmatch(string(md), -1)
	if len(refs) != res.Images {
		t.Fatalf("%d image references, %d images", len(refs), res.Images)
	}
	for _, m := range refs {
		if _, err := os.Stat(filepath.Join(dir, "out", m[1])); err != nil {
			t.Error(err)
		}
	}
}

// TestNoTextLost checks every test document: each run of text of the
// layers that are written appears in the Markdown (list markers aside,
// which the list's own replace).
func TestNoTextLost(t *testing.T) {
	files, _ := filepath.Glob("../../testdata/*/*.bdf")
	if len(files) == 0 {
		t.Skip("no test documents")
	}
	unescape := func(s string) string {
		return strings.NewReplacer("\\", "", "`", "", "<br>", " ", "<sup>", "", "</sup>", "", "<sub>", "", "</sub>", "").Replace(s)
	}
	for _, path := range files {
		r, err := bdf.OpenSingleFile(path)
		if err != nil {
			t.Errorf("%s: %v", path, err)
			continue
		}
		if r.Locked() {
			r.Close()
			continue
		}
		d, err := r.ToDocument()
		r.Close()
		if err != nil {
			t.Errorf("%s: %v", path, err)
			continue
		}
		var b bytes.Buffer
		if _, err := Write(&b, d, &Options{KeepRunning: true, Image: func(name string, _ []byte) (string, error) { return name, nil }}); err != nil {
			t.Errorf("%s: %v", path, err)
			continue
		}
		md := unescape(b.String())
		views, _ := selectViews(d, nil)
		skip := map[string]bool{}
		for _, role := range DefaultSkipRoles {
			skip[role] = true
		}
		objs := map[bdf.Hash]*bdf.ObjectPart{}
		resolve := func(h bdf.Hash) *bdf.ObjectPart {
			if o, ok := objs[h]; ok {
				return o
			}
			var o *bdf.ObjectPart
			if p := d.Part(h); p != nil {
				o, _ = bdf.DecodeObject(p.Data)
			}
			objs[h] = o
			return o
		}
		missing := 0
		check := func(h bdf.Hash, keep func(*bdf.TextRun) bool) {
			o := resolve(h)
			if o == nil {
				return
			}
			runs, _ := bdf.ExtractText(o, resolve)
			for i := range runs {
				r := &runs[i]
				text := strings.TrimSpace(r.Text)
				if text == "" || keep != nil && !keep(r) || bulletRe.MatchString(text) || orderedRe.MatchString(text) {
					continue
				}
				if !strings.Contains(md, unescape(text)) && missing < 5 {
					t.Errorf("%s: %q is not in the Markdown", path, text)
					missing++
				}
			}
		}
		for _, v := range views {
			if v.Kind == bdf.ViewSheet {
				tile := v.TileSize()
				for _, k := range tileKeys(v) {
					h, _ := bdf.ParseHash(v.Tiles[k])
					check(h, func(r *bdf.TextRun) bool { return r.X >= 0 && r.X < tile && r.Y >= 0 && r.Y < tile })
				}
				continue
			}
			for _, p := range v.Pages {
				for _, l := range p.Layers {
					if !skip[l.Role] {
						check(l.Obj, nil)
					}
				}
			}
		}
	}
}

// TestRunning drops the running heads, feet and page numbers of a paged
// document whose pages hold them (as a PDF's do), and keeps them when
// asked. The Word fixture's pages, with their header and footer layers
// kept, stand in for such a document.
func TestRunning(t *testing.T) {
	d := open(t, "docx/basic.bdf")
	md, _, _ := write(t, d, &Options{Views: []string{"pages"}, SkipRoles: []string{}})
	for _, s := range []string{"BDF test document", "Page 1", "\n1\n", "\n2\n"} {
		if strings.Contains(md, s) {
			t.Errorf("%q is in the Markdown:\n%s", s, md[:min(len(md), 600)])
		}
	}
	wantLines(t, md, "Footnotes go at the bottom", "# Text and paragraphs")
	md, _, _ = write(t, d, &Options{Views: []string{"pages"}, SkipRoles: []string{}, KeepRunning: true})
	wantLines(t, md, "BDF test document Page 1")
}
