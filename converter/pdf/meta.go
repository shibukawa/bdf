package pdf

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
)

// info reads the document information dictionary as Dublin Core, mapped the
// way XMP maps it: Author is the creator, Subject the description, Keywords
// the subjects, and CreationDate and ModDate the created and modified dates.
func (p *pdf) info() bdf.DublinCore {
	var dc bdf.DublinCore
	if p.ctx.Info == nil {
		return dc
	}
	info := p.dict(*p.ctx.Info)
	if info == nil {
		return dc
	}
	for _, e := range []struct {
		f   *bdf.DCValues
		key string
	}{{&dc.Title, "Title"}, {&dc.Creator, "Author"}, {&dc.Description, "Subject"}} {
		if s := strings.TrimSpace(p.text(info[e.key])); s != "" {
			*e.f = bdf.DCValues{s}
		}
	}
	dc.Subject = bdf.SplitKeywords(p.text(info["Keywords"]))
	if d := pdfDate(p.text(info["CreationDate"])); d != "" {
		dc.Created = bdf.DCValues{d}
	}
	if d := pdfDate(p.text(info["ModDate"])); d != "" {
		dc.Modified = bdf.DCValues{d}
	}
	return dc
}

// pdfDate converts a PDF date (D:YYYYMMDDHHmmSSOHH'mm, PDF 32000-1 §7.9.4)
// to the W3C date and time format, keeping its precision. It returns "" for
// a malformed date.
func pdfDate(s string) string {
	s = strings.TrimPrefix(strings.TrimSpace(s), "D:")
	n := 0
	for n < len(s) && n < 14 && s[n] >= '0' && s[n] <= '9' {
		n++
	}
	if n < 4 || n%2 != 0 {
		return ""
	}
	num, zone := s[:n], s[n:]
	field := func(i, lo, hi int) (string, bool) {
		v, _ := strconv.Atoi(num[i : i+2])
		return num[i : i+2], v >= lo && v <= hi
	}
	var b strings.Builder
	b.WriteString(num[:4])
	// month, day, hour, minute, second
	sep := []string{"-", "-", "T", ":", ":"}
	lo := []int{1, 1, 0, 0, 0}
	hi := []int{12, 31, 23, 59, 59}
	for i := 0; 4+2*i < n; i++ {
		f, ok := field(4+2*i, lo[i], hi[i])
		if !ok {
			return ""
		}
		b.WriteString(sep[i] + f)
	}
	if n < 10 {
		return b.String()
	}
	if n == 10 {
		b.WriteString(":00") // the W3C format has no hour without minutes
	}
	switch {
	case zone == "":
	case zone[0] == 'Z':
		b.WriteByte('Z')
	case zone[0] == '+' || zone[0] == '-':
		z := strings.ReplaceAll(zone[1:], "'", "")
		if len(z) != 2 && len(z) != 4 || strings.Trim(z, "0123456789") != "" {
			return ""
		}
		mm := "00"
		if len(z) == 4 {
			mm = z[2:]
		}
		b.WriteString(zone[:1] + z[:2] + ":" + mm)
	default:
		return ""
	}
	return b.String()
}

// pdfDocEncoding maps the PDFDocEncoding codes that differ from Latin-1
// (PDF 32000-1 Annex D.2); 0 marks an undefined code.
var pdfDocEncoding = map[byte]rune{
	0x18: '˘', 0x19: 'ˇ', 0x1a: 'ˆ', 0x1b: '˙', 0x1c: '˝', 0x1d: '˛', 0x1e: '˚', 0x1f: '˜',
	0x7f: 0, 0x80: '•', 0x81: '†', 0x82: '‡', 0x83: '…', 0x84: '—', 0x85: '–', 0x86: 'ƒ', 0x87: '⁄',
	0x88: '‹', 0x89: '›', 0x8a: '−', 0x8b: '‰', 0x8c: '„', 0x8d: '“', 0x8e: '”', 0x8f: '‘',
	0x90: '’', 0x91: '‚', 0x92: '™', 0x93: 'ﬁ', 0x94: 'ﬂ', 0x95: 'Ł', 0x96: 'Œ', 0x97: 'Š',
	0x98: 'Ÿ', 0x99: 'Ž', 0x9a: 'ı', 0x9b: 'ł', 0x9c: 'œ', 0x9d: 'š', 0x9e: 'ž', 0x9f: 0, 0xa0: '€', 0xad: 0,
}

// pdfDocText decodes a PDFDocEncoding string. Bytes that are valid UTF-8 are
// kept as they are: PDFDocEncoding agrees with ASCII, and producers also
// write UTF-8 without a byte order mark.
func pdfDocText(b []byte) string {
	if utf8.Valid(b) {
		return string(b)
	}
	out := make([]rune, 0, len(b))
	for _, c := range b {
		r, ok := pdfDocEncoding[c]
		switch {
		case !ok:
			r = rune(c)
		case r == 0:
			r = utf8.RuneError
		}
		out = append(out, r)
	}
	return string(out)
}
