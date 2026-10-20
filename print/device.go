package print

import (
	"bytes"
	"fmt"

	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/pdl"
	"github.com/shibukawa/bdf/print/raw"
)

// Formats of the data of a printer's own language, sent to an IPP queue
// that passes it on (a raw queue of CUPS).
const (
	FormatCUPSRaw     = "application/vnd.cups-raw"
	FormatOctetStream = "application/octet-stream"
)

// deviceLanguage returns the language of a printer of labels, receipts or
// forms that a format names ("zpl", "escpos" …), "" for others.
func deviceLanguage(format string) string {
	if _, ok := FormatName(format); ok {
		return ""
	}
	return pdl.Language(format)
}

// IsDeviceFormat reports whether a format (Options.Format) is the
// language of a printer of labels, receipts or forms.
func IsDeviceFormat(format string) bool { return deviceLanguage(format) != "" }

// deviceSettings are the settings of a ticket for a printer of labels,
// receipts or forms.
func (t *Ticket) deviceSettings() pdl.Settings {
	return pdl.Settings{Darkness: t.Darkness, Speed: t.Speed, Tracking: t.MediaTracking,
		Gap: float64(t.LabelGap) * mm / 100, Cut: t.Cut, Drawer: t.CashDrawer,
		PerforationSkip: float64(t.PerforationSkip) * mm / 100, NoPrinterBarcodes: t.NoPrinterBarcodes, Mode: t.PrintMode}
}

// hasDevice reports whether a ticket asks what only a printer of labels,
// receipts or forms is asked.
func (t *Ticket) hasDevice() bool {
	return t.LabelGap != 0 || t.Cut != "" || t.CashDrawer != 0 || t.PerforationSkip != 0 || t.NoPrinterBarcodes || t.PrintMode != ""
}

// deviceDocument writes the selected pages in the language of a profile:
// the copies, the paper and the settings of the ticket in it, with what
// the language does not carry in Result.Unsupported.
func (j *job) deviceDocument(p pdl.Profile) ([]byte, error) {
	t := j.t
	o := &pdl.Options{Profile: p, Settings: t.deviceSettings(), Pages: j.pages, Copies: t.copies(), Collate: t.Collate != "uncollated",
		Scaling: t.Scaling, Orientation: t.Orientation, SkipRoles: j.o.SkipRoles, Title: t.JobName,
		FontFS: j.o.FontFS, FontDirs: j.o.FontDirs, NoSystemFonts: j.o.NoSystemFonts,
		BandRows: j.o.BandRows, BandWorkers: j.o.BandWorkers, MaxPages: j.o.MaxImpressions}
	if t.Media != "" {
		m, err := ParseMedia(t.Media)
		if err != nil {
			return nil, err
		}
		o.MediaWidth, o.MediaLength = m.Width, m.Height
	}
	var unsupported []string
	no := func(name string, v any) { unsupported = append(unsupported, fmt.Sprintf("%s=%v", name, v)) }
	if t.Sides != "" && t.Sides != "one-sided" {
		no("sides", t.Sides)
	}
	if t.ColorMode == "color" {
		no("print-color-mode", t.ColorMode)
	}
	if t.Resolution != 0 && t.Resolution != p.DPI {
		no("printer-resolution", t.Resolution)
	}
	if t.MediaSource != "" {
		no("media-source", t.MediaSource)
	}
	if t.MediaType != "" {
		no("media-type", t.MediaType)
	}
	if t.Quality != 0 {
		no("print-quality", t.Quality)
	}
	for _, f := range t.Finishings {
		no("finishings", f)
	}
	if t.OutputBin != "" {
		no("output-bin", t.OutputBin)
	}
	if t.HoldUntil != "" {
		no("job-hold-until", t.HoldUntil)
	}
	if t.Priority != 0 {
		no("job-priority", t.Priority)
	}
	var buf bytes.Buffer
	res, err := pdl.Write(&buf, j.doc, j.view, o)
	j.res.Sheets = res.Pages
	j.res.Warnings = append(j.res.Warnings, res.Warnings...)
	unsupported = append(unsupported, res.Unsupported...)
	if len(unsupported) > 0 {
		if j.o.Strict {
			return nil, fmt.Errorf("%w: %v", ErrUnsupported, unsupported)
		}
		j.res.Unsupported = append(j.res.Unsupported, unsupported...)
	}
	if err != nil {
		return nil, err
	}
	j.res.Format, j.res.Media = pdl.Title(p.Language), t.Media
	return buf.Bytes(), nil
}

// profile returns the profile a job is written for: the options', else
// the default one of the language.
func (j *job) profile(lang string) (pdl.Profile, error) {
	if j.o.Profile != nil {
		p := *j.o.Profile
		if p.Language != lang {
			return p, fmt.Errorf("print: the profile %s is of %s, not %s", firstOr(p.Name, "given"), pdl.Title(p.Language), pdl.Title(lang))
		}
		return p, nil
	}
	return pdl.ParseProfile("", lang, "", "")
}

// runDevice sends a job in a printer's own language to an IPP queue that
// passes it on as it is: a raw queue of CUPS (application/vnd.cups-raw),
// or one that takes application/octet-stream.
func (j *job) runDevice(lang string) error {
	mime := ""
	for _, f := range []string{FormatCUPSRaw, FormatOctetStream} {
		if j.p.Supports(f) {
			mime = f
			break
		}
	}
	if mime == "" {
		if len(j.p.Formats) > 0 {
			return fmt.Errorf("print: the printer takes no data of its own language (%s or %s): a raw queue of CUPS does", FormatCUPSRaw, FormatOctetStream)
		}
		mime = FormatOctetStream
	}
	p, err := j.profile(lang)
	if err != nil {
		return err
	}
	data, err := j.deviceDocument(p)
	if err != nil {
		return err
	}
	format := j.res.Format
	defer func() { j.res.Format = format }()
	return j.send(mime, j.t.Media, data, nil, "")
}

// rawDevice writes a job for a raw printer of a language: the profile of
// its URL, or of the options.
func (j *job) rawDevice(p *raw.Printer, lang string) ([]byte, error) {
	prof := j.o.Profile
	if prof == nil {
		pp, err := pdl.ParseProfile(p.Profile, lang, p.DPI, p.Width)
		if err != nil {
			return nil, err
		}
		prof = &pp
	}
	if prof.Language != lang {
		return nil, fmt.Errorf("print: the profile %s is of %s, not %s", firstOr(prof.Name, "given"), pdl.Title(prof.Language), pdl.Title(lang))
	}
	return j.deviceDocument(*prof)
}

// deviceJobAttributes are the attributes of IPP that a ticket asks of a
// printer of labels: print-darkness and print-speed (PWG 5100.13), when it
// supports them, and media-tracking in media-col.
func deviceJobAttributes(t *Ticket, p *ipp.Printer, a *attrs) {
	if t.Darkness != 0 {
		_, ok := p.Has("print-darkness-supported")
		a.add("print-darkness", ok, ipp.Integer(t.Darkness))
	}
	if t.Speed != 0 {
		_, ok := p.Has("print-speed-supported")
		a.add("print-speed", ok, ipp.Integer(t.Speed))
	}
	if t.MediaTracking != "" {
		a.unsupported = append(a.unsupported, "media-tracking="+t.MediaTracking)
	}
	if t.hasDevice() {
		for _, s := range []struct {
			name string
			v    any
			set  bool
		}{
			{"bdf-label-gap", t.LabelGap, t.LabelGap != 0}, {"bdf-cut", t.Cut, t.Cut != ""}, {"bdf-cash-drawer", t.CashDrawer, t.CashDrawer != 0},
			{"bdf-perforation-skip", t.PerforationSkip, t.PerforationSkip != 0}, {"bdf-no-printer-barcodes", t.NoPrinterBarcodes, t.NoPrinterBarcodes},
			{"bdf-print-mode", t.PrintMode, t.PrintMode != ""},
		} {
			if s.set {
				a.unsupported = append(a.unsupported, fmt.Sprintf("%s=%v", s.name, s.v))
			}
		}
	}
}
