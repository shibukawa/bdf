// Package ipptest is a printer that is not there: an http.Handler that
// answers IPP as an IPP Everywhere printer does, keeps the jobs and
// documents it is sent, and can be told what to support and what to
// refuse. The tests of the client, of bdf print and of the spooler run
// against it.
package ipptest

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
)

// Printer is the printer that is not there. New gives one with the
// capabilities of a plain office printer; change the fields before
// serving.
type Printer struct {
	Name      string
	MakeModel string
	UUID      string
	// Formats are document-format-supported.
	Formats []string
	// Reject are formats refused with client-error-document-format-
	// not-supported although listed, to test a client's fallback.
	Reject      []string
	Media       []string
	MediaReady  []string
	Sides       []string
	ColorModes  []string
	Qualities   []int
	Resolutions []ipp.Resolution
	Copies      ipp.Range
	PageRanges  bool
	Finishings  []int
	Scaling     []string

	PWGRasterTypes       []string
	PWGRasterResolutions []ipp.Resolution
	SheetBack            string
	URF                  []string
	// PCLmStripHeight, when not 0, makes the printer say it takes PCLm,
	// with these strips, back side, resolutions and compressions.
	PCLmStripHeight int
	PCLmBackSide    string
	PCLmResolutions []ipp.Resolution
	PCLmCompression []string

	State     ipp.PrinterState
	Reasons   []string
	Accepting bool

	// Username and Password, when set, are required as Basic
	// authentication.
	Username, Password string
	// Only11 refuses IPP/2.x with server-error-version-not-supported.
	Only11 bool
	// Processing is how long a job stays processing after its last
	// document before it is completed.
	Processing time.Duration
	// Busy refuses so many of the next Create-Job and Print-Job with
	// server-error-busy, as a printer that is printing for others does.
	Busy int
	// Dir saves each document to a file "<job>-<n>.<ext>" when set.
	Dir string
	// Limits bound the decoding of requests.
	Limits *ipp.Limits

	mu      sync.Mutex
	jobs    []*Job
	nextJob int32
	// failAfter cuts the connection of the next document after so many
	// bytes (-1: no); dropCreate cuts it before the answer to the next
	// Create-Job, once the job is made.
	failAfter  int64
	dropCreate bool
	// Requests counts the operations received, by name.
	Requests map[string]int
}

// Job is a job the printer received.
type Job struct {
	ID        int32
	Name      string
	User      string
	Format    string // of the last document
	Attrs     []ipp.Attribute
	Documents [][]byte
	Created   time.Time
	// doneAt is when the last document arrived; the state is derived
	// from it and Processing.
	doneAt   time.Time
	canceled bool
	aborted  bool
}

// New returns a printer with the capabilities of a plain office
// printer: PDF, PWG Raster and JPEG, A4, Letter and B5, two-sided,
// colour, 300 and 600 dpi, up to 99 copies.
func New() *Printer {
	return &Printer{
		Name: "ipptest", MakeModel: "BDF Test Printer", UUID: "urn:uuid:4e7b5a0e-5c3e-4f1a-9f6b-000000000001",
		Formats:              []string{"application/pdf", "image/pwg-raster", "image/jpeg", "application/octet-stream"},
		Media:                []string{"iso_a4_210x297mm", "na_letter_8.5x11in", "jis_b5_182x257mm", "iso_a3_297x420mm"},
		MediaReady:           []string{"iso_a4_210x297mm"},
		Sides:                []string{"one-sided", "two-sided-long-edge", "two-sided-short-edge"},
		ColorModes:           []string{"auto", "color", "monochrome"},
		Qualities:            []int{3, 4, 5},
		Resolutions:          []ipp.Resolution{{X: 300, Y: 300, Units: ipp.UnitsDPI}, {X: 600, Y: 600, Units: ipp.UnitsDPI}},
		Copies:               ipp.Range{Lo: 1, Hi: 99},
		PageRanges:           true,
		Finishings:           []int{3},
		Scaling:              []string{"auto", "auto-fit", "fill", "fit", "none"},
		PWGRasterTypes:       []string{"black_1", "sgray_8", "srgb_8"},
		PWGRasterResolutions: []ipp.Resolution{{X: 300, Y: 300, Units: ipp.UnitsDPI}, {X: 600, Y: 600, Units: ipp.UnitsDPI}},
		SheetBack:            "normal",
		URF:                  []string{"V1.4", "W8", "SRGB24", "RS300-600", "DM1"},
		State:                ipp.PrinterIdle, Reasons: []string{"none"}, Accepting: true,
		Requests: map[string]int{}, failAfter: -1,
	}
}

// Stop stops the printer for reasons ("media-empty-error" …): it answers
// as stopped and does not finish its jobs until Resume.
func (p *Printer) Stop(reasons ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(reasons) == 0 {
		reasons = []string{"paused"}
	}
	p.State, p.Reasons = ipp.PrinterStopped, reasons
}

// Resume makes a stopped printer idle again.
func (p *Printer) Resume() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.State, p.Reasons = ipp.PrinterIdle, []string{"none"}
}

// Forget drops the jobs the printer knows, as a printer that keeps no
// history (or was restarted) does: their ids are not reused.
func (p *Printer) Forget() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.jobs = nil
}

// FailAfter cuts the connection of the next document sent after n bytes
// of it, as a network that fails does. The job keeps no document.
func (p *Printer) FailAfter(n int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.failAfter = n
}

// DropCreateResponse makes the next Create-Job make its job and then cut
// the connection before answering: the client does not learn the job.
func (p *Printer) DropCreateResponse() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.dropCreate = true
}

// Jobs returns the jobs received, oldest first.
func (p *Printer) Jobs() []*Job {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]*Job(nil), p.jobs...)
}

// Job returns a job by id, or nil.
func (p *Printer) Job(id int32) *Job {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.job(id)
}

func (p *Printer) job(id int32) *Job {
	for _, j := range p.jobs {
		if j.ID == id {
			return j
		}
	}
	return nil
}

// JobState returns the state of a job now.
func (p *Printer) JobState(j *Job) ipp.JobState {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.state(j)
}

func (p *Printer) state(j *Job) ipp.JobState {
	switch {
	case j.canceled:
		return ipp.JobCanceled
	case j.aborted:
		return ipp.JobAborted
	case j.doneAt.IsZero():
		return ipp.JobPending
	case time.Since(j.doneAt) < p.Processing, p.State == ipp.PrinterStopped:
		return ipp.JobProcessing
	}
	return ipp.JobCompleted
}

// Abort makes a job aborted, to test a client's handling.
func (p *Printer) Abort(id int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if j := p.job(id); j != nil {
		j.aborted = true
	}
}

func (p *Printer) uri(r *http.Request) string {
	scheme := "ipp"
	if r.TLS != nil {
		scheme = "ipps"
	}
	return scheme + "://" + r.Host + "/ipp/print"
}

// ServeHTTP answers one request.
func (p *Printer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "IPP is POST", http.StatusMethodNotAllowed)
		return
	}
	if p.Username != "" {
		u, pw, ok := r.BasicAuth()
		if !ok || u != p.Username || pw != p.Password {
			w.Header().Set("WWW-Authenticate", `Basic realm="ipptest"`)
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
	}
	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/ipp") {
		http.Error(w, "not application/ipp", http.StatusBadRequest)
		return
	}
	req, br, err := ipp.DecodeFrom(bufio.NewReader(r.Body), p.Limits)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := p.handle(req, br, r)
	b, err := res.Marshal()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/ipp")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}

func (p *Printer) handle(req *ipp.Message, body *bufio.Reader, r *http.Request) *ipp.Message {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Requests[req.Op().String()]++
	res := ipp.NewResponse(ipp.StatusOK, req.RequestID, "utf-8", "en")
	res.Version = req.Version
	fail := func(st ipp.Status, msg string) *ipp.Message {
		res.Code = uint16(st)
		if msg != "" {
			res.Add(ipp.TagOperation, ipp.Attr("status-message", ipp.Text(msg)))
		}
		return res
	}
	if p.Only11 && req.Version.Major >= 2 {
		res.Version = ipp.Version11
		return fail(ipp.StatusServerErrorVersionNotSupported, "IPP/1.1 only")
	}
	op := req.Group(ipp.TagOperation)
	user := op.Get("requesting-user-name").Text()
	switch req.Op() {
	case ipp.OpGetPrinterAttributes:
		res.Add(ipp.TagPrinter, p.attributes(r)...)
		return res
	case ipp.OpValidateJob:
		if st, msg, unsupported := p.check(req); st != ipp.StatusOK {
			fail(st, msg)
			if len(unsupported) > 0 {
				res.Add(ipp.TagUnsupportedGroup, unsupported...)
			}
		}
		return res
	case ipp.OpCreateJob, ipp.OpPrintJob:
		if !p.Accepting {
			return fail(ipp.StatusServerErrorNotAcceptingJobs, "not accepting jobs")
		}
		if p.Busy > 0 {
			p.Busy--
			return fail(ipp.StatusServerErrorBusy, "busy")
		}
		if st, msg, unsupported := p.check(req); st != ipp.StatusOK {
			fail(st, msg)
			if len(unsupported) > 0 {
				res.Add(ipp.TagUnsupportedGroup, unsupported...)
			}
			return res
		}
		p.nextJob++
		j := &Job{ID: p.nextJob, Name: op.Get("job-name").Text(), User: user, Created: time.Now()}
		if g := req.Group(ipp.TagJob); g != nil {
			j.Attrs = g.Attrs
		}
		p.jobs = append(p.jobs, j)
		if req.Op() == ipp.OpCreateJob && p.dropCreate {
			p.dropCreate = false
			panic(http.ErrAbortHandler)
		}
		if req.Op() == ipp.OpPrintJob {
			j.Format = op.Get("document-format").Text()
			if err := p.receive(j, body); err != nil {
				return fail(ipp.StatusClientErrorBadRequest, err.Error())
			}
			j.doneAt = time.Now()
		}
		res.Add(ipp.TagJob, p.jobAttrs(j, r)...)
		return res
	case ipp.OpSendDocument:
		id, _ := op.Get("job-id").Int()
		j := p.job(id)
		if j == nil {
			return fail(ipp.StatusClientErrorNotFound, "no such job")
		}
		if !j.doneAt.IsZero() {
			return fail(ipp.StatusClientErrorNotPossible, "the job has its last document")
		}
		f := op.Get("document-format").Text()
		if f != "" && !p.supports(f) {
			return fail(ipp.StatusClientErrorDocumentFormatNotSupported, f)
		}
		j.Format = f
		if err := p.receive(j, body); err != nil {
			return fail(ipp.StatusClientErrorBadRequest, err.Error())
		}
		if last, _ := op.Get("last-document").Bool(); last {
			j.doneAt = time.Now()
		}
		res.Add(ipp.TagJob, p.jobAttrs(j, r)...)
		return res
	case ipp.OpGetJobAttributes:
		id, _ := op.Get("job-id").Int()
		j := p.job(id)
		if j == nil {
			return fail(ipp.StatusClientErrorNotFound, "no such job")
		}
		res.Add(ipp.TagJob, p.jobAttrs(j, r)...)
		return res
	case ipp.OpGetJobs:
		which := op.Get("which-jobs").Text()
		mine, _ := op.Get("my-jobs").Bool()
		for _, j := range p.jobs {
			done := p.state(j).Done()
			if (which == "completed" && !done) || ((which == "" || which == "not-completed") && done) {
				continue
			}
			if mine && j.User != user {
				continue
			}
			res.Groups = append(res.Groups, ipp.Group{Tag: ipp.TagJob, Attrs: p.jobAttrs(j, r)})
		}
		return res
	case ipp.OpCancelJob:
		id, _ := op.Get("job-id").Int()
		j := p.job(id)
		if j == nil {
			return fail(ipp.StatusClientErrorNotFound, "no such job")
		}
		if p.state(j).Done() {
			return fail(ipp.StatusClientErrorNotPossible, "the job is done")
		}
		j.canceled = true
		return res
	case ipp.OpIdentifyPrinter:
		return res
	case ipp.OpCupsGetPrinters:
		res.Add(ipp.TagPrinter, ipp.Attr("printer-name", ipp.Name(p.Name)), ipp.Attr("printer-uri-supported", ipp.URI(p.uri(r))),
			ipp.Attr("printer-info", ipp.Text(p.Name)), ipp.Attr("printer-make-and-model", ipp.Text(p.MakeModel)),
			ipp.Attr("printer-state", ipp.Enum(int(p.State))), ipp.Attr("printer-state-reasons", keywords(p.Reasons)...),
			ipp.Attr("printer-is-accepting-jobs", ipp.Boolean(p.Accepting)), ipp.Attr("printer-type", ipp.Integer(0)))
		return res
	}
	return fail(ipp.StatusServerErrorOperationNotSupported, req.Op().String())
}

func (p *Printer) supports(format string) bool {
	for _, f := range p.Reject {
		if f == format {
			return false
		}
	}
	for _, f := range p.Formats {
		if f == format {
			return true
		}
	}
	return false
}

// check validates the document format and the job attributes of a
// request.
func (p *Printer) check(req *ipp.Message) (ipp.Status, string, []ipp.Attribute) {
	op := req.Group(ipp.TagOperation)
	if f := op.Get("document-format").Text(); f != "" && !p.supports(f) {
		return ipp.StatusClientErrorDocumentFormatNotSupported, f, nil
	}
	var unsupported []ipp.Attribute
	g := req.Group(ipp.TagJob)
	if g == nil {
		return ipp.StatusOK, "", nil
	}
	for _, a := range g.Attrs {
		bad := false
		switch a.Name {
		case "copies":
			n, _ := a.Int()
			bad = n < p.Copies.Lo || n > p.Copies.Hi
		case "sides":
			bad = !has(p.Sides, a.Text())
		case "media":
			bad = !has(p.Media, a.Text())
		case "print-color-mode":
			bad = !has(p.ColorModes, a.Text())
		case "print-scaling":
			bad = !has(p.Scaling, a.Text())
		case "page-ranges":
			bad = !p.PageRanges
		case "print-quality":
			n, _ := a.Int()
			bad = !hasInt(p.Qualities, int(n))
		case "printer-resolution":
			rs := a.Resolutions()
			bad = len(rs) == 0 || !hasRes(p.Resolutions, rs[0])
		case "finishings":
			for _, n := range a.Ints() {
				if !hasInt(p.Finishings, int(n)) {
					bad = true
				}
			}
		case "media-col", "job-name", "multiple-document-handling", "orientation-requested", "output-bin", "job-hold-until", "job-priority", "number-up":
		default:
			bad = true
		}
		if bad {
			unsupported = append(unsupported, a)
		}
	}
	if len(unsupported) > 0 {
		return ipp.StatusClientErrorAttributesOrValuesNotSupported, "", unsupported
	}
	return ipp.StatusOK, "", nil
}

func (p *Printer) receive(j *Job, body *bufio.Reader) error {
	if p.failAfter >= 0 {
		io.CopyN(io.Discard, body, p.failAfter)
		p.failAfter = -1
		panic(http.ErrAbortHandler)
	}
	b, err := io.ReadAll(io.LimitReader(body, 256<<20))
	if err != nil {
		return err
	}
	j.Documents = append(j.Documents, b)
	if p.Dir != "" {
		ext := "bin"
		switch j.Format {
		case "application/pdf":
			ext = "pdf"
		case "image/pwg-raster":
			ext = "pwg"
		case "image/urf":
			ext = "urf"
		case "application/postscript":
			ext = "ps"
		case "image/jpeg":
			ext = "jpg"
		}
		name := filepath.Join(p.Dir, fmt.Sprintf("%d-%d.%s", j.ID, len(j.Documents), ext))
		if err := os.WriteFile(name, b, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func (p *Printer) jobAttrs(j *Job, r *http.Request) []ipp.Attribute {
	st := p.state(j)
	reason := "none"
	switch st {
	case ipp.JobPending:
		reason = "job-incoming"
	case ipp.JobProcessing:
		reason = "job-printing"
	case ipp.JobCompleted:
		reason = "job-completed-successfully"
	case ipp.JobCanceled:
		reason = "job-canceled-by-user"
	case ipp.JobAborted:
		reason = "aborted-by-system"
	}
	copies := int32(1)
	for _, a := range j.Attrs {
		if a.Name == "copies" {
			copies, _ = a.Int()
		}
	}
	impressions := 0
	if st == ipp.JobCompleted {
		impressions = len(j.Documents) * int(copies)
	}
	return []ipp.Attribute{
		ipp.Attr("job-id", ipp.Integer(int(j.ID))),
		ipp.Attr("job-uri", ipp.URI(p.uri(r)+"/job/"+strconv.Itoa(int(j.ID)))),
		ipp.Attr("job-name", ipp.Name(j.Name)),
		ipp.Attr("job-originating-user-name", ipp.Name(j.User)),
		ipp.Attr("job-state", ipp.Enum(int(st))),
		ipp.Attr("job-state-reasons", ipp.Keyword(reason)),
		ipp.Attr("job-impressions-completed", ipp.Integer(impressions)),
	}
}

func (p *Printer) attributes(r *http.Request) []ipp.Attribute {
	uri := p.uri(r)
	auth := "none"
	if p.Username != "" {
		auth = "basic"
	}
	sec := "none"
	if r.TLS != nil {
		sec = "tls"
	}
	attrs := []ipp.Attribute{
		ipp.Attr("printer-uri-supported", ipp.URI(uri)),
		ipp.Attr("uri-authentication-supported", ipp.Keyword(auth)),
		ipp.Attr("uri-security-supported", ipp.Keyword(sec)),
		ipp.Attr("printer-name", ipp.Name(p.Name)),
		ipp.Attr("printer-info", ipp.Text(p.Name)),
		ipp.Attr("printer-make-and-model", ipp.Text(p.MakeModel)),
		ipp.Attr("printer-uuid", ipp.URI(p.UUID)),
		ipp.Attr("printer-state", ipp.Enum(int(p.State))),
		ipp.Attr("printer-state-reasons", keywords(p.Reasons)...),
		ipp.Attr("printer-is-accepting-jobs", ipp.Boolean(p.Accepting)),
		ipp.Attr("ipp-versions-supported", ipp.Keyword("1.1"), ipp.Keyword("2.0")),
		ipp.Attr("operations-supported", ipp.Enum(int(ipp.OpPrintJob)), ipp.Enum(int(ipp.OpValidateJob)), ipp.Enum(int(ipp.OpCreateJob)),
			ipp.Enum(int(ipp.OpSendDocument)), ipp.Enum(int(ipp.OpCancelJob)), ipp.Enum(int(ipp.OpGetJobAttributes)),
			ipp.Enum(int(ipp.OpGetJobs)), ipp.Enum(int(ipp.OpGetPrinterAttributes)), ipp.Enum(int(ipp.OpIdentifyPrinter))),
		ipp.Attr("printer-config-change-time", ipp.Integer(1)),
		ipp.Attr("document-format-supported", mimes(p.Formats)...),
		ipp.Attr("document-format-default", ipp.MimeType(p.Formats[0])),
		ipp.Attr("media-supported", keywords(p.Media)...),
		ipp.Attr("media-ready", keywords(p.MediaReady)...),
		ipp.Attr("media-default", ipp.Keyword(p.Media[0])),
		ipp.Attr("media-source-supported", ipp.Keyword("auto"), ipp.Keyword("tray-1")),
		ipp.Attr("media-type-supported", ipp.Keyword("auto"), ipp.Keyword("stationery")),
		ipp.Attr("media-col-supported", ipp.Keyword("media-size"), ipp.Keyword("media-source"), ipp.Keyword("media-type")),
		ipp.Attr("sides-supported", keywords(p.Sides)...),
		ipp.Attr("sides-default", ipp.Keyword(p.Sides[0])),
		ipp.Attr("print-color-mode-supported", keywords(p.ColorModes)...),
		ipp.Attr("print-color-mode-default", ipp.Keyword(p.ColorModes[0])),
		ipp.Attr("print-quality-supported", enums(p.Qualities)...),
		ipp.Attr("print-quality-default", ipp.Enum(4)),
		ipp.Attr("printer-resolution-supported", resolutions(p.Resolutions)...),
		ipp.Attr("printer-resolution-default", resolutions(p.Resolutions[:1])...),
		ipp.Attr("copies-supported", ipp.NewRange(int(p.Copies.Lo), int(p.Copies.Hi))),
		ipp.Attr("copies-default", ipp.Integer(1)),
		ipp.Attr("page-ranges-supported", ipp.Boolean(p.PageRanges)),
		ipp.Attr("print-scaling-supported", keywords(p.Scaling)...),
		ipp.Attr("print-scaling-default", ipp.Keyword("auto")),
		ipp.Attr("finishings-supported", enums(p.Finishings)...),
		ipp.Attr("finishings-default", ipp.Enum(3)),
		ipp.Attr("output-bin-supported", ipp.Keyword("face-down")),
		ipp.Attr("pwg-raster-document-resolution-supported", resolutions(p.PWGRasterResolutions)...),
		ipp.Attr("pwg-raster-document-type-supported", keywords(p.PWGRasterTypes)...),
		ipp.Attr("pwg-raster-document-sheet-back", ipp.Keyword(p.SheetBack)),
		ipp.Attr("urf-supported", keywords(p.URF)...),
		ipp.Attr("identify-actions-supported", ipp.Keyword("display"), ipp.Keyword("sound")),
		ipp.Attr("multiple-document-handling-supported", ipp.Keyword("separate-documents-collated-copies"), ipp.Keyword("separate-documents-uncollated-copies")),
	}
	if p.PCLmStripHeight > 0 {
		attrs = append(attrs, ipp.Attr("pclm-strip-height-preferred", ipp.Integer(p.PCLmStripHeight)),
			ipp.Attr("pclm-raster-back-side", ipp.Keyword(p.PCLmBackSide)),
			ipp.Attr("pclm-source-resolution-supported", resolutions(p.PCLmResolutions)...),
			ipp.Attr("pclm-compression-method-preferred", keywords(p.PCLmCompression)...))
	}
	return attrs
}

func keywords(s []string) []ipp.Value {
	out := make([]ipp.Value, len(s))
	for i, v := range s {
		out[i] = ipp.Keyword(v)
	}
	return out
}

func mimes(s []string) []ipp.Value {
	out := make([]ipp.Value, len(s))
	for i, v := range s {
		out[i] = ipp.MimeType(v)
	}
	return out
}

func enums(ns []int) []ipp.Value {
	out := make([]ipp.Value, len(ns))
	for i, n := range ns {
		out[i] = ipp.Enum(n)
	}
	return out
}

func resolutions(rs []ipp.Resolution) []ipp.Value {
	out := make([]ipp.Value, len(rs))
	for i, r := range rs {
		out[i] = ipp.NewResolution(int(r.X), int(r.Y), r.Units)
	}
	return out
}

func has(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func hasInt(list []int, n int) bool {
	for _, v := range list {
		if v == n {
			return true
		}
	}
	return false
}

func hasRes(list []ipp.Resolution, r ipp.Resolution) bool {
	for _, v := range list {
		if v == r {
			return true
		}
	}
	return false
}

// Document returns the bytes of a document of a job, for a test to
// look at.
func (j *Job) Document(n int) []byte {
	if n < 0 || n >= len(j.Documents) {
		return nil
	}
	return bytes.Clone(j.Documents[n])
}
