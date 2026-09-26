package wordproc

import (
	"fmt"
	"math"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// HTMLDocument is an HTML document to lay out in reader mode: its elements
// take the formatting of a fixed style sheet (see sheet), not the page's
// CSS.
type HTMLDocument struct {
	// Body holds the content: a parsed document, its body, or the element
	// that holds the article.
	Body *html.Node
	// DC is the document's metadata. Its language is that of the text
	// whose elements do not say.
	DC bdf.DublinCore
	// Source names the input format in the manifest ("html", "markdown").
	Source string
	// Width is the text width of the scroll view in points (0: 36 ems).
	Width float64
	// FontSize is the size of the text in points (0: 12).
	FontSize float64
	// Font and MonoFont are the font families of the text and of code (""
	// picks installed ones).
	Font, MonoFont string
	// Image returns the bytes of an image by its src attribute as written;
	// nil draws no images.
	Image func(src string) ([]byte, error)
	// Link resolves the href of a link to an absolute http:, https: or
	// mailto: URL, or "" to leave the link out; nil keeps absolute URLs.
	// Links to elements of the document itself ("#id") are resolved by the
	// layout.
	Link func(href string) string
}

// The reader style sheet. Lengths are in ems of the text size.
var (
	textColor     = bdf.RGB(0x1f, 0x23, 0x28)
	mutedColor    = bdf.RGB(0x59, 0x63, 0x6e)
	linkColor     = bdf.RGB(0x09, 0x69, 0xda)
	ruleColor     = bdf.RGB(0xd1, 0xd9, 0xe0)
	codeBlockBG   = bdf.RGB(0xf6, 0xf8, 0xfa)
	codeBG        = bdf.RGB(0xef, 0xf1, 0xf3)
	markBG        = bdf.RGB(0xff, 0xf0, 0x9a)
	headerCellBG  = bdf.RGB(0xf6, 0xf8, 0xfa)
	headingSizes  = [7]float64{0, 1.75, 1.45, 1.2, 1.05, 1, 0.9}
	defaultWidth  = 36.0 // text width of the scroll view
	defaultSize   = 12.0 // points
	bodyLH        = 1.7
	headingLH     = 1.3
	codeLH        = 1.45
	blockGap      = 0.8  // above and below paragraphs, lists, tables, boxes
	headingBefore = 1.4  // above headings
	headingAfter  = 0.6  // below headings
	itemGap       = 0.25 // between list items
	listIndent    = 1.75 // per list level
	labelGap      = 0.4  // between a list item's label and its text
	codeSize      = 0.875
	quoteBar      = 0.25 // width of the bar left of a quotation
	quotePad      = 1.0
	codePad       = 0.8
	cellPadV      = 0.35
	cellPadH      = 0.6
	ruleWidth     = 0.1
	pxToPt        = 0.75
)

// monoFamilies are tried in order for code when no family is given.
var monoFamilies = []string{"SF Mono", "Menlo", "Consolas", "Cascadia Mono", "DejaVu Sans Mono", "Liberation Mono",
	"Noto Sans Mono", "Source Code Pro", "Courier New"}

// maxDepth is how deeply nested elements keep their structure; the text of
// deeper ones is read as that of their parent.
const maxDepth = 200

// maxIndent is how many levels of lists are indented.
const maxIndent = 8

// maxTables is how deeply tables (and boxes) may nest; deeper ones are read
// as paragraphs.
const maxTables = 12

// ConvertHTML lays out an HTML document.
func ConvertHTML(d *HTMLDocument, opts *Options) (res *Result, err error) {
	defer func() {
		if e := recover(); e != nil {
			res, err = nil, fmt.Errorf("html: internal error: %v", e)
		}
	}()
	c, views, err := newConverter(opts, ViewsScroll)
	if err != nil {
		return nil, fmt.Errorf("html: %w", err)
	}
	c.checkFonts()
	c.css, c.compatMode = true, 15
	size := d.FontSize
	if size <= 0 {
		size = defaultSize
	}
	c.scrollWidth = d.Width
	if c.scrollWidth <= 0 {
		c.scrollWidth = defaultWidth * size
	}
	c.baseSize = size
	c.doc.Meta.Source = d.Source
	if c.doc.Meta.Source == "" {
		c.doc.Meta.Source = "html"
	}
	c.doc.Meta.DC = d.DC
	if c.opts.Title != "" {
		c.doc.Meta.DC.Title = bdf.DCValues{c.opts.Title}
	}
	r := &htmlReader{c: c, d: d, size: size, images: map[string]*htmlImage{}}
	r.font = d.Font
	if r.font == "" {
		r.font = fontdb.Sans
	}
	r.mono = d.MonoFont
	if r.mono == "" {
		r.mono = c.monoFamily()
	}
	lang := d.DC.Language.First()
	root := &hstyle{rp: &rprops{fonts: [4]string{r.font, r.font, r.font, r.font}, sz: size, color: textColor, u: "none", scale: 1,
		lang: [3]string{lang, lang, ""}}, jc: "left", lh: bodyLH}
	var blocks []block
	if d.Body != nil {
		blocks = r.read(d.Body, root, false)
	}
	// A4 pages with margins of 2 cm; the scroll view is c.scrollWidth wide
	sec := &section{pgW: 595.3, pgH: 841.9, top: 56.7, bottom: 56.7, left: 56.7, right: 56.7, header: 28.35, footer: 28.35,
		typ: "nextPage", pgStart: -1, pgFmt: "decimal", hdr: map[string]string{}, ftr: map[string]string{}, blocks: blocks}
	sec.cols = []column{{0, sec.textWidth()}}
	c.sections = []*section{sec}
	return c.finish(views), nil
}

// monoFamily picks an installed family for code: a known one, else the
// first installed family whose name says it is monospaced.
func (c *converter) monoFamily() string {
	for _, f := range monoFamilies {
		if len(c.db.Family(f)) > 0 {
			return f
		}
	}
	best := ""
	for _, f := range c.db.Faces {
		if fontdb.Classify(f.Family) == fontdb.Mono && (best == "" || f.Family < best) {
			best = f.Family
		}
	}
	if best != "" {
		return best
	}
	return fontdb.Mono
}

// hstyle is what an element passes on to its content: the character
// formatting (inherited as in CSS) and the formatting of the paragraphs
// that the content makes.
type hstyle struct {
	rp        *rprops
	link      string
	pre       bool     // white space is kept (pre)
	jc        string   // text alignment
	indL      float64  // left indent in the area (lists, dd)
	lh        float64  // line height
	before    float64  // margins of the paragraphs of a p or heading (inline
	after     float64  // elements pass them on, blocks set their own)
	heading   int      // 1-6 in headings
	frames    []*frame // structure the content is in
	listDepth int
	figAlt    string // the caption of the figure the content is in
	tables    int    // tables and boxes the content is in
	depth     int

	run *runStyle // derived from rp and link, when first needed
}

// child returns a copy for an element's content.
func (s *hstyle) child() *hstyle {
	n := *s
	n.run = nil
	n.depth++
	return &n
}

// block returns a copy for the content of a block element, whose
// paragraphs have no margins of their own.
func (s *hstyle) block() *hstyle {
	n := s.child()
	n.before, n.after = 0, 0
	return n
}

// chars returns a copy whose character formatting may change.
func (s *hstyle) chars() *hstyle {
	n := s.child()
	n.rp = s.rp.clone()
	return n
}

func (s *hstyle) inFrames(f ...*frame) []*frame {
	return append(slices.Clip(s.frames), f...)
}

// htmlReader reads the elements of an HTML document into blocks.
type htmlReader struct {
	c          *converter
	d          *HTMLDocument
	size       float64 // of the text
	font, mono string
	images     map[string]*htmlImage

	out        *[]block   // where finished blocks go
	p          *para      // the paragraph being filled
	space      bool       // collapsed white space is pending
	spaceNL    bool       // … and it held a line break
	spaceSt    *hstyle    // the style of the pending space
	pendBefore float64    // margin of elements opened since the last block
	pendLabel  *label     // the label of the list item whose first paragraph is next
	bookmarks  []string   // ids of elements waiting for their paragraph
	check      *inlineObj // the last checkbox, which the space after it goes into
}

// read reads the content of n into blocks; self reads n itself.
func (r *htmlReader) read(n *html.Node, st *hstyle, self bool) []block {
	var out []block
	saved := *r
	r.out, r.p, r.space, r.pendBefore, r.pendLabel = &out, nil, false, 0, nil
	if self {
		r.node(n, st)
	} else {
		r.children(n, st)
	}
	r.flush()
	bms := r.bookmarks
	*r = saved
	r.bookmarks = bms
	return out
}

func (r *htmlReader) children(n *html.Node, st *hstyle) {
	for k := n.FirstChild; k != nil; k = k.NextSibling {
		r.node(k, st)
	}
}

// em is a length in ems of the text.
func (r *htmlReader) em(v float64) float64 { return v * r.size }

// style returns the run style of st.
func (r *htmlReader) style(st *hstyle) *runStyle {
	if st.run == nil {
		st.run = r.c.runStyle(st.rp, st.link, false)
	}
	return st.run
}

// para returns the paragraph being filled, starting one in st's formatting.
func (r *htmlReader) para(st *hstyle) *para {
	if r.p != nil {
		return r.p
	}
	pp := defaultPProps()
	pp.jc, pp.indL, pp.lh = st.jc, st.indL, st.lh
	pp.before, pp.after = st.before, st.after
	pp.keepNext = st.heading > 0
	pp.widowCtl = true
	pp.overflowPunct = false
	mark := r.style(st)
	p := &para{pp: pp, mark: mark, markFace: r.c.face(mark, 'x'), heading: st.heading, ctx: st.inFrames(),
		bookmarks: r.bookmarks, label: r.pendLabel}
	if p.label != nil {
		pp.indFirst = -r.em(labelGap)
	}
	r.bookmarks, r.pendLabel = nil, nil
	r.p = p
	return p
}

// flush ends the paragraph being filled.
func (r *htmlReader) flush() {
	p := r.p
	r.p, r.space, r.spaceNL = nil, false, false
	if p == nil {
		return
	}
	for len(p.items) > 0 {
		it := p.items[len(p.items)-1]
		if it.kind != kChar || !isHTMLSpace(it.r) {
			break
		}
		p.items = p.items[:len(p.items)-1]
	}
	if len(p.items) == 0 && len(p.anchors) == 0 {
		// an empty paragraph: its label and bookmarks go to the next
		r.bookmarks = append(p.bookmarks, r.bookmarks...)
		if p.label != nil && r.pendLabel == nil {
			r.pendLabel = p.label
		}
		return
	}
	(&walker{c: r.c}).finishParagraph(p)
	r.add(p)
}

// add appends a finished block, with the margins of the elements opened
// since the last one (CSS margins collapse through them).
func (r *htmlReader) add(b block) {
	switch b := b.(type) {
	case *para:
		b.pp.before = math.Max(b.pp.before, r.pendBefore)
	case *table:
		b.before = math.Max(b.before, r.pendBefore)
	case *rule:
		b.before = math.Max(b.before, r.pendBefore)
	}
	r.pendBefore = 0
	*r.out = append(*r.out, b)
}

// block reads a block-level element: what comes before and after it is in
// other paragraphs, and its margins collapse with those of what it holds.
func (r *htmlReader) block(before, after float64, fn func()) {
	r.flush()
	r.pendBefore = math.Max(r.pendBefore, before)
	n := len(*r.out)
	fn()
	r.flush()
	if len(*r.out) == n {
		r.pendBefore = math.Max(r.pendBefore, after)
		return
	}
	switch b := (*r.out)[len(*r.out)-1].(type) {
	case *para:
		b.pp.after = math.Max(b.pp.after, after)
	case *table:
		b.after = math.Max(b.after, after)
	case *rule:
		b.after = math.Max(b.after, after)
	}
}

func isHTMLSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '\f'
}

// text adds the text of a text node.
func (r *htmlReader) text(s string, st *hstyle) {
	if st.pre {
		r.preText(s, st)
		return
	}
	for _, ch := range s {
		if isHTMLSpace(ch) {
			if !r.space {
				r.space, r.spaceSt = true, st
			}
			r.spaceNL = r.spaceNL || ch == '\n'
			continue
		}
		p := r.para(st)
		if r.space {
			if n := len(p.items); n > 0 && p.items[n-1].kind != kBreak && (r.check == nil || p.items[n-1].obj != r.check) {
				// a line break of the source between two East Asian
				// characters is no space
				prev := p.items[n-1].r
				if !(r.spaceNL && fontdb.IsCJK(prev) && fontdb.IsCJK(ch)) {
					r.c.addChar(p, r.style(r.spaceSt), ' ')
				}
			}
			r.space, r.spaceNL = false, false
		}
		r.c.addChar(p, r.style(st), ch)
	}
}

// preText adds text whose white space is kept: line breaks break lines,
// tabs move to the next multiple of four columns.
func (r *htmlReader) preText(s string, st *hstyle) {
	p := r.para(st)
	sty := r.style(st)
	col := 0
	for i := len(p.items) - 1; i >= 0 && p.items[i].kind != kBreak; i-- {
		col++
	}
	for _, ch := range s {
		switch ch {
		case '\r':
			continue
		case '\n':
			p.items = append(p.items, item{kind: kBreak, r: '\n', st: sty})
			col = 0
		case '\t':
			for n := 4 - col%4; n > 0; n-- {
				r.c.addChar(p, sty, ' ')
				col++
			}
		default:
			r.c.addChar(p, sty, ch)
			col++
		}
	}
}

func attr(n *html.Node, name string) (string, bool) {
	for _, a := range n.Attr {
		if a.Namespace == "" && a.Key == name {
			return a.Val, true
		}
	}
	return "", false
}

func attrStr(n *html.Node, name string) string {
	v, _ := attr(n, name)
	return strings.TrimSpace(v)
}

// cssDecls reads the declarations of a style attribute.
func cssDecls(n *html.Node) map[string]string {
	s, ok := attr(n, "style")
	if !ok {
		return nil
	}
	m := map[string]string{}
	for _, d := range strings.Split(s, ";") {
		k, v, ok := strings.Cut(d, ":")
		if ok {
			m[strings.ToLower(strings.TrimSpace(k))] = strings.ToLower(strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(v), "!important")))
		}
	}
	return m
}

// hidden reports whether an element is not rendered.
func hidden(n *html.Node, css map[string]string) bool {
	if _, ok := attr(n, "hidden"); ok {
		return true
	}
	return css["display"] == "none" || css["visibility"] == "hidden"
}

// align reads the text alignment an element asks for ("" when none).
func align(n *html.Node, css map[string]string) string {
	a := css["text-align"]
	if a == "" {
		a = strings.ToLower(attrStr(n, "align"))
	}
	switch a {
	case "left", "start":
		return "left"
	case "right", "end":
		return "right"
	case "center", "middle":
		return "center"
	case "justify":
		return "both"
	}
	return ""
}

// node reads a node of the document.
func (r *htmlReader) node(n *html.Node, st *hstyle) {
	switch n.Type {
	case html.TextNode:
		r.text(n.Data, st)
		return
	case html.DocumentNode:
		r.children(n, st)
		return
	case html.ElementNode:
	default:
		return
	}
	css := cssDecls(n)
	if hidden(n, css) || skipped(n.DataAtom) {
		return
	}
	if st.depth > maxDepth {
		r.children(n, st)
		return
	}
	if id := attrStr(n, "id"); id != "" {
		r.bookmark(id)
	}
	if n.DataAtom == atom.A {
		if name := attrStr(n, "name"); name != "" {
			r.bookmark(name)
		}
	}
	if l, ok := attr(n, "lang"); ok {
		st = st.chars()
		st.rp.lang = [3]string{strings.TrimSpace(l), strings.TrimSpace(l), ""}
	}
	switch n.DataAtom {
	case atom.Svg:
		r.c.warnOnce("svg", "SVG images are not drawn")
		return
	case atom.Math:
		r.c.warnOnce("math", "equations are drawn as plain text")
		r.children(n, st)
	case atom.P:
		r.block(0, 0, func() {
			s := r.blockStyle(n, css, st)
			s.before, s.after = r.em(blockGap), r.em(blockGap)
			r.children(n, s)
		})
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
		level := int(n.Data[1] - '0')
		r.block(0, 0, func() {
			s := r.blockStyle(n, css, st).chars()
			s.rp.sz = st.rp.sz * headingSizes[level]
			s.rp.setToggle(tB, true)
			if level == 6 {
				s.rp.color = mutedColor
			}
			s.lh, s.heading = headingLH, level
			s.before, s.after = r.em(headingBefore), r.em(headingAfter)
			r.children(n, s)
		})
	case atom.Div, atom.Section, atom.Article, atom.Main, atom.Header, atom.Footer, atom.Nav, atom.Aside, atom.Hgroup,
		atom.Form, atom.Fieldset, atom.Body, atom.Html, atom.Center, atom.Search, atom.Details, atom.Legend:
		r.block(0, 0, func() {
			s := r.blockStyle(n, css, st)
			if n.DataAtom == atom.Center {
				s.jc = "center"
			}
			r.children(n, s)
		})
	case atom.Address:
		r.block(0, 0, func() {
			s := r.blockStyle(n, css, st).chars()
			s.rp.setToggle(tI, true)
			r.children(n, s)
		})
	case atom.Summary, atom.Dt:
		r.block(0, 0, func() {
			s := r.blockStyle(n, css, st).chars()
			s.rp.setToggle(tB, true)
			r.children(n, s)
		})
	case atom.Dl:
		r.block(r.em(blockGap), r.em(blockGap), func() { r.children(n, r.blockStyle(n, css, st)) })
	case atom.Dd:
		r.block(0, r.em(itemGap), func() {
			s := r.blockStyle(n, css, st)
			s.indL += r.em(listIndent)
			r.children(n, s)
		})
	case atom.Figure:
		r.block(r.em(blockGap), r.em(blockGap), func() {
			s := r.blockStyle(n, css, st)
			if s.jc == "left" {
				s.jc = "center"
			}
			s.figAlt = caption(n)
			r.children(n, s)
		})
	case atom.Figcaption, atom.Caption:
		r.block(r.em(0.4), 0, func() {
			s := r.blockStyle(n, css, st).chars()
			s.rp.sz *= 0.875
			s.rp.color = mutedColor
			s.jc = "center"
			r.children(n, s)
		})
	case atom.Blockquote:
		r.quote(n, st)
	case atom.Pre, atom.Listing, atom.Xmp, atom.Plaintext:
		r.code(n, st)
	case atom.Ul, atom.Ol, atom.Menu, atom.Dir:
		r.list(n, css, st)
	case atom.Li:
		// an item outside a list
		r.block(0, 0, func() {
			s := r.nest(st)
			r.item(n, bullet("", s.listDepth), s)
		})
	case atom.Table:
		r.table(n, css, st)
	case atom.Hr:
		r.flush()
		r.add(&rule{before: r.em(1.2), after: r.em(1.2), ind: st.indL, width: r.em(ruleWidth), color: ruleColor})
	case atom.Br:
		p := r.para(st)
		p.items = append(p.items, item{kind: kBreak, r: '\n', st: r.style(st)})
		r.space, r.spaceNL = false, false
	case atom.Wbr:
		if r.p != nil {
			r.c.addChar(r.p, r.style(st), '\u200b') // a break opportunity
		}
	case atom.Img:
		r.img(n, css, st)
	case atom.Input:
		if strings.EqualFold(attrStr(n, "type"), "checkbox") {
			_, checked := attr(n, "checked")
			r.checkbox(checked, st)
		}
	case atom.A:
		s := st.child()
		if href, ok := attr(n, "href"); ok {
			if link := r.link(href); link != "" {
				s = st.chars()
				s.link = link
				s.rp.color = linkColor
				s.rp.u = "single"
			}
		}
		r.children(n, s)
	case atom.B, atom.Strong:
		s := st.chars()
		s.rp.setToggle(tB, true)
		r.children(n, s)
	case atom.I, atom.Em, atom.Cite, atom.Var, atom.Dfn:
		s := st.chars()
		s.rp.setToggle(tI, true)
		r.children(n, s)
	case atom.U, atom.Ins:
		s := st.chars()
		s.rp.u = "single"
		r.children(n, s)
	case atom.S, atom.Strike, atom.Del:
		s := st.chars()
		s.rp.setToggle(tStrike, true)
		r.children(n, s)
	case atom.Small:
		s := st.chars()
		s.rp.sz *= 0.85
		r.children(n, s)
	case atom.Big:
		s := st.chars()
		s.rp.sz *= 1.2
		r.children(n, s)
	case atom.Sup, atom.Sub:
		s := st.chars()
		s.rp.vertAlign = "superscript"
		if n.DataAtom == atom.Sub {
			s.rp.vertAlign = "subscript"
		}
		r.children(n, s)
	case atom.Mark:
		s := st.chars()
		s.rp.shd, s.rp.hasShd = markBG, true
		r.children(n, s)
	case atom.Code, atom.Kbd, atom.Samp, atom.Tt:
		s := st.chars()
		s.rp.fonts = [4]string{r.mono, r.mono, r.mono, r.mono}
		if !st.pre {
			s.rp.sz *= codeSize
			if n.DataAtom != atom.Samp && n.DataAtom != atom.Tt {
				s.rp.shd, s.rp.hasShd = codeBG, true
			}
		}
		r.children(n, s)
	case atom.Q:
		r.text("“", st)
		r.children(n, st.child())
		r.text("”", st)
	default:
		// span, abbr, time, label, font, custom elements …: their content
		r.children(n, st.child())
	}
}

// skipped reports whether the elements of a kind are not drawn: what is
// not content (head, script, style, template), what needs a browser to run
// (embedded documents, media, forms) and the annotations of ruby.
func skipped(a atom.Atom) bool {
	switch a {
	case atom.Head, atom.Script, atom.Style, atom.Template, atom.Noscript, atom.Iframe, atom.Object, atom.Embed,
		atom.Video, atom.Audio, atom.Canvas, atom.Select, atom.Textarea, atom.Button, atom.Map, atom.Dialog,
		atom.Rt, atom.Rp, atom.Source, atom.Track, atom.Param, atom.Colgroup, atom.Col, atom.Frameset:
		return true
	}
	return false
}

// blockStyle returns the style for the content of a block element, with
// the alignment it asks for.
func (r *htmlReader) blockStyle(n *html.Node, css map[string]string, st *hstyle) *hstyle {
	s := st.block()
	if a := align(n, css); a != "" {
		s.jc = a
	}
	return s
}

// bookmark makes an element's id a target of links.
func (r *htmlReader) bookmark(id string) {
	if r.p != nil {
		r.p.bookmarks = append(r.p.bookmarks, id)
		return
	}
	r.bookmarks = append(r.bookmarks, id)
}

// link resolves the target of a link.
func (r *htmlReader) link(href string) string {
	href = strings.TrimSpace(href)
	if frag, ok := strings.CutPrefix(href, "#"); ok {
		if frag == "" {
			return ""
		}
		if u, err := url.PathUnescape(frag); err == nil {
			frag = u
		}
		return "#bm:" + frag
	}
	if r.d.Link != nil {
		return r.d.Link(href)
	}
	return safeURL(href)
}

// caption returns the text of a figure's caption.
func caption(n *html.Node) string {
	for k := n.FirstChild; k != nil; k = k.NextSibling {
		if k.Type == html.ElementNode && k.DataAtom == atom.Figcaption {
			return strings.Join(strings.Fields(textContent(k)), " ")
		}
	}
	return ""
}

func textContent(n *html.Node) string {
	var b strings.Builder
	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
		}
		for k := n.FirstChild; k != nil; k = k.NextSibling {
			walk(k)
		}
	}
	walk(n)
	return b.String()
}

// list reads a ul or ol: its items are paragraphs indented one level more,
// the first of each with the item's bullet or number as its label.
func (r *htmlReader) list(n *html.Node, css map[string]string, st *hstyle) {
	gap := 0.0
	if st.listDepth == 0 {
		gap = r.em(blockGap)
	}
	r.block(gap, gap, func() {
		s := r.nest(r.blockStyle(n, css, st))
		ordered := n.DataAtom == atom.Ol
		typ := attrStr(n, "type")
		if t := css["list-style-type"]; t != "" {
			typ = t
		}
		num := 1
		if v, err := strconv.Atoi(attrStr(n, "start")); err == nil {
			num = v
		}
		_, reversed := attr(n, "reversed")
		if reversed {
			if _, ok := attr(n, "start"); !ok {
				num = 0
				for k := n.FirstChild; k != nil; k = k.NextSibling {
					if k.Type == html.ElementNode && k.DataAtom == atom.Li {
						num++
					}
				}
			}
		}
		for k := n.FirstChild; k != nil; k = k.NextSibling {
			if k.Type != html.ElementNode || k.DataAtom != atom.Li {
				r.node(k, s)
				continue
			}
			if v, err := strconv.Atoi(attrStr(k, "value")); err == nil {
				num = v
			}
			text := ""
			switch {
			case taskItem(k):
			case ordered || isOrderedType(typ):
				text = listNumber(num, typ) + "."
			default:
				text = bullet(typ, s.listDepth)
			}
			r.item(k, text, s)
			if reversed {
				num--
			} else {
				num++
			}
		}
	})
}

// nest returns the style of the items of a list in st: one level deeper
// (lists deeper than maxIndent levels are not indented further).
func (r *htmlReader) nest(st *hstyle) *hstyle {
	s := st.block()
	s.listDepth++
	if s.listDepth <= maxIndent {
		s.indL += r.em(listIndent)
	}
	s.frames = st.inFrames(&frame{kind: fList})
	return s
}

// item reads a list item.
func (r *htmlReader) item(n *html.Node, label string, st *hstyle) {
	css := cssDecls(n)
	if hidden(n, css) {
		return
	}
	if id := attrStr(n, "id"); id != "" {
		r.bookmark(id)
	}
	r.block(r.em(itemGap), 0, func() {
		s := r.blockStyle(n, css, st)
		s.frames = st.inFrames(&frame{kind: fItem})
		if label != "" {
			r.pendLabel = r.label(label, s)
		}
		r.children(n, s)
		r.flush()
		r.pendLabel = nil
	})
}

// taskItem reports whether a list item starts with a checkbox (a task
// list): the checkbox stands in for its bullet.
func taskItem(n *html.Node) bool {
	for k := n.FirstChild; k != nil; k = k.NextSibling {
		switch {
		case k.Type == html.TextNode && strings.TrimSpace(k.Data) == "":
			continue
		case k.Type == html.ElementNode && k.DataAtom == atom.Input:
			return strings.EqualFold(attrStr(k, "type"), "checkbox")
		case k.Type == html.ElementNode && k.DataAtom == atom.P:
			return taskItem(k)
		}
		return false
	}
	return false
}

func isOrderedType(t string) bool {
	switch t {
	case "1", "a", "A", "i", "I", "decimal", "lower-alpha", "upper-alpha", "lower-latin", "upper-latin", "lower-roman", "upper-roman":
		return true
	}
	return false
}

// listNumber formats the number of an ordered list's item.
func listNumber(n int, typ string) string {
	switch typ {
	case "a", "lower-alpha", "lower-latin":
		return formatNumber(n, "lowerLetter")
	case "A", "upper-alpha", "upper-latin":
		return formatNumber(n, "upperLetter")
	case "i", "lower-roman":
		return formatNumber(n, "lowerRoman")
	case "I", "upper-roman":
		return formatNumber(n, "upperRoman")
	}
	return strconv.Itoa(n)
}

// bullet returns the bullet of an unordered list's items.
func bullet(typ string, depth int) string {
	switch typ {
	case "circle":
		return "◦"
	case "square":
		return "▪"
	case "none":
		return ""
	case "disc":
		return "•"
	}
	switch depth {
	case 1:
		return "•"
	case 2:
		return "◦"
	}
	return "▪"
}

// label builds the label of a list item in the item's formatting.
func (r *htmlReader) label(text string, st *hstyle) *label {
	rp := st.rp.clone()
	rp.u, rp.hasShd, rp.vertAlign = "none", false, ""
	rp.setToggle(tStrike, false)
	sty := r.c.runStyle(rp, "", false)
	lb := &label{jc: "right", suffix: "tab"}
	tmp := &para{}
	r.c.addText(tmp, sty, text)
	lb.items = tmp.items
	for _, it := range lb.items {
		lb.w += it.w
	}
	return lb
}

// quote reads a blockquote: a box with a bar on its left, its text muted.
func (r *htmlReader) quote(n *html.Node, st *hstyle) {
	s := st.chars()
	s.rp.color = mutedColor
	bar := &border{style: "single", width: r.em(quoteBar), color: ruleColor}
	r.box(n, s, [4]float64{0, r.em(quotePad), 0, 0}, func(ce *cell) { ce.borders[sLeft] = bar })
}

// code reads a pre: a shaded box of monospaced text whose white space is
// kept. Lines longer than the box wrap.
func (r *htmlReader) code(n *html.Node, st *hstyle) {
	s := st.chars()
	s.rp.fonts = [4]string{r.mono, r.mono, r.mono, r.mono}
	s.rp.sz *= codeSize
	s.pre, s.lh, s.jc = true, codeLH, "left"
	pad := r.em(codePad)
	r.box(n, s, [4]float64{pad, pad, pad, pad}, func(ce *cell) { ce.shd, ce.hasShd = codeBlockBG, true })
}

// box lays an element's content out in a box of one cell (a table without
// structure): with padding (top, left, bottom, right), the width of the
// column and the margins of a block; style sets its borders and shading.
func (r *htmlReader) box(n *html.Node, st *hstyle, pad [4]float64, style func(*cell)) {
	if st.tables >= maxTables {
		r.block(r.em(blockGap), r.em(blockGap), func() { r.children(n, st.block()) })
		return
	}
	r.flush()
	s := st.block()
	s.indL, s.listDepth = 0, 0
	s.tables++
	blocks := r.read(n, s, false)
	if len(blocks) == 0 {
		return
	}
	ce := &cell{blocks: blocks, span: 1, rowSpan: 1, vAlign: "top", mar: pad, ctx: st.inFrames()}
	style(ce)
	t := &table{jc: "left", ind: st.indL, rowBand: 1, colBand: 1, rows: []*row{{cells: []*cell{ce}}}, grid: []float64{0},
		before: r.em(blockGap), after: r.em(blockGap)}
	if b := ce.borders[sLeft]; b != nil {
		// the border is centered on the box's edge
		t.ind += b.width / 2
	}
	t.auto = autoLayout(t)
	t.auto.fill = true
	r.add(t)
}

// checkbox adds a checkbox of a task list, drawn as a square (with a check
// mark) the size of the text.
func (r *htmlReader) checkbox(checked bool, st *hstyle) {
	p := r.para(st)
	size := st.rp.sz * 0.85
	gap := st.rp.sz * 0.35
	o := &inlineObj{w: size, h: size, ext: [4]float64{0, 0, gap, -st.rp.sz * 0.12}}
	col := st.rp.color
	o.paint = func(e *emitter, box drawingml.Box) {
		lw := size * 0.08
		e.cv.Obj.Save()
		e.cv.Obj.StrokeColor(col)
		e.cv.Obj.Line(float32(lw), bdf.CapRound, bdf.JoinRound, 10)
		sq := &bdf.Path{}
		sq.Rect(f32(box.X+lw/2), f32(box.Y+lw/2), f32(box.W-lw), f32(box.H-lw))
		e.cv.Obj.StrokePath(e.cv.Obj.AddPath(sq))
		if checked {
			ck := &bdf.Path{}
			ck.MoveTo(f32(box.X+box.W*0.22), f32(box.Y+box.H*0.52)).
				LineTo(f32(box.X+box.W*0.42), f32(box.Y+box.H*0.72)).
				LineTo(f32(box.X+box.W*0.78), f32(box.Y+box.H*0.28))
			e.cv.Obj.Line(float32(lw*1.6), bdf.CapRound, bdf.JoinRound, 10)
			e.cv.Obj.StrokePath(e.cv.Obj.AddPath(ck))
		}
		e.cv.Obj.Restore()
	}
	p.items = append(p.items, item{kind: kObject, obj: o, st: r.style(st), w: o.w + o.ext[0] + o.ext[2]})
	r.space, r.check = false, o
}
