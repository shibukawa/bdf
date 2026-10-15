package imagebdf

import (
	"errors"
	"fmt"
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
	doc := bdf.NewDocument()
	hash, _ := doc.AddObject(o)
	var ropts Options
	if opts != nil {
		ropts = *opts
	}
	// A one-shot helper has no later request to reuse a completed image.
	ropts.ImageCacheBytes = -1
	return New(doc, &ropts).Object(hash, scale)
}

// Object draws a document object, including its dependencies, at scale pixels
// per unit. Its bounds use the object's origin as the package-level Object
// helper does. A nil background is transparent. Reusing a Renderer also reuses
// the completed-image and text-mask caches.
func (r *Renderer) Object(hash bdf.Hash, scale float64) (img *image.RGBA, err error) {
	defer recovered(&img, &err)
	o := r.object(hash)
	if o == nil {
		return nil, fmt.Errorf("imagebdf: object %s is missing or invalid", hash)
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
	key, cacheable := r.imageKey(4, hash, b, int(fw), int(fh), scale)
	if cached := r.storedImage(key, cacheable); cached != nil {
		return cached, nil
	}
	m := matrix{scale, 0, 0, scale, -x0, -y0}
	d := r.canvasAt(m, int(fw), int(fh))
	clip := &mask{r: d.ctx.target.bounds()}
	if r.opts.Background != nil {
		d.ctx.target.fill(clip, r.background(), 1, bdf.BlendSourceOver, clip)
	}
	d.drawTop(hash, m, clip)
	img = d.finish()
	img.Rect = img.Rect.Add(image.Pt(int(x0), int(y0)))
	r.keepImage(key, cacheable, img)
	return img, nil
}
