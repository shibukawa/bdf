package drawingml

import (
	"encoding/xml"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// drawShapeText lays out and draws the text body of a shape in its text
// rectangle (or, for SmartArt drawings, the separate text transform).
func (s *Drawing) drawShapeText(cv *canvas.Canvas, sh *shape, xf xform, geo *geometry) {
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
	hasText := false
	for _, p := range paras {
		if p.hasText {
			hasText = true
			break
		}
	}
	if !hasText {
		return nil
	}
	b := &textBlock{}
	if tf.phType == "title" || tf.phType == "ctrTitle" {
		b.heading = "1"
	}
	b.ins = [4]float64{emuChain(bp, "lIns", 7.2), emuChain(bp, "tIns", 3.6), emuChain(bp, "rIns", 7.2), emuChain(bp, "bIns", 3.6)}
	b.anchor, _ = bp.attr("anchor")
	wrapV, _ := bp.attr("wrap")
	wrap := wrapV != "none"
	vert, _ := bp.attr("vert")
	fm := m
	if v, ok := bp.attr("rot"); ok {
		if rot := atof(v, 0) / 60000; rot != 0 {
			cx, cy := x+w/2, y+h/2
			fm = fm.Mul(canvas.Translate(cx, cy)).Mul(canvas.Rotate(rot)).Mul(canvas.Translate(-cx, -cy))
		}
	}
	switch vert {
	case "vert", "eaVert", "wordArtVert", "wordArtVertRtl", "mongolianVert":
		fm = fm.Mul(canvas.Translate(x+w, y)).Mul(canvas.Rotate(90))
		x, y, w, h = 0, 0, h, w
		b.upright = vert != "vert"
		b.vertical = true
	case "vert270":
		fm = fm.Mul(canvas.Translate(x, y+h)).Mul(canvas.Rotate(-90))
		x, y, w, h = 0, 0, h, w
		b.vertical = true
	}
	b.x, b.y, b.w, b.h, b.fm = x, y, w, h, fm
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
	b.lo = layoutBody(paras, width, wrap)
	return b
}

// height is the height the text needs, insets included.
func (b *textBlock) height() float64 { return b.lo.height + b.ins[1] + b.ins[3] }

func (s *Drawing) drawTextBlock(cv *canvas.Canvas, b *textBlock) {
	ix, iy := b.x+b.ins[0], b.y+b.ins[1]
	ih := b.h - b.ins[1] - b.ins[3]
	dy := iy
	switch b.anchor {
	case "ctr":
		dy = iy + (ih-b.lo.height)/2
	case "b":
		dy = iy + ih - b.lo.height
	}
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
