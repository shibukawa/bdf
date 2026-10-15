package imagebdf

import (
	stdcontext "context"
	"errors"
	"fmt"
	"image"
	"math"
	"runtime"
	"sync"

	"github.com/shibukawa/bdf"
)

// Pages draws selected zero-based page numbers and calls yield in selection
// order, on the calling goroutine. nil pages selects all pages; an empty slice
// draws none. The callback owns each image and may retain or modify it.
// Each worker uses a separate Renderer. Neither the receiver, document nor
// external font resources may be mutated concurrently with this call. FontFS
// must support concurrent reads when more than one worker is used.
//
// workers=0 chooses up to four workers; explicit values must be 1..16.
// Concurrency is further bounded by a 32-million-pixel aggregate base-canvas
// budget (a larger individual page is drawn alone). Layer/scratch allocations
// are additional. Output images in flight are bounded by the worker count;
// images retained by yield are the caller's responsibility.
// Cancellation and callback errors stop between drawings/batches, not during
// a rasterization already in progress. All workers finish before return.
func (r *Renderer) Pages(ctx stdcontext.Context, v *bdf.View, pages []int, scale float64, workers int, yield func(page int, img *image.RGBA) error) error {
	if ctx == nil || yield == nil {
		return errors.New("imagebdf: Pages requires a context and callback")
	}
	if workers < 0 || workers > 16 {
		return errors.New("imagebdf: workers must be 0..16")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if v == nil {
		return errors.New("imagebdf: no view")
	}
	if v.Kind == bdf.ViewSheet {
		return errors.New("imagebdf: Pages does not draw sheet views")
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
	maxPixels := 1
	// Validate the entire selection before invoking a callback.
	for _, page := range pages {
		p, err := drawablePage(v, page)
		if err != nil {
			return err
		}
		w, h := math.Max(math.Ceil(float64(p.W)*scale-1e-6), 1), math.Max(math.Ceil(float64(p.H)*scale-1e-6), 1)
		if !(w*h <= MaxPixels) {
			return ErrTooLarge
		}
		maxPixels = max(maxPixels, int(w*h))
	}
	if len(pages) == 0 {
		return nil
	}
	if workers == 0 {
		workers = min(4, runtime.GOMAXPROCS(0))
	}
	workers = min(workers, len(pages), max(1, (32<<20)/maxPixels))
	// Workers persist across calls. They share immutable text masks through a
	// locked LRU, while fonts, decoded geometry and scratch remain independent.
	if len(r.pageWorkers) != workers {
		opts := r.opts
		opts.ImageCacheBytes = -1
		r.pageWorkers = make([]*Renderer, workers)
		for i := range r.pageWorkers {
			child := New(r.doc, &opts)
			child.textMasks = r.textMasks
			child.freeLimit = maxFreeSurfaces / workers
			r.pageWorkers[i] = child
		}
	}
	pool := r.pageWorkers
	type result struct {
		img       *image.RGBA
		err       error
		key       imageKey
		cacheable bool
		hit       bool
	}
	for start := 0; start < len(pages); start += workers {
		if err := ctx.Err(); err != nil {
			return err
		}
		n := min(workers, len(pages)-start)
		results := make([]result, n)
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			page := pages[start+i]
			p := v.Pages[page]
			w, h := int(math.Max(math.Ceil(float64(p.W)*scale-1e-6), 1)), int(math.Max(math.Ceil(float64(p.H)*scale-1e-6), 1))
			key, ok := r.imageKey(1, p, bdf.Rect{W: p.W, H: p.H}, w, h, 0)
			results[i].key, results[i].cacheable = key, ok
			if img := r.storedImage(key, ok); img != nil {
				results[i].img = img
				results[i].hit = true
				continue
			}
			wg.Add(1)
			go func(i, page int) {
				defer wg.Done()
				if err := ctx.Err(); err != nil {
					results[i].err = err
					return
				}
				results[i].img, results[i].err = pool[i].Page(v, page, scale)
			}(i, page)
		}
		wg.Wait()
		// Merge warnings deterministically, even when a callback or draw fails.
		for i := 0; i < n; i++ {
			for _, warning := range pool[i].Warnings() {
				r.warnf("%s", warning)
			}
		}
		for i, res := range results {
			if err := ctx.Err(); err != nil {
				return err
			}
			if res.err != nil {
				return fmt.Errorf("imagebdf: page %d: %w", pages[start+i], res.err)
			}
			if !res.hit {
				r.keepImage(res.key, res.cacheable, res.img)
			}
			if err := yield(pages[start+i], res.img); err != nil {
				return err
			}
		}
	}
	return nil
}

// PagesDPI is Pages with dpi converted to device pixels per point.
func (r *Renderer) PagesDPI(ctx stdcontext.Context, v *bdf.View, pages []int, dpi float64, workers int, yield func(page int, img *image.RGBA) error) error {
	return r.Pages(ctx, v, pages, dpi/72, workers, yield)
}
