package audio

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
)

// legacy decodes the strings a format stores without saying how they are
// encoded: the "ISO-8859-1" frames of ID3v2, ID3v1, RIFF INFO and the text
// chunks of AIFF. Writers put UTF-8 and the system code page there, so the
// strings of a file are read together: as UTF-8 when all of them are valid
// UTF-8, else as Shift_JIS when all of them are, else as Windows-1252.
type legacy struct {
	enc encoding.Encoding // nil: UTF-8
}

// guessLegacy picks the encoding of a file's legacy strings.
func guessLegacy(samples [][]byte) legacy {
	utf, sjis := true, true
	for _, b := range samples {
		b = bytes.Trim(b, "\x00")
		if len(b) == 0 {
			continue
		}
		if utf && !utf8.Valid(b) {
			utf = false
		}
		if sjis && !utf8.Valid(b) {
			if _, ok := decodeSJIS(b); !ok {
				sjis = false
			}
		}
	}
	switch {
	case utf:
		return legacy{}
	case sjis:
		return legacy{japanese.ShiftJIS}
	}
	return legacy{charmap.Windows1252}
}

// decode decodes a string, up to its first NUL.
func (l legacy) decode(b []byte) string {
	if i := bytes.IndexByte(b, 0); i >= 0 {
		b = b[:i]
	}
	switch {
	case utf8.Valid(b):
		return string(b)
	case l.enc == nil:
		return string(bytes.ToValidUTF8(b, []byte("�")))
	}
	out, err := l.enc.NewDecoder().Bytes(b)
	if err != nil {
		return string(bytes.ToValidUTF8(b, []byte("�")))
	}
	return string(out)
}

// decodeSJIS decodes Shift_JIS, failing on bytes that are not.
func decodeSJIS(b []byte) (string, bool) {
	out, err := japanese.ShiftJIS.NewDecoder().Bytes(b)
	if err != nil {
		return "", false
	}
	for _, r := range string(out) {
		if r == utf8.RuneError || r >= 0x80 && r < 0xA0 {
			return "", false
		}
	}
	return string(out), true
}

// decodeUTF16 decodes UTF-16 text: by its byte order mark when it has one,
// else big-endian when bigEndian says so, else little-endian (the common
// case of "UTF-16" without a mark). It ends at the first NUL character.
func decodeUTF16(b []byte, bigEndian bool) string {
	if len(b) >= 2 {
		switch {
		case b[0] == 0xFF && b[1] == 0xFE:
			bigEndian, b = false, b[2:]
		case b[0] == 0xFE && b[1] == 0xFF:
			bigEndian, b = true, b[2:]
		}
	}
	u := make([]uint16, 0, len(b)/2)
	for i := 0; i+1 < len(b); i += 2 {
		c := uint16(b[i]) | uint16(b[i+1])<<8
		if bigEndian {
			c = uint16(b[i])<<8 | uint16(b[i+1])
		}
		if c == 0 {
			break
		}
		u = append(u, c)
	}
	return string(utf16.Decode(u))
}

// clean normalizes line breaks and drops NULs, control characters and
// the padding of fixed-size fields.
func clean(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	s = strings.Map(func(r rune) rune {
		switch {
		case r == '\n' || r == '\t':
			return r
		case r < 0x20 || r == 0x7F || r >= 0x80 && r < 0xA0 || r == 0xFEFF:
			return -1
		}
		return r
	}, s)
	return strings.TrimSpace(s)
}

// oneLine joins the lines and runs of spaces of a string with single
// spaces.
func oneLine(s string) string {
	return strings.Join(strings.FieldsFunc(s, unicode.IsSpace), " ")
}
