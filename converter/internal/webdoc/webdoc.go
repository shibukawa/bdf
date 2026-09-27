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
	"encoding/xml"
	"errors"
	"io"
	"mime"
	"net/url"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"golang.org/x/net/html/charset"
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
// from the XML declaration (UTF-8 by default).
func ParseXHTML(data []byte) (*html.Node, error) {
	d := xml.NewDecoder(bytes.NewReader(bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))))
	d.Strict = true
	d.Entity = xml.HTMLEntity
	d.CharsetReader = charset.NewReaderLabel
	doc := &html.Node{Type: html.DocumentNode}
	cur := doc
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if cur == doc && hasElement(doc) {
				return nil, errors.New("more than one root element")
			}
			n := element(t)
			cur.AppendChild(n)
			cur = n
		case xml.EndElement:
			if cur.Parent == nil {
				return nil, errors.New("unbalanced end tag")
			}
			cur = cur.Parent
		case xml.CharData:
			if cur == doc {
				continue // white space around the root element
			}
			if last := cur.LastChild; last != nil && last.Type == html.TextNode {
				last.Data += string(t)
			} else {
				cur.AppendChild(&html.Node{Type: html.TextNode, Data: string(t)})
			}
		case xml.Comment:
			cur.AppendChild(&html.Node{Type: html.CommentNode, Data: string(t)})
		}
	}
	if !hasElement(doc) {
		return nil, errors.New("no root element")
	}
	return doc, nil
}

func hasElement(n *html.Node) bool {
	for k := n.FirstChild; k != nil; k = k.NextSibling {
		if k.Type == html.ElementNode {
			return true
		}
	}
	return false
}

// element makes the node of an XML element.
func element(t xml.StartElement) *html.Node {
	n := &html.Node{Type: html.ElementNode, Data: t.Name.Local}
	foreign := false
	switch t.Name.Space {
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
		n.Namespace = t.Name.Space
	}
	lang, hasLang := "", false
	for _, a := range t.Attr {
		space, key := a.Name.Space, a.Name.Local
		if space == "xmlns" || space == "" && key == "xmlns" || space == nsXMLNS {
			continue // namespace declarations
		}
		attr := html.Attribute{Key: key, Val: a.Value}
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
			lang = a.Value
		}
		n.Attr = append(n.Attr, attr)
	}
	if lang != "" && !hasLang && !foreign && n.Namespace == "" {
		n.Attr = append(n.Attr, html.Attribute{Key: "lang", Val: lang})
	}
	return n
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
