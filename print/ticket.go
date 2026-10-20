package print

import (
	"fmt"
	"slices"
	"strings"

	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/raster"
)

// Ticket is what a job asks of a printer, named as IPP names it (its
// JSON keys are the IPP attributes' names). The zero value asks for
// nothing: the printer's defaults.
type Ticket struct {
	// JobName is the job-name (the document's file name, usually).
	JobName string `json:"job-name,omitempty"`
	// Copies of the document (copies); 0 or 1 for one.
	Copies int `json:"copies,omitempty"`
	// Collate is "collated" (1 2 3 1 2 3), "uncollated" (1 1 2 2 3 3)
	// or "" for the printer's default.
	Collate string `json:"collate,omitempty"`
	// Sides is "one-sided", "two-sided-long-edge" or
	// "two-sided-short-edge" (sides), "" for the default.
	Sides string `json:"sides,omitempty"`
	// Media is a PWG media name; "" picks the paper of the size of the
	// first page among the printer's (media-ready first).
	Media string `json:"media,omitempty"`
	// MediaSource ("tray-1", "manual" …) and MediaType ("stationery",
	// "labels" …) go in media-col.
	MediaSource string `json:"media-source,omitempty"`
	MediaType   string `json:"media-type,omitempty"`
	// ColorMode is "color", "monochrome" or "auto" (print-color-mode),
	// "" for the default.
	ColorMode string `json:"print-color-mode,omitempty"`
	// Bilevel draws black and white, without gray, when the colour mode
	// is monochrome and the document is sent as PWG Raster (black_1 or
	// sgray_1): line drawings and text, at a sixth of the size.
	Bilevel bool `json:"bilevel,omitempty"`
	// Quality is print-quality: 3 draft, 4 normal, 5 high, 0 the default.
	Quality int `json:"print-quality,omitempty"`
	// Resolution in dots per inch (printer-resolution, and the raster's);
	// 0 picks one.
	Resolution int `json:"printer-resolution,omitempty"`
	// Pages are the one-based pages to print, in order; nil prints all.
	Pages []int `json:"pages,omitempty"`
	// Scaling is print-scaling: how a page is fitted to the paper.
	Scaling raster.Scaling `json:"print-scaling,omitempty"`
	// Orientation is orientation-requested; 0 turns landscape pages
	// onto portrait paper by itself.
	Orientation raster.Orientation `json:"orientation-requested,omitempty"`
	// Finishings are the enums of finishings (4 staple, 5 punch …).
	Finishings []int `json:"finishings,omitempty"`
	// OutputBin is output-bin.
	OutputBin string `json:"output-bin,omitempty"`
	// HoldUntil is job-hold-until ("indefinite", "night" …).
	HoldUntil string `json:"job-hold-until,omitempty"`
	// Priority is job-priority, 1 to 100; 0 for the default.
	Priority int `json:"job-priority,omitempty"`
}

// Merge fills what a ticket does not ask with what defaults ask: the
// defaults of a queue under the ticket of a job.
func (t Ticket) Merge(defaults Ticket) Ticket {
	set := func(v *string, d string) {
		if *v == "" {
			*v = d
		}
	}
	setInt := func(v *int, d int) {
		if *v == 0 {
			*v = d
		}
	}
	set(&t.JobName, defaults.JobName)
	setInt(&t.Copies, defaults.Copies)
	set(&t.Collate, defaults.Collate)
	set(&t.Sides, defaults.Sides)
	set(&t.Media, defaults.Media)
	set(&t.MediaSource, defaults.MediaSource)
	set(&t.MediaType, defaults.MediaType)
	set(&t.ColorMode, defaults.ColorMode)
	t.Bilevel = t.Bilevel || defaults.Bilevel
	setInt(&t.Quality, defaults.Quality)
	setInt(&t.Resolution, defaults.Resolution)
	if t.Pages == nil {
		t.Pages = slices.Clone(defaults.Pages)
	}
	if t.Scaling == "" {
		t.Scaling = defaults.Scaling
	}
	if t.Orientation == 0 {
		t.Orientation = defaults.Orientation
	}
	if t.Finishings == nil {
		t.Finishings = slices.Clone(defaults.Finishings)
	}
	set(&t.OutputBin, defaults.OutputBin)
	set(&t.HoldUntil, defaults.HoldUntil)
	setInt(&t.Priority, defaults.Priority)
	return t
}

// Check validates the values of a ticket that do not depend on the
// printer.
func (t *Ticket) Check() error { return t.check() }

// copies returns the copies asked for, one at least.
func (t *Ticket) copies() int { return max(t.Copies, 1) }

// duplex tells the sides of the ticket: two-sided, and on the short edge.
func (t *Ticket) duplex() (duplex, tumble bool) {
	switch t.Sides {
	case "two-sided-long-edge":
		return true, false
	case "two-sided-short-edge":
		return true, true
	}
	return false, false
}

// attrs are the job attributes of a ticket for a printer, and the names
// of what the printer does not support (left out of the attributes).
type attrs struct {
	job         []ipp.Attribute
	unsupported []string
}

func (a *attrs) add(name string, ok bool, v ...ipp.Value) {
	if !ok {
		a.unsupported = append(a.unsupported, name+"="+valuesString(v))
		return
	}
	a.job = append(a.job, ipp.Attr(name, v...))
}

func valuesString(v []ipp.Value) string {
	s := make([]string, len(v))
	for i, x := range v {
		s[i] = x.String()
	}
	return strings.Join(s, ",")
}

// sendCopies tells whether the printer makes the copies (its
// copies-supported reaches them); otherwise the document is sent as
// many times.
func sendCopies(p *ipp.Printer, n int) bool {
	return n <= 1 || int(p.Copies.Hi) >= n && int(p.Copies.Lo) <= n
}

// jobAttributes turns a ticket into job attributes for a printer, for a
// document sent in format: the pages selected are the pages of the
// document sent (page-ranges is not asked, as printers differ in it);
// the placement of raster and PostScript pages is in the pages; copies the client makes
// by sending the pages again are not asked either.
func jobAttributes(t *Ticket, p *ipp.Printer, format, media string, printerCopies bool) attrs {
	var a attrs
	isPDF := format == FormatPDF
	// (the printer draws vectors at its resolution; PostScript pages are
	// placed on the paper, as raster pages are)
	vector := isPDF || format == FormatPostScript
	if n := t.copies(); n > 1 && printerCopies {
		a.add("copies", true, ipp.Integer(n))
	}
	if t.Collate != "" && t.copies() > 1 && printerCopies {
		v := "separate-documents-collated-copies"
		if t.Collate == "uncollated" {
			v = "separate-documents-uncollated-copies"
		}
		a.add("multiple-document-handling", has(p.DocumentHandling, v), ipp.Keyword(v))
	}
	if t.Sides != "" {
		a.add("sides", has(p.Sides, t.Sides), ipp.Keyword(t.Sides))
	}
	if media != "" {
		if t.MediaSource != "" || t.MediaType != "" {
			if m, err := ParseMedia(media); err == nil && p.MediaCol {
				members := []ipp.Attribute{ipp.Attr("media-size", ipp.Collection(
					ipp.Attr("x-dimension", ipp.Integer(int(m.Width/mm*100+0.5))),
					ipp.Attr("y-dimension", ipp.Integer(int(m.Height/mm*100+0.5)))))}
				if t.MediaSource != "" {
					ok := has(p.MediaSources, t.MediaSource)
					if !ok {
						a.unsupported = append(a.unsupported, "media-source="+t.MediaSource)
					} else {
						members = append(members, ipp.Attr("media-source", ipp.Keyword(t.MediaSource)))
					}
				}
				if t.MediaType != "" {
					ok := has(p.MediaTypes, t.MediaType)
					if !ok {
						a.unsupported = append(a.unsupported, "media-type="+t.MediaType)
					} else {
						members = append(members, ipp.Attr("media-type", ipp.Keyword(t.MediaType)))
					}
				}
				a.add("media-col", true, ipp.Collection(members...))
			} else {
				a.add("media", has(p.Media, media), ipp.Keyword(media))
				a.unsupported = append(a.unsupported, "media-col")
			}
		} else {
			a.add("media", has(p.Media, media), ipp.Keyword(media))
		}
	}
	if t.ColorMode != "" {
		a.add("print-color-mode", has(p.ColorModes, t.ColorMode), ipp.Keyword(t.ColorMode))
	}
	if t.Quality != 0 {
		a.add("print-quality", slices.Contains(p.Qualities, int32(t.Quality)), ipp.Enum(t.Quality))
	}
	if t.Resolution != 0 && vector {
		r := ipp.Resolution{X: int32(t.Resolution), Y: int32(t.Resolution), Units: ipp.UnitsDPI}
		a.add("printer-resolution", slices.Contains(p.Resolutions, r), ipp.NewResolution(t.Resolution, t.Resolution, ipp.UnitsDPI))
	}
	if isPDF && t.Scaling != "" {
		a.add("print-scaling", has(p.Scaling, string(t.Scaling)), ipp.Keyword(string(t.Scaling)))
	}
	if isPDF && t.Orientation != 0 {
		a.add("orientation-requested", slices.Contains(p.Orientations, int32(t.Orientation)) || len(p.Orientations) == 0, ipp.Enum(int(t.Orientation)))
	}
	if len(t.Finishings) > 0 {
		v := make([]ipp.Value, len(t.Finishings))
		ok := true
		for i, f := range t.Finishings {
			v[i] = ipp.Enum(f)
			ok = ok && slices.Contains(p.Finishings, int32(f))
		}
		a.add("finishings", ok, v...)
	}
	if t.OutputBin != "" {
		a.add("output-bin", has(p.OutputBins, t.OutputBin), ipp.Keyword(t.OutputBin))
	}
	if t.HoldUntil != "" {
		_, ok := p.Has("job-hold-until-supported")
		a.add("job-hold-until", ok, ipp.Keyword(t.HoldUntil))
	}
	if t.Priority != 0 {
		_, ok := p.Has("job-priority-supported")
		a.add("job-priority", ok, ipp.Integer(t.Priority))
	}
	return a
}

func has(list []string, s string) bool { return slices.Contains(list, s) }

// check validates the values of a ticket that do not depend on the
// printer.
func (t *Ticket) check() error {
	switch t.Collate {
	case "", "collated", "uncollated":
	default:
		return fmt.Errorf("print: collate %q: collated or uncollated", t.Collate)
	}
	switch t.Sides {
	case "", "one-sided", "two-sided-long-edge", "two-sided-short-edge":
	default:
		return fmt.Errorf("print: sides %q: one-sided, two-sided-long-edge or two-sided-short-edge", t.Sides)
	}
	switch t.ColorMode {
	case "", "color", "monochrome", "auto", "process-monochrome", "bi-level", "auto-monochrome":
	default:
		return fmt.Errorf("print: print-color-mode %q", t.ColorMode)
	}
	if t.Quality != 0 && (t.Quality < 3 || t.Quality > 5) {
		return fmt.Errorf("print: print-quality %d: 3 (draft), 4 (normal) or 5 (high)", t.Quality)
	}
	if t.Copies < 0 || t.Copies > 9999 {
		return fmt.Errorf("print: %d copies", t.Copies)
	}
	if t.Priority < 0 || t.Priority > 100 {
		return fmt.Errorf("print: job-priority %d: 1 to 100", t.Priority)
	}
	if t.Resolution < 0 || t.Resolution > 4800 {
		return fmt.Errorf("print: %d dpi", t.Resolution)
	}
	return nil
}
