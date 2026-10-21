package export

import (
	"bytes"
	"math"

	"github.com/shibukawa/bdf"
)

// Stop is a stop of a gradient: its offset and its color, as fractions.
type Stop struct {
	Off     float64
	R, G, B float64
	A       float64
}

// StopsOf reads the stops of a gradient in their order: a stop before the
// one before it is at its place, and one past 1 at 1.
func StopsOf(p *bdf.Paint) []Stop {
	var out []Stop
	last := 0.0
	for _, s := range p.Stops {
		off := float64(s.Offset)
		if math.IsNaN(off) {
			continue
		}
		off = math.Min(math.Max(off, last), 1)
		last = off
		c := uint32(s.Color)
		out = append(out, Stop{off, float64(c>>24) / 255, float64(c>>16&0xff) / 255, float64(c>>8&0xff) / 255, float64(c&0xff) / 255})
	}
	return out
}

// ClampShadow bounds a number of a shadow (its blur, its offsets).
func ClampShadow(v float64) float64 {
	if math.IsNaN(v) {
		return 0
	}
	return math.Min(math.Max(v, -1e6), 1e6)
}

// IsSVG tells an SVG document from the other images a document holds.
func IsSVG(b []byte) bool {
	b = bytes.TrimPrefix(b, []byte("\xef\xbb\xbf"))
	b = bytes.TrimLeft(b, " \t\r\n")
	return len(b) > 0 && b[0] == '<'
}

// F64 is a number argument of an instruction, 0 when it is not finite.
func F64(in bdf.Instr, i int) float64 {
	if v, ok := in.Args[i].(float32); ok {
		if f := float64(v); !math.IsNaN(f) && !math.IsInf(f, 0) {
			return f
		}
	}
	return 0
}

// U64 is an index argument of an instruction.
func U64(in bdf.Instr, i int) uint64 {
	if v, ok := in.Args[i].(uint64); ok {
		return v
	}
	return 0
}
