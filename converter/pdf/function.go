package pdf

import (
	"math"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// pdfFunction evaluates PDF functions of type 0 (sampled), 2 (exponential),
// 3 (stitching) and 4 (PostScript calculator; mid-range values when the
// program cannot be run).
type pdfFunction struct {
	kind   int
	domain []float64
	rng    []float64
	// type 2
	c0, c1 []float64
	n      float64
	// type 3
	funcs  []*pdfFunction
	bounds []float64
	encode []float64
	// type 0
	size    []int
	bps     int
	decode  []float64
	samples []byte
	nOut    int
	// type 4
	prog  []psOp
	cache map[float64][]float64 // results for one input (image samples repeat)
	// arrays of functions (one per output)
	parts []*pdfFunction
}

func (p *pdf) loadFunction(o types.Object) *pdfFunction {
	if a := p.array(o); a != nil && len(a) > 0 {
		if _, isNum := p.num(a[0]); !isNum {
			f := &pdfFunction{}
			for _, e := range a {
				if sub := p.loadFunction(e); sub != nil {
					f.parts = append(f.parts, sub)
				}
			}
			if len(f.parts) == 0 {
				return nil
			}
			f.domain = f.parts[0].domain
			return f
		}
	}
	d := p.dict(o)
	if d == nil {
		return nil
	}
	f := &pdfFunction{kind: p.intOr(d["FunctionType"], -1), domain: p.nums(d["Domain"]), rng: p.nums(d["Range"])}
	if len(f.domain) < 2 {
		f.domain = []float64{0, 1}
	}
	switch f.kind {
	case 2:
		f.c0 = p.nums(d["C0"])
		f.c1 = p.nums(d["C1"])
		if f.c0 == nil {
			f.c0 = []float64{0}
		}
		if f.c1 == nil {
			f.c1 = []float64{1}
		}
		f.n = p.numOr(d["N"], 1)
	case 3:
		for _, e := range p.array(d["Functions"]) {
			f.funcs = append(f.funcs, p.loadFunction(e))
		}
		f.bounds = p.nums(d["Bounds"])
		f.encode = p.nums(d["Encode"])
	case 0:
		sd := p.stream(o)
		if sd == nil {
			return nil
		}
		data, _, err := p.decodeStream(sd)
		if err != nil {
			return nil
		}
		f.samples = data
		for _, s := range p.nums(d["Size"]) {
			f.size = append(f.size, int(s))
		}
		f.bps = p.intOr(d["BitsPerSample"], 8)
		f.encode = p.nums(d["Encode"])
		f.decode = p.nums(d["Decode"])
		f.nOut = len(f.rng) / 2
		if len(f.size) == 0 || f.nOut == 0 {
			return nil
		}
	case 4:
		if sd := p.stream(o); sd != nil {
			if data, _, err := p.decodeStream(sd); err == nil {
				f.prog = compilePS(data)
			}
		}
	default:
		return nil
	}
	return f
}

func interpolate(x, xmin, xmax, ymin, ymax float64) float64 {
	if xmax == xmin {
		return ymin
	}
	return ymin + (x-xmin)*(ymax-ymin)/(xmax-xmin)
}

func clamp(v, lo, hi float64) float64 {
	if lo > hi {
		lo, hi = hi, lo
	}
	return math.Max(lo, math.Min(hi, v))
}

// eval evaluates the function for the given inputs.
func (f *pdfFunction) eval(in ...float64) []float64 {
	if f == nil {
		return nil
	}
	if f.parts != nil {
		var out []float64
		for _, sub := range f.parts {
			r := sub.eval(in...)
			if len(r) > 0 {
				out = append(out, r[0])
			}
		}
		return out
	}
	for i := range in {
		if 2*i+1 < len(f.domain) {
			in[i] = clamp(in[i], f.domain[2*i], f.domain[2*i+1])
		}
	}
	var out []float64
	switch f.kind {
	case 2:
		x := in[0]
		t := x
		if f.n != 1 {
			t = math.Pow(x, f.n)
		}
		n := max(len(f.c0), len(f.c1))
		out = make([]float64, n)
		for i := 0; i < n; i++ {
			c0, c1 := 0.0, 1.0
			if i < len(f.c0) {
				c0 = f.c0[i]
			}
			if i < len(f.c1) {
				c1 = f.c1[i]
			}
			out[i] = c0 + t*(c1-c0)
		}
	case 3:
		x := in[0]
		d0, d1 := f.domain[0], f.domain[1]
		k := 0
		for k < len(f.bounds) && x >= f.bounds[k] {
			k++
		}
		if k >= len(f.funcs) {
			k = len(f.funcs) - 1
		}
		if k < 0 {
			return nil
		}
		lo, hi := d0, d1
		if k > 0 {
			lo = f.bounds[k-1]
		}
		if k < len(f.bounds) {
			hi = f.bounds[k]
		}
		e0, e1 := 0.0, 1.0
		if 2*k+1 < len(f.encode) {
			e0, e1 = f.encode[2*k], f.encode[2*k+1]
		}
		out = f.funcs[k].eval(interpolate(x, lo, hi, e0, e1))
	case 0:
		out = f.evalSampled(in)
	case 4:
		out = f.evalPS(in)
	}
	if len(f.rng) >= 2 {
		for i := range out {
			if 2*i+1 < len(f.rng) {
				out[i] = clamp(out[i], f.rng[2*i], f.rng[2*i+1])
			}
		}
	}
	return out
}

func (f *pdfFunction) sample(idx []int, j int) float64 {
	// Sample index in row-major order with the first dimension varying fastest.
	offset := 0
	stride := 1
	for i, s := range f.size {
		offset += idx[i] * stride
		stride *= s
	}
	bit := (offset*f.nOut + j) * f.bps
	maxv := float64(uint64(1)<<uint(f.bps) - 1)
	switch f.bps {
	case 8:
		if bit/8 < len(f.samples) {
			return float64(f.samples[bit/8]) / maxv
		}
	case 16:
		if bit/8+1 < len(f.samples) {
			return float64(uint16(f.samples[bit/8])<<8|uint16(f.samples[bit/8+1])) / maxv
		}
	case 32:
		if bit/8+3 < len(f.samples) {
			v := uint32(f.samples[bit/8])<<24 | uint32(f.samples[bit/8+1])<<16 | uint32(f.samples[bit/8+2])<<8 | uint32(f.samples[bit/8+3])
			return float64(v) / maxv
		}
	default:
		var v uint64
		for k := 0; k < f.bps; k++ {
			b := bit + k
			if b/8 >= len(f.samples) {
				return 0
			}
			v = v<<1 | uint64(f.samples[b/8]>>(7-uint(b%8))&1)
		}
		return float64(v) / maxv
	}
	return 0
}

func (f *pdfFunction) evalSampled(in []float64) []float64 {
	m := len(f.size)
	if len(in) < m {
		return nil
	}
	// Multilinear interpolation on the first input only (common case); nearest for the rest.
	idx := make([]int, m)
	var frac0 float64
	for i := 0; i < m; i++ {
		d0, d1 := f.domain[2*i], f.domain[2*i+1]
		e0, e1 := 0.0, float64(f.size[i]-1)
		if 2*i+1 < len(f.encode) {
			e0, e1 = f.encode[2*i], f.encode[2*i+1]
		}
		e := clamp(interpolate(in[i], d0, d1, e0, e1), 0, float64(f.size[i]-1))
		if i == 0 {
			idx[i] = int(math.Floor(e))
			frac0 = e - float64(idx[i])
		} else {
			idx[i] = int(math.Round(e))
		}
	}
	out := make([]float64, f.nOut)
	for j := 0; j < f.nOut; j++ {
		s0 := f.sample(idx, j)
		s := s0
		if frac0 > 0 && idx[0]+1 < f.size[0] {
			idx2 := append([]int(nil), idx...)
			idx2[0]++
			s = s0 + frac0*(f.sample(idx2, j)-s0)
		}
		dmin, dmax := 0.0, 1.0
		if 2*j+1 < len(f.decode) {
			dmin, dmax = f.decode[2*j], f.decode[2*j+1]
		} else if 2*j+1 < len(f.rng) {
			dmin, dmax = f.rng[2*j], f.rng[2*j+1]
		}
		out[j] = dmin + s*(dmax-dmin)
	}
	return out
}

// evalPS runs a calculator function; its outputs are the top of the stack.
func (f *pdfFunction) evalPS(in []float64) []float64 {
	n := len(f.rng) / 2
	if len(in) == 1 && f.cache != nil {
		if r, ok := f.cache[in[0]]; ok {
			return append([]float64(nil), r...)
		}
	}
	out := runPS(f.prog, in)
	if f.prog == nil || len(out) < n {
		out = make([]float64, n)
		for i := range out {
			out[i] = (f.rng[2*i] + f.rng[2*i+1]) / 2
		}
		return out
	}
	out = out[len(out)-n:]
	if len(in) == 1 {
		if f.cache == nil {
			f.cache = map[float64][]float64{}
		}
		if len(f.cache) < 4096 {
			f.cache[in[0]] = append([]float64(nil), out...)
		}
	}
	return out
}
