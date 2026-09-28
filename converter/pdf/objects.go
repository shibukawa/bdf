package pdf

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"github.com/pdfcpu/pdfcpu/pkg/filter"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/shibukawa/bdf/imgconv"
)

// Thin helpers over pdfcpu's object model. All of them tolerate missing or
// malformed values and return zero values, because real-world PDFs are messy.

type pdf struct {
	ctx *model.Context
	// maxStream bounds the decoded bytes of a stream (0: maxStreamBytes),
	// maxSamples those of the samples of an image (0: maxImageBytes), and
	// maxPixels the pixels of an image that is not reduced (0:
	// imgconv.MaxDecodePixels).
	maxStream, maxSamples, maxPixels int64
	// treeReads counts the nodes of name and number trees that were read.
	treeReads int
}

// streamLimit returns the bound of the decoded bytes of a stream.
func (p *pdf) streamLimit() int64 {
	if p.maxStream > 0 {
		return p.maxStream
	}
	return maxStreamBytes
}

// sampleLimit returns the bound of the bytes of the samples of an image.
func (p *pdf) sampleLimit() int64 {
	if p.maxSamples > 0 {
		return p.maxSamples
	}
	return maxImageBytes
}

// pixelLimit returns the pixels beyond which an image is reduced.
func (p *pdf) pixelLimit() int64 {
	if p.maxPixels > 0 {
		return p.maxPixels
	}
	return imgconv.MaxDecodePixels
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

// text returns a PDF text string (PDFDocEncoding, or UTF-16BE or UTF-8 with
// a byte order mark) as UTF-8.
func (p *pdf) text(o types.Object) string { return decodeText(p.str(o)) }

// treeBudget bounds the nodes visited by one number or name tree lookup
// (a malformed tree may share or loop its kids).
const treeBudget = 4096

// numberTree returns the value of key in the number tree rooted at o.
func (p *pdf) numberTree(o types.Object, key int) types.Object {
	budget := treeBudget
	var find func(o types.Object) types.Object
	find = func(o types.Object) types.Object {
		d := p.dict(o)
		if d == nil || budget <= 0 {
			return nil
		}
		budget--
		p.treeReads++
		if lim := p.nums(d["Limits"]); len(lim) == 2 && (float64(key) < lim[0] || float64(key) > lim[1]) {
			return nil
		}
		nums := p.array(d["Nums"])
		for i := 0; i+1 < len(nums); i += 2 {
			if k, ok := p.num(nums[i]); ok && int(k) == key {
				return nums[i+1]
			}
		}
		for _, kid := range p.array(d["Kids"]) {
			if v := find(kid); v != nil {
				return v
			}
		}
		return nil
	}
	return find(o)
}

// nameTree returns the names of the name tree rooted at o with their
// values: for a name given more than once, the first value in the order of
// the tree. It is read once for all the lookups, which read it each from
// its root before.
func (p *pdf) nameTree(o types.Object) map[string]types.Object {
	budget := treeBudget
	out := map[string]types.Object{}
	var read func(o types.Object)
	read = func(o types.Object) {
		d := p.dict(o)
		if d == nil || budget <= 0 {
			return
		}
		budget--
		p.treeReads++
		names := p.array(d["Names"])
		for i := 0; i+1 < len(names); i += 2 {
			key := string(p.str(names[i]))
			if _, ok := out[key]; !ok && names[i+1] != nil {
				out[key] = names[i+1]
			}
		}
		for _, kid := range p.array(d["Kids"]) {
			read(kid)
		}
	}
	read(o)
	return out
}

// key identifies an object for caching: its indirect reference if it has one,
// otherwise a pointer-ish string built from the object's printed form.
func objKey(o types.Object) string {
	if r, ok := o.(types.IndirectRef); ok {
		return fmt.Sprintf("%d.%d", r.ObjectNumber, r.GenerationNumber)
	}
	return ""
}

// maxStreamBytes bounds the decoded bytes of a stream: deflate data of a
// thousand bytes can ask for a gigabyte, and twice deflated for far more.
const maxStreamBytes = 256 << 20

// decodeStream returns the decoded bytes of a stream, applying all filters
// except image codecs (DCT, JPX, JBIG2, CCITT), which are reported via the
// returned filter name so the caller can handle them. A stream that decodes
// to more than maxStreamBytes is an error.
func (p *pdf) decodeStream(sd *types.StreamDict) ([]byte, string, error) {
	return p.decodeStreamUpTo(sd, p.streamLimit(), false)
}

// decodeStreamUpTo is decodeStream with another limit for the decoded
// bytes. With cut set, the bytes after the limit are left out instead of
// being an error: the samples of an image, which has no use for more than
// its rows. The data between two filters is held to maxStreamBytes.
func (p *pdf) decodeStreamUpTo(sd *types.StreamDict, limit int64, cut bool) ([]byte, string, error) {
	if sd == nil {
		return nil, "", fmt.Errorf("nil stream")
	}
	data := sd.Raw
	for i, f := range sd.FilterPipeline {
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
		last := true
		for _, g := range sd.FilterPipeline[i+1:] {
			switch g.Name {
			case filter.DCT, filter.JPX, filter.JBIG2, filter.CCITTFax:
			default:
				last = false
			}
		}
		max, cutHere := p.streamLimit(), false
		if last {
			max, cutHere = limit, cut
		}
		out, err := decodeFilter(f.Name, parms, data, max)
		if err != nil {
			return nil, "", fmt.Errorf("%s: %w", f.Name, err)
		}
		if int64(len(out)) > max {
			if !cutHere {
				return nil, "", fmt.Errorf("%s: the stream decodes to more than %d bytes", f.Name, max)
			}
			out = out[:max]
		}
		data = out
	}
	return data, "", nil
}

// decodeFilter decodes data with one filter. It returns a little more than
// max bytes of a stream that is longer, and may leave the rest undecoded.
// pdfcpu's filters take the parameters and the data as they are, and panic
// or never return on some: those are checked or decoded here.
func decodeFilter(name string, parms map[string]int, data []byte, max int64) ([]byte, error) {
	switch name {
	case filter.RunLength:
		return decodeRunLength(data, max), nil
	case filter.ASCII85:
		if len(data) == 0 {
			return nil, nil // pdfcpu looks at the last byte
		}
	case filter.Flate:
		if err := checkPredictor(parms, len(data)); err != nil {
			return nil, err
		}
	}
	fi, err := filter.NewFilter(name, parms)
	if err != nil {
		return nil, err
	}
	switch name {
	case filter.Flate:
		// pdfcpu returns what it read with io.EOF for a stream shorter than
		// the length asked for.
		r, err := fi.DecodeLength(bytesReader(data), max+1)
		if b, ok := r.(*bytes.Buffer); ok && b != nil && (err == nil || err == io.EOF) {
			return b.Bytes(), nil
		}
		if err == nil {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	case filter.LZW:
		// pdfcpu returns nothing but io.EOF for a stream shorter than the
		// length asked for: then the stream is short enough to decode whole.
		r, err := fi.DecodeLength(bytesReader(data), max+1)
		if err == nil {
			return readAll(r)
		}
		if err != io.EOF {
			return nil, err
		}
	}
	// ASCIIHex and ASCII85 give fewer bytes than they take.
	r, err := fi.Decode(bytesReader(data))
	if err != nil {
		return nil, err
	}
	return readAll(r)
}

// checkPredictor checks the parameters of a Flate predictor, which pdfcpu
// takes as they are: it divides by a row of no bytes, or reads such rows for
// ever, and makes its rows before it reads any data. n is the number of
// bytes to decode: deflate gives at most 1032 bytes for one, so that a
// longer row cannot be a row of this stream.
func checkPredictor(parms map[string]int, n int) error {
	if pr, ok := parms["Predictor"]; !ok || pr <= 1 {
		return nil
	}
	get := func(key string, def int) int {
		if v, ok := parms[key]; ok {
			return v
		}
		return def
	}
	columns, colors, bpc := get("Columns", 1), get("Colors", 1), get("BitsPerComponent", 8)
	if columns < 1 || columns > 1<<24 {
		return fmt.Errorf("predictor with %d columns", columns)
	}
	if colors < 1 || colors > 32 {
		return fmt.Errorf("predictor with %d colours", colors)
	}
	switch bpc {
	case 1, 2, 4, 8, 16:
	default:
		return fmt.Errorf("predictor with %d bits per component", bpc)
	}
	if row := (int64(bpc)*int64(colors)*int64(columns) + 7) / 8; row > 1032*int64(n)+64 {
		return fmt.Errorf("predictor with rows of %d bytes for %d bytes of data", row, n)
	}
	return nil
}

// decodeRunLength decodes RunLengthDecode data up to a little more than max
// bytes. Data that ends within a run gives the bytes before its end.
func decodeRunLength(src []byte, max int64) []byte {
	var out []byte
	for i := 0; i < len(src) && int64(len(out)) <= max; {
		b := src[i]
		i++
		switch {
		case b == 0x80: // end of data
			return out
		case b < 0x80: // b+1 bytes as they are
			n := min(int(b)+1, len(src)-i)
			out = append(out, src[i:i+n]...)
			i += n
		case i < len(src): // the next byte 257-b times
			for k := 257 - int(b); k > 0; k-- {
				out = append(out, src[i])
			}
			i++
		}
	}
	return out
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
