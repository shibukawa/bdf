package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"image"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/ipp/ipptest"
	"github.com/shibukawa/bdf/print/raster"
	"github.com/shibukawa/bdf/print/spool"
)

// rasterConverter converts the rasters of driverless clients, as bdf
// printd's converter does.
type rasterConverter struct{}

func (rasterConverter) Detect(r io.ReaderAt, size int64, name string) (string, bool) {
	head := make([]byte, 8)
	n, _ := r.ReadAt(head, 0)
	switch raster.Sniff(head[:n]) {
	case raster.FormatPWG:
		return "pwg-raster", true
	case raster.FormatURF:
		return "urf", true
	}
	return "", false
}

func (rasterConverter) Convert(ctx context.Context, r io.ReaderAt, size int64, format, name, password string) (*bdf.Document, []string, error) {
	doc, _, err := raster.Document(ctx, io.NewSectionReader(r, 0, size), nil)
	return doc, nil, err
}

// ippFixture is a server whose queues are IPP printers, on a fake printer.
type ippFixture struct {
	*fixture
	ipp *httptest.Server
}

const ippConfigYAML = `
printers:
  mfp1: { url: "%s" }
queues:
  office: { printers: [mfp1], info: "Office", location: "2F" }
  warehouse: { printers: [mfp1] }
ipp: { dns-sd: false, default: office }
`

func newIPPFixture(t *testing.T, yaml string) *ippFixture {
	t.Helper()
	p := ipptest.New()
	p.Formats = append(p.Formats, "application/pdf")
	ps := httptest.NewServer(p)
	t.Cleanup(ps.Close)
	u, _ := url.Parse(ps.URL)
	cfg, err := ParseConfig([]byte(fmt.Sprintf(yaml, "ipp://"+u.Host+"/ipp/print")))
	if err != nil {
		t.Fatal(err)
	}
	cfg.dir = t.TempDir()
	cfg.TempDir = t.TempDir()
	if err := cfg.resolve(); err != nil {
		t.Fatal(err)
	}
	srv, err := New(context.Background(), cfg, Options{Grace: time.Second, Converter: rasterConverter{}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { srv.Close() })
	srv.Spool().Start()
	hs := httptest.NewServer(srv.IPPHandler())
	t.Cleanup(hs.Close)
	api := httptest.NewServer(srv.Handler())
	t.Cleanup(api.Close)
	return &ippFixture{fixture: &fixture{t: t, printer: p, srv: srv, http: api, cfg: cfg}, ipp: hs}
}

// client is an IPP client of a queue, as a user.
func (f *ippFixture) client(queue, user string) *ipp.Client {
	u, _ := url.Parse(f.ipp.URL)
	path := "/ipp/print"
	if queue != "" {
		path += "/" + queue
	}
	return &ipp.Client{URL: "ipp://" + u.Host + path, User: user}
}

// waitJob waits for a job of the spool to be in a state.
func (f *ippFixture) waitJob(id int32, st ipp.JobState) spool.Job {
	f.t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for {
		j, err := f.srv.Spool().Get(context.Background(), id)
		if err == nil && j.State == st {
			return j
		}
		if time.Now().After(deadline) {
			f.t.Fatalf("job %d: %v %v, not %v", id, j.State, err, st)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func statusOf(err error) ipp.Status {
	var se *ipp.StatusError
	if errors.As(err, &se) {
		return se.Status
	}
	return 0
}

// TestIPPPrintJob prints a .bdf with Print-Job: the job is the spool's,
// its user the requesting-user-name, its ticket the job's attributes, and
// it reaches the printer behind.
func TestIPPPrintJob(t *testing.T) {
	f := newIPPFixture(t, ippConfigYAML)
	c := f.client("office", "alice")
	ctx := context.Background()
	job := []ipp.Attribute{
		ipp.Attr("copies", ipp.Integer(2)),
		ipp.Attr("sides", ipp.Keyword("two-sided-long-edge")),
		ipp.Attr("media-col", ipp.Collection(ipp.Attr("media-size", ipp.Collection(
			ipp.Attr("x-dimension", ipp.Integer(21000)), ipp.Attr("y-dimension", ipp.Integer(29700)))))),
		ipp.Attr("page-ranges", ipp.NewRange(1, 1)),
		ipp.Attr("print-color-mode", ipp.Keyword("monochrome")),
		ipp.Attr("job-priority", ipp.Integer(80)),
	}
	js, err := c.PrintJob(ctx, "invoice.bdf", "application/x-bdf", job, bytes.NewReader(docBytes(t)))
	if err != nil {
		t.Fatal(err)
	}
	if js.ID <= 0 || !strings.HasSuffix(js.URI, fmt.Sprintf("/ipp/print/office/%d", js.ID)) {
		t.Fatalf("job %+v", js)
	}
	j := f.waitJob(js.ID, ipp.JobCompleted)
	tk := j.Ticket
	if j.User != "alice" || j.Queue != "office" || j.Name != "invoice.bdf" || j.Document.Type != "application/x-bdf" || tk.Copies != 2 ||
		tk.Sides != "two-sided-long-edge" || tk.Media != "iso_a4_210x297mm" || len(tk.Pages) != 1 || tk.ColorMode != "monochrome" || j.Priority != 80 {
		t.Fatalf("job %+v", j)
	}
	if got := f.printer.Jobs(); len(got) != 1 || got[0].Format != "application/pdf" {
		t.Fatalf("the printer has %d jobs", len(got))
	}
	// the job as IPP tells it
	st, err := c.GetJobAttributes(ctx, js.ID)
	if err != nil || st.State != ipp.JobCompleted || st.Name != "invoice.bdf" {
		t.Fatalf("%+v %v", st, err)
	}
	// the default queue is at /ipp/print too
	if _, err := f.client("", "alice").GetJobAttributes(ctx, js.ID); err != nil {
		t.Fatal(err)
	}
	// and not in the other queue
	if _, err := f.client("warehouse", "alice").GetJobAttributes(ctx, js.ID); statusOf(err) != ipp.StatusClientErrorNotFound {
		t.Fatalf("in another queue: %v", err)
	}
}

// TestIPPCreateJob makes jobs with Create-Job and their documents with
// Send-Document: compressed, in two requests with Close-Job, none, two,
// and too late.
func TestIPPCreateJob(t *testing.T) {
	f := newIPPFixture(t, ippConfigYAML)
	c := f.client("office", "bob")
	ctx := context.Background()
	doc := docBytes(t)

	js, err := c.CreateJob(ctx, "a", nil)
	if err != nil {
		t.Fatal(err)
	}
	if st, _ := c.GetJobAttributes(ctx, js.ID); st == nil || st.State != ipp.JobPending || len(st.Reasons) != 1 || st.Reasons[0] != "job-incoming" {
		t.Fatalf("before the document: %+v", st)
	}
	var gz bytes.Buffer
	zw := gzip.NewWriter(&gz)
	zw.Write(doc)
	zw.Close()
	m := c.NewRequest(ipp.OpSendDocument)
	m.Add(ipp.TagOperation, ipp.Attr("job-id", ipp.Integer(int(js.ID))), ipp.Attr("compression", ipp.Keyword("gzip")),
		ipp.Attr("last-document", ipp.Boolean(true)))
	if _, err := c.Do(ctx, m, &gz); err != nil {
		t.Fatal(err)
	}
	j := f.waitJob(js.ID, ipp.JobCompleted)
	if j.Document.Size != int64(len(doc)) || j.User != "bob" {
		t.Fatalf("%+v", j)
	}
	// a document, then Close-Job
	js, _ = c.CreateJob(ctx, "b", nil)
	if _, err := c.SendDocument(ctx, js.ID, "b.bdf", "", bytes.NewReader(doc), false); err != nil {
		t.Fatal(err)
	}
	if _, err := c.SendDocument(ctx, js.ID, "c.bdf", "", bytes.NewReader(doc), true); statusOf(err) != ipp.StatusServerErrorMultipleDocumentJobsNotSupported {
		t.Fatalf("a second document: %v", err)
	}
	cl := c.NewRequest(ipp.OpCloseJob)
	cl.Add(ipp.TagOperation, ipp.Attr("job-id", ipp.Integer(int(js.ID))))
	if _, err := c.Do(ctx, cl, nil); err != nil {
		t.Fatal(err)
	}
	if j := f.waitJob(js.ID, ipp.JobCompleted); j.Name != "b.bdf" {
		t.Fatalf("%+v", j)
	}
	// Close-Job without a document aborts the job
	js, _ = c.CreateJob(ctx, "d", nil)
	cl = c.NewRequest(ipp.OpCloseJob)
	cl.Add(ipp.TagOperation, ipp.Attr("job-id", ipp.Integer(int(js.ID))))
	c.Do(ctx, cl, nil)
	if st, _ := c.GetJobAttributes(ctx, js.ID); st == nil || st.State != ipp.JobAborted {
		t.Fatalf("closed without a document: %+v", st)
	}
	// another user may not send the document of bob's job
	js, _ = c.CreateJob(ctx, "e", nil)
	if _, err := f.client("office", "mallory").SendDocument(ctx, js.ID, "", "", bytes.NewReader(doc), true); statusOf(err) != ipp.StatusClientErrorNotAuthorized {
		t.Fatalf("another's job: %v", err)
	}
	// a document that does not come in time
	old := incomingTimeout
	incomingTimeout = time.Millisecond
	defer func() { incomingTimeout = old }()
	js, _ = c.CreateJob(ctx, "f", nil)
	time.Sleep(5 * time.Millisecond)
	if st, _ := c.GetJobAttributes(ctx, js.ID); st == nil || st.State != ipp.JobAborted {
		t.Fatalf("too late: %+v", st)
	}
	if _, err := c.SendDocument(ctx, js.ID, "", "", bytes.NewReader(doc), true); statusOf(err) != ipp.StatusClientErrorNotPossible {
		t.Fatalf("a document too late: %v", err)
	}
	// ids are the spool's: the next job is after them all
	js2, err := c.PrintJob(ctx, "g", "", nil, bytes.NewReader(doc))
	if err != nil || js2.ID <= js.ID {
		t.Fatalf("%+v %v", js2, err)
	}
}

// TestIPPAttributes checks what the job attributes ask: those not
// supported are told, and refuse the job with ipp-attribute-fidelity.
func TestIPPAttributes(t *testing.T) {
	f := newIPPFixture(t, ippConfigYAML)
	c := f.client("office", "alice")
	ctx := context.Background()
	bad := []ipp.Attribute{ipp.Attr("media", ipp.Keyword("iso_a0_841x1189mm")), ipp.Attr("x-unknown", ipp.Integer(1))}
	m := c.NewRequest(ipp.OpValidateJob)
	m.Add(ipp.TagOperation, ipp.Attr("ipp-attribute-fidelity", ipp.Boolean(true)))
	m.Add(ipp.TagJob, bad...)
	res, err := c.Do(ctx, m, nil)
	if statusOf(err) != ipp.StatusClientErrorAttributesOrValuesNotSupported {
		t.Fatalf("fidelity: %v", err)
	}
	if g := res.Group(ipp.TagUnsupportedGroup); g == nil || len(g.Attrs) != 2 || g.Get("x-unknown").OutOfBand() != ipp.TagUnsupported {
		t.Fatalf("unsupported: %+v", res.Groups)
	}
	// without fidelity, the job is taken without them
	js, err := c.PrintJob(ctx, "x", "", bad, bytes.NewReader(docBytes(t)))
	if err != nil {
		t.Fatal(err)
	}
	if j := f.waitJob(js.ID, ipp.JobCompleted); j.Ticket.Media != "" {
		t.Fatalf("%+v", j.Ticket)
	}
	for name, a := range map[string]ipp.Attribute{
		"copies":      ipp.Attr("copies", ipp.Integer(0)),
		"sides":       ipp.Attr("sides", ipp.Keyword("both")),
		"page-ranges": ipp.Attr("page-ranges", ipp.NewRange(3, 1)),
		"resolution":  ipp.Attr("printer-resolution", ipp.NewResolution(301, 301, ipp.UnitsDPI)),
		"hold":        ipp.Attr("job-hold-until", ipp.Keyword("tomorrow")),
		"number-up":   ipp.Attr("number-up", ipp.Integer(2)),
	} {
		if err := c.ValidateJob(ctx, []ipp.Attribute{a}, ""); err != nil {
			t.Errorf("%s: without fidelity: %v", name, err)
		}
		m := c.NewRequest(ipp.OpValidateJob)
		m.Add(ipp.TagOperation, ipp.Attr("ipp-attribute-fidelity", ipp.Boolean(true)))
		m.Add(ipp.TagJob, a)
		if _, err := c.Do(ctx, m, nil); statusOf(err) != ipp.StatusClientErrorAttributesOrValuesNotSupported {
			t.Errorf("%s: %v", name, err)
		}
	}
	if err := c.ValidateJob(ctx, nil, "application/vnd.ms-excel"); statusOf(err) != ipp.StatusClientErrorDocumentFormatNotSupported {
		t.Errorf("format: %v", err)
	}
	// a document the server does not print
	if _, err := c.PrintJob(ctx, "x", "", nil, strings.NewReader("hello")); statusOf(err) != ipp.StatusClientErrorDocumentFormatNotSupported {
		t.Errorf("text: %v", err)
	}
	// a queue that does not take jobs
	f.srv.Spool().RejectQueue(ctx, "office")
	if _, err := c.PrintJob(ctx, "x", "", nil, bytes.NewReader(docBytes(t))); statusOf(err) != ipp.StatusServerErrorNotAcceptingJobs {
		t.Errorf("rejecting: %v", err)
	}
	p, err := c.GetPrinterAttributes(ctx)
	if err != nil || p.AcceptingJobs || p.Info != "Office" || p.Location != "2F" || p.UUID == "" || !strings.HasPrefix(p.UUID, "urn:uuid:") {
		t.Fatalf("%+v %v", p, err)
	}
	if !p.SupportsOp(ipp.OpCloseJob) || len(p.URF) == 0 || p.MediaDefault == "" {
		t.Fatalf("%+v", p)
	}
}

// TestIPPRaster prints PWG Raster as a client draws it: the pages are
// pictures of the document sent on.
func TestIPPRaster(t *testing.T) {
	f := newIPPFixture(t, ippConfigYAML)
	c := f.client("office", "alice")
	var buf bytes.Buffer
	w := raster.NewWriter(&buf)
	w.StartPage(raster.Page{Width: 100, Height: 141, XDPI: 12, YDPI: 12, Type: raster.SGray8, PageWidth: 595, PageHeight: 842})
	w.WriteRows(image100x141())
	w.Close()
	js, err := c.PrintJob(context.Background(), "scan", "image/pwg-raster", nil, &buf)
	if err != nil {
		t.Fatal(err)
	}
	j := f.waitJob(js.ID, ipp.JobCompleted)
	if j.Document.Format != "pwg-raster" || j.Media != "iso_a4_210x297mm" {
		t.Fatalf("%+v", j)
	}
}

// TestIPPCancelAndList cancels jobs (one's own, and not another's) and
// lists them, the names of others' hidden.
func TestIPPCancelAndList(t *testing.T) {
	f := newIPPFixture(t, ippConfigYAML)
	ctx := context.Background()
	alice, bob := f.client("office", "alice"), f.client("office", "bob")
	f.srv.Spool().PauseQueue(ctx, "office")
	a, _ := alice.PrintJob(ctx, "alice.bdf", "", nil, bytes.NewReader(docBytes(t)))
	b, _ := bob.PrintJob(ctx, "bob.bdf", "", nil, bytes.NewReader(docBytes(t)))
	if err := bob.CancelJob(ctx, a.ID); statusOf(err) != ipp.StatusClientErrorNotAuthorized {
		t.Fatalf("another's job: %v", err)
	}
	jobs, err := alice.GetJobs(ctx, "not-completed", false, 0)
	if err != nil || len(jobs) != 2 {
		t.Fatalf("%d %v", len(jobs), err)
	}
	for _, js := range jobs {
		if js.ID == b.ID && js.Name != "" || js.ID == a.ID && js.Name != "alice.bdf" {
			t.Errorf("job %d is named %q", js.ID, js.Name)
		}
	}
	if mine, _ := alice.GetJobs(ctx, "not-completed", true, 0); len(mine) != 1 || mine[0].ID != a.ID {
		t.Fatalf("my jobs: %+v", mine)
	}
	if err := alice.CancelJob(ctx, a.ID); err != nil {
		t.Fatal(err)
	}
	m := bob.NewRequest(ipp.OpCancelMyJobs)
	if _, err := bob.Do(ctx, m, nil); err != nil {
		t.Fatal(err)
	}
	f.waitJob(a.ID, ipp.JobCanceled)
	f.waitJob(b.ID, ipp.JobCanceled)
	if done, _ := alice.GetJobs(ctx, "completed", false, 0); len(done) != 2 {
		t.Fatalf("completed: %d", len(done))
	}
	if err := alice.CancelJob(ctx, a.ID); statusOf(err) != ipp.StatusClientErrorNotPossible {
		t.Fatalf("canceled twice: %v", err)
	}
}

// TestIPPBasic asks for a token's name and secret, except to tell what
// the printer is; and refuses ipp from another machine.
func TestIPPBasic(t *testing.T) {
	f := newIPPFixture(t, ippConfigYAML+`
tokens:
  - { name: alice, value: `+aliceToken+`, queues: [office] }
`)
	ctx := context.Background()
	c := f.client("office", "")
	if _, err := c.GetPrinterAttributes(ctx); err != nil {
		t.Fatalf("anonymous Get-Printer-Attributes: %v", err)
	}
	if _, err := c.PrintJob(ctx, "x", "", nil, bytes.NewReader(docBytes(t))); !errors.Is(err, ipp.ErrAuth) {
		t.Fatalf("anonymous Print-Job: %v", err)
	}
	c.Username, c.Password = "alice", aliceToken
	js, err := c.PrintJob(ctx, "x", "", nil, bytes.NewReader(docBytes(t)))
	if err != nil {
		t.Fatal(err)
	}
	if j := f.waitJob(js.ID, ipp.JobCompleted); j.User != "alice" {
		t.Fatalf("user %q", j.User)
	}
	c.Password = "wrong-0123456789abcdef"
	if _, err := c.PrintJob(ctx, "x", "", nil, bytes.NewReader(docBytes(t))); !errors.Is(err, ipp.ErrAuth) {
		t.Fatalf("a wrong secret: %v", err)
	}
	// a queue the token may not print to
	w := f.client("warehouse", "")
	w.Username, w.Password = "alice", aliceToken
	if _, err := w.PrintJob(ctx, "x", "", nil, bytes.NewReader(docBytes(t))); statusOf(err) != ipp.StatusClientErrorNotAuthorized {
		t.Fatalf("another queue: %v", err)
	}
	// ipp from another machine is told to use ipps
	req := httptest.NewRequest("POST", "/ipp/print/office", strings.NewReader(""))
	req.RemoteAddr = "192.0.2.7:5555"
	req.Header.Set("Content-Type", "application/ipp")
	rec := httptest.NewRecorder()
	f.srv.IPPHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusUpgradeRequired {
		t.Fatalf("plain ipp from afar: %d", rec.Code)
	}
}

// TestIPPServeTLS serves ipp and ipps on one port: each is told the URIs
// of its own scheme.
func TestIPPServeTLS(t *testing.T) {
	p := ipptest.New()
	ps := httptest.NewServer(p)
	defer ps.Close()
	u, _ := url.Parse(ps.URL)
	cfg, err := ParseConfig([]byte(fmt.Sprintf(ippConfigYAML, "ipp://"+u.Host+"/ipp/print")))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	cfg.dir, cfg.Spool, cfg.TempDir = dir, dir+"/spool.db", dir
	cfg.Listen = "127.0.0.1:0"
	if err := cfg.resolve(); err != nil {
		t.Fatal(err)
	}
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	ippLn, _ := net.Listen("tcp", "127.0.0.1:0")
	srv, err := New(context.Background(), cfg, Options{Grace: time.Second, IPPListener: ippLn})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- srv.Serve(ctx, ln) }()
	defer func() {
		cancel()
		<-done
	}()
	addr := ippLn.Addr().String()
	for _, scheme := range []string{"ipp", "ipps"} {
		c := &ipp.Client{URL: scheme + "://" + addr + "/ipp/print/office", User: "alice"}
		c.TLS.Insecure = true
		var p *ipp.Printer
		for range 50 {
			if p, err = c.GetPrinterAttributes(ctx); err == nil {
				break
			}
			time.Sleep(20 * time.Millisecond)
		}
		if err != nil {
			t.Fatalf("%s: %v", scheme, err)
		}
		if len(p.URIs) != 2 || !strings.HasPrefix(p.URIs[0], "ipp://") || !strings.HasPrefix(p.URIs[1], "ipps://") {
			t.Fatalf("%s: %v", scheme, p.URIs)
		}
		js, err := c.PrintJob(ctx, "x", "", nil, bytes.NewReader(docBytes(t)))
		if err != nil || !strings.HasPrefix(js.URI, scheme+"://") {
			t.Fatalf("%s: %+v %v", scheme, js, err)
		}
	}
	// the certificate is a leaf, of the server's names
	conn, err := tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		t.Fatal(err)
	}
	if c := conn.ConnectionState().PeerCertificates[0]; c.IsCA || c.KeyUsage&x509.KeyUsageCertSign != 0 || len(c.DNSNames) == 0 {
		t.Errorf("certificate: CA %v, usage %v, names %v", c.IsCA, c.KeyUsage, c.DNSNames)
	}
	conn.Close()
	if srv.ipp.fingerprint == "" {
		t.Fatal("no fingerprint")
	}
	// the icons
	res, err := http.Get("http://" + addr + "/icons/128.png")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode != 200 || !bytes.HasPrefix(b, []byte("\x89PNG")) {
		t.Fatalf("icon: %d", res.StatusCode)
	}
	// the profile of AirPrint, for the address the device came to
	insecure := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	res, err = insecure.Get("http://" + addr + "/airprint.mobileconfig") // (redirected to https)
	if err != nil {
		t.Fatal(err)
	}
	if res.Request.URL.Scheme != "https" {
		t.Errorf("the profile over %s", res.Request.URL.Scheme)
	}
	b, _ = io.ReadAll(res.Body)
	res.Body.Close()
	host, port, _ := net.SplitHostPort(addr)
	for _, want := range []string{"<string>com.apple.airprint</string>", "<string>" + host + "</string>", "<integer>" + port + "</integer>",
		"<string>ipp/print/office</string>", "<string>ipp/print/warehouse</string>", "<key>ForceTLS</key>", "<true/>",
		"<string>com.apple.security.root</string>"} {
		if !bytes.Contains(b, []byte(want)) {
			t.Errorf("the profile has no %s", want)
		}
	}
	if res.Header.Get("Content-Type") != "application/x-apple-aspen-config" {
		t.Errorf("profile: %s", res.Header.Get("Content-Type"))
	}
	// the page of the printers
	res, err = http.Get("http://" + addr + "/ipp/print/office")
	if err != nil {
		t.Fatal(err)
	}
	b, _ = io.ReadAll(res.Body)
	res.Body.Close()
	if !bytes.Contains(b, []byte("Office")) || !bytes.Contains(b, []byte("ipp://"+addr+"/ipp/print/office")) || bytes.Contains(b, []byte("warehouse")) {
		t.Errorf("page:\n%s", b)
	}
}

func image100x141() *image.Gray {
	img := image.NewGray(image.Rect(0, 0, 100, 141))
	for y := 0; y < 141; y++ {
		for x := 0; x < 100; x++ {
			img.Pix[y*img.Stride+x] = uint8((x + y) * 2)
		}
	}
	return img
}

// TestZoneRecords writes the records of wide-area DNS-SD: the browsing
// domain, the types, and for each queue its PTR (and subtypes), SRV and
// TXT, at the host of the domain.
func TestZoneRecords(t *testing.T) {
	cfg, err := ParseConfig([]byte(`
printers: { p: { url: "ipp://192.0.2.1/ipp/print" } }
queues: { office: { printers: [p], info: "Main office" } }
ipp: { host: print.example.com, auth: none }
`))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ZoneRecords(cfg, "example.com", 8631, nil); err != nil {
		t.Fatal(err)
	}
	z, _ := ZoneRecords(cfg, "example.com.", 8631, map[string]*ipp.Printer{"p": {ColorModes: []string{"monochrome"}, Sides: []string{"one-sided"}}})
	for _, want := range []string{
		"b._dns-sd._udp.example.com.\tIN PTR\texample.com.",
		"_universal._sub._ipps._tcp.example.com.\tIN PTR\tMain\\032office\\032\\@\\032print._ipps._tcp.example.com.",
		"Main\\032office\\032\\@\\032print._ipp._tcp.example.com.\tIN SRV\t0 0 8631 print.example.com.",
		`"rp=ipp/print/office"`, `"Color=F"`, `"Duplex=F"`, `"URF=V1.4,`, `"adminurl=https://print.example.com:8631/ipp/print/office"`,
	} {
		if !strings.Contains(z, want) {
			t.Errorf("no %q in\n%s", want, z)
		}
	}
	cfg.IPP.Host = "server.local"
	if _, err := ZoneRecords(cfg, "example.com", 8631, nil); err == nil {
		t.Error("a host of .local in a zone")
	}
}

// TestIPPCancelWhileSending cancels a job while its document arrives: the
// document is dropped, and nothing is printed.
func TestIPPCancelWhileSending(t *testing.T) {
	f := newIPPFixture(t, ippConfigYAML)
	c := f.client("office", "bob")
	ctx := context.Background()
	js, err := c.CreateJob(ctx, "a", nil)
	if err != nil {
		t.Fatal(err)
	}
	pr, pw := io.Pipe()
	done := make(chan error, 1)
	go func() {
		_, err := c.SendDocument(ctx, js.ID, "a.bdf", "", pr, true)
		done <- err
	}()
	doc := docBytes(t)
	pw.Write(doc[:100])
	until(t, "the document arrives", func() bool {
		f.srv.ipp.mu.Lock()
		defer f.srv.ipp.mu.Unlock()
		return f.srv.ipp.incoming[js.ID].sending
	})
	if err := c.CancelJob(ctx, js.ID); err != nil {
		t.Fatal(err)
	}
	pw.Write(doc[100:])
	pw.Close()
	if err := <-done; statusOf(err) != ipp.StatusClientErrorNotPossible {
		t.Fatalf("the document of a canceled job: %v", err)
	}
	left, _ := os.ReadDir(f.cfg.TempDir)
	if len(left) != 0 {
		t.Errorf("left in the temporary directory: %v", left)
	}
	if st, _ := c.GetJobAttributes(ctx, js.ID); st == nil || st.State != ipp.JobCanceled {
		t.Fatalf("%+v", st)
	}
}
