package drawingml

import (
	"encoding/xml"
	"math"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// drawShapeText lays out and draws the text body of a shape in its text
// rectangle (or, for SmartArt drawings, the separate text transform).
func (s *Drawing) drawShapeText(cv *canvas.Canvas, sh *shape, xf xform, geo *geometry) {
	if sh.n.Child("txbx") != nil {
		s.drawTextBox(cv, sh, xf, geo)
		return
	}
	tb := sh.n.Child("txBody")
	if tb == nil {
		return
	}
	tf := &textFrame{body: tb, part: sh.part, cc: s.cc, own: tb.Child("lstStyle"), bodyPrs: chain{tb.Child("bodyPr")}}
	phType := ""
	if sh.ph != nil {
		if t, ok := sh.ph.Attr("type"); ok {
			phType = t
		}
	}
	for _, n := range sh.inh {
		tf.bodyPrs = append(tf.bodyPrs, n.Path("txBody", "bodyPr"))
		tf.lists = append(tf.lists, n.Path("txBody", "lstStyle"))
		if phType == "" {
			if t, ok := PlaceholderOf(n).Attr("type"); ok {
				phType = t
			}
		}
	}
	if sh.ph != nil {
		if phType == "" {
			phType = "obj"
		}
		tf.phType = phType
	}
	tf.lists = append(tf.lists, s.host.TextStyle(phType))
	if fr := sh.style().Child("fontRef"); fr != nil {
		tf.extra = append(tf.extra, fontRefProps(fr))
	}
	r := geo.textRect
	x, y, w, h := r[0], r[1], r[2]-r[0], r[3]-r[1]
	m := xf.textMatrix()
	if tx, ok := parseXfrm(sh.n.Child("txXfrm")); ok {
		// the text transform of a SmartArt drawing turns with its shape
		if own, ok := parseXfrm(sh.spPr().Child("xfrm")); ok {
			tx.Rot += own.Rot
		}
		tx = sh.grp.place(tx)
		m = tx.textMatrix()
		x, y, w, h = 0, 0, tx.W, tx.H
	}
	s.renderText(cv, tf, x, y, w, h, m)
}

// fontRefProps turns a p:style fontRef into run properties.
func fontRefProps(fr *ooxml.Node) *ooxml.Node {
	n := &ooxml.Node{Name: "defRPr"}
	switch fr.AttrStr("idx", "minor") {
	case "major":
		n.Kids = append(n.Kids, mkNode("latin", "typeface", "+mj-lt"), mkNode("ea", "typeface", "+mj-ea"), mkNode("cs", "typeface", "+mj-cs"))
	case "minor":
		n.Kids = append(n.Kids, mkNode("latin", "typeface", "+mn-lt"), mkNode("ea", "typeface", "+mn-ea"), mkNode("cs", "typeface", "+mn-cs"))
	}
	for _, k := range fr.Kids {
		if isColor(k) {
			n.Kids = append(n.Kids, &ooxml.Node{Name: "solidFill", Kids: []*ooxml.Node{k}})
			break
		}
	}
	return n
}

func mkNode(name string, kv ...string) *ooxml.Node {
	n := &ooxml.Node{Name: name}
	for i := 0; i+1 < len(kv); i += 2 {
		n.Attrs = append(n.Attrs, xml.Attr{Name: xml.Name{Local: kv[i]}, Value: kv[i+1]})
	}
	return n
}

func emuChain(c chain, name string, def float64) float64 {
	for _, n := range c {
		if _, ok := n.Attr(name); ok {
			return n.AttrEMU(name, def)
		}
	}
	return def
}

// textBlock is a laid-out text frame.
type textBlock struct {
	lo         *laidOut
	x, y, w, h float64 // box in frame coordinates
	ins        [4]float64
	anchor     string
	fm         canvas.Matrix
	upright    bool
	vertical   bool
	cols       int     // number of columns (0 or 1: one)
	colW       float64 // column width
	colGap     float64 // space between columns
	heading    string  // heading level of its paragraphs ("" for none)
	// inner maps the box (after rotation and vertical text) to the
	// coordinate space of the layout box, for Bounds
	inner canvas.Matrix
}

// renderText lays out a text frame in the box x,y,w,h of the coordinate
// space m and draws it.
func (s *Drawing) renderText(cv *canvas.Canvas, tf *textFrame, x, y, w, h float64, m canvas.Matrix) {
	if b := s.layoutText(tf, x, y, w, h, m); b != nil {
		s.drawTextBlock(cv, b)
	}
}

// layoutText lays out a text frame; nil when it has no text.
func (s *Drawing) layoutText(tf *textFrame, x, y, w, h float64, m canvas.Matrix) *textBlock {
	bp := tf.bodyPrs
	fontScale, lnReduce := 1.0, 0.0
	if af := bp.first("normAutofit", "noAutofit", "spAutoFit"); af != nil && af.Name == "normAutofit" {
		fontScale = af.AttrPct("fontScale", 1)
		lnReduce = af.AttrPct("lnSpcReduction", 0)
	}
	paras := s.paragraphs(tf, fontScale, lnReduce)
	if !anyText(paras) {
		return nil
	}
	return s.layoutBlock(tf, paras, x, y, w, h, m, 0)
}

func anyText(paras []*para) bool {
	for _, p := range paras {
		if p.hasText {
			return true
		}
	}
	return false
}

// layoutBlock lays paragraphs out in the box of a text frame (see
// layoutBodyLimit for limit).
func (s *Drawing) layoutBlock(tf *textFrame, paras []*para, x, y, w, h float64, m canvas.Matrix, limit float64) *textBlock {
	bp := tf.bodyPrs
	b := &textBlock{}
	if tf.phType == "title" || tf.phType == "ctrTitle" {
		b.heading = "1"
	}
	b.ins = [4]float64{emuChain(bp, "lIns", 7.2), emuChain(bp, "tIns", 3.6), emuChain(bp, "rIns", 7.2), emuChain(bp, "bIns", 3.6)}
	b.anchor, _ = bp.attr("anchor")
	wrapV, _ := bp.attr("wrap")
	wrap := wrapV != "none"
	vert, _ := bp.attr("vert")
	// fm maps the (turned) box to the page; inner is the same without m,
	// kept apart rather than factored out of fm so that fm rounds as before
	fm, inner := m, canvas.Identity
	if v, ok := bp.attr("rot"); ok {
		if rot := atof(v, 0) / 60000; rot != 0 {
			cx, cy := x+w/2, y+h/2
			fm = fm.Mul(canvas.Translate(cx, cy)).Mul(canvas.Rotate(rot)).Mul(canvas.Translate(-cx, -cy))
			inner = inner.Mul(canvas.Translate(cx, cy)).Mul(canvas.Rotate(rot)).Mul(canvas.Translate(-cx, -cy))
		}
	}
	switch vert {
	case "vert", "eaVert", "wordArtVert", "wordArtVertRtl", "mongolianVert":
		fm = fm.Mul(canvas.Translate(x+w, y)).Mul(canvas.Rotate(90))
		inner = inner.Mul(canvas.Translate(x+w, y)).Mul(canvas.Rotate(90))
		x, y, w, h = 0, 0, h, w
		b.upright = vert != "vert"
		b.vertical = true
	case "vert270":
		fm = fm.Mul(canvas.Translate(x, y+h)).Mul(canvas.Rotate(-90))
		inner = inner.Mul(canvas.Translate(x, y+h)).Mul(canvas.Rotate(-90))
		x, y, w, h = 0, 0, h, w
		b.vertical = true
	}
	b.x, b.y, b.w, b.h, b.fm, b.inner = x, y, w, h, fm, inner
	width := w - b.ins[0] - b.ins[2]
	if v, ok := bp.attr("numCol"); ok && atof(v, 1) > 1 && wrap {
		b.cols = min(int(atof(v, 1)), 16)
		b.colGap = emuChain(bp, "spcCol", 0)
		b.colW = (width - float64(b.cols-1)*b.colGap) / float64(b.cols)
		if b.colW > 0 {
			width = b.colW
		} else {
			b.cols = 0
		}
	}
	if limit > 0 && b.cols > 1 {
		limit *= float64(b.cols)
	}
	b.lo = layoutBodyLimit(paras, width, wrap, limit)
	return b
}

// height is the height the text needs, insets included.
func (b *textBlock) height() float64 { return b.lo.height + b.ins[1] + b.ins[3] }

// top returns where the first line starts: the text is anchored to the
// top, the middle or the bottom of the box inside its insets.
func (b *textBlock) top() float64 {
	iy := b.y + b.ins[1]
	ih := b.h - b.ins[1] - b.ins[3]
	switch b.anchor {
	case "ctr":
		return iy + (ih-b.lo.height)/2
	case "b":
		return iy + ih - b.lo.height
	}
	return iy
}

func (s *Drawing) drawTextBlock(cv *canvas.Canvas, b *textBlock) {
	ix, iy := b.x+b.ins[0], b.y+b.ins[1]
	ih := b.h - b.ins[1] - b.ins[3]
	dy := b.top()
	if b.cols > 1 {
		// Lines flow into the next column when they pass the bottom.
		dy = iy
		col, colTop, inCol := 0, 0.0, 0
		for _, ln := range b.lo.lines {
			top := ln.baseline + ln.desc - ln.height
			if ln.baseline+ln.desc-colTop > ih && inCol > 0 && col < b.cols-1 {
				col++
				colTop, inCol = top, 0
			}
			ln.cx, ln.cy = float64(col)*(b.colW+b.colGap), -colTop
			inCol++
		}
	}
	cv.Obj.Save()
	cv.Transform(b.fm)
	cv.Obj.Mark(bdf.MarkBox, "")
	em := &textEmitter{c: s.c, cv: cv, m: b.fm, upright: b.upright, heading: b.heading}
	em.emitLines(b.lo.lines, ix, dy)
	em.closeLists(-1)
	cv.Obj.Restore()
	for _, l := range em.links {
		cv.Obj.Link(f32(l.x0), f32(l.y0), f32(l.x1-l.x0), f32(l.y1-l.y0), l.url)
	}
}

// TextBody is a text body laid out in a box, ready to be drawn.
type TextBody struct {
	s *Drawing
	b *textBlock
}

// LayoutText lays out a text body (an element with the children of an
// a:txBody: a:bodyPr, a:lstStyle and a:p) in the box x,y,w,h of the
// coordinate space m, as the text of a shape that is not a placeholder: the
// host's TextStyle("") applies last. part is the part the body is in, for
// its relationships. Documents whose drawings are not DrawingML (Visio)
// describe their text this way to share the layout. It returns nil when the
// body has no text.
func (s *Drawing) LayoutText(body *ooxml.Node, part string, x, y, w, h float64, m canvas.Matrix) *TextBody {
	tf := &textFrame{body: body, part: part, cc: s.cc, own: body.Child("lstStyle"), bodyPrs: chain{body.Child("bodyPr")}}
	tf.lists = append(tf.lists, s.host.TextStyle(""))
	b := s.layoutText(tf, x, y, w, h, m)
	if b == nil {
		return nil
	}
	return &TextBody{s: s, b: b}
}

// Bounds returns the box the lines of text take with the insets around
// them, in the coordinate space of the layout box (before m).
func (t *TextBody) Bounds() (x0, y0, x1, y1 float64) {
	b := t.b
	lx0, lx1 := math.Inf(1), math.Inf(-1)
	for _, ln := range b.lo.lines {
		start := ln.start + ln.offset
		if bu := ln.pa.bullet; bu != nil && ln.first {
			start = min(start, ln.bulletX)
		}
		lx0, lx1 = min(lx0, start), max(lx1, ln.start+ln.offset+ln.width)
	}
	if lx0 > lx1 {
		lx0, lx1 = 0, 0
	}
	ix := b.x + b.ins[0]
	top := b.top()
	x0, y0 = math.Inf(1), math.Inf(1)
	x1, y1 = math.Inf(-1), math.Inf(-1)
	for _, p := range [][2]float64{
		{ix + lx0 - b.ins[0], top - b.ins[1]}, {ix + lx1 + b.ins[2], top - b.ins[1]},
		{ix + lx0 - b.ins[0], top + b.lo.height + b.ins[3]}, {ix + lx1 + b.ins[2], top + b.lo.height + b.ins[3]},
	} {
		px, py := b.inner.Apply(p[0], p[1])
		x0, y0, x1, y1 = min(x0, px), min(y0, py), max(x1, px), max(y1, py)
	}
	return x0, y0, x1, y1
}

// Draw draws the text.
func (t *TextBody) Draw(cv *canvas.Canvas) { t.s.drawTextBlock(cv, t.b) }

// TextFlow is a text body laid out over several boxes, as the story of a
// page layout program flows through its threaded text frames: each Fill
// takes the lines that fit its box, and the next continues with the rest.
type TextFlow struct {
	s     *Drawing
	tf    *textFrame
	paras []*para // what is left; the first may be the tail of a paragraph
}

// NewTextFlow starts the flow of a text body (an element with the children
// of an a:txBody, whose a:bodyPr the boxes override) of part; see
// LayoutText. It is nil when the body has no text.
func (s *Drawing) NewTextFlow(body *ooxml.Node, part string) *TextFlow {
	tf := &textFrame{body: body, part: part, cc: s.cc, own: body.Child("lstStyle"), bodyPrs: chain{body.Child("bodyPr")}}
	tf.lists = append(tf.lists, s.host.TextStyle(""))
	paras := s.paragraphs(tf, 1, 0)
	if !anyText(paras) {
		return nil
	}
	return &TextFlow{s: s, tf: tf, paras: paras}
}

// Done reports that no text is left to lay out.
func (f *TextFlow) Done() bool { return f == nil || len(f.paras) == 0 }

// Fill lays out as much of the text left as fits the box x,y,w,h of the
// coordinate space m, whose insets, anchor, columns and direction bodyPr
// gives (an a:bodyPr element), and keeps the rest for the next box. A box
// too small for a line still takes one. It returns nil when no text is
// left.
func (f *TextFlow) Fill(bodyPr *ooxml.Node, x, y, w, h float64, m canvas.Matrix) *TextBody {
	if f.Done() {
		return nil
	}
	tf := *f.tf
	tf.bodyPrs = chain{bodyPr}
	ih := h - emuChain(tf.bodyPrs, "tIns", 3.6) - emuChain(tf.bodyPrs, "bIns", 3.6)
	if v, _ := tf.bodyPrs.attr("vert"); v != "" && v != "horz" {
		ih = w - emuChain(tf.bodyPrs, "tIns", 3.6) - emuChain(tf.bodyPrs, "bIns", 3.6)
	}
	b := f.s.layoutBlock(&tf, f.paras, x, y, w, h, m, math.Max(ih, 1))
	lo := b.lo
	// the lines that fit, flowing into the next column as drawTextBlock does
	cols := max(b.cols, 1)
	ih = b.h - b.ins[1] - b.ins[3]
	n, col, colTop, inCol := 0, 0, 0.0, 0
	for _, ln := range lo.lines {
		if ln.baseline+ln.desc-colTop > ih && inCol > 0 {
			if col >= cols-1 {
				break
			}
			col++
			colTop, inCol = ln.baseline+ln.desc-ln.height, 0
		}
		inCol++
		n++
	}
	if n == 0 {
		// nothing fits: the box is empty, the text waits for the next one
		return nil
	}
	if n == len(lo.lines) {
		// every line laid out fits: paragraphs are laid out whole, so the
		// text goes on (when it does) with the next paragraph, which the
		// layout stopped before (the space after the last one passed the
		// limit)
		last := lo.lines[n-1]
		if last.pi >= len(f.paras)-1 {
			f.paras = nil
		} else {
			f.paras = f.paras[last.pi+1:]
		}
		return &TextBody{s: f.s, b: b}
	}
	cut := lo.lines[n]
	lo.lines = lo.lines[:n]
	last := lo.lines[n-1]
	lo.height = last.baseline + last.desc
	rest := f.paras[cut.pi:]
	if !cut.first {
		// the paragraph goes on in the next box, from the line that did
		// not fit: without its first-line indent, bullet or space before
		cp := *rest[0]
		cp.items = rest[0].items[cut.from:]
		cp.indent, cp.bullet, cp.befPct, cp.befPts = 0, nil, 0, 0
		rest = append([]*para{&cp}, rest[1:]...)
	}
	f.paras = rest
	return &TextBody{s: f.s, b: b}
}
