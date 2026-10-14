package pdfme

import (
	"math"
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

func TestExpressionResourceLimits(t *testing.T) {
	vars := map[string]any{"items": make([]any, 24), "large": strings.Repeat("x", maxExprText/2)}
	for name, src := range map[string]string{
		"replace":           "'x'.repeat(4096).replaceAll('x', 'y'.repeat(4096))",
		"empty replacement": "'x'.repeat(4096).replaceAll('', 'y'.repeat(4096))",
		"join":              "[large, large].join(large)",
		"split":             "'x'.repeat(65537).split('')",
		"concat":            "items.reduce(a => a.concat(a), [0])",
		"shared JSON":       "JSON.stringify(items.reduce(a => [a, a], [0]))",
		"shared text":       "String(items.reduce(a => [a, a], [0]))",
		"recursion":         "(f => f(f))(f => f(f))",
		"temporary texts":   "items.map(i => 'x'.repeat(1048576))",
	} {
		t.Run(name, func(t *testing.T) {
			n, err := parseExpr(src)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := evalExpr(n, vars); err == nil {
				t.Fatal("an expression exceeding its resource limit succeeded")
			}
			text := "{" + src + "}"
			if got, failed := expand(text, vars); !failed || got != text {
				t.Fatalf("failed expression was not left intact: failed = %v, length = %d", failed, len(got))
			}
		})
	}
	text := "{'x'.repeat(600000)}{'x'.repeat(600000)}"
	if got, failed := expand(text, nil); !failed || got != text {
		t.Fatal("the combined text of expressions exceeded its limit")
	}
	// A nested value returned directly is checked before text conversion,
	// including cyclic values supplied by a Go caller.
	cycle := make([]any, 1)
	cycle[0] = cycle
	if got, failed := expand("{cycle}", map[string]any{"cycle": cycle}); !failed || got != "{cycle}" {
		t.Fatal("a cyclic value was expanded")
	}
}

func TestParseFloatPrefix(t *testing.T) {
	for src, want := range map[string]float64{
		"12.5mm": 12.5, "+.5e2tail": 50, "1e+tail": 1, "1.2.3": 1.2,
		"-Infinityrest": math.Inf(-1), "1e999": math.Inf(1), "text": math.NaN(),
		strings.Repeat("0", maxExprText-1) + "1": 1,
	} {
		got := parseFloatPrefix(src)
		if got != want && !(math.IsNaN(got) && math.IsNaN(want)) {
			t.Errorf("parseFloat of a text of %d bytes = %v, want %v", len(src), got, want)
		}
	}
}

func TestExpressionStringScanBudget(t *testing.T) {
	vars := map[string]any{
		"rows":   make([]any, maxExprItems),
		"s":      strings.Repeat("\U0001F600", maxExprText/8),
		"blank":  strings.Repeat(" ", maxExprText/2),
		"number": strings.Repeat("0", maxExprText/2-1) + "1",
		"json":   strings.Repeat(" ", maxExprText/2-1) + "1",
	}
	for _, src := range []string{
		"rows.map(() => s.charAt(500000)).length",
		"rows.map(() => s[500000]).length",
		"rows.map(() => s.includes('missing')).length",
		"rows.map(() => s.indexOf('missing')).length",
		"rows.map(() => s.length).length",
		"rows.map(() => s.slice(1073741824)).length",
		"rows.map(() => blank.trim()).length",
		"rows.map(() => blank.replaceAll(' ', '')).length",
		"rows.map(() => parseFloat(number)).length",
		"rows.map(() => Number(number)).length",
		"rows.map(() => Number([number])).length",
		"rows.map(() => String(s)).length",
		"rows.map(() => JSON.parse(json)).length",
		"rows.map(() => rows.includes('missing')).length",
		"rows.map(() => rows.indexOf('missing')).length",
	} {
		n, err := parseExpr(src)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := evalExpr(n, vars); err == nil || !strings.Contains(err.Error(), "too long") {
			t.Errorf("a loop scanning large strings escaped the work budget: %s: %v", src, err)
		}
	}
	// An index beyond even the byte length needs no scan. Repeating this
	// empty result stays cheap and succeeds within the expression budget.
	src := "rows.map(() => s.charAt(1073741824)).length"
	n, err := parseExpr(src)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := evalExpr(n, vars); err != nil || got != float64(maxExprItems) {
		t.Fatalf("a cheap out-of-bounds lookup failed: %v, %v", got, err)
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
