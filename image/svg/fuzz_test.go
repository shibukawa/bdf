package svg

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// FuzzParse feeds damaged documents to the reader and reads the values of
// their attributes: it may read nothing, but must not panic or hang, and
// what it returns holds together — an svg root, parents and children that
// agree, and no more elements, nor deeper ones, than the limits allow. The
// seeds are the SVG files of the tests and documents with what the reader
// treats specially.
func FuzzParse(f *testing.F) {
	files, _ := filepath.Glob("testdata/*.svg")
	for _, name := range files {
		if b, err := os.ReadFile(name); err == nil {
			f.Add(b)
		}
	}
	for _, s := range []string{
		head + `<style>g .a, #b > rect.c { fill: red !important } /* x */ @media print { a { b: c } }</style><g class="a" id="b" style="stroke: blue; x"><rect class="c"/></g></svg>`,
		`<!DOCTYPE svg [<!ENTITY a "&#x41;"><!ENTITY b "&a;&a;">]><svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"><use xlink:href="#b" href="&b;"/><text>&b;<tspan>t</tspan>&nbsp;</text></svg>`,
		head + `<g transform="translate(1,2) rotate(30 4 5) skewX(6) matrix(1 0 0 1 7 8) scale(-1e3)" fill="hsl(120deg 50% 50% / 25%)" stroke="#abcd" width="1.5em" points="1,2 3-4.5.5"><foreignObject><p>a<br>b</foreignObject>`,
		"\xff\xfe<\x00s\x00v\x00g\x00/\x00>\x00",
		`<?xml version="1.0" encoding="ISO-8859-1"?><svg><text>caf` + "\xe9" + `</text></svg>`,
	} {
		f.Add([]byte(s))
	}
	opts := &Options{MaxDepth: 40, MaxNodes: 2000, MaxStyleMatches: 20000}
	f.Fuzz(func(t *testing.T, b []byte) {
		doc, err := Parse(b, opts)
		if err != nil {
			if doc != nil || !errors.Is(err, ErrNotSVG) {
				t.Fatalf("a document with the error %v", err)
			}
			return
		}
		if doc.Root == nil || doc.Root.Name != "svg" || doc.Root.Parent != nil {
			t.Fatalf("the root is %+v", doc.Root)
		}
		nodes := 0
		var walk func(n *Node, depth int)
		walk = func(n *Node, depth int) {
			// a run of text is in an element, which may be the deepest
			if nodes++; nodes > opts.MaxNodes || depth > opts.MaxDepth+1 || depth > opts.MaxDepth && n.Name != TextNode {
				t.Fatalf("%d nodes, a %s %d deep", nodes, n.Name, depth)
			}
			if n.Name == TextNode {
				if len(n.Children) != 0 || n.Attr != nil || n.Sheet != nil {
					t.Fatalf("a run of text with %d children and the attributes %v", len(n.Children), n.Attr)
				}
				return
			}
			if n.Attr == nil || n.Text != "" {
				t.Fatalf("the element %s has no attributes, or the text %q", n.Name, n.Text)
			}
			for k, v := range n.Attr {
				n.Declared(k)
				ParseColor(v, Color{A: 1})
				ParseLength(v, 100, 16)
				ParseNumbers(v)
				ParseTransform(v)
			}
			for k := range n.Sheet {
				n.Declared(k)
			}
			for _, c := range n.Children {
				if c.Parent != n {
					t.Fatalf("the parent of a %s is not the %s it is in", c.Name, n.Name)
				}
				walk(c, depth+1)
			}
		}
		walk(doc.Root, 1)
		for id, n := range doc.IDs {
			if id == "" || n.Attr["id"] != id {
				t.Fatalf("the id %q is of an element with the id %q", id, n.Attr["id"])
			}
		}
	})
}
