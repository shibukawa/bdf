package drawingml

import (
	"math"
	"slices"
	"testing"

	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// node and attr build a small element tree for the chart-cache tests.
func node(name string, attrs []ooxml.Attr, kids ...*ooxml.Node) *ooxml.Node {
	return &ooxml.Node{Name: name, Attrs: attrs, Kids: kids}
}

func attr(k, v string) []ooxml.Attr { return []ooxml.Attr{{Name: ooxml.Name{Local: k}, Value: v}} }

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

func TestChartAxisTicks(t *testing.T) {
	axisNode := func(step string) *ooxml.Node {
		return node("valAx", nil,
			node("scaling", nil, node("min", attr("val", "0")), node("max", attr("val", "1"))),
			node("majorUnit", attr("val", step)))
	}
	// Ordinary custom spacing is retained exactly.
	a := niceAxis(0, 1, axisNode("0.25"), 10)
	if got, want := a.ticks(), []float64{0, .25, .5, .75, 1}; !slices.Equal(got, want) {
		t.Fatalf("custom ticks = %v, want %v", got, want)
	}
	// A tiny custom unit would previously allocate labels without a useful
	// bound. Use the same spacing as an automatic axis instead.
	auto := niceAxis(0, 1, axisNode("0"), 10)
	for _, step := range []string{"1e-300", "NaN", "Inf"} {
		a := niceAxis(0, 1, axisNode(step), 10)
		if a.step != auto.step || !slices.Equal(a.ticks(), auto.ticks()) {
			t.Errorf("major unit %s: got step %g and ticks %v, want automatic step %g", step, a.step, a.ticks(), auto.step)
		}
	}
}

func TestChartAxisMalformedTicks(t *testing.T) {
	for _, tc := range []struct {
		name string
		a    axis
	}{
		{"rounding cannot advance", axis{lo: 1e16, hi: 1e16 + 10, step: .25}},
		{"tiny step", axis{lo: 0, hi: 1, step: 1e-300}},
		{"overflowing span", axis{lo: -math.MaxFloat64, hi: math.MaxFloat64, step: 1}},
		{"overflowing next tick", axis{lo: math.MaxFloat64, hi: math.MaxFloat64, step: math.MaxFloat64}},
		{"zero step", axis{lo: 0, hi: 1}},
		{"negative step", axis{lo: 0, hi: 1, step: -1}},
		{"NaN step", axis{lo: 0, hi: 1, step: math.NaN()}},
		{"infinite step", axis{lo: 0, hi: 1, step: math.Inf(1)}},
		{"NaN lower bound", axis{lo: math.NaN(), hi: 1, step: 1}},
		{"infinite upper bound", axis{lo: 0, hi: math.Inf(1), step: 1}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ticks := tc.a.ticks()
			if len(ticks) > maxAxisTicks {
				t.Fatalf("%d ticks exceeds limit %d", len(ticks), maxAxisTicks)
			}
			for i, v := range ticks {
				if math.IsNaN(v) || math.IsInf(v, 0) || i > 0 && v <= ticks[i-1] {
					t.Fatalf("tick %d = %g is nonfinite or does not advance", i, v)
				}
			}
		})
	}
}
