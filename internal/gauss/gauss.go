// Package gauss blurs the planes of shadows: the renderers of BDF
// (raster/imagebdf, vector/pdfbdf) draw the same shadow with it.
package gauss

import "math"

// Blur applies an approximate Gaussian blur (three box blurs) to a plane of
// w × h values, in place.
func Blur(a []float32, w, h int, sigma float64) {
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
		boxH(a, tmp, w, h, r)
		boxV(tmp, a, w, h, r)
	}
}

func boxH(src, dst []float32, w, h, r int) {
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

func boxV(src, dst []float32, w, h, r int) {
	k := 1 / float32(2*r+1)
	for x := 0; x < w; x++ {
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
