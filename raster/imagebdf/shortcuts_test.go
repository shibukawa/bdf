package imagebdf

import (
	"bytes"
	"image"
	"math"
	"math/rand"
	"slices"
	"sort"
	"testing"

	"github.com/shibukawa/bdf"
)

// The tests of this file compare the shortcuts of the renderer with the
// code they stand for; TestShortcutsChangeNoPixel compares them on the
// documents of testdata.

// both calls fn with the shortcuts and without them.
func both(t *testing.T, fn func() []byte) (fast, slow []byte) {
	t.Helper()
	defer func() { plain = false }()
	plain = false
	fast = fn()
	plain = true
	return fast, fn()
}

func maskBytes(m *mask) []byte {
	var b bytes.Buffer
	if m.empty() {
		return nil
	}
	for y := m.r.Min.Y; y < m.r.Max.Y; y++ {
		for x := m.r.Min.X; x < m.r.Max.X; x++ {
			v := math.Float32bits(m.at(x, y))
			b.Write([]byte{byte(x), byte(y), byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24)})
		}
	}
	return b.Bytes()
}

// covered is the coverage of a mask where it is not 0: a mask may have
// more or fewer pixels of none around it.
func covered(m *mask) []byte {
	b := maskBytes(m)
	var out []byte
	for i := 0; i < len(b); i += 6 {
		if !bytes.Equal(b[i+2:i+6], []byte{0, 0, 0, 0}) {
			out = append(out, b[i:i+6]...)
		}
	}
	return out
}

func TestSortOfCrossings(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	b := image.Rect(0, 0, 40, 40)
	for n := 0; n < 40; n++ {
		// polygons in no order, some with edges that cross in the same place
		var polys []polyline
		for k := 0; k < 1+rng.Intn(60); k++ {
			x, y := float64(rng.Intn(60)-10), float64(rng.Intn(60)-10)
			pl := polyline{closed: true}
			for c := 0; c < 3+rng.Intn(4); c++ {
				pl.pts = append(pl.pts, point{x + float64(rng.Intn(17)) - 8, y + rng.Float64()*16 - 8})
			}
			polys = append(polys, pl)
			if rng.Intn(3) == 0 {
				polys = append(polys, pl)
			}
		}
		for _, rule := range []byte{bdf.NonZero, bdf.EvenOdd} {
			sc := &scratch{}
			fast, slow := both(t, func() []byte { return maskBytes(sc.rasterize(polys, rule, b)) })
			if !bytes.Equal(fast, slow) {
				t.Fatalf("%d polygons, rule %d: sorted otherwise than by insertion", len(polys), rule)
			}
			// and once more, with what the scratch kept
			if again := maskBytes(sc.rasterize(polys, rule, b)); !bytes.Equal(again, slow) {
				t.Fatalf("%d polygons: the second coverage of a scratch differs", len(polys))
			}
		}
	}
}

func TestPathsOutsideTheCanvas(t *testing.T) {
	rng := rand.New(rand.NewSource(2))
	r := New(bdf.NewDocument(), nil)
	d, _ := r.canvas(bdf.Rect{W: 32, H: 32}, 32, 32)
	d.ctx.st.clip = &mask{r: image.Rect(3, 2, 30, 31)}
	for n := 0; n < 300; n++ {
		p := &path{}
		for k := 0; k < 1+rng.Intn(4); k++ {
			// about the canvas: inside, outside, across its sides
			x, y := rng.Float64()*120-40, rng.Float64()*120-40
			p.moveTo(x, y)
			for c := 0; c < 1+rng.Intn(3); c++ {
				switch rng.Intn(3) {
				case 0:
					p.lineTo(x+rng.Float64()*30-15, y+rng.Float64()*30-15)
				case 1:
					p.quadTo(x+rng.Float64()*30-15, y+rng.Float64()*30-15, x+rng.Float64()*30-15, y+rng.Float64()*30-15)
				default:
					p.cubicTo(x+rng.Float64()*30-15, y+rng.Float64()*30-15, x+rng.Float64()*30-15, y+rng.Float64()*30-15, x+rng.Float64()*30-15, y+rng.Float64()*30-15)
				}
			}
			if rng.Intn(2) == 0 {
				p.close()
			}
		}
		m := identity.translate(rng.Float64()*8-4, rng.Float64()*8-4).scale(0.5+rng.Float64(), 0.5+rng.Float64())
		if rng.Intn(3) == 0 {
			m = m.mul(matrix{1, 0.3, -0.2, 1, 0, 0})
		}
		fast, slow := both(t, func() []byte { return covered(d.coverage(p, bdf.NonZero, m)) })
		if !bytes.Equal(fast, slow) {
			t.Fatalf("path %d: the coverage differs", n)
		}
		ls := lineStyle{width: 0.5 + rng.Float64()*12, cap: byte(rng.Intn(3)), join: byte(rng.Intn(3)), miter: 1 + rng.Float64()*10}
		if rng.Intn(2) == 0 {
			ls.setDash([]float32{float32(1 + rng.Float64()*6), float32(0.5 + rng.Float64()*4)}, float32(rng.Float64()*5))
		}
		fast, slow = both(t, func() []byte {
			cov, alpha, cut := r.scratch().strokeMask(p, m, ls, d.bounds())
			if cut {
				t.Fatal("cut")
			}
			return append(covered(cov), byte(alpha*255))
		})
		if !bytes.Equal(fast, slow) {
			t.Fatalf("path %d: the coverage of its stroke differs (%+v)", n, ls)
		}
	}
}

func TestFillWithinTheCoverage(t *testing.T) {
	rng := rand.New(rand.NewSource(3))
	random := func(r image.Rectangle) *mask {
		m := &mask{r: r}
		if rng.Intn(3) > 0 {
			m.a = make([]float32, r.Dx()*r.Dy())
			for i := range m.a {
				if v := rng.Float32()*1.5 - 0.25; v > 0 {
					m.a[i] = min(v, 1)
				}
			}
		}
		return m
	}
	rect := func() image.Rectangle {
		x, y := rng.Intn(14)-2, rng.Intn(14)-2
		return image.Rect(x, y, x+1+rng.Intn(12), y+1+rng.Intn(12))
	}
	group := newSurface(9, 7)
	for i := range group.pix {
		group.pix[i] = rng.Float32()
	}
	paint := bdf.LinearGradient(0, 0, 12, 3, bdf.Stop{Offset: 0, Color: 0xff000080}, bdf.Stop{Offset: 1, Color: 0x0000ffff})
	shaders := []shader{solid{0.5, 0.25, 0, 0.5}, solid{0, 0, 1, 1}, surfaceShader{s: group, dx: 2, dy: 3}, newGradientShader(&paint, identity, nil)}
	for n := 0; n < 400; n++ {
		base := newSurface(12, 12)
		for i := range base.pix {
			base.pix[i] = rng.Float32()
		}
		cov, clip := random(rect()), random(rect().Intersect(base.bounds()))
		if rng.Intn(8) == 0 {
			cov = nil
		}
		sh, alpha, mode := shaders[rng.Intn(len(shaders))], rng.Float32(), byte(rng.Intn(len(bdf.BlendNames)))
		if rng.Intn(3) == 0 {
			mode = bdf.BlendSourceOver
		}
		fast, slow := both(t, func() []byte {
			s := &surface{w: base.w, h: base.h, pix: slices.Clone(base.pix)}
			s.fill(cov, sh, alpha, mode, clip)
			// twice: the row of colours is kept
			s.fill(cov, shaders[0], alpha, bdf.BlendSourceOver, clip)
			var b []byte
			for _, v := range s.pix {
				u := math.Float32bits(v)
				b = append(b, byte(u), byte(u>>8), byte(u>>16), byte(u>>24))
			}
			return b
		})
		if !bytes.Equal(fast, slow) {
			t.Fatalf("fill %d (mode %s): the pixels differ", n, bdf.BlendNames[mode])
		}
	}
}

func TestRounding(t *testing.T) {
	// the bytes of an image are the colours times 255, rounded as
	// math.Round does
	check := func(v float32) {
		s := &surface{w: 1, h: 1, pix: []float32{v, 0, 0, 0}}
		want := uint8(math.Round(float64(min(max(v, 0), 1)) * 255))
		if got := s.toRGBA().Pix[0]; got != want {
			t.Fatalf("%v (%#x) is %d, want %d", v, math.Float32bits(v), got, want)
		}
	}
	for i := 0; i <= 255; i++ {
		// about the middle between two bytes, where the rounding turns
		mid := math.Float32bits(float32((float64(i) + 0.5) / 255))
		for d := -64; d <= 64; d++ {
			check(math.Float32frombits(uint32(int(mid) + d)))
		}
		check(float32(i) / 255)
	}
	rng := rand.New(rand.NewSource(4))
	for i := 0; i < 200000; i++ {
		check(math.Float32frombits(rng.Uint32()&0x3fffffff + uint32(rng.Intn(1<<24)))) // to 2 and about
		check(rng.Float32())
	}
	for _, v := range []float32{0, 1, -1, 2, 1e-30, 1e-45, -1e-45, float32(math.Inf(1)), float32(math.Inf(-1)), math.Float32frombits(0x3f7fffff), math.Float32frombits(0x3f800001)} {
		check(v)
	}
}

func TestColoursOfGradients(t *testing.T) {
	// as the colours were worked out before: every stop looked at for
	// every colour
	before := func(stops []bdf.Stop) *gradientLUT {
		type stop struct {
			off float64
			c   [4]float64
		}
		ss := make([]stop, 0, len(stops))
		for _, s := range stops {
			c := uint32(s.Color)
			ss = append(ss, stop{math.Min(1, math.Max(0, float64(s.Offset))), [4]float64{float64(c>>24) / 255, float64(c>>16&0xff) / 255, float64(c>>8&0xff) / 255, float64(c&0xff) / 255}})
		}
		sort.SliceStable(ss, func(i, j int) bool { return ss[i].off < ss[j].off })
		lut := &gradientLUT{}
		for i := range lut {
			t := float64(i) / 1024
			var c [4]float64
			switch {
			case t <= ss[0].off:
				c = ss[0].c
			case t >= ss[len(ss)-1].off:
				c = ss[len(ss)-1].c
			default:
				for k := 1; k < len(ss); k++ {
					if t <= ss[k].off {
						a, b := ss[k-1], ss[k]
						u := 0.0
						if b.off > a.off {
							u = (t - a.off) / (b.off - a.off)
						}
						for j := range c {
							c[j] = a.c[j] + (b.c[j]-a.c[j])*u
						}
						break
					}
				}
			}
			lut[i] = [4]float32{float32(c[0] * c[3]), float32(c[1] * c[3]), float32(c[2] * c[3]), float32(c[3])}
		}
		return lut
	}
	same := func(a, b *gradientLUT) bool {
		for i := range a {
			for j := range a[i] {
				if math.Float32bits(a[i][j]) != math.Float32bits(b[i][j]) {
					return false
				}
			}
		}
		return true
	}
	rng := rand.New(rand.NewSource(5))
	for n := 0; n < 300; n++ {
		stops := make([]bdf.Stop, 1+rng.Intn(12))
		for i := range stops {
			stops[i] = bdf.Stop{Offset: rng.Float32()*1.4 - 0.2, Color: bdf.Color(rng.Uint32())}
			switch rng.Intn(12) {
			case 0:
				stops[i].Offset = float32(math.NaN())
			case 1:
				stops[i].Offset = stops[rng.Intn(i+1)].Offset // two stops in one place
			case 2:
				stops[i].Offset = float32(rng.Intn(3)) / 2
			}
		}
		if rng.Intn(4) == 0 {
			slices.SortFunc(stops, func(a, b bdf.Stop) int { return int(math.Float32bits(a.Offset)) - int(math.Float32bits(b.Offset)) })
		}
		if !same(buildLUT(stops), before(stops)) {
			t.Fatalf("the colours of %v differ", stops)
		}
	}
	// the colours of a gradient of an object are worked out once
	r := New(bdf.NewDocument(), nil)
	paints := make([]bdf.Paint, lutsKept+10)
	for i := range paints {
		paints[i] = bdf.LinearGradient(0, 0, 1, 0, bdf.Stop{Color: bdf.Color(i)<<8 | 0xff}, bdf.Stop{Offset: 1, Color: 0xff})
	}
	if a, b := r.lut(&paints[0]), r.lut(&paints[0]); a != b || !same(a, before(paints[0].Stops)) {
		t.Error("the colours of a gradient are not kept")
	}
	for i := range paints {
		if l := r.lut(&paints[i]); !same(l, before(paints[i].Stops)) || len(r.luts) > lutsKept {
			t.Fatalf("gradient %d: %d kept", i, len(r.luts))
		}
	}
	if r.lut(&bdf.Paint{}) != nil {
		t.Error("colours of a gradient without stops")
	}
}
