package kicad

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
)

// testOptions are options that read only the test fonts.
func testOptions() *Options {
	return &Options{FontDirs: []string{testFonts}, NoSystemFonts: true}
}

func readTestdata(t testing.TB, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", filepath.FromSlash(name)))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// convertBytes converts an input and reads the document back.
func convertBytes(t *testing.T, data []byte, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
	res, err := Convert(bytes.NewReader(data), int64(len(data)), opts)
	if err != nil {
		t.Fatal(err)
	}
	return res, reread(t, res.Doc)
}

func reread(t *testing.T, doc *bdf.Document) *bdf.Reader {
	t.Helper()
	var buf bytes.Buffer
	if err := doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// pageText returns the texts of a page, one a line.
func pageText(t *testing.T, r *bdf.Reader, view, page int) string {
	t.Helper()
	var sb strings.Builder
	for _, l := range r.Manifest.Views[view].Pages[page].Layers {
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
		for _, run := range rs {
			sb.WriteString(run.Text + "\n")
		}
	}
	return sb.String()
}

// pageLinks returns the links of a page.
func pageLinks(t *testing.T, r *bdf.Reader, view, page int) []string {
	t.Helper()
	var links []string
	for _, l := range r.Manifest.Views[view].Pages[page].Layers {
		o, err := r.Object(l.Obj)
		if err != nil {
			t.Fatal(err)
		}
		o.Walk(func(in bdf.Instr) {
			if in.Op == bdf.OpLink {
				links = append(links, in.Args[4].(string))
			}
		})
	}
	return links
}

func viewIDs(r *bdf.Reader) []string {
	var ids []string
	for _, v := range r.Manifest.Views {
		ids = append(ids, v.ID)
	}
	return ids
}

// hasLines reports whether every line is a line of the text.
func hasLines(text string, lines ...string) bool {
	all := strings.Split(text, "\n")
	for _, l := range lines {
		if !slices.Contains(all, l) {
			return false
		}
	}
	return true
}

var boardViews = []string{"board-front", "board-back", "layer-F.Cu", "layer-In1.Cu", "layer-In2.Cu", "layer-B.Cu",
	"layer-F.Paste", "layer-B.Paste", "layer-F.SilkS", "layer-B.SilkS", "layer-F.Mask", "layer-B.Mask",
	"layer-Dwgs.User", "layer-F.CrtYd", "layer-F.Fab", "layer-B.Fab"}

// checkDemo checks the conversion of the demo project.
func checkDemo(t *testing.T, res *Result, r *bdf.Reader) {
	t.Helper()
	if len(res.Warnings) > 0 {
		t.Errorf("warnings %q", res.Warnings)
	}
	if res.Sheets != 3 || res.Boards != 1 || res.EmbeddedFonts != 1 {
		t.Errorf("result: %d sheets, %d boards, %d fonts", res.Sheets, res.Boards, res.EmbeddedFonts)
	}
	m := r.Manifest.Meta
	if m.Source != "kicad" || m.DC.Title.First() != "Converter test" || m.DC.Date.First() != "2026-09-28" || m.DC.Publisher.First() != "bdf" {
		t.Errorf("meta %+v", m)
	}
	if ids := viewIDs(r); !slices.Equal(ids, append([]string{"schematic"}, boardViews...)) {
		t.Errorf("views %q", ids)
	}
	sch := r.Manifest.Views[0]
	if sch.Title != "Converter test" || len(sch.Pages) != 3 || sch.TextIndex == "" {
		t.Fatalf("schematic view %q: %d pages", sch.Title, len(sch.Pages))
	}
	// A4 landscape, then the sub-sheet's A5, in points
	for i, want := range [][2]float32{{841.896, 595.296}, {595.296, 419.544}, {595.296, 419.544}} {
		if p := sch.Pages[i]; abs32(p.W-want[0]) > 0.01 || abs32(p.H-want[1]) > 0.01 || len(p.Layers) != 2 || p.Layers[0].Role != bdf.RoleBackground {
			t.Errorf("page %d: %gx%g, %d layers", i+1, p.W, p.H, len(p.Layers))
		}
	}
	root := pageText(t, r, 0, 0)
	if !hasLines(root, "Title: Converter test", "Id: 1/3", "Sheet: /", "File: demo.kicad_sch", "KiCad E.D.A. 9.0",
		"Converter test rev B (demo)", "Channel A", "File: sub.kicad_sch", "RESET D0 x2 and plain", "IN0", "CLOCK_") {
		t.Errorf("root sheet texts:\n%s", root)
	}
	if !strings.Contains(root, "カナとかな、テスト。") {
		t.Error("the kana text is missing")
	}
	// the two instances of the sub-sheet, with their own references
	for i, want := range [][]string{
		{"Id: 2/3", "Sheet: /Channel A/", "Sheet 2 of 3: Channel A", "R101", "R102"},
		{"Id: 3/3", "Sheet: /Channel B/", "Sheet 3 of 3: Channel B", "R201", "R202"},
	} {
		if text := pageText(t, r, 0, i+1); !hasLines(text, want...) {
			t.Errorf("sheet %d texts:\n%s", i+2, text)
		}
	}
	// the sheets' boxes link to their pages
	if links := pageLinks(t, r, 0, 0); !slices.Equal(links, []string{"#page=2", "#page=3"}) {
		t.Errorf("links %q", links)
	}
	front := pageText(t, r, 1, 0)
	if !hasLines(front, "R1", "U1", "bdf KiCad test", "Converter test rev B", "60.0000 mm", "40.0000 mm") {
		t.Errorf("board texts:\n%s", front)
	}
	if silk := pageText(t, r, 1+slices.Index(boardViews, "layer-B.SilkS"), 0); !hasLines(silk, "R10", "BACK SIDE") || strings.Contains(silk, "R1\n") {
		t.Errorf("back silkscreen texts:\n%s", silk)
	}
}

func abs32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

// TestConvertZip converts the zip archive of a project: the web demo's
// way.
func TestConvertZip(t *testing.T) {
	res, r := convertBytes(t, readTestdata(t, "demo.zip"), testOptions())
	checkDemo(t, res, r)
}

// TestConvertFiles converts a project with the files beside it listed:
// the server's way, through the format's registration.
func TestConvertFiles(t *testing.T) {
	files := converter.FileMap{}
	for _, name := range []string{"demo.kicad_sch", "sub.kicad_sch", "demo.kicad_pcb"} {
		files[name] = readTestdata(t, "demo/"+name)
	}
	pro := readTestdata(t, "demo/demo.kicad_pro")
	opts := &converter.Options{FileName: "demo.kicad_pro", Files: files, FontDirs: []string{testFonts}, NoSystemFonts: true}
	res, err := converter.Convert(bytes.NewReader(pro), int64(len(pro)), "kicad", opts)
	if err != nil {
		t.Fatal(err)
	}
	if res.Summary != "3 sheet(s), a board, 1 embedded font(s)" {
		t.Errorf("summary %q", res.Summary)
	}
	checkDemo(t, &Result{Doc: res.Doc, Warnings: res.Warnings, Sheets: 3, Boards: 1, EmbeddedFonts: 1}, reread(t, res.Doc))

	// the schematic, its project file found beside it in a directory
	res, err = converter.ConvertFile(filepath.Join("testdata", "demo", "demo.kicad_sch"), "", &converter.Options{FontDirs: []string{testFonts}, NoSystemFonts: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) > 0 || len(res.Doc.Views) != 1 || len(res.Doc.Views[0].Pages) != 3 {
		t.Errorf("from a directory: %q, %d views", res.Warnings, len(res.Doc.Views))
	}
}

// TestConvertAlone converts a schematic without the files it refers to:
// its root sheet, with a warning.
func TestConvertAlone(t *testing.T) {
	opts := testOptions()
	opts.FileName = "demo.kicad_sch"
	res, r := convertBytes(t, readTestdata(t, "demo/demo.kicad_sch"), opts)
	if res.Sheets != 1 || len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "sub.kicad_sch") {
		t.Errorf("%d sheets, warnings %q", res.Sheets, res.Warnings)
	}
	// without the project file: the project is named after the file
	if text := pageText(t, r, 0, 0); !hasLines(text, "Converter test rev B (demo)", "Id: 1/1") {
		t.Errorf("texts:\n%s", text)
	}
	// a board alone
	res, r = convertBytes(t, readTestdata(t, "demo/demo.kicad_pcb"), testOptions())
	if res.Sheets != 0 || res.Boards != 1 || len(res.Warnings) > 0 || !slices.Equal(viewIDs(r), boardViews) {
		t.Errorf("board: %d sheets, %q, views %q", res.Sheets, res.Warnings, viewIDs(r))
	}
}

func TestConvertOptions(t *testing.T) {
	data := readTestdata(t, "demo.zip")
	for _, tc := range []struct {
		opts  Options
		views []string
		pages int // of the first view
	}{
		{Options{Views: "schematic"}, []string{"schematic"}, 3},
		{Options{Views: "board"}, boardViews, 1},
		{Options{Views: "board", Layers: "board"}, []string{"board-front", "board-back"}, 1},
		{Options{Layers: "layers"}, append([]string{"schematic"}, boardViews[2:]...), 3},
		{Options{Views: "schematic", Pages: converter.Pages{{From: 3, To: 3}, {From: 1, To: 1}}}, []string{"schematic"}, 2},
	} {
		opts := tc.opts
		opts.FontDirs, opts.NoSystemFonts = []string{testFonts}, true
		_, r := convertBytes(t, data, &opts)
		if ids := viewIDs(r); !slices.Equal(ids, tc.views) || len(r.Manifest.Views[0].Pages) != tc.pages {
			t.Errorf("%+v: views %q, %d pages", tc.opts, ids, len(r.Manifest.Views[0].Pages))
		}
	}
	// pages in the order asked for, their links to the pages they land on
	opts := testOptions()
	opts.Views, opts.Pages = "schematic", converter.Pages{{From: 3, To: 3}, {From: 1, To: 1}}
	_, r := convertBytes(t, data, opts)
	if !hasLines(pageText(t, r, 0, 0), "R201") || !slices.Equal(pageLinks(t, r, 0, 1), []string{"#page=1"}) {
		t.Errorf("selected pages: links %q", pageLinks(t, r, 0, 1))
	}
	opts.Title = "Given"
	if res, _ := convertBytes(t, data, opts); res.Doc.Meta.DC.Title.First() != "Given" {
		t.Error("the title was not overridden")
	}
	for _, bad := range []Options{{Views: "pcb"}, {Views: "schematic", Pages: converter.Pages{{From: 4, To: 4}}}} {
		if _, err := Convert(bytes.NewReader(data), int64(len(data)), &bad); err == nil {
			t.Errorf("%+v: no error", bad)
		}
	}
}

// TestConvertFrame converts a project with its own drawing sheet.
func TestConvertFrame(t *testing.T) {
	dir := filepath.Join("testdata", "frame")
	res, err := converter.ConvertFile(filepath.Join(dir, "frame.kicad_pro"), "", &converter.Options{FontDirs: []string{testFonts}, NoSystemFonts: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) > 0 {
		t.Errorf("warnings %q", res.Warnings)
	}
	r := reread(t, res.Doc)
	text := pageText(t, r, 0, 0)
	if !hasLines(text, "Demo project", "Designer: bdf tests", "Sheet 1/1  Custom frame  2026-09-28", "A sheet with the project's own drawing sheet") ||
		strings.Contains(text, "KiCad E.D.A.") {
		t.Errorf("texts:\n%s", text)
	}
	// the drawing sheet missing: the default one, with a warning
	files := converter.FileMap{"frame.kicad_sch": readTestdata(t, "frame/frame.kicad_sch")}
	pro := readTestdata(t, "frame/frame.kicad_pro")
	res, err = converter.Convert(bytes.NewReader(pro), int64(len(pro)), "kicad",
		&converter.Options{FileName: "frame.kicad_pro", Files: files, FontDirs: []string{testFonts}, NoSystemFonts: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "frame.kicad_wks") ||
		!hasLines(pageText(t, reread(t, res.Doc), 0, 0), "Title: Custom frame") {
		t.Errorf("without the drawing sheet: %q", res.Warnings)
	}
}
