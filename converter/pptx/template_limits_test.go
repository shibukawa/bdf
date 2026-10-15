package pptx

import (
	"strings"
	"testing"

	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

func TestInlineSectionDepth(t *testing.T) {
	for _, depth := range []int{maxSectionDepth, maxSectionDepth + 1} {
		p, err := ooxml.Parse([]byte(tp(strings.Repeat("{{#a}}", depth) + "x" + strings.Repeat("{{/a}}", depth))))
		if err != nil {
			t.Fatal(err)
		}
		f := newFiller(&TemplateData{})
		var stack []open
		err = f.scanPara(p, loc{name: "paragraph"}, &stack)
		if depth == maxSectionDepth && err != nil {
			t.Errorf("%d sections: %v", depth, err)
		}
		if depth > maxSectionDepth && (err == nil || !strings.Contains(err.Error(), "inside more than")) {
			t.Errorf("%d sections: %v, want a nesting error", depth, err)
		}
	}
}

func TestEmptyExpansionBudget(t *testing.T) {
	for _, inline := range []bool{false, true} {
		f := newFiller(&TemplateData{Data: map[string]any{"a": make([]any, 8), "empty": []any{}}})
		f.steps = maxExpansionSteps - 4
		p, err := ooxml.Parse([]byte(tp("{{#a}}{{/a}}")))
		if err != nil {
			t.Fatal(err)
		}
		var out []*ooxml.Node
		if inline {
			var stack []open
			if err = f.scanPara(p, loc{name: "paragraph"}, &stack); err != nil {
				t.Fatal(err)
			}
			_, err = f.fillPara(p, f.root, "paragraph")
		} else {
			s := &sec{name: "a", path: []string{"a"}, children: []*sec{{name: "empty", path: []string{"empty"}}}}
			err = f.fillParas([]*ooxml.Node{p}, []*sec{s}, 0, 0, f.root, "paragraph", &out)
		}
		if err == nil || !strings.Contains(err.Error(), "expansion steps") {
			t.Errorf("inline=%v: %v, want a work-limit error", inline, err)
		}
		if f.nodes != 0 || len(out) != 0 {
			t.Errorf("inline=%v: emitted %d nodes and %d paragraphs", inline, f.nodes, len(out))
		}
	}
}

func TestTextPrefix(t *testing.T) {
	s := strings.Repeat("界", maxParaText+100)
	got, n := textPrefix(s, maxParaText)
	if n != maxParaText || got != strings.Repeat("界", maxParaText) {
		t.Errorf("clipped text has %d runes and %d bytes", n, len(got))
	}
}
