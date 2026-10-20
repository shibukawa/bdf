package print

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"io/fs"
	"math"
	"strings"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/pdl"
	"github.com/shibukawa/bdf/print/raster"
	"github.com/shibukawa/bdf/raster/imagebdf"
	"github.com/shibukawa/bdf/vector/pdfbdf"
)

// Errors of a job, which bdf print tells apart by its exit status.
var (
	// ErrPrinterStopped is returned when the printer is stopped or not
	// accepting jobs: its state reasons tell why (media-empty …).
	ErrPrinterStopped = errors.New("print: the printer is stopped")
	// ErrUnsupported is returned, with Options.Strict, when the printer
	// does not support what the ticket asks.
	ErrUnsupported = errors.New("print: the printer does not support what the ticket asks")
	// ErrAborted is returned when the printer aborted the job.
	ErrAborted = errors.New("print: the printer aborted the job")
	// ErrCanceled is returned when the job was canceled, at the printer
	// or because the context was.
	ErrCanceled = errors.New("print: the job was canceled")
	// ErrTimeout is returned when the job is not done in Options.Wait.
	ErrTimeout = errors.New("print: the job is not done yet")
	// ErrJobLost is returned when the printer lost the job before it had
	// the document: it was restarted after Create-Job, and its job of the
	// id is gone or another's. Nothing was printed; the job can be sent
	// again.
	ErrJobLost = errors.New("print: the printer lost the job before it had the document")
)

// Options are a job: its ticket, how the document is sent, and how its
// end is waited for.
type Options struct {
	Ticket Ticket
	// Format is "pdf", "ps" (PostScript), "pwg-raster", "urf" (Apple
	// Raster), "pclm", or "" (or "auto") to choose: PDF when the printer
	// takes it, else PostScript, else PWG Raster, else Apple Raster, else
	// PCLm. It may be the language of a printer of labels, receipts or
	// forms too ("zpl", "sbpl", "tpcl", "escpos", "escp"; print/pdl),
	// which an IPP queue that passes data on as it is (a raw queue of
	// CUPS) is sent.
	Format string
	// Profile is the printer a language of print/pdl is written for: nil
	// for the profile of a raw printer's URL, or the language's default.
	Profile *pdl.Profile
	// RasterIfLossy sends PWG Raster instead of PDF when the PDF file
	// would leave out what the document draws (FILTER, composite
	// operations of Canvas, shadows past their limits …) and the printer
	// takes PWG Raster: the page is then printed as the viewer shows it.
	RasterIfLossy bool
	// View is the id of the view to print; "" prints the first view of
	// pages.
	View string
	// SkipRoles are the roles of the layers left out (bdf.RoleBackground,
	// bdf.RoleMaster): for paper that has them printed already.
	SkipRoles []string
	// Strict makes what the printer does not support an error; without
	// it, it is left out with a warning.
	Strict bool
	// FontFS, FontDirs and NoSystemFonts find the fonts of the text of
	// fonts referred to by name, for PDF files and raster pages alike.
	FontFS        fs.FS
	FontDirs      []string
	NoSystemFonts bool
	// Wait is how long to wait for the job to be done: 0 returns once
	// the document is sent; a negative value waits for as long as the
	// context lets it.
	Wait time.Duration
	// MaxImpressions bounds the pages sent (pages times copies the
	// client makes): 10,000 when 0.
	MaxImpressions int
	// Progress is told what happens: the job created, its document
	// sent, and each state the printer reports.
	Progress func(Event)
	// BandRows and BandWorkers are those of the renderer of raster pages
	// (imagebdf.Options).
	BandRows, BandWorkers int
	// Printer is what the printer supports, when the caller asked it
	// already: Get-Printer-Attributes is then not sent again. Whether the
	// printer is stopped is still checked, from its state in Printer.
	Printer *ipp.Printer
	// OnCreated is called with the job the printer made, before the
	// document is sent (after it when the printer has no Create-Job and
	// takes Print-Job). A spool records the job's id there, to know after
	// a crash whether the printer has the job already. An error cancels
	// the job at the printer and is returned.
	OnCreated func(*ipp.JobStatus) error
	// OnSent is called once the printer has the whole document. An
	// error is returned; the job is left to the printer.
	OnSent func(*ipp.JobStatus) error
}

// Event is something that happened to a job.
type Event struct {
	// Kind is "created", "sent" or "state".
	Kind string
	Job  *ipp.JobStatus
	// Bytes is the size of the document, once sent.
	Bytes int64
}

// Result is what happened to a job.
type Result struct {
	Printer *ipp.Printer
	// Format is the format the document was sent in, and Media the
	// paper asked for ("" for the printer's default).
	Format string
	Media  string
	// Sheets are the pages sent (as images or as pages of the PDF
	// file), copies the client made included.
	Sheets int
	// Bytes is the size of the document sent.
	Bytes int64
	// Job is the last state the printer reported of the job; nil when
	// it was not created.
	Job *ipp.JobStatus
	// Unsupported are the attributes of the ticket that the printer
	// does not support and that were left out (name=value).
	Unsupported []string
	Warnings    []string
}

func (r *Result) warn(format string, args ...any) {
	r.Warnings = append(r.Warnings, fmt.Sprintf(format, args...))
}

// job is a job being made.
type job struct {
	ctx   context.Context
	c     *ipp.Client
	doc   *bdf.Document
	o     Options
	t     *Ticket
	p     *ipp.Printer
	view  *bdf.View
	pages []int // zero-based pages of the view, in order
	res   *Result
	// raw is set for a printer of 9100 or LPD, which is told of the job
	// only in its file
	raw bool
}

// Print prints a document on the printer of a client: it asks the
// printer what it supports, makes the job's attributes of the ticket,
// sends the document as PDF or PWG Raster, and waits for the job to be
// done as Options.Wait says. When the context is canceled after the job
// is created, the job is canceled at the printer.
func Print(ctx context.Context, c *ipp.Client, doc *bdf.Document, opts *Options) (*Result, error) {
	var o Options
	if opts != nil {
		o = *opts
	}
	j := &job{ctx: ctx, c: c, doc: doc, o: o, t: &o.Ticket, res: &Result{}}
	if err := j.t.check(); err != nil {
		return j.res, err
	}
	v, err := PrintView(doc, o.View)
	if err != nil {
		return j.res, err
	}
	j.view = v
	if j.pages = pageNumbers(v, j.t.Pages); len(j.pages) == 0 {
		return j.res, fmt.Errorf("%w: no page of %d selected", ErrNoPages, len(v.Pages))
	}
	limit := o.MaxImpressions
	if limit <= 0 {
		limit = 10000
	}
	if n := len(j.pages) * j.t.copies(); n > limit {
		return j.res, fmt.Errorf("print: %d pages to print, more than %d", n, limit)
	}
	if j.p = o.Printer; j.p == nil {
		if j.p, err = c.GetPrinterAttributes(ctx); err != nil {
			return j.res, err
		}
	}
	j.res.Printer = j.p
	if err := j.ready(); err != nil {
		return j.res, err
	}
	return j.res, j.run()
}

// ready refuses a printer that is stopped or not accepting jobs.
func (j *job) ready() error {
	p := j.p
	reasons := strings.Join(nonNone(p.Reasons), ", ")
	if p.State == ipp.PrinterStopped {
		if reasons == "" {
			reasons = p.StateMessage
		}
		return fmt.Errorf("%w: %s", ErrPrinterStopped, reasons)
	}
	if _, said := p.Has("printer-is-accepting-jobs"); said && !p.AcceptingJobs {
		return fmt.Errorf("%w: it is not accepting jobs %s", ErrPrinterStopped, reasons)
	}
	return nil
}

func nonNone(reasons []string) []string {
	var out []string
	for _, r := range reasons {
		if r != "none" && r != "" {
			out = append(out, r)
		}
	}
	return out
}

func (j *job) run() error {
	if lang := deviceLanguage(j.o.Format); lang != "" {
		return j.runDevice(lang)
	}
	format, err := chooseFormat(j.p, j.o.Format)
	if err != nil {
		return err
	}
	first := j.view.Pages[j.pages[0]]
	media := j.t.Media
	if media == "" {
		media = ChooseMedia(float64(first.W), float64(first.H), j.p.MediaReady, j.p.Media)
		if media == "" {
			j.res.warn("the printer has no paper of %.0f × %.0f mm: the page is fitted to %s", float64(first.W)/mm, float64(first.H)/mm, firstOr(j.p.MediaDefault, "its default paper"))
		}
	}
	j.res.Media = media
	// the file of a vector format; raster pages are written as they are sent
	data, err := j.document(format, media)
	if err != nil {
		return err
	}
	if r := fallbackRaster(j.p); format == FormatPDF && j.o.RasterIfLossy && r != "" && lossy(j.res.Warnings) {
		j.res.warn("the PDF file would leave out what the page draws: sent as %s", formatTitle[r])
		format, data = r, nil
	}
	fallback := fallbackFormat(j.p)
	refused := func() error {
		j.res.warn("the printer refuses PDF: sent as %s", formatTitle[fallback])
		format = fallback
		d, err := j.document(format, media)
		data = d
		return err
	}
	attrs, err := j.attributes(format, media)
	if err != nil {
		return err
	}
	// Validate-Job tells what is wrong before a job is made; a printer
	// that has no Validate-Job says nothing
	err = j.c.ValidateJob(j.ctx, attrs, format)
	var se *ipp.StatusError
	switch {
	case err == nil, ipp.IsStatus(err, ipp.StatusServerErrorOperationNotSupported):
	case errors.As(err, &se) && se.Status == ipp.StatusClientErrorDocumentFormatNotSupported && format == FormatPDF && fallback != "":
		if err := refused(); err != nil {
			return err
		}
		if attrs, err = j.attributes(format, media); err != nil {
			return err
		}
	case errors.As(err, &se) && len(se.Unsupported) > 0:
		if j.o.Strict {
			return fmt.Errorf("%w: %v", ErrUnsupported, err)
		}
		attrs = without(attrs, se.Unsupported, j.res)
	default:
		return err
	}
	return j.send(format, media, data, attrs, fallback)
}

// send sends a job of a document in a format: the file of a vector
// format, or (data nil) raster pages drawn as they are sent. A printer
// that refuses PDF is sent the fallback format instead.
func (j *job) send(format, media string, data []byte, attrs []ipp.Attribute, fallback string) error {
	var err error
	refused := func() error {
		j.res.warn("the printer refuses PDF: sent as %s", formatTitle[fallback])
		format = fallback
		d, err := j.document(format, media)
		data = d
		return err
	}
	j.res.Format = format
	name := firstOr(j.t.JobName, "document")
	body := func() (io.Reader, func() error) {
		if data != nil {
			return bytes.NewReader(data), func() error { return nil }
		}
		pr, pw := io.Pipe()
		done := make(chan error, 1)
		go func() {
			err := j.raster(pw, media, format)
			pw.CloseWithError(err)
			done <- err
		}()
		return pr, func() error { return <-done }
	}
	counted := func(r io.Reader) *countingReader { return &countingReader{r: r} }
	var js *ipp.JobStatus
	if len(j.p.Operations) == 0 || j.p.SupportsOp(ipp.OpCreateJob) && j.p.SupportsOp(ipp.OpSendDocument) {
		if js, err = j.c.CreateJob(j.ctx, name, attrs); err != nil {
			return err
		}
		j.res.Job = js
		j.event(Event{Kind: "created", Job: js})
		if j.o.OnCreated != nil {
			if err := j.o.OnCreated(js); err != nil {
				j.cancel(js.ID)
				return err
			}
		}
		r, wait := body()
		cr := counted(r)
		sent, err := j.c.SendDocument(j.ctx, js.ID, name, format, cr, true)
		werr := wait()
		if err != nil {
			if ipp.IsStatus(err, ipp.StatusClientErrorDocumentFormatNotSupported) && format == FormatPDF && fallback != "" {
				if err := refused(); err != nil {
					j.cancel(js.ID)
					return err
				}
				j.res.Format = format
				r, wait = body()
				cr = counted(r)
				sent, err = j.c.SendDocument(j.ctx, js.ID, name, format, cr, true)
				werr = wait()
			}
			if err != nil {
				if j.lost(js.ID, name, err) {
					// (the id may be another's job now: it is not canceled)
					return fmt.Errorf("%w: %v", ErrJobLost, err)
				}
				j.cancel(js.ID)
				if werr != nil && !errors.Is(werr, io.ErrClosedPipe) {
					return werr
				}
				return j.interrupted(err)
			}
		}
		if werr != nil {
			return werr
		}
		js = merge(js, sent)
		j.res.Bytes = cr.n
	} else {
		r, wait := body()
		cr := counted(r)
		js, err = j.c.PrintJob(j.ctx, name, format, attrs, cr)
		werr := wait()
		if err != nil {
			if werr != nil && !errors.Is(werr, io.ErrClosedPipe) {
				return werr
			}
			return err
		}
		if werr != nil {
			return werr
		}
		j.res.Bytes = cr.n
		if j.o.OnCreated != nil {
			if err := j.o.OnCreated(js); err != nil {
				j.cancel(js.ID)
				return err
			}
		}
	}
	j.res.Job = js
	j.event(Event{Kind: "sent", Job: js, Bytes: j.res.Bytes})
	if j.o.OnSent != nil {
		if err := j.o.OnSent(js); err != nil {
			return err
		}
	}
	if j.o.Wait == 0 {
		return nil
	}
	return j.wait(js)
}

// merge fills a job's status from a later answer that may say less.
func merge(old, newer *ipp.JobStatus) *ipp.JobStatus {
	if newer == nil {
		return old
	}
	if newer.ID == 0 {
		newer.ID = old.ID
	}
	if newer.URI == "" {
		newer.URI = old.URI
	}
	if newer.Name == "" {
		newer.Name = old.Name
	}
	return newer
}

func firstOr(s, or string) string {
	if s != "" {
		return s
	}
	return or
}

func (j *job) event(e Event) {
	if j.o.Progress != nil {
		j.o.Progress(e)
	}
}

// attributes makes the job attributes of the ticket, leaving out (or,
// when strict, refusing) what the printer does not support.
func (j *job) attributes(format, media string) ([]ipp.Attribute, error) {
	a := jobAttributes(j.t, j.p, format, media, j.printerCopies())
	if len(a.unsupported) > 0 {
		if j.o.Strict {
			return nil, fmt.Errorf("%w: %s", ErrUnsupported, strings.Join(a.unsupported, ", "))
		}
		j.res.Unsupported = append(j.res.Unsupported, a.unsupported...)
	}
	return a.job, nil
}

// printerCopies tells whether the printer makes the copies; otherwise
// the document holds them.
func (j *job) printerCopies() bool { return sendCopies(j.p, j.t.copies()) }

// without leaves out of the attributes the ones the printer named
// unsupported.
func without(attrs []ipp.Attribute, unsupported []ipp.Attribute, res *Result) []ipp.Attribute {
	drop := map[string]bool{}
	for _, u := range unsupported {
		drop[u.Name] = true
		res.Unsupported = append(res.Unsupported, u.Name+"="+u.String())
	}
	var out []ipp.Attribute
	for _, a := range attrs {
		if !drop[a.Name] {
			out = append(out, a)
		}
	}
	return out
}

// order returns the pages of the document sent, the copies the client
// makes included: the pages again for each copy (collated), or each
// page as many times (uncollated); -1 is a blank page, which keeps the
// copies of a two-sided document of an odd number of pages from sharing
// a sheet.
func (j *job) order() []int {
	if j.printerCopies() {
		return j.pages
	}
	n := j.t.copies()
	duplex, _ := j.t.duplex()
	var out []int
	if j.t.Collate == "uncollated" && !duplex {
		for _, p := range j.pages {
			for range n {
				out = append(out, p)
			}
		}
		return out
	}
	for range n {
		out = append(out, j.pages...)
		if duplex && len(j.pages)%2 == 1 {
			out = append(out, -1)
		}
	}
	if duplex && len(out) > 0 && out[len(out)-1] == -1 {
		out = out[:len(out)-1]
	}
	return out
}

// document writes the file of a vector format (PDF, PostScript); nil for
// a raster format, whose pages are written as they are sent.
func (j *job) document(format, media string) ([]byte, error) {
	switch format {
	case FormatPDF:
		return j.pdf()
	case FormatPostScript:
		return j.ps(media)
	}
	return nil, nil
}

// pdf writes the selected pages of the view as a PDF file, with the
// copies the printer does not make.
func (j *job) pdf() ([]byte, error) {
	var sel []int
	order := j.order()
	all := len(order) == len(j.view.Pages)
	for i, p := range order {
		if p != i {
			all = false
		}
	}
	for _, p := range order {
		if p >= 0 {
			sel = append(sel, p)
		}
	}
	if len(sel) != len(order) {
		j.res.warn("the copies of a two-sided document of an odd number of pages share sheets: the printer makes no copies, and a PDF file has no blank sheet to put between them")
	}
	if all {
		sel = nil
	}
	var buf bytes.Buffer
	res, err := pdfbdf.Write(&buf, j.doc, &pdfbdf.Options{Views: []string{j.view.ID}, SkipRoles: j.o.SkipRoles, Pages: sel,
		FontFS: j.o.FontFS, FontDirs: j.o.FontDirs, NoSystemFonts: j.o.NoSystemFonts})
	if err != nil {
		return nil, err
	}
	j.res.Sheets = res.Pages
	j.res.Warnings = append(j.res.Warnings, res.Warnings...)
	return buf.Bytes(), nil
}

// lossy tells whether the warnings of a PDF file say it leaves out what
// the document draws, which a raster page would have.
func lossy(warnings []string) bool {
	for _, w := range warnings {
		for _, s := range []string{"is not applied", "is drawn as source-over", "is not drawn", "casts no shadow", "are not drawn"} {
			if strings.Contains(w, s) {
				return true
			}
		}
	}
	return false
}

// rasterWriter is a writer of raster pages: of PWG Raster, Apple Raster or
// PCLm.
type rasterWriter interface {
	StartPage(raster.Page) error
	WriteRows(image.Image) error
	Close() error
}

// raster writes the selected pages in a raster format on the paper of the
// media, drawing them in bands.
func (j *job) raster(w io.Writer, media, format string) error {
	order := j.order()
	var setup rasterSetup
	var back string
	var rw rasterWriter
	var err error
	switch format {
	case FormatURF:
		if setup, back, err = chooseURF(j.p, j.t); err != nil {
			return err
		}
		rw = raster.NewURFWriter(w, len(order))
	case FormatPCLm:
		var opts raster.PCLmOptions
		if setup, opts, err = choosePCLm(j.p, j.t); err != nil {
			return err
		}
		back = j.p.PCLmBackSide
		if rw, err = raster.NewPCLmWriter(w, opts); err != nil {
			return err
		}
	default:
		if setup, err = chooseRaster(j.p, j.t); err != nil {
			return err
		}
		back = j.p.PWGRasterSheetBack
		rw = raster.NewWriter(w)
	}
	name := media
	if name == "" {
		name = j.p.MediaDefault
	}
	m, err := ParseMedia(name)
	if err != nil {
		return fmt.Errorf("print: the paper %q: %w", name, err)
	}
	sheet := raster.Sheet{Width: m.Width, Height: m.Height, DPI: setup.dpi}
	scaling := j.t.Scaling
	if media == "" && scaling == "" {
		scaling = raster.ScaleFit
	}
	ro := &imagebdf.Options{FontFS: j.o.FontFS, FontDirs: j.o.FontDirs, NoSystemFonts: j.o.NoSystemFonts, ImageCacheBytes: -1,
		BandRows: j.o.BandRows, BandWorkers: j.o.BandWorkers}
	mono := imagebdf.Mono(0)
	switch {
	case setup.typ == raster.Black1, setup.typ == raster.SGray1:
		mono = imagebdf.Bilevel
		ro.Antialias = imagebdf.AntialiasOff
	case setup.mono:
		mono = imagebdf.Gray
	}
	r := imagebdf.New(j.doc, ro)
	view := skipLayers(j.view, j.o.SkipRoles)
	duplex, tumble := j.t.duplex()
	cf, ff, err := raster.Back(back, tumble)
	if err != nil {
		cf, ff = 1, 1
		j.res.warn("%v: the back of a sheet is sent as the front", err)
	}
	j.res.Sheets = len(order)
	pw, ph := sheet.Pixels()
	for i, page := range order {
		if err := j.ctx.Err(); err != nil {
			return err
		}
		back := duplex && i%2 == 1
		hdr := raster.Page{Width: pw, Height: ph, XDPI: setup.dpi, YDPI: setup.dpi, Type: setup.typ,
			PageWidth: int(math.Round(m.Width)), PageHeight: int(math.Round(m.Height)), PageSizeName: name, MediaType: j.t.MediaType,
			Duplex: duplex, Tumble: tumble, CrossFeed: 1, Feed: 1, Copies: 1, TotalPages: len(order), Quality: j.t.Quality}
		if back {
			hdr.CrossFeed, hdr.Feed = cf, ff
		}
		if err := rw.StartPage(hdr); err != nil {
			return err
		}
		if page < 0 {
			if err := blank(rw, pw, ph); err != nil {
				return err
			}
			continue
		}
		p := view.Pages[page]
		pl, err := raster.Fit(float64(p.W), float64(p.H), sheet, scaling, j.t.Orientation)
		if err != nil {
			return err
		}
		if back {
			pl = raster.Mirror(pl, cf, ff)
		}
		if err := r.PlaceRows(view, page, pl, mono, func(_ int, rows image.Image) error { return rw.WriteRows(rows) }); err != nil {
			return fmt.Errorf("print: page %d: %w", page+1, err)
		}
	}
	j.res.Warnings = append(j.res.Warnings, r.Warnings()...)
	return rw.Close()
}

// blank writes the rows of a blank sheet.
func blank(rw rasterWriter, w, h int) error {
	band := image.NewGray(image.Rect(0, 0, w, min(h, 64)))
	for i := range band.Pix {
		band.Pix[i] = 0xff
	}
	for y := 0; y < h; y += band.Rect.Dy() {
		rows := band
		if y+rows.Rect.Dy() > h {
			rows = band.SubImage(image.Rect(0, 0, w, h-y)).(*image.Gray)
		}
		if err := rw.WriteRows(rows); err != nil {
			return err
		}
	}
	return nil
}

// wait polls the job until it is done, the wait is over or the context
// is canceled, which cancels the job.
func (j *job) wait(js *ipp.JobStatus) error {
	ctx := j.ctx
	var cancel context.CancelFunc
	if j.o.Wait > 0 {
		ctx, cancel = context.WithTimeout(ctx, j.o.Wait)
		defer cancel()
	}
	delay := 500 * time.Millisecond
	last := js
	for !last.State.Done() {
		select {
		case <-ctx.Done():
			if j.ctx.Err() != nil {
				j.cancel(js.ID)
				return fmt.Errorf("%w: %v", ErrCanceled, j.ctx.Err())
			}
			return fmt.Errorf("%w: %s after %v", ErrTimeout, last, j.o.Wait)
		case <-time.After(delay):
		}
		delay = min(delay*3/2, 5*time.Second)
		now, err := Check(ctx, j.c, last)
		if now == nil {
			if ctx.Err() != nil {
				continue
			}
			return err
		}
		if now.State != last.State || strings.Join(now.Reasons, ",") != strings.Join(last.Reasons, ",") {
			j.event(Event{Kind: "state", Job: now})
		}
		last = now
		j.res.Job = last
		if err != nil {
			return err
		}
	}
	return nil
}

// Check asks the printer once for the state of a job it was sent, and
// returns it merged into the last state known. A job canceled or
// aborted at the printer returns its state with ErrCanceled or
// ErrAborted; a job not done yet, or completed, returns no error. When
// the printer could not be asked, the state is nil.
func Check(ctx context.Context, c *ipp.Client, last *ipp.JobStatus) (*ipp.JobStatus, error) {
	now, err := c.GetJobAttributes(ctx, last.ID)
	if err != nil {
		return nil, err
	}
	now = merge(last, now)
	switch now.State {
	case ipp.JobCanceled:
		return now, fmt.Errorf("%w: %s", ErrCanceled, now)
	case ipp.JobAborted:
		return now, fmt.Errorf("%w: %s", ErrAborted, now)
	}
	return now, nil
}

// cancel cancels a job at the printer, whatever the context of the job.
func (j *job) cancel(id int32) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	j.c.CancelJob(ctx, id)
}

// lost tells, after Send-Document failed with err, whether the printer
// lost the job it made: it no longer knows the id, or its job of the id
// has another name (the printer was restarted, counts its ids from 1
// again, and another client's job took the id). A printer that tells no
// name, or a name the job's begins with (one it cut), has the job.
func (j *job) lost(id int32, name string, err error) bool {
	if !ipp.IsStatus(err, ipp.StatusClientErrorNotPossible) && !ipp.IsStatus(err, ipp.StatusClientErrorNotFound) {
		return false
	}
	ctx, cancel := context.WithTimeout(context.WithoutCancel(j.ctx), 10*time.Second)
	defer cancel()
	js, gerr := j.c.GetJobAttributes(ctx, id)
	switch {
	case ipp.IsStatus(gerr, ipp.StatusClientErrorNotFound):
		return true
	case gerr != nil:
		return false
	}
	return js.Name != "" && !strings.HasPrefix(name, js.Name)
}

// interrupted turns the error of a request canceled with the context
// into ErrCanceled.
func (j *job) interrupted(err error) error {
	if j.ctx.Err() != nil {
		return fmt.Errorf("%w: %v", ErrCanceled, j.ctx.Err())
	}
	return err
}

type countingReader struct {
	r io.Reader
	n int64
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.n += int64(n)
	return n, err
}
