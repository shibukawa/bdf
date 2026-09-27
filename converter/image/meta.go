package image

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
)

// merge fills the elements of dc that have no value with those of src.
// Sources are merged in order of precedence (docs/spec.md §4.3).
func merge(dc *bdf.DublinCore, src bdf.DublinCore) {
	for _, name := range bdf.DCTerms {
		if f := dc.Field(name); len(*f) == 0 {
			*f = *src.Field(name)
		}
	}
}

// add appends the non-empty values to an element.
func add(f *bdf.DCValues, values ...string) {
	for _, v := range values {
		if v = clean(v); v != "" {
			*f = append(*f, v)
		}
	}
}

// clean trims spaces and the NULs that pad fixed-size fields.
func clean(s string) string {
	return strings.TrimSpace(strings.Trim(s, "\x00"))
}

// text decodes metadata text that should be ASCII (EXIF, IPTC) or is UTF-8
// in practice: bytes that are not valid UTF-8 are read as Latin-1. Text
// ends at the first NUL.
func text(b []byte) string {
	if i := strings.IndexByte(string(b), 0); i >= 0 {
		b = b[:i]
	}
	if utf8.Valid(b) {
		return clean(string(b))
	}
	return clean(latin1(b))
}

// latin1 decodes ISO 8859-1 (PNG tEXt and zTXt, GIF comments).
func latin1(b []byte) string {
	r := make([]rune, len(b))
	for i, c := range b {
		r[i] = rune(c)
	}
	return string(r)
}

// utf16Text decodes UTF-16 text in the given byte order, up to the first
// NUL (the Windows EXIF tags are UTF-16LE).
func utf16Text(b []byte, bigEndian bool) string {
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
	return clean(string(utf16.Decode(u)))
}

// splitList splits a list whose items are separated by semicolons (EXIF's
// Artist, the Windows author and keyword tags).
func splitList(s string) []string {
	var out []string
	for _, v := range strings.Split(s, ";") {
		if v = clean(v); v != "" {
			out = append(out, v)
		}
	}
	return out
}

// placeholder reports whether an image description is the text a camera
// writes when the user wrote none.
func placeholder(s string) bool {
	switch strings.ToUpper(strings.Join(strings.Fields(s), " ")) {
	case "", "OLYMPUS DIGITAL CAMERA", "SONY DSC", "DIGITAL CAMERA", "SAMSUNG DIGITAL CAMERA",
		"MINOLTA DIGITAL CAMERA", "KONICA MINOLTA DIGITAL CAMERA", "KODAK DIGITAL STILL CAMERA",
		"EXIF_JPEG_PICTURE", "DEFAULT":
		return true
	}
	return false
}

// w3cdtf writes a date in the W3C date and time format with the precision
// it has: n is the number of fields (1 year … 6 second) and zone is "",
// "Z" or "±hh:mm".
func w3cdtf(f [6]int, n int, zone string) string {
	s := fmt.Sprintf("%04d", f[0])
	for i, sep := range []string{"-", "-", "T", ":", ":"} {
		if i+1 >= n {
			break
		}
		s += sep + fmt.Sprintf("%02d", f[i+1])
	}
	if n == 4 {
		s += ":00" // the W3C format has no hour without minutes
	}
	if n >= 4 {
		s += zone
	}
	return s
}

// validDate checks the fields of a date: n fields, from the year on.
func validDate(f [6]int, n int) bool {
	lo := [6]int{1, 1, 1, 0, 0, 0}
	hi := [6]int{9999, 12, 31, 23, 59, 60}
	for i := 0; i < n; i++ {
		if f[i] < lo[i] || f[i] > hi[i] {
			return false
		}
	}
	return n > 0
}

// dateRE matches the dates of XMP (ISO 8601: 2026, 2026-09, 2026-09-01,
// 2026-09-01T10:20, …:30, …:30.25, with a zone), and EXIF-style dates that
// some tools write there (2026:09:01 10:20:30).
var dateRE = regexp.MustCompile(`^(\d{4})(?:[-:](\d{2})(?:[-:](\d{2})(?:[T ](\d{2}):(\d{2})(?::(\d{2})(\.\d+)?)?\s*(Z|[+-]\d{2}:?\d{2})?)?)?)?$`)

// isoDate normalizes an ISO 8601 date to the W3C date and time format;
// "" when it is not one.
func isoDate(s string) string {
	m := dateRE.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return ""
	}
	var f [6]int
	n := 0
	for i := 1; i <= 6 && m[i] != ""; i++ {
		f[i-1], _ = strconv.Atoi(m[i])
		n = i
	}
	if !validDate(f, n) {
		return ""
	}
	zone := m[8]
	if len(zone) == 5 { // +0900
		zone = zone[:3] + ":" + zone[3:]
	}
	out := w3cdtf(f, n, "")
	if n == 6 && m[7] != "" {
		out += m[7]
	}
	if n >= 4 {
		out += zone
	}
	return out
}

// exifDate converts an EXIF date and time ("2026:09:01 10:20:30", with
// spaces for unknown fields) and the offset of its OffsetTime tag
// ("+09:00") to the W3C format; "" when it is blank or malformed.
func exifDate(s, offset string) string {
	s = clean(s)
	if s == "" || strings.Trim(s, "0: ") == "" {
		return ""
	}
	d := isoDate(s)
	// a date with a time and no zone takes the offset
	if len(d) < 16 || strings.ContainsAny(d[10:], "Z+-") {
		return d
	}
	if o := clean(offset); len(o) == 6 && (o[0] == '+' || o[0] == '-') && o[3] == ':' && isoDate("2000-01-01T00:00"+o) != "" {
		d += o
	}
	return d
}

// textDate reads the free-form dates of PNG's "Creation Time" (RFC 1123 is
// recommended; ISO 8601 and EXIF-style dates are common).
func textDate(s string) string {
	s = clean(s)
	if d := isoDate(s); d != "" {
		return d
	}
	for _, layout := range []string{time.RFC1123Z, time.RFC1123, time.RFC822Z, time.RFC822, time.RFC850, time.ANSIC,
		"Mon, 2 Jan 2006 15:04:05 -0700", "2 Jan 2006 15:04:05 -0700", "Mon, 2 Jan 2006 15:04:05 MST"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.Format("2006-01-02T15:04:05Z07:00")
		}
	}
	return ""
}
