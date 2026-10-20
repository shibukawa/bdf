package print

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/raw"
)

// PrintRaw prints a document on a printer of 9100 or LPD (package
// print/raw), as PostScript or PDF, as the printer's URL or Options.Format
// says. Such a printer is told of the job only in its file: PostScript
// asks for the paper (the pages placed on it as on an IPP printer), the
// copies, the sides and the media type with setpagedevice; PDF holds its
// copies as pages again, and asks for the sides, the paper and gray in
// its PJL header. What neither carries is left out (or, with Strict,
// refused). The job ends when the printer has the document, which is all
// the protocols tell: Options.Wait is not waited for.
//
// OnCreated is called before the printer may begin to print the job (with
// a job status of no id), OnSent once it has it. An error after the first
// wraps raw.ErrPartial: the printer may have printed a part of the job.
func PrintRaw(ctx context.Context, p *raw.Printer, doc *bdf.Document, opts *Options) (*Result, error) {
	var o Options
	if opts != nil {
		o = *opts
	}
	j := &job{ctx: ctx, doc: doc, o: o, t: &o.Ticket, res: &Result{}, raw: true}
	t := j.t
	if err := t.check(); err != nil {
		return j.res, err
	}
	v, err := PrintView(doc, o.View)
	if err != nil {
		return j.res, err
	}
	j.view = v
	if j.pages = pageNumbers(v, t.Pages); len(j.pages) == 0 {
		return j.res, fmt.Errorf("%w: no page of %d selected", ErrNoPages, len(v.Pages))
	}
	limit := o.MaxImpressions
	if limit <= 0 {
		limit = 10000
	}
	if n := len(j.pages) * t.copies(); n > limit {
		return j.res, fmt.Errorf("print: %d pages to print, more than %d", n, limit)
	}
	format := p.Format
	switch f, ok := FormatName(o.Format); {
	case !ok:
		return j.res, fmt.Errorf("print: %q: pdf or ps", o.Format)
	case f == FormatPDF:
		format = raw.FormatPDF
	case f == FormatPostScript:
		format = raw.FormatPostScript
	case f != "":
		return j.res, fmt.Errorf("print: a 9100 or LPD printer is sent PostScript or PDF, not %s", FormatTitle(f))
	}
	mime := FormatPostScript
	if format == raw.FormatPDF {
		mime = FormatPDF
	}
	// the printer as the rest of the package sees it: it takes the
	// format, and makes the copies a PostScript file asks for
	j.p = &ipp.Printer{Formats: []string{mime}}
	if mime == FormatPostScript {
		j.p.Copies = ipp.Range{Lo: 1, Hi: 9999}
	}
	pjl, unsupported := rawTicket(t, mime, p.UsesPJL(format))
	if len(unsupported) > 0 {
		if o.Strict {
			return j.res, fmt.Errorf("%w: %s", ErrUnsupported, strings.Join(unsupported, ", "))
		}
		j.res.Unsupported = unsupported
	}
	j.res.Format, j.res.Media = mime, t.Media
	data, err := j.document(mime, t.Media)
	if err != nil {
		return j.res, err
	}
	name := firstOr(t.JobName, "document")
	rj := &raw.Job{Name: name, Format: format, Data: data, PJL: pjl, Begin: func() error {
		if o.OnCreated != nil {
			return o.OnCreated(&ipp.JobStatus{Name: name, State: ipp.JobProcessing})
		}
		return nil
	}}
	n, err := p.Send(ctx, rj)
	j.res.Bytes = n
	if err != nil {
		if ctx.Err() != nil && !errors.Is(err, raw.ErrPartial) {
			return j.res, fmt.Errorf("%w: %v", ErrCanceled, ctx.Err())
		}
		return j.res, err
	}
	js := &ipp.JobStatus{Name: name, State: ipp.JobProcessing}
	j.event(Event{Kind: "sent", Job: js, Bytes: n})
	if o.OnSent != nil {
		if err := o.OnSent(js); err != nil {
			return j.res, err
		}
	}
	return j.res, nil
}

// RawCapabilities is what a raw printer is known to support: the format of
// its URL, and no more.
func RawCapabilities(p *raw.Printer) *ipp.Printer {
	format := FormatPostScript
	if p.Format == raw.FormatPDF {
		format = FormatPDF
	}
	return &ipp.Printer{State: ipp.PrinterIdle, Formats: []string{format}}
}

// rawTicket returns the lines of the PJL header of a ticket (for PDF: the
// sides, the paper and gray; a PostScript file asks for them itself), and
// what of the ticket a raw printer is not told.
func rawTicket(t *Ticket, format string, pjl bool) (lines, unsupported []string) {
	no := func(name, value string) { unsupported = append(unsupported, name+"="+value) }
	if format == FormatPDF {
		switch {
		case t.Sides == "":
		case !pjl:
			no("sides", t.Sides)
		case t.Sides == "one-sided":
			lines = append(lines, "@PJL SET DUPLEX=OFF")
		default:
			binding := "LONGEDGE"
			if t.Sides == "two-sided-short-edge" {
				binding = "SHORTEDGE"
			}
			lines = append(lines, "@PJL SET DUPLEX=ON", "@PJL SET BINDING="+binding)
		}
		if t.Media != "" {
			if paper := raw.PJLPaper(t.Media); paper != "" && pjl {
				lines = append(lines, "@PJL SET PAPER="+paper)
			} else {
				no("media", t.Media)
			}
		}
		switch {
		case t.ColorMode != "monochrome":
		case pjl:
			lines = append(lines, "@PJL SET RENDERMODE=GRAYSCALE")
		default:
			no("print-color-mode", t.ColorMode)
		}
		if t.MediaType != "" {
			no("media-type", t.MediaType)
		}
		// the printer fits a page of PDF to its paper as it does
		if t.Scaling != "" {
			no("print-scaling", string(t.Scaling))
		}
		if t.Orientation != 0 {
			no("orientation-requested", fmt.Sprint(int(t.Orientation)))
		}
	}
	if t.MediaSource != "" {
		no("media-source", t.MediaSource)
	}
	if t.Quality != 0 {
		no("print-quality", fmt.Sprint(t.Quality))
	}
	for _, f := range t.Finishings {
		no("finishings", fmt.Sprint(f))
	}
	if t.OutputBin != "" {
		no("output-bin", t.OutputBin)
	}
	if t.HoldUntil != "" {
		no("job-hold-until", t.HoldUntil)
	}
	if t.Priority != 0 {
		no("job-priority", fmt.Sprint(t.Priority))
	}
	return lines, unsupported
}
