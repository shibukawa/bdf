package pdf

import "testing"

func TestPDFDate(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"D:20260925072116+00'00'", "2026-09-25T07:21:16+00:00"},
		{"D:19990209153925-08'00", "1999-02-09T15:39:25-08:00"},
		{"D:20240501123456+0930", "2024-05-01T12:34:56+09:30"},
		{"D:20240501123456+09", "2024-05-01T12:34:56+09:00"},
		{"D:20240501123456Z00'00'", "2024-05-01T12:34:56Z"},
		{"D:20240501123456", "2024-05-01T12:34:56"},
		{"20240501", "2024-05-01"},
		{"D:2024", "2024"},
		{"D:202405", "2024-05"},
		{"D:2024050112", "2024-05-01T12:00"},
		{"D:202405011234", "2024-05-01T12:34"},
		{"D:20241301", ""},         // month 13
		{"D:202405011", ""},        // odd number of digits
		{"D:2024050112345+09", ""}, // odd number of digits
		{"D:20240501123456+9", ""},
		{"D:20240501123456 PST", ""},
		{"", ""},
		{"yesterday", ""},
	} {
		if got := pdfDate(c.in); got != c.want {
			t.Errorf("pdfDate(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestPDFDocText(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"plain", "plain"},
		{"caf\xe9 \x93 \x80 \xa0 \x84", "café ﬁ • € —"},
		{"UTF-8 without BOM: café", "UTF-8 without BOM: café"},
		{"\x9f", "�"},
	} {
		if got := pdfDocText([]byte(c.in)); got != c.want {
			t.Errorf("pdfDocText(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
