package manager

import (
	"bufio"
	"context"
	"crypto/sha256"
	"crypto/subtle"
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

// Each location is an IPP printer at /ipp/print/<location>: a sender
// (a WMS, a business system, an office PC) prints to it as to any
// printer, over HTTP(S) on the manager's own port. The job is completed
// as soon as its slip is on the board. Only the operations of one request
// are taken (Print-Job, not Create-Job and Send-Document), so that a
// request is all there is to a job: no instance waits for a document.

// IPPPath is where the locations' printers are.
const IPPPath = "/ipp/print"

var ippOperations = []ipp.Op{ipp.OpPrintJob, ipp.OpValidateJob, ipp.OpCancelJob, ipp.OpGetJobAttributes, ipp.OpGetJobs, ipp.OpGetPrinterAttributes}

var ippFormats = []string{"application/octet-stream", "application/pdf", "application/x-bdf", "image/pwg-raster", "image/urf", "image/jpeg", "image/png"}

type ippRequest struct {
	w   http.ResponseWriter
	r   *http.Request
	req *ipp.Message
	op  *ipp.Group
	res *ipp.Message
	loc *Location
	// user is who sent it: the token's name, or what the sender said
	// without a token.
	user string
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

func (q *ippRequest) printerURI() string {
	return q.scheme() + "://" + q.r.Host + IPPPath + "/" + q.loc.Name
}

// serveIPP serves an IPP request to a location's printer.
func (h *handler) serveIPP(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/ipp") {
		http.Error(w, "an IPP request is application/ipp", http.StatusBadRequest)
		return
	}
	user, ok := h.submitter(r)
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="bdf print manager"`)
		http.Error(w, "a token's name and secret are required", http.StatusUnauthorized)
		return
	}
	ctx := r.Context()
	loc, err := h.m.Location(ctx, r.PathValue("location"))
	if err != nil {
		http.Error(w, "no such printer", http.StatusNotFound)
		return
	}
	req, body, err := ipp.DecodeFrom(bufio.NewReaderSize(r.Body, 64<<10), &ipp.Limits{MaxBytes: 1 << 20, MaxAttributes: 1000, MaxValues: 1000, MaxValueBytes: 64 << 10, MaxDepth: 4})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	q := &ippRequest{w: w, r: r, req: req, op: req.Group(ipp.TagOperation), loc: loc, user: user}
	q.res = ipp.NewResponse(ipp.StatusOK, req.RequestID, "utf-8", "en")
	q.res.Version = req.Version
	if user == "" {
		q.user = firstOr(q.text("requesting-user-name"), "anonymous")
	}
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
	case ipp.OpGetPrinterAttributes:
		h.printerAttributes(ctx, q)
	case ipp.OpValidateJob:
		if f := q.text("document-format"); f != "" && !contains(ippFormats, f) {
			q.fail(ipp.StatusClientErrorDocumentFormatNotSupported, "document-format %s", f)
		}
	case ipp.OpPrintJob:
		h.printJob(ctx, q, body)
	case ipp.OpGetJobAttributes:
		s := h.job(ctx, q)
		if s != nil {
			q.res.Add(ipp.TagJob, h.jobAttributes(q, s)...)
		}
	case ipp.OpGetJobs:
		h.jobs(ctx, q)
	case ipp.OpCancelJob:
		s := h.job(ctx, q)
		if s == nil {
			return
		}
		if _, err := h.m.Cancel(ctx, s.ID, Mark{Worker: q.user}); err != nil {
			q.fail(ipp.StatusClientErrorNotPossible, "the job is %s", s.State)
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

// printerUUID is a location's printer-uuid: the same for the same name.
func printerUUID(name string) string {
	sum := sha256.Sum256([]byte("bdf-print-manager\x00" + name))
	b := sum[:16]
	b[6] = b[6]&0x0f | 0x50
	b[8] = b[8]&0x3f | 0x80
	x := hex.EncodeToString(b)
	return "urn:uuid:" + x[:8] + "-" + x[8:12] + "-" + x[12:16] + "-" + x[16:20] + "-" + x[20:]
}

func (h *handler) printerAttributes(ctx context.Context, q *ippRequest) {
	open := 0
	if b, err := h.m.Board(ctx, q.loc.Name); err == nil {
		open = len(b.Open)
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
		ipp.Attr("printer-name", ipp.Name(q.loc.Name)),
		ipp.Attr("printer-info", ipp.Text(or(q.loc.Title, q.loc.Name))),
		ipp.Attr("printer-location", ipp.Text(or(q.loc.Title, q.loc.Name))),
		ipp.Attr("printer-make-and-model", ipp.Text("BDF Print Manager")),
		ipp.Attr("printer-uuid", ipp.URI(printerUUID(q.loc.Name))),
		ipp.Attr("printer-state", ipp.Enum(3)), // idle
		ipp.Attr("printer-state-reasons", ipp.Keyword("none")),
		ipp.Attr("printer-is-accepting-jobs", ipp.Boolean(true)),
		ipp.Attr("queued-job-count", ipp.Integer(open)),
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
	s, err := h.m.Submit(ctx, q.loc.Name, Input{Name: name, User: q.user, Source: "ipp", Doc: f, Size: n})
	switch {
	case errors.Is(err, ErrFormat):
		q.fail(ipp.StatusClientErrorDocumentFormatError, "%v", err)
		return
	case errors.Is(err, ErrInvalid):
		q.fail(ipp.StatusClientErrorBadRequest, "%v", err)
		return
	case err != nil:
		h.m.o.Log.Error("ipp: print-job", "location", q.loc.Name, "err", err)
		q.fail(ipp.StatusServerErrorInternal, "the job could not be kept")
		return
	}
	q.res.Add(ipp.TagJob,
		ipp.Attr("job-id", ipp.Integer(int(s.Job))),
		ipp.Attr("job-uri", ipp.URI(q.printerURI()+"/"+strconv.Itoa(int(s.Job)))),
		ipp.Attr("job-state", ipp.Enum(int(ipp.JobCompleted))),
		ipp.Attr("job-state-reasons", ipp.Keyword("job-completed-successfully")),
	)
}

// job finds the job a request names (job-id, or job-uri), of the
// request's location.
func (h *handler) job(ctx context.Context, q *ippRequest) *Slip {
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
	s, err := h.m.SlipOfJob(ctx, id)
	if err != nil || s.Location != q.loc.Name {
		q.fail(ipp.StatusClientErrorNotFound, "no job %d", id)
		return nil
	}
	return s
}

func ippState(s *Slip) (ipp.JobState, string) {
	if s.State == SlipCanceled {
		return ipp.JobCanceled, "job-canceled-by-user"
	}
	return ipp.JobCompleted, "job-completed-successfully"
}

func (h *handler) jobAttributes(q *ippRequest, s *Slip) []ipp.Attribute {
	st, reason := ippState(s)
	return []ipp.Attribute{
		ipp.Attr("job-id", ipp.Integer(int(s.Job))),
		ipp.Attr("job-uri", ipp.URI(q.printerURI()+"/"+strconv.Itoa(int(s.Job)))),
		ipp.Attr("job-printer-uri", ipp.URI(q.printerURI())),
		ipp.Attr("job-name", ipp.Name(firstOr(s.Name, "slip-"+s.ID))),
		ipp.Attr("job-originating-user-name", ipp.Name(firstOr(s.User, "anonymous"))),
		ipp.Attr("job-state", ipp.Enum(int(st))),
		ipp.Attr("job-state-reasons", ipp.Keyword(reason)),
		ipp.Attr("job-impressions-completed", ipp.Integer(s.Pages)),
		ipp.Attr("time-at-creation", ipp.Integer(int(s.Received.Unix()&0x7fffffff))),
		ipp.Attr("time-at-completed", ipp.Integer(int(s.Received.Unix()&0x7fffffff))),
	}
}

// jobs answers Get-Jobs: the board's slips, as completed (or canceled)
// jobs; which-jobs not-completed has none.
func (h *handler) jobs(ctx context.Context, q *ippRequest) {
	if w := q.text("which-jobs"); w == "not-completed" {
		return
	}
	limit := 1000
	if a := q.attr("limit"); a != nil {
		if n, ok := a.Int(); ok && n > 0 {
			limit = int(n)
		}
	}
	b, err := h.m.Board(ctx, q.loc.Name)
	if err != nil {
		q.fail(ipp.StatusServerErrorInternal, "%v", err)
		return
	}
	all := append(b.Open, b.Ended...)
	for i := range all {
		if i >= limit {
			break
		}
		// a group for each job (Add would join them)
		q.res.Groups = append(q.res.Groups, ipp.Group{Tag: ipp.TagJob, Attrs: h.jobAttributes(q, &all[i])})
	}
}

// submitter finds who sends jobs: with SubmitTokens, a token's name and
// secret (Basic) or the admin token; without, anyone ("").
func (h *handler) submitter(r *http.Request) (string, bool) {
	if h.admin(r) {
		return "admin", true
	}
	if len(h.m.o.SubmitTokens) == 0 {
		return "", true
	}
	name, secret, ok := r.BasicAuth()
	if !ok {
		if tok, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer "); found {
			name, secret, ok = strings.Cut(tok, ":")
		}
	}
	if !ok {
		return "", false
	}
	want, known := h.m.o.SubmitTokens[name]
	if !known || subtle.ConstantTimeCompare([]byte(want), []byte(secret)) != 1 {
		return "", false
	}
	return name, true
}
