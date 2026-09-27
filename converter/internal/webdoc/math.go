package webdoc

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Pages seldom hold formulas as bare MathML: KaTeX, MathJax and Wikipedia
// write the MathML beside a rendering of their own in HTML, SVG or an
// image (which a reader view would show a second time, and garbled), and
// hide one of them with their style sheets, which reader mode does not
// read. NormalizeMath rewrites these into the math element alone (the HTML
// converter does so before the article is picked out, which drops the
// class names they are told apart by). LaTeX sources (MathJax 2's scripts, GitHub's math-renderer
// elements) become math elements that hold the source as a TeX
// annotation, which the formula reader parses.
func NormalizeMath(n *html.Node) {
	for k := n.FirstChild; k != nil; {
		next := k.NextSibling
		if k.Type == html.ElementNode {
			normalizeMathElement(k)
		}
		k = next
	}
}

func normalizeMathElement(n *html.Node) {
	cls := " " + attrStr(n, "class") + " "
	switch {
	case strings.Contains(cls, " mwe-math-element "), n.Data == "mjx-container":
		// Wikipedia, MathJax 3: the MathML (hidden or assistive) beside an
		// image or glyphs drawn with CSS
		if m := findMath(n); m != nil {
			if n.Data == "mjx-container" && attrStr(n, "display") == "true" {
				setAttr(m, "display", "block")
			}
			m.Parent.RemoveChild(m)
			replaceChildren(n, m)
			n.Attr = nil
			return
		}
	case strings.Contains(cls, " katex-html "), strings.Contains(cls, " MathJax_Preview "),
		strings.Contains(cls, " MathJax_Display "), strings.Contains(cls, " MathJax "),
		strings.Contains(cls, " MathJax_SVG "), strings.Contains(cls, " MathJax_SVG_Display "),
		strings.Contains(cls, " MathJax_CHTML "):
		// a rendering beside the MathML or TeX source
		if hasMathSource(n) {
			n.Parent.RemoveChild(n)
			return
		}
	case n.DataAtom == atom.Script:
		t := strings.ToLower(attrStr(n, "type"))
		if strings.HasPrefix(t, "math/tex") {
			m := texMath(textContent(n), strings.Contains(t, "mode=display"))
			n.Parent.InsertBefore(m, n)
			n.Parent.RemoveChild(n)
		}
		return
	case n.Data == "math-renderer":
		// GitHub: the LaTeX source between dollar signs, typeset by a script
		src := strings.TrimSpace(textContent(n))
		display := strings.HasPrefix(src, "$$") || strings.Contains(cls, " js-display-math ")
		if strings.HasPrefix(src, "$`") && strings.HasSuffix(src, "`$") && len(src) >= 4 {
			src = src[2 : len(src)-2]
		} else {
			src = strings.Trim(src, "$")
		}
		replaceChildren(n, texMath(src, display))
		n.Attr = nil
		return
	case n.DataAtom == atom.Math:
		return
	}
	NormalizeMath(n)
}

// hasMathSource reports whether the element's parent holds the formula as
// MathML or TeX besides the element.
func hasMathSource(n *html.Node) bool {
	p := n.Parent
	if p == nil {
		return false
	}
	for k := p.FirstChild; k != nil; k = k.NextSibling {
		if k == n || k.Type != html.ElementNode {
			continue
		}
		if findMath(k) != nil || k.DataAtom == atom.Script && strings.HasPrefix(strings.ToLower(attrStr(k, "type")), "math/tex") {
			return true
		}
	}
	return false
}

// findMath returns the first math element in n (n itself included).
func findMath(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.DataAtom == atom.Math {
		return n
	}
	for k := n.FirstChild; k != nil; k = k.NextSibling {
		if m := findMath(k); m != nil {
			return m
		}
	}
	return nil
}

// texMath makes a math element holding LaTeX as its annotation.
func texMath(src string, display bool) *html.Node {
	m := &html.Node{Type: html.ElementNode, Data: "math", DataAtom: atom.Math, Namespace: "math"}
	if display {
		m.Attr = []html.Attribute{{Key: "display", Val: "block"}}
	}
	sem := &html.Node{Type: html.ElementNode, Data: "semantics", Namespace: "math"}
	ann := &html.Node{Type: html.ElementNode, Data: "annotation", Namespace: "math",
		Attr: []html.Attribute{{Key: "encoding", Val: "application/x-tex"}}}
	ann.AppendChild(&html.Node{Type: html.TextNode, Data: src})
	sem.AppendChild(ann)
	m.AppendChild(sem)
	return m
}

func replaceChildren(n, only *html.Node) {
	for k := n.FirstChild; k != nil; {
		next := k.NextSibling
		n.RemoveChild(k)
		k = next
	}
	n.AppendChild(only)
}

func setAttr(n *html.Node, key, val string) {
	for i, a := range n.Attr {
		if a.Key == key {
			n.Attr[i].Val = val
			return
		}
	}
	n.Attr = append(n.Attr, html.Attribute{Key: key, Val: val})
}

func textContent(n *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		for k := n.FirstChild; k != nil; k = k.NextSibling {
			if k.Type == html.TextNode {
				b.WriteString(k.Data)
			} else {
				walk(k)
			}
		}
	}
	walk(n)
	return b.String()
}
