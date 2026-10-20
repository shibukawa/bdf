// Package pdl writes BDF pages in the languages of printers of labels,
// receipts and continuous forms, which take neither PDF nor PostScript:
// ZPL (Zebra), SBPL (SATO) and TPCL (Toshiba TEC) for labels, ESC/POS
// (Epson TM and the printers like it) for receipts, and ESC/P (Epson VP
// and the dot-matrix printers like it) for continuous forms.
//
// A page is drawn by raster/imagebdf in black and white at the printer's
// dots: text and lines cut at the half, images dithered, and the modules
// of QR codes and barcodes on whole dots. The symbols that a BARCODE mark
// tells (docs/spec.md §7.8) are printed by a label printer with its own
// barcode commands instead, which its readers read best.
//
// The printer is described by a Profile (its language, dots, width and
// the longest page it takes) and the job by Options, whose Settings are
// what a label, receipt or form printer is asked besides what IPP asks:
// the darkness and speed, the media's tracking and gap, the cutter and
// the cash drawer.
package pdl

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print/raster"
)

// Languages.
const (
	ZPL    = "zpl"
	SBPL   = "sbpl"
	TPCL   = "tpcl"
	ESCPOS = "escpos"
	ESCP   = "escp"
)

// Languages lists the languages, in the order of the docs.
func Languages() []string { return []string{ZPL, SBPL, TPCL, ESCPOS, ESCP} }

// Language returns the language a name means ("zpl", "ZPL II", "esc/pos"
// …), "" for none.
func Language(name string) string {
	switch strings.NewReplacer("/", "", "-", "", "_", "", " ", "").Replace(strings.ToLower(name)) {
	case "zpl", "zplii", "zpl2":
		return ZPL
	case "sbpl":
		return SBPL
	case "tpcl":
		return TPCL
	case "escpos", "pos":
		return ESCPOS
	case "escp", "escpj", "escp2":
		return ESCP
	}
	return ""
}

// Title is the name of a language as people write it.
func Title(lang string) string {
	switch lang {
	case ZPL:
		return "ZPL"
	case SBPL:
		return "SBPL"
	case TPCL:
		return "TPCL"
	case ESCPOS:
		return "ESC/POS"
	case ESCP:
		return "ESC/P"
	}
	return lang
}

// IsLabel reports whether a language is that of label printers, which
// print a page as a label of a length and print barcodes themselves.
func IsLabel(lang string) bool { return lang == ZPL || lang == SBPL || lang == TPCL }

// Cuts.
const (
	CutNone = "none" // tear off
	CutPage = "page" // after each page (label, receipt)
	CutJob  = "job"  // after the last page
)

// Settings are what a printer of labels, receipts or forms is asked
// besides the paper and the copies: the zero value asks for nothing,
// which is the printer's own settings.
type Settings struct {
	// Darkness is print-darkness: -100 (lightest) to 100 (darkest), 0
	// for the printer's.
	Darkness int
	// Speed is print-speed, in hundredths of a millimetre a second; 0
	// for the printer's.
	Speed int
	// Tracking is media-tracking: "continuous", "gap" (labels with gaps
	// between them, or web), "mark" (a black mark on the back), or ""
	// for the printer's.
	Tracking string
	// Gap is the gap between labels, in points (TPCL asks it with the
	// label's size); 0 for 3 mm.
	Gap float64
	// Cut is CutNone, CutPage or CutJob, "" for the language's default
	// (a receipt is cut after each page, a label not at all).
	Cut string
	// Drawer opens the cash drawer of a receipt printer at the end of the
	// job: 1 its connector pin 2, 2 pin 5, 0 none.
	Drawer int
	// PerforationSkip is the skip over the perforation of a continuous
	// form, in points (ESC/P); 0 for none.
	PerforationSkip float64
	// NoPrinterBarcodes draws the symbols that a label printer would
	// print itself, as the other dots.
	NoPrinterBarcodes bool
	// Mode is ModeRaster, ModeText (ESC/P and ESC/POS: the characters of
	// the print image of a form of a template of text, in the printer's
	// fonts), or "" for the language's default: ESC/P prints a document
	// with a print image in the text mode, the others print dots.
	Mode string
}

// Check validates the settings.
func (s *Settings) Check() error {
	if s.Darkness < -100 || s.Darkness > 100 {
		return fmt.Errorf("pdl: print-darkness %d: -100 to 100", s.Darkness)
	}
	if s.Speed < 0 || s.Speed > 100000 {
		return fmt.Errorf("pdl: print-speed %d", s.Speed)
	}
	switch s.Tracking {
	case "", "continuous", "gap", "web", "mark":
	default:
		return fmt.Errorf("pdl: media-tracking %q: continuous, gap or mark", s.Tracking)
	}
	if s.Gap < 0 || s.Gap > 72*4 || s.PerforationSkip < 0 || s.PerforationSkip > 72*4 {
		return errors.New("pdl: a gap or a skip of more than 4 inches")
	}
	switch s.Cut {
	case "", CutNone, CutPage, CutJob:
	default:
		return fmt.Errorf("pdl: cut %q: none, page or job", s.Cut)
	}
	if s.Drawer < 0 || s.Drawer > 2 {
		return fmt.Errorf("pdl: cash drawer %d: 1 or 2", s.Drawer)
	}
	switch s.Mode {
	case "", ModeRaster, ModeText:
	default:
		return fmt.Errorf("pdl: print mode %q: raster or text", s.Mode)
	}
	return nil
}

// Options are a job for a printer of a language.
type Options struct {
	Profile  Profile
	Settings Settings
	// Pages are the zero-based pages of the view to print, in order; nil
	// prints all.
	Pages []int
	// Copies of the job, 1 when 0; Collate prints them as whole jobs
	// (1 2 3 1 2 3), else each page as many times (1 1 2 2 3 3).
	Copies  int
	Collate bool
	// MediaWidth and MediaLength are the size of the label or the form,
	// in points: the width is that of the profile when 0 (or when it is
	// wider), the length that of each page.
	MediaWidth, MediaLength float64
	// Scaling and Orientation place a page on the label as on a sheet;
	// a page is shrunk to the width when it is wider by default.
	Scaling     raster.Scaling
	Orientation raster.Orientation
	// SkipRoles are the roles of the layers left out.
	SkipRoles []string
	// Title names the job (the printers that keep a name of a job).
	Title string
	// FontFS, FontDirs and NoSystemFonts find the fonts of the text of
	// fonts referred to by name.
	FontFS        fs.FS
	FontDirs      []string
	NoSystemFonts bool
	// BandRows and BandWorkers are those of the renderer
	// (imagebdf.Options).
	BandRows, BandWorkers int
	// MaxPages bounds the pages written (pages times copies the writer
	// makes): 10,000 when 0.
	MaxPages int
}

// Result is what a job wrote.
type Result struct {
	// Pages are the pages (labels, receipts, forms) written, copies
	// included.
	Pages int
	// Symbols are the QR codes and barcodes the printer prints itself.
	Symbols int
	// Text is set when the pages were printed in the text mode.
	Text bool
	// Unsupported are the settings the language has no command for
	// (name=value), which were left out.
	Unsupported []string
	Warnings    []string
}

// Write writes pages of a view in the language of the profile.
func Write(w io.Writer, doc *bdf.Document, view *bdf.View, opts *Options) (*Result, error) {
	o := *opts
	res := &Result{}
	if err := o.Profile.Check(); err != nil {
		return res, err
	}
	if err := o.Settings.Check(); err != nil {
		return res, err
	}
	if o.Copies < 0 || o.Copies > 9999 {
		return res, fmt.Errorf("pdl: %d copies", o.Copies)
	}
	if o.Copies == 0 {
		o.Copies = 1
	}
	pages := o.Pages
	if pages == nil {
		for i := range view.Pages {
			pages = append(pages, i)
		}
	}
	for _, p := range pages {
		if p < 0 || p >= len(view.Pages) || view.Pages[p] == nil {
			return res, fmt.Errorf("pdl: no page %d", p+1)
		}
	}
	if len(pages) == 0 {
		return res, errors.New("pdl: no pages to print")
	}
	limit := o.MaxPages
	if limit <= 0 {
		limit = 10000
	}
	if n := len(pages) * o.Copies; n > limit {
		return res, fmt.Errorf("pdl: %d pages to print, more than %d", n, limit)
	}
	o.Pages = pages
	ew := &errWriter{w: w}
	im, err := textImage(doc, view, &o, res)
	if err != nil {
		return res, err
	}
	if im != nil {
		if done, err := writeText(ew, im, &o, res); done || err != nil {
			return res, err
		}
	}
	view = skipLayers(view, o.SkipRoles)
	rd := newRender(doc, view, &o, o.Profile)
	switch o.Profile.Language {
	case ZPL:
		err = writeZPL(ew, rd, res)
	case SBPL:
		err = writeSBPL(ew, rd, res)
	case TPCL:
		err = writeTPCL(ew, rd, res)
	case ESCPOS:
		err = writeESCPOS(ew, rd, res)
	case ESCP:
		err = writeESCP(ew, rd, res)
	default:
		err = fmt.Errorf("pdl: no language %q", o.Profile.Language)
	}
	if err == nil {
		err = ew.err
	}
	res.Warnings = append(res.Warnings, rd.r.Warnings()...)
	return res, err
}

// unsupported notes a setting the language has no command for.
func (r *Result) unsupported(name string, value any) {
	r.Unsupported = append(r.Unsupported, fmt.Sprintf("%s=%v", name, value))
}

// sequence returns the pages to write in order and how many times the
// printer prints each: the copies of a collated job are the pages again,
// those of an uncollated one the printer's count when it has one.
func (rd *render) sequence(printerCopies bool) (pages []int, each int) {
	return rd.o.sequence(printerCopies)
}

func (o *Options) sequence(printerCopies bool) (pages []int, each int) {
	switch {
	case o.Copies == 1:
		return o.Pages, 1
	case printerCopies && (!o.Collate || len(o.Pages) == 1):
		// (the copies of one page are collated as they are)
		return o.Pages, o.Copies
	case o.Collate:
		for range o.Copies {
			pages = append(pages, o.Pages...)
		}
	default:
		for _, p := range o.Pages {
			for range o.Copies {
				pages = append(pages, p)
			}
		}
	}
	return pages, 1
}

// errWriter keeps the first error of its writes.
type errWriter struct {
	w   io.Writer
	err error
	n   int64
}

func (e *errWriter) Write(b []byte) (int, error) {
	if e.err != nil {
		return 0, e.err
	}
	n, err := e.w.Write(b)
	e.n += int64(n)
	e.err = err
	return n, err
}

func (e *errWriter) str(s string) { io.WriteString(e, s) }

func (e *errWriter) printf(format string, args ...any) { fmt.Fprintf(e, format, args...) }

// skipLayers returns the view with the layers of the roles left out.
func skipLayers(v *bdf.View, roles []string) *bdf.View {
	if len(roles) == 0 {
		return v
	}
	skip := map[string]bool{}
	for _, r := range roles {
		skip[r] = true
	}
	out := *v
	out.Pages = make([]*bdf.Page, len(v.Pages))
	for i, p := range v.Pages {
		if p == nil {
			continue
		}
		q := *p
		q.Layers = nil
		for _, l := range p.Layers {
			if !skip[l.Role] {
				q.Layers = append(q.Layers, l)
			}
		}
		out.Pages[i] = &q
	}
	return &out
}
