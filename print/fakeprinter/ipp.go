package fakeprinter

import (
	"bufio"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
)

// The operations a fake printer supports: those IPP Everywhere (PWG
// 5100.14) requires, Hold-Job and Release-Job, and CUPS-Get-Printers
// (answered with the printer itself).
var operations = []ipp.Op{
	ipp.OpPrintJob, ipp.OpValidateJob, ipp.OpCreateJob, ipp.OpSendDocument, ipp.OpCancelJob, ipp.OpGetJobAttributes,
	ipp.OpGetJobs, ipp.OpGetPrinterAttributes, ipp.OpHoldJob, ipp.OpReleaseJob, ipp.OpCancelMyJobs, ipp.OpCloseJob,
	ipp.OpIdentifyPrinter,
}

// allAttrs is what a request without requested-attributes asks for.
var allAttrs = []string{"all"}

// jobGroups are the groups of job attributes requested-attributes names:
// "job-template" the job template attributes, "job-description" the
// others (filter lets every name not in a group through as itself).
var jobGroups = map[string][]string{
	"job-template": ipp.JobTemplateAttrs,
	"job-description": {"job-id", "job-uri", "job-uuid", "job-printer-uri", "job-name", "job-originating-user-name", "job-state",
		"job-state-reasons", "job-state-message", "number-of-documents", "job-printer-up-time", "job-k-octets",
		"job-impressions-completed", "job-media-sheets-completed", "time-at-creation", "time-at-processing", "time-at-completed",
		"date-time-at-creation", "date-time-at-processing", "date-time-at-completed", "document-format-supplied"},
}

// request is an IPP request as it is served.
type request struct {
	r      *http.Request
	req    *ipp.Message
	op     *ipp.Group
	body   *bufio.Reader
	res    *ipp.Message
	path   string // of the printer, as the client reached it
	secure bool
}

func (q *request) fail(st ipp.Status, format string, args ...any) *ipp.Message {
	q.res.Code = uint16(st)
	if format != "" {
		q.res.Add(ipp.TagOperation, ipp.Attr("status-message", ipp.Text(fmt.Sprintf(format, args...))))
	}
	return q.res
}

// printerURI is the printer's URI as the client reached it.
func (q *request) printerURI() string {
	scheme := "ipp"
	if q.secure {
		scheme = "ipps"
	}
	return scheme + "://" + q.r.Host + q.path
}

func (q *request) jobURI(id int32) string { return q.printerURI() + "/" + strconv.Itoa(int(id)) }

// ServeHTTP answers an IPP request at any path.
func (p *Printer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.serveIPP(w, r, "/ipp/print")
}

// serveIPP answers an IPP request to the printer at path.
func (p *Printer) serveIPP(w http.ResponseWriter, r *http.Request, path string) {
	if !p.admit(w, r) {
		return
	}
	q, ok := decode(w, r, p.throttle(r.Body), p.limits(), path)
	if !ok {
		return
	}
	p.handle(q)
	q.write(w)
}

// admit makes the printer's latency and asks for its Basic
// authentication; it writes the refusal itself.
func (p *Printer) admit(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		http.Error(w, "IPP is POST", http.StatusMethodNotAllowed)
		return false
	}
	p.lock()
	user, pass, latency, noChunked := p.Username, p.Password, p.Latency, p.NoChunked
	p.mu.Unlock()
	if latency > 0 {
		time.Sleep(latency)
	}
	if noChunked && r.ContentLength < 0 {
		http.Error(w, "a request needs its length", http.StatusLengthRequired)
		return false
	}
	if user != "" {
		u, pw, ok := r.BasicAuth()
		if !ok || u != user || pw != pass {
			w.Header().Set("WWW-Authenticate", `Basic realm="fakeprinter"`)
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return false
		}
	}
	return true
}

func (p *Printer) limits() *ipp.Limits {
	p.lock()
	defer p.mu.Unlock()
	return p.Limits
}

// decode reads an IPP request; it writes the refusal of one that is not.
func decode(w http.ResponseWriter, r *http.Request, body io.Reader, limits *ipp.Limits, path string) (*request, bool) {
	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/ipp") {
		http.Error(w, "not application/ipp", http.StatusBadRequest)
		return nil, false
	}
	if limits == nil {
		limits = &ipp.Limits{MaxBytes: 1 << 20, MaxAttributes: 1000, MaxValues: 1000, MaxValueBytes: 64 << 10, MaxDepth: 4}
	}
	req, br, err := ipp.DecodeFrom(bufio.NewReaderSize(body, 64<<10), limits)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, false
	}
	q := &request{r: r, req: req, body: br, path: path, secure: r.TLS != nil}
	q.op = req.Group(ipp.TagOperation)
	if q.op == nil {
		q.op = &ipp.Group{Tag: ipp.TagOperation}
	}
	q.res = ipp.NewResponse(ipp.StatusOK, req.RequestID, "utf-8", "en")
	q.res.Version = req.Version
	return q, true
}

// write writes the response of a request.
func (q *request) write(w http.ResponseWriter) {
	b, err := q.res.Marshal()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/ipp")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}

// throttle slows a body to the printer's bandwidth.
func (p *Printer) throttle(r io.Reader) io.Reader {
	p.lock()
	bw := p.Bandwidth
	p.mu.Unlock()
	if bw <= 0 {
		return r
	}
	return &slowReader{r: r, rate: bw}
}

// slowReader reads at most rate bytes a second.
type slowReader struct {
	r     io.Reader
	rate  int
	start time.Time
	n     int64
}

func (s *slowReader) Read(b []byte) (int, error) {
	if s.start.IsZero() {
		s.start = time.Now()
	}
	if len(b) > s.rate/10+1 {
		b = b[:s.rate/10+1]
	}
	n, err := s.r.Read(b)
	s.n += int64(n)
	if due := s.start.Add(time.Duration(float64(s.n) / float64(s.rate) * float64(time.Second))); time.Until(due) > 0 {
		time.Sleep(time.Until(due))
	}
	return n, err
}

// check checks what every request has (RFC 8011 §4.1): a version the
// printer speaks, a request id, the charset then the language first, no
// attribute twice in a group, and the target.
func (p *Printer) check(q *request) bool { return check(q, p.Only11, true) }

// check checks a request; target asks for its printer-uri or job-uri
// (CUPS's own operations have none).
func check(q *request, only11, target bool) bool {
	v := q.req.Version
	if v.Major < 1 || v.Major > 2 || only11 && v.Major >= 2 {
		q.res.Version = ipp.Version11
		if !only11 {
			q.res.Version = ipp.Version20
		}
		q.fail(ipp.StatusServerErrorVersionNotSupported, "IPP/%s is not supported", v)
		return false
	}
	if q.req.RequestID == 0 {
		q.fail(ipp.StatusClientErrorBadRequest, "request-id 0")
		return false
	}
	attrs := q.op.Attrs
	if len(attrs) < 2 || attrs[0].Name != "attributes-charset" || attrs[1].Name != "attributes-natural-language" {
		q.fail(ipp.StatusClientErrorBadRequest, "attributes-charset and attributes-natural-language come first")
		return false
	}
	if v, ok := attrs[0].First(); !ok || v.Tag != ipp.TagCharset {
		q.fail(ipp.StatusClientErrorBadRequest, "attributes-charset is a charset")
		return false
	}
	if cs := strings.ToLower(attrs[0].Text()); cs != "utf-8" && cs != "us-ascii" {
		q.fail(ipp.StatusClientErrorCharsetNotSupported, "charset %s: utf-8", cs)
		return false
	}
	if v, ok := attrs[1].First(); !ok || v.Tag != ipp.TagLanguage {
		q.fail(ipp.StatusClientErrorBadRequest, "attributes-natural-language is a naturalLanguage")
		return false
	}
	for _, g := range q.req.Groups {
		seen := map[string]bool{}
		for _, a := range g.Attrs {
			if seen[a.Name] {
				q.fail(ipp.StatusClientErrorBadRequest, "%s twice", a.Name)
				return false
			}
			seen[a.Name] = true
		}
	}
	if target && q.op.Get("printer-uri") == nil && q.op.Get("job-uri") == nil {
		q.fail(ipp.StatusClientErrorBadRequest, "no printer-uri")
		return false
	}
	return true
}

func (p *Printer) handle(q *request) {
	p.lock()
	defer p.mu.Unlock()
	p.Requests[q.req.Op().String()]++
	if !p.check(q) {
		return
	}
	switch q.req.Op() {
	case ipp.OpGetPrinterAttributes:
		p.getPrinterAttributes(q)
	case ipp.OpValidateJob:
		p.validate(q)
	case ipp.OpCreateJob, ipp.OpPrintJob:
		p.createJob(q)
	case ipp.OpSendDocument:
		p.sendDocument(q)
	case ipp.OpCloseJob:
		j := p.jobOf(q)
		if j == nil {
			return
		}
		if j.doneAt.IsZero() && !j.canceled && !j.aborted {
			p.done(j)
		}
	case ipp.OpGetJobAttributes:
		if j := p.jobOf(q); j != nil {
			q.res.Add(ipp.TagJob, ipp.Filter(p.jobAttrs(j, q), ipp.Requested(q.op, allAttrs), jobGroups)...)
		}
	case ipp.OpGetJobs:
		p.getJobs(q)
	case ipp.OpCancelJob:
		j := p.jobOf(q)
		if j == nil {
			return
		}
		if p.state(j).Done() {
			q.fail(ipp.StatusClientErrorNotPossible, "the job is done")
			return
		}
		j.canceled = true
		p.notify(Event{Kind: "job", Job: j.ID})
	case ipp.OpCancelMyJobs:
		user := q.op.Get("requesting-user-name").Text()
		for _, j := range p.jobs {
			if j.User == user && !p.state(j).Done() {
				j.canceled = true
			}
		}
		p.notify(Event{Kind: "jobs"})
	case ipp.OpHoldJob, ipp.OpReleaseJob:
		j := p.jobOf(q)
		if j == nil {
			return
		}
		if st := p.state(j); st != ipp.JobPending && st != ipp.JobPendingHeld {
			q.fail(ipp.StatusClientErrorNotPossible, "the job is %s", st)
			return
		}
		j.held = q.req.Op() == ipp.OpHoldJob
		p.notify(Event{Kind: "job", Job: j.ID})
	case ipp.OpIdentifyPrinter:
	case ipp.OpCupsGetPrinters:
		q.res.Add(ipp.TagPrinter, ipp.Filter(p.cupsAttrs(q, q.printerURI()), ipp.Requested(q.op, allAttrs), nil)...)
	default:
		q.fail(ipp.StatusServerErrorOperationNotSupported, "%s", q.req.Op())
	}
}

// jobOf finds the job of a request (job-id or job-uri), or fails it.
func (p *Printer) jobOf(q *request) *Job {
	id, ok := q.op.Get("job-id").Int()
	if !ok {
		if u := q.op.Get("job-uri").Text(); u != "" {
			if i := strings.LastIndexByte(u, '/'); i >= 0 {
				n, _ := strconv.Atoi(u[i+1:])
				id, ok = int32(n), n > 0
			}
		}
	}
	if !ok {
		q.fail(ipp.StatusClientErrorBadRequest, "no job-id")
		return nil
	}
	j := p.job(id)
	if j == nil {
		q.fail(ipp.StatusClientErrorNotFound, "no job %d", id)
	}
	return j
}

func (p *Printer) supports(format string) bool {
	return slices.Contains(p.Formats, format) && !slices.Contains(p.Reject, format)
}

// validate checks a job's document format and attributes (Validate-Job,
// and before Create-Job and Print-Job).
func (p *Printer) validate(q *request) bool {
	if f := q.op.Get("document-format").Text(); f != "" && !p.supports(f) {
		q.fail(ipp.StatusClientErrorDocumentFormatNotSupported, "%s", f)
		q.res.Add(ipp.TagUnsupportedGroup, *q.op.Get("document-format"))
		return false
	}
	if c := q.op.Get("compression").Text(); c != "" && c != "none" && c != "gzip" && c != "deflate" {
		q.fail(ipp.StatusClientErrorCompressionNotSupported, "compression %s", c)
		q.res.Add(ipp.TagUnsupportedGroup, *q.op.Get("compression"))
		return false
	}
	g := q.req.Group(ipp.TagJob)
	if g == nil {
		return true
	}
	var unsupported []ipp.Attribute
	for _, a := range g.Attrs {
		if !p.takes(a) {
			unsupported = append(unsupported, a)
		}
	}
	if len(unsupported) == 0 {
		return true
	}
	if fidelity, _ := q.op.Get("ipp-attribute-fidelity").Bool(); fidelity || !p.IgnoreUnsupported {
		q.fail(ipp.StatusClientErrorAttributesOrValuesNotSupported, "")
		q.res.Add(ipp.TagUnsupportedGroup, unsupported...)
		return false
	}
	q.res.Code = uint16(ipp.StatusOKIgnoredOrSubstituted)
	q.res.Add(ipp.TagUnsupportedGroup, unsupported...)
	return true
}

// takes tells whether the printer takes a job template attribute.
func (p *Printer) takes(a ipp.Attribute) bool {
	switch a.Name {
	case "copies":
		n, _ := a.Int()
		return n >= p.Copies.Lo && n <= p.Copies.Hi
	case "sides":
		return slices.Contains(p.Sides, a.Text())
	case "media":
		return slices.Contains(p.Media, a.Text())
	case "media-col":
		return p.takesMediaCol(a)
	case "print-color-mode":
		return slices.Contains(p.ColorModes, a.Text())
	case "print-scaling":
		return slices.Contains(p.Scaling, a.Text())
	case "page-ranges":
		return p.PageRanges
	case "print-quality":
		n, _ := a.Int()
		return slices.Contains(p.Qualities, int(n))
	case "printer-resolution":
		rs := a.Resolutions()
		return len(rs) == 1 && slices.Contains(p.Resolutions, rs[0])
	case "finishings":
		for _, n := range a.Ints() {
			if !slices.Contains(p.Finishings, int(n)) {
				return false
			}
		}
		return len(a.Values) > 0
	case "output-bin":
		return slices.Contains(p.OutputBins, a.Text())
	case "job-name", "multiple-document-handling", "orientation-requested", "job-hold-until", "job-priority", "number-up",
		"print-content-optimize", "print-rendering-intent", "overrides":
		return true
	}
	return false
}

// takesMediaCol tells whether the printer takes a media-col: its source
// and type among the printer's, its size one of the media's.
func (p *Printer) takesMediaCol(a ipp.Attribute) bool {
	v, ok := a.First()
	if !ok || v.Tag != ipp.TagBeginCollection {
		return false
	}
	if m := ipp.Get(v.Coll, "media-source"); m != nil && !slices.Contains(p.Sources, m.Text()) {
		return false
	}
	if m := ipp.Get(v.Coll, "media-type"); m != nil && !slices.Contains(p.Types, m.Text()) {
		return false
	}
	if m := ipp.Get(v.Coll, "media-size-name"); m != nil && !slices.Contains(p.Media, m.Text()) {
		return false
	}
	return true
}

func (p *Printer) createJob(q *request) {
	if !p.Accepting {
		q.fail(ipp.StatusServerErrorNotAcceptingJobs, "not accepting jobs")
		return
	}
	if p.Busy > 0 {
		p.Busy--
		q.fail(ipp.StatusServerErrorBusy, "busy")
		return
	}
	if !p.validate(q) {
		return
	}
	j := &Job{Transport: "ipp", Name: q.op.Get("job-name").Text(), User: q.op.Get("requesting-user-name").Text()}
	if g := q.req.Group(ipp.TagJob); g != nil {
		j.Attrs = g.Attrs
		if h := g.Get("job-hold-until").Text(); h != "" && h != "no-hold" {
			j.held = true
		}
	}
	p.add(j)
	if q.req.Op() == ipp.OpCreateJob && p.dropCreate {
		p.dropCreate = false
		p.mu.Unlock()
		defer p.mu.Lock()
		panic(http.ErrAbortHandler)
	}
	if q.req.Op() == ipp.OpPrintJob {
		j.Format = q.op.Get("document-format").Text()
		if err := p.receive(j, q); err != nil {
			j.aborted = true
			q.fail(ipp.StatusClientErrorBadRequest, "%v", err)
			return
		}
		p.done(j)
	}
	q.res.Add(ipp.TagJob, p.jobStatus(j, q)...)
}

func (p *Printer) sendDocument(q *request) {
	j := p.jobOf(q)
	if j == nil {
		return
	}
	if !j.doneAt.IsZero() || j.canceled || j.aborted {
		q.fail(ipp.StatusClientErrorNotPossible, "the job has its last document")
		return
	}
	f := q.op.Get("document-format").Text()
	if f != "" && !p.supports(f) {
		q.fail(ipp.StatusClientErrorDocumentFormatNotSupported, "%s", f)
		return
	}
	if c := q.op.Get("compression").Text(); c != "" && c != "none" && c != "gzip" && c != "deflate" {
		q.fail(ipp.StatusClientErrorCompressionNotSupported, "compression %s", c)
		return
	}
	last, ok := q.op.Get("last-document").Bool()
	if !ok {
		q.fail(ipp.StatusClientErrorBadRequest, "no last-document")
		return
	}
	j.Format = f
	peek, _ := q.body.Peek(1)
	if len(peek) > 0 || !last {
		if err := p.receive(j, q); err != nil {
			j.aborted = true
			q.fail(ipp.StatusClientErrorBadRequest, "%v", err)
			return
		}
	}
	if last {
		p.done(j)
	}
	q.res.Add(ipp.TagJob, p.jobStatus(j, q)...)
}

// receive reads a document of a job, as the request compresses it. The
// printer is locked; it is unlocked while the document comes.
func (p *Printer) receive(j *Job, q *request) error {
	cut, max := p.cutAfter, p.MaxDocument
	p.cutAfter = -1
	if max <= 0 {
		max = 256 << 20
	}
	p.mu.Unlock()
	defer p.mu.Lock()
	var body io.Reader = q.body
	if cut >= 0 {
		io.CopyN(io.Discard, body, cut)
		panic(http.ErrAbortHandler)
	}
	switch q.op.Get("compression").Text() {
	case "gzip":
		zr, err := gzip.NewReader(body)
		if err != nil {
			return err
		}
		body = zr
	case "deflate":
		body = flate.NewReader(body)
	}
	b, err := io.ReadAll(io.LimitReader(body, max+1))
	if err != nil {
		return err
	}
	if int64(len(b)) > max {
		return fmt.Errorf("a document of more than %d bytes", max)
	}
	p.mu.Lock()
	j.Documents = append(j.Documents, b)
	n, dir := len(j.Documents), p.Dir
	p.mu.Unlock()
	if dir != "" {
		if err := writeFile(dir, fmt.Sprintf("%d-%d.%s", j.ID, n, extension(j.Format, b)), b); err != nil {
			return err
		}
	}
	return nil
}

// extension is the extension of a document's file.
func extension(format string, b []byte) string {
	switch format {
	case "application/pdf", "pdf":
		return "pdf"
	case "image/pwg-raster":
		return "pwg"
	case "image/urf":
		return "urf"
	case "application/PCLm":
		return "pclm"
	case "application/postscript", "postscript":
		return "ps"
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	case "zpl", "sbpl", "tpcl", "escpos", "escp", "pcl":
		return format
	}
	if f := Sniff(b); f != "" {
		return extension(f, nil)
	}
	return "prn"
}

func (p *Printer) getJobs(q *request) {
	which := q.op.Get("which-jobs").Text()
	switch which {
	case "", "not-completed", "completed", "all":
	default:
		q.fail(ipp.StatusClientErrorAttributesOrValuesNotSupported, "which-jobs %s", which)
		q.res.Add(ipp.TagUnsupportedGroup, *q.op.Get("which-jobs"))
		return
	}
	mine, _ := q.op.Get("my-jobs").Bool()
	user := q.op.Get("requesting-user-name").Text()
	limit, _ := q.op.Get("limit").Int()
	req := ipp.Requested(q.op, []string{"job-id", "job-uri"})
	n := 0
	// the newest first, as CUPS and the printers list them
	for i := len(p.jobs) - 1; i >= 0; i-- {
		j := p.jobs[i]
		done := p.state(j).Done()
		if (which == "completed" && !done) || ((which == "" || which == "not-completed") && done) {
			continue
		}
		if mine && j.User != user {
			continue
		}
		if limit > 0 && n >= int(limit) {
			break
		}
		n++
		q.res.Groups = append(q.res.Groups, ipp.Group{Tag: ipp.TagJob, Attrs: ipp.Filter(p.jobAttrs(j, q), req, jobGroups)})
	}
}

// jobStatus are the attributes of a job in the answer to its creation.
func (p *Printer) jobStatus(j *Job, q *request) []ipp.Attribute {
	return ipp.Filter(p.jobAttrs(j, q), []string{"job-id", "job-uri", "job-state", "job-state-reasons", "job-state-message", "job-uuid"}, nil)
}

func (p *Printer) jobAttrs(j *Job, q *request) []ipp.Attribute {
	st := p.state(j)
	reason := "none"
	switch st {
	case ipp.JobPending:
		reason = "job-incoming"
		if !j.doneAt.IsZero() {
			reason = "none"
		}
	case ipp.JobPendingHeld:
		reason = "job-hold-until-specified"
	case ipp.JobProcessing:
		reason = "job-printing"
		if p.State == ipp.PrinterStopped {
			reason = "printer-stopped"
		}
	case ipp.JobCompleted:
		reason = "job-completed-successfully"
	case ipp.JobCanceled:
		reason = "job-canceled-by-user"
	case ipp.JobAborted:
		reason = "aborted-by-system"
	}
	copies := int32(1)
	if a := ipp.Get(j.Attrs, "copies"); a != nil {
		copies, _ = a.Int()
	}
	impressions := 0
	if st == ipp.JobCompleted {
		impressions = len(j.Documents) * int(copies)
	}
	size := 0
	for _, d := range j.Documents {
		size += len(d)
	}
	attrs := []ipp.Attribute{
		ipp.Attr("job-id", ipp.Integer(int(j.ID))),
		ipp.Attr("job-uri", ipp.URI(q.jobURI(j.ID))),
		ipp.Attr("job-uuid", ipp.URI(j.UUID)),
		ipp.Attr("job-printer-uri", ipp.URI(q.printerURI())),
		ipp.Attr("job-name", ipp.Name(j.Name)),
		ipp.Attr("job-originating-user-name", ipp.Name(j.User)),
		ipp.Attr("job-state", ipp.Enum(int(st))),
		ipp.Attr("job-state-reasons", ipp.Keyword(reason)),
		ipp.Attr("job-state-message", ipp.Text("")),
		ipp.Attr("number-of-documents", ipp.Integer(len(j.Documents))),
		ipp.Attr("job-printer-up-time", ipp.Integer(p.upTime(time.Now()))),
		ipp.Attr("job-k-octets", ipp.Integer((size+1023)/1024)),
		ipp.Attr("job-impressions-completed", ipp.Integer(impressions)),
		ipp.Attr("job-media-sheets-completed", ipp.Integer(impressions)),
		ipp.Attr("time-at-creation", ipp.Integer(p.upTime(j.Created))),
		ipp.Attr("date-time-at-creation", ipp.DateTime(j.Created)),
	}
	if !j.doneAt.IsZero() {
		attrs = append(attrs, ipp.Attr("time-at-processing", ipp.Integer(p.upTime(j.doneAt))),
			ipp.Attr("date-time-at-processing", ipp.DateTime(j.doneAt)))
	} else {
		attrs = append(attrs, ipp.Attr("time-at-processing", ipp.OutOfBand(ipp.TagNoValue)),
			ipp.Attr("date-time-at-processing", ipp.OutOfBand(ipp.TagNoValue)))
	}
	if st.Done() {
		end := j.doneAt
		if end.IsZero() {
			end = time.Now()
		}
		attrs = append(attrs, ipp.Attr("time-at-completed", ipp.Integer(p.upTime(end))), ipp.Attr("date-time-at-completed", ipp.DateTime(end)))
	} else {
		attrs = append(attrs, ipp.Attr("time-at-completed", ipp.OutOfBand(ipp.TagNoValue)),
			ipp.Attr("date-time-at-completed", ipp.OutOfBand(ipp.TagNoValue)))
	}
	if j.Format != "" && strings.Contains(j.Format, "/") {
		attrs = append(attrs, ipp.Attr("document-format-supplied", ipp.MimeType(j.Format)))
	}
	return append(attrs, j.Attrs...)
}

// upTime is printer-up-time at t: seconds since the printer started,
// from 1.
func (p *Printer) upTime(t time.Time) int {
	if t.Before(p.started) {
		return 1
	}
	return int(t.Sub(p.started)/time.Second) + 1
}
