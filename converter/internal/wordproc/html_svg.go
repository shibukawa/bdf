package wordproc

import (
	"bytes"
	"encoding/base64"
	"regexp"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/imgconv"
	"golang.org/x/net/html"
)

// inlineSVG draws an svg element of the document as an image: its markup,
// made an SVG document of its own, which the viewer draws (spec §6.2). It
// is sized as browsers size inline SVG: the width and height its
// attributes and style give, a missing one following from the view box's
// proportions; with only a view box it is as wide as the line.
func (r *htmlReader) inlineSVG(n *html.Node, css map[string]string, st *hstyle) {
	size := func(name string) (float64, bool) {
		v := css[name]
		if v == "" {
			v = attrStr(n, name)
		}
		return r.length(v), zeroLength(v)
	}
	aw, zw := size("width")
	ah, zh := size("height")
	if zw || zh || !drawsSomething(n) {
		return // an empty box, or a sprite sheet of symbols for other svg elements to use
	}
	// currentColor is the color of the text around it
	color := st.rp.color
	if color == 0 {
		color = textColor
	}
	data := svgDocument(n, color.CSS(), r.elementByID(n), r.embedImage)
	s := imgconv.ParseSVGSize("", "", attrStr(n, "viewBox"))
	im := &htmlImage{hash: r.c.doc.AddImage(data), ratio: s.Ratio()}
	r.picture(n, css, st, im, aw, ah, svgAlt(n, st))
}

// zeroLength reports a length that is written as zero ("0", "0px").
func zeroLength(v string) bool {
	v = strings.TrimSpace(v)
	num := strings.TrimRight(v, "abcdefghijklmnopqrstuvwxyz%")
	f, err := strconv.ParseFloat(num, 64)
	return v != "" && err == nil && f == 0
}

// drawsSomething reports whether an svg element has content besides
// definitions and descriptions.
func drawsSomething(n *html.Node) bool {
	for k := n.FirstChild; k != nil; k = k.NextSibling {
		if k.Type != html.ElementNode {
			continue
		}
		switch k.Data {
		case "defs", "symbol", "style", "title", "desc", "metadata", "script", "linearGradient", "radialGradient", "pattern",
			"clipPath", "mask", "filter", "marker":
		default:
			return true
		}
	}
	return false
}

// svgAlt is the alternative text of an svg element: aria-label, else its
// title, else the caption of its figure; "" for a decorative one.
func svgAlt(n *html.Node, st *hstyle) string {
	if attrStr(n, "aria-hidden") == "true" {
		return ""
	}
	switch attrStr(n, "role") {
	case "presentation", "none":
		return ""
	}
	if a := strings.Join(strings.Fields(attrStr(n, "aria-label")), " "); a != "" {
		return a
	}
	for k := n.FirstChild; k != nil; k = k.NextSibling {
		if k.Type == html.ElementNode && k.Data == "title" {
			if t := strings.Join(strings.Fields(textContent(k)), " "); t != "" {
				return t
			}
		}
	}
	return st.figAlt
}

// embedImage returns the image an image element of an inline SVG refers
// to as a data: URL, which an SVG document drawn as an image can show
// ("" when it cannot be read).
func (r *htmlReader) embedImage(href string) string {
	if href == "" || strings.HasPrefix(href, "data:") || strings.HasPrefix(href, "#") || r.d.Image == nil {
		return ""
	}
	data, err := r.d.Image(href)
	if err != nil {
		r.c.warnf("image %s in an SVG: %v", href, err)
		return ""
	}
	format := imgconv.Sniff(data)
	if format == "" {
		return ""
	}
	return "data:" + imgconv.MIME(format) + ";base64," + base64.StdEncoding.EncodeToString(data)
}

// urlRefRE finds references to elements by id in attribute values
// (fill="url(#g)", style="clip-path: url(#c)").
var urlRefRE = regexp.MustCompile(`url\(\s*['"]?#([^'")\s]+)['"]?\s*\)`)

// elementByID returns a lookup of the elements of the document n is in by
// their id, made once per document.
func (r *htmlReader) elementByID(n *html.Node) func(id string) *html.Node {
	top := n
	for top.Parent != nil {
		top = top.Parent
	}
	if r.ids == nil || r.idsOf != top {
		r.ids, r.idsOf = map[string]*html.Node{}, top
		walkElements(top, func(e *html.Node) {
			if id := attrStr(e, "id"); id != "" {
				if _, dup := r.ids[id]; !dup {
					r.ids[id] = e
				}
			}
		})
	}
	return func(id string) *html.Node { return r.ids[id] }
}

const (
	svgNS   = "http://www.w3.org/2000/svg"
	xlinkNS = "http://www.w3.org/1999/xlink"
	xhtmlNS = "http://www.w3.org/1999/xhtml"
)

// svgDocument writes an svg element of an HTML document as an SVG file:
// with the namespaces a file declares, the text color it inherits (for
// currentColor), the SVG elements it refers to by id from elsewhere in the
// page (the symbols of a sprite sheet, shared gradients) copied into it,
// and its images as data: URLs (embed; images an SVG file refers to are not
// loaded where it is drawn as an image). Scripts are left out.
func svgDocument(root *html.Node, color string, byID func(id string) *html.Node, embed func(href string) string) []byte {
	// the ids defined inside, and the elements referred to from outside
	defined := map[string]bool{}
	walkElements(root, func(n *html.Node) {
		if id := attrStr(n, "id"); id != "" {
			defined[id] = true
		}
	})
	var borrowed []*html.Node
	pending := []*html.Node{root}
	for len(pending) > 0 && len(borrowed) < 1000 {
		n := pending[0]
		pending = pending[1:]
		walkElements(n, func(e *html.Node) {
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

	var b bytes.Buffer
	var write func(n *html.Node, isRoot, inForeign bool)
	write = func(n *html.Node, isRoot, inForeign bool) {
		switch n.Type {
		case html.TextNode:
			escapeXML(&b, n.Data, false)
			return
		case html.ElementNode:
		default:
			return
		}
		if n.Data == "script" {
			return
		}
		b.WriteString("<" + n.Data)
		hasNS, hasXlink, hasColor := false, false, false
		for _, a := range n.Attr {
			name := a.Key
			switch {
			case a.Namespace != "":
				name = a.Namespace + ":" + a.Key
				hasXlink = hasXlink || a.Namespace == "xmlns" && a.Key == "xlink"
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
				b.WriteString(` xmlns="` + svgNS + `"`)
			}
			if !hasXlink {
				b.WriteString(` xmlns:xlink="` + xlinkNS + `"`)
			}
			if !hasColor && color != "" {
				b.WriteString(` color="` + color + `"`)
			}
		case n.Namespace == "" && !inForeign && !hasNS:
			// HTML in a foreignObject
			b.WriteString(` xmlns="` + xhtmlNS + `"`)
			inForeign = true
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
	return b.Bytes()
}

// walkElements calls fn for n and every element under it.
func walkElements(n *html.Node, fn func(*html.Node)) {
	if n.Type == html.ElementNode {
		fn(n)
	}
	for k := n.FirstChild; k != nil; k = k.NextSibling {
		walkElements(k, fn)
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
