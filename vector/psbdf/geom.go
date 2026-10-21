package psbdf

import (
	"math"

	"github.com/shibukawa/bdf/internal/affine"
)

// matrix is an affine transform: x' = a·x + c·y + e, y' = b·x + d·y + f,
// as [a b c d e f] (see affine.Matrix). The pen that writes paths and
// the outlines of glyphs as the procedures of the prolog (m, l, c, h) is
// export.Pen.
type matrix = affine.Matrix

var identity = affine.Identity

// union returns the box around two boxes; empty boxes (x1 < x0) are none.
func union(a, b [4]float64) [4]float64 {
	if !(a[2] >= a[0]) {
		return b
	}
	if !(b[2] >= b[0]) {
		return a
	}
	return [4]float64{math.Min(a[0], b[0]), math.Min(a[1], b[1]), math.Max(a[2], b[2]), math.Max(a[3], b[3])}
}

// intersect returns where two boxes overlap (empty when they do not).
func intersect(a, b [4]float64) [4]float64 {
	return [4]float64{math.Max(a[0], b[0]), math.Max(a[1], b[1]), math.Min(a[2], b[2]), math.Min(a[3], b[3])}
}

// grow returns a box grown by d on each side.
func grow(b [4]float64, d float64) [4]float64 {
	return [4]float64{b[0] - d, b[1] - d, b[2] + d, b[3] + d}
}

// none is the empty box.
var none = [4]float64{math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)}

// everywhere is the box of what is not cut.
var everywhere = [4]float64{math.Inf(-1), math.Inf(-1), math.Inf(1), math.Inf(1)}
