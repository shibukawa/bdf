package webdoc

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func nested(depth int) []byte {
	return []byte(`<?xml version="1.0"?><html xmlns="http://www.w3.org/1999/xhtml"><body>` +
		strings.Repeat("<b>", depth-2) + "x" + strings.Repeat("</b>", depth-2) + `</body></html>`)
}

// Elements may nest as deeply as the HTML parser lets them (512); a
// document nested deeper is refused, as the HTML parser refuses it: the
// readers of the tree would need a stack as deep.
func TestXHTMLDepth(t *testing.T) {
	doc, err := ParseXHTML(nested(512))
	if err != nil {
		t.Fatalf("512 levels: %v", err)
	}
	if got := text(doc); got != "x" {
		t.Errorf("512 levels: text %q", got)
	}
	if _, err := ParseXHTML(nested(513)); err == nil {
		t.Errorf("513 levels are read")
	}
	// Parse falls back to the HTML parser, which refuses it too
	if _, err := Parse(nested(513), ""); err == nil {
		t.Errorf("513 levels are read as HTML")
	}
}

// Character data split by CDATA and processing instructions remains one
// text node until a real child node separates it.
func TestXHTMLTextPieces(t *testing.T) {
	for _, c := range []struct {
		body, want, first string
		children          int
	}{
		{`<p>a<![CDATA[b<c]]>d</p>`, "ab<cd", "ab<cd", 1},
		{`<p>a<?pi x?>b&amp;c<![CDATA[d]]></p>`, "ab&cd", "ab&cd", 1},
		{`<p>a<!-- c -->d&amp;e&nbsp;f</p>`, "ad&e\u00a0f", "a", 3},
	} {
		src := []byte(`<?xml version="1.0"?><html xmlns="http://www.w3.org/1999/xhtml"><body>` + c.body + `</body></html>`)
		doc, err := ParseXHTML(src)
		if err != nil {
			t.Fatalf("%s: %v", c.body, err)
		}
		p := find(doc, "p")
		if p == nil {
			t.Fatalf("%s: no paragraph", c.body)
		}
		if got := text(p); got != c.want {
			t.Errorf("%s: text %q, want %q", c.body, got, c.want)
		}
		children := 0
		for child := p.FirstChild; child != nil; child = child.NextSibling {
			children++
		}
		if p.FirstChild == nil || p.FirstChild.Type != html.TextNode || p.FirstChild.Data != c.first || children != c.children {
			t.Errorf("%s: first child %+v, %d children", c.body, p.FirstChild, children)
		}
	}
}

// A text in many pieces is joined once: joining the pieces one by one
// copied the text for each piece (20000 pieces of 50 bytes: 10 GB).
func TestXHTMLManyPieces(t *testing.T) {
	const pieces = 20000
	piece := strings.Repeat("x", 50)
	src := []byte(`<?xml version="1.0"?><html xmlns="http://www.w3.org/1999/xhtml"><body><p>` +
		strings.Repeat(piece+"<?p?>", pieces) + `</p></body></html>`)
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	doc, err := ParseXHTML(src)
	runtime.ReadMemStats(&after)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(text(doc)); got != pieces*len(piece) {
		t.Errorf("text of %d bytes, want %d", got, pieces*len(piece))
	}
	if n := after.TotalAlloc - before.TotalAlloc; n > 64<<20 {
		t.Errorf("%d pieces allocated %d MiB", pieces, n>>20)
	}
}

// byIDs returns the lookup of the elements of a document by their id.
func byIDs(root *html.Node) func(id string) *html.Node {
	ids := map[string]*html.Node{}
	WalkElements(root, func(e *html.Node) {
		if id := attrStr(e, "id"); id != "" && ids[id] == nil {
			ids[id] = e
		}
	})
	return func(id string) *html.Node { return ids[id] }
}

// The document of an svg element that would be larger than the size given
// is not made: what it borrows from the page and its images can make it
// much larger than the element.
func TestSVGDocumentMax(t *testing.T) {
	path := `<path d="` + strings.Repeat("M0 0L9 9", 200) + `"></path>`
	root, err := Parse([]byte(`<!DOCTYPE html><body><svg width="0" height="0"><symbol id="big">`+path+`</symbol></svg>`+
		`<svg id="a" viewBox="0 0 9 9"><use href="#big"></use><image href="p.png"></image></svg>`), "")
	if err != nil {
		t.Fatal(err)
	}
	svg := byIDs(root)("a")
	embed := func(string) string { return "data:image/png;base64," + strings.Repeat("A", 4000) }
	whole := SVGDocument(svg, "", byIDs(root), embed)
	if !strings.Contains(string(whole), `<defs><symbol id="big">`+path) || len(whole) < 5600 {
		t.Fatalf("document of %d bytes: %.300s", len(whole), whole)
	}
	for _, c := range []struct {
		max int
		ok  bool
	}{{0, true}, {len(whole), true}, {len(whole) - 1, false}, {3000, false}, {100, false}} {
		got, ok := SVGDocumentMax(svg, "", byIDs(root), embed, c.max)
		if ok != c.ok || ok && string(got) != string(whole) || !ok && got != nil {
			t.Errorf("at most %d bytes: a document of %d bytes, %v", c.max, len(got), ok)
		}
	}
	// past the size, no more images are asked for
	asked := 0
	SVGDocumentMax(parseSVG(t, strings.Repeat(`<image href="p.png"></image>`, 50)), "", func(string) *html.Node { return nil }, func(string) string {
		asked++
		return embed("")
	}, 10000)
	if asked != 3 {
		t.Errorf("%d images asked for, want 3", asked)
	}
}

func parseSVG(t *testing.T, content string) *html.Node {
	t.Helper()
	root, err := Parse([]byte(`<!DOCTYPE html><body><svg>`+content+`</svg>`), "")
	if err != nil {
		t.Fatal(err)
	}
	return find(root, "svg")
}

// Elements borrowed from the page that are inside each other are each
// copied, in the order they are referred to; what the inner ones refer to
// is found when the outer one is searched, once.
func TestSVGBorrowedNested(t *testing.T) {
	const levels = 300
	var b strings.Builder
	b.WriteString(`<!DOCTYPE html><body><svg width="0" height="0">`)
	for i := range levels {
		fmt.Fprintf(&b, `<g id="g%d"><use href="#s%d"></use>`, i, i)
	}
	b.WriteString(strings.Repeat(`</g>`, levels))
	for i := range levels {
		fmt.Fprintf(&b, `<symbol id="s%d"></symbol>`, i)
	}
	b.WriteString(`</svg><svg id="a">`)
	for i := levels - 1; i >= 0; i -= 100 {
		fmt.Fprintf(&b, `<use href="#g%d"></use>`, i)
	}
	b.WriteString(`</svg>`)
	root, err := Parse([]byte(b.String()), "")
	if err != nil {
		t.Fatal(err)
	}
	lookups := 0
	ids := byIDs(root)
	out := string(SVGDocument(ids("a"), "", func(id string) *html.Node {
		lookups++
		return ids(id)
	}, func(string) string { return "" }))
	// g299, g199 and g99, then the symbols they use: those of g99 and of
	// what is in it (g100 to g299) first, as g99 is the outermost
	want := `<defs><g id="g299"><use href="#s299"></use></g><g id="g199">`
	if i := strings.Index(out, "<defs>"); i < 0 || !strings.HasPrefix(out[i:], want) {
		t.Errorf("document starts %.200s", out)
	}
	for _, s := range []string{`</g><symbol id="s299"></symbol><symbol id="s199"></symbol><symbol id="s200"></symbol>`, `<symbol id="s298"></symbol><symbol id="s99"></symbol><symbol id="s100"></symbol>`,
		`<symbol id="s198"></symbol></defs>`} {
		if !strings.Contains(out, s) {
			t.Errorf("document lacks %s", s)
		}
	}
	if n := strings.Count(out, "<symbol "); n != 201 || lookups != 3+201 {
		t.Errorf("%d symbols, %d lookups", n, lookups)
	}
}
