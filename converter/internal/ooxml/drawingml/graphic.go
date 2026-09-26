package drawingml

import (
	"encoding/xml"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// What word processing documents (WordprocessingML) need: their drawings
// are a:graphic elements placed by the text they are anchored in (the
// wp:inline and wp:anchor around them), their shapes hold text boxes of
// WordprocessingML paragraphs, and their text takes theme colors and fonts.

// Box is a rectangle on the page, in points.
type Box struct{ X, Y, W, H float64 }

// TextBoxHost is implemented by hosts whose shapes hold text in their own
// markup: the text boxes of word processing documents (wps:txbx, holding
// w:txbxContent). The host lays the text out and draws it.
type TextBoxHost interface {
	// DrawTextBox draws the content of a text box (the element inside
	// wps:txbx) of part.
	DrawTextBox(cv *canvas.Canvas, content *ooxml.Node, part string, box TextBox)
}

// TextBox is where the text of a shape goes.
type TextBox struct {
	// M maps the text area, whose top left corner is (0,0), onto the page.
	// It turns with the shape and, for vertical text, by 90°.
	M canvas.Matrix
	// W and H are the size of the text area: the shape's text rectangle
	// less the insets (W is the line length, also for vertical text).
	W, H float64
	// Anchor is the vertical alignment of the text: "t", "ctr" or "b".
	Anchor string
	// Wrap is false when lines do not wrap at the width (wrap="none").
	Wrap bool
	// Vert is the text direction of the shape's body properties ("horz",
	// "vert", "vert270", "eaVert" …).
	Vert string
}

// DrawGraphic draws the content of an a:graphic element of part placed in
// box: a picture, a shape (with its text box), a group, a drawing canvas,
// a chart, a SmartArt diagram, or the preview picture of another
// embedding. The content turns and mirrors as its own transform says. It
// returns false when there is nothing it can draw.
func (s *Drawing) DrawGraphic(cv *canvas.Canvas, graphic *ooxml.Node, part string, box Box) bool {
	gd := graphic.Child("graphicData")
	uri := gd.AttrStr("uri", "")
	var el *ooxml.Node
	for _, k := range gd.Elements() {
		el = k
		break
	}
	if el == nil {
		return false
	}
	own, _ := parseXfrm(el.Child("spPr").Child("xfrm"))
	if own.W == 0 && own.H == 0 {
		own, _ = parseXfrm(el.Child("grpSpPr").Child("xfrm"))
	}
	xf := xform{X: box.X, Y: box.Y, W: box.W, H: box.H, Rot: own.Rot, FlipH: own.FlipH, FlipV: own.FlipV}
	sh := &shape{n: el, part: part, box: &xf}
	switch {
	case el.Name == "pic":
		s.drawPicAt(cv, sh, xf)
	case el.Name == "wsp":
		s.drawSp(cv, sh)
	case el.Name == "wgp" || el.Name == "grpSp":
		s.drawGroupAt(cv, sh, xf)
	case el.Name == "wpc":
		// a drawing canvas: its shapes are placed relative to it
		if bg := el.Child("bg"); bg != nil {
			if e := fillElem(bg); e != nil {
				cv.Obj.Save()
				cv.Transform(xf.matrix())
				s.fillPage(cv, s.resolveFill(e, part, s.cc, nil), xf.W, xf.H)
				cv.Obj.Restore()
			}
		}
		g := &groupCtx{xf: xf, chExt: [2]float64{xf.W, xf.H}, part: part}
		for _, k := range el.Kids {
			s.drawElem(cv, k, part, g)
		}
	case strings.HasSuffix(uri, "/chart") || el.Name == "chart":
		s.drawChart(cv, sh, xf, el.RelID("id"))
	case strings.HasSuffix(uri, "/diagram") || el.Name == "relIds":
		s.drawDiagram(cv, sh, xf, el)
	default:
		pic := findPic(gd)
		if pic == nil {
			s.c.warnOnce("graphic:"+uri, "graphic content %s is not supported", uri)
			return false
		}
		s.drawPicAt(cv, &shape{n: pic, part: part, box: &xf}, xf)
	}
	return true
}

// findPic returns the first picture (pic:pic) under n.
func findPic(n *ooxml.Node) *ooxml.Node {
	for _, k := range n.Elements() {
		if k.Name == "pic" {
			return k
		}
		if p := findPic(k); p != nil {
			return p
		}
	}
	return nil
}

// drawGroupAt draws a group whose own box is xf (instead of its xfrm).
func (s *Drawing) drawGroupAt(cv *canvas.Canvas, sh *shape, xf xform) {
	gp := sh.n.Child("grpSpPr")
	x := gp.Child("xfrm")
	g := &groupCtx{parent: sh.grp, xf: xf, spPr: gp, part: sh.part}
	chOff, chExt := x.Child("chOff"), x.Child("chExt")
	ext := x.Child("ext")
	g.chOff = [2]float64{chOff.AttrEMU("x", 0), chOff.AttrEMU("y", 0)}
	g.chExt = [2]float64{chExt.AttrEMU("cx", ext.AttrEMU("cx", xf.W)), chExt.AttrEMU("cy", ext.AttrEMU("cy", xf.H))}
	for _, k := range sh.n.Kids {
		s.drawElem(cv, k, sh.part, g)
	}
}

// drawTextBox hands the text box of a word processing shape (wps:txbx) to
// the host, in the shape's text rectangle less the insets of its body
// properties (wps:bodyPr).
func (s *Drawing) drawTextBox(cv *canvas.Canvas, sh *shape, xf xform, geo *geometry) {
	h, ok := s.host.(TextBoxHost)
	if !ok {
		return
	}
	var content *ooxml.Node
	for _, k := range sh.n.Child("txbx").Elements() {
		content = k
		break
	}
	if content == nil {
		return
	}
	bp := sh.n.Child("bodyPr")
	r := geo.textRect
	x, y, w, ht := r[0], r[1], r[2]-r[0], r[3]-r[1]
	l, t := bp.AttrEMU("lIns", 7.2), bp.AttrEMU("tIns", 3.6)
	x, y = x+l, y+t
	w -= l + bp.AttrEMU("rIns", 7.2)
	ht -= t + bp.AttrEMU("bIns", 3.6)
	m := xf.textMatrix()
	vert := bp.AttrStr("vert", "horz")
	switch vert {
	case "vert", "eaVert", "wordArtVert", "wordArtVertRtl", "mongolianVert":
		m = m.Mul(canvas.Translate(x+w, y)).Mul(canvas.Rotate(90))
		w, ht = ht, w
	case "vert270":
		m = m.Mul(canvas.Translate(x, y+ht)).Mul(canvas.Rotate(-90))
		w, ht = ht, w
	default:
		m = m.Mul(canvas.Translate(x, y))
	}
	h.DrawTextBox(cv, content, sh.part, TextBox{M: m, W: max(w, 0), H: max(ht, 0),
		Anchor: bp.AttrStr("anchor", "t"), Wrap: bp.AttrStr("wrap", "square") != "none", Vert: vert})
}

// DrawImage draws the picture a relationship of part points to (a raster
// image or a Windows metafile) stretched over box. It returns false when
// the picture cannot be read.
func (s *Drawing) DrawImage(cv *canvas.Canvas, part, rid string, box Box) bool {
	blip := &ooxml.Node{Name: "blip", Attrs: []xml.Attr{{Name: xml.Name{Space: relNS, Local: "embed"}, Value: rid}}}
	bf := &ooxml.Node{Name: "blipFill", Kids: []*ooxml.Node{blip}}
	return s.drawBlip(cv, bf, part, s.cc, box.X, box.Y, box.W, box.H, false)
}

const relNS = "http://schemas.openxmlformats.org/officeDocument/2006/relationships"

// ThemeColor returns a color of the theme part relates to by its scheme
// name (dk1, lt1, dk2, lt2, accent1 … accent6, hlink, folHlink).
func (c *Renderer) ThemeColor(part, name string) (bdf.Color, bool) {
	col, ok := c.theme(part).colors[name]
	return col.bdf(), ok
}

// ThemeFont resolves a theme font reference (+mn-lt, +mj-ea, +mn-cs …) with
// the theme part relates to. script is the script tag (Jpan, Hans …)
// whose font stands in when the theme leaves the East Asian or complex
// script font empty.
func (c *Renderer) ThemeFont(part, typeface, script string) string {
	return c.theme(part).fontFor(typeface, script)
}
