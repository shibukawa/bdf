package xlsx

import (
	"strconv"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/form/tag"
)

func TestInlineSectionDepth(t *testing.T) {
	for _, mode := range []string{"", " | rows"} {
		for _, depth := range []int{maxSectionDepth, maxSectionDepth + 1} {
			src := strings.Repeat("{{#a"+mode+"}}", depth) + "x" + strings.Repeat("{{/a}}", depth)
			ts := &tplSheet{ws: &worksheet{name: "S", rows: []row{{idx: 0, cells: []cell{{kind: cellStr, text: &richText{plain: src}}}}}}, cells: map[[2]int]*tplCell{}}
			f := newFiller(&TemplateData{})
			err := f.scan(ts)
			if depth == maxSectionDepth && err != nil {
				t.Errorf("%d sections%s: %v", depth, mode, err)
			}
			if depth > maxSectionDepth && (err == nil || !strings.Contains(err.Error(), "inside more than")) {
				t.Errorf("%d sections%s: %v, want a nesting error", depth, mode, err)
			}
		}
	}
}

// Empty inline sections and repeated rows without any emitted source rows
// still consume work; neither can evade the limit by producing no cells.
func TestEmptyExpansionBudget(t *testing.T) {
	for _, inline := range []bool{false, true} {
		f := newFiller(&TemplateData{Data: map[string]any{"a": make([]any, 8)}})
		f.steps = maxExpansionSteps - 4
		ts := &tplSheet{ws: &worksheet{name: "S"}}
		x := &expander{f: f, ts: ts}
		var err error
		if inline {
			toks, perr := tag.Parse("{{#a}}{{/a}}")
			if perr != nil {
				t.Fatal(perr)
			}
			_, err = x.cell(srcCell{tc: &tplCell{toks: toks, rowLevel: make([]bool, len(toks))}}, 0, f.root)
		} else {
			s := &tplSection{name: "a", path: []string{"a"}, mode: tag.ModeRows}
			_, err = x.unit([]*tplSection{s}, 0, 0, f.root, &itemInst{}, 0)
		}
		if err == nil || !strings.Contains(err.Error(), "expansion steps") {
			t.Errorf("inline=%v: %v, want a work-limit error", inline, err)
		}
		if f.cells != 0 || len(x.out) != 0 {
			t.Errorf("inline=%v: emitted %d cells and %d rows", inline, f.cells, len(x.out))
		}
	}
}

func TestClipTextPrefix(t *testing.T) {
	s := strings.Repeat("界", maxCellText+100)
	if got := clipText(s); got != strings.Repeat("界", maxCellText) {
		t.Errorf("clipped text has %d bytes", len(got))
	}
	if n := testing.AllocsPerRun(10, func() { _ = clipText(s) }); n != 0 {
		t.Errorf("clipping an oversized value allocated %v times", n)
	}
}

func BenchmarkRepeatedInlineText(b *testing.B) {
	for _, n := range []int{1000, 4000} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			f := newFiller(&TemplateData{Data: map[string]any{"a": make([]any, n)}})
			toks, err := tag.Parse("{{#a}}x{{/a}}")
			if err != nil {
				b.Fatal(err)
			}
			font := &xfont{}
			src := srcCell{cl: cell{kind: cellStr, text: &richText{runs: []textRun{{text: "{{#a}}x{{/a}}", font: font}}}}, tc: &tplCell{toks: toks, rowLevel: make([]bool, len(toks))}}
			x := &expander{f: f, ts: &tplSheet{ws: &worksheet{name: "S"}}}
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				f.steps = 0
				if _, err := x.cell(src, 0, f.root); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
