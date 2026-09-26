// Package pptx converts PowerPoint presentations (.pptx) into BDF documents.
//
// Slides are rendered directly from the DrawingML markup: every slide page
// gets a background layer, the non-placeholder shapes of its slide master
// and slide layout as master layers, and its own shapes as the body layer.
// Master and layout layers are identical objects on every slide that uses
// them, so the container stores them once. The shapes are drawn by the
// DrawingML renderer that the Office converters share
// (converter/internal/ooxml/drawingml), which lays text out with the
// metrics of the fonts that are then embedded as subsets; this package adds
// what PresentationML defines: slides, layouts, masters and their
// placeholders. See docs/design.md §3.4.
package pptx

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
	"github.com/shibukawa/bdf/imgconv"
)

// Options controls the conversion.
type Options struct {
	// Slides selects 1-based slides (see converter.Pages); nil converts
	// every slide (hidden ones only with Hidden).
	Slides conv.Pages
	// Hidden includes slides marked as hidden when Slides is nil.
	Hidden bool
	// Title overrides the document title.
	Title string
	// Images controls whether raster images are re-encoded (see imgconv).
	// The zero value keeps images as they are.
	Images imgconv.Options
	// FontDirs are searched for fonts before the system font directories.
	FontDirs []string
	// NoSystemFonts restricts font lookup to FontDirs.
	NoSystemFonts bool
	// SystemFonts refers to fonts by family name instead of embedding the
	// fonts used for layout. Viewers then substitute their own fonts; the
	// text advances keep lines at the widths computed here.
	SystemFonts bool
	// NoSubset embeds whole fonts instead of the glyphs in use.
	NoSubset bool
	// NoWOFF2 stores embedded fonts as TrueType/OpenType instead of WOFF2
	// (builds tagged bdf_noconv never produce WOFF2).
	NoWOFF2 bool
	// IgnoreFSType embeds fonts whose OS/2 fsType forbids embedding or
	// subsetting. Without it such fonts are referred to by name or embedded
	// whole. Set it only when you hold the rights to embed the fonts.
	IgnoreFSType bool
	// NoTextIndex skips building the text index part.
	NoTextIndex bool
	// Warn receives non-fatal problems; when nil they are collected in Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc           *bdf.Document
	Warnings      []string
	Slides        int
	EmbeddedFonts int
}

type converter struct {
	pkg      *ooxml.Package
	opts     *Options
	doc      *bdf.Document
	fonts    *fontset.Set
	cvs      *canvas.Builder
	r        *drawingml.Renderer
	warnings []string
	warned   map[string]bool

	presPart     string
	pres         *ooxml.Node
	slideW       float64
	slideH       float64
	firstSlide   int
	defTextStyle *ooxml.Node

	pageOf        map[string]int // slide part → page number in the output
	pages         int
	embeddedFonts int
}

// ConvertFile converts a .pptx file.
func ConvertFile(path string, opts *Options) (*Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return Convert(f, st.Size(), opts)
}

// Convert converts a presentation read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	p, err := ooxml.Open(r, size)
	if err != nil {
		return nil, fmt.Errorf("pptx: %w", err)
	}
	c := &converter{pkg: p, opts: opts, doc: bdf.NewDocument(), warned: map[string]bool{}, pageOf: map[string]int{}}
	warn := func(msg string) { c.warnf("%s", msg) }
	db := fontdb.New(opts.FontDirs, !opts.NoSystemFonts)
	c.fonts = fontset.New(db, warn)
	c.cvs = canvas.NewBuilder(c.doc, c.fonts)
	c.r = drawingml.New(drawingml.Config{Package: p, Doc: c.doc, Fonts: c.fonts, Images: opts.Images, Warn: warn})
	if len(db.Faces) == 0 && !opts.SystemFonts {
		c.warnf("no fonts found; text is laid out with estimated metrics and not embedded")
	}

	c.presPart = "ppt/presentation.xml"
	if r, ok := p.RelOfType("", "/officeDocument"); ok {
		c.presPart = r.Target
	}
	c.pres, err = p.XML(c.presPart)
	if err != nil {
		return nil, fmt.Errorf("pptx: %w", err)
	}
	sz := c.pres.Child("sldSz")
	c.slideW, c.slideH = sz.AttrEMU("cx", 720), sz.AttrEMU("cy", 540)
	if c.slideW <= 0 || c.slideH <= 0 || c.slideW > 20000 || c.slideH > 20000 {
		c.slideW, c.slideH = 720, 540
	}
	c.firstSlide = int(c.pres.AttrInt("firstSlideNum", 1))
	c.defTextStyle = c.pres.Child("defaultTextStyle")
	c.r.LoadTableStyles(c.presPart)

	c.doc.Meta.Source = "pptx"
	c.doc.Meta.DC = p.CoreProperties()
	if opts.Title != "" {
		c.doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	if len(c.doc.Meta.DC.Language) == 0 {
		// the default of the text (spec §7.8), when the core properties name none
		if l := c.textStyleLang(); l != "" {
			c.doc.Meta.DC.Language = bdf.DCValues{l}
		}
	}
	view := c.doc.NewView("slides", bdf.ViewFixed, c.doc.Meta.DC.Title.First())

	var slides []string
	var hidden []bool
	for _, s := range c.pres.Path("sldIdLst").Children("sldId") {
		if r, ok := p.Target(c.presPart, s.RelID("id")); ok {
			slides = append(slides, r.Target)
			sn, err := p.XML(r.Target)
			hidden = append(hidden, err == nil && !sn.AttrBool("show", true))
		}
	}
	sel := opts.Slides.Numbers(len(slides))
	if sel == nil {
		for i := range slides {
			if !hidden[i] || opts.Hidden {
				sel = append(sel, i+1)
			}
		}
	}
	for i, n := range sel {
		if n >= 1 && n <= len(slides) {
			c.pageOf[slides[n-1]] = i + 1
		}
	}
	c.pages = len(sel)
	type pageRef struct {
		page   *bdf.Page
		layers []*canvas.Canvas
	}
	var pages []pageRef
	for _, n := range sel {
		if n < 1 || n > len(slides) {
			return nil, fmt.Errorf("pptx: slide %d out of range (1-%d)", n, len(slides))
		}
		layers, err := c.renderSlideSafe(slides[n-1], c.firstSlide+n-1)
		if err != nil {
			return nil, fmt.Errorf("pptx: slide %d: %w", n, err)
		}
		page := view.AddPage(f32(c.slideW), f32(c.slideH))
		var kept []*canvas.Canvas
		for _, l := range layers {
			if l.cv.Drawn {
				page.Layers = append(page.Layers, bdf.Layer{Role: l.role})
				kept = append(kept, l.cv)
			}
		}
		pages = append(pages, pageRef{page, kept})
	}
	c.finalize()
	c.fonts.ReportMissing()
	for _, pr := range pages {
		for i, cv := range pr.layers {
			pr.page.Layers[i].Obj = cv.Hash()
		}
		if len(pr.page.Layers) == 0 {
			// An empty slide still needs something to draw.
			o := bdf.NewObject()
			o.SetBBox(0, 0, pr.page.W, pr.page.H)
			h, _ := c.doc.AddObject(o)
			pr.page.Layers = append(pr.page.Layers, bdf.Layer{Role: bdf.RoleBody, Obj: h})
		}
	}
	if !opts.NoTextIndex {
		if _, err := c.doc.BuildTextIndex(view); err != nil {
			c.warnf("text index: %v", err)
		}
	}
	return &Result{Doc: c.doc, Warnings: c.warnings, Slides: len(sel), EmbeddedFonts: c.embeddedFonts}, nil
}

func (c *converter) warnf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if c.opts.Warn != nil {
		c.opts.Warn(msg)
		return
	}
	c.warnings = append(c.warnings, msg)
}

func (c *converter) warnOnce(key, format string, args ...any) {
	if c.warned[key] {
		return
	}
	c.warned[key] = true
	c.warnf(format, args...)
}

// textStyleLang returns the language of the presentation's default text
// style: that of text whose runs name no language.
func (c *converter) textStyleLang() string {
	for _, name := range []string{"defPPr", "lvl1pPr"} {
		if l := canvas.NormLang(c.defTextStyle.Path(name, "defRPr").AttrStr("lang", "")); l != "" {
			return l
		}
	}
	return ""
}

// layer is one page layer being drawn.
type layer struct {
	role string
	cv   *canvas.Canvas
}

// slideCtx is a slide with its layout and master. It is the host of the
// slide's drawing: placeholders inherit from the layout and the master,
// whose text styles apply.
type slideCtx struct {
	c                     *converter
	d                     *drawingml.Drawing
	part, layoutPart      string
	masterPart            string
	slide, layout, master *ooxml.Node
	num                   int
	layoutPh, masterPh    []*ooxml.Node
}

// renderSlideSafe renders a slide; a malformed slide that trips the
// renderer becomes an empty page and a warning instead of failing the
// whole conversion.
func (c *converter) renderSlideSafe(part string, num int) (layers []layer, err error) {
	defer func() {
		if r := recover(); r != nil {
			c.warnf("slide %s: internal error: %v", part, r)
			layers, err = nil, nil
		}
	}()
	return c.renderSlide(part, num)
}

func (c *converter) renderSlide(part string, num int) ([]layer, error) {
	p := c.pkg
	slide, err := p.XML(part)
	if err != nil {
		return nil, err
	}
	s := &slideCtx{c: c, part: part, slide: slide, num: num}
	if r, ok := p.RelOfType(part, "/slideLayout"); ok {
		s.layoutPart = r.Target
		s.layout, _ = p.XML(r.Target)
	}
	if s.layout != nil {
		if r, ok := p.RelOfType(s.layoutPart, "/slideMaster"); ok {
			s.masterPart = r.Target
			s.master, _ = p.XML(r.Target)
		}
	}
	clrMap := map[string]string{}
	if cm := s.master.Child("clrMap"); cm != nil {
		for _, a := range cm.Attrs {
			clrMap[a.Name.Local] = a.Value
		}
	}
	for _, src := range []*ooxml.Node{s.layout, s.slide} {
		if ov := src.Path("clrMapOvr", "overrideClrMapping"); ov != nil {
			for _, a := range ov.Attrs {
				clrMap[a.Name.Local] = a.Value
			}
		}
	}
	s.d = c.r.NewDrawing(s.masterPart, clrMap, s)
	s.layoutPh = placeholders(s.layout)
	s.masterPh = placeholders(s.master)

	var layers []layer
	bg := layer{bdf.RoleBackground, c.cvs.New()}
	s.background(bg.cv)
	layers = append(layers, bg)
	showMaster := slide.AttrBool("showMasterSp", true)
	if showMaster && s.layout.AttrBool("showMasterSp", true) && s.master != nil {
		l := layer{bdf.RoleMaster, c.cvs.New()}
		s.d.DrawTree(l.cv, s.master.Path("cSld", "spTree"), s.masterPart)
		layers = append(layers, l)
	}
	if showMaster && s.layout != nil {
		l := layer{bdf.RoleMaster, c.cvs.New()}
		s.d.DrawTree(l.cv, s.layout.Path("cSld", "spTree"), s.layoutPart)
		layers = append(layers, l)
	}
	body := layer{bdf.RoleBody, c.cvs.New()}
	s.d.DrawTree(body.cv, slide.Path("cSld", "spTree"), part)
	layers = append(layers, body)
	return layers, nil
}

// placeholders lists the placeholder shapes of a layout or master.
func placeholders(root *ooxml.Node) []*ooxml.Node {
	var out []*ooxml.Node
	for _, k := range root.Path("cSld", "spTree").Elements() {
		if phOf(k) != nil {
			out = append(out, k)
		}
	}
	return out
}

// phOf returns the p:ph element of a shape, or nil.
func phOf(sh *ooxml.Node) *ooxml.Node { return drawingml.PlaceholderOf(sh) }

func phType(ph *ooxml.Node) string { return ph.AttrStr("type", "obj") }

// masterPhType maps a placeholder type to the master placeholder it inherits from.
func masterPhType(t string) string {
	switch t {
	case "ctrTitle", "title":
		return "title"
	case "dt", "ftr", "sldNum", "hdr":
		return t
	}
	return "body"
}

// findPh finds the placeholder a shape inherits from: by index, then by type.
func findPh(list []*ooxml.Node, ph *ooxml.Node, master bool) *ooxml.Node {
	t := phType(ph)
	if !master {
		if idx, ok := ph.Attr("idx"); ok {
			for _, k := range list {
				if v, ok := phOf(k).Attr("idx"); ok && v == idx {
					return k
				}
			}
		}
		for _, k := range list {
			if phType(phOf(k)) == t {
				return k
			}
		}
		if t == "ctrTitle" || t == "title" {
			for _, k := range list {
				if tt := phType(phOf(k)); tt == "title" || tt == "ctrTitle" {
					return k
				}
			}
		}
		return nil
	}
	mt := masterPhType(t)
	for _, k := range list {
		if phType(phOf(k)) == mt {
			return k
		}
	}
	return nil
}

// txStyleFor returns the master text style a placeholder type uses.
func (s *slideCtx) txStyleFor(t string) *ooxml.Node {
	ts := s.master.Child("txStyles")
	switch masterPhType(t) {
	case "title":
		return ts.Child("titleStyle")
	case "body":
		return ts.Child("bodyStyle")
	}
	return ts.Child("otherStyle")
}

// background fills the page with the first background of slide, layout
// and master (the color scheme's bg1 when none has one).
func (s *slideCtx) background(cv *canvas.Canvas) {
	for _, src := range []struct {
		n    *ooxml.Node
		part string
	}{{s.slide, s.part}, {s.layout, s.layoutPart}, {s.master, s.masterPart}} {
		if bg := src.n.Path("cSld", "bg"); bg != nil {
			s.d.DrawBackground(cv, bg, src.part, s.c.slideW, s.c.slideH)
			return
		}
	}
	s.d.DrawBackground(cv, nil, "", s.c.slideW, s.c.slideH)
}

// Placeholder implements drawingml.Host: the placeholders of a slide
// inherit from the layout's placeholder, which inherits from the
// master's; those of layouts and masters are not drawn.
func (s *slideCtx) Placeholder(ph *ooxml.Node, part string) ([]drawingml.Inherited, bool) {
	if part != s.part {
		return nil, false
	}
	var inh []drawingml.Inherited
	mph := ph
	if lp := findPh(s.layoutPh, ph, false); lp != nil {
		inh = append(inh, drawingml.Inherited{Node: lp, Part: s.layoutPart})
		mph = phOf(lp)
	}
	if mp := findPh(s.masterPh, mph, true); mp != nil {
		inh = append(inh, drawingml.Inherited{Node: mp, Part: s.masterPart})
	}
	return inh, true
}

// TextStyle implements drawingml.Host: placeholders take the master's text
// style for their type, other shapes the presentation's default text style.
func (s *slideCtx) TextStyle(phType string) *ooxml.Node {
	if phType == "" {
		return s.c.defTextStyle
	}
	return s.txStyleFor(phType)
}

// Field implements drawingml.Host: slide numbers.
func (s *slideCtx) Field(typ string) (string, bool) {
	if strings.HasPrefix(typ, "slidenum") {
		return itoa(s.num), true
	}
	return "", false
}

// Link implements drawingml.Host: jumps to other slides.
func (s *slideCtx) Link(action string, target ooxml.Rel) string {
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
	if strings.HasSuffix(target.Type, "/slide") {
		return pageLink(s.c.pageOf[target.Target])
	}
	return ""
}

func pageLink(n int) string {
	if n < 1 {
		return ""
	}
	return "#page=" + itoa(n)
}

func itoa(i int) string { return strconv.Itoa(i) }

func f32(v float64) float32 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return float32(v)
}

// finalize builds the embedded fonts and encodes every canvas.
func (c *converter) finalize() {
	if !c.opts.SystemFonts {
		c.embeddedFonts = c.fonts.Embed(c.doc, fontset.EmbedOptions{NoSubset: c.opts.NoSubset, NoWOFF2: c.opts.NoWOFF2, IgnoreFSType: c.opts.IgnoreFSType})
	}
	c.cvs.Encode()
}
