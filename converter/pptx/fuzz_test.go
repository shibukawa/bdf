package pptx

import (
	"bytes"
	"encoding/json"
	"testing"
)

// FuzzFill fills a template of two slides (text boxes, a table and a
// picture) with data: whatever the tags of the texts and the data, the
// template is filled or refused.
func FuzzFill(f *testing.F) {
	f.Add("{{#a}}{{v}}", "x{{/a}}", "{{n}} {{_slide}}/{{_slides}}", "{{%p}}", `{"a":[{"v":1},{"v":2}],"n":3}`)
	f.Add("{{#p | pages}}{{#r}}{{v}}", "", "{{/r}}", "{{/p}}{{_page}}", `{"p":[{"r":[{"v":1}]},{"r":[]}]}`)
	f.Add("{{#a | rows}}{{#b}}", "{{.}}{{/b}}{{/a}}", "{{^a | pages}}none{{/a}}", "", `{"a":[{"b":[1,2]},{"b":[]}]}`)
	f.Add(`{{d | fmt:"yyyy/mm/dd"}}`, "{{%p | fit:cover}} alt", "{{#a | slots}}", `\{{x}}{{/a}}`, `{"d":"2026-10-13","p":"data:image/png;base64,AAAA"}`)
	f.Fuzz(func(t *testing.T, a, b, c, d, data string) {
		if len(a)+len(b)+len(c)+len(d)+len(data) > 4096 {
			return
		}
		var v map[string]any
		if json.Unmarshal([]byte(data), &v) != nil {
			v = map[string]any{}
		}
		deck := tplDeck(t,
			tsp("One", 36, 36, 300, 40, `<a:normAutofit/>`, tp(a, b), tp(c))+ttbl("Table", 36, 100, 400, 20, []string{a, b}, []string{c, d}, []string{b, a})+tpic("Pic", d, 500, 100, 80, 80),
			tsp("Two", 36, 36, 300, 40, `<a:spAutoFit/>`, tp(d), tp(c, a))+ttbl("Table", 36, 100, 400, 20, []string{d, c}))
		res, err := Fill(bytes.NewReader(deck), int64(len(deck)), &TemplateData{Data: v, Format: tplFormat}, testOptions())
		if err != nil {
			return
		}
		if res.Doc == nil || len(res.Doc.Views) != 1 {
			t.Fatal("a filled template without its view")
		}
	})
}
