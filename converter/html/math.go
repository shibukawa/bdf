package html

import (
	"strings"

	xhtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Pages seldom hold formulas as bare MathML: KaTeX, MathJax and Wikipedia
// write the MathML beside a rendering of their own in HTML, SVG or an
// image (which a reader view would show a second time, and garbled), and
// hide one of them with their style sheets, which reader mode does not
// read. normalizeMath rewrites these into the math element alone before
// the article is picked out (which drops the class names they are told
// apart by). LaTeX sources (MathJax 2's scripts, GitHub's math-renderer
// elements) become math elements that hold the source as a TeX
// annotation, which the formula reader parses.
func normalizeMath(n *xhtml.Node) {
	for k := n.FirstChild; k != nil; {
		next := k.NextSibling
		if k.Type == xhtml.ElementNode {
			normalizeMathElement(k)
		}
		k = next
	}
}

func normalizeMathElement(n *xhtml.Node) {
	cls := " " + attr(n, "class") + " "
	switch {
	case strings.Contains(cls, " mwe-math-element "), n.Data == "mjx-container":
		// Wikipedia, MathJax 3: the MathML (hidden or assistive) beside an
		// image or glyphs drawn with CSS
		if m := findMath(n); m != nil {
			if n.Data == "mjx-container" && attr(n, "display") == "true" {
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
		t := strings.ToLower(attr(n, "type"))
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
	normalizeMath(n)
}

// hasMathSource reports whether the element's parent holds the formula as
// MathML or TeX besides the element.
func hasMathSource(n *xhtml.Node) bool {
	p := n.Parent
	if p == nil {
		return false
	}
	for k := p.FirstChild; k != nil; k = k.NextSibling {
		if k == n || k.Type != xhtml.ElementNode {
			continue
		}
		if findMath(k) != nil || k.DataAtom == atom.Script && strings.HasPrefix(strings.ToLower(attr(k, "type")), "math/tex") {
			return true
		}
	}
	return false
}

// findMath returns the first math element in n (n itself included).
func findMath(n *xhtml.Node) *xhtml.Node {
	if n.Type == xhtml.ElementNode && n.DataAtom == atom.Math {
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
func texMath(src string, display bool) *xhtml.Node {
	m := &xhtml.Node{Type: xhtml.ElementNode, Data: "math", DataAtom: atom.Math, Namespace: "math"}
	if display {
		m.Attr = []xhtml.Attribute{{Key: "display", Val: "block"}}
	}
	sem := &xhtml.Node{Type: xhtml.ElementNode, Data: "semantics", Namespace: "math"}
	ann := &xhtml.Node{Type: xhtml.ElementNode, Data: "annotation", Namespace: "math",
		Attr: []xhtml.Attribute{{Key: "encoding", Val: "application/x-tex"}}}
	ann.AppendChild(&xhtml.Node{Type: xhtml.TextNode, Data: src})
	sem.AppendChild(ann)
	m.AppendChild(sem)
	return m
}

func replaceChildren(n, only *xhtml.Node) {
	for k := n.FirstChild; k != nil; {
		next := k.NextSibling
		n.RemoveChild(k)
		k = next
	}
	n.AppendChild(only)
}

func setAttr(n *xhtml.Node, key, val string) {
	for i, a := range n.Attr {
		if a.Key == key {
			n.Attr[i].Val = val
			return
		}
	}
	n.Attr = append(n.Attr, xhtml.Attribute{Key: key, Val: val})
}

func textContent(n *xhtml.Node) string {
	var b strings.Builder
	var walk func(*xhtml.Node)
	walk = func(n *xhtml.Node) {
		for k := n.FirstChild; k != nil; k = k.NextSibling {
			if k.Type == xhtml.TextNode {
				b.WriteString(k.Data)
			} else {
				walk(k)
			}
		}
	}
	walk(n)
	return b.String()
}
