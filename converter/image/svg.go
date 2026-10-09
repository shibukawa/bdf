package image

import (
	"bytes"
	"errors"
	"strings"

	"github.com/shibukawa/tinygodriver/encoding/xmlro"

	"github.com/shibukawa/bdf/converter/internal/xmp"
	"github.com/shibukawa/bdf/imgconv"
	"github.com/shibukawa/bdf/internal/xmltree"
)

// svgRoot reports whether a document's first element (after the XML
// declaration, comments, processing instructions and document type) is an
// svg element. known is false when b ends before it.
func svgRoot(b []byte) (isSVG, known bool) {
	b = xmltree.UTF8(b)
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

// readSVG reads the size of an SVG image as browsers give it to an image
// element, its title, description and RDF metadata, and its language.
func readSVG(data []byte) (*picture, error) {
	var warnings []string
	data, dropped := xmltree.BoundEntities(xmltree.UTF8(data))
	if dropped {
		warnings = append(warnings, "the entities of the SVG document stand for too much text: they are not replaced")
	}
	r := xmp.Reader(data)
	for {
		k, err := r.Next()
		if err != nil || k == xmlro.EOF {
			return nil, errors.New("the SVG document has no root element")
		}
		if k == xmlro.StartElement {
			break
		}
	}
	rn := &xmp.Node{Name: xmltree.ElementName(r), Attrs: xmltree.Attrs(r)}
	if rn.Name.Local != "svg" {
		return nil, errors.New("the root element is not svg")
	}
	p := &picture{format: "svg", warnings: warnings}
	if rn.Name.Space != nsSVG {
		p.warnings = append(p.warnings, `the root svg element is not in the SVG namespace (xmlns="`+nsSVG+`"): browsers do not draw it`)
	}
	p.w, p.h = imgconv.ParseSVGSize(rn.Attr("", "width"), rn.Attr("", "height"), rn.Attr("", "viewBox")).Pixels()
	if lang := rn.Attr(nsXML, "lang"); lang != "" {
		add(&p.native.Language, lang)
	} else {
		add(&p.native.Language, rn.Attr("", "lang"))
	}
	// title, desc and metadata among the root's children; the others are
	// passed over
	var meta *xmp.Node
	root := r.Element()
	for {
		ok, err := r.NextChild(root)
		if err != nil || !ok {
			break
		}
		name := xmltree.ElementName(r)
		if name.Space != nsSVG && name.Space != "" {
			continue
		}
		switch name.Local {
		case "title", "desc":
			n, _ := xmp.Subtree(r)
			f := &p.native.Title
			if name.Local == "desc" {
				f = &p.native.Description
			}
			p.setText(f, collapse(allText(n)))
		case "metadata":
			if meta == nil {
				meta, _ = xmp.Subtree(r)
			}
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
