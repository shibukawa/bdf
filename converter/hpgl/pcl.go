package hpgl

import (
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// The page set-up of a PCL printer, for HP-GL/2 sent to one (a PCL job that
// enters HP-GL/2 with ESC % # B): the page (size and orientation), its
// logical page and the picture frame, where HP-GL/2 draws, scaled from its
// plot size. Coordinates of the page are plotter units from its top left
// corner, y down.

type pcl struct {
	used        bool    // the job set up a page
	size        int     // page size code (ESC & l # A)
	orientation int     // 0 portrait, 1 landscape, 2 reverse portrait, 3 reverse landscape
	topMargin   float64 // plotter units
	units       float64 // PCL units per inch (ESC & u # D)
	frameW      float64 // picture frame (plotter units; 0: default)
	frameH      float64
	anchor      *[2]float64 // picture frame anchor (logical page coordinates)
	plotW       float64     // HP-GL/2 plot size (plotter units; 0: the frame's)
	plotH       float64
	cursor      cad.Point // the cursor, for the anchor of the picture frame
}

// pclPage is the physical page of a PCL printer in its orientation.
type pclPage struct{ w, h float64 }

// pageSizes are the PCL page sizes (portrait, inches × 1016).
var pageSizes = map[int][2]float64{
	1:   {7.25 * 1016, 10.5 * 1016}, // Executive
	2:   {8.5 * 1016, 11 * 1016},    // Letter
	3:   {8.5 * 1016, 14 * 1016},    // Legal
	6:   {11 * 1016, 17 * 1016},     // Ledger
	25:  {148 * 40, 210 * 40},       // A5
	26:  {210 * 40, 297 * 40},       // A4
	27:  {297 * 40, 420 * 40},       // A3
	45:  {182 * 40, 257 * 40},       // JIS B5
	46:  {257 * 40, 364 * 40},       // JIS B4
	71:  {100 * 40, 148 * 40},       // Hagaki
	100: {176 * 40, 250 * 40},       // ISO B5
}

// printer reports whether the job is a PCL printer's: it set up a page,
// and nothing made it a plotter's.
func (c *converter) printer() bool {
	return c.pcl != nil && c.pcl.used && !c.gl.standaloneSet
}

// pclState returns the PCL settings of the job.
func (c *converter) pclState() *pcl {
	if c.pcl == nil {
		c.pcl = &pcl{}
		c.pcl.reset()
	}
	return c.pcl
}

// pclSetup returns the page set-up of the job, making it a PCL printer's.
func (c *converter) pclSetup() *pcl {
	p := c.pclState()
	p.used = true
	return p
}

func (p *pcl) reset() {
	used := p.used
	*p = pcl{size: 2, units: 300, topMargin: 1016 / 2}
	p.used = used
}

// dims returns the physical page in its orientation.
func (p *pcl) dims() (w, h float64) {
	s, ok := pageSizes[p.size]
	if !ok {
		s = pageSizes[2]
	}
	w, h = s[0], s[1]
	if p.orientation%2 == 1 {
		w, h = h, w
	}
	return w, h
}

// offset returns the left edge of the logical page on the physical one.
func (p *pcl) offset() float64 {
	if p.orientation%2 == 1 {
		return 0.2 * 1016
	}
	return 0.25 * 1016
}

func (p *pcl) page() *pclPage {
	w, h := p.dims()
	return &pclPage{w: w, h: h}
}

// logical returns the transform of cursor coordinates into the page: from
// the left edge of the logical page and the top margin.
func (p *pcl) logical() canvas.Matrix { return canvas.Translate(p.offset(), p.topMargin) }

// frameInfo is where HP-GL/2 draws on the page.
type frameInfo struct {
	base canvas.Matrix // HP-GL/2 plotter units -> page
	w, h float64       // hard-clip limits (HP-GL/2 plotter units)
}

// frame returns the picture frame: by default the width of the logical page
// between the top and bottom margins, anchored at the top margin.
func (p *pcl) frame() frameInfo {
	pw, ph := p.dims()
	fw, fh := p.frameW, p.frameH
	if fw <= 0 {
		fw = pw - 2*p.offset()
	}
	if fh <= 0 {
		fh = ph - 2*p.topMargin
	}
	var ax, ay float64 // the anchor, in cursor coordinates
	if p.anchor != nil {
		ax, ay = p.anchor[0], p.anchor[1]
	}
	w, h := p.plotW, p.plotH
	if w <= 0 || h <= 0 {
		w, h = fw, fh
	}
	kx, ky := fw/w, fh/h
	base := p.logical().Mul(canvas.Matrix{kx, 0, 0, -ky, ax, ay + fh})
	return frameInfo{base: base, w: w, h: h}
}

// pictureFrame carries out the picture frame commands (ESC * c).
func (p *pcl) pictureFrame(cmd byte, v float64) {
	if !valid(v) || v < 0 {
		return
	}
	switch cmd {
	case 'X':
		p.frameW = v * 1016 / 720
	case 'Y':
		p.frameH = v * 1016 / 720
	case 'K':
		p.plotW = v * 1016
	case 'L':
		p.plotH = v * 1016
	case 'T':
		// the anchor at the cursor
		p.anchor = &[2]float64{p.cursor.X, p.cursor.Y}
	}
}

// pageCommand carries out the page commands (ESC & l, ESC & u): the page
// size, orientation and top margin make the job a PCL printer's.
func (c *converter) pageCommand(g, cmd byte, v float64) {
	if !valid(v) {
		return
	}
	switch {
	case g == 'l' && cmd == 'A':
		if _, ok := pageSizes[int(v)]; ok {
			c.pclSetup().size = int(v)
		}
	case g == 'l' && cmd == 'O':
		if v >= 0 && v <= 3 {
			c.pclSetup().orientation = int(v)
		}
	case g == 'l' && cmd == 'E':
		// the top margin in lines of 6 per inch
		c.pclSetup().topMargin = v * 1016 / 6
	case g == 'u' && cmd == 'D':
		if v > 0 && v <= 7200 {
			c.pclState().units = v
		}
	}
}
