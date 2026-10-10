package svg

import (
	"math"
	"strconv"
	"strings"
)

// ParseNumber reads a number at the start of s, after white space and
// commas: its value and the rest of s. It reports false when s does not
// start with a number.
func ParseNumber(s string) (float64, string, bool) {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == ',' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
		i++
	}
	s = s[i:]
	j := 0
	if j < len(s) && (s[j] == '+' || s[j] == '-') {
		j++
	}
	digits, dot := false, false
	for j < len(s) {
		c := s[j]
		if c >= '0' && c <= '9' {
			digits = true
		} else if c == '.' && !dot {
			dot = true
		} else {
			break
		}
		j++
	}
	if !digits {
		return 0, s, false
	}
	if j < len(s) && (s[j] == 'e' || s[j] == 'E') {
		k := j + 1
		if k < len(s) && (s[k] == '+' || s[k] == '-') {
			k++
		}
		if k < len(s) && s[k] >= '0' && s[k] <= '9' {
			for k < len(s) && s[k] >= '0' && s[k] <= '9' {
				k++
			}
			j = k
		}
	}
	v, err := strconv.ParseFloat(s[:j], 64)
	return v, s[j:], err == nil
}

// ParseNumbers reads a list of numbers separated by white space or commas
// (a viewBox, the points of a polygon), up to what is not a number.
func ParseNumbers(s string) []float64 {
	var out []float64
	for {
		v, rest, ok := ParseNumber(s)
		if !ok {
			return out
		}
		out = append(out, v)
		s = rest
	}
}

// unit sizes in px
var units = map[string]float64{"": 1, "px": 1, "pt": 96.0 / 72, "pc": 16, "in": 96, "cm": 96 / 2.54, "mm": 96 / 25.4, "q": 96 / 101.6}

// ParseLength reads a length in CSS px: a number alone or with a unit (px,
// pt, pc, in, cm, mm, q), a percentage of ref, or em and ex of the font
// size em. It reports false for what is no length.
func ParseLength(s string, ref, em float64) (float64, bool) {
	v, rest, ok := ParseNumber(s)
	if !ok {
		return 0, false
	}
	u := strings.ToLower(strings.TrimSpace(rest))
	switch u {
	case "%":
		return v * ref / 100, true
	case "em":
		return v * em, true
	case "ex":
		return v * em / 2, true
	}
	if k, ok := units[u]; ok {
		return v * k, true
	}
	return 0, false
}

// Matrix is an affine transform as Canvas 2D keeps it: x' = a x + c y + e,
// y' = b x + d y + f, for the matrix {a, b, c, d, e, f}.
type Matrix [6]float64

// Identity is the transform that changes nothing.
var Identity = Matrix{1, 0, 0, 1, 0, 0}

// Mul returns m × n: n applies first (the transform of an element inside
// one whose transform is m).
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

// ParseTransform reads a transform list (matrix, translate, scale, rotate,
// skewX, skewY) into one matrix, up to what is no transform.
func ParseTransform(s string) Matrix {
	m := Identity
	for {
		s = strings.TrimLeft(s, " ,\t\n\r")
		open := strings.IndexByte(s, '(')
		if open < 0 {
			return m
		}
		name := strings.TrimSpace(s[:open])
		close := strings.IndexByte(s[open:], ')')
		if close < 0 {
			return m
		}
		a := ParseNumbers(s[open+1 : open+close])
		s = s[open+close+1:]
		arg := func(i int, d float64) float64 {
			if i < len(a) {
				return a[i]
			}
			return d
		}
		var t Matrix
		switch name {
		case "matrix":
			if len(a) != 6 {
				return m
			}
			t = Matrix{a[0], a[1], a[2], a[3], a[4], a[5]}
		case "translate":
			t = Matrix{1, 0, 0, 1, arg(0, 0), arg(1, 0)}
		case "scale":
			sx := arg(0, 1)
			t = Matrix{sx, 0, 0, arg(1, sx), 0, 0}
		case "rotate":
			sn, cs := math.Sincos(arg(0, 0) * math.Pi / 180)
			cx, cy := arg(1, 0), arg(2, 0)
			t = Matrix{1, 0, 0, 1, cx, cy}.Mul(Matrix{cs, sn, -sn, cs, 0, 0}).Mul(Matrix{1, 0, 0, 1, -cx, -cy})
		case "skewX":
			t = Matrix{1, 0, math.Tan(arg(0, 0) * math.Pi / 180), 1, 0, 0}
		case "skewY":
			t = Matrix{1, math.Tan(arg(0, 0) * math.Pi / 180), 0, 1, 0, 0}
		default:
			return m
		}
		m = m.Mul(t)
	}
}
