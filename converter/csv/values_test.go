package csv

import (
	"strings"
	"testing"
)

func TestClassify(t *testing.T) {
	for want, vs := range map[kind][]string{
		kindNumber: {"0", "42", "-7", "+3", " 12 ", "3.14", ".5", "3.", "1,5", "1,234", "1,234,567.89", "1.234.567,89", "1 234 567",
			"1'234'567", "1e10", "6.02E+23", "12%", "-12.5 %", "$5", "-$5", "$-5", "$1,234.50", "¥1,000", "€ 3,50", "(123)", "(1,234.50)"},
		kindDate: {"2024-01-02", "2024/1/2", "2024.01.02", "1/2/2024", "31/12/2024", "12/31/24", "1/2", "2024-02", "2024/2",
			"2-Jan-2024", "02-Jan-24", "Jan-24", "Jan 2, 2024", "2 January 2024", "9:30", "09:30:15", "9:30:15.250", "9:30 PM",
			"2024-01-02T03:04:05Z", "2024-01-02T03:04:05+09:00", "2024-01-02 03:04", "2024/1/2 3:04:05",
			"2024年1月2日", "2024年1月", "1月2日", "令和6年1月2日", "令和元年5月1日", "R6.1.2", "H31.4.30"},
		kindText: {"007", "0012", "01.5", "1.2.3", "192.168.1.1", "03-1234-5678", "1,23,456", "2024-13-01", "2023-02-29",
			"25:00", "9:60", "abc", "12abc", "1-2", "$", "%", "-", "e5", "Jan", "3 apples", "13/13/2024", "2024年13月1日", "No. 5"},
		kindBool:  {"TRUE", "false", "True"},
		kindEmail: {"a@example.com", "first.last@mail.co.jp"},
		kindURL:   {"https://example.com/x", "http://a.b"},
		kindEmpty: {"", "   "},
	} {
		for _, v := range vs {
			if got := classify(v); got != want {
				t.Errorf("%q: kind %d, want %d", v, got, want)
			}
		}
	}
}

func TestHeader(t *testing.T) {
	recs := func(s string) [][]string {
		var out [][]string
		for _, l := range strings.Split(strings.TrimSpace(s), "\n") {
			out = append(out, strings.Split(l, ","))
		}
		return out
	}
	for _, c := range []struct {
		name string
		in   string
		want bool
	}{
		{"labels over numbers", "id,name,price\n1,Ann,3.5\n2,Bob,4\n3,Cy,5", true},
		{"numbers only", "1,2,3\n4,5,6\n7,8,9", false},
		{"data first", "1,Ann,3.5\n2,Bob,4\n3,Cy,5", false},
		{"labels over dates and mail", "when,who\n2024-01-02,a@example.com\n2024-01-03,b@example.com", true},
		{"years over numbers", "country,2022,2023,2024\nJapan,1.2,1.3,1.5\nFrance,2.1,2.0,2.2", true},
		{"years as data", "2022,Japan,1.2\n2023,Japan,1.3\n2024,France,2.2", false},
		{"text columns", "name,city\nAnn,Tokyo\nBob,Osaka\nCy,Nagoya", true},
		{"text data with repeats", "Ann,Tokyo\nBob,Osaka\nCy,Tokyo\nDan,Osaka", false},
		{"codes of one length", "name,code\nAnn,AB12\nBob,CD34\nCy,EF56", true},
		{"cross table corner", ",Jan,Feb\nA,1,2\nB,3,4", true},
		{"repeated values", "a,a,b\nc,d,e\nf,g,h", false},
		{"one record", "id,name", false},
		{"booleans", "active,name\ntrue,Ann\nfalse,Bob\ntrue,Cy", true},
		{"japanese", "氏名,年齢\n山田,30\n佐藤,25", true},
	} {
		if got := hasHeader(recs(c.in)); got != c.want {
			t.Errorf("%s: header %v, want %v", c.name, got, c.want)
		}
	}
}
