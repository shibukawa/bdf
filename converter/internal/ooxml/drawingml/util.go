package drawingml

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
)

func atof(s string, def float64) float64 {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return def
	}
	return f
}

func ftoa(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }

// parsePct parses an ST_Percentage value into a fraction.
func parsePct(v string) float64 {
	v = strings.TrimSpace(v)
	if strings.HasSuffix(v, "%") {
		return atof(strings.TrimSuffix(v, "%"), 0) / 100
	}
	return atof(v, 0) / 100000
}

func itoa(i int) string { return strconv.Itoa(i) }

func f32(v float64) float32 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return float32(v)
}

// pathB builds a bdf.Path from float64 coordinates.
type pathB struct{ p *bdf.Path }

func newPath() *pathB { return &pathB{&bdf.Path{}} }

func (b *pathB) moveTo(x, y float64) *pathB { b.p.MoveTo(f32(x), f32(y)); return b }
func (b *pathB) lineTo(x, y float64) *pathB { b.p.LineTo(f32(x), f32(y)); return b }
func (b *pathB) close() *pathB              { b.p.Close(); return b }
