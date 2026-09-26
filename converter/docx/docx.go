// Package docx converts Word documents (.docx) into BDF documents.
//
// BDF has no layout engine, so the document is laid out here, with the
// metrics of the fonts that are then embedded as subsets: paragraphs are
// broken into lines (tab stops, East Asian line breaking rules and
// spacing, justification, the document grid), tables into rows and cells,
// and the lines and rows into columns and pages (keeping lines and
// paragraphs together, widows and orphans, floating objects, footnotes).
// Drawings are rendered by the DrawingML renderer the Office converters
// share (converter/internal/ooxml/drawingml); their text boxes hold
// paragraphs laid out here.
//
// A document converts into two views (docs/spec.md §4.1): a flow view of
// pages, with header, body and footer layers and the body rectangle of
// each page, and a scroll view laid out once more without pages, as one
// long column of the text width (what Word's draft and web layouts show),
// stored as strips cut between lines. The views share fonts and images.
// Sections of East Asian vertical text are laid out as pages turned by 90°
// in the page view, and horizontally in the scroll view. See
// docs/design.md §3.9.
package docx

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"slices"
	"strconv"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
	"github.com/shibukawa/bdf/imgconv"
)

// View selections.
const (
	ViewsBoth   = "both"
	ViewsPages  = "pages"
	ViewsScroll = "scroll"
)

// Options controls the conversion.
type Options struct {
	// Pages selects 1-based pages of the page view; nil keeps every page.
	// The scroll view always holds the whole document.
	Pages []int
	// Views selects the views to make: ViewsBoth (the default when ""),
	// ViewsPages or ViewsScroll.
	Views string
	// Title overrides the document title.
	Title string
	// Images controls whether raster images are re-encoded (see imgconv).
	// The zero value keeps images as they are.
	Images imgconv.Options
	// FontFS holds fonts that are not in the local file system; it is
	// searched before FontDirs (see converter.Options.FontFS).
	FontFS fs.FS
	// FontDirs are searched for fonts before the system font directories.
	FontDirs []string
	// NoSystemFonts restricts font lookup to FontFS and FontDirs.
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
	// subsetting. Set it only when you hold the rights to embed the fonts.
	IgnoreFSType bool
	// NoTextIndex skips building the text index parts.
	NoTextIndex bool
	// Warn receives non-fatal problems; when nil they are collected in
	// Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc           *bdf.Document
	Warnings      []string
	Pages         int // pages of the page view
	Strips        int // strips of the scroll view
	EmbeddedFonts int
}

type converter struct {
	pkg      *ooxml.Package
	opts     *Options
	doc      *bdf.Document
	fonts    *fontset.Set
	cvs      *canvas.Builder
	r        *drawingml.Renderer
	dr       *drawingml.Drawing
	warnings []string
	warned   map[string]bool

	main        string
	st          *styles
	num         *numbering
	sections    []*section
	lastSection *section
	clrMap      map[string]string
	eaScript    string
	baseSize    float64 // the default font size: the character grid's cell less its extra pitch

	// settings
	defTab                 float64
	evenOdd                bool
	compatMode             int
	noExpandShiftReturn    bool
	suppressSpBfAfterPgBrk bool
	background             bdf.Color
	hasBackground          bool
	footFmt, endFmt        string

	notes      map[string]*note
	noteOrder  []*note // referenced notes in document order
	notesUsed  map[string]int
	hf         map[string][]block // header and footer blocks by part
	boxes      map[*ooxml.Node][]block
	noteCache  map[noteKey]*flow
	hfCache    map[hfKey]*flow
	bmPage     map[string][2]int // bookmark → physical page index, printed number
	bmY        map[string]float64
	curEmitter *emitter
	cjk, latin int
	tcyGroups  int
}

// note is a footnote or an endnote.
type note struct {
	kind   string // "f" or "e"
	id     string
	num    string
	node   *ooxml.Node
	part   string
	blocks []block
}

type noteUse struct {
	n *note
	h float64
}

type noteKey struct {
	n        *note
	w        float64
	vertical bool
}

type hfKey struct {
	part string
	sec  *section
}

// ConvertFile converts a .docx file.
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

// supported are the markup compatibility namespaces whose content is
// drawn (their mc:Choice is taken over the fallback).
var supported = map[string]bool{"wps": true, "wpg": true, "wpc": true, "wp14": true, "w14": true, "w15": true, "a14": true, "pic14": true}

// Convert converts a document read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (res *Result, err error) {
	defer func() {
		if e := recover(); e != nil {
			res, err = nil, fmt.Errorf("docx: internal error: %v", e)
		}
	}()
	if opts == nil {
		opts = &Options{}
	}
	views := opts.Views
	switch views {
	case "":
		views = ViewsBoth
	case ViewsBoth, ViewsPages, ViewsScroll:
	default:
		return nil, fmt.Errorf("docx: unknown views %q (want both, pages or scroll)", views)
	}
	p, err := ooxml.Open(r, size)
	if err != nil {
		return nil, fmt.Errorf("docx: %w", err)
	}
	p.Supported = func(prefix string) bool { return supported[prefix] }
	c := &converter{pkg: p, opts: opts, doc: bdf.NewDocument(), warned: map[string]bool{}, clrMap: map[string]string{},
		notes: map[string]*note{}, notesUsed: map[string]int{}, hf: map[string][]block{}, boxes: map[*ooxml.Node][]block{},
		noteCache: map[noteKey]*flow{}, hfCache: map[hfKey]*flow{}, bmPage: map[string][2]int{}, bmY: map[string]float64{},
		defTab: 36, compatMode: 12, footFmt: "decimal", endFmt: "lowerRoman"}
	warn := func(msg string) { c.warnf("%s", msg) }
	db := fontdb.New(opts.FontFS, opts.FontDirs, !opts.NoSystemFonts)
	c.fonts = fontset.New(db, warn)
	c.cvs = canvas.NewBuilder(c.doc, c.fonts)
	c.r = drawingml.New(drawingml.Config{Package: p, Doc: c.doc, Fonts: c.fonts, Images: opts.Images, Warn: warn})
	if len(db.Faces) == 0 && !opts.SystemFonts {
		c.warnf("no fonts found; text is laid out with estimated metrics and not embedded")
	}
	c.main = "word/document.xml"
	if rel, ok := p.RelOfType("", "/officeDocument"); ok {
		c.main = rel.Target
	}
	root, err := p.XML(c.main)
	if err != nil {
		return nil, fmt.Errorf("docx: %w", err)
	}
	c.dr = c.r.NewDrawing(c.main, nil, c)
	c.loadSettings()
	c.loadStyles()
	c.baseSize = c.runProps(c.st.defPara, "", nil).sz
	c.eaScript = fontset.Script(c.themeFontLang())
	c.loadNumbering()
	c.loadNotes()

	c.doc.Meta.Source = "docx"
	c.doc.Meta.DC = p.CoreProperties()
	if opts.Title != "" {
		c.doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	if bg := root.Child("background"); bg != nil {
		c.background, c.hasBackground = c.color(bg.AttrStr("color", "auto"), bg.AttrStr("themeColor", ""), bg.AttrStr("themeTint", ""), bg.AttrStr("themeShade", ""))
		c.hasBackground = c.hasBackground && c.displayBackground()
	}
	body := root.Child("body")
	w := &walker{c: c, part: c.main}
	blocks := w.blocks(body)
	c.splitSections(blocks, body.Child("sectPr"))
	if len(c.doc.Meta.DC.Language) == 0 {
		if l := c.docLang(); l != "" {
			c.doc.Meta.DC.Language = bdf.DCValues{l}
		}
	}

	res = &Result{Doc: c.doc}
	title := c.doc.Meta.DC.Title.First()
	var jobs []func()
	if views != ViewsScroll {
		v := c.doc.NewView("pages", bdf.ViewFlow, title)
		v.Continuous = &bdf.Continuous{Gap: 24}
		jobs = append(jobs, c.pageView(v, &res.Pages))
	}
	if views != ViewsPages {
		v := c.doc.NewView("scroll", bdf.ViewScroll, title)
		jobs = append(jobs, c.scrollView(v, &res.Strips))
	}
	if !opts.SystemFonts {
		res.EmbeddedFonts = c.fonts.Embed(c.doc, fontset.EmbedOptions{NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType})
	}
	c.cvs.Encode()
	c.fonts.ReportMissing()
	for _, job := range jobs {
		job()
	}
	if !opts.NoTextIndex {
		for _, v := range c.doc.Views {
			if _, err := c.doc.BuildTextIndex(v); err != nil {
				c.warnf("text index: %v", err)
			}
		}
	}
	res.Warnings = c.warnings
	return res, nil
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

func (c *converter) settings() *ooxml.Node {
	r, ok := c.pkg.RelOfType(c.main, "/settings")
	if !ok {
		return nil
	}
	n, err := c.pkg.XML(r.Target)
	if err != nil {
		return nil
	}
	return n
}

func (c *converter) loadSettings() {
	s := c.settings()
	if s == nil {
		return
	}
	if t := s.Child("defaultTabStop"); t != nil {
		if v := twips(t, "val", 36); v > 0 {
			c.defTab = v
		}
	}
	c.evenOdd = onOff(s.Child("evenAndOddHeaders"))
	if m := s.Child("clrSchemeMapping"); m != nil {
		for _, a := range []struct{ attr, name string }{{"bg1", "background1"}, {"t1", "text1"}, {"bg2", "background2"}, {"t2", "text2"}} {
			if v, ok := m.Attr(a.attr); ok {
				c.clrMap[a.name] = v
			}
		}
	}
	for k, v := range map[string]string{"background1": "light1", "text1": "dark1", "background2": "light2", "text2": "dark2"} {
		if _, ok := c.clrMap[k]; !ok {
			c.clrMap[k] = v
		}
	}
	compat := s.Child("compat")
	c.noExpandShiftReturn = onOff(compat.Child("doNotExpandShiftReturn"))
	c.suppressSpBfAfterPgBrk = onOff(compat.Child("suppressSpBfAfterPgBrk"))
	for _, cs := range compat.Children("compatSetting") {
		if cs.AttrStr("name", "") == "compatibilityMode" {
			if v, err := strconv.Atoi(cs.AttrStr("val", "")); err == nil {
				c.compatMode = v
			}
		}
	}
	if f := s.Path("footnotePr", "numFmt"); f != nil {
		c.footFmt = val(f)
	}
	if f := s.Path("endnotePr", "numFmt"); f != nil {
		c.endFmt = val(f)
	}
}

func (c *converter) displayBackground() bool {
	return onOff(c.settings().Child("displayBackgroundShape"))
}

// themeFontLang returns the language whose script picks the theme's East
// Asian fonts.
func (c *converter) themeFontLang() string {
	if l := c.settings().Child("themeFontLang").AttrStr("eastAsia", ""); l != "" {
		return l
	}
	return c.st.docRPr.Child("lang").AttrStr("eastAsia", "")
}

// docLang is the document's language when the core properties name none:
// that of the default run properties, East Asian when most text is.
func (c *converter) docLang() string {
	l := c.st.docRPr.Child("lang")
	latin, ea := canvas.NormLang(l.AttrStr("val", "")), canvas.NormLang(l.AttrStr("eastAsia", ""))
	if ea != "" && c.cjk > c.latin {
		return ea
	}
	if latin != "" {
		return latin
	}
	return ea
}

// splitSections groups the blocks into sections: each ends with the
// paragraph that holds its section properties, the last with the body's.
func (c *converter) splitSections(blocks []block, last *ooxml.Node) {
	var cur []block
	for _, b := range blocks {
		cur = append(cur, b)
		if p, ok := b.(*para); ok && p.sect != nil {
			p.sect.blocks = cur
			c.sections = append(c.sections, p.sect)
			cur = nil
		}
	}
	s := c.parseSection(last, c.lastSection)
	s.blocks = cur
	c.sections = append(c.sections, s)
	// header and footer references carry over from the previous section
	for i, s := range c.sections {
		if i == 0 {
			continue
		}
		prev := c.sections[i-1]
		for k, v := range prev.hdr {
			if _, ok := s.hdr[k]; !ok {
				s.hdr[k] = v
			}
		}
		for k, v := range prev.ftr {
			if _, ok := s.ftr[k]; !ok {
				s.ftr[k] = v
			}
		}
	}
}

// loadNotes reads the footnotes and endnotes.
func (c *converter) loadNotes() {
	for _, k := range []struct{ kind, rel, elem string }{{"f", "/footnotes", "footnote"}, {"e", "/endnotes", "endnote"}} {
		r, ok := c.pkg.RelOfType(c.main, k.rel)
		if !ok {
			continue
		}
		root, err := c.pkg.XML(r.Target)
		if err != nil {
			c.warnf("%s: %v", k.elem, err)
			continue
		}
		for _, n := range root.Children(k.elem) {
			if t := n.AttrStr("type", "normal"); t != "normal" {
				continue
			}
			id := n.AttrStr("id", "")
			c.notes[k.kind+id] = &note{kind: k.kind, id: id, node: n, part: r.Target}
		}
	}
}

// noteReference numbers a note where it is referenced and returns its
// number.
func (c *converter) noteReference(kind, id string) string {
	n := c.notes[kind+id]
	if n == nil {
		return ""
	}
	if n.num == "" {
		c.notesUsed[kind]++
		f := c.footFmt
		if kind == "e" {
			f = c.endFmt
		}
		n.num = formatNumber(c.notesUsed[kind], f)
		c.noteOrder = append(c.noteOrder, n)
		w := &walker{c: c, part: n.part, noteNo: n.num}
		n.blocks = w.blocks(n.node)
	}
	return n.num
}

// noteFlow lays out a note's content in the line length of a section.
func (c *converter) noteFlow(n *note, sec *section) *flow {
	w, vertical := 468.0, false
	if sec != nil {
		w, vertical = sec.lineLength(), sec.vertical
	}
	k := noteKey{n, w, vertical}
	if f, ok := c.noteCache[k]; ok {
		return f
	}
	f := c.subflow(w, sec, vertical)
	f.blocks(n.blocks, nil)
	c.noteCache[k] = f
	return f
}

func (c *converter) noteHeight(n *note, sec *section) float64 { return c.noteFlow(n, sec).y }

// hfFlow lays out a header or footer part for a section.
func (c *converter) hfFlow(part string, sec *section) *flow {
	k := hfKey{part, sec}
	if f, ok := c.hfCache[k]; ok {
		return f
	}
	blocks, ok := c.hf[part]
	if !ok {
		if root, err := c.pkg.XML(part); err == nil {
			w := &walker{c: c, part: part}
			blocks = w.blocks(root)
		} else {
			c.warnf("header/footer: %v", err)
		}
		c.hf[part] = blocks
	}
	// headers and footers are horizontal also in sections of vertical text
	f := c.subflow(sec.textWidth(), sec, false)
	f.origin = &[2]float64{sec.left, sec.header}
	f.blocks(blocks, nil)
	c.hfCache[k] = f
	return f
}

// hfPart returns the header (or footer) part of a page.
func (c *converter) hfPart(pg *page, refs map[string]string) string {
	switch {
	case pg.first && pg.sec.titlePg:
		return refs["first"]
	case c.evenOdd && pg.num%2 == 0:
		return refs["even"]
	}
	return refs["default"]
}

// bodyArea returns the top and bottom of a page's body: the margins,
// pushed in by a header or footer taller than they leave room for.
func (c *converter) bodyArea(pg *page) (top, bottom float64) {
	s := pg.sec
	top, bottom = s.top, s.pgH-s.bottom
	if part := c.hfPart(pg, s.hdr); part != "" && !s.topExact {
		if h := c.hfFlow(part, s).y; s.header+h > top {
			top = s.header + h
		}
	}
	if part := c.hfPart(pg, s.ftr); part != "" && !s.bottomExact {
		if h := c.hfFlow(part, s).y; s.pgH-s.footer-h < bottom {
			bottom = s.pgH - s.footer - h
		}
	}
	if bottom-top < 36 {
		bottom = top + 36
	}
	return
}

// recordBookmarks notes where the bookmarks of a paragraph went.
func (c *converter) recordBookmarks(f *flow, p *para, y float64) {
	if !f.main || len(p.bookmarks) == 0 {
		return
	}
	for _, b := range p.bookmarks {
		if f.paged {
			c.bmPage[b] = [2]int{len(f.pages) - 1, f.pg.num}
		} else {
			c.bmY[b] = y
		}
	}
}

// pageView lays out the pages. It returns the job that fills in the
// object hashes of the layers once the canvases are encoded.
func (c *converter) pageView(v *bdf.View, count *int) func() {
	f := &flow{c: c, paged: true, main: true, atTop: true, pageStart: -1}
	f.sections(c.sections)
	c.endnotes(f)
	pages := f.pages
	sel := c.opts.Pages
	keep := make([]bool, len(pages))
	for i := range pages {
		keep[i] = sel == nil || slices.Contains(sel, i+1)
	}
	index := map[int]int{} // physical page → page in the view
	for i := range pages {
		if keep[i] {
			index[i] = len(index) + 1
		}
	}
	secPages := map[int]int{}
	for _, pg := range pages {
		if !pg.blank {
			secPages[pg.secNo]++
		}
	}
	vi := &viewInfo{paged: true, numPages: len(pages), secPages: secPages,
		pageNum: func(b string) (int, bool) {
			p, ok := c.bmPage[b]
			return p[1], ok
		},
		link: func(b string) string {
			if p, ok := c.bmPage[b]; ok {
				if n, ok := index[p[0]]; ok {
					return "#page=" + strconv.Itoa(n)
				}
			}
			return ""
		}}
	type layerRef struct {
		page *bdf.Page
		cvs  []*canvas.Canvas
	}
	var refs []layerRef
	for i, pg := range pages {
		if !keep[i] {
			continue
		}
		s := pg.sec
		bp := v.AddPage(f32(s.pgW), f32(s.pgH))
		bp.Body = &bdf.RectDef{X: f32(s.left), Y: f32(pg.top), W: f32(s.textWidth()), H: f32(pg.bottom - pg.top)}
		var cvs []*canvas.Canvas
		add := func(role string, cv *canvas.Canvas, always bool) {
			if cv.Drawn || always {
				bp.Layers = append(bp.Layers, bdf.Layer{Role: role})
				cvs = append(cvs, cv)
			}
		}
		add(bdf.RoleBackground, c.pageBackground(pg), false)
		add(bdf.RoleHeader, c.hfLayer(pg, vi, s.hdr, true), false)
		add(bdf.RoleFooter, c.hfLayer(pg, vi, s.ftr, false), false)
		add(bdf.RoleBody, c.bodyLayer(pg, vi), true)
		refs = append(refs, layerRef{bp, cvs})
	}
	*count = len(refs)
	return func() {
		for _, r := range refs {
			for i, cv := range r.cvs {
				r.page.Layers[i].Obj = cv.Hash()
			}
		}
	}
}

// endnotes places the endnotes after the last section.
func (c *converter) endnotes(f *flow) {
	var ends []*note
	for _, n := range c.noteOrder {
		if n.kind == "e" {
			ends = append(ends, n)
		}
	}
	if len(ends) == 0 {
		return
	}
	f.separator()
	for _, n := range ends {
		f.blocks(n.blocks, nil)
	}
}

// separator draws the short line that separates notes from the text.
func (f *flow) separator() {
	y := f.y + f.prevAfter + 6
	x := f.x0
	w := math.Min(144, f.w/3)
	f.emit(op{y0: y, y1: y + 0.5, fn: func(e *emitter, dx, dy float64) {
		e.fillRect(bdf.RGB(0, 0, 0), x+dx, y+dy, w, 0.5)
	}})
	f.y = y + 6
	f.prevAfter = 0
	f.prevPara = nil
}

func (c *converter) pageBackground(pg *page) *canvas.Canvas {
	cv := c.cvs.New()
	e := &emitter{c: c, cv: cv, pg: pg}
	s := pg.sec
	if c.hasBackground {
		e.fillRect(c.background, 0, 0, s.pgW, s.pgH)
	}
	if !(s.bordersFirst && !pg.first) && !(s.bordersNotFirst && pg.first) {
		x0, y0, x1, y1 := 0.0, 0.0, s.pgW, s.pgH
		b := s.borders
		if s.bordersOffsetFromText {
			x0, y0, x1, y1 = s.left, pg.top, s.pgW-s.right, pg.bottom
			if b[0] != nil {
				y0 -= b[0].space
			}
			if b[1] != nil {
				x0 -= b[1].space
			}
			if b[2] != nil {
				y1 += b[2].space
			}
			if b[3] != nil {
				x1 += b[3].space
			}
		} else {
			if b[0] != nil {
				y0 += b[0].space
			}
			if b[1] != nil {
				x0 += b[1].space
			}
			if b[2] != nil {
				y1 -= b[2].space
			}
			if b[3] != nil {
				x1 -= b[3].space
			}
		}
		e.rule(b[0], x0, y0, x1, y0)
		e.rule(b[2], x0, y1, x1, y1)
		e.rule(b[1], x0, y0, x0, y1)
		e.rule(b[3], x1, y0, x1, y1)
	}
	return cv
}

// hfLayer draws a page's header or footer.
func (c *converter) hfLayer(pg *page, vi *viewInfo, refs map[string]string, header bool) *canvas.Canvas {
	cv := c.cvs.New()
	part := c.hfPart(pg, refs)
	if part == "" || pg.blank {
		return cv
	}
	s := pg.sec
	f := c.hfFlow(part, s)
	y := s.header
	if !header {
		y = s.pgH - s.footer - f.y
	}
	e := &emitter{c: c, cv: cv, pg: pg, view: vi}
	for _, o := range f.all() {
		o.shift(s.left, y)
		e.op(&o)
	}
	e.finish()
	return cv
}

// bodyLayer draws a page's body: the floating objects behind the text,
// the text, the footnotes, and the floating objects in front.
func (c *converter) bodyLayer(pg *page, vi *viewInfo) *canvas.Canvas {
	cv := c.cvs.New()
	o := bdf.Rect{X: 0, Y: 0, W: f32(pg.sec.pgW), H: f32(pg.sec.pgH)}
	cv.Obj.SetBBox(o.X, o.Y, o.W, o.H)
	e := &emitter{c: c, cv: cv, pg: pg, view: vi}
	for i := range pg.behind {
		e.op(&pg.behind[i])
	}
	s := pg.sec
	if s.vertical {
		// the text in the turned coordinates of the flow
		cv.Obj.Save()
		cv.Transform(pg.rotation())
		e.vertical = true
	}
	for i := range pg.body {
		e.op(&pg.body[i])
	}
	if len(pg.notes) > 0 {
		e.close(0)
		x, _ := s.column(0, pg.lineLen())
		_, bottom := pg.block()
		y := bottom - pg.noteH
		e.fillRect(bdf.RGB(0, 0, 0), x, y+noteSeparator/2, math.Min(144, pg.lineLen()/3), 0.5)
		y += noteSeparator
		for _, u := range pg.notes {
			for _, op := range c.noteFlow(u.n, s).all() {
				op.shift(x, y)
				e.op(&op)
			}
			e.close(0)
			y += u.h
		}
	}
	if s.vertical {
		cv.Obj.Restore()
		e.invalidate()
		e.vertical = false
	}
	for i := range pg.front {
		e.op(&pg.front[i])
	}
	e.finish()
	return cv
}

// scrollPad is the margin around the scroll view's column.
const scrollPad = 36

// stripHeight is the height the scroll view's strips are cut at (the
// cuts move up to the nearest line boundary).
const stripHeight = 1024

// scrollView lays the document out as one column without pages and cuts
// it into strips.
func (c *converter) scrollView(v *bdf.View, count *int) func() {
	width := 0.0
	for _, s := range c.sections {
		width = math.Max(width, s.textWidth())
	}
	f := &flow{c: c, main: true, atTop: true, pageStart: -1, bottom: math.Inf(1)}
	f.y = scrollPad
	for _, s := range c.sections {
		f.sec = s
		// vertical text too is laid out horizontally, as Word's draft
		// layout shows it
		f.lc = c.lineCtx(s, false)
		f.x0, f.w = scrollPad, s.textWidth()
		f.blocks(s.blocks, nil)
	}
	// the footnotes, then the endnotes, at the end
	if len(c.noteOrder) > 0 {
		f.separator()
		for _, kind := range []string{"f", "e"} {
			for _, n := range c.noteOrder {
				if n.kind == kind {
					f.blocks(n.blocks, nil)
				}
			}
		}
	}
	ops := f.all()
	total := f.y + scrollPad
	for _, o := range ops {
		total = math.Max(total, o.y1)
	}
	cuts := stripCuts(ops, total)
	W := width + 2*scrollPad
	stripOf := func(y float64) int {
		for i := 1; i < len(cuts); i++ {
			if y < cuts[i] {
				return i
			}
		}
		return len(cuts) - 1
	}
	vi := &viewInfo{link: func(b string) string {
		if y, ok := c.bmY[b]; ok {
			return "#page=" + strconv.Itoa(stripOf(y))
		}
		return ""
	}}
	var pages []*bdf.Page
	var cvs []*canvas.Canvas
	for i := 0; i+1 < len(cuts); i++ {
		a, b := cuts[i], cuts[i+1]
		cv := c.cvs.New()
		cv.Obj.SetBBox(0, 0, f32(W), f32(b-a))
		e := &emitter{c: c, cv: cv, view: vi}
		for _, o := range ops {
			in := o.y1 > a && o.y0 < b
			if o.text {
				in = o.y0 >= a-0.01 && o.y0 < b-0.01 || o.y0 == o.y1 && o.y0 >= a && o.y0 < b
			}
			if !in {
				continue
			}
			o.shift(0, -a)
			e.op(&o)
		}
		e.finish()
		pages = append(pages, v.AddPage(f32(W), f32(b-a)))
		cvs = append(cvs, cv)
	}
	*count = len(pages)
	return func() {
		for i, p := range pages {
			p.Layers = []bdf.Layer{{Role: bdf.RoleBody, Obj: cvs[i].Hash()}}
		}
	}
}

// stripCuts returns where the scroll view is cut into strips: about every
// stripHeight, moved up to a position no line of text crosses.
func stripCuts(ops []op, total float64) []float64 {
	cuts := []float64{0}
	y := 0.0
	for y+stripHeight*1.25 < total {
		c := safeCutOps(ops, y+stripHeight, y+stripHeight/2)
		if c <= y+1 {
			c = y + stripHeight
		}
		cuts = append(cuts, c)
		y = c
	}
	return append(cuts, total)
}

// safeCutOps returns the lowest position at or above y (and above min)
// that no text op crosses.
func safeCutOps(ops []op, y, min float64) float64 {
	for range 1000 {
		moved := false
		for _, o := range ops {
			if o.text && o.y0 < y-0.01 && o.y1 > y+0.01 {
				y = o.y0
				moved = true
			}
		}
		if !moved || y <= min {
			break
		}
	}
	return y
}

func boxAt(x, y, w, h float64) drawingml.Box { return drawingml.Box{X: x, Y: y, W: w, H: h} }

// paraBorderSides returns the top and bottom borders a paragraph draws:
// paragraphs with the same borders form one box, with the between border
// inside.
func (c *converter) paraBorderSides(p, prev *para, next block) (top, bottom *border) {
	b := p.pp.bdr
	if !sameBorders(p, prev) {
		top = b[bTop]
	} else {
		top = b[bBetween]
	}
	np, _ := next.(*para)
	if !sameBorders(p, np) {
		bottom = b[bBottom]
	}
	if top.none() {
		top = nil
	}
	if bottom.none() {
		bottom = nil
	}
	return
}

func sameBorders(a, b *para) bool {
	if a == nil || b == nil {
		return false
	}
	has := false
	for i := range 5 {
		x, y := a.pp.bdr[i], b.pp.bdr[i]
		if x.none() != y.none() || (!x.none() && *x != *y) {
			return false
		}
		has = has || !x.none()
	}
	al, ar, _ := a.indents()
	bl, br, _ := b.indents()
	return has && al == bl && ar == br
}

// paraBackground adds the shading and borders of a paragraph behind its
// lines (the ops from index from on).
func (f *flow) paraBackground(p *para, segs []paraSeg, from int, top, bottom *border) {
	pp := p.pp
	b := pp.bdr
	left, right := b[bLeft], b[bRight]
	if left.none() {
		left = nil
	}
	if right.none() {
		right = nil
	}
	if !pp.hasShd && top == nil && bottom == nil && left == nil && right == nil {
		return
	}
	indL, indR, _ := p.indents()
	for i, s := range segs {
		x0 := s.x0 + indL
		x1 := s.x0 + s.w - indR
		if left != nil {
			x0 -= left.space + left.width/2
		}
		if right != nil {
			x1 += right.space + right.width/2
		}
		y0, y1 := s.y0, s.y1
		var t, bt *border
		if i == 0 {
			t = top
		}
		if i == len(segs)-1 {
			bt = bottom
		}
		shd, hasShd := pp.shd, pp.hasShd
		fn := func(e *emitter, dx, dy float64) {
			if hasShd {
				e.fillRect(shd, x0+dx, y0+dy, x1-x0, y1-y0)
			}
			if t != nil {
				e.rule(t, x0+dx, y0+dy+t.width/2, x1+dx, y0+dy+t.width/2)
			}
			if bt != nil {
				e.rule(bt, x0+dx, y1+dy-bt.width/2, x1+dx, y1+dy-bt.width/2)
			}
			if left != nil {
				e.rule(left, x0+dx, y0+dy, x0+dx, y1+dy)
			}
			if right != nil {
				e.rule(right, x1+dx, y0+dy, x1+dx, y1+dy)
			}
		}
		o := op{y0: y0, y1: y1, fn: fn}
		// behind the lines of the segment
		ops := s.page.opsOf(f)
		at := from
		if s.page != segs[0].page {
			at = 0
		}
		at = min(at, len(*ops))
		*ops = slices.Insert(*ops, at, o)
	}
}

// opsOf returns the body ops of a page (the flow's ops when unpaged).
func (pg *page) opsOf(f *flow) *[]op {
	if pg == nil || !f.paged {
		return &f.ops
	}
	return &pg.body
}
