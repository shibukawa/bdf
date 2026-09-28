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
//
// A document decides what it asks of a renderer, so the drawing has limits:
// of the objects drawn again (maxReusedInstructions), of the groups, soft
// masks, states and clips that are open at a time (maxLayers, maxStates,
// maxLive), of the outline of a stroke (strokePoints, strokeDashes), of the
// blur of a shadow (maxShadowBlur), of the pictures kept decoded
// (maxPicturePixels), of SVG images (svgDepth, svgNodes, svgBudget,
// svgLayers and the constants beside them) and of the gridlines of a sheet
// (maxGridLines). What passes a limit is left out with a warning, and the
// rest is drawn. A panic of a drawing is the error of Page and Region.
package raster

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io/fs"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/imgconv"
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
// paths, images, fonts) for the next drawing, the images up to
// maxPicturePixels, past which those of the drawings before make room; it
// is not safe for concurrent use.
type Renderer struct {
	doc  *bdf.Document
	opts Options

	objects  map[bdf.Hash]*bdf.ObjectPart
	coll     map[bdf.Hash][]*bdf.Path
	paths    map[pathKey]*path
	pictures map[bdf.Hash]*picture
	svgs     map[bdf.Hash]*svgImage
	luts     map[*bdf.Paint]*gradientLUT
	fonts    *fonts
	warnings []string
	warned   map[string]bool

	// drawing counts the images drawn, to tell the pictures the one being
	// drawn uses from those kept from before (see keep)
	drawing int
	kept    []*picture
	pixels  int // of the pictures kept
	room    int // for the pixels of pictures: maxPicturePixels
	sc      scratch
}

// plain draws without the shortcuts that change no pixel: what lies outside
// the canvas is drawn like the rest, the loops of the common cases and the
// memory kept from one drawing to the next are not used. The tests compare
// both ways of drawing.
var plain bool

// scratch returns the working memory of the rasterizer, emptied for a
// drawing instruction.
func (r *Renderer) scratch() *scratch {
	r.sc.reset()
	return &r.sc
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
		luts:     map[*bdf.Paint]*gradientLUT{},
		warned:   map[string]bool{},
		room:     maxPicturePixels,
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

// maxPicturePixels bounds the pixels of the pictures a renderer keeps: the
// images it decoded and the rasters of SVG images (4 bytes a pixel, and a
// third more for the smaller copies they are scaled down from). One image
// may have imgconv.MaxDecodePixels; a document decides how many it draws.
const maxPicturePixels = 256 << 20

// makeRoom makes room for a picture of n pixels among those kept: pictures
// that the image being drawn has not used are let go, the ones used longest
// ago first. It reports whether there is room.
func (r *Renderer) makeRoom(n int) bool {
	if n > r.room {
		return false
	}
	for r.pixels+n > r.room {
		oldest := -1
		for i, p := range r.kept {
			if p.used < r.drawing && (oldest < 0 || p.used < r.kept[oldest].used) {
				oldest = i
			}
		}
		if oldest < 0 {
			return false
		}
		p := r.kept[oldest]
		p.drop()
		r.pixels -= p.w * p.h
		r.kept = append(r.kept[:oldest], r.kept[oldest+1:]...)
	}
	return true
}

// keep counts a picture among those kept; drop takes it from where it is
// kept when its room is needed.
func (r *Renderer) keep(p *picture, drop func()) {
	p.used, p.drop = r.drawing, drop
	r.kept = append(r.kept, p)
	r.pixels += p.w * p.h
}

// tooManyPixels warns of pictures past maxPicturePixels.
const tooManyPixels = "the images of a page have more than %d pixels in all: the rest of them is not drawn"

// decode decodes a bitmap image to keep, making room for it first: nil
// for one that cannot be decoded, or that there is no room for while this
// image is drawn (full).
func (r *Renderer) decode(data []byte, drop func()) (pic *picture, full bool) {
	if c, _, err := image.DecodeConfig(bytes.NewReader(data)); err == nil && c.Width > 0 && c.Height > 0 &&
		int64(c.Width)*int64(c.Height) <= imgconv.MaxDecodePixels && !r.makeRoom(c.Width*c.Height) {
		r.warnf(tooManyPixels, r.room)
		return nil, true
	}
	p, format, ok := decodePicture(data)
	if !ok {
		switch format {
		case "avif":
			r.warnf("AVIF images are not drawn")
		case "too large":
			r.warnf("an image of more than %d pixels is not drawn", imgconv.MaxDecodePixels)
		default:
			r.warnf("an image (%s) cannot be decoded", format)
		}
		return nil, false
	}
	r.keep(p, drop)
	return p, false
}

// picture returns a decoded bitmap image part, or nil (for an SVG image,
// see svg, and for an image that cannot be decoded).
func (r *Renderer) picture(h bdf.Hash) *picture {
	if p, ok := r.pictures[h]; ok {
		if p != nil {
			p.used = r.drawing
		}
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
		for _, w := range si.doc.warnings {
			r.warnf("%s", w)
		}
		r.svgs[h] = si
		return nil
	}
	p, full := r.decode(part.Data, func() { delete(r.pictures, h) })
	if full {
		delete(r.pictures, h) // there may be room for it in the next image drawn
	}
	if p == nil {
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

// recovered makes an error of a panic of a drawing. The limits of this
// package keep a document from what documents are known to ask for; this
// keeps what nobody thought of from ending the program.
func recovered(img **image.RGBA, err *error) {
	if p := recover(); p != nil {
		*img, *err = nil, fmt.Errorf("raster: the drawing failed: %v", p)
	}
}

// Page draws a whole page of a fixed or flow view (or a strip of a scroll
// view) at scale device pixels per unit.
func (r *Renderer) Page(v *bdf.View, page int, scale float64) (img *image.RGBA, err error) {
	defer recovered(&img, &err)
	if v == nil {
		return nil, errors.New("raster: no view")
	}
	if v.Kind == bdf.ViewSheet {
		return nil, errors.New("raster: a sheet has no pages; draw a Region")
	}
	if page < 0 || page >= len(v.Pages) || v.Pages[page] == nil {
		return nil, ErrNoPage
	}
	p := v.Pages[page]
	if !(p.W > 0) || !(p.H > 0) || !(scale > 0) {
		return nil, errors.New("raster: empty page")
	}
	fw, fh := math.Max(math.Ceil(float64(p.W)*scale-1e-6), 1), math.Max(math.Ceil(float64(p.H)*scale-1e-6), 1)
	if !(fw*fh <= MaxPixels) {
		return nil, ErrTooLarge
	}
	return r.drawPage(p, bdf.Rect{W: p.W, H: p.H}, int(fw), int(fh))
}

// Region draws a rectangle of a view into a w × h image: of the page with
// index page in a fixed or flow view (in page units), of the continuous
// layout of a scroll view (its strips stacked, in units from the top), or
// of a sheet (page is ignored in the last two).
func (r *Renderer) Region(v *bdf.View, page int, region bdf.Rect, w, h int) (img *image.RGBA, err error) {
	defer recovered(&img, &err)
	if v == nil {
		return nil, errors.New("raster: no view")
	}
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
	if page < 0 || page >= len(v.Pages) || v.Pages[page] == nil {
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
	r.drawing++
	d := &drawer{r: r, ctx: &context{target: s}, limit: max(maxLive*16*w*h, minLive), most: maxReusedInstructions}
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
		d.drawTop(l.Obj, m, clip)
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
		pm := m.translate(0, top)
		clip := d.clipTo(pm, 0, 0, float64(b.W), float64(b.H))
		d.ctx.target.fill(clip, bg, 1, bdf.BlendSourceOver, &mask{r: d.ctx.target.bounds()})
		pm = pm.translate(-float64(b.X), -float64(b.Y))
		for _, l := range p.Layers {
			if l.Role != bdf.RoleBody && l.Role != bdf.RoleAnnotation {
				continue
			}
			d.drawTop(l.Obj, pm, clip)
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
	if !(tile >= 1) || math.IsInf(tile, 0) {
		// a region would hold any number of them
		r.warnf("the tiles of a sheet are %g units wide: they are not drawn", tile)
		return d.ctx.target.toRGBA()
	}
	index := func(v float64) int { return int(math.Floor(clampF(v / tile))) }
	tx0, ty0 := index(float64(region.X)), index(float64(region.Y))
	tx1, ty1 := index(float64(region.X+region.W)-1e-6), index(float64(region.Y+region.H)-1e-6)
	type place struct {
		tx, ty int
		hash   string
	}
	var places []place
	if (tx1-tx0+1)*(ty1-ty0+1) > len(v.Tiles) {
		// more places than tiles: the tiles tell theirs
		for k, hs := range v.Tiles {
			x, y, _ := strings.Cut(k, ",")
			tx, errx := strconv.Atoi(x)
			ty, erry := strconv.Atoi(y)
			if errx != nil || erry != nil || tx < tx0 || tx > tx1 || ty < ty0 || ty > ty1 || strconv.Itoa(tx)+","+strconv.Itoa(ty) != k {
				continue
			}
			places = append(places, place{tx, ty, hs})
		}
		slices.SortFunc(places, func(a, b place) int { return cmp.Or(cmp.Compare(a.ty, b.ty), cmp.Compare(a.tx, b.tx)) })
	} else {
		for ty := ty0; ty <= ty1; ty++ {
			for tx := tx0; tx <= tx1; tx++ {
				if hs, ok := v.Tiles[strconv.Itoa(tx)+","+strconv.Itoa(ty)]; ok {
					places = append(places, place{tx, ty, hs})
				}
			}
		}
	}
	for _, p := range places {
		hash, err := bdf.ParseHash(p.hash)
		if err != nil {
			r.warnf("bad tile hash %q", p.hash)
			continue
		}
		tm := m.translate(float64(p.tx)*tile, float64(p.ty)*tile)
		d.drawTop(hash, tm, d.clipTo(tm, 0, 0, tile, tile))
	}
	return d.ctx.target.toRGBA()
}

// Limits of the gridlines of a sheet, across and down: a document decides
// how many rows and columns a sheet has.
const (
	// maxGridSteps bounds the rows or columns gone through for their lines.
	maxGridSteps = 1 << 24
	// maxGridLines bounds the lines drawn. A region shows a line a pixel
	// at most; hidden rows and columns draw theirs again at the same place.
	maxGridLines = 1 << 16
)

// gridlines draws the lines between the cells of a sheet as the viewer
// does: one device pixel wide, snapped to pixels, in #d9d9d9.
func (r *Renderer) gridlines(d *drawer, v *bdf.View, region bdf.Rect, m matrix) {
	sw, sh := SheetSize(v)
	x1 := math.Min(float64(region.X+region.W), sw)
	y1 := math.Min(float64(region.Y+region.H), sh)
	color := premul(0xd9d9d9ff)
	t := d.ctx.target
	all := &mask{r: t.bounds()}
	var before []float32
	// line draws a line; again says that the line before was the same, and
	// the result whether drawing it once more changed a pixel (it changes
	// those the line covers in part, until they have its colour)
	line := func(ax, ay, bx, by float64, again bool) bool {
		cov := rectMask(ax, ay, bx, by, all.r)
		if !again {
			t.fill(cov, color, 1, bdf.BlendSourceOver, all)
			return true
		}
		if cov.empty() {
			return false
		}
		if cov.a == nil {
			// it covers its pixels whole, which have its colour
			if plain {
				t.fill(cov, color, 1, bdf.BlendSourceOver, all)
			}
			return false
		}
		before = before[:0]
		for y := cov.r.Min.Y; y < cov.r.Max.Y; y++ {
			before = append(before, t.pix[4*(y*t.w+cov.r.Min.X):4*(y*t.w+cov.r.Max.X)]...)
		}
		t.fill(cov, color, 1, bdf.BlendSourceOver, all)
		n := 4 * cov.r.Dx()
		for i, y := 0, cov.r.Min.Y; y < cov.r.Max.Y; i, y = i+1, y+1 {
			if !slices.Equal(before[i*n:(i+1)*n], t.pix[4*(y*t.w+cov.r.Min.X):4*(y*t.w+cov.r.Max.X)]) {
				return true
			}
		}
		return false
	}
	// walk goes through the rows or columns from the first; at gives the
	// pixel of the line after one and draw draws it
	walk := func(runs []bdf.Run, from, to float64, at func(float64) float64, draw func(px float64, again bool) bool) {
		pos, steps, lines := 0.0, 0, 0
		last, settled := math.NaN(), false
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
				if !(pos >= from) {
					if size == 0 {
						break // so are the rest of the run
					}
					continue
				}
				px := at(pos)
				again := px == last
				if again && settled {
					if plain {
						draw(px, again)
					} else if size == 0 {
						break
					}
					continue
				}
				if lines++; lines > maxGridLines {
					r.warnf("a sheet has more than %d gridlines in a region: the rest are not drawn", maxGridLines)
					return
				}
				changed := draw(px, again)
				last, settled = px, again && !changed
			}
		}
	}
	_, top := m.apply(0, float64(region.Y))
	_, bottom := m.apply(0, y1)
	walk(v.Cols, float64(region.X), x1, func(x float64) float64 {
		px, _ := m.apply(x, 0)
		return math.Round(px)
	}, func(px float64, again bool) bool {
		return line(px, top, px+1, bottom, again)
	})
	left, _ := m.apply(float64(region.X), 0)
	right, _ := m.apply(x1, 0)
	walk(v.Rows, float64(region.Y), y1, func(y float64) float64 {
		_, py := m.apply(0, y)
		return math.Round(py)
	}, func(py float64, again bool) bool {
		return line(left, py, right, py+1, again)
	})
}
