//go:build !tinygo

package htmlro

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/html"
)

// countNodes counts the elements, attributes and text nodes of a tree.
func countNodes(n *html.Node) (elements, attrs, texts int) {
	switch n.Type {
	case html.ElementNode:
		elements++
		attrs += len(n.Attr)
	case html.TextNode:
		texts++
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		e, a, t := countNodes(c)
		elements += e
		attrs += a
		texts += t
	}
	return
}

// A parse allocates what the tree keeps and little else: a node and an
// attribute slice per element, a string per attribute value and per text
// node, which the interning of repeated values and names reduces further.
// Pinned here against a bound, with the tokenizer of golang.org/x/net/html
// alongside for the comparison that is the point of the package.
func TestParseAllocations(t *testing.T) {
	data := benchDoc(2000)
	doc, err := ParseBytes(data, Options{})
	if err != nil {
		t.Fatal(err)
	}
	elements, attrs, texts := countNodes(doc)
	allocs := testing.AllocsPerRun(5, func() {
		ParseBytes(data, Options{})
	})
	xnet := testing.AllocsPerRun(5, func() {
		html.Parse(bytes.NewReader(data))
	})
	// Each element: a node, and an attribute slice when it has attributes
	// (here every element has). Each text node: a node and its string. Each
	// attribute: its value string.
	bound := float64(2*elements + 2*texts + attrs)
	t.Logf("%d elements, %d attributes, %d text nodes: %.0f allocations (%.2f per element), x/net/html %.0f (%.2f per element)",
		elements, attrs, texts, allocs, allocs/float64(elements), xnet, xnet/float64(elements))
	if allocs > bound {
		t.Errorf("%.0f allocations, want at most %.0f", allocs, bound)
	}
	if allocs*2 > xnet {
		t.Errorf("%.0f allocations, x/net/html %.0f: not half", allocs, xnet)
	}
}

// TestCorpusAllocations reports what the two parsers allocate over the HTML
// files under HTMLRO_CORPUS; it is skipped without it.
func TestCorpusAllocations(t *testing.T) {
	dir := os.Getenv("HTMLRO_CORPUS")
	if dir == "" {
		t.Skip("HTMLRO_CORPUS not set")
		return // TinyGo's Skip does not stop the test
	}
	var docs [][]byte
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(path, ".html") {
			if data, err := os.ReadFile(path); err == nil {
				docs = append(docs, data)
			}
		}
		return nil
	})
	var elements, attrs, texts, size int
	for _, data := range docs {
		doc, err := ParseBytes(data, Options{})
		if err != nil {
			t.Fatal(err)
		}
		e, a, x := countNodes(doc)
		elements += e
		attrs += a
		texts += x
		size += len(data)
	}
	measure := func(parse func([]byte)) (allocs float64, bytes uint64, d time.Duration) {
		allocs = testing.AllocsPerRun(1, func() {
			for _, data := range docs {
				parse(data)
			}
		})
		start := time.Now()
		bytes = allocatedBytes(1, func() {
			for _, data := range docs {
				parse(data)
			}
		})
		return allocs, bytes, time.Since(start)
	}
	a1, b1, d1 := measure(func(data []byte) { ParseBytes(data, Options{}) })
	a2, b2, d2 := measure(func(data []byte) { html.Parse(bytes.NewReader(data)) })
	t.Logf("%d documents, %.1f MB, %d elements, %d attributes, %d text nodes", len(docs), float64(size)/1e6, elements, attrs, texts)
	t.Logf("htmlro:     %10.0f allocations (%.2f per element), %6.1f MB, %v", a1, a1/float64(elements), float64(b1)/1e6, d1)
	t.Logf("x/net/html: %10.0f allocations (%.2f per element), %6.1f MB, %v", a2, a2/float64(elements), float64(b2)/1e6, d2)
}
