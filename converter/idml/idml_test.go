package idml

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
)

// The test documents are made by test/idml/gen.py; text is laid out with
// the test fonts of the PowerPoint converter.

func testOptions() *Options {
	return &Options{FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true}
}

func convert(t *testing.T, name string, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
	if opts == nil {
		opts = testOptions()
	}
	if opts.Dir == "" {
		opts.Dir = "testdata"
	}
	res, err := ConvertFile(filepath.Join("testdata", name), opts)
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
	return res, r
}

// pageTexts returns the text of each page of the first view.
func pageTexts(t *testing.T, r *bdf.Reader) []string {
	t.Helper()
	h, err := bdf.ParseHash(r.Manifest.Views[0].TextIndex)
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.Part(h)
	if err != nil {
		t.Fatal(err)
	}
	runs, err := bdf.DecodeTextIndex(b)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for i := 0; i < len(runs); {
		j := i
		for j < len(runs) && runs[j].A == runs[i].A {
			j++
		}
		for int(runs[i].A) >= len(out) {
			out = append(out, "")
		}
		out[runs[i].A] = bdf.PlainText(runs[i:j])
		i = j
	}
	return out
}

// ops counts the instructions of an object with a given opcode.
func ops(t *testing.T, r *bdf.Reader, h bdf.Hash, op byte) int {
	t.Helper()
	o, err := r.Object(h)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	o.Walk(func(in bdf.Instr) {
		if in.Op == op {
			n++
		}
	})
	return n
}

func layer(p *bdf.Page, role string) (bdf.Hash, bool) {
	for _, l := range p.Layers {
		if l.Role == role {
			return l.Obj, true
		}
	}
	return bdf.Hash{}, false
}

func TestBasic(t *testing.T) {
	res, r := convert(t, "basic.idml", nil)
	if res.Pages != 3 || len(r.Manifest.Views) != 1 || len(r.Manifest.Views[0].Pages) != 3 {
		t.Fatalf("%d pages, %d views", res.Pages, len(r.Manifest.Views))
	}
	v := r.Manifest.Views[0]
	if v.Kind != bdf.ViewFixed || v.Direction != "" {
		t.Errorf("view %s %q", v.Kind, v.Direction)
	}
	for i, p := range v.Pages {
		if p.W != 400 || p.H != 560 {
			t.Errorf("page %d: %g × %g", i+1, p.W, p.H)
		}
		if _, ok := layer(p, bdf.RoleMaster); !ok {
			t.Errorf("page %d: no master layer", i+1)
		}
		if _, ok := layer(p, bdf.RoleBody); !ok {
			t.Errorf("page %d: no body layer", i+1)
		}
	}
	if res.EmbeddedFonts != 2 {
		t.Errorf("%d embedded fonts, want the regular and bold subsets", res.EmbeddedFonts)
	}
	dc := r.Manifest.Meta.DC
	if dc.Title.First() != "IDML のテスト" || dc.Creator.First() != "BDF tests" || dc.Language.First() != "ja-JP" {
		t.Errorf("metadata %+v", dc)
	}
	texts := pageTexts(t, r)
	// the story flows through the three threaded frames in page order
	for i, want := range []string{"InDesign サンプル", "ながい文章", "the story ends"} {
		if !strings.Contains(texts[i], want) {
			t.Errorf("page %d lacks %q:\n%s", i+1, want, texts[i])
		}
	}
	if strings.Contains(texts[0], "the story ends") || strings.Contains(texts[2], "InDesign サンプル") {
		t.Error("the story did not flow across the frames")
	}
	// page numbers from the master's frame, numbered by the section
	for i, want := range []string{"- 1 -", "- 2 -", "- 3 -"} {
		if !strings.Contains(texts[i], want) {
			t.Errorf("page %d lacks its number %q", i+1, want)
		}
	}
	// a continued paragraph is whole: the last words of the paragraph cut
	// between pages 2 and 3 are on page 3
	if !strings.Contains(texts[1], "日本語の文章") || !strings.Contains(texts[2], "行頭に来ません") {
		t.Error("the paragraph cut between pages lost text")
	}
	// the overset frame shows what fits and warns about the rest
	overset := false
	for _, w := range res.Warnings {
		if strings.Contains(w, "overset") {
			overset = true
		}
	}
	if !overset {
		t.Errorf("no overset warning in %q", res.Warnings)
	}
	if len(res.Warnings) != 1 {
		t.Errorf("warnings %q", res.Warnings)
	}
	// the master layer is the same object on every page but for the page
	// number: the rule is shared through the text frame's own object
	m1, _ := layer(v.Pages[0], bdf.RoleMaster)
	if ops(t, r, m1, bdf.OpStrokePath) != 1 {
		t.Errorf("master layer strokes %d paths, want the rule", ops(t, r, m1, bdf.OpStrokePath))
	}
	// two placed images: the embedded one and the linked one
	images := 0
	for _, p := range r.Manifest.Parts {
		if p.T == bdf.PartImage {
			images++
		}
	}
	if images != 2 {
		t.Errorf("%d image parts, want 2", images)
	}
	b2, _ := layer(v.Pages[1], bdf.RoleBody)
	if n := ops(t, r, b2, bdf.OpImage); n != 2 {
		t.Errorf("page 2 draws %d images, want 2", n)
	}
	// page 1: the items on the hidden layer and the invisible item are not
	// drawn: the shapes with fills are the rounded rectangle, two ovals,
	// the triangle, the group's square and circle, and the frame of the
	// text (none)
	b1, _ := layer(v.Pages[0], bdf.RoleBody)
	if n := ops(t, r, b1, bdf.OpFillPath); n != 6 {
		t.Errorf("page 1 fills %d paths, want 6", n)
	}
}

func TestLinkedImageMissing(t *testing.T) {
	// without the directory of the document, the linked image is not
	// found and its frame stays empty
	opts := testOptions()
	opts.Dir = "."
	res, r := convert(t, "basic.idml", opts)
	found := false
	for _, w := range res.Warnings {
		if strings.Contains(w, "link.png") {
			found = true
		}
	}
	if !found {
		t.Errorf("no warning about the missing link in %q", res.Warnings)
	}
	b2, _ := layer(r.Manifest.Views[0].Pages[1], bdf.RoleBody)
	if n := ops(t, r, b2, bdf.OpImage); n != 1 {
		t.Errorf("page 2 draws %d images, want the embedded one", n)
	}
	// given as a file, it is found by name
	data, err := os.ReadFile("testdata/link.png")
	if err != nil {
		t.Fatal(err)
	}
	opts = testOptions()
	opts.Dir = "."
	opts.Files = conv.FileMap{"uploads/link.png": data}
	res, r = convert(t, "basic.idml", opts)
	b2, _ = layer(r.Manifest.Views[0].Pages[1], bdf.RoleBody)
	if n := ops(t, r, b2, bdf.OpImage); n != 2 {
		t.Errorf("page 2 draws %d images with the file given, want 2", n)
	}
}

func TestVertical(t *testing.T) {
	res, r := convert(t, "vertical.idml", nil)
	v := r.Manifest.Views[0]
	if res.Pages != 2 || v.Direction != bdf.DirectionRTL {
		t.Fatalf("%d pages, direction %q", res.Pages, v.Direction)
	}
	texts := pageTexts(t, r)
	if !strings.Contains(texts[0], "縦書きのテキスト") || !strings.Contains(texts[1], "行頭に来ません") {
		t.Errorf("texts %q", texts)
	}
	// vertical lines are child objects placed with USE_AT
	b1, _ := layer(v.Pages[0], bdf.RoleBody)
	if ops(t, r, b1, bdf.OpUseAt) < 10 {
		t.Errorf("page 1 places %d vertical lines", ops(t, r, b1, bdf.OpUseAt))
	}
	if len(res.Warnings) != 0 {
		t.Errorf("warnings %q", res.Warnings)
	}
}

func TestPagesAndRegistry(t *testing.T) {
	opts := testOptions()
	opts.Pages = conv.PageList(3, 1)
	res, r := convert(t, "basic.idml", opts)
	if res.Pages != 2 {
		t.Fatalf("%d pages", res.Pages)
	}
	texts := pageTexts(t, r)
	if !strings.Contains(texts[0], "- 3 -") || !strings.Contains(texts[1], "- 1 -") {
		t.Errorf("selected pages %q", texts)
	}
	opts.Pages = conv.PageList(4)
	if _, err := ConvertFile("testdata/basic.idml", opts); err == nil || !strings.Contains(err.Error(), "out of range") {
		t.Errorf("page 4: %v", err)
	}
	f, err := conv.DetectFile("testdata/basic.idml")
	if err != nil || f == nil || f.Name != "idml" {
		t.Fatalf("detected %v, %v", f, err)
	}
	if f, _ := conv.DetectFile("../pptx/testdata/basic.pptx"); f != nil && f.Name == "idml" {
		t.Error("a PowerPoint deck detected as IDML")
	}
	cres, err := conv.ConvertFile("testdata/vertical.idml", "", &conv.Options{FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true})
	if err != nil {
		t.Fatal(err)
	}
	if cres.Summary != "2 page(s), 2 embedded font(s)" {
		t.Errorf("summary %q", cres.Summary)
	}
}

func TestStyleChain(t *testing.T) {
	_, r := convert(t, "basic.idml", nil)
	// the title takes its size and color from Title → Centered → Body →
	// NormalParagraphStyle: 24 pt, cyan, bold
	b1, _ := layer(r.Manifest.Views[0].Pages[0], bdf.RoleBody)
	o, err := r.Object(b1)
	if err != nil {
		t.Fatal(err)
	}
	var sizes []any
	has24, has14 := false, false
	o.Walk(func(in bdf.Instr) {
		if in.Op == bdf.OpFont && len(in.Args) == 2 {
			sizes = append(sizes, in.Args[1])
			has24 = has24 || in.Args[1] == float32(24)
			has14 = has14 || in.Args[1] == float32(14)
		}
	})
	if !has24 || !has14 {
		t.Errorf("no 24 pt title or 14 pt heading among sizes %v", sizes)
	}
}

// TestInDesignExports converts documents InDesign itself wrote (see
// testdata/simpleidml/README.md).
func TestInDesignExports(t *testing.T) {
	for _, tc := range []struct {
		name     string
		pages    int
		want     []string // on the pages, in order
		warnings int      // besides overset text, which depends on the fonts
	}{
		{"4-pages.idml", 4, []string{"Lorem ipsum", "", "", ""}, 0},
		{"4-pages-layers-with-guides.idml", 4, []string{"Lorem ipsum", "", "", ""}, 0},
		{"2articles-1photo.idml", 1, []string{"THE HEADLINE HERE"}, 2},
		{"magazineA-courrier-des-lecteurs-3pages.idml", 3, []string{"DES LECTEURS", "", ""}, 0},
		{"magazineA-bloc-notes.idml", 2, []string{"What whales really want", ""}, 1},
		{"interview.idml", 1, []string{"Cabinet De La Pyramide"}, 0},
		{"page-9modules.idml", 1, []string{"HEADLINE"}, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opts := testOptions()
			opts.Dir = "testdata/simpleidml"
			res, r := convert(t, "simpleidml/"+tc.name, opts)
			if res.Pages != tc.pages {
				t.Fatalf("%d pages, want %d", res.Pages, tc.pages)
			}
			var others []string
			for _, w := range res.Warnings {
				if strings.Contains(w, "internal error") {
					t.Error(w)
				}
				if !strings.Contains(w, "overset") {
					others = append(others, w)
				}
			}
			if len(others) != tc.warnings {
				t.Errorf("warnings %q", res.Warnings)
			}
			texts := pageTexts(t, r)
			for i, want := range tc.want {
				if want != "" && (i >= len(texts) || !strings.Contains(texts[i], want)) {
					t.Errorf("page %d lacks %q", i+1, want)
				}
			}
			for i, p := range r.Manifest.Views[0].Pages {
				if p.W <= 0 || p.H <= 0 {
					t.Errorf("page %d: %g × %g", i+1, p.W, p.H)
				}
			}
		})
	}
	// the frame whose insets leave no room shows nothing: its story's
	// text is not on the page
	opts := testOptions()
	opts.Dir = "testdata/simpleidml"
	_, r := convert(t, "simpleidml/magazineA-bloc-notes.idml", opts)
	if texts := pageTexts(t, r); strings.Contains(texts[0], "Le bloc-notes de") {
		t.Error("text of a frame without room for it was drawn")
	}
}
