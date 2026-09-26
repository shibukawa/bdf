package jpx

// Inverse discrete wavelet transform (Annex F).
//
// The float32 conversions around products keep the compiler from fusing
// them with the following addition, so that results do not depend on the
// platform.
//
// A tile-component buffer holds the sub-bands of each resolution level in
// quadrants: the lower resolution (LL) at the top left, HL to its right, LH
// below and HH diagonally. Every level is reconstructed in place by
// filtering all rows (HOR_SR) and then all columns (VER_SR). The low-pass
// samples of a line sit at even coordinates of the reference grid, so the
// parity of the first coordinate decides how the two halves interleave.

// Lifting constants of the 9/7 filter (Table F.4).
const (
	alpha97 = -1.586134342059924
	beta97  = -0.052980118572961
	gamma97 = 0.882911075530934
	delta97 = 0.443506852043971
	k97     = 1.230174104914001
)

// invK97 scales the high-pass samples of the 9/7 filter. OpenJPEG uses
// 1.625732422/2, a fixed-point approximation, and the tests switch to it
// to compare bit for bit.
var invK97 float32 = 1 / k97

// level is the geometry of one resolution level: its size, the size of the
// next lower level (the low-pass halves) and the parity of its origin.
type level struct {
	w, h   int
	lw, lh int
	px, py int
}

// stripWidth is the number of columns that the vertical pass filters at a
// time; the scratch buffer needs max(width, stripWidth×height) samples.
const stripWidth = 64

// scratchSize returns the size of the scratch buffer for a w×h
// tile-component.
func scratchSize(w, h int) int { return max(w, min(w, stripWidth)*h) }

// idwt53 reconstructs a tile-component with the reversible 5/3 filter.
func idwt53(buf []int32, stride int, levels []level, tmp []int32) {
	for _, l := range levels {
		if l.w == 0 || l.h == 0 {
			continue
		}
		row := tmp[:l.w]
		for y := range l.h {
			line := buf[y*stride : y*stride+l.w]
			interleave(row, line[:l.lw], line[l.lw:], l.px)
			lift53(row, l.px)
			copy(line, row)
		}
		if l.h == 1 {
			if l.py == 1 {
				for x := range l.w {
					buf[x] /= 2
				}
			}
			continue
		}
		// Interleave the rows of a strip into tmp, lift whole rows, copy
		// back.
		n := l.h
		for x0 := 0; x0 < l.w; x0 += stripWidth {
			w := min(stripWidth, l.w-x0)
			t := tmp[:w*n]
			for y := range n {
				k := interleavedRow(y, l.lh, l.py)
				copy(t[k*w:k*w+w], buf[y*stride+x0:])
			}
			for k := l.py; k < n; k += 2 {
				a, b := mirror(k-1, n), mirror(k+1, n)
				dst, ra, rb := t[k*w:k*w+w], t[a*w:a*w+w], t[b*w:b*w+w]
				for i := range dst {
					dst[i] -= (ra[i] + rb[i] + 2) >> 2
				}
			}
			for k := 1 - l.py; k < n; k += 2 {
				a, b := mirror(k-1, n), mirror(k+1, n)
				dst, ra, rb := t[k*w:k*w+w], t[a*w:a*w+w], t[b*w:b*w+w]
				for i := range dst {
					dst[i] += (ra[i] + rb[i]) >> 1
				}
			}
			for y := range n {
				copy(buf[y*stride+x0:y*stride+x0+w], t[y*w:y*w+w])
			}
		}
	}
}

// interleavedRow returns where row y of a level, whose first lh rows are
// low-pass, goes when the rows are interleaved; p is the parity of the
// first row.
func interleavedRow(y, lh, p int) int {
	if y < lh {
		return 2*y + p
	}
	return 2*(y-lh) + 1 - p
}

// lift53 is the 1-D 5/3 synthesis of an interleaved line whose first
// sample has parity p. The symmetric extension mirrors x[1] to x[-1] and
// x[n-2] to x[n].
func lift53(x []int32, p int) {
	n := len(x)
	if n == 1 {
		if p == 1 {
			x[0] /= 2
		}
		return
	}
	k := p
	if k == 0 {
		x[0] -= (2*x[1] + 2) >> 2
		k = 2
	}
	for ; k < n-1; k += 2 {
		x[k] -= (x[k-1] + x[k+1] + 2) >> 2
	}
	if k == n-1 {
		x[k] -= (2*x[k-1] + 2) >> 2
	}
	k = 1 - p
	if k == 0 {
		x[0] += x[1]
		k = 2
	}
	for ; k < n-1; k += 2 {
		x[k] += (x[k-1] + x[k+1]) >> 1
	}
	if k == n-1 {
		x[k] += x[k-1]
	}
}

// idwt97 reconstructs a tile-component with the irreversible 9/7 filter.
func idwt97(buf []float32, stride int, levels []level, tmp []float32) {
	for _, l := range levels {
		if l.w == 0 || l.h == 0 {
			continue
		}
		row := tmp[:l.w]
		for y := range l.h {
			line := buf[y*stride : y*stride+l.w]
			interleave(row, line[:l.lw], line[l.lw:], l.px)
			lift97(row, l.px)
			copy(line, row)
		}
		if l.h == 1 {
			if l.py == 1 {
				for x := range l.w {
					buf[x] *= 0.5
				}
			}
			continue
		}
		n := l.h
		for x0 := 0; x0 < l.w; x0 += stripWidth {
			w := min(stripWidth, l.w-x0)
			t := tmp[:w*n]
			for y := range n {
				k := interleavedRow(y, l.lh, l.py)
				var s float32 = k97
				if y >= l.lh {
					s = invK97
				}
				d, r := t[k*w:k*w+w], buf[y*stride+x0:y*stride+x0+w]
				for i := range d {
					d[i] = r[i] * s
				}
			}
			liftRows97(t, w, n, l.py, -delta97)
			liftRows97(t, w, n, 1-l.py, -gamma97)
			liftRows97(t, w, n, l.py, -beta97)
			liftRows97(t, w, n, 1-l.py, -alpha97)
			for y := range n {
				copy(buf[y*stride+x0:y*stride+x0+w], t[y*w:y*w+w])
			}
		}
	}
}

func liftRows97(t []float32, w, n, start int, c float32) {
	for k := start; k < n; k += 2 {
		a, b := mirror(k-1, n), mirror(k+1, n)
		dst, ra, rb := t[k*w:k*w+w], t[a*w:a*w+w], t[b*w:b*w+w]
		for i := range dst {
			dst[i] += float32(c * (ra[i] + rb[i]))
		}
	}
}

// lift97 is the 1-D 9/7 synthesis of an interleaved line whose first
// sample has parity p.
func lift97(x []float32, p int) {
	n := len(x)
	if n == 1 {
		if p == 1 {
			x[0] *= 0.5
		}
		return
	}
	for k := p; k < n; k += 2 {
		x[k] *= k97
	}
	for k := 1 - p; k < n; k += 2 {
		x[k] *= invK97
	}
	lift97Step(x, p, -delta97)
	lift97Step(x, 1-p, -gamma97)
	lift97Step(x, p, -beta97)
	lift97Step(x, 1-p, -alpha97)
}

func lift97Step(x []float32, k int, c float32) {
	n := len(x)
	if k == 0 {
		x[0] += float32(c * (x[1] + x[1]))
		k = 2
	}
	for ; k < n-1; k += 2 {
		x[k] += float32(c * (x[k-1] + x[k+1]))
	}
	if k == n-1 {
		x[k] += float32(c * (x[k-1] + x[k-1]))
	}
}

// mirror applies the symmetric extension of a line of n > 1 samples.
func mirror(k, n int) int {
	if k < 0 {
		return -k
	}
	if k >= n {
		return 2*(n-1) - k
	}
	return k
}

// interleave merges the low-pass and high-pass halves of a line; p is the
// parity of its first sample.
func interleave[T int32 | float32](dst, lo, hi []T, p int) {
	for i, v := range lo {
		dst[2*i+p] = v
	}
	for i, v := range hi {
		dst[2*i+1-p] = v
	}
}
