package gauss

import (
	"math"
	"math/rand"
	"testing"
)

// TestBlurOfAPart checks that the blur of a part of a plane is that of
// the whole wherever the part reaches as far as the blur does: the sums
// of the boxes are exact, so a value does not depend on where the plane
// starts.
func TestBlurOfAPart(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for _, sigma := range []float64{0.4, 1, 2.5, 7, 20} {
		w, h := 37, 200
		whole := make([]float32, w*h)
		for i := range whole {
			if rng.Intn(3) > 0 {
				whole[i] = rng.Float32()
			}
		}
		part := make([]float32, w*h)
		copy(part, whole)
		y0, y1 := 60, 150
		part = part[y0*w : y1*w]
		BlurN(whole, w, h, sigma, 4)
		Blur(part, w, y1-y0, sigma)
		reach := Reach(sigma)
		for y := y0 + reach; y < y1-reach; y++ {
			for x := 0; x < w; x++ {
				if a, b := whole[y*w+x], part[(y-y0)*w+x]; a != b {
					t.Fatalf("sigma %g: (%d, %d) is %v in the whole, %v in the part", sigma, x, y, a, b)
				}
			}
		}
		// a value is a mean of what is around it: within 0 and 1
		var sum float64
		for _, v := range whole {
			sum += float64(v)
			if v < 0 || v > 1 {
				t.Fatalf("sigma %g: a value of %v", sigma, v)
			}
		}
		if sum == 0 {
			t.Fatalf("sigma %g: blurred to nothing", sigma)
		}
	}
}

// TestBlurOfAPoint checks that a point blurs to a bump that reaches no
// further than Reach says.
func TestBlurOfAPoint(t *testing.T) {
	for _, sigma := range []float64{1, 3, 10} {
		w := 2*Reach(sigma) + 11
		a := make([]float32, w*w)
		a[(w/2)*w+w/2] = 1
		Blur(a, w, w, sigma)
		reach := Reach(sigma)
		for y := 0; y < w; y++ {
			for x := 0; x < w; x++ {
				v := a[y*w+x]
				far := math.Abs(float64(x-w/2)) > float64(reach) || math.Abs(float64(y-w/2)) > float64(reach)
				if far && v != 0 {
					t.Fatalf("sigma %g: (%d, %d) is %v beyond the reach %d", sigma, x, y, v, reach)
				}
			}
		}
		if a[(w/2)*w+w/2] == 0 || a[(w/2)*w+w/2] >= 1 {
			t.Fatalf("sigma %g: the middle is %v", sigma, a[(w/2)*w+w/2])
		}
	}
}
