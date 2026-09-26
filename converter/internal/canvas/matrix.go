package canvas

import "math"

// Matrix is an affine transform [a b c d e f] (Canvas order).
type Matrix [6]float64

// Identity is the transform that changes nothing.
var Identity = Matrix{1, 0, 0, 1, 0, 0}

// Mul returns m·n: n applied first, then m.
func (m Matrix) Mul(n Matrix) Matrix {
	return Matrix{
		m[0]*n[0] + m[2]*n[1], m[1]*n[0] + m[3]*n[1],
		m[0]*n[2] + m[2]*n[3], m[1]*n[2] + m[3]*n[3],
		m[0]*n[4] + m[2]*n[5] + m[4], m[1]*n[4] + m[3]*n[5] + m[5],
	}
}

// Apply transforms a point.
func (m Matrix) Apply(x, y float64) (float64, float64) {
	return m[0]*x + m[2]*y + m[4], m[1]*x + m[3]*y + m[5]
}

// Translate returns a translation.
func Translate(x, y float64) Matrix { return Matrix{1, 0, 0, 1, x, y} }

// Scale returns a scaling.
func Scale(x, y float64) Matrix { return Matrix{x, 0, 0, y, 0, 0} }

// Rotate returns a rotation by deg degrees (clockwise in y-down space).
func Rotate(deg float64) Matrix {
	s, c := math.Sincos(deg * math.Pi / 180)
	return Matrix{c, s, -s, c, 0, 0}
}
