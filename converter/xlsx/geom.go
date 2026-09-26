package xlsx

import (
	"math"
	"sort"
	"strings"

	"github.com/shibukawa/bdf"
)

// Excel sizes rows in points and columns in characters of the default
// font's widest digit (the maximum digit width, MDW, in pixels at 96 dpi),
// and draws both on whole pixels; BDF units are points (0.75 px).

const pxPt = 0.75 // points per pixel at 96 dpi

// axis is the rows or the columns of a sheet: a default size with
// exceptions, positions by prefix sums.
type axis struct {
	n    int // entries in the view
	def  float64
	idx  []int     // indexes with their own size, ascending
	size []float64 // their sizes
	pre  []float64 // pre[i] = sum over idx[:i] of (size - def)
}

func newAxis(n int, def float64, sizes map[int]float64) *axis {
	a := &axis{n: n, def: def}
	for i := range sizes {
		if i < n && sizes[i] != def {
			a.idx = append(a.idx, i)
		}
	}
	sort.Ints(a.idx)
	a.size = make([]float64, len(a.idx))
	a.pre = make([]float64, len(a.idx)+1)
	for k, i := range a.idx {
		a.size[k] = sizes[i]
		a.pre[k+1] = a.pre[k] + sizes[i] - def
	}
	return a
}

// at returns the size of entry i.
func (a *axis) at(i int) float64 {
	k := sort.SearchInts(a.idx, i)
	if k < len(a.idx) && a.idx[k] == i {
		return a.size[k]
	}
	return a.def
}

// pos returns the start of entry i (i may be n: the end).
func (a *axis) pos(i int) float64 {
	k := sort.SearchInts(a.idx, i) // exceptions before i
	return float64(i)*a.def + a.pre[k]
}

// total is the length of the axis.
func (a *axis) total() float64 { return a.pos(a.n) }

// index returns the entry that contains position p (clamped to 0..n-1).
func (a *axis) index(p float64) int {
	i := sort.Search(a.n, func(i int) bool { return a.pos(i+1) > p })
	return min(i, a.n-1)
}

// runs returns the sizes as the manifest's run-length list.
func (a *axis) runs() []bdf.Run {
	var out []bdf.Run
	add := func(count int, size float64) {
		if count <= 0 {
			return
		}
		s := float32(size)
		if n := len(out); n > 0 && out[n-1][1] == s {
			out[n-1][0] += float32(count)
			return
		}
		out = append(out, bdf.Run{float32(count), s})
	}
	prev := 0
	for k, i := range a.idx {
		add(i-prev, a.def)
		add(1, a.size[k])
		prev = i + 1
	}
	add(a.n-prev, a.def)
	return out
}

// maxDigitWidth measures the widest digit of the default font in pixels
// the way Windows draws it at 96 dpi: at a whole pixel size, rounded.
// Calibri's hinting draws its digits a pixel narrower than that, and it is
// the default font of most workbooks, so its advance is taken unrounded.
func (c *converter) maxDigitWidth(f *xfont) float64 {
	name := c.fontName(f)
	fc := c.fonts.Choose(name, false, false, false)
	w := 0.0
	for r := '0'; r <= '9'; r++ {
		w = math.Max(w, c.fonts.Advance(fc, r))
	}
	if fc.Loaded == nil || w == 0 {
		return 7
	}
	ppem := math.Round(f.size * 96 / 72)
	if strings.EqualFold(name, "Calibri") || strings.EqualFold(name, "Carlito") {
		ppem = f.size * 96 / 72
	}
	return math.Max(1, math.Round(w*ppem))
}

// colWidthPt converts a column width in characters to points
// (§18.3.1.13: pixels = trunc(((256·width + trunc(128/MDW)) / 256) · MDW)).
func colWidthPt(width, mdw float64) float64 {
	px := math.Trunc((256*width + math.Trunc(128/mdw)) / 256 * mdw)
	return math.Max(px, 0) * pxPt
}

// defaultColPt is the width of columns that set none: defaultColWidth, or
// baseColWidth characters plus 5 pixels of padding and gridline, rounded up
// to a multiple of 8 pixels as Excel does.
func (ws *worksheet) defaultColPt(mdw float64) float64 {
	if ws.defColW > 0 {
		return colWidthPt(ws.defColW, mdw)
	}
	px := math.Ceil((ws.baseColW*mdw+5)/8) * 8
	return px * pxPt
}

// fontLineHeight is the height of a line of text in a font, in points:
// the Windows ascent and descent of the face.
func (c *converter) fontLineHeight(f *xfont) float64 {
	fc := c.fonts.Choose(c.fontName(f), f.bold, f.italic, false)
	return (fc.Asc + fc.Desc) * f.size
}

// autoRowHeight is the height Excel gives a row for text of a given height
// in points: whole pixels plus 2 pixels of padding.
func autoRowHeight(textPt float64) float64 {
	return (math.Ceil(textPt/pxPt-0.05) + 2) * pxPt
}
