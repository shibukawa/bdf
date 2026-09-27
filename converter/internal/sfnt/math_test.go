package sfnt

import (
	"os"
	"testing"
)

// The expected values come from fontTools (BoundsPen and the MATH table)
// on testdata/STIXTwoMath-cff-subset.otf, a subset of STIX Two Math.
func TestMathTableAndCFFBounds(t *testing.T) {
	data, err := os.ReadFile("testdata/STIXTwoMath-cff-subset.otf")
	if err != nil {
		t.Fatal(err)
	}
	f, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	m := ParseMath(f)
	if m == nil {
		t.Fatal("no MATH table")
	}
	c := m.Constants
	if c.ScriptPercentScaleDown != 70 || c.ScriptScriptPercentScaleDown != 55 || c.AxisHeight != 258 || c.FractionRuleThickness != 68 {
		t.Errorf("constants: %+v", c)
	}
	if m.MinConnectorOverlap != 100 {
		t.Errorf("MinConnectorOverlap = %d", m.MinConnectorOverlap)
	}
	paren := f.Cmap['(']
	v := m.Vertical(paren)
	if v == nil || len(v.Variants) != 13 || v.Variants[1].Advance != 1187 || len(v.Parts) != 3 {
		t.Fatalf("( construction: %+v", v)
	}
	if p := v.Parts[1]; !p.Extender || p.StartConnector != 1000 || p.FullAdvance != 1252 {
		t.Errorf("( extender: %+v", p)
	}
	integral := f.Cmap['∫']
	if it, ok := m.Italic(integral); !ok || it != 230 {
		t.Errorf("∫ italic = %d, %v", it, ok)
	}
	if ta, ok := m.TopAccent(f.Cmap['x']); !ok || ta != 250 {
		t.Errorf("x top accent = %d, %v", ta, ok)
	}

	o := NewOutlines(f)
	for _, tc := range []struct {
		r    rune
		want Rect
	}{
		{'(', Rect{45, -196, 327, 736}},
		{'x', Rect{-2, 0, 482, 473}},
		{'a', Rect{38, -7, 476, 485}},
		{'∑', Rect{62, -248, 861, 782}},
		{'∫', Rect{30, -226, 654, 727}},
		{'√', Rect{18, -265, 829, 922}},
		{'𝑥', Rect{23, -10, 540, 479}},
		{0x302, Rect{-371, 521, -89, 674}},
	} {
		got, ok := o.Bounds(f.Cmap[uint32(tc.r)])
		if !ok || !near(got, tc.want) {
			t.Errorf("bounds of %q = %v, %v; want %v", tc.r, got, ok, tc.want)
		}
	}
}

func near(a, b Rect) bool {
	d := func(x, y float64) bool { return x-y < 1 && y-x < 1 }
	return d(a.XMin, b.XMin) && d(a.YMin, b.YMin) && d(a.XMax, b.XMax) && d(a.YMax, b.YMax)
}

func TestTrueTypeBounds(t *testing.T) {
	data, err := os.ReadFile("../../docx/testdata/fonts/STIXTwoMath-subset.ttf")
	if err != nil {
		t.Fatal(err)
	}
	f, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := NewOutlines(f).Bounds(f.Cmap['('])
	if !ok || !near(got, Rect{45, -196, 327, 736}) {
		t.Errorf("bounds of ( = %v, %v", got, ok)
	}
	if ParseMath(f) == nil {
		t.Error("no MATH table")
	}
}
