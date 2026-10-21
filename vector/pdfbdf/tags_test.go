package pdfbdf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/vector/internal/vectortest"
)

// pdfObjects reads the objects of a file written without compression, by
// their numbers.
func pdfObjects(file string) map[int]string {
	out := map[int]string{}
	for _, m := range regexp.MustCompile(`(?ms)^(\d+) 0 obj\n(.*?)\nendobj`).FindAllStringSubmatch(file, -1) {
		n, _ := strconv.Atoi(m[1])
		out[n] = m[2]
	}
	return out
}

var (
	reRef     = regexp.MustCompile(`(\d+) 0 R`)
	reDict    = regexp.MustCompile(`<<.*?>>`)
	reTag     = regexp.MustCompile(`/S /(\w+)`)
	reKids    = regexp.MustCompile(`/K \[(.*)\]`)
	reMCID    = regexp.MustCompile(`/MCID (\d+)`)
	reTree    = regexp.MustCompile(`/StructTreeRoot (\d+) 0 R`)
	reAlt     = regexp.MustCompile(`/(Alt|Lang) \((.*?)\)`)
	reAttrs   = regexp.MustCompile(`/A << /O /Table(.*?) >>`)
	reParents = regexp.MustCompile(`/StructParents (\d+)`)
)

// structTree writes the structure tree of a file as lines: the elements
// indented under those they are in, each with what it says of itself and
// the number of its pieces of content.
func structTree(t *testing.T, file string) []string {
	t.Helper()
	objs := pdfObjects(file)
	m := reTree.FindStringSubmatch(file)
	if m == nil {
		t.Fatal("the file has no structure tree")
	}
	root, _ := strconv.Atoi(m[1])
	var out []string
	var walk func(n, depth int, parent int)
	walk = func(n, depth, parent int) {
		o := objs[n]
		if !strings.Contains(o, fmt.Sprintf("/P %d 0 R", parent)) {
			t.Errorf("element %d does not name its parent %d: %s", n, parent, o)
		}
		line := strings.Repeat("  ", depth) + reTag.FindStringSubmatch(o)[1]
		for _, a := range reAlt.FindAllStringSubmatch(o, -1) {
			line += fmt.Sprintf(" %s=%q", a[1], a[2])
		}
		if a := reAttrs.FindStringSubmatch(o); a != nil {
			line += " [" + strings.TrimSpace(a[1]) + "]"
		}
		kids := reKids.FindStringSubmatch(o)[1]
		content := len(reMCID.FindAllString(kids, -1))
		bare := reDict.ReplaceAllString(kids, "")
		bare = reRef.ReplaceAllString(bare, "")
		content += len(strings.Fields(bare))
		if content > 0 {
			line += fmt.Sprintf(" (%d)", content)
		}
		if strings.Contains(kids, "/OBJR") {
			line += " +annotation"
		}
		out = append(out, line)
		for _, k := range reRef.FindAllStringSubmatch(reDict.ReplaceAllString(kids, ""), -1) {
			c, _ := strconv.Atoi(k[1])
			walk(c, depth+1, n)
		}
	}
	for _, k := range reRef.FindAllStringSubmatch(reKids.FindStringSubmatch(objs[root])[1], -1) {
		c, _ := strconv.Atoi(k[1])
		walk(c, 0, root)
	}
	return out
}

// checkMarked checks the marked content of every stream of a file: its
// sequences are closed, none is in another but those of alternative texts,
// and each piece of content has an element in the parent tree.
func checkMarked(t *testing.T, file string) (artifacts, pieces int) {
	t.Helper()
	objs := pdfObjects(file)
	root, _ := strconv.Atoi(reTree.FindStringSubmatch(file)[1])
	tree, _ := strconv.Atoi(regexp.MustCompile(`/ParentTree (\d+) 0 R`).FindStringSubmatch(objs[root])[1])
	nums := map[int]string{}
	for _, m := range regexp.MustCompile(`(\d+) (\d+) 0 R`).FindAllStringSubmatch(objs[tree], -1) {
		key, _ := strconv.Atoi(m[1])
		ref, _ := strconv.Atoi(m[2])
		nums[key] = objs[ref]
	}
	contents := map[int]int{} // the content stream of a page, by the page
	for n, o := range objs {
		if m := regexp.MustCompile(`/Type /Page .*/Contents (\d+) 0 R`).FindStringSubmatch(o); m != nil {
			c, _ := strconv.Atoi(m[1])
			contents[c] = n
		}
	}
	for n, o := range objs {
		i := strings.Index(o, "stream\n")
		if i < 0 || !strings.Contains(o[i:], "BDC") && !strings.Contains(o[i:], "BMC") {
			continue
		}
		head, body := o[:i], o[i:]
		// the key of the stream: of its form, or of the page it is the content of
		key := -1
		if m := reParents.FindStringSubmatch(head); m != nil {
			key, _ = strconv.Atoi(m[1])
		} else if pg, ok := contents[n]; ok {
			if m := reParents.FindStringSubmatch(objs[pg]); m != nil {
				key, _ = strconv.Atoi(m[1])
			}
		}
		depth, tagged := 0, 0
		for _, line := range strings.Split(body, "\n") {
			switch {
			case strings.HasSuffix(line, "BDC") || strings.HasSuffix(line, "BMC"):
				if depth > 0 && !strings.Contains(line, "/ActualText") {
					t.Errorf("object %d: marked content in marked content: %s", n, line)
				}
				depth++
				if line == "/Artifact BMC" {
					artifacts++
				}
				if m := reMCID.FindStringSubmatch(line); m != nil {
					id, _ := strconv.Atoi(m[1])
					if id != tagged {
						t.Errorf("object %d: the content %d comes as the %d", n, id, tagged)
					}
					tagged++
				}
			case line == "EMC":
				if depth--; depth < 0 {
					t.Errorf("object %d: EMC without marked content", n)
				}
			}
		}
		if depth != 0 {
			t.Errorf("object %d: %d sequences of marked content are left open", n, depth)
		}
		pieces += tagged
		if tagged == 0 {
			continue
		}
		arr, ok := nums[key]
		if !ok {
			t.Errorf("object %d: its marked content has no entry (key %d) in the parent tree", n, key)
			continue
		}
		if got := len(reRef.FindAllString(arr, -1)); got != tagged {
			t.Errorf("object %d: %d pieces of content, %d elements in the parent tree", n, tagged, got)
		}
	}
	return artifacts, pieces
}

// TestTagged writes a tagged file: the structure tree of the MARKs of the
// document, with the content of the pages marked as that of its elements
// and the rest as artifacts.
func TestTagged(t *testing.T) {
	font, err := os.ReadFile(filepath.Join(vectortest.Fonts[0], "MPLUS1p-Regular-subset.ttf"))
	if err != nil {
		t.Fatal(err)
	}
	doc := bdf.NewDocument()
	doc.Meta.DC.Title = []string{"A tagged document"}
	doc.Meta.DC.Language = []string{"en"}
	fh := doc.AddFont(font)
	// a master that two pages share, and the body of each
	master := bdf.NewObject()
	master.FillColor(bdf.RGB(230, 236, 245)).FillRect(0, 0, 300, 20)
	master.Font(master.AddFont(bdf.EmbeddedFont(fh, 400, bdf.StyleNormal)), 8).FillColor(bdf.RGB(0, 0, 0)).FillText("Master", 10, 14, 0)
	master.SetBBox(0, 0, 300, 400)
	mh, _ := doc.AddObject(master)

	body := bdf.NewObject()
	f := body.AddFont(bdf.EmbeddedFont(fh, 400, bdf.StyleNormal))
	body.Font(f, 12).FillColor(bdf.RGB(0, 0, 0))
	text := func(o *bdf.Object, s string, y float32) { o.FillText(s, 20, y, 0) }
	body.Mark(bdf.MarkHeading, "2")
	text(body, "Heading", 40)
	// (a rule under the heading: not content)
	body.FillRect(20, 44, 200, 1)
	body.Mark(bdf.MarkParagraph, "")
	text(body, "First line", 60)
	body.Mark(bdf.MarkLine, "")
	text(body, "second line", 74)
	body.Link(20, 64, 60, 12, "https://example.com/")
	body.Mark(bdf.MarkLang, "ja")
	text(body, "abc", 88)
	body.Mark(bdf.MarkLang, "")
	body.Mark(bdf.MarkList, "")
	body.Mark(bdf.MarkListItem, "")
	text(body, "one", 110)
	body.Mark(bdf.MarkListItem, "")
	text(body, "two", 124)
	body.Mark(bdf.MarkEnd, "")
	body.Mark(bdf.MarkTable, "")
	text(body, "Caption", 150)
	body.Mark(bdf.MarkCell, "A1 col")
	text(body, "Name", 166)
	body.Mark(bdf.MarkCell, "B1:C1 col")
	text(body, "Values", 166)
	body.Mark(bdf.MarkCell, "A2:A3 row")
	text(body, "x", 180)
	body.Mark(bdf.MarkCell, "C2")
	text(body, "1", 180)
	body.Mark(bdf.MarkCell, "B3")
	text(body, "2", 194)
	body.Mark(bdf.MarkEnd, "")
	body.Mark(bdf.MarkFigure, "A red square")
	body.FillColor(bdf.RGB(200, 0, 0)).FillRect(20, 210, 30, 30)
	body.FillColor(bdf.RGB(0, 0, 0))
	text(body, "in the figure", 236)
	body.Mark(bdf.MarkEnd, "")
	body.Mark(bdf.MarkFigure, "")
	body.FillRect(100, 210, 30, 30) // decoration
	body.Mark(bdf.MarkEnd, "")
	// text in a group: marked in the form of the group
	body.Mark(bdf.MarkBox, "")
	body.GroupBegin(0.5, bdf.BlendSourceOver, 0, 250, 300, 40)
	text(body, "In a group", 270)
	body.GroupEnd()
	// a drawing that stands for text
	glyph := bdf.NewObject()
	glyph.FillRect(0, 0, 8, 8).FillRect(10, 0, 8, 8)
	glyph.SetBBox(0, 0, 20, 10)
	gh, gbox := doc.AddObject(glyph)
	body.Mark(bdf.MarkParagraph, "")
	body.Mark(bdf.MarkAltText, "drawn")
	body.UseAt(body.AddObject(gh, gbox), 20, 300)
	// cells of a sheet, outside a table
	body.Mark(bdf.MarkCell, "B5")
	text(body, "b5", 330)
	body.Mark(bdf.MarkCell, "D5")
	text(body, "d5", 330)
	body.Mark(bdf.MarkCell, "C6")
	text(body, "c6", 344)
	body.SetBBox(0, 0, 300, 400)
	bh, _ := doc.AddObject(body)

	second := bdf.NewObject()
	second.Font(second.AddFont(bdf.EmbeddedFont(fh, 400, bdf.StyleNormal)), 12).FillColor(bdf.RGB(0, 0, 0)).FillText("Page two", 20, 40, 0)
	second.SetBBox(0, 0, 300, 400)
	sh, _ := doc.AddObject(second)
	v := doc.NewView("pages", bdf.ViewFixed, "")
	for _, h := range []bdf.Hash{bh, sh} {
		p := v.AddPage(300, 400)
		p.Layers = append(p.Layers, bdf.Layer{Role: bdf.RoleMaster, Obj: mh}, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	}

	o := testOptions()
	o.NoCompression, o.Tagged = true, true
	data, res := write(t, doc, o)
	file := string(data)
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	for _, want := range []string{"/MarkInfo << /Marked true >>", "/Lang (en)", "/DisplayDocTitle true", "/Tabs /S", "/Type /StructTreeRoot", "/ParentTreeNextKey"} {
		if !strings.Contains(file, want) {
			t.Errorf("the file lacks %q", want)
		}
	}
	got := structTree(t, file)
	want := []string{
		"Document",
		"  H2 (1)",
		"  P (1)", // (its two lines are one piece of content)
		"    Link +annotation",
		"    Span Lang=\"ja\" (1)",
		"  L",
		"    LI",
		"      LBody",
		"        P (1)",
		"    LI",
		"      LBody",
		"        P (1)",
		"  Table",
		"    Caption (1)",
		"    TR",
		"      TH [/Scope /Column]",
		"        P (1)",
		"      TH [/ColSpan 2 /Scope /Column]",
		"        P (1)",
		"    TR",
		"      TH [/RowSpan 2 /Scope /Row]",
		"        P (1)",
		"      TD", // (a cell that the document has no text for)
		"      TD",
		"        P (1)",
		"    TR",
		"      TD",
		"        P (1)",
		"      TD",
		"  Figure Alt=\"A red square\" (1)",
		"  P (1)", // in the group
		"  P (1)", // the drawing
		"  Table", // the cells of the sheet
		"    TR",
		"      TD (1)",
		"      TD",
		"      TD (1)",
		"    TR",
		"      TD",
		"      TD (1)",
		"      TD",
		"  P (1)", // the second page
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("the structure tree:\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	artifacts, pieces := checkMarked(t, file)
	if artifacts < 4 || pieces != 18 {
		t.Errorf("%d artifacts and %d pieces of content", artifacts, pieces)
	}
	// the master is an artifact where it is drawn; the content of the group is marked in its form
	if !strings.Contains(file, "/Artifact BMC\nq 1 0 0 -1 0 400 cm /Fm") {
		t.Error("the form of the master is not drawn as an artifact")
	}
	if n := len(regexp.MustCompile(`/Subtype /Form[^\n]*/StructParents \d+`).FindAllString(file, -1)); n != 1 {
		t.Errorf("%d forms with marked content, want that of the group", n)
	}
	if !regexp.MustCompile(`/Type /MCR /Pg \d+ 0 R /Stm \d+ 0 R /MCID 0`).MatchString(file) {
		t.Error("the content of the group is not that of its form")
	}
	// the drawing that stands for text, with the text
	if !regexp.MustCompile(`/P << /MCID \d+ >> BDC\n/Span << /ActualText \(drawn\) >> BDC`).MatchString(file) {
		t.Error("the drawing is not the content of its paragraph, with its text")
	}
	// the link: an annotation of the tree
	if !regexp.MustCompile(`/Subtype /Link [^\n]*/Contents \(https://example.com/\) /StructParent \d+`).MatchString(file) {
		t.Error("the link is not of the structure tree")
	}
	// the pages show what those of a file without tags show
	o.Tagged = false
	plain, _ := write(t, doc, o)
	for i := 0; i < 2; i++ {
		if d := vectortest.Distance(render(t, back(t, data), 0, i), render(t, back(t, plain), 0, i)); d != 0 {
			t.Errorf("page %d differs by %.3f from that of the file without tags", i+1, d)
		}
	}
	if strings.Contains(string(plain), "BDC") && !strings.Contains(string(plain), "/ActualText") || strings.Contains(string(plain), "StructTreeRoot") {
		t.Error("a file without tags has marked content")
	}
	// the same file every time
	again, _ := write(t, doc, &Options{FontDirs: vectortest.Fonts, NoSystemFonts: true, NoCompression: true, Tagged: true})
	if !bytes.Equal(again, data) {
		t.Error("the tagged file is not the same when it is written again")
	}
}

// TestTaggedDocuments writes the test documents of the converters as
// tagged files: every stream has its marked content closed and in the
// parent tree, and the kinds of elements are those of the documents.
func TestTaggedDocuments(t *testing.T) {
	for name, want := range map[string][]string{
		"docx/basic.bdf":       {"H1", "P", "L", "LI", "LBody", "Table", "TR", "TH", "TD", "Figure", "Link"},
		"pptx/basic.bdf":       {"P", "H1"},
		"xlsx/print-pages.bdf": {"Table", "TR", "TD"},
		"form/pptx-rows.bdf":   {"Table", "TR", "TD", "P"},
		"drawio/showcase.bdf":  {"P"},
		"pdf/tagged.bdf":       {"P"},
		"demo.bdf":             {"P"},
	} {
		if _, err := os.Stat(filepath.Join("../../testdata", name)); err != nil {
			continue // (a document this build has no fixture of)
		}
		doc := vectortest.Load(t, name)
		o := testOptions()
		o.NoCompression, o.Tagged = true, true
		data, _ := write(t, doc, o)
		file := string(data)
		tree := strings.Join(structTree(t, file), "\n")
		for _, tag := range want {
			if !regexp.MustCompile(`(?m)^ *` + tag + `\b`).MatchString(tree) {
				t.Errorf("%s: no %s in the structure tree", name, tag)
			}
		}
		if _, pieces := checkMarked(t, file); pieces == 0 {
			t.Errorf("%s: no marked content", name)
		}
	}
}

// TestPDFA writes PDF/A files: with the output intent of sRGB and the
// name of the standard, and without what the standard has no place for.
func TestPDFA(t *testing.T) {
	doc := vectortest.Load(t, "form/pptx-fixed.bdf")
	o := testOptions()
	o.NoCompression = true
	plain, _ := write(t, doc, o)
	if !strings.Contains(string(plain), "/Interpolate true") {
		t.Fatal("the document of the test has no image that is smoothed")
	}
	for level, want := range map[string][2]string{
		"2b": {"2", "B"}, "2u": {"2", "U"}, "2a": {"2", "A"}, "3b": {"3", "B"}, "3a": {"3", "A"}, "PDF/A-2B": {"2", "B"}, "pdfa": {"2", "B"},
	} {
		o.PDFA = level
		data, res := write(t, doc, o)
		file := string(data)
		if len(res.Warnings) != 0 {
			t.Errorf("%s: warnings: %q", level, res.Warnings)
		}
		for _, s := range []string{
			"<pdfaid:part>" + want[0] + "</pdfaid:part>", "<pdfaid:conformance>" + want[1] + "</pdfaid:conformance>",
			"/OutputIntents [<< /Type /OutputIntent /S /GTS_PDFA1", "/DestOutputProfile", "%PDF-1.7\n%\xe2\xe3\xcf\xd3\n", "/ID [<",
		} {
			if !strings.Contains(file, s) {
				t.Errorf("%s: the file lacks %q", level, s)
			}
		}
		if strings.Contains(file, "/Interpolate") {
			t.Errorf("%s: an image is to be smoothed", level)
		}
		// level a is tagged
		if tagged := strings.Contains(file, "/StructTreeRoot"); tagged != (want[1] == "A") {
			t.Errorf("%s: tagged = %v", level, tagged)
		}
		// the pages are those of the file that is no PDF/A (but for the smoothing of images)
		if d := vectortest.Distance(render(t, back(t, data), 0, 0), render(t, back(t, plain), 0, 0)); d > 0.5 {
			t.Errorf("%s: the page differs by %.3f", level, d)
		}
	}
	o.PDFA = "1b"
	if _, err := Write(&bytes.Buffer{}, doc, o); err == nil || !strings.Contains(err.Error(), "1b") {
		t.Errorf("PDF/A-1b: %v", err)
	}

	// the profile of the output intent: an ICC profile of a display, of RGB, with its tags inside it
	p := srgbProfile()
	if int(binary.BigEndian.Uint32(p)) != len(p) || string(p[12:24]) != "mntrRGB XYZ " || string(p[36:40]) != "acsp" || len(p)%4 != 0 {
		t.Fatalf("the header of the profile: % x", p[:40])
	}
	n := int(binary.BigEndian.Uint32(p[128:]))
	tags := map[string]bool{}
	for i := 0; i < n; i++ {
		e := p[132+12*i:]
		off, size := int(binary.BigEndian.Uint32(e[4:])), int(binary.BigEndian.Uint32(e[8:]))
		if off < 132+12*n || off+size > len(p) || off%4 != 0 {
			t.Errorf("the tag %s is at %d+%d of %d", e[:4], off, size, len(p))
		}
		tags[string(e[:4])] = true
	}
	for _, sig := range []string{"desc", "cprt", "wtpt", "rXYZ", "gXYZ", "bXYZ", "rTRC", "gTRC", "bTRC"} {
		if !tags[sig] {
			t.Errorf("the profile has no %s", sig)
		}
	}
	if !bytes.Equal(p, srgbProfile()) {
		t.Error("the profile is not the same every time")
	}

	// a character that no font has is left out of a PDF/A file, with a warning
	font, err := os.ReadFile(filepath.Join(vectortest.Fonts[0], "MPLUS1p-Regular-subset.ttf"))
	if err != nil {
		t.Fatal(err)
	}
	missing := vectortest.One(t, 200, 100, func(doc *bdf.Document, o *bdf.Object) {
		o.Font(o.AddFont(bdf.EmbeddedFont(doc.AddFont(font), 400, bdf.StyleNormal)), 20).FillText("AกB", 10, 50, 0)
	})
	o = testOptions()
	o.NoCompression, o.PDFA = true, "2u"
	data, res := write(t, missing, o)
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "no font has") {
		t.Errorf("warnings: %q", res.Warnings)
	}
	if regexp.MustCompile(`<(?:[0-9A-F]{4})*0000(?:[0-9A-F]{4})*> Tj`).MatchString(string(data)) {
		t.Error("the glyph of a character that the font lacks is written")
	}
	if n := strings.Count(string(data), " Tj"); n != 2 {
		t.Errorf("%d runs of text, want one on each side of the character", n)
	}
	// (private use text at level a: said, not hidden)
	pua := vectortest.One(t, 200, 100, func(doc *bdf.Document, o *bdf.Object) {
		o.Font(o.AddFont(bdf.EmbeddedFont(doc.AddFont(font), 400, bdf.StyleNormal)), 20).FillText("A", 10, 50, 0)
	})
	o.PDFA = "2a"
	if _, res := write(t, pua, o); !strings.Contains(strings.Join(res.Warnings, "\n"), "private use") {
		t.Errorf("warnings: %q", res.Warnings)
	}
}
