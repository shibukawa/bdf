package print

import (
	"fmt"
	"math"
	"strings"

	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/raster"
)

// Document formats.
const (
	FormatPDF       = "application/pdf"
	FormatPWGRaster = "image/pwg-raster"
	FormatURF       = "image/urf"
	FormatPCLm      = "application/PCLm"
)

// chooseFormat picks what to send a printer: the format asked, else PDF
// (the document as vectors, which the printer draws at its resolution)
// when the printer takes it, else PWG Raster.
func chooseFormat(p *ipp.Printer, asked string) (string, error) {
	switch asked {
	case "pdf":
		asked = FormatPDF
	case "pwg-raster", "pwg":
		asked = FormatPWGRaster
	case "urf", "pclm":
		return "", fmt.Errorf("print: %s is not written yet (print-m3)", asked)
	}
	if asked != "" && asked != "auto" {
		if asked != FormatPDF && asked != FormatPWGRaster {
			return "", fmt.Errorf("print: %q: pdf or pwg-raster", asked)
		}
		if len(p.Formats) > 0 && !p.Supports(asked) {
			return "", fmt.Errorf("print: the printer does not take %s (it takes %s)", asked, strings.Join(p.Formats, ", "))
		}
		return asked, nil
	}
	switch {
	case p.Supports(FormatPDF):
		return FormatPDF, nil
	case p.Supports(FormatPWGRaster):
		return FormatPWGRaster, nil
	case p.Supports(FormatURF), p.Supports(FormatPCLm):
		return "", fmt.Errorf("print: the printer takes Apple Raster or PCLm but neither PDF nor PWG Raster; those are written in print-m3")
	case len(p.Formats) == 0:
		// a printer that says nothing of its formats: PDF, which most take
		return FormatPDF, nil
	}
	return "", fmt.Errorf("print: the printer takes neither PDF nor PWG Raster (it takes %s)", strings.Join(p.Formats, ", "))
}

// rasterSetup is how pages are drawn as PWG Raster for a printer.
type rasterSetup struct {
	dpi  int
	typ  raster.Type
	mono bool // drawn in gray or black and white
}

// chooseRaster picks the resolution and the type of the raster: the
// resolution asked when the printer has it, else 600 dpi for black and
// white or 300 for colour or the nearest the printer has; colour when
// the ticket asks for it (or for auto) and the printer has sRGB, gray
// otherwise, black and white when the ticket asks for it and the
// printer has it.
func chooseRaster(p *ipp.Printer, t *Ticket) (rasterSetup, error) {
	var s rasterSetup
	types := p.PWGRasterTypes
	hasType := func(tt raster.Type) bool { return len(types) == 0 || has(types, string(tt)) }
	colour := t.ColorMode == "color" || t.ColorMode == "auto" || t.ColorMode == ""
	if colour && has(p.ColorModes, "color") || colour && len(p.ColorModes) == 0 {
		if hasType(raster.SRGB8) {
			s.typ = raster.SRGB8
		}
	}
	if s.typ == "" {
		s.mono = true
		switch {
		case t.Bilevel && has(types, string(raster.Black1)):
			s.typ = raster.Black1
		case t.Bilevel && has(types, string(raster.SGray1)):
			s.typ = raster.SGray1
		case hasType(raster.SGray8):
			s.typ = raster.SGray8
		case has(types, string(raster.Black1)):
			s.typ = raster.Black1
		default:
			return s, fmt.Errorf("print: the printer's raster types (%s) are not written here", strings.Join(types, ", "))
		}
	}
	want := t.Resolution
	if want == 0 {
		want = 600
		if !s.mono {
			want = 300
		}
	}
	var square []int
	for _, r := range p.PWGRasterResolutions {
		if r.X == r.Y && r.Units == ipp.UnitsDPI {
			square = append(square, int(r.X))
		}
	}
	if len(square) == 0 {
		if t.Resolution != 0 {
			s.dpi = t.Resolution
		} else {
			s.dpi = want
		}
		return s, nil
	}
	best := square[0]
	for _, d := range square {
		if math.Abs(float64(d-want)) < math.Abs(float64(best-want)) || d == want {
			best = d
		}
	}
	if t.Resolution != 0 && best != t.Resolution {
		return s, fmt.Errorf("print: the printer draws PWG Raster at %v dpi, not %d", square, t.Resolution)
	}
	s.dpi = best
	return s, nil
}
