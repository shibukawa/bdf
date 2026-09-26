package xlsx

import (
	"strings"
	"testing"
)

// shown renders a formatted value as text, with gaps as "_" and the fill
// character as "*x".
func shown(f formatted) string {
	var b strings.Builder
	for _, p := range f.pieces {
		switch p.kind {
		case pieceText:
			b.WriteString(p.s)
		case pieceGap:
			b.WriteString("_")
		case pieceFill:
			b.WriteString("*" + p.s)
		}
	}
	return b.String()
}

func TestNumberFormats(t *testing.T) {
	for _, c := range []struct {
		code string
		v    float64
		want string
	}{
		{"General", 0, "0"},
		{"General", 1234567, "1234567"},
		{"General", 1.0 / 3, "0.333333333"},
		{"General", -2.0 / 3, "-0.666666667"},
		{"General", 123.456, "123.456"},
		{"General", 0.1 + 0.2, "0.3"},
		{"General", 123456789012, "1.23457E+11"},
		{"General", 0.00000000012, "1.2E-10"},
		{"0", 2.5, "3"},
		{"0", -2.5, "-3"},
		{"0.00", 3.14159, "3.14"},
		{"0.00", -0.001, "0.00"},
		{"#,##0", 1234567.8, "1,234,568"},
		{"#,##0.00", -1234.5, "-1,234.50"},
		{"#,##0,", 1234567, "1,235"},
		{"0.0,,\"M\"", 12345678, "12.3M"},
		{"0%", 0.256, "26%"},
		{"0.00%", 0.0123, "1.23%"},
		{"0.00E+00", 12345, "1.23E+04"},
		{"0.00E+00", 0.00012, "1.20E-04"},
		{"##0.0E+0", 12345, "12.3E+3"},
		{"# ?/?", 1.25, "1 1/4"},
		{"# ??/??", 3.14159, "3 14/99"},
		{"# ???/???", 3.14159, "3 _16/113"},
		{"?/8", 0.375, "3/8"},
		{"# ?/?", 2, "2 ___"},
		{"000-0000", 1234567, "123-4567"},
		{"00000", 42, "00042"},
		{"#.##", 0.5, ".5"},
		{"0.0#", 1.5, "1.5"},
		{"0.0?", 1.5, "1.5_"},
		{`#,##0;[Red]\-#,##0`, -1234, "-1,234"},
		{`#,##0_);(#,##0)`, -1234, "(1,234)"},
		{`#,##0_);(#,##0)`, 1234, "1,234_"},
		{`_(* #,##0_);_(* \(#,##0\);_(* "-"_);_(@_)`, 0, "_* -_"},
		{`_(* #,##0.00_);_(* \(#,##0.00\);_(* "-"??_);_(@_)`, -5, "_* (5.00)"},
		{`"$"#,##0.00`, 1234.5, "$1,234.50"},
		{`[$€-407]#,##0.00`, 5, "€5.00"},
		{`[>=100]"big";[<0]"neg";"small"`, 150, "big"},
		{`[>=100]"big";[<0]"neg";"small"`, -3, "neg"},
		{`[>=100]"big";[<0]"neg";"small"`, 5, "small"},
		{`0;-0;"zero"`, 0, "zero"},
		{`0;;`, -5, ""},
		{"@", 12, "12"},
		{"m/d/yyyy", 45000, "3/15/2023"},
		{"yyyy/m/d", 45000, "2023/3/15"},
		{"d-mmm-yy", 45000, "15-Mar-23"},
		{"dddd, mmmm dd", 45000, "Wednesday, March 15"},
		{"yyyy-mm-dd hh:mm:ss", 45000.5, "2023-03-15 12:00:00"},
		{"h:mm AM/PM", 0.75, "6:00 PM"},
		{"h:mm:ss", 0.999999, "0:00:00"},
		{"[h]:mm:ss", 1.5, "36:00:00"},
		{"mm:ss.0", 0.00001, "00:00.9"},
		{`yyyy"年"m"月"d"日"`, 45000, "2023年3月15日"},
		{`[$-411]ggge"年"m"月"d"日"`, 45000, "令和5年3月15日"},
		{`[$-411]ge.m.d`, 32874, "H2.1.1"},
		{"aaaa", 45000, "水曜日"},
		{"m/d/yyyy", 60, "2/29/1900"},
		{"m/d/yyyy", 61, "3/1/1900"},
		{"m/d/yyyy", 1, "1/1/1900"},
		{"m/d/yyyy", -1, "#"},
		{"[DBNum1]0", 123, "一二三"},
	} {
		nf := parseNumFormat(c.code, defaultPalette)
		if got := shown(nf.formatNumber(c.v, false)); got != c.want {
			t.Errorf("format(%v, %q) = %q, want %q", c.v, c.code, got, c.want)
		}
	}
	// date1904
	if got := shown(parseNumFormat("yyyy-mm-dd", nil).formatNumber(0, true)); got != "1904-01-01" {
		t.Errorf("1904 = %q", got)
	}
	// colors and text sections
	nf := parseNumFormat(`0;[Red]-0;0;"text: "@`, defaultPalette)
	if f := nf.formatNumber(-3, false); f.color == nil || *f.color != (rgb{1, 0, 0}) || shown(f) != "-3" {
		t.Errorf("red negative = %q %v", shown(f), f.color)
	}
	if got := shown(nf.formatText("abc")); got != "text: abc" {
		t.Errorf("text section = %q", got)
	}
	if got := shown(parseNumFormat("0.00", nil).formatText("abc")); got != "abc" {
		t.Errorf("text in number format = %q", got)
	}
	if nf := parseNumFormat("[Color10]0", defaultPalette); nf.formatNumber(1, false).color == nil {
		t.Error("[Color10] ignored")
	}
}

func TestGeneralNarrow(t *testing.T) {
	for _, c := range []struct {
		v    float64
		n    int
		want string
	}{
		{1.0 / 3, 5, "0.333"},
		{123456.789, 6, "123457"},
		{123456789, 6, "1E+08"},
		{123456789, 3, ""},
		{-1.5, 3, "-1.5"},
	} {
		if got := formatGeneral(c.v, c.n); got != c.want {
			t.Errorf("formatGeneral(%v, %d) = %q, want %q", c.v, c.n, got, c.want)
		}
	}
}

func TestTint(t *testing.T) {
	// accent1 of Office 2013+ (4472C4) lightened by 0.8 and darkened by 0.25
	c := applyTint(hexRGB(0x4472C4), 0.7999816888943144)
	if got := c.bdf().CSS(); got != "#dae3f3" {
		t.Errorf("tint 0.8 = %s", got)
	}
	c = applyTint(hexRGB(0x4472C4), -0.249977111117893)
	if got := c.bdf().CSS(); got != "#2f5597" {
		t.Errorf("tint -0.25 = %s", got)
	}
}
