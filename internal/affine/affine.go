// Package affine has the affine transform of Canvas 2D, as the rasterizer
// and the vector exporters keep it.
package affine

import (
	"math"
	"strconv"
	"strings"
)

// Matrix is an affine transform as Canvas 2D keeps it: x' = a·x + c·y + e,
// y' = b·x + d·y + f, as [a b c d e f].
type Matrix [6]float64

// Identity is the transform that moves nothing.
var Identity = Matrix{1, 0, 0, 1, 0, 0}

// Mul returns m × n: the transform that applies n, then m
// (ctx.transform(n) on a context at m).
func (m Matrix) Mul(n Matrix) Matrix {
	return Matrix{
		m[0]*n[0] + m[2]*n[1],
		m[1]*n[0] + m[3]*n[1],
		m[0]*n[2] + m[2]*n[3],
		m[1]*n[2] + m[3]*n[3],
		m[0]*n[4] + m[2]*n[5] + m[4],
		m[1]*n[4] + m[3]*n[5] + m[5],
	}
}

// Apply returns a point under the transform.
func (m Matrix) Apply(x, y float64) (float64, float64) {
	return m[0]*x + m[2]*y + m[4], m[1]*x + m[3]*y + m[5]
}

// Translate returns the transform moved by (x, y) in its own space.
func (m Matrix) Translate(x, y float64) Matrix { return m.Mul(Matrix{1, 0, 0, 1, x, y}) }

// Scale returns the transform scaled by (x, y) in its own space.
func (m Matrix) Scale(x, y float64) Matrix { return m.Mul(Matrix{x, 0, 0, y, 0, 0}) }

// Det is the determinant of the transform: how much it scales areas.
func (m Matrix) Det() float64 { return m[0]*m[3] - m[1]*m[2] }

// Invert returns the inverse of the transform, when it has one that is
// finite.
func (m Matrix) Invert() (Matrix, bool) {
	det := m.Det()
	if det == 0 || math.IsNaN(det) || math.IsInf(det, 0) || math.IsInf(1/det, 0) {
		return Matrix{}, false
	}
	return Matrix{
		m[3] / det, -m[1] / det, -m[2] / det, m[0] / det,
		(m[2]*m[5] - m[3]*m[4]) / det, (m[1]*m[4] - m[0]*m[5]) / det,
	}, true
}

// MaxScale is the largest factor the transform stretches a length by.
func (m Matrix) MaxScale() float64 {
	return math.Max(math.Hypot(m[0], m[1]), math.Hypot(m[2], m[3]))
}

// AreaScale is how much the transform scales lengths, as the root of how
// much it scales areas (the zoom of a uniform transform).
func (m Matrix) AreaScale() float64 { return math.Sqrt(math.Abs(m.Det())) }

// Finite reports whether every number of the transform is one.
func (m Matrix) Finite() bool {
	for _, v := range m {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return false
		}
	}
	return true
}

// Bounds returns the box around a rectangle under the transform.
func (m Matrix) Bounds(x0, y0, x1, y1 float64) [4]float64 {
	out := [4]float64{math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)}
	for _, c := range [][2]float64{{x0, y0}, {x1, y0}, {x0, y1}, {x1, y1}} {
		x, y := m.Apply(c[0], c[1])
		out[0], out[1], out[2], out[3] = math.Min(out[0], x), math.Min(out[1], y), math.Max(out[2], x), math.Max(out[3], y)
	}
	return out
}

// Exact writes the matrix with the digits that a transform and its
// inverse, one after the other, need to be no transform.
func (m Matrix) Exact() string {
	var b strings.Builder
	for i, v := range m {
		if i > 0 {
			b.WriteByte(' ')
		}
		if math.IsNaN(v) || math.IsInf(v, 0) || v == 0 {
			v = 0 // (and no -0)
		}
		b.WriteString(strconv.FormatFloat(v, 'f', -1, 32))
	}
	return b.String()
}
