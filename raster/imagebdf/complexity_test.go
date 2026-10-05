package imagebdf

import (
	"testing"

	"github.com/shibukawa/bdf"
)

func TestFlattenBudgetAcrossPaths(t *testing.T) {
	p := &path{}
	p.moveTo(0, 0)
	for i := 0; i < 2100; i++ {
		p.cubicTo(1e6, 0, 0, 1e6, 8, 8)
	}
	sc := &scratch{}
	polys := p.flattenTo(sc, nil, identity, flatTol)
	if sc.limited || len(sc.pts) < maxPathPoints/2 {
		t.Fatalf("first path: limited %v, points %d", sc.limited, len(sc.pts))
	}
	p.flattenTo(sc, polys, identity, flatTol)
	if !sc.limited || len(sc.pts) > maxPathPoints {
		t.Fatalf("two paths: limited %v, points %d", sc.limited, len(sc.pts))
	}
	sc.reset()
	p = &path{}
	p.rect(0, 0, 4, 4)
	if p.flattenTo(sc, nil, identity, flatTol) == nil || sc.limited {
		t.Fatal("budget did not reset for the next instruction")
	}
}

func TestComplexPathsAreSkippedAndDrawingContinues(t *testing.T) {
	curves := &bdf.Path{}
	curves.MoveTo(0, 0)
	for i := 0; i < 4300; i++ {
		curves.CubicTo(1e6, 0, 0, 1e6, 8, 8)
	}
	for _, mode := range []string{"fill", "stroke", "scan"} {
		t.Run(mode, func(t *testing.T) {
			o := bdf.NewObject()
			o.FillColor(0xff0000ff).StrokeColor(0xff0000ff)
			p := curves
			if mode == "scan" {
				// Every one of 512 sampled scanlines crosses 36,000 edges:
				// the point count is small but the scan work is excessive.
				p = slivers(18000, 32, 32)
			}
			ref := o.AddPath(p)
			if mode == "stroke" {
				o.StrokePath(ref)
			} else {
				o.FillPath(ref, bdf.NonZero)
			}
			o.FillColor(0x0000ffff).FillRect(28, 28, 4, 4)
			r, at := drawn(t, o, 32, 32)
			if !warned(r, "too complex") || at(1, 1) != white || at(29, 29) != blue {
				t.Errorf("warnings %v, skipped path %v, later rectangle %v", r.Warnings(), at(1, 1), at(29, 29))
			}
		})
	}
}
