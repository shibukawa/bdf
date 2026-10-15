package imagebdf

import (
	"errors"
	"image"
	"math"

	"github.com/shibukawa/bdf"
)

// ErrNeedsDocument is returned by Object for an object that refers to other
// parts (an embedded font, an image, a path collection or another object):
// add it to a document with them and draw a page.
var ErrNeedsDocument = errors.New("imagebdf: the object refers to other parts of a document")

// Object draws an object that needs no other part, such as a formula laid
// out as paths, at scale device pixels per unit. The image holds the
// bounding box of the object grown to whole pixels, and its bounds are in
// device pixels from the origin of the object: the pixel at (0, 0) is where
// the origin is drawn, wherever the bounding box lies around it. The
// background is opts.Background, and transparent without one.
func Object(o *bdf.Object, scale float64, opts *Options) (img *image.RGBA, err error) {
	defer recovered(&img, &err)
	if o == nil {
		return nil, errors.New("imagebdf: no object")
	}
	if len(o.Deps()) > 0 {
		return nil, ErrNeedsDocument
	}
	b := o.BBox
	if !(b.W > 0) || !(b.H > 0) || !(scale > 0) {
		return nil, errors.New("imagebdf: empty object")
	}
	x0, y0 := math.Floor(float64(b.X)*scale+1e-6), math.Floor(float64(b.Y)*scale+1e-6)
	fw := math.Max(math.Ceil(float64(b.X+b.W)*scale-1e-6)-x0, 1)
	fh := math.Max(math.Ceil(float64(b.Y+b.H)*scale-1e-6)-y0, 1)
	if !(fw*fh <= MaxPixels) || !(math.Abs(x0) < math.MaxInt32) || !(math.Abs(y0) < math.MaxInt32) {
		return nil, ErrTooLarge
	}
	doc := bdf.NewDocument()
	h, _ := doc.AddObject(o)
	r := New(doc, opts)
	m := matrix{scale, 0, 0, scale, -x0, -y0}
	d := r.canvasAt(m, int(fw), int(fh))
	clip := &mask{r: d.ctx.target.bounds()}
	if r.opts.Background != nil {
		d.ctx.target.fill(clip, r.background(), 1, bdf.BlendSourceOver, clip)
	}
	d.drawTop(h, m, clip)
	img = d.finish()
	img.Rect = img.Rect.Add(image.Pt(int(x0), int(y0)))
	return img, nil
}
