package psd

import (
	"image"

	"github.com/shibukawa/bdf/imgconv"
)

// The layers are composited at the resolution the pages are stored at, not
// at the resolution of the document and reduced afterwards: a canvas that
// the resolution cap of image inputs (imgconv.Options.MaxDPI and MaxPixels)
// reduces takes the memory and the time of the reduced one. The buffers of
// the canvas, the groups and the fill layers have the reduced size. The
// pixels and the masks of a layer are decoded as they are stored and
// averaged into their place on the reduced canvas before they are blended.
//
// Rectangles (layers, masks, artboards, pages) are placed by their edges,
// each at the nearest pixel boundary of the reduced canvas: what shares an
// edge in the document shares it there, without a gap or an overlap. What
// is composited this way differs from a composite reduced afterwards where
// the pixels that one reduced pixel stands for differ among themselves in
// transparency, mask or backdrop: at the edges of layers, masks and
// artboards, which lie between reduced pixels instead of being averaged
// into them, and in blend modes other than normal, which blend the
// averages.

// fit is how a page is stored: its size in points and in pixels.
type fit struct {
	w, h   float32
	tw, th int
}

// fit returns how the page of a region is stored under the resolution cap.
func (rg region) fit(opts imgconv.Options, xres, yres float64) fit {
	w, h := float32(float64(rg.r.Dx())*72/xres), float32(float64(rg.r.Dy())*72/yres)
	tw, th := opts.FitSize(rg.r.Dx(), rg.r.Dy(), float64(w), float64(h))
	return fit{w, h, tw, th}
}

// compositeSize returns the size to composite the layers at for pages that
// are stored as fits says: the canvas reduced, along each axis, as the page
// that is reduced least, so that every page has the pixels it is stored
// with. A document within the resolution cap keeps the size of its canvas.
func compositeSize(canvas image.Point, pages []int, regions []region, fits []fit) image.Point {
	// the largest of the factors tw/w and th/h, as fractions
	x, y := [2]int64{0, 1}, [2]int64{0, 1}
	for i, n := range pages {
		r := regions[n-1].r
		if tw, w := int64(fits[i].tw), int64(r.Dx()); tw*x[1] > x[0]*w {
			x = [2]int64{tw, w}
		}
		if th, h := int64(fits[i].th), int64(r.Dy()); th*y[1] > y[0]*h {
			y = [2]int64{th, h}
		}
	}
	up := func(n int, f [2]int64) int { return max(1, min(n, int((int64(n)*f[0]+f[1]-1)/f[1]))) }
	return image.Pt(up(canvas.X, x), up(canvas.Y, y))
}

// scale places what the document holds on the canvas the layers are
// composited on: full is the size of the document's canvas, to the size it
// is composited at.
type scale struct {
	full, to image.Point
}

// reduced reports whether the layers are composited smaller than they are.
func (s scale) reduced() bool { return s.to != s.full }

// edge returns where a coordinate of the document is along an axis of full
// pixels composited as to pixels: at the nearest pixel boundary (halfway,
// the next).
func edge(v, to, full int) int {
	return int(floorDiv(2*int64(v)*int64(to)+int64(full), 2*int64(full)))
}

// floorDiv divides by a positive number and rounds towards minus infinity:
// layers start left of the canvas and above it too.
func floorDiv(n, d int64) int64 {
	q := n / d
	if n%d < 0 {
		q--
	}
	return q
}

// place returns the rectangle that a rectangle of the document takes.
func (s scale) place(r image.Rectangle) image.Rectangle {
	if !s.reduced() {
		return r
	}
	return image.Rect(edge(r.Min.X, s.to.X, s.full.X), edge(r.Min.Y, s.to.Y, s.full.Y),
		edge(r.Max.X, s.to.X, s.full.X), edge(r.Max.Y, s.to.Y, s.full.Y))
}

// weight is a pixel of a layer or mask and how much of a pixel of the
// reduced canvas it covers.
type weight struct{ i, w int }

// span is a side of a layer or mask along an axis of the reduced canvas:
// the pixel it starts at and, for each pixel it takes, the pixels averaged
// into it with their weights, and the weight of a pixel that is covered
// whole.
type span struct {
	start int
	px    [][]weight
	whole []int64
}

// span places the side from lo to hi of the document on an axis. A side at
// least as long as a pixel takes the pixels between its edges and fills
// them, each with the average of the part of the side that lies on it (as
// imgconv.Resize averages, in whole numbers): what the layer shows stays
// where it is. A shorter side takes the pixel its middle lies in, as far
// as it covers a pixel: a hairline becomes fainter, neither invisible nor
// as wide as a pixel.
func (s scale) span(lo, hi, to, full int) span {
	n, to64, full64 := hi-lo, int64(to), int64(full)
	if int64(n)*to64 < full64 {
		all := make([]weight, n)
		for i := range all {
			all[i] = weight{i, to}
		}
		mid := floorDiv((int64(lo)+int64(hi))*to64, 2*full64)
		return span{int(mid), [][]weight{all}, []int64{full64}}
	}
	// In units of 1/to pixel of the document, a pixel of the document is to
	// long and one of the reduced canvas full.
	a, b := edge(lo, to, full), edge(hi, to, full)
	sp := span{a, make([][]weight, b-a), make([]int64, b-a)}
	for t := range sp.px {
		from, until := max(int64(a+t)*full64, int64(lo)*to64), min(int64(a+t+1)*full64, int64(hi)*to64)
		for i := floorDiv(from, to64); i*to64 < until; i++ {
			if w := min(until, (i+1)*to64) - max(from, i*to64); w > 0 {
				sp.px[t] = append(sp.px[t], weight{int(i) - lo, int(w)})
				sp.whole[t] += w
			}
		}
	}
	return sp
}

// reduce returns the pixels of a layer, img over lr in the document, where
// they are on the reduced canvas: premultiplied averages. The sums are
// whole numbers, so that the result is the same on every platform.
func (cp *compositor) reduce(img *image.NRGBA, lr image.Rectangle) *rgba {
	if lr.Empty() {
		return nil
	}
	xs := cp.span(lr.Min.X, lr.Max.X, cp.to.X, cp.full.X)
	ys := cp.span(lr.Min.Y, lr.Max.Y, cp.to.Y, cp.full.Y)
	at := image.Rect(xs.start, ys.start, xs.start+len(xs.px), ys.start+len(ys.px))
	r := at.Intersect(cp.canvas)
	if r.Empty() {
		return nil
	}
	buf := newRGBA(r)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			var sum [4]uint64
			for _, wy := range ys.px[y-at.Min.Y] {
				row := img.Pix[wy.i*img.Stride:]
				for _, wx := range xs.px[x-at.Min.X] {
					p := row[4*wx.i : 4*wx.i+4]
					a := uint64(wy.w) * uint64(wx.w) * uint64(p[3])
					sum[0] += a * uint64(p[0])
					sum[1] += a * uint64(p[1])
					sum[2] += a * uint64(p[2])
					sum[3] += a
				}
			}
			whole := float64(xs.whole[x-at.Min.X]) * float64(ys.whole[y-at.Min.Y])
			i := buf.at(x, y)
			buf.pix[i] = float32(float64(sum[0]) / (whole * 255 * 255))
			buf.pix[i+1] = float32(float64(sum[1]) / (whole * 255 * 255))
			buf.pix[i+2] = float32(float64(sum[2]) / (whole * 255 * 255))
			buf.pix[i+3] = float32(float64(sum[3]) / (whole * 255))
		}
	}
	return buf
}

// reduceMask returns a mask, pix over r in the document, where it is on the
// reduced canvas: averages, with the mask's default value for what a mask
// thinner than a pixel leaves of it.
func (cp *compositor) reduceMask(pix []uint8, r image.Rectangle, deflt uint8) ([]uint8, image.Rectangle) {
	if r.Empty() {
		return nil, image.Rectangle{}
	}
	xs := cp.span(r.Min.X, r.Max.X, cp.to.X, cp.full.X)
	ys := cp.span(r.Min.Y, r.Max.Y, cp.to.Y, cp.full.Y)
	at := image.Rect(xs.start, ys.start, xs.start+len(xs.px), ys.start+len(ys.px))
	out := make([]uint8, at.Dx()*at.Dy())
	for y, wys := range ys.px {
		for x, wxs := range xs.px {
			var sum, covered uint64
			for _, wy := range wys {
				row := pix[wy.i*r.Dx():]
				for _, wx := range wxs {
					w := uint64(wy.w) * uint64(wx.w)
					sum += w * uint64(row[wx.i])
					covered += w
				}
			}
			whole := uint64(xs.whole[x]) * uint64(ys.whole[y])
			sum += (whole - covered) * uint64(deflt)
			out[y*at.Dx()+x] = uint8((sum + whole/2) / whole)
		}
	}
	return out, at
}
