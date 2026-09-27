package image

import (
	"bytes"
	"encoding/xml"
	"errors"
	"regexp"
	"strings"
	"unicode/utf16"

	"github.com/shibukawa/bdf/converter/internal/xmp"
	"github.com/shibukawa/bdf/imgconv"
)

// utf8Text converts a document with a UTF-16 byte order mark to UTF-8 and
// drops a UTF-8 byte order mark.
func utf8Text(b []byte) []byte {
	bigEndian := bytes.HasPrefix(b, []byte{0xfe, 0xff})
	if !bigEndian && !bytes.HasPrefix(b, []byte{0xff, 0xfe}) {
		return bytes.TrimPrefix(b, []byte("\xef\xbb\xbf"))
	}
	u := make([]uint16, 0, len(b)/2)
	for i := 2; i+1 < len(b); i += 2 {
		if bigEndian {
			u = append(u, be.Uint16(b[i:]))
		} else {
			u = append(u, le.Uint16(b[i:]))
		}
	}
	return []byte(string(utf16.Decode(u)))
}

// svgRoot reports whether a document's first element (after the XML
// declaration, comments, processing instructions and document type) is an
// svg element. known is false when b ends before it.
func svgRoot(b []byte) (isSVG, known bool) {
	b = utf8Text(b)
	for {
		b = bytes.TrimLeft(b, " \t\r\n")
		var end []byte
		switch {
		case len(b) == 0:
			return false, false
		case b[0] != '<':
			return false, true
		case bytes.HasPrefix(b, []byte("<?")):
			end = []byte("?>")
		case bytes.HasPrefix(b, []byte("<!--")):
			end = []byte("-->")
		case bytes.HasPrefix(b, []byte("<!DOCTYPE")), bytes.HasPrefix(b, []byte("<!doctype")):
			// the internal subset holds declarations that end with '>'
			i := bytes.IndexAny(b, "[>")
			if i >= 0 && b[i] == '[' {
				j := bytes.Index(b[i:], []byte("]"))
				if j < 0 {
					return false, false
				}
				i += j
			}
			j := -1
			if i >= 0 {
				j = bytes.IndexByte(b[i:], '>')
			}
			if j < 0 {
				return false, false
			}
			b = b[i+j+1:]
			continue
		default:
			name := b[1:]
			if i := bytes.IndexAny(name, " \t\r\n/>"); i >= 0 {
				name = name[:i]
			} else {
				return false, false
			}
			if i := bytes.IndexByte(name, ':'); i >= 0 {
				name = name[i+1:]
			}
			return string(name) == "svg", true
		}
		i := bytes.Index(b, end)
		if i < 0 {
			return false, false
		}
		b = b[i+len(end):]
	}
}

// entityRE matches the general entities of an internal DTD subset.
var entityRE = regexp.MustCompile(`<!ENTITY\s+([A-Za-z_][\w.-]*)\s+(?:"([^"]*)"|'([^']*)')\s*>`)

// readSVG reads the size of an SVG image as browsers give it to an image
// element, its title, description and RDF metadata, and its language.
func readSVG(data []byte) (*picture, error) {
	data = utf8Text(data)
	entities := map[string]string{}
	head := data
	if i := bytes.Index(head, []byte("<svg")); i >= 0 {
		head = head[:i]
	}
	for _, m := range entityRE.FindAllSubmatch(head, -1) {
		entities[string(m[1])] = string(m[2]) + string(m[3])
	}
	d := xmp.Decoder(data, entities)
	var root xml.StartElement
	for {
		tok, err := d.Token()
		if err != nil {
			return nil, errors.New("the SVG document has no root element")
		}
		if s, ok := tok.(xml.StartElement); ok {
			root = s
			break
		}
	}
	if root.Name.Local != "svg" {
		return nil, errors.New("the root element is not svg")
	}
	p := &picture{format: "svg"}
	rn := &xmp.Node{Name: root.Name, Attrs: root.Attr}
	if root.Name.Space != nsSVG {
		p.warnings = append(p.warnings, `the root svg element is not in the SVG namespace (xmlns="`+nsSVG+`"): browsers do not draw it`)
	}
	p.w, p.h = imgconv.ParseSVGSize(rn.Attr("", "width"), rn.Attr("", "height"), rn.Attr("", "viewBox")).Pixels()
	if lang := rn.Attr(nsXML, "lang"); lang != "" {
		add(&p.native.Language, lang)
	} else {
		add(&p.native.Language, rn.Attr("", "lang"))
	}
	// title, desc and metadata among the root's children
	var meta *xmp.Node
	for {
		tok, err := d.Token()
		if err != nil {
			break
		}
		s, ok := tok.(xml.StartElement)
		if !ok {
			if _, end := tok.(xml.EndElement); end {
				break // the root's end
			}
			continue
		}
		if s.Name.Space != nsSVG && s.Name.Space != "" {
			d.Skip()
			continue
		}
		switch s.Name.Local {
		case "title", "desc":
			n, _ := xmp.Subtree(d, s)
			f := &p.native.Title
			if s.Name.Local == "desc" {
				f = &p.native.Description
			}
			p.setText(f, collapse(allText(n)))
		case "metadata":
			if meta == nil {
				meta, _ = xmp.Subtree(d, s)
			} else {
				d.Skip()
			}
		default:
			d.Skip()
		}
	}
	if meta != nil {
		// RDF metadata (Inkscape's document properties) comes before the
		// title and description, as XMP does in the raster formats
		rdf := xmp.RDF(meta)
		xmp.Merge(&rdf, p.native)
		p.native = rdf
	}
	return p, nil
}

// allText returns the text of an element and its descendants.
func allText(n *xmp.Node) string {
	s := n.Text
	for _, c := range n.Children {
		s += " " + allText(c)
	}
	return s
}

// collapse trims text and collapses its runs of white space.
func collapse(s string) string { return strings.Join(strings.Fields(s), " ") }

// Namespaces of the elements and attributes of an SVG root read.
const (
	nsSVG = "http://www.w3.org/2000/svg"
	nsXML = "http://www.w3.org/XML/1998/namespace"
)
