package tag

import (
	"strings"
	"testing"
)

// FuzzParse checks that any text splits into tokens that cover it, in
// order, or is refused with the place of the mistake.
func FuzzParse(f *testing.F) {
	for _, s := range []string{"", "plain", "{{a}}", "{{#a | pages}}x{{/a}}", `{{t | fmt:"#,##0"}}`, "｛｛＃明細｜slots｝｝", `\{{x}}`, "{{a", "{{}}", `{{a|fmt:"}}`, "{{%p|fit:cover}}{{^n}}{{/}}"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		toks, err := Parse(s)
		if err != nil {
			e, ok := err.(*Error)
			if !ok || e.Pos < 0 || e.Pos > len(s) || e.Msg == "" {
				t.Fatalf("Parse(%q): the error %v is not at a place of the text", s, err)
			}
			return
		}
		end := 0
		for _, tk := range toks {
			if tk.Pos != end || tk.End <= tk.Pos || tk.End > len(s) {
				t.Fatalf("Parse(%q): the token %+v does not follow the one before it (%d)", s, tk, end)
			}
			end = tk.End
			if tk.Kind != Text && tk.Kind != Close && (tk.Name == "" || len(tk.Path) == 0) {
				t.Fatalf("Parse(%q): the tag %+v has no name", s, tk)
			}
			if tk.Kind == Text && tk.Text == "" {
				t.Fatalf("Parse(%q): an empty text", s)
			}
		}
		if end != len(s) {
			t.Fatalf("Parse(%q): the tokens end at %d", s, end)
		}
		if len(toks) > 0 && !Has(s) && (len(toks) != 1 || toks[0].Kind != Text) && !strings.Contains(s, "\\") {
			t.Fatalf("Parse(%q): tags in a text that Has says has none", s)
		}
	})
}
