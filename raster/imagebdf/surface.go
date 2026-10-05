package imagebdf

import (
	"image"
	"math"

	"github.com/shibukawa/bdf"
)

// surface is a canvas of premultiplied RGBA pixels, 0 to 1.
type surface struct {
	w, h int
	pix  []float32
	row  []float32 // the colours of a row being filled
}

func newSurface(w, h int) *surface {
	return &surface{w: w, h: h, pix: make([]float32, 4*w*h)}
}

func (s *surface) bounds() image.Rectangle { return image.Rect(0, 0, s.w, s.h) }

// shader gives the premultiplied colours of a run of pixels: those of
// [x, x+n) in row y, 4 floats each, into dst.
type shader interface {
	shade(x, y, n int, dst []float32)
}

// solid is a shader of one premultiplied colour.
type solid [4]float32

func (c solid) shade(x, y, n int, dst []float32) {
	for i := 0; i < n; i++ {
		copy(dst[4*i:4*i+4], c[:])
	}
}

// premul converts a packed 0xRRGGBBAA colour.
func premul(c bdf.Color) solid {
	r, g, b, a := float32(uint32(c)>>24)/255, float32(uint32(c)>>16&0xff)/255, float32(uint32(c)>>8&0xff)/255, float32(uint32(c)&0xff)/255
	return solid{r * a, g * a, b * a, a}
}

// unbounded reports the composite operations that change pixels the source
// does not cover (Canvas applies them to the whole clip region).
func unbounded(mode byte) bool {
	switch mode {
	case bdf.BlendSourceIn, bdf.BlendSourceOut, bdf.BlendDestinationIn, bdf.BlendDestinationAtop, bdf.BlendCopy:
		return true
	}
	return false
}

// fill composites a shader through a coverage mask, scaled by alpha, with
// a composite operation, inside a clip. A nil cov covers everything.
func (s *surface) fill(cov *mask, sh shader, alpha float32, mode byte, clip *mask) {
	if clip.empty() || alpha <= 0 && !unbounded(mode) {
		return
	}
	r := clip.r.Intersect(s.bounds())
	if cov != nil && !unbounded(mode) {
		r = r.Intersect(cov.r)
	}
	if r.Empty() {
		return
	}
	if unbounded(mode) || plain {
		s.fillAll(cov, sh, alpha, mode, clip, r)
		return
	}
	// the pixels are all within the coverage: those of a row are told
	// without asking where each one is
	n := r.Dx()
	if cap(s.row) < 4*n {
		s.row = make([]float32, 4*n)
	}
	buf := s.row[:4*n]
	_, flat := sh.(solid)
	if flat {
		sh.shade(r.Min.X, r.Min.Y, n, buf) // the same in every row
	}
	for y := r.Min.Y; y < r.Max.Y; y++ {
		if !flat {
			sh.shade(r.Min.X, y, n, buf)
		}
		krow := clip.row(y, r.Min.X, r.Max.X)
		var crow []float32
		if cov != nil {
			crow = cov.row(y, r.Min.X, r.Max.X)
		}
		d := s.pix[4*(y*s.w+r.Min.X) : 4*(y*s.w+r.Max.X)]
		for i := 0; i < n; i++ {
			c := alpha
			if crow != nil {
				c = crow[i] * alpha
			}
			k := float32(1)
			if krow != nil {
				k = krow[i]
			}
			if k == 0 || c == 0 {
				continue
			}
			src := [4]float32{buf[4*i] * c, buf[4*i+1] * c, buf[4*i+2] * c, buf[4*i+3] * c}
			dst := [4]float32{d[4*i], d[4*i+1], d[4*i+2], d[4*i+3]}
			var out [4]float32
			if mode == bdf.BlendSourceOver {
				ia := 1 - src[3]
				out = [4]float32{src[0] + dst[0]*ia, src[1] + dst[1]*ia, src[2] + dst[2]*ia, src[3] + dst[3]*ia}
			} else {
				out = composite(mode, src, dst)
			}
			if k == 1 {
				copy(d[4*i:4*i+4], out[:])
				continue
			}
			for j := 0; j < 4; j++ {
				d[4*i+j] = dst[j] + k*(out[j]-dst[j])
			}
		}
	}
}

// fillBackground draws the page or sheet background. Integer-aligned
// rectangles need no shader row or per-pixel coverage/clip checks. Keep
// this outside fill so small paths pay no cost for this common large fill.
func (s *surface) fillBackground(cov *mask, c solid) {
	if plain || cov == nil || cov.a != nil {
		s.fill(cov, c, 1, bdf.BlendSourceOver, &mask{r: s.bounds()})
		return
	}
	r := cov.r.Intersect(s.bounds())
	if r.Empty() {
		return
	}
	// Reserve the row for later draws, as the ordinary background fill
	// does; otherwise progressively wider text runs keep reallocating it.
	if cap(s.row) < 4*r.Dx() {
		s.row = make([]float32, 4*r.Dx())
	}
	ia := 1 - c[3]
	for y := r.Min.Y; y < r.Max.Y; y++ {
		d := s.pix[4*(y*s.w+r.Min.X) : 4*(y*s.w+r.Max.X)]
		for i := 0; i < len(d); i += 4 {
			p := d[i : i+4]
			p[0] = c[0] + p[0]*ia
			p[1] = c[1] + p[1]*ia
			p[2] = c[2] + p[2]*ia
			p[3] = c[3] + p[3]*ia
		}
	}
}

// fillAll is fill for the pixels r of the clip, inside the coverage or not:
// the composite operations that change pixels the source does not cover.
func (s *surface) fillAll(cov *mask, sh shader, alpha float32, mode byte, clip *mask, r image.Rectangle) {
	n := r.Dx()
	buf := make([]float32, 4*n)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		sh.shade(r.Min.X, y, n, buf)
		krow := clip.row(y, r.Min.X, r.Max.X)
		var crow []float32
		covered := func(x int) float32 { return 1 }
		if cov != nil {
			if y >= cov.r.Min.Y && y < cov.r.Max.Y {
				lo, hi := max(r.Min.X, cov.r.Min.X), min(r.Max.X, cov.r.Max.X)
				if lo < hi {
					crow = cov.row(y, lo, hi)
					covered = func(x int) float32 {
						if x < lo || x >= hi {
							return 0
						}
						if crow == nil {
							return 1
						}
						return crow[x-lo]
					}
				} else {
					covered = func(int) float32 { return 0 }
				}
			} else {
				covered = func(int) float32 { return 0 }
			}
		}
		d := s.pix[4*(y*s.w+r.Min.X) : 4*(y*s.w+r.Max.X)]
		for i := 0; i < n; i++ {
			x := r.Min.X + i
			c := covered(x) * alpha
			k := float32(1)
			if krow != nil {
				k = krow[i]
			}
			if k == 0 || c == 0 && !unbounded(mode) {
				continue
			}
			src := [4]float32{buf[4*i] * c, buf[4*i+1] * c, buf[4*i+2] * c, buf[4*i+3] * c}
			dst := [4]float32{d[4*i], d[4*i+1], d[4*i+2], d[4*i+3]}
			var out [4]float32
			if mode == bdf.BlendSourceOver {
				ia := 1 - src[3]
				out = [4]float32{src[0] + dst[0]*ia, src[1] + dst[1]*ia, src[2] + dst[2]*ia, src[3] + dst[3]*ia}
			} else {
				out = composite(mode, src, dst)
			}
			if k == 1 {
				copy(d[4*i:4*i+4], out[:])
				continue
			}
			for j := 0; j < 4; j++ {
				d[4*i+j] = dst[j] + k*(out[j]-dst[j])
			}
		}
	}
}

// composite applies a composite operation or blend mode to premultiplied
// source and destination colours (W3C Compositing and Blending).
func composite(mode byte, s, d [4]float32) [4]float32 {
	as, ab := s[3], d[3]
	pd := func(fs, fd float32) [4]float32 {
		return [4]float32{s[0]*fs + d[0]*fd, s[1]*fs + d[1]*fd, s[2]*fs + d[2]*fd, s[3]*fs + d[3]*fd}
	}
	switch mode {
	case bdf.BlendSourceOver:
		return pd(1, 1-as)
	case bdf.BlendDestinationOver:
		return pd(1-ab, 1)
	case bdf.BlendSourceIn:
		return pd(ab, 0)
	case bdf.BlendDestinationIn:
		return pd(0, as)
	case bdf.BlendSourceOut:
		return pd(1-ab, 0)
	case bdf.BlendDestinationOut:
		return pd(0, 1-as)
	case bdf.BlendSourceAtop:
		return pd(ab, 1-as)
	case bdf.BlendDestinationAtop:
		return pd(1-ab, as)
	case bdf.BlendXor:
		return pd(1-ab, 1-as)
	case bdf.BlendCopy:
		return s
	case bdf.BlendLighter:
		return [4]float32{min(1, s[0]+d[0]), min(1, s[1]+d[1]), min(1, s[2]+d[2]), min(1, s[3]+d[3])}
	}
	// blend modes, composited source-over
	var cs, cb [3]float32
	for i := 0; i < 3; i++ {
		if as > 0 {
			cs[i] = s[i] / as
		}
		if ab > 0 {
			cb[i] = d[i] / ab
		}
	}
	var b [3]float32
	if mode >= bdf.BlendHue && mode <= bdf.BlendLuminosity {
		b = nonSeparable(mode, cb, cs)
	} else {
		for i := 0; i < 3; i++ {
			b[i] = separable(mode, cb[i], cs[i])
		}
	}
	var out [4]float32
	for i := 0; i < 3; i++ {
		out[i] = s[i]*(1-ab) + d[i]*(1-as) + as*ab*b[i]
	}
	out[3] = as + ab*(1-as)
	return out
}

func separable(mode byte, cb, cs float32) float32 {
	switch mode {
	case bdf.BlendMultiply:
		return cb * cs
	case bdf.BlendScreen:
		return cb + cs - cb*cs
	case bdf.BlendOverlay:
		return hardLight(cs, cb)
	case bdf.BlendDarken:
		return min(cb, cs)
	case bdf.BlendLighten:
		return max(cb, cs)
	case bdf.BlendColorDodge:
		switch {
		case cb == 0:
			return 0
		case cs >= 1:
			return 1
		}
		return min(1, cb/(1-cs))
	case bdf.BlendColorBurn:
		switch {
		case cb >= 1:
			return 1
		case cs <= 0:
			return 0
		}
		return 1 - min(1, (1-cb)/cs)
	case bdf.BlendHardLight:
		return hardLight(cb, cs)
	case bdf.BlendSoftLight:
		if cs <= 0.5 {
			return cb - (1-2*cs)*cb*(1-cb)
		}
		var dd float32
		if cb <= 0.25 {
			dd = ((16*cb-12)*cb + 4) * cb
		} else {
			dd = float32(math.Sqrt(float64(cb)))
		}
		return cb + (2*cs-1)*(dd-cb)
	case bdf.BlendDifference:
		if cb > cs {
			return cb - cs
		}
		return cs - cb
	case bdf.BlendExclusion:
		return cb + cs - 2*cb*cs
	}
	return cs
}

func hardLight(cb, cs float32) float32 {
	if cs <= 0.5 {
		return cb * 2 * cs
	}
	s := 2*cs - 1
	return cb + s - cb*s
}

func lum(c [3]float32) float32 { return 0.3*c[0] + 0.59*c[1] + 0.11*c[2] }

func clipColor(c [3]float32) [3]float32 {
	l := lum(c)
	n := min(c[0], c[1], c[2])
	x := max(c[0], c[1], c[2])
	for i := range c {
		if n < 0 && l-n != 0 {
			c[i] = l + (c[i]-l)*l/(l-n)
		}
		if x > 1 && x-l != 0 {
			c[i] = l + (c[i]-l)*(1-l)/(x-l)
		}
	}
	return c
}

func setLum(c [3]float32, l float32) [3]float32 {
	d := l - lum(c)
	return clipColor([3]float32{c[0] + d, c[1] + d, c[2] + d})
}

func sat(c [3]float32) float32 { return max(c[0], c[1], c[2]) - min(c[0], c[1], c[2]) }

func setSat(c [3]float32, s float32) [3]float32 {
	idx := [3]int{0, 1, 2}
	// order the channels: min, mid, max
	if c[idx[0]] > c[idx[1]] {
		idx[0], idx[1] = idx[1], idx[0]
	}
	if c[idx[1]] > c[idx[2]] {
		idx[1], idx[2] = idx[2], idx[1]
	}
	if c[idx[0]] > c[idx[1]] {
		idx[0], idx[1] = idx[1], idx[0]
	}
	var out [3]float32
	mn, md, mx := c[idx[0]], c[idx[1]], c[idx[2]]
	if mx > mn {
		out[idx[1]] = (md - mn) * s / (mx - mn)
		out[idx[2]] = s
	}
	return out
}

func nonSeparable(mode byte, cb, cs [3]float32) [3]float32 {
	switch mode {
	case bdf.BlendHue:
		return setLum(setSat(cs, sat(cb)), lum(cb))
	case bdf.BlendSaturation:
		return setLum(setSat(cb, sat(cs)), lum(cb))
	case bdf.BlendColor:
		return setLum(cs, lum(cb))
	default: // luminosity
		return setLum(cb, lum(cs))
	}
}

// surfaceShader reads the pixels of another surface placed at (dx, dy).
type surfaceShader struct {
	s      *surface
	dx, dy int
}

func (ss surfaceShader) shade(x, y, n int, dst []float32) {
	sy := y - ss.dy
	if sx := x - ss.dx; sx >= 0 && sy >= 0 && sx+n <= ss.s.w && sy < ss.s.h && !plain {
		copy(dst[:4*n], ss.s.pix[4*(sy*ss.s.w+sx):]) // the whole row is of the surface
		return
	}
	for i := 0; i < n; i++ {
		sx := x + i - ss.dx
		if sx < 0 || sy < 0 || sx >= ss.s.w || sy >= ss.s.h {
			dst[4*i], dst[4*i+1], dst[4*i+2], dst[4*i+3] = 0, 0, 0, 0
			continue
		}
		copy(dst[4*i:4*i+4], ss.s.pix[4*(sy*ss.s.w+sx):])
	}
}

// toRGBA converts the surface to an image of premultiplied 8-bit pixels.
func (s *surface) toRGBA() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, s.w, s.h))
	for i, v := range s.pix {
		// rounded to the nearest: a float32 times 255 is exact as a
		// float64, and so is the half added to what is not far below it
		img.Pix[i] = uint8(float64(min(max(v, 0), 1))*255 + 0.5)
	}
	return img
}
