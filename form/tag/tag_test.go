package tag

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	type tok struct {
		kind   Kind
		text   string // of Text, or the name
		format string
		opt    string // mode or fit
	}
	for in, want := range map[string][]tok{
		"":                                 nil,
		"plain":                            {{Text, "plain", "", ""}},
		"{{name}}":                         {{Value, "name", "", ""}},
		"a {{ b . c }} d":                  {{Text, "a ", "", ""}, {Value, "b.c", "", ""}, {Text, " d", "", ""}},
		"{{.}}":                            {{Value, ".", "", ""}},
		"{{#items}}x{{/items}}":            {{Open, "items", "", ""}, {Text, "x", "", ""}, {Close, "items", "", ""}},
		"{{^paid}}{{/}}":                   {{Inverted, "paid", "", ""}, {Close, "", "", ""}},
		"{{%photo}}{{% logo | fit:cover}}": {{Image, "photo", "", ""}, {Image, "logo", "", "cover"}},
		`{{total | fmt:"#,##0"}}`:          {{Value, "total", "#,##0", ""}},
		`{{d|fmt:'yyyy"年"m"月"d"日"'}}`:      {{Value, "d", `yyyy"年"m"月"d"日"`, ""}},
		"{{total | fmt:“#,##0”}}":          {{Value, "total", "#,##0", ""}}, // quotes a word processor curled
		`{{x | fmt:"}}"}}`:                 {{Value, "x", "}}", ""}},
		`{{x | fmt:"a|b"}}`:                {{Value, "x", "a|b", ""}},
		`{{x | fmt:""}}`:                   {{Value, "x", "", ""}},
		"{{#rows | slots}}{{#p|pages}}":    {{Open, "rows", "", "slots"}, {Open, "p", "", "pages"}},
		"{{#g | block}}{{#r|ROWS}}":        {{Open, "g", "", "rows"}, {Open, "r", "", "rows"}},
		"｛｛顧客名｝｝様":                         {{Value, "顧客名", "", ""}, {Text, "様", "", ""}},
		"｛｛＃明細｜slots｝｝{{／明細}}": {{Open, "明細", "", "slots"}, {Close, "明細", "", ""}},
		"{{金額 ｜ fmt：「#,##0」}}": {{Value, "金額", "#,##0", ""}},
		"{{　name　}}":           {{Value, "name", "", ""}},
		`\{{not a tag}} {{a}}`: {{Text, "{{not a tag}} ", "", ""}, {Value, "a", "", ""}},
		"a } b { c }} d":       {{Text, "a } b { c }} d", "", ""}},
		"{{unit price}}":       {{Value, "unit price", "", ""}},
	} {
		got, err := Parse(in)
		if err != nil {
			t.Errorf("Parse(%q): %v", in, err)
			continue
		}
		var flat []tok
		for _, g := range got {
			x := tok{kind: g.Kind, text: g.Name, format: g.Format, opt: g.Mode + g.Fit}
			if g.Kind == Text {
				x.text = g.Text
			}
			flat = append(flat, x)
			if in[g.Pos:g.End] == "" {
				t.Errorf("Parse(%q): an empty token at %d", in, g.Pos)
			}
		}
		if !reflect.DeepEqual(flat, want) {
			t.Errorf("Parse(%q) = %+v, want %+v", in, flat, want)
		}
	}
	if toks, _ := Parse(`{{x | fmt:""}}`); !toks[0].HasFormat {
		t.Error("an empty format is not kept")
	}
}

func TestParseErrors(t *testing.T) {
	for in, want := range map[string]string{
		"a {{name":           "not closed",
		"{{a {{b}}":          "not closed",
		"{{}}":               "without a name",
		"{{#}}":              "without a name",
		"{{a..b}}":           "empty part",
		"{{a | upper}}":      "option",
		"{{a | fmt}}":        "option",
		"{{#a | fmt:\"0\"}}": "option",
		"{{%a | fit:tile}}":  "option",
		"{{a | pages}}":      "option",
		"{{a#b}}":            "character",
		"{{" + strings.Repeat("a", MaxTagLen+1) + "}}": "longer",
	} {
		_, err := Parse(in)
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("Parse(%.30q) = %v, want an error with %q", in, err, want)
		}
		if e, ok := err.(*Error); ok && (e.Pos < 0 || e.Pos >= len(in)) {
			t.Errorf("Parse(%.30q): the error is at %d", in, e.Pos)
		}
	}
	if !Has("x {{a}}") || !Has("｛｛a｝｝") || Has("{ {a} }") {
		t.Error("Has")
	}
}

func TestScope(t *testing.T) {
	var data any
	if err := json.Unmarshal([]byte(`{"name":"root","customer":{"name":"ACME","zip":null},"items":[{"name":"a","qty":2},{"qty":0}],"n":3}`), &data); err != nil {
		t.Fatal(err)
	}
	root := NewScope(Normalize(data))
	get := func(s *Scope, name string) any {
		v, ok := s.Lookup(strings.Split(name, "."))
		if name == "." {
			v, ok = s.Lookup([]string{"."})
		}
		if !ok {
			return "<none>"
		}
		return v
	}
	items := Items(get(root, "items"))
	first := root.Item(items[0], map[string]any{"_index": 1.0, "_first": true})
	second := first.Item(items[1], map[string]any{"_index": 2.0})
	blank := root.Item(items[0], map[string]any{"_page": 7.0}).Blank(map[string]any{"_blank": true})
	for _, c := range []struct {
		s    *Scope
		name string
		want any
	}{
		{root, "name", "root"}, {root, "customer.name", "ACME"}, {root, "customer.zip", nil}, {root, "customer.none", "<none>"},
		{root, "name.x", "<none>"}, {root, "none", "<none>"}, {root, "n", 3.0},
		{first, "name", "a"}, {first, "qty", 2.0}, {first, "n", 3.0}, {first, "_index", 1.0}, {first, "_first", true},
		{second, "name", "a"}, {second, "qty", 0.0}, {second, "_index", 2.0}, {second, "_first", true},
		{second, ".", items[1]},
		// from an empty slot only the variables are found
		{blank, "name", "<none>"}, {blank, "n", "<none>"}, {blank, "_blank", true}, {blank, "_page", 7.0}, {blank, ".", "<none>"},
	} {
		if got := get(c.s, c.name); !reflect.DeepEqual(got, c.want) {
			t.Errorf("Lookup(%q) = %v, want %v", c.name, got, c.want)
		}
	}
	for v, want := range map[any]bool{nil: false, false: false, true: true, 0.0: false, 0.5: true, "": false, "x": true} {
		if Truthy(v) != want {
			t.Errorf("Truthy(%v) = %v", v, !want)
		}
	}
	if Truthy([]any{}) || !Truthy([]any{nil}) || !Truthy(map[string]any{}) || !Truthy(time.Time{}) {
		t.Error("Truthy of lists, objects and dates")
	}
	if len(Items(nil)) != 0 || len(Items(false)) != 0 || len(Items([]any{1.0, 2.0})) != 2 || len(Items(map[string]any{})) != 1 || len(Items("x")) != 1 {
		t.Error("Items")
	}
}

func TestNormalize(t *testing.T) {
	in := map[any]any{"a": 1, "b": []map[string]any{{"c": int64(2)}}, 3: []string{"x"}, "d": json.Number("1.5"), "e": uint8(4), "f": (*int)(nil)}
	want := map[string]any{"a": 1.0, "b": []any{map[string]any{"c": 2.0}}, "3": []any{"x"}, "d": 1.5, "e": 4.0, "f": nil}
	if got := Normalize(in); !reflect.DeepEqual(got, want) {
		t.Errorf("Normalize = %#v", got)
	}
}

func TestParseDate(t *testing.T) {
	for in, want := range map[string]string{
		"2026-10-13":                "2026-10-13T00:00:00Z",
		"2026/1/2":                  "2026-01-02T00:00:00Z",
		"2026-10-13T09:30:00":       "2026-10-13T09:30:00Z",
		"2026-10-13T09:30:00+09:00": "2026-10-13T09:30:00Z", // the time on the wall
		"2026-10-13 09:30":          "2026-10-13T09:30:00Z",
		"20261013":                  "",
		"October 13":                "",
		"2026-13-40":                "",
		"12345":                     "",
	} {
		got, ok := ParseDate(in)
		if s := got.Format(time.RFC3339); ok != (want != "") || ok && s != want {
			t.Errorf("ParseDate(%q) = %s, %v; want %q", in, s, ok, want)
		}
	}
}

// TestCodes reads the options of picture tags that make QR codes and
// barcodes of their values.
func TestCodes(t *testing.T) {
	for in, want := range map[string]Token{
		"{{%url | qr}}":                     {Code: "qr"},
		"{{%url | qr:h | fit:contain}}":     {Code: "qr", Level: "H", Fit: FitContain},
		"{{％url ｜ ｑｒ}}":                     {},
		"{{%jan | barcode:JAN-13}}":         {Code: "ean13"},
		"{{%lot | barcode:\"gs1-128\"}}":    {Code: "gs1-128"},
		"{{%url | qr | quiet:0}}":           {Code: "qr", HasQuiet: true},
		"{{%code | barcode:nw7 | quiet:5}}": {Code: "nw7", Quiet: 5, HasQuiet: true},
	} {
		toks, err := Parse(in)
		if in == "{{％url ｜ ｑｒ}}" {
			// (the names of options are not read in full-width letters)
			if err == nil {
				t.Errorf("%s: no error", in)
			}
			continue
		}
		if err != nil || len(toks) != 1 || toks[0].Kind != Image {
			t.Errorf("%s: %v, %v", in, toks, err)
			continue
		}
		got := toks[0]
		if got.Code != want.Code || got.Level != want.Level || got.Fit != want.Fit || got.Quiet != want.Quiet || got.HasQuiet != want.HasQuiet {
			t.Errorf("%s: %+v, want %+v", in, got, want)
		}
	}
	for _, in := range []string{"{{url | qr}}", "{{%url | qr:Z}}", "{{%x | barcode}}", "{{%x | barcode:pdf417}}", "{{%x | quiet:65}}", "{{%x | quiet:-1}}", "{{#x | qr}}"} {
		if _, err := Parse(in); err == nil || !strings.Contains(err.Error(), "an option it does not take") {
			t.Errorf("%s: %v", in, err)
		}
	}
	// the symbols of values: a number is written as its digits
	toks, _ := Parse("{{%jan | barcode:jan}}{{%url | qr:Q | quiet:0}}")
	if s, err := toks[0].Symbol(4901234567894.0); err != nil || s.Text != "4901234567894" || !s.Linear {
		t.Errorf("a number: %+v, %v", s, err)
	}
	if s, err := toks[1].Symbol("https://example.com/"); err != nil || s.W != 25 || s.Linear {
		t.Errorf("a QR code without a quiet zone: %+v, %v", s, err)
	}
	if _, err := toks[0].Symbol("49012"); err == nil || !strings.Contains(err.Error(), "takes 13 digits") {
		t.Errorf("a short code: %v", err)
	}
	if _, err := toks[0].Symbol(true); err == nil {
		t.Error("a value that is no text")
	}
}

func TestParseText(t *testing.T) {
	type widths struct {
		width           int
		align, overflow string
	}
	for in, want := range map[string]widths{
		"{{a | width:20}}":                     {20, AlignLeft, ""},
		"{{a | right:8}}":                      {8, AlignRight, ""},
		"{{a|center：4|wrap}}":                  {4, AlignCenter, OverflowWrap},
		`{{a | fmt:"#,##0" | right:10 | cut}}`: {10, AlignRight, OverflowCut},
		"{{a}}":                                {0, "", ""},
	} {
		toks, err := ParseText(in)
		if err != nil {
			t.Errorf("ParseText(%q): %v", in, err)
			continue
		}
		if got := (widths{toks[0].Width, toks[0].Align, toks[0].Overflow}); got != want {
			t.Errorf("ParseText(%q) = %+v, want %+v", in, got, want)
		}
	}
	for _, in := range []string{"{{a | width:0}}", "{{a | width:x}}", "{{a | right:5000}}", "{{a | cut}}", "{{a | width:3 | right:4}}",
		"{{a | width:3 | cut | cut}}", "{{#a | width:3}}", "{{%a | right:3}}"} {
		if _, err := ParseText(in); err == nil {
			t.Errorf("ParseText(%q) takes it", in)
		}
	}
	// the templates of Office have no columns
	if _, err := Parse("{{a | width:20}}"); err == nil {
		t.Error("Parse takes width")
	}
}
