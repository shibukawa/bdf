package shapes_test

import (
	"errors"
	"math"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/raster/imagebdf"
	"github.com/shibukawa/bdf/raster/internal/shapes"
)

func shapesOf(t *testing.T, o *bdf.Object) []shapes.Shape {
	t.Helper()
	part, err := bdf.DecodeObject(o.Encode())
	if err != nil {
		t.Fatal(err)
	}
	s, err := shapes.Of(part)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// sample draws what a renderer of paths draws: every verb of a path, fills
// under both rules, strokes with every cap and join, transforms, an alpha.
func sample() *bdf.Object {
	o := bdf.NewObject()
	o.SetBBox(0, 0, 200, 160)
	o.FillColor(bdf.RGBA(200, 30, 30, 255))
	o.FillRect(10, 10, 40, 30)
	o.Save()
	o.Translate(60, 10)
	o.Scale(2, 1.5)
	o.FillColor(bdf.RGBA(30, 30, 200, 255))
	p := o.AddPath((&bdf.Path{}).RoundRect(0, 0, 30, 20, 6).Circle(15, 10, 5))
	o.FillPath(p, bdf.EvenOdd)
	o.Restore()
	o.Alpha(0.5)
	o.FillPath(o.AddPath((&bdf.Path{}).Ellipse(160, 30, 30, 15, 0.5, 0, 2*math.Pi, false)), bdf.NonZero)
	o.Alpha(1)
	o.StrokeColor(bdf.RGBA(0, 120, 0, 255))
	o.Line(4, bdf.CapRound, bdf.JoinRound, 10)
	o.StrokePath(o.AddPath((&bdf.Path{}).MoveTo(10, 70).LineTo(40, 100).ArcTo(70, 130, 100, 100, 20).LineTo(100, 70)))
	o.Line(3, bdf.CapSquare, bdf.JoinMiter, 4)
	o.StrokeRect(120, 70, 50, 30)
	o.Line(2, bdf.CapButt, bdf.JoinBevel, 10)
	o.StrokePath(o.AddPath((&bdf.Path{}).MoveTo(10, 140).QuadTo(40, 110, 70, 140).CubicTo(90, 160, 110, 120, 130, 140).Rect(140, 120, 20, 20).Close()))
	tri := o.AddPath((&bdf.Path{}).MoveTo(0, 0).LineTo(8, 0).LineTo(4, -8).Close())
	o.FillColor(bdf.RGBA(0, 0, 0, 255))
	o.FillPathRun(bdf.NonZero, []bdf.Glyph{{Path: tri, X: 170, Y: 150}, {Path: tri, X: 180, Y: 150}})
	o.FillPathAt(tri, bdf.NonZero, 190, 150)
	return o
}

func TestShapes(t *testing.T) {
	s := shapesOf(t, sample())
	if len(s) != 8 {
		t.Fatalf("%d shapes, want 8", len(s))
	}
	for i, sh := range s {
		for _, v := range sh.Path.Verbs {
			if v > bdf.VerbClose {
				t.Errorf("shape %d: verb %d is left", i, v)
			}
		}
	}
	// the rectangle, as it is
	if got := s[0].Path.Args; len(got) != 8 || got[0] != 10 || got[1] != 10 || got[4] != 50 || got[5] != 40 {
		t.Errorf("the rectangle: %v", got)
	}
	// the transform is applied: the rounded rectangle starts at (6, 0), which is (72, 10)
	if got := s[1].Path.Args; got[0] != 72 || got[1] != 10 || s[1].Rule != bdf.EvenOdd {
		t.Errorf("the rounded rectangle starts at (%v, %v) under rule %d", got[0], got[1], s[1].Rule)
	}
	// the restored state fills in the color from before, at half the alpha
	if s[1].Color != bdf.RGBA(30, 30, 200, 255) || s[2].Color != bdf.RGBA(200, 30, 30, 128) {
		t.Errorf("colors %08x and, restored and at half the alpha, %08x", uint32(s[1].Color), uint32(s[2].Color))
	}
	if !s[3].Stroke || s[3].Width != 4 || s[3].Cap != bdf.CapRound || s[3].Join != bdf.JoinRound {
		t.Errorf("the first line: %+v", s[3])
	}
	if !s[4].Stroke || s[4].Miter != 4 || s[5].Join != bdf.JoinBevel {
		t.Errorf("the lines: miter %v, join %d", s[4].Miter, s[5].Join)
	}
	// a run is one shape of all its paths
	if n := len(s[6].Path.Verbs); n != 8 {
		t.Errorf("the run of two triangles has %d verbs", n)
	}
}

// A line is as wide as the transform makes it.
func TestStrokeScale(t *testing.T) {
	o := bdf.NewObject()
	o.Scale(3, 3)
	o.Line(2, 0, 0, 10)
	o.StrokeRect(0, 0, 10, 10)
	if s := shapesOf(t, o); len(s) != 1 || s[0].Width != 6 || s[0].Path.Args[2] != 30 {
		t.Errorf("shapes %+v", s)
	}
}

func TestUnsupported(t *testing.T) {
	for name, draw := range map[string]func(o *bdf.Object){
		"FILL_TEXT":   func(o *bdf.Object) { o.FillText("x", 0, 0, 0) },
		"CLIP_RECT":   func(o *bdf.Object) { o.ClipRect(0, 0, 1, 1) },
		"GROUP_BEGIN": func(o *bdf.Object) { o.GroupBegin(0.5, 0, 0, 0, 1, 1).GroupEnd() },
		"DASH":        func(o *bdf.Object) { o.Dash([]float32{2, 1}, 0) },
		"BLEND":       func(o *bdf.Object) { o.Blend(bdf.BlendMultiply) },
	} {
		o := bdf.NewObject()
		o.FillRect(0, 0, 1, 1)
		draw(o)
		part, err := bdf.DecodeObject(o.Encode())
		if err != nil {
			t.Fatal(err)
		}
		var ue *shapes.UnsupportedError
		if _, err := shapes.Of(part); !errors.As(err, &ue) || bdf.OpName(ue.Op) != name {
			t.Errorf("%s: %v", name, err)
		}
	}
	// what sets a state that paths do not use, and what draws nothing
	o := bdf.NewObject()
	o.Dash(nil, 0).Blend(bdf.BlendSourceOver).TextStyle(0, 0, 0, 0).Mark(bdf.MarkAltText, "x").Link(0, 0, 1, 1, "https://example.com/")
	o.FillRect(0, 0, 1, 1)
	if s := shapesOf(t, o); len(s) != 1 {
		t.Errorf("%d shapes", len(s))
	}
}

// The shapes draw what the object draws: the sample drawn by imagebdf, and
// an object of its shapes drawn by imagebdf.
func TestSameAsImage(t *testing.T) {
	o := sample()
	again := bdf.NewObject()
	again.SetBBox(o.BBox.X, o.BBox.Y, o.BBox.W, o.BBox.H)
	for _, s := range shapesOf(t, o) {
		p := s.Path
		if s.Stroke {
			again.StrokeColor(s.Color)
			again.Line(float32(s.Width), s.Cap, s.Join, float32(s.Miter))
			again.StrokePath(again.AddPath(&p))
		} else {
			again.FillColor(s.Color)
			again.FillPath(again.AddPath(&p), s.Rule)
		}
	}
	a, err := imagebdf.Object(o, 2, nil)
	if err != nil {
		t.Fatal(err)
	}
	b, err := imagebdf.Object(again, 2, nil)
	if err != nil {
		t.Fatal(err)
	}
	if a.Bounds() != b.Bounds() {
		t.Fatalf("bounds %v and %v", a.Bounds(), b.Bounds())
	}
	differ, drawn := 0, 0
	for i := range a.Pix {
		if d := int(a.Pix[i]) - int(b.Pix[i]); d > 2 || d < -2 {
			differ++
		}
		if i%4 == 3 && a.Pix[i] > 0 {
			drawn++
		}
	}
	if differ > 0 || drawn < 5000 {
		t.Errorf("%d of %d values differ (%d pixels drawn)", differ, len(a.Pix), drawn)
	}
}
