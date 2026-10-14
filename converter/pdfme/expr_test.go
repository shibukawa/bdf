package pdfme

import (
	"strings"
	"testing"
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
	if _, err := evalExpr(n, map[string]any{"items": big}); err == nil || !strings.Contains(err.Error(), "too long") {
		t.Errorf("an expression of 1.6e9 steps: %v", err)
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
