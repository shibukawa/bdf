// Package drawio converts draw.io (diagrams.net) diagrams into BDF documents.
//
// Every page of the diagram becomes a view of its own, so that viewers
// switch between pages the way spreadsheet viewers switch between sheets.
// A view holds one fixed page the size of the drawing (plus a border), with
// a layer object per draw.io layer. Shapes, edges and labels are drawn from
// the mxGraphModel XML the way draw.io itself renders them: the cell
// geometry, edge routing, shapes and label layout follow mxGraph and
// draw.io's Graph and Shapes code, ported from JavaScript (Apache License
// 2.0, see NOTICE). See docs/design.md §3.5.
package drawio

import (
	"fmt"
	"math"
	"net/url"
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
	// Pages selects 1-based page numbers; nil converts every page.
	Pages []int
	// Title overrides the document title.
	Title string
	// Border is the margin around the drawing in diagram units (pixels);
	// negative values mean none. Zero selects the default (10).
	Border float64
	// Images controls whether raster images are re-encoded (see imgconv).
	// The zero value keeps images as they are.
	Images imgconv.Options
	// FontDirs are searched for fonts before the system font directories.
	FontDirs []string
	// NoSystemFonts restricts font lookup to FontDirs.
	NoSystemFonts bool
	// SystemFonts refers to fonts by family name instead of embedding the
	// fonts used for layout.
	SystemFonts bool
	// NoSubset embeds whole fonts instead of the glyphs in use.
	NoSubset bool
	// NoWOFF2 stores embedded fonts as TrueType/OpenType instead of WOFF2.
	NoWOFF2 bool
	// IgnoreFSType embeds fonts whose OS/2 fsType forbids embedding or
	// subsetting (only with the rights to do so).
	IgnoreFSType bool
	// NoTextIndex skips building the text index parts.
	NoTextIndex bool
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

// ptPerUnit converts diagram units (CSS pixels) to points.
const ptPerUnit = 0.75

const defaultBorder = 10

type converter struct {
	opts     *Options
	doc      *bdf.Document
	db       *fontdb.DB
	warnings []string
	warned   map[string]bool

	canvases      []*canvas
	faceRunes     map[*fontdb.Face]map[rune]bool
	choices       map[resolveKey]*faceChoice
	fallback      map[fallbackKey]*faceChoice
	fallbackLists map[fallbackListKey][]fontdb.Resolved
	missing       map[rune]bool
	images        map[string]*imageRef
	embeddedFonts int

	// viewOf maps diagram ids to the ids of the views made from them, for
	// links between pages.
	viewOf map[string]string
	// awsLegacy counts the shapes of older AWS icon sets drawn with their
	// current counterparts (aws_legacy.go).
	awsLegacy int

	// per page
	m          *model
	page       *pageInfo
	pageShadow bool
}

func newConverter(opts *Options) *converter {
	c := &converter{opts: opts, doc: bdf.NewDocument(), warned: map[string]bool{},
		faceRunes: map[*fontdb.Face]map[rune]bool{}, choices: map[resolveKey]*faceChoice{},
		fallback: map[fallbackKey]*faceChoice{}, fallbackLists: map[fallbackListKey][]fontdb.Resolved{},
		missing: map[rune]bool{}, images: map[string]*imageRef{}, viewOf: map[string]string{}}
	c.db = fontdb.New(opts.FontDirs, !opts.NoSystemFonts)
	if len(c.db.Faces) == 0 && !opts.SystemFonts {
		c.warnf("no fonts found; text is laid out with estimated metrics and not embedded")
	}
	return c
}

// ConvertFile converts a .drawio (or .drawio.svg, .drawio.png, .xml) file.
func ConvertFile(path string, opts *Options) (*Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Convert(data, opts)
}

// Convert converts draw.io data.
func Convert(data []byte, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	f, err := readFile(data)
	if err != nil {
		return nil, fmt.Errorf("drawio: %w", err)
	}
	c := newConverter(opts)
	c.doc.Meta.Source = "drawio"
	if opts.Title != "" {
		c.doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	if f.modified != "" {
		c.doc.Meta.DC.Modified = bdf.DCValues{f.modified}
	}

	sel := opts.Pages
	if sel == nil {
		for i := range f.pages {
			sel = append(sel, i+1)
		}
	}
	for _, n := range sel {
		if n < 1 || n > len(f.pages) {
			return nil, fmt.Errorf("drawio: page %d out of range (1-%d)", n, len(f.pages))
		}
	}
	// view ids: the diagram's id when it is usable, else pageN
	used := map[string]bool{}
	ids := make([]string, len(sel))
	for i, n := range sel {
		p := f.pages[n-1]
		id := p.id
		if id == "" || used[id] {
			id = "page" + strconv.Itoa(n)
			for used[id] {
				id += "_"
			}
		}
		used[id] = true
		ids[i] = id
		if p.id != "" {
			c.viewOf[p.id] = id
		}
	}
	type pending struct {
		view   *bdf.View
		page   *bdf.Page
		layers []*canvas
	}
	var views []pending
	for i, n := range sel {
		p := f.pages[n-1]
		title := p.name
		if title == "" {
			title = "Page-" + strconv.Itoa(n)
		}
		view := c.doc.NewView(ids[i], bdf.ViewFixed, title)
		c.page = &pageInfo{name: title, number: n, count: len(f.pages)}
		pg, layers, err := c.renderPageSafe(p)
		if err != nil {
			return nil, fmt.Errorf("drawio: page %d: %w", n, err)
		}
		view.Pages = append(view.Pages, pg)
		views = append(views, pending{view, pg, layers})
	}
	c.finalize()
	c.reportMissing()
	c.reportAWSLegacy()
	for _, v := range views {
		for i, cv := range v.layers {
			v.page.Layers[i].Obj = cv.hash
		}
		if len(v.page.Layers) == 0 {
			o := bdf.NewObject()
			o.SetBBox(0, 0, v.page.W, v.page.H)
			h, _ := c.doc.AddObject(o)
			v.page.Layers = append(v.page.Layers, bdf.Layer{Role: bdf.RoleBody, Obj: h})
		}
		if !opts.NoTextIndex {
			if _, err := c.doc.BuildTextIndex(v.view); err != nil {
				c.warnf("text index: %v", err)
			}
		}
	}
	return &Result{Doc: c.doc, Warnings: c.warnings, Pages: len(sel), EmbeddedFonts: c.embeddedFonts}, nil
}

func (c *converter) reportMissing() {
	if len(c.missing) == 0 {
		return
	}
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

// renderPageSafe renders a page; a malformed diagram that trips the
// renderer becomes an empty page and a warning instead of failing the
// whole conversion.
func (c *converter) renderPageSafe(p *page) (pg *bdf.Page, layers []*canvas, err error) {
	defer func() {
		if r := recover(); r != nil {
			c.warnf("page %q: internal error: %v", c.page.name, r)
			pg, layers, err = &bdf.Page{W: 100, H: 100}, nil, nil
		}
	}()
	return c.renderPage(p)
}

// item is something drawn for a cell: its shape and label.
type item struct {
	st    *cellState
	shape *shape
	label *labelBox
}

func (c *converter) renderPage(p *page) (*bdf.Page, []*canvas, error) {
	m := parseModel(p.model)
	c.m = m
	c.pageShadow = m.attrs["shadow"] == "1"
	if m.attrs["math"] == "1" {
		c.warnOnce("math", "mathematical typesetting (math=1) is drawn as plain text")
	}
	c.awsLegacy += c.substituteAWSLegacy(m)
	v := newView(m, c.warnOnce, c.stencil)

	// shapes and label layouts, then the bounds of the drawing
	var items []*item
	var bounds rect
	for _, st := range v.order {
		if st.layer == nil || st.cell.geo == nil && !st.cell.edge {
			continue // layers themselves, cells without geometry
		}
		if st.cell.edge && st.absPoints == nil {
			continue
		}
		c.checkUnsupported(st)
		it := &item{st: st, shape: c.newShape(st)}
		it.label = c.layoutLabel(st, it.shape)
		items = append(items, it)
		bounds = bounds.union(it.shape.extent())
		if it.label != nil {
			bounds = bounds.union(it.label.extent())
		}
	}
	border := c.opts.Border
	if border == 0 {
		border = defaultBorder
	} else if border < 0 {
		border = 0
	}
	if bounds == (rect{}) {
		// an empty page: the page format of the diagram
		w, h := attrNum(p.model, "pageWidth", 827), attrNum(p.model, "pageHeight", 1169)
		bounds = rect{0, 0, math.Max(1, w), math.Max(1, h)}
		border = 0
	}
	bounds = bounds.grow(border)
	pg := &bdf.Page{W: f32(math.Ceil(bounds.w * ptPerUnit)), H: f32(math.Ceil(bounds.h * ptPerUnit))}
	origin := matrix{ptPerUnit, 0, 0, ptPerUnit, -bounds.x * ptPerUnit, -bounds.y * ptPerUnit}

	var layers []*canvas
	if bg, ok := parseColor(m.attrs["background"]); ok && bg.a > 0 {
		cv := c.newCanvas()
		cv.obj.SetBBox(0, 0, pg.W, pg.H)
		cv.obj.FillColor(bg.bdf()).FillRect(0, 0, pg.W, pg.H)
		cv.drawn = true
		pg.Layers = append(pg.Layers, bdf.Layer{Role: bdf.RoleBackground})
		layers = append(layers, cv)
	}
	for _, l := range m.root.children {
		if !l.visible || v.state(l) == nil {
			continue
		}
		cv := c.newCanvas()
		cv.obj.SetBBox(0, 0, pg.W, pg.H)
		cv.transform(origin)
		c2 := newC2D(cv)
		for _, it := range items {
			if it.st.layer != l {
				continue
			}
			c.drawItem(c2, it)
		}
		if cv.drawn {
			pg.Layers = append(pg.Layers, bdf.Layer{Role: bdf.RoleBody})
			layers = append(layers, cv)
		}
	}
	return pg, layers, nil
}

// checkUnsupported warns about styles drawn differently from draw.io.
func (c *converter) checkUnsupported(st *cellState) {
	s := st.style
	if s.is("sketch") || s.get("comic", "") == "1" {
		c.warnOnce("sketch", "hand-drawn styles (sketch) are drawn with straight lines")
	}
	// the patterns of mxSvgCanvas2D.getFillPattern; draw.io fills solid
	// for other values (such as those of mxgraph.basic.patternFillRect)
	switch v := s.get("fillStyle", ""); v {
	case "hatch", "dots", "cross-hatch", "dashed", "zigzag", "zigzag-line":
		c.warnOnce("fillStyle", "fill style %q is drawn as a solid fill", v)
	}
}

// drawItem draws a cell's shape and label, and its link.
func (c *converter) drawItem(c2 *c2d, it *item) {
	it.shape.paint(c2)
	if it.label != nil {
		c.drawLabel(c2, it.label)
	}
	if link := it.st.cell.attrs["link"]; link != "" {
		if u := c.linkURL(link); u != "" {
			r := it.st.bounds()
			c2.obj.Link(f32(r.x), f32(r.y), f32(r.w), f32(r.h), u)
		}
	}
}

// linkURL maps a draw.io link to a BDF LINK url: web and mail links stay,
// links to another page ("data:page/id,…") become "#view=" links, and
// everything else (actions, scripts, relative paths) is dropped.
func (c *converter) linkURL(link string) string {
	link = strings.TrimSpace(link)
	if id, ok := strings.CutPrefix(link, "data:page/id,"); ok {
		if v, ok := c.viewOf[id]; ok {
			return "#view=" + url.QueryEscape(v)
		}
		return ""
	}
	u, err := url.Parse(link)
	if err != nil {
		return ""
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https", "mailto":
		return u.String()
	}
	return ""
}

// extent is the rectangle the shape paints over (for the page size):
// its rotated bounds or its points, with the stroke and end markers.
func (s *shape) extent() rect {
	sw := s.strokewidth / 2
	if s.stroke == "" {
		sw = 0
	}
	if s.edge {
		r, ok := pointsExtent(s.points)
		if !ok {
			return rect{}
		}
		m := 0.0
		if s.startArrow != "" && s.startArrow != "none" || s.endArrow != "" && s.endArrow != "none" {
			m = math.Max(s.style.num("startSize", defaultMarkerSize), s.style.num("endSize", defaultMarkerSize)) + s.strokewidth
		}
		return s.augmented(r.grow(math.Max(sw, m/2)))
	}
	r := s.bounds
	// the paint box of north and south directions (mxShape.createBoundingBox);
	// stencils turn for directions themselves
	if s.isPaintBoundsInverted() || s.stencil != nil && (s.direction == "north" || s.direction == "south") {
		r = r.rotate90()
	}
	if rot := s.shapeRotation(); math.Mod(rot, 360) != 0 {
		// bounding box of the rotated rectangle
		rad := toRadians(rot)
		cos, sin := math.Abs(math.Cos(rad)), math.Abs(math.Sin(rad))
		w := r.w*cos + r.h*sin
		h := r.w*sin + r.h*cos
		r = rect{r.cx() - w/2, r.cy() - h/2, w, h}
	}
	if !s.hasPaint() {
		return rect{}
	}
	return s.augmented(r.grow(sw))
}

func itoa(i int) string { return strconv.Itoa(i) }

// fmtNum formats a number without trailing zeros.
func fmtNum(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }
