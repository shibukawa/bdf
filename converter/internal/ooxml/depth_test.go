package ooxml

import (
	"strings"
	"testing"
)

// A part whose elements nest deeper than maxDepth is rejected while the tree
// is built, so no later walk over it can overflow the stack. Nesting up to
// the limit still parses.
func TestParseDepthLimit(t *testing.T) {
	atLimit := []byte(strings.Repeat("<a>", maxDepth) + strings.Repeat("</a>", maxDepth))
	if _, err := Parse(atLimit); err != nil {
		t.Fatalf("nesting %d deep should parse: %v", maxDepth, err)
	}
	tooDeep := []byte(strings.Repeat("<a>", maxDepth+1) + strings.Repeat("</a>", maxDepth+1))
	if _, err := Parse(tooDeep); err == nil {
		t.Fatalf("nesting %d deep should be rejected", maxDepth+1)
	}
}
