package image

import (
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/xmp"
)

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

// isoDate normalizes an ISO 8601 date to the W3C date and time format;
// "" when it is not one.
func isoDate(s string) string { return xmp.Date(s) }

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
