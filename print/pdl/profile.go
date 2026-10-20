package pdl

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/image/barcode"
)

// Profile is a model of printer: its language and its dots.
type Profile struct {
	// Name is the profile's name ("zd421-203"), Title the model's.
	Name, Title string
	// Language is the language the printer takes.
	Language string
	// DPI is its dots an inch across the paper; YDPI along it when they
	// differ (0 when they do not). The dots of a page are square: a
	// printer of other dots along is printed at DPI along too.
	DPI, YDPI int
	// Width is the dots it prints across.
	Width int
	// MaxLength is the most dots a page may be long: 0 for no bound.
	MaxLength int
	// Barcodes prints the symbols of BARCODE marks with the printer's
	// own barcode commands (the label languages).
	Barcodes bool
	// Ribbon is set for a printer of thermal transfer, with a ribbon
	// (TPCL asks which).
	Ribbon bool
}

// Check validates a profile.
func (p *Profile) Check() error {
	if slices.Index(Languages(), p.Language) < 0 {
		return fmt.Errorf("pdl: no language %q (%s)", p.Language, strings.Join(Languages(), ", "))
	}
	if p.DPI < 60 || p.DPI > 1200 {
		return fmt.Errorf("pdl: %d dpi", p.DPI)
	}
	if p.Width < 8 || p.Width > 9600 {
		return fmt.Errorf("pdl: %d dots wide", p.Width)
	}
	if p.MaxLength < 0 {
		return fmt.Errorf("pdl: a length of %d dots", p.MaxLength)
	}
	return nil
}

// mmDots is the dots of a length in millimetres at a resolution.
func mmDots(mm float64, dpi int) int { return int(math.Round(mm * float64(dpi) / 25.4)) }

// profiles are the models known by name. The widths are the dots the
// models print across, from their makers' specifications.
var profiles = []Profile{
	{Name: "zd421-203", Title: "Zebra ZD421 203 dpi", Language: ZPL, DPI: 203, Width: 832, MaxLength: 39 * 203, Barcodes: true},
	{Name: "zd421-300", Title: "Zebra ZD421 300 dpi", Language: ZPL, DPI: 300, Width: 1280, MaxLength: 39 * 300, Barcodes: true},
	{Name: "gk420d", Title: "Zebra GK420d", Language: ZPL, DPI: 203, Width: 832, MaxLength: 39 * 203, Barcodes: true},
	{Name: "zt411-203", Title: "Zebra ZT411 203 dpi", Language: ZPL, DPI: 203, Width: 832, MaxLength: 39 * 203, Barcodes: true},
	{Name: "zt411-300", Title: "Zebra ZT411 300 dpi", Language: ZPL, DPI: 300, Width: 1248, MaxLength: 39 * 300, Barcodes: true},
	{Name: "cl4nx-203", Title: "SATO CL4NX Plus 203 dpi", Language: SBPL, DPI: 203, Width: 832, MaxLength: mmDots(2500, 203), Barcodes: true},
	{Name: "cl4nx-305", Title: "SATO CL4NX Plus 305 dpi", Language: SBPL, DPI: 305, Width: 1248, MaxLength: mmDots(1500, 305), Barcodes: true},
	{Name: "b-ex4t1-203", Title: "Toshiba TEC B-EX4T1 203 dpi", Language: TPCL, DPI: 203, Width: 832, MaxLength: mmDots(1498, 203), Barcodes: true, Ribbon: true},
	{Name: "b-ex4t1-305", Title: "Toshiba TEC B-EX4T1 305 dpi", Language: TPCL, DPI: 305, Width: 1248, MaxLength: mmDots(1498, 305), Barcodes: true, Ribbon: true},
	{Name: "tm-t88", Title: "Epson TM-T88VI / VII (80 mm)", Language: ESCPOS, DPI: 180, Width: 512},
	{Name: "tm-m30", Title: "Epson TM-m30III (80 mm)", Language: ESCPOS, DPI: 203, Width: 576},
	{Name: "pos-80", Title: "Receipt printer of 80 mm, 203 dpi", Language: ESCPOS, DPI: 203, Width: 576},
	{Name: "pos-58", Title: "Receipt printer of 58 mm, 203 dpi", Language: ESCPOS, DPI: 203, Width: 384},
	{Name: "vp-f2000", Title: "Epson VP-F2000 (136 columns)", Language: ESCP, DPI: 180, Width: 136 * 18},
	{Name: "vp-d1300", Title: "Epson VP-D1300 (106 columns)", Language: ESCP, DPI: 180, Width: 106 * 18},
	{Name: "escp-80", Title: "Dot-matrix printer of 80 columns", Language: ESCP, DPI: 180, Width: 80 * 18},
}

// Profiles lists the profiles known by name.
func Profiles() []Profile { return slices.Clone(profiles) }

// defaults are the profiles of a language alone: a printer of 4 inches
// for labels, of 80 mm for receipts, of 136 columns for forms.
var defaults = map[string]string{ZPL: "zd421-203", SBPL: "cl4nx-203", TPCL: "b-ex4t1-203", ESCPOS: "pos-80", ESCP: "vp-f2000"}

// LookupProfile returns the profile of a name, or the default one of a
// language named instead.
func LookupProfile(name string) (Profile, bool) {
	n := strings.ToLower(name)
	if l := Language(n); l != "" {
		n = defaults[l]
	}
	for _, p := range profiles {
		if p.Name == n {
			return p, true
		}
	}
	return Profile{}, false
}

// ParseProfile makes a profile of the parameters of a printer's URL:
// profile (a name), lang (a language), dpi, and width (dots, or a length
// in mm or in: "104mm", "4in"), each overriding the one before. With lang
// and no profile, the language's default profile is the start.
func ParseProfile(profile, lang, dpi, width string) (Profile, error) {
	var p Profile
	switch {
	case profile != "":
		var ok bool
		if p, ok = LookupProfile(profile); !ok {
			return p, fmt.Errorf("pdl: no profile %q", profile)
		}
	case lang != "":
		l := Language(lang)
		if l == "" {
			return p, fmt.Errorf("pdl: no language %q (%s)", lang, strings.Join(Languages(), ", "))
		}
		p, _ = LookupProfile(defaults[l])
	default:
		return p, fmt.Errorf("pdl: a profile or a language (%s) is needed", strings.Join(Languages(), ", "))
	}
	if lang != "" {
		l := Language(lang)
		if l == "" {
			return p, fmt.Errorf("pdl: no language %q (%s)", lang, strings.Join(Languages(), ", "))
		}
		if l != p.Language {
			p.Language, p.Barcodes = l, IsLabel(l)
		}
	}
	if dpi != "" {
		n, err := strconv.Atoi(dpi)
		if err != nil {
			return p, fmt.Errorf("pdl: dpi %q", dpi)
		}
		if p.DPI != n {
			// the same width in inches at the new resolution
			p.Width = int(math.Round(float64(p.Width) * float64(n) / float64(p.DPI)))
			if p.MaxLength > 0 {
				p.MaxLength = int(math.Round(float64(p.MaxLength) * float64(n) / float64(p.DPI)))
			}
			p.DPI, p.YDPI = n, 0
		}
	}
	if width != "" {
		n, err := parseDots(width, p.DPI)
		if err != nil {
			return p, err
		}
		p.Width = n
	}
	if profile == "" {
		p.Name, p.Title = "", ""
	}
	return p, p.Check()
}

// parseDots reads a length in dots, or in mm or in at a resolution.
func parseDots(s string, dpi int) (int, error) {
	unit := 0.0
	v := strings.TrimSpace(strings.ToLower(s))
	switch {
	case strings.HasSuffix(v, "mm"):
		unit, v = float64(dpi)/25.4, v[:len(v)-2]
	case strings.HasSuffix(v, "in"):
		unit, v = float64(dpi), v[:len(v)-2]
	case strings.HasSuffix(v, "dots"):
		unit, v = 1, v[:len(v)-4]
	default:
		unit = 1
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil || !(f > 0) || f*unit > 9600 {
		return 0, fmt.Errorf("pdl: a width of %q", s)
	}
	return int(math.Round(f * unit)), nil
}

// printable reports whether a label printer of a language prints a
// symbol with its own commands.
func printable(lang string, s *Symbol) bool {
	switch s.Symbol.Kind {
	case barcode.QR, barcode.EAN13, barcode.EAN8, barcode.UPCA, barcode.Code39, barcode.Code128, barcode.GS1128, barcode.ITF, barcode.ITF14, barcode.NW7:
	default:
		return false
	}
	switch lang {
	case ZPL:
		return true
	case SBPL:
		return sbplPrintable(s)
	case TPCL:
		return tpclPrintable(s)
	}
	return false
}
