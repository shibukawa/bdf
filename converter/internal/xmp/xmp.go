// Package xmp reads the Dublin Core of XMP metadata packets (ISO 16684-1),
// as image and design files embed them, and of other RDF/XML, such as the
// metadata of SVG documents (docs/spec.md §4.3). Besides the dc schema, it
// takes the created and modified dates from the XMP, Photoshop, EXIF,
// TIFF and DCMI Terms schemas, normalized to the W3C date and time format,
// and the description, creators and rights from the TIFF schema (EXIF
// mirrored in XMP) when the dc schema has none.
package xmp

import (
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/shibukawa/tinygodriver/encoding/xmlro"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/xmltree"
)

// Namespaces of the metadata read.
const (
	nsRDF       = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	nsXML       = "http://www.w3.org/XML/1998/namespace"
	nsDC        = "http://purl.org/dc/elements/1.1/"
	nsDCTerms   = "http://purl.org/dc/terms/"
	nsXMP       = "http://ns.adobe.com/xap/1.0/"
	nsPhotoshop = "http://ns.adobe.com/photoshop/1.0/"
	nsEXIF      = "http://ns.adobe.com/exif/1.0/"
	nsTIFF      = "http://ns.adobe.com/tiff/1.0/"
)

// DublinCore reads the Dublin Core of an XMP packet, or of any XML document
// that holds RDF/XML. A language alternative gives its default (x-default)
// value, or its first. A malformed packet gives what was read before the
// error.
func DublinCore(packet []byte) bdf.DublinCore {
	r := Reader(packet)
	for {
		k, err := r.Next()
		if err != nil || k == xmlro.EOF {
			return bdf.DublinCore{}
		}
		if k == xmlro.StartElement {
			root, _ := Subtree(r)
			return RDF(root)
		}
	}
}

// Node is an XML element.
type Node struct {
	Name     xmltree.Name
	Attrs    []xmltree.Attr
	Children []*Node
	// Text is the character data directly in the element.
	Text string
}

// Attr returns an attribute's value ("" when absent).
func (n *Node) Attr(space, local string) string {
	for _, a := range n.Attrs {
		if a.Name.Local == local && a.Name.Space == space {
			return a.Value
		}
	}
	return ""
}

// Is reports whether the element has the name.
func (n *Node) Is(space, local string) bool { return n.Name.Space == space && n.Name.Local == local }

// Reader makes an XML reader that tolerates what browsers tolerate in SVG
// files and what image files carry: markup that is not well formed,
// entities (those of an internal DTD subset, as Illustrator declares its
// namespaces, within xmltree.MaxEntityText), and charsets other than UTF-8
// (UTF-16 and Latin-1 are decoded; others keep their ASCII and replace the
// rest). Elements may be readerDepth deep, far deeper than Subtree keeps
// them, so that what follows the deep ones is read.
func Reader(data []byte) *xmlro.Reader {
	return xmltree.Open(data, xmlro.Options{Lenient: true, MaxDepth: readerDepth})
}

// maxDepth bounds the elements within one another that Subtree reads:
// metadata and drawings nest a few dozen deep, a damaged file as deep as
// it is long. readerDepth bounds those the reader passes over, which it
// keeps a few bytes for.
const (
	maxDepth    = 512
	readerDepth = 1 << 20
)

// Subtree reads the element whose start a reader is on, up to its end.
// What was read before an error is returned with it. An element maxDepth
// deep is read without what it contains.
func Subtree(r *xmlro.Reader) (*Node, error) {
	return subtree(r, 0)
}

func subtree(r *xmlro.Reader, depth int) (*Node, error) {
	n := &Node{Name: xmltree.ElementName(r), Attrs: xmltree.Attrs(r)}
	if depth >= maxDepth {
		return n, r.Skip()
	}
	var text []byte
	for {
		k, err := r.Next()
		if err == nil && k == xmlro.EOF {
			err = io.ErrUnexpectedEOF
		}
		if err != nil {
			n.Text = string(text)
			return n, err
		}
		switch k {
		case xmlro.StartElement:
			c, err := subtree(r, depth+1)
			n.Children = append(n.Children, c)
			if err != nil {
				return n, err
			}
		case xmlro.Text, xmlro.CData:
			text = xmltree.AppendText(text, r)
		case xmlro.EndElement:
			n.Text = string(text)
			return n, nil
		}
	}
}

// RDF reads the Dublin Core of the RDF/XML (rdf:RDF elements) under an
// element: an XMP packet, the metadata of an SVG document. Node elements
// may be typed (Inkscape's cc:Work), and a property's value may be a
// resource described in place (cc:Agent) or named (rdf:resource).
func RDF(root *Node) bdf.DublinCore {
	var dc, tiffDC bdf.DublinCore
	var created, digitized []string
	prop := func(space, local string, values []string) {
		switch space {
		case nsDC:
			if f := dc.Field(local); f != nil && local != "created" && local != "modified" {
				if local == "date" {
					values = dates(values)
				}
				// an alternative text keeps its default only
				if local == "title" || local == "description" || local == "rights" {
					values = first(values)
				}
				add(f, values...)
			}
		case nsDCTerms:
			switch local {
			case "created":
				add(&dc.Created, dates(values)...)
			case "modified":
				add(&dc.Modified, dates(values)...)
			}
		case nsPhotoshop:
			if local == "DateCreated" {
				created = append(created, values...)
			}
		case nsEXIF:
			switch local {
			case "DateTimeOriginal":
				created = append(created, values...)
			case "DateTimeDigitized":
				digitized = append(digitized, values...)
			}
		case nsXMP:
			switch local {
			case "CreateDate":
				digitized = append(digitized, values...)
			case "ModifyDate":
				add(&dc.Modified, dates(values)...)
			}
		case nsTIFF:
			switch local {
			case "ImageDescription":
				add(&tiffDC.Description, first(values)...)
			case "Artist":
				for _, v := range values {
					add(&tiffDC.Creator, SplitList(v)...)
				}
			case "Copyright":
				add(&tiffDC.Rights, first(values)...)
			case "DateTime":
				add(&tiffDC.Modified, dates(values)...)
			}
		}
	}
	var walk func(n *Node)
	walk = func(n *Node) {
		if n == nil {
			return
		}
		if !n.Is(nsRDF, "RDF") {
			for _, c := range n.Children {
				walk(c)
			}
			return
		}
		// node elements (rdf:Description, cc:Work …): their properties are
		// attributes and child elements
		for _, desc := range n.Children {
			for _, a := range desc.Attrs {
				if a.Name.Space != nsRDF && a.Name.Space != "xmlns" && a.Name.Space != nsXML {
					prop(a.Name.Space, a.Name.Local, []string{a.Value})
				}
			}
			for _, p := range desc.Children {
				prop(p.Name.Space, p.Name.Local, values(p))
			}
		}
	}
	walk(root)
	// the original's date, else the digitized one
	if len(dc.Created) == 0 {
		if dc.Created = first(dates(created)); len(dc.Created) == 0 {
			dc.Created = first(dates(digitized))
		}
	}
	if len(dc.Description) > 0 && Placeholder(dc.Description[0]) {
		dc.Description = nil
	}
	if len(tiffDC.Description) > 0 && Placeholder(tiffDC.Description[0]) {
		tiffDC.Description = nil
	}
	Merge(&dc, tiffDC)
	return dc
}

// values returns the values of an RDF property element: the items of an
// rdf:Alt (the default language first), rdf:Seq or rdf:Bag, the resource it
// names, the title of the resource it describes (Inkscape's cc:Agent), or
// its text.
func values(p *Node) []string {
	if r := p.Attr(nsRDF, "resource"); r != "" {
		return []string{r}
	}
	if v := p.Attr(nsRDF, "value"); v != "" {
		return []string{v}
	}
	for _, c := range p.Children {
		switch {
		case c.Is(nsRDF, "Alt"), c.Is(nsRDF, "Seq"), c.Is(nsRDF, "Bag"):
			// The default of an alternative comes first (the last of
			// several defaults before the others).
			var defaults [][]string
			var others []string
			for _, li := range c.Children {
				if !li.Is(nsRDF, "li") {
					continue
				}
				v := values(li)
				if c.Is(nsRDF, "Alt") && li.Attr(nsXML, "lang") == "x-default" {
					defaults = append(defaults, v)
				} else {
					others = append(others, v...)
				}
			}
			var out []string
			for i := len(defaults) - 1; i >= 0; i-- {
				out = append(out, defaults[i]...)
			}
			return append(out, others...)
		case c.Is(nsRDF, "value"):
			return []string{collapse(c.Text)}
		default:
			// a resource described in place: its title or value
			for _, gc := range c.Children {
				if gc.Is(nsDC, "title") || gc.Is(nsRDF, "value") {
					return values(gc)
				}
			}
		}
	}
	if s := collapse(p.Text); s != "" {
		return []string{s}
	}
	return nil
}

// collapse trims text and collapses its runs of white space.
func collapse(s string) string { return strings.Join(strings.Fields(s), " ") }

// dates normalizes dates, dropping those that are not.
func dates(values []string) []string {
	var out []string
	for _, v := range values {
		if d := Date(v); d != "" {
			out = append(out, d)
		}
	}
	return out
}

// first returns the first value only.
func first(values []string) bdf.DCValues {
	if len(values) == 0 {
		return nil
	}
	return bdf.DCValues{values[0]}
}

// add appends the values that are not empty (trimmed of spaces and NULs).
func add(f *bdf.DCValues, values ...string) {
	for _, v := range values {
		if v = strings.TrimSpace(strings.Trim(v, "\x00")); v != "" {
			*f = append(*f, v)
		}
	}
}

// Merge fills the elements of dc that have no value with those of src.
// Sources are merged in order of precedence.
func Merge(dc *bdf.DublinCore, src bdf.DublinCore) {
	for _, name := range bdf.DCTerms {
		if f := dc.Field(name); len(*f) == 0 {
			*f = *src.Field(name)
		}
	}
}

// SplitList splits a list whose items are separated by semicolons (EXIF's
// and TIFF's Artist, the Windows author and keyword tags).
func SplitList(s string) []string {
	var out []string
	for _, v := range strings.Split(s, ";") {
		if v = strings.TrimSpace(strings.Trim(v, "\x00")); v != "" {
			out = append(out, v)
		}
	}
	return out
}

// Placeholder reports whether an image description is the text a camera
// writes when the user wrote none.
func Placeholder(s string) bool {
	switch strings.ToUpper(strings.Join(strings.Fields(s), " ")) {
	case "", "OLYMPUS DIGITAL CAMERA", "SONY DSC", "DIGITAL CAMERA", "SAMSUNG DIGITAL CAMERA",
		"MINOLTA DIGITAL CAMERA", "KONICA MINOLTA DIGITAL CAMERA", "KODAK DIGITAL STILL CAMERA",
		"EXIF_JPEG_PICTURE", "DEFAULT":
		return true
	}
	return false
}

// dateRE matches the dates of XMP (ISO 8601: 2026, 2026-09, 2026-09-01,
// 2026-09-01T10:20, …:30, …:30.25, with a zone), and EXIF-style dates that
// some tools write there (2026:09:01 10:20:30).
var dateRE = regexp.MustCompile(`^(\d{4})(?:[-:](\d{2})(?:[-:](\d{2})(?:[T ](\d{2}):(\d{2})(?::(\d{2})(\.\d+)?)?\s*(Z|[+-]\d{2}:?\d{2})?)?)?)?$`)

// Date normalizes a date written in ISO 8601 (or as EXIF writes dates) to
// the W3C date and time format, keeping its precision; "" when it is not a
// date.
func Date(s string) string {
	m := dateRE.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return ""
	}
	var f [6]int
	n := 0
	for i := 1; i <= 6 && m[i] != ""; i++ {
		f[i-1], _ = strconv.Atoi(m[i])
		n = i
	}
	lo := [6]int{1, 1, 1, 0, 0, 0}
	hi := [6]int{9999, 12, 31, 23, 59, 60}
	for i := range n {
		if f[i] < lo[i] || f[i] > hi[i] {
			return ""
		}
	}
	var b strings.Builder
	b.WriteString(m[1])
	for i, sep := range []string{"-", "-", "T", ":", ":"} {
		if i+2 > n {
			break
		}
		b.WriteString(sep + m[i+2])
	}
	if n == 6 {
		b.WriteString(m[7])
	}
	if n >= 4 {
		zone := m[8]
		if len(zone) == 5 { // +0900
			zone = zone[:3] + ":" + zone[3:]
		}
		b.WriteString(zone)
	}
	return b.String()
}
