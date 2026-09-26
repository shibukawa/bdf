package pdf

import (
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/shibukawa/bdf"
)

// shading is a converted PDF shading dictionary.
type shading struct {
	kind   int
	paint  bdf.Paint // for axial/radial
	avg    bdf.Color // fallback colour for unsupported types
	bbox   *rect
	hasAvg bool
}

// loadShading converts shading types 2 and 3 to gradients; other types get an
// average colour sampled from their function (or gray).
func (c *converter) loadShading(o types.Object, res types.Dict) *shading {
	p := c.pdf
	d := p.dict(o)
	if d == nil {
		return nil
	}
	sh := &shading{kind: p.intOr(d["ShadingType"], 0)}
	cs := c.loadColorSpace(d["ColorSpace"], res)
	fn := p.loadFunction(d["Function"])
	if bb := p.nums(d["BBox"]); len(bb) == 4 {
		r := rect{bb[0], bb[1], bb[2], bb[3]}
		sh.bbox = &r
	}
	coords := p.nums(d["Coords"])
	domain := p.nums(d["Domain"])
	t0, t1 := 0.0, 1.0
	if len(domain) >= 2 {
		t0, t1 = domain[0], domain[1]
	}
	colorAt := func(t float64) bdf.Color {
		if fn == nil {
			return bdf.RGB(128, 128, 128)
		}
		return cs.rgb(fn.eval(t))
	}
	stops := func() []bdf.Stop {
		var out []bdf.Stop
		add := func(s float64) {
			out = append(out, bdf.Stop{Offset: float32(s), Color: colorAt(t0 + s*(t1-t0))})
		}
		if fn != nil && fn.kind == 3 && len(fn.bounds) > 0 && allExponential1(fn) {
			// Exact stops at stitching bounds.
			add(0)
			for _, b := range fn.bounds {
				s := (b - t0) / (t1 - t0)
				out = append(out, bdf.Stop{Offset: float32(s), Color: cs.rgb(fn.eval(b - 1e-9))})
				out = append(out, bdf.Stop{Offset: float32(s), Color: cs.rgb(fn.eval(b))})
			}
			add(1)
			return out
		}
		if fn != nil && fn.kind == 2 && fn.n == 1 {
			add(0)
			add(1)
			return out
		}
		const n = 32
		for i := 0; i <= n; i++ {
			add(float64(i) / n)
		}
		return out
	}
	switch sh.kind {
	case 2:
		if len(coords) >= 4 {
			sh.paint = bdf.LinearGradient(float32(coords[0]), float32(coords[1]), float32(coords[2]), float32(coords[3]), stops()...)
			return sh
		}
	case 3:
		if len(coords) >= 6 {
			sh.paint = bdf.RadialGradient(float32(coords[0]), float32(coords[1]), float32(coords[2]), float32(coords[3]), float32(coords[4]), float32(coords[5]), stops()...)
			return sh
		}
	}
	// Fallback: a flat colour. Function-based shadings are sampled at the centre.
	sh.hasAvg = true
	if fn != nil {
		if sh.kind == 1 {
			dom := p.nums(d["Domain"])
			if len(dom) < 4 {
				dom = []float64{0, 1, 0, 1}
			}
			sh.avg = cs.rgb(fn.eval((dom[0]+dom[1])/2, (dom[2]+dom[3])/2))
		} else {
			sh.avg = colorAt((t0 + t1) / 2)
		}
	} else if bg := p.nums(d["Background"]); len(bg) > 0 {
		sh.avg = cs.rgb(bg)
	} else {
		sh.avg = bdf.RGB(128, 128, 128)
	}
	c.warnf("shading type %d approximated by a flat colour", sh.kind)
	return sh
}

func allExponential1(f *pdfFunction) bool {
	for _, sub := range f.funcs {
		if sub == nil || sub.kind != 2 || sub.n != 1 {
			return false
		}
	}
	return true
}
