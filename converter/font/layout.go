package font

import (
	"strconv"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// The page of the views: a column of this width, cut into strips.
const (
	pageW       = 720.0
	margin      = 36.0
	contentW    = pageW - 2*margin
	stripHeight = 1024.0
)

// scroll lays a view out as rows stacked from the top, then cuts it into
// strips between rows (docs/spec.md §4.1, scroll views).
type scroll struct {
	c       *converter
	rows    []*band
	y       float64
	anchors map[string]float64
}

// band is a row of a view: its height and how it is drawn with its top at y.
type band struct {
	y, h float64
	draw func(p *pen, y float64)
	// keep keeps the band in the strip of the next one (a heading and what
	// it heads).
	keep bool
}

func newScroll(c *converter) *scroll {
	return &scroll{c: c, y: margin, anchors: map[string]float64{}}
}

// add appends a band of height h.
func (s *scroll) add(h float64, draw func(p *pen, y float64)) *band {
	b := &band{y: s.y, h: h, draw: draw}
	s.rows = append(s.rows, b)
	s.y += h
	return b
}

// space adds empty space.
func (s *scroll) space(h float64) { s.y += h }

// anchor names the position the next band starts at, for links.
func (s *scroll) anchor(name string) { s.anchors[name] = s.y }

// finish cuts the bands into strips and draws them as the pages of a view.
func (s *scroll) finish(v *bdf.View) []*strip {
	total := s.y + margin
	// the cuts: at a band boundary after about stripHeight
	cuts := []float64{0}
	start := 0.0
	for i, b := range s.rows {
		if i == 0 || b.y-start < stripHeight || s.rows[i-1].keep {
			continue
		}
		if total-b.y < stripHeight/4 {
			break // not a strip of its own for what is left
		}
		cuts = append(cuts, b.y)
		start = b.y
	}
	cuts = append(cuts, total)
	stripOf := func(y float64) int {
		for i := 1; i < len(cuts); i++ {
			if y < cuts[i] {
				return i
			}
		}
		return len(cuts) - 1
	}
	link := func(a string) string {
		if y, ok := s.anchors[a]; ok {
			return "#page=" + strconv.Itoa(stripOf(y))
		}
		return ""
	}
	var out []*strip
	k := 0
	for i := 0; i+1 < len(cuts); i++ {
		a, b := cuts[i], cuts[i+1]
		cv := s.c.cvs.New()
		cv.Obj.SetBBox(0, 0, f32(pageW), f32(b-a))
		p := &pen{c: s.c, cv: cv, link: link}
		p.rect(0, 0, pageW, b-a, colWhite)
		for ; k < len(s.rows) && s.rows[k].y < b; k++ {
			r := s.rows[k]
			r.draw(p, r.y-a)
		}
		st := &strip{page: v.AddPage(f32(pageW), f32(b-a)), cv: cv}
		out = append(out, st)
	}
	return out
}

// strip is a page of a scroll view, whose layer is set once the objects
// are encoded.
type strip struct {
	page *bdf.Page
	cv   *canvas.Canvas
}
