// Package flatten draws as pictures what a vector format without
// transparency cannot draw (PostScript, vector/psbdf): the rectangles of a
// page where something is drawn with alpha, a blend mode, a group or a
// soft mask, a shadow or a filter are collected, merged until none
// overlap, and drawn by the image renderer with all that the page draws
// there. The rest of the page stays vectors, under a clip that leaves the
// rectangles out.
package flatten

import (
	"image"
	"math"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

// Limits past which a page is one picture: the rectangles found, the
// rectangles left after merging, and the part of the page they cover.
const (
	maxBoxes   = 4096
	maxRegions = 64
	maxCover   = 0.5
)

// Regions collects the rectangles of a page that are drawn as pictures.
type Regions struct {
	boxes [][4]float64
	whole bool
}

// Add adds a box (x0, y0, x1, y1) in the units of the page; boxes that are
// empty or not numbers are left out.
func (r *Regions) Add(b [4]float64) {
	for _, v := range b {
		if math.IsNaN(v) {
			return
		}
	}
	if !(b[2] > b[0]) || !(b[3] > b[1]) {
		return
	}
	if len(r.boxes) >= maxBoxes {
		r.whole = true
		return
	}
	r.boxes = append(r.boxes, b)
}

// Empty reports whether no box was added.
func (r *Regions) Empty() bool { return len(r.boxes) == 0 && !r.whole }

// Merge returns the rectangles of a page of w × h: the boxes grown to whole
// units and by margin (for the edges a renderer smooths), cut to the page,
// and merged with those they overlap or come within margin of, until none
// do. Too many of them, or too much of the page, is the whole page.
func (r *Regions) Merge(w, h, margin float64) [][4]float64 {
	if r.Empty() {
		return nil
	}
	page := [4]float64{0, 0, w, h}
	if r.whole {
		return [][4]float64{page}
	}
	var out [][4]float64
	for _, b := range r.boxes {
		b = [4]float64{math.Floor(b[0] - margin), math.Floor(b[1] - margin), math.Ceil(b[2] + margin), math.Ceil(b[3] + margin)}
		b = [4]float64{math.Max(b[0], 0), math.Max(b[1], 0), math.Min(b[2], w), math.Min(b[3], h)}
		if b[2] > b[0] && b[3] > b[1] {
			out = append(out, b)
		}
	}
	near := func(a, b [4]float64) bool {
		return a[0] <= b[2]+margin && b[0] <= a[2]+margin && a[1] <= b[3]+margin && b[1] <= a[3]+margin
	}
	for merged := true; merged; {
		merged = false
		for i := 0; i < len(out); i++ {
			for j := i + 1; j < len(out); j++ {
				if near(out[i], out[j]) {
					a, b := out[i], out[j]
					out[i] = [4]float64{math.Min(a[0], b[0]), math.Min(a[1], b[1]), math.Max(a[2], b[2]), math.Max(a[3], b[3])}
					out = append(out[:j], out[j+1:]...)
					merged = true
					j = i
				}
			}
		}
	}
	area := 0.0
	for _, b := range out {
		area += (b[2] - b[0]) * (b[3] - b[1])
	}
	if len(out) > maxRegions || area > maxCover*w*h {
		return [][4]float64{page}
	}
	return out
}

// Picture draws a rectangle (x0, y0, x1, y1) of a page of a view at a
// resolution (dots per unit of 1/72 inch × 72).
func Picture(rd *imagebdf.Renderer, v *bdf.View, page int, box [4]float64, dpi float64) (*image.RGBA, error) {
	k := dpi / 72
	w := max(int(math.Ceil((box[2]-box[0])*k)), 1)
	h := max(int(math.Ceil((box[3]-box[1])*k)), 1)
	return rd.Region(v, page, bdf.Rect{X: float32(box[0]), Y: float32(box[1]), W: float32(box[2] - box[0]), H: float32(box[3] - box[1])}, w, h)
}
