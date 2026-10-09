// Package webdoc reads web documents for the converters that lay them out
// in reader mode (converter/html, converter/markdown, converter/epub): HTML
// and XHTML into the node tree of golang.org/x/net/html, svg elements into
// SVG documents of their own, and data: URLs.
//
// XHTML is XML, and the HTML parser misreads it where the two differ: an
// element written as empty (<a id="p5"/>, <div/>, <script src="…"/>) stays
// open there and takes in what follows it, which for a script or a title
// is the rest of the document. XHTML documents are therefore parsed as XML
// into the same tree the HTML parser makes, and only those that are not
// well-formed fall back to the HTML parser.
package webdoc

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/shibukawa/tinygodriver/encoding/xmlro"
	"github.com/shibukawa/tinygodriver/encoding/xmlro/htmlentity"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"golang.org/x/net/html/charset"

	"github.com/shibukawa/bdf/internal/xmltree"
)

// Namespaces of XHTML documents.
const (
	nsXHTML  = "http://www.w3.org/1999/xhtml"
	nsSVG    = "http://www.w3.org/2000/svg"
	nsMathML = "http://www.w3.org/1998/Math/MathML"
	nsXLink  = "http://www.w3.org/1999/xlink"
	nsXML    = "http://www.w3.org/XML/1998/namespace"
	nsXMLNS  = "http://www.w3.org/2000/xmlns/"
	nsOPS    = "http://www.idpf.org/2007/ops"
)

// MaxDepth is how deeply the elements of an XHTML document may nest, as
// many as the HTML parser allows. The readers of the tree call themselves
// for the content of an element, and a stack that overflows cannot be
// recovered from.
const MaxDepth = 512

// prefixes are the usual prefixes of the namespaces attributes are in.
var prefixes = map[string]string{nsXLink: "xlink", nsXML: "xml", nsOPS: "epub", nsXMLNS: "xmlns"}

// Parse parses an HTML or XHTML document. A document is XHTML when its
// content type says so or it starts with an XML declaration; one that is
// not well-formed XML is parsed as HTML.
func Parse(data []byte, contentType string) (*html.Node, error) {
	if IsXML(data, contentType) {
		if doc, err := ParseXHTML(data); err == nil {
			return doc, nil
		}
	}
	return ParseHTML(data, contentType)
}

// IsXML reports whether a document is XML: its content type is an XML
// type (application/xhtml+xml), or it starts with an XML declaration.
func IsXML(data []byte, contentType string) bool {
	if mt, _, err := mime.ParseMediaType(contentType); err == nil && (strings.HasSuffix(mt, "+xml") || strings.HasSuffix(mt, "/xml")) {
		return true
	}
	b := bytes.TrimPrefix(data[:min(len(data), 512)], []byte("\xef\xbb\xbf"))
	return bytes.HasPrefix(bytes.TrimLeft(b, " \t\r\n"), []byte("<?xml"))
}

// ParseHTML decodes a document to UTF-8 by its byte order mark, the
// charset of its content type or of its meta element, and parses it as
// HTML.
func ParseHTML(data []byte, contentType string) (*html.Node, error) {
	enc, _, certain := charset.DetermineEncoding(data, contentType)
	if !certain && utf8.Valid(data) {
		enc = nil
	}
	var r io.Reader = bytes.NewReader(data)
	if enc != nil {
		r = enc.NewDecoder().Reader(r)
	}
	return html.Parse(r)
}

// ParseXHTML parses an XHTML document into the tree the HTML parser makes:
// elements of XHTML (and of no namespace) are HTML elements, those of SVG
// and MathML are in the "svg" and "math" namespaces, and attributes with a
// prefix keep it in their key ("epub:type", "xml:lang"), as the HTML parser
// leaves them, except on SVG and MathML elements, where they are in the
// "xlink" and "xml" namespaces. An xml:lang attribute also gives the
// element the lang attribute HTML readers look for. Elements of other
// namespaces keep their namespace's URL and no atom. The encoding comes
// from the XML declaration (UTF-8 by default). A document whose elements
// nest deeper than MaxDepth is refused.
func ParseXHTML(data []byte) (*html.Node, error) {
	// one level more than the tree may have, so that the document that has
	// it is refused here
	r := xmltree.Open(data, xmlro.Options{Entities: htmlentity.Lookup, CharsetReader: charset.NewReaderLabel, MaxDepth: MaxDepth + 1})
	doc := &html.Node{Type: html.DocumentNode}
	cur := doc
	depth := 0
	// the text being read: character data comes in pieces (around CDATA
	// sections and processing instructions), which are joined when the
	// text ends
	var text *html.Node
	var pieces []byte
	endText := func() {
		if text != nil {
			text.Data = string(pieces)
			text, pieces = nil, pieces[:0]
		}
	}
read:
	for {
		k, err := r.Next()
		if err != nil {
			return nil, err
		}
		switch k {
		case xmlro.EOF:
			break read
		case xmlro.StartElement:
			if cur == doc && hasElement(doc) {
				return nil, errors.New("more than one root element")
			}
			if depth++; depth > MaxDepth {
				return nil, fmt.Errorf("elements nested deeper than %d", MaxDepth)
			}
			endText()
			n, err := element(r)
			if err != nil {
				return nil, err
			}
			cur.AppendChild(n)
			cur = n
		case xmlro.EndElement:
			if cur.Parent == nil {
				return nil, errors.New("unbalanced end tag")
			}
			endText()
			cur = cur.Parent
			depth--
		case xmlro.Text, xmlro.CData:
			if !wellFormedData(r.Text(), k == xmlro.Text) || bytes.Contains(r.Text(), []byte("]]>")) {
				return nil, errors.New("character data is not well-formed")
			}
			if cur == doc {
				continue // white space around the root element
			}
			if text == nil {
				text = &html.Node{Type: html.TextNode}
				cur.AppendChild(text)
			}
			pieces = xmltree.AppendText(pieces, r)
		case xmlro.Comment:
			if !wellFormedData(r.Text(), false) || bytes.Contains(r.Text(), []byte("--")) {
				return nil, errors.New("comment is not well-formed")
			}
			endText()
			cur.AppendChild(&html.Node{Type: html.CommentNode, Data: string(r.Text())})
		}
	}
	endText()
	if !hasElement(doc) {
		return nil, errors.New("no root element")
	}
	return doc, nil
}

// wellFormedData reports whether character data, a comment or the value of an
// attribute is what XML lets it be, where the reader does not look: UTF-8
// without control characters and, with references, every "&" the start of
// one that the reader replaces (those of an entity it does not know stay
// as they are written). A document that fails is one for the HTML parser,
// which knows more entities and reads "&" alone.
func wellFormedData(v []byte, references bool) bool {
	if !utf8.Valid(v) {
		return false
	}
	for i := 0; i < len(v); i++ {
		switch c := v[i]; {
		case c < 0x20 && c != '\t' && c != '\n' && c != '\r':
			return false
		case c == '&' && references:
			// &#x10FFFF; and the zeros before a number
			j := bytes.IndexByte(v[i:min(len(v), i+32)], ';')
			if j < 0 || !reference(v[i+1:i+j]) {
				return false
			}
			i += j
		}
	}
	return true
}

// reference reports whether a name between "&" and ";" is one of the
// references of XML: a predefined entity or a character by its number.
func reference(name []byte) bool {
	switch string(name) {
	case "lt", "gt", "amp", "apos", "quot":
		return true
	}
	if len(name) < 2 || name[0] != '#' {
		return false
	}
	digits, hex := name[1:], false
	if digits[0] == 'x' {
		digits, hex = digits[1:], true
	}
	for _, c := range digits {
		if !('0' <= c && c <= '9' || hex && ('a' <= c && c <= 'f' || 'A' <= c && c <= 'F')) {
			return false
		}
	}
	return len(digits) > 0
}

func hasElement(n *html.Node) bool {
	for k := n.FirstChild; k != nil; k = k.NextSibling {
		if k.Type == html.ElementNode {
			return true
		}
	}
	return false
}

// element makes the node of the XML element whose start a reader is on.
func element(r *xmlro.Reader) (*html.Node, error) {
	name := xmltree.ElementName(r)
	n := &html.Node{Type: html.ElementNode, Data: name.Local}
	foreign := false
	switch name.Space {
	case "", nsXHTML:
		n.Data = strings.ToLower(n.Data)
		n.DataAtom = atom.Lookup([]byte(n.Data))
	case nsSVG:
		n.Namespace, foreign = "svg", true
		n.DataAtom = atom.Lookup([]byte(strings.ToLower(n.Data)))
	case nsMathML:
		n.Namespace, foreign = "math", true
		n.DataAtom = atom.Lookup([]byte(strings.ToLower(n.Data)))
	default:
		n.Namespace = name.Space
	}
	lang, hasLang := "", false
	for {
		qname, val, ok := r.NextAttr()
		if !ok {
			break
		}
		if !wellFormedData(val, true) || bytes.IndexByte(val, '<') >= 0 {
			return nil, errors.New("attribute value is not well-formed")
		}
		a := xmltree.AttrName(r, qname)
		space, key := a.Space, a.Local
		if space == "xmlns" || space == "" && key == "xmlns" || space == nsXMLNS {
			continue // namespace declarations
		}
		attr := html.Attribute{Key: key, Val: val.String()}
		if space != "" {
			prefix, known := prefixes[space]
			switch {
			case foreign && known:
				attr.Namespace = prefix
			case known:
				attr.Key = prefix + ":" + key
			default:
				attr.Namespace = space
			}
		}
		switch {
		case space == "" && key == "lang":
			hasLang = true
		case space == nsXML && key == "lang":
			lang = attr.Val
		}
		n.Attr = append(n.Attr, attr)
	}
	if lang != "" && !hasLang && !foreign && n.Namespace == "" {
		n.Attr = append(n.Attr, html.Attribute{Key: "lang", Val: lang})
	}
	return n, nil
}

// DataURL decodes a data: URL.
func DataURL(s string) ([]byte, error) {
	_, rest, _ := strings.Cut(s, ":")
	meta, data, ok := strings.Cut(rest, ",")
	if !ok {
		return nil, errors.New("malformed data: URL")
	}
	if strings.HasSuffix(strings.ToLower(meta), ";base64") {
		data = strings.Map(func(r rune) rune {
			if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
				return -1
			}
			return r
		}, data)
		if b, err := base64.StdEncoding.DecodeString(data); err == nil {
			return b, nil
		}
		return base64.RawStdEncoding.DecodeString(strings.TrimRight(data, "="))
	}
	d, err := url.PathUnescape(data)
	return []byte(d), err
}

// IsDataURL reports whether a reference is a data: URL.
func IsDataURL(s string) bool {
	return len(s) >= 5 && strings.EqualFold(s[:5], "data:")
}
