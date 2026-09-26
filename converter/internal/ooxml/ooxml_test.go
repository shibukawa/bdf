package ooxml

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func testPackage(t *testing.T, parts map[string]string) *Package {
	t.Helper()
	var b bytes.Buffer
	zw := zip.NewWriter(&b)
	for name, data := range parts {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(data))
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	p, err := Open(bytes.NewReader(b.Bytes()), int64(b.Len()))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

const relsNS = `xmlns="http://schemas.openxmlformats.org/package/2006/relationships"`

func TestPackage(t *testing.T) {
	p := testPackage(t, map[string]string{
		"_rels/.rels":       `<Relationships ` + relsNS + `><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="/word/document.xml"/></Relationships>`,
		"word/document.xml": `<w:document xmlns:w="w"/>`,
		"word/_rels/document.xml.rels": `<Relationships ` + relsNS + `>` +
			`<Relationship Id="rId10" Type="http://x/relationships/image" Target="media/b.png"/>` +
			`<Relationship Id="rId2" Type="http://x/relationships/image" Target="../media/a.png"/>` +
			`<Relationship Id="rId3" Type="http://x/relationships/hyperlink" Target="https://example.com/" TargetMode="External"/>` +
			`</Relationships>`,
		"Media/A.png": "png",
	})
	r, ok := p.RelOfType("", "/officeDocument")
	if !ok || r.Target != "word/document.xml" {
		t.Fatalf("officeDocument = %+v, %v", r, ok)
	}
	// rId2 comes before rId10
	if r, ok := p.RelOfType(r.Target, "/image"); !ok || r.ID != "rId2" || r.Target != "media/a.png" {
		t.Errorf("first image = %+v, %v", r, ok)
	}
	if r, ok := p.Target("word/document.xml", "rId10"); !ok || r.Target != "word/media/b.png" {
		t.Errorf("rId10 = %+v, %v", r, ok)
	}
	if r := p.Rels("word/document.xml")["rId3"]; !r.External || r.Target != "https://example.com/" {
		t.Errorf("external link = %+v", r)
	}
	if _, ok := p.Target("word/document.xml", "rId9"); ok {
		t.Error("unknown relationship found")
	}
	// part names are case-insensitive
	if !p.Has("/media/a.png") {
		t.Error("Has(/media/a.png) = false")
	}
	if b, err := p.Read("media/a.png"); err != nil || string(b) != "png" {
		t.Errorf("Read = %q, %v", b, err)
	}
	if _, err := p.Read("missing.xml"); err == nil {
		t.Error("reading a missing part succeeded")
	}
	n1, err := p.XML("word/document.xml")
	if err != nil {
		t.Fatal(err)
	}
	if n2, _ := p.XML("/WORD/document.xml"); n2 != n1 {
		t.Error("XML is not cached")
	}
}

func TestParse(t *testing.T) {
	n, err := Parse([]byte(`<p:sp xmlns:p="p" xmlns:a="a" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006">` +
		`<a:t>text</a:t>` +
		`<mc:AlternateContent><mc:Choice Requires="x"><a:new/></mc:Choice><mc:Fallback><a:old/><a:old/></mc:Fallback></mc:AlternateContent>` +
		`<a:blip r:embed="rId4" embed="plain" amt="50000" pct="25%" x="12700" y="1in" b="on" i="0" n="3.7"/>` +
		`</p:sp>`))
	if err != nil {
		t.Fatal(err)
	}
	if n.Name != "sp" || n.Space != "p" {
		t.Errorf("root = %s %s", n.Space, n.Name)
	}
	var names []string
	for _, k := range n.Elements() {
		names = append(names, k.Name)
	}
	if got := len(n.Children("old")); got != 2 || n.Child("new") != nil {
		t.Errorf("alternate content not replaced by its fallback: %v", names)
	}
	if got := n.Child("t").Content(); got != "text" {
		t.Errorf("text = %q", got)
	}
	b := n.Path("blip")
	if got := b.RelID("embed"); got != "rId4" {
		t.Errorf("RelID = %q", got)
	}
	if got, _ := b.Attr("embed"); got != "plain" {
		t.Errorf("Attr(embed) = %q, the relationship attribute leaks in", got)
	}
	for _, c := range []struct {
		name      string
		got, want float64
	}{
		{"AttrPct", b.AttrPct("amt", 0), 0.5},
		{"AttrPct %", b.AttrPct("pct", 0), 0.25},
		{"AttrPct default", b.AttrPct("none", 0.75), 0.75},
		{"AttrEMU", b.AttrEMU("x", 0), 1},
		{"AttrEMU in", b.AttrEMU("y", 0), 72},
		{"AttrInt", float64(b.AttrInt("n", 0)), 3},
		{"AttrFloat", b.AttrFloat("n", 0), 3.7},
	} {
		if c.got != c.want {
			t.Errorf("%s = %g, want %g", c.name, c.got, c.want)
		}
	}
	if !b.AttrBool("b", false) || b.AttrBool("i", true) || !b.AttrBool("none", true) {
		t.Error("AttrBool")
	}
	// nil elements read as empty
	var none *Node
	if none.Path("a", "b") != nil || none.Content() != "" || none.Elements() != nil || none.AttrStr("x", "d") != "d" || none.RelID("id") != "" ||
		none.Segments() != nil {
		t.Error("nil node accessors")
	}
}

func TestSegments(t *testing.T) {
	// Visio text: markers among the characters they format
	n, err := Parse([]byte(`<Text xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006"><cp IX="0"/>one <cp IX="1"/>two` +
		`<mc:AlternateContent><mc:Fallback><fld/><fld/></mc:Fallback></mc:AlternateContent>&amp;<pp/>three</Text>`))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, s := range n.Segments() {
		if s.Elem != nil {
			got = append(got, "<"+s.Elem.Name+">")
		} else {
			got = append(got, s.Text)
		}
	}
	want := []string{"<cp>", "one ", "<cp>", "two", "<fld>", "<fld>", "&", "<pp>", "three"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("segments = %q, want %q", got, want)
	}
	if n.Content() != "one two&three" {
		t.Errorf("content = %q", n.Content())
	}
	// text alone, and an empty element
	n, _ = Parse([]byte(`<t>plain</t>`))
	if s := n.Segments(); len(s) != 1 || s[0].Text != "plain" || s[0].Elem != nil {
		t.Errorf("text-only segments = %+v", s)
	}
	n, _ = Parse([]byte(`<t/>`))
	if s := n.Segments(); s != nil {
		t.Errorf("empty element segments = %+v", s)
	}
}

func TestNilPackage(t *testing.T) {
	var p *Package
	if p.Has("a.xml") || len(p.Rels("")) != 0 {
		t.Error("a nil package has parts")
	}
	if _, err := p.XML("a.xml"); err == nil {
		t.Error("a nil package reads a part")
	}
	if _, ok := p.RelOfType("", "/theme"); ok {
		t.Error("a nil package has relationships")
	}
	if dc := p.CoreProperties(); len(dc.Title) != 0 {
		t.Error("a nil package has core properties")
	}
}

func TestResolvePart(t *testing.T) {
	for _, c := range [][3]string{
		{"ppt/slides/", "../media/image1.png", "ppt/media/image1.png"},
		{"ppt/", "slides/slide1.xml", "ppt/slides/slide1.xml"},
		{"ppt/slides/", "/ppt/theme/theme1.xml", "ppt/theme/theme1.xml"},
		{"", "word/document.xml", "word/document.xml"},
	} {
		if got := ResolvePart(c[0], c[1]); got != c[2] {
			t.Errorf("ResolvePart(%q, %q) = %q, want %q", c[0], c[1], got, c[2])
		}
	}
}

func TestCoreProperties(t *testing.T) {
	p := testPackage(t, map[string]string{
		"_rels/.rels": `<Relationships ` + relsNS + `><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/></Relationships>`,
		"docProps/core.xml": `<cp:coreProperties xmlns:cp="cp" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/">` +
			`<dc:title>Report</dc:title><dc:creator>Alice</dc:creator><cp:keywords>alpha; beta</cp:keywords>` +
			`<dc:language>ja-JP</dc:language><dcterms:created>2026-01-02T03:04:05Z</dcterms:created></cp:coreProperties>`,
	})
	dc := p.CoreProperties()
	if dc.Title.First() != "Report" || dc.Creator.First() != "Alice" || dc.Language.First() != "ja-JP" || dc.Created.First() != "2026-01-02T03:04:05Z" {
		t.Errorf("core properties = %+v", dc)
	}
	if len(dc.Subject) != 2 || dc.Subject[0] != "alpha" || dc.Subject[1] != "beta" {
		t.Errorf("keywords = %q", dc.Subject)
	}
	if dc := testPackage(t, map[string]string{"a.xml": "<a/>"}).CoreProperties(); dc.Title != nil {
		t.Errorf("no core properties: %+v", dc)
	}
}
