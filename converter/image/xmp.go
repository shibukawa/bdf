package image

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
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
	nsSVG       = "http://www.w3.org/2000/svg"
)

// node is an XML element.
type node struct {
	name     xml.Name
	attrs    []xml.Attr
	children []*node
	text     string
}

// attr returns an attribute's value ("" when absent).
func (n *node) attr(space, local string) string {
	for _, a := range n.attrs {
		if a.Name.Local == local && a.Name.Space == space {
			return a.Value
		}
	}
	return ""
}

// is reports whether the element has the name.
func (n *node) is(space, local string) bool { return n.name.Space == space && n.name.Local == local }

// decoder makes an XML decoder that tolerates what browsers tolerate in
// SVG files: the entities of the internal DTD subset (Illustrator declares
// its namespaces that way), and charsets other than UTF-8 (Latin-1 is
// decoded; others keep their ASCII and replace the rest).
func decoder(data []byte, entities map[string]string) *xml.Decoder {
	d := xml.NewDecoder(bytes.NewReader(data))
	d.Strict = false
	d.Entity = entities
	d.CharsetReader = func(charset string, r io.Reader) (io.Reader, error) {
		b, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		switch strings.ToLower(charset) {
		case "iso-8859-1", "latin1", "latin-1", "iso_8859-1", "l1", "windows-1252", "cp1252", "us-ascii", "ascii":
			return strings.NewReader(latin1(b)), nil
		}
		return bytes.NewReader(bytes.ToValidUTF8(b, []byte(string(utf8.RuneError)))), nil
	}
	return d
}

// subtree reads the element that start opens, up to its end.
func subtree(d *xml.Decoder, start xml.StartElement) (*node, error) {
	n := &node{name: start.Name, attrs: start.Attr}
	var text strings.Builder
	for {
		tok, err := d.Token()
		if err != nil {
			return n, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			c, err := subtree(d, t)
			n.children = append(n.children, c)
			if err != nil {
				return n, err
			}
		case xml.CharData:
			text.Write(t)
		case xml.EndElement:
			n.text = text.String()
			return n, nil
		}
	}
}

// parseXML reads a document's root element.
func parseXML(data []byte) (*node, error) {
	d := decoder(data, nil)
	for {
		tok, err := d.Token()
		if err != nil {
			return nil, err
		}
		if s, ok := tok.(xml.StartElement); ok {
			return subtree(d, s)
		}
	}
}

// rdfValues returns the values of an RDF property element: the items of an
// rdf:Alt (the default language first), rdf:Seq or rdf:Bag, the resource it
// names, the title of the resource it describes (Inkscape's cc:Agent), or
// its text.
func rdfValues(p *node) []string {
	if r := p.attr(nsRDF, "resource"); r != "" {
		return []string{r}
	}
	if v := p.attr(nsRDF, "value"); v != "" {
		return []string{v}
	}
	for _, c := range p.children {
		switch {
		case c.is(nsRDF, "Alt"), c.is(nsRDF, "Seq"), c.is(nsRDF, "Bag"):
			var out []string
			for _, li := range c.children {
				if !li.is(nsRDF, "li") {
					continue
				}
				v := rdfValues(li)
				if c.is(nsRDF, "Alt") && li.attr(nsXML, "lang") == "x-default" {
					out = append(v, out...)
				} else {
					out = append(out, v...)
				}
			}
			return out
		case c.is(nsRDF, "value"):
			return []string{collapse(c.text)}
		default:
			// a resource described in place: its title or value
			for _, gc := range c.children {
				if gc.is(nsDC, "title") || gc.is(nsRDF, "value") {
					return rdfValues(gc)
				}
			}
		}
	}
	if s := collapse(p.text); s != "" {
		return []string{s}
	}
	return nil
}

// collapse trims text and collapses its runs of white space.
func collapse(s string) string { return strings.Join(strings.Fields(s), " ") }

// rdfDC reads the Dublin Core of the RDF/XML under an element (an XMP
// packet, SVG metadata). The dates of XMP's other schemas become the
// created and modified dates, and the TIFF schema (EXIF mirrored in XMP)
// fills in the elements the dc schema leaves out.
func rdfDC(root *node) bdf.DublinCore {
	var dc, tiffDC bdf.DublinCore
	var created, digitized []string
	prop := func(space, local string, values []string) {
		switch space {
		case nsDC:
			if f := dc.Field(local); f != nil && local != "format" {
				if local == "date" {
					values = mapDates(values)
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
				add(&dc.Created, mapDates(values)...)
			case "modified":
				add(&dc.Modified, mapDates(values)...)
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
				add(&dc.Modified, mapDates(values)...)
			}
		case nsTIFF:
			switch local {
			case "ImageDescription":
				add(&tiffDC.Description, first(values)...)
			case "Artist":
				for _, v := range values {
					add(&tiffDC.Creator, splitList(v)...)
				}
			case "Copyright":
				add(&tiffDC.Rights, first(values)...)
			case "DateTime":
				add(&tiffDC.Modified, mapDates(values)...)
			}
		}
	}
	var walk func(n *node)
	walk = func(n *node) {
		if !n.is(nsRDF, "RDF") {
			for _, c := range n.children {
				walk(c)
			}
			return
		}
		// node elements (rdf:Description, cc:Work …): their properties are
		// attributes and child elements
		for _, desc := range n.children {
			for _, a := range desc.attrs {
				if a.Name.Space != nsRDF && a.Name.Space != "xmlns" && a.Name.Space != nsXML {
					prop(a.Name.Space, a.Name.Local, []string{a.Value})
				}
			}
			for _, p := range desc.children {
				prop(p.name.Space, p.name.Local, rdfValues(p))
			}
		}
	}
	walk(root)
	if len(dc.Created) == 0 {
		if dc.Created = first(mapDates(created)); len(dc.Created) == 0 {
			dc.Created = first(mapDates(digitized))
		}
	}
	if len(dc.Description) > 0 && placeholder(dc.Description[0]) {
		dc.Description = nil
	}
	if len(tiffDC.Description) > 0 && placeholder(tiffDC.Description[0]) {
		tiffDC.Description = nil
	}
	merge(&dc, tiffDC)
	return dc
}

// mapDates normalizes dates, dropping those that are not.
func mapDates(values []string) []string {
	var out []string
	for _, v := range values {
		if d := isoDate(v); d != "" {
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

// parseXMP reads an XMP packet; what was read before an error is kept.
func parseXMP(b []byte) bdf.DublinCore {
	root, _ := parseXML(b)
	if root == nil {
		return bdf.DublinCore{}
	}
	return rdfDC(root)
}
