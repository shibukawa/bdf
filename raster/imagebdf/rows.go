package imagebdf

import (
	stdcontext "context"
	"errors"
	"fmt"
	"image"

	"github.com/shibukawa/bdf"
)

// PageRows draws a whole page of a fixed or flow view (or a strip of a
// scroll view) at scale device pixels per unit in bands of rows, from the
// top, and calls yield for each: with the first row of the band and the
// image of its rows, whose bounds are those rows of the page (Min.Y is
// y0). The image is an *image.RGBA for a mono of 0, and the image of the
// Mono otherwise (see PageMono); it is the renderer's, reused for the
// next band, so yield copies what it keeps. The pixels are those of Page
// and PageMono.
//
// A band takes a canvas of its rows (Options.BandRows of them, and the
// rows above and below that the shadows of the page reach from) instead
// of the whole page, and the image of its rows: an A4 page at 600 dpi
// takes tens of megabytes instead of hundreds. The instructions of the
// page are run once a band, skipping from the second band the draws that
// lie outside it. An error of yield ends the drawing with that error.
// The image is not kept in the completed-image cache.
func (r *Renderer) PageRows(v *bdf.View, page int, scale float64, mono Mono, yield func(y0 int, rows image.Image) error) (err error) {
	defer recoveredErr(&err)
	if err := mono.check(); err != nil {
		return err
	}
	if yield == nil {
		return errors.New("imagebdf: PageRows requires a callback")
	}
	p, err := drawablePage(v, page)
	if err != nil {
		return err
	}
	w, h, err := pagePixels(p, scale)
	if err != nil {
		return err
	}
	return r.pageRows(p, w, h, mono, yield)
}

// PagesRows is PageRows for the selected zero-based page numbers, in
// selection order, one after the other: yield takes the page number too.
// nil pages selects all pages; an empty slice draws none. The selection
// is checked before any page is drawn. Cancellation of ctx is seen
// between bands.
func (r *Renderer) PagesRows(ctx stdcontext.Context, v *bdf.View, pages []int, scale float64, mono Mono, yield func(page, y0 int, rows image.Image) error) (err error) {
	defer recoveredErr(&err)
	if ctx == nil || yield == nil {
		return errors.New("imagebdf: PagesRows requires a context and callback")
	}
	if err := mono.check(); err != nil {
		return err
	}
	if v == nil {
		return errors.New("imagebdf: no view")
	}
	if v.Kind == bdf.ViewSheet {
		return errors.New("imagebdf: PagesRows does not draw sheet views")
	}
	if !(scale > 0) {
		return errors.New("imagebdf: empty page")
	}
	if pages == nil {
		pages = make([]int, len(v.Pages))
		for i := range pages {
			pages[i] = i
		}
	}
	for _, page := range pages {
		if _, _, err := PagePixels(v, page, scale); err != nil {
			return err
		}
	}
	for _, page := range pages {
		if err := ctx.Err(); err != nil {
			return err
		}
		p := v.Pages[page]
		w, h, _ := pagePixels(p, scale)
		err := r.pageRows(p, w, h, mono, func(y0 int, rows image.Image) error {
			if err := ctx.Err(); err != nil {
				return err
			}
			return yield(page, y0, rows)
		})
		if err != nil {
			return fmt.Errorf("imagebdf: page %d: %w", page, err)
		}
	}
	return nil
}

// recoveredErr makes an error of a panic of a drawing (see recovered).
func recoveredErr(err *error) {
	if p := recover(); p != nil {
		*err = fmt.Errorf("imagebdf: the drawing failed: %v", p)
	}
}

// bandHook is called with the record of each band drawn, for the tests.
var bandHook func(*bandRecord)

// bandRows is the rows of a band of a renderer (Options.BandRows).
func (r *Renderer) bandRows() int {
	if r.opts.BandRows > 0 {
		return r.opts.BandRows
	}
	return DefaultBandRows
}

// pageRows draws a page of w × h pixels in bands of rows (see PageRows).
//
// Each band is drawn as drawPage draws the page, on a canvas of the rows
// of the band and a margin above and below it for the shadows, at those
// rows of the device: the transforms are those of the whole page, so
// the pixels are the same. What the shadows need of the margin is not
// known before the page is drawn: a band whose shadows reach past its
// margin is drawn again with as much, and the bands after it with as
// much from the start. The first band notes the extent of each draw for
// the bands after it, which skip the draws outside them (see
// bandRecord).
func (r *Renderer) pageRows(p *bdf.Page, w, h int, mono Mono, yield func(y0 int, rows image.Image) error) error {
	bg := r.background()
	sx, sy := float64(w)/float64(p.W), float64(h)/float64(p.H)
	pm := matrix{sx, 0, 0, sy, 0, 0}
	whole := image.Rect(0, 0, w, h)
	// the background is the canvas when the page covers the whole image
	ax, ay := pm.apply(0, 0)
	bx, by := pm.apply(float64(p.W), float64(p.H))
	cover := r.rectMask(ax, ay, bx, by, whole)
	var filled *solid
	if cover.a == nil && cover.r == whole {
		filled = &bg
	}
	r.releaseMask(cover)
	rows := r.bandRows()
	var rec *bandRecord
	margin := 0
	var out image.Image
	for y0 := 0; y0 < h; y0 += rows {
		y1 := min(y0+rows, h)
		for {
			canvas := image.Rect(0, max(0, y0-margin), w, min(h, y1+margin))
			d := r.canvasAt(pm, canvas, whole, filled)
			if rec == nil {
				rec = &bandRecord{recording: true, margin: margin}
			}
			rec.visit = 0
			d.band = rec
			clip := d.clipTo(pm, 0, 0, float64(p.W), float64(p.H))
			if filled == nil {
				d.ctx.target.fillBackground(clip, bg)
			}
			for _, l := range p.Layers {
				d.drawTop(l.Obj, pm, clip)
			}
			s := d.ctx.target
			if rec.margin > margin {
				// the shadows reach past the margin: once more with as
				// much, noting the extents anew
				margin = rec.margin
				rec = nil
				r.release(s)
				continue
			}
			rec.recording = false
			if bandHook != nil {
				bandHook(rec)
			}
			band := image.Rect(0, y0, w, y1)
			switch mono {
			case Gray:
				img, _ := out.(*image.Gray)
				out = s.toGrayInto(img, band)
			case Bilevel:
				img, _ := out.(*image.Paletted)
				out = s.toBilevelInto(img, band)
			default:
				img, _ := out.(*image.RGBA)
				out = s.toRGBAInto(img, band)
			}
			r.release(s)
			break
		}
		if err := yield(y0, out); err != nil {
			return err
		}
	}
	return nil
}
