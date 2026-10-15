package imagebdf

import (
	"image"
	"image/color"
	"math"

	"github.com/shibukawa/bdf"
)

// surface is a canvas of premultiplied RGBA pixels, 0 to 1.
type surface struct {
	w, h int
	pix  []float32
	row  []float32 // the colours of a row being filled
	// rows holds the colours of the rows the other goroutines of a pass
	// fill (see splitRows)
	rows [][]float32
	// fresh says that the pixels are all 0 still: those of a new canvas,
	// or of one cleared for a drawing (see fillBackground)
	fresh bool
	// threads is how many goroutines the passes over the canvas are
	// split among (see splitRows)
	threads int
}

func newSurface(w, h int) *surface {
	return &surface{w: w, h: h, pix: make([]float32, 4*w*h), fresh: true, threads: 1}
}

func (s *surface) bounds() image.Rectangle { return image.Rect(0, 0, s.w, s.h) }

// bands calls fn for bands of the rows of r, in parallel for a large r
// (see splitRows).
func (s *surface) bands(r image.Rectangle, fn func(y0, y1, k int)) {
	splitRows(s.threads, r.Min.Y, r.Max.Y, r.Dx(), fn)
}

// rowBuf returns the colours of a row for band k of a pass, n pixels wide.
// The rows of the bands after the first are made before the pass (see
// rowBufs): the goroutines of the bands share the surface.
func (s *surface) rowBuf(k, n int) []float32 {
	if k == 0 {
		if cap(s.row) < 4*n {
			s.row = make([]float32, 4*n)
		}
		return s.row[:4*n]
	}
	return s.rows[k-1][:4*n]
}

// rowBufs makes the rows of colours of the bands of a pass, n pixels wide.
func (s *surface) rowBufs(n int) {
	for len(s.rows) < s.threads-1 {
		s.rows = append(s.rows, nil)
	}
	for i := range s.rows {
		if cap(s.rows[i]) < 4*n {
			s.rows[i] = make([]float32, 4*n)
		}
	}
}

// setRows sets the pixels of r to the colours of one row, which are those
// of all of them.
func (s *surface) setRows(r image.Rectangle, row []float32) {
	s.bands(r, func(y0, y1, _ int) {
		for y := y0; y < y1; y++ {
			copy(s.pix[4*(y*s.w+r.Min.X):4*(y*s.w+r.Max.X)], row)
		}
	})
}

// repeat fills a row with one colour: the first pixel is set, and copied
// to the rest, doubling what is copied each time.
func repeat(row []float32, c solid) {
	if len(row) < 4 {
		return
	}
	copy(row, c[:])
	for n := 4; n < len(row); n *= 2 {
		copy(row[n:], row[:n])
	}
}

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
	s.fresh = false
	if unbounded(mode) || plain {
		s.fillAll(cov, sh, alpha, mode, clip, r)
		return
	}
	// the pixels are all within the coverage: those of a row are told
	// without asking where each one is
	n := r.Dx()
	buf := s.rowBuf(0, n)
	c, flat := sh.(solid)
	if flat {
		sh.shade(r.Min.X, r.Min.Y, n, buf) // the same in every row
		if mode == bdf.BlendSourceOver && alpha == 1 && c[3] == 1 && clip.a == nil && (cov == nil || cov.a == nil) {
			// an opaque colour over pixels it covers whole: it is what
			// they become, whatever they were
			s.setRows(r, buf)
			return
		}
	}
	if !flat {
		s.rowBufs(n)
	}
	s.bands(r, func(y0, y1, k int) {
		buf := buf
		if !flat && k > 0 {
			buf = s.rowBuf(k, n)
		}
		for y := y0; y < y1; y++ {
			if !flat {
				sh.shade(r.Min.X, y, n, buf)
			}
			krow := clip.row(y, r.Min.X, r.Max.X)
			var crow []float32
			if cov != nil {
				crow = cov.row(y, r.Min.X, r.Max.X)
			}
			d := s.pix[4*(y*s.w+r.Min.X) : 4*(y*s.w+r.Max.X)]
			if mode == bdf.BlendSourceOver && krow == nil {
				// the common case, without asking each pixel what the
				// clip and the composite operation are
				if crow == nil {
					fillRowOver(d, buf, alpha)
					continue
				}
				fillRowCovered(d, buf, crow, alpha)
				continue
			}
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
				src := [4]float32{float32(buf[4*i] * c), float32(buf[4*i+1] * c), float32(buf[4*i+2] * c), float32(buf[4*i+3] * c)}
				dst := [4]float32{d[4*i], d[4*i+1], d[4*i+2], d[4*i+3]}
				var out [4]float32
				if mode == bdf.BlendSourceOver {
					ia := 1 - src[3]
					out = [4]float32{src[0] + float32(dst[0]*ia), src[1] + float32(dst[1]*ia), src[2] + float32(dst[2]*ia), src[3] + float32(dst[3]*ia)}
				} else {
					out = composite(mode, src, dst)
				}
				if k == 1 {
					copy(d[4*i:4*i+4], out[:])
					continue
				}
				for j := 0; j < 4; j++ {
					d[4*i+j] = dst[j] + float32(k*(out[j]-dst[j]))
				}
			}
		}
	})
}

// fillRowOver composites the colours of a row, times alpha, over the
// pixels d, source-over.
func fillRowOver(d, buf []float32, alpha float32) {
	d, buf = d[:len(d)&^3], buf[:len(d)]
	for i := 0; i+3 < len(d); i += 4 {
		s0, s1, s2, s3 := float32(buf[i]*alpha), float32(buf[i+1]*alpha), float32(buf[i+2]*alpha), float32(buf[i+3]*alpha)
		ia := 1 - s3
		d[i], d[i+1], d[i+2], d[i+3] = s0+float32(d[i]*ia), s1+float32(d[i+1]*ia), s2+float32(d[i+2]*ia), s3+float32(d[i+3]*ia)
	}
}

// fillRowCovered composites the colours of a row, times the coverage of
// each pixel and alpha, over the pixels d, source-over.
func fillRowCovered(d, buf, crow []float32, alpha float32) {
	crow = crow[:len(d)/4]
	for x, c := range crow {
		c *= alpha
		if c == 0 {
			continue
		}
		i := 4 * x
		s0, s1, s2, s3 := float32(buf[i]*c), float32(buf[i+1]*c), float32(buf[i+2]*c), float32(buf[i+3]*c)
		ia := 1 - s3
		d[i], d[i+1], d[i+2], d[i+3] = s0+float32(d[i]*ia), s1+float32(d[i+1]*ia), s2+float32(d[i+2]*ia), s3+float32(d[i+3]*ia)
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
	row := s.rowBuf(0, r.Dx())
	if s.fresh {
		// over pixels of 0, the colour is what they become: a row of
		// it is made once and copied to the rest
		s.fresh = false
		repeat(row, c)
		s.setRows(r, row)
		return
	}
	ia := 1 - c[3]
	s.bands(r, func(y0, y1, _ int) {
		for y := y0; y < y1; y++ {
			d := s.pix[4*(y*s.w+r.Min.X) : 4*(y*s.w+r.Max.X)]
			for i := 0; i < len(d); i += 4 {
				p := d[i : i+4]
				p[0] = c[0] + float32(p[0]*ia)
				p[1] = c[1] + float32(p[1]*ia)
				p[2] = c[2] + float32(p[2]*ia)
				p[3] = c[3] + float32(p[3]*ia)
			}
		}
	})
}

// fillAll is fill for the pixels r of the clip, inside the coverage or not:
// the composite operations that change pixels the source does not cover.
func (s *surface) fillAll(cov *mask, sh shader, alpha float32, mode byte, clip *mask, r image.Rectangle) {
	s.fresh = false
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
			src := [4]float32{float32(buf[4*i] * c), float32(buf[4*i+1] * c), float32(buf[4*i+2] * c), float32(buf[4*i+3] * c)}
			dst := [4]float32{d[4*i], d[4*i+1], d[4*i+2], d[4*i+3]}
			var out [4]float32
			if mode == bdf.BlendSourceOver {
				ia := 1 - src[3]
				out = [4]float32{src[0] + float32(dst[0]*ia), src[1] + float32(dst[1]*ia), src[2] + float32(dst[2]*ia), src[3] + float32(dst[3]*ia)}
			} else {
				out = composite(mode, src, dst)
			}
			if k == 1 {
				copy(d[4*i:4*i+4], out[:])
				continue
			}
			for j := 0; j < 4; j++ {
				d[4*i+j] = dst[j] + float32(k*(out[j]-dst[j]))
			}
		}
	}
}

// composite applies a composite operation or blend mode to premultiplied
// source and destination colours (W3C Compositing and Blending).
func composite(mode byte, s, d [4]float32) [4]float32 {
	as, ab := s[3], d[3]
	pd := func(fs, fd float32) [4]float32 {
		return [4]float32{float32(s[0]*fs) + float32(d[0]*fd), float32(s[1]*fs) + float32(d[1]*fd), float32(s[2]*fs) + float32(d[2]*fd), float32(s[3]*fs) + float32(d[3]*fd)}
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
		out[i] = float32(s[i]*(1-ab)) + float32(d[i]*(1-as)) + float32(as*ab*b[i])
	}
	out[3] = as + float32(ab*(1-as))
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
	s.bands(s.bounds(), func(y0, y1, _ int) {
		pix := s.pix[4*y0*s.w : 4*y1*s.w]
		out := img.Pix[4*y0*s.w : 4*y1*s.w]
		out = out[:len(pix)]
		for i, v := range pix {
			out[i] = byteOf(v)
		}
	})
	return img
}

// byteOf is a colour value as a byte: rounded to the nearest, as a float32
// times 255 is exact as a float64, and so is the half added to what is
// not far below it.
func byteOf(v float32) uint8 {
	return uint8(float64(min(max(v, 0), 1))*255 + 0.5)
}

// grayOf is the luminance of a premultiplied pixel over white, as a
// byte: the weights are those image/color converts to gray with.
func grayOf(p []float32) uint8 {
	return byteOf(0.299*p[0] + 0.587*p[1] + 0.114*p[2] + (1 - p[3]))
}

// toGray converts the surface to an image of the luminance of its pixels
// over white (see Gray).
func (s *surface) toGray() *image.Gray {
	img := image.NewGray(image.Rect(0, 0, s.w, s.h))
	s.bands(s.bounds(), func(y0, y1, _ int) {
		pix := s.pix[4*y0*s.w : 4*y1*s.w]
		out := img.Pix[y0*s.w : y1*s.w]
		for i := range out {
			out[i] = grayOf(pix[4*i : 4*i+4])
		}
	})
	return img
}

// bilevelPalette is the colours of a Bilevel image: black, then white.
var bilevelPalette = color.Palette{color.Gray{0}, color.Gray{255}}

// toBilevel converts the surface to an image of black and white: black
// where the gray value is below 128 (see Bilevel).
func (s *surface) toBilevel() *image.Paletted {
	img := image.NewPaletted(image.Rect(0, 0, s.w, s.h), bilevelPalette)
	s.bands(s.bounds(), func(y0, y1, _ int) {
		pix := s.pix[4*y0*s.w : 4*y1*s.w]
		out := img.Pix[y0*s.w : y1*s.w]
		for i := range out {
			if grayOf(pix[4*i:4*i+4]) >= 128 {
				out[i] = 1
			}
		}
	})
	return img
}
