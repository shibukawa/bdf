package imagebdf

import (
	stdcontext "context"
	"errors"
	"fmt"
	"image"
	"math"
	"sync"
	"sync/atomic"

	"github.com/shibukawa/bdf"
)

// PageRows draws a whole page of a fixed or flow view (or a strip of a
// scroll view) at scale device pixels per unit in bands of rows, from the
// top, and calls yield for each: with the first row of the band and the
// image of its rows, whose bounds are those rows of the page (Min.Y is
// y0). The image is an *image.RGBA for a mono of 0, and the image of the
// Mono otherwise (see PageMono); it is the renderer's, reused for a
// band after it, so yield copies what it keeps. The pixels are those of
// Page and PageMono.
//
// A band takes a canvas of its rows (Options.BandRows of them, and the
// rows above and below that the shadows of the page reach from) instead
// of the whole page, and the image of its rows: an A4 page at 600 dpi
// takes tens of megabytes instead of hundreds. The instructions of the
// page are run once a band, skipping from the second band the draws that
// lie outside it. The bands after the first are drawn in parallel, as
// many at a time as the renderer has threads (up to eight, those of the
// machine) and within MaxBandCanvases of canvases, each by a worker of
// its own with the memory of the renderer's options; yield is called in
// order, on the calling goroutine. An error of yield ends the drawing
// with that error once the bands begun are drawn. The image is not kept
// in the completed-image cache.
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
// known before the page is drawn: the first band is drawn again with as
// much when its shadows reach past its margin, and the bands after it
// with as much from the start (and as tall again, so that the margin is
// not most of a band). The first band notes the extent of each draw for
// the bands after it, which skip the draws outside them (see
// bandRecord). The bands after the first are drawn by the band workers
// of the renderer in parallel (see bandWorkers), and given to yield in
// order.
func (r *Renderer) pageRows(p *bdf.Page, w, h int, mono Mono, yield func(y0 int, rows image.Image) error) error {
	bg := r.background()
	sx, sy := float64(w)/float64(p.W), float64(h)/float64(p.H)
	pg := &pageBands{p: p, w: w, h: h, mono: mono, bg: bg, pm: matrix{sx, 0, 0, sy, 0, 0}}
	whole := image.Rect(0, 0, w, h)
	// the background is the canvas when the page covers the whole image
	ax, ay := pg.pm.apply(0, 0)
	bx, by := pg.pm.apply(float64(p.W), float64(p.H))
	cover := r.rectMask(ax, ay, bx, by, whole)
	pg.filled, pg.borders = pageBackground(cover, whole, &bg)
	r.releaseMask(cover)
	return r.bandsOf(pg, yield)
}

// Place says where a page goes on an image, to print it on paper (see
// PlaceRows): the image is W × H pixels; the page is drawn at Scale
// pixels per unit, turned clockwise by Turns quarter turns (0 to 3), and
// the box it then takes has its top-left corner at (X, Y) pixels; then
// the image is mirrored left to right (MirrorX) or top to bottom
// (MirrorY), as the back of a sheet is for some printers that turn the
// sheet over (pwg-raster-document-sheet-back). The rest of the image is
// the background, as the page's is.
type Place struct {
	W, H             int
	Scale            float64
	Turns            int
	X, Y             float64
	MirrorX, MirrorY bool
}

// MaxPlacePixels is the largest image PlaceRows draws. It is larger than
// MaxPixels because only a band of it is drawn at a time: an A0 sheet at
// 300 dpi, or A3 at 600, fits.
const MaxPlacePixels = 256 << 20

// matrix returns the device transform of a page of w × h units.
func (pl Place) matrix(w, h float64) matrix {
	s := pl.Scale
	pw, ph := w*s, h*s
	var m matrix
	switch pl.Turns & 3 {
	case 1:
		m = matrix{0, s, -s, 0, pl.X + ph, pl.Y}
	case 2:
		m = matrix{-s, 0, 0, -s, pl.X + pw, pl.Y + ph}
	case 3:
		m = matrix{0, -s, s, 0, pl.X, pl.Y + pw}
	default:
		m = matrix{s, 0, 0, s, pl.X, pl.Y}
	}
	if pl.MirrorX {
		m = matrix{-1, 0, 0, 1, float64(pl.W), 0}.mul(m)
	}
	if pl.MirrorY {
		m = matrix{1, 0, 0, -1, 0, float64(pl.H)}.mul(m)
	}
	return m
}

// PlaceRows draws a page of a fixed or flow view (or a strip of a
// scroll view) placed on an image as pl says, in bands of rows from the
// top, as PageRows draws a page: for printing, where the image is the
// sheet of paper at the printer's resolution and the page is scaled to
// it, turned to its orientation, and centred or put at its margins. The
// shadows of the page turn with it. The pixels of a Place of the page's
// own size, at no turn and at the origin, are those of PageRows.
func (r *Renderer) PlaceRows(v *bdf.View, page int, pl Place, mono Mono, yield func(y0 int, rows image.Image) error) (err error) {
	defer recoveredErr(&err)
	if err := mono.check(); err != nil {
		return err
	}
	if yield == nil {
		return errors.New("imagebdf: PlaceRows requires a callback")
	}
	p, err := drawablePage(v, page)
	if err != nil {
		return err
	}
	if pl.W <= 0 || pl.H <= 0 || !(pl.Scale > 0) || math.IsInf(pl.Scale, 0) || math.IsNaN(pl.X) || math.IsNaN(pl.Y) ||
		math.IsInf(pl.X, 0) || math.IsInf(pl.Y, 0) {
		return errors.New("imagebdf: empty placement")
	}
	if pl.Turns < 0 || pl.Turns > 3 {
		return fmt.Errorf("imagebdf: %d quarter turns", pl.Turns)
	}
	if float64(pl.W)*float64(pl.H) > MaxPlacePixels {
		return ErrTooLarge
	}
	bg := r.background()
	// the image is the background where the page is not: a sheet of paper
	pm := pl.matrix(float64(p.W), float64(p.H))
	// the turn and mirror of the page, without its scale, which the
	// offsets of its shadows go by
	orient := matrix{pm[0] / pl.Scale, pm[1] / pl.Scale, pm[2] / pl.Scale, pm[3] / pl.Scale, 0, 0}
	pg := &pageBands{p: p, w: pl.W, h: pl.H, mono: mono, bg: bg, pm: pm, orient: &orient, filled: &bg}
	return r.bandsOf(pg, yield)
}

// bandsOf draws the bands of a page (see pageRows).
func (r *Renderer) bandsOf(pg *pageBands, yield func(y0 int, rows image.Image) error) error {
	w, h := pg.w, pg.h
	rows := r.bandRows()
	// the first band, drawn again with the margin its shadows need
	var rec *bandRecord
	var out image.Image
	y1 := min(rows, h)
	for {
		rec = &bandRecord{recording: true, margin: pg.margin}
		out = r.band(pg, 0, y1, rec, out)
		if rec.margin > pg.margin {
			pg.margin = rec.margin
			continue
		}
		rec.recording = false
		break
	}
	if bandHook != nil {
		bandHook(rec)
	}
	if err := yield(0, out); err != nil {
		return err
	}
	if y1 >= h {
		return nil
	}
	// the bands after it: at least twice as tall as the margin, so that
	// the margin drawn on either side is no more than the band
	rows = max(rows, 2*pg.margin)
	var bands []image.Rectangle
	for y0 := y1; y0 < h; y0 += rows {
		bands = append(bands, image.Rect(0, y0, w, min(y0+rows, h)))
	}
	workers := r.bandWorkers(len(bands), 16*w*min(h, rows+2*pg.margin))
	if workers <= 1 {
		for _, b := range bands {
			out = r.band(pg, b.Min.Y, b.Max.Y, rec, out)
			if bandHook != nil {
				bandHook(rec)
			}
			if err := yield(b.Min.Y, out); err != nil {
				return err
			}
		}
		return nil
	}
	return r.bandsInParallel(pg, bands, rec, workers, yield)
}

// pageBands is a page drawn in bands (see pageRows): what every band is
// drawn from.
type pageBands struct {
	p    *bdf.Page
	w, h int
	mono Mono
	bg   solid
	pm   matrix
	// filled is the background the canvas of a band begins as, when the
	// page covers it, and borders says that the fractional edges of the
	// page are written after (see pageBackground)
	filled  *solid
	borders bool
	// margin is the rows above and below a band that are drawn with it
	// for the shadows (see bandRecord.margin)
	margin int
	// orient is the turn and mirror of the page on the image (Place),
	// which the offsets of its shadows go by; nil when it is upright
	orient *matrix
}

// band draws the rows [y0, y1) of a page on the renderer, with the
// margin of the page above and below them for the shadows, noting the
// draws when rec is recording, and returns the image of the rows: out
// when it is large enough, else a new one.
func (r *Renderer) band(pg *pageBands, y0, y1 int, rec *bandRecord, out image.Image) image.Image {
	whole := image.Rect(0, 0, pg.w, pg.h)
	canvas := image.Rect(0, max(0, y0-pg.margin), pg.w, min(pg.h, y1+pg.margin))
	d := r.canvasAt(pg.pm, canvas, whole, pg.filled)
	d.band = rec
	d.orient = pg.orient
	clip := d.clipTo(pg.pm, 0, 0, float64(pg.p.W), float64(pg.p.H))
	switch {
	case pg.borders:
		d.ctx.target.setBorders(clip, pg.bg)
	case pg.filled == nil:
		d.ctx.target.fillBackground(clip, pg.bg)
	}
	for _, l := range pg.p.Layers {
		d.drawTop(l.Obj, pg.pm, clip)
	}
	s := d.ctx.target
	band := image.Rect(0, y0, pg.w, y1)
	switch pg.mono {
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
	return out
}

// MaxBandCanvases bounds the memory of the canvases of the bands of a
// page drawn at a time (Options.BandWorkers): fewer bands are drawn at a
// time when they are large. An A4 page at 600 dpi takes 10 MiB a band
// of DefaultBandRows.
const MaxBandCanvases = 192 << 20

// bandWorkers returns how many bands of a page are drawn at a time, for
// n bands whose canvases take canvasBytes each: Options.BandWorkers, or
// as many as the threads of the renderer (see splitRows), within
// MaxBandCanvases, and makes the renderers of the workers (see
// bandsInParallel). One means the bands are drawn one after the other by
// the renderer itself.
func (r *Renderer) bandWorkers(n, canvasBytes int) int {
	workers := r.opts.BandWorkers
	if workers <= 0 {
		workers = min(r.threads, max(1, MaxBandCanvases/max(canvasBytes, 1)))
	}
	workers = min(workers, n, maxThreads)
	if workers <= 1 || plain {
		return 1
	}
	if len(r.rowWorkers) != workers {
		// the workers draw their bands alone: the renderer's threads are
		// theirs; they share the text masks, and decode what they draw
		// for themselves
		opts := r.opts
		opts.ImageCacheBytes = -1
		r.rowWorkers = make([]*Renderer, workers)
		for i := range r.rowWorkers {
			child := New(r.doc, &opts)
			child.textMasks = r.textMasks
			child.threads = 1
			r.rowWorkers[i] = child
		}
	}
	return workers
}

// bandsInParallel draws the bands after the first of a page with the
// band workers of the renderer, each a band at a time, and gives them
// to yield in the order of the bands: at most as many bands are drawn
// ahead of yield as there are workers, and a worker draws its next band
// once yield is done with the image of the one before (which it reuses).
// An error of yield ends the drawing once the bands begun are drawn.
func (r *Renderer) bandsInParallel(pg *pageBands, bands []image.Rectangle, rec *bandRecord, workers int, yield func(y0 int, rows image.Image) error) error {
	type result struct {
		img  image.Image
		err  error
		done chan struct{}
	}
	slots := make([]chan result, len(bands))
	for i := range slots {
		slots[i] = make(chan result, 1)
	}
	jobs := make(chan int)
	stop, fed := make(chan struct{}), make(chan struct{})
	var dispatched atomic.Int32
	var wg sync.WaitGroup
	pool := r.rowWorkers[:workers]
	for _, wk := range pool {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out := wk.bandImage
			defer func() { wk.bandImage = out }()
			for i := range jobs {
				res := result{done: make(chan struct{})}
				func() {
					defer func() {
						if p := recover(); p != nil {
							res.err = fmt.Errorf("imagebdf: the drawing failed: %v", p)
						}
					}()
					b := bands[i]
					out = wk.band(pg, b.Min.Y, b.Max.Y, rec, out)
					res.img = out
				}()
				slots[i] <- res
				<-res.done
			}
		}()
	}
	// the feeder hands out the bands in order, until told to stop
	go func() {
		defer close(fed)
		defer close(jobs)
		for i := range bands {
			select {
			case jobs <- i:
				dispatched.Add(1)
			case <-stop:
				return
			}
		}
	}()
	var err error
	taken := 0
	for ; taken < len(bands); taken++ {
		res := <-slots[taken]
		if res.err != nil {
			err = res.err
		} else {
			if bandHook != nil {
				bandHook(rec)
			}
			err = yield(bands[taken].Min.Y, res.img)
		}
		close(res.done)
		if err != nil {
			taken++
			break
		}
	}
	close(stop)
	<-fed
	// the bands handed out are drawn still: take them, so that their
	// workers may finish
	for i := taken; i < int(dispatched.Load()); i++ {
		res := <-slots[i]
		close(res.done)
	}
	wg.Wait()
	for _, wk := range pool {
		for _, warning := range wk.Warnings() {
			r.warnf("%s", warning)
		}
	}
	return err
}
