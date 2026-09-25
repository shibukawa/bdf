// Package pdf2bdf converts PDF files into BDF documents.
//
// Page content streams are interpreted and re-emitted as BDF instructions:
// paths, images, shadings, form XObjects (as shared objects) and text with
// embedded fonts rebuilt so browsers can load them. See docs/design.md.
package pdf2bdf

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/imgconv"
)

// Options controls the conversion.
type Options struct {
	// Pages selects 1-based page numbers; nil converts every page.
	Pages []int
	// Title overrides the document title.
	Title string
	// Kind is the view kind: "fixed" (default) or "flow".
	Kind string
	// NoTextIndex skips building the text index part.
	NoTextIndex bool
	// NoAnnotations skips annotation appearance streams and links.
	NoAnnotations bool
	// NoSubset keeps every glyph of embedded TrueType fonts instead of
	// dropping the outlines of unused ones (CFF fonts are never subset).
	NoSubset bool
	// Images controls whether raster images are re-encoded (see imgconv).
	// The zero value keeps images as they are.
	Images imgconv.Options
	// NoSharePrefix keeps each page body whole instead of moving the
	// instruction prefix that two or more pages have in common (a master
	// slide, a letterhead) into one shared object.
	NoSharePrefix bool
	// Warn receives non-fatal problems; when nil they are collected in Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc      *bdf.Document
	Warnings []string
	Pages    int
	// SharedPrefixes counts the prefix objects shared between pages, and
	// SharedBytes the op-stream bytes the pages no longer carry.
	SharedPrefixes int
	SharedBytes    int
}

type imageEntry struct {
	hash bdf.Hash
	ok   bool
	size [2]int
}

type converter struct {
	pdf      *pdf
	doc      *bdf.Document
	opts     *Options
	warnings []string
	warned   map[string]bool

	fonts          map[string]*pdfFont
	forms          map[string]*pending
	images         map[string]*imageEntry
	shadings       map[string]*shading
	pendings       []*pending
	pageBox        rect // current page box in object space
	deflt          *pdfFont
	pageBodies     []pageRef
	sharedPrefixes int
	sharedBytes    int
}

// ConvertFile converts a PDF file.
func ConvertFile(path string, opts *Options) (*Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Convert(f, opts)
}

// Convert converts a PDF read from rs.
func Convert(rs io.ReadSeeker, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed
	conf.DecodeAllStreams = false
	ctx, err := api.ReadContext(rs, conf)
	if err != nil {
		return nil, fmt.Errorf("pdf2bdf: %w", err)
	}
	if err := ctx.EnsurePageCount(); err != nil {
		return nil, fmt.Errorf("pdf2bdf: %w", err)
	}
	c := &converter{pdf: &pdf{ctx: ctx}, doc: bdf.NewDocument(), opts: opts, warned: map[string]bool{},
		fonts: map[string]*pdfFont{}, forms: map[string]*pending{}, images: map[string]*imageEntry{}, shadings: map[string]*shading{}}
	c.doc.Meta.Source = "pdf"
	c.doc.Meta.Title = opts.Title
	if c.doc.Meta.Title == "" {
		if info := c.pdf.dict(ctx.Info); info != nil {
			c.doc.Meta.Title = c.pdf.text(info["Title"])
		}
	}
	kind := opts.Kind
	if kind == "" {
		kind = bdf.ViewFixed
	}
	view := c.doc.NewView("pages", kind, c.doc.Meta.Title)
	if kind == bdf.ViewFlow {
		view.Continuous = &bdf.Continuous{Gap: 16}
	}
	pages := opts.Pages
	if pages == nil {
		for i := 1; i <= ctx.PageCount; i++ {
			pages = append(pages, i)
		}
	}
	for _, n := range pages {
		if n < 1 || n > ctx.PageCount {
			return nil, fmt.Errorf("pdf2bdf: page %d out of range (1-%d)", n, ctx.PageCount)
		}
		if err := c.convertPage(view, n); err != nil {
			return nil, fmt.Errorf("pdf2bdf: page %d: %w", n, err)
		}
	}
	if !opts.NoSharePrefix {
		c.sharePagePrefixes()
	}
	c.finalize()
	if !opts.NoTextIndex {
		if _, err := c.doc.BuildTextIndex(view); err != nil {
			c.warnf("text index: %v", err)
		}
	}
	return &Result{Doc: c.doc, Warnings: c.warnings, Pages: len(pages), SharedPrefixes: c.sharedPrefixes, SharedBytes: c.sharedBytes}, nil
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

func (c *converter) defaultFont() *pdfFont {
	if c.deflt == nil {
		c.deflt = &pdfFont{subtype: "Type1", baseFont: "Helvetica", used: map[uint32]*usedGlyph{}, uniOwner: map[rune]int{}, puaNext: 0xE000,
			warned: map[string]bool{}, family: "sans-serif", weight: 400, fontMatrix: matrix{0.001, 0, 0, 0.001, 0, 0}, defaultWidth: 1000}
		c.deflt.baseNames = standardEncoding
		c.deflt.hasBaseEnc = true
	}
	return c.deflt
}

// convertPage converts one page into a body object and appends it to the view.
func (c *converter) convertPage(view *bdf.View, pageNr int) error {
	p := c.pdf
	pageDict, _, attrs, err := c.pdf.ctx.PageDict(pageNr, false)
	if err != nil {
		return err
	}
	if pageDict == nil {
		return fmt.Errorf("page dict missing")
	}
	// Boxes and rotation.
	box := rect{0, 0, 612, 792}
	if attrs.MediaBox != nil {
		box = rect{attrs.MediaBox.LL.X, attrs.MediaBox.LL.Y, attrs.MediaBox.UR.X, attrs.MediaBox.UR.Y}
	}
	if attrs.CropBox != nil {
		cb := rect{attrs.CropBox.LL.X, attrs.CropBox.LL.Y, attrs.CropBox.UR.X, attrs.CropBox.UR.Y}
		if !cb.intersect(box).empty() {
			box = cb.intersect(box)
		}
	}
	if box.empty() || box.x1-box.x0 > 20000 || box.y1-box.y0 > 20000 {
		box = rect{0, 0, 612, 792}
	}
	w, h := box.x1-box.x0, box.y1-box.y0
	rotate := ((attrs.Rotate % 360) + 360) % 360
	// PDF user space → page space (y down), then rotation.
	base := matrix{1, 0, 0, -1, -box.x0, box.y1}
	pw, ph := w, h
	switch rotate {
	case 90:
		base = base.mul(matrix{0, 1, -1, 0, h, 0})
		pw, ph = h, w
	case 180:
		base = base.mul(matrix{-1, 0, 0, -1, w, h})
	case 270:
		base = base.mul(matrix{0, -1, 1, 0, 0, w})
		pw, ph = h, w
	}
	c.pageBox = rect{0, 0, pw, ph}

	body := newPending(bdf.Rect{X: 0, Y: 0, W: float32(pw), H: float32(ph)})
	c.pendings = append(c.pendings, body)
	res := attrs.Resources
	if res == nil {
		res = p.dict(pageDict["Resources"])
	}
	content, err := c.pdf.ctx.PageContent(pageDict, pageNr)
	if err != nil && err != model.ErrNoContent {
		c.warnf("page %d content: %v", pageNr, err)
	}
	in := c.newInterp(body, res, base, 0)
	in.transform(base)
	if len(content) > 0 {
		in.run(content)
	}
	if !c.opts.NoAnnotations {
		c.annotations(in, pageDict, base)
	}
	page := view.AddPage(float32(pw), float32(ph), bdf.Layer{Role: bdf.RoleBody, Obj: bdf.Hash{}})
	page.Layers[0].Obj = bdf.Hash{} // patched in finalize
	c.pageBodies = append(c.pageBodies, pageRef{page: page, body: body})
	if view.Kind == bdf.ViewFlow {
		page.Body = &bdf.RectDef{X: 0, Y: 0, W: float32(pw), H: float32(ph)}
	}
	return nil
}

type pageRef struct {
	page *bdf.Page
	body *pending
}

// annotations draws annotation appearance streams and emits link areas.
func (c *converter) annotations(in *interp, pageDict types.Dict, base matrix) {
	p := c.pdf
	for _, a := range p.array(pageDict["Annots"]) {
		ad := p.dict(a)
		if ad == nil {
			continue
		}
		sub := p.name(ad["Subtype"])
		flags := p.intOr(ad["F"], 0)
		if sub == "Popup" || flags&2 != 0 || flags&32 != 0 {
			continue
		}
		r := p.rectOr(ad["Rect"], rect{})
		if r.empty() {
			continue
		}
		if sub == "Link" {
			target := ""
			if act := p.dict(ad["A"]); act != nil {
				if p.name(act["S"]) == "URI" {
					target = string(p.str(act["URI"]))
				}
			}
			if target == "" {
				if dest := p.array(ad["Dest"]); len(dest) > 0 {
					if ref, ok := dest[0].(types.IndirectRef); ok {
						if n, err := c.pdf.ctx.PageNumber(ref.ObjectNumber.Value()); err == nil {
							target = fmt.Sprintf("#page=%d", n)
						}
					}
				}
			}
			if target != "" {
				lr := base.transformRect(r)
				in.obj.Link(float32(lr.x0), float32(lr.y0), float32(lr.x1-lr.x0), float32(lr.y1-lr.y0), target)
			}
			continue
		}
		ap := p.dict(ad["AP"])
		if ap == nil {
			continue
		}
		n := ap["N"]
		sd := p.stream(n)
		if sd == nil {
			if nd := p.dict(n); nd != nil {
				as := p.name(ad["AS"])
				if as == "" && len(nd) == 1 {
					for k := range nd {
						as = k
					}
				}
				sd = p.stream(nd[as])
			}
		}
		if sd == nil {
			continue
		}
		// Algorithm 8.1: map the transformed BBox to Rect.
		bbox := p.rectOr(sd.Dict["BBox"], rect{})
		mat := p.matrixOr(sd.Dict["Matrix"], identity)
		tb := mat.transformRect(bbox)
		extra := identity
		if !tb.empty() {
			sx, sy := (r.x1-r.x0)/(tb.x1-tb.x0), (r.y1-r.y0)/(tb.y1-tb.y0)
			extra = matrix{sx, 0, 0, sy, r.x0 - tb.x0*sx, r.y0 - tb.y0*sy}
		}
		in.gs = initialGState(in.gs.ctm)
		in.drawForm(objKey(n)+"/"+p.name(ad["AS"]), sd, extra)
	}
}

// formObject converts a form XObject (or pattern cell) into a shared pending object.
func (c *converter) formObject(key string, sd *types.StreamDict, parentRes types.Dict, depth int) *pending {
	p := c.pdf
	if key != "" {
		if f, ok := c.forms[key]; ok {
			if f == nil || f.building {
				return nil // in progress (cycle) or failed
			}
			return f
		}
	}
	bbox := p.rectOr(sd.Dict["BBox"], rect{0, 0, 1, 1})
	if bbox.empty() {
		bbox = rect{0, 0, 1, 1}
	}
	f := newPending(bdf.Rect{X: float32(bbox.x0), Y: float32(bbox.y0), W: float32(bbox.x1 - bbox.x0), H: float32(bbox.y1 - bbox.y0)})
	f.building = true
	if key != "" {
		c.forms[key] = f
	}
	res := p.dict(sd.Dict["Resources"])
	if res == nil {
		res = parentRes
	}
	data, _, err := p.decodeStream(sd)
	if err != nil {
		c.warnf("form XObject: %v", err)
		if key != "" {
			c.forms[key] = nil
		}
		return nil
	}
	in := c.newInterp(f, res, identity, depth)
	in.run(data)
	f.building = false
	c.pendings = append(c.pendings, f)
	return f
}

// type3Glyph converts a Type3 glyph procedure into a pending object.
func (c *converter) type3Glyph(f *pdfFont, name string, depth int) *pending {
	if name == "" || f.charProcs == nil {
		return nil
	}
	if g, ok := f.t3Glyphs[name]; ok {
		return g
	}
	if depth > 8 {
		return nil
	}
	p := c.pdf
	sd := p.stream(f.charProcs[name])
	if sd == nil {
		f.t3Glyphs[name] = nil
		return nil
	}
	data, _, err := p.decodeStream(sd)
	if err != nil {
		f.t3Glyphs[name] = nil
		return nil
	}
	bb := p.rectOr(f.dict["FontBBox"], rect{})
	if bb.empty() {
		inv, _ := f.fontMatrix.inverse()
		bb = inv.transformRect(rect{-0.2, -0.3, 1.2, 1.0})
	}
	g := newPending(bdf.Rect{X: float32(bb.x0), Y: float32(bb.y0), W: float32(bb.x1 - bb.x0), H: float32(bb.y1 - bb.y0)})
	f.t3Glyphs[name] = g
	in := c.newInterp(g, f.t3Resources, identity, depth)
	in.type3Glyph = true
	in.run(data)
	c.pendings = append(c.pendings, g)
	return g
}

func (c *converter) shadingFor(key string, o types.Object, res types.Dict) *shading {
	if key != "" {
		if s, ok := c.shadings[key]; ok {
			return s
		}
	}
	s := c.loadShading(o, res)
	if key != "" {
		c.shadings[key] = s
	}
	return s
}

// finalize builds font programs, encodes pending objects (children first) and
// patches page layers.
func (c *converter) finalize() {
	var encode func(p *pending) bdf.Hash
	encode = func(p *pending) bdf.Hash {
		if p.encoded {
			return p.hash
		}
		p.encoded = true
		fontRefs := make([]int, 0, len(p.fonts))
		for r := range p.fonts {
			fontRefs = append(fontRefs, int(r))
		}
		sort.Ints(fontRefs)
		for _, r := range fontRefs {
			p.obj.UpdateFont(bdf.FontRef(r), c.finalizeFont(p.fonts[bdf.FontRef(r)]))
		}
		// Encode children in a stable order.
		refs := make([]int, 0, len(p.children))
		for r := range p.children {
			refs = append(refs, int(r))
		}
		sort.Ints(refs)
		for _, r := range refs {
			child := p.children[bdf.ObjRef(r)]
			p.obj.UpdateObject(bdf.ObjRef(r), encode(child))
		}
		p.hash, _ = c.doc.AddObject(p.obj)
		return p.hash
	}
	for _, p := range c.pendings {
		encode(p)
	}
	for _, pr := range c.pageBodies {
		pr.page.Layers[0].Obj = pr.body.hash
	}
}

// PageRange parses "1-3,5,8-" style selections.
func PageRange(spec string, count int) ([]int, error) {
	var out []int
	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		lo, hi := 1, count
		if i := strings.IndexByte(part, '-'); i >= 0 {
			if i > 0 {
				if _, err := fmt.Sscanf(part[:i], "%d", &lo); err != nil {
					return nil, fmt.Errorf("bad page range %q", part)
				}
			}
			if i+1 < len(part) {
				if _, err := fmt.Sscanf(part[i+1:], "%d", &hi); err != nil {
					return nil, fmt.Errorf("bad page range %q", part)
				}
			}
		} else {
			if _, err := fmt.Sscanf(part, "%d", &lo); err != nil {
				return nil, fmt.Errorf("bad page %q", part)
			}
			hi = lo
		}
		lo, hi = max(1, lo), min(count, hi)
		for n := lo; n <= hi; n++ {
			out = append(out, n)
		}
	}
	return out, nil
}
