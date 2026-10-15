// Package gauss blurs the planes of shadows: the renderers of BDF
// (raster/imagebdf, vector/pdfbdf) draw the same shadow with it.
package gauss

import (
	"math"
	"sync"
)

// Blur applies an approximate Gaussian blur (three box blurs) to a plane of
// w × h values, 0 to 1, in place.
//
// The sums of the boxes are of 16-bit fixed-point values, exact in 32
// bits, so a blurred value depends on the values within its reach only:
// the blur of a part of a plane (the rows of a band of a page) is that of
// the whole wherever the part reaches as far as the blur does. A running
// sum of float32 would round by where the plane starts.
func Blur(a []float32, w, h int, sigma float64) { BlurN(a, w, h, sigma, 1) }

// BlurN is Blur with the rows and columns of a large plane divided among
// up to threads goroutines.
func BlurN(a []float32, w, h int, sigma float64, threads int) {
	// box sizes for three passes (Kovesi)
	wIdeal := math.Sqrt(12*sigma*sigma/3 + 1)
	wl := int(math.Floor(wIdeal))
	if wl%2 == 0 {
		wl--
	}
	wu := wl + 2
	mIdeal := (12*sigma*sigma - 3*float64(wl*wl) - 12*float64(wl) - 9) / (-4*float64(wl) - 4)
	mm := int(math.Round(mIdeal))
	n := w * h
	if n <= 0 || len(a) < n {
		return
	}
	a = a[:n]
	p, q := make([]uint16, n), make([]uint16, n)
	split(threads, h, w, func(y0, y1 int) {
		for i := y0 * w; i < y1*w; i++ {
			v := a[i]
			switch {
			case !(v > 0):
				p[i] = 0
			case v >= 1:
				p[i] = fixedOne
			default:
				p[i] = uint16(v*fixedOne + 0.5)
			}
		}
	})
	for i := 0; i < 3; i++ {
		size := wl
		if i >= mm {
			size = wu
		}
		r := (size - 1) / 2
		if r <= 0 {
			continue
		}
		split(threads, h, w, func(y0, y1 int) { boxH(p[y0*w:y1*w], q[y0*w:y1*w], w, y1-y0, r) })
		split(threads, w, h, func(x0, x1 int) { boxV(q, p, w, h, x0, x1, r) })
	}
	split(threads, h, w, func(y0, y1 int) {
		for i := y0 * w; i < y1*w; i++ {
			a[i] = float32(p[i]) / fixedOne
		}
	})
}

// fixedOne is 1 in the fixed-point values of the boxes: the sum of a box
// of up to 65536 values fits in 32 bits.
const fixedOne = 65535

// Reach is how far a blur of sigma reaches from a value, in values: the
// radius of its three boxes together. What is drawn past it from a value
// does not change the value.
func Reach(sigma float64) int {
	wIdeal := math.Sqrt(12*sigma*sigma/3 + 1)
	wl := int(math.Floor(wIdeal))
	if wl%2 == 0 {
		wl--
	}
	return 3 * ((wl + 2 - 1) / 2)
}

// minSplit is the least values of a plane whose blur is divided among
// goroutines.
const minSplit = 1 << 16

// split calls fn for bands of [0, n) in up to threads goroutines, and
// waits for them; each band is of the n of a plane of n × m values.
func split(threads, n, m int, fn func(lo, hi int)) {
	if threads <= 1 || n*m < minSplit || n < 2 {
		fn(0, n)
		return
	}
	threads = min(threads, n)
	var wg sync.WaitGroup
	wg.Add(threads)
	for k := 0; k < threads; k++ {
		lo, hi := n*k/threads, n*(k+1)/threads
		go func() {
			defer wg.Done()
			fn(lo, hi)
		}()
	}
	wg.Wait()
}

// boxMean is the mean of a box of 2r+1 values whose sum is sum, rounded
// to the nearest.
func boxMean(sum uint32, r int) uint16 {
	n := uint32(2*r + 1)
	return uint16((sum + n/2) / n)
}

// boxH averages each value of the rows of src with the r values on either
// side of it in its row, into dst.
func boxH(src, dst []uint16, w, h, r int) {
	for y := 0; y < h; y++ {
		row := src[y*w : (y+1)*w]
		out := dst[y*w : (y+1)*w]
		var sum uint32
		for x := 0; x <= r && x < w; x++ {
			sum += uint32(row[x])
		}
		for x := 0; x < w; x++ {
			out[x] = boxMean(sum, r)
			if x+r+1 < w {
				sum += uint32(row[x+r+1])
			}
			if x-r >= 0 {
				sum -= uint32(row[x-r])
			}
		}
	}
}

// boxV averages each value of the columns [x0, x1) of src with the r
// values above and below it, into dst.
func boxV(src, dst []uint16, w, h, x0, x1, r int) {
	for x := x0; x < x1; x++ {
		var sum uint32
		for y := 0; y <= r && y < h; y++ {
			sum += uint32(src[y*w+x])
		}
		for y := 0; y < h; y++ {
			dst[y*w+x] = boxMean(sum, r)
			if y+r+1 < h {
				sum += uint32(src[(y+r+1)*w+x])
			}
			if y-r >= 0 {
				sum -= uint32(src[(y-r)*w+x])
			}
		}
	}
}
