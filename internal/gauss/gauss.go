// Package gauss blurs the planes of shadows: the renderers of BDF
// (raster/imagebdf, vector/pdfbdf) draw the same shadow with it.
package gauss

import (
	"math"
	"sync"
)

// Blur applies an approximate Gaussian blur (three box blurs) to a plane of
// w × h values, in place.
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
	tmp := make([]float32, len(a))
	for i := 0; i < 3; i++ {
		size := wl
		if i >= mm {
			size = wu
		}
		r := (size - 1) / 2
		if r <= 0 {
			continue
		}
		split(threads, h, w, func(y0, y1 int) { boxH(src(a, tmp, w, y0, y1), r) })
		split(threads, w, h, func(x0, x1 int) { boxV(tmp, a, w, h, x0, x1, r) })
	}
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

// rows is the rows [y0, y1) of a source and a destination plane.
type rows struct {
	src, dst []float32
	w, h     int
}

func src(a, b []float32, w, y0, y1 int) rows {
	return rows{src: a[y0*w : y1*w], dst: b[y0*w : y1*w], w: w, h: y1 - y0}
}

func boxH(p rows, r int) {
	src, dst, w, h := p.src, p.dst, p.w, p.h
	k := 1 / float32(2*r+1)
	for y := 0; y < h; y++ {
		row := src[y*w : (y+1)*w]
		out := dst[y*w : (y+1)*w]
		var sum float32
		for x := -r; x <= r; x++ {
			if x >= 0 && x < w {
				sum += row[x]
			}
		}
		for x := 0; x < w; x++ {
			out[x] = sum * k
			if x+r+1 < w {
				sum += row[x+r+1]
			}
			if x-r >= 0 {
				sum -= row[x-r]
			}
		}
	}
}

func boxV(src, dst []float32, w, h, x0, x1, r int) {
	k := 1 / float32(2*r+1)
	for x := x0; x < x1; x++ {
		var sum float32
		for y := -r; y <= r; y++ {
			if y >= 0 && y < h {
				sum += src[y*w+x]
			}
		}
		for y := 0; y < h; y++ {
			dst[y*w+x] = sum * k
			if y+r+1 < h {
				sum += src[(y+r+1)*w+x]
			}
			if y-r >= 0 {
				sum -= src[(y-r)*w+x]
			}
		}
	}
}
