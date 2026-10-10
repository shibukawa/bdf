package ebitenginebdf

import (
	"errors"
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/image/formula"
)

// The tests compile objects and look at the paths; drawing needs a running
// game (see example, which saves what it draws with -shot).

func TestCompile(t *testing.T) {
	o := bdf.NewObject()
	o.FillColor(bdf.RGBA(255, 0, 0, 128))
	o.FillRect(10, 20, 30, 40)
	o.StrokeColor(bdf.RGBA(0, 0, 255, 255))
	o.Line(3, bdf.CapRound, bdf.JoinBevel, 4)
	o.StrokePath(o.AddPath((&bdf.Path{}).Circle(100, 100, 50)))
	o.FillPath(o.AddPath((&bdf.Path{}).Rect(0, 0, 8, 8).Rect(2, 2, 4, 4)), bdf.EvenOdd)
	d, err := Compile(o)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.shapes) != 3 {
		t.Fatalf("%d shapes, want 3", len(d.shapes))
	}
	rect, circle, holed := &d.shapes[0], &d.shapes[1], &d.shapes[2]
	if b := rect.path.Bounds(); b != image.Rect(10, 20, 40, 60) {
		t.Errorf("the rectangle is %v", b)
	}
	// half transparent red, premultiplied
	if rect.stroke || rect.r < 0.5 || rect.r > 0.51 || rect.g != 0 || rect.a < 0.5 || rect.a > 0.51 {
		t.Errorf("the rectangle is drawn in (%v, %v, %v, %v)", rect.r, rect.g, rect.b, rect.a)
	}
	if b := circle.path.Bounds(); b != image.Rect(50, 50, 150, 150) {
		t.Errorf("the circle is %v", b)
	}
	if want := (vector.StrokeOptions{Width: 3, LineCap: vector.LineCapRound, LineJoin: vector.LineJoinBevel, MiterLimit: 4}); !circle.stroke || circle.line != want {
		t.Errorf("the circle is stroked with %+v", circle.line)
	}
	if holed.fill.FillRule != vector.FillRuleEvenOdd || rect.fill.FillRule != vector.FillRuleNonZero {
		t.Errorf("fill rules %v and %v", rect.fill.FillRule, holed.fill.FillRule)
	}
	if d.Bounds() != o.BBox {
		t.Errorf("bounds %v, want %v", d.Bounds(), o.BBox)
	}
}

// The paths are transformed once for the drawings at one place.
func TestPlace(t *testing.T) {
	o := bdf.NewObject()
	o.FillRect(0, 0, 10, 10)
	d, err := Compile(o)
	if err != nil {
		t.Fatal(err)
	}
	var g ebiten.GeoM
	g.Scale(2, 3)
	g.Translate(100, 200)
	p := d.place(g)
	if b := p[0].Bounds(); b != image.Rect(100, 200, 120, 230) {
		t.Errorf("placed at %v", b)
	}
	if q := d.place(g); &q[0] != &p[0] {
		t.Error("the paths were made again for the same transform")
	}
	g.Translate(1, 0)
	if b := d.place(g)[0].Bounds(); b != image.Rect(101, 200, 121, 230) {
		t.Errorf("moved to %v", b)
	}
}

func TestFormula(t *testing.T) {
	ts, err := formula.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	l := ts.Layout(formula.ParseTeX(`x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}`), formula.Style{Size: 40, Display: true})
	d, err := Compile(l.Object())
	if err != nil {
		t.Fatal(err)
	}
	if len(d.shapes) < 3 {
		t.Fatalf("%d shapes", len(d.shapes))
	}
	// the glyphs reach above and below the baseline, within the ink
	var all image.Rectangle
	for i := range d.shapes {
		all = all.Union(d.shapes[i].path.Bounds())
	}
	ink := l.Ink()
	box := image.Rect(int(ink.X)-1, int(ink.Y)-1, int(ink.X+ink.W)+2, int(ink.Y+ink.H)+2)
	if all.Min.Y >= 0 || all.Max.Y <= 0 || !all.In(box) || all.Dx() < box.Dx()*9/10 || all.Dy() < box.Dy()*9/10 {
		t.Errorf("the paths cover %v, the ink is %v", all, ink)
	}
}

func TestErrors(t *testing.T) {
	if _, err := Compile(nil); err == nil {
		t.Error("no object: no error")
	}
	text := bdf.NewObject()
	text.FillText("x", 0, 0, 0)
	var ue *UnsupportedError
	if _, err := Compile(text); !errors.As(err, &ue) || ue.Op != bdf.OpFillText {
		t.Errorf("text: %v", err)
	}
	img := bdf.NewObject()
	img.Image(img.AddImage(bdf.Hash{1}), 0, 0, 1, 1)
	if _, err := Compile(img); !errors.Is(err, ErrNeedsDocument) {
		t.Errorf("an image: %v", err)
	}
}
