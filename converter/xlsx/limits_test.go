package xlsx

import (
	"bytes"
	"math"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// TestFormulaDepth checks that a formula nested past the depth limit is left
// unevaluated (with a warning) rather than overflowing the stack, in both
// the parser (nested parentheses) and the evaluator (a long operator chain),
// while a formula just under the limit still works.
func TestFormulaDepth(t *testing.T) {
	var warnings []string
	c := &converter{parsed: map[string]*worksheet{}, warned: map[string]bool{}, opts: &Options{Warn: func(w string) { warnings = append(warnings, w) }}}
	s := &sheetCtx{c: c, ws: &worksheet{name: "S"}, formulas: map[string]fnode{}}

	deep := []string{
		strings.Repeat("(", maxFormulaDepth+1) + "1" + strings.Repeat(")", maxFormulaDepth+1), // parser recursion
		"1" + strings.Repeat("+1", maxFormulaDepth+1),                                         // evaluator recursion (a left-deep tree)
	}
	for _, f := range deep {
		if _, ok := s.evalFormula(f, [2]int{}, 0, 0); ok {
			t.Errorf("a formula of depth %d was evaluated", maxFormulaDepth+1)
		}
	}
	if len(warnings) != len(deep) {
		t.Errorf("warnings = %q", warnings)
	}

	// a formula comfortably under the limit is still evaluated
	near := strings.Repeat("(", maxFormulaDepth/2) + "2+3" + strings.Repeat(")", maxFormulaDepth/2)
	if v, ok := s.evalFormula(near, [2]int{}, 0, 0); !ok || toStr(v) != "5" {
		t.Errorf("nested formula = %q %v", toStr(v), ok)
	}
}

// Wildcards match Unicode characters, with * matching any number and ? one.
func TestWildMatch(t *testing.T) {
	for _, c := range []struct {
		pattern, value string
		want           bool
	}{
		{"", "", true}, {"", "a", false}, {"*", "日本語", true},
		{"a?c", "abc", true}, {"a?c", "ac", false},
		{"a*c", "abbbc", true}, {"a*c", "abbd", false},
		{"a**b", "axxb", true}, {"?", "あ", true}, {"?", "", false},
	} {
		if got := wildMatch(c.pattern, c.value); got != c.want {
			t.Errorf("wildMatch(%q, %q) = %v, want %v", c.pattern, c.value, got, c.want)
		}
	}
}

// TestMergeExtent checks that a merged range of more cells than merges may
// have does not enlarge the view: a merge that spans the whole sheet leaves
// the view the size the cells call for, so a few bytes cannot make a grid
// of a million rows. A title merged over more columns than have cells does.
func TestMergeExtent(t *testing.T) {
	ws := &worksheet{rows: []row{{idx: 0, cells: []cell{{col: 0, kind: cellNum, num: 1}}}}}
	base := &sheetCtx{c: &converter{}, ws: ws, comments: map[[2]int]bool{}}
	base.extent()

	ws.merges = []cellRange{{0, 0, maxRows - 1, maxCols - 1}}
	s := &sheetCtx{c: &converter{}, ws: ws, comments: map[[2]int]bool{}}
	s.extent()
	if s.nRows != base.nRows || s.nCols != base.nCols {
		t.Errorf("a whole-sheet merge made the view %d×%d, not %d×%d", s.nRows, s.nCols, base.nRows, base.nCols)
	}
	if s.nRows > minRows+maxCols || s.nCols > minCols+maxCols {
		t.Errorf("view unexpectedly large: %d×%d", s.nRows, s.nCols)
	}
	ws.merges = []cellRange{{0, 0, 0, 39}, {0, 0, maxRows - 1, maxCols - 1}}
	s = &sheetCtx{c: &converter{}, ws: ws, comments: map[[2]int]bool{}}
	s.extent()
	if s.nCols < 40 || s.nRows != base.nRows {
		t.Errorf("a title merged over 40 columns gives a view of %d×%d", s.nRows, s.nCols)
	}
}

// TestSizeClamps checks the row-height and column-width clamps: non-finite
// or out-of-range values are brought back to Excel's own limits.
func TestSizeClamps(t *testing.T) {
	for _, c := range []struct {
		in, want float64
	}{
		{100, 100}, {maxRowPt, maxRowPt}, {500, maxRowPt}, {-1, -1},
		{math.Inf(1), -1}, {math.NaN(), -1},
	} {
		if got := clampRowHt(c.in); got != c.want && !(math.IsNaN(c.in) && got == -1) {
			t.Errorf("clampRowHt(%v) = %v, want %v", c.in, got, c.want)
		}
	}
	for _, c := range []struct {
		in, want float64
	}{
		{50, 50}, {maxColChar, maxColChar}, {300, maxColChar}, {-5, 0},
		{math.Inf(1), maxColChar}, {math.NaN(), 0},
	} {
		if got := clampColWidth(c.in); got != c.want {
			t.Errorf("clampColWidth(%v) = %v, want %v", c.in, got, c.want)
		}
	}
	// end to end: a row height and a column width past the limits are drawn
	// at the limits (a fresh row 50 keeps its height from being overridden by
	// the workbook's own row 1, and column A is widened past the maximum)
	_, r := convertEdited(t, "basic.xlsx", map[string]func(string) string{
		"xl/worksheets/sheet1.xml": func(s string) string {
			s = strings.Replace(s, "<sheetData>", `<sheetData><row r="50" ht="900" customHeight="1"><c r="A50"><v>1</v></c></row>`, 1)
			return strings.Replace(s, `<col width="14" customWidth="1" min="1" max="1"/>`, `<col width="500" customWidth="1" min="1" max="1"/>`, 1)
		},
	}, testOptions())
	for _, h := range sizes(r.Manifest.Views[0].Rows) {
		if float64(h) > maxRowPt+0.01 {
			t.Errorf("row height %v exceeds the limit %v", h, maxRowPt)
		}
	}
	// column A at 255 characters is under ~2000 pt; unclamped, 500 characters
	// would be well over it
	if w := sizes(r.Manifest.Views[0].Cols); w[0] > 2000 {
		t.Errorf("column A width %v not clamped", w[0])
	}
}

// TestNewlineLines checks that a wrapped cell of nothing but line breaks is
// laid out with at most maxShown lines: the shown-character limit counts the
// breaks too, so a value of pure breaks cannot make an unbounded number of
// lines.
func TestNewlineLines(t *testing.T) {
	c := newConverter(nil, testOptions())
	c.st = loadStyles(c, nil, defaultThemeColor)
	c.st.fonts = []xfont{{name: "Calibri", size: 11}} // a plain font, as a grid uses (no theme)
	c.eaScript = "Jpan"
	c.mdw = c.maxDigitWidth(c.st.font(0))
	s := &sheetCtx{c: c, base: map[int]*cellFmt{}, advances: map[advanceKey]float64{}}
	f := *s.fmtOf(0)
	f.align.wrap = true

	n := maxShown * 2
	cl := &cell{kind: cellStr, text: &richText{plain: strings.Repeat("\n", n)}}
	lay := s.layoutCell(0, 0, cl, &f, box{0, 0, 50, 1e6}, false)
	if lay == nil {
		t.Fatal("no layout")
	}
	if len(lay.lines) > maxShown+2 {
		t.Errorf("%d lines from a cell of %d breaks (limit %d)", len(lay.lines), n, maxShown)
	}
}

// TestFormatTokens checks that a number format of absurdly many placeholders
// is truncated, so formatting a cell with it stays cheap, and that a normal
// format is untouched.
func TestFormatTokens(t *testing.T) {
	nf := parseNumFormat("0."+strings.Repeat("0", 5000), defaultPalette)
	if n := len(nf.sections[0].toks); n > maxSectionTokens {
		t.Errorf("section kept %d tokens (limit %d)", n, maxSectionTokens)
	}
	// it still formats a value without panicking
	if nf.formatNumber(1.5, false).text() == "" {
		t.Error("the truncated format produced nothing")
	}
	normal := parseNumFormat("#,##0.00", defaultPalette)
	if got := normal.formatNumber(1234.5, false).text(); got != "1,234.50" {
		t.Errorf("normal format = %q", got)
	}
}

// TestCFBudget checks that conditional format evaluation stops after its cell
// budget and warns once, rather than reading a whole wide range per rule.
func TestCFBudget(t *testing.T) {
	c := &converter{warned: map[string]bool{}, opts: &Options{}}
	s := &sheetCtx{c: c, ws: &worksheet{name: "S"}, nRows: 1000, nCols: 1000, cfCells: 100}
	out := s.rangeValues([]cellRange{{0, 0, 999, 999}})
	if len(out) != 100 {
		t.Errorf("read %d cells, want the budget of 100", len(out))
	}
	if len(c.warnings) != 1 || !strings.Contains(c.warnings[0], "cover more than") {
		t.Errorf("warnings = %q", c.warnings)
	}
}

// TestMeasureCache checks that caching the character advance keeps it exact:
// the same face at twice the size advances twice as far (the size is not part
// of the cache key, the per-em advance is).
func TestMeasureCache(t *testing.T) {
	c := newConverter(nil, testOptions())
	c.st = loadStyles(c, nil, defaultThemeColor)
	c.eaScript = "Jpan"
	s := &sheetCtx{c: c, advances: map[advanceKey]float64{}}
	st := &tstyle{latin: "Calibri", size: 10}
	big := &tstyle{latin: "Calibri", size: 20}
	small := s.item(st, 'M')
	large := s.item(big, 'M')
	if small.w == 0 || math.Abs(large.w-2*small.w) > 1e-9 {
		t.Errorf("advance at 10pt %v, at 20pt %v", small.w, large.w)
	}
	if len(s.advances) != 1 { // one advance shared by both sizes
		t.Errorf("cache holds %d entries, want 1", len(s.advances))
	}
}

// layoutContext is a sheet context that can be laid out: a converter with
// the test fonts and the formats of a grid (a plain font, no theme).
func layoutContext(ws *worksheet, drawings ...*anchored) *sheetCtx {
	c := newConverter(nil, testOptions())
	c.st = loadStyles(c, nil, defaultThemeColor)
	c.st.fonts = []xfont{{name: "Calibri", size: 11}}
	c.eaScript = "Jpan"
	c.mdw = c.maxDigitWidth(c.st.font(0))
	return &sheetCtx{c: c, ws: ws, tiles: map[[2]int]*tileCv{}, base: map[int]*cellFmt{}, comments: map[[2]int]bool{},
		advances: map[advanceKey]float64{}, drawings: drawings}
}

func drawingWarnings(c *converter) int {
	n := 0
	for _, w := range c.warnings {
		if strings.Contains(w, "drawings cover more than") {
			n++
		}
	}
	return n
}

// TestDrawingTiles checks the budget of tiles of the drawings: a drawing
// anchored over the whole sheet is left out, with a warning, and does not
// make the view larger; a small one anchored far away stays, and the view
// reaches it; and the tiles of the drawings count together.
func TestDrawingTiles(t *testing.T) {
	sheet := func() *worksheet {
		return &worksheet{name: "S", baseColW: 8, rows: []row{{idx: 0, ht: -1, style: -1, cells: []cell{{col: 0, kind: cellNum, num: 1}}}}}
	}
	whole := &anchored{kind: "twoCellAnchor", to: anchorPt{col: maxCols - 1, row: maxRows - 1}, toRow: maxRows - 1, toCol: maxCols - 1}
	far := &anchored{kind: "twoCellAnchor", from: anchorPt{col: 2, row: 900000}, to: anchorPt{col: 4, row: 900003}, toRow: 900003, toCol: 4}
	s := layoutContext(sheet(), whole, far)
	s.layout()
	if !whole.out || far.out {
		t.Errorf("left out: the drawing over the sheet %v, the small one far away %v", whole.out, far.out)
	}
	if s.nRows != far.toRow+1+marginRows || s.nCols != minCols {
		t.Errorf("view of %d×%d", s.nRows, s.nCols)
	}
	if n := drawingWarnings(s.c); n != 1 {
		t.Errorf("%d warnings: %q", n, s.c.warnings)
	}
	if n := s.drawingTiles(s.anchorBox(far)); n < 1 || n > 4 {
		t.Errorf("the small drawing is given to %d tiles", n)
	}

	// drawings of some thousand tiles each: those that fit the budget
	// together stay, the next one is left out
	large := func() *anchored {
		return &anchored{kind: "absoluteAnchor", ext: [2]float64{100000, 100000}, toRow: 14000, toCol: 2100}
	}
	one := layoutContext(sheet(), large())
	one.layout()
	each := one.drawingTiles(one.anchorBox(one.drawings[0]))
	if each < 1000 || each > maxDrawingTiles/2 {
		t.Fatalf("a drawing of %d tiles", each)
	}
	var many []*anchored
	for range maxDrawingTiles/each + 1 {
		many = append(many, large())
	}
	s = layoutContext(sheet(), many...)
	s.layout()
	for i, a := range many {
		if want := i >= maxDrawingTiles/each; a.out != want {
			t.Errorf("drawing %d of %d tiles: left out %v", i, each, a.out)
		}
	}
	if n := drawingWarnings(s.c); n != 1 {
		t.Errorf("%d warnings: %q", n, s.c.warnings)
	}
}

// TestDrawingFarAway converts a workbook whose picture is anchored at row
// 900,001 and whose shape is anchored over the whole sheet: the picture is
// drawn in the tile it is in, the shape is left out, and the tiles are few.
func TestDrawingFarAway(t *testing.T) {
	res, r := convertEdited(t, "features.xlsx", map[string]func(string) string{
		"xl/drawings/drawing1.xml": func(s string) string {
			s = strings.Replace(s, "<from><col>8</col><colOff>0</colOff><row>2</row>", "<from><col>8</col><colOff>0</colOff><row>900000</row>", 1)
			return strings.Replace(s, "<from><col>8</col><colOff>0</colOff><row>31</row><rowOff>0</rowOff></from><to><col>11</col><colOff>304800</colOff><row>36</row>",
				"<from><col>0</col><colOff>0</colOff><row>0</row><rowOff>0</rowOff></from><to><col>16383</col><colOff>0</colOff><row>1048575</row>", 1)
		},
	}, testOptions())
	found := 0
	for _, w := range res.Warnings {
		if strings.Contains(w, "drawings cover more than") {
			found++
		}
	}
	if found != 1 {
		t.Errorf("warnings %q", res.Warnings)
	}
	v := r.Manifest.Views[0]
	rows := 0
	for _, run := range v.Rows {
		rows += int(run[0])
	}
	if rows < 900000 || rows > 900100 {
		t.Errorf("view of %d rows", rows)
	}
	if len(v.Tiles) > 8 {
		t.Errorf("%d tiles", len(v.Tiles))
	}
	// the picture, in the tile of its row
	far := ""
	for key := range v.Tiles {
		if !strings.HasSuffix(key, ",0") {
			far = key
		}
	}
	if far == "" {
		t.Fatalf("no tile far away: %v", v.Tiles)
	}
	ms, _ := marks(t, r, tile(t, v, far))
	if !contains(ms, "F Picture") {
		t.Errorf("tile %s: marks %q", far, ms)
	}
}

// TestTileSpan checks the tiles of rectangles: those that tilesIn makes are
// those that tileSpan counts, and a rectangle whose coordinates are not
// numbers has none (the conversion of NaN to an integer depends on the
// machine: the tiles would start anywhere).
func TestTileSpan(t *testing.T) {
	s := layoutContext(&worksheet{name: "S"})
	s.cols, s.rows = newAxis(maxCols, 48, nil), newAxis(1000, 15, nil)
	for _, c := range []struct {
		x0, y0, x1, y1 float64
		want           int
	}{
		{0, 0, 10, 10, 1},
		{0, 0, tileSize, tileSize, 1},
		{0, 0, tileSize + 1, 10, 2},
		{tileSize - 1, tileSize - 1, tileSize + 1, tileSize + 1, 4},
		{-100, -100, 3 * tileSize, 1e9, 3 * 8}, // cut to the 15,000 units of the rows
		{10, 10, 10, 20, 0},
		{1e9, 0, 2e9, 10, 0},
		{math.NaN(), 0, 10, 10, 0},
		{0, 0, 10, math.NaN(), 0},
		{0, math.Inf(-1), 10, math.Inf(1), 8},
	} {
		tx0, ty0, tx1, ty1, ok := s.tileSpan(c.x0, c.y0, c.x1, c.y1)
		n := 0
		if ok {
			n = (tx1 - tx0) * (ty1 - ty0)
		}
		if got := len(s.tilesIn(c.x0, c.y0, c.x1, c.y1)); got != c.want || n != c.want {
			t.Errorf("rectangle %v %v %v %v: %d tiles made, %d counted, want %d", c.x0, c.y0, c.x1, c.y1, got, n, c.want)
		}
	}
}

// TestColumnSteps checks that the work for the columns does not grow with
// the columns of the view or with the col elements over a column: a row of
// a view of all the columns of a sheet that col elements give a format is
// walked in some steps, not in one for each column, and of col elements one
// over the other each column is given to the last, once.
func TestColumnSteps(t *testing.T) {
	ws := &worksheet{name: "S", baseColW: 8, cols: []colDef{{min: 0, max: maxCols - 1, width: -1, style: 1}, {min: 7, max: 8, width: -1, hidden: true, style: 1}},
		rows: []row{{idx: 0, ht: -1, style: -1, cells: []cell{{col: maxCols - 1, kind: cellNum, num: 1}}}}, merges: []cellRange{{20, 2, 22, 3}}}
	s := layoutContext(ws)
	s.c.st.cellXfs = []xf{{}, {fill: 1}}
	s.layout()
	if s.nCols != maxCols {
		t.Fatalf("view of %d columns", s.nCols)
	}
	for r := 0; r < 100; r++ {
		s.eachFormatted(r, func(c0, c1 int, f *cellFmt) {})
	}
	if s.steps > 100*8 {
		t.Errorf("%d steps for 100 rows of %d columns", s.steps, s.nCols)
	}

	// many col elements over the same columns
	many := &worksheet{name: "S"}
	for i := range 3000 {
		many.cols = append(many.cols, colDef{min: 0, max: maxCols - 1, width: float64(i % 7), style: i % 3})
	}
	last, steps := lastOver(many.cols, func(colDef) bool { return true })
	if steps > len(many.cols)+maxCols {
		t.Errorf("%d steps for %d col elements", steps, len(many.cols))
	}
	for c := 0; c < maxCols; c += 97 {
		if int(last[c]) != len(many.cols)-1 || many.colStyle(c) != 2 {
			t.Fatalf("column %d: col element %d, style %d", c, last[c], many.colStyle(c))
		}
	}
}

// TestFormatEdges checks the budget of the edges of the borders that the
// formats of rows and columns give cells that are not in the file: past it
// they are left out, with a warning, and the cells of the file keep theirs.
func TestFormatEdges(t *testing.T) {
	sheet := func(room int) (*sheetCtx, *edgeSet) {
		ws := &worksheet{name: "S", baseColW: 8, cols: []colDef{{min: 0, max: 4, width: -1, style: 1}}, rows: []row{
			{idx: 150, ht: -1, style: -1, cells: []cell{{col: 1, style: 1, kind: cellNum, num: 1}, {col: 2, style: 1}}},
			{idx: 199, ht: -1, style: -1, cells: []cell{{col: 0, kind: cellNum, num: 1}}},
		}}
		s := layoutContext(ws)
		thin := borderSide{style: "thin"}
		s.c.st.borders = []xborder{{}, {left: thin, right: thin, top: thin, bottom: thin}}
		s.c.st.cellXfs = []xf{{}, {border: 1}}
		s.edgeRoom = room
		s.layout()
		return s, s.edges()
	}
	warnings := func(s *sheetCtx) (n int) {
		for _, w := range s.c.warnings {
			if strings.Contains(w, "give borders to more than") {
				n++
			}
		}
		return n
	}
	s, es := sheet(200)
	inFile := 8 // the sides of the two cells of the file, at most
	if n := len(es.h) + len(es.v); n > 200+inFile || n < 200-cellEdges {
		t.Errorf("%d edges kept", n)
	}
	if warnings(s) != 1 {
		t.Errorf("warnings %q", s.c.warnings)
	}
	for _, c := range []int{1, 2} {
		if _, ok := es.h[edgeKey(150, c)]; !ok {
			t.Errorf("no top edge of the cell of the file in column %d", c)
		}
	}
	if _, ok := es.v[edgeKey(150, 0)]; ok {
		t.Error("an edge of a cell that is not in the file, past the budget")
	}
	// all of them within the budget of a sheet
	s, es = sheet(maxFormatEdges)
	if _, ok := es.v[edgeKey(150, 0)]; !ok || warnings(s) != 0 || s.edgeRoom >= maxFormatEdges {
		t.Errorf("within the budget: the edge is kept %v, %d left, warnings %q", ok, s.edgeRoom, s.c.warnings)
	}
}

// TestTableCells checks the budget of the cells of tables that are not in
// the file: a table of cells of the file has its style whatever its size,
// the others while their cells are within the budget.
func TestTableCells(t *testing.T) {
	ws := &worksheet{name: "S"}
	for r := 0; r < 3; r++ {
		ws.rows = append(ws.rows, row{idx: r, ht: -1, style: -1, cells: []cell{{col: 0, kind: cellNum}, {col: 1}, {col: 2, kind: cellNum}}})
	}
	if n := ws.cellsIn(cellRange{1, 1, 5, 5}); n != 4 {
		t.Errorf("%d cells of the file in B2:F6", n)
	}
	s := layoutContext(ws)
	s.tableRoom = 100
	for _, c := range []struct {
		ref  cellRange
		want bool
		left int
	}{
		{cellRange{0, 0, 2, 2}, true, 100}, // its cells are all in the file
		{cellRange{0, 0, 9, 9}, true, 9},   // 91 that are not
		{cellRange{10, 0, 10, 9}, false, 9},
		{cellRange{0, 0, 3, 2}, true, 6},
		{cellRange{0, 0, maxRows - 1, maxCols - 1}, false, 6},
		{cellRange{0, 0, 2, 2}, true, 6},
	} {
		if got := s.admitTable(c.ref); got != c.want || s.tableRoom != c.left {
			t.Errorf("table over %v: styled %v, %d left; want %v, %d", c.ref, got, s.tableRoom, c.want, c.left)
		}
	}
	if len(s.c.warnings) != 1 || !strings.Contains(s.c.warnings[0], "tables cover more than") {
		t.Errorf("warnings %q", s.c.warnings)
	}
}

// TestLinkTiles checks the budget of the tiles of links after their first:
// a link over more tiles than are left is in its first tile only.
func TestLinkTiles(t *testing.T) {
	ws := &worksheet{name: "S", baseColW: 8, rows: []row{{idx: 299, ht: -1, style: -1, cells: []cell{{col: 99, kind: cellNum, num: 1}}}}}
	s := layoutContext(ws)
	s.layout()
	whole := cellRange{0, 0, maxRows - 1, maxCols - 1}
	b := s.rangeBox(cellRange{0, 0, s.nRows - 1, s.nCols - 1})
	tx0, ty0, tx1, ty1, _ := s.tileSpan(b.x, b.y, b.x+b.w, b.y+b.h)
	all := (tx1 - tx0) * (ty1 - ty0)
	if all < 6 {
		t.Fatalf("a view of %d tiles", all)
	}
	s.linkRoom = 4
	s.putLink(cellRange{0, 0, 0, 0}, "https://example.com/")
	s.putLink(whole, "https://example.com/")
	if len(s.tiles) != 1 || s.linkRoom != 4 || len(s.c.warnings) != 1 || !strings.Contains(s.c.warnings[0], "links cover more than") {
		t.Errorf("%d tiles, %d left, warnings %q", len(s.tiles), s.linkRoom, s.c.warnings)
	}
	s.linkRoom = all
	s.putLink(whole, "https://example.com/")
	if len(s.tiles) != all || s.linkRoom != 1 || len(s.c.warnings) != 1 {
		t.Errorf("%d tiles of %d, %d left, warnings %q", len(s.tiles), all, s.linkRoom, s.c.warnings)
	}
}

// TestBudgetsOfSheets converts workbooks that state much about cells they
// do not hold: a table over A16:XFD4096, a link over the whole sheet in a
// view of 300,000 rows by 16,384 columns. The sheets are drawn without
// them, in a few tiles; and the borders of the format of a column, which
// are within their budget, are drawn.
func TestBudgetsOfSheets(t *testing.T) {
	warned := func(res *Result, what string) bool {
		for _, w := range res.Warnings {
			if strings.Contains(w, what) {
				return true
			}
		}
		return false
	}
	rows := func(v *bdf.View) (n int) {
		for _, run := range v.Rows {
			n += int(run[0])
		}
		return n
	}
	res, r := convertEdited(t, "features.xlsx", map[string]func(string) string{
		"xl/tables/table1.xml": func(s string) string { return strings.Replace(s, `ref="A16:C21"`, `ref="A16:XFD4096"`, 1) },
	}, testOptions())
	if v := r.Manifest.Views[0]; !warned(res, "tables cover more than") || rows(v) > 100 || len(v.Tiles) > 4 {
		t.Errorf("table: %d rows, %d tiles, warnings %q", rows(v), len(v.Tiles), res.Warnings)
	}

	res, r = convertEdited(t, "basic.xlsx", map[string]func(string) string{
		"xl/worksheets/sheet1.xml": func(s string) string {
			s = strings.Replace(s, `ref="A26" r:id="rId1"`, `ref="A1:XFD1048576" r:id="rId1"`, 1)
			return strings.Replace(s, "</sheetData>", `<row r="300000"><c r="XFD300000"><v>1</v></c></row></sheetData>`, 1)
		},
	}, testOptions())
	if v := r.Manifest.Views[0]; !warned(res, "links cover more than") || rows(v) < 300000 || len(v.Tiles) > 8 {
		t.Errorf("link: %d rows, %d tiles, warnings %q", rows(v), len(v.Tiles), res.Warnings)
	}

	size := func(style string) int {
		res, _ := convertEdited(t, "basic.xlsx", map[string]func(string) string{
			"xl/worksheets/sheet1.xml": func(s string) string {
				return strings.Replace(s, "</cols>", `<col min="11" max="12" width="9" style="`+style+`"/></cols>`, 1)
			},
		}, testOptions())
		if warned(res, "more than") {
			t.Errorf("warnings %q", res.Warnings)
		}
		var buf bytes.Buffer
		if err := res.Doc.WriteSingle(&buf); err != nil {
			t.Fatal(err)
		}
		return buf.Len()
	}
	// xf 2 of the workbook has a border, xf 1 has none
	if with, without := size("2"), size("1"); with <= without {
		t.Errorf("the borders of the format of a column are not drawn: %d bytes with them, %d without", with, without)
	}
}
