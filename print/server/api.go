package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/spool"
)

// routes are the server's paths.
func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.page)
	mux.HandleFunc("GET /app.js", s.page)
	mux.HandleFunc("GET /healthz", s.healthz)
	mux.HandleFunc("GET /metrics", s.authed(s.metrics))

	const v1 = "/api/v1"
	mux.HandleFunc("POST "+v1+"/queues/{queue}/jobs", s.authed(s.submit))
	mux.HandleFunc("GET "+v1+"/jobs", s.authed(s.listJobs))
	mux.HandleFunc("GET "+v1+"/jobs/{id}", s.authed(s.getJob))
	mux.HandleFunc("PATCH "+v1+"/jobs/{id}", s.authed(s.patchJob))
	mux.HandleFunc("POST "+v1+"/jobs/{id}/{op}", s.authed(s.jobOp))
	mux.HandleFunc("GET "+v1+"/queues", s.authed(s.listQueues))
	mux.HandleFunc("GET "+v1+"/queues/{queue}", s.authed(s.getQueue))
	mux.HandleFunc("POST "+v1+"/queues/{queue}/{op}", s.authed(s.queueOp))
	mux.HandleFunc("GET "+v1+"/printers", s.authed(s.listPrinters))
	mux.HandleFunc("GET "+v1+"/printers/{printer}", s.authed(s.getPrinter))
	mux.HandleFunc("GET "+v1+"/discover", s.authed(s.discover))
	mux.HandleFunc("GET "+v1+"/accounting", s.authed(s.accounting))
	mux.HandleFunc("GET "+v1+"/backup", s.authed(s.backup))
	mux.HandleFunc("GET "+v1+"/events", s.authed(s.events))
	mux.HandleFunc(v1+"/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "not-found", "no such path: "+r.URL.Path)
	})
	return mux
}

// authed passes the caller of a request to a handler, or refuses it.
func (s *Server) authed(h func(http.ResponseWriter, *http.Request, *caller)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		c := s.auth.caller(r)
		if c == nil {
			w.Header().Set("WWW-Authenticate", `Bearer realm="bdf"`)
			writeError(w, http.StatusUnauthorized, "unauthorized", "a valid token is required (Authorization: Bearer …)")
			return
		}
		h(w, r, c)
	}
}

// apiError is the body of an error.
type apiError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	var e apiError
	e.Error.Code, e.Error.Message = code, message
	writeJSON(w, status, e)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(append(b, '\n'))
}

// spoolErrors are the spool's errors as HTTP statuses and codes.
var spoolErrors = []struct {
	err    error
	status int
	code   string
}{
	{spool.ErrInvalid, http.StatusBadRequest, "invalid"},
	{spool.ErrBadDocument, http.StatusBadRequest, "bad-document"},
	{spool.ErrUnknownFormat, http.StatusUnsupportedMediaType, "unknown-format"},
	{spool.ErrNotFound, http.StatusNotFound, "not-found"},
	{spool.ErrNoQueue, http.StatusNotFound, "no-queue"},
	{spool.ErrNotAccepting, http.StatusServiceUnavailable, "not-accepting"},
	{spool.ErrFull, http.StatusInsufficientStorage, "full"},
	{spool.ErrTooLarge, http.StatusRequestEntityTooLarge, "too-large"},
	{spool.ErrNotPossible, http.StatusConflict, "not-possible"},
	{spool.ErrWrongPassword, http.StatusForbidden, "wrong-password"},
	{spool.ErrPasswordRequired, http.StatusForbidden, "password-required"},
	{spool.ErrClosed, http.StatusServiceUnavailable, "closed"},
}

// fail answers an error of the spool, or an internal error, which is
// logged and not told.
func (s *Server) fail(w http.ResponseWriter, r *http.Request, err error) {
	for _, e := range spoolErrors {
		if errors.Is(err, e.err) {
			writeError(w, e.status, e.code, err.Error())
			return
		}
	}
	if r.Context().Err() != nil {
		return // the client went away
	}
	s.log.Error("request", "method", r.Method, "path", r.URL.Path, "err", err)
	writeError(w, http.StatusInternalServerError, "internal", "the server failed; its log tells why")
}

func forbidden(w http.ResponseWriter, message string) {
	writeError(w, http.StatusForbidden, "forbidden", message)
}

func badRequest(w http.ResponseWriter, message string) {
	writeError(w, http.StatusBadRequest, "invalid", message)
}

// readBody reads a small JSON body into v; an empty body leaves v.
func readBody(w http.ResponseWriter, r *http.Request, v any) bool {
	b, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 64<<10))
	if err != nil {
		writeError(w, http.StatusRequestEntityTooLarge, "too-large", "the body is larger than 64 KiB")
		return false
	}
	if len(bytes.TrimSpace(b)) == 0 {
		return true
	}
	if err := decodeStrict(b, v); err != nil {
		badRequest(w, "the body: "+err.Error())
		return false
	}
	return true
}

func decodeStrict(b []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// The parts of a submission besides the document are small, and few.
const (
	maxField  = 1 << 20
	maxFields = 16
	// uploadTime is how long a submission's request may take to arrive.
	uploadTime = 30 * time.Minute
)

// submit takes a job: a multipart/form-data request with the document
// (document), and ticket and options (JSON), name, password, priority,
// hold-until and hold-until-time, in any order. The document is received
// as it comes; the job is recorded once the request has arrived.
func (s *Server) submit(w http.ResponseWriter, r *http.Request, c *caller) {
	queue := r.PathValue("queue")
	if !c.can(roleSubmit) {
		forbidden(w, "the token may not submit jobs")
		return
	}
	if !c.mayQueue(queue) {
		forbidden(w, "the token may not submit to queue "+queue)
		return
	}
	// a queue that does not take jobs is told before the document is read
	if q, err := s.sp.Queue(queue); err != nil {
		s.fail(w, r, err)
		return
	} else if q.Rejecting {
		s.fail(w, r, fmt.Errorf("%w: %s", spool.ErrNotAccepting, queue))
		return
	}
	select {
	case s.uploads <- struct{}{}:
		defer func() { <-s.uploads }()
	default:
		w.Header().Set("Retry-After", "5")
		writeError(w, http.StatusServiceUnavailable, "busy", "the server is receiving as many documents as it takes at once")
		return
	}
	http.NewResponseController(w).SetReadDeadline(time.Now().Add(uploadTime))
	mr, err := r.MultipartReader()
	if err != nil {
		badRequest(w, "a submission is multipart/form-data: "+err.Error())
		return
	}
	sub := spool.Submission{Queue: queue, User: c.user}
	defer func() {
		// Submit takes the document; one not submitted is dropped
		if sub.Upload != nil {
			sub.Upload.Discard()
		}
	}()
	var fileName string
	fields := 0
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			badRequest(w, "the multipart body: "+err.Error())
			return
		}
		name := part.FormName()
		if name == "document" {
			if sub.Upload != nil {
				badRequest(w, "more than one document")
				return
			}
			up, err := s.sp.Receive(r.Context(), part)
			if err != nil {
				s.fail(w, r, err)
				return
			}
			sub.Upload, fileName = up, filepath.Base(filepath.FromSlash(part.FileName()))
			continue
		}
		if fields++; fields > maxFields {
			badRequest(w, "too many fields")
			return
		}
		b, err := io.ReadAll(io.LimitReader(part, maxField+1))
		if err != nil {
			badRequest(w, "field "+name+": "+err.Error())
			return
		}
		if len(b) > maxField {
			writeError(w, http.StatusRequestEntityTooLarge, "too-large", "field "+name+" is larger than 1 MiB")
			return
		}
		v := string(b)
		switch name {
		case "ticket":
			err = decodeStrict(b, &sub.Ticket)
		case "options":
			err = decodeStrict(b, &sub.Options)
		case "name":
			sub.Name = v
		case "password":
			sub.Password = v
		case "priority":
			sub.Priority, err = strconv.Atoi(strings.TrimSpace(v))
		case "hold-until":
			sub.HoldUntil = strings.TrimSpace(v)
		case "hold-until-time":
			sub.HoldUntilTime, err = time.Parse(time.RFC3339, strings.TrimSpace(v))
		default:
			err = errors.New("unknown field")
		}
		if err != nil {
			badRequest(w, "field "+name+": "+err.Error())
			return
		}
	}
	if sub.Upload == nil {
		badRequest(w, "no document (the part named document)")
		return
	}
	if sub.Name == "" && fileName != "." && fileName != string(filepath.Separator) {
		sub.Name = fileName
	}
	j, err := s.sp.Submit(r.Context(), sub)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/api/v1/jobs/%d", j.ID))
	writeJSON(w, http.StatusCreated, j)
}

// jobList is a page of jobs; Next is the after of the next page, when
// there may be one.
type jobList struct {
	Jobs []spool.Job `json:"jobs"`
	Next int32       `json:"next,omitempty"`
}

// listJobs lists jobs: those not done by default, state=all for those done
// within the retention of documents too, or state=<job-state>; queue= and
// user= select; after= and limit= page. A caller without manage sees its
// own.
func (s *Server) listJobs(w http.ResponseWriter, r *http.Request, c *caller) {
	q := r.URL.Query()
	f := spool.Filter{Queue: q.Get("queue"), User: q.Get("user")}
	if !c.can(roleManage) {
		if f.User != "" && f.User != c.user {
			writeJSON(w, http.StatusOK, jobList{Jobs: []spool.Job{}})
			return
		}
		f.User = c.user
	}
	match := func(*spool.Job) bool { return true }
	switch st := q.Get("state"); st {
	case "":
	case "all":
		f.Done = true
	default:
		var want ipp.JobState
		if err := want.UnmarshalText([]byte(st)); err != nil {
			badRequest(w, "state: "+err.Error())
			return
		}
		f.Done = want.Done()
		match = func(j *spool.Job) bool { return j.State == want }
	}
	limit := 100
	if v := q.Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 1000 {
			badRequest(w, "limit: 1 to 1000")
			return
		}
		limit = n
	}
	if v := q.Get("after"); v != "" {
		n, err := strconv.ParseInt(v, 10, 32)
		if err != nil || n < 0 {
			badRequest(w, "after: a job id")
			return
		}
		f.After = int32(n)
	}
	out := jobList{Jobs: []spool.Job{}}
	f.Limit = 1000
	for len(out.Jobs) < limit {
		page := s.sp.List(f)
		for i := range page {
			if match(&page[i]) && len(out.Jobs) < limit {
				out.Jobs = append(out.Jobs, page[i])
			}
		}
		if len(page) < f.Limit {
			break
		}
		f.After = page[len(page)-1].ID
	}
	if len(out.Jobs) == limit {
		out.Next = out.Jobs[limit-1].ID
	}
	writeJSON(w, http.StatusOK, out)
}

// job returns the job of a request's path that the caller may see.
func (s *Server) job(w http.ResponseWriter, r *http.Request, c *caller) (spool.Job, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || id <= 0 {
		writeError(w, http.StatusNotFound, "not-found", "no such job: "+r.PathValue("id"))
		return spool.Job{}, false
	}
	j, err := s.sp.Get(r.Context(), int32(id))
	if err == nil && !c.sees(&j) {
		// others' jobs are not told to exist
		err = fmt.Errorf("%w: %d", spool.ErrNotFound, id)
	}
	if err != nil {
		s.fail(w, r, err)
		return spool.Job{}, false
	}
	return j, true
}

func (s *Server) getJob(w http.ResponseWriter, r *http.Request, c *caller) {
	if j, ok := s.job(w, r, c); ok {
		writeJSON(w, http.StatusOK, j)
	}
}

// jobOp does cancel, hold, release ({"password"}), password ({"password"}),
// retry and move ({"queue"}).
func (s *Server) jobOp(w http.ResponseWriter, r *http.Request, c *caller) {
	j, ok := s.job(w, r, c)
	if !ok {
		return
	}
	if !c.acts(&j) {
		forbidden(w, "the token may not change others' jobs")
		return
	}
	var body struct {
		Password string `json:"password"`
		Queue    string `json:"queue"`
	}
	if !readBody(w, r, &body) {
		return
	}
	ctx := r.Context()
	var err error
	switch op := r.PathValue("op"); op {
	case "cancel":
		reason := "job-canceled-by-user"
		if !c.owns(&j) {
			reason = "job-canceled-by-operator"
		}
		j, err = s.sp.Cancel(ctx, j.ID, reason)
	case "hold":
		j, err = s.sp.Hold(ctx, j.ID)
	case "release":
		j, err = s.sp.Release(ctx, j.ID, body.Password)
	case "password":
		if body.Password == "" {
			badRequest(w, "no password")
			return
		}
		j, err = s.sp.SetPassword(ctx, j.ID, body.Password)
	case "retry":
		j, err = s.sp.Retry(ctx, j.ID)
	case "move":
		if body.Queue == "" {
			badRequest(w, "no queue")
			return
		}
		if !c.mayQueue(body.Queue) {
			forbidden(w, "the token may not submit to queue "+body.Queue)
			return
		}
		j, err = s.sp.Move(ctx, j.ID, body.Queue)
	default:
		writeError(w, http.StatusNotFound, "not-found", "no operation "+op+" of jobs")
		return
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, j)
}

// patchJob changes a job's priority ({"priority": n}).
func (s *Server) patchJob(w http.ResponseWriter, r *http.Request, c *caller) {
	j, ok := s.job(w, r, c)
	if !ok {
		return
	}
	if !c.acts(&j) {
		forbidden(w, "the token may not change others' jobs")
		return
	}
	var body struct {
		Priority *int `json:"priority"`
	}
	if !readBody(w, r, &body) {
		return
	}
	if body.Priority == nil {
		badRequest(w, "nothing to change (priority)")
		return
	}
	j, err := s.sp.SetPriority(r.Context(), j.ID, *body.Priority)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, j)
}

func (s *Server) listQueues(w http.ResponseWriter, r *http.Request, c *caller) {
	writeJSON(w, http.StatusOK, map[string]any{"queues": s.sp.Queues()})
}

func (s *Server) getQueue(w http.ResponseWriter, r *http.Request, c *caller) {
	q, err := s.sp.Queue(r.PathValue("queue"))
	if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, q)
}

// queueOp pauses, resumes, rejects or accepts a queue (admin).
func (s *Server) queueOp(w http.ResponseWriter, r *http.Request, c *caller) {
	if !c.can(roleAdmin) {
		forbidden(w, "queues are changed with the admin role")
		return
	}
	name, ctx := r.PathValue("queue"), r.Context()
	var q spool.QueueInfo
	var err error
	switch op := r.PathValue("op"); op {
	case "pause":
		q, err = s.sp.PauseQueue(ctx, name)
	case "resume":
		q, err = s.sp.ResumeQueue(ctx, name)
	case "reject":
		q, err = s.sp.RejectQueue(ctx, name)
	case "accept":
		q, err = s.sp.AcceptQueue(ctx, name)
	default:
		writeError(w, http.StatusNotFound, "not-found", "no operation "+op+" of queues")
		return
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	s.log.Info("queue changed", "queue", name, "op", r.PathValue("op"), "by", c.user)
	writeJSON(w, http.StatusOK, q)
}

// printers returns the printers, without the users and passwords their
// URLs may hold.
func (s *Server) printers() []spool.PrinterInfo {
	ps := s.sp.Printers()
	for i := range ps {
		ps[i].URL = redact(ps[i].URL)
	}
	return ps
}

func (s *Server) listPrinters(w http.ResponseWriter, r *http.Request, c *caller) {
	writeJSON(w, http.StatusOK, map[string]any{"printers": s.printers()})
}

// capabilities are what a printer supports, with IPP's names.
type capabilities struct {
	Formats          []string `json:"document-format-supported,omitempty"`
	Media            []string `json:"media-supported,omitempty"`
	MediaReady       []string `json:"media-ready,omitempty"`
	MediaDefault     string   `json:"media-default,omitempty"`
	MediaSources     []string `json:"media-source-supported,omitempty"`
	MediaTypes       []string `json:"media-type-supported,omitempty"`
	Sides            []string `json:"sides-supported,omitempty"`
	SidesDefault     string   `json:"sides-default,omitempty"`
	ColorModes       []string `json:"print-color-mode-supported,omitempty"`
	ColorModeDefault string   `json:"print-color-mode-default,omitempty"`
	Qualities        []int32  `json:"print-quality-supported,omitempty"`
	Finishings       []int32  `json:"finishings-supported,omitempty"`
	OutputBins       []string `json:"output-bin-supported,omitempty"`
	CopiesMax        int32    `json:"copies-max,omitempty"`
	Location         string   `json:"printer-location,omitempty"`
	Info             string   `json:"printer-info,omitempty"`
}

func capabilitiesOf(p *ipp.Printer) *capabilities {
	if p == nil {
		return nil
	}
	return &capabilities{Formats: p.Formats, Media: p.Media, MediaReady: p.MediaReady, MediaDefault: p.MediaDefault,
		MediaSources: p.MediaSources, MediaTypes: p.MediaTypes, Sides: p.Sides, SidesDefault: p.SidesDefault,
		ColorModes: p.ColorModes, ColorModeDefault: p.ColorModeDefault, Qualities: p.Qualities, Finishings: p.Finishings,
		OutputBins: p.OutputBins, CopiesMax: p.Copies.Hi, Location: p.Location, Info: p.Info}
}

func (s *Server) getPrinter(w http.ResponseWriter, r *http.Request, c *caller) {
	name := r.PathValue("printer")
	for _, p := range s.printers() {
		if p.Name != name {
			continue
		}
		caps, _ := s.sp.PrinterCapabilities(name)
		writeJSON(w, http.StatusOK, struct {
			spool.PrinterInfo
			Capabilities *capabilities `json:"capabilities,omitempty"`
		}{p, capabilitiesOf(caps)})
		return
	}
	writeError(w, http.StatusNotFound, "no-printer", "no such printer: "+name)
}

// discover looks for printers on the network with DNS-SD (admin), for
// timeout= (2s by default, at most 10s).
func (s *Server) discover(w http.ResponseWriter, r *http.Request, c *caller) {
	if !c.can(roleAdmin) {
		forbidden(w, "discovery is for the admin role")
		return
	}
	timeout := 2 * time.Second
	if v := r.URL.Query().Get("timeout"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil || d <= 0 || d > 10*time.Second {
			badRequest(w, "timeout: up to 10s")
			return
		}
		timeout = d
	}
	svcs, err := ipp.Browse(r.Context(), timeout)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	type found struct {
		Name    string            `json:"name"`
		URL     string            `json:"url"`
		Host    string            `json:"host"`
		Port    int               `json:"port"`
		Secure  bool              `json:"secure"`
		Formats []string          `json:"formats,omitempty"`
		TXT     map[string]string `json:"txt,omitempty"`
	}
	out := []found{}
	for _, sv := range svcs {
		out = append(out, found{sv.Name, sv.URL, sv.Host, sv.Port, sv.Secure, sv.Formats(), sv.TXT})
	}
	writeJSON(w, http.StatusOK, map[string]any{"printers": out})
}

// accounting writes the jobs done from= to= (dates, inclusive, or RFC
// 3339 times) as format=jsonl or csv (manage).
func (s *Server) accounting(w http.ResponseWriter, r *http.Request, c *caller) {
	if !c.can(roleManage) {
		forbidden(w, "the accounting is for the manage role")
		return
	}
	q := r.URL.Query()
	from, err := parseTime(q.Get("from"), false)
	if err != nil {
		badRequest(w, "from: "+err.Error())
		return
	}
	to, err := parseTime(q.Get("to"), true)
	if err != nil {
		badRequest(w, "to: "+err.Error())
		return
	}
	format := q.Get("format")
	switch format {
	case "", "jsonl":
		format = "jsonl"
		w.Header().Set("Content-Type", "application/jsonl")
	case "csv":
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	default:
		badRequest(w, "format: jsonl or csv")
		return
	}
	if err := s.sp.Accounting(r.Context(), from, to, format, w); err != nil && r.Context().Err() == nil {
		// the status went with the first line
		s.log.Error("accounting", "err", err)
	}
}

// parseTime reads a date (the start of the day; of the next day for the
// end of a range) or an RFC 3339 time; "" is the zero time.
func parseTime(v string, end bool) (time.Time, error) {
	if v == "" {
		return time.Time{}, nil
	}
	if t, err := time.Parse(time.DateOnly, v); err == nil {
		if end {
			t = t.AddDate(0, 0, 1)
		}
		return t, nil
	}
	return time.Parse(time.RFC3339, v)
}

// backup sends a consistent copy of the spool's file (admin). The copy is
// made to a temporary file first: a slow client does not hold the store.
func (s *Server) backup(w http.ResponseWriter, r *http.Request, c *caller) {
	if !c.can(roleAdmin) {
		forbidden(w, "backups are for the admin role")
		return
	}
	f, err := os.CreateTemp(s.cfg.TempDir, "bdf-backup-*")
	if err != nil {
		s.fail(w, r, err)
		return
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if err := s.sp.Backup(r.Context(), f); err != nil {
		s.fail(w, r, err)
		return
	}
	size, err := f.Seek(0, io.SeekCurrent)
	if err == nil {
		_, err = f.Seek(0, io.SeekStart)
	}
	if err != nil {
		s.fail(w, r, err)
		return
	}
	name := "spool-" + time.Now().UTC().Format("20060102T150405Z") + ".db"
	h := w.Header()
	h.Set("Content-Type", "application/octet-stream")
	h.Set("Content-Disposition", `attachment; filename="`+name+`"`)
	h.Set("Content-Length", strconv.FormatInt(size, 10))
	s.log.Info("backup", "by", c.user, "bytes", size)
	io.Copy(w, f)
}

// healthz answers 200 while the server serves.
func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	select {
	case <-s.done:
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "stopping"})
	default:
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}
