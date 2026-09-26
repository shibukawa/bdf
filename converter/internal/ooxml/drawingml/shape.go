package drawingml

import (
	"math"
	"net/url"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// xform is a shape's placement on the page (points, degrees clockwise).
type xform struct {
	X, Y, W, H   float64
	Rot          float64
	FlipH, FlipV bool
}

func parseXfrm(x *ooxml.Node) (xform, bool) {
	if x == nil {
		return xform{}, false
	}
	off, ext := x.Child("off"), x.Child("ext")
	if off == nil && ext == nil {
		return xform{}, false
	}
	return xform{
		X: off.AttrEMU("x", 0), Y: off.AttrEMU("y", 0), W: ext.AttrEMU("cx", 0), H: ext.AttrEMU("cy", 0),
		Rot:   float64(x.AttrInt("rot", 0)) / 60000,
		FlipH: x.AttrBool("flipH", false), FlipV: x.AttrBool("flipV", false),
	}, true
}

// matrix maps the shape's local box (0,0)-(W,H) to the page.
func (x xform) matrix() canvas.Matrix {
	fh, fv := 1.0, 1.0
	if x.FlipH {
		fh = -1
	}
	if x.FlipV {
		fv = -1
	}
	if x.Rot == 0 && !x.FlipH && !x.FlipV {
		return canvas.Translate(x.X, x.Y)
	}
	return canvas.Translate(x.X+x.W/2, x.Y+x.H/2).Mul(canvas.Rotate(x.Rot)).Mul(canvas.Scale(fh, fv)).Mul(canvas.Translate(-x.W/2, -x.H/2))
}

// textMatrix is matrix without mirroring: text is never drawn mirrored; a
// vertical flip turns it upside down.
func (x xform) textMatrix() canvas.Matrix {
	rot := x.Rot
	if x.FlipV {
		rot += 180
	}
	if rot == 0 {
		return canvas.Translate(x.X, x.Y)
	}
	return canvas.Translate(x.X+x.W/2, x.Y+x.H/2).Mul(canvas.Rotate(rot)).Mul(canvas.Translate(-x.W/2, -x.H/2))
}

// bounds returns the axis-aligned bounding box on the page.
func (x xform) bounds() (x0, y0, x1, y1 float64) {
	m := x.matrix()
	x0, y0 = math.Inf(1), math.Inf(1)
	x1, y1 = math.Inf(-1), math.Inf(-1)
	for _, p := range [][2]float64{{0, 0}, {x.W, 0}, {0, x.H}, {x.W, x.H}} {
		px, py := m.Apply(p[0], p[1])
		x0, y0, x1, y1 = math.Min(x0, px), math.Min(y0, py), math.Max(x1, px), math.Max(y1, py)
	}
	return
}

// groupCtx maps the child coordinate space of a group onto the page.
type groupCtx struct {
	parent *groupCtx
	xf     xform // the group's placement on the page
	chOff  [2]float64
	chExt  [2]float64
	spPr   *ooxml.Node // grpSpPr, for grpFill
	part   string
}

// place converts a child placement into page coordinates. Groups are
// flattened: the child's box is scaled into the group, then the group's
// flips and rotation are applied around its center and folded into the
// child's own rotation and flips.
func (g *groupCtx) place(c xform) xform {
	if g == nil {
		return c
	}
	sx, sy := 1.0, 1.0
	if g.chExt[0] > 0 {
		sx = g.xf.W / g.chExt[0]
	}
	if g.chExt[1] > 0 {
		sy = g.xf.H / g.chExt[1]
	}
	bx := g.xf.X + (c.X-g.chOff[0])*sx
	by := g.xf.Y + (c.Y-g.chOff[1])*sy
	cx, cy := bx+c.W*sx/2, by+c.H*sy/2
	// A child turned sideways is stretched along its own axes the other way.
	r := math.Mod(math.Abs(c.Rot), 180)
	if r > 45 && r < 135 {
		sx, sy = sy, sx
	}
	w, h := c.W*sx, c.H*sy
	gcx, gcy := g.xf.X+g.xf.W/2, g.xf.Y+g.xf.H/2
	dx, dy := cx-gcx, cy-gcy
	if g.xf.FlipH {
		dx = -dx
	}
	if g.xf.FlipV {
		dy = -dy
	}
	s, co := math.Sincos(g.xf.Rot * math.Pi / 180)
	ncx, ncy := gcx+dx*co-dy*s, gcy+dx*s+dy*co
	rot := c.Rot
	if g.xf.FlipH != g.xf.FlipV {
		rot = -rot
	}
	rot += g.xf.Rot
	out := xform{X: ncx - w/2, Y: ncy - h/2, W: w, H: h, Rot: rot, FlipH: c.FlipH != g.xf.FlipH, FlipV: c.FlipV != g.xf.FlipV}
	return g.parent.place(out)
}

// shape is a shape element with what it inherits from placeholders.
type shape struct {
	n      *ooxml.Node
	part   string
	inh    []*ooxml.Node // layout placeholder, master placeholder (most specific first)
	inhPrt []string
	ph     *ooxml.Node
	grp    *groupCtx
	box    *xform // placement given by the host (word processing drawings)
}

func (sh *shape) spPr() *ooxml.Node {
	if p := sh.n.Child("spPr"); p != nil {
		return p
	}
	return sh.n.Child("grpSpPr")
}

// spPrs returns the shape properties of the shape and its placeholders.
func (sh *shape) spPrs() ([]*ooxml.Node, []string) {
	out := []*ooxml.Node{sh.spPr()}
	parts := []string{sh.part}
	for i, n := range sh.inh {
		out = append(out, n.Child("spPr"))
		parts = append(parts, sh.inhPrt[i])
	}
	return out, parts
}

func (sh *shape) xform() (xform, bool) {
	if sh.box != nil {
		return *sh.box, true
	}
	cands := []*ooxml.Node{sh.spPr().Child("xfrm"), sh.n.Child("xfrm")}
	for _, n := range sh.inh {
		cands = append(cands, n.Path("spPr", "xfrm"), n.Child("xfrm"))
	}
	// Offset and size are taken separately, so that a shape that only
	// resizes its placeholder keeps the inherited position.
	var first, off, ext *ooxml.Node
	for _, x := range cands {
		if x == nil {
			continue
		}
		if first == nil {
			first = x
		}
		if off == nil {
			off = x.Child("off")
		}
		if ext == nil {
			ext = x.Child("ext")
		}
	}
	if off == nil && ext == nil {
		return xform{}, false
	}
	xf := xform{
		X: off.AttrEMU("x", 0), Y: off.AttrEMU("y", 0), W: ext.AttrEMU("cx", 0), H: ext.AttrEMU("cy", 0),
		Rot:   float64(first.AttrInt("rot", 0)) / 60000,
		FlipH: first.AttrBool("flipH", false), FlipV: first.AttrBool("flipV", false),
	}
	return sh.grp.place(xf), true
}

func (sh *shape) style() *ooxml.Node {
	if st := sh.n.Child("style"); st != nil {
		return st
	}
	for _, n := range sh.inh {
		if st := n.Child("style"); st != nil {
			return st
		}
	}
	return nil
}

func (sh *shape) geometry(w, h float64) (*geometry, string) {
	sps, _ := sh.spPrs()
	for _, sp := range sps {
		if sp.Child("prstGeom") != nil || sp.Child("custGeom") != nil {
			return shapeGeometry(sp, w, h)
		}
	}
	return rectGeometry(w, h), ""
}

func (s *Drawing) shapeFill(sh *shape) fill {
	sps, parts := sh.spPrs()
	for i, sp := range sps {
		if e := fillElem(sp); e != nil {
			return s.resolveFill(e, parts[i], s.cc, sh.grp)
		}
	}
	return s.styleFill(sh.style().Child("fillRef"), s.cc)
}

func (s *Drawing) shapeLine(sh *shape) *line {
	sps, parts := sh.spPrs()
	var lns []*ooxml.Node
	var lparts []string
	for i, sp := range sps {
		if ln := sp.Child("ln"); ln != nil {
			lns = append(lns, ln)
			lparts = append(lparts, parts[i])
		}
	}
	return s.resolveLine(lns, lparts, sh.style().Child("lnRef"), s.cc)
}

func (s *Drawing) shapeShadow(sh *shape) *shadow {
	sps, _ := sh.spPrs()
	var effects []*ooxml.Node
	for _, sp := range sps {
		effects = append(effects, sp.Child("effectLst"))
	}
	return s.resolveShadow(effects, sh.style().Child("effectRef"), s.cc)
}

func cNvPr(n *ooxml.Node) *ooxml.Node {
	for _, k := range n.Kids {
		if strings.HasPrefix(k.Name, "nv") {
			if c := k.Child("cNvPr"); c != nil {
				return c
			}
		}
	}
	return nil
}

// altText returns the alternative text of a shape, picture, graphic frame
// or group (its description, else its title); false when it has none or is
// marked as decorative.
func altText(n *ooxml.Node) (string, bool) {
	pr := cNvPr(n)
	for _, ext := range pr.Path("extLst").Children("ext") {
		if ext.Child("decorative").AttrBool("val", false) {
			return "", false
		}
	}
	for _, name := range []string{"descr", "title"} {
		if v := strings.TrimSpace(pr.AttrStr(name, "")); v != "" {
			return v, true
		}
	}
	return "", false
}

// beginFigure opens a FIGURE for what n draws when it has alternative text;
// endFigure closes it.
func beginFigure(cv *canvas.Canvas, n *ooxml.Node) bool {
	alt, ok := altText(n)
	if ok {
		cv.Obj.Mark(bdf.MarkFigure, alt)
	}
	return ok
}

func endFigure(cv *canvas.Canvas, open bool) {
	if open {
		cv.Obj.Mark(bdf.MarkEnd, "")
	}
}

func (s *Drawing) drawElem(cv *canvas.Canvas, k *ooxml.Node, part string, grp *groupCtx) {
	switch k.Name {
	case "sp", "cxnSp", "pic", "graphicFrame", "grpSp", "wsp", "wgp":
	case "contentPart":
		s.c.warnOnce("ink", "ink (content parts) is not supported")
		return
	default:
		return
	}
	if cNvPr(k).AttrBool("hidden", false) {
		return
	}
	sh := &shape{n: k, part: part, grp: grp, ph: PlaceholderOf(k)}
	if sh.ph != nil {
		inh, draw := s.host.Placeholder(sh.ph, part)
		if !draw {
			return
		}
		for _, i := range inh {
			sh.inh = append(sh.inh, i.Node)
			sh.inhPrt = append(sh.inhPrt, i.Part)
		}
	}
	switch k.Name {
	case "grpSp", "wgp":
		s.drawGroup(cv, sh)
	case "sp", "cxnSp", "wsp":
		s.drawSp(cv, sh)
	case "pic":
		s.drawPic(cv, sh)
	case "graphicFrame":
		s.drawFrame(cv, sh)
	}
}

func (s *Drawing) drawGroup(cv *canvas.Canvas, sh *shape) {
	gp := sh.n.Child("grpSpPr")
	x := gp.Child("xfrm")
	own, ok := parseXfrm(x)
	g := &groupCtx{parent: sh.grp, spPr: gp, part: sh.part}
	if ok {
		// The group box stays in its parent's coordinates: place() maps
		// children into it and then recurses into the parent group.
		g.xf = own
		chOff, chExt := x.Child("chOff"), x.Child("chExt")
		g.chOff = [2]float64{chOff.AttrEMU("x", own.X), chOff.AttrEMU("y", own.Y)}
		g.chExt = [2]float64{chExt.AttrEMU("cx", own.W), chExt.AttrEMU("cy", own.H)}
	} else {
		// no placement: children keep their coordinates
		g.xf = xform{W: 1, H: 1}
		g.chExt = [2]float64{1, 1}
	}
	fig := beginFigure(cv, sh.n)
	for _, k := range sh.n.Kids {
		s.drawElem(cv, k, sh.part, g)
	}
	endFigure(cv, fig)
}

func (s *Drawing) drawSp(cv *canvas.Canvas, sh *shape) {
	xf, ok := sh.xform()
	if !ok {
		return
	}
	geo, unknown := sh.geometry(xf.W, xf.H)
	if unknown != "" {
		s.c.warnOnce("prst:"+unknown, "unknown preset shape %q drawn as a rectangle", unknown)
	}
	fl := s.shapeFill(sh)
	ln := s.shapeLine(sh)
	if fl.kind != fillNone || ln != nil {
		// the alternative text describes the drawing; the text stays text
		fig := beginFigure(cv, sh.n)
		cv.Obj.Save()
		cv.Transform(xf.matrix())
		s.paintGeometry(cv, geo, fl, ln, s.shapeShadow(sh), xf.W, xf.H)
		cv.Obj.Restore()
		endFigure(cv, fig)
		cv.Drawn = true
	}
	s.drawShapeText(cv, sh, xf, geo)
	s.shapeLink(cv, sh, xf)
}

// paintGeometry fills and strokes the paths of a geometry in local coordinates.
func (s *Drawing) paintGeometry(cv *canvas.Canvas, geo *geometry, fl fill, ln *line, shd *shadow, w, h float64) {
	refs := make([]bdf.PathRef, len(geo.paths))
	for i, p := range geo.paths {
		refs[i] = cv.Obj.AddPath(p.path)
	}
	hasFill := false
	if fl.kind != fillNone {
		for _, p := range geo.paths {
			if p.fill != "none" {
				hasFill = true
			}
		}
	}
	if shd != nil && (hasFill || ln != nil) {
		cv.Obj.Save()
		cv.Obj.Shadow(shd.c.bdf(), f32(shd.blur), f32(shd.dx), f32(shd.dy))
	} else {
		shd = nil
	}
	if hasFill {
		for i, p := range geo.paths {
			if p.fill == "none" {
				continue
			}
			s.fillPath(cv, refs[i], fl.adjust(p.fill), w, h)
		}
		if shd != nil {
			cv.Obj.Restore()
			shd = nil
		}
	}
	if ln != nil {
		cv.Obj.Save()
		setLine(cv, ln, w, h)
		for i, p := range geo.paths {
			if !p.stroke || p.empty {
				continue
			}
			cv.Obj.StrokePath(refs[i])
		}
		for _, p := range geo.paths {
			if p.stroke && !p.empty && !p.closed {
				cv.Obj.Save()
				cv.Obj.Dash(nil, 0)
				arrowHead(cv, ln, ln.head, p.start[0], p.start[1], -p.startDir[0], -p.startDir[1])
				arrowHead(cv, ln, ln.tail, p.end[0], p.end[1], p.endDir[0], p.endDir[1])
				cv.Obj.Restore()
			}
		}
		cv.Obj.Restore()
	}
	if shd != nil {
		cv.Obj.Restore()
	}
}

func (s *Drawing) fillPath(cv *canvas.Canvas, ref bdf.PathRef, f fill, w, h float64) {
	if f.kind == fillBlip {
		cv.Obj.Save()
		cv.Obj.ClipPath(ref, 0)
		s.drawBlip(cv, f.blip, f.part, f.cc, 0, 0, w, h, true)
		cv.Obj.Restore()
		return
	}
	if s.setFill(cv, f, w, h) {
		cv.Obj.FillPath(ref, 0)
	}
}

// drawBlip draws the image of an a:blipFill / p:blipFill into x,y,w,h
// (stretched, cropped by srcRect, or tiled).
func (s *Drawing) drawBlip(cv *canvas.Canvas, bf *ooxml.Node, part string, cc *colorCtx, x, y, w, h float64, clipped bool) bool {
	blip := bf.Child("blip")
	rid := blip.RelID("embed")
	if rid == "" {
		rid = blip.RelID("link")
	}
	e := s.c.image(part, rid)
	if !e.ok {
		return false
	}
	cv.Drawn = true
	if cc == nil {
		cc = s.cc
	}
	rc := blipRecolor(blip, cc)
	if e.mf != nil {
		s.drawMetafileBlip(cv, bf, e, rc, x, y, w, h)
		return true
	}
	if rc != nil {
		e = s.c.recolored(e, rc)
	}
	ref := cv.Image(e.hash)
	if a := blip.Child("alphaModFix"); a != nil {
		cv.Obj.Alpha(f32(a.AttrPct("amt", 1)))
	}
	if t := bf.Child("tile"); t != nil && e.w > 0 {
		sx, sy := t.AttrPct("sx", 1), t.AttrPct("sy", 1)
		tw, th := e.w*0.75*sx, e.h*0.75*sy
		ox, oy := x+t.AttrEMU("tx", 0), y+t.AttrEMU("ty", 0)
		algn := t.AttrStr("algn", "tl")
		switch {
		case strings.Contains(algn, "r"):
			ox += w - tw
		case algn == "t" || algn == "b" || algn == "ctr":
			ox += (w - tw) / 2
		}
		switch {
		case strings.HasPrefix(algn, "b"):
			oy += h - th
		case algn == "l" || algn == "r" || algn == "ctr":
			oy += (h - th) / 2
		}
		p := bdf.Pattern(ref, bdf.RepeatBoth)
		p.Matrix = [6]float32{f32(0.75 * sx), 0, 0, f32(0.75 * sy), f32(ox), f32(oy)}
		cv.Obj.FillPaint(cv.Obj.AddPaint(p))
		if !clipped {
			cv.Obj.FillRect(f32(x), f32(y), f32(w), f32(h))
		} else {
			cv.Obj.FillRect(f32(x-1), f32(y-1), f32(w+2), f32(h+2))
		}
		return true
	}
	// stretch, inset by fillRect
	if fr := bf.Path("stretch", "fillRect"); fr != nil {
		l, t, r, b := fr.AttrPct("l", 0), fr.AttrPct("t", 0), fr.AttrPct("r", 0), fr.AttrPct("b", 0)
		x, y, w, h = x+l*w, y+t*h, w*(1-l-r), h*(1-t-b)
	}
	sr := bf.Child("srcRect")
	if sr == nil || e.w <= 0 || e.h <= 0 {
		cv.Obj.Image(ref, f32(x), f32(y), f32(w), f32(h))
		return true
	}
	l, t, r, b := sr.AttrPct("l", 0), sr.AttrPct("t", 0), sr.AttrPct("r", 0), sr.AttrPct("b", 0)
	if l == 0 && t == 0 && r == 0 && b == 0 {
		cv.Obj.Image(ref, f32(x), f32(y), f32(w), f32(h))
		return true
	}
	// Source window in pixels; parts outside the image (negative crops)
	// leave the corresponding destination empty.
	u0, u1 := l*e.w, (1-r)*e.w
	v0, v1 := t*e.h, (1-b)*e.h
	if u1 <= u0 || v1 <= v0 {
		return true
	}
	cu0, cu1 := math.Max(u0, 0), math.Min(u1, e.w)
	cv0, cv1 := math.Max(v0, 0), math.Min(v1, e.h)
	if cu1 <= cu0 || cv1 <= cv0 {
		return true
	}
	dx0 := x + (cu0-u0)/(u1-u0)*w
	dx1 := x + (cu1-u0)/(u1-u0)*w
	dy0 := y + (cv0-v0)/(v1-v0)*h
	dy1 := y + (cv1-v0)/(v1-v0)*h
	cv.Obj.ImageSub(ref, f32(cu0), f32(cv0), f32(cu1-cu0), f32(cv1-cv0), f32(dx0), f32(dy0), f32(dx1-dx0), f32(dy1-dy0))
	return true
}

// drawMetafileBlip replays a metafile picture into x,y,w,h, cropped by
// srcRect and inset by fillRect.
func (s *Drawing) drawMetafileBlip(cv *canvas.Canvas, bf *ooxml.Node, e *imageEntry, rc *recolor, x, y, w, h float64) {
	if fr := bf.Path("stretch", "fillRect"); fr != nil {
		l, t, r, b := fr.AttrPct("l", 0), fr.AttrPct("t", 0), fr.AttrPct("r", 0), fr.AttrPct("b", 0)
		x, y, w, h = x+l*w, y+t*h, w*(1-l-r), h*(1-t-b)
	}
	cx, cy, cw, ch := x, y, w, h
	if sr := bf.Child("srcRect"); sr != nil {
		// the whole picture is scaled so that the source window fills the box
		l, t, r, b := sr.AttrPct("l", 0), sr.AttrPct("t", 0), sr.AttrPct("r", 0), sr.AttrPct("b", 0)
		if fw, fh := 1-l-r, 1-t-b; fw > 0 && fh > 0 {
			w, h = w/fw, h/fh
			x, y = x-l*w, y-t*h
		}
	}
	if a := bf.Path("blip", "alphaModFix"); a != nil {
		cv.Obj.Save()
		cv.Obj.Alpha(f32(a.AttrPct("amt", 1)))
		defer cv.Obj.Restore()
	}
	cv.Obj.Save()
	cv.Obj.ClipRect(f32(cx), f32(cy), f32(cw), f32(ch))
	s.drawMetafile(cv, e.mf, rc, x, y, w, h)
	cv.Obj.Restore()
}

func (s *Drawing) drawPic(cv *canvas.Canvas, sh *shape) {
	if xf, ok := sh.xform(); ok {
		fig := beginFigure(cv, sh.n)
		s.drawPicAt(cv, sh, xf)
		endFigure(cv, fig)
	}
}

func (s *Drawing) drawPicAt(cv *canvas.Canvas, sh *shape, xf xform) {
	geo, _ := sh.geometry(xf.W, xf.H)
	bf := sh.n.Child("blipFill")
	if bf == nil {
		for _, n := range sh.inh {
			if bf = n.Child("blipFill"); bf != nil {
				break
			}
		}
	}
	fl := s.shapeFill(sh)
	ln := s.shapeLine(sh)
	shd := s.shapeShadow(sh)
	cv.Obj.Save()
	cv.Transform(xf.matrix())
	if fl.kind != fillNone {
		s.paintGeometry(cv, geo, fl, nil, shd, xf.W, xf.H)
		shd = nil
	}
	if bf != nil {
		cv.Obj.Save()
		rect := len(geo.paths) == 1 && len(geo.paths[0].path.Verbs) == 1 && geo.paths[0].path.Verbs[0] == bdf.VerbRect
		if shd != nil {
			cv.Obj.Shadow(shd.c.bdf(), f32(shd.blur), f32(shd.dx), f32(shd.dy))
		}
		if !rect {
			for _, p := range geo.paths {
				if p.fill != "none" {
					cv.Obj.ClipPath(cv.Obj.AddPath(p.path), 0)
					break
				}
			}
		}
		s.drawBlip(cv, bf, sh.part, s.cc, 0, 0, xf.W, xf.H, !rect)
		cv.Obj.Restore()
	}
	if ln != nil {
		s.paintGeometry(cv, geo, fill{}, ln, nil, xf.W, xf.H)
	}
	cv.Obj.Restore()
	cv.Drawn = true
	s.shapeLink(cv, sh, xf)
}

// shapeLink emits the click hyperlink of a shape over its bounds.
func (s *Drawing) shapeLink(cv *canvas.Canvas, sh *shape, xf xform) {
	h := cNvPr(sh.n).Child("hlinkClick")
	if h == nil {
		return
	}
	target := s.linkTarget(h, sh.part)
	if target == "" {
		return
	}
	x0, y0, x1, y1 := xf.bounds()
	cv.Obj.Link(f32(x0), f32(y0), f32(x1-x0), f32(y1-y0), target)
}

// linkTarget resolves an a:hlinkClick to a URL or, through the host, a
// link within the document.
func (s *Drawing) linkTarget(h *ooxml.Node, part string) string {
	action := h.AttrStr("action", "")
	if strings.HasPrefix(action, "ppaction://hlinkshowjump") {
		return s.host.Link(action, ooxml.Rel{})
	}
	r, ok := s.c.pkg.Target(part, h.RelID("id"))
	if !ok {
		return ""
	}
	if r.External {
		if u, err := url.Parse(r.Target); err == nil && (u.Scheme == "http" || u.Scheme == "https" || u.Scheme == "mailto") {
			return r.Target
		}
		return ""
	}
	return s.host.Link(action, r)
}

// fillPage fills a page of w×h.
func (s *Drawing) fillPage(cv *canvas.Canvas, f fill, w, h float64) {
	switch f.kind {
	case fillNone:
		return
	case fillBlip:
		s.drawBlip(cv, f.blip, f.part, f.cc, 0, 0, w, h, false)
	default:
		if s.setFill(cv, f, w, h) {
			cv.Obj.FillRect(0, 0, f32(w), f32(h))
		}
	}
	cv.Drawn = true
}

// drawFrame draws a graphic frame: a table, a chart, a SmartArt diagram
// (from its pre-rendered drawing) or an OLE object's preview picture.
func (s *Drawing) drawFrame(cv *canvas.Canvas, sh *shape) {
	xf, ok := sh.xform()
	if !ok {
		return
	}
	gd := sh.n.Path("graphic", "graphicData")
	uri := gd.AttrStr("uri", "")
	if gd.Child("tbl") != nil {
		// tables have their own structure
		s.drawTable(cv, sh, xf, gd.Child("tbl"))
		return
	}
	fig := beginFigure(cv, sh.n)
	defer func() { endFigure(cv, fig) }()
	switch {
	case strings.HasSuffix(uri, "/chart"):
		s.drawChart(cv, sh, xf, gd.Child("chart").RelID("id"))
	case strings.HasSuffix(uri, "/diagram"):
		s.drawDiagram(cv, sh, xf, gd.Child("relIds"))
	default:
		// OLE objects and other embeddings: draw their preview picture.
		var pic *ooxml.Node
		var find func(n *ooxml.Node)
		find = func(n *ooxml.Node) {
			for _, k := range n.Kids {
				if pic != nil {
					return
				}
				if k.Name == "pic" {
					pic = k
					return
				}
				find(k)
			}
		}
		find(gd)
		if pic == nil {
			s.c.warnOnce("frame:"+uri, "graphic frame content %s is not supported", uri)
			return
		}
		ps := &shape{n: pic, part: sh.part, grp: sh.grp}
		if !fig {
			fig = beginFigure(cv, pic) // the preview's own description
		}
		if pxf, ok := ps.xform(); ok {
			s.drawPicAt(cv, ps, pxf)
		} else {
			s.drawPicAt(cv, ps, xf) // the preview fills the frame
		}
	}
}

// drawDiagram draws a SmartArt diagram from the drawing part PowerPoint
// stores next to its data model (the shapes as last laid out).
func (s *Drawing) drawDiagram(cv *canvas.Canvas, sh *shape, xf xform, relIds *ooxml.Node) {
	var drawingPart string
	if r, ok := s.c.pkg.Target(sh.part, relIds.RelID("dm")); ok {
		if dm, err := s.c.pkg.XML(r.Target); err == nil {
			for _, ext := range dm.Path("extLst").Children("ext") {
				if d := ext.Child("dataModelExt"); d != nil {
					if dr, ok := s.c.pkg.Target(sh.part, d.AttrStr("relId", "")); ok {
						drawingPart = dr.Target
					}
				}
			}
		}
	}
	if drawingPart == "" {
		s.c.warnOnce("smartart", "SmartArt diagram without a drawing part is not drawn")
		return
	}
	d, err := s.c.pkg.XML(drawingPart)
	if err != nil {
		s.c.warnf("SmartArt drawing: %v", err)
		return
	}
	// xf is already on the page; the drawing's coordinates are relative to
	// the frame.
	g := &groupCtx{xf: xform{X: xf.X, Y: xf.Y, W: xf.W, H: xf.H}, chExt: [2]float64{xf.W, xf.H}, part: drawingPart}
	for _, k := range d.Path("spTree").Elements() {
		s.drawElem(cv, k, drawingPart, g)
	}
}
