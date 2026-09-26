package csv

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/japanese"
	xunicode "golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

// The character encoding of a file is told by its byte order mark, when it
// has one (which also wins over an encoding the caller names, as in
// browsers). Otherwise: UTF-16 by the zero bytes of ASCII characters, the
// 7-bit escape sequences of ISO-2022-JP, UTF-8 when the bytes are UTF-8,
// and else the legacy encoding whose decoding reads most like text:
// Shift_JIS (Windows-31J, as Excel writes Japanese CSV), EUC-JP or
// Windows-1252. Other legacy encodings (GB18030, Big5, EUC-KR …) are
// read when named.

// charset is an encoding a file is read in.
type charset struct {
	name string
	enc  encoding.Encoding // nil for UTF-8
	bom  int               // length of the byte order mark
}

var boms = []struct {
	mark []byte
	cs   charset
}{
	{[]byte{0xEF, 0xBB, 0xBF}, charset{"UTF-8", nil, 3}},
	{[]byte{0xFF, 0xFE, 0, 0}, charset{"UTF-32LE", utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM), 4}},
	{[]byte{0, 0, 0xFE, 0xFF}, charset{"UTF-32BE", utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM), 4}},
	{[]byte{0xFF, 0xFE}, charset{"UTF-16LE", xunicode.UTF16(xunicode.LittleEndian, xunicode.IgnoreBOM), 2}},
	{[]byte{0xFE, 0xFF}, charset{"UTF-16BE", xunicode.UTF16(xunicode.BigEndian, xunicode.IgnoreBOM), 2}},
}

// bomCharset returns the encoding a byte order mark tells.
func bomCharset(data []byte) (charset, bool) {
	for _, b := range boms {
		if bytes.HasPrefix(data, b.mark) {
			return b.cs, true
		}
	}
	return charset{}, false
}

// namedCharset returns an encoding by its WHATWG label ("shift_jis",
// "euc-kr", "utf-16le" …).
func namedCharset(label string) (charset, error) {
	if l, ok := charsetAliases[strings.ToLower(strings.TrimSpace(label))]; ok {
		label = l
	}
	enc, err := htmlindex.Get(label)
	if err != nil {
		return charset{}, fmt.Errorf("unknown character encoding %q", label)
	}
	name, _ := htmlindex.Name(enc)
	if name == "utf-8" {
		return charset{name: "UTF-8"}, nil
	}
	return charset{name: canonicalName(name), enc: enc}, nil
}

// charsetAliases are the Windows code page names of the encodings that
// WHATWG labels otherwise.
var charsetAliases = map[string]string{"cp932": "shift_jis", "cp936": "gbk", "cp949": "euc-kr", "cp950": "big5",
	"cp1252": "windows-1252", "eucjp": "euc-jp", "euckr": "euc-kr"}

func canonicalName(whatwg string) string {
	switch whatwg {
	case "shift_jis":
		return "Shift_JIS"
	case "utf-16le", "utf-16be", "euc-jp", "euc-kr", "iso-2022-jp", "windows-1252", "gbk", "gb18030", "big5":
		return strings.ToUpper(whatwg)
	}
	return whatwg
}

// sniffLen is how much of a file the guesses look at.
const sniffLen = 64 << 10

// detectCharset guesses the encoding of a file without a byte order mark.
func detectCharset(data []byte) charset {
	sample := data[:min(len(data), sniffLen)]
	if cs, ok := utf16Guess(sample); ok {
		return cs
	}
	if isISO2022JP(sample) {
		return charset{name: "ISO-2022-JP", enc: japanese.ISO2022JP}
	}
	if mostlyUTF8(data) {
		return charset{name: "UTF-8"}
	}
	best, bestScore := charset{}, -1e9
	for _, cs := range []charset{
		{name: "Shift_JIS", enc: japanese.ShiftJIS},
		{name: "EUC-JP", enc: japanese.EUCJP},
		{name: "WINDOWS-1252", enc: charmap.Windows1252},
	} {
		text, err := cs.enc.NewDecoder().Bytes(trimPartial(sample, len(data) > len(sample)))
		if err != nil {
			continue
		}
		if s := textScore(string(text), cs.name); s > bestScore {
			best, bestScore = cs, s
		}
	}
	return best
}

// trimPartial drops the last bytes of a sample cut short, which may be
// part of a character.
func trimPartial(sample []byte, cut bool) []byte {
	if !cut {
		return sample
	}
	if i := bytes.LastIndexAny(sample, "\n\r"); i > 0 {
		return sample[:i]
	}
	return sample
}

// utf16Guess recognizes UTF-16 without a byte order mark by the zero high
// byte of the characters below U+0100, which every line of CSV has (its
// delimiters and line breaks): text in other encodings has no zero bytes,
// and binary data has them in both bytes of the pairs.
func utf16Guess(sample []byte) (charset, bool) {
	n := len(sample) / 2
	even, odd := 0, 0
	for i := 0; i+1 < len(sample); i += 2 {
		if sample[i] == 0 {
			even++
		}
		if sample[i+1] == 0 {
			odd++
		}
	}
	switch {
	case even >= max(2, n/50) && odd*10 <= even:
		return charset{"UTF-16BE", xunicode.UTF16(xunicode.BigEndian, xunicode.IgnoreBOM), 0}, true
	case odd >= max(2, n/50) && even*10 <= odd:
		return charset{"UTF-16LE", xunicode.UTF16(xunicode.LittleEndian, xunicode.IgnoreBOM), 0}, true
	}
	return charset{}, false
}

// isISO2022JP reports whether 7-bit text switches to JIS X 0208 (or half-
// width katakana) with the escape sequences of ISO-2022-JP.
func isISO2022JP(sample []byte) bool {
	for _, b := range sample {
		if b >= 0x80 {
			return false
		}
	}
	for _, esc := range []string{"\x1b$B", "\x1b$@", "\x1b(I", "\x1b$(D"} {
		if bytes.Contains(sample, []byte(esc)) {
			return true
		}
	}
	return false
}

// mostlyUTF8 reports whether data is UTF-8: valid, or with a stray byte
// for at least twenty valid multibyte characters (a file edited in another
// encoding once).
func mostlyUTF8(data []byte) bool {
	if utf8.Valid(data) {
		return true
	}
	multi, bad := 0, 0
	for i := 0; i < len(data); {
		if data[i] < utf8.RuneSelf {
			i++
			continue
		}
		r, size := utf8.DecodeRune(data[i:])
		if r == utf8.RuneError && size <= 1 {
			if len(data)-i < utf8.UTFMax && !utf8.FullRune(data[i:]) {
				break // cut short at the end
			}
			bad++
		} else {
			multi++
		}
		i += max(size, 1)
	}
	return multi > 0 && bad*20 <= multi
}

// textScore rates how much a decoding reads like text in its encoding's
// languages: the mean weight of the characters outside ASCII, which is
// high for the kana and common kanji of Japanese and for the accented
// letters and punctuation of Western languages, and low for what a wrong
// decoding makes of them (half-width katakana, rare kanji, C1 controls,
// user-defined characters, U+FFFD). Western text has its accented letters
// among ASCII ones; a run of them is the two bytes of East Asian
// characters read one by one.
func textScore(text, name string) float64 {
	total, n := 0.0, 0
	enc := japanese.ShiftJIS.NewEncoder()
	prevHigh := false
	for _, r := range text {
		if r < utf8.RuneSelf {
			if r < 0x20 && r != '\t' && r != '\n' && r != '\r' {
				total -= 1
				n++
			}
			prevHigh = false
			continue
		}
		n++
		run := prevHigh
		prevHigh = true
		switch {
		case r == utf8.RuneError:
			total -= 5
		case r >= 0x80 && r <= 0x9F, r >= 0xE000 && r <= 0xF8FF:
			total -= 1 // C1 controls, private use
		case name == "WINDOWS-1252":
			switch {
			case run:
				total -= 0.5
			case unicode.IsLetter(r):
				total += 0.8
			case r >= 0x2010 && r <= 0x2044, r == 0x20AC:
				total += 0.8 // dashes, quotes, bullets, ellipsis, euro
			default:
				total += 0.3
			}
		case r >= 0x3041 && r <= 0x309F:
			total += 1 // hiragana
		case r >= 0x30A0 && r <= 0x30FF, r >= 0x3000 && r <= 0x303F:
			total += 0.8 // katakana, CJK punctuation
		case r >= 0xFF01 && r <= 0xFF5E:
			total += 0.6 // full-width ASCII
		case r >= 0xFF61 && r <= 0xFF9F:
			total += 0.2 // half-width katakana
		case unicode.Is(unicode.Han, r):
			if b, err := enc.Bytes([]byte(string(r))); err == nil && len(b) == 2 && b[0] >= 0x88 && b[0] <= 0x98 {
				total += 0.7 // JIS level 1
			} else {
				total += 0.1
			}
		}
	}
	if n == 0 {
		return 0
	}
	return total / float64(n)
}

// decode converts data in an encoding to text, without the byte order
// mark. Bytes that are not text in it become U+FFFD.
func decode(data []byte, cs charset) string {
	data = data[cs.bom:]
	if cs.enc == nil {
		if utf8.Valid(data) {
			return string(data)
		}
		return strings.ToValidUTF8(string(data), "�")
	}
	out, err := cs.enc.NewDecoder().Bytes(data)
	if err != nil {
		return strings.ToValidUTF8(string(out), "�")
	}
	return string(out)
}
