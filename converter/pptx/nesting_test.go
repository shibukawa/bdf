package pptx

import (
	"strings"
	"testing"
)

// nestGroups wraps a slide's shape tree in depth nested groups.
func nestGroups(depth int) func(string) string {
	grp := `<p:grpSp><p:nvGrpSpPr><p:cNvPr id="10" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr>` +
		`<p:grpSpPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="100" cy="100"/>` +
		`<a:chOff x="0" y="0"/><a:chExt cx="100" cy="100"/></a:xfrm></p:grpSpPr>`
	return func(s string) string {
		nest := strings.Repeat(grp, depth) + strings.Repeat(`</p:grpSp>`, depth)
		return strings.Replace(s, `</p:spTree>`, nest+`</p:spTree>`, 1)
	}
}

// Deeply nested groups are dropped past the depth limit with a warning
// instead of overflowing the stack; a shallow nesting draws without one.
func TestNestedGroupDepth(t *testing.T) {
	warned := func(depth int) bool {
		res, _ := convertEdited(t, "basic.pptx", map[string]func(string) string{
			"ppt/slides/slide1.xml": nestGroups(depth),
		}, testOptions())
		for _, w := range res.Warnings {
			if strings.Contains(w, "nested deeper") {
				return true
			}
		}
		return false
	}
	if warned(20) {
		t.Errorf("20 nested groups should draw without a depth warning")
	}
	if !warned(150) {
		t.Errorf("150 nested groups should warn and be dropped, not crash")
	}
}
