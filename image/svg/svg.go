// Package svg reads SVG documents: the tree of their elements with the
// style sheets of the document applied, and the values their attributes
// hold — colours, numbers, lengths and transforms. It draws nothing;
// raster/imagebdf is a renderer of the documents it reads.
//
//	doc, err := svg.Parse(data, nil)
//	if err != nil {
//		return err // no svg element
//	}
//	for _, n := range doc.Root.Children {
//		if fill, ok := n.Declared("fill"); ok {
//			c, _ := svg.ParseColor(fill, svg.Color{A: 1})
//			...
//		}
//	}
//
// The reader is lenient, as browsers are with the SVG images of web pages:
// it reads the entities and the void elements of HTML, documents in the
// encodings they declare, and a document with an error up to the error.
// Only a document without an svg root element is refused.
//
// A document decides how much its reader keeps in memory, so reading has
// limits (Options): what passes one ends the reading there, and
// Document.Warnings says so. The entities a DOCTYPE declares are replaced
// while together they stand for a megabyte of text at most.
//
// What is read of a style sheet is what images made for previews use:
// rules with selectors of type, class and id, combined and inside each
// other. Rules with attribute selectors or pseudo-classes and at-rules
// are left out.
//
// Path data is not read here: the d attribute of a path is the string of
// the document.
package svg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shibukawa/tinygodriver/encoding/xmlro"
	"github.com/shibukawa/tinygodriver/encoding/xmlro/htmlentity"
	"golang.org/x/net/html/charset"

	"github.com/shibukawa/bdf/internal/xmltree"
)

// TextNode is the Name of a Node that is a run of text.
const TextNode = "#text"

// Node is an element of an SVG document, or a run of text (TextNode).
type Node struct {
	// Name is the local name of the element (rect, linearGradient),
	// whatever prefix the document gives its namespace.
	Name string
	// Attr holds the attributes by their local names. The declarations of
	// namespaces are left out, and xlink:href is href (an href attribute
	// of the element itself wins).
	Attr     map[string]string
	Children []*Node
	Parent   *Node
	// Sheet holds the declarations of the rules of the document's style
	// sheets that match the element, the more specific and the later ones
	// winning; nil when none matches. Declared looks here.
	Sheet map[string]string
	// Text is the text of a run of text. Runs are kept where SVG draws
	// them: in text, tspan, textPath and a elements.
	Text string
}

// Document is an SVG document that was read.
type Document struct {
	// Root is the svg element.
	Root *Node
	// IDs holds the elements by their id attributes, the first of the
	// elements that share one.
	IDs map[string]*Node
	// Warnings say what was left out of a document past the limits.
	Warnings []string
}

// The limits of a document when Options name no others. Its elements are
// kept in memory, each many times larger than its markup, and the
// functions that go through them call themselves for the elements inside.
const (
	// DefaultMaxDepth is how deep elements may be in elements (the limit
	// of libxml2, the parser of browsers, is 256 too).
	DefaultMaxDepth = 256
	// DefaultMaxNodes is the most elements and runs of text.
	DefaultMaxNodes = 1 << 20
	// DefaultMaxStyleMatches is the most elements that the rules of the
	// style sheets are matched against, one rule and one element at a
	// time.
	DefaultMaxStyleMatches = 1 << 22
)

// Options are the limits of reading a document. A zero or negative limit
// is the default one.
type Options struct {
	// MaxDepth is how deep elements may be in elements (DefaultMaxDepth).
	MaxDepth int
	// MaxNodes is the most elements and runs of text (DefaultMaxNodes).
	MaxNodes int
	// MaxStyleMatches is the most elements that the rules of the style
	// sheets are matched against (DefaultMaxStyleMatches): the elements
	// past it have no Sheet.
	MaxStyleMatches int
}

// ErrNotSVG is the error of Parse for a document that has no svg root
// element (errors.Is tells; the error may say why there is none).
var ErrNotSVG = errors.New("svg: no svg root element")

// Parse reads an SVG document; opts may be nil. It keeps what it read of
// a document with an error or past a limit, as browsers draw up to an
// error, and returns an error only for a document without an svg root
// element.
func Parse(data []byte, opts *Options) (*Document, error) {
	depth, most, matches := DefaultMaxDepth, DefaultMaxNodes, DefaultMaxStyleMatches
	if opts != nil {
		if opts.MaxDepth > 0 {
			depth = opts.MaxDepth
		}
		if opts.MaxNodes > 0 {
			most = opts.MaxNodes
		}
		if opts.MaxStyleMatches > 0 {
			matches = opts.MaxStyleMatches
		}
	}
	// one level more than is read, so that the element that has it ends
	// the reading here
	r := xmltree.Open(data, xmlro.Options{Lenient: true, Entities: htmlentity.Lookup, AutoClose: htmlentity.AutoClose,
		CharsetReader: charset.NewReaderLabel, MaxDepth: depth + 1})
	doc := &Document{IDs: map[string]*Node{}}
	var stack []*Node
	var styles []string
	var text []byte
	nodes := 0
read:
	for {
		k, err := r.Next()
		if err != nil {
			if doc.Root == nil {
				return nil, fmt.Errorf("%w: %v", ErrNotSVG, err)
			}
			break // keep what was read, as browsers draw up to an error
		}
		switch k {
		case xmlro.EOF:
			break read
		case xmlro.StartElement:
			// what passes a limit ends the reading as an error does
			if len(stack) >= depth {
				doc.Warnings = append(doc.Warnings, fmt.Sprintf("the elements of an SVG image are more than %d deep: the rest are not read", depth))
				break read
			}
			if nodes++; nodes > most {
				doc.Warnings = append(doc.Warnings, fmt.Sprintf("an SVG image has more than %d elements: the rest are not read", most))
				break read
			}
			n := &Node{Name: string(r.LocalName()), Attr: map[string]string{}}
			for {
				name, val, ok := r.NextAttr()
				if !ok {
					break
				}
				local := xmltree.Local(name)
				prefixed := len(local) < len(name)
				if xmlro.Equal(local, "xmlns") || prefixed && xmlro.Equal(name[:len(name)-len(local)], "xmlns:") {
					continue
				}
				if prefixed && xmlro.Equal(local, "href") {
					if _, ok := n.Attr["href"]; ok {
						continue // href wins over xlink:href
					}
				}
				n.Attr[string(local)] = val.String()
			}
			if len(stack) > 0 {
				p := stack[len(stack)-1]
				n.Parent = p
				p.Children = append(p.Children, n)
			} else if doc.Root == nil {
				doc.Root = n
			}
			if id := n.Attr["id"]; id != "" {
				if _, dup := doc.IDs[id]; !dup {
					doc.IDs[id] = n
				}
			}
			stack = append(stack, n)
		case xmlro.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		case xmlro.Text, xmlro.CData:
			if len(stack) == 0 {
				continue
			}
			p := stack[len(stack)-1]
			switch p.Name {
			case "style":
				text = xmltree.AppendText(text[:0], r)
				styles = append(styles, string(text))
			case "text", "tspan", "textPath", "a":
				if nodes++; nodes > most {
					continue
				}
				text = xmltree.AppendText(text[:0], r)
				p.Children = append(p.Children, &Node{Name: TextNode, Text: string(text), Parent: p})
			}
		}
	}
	if doc.Root == nil || doc.Root.Name != "svg" {
		return nil, ErrNotSVG
	}
	if len(styles) > 0 {
		sheet := newStyleSheet(parseCSS(strings.Join(styles, "\n")), matches)
		sheet.apply(doc.Root)
		if sheet.matches > sheet.most {
			doc.Warnings = append(doc.Warnings, "the style sheets of an SVG image have too many rules for its elements: the rest of the elements are drawn without them")
		}
	}
	return doc, nil
}
