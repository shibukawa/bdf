package svg

import (
	"errors"
	"fmt"
	"math"
	"os"
	"runtime"
	"slices"
	"strings"
	"testing"
	"unicode/utf16"
)

const head = `<svg xmlns="http://www.w3.org/2000/svg">`

func mustParse(t *testing.T, src string) *Document {
	t.Helper()
	doc, err := Parse([]byte(src), nil)
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

func near(a, b, eps float64) bool { return math.Abs(a-b) <= eps }

func TestParse(t *testing.T) {
	doc := mustParse(t, `<?xml version="1.0"?>
<s:svg xmlns:s="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="10" s:height="20">
  <s:defs><s:rect id="a" width="1"/><s:circle id="a" r="2"/></s:defs>
  <s:use xlink:href="#a"/><s:use xlink:href="#old" href="#new"/><s:use href="#new" xlink:href="#old"/>
  <s:g>between<s:text>one<s:tspan>two</s:tspan>three</s:text><s:a>link</s:a><s:title>not kept</s:title></s:g>
</s:svg>`)
	root := doc.Root
	if root.Name != "svg" || root.Parent != nil || root.Attr["width"] != "10" || root.Attr["height"] != "20" || len(root.Attr) != 2 {
		t.Errorf("root %q %v", root.Name, root.Attr)
	}
	if len(root.Children) != 5 {
		t.Fatalf("%d children of the root", len(root.Children))
	}
	for _, c := range root.Children {
		if c.Parent != root {
			t.Errorf("the parent of %s is not the root", c.Name)
		}
	}
	// the first of the elements that share an id
	if a := doc.IDs["a"]; a == nil || a.Name != "rect" || len(doc.IDs) != 1 {
		t.Errorf("ids %v", doc.IDs)
	}
	// xlink:href is href, and an href of the element wins wherever it is
	for i, want := range []string{"#a", "#new", "#new"} {
		if u := root.Children[1+i]; u.Name != "use" || u.Attr["href"] != want || len(u.Attr) != 1 {
			t.Errorf("use %d: %v, want href %s", i, u.Attr, want)
		}
	}
	// runs of text are kept where SVG draws them
	g := root.Children[4]
	var names []string
	for _, c := range g.Children {
		names = append(names, c.Name)
	}
	if fmt.Sprint(names) != "[text a title]" {
		t.Errorf("children of g: %v", names)
	}
	text := g.Children[0]
	if len(text.Children) != 3 || text.Children[0].Name != TextNode || text.Children[0].Text != "one" || text.Children[0].Parent != text ||
		text.Children[1].Name != "tspan" || text.Children[1].Children[0].Text != "two" || text.Children[2].Text != "three" {
		t.Errorf("the runs of text: %+v", text.Children)
	}
	if a := g.Children[1]; len(a.Children) != 1 || a.Children[0].Text != "link" {
		t.Errorf("the text of a: %+v", a.Children)
	}
	if len(g.Children[2].Children) != 0 || len(doc.Warnings) != 0 {
		t.Errorf("title %+v, warnings %v", g.Children[2].Children, doc.Warnings)
	}
	// a run of text declares nothing
	if v, ok := text.Children[0].Declared("fill"); ok || v != "" {
		t.Errorf("a run of text declares fill %q", v)
	}
}

func TestNotSVG(t *testing.T) {
	for _, src := range []string{"", "   ", "plain text", `<html><body><svg/></body></html>`, `<?xml version="1.0"?>`, `<!-- only a comment -->`, "\x00\x01\x02"} {
		doc, err := Parse([]byte(src), nil)
		if doc != nil || !errors.Is(err, ErrNotSVG) {
			t.Errorf("%q: %v, %v", src, doc, err)
		}
	}
	if _, err := Parse(nil, &Options{}); !errors.Is(err, ErrNotSVG) {
		t.Errorf("no document: %v", err)
	}
}

func TestLenient(t *testing.T) {
	// a document that ends in the middle keeps what was read
	doc := mustParse(t, head+`<rect id="a" fill="green"/><g><rect id="b" fill=blue><circle r=5 fill="red"`)
	if doc.IDs["a"] == nil || doc.IDs["b"] == nil || doc.IDs["b"].Attr["fill"] != "blue" || doc.IDs["b"].Parent.Name != "g" {
		t.Errorf("a document cut short: %v", doc.IDs)
	}
	// the entities of a DOCTYPE (Illustrator names namespaces so) and of HTML
	doc = mustParse(t, `<!DOCTYPE svg [<!ENTITY ns_svg "http://www.w3.org/2000/svg"><!ENTITY col "#c22">]>
<svg xmlns="&ns_svg;"><rect fill="&col;"/><text>&lt;A&amp;B&gt;&copy;&nbsp;&#x41;</text></svg>`)
	if f := doc.Root.Children[0].Attr["fill"]; f != "#c22" {
		t.Errorf("an entity in an attribute: %q", f)
	}
	if s := doc.Root.Children[1].Children[0].Text; s != "<A&B>© A" {
		t.Errorf("entities in text: %q", s)
	}
	// encodings: the one the declaration names, and UTF-16 by its mark
	latin1 := append([]byte(`<?xml version="1.0" encoding="ISO-8859-1"?>`+head+`<text>caf`), 0xe9)
	latin1 = append(latin1, `</text></svg>`...)
	if doc, err := Parse(latin1, nil); err != nil || doc.Root.Children[0].Children[0].Text != "café" {
		t.Errorf("Latin-1: %v", err)
	}
	u16 := []byte{0xff, 0xfe}
	for _, u := range utf16.Encode([]rune(head + `<text>日本</text></svg>`)) {
		u16 = append(u16, byte(u), byte(u>>8))
	}
	if doc, err := Parse(u16, nil); err != nil || doc.Root.Children[0].Children[0].Text != "日本" {
		t.Errorf("UTF-16: %v", err)
	}
	// the void elements of HTML inside a foreignObject close themselves
	doc = mustParse(t, head+`<foreignObject><p>a<br>b</p></foreignObject><rect id="after"/></svg>`)
	if n := doc.IDs["after"]; n == nil || n.Parent != doc.Root {
		t.Errorf("the element after HTML: %+v", n)
	}
}

func TestDeclared(t *testing.T) {
	doc := mustParse(t, head+`<style>
		rect { fill: blue } .a { fill: green } #x { fill: yellow } g .b { stroke: red }
		/* left out */ rect[width] { fill: pink } rect:hover { fill: pink } @media print { rect { fill: none } }
		circle { FILL: Red !important; stroke : url(#g) }
	</style><g><rect id="x" class="a b" fill="black"/><rect class="a" style="fill: #fff; stroke-width: 2 !important" stroke-width=" 9 "/></g>
	<circle fill="blue"/><rect width="1" stroke=" navy "/><ellipse/></svg>`)
	x := doc.IDs["x"]
	if v, _ := x.Declared("fill"); v != "yellow" {
		t.Errorf("fill of #x: %q", v)
	}
	if v, _ := x.Declared("stroke"); v != "red" {
		t.Errorf("stroke of #x: %q", v)
	}
	second := doc.Root.Children[1].Children[1]
	if v, _ := second.Declared("fill"); v != "#fff" {
		t.Errorf("inline style: %q", v)
	}
	if v, _ := second.Declared("stroke-width"); v != "2" {
		t.Errorf("inline style with !important: %q", v)
	}
	circle := doc.Root.Children[2]
	if v, _ := circle.Declared("fill"); v != "Red" {
		t.Errorf("a style sheet over a presentation attribute: %q", v)
	}
	if v, _ := circle.Declared("stroke"); v != "url(#g)" {
		t.Errorf("stroke of circle: %q", v)
	}
	// the rules left out are not applied, and an attribute is trimmed
	plain := doc.Root.Children[3]
	if v, _ := plain.Declared("fill"); v != "blue" {
		t.Errorf("fill of the plain rect: %q", v)
	}
	if v, ok := plain.Declared("stroke"); v != "navy" || !ok {
		t.Errorf("a presentation attribute: %q", v)
	}
	if e := doc.Root.Children[4]; e.Sheet != nil {
		t.Errorf("an element no rule matches has declarations: %v", e.Sheet)
	}
	if v, ok := doc.Root.Children[4].Declared("fill"); ok || v != "" {
		t.Errorf("nothing declared: %q", v)
	}
}

func TestColors(t *testing.T) {
	cur := Color{R: 0.25, G: 0.5, B: 0.75, A: 1}
	for s, want := range map[string]Color{
		"#f00": {1, 0, 0, 1}, "#F00a": {1, 0, 0, 170.0 / 255}, "#0080ff": {0, 128.0 / 255, 1, 1}, "#0080ff80": {0, 128.0 / 255, 1, 128.0 / 255},
		"rgb(0, 128, 255)": {0, 128.0 / 255, 1, 1}, "rgba(0,0,0,.5)": {0, 0, 0, 0.5}, "rgb(100% 0% 50% / 25%)": {1, 0, 0.5, 0.25}, "rgb(300, -5, 0)": {1, 0, 0, 1},
		"hsl(120, 100%, 50%)": {0, 1, 0, 1}, "hsla(240deg 100% 50% / .5)": {0, 0, 1, 0.5}, "hsl(-120, 100%, 50%)": {0, 0, 1, 1}, "hsl(0, 0%, 40%)": {0.4, 0.4, 0.4, 1},
		"rebeccapurple": {0x66 / 255.0, 0x33 / 255.0, 0x99 / 255.0, 1}, " ReD ": {1, 0, 0, 1}, "transparent": {}, "currentColor": cur,
	} {
		got, ok := ParseColor(s, cur)
		if !ok || !near(got.R, want.R, 1e-3) || !near(got.G, want.G, 1e-3) || !near(got.B, want.B, 1e-3) || !near(got.A, want.A, 1e-3) {
			t.Errorf("%s: %v, want %v", s, got, want)
		}
	}
	for _, s := range []string{"", "none", "url(#g)", "#ff", "#ggg", "#12345", "rgb(1, 2)", "rgb(", "notacolour"} {
		if c, ok := ParseColor(s, cur); ok {
			t.Errorf("%q is the colour %v", s, c)
		}
	}
}

func TestNumbersAndLengths(t *testing.T) {
	for _, c := range []struct {
		s    string
		v    float64
		rest string
		ok   bool
	}{
		{"12", 12, "", true}, {" ,\t-1.5e2px", -150, "px", true}, {"+.5.5", 0.5, ".5", true}, {"3e", 3, "e", true}, {"1em", 1, "em", true},
		{"", 0, "", false}, {"abc", 0, "abc", false}, {".", 0, ".", false}, {"-", 0, "-", false},
	} {
		if v, rest, ok := ParseNumber(c.s); v != c.v || rest != c.rest || ok != c.ok {
			t.Errorf("%q: %v %q %v, want %v %q %v", c.s, v, rest, ok, c.v, c.rest, c.ok)
		}
	}
	if v := ParseNumbers("0 0,100 50"); !slices.Equal(v, []float64{0, 0, 100, 50}) {
		t.Errorf("a view box: %v", v)
	}
	if v := ParseNumbers("10-5.5.5 1e1,x 7"); !slices.Equal(v, []float64{10, -5.5, 0.5, 10}) {
		t.Errorf("numbers without separators, up to what is none: %v", v)
	}
	if v := ParseNumbers(""); v != nil {
		t.Errorf("no numbers: %v", v)
	}
	for s, want := range map[string]float64{
		"10": 10, "10px": 10, "72pt": 96, "1in": 96, "2.54cm": 96, "25.4mm": 96, "6pc": 96, "101.6q": 96, " 3 PX ": 3,
		"50%": 100, "2em": 24, "2ex": 12, "-1.5em": -18,
	} {
		if v, ok := ParseLength(s, 200, 12); !ok || !near(v, want, 1e-9) {
			t.Errorf("%q: %v %v, want %v", s, v, ok, want)
		}
	}
	for _, s := range []string{"", "auto", "10furlongs", "px", "1 2"} {
		if v, ok := ParseLength(s, 200, 12); ok {
			t.Errorf("%q is the length %v", s, v)
		}
	}
}

func TestTransform(t *testing.T) {
	apply := func(m Matrix, x, y float64) (float64, float64) {
		return m[0]*x + m[2]*y + m[4], m[1]*x + m[3]*y + m[5]
	}
	for _, c := range []struct {
		s          string
		x, y, u, v float64
	}{
		{"", 1, 2, 1, 2}, {"translate(10)", 1, 2, 11, 2}, {"translate(10,20)", 1, 2, 11, 22}, {"scale(2)", 1, 2, 2, 4}, {"scale(2 3)", 1, 2, 2, 6},
		{"rotate(90)", 1, 0, 0, 1}, {"rotate(90, 10, 10)", 10, 0, 20, 10}, {"skewX(45)", 0, 1, 1, 1}, {"skewY(45)", 1, 0, 1, 1},
		{"matrix(1 2 3 4 5 6)", 1, 1, 9, 12},
		// the transform on the right applies first
		{"translate(10,0) scale(2)", 1, 1, 12, 2}, {"scale(2),translate(10,0)", 1, 1, 22, 2},
		// up to what is no transform
		{"translate(5) bogus(1) scale(9)", 1, 1, 6, 1}, {"matrix(1 2 3)", 1, 2, 1, 2}, {"scale(2", 1, 2, 1, 2},
	} {
		if u, v := apply(ParseTransform(c.s), c.x, c.y); !near(u, c.u, 1e-9) || !near(v, c.v, 1e-9) {
			t.Errorf("%q of (%v, %v): (%v, %v), want (%v, %v)", c.s, c.x, c.y, u, v, c.u, c.v)
		}
	}
	if m := Identity.Mul(Matrix{1, 2, 3, 4, 5, 6}); m != (Matrix{1, 2, 3, 4, 5, 6}) {
		t.Errorf("the identity times a matrix: %v", m)
	}
}

// TestLimits checks the limits of reading: elements in elements, their
// number, and the rules of style sheets matched against them.
func TestLimits(t *testing.T) {
	depth := func(n *Node) int {
		d := 0
		for ; len(n.Children) > 0; n = n.Children[0] {
			d++
		}
		return d
	}
	// elements in elements, with a style sheet and text that are gone
	// through by functions that call themselves
	doc := mustParse(t, head+`<style>g{fill:red}</style><text>`+strings.Repeat(`<tspan>`, DefaultMaxDepth+50)+`x</text>`+strings.Repeat(`<g>`, DefaultMaxDepth+50))
	if len(doc.Warnings) != 1 || depth(doc.Root.Children[1]) != DefaultMaxDepth-2 {
		t.Errorf("elements %d deep: read %d deep, %v", DefaultMaxDepth+50, depth(doc.Root.Children[1]), doc.Warnings)
	}
	doc, err := Parse([]byte(head+strings.Repeat(`<g>`, 40)), &Options{MaxDepth: 10})
	if err != nil || len(doc.Warnings) != 1 || depth(doc.Root) != 9 || !strings.Contains(doc.Warnings[0], "more than 10 deep") {
		t.Errorf("elements 40 deep where 10 are the limit: %v, read %d deep, %v", err, depth(doc.Root), doc.Warnings)
	}
	doc, err = Parse([]byte(head+strings.Repeat(`<g/>`, 100)+`<text>a<tspan>b</tspan>c</text></svg>`), &Options{MaxNodes: 50})
	if err != nil || len(doc.Root.Children) != 49 || len(doc.Warnings) != 1 {
		t.Errorf("100 elements where 50 are the limit: %v, %d read, %v", err, len(doc.Root.Children), doc.Warnings)
	}
	doc, err = Parse([]byte(head+`<text>a<tspan>b</tspan>c</text></svg>`), &Options{MaxNodes: 3})
	if err != nil || len(doc.Root.Children[0].Children) != 1 {
		t.Errorf("runs of text past the limit: %v, %+v", err, doc.Root.Children[0].Children)
	}
	// a limit of zero or less is the default one
	doc, err = Parse([]byte(head+strings.Repeat(`<g>`, 40)), &Options{MaxDepth: -1, MaxNodes: 0, MaxStyleMatches: -5})
	if err != nil || len(doc.Warnings) != 0 || depth(doc.Root) != 40 {
		t.Errorf("no limits named: %v, read %d deep, %v", err, depth(doc.Root), doc.Warnings)
	}

	// comments: taken out in one pass
	css := `rect{fill:red}` + strings.Repeat(`/**/`, 20000) + `circle{fill:blue}/*`
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	before := ms.TotalAlloc
	rules := parseCSS(css)
	runtime.ReadMemStats(&ms)
	if n := ms.TotalAlloc - before; len(rules) != 2 || n > 1<<20 {
		t.Errorf("20000 comments: %d rules, %d bytes allocated", len(rules), n)
	}
	for src, want := range map[string]int{"a{}/*x*/b{}": 2, "/*a{}*/b{}": 1, "a{}/*b{}": 1, "/**/a{}/**//**/b{}c/**/{}": 3, "*/a{}": 1} {
		if got := len(parseCSS(src)); got != want {
			t.Errorf("%q: %d rules, want %d", src, got, want)
		}
	}

	// rules matched against elements
	var sheet strings.Builder
	for i := 0; i < 100; i++ {
		fmt.Fprintf(&sheet, ".c{stroke-width:%d}", i)
	}
	src := head + `<style>` + sheet.String() + `</style>` + strings.Repeat(`<g class="c"/>`, 100) + `</svg>`
	doc, err = Parse([]byte(src), &Options{MaxStyleMatches: 5000})
	if err != nil || len(doc.Warnings) != 1 || doc.Root.Children[10].Sheet["stroke-width"] != "99" || doc.Root.Children[90].Sheet != nil {
		t.Errorf("100 rules for 100 elements where 5000 matches are the limit: %v, %v", err, doc.Warnings)
	}
	if doc = mustParse(t, src); len(doc.Warnings) != 0 || doc.Root.Children[100].Sheet["stroke-width"] != "99" {
		t.Errorf("within the limit: %v", doc.Warnings)
	}
}

// TestStyleSheetIndex checks that the rules found by what their selectors
// end with are those found by matching every rule against every element.
func TestStyleSheetIndex(t *testing.T) {
	const css = `* { a: 1 } g * { b: 1 } rect { c: 1 } .k { d: 1 } .k.l { e: 1 } rect.l { f: 1 } #i { g: 1 } g > #i.k { h: 1 }
		svg g rect.l.k#i { i: 1 } .m, circle, #j { j: 1 } . { k: 1 } # { l: 1 } .l { c: 2 } rect { d: 2 } #j { c: 3 } q#i { m: 1 }
		g .k { n: 1 } g g * { o: 1 } *.l { p: 1 } rect.k, rect.k { q: 1 }`
	doc := mustParse(t, head+`<style>`+css+`</style>
		<g class="l"><rect id="i" class="k l k"/><rect class="l m"/><circle id="j" class="k"/><g><g id="i"/></g></g><rect/><text>a<tspan class="k">b</tspan></text></svg>`)
	rules := parseCSS(css)
	elements := 0
	var check func(n *Node)
	check = func(n *Node) {
		if n.Name == TextNode {
			return
		}
		elements++
		// as the style sheets were applied before there was an index
		var hit []*cssRule
		for i := range rules {
			if rules[i].matches(n) {
				hit = append(hit, &rules[i])
			}
		}
		slices.SortStableFunc(hit, func(a, b *cssRule) int {
			if a.spec != b.spec {
				return a.spec - b.spec
			}
			return a.order - b.order
		})
		want := map[string]string{}
		for _, r := range hit {
			for _, d := range r.decls {
				want[d[0]] = d[1]
			}
		}
		if len(want) == 0 && n.Sheet != nil || len(want) > 0 && fmt.Sprint(want) != fmt.Sprint(n.Sheet) {
			t.Errorf("%s id=%q class=%q: %v, want %v", n.Name, n.Attr["id"], n.Attr["class"], n.Sheet, want)
		}
		for _, c := range n.Children {
			check(c)
		}
	}
	check(doc.Root)
	if i := doc.IDs["i"].Sheet; elements != 11 || i["i"] != "1" || i["c"] != "2" || i["q"] != "1" || i["k"] != "" {
		t.Errorf("%d elements, #i has %v", elements, i)
	}
}

// TestFiles reads the SVG files of the tests: images as tools write them.
func TestFiles(t *testing.T) {
	for _, name := range []string{"drawing.svg", "badge.svg"} {
		data, err := os.ReadFile("testdata/" + name)
		if err != nil {
			t.Fatal(err)
		}
		doc, err := Parse(data, nil)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		n := 0
		var walk func(*Node)
		walk = func(e *Node) {
			n++
			for _, c := range e.Children {
				if c.Parent != e {
					t.Errorf("%s: the parent of a %s is not the %s it is in", name, c.Name, e.Name)
				}
				walk(c)
			}
		}
		walk(doc.Root)
		if doc.Root.Name != "svg" || n < 3 || len(doc.Warnings) != 0 || doc.Root.Attr["viewBox"] == "" && doc.Root.Attr["width"] == "" {
			t.Errorf("%s: root %s, %d nodes, %v", name, doc.Root.Name, n, doc.Warnings)
		}
	}
}
