package xlsx

import (
	"bytes"
	"encoding/json"
	"testing"
)

// FuzzFill fills a template of four rows of cells with data: whatever the
// tags of the cells and the data, the template is filled or refused.
func FuzzFill(f *testing.F) {
	f.Add("{{#a}}{{v}}", "x{{/a}}", "=SUM(A1:A1)", "{{n}}", `{"a":[{"v":1},{"v":2}],"n":3}`)
	f.Add("{{#p | pages}}{{#r}}{{v}}", "", "{{/r}}", "{{/p}}{{_page}}", `{"p":[{"r":[{"v":1}]},{"r":[]}]}`)
	f.Add("{{#a | rows}}{{#b}}", "{{.}}{{/b}}{{/a}}", "=A1+B2", "{{^a}}none{{/a}}", `{"a":[{"b":[1,2]},{"b":[]}]}`)
	f.Add(`{{d | fmt:"yyyy/mm/dd"}}`, "{{%p}}", "=A1*2", `\{{x}}`, `{"d":"2026-10-13","p":"data:image/png;base64,AAAA"}`)
	f.Fuzz(func(t *testing.T, a, b, c, d, data string) {
		if len(a)+len(b)+len(c)+len(d)+len(data) > 4096 {
			return
		}
		var v map[string]any
		if json.Unmarshal([]byte(data), &v) != nil {
			v = map[string]any{}
		}
		sheet := rowXML(1, a, b) + rowXML(2, c, d) + rowXML(3, b, a) + rowXML(4, d, c)
		wb := book([]string{sheet}, `<definedNames><definedName name="_xlnm.Print_Area" localSheetId="0">S1!$A$1:$B$4</definedName></definedNames>`)
		opts := testOptions()
		res, err := Fill(bytes.NewReader(wb), int64(len(wb)), &TemplateData{Data: v}, opts)
		if err != nil {
			return
		}
		if res.Doc == nil || len(res.Doc.Views) == 0 {
			t.Fatal("a filled template without views")
		}
		for _, view := range res.Doc.Views {
			if len(view.Pages) == 0 {
				t.Fatalf("view %s has no pages", view.ID)
			}
		}
	})
}
