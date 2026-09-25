// Package pptx converts PowerPoint presentations (.pptx) into BDF documents.
//
// Slides are rendered directly from the DrawingML markup: every slide page
// gets a background layer, the non-placeholder shapes of its slide master
// and slide layout as master layers, and its own shapes as the body layer.
// Master and layout layers are identical objects on every slide that uses
// them, so the container stores them once. Text is laid out here (BDF has
// no layout engine) with the metrics of the fonts that are then embedded as
// subsets. See docs/design.md §3.3.
package pptx

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/imgconv"
)

// Options controls the conversion.
type Options struct {
	// Slides selects 1-based slide numbers; nil converts every slide
	// (hidden ones only with Hidden).
	Slides []int
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
	pkg      *pkg
	opts     *Options
	doc      *bdf.Document
	db       *fontdb.DB
	warnings []string
	warned   map[string]bool

	presPart     string
	pres         *node
	slideW       float64
	slideH       float64
	firstSlide   int
	defTextStyle *node
	themes       map[string]*theme

	canvases      []*canvas
	faceRunes     map[*fontdb.Face]map[rune]bool
	choices       map[resolveKey]*faceChoice
	fallback      map[fallbackKey]*faceChoice
	fallbackLists map[fallbackListKey][]fontdb.Resolved
	missing       map[rune]bool  // characters no available font has
	pageOf        map[string]int // slide part → page number in the output
	pages         int
	images        map[string]*imageEntry
	patterns      map[string]bdf.Hash
	tableStyles   map[string]*node
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
	p, err := openPkg(r, size)
	if err != nil {
		return nil, fmt.Errorf("pptx: %w", err)
	}
	c := &converter{pkg: p, opts: opts, doc: bdf.NewDocument(), warned: map[string]bool{},
		themes: map[string]*theme{}, faceRunes: map[*fontdb.Face]map[rune]bool{},
		choices: map[resolveKey]*faceChoice{}, fallback: map[fallbackKey]*faceChoice{},
		fallbackLists: map[fallbackListKey][]fontdb.Resolved{}, pageOf: map[string]int{}, missing: map[rune]bool{},
		images: map[string]*imageEntry{}, patterns: map[string]bdf.Hash{}}
	c.db = fontdb.New(opts.FontDirs, !opts.NoSystemFonts)
	if len(c.db.Faces) == 0 && !opts.SystemFonts {
		c.warnf("no fonts found; text is laid out with estimated metrics and not embedded")
	}

	c.presPart = "ppt/presentation.xml"
	if r, ok := p.relOfType("", "/officeDocument"); ok {
		c.presPart = r.Target
	}
	c.pres, err = p.xml(c.presPart)
	if err != nil {
		return nil, fmt.Errorf("pptx: %w", err)
	}
	sz := c.pres.child("sldSz")
	c.slideW, c.slideH = sz.emuAttr("cx", 720), sz.emuAttr("cy", 540)
	if c.slideW <= 0 || c.slideH <= 0 || c.slideW > 20000 || c.slideH > 20000 {
		c.slideW, c.slideH = 720, 540
	}
	c.firstSlide = int(c.pres.attrInt("firstSlideNum", 1))
	c.defTextStyle = c.pres.child("defaultTextStyle")
	c.loadTableStyles()

	c.doc.Meta.Source = "pptx"
	c.doc.Meta.Title = opts.Title
	if c.doc.Meta.Title == "" {
		c.doc.Meta.Title = c.coreTitle()
	}
	view := c.doc.NewView("slides", bdf.ViewFixed, c.doc.Meta.Title)

	var slides []string
	var hidden []bool
	for _, s := range c.pres.path("sldIdLst").children("sldId") {
		if r, ok := p.target(c.presPart, s.rid("id")); ok {
			slides = append(slides, r.Target)
			sn, err := p.xml(r.Target)
			hidden = append(hidden, err == nil && !sn.attrBool("show", true))
		}
	}
	sel := opts.Slides
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
		layers []*canvas
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
		var kept []*canvas
		for _, l := range layers {
			if l.cv.drawn {
				page.Layers = append(page.Layers, bdf.Layer{Role: l.role})
				kept = append(kept, l.cv)
			}
		}
		pages = append(pages, pageRef{page, kept})
	}
	c.finalize()
	if len(c.missing) > 0 {
		var rs []rune
		for r := range c.missing {
			rs = append(rs, r)
		}
		slices.Sort(rs)
		var b strings.Builder
		for i, r := range rs {
			if i == 20 {
				fmt.Fprintf(&b, " … (%d more)", len(rs)-20)
				break
			}
			fmt.Fprintf(&b, " %c U+%04X", r, r)
		}
		c.warnf("no available font has glyphs for:%s; viewers draw them with their own fonts", b.String())
	}
	for _, pr := range pages {
		for i, cv := range pr.layers {
			pr.page.Layers[i].Obj = cv.hash
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

// coreTitle reads dc:title from the core properties.
func (c *converter) coreTitle() string {
	r, ok := c.pkg.relOfType("", "/core-properties")
	if !ok {
		return ""
	}
	n, err := c.pkg.xml(r.Target)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(n.child("title").text())
}

func (c *converter) theme(masterPart string) *theme {
	r, ok := c.pkg.relOfType(masterPart, "/theme")
	if !ok {
		return defaultTheme
	}
	if th, ok := c.themes[r.Target]; ok {
		return th
	}
	n, err := c.pkg.xml(r.Target)
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

// layer is one page layer being drawn.
type layer struct {
	role string
	cv   *canvas
}

// slideCtx carries what rendering a slide needs from its layout, master
// and theme.
type slideCtx struct {
	c                     *converter
	part, layoutPart      string
	masterPart            string
	slide, layout, master *node
	th                    *theme
	cc                    *colorCtx
	num                   int
	layoutPh, masterPh    []*node
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
	slide, err := p.xml(part)
	if err != nil {
		return nil, err
	}
	s := &slideCtx{c: c, part: part, slide: slide, num: num}
	if r, ok := p.relOfType(part, "/slideLayout"); ok {
		s.layoutPart = r.Target
		s.layout, _ = p.xml(r.Target)
	}
	if s.layout != nil {
		if r, ok := p.relOfType(s.layoutPart, "/slideMaster"); ok {
			s.masterPart = r.Target
			s.master, _ = p.xml(r.Target)
		}
	}
	s.th = c.theme(s.masterPart)
	clrMap := map[string]string{}
	for k, v := range defaultClrMap {
		clrMap[k] = v
	}
	for _, a := range s.master.child("clrMap").attrs() {
		clrMap[a.Name.Local] = a.Value
	}
	for _, src := range []*node{s.layout, s.slide} {
		if ov := src.path("clrMapOvr", "overrideClrMapping"); ov != nil {
			for _, a := range ov.Attr {
				clrMap[a.Name.Local] = a.Value
			}
		}
	}
	s.cc = &colorCtx{scheme: s.th.colors, clrMap: clrMap}
	s.layoutPh = placeholders(s.layout)
	s.masterPh = placeholders(s.master)

	var layers []layer
	bg := layer{bdf.RoleBackground, c.newCanvas()}
	s.background(bg.cv)
	layers = append(layers, bg)
	showMaster := slide.attrBool("showMasterSp", true)
	if showMaster && s.layout.attrBool("showMasterSp", true) && s.master != nil {
		l := layer{bdf.RoleMaster, c.newCanvas()}
		s.drawTree(l.cv, s.master.path("cSld", "spTree"), s.masterPart, false)
		layers = append(layers, l)
	}
	if showMaster && s.layout != nil {
		l := layer{bdf.RoleMaster, c.newCanvas()}
		s.drawTree(l.cv, s.layout.path("cSld", "spTree"), s.layoutPart, false)
		layers = append(layers, l)
	}
	body := layer{bdf.RoleBody, c.newCanvas()}
	s.drawTree(body.cv, slide.path("cSld", "spTree"), part, true)
	layers = append(layers, body)
	return layers, nil
}

// placeholders lists the placeholder shapes of a layout or master.
func placeholders(root *node) []*node {
	var out []*node
	for _, k := range root.path("cSld", "spTree").kids() {
		if phOf(k) != nil {
			out = append(out, k)
		}
	}
	return out
}

// phOf returns the p:ph element of a shape, or nil.
func phOf(sh *node) *node {
	for _, nv := range []string{"nvSpPr", "nvPicPr", "nvGraphicFramePr", "nvGrpSpPr", "nvCxnSpPr"} {
		if ph := sh.path(nv, "nvPr", "ph"); ph != nil {
			return ph
		}
	}
	return nil
}

func phType(ph *node) string { return ph.attrStr("type", "obj") }

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
func findPh(list []*node, ph *node, master bool) *node {
	t := phType(ph)
	if !master {
		if idx, ok := ph.attr("idx"); ok {
			for _, k := range list {
				if v, ok := phOf(k).attr("idx"); ok && v == idx {
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
func (s *slideCtx) txStyleFor(t string) *node {
	ts := s.master.child("txStyles")
	switch masterPhType(t) {
	case "title":
		return ts.child("titleStyle")
	case "body":
		return ts.child("bodyStyle")
	}
	return ts.child("otherStyle")
}

func itoa(i int) string { return strconv.Itoa(i) }
