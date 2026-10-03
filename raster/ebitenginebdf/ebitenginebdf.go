// Package ebitenginebdf draws BDF objects of paths with Ebitengine: the
// formulas that package formula lays out, and other objects that fill and
// stroke paths in plain colors.
//
//	ts, _ := formula.New(nil)
//	l := ts.Layout(formula.ParseTeX(`e^{i\pi} + 1 = 0`), formula.Style{Size: 48})
//	drawing, _ := ebitenginebdf.Compile(l.Object())
//
//	func (g *Game) Draw(screen *ebiten.Image) {
//		var opts ebitenginebdf.DrawOptions
//		opts.GeoM.Translate(40, 120) // where the origin of the object goes
//		opts.AntiAlias = true
//		drawing.Draw(screen, &opts)
//	}
//
// An object is compiled once into paths of Ebitengine's vector package and
// drawn every frame, at any position, scale and rotation: it is vector
// graphics, sharp at every size. What an object draws besides paths in
// plain colors (text in fonts, images, gradients, clips, groups, masks)
// is not drawn by this package; the pages of documents, which hold all of
// those, are drawn into an image by raster/imagebdf and shown as an
// ebiten.Image (see Image).
package ebitenginebdf

import (
	"errors"
	"fmt"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/raster/internal/shapes"
)

// ErrNeedsDocument is returned by Compile for an object that refers to
// other parts (an embedded font, an image, a path collection or another
// object).
var ErrNeedsDocument = errors.New("ebitenginebdf: the object refers to other parts of a document")

// UnsupportedError is the error of Compile for an object that draws more
// than paths in plain colors.
type UnsupportedError struct {
	Op byte // the first instruction that is not drawn (bdf.OpFillText …)
}

func (e *UnsupportedError) Error() string {
	return "ebitenginebdf: the object draws more than paths in plain colors: " + bdf.OpName(e.Op)
}

// Drawing is an object compiled for drawing. It is drawn from the goroutine
// that draws the game, like the images of Ebitengine.
type Drawing struct {
	shapes []shape
	bounds bdf.Rect

	// the paths under the transform of the last drawing: a drawing that
	// stays where it is is not transformed again every frame
	geoM   ebiten.GeoM
	placed []vector.Path
}

// shape is a path that is filled or stroked in one color.
type shape struct {
	path       vector.Path
	stroke     bool
	fill       vector.FillOptions
	line       vector.StrokeOptions
	r, g, b, a float32 // premultiplied
}

// Compile compiles an object that needs no other part. An object that
// draws more than paths in plain colors is an error (see the package
// comment).
func Compile(o *bdf.Object) (*Drawing, error) {
	if o == nil {
		return nil, errors.New("ebitenginebdf: no object")
	}
	if len(o.Deps()) > 0 {
		return nil, ErrNeedsDocument
	}
	part, err := bdf.DecodeObject(o.Encode())
	if err != nil {
		return nil, fmt.Errorf("ebitenginebdf: %w", err)
	}
	return CompilePart(part)
}

// CompilePart is Compile for an object read from a document
// (bdf.Reader.Object). The paths it refers to in path collections are not
// drawn.
func CompilePart(o *bdf.ObjectPart) (*Drawing, error) {
	if o == nil {
		return nil, errors.New("ebitenginebdf: no object")
	}
	list, err := shapes.Of(o)
	if ue := (*shapes.UnsupportedError)(nil); errors.As(err, &ue) {
		return nil, &UnsupportedError{Op: ue.Op}
	}
	if err != nil {
		return nil, fmt.Errorf("ebitenginebdf: %w", err)
	}
	d := &Drawing{bounds: o.BBox, shapes: make([]shape, len(list))}
	for i := range list {
		s, sh := &list[i], &d.shapes[i]
		addPath(&sh.path, &s.Path)
		c := uint32(s.Color)
		a := float32(c&0xff) / 255
		sh.r, sh.g, sh.b, sh.a = float32(c>>24)/255*a, float32(c>>16&0xff)/255*a, float32(c>>8&0xff)/255*a, a
		if !s.Stroke {
			if s.Rule == bdf.EvenOdd {
				sh.fill.FillRule = vector.FillRuleEvenOdd
			}
			continue
		}
		sh.stroke = true
		sh.line = vector.StrokeOptions{Width: float32(s.Width), MiterLimit: float32(s.Miter)}
		switch s.Cap {
		case bdf.CapRound:
			sh.line.LineCap = vector.LineCapRound
		case bdf.CapSquare:
			sh.line.LineCap = vector.LineCapSquare
		}
		switch s.Join {
		case bdf.JoinRound:
			sh.line.LineJoin = vector.LineJoinRound
		case bdf.JoinBevel:
			sh.line.LineJoin = vector.LineJoinBevel
		}
	}
	return d, nil
}

// addPath adds a path of moves, lines, curves and closes to dst.
func addPath(dst *vector.Path, p *bdf.Path) {
	a := p.Args
	for _, v := range p.Verbs {
		switch v {
		case bdf.VerbMove:
			dst.MoveTo(a[0], a[1])
			a = a[2:]
		case bdf.VerbLine:
			dst.LineTo(a[0], a[1])
			a = a[2:]
		case bdf.VerbQuad:
			dst.QuadTo(a[0], a[1], a[2], a[3])
			a = a[4:]
		case bdf.VerbCubic:
			dst.CubicTo(a[0], a[1], a[2], a[3], a[4], a[5])
			a = a[6:]
		case bdf.VerbClose:
			dst.Close()
		}
	}
}

// Bounds returns the bounding box of the object, around its origin.
func (d *Drawing) Bounds() bdf.Rect { return d.bounds }

// DrawOptions are how a drawing is drawn.
type DrawOptions struct {
	// GeoM maps the units of the object to the pixels of the destination:
	// where the origin of the object goes, and the scale and rotation.
	GeoM ebiten.GeoM
	// ColorScale scales the colors of the object: its alpha fades the
	// drawing, and a color tints an object drawn in white.
	ColorScale ebiten.ColorScale
	// AntiAlias smooths the edges, as in vector.DrawPathOptions.
	AntiAlias bool
	Blend     ebiten.Blend
}

// Draw draws the drawing on dst. A nil opts draws it with its origin at
// the top left corner of dst, without anti-aliasing.
func (d *Drawing) Draw(dst *ebiten.Image, opts *DrawOptions) {
	var o DrawOptions
	if opts != nil {
		o = *opts
	}
	paths := d.place(o.GeoM)
	// a line is as wide as the transform makes it (the mean scale where
	// it does not scale evenly)
	k := float32(math.Sqrt(math.Abs(o.GeoM.Element(0, 0)*o.GeoM.Element(1, 1) - o.GeoM.Element(0, 1)*o.GeoM.Element(1, 0))))
	for i := range d.shapes {
		s := &d.shapes[i]
		dp := vector.DrawPathOptions{AntiAlias: o.AntiAlias, Blend: o.Blend}
		dp.ColorScale.Scale(s.r, s.g, s.b, s.a)
		dp.ColorScale.ScaleWithColorScale(o.ColorScale)
		if s.stroke {
			line := s.line
			line.Width *= k
			vector.StrokePath(dst, &paths[i], &line, &dp)
		} else {
			vector.FillPath(dst, &paths[i], &s.fill, &dp)
		}
	}
}

// place returns the paths of the shapes under a transform.
func (d *Drawing) place(g ebiten.GeoM) []vector.Path {
	if d.placed != nil && d.geoM == g {
		return d.placed
	}
	if d.placed == nil {
		d.placed = make([]vector.Path, len(d.shapes))
	}
	for i := range d.shapes {
		d.placed[i].Reset()
		d.placed[i].AddPath(&d.shapes[i].path, &vector.AddPathOptions{GeoM: g})
	}
	d.geoM = g
	return d.placed
}

// Image returns an image drawn by raster/imagebdf (a page or a region of a
// document, or an object) as an image of Ebitengine, with where its top
// left corner goes when the origin of what it shows is at (0, 0): the
// images of imagebdf.Object keep their origin in their bounds.
func Image(img *image.RGBA) (*ebiten.Image, image.Point) {
	return ebiten.NewImageFromImage(img), img.Bounds().Min
}
