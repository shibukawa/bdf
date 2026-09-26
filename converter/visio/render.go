package visio

import (
	"math"
	"net/url"
	"strings"

	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
)

// pageCtx is the state of drawing the shapes of one page.
type pageCtx struct {
	c      *converter
	pg     *page
	th     *theme
	scale  float64 // drawing units → paper inches
	height float64 // height of the output page in points
	dm     *drawingml.Drawing
	qs     map[*shape]*quickStyle
	hidden map[string]bool   // indices of the layers that are not shown
	jumps  map[*shape][]jump // line jumps of the connectors (jumps.go)
}

func (c *converter) newPageCtx(pg *page, height float64) *pageCtx {
	p := &pageCtx{c: c, pg: pg, th: pg.theme, scale: pageScale(pg), height: height, dm: c.dm,
		qs: map[*shape]*quickStyle{}, hidden: map[string]bool{}, jumps: map[*shape][]jump{}}
	if p.th == nil {
		p.th = c.d.theme
	}
	if sec := pg.sheet.sections[secKey{"Layer", -1}]; sec != nil {
		for _, k := range sec.order {
			r := sec.rows[k]
			if c, ok := r.cells["Visible"]; ok && num(c.v) == 0 {
				p.hidden[k] = true
			}
		}
	}
	return p
}

// pageScale is the ratio of paper to drawing units of a scaled drawing
// (1 for drawings at full size).
func pageScale(pg *page) float64 {
	ps, ds := num(pg.sheet.cells["PageScale"].v), num(pg.sheet.cells["DrawingScale"].v)
	if ps > 0 && ds > 0 {
		return ps / ds
	}
	return 1
}

// pageSize returns the paper size of a page in points.
func pageSize(pg *page) (float64, float64) {
	s := pageScale(pg)
	w, h := num(pg.sheet.cells["PageWidth"].v)*s*72, num(pg.sheet.cells["PageHeight"].v)*s*72
	if w <= 0 || h <= 0 || w > 1e6 || h > 1e6 {
		return 612, 792
	}
	return w, h
}

// pageMatrix maps page coordinates (drawing units, y up) to points (y down).
func (p *pageCtx) pageMatrix() canvas.Matrix {
	k := 72 * p.scale
	return canvas.Matrix{k, 0, 0, -k, 0, p.height}
}

// localMatrix maps a shape's local coordinates to its parent's: the pin
// point is at PinX, PinY, the shape turns by Angle about it and flips.
func (p *pageCtx) localMatrix(s *shape) canvas.Matrix {
	w, h := p.num(s, "Width", 0), p.num(s, "Height", 0)
	px, py := p.num(s, "PinX", 0), p.num(s, "PinY", 0)
	lx, ly := p.dim(s, "LocPinX", w/2), p.dim(s, "LocPinY", h/2)
	fx, fy := 1.0, 1.0
	if p.flag(s, "FlipX") {
		fx = -1
	}
	if p.flag(s, "FlipY") {
		fy = -1
	}
	return canvas.Translate(px, py).Mul(canvas.Rotate(p.num(s, "Angle", 0) * 180 / math.Pi)).
		Mul(canvas.Scale(fx, fy)).Mul(canvas.Translate(-lx, -ly))
}

// layerHidden reports whether every layer a shape belongs to is hidden.
func (p *pageCtx) layerHidden(s *shape) bool {
	if len(p.hidden) == 0 {
		return false
	}
	v, ok := p.val(s, "LayerMember")
	if !ok || strings.TrimSpace(v) == "" {
		return false
	}
	for _, l := range strings.Split(v, ";") {
		if !p.hidden[strings.TrimSpace(l)] {
			return false
		}
	}
	return true
}

// drawShapes draws shapes in order.
func (p *pageCtx) drawShapes(cv *canvas.Canvas, shapes []*shape, parent canvas.Matrix, top bool, depth int) {
	for _, s := range shapes {
		p.drawShape(cv, s, parent, top, depth)
	}
}

func (p *pageCtx) drawShape(cv *canvas.Canvas, s *shape, parent canvas.Matrix, top bool, depth int) {
	if s.del || depth > maxChain || s.typ == "Guide" {
		return
	}
	m := parent.Mul(p.localMatrix(s))
	hidden := p.layerHidden(s)
	if len(s.kids) == 0 {
		if !hidden {
			p.drawOwn(cv, s, m, top)
		}
		return
	}
	// a group: its own geometry and text in front of or behind its members
	mode := int(p.num(s, "DisplayMode", 2))
	if mode == 0 {
		return
	}
	if mode == 1 && !hidden {
		p.drawOwn(cv, s, m, top)
	}
	p.drawShapes(cv, s.kids, m, false, depth+1)
	if mode != 1 && !hidden {
		p.drawOwn(cv, s, m, top)
	}
}

// drawOwn draws a shape's geometry, picture and text.
func (p *pageCtx) drawOwn(cv *canvas.Canvas, s *shape, m canvas.Matrix, top bool) {
	p.drawGeometry(cv, s, m, top)
	if fd, _ := s.foreignNode(); fd != nil {
		p.drawForeign(cv, s, m)
	}
	p.drawText(cv, s, m)
	p.shapeLink(cv, s, m)
}

// drawGeometry fills and strokes the geometry sections of a shape.
func (p *pageCtx) drawGeometry(cv *canvas.Canvas, s *shape, m canvas.Matrix, top bool) {
	geo := s.geometry()
	if len(geo) == 0 {
		return
	}
	w, h := p.num(s, "Width", 0), p.num(s, "Height", 0)
	var fills [][]*subpath
	var strokes []*subpath
	for _, g := range geo {
		if g.flag("NoShow") {
			continue
		}
		subs := buildGeometry(g, m, w, h)
		if !g.flag("NoFill") {
			var closed []*subpath
			for _, sp := range subs {
				if sp.closed() {
					closed = append(closed, sp)
				}
			}
			if len(closed) > 0 {
				fills = append(fills, closed)
			}
		}
		if !g.flag("NoLine") {
			strokes = append(strokes, subs...)
		}
	}
	if len(fills) == 0 && len(strokes) == 0 {
		return
	}
	f := p.fillStyle(s)
	if f.kind == fillNone || (f.fg&0xff == 0 && f.kind == fillSolid) {
		fills = nil
	}
	l := p.lineStyle(s)
	if l == nil || l.color&0xff == 0 {
		strokes, l = nil, nil
	}
	if len(fills) == 0 && len(strokes) == 0 {
		return
	}
	if l != nil && l.rounding > 0 {
		k := math.Sqrt(math.Abs(m[0]*m[3] - m[1]*m[2]))
		for _, sp := range strokes {
			roundCorners(sp, l.rounding*k)
		}
		for _, fs := range fills {
			for _, sp := range fs {
				roundCorners(sp, l.rounding*k)
			}
		}
	}
	if js := p.jumps[s]; len(js) > 0 && len(strokes) > 0 {
		strokes = applyJumps(strokes, js)
	}
	shd := p.shadowStyle(s, top, true)
	if shd != nil {
		cv.Obj.Save()
		cv.Obj.Shadow(shd.color, f32(shd.blur), f32(shd.dx), f32(shd.dy))
	}
	for _, fs := range fills {
		p.setFill(cv, f, m, w, h)
		cv.Obj.FillPath(cv.Obj.AddPath(toPath(fs, true)), 1)
	}
	if shd != nil && len(fills) > 0 {
		cv.Obj.Restore()
		shd = nil
	}
	if len(strokes) > 0 {
		p.strokeWithArrows(cv, l, strokes)
	}
	if shd != nil {
		cv.Obj.Restore()
	}
	cv.Drawn = true
}

// strokeWithArrows strokes figures and puts the arrowheads at the start of
// the first one and the end of the last one when they are open.
func (p *pageCtx) strokeWithArrows(cv *canvas.Canvas, l *lineStyle, subs []*subpath) {
	first, last := subs[0], subs[len(subs)-1]
	var begin, end bool
	var bTip, bDir, eTip, eDir pt
	if def, ok := arrowDefOf(l.begin); ok && !first.closed() && len(first.segs) > 0 {
		begin, bTip, bDir = true, first.start, first.startDir().mul(-1)
		trimStart(first, def.trim*arrowLength(l.beginSize, l.width))
	}
	if def, ok := arrowDefOf(l.end); ok && !last.closed() && len(last.segs) > 0 {
		end, eTip, eDir = true, last.end(), last.endDir()
		trimEnd(last, def.trim*arrowLength(l.endSize, l.width))
	}
	cv.Obj.Save()
	setLine(cv, l)
	cv.Obj.StrokePath(cv.Obj.AddPath(toPath(subs, false)))
	cv.Obj.Restore()
	if begin {
		drawArrow(cv, l, l.begin, l.beginSize, bTip, bDir)
	}
	if end {
		drawArrow(cv, l, l.end, l.endSize, eTip, eDir)
	}
}

// roundCorners rounds the corners between line segments of a figure with
// arcs of radius r (points), as the Rounding cell asks.
func roundCorners(sp *subpath, r float64) {
	pts := []pt{sp.start}
	for _, s := range sp.segs {
		if s.cubic {
			return // only polygons are rounded
		}
		pts = append(pts, s.p)
	}
	closed := sp.closed()
	if closed {
		pts = pts[:len(pts)-1]
	}
	n := len(pts)
	if n < 3 {
		return
	}
	// corner returns the points where the rounding of corner i starts and
	// ends (both the corner itself when it is not rounded).
	corner := func(i int) (pt, pt) {
		prev, cur, next := pts[(i-1+n)%n], pts[i], pts[(i+1)%n]
		a, b := prev.sub(cur), next.sub(cur)
		la, lb := a.len(), b.len()
		if la == 0 || lb == 0 {
			return cur, cur
		}
		ang := math.Acos(math.Max(-1, math.Min(1, (a.x*b.x+a.y*b.y)/(la*lb))))
		if ang < 1e-3 || math.Pi-ang < 1e-3 {
			return cur, cur
		}
		// no farther than half of each segment
		d := math.Min(r/math.Tan(ang/2), math.Min(la, lb)/2)
		return cur.add(a.mul(d / la)), cur.add(b.mul(d / lb))
	}
	var out subpath
	lineTo := func(q pt) { out.segs = append(out.segs, seg{p: q}) }
	round := func(i int) {
		t0, t1 := corner(i)
		lineTo(t0)
		if t0 != t1 {
			// a cubic close to the arc tangent to both segments
			const k = 0.5523
			c := pts[i]
			out.segs = append(out.segs, seg{cubic: true, c1: t0.add(c.sub(t0).mul(k)), c2: t1.add(c.sub(t1).mul(k)), p: t1})
		}
	}
	if closed {
		// from the end of corner 0 round to its start
		_, t1 := corner(0)
		out.start = t1
		for i := 1; i < n; i++ {
			round(i)
		}
		round(0)
	} else {
		out.start = pts[0]
		for i := 1; i < n-1; i++ {
			round(i)
		}
		lineTo(pts[n-1])
	}
	*sp = out
}

// shapeLink makes the first hyperlink of a shape a link over its bounds.
func (p *pageCtx) shapeLink(cv *canvas.Canvas, s *shape, m canvas.Matrix) {
	for _, k := range s.rowKeys("Hyperlink") {
		addr, _ := p.rowVal(s, "Hyperlink", k, "Address")
		sub, _ := p.rowVal(s, "Hyperlink", k, "SubAddress")
		target := ""
		if addr = strings.TrimSpace(addr); addr != "" {
			if u, err := url.Parse(addr); err == nil && (u.Scheme == "http" || u.Scheme == "https" || u.Scheme == "mailto") {
				target = addr
			}
		} else if sub = strings.TrimSpace(sub); sub != "" {
			// a page, or a shape on a page ("Page-2/Sheet.4")
			name, _, _ := strings.Cut(sub, "/")
			for _, pg := range p.c.d.pages {
				if pg.name == name && p.c.pageNum[pg.id] > 0 {
					target = "#page=" + itoa(p.c.pageNum[pg.id])
				}
			}
		}
		if target == "" {
			continue
		}
		w, h := p.num(s, "Width", 0), p.num(s, "Height", 0)
		x0, y0, x1, y1 := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
		for _, q := range []pt{{0, 0}, {w, 0}, {0, h}, {w, h}} {
			a := apply(m, q)
			x0, y0, x1, y1 = math.Min(x0, a.x), math.Min(y0, a.y), math.Max(x1, a.x), math.Max(y1, a.y)
		}
		cv.Obj.Link(f32(x0), f32(y0), f32(x1-x0), f32(y1-y0), target)
		return
	}
}

// drawPage draws the shapes of a page (not its background pages).
func (p *pageCtx) drawPage(cv *canvas.Canvas) {
	p.collectJumps()
	p.drawShapes(cv, p.pg.shapes, p.pageMatrix(), true, 0)
}
