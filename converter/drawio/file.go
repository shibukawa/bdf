package drawio

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/zlib"
	"github.com/shibukawa/tinygodriver/encoding/xmlro"
	"github.com/shibukawa/tinygodriver/encoding/xmlro/htmlentity"

	"github.com/shibukawa/bdf/internal/xmltree"
)

// A draw.io file holds one or more diagrams (pages). They come as
//
//   - .drawio / .xml: an <mxfile> with a <diagram> per page, each holding an
//     <mxGraphModel> element, or the model compressed (raw deflate, base64,
//     with the XML URL-encoded before compression); a bare <mxGraphModel>
//     is a single page;
//   - .drawio.svg: an SVG whose root "content" attribute holds the mxfile;
//   - .drawio.png: a PNG with the mxfile in a tEXt (or zTXt) chunk named
//     "mxfile", URL-encoded.

// file is a parsed draw.io file.
type file struct {
	modified string // mxfile modified attribute (an ISO 8601 date)
	pages    []*page
}

// page is one diagram of a file.
type page struct {
	id, name string
	model    *xmlNode // the mxGraphModel element
}

// errNotDrawio is returned for input that holds no draw.io diagram.
var errNotDrawio = errors.New("not a draw.io file")

// maxDecompressed bounds the size a compressed diagram or inline stencil in
// the input may inflate to, so that a small but highly compressed payload
// cannot exhaust memory (a decompression bomb). It is a var so tests can
// lower it; the default is far above any real diagram.
var maxDecompressed = 64 << 20

// readAllLimited reads r up to maxDecompressed bytes and errors when the
// stream is larger.
func readAllLimited(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, int64(maxDecompressed)+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxDecompressed {
		return nil, fmt.Errorf("decompressed data exceeds %d bytes", maxDecompressed)
	}
	return data, nil
}

// readFile parses draw.io data in any of the supported containers.
func readFile(data []byte) (*file, error) {
	if bytes.HasPrefix(data, []byte("\x89PNG\r\n\x1a\n")) {
		x, err := pngText(data)
		if err != nil {
			return nil, err
		}
		data = x
	}
	root, err := parseXML(data)
	if err != nil {
		return nil, err
	}
	switch root.name {
	case "svg":
		content, ok := root.attrs["content"]
		if !ok {
			return nil, fmt.Errorf("%w: the SVG has no embedded diagram (content attribute)", errNotDrawio)
		}
		return readContent(content)
	case "html":
		// draw.io's "HTML" export embeds the mxfile in a div's data-mxgraph JSON
		return nil, fmt.Errorf("%w: HTML exports are not supported; export the diagram as .drawio", errNotDrawio)
	}
	return readRoot(root)
}

// readContent reads the diagram text of an SVG or PNG export: the mxfile
// XML, or (older exports) the whole of it compressed.
func readContent(s string) (*file, error) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "<") {
		d, err := decompress(s)
		if err != nil {
			return nil, fmt.Errorf("%w: cannot read the embedded diagram: %v", errNotDrawio, err)
		}
		s = d
	}
	root, err := parseXML([]byte(s))
	if err != nil {
		return nil, err
	}
	return readRoot(root)
}

func readRoot(root *xmlNode) (*file, error) {
	f := &file{}
	switch root.name {
	case "mxfile":
		f.modified = root.attrs["modified"]
		for _, d := range root.childrenNamed("diagram") {
			p := &page{id: d.attrs["id"], name: d.attrs["name"]}
			if m := d.child("mxGraphModel"); m != nil {
				p.model = m
			} else if txt := strings.TrimSpace(d.text); txt != "" {
				x, err := decompress(txt)
				if err != nil {
					return nil, fmt.Errorf("drawio: page %q: %v", p.name, err)
				}
				n, err := parseXML([]byte(x))
				if err != nil {
					return nil, fmt.Errorf("drawio: page %q: %v", p.name, err)
				}
				if n.name != "mxGraphModel" {
					return nil, fmt.Errorf("drawio: page %q: expected mxGraphModel, found %s", p.name, n.name)
				}
				p.model = n
			} else {
				p.model = &xmlNode{name: "mxGraphModel"}
			}
			f.pages = append(f.pages, p)
		}
		if len(f.pages) == 0 {
			return nil, fmt.Errorf("%w: the mxfile has no diagram", errNotDrawio)
		}
	case "mxGraphModel":
		f.pages = []*page{{model: root}}
	default:
		return nil, fmt.Errorf("%w: unexpected root element <%s>", errNotDrawio, root.name)
	}
	return f, nil
}

// decompress reverses Graph.compress: base64, raw deflate, then URL
// decoding (files written by some older versions skip the last step).
func decompress(s string) (string, error) {
	s = strings.Map(func(r rune) rune {
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
			return -1
		}
		return r
	}, s)
	raw, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		if raw, err = base64.RawStdEncoding.DecodeString(strings.TrimRight(s, "=")); err != nil {
			return "", fmt.Errorf("bad base64: %v", err)
		}
	}
	inflated, err := readAllLimited(flate.NewReader(bytes.NewReader(raw)))
	if err != nil {
		// some writers use a zlib stream
		zr, zerr := zlib.NewReader(bytes.NewReader(raw))
		if zerr != nil {
			return "", fmt.Errorf("bad deflate data: %v", err)
		}
		if inflated, err = readAllLimited(zr); err != nil {
			return "", fmt.Errorf("bad deflate data: %v", err)
		}
	}
	if u, err := url.PathUnescape(string(inflated)); err == nil {
		return u, nil
	}
	return string(inflated), nil
}

// pngText returns the diagram stored in a PNG export: the tEXt or zTXt
// chunk with the keyword "mxfile" (or "mxGraphModel" in old exports),
// URL-decoded.
func pngText(data []byte) ([]byte, error) {
	p := 8
	for p+12 <= len(data) {
		n := int(binary.BigEndian.Uint32(data[p:]))
		typ := string(data[p+4 : p+8])
		if n < 0 || p+12+n > len(data) {
			break
		}
		body := data[p+8 : p+8+n]
		p += 12 + n
		if typ != "tEXt" && typ != "zTXt" {
			if typ == "IEND" {
				break
			}
			continue
		}
		key, val, ok := bytes.Cut(body, []byte{0})
		if !ok || (string(key) != "mxfile" && string(key) != "mxGraphModel") {
			continue
		}
		if typ == "zTXt" {
			if len(val) < 1 {
				continue
			}
			zr, err := zlib.NewReader(bytes.NewReader(val[1:]))
			if err != nil {
				return nil, fmt.Errorf("drawio: PNG zTXt chunk: %v", err)
			}
			if val, err = readAllLimited(zr); err != nil {
				return nil, fmt.Errorf("drawio: PNG zTXt chunk: %v", err)
			}
		}
		s := string(val)
		if u, err := url.QueryUnescape(strings.ReplaceAll(s, "+", "%2B")); err == nil {
			s = u
		}
		return []byte(s), nil
	}
	return nil, fmt.Errorf("%w: the PNG has no embedded diagram", errNotDrawio)
}

// xmlNode is a generic XML element.
type xmlNode struct {
	name  string
	attrs map[string]string
	kids  []*xmlNode
	text  string // character data directly inside the element
}

func (n *xmlNode) child(name string) *xmlNode {
	if n == nil {
		return nil
	}
	for _, k := range n.kids {
		if k.name == name {
			return k
		}
	}
	return nil
}

func (n *xmlNode) childrenNamed(name string) []*xmlNode {
	if n == nil {
		return nil
	}
	var out []*xmlNode
	for _, k := range n.kids {
		if k.name == name {
			out = append(out, k)
		}
	}
	return out
}

// parseXML reads an XML document into an element tree (names without
// namespace prefixes). It is lenient, as draw.io's files are HTML by
// habit: the entities and the void elements of HTML are known.
func parseXML(data []byte) (*xmlNode, error) {
	r := xmltree.Open(data, xmlro.Options{Lenient: true, Entities: htmlentity.Lookup, AutoClose: htmlentity.AutoClose})
	var stack []*xmlNode
	var root *xmlNode
	var text []byte
read:
	for {
		k, err := r.Next()
		if err != nil {
			if root != nil {
				break
			}
			return nil, fmt.Errorf("%w: %v", errNotDrawio, err)
		}
		switch k {
		case xmlro.EOF:
			break read
		case xmlro.StartElement:
			n := &xmlNode{name: string(r.LocalName()), attrs: map[string]string{}}
			for {
				name, val, ok := r.NextAttr()
				if !ok {
					break
				}
				n.attrs[string(xmltree.Local(name))] = val.String()
			}
			if len(stack) > 0 {
				p := stack[len(stack)-1]
				p.kids = append(p.kids, n)
			} else if root == nil {
				root = n
			}
			stack = append(stack, n)
		case xmlro.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		case xmlro.Text, xmlro.CData:
			if len(stack) > 0 {
				text = xmltree.AppendText(text[:0], r)
				stack[len(stack)-1].text += string(text)
			}
		}
	}
	if root == nil {
		return nil, fmt.Errorf("%w: no XML element", errNotDrawio)
	}
	return root, nil
}
