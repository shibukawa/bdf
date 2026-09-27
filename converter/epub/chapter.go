package epub

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/converter/internal/webdoc"
	"github.com/shibukawa/bdf/converter/internal/wordproc"
	"github.com/shibukawa/bdf/imgconv"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// chapter is a content document of the spine.
type chapter struct {
	path   string // in the container
	item   *item
	linear bool
	fixed  bool // pre-paginated

	body     *html.Node
	mode     int        // writing mode (modeUnset …)
	viewport [2]float64 // fixed layout: the viewport's width and height in CSS pixels
	pictures []picture
	text     bool // holds text
}

// picture is a picture of a chapter: an image by its address, or an svg
// element that draws more than one image.
type picture struct {
	src  string     // as rewritten: a path in the container or a data: URL
	svg  *html.Node // an svg element drawn as an SVG document
	alt  string
	w, h float64 // the size it is drawn at in CSS pixels (0: unknown)
}

// contentTypes are the media types of content documents.
var contentTypes = map[string]bool{"application/xhtml+xml": true, "text/html": true, "image/svg+xml": true,
	"application/x-dtbook+xml": true, "text/x-oeb1-document": true}

// spine lists the chapters in reading order: the linear items of the
// spine, then those that are not linear (reachable by links only).
func (c *converter) spine() error {
	fixed := c.props["rendition:layout"] == "pre-paginated" || strings.EqualFold(c.props["fixed-layout"], "true")
	var linear, other []*chapter
	for _, ref := range c.pub.pkg.Spine.Itemrefs {
		it := c.pub.items[ref.IDRef]
		if it == nil {
			c.warn(fmt.Sprintf("spine item %q is not in the manifest", ref.IDRef))
			continue
		}
		// a publication-specific media type falls back to a content document
		for n := 0; !contentTypes[it.MediaType] && !strings.HasPrefix(it.MediaType, "image/") && it.Fallback != "" && n < 8; n++ {
			if fb := c.pub.items[it.Fallback]; fb != nil {
				it = fb
			} else {
				break
			}
		}
		if c.pub.encrypted[c.pub.canonical(it.path)] {
			return ErrDRM
		}
		path := c.pub.canonical(it.path)
		if path == "" {
			c.warn(fmt.Sprintf("%s: not in the publication; left out", it.path))
			continue
		}
		if c.byPath[path] != nil {
			continue
		}
		ch := &chapter{path: path, item: it, linear: ref.Linear != "no", fixed: fixed}
		switch {
		case hasToken(ref.Properties, "rendition:layout-pre-paginated"):
			ch.fixed = true
		case hasToken(ref.Properties, "rendition:layout-reflowable"):
			ch.fixed = false
		}
		c.byPath[path] = ch
		if ch.linear {
			linear = append(linear, ch)
		} else {
			other = append(other, ch)
		}
	}
	c.chapters = append(linear, other...)
	if len(c.chapters) == 0 {
		return fmt.Errorf("epub: the spine is empty")
	}
	return nil
}

// load reads a chapter and prepares it for the layout: the references of
// its links, targets and pictures are rewritten (see rewrite), and the
// honored properties of its style sheets are put in its elements' style
// attributes.
func (c *converter) load(ch *chapter) error {
	mt := ch.item.MediaType
	if strings.HasPrefix(mt, "image/") && mt != "image/svg+xml" {
		// a picture in the spine: a page of it
		img := &html.Node{Type: html.ElementNode, Data: "img", DataAtom: atom.Img,
			Attr: []html.Attribute{{Key: "src", Val: ch.path}, {Key: "alt", Val: ""}}}
		body := &html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body}
		body.AppendChild(img)
		ch.body = body
		ch.pictures = []picture{{src: ch.path}}
		return nil
	}
	if !contentTypes[mt] {
		return fmt.Errorf("media type %q is not a content document", mt)
	}
	data, err := c.pub.read(ch.path)
	if err != nil {
		return err
	}
	var doc *html.Node
	if mt == "text/html" {
		doc, err = webdoc.Parse(data, "text/html")
	} else if doc, err = webdoc.ParseXHTML(data); err != nil {
		c.warn(fmt.Sprintf("%s: not well-formed XML (%v); read as HTML", ch.path, err))
		doc, err = webdoc.ParseHTML(data, "")
	}
	if err != nil {
		return err
	}
	var rules []cssRule
	c.rewrite(ch, doc, &rules)
	st := newStyler(rules)
	root, body := documentElement(doc), findElement(doc, "body")
	if root != nil && root.Namespace == "svg" {
		// an SVG content document: its svg in a body
		body = &html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body}
		doc.RemoveChild(root)
		body.AppendChild(root)
		doc.AppendChild(body)
		root = nil
	}
	// the honored properties of the style sheets, then the ids (which the
	// rules may name) as bookmarks
	var apply func(n *html.Node)
	apply = func(n *html.Node) {
		if n.Type == html.ElementNode {
			st.apply(n)
			targets(ch, n)
		}
		for k := n.FirstChild; k != nil; k = k.NextSibling {
			apply(k)
		}
	}
	if root != nil {
		st.apply(root)
	}
	ch.body = doc
	if body != nil {
		ch.body = body
		ch.mode = writingMode(body)
	}
	apply(ch.body)
	if ch.mode == modeUnset && root != nil {
		ch.mode = writingMode(root)
	}
	c.analyze(ch)
	return nil
}

// rewrite rewrites the references of a chapter: links to chapters link to
// their bookmarks (see targets; "#" + the escaped bookmark; links to files
// that are not laid out are dropped), and pictures refer to paths in the
// container. It also reads the chapter's style sheets (link and style
// elements) and its viewport.
func (c *converter) rewrite(ch *chapter, n *html.Node, rules *[]cssRule) {
	for k := n.FirstChild; k != nil; {
		if k.Type == html.ElementNode && isOPS(k, "switch") {
			// epub:switch: its default content in its place (the cases need
			// MathML and the like)
			var first *html.Node
			for d := k.FirstChild; d != nil; d = d.NextSibling {
				if d.Type != html.ElementNode || !isOPS(d, "default") {
					continue
				}
				for g := d.FirstChild; g != nil; {
					next := g.NextSibling
					d.RemoveChild(g)
					n.InsertBefore(g, k)
					if first == nil {
						first = g
					}
					g = next
				}
			}
			next := k.NextSibling
			n.RemoveChild(k)
			if first != nil {
				next = first
			}
			k = next
			continue
		}
		if k.Type == html.ElementNode {
			c.element(ch, k, rules)
			c.rewrite(ch, k, rules)
		}
		k = k.NextSibling
	}
}

func isOPS(n *html.Node, local string) bool {
	return n.Data == local && n.Namespace == "http://www.idpf.org/2007/ops" || n.Data == "epub:"+local
}

// targets makes the id of an element (and the name of an a element) a
// bookmark unique in the book: the chapter's path, "#", the id. The ids of
// SVG and MathML elements stay as they are: SVG refers to them ("#icon",
// "url(#gradient)").
func targets(ch *chapter, n *html.Node) {
	if n.Namespace != "" {
		return
	}
	for i := range n.Attr {
		if a := &n.Attr[i]; a.Namespace == "" && (a.Key == "id" || a.Key == "name" && n.DataAtom == atom.A) {
			a.Val = bookmark(ch.path, strings.TrimSpace(a.Val))
		}
	}
}

// element rewrites the references of one element.
func (c *converter) element(ch *chapter, n *html.Node, rules *[]cssRule) {
	switch {
	case n.Namespace == "svg" && n.Data == "image":
		for i := range n.Attr {
			if a := &n.Attr[i]; a.Key == "href" && (a.Namespace == "" || a.Namespace == "xlink") || a.Key == "xlink:href" {
				a.Val = c.pictureRef(ch.path, a.Val)
			}
		}
	case n.Namespace != "":
	case n.DataAtom == atom.A || n.DataAtom == atom.Area:
		for i := 0; i < len(n.Attr); i++ {
			if a := &n.Attr[i]; a.Namespace == "" && a.Key == "href" {
				if ref, ok := c.linkRef(ch.path, a.Val); ok {
					a.Val = ref
				} else {
					n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
					i--
				}
			}
		}
	case n.DataAtom == atom.Img:
		src := attrVal(n, "src")
		if src == "" {
			if f := strings.Fields(attrVal(n, "srcset")); len(f) > 0 {
				src = strings.TrimSuffix(f[0], ",")
			}
		}
		setAttr(n, "src", c.pictureRef(ch.path, src))
		removeAttr(n, "srcset")
	case n.DataAtom == atom.Link:
		if rel := strings.ToLower(attrVal(n, "rel")); hasToken(rel, "stylesheet") && !hasToken(rel, "alternate") {
			if p := c.pub.canonical(c.pub.resolve(ch.path, attrVal(n, "href"))); p != "" {
				c.styleSheet(p, c.css(p), rules)
			}
		}
	case n.DataAtom == atom.Style:
		c.styleSheet(ch.path, textOf(n), rules)
	case n.DataAtom == atom.Meta:
		if strings.EqualFold(attrVal(n, "name"), "viewport") {
			ch.viewport = parseViewport(attrVal(n, "content"))
		}
	}
}

// styleSheet reads the rules of a style sheet at path (its imports
// resolve against it).
func (c *converter) styleSheet(path, text string, rules *[]cssRule) {
	if text == "" {
		return
	}
	n := len(*rules)
	seen := map[string]bool{}
	parseCSS(text, func(u string) string {
		// a sheet is imported once
		p := c.pub.canonical(c.pub.resolve(path, u))
		if seen[p] {
			return ""
		}
		seen[p] = true
		return c.css(p)
	}, 0, rules)
	for _, r := range (*rules)[n:] {
		for _, d := range r.decls {
			if strings.HasSuffix(d[0], "writing-mode") {
				c.modes = true
			}
		}
	}
}

// css returns the text of a style sheet of the container.
func (c *converter) css(path string) string {
	if path == "" {
		return ""
	}
	if s, ok := c.sheets[path]; ok {
		return s
	}
	b, err := c.pub.read(path)
	s := ""
	if err == nil {
		s = strings.TrimPrefix(string(b), "\ufeff")
	}
	c.sheets[path] = s
	return s
}

// bookmark is the name of a target in the book: the chapter's path, and
// the id of the element in it.
func bookmark(path, id string) string {
	if id == "" {
		return path
	}
	return path + "#" + id
}

// linkRef rewrites the target of a link: to a bookmark for a chapter or an
// element in one, kept for absolute URLs; false drops the link.
func (c *converter) linkRef(base, href string) (string, bool) {
	u, err := url.Parse(strings.TrimSpace(href))
	if err != nil {
		return "", false
	}
	if u.Scheme != "" {
		return href, true
	}
	target := base
	if u.Path != "" {
		target = c.pub.canonical(c.pub.resolve(base, u.Path))
	}
	if c.byPath[target] == nil {
		return "", false
	}
	return "#" + url.PathEscape(bookmark(target, u.Fragment)), true
}

// pictureRef rewrites the address of a picture to its path in the
// container; data: URLs and absolute URLs stay as they are.
func (c *converter) pictureRef(base, src string) string {
	src = strings.TrimSpace(src)
	if src == "" || webdoc.IsDataURL(src) || strings.HasPrefix(src, "#") {
		return src
	}
	if p := c.pub.canonical(c.pub.resolve(base, src)); p != "" {
		return p
	}
	return src
}

// parseViewport reads the width and height of a viewport meta element
// ("width=1200, height=1600").
func parseViewport(s string) (vp [2]float64) {
	for _, f := range strings.FieldsFunc(s, func(r rune) bool { return r == ',' || r == ';' }) {
		k, v, ok := strings.Cut(f, "=")
		if !ok {
			continue
		}
		n, err := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(v), "px"), 64)
		if err != nil || n <= 0 {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(k)) {
		case "width":
			vp[0] = n
		case "height":
			vp[1] = n
		}
	}
	return vp
}

// analyze finds the pictures of a chapter and whether it holds text.
func (c *converter) analyze(ch *chapter) {
	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			if strings.TrimSpace(strings.ReplaceAll(n.Data, " ", " ")) != "" {
				ch.text = true
			}
			return
		case html.ElementNode:
			if hiddenElement(n) {
				return
			}
			switch {
			case n.Namespace == "svg" && n.Data == "svg":
				// a picture (the text in it is the picture's)
				if p, ok := svgPicture(n); ok {
					ch.pictures = append(ch.pictures, p)
				}
				return
			case n.DataAtom == atom.Img && n.Namespace == "":
				w, _ := strconv.ParseFloat(attrVal(n, "width"), 64)
				h, _ := strconv.ParseFloat(attrVal(n, "height"), 64)
				ch.pictures = append(ch.pictures, picture{src: attrVal(n, "src"), alt: attrVal(n, "alt"), w: w, h: h})
			case n.Namespace == "" && (n.DataAtom == atom.Head || n.DataAtom == atom.Script || n.DataAtom == atom.Style ||
				n.DataAtom == atom.Title || n.DataAtom == atom.Noscript || n.DataAtom == atom.Template):
				return
			}
		}
		for k := n.FirstChild; k != nil; k = k.NextSibling {
			walk(k)
		}
	}
	walk(ch.body)
}

// svgPicture reads the picture an svg element is: the one image it fits
// into its view box (the way EPUB covers and the pages of comics are
// written), drawn as that image, or else the svg itself, drawn as an SVG
// document. The size is the svg's (its view box's when it has no absolute
// one). False for an svg that draws nothing (a hidden sprite sheet).
func svgPicture(n *html.Node) (picture, bool) {
	if inlineStyle(n, "display") == "none" || zeroLength(attrVal(n, "width")) || zeroLength(attrVal(n, "height")) {
		return picture{}, false
	}
	var images []*html.Node
	other, title := false, ""
	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		for k := n.FirstChild; k != nil; k = k.NextSibling {
			if k.Type != html.ElementNode {
				continue
			}
			switch k.Data {
			case "image":
				images = append(images, k)
			case "g", "a", "switch":
				walk(k)
			case "title":
				if title == "" {
					title = strings.Join(strings.Fields(textOf(k)), " ")
				}
			case "desc", "metadata", "defs", "style", "script", "symbol", "linearGradient", "radialGradient", "pattern",
				"clipPath", "mask", "filter", "marker":
			default:
				other = true
			}
		}
	}
	walk(n)
	if len(images) == 0 && !other {
		return picture{}, false
	}
	p := picture{alt: attrVal(n, "aria-label")}
	if p.alt == "" {
		p.alt = title
	}
	p.w, p.h = imgconv.ParseSVGSize(attrVal(n, "width"), attrVal(n, "height"), attrVal(n, "viewBox")).Pixels()
	if len(images) == 1 && !other {
		p.src = svgHref(images[0])
	}
	if p.src == "" {
		p.svg = n
	}
	return p, true
}

// zeroLength reports a length written as zero ("0", "0px").
func zeroLength(v string) bool {
	f, err := strconv.ParseFloat(strings.TrimRight(strings.TrimSpace(v), "abcdefghijklmnopqrstuvwxyz%"), 64)
	return err == nil && f == 0
}

// svgHref returns the address of an SVG image element: href, else
// xlink:href.
func svgHref(n *html.Node) string {
	xlink := ""
	for _, a := range n.Attr {
		switch {
		case a.Key == "href" && a.Namespace == "":
			return strings.TrimSpace(a.Val)
		case a.Key == "href" && a.Namespace == "xlink", a.Key == "xlink:href":
			xlink = strings.TrimSpace(a.Val)
		}
	}
	return xlink
}

// hiddenElement reports whether an element is not rendered.
func hiddenElement(n *html.Node) bool {
	if n.Namespace != "" {
		return false
	}
	for _, a := range n.Attr {
		if a.Namespace == "" && a.Key == "hidden" {
			return true
		}
	}
	return inlineStyle(n, "display") == "none"
}

// singlePicture returns a chapter's picture when the chapter is one
// picture and nothing else.
func (ch *chapter) singlePicture() (picture, bool) {
	if ch.text || len(ch.pictures) != 1 {
		return picture{}, false
	}
	p := ch.pictures[0]
	return p, p.src != "" || p.svg != nil
}

// pictureOnly reports whether a chapter holds pictures and no text (a
// cover, an illustration): it is laid out centered and without a page
// number.
func (ch *chapter) pictureOnly() bool {
	return !ch.text && len(ch.pictures) > 0
}

// cjk reports whether a language is written in East Asian scripts.
func cjk(lang string) bool {
	l := strings.ToLower(lang)
	for _, p := range []string{"ja", "zh", "ko"} {
		if l == p || strings.HasPrefix(l, p+"-") {
			return true
		}
	}
	return false
}

// reflow lays the book out in reader mode.
func (c *converter) reflow() (*Result, error) {
	opts := c.opts
	rtl := strings.EqualFold(c.pub.pkg.Spine.PPD, "rtl")
	// the book's writing mode, for chapters whose style does not say: that
	// of Kindle's primary-writing-mode, else vertical for East Asian books
	// bound on the right whose style sheets say no writing mode at all
	def := modeHorizontal
	switch c.props["primary-writing-mode"] {
	case "vertical-rl":
		def = modeVertical
	case "":
		if rtl && !c.modes && cjk(c.dc.Language.First()) {
			def = modeVertical
		}
	}
	fixed := 0
	res := &Result{}
	d := &wordproc.HTMLDocument{DC: c.dc, Source: "epub", Width: opts.Width, FontSize: opts.FontSize, Font: opts.Font,
		MonoFont: opts.MonoFont, Link: link, Folios: true,
		Image: func(src string) ([]byte, error) {
			b, err := c.image(src)
			if err == nil {
				res.Images++
			}
			return b, err
		}}
	paper := opts.Paper
	if paper.Width <= 0 || paper.Height <= 0 {
		paper = A5
	}
	d.Page = wordproc.PageSetup{Width: paper.Width, Height: paper.Height, Left: paper.Width * 0.1, Right: paper.Width * 0.1,
		Top: paper.Height * 0.09, Bottom: paper.Height * 0.09}
	for _, ch := range c.chapters {
		if ch.fixed {
			fixed++
		}
		mode := ch.mode
		if mode == modeUnset {
			mode = def
		}
		if mode == modeVerticalLR {
			c.warnOnce("vertical-lr", "vertical-lr text (Mongolian) is laid out horizontally")
		}
		hc := wordproc.HTMLChapter{Body: ch.body, ID: ch.path, Vertical: mode == modeVertical}
		if ch.pictureOnly() {
			// a cover or an illustration: centered, horizontal, no page
			// number; one picture fills the page
			hc.Vertical, hc.NoFolio = false, true
			_, hc.Picture = ch.singlePicture()
			if ch.body.Type == html.ElementNode && ch.body.Namespace == "" {
				setAttr(ch.body, "style", attrVal(ch.body, "style")+"; text-align: center")
			}
		}
		if hc.Vertical {
			res.Vertical = true
		}
		d.Chapters = append(d.Chapters, hc)
	}
	if fixed > 0 {
		c.warn(fmt.Sprintf("%d fixed-layout page(s) hold more than a picture; the book is laid out as reflowable text", fixed))
	}
	views := opts.Views
	if views == "" {
		views = ViewsPages
	}
	wo := &wordproc.Options{Pages: opts.Pages, Views: views, Title: opts.Title, Images: opts.Images,
		FontFS: opts.FontFS, FontDirs: opts.FontDirs, NoSystemFonts: opts.NoSystemFonts, SystemFonts: !opts.EmbedFonts,
		NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType, NoTextIndex: opts.NoTextIndex,
		Warn: opts.Warn}
	r, err := wordproc.ConvertHTML(d, wo)
	if err != nil {
		return nil, fmt.Errorf("epub: %w", err)
	}
	res.Doc, res.Pages, res.Strips, res.EmbeddedFonts = r.Doc, r.Pages, r.Strips, r.EmbeddedFonts
	res.Chapters = len(c.chapters)
	res.Warnings = append(c.warnings, r.Warnings...)
	return res, nil
}

// documentElement returns the root element of a document.
func documentElement(doc *html.Node) *html.Node {
	for k := doc.FirstChild; k != nil; k = k.NextSibling {
		if k.Type == html.ElementNode {
			return k
		}
	}
	return nil
}

// findElement returns the first HTML element of a kind, in document order.
func findElement(n *html.Node, tag string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tag && n.Namespace == "" {
		return n
	}
	for k := n.FirstChild; k != nil; k = k.NextSibling {
		if f := findElement(k, tag); f != nil {
			return f
		}
	}
	return nil
}

func attrVal(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Namespace == "" && a.Key == key {
			return strings.TrimSpace(a.Val)
		}
	}
	return ""
}

func setAttr(n *html.Node, key, val string) {
	for i, a := range n.Attr {
		if a.Namespace == "" && a.Key == key {
			n.Attr[i].Val = val
			return
		}
	}
	n.Attr = append(n.Attr, html.Attribute{Key: key, Val: val})
}

func removeAttr(n *html.Node, key string) {
	for i, a := range n.Attr {
		if a.Namespace == "" && a.Key == key {
			n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
			return
		}
	}
}

func textOf(n *html.Node) string {
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
