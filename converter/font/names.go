package font

import (
	"strconv"
	"strings"
	"unicode/utf16"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"

	"github.com/shibukawa/bdf/internal/cff"
	"github.com/shibukawa/bdf/internal/sfnt"
)

// nameRec is a string of the name table.
type nameRec struct {
	id   uint16
	lang string // BCP 47 when known, else the platform's language ID
	en   bool   // an English record
	text string
	rank int // how good a record of its ID it is: Windows, then Unicode, then Macintosh
}

// parseNames reads the strings of a name table, of every platform and
// language that can be decoded.
func parseNames(b []byte) []nameRec {
	if len(b) < 6 {
		return nil
	}
	count, strOff := int(sfnt.BE16(b, 2)), int(sfnt.BE16(b, 4))
	var langTags []string
	if sfnt.BE16(b, 0) == 1 {
		p := 6 + count*12
		n := int(sfnt.BE16(b, p))
		for i := range n {
			l, o := int(sfnt.BE16(b, p+2+i*4)), strOff+int(sfnt.BE16(b, p+4+i*4))
			if o+l <= len(b) {
				langTags = append(langTags, utf16be(b[o:o+l]))
			} else {
				langTags = append(langTags, "")
			}
		}
	}
	var out []nameRec
	seen := map[string]bool{}
	for i := range count {
		rec := 6 + i*12
		if rec+12 > len(b) {
			break
		}
		pid, eid, lid, id := sfnt.BE16(b, rec), sfnt.BE16(b, rec+2), sfnt.BE16(b, rec+4), sfnt.BE16(b, rec+6)
		l, o := int(sfnt.BE16(b, rec+8)), strOff+int(sfnt.BE16(b, rec+10))
		if o+l > len(b) || l == 0 {
			continue
		}
		raw := b[o : o+l]
		r := nameRec{id: id}
		switch pid {
		case 0:
			r.text, r.rank = utf16be(raw), 2
			r.lang = langTag(lid, langTags, nil)
		case 3:
			switch eid {
			case 0, 1, 10:
				r.text = utf16be(raw)
			case 2, 3, 4, 5, 6:
				r.text = legacy(raw, eid)
			default:
				continue
			}
			r.rank = 3
			r.lang = langTag(lid, langTags, windowsLangs)
			r.en = lid&0x3FF == 0x09
		case 1:
			var enc encoding.Encoding
			switch eid {
			case 0:
				enc = charmap.Macintosh
			case 1:
				enc = japanese.ShiftJIS
			case 2:
				enc = traditionalchinese.Big5
			case 3:
				enc = korean.EUCKR
			case 25:
				enc = simplifiedchinese.GBK
			default:
				continue
			}
			t, err := enc.NewDecoder().Bytes(raw)
			if err != nil {
				continue
			}
			r.text, r.rank = string(t), 1
			r.lang = langTag(lid, langTags, macLangs)
			r.en = lid == 0
		default:
			continue
		}
		r.text = strings.TrimRight(r.text, "\x00")
		if strings.TrimSpace(r.text) == "" {
			continue
		}
		// the same string of the same language on several platforms once
		k := strconv.Itoa(int(id)) + "\x00" + r.lang + "\x00" + r.text
		if seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, r)
	}
	return out
}

func utf16be(b []byte) string {
	u := make([]uint16, len(b)/2)
	for k := range u {
		u[k] = sfnt.BE16(b, k*2)
	}
	return string(utf16.Decode(u))
}

// legacy decodes a Windows record in a legacy East Asian encoding: the
// bytes of the encoding stored in 16-bit units.
func legacy(raw []byte, eid uint16) string {
	var bs []byte
	for k := 0; k+1 < len(raw); k += 2 {
		if raw[k] != 0 {
			bs = append(bs, raw[k])
		}
		bs = append(bs, raw[k+1])
	}
	var enc encoding.Encoding
	switch eid {
	case 2:
		enc = japanese.ShiftJIS
	case 3:
		enc = simplifiedchinese.GBK
	case 4:
		enc = traditionalchinese.Big5
	case 5:
		enc = korean.EUCKR
	default:
		return ""
	}
	t, err := enc.NewDecoder().Bytes(bs)
	if err != nil {
		return ""
	}
	return string(t)
}

// langTag names the language of a record.
func langTag(lid uint16, tags []string, known map[uint16]string) string {
	if lid >= 0x8000 {
		if i := int(lid - 0x8000); i < len(tags) {
			return tags[i]
		}
		return ""
	}
	if t, ok := known[lid]; ok {
		return t
	}
	return "0x" + strings.ToUpper(strconv.FormatUint(uint64(lid), 16))
}

var windowsLangs = map[uint16]string{
	0x0401: "ar", 0x0402: "bg", 0x0403: "ca", 0x0404: "zh-TW", 0x0405: "cs", 0x0406: "da", 0x0407: "de",
	0x0408: "el", 0x0409: "en", 0x040A: "es", 0x040B: "fi", 0x040C: "fr", 0x040D: "he", 0x040E: "hu",
	0x040F: "is", 0x0410: "it", 0x0411: "ja", 0x0412: "ko", 0x0413: "nl", 0x0414: "nb", 0x0415: "pl",
	0x0416: "pt-BR", 0x0418: "ro", 0x0419: "ru", 0x041A: "hr", 0x041B: "sk", 0x041D: "sv", 0x041E: "th",
	0x041F: "tr", 0x0421: "id", 0x0422: "uk", 0x0423: "be", 0x0424: "sl", 0x0425: "et", 0x0426: "lv",
	0x0427: "lt", 0x0429: "fa", 0x042A: "vi", 0x042D: "eu", 0x0439: "hi", 0x043E: "ms", 0x0456: "gl",
	0x0804: "zh-CN", 0x0809: "en-GB", 0x080A: "es-MX", 0x080C: "fr-BE", 0x0816: "pt-PT", 0x0C04: "zh-HK",
	0x0C07: "de-AT", 0x0C09: "en-AU", 0x0C0A: "es", 0x0C0C: "fr-CA", 0x1004: "zh-SG", 0x1009: "en-CA",
	0x1404: "zh-MO", 0x100C: "fr-CH", 0x0807: "de-CH", 0x0420: "ur", 0x0C01: "ar-EG", 0x0801: "ar-IQ",
	0x0445: "bn", 0x0449: "ta", 0x044A: "te", 0x044B: "kn", 0x044C: "ml", 0x0447: "gu", 0x0446: "pa",
	0x0448: "or", 0x044E: "mr", 0x045B: "si", 0x0454: "lo", 0x0453: "km", 0x0455: "my",
}

var macLangs = map[uint16]string{
	0: "en", 1: "fr", 2: "de", 3: "it", 4: "nl", 5: "sv", 6: "es", 7: "da", 8: "pt", 9: "no", 10: "he",
	11: "ja", 12: "ar", 13: "fi", 14: "el", 15: "is", 17: "tr", 18: "hr", 19: "zh-Hant", 21: "hi",
	22: "th", 23: "ko", 25: "pl", 26: "hu", 27: "et", 28: "lv", 30: "fo", 31: "fa", 32: "ru",
	33: "zh-Hans",
}

// name returns the best string of a name ID: English from Windows first,
// then any Windows or Unicode record, then Macintosh ones.
func (fc *face) name(id uint16) string {
	best, score := "", -1
	for _, r := range fc.names {
		if r.id != id {
			continue
		}
		s := r.rank
		if r.en {
			s += 10
		}
		if s > score {
			best, score = r.text, s
		}
	}
	return best
}

// localNames returns the strings of a name ID in languages other than the
// one name returns, each language once.
func (fc *face) localNames(id uint16) []nameRec {
	main := fc.name(id)
	var out []nameRec
	seen := map[string]bool{main: true}
	for _, r := range fc.names {
		if r.id == id && !r.en && !seen[r.text] {
			seen[r.text] = true
			out = append(out, r)
		}
	}
	return out
}

// family is the name the font is known by: its typographic family, or its
// family, or its full name.
func (fc *face) family() string {
	for _, id := range []uint16{16, 1, 4} {
		if s := fc.name(id); s != "" {
			return s
		}
	}
	return ""
}

// fullName is the font's full name, or its family and style.
func (fc *face) fullName() string {
	if s := fc.name(4); s != "" {
		return s
	}
	fam, st := fc.family(), fc.name(17)
	if st == "" {
		st = fc.name(2)
	}
	if st != "" && fam != "" {
		return fam + " " + st
	}
	return fam
}

// glyphNames returns the glyph names of a font: those of its post table,
// or of its CFF charset (cidNNNNN for a CID-keyed font).
func glyphNames(f *sfnt.Font) []string {
	names := make([]string, f.NumGlyphs)
	if len(f.PostNames) > 0 {
		for n, g := range f.PostNames {
			if int(g) < len(names) {
				names[g] = n
			}
		}
		return names
	}
	if b := f.Tables["CFF "]; b != nil {
		c, err := cff.Parse(b)
		if err != nil {
			return names
		}
		for g := range min(len(names), len(c.Charset)) {
			if c.IsCID {
				names[g] = "cid" + strconv.Itoa(c.Charset[g])
			} else {
				names[g] = c.SIDName(c.Charset[g])
			}
		}
	}
	return names
}
