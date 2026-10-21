package docx

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/shibukawa/bdf/converter/internal/convtest"
)

// FuzzFill fills a template (paragraphs, a table, a picture and a header)
// with data: whatever the tags of the texts and the data, the template is
// filled or refused.
func FuzzFill(f *testing.F) {
	f.Add("{{#a}}{{v}}", "x{{/a}}", "{{n}}", "{{%p}}", `{"a":[{"v":1},{"v":2}],"n":3}`)
	f.Add("{{#p | pages}}{{#r}}{{v}}", "", "{{/r}}", "{{/p}}{{_page}}", `{"p":[{"r":[{"v":1}]},{"r":[]}]}`)
	f.Add("{{#a | rows}}{{#b}}", "{{.}}{{/b}}{{/a}}", "{{^a | rows}}none{{/a}}", "", `{"a":[{"b":[1,2]},{"b":[]}]}`)
	f.Add(`{{d | fmt:"yyyy/mm/dd"}}`, "{{%p | fit:cover}} alt", "{{#a | slots}}", `\{{x}}{{/a}}`, `{"d":"2026-10-13","p":"data:image/png;base64,AAAA"}`)
	f.Fuzz(func(t *testing.T, a, b, c, d, data string) {
		if len(a)+len(b)+len(c)+len(d)+len(data) > 4096 {
			return
		}
		var v map[string]any
		if json.Unmarshal([]byte(data), &v) != nil {
			v = map[string]any{}
		}
		body := tp(a, b) + tp(c) + ttbl([]string{a, b}, []string{c, d}, []string{b, a}) + tp(tpic(d), c) +
			tp(`<w:hyperlink w:anchor="x"><w:r><w:t>`+"link"+`</w:t></w:r></w:hyperlink>`, d) + ttbl([]string{d + "\n" + c}) + tp(b, a) +
			`<w:sectPr><w:headerReference w:type="default" r:id="rIdheader"/><w:pgSz w:w="6000" w:h="8000"/><w:pgMar w:top="1000" w:bottom="1000" w:left="1000" w:right="1000" w:header="500" w:footer="500"/></w:sectPr>`
		parts := map[string]string{
			"word/header.xml":       `<w:hdr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">` + tp(c) + `</w:hdr>`,
			"word/media/image1.png": string(convtest.PNG(4, 4)),
		}
		file := docxBytes(t, body, parts)
		res, err := Fill(bytes.NewReader(file), int64(len(file)), &TemplateData{Data: v, Format: tplFormat}, testOptions())
		if err != nil {
			return
		}
		if res.Doc == nil || len(res.Doc.Views) != 1 {
			t.Fatal("a filled template without its view")
		}
	})
}
