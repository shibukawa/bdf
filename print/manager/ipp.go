package manager

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
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/print/ipp"
)

// Every node is an IPP printer at /ipp/print/<name>: a sender (a WMS, a
// business system, an office PC) prints to it as to any printer, over
// HTTP(S) on the manager's own port, and the job follows the topology
// from there. Only the operations of one request are taken (Print-Job,
// not Create-Job and Send-Document), so that a request is all there is to
// a job: no instance waits for a document. /ipp/print itself answers
// CUPS-Get-Printers and CUPS-Get-Classes with every node, as a CUPS
// server lists its destinations.

// IPPPath is where the nodes' printers are.
const IPPPath = "/ipp/print"

var ippOperations = []ipp.Op{ipp.OpPrintJob, ipp.OpValidateJob, ipp.OpCancelJob, ipp.OpGetJobAttributes, ipp.OpGetJobs, ipp.OpGetPrinterAttributes,
	ipp.OpCupsGetPrinters, ipp.OpCupsGetClasses, ipp.OpCupsGetDefault}

var ippFormats = []string{"application/octet-stream", "application/pdf", "application/x-bdf", "image/pwg-raster", "image/urf", "image/jpeg", "image/png"}

type ippRequest struct {
	w   http.ResponseWriter
	r   *http.Request
	req *ipp.Message
	op  *ipp.Group
	res *ipp.Message
	// dest is the node the request is to ("" for the server itself),
	// kind its kind and title its title.
	dest  string
	kind  string
	title string
	// by is who sent it, as verified (a token, or anonymous).
	by *Principal
}

// user is who the request is from, to show: the token's name, or else
// requesting-user-name.
func (q *ippRequest) user() string {
	if q.by.Kind == PrincipalToken {
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

func (q *ippRequest) uri(dest string) string {
	return q.scheme() + "://" + q.r.Host + IPPPath + "/" + dest
}

func (q *ippRequest) printerURI() string { return q.uri(q.dest) }

// destInfo tells the kind and the title of a destination.
func (m *Manager) destInfo(ctx context.Context, name string) (kind, title string, err error) {
	n, err := m.Node(ctx, name)
	if err != nil {
		return "", "", err
	}
	return n.Kind, or(n.Title, name), nil
}

// serveIPP serves an IPP request to a node's printer, or to the server
// (no node).
func (h *handler) serveIPP(w http.ResponseWriter, r *http.Request) {
	h.serveIPPTo(w, r, r.PathValue("dest"))
}

// serveIPPTo serves an IPP request to a destination ("" the server), on
// the manager's server or a CUPS server's port.
func (h *handler) serveIPPTo(w http.ResponseWriter, r *http.Request, dest string) {
	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/ipp") {
		http.Error(w, "an IPP request is application/ipp", http.StatusBadRequest)
		return
	}
	by, ok := h.submitter(r)
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="bdf print manager"`)
		http.Error(w, "a token's name and secret are required", http.StatusUnauthorized)
		return
	}
	ctx := r.Context()
	q := &ippRequest{w: w, r: r, by: by}
	if dest != "" {
		kind, title, err := h.m.destInfo(ctx, dest)
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
		h.dispatch(ctx, q, body)
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

func (h *handler) dispatch(ctx context.Context, q *ippRequest, body io.Reader) {
	switch q.req.Op() {
	case ipp.OpCupsGetPrinters:
		h.cupsPrinters(ctx, q, false)
		return
	case ipp.OpCupsGetClasses:
		h.cupsPrinters(ctx, q, true)
		return
	case ipp.OpCupsGetDefault:
		q.fail(ipp.StatusClientErrorNotFound, "no default destination: name one")
		return
	}
	if q.dest == "" {
		q.fail(ipp.StatusClientErrorNotFound, "the server is not a printer: ask %s/<name>", IPPPath)
		return
	}
	switch q.req.Op() {
	case ipp.OpGetPrinterAttributes:
		h.printerAttributes(ctx, q)
	case ipp.OpValidateJob:
		if f := q.text("document-format"); f != "" && !contains(ippFormats, f) {
			q.fail(ipp.StatusClientErrorDocumentFormatNotSupported, "document-format %s", f)
		}
	case ipp.OpPrintJob:
		h.printJob(ctx, q, body)
	case ipp.OpGetJobAttributes:
		j := h.job(ctx, q)
		if j != nil {
			q.res.Add(ipp.TagJob, h.jobAttributes(q, j)...)
		}
	case ipp.OpGetJobs:
		h.jobs(ctx, q)
	case ipp.OpCancelJob:
		j := h.job(ctx, q)
		if j == nil {
			return
		}
		if _, err := h.m.Cancel(ctx, j.ID, Mark{Worker: clip(q.user(), 80)}); err != nil {
			q.fail(ipp.StatusClientErrorNotPossible, "the job is %s", j.State)
		}
	default:
		q.fail(ipp.StatusServerErrorOperationNotSupported, "%s is not supported", q.req.Op())
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
	sum := sha256.Sum256([]byte("bdf-print-manager\x00" + name))
	b := sum[:16]
	b[6] = b[6]&0x0f | 0x50
	b[8] = b[8]&0x3f | 0x80
	x := hex.EncodeToString(b)
	return "urn:uuid:" + x[:8] + "-" + x[8:12] + "-" + x[12:16] + "-" + x[16:20] + "-" + x[20:]
}

// cupsPrinterType is CUPS's printer-type of a node: the bit of a class
// for classes, none for the rest.
func cupsPrinterType(kind string) int {
	if kind == KindClass {
		return 1 // CUPS_PRINTER_CLASS
	}
	return 0
}

// deviceURI is the device-uri a node is shown with: a printer's URL, or
// the kind of the node as a scheme.
func deviceURI(n *TopoNode) string {
	if n.Kind == KindPrinter {
		return n.URL
	}
	return n.Kind + "://" + n.Name
}

// cupsPrinters answers CUPS-Get-Printers (every node) or CUPS-Get-Classes
// (the classes, with their members).
func (h *handler) cupsPrinters(ctx context.Context, q *ippRequest, classes bool) {
	nodes, err := h.m.Topology(ctx)
	if err != nil {
		q.fail(ipp.StatusServerErrorInternal, "%v", err)
		return
	}
	for i := range nodes {
		n := &nodes[i]
		if classes && n.Kind != KindClass {
			continue
		}
		attrs := []ipp.Attribute{
			ipp.Attr("printer-name", ipp.Name(n.Name)),
			ipp.Attr("printer-uri-supported", ipp.URI(q.uri(n.Name))),
			ipp.Attr("printer-info", ipp.Text(or(n.Title, n.Name))),
			ipp.Attr("printer-location", ipp.Text(or(n.Title, n.Name))),
			ipp.Attr("printer-make-and-model", ipp.Text("BDF Print Manager "+n.Kind)),
			ipp.Attr("printer-uuid", ipp.URI(printerUUID(n.Name))),
			ipp.Attr("printer-state", ipp.Enum(3)),
			ipp.Attr("printer-state-reasons", ipp.Keyword("none")),
			ipp.Attr("printer-is-accepting-jobs", ipp.Boolean(true)),
			ipp.Attr("printer-type", ipp.Integer(cupsPrinterType(n.Kind))),
			ipp.Attr("device-uri", ipp.URI(deviceURI(n))),
		}
		if n.Kind == KindClass {
			names := make([]ipp.Value, len(n.Members))
			uris := make([]ipp.Value, len(n.Members))
			for i, mem := range n.Members {
				names[i], uris[i] = ipp.Name(mem), ipp.URI(q.uri(mem))
			}
			attrs = append(attrs, ipp.Attr("member-names", names...), ipp.Attr("member-uris", uris...))
		}
		q.res.Groups = append(q.res.Groups, ipp.Group{Tag: ipp.TagPrinter, Attrs: attrs})
	}
	if len(q.res.Groups) == 1 { // the operation group alone
		q.fail(ipp.StatusClientErrorNotFound, "no destinations")
	}
}

func (h *handler) printerAttributes(ctx context.Context, q *ippRequest) {
	queued := 0
	if jobs, err := h.m.JobsOf(ctx, q.dest, 0); err == nil {
		for _, j := range jobs {
			if !j.State.ended() {
				queued++
			}
		}
	}
	ops := make([]ipp.Value, len(ippOperations))
	for i, op := range ippOperations {
		ops[i] = ipp.Enum(int(op))
	}
	formats := make([]ipp.Value, len(ippFormats))
	for i, f := range ippFormats {
		formats[i] = ipp.MimeType(f)
	}
	auth := "none"
	if len(h.m.o.SubmitTokens) > 0 {
		auth = "basic"
	}
	security := "none"
	if q.scheme() == "ipps" {
		security = "tls"
	}
	q.res.Add(ipp.TagPrinter,
		ipp.Attr("printer-uri-supported", ipp.URI(q.printerURI())),
		ipp.Attr("uri-security-supported", ipp.Keyword(security)),
		ipp.Attr("uri-authentication-supported", ipp.Keyword(auth)),
		ipp.Attr("printer-name", ipp.Name(q.dest)),
		ipp.Attr("printer-info", ipp.Text(q.title)),
		ipp.Attr("printer-location", ipp.Text(q.title)),
		ipp.Attr("printer-make-and-model", ipp.Text("BDF Print Manager "+q.kind)),
		ipp.Attr("printer-uuid", ipp.URI(printerUUID(q.dest))),
		ipp.Attr("printer-type", ipp.Integer(cupsPrinterType(q.kind))),
		ipp.Attr("printer-state", ipp.Enum(3)), // idle
		ipp.Attr("printer-state-reasons", ipp.Keyword("none")),
		ipp.Attr("printer-is-accepting-jobs", ipp.Boolean(true)),
		ipp.Attr("queued-job-count", ipp.Integer(queued)),
		ipp.Attr("printer-up-time", ipp.Integer(int(h.m.now().Unix()&0x7fffffff))),
		ipp.Attr("ipp-versions-supported", ipp.Keyword("1.1"), ipp.Keyword("2.0")),
		ipp.Attr("operations-supported", ops...),
		ipp.Attr("charset-configured", ipp.Charset("utf-8")),
		ipp.Attr("charset-supported", ipp.Charset("utf-8")),
		ipp.Attr("natural-language-configured", ipp.Language("en")),
		ipp.Attr("generated-natural-language-supported", ipp.Language("en")),
		ipp.Attr("document-format-default", ipp.MimeType("application/octet-stream")),
		ipp.Attr("document-format-supported", formats...),
		ipp.Attr("compression-supported", ipp.Keyword("none")),
		ipp.Attr("pdl-override-supported", ipp.Keyword("not-attempted")),
		ipp.Attr("multiple-document-jobs-supported", ipp.Boolean(false)),
		ipp.Attr("job-ids-supported", ipp.Boolean(true)),
		ipp.Attr("copies-default", ipp.Integer(1)),
		ipp.Attr("copies-supported", ipp.NewRange(1, 1)),
		ipp.Attr("sides-default", ipp.Keyword("one-sided")),
		ipp.Attr("sides-supported", ipp.Keyword("one-sided")),
		ipp.Attr("media-default", ipp.Keyword("iso_a4_210x297mm")),
		ipp.Attr("media-supported", ipp.Keyword("iso_a4_210x297mm"), ipp.Keyword("na_letter_8.5x11in")),
		ipp.Attr("print-color-mode-default", ipp.Keyword("auto")),
		ipp.Attr("print-color-mode-supported", ipp.Keyword("auto"), ipp.Keyword("color"), ipp.Keyword("monochrome")),
		ipp.Attr("color-supported", ipp.Boolean(true)),
	)
}

// printJob takes a Print-Job: its document follows the request.
func (h *handler) printJob(ctx context.Context, q *ippRequest, body io.Reader) {
	if c := q.text("compression"); c != "" && c != "none" {
		q.fail(ipp.StatusClientErrorCompressionNotSupported, "compression %s", c)
		return
	}
	if f := q.text("document-format"); f != "" && !contains(ippFormats, f) {
		q.fail(ipp.StatusClientErrorDocumentFormatNotSupported, "document-format %s", f)
		return
	}
	f, err := os.CreateTemp(h.m.o.TempDir, "bdf-pm-in-*")
	if err != nil {
		q.fail(ipp.StatusServerErrorInternal, "%v", err)
		return
	}
	defer os.Remove(f.Name())
	defer f.Close()
	n, err := io.Copy(f, io.LimitReader(body, h.m.o.MaxDocument+1))
	if err != nil {
		q.fail(ipp.StatusClientErrorDocumentAccessError, "the document: %v", err)
		return
	}
	if n > h.m.o.MaxDocument {
		q.fail(ipp.StatusClientErrorRequestEntityTooLarge, "the document is larger than %d bytes", h.m.o.MaxDocument)
		return
	}
	name := firstOr(q.text("job-name"), q.text("document-name"))
	if q.by.Kind == PrincipalApp {
		// an app prints to the endpoints that list it
		if v, ok := h.viewerOf(q.r); ok {
			if err := h.m.appMaySend(ctx, v, q.dest, false); err != nil {
				q.fail(ipp.StatusClientErrorNotAuthorized, "%v", err)
				return
			}
		}
	}
	j, err := h.m.Submit(ctx, q.dest, Input{Name: name, By: q.by, ClaimedUser: q.text("requesting-user-name"), Source: "ipp", Doc: f, Size: n})
	switch {
	case errors.Is(err, ErrFormat):
		q.fail(ipp.StatusClientErrorDocumentFormatError, "%v", err)
		return
	case errors.Is(err, ErrInvalid):
		q.fail(ipp.StatusClientErrorBadRequest, "%v", err)
		return
	case err != nil:
		h.m.o.Log.Error("ipp: print-job", "dest", q.dest, "err", err)
		q.fail(ipp.StatusServerErrorInternal, "the job could not be kept")
		return
	}
	st, reason := ippState(j)
	q.res.Add(ipp.TagJob,
		ipp.Attr("job-id", ipp.Integer(int(j.Number))),
		ipp.Attr("job-uri", ipp.URI(q.printerURI()+"/"+strconv.Itoa(int(j.Number)))),
		ipp.Attr("job-state", ipp.Enum(int(st))),
		ipp.Attr("job-state-reasons", ipp.Keyword(reason)),
	)
}

// job finds the job a request names (job-id, or job-uri), of the
// request's destination.
func (h *handler) job(ctx context.Context, q *ippRequest) *Job {
	var id int32
	if a := q.attr("job-id"); a != nil {
		id, _ = a.Int()
	} else if uri := q.text("job-uri"); uri != "" {
		n, _ := strconv.Atoi(uri[strings.LastIndexByte(uri, '/')+1:])
		id = int32(n)
	}
	if id <= 0 {
		q.fail(ipp.StatusClientErrorBadRequest, "no job-id")
		return nil
	}
	j, err := h.m.JobOfNumber(ctx, id)
	if err != nil || j.Dest != q.dest {
		q.fail(ipp.StatusClientErrorNotFound, "no job %d", id)
		return nil
	}
	return j
}

// ippState is a job's job-state and job-state-reasons.
func ippState(j *Job) (ipp.JobState, string) {
	switch j.State {
	case JobProcessing:
		return ipp.JobProcessing, "job-printing"
	case JobHeld:
		return ipp.JobPendingHeld, "job-hold-until-specified"
	case JobCanceled:
		return ipp.JobCanceled, "job-canceled-by-user"
	case JobAborted:
		return ipp.JobAborted, "aborted-by-system"
	}
	return ipp.JobCompleted, "job-completed-successfully"
}

func (h *handler) jobAttributes(q *ippRequest, j *Job) []ipp.Attribute {
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
	}
	if !j.Ended.IsZero() {
		attrs = append(attrs, ipp.Attr("time-at-completed", ipp.Integer(int(j.Ended.Unix()&0x7fffffff))))
	}
	return attrs
}

// jobs answers Get-Jobs: the jobs sent to the destination, those not
// completed by default (RFC 8011), the completed ones with which-jobs
// completed, all with which-jobs all.
func (h *handler) jobs(ctx context.Context, q *ippRequest) {
	which := firstOr(q.text("which-jobs"), "not-completed")
	limit := 1000
	if a := q.attr("limit"); a != nil {
		if n, ok := a.Int(); ok && n > 0 {
			limit = int(n)
		}
	}
	jobs, err := h.m.JobsOf(ctx, q.dest, 0)
	if err != nil {
		q.fail(ipp.StatusServerErrorInternal, "%v", err)
		return
	}
	n := 0
	for i := range jobs {
		j := &jobs[i]
		switch which {
		case "completed":
			if !j.State.ended() {
				continue
			}
		case "not-completed":
			if j.State.ended() {
				continue
			}
		}
		if n++; n > limit {
			break
		}
		// a group for each job (Add would join them)
		q.res.Groups = append(q.res.Groups, ipp.Group{Tag: ipp.TagJob, Attrs: h.jobAttributes(q, j)})
	}
}

// submitter finds who sends jobs over IPP and the API: the admin token,
// a submit token's name and secret (Basic, or bearer name:secret), an
// app (its token, or Basic client_id:secret), or, with no SubmitTokens,
// anyone (anonymous). A person's session sends from the pages (submit).
func (h *handler) submitter(r *http.Request) (*Principal, bool) {
	if v, ok := h.viewerOf(r); ok && (v.Kind == PrincipalToken || v.Kind == PrincipalApp) {
		return v.Principal(), true
	}
	if r.Header.Get("Authorization") != "" {
		return nil, false // credentials given, and not known
	}
	if len(h.m.o.SubmitTokens) == 0 {
		return &Principal{Kind: PrincipalAnonymous, Via: ViaNone}, true
	}
	return nil, false
}
