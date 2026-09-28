package drawingml

import (
	"encoding/xml"
	"testing"

	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// node and attr build a small element tree for the chart-cache tests.
func node(name string, attrs []xml.Attr, kids ...*ooxml.Node) *ooxml.Node {
	return &ooxml.Node{Name: name, Attrs: attrs, Kids: kids}
}

func attr(k, v string) []xml.Attr { return []xml.Attr{{Name: xml.Name{Local: k}, Value: v}} }

// A cached point count (or a point index) far above the data present is not
// trusted as a slice size; a normal count is kept exactly.
func TestChartCachePointsCapped(t *testing.T) {
	ch := &chartCtx{s: New(Config{}).NewDrawing("", nil, nil)}

	huge := node("val", nil, node("numRef", nil, node("numCache", nil, node("ptCount", attr("val", "2000000000")))))
	if got := len(ch.cacheNumbers(huge, nil)); got != maxCachePoints {
		t.Errorf("cacheNumbers with a huge ptCount: %d entries, want the cap %d", got, maxCachePoints)
	}

	idx := node("val", nil, node("numRef", nil, node("numCache", nil, node("pt", attr("idx", "2000000000"), node("v", nil)))))
	if got := len(ch.cacheNumbers(idx, nil)); got != maxCachePoints {
		t.Errorf("cacheNumbers with a huge pt idx: %d entries, want the cap %d", got, maxCachePoints)
	}

	strs := node("cat", nil, node("strRef", nil, node("strCache", nil, node("ptCount", attr("val", "2000000000")))))
	if got := len(ch.cacheStrings(strs)); got != maxCachePoints {
		t.Errorf("cacheStrings with a huge ptCount: %d entries, want the cap %d", got, maxCachePoints)
	}

	// a normal count is unchanged
	ok := node("val", nil, node("numRef", nil, node("numCache", nil, node("ptCount", attr("val", "3")))))
	if got := len(ch.cacheNumbers(ok, nil)); got != 3 {
		t.Errorf("cacheNumbers with ptCount 3: %d entries, want 3", got)
	}
}

// A bullet number falls back to a decimal number when a hostile startAt
// would repeat a letter far more than any real list.
func TestAlphaCapped(t *testing.T) {
	if got := alpha(28); got != "BB" {
		t.Errorf("alpha(28) = %q, want BB", got)
	}
	if got := alpha(2000000000); got != itoa(2000000000) {
		t.Errorf("alpha(2000000000) = %q, want the decimal fallback", got)
	}
}
