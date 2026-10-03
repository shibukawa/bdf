package formula

import (
	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/mathlayout"
)

// Object returns the formula as an object of paths: the outlines of its
// glyphs, the rules of its fractions and radicals, and the lines of its
// enclosures. The origin of the object is the left end of the baseline and
// its bounding box is what the formula draws (Ink). The object refers to no
// font and no other part, so it is drawn as it is (imagebdf.Object,
// ebitenginebdf.Compile) or added to a document.
func (l *Layout) Object() *bdf.Object {
	o := bdf.NewObject()
	ink := l.Ink()
	o.SetBBox(ink.X, ink.Y, ink.W, ink.H)
	l.Draw(o, 0, 0)
	return o
}

// Draw draws the formula into o with the left end of its baseline at
// (x, y), as paths (see Object). It leaves the fill and stroke styles of o
// as the formula set them.
func (l *Layout) Draw(o *bdf.Object, x, y float64) {
	l.t.mu.Lock()
	defer l.t.mu.Unlock()
	d := drawer{o: o, glyphs: map[glyphKey]bdf.PathRef{}}
	var run []bdf.Glyph
	flush := func() {
		if len(run) > 0 {
			o.FillPathRun(bdf.NonZero, run)
			run = nil
		}
	}
	for _, it := range l.box.Items() {
		switch it.Kind {
		case mathlayout.ItemRule:
			flush()
			d.fill(it.Color)
			o.FillRect(float32(x+it.X), float32(y+it.Y), float32(it.W), float32(it.H))
		case mathlayout.ItemStroke:
			flush()
			d.stroke(it.Color, it.W)
			o.StrokePath(o.AddPath(translate(it.Path, float32(x+it.X), float32(y+it.Y))))
		default:
			ref, ok := d.glyph(it)
			if !ok {
				continue
			}
			if !d.hasFill || d.fillColor != it.Color {
				flush()
				d.fill(it.Color)
			}
			run = append(run, bdf.Glyph{Path: ref, X: float32(x + it.X), Y: float32(y + it.Y)})
		}
	}
	flush()
}

// glyphKey is a glyph at a size: one path of the object.
type glyphKey struct {
	face     *mathlayout.Face
	glyph    uint16
	size     float64
	sx, sy   float64
	hasGlyph bool
}

// drawer keeps the paths an object has of the glyphs drawn into it and the
// styles its instructions have set.
type drawer struct {
	o           *bdf.Object
	glyphs      map[glyphKey]bdf.PathRef
	empty       map[glyphKey]bool
	fillColor   bdf.Color
	hasFill     bool
	strokeColor bdf.Color
	strokeWidth float64
	hasStroke   bool
}

func (d *drawer) fill(c bdf.Color) {
	if !d.hasFill || c != d.fillColor {
		d.o.FillColor(c)
		d.fillColor, d.hasFill = c, true
	}
}

func (d *drawer) stroke(c bdf.Color, w float64) {
	if !d.hasStroke || c != d.strokeColor || w != d.strokeWidth {
		d.o.StrokeColor(c)
		d.o.Line(float32(w), 1, 1, 10)
		d.strokeColor, d.strokeWidth, d.hasStroke = c, w, true
	}
}

// glyph returns the path of the glyph of an item at its size, with its
// origin on the baseline; ok is false for what draws nothing (a space, a
// character no font has).
func (d *drawer) glyph(it mathlayout.Item) (ref bdf.PathRef, ok bool) {
	if it.Face == nil || it.Face.Loaded == nil || !it.HasGlyph {
		return 0, false
	}
	k := glyphKey{face: it.Face, glyph: it.Glyph, size: it.Size, sx: it.ScaleX, sy: it.ScaleY, hasGlyph: true}
	if ref, ok := d.glyphs[k]; ok {
		return ref, true
	}
	if d.empty[k] {
		return 0, false
	}
	// font units (y up) to the units of the formula (y down)
	s := it.Size / it.Face.Loaded.UnitsPerEm()
	pen := &pathPen{sx: s * it.ScaleX, sy: -s * it.ScaleY}
	if !it.Face.Loaded.Outline(it.Glyph, pen) || len(pen.p.Verbs) == 0 {
		if d.empty == nil {
			d.empty = map[glyphKey]bool{}
		}
		d.empty[k] = true
		return 0, false
	}
	ref = d.o.AddPath(&pen.p)
	d.glyphs[k] = ref
	return ref, true
}

// pathPen makes a path of the outline of a glyph, scaled.
type pathPen struct {
	p      bdf.Path
	sx, sy float64
}

func (pp *pathPen) pt(x, y float64) (float32, float32) {
	return float32(x * pp.sx), float32(y * pp.sy)
}

func (pp *pathPen) MoveTo(x, y float64) {
	px, py := pp.pt(x, y)
	pp.p.MoveTo(px, py)
}

func (pp *pathPen) LineTo(x, y float64) {
	px, py := pp.pt(x, y)
	pp.p.LineTo(px, py)
}

func (pp *pathPen) QuadTo(cx, cy, x, y float64) {
	qx, qy := pp.pt(cx, cy)
	px, py := pp.pt(x, y)
	pp.p.QuadTo(qx, qy, px, py)
}

func (pp *pathPen) CubeTo(ax, ay, bx, by, x, y float64) {
	c1x, c1y := pp.pt(ax, ay)
	c2x, c2y := pp.pt(bx, by)
	px, py := pp.pt(x, y)
	pp.p.CubicTo(c1x, c1y, c2x, c2y, px, py)
}

func (pp *pathPen) Close() { pp.p.Close() }

// translate returns p moved by (dx, dy).
func translate(p *bdf.Path, dx, dy float32) *bdf.Path {
	out := &bdf.Path{Verbs: p.Verbs, Args: append([]float32(nil), p.Args...)}
	a := out.Args
	i := 0
	for _, v := range p.Verbs {
		// the arguments that are points, then the others (sizes, radii, angles)
		points, rest := 0, 0
		switch v {
		case bdf.VerbMove, bdf.VerbLine:
			points = 1
		case bdf.VerbQuad:
			points = 2
		case bdf.VerbCubic:
			points = 3
		case bdf.VerbRect:
			points, rest = 1, 2
		case bdf.VerbEllipse:
			points, rest = 1, 6
		case bdf.VerbArcTo:
			points, rest = 2, 1
		case bdf.VerbRoundRect:
			points, rest = 1, 3
		}
		for k := 0; k < points && i+1 < len(a); k++ {
			a[i] += dx
			a[i+1] += dy
			i += 2
		}
		i += rest
	}
	return out
}

// Document returns a document of one page that holds the formula with a
// margin around it: what a browser draws on a canvas with @bdfkit/render
// or shows with the viewer, and what raster/imagebdf draws as a page. The
// page carries the formula in linear notation (Text) as its text, for
// search, copy and screen readers.
func (l *Layout) Document(margin float64) *bdf.Document {
	doc := bdf.NewDocument()
	ink := l.Ink()
	h, bbox := doc.AddObject(l.Object())
	page := bdf.NewObject()
	page.Mark(bdf.MarkAltText, l.text)
	page.UseAt(page.AddObject(h, bbox), float32(margin)-ink.X, float32(margin)-ink.Y)
	ph, _ := doc.AddObject(page)
	v := doc.NewView("formula", bdf.ViewFixed, "")
	v.AddPage(ink.W+2*float32(margin), ink.H+2*float32(margin), bdf.Layer{Role: bdf.RoleBody, Obj: ph})
	return doc
}
