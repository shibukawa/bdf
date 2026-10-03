package ooxml

import (
	"strings"
	"testing"

	"github.com/shibukawa/bdf/internal/xmltree"
)

// A part whose elements nest deeper than xmltree.MaxDepth is rejected while the tree
// is built, so no later walk over it can overflow the stack. Nesting up to
// the limit still parses.
func TestParseDepthLimit(t *testing.T) {
	atLimit := []byte(strings.Repeat("<a>", xmltree.MaxDepth) + strings.Repeat("</a>", xmltree.MaxDepth))
	if _, err := Parse(atLimit); err != nil {
		t.Fatalf("nesting %d deep should parse: %v", xmltree.MaxDepth, err)
	}
	tooDeep := []byte(strings.Repeat("<a>", xmltree.MaxDepth+1) + strings.Repeat("</a>", xmltree.MaxDepth+1))
	if _, err := Parse(tooDeep); err == nil {
		t.Fatalf("nesting %d deep should be rejected", xmltree.MaxDepth+1)
	}
}
