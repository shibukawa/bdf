package print

import (
	"fmt"
	"math"
	"strconv"
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

// formats are the formats written, in the order chosen: PDF (the document
// as vectors, which the printer draws at its resolution), then the
// rasters of IPP Everywhere, AirPrint and Mopria.
var formats = []string{FormatPDF, FormatPWGRaster, FormatURF, FormatPCLm}

// formatNames are the short names of the formats.
var formatNames = map[string]string{"pdf": FormatPDF, "pwg-raster": FormatPWGRaster, "pwg": FormatPWGRaster, "urf": FormatURF, "pclm": FormatPCLm}

// chooseFormat picks what to send a printer: the format asked, else the
// first of formats the printer takes.
func chooseFormat(p *ipp.Printer, asked string) (string, error) {
	if f, ok := formatNames[strings.ToLower(asked)]; ok {
		asked = f
	}
	if asked != "" && asked != "auto" {
		known := false
		for _, f := range formats {
			if strings.EqualFold(f, asked) {
				asked, known = f, true
			}
		}
		if !known {
			return "", fmt.Errorf("print: %q: pdf, pwg-raster, urf or pclm", asked)
		}
		if len(p.Formats) > 0 && !p.Supports(asked) {
			return "", fmt.Errorf("print: the printer does not take %s (it takes %s)", asked, strings.Join(p.Formats, ", "))
		}
		return asked, nil
	}
	for _, f := range formats {
		if p.Supports(f) {
			return f, nil
		}
	}
	if len(p.Formats) == 0 {
		// a printer that says nothing of its formats: PDF, which most take
		return FormatPDF, nil
	}
	return "", fmt.Errorf("print: the printer takes none of PDF, PWG Raster, Apple Raster and PCLm (it takes %s)", strings.Join(p.Formats, ", "))
}

// formatTitle names the formats for people.
var formatTitle = map[string]string{FormatPDF: "PDF", FormatPWGRaster: "PWG Raster", FormatURF: "Apple Raster", FormatPCLm: "PCLm"}

// FormatTitle names a format (a MIME type) for people: "PDF", "PWG
// Raster", "Apple Raster", "PCLm"; the type itself for others.
func FormatTitle(format string) string {
	if t, ok := formatTitle[format]; ok {
		return t
	}
	return format
}

// FormatName reads a format as Options.Format takes it: a short name
// ("pdf", "pwg-raster", "urf", "pclm") or a MIME type of a format that is
// written, as its MIME type; "" or "auto" as "". ok is false for others.
func FormatName(s string) (format string, ok bool) {
	if s == "" || s == "auto" {
		return "", true
	}
	if f, ok := formatNames[strings.ToLower(s)]; ok {
		return f, true
	}
	for _, f := range formats {
		if strings.EqualFold(f, s) {
			return f, true
		}
	}
	return "", false
}

// fallbackRaster is the raster sent to a printer that refuses PDF: the
// first it takes of PWG Raster, Apple Raster and PCLm ("" for none).
func fallbackRaster(p *ipp.Printer) string {
	for _, f := range formats[1:] {
		if p.Supports(f) {
			return f
		}
	}
	return ""
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
	if colour(p, t) && hasType(raster.SRGB8) {
		s.typ = raster.SRGB8
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
	var square []int
	for _, r := range p.PWGRasterResolutions {
		if r.X == r.Y && r.Units == ipp.UnitsDPI {
			square = append(square, int(r.X))
		}
	}
	var err error
	s.dpi, err = nearestDPI(square, t.Resolution, s.mono, "PWG Raster")
	return s, err
}

// colour tells whether pages are printed in colour: when the ticket asks
// for it (or for auto) and the printer has it.
func colour(p *ipp.Printer, t *Ticket) bool {
	asked := t.ColorMode == "color" || t.ColorMode == "auto" || t.ColorMode == ""
	return asked && (has(p.ColorModes, "color") || len(p.ColorModes) == 0)
}

// nearestDPI picks among the square resolutions of a printer the one asked
// (an error when it has not it), else the one nearest 600 dpi for black
// and white or 300 for colour; with none, the one asked or the default.
func nearestDPI(square []int, asked int, mono bool, format string) (int, error) {
	want := asked
	if want == 0 {
		want = 600
		if !mono {
			want = 300
		}
	}
	if len(square) == 0 {
		return want, nil
	}
	best := square[0]
	for _, d := range square {
		if math.Abs(float64(d-want)) < math.Abs(float64(best-want)) || d == want {
			best = d
		}
	}
	if asked != 0 && best != asked {
		return 0, fmt.Errorf("print: the printer draws %s at %v dpi, not %d", format, square, asked)
	}
	return best, nil
}

// chooseURF picks how pages are drawn as Apple Raster, from the printer's
// urf-supported: gray (W8) or sRGB (SRGB24), the resolutions (RS300-600)
// and the back of a sheet (DM1 to DM4: normal, flipped, rotated,
// manual-tumble).
func chooseURF(p *ipp.Printer, t *Ticket) (s rasterSetup, back string, err error) {
	var gray, rgb bool
	var square []int
	back = "normal"
	for _, k := range p.URF {
		switch {
		case k == "W8":
			gray = true
		case k == "SRGB24":
			rgb = true
		case strings.HasPrefix(k, "RS"):
			for _, n := range strings.Split(k[2:], "-") {
				if d, err := strconv.Atoi(n); err == nil && d > 0 {
					square = append(square, d)
				}
			}
		case k == "DM2":
			back = "flipped"
		case k == "DM3":
			back = "rotated"
		case k == "DM4":
			back = "manual-tumble"
		}
	}
	switch {
	case colour(p, t) && rgb:
		s.typ = raster.SRGB8
	case gray:
		s.typ, s.mono = raster.SGray8, true
	case rgb:
		// gray pages in sRGB
		s.typ, s.mono = raster.SRGB8, true
	default:
		return s, back, fmt.Errorf("print: the printer's Apple Raster (%s) has neither W8 nor SRGB24", strings.Join(p.URF, ","))
	}
	s.dpi, err = nearestDPI(square, t.Resolution, s.mono, "Apple Raster")
	return s, back, err
}

// choosePCLm picks how pages are drawn as PCLm: gray or sRGB, the
// resolution of pclm-source-resolution-supported, the strips and the
// compression the printer prefers.
func choosePCLm(p *ipp.Printer, t *Ticket) (s rasterSetup, opts raster.PCLmOptions, err error) {
	s.typ, s.mono = raster.SGray8, true
	if colour(p, t) {
		s.typ, s.mono = raster.SRGB8, false
	}
	var square []int
	for _, r := range p.PCLmResolutions {
		if r.X == r.Y && r.Units == ipp.UnitsDPI {
			square = append(square, int(r.X))
		}
	}
	if s.dpi, err = nearestDPI(square, t.Resolution, s.mono, "PCLm"); err != nil {
		return s, opts, err
	}
	opts.StripHeight = int(p.PCLmStripHeight)
	for _, c := range p.PCLmCompression {
		if c == "flate" || c == "rle" || c == "jpeg" {
			opts.Compression = c
			break
		}
	}
	return s, opts, nil
}
