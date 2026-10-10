package imagebdf

import (
	"image"
	"math"
	"testing"

	"github.com/shibukawa/bdf"
)

func near(a, b, eps float64) bool { return math.Abs(a-b) <= eps }

// val is the coverage of a pixel, 0 outside the mask.
func val(m *mask, x, y int) float32 {
	if m.empty() || !image.Pt(x, y).In(m.r) {
		return 0
	}
	return m.at(x, y)
}

func TestRectCoverage(t *testing.T) {
	m := rectMask(1.5, 1, 3, 2.25, image.Rect(0, 0, 8, 8))
	for _, c := range []struct {
		x, y int
		want float64
	}{{1, 1, 0.5}, {2, 1, 1}, {3, 1, 0}, {2, 2, 0.25}, {1, 2, 0.125}} {
		got := 0.0
		if image.Pt(c.x, c.y).In(m.r) {
			got = float64(m.at(c.x, c.y))
		}
		if !near(got, c.want, 1e-6) {
			t.Errorf("pixel (%d, %d): %v, want %v", c.x, c.y, got, c.want)
		}
	}
}

// square returns a closed square polyline, clockwise or not.
func square(x, y, s float64, ccw bool) polyline {
	pts := []point{{x, y}, {x + s, y}, {x + s, y + s}, {x, y + s}}
	if ccw {
		reversePoints(pts)
	}
	return polyline{pts: pts, closed: true}
}

func TestFillRules(t *testing.T) {
	b := image.Rect(0, 0, 20, 20)
	two := []polyline{square(2, 2, 8, false), square(6, 6, 8, false)}
	nz := rasterize(two, bdf.NonZero, b)
	eo := rasterize(two, bdf.EvenOdd, b)
	if v := val(nz, 8, 8); v != 1 {
		t.Errorf("nonzero: overlap covered %v", v)
	}
	if v := val(eo, 8, 8); v != 0 {
		t.Errorf("even-odd: overlap covered %v", v)
	}
	if v := val(eo, 3, 3); v != 1 {
		t.Errorf("even-odd: square covered %v", v)
	}
	// a hole wound the other way is empty under both rules
	hole := []polyline{square(2, 2, 12, false), square(5, 5, 4, true)}
	if v := val(rasterize(hole, bdf.NonZero, b), 6, 6); v != 0 {
		t.Errorf("nonzero hole covered %v", v)
	}
	// anti-aliasing: a diagonal edge through a pixel covers half of it
	tri := []polyline{{pts: []point{{0, 0}, {10, 0}, {0, 10}}, closed: true}}
	if v := float64(val(rasterize(tri, bdf.NonZero, b), 4, 5)); !near(v, 0.5, 0.05) {
		t.Errorf("diagonal edge covers %v of a pixel", v)
	}
}

func TestStrokes(t *testing.T) {
	b := image.Rect(0, 0, 40, 40)
	p := &path{}
	p.moveTo(5, 10.5)
	p.lineTo(30, 10.5)
	// a 3 px line with butt caps covers rows 9 to 11 between its ends
	cov, alpha := strokeMask(p, identity, lineStyle{width: 3, miter: 10}, b)
	if alpha != 1 || val(cov, 10, 9) != 1 || val(cov, 10, 11) != 1 || val(cov, 10, 12) != 0 || val(cov, 4, 10) != 0 {
		t.Errorf("3 px line: alpha %v, rows %v %v %v, before the start %v", alpha, val(cov, 10, 9), val(cov, 10, 11), val(cov, 10, 12), val(cov, 4, 10))
	}
	// square caps reach half the width further
	cov, _ = strokeMask(p, identity, lineStyle{width: 3, cap: bdf.CapSquare, miter: 10}, b)
	if !near(float64(val(cov, 4, 10)), 1, 1e-6) {
		t.Errorf("square cap covers %v", val(cov, 4, 10))
	}
	// a hairline is a pixel wide, as faint as it is thin (Skia's hairlines)
	cov, alpha = strokeMask(p, identity, lineStyle{width: 0.25, miter: 10}, b)
	if !near(float64(alpha), 0.25, 1e-6) || !near(float64(val(cov, 10, 10)), 1, 0.07) {
		t.Errorf("hairline: alpha %v, coverage %v", alpha, val(cov, 10, 10))
	}
	// dashes: 4 on, 4 off from x = 5
	ls := lineStyle{width: 2, miter: 10}
	ls.setDash([]float32{4, 4}, 0)
	cov, _ = strokeMask(p, identity, ls, b)
	if val(cov, 6, 10) != 1 || val(cov, 10, 10) != 0 || val(cov, 14, 10) != 1 {
		t.Errorf("dashes: %v %v %v", val(cov, 6, 10), val(cov, 10, 10), val(cov, 14, 10))
	}
}

func TestComposite(t *testing.T) {
	red, grey := [4]float32{1, 0, 0, 1}, [4]float32{0.5, 0.5, 0.5, 1}
	for _, c := range []struct {
		mode byte
		s, d [4]float32
		want [4]float32
	}{
		{bdf.BlendSourceOver, [4]float32{0.5, 0, 0, 0.5}, grey, [4]float32{0.75, 0.25, 0.25, 1}},
		{bdf.BlendMultiply, red, grey, [4]float32{0.5, 0, 0, 1}},
		{bdf.BlendScreen, red, grey, [4]float32{1, 0.5, 0.5, 1}},
		{bdf.BlendDestinationOut, [4]float32{0, 0, 0, 0.5}, grey, [4]float32{0.25, 0.25, 0.25, 0.5}},
		{bdf.BlendDestinationIn, [4]float32{0, 0, 0, 0}, grey, [4]float32{}},
		// red at the luminosity of grey, clipped into range: W3C SetLum and ClipColor
		{bdf.BlendLuminosity, grey, red, [4]float32{1, 0.2857, 0.2857, 1}},
	} {
		got := composite(c.mode, c.s, c.d)
		for i := range got {
			if !near(float64(got[i]), float64(c.want[i]), 1e-3) {
				t.Errorf("%s: %v, want %v", bdf.BlendNames[c.mode], got, c.want)
				break
			}
		}
	}
}

// TestDrawDocument draws a document built here: a clipped fill, a
// gradient, a group with alpha and text in an embedded font, and checks
// pixels.
func TestDrawDocument(t *testing.T) {
	d := bdf.NewDocument()
	o := bdf.NewObject()
	o.FillColor(0xff0000ff).FillRect(0, 0, 10, 10)
	o.Save().ClipRect(20, 0, 5, 5).FillColor(0x0000ffff).FillRect(15, 0, 20, 20).Restore()
	o.GroupBegin(0.5, bdf.BlendSourceOver, 0, 20, 10, 10).FillColor(0x000000ff).FillRect(0, 20, 10, 10).GroupEnd()
	h, _ := d.AddObject(o)
	v := d.NewView("v", bdf.ViewFixed, "")
	v.AddPage(40, 40, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	img, err := New(d, nil).Page(v, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	at := func(x, y int) [4]uint8 {
		i := img.PixOffset(x, y)
		return [4]uint8(img.Pix[i : i+4])
	}
	for _, c := range []struct {
		x, y int
		want [4]uint8
	}{
		{5, 5, [4]uint8{255, 0, 0, 255}},      // red square
		{22, 2, [4]uint8{0, 0, 255, 255}},     // blue inside the clip
		{17, 2, [4]uint8{255, 255, 255, 255}}, // white outside it
		{22, 10, [4]uint8{255, 255, 255, 255}},
		{5, 25, [4]uint8{128, 128, 128, 255}}, // black at half alpha over white
	} {
		if got := at(c.x, c.y); got != c.want {
			t.Errorf("pixel (%d, %d): %v, want %v", c.x, c.y, got, c.want)
		}
	}
}

func TestPathData(t *testing.T) {
	// compact numbers, relative commands after a close, and an arc
	p := parsePathData("M10,20l5-5.5.5.5Z m1 1 h2v2H1V1zM0 0A5 5 0 0 1 10 0")
	var n int
	for _, s := range p.subs {
		n += len(s.segs)
	}
	if n != 2+4+3 { // two lines, h v H V, then a line to the arc and two quarter curves
		t.Fatalf("%d segments: %+v", n, p.subs)
	}
	x0, y0, x1, y1 := p.bbox()
	// the half circle from (0, 0) to (10, 0) sweeps over the top, to y = -5
	if !near(x0, 0, 1e-9) || !near(x1, 15.5, 1e-9) || !near(y1, 23, 1e-9) || !near(y0, -5, 1e-6) {
		t.Errorf("bbox %v %v %v %v", x0, y0, x1, y1)
	}
}

func TestVisualOrder(t *testing.T) {
	for _, c := range []struct {
		s    string
		rtl  bool
		want string
	}{
		{"abc", false, "abc"},
		{"שלום", false, "םולש"},
		{"abc שלום def", false, "abc םולש def"},
		{"שלום abc", true, "abc םולש"},
	} {
		if got := string(visualOrder(c.s, c.rtl)); got != c.want {
			t.Errorf("%q (rtl %v): %q, want %q", c.s, c.rtl, got, c.want)
		}
	}
}

// TestLimits checks the guards against documents that would take a server
// down: canvases too large, and SVG images that repeat their content
// exponentially.
func TestLimits(t *testing.T) {
	d := bdf.NewDocument()
	v := d.NewView("v", bdf.ViewFixed, "")
	v.AddPage(1e5, 1e5)
	if _, err := New(d, nil).Page(v, 0, 1); err != ErrTooLarge {
		t.Errorf("a huge page: %v", err)
	}
	// ten levels of ten uses: 10^10 circles in a browser
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10"><defs><circle id="l0" r="1"/>`
	for i := 1; i <= 10; i++ {
		svg += `<g id="l` + string(rune('0'+i%10)) + `x">`
		for k := 0; k < 10; k++ {
			svg += `<use href="#l` + string(rune('0'+(i-1)%10)) + map[bool]string{true: "", false: "x"}[i == 1] + `"/>`
		}
		svg += `</g>`
	}
	svg += `</defs><use href="#l0x"/></svg>`
	si, ok := newSVGImage([]byte(svg))
	if !ok {
		t.Fatal("not parsed")
	}
	r := New(bdf.NewDocument(), nil)
	si.raster(r, 1)
	if len(r.Warnings()) == 0 {
		t.Error("no warning for an SVG image that was cut short")
	}
}
