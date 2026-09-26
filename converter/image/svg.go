package image

import (
	"bytes"
	"encoding/xml"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"
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
	d := decoder(data, entities)
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
	rn := &node{name: root.Name, attrs: root.Attr}
	if root.Name.Space != nsSVG {
		p.warnings = append(p.warnings, `the root svg element is not in the SVG namespace (xmlns="`+nsSVG+`"): browsers do not draw it`)
	}
	p.w, p.h = svgSize(rn.attr("", "width"), rn.attr("", "height"), rn.attr("", "viewBox"))
	if lang := rn.attr(nsXML, "lang"); lang != "" {
		add(&p.native.Language, lang)
	} else {
		add(&p.native.Language, rn.attr("", "lang"))
	}
	// title, desc and metadata among the root's children
	var meta *node
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
			n, _ := subtree(d, s)
			f := &p.native.Title
			if s.Name.Local == "desc" {
				f = &p.native.Description
			}
			p.setText(f, collapse(allText(n)))
		case "metadata":
			if meta == nil {
				meta, _ = subtree(d, s)
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
		rdf := rdfDC(meta)
		merge(&rdf, p.native)
		p.native = rdf
	}
	return p, nil
}

// allText returns the text of an element and its descendants.
func allText(n *node) string {
	s := n.text
	for _, c := range n.children {
		s += " " + allText(c)
	}
	return s
}

// svgSize computes the size of an SVG image from the width, height and
// viewBox of its root, as for an image element (CSS px): a missing or
// relative width or height follows from the other and the view box's
// proportions, or is the view box's; without a view box it is 300 × 150
// (the default size of replaced elements).
func svgSize(width, height, viewBox string) (w, h float64) {
	w, okW := svgLength(width)
	h, okH := svgLength(height)
	var vw, vh float64
	if f := strings.FieldsFunc(viewBox, func(r rune) bool { return r == ' ' || r == ',' || r == '\t' || r == '\n' || r == '\r' }); len(f) == 4 {
		vw, _ = strconv.ParseFloat(f[2], 64)
		vh, _ = strconv.ParseFloat(f[3], 64)
	}
	ratio := vw > 0 && vh > 0
	switch {
	case okW && okH:
	case okW && ratio:
		h = w * vh / vw
	case okH && ratio:
		w = h * vw / vh
	case ratio:
		w, h = vw, vh
	default:
		if !okW {
			w = 300
		}
		if !okH {
			h = 150
		}
	}
	return w, h
}

// lengthRE matches a CSS length: a number and a unit.
var lengthRE = regexp.MustCompile(`^\s*([+-]?(?:\d+\.?\d*|\.\d+)(?:[eE][+-]?\d+)?)\s*([A-Za-z%]*)\s*$`)

// svgLength reads an absolute length in CSS px; percentages and invalid
// lengths are not.
func svgLength(s string) (float64, bool) {
	m := lengthRE.FindStringSubmatch(s)
	if m == nil {
		return 0, false
	}
	v, err := strconv.ParseFloat(m[1], 64)
	if err != nil || v <= 0 {
		return 0, false
	}
	switch strings.ToLower(m[2]) {
	case "", "px":
	case "pt":
		v *= 96.0 / 72
	case "pc":
		v *= 16
	case "in":
		v *= 96
	case "cm":
		v *= 96 / 2.54
	case "mm":
		v *= 96 / 25.4
	case "q":
		v *= 96 / 101.6
	case "em":
		v *= 16
	case "ex":
		v *= 8
	default:
		return 0, false
	}
	return v, true
}
