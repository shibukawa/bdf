package export

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/internal/affine"
)

// Num writes a number of a content stream, a dictionary or a procedure:
// without an exponent, with as many decimals as a place on a page or a
// color needs (four), within what single precision holds (the numbers of
// PDF and the reals of PostScript are single precision).
func Num(v float64) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "0"
	}
	v = math.Min(math.Max(v, -1e9), 1e9)
	r := math.Round(v*1e4) / 1e4
	if r == math.Trunc(r) {
		return strconv.FormatInt(int64(r), 10)
	}
	return strconv.FormatFloat(r, 'f', -1, 64)
}

// Nums writes numbers with spaces between them.
func Nums(vs ...float64) string {
	var b strings.Builder
	for i, v := range vs {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(Num(v))
	}
	return b.String()
}

// Fine writes a number that is multiplied by a size (a scale of a matrix,
// a width of a glyph, a component of a color), with more decimals than
// Num (six).
func Fine(v float64) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "0"
	}
	v = math.Min(math.Max(v, -1e9), 1e9)
	r := math.Round(v*1e6) / 1e6
	if r == math.Trunc(r) {
		return strconv.FormatInt(int64(r), 10)
	}
	return strconv.FormatFloat(r, 'f', -1, 64)
}

// Mat writes a matrix as the operands of a cm, a text, font, pattern or
// image matrix: its scales finely, its offsets as places.
func Mat(m affine.Matrix) string {
	return Fine(m[0]) + " " + Fine(m[1]) + " " + Fine(m[2]) + " " + Fine(m[3]) + " " + Num(m[4]) + " " + Num(m[5])
}
