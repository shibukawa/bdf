package xlsx

import (
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// Spreadsheet drawings (ECMA-376 Part 1 §20.5) anchor DrawingML shapes to
// cells: from one cell corner to another (twoCellAnchor), from a cell with
// a size (oneCellAnchor), or at a position (absoluteAnchor). Offsets are in
// EMU within the cell. Each anchored shape is drawn once into an object of
// its own, in its anchor box, and every tile it covers uses that object.

type anchorPt struct {
	col, row       int
	colOff, rowOff float64 // points
}

type anchored struct {
	kind         string // twoCellAnchor, oneCellAnchor, absoluteAnchor
	from, to     anchorPt
	pos, ext     [2]float64
	elem         *ooxml.Node
	part         string
	toRow, toCol int  // extent of the sheet the anchor needs
	out          bool // left out: it would be drawn into too many tiles
}

func readAnchorPt(n *ooxml.Node) anchorPt {
	num := func(name string) float64 {
		v, _ := strconv.ParseFloat(strings.TrimSpace(n.Child(name).Content()), 64)
		return v
	}
	return anchorPt{col: int(num("col")), row: int(num("row")), colOff: num("colOff") / ooxml.EMUPerPoint, rowOff: num("rowOff") / ooxml.EMUPerPoint}
}

// readAnchors reads the anchored shapes of a drawing part.
func readAnchors(root *ooxml.Node, part string) []*anchored {
	var out []*anchored
	for _, k := range root.Elements() {
		a := &anchored{kind: k.Name, part: part}
		switch k.Name {
		case "twoCellAnchor":
			a.from, a.to = readAnchorPt(k.Child("from")), readAnchorPt(k.Child("to"))
			a.toRow, a.toCol = a.to.row, a.to.col
		case "oneCellAnchor":
			a.from = readAnchorPt(k.Child("from"))
			a.ext = [2]float64{k.Child("ext").AttrEMU("cx", 0), k.Child("ext").AttrEMU("cy", 0)}
			a.toRow, a.toCol = a.from.row+int(a.ext[1]/15)+1, a.from.col+int(a.ext[0]/48)+1
		case "absoluteAnchor":
			a.pos = [2]float64{k.Child("pos").AttrEMU("x", 0), k.Child("pos").AttrEMU("y", 0)}
			a.ext = [2]float64{k.Child("ext").AttrEMU("cx", 0), k.Child("ext").AttrEMU("cy", 0)}
			a.toRow, a.toCol = int((a.pos[1]+a.ext[1])/15)+1, int((a.pos[0]+a.ext[0])/48)+1
		default:
			continue
		}
		for _, e := range k.Elements() {
			switch e.Name {
			case "sp", "grpSp", "graphicFrame", "cxnSp", "pic", "contentPart":
				a.elem = e
			}
		}
		if a.elem != nil {
			a.toRow, a.toCol = min(max(a.toRow, 0), maxRows-1), min(max(a.toCol, 0), maxCols-1)
			out = append(out, a)
		}
	}
	return out
}

func (s *sheetCtx) loadDrawings() {
	if s.ws.drawing == "" {
		return
	}
	r, ok := s.c.pkg.Target(s.ws.part, s.ws.drawing)
	if !ok {
		return
	}
	root, err := s.c.pkg.XML(r.Target)
	if err != nil {
		s.c.warnf("drawing %s: %v", r.Target, err)
		return
	}
	s.drawings = readAnchors(root, r.Target)
	s.c.fillChartCaches(r.Target)
}

// anchorBox places an anchor on the sheet's grid.
func (s *sheetCtx) anchorBox(a *anchored) box {
	at := func(p anchorPt) (float64, float64) {
		c, r := min(p.col, s.nCols), min(p.row, s.nRows)
		return s.cols.pos(c) + math.Min(p.colOff, s.cols.at(c)), s.rows.pos(r) + math.Min(p.rowOff, s.rows.at(r))
	}
	switch a.kind {
	case "twoCellAnchor":
		x0, y0 := at(a.from)
		x1, y1 := at(a.to)
		return box{x0, y0, x1 - x0, y1 - y0}
	case "oneCellAnchor":
		x0, y0 := at(a.from)
		return box{x0, y0, a.ext[0], a.ext[1]}
	}
	return box{a.pos[0], a.pos[1], a.ext[0], a.ext[1]}
}

// maxDrawingTiles bounds the tiles that the drawings of a sheet are drawn
// into, in all. A drawing is an object of its own, but every tile it touches
// uses it, and how far it reaches is the file's to state: one anchored from
// A1 to the last cell of the sheet touches millions of tiles.
const maxDrawingTiles = 1 << 14

// reach is how far past its box a drawing may paint: rotated shapes,
// shadows and line widths reach past the box.
func (b box) reach() float64 { return math.Max(b.w, b.h)*0.5 + 8 }

// drawingTiles returns the number of tiles a drawing in box b is given to;
// none for a box of a negative size, which is not drawn.
func (s *sheetCtx) drawingTiles(b box) int {
	if !(b.w >= 0 && b.h >= 0) {
		return 0
	}
	m := b.reach()
	tx0, ty0, tx1, ty1, ok := s.tileSpan(b.x-m, b.y-m, b.x+b.w+m, b.y+b.h+m)
	if !ok {
		return 0
	}
	return (tx1 - tx0) * (ty1 - ty0)
}

// limitDrawings leaves out the drawings that would be drawn into more tiles
// than are left of maxDrawingTiles after the drawings before them, with a
// warning. It reports whether it left any out: they do not make the view
// larger, which is then set again.
func (s *sheetCtx) limitDrawings() bool {
	room, any := maxDrawingTiles, false
	for _, a := range s.drawings {
		if a.out {
			continue
		}
		n := s.drawingTiles(s.anchorBox(a))
		if n > room {
			s.leaveOut(a)
			any = true
			continue
		}
		room -= n
	}
	return any
}

func (s *sheetCtx) leaveOut(a *anchored) {
	a.out = true
	s.c.warnOnce("drawings:"+s.ws.name, "sheet %q: the drawings cover more than %d tiles; the rest are left out", s.ws.name, maxDrawingTiles)
}

// paintDrawings draws each anchored shape into an object that the tiles
// under it use.
func (s *sheetCtx) paintDrawings() {
	room := maxDrawingTiles
	for _, a := range s.drawings {
		if a.out {
			continue
		}
		b := s.anchorBox(a)
		if !(b.w >= 0 && b.h >= 0) {
			continue
		}
		// the tiles are counted where the drawing is given to them (the
		// view, and with it the boxes, was set again if limitDrawings left
		// drawings out)
		n := s.drawingTiles(b)
		if n > room {
			s.leaveOut(a)
			continue
		}
		cv := s.c.cvs.New()
		s.c.d.DrawAnchored(cv, a.elem, a.part, 0, 0, b.w, b.h)
		if !cv.Drawn {
			continue
		}
		room -= n
		m := b.reach()
		bbox := bdf.Rect{X: f32(-m), Y: f32(-m), W: f32(b.w + 2*m), H: f32(b.h + 2*m)}
		cv.Obj.SetBBox(bbox.X, bbox.Y, bbox.W, bbox.H)
		for _, t := range s.tilesIn(b.x-m, b.y-m, b.x+b.w+m, b.y+b.h+m) {
			ref := t.cv.Share(cv, bbox)
			t.cv.Obj.UseAt(ref, f32(b.x-t.ox), f32(b.y-t.oy))
			t.cv.Used(cv)
			t.cv.Drawn = true
		}
	}
}

// loadComments finds the cells that have notes.
func (s *sheetCtx) loadComments() {
	s.comments = map[[2]int]bool{}
	if s.c.pkg == nil {
		return // a grid
	}
	r, ok := s.c.pkg.RelOfType(s.ws.part, "/comments")
	if !ok {
		return
	}
	n, err := s.c.pkg.XML(r.Target)
	if err != nil {
		return
	}
	for _, cm := range n.Path("commentList").Children("comment") {
		if col, row, ok := parseRef(cm.AttrStr("ref", "")); ok {
			s.comments[[2]int{row, col}] = true
		}
	}
}

// paintMarks draws what Excel shows over cells: the red corner of cells
// with notes and the buttons of AutoFilters.
func (s *sheetCtx) paintMarks() {
	for k := range s.comments {
		if k[0] >= s.nRows || k[1] >= s.nCols {
			continue
		}
		b := s.cellBox(k[0], k[1])
		if m, ok := s.mergeAt(k[0], k[1]); ok {
			b = s.rangeBox(m)
		}
		if b.w <= 0 || b.h <= 0 {
			continue
		}
		sz := 5 * pxPt
		x1, y0 := b.x+b.w, b.y
		for _, t := range s.tilesIn(x1-sz, y0, x1, y0+sz) {
			p := &bdf.Path{}
			p.MoveTo(f32(x1-sz-t.ox), f32(y0-t.oy)).LineTo(f32(x1-t.ox), f32(y0-t.oy)).LineTo(f32(x1-t.ox), f32(y0+sz-t.oy)).Close()
			t.fillColor(bdf.RGB(0xe0, 0x20, 0x20))
			t.cv.Obj.FillPath(t.cv.Obj.AddPath(p), 0)
			t.cv.Drawn = true
		}
	}
	for _, t := range s.tables {
		if !t.filter {
			continue
		}
		for c := t.ref.c0; c <= t.ref.c1 && c < s.nCols; c++ {
			if t.filterButtons[c-t.ref.c0] || t.ref.r0 >= s.nRows {
				continue
			}
			if b := s.cellBox(t.ref.r0, c); b.w > 0 && b.h > 0 {
				s.drawFilterButton(b)
			}
		}
	}
	if af := s.ws.autoFilter; af != nil && af.r0 < s.nRows {
		for c := af.c0; c <= af.c1 && c < s.nCols; c++ {
			if b := s.cellBox(af.r0, c); b.w > 0 && b.h > 0 {
				s.drawFilterButton(b)
			}
		}
	}
}

// paintLinks makes the hyperlinks of cells to web and mail addresses into
// links; links to places in the workbook are left out (no sheet view can
// be a link target).
func (s *sheetCtx) paintLinks() {
	for _, h := range s.ws.links {
		rid := h.RelID("id")
		if rid == "" {
			continue
		}
		r, ok := s.c.pkg.Target(s.ws.part, rid)
		if !ok || !r.External {
			continue
		}
		target := r.Target
		if loc := h.AttrStr("location", ""); loc != "" {
			target += "#" + loc
		}
		u, err := url.Parse(target)
		if err != nil || u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "mailto" {
			continue
		}
		for _, rg := range parseSqref(h.AttrStr("ref", "")) {
			s.putLink(rg, target)
		}
	}
}

// maxLinkTiles bounds the tiles that the links of a sheet are put into
// after their first: a link is put into every tile its range touches, and
// the range is the file's to state.
const maxLinkTiles = 1 << 14

// putLink puts a link over a range into the tiles the range touches, or
// into the first of them only when the others are more than what is left of
// maxLinkTiles after the links before it. (The tiles of the cells of the
// range that are in the file would take, for every such link, a walk through
// those cells.)
func (s *sheetCtx) putLink(rg cellRange, target string) {
	if rg.r0 >= s.nRows || rg.c0 >= s.nCols {
		return
	}
	rg.r1, rg.c1 = min(rg.r1, s.nRows-1), min(rg.c1, s.nCols-1)
	b := s.rangeBox(rg)
	if b.w <= 0 || b.h <= 0 {
		return
	}
	link := func(t *tileCv) {
		t.cv.Obj.Link(f32(b.x-t.ox), f32(b.y-t.oy), f32(b.w), f32(b.h), target)
	}
	tx0, ty0, tx1, ty1, ok := s.tileSpan(b.x, b.y, b.x+b.w, b.y+b.h)
	if !ok {
		return
	}
	more := (tx1-tx0)*(ty1-ty0) - 1
	if more > s.linkRoom {
		s.c.warnOnce("links:"+s.ws.name, "sheet %q: the links cover more than %d tiles; the rest are in their first tile only", s.ws.name, maxLinkTiles)
		link(s.tileAt(tx0, ty0))
		return
	}
	s.linkRoom -= more
	for _, t := range s.tilesIn(b.x, b.y, b.x+b.w, b.y+b.h) {
		link(t)
	}
}

// chartSheetSafe renders a chart sheet as a one-page fixed view.
func (c *converter) chartSheetSafe(ref sheetRef, id string) (v *bdf.View, layers []*canvas.Canvas) {
	v = c.doc.NewView(id, bdf.ViewFixed, ref.name)
	w, h := 720.0, 540.0
	cv := c.cvs.New()
	func() {
		defer func() {
			if r := recover(); r != nil {
				c.warnf("chart sheet %q: internal error: %v", ref.name, r)
			}
		}()
		n, err := c.pkg.XML(ref.part)
		if err != nil {
			c.warnf("chart sheet %q: %v", ref.name, err)
			return
		}
		dr, ok := c.pkg.Target(ref.part, n.Child("drawing").RelID("id"))
		if !ok {
			return
		}
		root, err := c.pkg.XML(dr.Target)
		if err != nil {
			c.warnf("chart sheet %q: %v", ref.name, err)
			return
		}
		anchors := readAnchors(root, dr.Target)
		c.fillChartCaches(dr.Target)
		for _, a := range anchors {
			if a.kind == "absoluteAnchor" && a.ext[0] > 0 && a.ext[1] > 0 {
				w, h = a.pos[0]+a.ext[0], a.pos[1]+a.ext[1]
				break
			}
		}
		cv.Obj.FillColor(bdf.RGB(255, 255, 255))
		cv.Obj.FillRect(0, 0, f32(w), f32(h))
		for _, a := range anchors {
			// a chart sheet's chart fills the sheet unless it has a size
			b := box{a.pos[0], a.pos[1], a.ext[0], a.ext[1]}
			if a.kind != "absoluteAnchor" || b.w <= 0 || b.h <= 0 {
				b = box{0, 0, w, h}
			}
			c.d.DrawAnchored(cv, a.elem, a.part, b.x, b.y, b.w, b.h)
		}
	}()
	cv.Obj.SetBBox(0, 0, f32(w), f32(h))
	cv.Drawn = true
	v.AddPage(f32(w), f32(h), bdf.Layer{Role: bdf.RoleBody})
	return v, []*canvas.Canvas{cv}
}
