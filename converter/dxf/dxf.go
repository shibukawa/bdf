// Package dxf converts DXF drawings (AutoCAD's Drawing Exchange Format,
// text and binary, R12 to 2018) into BDF documents.
//
// Model space becomes a view of one page that fits the drawing, on the dark
// background CAD programs show it on; each paper space layout with
// something on it becomes a view of one page of its paper, where the
// viewports show model space at their scale. Lines are drawn with their
// lineweights as a plotter would, text with fonts that stand in for
// AutoCAD's (SHX fonts become sans-serif ones, big fonts East Asian ones),
// and the fonts in use are embedded as subsets. See docs/design.md §3.12.
package dxf

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
)

// Options controls the conversion.
type Options struct {
	// Pages selects 1-based layouts in tab order, model space first; nil
	// converts every one that has something on it.
	Pages []int
	// Title overrides the document title.
	Title string
	// Views selects what to convert: "all" (the default), "model" or
	// "layouts".
	Views string
	// Light draws model space on white paper instead of the dark
	// background of CAD programs.
	Light bool
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
	// Warn receives non-fatal problems; when nil they are collected in
	// Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc           *bdf.Document
	Warnings      []string
	Views         int
	EmbeddedFonts int
}

// darkBackground is the background of model space (AutoCAD's default).
var darkBackground = bdf.RGB(33, 40, 48)

// modelSize is the longer side of the page of model space (A3, in pt).
const modelSize = 420 * 72 / 25.4

type converter struct {
	opts      *Options
	d         *document
	doc       *bdf.Document
	fonts     *cad.Fonts
	cvs       *canvas.Builder
	plotter   *cad.Plotter
	warnings  []string
	warned    map[string]bool
	ltscale   float64
	psltscale bool
	model     []*entity // the entities of model space
	count     int       // entities drawn for the view being made
}

// ConvertFile converts a DXF file.
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

// maxSize bounds the files read.
const maxSize = 1 << 30

// Convert converts a drawing read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if size > maxSize {
		return nil, fmt.Errorf("dxf: the file is larger than %d bytes", maxSize)
	}
	data := make([]byte, size)
	if _, err := r.ReadAt(data, 0); err != nil && err != io.EOF {
		return nil, fmt.Errorf("dxf: %w", err)
	}
	c := &converter{opts: opts, doc: bdf.NewDocument(), warned: map[string]bool{}}
	tags, err := readTags(data, func(msg string) { c.warnf("%s", msg) })
	if err != nil {
		if len(tags) == 0 {
			return nil, fmt.Errorf("dxf: %w", err)
		}
		c.warnf("%v; the rest of the file is ignored", err)
	}
	c.d = parse(tags)
	c.doc.Meta.Source = "dxf"
	c.ltscale = c.d.hnum("$LTSCALE", 1)
	if c.ltscale <= 0 {
		c.ltscale = 1
	}
	c.psltscale = c.d.hnum("$PSLTSCALE", 1) != 0
	switch c.d.encoding {
	case "ansi_932", "shift_jis":
		c.doc.Meta.DC.Language = bdf.DCValues{"ja"}
	case "ansi_936":
		c.doc.Meta.DC.Language = bdf.DCValues{"zh-Hans"}
	case "ansi_949":
		c.doc.Meta.DC.Language = bdf.DCValues{"ko"}
	case "ansi_950":
		c.doc.Meta.DC.Language = bdf.DCValues{"zh-Hant"}
	}
	if opts.Title != "" {
		c.doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	db := fontdb.New(opts.FontFS, opts.FontDirs, !opts.NoSystemFonts)
	set := fontset.New(db, func(msg string) { c.warnf("%s", msg) })
	c.fonts = &cad.Fonts{Set: set}
	c.plotter = &cad.Plotter{Fonts: c.fonts, Thin: 0.1 * 72 / 25.4}
	c.cvs = canvas.NewBuilder(c.doc, set)
	if len(db.Faces) == 0 && !opts.SystemFonts {
		c.warnf("no fonts found; text is laid out with estimated metrics and not embedded")
	}

	type pageRef struct {
		view *bdf.View
		page *bdf.Page
		cv   *canvas.Canvas
	}
	var pages []pageRef
	sheets := c.sheets()
	sel := opts.Pages
	if sel == nil {
		for i := range sheets {
			sel = append(sel, i+1)
		}
	}
	for _, n := range sel {
		if n < 1 || n > len(sheets) {
			return nil, fmt.Errorf("dxf: layout %d out of range (1-%d)", n, len(sheets))
		}
		s := sheets[n-1]
		view := c.doc.NewView(s.id, bdf.ViewFixed, s.title)
		page, cv := c.renderSafe(s, view)
		pages = append(pages, pageRef{view, page, cv})
	}
	res := &Result{Doc: c.doc, Views: len(pages)}
	if !opts.SystemFonts {
		res.EmbeddedFonts = set.Embed(c.doc, fontset.EmbedOptions{NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType})
	}
	c.cvs.Encode()
	set.ReportMissing()
	for _, pr := range pages {
		pr.page.Layers = []bdf.Layer{{Role: bdf.RoleBody, Obj: pr.cv.Hash()}}
		if !opts.NoTextIndex {
			if _, err := c.doc.BuildTextIndex(pr.view); err != nil {
				c.warnf("text index: %v", err)
			}
		}
	}
	res.Warnings = c.warnings
	return res, nil
}

// sheet is a view to make: model space or a layout.
type sheet struct {
	id, title string
	model     bool
	layout    *layout
	entities  []*entity
}

// sheets lists model space and the layouts that have something on them,
// in tab order.
func (c *converter) sheets() []sheet {
	var model, paper []*entity
	for _, e := range c.d.entities {
		if e.int(67, 0) == 1 {
			paper = append(paper, e)
		} else {
			model = append(model, e)
		}
	}
	if b := c.d.blocks["*MODEL_SPACE"]; b != nil {
		model = append(model, b.entities...)
	}
	c.model = model
	views := strings.ToLower(c.opts.Views)
	var out []sheet
	if views != "layouts" {
		out = append(out, sheet{id: "model", title: "Model", model: true, entities: model})
	}
	if views == "model" {
		return out
	}
	n := 0
	add := func(title string, l *layout, ents []*entity) {
		if !hasContent(ents) {
			return
		}
		n++
		out = append(out, sheet{id: fmt.Sprintf("layout%d", n), title: title, layout: l, entities: ents})
	}
	if len(c.d.layouts) == 0 {
		add("Paper space", nil, paper)
		return out
	}
	for _, l := range c.d.layouts {
		if strings.EqualFold(l.name, "Model") || l.blockKey == "*MODEL_SPACE" {
			continue
		}
		var ents []*entity
		if l.blockKey == "*PAPER_SPACE" {
			ents = paper
		}
		if b := c.d.blocks[l.blockKey]; b != nil {
			ents = append(ents, b.entities...)
		}
		add(l.name, l, ents)
	}
	return out
}

// hasContent reports whether a layout has more than its own viewport.
func hasContent(ents []*entity) bool {
	for _, e := range ents {
		if e.typ == "VIEWPORT" && e.int(69, 0) == 1 {
			continue
		}
		return true
	}
	return false
}

func (c *converter) renderSafe(s sheet, view *bdf.View) (page *bdf.Page, cv *canvas.Canvas) {
	cv = c.cvs.New()
	c.count = 0
	defer func() {
		if r := recover(); r != nil {
			c.warnf("%s: internal error: %v", s.title, r)
			if page == nil {
				page = view.AddPage(595.28, 841.89)
			}
			cv = c.cvs.New()
			cv.Obj.SetBBox(0, 0, page.W, page.H)
		}
	}()
	if s.model {
		return c.modelPage(s, view, cv), cv
	}
	return c.layoutPage(s, view, cv), cv
}

func (c *converter) topCtx(out *cad.Drawing, dark, plot bool, deferred *[]func(cad.Rect)) *ctx {
	return &ctx{m: canvas.Identity, out: out, dark: dark, plot: plot, ltFactor: 1, deferred: deferred,
		byBlock: props{color: c.aci(7, dark), ltype: "CONTINUOUS", lw: -3}}
}

// draw draws the entities of a space; a malformed entity is skipped with
// a warning.
func (c *converter) draw(ents []*entity, x *ctx) {
	for _, e := range ents {
		func() {
			defer func() {
				if r := recover(); r != nil {
					c.warnOnce("panic:"+e.typ, "%s: internal error: %v", e.typ, r)
				}
			}()
			c.entity(e, x)
		}()
	}
}

func (c *converter) modelPage(s sheet, view *bdf.View, cv *canvas.Canvas) *bdf.Page {
	dark := !c.opts.Light
	d := &cad.Drawing{}
	var deferred []func(cad.Rect)
	c.draw(s.entities, c.topCtx(d, dark, false, &deferred))
	b := d.Bounds()
	if !b.Valid() {
		if lo, ok := c.d.hpt("$EXTMIN"); ok {
			if hi, ok := c.d.hpt("$EXTMAX"); ok && hi.X > lo.X && hi.Y > lo.Y {
				b = b.Add(lo).Add(hi)
			}
		}
	}
	if !b.Valid() {
		b = b.Add(cad.Point{}).Add(cad.Point{X: 297, Y: 210})
	}
	for _, f := range deferred {
		f(b)
	}
	w, h := b.W(), b.H()
	if w <= 0 && h <= 0 {
		w, h = 1, 1
	}
	margin := modelSize * 0.03
	scale := (modelSize - 2*margin) / math.Max(w, h)
	pw, ph := w*scale+2*margin, h*scale+2*margin
	pw, ph = math.Max(pw, 4*margin), math.Max(ph, 4*margin)
	page := view.AddPage(float32(pw), float32(ph))
	// centre a drawing narrower than the minimum page
	ox := (pw - w*scale) / 2
	oy := (ph - h*scale) / 2
	m := canvas.Matrix{scale, 0, 0, -scale, ox - b.Min.X*scale, oy + b.Max.Y*scale}
	cv.Obj.SetBBox(0, 0, page.W, page.H)
	bg := bdf.RGB(255, 255, 255)
	if dark {
		bg = darkBackground
	}
	cv.Obj.FillColor(bg)
	cv.Obj.FillRect(0, 0, page.W, page.H)
	cv.Drawn = true
	c.plotter.Plot(cv, d, m, pageRect(pw, ph))
	return page
}

func pageRect(w, h float64) cad.Rect {
	return cad.Rect{}.Add(cad.Point{}).Add(cad.Point{X: w, Y: h})
}

// paper returns the size of a layout's paper in mm (as it is shown, turned
// by the plot rotation), and the position on it of the origin of paper
// space.
func (c *converter) paper(s sheet) (w, h float64, origin cad.Point, unit float64) {
	unit = 1
	if s.layout != nil {
		ps := s.layout.plot
		w, h = ps.num(44, 0), ps.num(45, 0)
		left, bottom, right, top := ps.num(40, 0), ps.num(41, 0), ps.num(42, 0), ps.num(43, 0)
		switch ps.int(73, 0) {
		case 1:
			w, h = h, w
			left, bottom = top, left
		case 2:
			left, bottom = right, top
		case 3:
			w, h = h, w
			left, bottom = bottom, right
		}
		if ps.int(72, 1) == 0 {
			unit = 25.4
		}
		if num, den := ps.num(142, 1), ps.num(143, 1); num > 0 && den > 0 {
			unit *= num / den
		}
		origin = cad.Point{X: left + ps.num(46, 0), Y: bottom + ps.num(47, 0)}
		if w > 0 && h > 0 {
			return
		}
		// no paper: the limits of the layout
		lo, hi := s.layout.lay.pt(10), s.layout.lay.pt(11)
		if hi.X > lo.X && hi.Y > lo.Y {
			return (hi.X - lo.X) * unit, (hi.Y - lo.Y) * unit, cad.Point{X: -lo.X * unit, Y: -lo.Y * unit}, unit
		}
	}
	// R12: the paper space limits
	if lo, ok := c.d.hpt("$PLIMMIN"); ok {
		if hi, ok := c.d.hpt("$PLIMMAX"); ok && hi.X > lo.X && hi.Y > lo.Y {
			return hi.X - lo.X, hi.Y - lo.Y, cad.Point{X: -lo.X, Y: -lo.Y}, 1
		}
	}
	return 0, 0, cad.Point{}, 1
}

func (c *converter) layoutPage(s sheet, view *bdf.View, cv *canvas.Canvas) *bdf.Page {
	d := &cad.Drawing{}
	var deferred []func(cad.Rect)
	x := c.topCtx(d, false, true, &deferred)
	var vps []*entity
	for _, e := range s.entities {
		if e.typ == "VIEWPORT" {
			vps = append(vps, e)
		}
	}
	// viewports first: paper space is drawn over what they show
	for _, vp := range vps {
		c.viewport(vp, s, d)
	}
	c.draw(s.entities, x)
	w, h, origin, unit := c.paper(s)
	const mm = 72 / 25.4
	if w <= 0 || h <= 0 {
		// no paper: fit the content
		b := d.Bounds()
		if !b.Valid() {
			b = pageRect(420, 297)
		}
		w, h = b.W()*unit+20, b.H()*unit+20
		origin = cad.Point{X: 10 - b.Min.X*unit, Y: 10 - b.Min.Y*unit}
	}
	for _, f := range deferred {
		f(cad.Rect{}.Add(cad.Point{X: -origin.X / unit, Y: -origin.Y / unit}).Add(cad.Point{X: (w - origin.X) / unit, Y: (h - origin.Y) / unit}))
	}
	pw, ph := w*mm, h*mm
	page := view.AddPage(float32(pw), float32(ph))
	k := unit * mm
	m := canvas.Matrix{k, 0, 0, -k, origin.X * mm, ph - origin.Y*mm}
	cv.Obj.SetBBox(0, 0, page.W, page.H)
	c.plotter.Plot(cv, d, m, pageRect(pw, ph))
	return page
}

// viewport draws model space as a VIEWPORT of a layout shows it.
func (c *converter) viewport(vp *entity, s sheet, out *cad.Drawing) {
	id := vp.int(69, 0)
	if id == 1 {
		return
	}
	if vp.has(68) {
		if vp.int(68, 0) <= 0 {
			return
		}
	} else if id <= 1 {
		return
	}
	view := viewportView(vp)
	dir := view.dir
	if math.Abs(dir[0]) > 1e-9 || math.Abs(dir[1]) > 1e-9 || dir[2] <= 0 {
		c.warnOnce("vp3d", "viewports that show model space other than from the top are not drawn")
		return
	}
	ctr := vp.pt(10)
	vw, vh := vp.num(40, 0), vp.num(41, 0)
	viewH := view.height
	if vw <= 0 || vh <= 0 || viewH <= 0 {
		return
	}
	scale := vh / viewH
	viewCtr := view.center
	target := view.target
	twist := view.twist
	m := canvas.Translate(ctr.X-viewCtr.X*scale, ctr.Y-viewCtr.Y*scale).Mul(canvas.Rotate(twist)).
		Mul(canvas.Scale(scale, scale)).Mul(canvas.Translate(-target[0], -target[1]))
	clip := (&cad.Path{}).Polyline([]cad.Point{
		{X: ctr.X - vw/2, Y: ctr.Y - vh/2}, {X: ctr.X + vw/2, Y: ctr.Y - vh/2},
		{X: ctr.X + vw/2, Y: ctr.Y + vh/2}, {X: ctr.X - vw/2, Y: ctr.Y + vh/2}}, true)
	if vp.int(90, 0)&0x10000 != 0 {
		if b := c.d.byHandle[strings.ToUpper(vp.str(340))]; b != nil {
			if p := outline(b); p != nil {
				clip = p
			}
		}
	}
	frozen := map[string]bool{}
	for _, t := range vp.tags {
		if t.code == 331 {
			if l := c.d.layerHandles[strings.ToUpper(t.s)]; l != nil {
				frozen[key(l.name)] = true
			}
		}
	}
	for _, name := range view.frozen {
		frozen[key(name)] = true
	}
	var deferred []func(cad.Rect)
	x := c.topCtx(out, false, true, &deferred)
	x.m = m
	x.frozen = frozen
	if c.psltscale {
		x.ltFactor = 1 / scale
	}
	out.Begin(clip)
	c.draw(c.model, x)
	for _, f := range deferred {
		f(clip.Bounds())
	}
	out.End()
}

// vpView is what a viewport shows of model space.
type vpView struct {
	target, dir   [3]float64
	center        cad.Point
	height, twist float64
	frozen        []string // layer names (R12)
}

// viewportView reads the view of a viewport: from its codes, or in R12
// files from the MVIEW list of its extended data (target, direction, twist,
// height, centre, ..., frozen layers).
func viewportView(vp *entity) vpView {
	v := vpView{target: vp.vec3(17, [3]float64{}), dir: vp.vec3(16, [3]float64{0, 0, 1}), center: vp.pt(12),
		height: vp.num(45, 0), twist: vp.num(51, 0)}
	if vp.has(45) {
		return v
	}
	xd := vp.xdata("ACAD")
	i := 0
	for i < len(xd) && !(xd[i].code == 1000 && strings.EqualFold(xd[i].s, "MVIEW")) {
		i++
	}
	if i == len(xd) {
		return v
	}
	// the values of the list, a point (1010, 1020, 1030) counting as one
	var vals []tag
	var pts [][3]float64
	depth := 0
	for i++; i < len(xd); i++ {
		t := xd[i]
		switch {
		case t.code == 1002 && strings.TrimSpace(t.s) == "{":
			depth++
			continue
		case t.code == 1002:
			depth--
			continue
		case t.code == 1003 && depth == 2:
			v.frozen = append(v.frozen, t.s)
			continue
		case t.code == 1010:
			pts = append(pts, [3]float64{t.f})
			vals = append(vals, tag{code: 1010, f: float64(len(pts) - 1)})
			continue
		case t.code == 1020 && len(pts) > 0:
			pts[len(pts)-1][1] = t.f
			continue
		case t.code == 1030 && len(pts) > 0:
			pts[len(pts)-1][2] = t.f
			continue
		}
		if depth == 1 {
			vals = append(vals, t)
		}
	}
	// version, target, direction, twist, height, centre x, centre y
	if len(vals) < 7 || vals[1].code != 1010 || vals[2].code != 1010 {
		return v
	}
	v.target, v.dir = pts[int(vals[1].f)], pts[int(vals[2].f)]
	v.twist, v.height = vals[3].f, vals[4].f
	v.center = cad.Point{X: vals[5].f, Y: vals[6].f}
	return v
}

// outline returns the outline of a paper space entity that clips a
// viewport.
func outline(e *entity) *cad.Path {
	p := &cad.Path{}
	switch e.typ {
	case "LWPOLYLINE":
		vs, _, _ := lwVertices(e)
		for i, v := range vs {
			if i == 0 {
				p.MoveTo(v.p.X, v.p.Y)
			}
			p.BulgeTo(vs[(i+1)%len(vs)].p, v.bulge)
		}
		p.Close()
		return p.Transform(ocs(e, e.num(38, 0)))
	case "POLYLINE":
		var vs []vertex
		for _, k := range e.kids {
			vs = append(vs, vertex{p: k.pt(10), bulge: k.num(42, 0)})
		}
		for i, v := range vs {
			if i == 0 {
				p.MoveTo(v.p.X, v.p.Y)
			}
			p.BulgeTo(vs[(i+1)%len(vs)].p, v.bulge)
		}
		return p.Close()
	case "CIRCLE":
		ctr := e.vec3(10, [3]float64{})
		return p.Circle(pt2(ctr), e.num(40, 0)).Transform(ocs(e, ctr[2]))
	case "ELLIPSE":
		ctr := e.vec3(10, [3]float64{})
		maj := e.vec3(11, [3]float64{1, 0, 0})
		return p.Ellipse(pt2(ctr), pt2(maj), e.num(40, 1))
	}
	return nil
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
