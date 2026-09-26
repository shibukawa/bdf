// Package pdf converts PDF files into BDF documents.
//
// Page content streams are interpreted and re-emitted as BDF instructions:
// paths, images, shadings, form XObjects (as shared objects) and text with
// embedded fonts rebuilt so browsers can load them. See docs/design.md.
package pdf

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
	"github.com/shibukawa/bdf/imgconv"
)

// Options controls the conversion.
type Options struct {
	// Pages selects 1-based pages (see converter.Pages); nil converts every
	// page.
	Pages conv.Pages
	// Title overrides the document title.
	Title string
	// Kind is the view kind: "fixed" (default) or "flow".
	Kind string
	// NoTextIndex skips building the text index part.
	NoTextIndex bool
	// NoAnnotations skips annotation appearance streams and links.
	NoAnnotations bool
	// NoSubset keeps every glyph of embedded fonts instead of dropping the
	// ones the document does not use.
	NoSubset bool
	// NoWOFF2 stores rebuilt fonts as TrueType/OpenType instead of WOFF2
	// (builds tagged bdf_noconv never produce WOFF2).
	NoWOFF2 bool
	// IgnoreFSType embeds fonts whose OS/2 fsType forbids it (Restricted
	// License, bitmap only) or forbids subsetting. Without it such fonts fall
	// back to a system font or are embedded whole. Set it only when you hold
	// the rights to embed the fonts of the document.
	IgnoreFSType bool
	// Images controls whether raster images are re-encoded (see imgconv).
	// The zero value keeps images as they are.
	Images imgconv.Options
	// NoSharePrefix keeps each page body whole instead of moving the
	// instruction prefix that two or more pages have in common (a master
	// slide, a letterhead) into one shared object.
	NoSharePrefix bool
	// Password opens a PDF encrypted with a user (open) password; the owner
	// password opens it too.
	Password string
	// Warn receives non-fatal problems; when nil they are collected in Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc      *bdf.Document
	Warnings []string
	Pages    int
	// Protected reports that the PDF opened only with Options.Password.
	Protected bool
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
	pageBodies     []*pageRef // the pages of the view, in order
	sharedPrefixes int
	sharedBytes    int
	// pageCodes collects the codes each font draws while a Stream converts
	// a page (nil otherwise).
	pageCodes map[*pdfFont]map[uint32]bool

	tree       *structTree // nil for untagged documents
	docLang    string      // catalog /Lang
	outPages   map[int]int // page object number → 1-based page index in the view
	outPageNrs map[int]int // PDF page number → 1-based page index in the view
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

// readContext reads the PDF, first without a password. A PDF encrypted with
// a user password is read again with password, and reported protected.
// Wrong or missing passwords are converter.ErrWrongPassword and
// converter.ErrPasswordRequired.
func readContext(rs io.ReadSeeker, password string) (ctx *model.Context, protected bool, err error) {
	read := func(pw string) (*model.Context, error) {
		if _, err := rs.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		conf := model.NewDefaultConfiguration()
		conf.ValidationMode = model.ValidationRelaxed
		conf.DecodeAllStreams = false
		conf.UserPW, conf.OwnerPW = pw, pw
		return api.ReadContext(rs, conf)
	}
	ctx, err = read("")
	if errors.Is(err, pdfcpu.ErrWrongPassword) {
		protected = true
		if password == "" {
			return nil, true, fmt.Errorf("pdf: %w", conv.ErrPasswordRequired)
		}
		ctx, err = read(password)
		if errors.Is(err, pdfcpu.ErrWrongPassword) {
			return nil, true, fmt.Errorf("pdf: %w", conv.ErrWrongPassword)
		}
	}
	if err != nil {
		return nil, protected, fmt.Errorf("pdf: %w", err)
	}
	return ctx, protected, nil
}

// Convert converts a PDF read from rs.
func Convert(rs io.ReadSeeker, opts *Options) (*Result, error) {
	s, err := NewStream(rs, opts)
	if err != nil {
		return nil, err
	}
	return s.Finish()
}

// newConverter reads the PDF and sets up its view: every page selected,
// sized, with a body layer that finalize fills in.
func newConverter(rs io.ReadSeeker, opts *Options) (*converter, bool, error) {
	if opts == nil {
		opts = &Options{}
	}
	ctx, protected, err := readContext(rs, opts.Password)
	if err != nil {
		return nil, protected, err
	}
	if err := ctx.EnsurePageCount(); err != nil {
		return nil, protected, fmt.Errorf("pdf: %w", err)
	}
	c := &converter{pdf: &pdf{ctx: ctx}, doc: bdf.NewDocument(), opts: opts, warned: map[string]bool{},
		fonts: map[string]*pdfFont{}, forms: map[string]*pending{}, images: map[string]*imageEntry{}, shadings: map[string]*shading{}}
	c.doc.Meta.Source = "pdf"
	c.doc.Meta.DC = c.pdf.info()
	if opts.Title != "" {
		c.doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	if cat, err := ctx.Catalog(); err == nil && cat != nil {
		// the catalog's /Lang is the document's language (the default of
		// its text, spec §7.8), unless the metadata already names one
		if lang := strings.TrimSpace(c.pdf.text(cat["Lang"])); lang != "" && len(c.doc.Meta.DC.Language) == 0 {
			c.doc.Meta.DC.Language = bdf.DCValues{lang}
		}
		c.docLang = c.doc.Meta.DC.Language.First()
		c.tree = newStructTree(c.pdf, cat, c.docLang)
	}
	kind := opts.Kind
	if kind == "" {
		kind = bdf.ViewFixed
	}
	view := c.doc.NewView("pages", kind, c.doc.Meta.DC.Title.First())
	if kind == bdf.ViewFlow {
		view.Continuous = &bdf.Continuous{Gap: 16}
	}
	pages := opts.Pages.Numbers(ctx.PageCount)
	if pages == nil {
		for i := 1; i <= ctx.PageCount; i++ {
			pages = append(pages, i)
		}
	}
	c.outPages, c.outPageNrs = map[int]int{}, map[int]int{}
	for i, n := range pages {
		if n < 1 || n > ctx.PageCount {
			return nil, protected, fmt.Errorf("pdf: page %d out of range (1-%d)", n, ctx.PageCount)
		}
		pr := &pageRef{nr: n}
		pr.dict, pr.ref, pr.attrs, pr.err = ctx.PageDict(n, false)
		if pr.err == nil && pr.dict == nil {
			pr.err = fmt.Errorf("page dict missing")
		}
		if _, ok := c.outPageNrs[n]; !ok {
			c.outPageNrs[n] = i + 1
			if pr.err == nil && pr.ref != nil {
				c.outPages[pr.ref.ObjectNumber.Value()] = i + 1
			}
		}
		pr.geometry()
		pr.page = view.AddPage(float32(pr.w), float32(pr.h), bdf.Layer{Role: bdf.RoleBody, Obj: bdf.Hash{}}) // Obj: patched in finalize
		if view.Kind == bdf.ViewFlow {
			pr.page.Body = &bdf.RectDef{X: 0, Y: 0, W: float32(pr.w), H: float32(pr.h)}
		}
		c.pageBodies = append(c.pageBodies, pr)
	}
	return c, protected, nil
}

// finish makes the document once every page is converted.
func (c *converter) finish(protected bool) *Result {
	if !c.opts.NoSharePrefix {
		c.sharePagePrefixes()
	}
	c.finalize()
	if !c.opts.NoTextIndex {
		if _, err := c.doc.BuildTextIndex(c.doc.Views[0]); err != nil {
			c.warnf("text index: %v", err)
		}
	}
	return &Result{Doc: c.doc, Warnings: c.warnings, Pages: len(c.pageBodies), Protected: protected, SharedPrefixes: c.sharedPrefixes, SharedBytes: c.sharedBytes}
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

// geometry sets the page's size and the transform from PDF user space
// (y up) to page space (y down): the crop box, clipped to the media box,
// turned by /Rotate. A page whose dictionary did not load gets a Letter page.
func (pr *pageRef) geometry() {
	box := rect{0, 0, 612, 792}
	rotate := 0
	if attrs := pr.attrs; attrs != nil {
		if attrs.MediaBox != nil {
			box = rect{attrs.MediaBox.LL.X, attrs.MediaBox.LL.Y, attrs.MediaBox.UR.X, attrs.MediaBox.UR.Y}
		}
		if attrs.CropBox != nil {
			cb := rect{attrs.CropBox.LL.X, attrs.CropBox.LL.Y, attrs.CropBox.UR.X, attrs.CropBox.UR.Y}
			if !cb.intersect(box).empty() {
				box = cb.intersect(box)
			}
		}
		rotate = ((attrs.Rotate % 360) + 360) % 360
	}
	if box.empty() || box.x1-box.x0 > 20000 || box.y1-box.y0 > 20000 {
		box = rect{0, 0, 612, 792}
	}
	w, h := box.x1-box.x0, box.y1-box.y0
	// PDF user space → page space (y down), then rotation.
	pr.base = matrix{1, 0, 0, -1, -box.x0, box.y1}
	pr.w, pr.h = w, h
	switch rotate {
	case 90:
		pr.base = pr.base.mul(matrix{0, 1, -1, 0, h, 0})
		pr.w, pr.h = h, w
	case 180:
		pr.base = pr.base.mul(matrix{-1, 0, 0, -1, w, h})
	case 270:
		pr.base = pr.base.mul(matrix{0, -1, 1, 0, 0, w})
		pr.w, pr.h = h, w
	}
}

// convertPage converts one page of the view into its body object.
func (c *converter) convertPage(pr *pageRef) error {
	if pr.err != nil {
		return pr.err
	}
	p := c.pdf
	pageNr, pageDict, pageIRef, attrs := pr.nr, pr.dict, pr.ref, pr.attrs
	pw, ph, base := pr.w, pr.h, pr.base
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
	in.st = &structState{}
	if c.tree != nil {
		sp, ok := p.num(pageDict["StructParents"])
		pageObj := 0
		if pageIRef != nil {
			pageObj = pageIRef.ObjectNumber.Value()
		}
		in.mcids = c.tree.contentTable(int(sp), ok, pageObj)
	}
	in.transform(base)
	userSpace := in.gs.ctm
	if len(content) > 0 {
		in.run(content)
	}
	in.closeStructure()
	in.st = nil
	if !c.opts.NoAnnotations && len(c.pdf.array(pageDict["Annots"])) > 0 {
		// Annotations are in default user space, but the content may leave a
		// cm in effect outside any q (Chrome's does): undo it around them.
		undo := in.gs.ctm != userSpace
		if undo {
			inv, ok := in.gs.ctm.inverse()
			undo = ok
			if ok {
				in.save()
				in.transform(userSpace.mul(inv))
			}
		}
		c.annotations(in, pageDict)
		if undo {
			in.restore()
		}
	}
	pr.body = body
	return nil
}

// pageRef is a page of the view: the PDF page it converts, and its body
// object once converted.
type pageRef struct {
	nr    int // PDF page number
	dict  types.Dict
	ref   *types.IndirectRef
	attrs *model.InheritedPageAttrs
	err   error // the page dictionary did not load

	w, h float64 // size in page space
	base matrix  // PDF user space → page space

	page *bdf.Page
	body *pending // nil until converted
}

// annotations draws annotation appearance streams and emits link areas.
// The current transform is the page's user space.
func (c *converter) annotations(in *interp, pageDict types.Dict) {
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
			if target := c.linkTarget(ad); target != "" {
				in.obj.Link(float32(r.x0), float32(r.y0), float32(r.x1-r.x0), float32(r.y1-r.y0), target)
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

// linkTarget returns the url of a link annotation: a URI action, or
// "#page=N" for a destination on a converted page.
func (c *converter) linkTarget(ad types.Dict) string {
	p := c.pdf
	if act := p.dict(ad["A"]); act != nil {
		switch p.name(act["S"]) {
		case "URI":
			return string(p.str(act["URI"]))
		case "GoTo":
			return c.destTarget(act["D"])
		}
		return ""
	}
	return c.destTarget(ad["Dest"])
}

// destTarget resolves an explicit or named destination to "#page=N".
func (c *converter) destTarget(o types.Object) string {
	p := c.pdf
	for i := 0; i < 4; i++ {
		switch v := p.deref(o).(type) {
		case types.Name:
			o = c.namedDest(v.Value())
		case types.StringLiteral, types.HexLiteral:
			o = c.namedDest(string(p.str(v)))
		case types.Dict:
			o = v["D"]
		case types.Array:
			if len(v) == 0 {
				return ""
			}
			n, ok := 0, false
			switch pg := v[0].(type) {
			case types.IndirectRef:
				n, ok = c.outPages[pg.ObjectNumber.Value()]
			case types.Integer: // a page index, as some producers write
				n, ok = c.outPageNrs[int(pg)+1]
			}
			if !ok {
				return ""
			}
			return fmt.Sprintf("#page=%d", n)
		default:
			return ""
		}
	}
	return ""
}

// namedDest looks a destination name up in the catalog's /Dests name tree
// (PDF 1.2) and /Dests dictionary (PDF 1.1).
func (c *converter) namedDest(name string) types.Object {
	p := c.pdf
	cat, err := p.ctx.Catalog()
	if err != nil || cat == nil {
		return nil
	}
	if v := p.nameTree(p.dict(cat["Names"])["Dests"], name); v != nil {
		return v
	}
	return p.dict(cat["Dests"])[name]
}

// formObject converts a form XObject (or pattern cell) into a shared pending
// object. fs is the caller's structure state (nil: no structure MARKs).
func (c *converter) formObject(key string, sd *types.StreamDict, parentRes types.Dict, depth int, fs *formStruct) *pending {
	p := c.pdf
	if key != "" {
		if f, ok := c.forms[key]; ok {
			if f == nil || f.building {
				return nil // in progress (cycle) or failed
			}
			if f.st == nil || !f.st.used || f.st.reusableFor(fs) {
				return f
			}
			key = "" // its MARKs were made for another state: convert it again
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
	if fs != nil {
		st := fs.entry.clone()
		in.st, in.mcids, in.inherit = &st, fs.mcids, fs.inherit
	}
	in.run(data)
	if fs != nil {
		fs.used, fs.exit = in.structUsed, in.st.clone()
		f.st = fs
	}
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
