package pdf2bdf

import (
	"fmt"
	"math"

	"github.com/pdfcpu/pdfcpu/pkg/filter"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// Thin helpers over pdfcpu's object model. All of them tolerate missing or
// malformed values and return zero values, because real-world PDFs are messy.

type pdf struct {
	ctx *model.Context
}

func (p *pdf) deref(o types.Object) types.Object {
	if o == nil {
		return nil
	}
	v, err := p.ctx.Dereference(o)
	if err != nil {
		return nil
	}
	return v
}

func (p *pdf) dict(o types.Object) types.Dict {
	switch v := p.deref(o).(type) {
	case types.Dict:
		return v
	case types.StreamDict:
		return v.Dict
	case *types.StreamDict:
		return v.Dict
	}
	return nil
}

func (p *pdf) array(o types.Object) types.Array {
	if a, ok := p.deref(o).(types.Array); ok {
		return a
	}
	return nil
}

func (p *pdf) name(o types.Object) string {
	if n, ok := p.deref(o).(types.Name); ok {
		return n.Value()
	}
	return ""
}

func (p *pdf) num(o types.Object) (float64, bool) {
	switch v := p.deref(o).(type) {
	case types.Integer:
		return float64(v), true
	case types.Float:
		return float64(v), true
	}
	return 0, false
}

func (p *pdf) numOr(o types.Object, def float64) float64 {
	if v, ok := p.num(o); ok {
		return v
	}
	return def
}

func (p *pdf) intOr(o types.Object, def int) int {
	if v, ok := p.num(o); ok {
		return int(v)
	}
	return def
}

func (p *pdf) boolOr(o types.Object, def bool) bool {
	if b, ok := p.deref(o).(types.Boolean); ok {
		return bool(b)
	}
	return def
}

func (p *pdf) nums(o types.Object) []float64 {
	a := p.array(o)
	if a == nil {
		return nil
	}
	out := make([]float64, 0, len(a))
	for _, e := range a {
		v, _ := p.num(e)
		out = append(out, v)
	}
	return out
}

func (p *pdf) stream(o types.Object) *types.StreamDict {
	switch v := p.deref(o).(type) {
	case types.StreamDict:
		return &v
	case *types.StreamDict:
		return v
	}
	return nil
}

// str returns the bytes of a string object.
func (p *pdf) str(o types.Object) []byte {
	switch v := p.deref(o).(type) {
	case types.StringLiteral:
		b, err := types.Unescape(v.Value())
		if err != nil {
			return []byte(v.Value())
		}
		return b
	case types.HexLiteral:
		b, err := v.Bytes()
		if err != nil {
			return nil
		}
		return b
	}
	return nil
}

// text returns a PDF text string (PDFDocEncoding or UTF-16) as UTF-8.
func (p *pdf) text(o types.Object) string {
	b := p.str(o)
	if len(b) >= 2 && b[0] == 0xfe && b[1] == 0xff {
		var out []rune
		for i := 2; i+1 < len(b); i += 2 {
			u := rune(b[i])<<8 | rune(b[i+1])
			if u >= 0xd800 && u < 0xdc00 && i+3 < len(b) {
				lo := rune(b[i+2])<<8 | rune(b[i+3])
				u = 0x10000 + (u-0xd800)<<10 + (lo - 0xdc00)
				i += 2
			}
			out = append(out, u)
		}
		return string(out)
	}
	return string(b)
}

// key identifies an object for caching: its indirect reference if it has one,
// otherwise a pointer-ish string built from the object's printed form.
func objKey(o types.Object) string {
	if r, ok := o.(types.IndirectRef); ok {
		return fmt.Sprintf("%d.%d", r.ObjectNumber, r.GenerationNumber)
	}
	return ""
}

// decodeStream returns the decoded bytes of a stream, applying all filters
// except image codecs (DCT, JPX, JBIG2, CCITT), which are reported via the
// returned filter name so the caller can handle them.
func (p *pdf) decodeStream(sd *types.StreamDict) ([]byte, string, error) {
	if sd == nil {
		return nil, "", fmt.Errorf("nil stream")
	}
	data := sd.Raw
	for _, f := range sd.FilterPipeline {
		switch f.Name {
		case filter.DCT, filter.JPX, filter.JBIG2:
			return data, f.Name, nil
		case filter.CCITTFax:
			// pdfcpu can decode CCITT but needs image parameters; let the image code do it.
			return data, f.Name, nil
		}
		parms := map[string]int{}
		if f.DecodeParms != nil {
			for k, v := range f.DecodeParms {
				if n, ok := p.num(v); ok {
					parms[k] = int(n)
				}
			}
		}
		fi, err := filter.NewFilter(f.Name, parms)
		if err != nil {
			return nil, "", err
		}
		r, err := fi.Decode(bytesReader(data))
		if err != nil {
			return nil, "", fmt.Errorf("%s: %w", f.Name, err)
		}
		data, err = readAll(r)
		if err != nil {
			return nil, "", err
		}
	}
	return data, "", nil
}

// matrix is a PDF/canvas affine matrix [a b c d e f].
type matrix [6]float64

var identity = matrix{1, 0, 0, 1, 0, 0}

// mul returns m × n: apply m first, then n (PDF convention).
func (m matrix) mul(n matrix) matrix {
	return matrix{
		m[0]*n[0] + m[1]*n[2],
		m[0]*n[1] + m[1]*n[3],
		m[2]*n[0] + m[3]*n[2],
		m[2]*n[1] + m[3]*n[3],
		m[4]*n[0] + m[5]*n[2] + n[4],
		m[4]*n[1] + m[5]*n[3] + n[5],
	}
}

func (m matrix) apply(x, y float64) (float64, float64) {
	return m[0]*x + m[2]*y + m[4], m[1]*x + m[3]*y + m[5]
}

func (m matrix) inverse() (matrix, bool) {
	det := m[0]*m[3] - m[1]*m[2]
	if det == 0 || math.IsNaN(det) || math.IsInf(det, 0) {
		return identity, false
	}
	return matrix{
		m[3] / det, -m[1] / det,
		-m[2] / det, m[0] / det,
		(m[2]*m[5] - m[3]*m[4]) / det,
		(m[1]*m[4] - m[0]*m[5]) / det,
	}, true
}

func (m matrix) isIdentity() bool { return m == identity }

// scaleFactor is the geometric mean scale of the matrix.
func (m matrix) scaleFactor() float64 {
	return math.Sqrt(math.Abs(m[0]*m[3] - m[1]*m[2]))
}

// rect is an axis-aligned box.
type rect struct{ x0, y0, x1, y1 float64 }

func (r rect) empty() bool { return r.x1 <= r.x0 || r.y1 <= r.y0 }

func (r rect) intersect(o rect) rect {
	return rect{math.Max(r.x0, o.x0), math.Max(r.y0, o.y0), math.Min(r.x1, o.x1), math.Min(r.y1, o.y1)}
}

func (r rect) union(o rect) rect {
	if r.empty() {
		return o
	}
	if o.empty() {
		return r
	}
	return rect{math.Min(r.x0, o.x0), math.Min(r.y0, o.y0), math.Max(r.x1, o.x1), math.Max(r.y1, o.y1)}
}

// transformRect returns the bounding box of the rectangle's corners after m.
func (m matrix) transformRect(r rect) rect {
	xs := [4]float64{}
	ys := [4]float64{}
	xs[0], ys[0] = m.apply(r.x0, r.y0)
	xs[1], ys[1] = m.apply(r.x1, r.y0)
	xs[2], ys[2] = m.apply(r.x0, r.y1)
	xs[3], ys[3] = m.apply(r.x1, r.y1)
	out := rect{xs[0], ys[0], xs[0], ys[0]}
	for i := 1; i < 4; i++ {
		out.x0 = math.Min(out.x0, xs[i])
		out.y0 = math.Min(out.y0, ys[i])
		out.x1 = math.Max(out.x1, xs[i])
		out.y1 = math.Max(out.y1, ys[i])
	}
	return out
}

func (p *pdf) matrixOr(o types.Object, def matrix) matrix {
	v := p.nums(o)
	if len(v) != 6 {
		return def
	}
	return matrix{v[0], v[1], v[2], v[3], v[4], v[5]}
}

func (p *pdf) rectOr(o types.Object, def rect) rect {
	v := p.nums(o)
	if len(v) != 4 {
		return def
	}
	return rect{math.Min(v[0], v[2]), math.Min(v[1], v[3]), math.Max(v[0], v[2]), math.Max(v[1], v[3])}
}
