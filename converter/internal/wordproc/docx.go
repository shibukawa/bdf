package wordproc

import (
	"fmt"
	"io"
	"strconv"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
)

// Reading WordprocessingML: the document, its settings, sections,
// footnotes and endnotes, headers and footers.

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

// supported are the markup compatibility namespaces whose content is
// drawn (their mc:Choice is taken over the fallback).
var supported = map[string]bool{"wps": true, "wpg": true, "wpc": true, "wp14": true, "w14": true, "w15": true, "a14": true, "pic14": true}

// ConvertDOCX converts a Word document read from r.
func ConvertDOCX(r io.ReaderAt, size int64, opts *Options) (res *Result, err error) {
	defer func() {
		if e := recover(); e != nil {
			res, err = nil, fmt.Errorf("docx: internal error: %v", e)
		}
	}()
	c, views, err := newConverter(opts, ViewsBoth)
	if err != nil {
		return nil, fmt.Errorf("docx: %w", err)
	}
	opts = c.opts
	p, err := ooxml.Open(r, size)
	if err != nil {
		return nil, fmt.Errorf("docx: %w", err)
	}
	p.Supported = func(prefix string) bool { return supported[prefix] }
	c.pkg = p
	c.r = drawingml.New(drawingml.Config{Package: p, Doc: c.doc, Fonts: c.fonts, Images: opts.Images, Warn: c.warn})
	c.checkFonts()
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
	return c.finish(views), nil
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
