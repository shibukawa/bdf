package xmltree

import (
	"bytes"
	"errors"
	"io"
	"math"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/shibukawa/tinygodriver/encoding/xmlro"
)

// The XML that is not a part of an Office package is written by many
// tools and by hand: SVG with the entities of Illustrator and the HTML of
// a foreignObject, draw.io files, the package documents of EPUB, MusicXML,
// XMP packets. Its readers set the options of the reader they want
// (Lenient for what encoding/xml reads with Strict false, the entities and
// the void elements of the htmlentity package for HTML) and open it here,
// where what the options do not bound is bounded.

// MaxEntityText bounds the text that the entities a document declares
// stand for, all their references taken together. Illustrator's entities
// name namespaces; an entity of many bytes referred to many times makes a
// small file a large text.
const MaxEntityText = 1 << 20

// Open returns a reader of a document that is in memory, with the options
// of its reader (the zero Options for a strict one). A document with a
// UTF-16 byte order mark is decoded, and one in another encoding by
// Options.CharsetReader (Latin1 without one). A DOCTYPE is read, and the
// entities it declares are known unless they stand for more than
// MaxEntityText (BoundEntities). What a Lenient reader gives is valid
// UTF-8: the bytes that are not become U+FFFD. A token may be as long as
// the document.
func Open(data []byte, o xmlro.Options) *xmlro.Reader {
	marked := hasUTF16Mark(data)
	data = UTF8(data)
	if o.Lenient && !utf8.Valid(data) {
		// a document in another encoding is converted first
		if enc := DeclaredEncoding(data); enc == "" || enc == "utf-8" {
			data = bytes.ToValidUTF8(data, []byte(string(utf8.RuneError)))
		}
	}
	data, _ = BoundEntities(data)
	convert := o.CharsetReader
	if convert == nil {
		convert = Latin1
	}
	lenient, converted := o.Lenient, false
	o.CharsetReader = func(label string, src io.Reader) (io.Reader, error) {
		if converted {
			return nil, errDeclarations
		}
		converted = true
		// the rest of the document, which is in memory: it is converted
		// whole, and bounded as the document was
		if !marked { // the byte order mark says more than the declaration
			var err error
			if src, err = convert(label, src); err != nil {
				return nil, err
			}
		}
		b, err := io.ReadAll(src)
		if err != nil {
			return nil, err
		}
		if lenient {
			b = bytes.ToValidUTF8(b, []byte(string(utf8.RuneError)))
		}
		b, _ = BoundEntities(b)
		undeclare(b)
		return bytes.NewReader(b), nil
	}
	o.AllowDoctype = true
	if o.MaxBufferBytes == 0 {
		// a converted document is read through a buffer: no encoding takes
		// four times its bytes in UTF-8
		o.MaxBufferBytes = int(min(4*int64(len(data))+MaxEntityText, math.MaxInt32))
	}
	return xmlro.NewBytesReader(data, o)
}

// errDeclarations is the error of a document that declares an encoding
// again after it was converted.
var errDeclarations = errors.New("xml: the document declares its encoding twice")

// undeclare makes the XML declarations in the rest of a document that was
// converted processing instructions of another name. A document has one
// declaration, at its start; each of those that follow (in documents put
// together from others) would have the rest converted again, and the
// reader's buffer copied. Comments and CDATA sections are passed over:
// what they hold is their text. A declaration this does not find (after
// markup that is not well formed) ends the reading (errDeclarations).
func undeclare(b []byte) {
	if !bytes.Contains(b, []byte("<?xml")) {
		return
	}
	skip := func(i int, open, end string) int {
		if j := bytes.Index(b[i+len(open):], []byte(end)); j >= 0 {
			return i + len(open) + j + len(end)
		}
		return i + len(open)
	}
	for i := 0; ; {
		j := bytes.IndexByte(b[i:], '<')
		if j < 0 {
			return
		}
		i += j
		switch rest := b[i:]; {
		case bytes.HasPrefix(rest, []byte("<!--")):
			i = skip(i, "<!--", "-->")
		case bytes.HasPrefix(rest, []byte("<![CDATA[")):
			i = skip(i, "<![CDATA[", "]]>")
		case len(rest) > 5 && string(rest[:5]) == "<?xml" && (rest[5] == ' ' || rest[5] == '\t' || rest[5] == '\r' || rest[5] == '\n'):
			b[i+4] = 'L'
			i += 5
		default:
			i++
		}
	}
}

// Latin1 is the CharsetReader of the readers that have none of their own:
// Latin-1 (and Windows-1252, read as Latin-1) is decoded, and the other
// encodings keep their ASCII.
func Latin1(label string, src io.Reader) (io.Reader, error) {
	switch strings.ToLower(label) {
	case "iso-8859-1", "latin1", "latin-1", "iso_8859-1", "l1", "windows-1252", "cp1252", "us-ascii", "ascii":
		b, err := io.ReadAll(src)
		if err != nil {
			return nil, err
		}
		rs := make([]rune, len(b))
		for i, c := range b {
			rs[i] = rune(c)
		}
		return strings.NewReader(string(rs)), nil
	}
	return src, nil
}

func hasUTF16Mark(b []byte) bool {
	return bytes.HasPrefix(b, []byte{0xfe, 0xff}) || bytes.HasPrefix(b, []byte{0xff, 0xfe})
}

// UTF8 converts a document with a UTF-16 byte order mark to UTF-8 and
// drops a UTF-8 byte order mark.
func UTF8(b []byte) []byte {
	if !hasUTF16Mark(b) {
		return bytes.TrimPrefix(b, []byte("\xef\xbb\xbf"))
	}
	bigEndian := b[0] == 0xfe
	u := make([]uint16, 0, len(b)/2)
	for i := 2; i+1 < len(b); i += 2 {
		if bigEndian {
			u = append(u, uint16(b[i])<<8|uint16(b[i+1]))
		} else {
			u = append(u, uint16(b[i+1])<<8|uint16(b[i]))
		}
	}
	return []byte(string(utf16.Decode(u)))
}

// DeclaredEncoding returns the encoding of the XML declaration a document
// starts with (after white space, which some have before it), lower case,
// or "".
func DeclaredEncoding(b []byte) string {
	b = bytes.TrimLeft(b[:min(len(b), 1024)], " \t\r\n")
	b = b[:min(len(b), 256)]
	if !bytes.HasPrefix(b, []byte("<?xml")) {
		return ""
	}
	end := bytes.Index(b, []byte("?>"))
	if end < 0 {
		return ""
	}
	decl := string(b[:end])
	i := strings.Index(decl, "encoding")
	if i < 0 {
		return ""
	}
	v := strings.TrimLeft(decl[i+len("encoding"):], " \t\r\n=")
	if v == "" || v[0] != '"' && v[0] != '\'' {
		return ""
	}
	q := v[0]
	v = v[1:]
	if j := strings.IndexByte(v, q); j >= 0 {
		v = v[:j]
	}
	return strings.ToLower(strings.TrimSpace(v))
}

// BoundEntities returns a document whose entities (those of the internal
// subset of its DOCTYPE, which the reader replaces the references to)
// stand for MaxEntityText at most, all their references taken together.
// The declarations of a document whose entities stand for more are
// unmade, in a copy, so that its references stay as they are written;
// dropped reports it.
func BoundEntities(data []byte) (bounded []byte, dropped bool) {
	if !bytes.Contains(data, entityDecl) {
		return data, false
	}
	// Every "<!ENTITY" counts, wherever it is: more than the reader takes
	// for declarations, never less.
	sizes := map[string]int{}
	longest := 0
	for rest := data; ; {
		i := bytes.Index(rest, entityDecl)
		if i < 0 {
			break
		}
		rest = rest[i+len(entityDecl):]
		if name, size, ok := declaredEntity(rest); ok {
			sizes[string(name)] = max(sizes[string(name)], size)
			longest = max(longest, len(name))
		}
	}
	if len(sizes) == 0 {
		return data, false
	}
	total := 0
	for rest := data; ; {
		i := bytes.IndexByte(rest, '&')
		if i < 0 {
			return data, false
		}
		rest = rest[i+1:]
		if j := bytes.IndexByte(rest[:min(len(rest), longest+1)], ';'); j > 0 {
			if total += sizes[string(rest[:j])]; total > MaxEntityText {
				break
			}
		}
	}
	// a declaration in lower case is one the reader passes over
	return bytes.ReplaceAll(data, entityDecl, []byte("<!entity")), true
}

// declaredEntity reads the name of a general entity and the size of its
// value from what follows "<!ENTITY", as the reader reads a declaration.
func declaredEntity(b []byte) (name []byte, size int, ok bool) {
	space := func(c byte) bool { return c == ' ' || c == '\t' || c == '\r' || c == '\n' }
	if len(b) == 0 || !space(b[0]) {
		return nil, 0, false
	}
	i := 0
	for i < len(b) && space(b[i]) {
		i++
	}
	start := i
	for i < len(b) && (b[i] >= 0x80 || 'a' <= b[i] && b[i] <= 'z' || 'A' <= b[i] && b[i] <= 'Z' || '0' <= b[i] && b[i] <= '9' || b[i] == '_' || b[i] == '-' || b[i] == '.' || b[i] == ':') {
		i++
	}
	name = b[start:i]
	for i < len(b) && space(b[i]) {
		i++
	}
	if len(name) == 0 || i >= len(b) || b[i] != '"' && b[i] != '\'' {
		return nil, 0, false
	}
	size = bytes.IndexByte(b[i+1:], b[i])
	return name, size, size >= 0
}

// ElementName returns the name of the StartElement a reader is on.
func ElementName(r *xmlro.Reader) Name {
	n := Name{Local: string(r.LocalName())}
	if p := r.Prefix(); p != nil {
		if uri, ok := r.LookupNamespace(p); ok {
			n.Space = uri
		} else {
			n.Space = string(p)
		}
	} else {
		n.Space = r.Namespace()
	}
	return n
}

// Attrs returns the attributes of the StartElement a reader is on, in
// their order, with the entities of their values replaced.
func Attrs(r *xmlro.Reader) []Attr {
	var out []Attr
	for {
		name, val, ok := r.NextAttr()
		if !ok {
			return out
		}
		out = append(out, Attr{AttrName(r, name), val.String()})
	}
}

// AttrName returns the name of an attribute of the StartElement a reader
// is on, from its name as NextAttr gives it.
func AttrName(r *xmlro.Reader, name []byte) Name {
	i := bytes.IndexByte(name, ':')
	if i <= 0 || i == len(name)-1 {
		return Name{Local: string(name)} // no default namespace for attributes
	}
	n := Name{Local: string(name[i+1:])}
	prefix := name[:i]
	switch uri, ok := r.LookupNamespace(prefix); {
	case xmlro.Equal(prefix, "xmlns"):
		n.Space = "xmlns"
	case ok:
		n.Space = uri
	default:
		n.Space = string(prefix)
	}
	return n
}

// Local returns the part of a qualified name after its prefix: the name
// that encoding/xml calls local.
func Local(name []byte) []byte {
	if i := bytes.IndexByte(name, ':'); i > 0 && i < len(name)-1 {
		return name[i+1:]
	}
	return name
}

// AppendText appends the character data of the Text or CData token a
// reader is on: the entities of text replaced, and the line ends of both
// made "\n".
func AppendText(dst []byte, r *xmlro.Reader) []byte {
	if r.Kind() == xmlro.Text {
		return r.Text().AppendTo(dst)
	}
	t := r.Text()
	if bytes.IndexByte(t, '\r') < 0 {
		return append(dst, t...)
	}
	for i := 0; i < len(t); i++ {
		c := t[i]
		if c == '\r' {
			c = '\n'
			if i+1 < len(t) && t[i+1] == '\n' {
				i++
			}
		}
		dst = append(dst, c)
	}
	return dst
}
