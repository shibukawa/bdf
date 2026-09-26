package dxf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

// tag is a group code and its value.
type tag struct {
	code int
	s    string  // strings (undecoded until decodeStrings); the text of numbers in text files
	f    float64 // numbers
}

// valueType is the type of a group code's value.
type valueType uint8

const (
	tString valueType = iota
	tFloat
	tInt16
	tInt32
	tInt64
	tBool
	tBinary // hexadecimal text in text files, length-prefixed bytes in binary ones
)

func typeOf(code int) valueType {
	switch {
	case code >= 10 && code < 60, code >= 110 && code < 150, code >= 210 && code < 240,
		code >= 460 && code < 470, code >= 1010 && code < 1060:
		return tFloat
	case code >= 60 && code < 80, code >= 170 && code < 180, code >= 270 && code < 290,
		code >= 370 && code < 390, code >= 400 && code < 410, code >= 1060 && code < 1071:
		return tInt16
	case code >= 90 && code < 100, code >= 420 && code < 430, code >= 440 && code < 460, code == 1071:
		return tInt32
	case code >= 160 && code < 170:
		return tInt64
	case code >= 290 && code < 300:
		return tBool
	case code >= 310 && code < 320, code == 1004:
		return tBinary
	}
	return tString
}

var binarySentinel = []byte("AutoCAD Binary DXF\r\n\x1a\x00")

var errNotDXF = errors.New("not a DXF file")

// readTags splits a DXF file into tags. Numbers are parsed; strings are
// kept as bytes until the encoding is known (decodeStrings).
func readTags(data []byte, warn func(string)) ([]tag, error) {
	if bytes.HasPrefix(data, binarySentinel) {
		return readBinary(data)
	}
	return readText(data, warn)
}

func readText(data []byte, warn func(string)) ([]tag, error) {
	// a UTF-8 byte order mark
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
	var tags []tag
	bad := 0
	line := func() ([]byte, bool) {
		if len(data) == 0 {
			return nil, false
		}
		i := bytes.IndexByte(data, '\n')
		var l []byte
		if i < 0 {
			l, data = data, nil
		} else {
			l, data = data[:i], data[i+1:]
		}
		return bytes.TrimSuffix(l, []byte("\r")), true
	}
	for {
		cl, ok := line()
		if !ok {
			break
		}
		cs := strings.TrimSpace(string(cl))
		if cs == "" && len(data) == 0 {
			break
		}
		code, err := strconv.Atoi(cs)
		if err != nil {
			if len(tags) == 0 {
				return nil, errNotDXF
			}
			return tags, errors.New("invalid group code " + strconv.Quote(cs))
		}
		vl, ok := line()
		if !ok {
			break
		}
		t := tag{code: code}
		switch typeOf(code) {
		case tString, tBinary:
			t.s = string(vl)
		default:
			v := strings.TrimSpace(string(vl))
			t.s = v
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				bad++
			}
			t.f = f
		}
		tags = append(tags, t)
		if code == 0 && t.s == "EOF" {
			break
		}
	}
	if bad > 0 && warn != nil {
		warn(strconv.Itoa(bad) + " numeric value(s) could not be read and were taken as 0")
	}
	if len(tags) == 0 {
		return nil, errNotDXF
	}
	return tags, nil
}

func readBinary(data []byte) ([]tag, error) {
	// R12 binary files have 1-byte group codes (255 escapes a 2-byte one)
	// (the group code 1 after "$ACADVER\0" is 01 in them, 01 00 in later ones)
	r12 := true
	if i := bytes.Index(data[:min(len(data), 1024)], []byte("$ACADVER\x00")); i >= 0 {
		if j := i + 9; j+1 < len(data) && data[j] == 1 && data[j+1] == 0 {
			r12 = false
		}
	}
	var tags []tag
	i := len(binarySentinel)
	errShort := errors.New("the binary DXF data ends in the middle of a value")
	for i < len(data) {
		var code int
		if r12 {
			code = int(data[i])
			i++
			if code == 255 {
				if i+2 > len(data) {
					return tags, errShort
				}
				code = int(binary.LittleEndian.Uint16(data[i:]))
				i += 2
			}
		} else {
			if i+2 > len(data) {
				return tags, errShort
			}
			code = int(binary.LittleEndian.Uint16(data[i:]))
			i += 2
		}
		t := tag{code: code}
		need := func(n int) bool { return i+n <= len(data) }
		switch typeOf(code) {
		case tFloat:
			if !need(8) {
				return tags, errShort
			}
			t.f = math.Float64frombits(binary.LittleEndian.Uint64(data[i:]))
			i += 8
		case tInt16:
			if !need(2) {
				return tags, errShort
			}
			t.f = float64(int16(binary.LittleEndian.Uint16(data[i:])))
			i += 2
		case tInt32:
			if !need(4) {
				return tags, errShort
			}
			t.f = float64(int32(binary.LittleEndian.Uint32(data[i:])))
			i += 4
		case tInt64:
			if !need(8) {
				return tags, errShort
			}
			t.f = float64(int64(binary.LittleEndian.Uint64(data[i:])))
			i += 8
		case tBool:
			if !need(1) {
				return tags, errShort
			}
			t.f = float64(data[i])
			i++
		case tBinary:
			if !need(1) {
				return tags, errShort
			}
			n := int(data[i])
			i++
			if !need(n) {
				return tags, errShort
			}
			t.s = string(data[i : i+n])
			i += n
		default:
			j := bytes.IndexByte(data[i:], 0)
			if j < 0 {
				return tags, errShort
			}
			t.s = string(data[i : i+j])
			i += j + 1
		}
		if typeOf(code) != tString && typeOf(code) != tBinary {
			t.s = strconv.FormatFloat(t.f, 'g', -1, 64)
		}
		tags = append(tags, t)
		if code == 0 && t.s == "EOF" {
			break
		}
	}
	return tags, nil
}

// codepages maps $DWGCODEPAGE values to encodings.
var codepages = map[string]encoding.Encoding{
	"ANSI_874":  charmap.Windows874,
	"ANSI_932":  japanese.ShiftJIS,
	"DOS932":    japanese.ShiftJIS,
	"ANSI_936":  simplifiedchinese.GBK,
	"GB2312":    simplifiedchinese.GBK,
	"ANSI_949":  korean.EUCKR,
	"ANSI_950":  traditionalchinese.Big5,
	"BIG5":      traditionalchinese.Big5,
	"ANSI_1250": charmap.Windows1250,
	"ANSI_1251": charmap.Windows1251,
	"ANSI_1252": charmap.Windows1252,
	"ANSI_1253": charmap.Windows1253,
	"ANSI_1254": charmap.Windows1254,
	"ANSI_1255": charmap.Windows1255,
	"ANSI_1256": charmap.Windows1256,
	"ANSI_1257": charmap.Windows1257,
	"ANSI_1258": charmap.Windows1258,
	"DOS437":    charmap.CodePage437,
	"DOS850":    charmap.CodePage850,
	"DOS852":    charmap.CodePage852,
	"DOS855":    charmap.CodePage855,
	"DOS860":    charmap.CodePage860,
	"DOS861":    charmap.CodePage850,
	"DOS863":    charmap.CodePage863,
	"DOS865":    charmap.CodePage865,
	"DOS866":    charmap.CodePage866,

	"ISO8859-1": charmap.ISO8859_1,
	"ISO8859-2": charmap.ISO8859_2,
	"ISO8859-5": charmap.ISO8859_5,
	"ISO8859-7": charmap.ISO8859_7,
}

// headerValue returns the first value of a header variable in the tags.
func headerValue(tags []tag, name string) (string, bool) {
	for i := 0; i+1 < len(tags); i++ {
		t := tags[i]
		if t.code == 0 && t.s == "ENDSEC" {
			break
		}
		if t.code == 9 && t.s == name {
			return strings.TrimSpace(tags[i+1].s), true
		}
	}
	return "", false
}

// decodeStrings converts the string values into UTF-8: files of AutoCAD
// 2007 and later are UTF-8, older ones are in $DWGCODEPAGE. Files that do
// not say but are valid UTF-8 (or Shift_JIS, when they are not) are read
// as such. It returns the name of the encoding used.
func decodeStrings(tags []tag) string {
	ver, _ := headerValue(tags, "$ACADVER")
	cp, hasCP := headerValue(tags, "$DWGCODEPAGE")
	cp = strings.ToUpper(cp)
	var enc encoding.Encoding
	name := "utf-8"
	switch {
	case ver >= "AC1021":
	case allUTF8(tags):
	case (!hasCP || cp == "ANSI_1252") && validSJIS(tags):
		// Japanese text written without saying so (or saying the default)
		enc, name = japanese.ShiftJIS, "shift_jis"
	default:
		if e, ok := codepages[cp]; ok {
			enc, name = e, strings.ToLower(cp)
		} else {
			enc, name = charmap.Windows1252, "ansi_1252"
		}
	}
	var dec *encoding.Decoder
	if enc != nil {
		dec = enc.NewDecoder()
	}
	for i := range tags {
		t := &tags[i]
		if typeOf(t.code) != tString || t.s == "" {
			continue
		}
		s := t.s
		if dec != nil && !isASCII(s) {
			if out, err := dec.String(s); err == nil {
				s = out
			} else {
				s = strings.ToValidUTF8(s, "�")
			}
		} else if !utf8.ValidString(s) {
			s = strings.ToValidUTF8(s, "�")
		}
		t.s = unescapeUnicode(s)
	}
	return name
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

// allUTF8 reports whether the strings are valid UTF-8 and some are not
// ASCII.
func allUTF8(tags []tag) bool {
	multi := false
	for _, t := range tags {
		if typeOf(t.code) != tString || isASCII(t.s) {
			continue
		}
		if !utf8.ValidString(t.s) {
			return false
		}
		multi = true
	}
	return multi
}

// validSJIS reports whether the non-ASCII strings decode as Shift_JIS into
// Japanese text: kana or kanji, not only the half-width katakana that the
// letters of Western code pages would also decode to.
func validSJIS(tags []tag) bool {
	dec := japanese.ShiftJIS.NewDecoder()
	found := 0
	for _, t := range tags {
		if typeOf(t.code) != tString || isASCII(t.s) {
			continue
		}
		out, err := dec.String(t.s)
		if err != nil || strings.ContainsRune(out, utf8.RuneError) {
			return false
		}
		for _, r := range out {
			if r >= 0x3040 && r < 0x3100 || r >= 0x4e00 && r < 0xa000 {
				found++
			}
		}
	}
	return found >= 2
}

// unescapeUnicode replaces the \U+XXXX escapes of DXF strings and the
// \M+nXXXX multibyte escapes of old files.
func unescapeUnicode(s string) string {
	if !strings.Contains(s, "\\U+") && !strings.Contains(s, "\\u+") && !strings.Contains(s, "\\M+") && !strings.Contains(s, "\\m+") {
		return s
	}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+2 < len(s) && s[i+2] == '+' {
			switch s[i+1] {
			case 'U', 'u':
				if i+7 <= len(s) {
					if v, err := strconv.ParseUint(s[i+3:i+7], 16, 32); err == nil {
						b.WriteRune(rune(v))
						i += 6
						continue
					}
				}
			case 'M', 'm':
				if i+8 <= len(s) {
					if r, ok := mif(s[i+3], s[i+4:i+8]); ok {
						b.WriteRune(r)
						i += 7
						continue
					}
				}
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

// mif decodes a character of the multibyte interchange format: n selects
// the code page and hex holds its two bytes.
func mif(n byte, hex string) (rune, bool) {
	v, err := strconv.ParseUint(hex, 16, 16)
	if err != nil {
		return 0, false
	}
	var enc encoding.Encoding
	switch n {
	case '1':
		enc = japanese.ShiftJIS
	case '2':
		enc = traditionalchinese.Big5
	case '3', '4':
		enc = korean.EUCKR
	case '5':
		enc = simplifiedchinese.GBK
	default:
		return 0, false
	}
	out, err := enc.NewDecoder().Bytes([]byte{byte(v >> 8), byte(v)})
	if err != nil {
		return 0, false
	}
	r, _ := utf8.DecodeRune(out)
	return r, r != utf8.RuneError
}
