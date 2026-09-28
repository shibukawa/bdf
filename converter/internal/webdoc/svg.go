package webdoc

import (
	"bytes"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// urlRefRE finds references to elements by id in attribute values
// (fill="url(#g)", style="clip-path: url(#c)").
var urlRefRE = regexp.MustCompile(`url\(\s*['"]?#([^'")\s]+)['"]?\s*\)`)

// SVGDocument writes an svg element of an HTML document as an SVG file:
// with the namespaces a file declares, the text color it inherits (for
// currentColor), the SVG elements it refers to by id from elsewhere in the
// page (the symbols of a sprite sheet, shared gradients) copied into it,
// and its images as data: URLs (embed; images an SVG file refers to are not
// loaded where it is drawn as an image). Scripts are left out.
func SVGDocument(root *html.Node, color string, byID func(id string) *html.Node, embed func(href string) string) []byte {
	data, _ := SVGDocumentMax(root, color, byID, embed, 0)
	return data
}

// SVGDocumentMax is SVGDocument for a document of at most max bytes (0: of
// any size); it returns false, and no document, for an svg element that
// makes a larger one. The document of an svg element can be much larger
// than the element: it holds the elements it borrows and its images.
func SVGDocumentMax(root *html.Node, color string, byID func(id string) *html.Node, embed func(href string) string, max int) ([]byte, bool) {
	// the ids defined inside, and the elements referred to from outside
	defined := map[string]bool{}
	WalkElements(root, func(n *html.Node) {
		if id := attrStr(n, "id"); id != "" {
			defined[id] = true
		}
	})
	var borrowed []*html.Node
	pending := []*html.Node{root}
	searched := map[*html.Node]bool{}
	for len(pending) > 0 && len(borrowed) < 1000 {
		n := pending[0]
		pending = pending[1:]
		if inside(n, searched) {
			continue // its references were found in the element it is in
		}
		searched[n] = true
		WalkElements(n, func(e *html.Node) {
			for _, a := range e.Attr {
				var ids []string
				if a.Key == "href" {
					if id, ok := strings.CutPrefix(strings.TrimSpace(a.Val), "#"); ok {
						ids = append(ids, id)
					}
				}
				for _, m := range urlRefRE.FindAllStringSubmatch(a.Val, -1) {
					ids = append(ids, m[1])
				}
				for _, id := range ids {
					if defined[id] {
						continue
					}
					defined[id] = true // looked for once
					if el := byID(id); el != nil && el.Namespace == "svg" {
						borrowed = append(borrowed, el)
						pending = append(pending, el)
					}
				}
			}
		})
	}

	// the prefixes of attributes that the document declares (attributes of
	// other prefixes, such as epub:type on an svg element, or of namespaces
	// without a prefix here would make it malformed, and are left out)
	prefixes := map[string]bool{"xml": true, "xlink": true}
	WalkElements(root, func(n *html.Node) {
		for _, a := range n.Attr {
			if p, ok := strings.CutPrefix(a.Key, "xmlns:"); ok && a.Namespace == "" {
				prefixes[p] = true
			} else if a.Namespace == "xmlns" {
				prefixes[a.Key] = true
			}
		}
	})

	var b bytes.Buffer
	var write func(n *html.Node, isRoot, inForeign bool)
	write = func(n *html.Node, isRoot, inForeign bool) {
		if max > 0 && b.Len() > max {
			return
		}
		switch n.Type {
		case html.TextNode:
			escapeXML(&b, n.Data, false)
			return
		case html.ElementNode:
		default:
			return
		}
		if n.Data == "script" || strings.Contains(n.Data, ":") || n.Namespace != "" && n.Namespace != "svg" && n.Namespace != "math" {
			return // scripts, and elements of other vocabularies (an editor's metadata)
		}
		b.WriteString("<" + n.Data)
		hasNS, hasXlink, hasColor := false, false, false
		for _, a := range n.Attr {
			name := a.Key
			if p, _, ok := strings.Cut(a.Key, ":"); ok && a.Namespace == "" && p != "xmlns" && !prefixes[p] {
				continue
			}
			switch {
			case a.Namespace == "xlink" || a.Namespace == "xml" || a.Namespace == "xmlns":
				name = a.Namespace + ":" + a.Key
				hasXlink = hasXlink || a.Namespace == "xmlns" && a.Key == "xlink"
			case a.Namespace != "":
				continue
			case a.Key == "xmlns":
				hasNS = true
			case a.Key == "color":
				hasColor = true
			case strings.HasPrefix(a.Key, "on"):
				continue // event handlers
			}
			val := a.Val
			if (a.Key == "href" || name == "xlink:href") && (n.Data == "image" || n.Data == "feImage") {
				if u := embed(val); u != "" {
					val = u
				}
			}
			b.WriteString(" " + name + `="`)
			escapeXML(&b, val, true)
			b.WriteByte('"')
		}
		switch {
		case isRoot:
			if !hasNS {
				b.WriteString(` xmlns="` + nsSVG + `"`)
			}
			if !hasXlink {
				b.WriteString(` xmlns:xlink="` + nsXLink + `"`)
			}
			if !hasColor && color != "" {
				b.WriteString(` color="` + color + `"`)
			}
		case n.Namespace == "" && !inForeign && !hasNS:
			// HTML in a foreignObject
			b.WriteString(` xmlns="` + nsXHTML + `"`)
			inForeign = true
		case n.Namespace == "math" && n.Data == "math" && !hasNS:
			b.WriteString(` xmlns="` + nsMathML + `"`)
		}
		b.WriteByte('>')
		if isRoot && len(borrowed) > 0 {
			b.WriteString("<defs>")
			for _, e := range borrowed {
				write(e, false, false)
			}
			b.WriteString("</defs>")
		}
		for k := n.FirstChild; k != nil; k = k.NextSibling {
			write(k, false, inForeign)
		}
		b.WriteString("</" + n.Data + ">")
	}
	write(root, true, false)
	if max > 0 && b.Len() > max {
		return nil, false
	}
	return b.Bytes(), true
}

// inside reports whether an element is in one of the elements of a set.
func inside(n *html.Node, set map[*html.Node]bool) bool {
	for p := n.Parent; p != nil; p = p.Parent {
		if set[p] {
			return true
		}
	}
	return false
}

// WalkElements calls fn for n and every element under it.
func WalkElements(n *html.Node, fn func(*html.Node)) {
	if n.Type == html.ElementNode {
		fn(n)
	}
	for k := n.FirstChild; k != nil; k = k.NextSibling {
		WalkElements(k, fn)
	}
}

// escapeXML writes text escaped for XML content, or for a double-quoted
// attribute value.
func escapeXML(b *bytes.Buffer, s string, attr bool) {
	for _, c := range s {
		switch {
		case c == '&':
			b.WriteString("&amp;")
		case c == '<':
			b.WriteString("&lt;")
		case c == '>':
			b.WriteString("&gt;")
		case c == '"' && attr:
			b.WriteString("&quot;")
		case c == '\n' && attr:
			b.WriteString("&#10;")
		case c < 0x20 && c != '\t' && c != '\n' && c != '\r':
			// not allowed in XML 1.0
		default:
			b.WriteRune(c)
		}
	}
}

// attrStr returns the value of an attribute without a namespace, trimmed.
func attrStr(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if a.Namespace == "" && a.Key == name {
			return strings.TrimSpace(a.Val)
		}
	}
	return ""
}
