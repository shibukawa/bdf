package csv

import (
	"reflect"
	"strings"
	"testing"
)

func readAll(s string, d Dialect) [][]string {
	r := newReader(s, d)
	var out [][]string
	for {
		f, ok := r.record(nil)
		if !ok {
			return out
		}
		out = append(out, f)
	}
}

func TestReader(t *testing.T) {
	std := Dialect{Delimiter: ',', Quote: '"'}
	for _, c := range []struct {
		in   string
		d    Dialect
		want [][]string
	}{
		{"a,b\r\nc,d\r\n", std, [][]string{{"a", "b"}, {"c", "d"}}},
		{"a,b\rc,d", std, [][]string{{"a", "b"}, {"c", "d"}}},
		{"a,,\n,b", std, [][]string{{"a", "", ""}, {"", "b"}}},
		{"a\n\nb\n", std, [][]string{{"a"}, nil, {"b"}}},
		{`"a,b","c""d",e`, std, [][]string{{"a,b", `c"d`, "e"}}},
		{"\"two\r\nlines\",x\ny,z", std, [][]string{{"two\r\nlines", "x"}, {"y", "z"}}},
		{`a, "b, c" ,d`, std, [][]string{{"a", "b, c", "d"}}},
		{`5" screen,"say "hi" now",x`, std, [][]string{{`5" screen`, `say "hi" now`, "x"}}},
		{`"open,x` + "\ny", std, [][]string{{"open,x\ny"}}},
		{`"a\"b",c`, Dialect{Delimiter: ',', Quote: '"', Escape: '\\'}, [][]string{{`a"b`, "c"}}},
		{`'a;b';'it''s'`, Dialect{Delimiter: ';', Quote: '\''}, [][]string{{"a;b", "it's"}}},
		{"\"a\"\tb\n", Dialect{Delimiter: '\t'}, [][]string{{`"a"`, "b"}}},
	} {
		if got := readAll(c.in, c.d); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%q: %q, want %q", c.in, got, c.want)
		}
	}
}

func TestSniff(t *testing.T) {
	lines := func(ls ...string) string { return strings.Join(ls, "\n") + "\n" }
	for _, c := range []struct {
		name string
		in   string
		want Dialect
	}{
		{"csv", lines("a,b,c", "1,2,3", "4,5,6"), Dialect{',', '"', 0}},
		{"tsv", lines("a\tb\tc", "1\t2\t3"), Dialect{'\t', '"', 0}},
		{"semicolon, decimal commas", lines("a;b;c", "1,5;2;3", "4,25;5;6,5"), Dialect{';', '"', 0}},
		{"pipe", lines("a|b|c", "1|2|3"), Dialect{'|', '"', 0}},
		{"quoted delimiters", lines(`id,place,n`, `1,"Tokyo, Japan",3`, `2,"Paris, France",4`), Dialect{',', '"', 0}},
		{"tsv with commas", lines("name\tplace", "Ann\tTokyo, Japan", "Bob\tParis, France"), Dialect{'\t', '"', 0}},
		{"single quotes", lines(`'a','b,c'`, `'d','e,f'`, `'g','h'`), Dialect{',', '\'', 0}},
		{"apostrophes", lines("name,note", "Ann,it's fine", "Bob,that's it", "Cy,ok"), Dialect{',', '"', 0}},
		{"backslash escapes", lines(`1,"say \"hi\", then go",x`, `2,"plain",y`), Dialect{',', '"', '\\'}},
		{"unquoted tsv", lines("a\tb", "\"quoted\" start\tx", "y\t\"z"), Dialect{'\t', 0, 0}},
		{"one column", lines("a", "b", "c"), Dialect{',', '"', 0}},
	} {
		if got := sniff(c.in, false, 0, 0).d; got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
	// a sample cut in the middle of a quoted field: the record is left out
	s := "a,b\n1,\"x\n" + strings.Repeat("y", 10)
	if got := sniff(s, true, 0, 0); got.d != (Dialect{',', '"', 0}) || got.records != 1 || got.score != 1 {
		t.Errorf("cut sample: %+v", got)
	}
	// given options narrow the guess
	if got := sniff("a;b,c\n", false, ';', noQuote).d; got != (Dialect{';', 0, 0}) {
		t.Errorf("given: %q", got)
	}
}
