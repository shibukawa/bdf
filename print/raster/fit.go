package raster

import (
	"errors"
	"fmt"
	"math"

	"github.com/shibukawa/bdf/raster/imagebdf"
)

// Scaling is how a page is scaled to a sheet, as print-scaling names it.
type Scaling string

// Scalings.
const (
	// ScaleAutoFit is ScaleFit for a page larger than the sheet and
	// ScaleNone for the others: the page is never enlarged.
	ScaleAutoFit Scaling = "auto-fit"
	// ScaleAuto is ScaleAutoFit (print-scaling's "auto" is fill for
	// borderless documents, which BDF does not tell).
	ScaleAuto Scaling = "auto"
	// ScaleFit scales the page to fit the sheet, centred.
	ScaleFit Scaling = "fit"
	// ScaleFill scales the page to cover the sheet, centred, cutting
	// what passes it.
	ScaleFill Scaling = "fill"
	// ScaleNone keeps the page at its size, centred.
	ScaleNone Scaling = "none"
)

// Orientation is the orientation-requested of a job: how a page is
// turned on the sheet.
type Orientation int

// Orientations.
const (
	// OrientAuto turns a page whose long side is across onto a sheet
	// whose long side is along, as landscape does.
	OrientAuto Orientation = 0
	// Portrait does not turn the page.
	Portrait Orientation = 3
	// Landscape turns the page a quarter turn counterclockwise.
	Landscape Orientation = 4
	// ReverseLandscape turns it a quarter turn clockwise.
	ReverseLandscape Orientation = 5
	// ReversePortrait turns it a half turn.
	ReversePortrait Orientation = 6
)

// Sheet is a sheet of paper at a printer's resolution.
type Sheet struct {
	// Width and Height are its size in points (1/72 inch).
	Width, Height float64
	// DPI is its resolution, across and along.
	DPI int
}

// Pixels returns the size of the sheet in pixels, rounded down as CUPS
// rounds it (an A4 sheet at 300 dpi is 2480 × 3507 pixels).
func (s Sheet) Pixels() (w, h int) {
	scale := float64(s.DPI) / 72
	return int(math.Max(math.Floor(s.Width*scale+1e-6), 1)), int(math.Max(math.Floor(s.Height*scale+1e-6), 1))
}

// Fit returns where a page of w × h points goes on a sheet: scaled as
// scaling says, turned as orientation says, and centred.
func Fit(w, h float64, s Sheet, scaling Scaling, orientation Orientation) (imagebdf.Place, error) {
	if !(w > 0) || !(h > 0) || !(s.Width > 0) || !(s.Height > 0) || s.DPI <= 0 {
		return imagebdf.Place{}, errors.New("raster: an empty page or sheet")
	}
	var turns int
	switch orientation {
	case OrientAuto:
		if (w > h) != (s.Width > s.Height) && w != h && s.Width != s.Height {
			turns = 3
		}
	case Portrait:
	case Landscape:
		turns = 3
	case ReverseLandscape:
		turns = 1
	case ReversePortrait:
		turns = 2
	default:
		return imagebdf.Place{}, fmt.Errorf("raster: orientation %d", orientation)
	}
	tw, th := w, h
	if turns&1 == 1 {
		tw, th = h, w
	}
	f := 1.0
	fit := math.Min(s.Width/tw, s.Height/th)
	switch scaling {
	case ScaleNone:
	case ScaleFit:
		f = fit
	case ScaleFill:
		f = math.Max(s.Width/tw, s.Height/th)
	case ScaleAuto, ScaleAutoFit, "":
		f = math.Min(1, fit)
	default:
		return imagebdf.Place{}, fmt.Errorf("raster: print-scaling %q", scaling)
	}
	if math.Abs(f-1) < 1e-9 {
		f = 1
	}
	dpi := float64(s.DPI) / 72
	pw, ph := s.Pixels()
	scale := dpi * f
	// centred on whole pixels, so that the lines of a page at its size
	// fall on the pixels as they do on its own image
	x := math.Round((float64(pw) - tw*scale) / 2)
	y := math.Round((float64(ph) - th*scale) / 2)
	return imagebdf.Place{W: pw, H: ph, Scale: scale, Turns: turns, X: x, Y: y}, nil
}

// Back returns the transform of the back of a sheet printed two-sided,
// as a printer's pwg-raster-document-sheet-back asks ("normal",
// "flipped", "rotated", "manual-tumble"), for the long edge or the short
// (tumble): crossFeed and feed are -1 where the image is mirrored across
// and along the feed (the header's CrossFeedTransform and
// FeedTransform), as CUPS makes them.
func Back(sheetBack string, tumble bool) (crossFeed, feed int, err error) {
	crossFeed, feed = 1, 1
	switch sheetBack {
	case "normal", "":
	case "flipped":
		if tumble {
			crossFeed = -1
		} else {
			feed = -1
		}
	case "manual-tumble":
		if tumble {
			crossFeed, feed = -1, -1
		}
	case "rotated":
		if !tumble {
			crossFeed, feed = -1, -1
		}
	default:
		return 1, 1, fmt.Errorf("raster: pwg-raster-document-sheet-back %q", sheetBack)
	}
	return crossFeed, feed, nil
}

// Mirror sets the mirrors of a placement for the transforms of Back.
func Mirror(pl imagebdf.Place, crossFeed, feed int) imagebdf.Place {
	pl.MirrorX = pl.MirrorX != (crossFeed < 0)
	pl.MirrorY = pl.MirrorY != (feed < 0)
	return pl
}
