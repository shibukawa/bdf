package xmp

import (
	"runtime"
	"slices"
	"strings"
	"testing"
)

func TestDublinCore(t *testing.T) {
	packet := `<?xpacket begin="` + "\ufeff" + `" id="W5M0MpCehiHzreSzNTczkc9d"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
 <rdf:Description rdf:about="" xmlns:xmp="http://ns.adobe.com/xap/1.0/" xmp:CreateDate="2026-09-26T12:00:00+09:00"
   xmlns:dc="http://purl.org/dc/elements/1.1/" dc:format="image/vnd.adobe.photoshop">
  <dc:title><rdf:Alt><rdf:li xml:lang="en">Poster</rdf:li><rdf:li xml:lang="x-default">ポスター</rdf:li></rdf:Alt></dc:title>
  <dc:creator><rdf:Seq><rdf:li>Ann</rdf:li><rdf:li>Bo</rdf:li></rdf:Seq></dc:creator>
  <dc:subject><rdf:Bag><rdf:li>red</rdf:li><rdf:li>blue</rdf:li></rdf:Bag></dc:subject>
  <dc:description><rdf:Alt><rdf:li xml:lang="en">A test</rdf:li></rdf:Alt></dc:description>
  <xmp:ModifyDate>2026-09-27T08:00:00Z</xmp:ModifyDate>
  <xmp:CreatorTool>ignored</xmp:CreatorTool>
 </rdf:Description>
</rdf:RDF></x:xmpmeta><?xpacket end="w"?>`
	dc := DublinCore([]byte(packet))
	for _, c := range []struct {
		name string
		got  []string
		want []string
	}{
		{"title", dc.Title, []string{"ポスター"}},
		{"creator", dc.Creator, []string{"Ann", "Bo"}},
		{"subject", dc.Subject, []string{"red", "blue"}},
		{"description", dc.Description, []string{"A test"}},
		{"format", dc.Format, []string{"image/vnd.adobe.photoshop"}},
		{"created", dc.Created, []string{"2026-09-26T12:00:00+09:00"}},
		{"modified", dc.Modified, []string{"2026-09-27T08:00:00Z"}},
	} {
		if !slices.Equal(c.got, c.want) {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}
	if !DublinCore([]byte("<not xmp")).IsZero() {
		t.Error("malformed packet gave values")
	}
}

// TestDates: the created date comes from the original's date (Photoshop,
// EXIF) before XMP's CreateDate (the digitizing); dates are normalized, and
// the dc schema has no created element.
func TestDates(t *testing.T) {
	packet := `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
 <rdf:Description xmlns:xmp="http://ns.adobe.com/xap/1.0/" xmlns:photoshop="http://ns.adobe.com/photoshop/1.0/"
   xmlns:dc="http://purl.org/dc/elements/1.1/" xmp:CreateDate="2026-09-03T12:00:00Z" xmp:ModifyDate="2026:09:04 08:00:00+0900"
   photoshop:DateCreated="2026-09-01">
  <dc:date><rdf:Seq><rdf:li>2026-09</rdf:li><rdf:li>yesterday</rdf:li></rdf:Seq></dc:date>
  <dc:created>2026-01-01</dc:created>
 </rdf:Description></rdf:RDF></x:xmpmeta>`
	dc := DublinCore([]byte(packet))
	if !slices.Equal(dc.Created, []string{"2026-09-01"}) || !slices.Equal(dc.Modified, []string{"2026-09-04T08:00:00+09:00"}) ||
		!slices.Equal(dc.Date, []string{"2026-09"}) {
		t.Errorf("created %q, modified %q, date %q", dc.Created, dc.Modified, dc.Date)
	}
	for _, c := range []struct{ in, want string }{
		{"2026", "2026"},
		{"2026-09-01T10:20", "2026-09-01T10:20"},
		{"2026-09-01T10:20:30.25+0900", "2026-09-01T10:20:30.25+09:00"},
		{"2026:09:01 10:20:30Z", "2026-09-01T10:20:30Z"},
		{"2026-13-01", ""},
		{"0000", ""},
	} {
		if got := Date(c.in); got != c.want {
			t.Errorf("Date(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestRDF: typed node elements and resources described in place
// (Inkscape's metadata in SVG files), named resources, and the TIFF schema
// filling in what the dc schema leaves out; a camera's placeholder
// description is no description.
func TestRDF(t *testing.T) {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
  xmlns:cc="http://creativecommons.org/ns#" xmlns:dc="http://purl.org/dc/elements/1.1/"><metadata><rdf:RDF>
  <cc:Work rdf:about=""><dc:title>Shapes</dc:title><dc:creator><cc:Agent><dc:title>Hanako Sato</dc:title></cc:Agent></dc:creator>
   <dc:source rdf:resource="https://example.com/src"/></cc:Work></rdf:RDF></metadata><rect/></svg>`
	dc := DublinCore([]byte(svg))
	if dc.Title.First() != "Shapes" || dc.Creator.First() != "Hanako Sato" || dc.Source.First() != "https://example.com/src" {
		t.Errorf("svg: %+v", dc)
	}
	packet := `<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description
  xmlns:tiff="http://ns.adobe.com/tiff/1.0/" xmlns:dc="http://purl.org/dc/elements/1.1/"
  tiff:Artist="Ann; Bob" tiff:Copyright="CC0" tiff:ImageDescription="A cat">
  <dc:description><rdf:Alt><rdf:li xml:lang="x-default">OLYMPUS DIGITAL CAMERA</rdf:li></rdf:Alt></dc:description>
  <dc:rights><rdf:Alt><rdf:li xml:lang="x-default">© Ann</rdf:li></rdf:Alt></dc:rights></rdf:Description></rdf:RDF>`
	dc = DublinCore([]byte(packet))
	if !slices.Equal(dc.Creator, []string{"Ann", "Bob"}) || !slices.Equal(dc.Rights, []string{"© Ann"}) || !slices.Equal(dc.Description, []string{"A cat"}) {
		t.Errorf("tiff: %+v", dc)
	}
}

// TestSubtreeDepth: elements are read as deep as maxDepth, and what follows
// the deeper ones is read.
func TestSubtreeDepth(t *testing.T) {
	const deep = maxDepth + 100
	data := "<r>" + strings.Repeat("<a>", deep) + "text" + strings.Repeat("</a>", deep) + "<b>after</b></r>"
	r := Reader([]byte(data))
	if _, err := r.Next(); err != nil {
		t.Fatal(err)
	}
	root, err := Subtree(r)
	if err != nil || len(root.Children) != 2 || root.Children[1].Text != "after" {
		t.Fatalf("%v, children %+v", err, root.Children)
	}
	depth := 0
	n := root
	for ; len(n.Children) > 0; n = n.Children[0] {
		depth++
	}
	if depth != maxDepth || n.Text != "" {
		t.Errorf("elements %d deep, want %d; the last has the text %q", depth, maxDepth, n.Text)
	}
	// without the elements that end them
	if !DublinCore([]byte(strings.Repeat("<a>", deep))).IsZero() {
		t.Error("values of elements without end")
	}
}

// TestDefaultAlternatives: the default of a language alternative comes
// first, the last of several before the others; many of them are put in
// order once (200 MB of copies for these before).
func TestDefaultAlternatives(t *testing.T) {
	alt := func(items string) *Node {
		data := `<dc:title xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Alt>` + items + `</rdf:Alt></dc:title>`
		r := Reader([]byte(data))
		if _, err := r.Next(); err != nil {
			t.Fatal(err)
		}
		n, err := Subtree(r)
		if err != nil {
			t.Fatal(err)
		}
		return n
	}
	n := alt(`<rdf:li xml:lang="en">A</rdf:li><rdf:li xml:lang="x-default">B</rdf:li><rdf:li xml:lang="fr">C</rdf:li><rdf:li xml:lang="x-default">D</rdf:li>`)
	if got := values(n); !slices.Equal(got, []string{"D", "B", "A", "C"}) {
		t.Errorf("values %q", got)
	}
	if got := values(alt("")); got != nil {
		t.Errorf("values of no item: %q", got)
	}
	const items = 5000
	n = alt(strings.Repeat(`<rdf:li xml:lang="x-default">t</rdf:li>`, items))
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	got := values(n)
	runtime.ReadMemStats(&after)
	if len(got) != items || after.TotalAlloc-before.TotalAlloc > 16<<20 {
		t.Errorf("%d values, %d bytes allocated", len(got), after.TotalAlloc-before.TotalAlloc)
	}
}
