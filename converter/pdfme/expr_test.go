package pdfme

import (
	"strings"
	"testing"
	"time"
)

// TestExpressions evaluates the expressions that the text of a read-only
// field holds, as pdfme does: with the values of the data, the operators
// and the functions of its list.
func TestExpressions(t *testing.T) {
	vars := map[string]any{
		"name": "ACME", "qty": 3.0, "price": 1200.0, "rate": "10", "paid": true, "empty": "",
		"items":       []any{[]any{"Pen", "2", "150"}, []any{"Ink", "1", "800"}},
		"customer":    map[string]any{"name": "Taro", "tags": []any{"a", "b"}},
		"currentPage": 2.0, "totalPages": 5.0, "date": "2026/10/14",
	}
	for text, want := range map[string]string{
		"plain":                                       "plain",
		"{name}":                                      "ACME",
		"Page {currentPage} / {totalPages}":           "Page 2 / 5",
		"{qty * price}":                               "3600",
		"{Number(rate) / 100 * price}":                "120",
		"{0.1 + 0.2} {price * 1.1}":                   "0.30000000000000004 1320",
		"{(price * 1.1).toFixed(0)}":                  "1320",
		"{Math.round(price * 1.08)}":                  "1296",
		"{paid ? 'PAID' : 'DUE'}":                     "PAID",
		"{empty || '-'}":                              "-",
		"{name + ' Inc.'}":                            "ACME Inc.",
		"{name.toLowerCase().length}":                 "4",
		"{customer.name} ({customer.tags.join('/')})": "Taro (a/b)",
		"{items.length}":                              "2",
		"{items.reduce((sum, r) => sum + Number(r[1]) * parseFloat(r[2]), 0)}": "1100",
		"{items.map(r => r[0]).join(', ')}":                                    "Pen, Ink",
		"{items.filter(r => Number(r[2]) > 200).length}":                       "1",
		"{String(qty).padStart(4, '0')}":                                       "0003",
		"{2 ** 3 ** 2}":                                                        "512",
		"{-qty + 10 % 4}":                                                      "-1",
		"{qty >= 3 && price < 2000}":                                           "true",
		"{'5' == 5} {'5' === 5} {null == undefined}":                           "true false true",
		"{[1, 2, 3].includes(2)} {({a: 1}).a}":                                 "true 1",
		"{date}":                                                               "2026/10/14",
		"a {  qty  } b {price}c":                                               "a 3 b 1200c",
		// what is not evaluated stays as it is
		"{unknown}":                 "{unknown}",
		"{qty +}":                   "{qty +}",
		"{new Date()}":              "{new Date()}",
		"{name.constructor}":        "{name.constructor}",
		"{items.map(}":              "{items.map(}",
		"an open { brace":           "an open { brace",
		"{typeof qty}":              "{typeof qty}",
		"{name.toLocaleString()}":   "{name.toLocaleString()}",
		"{customer.missing.deeper}": "{customer.missing.deeper}",
	} {
		got, failed := expand(text, vars)
		if got != want {
			t.Errorf("%s = %q, want %q", text, got, want)
		}
		if failed != (strings.Contains(want, "{") && want != "an open { brace") {
			t.Errorf("%s: failed = %v", text, failed)
		}
	}
	// an expression that does not end, and one that is too deep, fail
	if _, err := parseExpr(strings.Repeat("(", 100) + "1" + strings.Repeat(")", 100)); err == nil {
		t.Error("an expression of a hundred parentheses is parsed")
	}
	n, err := parseExpr("items.map(a => items.map(b => items.map(c => items.map(d => 1))))")
	if err != nil {
		t.Fatal(err)
	}
	big := make([]any, 200)
	if _, err := evalExpr(n, map[string]any{"items": big}); err == nil || !strings.Contains(err.Error(), "too long") && !strings.Contains(err.Error(), "too many items") {
		t.Errorf("an expression of 1.6e9 steps: %v", err)
	}
}

// TestExpressionLimits refuses the expressions that would take too much
// room, each quickly: a list that holds itself written out as many times
// as it is nested, a list joined with itself over and over, more items
// than a list holds, and lists nested too deep.
func TestExpressionLimits(t *testing.T) {
	nest := func(n int, inner string) string {
		return "(f => " + strings.Repeat("f(", n) + inner + strings.Repeat(")", n) + ")"
	}
	vars := map[string]any{"a": "x", "l": []any{1.0}, "big": strings.Repeat("y", 1<<19+8), "items": make([]any, 1000)}
	for _, c := range []struct{ expr, want string }{
		{nest(30, "a") + "(y => [y, y])", "too long"},                      // 2^30 times x
		{"JSON.stringify(" + nest(30, "a") + "(y => [y, y]))", "too long"}, // and as JSON
		{nest(30, "a") + "(y => [y, y]).join('')", "too long"},             // and joined
		{"String(" + nest(30, "a") + "(y => [y, y]))", "too long"},
		{nest(40, "l") + "(y => y.concat(y))", "too many items"}, // 2^40 items
		{"big.split('').map(c => c + c)", "too many items"},      // half a million items, twice
		{"big + big + big", "too long"},
		{"items.reduce((acc, x) => [acc], [])", "nested too deep"}, // a thousand lists in each other
		{"JSON.stringify(items.reduce((acc, x) => [acc], []))", "nested too deep"},
		{"JSON.parse('[' + '1,'.repeat(400000) + '1]').concat(JSON.parse('[' + '1,'.repeat(400000) + '1]'))", "too many items"},
	} {
		start := time.Now()
		n, err := parseExpr(c.expr)
		if err != nil {
			t.Fatalf("%s: %v", c.expr, err)
		}
		v, err := evalExpr(n, vars)
		if err == nil {
			_, err = exprText(v)
		}
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v, want %q", c.expr, err, c.want)
		}
		if d := time.Since(start); d > 2*time.Second {
			t.Errorf("%s took %v", c.expr, d)
		}
	}
	// in a text, such an expression is left as it is, and the others are evaluated
	if got, failed := expand("{"+nest(30, "a")+"(y => [y, y])} {a}", vars); got != "{"+nest(30, "a")+"(y => [y, y])} x" || !failed {
		t.Errorf("expand: %q, %v", got, failed)
	}
	// what is within the limits is as JavaScript has it
	for expr, want := range map[string]string{
		"[1, [2, 3], [[4]]].join('/')":                                   "1/2,3/4",
		"String([1, [2, [3]]]) + [4]":                                    "1,2,34",
		"JSON.stringify({b: [1, 'x', null, true], a: 1.5})":              `{"a":1.5,"b":[1,"x",null,true]}`,
		"JSON.stringify('<a>') + JSON.stringify(0/0)":                    `"<a>"null`,
		"Object.entries({a: 1})[0].join('=')":                            "a=1",
		"Object.keys({b: 1, a: 2}).concat(Object.values({c: 3})).join()": "a,b,3",
		"JSON.parse('[1,[2]]').length":                                   "2",
		"'ab'.padStart(5, 'xy') + 'ab'.repeat(2)":                        "xyxababab",
		"'héllo'.slice(1, 3) + 'héllo'.charAt(1)":                        "élé",
	} {
		if got, failed := expand("{"+expr+"}", vars); got != want || failed {
			t.Errorf("%s = %q (failed %v), want %q", expr, got, failed, want)
		}
	}
}

// FuzzExpr evaluates any text as an expression: it ends, with a value or
// an error.
func FuzzExpr(f *testing.F) {
	for _, s := range []string{"a + b * 2", "items.reduce((s, r) => s + r[1], 0)", "x ? 'a' : \"b\"", "{a: [1, 2].map(v => v ** 2)}", "Math.max(1, 2).toFixed(2)", "!a && -b || +c"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		n, err := parseExpr(s)
		if err != nil {
			return
		}
		evalExpr(n, map[string]any{"a": 1.0, "b": "2", "c": nil, "x": true, "items": []any{[]any{"p", 2.0}, []any{"q", 3.0}}})
	})
}
