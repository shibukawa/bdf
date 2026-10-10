package htmlro

import (
	"bytes"
	"strings"
	"unicode/utf8"
	"unsafe"
)

// Equal reports whether b holds exactly the bytes of s. It allocates on
// neither compiler: the Go compiler elides the conversion in string(b) == s,
// but TinyGo does not, and copies b for every comparison, every switch on
// string(b) and every map index by string(b). Compare names through this,
// NameIs or Value.Equal when the binary is a TinyGo one.
func Equal(b []byte, s string) bool {
	return len(b) == len(s) && unsafe.String(unsafe.SliceData(b), len(b)) == s
}

// EqualFold reports whether b equals s under ASCII case folding of b. s is
// lower case.
func EqualFold(b []byte, s string) bool {
	if len(b) != len(s) {
		return false
	}
	for i := range b {
		c := b[i]
		if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
		}
		if c != s[i] {
			return false
		}
	}
	return true
}

// The bits of Value.flags: what decoding does to the bytes, which depends on
// where they were read.
const (
	valEntities uint8 = 1 << iota // character references are decoded
	valAttr                       // with the rule for legacy references in attribute values
	valNUL                        // NUL becomes U+FFFD
	valLineEnds                   // "\r\n" and "\r" become "\n"
)

// valReader is what every Value a Reader hands out decodes: the line ends,
// and the rest as the token says.
const valReader = valLineEnds

// Value is text, comment or attribute content as it appears in the
// document: character references undecoded, carriage returns and NULs as
// written. Its methods decode on demand, so content with none of those,
// which is nearly all of it, is never copied. A Value aliases the reader's
// buffer and is valid until the reader advances.
//
// Decoding is what the HTML tokenizer does to the bytes: "\r\n" and a lone
// "\r" become "\n"; a NUL becomes U+FFFD in an attribute value, a comment
// and the text of a raw text element (the parser handles the NULs of other
// text); and in an attribute value or text that is not raw, a character
// reference, named or numeric, becomes the characters it names, with the
// legacy references that need no semicolon and the Windows-1252 numeric
// remappings of the specification.
type Value struct {
	raw   []byte
	flags uint8
}

// Raw returns the bytes as written.
func (v Value) Raw() []byte { return v.raw }

// NeedsDecoding reports whether decoding would change the bytes: a
// carriage return, a NUL to replace, or an ampersand where references are
// decoded.
func (v Value) NeedsDecoding() bool {
	if v.flags&valLineEnds != 0 && bytes.IndexByte(v.raw, '\r') >= 0 {
		return true
	}
	if v.flags&valEntities != 0 && bytes.IndexByte(v.raw, '&') >= 0 {
		return true
	}
	return v.flags&valNUL != 0 && bytes.IndexByte(v.raw, 0) >= 0
}

// AppendTo appends the decoded content to dst.
func (v Value) AppendTo(dst []byte) []byte {
	if !v.NeedsDecoding() {
		return append(dst, v.raw...)
	}
	return decode(dst, v.raw, v.flags)
}

// String returns the decoded content as a new string.
func (v Value) String() string {
	if !v.NeedsDecoding() {
		return string(v.raw)
	}
	var tmp [128]byte
	return string(decode(tmp[:0], v.raw, v.flags))
}

// Equal reports whether the decoded content equals s.
func (v Value) Equal(s string) bool {
	if !v.NeedsDecoding() {
		return Equal(v.raw, s)
	}
	var tmp [128]byte
	return Equal(decode(tmp[:0], v.raw, v.flags), s)
}

// decode appends b decoded as flags say to dst.
func decode(dst, b []byte, flags uint8) []byte {
	for i := 0; i < len(b); {
		c := b[i]
		switch {
		case c == '\r' && flags&valLineEnds != 0:
			dst = append(dst, '\n')
			i++
			if i < len(b) && b[i] == '\n' {
				i++
			}
		case c == 0 && flags&valNUL != 0:
			dst = append(dst, "�"...)
			i++
		case c == '&' && flags&valEntities != 0:
			s, r, n := entityAt(b[i:], flags&valAttr != 0)
			switch {
			case n == 1:
				dst = append(dst, '&')
			case s != "":
				dst = append(dst, s...)
			default:
				dst = utf8.AppendRune(dst, r)
			}
			i += n
		default:
			dst = append(dst, c)
			i++
		}
	}
	return dst
}

// UnescapeString decodes the character references of s as the tokenizer
// decodes text outside an attribute value, which is what
// golang.org/x/net/html.UnescapeString does: "a&lt;b" becomes "a<b". Line
// ends and NULs are left as they are. It returns s itself when s has no
// ampersand.
func UnescapeString(s string) string {
	if strings.IndexByte(s, '&') < 0 {
		return s
	}
	// A byte view of s, read and not written; it does not outlive the call.
	b := unsafe.Slice(unsafe.StringData(s), len(s))
	return string(decode(make([]byte, 0, len(s)), b, valEntities))
}
