package server

import (
	"bufio"
	"compress/flate"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/spool"
)

// Authentication of the IPP server.
const (
	authNone  = "none"
	authBasic = "basic"
)

// ippPath is where the printers of the queues are: /ipp/print/<queue>,
// and the default queue's at /ipp/print.
const ippPath = "/ipp/print"

// incomingTimeout is multiple-operation-time-out: how long a job made by
// Create-Job waits for its document.
var incomingTimeout = 5 * time.Minute

// goneFor is how long a job that ended before it reached the spool (its
// document did not come, or it was canceled first) is still told of.
const goneFor = 10 * time.Minute

// ippServer shows the queues of a server as IPP printers. The jobs it
// takes go to the spool, and their states are the spool's: IPP's own.
type ippServer struct {
	s    *Server
	cfg  *IPPConfig
	host string
	// queues are the queues shown, by name; def is the one at /ipp/print.
	queues map[string]*ippQueue
	order  []*ippQueue
	def    *ippQueue
	start  time.Time
	// secureOnly refuses ipp (basic, or require-tls).
	secureOnly bool
	// zone makes the records of a DNS zone, not of mDNS.
	zone bool

	mu sync.Mutex
	// incoming are the jobs made by Create-Job whose document has not
	// come; they are in memory only, and lost with the server.
	incoming map[int32]*incoming
	// changed is when each queue's state last changed, for
	// printer-state-change-time.
	changed map[string]time.Time

	// port and fingerprint are the listener's and the certificate's,
	// once served; cancel ends the watching and the advertising.
	port        int
	fingerprint string
	cancel      context.CancelFunc
	responder   *ipp.Responder
}

// ippQueue is a queue as an IPP printer.
type ippQueue struct {
	name, path string
	uuid       string
	info       string
	location   string
}

// incoming is a job that Create-Job made, waiting for its document, or one
// that ended before it reached the spool.
type incoming struct {
	id       int32
	uuid     string
	queue    *ippQueue
	user     string
	name     string
	job      jobRequest
	created  time.Time
	deadline time.Time
	// upload is the document received with a Send-Document that was not
	// the last (Close-Job submits it).
	upload  *spool.Upload
	docName string
	docType string
	// sending is set while a document arrives; state and reasons are
	// those of a job that ended here (canceled, aborted), at ended.
	sending bool
	state   ipp.JobState
	reasons []string
	ended   time.Time
}

func newIPPServer(s *Server) (*ippServer, error) {
	cfg := s.cfg.IPP
	x := &ippServer{s: s, cfg: cfg, queues: map[string]*ippQueue{}, start: time.Now(), incoming: map[int32]*incoming{},
		changed: map[string]time.Time{}, secureOnly: cfg.Auth == authBasic || cfg.RequireTLS}
	if err := x.init(); err != nil {
		return nil, err
	}
	for _, q := range x.order {
		x.changed[q.name] = x.start
	}
	return x, nil
}

// init makes the printers of the queues shown, and the host's name.
func (x *ippServer) init() error {
	cfg := x.cfg
	x.host = cfg.Host
	if x.host == "" {
		x.host = localHostName()
	}
	names := cfg.Queues
	if len(names) == 0 {
		names = sortedKeys(x.s.cfg.Queues)
	}
	for _, name := range names {
		qc := x.s.cfg.Queues[name]
		q := &ippQueue{name: name, path: ippPath + "/" + name, uuid: queueUUID(x.s.cfg.Spool, name), info: qc.Info, location: qc.Location}
		if q.info == "" {
			q.info = name
		}
		if q.location == "" {
			q.location = cfg.Location
		}
		x.queues[name] = q
		x.order = append(x.order, q)
	}
	if cfg.Default != "" {
		x.def = x.queues[cfg.Default]
	} else if len(x.order) == 1 {
		x.def = x.order[0]
	}
	return nil
}

// localHostName is the machine's name in .local: the name of the system,
// up to its first dot.
func localHostName() string {
	h, err := os.Hostname()
	if err != nil || h == "" {
		return "localhost"
	}
	h, _, _ = strings.Cut(h, ".")
	return h + ".local"
}

// queueUUID is the printer-uuid of a queue: the same for the same spool
// and queue, across restarts, as clients keep the printers they added by
// it.
func queueUUID(spool, queue string) string {
	sum := sha256.Sum256([]byte("bdf-printd\x00" + spool + "\x00" + queue))
	b := sum[:16]
	b[6] = b[6]&0x0f | 0x50
	b[8] = b[8]&0x3f | 0x80
	h := hex.EncodeToString(b)
	return "urn:uuid:" + h[:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:]
}

func newJobUUID() string {
	var b [16]byte
	rand.Read(b[:])
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	h := hex.EncodeToString(b[:])
	return "urn:uuid:" + h[:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:]
}

// upTime is the seconds since the server started, as printer-up-time,
// which is at least 1.
func (x *ippServer) upTime(t time.Time) int {
	if t.IsZero() {
		return 0
	}
	return max(int(t.Sub(x.start)/time.Second), 0) + 1
}

// handler is the IPP server's HTTP handler: the printers, their icons
// and a page of their state.
func (x *ippServer) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+ippPath, x.serveIPP)
	mux.HandleFunc("POST "+ippPath+"/{queue}", x.serveIPP)
	mux.HandleFunc("POST "+ippPath+"/{queue}/{job}", x.serveIPP)
	mux.HandleFunc("GET /icons/{file}", x.icon)
	mux.HandleFunc("GET /airprint.mobileconfig", x.serveProfile)
	mux.HandleFunc("GET /{$}", x.page)
	mux.HandleFunc("GET "+ippPath, x.page)
	mux.HandleFunc("GET "+ippPath+"/{queue}", x.page)
	return mux
}

// ippRequest is a request as it is served.
type ippRequest struct {
	w      http.ResponseWriter
	r      *http.Request
	req    *ipp.Message
	op     *ipp.Group
	body   *bufio.Reader
	res    *ipp.Message
	queue  *ippQueue
	secure bool
	caller *caller
	// after writes what follows the response: a document fetched.
	after func(io.Writer) error
}

// scheme is the scheme of the URIs told: that of the request.
func (q *ippRequest) scheme() string {
	if q.secure {
		return "ipps"
	}
	return "ipp"
}

// printerURI is the queue's URI as the client reached it.
func (q *ippRequest) printerURI() string { return q.scheme() + "://" + q.r.Host + q.queue.path }

func (q *ippRequest) jobURI(id int32) string {
	return q.printerURI() + "/" + strconv.Itoa(int(id))
}

func (q *ippRequest) fail(st ipp.Status, format string, args ...any) {
	q.res.Code = uint16(st)
	if format != "" {
		q.res.Add(ipp.TagOperation, ipp.Attr("status-message", ipp.Text(clip127(fmt.Sprintf(format, args...)))))
	}
}

func clip127(s string) string {
	if len(s) <= 255 {
		return s
	}
	for i := 255; i > 0; i-- {
		if (s[i] & 0xc0) != 0x80 {
			return s[:i]
		}
	}
	return s[:255]
}

// serveIPP serves an IPP request.
func (x *ippServer) serveIPP(w http.ResponseWriter, r *http.Request) {
	secure := isSecure(r)
	if !secure && x.secureOnly && !loopbackRemote(r) {
		// RFC 2817: the client is to come back with TLS
		w.Header().Set("Upgrade", "TLS/1.2, HTTP/1.1")
		w.Header().Set("Connection", "Upgrade")
		http.Error(w, "this printer is reached with ipps", http.StatusUpgradeRequired)
		return
	}
	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/ipp") {
		http.Error(w, "an IPP request is application/ipp", http.StatusBadRequest)
		return
	}
	q := &ippRequest{w: w, r: r, secure: secure}
	if name := r.PathValue("queue"); name != "" {
		q.queue = x.queues[name]
	} else {
		q.queue = x.def
	}
	if q.queue == nil {
		http.Error(w, "no such printer", http.StatusNotFound)
		return
	}
	req, body, err := ipp.DecodeFrom(bufio.NewReaderSize(r.Body, 64<<10), &ipp.Limits{MaxBytes: 1 << 20, MaxAttributes: 1000, MaxValues: 1000, MaxValueBytes: 64 << 10, MaxDepth: 4})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	q.req, q.body, q.op = req, body, req.Group(ipp.TagOperation)
	if q.op == nil {
		q.op = &ipp.Group{Tag: ipp.TagOperation}
	}
	q.res = ipp.NewResponse(ipp.StatusOK, req.RequestID, "utf-8", "en")
	q.res.Version = req.Version
	if req.Version.Major < 1 || req.Version.Major > 2 {
		q.res.Version = ipp.Version20
		q.fail(ipp.StatusServerErrorVersionNotSupported, "IPP/%s is not supported: 1.1 and 2.x are", req.Version)
	} else if x.checkRequest(q) {
		if x.authorize(w, q) {
			x.dispatch(q)
		} else {
			return
		}
	}
	b, err := q.res.Marshal()
	if err != nil {
		x.s.log.Error("ipp: the answer cannot be encoded", "op", req.Op().String(), "err", err)
		http.Error(w, "the answer cannot be encoded", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/ipp")
	if q.after == nil || !q.res.Status().OK() {
		w.Header().Set("Content-Length", strconv.Itoa(len(b)))
		w.Write(b)
		return
	}
	// the document follows the response, chunked; a failure cuts it
	w.Write(b)
	if err := q.after(w); err != nil {
		x.s.log.Error("ipp: the document could not be sent", "op", req.Op().String(), "queue", q.queue.name, "err", err)
		panic(http.ErrAbortHandler)
	}
}

// checkRequest checks what every request has (RFC 8011 §4.1): a request
// id, the charset then the language first, and the target.
func (x *ippServer) checkRequest(q *ippRequest) bool {
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
		q.res.Groups[0].Attrs[0] = ipp.Attr("attributes-charset", ipp.Charset("utf-8"))
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
	if q.op.Get("printer-uri") == nil && q.op.Get("job-uri") == nil {
		q.fail(ipp.StatusClientErrorBadRequest, "no printer-uri")
		return false
	}
	return true
}

// authorize finds who makes a request. With basic authentication, every
// operation but Get-Printer-Attributes (which tells how to authenticate)
// needs a token's name and secret; it writes the 401 itself.
func (x *ippServer) authorize(w http.ResponseWriter, q *ippRequest) bool {
	if isInfra(q) {
		return x.authorizeProxy(w, q)
	}
	named := q.op.Get("requesting-user-name").Text()
	if named != "" && !validUser(named) {
		named = ""
	}
	if x.cfg.Auth != authBasic {
		if named == "" {
			named = "anonymous"
		}
		q.caller = &caller{user: named, roles: []string{roleSubmit}}
		return true
	}
	name, secret, ok := q.r.BasicAuth()
	if ok {
		q.caller = x.s.auth.basic(name, secret)
	}
	if q.caller == nil {
		if q.req.Op() == ipp.OpGetPrinterAttributes {
			q.caller = &caller{user: "anonymous"}
			return true
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="bdf"`)
		http.Error(w, "a token's name and secret are required", http.StatusUnauthorized)
		return false
	}
	if !q.caller.mayQueue(q.queue.name) && q.req.Op() != ipp.OpGetPrinterAttributes {
		q.caller = &caller{user: q.caller.user}
	}
	return true
}

// isSecure tells a request that came with TLS, directly or through the
// listener that tells ipp from ipps on one port.
func isSecure(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	c, ok := r.Context().Value(connKey{}).(*sniffConn)
	return ok && c.secure
}

// loopbackRemote tells a request of the machine itself.
func loopbackRemote(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// operations are those the printers support.
var operations = []ipp.Op{
	ipp.OpPrintJob, ipp.OpValidateJob, ipp.OpCreateJob, ipp.OpSendDocument, ipp.OpCancelJob, ipp.OpGetJobAttributes,
	ipp.OpGetJobs, ipp.OpGetPrinterAttributes, ipp.OpHoldJob, ipp.OpReleaseJob, ipp.OpCancelMyJobs, ipp.OpCloseJob,
	ipp.OpIdentifyPrinter,
}

func (x *ippServer) dispatch(q *ippRequest) {
	op := q.req.Op()
	if q.caller.roles == nil && op != ipp.OpGetPrinterAttributes {
		q.fail(ipp.StatusClientErrorNotAuthorized, "the token may not print to %s", q.queue.name)
		return
	}
	if isInfra(q) {
		x.infra(q)
		return
	}
	switch op {
	case ipp.OpGetPrinterAttributes:
		x.getPrinterAttributes(q)
	case ipp.OpValidateJob:
		x.validateJob(q)
	case ipp.OpPrintJob:
		x.printJob(q)
	case ipp.OpCreateJob:
		x.createJob(q)
	case ipp.OpSendDocument:
		x.sendDocument(q)
	case ipp.OpCloseJob:
		x.closeJob(q)
	case ipp.OpCancelJob:
		x.cancelJob(q)
	case ipp.OpCancelMyJobs:
		x.cancelMyJobs(q)
	case ipp.OpGetJobAttributes:
		x.getJobAttributes(q)
	case ipp.OpGetJobs:
		x.getJobs(q)
	case ipp.OpHoldJob, ipp.OpReleaseJob:
		x.holdRelease(q)
	case ipp.OpIdentifyPrinter:
		x.s.log.Info("ipp: identify the printer", "queue", q.queue.name, "user", q.caller.user,
			"actions", strings.Join(q.op.Get("identify-actions").Strings(), ","), "message", q.op.Get("message").Text())
	default:
		q.fail(ipp.StatusServerErrorOperationNotSupported, "%s", op)
	}
}

// checkFormat checks the document-format and compression of a request.
func (x *ippServer) checkFormat(q *ippRequest) bool {
	if f := q.op.Get("document-format"); f != nil {
		if !slices.Contains(documentFormats, f.Text()) {
			q.fail(ipp.StatusClientErrorDocumentFormatNotSupported, "document-format %s", f.Text())
			q.res.Add(ipp.TagUnsupportedGroup, *f)
			return false
		}
	}
	if c := q.op.Get("compression"); c != nil {
		switch c.Text() {
		case "none", "gzip", "deflate":
		default:
			q.fail(ipp.StatusClientErrorCompressionNotSupported, "compression %s", c.Text())
			q.res.Add(ipp.TagUnsupportedGroup, *c)
			return false
		}
	}
	return true
}

// accepting refuses a job for a queue that does not take jobs.
func (x *ippServer) accepting(q *ippRequest) bool {
	info, err := x.s.sp.Queue(q.queue.name)
	if err != nil {
		q.fail(ipp.StatusServerErrorInternal, "%v", err)
		return false
	}
	if info.Rejecting {
		q.fail(ipp.StatusServerErrorNotAcceptingJobs, "the queue does not take jobs")
		return false
	}
	return true
}

func (x *ippServer) validateJob(q *ippRequest) {
	if !x.checkFormat(q) {
		return
	}
	if _, ok := x.jobRequest(q); !ok {
		return
	}
	x.accepting(q)
}

// printJob takes a job and its document.
func (x *ippServer) printJob(q *ippRequest) {
	if !x.checkFormat(q) || !x.accepting(q) {
		return
	}
	jr, ok := x.jobRequest(q)
	if !ok {
		return
	}
	if !x.slot(q) {
		return
	}
	defer x.unslot()
	up, ok := x.receive(q)
	if !ok {
		return
	}
	j, err := x.submit(q, jr, 0, "", up, q.op.Get("document-name").Text(), q.op.Get("document-format").Text())
	if err != nil {
		x.failSpool(q, err)
		return
	}
	q.res.Add(ipp.TagJob, x.jobStatusAttrs(q, &j)...)
}

// createJob makes a job whose document comes with Send-Document: the
// spool's id is taken now, the job is in memory until then.
func (x *ippServer) createJob(q *ippRequest) {
	if !x.accepting(q) {
		return
	}
	jr, ok := x.jobRequest(q)
	if !ok {
		return
	}
	x.expire()
	id, err := x.s.sp.ReserveID(q.r.Context())
	if err != nil {
		x.failSpool(q, err)
		return
	}
	now := time.Now()
	in := &incoming{id: id, uuid: newJobUUID(), queue: q.queue, user: q.caller.user, name: jr.ticket.JobName, job: jr,
		created: now, deadline: now.Add(incomingTimeout), state: ipp.JobPending, reasons: []string{"job-incoming"}}
	x.mu.Lock()
	if len(x.incoming) >= 1000 {
		x.mu.Unlock()
		q.fail(ipp.StatusServerErrorTooManyJobs, "too many jobs wait for their documents")
		return
	}
	x.incoming[id] = in
	x.mu.Unlock()
	x.s.log.Info("ipp: job created", "job", id, "queue", q.queue.name, "user", in.user)
	q.res.Add(ipp.TagJob, x.incomingAttrs(q, in)...)
}

// jobOf finds the job of a request (job-id with printer-uri, or
// job-uri): in memory (incoming), or in the spool.
func (x *ippServer) jobOf(q *ippRequest) (in *incoming, j *spool.Job, ok bool) {
	var id int32
	if a := q.op.Get("job-uri"); a != nil {
		u := a.Text()
		i := strings.LastIndexByte(u, '/')
		n, err := strconv.Atoi(u[i+1:])
		if i < 0 || err != nil || n <= 0 {
			q.fail(ipp.StatusClientErrorBadRequest, "job-uri %s", u)
			return nil, nil, false
		}
		id = int32(n)
	} else if a := q.op.Get("job-id"); a != nil {
		n, isInt := a.Int()
		if v, _ := a.First(); !isInt || v.Tag != ipp.TagInteger || n <= 0 {
			q.fail(ipp.StatusClientErrorBadRequest, "job-id is a positive integer")
			return nil, nil, false
		}
		id = n
	} else if n, err := strconv.Atoi(q.r.PathValue("job")); err == nil && n > 0 {
		id = int32(n)
	} else {
		q.fail(ipp.StatusClientErrorBadRequest, "no job-id")
		return nil, nil, false
	}
	x.expire()
	x.mu.Lock()
	in = x.incoming[id]
	x.mu.Unlock()
	if in != nil {
		if in.queue != q.queue {
			q.fail(ipp.StatusClientErrorNotFound, "no job %d in %s", id, q.queue.name)
			return nil, nil, false
		}
		return in, nil, true
	}
	job, err := x.s.sp.Get(q.r.Context(), id)
	if err != nil || job.Queue != q.queue.name {
		q.fail(ipp.StatusClientErrorNotFound, "no job %d", id)
		return nil, nil, false
	}
	return nil, &job, true
}

// sendDocument takes the document of a job Create-Job made.
func (x *ippServer) sendDocument(q *ippRequest) {
	in, j, ok := x.jobOf(q)
	if !ok {
		return
	}
	last, isBool := q.op.Get("last-document").Bool()
	if q.op.Get("last-document") == nil || !isBool {
		q.fail(ipp.StatusClientErrorBadRequest, "last-document is required")
		return
	}
	if j != nil || in.state.Done() {
		q.fail(ipp.StatusClientErrorNotPossible, "the job has its document")
		if in != nil && in.state.Done() {
			q.fail(ipp.StatusClientErrorNotPossible, "the job is %s", in.state)
		}
		return
	}
	if !q.caller.acts(&spool.Job{User: in.user}) {
		q.fail(ipp.StatusClientErrorNotAuthorized, "the job is another's")
		return
	}
	if !x.checkFormat(q) {
		return
	}
	x.mu.Lock()
	if in.sending {
		x.mu.Unlock()
		q.fail(ipp.StatusServerErrorBusy, "a document of the job is arriving")
		return
	}
	in.sending = true
	had := in.upload != nil
	x.mu.Unlock()
	defer func() {
		x.mu.Lock()
		in.sending = false
		x.mu.Unlock()
	}()
	// a Send-Document without data that is the last closes the job
	if peek, _ := q.body.Peek(1); len(peek) == 0 {
		if !last {
			q.res.Add(ipp.TagJob, x.incomingAttrs(q, in)...)
			return
		}
		x.close(q, in)
		return
	}
	if had {
		q.fail(ipp.StatusServerErrorMultipleDocumentJobsNotSupported, "a job has one document")
		return
	}
	if !x.slot(q) {
		return
	}
	defer x.unslot()
	up, ok := x.receive(q)
	if !ok {
		return
	}
	x.mu.Lock()
	if in.state.Done() {
		// canceled (or ended) while the document arrived
		x.mu.Unlock()
		up.Discard()
		q.fail(ipp.StatusClientErrorNotPossible, "the job is %s", in.state)
		return
	}
	in.upload, in.docName, in.docType = up, q.op.Get("document-name").Text(), q.op.Get("document-format").Text()
	in.deadline = time.Now().Add(incomingTimeout)
	if pw := q.op.Get("document-password"); pw != nil {
		in.job.password = pw.Text()
	}
	x.mu.Unlock()
	if !last {
		q.res.Add(ipp.TagJob, x.incomingAttrs(q, in)...)
		return
	}
	x.close(q, in)
}

// closeJob ends the documents of a job Create-Job made.
func (x *ippServer) closeJob(q *ippRequest) {
	in, j, ok := x.jobOf(q)
	if !ok {
		return
	}
	if j != nil {
		q.res.Add(ipp.TagJob, x.jobStatusAttrs(q, j)...)
		return
	}
	if !q.caller.acts(&spool.Job{User: in.user}) {
		q.fail(ipp.StatusClientErrorNotAuthorized, "the job is another's")
		return
	}
	x.mu.Lock()
	busy := in.sending
	x.mu.Unlock()
	if busy {
		q.fail(ipp.StatusServerErrorBusy, "a document of the job is arriving")
		return
	}
	x.close(q, in)
}

// close submits a job Create-Job made, with its document; one without a
// document is aborted.
func (x *ippServer) close(q *ippRequest, in *incoming) {
	x.mu.Lock()
	up := in.upload
	in.upload = nil
	done := in.state.Done()
	x.mu.Unlock()
	if done {
		if up != nil {
			up.Discard()
		}
		q.res.Add(ipp.TagJob, x.incomingAttrs(q, in)...)
		return
	}
	if up == nil {
		x.end(in, ipp.JobAborted, "job-data-insufficient")
		q.res.Add(ipp.TagJob, x.incomingAttrs(q, in)...)
		return
	}
	j, err := x.submit(q, in.job, in.id, in.uuid, up, in.docName, in.docType)
	if err != nil {
		x.end(in, ipp.JobAborted, "aborted-by-system")
		x.failSpool(q, err)
		return
	}
	x.mu.Lock()
	canceled := in.state.Done() // by Cancel-Job while it was submitted
	delete(x.incoming, in.id)
	x.mu.Unlock()
	if canceled {
		if c, err := x.s.sp.Cancel(q.r.Context(), j.ID, "job-canceled-by-user"); err == nil {
			j = c
		}
	}
	q.res.Add(ipp.TagJob, x.jobStatusAttrs(q, &j)...)
}

// end ends a job that did not reach the spool.
func (x *ippServer) end(in *incoming, st ipp.JobState, reason string) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if in.upload != nil {
		in.upload.Discard()
		in.upload = nil
	}
	in.state, in.reasons, in.ended = st, []string{reason}, time.Now()
}

// expire ends the jobs whose documents did not come in time, and forgets
// those ended long ago.
func (x *ippServer) expire() {
	now := time.Now()
	x.mu.Lock()
	var late []*incoming
	for id, in := range x.incoming {
		switch {
		case in.state.Done() && now.Sub(in.ended) > goneFor:
			delete(x.incoming, id)
		case !in.state.Done() && !in.sending && now.After(in.deadline):
			late = append(late, in)
		}
	}
	x.mu.Unlock()
	for _, in := range late {
		x.s.log.Info("ipp: the document of a job did not come", "job", in.id, "queue", in.queue.name)
		x.end(in, ipp.JobAborted, "aborted-by-system")
	}
}

// slot takes one of the server's submissions at once.
func (x *ippServer) slot(q *ippRequest) bool {
	select {
	case x.s.uploads <- struct{}{}:
		return true
	default:
		q.fail(ipp.StatusServerErrorBusy, "the server is receiving as many documents as it takes at once")
		return false
	}
}

func (x *ippServer) unslot() { <-x.s.uploads }

// receive reads the document that follows a request, decompressed.
func (x *ippServer) receive(q *ippRequest) (*spool.Upload, bool) {
	// a document has longer than a request to arrive
	http.NewResponseController(q.w).SetReadDeadline(time.Now().Add(uploadTime))
	var r io.Reader = q.body
	switch q.op.Get("compression").Text() {
	case "gzip":
		zr, err := gzip.NewReader(q.body)
		if err != nil {
			q.fail(ipp.StatusClientErrorCompressionError, "gzip: %v", err)
			return nil, false
		}
		defer zr.Close()
		r = zr
	case "deflate":
		zr := flate.NewReader(q.body)
		defer zr.Close()
		r = zr
	}
	up, err := x.s.sp.Receive(q.r.Context(), r)
	if err != nil {
		var corrupt flate.CorruptInputError
		if errors.Is(err, gzip.ErrHeader) || errors.Is(err, gzip.ErrChecksum) || errors.As(err, &corrupt) {
			q.fail(ipp.StatusClientErrorCompressionError, "%v", err)
			return nil, false
		}
		x.failSpool(q, err)
		return nil, false
	}
	return up, true
}

// submit submits a job to the spool.
func (x *ippServer) submit(q *ippRequest, jr jobRequest, id int32, uuid string, up *spool.Upload, docName, docType string) (spool.Job, error) {
	name := jr.ticket.JobName
	if docName != "" {
		name = docName
	}
	if name == "" {
		name = "untitled"
	}
	if docType == "application/octet-stream" || docType == "" {
		docType = ""
	}
	sub := spool.Submission{ID: id, UUID: uuid, Queue: q.queue.name, Name: name, User: q.caller.user, Type: docType, Upload: up,
		Ticket: jr.ticket, Options: jr.options, Priority: jr.ticket.Priority, HoldUntil: jr.ticket.HoldUntil, Password: jr.password}
	if pw := q.op.Get("document-password"); pw != nil {
		sub.Password = pw.Text()
	}
	j, err := x.s.sp.Submit(q.r.Context(), sub)
	if err != nil {
		return j, err
	}
	x.s.log.Info("ipp: job received", "job", j.ID, "queue", j.Queue, "user", j.User, "format", j.Document.Format, "size", j.Document.Size)
	return j, nil
}

// failSpool answers an error of the spool.
func (x *ippServer) failSpool(q *ippRequest, err error) {
	switch {
	case errors.Is(err, spool.ErrUnknownFormat):
		q.fail(ipp.StatusClientErrorDocumentFormatNotSupported, "the document's format is not one the server prints")
	case errors.Is(err, spool.ErrBadDocument):
		q.fail(ipp.StatusClientErrorDocumentFormatError, "%v", err)
	case errors.Is(err, spool.ErrTooLarge):
		q.fail(ipp.StatusClientErrorRequestEntityTooLarge, "%v", err)
	case errors.Is(err, spool.ErrWrongPassword):
		q.fail(ipp.StatusClientErrorDocumentPasswordError, "the password does not open the document")
	case errors.Is(err, spool.ErrPasswordRequired):
		q.fail(ipp.StatusClientErrorDocumentPasswordError, "the document is encrypted: document-password")
	case errors.Is(err, spool.ErrNotAccepting):
		q.fail(ipp.StatusServerErrorNotAcceptingJobs, "the queue does not take jobs")
	case errors.Is(err, spool.ErrFull):
		q.fail(ipp.StatusServerErrorBusy, "the spool is full")
	case errors.Is(err, spool.ErrClosed):
		q.fail(ipp.StatusServerErrorServiceUnavailable, "the server is stopping")
	case errors.Is(err, spool.ErrInvalid):
		q.fail(ipp.StatusClientErrorBadRequest, "%v", err)
	case errors.Is(err, spool.ErrNotPossible):
		q.fail(ipp.StatusClientErrorNotPossible, "%v", err)
	case errors.Is(err, spool.ErrNotFound):
		q.fail(ipp.StatusClientErrorNotFound, "%v", err)
	case q.r.Context().Err() != nil:
		q.fail(ipp.StatusClientErrorBadRequest, "the document did not arrive")
	default:
		x.s.log.Error("ipp", "op", q.req.Op().String(), "queue", q.queue.name, "err", err)
		q.fail(ipp.StatusServerErrorInternal, "the server failed; its log tells why")
	}
}

func (x *ippServer) cancelJob(q *ippRequest) {
	in, j, ok := x.jobOf(q)
	if !ok {
		return
	}
	if in != nil {
		if !q.caller.acts(&spool.Job{User: in.user}) {
			q.fail(ipp.StatusClientErrorNotAuthorized, "the job is another's")
			return
		}
		if in.state.Done() {
			q.fail(ipp.StatusClientErrorNotPossible, "the job is %s", in.state)
			return
		}
		x.end(in, ipp.JobCanceled, "job-canceled-by-user")
		return
	}
	if !q.caller.acts(j) {
		q.fail(ipp.StatusClientErrorNotAuthorized, "the job is another's")
		return
	}
	if j.State.Done() {
		q.fail(ipp.StatusClientErrorNotPossible, "the job is %s", j.State)
		return
	}
	reason := "job-canceled-by-user"
	if !q.caller.owns(j) {
		reason = "job-canceled-by-operator"
	}
	if _, err := x.s.sp.Cancel(q.r.Context(), j.ID, reason); err != nil {
		x.failSpool(q, err)
	}
}

// cancelMyJobs cancels the caller's jobs of the queue that are not done,
// or those of job-ids.
func (x *ippServer) cancelMyJobs(q *ippRequest) {
	var ids []int32
	if a := q.op.Get("job-ids"); a != nil {
		ids = a.Ints()
	}
	pick := func(id int32) bool { return ids == nil || slices.Contains(ids, id) }
	x.mu.Lock()
	var mine []*incoming
	for _, in := range x.incoming {
		if in.queue == q.queue && in.user == q.caller.user && !in.state.Done() && pick(in.id) {
			mine = append(mine, in)
		}
	}
	x.mu.Unlock()
	for _, in := range mine {
		x.end(in, ipp.JobCanceled, "job-canceled-by-user")
	}
	for _, j := range x.s.sp.List(spool.Filter{Queue: q.queue.name, User: q.caller.user, Limit: 1000}) {
		if !j.State.Done() && pick(j.ID) {
			x.s.sp.Cancel(q.r.Context(), j.ID, "job-canceled-by-user")
		}
	}
}

func (x *ippServer) holdRelease(q *ippRequest) {
	in, j, ok := x.jobOf(q)
	if !ok {
		return
	}
	if in != nil {
		q.fail(ipp.StatusClientErrorNotPossible, "the job has no document yet")
		return
	}
	if !q.caller.acts(j) {
		q.fail(ipp.StatusClientErrorNotAuthorized, "the job is another's")
		return
	}
	var err error
	if q.req.Op() == ipp.OpHoldJob {
		_, err = x.s.sp.Hold(q.r.Context(), j.ID)
	} else {
		_, err = x.s.sp.Release(q.r.Context(), j.ID, q.op.Get("document-password").Text())
	}
	if err != nil {
		x.failSpool(q, err)
	}
}

func (x *ippServer) getJobAttributes(q *ippRequest) {
	in, j, ok := x.jobOf(q)
	if !ok {
		return
	}
	var attrs []ipp.Attribute
	if in != nil {
		attrs = x.incomingAttrs(q, in)
	} else {
		attrs = x.jobAttrs(q, j, q.caller.sees(j))
	}
	q.res.Add(ipp.TagJob, ipp.Filter(attrs, ipp.Requested(q.op, jobAll), jobGroups)...)
}

// getJobs lists the jobs of the queue.
func (x *ippServer) getJobs(q *ippRequest) {
	which := q.op.Get("which-jobs").Text()
	switch which {
	case "", "not-completed", "completed", "all":
	default:
		q.fail(ipp.StatusClientErrorAttributesOrValuesNotSupported, "which-jobs %s", which)
		q.res.Add(ipp.TagUnsupportedGroup, *q.op.Get("which-jobs"))
		return
	}
	mine, _ := q.op.Get("my-jobs").Bool()
	limit := 0
	if a := q.op.Get("limit"); a != nil {
		n, _ := a.Int()
		if n <= 0 {
			q.fail(ipp.StatusClientErrorBadRequest, "limit %d", n)
			return
		}
		limit = int(n)
	}
	var ids []int32
	if a := q.op.Get("job-ids"); a != nil {
		ids = a.Ints()
		which = "all"
	}
	f := spool.Filter{Queue: q.queue.name, Done: which == "completed" || which == "all", Limit: 1000}
	if mine {
		f.User = q.caller.user
	}
	keep := func(id int32, st ipp.JobState) bool {
		if ids != nil && !slices.Contains(ids, id) {
			return false
		}
		switch which {
		case "completed":
			return st.Done()
		case "all":
			return true
		}
		return !st.Done()
	}
	req := ipp.Requested(q.op, []string{"job-id", "job-uri"})
	var groups [][]ipp.Attribute
	x.expire()
	x.mu.Lock()
	var ins []*incoming
	for _, in := range x.incoming {
		if in.queue == q.queue && (!mine || in.user == q.caller.user) && keep(in.id, in.state) {
			ins = append(ins, in)
		}
	}
	x.mu.Unlock()
	slices.SortFunc(ins, func(a, b *incoming) int { return int(a.id - b.id) })
	jobs := x.s.sp.List(f)
	// the newest first, as printers list them
	slices.Reverse(jobs)
	for _, in := range ins {
		groups = append(groups, ipp.Filter(x.incomingAttrs(q, in), req, jobGroups))
	}
	for i := range jobs {
		if keep(jobs[i].ID, jobs[i].State) {
			groups = append(groups, ipp.Filter(x.jobAttrs(q, &jobs[i], q.caller.sees(&jobs[i])), req, jobGroups))
		}
	}
	if limit > 0 && len(groups) > limit {
		groups = groups[:limit]
	}
	for _, g := range groups {
		q.res.Groups = append(q.res.Groups, ipp.Group{Tag: ipp.TagJob, Attrs: g})
	}
}

// basic finds the token of a name and secret (IPP's Basic
// authentication).
func (a *authenticator) basic(name, secret string) *caller {
	sum := sha256.Sum256([]byte(secret))
	var found *token
	for i := range a.tokens {
		t := &a.tokens[i]
		if subtle.ConstantTimeCompare(sum[:], t.sum[:]) == 1 && subtle.ConstantTimeCompare([]byte(name), []byte(t.Name)) == 1 {
			found = t
		}
	}
	if found == nil {
		return nil
	}
	c := &caller{token: found.Name, user: found.Name, roles: found.Roles, queues: found.Queues}
	if len(c.roles) == 0 {
		c.roles = []string{roleSubmit}
	}
	return c
}

// watchQueues notes when the queues' states change, for
// printer-state-change-time, until ctx ends.
func (x *ippServer) watchQueues(ctx context.Context) {
	events, stop := x.s.sp.Subscribe(64)
	defer stop()
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-events:
			if !ok {
				return
			}
			now := time.Now()
			x.mu.Lock()
			switch {
			case ev.Queue != nil:
				x.changed[ev.Queue.Name] = now
			case ev.Job != nil:
				x.changed[ev.Job.Queue] = now
			case ev.Printer != nil:
				for name := range x.changed {
					if qc, ok := x.s.cfg.Queues[name]; ok && slices.Contains(qc.Printers, ev.Printer.Name) {
						x.changed[name] = now
					}
				}
			}
			x.mu.Unlock()
			if ev.Printer != nil {
				x.readvertise()
			}
		}
	}
}
