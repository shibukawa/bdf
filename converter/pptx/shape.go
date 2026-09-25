package pptx

import (
	"math"
	"net/url"
	"strings"

	"github.com/shibukawa/bdf"
)

// xform is a shape's placement on the slide (points, degrees clockwise).
type xform struct {
	X, Y, W, H   float64
	Rot          float64
	FlipH, FlipV bool
}

func parseXfrm(x *node) (xform, bool) {
	if x == nil {
		return xform{}, false
	}
	off, ext := x.child("off"), x.child("ext")
	if off == nil && ext == nil {
		return xform{}, false
	}
	return xform{
		X: off.emuAttr("x", 0), Y: off.emuAttr("y", 0), W: ext.emuAttr("cx", 0), H: ext.emuAttr("cy", 0),
		Rot:   float64(x.attrInt("rot", 0)) / 60000,
		FlipH: x.attrBool("flipH", false), FlipV: x.attrBool("flipV", false),
	}, true
}

// matrix maps the shape's local box (0,0)-(W,H) to the slide.
func (x xform) matrix() matrix {
	fh, fv := 1.0, 1.0
	if x.FlipH {
		fh = -1
	}
	if x.FlipV {
		fv = -1
	}
	if x.Rot == 0 && !x.FlipH && !x.FlipV {
		return translate(x.X, x.Y)
	}
	return translate(x.X+x.W/2, x.Y+x.H/2).mul(rotate(x.Rot)).mul(scale(fh, fv)).mul(translate(-x.W/2, -x.H/2))
}

// textMatrix is matrix without mirroring: text is never drawn mirrored; a
// vertical flip turns it upside down.
func (x xform) textMatrix() matrix {
	rot := x.Rot
	if x.FlipV {
		rot += 180
	}
	if rot == 0 {
		return translate(x.X, x.Y)
	}
	return translate(x.X+x.W/2, x.Y+x.H/2).mul(rotate(rot)).mul(translate(-x.W/2, -x.H/2))
}

// bounds returns the axis-aligned bounding box on the slide.
func (x xform) bounds() (x0, y0, x1, y1 float64) {
	m := x.matrix()
	x0, y0 = math.Inf(1), math.Inf(1)
	x1, y1 = math.Inf(-1), math.Inf(-1)
	for _, p := range [][2]float64{{0, 0}, {x.W, 0}, {0, x.H}, {x.W, x.H}} {
		px, py := m.apply(p[0], p[1])
		x0, y0, x1, y1 = math.Min(x0, px), math.Min(y0, py), math.Max(x1, px), math.Max(y1, py)
	}
	return
}

// groupCtx maps the child coordinate space of a group onto the slide.
type groupCtx struct {
	parent *groupCtx
	xf     xform // the group's placement on the slide
	chOff  [2]float64
	chExt  [2]float64
	spPr   *node // grpSpPr, for grpFill
	part   string
}

// place converts a child placement into slide coordinates. Groups are
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
	n      *node
	part   string
	inh    []*node // layout placeholder, master placeholder (most specific first)
	inhPrt []string
	ph     *node
	grp    *groupCtx
}

func (sh *shape) spPr() *node {
	if p := sh.n.child("spPr"); p != nil {
		return p
	}
	return sh.n.child("grpSpPr")
}

// spPrs returns the shape properties of the shape and its placeholders.
func (sh *shape) spPrs() ([]*node, []string) {
	out := []*node{sh.spPr()}
	parts := []string{sh.part}
	for i, n := range sh.inh {
		out = append(out, n.child("spPr"))
		parts = append(parts, sh.inhPrt[i])
	}
	return out, parts
}

func (sh *shape) xform() (xform, bool) {
	cands := []*node{sh.spPr().child("xfrm"), sh.n.child("xfrm")}
	for _, n := range sh.inh {
		cands = append(cands, n.path("spPr", "xfrm"), n.child("xfrm"))
	}
	// Offset and size are taken separately, so that a shape that only
	// resizes its placeholder keeps the inherited position.
	var first, off, ext *node
	for _, x := range cands {
		if x == nil {
			continue
		}
		if first == nil {
			first = x
		}
		if off == nil {
			off = x.child("off")
		}
		if ext == nil {
			ext = x.child("ext")
		}
	}
	if off == nil && ext == nil {
		return xform{}, false
	}
	xf := xform{
		X: off.emuAttr("x", 0), Y: off.emuAttr("y", 0), W: ext.emuAttr("cx", 0), H: ext.emuAttr("cy", 0),
		Rot:   float64(first.attrInt("rot", 0)) / 60000,
		FlipH: first.attrBool("flipH", false), FlipV: first.attrBool("flipV", false),
	}
	return sh.grp.place(xf), true
}

func (sh *shape) style() *node {
	if st := sh.n.child("style"); st != nil {
		return st
	}
	for _, n := range sh.inh {
		if st := n.child("style"); st != nil {
			return st
		}
	}
	return nil
}

func (sh *shape) geometry(w, h float64) (*geometry, string) {
	sps, _ := sh.spPrs()
	for _, sp := range sps {
		if sp.child("prstGeom") != nil || sp.child("custGeom") != nil {
			return shapeGeometry(sp, w, h)
		}
	}
	return rectGeometry(w, h), ""
}

func (s *slideCtx) shapeFill(sh *shape) fill {
	sps, parts := sh.spPrs()
	for i, sp := range sps {
		if e := fillElem(sp); e != nil {
			return s.resolveFill(e, parts[i], s.cc, sh.grp)
		}
	}
	return s.styleFill(sh.style().child("fillRef"), s.cc)
}

func (s *slideCtx) shapeLine(sh *shape) *line {
	sps, parts := sh.spPrs()
	var lns []*node
	var lparts []string
	for i, sp := range sps {
		if ln := sp.child("ln"); ln != nil {
			lns = append(lns, ln)
			lparts = append(lparts, parts[i])
		}
	}
	return s.resolveLine(lns, lparts, sh.style().child("lnRef"), s.cc)
}

func (s *slideCtx) shapeShadow(sh *shape) *shadow {
	sps, _ := sh.spPrs()
	var effects []*node
	for _, sp := range sps {
		effects = append(effects, sp.child("effectLst"))
	}
	return s.resolveShadow(effects, sh.style().child("effectRef"), s.cc)
}

func cNvPr(n *node) *node {
	for _, k := range n.Kids {
		if strings.HasPrefix(k.Name, "nv") {
			if c := k.child("cNvPr"); c != nil {
				return c
			}
		}
	}
	return nil
}

// drawTree draws the shapes of a p:spTree. Placeholders are only drawn on
// slides: on masters and layouts they are prompts for the slides.
func (s *slideCtx) drawTree(cv *canvas, tree *node, part string, isSlide bool) {
	for _, k := range tree.kids() {
		s.drawElem(cv, k, part, isSlide, nil)
	}
}

func (s *slideCtx) drawElem(cv *canvas, k *node, part string, isSlide bool, grp *groupCtx) {
	switch k.Name {
	case "sp", "cxnSp", "pic", "graphicFrame", "grpSp":
	case "contentPart":
		s.c.warnOnce("ink", "ink (content parts) is not supported")
		return
	default:
		return
	}
	if cNvPr(k).attrBool("hidden", false) {
		return
	}
	sh := &shape{n: k, part: part, grp: grp, ph: phOf(k)}
	if sh.ph != nil {
		if !isSlide {
			return
		}
		mph := sh.ph
		if lp := findPh(s.layoutPh, sh.ph, false); lp != nil {
			sh.inh = append(sh.inh, lp)
			sh.inhPrt = append(sh.inhPrt, s.layoutPart)
			mph = phOf(lp)
		}
		if mp := findPh(s.masterPh, mph, true); mp != nil {
			sh.inh = append(sh.inh, mp)
			sh.inhPrt = append(sh.inhPrt, s.masterPart)
		}
	}
	switch k.Name {
	case "grpSp":
		s.drawGroup(cv, sh, isSlide)
	case "sp", "cxnSp":
		s.drawSp(cv, sh)
	case "pic":
		s.drawPic(cv, sh)
	case "graphicFrame":
		s.drawFrame(cv, sh)
	}
}

func (s *slideCtx) drawGroup(cv *canvas, sh *shape, isSlide bool) {
	gp := sh.n.child("grpSpPr")
	x := gp.child("xfrm")
	own, ok := parseXfrm(x)
	g := &groupCtx{parent: sh.grp, spPr: gp, part: sh.part}
	if ok {
		// The group box stays in its parent's coordinates: place() maps
		// children into it and then recurses into the parent group.
		g.xf = own
		chOff, chExt := x.child("chOff"), x.child("chExt")
		g.chOff = [2]float64{chOff.emuAttr("x", own.X), chOff.emuAttr("y", own.Y)}
		g.chExt = [2]float64{chExt.emuAttr("cx", own.W), chExt.emuAttr("cy", own.H)}
	} else {
		// no placement: children keep their coordinates
		g.xf = xform{W: 1, H: 1}
		g.chExt = [2]float64{1, 1}
	}
	for _, k := range sh.n.Kids {
		s.drawElem(cv, k, sh.part, isSlide, g)
	}
}

func (s *slideCtx) drawSp(cv *canvas, sh *shape) {
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
		cv.obj.Save()
		cv.transform(xf.matrix())
		s.paintGeometry(cv, geo, fl, ln, s.shapeShadow(sh), xf.W, xf.H)
		cv.obj.Restore()
		cv.drawn = true
	}
	s.drawShapeText(cv, sh, xf, geo)
	s.shapeLink(cv, sh, xf)
}

// paintGeometry fills and strokes the paths of a geometry in local coordinates.
func (s *slideCtx) paintGeometry(cv *canvas, geo *geometry, fl fill, ln *line, shd *shadow, w, h float64) {
	refs := make([]bdf.PathRef, len(geo.paths))
	for i, p := range geo.paths {
		refs[i] = cv.obj.AddPath(p.path)
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
		cv.obj.Save()
		cv.obj.Shadow(shd.c.bdf(), f32(shd.blur), f32(shd.dx), f32(shd.dy))
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
			cv.obj.Restore()
			shd = nil
		}
	}
	if ln != nil {
		cv.obj.Save()
		cv.setLine(ln, w, h)
		for i, p := range geo.paths {
			if !p.stroke || p.empty {
				continue
			}
			cv.obj.StrokePath(refs[i])
		}
		for _, p := range geo.paths {
			if p.stroke && !p.empty && !p.closed {
				cv.obj.Save()
				cv.obj.Dash(nil, 0)
				cv.arrowHead(ln, ln.head, p.start[0], p.start[1], -p.startDir[0], -p.startDir[1])
				cv.arrowHead(ln, ln.tail, p.end[0], p.end[1], p.endDir[0], p.endDir[1])
				cv.obj.Restore()
			}
		}
		cv.obj.Restore()
	}
	if shd != nil {
		cv.obj.Restore()
	}
}

func (s *slideCtx) fillPath(cv *canvas, ref bdf.PathRef, f fill, w, h float64) {
	if f.kind == fillBlip {
		cv.obj.Save()
		cv.obj.ClipPath(ref, 0)
		s.drawBlip(cv, f.blip, f.part, 0, 0, w, h, true)
		cv.obj.Restore()
		return
	}
	if cv.setFill(f, w, h) {
		cv.obj.FillPath(ref, 0)
	}
}

// drawBlip draws the image of an a:blipFill / p:blipFill into x,y,w,h
// (stretched, cropped by srcRect, or tiled).
func (s *slideCtx) drawBlip(cv *canvas, bf *node, part string, x, y, w, h float64, clipped bool) bool {
	blip := bf.child("blip")
	rid := blip.rid("embed")
	if rid == "" {
		rid = blip.rid("link")
	}
	e := s.c.image(part, rid)
	if !e.ok {
		return false
	}
	ref := cv.image(e.hash)
	cv.drawn = true
	if a := blip.child("alphaModFix"); a != nil {
		cv.obj.Alpha(f32(pct(a, "amt", 1)))
	}
	var filters []string
	if blip.child("grayscl") != nil {
		filters = append(filters, "grayscale(1)")
	}
	if l := blip.child("lum"); l != nil {
		if b := pct(l, "bright", 0); b != 0 {
			filters = append(filters, "brightness("+ftoa(1+b)+")")
		}
		if c := pct(l, "contrast", 0); c != 0 {
			filters = append(filters, "contrast("+ftoa(1+c)+")")
		}
	}
	if blip.child("duotone") != nil || blip.child("clrChange") != nil || blip.child("biLevel") != nil {
		s.c.warnOnce("blipfx", "picture recoloring effects are not supported")
	}
	if len(filters) > 0 {
		cv.obj.Filter(strings.Join(filters, " "))
	}
	if t := bf.child("tile"); t != nil && e.w > 0 {
		sx, sy := pct(t, "sx", 1), pct(t, "sy", 1)
		tw, th := e.w*0.75*sx, e.h*0.75*sy
		ox, oy := x+t.emuAttr("tx", 0), y+t.emuAttr("ty", 0)
		algn := t.attrStr("algn", "tl")
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
		cv.obj.FillPaint(cv.obj.AddPaint(p))
		if !clipped {
			cv.obj.FillRect(f32(x), f32(y), f32(w), f32(h))
		} else {
			cv.obj.FillRect(f32(x-1), f32(y-1), f32(w+2), f32(h+2))
		}
		return true
	}
	// stretch, inset by fillRect
	if fr := bf.path("stretch", "fillRect"); fr != nil {
		l, t, r, b := pct(fr, "l", 0), pct(fr, "t", 0), pct(fr, "r", 0), pct(fr, "b", 0)
		x, y, w, h = x+l*w, y+t*h, w*(1-l-r), h*(1-t-b)
	}
	sr := bf.child("srcRect")
	if sr == nil || e.w <= 0 || e.h <= 0 {
		cv.obj.Image(ref, f32(x), f32(y), f32(w), f32(h))
		return true
	}
	l, t, r, b := pct(sr, "l", 0), pct(sr, "t", 0), pct(sr, "r", 0), pct(sr, "b", 0)
	if l == 0 && t == 0 && r == 0 && b == 0 {
		cv.obj.Image(ref, f32(x), f32(y), f32(w), f32(h))
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
	cv.obj.ImageSub(ref, f32(cu0), f32(cv0), f32(cu1-cu0), f32(cv1-cv0), f32(dx0), f32(dy0), f32(dx1-dx0), f32(dy1-dy0))
	return true
}

func (s *slideCtx) drawPic(cv *canvas, sh *shape) {
	if xf, ok := sh.xform(); ok {
		s.drawPicAt(cv, sh, xf)
	}
}

func (s *slideCtx) drawPicAt(cv *canvas, sh *shape, xf xform) {
	geo, _ := sh.geometry(xf.W, xf.H)
	bf := sh.n.child("blipFill")
	if bf == nil {
		for _, n := range sh.inh {
			if bf = n.child("blipFill"); bf != nil {
				break
			}
		}
	}
	fl := s.shapeFill(sh)
	ln := s.shapeLine(sh)
	shd := s.shapeShadow(sh)
	cv.obj.Save()
	cv.transform(xf.matrix())
	if fl.kind != fillNone {
		s.paintGeometry(cv, geo, fl, nil, shd, xf.W, xf.H)
		shd = nil
	}
	if bf != nil {
		cv.obj.Save()
		rect := len(geo.paths) == 1 && len(geo.paths[0].path.Verbs) == 1 && geo.paths[0].path.Verbs[0] == bdf.VerbRect
		if shd != nil {
			cv.obj.Shadow(shd.c.bdf(), f32(shd.blur), f32(shd.dx), f32(shd.dy))
		}
		if !rect {
			for _, p := range geo.paths {
				if p.fill != "none" {
					cv.obj.ClipPath(cv.obj.AddPath(p.path), 0)
					break
				}
			}
		}
		s.drawBlip(cv, bf, sh.part, 0, 0, xf.W, xf.H, !rect)
		cv.obj.Restore()
	}
	if ln != nil {
		s.paintGeometry(cv, geo, fill{}, ln, nil, xf.W, xf.H)
	}
	cv.obj.Restore()
	cv.drawn = true
	s.shapeLink(cv, sh, xf)
}

// shapeLink emits the click hyperlink of a shape over its bounds.
func (s *slideCtx) shapeLink(cv *canvas, sh *shape, xf xform) {
	h := cNvPr(sh.n).child("hlinkClick")
	if h == nil {
		return
	}
	target := s.linkTarget(h, sh.part)
	if target == "" {
		return
	}
	x0, y0, x1, y1 := xf.bounds()
	cv.obj.Link(f32(x0), f32(y0), f32(x1-x0), f32(y1-y0), target)
}

// linkTarget resolves an a:hlinkClick to a URL or "#page=N".
func (s *slideCtx) linkTarget(h *node, part string) string {
	action := h.attrStr("action", "")
	if strings.HasPrefix(action, "ppaction://hlinkshowjump") {
		cur := s.c.pageOf[s.part]
		switch {
		case strings.Contains(action, "nextslide"):
			return pageLink(cur + 1)
		case strings.Contains(action, "previousslide"):
			return pageLink(cur - 1)
		case strings.Contains(action, "firstslide"):
			return pageLink(1)
		case strings.Contains(action, "lastslide"):
			return pageLink(s.c.pages)
		}
		return ""
	}
	r, ok := s.c.pkg.target(part, h.rid("id"))
	if !ok {
		return ""
	}
	if r.External {
		if u, err := url.Parse(r.Target); err == nil && (u.Scheme == "http" || u.Scheme == "https" || u.Scheme == "mailto") {
			return r.Target
		}
		return ""
	}
	if strings.HasSuffix(r.Type, "/slide") {
		return pageLink(s.c.pageOf[r.Target])
	}
	return ""
}

func pageLink(n int) string {
	if n < 1 {
		return ""
	}
	return "#page=" + itoa(n)
}

// background fills the page with the first background of slide, layout
// and master (white when none has one).
func (s *slideCtx) background(cv *canvas) {
	type src struct {
		n    *node
		part string
	}
	for _, sr := range []src{{s.slide, s.part}, {s.layout, s.layoutPart}, {s.master, s.masterPart}} {
		bg := sr.n.path("cSld", "bg")
		if bg == nil {
			continue
		}
		var f fill
		if pr := bg.child("bgPr"); pr != nil {
			f = s.resolveFill(fillElem(pr), sr.part, s.cc, nil)
		} else if ref := bg.child("bgRef"); ref != nil {
			f = s.styleFill(ref, s.cc)
		}
		s.fillPage(cv, f)
		return
	}
	c, _ := s.cc.scheme1("bg1")
	s.fillPage(cv, fill{kind: fillSolid, color: c})
}

func (s *slideCtx) fillPage(cv *canvas, f fill) {
	w, h := s.c.slideW, s.c.slideH
	switch f.kind {
	case fillNone:
		return
	case fillBlip:
		s.drawBlip(cv, f.blip, f.part, 0, 0, w, h, false)
	default:
		if cv.setFill(f, w, h) {
			cv.obj.FillRect(0, 0, f32(w), f32(h))
		}
	}
	cv.drawn = true
}

// drawFrame draws a graphic frame: a table, a chart, a SmartArt diagram
// (from its pre-rendered drawing) or an OLE object's preview picture.
func (s *slideCtx) drawFrame(cv *canvas, sh *shape) {
	xf, ok := sh.xform()
	if !ok {
		return
	}
	gd := sh.n.path("graphic", "graphicData")
	uri := gd.attrStr("uri", "")
	switch {
	case gd.child("tbl") != nil:
		s.drawTable(cv, sh, xf, gd.child("tbl"))
	case strings.HasSuffix(uri, "/chart"):
		s.drawChart(cv, sh, xf, gd.child("chart").rid("id"))
	case strings.HasSuffix(uri, "/diagram"):
		s.drawDiagram(cv, sh, xf, gd.child("relIds"))
	default:
		// OLE objects and other embeddings: draw their preview picture.
		var pic *node
		var find func(n *node)
		find = func(n *node) {
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
		if pxf, ok := ps.xform(); ok {
			s.drawPicAt(cv, ps, pxf)
		} else {
			s.drawPicAt(cv, ps, xf) // the preview fills the frame
		}
	}
}

// drawDiagram draws a SmartArt diagram from the drawing part PowerPoint
// stores next to its data model (the shapes as last laid out).
func (s *slideCtx) drawDiagram(cv *canvas, sh *shape, xf xform, relIds *node) {
	var drawingPart string
	if r, ok := s.c.pkg.target(sh.part, relIds.rid("dm")); ok {
		if dm, err := s.c.pkg.xml(r.Target); err == nil {
			for _, ext := range dm.path("extLst").children("ext") {
				if d := ext.child("dataModelExt"); d != nil {
					if dr, ok := s.c.pkg.target(sh.part, d.attrStr("relId", "")); ok {
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
	d, err := s.c.pkg.xml(drawingPart)
	if err != nil {
		s.c.warnf("SmartArt drawing: %v", err)
		return
	}
	// xf is already on the slide; the drawing's coordinates are relative to
	// the frame.
	g := &groupCtx{xf: xform{X: xf.X, Y: xf.Y, W: xf.W, H: xf.H}, chExt: [2]float64{xf.W, xf.H}, part: drawingPart}
	for _, k := range d.path("spTree").kids() {
		s.drawElem(cv, k, drawingPart, false, g)
	}
}
