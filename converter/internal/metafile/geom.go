package metafile

import (
	"math"

	"github.com/shibukawa/bdf"
)

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
