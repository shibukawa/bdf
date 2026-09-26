package imgconv

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// The resolution cap of raster inputs (Options.MaxDPI, Options.MaxPixels).
const (
	// DefaultMaxDPI is twice the CSS reference pixel density of 96 per inch:
	// a page at its real size on a 2× (Retina) display shows one image pixel
	// per device pixel.
	DefaultMaxDPI = 192
	// DefaultMaxPixels is 3840 × 3840 pixels: the width of a 4K display in
	// both directions (56 MiB decoded as RGBA).
	DefaultMaxPixels = 3840 * 3840
)

func (o Options) maxDPI() float64 {
	switch {
	case o.MaxDPI == 0:
		return DefaultMaxDPI
	case o.MaxDPI < 0:
		return 0
	}
	return o.MaxDPI
}

func (o Options) maxPixels() int {
	switch {
	case o.MaxPixels == 0:
		return DefaultMaxPixels
	case o.MaxPixels < 0:
		return 0
	}
	return o.MaxPixels
}

// FitSize returns the size in pixels to store a raster input of w×h pixels
// that is shown wPt×hPt points large (1/72 inch): each axis is scaled down
// to at most MaxDPI pixels per inch, then both by the same factor to at most
// MaxPixels pixels. It never scales up, and an axis whose size on the page
// is unknown (0) is capped by MaxPixels only.
func (o Options) FitSize(w, h int, wPt, hPt float64) (int, int) {
	if w <= 0 || h <= 0 {
		return w, h
	}
	fit := func(n int, pt float64) int {
		dpi := o.maxDPI()
		if dpi <= 0 || pt <= 0 {
			return n
		}
		// The explicit conversions keep a fused multiply-add from changing
		// the result between architectures.
		s := float64(dpi * pt / 72 / float64(n))
		if s >= 1-1e-9 {
			return n
		}
		return max(1, int(float64(float64(n)*s)))
	}
	nw, nh := fit(w, wPt), fit(h, hPt)
	if mp := int64(o.maxPixels()); mp > 0 && int64(nw)*int64(nh) > mp {
		// math.Sqrt is correctly rounded, so this is the same everywhere.
		s := math.Sqrt(float64(mp) / float64(int64(nw)*int64(nh)))
		nw, nh = max(1, int(float64(float64(nw)*s))), max(1, int(float64(float64(nh)*s)))
		for int64(nw)*int64(nh) > mp { // rounding
			if nw >= nh {
				nw--
			} else {
				nh--
			}
		}
	}
	return nw, nh
}

// Resize scales img down to w×h pixels by averaging the source pixels each
// target pixel covers (a box filter, the usual way to downsample scans).
// The arithmetic is exact integer arithmetic, so the result is the same on
// every platform.
//
// A two-colour paletted image (a bilevel scan) stays two-colour, which
// keeps it small: a target pixel takes the less frequent colour (the ink)
// when that covers at least two fifths of it, a little less than half, so
// that thin strokes survive without text turning bold.
// Grey images stay grey; everything else becomes RGBA (premultiplied, so
// that transparent pixels do not bleed into their neighbours). img is
// returned as it is when it already has that size or when w or h is not
// smaller than the source.
func Resize(img image.Image, w, h int) image.Image {
	b := img.Bounds()
	sw, sh := b.Dx(), b.Dy()
	if w <= 0 || h <= 0 || (w >= sw && h >= sh) {
		return img
	}
	w, h = min(w, sw), min(h, sh)
	switch src := img.(type) {
	case *image.Paletted:
		if len(src.Palette) == 2 {
			return resizeBilevel(src, w, h)
		}
	case *image.Gray:
		dst := image.NewGray(image.Rect(0, 0, w, h))
		boxResample(src.Pix[src.PixOffset(b.Min.X, b.Min.Y):], src.Stride, sw, sh, 1, dst.Pix, dst.Stride, w, h)
		return dst
	}
	src, ok := img.(*image.RGBA)
	if !ok {
		src = image.NewRGBA(image.Rect(0, 0, sw, sh))
		draw.Draw(src, src.Bounds(), img, b.Min, draw.Src)
		b = src.Bounds()
	}
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	boxResample(src.Pix[src.PixOffset(b.Min.X, b.Min.Y):], src.Stride, sw, sh, 4, dst.Pix, dst.Stride, w, h)
	return dst
}

// resizeBilevel scales a two-colour image, keeping the ink (the less
// frequent colour) where it covers at least two fifths of a target pixel.
func resizeBilevel(src *image.Paletted, w, h int) *image.Paletted {
	b := src.Bounds()
	sw, sh := b.Dx(), b.Dy()
	cov := make([]uint8, sw*sh) // 255 where the source pixel is ink
	var n1 int
	for y := range sh {
		for _, v := range src.Pix[src.PixOffset(b.Min.X, b.Min.Y+y):][:sw] {
			if v == 1 {
				n1++
			}
		}
	}
	ink := uint8(1)
	if 2*n1 > sw*sh {
		ink = 0
	}
	for y := range sh {
		row := src.Pix[src.PixOffset(b.Min.X, b.Min.Y+y):][:sw]
		for x, v := range row {
			if v == ink {
				cov[y*sw+x] = 255
			}
		}
	}
	avg := make([]uint8, w*h)
	boxResample(cov, sw, sw, sh, 1, avg, w, w, h)
	dst := image.NewPaletted(image.Rect(0, 0, w, h), color.Palette{src.Palette[0], src.Palette[1]})
	for i, v := range avg {
		if v >= 102 { // 2/5 of 255
			dst.Pix[i] = ink
		} else {
			dst.Pix[i] = 1 - ink
		}
	}
	return dst
}

// boxResample averages the sw×sh pixels of src (ch interleaved 8-bit
// channels per pixel, rows stride bytes apart) down to the dw×dh pixels of
// dst. Target pixel x covers source columns [x·sw/dw, (x+1)·sw/dw); in
// units of 1/dw of a source pixel the overlaps are integers that add up to
// sw, so the weighted sums are exact and are rounded once at the end.
func boxResample(src []byte, stride, sw, sh, ch int, dst []byte, dstride, dw, dh int) {
	xw := boxWeights(sw, dw)
	yw := boxWeights(sh, dh)
	den := uint64(sw) * uint64(sh)
	row := make([]uint64, dw*ch) // one source row, summed horizontally
	acc := make([]uint64, dw*ch)
	for y := range dh {
		clear(acc)
		for _, ry := range yw[y] {
			s := src[ry.i*stride:]
			clear(row)
			for x := range dw {
				o := row[x*ch : x*ch+ch]
				for _, rx := range xw[x] {
					p := s[rx.i*ch : rx.i*ch+ch]
					for c := range ch {
						o[c] += uint64(rx.w) * uint64(p[c])
					}
				}
			}
			for k, v := range row {
				acc[k] += uint64(ry.w) * v
			}
		}
		d := dst[y*dstride : y*dstride+dw*ch]
		for k, v := range acc {
			d[k] = uint8((v + den/2) / den)
		}
	}
}

// boxWeight is a source pixel and how much of a target pixel it covers.
type boxWeight struct{ i, w int }

// boxWeights returns, for each of the m target pixels along an axis of n
// source pixels, the source pixels it covers and their overlaps in units of
// 1/m source pixel (they add up to n).
func boxWeights(n, m int) [][]boxWeight {
	out := make([][]boxWeight, m)
	n64, m64 := int64(n), int64(m) // products overflow 32-bit ints (TinyGo) on big images
	for t := range m64 {
		lo, hi := t*n64, (t+1)*n64
		for i := lo / m64; i*m64 < hi; i++ {
			w := min(hi, (i+1)*m64) - max(lo, i*m64)
			if w > 0 {
				out[t] = append(out[t], boxWeight{int(i), int(w)})
			}
		}
	}
	return out
}
