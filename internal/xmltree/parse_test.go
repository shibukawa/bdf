package xmltree

import (
	"reflect"
	"strings"
	"testing"
)

func checkNode(t *testing.T, got, want *Node) {
	t.Helper()
	if got == nil || got.Space != want.Space || got.Name != want.Name || got.Text != want.Text || !reflect.DeepEqual(got.Attrs, want.Attrs) || len(got.Kids) != len(want.Kids) {
		t.Fatalf("node: got %+v, want %+v", got, want)
	}
	for i := range want.Kids {
		checkNode(t, got.Kids[i], want.Kids[i])
	}
}

func TestParseContent(t *testing.T) {
	for _, tc := range []struct {
		name, doc string
		want      *Node
	}{
		{"namespaces", `<r xmlns="urn:root" xmlns:w="urn:word"><w:p plain="1" w:val="2" xml:space="preserve"/><u:x u:y="3"/><none xmlns=""/></r>`, &Node{
			Space: "urn:root", Name: "r", Attrs: []Attr{{Name{"", "xmlns"}, "urn:root"}, {Name{"xmlns", "w"}, "urn:word"}},
			Kids: []*Node{
				{Space: "urn:word", Name: "p", Attrs: []Attr{{Name{"", "plain"}, "1"}, {Name{"urn:word", "val"}, "2"}, {Name{"http://www.w3.org/XML/1998/namespace", "space"}, "preserve"}}},
				{Space: "u", Name: "x", Attrs: []Attr{{Name{"u", "y"}, "3"}}},
				{Name: "none", Attrs: []Attr{{Name{"", "xmlns"}, ""}}},
			},
		}},
		{"text", "<r>a&amp;b &lt;c&gt; &#65;&#x42; &unknown; <![CDATA[<raw>&amp;]]> line\r\nend\rtoo<!-- ignored --><?pi ignored?><?ACE 1F600?></r>", &Node{Name: "r", Text: "a&b <c> AB &unknown; <raw>&amp; line\nend\ntoo😀"}},
		{"attributes", `<r a='single' b="dou&quot;ble" c="&#10;nl" d="  spaced  " e="tab	tab"/>`, &Node{Name: "r", Attrs: []Attr{{Name{"", "a"}, "single"}, {Name{"", "b"}, `dou"ble`}, {Name{"", "c"}, "\nnl"}, {Name{"", "d"}, "  spaced  "}, {Name{"", "e"}, "tab\ttab"}}}},
		{"mixed", `<r>one<a>two<b/>three</a>four<c/>five</r>`, &Node{Name: "r", Text: "onefourfive", Kids: []*Node{{Name: "a", Text: "twothree", Kids: []*Node{{Name: "b"}}}, {Name: "c"}}}},
		{"empty", `<r><a/><b></b><c x=""/></r>`, &Node{Name: "r", Kids: []*Node{{Name: "a"}, {Name: "b"}, {Name: "c", Attrs: []Attr{{Name{"", "x"}, ""}}}}}},
		{"bom and doctype", "\xef\xbb\xbf<?xml version=\"1.0\"?><!DOCTYPE r><r><a>x</a></r>", &Node{Name: "r", Kids: []*Node{{Name: "a", Text: "x"}}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, mode := range []string{"bytes", "reader"} {
				t.Run(mode, func(t *testing.T) {
					var got *Node
					var err error
					if mode == "bytes" {
						got, err = Parse([]byte(tc.doc))
					} else {
						got, err = ParsePickingReader(strings.NewReader(tc.doc), nil)
					}
					if err != nil {
						t.Fatal(err)
					}
					checkNode(t, got, tc.want)
					if tc.name == "mixed" {
						want := []Segment{{Text: "one"}, {Elem: got.Kids[0]}, {Text: "four"}, {Elem: got.Kids[1]}, {Text: "five"}}
						if !reflect.DeepEqual(got.Segments(), want) {
							t.Errorf("mixed content: got %+v, want %+v", got.Segments(), want)
						}
					}
				})
			}
		})
	}
}

func TestParseAlternateContent(t *testing.T) {
	const doc = `<r xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006">before<mc:AlternateContent><mc:Choice Requires="w14 extra"><new/></mc:Choice><mc:Fallback><old/></mc:Fallback></mc:AlternateContent>after<mc:AlternateContent><mc:Choice Requires="x"><only/></mc:Choice></mc:AlternateContent></r>`
	for _, tc := range []struct {
		name      string
		supported func(string) bool
		want      string
	}{
		{"fallback", nil, "old"},
		{"unsupported", func(string) bool { return false }, "old"},
		{"partial support", func(p string) bool { return p == "w14" }, "old"},
		{"supported", func(p string) bool { return p == "w14" || p == "extra" }, "new"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, mode := range []string{"bytes", "reader"} {
				var got *Node
				var err error
				pick := SupportedChoice(tc.supported)
				if mode == "bytes" {
					got, err = ParsePicking([]byte(doc), pick)
				} else {
					got, err = ParsePickingReader(strings.NewReader(doc), pick)
				}
				if err != nil {
					t.Fatal(err)
				}
				if len(got.Kids) != 2 || got.Kids[0].Name != tc.want || got.Kids[1].Name != "only" {
					t.Fatalf("%s: alternate content: %+v", mode, got.Kids)
				}
				want := []Segment{{Text: "before"}, {Elem: got.Kids[0]}, {Text: "after"}, {Elem: got.Kids[1]}}
				if !reflect.DeepEqual(got.Segments(), want) {
					t.Errorf("%s: mixed content after selection: %+v", mode, got.Segments())
				}
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
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
