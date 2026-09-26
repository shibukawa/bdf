package csv

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// Values are classified the way Excel reads what it opens: numbers (with
// signs, thousands separators, decimal commas, exponents, percent signs,
// currency symbols and parentheses for negatives), dates and times become
// numbers, which Excel aligns to the right; the rest is text. Unlike
// Excel, a number with leading zeros (a code such as 007 or a postal
// code) stays text, and so do the other integers of a column with such
// codes. The classes also tell header rows from data.

type kind uint8

const (
	kindEmpty kind = iota
	kindText
	kindNumber
	kindDate // dates, times, or both
	kindBool
	kindEmail
	kindURL
)

// numeric reports whether a class is aligned as a number.
func (k kind) numeric() bool { return k == kindNumber || k == kindDate }

// classify returns the class of a value.
func classify(v string) kind {
	s := strings.TrimSpace(v)
	switch {
	case s == "":
		return kindEmpty
	case isNumber(s):
		return kindNumber
	case isDateTime(s):
		return kindDate
	}
	switch strings.ToLower(s) {
	case "true", "false":
		return kindBool
	}
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return kindURL
	}
	if at := strings.IndexByte(s, '@'); at > 0 && !strings.ContainsAny(s, " \t,;") && strings.Contains(s[at+1:], ".") {
		return kindEmail
	}
	return kindText
}

var currencies = []string{"$", "€", "£", "¥", "￥", "₩", "₹", "¢"}

// isNumber reports whether s (trimmed) is a number.
func isNumber(s string) bool {
	neg := false
	if len(s) > 2 && s[0] == '(' && s[len(s)-1] == ')' {
		s, neg = strings.TrimSpace(s[1:len(s)-1]), true
	}
	if s != "" && (s[0] == '+' || s[0] == '-') && !neg {
		s = strings.TrimSpace(s[1:])
	}
	for _, c := range currencies {
		if strings.HasPrefix(s, c) {
			s = strings.TrimSpace(s[len(c):])
			break
		}
	}
	if s != "" && (s[0] == '+' || s[0] == '-') && !neg {
		s = s[1:] // -$5 or $-5
	}
	if strings.HasSuffix(s, "%") {
		s = strings.TrimSpace(s[:len(s)-1])
	}
	mant, exp := s, ""
	if i := strings.IndexAny(s, "eE"); i > 0 {
		mant, exp = s[:i], s[i+1:]
		if exp != "" && (exp[0] == '+' || exp[0] == '-') {
			exp = exp[1:]
		}
		if !allDigits(exp) {
			return false
		}
	}
	return isDecimal(mant, exp == "")
}

func allDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// isDecimal reports whether s is digits with an optional fraction after
// '.' or ',' (1234, 0.5, .5, 3., 1,5), or, when grouped, with thousands
// separators (',', '.', an apostrophe, a space) between groups of three:
// 1,234,567.8 and 1.234.567,8. Either reading of 1,234 makes it a number.
// Integers with leading zeros are codes (007), not numbers.
func isDecimal(s string, grouped bool) bool {
	if s == "" {
		return false
	}
	var groups, seps []string
	start := 0
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r >= '0' && r <= '9' {
			i++
			continue
		}
		switch r {
		case '.', ',', '\'', ' ', ' ', ' ':
		default:
			return false
		}
		groups, seps = append(groups, s[start:i]), append(seps, s[i:i+size])
		i += size
		start = i
	}
	groups = append(groups, s[start:])
	code := func(g string) bool { return len(g) > 1 && g[0] == '0' }
	if len(seps) == 0 {
		return !code(groups[0])
	}
	if len(seps) == 1 && (seps[0] == "." || seps[0] == ",") {
		// a fraction
		return (groups[0] != "" || groups[1] != "") && !code(groups[0])
	}
	if !grouped {
		return false
	}
	k := len(seps) // thousands separators
	if last := seps[k-1]; last != seps[0] {
		if last != "." && last != "," || groups[k] == "" {
			return false
		}
		k-- // then a fraction
	}
	if len(groups[0]) == 0 || len(groups[0]) > 3 || groups[0][0] == '0' {
		return false
	}
	for i := 1; i <= k; i++ {
		if seps[i-1] != seps[0] || len(groups[i]) != 3 {
			return false
		}
	}
	return true
}

var months = []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}

// isDateTime reports whether s (trimmed) is a date, a time, or a date and
// a time.
func isDateTime(s string) bool {
	if isTime(s) || isDate(s) {
		return true
	}
	// a date and a time: 2024-01-02T03:04:05Z, 2024/1/2 3:04
	if i := strings.IndexAny(s, "T "); i > 0 {
		return isDate(s[:i]) && isTime(strings.TrimSpace(s[i+1:]))
	}
	return false
}

// isTime recognizes 9:30, 09:30:15, 9:30:15.250 with AM/PM or a time
// zone (Z, +09:00).
func isTime(s string) bool {
	u := strings.ToUpper(s)
	for _, suf := range []string{" AM", " PM", "AM", "PM"} {
		if strings.HasSuffix(u, suf) {
			s = s[:len(s)-len(suf)]
			break
		}
	}
	if strings.HasSuffix(s, "Z") {
		s = s[:len(s)-1]
	} else if i := strings.LastIndexAny(s, "+-"); i > 0 {
		if !isClock(s[i+1:], 2) && !(len(s[i+1:]) == 4 && allDigits(s[i+1:])) {
			return false
		}
		s = s[:i]
	}
	return isClock(s, 3)
}

// isClock recognizes h:mm, h:mm:ss and h:mm:ss.fff (at most parts parts).
func isClock(s string, parts int) bool {
	if i := strings.IndexByte(s, '.'); i > 0 {
		if !allDigits(s[i+1:]) {
			return false
		}
		s = s[:i]
	}
	f := strings.Split(s, ":")
	if len(f) < 2 || len(f) > parts {
		return false
	}
	for i, p := range f {
		if !allDigits(p) || len(p) > 2 || i > 0 && len(p) != 2 {
			return false
		}
		if n, _ := strconv.Atoi(p); i > 0 && n > 59 || i == 0 && n > 24 {
			return false
		}
	}
	return true
}

// isDate recognizes 2024-01-02, 2024/1/2, 2024.1.2, 1/2/2024, 2/1/24,
// 1/2, 2-Jan-2024, Jan 2, 2024, 2 Jan 2024, 2024年1月2日, 令和6年1月2日,
// 1月2日 and 2024年1月.
func isDate(s string) bool {
	if isJapaneseDate(s) {
		return true
	}
	if isNamedMonthDate(s) {
		return true
	}
	sep := strings.IndexAny(s, "-/.")
	if sep <= 0 {
		return false
	}
	f := strings.Split(s, s[sep:sep+1])
	for _, p := range f {
		if !allDigits(p) || len(p) > 4 {
			return false
		}
	}
	n := func(i int) int { v, _ := strconv.Atoi(f[i]); return v }
	switch len(f) {
	case 3:
		if len(f[0]) == 4 { // y-m-d
			return validDay(n(0), n(1), n(2))
		}
		if len(f[2]) == 4 || len(f[2]) == 2 && s[sep] == '/' { // m/d/y or d/m/y
			y := n(2)
			return len(f[0]) <= 2 && len(f[1]) <= 2 && (validDay(y, n(0), n(1)) || validDay(y, n(1), n(0)))
		}
	case 2:
		if s[sep] == '/' && len(f[0]) <= 2 && len(f[1]) <= 2 { // m/d
			return validDay(2000, n(0), n(1))
		}
		if len(f[0]) == 4 && s[sep] != '.' && len(f[1]) <= 2 { // y-m, y/m
			m := n(1)
			return m >= 1 && m <= 12
		}
	}
	return false
}

func validDay(y, m, d int) bool {
	if m < 1 || m > 12 || d < 1 {
		return false
	}
	days := [...]int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}[m-1]
	if m == 2 && !(y%4 == 0 && (y%100 != 0 || y%400 == 0)) {
		days = 28
	}
	return d <= days
}

// isNamedMonthDate recognizes 2-Jan-2024, 02-Jan-24, Jan-24, Jan 2, 2024
// and 2 January 2024.
func isNamedMonthDate(s string) bool {
	f := strings.FieldsFunc(strings.ToLower(s), func(r rune) bool { return r == ' ' || r == '-' || r == ',' || r == '/' })
	if len(f) < 2 || len(f) > 3 {
		return false
	}
	month, nums := -1, 0
	for _, p := range f {
		if allDigits(p) && len(p) <= 4 {
			nums++
			continue
		}
		for i, m := range months {
			if len(p) >= 3 && strings.HasPrefix(p, m) && month < 0 {
				month = i
			}
		}
		if month < 0 {
			return false
		}
	}
	return month >= 0 && nums == len(f)-1
}

var eras = []string{"明治", "大正", "昭和", "平成", "令和", "M", "T", "S", "H", "R"}

// isJapaneseDate recognizes 2024年1月2日, 令和6年1月2日, R6.1.2, 2024年1月
// and 1月2日.
func isJapaneseDate(s string) bool {
	for _, e := range eras {
		rest, ok := strings.CutPrefix(s, e)
		if !ok {
			continue
		}
		if r, ok := strings.CutPrefix(rest, "元年"); ok {
			rest = "1年" + r // the first year of an era
		}
		if rest == "" || rest[0] < '0' || rest[0] > '9' {
			continue
		}
		if strings.Contains(rest, ".") {
			f := strings.Split(rest, ".")
			return len(f) == 3 && allDigits(f[0]) && allDigits(f[1]) && allDigits(f[2])
		}
		return kanjiDate(rest, true)
	}
	return kanjiDate(s, false)
}

// kanjiDate recognizes y年m月d日, y年m月 and m月d日 (era years when era).
func kanjiDate(s string, era bool) bool {
	var y, m, d string
	rest := s
	if i := strings.Index(rest, "年"); i > 0 {
		y, rest = rest[:i], rest[i+len("年"):]
	}
	if i := strings.Index(rest, "月"); i > 0 {
		m, rest = rest[:i], rest[i+len("月"):]
	} else {
		return false
	}
	if i := strings.Index(rest, "日"); i > 0 {
		d, rest = rest[:i], rest[i+len("日"):]
	}
	if rest != "" || !allDigits(m) || y != "" && !allDigits(y) || d != "" && !allDigits(d) || y == "" && d == "" {
		return false
	}
	mm, _ := strconv.Atoi(m)
	if d == "" {
		return mm >= 1 && mm <= 12
	}
	yy, _ := strconv.Atoi(y)
	if y == "" || era {
		yy = 2000
	}
	dd, _ := strconv.Atoi(d)
	return validDay(yy, mm, dd)
}
