package webdoc

// The character encoding of an HTML document, found as
// golang.org/x/net/html/charset.DetermineEncoding finds it, with htmlro
// reading the meta elements in place of that package's tokenizer, so that
// the x/net/html tokenizer is not linked beside htmlro for this alone.

import (
	"bytes"
	"mime"
	"strings"
	"unicode/utf8"

	"github.com/shibukawa/bdf/encoding/htmlro"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

var boms = []struct {
	bom []byte
	enc string
}{
	{[]byte{0xfe, 0xff}, "utf-16be"},
	{[]byte{0xff, 0xfe}, "utf-16le"},
	{[]byte{0xef, 0xbb, 0xbf}, "utf-8"},
}

// DetermineEncoding determines the encoding of an HTML document by
// examining up to the first 1024 bytes of content and the declared
// Content-Type, as the HTML specification and
// golang.org/x/net/html/charset.DetermineEncoding have it: a byte order
// mark, then the charset of the content type, then a meta element, then
// UTF-8 when the bytes are valid UTF-8 with a byte above 0x7F, and
// Windows-1252 otherwise. certain is whether the first two found it.
func DetermineEncoding(content []byte, contentType string) (e encoding.Encoding, name string, certain bool) {
	if len(content) > 1024 {
		content = content[:1024]
	}

	for _, b := range boms {
		if bytes.HasPrefix(content, b.bom) {
			e, name = charset.Lookup(b.enc)
			return e, name, true
		}
	}

	if _, params, err := mime.ParseMediaType(contentType); err == nil {
		if cs, ok := params["charset"]; ok {
			if e, name = charset.Lookup(cs); e != nil {
				return e, name, true
			}
		}
	}

	if len(content) > 0 {
		e, name = prescan(content)
		if e != nil {
			return e, name, false
		}
	}

	// Try to detect UTF-8.
	// First eliminate any partial rune at the end.
	for i := len(content) - 1; i >= 0 && i > len(content)-4; i-- {
		b := content[i]
		if b < 0x80 {
			break
		}
		if utf8.RuneStart(b) {
			content = content[:i]
			break
		}
	}
	hasHighBit := false
	for _, c := range content {
		if c >= 0x80 {
			hasHighBit = true
			break
		}
	}
	if hasHighBit && utf8.Valid(content) {
		return encoding.Nop, "utf-8", false
	}

	return charmap.Windows1252, "windows-1252", false
}

// prescan looks for a meta element that names the encoding, in the first
// bytes of a document, as the specification's prescan does.
func prescan(content []byte) (e encoding.Encoding, name string) {
	const (
		dontKnow = iota
		doNeedPragma
		doNotNeedPragma
	)
	r := htmlro.NewBytesReader(content, htmlro.Options{})
	for {
		k, err := r.Next()
		if err != nil || k == htmlro.EOF {
			return nil, ""
		}
		if k != htmlro.StartTag && k != htmlro.SelfClosingTag || !r.NameIs("meta") {
			continue
		}
		gotPragma := false
		needPragma := dontKnow
		name = ""
		e = nil
		for {
			// A name the tag repeats comes once, with its first value, as
			// the tokenizer has it.
			key, v, ok := r.NextAttr()
			if !ok {
				break
			}
			val := lowerASCII(v.String())
			switch {
			case htmlro.Equal(key, "http-equiv"):
				if val == "content-type" {
					gotPragma = true
				}
			case htmlro.Equal(key, "content"):
				if e == nil {
					name = fromMetaElement(val)
					if name != "" {
						e, name = charset.Lookup(name)
						if e != nil {
							needPragma = doNeedPragma
						}
					}
				}
			case htmlro.Equal(key, "charset"):
				e, name = charset.Lookup(val)
				needPragma = doNotNeedPragma
			}
		}

		if needPragma == dontKnow || needPragma == doNeedPragma && !gotPragma {
			continue
		}

		if strings.HasPrefix(name, "utf-16") {
			name = "utf-8"
			e = encoding.Nop
		}

		if e != nil {
			return e, name
		}
	}
}

// lowerASCII lower-cases the ASCII letters of s, and only those, as the
// prescan compares attribute values.
func lowerASCII(s string) string {
	for i := 0; i < len(s); i++ {
		if c := s[i]; 'A' <= c && c <= 'Z' {
			b := []byte(s)
			for j := i; j < len(b); j++ {
				if c := b[j]; 'A' <= c && c <= 'Z' {
					b[j] = c + 0x20
				}
			}
			return string(b)
		}
	}
	return s
}

// fromMetaElement extracts the charset from the content of a meta element
// with http-equiv="content-type": `text/html; charset=utf-8` gives utf-8.
func fromMetaElement(s string) string {
	for s != "" {
		csLoc := strings.Index(s, "charset")
		if csLoc == -1 {
			return ""
		}
		s = s[csLoc+len("charset"):]
		s = strings.TrimLeft(s, " \t\n\f\r")
		if !strings.HasPrefix(s, "=") {
			continue
		}
		s = s[1:]
		s = strings.TrimLeft(s, " \t\n\f\r")
		if s == "" {
			return ""
		}
		if q := s[0]; q == '"' || q == '\'' {
			s = s[1:]
			closeQuote := strings.IndexRune(s, rune(q))
			if closeQuote == -1 {
				return ""
			}
			return s[:closeQuote]
		}

		end := strings.IndexAny(s, "; \t\n\f\r")
		if end == -1 {
			end = len(s)
		}
		return s[:end]
	}
	return ""
}
