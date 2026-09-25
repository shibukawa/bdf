package bdf

// Paint kinds.
const (
	PaintLinear  byte = 0
	PaintRadial  byte = 1
	PaintConic   byte = 2
	PaintPattern byte = 3
)

// Pattern repeat modes.
const (
	RepeatBoth byte = 0
	RepeatX    byte = 1
	RepeatY    byte = 2
	NoRepeat   byte = 3
)

// Stop is a gradient color stop.
type Stop struct {
	Offset float32
	Color  Color
}

// Paint is a gradient or pattern fill/stroke style.
type Paint struct {
	Kind   byte
	Coords []float32 // linear: x0 y0 x1 y1; radial: x0 y0 r0 x1 y1 r1; conic: angle x y
	Stops  []Stop
	Image  ImageRef   // pattern only
	Repeat byte       // pattern only
	Matrix [6]float32 // pattern only
}

// LinearGradient builds a linear gradient paint.
func LinearGradient(x0, y0, x1, y1 float32, stops ...Stop) Paint {
	return Paint{Kind: PaintLinear, Coords: []float32{x0, y0, x1, y1}, Stops: stops}
}

// RadialGradient builds a radial gradient paint.
func RadialGradient(x0, y0, r0, x1, y1, r1 float32, stops ...Stop) Paint {
	return Paint{Kind: PaintRadial, Coords: []float32{x0, y0, r0, x1, y1, r1}, Stops: stops}
}

// ConicGradient builds a conic gradient paint.
func ConicGradient(angle, x, y float32, stops ...Stop) Paint {
	return Paint{Kind: PaintConic, Coords: []float32{angle, x, y}, Stops: stops}
}

// Pattern builds an image pattern paint with an identity matrix.
func Pattern(img ImageRef, repeat byte) Paint {
	return Paint{Kind: PaintPattern, Image: img, Repeat: repeat, Matrix: [6]float32{1, 0, 0, 1, 0, 0}}
}

var paintCoords = [...]int{4, 6, 3}

func (p *Paint) encode(w *buf) {
	w.u8(p.Kind)
	if p.Kind == PaintPattern {
		w.varuint(uint64(p.Image))
		w.u8(p.Repeat)
		w.f32s(p.Matrix[:]...)
		return
	}
	w.f32s(p.Coords...)
	w.varuint(uint64(len(p.Stops)))
	for _, s := range p.Stops {
		w.f32(s.Offset)
		w.u32(uint32(s.Color))
	}
}

func decodePaint(r *reader) Paint {
	p := Paint{Kind: r.u8()}
	if p.Kind == PaintPattern {
		p.Image = ImageRef(r.varuint())
		p.Repeat = r.u8()
		for i := range p.Matrix {
			p.Matrix[i] = r.f32()
		}
		return p
	}
	if int(p.Kind) >= len(paintCoords) {
		r.fail("unknown paint kind")
		return p
	}
	p.Coords = make([]float32, paintCoords[p.Kind])
	for i := range p.Coords {
		p.Coords[i] = r.f32()
	}
	n := int(r.varuint())
	if n > len(r.b) {
		r.fail("bad stop count")
		return p
	}
	p.Stops = make([]Stop, n)
	for i := range p.Stops {
		p.Stops[i] = Stop{Offset: r.f32(), Color: Color(r.u32())}
	}
	return p
}
