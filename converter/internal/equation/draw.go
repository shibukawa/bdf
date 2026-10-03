package equation

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/internal/mathlayout"
)

// Place draws a laid out formula into a child object of cv whose origin is
// the left end of the formula's baseline, and uses it at (x, y) after an
// ALT_TEXT of text (the formula in linear notation), so that extraction
// reads the formula as that text spread over its width.
func Place(cv *canvas.Canvas, b *Box, x, y float64, text string) {
	ink := b.Ink()
	m := float32(math.Max(b.H+b.D, 1) * 0.1)
	bbox := bdf.Rect{X: ink.X - m, Y: ink.Y - m, W: ink.W + 2*m, H: ink.H + 2*m}
	ch, ref := cv.Child(bbox)
	ch.Obj.SetBBox(bbox.X, bbox.Y, bbox.W, bbox.H)
	// the object starts in the state of the place that uses it
	ch.Obj.TextStyle(bdf.AlignLeft, bdf.BaselineAlphabetic, bdf.DirInherit, 0)
	Draw(ch, b, 0, 0)
	cv.Obj.Mark(bdf.MarkAltText, text)
	cv.Obj.UseAt(ref, f32(x), f32(y))
	cv.Drawn = true
}

// Draw draws b into cv with its origin (the left end of its baseline) at
// (x, y).
func Draw(cv *canvas.Canvas, b *Box, x, y float64) {
	d := drawer{cv: cv}
	items := b.Items()
	for i := 0; i < len(items); {
		it := &items[i]
		switch it.Kind {
		case mathlayout.ItemRule:
			d.fill(it.Color)
			cv.Obj.FillRect(f32(x+it.X), f32(y+it.Y), f32(it.W), f32(it.H))
			i++
		case mathlayout.ItemStroke:
			d.strokeStyle(it.Color, it.W)
			cv.Obj.Save()
			cv.Obj.Translate(f32(x+it.X), f32(y+it.Y))
			cv.Obj.StrokePath(cv.Obj.AddPath(it.Path))
			cv.Obj.Restore()
			i++
		default:
			if it.Face == nil {
				i++
				continue
			}
			d.font(choice(it.Face), it.Size)
			d.fill(it.Color)
			if it.ScaleX != 1 || it.ScaleY != 1 {
				cv.Obj.Save()
				cv.Obj.Transform(f32(it.ScaleX), 0, 0, f32(it.ScaleY), f32(x+it.X), f32(y+it.Y))
				cv.Obj.FillText(it.Text, 0, 0, f32(it.Advance/it.ScaleX))
				cv.Obj.Restore()
				i++
				continue
			}
			// glyphs that follow each other on a baseline in one font are
			// one run
			adv := it.Advance
			j := i + 1
			for j < len(items) {
				n := &items[j]
				p := &items[j-1]
				if n.Kind != mathlayout.ItemGlyph || n.Face != it.Face || n.Size != it.Size || n.Color != it.Color || n.Y != it.Y ||
					n.ScaleX != 1 || n.ScaleY != 1 || math.Abs(p.X+p.Advance-n.X) > 1e-3 {
					break
				}
				adv += n.Advance
				j++
			}
			text := it.Text
			if j > i+1 {
				var b strings.Builder
				for k := i; k < j; k++ {
					b.WriteString(items[k].Text)
				}
				text = b.String()
			}
			cv.Obj.FillText(text, f32(x+it.X), f32(y+it.Y), f32(adv))
			i = j
		}
	}
	cv.Drawn = cv.Drawn || len(items) > 0
}

// drawer keeps the state an object's instructions have set.
type drawer struct {
	cv          *canvas.Canvas
	fontRef     bdf.FontRef
	fontSize    float64
	hasFont     bool
	fillColor   bdf.Color
	hasFill     bool
	strokeColor bdf.Color
	strokeWidth float64
	hasStroke   bool
}

func (d *drawer) font(fc *fontset.Choice, size float64) {
	ref := d.cv.Font(fc.Use)
	if !d.hasFont || ref != d.fontRef || size != d.fontSize {
		d.cv.Obj.Font(ref, f32(size))
		d.fontRef, d.fontSize, d.hasFont = ref, size, true
	}
}

func (d *drawer) fill(c bdf.Color) {
	if !d.hasFill || c != d.fillColor {
		d.cv.Obj.FillColor(c)
		d.fillColor, d.hasFill = c, true
	}
}

func (d *drawer) strokeStyle(c bdf.Color, w float64) {
	if !d.hasStroke || c != d.strokeColor || w != d.strokeWidth {
		d.cv.Obj.StrokeColor(c)
		d.cv.Obj.Line(f32(w), 1, 1, 10)
		d.strokeColor, d.strokeWidth, d.hasStroke = c, w, true
	}
}

func f32(v float64) float32 { return float32(v) }
