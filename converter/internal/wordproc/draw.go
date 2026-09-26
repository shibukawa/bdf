package wordproc

import (
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
)

// inlineObj is a drawing: a DrawingML graphic (picture, shape, group,
// chart, SmartArt) or a VML picture.
type inlineObj struct {
	graphic *ooxml.Node // a:graphic
	vmlRID  string      // VML picture: the image relationship
	part    string
	w, h    float64
	ext     [4]float64 // effect extent (left, top, right, bottom): room for shadows and glow
	alt     string     // alternative text ("" for none or decorative)
	link    string

	// HTML: paint draws the object into its box (instead of a graphic), and
	// fit shrinks it to the width of the line when it is wider
	paint func(e *emitter, box drawingml.Box)
	fit   bool
}

// floatObj is a drawing anchored in a paragraph and positioned on the page
// (wp:anchor).
type floatObj struct {
	inlineObj
	hRel, vRel     string // relativeFrom
	hAlign, vAlign string
	hOff, vOff     float64
	behind         bool
	wrap           string // none, square, tight, through, topAndBottom
	side           string // bothSides, left, right, largest
	dist           [4]float64
	z              int64
}

// drawing reads a w:drawing: an inline drawing becomes an item of the
// line, an anchored one a floating object of the paragraph.
func (w *walker) drawing(p *para, n *ooxml.Node, st *runStyle) {
	for _, k := range n.Elements() {
		switch k.Name {
		case "inline":
			o := w.drawingObj(k)
			if o == nil {
				continue
			}
			p.items = append(p.items, item{kind: kObject, obj: o, st: st, w: o.w + o.ext[0] + o.ext[2]})
		case "anchor":
			o := w.drawingObj(k)
			if o == nil {
				continue
			}
			f := &floatObj{inlineObj: *o, behind: k.AttrBool("behindDoc", false), z: k.AttrInt("relativeHeight", 0),
				dist: [4]float64{emu(k, "distT"), emu(k, "distB"), emu(k, "distL"), emu(k, "distR")}}
			ph, pv := k.Child("positionH"), k.Child("positionV")
			f.hRel, f.vRel = ph.AttrStr("relativeFrom", "column"), pv.AttrStr("relativeFrom", "paragraph")
			f.hAlign, f.vAlign = ph.Child("align").Content(), pv.Child("align").Content()
			f.hOff, f.vOff = emuText(ph.Child("posOffset")), emuText(pv.Child("posOffset"))
			if k.AttrBool("simplePos", false) {
				sp := k.Child("simplePos")
				f.hRel, f.vRel, f.hAlign, f.vAlign = "page", "page", "", ""
				f.hOff, f.vOff = sp.AttrEMU("x", 0), sp.AttrEMU("y", 0)
			}
			f.wrap = "none"
			for _, wk := range k.Elements() {
				switch wk.Name {
				case "wrapSquare", "wrapTight", "wrapThrough":
					f.wrap = "square"
					f.side = wk.AttrStr("wrapText", "bothSides")
				case "wrapTopAndBottom":
					f.wrap = "topAndBottom"
				}
			}
			p.anchors = append(p.anchors, f)
		}
	}
}

func emu(n *ooxml.Node, name string) float64 { return n.AttrEMU(name, 0) }

func emuText(n *ooxml.Node) float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(n.Content()), 64)
	if err != nil {
		return 0
	}
	return v / ooxml.EMUPerPoint
}

// drawingObj reads what an inline and an anchored drawing share.
func (w *walker) drawingObj(k *ooxml.Node) *inlineObj {
	ext := k.Child("extent")
	o := &inlineObj{graphic: k.Child("graphic"), part: w.part, w: ext.AttrEMU("cx", 0), h: ext.AttrEMU("cy", 0)}
	if ee := k.Child("effectExtent"); ee != nil {
		o.ext = [4]float64{emu(ee, "l"), emu(ee, "t"), emu(ee, "r"), emu(ee, "b")}
		for i := range o.ext {
			o.ext[i] = max(o.ext[i], 0)
		}
	}
	if o.graphic == nil || o.w <= 0 || o.h <= 0 {
		return nil
	}
	doc := k.Child("docPr")
	if doc.AttrBool("hidden", false) {
		return nil
	}
	decorative := false
	for _, e := range doc.Path("extLst").Children("ext") {
		if e.Child("decorative").AttrBool("val", false) {
			decorative = true
		}
	}
	if !decorative {
		for _, name := range []string{"descr", "title"} {
			if v := strings.TrimSpace(doc.AttrStr(name, "")); v != "" {
				o.alt = v
				break
			}
		}
	}
	if h := doc.Child("hlinkClick"); h != nil {
		if r, ok := w.c.pkg.Target(w.part, h.RelID("id")); ok && r.External {
			o.link = safeURL(r.Target)
		}
	}
	return o
}

// vml reads a legacy VML picture (w:pict, w:object): its image, inline.
func (w *walker) vml(p *para, n *ooxml.Node, st *runStyle) {
	var find func(n *ooxml.Node) (*ooxml.Node, *ooxml.Node)
	find = func(n *ooxml.Node) (*ooxml.Node, *ooxml.Node) {
		for _, k := range n.Elements() {
			if k.Name == "imagedata" {
				return n, k
			}
			if s, d := find(k); d != nil {
				return s, d
			}
		}
		return nil, nil
	}
	shape, data := find(n)
	if data == nil {
		if hasElem(n, "textbox") {
			w.c.warnOnce("vmltext", "legacy VML text boxes are not drawn")
		}
		return
	}
	rid := data.RelID("id")
	if rid == "" {
		rid = data.RelID("pict")
	}
	css := cssProps(shape.AttrStr("style", ""))
	o := &inlineObj{vmlRID: rid, part: w.part, w: cssLength(css["width"]), h: cssLength(css["height"]), alt: strings.TrimSpace(shape.AttrStr("alt", ""))}
	if o.w <= 0 || o.h <= 0 || rid == "" {
		return
	}
	if css["position"] == "absolute" {
		f := &floatObj{inlineObj: *o, wrap: "none", hRel: "column", vRel: "paragraph",
			hOff: cssLength(css["margin-left"]) + cssLength(css["left"]), vOff: cssLength(css["margin-top"]) + cssLength(css["top"])}
		if rel := css["mso-position-horizontal-relative"]; rel != "" {
			f.hRel = strings.TrimSuffix(rel, "-margin-area")
		}
		if rel := css["mso-position-vertical-relative"]; rel != "" {
			f.vRel = strings.TrimSuffix(rel, "-margin-area")
		}
		if z, err := strconv.ParseInt(css["z-index"], 10, 64); err == nil && z < 0 {
			f.behind = true
		}
		p.anchors = append(p.anchors, f)
		return
	}
	p.items = append(p.items, item{kind: kObject, obj: o, st: st, w: o.w})
}

func hasElem(n *ooxml.Node, name string) bool {
	for _, k := range n.Elements() {
		if k.Name == name || hasElem(k, name) {
			return true
		}
	}
	return false
}

func cssProps(s string) map[string]string {
	m := map[string]string{}
	for _, decl := range strings.Split(s, ";") {
		k, v, ok := strings.Cut(decl, ":")
		if ok {
			m[strings.TrimSpace(strings.ToLower(k))] = strings.TrimSpace(v)
		}
	}
	return m
}

// cssLength reads a VML length (pt, in, cm, mm, px; points without a unit).
func cssLength(v string) float64 {
	v = strings.TrimSpace(v)
	if strings.HasSuffix(v, "px") {
		f, _ := strconv.ParseFloat(strings.TrimSuffix(v, "px"), 64)
		return f * 0.75
	}
	return measure(v, 1, 0)
}

// drawObject draws a drawing into box, inside a FIGURE when it has
// alternative text.
func (e *emitter) drawObject(o *inlineObj, box drawingml.Box) {
	if o.alt != "" {
		e.cv.Obj.Mark(bdf.MarkFigure, o.alt)
	}
	c := e.c
	prev := c.curEmitter
	c.curEmitter = e
	switch {
	case o.paint != nil:
		o.paint(e, box)
	case o.graphic != nil:
		c.dr.DrawGraphic(e.cv, o.graphic, o.part, box)
	default:
		c.dr.DrawImage(e.cv, o.part, o.vmlRID, box)
	}
	c.curEmitter = prev
	e.cv.Drawn = true
	e.invalidate()
	if o.alt != "" {
		e.cv.Obj.Mark(bdf.MarkEnd, "")
	}
	if o.link != "" {
		e.cv.Obj.Link(f32(box.X), f32(box.Y), f32(box.W), f32(box.H), o.link)
	}
}

// The converter is the drawingml.Host of the document's drawings: Word
// shapes have no placeholders or master text styles, and their text boxes
// hold WordprocessingML paragraphs.

// Placeholder implements drawingml.Host.
func (c *converter) Placeholder(ph *ooxml.Node, part string) ([]drawingml.Inherited, bool) {
	return nil, true
}

// TextStyle implements drawingml.Host.
func (c *converter) TextStyle(phType string) *ooxml.Node { return nil }

// Field implements drawingml.Host.
func (c *converter) Field(typ string) (string, bool) { return "", false }

// Link implements drawingml.Host.
func (c *converter) Link(action string, target ooxml.Rel) string { return "" }

// DrawTextBox implements drawingml.TextBoxHost: the paragraphs of a text
// box are laid out in its text area and drawn in a BOX.
func (c *converter) DrawTextBox(cv *canvas.Canvas, content *ooxml.Node, part string, box drawingml.TextBox) {
	blocks, ok := c.boxes[content]
	if !ok {
		w := &walker{c: c, part: part, base: []*frame{{kind: fBox}}}
		blocks = w.blocks(content)
		c.boxes[content] = blocks
	}
	width := box.W
	if !box.Wrap {
		width = 1e5
	}
	// East Asian vertical text: upright characters in the turned box
	vertical := box.Vert == "eaVert" || box.Vert == "wordArtVert" || box.Vert == "mongolianVert"
	f := c.subflow(width, nil, vertical)
	f.blocks(blocks, nil)
	h := f.y
	dy := 0.0
	switch box.Anchor {
	case "ctr":
		dy = (box.H - h) / 2
	case "b":
		dy = box.H - h
	}
	e := c.curEmitter
	sub := &emitter{c: c, cv: cv, vertical: vertical}
	if e != nil {
		sub.pg = e.pg
		sub.view = e.view
	}
	cv.Obj.Save()
	cv.Transform(box.M)
	for _, o := range f.all() {
		o.dy += dy
		sub.op(&o)
	}
	sub.finish()
	cv.Obj.Restore()
	if e != nil {
		e.invalidate()
	}
}

func f32(v float64) float32 { return float32(v) }
