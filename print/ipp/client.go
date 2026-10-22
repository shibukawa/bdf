package ipp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Client speaks IPP to one printer, or to one queue of a CUPS server, at
// a URL. Its zero value is not usable: set URL. The other fields are
// optional.
type Client struct {
	// URL is the printer's: ipp://host[:631]/ipp/print, ipps://… for
	// TLS, or the http:// and https:// forms of the same. A CUPS queue is
	// ipp://localhost:631/printers/<name>.
	URL string
	// HTTP sends the requests. nil uses a client with a dialing timeout
	// and the TLS configuration of TLS.
	HTTP *http.Client
	// TLS verifies the certificate of an ipps printer. nil verifies with
	// the roots of the system.
	TLS TLSPolicy
	// User is the requesting-user-name of the requests: the user of the
	// process when "".
	User string
	// Language is the attributes-natural-language: "en" when "".
	Language string
	// Username and Password answer a printer that asks for Basic
	// authentication.
	Username, Password string
	// Limits bound the decoding of responses: DefaultLimits when nil.
	Limits *Limits
	// Version is the version requests name: 2.0 when zero. A printer
	// that answers server-error-version-not-supported is asked again with
	// 1.1, once.
	Version Version
	// Timeout bounds a request that sends no document: 30 seconds when
	// zero. A request that sends a document is bounded by its context
	// alone.
	Timeout time.Duration

	nextID atomic.Uint32
	// fellBack is set once a printer refused the version; noChunked once
	// it refused a chunked request.
	fellBack  atomic.Bool
	noChunked atomic.Bool
	once      sync.Once
	client    *http.Client
}

// Errors of talking to a printer.
var (
	// ErrAuth is returned when the printer wants a user and password the
	// client does not have, or refuses the ones it has.
	ErrAuth = errors.New("ipp: the printer asks for authentication")
	// ErrProtocol is returned when the answer is not an IPP message.
	ErrProtocol = errors.New("ipp: the answer is not an IPP message")
)

// StatusError is the error of a response whose status is not a
// success.
type StatusError struct {
	Op      Op
	Status  Status
	Message string // status-message, when the printer gave one
	// Unsupported are the attributes of the unsupported-attributes
	// group, when there is one.
	Unsupported []Attribute
}

func (e *StatusError) Error() string {
	s := "ipp: " + e.Op.String() + ": " + e.Status.String()
	if e.Message != "" {
		s += " (" + e.Message + ")"
	}
	if len(e.Unsupported) > 0 {
		names := make([]string, len(e.Unsupported))
		for i, a := range e.Unsupported {
			names[i] = a.Name
		}
		s += ": " + strings.Join(names, ", ")
	}
	return s
}

// IsStatus reports whether err is a StatusError of the status.
func IsStatus(err error, st Status) bool {
	var se *StatusError
	return errors.As(err, &se) && se.Status == st
}

// ParseURL turns a printer URL into the URL requests go to, and the
// printer-uri the requests name. ipp and http become http on port 631
// unless a port is given; ipps and https become https.
func ParseURL(s string) (httpURL, printerURI string, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", "", fmt.Errorf("ipp: %w", err)
	}
	var scheme, ippScheme string
	switch strings.ToLower(u.Scheme) {
	case "ipp", "http":
		scheme, ippScheme = "http", "ipp"
	case "ipps", "https":
		scheme, ippScheme = "https", "ipps"
	default:
		return "", "", fmt.Errorf("ipp: %s: the scheme must be ipp, ipps, http or https", s)
	}
	if u.Host == "" {
		return "", "", fmt.Errorf("ipp: %s: no host", s)
	}
	host := u.Host
	if u.Port() == "" {
		host = net.JoinHostPort(u.Hostname(), "631")
	}
	path := u.Path
	if path == "" {
		path = "/"
	}
	h := url.URL{Scheme: scheme, Host: host, Path: path, RawQuery: u.RawQuery}
	p := url.URL{Scheme: ippScheme, Host: host, Path: path, RawQuery: u.RawQuery}
	return h.String(), p.String(), nil
}

// PrinterURI returns the printer-uri the client names in requests.
func (c *Client) PrinterURI() string {
	_, p, err := ParseURL(c.URL)
	if err != nil {
		return c.URL
	}
	return p
}

func (c *Client) userName() string {
	if c.User != "" {
		return c.User
	}
	if u, err := user.Current(); err == nil && u.Username != "" {
		return u.Username
	}
	if u := os.Getenv("USER"); u != "" {
		return u
	}
	return "bdf"
}

func (c *Client) httpClient() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	c.once.Do(func() {
		host := ""
		if u, err := url.Parse(c.URL); err == nil {
			host = u.Hostname()
		}
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.DialContext = (&net.Dialer{Timeout: 15 * time.Second}).DialContext
		t.ExpectContinueTimeout = 2 * time.Second
		t.TLSClientConfig = c.TLS.config(host)
		c.client = &http.Client{Transport: t}
	})
	return c.client
}

// NewRequest makes a request of the operation with the client's
// charset, language, printer-uri and requesting-user-name.
func (c *Client) NewRequest(op Op) *Message {
	v := c.Version
	if v == (Version{}) {
		v = Version20
	}
	if c.fellBack.Load() {
		v = Version11
	}
	m := NewRequest(op, c.nextID.Add(1), "utf-8", c.Language)
	m.Version = v
	m.Add(TagOperation, Attr("printer-uri", URI(c.PrinterURI())), Attr("requesting-user-name", Name(c.userName())))
	return m
}

// Do sends a request, with the document that follows it in the body
// when doc is not nil, and returns the response. A response whose
// status is not a success is returned with a StatusError. The document
// is streamed: the request body is chunked and waits for the printer's
// 100 Continue before it is sent, so a refused request sends none of
// it.
func (c *Client) Do(ctx context.Context, m *Message, doc io.Reader) (*Message, error) {
	res, err := c.do(ctx, m, doc)
	if err == nil && res.Status() == StatusServerErrorVersionNotSupported && m.Version != Version11 && doc == nil {
		c.fellBack.Store(true)
		m.Version = Version11
		res, err = c.do(ctx, m, nil)
	}
	if err != nil {
		return nil, err
	}
	if st := res.Status(); !st.OK() {
		se := &StatusError{Op: Op(m.Code), Status: st}
		if a := res.Attr(TagOperation, "status-message"); a != nil {
			se.Message = a.Text()
		}
		if g := res.Group(TagUnsupportedGroup); g != nil {
			se.Unsupported = g.Attrs
		}
		return res, se
	}
	return res, nil
}

// DoStream sends a request and returns the response with what follows it
// in the body: the document of Fetch-Document. The caller closes the
// body; a response that is not a success is returned with a StatusError
// and no body.
func (c *Client) DoStream(ctx context.Context, m *Message) (*Message, io.ReadCloser, error) {
	resp, err := c.post(ctx, m, nil)
	if err != nil {
		return nil, nil, err
	}
	res, br, err := DecodeFrom(resp.Body, c.Limits)
	if err != nil {
		resp.Body.Close()
		return nil, nil, err
	}
	if res.RequestID != m.RequestID {
		resp.Body.Close()
		return nil, nil, fmt.Errorf("%w: request id %d answered with %d", ErrProtocol, m.RequestID, res.RequestID)
	}
	if st := res.Status(); !st.OK() {
		resp.Body.Close()
		se := &StatusError{Op: Op(m.Code), Status: st}
		if a := res.Attr(TagOperation, "status-message"); a != nil {
			se.Message = a.Text()
		}
		return res, nil, se
	}
	return res, struct {
		io.Reader
		io.Closer
	}{br, resp.Body}, nil
}

func (c *Client) do(ctx context.Context, m *Message, doc io.Reader) (*Message, error) {
	if doc == nil && c.Timeout >= 0 {
		t := c.Timeout
		if t == 0 {
			t = 30 * time.Second
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, t)
		defer cancel()
	}
	resp, err := c.post(ctx, m, doc)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res, err := Decode(resp.Body, c.Limits)
	if err != nil {
		return nil, err
	}
	if res.RequestID != m.RequestID {
		return nil, fmt.Errorf("%w: request id %d answered with %d", ErrProtocol, m.RequestID, res.RequestID)
	}
	return res, nil
}

// post sends a request and returns the HTTP response of an IPP message,
// or why there is none. A document is streamed (chunked); a printer that
// takes requests of a known length only (HTTP 411, as old ones do) is sent
// it again from a temporary file, and so from then on.
func (c *Client) post(ctx context.Context, m *Message, doc io.Reader) (*http.Response, error) {
	httpURL, _, err := ParseURL(c.URL)
	if err != nil {
		return nil, err
	}
	head, err := m.Marshal()
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return c.send(ctx, httpURL, bytes.NewReader(head), int64(len(head)))
	}
	if c.noChunked.Load() {
		return c.postSized(ctx, httpURL, head, doc)
	}
	read := &countReader{r: doc}
	resp, err := c.send(ctx, httpURL, io.MultiReader(bytes.NewReader(head), read), -1)
	if errors.Is(err, errLengthRequired) && read.n == 0 {
		// the printer refused the request before its body: the document
		// is whole yet
		c.noChunked.Store(true)
		return c.postSized(ctx, httpURL, head, doc)
	}
	return resp, err
}

// postSized sends a request and its document with its length: the
// document is copied to a temporary file first.
func (c *Client) postSized(ctx context.Context, httpURL string, head []byte, doc io.Reader) (*http.Response, error) {
	f, err := os.CreateTemp("", "bdf-ipp-*")
	if err != nil {
		return nil, fmt.Errorf("ipp: %w", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()
	n, err := io.Copy(f, doc)
	if err == nil {
		_, err = f.Seek(0, io.SeekStart)
	}
	if err != nil {
		return nil, fmt.Errorf("ipp: the document: %w", err)
	}
	return c.send(ctx, httpURL, io.MultiReader(bytes.NewReader(head), f), int64(len(head))+n)
}

// errLengthRequired is the HTTP 411 of a printer that takes no chunked
// requests.
var errLengthRequired = fmt.Errorf("%w: HTTP 411 Length Required: the printer takes no chunked requests", ErrProtocol)

// countReader counts what is read of a reader.
type countReader struct {
	r io.Reader
	n int64
}

func (c *countReader) Read(b []byte) (int, error) {
	n, err := c.r.Read(b)
	c.n += int64(n)
	return n, err
}

// send posts a body of a length (-1: chunked, after the printer's 100
// Continue) and checks the response is an IPP message.
func (c *Client) send(ctx context.Context, httpURL string, body io.Reader, length int64) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, httpURL, body)
	if err != nil {
		return nil, fmt.Errorf("ipp: %w", err)
	}
	req.Header.Set("Content-Type", "application/ipp")
	req.Header.Set("Accept", "application/ipp")
	req.Header.Set("User-Agent", "bdf (IPP client)")
	req.ContentLength = length
	if length < 0 {
		req.Header.Set("Expect", "100-continue")
	}
	if c.Username != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("ipp: %w", err)
	}
	fail := func(err error) (*http.Response, error) {
		resp.Body.Close()
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusUnauthorized, http.StatusForbidden:
		return fail(ErrAuth)
	case http.StatusUpgradeRequired:
		return fail(fmt.Errorf("ipp: the printer requires TLS (HTTP 426): use ipps://"))
	case http.StatusLengthRequired:
		return fail(errLengthRequired)
	default:
		return fail(fmt.Errorf("%w: HTTP %s", ErrProtocol, resp.Status))
	}
	if ct, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type")); ct != "application/ipp" {
		return fail(fmt.Errorf("%w: Content-Type %q", ErrProtocol, resp.Header.Get("Content-Type")))
	}
	return resp, nil
}

// JobStatus is what a printer tells of a job.
type JobStatus struct {
	ID          int32
	URI         string
	Name        string
	State       JobState
	Reasons     []string
	Message     string // job-state-message
	Impressions int32  // job-impressions-completed, or -1 when the printer does not say
	// Attrs are the job attributes as received.
	Attrs []Attribute
	// Unsupported are the attributes the printer did not take, as the
	// answer to a request that makes a job listed them (its
	// unsupported-attributes group).
	Unsupported []Attribute
}

func (j *JobStatus) String() string {
	s := "job " + strconv.Itoa(int(j.ID)) + " " + j.State.String()
	if len(j.Reasons) > 0 {
		s += " (" + strings.Join(j.Reasons, ", ") + ")"
	}
	if j.Impressions >= 0 {
		s += " " + strconv.Itoa(int(j.Impressions)) + " impressions"
	}
	return s
}

func jobStatus(attrs []Attribute) *JobStatus {
	j := &JobStatus{Impressions: -1, Attrs: attrs}
	for _, a := range attrs {
		a := a
		switch a.Name {
		case "job-id":
			j.ID, _ = a.Int()
		case "job-uri":
			j.URI = a.Text()
		case "job-name":
			j.Name = a.Text()
		case "job-state":
			n, _ := a.Int()
			j.State = JobState(n)
		case "job-state-reasons":
			j.Reasons = a.Strings()
		case "job-state-message":
			j.Message = a.Text()
		case "job-impressions-completed":
			j.Impressions, _ = a.Int()
		}
	}
	return j
}

// jobAttrsRequested are the job attributes the client asks about.
var jobAttrsRequested = []Value{Keyword("job-id"), Keyword("job-uri"), Keyword("job-name"), Keyword("job-state"),
	Keyword("job-state-reasons"), Keyword("job-state-message"), Keyword("job-impressions-completed")}

// GetPrinterAttributes asks the printer about itself: the attributes
// named, or the ones Printer reads when none are (see
// PrinterAttributes).
func (c *Client) GetPrinterAttributes(ctx context.Context, requested ...string) (*Printer, error) {
	m := c.NewRequest(OpGetPrinterAttributes)
	if len(requested) == 0 {
		requested = PrinterAttributes
	}
	vals := make([]Value, len(requested))
	for i, r := range requested {
		vals[i] = Keyword(r)
	}
	m.Add(TagOperation, Attr("requested-attributes", vals...))
	res, err := c.Do(ctx, m, nil)
	if err != nil {
		return nil, err
	}
	p := ParsePrinter(res)
	p.URI = c.PrinterURI()
	return p, nil
}

// ValidateJob asks whether a job of the attributes, in the document
// format, would be accepted. A nil error means yes; a StatusError says
// what is not supported.
func (c *Client) ValidateJob(ctx context.Context, job []Attribute, format string) error {
	m := c.NewRequest(OpValidateJob)
	if format != "" {
		m.Add(TagOperation, Attr("document-format", MimeType(format)))
	}
	if len(job) > 0 {
		m.Add(TagJob, job...)
	}
	_, err := c.Do(ctx, m, nil)
	return err
}

// CreateJob creates a job of the attributes, to send documents to.
func (c *Client) CreateJob(ctx context.Context, name string, job []Attribute) (*JobStatus, error) {
	m := c.NewRequest(OpCreateJob)
	if name != "" {
		m.Add(TagOperation, Attr("job-name", Name(name)))
	}
	if len(job) > 0 {
		m.Add(TagJob, job...)
	}
	res, err := c.Do(ctx, m, nil)
	if err != nil {
		return nil, err
	}
	return madeJob(res), nil
}

// SendDocument sends a document of the format to a job, streaming doc.
// last tells the printer the job has no more documents. op are more
// operation attributes of the document (document-metadata …).
func (c *Client) SendDocument(ctx context.Context, jobID int32, name, format string, doc io.Reader, last bool, op ...Attribute) (*JobStatus, error) {
	m := c.NewRequest(OpSendDocument)
	m.Add(TagOperation, Attr("job-id", Integer(int(jobID))))
	if name != "" {
		m.Add(TagOperation, Attr("document-name", Name(name)))
	}
	if format != "" {
		m.Add(TagOperation, Attr("document-format", MimeType(format)))
	}
	m.Add(TagOperation, Attr("last-document", Boolean(last)))
	if len(op) > 0 {
		m.Add(TagOperation, op...)
	}
	res, err := c.Do(ctx, m, doc)
	if err != nil {
		return nil, err
	}
	return madeJob(res), nil
}

// PrintJob prints one document as a job of the attributes, streaming
// doc. op are more operation attributes of the document
// (document-metadata …).
func (c *Client) PrintJob(ctx context.Context, name, format string, job []Attribute, doc io.Reader, op ...Attribute) (*JobStatus, error) {
	m := c.NewRequest(OpPrintJob)
	if name != "" {
		m.Add(TagOperation, Attr("job-name", Name(name)))
	}
	if format != "" {
		m.Add(TagOperation, Attr("document-format", MimeType(format)))
	}
	if len(op) > 0 {
		m.Add(TagOperation, op...)
	}
	if len(job) > 0 {
		m.Add(TagJob, job...)
	}
	res, err := c.Do(ctx, m, doc)
	if err != nil {
		return nil, err
	}
	return madeJob(res), nil
}

// madeJob is the status of a job a request made, with what the printer
// did not take.
func madeJob(res *Message) *JobStatus {
	js := jobStatus(groupAttrs(res, TagJob))
	js.Unsupported = groupAttrs(res, TagUnsupportedGroup)
	return js
}

// GetJobAttributes asks about a job.
func (c *Client) GetJobAttributes(ctx context.Context, jobID int32) (*JobStatus, error) {
	m := c.NewRequest(OpGetJobAttributes)
	m.Add(TagOperation, Attr("job-id", Integer(int(jobID))), Attr("requested-attributes", jobAttrsRequested...))
	res, err := c.Do(ctx, m, nil)
	if err != nil {
		return nil, err
	}
	return jobStatus(groupAttrs(res, TagJob)), nil
}

// GetJobs lists jobs: which is "completed", "not-completed" or "all"
// ("" asks for the ones not completed), mine limits them to the user's.
func (c *Client) GetJobs(ctx context.Context, which string, mine bool, limit int) ([]*JobStatus, error) {
	m := c.NewRequest(OpGetJobs)
	if which != "" {
		m.Add(TagOperation, Attr("which-jobs", Keyword(which)))
	}
	if mine {
		m.Add(TagOperation, Attr("my-jobs", Boolean(true)))
	}
	if limit > 0 {
		m.Add(TagOperation, Attr("limit", Integer(limit)))
	}
	m.Add(TagOperation, Attr("requested-attributes", jobAttrsRequested...))
	res, err := c.Do(ctx, m, nil)
	if err != nil {
		return nil, err
	}
	var jobs []*JobStatus
	for i := range res.Groups {
		if res.Groups[i].Tag == TagJob {
			jobs = append(jobs, jobStatus(res.Groups[i].Attrs))
		}
	}
	return jobs, nil
}

// CancelJob cancels a job.
func (c *Client) CancelJob(ctx context.Context, jobID int32) error {
	m := c.NewRequest(OpCancelJob)
	m.Add(TagOperation, Attr("job-id", Integer(int(jobID))))
	_, err := c.Do(ctx, m, nil)
	return err
}

// IdentifyPrinter asks the printer to make itself known: by the
// actions ("flash", "sound", "display", "speak"), or its default.
func (c *Client) IdentifyPrinter(ctx context.Context, actions ...string) error {
	m := c.NewRequest(OpIdentifyPrinter)
	if len(actions) > 0 {
		vals := make([]Value, len(actions))
		for i, a := range actions {
			vals[i] = Keyword(a)
		}
		m.Add(TagOperation, Attr("identify-actions", vals...))
	}
	_, err := c.Do(ctx, m, nil)
	return err
}

// CupsPrinter is a queue of a CUPS server, as CUPS-Get-Printers lists
// them.
type CupsPrinter struct {
	Name      string
	URI       string // printer-uri-supported: the queue's URL
	Info      string
	MakeModel string
	Location  string
	State     PrinterState
	Reasons   []string
	Accepting bool
	IsClass   bool
	DeviceURI string
}

// CupsGetPrinters lists the queues of a CUPS server. The client's URL
// is the server's (ipp://localhost:631/).
func (c *Client) CupsGetPrinters(ctx context.Context) ([]CupsPrinter, error) {
	m := c.NewRequest(OpCupsGetPrinters)
	m.Add(TagOperation, Attr("requested-attributes", Keyword("printer-name"), Keyword("printer-uri-supported"), Keyword("printer-info"),
		Keyword("printer-make-and-model"), Keyword("printer-location"), Keyword("printer-state"), Keyword("printer-state-reasons"),
		Keyword("printer-is-accepting-jobs"), Keyword("printer-type"), Keyword("device-uri")))
	res, err := c.Do(ctx, m, nil)
	if err != nil {
		if IsStatus(err, StatusClientErrorNotFound) {
			return nil, nil
		}
		return nil, err
	}
	var out []CupsPrinter
	for i := range res.Groups {
		g := &res.Groups[i]
		if g.Tag != TagPrinter {
			continue
		}
		p := CupsPrinter{Name: g.Get("printer-name").Text(), Info: g.Get("printer-info").Text(), MakeModel: g.Get("printer-make-and-model").Text(),
			Location: g.Get("printer-location").Text(), Reasons: g.Get("printer-state-reasons").Strings(), DeviceURI: g.Get("device-uri").Text()}
		for _, u := range g.Get("printer-uri-supported").Strings() {
			if p.URI == "" || strings.HasPrefix(u, "ipp:") {
				p.URI = u
			}
		}
		if n, ok := g.Get("printer-state").Int(); ok {
			p.State = PrinterState(n)
		}
		p.Accepting, _ = g.Get("printer-is-accepting-jobs").Bool()
		if t, ok := g.Get("printer-type").Int(); ok {
			p.IsClass = t&1 != 0 // CUPS_PRINTER_CLASS
		}
		out = append(out, p)
	}
	return out, nil
}

func groupAttrs(m *Message, tag Tag) []Attribute {
	if g := m.Group(tag); g != nil {
		return g.Attrs
	}
	return nil
}
