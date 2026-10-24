package ipp

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/ipp"
)

// Path is where the nodes' printers are.
const Path = "/ipp/print"

// ippOperations are the operations of the nodes' printers on the
// manager's own server, where one request is all there is to a job (no
// instance waits for a document); portOperations those on a CUPS server
// node's own port, where one process listens: Create-Job, Send-Document
// and Close-Job too, as IPP Everywhere (and the clients of Mopria and
// Windows) ask (created.go).
var (
	ippOperations = []ipp.Op{ipp.OpPrintJob, ipp.OpValidateJob, ipp.OpCancelJob, ipp.OpGetJobAttributes, ipp.OpGetJobs, ipp.OpGetPrinterAttributes,
		ipp.OpCancelMyJobs, ipp.OpIdentifyPrinter, ipp.OpCupsGetPrinters, ipp.OpCupsGetClasses, ipp.OpCupsGetDefault}
	portOperations = append(slices.Clone(ippOperations), ipp.OpCreateJob, ipp.OpSendDocument, ipp.OpCloseJob)
)

// ippFormats are the formats the printers take: PDF and the rasters of
// AirPrint (image/urf), IPP Everywhere (image/pwg-raster) and Mopria
// (application/PCLm, a PDF of strips of image, which the PDF converter
// reads), .bdf, images and data.
var ippFormats = []string{"application/octet-stream", "application/pdf", "application/x-bdf", "image/pwg-raster", "image/urf", "application/PCLm", "image/jpeg", "image/png",
	"text/csv", "application/json"}

// maxIPPMeta bounds the job attributes of a request taken as metadata.
const maxIPPMeta = 32

// ippMeta is the metadata of a Print-Job request: its job template
// attributes (the job group: media, copies, sides …) and its
// document-format, as ipp.<name>.
func ippMeta(req *ipp.Message) map[string]any {
	out := map[string]any{}
	add := func(a *ipp.Attribute) {
		if len(out) >= maxIPPMeta || !docflow.MetaKeyOK(a.Name) {
			return
		}
		if n, ok := a.Int(); ok && len(a.Values) == 1 && (a.Values[0].Tag == ipp.TagInteger || a.Values[0].Tag == ipp.TagEnum) {
			out["ipp."+a.Name] = n
			return
		}
		if t := a.Text(); t != "" {
			out["ipp."+a.Name] = clip(t, 255)
			return
		}
		out["ipp."+a.Name] = clip(a.String(), 255)
	}
	if g := req.Group(ipp.TagOperation); g != nil {
		if a := g.Get("document-format"); a != nil {
			add(a)
		}
	}
	if g := req.Group(ipp.TagJob); g != nil {
		for i := range g.Attrs {
			add(&g.Attrs[i])
		}
	}
	return out
}

type ippRequest struct {
	w   http.ResponseWriter
	r   *http.Request
	req *ipp.Message
	op  *ipp.Group
	res *ipp.Message
	// dest is the node the request is to ("" for the server itself),
	// kind its kind and title its title. server is the CUPS server node
	// whose own port the request came on ("" on the manager's own
	// server): its one printer is the node, at /printers/<name>.
	dest   string
	kind   string
	title  string
	server string
	// by is who sent it, as verified (a token, or anonymous).
	by *docflow.Principal
}

// user is who the request is from, to show: the token's name, or else
// requesting-user-name.
func (q *ippRequest) user() string {
	if q.by.Kind == docflow.PrincipalToken {
		return q.by.ID
	}
	return firstOr(q.text("requesting-user-name"), "anonymous")
}

func (q *ippRequest) attr(name string) *ipp.Attribute {
	if q.op == nil {
		return nil
	}
	return q.op.Get(name)
}

func (q *ippRequest) text(name string) string {
	if a := q.attr(name); a != nil {
		return a.Text()
	}
	return ""
}

func (q *ippRequest) fail(st ipp.Status, format string, args ...any) {
	q.res.Code = uint16(st)
	q.res.Add(ipp.TagOperation, ipp.Attr("status-message", ipp.Text(clip(fmt.Sprintf(format, args...), 255))))
}

func (q *ippRequest) scheme() string {
	if q.r.TLS != nil || strings.EqualFold(q.r.Header.Get("X-Forwarded-Proto"), "https") {
		return "ipps"
	}
	return "ipp"
}

// uri is a destination's printer-uri: on the manager's own server
// /ipp/print/<name>, on a CUPS server's port /printers/<name>, as CUPS
// names its printers.
func (q *ippRequest) uri(dest string) string {
	if q.server != "" {
		return q.scheme() + "://" + q.r.Host + "/printers/" + dest
	}
	return q.scheme() + "://" + q.r.Host + Path + "/" + dest
}

func (q *ippRequest) printerURI() string { return q.uri(q.dest) }

// operations are the operations of the request's printer: those of a
// CUPS server node's own port, or of the manager's own server.
func (q *ippRequest) operations() []ipp.Op {
	if q.server != "" {
		return portOperations
	}
	return ippOperations
}

// destInfo tells the kind and the title of a destination: a CUPS server
// node, the one kind that takes IPP.
func (s *server) destInfo(ctx context.Context, name string) (kind, title string, err error) {
	n, err := s.m.Node(ctx, name)
	if err != nil {
		return "", "", err
	}
	if n.Kind != KindCups {
		return "", "", docflow.ErrNotFound
	}
	return n.Kind, or(n.Title, name), nil
}

// serveIPP serves an IPP request to a node's printer, or to the server
// (no node).
func (s *server) serveIPP(w http.ResponseWriter, r *http.Request) {
	s.serveIPPTo(w, r, r.PathValue("dest"), "")
}

// serveIPPTo serves an IPP request to a destination ("" the server), on
// the manager's server (server "") or a CUPS server node's own port
// (server its name: the one printer there).
func (s *server) serveIPPTo(w http.ResponseWriter, r *http.Request, dest, server string) {
	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/ipp") {
		http.Error(w, "an IPP request is application/ipp", http.StatusBadRequest)
		return
	}
	by, ok := s.m.Sender(r)
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="bdf docflow"`)
		http.Error(w, "a token's name and secret are required", http.StatusUnauthorized)
		return
	}
	ctx := r.Context()
	q := &ippRequest{w: w, r: r, by: by, server: server}
	if dest != "" {
		kind, title, err := s.destInfo(ctx, dest)
		if err != nil {
			http.Error(w, "no such printer", http.StatusNotFound)
			return
		}
		q.dest, q.kind, q.title = dest, kind, title
	}
	req, body, err := ipp.DecodeFrom(bufio.NewReaderSize(r.Body, 64<<10), &ipp.Limits{MaxBytes: 1 << 20, MaxAttributes: 1000, MaxValues: 1000, MaxValueBytes: 64 << 10, MaxDepth: 4})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	q.req, q.op = req, req.Group(ipp.TagOperation)
	q.res = ipp.NewResponse(ipp.StatusOK, req.RequestID, "utf-8", "en")
	q.res.Version = req.Version
	switch {
	case req.Version.Major < 1 || req.Version.Major > 2:
		q.res.Version = ipp.Version20
		q.fail(ipp.StatusServerErrorVersionNotSupported, "IPP/%s is not supported: 1.1 and 2.x are", req.Version)
	case req.RequestID == 0 || q.op == nil || len(q.op.Attrs) < 2 ||
		q.op.Attrs[0].Name != "attributes-charset" || q.op.Attrs[1].Name != "attributes-natural-language":
		q.fail(ipp.StatusClientErrorBadRequest, "a request id, then attributes-charset and attributes-natural-language")
	default:
		s.dispatch(ctx, q, body)
	}
	discard(body)
	b, err := q.res.Marshal()
	if err != nil {
		http.Error(w, "the answer cannot be encoded", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/ipp")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}

func (s *server) dispatch(ctx context.Context, q *ippRequest, body *bufio.Reader) {
	switch q.req.Op() {
	case ipp.OpCupsGetPrinters:
		s.cupsPrinters(ctx, q, false)
		return
	case ipp.OpCupsGetClasses:
		s.cupsPrinters(ctx, q, true)
		return
	case ipp.OpCupsGetDefault:
		q.fail(ipp.StatusClientErrorNotFound, "no default destination: name one")
		return
	}
	if q.attr("printer-uri") == nil {
		// a printer's operations name the printer (RFC 8011 §4.2)
		q.fail(ipp.StatusClientErrorBadRequest, "no printer-uri")
		return
	}
	if q.dest == "" {
		if q.server != "" {
			q.fail(ipp.StatusClientErrorNotFound, "the server is not a printer: ask /printers/%s", q.server)
			return
		}
		q.fail(ipp.StatusClientErrorNotFound, "the server is not a printer: ask %s/<name>", Path)
		return
	}
	if !slices.Contains(q.operations(), q.req.Op()) {
		q.fail(ipp.StatusServerErrorOperationNotSupported, "%s is not supported", q.req.Op())
		return
	}
	switch q.req.Op() {
	case ipp.OpGetPrinterAttributes:
		s.printerAttributes(ctx, q)
	case ipp.OpValidateJob:
		if f := q.text("document-format"); f != "" && !contains(ippFormats, f) {
			q.fail(ipp.StatusClientErrorDocumentFormatNotSupported, "document-format %s", f)
		}
	case ipp.OpPrintJob:
		s.printJob(ctx, q, body)
	case ipp.OpCreateJob:
		s.createJob(ctx, q)
	case ipp.OpSendDocument:
		s.sendDocument(ctx, q, body)
	case ipp.OpCloseJob:
		s.closeJob(ctx, q)
	case ipp.OpGetJobAttributes:
		c, j := s.jobOrCreated(ctx, q)
		switch {
		case c != nil:
			q.res.Add(ipp.TagJob, s.createdAttributes(q, c)...)
		case j != nil:
			q.res.Add(ipp.TagJob, s.jobAttributes(q, j)...)
		}
	case ipp.OpGetJobs:
		s.jobs(ctx, q)
	case ipp.OpCancelJob:
		c, j := s.jobOrCreated(ctx, q)
		if c != nil {
			if !createdOf(s.m).end(c, ipp.JobCanceled, "job-canceled-by-user") {
				q.fail(ipp.StatusClientErrorNotPossible, "the job is %s", c.state)
			}
			return
		}
		if j == nil {
			return
		}
		if _, err := s.m.Cancel(ctx, j.ID, docflow.Mark{Worker: clip(q.user(), 80)}); err != nil {
			q.fail(ipp.StatusClientErrorNotPossible, "the job is %s", j.State)
		}
	case ipp.OpCancelMyJobs:
		s.cancelMyJobs(ctx, q)
	case ipp.OpIdentifyPrinter:
		// nothing to flash: the log tells it
		s.m.Options().Log.Info("ipp: identify the printer", "dest", q.dest, "user", q.user(),
			"actions", strings.Join(q.attr("identify-actions").Strings(), ","), "message", q.text("message"))
	}
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if strings.EqualFold(v, s) {
			return true
		}
	}
	return false
}

// printerUUID is a node's printer-uuid: the same for the same name.
func printerUUID(name string) string {
	// The label is from the days of the Print Manager; it stays, so that
	// the clients that added a node keep its uuid.
	sum := sha256.Sum256([]byte("bdf-print-manager\x00" + name))
	b := sum[:16]
	b[6] = b[6]&0x0f | 0x50
	b[8] = b[8]&0x3f | 0x80
	x := hex.EncodeToString(b)
	return "urn:uuid:" + x[:8] + "-" + x[8:12] + "-" + x[12:16] + "-" + x[16:20] + "-" + x[20:]
}

// cupsPrinters answers CUPS-Get-Printers (the CUPS server nodes on the
// manager's own server; on a CUPS server node's port, that node alone) or
// CUPS-Get-Classes (none: a class is no IPP printer).
func (s *server) cupsPrinters(ctx context.Context, q *ippRequest, classes bool) {
	nodes, err := s.m.Topology(ctx)
	if err != nil {
		q.fail(ipp.StatusServerErrorInternal, "%v", err)
		return
	}
	for i := range nodes {
		n := &nodes[i]
		if classes || n.Kind != KindCups || q.server != "" && n.Name != q.server {
			continue
		}

		attrs := []ipp.Attribute{
			ipp.Attr("printer-name", ipp.Name(n.Name)),
			ipp.Attr("printer-uri-supported", ipp.URI(q.uri(n.Name))),
			ipp.Attr("printer-info", ipp.Text(or(n.Title, n.Name))),
			ipp.Attr("printer-location", ipp.Text(or(n.Title, n.Name))),
			ipp.Attr("printer-make-and-model", ipp.Text(makeModel)),
			ipp.Attr("printer-uuid", ipp.URI(printerUUID(n.Name))),
			ipp.Attr("printer-state", ipp.Enum(3)),
			ipp.Attr("printer-state-reasons", ipp.Keyword("none")),
			ipp.Attr("printer-is-accepting-jobs", ipp.Boolean(true)),
			ipp.Attr("printer-type", ipp.Integer(0)),
			ipp.Attr("device-uri", ipp.URI(n.Kind+"://"+n.Name)),
		}
		q.res.Groups = append(q.res.Groups, ipp.Group{Tag: ipp.TagPrinter, Attrs: attrs})
	}
	if len(q.res.Groups) == 1 { // the operation group alone
		q.fail(ipp.StatusClientErrorNotFound, "no destinations")
	}
}

// printerAttributes answers Get-Printer-Attributes: the node as an
// AirPrint and IPP Everywhere printer (airPrintAttrs), the attributes
// requested (all by default); media-col-database when asked by name, as
// printers tell it.
func (s *server) printerAttributes(ctx context.Context, q *ippRequest) {
	if f := q.attr("document-format"); f != nil && !contains(ippFormats, f.Text()) {
		q.fail(ipp.StatusClientErrorDocumentFormatNotSupported, "document-format %s", f.Text())
		q.res.Add(ipp.TagUnsupportedGroup, *f)
		return
	}
	queued := 0
	if jobs, err := s.m.JobsOf(ctx, q.dest, 0); err == nil {
		for _, j := range jobs {
			if !j.State.Ended() {
				queued++
			}
		}
	}
	auth := "none"
	if len(s.m.Options().SubmitTokens) > 0 {
		auth = "basic"
	}
	security := "none"
	if q.scheme() == "ipps" {
		security = "tls"
	}
	attrs := s.airPrintAttrs(q, queued, auth, security)
	req := ipp.Requested(q.op, []string{"all"})
	if !slices.Contains(req, "media-col-database") {
		attrs = slices.DeleteFunc(attrs, func(a ipp.Attribute) bool { return a.Name == "media-col-database" })
	}
	q.res.Add(ipp.TagPrinter, ipp.Filter(attrs, req, ipp.PrinterGroups(attrs))...)
}

// printJob takes a Print-Job: its document follows the request.
func (s *server) printJob(ctx context.Context, q *ippRequest, body io.Reader) {
	if q.by.Kind == docflow.PrincipalApp {
		// an app prints to the entries that list it
		if err := s.m.MaySend(ctx, q.r, q.dest); err != nil {
			q.fail(ipp.StatusClientErrorNotAuthorized, "%v", err)
			return
		}
	}
	s.take(ctx, q, body, docflow.Input{Name: firstOr(q.text("job-name"), q.text("document-name")), By: q.by, ClaimedUser: q.text("requesting-user-name"),
		Meta: s.meta(q)})
}

// take takes the document that follows a request (Print-Job,
// Send-Document) as the job of in (its name, sender and metadata, and
// its number when Create-Job took one), and answers the job's state.
// What the job is not is the error.
func (s *server) take(ctx context.Context, q *ippRequest, body io.Reader, in docflow.Input) error {
	if c := q.text("compression"); c != "" && c != "none" {
		q.fail(ipp.StatusClientErrorCompressionNotSupported, "compression %s", c)
		return docflow.ErrInvalid
	}
	if f := q.text("document-format"); f != "" && !contains(ippFormats, f) {
		q.fail(ipp.StatusClientErrorDocumentFormatNotSupported, "document-format %s", f)
		return docflow.ErrFormat
	}
	f, err := os.CreateTemp(s.m.Options().TempDir, "bdf-pm-in-*")
	if err != nil {
		q.fail(ipp.StatusServerErrorInternal, "%v", err)
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	n, err := io.Copy(f, io.LimitReader(body, s.m.Options().MaxDocument+1))
	if err != nil {
		q.fail(ipp.StatusClientErrorDocumentAccessError, "the document: %v", err)
		return err
	}
	if n > s.m.Options().MaxDocument {
		q.fail(ipp.StatusClientErrorRequestEntityTooLarge, "the document is larger than %d bytes", s.m.Options().MaxDocument)
		return docflow.ErrInvalid
	}
	in.Source, in.Doc, in.Size, in.ContentType = "ipp", f, n, q.text("document-format")
	j, err := s.m.Submit(ctx, q.dest, in)
	if j == nil && (errors.Is(err, docflow.ErrFormat) || errors.Is(err, docflow.ErrInvalid)) {
		// no job tells it: a client retries what is refused (a Mac's queue
		// every few seconds), so the log does
		s.m.Options().Log.Warn("ipp: print-job refused", "dest", q.dest, "name", in.Name, "format", q.text("document-format"), "user", q.user(), "err", err)
	}
	switch {
	case errors.Is(err, docflow.ErrFormat):
		q.fail(ipp.StatusClientErrorDocumentFormatError, "%v", err)
		return err
	case errors.Is(err, docflow.ErrInvalid), errors.Is(err, docflow.ErrRefused):
		// a refusal of an assert after the node: the job is kept, aborted
		q.fail(ipp.StatusClientErrorBadRequest, "%v", err)
		return err
	case err != nil:
		s.m.Options().Log.Error("ipp: print-job", "dest", q.dest, "err", err)
		q.fail(ipp.StatusServerErrorInternal, "the job could not be kept")
		return err
	}
	q.res.Add(ipp.TagJob, s.jobStatus(q, j)...)
	return nil
}

// jobStatus are the attributes a job's operations answer of it: its id,
// its uri and its state.
func (s *server) jobStatus(q *ippRequest, j *docflow.Job) []ipp.Attribute {
	st, reason := ippState(j)
	return []ipp.Attribute{
		ipp.Attr("job-id", ipp.Integer(int(j.Number))),
		ipp.Attr("job-uri", ipp.URI(q.printerURI()+"/"+strconv.Itoa(int(j.Number)))),
		ipp.Attr("job-state", ipp.Enum(int(st))),
		ipp.Attr("job-state-reasons", ipp.Keyword(reason)),
	}
}

// meta is the metadata of a Print-Job (ippMeta), with document-metadata
// (an operation attribute, PWG 5100.13) as ipp.document-metadata. The
// attributes a printer node asks the printer again (job-account-id,
// job-accounting-user-id, document-metadata and x-…, see printMeta) are
// checked as print.Meta takes them: one whose value cannot be taken is
// answered unsupported and not kept; one of several values is kept as
// the list of them.
func (s *server) meta(q *ippRequest) map[string]any {
	meta := ippMeta(q.req)
	var m print.Meta
	var unsupported []ipp.Attribute
	if a := q.op.Get("document-metadata"); a != nil {
		if _, ok := m.Take(*a); ok {
			meta["ipp.document-metadata"] = a.Strings()
		} else {
			unsupported = append(unsupported, *a)
		}
	}
	if g := q.req.Group(ipp.TagJob); g != nil {
		for _, a := range g.Attrs {
			known, ok := m.Take(a)
			switch {
			case known && !ok:
				unsupported = append(unsupported, a)
				delete(meta, "ipp."+a.Name)
			case known && len(a.Values) > 1:
				if _, kept := meta["ipp."+a.Name]; kept {
					meta["ipp."+a.Name] = a.Strings()
				}
			}
		}
	}
	if len(unsupported) > 0 {
		q.res.Add(ipp.TagUnsupportedGroup, unsupported...)
		q.res.Code = uint16(ipp.StatusOKIgnoredOrSubstituted)
	}
	return meta
}

// jobID is the job a request names (job-id, or job-uri), 0 (and the
// request failed) when none.
func (q *ippRequest) jobID() int32 {
	var id int32
	if a := q.attr("job-id"); a != nil {
		id, _ = a.Int()
	} else if uri := q.text("job-uri"); uri != "" {
		n, _ := strconv.Atoi(uri[strings.LastIndexByte(uri, '/')+1:])
		id = int32(n)
	}
	if id <= 0 {
		q.fail(ipp.StatusClientErrorBadRequest, "no job-id")
		return 0
	}
	return id
}

// jobOrCreated finds the job a request names, of the request's
// destination: one Create-Job made that waits for its document (or
// ended without one), or a job.
func (s *server) jobOrCreated(ctx context.Context, q *ippRequest) (*created, *docflow.Job) {
	id := q.jobID()
	if id == 0 {
		return nil, nil
	}
	if c := createdOf(s.m).get(id); c != nil && c.dest == q.dest {
		return c, nil
	}
	j, err := s.m.JobOfNumber(ctx, id)
	if err != nil || j.Dest != q.dest {
		q.fail(ipp.StatusClientErrorNotFound, "no job %d", id)
		return nil, nil
	}
	return nil, j
}

// ippState is a job's job-state and job-state-reasons.
func ippState(j *docflow.Job) (ipp.JobState, string) {
	switch j.State {
	case docflow.JobProcessing:
		return ipp.JobProcessing, "job-printing"
	case docflow.JobHeld:
		return ipp.JobPendingHeld, "job-hold-until-specified"
	case docflow.JobCanceled:
		return ipp.JobCanceled, "job-canceled-by-user"
	case docflow.JobAborted:
		return ipp.JobAborted, "aborted-by-system"
	}
	return ipp.JobCompleted, "job-completed-successfully"
}

func (s *server) jobAttributes(q *ippRequest, j *docflow.Job) []ipp.Attribute {
	st, reason := ippState(j)
	attrs := []ipp.Attribute{
		ipp.Attr("job-id", ipp.Integer(int(j.Number))),
		ipp.Attr("job-uri", ipp.URI(q.printerURI()+"/"+strconv.Itoa(int(j.Number)))),
		ipp.Attr("job-printer-uri", ipp.URI(q.printerURI())),
		ipp.Attr("job-name", ipp.Name(firstOr(j.Name, "job-"+j.ID))),
		ipp.Attr("job-originating-user-name", ipp.Name(firstOr(j.ClaimedUser, firstOr(j.SubmittedBy.Label(), "anonymous")))),
		ipp.Attr("job-state", ipp.Enum(int(st))),
		ipp.Attr("job-state-reasons", ipp.Keyword(reason)),
		ipp.Attr("job-impressions-completed", ipp.Integer(j.Pages)),
		ipp.Attr("time-at-creation", ipp.Integer(int(j.Received.Unix()&0x7fffffff))),
		// a job is processed as it is received
		ipp.Attr("time-at-processing", ipp.Integer(int(j.Received.Unix()&0x7fffffff))),
		ipp.Attr("job-printer-up-time", ipp.Integer(int(s.m.Now().Unix()&0x7fffffff))),
	}
	if !j.Ended.IsZero() {
		attrs = append(attrs, ipp.Attr("time-at-completed", ipp.Integer(int(j.Ended.Unix()&0x7fffffff))))
	}
	pm := printMeta(j.Meta)
	return append(attrs, pm.Attrs()...)
}

// jobs answers Get-Jobs: the jobs sent to the destination, those not
// completed by default (RFC 8011), the completed ones with which-jobs
// completed, all with which-jobs all; their job-id and job-uri unless
// requested-attributes asks for more.
func (s *server) jobs(ctx context.Context, q *ippRequest) {
	req := ipp.Requested(q.op, []string{"job-id", "job-uri"})
	which := firstOr(q.text("which-jobs"), "not-completed")
	limit := 1000
	if a := q.attr("limit"); a != nil {
		if n, ok := a.Int(); ok && n > 0 {
			limit = int(n)
		}
	}
	jobs, err := s.m.JobsOf(ctx, q.dest, 0)
	if err != nil {
		q.fail(ipp.StatusServerErrorInternal, "%v", err)
		return
	}
	n := 0
	// the jobs Create-Job made, newest first: before the jobs
	for _, c := range createdOf(s.m).of(q.dest) {
		if which == "completed" && !c.state.Done() || which == "not-completed" && c.state.Done() {
			continue
		}
		if n++; n > limit {
			return
		}
		q.res.Groups = append(q.res.Groups, ipp.Group{Tag: ipp.TagJob, Attrs: ipp.Filter(s.createdAttributes(q, c), req, nil)})
	}
	for i := range jobs {
		j := &jobs[i]
		switch which {
		case "completed":
			if !j.State.Ended() {
				continue
			}
		case "not-completed":
			if j.State.Ended() {
				continue
			}
		}
		if n++; n > limit {
			break
		}
		// a group for each job (Add would join them)
		q.res.Groups = append(q.res.Groups, ipp.Group{Tag: ipp.TagJob, Attrs: ipp.Filter(s.jobAttributes(q, j), req, nil)})
	}
}
