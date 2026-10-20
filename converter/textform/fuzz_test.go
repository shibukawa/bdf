package textform

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/form/textpart"
)

// FuzzFill fills what may be a template of text with data: it is filled or
// refused, and a form has a print image that agrees with its pages.
func FuzzFill(f *testing.F) {
	f.Add("---\nlines: 6\ntitle: 1\n---\nT\n{{#a | rows}}\n{{v | width:3 | wrap}}|\n{{/a}}", `{"a":[{"v":"abcdefg"},{"v":"漢字かな"}]}`)
	f.Add("{{#_big}}大{{/_big}}{{#preprinted}}│{{/preprinted}}\n{{#s}}\n{{.}}\n\n{{/s}}\f{{^s}}none{{/s}}", `{"s":[1,2]}`)
	f.Add("---\nlines: 0\nmargin: {left: 1, top: 1}\n---\n{{#p | pages}}\n{{_page}}/{{_pages}} {{#_tall}}{{n | right:5 | cut}}{{/_tall}}\n{{/p}}", `{"p":[{"n":123456},{"n":"x"}]}`)
	f.Add("{{#a}}{{#b}}\n{{/b}}{{/a}}\t\x01{{x | center:4}}", `{"x":"a\nb"}`)
	f.Fuzz(func(t *testing.T, tpl, data string) {
		if len(tpl) > 1<<12 || len(data) > 1<<12 {
			return
		}
		var v map[string]any
		if json.Unmarshal([]byte(data), &v) != nil {
			v = map[string]any{}
		}
		Inspect(strings.NewReader(tpl), int64(len(tpl)))
		o := testOptions()
		o.NoTextIndex = true
		res, err := Fill(strings.NewReader(tpl), int64(len(tpl)), &TemplateData{Data: v}, o)
		if err != nil {
			return
		}
		im := res.Image
		if len(res.Doc.Views) != 1 || len(res.Doc.Views[0].Pages) != len(im.Pages) || len(im.Pages) == 0 {
			t.Fatalf("%d pages of the print image", len(im.Pages))
		}
		b, err := im.Encode()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := textpart.Decode(b); err != nil {
			t.Fatal(err)
		}
		if n := strings.Count(im.Text(textpart.TextOptions{}), "\f"); n != len(im.Pages)-1 {
			t.Fatalf("%d form feeds for %d pages", n, len(im.Pages))
		}
		for _, p := range im.Pages {
			for _, l := range p.Lines {
				cols := 0
				for _, c := range im.Cells(l) {
					if c.Col != cols {
						t.Fatalf("a cell of %q at column %d, after %d", l.Text, c.Col, cols)
					}
					cols += c.Cols
				}
			}
		}
	})
}
