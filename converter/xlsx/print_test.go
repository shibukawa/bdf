package xlsx

import (
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
)

// printText converts print.xlsx as its printed pages and returns the views
// with the text of their pages (pages without text are "").
func printText(t *testing.T, opts *Options) (*Result, map[string][]string) {
	t.Helper()
	opts.Views = ViewsPrint
	res, _ := convert(t, "print.xlsx", opts)
	st, err := res.Doc.SearchText()
	if err != nil {
		t.Fatal(err)
	}
	out := map[string][]string{}
	for _, v := range res.Doc.Views {
		out[v.Title] = make([]string, len(v.Pages))
	}
	for _, v := range st.Views {
		for _, pg := range v.Pages {
			out[v.Title][pg.Page-1] = pg.Text
		}
	}
	return res, out
}

func view(t *testing.T, doc *bdf.Document, title string) *bdf.View {
	t.Helper()
	for _, v := range doc.Views {
		if v.Title == title {
			return v
		}
	}
	t.Fatalf("no view %q", title)
	return nil
}

func near(a float32, b float64) bool { return math.Abs(float64(a)-b) < 0.01 }

func TestPrintView(t *testing.T) {
	opts := testOptions()
	opts.PrintTime = time.Date(2026, 9, 1, 9, 30, 0, 0, time.UTC)
	opts.FileName = "print.xlsx"
	res, text := printText(t, opts)
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	for _, v := range res.Doc.Views {
		if v.Kind != bdf.ViewFlow || !strings.HasPrefix(v.ID, "print") {
			t.Errorf("view %s is %s", v.ID, v.Kind)
		}
	}

	// Report: A4 upright; the manual break after row 30 ends the first
	// page; the title rows are on the pages after them; header and footer
	pages := text["Report"]
	v := view(t, res.Doc, "Report")
	if len(pages) != 3 {
		t.Fatalf("Report: %d pages", len(pages))
	}
	if pg := v.Pages[0]; !near(pg.W, 210*mm) || !near(pg.H, 297*mm) || !near(pg.Body.X, 36) || !near(pg.Body.Y, 0.9*72) ||
		!near(pg.Body.W, 210*mm-72) || !near(pg.Body.H, 297*mm-1.7*72) {
		t.Errorf("Report: page %vx%v body %+v", pg.W, pg.H, *pg.Body)
	}
	if !strings.Contains(pages[0], "Item 028") || strings.Contains(pages[0], "Item 029") || !strings.Contains(pages[1], "Item 029") {
		t.Errorf("Report: the manual break is not after row 30:\n%s", pages[0])
	}
	if !strings.Contains(pages[2], "Item 088") {
		t.Errorf("Report: the last page does not end the table:\n%s", pages[2])
	}
	for i, p := range pages {
		if n := strings.Count(p, "Items by line"); n != 1 {
			t.Errorf("Report: page %d has the title rows %d times", i+1, n)
		}
		for _, want := range []string{"本書 report", "Report", "2026/9/1 print.xlsx", string(rune('1'+i)) + " / 3", "Line one", "line two &"} {
			if want == "2026/9/1 print.xlsx" && !strings.Contains(p, want) {
				want = "9/1/2026 print.xlsx" // not a Japanese workbook
			}
			if !strings.Contains(p, want) {
				t.Errorf("Report: page %d lacks %q:\n%s", i+1, want, p)
			}
		}
		if strings.Contains(p, "outside the print area") {
			t.Errorf("Report: page %d shows a cell outside the print area", i+1)
		}
		if roles := roles(v.Pages[i]); !reflect.DeepEqual(roles, []string{bdf.RoleHeader, bdf.RoleFooter, bdf.RoleBody}) {
			t.Errorf("Report: page %d has the layers %v", i+1, roles)
		}
	}
	// the title rows of the second and the third page are one object
	if a, b := titleObject(t, res.Doc, v.Pages[1]), titleObject(t, res.Doc, v.Pages[2]); a != b || a == (bdf.Hash{}) {
		t.Errorf("Report: the title rows of pages 2 and 3 are %s and %s", a, b)
	}

	// Wide: A4 turned, 150%; the columns break after G; the pages go
	// across before they go down; title columns and headings
	pages = text["Wide"]
	v = view(t, res.Doc, "Wide")
	if pg := v.Pages[0]; !near(pg.W, 297*mm) || !near(pg.H, 210*mm) {
		t.Errorf("Wide: page %vx%v", pg.W, pg.H)
	}
	if len(pages) < 4 || len(pages)%2 != 0 && len(pages)%3 != 0 {
		t.Fatalf("Wide: %d pages", len(pages))
	}
	if !strings.Contains(pages[0], "C6") || strings.Contains(pages[0], "C7") || !strings.Contains(pages[1], "C7") {
		t.Errorf("Wide: the columns do not break after G, or the pages do not go across first:\n%s\n--\n%s", pages[0], pages[1])
	}
	for i, p := range pages {
		if !strings.Contains(p, "Row ") {
			t.Errorf("Wide: page %d lacks the title column:\n%s", i+1, p)
		}
		if !strings.Contains(p, "Page "+itoa(i+1)+" of "+itoa(len(pages))) {
			t.Errorf("Wide: page %d lacks its number:\n%s", i+1, p)
		}
	}
	// the headings: the letter of the title column on every page
	if !strings.Contains(pages[1], "A\n") && !strings.Contains(pages[1], "A ") {
		t.Errorf("Wide: page 2 lacks the heading of column A:\n%s", pages[1])
	}

	// Fit: every column on one page across
	pages = text["Fit"]
	if len(pages) == 0 || !strings.Contains(pages[0], "Column 1") || !strings.Contains(pages[0], "Column 12") {
		t.Errorf("Fit: the columns are not on one page: %d pages\n%s", len(pages), pages[0])
	}
	if pg := view(t, res.Doc, "Fit").Pages[0]; !near(pg.W, 612) || !near(pg.H, 792) {
		t.Errorf("Fit: page %vx%v", pg.W, pg.H)
	}

	// Areas: a page for each area, numbered from 5; another header first
	pages = text["Areas"]
	v = view(t, res.Doc, "Areas")
	if len(pages) != 2 {
		t.Fatalf("Areas: %d pages", len(pages))
	}
	if pg := v.Pages[0]; !near(pg.W, 148*mm) || !near(pg.H, 210*mm) {
		t.Errorf("Areas: page %vx%v", pg.W, pg.H)
	}
	for i, want := range [][]string{{"First page", "A1", "C6", "5"}, {"Other pages", "E1", "G6", "6"}} {
		for _, w := range want {
			if !strings.Contains(pages[i], w) {
				t.Errorf("Areas: page %d lacks %q:\n%s", i+1, w, pages[i])
			}
		}
		if strings.Contains(pages[i], "no") {
			t.Errorf("Areas: page %d shows the column between the areas:\n%s", i+1, pages[i])
		}
	}
	if strings.Contains(pages[0], "E1") || strings.Contains(pages[1], "A1") || strings.Contains(pages[1], "First page") {
		t.Errorf("Areas: the areas are not a page each:\n%s\n--\n%s", pages[0], pages[1])
	}

	// Used: what the sheet uses, from A1, with the text that flows out
	pages = text["Used"]
	if len(pages) == 0 || !strings.Contains(pages[0], "From B2") || !strings.Contains(pages[0], "merged") {
		t.Errorf("Used: %d pages\n%s", len(pages), strings.Join(pages, "\n--\n"))
	}
}

func itoa(n int) string { return string(rune('0' + n)) }

func roles(pg *bdf.Page) []string {
	var out []string
	for _, l := range pg.Layers {
		out = append(out, l.Role)
	}
	return out
}

// titleObject returns the first object the body layer of a page uses: the
// title rows of a page that shows them.
func titleObject(t *testing.T, doc *bdf.Document, pg *bdf.Page) bdf.Hash {
	t.Helper()
	for _, l := range pg.Layers {
		if l.Role != bdf.RoleBody {
			continue
		}
		part := doc.Part(l.Obj)
		if part == nil {
			t.Fatal("no body object")
		}
		obj, err := bdf.DecodeObject(part.Data)
		if err != nil {
			t.Fatal(err)
		}
		if len(obj.Objects) > 0 {
			return obj.Objects[0]
		}
	}
	return bdf.Hash{}
}

// TestPrintViews checks the selection of views: the sheets alone by
// default, the pages alone, or both with the sheet first.
func TestPrintViews(t *testing.T) {
	kinds := func(views string) []string {
		opts := testOptions()
		opts.Views = views
		res, _ := convert(t, "print.xlsx", opts)
		var out []string
		for _, v := range res.Doc.Views[:min(len(res.Doc.Views), 2)] {
			out = append(out, v.ID+" "+v.Kind+" "+v.Title)
		}
		return out
	}
	if got, want := kinds(""), []string{"sheet1 sheet Report", "sheet2 sheet Wide"}; !reflect.DeepEqual(got, want) {
		t.Errorf("default views: %q", got)
	}
	if got, want := kinds(ViewsPrint), []string{"print1 flow Report", "print2 flow Wide"}; !reflect.DeepEqual(got, want) {
		t.Errorf("print views: %q", got)
	}
	if got, want := kinds(ViewsBoth), []string{"sheet1 sheet Report", "print1 flow Report (print)"}; !reflect.DeepEqual(got, want) {
		t.Errorf("both views: %q", got)
	}
	if _, err := ConvertFile("testdata/print.xlsx", &Options{Views: "pages"}); err == nil {
		t.Error("an unknown selection of views is accepted")
	}
	if _, err := ConvertFile("testdata/print.xlsx", &Options{Paper: "a9"}); err == nil {
		t.Error("an unknown paper is accepted")
	}
}

// TestPrintSameSheet checks that drawing the pages of a sheet leaves the
// sheet view as it is without them: the same tiles with the same
// instructions (their fonts are other subsets, with the characters of the
// headers and footers).
func TestPrintSameSheet(t *testing.T) {
	alone, _ := convert(t, "print.xlsx", testOptions())
	opts := testOptions()
	opts.Views = ViewsBoth
	both, _ := convert(t, "print.xlsx", opts)
	ops := func(doc *bdf.Document, hash string) (string, []string) {
		h, err := bdf.ParseHash(hash)
		if err != nil {
			t.Fatal(err)
		}
		obj, err := bdf.DecodeObject(doc.Part(h).Data)
		if err != nil {
			t.Fatal(err)
		}
		return string(obj.Ops), obj.Strings
	}
	for _, v := range alone.Doc.Views {
		var w *bdf.View
		for _, x := range both.Doc.Views {
			if x.ID == v.ID {
				w = x
			}
		}
		if w == nil || len(v.Tiles) != len(w.Tiles) || len(v.Tiles) == 0 || !reflect.DeepEqual(v.Cols, w.Cols) || !reflect.DeepEqual(v.Rows, w.Rows) {
			t.Errorf("view %s differs when the pages are drawn too", v.ID)
			continue
		}
		for k, h := range v.Tiles {
			a, as := ops(alone.Doc, h)
			b, bs := ops(both.Doc, w.Tiles[k])
			if a != b || !reflect.DeepEqual(as, bs) {
				t.Errorf("view %s: tile %s differs when the pages are drawn too", v.ID, k)
			}
		}
	}
}

func TestSplitRefs(t *testing.T) {
	for in, want := range map[string][]string{
		"Sheet1!$A$1:$C$5":                      {"$A$1:$C$5"},
		"'My Sheet'!$A$1:$C$5,'My Sheet'!$E:$F": {"$A$1:$C$5", "$E:$F"},
		"'a,b!c'!$1:$2":                         {"$1:$2"},
		"'it''s'!$B$2":                          {"$B$2"},
		"#REF!":                                 nil,
		"":                                      nil,
	} {
		if got := splitRefs(in); !reflect.DeepEqual(got, want) {
			t.Errorf("splitRefs(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestClipSegment(t *testing.T) {
	b := &box{10, 10, 20, 10} // x 10..30, y 10..20
	for _, c := range []struct {
		in   segment
		want segment
		ok   bool
	}{
		{segment{0, 15, 40, 15}, segment{10, 15, 30, 15}, true},  // across
		{segment{30, 0, 30, 40}, segment{30, 10, 30, 20}, true},  // on the right edge
		{segment{31, 0, 31, 40}, segment{}, false},               // beside it
		{segment{0, 20, 40, 20}, segment{10, 20, 30, 20}, true},  // on the bottom edge
		{segment{12, 12, 18, 18}, segment{12, 12, 18, 18}, true}, // inside
		{segment{0, 0, 40, 20}, segment{20, 10, 30, 15}, true},   // diagonal
		{segment{0, 30, 40, 30}, segment{}, false},               // below
		{segment{15, 15, 15, 15}, segment{15, 15, 15, 15}, true}, // a point inside
		{segment{35, 15, 35, 15}, segment{}, false},              // a point outside
		{segment{40, 15, 0, 15}, segment{30, 15, 10, 15}, true},  // right to left
		{segment{30, 20, 50, 40}, segment{30, 20, 30, 20}, true}, // touching a corner
		{segment{31, 20, 50, 40}, segment{}, false},              // leaving from outside
		{segment{10, 0, 10, 10}, segment{10, 10, 10, 10}, true},  // ending on the corner
		{segment{0, 10, 10, 10}, segment{10, 10, 10, 10}, true},  // ending on the corner
		{segment{0, 9.99, 40, 9.99}, segment{}, false},           // just above
		{segment{10, 10, 30, 20}, segment{10, 10, 30, 20}, true}, // corner to corner
		{segment{-10, 15, 5, 15}, segment{}, false},              // left of it
		{segment{20, 5, 20, 25}, segment{20, 10, 20, 20}, true},  // down through it
		{segment{20, 25, 20, 5}, segment{20, 20, 20, 10}, true},  // up through it
		{segment{20, 20, 20, 25}, segment{20, 20, 20, 20}, true}, // leaving from the edge
		{segment{20, 20.01, 20, 25}, segment{}, false},           // leaving from below the edge
	} {
		got, ok := clipSegment(c.in, b)
		if ok != c.ok || ok && (math.Abs(got.x0-c.want.x0) > 1e-9 || math.Abs(got.y0-c.want.y0) > 1e-9 || math.Abs(got.x1-c.want.x1) > 1e-9 || math.Abs(got.y1-c.want.y1) > 1e-9) {
			t.Errorf("clipSegment(%v) = %v, %v; want %v, %v", c.in, got, ok, c.want, c.ok)
		}
	}
}

// TestPrintBands checks how rows are cut into pages: where the next does
// not fit, before manual breaks, with less room under repeated titles.
func TestPrintBands(t *testing.T) {
	rows := newAxis(1000, 10, map[int]float64{3: 35, 20: 0})
	plan := func() *printPlan {
		return &printPlan{ps: &printSetup{}, scale: 1, ph: 100, pw: 100}
	}
	cut := func(p *printPlan, first, last int) []printBand {
		b, _ := p.cut(rows, first, last, true, 0)
		return b
	}
	// 100 high: rows 0-2 (30) and row 3 (35) and rows 4-6 (30) fit; row 7 does not
	p := plan()
	if got, want := cut(p, 0, 15), []printBand{{0, 7}, {7, 16}}; !reflect.DeepEqual(got, want) {
		t.Errorf("bands = %v, want %v", got, want)
	}
	// a hidden row takes no room
	if got, want := cut(p, 12, 23), []printBand{{12, 23}, {23, 24}}; !reflect.DeepEqual(got, want) {
		t.Errorf("bands over a hidden row = %v, want %v", got, want)
	}
	// manual breaks: before rows 2 and 5; one at the first row does nothing
	p.ps.rowBreaks = []int{0, 2, 5, 500}
	if got, want := cut(p, 0, 15), []printBand{{0, 2}, {2, 5}, {5, 15}, {15, 16}}; !reflect.DeepEqual(got, want) {
		t.Errorf("bands with breaks = %v, want %v", got, want)
	}
	// title rows 0-1 (20) leave 80 on the pages after them
	p = plan()
	p.titleRows = &[2]int{0, 1}
	p.ty0, p.ty1 = 0, 20
	if got, want := cut(p, 0, 30), []printBand{{0, 7}, {7, 15}, {15, 24}, {24, 31}}; !reflect.DeepEqual(got, want) {
		t.Errorf("bands under titles = %v, want %v", got, want)
	}
	// at half the scale a page holds twice the rows
	p = plan()
	p.scale = 0.5
	if got, want := cut(p, 4, 40), []printBand{{4, 25}, {25, 41}}; !reflect.DeepEqual(got, want) {
		t.Errorf("bands at 50%% = %v, want %v", got, want)
	}
	// a row taller than the page is a page
	tall := newAxis(10, 10, map[int]float64{1: 300})
	p = plan()
	if got, _ := p.cut(tall, 0, 3, true, 0); !reflect.DeepEqual(got, []printBand{{0, 1}, {1, 2}, {2, 4}}) {
		t.Errorf("bands with a tall row = %v", got)
	}
	// a limit stops the cutting
	if got, ok := plan().cut(rows, 0, 999, true, 3); ok || len(got) != 3 {
		t.Errorf("limited bands = %v, %v", got, ok)
	}
}

func TestBandsIn(t *testing.T) {
	pos := []float64{0, 10, 30, 60}
	for _, c := range []struct {
		a, b   float64
		k0, k1 int
	}{{0, 5, 0, 1}, {5, 15, 0, 2}, {10, 30, 1, 2}, {9.99, 30.01, 0, 3}, {60, 70, 3, 3}, {-5, 0, 0, 0}, {-5, 0.01, 0, 1}, {59, 100, 2, 3}} {
		if k0, k1 := bandsIn(pos, c.a, c.b); k0 != c.k0 || k1 != c.k1 {
			t.Errorf("bandsIn(%v, %v) = %d, %d; want %d, %d", c.a, c.b, k0, k1, c.k0, c.k1)
		}
	}
}

// TestHeaderFooterCodes checks the codes of headers and footers.
func TestHeaderFooterCodes(t *testing.T) {
	c := newConverter(nil, testOptions())
	c.opts.FileName = "book.xlsx"
	c.opts.PrintTime = time.Date(2026, 1, 2, 3, 4, 0, 0, time.UTC)
	c.st = loadStyles(c, nil, func(string) (rgb, bool) { return rgb{}, false })
	c.st.fonts = []xfont{{name: "M PLUS 1p", size: 11}}
	s := &sheetCtx{c: c, ws: &worksheet{name: "Sheet A"}, advances: map[advanceKey]float64{}}
	p := &printPlan{s: s, ps: &printSetup{hf: &headerFooter{scaleWithDoc: true}}, scale: 0.5}
	type run struct {
		text         string
		size         float64
		bold, italic bool
		font         string
	}
	flat := func(code string) [3][]string {
		var out [3][]string
		for k, lines := range p.hfSections(code, 3, 9) {
			for _, runs := range lines {
				var b strings.Builder
				for _, r := range runs {
					b.WriteString(r.text)
				}
				out[k] = append(out[k], b.String())
			}
		}
		return out
	}
	for code, want := range map[string][3][]string{
		"plain":                           {nil, {"plain"}, nil},
		"&Lleft&Ccenter&Rright":           {{"left"}, {"center"}, {"right"}},
		"&P of &N":                        {nil, {"3 of 9"}, nil},
		"&P+2 &P-1 &P+":                   {nil, {"5 2 3+"}, nil},
		"&L&A&R&F":                        {{"Sheet A"}, nil, {"book.xlsx"}},
		"&D &T":                           {nil, {"1/2/2026 03:04"}, nil},
		"a && b":                          {nil, {"a & b"}, nil},
		"one\ntwo\n\nfour":                {nil, {"one", "two", "", "four"}, nil},
		"&\"Arial,Bold\"&12x&Ky&K00FF00z": {nil, {"xyz"}, nil},
		"&Z&G&O&H&Qtext&":                 {nil, {"text&"}, nil},
		"&K01+000themed":                  {nil, {"themed"}, nil},
	} {
		if got := flat(code); !reflect.DeepEqual(got, want) {
			t.Errorf("sections of %q = %q, want %q", code, got, want)
		}
	}
	// the formats of runs: sizes scale with the sheet
	secs := p.hfSections("&\"M PLUS 1p,Bold Italic\"&20big&\"-,Regular\"&Bb&B&In&8small", 1, 1)
	var got []run
	for _, r := range secs[1][0] {
		got = append(got, run{r.text, r.st.size, r.st.bold, r.st.italic, r.st.latin})
	}
	def := s.textStyle(c.st.font(0), nil, nil).latin
	want := []run{{"big", 10, true, true, "M PLUS 1p"}, {"b", 10, true, false, def}, {"n", 10, false, true, def}, {"small", 4, false, true, def}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("runs = %+v, want %+v", got, want)
	}
}

// TestPrintLimits checks that print areas a file of a few bytes states (a
// million rows, hundreds of times) do not make more pages than a sheet may
// have, nor take long to cut into them.
func TestPrintLimits(t *testing.T) {
	var areas []string
	for i := 0; i < 400; i++ {
		areas = append(areas, "Sales!$A$2:$C$1048576")
	}
	names := `<definedNames><definedName name="_xlnm.Print_Area" localSheetId="0">` + strings.Join(areas, ",") + `</definedName>` +
		`<definedName name="_xlnm.Print_Titles" localSheetId="0">Sales!$1:$1048576,Sales!$A:$XFD</definedName></definedNames>`
	opts := testOptions()
	opts.Views = ViewsPrint
	opts.Sheets = []conv.PageSpan{{From: 1, To: 1}}
	start := time.Now()
	res, _ := convertEdited(t, "basic.xlsx", map[string]func(string) string{
		"xl/workbook.xml": func(s string) string { return strings.Replace(s, "</sheets>", "</sheets>"+names, 1) },
		"xl/worksheets/sheet1.xml": func(s string) string {
			s = strings.Replace(s, "<sheetData>", `<sheetPr><pageSetUpPr fitToPage="1"/></sheetPr><sheetData>`, 1)
			return strings.Replace(s, "</sheetData>", `</sheetData><pageSetup fitToWidth="9999999" fitToHeight="9999999" scale="1e300" paperWidth="1e300mm" paperHeight="NaNmm"/>`, 1)
		},
	}, opts)
	if d := time.Since(start); d > 20*time.Second {
		t.Errorf("the conversion took %v", d)
	}
	v := res.Doc.Views[0]
	if len(v.Pages) == 0 || len(v.Pages) > maxPrintPages {
		t.Errorf("%d pages", len(v.Pages))
	}
	found := false
	for _, w := range res.Warnings {
		found = found || strings.Contains(w, "pages")
	}
	if !found {
		t.Errorf("no warning about the pages: %q", res.Warnings)
	}
	for _, pg := range v.Pages[:1] {
		if !(pg.W > 0 && pg.H > 0 && pg.W <= 14400 && pg.H <= 14400) {
			t.Errorf("page %vx%v", pg.W, pg.H)
		}
	}
}
