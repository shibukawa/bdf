package imagebdf

import (
	"image"
	"image/color"
	"math"

	"github.com/shibukawa/bdf"
)

// surface is a canvas of premultiplied RGBA pixels, 0 to 1, at (ox, oy)
// of the device: the canvases of the groups and soft masks of a drawing
// lie on the same device as the canvas of the image, and a band of a page
// (see pageRows) is the rows of the page it draws.
//
// A canvas of a renderer is lazy: it begins as its background colour
// (that of the page, or none for a group) without the colour being
// written, and the pixels are written as the draws reach them, in
// blocks of blockPixels (see touch). What is not drawn on is never
// written, nor read: the conversions of the canvas to an image write the
// background for it. A page is mostly background, so a lazy canvas
// takes a fraction of the memory passes of one filled first.
type surface struct {
	w, h   int
	ox, oy int
	pix    []float32
	// page is the whole image a band is of (see pageRows), the bounds
	// when the canvas is the whole of it
	page image.Rectangle
	row  []float32 // the colours of a row being filled
	// rows holds the colours of the rows the other goroutines of a pass
	// fill (see splitRows)
	rows [][]float32
	// threads is how many goroutines the passes over the canvas are
	// split among (see splitRows)
	threads int
	// lazy says that the pixels of the blocks that are not initialized
	// (see init) are not in pix: they are bg
	lazy bool
	bg   solid
	// init has a bit a block of blockPixels pixels of a row, words a row
	// (see blockWords): set once the block's pixels are in pix
	init  []uint64
	words int
}

// blockPixels is the pixels of a block of a lazy canvas (see surface): a
// block is written whole with the background before a draw reads or
// writes part of it.
const blockPixels = 64

// newSurface returns a canvas of w × h pixels of 0, written: the canvases
// of a renderer are lazy instead (see Renderer.surfaceFor).
func newSurface(w, h int) *surface {
	return &surface{w: w, h: h, pix: make([]float32, 4*w*h), threads: 1}
}

// begin starts a lazy canvas anew with a background colour: no pixel is
// written.
func (s *surface) begin(bg solid) {
	s.lazy, s.bg = true, bg
	s.words = (blocks(s.w) + 63) / 64
	n := s.words * s.h
	if cap(s.init) < n {
		s.init = make([]uint64, n)
	} else {
		s.init = s.init[:n]
		clear(s.init)
	}
}

// blocks is the blocks of a row of w pixels.
func blocks(w int) int { return (w + blockPixels - 1) / blockPixels }

func (s *surface) bounds() image.Rectangle { return image.Rect(s.ox, s.oy, s.ox+s.w, s.oy+s.h) }

// at returns the index in pix of the pixel at (x, y) of the device, which
// must lie on the canvas.
func (s *surface) at(x, y int) int { return 4 * ((y-s.oy)*s.w + x - s.ox) }

// initRow returns the words of the blocks of row y of the device.
func (s *surface) initRow(y int) []uint64 {
	i := (y - s.oy) * s.words
	return s.init[i : i+s.words]
}

// initBlock writes the background over block b of row y, and notes it.
func (s *surface) initBlock(row []uint64, y, b int) {
	row[b>>6] |= 1 << (b & 63)
	x0 := b * blockPixels
	x1 := min(x0+blockPixels, s.w)
	p := s.pix[4*((y-s.oy)*s.w+x0) : 4*((y-s.oy)*s.w+x1)]
	repeat(p, s.bg)
}

// touch makes the pixels [x0, x1) of row y, in device coordinates, hold
// what they show, for a draw that reads or writes part of them: the
// blocks of a lazy canvas that hold them are written with the background
// if they are not initialized. With overwrite, the pixels are about to
// be written whole: the blocks within them are noted without being
// written.
func (s *surface) touch(y, x0, x1 int, overwrite bool) {
	if !s.lazy || x1 <= x0 {
		return
	}
	row := s.initRow(y)
	x0, x1 = x0-s.ox, x1-s.ox
	b0, b1 := x0/blockPixels, (x1-1)/blockPixels
	for b := b0; b <= b1; b++ {
		if row[b>>6]>>(b&63)&1 != 0 {
			continue
		}
		if overwrite && b*blockPixels >= x0 && min((b+1)*blockPixels, s.w) <= x1 {
			row[b>>6] |= 1 << (b & 63)
			continue
		}
		s.initBlock(row, y, b)
	}
}

// touchAll makes every pixel of the canvas hold what it shows.
func (s *surface) touchAll() {
	if !s.lazy {
		return
	}
	s.bands(s.bounds(), func(y0, y1, _ int) {
		for y := y0; y < y1; y++ {
			s.touch(y, s.ox, s.ox+s.w, false)
		}
	})
}

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
			s.touch(y, r.Min.X, r.Max.X, true)
			copy(s.pix[s.at(r.Min.X, y):s.at(r.Max.X, y)], row)
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
			d := s.pix[s.at(r.Min.X, y):s.at(r.Max.X, y)]
			// generic composites the pixels [i0, i1) of the row one by
			// one, asking each what the clip and the coverage are
			generic := func(i0, i1 int) {
				if i1 <= i0 {
					return
				}
				s.touch(y, r.Min.X+i0, r.Min.X+i1, false)
				for i := i0; i < i1; i++ {
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
			if mode != bdf.BlendSourceOver {
				generic(0, n)
				continue
			}
			// covered composites [i0, i1) through the coverage alone: the
			// clip covers them whole
			covered := func(i0, i1 int) {
				if i1 <= i0 {
					return
				}
				if s.lazy {
					s.fillRowCoveredLazy(y, r.Min.X+i0, d[4*i0:4*i1], buf[4*i0:4*i1], crow[i0:i1], alpha)
					return
				}
				fillRowCovered(d[4*i0:4*i1], buf[4*i0:4*i1], crow[i0:i1], alpha)
			}
			// fast composites [i0, i1), which the clip and the coverage
			// cover whole, without asking each pixel
			fast := func(i0, i1 int) {
				if i1 <= i0 {
					return
				}
				if flat && alpha == 1 && c[3] == 1 {
					// an opaque colour is what the pixels become
					s.touch(y, r.Min.X+i0, r.Min.X+i1, true)
					copy(d[4*i0:4*i1], buf[4*i0:4*i1])
					return
				}
				s.touch(y, r.Min.X+i0, r.Min.X+i1, false)
				fillRowOver(d[4*i0:4*i1], buf[4*i0:4*i1], alpha)
			}
			// the columns the clip covers whole, and within them those
			// the coverage does: the rest pixel by pixel
			k0, k1 := 0, n
			if krow != nil {
				x0, x1 := clip.within(y, r.Min.X, r.Max.X)
				k0, k1 = x0-r.Min.X, x1-r.Min.X
			}
			generic(0, k0)
			if crow == nil {
				fast(k0, k1)
			} else {
				c0, c1 := k0, k0
				if !cov.full.Empty() {
					x0, x1 := cov.within(y, r.Min.X+k0, r.Min.X+k1)
					c0, c1 = x0-r.Min.X, x1-r.Min.X
				}
				covered(k0, c0)
				fast(c0, c1)
				covered(c1, k1)
			}
			generic(k1, n)
		}
	})
}

// setBorders writes the pixels of a lazy canvas of the background that
// a rectangular clip with fractional edges (see rectMaskIn) covers in
// part: those of its rectangle outside the rectangle it covers whole, as
// the background times their coverage, over nothing (what fill writes
// over a canvas of nothing). The clip must lie on the canvas.
func (s *surface) setBorders(clip *mask, bg solid) {
	r, f := clip.r, clip.full
	s.bands(r, func(y0, y1, _ int) {
		for y := y0; y < y1; y++ {
			segs := [2][2]int{{r.Min.X, r.Max.X}, {r.Max.X, r.Max.X}}
			if y >= f.Min.Y && y < f.Max.Y {
				segs = [2][2]int{{r.Min.X, f.Min.X}, {f.Max.X, r.Max.X}}
			}
			for _, sg := range segs {
				x0, x1 := sg[0], sg[1]
				if x1 <= x0 {
					continue
				}
				s.touch(y, x0, x1, true)
				crow := clip.row(y, x0, x1)
				d := s.pix[s.at(x0, y):s.at(x1, y)]
				for i, c := range crow {
					d[4*i], d[4*i+1], d[4*i+2], d[4*i+3] = float32(bg[0]*c), float32(bg[1]*c), float32(bg[2]*c), float32(bg[3]*c)
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

// fillRowCoveredLazy is fillRowCovered on a lazy canvas: d is the pixels
// from x0 of row y, and the blocks of the pixels covered are initialized
// as they are reached, so that the ones between the glyphs of a line of
// text, say, are not.
func (s *surface) fillRowCoveredLazy(y, x0 int, d, buf, crow []float32, alpha float32) {
	row := s.initRow(y)
	rx := x0 - s.ox
	crow = crow[:len(d)/4]
	for x, c := range crow {
		c *= alpha
		if c == 0 {
			continue
		}
		if b := (rx + x) / blockPixels; row[b>>6]>>(b&63)&1 == 0 {
			s.initBlock(row, y, b)
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
	if c[3] == 1 {
		// an opaque colour is what the pixels become
		repeat(row, c)
		s.setRows(r, row)
		return
	}
	ia := 1 - c[3]
	s.bands(r, func(y0, y1, _ int) {
		for y := y0; y < y1; y++ {
			s.touch(y, r.Min.X, r.Max.X, false)
			d := s.pix[s.at(r.Min.X, y):s.at(r.Max.X, y)]
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
	n := r.Dx()
	buf := make([]float32, 4*n)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		s.touch(y, r.Min.X, r.Max.X, false)
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
		d := s.pix[s.at(r.Min.X, y):s.at(r.Max.X, y)]
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

// surfaceShader reads the pixels of another surface where it lies on the
// device: none outside it, and the background of the blocks of a lazy
// canvas that are not initialized.
type surfaceShader struct {
	s *surface
}

func (ss surfaceShader) shade(x, y, n int, dst []float32) {
	s := ss.s
	sy := y - s.oy
	if sx := x - s.ox; sx >= 0 && sy >= 0 && sx+n <= s.w && sy < s.h && !plain {
		// the whole run is of the surface
		if !s.lazy {
			copy(dst[:4*n], s.pix[4*(sy*s.w+sx):])
			return
		}
		row := s.initRow(y)
		for i := 0; i < n; {
			b := (sx + i) / blockPixels
			m := min(n, (b+1)*blockPixels-sx) // the pixels of the run in the block
			if row[b>>6]>>(b&63)&1 != 0 {
				copy(dst[4*i:4*m], s.pix[4*(sy*s.w+sx+i):])
			} else {
				repeat(dst[4*i:4*m], s.bg)
			}
			i = m
		}
		return
	}
	for i := 0; i < n; i++ {
		sx := x + i - s.ox
		if sx < 0 || sy < 0 || sx >= s.w || sy >= s.h {
			dst[4*i], dst[4*i+1], dst[4*i+2], dst[4*i+3] = 0, 0, 0, 0
			continue
		}
		if s.lazy {
			if row, b := s.initRow(y), sx/blockPixels; row[b>>6]>>(b&63)&1 == 0 {
				copy(dst[4*i:4*i+4], s.bg[:])
				continue
			}
		}
		copy(dst[4*i:4*i+4], s.pix[4*(sy*s.w+sx):])
	}
}

// toRGBA converts the surface to an image of premultiplied 8-bit pixels,
// whose bounds are those of the canvas.
func (s *surface) toRGBA() *image.RGBA {
	return s.toRGBAInto(nil, s.bounds())
}

// rgbaFor returns an image of the pixels of rows, which must be as wide
// as the canvas: img when it holds as many, else a new one.
func rgbaFor(img *image.RGBA, rows image.Rectangle) *image.RGBA {
	n := 4 * rows.Dx() * rows.Dy()
	if img == nil || cap(img.Pix) < n {
		return image.NewRGBA(rows)
	}
	return &image.RGBA{Pix: img.Pix[:n], Stride: 4 * rows.Dx(), Rect: rows}
}

// convert converts the rows of the canvas (its whole width) to the bytes
// of an image, bpp a pixel, with fn for the pixels of the blocks that are
// initialized, and bg (the bytes of a pixel of the background) for the
// rest of a lazy canvas.
func (s *surface) convert(rows image.Rectangle, out []byte, bpp int, bg []byte, fn func(pix []float32, out []byte)) {
	s.bands(rows, func(y0, y1, _ int) {
		var bgBlock []byte
		if s.lazy {
			bgBlock = make([]byte, bpp*blockPixels)
			copy(bgBlock, bg)
			for n := bpp; n < len(bgBlock); n *= 2 {
				copy(bgBlock[n:], bgBlock[:n])
			}
		}
		for y := y0; y < y1; y++ {
			pix := s.pix[s.at(s.ox, y):s.at(s.ox+s.w, y)]
			line := out[bpp*(y-rows.Min.Y)*s.w : bpp*(y-rows.Min.Y+1)*s.w]
			if !s.lazy {
				fn(pix, line)
				continue
			}
			row := s.initRow(y)
			for b, x0 := 0, 0; x0 < s.w; b, x0 = b+1, x0+blockPixels {
				x1 := min(x0+blockPixels, s.w)
				if row[b>>6]>>(b&63)&1 != 0 {
					fn(pix[4*x0:4*x1], line[bpp*x0:bpp*x1])
				} else {
					copy(line[bpp*x0:bpp*x1], bgBlock)
				}
			}
		}
	})
}

// toRGBAInto converts the rows of the canvas (its whole width) to an
// image of premultiplied 8-bit pixels with those bounds, into the memory
// of img when it is large enough (see rgbaFor).
func (s *surface) toRGBAInto(img *image.RGBA, rows image.Rectangle) *image.RGBA {
	img = rgbaFor(img, rows)
	bg := []byte{byteOf(s.bg[0]), byteOf(s.bg[1]), byteOf(s.bg[2]), byteOf(s.bg[3])}
	s.convert(rows, img.Pix, 4, bg, func(pix []float32, out []byte) {
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

// lumOf is the luminance of a premultiplied pixel over white, 0 to 1:
// the weights are those image/color converts to gray with. Each product
// is rounded, so that every conversion gets the same value.
func lumOf(p []float32) float32 {
	return float32(0.299*p[0]) + float32(0.587*p[1]) + float32(0.114*p[2]) + (1 - p[3])
}

// grayOf is the luminance of a premultiplied pixel over white, as a byte.
func grayOf(p []float32) uint8 { return byteOf(lumOf(p)) }

// toGray converts the surface to an image of the luminance of its pixels
// over white (see Gray).
func (s *surface) toGray() *image.Gray { return s.toGrayInto(nil, s.bounds()) }

// grayFor is rgbaFor for a gray image.
func grayFor(img *image.Gray, rows image.Rectangle) *image.Gray {
	n := rows.Dx() * rows.Dy()
	if img == nil || cap(img.Pix) < n {
		return image.NewGray(rows)
	}
	return &image.Gray{Pix: img.Pix[:n], Stride: rows.Dx(), Rect: rows}
}

// toGrayInto is toGray for the rows of the canvas, into the memory of img
// when it is large enough.
func (s *surface) toGrayInto(img *image.Gray, rows image.Rectangle) *image.Gray {
	img = grayFor(img, rows)
	s.convert(rows, img.Pix, 1, []byte{grayOf(s.bg[:])}, func(pix []float32, out []byte) {
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
func (s *surface) toBilevel() *image.Paletted { return s.toBilevelInto(nil, s.bounds()) }

// bilevelFor is rgbaFor for an image of black and white.
func bilevelFor(img *image.Paletted, rows image.Rectangle) *image.Paletted {
	n := rows.Dx() * rows.Dy()
	if img == nil || cap(img.Pix) < n {
		return image.NewPaletted(rows, bilevelPalette)
	}
	return &image.Paletted{Pix: img.Pix[:n], Stride: rows.Dx(), Rect: rows, Palette: bilevelPalette}
}

// bitOf is a pixel as black (0) or white (1): white from a gray value of
// 128 or more, a luminance of a half or more (127.5 / 255 is exactly 0.5).
func bitOf(p []float32) uint8 {
	if lumOf(p) >= 0.5 {
		return 1
	}
	return 0
}

// toBilevelInto is toBilevel for the rows of the canvas, into the memory
// of img when it is large enough.
func (s *surface) toBilevelInto(img *image.Paletted, rows image.Rectangle) *image.Paletted {
	img = bilevelFor(img, rows)
	s.convert(rows, img.Pix, 1, []byte{bitOf(s.bg[:])}, func(pix []float32, out []byte) {
		for i := range out {
			out[i] = bitOf(pix[4*i : 4*i+4])
		}
	})
	return img
}
