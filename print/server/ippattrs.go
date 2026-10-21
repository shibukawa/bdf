package server

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/raster"
	"github.com/shibukawa/bdf/print/spool"
)

// documentFormats are the formats the printers take: PDF first, so that
// clients send vectors; the rasters of driverless clients; images; .bdf;
// and anything else, which the server tells by its content (the formats
// of the converters).
var documentFormats = []string{
	"application/pdf", "image/urf", "image/pwg-raster", "image/jpeg", "image/png", "application/x-bdf", "application/octet-stream",
}

// The resolutions and types of the rasters the printers take: those
// clients draw at, not those of the printers behind (the server draws
// for them).
var (
	rasterResolutions = []int{150, 300, 600}
	rasterTypes       = []string{string(raster.Black1), string(raster.SGray8), string(raster.SRGB8)}
	// urfSupported is urf-supported: Apple Raster version 1.4, 8-bit gray
	// and 24-bit sRGB, 150 to 600 dpi, both sides (DM1: the back is not
	// turned), copies, input slots, media types, quality and finishings
	urfSupported = []string{"V1.4", "W8", "SRGB24", "RS150-300-600", "DM1", "CP255", "IS1", "MT1-2-3-4-5-8-11", "OB10", "PQ3-4-5", "FN3"}
)

// defaultMedia is the paper of a queue whose printer does not say: an
// office's, A4 first.
var defaultMedia = []string{
	"iso_a4_210x297mm", "iso_a3_297x420mm", "iso_a5_148x210mm", "jis_b4_257x364mm", "jis_b5_182x257mm", "na_letter_8.5x11in",
	"na_legal_8.5x14in", "jpn_hagaki_100x148mm", "na_index-4x6_4x6in", "iso_dl_110x220mm", "jpn_chou3_120x235mm",
}

// the qualities, modes, scalings and holds the printers take: the server
// makes the pages, so these do not depend on the printer behind
var (
	colorModes       = []string{"auto", "color", "monochrome"}
	scalings         = []string{"auto", "auto-fit", "fill", "fit", "none"}
	holds            = []string{"no-hold", "indefinite", "day-time", "evening", "night", "second-shift", "third-shift", "weekend"}
	contentOptimizes = []string{"auto", "graphic", "photo", "text", "text-and-graphic"}
	renderingIntents = []string{"auto", "perceptual", "relative", "relative-bpc", "saturation", "absolute"}
	handlings        = []string{"separate-documents-collated-copies", "separate-documents-uncollated-copies"}
)

// jobCreationAttrs are the job template attributes taken.
var jobCreationAttrs = []string{
	"copies", "finishings", "ipp-attribute-fidelity", "job-hold-until", "job-name", "job-priority", "media", "media-col",
	"multiple-document-handling", "number-up", "orientation-requested", "output-bin", "overrides", "page-ranges",
	"print-color-mode", "print-content-optimize", "print-quality", "print-rendering-intent", "print-scaling",
	"printer-resolution", "sides",
}

// queueCaps is what a queue prints on: its first printer's capabilities,
// or none when it has not answered (or tells nothing, as 9100 and LPD).
func (x *ippServer) queueCaps(name string) *ipp.Printer {
	qc := x.s.cfg.Queues[name]
	if len(qc.Printers) == 0 {
		return nil
	}
	caps, _ := x.s.sp.PrinterCapabilities(qc.Printers[0])
	return caps
}

// media are the media-supported of a queue.
func mediaOf(caps *ipp.Printer) (supported, ready []string) {
	if caps != nil {
		for _, m := range caps.Media {
			if _, err := print.ParseMedia(m); err == nil {
				supported = append(supported, m)
			}
		}
		for _, m := range caps.MediaReady {
			if slices.Contains(supported, m) {
				ready = append(ready, m)
			}
		}
	}
	if len(supported) == 0 {
		supported = defaultMedia
	}
	if len(ready) == 0 {
		ready = supported[:1]
	}
	return supported, ready
}

// jobRequest is what a request asks of a job.
type jobRequest struct {
	ticket   print.Ticket
	options  spool.JobOptions
	password string
}

// jobRequest reads the job template attributes of a request into a
// ticket. What is not supported is returned in the unsupported group;
// with ipp-attribute-fidelity the request is refused, without it the
// status says the attributes were ignored.
func (x *ippServer) jobRequest(q *ippRequest) (jobRequest, bool) {
	var jr jobRequest
	t := &jr.ticket
	caps := x.queueCaps(q.queue.name)
	media, _ := mediaOf(caps)
	t.JobName = q.op.Get("job-name").Text()
	if len(t.JobName) > 255 {
		t.JobName = t.JobName[:255]
	}
	fidelity, _ := q.op.Get("ipp-attribute-fidelity").Bool()
	var unsupported []ipp.Attribute
	bad := func(a ipp.Attribute) { unsupported = append(unsupported, a) }
	keyword := func(a ipp.Attribute, list []string) (string, bool) {
		v, ok := a.First()
		if !ok || len(a.Values) != 1 || (v.Tag != ipp.TagKeyword && v.Tag != ipp.TagName && v.Tag != ipp.TagNameLang) || !slices.Contains(list, v.Str) {
			bad(a)
			return "", false
		}
		return v.Str, true
	}
	integer := func(a ipp.Attribute, tag ipp.Tag, lo, hi int) (int, bool) {
		v, ok := a.First()
		if !ok || len(a.Values) != 1 || v.Tag != tag || int(v.Int) < lo || int(v.Int) > hi {
			bad(a)
			return 0, false
		}
		return int(v.Int), true
	}
	var job []ipp.Attribute
	if g := q.req.Group(ipp.TagJob); g != nil {
		job = g.Attrs
	}
	for _, a := range job {
		switch a.Name {
		case "copies":
			if n, ok := integer(a, ipp.TagInteger, 1, 999); ok {
				t.Copies = n
			}
		case "sides":
			list := []string{"one-sided", "two-sided-long-edge", "two-sided-short-edge"}
			if caps != nil && len(caps.Sides) > 0 {
				list = caps.Sides
			}
			t.Sides, _ = keyword(a, list)
		case "media":
			t.Media, _ = keyword(a, media)
		case "media-col":
			if !x.mediaCol(a, t, caps, media) {
				bad(a)
			}
		case "print-color-mode":
			if v, ok := keyword(a, colorModes); ok && v != "auto" {
				t.ColorMode = v
			}
		case "print-quality":
			if n, ok := integer(a, ipp.TagEnum, 3, 5); ok {
				t.Quality = n
			}
		case "printer-resolution":
			rs := a.Resolutions()
			if len(rs) != 1 || rs[0].X != rs[0].Y || rs[0].Units != ipp.UnitsDPI || !slices.Contains(rasterResolutions, int(rs[0].X)) {
				bad(a)
			} else {
				t.Resolution = int(rs[0].X)
			}
		case "page-ranges":
			pages, ok := pageRanges(a)
			if !ok {
				bad(a)
			} else {
				t.Pages = pages
			}
		case "print-scaling":
			if v, ok := keyword(a, scalings); ok {
				t.Scaling = raster.Scaling(v)
			}
		case "orientation-requested":
			if v, ok := a.First(); ok && v.Tag == ipp.TagNoValue {
				continue
			}
			if n, ok := integer(a, ipp.TagEnum, 3, 6); ok {
				t.Orientation = raster.Orientation(n)
			}
		case "finishings":
			ok := len(a.Values) > 0
			var fs []int
			for _, v := range a.Values {
				if v.Tag != ipp.TagEnum || v.Int != 3 && (caps == nil || !slices.Contains(caps.Finishings, v.Int)) {
					ok = false
				} else if v.Int != 3 {
					fs = append(fs, int(v.Int))
				}
			}
			if !ok {
				bad(a)
			} else {
				t.Finishings = fs
			}
		case "output-bin":
			if caps == nil || len(caps.OutputBins) == 0 {
				bad(a)
			} else {
				t.OutputBin, _ = keyword(a, caps.OutputBins)
			}
		case "job-hold-until":
			if v, ok := keyword(a, holds); ok && v != "no-hold" {
				t.HoldUntil = v
			}
		case "job-priority":
			if n, ok := integer(a, ipp.TagInteger, 1, 100); ok {
				t.Priority = n
			}
		case "multiple-document-handling":
			if v, ok := keyword(a, handlings); ok {
				t.Collate = "collated"
				if v == "separate-documents-uncollated-copies" {
					t.Collate = "uncollated"
				}
			}
		case "number-up":
			integer(a, ipp.TagInteger, 1, 1)
		case "print-content-optimize":
			keyword(a, contentOptimizes)
		case "print-rendering-intent":
			keyword(a, renderingIntents)
		case "job-sheets":
			keyword(a, []string{"none"})
		case "overrides":
			// the pages and documents of overrides are not drawn
			// differently: taken, and left
		default:
			bad(ipp.Attr(a.Name, ipp.OutOfBand(ipp.TagUnsupported)))
		}
	}
	if len(unsupported) > 0 {
		q.res.Add(ipp.TagUnsupportedGroup, unsupported...)
		if fidelity {
			q.fail(ipp.StatusClientErrorAttributesOrValuesNotSupported, "the job asks for what the printer does not support")
			return jr, false
		}
		q.res.Code = uint16(ipp.StatusOKIgnoredOrSubstituted)
	}
	if err := t.Check(); err != nil {
		q.fail(ipp.StatusClientErrorBadRequest, "%v", err)
		return jr, false
	}
	return jr, true
}

// mediaCol reads a media-col into a ticket: the size (by its name, or its
// dimensions matched to a size of the queue), the source and the type.
func (x *ippServer) mediaCol(a ipp.Attribute, t *print.Ticket, caps *ipp.Printer, media []string) bool {
	cs := a.Collections()
	if len(cs) != 1 {
		return false
	}
	m := ipp.ParseMedia(cs[0])
	if m.SizeName != "" {
		if !slices.Contains(media, m.SizeName) {
			return false
		}
		t.Media = m.SizeName
	} else if m.X > 0 && m.Y > 0 {
		w, h := float64(m.X)*72/2540, float64(m.Y)*72/2540
		t.Media = print.ChooseMedia(w, h, nil, media)
		if t.Media == "" {
			return false
		}
	}
	if m.Source != "" {
		if caps == nil || !slices.Contains(caps.MediaSources, m.Source) {
			if m.Source != "auto" {
				return false
			}
		} else if m.Source != "auto" {
			t.MediaSource = m.Source
		}
	}
	if m.Type != "" && m.Type != "auto" {
		if caps == nil || !slices.Contains(caps.MediaTypes, m.Type) {
			return false
		}
		t.MediaType = m.Type
	}
	return true
}

// pageRanges reads page-ranges into the pages, in order, 10,000 at most.
func pageRanges(a ipp.Attribute) ([]int, bool) {
	var pages []int
	last := 0
	for _, v := range a.Values {
		if v.Tag != ipp.TagRange || v.Range.Lo < 1 || v.Range.Hi < v.Range.Lo || int(v.Range.Lo) <= last {
			return nil, false
		}
		for p := v.Range.Lo; p <= v.Range.Hi; p++ {
			if len(pages) >= 10000 {
				return nil, false
			}
			pages = append(pages, int(p))
		}
		last = int(v.Range.Hi)
	}
	return pages, len(pages) > 0
}

// The attributes of jobs, as requested-attributes names their groups.
var (
	jobDescription = []string{
		"job-id", "job-uri", "job-uuid", "job-printer-uri", "job-more-info", "job-name", "job-originating-user-name",
		"job-state", "job-state-reasons", "job-state-message", "number-of-documents", "time-at-creation", "time-at-processing",
		"time-at-completed", "job-printer-up-time", "date-time-at-creation", "date-time-at-processing", "date-time-at-completed",
		"job-k-octets", "job-impressions-completed", "job-media-sheets-completed", "document-format-supplied", "document-name-supplied",
		"number-of-intervening-jobs",
	}
	jobTemplate = []string{
		"copies", "finishings", "job-hold-until", "job-priority", "media", "media-col", "multiple-document-handling", "orientation-requested",
		"output-bin", "page-ranges", "print-color-mode", "print-quality", "print-scaling", "printer-resolution", "sides",
	}
	jobGroups = map[string][]string{"job-description": jobDescription, "job-template": jobTemplate}
	jobAll    = []string{"all"}
)

// jobStatusAttrs are the attributes of a job a request that makes one is
// answered with.
func (x *ippServer) jobStatusAttrs(q *ippRequest, j *spool.Job) []ipp.Attribute {
	attrs := x.jobAttrs(q, j, true)
	return filterAttrs(attrs, []string{"job-id", "job-uri", "job-uuid", "job-state", "job-state-reasons", "job-state-message"}, nil)
}

// jobAttrs are the attributes of a job of the spool. A job of another
// that the caller may not see is told without its name and user.
func (x *ippServer) jobAttrs(q *ippRequest, j *spool.Job, visible bool) []ipp.Attribute {
	name, user := j.Name, j.User
	if !visible {
		name, user = "", "anonymous"
	}
	reasons := slices.Clone(j.Reasons)
	if len(reasons) == 0 {
		reasons = []string{"none"}
	}
	msg := j.Message
	if !visible {
		msg = ""
	}
	attrs := []ipp.Attribute{
		ipp.Attr("job-id", ipp.Integer(int(j.ID))),
		ipp.Attr("job-uri", ipp.URI(q.jobURI(j.ID))),
		ipp.Attr("job-uuid", ipp.URI(j.UUID)),
		ipp.Attr("job-printer-uri", ipp.URI(q.printerURI())),
		ipp.Attr("job-more-info", ipp.URI(q.moreInfo())),
		ipp.Attr("job-name", ipp.Name(name)),
		ipp.Attr("job-originating-user-name", ipp.Name(user)),
		ipp.Attr("job-state", ipp.Enum(int(j.State))),
		ipp.Attr("job-state-reasons", keywordValues(publicReasons(reasons))...),
		ipp.Attr("job-state-message", ipp.Text(clip127(msg))),
		ipp.Attr("number-of-documents", ipp.Integer(1)),
		ipp.Attr("job-printer-up-time", ipp.Integer(x.upTime(time.Now()))),
		ipp.Attr("job-k-octets", ipp.Integer(int(min((j.Document.Size+1023)/1024, math.MaxInt32)))),
		ipp.Attr("job-impressions-completed", ipp.Integer(j.Impressions)),
		ipp.Attr("job-media-sheets-completed", ipp.Integer(sheetsDone(j))),
		ipp.Attr("number-of-intervening-jobs", ipp.Integer(0)),
	}
	attrs = append(attrs, x.timeAttrs("creation", j.Submitted)...)
	attrs = append(attrs, x.timeAttrs("processing", j.Started)...)
	attrs = append(attrs, x.timeAttrs("completed", j.Finished)...)
	if visible {
		if j.Document.Type != "" {
			attrs = append(attrs, ipp.Attr("document-format-supplied", ipp.MimeType(j.Document.Type)))
		}
		attrs = append(attrs, ipp.Attr("document-name-supplied", ipp.Name(j.Document.Name)))
	}
	return append(attrs, ticketAttrs(&j.Ticket)...)
}

// sheetsDone is the sheets of a job done.
func sheetsDone(j *spool.Job) int {
	if j.State == ipp.JobCompleted {
		return j.Sheets
	}
	return 0
}

// publicReasons turns the spool's own reasons (bdf-…) into IPP's for
// clients that know these only.
func publicReasons(reasons []string) []string {
	out := make([]string, 0, len(reasons))
	for _, r := range reasons {
		switch r {
		case spool.ReasonRetryWait:
			r = "printer-stopped"
		case spool.ReasonRetryExhausted, spool.ReasonUnsupported:
			r = "aborted-by-system"
		case spool.ReasonDuplicateCheck:
			r = "job-hold-until-specified"
		}
		if !slices.Contains(out, r) {
			out = append(out, r)
		}
	}
	return out
}

// incomingAttrs are the attributes of a job that has not reached the
// spool.
func (x *ippServer) incomingAttrs(q *ippRequest, in *incoming) []ipp.Attribute {
	x.mu.Lock()
	defer x.mu.Unlock()
	visible := q.caller.sees(&spool.Job{User: in.user})
	name, user := in.name, in.user
	if !visible {
		name, user = "", "anonymous"
	}
	attrs := []ipp.Attribute{
		ipp.Attr("job-id", ipp.Integer(int(in.id))),
		ipp.Attr("job-uri", ipp.URI(q.jobURI(in.id))),
		ipp.Attr("job-uuid", ipp.URI(in.uuid)),
		ipp.Attr("job-printer-uri", ipp.URI(q.printerURI())),
		ipp.Attr("job-more-info", ipp.URI(q.moreInfo())),
		ipp.Attr("job-name", ipp.Name(name)),
		ipp.Attr("job-originating-user-name", ipp.Name(user)),
		ipp.Attr("job-state", ipp.Enum(int(in.state))),
		ipp.Attr("job-state-reasons", keywordValues(in.reasons)...),
		ipp.Attr("job-state-message", ipp.Text("")),
		ipp.Attr("number-of-documents", ipp.Integer(0)),
		ipp.Attr("job-printer-up-time", ipp.Integer(x.upTime(time.Now()))),
		ipp.Attr("job-k-octets", ipp.Integer(0)),
		ipp.Attr("job-impressions-completed", ipp.Integer(0)),
		ipp.Attr("job-media-sheets-completed", ipp.Integer(0)),
		ipp.Attr("number-of-intervening-jobs", ipp.Integer(0)),
	}
	attrs = append(attrs, x.timeAttrs("creation", in.created)...)
	attrs = append(attrs, x.timeAttrs("processing", time.Time{})...)
	attrs = append(attrs, x.timeAttrs("completed", in.ended)...)
	return append(attrs, ticketAttrs(&in.job.ticket)...)
}

// timeAttrs are time-at-<what> and date-time-at-<what>, no-value when the
// job is not there yet.
func (x *ippServer) timeAttrs(what string, t time.Time) []ipp.Attribute {
	if t.IsZero() {
		return []ipp.Attribute{ipp.Attr("time-at-"+what, ipp.OutOfBand(ipp.TagNoValue)), ipp.Attr("date-time-at-"+what, ipp.OutOfBand(ipp.TagNoValue))}
	}
	return []ipp.Attribute{ipp.Attr("time-at-"+what, ipp.Integer(x.upTime(t))), ipp.Attr("date-time-at-"+what, ipp.DateTime(t))}
}

// ticketAttrs are the job template attributes of a ticket, as asked.
func ticketAttrs(t *print.Ticket) []ipp.Attribute {
	attrs := []ipp.Attribute{ipp.Attr("copies", ipp.Integer(max(t.Copies, 1)))}
	if t.Sides != "" {
		attrs = append(attrs, ipp.Attr("sides", ipp.Keyword(t.Sides)))
	}
	if t.Media != "" {
		attrs = append(attrs, ipp.Attr("media", ipp.Keyword(t.Media)))
	}
	if t.ColorMode != "" {
		attrs = append(attrs, ipp.Attr("print-color-mode", ipp.Keyword(t.ColorMode)))
	}
	if t.Quality != 0 {
		attrs = append(attrs, ipp.Attr("print-quality", ipp.Enum(t.Quality)))
	}
	if t.Resolution != 0 {
		attrs = append(attrs, ipp.Attr("printer-resolution", ipp.NewResolution(t.Resolution, t.Resolution, ipp.UnitsDPI)))
	}
	if t.Scaling != "" {
		attrs = append(attrs, ipp.Attr("print-scaling", ipp.Keyword(string(t.Scaling))))
	}
	if t.Orientation != 0 {
		attrs = append(attrs, ipp.Attr("orientation-requested", ipp.Enum(int(t.Orientation))))
	}
	if len(t.Finishings) > 0 {
		attrs = append(attrs, ipp.Attr("finishings", enumValues(t.Finishings)...))
	}
	if t.OutputBin != "" {
		attrs = append(attrs, ipp.Attr("output-bin", ipp.Keyword(t.OutputBin)))
	}
	hold := t.HoldUntil
	if hold == "" {
		hold = "no-hold"
	}
	attrs = append(attrs, ipp.Attr("job-hold-until", ipp.Keyword(hold)), ipp.Attr("job-priority", ipp.Integer(cmpOr(t.Priority, 50))))
	switch t.Collate {
	case "collated":
		attrs = append(attrs, ipp.Attr("multiple-document-handling", ipp.Keyword("separate-documents-collated-copies")))
	case "uncollated":
		attrs = append(attrs, ipp.Attr("multiple-document-handling", ipp.Keyword("separate-documents-uncollated-copies")))
	}
	if len(t.Pages) > 0 {
		attrs = append(attrs, ipp.Attr("page-ranges", pageRangeValues(t.Pages)...))
	}
	return attrs
}

func cmpOr(n, def int) int {
	if n == 0 {
		return def
	}
	return n
}

// pageRangeValues turns pages back into ranges.
func pageRangeValues(pages []int) []ipp.Value {
	var out []ipp.Value
	for i := 0; i < len(pages); {
		j := i
		for j+1 < len(pages) && pages[j+1] == pages[j]+1 {
			j++
		}
		out = append(out, ipp.NewRange(pages[i], pages[j]))
		i = j + 1
	}
	return out
}

func keywordValues(s []string) []ipp.Value {
	out := make([]ipp.Value, len(s))
	for i, v := range s {
		out[i] = ipp.Keyword(v)
	}
	return out
}

func enumValues(ns []int) []ipp.Value {
	out := make([]ipp.Value, len(ns))
	for i, n := range ns {
		out[i] = ipp.Enum(n)
	}
	return out
}

func resolutionValues(dpis []int) []ipp.Value {
	out := make([]ipp.Value, len(dpis))
	for i, d := range dpis {
		out[i] = ipp.NewResolution(d, d, ipp.UnitsDPI)
	}
	return out
}

// moreInfo is the page of the queue on the IPP server.
func (q *ippRequest) moreInfo() string {
	scheme := "http"
	if q.secure {
		scheme = "https"
	}
	return scheme + "://" + q.r.Host + q.queue.path
}

// queueState is the printer-state of a queue: stopped when it is paused
// or its printers are, processing while it has a job at a printer.
func (x *ippServer) queueState(name string) (ipp.PrinterState, []string, string, bool) {
	return x.s.queueState(name)
}

// queueState is the printer-state of a queue, its reasons and message, and
// whether it takes jobs: as an IPP printer, and as a proxy's printer.
func (s *Server) queueState(name string) (ipp.PrinterState, []string, string, bool) {
	x := struct{ s *Server }{s}
	info, err := x.s.sp.Queue(name)
	if err != nil {
		return ipp.PrinterStopped, []string{"other-error"}, err.Error(), false
	}
	reasons := []string{}
	msg := ""
	st := ipp.PrinterIdle
	if info.Jobs[ipp.JobProcessing.String()] > 0 || info.Jobs[ipp.JobProcessingStopped.String()] > 0 {
		st = ipp.PrinterProcessing
	}
	if info.Paused {
		st, reasons = ipp.PrinterStopped, append(reasons, "paused")
		msg = "the queue is paused"
	}
	qc := x.s.cfg.Queues[name]
	down := 0
	for _, p := range x.s.sp.Printers() {
		if !slices.Contains(qc.Printers, p.Name) {
			continue
		}
		if !p.Up && !p.Checked.IsZero() {
			down++
		}
		// the printers' warnings and errors, as a server reports those of
		// its printers
		for _, r := range p.Reasons {
			if r != "none" && (strings.HasSuffix(r, "-error") || strings.HasSuffix(r, "-warning") || strings.HasSuffix(r, "-report")) && !slices.Contains(reasons, r) {
				reasons = append(reasons, r)
			}
		}
		if p.Message != "" && msg == "" {
			msg = p.Message
		}
	}
	if down > 0 && down == len(qc.Printers) {
		reasons = append(reasons, "offline-report")
		if msg == "" {
			msg = "the printer is not reached"
		}
	}
	if len(reasons) == 0 {
		reasons = []string{"none"}
	}
	return st, reasons, msg, !info.Rejecting
}

// getPrinterAttributes answers what a queue is and supports.
func (x *ippServer) getPrinterAttributes(q *ippRequest) {
	if f := q.op.Get("document-format"); f != nil && !slices.Contains(documentFormats, f.Text()) {
		q.fail(ipp.StatusClientErrorDocumentFormatNotSupported, "document-format %s", f.Text())
		q.res.Add(ipp.TagUnsupportedGroup, *f)
		return
	}
	attrs := x.printerAttrs(q)
	req := requested(q.op, []string{"all"})
	// media-col-database is told when asked by name, as printers do
	if !slices.Contains(req, "media-col-database") {
		attrs = slices.DeleteFunc(attrs, func(a ipp.Attribute) bool { return a.Name == "media-col-database" })
	}
	if slices.Contains(req, "all") {
		q.res.Add(ipp.TagPrinter, attrs...)
		return
	}
	q.res.Add(ipp.TagPrinter, filterAttrs(attrs, req, x.printerAttrGroups(attrs))...)
}

// printerAttrGroups are the groups of printer attributes: the template
// ones are the -default, -supported, -ready of job template attributes,
// the description ones the others.
func (x *ippServer) printerAttrGroups(attrs []ipp.Attribute) map[string][]string {
	var tmpl, desc []string
	for _, a := range attrs {
		base := a.Name
		for _, suf := range []string{"-default", "-supported", "-ready", "-database"} {
			base = strings.TrimSuffix(base, suf)
		}
		if base != a.Name && (slices.Contains(jobCreationAttrs, base) || base == "media-col" || base == "media") {
			tmpl = append(tmpl, a.Name)
		} else {
			desc = append(desc, a.Name)
		}
	}
	return map[string][]string{"job-template": tmpl, "printer-description": desc}
}

// printerAttrs are the attributes of a queue as an IPP Everywhere and
// AirPrint printer.
func (x *ippServer) printerAttrs(q *ippRequest) []ipp.Attribute {
	caps := x.queueCaps(q.queue.name)
	st, reasons, msg, accepting := x.queueState(q.queue.name)
	x.mu.Lock()
	changed := x.changed[q.queue.name]
	x.mu.Unlock()
	host := q.r.Host
	var uris, auth, security []ipp.Value
	authKw := "requesting-user-name"
	if x.cfg.Auth == authBasic {
		authKw = "basic"
	}
	if !x.secureOnly {
		uris = append(uris, ipp.URI("ipp://"+host+q.queue.path))
		auth, security = append(auth, ipp.Keyword(authKw)), append(security, ipp.Keyword("none"))
	}
	uris = append(uris, ipp.URI("ipps://"+host+q.queue.path))
	auth, security = append(auth, ipp.Keyword(authKw)), append(security, ipp.Keyword("tls"))
	media, ready := mediaOf(caps)
	sides := []string{"one-sided", "two-sided-long-edge", "two-sided-short-edge"}
	color := true
	makeModel := "BDF Print Server"
	if caps != nil {
		if len(caps.Sides) > 0 {
			sides = caps.Sides
		}
		if len(caps.ColorModes) > 0 && !slices.Contains(caps.ColorModes, "color") {
			color = false
		}
		if caps.MakeModel != "" {
			makeModel = "BDF Print Server (" + caps.MakeModel + ")"
		}
	}
	modes := colorModes
	if !color {
		modes = []string{"auto", "monochrome"}
	}
	finishings := []int{3}
	if caps != nil {
		for _, f := range caps.Finishings {
			if f != 3 {
				finishings = append(finishings, int(f))
			}
		}
	}
	sources, types := []string{"auto"}, []string{"auto"}
	if caps != nil {
		for _, s := range caps.MediaSources {
			if s != "auto" {
				sources = append(sources, s)
			}
		}
		for _, t := range caps.MediaTypes {
			if t != "auto" {
				types = append(types, t)
			}
		}
	}
	top, bottom, left, right := margins(caps)
	var mediaDB, mediaReady, mediaSizes []ipp.Value
	for _, m := range media {
		pm, _ := print.ParseMedia(m)
		w, h := int(math.Round(pm.Width/72*2540)), int(math.Round(pm.Height/72*2540))
		size := ipp.Collection(ipp.Attr("x-dimension", ipp.Integer(w)), ipp.Attr("y-dimension", ipp.Integer(h)))
		mediaSizes = append(mediaSizes, size)
		col := ipp.Collection(ipp.Attr("media-key", ipp.Keyword(m)), ipp.Attr("media-size", size), ipp.Attr("media-size-name", ipp.Keyword(m)),
			ipp.Attr("media-source", ipp.Keyword("auto")), ipp.Attr("media-type", ipp.Keyword("auto")),
			ipp.Attr("media-top-margin", ipp.Integer(top)), ipp.Attr("media-bottom-margin", ipp.Integer(bottom)),
			ipp.Attr("media-left-margin", ipp.Integer(left)), ipp.Attr("media-right-margin", ipp.Integer(right)))
		mediaDB = append(mediaDB, col)
		if slices.Contains(ready, m) {
			mediaReady = append(mediaReady, col)
		}
	}
	ops := make([]ipp.Value, len(operations))
	for i, op := range operations {
		ops[i] = ipp.Enum(int(op))
	}
	features := []ipp.Value{ipp.Keyword("ipp-everywhere"), ipp.Keyword("ipp-everywhere-server"), ipp.Keyword("page-overrides")}
	which := []ipp.Value{ipp.Keyword("completed"), ipp.Keyword("not-completed"), ipp.Keyword("all")}
	if qc := x.s.cfg.Queues[q.queue.name]; len(qc.Printers) == 1 && x.s.cfg.Printers[qc.Printers[0]].Proxy {
		// an infrastructure printer: proxies fetch its jobs
		for _, op := range []ipp.Op{ipp.OpAcknowledgeDocument, ipp.OpAcknowledgeJob, ipp.OpFetchDocument, ipp.OpFetchJob, ipp.OpUpdateActiveJobs,
			ipp.OpDeregisterOutputDevice, ipp.OpUpdateDocumentStatus, ipp.OpUpdateJobStatus, ipp.OpUpdateOutputDeviceAttributes} {
			ops = append(ops, ipp.Enum(int(op)))
		}
		features = append(features, ipp.Keyword("infrastructure-printer"))
		which = append(which, ipp.Keyword("fetchable"))
	}
	queued := 0
	if info, err := x.s.sp.Queue(q.queue.name); err == nil {
		for state, n := range info.Jobs {
			if state == ipp.JobPending.String() || state == ipp.JobPendingHeld.String() || state == ipp.JobProcessing.String() || state == ipp.JobProcessingStopped.String() {
				queued += n
			}
		}
	}
	paperMax := paperMax(media)
	cmd := "PDF,PWG,URF,JPEG,PNG"
	scheme := "http"
	if q.secure {
		scheme = "https"
	}
	icons := []ipp.Value{}
	for _, n := range iconSizes {
		icons = append(icons, ipp.URI(fmt.Sprintf("%s://%s/icons/%d.png", scheme, host, n)))
	}
	pdfVersions := []ipp.Value{ipp.Keyword("adobe-1.3"), ipp.Keyword("adobe-1.4"), ipp.Keyword("adobe-1.5"), ipp.Keyword("adobe-1.6"),
		ipp.Keyword("adobe-1.7"), ipp.Keyword("iso-32000-1_2008"), ipp.Keyword("iso-32000-2_2017"), ipp.Keyword("pwg-5102.3")}
	attrs := []ipp.Attribute{
		ipp.Attr("printer-uri-supported", uris...),
		ipp.Attr("uri-authentication-supported", auth...),
		ipp.Attr("uri-security-supported", security...),
		ipp.Attr("printer-name", ipp.Name(q.queue.name)),
		ipp.Attr("printer-info", ipp.Text(q.queue.info)),
		ipp.Attr("printer-location", ipp.Text(q.queue.location)),
		ipp.Attr("printer-make-and-model", ipp.Text(clip127(makeModel))),
		ipp.Attr("printer-more-info", ipp.URI(q.moreInfo())),
		ipp.Attr("printer-supply-info-uri", ipp.URI(q.moreInfo())),
		ipp.Attr("printer-device-id", ipp.Text("MFG:BDF;MDL:Print Server;CMD:"+cmd+";")),
		ipp.Attr("printer-uuid", ipp.URI(q.queue.uuid)),
		ipp.Attr("printer-dns-sd-name", ipp.Name(x.instanceName(q.queue))),
		ipp.Attr("printer-kind", ipp.Keyword("document")),
		ipp.Attr("printer-organization", ipp.Text(x.cfg.Organization)),
		ipp.Attr("printer-organizational-unit", ipp.Text(x.cfg.OrganizationalUnit)),
		ipp.Attr("printer-geo-location", ipp.OutOfBand(ipp.TagUnknown)),
		ipp.Attr("printer-icons", icons...),
		ipp.Attr("printer-state", ipp.Enum(int(st))),
		ipp.Attr("printer-state-reasons", keywordValues(reasons)...),
		ipp.Attr("printer-state-message", ipp.Text(clip127(msg))),
		ipp.Attr("printer-is-accepting-jobs", ipp.Boolean(accepting)),
		ipp.Attr("printer-up-time", ipp.Integer(x.upTime(time.Now()))),
		ipp.Attr("printer-current-time", ipp.DateTime(time.Now())),
		ipp.Attr("printer-state-change-time", ipp.Integer(x.upTime(changed))),
		ipp.Attr("printer-state-change-date-time", ipp.DateTime(changed)),
		ipp.Attr("printer-config-change-time", ipp.Integer(1)),
		ipp.Attr("printer-config-change-date-time", ipp.DateTime(x.start)),
		ipp.Attr("queued-job-count", ipp.Integer(queued)),
		ipp.Attr("ipp-versions-supported", ipp.Keyword("1.1"), ipp.Keyword("2.0")),
		ipp.Attr("ipp-features-supported", features...),
		ipp.Attr("operations-supported", ops...),
		ipp.Attr("charset-configured", ipp.Charset("utf-8")),
		ipp.Attr("charset-supported", ipp.Charset("utf-8"), ipp.Charset("us-ascii")),
		ipp.Attr("natural-language-configured", ipp.Language("en")),
		ipp.Attr("generated-natural-language-supported", ipp.Language("en")),
		ipp.Attr("compression-supported", ipp.Keyword("none"), ipp.Keyword("deflate"), ipp.Keyword("gzip")),
		ipp.Attr("document-format-supported", mimeValues(documentFormats)...),
		ipp.Attr("document-format-default", ipp.MimeType("application/octet-stream")),
		ipp.Attr("document-format-preferred", ipp.MimeType("application/pdf")),
		ipp.Attr("pdf-versions-supported", pdfVersions...),
		ipp.Attr("pdl-override-supported", ipp.Keyword("attempted")),
		ipp.Attr("multiple-document-jobs-supported", ipp.Boolean(false)),
		ipp.Attr("multiple-operation-time-out", ipp.Integer(int(incomingTimeout/time.Second))),
		ipp.Attr("multiple-operation-time-out-action", ipp.Keyword("abort-job")),
		ipp.Attr("job-ids-supported", ipp.Boolean(true)),
		ipp.Attr("preferred-attributes-supported", ipp.Boolean(false)),
		ipp.Attr("which-jobs-supported", which...),
		ipp.Attr("job-creation-attributes-supported", keywordValues(jobCreationAttrs)...),
		ipp.Attr("printer-get-attributes-supported", ipp.Keyword("document-format")),
		ipp.Attr("overrides-supported", ipp.Keyword("document-number"), ipp.Keyword("pages")),
		ipp.Attr("identify-actions-supported", ipp.Keyword("display")),
		ipp.Attr("identify-actions-default", ipp.Keyword("display")),
		ipp.Attr("job-privacy-attributes", ipp.Keyword("default")),
		ipp.Attr("job-privacy-scope", ipp.Keyword("owner")),
		ipp.Attr("copies-supported", ipp.NewRange(1, 999)),
		ipp.Attr("copies-default", ipp.Integer(1)),
		ipp.Attr("page-ranges-supported", ipp.Boolean(true)),
		ipp.Attr("job-hold-until-supported", keywordValues(holds)...),
		ipp.Attr("job-hold-until-default", ipp.Keyword("no-hold")),
		ipp.Attr("job-priority-supported", ipp.Integer(100)),
		ipp.Attr("job-priority-default", ipp.Integer(50)),
		ipp.Attr("job-sheets-supported", ipp.Keyword("none")),
		ipp.Attr("job-sheets-default", ipp.Keyword("none")),
		ipp.Attr("multiple-document-handling-supported", keywordValues(handlings)...),
		ipp.Attr("multiple-document-handling-default", ipp.Keyword(handlings[0])),
		ipp.Attr("number-up-supported", ipp.Integer(1)),
		ipp.Attr("number-up-default", ipp.Integer(1)),
		ipp.Attr("orientation-requested-supported", ipp.Enum(3), ipp.Enum(4), ipp.Enum(5), ipp.Enum(6)),
		ipp.Attr("orientation-requested-default", ipp.OutOfBand(ipp.TagNoValue)),
		ipp.Attr("print-quality-supported", ipp.Enum(3), ipp.Enum(4), ipp.Enum(5)),
		ipp.Attr("print-quality-default", ipp.Enum(4)),
		ipp.Attr("printer-resolution-supported", resolutionValues(rasterResolutions)...),
		ipp.Attr("printer-resolution-default", ipp.NewResolution(300, 300, ipp.UnitsDPI)),
		ipp.Attr("print-scaling-supported", keywordValues(scalings)...),
		ipp.Attr("print-scaling-default", ipp.Keyword("auto")),
		ipp.Attr("print-color-mode-supported", keywordValues(modes)...),
		ipp.Attr("print-color-mode-default", ipp.Keyword("auto")),
		ipp.Attr("color-supported", ipp.Boolean(color)),
		ipp.Attr("print-content-optimize-supported", keywordValues(contentOptimizes)...),
		ipp.Attr("print-content-optimize-default", ipp.Keyword("auto")),
		ipp.Attr("print-rendering-intent-supported", keywordValues(renderingIntents)...),
		ipp.Attr("print-rendering-intent-default", ipp.Keyword("auto")),
		ipp.Attr("sides-supported", keywordValues(sides)...),
		ipp.Attr("sides-default", ipp.Keyword("one-sided")),
		ipp.Attr("finishings-supported", enumValues(finishings)...),
		ipp.Attr("finishings-default", ipp.Enum(3)),
		ipp.Attr("media-supported", keywordValues(media)...),
		ipp.Attr("media-default", ipp.Keyword(ready[0])),
		ipp.Attr("media-ready", keywordValues(ready)...),
		ipp.Attr("media-col-supported", ipp.Keyword("media-bottom-margin"), ipp.Keyword("media-left-margin"), ipp.Keyword("media-right-margin"),
			ipp.Keyword("media-size"), ipp.Keyword("media-size-name"), ipp.Keyword("media-source"), ipp.Keyword("media-top-margin"), ipp.Keyword("media-type")),
		ipp.Attr("media-col-database", mediaDB...),
		ipp.Attr("media-col-ready", mediaReady...),
		ipp.Attr("media-col-default", mediaReady[0]),
		ipp.Attr("media-size-supported", mediaSizes...),
		ipp.Attr("media-source-supported", keywordValues(sources)...),
		ipp.Attr("media-type-supported", keywordValues(types)...),
		ipp.Attr("media-top-margin-supported", ipp.Integer(top)),
		ipp.Attr("media-bottom-margin-supported", ipp.Integer(bottom)),
		ipp.Attr("media-left-margin-supported", ipp.Integer(left)),
		ipp.Attr("media-right-margin-supported", ipp.Integer(right)),
		ipp.Attr("pwg-raster-document-resolution-supported", resolutionValues(rasterResolutions)...),
		ipp.Attr("pwg-raster-document-type-supported", keywordValues(rasterTypes)...),
		ipp.Attr("pwg-raster-document-sheet-back", ipp.Keyword("normal")),
		ipp.Attr("urf-supported", keywordValues(urfSupported)...),
		ipp.Attr("printer-supply", ipp.OctetString([]byte("index=1;class=unknown;type=other;unit=percent;maxcapacity=100;level=-2;colorantname=unknown;"))),
		ipp.Attr("printer-supply-description", ipp.Text("supplies of the printers behind the server")),
		ipp.Attr("job-k-octets-supported", ipp.NewRange(0, int(min(x.s.maxDocument()/1024, math.MaxInt32)))),
		ipp.Attr("document-password-supported", ipp.Integer(1023)),
		ipp.Attr("printer-firmware-name", ipp.Name("bdf")),
		ipp.Attr("printer-firmware-string-version", ipp.Text("bdf printd")),
		ipp.Attr("printer-firmware-version", ipp.OctetString([]byte{1})),
		ipp.Attr("landscape-orientation-requested-preferred", ipp.Enum(5)),
		ipp.Attr("mopria-certified", ipp.Text("2.0")),
	}
	ppm := 20
	if caps != nil {
		if a := attrOf(caps.Attrs, "pages-per-minute"); a != nil {
			if n, ok := a.Int(); ok && n > 0 {
				ppm = int(n)
			}
		}
	}
	attrs = append(attrs, ipp.Attr("pages-per-minute", ipp.Integer(ppm)))
	if color {
		attrs = append(attrs, ipp.Attr("pages-per-minute-color", ipp.Integer(ppm)))
	}
	if paperMax != "" {
		attrs = append(attrs, ipp.Attr("printer-paper-max", ipp.Keyword(paperMax)))
	}
	if caps != nil && len(caps.OutputBins) > 0 {
		attrs = append(attrs, ipp.Attr("output-bin-supported", keywordValues(caps.OutputBins)...),
			ipp.Attr("output-bin-default", ipp.Keyword(cmp(caps.OutputBinDefault, caps.OutputBins[0]))))
	}
	return attrs
}

func cmp(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func mimeValues(s []string) []ipp.Value {
	out := make([]ipp.Value, len(s))
	for i, v := range s {
		out[i] = ipp.MimeType(v)
	}
	return out
}

// margins are the margins of the queue's printer in hundredths of a
// millimetre (its media-col-database's, or its -margin-supported), 4.23 mm
// (a sixth of an inch) when it does not say.
func margins(caps *ipp.Printer) (top, bottom, left, right int) {
	top, bottom, left, right = 423, 423, 423, 423
	if caps == nil {
		return
	}
	get := func(name string, v *int) {
		if a := attrOf(caps.Attrs, name); a != nil {
			if ns := a.Ints(); len(ns) > 0 {
				*v = int(slices.Max(ns))
			}
		}
	}
	get("media-top-margin-supported", &top)
	get("media-bottom-margin-supported", &bottom)
	get("media-left-margin-supported", &left)
	get("media-right-margin-supported", &right)
	return
}

func attrOf(attrs []ipp.Attribute, name string) *ipp.Attribute { return ipp.Get(attrs, name) }

// paperMax is the PaperMax of the TXT record (and printer-paper-max): the
// class of the largest of the media.
func paperMax(media []string) string {
	most := 0.0
	for _, m := range media {
		if pm, err := print.ParseMedia(m); err == nil {
			most = max(most, max(pm.Width, pm.Height)*min(pm.Width, pm.Height))
		}
	}
	a4 := 595.28 * 841.89
	switch {
	case most == 0:
		return ""
	case most <= a4*1.05:
		return "legal-A4"
	case most <= a4*2.1:
		return "tabloid-A3"
	case most <= a4*4.2:
		return "isoC-A2"
	}
	return ">isoC-A2"
}

// maxDocument is the size of a document the spool takes.
func (s *Server) maxDocument() int64 {
	if n := int64(s.cfg.Limits.MaxDocument); n > 0 {
		return n
	}
	return 256 << 20
}

var _ = strconv.Itoa
