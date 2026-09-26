// Package drawingml draws DrawingML (ECMA-376 Part 1 §20, §21): the shapes,
// pictures, text bodies, tables and charts that Office documents share,
// with their themes and color maps. Text is laid out here, with the fonts
// that are then embedded as subsets (see fontset); Windows metafile
// pictures are replayed (see metafile).
//
// What depends on the kind of document (PresentationML placeholders and
// text styles, links between slides, slide numbers) comes from a Host.
package drawingml

import (
	"encoding/xml"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/imgconv"
)

// Config is what a Renderer draws with.
type Config struct {
	// Package is the document the drawings are in.
	Package *ooxml.Package
	// Doc receives the images and patterns.
	Doc *bdf.Document
	// Fonts lays out the text.
	Fonts *fontset.Set
	// Images controls whether raster images are re-encoded.
	Images imgconv.Options
	// Warn receives non-fatal problems.
	Warn func(msg string)
}

// Renderer draws the DrawingML of one document. It keeps what the pages
// share: images, pattern tiles, themes and table styles.
type Renderer struct {
	pkg         *ooxml.Package
	doc         *bdf.Document
	fonts       *fontset.Set
	imgOpts     imgconv.Options
	warn        func(msg string)
	warned      map[string]bool
	themes      map[string]*theme
	images      map[string]*imageEntry
	patterns    map[string]bdf.Hash
	tableStyles map[string]*ooxml.Node
}

// New returns a renderer for a document.
func New(cfg Config) *Renderer {
	return &Renderer{pkg: cfg.Package, doc: cfg.Doc, fonts: cfg.Fonts, imgOpts: cfg.Images, warn: cfg.Warn,
		warned: map[string]bool{}, themes: map[string]*theme{}, images: map[string]*imageEntry{},
		patterns: map[string]bdf.Hash{}, tableStyles: map[string]*ooxml.Node{}}
}

func (c *Renderer) warnf(format string, args ...any) {
	if c.warn != nil {
		c.warn(fmt.Sprintf(format, args...))
	}
}

func (c *Renderer) warnOnce(key, format string, args ...any) {
	if c.warned[key] {
		return
	}
	c.warned[key] = true
	c.warnf(format, args...)
}

// LoadTableStyles reads the table styles of the tableStyles part that part
// (the presentation) relates to.
func (c *Renderer) LoadTableStyles(part string) {
	if r, ok := c.pkg.RelOfType(part, "/tableStyles"); ok {
		if n, err := c.pkg.XML(r.Target); err == nil {
			for _, st := range n.Children("tblStyle") {
				c.tableStyles[strings.ToUpper(st.AttrStr("styleId", ""))] = st
			}
		}
	}
}

// theme returns the theme part relates to (a slide master, a workbook, a
// document), or the default theme.
func (c *Renderer) theme(part string) *theme {
	r, ok := c.pkg.RelOfType(part, "/theme")
	if !ok {
		return defaultTheme
	}
	if th, ok := c.themes[r.Target]; ok {
		return th
	}
	n, err := c.pkg.XML(r.Target)
	if err != nil {
		c.warnf("theme: %v", err)
		n = nil
	}
	th := parseTheme(n)
	if th != defaultTheme {
		th.part = r.Target
	}
	c.themes[r.Target] = th
	return th
}

// Host supplies what depends on the kind of document a drawing is in.
type Host interface {
	// Placeholder resolves a shape that is a placeholder (ph is its
	// p:ph element; part the part it is in): the shapes it inherits
	// properties from, most specific first, and whether it is drawn at
	// all (the placeholders of masters and layouts are prompts).
	Placeholder(ph *ooxml.Node, part string) (inherited []Inherited, draw bool)
	// TextStyle returns the list style the text of a shape inherits last:
	// the master's style for a placeholder type, or for phType "" (other
	// shapes) the document's default text style. It may return nil.
	TextStyle(phType string) *ooxml.Node
	// Field returns the text of a field (a:fld) of the given type, or false
	// to keep the text the field was saved with.
	Field(typ string) (string, bool)
	// Link returns the link ("#page=N") of a hyperlink that stays in the
	// document: an action (ppaction://…) or a relationship to another part;
	// "" when there is none.
	Link(action string, target ooxml.Rel) string
}

// Inherited is a shape a placeholder inherits from, with its part.
type Inherited struct {
	Node *ooxml.Node
	Part string
}

// Drawing draws the shapes of one page (a slide with its layout and master)
// with its theme and color map.
type Drawing struct {
	c    *Renderer
	th   *theme
	cc   *colorCtx
	host Host
}

// NewDrawing starts a page whose theme is the one part relates to (a slide
// master) and whose color map (bg1 → lt1 …) is the default one changed by
// clrMap.
func (c *Renderer) NewDrawing(part string, clrMap map[string]string, host Host) *Drawing {
	s := &Drawing{c: c, th: c.theme(part), host: host}
	m := map[string]string{}
	for k, v := range defaultClrMap {
		m[k] = v
	}
	for k, v := range clrMap {
		m[k] = v
	}
	s.cc = &colorCtx{scheme: s.th.colors, clrMap: m}
	return s
}

// DrawTree draws the shapes of a shape tree (p:spTree) of part.
func (s *Drawing) DrawTree(cv *canvas.Canvas, tree *ooxml.Node, part string) {
	for _, k := range tree.Elements() {
		s.drawElem(cv, k, part, nil)
	}
}

// DrawBackground fills a page of w×h with a background (p:bg, of part), or
// with the color scheme's bg1 when bg is nil.
func (s *Drawing) DrawBackground(cv *canvas.Canvas, bg *ooxml.Node, part string, w, h float64) {
	var f fill
	switch {
	case bg == nil:
		c, _ := s.cc.scheme1("bg1")
		f = fill{kind: fillSolid, color: c}
	case bg.Child("bgPr") != nil:
		f = s.resolveFill(fillElem(bg.Child("bgPr")), part, s.cc, nil)
	case bg.Child("bgRef") != nil:
		f = s.styleFill(bg.Child("bgRef"), s.cc)
	}
	s.fillPage(cv, f, w, h)
}

// PlaceholderOf returns the p:ph element of a shape, or nil.
func PlaceholderOf(sh *ooxml.Node) *ooxml.Node {
	for _, nv := range []string{"nvSpPr", "nvPicPr", "nvGraphicFramePr", "nvGrpSpPr", "nvCxnSpPr"} {
		if ph := sh.Path(nv, "nvPr", "ph"); ph != nil {
			return ph
		}
	}
	return nil
}

// ThemeColor returns a color of the drawing's color scheme by name (dk1,
// lt1, accent1 …, or a mapped name such as bg1), with components in 0..1.
func (s *Drawing) ThemeColor(name string) (r, g, b float64, ok bool) {
	c, ok := s.cc.scheme1(name)
	return c.R, c.G, c.B, ok
}

// ThemeFont resolves a theme font reference ("+mn-lt", "+mj-ea" …); script
// picks the font the theme lists for a script ("Jpan", "Hang" …) when its
// East Asian font is empty. It returns "" when the theme has none.
func (s *Drawing) ThemeFont(ref, script string) string { return s.th.fontFor(ref, script) }

// DrawAnchored draws a shape element of part (sp, cxnSp, pic, grpSp or
// graphicFrame) into the box x, y, w, h, the way a spreadsheet drawing's
// anchor places it: the anchor takes precedence over the offset and size
// of the shape's own transform (which the frames of charts often leave
// empty); its rotation and flips stay.
func (s *Drawing) DrawAnchored(cv *canvas.Canvas, k *ooxml.Node, part string, x, y, w, h float64) {
	own := ownXfrm(k)
	if xf, ok := parseXfrm(own); ok && xf.W > 0 && xf.H > 0 {
		g := &groupCtx{xf: xform{X: x, Y: y, W: w, H: h}, chOff: [2]float64{xf.X, xf.Y}, chExt: [2]float64{xf.W, xf.H}, part: part}
		s.drawElem(cv, k, part, g)
		return
	}
	s.drawElem(cv, withXfrm(k, own, x, y, w, h), part, nil)
}

// ownXfrm returns the transform element of a shape.
func ownXfrm(k *ooxml.Node) *ooxml.Node {
	switch k.Name {
	case "graphicFrame":
		return k.Child("xfrm")
	case "grpSp":
		return k.Path("grpSpPr", "xfrm")
	}
	return k.Path("spPr", "xfrm")
}

// withXfrm returns a copy of a shape whose transform places it in a box.
func withXfrm(k, own *ooxml.Node, x, y, w, h float64) *ooxml.Node {
	emu := func(v float64) string { return strconv.FormatInt(int64(math.Round(v*ooxml.EMUPerPoint)), 10) }
	xfrm := &ooxml.Node{Name: "xfrm", Kids: []*ooxml.Node{
		{Name: "off", Attrs: []xml.Attr{{Name: xml.Name{Local: "x"}, Value: emu(x)}, {Name: xml.Name{Local: "y"}, Value: emu(y)}}},
		{Name: "ext", Attrs: []xml.Attr{{Name: xml.Name{Local: "cx"}, Value: emu(w)}, {Name: xml.Name{Local: "cy"}, Value: emu(h)}}},
	}}
	if own != nil {
		xfrm.Space, xfrm.Attrs = own.Space, own.Attrs
	}
	for _, kid := range own.Elements() {
		if kid.Name != "off" && kid.Name != "ext" {
			xfrm.Kids = append(xfrm.Kids, kid)
		}
	}
	replace := func(n *ooxml.Node, name string, with *ooxml.Node) *ooxml.Node {
		cp := *n
		cp.Kids = nil
		found := false
		for _, kid := range n.Kids {
			if kid.Name == name {
				cp.Kids = append(cp.Kids, with)
				found = true
				continue
			}
			cp.Kids = append(cp.Kids, kid)
		}
		if !found {
			cp.Kids = append([]*ooxml.Node{with}, cp.Kids...)
		}
		return &cp
	}
	switch k.Name {
	case "graphicFrame":
		return replace(k, "xfrm", xfrm)
	case "grpSp":
		pr := k.Child("grpSpPr")
		if pr == nil {
			pr = &ooxml.Node{Name: "grpSpPr"}
		}
		return replace(k, "grpSpPr", replace(pr, "xfrm", xfrm))
	}
	pr := k.Child("spPr")
	if pr == nil {
		pr = &ooxml.Node{Name: "spPr"}
	}
	return replace(k, "spPr", replace(pr, "xfrm", xfrm))
}
