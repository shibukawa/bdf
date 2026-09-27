// Package raster draws BDF documents into images without a browser: the
// pages of fixed and flow views, stretches of scroll views and regions of
// sheets, for thumbnails and other server-side previews.
//
// It runs the same instructions as the viewer's Canvas 2D renderer
// (docs/spec.md §7) with a software rasterizer: anti-aliased paths with
// nonzero and even-odd fills, strokes with joins, caps and dashes, clips,
// gradients and patterns, images (PNG, JPEG, GIF, BMP, WebP), text in the
// embedded fonts (WOFF2, TrueType, OpenType) and in the fonts of the
// system for fonts referred to by name, groups with their alpha and blend
// mode, soft masks and shadows. The result is close to what a browser
// draws but not identical: text is not hinted, shaped (no kerning,
// ligatures or Arabic joining) or laid out by a full bidirectional
// algorithm; SVG and AVIF images and FILTER are not drawn. Warnings list
// what was left out.
package raster

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io/fs"
	"math"
	"strconv"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/fontdb"
)

// Options configures a Renderer.
type Options struct {
	// FontFS holds fonts searched before FontDirs, for the text of fonts
	// referred to by name (HTML, Markdown and EPUB documents, and others
	// converted with system fonts) and for characters an embedded font
	// lacks.
	FontFS fs.FS
	// FontDirs are searched for fonts before the system font directories.
	FontDirs []string
	// NoSystemFonts restricts font lookup to FontFS and FontDirs.
	NoSystemFonts bool
	// Background fills the page behind its content (default white). Use
	// color.Transparent for a transparent page.
	Background color.Color
}

// Renderer draws the views of a document. It keeps what it decodes (objects,
// paths, images, fonts) for the next drawing; it is not safe for
// concurrent use.
type Renderer struct {
	doc  *bdf.Document
	opts Options

	objects  map[bdf.Hash]*bdf.ObjectPart
	coll     map[bdf.Hash][]*bdf.Path
	paths    map[pathKey]*path
	pictures map[bdf.Hash]*picture
	svgs     map[bdf.Hash]*svgImage
	fonts    *fonts
	warnings []string
	warned   map[string]bool
}

type pathKey struct {
	obj *bdf.ObjectPart
	i   int
}

// New returns a renderer of a document.
func New(doc *bdf.Document, opts *Options) *Renderer {
	r := &Renderer{
		doc:      doc,
		objects:  map[bdf.Hash]*bdf.ObjectPart{},
		coll:     map[bdf.Hash][]*bdf.Path{},
		paths:    map[pathKey]*path{},
		pictures: map[bdf.Hash]*picture{},
		svgs:     map[bdf.Hash]*svgImage{},
		warned:   map[string]bool{},
	}
	if opts != nil {
		r.opts = *opts
	}
	r.fonts = &fonts{r: r, parts: map[bdf.Hash]*fontFace{}, files: map[*fontdb.Face]*fontFace{}, chains: map[bdf.Font]*fontChain{}}
	return r
}

// Warnings returns what could not be drawn as the viewer would, once each.
func (r *Renderer) Warnings() []string { return r.warnings }

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

func (r *Renderer) path(o *bdf.ObjectPart, i int) *path {
	k := pathKey{o, i}
	if p, ok := r.paths[k]; ok {
		return p
	}
	e := o.Paths[i]
	var bp *bdf.Path
	if e.Inline != nil {
		bp = e.Inline
	} else {
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
			bp = ps[e.Index]
		}
	}
	var p *path
	if bp != nil {
		p = buildPath(bp)
	}
	r.paths[k] = p
	return p
}

// picture returns a decoded bitmap image part, or nil (for an SVG image,
// see svg, and for an image that cannot be decoded).
func (r *Renderer) picture(h bdf.Hash) *picture {
	if p, ok := r.pictures[h]; ok {
		return p
	}
	r.pictures[h] = nil
	part := r.doc.Part(h)
	if part == nil {
		r.warnf("image %s is missing", h)
		return nil
	}
	if isSVG(part.Data) {
		si, ok := newSVGImage(part.Data)
		if !ok {
			r.warnf("an SVG image cannot be read")
			return nil
		}
		r.svgs[h] = si
		return nil
	}
	p, format, ok := decodePicture(part.Data)
	if !ok {
		switch format {
		case "avif":
			r.warnf("AVIF images are not drawn")
		case "too large":
			r.warnf("an image of more than %d pixels is not drawn", maxImagePixels)
		default:
			r.warnf("an image (%s) cannot be decoded", format)
		}
		return nil
	}
	r.pictures[h] = p
	return p
}

// svg returns an SVG image part, or nil.
func (r *Renderer) svg(h bdf.Hash) *svgImage {
	r.picture(h)
	return r.svgs[h]
}

// ErrNoPage is returned for a page a view does not have.
var ErrNoPage = errors.New("raster: no such page")

// MaxPixels is the largest image Page and Region draw (a canvas takes 16
// bytes a pixel while it is drawn).
const MaxPixels = 64 << 20

// ErrTooLarge is returned for an image of more than MaxPixels pixels.
var ErrTooLarge = errors.New("raster: the image would be too large")

// Page draws a whole page of a fixed or flow view (or a strip of a scroll
// view) at scale device pixels per unit.
func (r *Renderer) Page(v *bdf.View, page int, scale float64) (*image.RGBA, error) {
	if v.Kind == bdf.ViewSheet {
		return nil, errors.New("raster: a sheet has no pages; draw a Region")
	}
	if page < 0 || page >= len(v.Pages) {
		return nil, ErrNoPage
	}
	p := v.Pages[page]
	fw, fh := math.Ceil(float64(p.W)*scale-1e-6), math.Ceil(float64(p.H)*scale-1e-6)
	if !(fw*fh <= MaxPixels) {
		return nil, ErrTooLarge
	}
	return r.drawPage(p, bdf.Rect{W: p.W, H: p.H}, max(int(fw), 1), max(int(fh), 1))
}

// Region draws a rectangle of a view into a w × h image: of the page with
// index page in a fixed or flow view (in page units), of the continuous
// layout of a scroll view (its strips stacked, in units from the top), or
// of a sheet (page is ignored in the last two).
func (r *Renderer) Region(v *bdf.View, page int, region bdf.Rect, w, h int) (*image.RGBA, error) {
	if w <= 0 || h <= 0 || !(region.W > 0) || !(region.H > 0) {
		return nil, errors.New("raster: empty region")
	}
	if int64(w)*int64(h) > MaxPixels {
		return nil, ErrTooLarge
	}
	switch v.Kind {
	case bdf.ViewSheet:
		return r.drawSheet(v, region, w, h), nil
	case bdf.ViewScroll:
		return r.drawContinuous(v, region, w, h), nil
	}
	if page < 0 || page >= len(v.Pages) {
		return nil, ErrNoPage
	}
	return r.drawPage(v.Pages[page], region, w, h)
}

// canvas starts an image for a region with the device transform that maps
// the region onto it.
func (r *Renderer) canvas(region bdf.Rect, w, h int) (*drawer, matrix) {
	sx, sy := float64(w)/float64(region.W), float64(h)/float64(region.H)
	m := matrix{sx, 0, 0, sy, -float64(region.X) * sx, -float64(region.Y) * sy}
	s := newSurface(w, h)
	d := &drawer{r: r, ctx: &context{target: s}}
	d.ctx.st = initialState(m, &mask{r: s.bounds()})
	return d, m
}

func (r *Renderer) background() solid {
	if r.opts.Background == nil {
		return solid{1, 1, 1, 1}
	}
	c, g, b, a := r.opts.Background.RGBA()
	return solid{float32(c) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff}
}

// clipTo returns the coverage of a rectangle in units under m, within the
// canvas.
func (d *drawer) clipTo(m matrix, x, y, w, h float64) *mask {
	ax, ay := m.apply(x, y)
	bx, by := m.apply(x+w, y+h)
	return rectMask(ax, ay, bx, by, d.ctx.target.bounds())
}

func (r *Renderer) drawPage(p *bdf.Page, region bdf.Rect, w, h int) (*image.RGBA, error) {
	d, m := r.canvas(region, w, h)
	clip := d.clipTo(m, 0, 0, float64(p.W), float64(p.H))
	d.ctx.target.fill(clip, r.background(), 1, bdf.BlendSourceOver, &mask{r: d.ctx.target.bounds()})
	for _, l := range p.Layers {
		if o := r.object(l.Obj); o != nil {
			d.drawTop(o, m, clip)
		}
	}
	return d.ctx.target.toRGBA(), nil
}

// drawContinuous draws the strips of a scroll view stacked without gaps,
// each cut to its body, as the viewer scrolls them.
func (r *Renderer) drawContinuous(v *bdf.View, region bdf.Rect, w, h int) *image.RGBA {
	d, m := r.canvas(region, w, h)
	bg := r.background()
	y := 0.0
	for _, p := range v.Pages {
		b := bdf.RectDef{W: p.W, H: p.H}
		if p.Body != nil {
			b = *p.Body
		}
		top := y
		y += float64(b.H)
		if top >= float64(region.Y+region.H) || y <= float64(region.Y) {
			continue
		}
		pm := m.translate(0, top)
		clip := d.clipTo(pm, 0, 0, float64(b.W), float64(b.H))
		d.ctx.target.fill(clip, bg, 1, bdf.BlendSourceOver, &mask{r: d.ctx.target.bounds()})
		pm = pm.translate(-float64(b.X), -float64(b.Y))
		for _, l := range p.Layers {
			if l.Role != bdf.RoleBody && l.Role != bdf.RoleAnnotation {
				continue
			}
			if o := r.object(l.Obj); o != nil {
				d.drawTop(o, pm, clip)
			}
		}
	}
	return d.ctx.target.toRGBA()
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

// drawSheet draws a region of a sheet: the background, the gridlines when
// the view asks for them, and the tiles.
func (r *Renderer) drawSheet(v *bdf.View, region bdf.Rect, w, h int) *image.RGBA {
	d, m := r.canvas(region, w, h)
	all := &mask{r: d.ctx.target.bounds()}
	d.ctx.target.fill(all, r.background(), 1, bdf.BlendSourceOver, all)
	if v.Gridlines {
		r.gridlines(d, v, region, m)
	}
	tile := float64(v.Tile)
	if tile <= 0 {
		tile = 2048
	}
	tx0, ty0 := int(math.Floor(float64(region.X)/tile)), int(math.Floor(float64(region.Y)/tile))
	tx1 := int(math.Floor((float64(region.X+region.W) - 1e-6) / tile))
	ty1 := int(math.Floor((float64(region.Y+region.H) - 1e-6) / tile))
	for ty := ty0; ty <= ty1; ty++ {
		for tx := tx0; tx <= tx1; tx++ {
			hs, ok := v.Tiles[strconv.Itoa(tx)+","+strconv.Itoa(ty)]
			if !ok {
				continue
			}
			hash, err := bdf.ParseHash(hs)
			if err != nil {
				r.warnf("bad tile hash %q", hs)
				continue
			}
			o := r.object(hash)
			if o == nil {
				continue
			}
			tm := m.translate(float64(tx)*tile, float64(ty)*tile)
			d.drawTop(o, tm, d.clipTo(tm, 0, 0, tile, tile))
		}
	}
	return d.ctx.target.toRGBA()
}

// gridlines draws the lines between the cells of a sheet as the viewer
// does: one device pixel wide, snapped to pixels, in #d9d9d9.
func (r *Renderer) gridlines(d *drawer, v *bdf.View, region bdf.Rect, m matrix) {
	sw, sh := SheetSize(v)
	x1 := math.Min(float64(region.X+region.W), sw)
	y1 := math.Min(float64(region.Y+region.H), sh)
	color := premul(0xd9d9d9ff)
	all := &mask{r: d.ctx.target.bounds()}
	line := func(ax, ay, bx, by float64) {
		d.ctx.target.fill(rectMask(ax, ay, bx, by, all.r), color, 1, bdf.BlendSourceOver, all)
	}
	walk := func(runs []bdf.Run, from, to float64, fn func(float64)) {
		pos := 0.0
		for _, run := range runs {
			for i := 0; i < int(run[0]); i++ {
				pos += float64(run[1])
				if pos > to {
					return
				}
				if pos >= from {
					fn(pos)
				}
			}
		}
	}
	_, top := m.apply(0, float64(region.Y))
	_, bottom := m.apply(0, y1)
	walk(v.Cols, float64(region.X), x1, func(x float64) {
		px, _ := m.apply(x, 0)
		px = math.Round(px)
		line(px, top, px+1, bottom)
	})
	left, _ := m.apply(float64(region.X), 0)
	right, _ := m.apply(x1, 0)
	walk(v.Rows, float64(region.Y), y1, func(y float64) {
		_, py := m.apply(0, y)
		py = math.Round(py)
		line(left, py, right, py+1)
	})
}
