// Package thumbnail draws a thumbnail image of a BDF document: the top of
// the first page of a text document, the whole first slide or drawing, or
// the top-left corner of a spreadsheet.
//
// The layout follows the kind of document (Meta.Source, then the view):
//
//   - Crop (Word, HTML, Markdown, Excel, CSV; PDF and TIFF pages taller
//     than wide): a square from the top-left corner of the first page, as
//     wide as the page, scaled to Size × Size. A scroll view is cut from
//     its top; a sheet from cell A1, the square as large as the smaller
//     side of the sheet but at most MaxSheetSide units, with gridlines.
//   - Fit (PowerPoint, Visio, draw.io, CAD drawings and plots, circuit
//     boards, Illustrator, Photoshop, images, EPUB covers; other PDF and
//     TIFF pages): the whole first page scaled so that its longer side is
//     Size pixels.
//
// Thumbnails show the content of a document outside the document: whoever
// makes one of a password-protected input should keep it as protected as
// the input (the bdf command does not make one unless asked to).
package thumbnail

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"path/filepath"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/imgconv"
	"github.com/shibukawa/bdf/raster"
)

// Mode is how a page becomes a thumbnail.
type Mode int

const (
	// Auto picks Crop or Fit from the kind of document.
	Auto Mode = iota
	// Crop cuts a square from the top-left corner of the first page.
	Crop
	// Fit scales the whole first page into the thumbnail.
	Fit
)

// ParseMode reads a mode name: auto, crop or fit.
func ParseMode(s string) (Mode, error) {
	switch s {
	case "", "auto":
		return Auto, nil
	case "crop":
		return Crop, nil
	case "fit":
		return Fit, nil
	}
	return Auto, errors.New("thumbnail: mode must be auto, crop or fit")
}

func (m Mode) String() string {
	switch m {
	case Crop:
		return "crop"
	case Fit:
		return "fit"
	}
	return "auto"
}

// DefaultSize is the default side of a thumbnail in pixels.
const DefaultSize = 256

// MaxSheetSide is the largest square of a sheet a thumbnail shows, in units.
const MaxSheetSide = 480

// Options configures a thumbnail.
type Options struct {
	// Size is the side of a cropped thumbnail and the longer side of a
	// fitted one, in pixels (default DefaultSize).
	Size int
	// Mode picks the layout (default Auto).
	Mode Mode
	// View is the id of the view to draw (default: the first one).
	View string
	// Raster configures the drawing: fonts and background.
	Raster raster.Options
}

// Result is a thumbnail and how it was made.
type Result struct {
	Image *image.RGBA
	// Mode is the layout used: Crop or Fit.
	Mode Mode
	// Warnings lists what the drawing left out (see raster).
	Warnings []string
}

// cropSources are the formats whose first page is text read from the top.
var cropSources = map[string]bool{"docx": true, "html": true, "markdown": true, "xlsx": true, "csv": true, "parquet": true}

// shapeSources are the formats whose pages may be documents or slides:
// those taller than wide are cropped.
var shapeSources = map[string]bool{"pdf": true, "tiff": true, "": true}

// Make draws the thumbnail of a document.
func Make(doc *bdf.Document, opts *Options) (*Result, error) {
	o := Options{}
	if opts != nil {
		o = *opts
	}
	if o.Size <= 0 {
		o.Size = DefaultSize
	}
	v, err := pickView(doc, o.View)
	if err != nil {
		return nil, err
	}
	mode := o.Mode
	if mode == Auto {
		mode = autoMode(doc.Meta.Source, v)
	}
	r := raster.New(doc, &o.Raster)
	var img *image.RGBA
	switch {
	case v.Kind == bdf.ViewSheet:
		w, h := raster.SheetSize(v)
		// a sheet has no page to fit: both layouts show the corner at A1
		side := math.Min(math.Min(w, h), MaxSheetSide)
		if !(side > 0) {
			side = MaxSheetSide
		}
		img, err = r.Region(v, 0, bdf.Rect{W: float32(side), H: float32(side)}, o.Size, o.Size)
	case v.Kind == bdf.ViewScroll:
		pw, ph := scrollSize(v)
		if mode == Fit {
			if ph > 3*pw {
				ph = 3 * pw // a long page: its top as a tall thumbnail
			}
			w, h := fitSize(pw, ph, o.Size)
			img, err = r.Region(v, 0, bdf.Rect{W: float32(pw), H: float32(ph)}, w, h)
			break
		}
		img, err = r.Region(v, 0, bdf.Rect{W: float32(pw), H: float32(pw)}, o.Size, o.Size)
	default:
		if len(v.Pages) == 0 {
			return nil, errors.New("thumbnail: the view has no pages")
		}
		p := v.Pages[0]
		pw, ph := float64(p.W), float64(p.H)
		if mode == Fit {
			w, h := fitSize(pw, ph, o.Size)
			img, err = r.Region(v, 0, bdf.Rect{W: p.W, H: p.H}, w, h)
			break
		}
		side := math.Min(pw, ph)
		img, err = r.Region(v, 0, bdf.Rect{W: float32(side), H: float32(side)}, o.Size, o.Size)
	}
	if err != nil {
		return nil, err
	}
	return &Result{Image: img, Mode: mode, Warnings: r.Warnings()}, nil
}

func pickView(doc *bdf.Document, id string) (*bdf.View, error) {
	if len(doc.Views) == 0 {
		return nil, errors.New("thumbnail: the document has no views")
	}
	if id == "" {
		return doc.Views[0], nil
	}
	for _, v := range doc.Views {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, errors.New("thumbnail: no view " + id)
}

// autoMode picks the layout for a view of a document converted from source.
func autoMode(source string, v *bdf.View) Mode {
	switch {
	case v.Kind == bdf.ViewSheet || v.Kind == bdf.ViewScroll || cropSources[source]:
		return Crop
	case shapeSources[source] && len(v.Pages) > 0 && v.Pages[0].H > v.Pages[0].W:
		return Crop
	}
	return Fit
}

// scrollSize is the width of a scroll view and the height of its strips.
func scrollSize(v *bdf.View) (w, h float64) {
	for _, p := range v.Pages {
		b := bdf.RectDef{W: p.W, H: p.H}
		if p.Body != nil {
			b = *p.Body
		}
		w = math.Max(w, float64(b.W))
		h += float64(b.H)
	}
	return w, h
}

// fitSize scales w × h so that the longer side is size pixels.
func fitSize(w, h float64, size int) (int, int) {
	if !(w > 0 && h > 0) {
		return size, size
	}
	if w >= h {
		return size, max(1, int(math.Round(float64(size)*h/w)))
	}
	return max(1, int(math.Round(float64(size)*w/h))), size
}

// Formats of Encode.
const (
	PNG  = "png"
	JPEG = "jpeg"
	WebP = "webp"
)

// FormatOf returns the image format a file name asks for by its extension
// (.png, .jpg or .jpeg, .webp), or "".
func FormatOf(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".png":
		return PNG
	case ".jpg", ".jpeg":
		return JPEG
	case ".webp":
		return WebP
	}
	return ""
}

// Encode writes a thumbnail as PNG, JPEG (quality 85) or WebP (lossy,
// quality 80).
func Encode(w io.Writer, img *image.RGBA, format string) error {
	switch format {
	case PNG:
		return (&png.Encoder{CompressionLevel: png.BestCompression}).Encode(w, img)
	case JPEG:
		return jpeg.Encode(w, img, &jpeg.Options{Quality: 85})
	case WebP:
		b, err := imgconv.EncodeWebP(toNRGBA(img), 80, false)
		if err != nil {
			return err
		}
		_, err = w.Write(b)
		return err
	}
	return errors.New("thumbnail: format must be png, jpeg or webp")
}

func toNRGBA(img *image.RGBA) *image.NRGBA {
	b := img.Bounds()
	out := image.NewNRGBA(b)
	for i := 0; i+3 < len(img.Pix); i += 4 {
		a := img.Pix[i+3]
		out.Pix[i+3] = a
		if a == 0 {
			continue
		}
		for c := 0; c < 3; c++ {
			out.Pix[i+c] = uint8(min(255, (int(img.Pix[i+c])*255+int(a)/2)/int(a)))
		}
	}
	return out
}
