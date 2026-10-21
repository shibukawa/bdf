// Package imagebdf draws BDF documents into images without a browser: the
// pages of fixed and flow views, stretches of scroll views and regions of
// sheets, for thumbnails and other server-side previews, and objects that
// are in no document (Object).
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
// algorithm; AVIF images and FILTER are not drawn. SVG images are drawn by
// a renderer of the SVG that previews meet, from the tree image/svg reads.
// Warnings list what was left out.
//
// A document decides what it asks of a renderer, so the drawing has limits:
// of the objects drawn again (maxReusedInstructions), of the groups, soft
// masks, states and clips that are open at a time (maxLayers, maxStates,
// maxLive), of the outline of a stroke (strokePoints, strokeDashes), of the
// blur of a shadow (maxShadowBlur), of the pictures kept decoded
// (maxPicturePixels), of SVG images (the limits of reading of image/svg,
// svgBudget, svgLayers and the constants beside them) and of the gridlines
// of a sheet (maxGridLines). What passes a limit is left out with a
// warning, and the rest is drawn. A panic of a drawing is the error of
// Page, Region and Object.
package imagebdf

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
	"github.com/shibukawa/bdf/image/imgconv"
)

// AntialiasMode selects how geometric edges are rasterized. It does not
// change source opacity, image interpolation, shadows or soft masks.
type AntialiasMode uint8

const (
	// AntialiasDefault preserves the existing antialiased rendering at any DPI.
	AntialiasDefault AntialiasMode = iota
	// AntialiasOn explicitly enables antialiased geometric edges.
	AntialiasOn
	// AntialiasOff samples geometric edges at pixel centers, without coverage
	// antialiasing. Subpixel-width strokes still use the existing width-to-alpha
	// compensation, so their pixels need not be opaque.
	AntialiasOff
)

// Mono selects a monochrome image in place of the RGBA one (PageMono,
// RegionMono, PagesMono): for printing in black, where the colours of
// the page are not wanted and the image is a quarter the size or less.
type Mono uint8

const (
	// Gray draws 8 bits of luminance a pixel, as an *image.Gray: the
	// luminance (with the weights of Rec. 601, as image/color converts to
	// gray) of each pixel over white, rounded to the nearest.
	Gray Mono = iota + 1
	// Bilevel draws black and white, as an *image.Paletted of those two
	// colours (black first), which image/png writes at 1 bit a pixel:
	// black where the gray value is below 128. With AntialiasOff the
	// edges are those of the drawing; with antialiasing, the coverage of
	// the edges is cut at the half, not dithered.
	Bilevel
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
	// Antialias controls path, text and geometric clip edges. The zero value
	// retains antialiasing, independent of DPI. Unknown values also retain it
	// and add a renderer warning.
	Antialias AntialiasMode
	// ImageCacheBytes bounds retained completed Page, Region and Renderer.Object
	// images. Zero uses DefaultImageCacheBytes; a negative value disables it.
	ImageCacheBytes int
	// TextCacheBytes bounds retained whole text-run coverage masks. Zero uses
	// DefaultTextCacheBytes; a negative value disables it. Neither cache changes
	// rendering quality. Call Invalidate after changing document parts or FontFS.
	TextCacheBytes int
	// CanvasCacheBytes bounds the working canvases kept from one drawing
	// for the next (16 bytes a pixel of the image drawn, and of each open
	// group), and as much again the coverage of large masks kept from one
	// fill for the next. Zero uses DefaultCanvasCacheBytes; a negative
	// value keeps none. A page at 600 dpi takes a canvas of hundreds of
	// megabytes, which pages drawn one after the other with Page reuse
	// within this budget; Pages keeps one canvas a worker for the length
	// of a batch whatever the budget.
	CanvasCacheBytes int
	// BandRows is the rows of a band of a page drawn in bands (PageRows,
	// PagesRows): the canvas of a band takes 16 bytes a pixel of it, and
	// the instructions of the page are run once a band. Zero uses
	// DefaultBandRows.
	BandRows int
	// BandWorkers is how many bands of a page are drawn at a time
	// (PageRows, PagesRows), each by a worker with a canvas of its own:
	// zero uses the threads of the machine, up to eight, within
	// MaxBandCanvases of canvases; one draws the bands one after the
	// other, with the memory of one.
	BandWorkers int
	// Barcode is called with each QR code and barcode that a BARCODE mark
	// tells (docs/spec.md §7.8) whose box the transform keeps upright or
	// turns by quarters: when it returns true, the drawing of the symbol
	// is left out (a printer of labels prints it with its own commands).
	// It is called again for each band of a page drawn in bands that the
	// symbol may show in, from the workers that draw them at once, and not
	// for an image the image cache had.
	Barcode func(Barcode) bool
	// AlignBarcodes draws those symbols again with their modules on whole
	// pixels, as many pixels wide as fit in their box (and the symbol in
	// its middle): at the few dots an inch of a label or a receipt, bars
	// a fraction of a dot wider or narrower are what a reader cannot
	// read. A symbol whose modules would be narrower than a pixel is drawn
	// as it is.
	AlignBarcodes bool
	// DitherImages draws images in black and white dots (an ordered
	// dither), for an image of black and white (Bilevel): photographs and
	// logos keep their shades, where text and lines, cut at the half, keep
	// their edges.
	DitherImages bool
}

// DefaultBandRows is the rows of a band of a page drawn in bands
// (Options.BandRows): an A4 page at 600 dpi (4960 pixels wide) takes a
// canvas of 10 MiB a band of 128 rows, and is drawn in 55 bands as fast
// as in fewer.
const DefaultBandRows = 128

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
	// free holds the canvases of the drawings before, for the next ones
	// (see newSurface), smallest first, and freeMasks the coverage of
	// masks that drawing is done with (see newMaskBuf); freeBytes and
	// freeMaskBytes are their memory, each within the budget (freeBudget)
	free          []*surface
	freeMasks     [][]float32
	freeBytes     int
	freeMaskBytes int
	images        lru[imageKey, image.Image]
	textMasks     *textMaskCache
	pageWorkers   []*Renderer
	// rowWorkers draw the bands of a page in parallel (see bandWorkers),
	// each into its bandImage, kept for the next page
	rowWorkers []*Renderer
	bandImage  image.Image
	freeLimit  int
	// threads is how many goroutines the passes over a large canvas are
	// split among (see splitRows): those of the machine, divided among
	// the workers of Pages
	threads int
}

// DefaultCanvasCacheBytes bounds the memory of the canvases a renderer
// keeps for its next drawings (Options.CanvasCacheBytes): a page takes 16
// bytes a pixel, and one made for each drawing costs the system as much
// again in giving the memory back and faulting it in.
const DefaultCanvasCacheBytes = 64 << 20

// newSurface returns a cleared canvas of w × h pixels: one of the drawings
// before when it has one large enough, else a new one.
func (r *Renderer) newSurface(w, h int) *surface {
	return r.filledSurface(w, h, solid{})
}

// surfaceFor returns a canvas of w × h pixels with the pixels it was left
// with: one of the drawings before when it has one large enough, else a
// new one (of zeros). The canvas is drawn on when it is begun (see
// surface.begin).
func (r *Renderer) surfaceFor(w, h int) *surface {
	n := 4 * w * h
	if !plain {
		for i, s := range r.free {
			if cap(s.pix) >= n {
				r.free = slices.Delete(r.free, i, i+1)
				r.freeBytes -= 4 * cap(s.pix)
				s.w, s.h, s.pix = w, h, s.pix[:n]
				s.ox, s.oy, s.page = 0, 0, image.Rectangle{}
				s.threads = r.threads
				if poisonSurfaces {
					for i := range s.pix {
						s.pix[i] = float32(math.NaN())
					}
				}
				return s
			}
		}
	}
	s := newSurface(w, h)
	s.threads = r.threads
	return s
}

// poisonSurfaces makes the canvases reused by a renderer hold what no
// drawing leaves in them, for the tests of lazy canvases: a pixel read
// before it is written shows.
var poisonSurfaces bool

// filledSurface returns a canvas of w × h pixels of one colour, as a
// lazy canvas (see surface): no pixel is written. Drawn plainly, the
// canvas is written whole.
func (r *Renderer) filledSurface(w, h int, c solid) *surface {
	s := r.surfaceFor(w, h)
	if plain {
		s.lazy = false
		row := s.rowBuf(0, w)
		repeat(row, c)
		s.setRows(s.bounds(), row)
		return s
	}
	s.begin(c)
	return s
}

// release keeps a canvas that a drawing is done with for the next ones,
// within the budget of the renderer (freeBudget); the smallest make room
// for it.
func (r *Renderer) release(s *surface) {
	if s == nil || plain {
		return
	}
	n := 4 * cap(s.pix)
	if n > r.freeBudget() {
		return
	}
	for r.freeBytes+n > r.freeBudget() && len(r.free) > 0 {
		r.freeBytes -= 4 * cap(r.free[0].pix)
		r.free = slices.Delete(r.free, 0, 1)
	}
	i, _ := slices.BinarySearchFunc(r.free, s, func(a, b *surface) int { return cmp.Compare(cap(a.pix), cap(b.pix)) })
	r.free = slices.Insert(r.free, i, s)
	r.freeBytes += n
}

// trimFree lets go of the canvases and masks kept past a budget, the
// largest first.
func (r *Renderer) trimFree(budget int) {
	for r.freeBytes > budget && len(r.free) > 0 {
		last := len(r.free) - 1
		r.freeBytes -= 4 * cap(r.free[last].pix)
		r.free = r.free[:last]
	}
	for r.freeMaskBytes > budget && len(r.freeMasks) > 0 {
		last := len(r.freeMasks) - 1
		r.freeMaskBytes -= 4 * cap(r.freeMasks[last])
		r.freeMasks = r.freeMasks[:last]
	}
}

// minPooledMask is the least coverage values of a mask whose memory is
// kept for the masks after it: a large mask costs the system more to
// fault in than to clear, a small one is the allocator's to reuse.
const minPooledMask = 1 << 16

// newMaskBuf returns the coverage of a mask of n pixels, all 0: the
// memory of a mask that drawing is done with when there is one large
// enough (see releaseMask), else new.
func (r *Renderer) newMaskBuf(n int) []float32 {
	if r != nil && n >= minPooledMask && !plain {
		for i, a := range r.freeMasks {
			if cap(a) >= n {
				r.freeMasks = slices.Delete(r.freeMasks, i, i+1)
				r.freeMaskBytes -= 4 * cap(a)
				a = a[:n]
				splitRows(r.threads, 0, n, 1, func(y0, y1, _ int) { clear(a[y0:y1]) })
				return a
			}
		}
	}
	return make([]float32, n)
}

// releaseMask keeps the memory of a mask that drawing is done with for
// the masks after it, within the budget of the renderer, and empties the
// mask: what is drawn with it after draws nothing.
func (r *Renderer) releaseMask(m *mask) {
	if r == nil || m == nil || m.a == nil || plain {
		return
	}
	a := m.a
	m.a, m.r, m.full = nil, image.Rectangle{}, image.Rectangle{}
	n := 4 * cap(a)
	if len(a) < minPooledMask || n > r.freeBudget() {
		return
	}
	for r.freeMaskBytes+n > r.freeBudget() && len(r.freeMasks) > 0 {
		r.freeMaskBytes -= 4 * cap(r.freeMasks[0])
		r.freeMasks = slices.Delete(r.freeMasks, 0, 1)
	}
	i, _ := slices.BinarySearchFunc(r.freeMasks, a, func(a, b []float32) int { return cmp.Compare(cap(a), cap(b)) })
	r.freeMasks = slices.Insert(r.freeMasks, i, a)
	r.freeMaskBytes += n
}

// plain draws without the shortcuts that change no pixel: what lies outside
// the canvas is drawn like the rest, the loops of the common cases and the
// memory kept from one drawing to the next are not used. The tests compare
// both ways of drawing.
var plain bool

// scratch returns the working memory of the rasterizer, emptied for a
// drawing instruction.
func (r *Renderer) scratch() *scratch {
	if plain {
		r.sc = scratch{}
	}
	r.sc.reset()
	r.sc.r = r
	r.sc.noAA = r.opts.Antialias == AntialiasOff
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
		threads:  defaultThreads(),
	}
	if opts != nil {
		r.opts = *opts
		r.opts.FontDirs = slices.Clone(opts.FontDirs)
	}
	if r.opts.Antialias > AntialiasOff {
		r.warnf("unknown antialias mode %d: using antialiasing", r.opts.Antialias)
		r.opts.Antialias = AntialiasDefault
	}
	r.fonts = &fonts{Doc: r.doc, FontFS: r.opts.FontFS, FontDirs: r.opts.FontDirs, NoSystemFonts: r.opts.NoSystemFonts, New: newFace, Warnf: r.warnf, NotDone: "not drawn"}
	r.images.limit = cacheLimit(r.opts.ImageCacheBytes, DefaultImageCacheBytes)
	r.textMasks = &textMaskCache{lru: lru[textMaskKey, *mask]{limit: cacheLimit(r.opts.TextCacheBytes, DefaultTextCacheBytes)}}
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
	if _, cut := imgconv.LimitJPEGScans(data); cut {
		r.warnf("a JPEG image has more than %d scans: the rest of them are not read", imgconv.MaxJPEGScans)
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
		for _, w := range si.doc.Warnings {
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
var ErrNoPage = errors.New("imagebdf: no such page")

// MaxPixels is the largest image Page and Region draw (a canvas takes 16
// bytes a pixel while it is drawn).
const MaxPixels = 64 << 20

// ErrTooLarge is returned for an image of more than MaxPixels pixels.
var ErrTooLarge = errors.New("imagebdf: the image would be too large")

// recovered makes an error of a panic of a drawing. The limits of this
// package keep a document from what documents are known to ask for; this
// keeps what nobody thought of from ending the program.
func recovered(img **image.RGBA, err *error) {
	if p := recover(); p != nil {
		*img, *err = nil, fmt.Errorf("imagebdf: the drawing failed: %v", p)
	}
}

// recoveredImage is recovered for an image of any kind.
func recoveredImage(img *image.Image, err *error) {
	if p := recover(); p != nil {
		*img, *err = nil, fmt.Errorf("imagebdf: the drawing failed: %v", p)
	}
}

// Page draws a whole page of a fixed or flow view (or a strip of a scroll
// view) at scale device pixels per unit.
func (r *Renderer) Page(v *bdf.View, page int, scale float64) (*image.RGBA, error) {
	img, err := r.pageMono(v, page, scale, 0)
	if err != nil {
		return nil, err
	}
	return img.(*image.RGBA), nil
}

// PageMono is Page in a monochrome image: an *image.Gray for Gray, an
// *image.Paletted of black and white for Bilevel. It draws the page as
// Page does and converts it; the renderer's antialiasing mode decides the
// edges. An unknown mono is an error.
func (r *Renderer) PageMono(v *bdf.View, page int, scale float64, mono Mono) (image.Image, error) {
	if err := mono.check(); err != nil {
		return nil, err
	}
	return r.pageMono(v, page, scale, mono)
}

// check reports an unknown Mono (0 is the RGBA image).
func (m Mono) check() error {
	if m > Bilevel {
		return fmt.Errorf("imagebdf: unknown monochrome format %d", m)
	}
	return nil
}

func (r *Renderer) pageMono(v *bdf.View, page int, scale float64, mono Mono) (img image.Image, err error) {
	defer recoveredImage(&img, &err)
	p, err := drawablePage(v, page)
	if err != nil {
		return nil, err
	}
	w, h, err := pagePixels(p, scale)
	if err != nil {
		return nil, err
	}
	return r.cachedPage(p, bdf.Rect{W: p.W, H: p.H}, w, h, mono), nil
}

// PagePixels returns the size of the image Page draws of a page of a
// view at scale device pixels per unit, with the errors of Page: the
// units of the page rounded up to pixels, a pixel at least.
func PagePixels(v *bdf.View, page int, scale float64) (w, h int, err error) {
	p, err := drawablePage(v, page)
	if err != nil {
		return 0, 0, err
	}
	return pagePixels(p, scale)
}

// pagePixels is PagePixels for a page (see drawablePage).
func pagePixels(p *bdf.Page, scale float64) (w, h int, err error) {
	if !(scale > 0) {
		return 0, 0, errors.New("imagebdf: empty page")
	}
	fw, fh := math.Max(math.Ceil(float64(p.W)*scale-1e-6), 1), math.Max(math.Ceil(float64(p.H)*scale-1e-6), 1)
	if !(fw*fh <= MaxPixels) {
		return 0, 0, ErrTooLarge
	}
	return int(fw), int(fh), nil
}

// PageDPI draws a whole page at dpi pixels per inch. BDF page units are
// points (1/72 inch), so 72 dpi draws one pixel per page unit.
func (r *Renderer) PageDPI(v *bdf.View, page int, dpi float64) (*image.RGBA, error) {
	return r.Page(v, page, dpi/72)
}

// PageSize draws a whole page into exactly w × h pixels. The page is
// stretched if the requested size has a different aspect ratio.
func (r *Renderer) PageSize(v *bdf.View, page, w, h int) (img *image.RGBA, err error) {
	defer recovered(&img, &err)
	p, err := drawablePage(v, page)
	if err != nil {
		return nil, err
	}
	if w <= 0 || h <= 0 {
		return nil, errors.New("imagebdf: empty image")
	}
	if w > MaxPixels/h {
		return nil, ErrTooLarge
	}
	return r.cachedPage(p, bdf.Rect{W: p.W, H: p.H}, w, h, 0).(*image.RGBA), nil
}

func drawablePage(v *bdf.View, page int) (*bdf.Page, error) {
	if v == nil {
		return nil, errors.New("imagebdf: no view")
	}
	if v.Kind == bdf.ViewSheet {
		return nil, errors.New("imagebdf: a sheet has no pages; draw a Region")
	}
	if page < 0 || page >= len(v.Pages) || v.Pages[page] == nil {
		return nil, ErrNoPage
	}
	p := v.Pages[page]
	if !(p.W > 0) || !(p.H > 0) || math.IsInf(float64(p.W), 0) || math.IsInf(float64(p.H), 0) {
		return nil, errors.New("imagebdf: empty page")
	}
	return p, nil
}

// Region draws a rectangle of a view into a w × h image: of the page with
// index page in a fixed or flow view (in page units), of the continuous
// layout of a scroll view (its strips stacked, in units from the top), or
// of a sheet (page is ignored in the last two).
func (r *Renderer) Region(v *bdf.View, page int, region bdf.Rect, w, h int) (*image.RGBA, error) {
	img, err := r.regionMono(v, page, region, w, h, 0)
	if err != nil {
		return nil, err
	}
	return img.(*image.RGBA), nil
}

// RegionMono is Region in a monochrome image (see PageMono).
func (r *Renderer) RegionMono(v *bdf.View, page int, region bdf.Rect, w, h int, mono Mono) (image.Image, error) {
	if err := mono.check(); err != nil {
		return nil, err
	}
	return r.regionMono(v, page, region, w, h, mono)
}

func (r *Renderer) regionMono(v *bdf.View, page int, region bdf.Rect, w, h int, mono Mono) (img image.Image, err error) {
	defer recoveredImage(&img, &err)
	if v == nil {
		return nil, errors.New("imagebdf: no view")
	}
	if w <= 0 || h <= 0 || !(region.W > 0) || !(region.H > 0) {
		return nil, errors.New("imagebdf: empty region")
	}
	if w > MaxPixels/h {
		return nil, ErrTooLarge
	}
	switch v.Kind {
	case bdf.ViewSheet:
		return r.cachedView(v, region, w, h, true, mono), nil
	case bdf.ViewScroll:
		return r.cachedView(v, region, w, h, false, mono), nil
	}
	if page < 0 || page >= len(v.Pages) || v.Pages[page] == nil {
		return nil, ErrNoPage
	}
	return r.cachedPage(v.Pages[page], region, w, h, mono), nil
}

// canvas starts an image for a region with the device transform that maps
// the region onto it.
func (r *Renderer) canvas(region bdf.Rect, w, h int, bg *solid) (*drawer, matrix) {
	sx, sy := float64(w)/float64(region.W), float64(h)/float64(region.H)
	m := matrix{sx, 0, 0, sy, -float64(region.X) * sx, -float64(region.Y) * sy}
	whole := image.Rect(0, 0, w, h)
	return r.canvasAt(m, whole, whole, bg), m
}

// canvasAt starts a drawing of the pixels rows of an image whole, with
// the device transform m, on a canvas of the background colour when it
// covers the whole image (see filledSurface), else of nothing.
// The rows are the whole image, or a band of it (see pageRows).
func (r *Renderer) canvasAt(m matrix, rows, whole image.Rectangle, bg *solid) *drawer {
	w, h := rows.Dx(), rows.Dy()
	var s *surface
	switch {
	case bg == nil:
		s = r.newSurface(w, h)
	case plain:
		s = r.newSurface(w, h)
		s.fillBackground(&mask{r: s.bounds()}, *bg)
	default:
		s = r.filledSurface(w, h, *bg)
	}
	s.ox, s.oy = rows.Min.X, rows.Min.Y
	if rows != whole {
		s.page = whole
	}
	r.drawing++
	// the groups and clips of a band may take what those of the whole
	// image may: the canvases of a band are its size at most
	d := &drawer{r: r, ctx: &context{target: s}, limit: max(maxLive*16*whole.Dx()*whole.Dy(), minLive), most: maxReusedInstructions, cur: -1}
	d.ctx.st = initialState(m, &mask{r: s.bounds()})
	return d
}

// finish returns the image drawn on the canvas of a drawer, whose canvas
// is kept for the next drawing.
func (d *drawer) finish() *image.RGBA {
	return d.finishMono(0).(*image.RGBA)
}

// finishMono is finish in the image of a Mono (0 for RGBA).
func (d *drawer) finishMono(mono Mono) image.Image {
	s := d.ctx.target
	var img image.Image
	switch mono {
	case Gray:
		img = s.toGray()
	case Bilevel:
		img = s.toBilevel()
	default:
		img = s.toRGBA()
	}
	d.r.release(s)
	return img
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
	ax, ay := m.Apply(x, y)
	bx, by := m.Apply(x+w, y+h)
	return d.r.rectMask(ax, ay, bx, by, d.ctx.target.bounds())
}

func (r *Renderer) drawPage(p *bdf.Page, region bdf.Rect, w, h int, mono Mono) image.Image {
	bg := r.background()
	// the background is the canvas when the page covers the whole image
	sx, sy := float64(w)/float64(region.W), float64(h)/float64(region.H)
	pm := matrix{sx, 0, 0, sy, -float64(region.X) * sx, -float64(region.Y) * sy}
	ax, ay := pm.Apply(0, 0)
	bx, by := pm.Apply(float64(p.W), float64(p.H))
	whole := r.rectMask(ax, ay, bx, by, image.Rect(0, 0, w, h))
	filled, borders := pageBackground(whole, image.Rect(0, 0, w, h), &bg)
	r.releaseMask(whole)
	d, m := r.canvas(region, w, h, filled)
	clip := d.clipTo(m, 0, 0, float64(p.W), float64(p.H))
	switch {
	case borders:
		d.ctx.target.setBorders(clip, bg)
	case filled == nil:
		d.ctx.target.fillBackground(clip, bg)
	}
	for _, l := range p.Layers {
		d.drawTop(l.Obj, m, clip)
	}
	return d.finishMono(mono)
}

// pageBackground tells how the background of a page is drawn on a canvas
// of bounds, from the coverage of the page: as the canvas itself (filled)
// when the page covers the canvas, with its fractional edges written
// after (borders) when it covers the canvas but for the coverage of its
// edges (see surface.setBorders); else drawn on a canvas of nothing.
func pageBackground(cover *mask, bounds image.Rectangle, bg *solid) (filled *solid, borders bool) {
	if cover.r != bounds {
		return nil, false
	}
	if cover.a == nil {
		return bg, false
	}
	if cover.full.Empty() {
		return nil, false
	}
	return bg, true
}

// drawContinuous draws the strips of a scroll view stacked without gaps,
// each cut to its body, as the viewer scrolls them.
func (r *Renderer) drawContinuous(v *bdf.View, region bdf.Rect, w, h int, mono Mono) image.Image {
	d, m := r.canvas(region, w, h, nil)
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
		pm := m.Translate(0, top)
		clip := d.clipTo(pm, 0, 0, float64(b.W), float64(b.H))
		d.ctx.target.fillBackground(clip, bg)
		pm = pm.Translate(-float64(b.X), -float64(b.Y))
		for _, l := range p.Layers {
			if l.Role != bdf.RoleBody && l.Role != bdf.RoleAnnotation {
				continue
			}
			d.drawTop(l.Obj, pm, clip)
		}
	}
	return d.finishMono(mono)
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
func (r *Renderer) drawSheet(v *bdf.View, region bdf.Rect, w, h int, mono Mono) image.Image {
	bg := r.background()
	var filled *solid
	if bg[3] == 1 {
		filled = &bg
	}
	d, m := r.canvas(region, w, h, filled)
	all := &mask{r: d.ctx.target.bounds()}
	if filled == nil {
		d.ctx.target.fillBackground(all, bg)
	}
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
		return d.finishMono(mono)
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
		tm := m.Translate(float64(p.tx)*tile, float64(p.ty)*tile)
		d.drawTop(hash, tm, d.clipTo(tm, 0, 0, tile, tile))
	}
	return d.finishMono(mono)
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
		cov := r.rectMask(ax, ay, bx, by, all.r)
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
			t.touch(y, cov.r.Min.X, cov.r.Max.X, false)
			before = append(before, t.pix[t.at(cov.r.Min.X, y):t.at(cov.r.Max.X, y)]...)
		}
		t.fill(cov, color, 1, bdf.BlendSourceOver, all)
		n := 4 * cov.r.Dx()
		for i, y := 0, cov.r.Min.Y; y < cov.r.Max.Y; i, y = i+1, y+1 {
			if !slices.Equal(before[i*n:(i+1)*n], t.pix[t.at(cov.r.Min.X, y):t.at(cov.r.Max.X, y)]) {
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
	_, top := m.Apply(0, float64(region.Y))
	_, bottom := m.Apply(0, y1)
	walk(v.Cols, float64(region.X), x1, func(x float64) float64 {
		px, _ := m.Apply(x, 0)
		return math.Round(px)
	}, func(px float64, again bool) bool {
		return line(px, top, px+1, bottom, again)
	})
	left, _ := m.Apply(float64(region.X), 0)
	right, _ := m.Apply(x1, 0)
	walk(v.Rows, float64(region.Y), y1, func(y float64) float64 {
		_, py := m.Apply(0, y)
		return math.Round(py)
	}, func(py float64, again bool) bool {
		return line(left, py, right, py+1, again)
	})
}
