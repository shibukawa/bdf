// Package idml converts InDesign documents saved as IDML (InDesign Markup
// Language, the XML interchange format InDesign exports and opens) into
// BDF documents.
//
// An IDML package is a zip of XML parts: a design map naming the others,
// spreads holding pages and the items on them (rectangles, ovals,
// polygons, lines, text frames, groups and placed images), master spreads
// whose items pages inherit, stories holding the text that flows through
// chains of text frames, and resources (swatches, stroke styles, styles).
// Each page becomes a page of a fixed view, with the items of its master
// as a master layer and its own items as the body layer. Text is laid out
// here, with the fonts that are then embedded as subsets, by the DrawingML
// text engine the Office converters share: a story is described as a
// DrawingML text body and flows through its frames (drawingml.TextFlow).
// The binary .indd format is not read. See docs/design.md §3.32.
package idml

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"math"
	"net/url"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
	"github.com/shibukawa/bdf/converter/internal/tiff"
	"github.com/shibukawa/bdf/converter/internal/xmp"
	"github.com/shibukawa/bdf/image/imgconv"
	"github.com/shibukawa/bdf/internal/fontdb"
)

// mimetype is the media type the package's mimetype part names.
const mimetype = "application/vnd.adobe.indesign-idml-package"

// Options controls the conversion.
type Options struct {
	// Pages selects 1-based pages (see converter.Pages); nil converts
	// every page.
	Pages conv.Pages
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
	// fonts used for layout.
	SystemFonts bool
	// NoSubset embeds whole fonts instead of the glyphs in use.
	NoSubset bool
	// NoWOFF2 stores embedded fonts as TrueType/OpenType instead of WOFF2.
	NoWOFF2 bool
	// IgnoreFSType embeds fonts whose OS/2 fsType forbids embedding or
	// subsetting. Set it only when you hold the rights to embed the fonts.
	IgnoreFSType bool
	// NoTextIndex skips building the text index part.
	NoTextIndex bool
	// Files holds the images the document places as links, by file name
	// (see converter.Options.Files); Dir is a directory to look for them
	// in (converter.Options.Dir). Links are matched by their file name.
	Files fs.FS
	Dir   string
	// Warn receives non-fatal problems; when nil they are collected in Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc           *bdf.Document
	Warnings      []string
	Pages         int
	EmbeddedFonts int
}

// doc is an IDML package read into its parts.
type doc struct {
	pkg          *ooxml.Package
	swatches     map[string]*ooxml.Node // Color, Tint, Gradient, MixedInk by Self
	strokeStyles map[string]*ooxml.Node
	pstyles      map[string]*ooxml.Node
	cstyles      map[string]*ooxml.Node
	layerZ       map[string]int // Self → stacking order, 0 at the back
	layerHidden  map[string]bool
	stories      map[string]*ooxml.Node
	masters      map[string]*spread
	spreads      []*spread
	sections     []*ooxml.Node
	rtl          bool
	dc           bdf.DublinCore
	lang         string
}

// spread is a spread or a master spread: its pages and items.
type spread struct {
	n       *ooxml.Node
	self    string
	pages   []*ooxml.Node
	items   []*ooxml.Node
	master  string // the applied master of a master spread ("" for none)
	binding int    // BindingLocation: the pages before the spine
}

// pageInfo is a page of the document.
type pageInfo struct {
	n       *ooxml.Node
	sp      *spread
	idx     int           // index in its spread
	inv     canvas.Matrix // spread → page
	bounds  rect          // in page coordinates
	spreadB rect          // in spread coordinates
	number  string
	left    bool
}

// frame is a text frame placed on a page.
type frame struct {
	n    *ooxml.Node
	m    canvas.Matrix // frame → page
	page *pageInfo
	key  string
	body *drawingml.TextBody
}

type converter struct {
	opts     *Options
	d        *doc
	doc      *bdf.Document
	fonts    *fontset.Set
	cvs      *canvas.Builder
	dm       *drawingml.Drawing
	pages    []*pageInfo
	frames   map[string]*frame // by key
	images   map[string]*picture
	warnings []string
	warned   map[string]bool
	refs     conv.Options
}

// picture is a placed image loaded into the document.
type picture struct {
	hash bdf.Hash
	ok   bool
}

// ConvertFile converts an .idml file.
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

// Convert converts a document read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	c := &converter{opts: opts, doc: bdf.NewDocument(), warned: map[string]bool{}, frames: map[string]*frame{},
		images: map[string]*picture{}, refs: conv.Options{Files: opts.Files, Dir: opts.Dir}}
	p, err := ooxml.Open(r, size)
	if err != nil {
		return nil, fmt.Errorf("idml: %w", err)
	}
	if c.d, err = c.read(p); err != nil {
		return nil, fmt.Errorf("idml: %w", err)
	}
	warn := func(msg string) { c.warnf("%s", msg) }
	db := fontdb.New(opts.FontFS, opts.FontDirs, !opts.NoSystemFonts)
	c.fonts = fontset.New(db, warn)
	c.cvs = canvas.NewBuilder(c.doc, c.fonts)
	rd := drawingml.New(drawingml.Config{Package: p, Doc: c.doc, Fonts: c.fonts, Images: opts.Images, Warn: warn})
	c.dm = rd.NewDrawing("", nil, textHost{})
	if len(db.Faces) == 0 && !opts.SystemFonts {
		c.warnf("no fonts found; text is laid out with estimated metrics and not embedded")
	}
	c.doc.Meta.Source = "idml"
	c.doc.Meta.DC = c.d.dc
	if opts.Title != "" {
		c.doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	if len(c.doc.Meta.DC.Language) == 0 && c.d.lang != "" {
		c.doc.Meta.DC.Language = bdf.DCValues{c.d.lang}
	}
	view := c.doc.NewView("pages", bdf.ViewFixed, c.doc.Meta.DC.Title.First())
	if c.d.rtl {
		view.Direction = bdf.DirectionRTL
	}

	c.pages = c.listPages()
	if len(c.pages) == 0 {
		return nil, fmt.Errorf("idml: the document has no pages")
	}
	sel := opts.Pages.Numbers(len(c.pages))
	if sel == nil {
		for i := range c.pages {
			sel = append(sel, i+1)
		}
	}
	for _, n := range sel {
		if n < 1 || n > len(c.pages) {
			return nil, fmt.Errorf("idml: page %d out of range (1-%d)", n, len(c.pages))
		}
	}
	// the text of every chain of frames flows across pages: lay it out
	// before any page is drawn
	c.layoutStories()
	type pageRef struct {
		page   *bdf.Page
		layers []*canvas.Canvas
	}
	var pages []pageRef
	for _, n := range sel {
		pg := c.pages[n-1]
		page := view.AddPage(f32(pg.bounds.w()), f32(pg.bounds.h()))
		pr := pageRef{page: page}
		if cv := c.renderSafe(pg, true, "master of page "+strconv.Itoa(n)); cv.Drawn {
			page.Layers = append(page.Layers, bdf.Layer{Role: bdf.RoleMaster})
			pr.layers = append(pr.layers, cv)
		}
		if cv := c.renderSafe(pg, false, "page "+strconv.Itoa(n)); cv.Drawn {
			page.Layers = append(page.Layers, bdf.Layer{Role: bdf.RoleBody})
			pr.layers = append(pr.layers, cv)
		}
		pages = append(pages, pr)
	}
	res := &Result{Doc: c.doc, Pages: len(sel)}
	if !opts.SystemFonts {
		res.EmbeddedFonts = c.fonts.Embed(c.doc, fontset.EmbedOptions{NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType})
	}
	c.cvs.Encode()
	c.fonts.ReportMissing()
	for _, pr := range pages {
		for i, cv := range pr.layers {
			pr.page.Layers[i].Obj = cv.Hash()
		}
		if len(pr.page.Layers) == 0 {
			// an empty page still needs something to draw
			o := bdf.NewObject()
			o.SetBBox(0, 0, pr.page.W, pr.page.H)
			hash, _ := c.doc.AddObject(o)
			pr.page.Layers = append(pr.page.Layers, bdf.Layer{Role: bdf.RoleBody, Obj: hash})
		}
	}
	if !opts.NoTextIndex {
		if _, err := c.doc.BuildTextIndex(view); err != nil {
			c.warnf("text index: %v", err)
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

// textHost is the drawingml host of the text: nothing is inherited.
type textHost struct{}

func (textHost) Placeholder(*ooxml.Node, string) ([]drawingml.Inherited, bool) { return nil, true }
func (textHost) TextStyle(string) *ooxml.Node                                  { return nil }
func (textHost) Field(string) (string, bool)                                   { return "", false }
func (textHost) Link(string, ooxml.Rel) string                                 { return "" }

// maxParts bounds the spreads and stories read, so that a design map
// listing the same part over and over does not grow the document without
// end.
const maxParts = 10000

// read reads the parts the design map names.
func (c *converter) read(p *ooxml.Package) (*doc, error) {
	d := &doc{pkg: p, swatches: map[string]*ooxml.Node{}, strokeStyles: map[string]*ooxml.Node{}, pstyles: map[string]*ooxml.Node{},
		cstyles: map[string]*ooxml.Node{}, layerZ: map[string]int{}, layerHidden: map[string]bool{}, stories: map[string]*ooxml.Node{},
		masters: map[string]*spread{}}
	root, err := p.XML("designmap.xml")
	if err != nil {
		return nil, err
	}
	if root.Name != "Document" {
		return nil, fmt.Errorf("designmap.xml: not an InDesign document")
	}
	if b, err := p.Read("META-INF/metadata.xml"); err == nil {
		d.dc = xmp.DublinCore(b)
	}
	part := func(k *ooxml.Node) (*ooxml.Node, bool) {
		src := k.AttrStr("src", "")
		if src == "" {
			return nil, false
		}
		n, err := p.XML(src)
		if err != nil {
			c.warnf("%s: %v", src, err)
			return nil, false
		}
		return n, true
	}
	// resources first: spreads and stories refer to them
	for _, k := range root.Elements() {
		switch {
		case k.Name == "Graphic" && k.AttrStr("src", "") != "":
			if n, ok := part(k); ok {
				for _, e := range n.Elements() {
					self := e.AttrStr("Self", "")
					switch e.Name {
					case "Color", "Tint", "Gradient", "MixedInk", "Swatch":
						d.swatches[self] = e
					case "StrokeStyle", "DashedStrokeStyle", "DottedStrokeStyle", "StripedStrokeStyle":
						d.strokeStyles[self] = e
					}
				}
			}
		case k.Name == "Styles" && k.AttrStr("src", "") != "":
			if n, ok := part(k); ok {
				collectStyles(n.Child("RootParagraphStyleGroup"), "ParagraphStyle", "ParagraphStyleGroup", d.pstyles, 0)
				collectStyles(n.Child("RootCharacterStyleGroup"), "CharacterStyle", "CharacterStyleGroup", d.cstyles, 0)
			}
		case k.Name == "Preferences" && k.AttrStr("src", "") != "":
			if n, ok := part(k); ok {
				dp := n.Child("DocumentPreference")
				d.rtl = dp.AttrStr("PageBinding", "LeftToRight") == "RightToLeft"
			}
		case k.Name == "Layer":
			d.layerHidden[k.AttrStr("Self", "")] = !k.AttrBool("Visible", true)
		case k.Name == "Section":
			d.sections = append(d.sections, k)
		}
	}
	// layers are listed from the front: the last is at the back
	var layers []string
	for _, k := range root.Children("Layer") {
		layers = append(layers, k.AttrStr("Self", ""))
	}
	for i, self := range layers {
		d.layerZ[self] = len(layers) - 1 - i
	}
	for _, s := range d.pstyles {
		if strings.HasSuffix(s.AttrStr("Self", ""), "/[No paragraph style]") {
			d.lang = langTag(func() string { v, _ := prop(s, "AppliedLanguage"); return v }())
		}
	}
	parts := 0
	for _, k := range root.Elements() {
		if k.AttrStr("src", "") == "" {
			continue
		}
		switch k.Name {
		case "Spread", "MasterSpread":
			if parts++; parts > maxParts {
				return nil, fmt.Errorf("the document has more than %d spreads and stories", maxParts)
			}
			n, ok := part(k)
			if !ok {
				continue
			}
			sn := n.Child(k.Name)
			if sn == nil {
				continue
			}
			sp := &spread{n: sn, self: sn.AttrStr("Self", ""), binding: int(sn.AttrInt("BindingLocation", 1))}
			for _, e := range sn.Elements() {
				if e.Name == "Page" {
					sp.pages = append(sp.pages, e)
				} else if isItem(e) {
					sp.items = append(sp.items, e)
				}
			}
			sortByLayer(sp.items, d.layerZ)
			if k.Name == "MasterSpread" {
				if m := sn.AttrStr("AppliedMaster", "n"); m != "n" {
					sp.master = m
				}
				d.masters[sp.self] = sp
			} else {
				d.spreads = append(d.spreads, sp)
			}
		case "Story":
			if parts++; parts > maxParts {
				return nil, fmt.Errorf("the document has more than %d spreads and stories", maxParts)
			}
			if n, ok := part(k); ok {
				if st := n.Child("Story"); st != nil {
					d.stories[st.AttrStr("Self", "")] = st
				}
			}
		}
	}
	return d, nil
}

// collectStyles gathers the styles of a style group and its subgroups.
func collectStyles(g *ooxml.Node, style, group string, out map[string]*ooxml.Node, depth int) {
	if g == nil || depth > 32 {
		return
	}
	for _, k := range g.Elements() {
		switch k.Name {
		case style:
			out[k.AttrStr("Self", "")] = k
		case group:
			collectStyles(k, style, group, out, depth+1)
		}
	}
}

// isItem reports whether an element is a page item that is drawn.
func isItem(n *ooxml.Node) bool {
	switch n.Name {
	case "Rectangle", "Oval", "Polygon", "GraphicLine", "TextFrame", "Group", "Button", "MultiStateObject":
		return true
	}
	return false
}

// sortByLayer orders page items by their layers, the backmost first,
// keeping the order within a layer.
func sortByLayer(items []*ooxml.Node, z map[string]int) {
	slices.SortStableFunc(items, func(a, b *ooxml.Node) int {
		return z[a.AttrStr("ItemLayer", "")] - z[b.AttrStr("ItemLayer", "")]
	})
}

// listPages lists the pages of the document in order, numbered by its
// sections.
func (c *converter) listPages() []*pageInfo {
	var out []*pageInfo
	for _, sp := range c.d.spreads {
		for i, pn := range sp.pages {
			b, ok := bounds(pn)
			if !ok || b.w() <= 0 || b.h() <= 0 || b.w() > bdf.MaxPageSize || b.h() > bdf.MaxPageSize {
				c.warnf("page %s: no usable size", pn.AttrStr("Self", ""))
				continue
			}
			m := parseTransform(pn)
			pg := &pageInfo{n: pn, sp: sp, idx: i, inv: invert(m), bounds: b, left: i < sp.binding}
			x0, y0 := m.Apply(b.x0, b.y0)
			x1, y1 := m.Apply(b.x1, b.y1)
			pg.spreadB = rect{math.Min(x0, x1), math.Min(y0, y1), math.Max(x0, x1), math.Max(y0, y1)}
			out = append(out, pg)
		}
	}
	// numbering: each section starts its own, unless it continues
	n, style := 1, "Arabic"
	starts := map[string]*ooxml.Node{}
	for _, s := range c.d.sections {
		starts[s.AttrStr("PageStart", "")] = s
	}
	for _, pg := range out {
		if s := starts[pg.n.AttrStr("Self", "")]; s != nil {
			if !s.AttrBool("ContinueNumbering", false) {
				n = int(s.AttrInt("PageNumberStart", 1))
			}
			style = s.AttrStr("PageNumberStyle", "Arabic")
		}
		pg.number = formatNumber(n, style)
		n++
	}
	return out
}

// formatNumber formats a page number in a section's style.
func formatNumber(n int, style string) string {
	if n < 1 || n > 10000 {
		return strconv.Itoa(n)
	}
	switch style {
	case "UpperRoman":
		return roman(n)
	case "LowerRoman":
		return strings.ToLower(roman(n))
	case "UpperLetters":
		return letters(n)
	case "LowerLetters":
		return strings.ToLower(letters(n))
	case "SingleLeadingZeros":
		return fmt.Sprintf("%02d", n)
	case "DoubleLeadingZeros":
		return fmt.Sprintf("%03d", n)
	}
	return strconv.Itoa(n)
}

func roman(n int) string {
	var b strings.Builder
	for _, v := range []struct {
		n int
		s string
	}{{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"}, {100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"}, {10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"}} {
		for n >= v.n {
			b.WriteString(v.s)
			n -= v.n
		}
	}
	return b.String()
}

func letters(n int) string {
	// A..Z, then AA..ZZ as InDesign counts
	r := rune('A' + (n-1)%26)
	return strings.Repeat(string(r), (n-1)/26+1)
}

// placedItem is an item with the transform placing it on a page.
type placedItem struct {
	n *ooxml.Node
	m canvas.Matrix // item → page
}

// walk visits the items of a spread (or group) placed on a page: m maps
// their parent's coordinates to the page. Invisible items and those on
// hidden layers are skipped; groups are entered.
func (c *converter) walk(items []*ooxml.Node, m canvas.Matrix, depth int, fn func(it placedItem)) {
	if depth > 64 {
		return
	}
	for _, n := range items {
		if !n.AttrBool("Visible", true) {
			continue
		}
		if l, ok := n.Attr("ItemLayer"); ok && c.d.layerHidden[l] {
			continue
		}
		im := m.Mul(parseTransform(n))
		if n.Name == "Group" || n.Name == "Button" || n.Name == "MultiStateObject" {
			var kids []*ooxml.Node
			for _, k := range n.Elements() {
				if isItem(k) {
					kids = append(kids, k)
				}
			}
			if n.Name == "MultiStateObject" && len(kids) > 0 {
				// the first state alone shows
				kids = kids[:1]
			}
			c.walk(kids, im, depth+1, fn)
			continue
		}
		fn(placedItem{n: n, m: im})
	}
}

// onPage reports whether an item placed by m on a page touches the page
// (or, when its geometry is empty, is skipped).
func onPage(n *ooxml.Node, m canvas.Matrix, pg *pageInfo) bool {
	g := geometry(n)
	if len(g) == 0 {
		return false
	}
	b := bbox(g)
	r := emptyRect
	for _, p := range [][2]float64{{b.x0, b.y0}, {b.x1, b.y0}, {b.x0, b.y1}, {b.x1, b.y1}} {
		x, y := m.Apply(p[0], p[1])
		r.add(x, y)
	}
	pb := pg.bounds
	return r.x1 >= pb.x0 && r.x0 <= pb.x1 && r.y1 >= pb.y0 && r.y0 <= pb.y1
}

// masterPage returns the page of a master spread that a page takes its
// items from: the one on the same side of the spine, else the last.
func (c *converter) masterPage(ms *spread, pg *pageInfo) (*ooxml.Node, int) {
	if len(ms.pages) == 0 {
		return nil, -1
	}
	for i, mp := range ms.pages {
		if (i < ms.binding) == pg.left {
			return mp, i
		}
	}
	return ms.pages[len(ms.pages)-1], len(ms.pages) - 1
}

// masterItems visits the items a page inherits from its master (and the
// master's own master, farthest first), placed on the page. Items the
// page overrides are left out.
func (c *converter) masterItems(pg *pageInfo, fn func(it placedItem)) {
	overridden := map[string]bool{}
	c.walk(pg.sp.items, canvas.Identity, 0, func(it placedItem) {
		if o, ok := it.n.Attr("OverriddenMasterPageItem"); ok {
			overridden[o] = true
		}
	})
	var visit func(self string, depth int)
	visit = func(self string, depth int) {
		ms := c.d.masters[self]
		if ms == nil || depth > 16 {
			return
		}
		if ms.master != "" {
			visit(ms.master, depth+1)
		}
		mp, _ := c.masterPage(ms, pg)
		if mp == nil {
			return
		}
		mb, ok := bounds(mp)
		if !ok {
			return
		}
		inv := invert(parseTransform(mp))
		mpg := &pageInfo{bounds: mb}
		c.walk(ms.items, inv, 0, func(it placedItem) {
			if overridden[it.n.AttrStr("Self", "")] || !onPage(it.n, it.m, mpg) {
				return
			}
			fn(it)
		})
	}
	if m := pg.n.AttrStr("AppliedMaster", "n"); m != "n" {
		visit(m, 0)
	}
}

// pageItems visits the page's own items.
func (c *converter) pageItems(pg *pageInfo, fn func(it placedItem)) {
	c.walk(pg.sp.items, pg.inv, 0, func(it placedItem) {
		if onPage(it.n, it.m, pg) {
			fn(it)
		}
	})
}

// renderSafe draws the master items (master) or the page's own items of a
// page; a malformed page that trips the renderer becomes an empty layer
// and a warning instead of failing the whole conversion.
func (c *converter) renderSafe(pg *pageInfo, master bool, what string) (cv *canvas.Canvas) {
	cv = c.cvs.New()
	defer func() {
		if r := recover(); r != nil {
			c.warnf("%s: internal error: %v", what, r)
			cv = c.cvs.New()
		}
	}()
	draw := func(it placedItem) { c.drawItem(cv, it, pg, master) }
	if master {
		c.masterItems(pg, draw)
	} else {
		c.pageItems(pg, draw)
	}
	return cv
}

// frameKey identifies a text frame's layout: the frame itself, and for a
// master frame the page it is drawn on (its page number differs).
func frameKey(n *ooxml.Node, pg *pageInfo, master bool) string {
	if master {
		return n.AttrStr("Self", "") + "@" + pg.n.AttrStr("Self", "")
	}
	return n.AttrStr("Self", "")
}

// layoutStories lays the text of every story out through its chain of
// frames, in page order for the frames that start a chain.
func (c *converter) layoutStories() {
	var all []*frame
	for _, pg := range c.pages {
		collect := func(master bool) func(it placedItem) {
			return func(it placedItem) {
				if it.n.Name != "TextFrame" {
					return
				}
				key := frameKey(it.n, pg, master)
				if _, ok := c.frames[key]; ok {
					return // a frame over two pages belongs to the first
				}
				f := &frame{n: it.n, m: it.m, page: pg, key: key}
				c.frames[key] = f
				all = append(all, f)
			}
		}
		c.masterItems(pg, collect(true))
		c.pageItems(pg, collect(false))
	}
	// chains: from each frame without a previous one, along NextTextFrame
	byKey := map[string]*frame{}
	for _, f := range all {
		byKey[f.key] = f
	}
	next := func(f *frame) *frame {
		n := f.n.AttrStr("NextTextFrame", "n")
		if n == "n" {
			return nil
		}
		if strings.Contains(f.key, "@") {
			n += f.key[strings.Index(f.key, "@"):]
		}
		return byKey[n]
	}
	done := map[*frame]bool{}
	for pass := 0; pass < 2; pass++ {
		for _, f := range all {
			if done[f] {
				continue
			}
			// a chain starts where no frame leads (or, in the second
			// pass, anywhere: a ring of frames)
			if pass == 0 {
				if p := f.n.AttrStr("PreviousTextFrame", "n"); p != "n" {
					if strings.Contains(f.key, "@") {
						p += f.key[strings.Index(f.key, "@"):]
					}
					if _, ok := byKey[p]; ok {
						continue
					}
				}
			}
			var chain []*frame
			for g := f; g != nil && !done[g] && len(chain) < 100000; g = next(g) {
				done[g] = true
				chain = append(chain, g)
			}
			c.layoutChain(chain)
		}
	}
}

// layoutChain flows a story through its frames.
func (c *converter) layoutChain(chain []*frame) {
	story := c.d.stories[chain[0].n.AttrStr("ParentStory", "")]
	if story == nil {
		return
	}
	st := c.storyBody(story, chain[0].page.number)
	flow := c.dm.NewTextFlow(st.body, "")
	for _, f := range chain {
		if flow.Done() {
			break
		}
		g := geometry(f.n)
		if len(g) == 0 {
			continue
		}
		box := anchorBox(g)
		if box.w() <= 0 || box.h() <= 0 {
			continue
		}
		pref := f.n.Child("TextFramePreference")
		bodyPr := frameBodyPr(pref, st.vertical)
		if insetsFill(bodyPr, box) {
			// insets larger than the frame leave no room: InDesign shows
			// nothing in such a frame (its text is overset)
			continue
		}
		f.body = flow.Fill(bodyPr, box.x0, box.y0, box.w(), box.h(), f.m)
	}
	if !flow.Done() {
		c.warnOnce("overset:"+story.AttrStr("Self", ""), "a story has more text than its frames hold (overset text is not shown)")
	}
}

// insetsFill reports whether the insets of a frame's body properties take
// all of the box, leaving no room for text.
func insetsFill(bodyPr *ooxml.Node, box rect) bool {
	ins := func(name string) float64 { return bodyPr.AttrEMU(name, 0) }
	w, h := box.w(), box.h()
	if v := bodyPr.AttrStr("vert", "horz"); v != "horz" {
		w, h = h, w
	}
	return w-ins("lIns")-ins("rIns") <= 0 || h-ins("tIns")-ins("bIns") <= 0
}

// drawItem draws a page item placed on a page.
func (c *converter) drawItem(cv *canvas.Canvas, it placedItem, pg *pageInfo, master bool) {
	n := it.n
	g := geometry(n)
	if len(g) == 0 {
		return
	}
	alpha := opacity(n.Child("TransparencySetting"))
	fillAlpha := alpha * opacity(n.Child("FillTransparencySetting"))
	if bm := n.Path("TransparencySetting", "BlendingSetting").AttrStr("BlendMode", "Normal"); bm != "Normal" {
		c.warnOnce("blend:"+bm, "the %s blend mode is drawn as normal", bm)
	}
	box := bbox(g)
	var p *bdf.Path
	if n.Name == "GraphicLine" {
		p = buildPath(g, false)
	} else if r := n.AttrFloat("CornerRadius", 0); r > 0 && isRect(g) && allRounded(n) {
		p = roundedRect(anchorBox(g), r)
	} else {
		p = buildPath(g, true)
	}
	fill := paintOf{}
	if n.Name != "GraphicLine" {
		fill = c.paint(n.AttrStr("FillColor", "Swatch/None"), n.AttrFloat("FillTint", -1), fillAlpha)
	}
	stroke, hasStroke := c.stroke(n)
	cv.Obj.Save()
	cv.Transform(it.m)
	ref := cv.Obj.AddPath(p)
	if hasStroke && stroke.align == "OutsideAlignment" {
		// the outer half of a doubled line shows beyond the fill
		stroke.width *= 2
		c.strokePath(cv, ref, stroke)
		stroke.width /= 2
	}
	if fill.ok() {
		if fill.grad != nil {
			cv.Obj.FillPaint(cv.Obj.AddPaint(c.gradientPaint(fill, n, box, "GradientFill")))
		} else {
			cv.Obj.FillColor(fill.solid.bdf())
		}
		cv.Obj.FillPath(ref, bdf.NonZero)
		cv.Drawn = true
	}
	// a placed graphic is clipped by the frame
	for _, k := range n.Elements() {
		switch k.Name {
		case "Image":
			c.drawImage(cv, k, ref)
		case "EPS", "PDF", "PICT", "WMF", "ImportedPage":
			c.warnOnce("placed:"+k.Name, "placed %s files are not drawn", k.Name)
		}
	}
	switch {
	case !hasStroke:
	case stroke.align == "InsideAlignment":
		cv.Obj.Save()
		cv.Obj.ClipPath(ref, bdf.NonZero)
		stroke.width *= 2
		c.strokePath(cv, ref, stroke)
		cv.Obj.Restore()
	case stroke.align == "OutsideAlignment":
	default:
		c.strokePath(cv, ref, stroke)
	}
	cv.Obj.Restore()
	if n.Name == "TextFrame" {
		// a frame over two pages has its text laid out on the first of them
		if f := c.frames[frameKey(n, pg, master)]; f != nil && f.body != nil && f.page == pg {
			f.body.Draw(cv)
			cv.Drawn = true
		}
	}
	if n.AttrStr("LeftLineEnd", "None") != "None" || n.AttrStr("RightLineEnd", "None") != "None" {
		c.warnOnce("lineends", "arrowheads and other line ends are not drawn")
	}
}

// allRounded reports whether every corner of an item takes the rounded
// corner option.
func allRounded(n *ooxml.Node) bool {
	for _, k := range []string{"TopLeftCornerOption", "TopRightCornerOption", "BottomLeftCornerOption", "BottomRightCornerOption"} {
		if n.AttrStr(k, "None") != "RoundedCorner" {
			return false
		}
	}
	return true
}

// strokePath strokes a path with a stroke style.
func (c *converter) strokePath(cv *canvas.Canvas, ref bdf.PathRef, s strokeStyle) {
	col := c.average(s.paint)
	cv.Obj.StrokeColor(col.bdf())
	cv.Obj.Line(f32(s.width), s.cap, s.join, f32(s.miter))
	cv.Obj.Dash(s.dash, 0)
	cv.Obj.StrokePath(ref)
	cv.Obj.Dash(nil, 0)
	cv.Drawn = true
}

// drawImage draws a placed raster image inside its frame, whose path
// (already added to the object) clips it.
func (c *converter) drawImage(cv *canvas.Canvas, img *ooxml.Node, clip bdf.PathRef) {
	pc := c.picture(img)
	if !pc.ok {
		return
	}
	gb := img.Path("Properties", "GraphicBounds")
	x0, y0 := gb.AttrFloat("Left", 0), gb.AttrFloat("Top", 0)
	x1, y1 := gb.AttrFloat("Right", 0), gb.AttrFloat("Bottom", 0)
	if x1 <= x0 || y1 <= y0 {
		return
	}
	cv.Obj.Save()
	cv.Obj.ClipPath(clip, bdf.NonZero)
	cv.Transform(parseTransform(img))
	alpha := opacity(img.Child("TransparencySetting"))
	if alpha < 1 {
		cv.Obj.Alpha(f32(alpha))
	}
	cv.Obj.Image(cv.Image(pc.hash), f32(x0), f32(y0), f32(x1-x0), f32(y1-y0))
	cv.Obj.Restore()
	cv.Drawn = true
}

// picture loads a placed image: the data embedded in the document, or the
// file it links to, by name, among the files given with the input.
func (c *converter) picture(img *ooxml.Node) *picture {
	link := img.Child("Link")
	uri := link.AttrStr("LinkResourceURI", "")
	key := uri
	if key == "" {
		key = "\x00" + img.AttrStr("Self", "")
	}
	if pc, ok := c.images[key]; ok {
		return pc
	}
	pc := &picture{}
	c.images[key] = pc
	var data []byte
	// Contents holds the XMP packet of a linked image; only an embedded
	// image (StoredState Embedded) keeps its data there, encoded
	if contents := img.Path("Properties", "Contents"); contents != nil && (link == nil || link.AttrStr("StoredState", "") == "Embedded") {
		text := strings.TrimSpace(contents.Content())
		if !strings.HasPrefix(text, "<") {
			clean := strings.Map(func(r rune) rune {
				if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
					return -1
				}
				return r
			}, text)
			if b, err := base64.StdEncoding.DecodeString(clean); err == nil && len(b) > 0 {
				data = b
			} else if b, err := hex.DecodeString(clean); err == nil && len(b) > 0 {
				data = b
			}
		}
	}
	if data == nil && uri != "" {
		name := linkName(uri)
		b, _, err := c.refs.ReadRef("", name)
		if err != nil {
			c.warnOnce("link:"+uri, "linked image %s not found (pass the file with the document): its frame is left empty", name)
			return pc
		}
		data = b
	}
	if data == nil {
		c.warnOnce("nodata", "a placed image has neither embedded data nor a link: its frame is left empty")
		return pc
	}
	switch {
	case tiff.Sniff(data):
		pic, damaged, err := tiff.Picture(data, c.opts.Images)
		if err != nil {
			c.warnf("placed image: %v", err)
			return pc
		}
		if damaged {
			c.warnOnce("tiffdamaged", "a TIFF image's pixel data is damaged; what is missing is left blank")
		}
		data = pic
		if res, err := imgconv.Optimize(data, c.opts.Images); err == nil || res.Data != nil {
			data = res.Data
		}
	case imgconv.Sniff(data) == "":
		c.warnOnce("picfmt", "a placed image in an unsupported format is not drawn")
		return pc
	default:
		if res, err := imgconv.Optimize(data, c.opts.Images); err == nil || res.Data != nil {
			data = res.Data
		}
	}
	pc.hash, pc.ok = c.doc.AddImage(data), true
	return pc
}

// linkName returns the file name a link URI names ("file:/Users/me/Links/
// photo.jpg", "file:///C:/Links/photo.jpg").
func linkName(uri string) string {
	s := uri
	if u, err := url.Parse(uri); err == nil && u.Path != "" {
		s = u.Path
	} else if i := strings.Index(uri, ":"); i >= 0 {
		s = uri[i+1:]
	}
	if un, err := url.PathUnescape(s); err == nil {
		s = un
	}
	s = strings.ReplaceAll(s, `\`, "/")
	return path.Base(s)
}
