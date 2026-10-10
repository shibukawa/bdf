// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package htmlro

// The tree construction tests of golang.org/x/net/html, run on this
// parser: the html5lib test suite and that package's own cases, each a
// document and the dump of the tree it must give.

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type testAttrs struct {
	text, want, context string
	scripting           bool
}

func readParseTest(r *bufio.Reader) (*testAttrs, error) {
	ta := &testAttrs{scripting: true}
	line, err := r.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	var b []byte

	// Read the HTML.
	if string(line) != "#data\n" {
		return nil, fmt.Errorf(`got %q want "#data\n"`, line)
	}
	for {
		line, err = r.ReadSlice('\n')
		if err != nil {
			return nil, err
		}
		if line[0] == '#' {
			break
		}
		b = append(b, line...)
	}
	ta.text = strings.TrimSuffix(string(b), "\n")
	b = b[:0]

	// Skip the error list.
	if string(line) != "#errors\n" {
		return nil, fmt.Errorf(`got %q want "#errors\n"`, line)
	}
	for {
		line, err = r.ReadSlice('\n')
		if err != nil {
			return nil, err
		}
		if line[0] == '#' {
			break
		}
	}

	// Skip the new-errors list.
	if string(line) == "#new-errors\n" {
		for {
			line, err = r.ReadSlice('\n')
			if err != nil {
				return nil, err
			}
			if line[0] == '#' {
				break
			}
		}
	}

	if ls := string(line); strings.HasPrefix(ls, "#script-") {
		switch {
		case strings.HasSuffix(ls, "-on\n"):
			ta.scripting = true
		case strings.HasSuffix(ls, "-off\n"):
			ta.scripting = false
		default:
			return nil, fmt.Errorf(`got %q, want "#script-on" or "#script-off"`, line)
		}
		for {
			line, err = r.ReadSlice('\n')
			if err != nil {
				return nil, err
			}
			if line[0] == '#' {
				break
			}
		}
	}

	if string(line) == "#document-fragment\n" {
		line, err = r.ReadSlice('\n')
		if err != nil {
			return nil, err
		}
		ta.context = strings.TrimSpace(string(line))
		line, err = r.ReadSlice('\n')
		if err != nil {
			return nil, err
		}
	}

	// Read the dump of what the parse tree should be.
	if string(line) != "#document\n" {
		return nil, fmt.Errorf(`got %q want "#document\n"`, line)
	}
	inQuote := false
	for {
		line, err = r.ReadSlice('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		trimmed := bytes.Trim(line, "| \n")
		if len(trimmed) > 0 {
			if line[0] == '|' && trimmed[0] == '"' {
				inQuote = true
			}
			if trimmed[len(trimmed)-1] == '"' && !(line[0] == '|' && len(trimmed) == 1) {
				inQuote = false
			}
		}
		if len(line) == 0 || len(line) == 1 && line[0] == '\n' && !inQuote {
			break
		}
		b = append(b, line...)
	}
	ta.want = string(b)
	return ta, nil
}

func dumpIndent(w io.Writer, level int) {
	io.WriteString(w, "| ")
	for i := 0; i < level; i++ {
		io.WriteString(w, "  ")
	}
}

type sortedAttributes []html.Attribute

func (a sortedAttributes) Len() int {
	return len(a)
}

func (a sortedAttributes) Less(i, j int) bool {
	return a[i].Key < a[j].Key
}

func (a sortedAttributes) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func dumpLevel(w io.Writer, n *html.Node, level int) error {
	dumpIndent(w, level)
	level++
	switch n.Type {
	case html.ErrorNode:
		return errors.New("unexpected ErrorNode")
	case html.DocumentNode:
		return errors.New("unexpected DocumentNode")
	case html.ElementNode:
		if n.Namespace != "" {
			fmt.Fprintf(w, "<%s %s>", n.Namespace, n.Data)
		} else {
			fmt.Fprintf(w, "<%s>", n.Data)
		}
		attr := sortedAttributes(n.Attr)
		sort.Sort(attr)
		for _, a := range attr {
			io.WriteString(w, "\n")
			dumpIndent(w, level)
			if a.Namespace != "" {
				fmt.Fprintf(w, `%s %s="%s"`, a.Namespace, a.Key, a.Val)
			} else {
				fmt.Fprintf(w, `%s="%s"`, a.Key, a.Val)
			}
		}
		if n.Namespace == "" && n.DataAtom == atom.Template {
			io.WriteString(w, "\n")
			dumpIndent(w, level)
			level++
			io.WriteString(w, "content")
		}
	case html.TextNode:
		fmt.Fprintf(w, `"%s"`, n.Data)
	case html.CommentNode:
		fmt.Fprintf(w, "<!-- %s -->", n.Data)
	case html.DoctypeNode:
		fmt.Fprintf(w, "<!DOCTYPE %s", n.Data)
		if n.Attr != nil {
			var p, s string
			for _, a := range n.Attr {
				switch a.Key {
				case "public":
					p = a.Val
				case "system":
					s = a.Val
				}
			}
			if p != "" || s != "" {
				fmt.Fprintf(w, ` "%s"`, p)
				fmt.Fprintf(w, ` "%s"`, s)
			}
		}
		io.WriteString(w, ">")
	case scopeMarkerNode:
		return errors.New("unexpected scopeMarkerNode")
	default:
		return errors.New("unknown node type")
	}
	io.WriteString(w, "\n")
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := dumpLevel(w, c, level); err != nil {
			return err
		}
	}
	return nil
}

// dump writes a tree in the format of the html5lib tests, from the
// children of n.
func dump(n *html.Node) (string, error) {
	if n == nil || n.FirstChild == nil {
		return "", nil
	}
	var b bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := dumpLevel(&b, c, 0); err != nil {
			return "", err
		}
	}
	return b.String(), nil
}

// checkTreeConsistency checks that the links of a tree agree with one
// another.
func checkTreeConsistency(n *html.Node) error {
	var prev *html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Parent != n {
			return fmt.Errorf("%v: child %v has parent %v", n, c, c.Parent)
		}
		if c.PrevSibling != prev {
			return fmt.Errorf("%v: child %v has PrevSibling %v, want %v", n, c, c.PrevSibling, prev)
		}
		if err := checkTreeConsistency(c); err != nil {
			return err
		}
		prev = c
	}
	if n.LastChild != prev {
		return fmt.Errorf("%v: LastChild is %v, want %v", n, n.LastChild, prev)
	}
	return nil
}

var testDataDirs = []string{"testdata/html5lib-tests/tree-construction/", "testdata/go/"}

// parseTestCases yields every case of the tree construction suites.
func parseTestCases(t *testing.T, fn func(file string, i int, ta *testAttrs)) {
	for _, testDataDir := range testDataDirs {
		testFiles, err := filepath.Glob(testDataDir + "*.dat")
		if err != nil {
			t.Fatal(err)
		}
		for _, tf := range testFiles {
			f, err := os.Open(tf)
			if err != nil {
				t.Fatal(err)
			}
			r := bufio.NewReader(f)
			for i := 0; ; i++ {
				ta, err := readParseTest(r)
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatal(err)
				}
				fn(tf, i, ta)
			}
			f.Close()
		}
	}
}

func TestParser(t *testing.T) {
	parseTestCases(t, func(tf string, i int, ta *testAttrs) {
		if parseTestBlacklist[ta.text] {
			return
		}
		testId := fmt.Sprintf("%s/%d", strings.ReplaceAll(tf, string(os.PathSeparator), "/"), i)
		t.Run(testId, func(t *testing.T) {
			if err := testParseCase(ta.text, ta.want, ta.context, ta.scripting); err != nil {
				t.Errorf("%s test #%d %q, %s", tf, i, ta.text, err)
			}
		})
	})
}

// Issue 16318
func TestParserWithoutScripting(t *testing.T) {
	text := `<noscript><img src='https://golang.org/doc/gopher/frontpage.png' /></noscript><p><img src='https://golang.org/doc/gopher/doc.png' /></p>`
	want := `| <html>
|   <head>
|     <noscript>
|   <body>
|     <img>
|       src="https://golang.org/doc/gopher/frontpage.png"
|     <p>
|       <img>
|         src="https://golang.org/doc/gopher/doc.png"
`

	if err := testParseCase(text, want, "", false); err != nil {
		t.Errorf("test with scripting is disabled, %q, %s", text, err)
	}
}

// testParseCase tests one test case from the test files. If the test does not
// pass, it returns an error that explains the failure.
// text is the HTML to be parsed, want is a dump of the correct parse tree,
// and context is the name of the context node, if any.
func testParseCase(text, want, context string, scripting bool) (err error) {
	defer func() {
		if x := recover(); x != nil {
			switch e := x.(type) {
			case error:
				err = e
			default:
				err = fmt.Errorf("%v", e)
			}
		}
	}()

	opts := Options{NoScripting: !scripting}
	var doc *html.Node
	if context == "" {
		doc, err = Parse(strings.NewReader(text), opts)
		if err != nil {
			return err
		}
	} else {
		namespace := ""
		if i := strings.IndexByte(context, ' '); i >= 0 {
			namespace, context = context[:i], context[i+1:]
		}
		contextNode := &html.Node{
			Data:      context,
			DataAtom:  atom.Lookup([]byte(context)),
			Namespace: namespace,
			Type:      html.ElementNode,
		}
		nodes, err := ParseFragment(strings.NewReader(text), contextNode, opts)
		if err != nil {
			return err
		}
		doc = &html.Node{
			Type: html.DocumentNode,
		}
		for _, n := range nodes {
			doc.AppendChild(n)
		}
	}

	if err := checkTreeConsistency(doc); err != nil {
		return err
	}

	got, err := dump(doc)
	if err != nil {
		return err
	}
	// Compare the parsed tree to the #document section.
	if got != want {
		return fmt.Errorf("got vs want:\n----\n%s----\n%s----", got, want)
	}

	if renderTestBlacklist[text] || context != "" {
		return nil
	}

	// Check that rendering and re-parsing results in an identical tree.
	var rendered bytes.Buffer
	if err := html.Render(&rendered, doc); err != nil {
		return err
	}
	doc1, err := ParseBytes(rendered.Bytes(), opts)
	if err != nil {
		return err
	}
	got1, err := dump(doc1)
	if err != nil {
		return err
	}
	if got != got1 {
		return fmt.Errorf("got vs got1:\n----\n%s----\n%s----", got, got1)
	}

	return nil
}

// Some test inputs are simply skipped - we would otherwise fail the test. We
// blacklist such inputs from the parse test.
var parseTestBlacklist = map[string]bool{
	// See the a.Template TODO in inHeadIM.
	`<math><template><mo><template>`:                                     true,
	`<template><svg><foo><template><foreignObject><div></template><div>`: true,
	// We don't support "element insertion steps"
	`<select><button><selectedcontent></button><option>X`:                   true,
	`<select><button><selectedcontent></button><option>x<i>i<b>ib</i>b`:     true,
	`<select><button><selectedcontent></button><option>X<option>Y`:          true,
	`<select><button><selectedcontent></button><option>X<option selected>Y`: true,
}

// Some test input result in parse trees are not 'well-formed' despite
// following the HTML5 recovery algorithms. Rendering and re-parsing such a
// tree will not result in an exact clone of that tree. We blacklist such
// inputs from the render test.
var renderTestBlacklist = map[string]bool{
	// The second <a> will be reparented to the first <table>'s parent. This
	// results in an <a> whose parent is an <a>, which is not 'well-formed'.
	`<a><table><td><a><table></table><a></tr><a></table><b>X</b>C<a>Y`: true,
	// The same thing with a <p>:
	`<p><table></p>`: true,
	// More cases of <a> being reparented:
	`<a href="blah">aba<table><a href="foo">br<tr><td></td></tr>x</table>aoe`: true,
	`<a><table><a></table><p><a><div><a>`:                                     true,
	`<a><table><td><a><table></table><a></tr><a></table><a>`:                  true,
	`<template><a><table><a>`:                                                 true,
	// A similar reparenting situation involving <nobr>:
	`<!DOCTYPE html><body><b><nobr>1<table><nobr></b><i><nobr>2<nobr></i>3`: true,
	// A <plaintext> element is reparented, putting it before a table.
	// A <plaintext> element can't have anything after it in HTML.
	`<table><plaintext><td>`:                                   true,
	`<!doctype html><table><plaintext></plaintext>`:            true,
	`<!doctype html><table><tbody><plaintext></plaintext>`:     true,
	`<!doctype html><table><tbody><tr><plaintext></plaintext>`: true,
	// A form inside a table inside a form doesn't work either.
	`<!doctype html><form><table></form><form></table></form>`: true,
	// A script that ends at EOF may escape its own closing tag when rendered.
	`<!doctype html><script><!--<script `:          true,
	`<!doctype html><script><!--<script <`:         true,
	`<!doctype html><script><!--<script <a`:        true,
	`<!doctype html><script><!--<script </`:        true,
	`<!doctype html><script><!--<script </s`:       true,
	`<!doctype html><script><!--<script </script`:  true,
	`<!doctype html><script><!--<script </scripta`: true,
	`<!doctype html><script><!--<script -`:         true,
	`<!doctype html><script><!--<script -a`:        true,
	`<!doctype html><script><!--<script -<`:        true,
	`<!doctype html><script><!--<script --`:        true,
	`<!doctype html><script><!--<script --a`:       true,
	`<!doctype html><script><!--<script --<`:       true,
	`<script><!--<script `:                         true,
	`<script><!--<script <a`:                       true,
	`<script><!--<script </script`:                 true,
	`<script><!--<script </scripta`:                true,
	`<script><!--<script -`:                        true,
	`<script><!--<script -a`:                       true,
	`<script><!--<script --`:                       true,
	`<script><!--<script --a`:                      true,
	`<script><!--<script <`:                        true,
	`<script><!--<script </`:                       true,
	`<script><!--<script </s`:                      true,
	// Reconstructing the active formatting elements results in a <plaintext>
	// element that contains an <a> element.
	`<!doctype html><p><a><plaintext>b`:                       true,
	`<table><math><select><mi><select></table>`:               true,
	`<!doctype html><table><colgroup><plaintext></plaintext>`: true,
	`<!doctype html><svg><plaintext>a</plaintext>b`:           true,
	// Due to fostering, parsing the rendered output produces a different tree.
	`<math><mtext><table><mglyph><style><img>`: true,
	// Confusing plaintext behavior
	`<!doctype html><table><select><plaintext>a<caption>b`: true,
}

func TestParseFragmentWithNilContext(t *testing.T) {
	// This shouldn't panic.
	nodes, err := ParseFragment(strings.NewReader("<p>x</p>"), nil, Options{})
	if err != nil || len(nodes) == 0 {
		t.Fatalf("ParseFragment(nil context): %v, %d nodes", err, len(nodes))
	}
}

func TestParseFragmentForeignContentTemplates(t *testing.T) {
	srcs := []string{
		"<math><html><template><mn><template></template></template>",
		"<math><math><head><template><mn><template>",
	}
	for _, src := range srcs {
		// The next line shouldn't infinite-loop.
		ParseFragment(strings.NewReader(src), nil, Options{})
	}
}

func TestDepthLimit(t *testing.T) {
	deep := strings.Repeat("<div>", 600)
	if _, err := ParseBytes([]byte(deep), Options{}); !errors.Is(err, ErrTooDeep) {
		t.Errorf("600 open divs: got %v, want ErrTooDeep", err)
	}
	if _, err := ParseBytes([]byte(deep), Options{MaxDepth: 1000}); err != nil {
		t.Errorf("600 open divs within MaxDepth 1000: %v", err)
	}
	if _, err := ParseBytes([]byte(strings.Repeat("<div>", 500)), Options{}); err != nil {
		t.Errorf("500 open divs: %v", err)
	}
}

func TestMaxBufferBytes(t *testing.T) {
	long := "<p>" + strings.Repeat("x", 100) + "</p>"
	if _, err := Parse(strings.NewReader(long), Options{MaxBufferBytes: 64}); !errors.Is(err, ErrTooLarge) {
		t.Errorf("100-byte text with MaxBufferBytes 64: got %v, want ErrTooLarge", err)
	}
	if _, err := ParseBytes([]byte(long), Options{MaxBufferBytes: 64}); !errors.Is(err, ErrTooLarge) {
		t.Errorf("ParseBytes, 100-byte text with MaxBufferBytes 64: got %v, want ErrTooLarge", err)
	}
	if _, err := Parse(strings.NewReader(long), Options{MaxBufferBytes: 128}); err != nil {
		t.Errorf("100-byte text with MaxBufferBytes 128: %v", err)
	}
}

// Text that several tokens make up, and that the parser moves, ends up in
// the node golang.org/x/net/html puts it in, whole.
func TestTextGathering(t *testing.T) {
	for _, in := range []string{
		"a<!---->b</p>c",
		"<table>x<tr>y</tr>z</table>",
		"<b>1<table><b>2</b>3",
		"<a><table><td><a><table></table><a></tr><a></table><b>X</b>C<a>Y",
		"<p>a\x00b<table>c\x00d<tr>e",
		"<svg>\x00a\x00b</svg>c",
		"<pre>\r\nx</pre><textarea>\ny</textarea><listing>\nz",
	} {
		want, err := html.Parse(strings.NewReader(in))
		if err != nil {
			t.Fatal(err)
		}
		got, err := ParseBytes([]byte(in), Options{})
		if err != nil {
			t.Fatal(err)
		}
		w, _ := dump(want)
		g, _ := dump(got)
		if g != w {
			t.Errorf("%q:\n%s\nwant:\n%s", in, g, w)
		}
	}
}
