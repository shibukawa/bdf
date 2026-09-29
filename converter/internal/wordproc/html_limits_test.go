package wordproc

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"runtime"
	"runtime/debug"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// htmlOptions lays text out with the Word converter's test fonts.
func htmlOptions(views string) *Options {
	return &Options{FontDirs: []string{"../../docx/testdata/fonts"}, NoSystemFonts: true, SystemFonts: true, Views: views, NoTextIndex: true}
}

// layoutHTML converts a parsed document.
func layoutHTML(t *testing.T, d *HTMLDocument, views string) *Result {
	t.Helper()
	res, err := ConvertHTML(d, htmlOptions(views))
	if err != nil {
		t.Fatal(err)
	}
	return res
}

// parseHTML parses a document for layoutHTML.
func parseHTML(t *testing.T, src string) *html.Node {
	t.Helper()
	doc, err := html.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

// viewText returns the text of a view of a converted document, its lines
// and paragraphs separated by spaces.
func viewText(t *testing.T, res *Result, view int) string {
	t.Helper()
	var buf bytes.Buffer
	if err := res.Doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	var words []string
	for _, pg := range r.Manifest.Views[view].Pages {
		for _, l := range pg.Layers {
			o, err := r.Object(l.Obj)
			if err != nil {
				t.Fatal(err)
			}
			runs, err := bdf.ExtractText(o, func(h bdf.Hash) *bdf.ObjectPart {
				ch, _ := r.Object(h)
				return ch
			})
			if err != nil {
				t.Fatal(err)
			}
			for _, run := range runs {
				words = append(words, run.Text)
			}
		}
	}
	return strings.Join(strings.Fields(strings.Join(words, " ")), " ")
}

func hasWarning(res *Result, part string) bool {
	for _, w := range res.Warnings {
		if strings.Contains(w, part) {
			return true
		}
	}
	return false
}

func elem(tag string, a atom.Atom, kids ...*html.Node) *html.Node {
	n := &html.Node{Type: html.ElementNode, Data: tag, DataAtom: a}
	for _, k := range kids {
		n.AppendChild(k)
	}
	return n
}

func textNode(s string) *html.Node { return &html.Node{Type: html.TextNode, Data: s} }

// The text of elements nested deeper than maxDepth is read without a call
// for each level: a tree built deeper than a parser makes it (ConvertNode
// takes any tree) needs no stack as deep as it is.
func TestHTMLDeepElements(t *testing.T) {
	const depth = 30000
	inner := elem("span", atom.Span, textNode("bottom"), elem("script", atom.Script, textNode("code")))
	inner.AppendChild(&html.Node{Type: html.ElementNode, Data: "i", DataAtom: atom.I,
		Attr: []html.Attribute{{Key: "hidden"}}})
	inner.LastChild.AppendChild(textNode("hidden"))
	n := inner
	for range depth {
		n = elem("span", atom.Span, textNode("a"), n, textNode("z"))
	}
	body := elem("body", atom.Body, elem("p", atom.P, n))

	defer debug.SetGCPercent(debug.SetGCPercent(-1)) // stacks shrink when the memory is collected
	done := make(chan uint64)
	var res *Result
	var err error
	go func() {
		var before, after runtime.MemStats
		runtime.ReadMemStats(&before)
		res, err = ConvertHTML(&HTMLDocument{Body: body}, htmlOptions(ViewsScroll))
		runtime.ReadMemStats(&after)
		done <- after.StackInuse - min(before.StackInuse, after.StackInuse)
	}()
	grown := <-done
	if err != nil {
		t.Fatal(err)
	}
	if grown > 8<<20 {
		t.Errorf("%d levels: the stack grew by %d MiB", depth, grown>>20)
	}
	text := strings.ReplaceAll(viewText(t, res, 0), " ", "")
	if want := strings.Repeat("a", depth) + "bottom" + strings.Repeat("z", depth); text != want {
		t.Errorf("text of %d characters, want %d; starts %.20q, middle %q", len(text), len(want), text,
			text[min(len(text), depth-3):min(len(text), depth+9)])
	}
}

// An element renamed after it was parsed (Readability makes a div of a
// heading by its name, which leaves the atom of the heading) is a block.
func TestHTMLRenamedHeading(t *testing.T) {
	h := elem("div", atom.H3, textNode("renamed"))
	body := elem("body", atom.Body, elem("h2", atom.H1, textNode("title")), h, elem("p", atom.P, textNode("text")))
	res := layoutHTML(t, &HTMLDocument{Body: body}, ViewsScroll)
	if got := viewText(t, res, 0); got != "title renamed text" {
		t.Errorf("text %q", got)
	}
}

// pngPixel is a PNG of one pixel.
const pngPixel = "\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00\x00\x1f\x15\xc4\x89" +
	"\x00\x00\x00\rIDATx\xdac\xfc\xcf\xc0P\x0f\x00\x04\x85\x01\x80\x84\xa9\x8c!\x00\x00\x00\x00IEND\xaeB`\x82"

func onePixel(string) ([]byte, error) { return []byte(pngPixel), nil }

// within runs fn and fails the test when it takes longer than the seconds
// given: what the test guards against does not end.
func within(t *testing.T, seconds int, fn func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		fn()
	}()
	select {
	case <-done:
	case <-time.After(time.Duration(seconds) * time.Second):
		t.Fatalf("not done after %d seconds", seconds)
	}
}

// A picture that says a height no number of strips reaches (the cuts of
// the scroll view were made without an end) is as high as maxLength.
func TestHTMLHugeHeight(t *testing.T) {
	for _, height := range []string{"1e300", "Infinity", "NaN", "1e400"} {
		doc := parseHTML(t, `<p>before</p><p><img src="x.png" width="10" height="`+height+`"></p><p>after</p>`)
		var res *Result
		var err error
		within(t, 5, func() {
			res, err = ConvertHTML(&HTMLDocument{Body: doc, Image: onePixel}, htmlOptions(ViewsBoth))
		})
		if err != nil {
			t.Fatalf("height %s: %v", height, err)
		}
		if want := int(math.Ceil(maxLength/stripHeight)) + 2; res.Strips > want {
			t.Errorf("height %s: %d strips, want at most %d", height, res.Strips, want)
		}
		if got := viewText(t, res, 0); got != "before after" {
			t.Errorf("height %s: text %q", height, got)
		}
	}
	r := &htmlReader{size: 12}
	for _, c := range []struct {
		v    string
		want float64
	}{{"10", 7.5}, {"10pt", 10}, {"2em", 24}, {"1e300", maxLength}, {"inf", maxLength}, {"nan", 0}, {"-1", 0}, {"50%", 0}} {
		if got := r.length(c.v); got != c.want {
			t.Errorf("length %q = %v, want %v", c.v, got, c.want)
		}
	}
}

// The scroll view is at most maxStrips strips high; what is below is left
// out with a warning.
func TestHTMLStripsLimit(t *testing.T) {
	defer func(n int) { maxStrips = n }(maxStrips)
	maxStrips = 20
	doc := parseHTML(t, `<p>first</p>`+strings.Repeat(`<p><img src="x.png" width="10" height="8000"></p>`, 5)+`<p>last</p>`)
	res := layoutHTML(t, &HTMLDocument{Body: doc, Image: onePixel}, ViewsScroll)
	if res.Strips < 20 || res.Strips > 24 || !hasWarning(res, "(20 strips)") {
		t.Errorf("%d strips, warnings %q", res.Strips, res.Warnings)
	}
	if got := viewText(t, res, 0); got != "first" {
		t.Errorf("text %q", got)
	}
	// a document that is lower has all its strips and no warning
	maxStrips = 40
	res = layoutHTML(t, &HTMLDocument{Body: doc, Image: onePixel}, ViewsScroll)
	if res.Strips < 25 || res.Strips > 40 || len(res.Warnings) != 0 {
		t.Errorf("%d strips, warnings %q", res.Strips, res.Warnings)
	}
	if got := viewText(t, res, 0); got != "first last" {
		t.Errorf("text %q", got)
	}
}

// Text is assigned by its top edge, while a background crossing a cut is
// drawn in every strip it reaches. A cut moves up to avoid splitting text.
func TestStripCutsAndContents(t *testing.T) {
	h := float64(stripHeight)
	ops := []op{
		{text: true, y0: h - 24, y1: h + 16},
		{text: true, y0: 2*h - 28, y1: 2*h + 22},
		{y0: h - 124, y1: 2*h + 52},
		{text: true, y0: h - 24, y1: h - 24},
	}
	cuts := stripCuts(ops, 3*h)
	if want := []float64{0, h - 24, 2*h - 28, 3 * h}; !slices.Equal(cuts, want) {
		t.Fatalf("cuts %v, want %v", cuts, want)
	}
	in := stripOps(ops, cuts)
	for i, want := range [][]int32{{2}, {0, 2, 3}, {1, 2}} {
		if !slices.Equal(in[i], want) {
			t.Errorf("strip %d has ops %v, want %v", i, in[i], want)
		}
	}
	for _, c := range []struct {
		y    float64
		want int
	}{{-1, 0}, {0, 0}, {h - 25, 0}, {h - 24, 1}, {2*h - 28, 2}, {3 * h, 2}, {math.NaN(), 2}} {
		if got := stripAt(cuts, c.y); got != c.want {
			t.Errorf("stripAt(%v) = %d, want %d", c.y, got, c.want)
		}
	}
}

func longStripOps(lines int, high bool) []op {
	ops := make([]op, lines)
	for i := range ops {
		ops[i] = op{text: true, y0: float64(i) * 20, y1: float64(i)*20 + 20}
	}
	if high {
		ops[0].y1 = float64(lines) * 20
	}
	return ops
}

// The long-document case also covers a line as high as the whole document.
func TestStripsOfLongDocument(t *testing.T) {
	const lines = 100000
	for _, high := range []bool{false, true} {
		ops := longStripOps(lines, high)
		in := stripOps(ops, stripCuts(ops, float64(lines)*20))
		n := 0
		for _, s := range in {
			n += len(s)
		}
		if n != lines {
			t.Errorf("high line %v: %d lines in the strips, want %d", high, n, lines)
		}
	}
}

var benchmarkStrips [][]int32

// Compare sizes outside CI's correctness tests when checking the cost of
// making cuts and assigning lines to strips.
func BenchmarkStripsOfLongDocument(b *testing.B) {
	for _, high := range []bool{false, true} {
		for _, lines := range []int{25000, 100000} {
			b.Run(fmt.Sprintf("high=%v/lines=%d", high, lines), func(b *testing.B) {
				ops := longStripOps(lines, high)
				b.ResetTimer()
				for range b.N {
					benchmarkStrips = stripOps(ops, stripCuts(ops, float64(lines)*20))
				}
			})
		}
	}
}

// viewInstrs calls fn for the instructions of the objects of a view of a
// converted document.
func viewInstrs(t *testing.T, res *Result, view int, fn func(in bdf.Instr)) {
	t.Helper()
	var buf bytes.Buffer
	if err := res.Doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	for _, pg := range r.Manifest.Views[view].Pages {
		for _, l := range pg.Layers {
			o, err := r.Object(l.Obj)
			if err != nil {
				t.Fatal(err)
			}
			o.Walk(fn)
		}
	}
}

// viewMarks returns the marks of a kind in a view of a converted document.
func viewMarks(t *testing.T, res *Result, view int, kind byte) []string {
	t.Helper()
	var marks []string
	viewInstrs(t, res, view, func(in bdf.Instr) {
		if in.Op == bdf.OpMark && byte(in.Args[0].(uint64)) == kind {
			marks = append(marks, in.Args[1].(string))
		}
	})
	return marks
}

// viewImages returns how many images a view of a converted document draws.
func viewImages(t *testing.T, res *Result, view int) int {
	t.Helper()
	n := 0
	viewInstrs(t, res, view, func(in bdf.Instr) {
		if in.Op == bdf.OpImage {
			n++
		}
	})
	return n
}

// The cells that span rows and columns take the slots of their spans until
// the tables of the document have taken maxSpanSlots; the cells after them
// span none. A row of cells that span all the rows below took a slot for
// each cell and row.
func TestHTMLSpanSlots(t *testing.T) {
	defer func(n int) { maxSpanSlots = n }(maxSpanSlots)
	maxSpanSlots = 60
	table := `<table><tr><td colspan="3" rowspan="2">a</td><td rowspan="0">b</td><td>c</td></tr>` + strings.Repeat(`<tr><td>d</td></tr>`, 40) + `</table>`
	// 6 slots, then the 41 rows of a column
	res := layoutHTML(t, &HTMLDocument{Body: parseHTML(t, table)}, ViewsScroll)
	if got := viewMarks(t, res, 0, bdf.MarkCell)[:3]; !slices.Equal(got, []string{"A1:C2", "D1:D41", "E1"}) || len(res.Warnings) != 0 {
		t.Errorf("cells %q, warnings %q", got, res.Warnings)
	}
	// the second table of the document: 6 slots more, and 41 that are left
	res = layoutHTML(t, &HTMLDocument{Body: parseHTML(t, table+table)}, ViewsScroll)
	cells := viewMarks(t, res, 0, bdf.MarkCell)
	second := slices.Index(cells[1:], "A1:C2") + 1
	if got := cells[second : second+3]; second < 3 || !slices.Equal(got, []string{"A1:C2", "D1", "E1"}) || !hasWarning(res, "span more than 60 rows and columns") {
		t.Errorf("cells of the second table %q, warnings %q", got, res.Warnings)
	}
}

// The svg elements of a document make SVG documents of MaxSVG bytes
// together (256 MiB when it does not say); the svg elements after them are
// left out with a warning. Each of them holds the symbol it uses, and the
// picture it shows.
func TestHTMLSVGBytes(t *testing.T) {
	var b strings.Builder
	b.WriteString(`<svg width="0" height="0"><symbol id="big"><path d="` + strings.Repeat("M0 0L9 9", 300) + `"/></symbol></svg>`)
	for i := range 8 {
		fmt.Fprintf(&b, `<p>%d <svg width="10" height="10" data-n="%d"><use href="#big"/></svg></p>`, i, i)
	}
	body := parseHTML(t, b.String())
	res := layoutHTML(t, &HTMLDocument{Body: body, MaxSVG: 10000}, ViewsScroll)
	// each document is some 2600 bytes: three of them
	if n := viewImages(t, res, 0); n != 3 || len(res.Warnings) != 1 || !hasWarning(res, "svg elements after them are left out") {
		t.Errorf("%d svg elements drawn, warnings %q", n, res.Warnings)
	}
	if got := viewText(t, res, 0); got != "0 1 2 3 4 5 6 7" {
		t.Errorf("text %q", got)
	}
	// pictures count as they are written into the documents
	read := 0
	picture := func(string) ([]byte, error) {
		read++
		return append([]byte(pngPixel), make([]byte, 3000)...), nil
	}
	body = parseHTML(t, strings.Repeat(`<svg width="10" height="10"><rect width="1" height="1"/><image href="p.png" width="9" height="9"/></svg>`, 8))
	res = layoutHTML(t, &HTMLDocument{Body: body, Image: picture, MaxSVG: 10000}, ViewsScroll)
	if n := viewImages(t, res, 0); n != 2 || read != 3 || len(res.Warnings) != 1 {
		t.Errorf("%d svg elements drawn, the picture read %d times, warnings %q", n, read, res.Warnings)
	}
	res = layoutHTML(t, &HTMLDocument{Body: body, Image: picture}, ViewsScroll)
	if n := viewImages(t, res, 0); n != 8 || len(res.Warnings) != 0 {
		t.Errorf("%d svg elements drawn, warnings %q", n, res.Warnings)
	}
}

// Pictures left out past a limit of the document are warned of once, with
// the warning of the limit; the others each, by their address without the
// user and password of a URL. They show their alternative text.
func TestHTMLImageWarnings(t *testing.T) {
	image := func(src string) ([]byte, error) {
		switch {
		case strings.Contains(src, "limit"):
			return nil, &LimitError{Warning: "more pictures than may be fetched; the pictures after them are left out"}
		case strings.Contains(src, "quiet"):
			return nil, fmt.Errorf("left out: %w", &LimitError{})
		}
		return nil, errors.New("not found")
	}
	body := parseHTML(t, `<p><img src="https://user:secret@example.com/a.png" alt="a"> <img src="limit1.png" alt="b"> <img src="limit2.png" alt="c">
<img src="quiet.png" alt="d"> <svg width="9" height="9"><image href="http://me:pw@example.com/limit3.png"/><image href="http://me:pw@example.com/e.png"/></svg></p>`)
	res := layoutHTML(t, &HTMLDocument{Body: body, Image: image}, ViewsScroll)
	want := []string{"image https://example.com/a.png: not found", "more pictures than may be fetched; the pictures after them are left out",
		"image http://example.com/e.png in an SVG: not found"}
	if !slices.Equal(res.Warnings, want) {
		t.Errorf("warnings %q, want %q", res.Warnings, want)
	}
	if got := viewText(t, res, 0); got != "a b c d" {
		t.Errorf("text %q", got)
	}
	for src, want := range map[string]string{"a.png": "a.png", "http://h/a@b.png": "http://h/a@b.png", "https://u@h/p?q=@": "https://h/p?q=@",
		"//u:p@h/x": "//u:p@h/x", "data:image/png;base64," + strings.Repeat("A", 100): "data:image/png;base64," + strings.Repeat("A", 55) + "…"} {
		if got := imageName(src); got != want {
			t.Errorf("%s is named %s, want %s", src, got, want)
		}
	}
}

// Spaces at wrap points stay on the preceding line; an explicit break ends
// its line and the last line is marked as such.
func TestLineBreaks(t *testing.T) {
	sty := &runStyle{size: 10}
	ch := func(r rune, w float64, brk bool) item {
		return item{kind: kChar, r: r, st: sty, size: 10, w: w, brk: brk}
	}
	p := &para{pp: defaultPProps(), mark: sty, items: []item{
		ch('a', 5, false), ch(' ', 2, true),
		ch('b', 5, false), ch(' ', 2, true),
		ch('c', 5, false), {kind: kBreak, r: '\n', st: sty},
		ch('d', 5, false),
	}}
	lc := &lineCtx{}
	for _, c := range []struct {
		start     int
		text      string
		last      bool
		endsBreak bool
	}{
		{0, "a ", false, false},
		{2, "b ", false, false},
		{4, "c\n", false, true},
		{6, "d", true, false},
	} {
		ln := lc.breakLine(p, c.start, c.start == 0, 0, 11)
		var text strings.Builder
		for _, it := range ln.items {
			text.WriteRune(it.r)
		}
		if text.String() != c.text || ln.last != c.last || ln.endsBreak != c.endsBreak || ln.used != 5 {
			t.Errorf("line at %d: text %q, last %v, break %v, used %v", c.start, text.String(), ln.last, ln.endsBreak, ln.used)
		}
	}
}

// A line longer than the first batch of lineItems is still laid out fully.
func TestLineBeyondInitialBatch(t *testing.T) {
	sty := &runStyle{size: 10}
	p := &para{pp: defaultPProps(), mark: sty, items: make([]item, 300)}
	for i := range p.items {
		p.items[i] = item{kind: kChar, r: 'a', st: sty, size: 10, w: 5}
	}
	lc := &lineCtx{}
	first := lc.breakLine(p, 0, true, 0, 1000)
	second := lc.breakLine(p, len(first.items), false, 0, 1000)
	if len(first.items) != 200 || first.last || len(second.items) != 100 || !second.last {
		t.Errorf("line lengths %d and %d, last %v and %v", len(first.items), len(second.items), first.last, second.last)
	}
}

func longParagraph(size int, brk bool) *para {
	sty := &runStyle{size: 10}
	p := &para{pp: defaultPProps(), mark: sty, items: make([]item, size)}
	for i := range p.items {
		it := item{kind: kChar, r: 'a', st: sty, size: 10, w: 5}
		switch {
		case i%6 == 5:
			it.r, it.brk = ' ', true
		case brk && i%40 == 39:
			it = item{kind: kBreak, r: '\n', st: sty}
		}
		p.items[i] = it
	}
	return p
}

func paragraphLineCount(p *para) (lines, consumed int) {
	lc := &lineCtx{}
	for consumed < len(p.items) {
		n := len(lc.breakLine(p, consumed, consumed == 0, 0, 432).items)
		if n == 0 {
			break
		}
		consumed += n
		lines++
	}
	return lines, consumed
}

// Long paragraphs with filled lines and explicit breaks both make progress.
func TestLinesOfLongParagraph(t *testing.T) {
	const size = 80000
	for _, brk := range []bool{false, true} {
		p := longParagraph(size, brk)
		n, consumed := paragraphLineCount(p)
		if consumed != size || n < size/100 || n > size/30 {
			t.Errorf("breaks %v: %d of %d items consumed in %d lines", brk, consumed, size, n)
		}
	}
}

var benchmarkParagraphLines int

// Benchmark the line-breaking work without timing input construction or
// depending on the load of a CI runner.
func BenchmarkLinesOfLongParagraph(b *testing.B) {
	for _, brk := range []bool{false, true} {
		for _, size := range []int{20000, 80000} {
			b.Run(fmt.Sprintf("breaks=%v/items=%d", brk, size), func(b *testing.B) {
				p := longParagraph(size, brk)
				if _, consumed := paragraphLineCount(p); consumed != size {
					b.Fatalf("consumed %d of %d items", consumed, size)
				}
				b.ResetTimer()
				for range b.N {
					benchmarkParagraphLines, _ = paragraphLineCount(p)
				}
			})
		}
	}
}
