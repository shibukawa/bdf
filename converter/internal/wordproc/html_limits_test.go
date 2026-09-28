package wordproc

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/rand"
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

// stripCutsScanning, safeCutOps and stripOpsScanning are how the scroll
// view was cut into strips: every op was looked at for every cut and for
// every strip.
func stripCutsScanning(ops []op, total float64) []float64 {
	cuts := []float64{0}
	y := 0.0
	for y+stripHeight*1.25 < total {
		c := safeCutOps(ops, y+stripHeight, y+stripHeight/2)
		if c <= y+1 {
			c = y + stripHeight
		}
		cuts = append(cuts, c)
		y = c
	}
	return append(cuts, total)
}

func safeCutOps(ops []op, y, min float64) float64 {
	for range 1000 {
		moved := false
		for _, o := range ops {
			if o.text && o.y0 < y-0.01 && o.y1 > y+0.01 {
				y = o.y0
				moved = true
			}
		}
		if !moved || y <= min {
			break
		}
	}
	return y
}

func stripOpsScanning(ops []op, cuts []float64) [][]int32 {
	in := make([][]int32, len(cuts)-1)
	for i := 0; i+1 < len(cuts); i++ {
		a, b := cuts[i], cuts[i+1]
		for k, o := range ops {
			drawn := o.y1 > a && o.y0 < b
			if o.text {
				drawn = o.y0 >= a-0.01 && o.y0 < b-0.01 || o.y0 == o.y1 && o.y0 >= a && o.y0 < b
			}
			if !drawn {
				continue
			}
			in[i] = append(in[i], int32(k))
		}
	}
	return in
}

// The strips are cut where they were, and hold the ops they held, when
// every op was looked at for each of them.
func TestStripsAsScanned(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for round := range 300 {
		total := rnd.Float64() * []float64{30000, 30000, 200000}[rnd.Intn(3)]
		var ops []op
		y := 0.0
		for range rnd.Intn(400) {
			o := op{text: rnd.Intn(3) > 0}
			switch rnd.Intn(12) {
			case 0: // a high line (a picture in the text), over several strips or many
				o.y0, o.y1 = y, y+rnd.Float64()*[]float64{6000, 6000, 60000}[rnd.Intn(3)]
			case 1: // a background behind much of the document
				o.y0 = rnd.Float64() * total
				o.y1 = o.y0 + rnd.Float64()*total
			case 2: // no height: the structure of an empty cell
				o.y0, o.y1 = y, y
			case 3: // at a cut, or next to one
				o.y0 = float64(rnd.Intn(30))*stripHeight + float64(rnd.Intn(5)-2)*0.01
				o.y1 = o.y0 + rnd.Float64()*30
			case 4: // outside the view, or without a position
				o.y0 = []float64{-50, -0.005, total, total + 100, math.Inf(1), math.Inf(-1), math.NaN()}[rnd.Intn(7)]
				o.y1 = []float64{o.y0 + 20, math.Inf(1), math.NaN(), -10}[rnd.Intn(4)]
			case 5: // a cell beside the last lines: back up
				y = math.Max(0, y-rnd.Float64()*200)
				o.y0, o.y1 = y, y+rnd.Float64()*40
			default:
				o.y0, o.y1 = y, y+10+rnd.Float64()*30
			}
			ops = append(ops, o)
			if o.y1 > y && o.y1 < total && rnd.Intn(4) > 0 {
				y = o.y1
			}
		}
		want := stripCutsScanning(ops, total)
		cuts := stripCuts(ops, total)
		if !slices.Equal(cuts, want) {
			t.Fatalf("round %d: cuts %v, want %v", round, cuts, want)
		}
		got, wantIn := stripOps(ops, cuts), stripOpsScanning(ops, cuts)
		for i := range wantIn {
			if !slices.Equal(got[i], wantIn[i]) {
				t.Fatalf("round %d: strip %d (%v to %v) holds the ops %v, want %v", round, i, cuts[i], cuts[i+1], got[i], wantIn[i])
			}
		}
		// the strip of a link's target
		for range 20 {
			pos := []float64{rnd.Float64() * total, -1, 0, total, total + 1, math.NaN(), cuts[rnd.Intn(len(cuts))]}[rnd.Intn(7)]
			want := len(cuts) - 1
			for i := 1; i < len(cuts); i++ {
				if pos < cuts[i] {
					want = i
					break
				}
			}
			if got := stripAt(cuts, pos) + 1; got != want && !(pos < 0) {
				t.Fatalf("round %d: position %v is in strip %d, want %d", round, pos, got, want)
			}
		}
	}
}

// fastest returns the shortest of three runs of fn.
func fastest(fn func()) time.Duration {
	best := time.Duration(math.MaxInt64)
	for range 3 {
		start := time.Now()
		fn()
		best = min(best, time.Since(start))
	}
	return best
}

// growsWithSize fails the test when fn takes much more than four times as
// long for a size four times as large: the time grows with the square of
// the size (sixteen times), not with the size.
func growsWithSize(t *testing.T, size int, fn func(size int)) {
	t.Helper()
	small := fastest(func() { fn(size) })
	large := fastest(func() { fn(4 * size) })
	if large > 8*small+20*time.Millisecond {
		t.Errorf("size %d took %v, size %d took %v", size, small, 4*size, large)
	}
}

// A long document is cut into strips in a time that grows with its length:
// each cut looks at the lines around it, and each strip at its own, also
// beside a line as high as the document (a picture in a cell of a table).
func TestStripsOfLongDocument(t *testing.T) {
	for _, high := range []bool{false, true} {
		growsWithSize(t, 25000, func(lines int) {
			ops := make([]op, lines)
			for i := range ops {
				ops[i] = op{text: true, y0: float64(i) * 20, y1: float64(i)*20 + 20}
			}
			if high {
				ops[0].y1 = float64(lines) * 20
			}
			in := stripOps(ops, stripCuts(ops, float64(lines)*20))
			n := 0
			for _, s := range in {
				n += len(s)
			}
			if n != lines {
				t.Errorf("%d lines in the strips, want %d", n, lines)
			}
		})
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

// breakLineWhole is breakLine as it was when it placed the rest of the
// paragraph for every line.
func (lc *lineCtx) breakLineWhole(p *para, start int, first bool, left, right float64) *line {
	indL, _, indFirst := p.indents()
	ln := &line{p: p, first: first, left: left, right: right}
	x0 := left
	if first {
		x0 += indFirst
		if p.label != nil {
			ln.label = append([]item(nil), p.label.items...)
		}
	}
	items := lc.fitObjects(p.items[start:], right-x0)
	end := len(items)
	pos := lc.place(p, ln, items, x0, indL, first, right+2*math.Max(right-x0, 100))
	if len(pos) < len(items) {
		end = len(pos)
	}
	lastBrk := -1
	forced := false
	for i := range pos {
		it := &pos[i]
		if it.kind == kBreak || it.kind == kPage || it.kind == kColumn {
			end = i + 1
			forced = true
			break
		}
		overflow := it.x+it.w+it.pad > right+0.01 && i > 0 && !(it.kind == kChar && lbSpace(it.r))
		if overflow && p.pp.overflowPunct && it.kind == kChar && hangingPunct(it.r) && it.x <= right+0.01 {
			overflow = false
		}
		if overflow {
			if lastBrk >= 0 {
				end = lastBrk + 1
			} else {
				end = i
			}
			break
		}
		if it.brk {
			lastBrk = i
		}
	}
	if end == 0 && len(items) > 0 {
		end = 1
	}
	for !forced && end < len(items) && items[end].kind == kChar && lbSpace(items[end].r) {
		end++
	}
	ln.items = lc.place(p, ln, items[:end], x0, indL, first, math.Inf(1))
	if n := len(ln.items); n > 0 {
		switch last := ln.items[n-1]; last.kind {
		case kBreak:
			ln.endsBreak = true
		case kPage, kColumn:
			ln.after = last.kind
		}
	}
	ln.last = start+end >= len(p.items)
	ln.used = x0
	for i := len(ln.items) - 1; i >= 0; i-- {
		it := ln.items[i]
		if (it.kind == kChar && lbSpace(it.r)) || it.kind == kBreak || it.kind == kPage || it.kind == kColumn {
			continue
		}
		ln.used = it.x + it.w + it.pad
		break
	}
	if len(ln.label) > 0 {
		ln.used = math.Max(ln.used, ln.label[len(ln.label)-1].x+ln.label[len(ln.label)-1].w)
	}
	ln.align(lc)
	ln.metrics(lc)
	return ln
}

// lineString writes what a line is, to compare lines.
func lineString(ln *line) string {
	var b strings.Builder
	fmt.Fprintf(&b, "first %v last %v break %v after %d sep %d asc %v desc %v height %v baseline %v left %v right %v used %v\n",
		ln.first, ln.last, ln.endsBreak, ln.after, ln.sep, ln.asc, ln.desc, ln.height, ln.baseline, ln.left, ln.right, ln.used)
	for _, list := range [][]item{ln.label, ln.items} {
		for _, it := range list {
			o := inlineObj{}
			if it.obj != nil {
				o = *it.obj
				o.paint = nil
			}
			it.obj = nil
			fmt.Fprintf(&b, "%+v %+v\n", it, o)
		}
		b.WriteString("\n")
	}
	return b.String()
}

// randomPara makes a paragraph of words, breaks, tabs, pictures and
// characters of vertical text, of a size and a kind that differ.
func randomPara(rnd *rand.Rand) *para {
	pp := defaultPProps()
	pp.lh = []float64{0, 1.7}[rnd.Intn(2)]
	pp.jc = []string{"left", "right", "center", "both", "distribute"}[rnd.Intn(5)]
	pp.overflowPunct = rnd.Intn(2) == 0
	pp.snapToGrid = rnd.Intn(2) == 0
	pp.indL, pp.indFirst = float64(rnd.Intn(3))*10, float64(rnd.Intn(3)-1)*12
	if rnd.Intn(4) == 0 {
		pp.tabs = []tabStop{{pos: 100, align: "right", leader: "dot"}, {pos: 180, align: "decimal", leader: "none"}, {pos: 300, align: "left", leader: "none"}}
	}
	sty := &runStyle{size: 10}
	p := &para{pp: pp, mark: sty}
	if rnd.Intn(4) == 0 {
		p.label = &label{jc: []string{"left", "right", "center"}[rnd.Intn(3)], suffix: []string{"tab", "space", "nothing"}[rnd.Intn(3)], w: 12,
			items: []item{{kind: kChar, r: '1', w: 6, size: 10, st: sty}, {kind: kChar, r: '.', w: 6, size: 10, st: sty}}}
	}
	narrow := rnd.Intn(3) == 0 // many characters to a line
	special := rnd.Intn(3)     // how often something other than a word comes
	tabs := rnd.Intn(3) == 0
	tcy := 0
	for n := rnd.Intn(1000); len(p.items) < n; {
		switch k := rnd.Intn(40); {
		case k == 0 && special > 0:
			p.items = append(p.items, item{kind: []uint8{kBreak, kBreak, kPage, kColumn}[rnd.Intn(4)], r: '\n', st: sty, clear: rnd.Intn(2) == 0})
		case k == 1 && special > 0 && tabs:
			p.items = append(p.items, item{kind: kTab, r: '\t', st: sty, size: 10})
		case k == 2 && special > 0:
			o := &inlineObj{w: 5 + rnd.Float64()*[]float64{30, 300, 900}[rnd.Intn(3)], h: 5 + rnd.Float64()*200, fit: rnd.Intn(2) == 0, fill: rnd.Intn(4) == 0,
				desc: float64(rnd.Intn(2)) * 3, central: rnd.Intn(2) == 0}
			p.items = append(p.items, item{kind: kObject, obj: o, st: sty, w: o.w, brk: rnd.Intn(2) == 0})
		case k == 3 && special > 1:
			tcy++
			for range 1 + rnd.Intn(3) {
				p.items = append(p.items, item{kind: kChar, r: '0' + rune(rnd.Intn(10)), st: sty, size: 10, w: 5, tcy: tcy})
			}
			p.items[len(p.items)-1].brk = true
		case k < 8: // spaces
			for range 1 + rnd.Intn(3)*rnd.Intn(2)*rnd.Intn(40) {
				p.items = append(p.items, item{kind: kChar, r: ' ', st: sty, size: 10, w: 2.5})
			}
			p.items[len(p.items)-1].brk = true
		case k < 12: // East Asian text: a break after each character
			for range 1 + rnd.Intn(30) {
				r := []rune("あ漢字、。")[rnd.Intn(5)]
				p.items = append(p.items, item{kind: kChar, r: r, st: sty, size: 10, w: 10, gap: float64(rnd.Intn(2)) * 2.5, brk: r != '、' && rnd.Intn(8) > 0})
			}
		default: // a word
			w := 5.0
			if narrow {
				w = 0.5
			}
			for range 1 + rnd.Intn(12)*(1+rnd.Intn(2)*rnd.Intn(30)) {
				p.items = append(p.items, item{kind: kChar, r: 'a' + rune(rnd.Intn(26)), st: sty, size: 10, w: w + float64(rnd.Intn(3))})
			}
			if rnd.Intn(6) == 0 {
				p.items = append(p.items, item{kind: kChar, r: []rune(",.-")[rnd.Intn(3)], st: sty, size: 10, w: 3, brk: rnd.Intn(2) == 0})
			}
		}
	}
	return p
}

// The lines of a paragraph break where they broke, and hold what they
// held, when the rest of the paragraph was placed for every line.
func TestLinesAsPlacedWhole(t *testing.T) {
	rnd := rand.New(rand.NewSource(7))
	lines, long := 0, 0
	for round := range 150 {
		p := randomPara(rnd)
		lc := &lineCtx{defTab: 36, vertical: rnd.Intn(4) == 0, expand: rnd.Intn(2) == 0, maxObj: float64(rnd.Intn(2)) * 300}
		if rnd.Intn(3) == 0 {
			lc.charPitch, lc.allChars = 10.5, rnd.Intn(2) == 0
		}
		if rnd.Intn(3) == 0 {
			lc.grid = 18
		}
		left := float64(rnd.Intn(3)) * 20
		right := left + []float64{0.5, 40, 200, 432, 451.3}[rnd.Intn(5)]
		for start, first := 0, true; first || start < len(p.items); first = false {
			got := lc.breakLine(p, start, first, left, right)
			want := lc.breakLineWhole(p, start, first, left, right)
			if g, w := lineString(got), lineString(want); g != w {
				t.Fatalf("round %d, the line at item %d of %d:\n%s\nwant\n%s", round, start, len(p.items), g, w)
			}
			if len(got.items) == 0 {
				break
			}
			if lines++; len(got.items) > lineItems {
				long++
			}
			start += len(got.items)
		}
	}
	if lines < 1500 || long < 40 {
		t.Errorf("%d lines compared, %d of them of more than %d items", lines, long, lineItems)
	}
}

// The lines of a long paragraph are broken in a time that grows with its
// length: a line is looked for in the items that follow its start, not in
// the rest of the paragraph.
func TestLinesOfLongParagraph(t *testing.T) {
	for _, brk := range []bool{false, true} { // words on lines that fill, and lines that end with a break (pre)
		growsWithSize(t, 20000, func(size int) {
			sty := &runStyle{size: 10}
			p := &para{pp: defaultPProps(), mark: sty}
			for i := range size {
				it := item{kind: kChar, r: 'a', st: sty, size: 10, w: 5}
				switch {
				case i%6 == 5:
					it.r, it.brk = ' ', true
				case brk && i%40 == 39:
					it = item{kind: kBreak, r: '\n', st: sty}
				}
				p.items = append(p.items, it)
			}
			lc := &lineCtx{}
			n := 0
			for start := 0; start < len(p.items); n++ {
				start += len(lc.breakLine(p, start, start == 0, 0, 432).items)
			}
			if n < size/100 || n > size/30 {
				t.Errorf("%d items in %d lines", size, n)
			}
		})
	}
}
