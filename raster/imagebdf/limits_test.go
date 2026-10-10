package imagebdf

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"math"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/image/svg"
)

// These tests draw what a document made to take a server down asks for:
// each would run out of memory, of time or into a panic without the limit
// it checks.

// onePage returns a document with one page of w × h that an object draws.
func onePage(o *bdf.Object, w, h float32) (*bdf.Document, *bdf.View) {
	d := bdf.NewDocument()
	hash, _ := d.AddObject(o)
	v := d.NewView("v", bdf.ViewFixed, "")
	v.AddPage(w, h, bdf.Layer{Role: bdf.RoleBody, Obj: hash})
	return d, v
}

// drawn draws the page of onePage at scale 1 and returns its pixels.
func drawn(t *testing.T, o *bdf.Object, w, h float32) (*Renderer, func(x, y int) [4]uint8) {
	t.Helper()
	d, v := onePage(o, w, h)
	r := New(d, &Options{NoSystemFonts: true})
	img, err := r.Page(v, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	return r, func(x, y int) [4]uint8 {
		i := img.PixOffset(x, y)
		return [4]uint8(img.Pix[i : i+4])
	}
}

func warned(r *Renderer, part string) bool {
	return slices.ContainsFunc(r.Warnings(), func(w string) bool { return strings.Contains(w, part) })
}

var (
	white = [4]uint8{255, 255, 255, 255}
	red   = [4]uint8{255, 0, 0, 255}
	blue  = [4]uint8{0, 0, 255, 255}
)

func TestLayersPastTheLimit(t *testing.T) {
	for _, mask := range []bool{false, true} {
		o := bdf.NewObject()
		o.FillColor(0xff0000ff)
		if mask {
			o.GroupBegin(1, bdf.BlendSourceOver, 0, 0, 8, 8)
		}
		for i := 0; i < 100; i++ {
			if i == 10 {
				o.FillColor(0x0000ffff).FillRect(0, 0, 2, 2) // in a group within the limit
			}
			if mask {
				o.MaskBegin(bdf.MaskAlpha, 0, nil)
			} else {
				o.GroupBegin(1, bdf.BlendSourceOver, 0, 0, 8, 8)
			}
		}
		o.FillColor(0x0000ffff).FillRect(4, 0, 2, 2) // in one past it
		for i := 0; i < 100; i++ {
			if mask {
				o.MaskEnd()
			} else {
				o.GroupEnd()
			}
		}
		if mask {
			o.GroupEnd()
		}
		o.FillRect(4, 4, 2, 2) // after them, in the state from before
		r, at := drawn(t, o, 8, 8)
		if !warned(r, "groups and soft masks") {
			t.Errorf("mask %v: no warning for 100 open at a time: %v", mask, r.Warnings())
		}
		if !mask && at(1, 1) != blue {
			t.Errorf("a group within the limit is drawn %v", at(1, 1))
		}
		if at(5, 1) != white {
			t.Errorf("mask %v: what is in a group past the limit is drawn %v", mask, at(5, 1))
		}
		if at(5, 5) != red {
			t.Errorf("mask %v: after the groups: %v, want the fill of the state before them", mask, at(5, 5))
		}
	}
}

// canvasOf starts the image of a page of onePage, to draw its object with
// other limits than those of a Page.
func canvasOf(d *bdf.Document, v *bdf.View) (*Renderer, *drawer, func()) {
	r := New(d, &Options{NoSystemFonts: true})
	p := v.Pages[0]
	dr, m := r.canvas(bdf.Rect{W: p.W, H: p.H}, int(p.W), int(p.H))
	return r, dr, func() { dr.drawTop(p.Layers[0].Obj, m, dr.clipTo(m, 0, 0, float64(p.W), float64(p.H))) }
}

func TestMemoryOfLayersAndClips(t *testing.T) {
	const canvas = 16 * 16 * 16 // bytes
	// groups as large as the page, in each other
	o := bdf.NewObject()
	for i := 0; i < 5; i++ {
		o.GroupBegin(1, bdf.BlendSourceOver, 0, 0, 16, 16)
	}
	for i := 0; i < 5; i++ {
		o.GroupEnd()
	}
	d, v := onePage(o, 16, 16)
	r, dr, draw := canvasOf(d, v)
	dr.limit = 3 * canvas
	draw()
	if !warned(r, "too much memory") {
		t.Errorf("five canvases where three are the limit: %v", r.Warnings())
	}
	// clips that are no rectangles of whole pixels, each kept by a state
	o = bdf.NewObject()
	for i := 0; i < 40; i++ {
		o.Save().ClipRect(0.5, 0.5, 15, 15)
	}
	o.FillRect(0, 0, 16, 16)
	d, v = onePage(o, 16, 16)
	r, dr, draw = canvasOf(d, v)
	dr.limit = 3 * canvas
	draw()
	if !warned(r, "too much memory") {
		t.Errorf("40 clips, a quarter of a canvas each, where three canvases are the limit: %v", r.Warnings())
	}
	if got := dr.ctx.target.pix[4*(8*16+8)+3]; got != 0 {
		t.Errorf("what is drawn in a clip past the limit: alpha %v", got)
	}

	// what is counted is let go of: groups, masks, states with clips and
	// objects that leave theirs behind
	child := bdf.NewObject()
	child.Save().ClipRect(0.5, 0.5, 9, 9).Save().ClipRect(1.5, 1.5, 7, 7).FillRect(0, 0, 16, 16)
	d = bdf.NewDocument()
	ch, bb := d.AddObject(child)
	o = bdf.NewObject()
	ref := o.AddObject(ch, bb)
	for i := 0; i < 3; i++ {
		o.Save().ClipRect(0.25, 0.25, 12, 12)
		o.GroupBegin(0.5, bdf.BlendSourceOver, 0, 0, 16, 16).ClipRect(0.5, 0.5, 8, 8)
		o.Use(ref)
		o.MaskBegin(bdf.MaskAlpha, 0, nil).ClipRect(0.5, 0.5, 3, 3).FillRect(0, 0, 8, 8).MaskEnd()
		o.GroupEnd()
		o.ClipRect(0.75, 0.75, 8, 8).Restore()
	}
	hash, _ := d.AddObject(o)
	r = New(d, &Options{NoSystemFonts: true})
	dr, m := r.canvas(bdf.Rect{W: 16, H: 16}, 16, 16)
	dr.ctx.st = initialState(m, &mask{r: dr.ctx.target.bounds()})
	dr.run(r.object(hash))
	if dr.live != 0 || dr.ctx.held != 0 || len(r.Warnings()) != 0 {
		t.Errorf("after an object that ends what it begins: %d bytes counted (%d of the page), warnings %v", dr.live, dr.ctx.held, r.Warnings())
	}
}

func TestStatesPastTheLimit(t *testing.T) {
	o := bdf.NewObject()
	o.FillColor(0xff0000ff)
	for i := 0; i < 2*maxStates; i++ {
		o.Save().FillColor(0x0000ffff)
	}
	for i := 0; i < 2*maxStates; i++ {
		o.Restore()
	}
	o.FillRect(0, 0, 4, 4)
	d, v := onePage(o, 4, 4)
	r, dr, draw := canvasOf(d, v)
	draw()
	if !warned(r, "states are saved") {
		t.Errorf("no warning for %d states: %v", 2*maxStates, r.Warnings())
	}
	// every RESTORE ends a SAVE: the last one restores the first state
	if got := dr.ctx.target.pix[0:4]; got[0] != 1 || got[2] != 0 {
		t.Errorf("filled %v after as many RESTOREs as SAVEs, want the red of the first state", got)
	}
}

func TestObjectsDrawnAgain(t *testing.T) {
	// four objects that each draw the next one ten times: 10^4 fills
	d := bdf.NewDocument()
	leaf := bdf.NewObject()
	leaf.FillColor(0x0000ffff).FillRect(0, 0, 1, 1)
	h, bb := d.AddObject(leaf)
	for l := 0; l < 4; l++ {
		o := bdf.NewObject()
		ref := o.AddObject(h, bb)
		for k := 0; k < 10; k++ {
			o.UseAt(ref, float32(k)/4, 0)
		}
		if l == 3 {
			o.FillColor(0xff0000ff).FillRect(6, 6, 2, 2) // of an object drawn once
		}
		h, bb = d.AddObject(o)
	}
	v := d.NewView("v", bdf.ViewFixed, "")
	v.AddPage(8, 8, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	r, dr, draw := canvasOf(d, v)
	dr.most = 500
	draw()
	if !warned(r, "drawn again too many times") || dr.reused != 500 {
		t.Errorf("%d instructions of objects drawn again where 500 are the limit: %v", dr.reused, r.Warnings())
	}
	at := func(x, y int) []float32 { return dr.ctx.target.pix[4*(y*8+x):][:4] }
	if p := at(0, 0); p[2] != 1 || p[0] != 0 {
		t.Errorf("what was drawn within the limit: %v", p)
	}
	if p := at(7, 7); p[0] != 1 || p[2] != 0 {
		t.Errorf("what is drawn once, after the limit was passed: %v", p)
	}
	// and within the limit of a Page nothing is left out
	r = New(d, &Options{NoSystemFonts: true})
	if _, err := r.Page(v, 0, 1); err != nil || len(r.Warnings()) != 0 {
		t.Errorf("10^4 fills: %v %v", err, r.Warnings())
	}

	// objects that draw each other deeper than bdf.MaxUseDepth
	d = bdf.NewDocument()
	h, bb = d.AddObject(leaf)
	for l := 0; l < bdf.MaxUseDepth+8; l++ {
		o := bdf.NewObject()
		o.Translate(0, 0.001) // no two objects are the same
		o.Use(o.AddObject(h, bb))
		h, bb = d.AddObject(o)
	}
	v = d.NewView("v", bdf.ViewFixed, "")
	v.AddPage(8, 8, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	r = New(d, &Options{NoSystemFonts: true})
	if _, err := r.Page(v, 0, 1); err != nil || !warned(r, "nested too deep") {
		t.Errorf("objects %d deep: %v %v", bdf.MaxUseDepth+8, err, r.Warnings())
	}
}

// TestReferencesPastTheTables draws instructions that refer to objects and
// fonts by numbers no table has, which no int holds either.
func TestReferencesPastTheTables(t *testing.T) {
	huge := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01} // 2^64-1 as a varuint
	var ops []byte
	ops = append(append(ops, bdf.OpUse), huge...)
	ops = append(append(append(ops, bdf.OpUseAt), huge...), make([]byte, 8)...)
	ops = append(append(append(ops, bdf.OpFont), huge...), 0, 0, 0x20, 0x41) // size 10
	ops = append(append(append(ops, bdf.OpFillPath), huge...), 0)
	ops = append(append(ops, bdf.OpFillPaint), huge...)
	ops = append(append(append(ops, bdf.OpImage), huge...), make([]byte, 16)...)
	o := &bdf.ObjectPart{Ops: ops, Objects: make([]bdf.Hash, 1), Fonts: make([]bdf.Font, 1)}
	r := New(bdf.NewDocument(), &Options{NoSystemFonts: true})
	dr, _ := r.canvas(bdf.Rect{W: 4, H: 4}, 4, 4)
	dr.run(o) // panics on an index out of range without the checks
	if !warned(r, "bad path reference") {
		t.Errorf("warnings %v", r.Warnings())
	}
}

func TestShadowValues(t *testing.T) {
	r := New(bdf.NewDocument(), nil)
	dr, _ := r.canvas(bdf.Rect{W: 4, H: 4}, 4, 4)
	shadow := func(blur, dx, dy float32) *state {
		dr.instr(bdf.Instr{Op: bdf.OpShadow, Args: []any{uint64(0xff), blur, dx, dy}})
		return &dr.ctx.st
	}
	if st := shadow(1e30, -1e30, 3); st.shadowBlur != maxShadowBlur || st.shadowDX != -(1<<24) || st.shadowDY != 3 {
		t.Errorf("blur %v, offset %v %v", st.shadowBlur, st.shadowDX, st.shadowDY)
	}
	nan, inf := float32(math.NaN()), float32(math.Inf(1))
	if st := shadow(nan, inf, nan); st.shadowBlur != maxShadowBlur || st.shadowDX != -(1<<24) || st.shadowDY != 3 {
		t.Errorf("values that are no numbers are not ignored: blur %v, offset %v %v", st.shadowBlur, st.shadowDX, st.shadowDY)
	}
	// shadows that fall outside the page take no memory for their blur
	for _, c := range [][3]float32{{1e12, 1e30, 0}, {1e30, 0, -1e30}, {2, 1e30, 0}, {0, 0, 1e9}} {
		o := bdf.NewObject()
		o.Shadow(0x000000ff, c[0], c[1], c[2]).FillColor(0xff0000ff).FillRect(1, 1, 2, 2)
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		before := ms.TotalAlloc
		_, at := drawn(t, o, 4, 4) // panics without the limits
		runtime.ReadMemStats(&ms)
		if at(1, 1) != red || at(0, 0) != white {
			t.Errorf("shadow %v: %v %v", c, at(1, 1), at(0, 0))
		}
		if n := ms.TotalAlloc - before; n > 1<<20 {
			t.Errorf("shadow %v: %d bytes allocated", c, n)
		}
	}
	// one that falls on the page is drawn as before
	o := bdf.NewObject()
	o.Shadow(0x000000ff, 0, 4, 0).FillColor(0xff0000ff).FillRect(1, 1, 2, 2)
	if _, at := drawn(t, o, 8, 4); at(5, 1) != [4]uint8{0, 0, 0, 255} {
		t.Errorf("a shadow 4 to the right: %v", at(5, 1))
	}
}

func TestStrokeLimits(t *testing.T) {
	line := []polyline{{pts: []point{{0, 0}, {100, 0}}}}
	// a line far wider than the page, with round caps on its dashes
	ls := lineStyle{width: 6000, cap: bdf.CapRound, join: bdf.JoinRound, miter: 10}
	ls.setDash([]float32{0.001, 0.001}, 0)
	s := newStroker(&ls, 1)
	s.points = 5000
	s.stroke(line)
	n := 0
	for _, p := range s.out {
		n += len(p.pts)
	}
	if !s.cut || n > 5000+strokePointsEach*2 {
		t.Errorf("cut %v, %d points where 5000 are the limit", s.cut, n)
	}
	// dashes far too fine for the line, which lie outside of what shows
	s = newStroker(&ls, 1).within(identity.translate(1e6, 0), image.Rect(0, 0, 16, 16))
	start := time.Now()
	s.stroke([]polyline{{pts: []point{{0, 0}, {1e6, 0}}}})
	if !s.cut || len(s.out) != 0 || s.dashes >= 0 {
		t.Errorf("cut %v, %d polygons, %d dashes left", s.cut, len(s.out), s.dashes)
	}
	if e := time.Since(start); e > 2*time.Second {
		t.Errorf("%d dashes outside the page took %v", strokeDashes, e)
	}
	// within the limits the outline is whole
	ls = lineStyle{width: 2, join: bdf.JoinRound, miter: 10}
	s = newStroker(&ls, 1)
	s.stroke([]polyline{{pts: []point{{0, 0}, {10, 0}, {10, 10}}}})
	if s.cut || len(s.out) != 3 {
		t.Errorf("two segments and a join: cut %v, %d polygons", s.cut, len(s.out))
	}
	// a circle far larger than the tolerance has 256 points on every machine
	if n := len(circle(point{}, 1e30, 0.2)); n != 256 {
		t.Errorf("a circle of radius 1e30 has %d points", n)
	}
}

// slivers is a path of n thin triangles as high as the page, drawn from
// right to left: every scanline crosses each.
func slivers(n int, w, h float32) *bdf.Path {
	p := &bdf.Path{}
	for i := 0; i < n; i++ {
		x := w - float32(i)*w/float32(n)
		p.MoveTo(x, 0).LineTo(x+0.01, h).LineTo(x+0.02, 0).Close()
	}
	return p
}

func TestCrossingsInNoOrder(t *testing.T) {
	o := bdf.NewObject()
	o.FillColor(0xff0000ff).FillPath(o.AddPath(slivers(1500, 200, 200)), bdf.NonZero)
	d, v := onePage(o, 200, 200)
	start := time.Now()
	if _, err := New(d, nil).Page(v, 0, 1); err != nil {
		t.Fatal(err)
	}
	// a minute when the crossings are sorted by insertion
	if e := time.Since(start); e > 10*time.Second {
		t.Errorf("1500 slivers drawn from right to left took %v", e)
	}
}

func TestSVGRasterSize(t *testing.T) {
	for _, c := range []struct {
		w, h, k float64
	}{{0.0001, 10000, 1e9}, {1e-9, 1e9, 1e9}, {1e9, 1e-9, 1}, {1e-300, 1e300, 1e300}, {1e200, 1e200, 1}, {3, 1e-320, 1e308}} {
		si := &svgImage{w: c.w, h: c.h}
		_, k, w, h := si.rasterSize(c.k)
		if w*h > 2*svgRasterPixels || w < 0 || h < 0 || math.IsNaN(k) {
			t.Errorf("an image of %g × %g at %g: a raster of %d × %d at %g", c.w, c.h, c.k, w, h, k)
		}
	}
	// the images that are a pixel wide and high at least are drawn at the
	// scale they were before there was a limit
	before := func(w, h, k float64) (int, int) {
		if max := math.Sqrt(svgRasterPixels / (w * h)); k > max {
			k = max
		}
		k = math.Pow(2, float64(int(math.Ceil(2*math.Log2(k))))/2)
		return max(1, int(math.Ceil(w*k))), max(1, int(math.Ceil(h*k)))
	}
	for _, c := range [][3]float64{{300, 150, 1}, {300, 150, 3}, {300, 150, 0.01}, {1000, 1000, 100}, {10, 100000, 50}, {1, 16777216, 9}, {4096, 4096, 1.5},
		{5792, 5793, 1}, {0.5, 100, 1}, {0.01, 1e4, 8}, {16777216, 2, 7}, {1e6, 1e6, 1}, {1.5, 2e7, 3}} {
		si := &svgImage{w: c[0], h: c[1]}
		pw, ph := before(c[0], c[1], c[2])
		if _, _, w, h := si.rasterSize(c[2]); w != pw || h != ph {
			t.Errorf("an image of %g × %g at %g: a raster of %d × %d, want %d × %d", c[0], c[1], c[2], w, h, pw, ph)
		}
	}
	svg := `<svg xmlns="http://www.w3.org/2000/svg" width="1e200" height="1e200" viewBox="0 0 1e-200 1"/>`
	if _, ok := newSVGImage([]byte(svg)); !ok {
		t.Error("an image of 1e200 × 1e200 is not read")
	}
	if _, ok := newSVGImage([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="1e200" viewBox="0 0 1e-200 1e200"/>`)); ok {
		t.Error("an image of no height a float64 holds is read")
	}
}

func svgOf(t *testing.T, svg string) *svgImage {
	t.Helper()
	si, ok := newSVGImage([]byte(svg))
	if !ok {
		t.Fatal("not read")
	}
	return si
}

func TestSVGLayers(t *testing.T) {
	const head = `<svg xmlns="http://www.w3.org/2000/svg" width="8" height="8">`
	// translucent elements in each other, deeper than svgLayers
	si := svgOf(t, head+strings.Repeat(`<g opacity=".9">`, svgLayers+2)+`<rect width="8" height="8"/>`+strings.Repeat(`</g>`, svgLayers+2)+
		`<rect opacity=".5" width="4" height="4" fill="red"/></svg>`)
	r := New(bdf.NewDocument(), nil)
	pic, _ := si.raster(r, 1)
	if !warned(r, "translucent") {
		t.Errorf("%d translucent groups in each other: %v", svgLayers+2, r.Warnings())
	}
	if c := pic.levels[0].RGBAAt(6, 6); c.A != 0 {
		t.Errorf("what is in a layer past the limit is drawn: %v", c)
	}
	if c := pic.levels[0].RGBAAt(1, 1); c != (color.RGBA{128, 0, 0, 128}) {
		t.Errorf("the element after them: %v", c)
	}
	// the pixels of all the layers of a raster
	si = svgOf(t, head+strings.Repeat(`<rect opacity=".5" width="8" height="8"/>`, 6)+`</svg>`)
	r = New(bdf.NewDocument(), nil)
	s := newSurface(8, 8)
	sr := &svgRenderer{r: r, doc: si.doc, budget: svgBudget, layerPixels: svgLayerPixels - 3*8*8}
	sr.render(s, identity, 8, 8)
	if !warned(r, "translucent") || sr.layerPixels != svgLayerPixels {
		t.Errorf("six layers where three are left: %d pixels, %v", sr.layerPixels, r.Warnings())
	}
	if a := s.pix[3]; a != 0.875 {
		t.Errorf("three layers of half opacity: alpha %v", a)
	}
	// the shapes of clip paths count as elements drawn
	si = svgOf(t, head+`<clipPath id="c">`+strings.Repeat(`<rect width="1" height="1"/>`, 50)+`</clipPath>`+
		strings.Repeat(`<rect clip-path="url(#c)" width="8" height="8"/>`, 50)+`</svg>`)
	r = New(bdf.NewDocument(), nil)
	sr = &svgRenderer{r: r, doc: si.doc, budget: 300}
	sr.render(newSurface(8, 8), identity, 8, 8)
	if !warned(r, "too many elements") {
		t.Errorf("50 elements with clip paths of 50 shapes where 300 are the limit: %v", r.Warnings())
	}
	// SVG images in SVG images
	for _, c := range []struct {
		deep  int
		drawn uint8
	}{{svgInside, 255}, {svgInside + 2, 0}} {
		svg := head + `<rect width="8" height="8"/></svg>`
		for i := 0; i < c.deep; i++ {
			svg = head + `<image width="8" height="8" href="data:image/svg+xml;base64,` + base64Of([]byte(svg)) + `"/></svg>`
		}
		r = New(bdf.NewDocument(), nil)
		if pic, _ := svgOf(t, svg).raster(r, 1); warned(r, "deep in SVG images") != (c.drawn == 0) || pic.levels[0].RGBAAt(4, 4).A != c.drawn {
			t.Errorf("images %d deep: alpha %d, %v", c.deep, pic.levels[0].RGBAAt(4, 4).A, r.Warnings())
		}
	}
}

// TestSVGWarningsOfReading checks that what the reader of an SVG image
// left out (image/svg has the limits) is a warning of the drawing.
func TestSVGWarningsOfReading(t *testing.T) {
	const head = `<svg xmlns="http://www.w3.org/2000/svg">`
	d := bdf.NewDocument()
	o := bdf.NewObject()
	o.Image(o.AddImage(d.AddImage([]byte(head+strings.Repeat(`<g>`, svg.DefaultMaxDepth+1)))), 0, 0, 4, 4)
	hash, _ := d.AddObject(o)
	v := d.NewView("v", bdf.ViewFixed, "")
	v.AddPage(4, 4, bdf.Layer{Role: bdf.RoleBody, Obj: hash})
	r := New(d, nil)
	if _, err := r.Page(v, 0, 1); err != nil || !warned(r, "deep") {
		t.Errorf("%v, warnings %v", err, r.Warnings())
	}
}

func TestSheetsPastTheLimits(t *testing.T) {
	d := bdf.NewDocument()
	o := bdf.NewObject()
	o.FillColor(0xff0000ff).FillRect(0, 0, 4, 4)
	h, _ := d.AddObject(o)
	sheet := func(tile float32, cols []bdf.Run) *bdf.View {
		return &bdf.View{Kind: bdf.ViewSheet, Tile: tile, Gridlines: true, Cols: cols, Rows: []bdf.Run{{4, 8}},
			Tiles: map[string]string{"2,1": h.String(), "0,0": h.String(), "03,1": h.String(), "x": h.String()}}
	}
	draw := func(v *bdf.View, region bdf.Rect, w, h int) (*Renderer, *image.RGBA) {
		t.Helper()
		r := New(d, nil)
		start := time.Now()
		img, err := r.Region(v, 0, region, w, h)
		if e := time.Since(start); err != nil || e > 5*time.Second {
			t.Fatalf("%v, in %v", err, e)
		}
		return r, img
	}
	// tiles of no width: a region would hold any number (and take any time)
	for _, tile := range []float32{1e-6, 0.5, float32(math.Inf(1)), float32(math.NaN())} {
		if r, img := draw(sheet(tile, nil), bdf.Rect{W: 32, H: 32}, 32, 32); !warned(r, "tiles") || img.RGBAAt(1, 1) != (color.RGBA{255, 255, 255, 255}) {
			t.Errorf("tiles of %g: %v", tile, r.Warnings())
		}
	}
	// more places than tiles: 10^12 in this region
	if r, _ := draw(sheet(8, nil), bdf.Rect{W: 8e6, H: 8e6}, 16, 16); len(r.Warnings()) != 0 {
		t.Errorf("warnings %v", r.Warnings())
	}
	_, near := draw(sheet(8, nil), bdf.Rect{W: 32, H: 32}, 32, 32)
	for _, c := range []struct {
		x, y int
		want color.RGBA
	}{{1, 1, color.RGBA{255, 0, 0, 255}}, {17, 9, color.RGBA{255, 0, 0, 255}}, {9, 1, color.RGBA{255, 255, 255, 255}}, {25, 9, color.RGBA{255, 255, 255, 255}}} {
		if got := near.RGBAAt(c.x, c.y); got != c.want {
			t.Errorf("tile at (%d, %d): %v, want %v", c.x, c.y, got, c.want)
		}
	}

	// columns without end, of no width, of a width that is no number
	for _, cols := range [][]bdf.Run{{{1e30, 0}, {4, 8}}, {{float32(math.Inf(1)), 1e-9}}, {{1e9, float32(math.NaN())}}, {{float32(math.NaN()), 8}, {1e9, -1}}} {
		draw(sheet(8, cols), bdf.Rect{W: 32, H: 32}, 32, 32)
	}
	if r, img := draw(sheet(8, []bdf.Run{{1e30, 0}, {4, 8}}), bdf.Rect{W: 32, H: 32}, 32, 32); img.RGBAAt(8, 12) != (color.RGBA{0xd9, 0xd9, 0xd9, 255}) {
		t.Errorf("the line after the columns of no width: %v %v", img.RGBAAt(8, 12), r.Warnings())
	}
	if r, _ := draw(sheet(8, []bdf.Run{{1e30, 1e-12}}), bdf.Rect{W: 32, H: 32}, 32, 32); !warned(r, "rows or columns") {
		t.Errorf("10^30 columns: %v", r.Warnings())
	}
}

// TestHiddenColumns checks that the lines of hidden rows and columns, which
// are drawn where the line before was, are drawn as one after the other is.
func TestHiddenColumns(t *testing.T) {
	defer func() { plain = false }()
	d := bdf.NewDocument()
	for _, c := range []struct {
		region bdf.Rect
		w, h   int
	}{{bdf.Rect{W: 40, H: 40}, 40, 40}, {bdf.Rect{X: 3.3, Y: 2.1, W: 40, H: 33.3}, 37, 41}, {bdf.Rect{W: 27.27, H: 100}, 64, 64}} {
		v := &bdf.View{Kind: bdf.ViewSheet, Gridlines: true,
			Cols: []bdf.Run{{2, 7.3}, {300, 0}, {1, 0.004}, {2, 5.55}, {70, 0}},
			Rows: []bdf.Run{{1, 3}, {5, 0}, {3, 6.1}, {1000, 0.001}, {2, 4}}}
		var got [2]*image.RGBA
		for i := range got {
			plain = i == 1
			var err error
			if got[i], err = New(d, nil).Region(v, 0, c.region, c.w, c.h); err != nil {
				t.Fatal(err)
			}
		}
		if !bytes.Equal(got[0].Pix, got[1].Pix) {
			t.Errorf("region %v: the lines drawn once each differ from those drawn every time", c.region)
		}
	}
}

func TestPageWithoutSize(t *testing.T) {
	d := bdf.NewDocument()
	v := d.NewView("v", bdf.ViewFixed, "")
	nan := float32(math.NaN())
	for _, s := range [][2]float32{{-5, 8e7}, {0, 1e9}, {1e-10, 1e10}, {nan, 10}, {10, nan}, {10, float32(math.Inf(1))}} {
		v.Pages = []*bdf.Page{{W: s[0], H: s[1]}}
		if img, err := New(d, nil).Page(v, 0, 1); err == nil {
			t.Errorf("a page of %v × %v: an image of %v", s[0], s[1], img.Rect)
		}
	}
	v.Pages = []*bdf.Page{{W: 10, H: 10}}
	for _, scale := range []float64{0, -1, math.NaN(), math.Inf(1)} {
		if _, err := New(d, nil).Page(v, 0, scale); err == nil {
			t.Errorf("no error at scale %v", scale)
		}
	}
	// nothing of a view, a page or a document
	v.Pages = []*bdf.Page{nil, {W: 10, H: 10}}
	r := New(d, nil)
	if _, err := r.Page(v, 0, 1); err != ErrNoPage {
		t.Errorf("Page of a nil page: %v", err)
	}
	if _, err := r.Region(v, 0, bdf.Rect{W: 1, H: 1}, 4, 4); err != ErrNoPage {
		t.Errorf("Region of a nil page: %v", err)
	}
	if _, err := r.Page(nil, 0, 1); err == nil {
		t.Error("Page of a nil view")
	}
	if _, err := r.Region(nil, 0, bdf.Rect{W: 1, H: 1}, 4, 4); err == nil {
		t.Error("Region of a nil view")
	}
	v.Kind = bdf.ViewScroll
	if img, err := r.Region(v, 0, bdf.Rect{W: 10, H: 10}, 4, 4); err != nil || img.RGBAAt(1, 1).A != 255 {
		t.Errorf("a scroll view with a nil page: %v", err)
	}
	// a panic is an error
	v.Kind = bdf.ViewFixed
	v.Pages[1].Layers = []bdf.Layer{{Role: bdf.RoleBody}}
	if img, err := New(nil, nil).Page(v, 1, 1); err == nil || img != nil || !strings.Contains(err.Error(), "nil pointer") {
		t.Errorf("a renderer of no document: %v", err)
	}
}

func TestNumbersNoFloatHolds(t *testing.T) {
	// the ends of an edge so far apart that its slope is no number
	b := image.Rect(0, 0, 8, 8)
	for _, pl := range [][]point{
		{{-1e308, -1e308}, {1e308, 1e308}, {1e308, -1e308}},
		{{-1e308, 0.03125}, {1e308, 5}, {0, 8}},
		{{0, 0}, {4, 5e-324}, {8, 8}, {0, 8}},
	} {
		m := rasterize([]polyline{{pts: pl, closed: true}}, bdf.NonZero, b) // panics on amd64 without the check
		for _, v := range m.a {
			if !(v >= 0 && v <= 1) {
				t.Fatalf("coverage %v", v)
			}
		}
	}
	// an image drawn in no place at all, under a composite operation that
	// draws all the same
	var png8 bytes.Buffer
	png.Encode(&png8, image.NewRGBA(image.Rect(0, 0, 4, 4)))
	d := bdf.NewDocument()
	o := bdf.NewObject()
	ref := o.AddImage(d.AddImage(png8.Bytes()))
	o.Blend(bdf.BlendCopy)
	for i := 0; i < 10; i++ {
		o.Scale(1e-30, 1)
	}
	o.Image(ref, 0, 0, 4e-9, 4e10)
	hash, _ := d.AddObject(o)
	v := d.NewView("v", bdf.ViewFixed, "")
	v.AddPage(8, 8, bdf.Layer{Role: bdf.RoleBody, Obj: hash})
	if _, err := New(d, nil).Page(v, 0, 1); err != nil {
		t.Error(err)
	}
	p := &picture{w: 4, h: 4, levels: []*image.RGBA{image.NewRGBA(image.Rect(0, 0, 4, 4))}}
	for _, s := range []float64{math.Inf(1), 1e300, 1e10, 3, math.NaN(), 0} {
		if l := p.levelFor(matrix{s, 0, 0, s, 0, 0}); l < 0 || l >= len(p.levels) {
			t.Errorf("level %d of %d for a scale of %v", l, len(p.levels), s)
		}
	}
}

func base64Of(b []byte) string { return base64.StdEncoding.EncodeToString(b) }

func pngOf(w, h int, c color.RGBA) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < len(img.Pix); i += 4 {
		copy(img.Pix[i:], []byte{c.R, c.G, c.B, c.A})
	}
	var b bytes.Buffer
	png.Encode(&b, img)
	return b.Bytes()
}

func TestRoomForPictures(t *testing.T) {
	d := bdf.NewDocument()
	v := d.NewView("v", bdf.ViewFixed, "")
	for p := 0; p < 2; p++ {
		o := bdf.NewObject()
		for i := 0; i < 3; i++ {
			o.Image(o.AddImage(d.AddImage(pngOf(32, 32, color.RGBA{uint8(p), uint8(i), 0, 255}))), float32(4*i), 0, 4, 4)
		}
		// outside the page: not worth decoding
		o.Image(o.AddImage(d.AddImage(pngOf(32, 32, color.RGBA{uint8(p), 9, 0, 255}))), 40, 40, 4, 4)
		h, _ := d.AddObject(o)
		v.AddPage(12, 4, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	}
	r := New(d, nil)
	r.room = 2 * 32 * 32
	img, err := r.Page(v, 0, 1)
	if err != nil || !warned(r, "pixels in all") {
		t.Fatalf("three images where there is room for two: %v %v", err, r.Warnings())
	}
	if img.RGBAAt(5, 1) != (color.RGBA{0, 1, 0, 255}) || img.RGBAAt(9, 1) != (color.RGBA{255, 255, 255, 255}) {
		t.Errorf("the second image %v, the third %v", img.RGBAAt(5, 1), img.RGBAAt(9, 1))
	}
	if len(r.pictures) != 2 || r.pixels != 2*32*32 {
		t.Errorf("%d pictures of %d pixels kept", len(r.pictures), r.pixels)
	}
	// the next page has the room of the pictures of the page before
	img, err = r.Page(v, 1, 1)
	if err != nil || img.RGBAAt(5, 1) != (color.RGBA{1, 1, 0, 255}) || len(r.kept) != 2 || r.pixels != 2*32*32 {
		t.Errorf("the second page: %v, %v, %d pictures of %d pixels kept", err, img.RGBAAt(5, 1), len(r.kept), r.pixels)
	}
	// and the first page draws the same again
	if img, err = r.Page(v, 0, 1); err != nil || img.RGBAAt(5, 1) != (color.RGBA{0, 1, 0, 255}) || r.pixels != 2*32*32 {
		t.Errorf("the first page again: %v, %v, %d pixels kept", err, img.RGBAAt(5, 1), r.pixels)
	}

	// an image of an SVG image is decoded once, however often it is used
	src := `<svg xmlns="http://www.w3.org/2000/svg" width="8" height="8"><defs><image id="i" width="4" height="4" href="data:image/png;base64,` +
		base64Of(pngOf(16, 16, color.RGBA{0, 0, 255, 255})) + `"/><path id="p" d="M4 4h4v4h-4z"/></defs>` + strings.Repeat(`<use href="#i"/><use href="#p"/>`, 5) + `</svg>`
	si := svgOf(t, src)
	r = New(bdf.NewDocument(), nil)
	pic, _ := si.raster(r, 1)
	if e := si.doc.embeds[si.doc.IDs["i"]]; e == nil || e.pic == nil || len(r.kept) != 2 || si.doc.paths[si.doc.IDs["p"]] == nil || len(si.doc.embeds) != 1 || len(si.doc.paths) != 1 {
		t.Errorf("%d pictures kept, %d images and %d paths of elements", len(r.kept), len(si.doc.embeds), len(si.doc.paths))
	}
	if pic.levels[0].RGBAAt(1, 1) != (color.RGBA{0, 0, 255, 255}) || pic.levels[0].RGBAAt(5, 5) != (color.RGBA{0, 0, 0, 255}) {
		t.Errorf("drawn %v %v", pic.levels[0].RGBAAt(1, 1), pic.levels[0].RGBAAt(5, 5))
	}
}
