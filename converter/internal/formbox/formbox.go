// Package formbox draws the values of the fields of a form in their places
// on pages: what the fillers of forms whose fields have places share (the
// form of a PDF, a template of pdfme). It has the fonts (chosen for each
// character, measured, and embedded when the pages are done), the breaking
// of text into lines, and the drawing of lines of text, pictures, the
// symbols of codes and marks. Where the text of a field goes in its place
// is for its filler to say: each kind of form has rules of its own.
package formbox

import (
	"bytes"
	"image"
	"io/fs"
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/linebreak"
	"github.com/shibukawa/bdf/image/barcode"
	"github.com/shibukawa/bdf/image/imgconv"
	"github.com/shibukawa/bdf/internal/fontdb"
)

// Options are the fonts that values are drawn with.
type Options struct {
	FontFS        fs.FS
	FontDirs      []string
	NoSystemFonts bool
	// SystemFonts refers to the fonts by their names instead of
	// embedding them.
	SystemFonts  bool
	NoSubset     bool
	NoWOFF2      bool
	IgnoreFSType bool
	Images       imgconv.Options
	Warn         func(msg string)
}

// Renderer draws on the pages of a document.
type Renderer struct {
	Doc   *bdf.Document
	Fonts *fontset.Set
	opts  Options
	db    *fontdb.DB
	cvs   *canvas.Builder
	// images are the pictures added, by the hash of their source data.
	// Hashes avoid retaining another full copy of each encoded image.
	images map[bdf.Hash]picture
}

type picture struct {
	hash bdf.Hash
	w, h float64
	ok   bool
}

// New makes a renderer for a document.
func New(doc *bdf.Document, opts Options) *Renderer {
	warn := opts.Warn
	if warn == nil {
		warn = func(string) {}
		opts.Warn = warn
	}
	db := fontdb.New(opts.FontFS, opts.FontDirs, !opts.NoSystemFonts)
	set := fontset.New(db, warn)
	return &Renderer{Doc: doc, Fonts: set, opts: opts, db: db, cvs: canvas.NewBuilder(doc, set), images: map[bdf.Hash]picture{}}
}

// Layer is what is drawn over a page: an object of the size of the page.
type Layer struct {
	r  *Renderer
	cv *canvas.Canvas
	// Obj is the object of the layer, for what the layer has no drawing
	// of.
	Obj *bdf.Object
}

// Layer starts a layer of a page of a size.
func (r *Renderer) Layer(w, h float64) *Layer {
	cv := r.cvs.New()
	cv.Obj.SetBBox(0, 0, float32(w), float32(h))
	return &Layer{r: r, cv: cv, Obj: cv.Obj}
}

// Drawn reports whether the layer has something drawn on it.
func (l *Layer) Drawn() bool { return l.cv.Drawn }

// Hash is the object of the layer in the document, once the renderer is
// done (Finish).
func (l *Layer) Hash() bdf.Hash { return l.cv.Hash() }

// Finish embeds the fonts and writes the layers into the document. It
// returns the number of fonts embedded.
func (r *Renderer) Finish() int {
	n := 0
	if !r.opts.SystemFonts {
		n = r.Fonts.Embed(r.Doc, fontset.EmbedOptions{NoSubset: r.opts.NoSubset, NoWOFF2: r.opts.NoWOFF2, IgnoreFSType: r.opts.IgnoreFSType})
	}
	r.cvs.Encode()
	r.Fonts.ReportMissing()
	return n
}

// Style is how a line of text is drawn.
type Style struct {
	// Font is the family of the font ("": a sans-serif font).
	Font         string
	Bold, Italic bool
	// Size is the size of the font in points, and Spacing what is added
	// between characters.
	Size, Spacing float64
	Color         bdf.Color
}

func (st *Style) family() string {
	if st.Font == "" {
		return "sans-serif"
	}
	return st.Font
}

// face is the font that draws a character of a style.
func (r *Renderer) face(st *Style, c rune) *fontset.Choice {
	return r.Fonts.FaceFor(st.family(), st.family(), st.Bold, st.Italic, c)
}

// Metrics returns how far the font of a style reaches over and under the
// line of its text, in points.
func (r *Renderer) Metrics(st *Style) (ascent, descent float64) {
	fc := r.Fonts.Choose(st.family(), st.Bold, st.Italic, false)
	asc, desc := fc.Asc, fc.Desc
	if asc <= 0 {
		asc, desc = 0.9, 0.25
	}
	return asc * st.Size, math.Abs(desc) * st.Size
}

// Width measures a line of text.
func (r *Renderer) Width(text string, st *Style) float64 {
	w, n := 0.0, 0
	for _, c := range text {
		w += r.Fonts.Advance(r.face(st, c), c) * st.Size
		n++
	}
	if n > 1 {
		w += st.Spacing * float64(n-1)
	}
	return w
}

// Wrap breaks text into the lines of a width: at its line breaks, and
// where a line is full, after spaces and between the characters that
// lines break between (see linebreak); a word longer than a line is broken
// where the line is full. A width of 0 breaks at line breaks alone.
func (r *Renderer) Wrap(text string, st *Style, width float64) []string {
	text = strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(text)
	var out []string
	for _, para := range strings.Split(text, "\n") {
		if width <= 0 {
			out = append(out, para)
			continue
		}
		rs := []rune(para)
		if len(rs) == 0 {
			out = append(out, "")
			continue
		}
		for start := 0; start < len(rs); {
			w, end, brk := 0.0, start, -1
			for end < len(rs) {
				add := r.Fonts.Advance(r.face(st, rs[end]), rs[end]) * st.Size
				if end > start {
					add += st.Spacing
				}
				// (the spaces at the end of a line take no room)
				if end > start && !linebreak.IsSpace(rs[end]) && w+add > width+1e-6 {
					break
				}
				w += add
				end++
				if end < len(rs) && linebreak.Allowed(rs[end-1], rs[end]) {
					brk = end
				}
			}
			full := end < len(rs)
			if full && brk > start {
				end = brk
			}
			out = append(out, strings.TrimRightFunc(string(rs[start:end]), linebreak.IsSpace))
			start = end
			for full && start < len(rs) && linebreak.IsSpace(rs[start]) {
				start++
			}
		}
	}
	return out
}

// Mark writes a MARK: what the text that follows is of.
func (l *Layer) Mark(kind byte, payload string) { l.Obj.Mark(kind, payload) }

// Text draws a line of text with its left end at x and its line at y, and
// returns its width. The characters are drawn in runs of the fonts that
// have them.
func (l *Layer) Text(text string, st *Style, x, y float64) float64 {
	if text == "" || st.Size <= 0 {
		return 0
	}
	r, o := l.r, l.Obj
	o.FillColor(st.Color)
	o.TextStyle(bdf.AlignLeft, bdf.BaselineAlphabetic, bdf.DirLTR, float32(st.Spacing))
	var run []rune
	var face *fontset.Choice
	runW, total := 0.0, 0.0
	flush := func() {
		if len(run) == 0 {
			return
		}
		o.Font(l.cv.Font(face.Use), float32(st.Size))
		o.FillText(string(run), float32(x+total), float32(y), float32(runW))
		total += runW
		run, runW = run[:0], 0
	}
	for _, c := range text {
		fc := r.face(st, c)
		if fc != face {
			flush()
			face = fc
		}
		// (the spacing is after each character, the last of a run too: a
		// run is as wide as its characters and their spacing)
		runW += r.Fonts.Advance(fc, c)*st.Size + st.Spacing
		run = append(run, c)
	}
	flush()
	l.cv.Drawn = true
	return total - st.Spacing
}

// Rect fills a rectangle.
func (l *Layer) Rect(x, y, w, h float64, c bdf.Color) {
	l.Obj.FillColor(c).FillRect(float32(x), float32(y), float32(w), float32(h))
	l.cv.Drawn = true
}

// Line draws a line of a width.
func (l *Layer) Line(x0, y0, x1, y1, width float64, c bdf.Color) {
	o := l.Obj
	o.StrokeColor(c).Line(float32(width), bdf.CapButt, bdf.JoinMiter, 10).Dash(nil, 0)
	o.StrokePath(o.AddPath((&bdf.Path{}).MoveTo(float32(x0), float32(y0)).LineTo(float32(x1), float32(y1))))
	l.cv.Drawn = true
}

// Transform sets a transform of what is drawn until the Restore of a Save
// around it.
func (l *Layer) Transform(m [6]float64) {
	l.Obj.Transform(float32(m[0]), float32(m[1]), float32(m[2]), float32(m[3]), float32(m[4]), float32(m[5]))
}

// Picture is a picture of the document.
type Picture struct {
	hash bdf.Hash
	// W and H are the size of the picture in pixels.
	W, H float64
}

// Picture adds a picture to the document: a PNG, JPEG, GIF, BMP or WebP
// image, or an SVG image (svg). ok is false for data that is none.
func (r *Renderer) Picture(data []byte) (pic Picture, ok bool) {
	key := bdf.HashOf(data)
	if p, seen := r.images[key]; seen {
		return Picture{p.hash, p.w, p.h}, p.ok
	}
	p := picture{}
	defer func() { r.images[key] = p }()
	format := imgconv.Sniff(data)
	if format == "svg" || format == "" && bytes.Contains(data[:min(len(data), 1024)], []byte("<svg")) {
		p.hash, p.w, p.h, p.ok = r.Doc.AddImage(data), 0, 0, true
		return Picture{p.hash, 0, 0}, true
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || format == "" || cfg.Width <= 0 || cfg.Height <= 0 {
		return Picture{}, false
	}
	if res, err := imgconv.Optimize(data, r.opts.Images); err == nil || res.Data != nil {
		data = res.Data
	}
	p.hash, p.w, p.h, p.ok = r.Doc.AddImage(data), float64(cfg.Width), float64(cfg.Height), true
	return Picture{p.hash, p.w, p.h}, true
}

// Image draws a picture in a rectangle, stretched to it.
func (l *Layer) Image(pic Picture, x, y, w, h float64) {
	l.Obj.Image(l.cv.Image(pic.hash), float32(x), float32(y), float32(w), float32(h))
	l.cv.Drawn = true
}

// Code draws the symbol of a QR code or of a barcode in a rectangle, in a
// color.
func (l *Layer) Code(sym *barcode.Symbol, x, y, w, h float64, c bdf.Color) {
	o := l.Obj
	o.FillColor(c)
	sym.Fill(o, float32(x), float32(y), float32(w), float32(h))
	l.cv.Drawn = true
}

// Check draws a check mark in a rectangle: two strokes, in a color.
func (l *Layer) Check(x, y, w, h float64, c bdf.Color) {
	o := l.Obj
	s := math.Min(w, h)
	x, y = x+(w-s)/2, y+(h-s)/2
	p := (&bdf.Path{}).MoveTo(float32(x+0.18*s), float32(y+0.52*s)).LineTo(float32(x+0.42*s), float32(y+0.76*s)).LineTo(float32(x+0.84*s), float32(y+0.24*s))
	o.StrokeColor(c).Line(float32(math.Max(0.11*s, 0.5)), bdf.CapRound, bdf.JoinRound, 10).Dash(nil, 0)
	o.StrokePath(o.AddPath(p))
	l.cv.Drawn = true
}

// Fill fills a path in a color.
func (l *Layer) Fill(p *bdf.Path, c bdf.Color) {
	o := l.Obj
	o.FillColor(c).FillPath(o.AddPath(p), bdf.NonZero)
	l.cv.Drawn = true
}

// Stroke draws the outline of a path with a line of a width, with round
// ends and corners or with square ones.
func (l *Layer) Stroke(p *bdf.Path, width float64, c bdf.Color, round bool) {
	o := l.Obj
	lineCap, join := byte(bdf.CapButt), byte(bdf.JoinMiter)
	if round {
		lineCap, join = bdf.CapRound, bdf.JoinRound
	}
	o.StrokeColor(c).Line(float32(width), lineCap, join, 10).Dash(nil, 0)
	o.StrokePath(o.AddPath(p))
	l.cv.Drawn = true
}
