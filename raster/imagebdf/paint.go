package imagebdf

import (
	"math"
	"sort"

	"github.com/shibukawa/bdf"
)

// style is a fill or stroke style: a colour, or a gradient or pattern
// captured with the image it draws (it outlives the object that set it).
type style struct {
	color bdf.Color
	paint *bdf.Paint
	image *picture  // pattern image; nil when it could not be decoded
	svg   *svgImage // or pattern SVG image, drawn at the scale it is used at
}

// gradientLUT is a gradient's colours at 1024 steps from 0 to 1, premultiplied.
type gradientLUT [1025][4]float32

func buildLUT(stops []bdf.Stop) *gradientLUT {
	type stop struct {
		off float64
		c   [4]float64 // unpremultiplied
	}
	ss := make([]stop, 0, len(stops))
	for _, s := range stops {
		c := uint32(s.Color)
		ss = append(ss, stop{math.Min(1, math.Max(0, float64(s.Offset))), [4]float64{float64(c>>24) / 255, float64(c>>16&0xff) / 255, float64(c>>8&0xff) / 255, float64(c&0xff) / 255}})
	}
	// addColorStop keeps stops with the same offset in the order they came
	sort.SliceStable(ss, func(i, j int) bool { return ss[i].off < ss[j].off })
	lut := &gradientLUT{}
	// k is the first stop from the second that t does not pass: it moves
	// on as t grows
	k := 1
	for i := range lut {
		t := float64(i) / 1024
		var c [4]float64
		switch {
		case t <= ss[0].off:
			c = ss[0].c
		case t >= ss[len(ss)-1].off:
			c = ss[len(ss)-1].c
		default:
			for k < len(ss) && !(t <= ss[k].off) {
				k++
			}
			if k < len(ss) {
				a, b := ss[k-1], ss[k]
				u := 0.0
				if b.off > a.off {
					u = (t - a.off) / (b.off - a.off)
				}
				for j := range c {
					c[j] = a.c[j] + (b.c[j]-a.c[j])*u
				}
			}
		}
		lut[i] = [4]float32{float32(c[0] * c[3]), float32(c[1] * c[3]), float32(c[2] * c[3]), float32(c[3])}
	}
	return lut
}

func (l *gradientLUT) at(t float64) [4]float32 {
	switch {
	case t <= 0 || math.IsNaN(t):
		return l[0]
	case t >= 1:
		return l[1024]
	}
	return l[int(t*1024+0.5)]
}

// gradientShader samples a gradient in user space through the inverse of
// the transform it is drawn with.
type gradientShader struct {
	inv matrix
	p   *bdf.Paint
	lut *gradientLUT
}

func newGradientShader(p *bdf.Paint, m matrix, lut *gradientLUT) shader {
	inv, ok := m.invert()
	if !ok || len(p.Stops) == 0 {
		return solid{}
	}
	if lut == nil {
		lut = buildLUT(p.Stops)
	}
	return &gradientShader{inv: inv, p: p, lut: lut}
}

// lutsKept is the most gradients whose colours a renderer keeps (16 KiB
// each).
const lutsKept = 1024

// lut returns the colours of a gradient of an object, which are worked out
// once for all it fills.
func (r *Renderer) lut(p *bdf.Paint) *gradientLUT {
	if len(p.Stops) == 0 || plain {
		return nil
	}
	l, ok := r.luts[p]
	if !ok {
		if len(r.luts) >= lutsKept {
			clear(r.luts)
		}
		l = buildLUT(p.Stops)
		r.luts[p] = l
	}
	return l
}

func (g *gradientShader) shade(x, y, n int, dst []float32) {
	c := g.p.Coords
	for i := 0; i < n; i++ {
		ux, uy := g.inv.apply(float64(x+i)+0.5, float64(y)+0.5)
		var col [4]float32
		switch g.p.Kind {
		case bdf.PaintLinear:
			x0, y0, x1, y1 := float64(c[0]), float64(c[1]), float64(c[2]), float64(c[3])
			dx, dy := x1-x0, y1-y0
			if l := dx*dx + dy*dy; l > 0 {
				col = g.lut.at(((ux-x0)*dx + (uy-y0)*dy) / l)
			}
		case bdf.PaintRadial:
			if t, ok := radialT(ux, uy, c); ok {
				col = g.lut.at(t)
			}
		case bdf.PaintConic:
			a := math.Atan2(uy-float64(c[2]), ux-float64(c[1])) - float64(c[0])
			t := math.Mod(a/(2*math.Pi), 1)
			if t < 0 {
				t++
			}
			col = g.lut.at(t)
		}
		copy(dst[4*i:4*i+4], col[:])
	}
}

// radialT solves the two-circle gradient of Canvas for the largest ω whose
// circle passes through (x, y) with a radius that is not negative.
func radialT(x, y float64, c []float32) (float64, bool) {
	x0, y0, r0 := float64(c[0]), float64(c[1]), float64(c[2])
	x1, y1, r1 := float64(c[3]), float64(c[4]), float64(c[5])
	cdx, cdy, dr := x1-x0, y1-y0, r1-r0
	pdx, pdy := x-x0, y-y0
	a := cdx*cdx + cdy*cdy - dr*dr
	b := pdx*cdx + pdy*cdy + r0*dr
	cc := pdx*pdx + pdy*pdy - r0*r0
	ok := func(w float64) bool { return r0+w*dr >= 0 }
	if math.Abs(a) < 1e-12 {
		if b == 0 {
			return 0, false
		}
		w := cc / (2 * b)
		return w, ok(w)
	}
	disc := b*b - a*cc
	if disc < 0 {
		return 0, false
	}
	sq := math.Sqrt(disc)
	w1, w2 := (b+sq)/a, (b-sq)/a
	if w1 < w2 {
		w1, w2 = w2, w1
	}
	if ok(w1) {
		return w1, true
	}
	if ok(w2) {
		return w2, true
	}
	return 0, false
}

// patternShader repeats an image in the pattern space of its paint.
type patternShader struct {
	inv    matrix // device → image pixels
	img    *picture
	repeat byte
	level  int
}

func newPatternShader(r *Renderer, st *style, m matrix) shader {
	pm := st.paint.Matrix
	full := m.mul(matrix{float64(pm[0]), float64(pm[1]), float64(pm[2]), float64(pm[3]), float64(pm[4]), float64(pm[5])})
	img := st.image
	if st.svg != nil {
		// an SVG cell as large as it is drawn
		img, _ = st.svg.raster(r, full.maxScale())
		full = full.scale(st.svg.w/float64(img.w), st.svg.h/float64(img.h))
	}
	if img == nil {
		return solid{}
	}
	inv, ok := full.invert()
	if !ok {
		return solid{}
	}
	return &patternShader{inv: inv, img: img, repeat: st.paint.Repeat, level: img.levelFor(inv)}
}

func (p *patternShader) shade(x, y, n int, dst []float32) {
	w, h := float64(p.img.w), float64(p.img.h)
	for i := 0; i < n; i++ {
		ix, iy := p.inv.apply(float64(x+i)+0.5, float64(y)+0.5)
		var col [4]float32
		inX, inY := ix >= 0 && ix < w, iy >= 0 && iy < h
		switch p.repeat {
		case bdf.RepeatBoth:
			ix, iy, inX, inY = wrap(ix, w), wrap(iy, h), true, true
		case bdf.RepeatX:
			ix, inX = wrap(ix, w), true
		case bdf.RepeatY:
			iy, inY = wrap(iy, h), true
		}
		if inX && inY {
			col = p.img.sample(p.level, ix, iy, p.repeat == bdf.RepeatBoth || p.repeat == bdf.RepeatX, p.repeat == bdf.RepeatBoth || p.repeat == bdf.RepeatY)
		}
		copy(dst[4*i:4*i+4], col[:])
	}
}

func wrap(v, n float64) float64 {
	v = math.Mod(v, n)
	if v < 0 {
		v += n
	}
	return v
}

// imageShader draws an image (or a part of it) into a destination
// rectangle; the rectangle's outline is the coverage.
type imageShader struct {
	inv    matrix // device → image pixels
	img    *picture
	level  int
	smooth bool
	// the source rectangle sampled from, in image pixels
	sx0, sy0, sx1, sy1 float64
}

func (s *imageShader) shade(x, y, n int, dst []float32) {
	for i := 0; i < n; i++ {
		ix, iy := s.inv.apply(float64(x+i)+0.5, float64(y)+0.5)
		ix = math.Min(math.Max(ix, s.sx0), s.sx1)
		iy = math.Min(math.Max(iy, s.sy0), s.sy1)
		var col [4]float32
		if s.smooth || s.level > 0 {
			col = s.img.sample(s.level, ix, iy, false, false)
		} else {
			col = s.img.nearest(ix, iy)
		}
		copy(dst[4*i:4*i+4], col[:])
	}
}
