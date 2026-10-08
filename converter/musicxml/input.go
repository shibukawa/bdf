package musicxml

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
	"unicode/utf8"

	"github.com/shibukawa/bdf/converter/internal/ziputil"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/unicode"
)

// maxInput caps the size of the score read: an uncompressed file, or the
// score inside a compressed one.
const maxInput = 64 << 20

// mxlType is the media type of compressed MusicXML (the mimetype entry),
// and rootType that of the score in its container.
const (
	mxlType  = "application/vnd.recordare.musicxml"
	rootType = "application/vnd.recordare.musicxml+xml"
)

// Detect reports whether an input is a MusicXML score: an uncompressed
// file whose first bytes name a score-partwise or score-timewise root (or
// carry the MusicXML document type), whatever its extension, or a
// compressed file (.mxl), a ZIP archive whose mimetype entry or container
// names a MusicXML score.
func Detect(head []byte, r io.ReaderAt, size int64) bool {
	if bytes.HasPrefix(head, []byte("PK\x03\x04")) {
		return detectMXL(head, r, size)
	}
	return isScoreText(headText(head))
}

// isScoreText reports whether the start of a document (as UTF-8) is that
// of a MusicXML score.
func isScoreText(s string) bool {
	s = strings.TrimLeft(s, " \t\r\n")
	if !strings.HasPrefix(s, "<") {
		return false
	}
	if strings.Contains(s, "<score-partwise") || strings.Contains(s, "<score-timewise") {
		return true
	}
	// the public identifiers of the partwise and timewise DTDs, not those
	// of opus, container and sounds files
	i := strings.Index(s, "-//Recordare//DTD MusicXML")
	if i < 0 {
		return false
	}
	id := s[i:]
	if j := strings.Index(id, "//EN"); j >= 0 {
		id = id[:j]
	}
	return strings.Contains(id, "Partwise") || strings.Contains(id, "Timewise")
}

// headText returns the first bytes of an input as UTF-8 text: without a
// UTF-8 byte order mark, decoded from UTF-16.
func headText(head []byte) string {
	if le, ok := utf16Order(head); ok {
		b := head[:len(head)&^1]
		var sb strings.Builder
		for i := 0; i+1 < len(b); i += 2 {
			var c uint16
			if le {
				c = binary.LittleEndian.Uint16(b[i:])
			} else {
				c = binary.BigEndian.Uint16(b[i:])
			}
			if c == 0xFEFF {
				continue
			}
			if c < 0x80 {
				sb.WriteByte(byte(c))
			} else {
				sb.WriteByte('?')
			}
		}
		return sb.String()
	}
	return string(bytes.TrimPrefix(head, []byte("\xef\xbb\xbf")))
}

// utf16Order recognizes UTF-16 text by its byte order mark, or by the
// zero byte of a "<" that starts it without one.
func utf16Order(b []byte) (littleEndian, ok bool) {
	switch {
	case bytes.HasPrefix(b, []byte{0xFF, 0xFE}):
		return true, true
	case bytes.HasPrefix(b, []byte{0xFE, 0xFF}):
		return false, true
	case len(b) >= 4 && b[0] == '<' && b[1] == 0 && b[2] != 0 && b[3] == 0:
		return true, true
	case len(b) >= 4 && b[0] == 0 && b[1] == '<' && b[2] == 0 && b[3] != 0:
		return false, true
	}
	return false, false
}

// detectMXL recognizes a compressed MusicXML file: by the mimetype entry
// stored first, or by META-INF/container.xml naming a MusicXML score.
// Other ZIP packages (Office documents, EPUB) have no such entry or name
// other roots.
func detectMXL(head []byte, r io.ReaderAt, size int64) bool {
	if len(head) >= 30 {
		n, extra := int(binary.LittleEndian.Uint16(head[26:])), int(binary.LittleEndian.Uint16(head[28:]))
		if 30+n <= len(head) && string(head[30:30+n]) == "mimetype" {
			if data := head[min(30+n+extra, len(head)):]; bytes.HasPrefix(data, []byte(mxlType)) {
				return true
			}
		}
	}
	zr, err := ziputil.NewReader(r, size)
	if err != nil {
		return false
	}
	f := zipFile(zr, "META-INF/container.xml")
	if f == nil || f.UncompressedSize64 > 1<<20 {
		return false
	}
	roots, err := containerRoots(f)
	if err != nil {
		return false
	}
	for _, rf := range roots {
		switch {
		case rf.mediaType == rootType:
			return true
		case rf.mediaType == "" && isScorePath(rf.path):
			if e := zipFile(zr, rf.path); e != nil {
				if rc, err := e.Open(); err == nil {
					b, _ := io.ReadAll(io.LimitReader(rc, 1024))
					rc.Close()
					if isScoreText(headText(b)) {
						return true
					}
				}
			}
		}
	}
	return false
}

// isScorePath reports whether a file name is that of a MusicXML score.
func isScorePath(name string) bool {
	ext := strings.ToLower(path.Ext(name))
	return ext == ".musicxml" || ext == ".xml"
}

// zipFile returns the entry of an archive with a name, or nil.
func zipFile(zr *zip.Reader, name string) *zip.File {
	name = cleanZipName(name)
	for _, f := range zr.File {
		if cleanZipName(f.Name) == name {
			return f
		}
	}
	return nil
}

func cleanZipName(name string) string {
	return strings.TrimPrefix(path.Clean("/"+strings.ReplaceAll(name, "\\", "/")), "/")
}

// rootFile is a rootfile of a container.
type rootFile struct {
	path, mediaType string
}

// containerRoots reads the rootfiles of META-INF/container.xml.
func containerRoots(f *zip.File) ([]rootFile, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	d := newDecoder(io.LimitReader(rc, 1<<20))
	var roots []rootFile
	for {
		tok, err := d.Token()
		if err != nil {
			if err == io.EOF {
				return roots, nil
			}
			return roots, err
		}
		if se, ok := tok.(xml.StartElement); ok && se.Name.Local == "rootfile" {
			var rf rootFile
			for _, a := range se.Attr {
				switch a.Name.Local {
				case "full-path":
					rf.path = a.Value
				case "media-type":
					rf.mediaType = strings.TrimSpace(a.Value)
				}
			}
			if rf.path != "" {
				roots = append(roots, rf)
			}
		}
	}
}

// readInput reads the score of an input: the file itself, or the score a
// compressed file holds.
func readInput(r io.ReaderAt, size int64) ([]byte, error) {
	head := make([]byte, 4)
	n, _ := r.ReadAt(head, 0)
	if bytes.HasPrefix(head[:n], []byte("PK\x03\x04")) {
		return readMXL(r, size)
	}
	if size > maxInput {
		return nil, fmt.Errorf("musicxml: the file is larger than %d MiB", maxInput>>20)
	}
	b := make([]byte, size)
	n, err := r.ReadAt(b, 0)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return b[:n], nil
}

// readMXL reads the score of a compressed MusicXML file: the first
// MusicXML rootfile of its container, or without one the first .musicxml
// or .xml file outside META-INF.
func readMXL(r io.ReaderAt, size int64) ([]byte, error) {
	zr, err := ziputil.NewReader(r, size)
	if err != nil {
		return nil, fmt.Errorf("musicxml: compressed file: %w", err)
	}
	var score *zip.File
	if f := zipFile(zr, "META-INF/container.xml"); f != nil && f.UncompressedSize64 <= 1<<20 {
		roots, _ := containerRoots(f)
		for _, rf := range roots {
			if rf.mediaType == rootType || rf.mediaType == "" && isScorePath(rf.path) {
				if score = zipFile(zr, rf.path); score != nil {
					break
				}
			}
		}
	}
	if score == nil {
		for _, f := range zr.File {
			name := cleanZipName(f.Name)
			if !strings.HasPrefix(name, "META-INF/") && isScorePath(name) {
				score = f
				break
			}
		}
	}
	if score == nil {
		return nil, errors.New("musicxml: the compressed file holds no score")
	}
	if score.UncompressedSize64 > maxInput {
		return nil, fmt.Errorf("musicxml: the score is larger than %d MiB", maxInput>>20)
	}
	rc, err := score.Open()
	if err != nil {
		return nil, fmt.Errorf("musicxml: %s: %w", score.Name, err)
	}
	defer rc.Close()
	b, err := io.ReadAll(io.LimitReader(rc, maxInput+1))
	if err != nil {
		return nil, fmt.Errorf("musicxml: %s: %w", score.Name, err)
	}
	if len(b) > maxInput {
		return nil, fmt.Errorf("musicxml: the score is larger than %d MiB", maxInput>>20)
	}
	return b, nil
}

// toUTF8 converts a document to UTF-8 for the XML decoder: UTF-16 (with
// a byte order mark or starting with "<") is decoded, a UTF-8 byte order
// mark dropped, and a document said to be UTF-8 (or saying nothing) that
// is not is read as Windows-1252. Control characters XML forbids become
// spaces. It reports what it had to repair.
func toUTF8(b []byte) ([]byte, string) {
	var note string
	if le, ok := utf16Order(b); ok {
		order := unicode.BigEndian
		if le {
			order = unicode.LittleEndian
		}
		if d, err := unicode.UTF16(order, unicode.IgnoreBOM).NewDecoder().Bytes(b[:len(b)&^1]); err == nil {
			b = d
		}
	}
	b = bytes.TrimPrefix(b, []byte("\xef\xbb\xbf"))
	if enc := declaredEncoding(b); (enc == "" || isUnicodeLabel(enc)) && !utf8.Valid(b) {
		if d, err := charmap.Windows1252.NewDecoder().Bytes(b); err == nil {
			b, note = d, "the text is not UTF-8; it is read as Windows-1252"
		}
	}
	for i, c := range b {
		if c < 0x20 && c != '\t' && c != '\n' && c != '\r' {
			if note == "" {
				note = "control characters are replaced with spaces"
			}
			b[i] = ' '
		}
	}
	return b, note
}

// declaredEncoding returns the encoding of the XML declaration, lower
// case, or "".
func declaredEncoding(b []byte) string {
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

func isUnicodeLabel(label string) bool {
	switch strings.ToLower(label) {
	case "utf-8", "utf8", "utf-16", "utf16", "utf-16le", "utf-16be", "ucs-2", "iso-10646-ucs-2", "unicode":
		return true
	}
	return false
}

// newDecoder returns a lenient XML decoder: the HTML entities are known,
// unknown entities are kept as text, and the document's encoding is
// decoded (the text is already UTF-8 when it says UTF-8 or UTF-16).
func newDecoder(r io.Reader) *xml.Decoder {
	d := xml.NewDecoder(r)
	d.Strict = false
	d.Entity = xml.HTMLEntity
	d.CharsetReader = func(label string, in io.Reader) (io.Reader, error) {
		if isUnicodeLabel(label) {
			return in, nil
		}
		switch strings.ToLower(label) {
		case "iso-8859-1", "iso8859-1", "latin1", "latin-1", "l1", "iso_8859-1":
			return charmap.ISO8859_1.NewDecoder().Reader(in), nil
		case "us-ascii", "ascii":
			return in, nil
		}
		if e, err := htmlindex.Get(label); err == nil {
			return e.NewDecoder().Reader(in), nil
		}
		return in, nil
	}
	return d
}

// node is an element of the document: its name, attributes, children and
// the text directly in it.
type node struct {
	name  string
	attrs []xml.Attr
	kids  []*node
	text  string
}

// errTooLarge stops reading an element with too many descendants.
var errTooLarge = errors.New("musicxml: an element holds too many elements")

// maxNodes caps the elements read into one tree (a measure, or an element
// of the header).
const maxNodes = 1 << 20

// readNode reads the element that start begins, with its descendants.
func readNode(d *xml.Decoder, start xml.StartElement) (*node, error) {
	root := &node{name: start.Name.Local, attrs: start.Attr}
	stack := []*node{root}
	text := [][]byte{nil}
	count := 1
	for len(stack) > 0 {
		tok, err := d.Token()
		if err != nil {
			return root, err
		}
		top := stack[len(stack)-1]
		switch t := tok.(type) {
		case xml.StartElement:
			if count++; count > maxNodes {
				return root, errTooLarge
			}
			n := &node{name: t.Name.Local, attrs: t.Attr}
			top.kids = append(top.kids, n)
			stack = append(stack, n)
			text = append(text, nil)
		case xml.CharData:
			if b := text[len(text)-1]; len(b) < 1<<16 {
				text[len(text)-1] = append(b, t...)
			}
		case xml.EndElement:
			// the white space between child elements is not kept
			if b := text[len(text)-1]; len(top.kids) == 0 || len(bytes.TrimSpace(b)) > 0 {
				top.text = string(b)
			}
			stack = stack[:len(stack)-1]
			text = text[:len(text)-1]
		}
	}
	return root, nil
}

// attr returns the value of an attribute, or "".
func (n *node) attr(name string) string {
	v, _ := n.attrOK(name)
	return v
}

// attrOK returns the value of an attribute and whether it is there.
func (n *node) attrOK(name string) (string, bool) {
	if n == nil {
		return "", false
	}
	for _, a := range n.attrs {
		if a.Name.Local == name {
			return strings.TrimSpace(a.Value), true
		}
	}
	return "", false
}

// child returns the first child with a name, or nil.
func (n *node) child(name string) *node {
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

// textOf returns the trimmed text of the first child with a name.
func (n *node) textOf(name string) string {
	return n.child(name).trim()
}

// trim returns the text of the element without surrounding white space.
func (n *node) trim() string {
	if n == nil {
		return ""
	}
	return strings.TrimSpace(n.text)
}

// yes reports whether an attribute is "yes".
func (n *node) yes(name string) bool {
	return n.attr(name) == "yes"
}
