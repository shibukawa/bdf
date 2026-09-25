package pdf2bdf

import (
	"math"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/shibukawa/bdf"
)

// colorSpace converts component values to sRGB.
type colorSpace struct {
	family string // DeviceGray, DeviceRGB, DeviceCMYK, Indexed, Separation, Lab, Pattern
	n      int    // number of components
	// Indexed
	base   *colorSpace
	lookup []byte
	hival  int
	// Separation / DeviceN
	tint *pdfFunction
	alt  *colorSpace
	// Lab
	wp     []float64
	labRng []float64
	// Pattern with an underlying space (uncolored patterns)
	under *colorSpace
}

var (
	csGray = &colorSpace{family: "DeviceGray", n: 1}
	csRGB  = &colorSpace{family: "DeviceRGB", n: 3}
	csCMYK = &colorSpace{family: "DeviceCMYK", n: 4}
)

func (cs *colorSpace) initial() []float64 {
	switch cs.family {
	case "DeviceRGB":
		return []float64{0, 0, 0}
	case "DeviceCMYK":
		return []float64{0, 0, 0, 1}
	case "Lab":
		return []float64{0, 0, 0}
	case "Indexed":
		return []float64{0}
	case "Separation", "DeviceN":
		out := make([]float64, cs.n)
		for i := range out {
			out[i] = 1
		}
		return out
	}
	return make([]float64, max(cs.n, 1))
}

func to8(v float64) uint8 {
	return uint8(math.Round(clamp(v, 0, 1) * 255))
}

// rgb converts components to an opaque color.
func (cs *colorSpace) rgb(v []float64) bdf.Color {
	r, g, b := cs.rgbf(v)
	return bdf.RGB(to8(r), to8(g), to8(b))
}

func (cs *colorSpace) rgbf(v []float64) (float64, float64, float64) {
	switch cs.family {
	case "DeviceGray":
		if len(v) < 1 {
			return 0, 0, 0
		}
		return v[0], v[0], v[0]
	case "DeviceRGB":
		if len(v) < 3 {
			return 0, 0, 0
		}
		return v[0], v[1], v[2]
	case "DeviceCMYK":
		if len(v) < 4 {
			return 0, 0, 0
		}
		c, m, y, k := clamp(v[0], 0, 1), clamp(v[1], 0, 1), clamp(v[2], 0, 1), clamp(v[3], 0, 1)
		return (1 - c) * (1 - k), (1 - m) * (1 - k), (1 - y) * (1 - k)
	case "Indexed":
		if len(v) < 1 || cs.base == nil {
			return 0, 0, 0
		}
		i := int(math.Round(v[0]))
		if i < 0 {
			i = 0
		}
		if i > cs.hival {
			i = cs.hival
		}
		n := cs.base.n
		comps := make([]float64, n)
		for k := 0; k < n; k++ {
			if i*n+k < len(cs.lookup) {
				comps[k] = float64(cs.lookup[i*n+k]) / 255
			}
		}
		if cs.base.family == "Lab" {
			comps = cs.base.decodeLabBytes(comps)
		}
		return cs.base.rgbf(comps)
	case "Separation", "DeviceN":
		if cs.tint != nil && cs.alt != nil {
			out := cs.tint.eval(append([]float64(nil), v...)...)
			if len(out) >= cs.alt.n {
				return cs.alt.rgbf(out)
			}
		}
		// Approximate: 1 = full ink = dark.
		t := 0.0
		if len(v) > 0 {
			t = clamp(v[0], 0, 1)
		}
		return 1 - t, 1 - t, 1 - t
	case "Lab":
		return labToRGB(v, cs.wp)
	case "Pattern":
		return 0.5, 0.5, 0.5
	}
	return 0, 0, 0
}

// decodeLabBytes maps 0..1 byte-scaled Lab components to Lab ranges.
func (cs *colorSpace) decodeLabBytes(c []float64) []float64 {
	r := cs.labRng
	if len(r) < 4 {
		r = []float64{-100, 100, -100, 100}
	}
	if len(c) < 3 {
		return c
	}
	return []float64{c[0] * 100, r[0] + c[1]*(r[1]-r[0]), r[2] + c[2]*(r[3]-r[2])}
}

func labToRGB(v []float64, wp []float64) (float64, float64, float64) {
	if len(v) < 3 {
		return 0, 0, 0
	}
	xw, yw, zw := 0.9505, 1.0, 1.089
	if len(wp) >= 3 {
		xw, yw, zw = wp[0], wp[1], wp[2]
	}
	L, a, b := v[0], v[1], v[2]
	m := (L + 16) / 116
	l := m + a/500
	n := m - b/200
	g := func(x float64) float64 {
		if x >= 6.0/29 {
			return x * x * x
		}
		return 108.0 / 841 * (x - 4.0/29)
	}
	X, Y, Z := xw*g(l), yw*g(m), zw*g(n)
	r := 3.2406*X - 1.5372*Y - 0.4986*Z
	gg := -0.9689*X + 1.8758*Y + 0.0415*Z
	bb := 0.0557*X - 0.2040*Y + 1.0570*Z
	gam := func(c float64) float64 {
		c = clamp(c, 0, 1)
		if c <= 0.0031308 {
			return 12.92 * c
		}
		return 1.055*math.Pow(c, 1/2.4) - 0.055
	}
	return gam(r), gam(gg), gam(bb)
}

// loadColorSpace resolves a colour space object or name (looked up in resources).
func (c *converter) loadColorSpace(o types.Object, res types.Dict) *colorSpace {
	p := c.pdf
	switch v := p.deref(o).(type) {
	case types.Name:
		switch v.Value() {
		case "DeviceGray", "G", "CalGray":
			return csGray
		case "DeviceRGB", "RGB", "CalRGB":
			return csRGB
		case "DeviceCMYK", "CMYK":
			return csCMYK
		case "Pattern":
			return &colorSpace{family: "Pattern", n: 1}
		case "Indexed", "I":
			return csGray
		}
		if res != nil {
			if csd := p.dict(res["ColorSpace"]); csd != nil {
				if e, ok := csd[v.Value()]; ok {
					return c.loadColorSpace(e, nil)
				}
			}
		}
		c.warnf("unknown colour space %s", v.Value())
		return csGray
	case types.Array:
		if len(v) == 0 {
			return csGray
		}
		fam := p.name(v[0])
		switch fam {
		case "DeviceGray", "G", "DeviceRGB", "RGB", "DeviceCMYK", "CMYK", "CalGray", "CalRGB":
			return c.loadColorSpace(v[0], res)
		case "ICCBased":
			n := 3
			if len(v) > 1 {
				if sd := p.stream(v[1]); sd != nil {
					n = p.intOr(sd.Dict["N"], 0)
					if n == 0 {
						if alt, ok := sd.Dict["Alternate"]; ok {
							return c.loadColorSpace(alt, res)
						}
						n = 3
					}
				}
			}
			switch n {
			case 1:
				return csGray
			case 4:
				return csCMYK
			default:
				return csRGB
			}
		case "Indexed", "I":
			if len(v) < 4 {
				return csGray
			}
			cs := &colorSpace{family: "Indexed", n: 1, base: c.loadColorSpace(v[1], res), hival: p.intOr(v[2], 0)}
			if sd := p.stream(v[3]); sd != nil {
				data, _, err := p.decodeStream(sd)
				if err == nil {
					cs.lookup = data
				}
			} else {
				cs.lookup = p.str(v[3])
			}
			return cs
		case "Separation", "DeviceN":
			cs := &colorSpace{family: fam, n: 1}
			if fam == "DeviceN" && len(v) > 1 {
				cs.n = len(p.array(v[1]))
			}
			if len(v) > 2 {
				cs.alt = c.loadColorSpace(v[2], res)
			}
			if len(v) > 3 {
				cs.tint = p.loadFunction(v[3])
			}
			if fam == "Separation" && len(v) > 1 && p.name(v[1]) == "None" {
				cs.family = "None"
			}
			return cs
		case "Lab":
			cs := &colorSpace{family: "Lab", n: 3}
			if len(v) > 1 {
				d := p.dict(v[1])
				cs.wp = p.nums(d["WhitePoint"])
				cs.labRng = p.nums(d["Range"])
			}
			return cs
		case "Pattern":
			cs := &colorSpace{family: "Pattern", n: 1}
			if len(v) > 1 {
				cs.under = c.loadColorSpace(v[1], res)
			}
			return cs
		}
		c.warnf("unsupported colour space %s", fam)
		return csRGB
	}
	return csGray
}
