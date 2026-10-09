package xmltree

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// The trees are read with xmlro; stdTree reads the same trees with
// encoding/xml, which they were read with before. These tests compare the
// two on the markup that tells them apart, and on every XML part of the
// Office documents of the repository.

// stdTree reads a document into a tree with encoding/xml.
func stdTree(data []byte, pick func(*Node) bool) (*Node, error) {
	d := xml.NewDecoder(bytes.NewReader(data))
	d.Strict = false
	for {
		tok, err := d.Token()
		if err == io.EOF {
			return nil, io.ErrUnexpectedEOF
		}
		if err != nil {
			return nil, err
		}
		if start, ok := tok.(xml.StartElement); ok {
			return stdElement(d, start, pick)
		}
	}
}

// stdElement reads the element that start opens from an encoding/xml
// decoder, as ReadFrom reads one.
func stdElement(d *xml.Decoder, start xml.StartElement, pick func(*Node) bool) (*Node, error) {
	node := func(t xml.StartElement) *Node {
		n := &Node{Space: t.Name.Space, Name: t.Name.Local}
		for _, a := range t.Attr {
			n.Attrs = append(n.Attrs, Attr{Name{a.Name.Space, a.Name.Local}, a.Value})
		}
		return n
	}
	root := node(start)
	stack := []*Node{root}
	for len(stack) > 0 {
		tok, err := d.Token()
		if err == io.EOF {
			return nil, io.ErrUnexpectedEOF
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			n := node(t)
			p := stack[len(stack)-1]
			p.Kids = append(p.Kids, n)
			stack = append(stack, n)
		case xml.EndElement:
			if n := stack[len(stack)-1]; len(n.Kids) == 0 {
				n.runs = nil
			}
			stack = stack[:len(stack)-1]
		case xml.CharData:
			stack[len(stack)-1].addText(string(t))
		case xml.ProcInst:
			if t.Target == "ACE" {
				if code, err := strconv.ParseUint(strings.TrimSpace(string(t.Inst)), 16, 32); err == nil && code < 0x110000 {
					stack[len(stack)-1].addText(string(rune(code)))
				}
			}
		}
	}
	root.resolveAlternates(pick)
	return root, nil
}

// same reports where two trees differ ("" when they do not): names,
// namespaces, attributes, text and its runs, in order.
func same(a, b *Node, path string) string {
	switch {
	case a == nil || b == nil:
		if a != b {
			return path + ": one tree has no element here"
		}
		return ""
	case a.Space != b.Space || a.Name != b.Name:
		return path + ": " + a.Space + ":" + a.Name + " vs " + b.Space + ":" + b.Name
	case a.Text != b.Text:
		return path + ": text " + a.Text + " vs " + b.Text
	case len(a.Attrs) != len(b.Attrs):
		return path + ": attributes differ in number"
	case len(a.Kids) != len(b.Kids):
		return path + ": children differ in number"
	}
	for i := range a.Attrs {
		if a.Attrs[i] != b.Attrs[i] {
			return path + ": attribute " + a.Attrs[i].Name.Local + " vs " + b.Attrs[i].Name.Local
		}
	}
	as, bs := a.Segments(), b.Segments()
	if len(as) != len(bs) {
		return path + ": segments differ in number"
	}
	for i := range as {
		if as[i].Text != bs[i].Text || (as[i].Elem == nil) != (bs[i].Elem == nil) {
			return path + ": segments differ"
		}
	}
	for i := range a.Kids {
		if d := same(a.Kids[i], b.Kids[i], path+"/"+a.Kids[i].Name); d != "" {
			return d
		}
	}
	return ""
}

func TestTreeMatchesEncodingXML(t *testing.T) {
	for name, doc := range map[string]string{
		"namespaces": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://w" xmlns:r="http://r" xmlns="http://default" mc:Ignorable="w14" xmlns:mc="http://mc"><w:body><p a="1" w:b="2" r:id="rId1" xml:space="preserve">text</p><q:x xmlns:q="http://q" q:y="3"/><unbound:z unbound:a="1">t</unbound:z></w:body></w:document>`,
		"text":      "<r>a&amp;b &lt;c&gt; &#65;&#x42; &unknown; <![CDATA[<raw>&amp;]]> line\r\nend\rtoo<k/>tail<!-- comment --><?pi target?>after</r>",
		"mixed":     `<r>one<a>two<b/>three</a>four<c/>five</r>`,
		"empty":     `<r></r>`,
		"self":      `<r><a/><b></b><c x=""/></r>`,
		"bom":       "\xef\xbb\xbf<r><a>x</a></r>",
		"doctype":   `<!DOCTYPE r><r><a/></r>`,
		"deepish":   strings.Repeat("<a>", 100) + "x" + strings.Repeat("</a>", 100),
		"alternate": `<r xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006"><mc:AlternateContent><mc:Choice Requires="w14"><new/></mc:Choice><mc:Fallback><old/></mc:Fallback></mc:AlternateContent>t<mc:AlternateContent><mc:Choice Requires="x"><nope/></mc:Choice></mc:AlternateContent></r>`,
		"attrs":     `<r a='single' b="dou&quot;ble" c="&#10;nl" d="  spaced  " e="tab	tab"/>`,
	} {
		for _, pick := range []func(*Node) bool{nil, SupportedChoice(func(p string) bool { return p == "w14" })} {
			want, err := stdTree([]byte(doc), pick)
			if err != nil {
				t.Fatalf("%s: encoding/xml: %v", name, err)
			}
			got, err := ParsePicking([]byte(doc), pick)
			if err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			if d := same(want, got, want.Name); d != "" {
				t.Errorf("%s: %s", name, d)
			}
			got, err = ParsePickingReader(strings.NewReader(doc), pick)
			if err != nil {
				t.Fatalf("%s (reader): %v", name, err)
			}
			if d := same(want, got, want.Name); d != "" {
				t.Errorf("%s (reader): %s", name, d)
			}
		}
	}
}

func TestTreeErrors(t *testing.T) {
	for name, doc := range map[string]string{
		"truncated": `<r><a>`,
		"mismatch":  `<r><a></b></r>`,
		"empty":     ``,
		"no root":   `<?xml version="1.0"?>`,
		"too deep":  strings.Repeat("<a>", MaxDepth+1) + strings.Repeat("</a>", MaxDepth+1),
	} {
		if _, err := Parse([]byte(doc)); err == nil {
			t.Errorf("%s: parsed", name)
		}
	}
	if _, err := Parse([]byte(strings.Repeat("<a>", MaxDepth) + strings.Repeat("</a>", MaxDepth))); err != nil {
		t.Errorf("nesting %d deep: %v", MaxDepth, err)
	}
}

// TestTreesOfDocumentsMatch compares the trees of every XML part of the
// Office documents in the repository, read both ways.
func TestTreesOfDocumentsMatch(t *testing.T) {
	var files []string
	for _, pattern := range []string{"../../converter/pptx/testdata/*.pptx", "../../converter/docx/testdata/*.docx", "../../converter/xlsx/testdata/*.xlsx", "../../converter/visio/testdata/*.vsdx"} {
		m, _ := filepath.Glob(pattern)
		files = append(files, m...)
	}
	if len(files) == 0 {
		t.Skip("no Office documents found")
	}
	parts := 0
	for _, name := range files {
		zr, err := zip.OpenReader(name)
		if err != nil {
			continue // an encrypted document is not a zip
		}
		for _, f := range zr.File {
			if !strings.HasSuffix(f.Name, ".xml") && !strings.HasSuffix(f.Name, ".rels") {
				continue
			}
			rc, err := f.Open()
			if err != nil {
				t.Fatal(err)
			}
			data, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				t.Fatal(err)
			}
			want, err := stdTree(data, nil)
			if err != nil {
				continue // what encoding/xml could not read either
			}
			got, err := ParsePickingReader(bytes.NewReader(data), nil)
			if err != nil {
				t.Errorf("%s %s: %v", name, f.Name, err)
				continue
			}
			if d := same(want, got, want.Name); d != "" {
				t.Errorf("%s %s: %s", name, f.Name, d)
			}
			parts++
		}
		zr.Close()
	}
	t.Logf("%d parts of %d documents compared", parts, len(files))
}
