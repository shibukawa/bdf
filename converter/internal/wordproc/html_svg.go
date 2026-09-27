package wordproc

import (
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/converter/internal/webdoc"
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
	data := webdoc.SVGDocument(n, color.CSS(), r.elementByID(n), r.embedImage)
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

// elementByID returns a lookup of the elements of the document n is in by
// their id, made once per document.
func (r *htmlReader) elementByID(n *html.Node) func(id string) *html.Node {
	top := n
	for top.Parent != nil {
		top = top.Parent
	}
	if r.ids == nil || r.idsOf != top {
		r.ids, r.idsOf = map[string]*html.Node{}, top
		webdoc.WalkElements(top, func(e *html.Node) {
			if id := attrStr(e, "id"); id != "" {
				if _, dup := r.ids[id]; !dup {
					r.ids[id] = e
				}
			}
		})
	}
	return func(id string) *html.Node { return r.ids[id] }
}
