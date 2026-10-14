package wordproc

import (
	"strconv"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

func templatePara(t testing.TB, text string) *ooxml.Node {
	t.Helper()
	p, err := ooxml.Parse([]byte("<p><r><t>" + text + "</t></r></p>"))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestInlineSectionDepth(t *testing.T) {
	for _, depth := range []int{maxSectionDepth, maxSectionDepth + 1} {
		p := templatePara(t, strings.Repeat("{{#a}}", depth)+"x"+strings.Repeat("{{/a}}", depth))
		f := newFiller(&TemplateData{})
		var stack []openSec
		err := f.scanInline(p, tplLoc{name: "paragraph"}, &stack, nil)
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
		p := templatePara(t, "{{#a}}{{/a}}")
		var err error
		var out []*ooxml.Node
		if inline {
			var stack []openSec
			if err = f.scanInline(p, tplLoc{name: "paragraph"}, &stack, nil); err != nil {
				t.Fatal(err)
			}
			_, err = f.fillInline(p, f.inl[p], f.root)
		} else {
			s := &tplSec{name: "a", path: []string{"a"}, children: []*tplSec{{name: "empty", path: []string{"empty"}}}}
			err = f.fillKids(&ooxml.Node{Kids: []*ooxml.Node{p}}, []*tplSec{s}, 0, 0, f.root, &out)
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
	if n := testing.AllocsPerRun(10, func() { _, _ = textPrefix(s, maxParaText) }); n != 0 {
		t.Errorf("clipping an oversized value allocated %v times", n)
	}
}

func BenchmarkRepeatedInlineText(b *testing.B) {
	for _, n := range []int{1000, 4000} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			f := newFiller(&TemplateData{Data: map[string]any{"a": make([]any, n)}})
			p := templatePara(b, "{{#a}}x{{/a}}")
			var stack []openSec
			if err := f.scanInline(p, tplLoc{name: "paragraph"}, &stack, nil); err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				f.steps, f.nodes = 0, 0
				if _, err := f.fillInline(p, f.inl[p], f.root); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
