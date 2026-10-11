// Package svgbdf writes the pages of BDF documents as SVG files, for where
// no viewer can run: a page is a file a browser displays by itself.
//
// It runs the same instructions as the viewer's Canvas 2D renderer
// (docs/spec.md §7) and writes the elements of SVG they correspond to:
// paths with their fills, strokes and dashes, clip paths, gradients and
// patterns, images, groups with their alpha and blend mode, soft masks,
// shadows and links. Text stays text: text elements in the fonts of the
// document, which a browser lays out, selects and searches, fitted to the
// advance the document states as the viewer fits it. The fonts and images
// are files beside the SVG, named after their parts and referred to by
// relative URL, so that the pages of a document share them (Files, Write);
// Options.Inline puts them into each SVG instead, for a file that stands
// alone.
//
// A browser draws the SVG as it draws the Canvas of the viewer, with these
// differences. Conic gradients have no element: what they paint has the
// color of their first stop. The FILTER instruction is left out. Of the
// composite operations, the blend modes, lighter, destination-out,
// destination-in and destination-over are drawn (the last three as masks
// of what was drawn before, and by drawing first); the others (source-in,
// xor …) are drawn as source-over. A soft mask of luminance is weighted
// as SVG weights it, which is the same for the grays masks are made of.
// The rectangle of a group does not cut what is drawn in it. Warnings
// list what was left out.
//
// An SVG image of the document is always in the SVG that shows it, as a
// data: URL: as a file beside it, it would be a document of the site the
// files are served from, with the scripts it may have.
//
// Fonts beside the SVG are loaded when the SVG is the document a browser
// shows (opened by itself, or in an object or iframe element); an SVG shown
// as an image (an img element, a CSS background) loads no other file, and
// needs Options.Inline. Programs that are no browsers mostly do not load
// fonts by URL at all.
//
// A document decides what it asks of a renderer, so the drawing has limits:
// of the objects drawn again (maxReusedInstructions), of the groups, soft
// masks and states that are open at a time (maxLayers, maxStates) and of
// the gridlines of a sheet (maxGridLines). What passes a limit is left out
// with a warning, and the rest is drawn.
package svgbdf

import (
	"bytes"
	"cmp"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	neturl "net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/image/bmp"
	"golang.org/x/image/webp"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/image/imgconv"
)

// Options configures a Renderer.
type Options struct {
	// Inline puts the fonts and images of a page into its SVG as data:
	// URLs: a file that stands alone, as large as all it shows. Without
	// it they are files beside the SVG (Files).
	Inline bool
	// Background fills the page behind its content (default white). Use
	// color.Transparent for none.
	Background color.Color
	// PageURL returns the URL of the SVG of a page (both from 0), for the
	// links of pages to pages; "" makes no link of one. The default is
	// the names Write gives the files (Name).
	PageURL func(view, page int) string
}

// File is a file the SVGs of a renderer refer to: a font or an image.
type File struct {
	Name string // the relative URL the SVGs have for it
	Data []byte
}

// Renderer writes the views of a document as SVG. It keeps what it decodes
// for the next page; it is not safe for concurrent use.
type Renderer struct {
	doc  *bdf.Document
	opts Options

	objects  map[bdf.Hash]*bdf.ObjectPart
	coll     map[bdf.Hash][]*bdf.Path
	images   map[bdf.Hash]*imageInfo
	fonts    map[bdf.Hash]string
	files    []File
	warnings []string
	warned   map[string]bool
}

// New returns a renderer of a document.
func New(doc *bdf.Document, opts *Options) *Renderer {
	r := &Renderer{
		doc:     doc,
		objects: map[bdf.Hash]*bdf.ObjectPart{},
		coll:    map[bdf.Hash][]*bdf.Path{},
		images:  map[bdf.Hash]*imageInfo{},
		fonts:   map[bdf.Hash]string{},
		warned:  map[string]bool{},
	}
	if opts != nil {
		r.opts = *opts
	}
	return r
}

// Warnings returns what could not be drawn as the viewer would, once each.
func (r *Renderer) Warnings() []string { return r.warnings }

// Files returns the fonts and images the SVGs written so far refer to, to
// be stored beside them under their names. There are none with
// Options.Inline.
func (r *Renderer) Files() []File { return r.files }

func (r *Renderer) warnf(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	if !r.warned[s] {
		r.warned[s] = true
		r.warnings = append(r.warnings, s)
	}
}

func (r *Renderer) object(h bdf.Hash) *bdf.ObjectPart {
	if o, ok := r.objects[h]; ok {
		return o
	}
	r.objects[h] = nil
	p := r.doc.Part(h)
	if p == nil {
		r.warnf("object %s is missing", h)
		return nil
	}
	o, err := bdf.DecodeObject(p.Data)
	if err != nil {
		r.warnf("object %s: %v", h, err)
		return nil
	}
	r.objects[h] = o
	return o
}

func (r *Renderer) path(o *bdf.ObjectPart, i int) *bdf.Path {
	e := o.Paths[i]
	if e.Inline != nil {
		return e.Inline
	}
	ps, ok := r.coll[e.Hash]
	if !ok {
		if part := r.doc.Part(e.Hash); part != nil {
			var err error
			if ps, err = bdf.DecodePathCollection(part.Data); err != nil {
				r.warnf("path collection %s: %v", e.Hash, err)
			}
		} else {
			r.warnf("path collection %s is missing", e.Hash)
		}
		r.coll[e.Hash] = ps
	}
	if int(e.Index) < len(ps) {
		return ps[e.Index]
	}
	return nil
}

// url returns the URL an SVG has for a part: a data: URL of it, or the
// name of a file beside the SVG, which it adds to the files.
func (r *Renderer) url(h bdf.Hash, data []byte, ext, mime string, inline bool) string {
	if inline {
		return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data)
	}
	name := h.String() + "." + ext
	r.files = append(r.files, File{Name: name, Data: data})
	return name
}

// font returns the URL of an embedded font, or "" for one that is missing.
func (r *Renderer) font(h bdf.Hash) string {
	if u, ok := r.fonts[h]; ok {
		return u
	}
	u := ""
	if p := r.doc.Part(h); p == nil {
		r.warnf("font %s is missing", h)
	} else {
		ext, mime := "ttf", "font/ttf"
		switch {
		case bytes.HasPrefix(p.Data, []byte("wOF2")):
			ext, mime = "woff2", "font/woff2"
		case bytes.HasPrefix(p.Data, []byte("wOFF")):
			ext, mime = "woff", "font/woff"
		case bytes.HasPrefix(p.Data, []byte("OTTO")):
			ext, mime = "otf", "font/otf"
		case bytes.HasPrefix(p.Data, []byte("ttcf")):
			ext, mime = "ttc", "font/collection"
		}
		u = r.url(h, p.Data, ext, mime, r.opts.Inline)
	}
	r.fonts[h] = u
	return u
}

// imageInfo is an image as an SVG refers to it.
type imageInfo struct {
	hash bdf.Hash
	url  string
	w, h float64 // its size in pixels, 0 when it cannot be told
}

var configs = map[string]func(io.Reader) (image.Config, error){
	"png": png.DecodeConfig, "jpeg": jpeg.DecodeConfig, "gif": gif.DecodeConfig, "bmp": bmp.DecodeConfig, "webp": webp.DecodeConfig,
}

// image returns an image part, or nil for one that is missing.
func (r *Renderer) image(h bdf.Hash) *imageInfo {
	if info, ok := r.images[h]; ok {
		return info
	}
	r.images[h] = nil
	p := r.doc.Part(h)
	if p == nil {
		r.warnf("image %s is missing", h)
		return nil
	}
	data := p.Data
	format := imgconv.Sniff(data)
	if format == "" {
		if t := bytes.TrimLeft(bytes.TrimPrefix(data, []byte("\xef\xbb\xbf")), " \t\r\n"); len(t) > 0 && t[0] == '<' {
			format = "svg"
		}
	}
	info := &imageInfo{hash: h}
	switch format {
	case "":
		r.warnf("an image of an unknown format is not drawn")
		return nil
	case "svg":
		// the size of the root, a missing one following from the other
		// and the view box, else the view box's, else 300 × 150
		info.w, info.h = 300, 150
		if width, height, viewBox, ok := imgconv.SVGRoot(data); ok {
			s := imgconv.ParseSVGSize(width, height, viewBox)
			switch {
			case s.W > 0 && s.H > 0:
				info.w, info.h = s.W, s.H
			case s.ViewW > 0 && s.ViewH > 0:
				info.w, info.h = s.ViewW, s.ViewH
			}
		}
	case "avif":
		info.w, info.h = avifSize(data)
	default:
		func() {
			defer func() { recover() }() // of a header reader, on a damaged file
			if c, err := configs[format](bytes.NewReader(data)); err == nil && c.Width > 0 && c.Height > 0 {
				info.w, info.h = float64(c.Width), float64(c.Height)
				if imgconv.Orientation(data) >= 5 {
					info.w, info.h = info.h, info.w // turned on its side, as browsers show it
				}
			}
		}()
	}
	ext := format
	if format == "jpeg" {
		ext = "jpg"
	}
	// An SVG image is always in the SVG that shows it: as a file beside
	// it, it would be a document of the site the files are served from,
	// with the scripts it may have.
	info.url = r.url(h, data, ext, imgconv.MIME(format), r.opts.Inline || format == "svg")
	r.images[h] = info
	return info
}

// avifSize returns the size the header of an AVIF image states: that of
// the largest of its image items (an image of tiles has those too).
func avifSize(data []byte) (w, h float64) {
	head := data[:min(len(data), 1<<16)]
	for {
		i := bytes.Index(head, []byte("ispe"))
		if i < 0 || i+16 > len(head) {
			return w, h
		}
		iw, ih := float64(binary.BigEndian.Uint32(head[i+8:])), float64(binary.BigEndian.Uint32(head[i+12:]))
		if iw*ih > w*h {
			w, h = iw, ih
		}
		head = head[i+16:]
	}
}

// Name returns the name Write gives the SVG of a page of a view (both from
// 0): page-N.svg in a document of one view, view-V-page-N.svg in one of
// more. A sheet is one page.
func Name(doc *bdf.Document, view, page int) string {
	if len(doc.Views) <= 1 {
		return "page-" + strconv.Itoa(page+1) + ".svg"
	}
	return "view-" + strconv.Itoa(view+1) + "-page-" + strconv.Itoa(page+1) + ".svg"
}

// href returns the URL an SVG has for the target of a LINK, or "" for
// what is no link (docs/spec.md §7.7).
func (r *Renderer) href(view int, url string) string {
	if rest, ok := strings.CutPrefix(url, "#"); ok {
		q, err := neturl.ParseQuery(rest)
		if err != nil || !q.Has("view") && !q.Has("page") {
			return ""
		}
		if q.Has("view") {
			view = slices.IndexFunc(r.doc.Views, func(v *bdf.View) bool { return v.ID == q.Get("view") })
		}
		page := 0
		if q.Has("page") {
			n, err := strconv.Atoi(q.Get("page"))
			if err != nil || n < 1 {
				return ""
			}
			page = n - 1
		}
		if view < 0 {
			return ""
		}
		if r.opts.PageURL != nil {
			return r.opts.PageURL(view, page)
		}
		return Name(r.doc, view, page)
	}
	lower := strings.ToLower(url[:min(len(url), 8)])
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") || strings.HasPrefix(lower, "mailto:") {
		return url
	}
	return ""
}

// --- paints ---

// paintServer returns the value of fill or stroke for a gradient or
// pattern, which it defines at its first use: a reference to it, or none
// for one that paints nothing.
func (d *drawer) paintServer(p *paint) string {
	if v, ok := d.paints[p.p]; ok {
		return v
	}
	v := d.definePaint(p)
	d.paints[p.p] = v
	return v
}

func (d *drawer) definePaint(p *paint) string {
	pt := p.p
	c := pt.Coords
	for _, v := range c {
		if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
			return "none"
		}
	}
	at := func(i int) string { return factor(float64(c[i])) }
	var open, end string
	switch pt.Kind {
	case bdf.PaintLinear:
		if len(c) < 4 || c[0] == c[2] && c[1] == c[3] {
			return "none" // Canvas paints nothing with it
		}
		open = `<linearGradient id="%s" gradientUnits="userSpaceOnUse" x1="` + at(0) + `" y1="` + at(1) + `" x2="` + at(2) + `" y2="` + at(3) + `">`
		end = "</linearGradient>\n"
	case bdf.PaintRadial:
		if len(c) < 6 || c[2] < 0 || c[5] < 0 || c[0] == c[3] && c[1] == c[4] && c[2] == c[5] {
			return "none"
		}
		// the circles of Canvas: the end circle, and the start circle as
		// the focal one
		open = `<radialGradient id="%s" gradientUnits="userSpaceOnUse" cx="` + at(3) + `" cy="` + at(4) + `" r="` + at(5) +
			`" fx="` + at(0) + `" fy="` + at(1) + `"`
		if c[2] != 0 {
			open += ` fr="` + at(2) + `"`
		}
		open += ">"
		end = "</radialGradient>\n"
	case bdf.PaintConic:
		d.r.warnf("conic gradients are not drawn: what they paint has the color of their first stop")
		if len(pt.Stops) == 0 {
			return "none"
		}
		return hexColor(pt.Stops[0].Color)
	case bdf.PaintPattern:
		return d.definePattern(p)
	default:
		return "none"
	}
	if len(pt.Stops) == 0 {
		return "none"
	}
	// in the order of their offsets, as Canvas keeps them
	stops := slices.Clone(pt.Stops)
	for i := range stops {
		stops[i].Offset = float32(math.Min(1, math.Max(0, float64(stops[i].Offset))))
	}
	slices.SortStableFunc(stops, func(a, b bdf.Stop) int { return cmp.Compare(a.Offset, b.Offset) })
	id := d.id("g")
	fmt.Fprintf(&d.defs, open, id)
	for _, s := range stops {
		d.defs.WriteString(`<stop offset="` + num(float64(s.Offset), 5) + `" stop-color="` + hexColor(s.Color) + `"`)
		if a := uint32(s.Color) & 0xff; a != 0xff {
			d.defs.WriteString(` stop-opacity="` + factor(float64(a)/255) + `"`)
		}
		d.defs.WriteString("/>")
	}
	d.defs.WriteString(end)
	return "url(#" + id + ")"
}

// unrepeated is the size of the cell of a pattern in a direction it is not
// repeated in: further than a page reaches.
const unrepeated = 1 << 15

func (d *drawer) definePattern(p *paint) string {
	pt := p.p
	if int(pt.Image) >= len(p.obj.Images) {
		d.r.warnf("bad image reference")
		return "none"
	}
	info := d.r.image(p.obj.Images[pt.Image])
	if info == nil {
		return "none"
	}
	if info.w == 0 {
		d.r.warnf("a pattern of an image whose size cannot be told is not drawn")
		return "none"
	}
	m := matrix{}
	for i, v := range pt.Matrix {
		m[i] = float64(v)
	}
	if !m.finite() {
		return "none"
	}
	w, h := info.w, info.h
	if pt.Repeat == bdf.RepeatY || pt.Repeat == bdf.NoRepeat {
		w = max(w, unrepeated)
	}
	if pt.Repeat == bdf.RepeatX || pt.Repeat == bdf.NoRepeat {
		h = max(h, unrepeated)
	}
	id := d.id("g")
	tf := ""
	if m != identity {
		tf = ` patternTransform="matrix(` + factor(m[0]) + " " + factor(m[1]) + " " + factor(m[2]) + " " + factor(m[3]) + " " + factor(m[4]) + " " + factor(m[5]) + `)"`
	}
	// the image first: it may define a symbol
	el := d.imageElement(info, 0, 0, info.w, info.h, 3, "")
	d.defs.WriteString(`<pattern id="` + id + `" patternUnits="userSpaceOnUse" width="` + num(w, 3) + `" height="` + num(h, 3) + `"` + tf + ">" + el + "</pattern>\n")
	return "url(#" + id + ")"
}

// family returns the font-family of a font: the family of its embedded
// font, which the SVG declares at its first use without weight or style
// (so that those are synthesized, as in the viewer), before the families
// the document names.
func (d *drawer) family(f *bdf.Font) string {
	if f.Kind == bdf.FontEmbedded {
		name, ok := d.fonts[f.Hash]
		if !ok {
			if u := d.r.font(f.Hash); u != "" {
				name = "bdf-" + f.Hash.String()[:10]
				d.css.WriteString(`@font-face{font-family:"` + name + `";src:url("` + u + `")}` + "\n")
			}
			d.fonts[f.Hash] = name
		}
		if name != "" {
			if f.Family != "" {
				return name + "," + f.Family
			}
			return name
		}
	}
	if f.Family == "" {
		return "sans-serif"
	}
	return f.Family
}

// --- pages ---

// ErrNoPage is returned for a page a view does not have.
var ErrNoPage = errors.New("svgbdf: no such page")

// recovered makes an error of a panic of a drawing. The limits of this
// package keep a document from what documents are known to ask for; this
// keeps what nobody thought of from ending the program.
func recovered(svg *[]byte, err *error) {
	if p := recover(); p != nil {
		*svg, *err = nil, fmt.Errorf("svgbdf: the drawing failed: %v", p)
	}
}

// start begins the SVG of a rectangle of a view, in units.
func (r *Renderer) start(v *bdf.View, x, y, w, h float64) *drawer {
	return &drawer{
		r: r, view: slices.Index(r.doc.Views, v), x: x, y: y, w: w, h: h,
		ctx:     &context{st: initialState(identity, 0), buf: &bytes.Buffer{}},
		clips:   map[string]int{},
		paints:  map[*bdf.Paint]string{},
		glyphs:  map[glyphKey]string{},
		filters: map[string]string{},
		fonts:   map[bdf.Hash]string{},
		symbols: map[bdf.Hash]string{},
	}
}

// background fills a rectangle with the background of the pages.
func (d *drawer) background(x, y, w, h float64) {
	var c color.Color = color.White
	if d.r.opts.Background != nil {
		c = d.r.opts.Background
	}
	n := color.NRGBAModel.Convert(c).(color.NRGBA)
	if n.A == 0 {
		return
	}
	data, _ := rectData(x, y, w, h, pageDigits)
	el := `<path d="` + data + `" fill="` + hexColor(bdf.RGBA(n.R, n.G, n.B, 255)) + `"`
	if n.A != 0xff {
		el += ` fill-opacity="` + factor(float64(n.A)/255) + `"`
	}
	d.write(d.ctx, d.ctx.st.clip, el+"/>")
}

// clipRect returns the clip path of a rectangle under m.
func (d *drawer) clipRect(m matrix, x, y, w, h float64) int {
	st := &d.ctx.st
	*st = initialState(m, 0)
	data, _ := rectData(x, y, w, h, m.digits())
	d.clip(data, bdf.NonZero)
	return st.clip
}

// finish returns the SVG a drawer has drawn.
func (d *drawer) finish() []byte {
	root := d.ctx
	d.closeAll(root)
	d.settle(root)
	var b bytes.Buffer
	b.Grow(root.buf.Len() + d.defs.Len() + d.css.Len() + d.links.Len() + 256)
	b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="` + num(d.x, pageDigits) + " " + num(d.y, pageDigits) + " " + num(d.w, pageDigits) + " " + num(d.h, pageDigits) +
		`" width="` + num(d.w, pageDigits) + `pt" height="` + num(d.h, pageDigits) + `pt"`)
	if root.blended {
		b.WriteString(` style="isolation:isolate"`)
	}
	b.WriteString(">\n")
	if d.css.Len() > 0 || d.defs.Len() > 0 {
		b.WriteString("<defs>\n")
		if d.css.Len() > 0 {
			b.WriteString("<style>\n")
			b.Write(d.css.Bytes())
			b.WriteString("</style>\n")
		}
		b.Write(d.defs.Bytes())
		b.WriteString("</defs>\n")
	}
	b.Write(root.buf.Bytes())
	b.Write(d.links.Bytes())
	b.WriteString("</svg>\n")
	return b.Bytes()
}

func sized(w, h float32) bool {
	return w > 0 && h > 0 && !math.IsInf(float64(w), 0) && !math.IsInf(float64(h), 0)
}

// Page writes a page of a fixed or flow view, or a strip of a scroll view
// cut to its body, as the viewer scrolls the strips: the SVGs of the
// strips stacked are the view.
func (r *Renderer) Page(v *bdf.View, page int) (svg []byte, err error) {
	defer recovered(&svg, &err)
	if v == nil {
		return nil, errors.New("svgbdf: no view")
	}
	if v.Kind == bdf.ViewSheet {
		return nil, errors.New("svgbdf: a sheet has no pages; draw a Region")
	}
	if page < 0 || page >= len(v.Pages) || v.Pages[page] == nil {
		return nil, ErrNoPage
	}
	p := v.Pages[page]
	if v.Kind == bdf.ViewScroll {
		b := bdf.RectDef{W: p.W, H: p.H}
		if p.Body != nil {
			b = *p.Body
		}
		if !sized(b.W, b.H) {
			return nil, errors.New("svgbdf: empty page")
		}
		d := r.start(v, float64(b.X), float64(b.Y), float64(b.W), float64(b.H))
		d.background(d.x, d.y, d.w, d.h)
		for _, l := range p.Layers {
			if l.Role == bdf.RoleBody || l.Role == bdf.RoleAnnotation {
				d.drawTop(l.Obj, identity, 0)
			}
		}
		return d.finish(), nil
	}
	if !sized(p.W, p.H) {
		return nil, errors.New("svgbdf: empty page")
	}
	return r.drawPage(v, p, bdf.Rect{W: p.W, H: p.H}), nil
}

// Region writes a rectangle of a view: of the page with index page in a
// fixed or flow view (in page units), of the continuous layout of a scroll
// view (its strips stacked, in units from the top), or of a sheet (page is
// ignored in the last two).
func (r *Renderer) Region(v *bdf.View, page int, region bdf.Rect) (svg []byte, err error) {
	defer recovered(&svg, &err)
	if v == nil {
		return nil, errors.New("svgbdf: no view")
	}
	if !sized(region.W, region.H) || math.IsNaN(float64(region.X+region.Y)) || math.IsInf(float64(region.X+region.Y), 0) {
		return nil, errors.New("svgbdf: empty region")
	}
	switch v.Kind {
	case bdf.ViewSheet:
		return r.drawSheet(v, region), nil
	case bdf.ViewScroll:
		return r.drawContinuous(v, region), nil
	}
	if page < 0 || page >= len(v.Pages) || v.Pages[page] == nil {
		return nil, ErrNoPage
	}
	return r.drawPage(v, v.Pages[page], region), nil
}

func (r *Renderer) drawPage(v *bdf.View, p *bdf.Page, region bdf.Rect) []byte {
	d := r.start(v, float64(region.X), float64(region.Y), float64(region.W), float64(region.H))
	// the SVG shows no more than the region: a clip only where that is
	// more than the page
	clip := 0
	if region.X < 0 || region.Y < 0 || region.X+region.W > p.W || region.Y+region.H > p.H {
		clip = d.clipRect(identity, 0, 0, float64(p.W), float64(p.H))
	}
	d.background(0, 0, float64(p.W), float64(p.H))
	for _, l := range p.Layers {
		d.drawTop(l.Obj, identity, clip)
	}
	return d.finish()
}

// drawContinuous writes the strips of a scroll view stacked without gaps,
// each cut to its body, as the viewer scrolls them.
func (r *Renderer) drawContinuous(v *bdf.View, region bdf.Rect) []byte {
	d := r.start(v, float64(region.X), float64(region.Y), float64(region.W), float64(region.H))
	y := 0.0
	for _, p := range v.Pages {
		if p == nil {
			continue
		}
		b := bdf.RectDef{W: p.W, H: p.H}
		if p.Body != nil {
			b = *p.Body
		}
		top := y
		y += float64(b.H)
		if top >= float64(region.Y+region.H) || y <= float64(region.Y) {
			continue
		}
		clip := d.clipRect(identity.translate(0, top), 0, 0, float64(b.W), float64(b.H))
		d.background(0, top, float64(b.W), float64(b.H))
		m := identity.translate(-float64(b.X), top-float64(b.Y))
		for _, l := range p.Layers {
			if l.Role == bdf.RoleBody || l.Role == bdf.RoleAnnotation {
				d.drawTop(l.Obj, m, clip)
			}
		}
	}
	return d.finish()
}

// SheetSize returns the size of a sheet in units, from its rows and columns.
func SheetSize(v *bdf.View) (w, h float64) {
	for _, c := range v.Cols {
		w += float64(c[0]) * float64(c[1])
	}
	for _, rr := range v.Rows {
		h += float64(rr[0]) * float64(rr[1])
	}
	return w, h
}

// tilePlace returns the place of a tile from its key ("x,y").
func tilePlace(k string) (tx, ty int, ok bool) {
	x, y, _ := strings.Cut(k, ",")
	tx, errx := strconv.Atoi(x)
	ty, erry := strconv.Atoi(y)
	return tx, ty, errx == nil && erry == nil && tx >= 0 && ty >= 0 && strconv.Itoa(tx)+","+strconv.Itoa(ty) == k
}

func tileSize(v *bdf.View) float64 {
	if v.Tile <= 0 {
		return bdf.DefaultTile
	}
	return float64(v.Tile)
}

// SheetBounds returns the rectangle of a sheet that has content, from A1:
// as far as its tiles reach, within its rows and columns.
func SheetBounds(v *bdf.View) bdf.Rect {
	tile := tileSize(v)
	var w, h float64
	for k := range v.Tiles {
		if tx, ty, ok := tilePlace(k); ok {
			w, h = math.Max(w, float64(tx+1)*tile), math.Max(h, float64(ty+1)*tile)
		}
	}
	if sw, sh := SheetSize(v); sw > 0 && sh > 0 {
		w, h = math.Min(w, sw), math.Min(h, sh)
	}
	return bdf.Rect{W: float32(w), H: float32(h)}
}

// drawSheet writes a region of a sheet: the background, the gridlines when
// the view asks for them, and the tiles.
func (r *Renderer) drawSheet(v *bdf.View, region bdf.Rect) []byte {
	d := r.start(v, float64(region.X), float64(region.Y), float64(region.W), float64(region.H))
	d.background(d.x, d.y, d.w, d.h)
	if v.Gridlines {
		r.gridlines(d, v)
	}
	tile := tileSize(v)
	if !(tile >= 1) || math.IsInf(tile, 0) {
		// a region would hold any number of them
		r.warnf("the tiles of a sheet are %g units wide: they are not drawn", tile)
		return d.finish()
	}
	type place struct {
		tx, ty int
		hash   string
	}
	var places []place
	for k, hs := range v.Tiles {
		tx, ty, ok := tilePlace(k)
		if !ok || float64(tx)*tile >= d.x+d.w || float64(tx+1)*tile <= d.x || float64(ty)*tile >= d.y+d.h || float64(ty+1)*tile <= d.y {
			continue
		}
		places = append(places, place{tx, ty, hs})
	}
	slices.SortFunc(places, func(a, b place) int { return cmp.Or(cmp.Compare(a.ty, b.ty), cmp.Compare(a.tx, b.tx)) })
	for _, p := range places {
		hash, err := bdf.ParseHash(p.hash)
		if err != nil {
			r.warnf("bad tile hash %q", p.hash)
			continue
		}
		m := identity.translate(float64(p.tx)*tile, float64(p.ty)*tile)
		d.drawTop(hash, m, d.clipRect(m, 0, 0, tile, tile))
	}
	return d.finish()
}

// Limits of the gridlines of a sheet, across and down: a document decides
// how many rows and columns a sheet has.
const (
	// maxGridSteps bounds the rows or columns gone through for their lines.
	maxGridSteps = 1 << 24
	// maxGridLines bounds the lines drawn.
	maxGridLines = 1 << 16
)

// gridlines draws the lines between the cells of a sheet as the viewer
// does: a pixel wide at every size, in #d9d9d9.
func (r *Renderer) gridlines(d *drawer, v *bdf.View) {
	sw, sh := SheetSize(v)
	x1, y1 := math.Min(d.x+d.w, sw), math.Min(d.y+d.h, sh)
	var b []byte
	// walk goes through the rows or columns from the first, and draws the
	// line after each
	walk := func(runs []bdf.Run, from, to float64, line func(pos float64)) {
		pos, steps, lines, last := 0.0, 0, 0, math.NaN()
		for _, run := range runs {
			n, size := float64(run[0]), float64(run[1])
			if !(n >= 1) {
				continue
			}
			for i := 0; float64(i) < math.Floor(n); i++ {
				if steps++; steps > maxGridSteps {
					r.warnf("a sheet has more than %d rows or columns: the gridlines of the rest are not drawn", maxGridSteps)
					return
				}
				pos += size
				if pos > to {
					return
				}
				if !(pos >= from) || pos == last {
					if size == 0 {
						break // so are the rest of the run
					}
					continue
				}
				if lines++; lines > maxGridLines {
					r.warnf("a sheet has more than %d gridlines in a region: the rest are not drawn", maxGridLines)
					return
				}
				line(pos)
				last = pos
			}
		}
	}
	walk(v.Cols, d.x, x1, func(x float64) {
		b = append(b, 'M')
		b = appendNum(b, x, pageDigits)
		b = append(b, ' ')
		b = appendNum(b, d.y, pageDigits)
		b = append(b, 'V')
		b = appendNum(b, y1, pageDigits)
	})
	walk(v.Rows, d.y, y1, func(y float64) {
		b = append(b, 'M')
		b = appendNum(b, d.x, pageDigits)
		b = append(b, ' ')
		b = appendNum(b, y, pageDigits)
		b = append(b, 'H')
		b = appendNum(b, x1, pageDigits)
	})
	if len(b) > 0 {
		d.write(d.ctx, 0, `<path d="`+string(b)+`" fill="none" stroke="#d9d9d9" vector-effect="non-scaling-stroke" shape-rendering="crispEdges"/>`)
	}
}

// --- documents ---

// Result is what Write wrote.
type Result struct {
	// Pages are the names of the SVGs, in the order of the views and their
	// pages.
	Pages []string
	// Files are the names of the fonts and images beside them.
	Files []string
	// Warnings list what could not be drawn as the viewer would.
	Warnings []string
}

// Write writes every page of a document into a directory, which it makes:
// an SVG a page (Name), a sheet as one page of all it holds (SheetBounds),
// and the fonts and images they share.
func Write(dir string, doc *bdf.Document, opts *Options) (*Result, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	r := New(doc, opts)
	res := &Result{}
	put := func(view, page int, svg []byte, err error) error {
		if err != nil {
			return fmt.Errorf("view %d, page %d: %w", view+1, page+1, err)
		}
		name := Name(doc, view, page)
		res.Pages = append(res.Pages, name)
		return os.WriteFile(filepath.Join(dir, name), svg, 0o644)
	}
	for vi, v := range doc.Views {
		if v.Kind == bdf.ViewSheet {
			if b := SheetBounds(v); sized(b.W, b.H) {
				svg, err := r.Region(v, 0, b)
				if err := put(vi, 0, svg, err); err != nil {
					return nil, err
				}
			}
			continue
		}
		for pi, p := range v.Pages {
			if p == nil {
				continue
			}
			svg, err := r.Page(v, pi)
			if err := put(vi, pi, svg, err); err != nil {
				return nil, err
			}
		}
	}
	for _, f := range r.Files() {
		res.Files = append(res.Files, f.Name)
		if err := os.WriteFile(filepath.Join(dir, f.Name), f.Data, 0o644); err != nil {
			return nil, err
		}
	}
	res.Warnings = r.Warnings()
	return res, nil
}
