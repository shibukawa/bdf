package imagebdf

import (
	"math"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/image/barcode"
)

// Barcode is a QR code or a barcode that a BARCODE mark of a page tells
// (docs/spec.md §7.8), where the image has it: see Options.Barcode.
type Barcode struct {
	// Mark is what the mark says: the symbol, and its box on the page.
	Mark *barcode.Mark
	// X, Y, W and H are the box of the mark on the image, in pixels: its
	// top left corner as the symbol is upright, and its size along the
	// symbol's bars (W) and up them (H), before it is turned.
	X, Y, W, H float64
	// Turns is the quarter turns clockwise of the symbol on the image:
	// 1 has the top of the symbol on the right, its box going down from
	// (X, Y) and to the left. The box's other corners are where the turns
	// take (X+W, Y) and (X, Y+H).
	Turns int
}

// isDraw tells the instructions that draw: the one after a BARCODE mark
// is the drawing of its symbol.
func isDraw(op byte) bool {
	switch op {
	case bdf.OpFillRect, bdf.OpStrokeRect, bdf.OpClearRect, bdf.OpFillPath, bdf.OpStrokePath, bdf.OpFillPathAt,
		bdf.OpFillPathRun, bdf.OpFillText, bdf.OpStrokeText, bdf.OpImage, bdf.OpImageSub, bdf.OpUse, bdf.OpUseAt:
		return true
	}
	return false
}

// barcodeMark notes the mark of the symbol the next drawing draws, when
// the options do something with symbols.
func (d *drawer) barcodeMark(payload string) {
	if d.r.opts.Barcode == nil && !d.r.opts.AlignBarcodes {
		return
	}
	m, err := barcode.ParseMark(payload)
	if err != nil {
		d.r.warnf("%v", err)
		return
	}
	d.symbol = m
}

// frame returns the transform from the upright box of a mark (its corner
// at 0, 0, one unit a pixel) to the image, its corner on a whole pixel,
// the box's size in pixels, and the quarter turns of the symbol: for a
// transform that keeps the axes and does not mirror them.
func frame(m matrix, box [4]float32) (f matrix, w, h float64, turns int, ok bool) {
	if !m.Finite() || m.Det() <= 0 {
		return matrix{}, 0, 0, 0, false
	}
	sign := func(v float64) float64 {
		if v < 0 {
			return -1
		}
		return 1
	}
	var ex, ey [2]float64
	switch {
	case m[1] == 0 && m[2] == 0:
		ex, ey = [2]float64{sign(m[0]), 0}, [2]float64{0, sign(m[3])}
		w, h = math.Abs(m[0])*float64(box[2]), math.Abs(m[3])*float64(box[3])
	case m[0] == 0 && m[3] == 0:
		ex, ey = [2]float64{0, sign(m[1])}, [2]float64{sign(m[2]), 0}
		w, h = math.Abs(m[1])*float64(box[2]), math.Abs(m[2])*float64(box[3])
	default:
		return matrix{}, 0, 0, 0, false
	}
	switch {
	case ex[0] > 0:
		turns = 0
	case ex[1] > 0:
		turns = 1
	case ex[0] < 0:
		turns = 2
	default:
		turns = 3
	}
	x, y := m.Apply(float64(box[0]), float64(box[1]))
	return matrix{ex[0], ex[1], ey[0], ey[1], math.Round(x), math.Round(y)}, w, h, turns, true
}

// drawSymbol does what the options say with the symbol of a mark that the
// drawing after it draws, and reports whether the drawing is to be left
// out: the printer prints the symbol (Options.Barcode), or it is drawn
// again with its modules on whole pixels (Options.AlignBarcodes).
func (d *drawer) drawSymbol(mk *barcode.Mark) bool {
	st := &d.ctx.st
	f, w, h, turns, ok := frame(st.m, mk.Box)
	if !ok || w < 1 || h < 1 {
		return false
	}
	if cb := d.r.opts.Barcode; cb != nil {
		x, y := float64(mk.Box[0]), float64(mk.Box[1])
		x, y = st.m.Apply(x, y)
		if cb(Barcode{Mark: mk, X: x, Y: y, W: w, H: h, Turns: turns}) {
			return true
		}
	}
	if !d.r.opts.AlignBarcodes {
		return false
	}
	sym, err := mk.Symbol()
	if err != nil {
		d.r.warnf("%v", err)
		return false
	}
	// the modules a whole number of pixels wide (and high, for a QR
	// code), as many as fit in the box, and the symbol in its middle
	mod := math.Floor(w / float64(sym.W))
	if !sym.Linear {
		mod = math.Min(mod, math.Floor(h/float64(sym.H)))
	}
	if mod < 1 {
		return false
	}
	sw, sh := mod*float64(sym.W), h
	if !sym.Linear {
		sh = mod * float64(sym.H)
	}
	ox, oy := math.Floor((w-sw)/2), math.Floor((h-sh)/2)
	p := buildPath(sym.Path(float32(ox), float32(oy), float32(sw), float32(sh)))
	d.fillPath(p, bdf.NonZero, f)
	return true
}

// bayer8 is the thresholds of an ordered dither of 8 × 8 pixels, in
// 64ths.
var bayer8 = [64]uint8{
	0, 32, 8, 40, 2, 34, 10, 42,
	48, 16, 56, 24, 50, 18, 58, 26,
	12, 44, 4, 36, 14, 46, 6, 38,
	60, 28, 52, 20, 62, 30, 54, 22,
	3, 35, 11, 43, 1, 33, 9, 41,
	51, 19, 59, 27, 49, 17, 57, 25,
	15, 47, 7, 39, 13, 45, 5, 37,
	63, 31, 55, 23, 61, 29, 53, 21,
}

// ditherShader draws what a shader draws in black and white dots, by an
// ordered dither on the pixels of the image (Options.DitherImages): a
// pixel is left as it was where the shader's alpha is below the
// threshold of its place, black where its luminance is, white elsewhere.
// It is the same in every band, which an error diffusion would not be.
type ditherShader struct{ shader }

func (s ditherShader) shade(x, y, n int, dst []float32) {
	s.shader.shade(x, y, n, dst)
	row := bayer8[(y&7)*8:]
	for i := 0; i < n; i++ {
		px := dst[4*i : 4*i+4]
		t := (float32(row[(x+i)&7]) + 0.5) / 64
		a := px[3]
		if a <= 0 || a < t {
			px[0], px[1], px[2], px[3] = 0, 0, 0, 0
			continue
		}
		// the luminance of the colour, unpremultiplied
		l := (0.299*px[0] + 0.587*px[1] + 0.114*px[2]) / a
		if l < t {
			px[0], px[1], px[2] = 0, 0, 0
		} else {
			px[0], px[1], px[2] = 1, 1, 1
		}
		px[3] = 1
	}
}
