package fakeprinter_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/print/fakeprinter"
	"github.com/shibukawa/bdf/print/ipp"
)

func client(url string) *ipp.Client { return &ipp.Client{URL: url, User: "tester"} }

func TestIPPJobs(t *testing.T) {
	p := fakeprinter.New()
	p.Dir = t.TempDir()
	c := client(fakeprinter.Start(t, p))
	ctx := context.Background()
	st, err := c.PrintJob(ctx, "one", "application/pdf", []ipp.Attribute{ipp.Attr("copies", ipp.Integer(2))}, bytes.NewReader([]byte("%PDF-1.7 one")))
	if err != nil {
		t.Fatal(err)
	}
	if st.ID != 1 {
		t.Errorf("job %d", st.ID)
	}
	j := p.Job(st.ID)
	if j == nil || string(j.Document(0)) != "%PDF-1.7 one" || j.Attr("copies") != "2" || j.User != "tester" || j.Transport != "ipp" {
		t.Fatalf("%+v", j)
	}
	if p.JobState(j) != ipp.JobCompleted {
		t.Errorf("state %s", p.JobState(j))
	}
	if b, err := os.ReadFile(p.Dir + "/1-1.pdf"); err != nil || string(b) != "%PDF-1.7 one" {
		t.Errorf("the saved document: %q %v", b, err)
	}
	// Create-Job, a gzip document, Close-Job
	st, err = c.CreateJob(ctx, "two", nil)
	if err != nil {
		t.Fatal(err)
	}
	var z bytes.Buffer
	zw := gzip.NewWriter(&z)
	zw.Write([]byte("RaS2 page"))
	zw.Close()
	m := c.NewRequest(ipp.OpSendDocument)
	m.Add(ipp.TagOperation, ipp.Attr("job-id", ipp.Integer(int(st.ID))), ipp.Attr("document-format", ipp.MimeType("image/pwg-raster")),
		ipp.Attr("compression", ipp.Keyword("gzip")), ipp.Attr("last-document", ipp.Boolean(false)))
	if _, err := c.Do(ctx, m, &z); err != nil {
		t.Fatal(err)
	}
	if got := string(p.Job(st.ID).Document(0)); got != "RaS2 page" {
		t.Errorf("gzip: %q", got)
	}
	if js, _ := c.GetJobAttributes(ctx, st.ID); js.State != ipp.JobPending {
		t.Errorf("before Close-Job: %s", js.State)
	}
	m = c.NewRequest(ipp.OpCloseJob)
	m.Add(ipp.TagOperation, ipp.Attr("job-id", ipp.Integer(int(st.ID))))
	if _, err := c.Do(ctx, m, nil); err != nil {
		t.Fatal(err)
	}
	if js, _ := c.GetJobAttributes(ctx, st.ID); js.State != ipp.JobCompleted {
		t.Errorf("after Close-Job: %s", js.State)
	}
	// the newest first, completed only when asked
	jobs, err := c.GetJobs(ctx, "completed", false, 0)
	if err != nil || len(jobs) != 2 || jobs[0].ID != 2 {
		t.Errorf("Get-Jobs: %v %v", jobs, err)
	}
	if jobs, _ := c.GetJobs(ctx, "", false, 0); len(jobs) != 0 {
		t.Errorf("not completed: %v", jobs)
	}
}

func TestIPPRefusals(t *testing.T) {
	p := fakeprinter.New()
	p.Reject = []string{"application/pdf"}
	c := client(fakeprinter.Start(t, p))
	ctx := context.Background()
	_, err := c.PrintJob(ctx, "x", "application/pdf", nil, strings.NewReader("%PDF"))
	if !ipp.IsStatus(err, ipp.StatusClientErrorDocumentFormatNotSupported) {
		t.Errorf("a rejected format: %v", err)
	}
	// what the printer does not take is refused, and told
	_, err = c.PrintJob(ctx, "x", "image/jpeg", []ipp.Attribute{ipp.Attr("sides", ipp.Keyword("two-sided-sideways")), ipp.Attr("copies", ipp.Integer(2))}, strings.NewReader("x"))
	var se *ipp.StatusError
	if !errors.As(err, &se) || se.Status != ipp.StatusClientErrorAttributesOrValuesNotSupported || len(se.Unsupported) != 1 || se.Unsupported[0].Name != "sides" {
		t.Errorf("unsupported: %v", err)
	}
	// or ignored by a printer of RFC 8011's kind
	p.IgnoreUnsupported = true
	if _, err := c.PrintJob(ctx, "x", "image/jpeg", []ipp.Attribute{ipp.Attr("sides", ipp.Keyword("two-sided-sideways"))}, strings.NewReader("x")); err != nil {
		t.Errorf("ignored: %v", err)
	}
	p.Busy = 1
	if _, err := c.PrintJob(ctx, "x", "image/jpeg", nil, strings.NewReader("x")); !ipp.IsStatus(err, ipp.StatusServerErrorBusy) {
		t.Errorf("busy: %v", err)
	}
	p.SetAccepting(false)
	if _, err := c.PrintJob(ctx, "x", "image/jpeg", nil, strings.NewReader("x")); !ipp.IsStatus(err, ipp.StatusServerErrorNotAcceptingJobs) {
		t.Errorf("not accepting: %v", err)
	}
	p.SetAccepting(true)
	p.Username, p.Password = "alice", "secret"
	if _, err := c.GetPrinterAttributes(ctx); err == nil {
		t.Error("no password")
	}
	c.Username, c.Password = "alice", "secret"
	if _, err := c.GetPrinterAttributes(ctx); err != nil {
		t.Errorf("with the password: %v", err)
	}
}

func TestIPPOldPrinter(t *testing.T) {
	ps, _ := fakeprinter.LookupPreset("legacy-ipp11")
	p := ps.New("old")
	c := client(fakeprinter.Start(t, p))
	// IPP/2.0 falls back to 1.1
	pr, err := c.GetPrinterAttributes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if pr.MakeModel != "BDF Fake Old Printer" {
		t.Errorf("%+v", pr)
	}
	// a chunked document is refused with 411
	if _, err := c.PrintJob(context.Background(), "x", "application/pdf", nil, strings.NewReader("%PDF")); err == nil {
		t.Error("a chunked request was taken")
	}
}

func TestIPPAttributes(t *testing.T) {
	ps, _ := fakeprinter.LookupPreset("office-mfp")
	p := ps.New("mfp")
	c := client(fakeprinter.Start(t, p))
	pr, err := c.GetPrinterAttributes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if pr.Name != "mfp" || !pr.Supports("image/urf") || len(pr.Media) != 10 || !strings.Contains(strings.Join(pr.Sides, ","), "two-sided-long-edge") {
		t.Errorf("%+v", pr)
	}
	// only what was asked, and media-col-database only by name
	res, err := c.Do(context.Background(), func() *ipp.Message {
		m := c.NewRequest(ipp.OpGetPrinterAttributes)
		m.Add(ipp.TagOperation, ipp.Attr("requested-attributes", ipp.Keyword("printer-state"), ipp.Keyword("job-template")))
		return m
	}(), nil)
	if err != nil {
		t.Fatal(err)
	}
	g := res.Group(ipp.TagPrinter)
	if g.Get("printer-state") == nil || g.Get("copies-supported") == nil || g.Get("printer-name") != nil || g.Get("media-col-database") != nil {
		t.Errorf("requested-attributes: %v", g.Attrs)
	}
}

func TestSocketAndLPD(t *testing.T) {
	ps, _ := fakeprinter.LookupPreset("zebra-zd421-203")
	p := ps.New("label")
	l := fakeprinter.StartSocket(t, p)
	send := func(addr string, data string) string {
		c, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		defer c.Close()
		c.Write([]byte(data))
		c.(*net.TCPConn).CloseWrite()
		b, _ := io.ReadAll(c)
		return string(b)
	}
	// a status query is answered while the job comes
	back := send(l.Addr(), "~HS^XA^PW10^LL1^XZ")
	if !strings.HasPrefix(back, "\x02030,0,0,") {
		t.Errorf("~HS: %q", back)
	}
	jobs := p.Jobs()
	if len(jobs) != 1 || jobs[0].Format != "zpl" || jobs[0].Transport != "socket" {
		t.Fatalf("%+v", jobs)
	}
	p.Stop("media-empty-error")
	if back := send(l.Addr(), "~HS"); !strings.HasPrefix(back, "\x02030,1,1,") {
		t.Errorf("~HS of a printer out of paper: %q", back)
	}
	p.Resume()
	// ESC/POS: DLE EOT 4, the paper sensors
	rs, _ := fakeprinter.LookupPreset("epson-tm-t88")
	receipt := rs.New("receipt")
	rl := fakeprinter.StartSocket(t, receipt)
	if back := send(rl.Addr(), "\x10\x04\x04"); back != "\x12" {
		t.Errorf("DLE EOT 4: %q", back)
	}
	receipt.Stop("media-empty-error")
	if back := send(rl.Addr(), "\x10\x04\x04"); back != "\x72" {
		t.Errorf("DLE EOT 4 out of paper: %q", back)
	}
	// a connection that sends nothing is no job; a cut job keeps nothing
	c, _ := net.Dial("tcp", l.Addr())
	c.Close()
	p.CutAfter(3)
	send(l.Addr(), "^XA^PW10^LL1^XZ")
	time.Sleep(50 * time.Millisecond)
	if n := len(p.Jobs()); n != 2 {
		t.Errorf("%d jobs", n)
	}

	lpd := fakeprinter.New()
	lq := fakeprinter.StartLPD(t, lpd, "lp")
	// CUPS's order: the control file first, then the data
	conn, err := net.Dial("tcp", lq.Addr())
	if err != nil {
		t.Fatal(err)
	}
	ack := func() {
		b := make([]byte, 1)
		if _, err := conn.Read(b); err != nil || b[0] != 0 {
			t.Fatalf("no ack: %v %v", b, err)
		}
	}
	io.WriteString(conn, "\x02lp\n")
	ack()
	ctl := "Hhost\nPalice\nJreport\nldfA001host\n"
	io.WriteString(conn, "\x02"+strconv.Itoa(len(ctl))+" cfA001host\n")
	ack()
	io.WriteString(conn, ctl+"\x00")
	ack()
	io.WriteString(conn, "\x03"+strconv.Itoa(5)+" dfA001host\n")
	ack()
	io.WriteString(conn, "%!PS\n\x00")
	ack()
	conn.Close()
	time.Sleep(100 * time.Millisecond)
	jobs = lpd.Jobs()
	if len(jobs) != 1 || jobs[0].User != "alice" || jobs[0].Name != "report" || string(jobs[0].Data()) != "%!PS\n" || jobs[0].Format != "application/postscript" {
		t.Fatalf("LPD: %+v", jobs)
	}
	lpd.Refuse(true)
	conn, _ = net.Dial("tcp", lq.Addr())
	io.WriteString(conn, "\x02lp\n")
	b := make([]byte, 1)
	conn.Read(b)
	conn.Close()
	if b[0] == 0 {
		t.Error("a refused queue took a job")
	}
}

func TestSerial(t *testing.T) {
	ps, _ := fakeprinter.LookupPreset("epson-tm-t88")
	p := ps.New("receipt")
	name := fakeprinter.StartSerial(t, p)
	f, err := os.OpenFile(name, os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	// (the status a serial printer says back is tested over 9100: the
	// near end of the terminal is not raw here)
	f.Write([]byte("\x1b@\x1dP\xb4\xb4\x1bJ\x10"))
	// the job ends when the line is quiet
	for range 50 {
		if len(p.Jobs()) == 1 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if jobs := p.Jobs(); len(jobs) != 1 || jobs[0].Format != "escpos" || jobs[0].Transport != "serial" {
		t.Fatalf("%+v", jobs)
	}
}

func TestFaults(t *testing.T) {
	p := fakeprinter.New()
	c := client(fakeprinter.Start(t, p))
	ctx := context.Background()
	print := func() (*ipp.JobStatus, error) {
		return c.PrintJob(ctx, "x", "image/jpeg", nil, strings.NewReader("x"))
	}
	// out of paper after the second job, for a while
	p.Inject(fakeprinter.Fault{AfterJobs: 2, Stop: []string{"media-empty-error"}, For: 200 * time.Millisecond})
	print()
	if pr, _ := c.GetPrinterAttributes(ctx); pr.State != ipp.PrinterIdle {
		t.Fatalf("after one job: %s", pr.State)
	}
	print()
	pr, _ := c.GetPrinterAttributes(ctx)
	if pr.State != ipp.PrinterStopped || strings.Join(pr.Reasons, ",") != "media-empty-error" {
		t.Fatalf("after two: %s %v", pr.State, pr.Reasons)
	}
	time.Sleep(300 * time.Millisecond)
	if pr, _ := c.GetPrinterAttributes(ctx); pr.State != ipp.PrinterIdle {
		t.Fatalf("after the paper came: %s", pr.State)
	}
	// a restart forgets the jobs and counts from 1
	p.Inject(fakeprinter.Fault{Restart: true})
	st, _ := print()
	if st.ID != 1 {
		t.Errorf("after a restart: job %d", st.ID)
	}
	// the next job is aborted; then one is cut
	p.Inject(fakeprinter.Fault{Abort: true})
	st, _ = print()
	if js, _ := c.GetJobAttributes(ctx, st.ID); js.State != ipp.JobAborted {
		t.Errorf("aborted: %s", js.State)
	}
	cut := int64(0)
	p.Inject(fakeprinter.Fault{CutAfter: &cut})
	if _, err := print(); err == nil {
		t.Error("a cut job")
	}
	// jobs that never complete
	never := true
	p.Inject(fakeprinter.Fault{NeverComplete: &never})
	st, _ = print()
	if js, _ := c.GetJobAttributes(ctx, st.ID); js.State != ipp.JobProcessing {
		t.Errorf("never complete: %s", js.State)
	}
}

func TestCertificate(t *testing.T) {
	p := fakeprinter.New()
	url := fakeprinter.StartTLS(t, p)
	fingerprint := func() [32]byte {
		conn, err := tls.Dial("tcp", strings.TrimSuffix(strings.TrimPrefix(url, "ipps://"), "/ipp/print"), &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		return sha256.Sum256(conn.ConnectionState().PeerCertificates[0].Raw)
	}
	before := fingerprint()
	if p.Certificate() == nil || sha256.Sum256(p.Certificate().Raw) != before {
		t.Fatal("Certificate is not the one presented")
	}
	p.Inject(fakeprinter.Fault{NewCertificate: true})
	if fingerprint() == before {
		t.Error("the certificate is the same")
	}
	p.Inject(fakeprinter.Fault{ExpiredCertificate: true})
	if !p.Certificate().NotAfter.Before(time.Now()) {
		t.Error("not expired")
	}
	_ = http.StatusOK
}

func TestCUPS(t *testing.T) {
	c := fakeprinter.CUPSPreset()
	base := fakeprinter.StartCUPS(t, c)
	ctx := context.Background()
	server := client(base + "/")
	printers, err := server.CupsGetPrinters(ctx)
	if err != nil || len(printers) != 3 || printers[0].Name != "office" || !strings.HasSuffix(printers[0].URI, "/printers/office") {
		t.Fatalf("%+v %v", printers, err)
	}
	// a class sends its jobs to its members in turn; the ids are the
	// server's
	floor := client(base + "/classes/floor")
	a, err := floor.PrintJob(ctx, "a", "image/pwg-raster", nil, strings.NewReader("RaS2"))
	if err != nil {
		t.Fatal(err)
	}
	b, err := floor.PrintJob(ctx, "b", "image/pwg-raster", nil, strings.NewReader("RaS2"))
	if err != nil {
		t.Fatal(err)
	}
	if a.ID == b.ID || len(c.Printer("office").Jobs())+len(c.Printer("laser").Jobs()) != 2 || len(c.Printer("office").Jobs()) != 1 {
		t.Errorf("jobs %d %d: office %d, laser %d", a.ID, b.ID, len(c.Printer("office").Jobs()), len(c.Printer("laser").Jobs()))
	}
	if js, err := floor.GetJobAttributes(ctx, b.ID); err != nil || js.ID != b.ID {
		t.Errorf("a job of the class: %+v %v", js, err)
	}
	// a stopped member is passed over
	c.Printer("office").Stop("media-jam-error")
	for range 2 {
		if _, err := floor.PrintJob(ctx, "c", "image/pwg-raster", nil, strings.NewReader("RaS2")); err != nil {
			t.Fatal(err)
		}
	}
	if n := len(c.Printer("office").Jobs()); n != 1 {
		t.Errorf("the stopped printer got %d jobs", n)
	}
	// CUPS-Move-Job
	office := client(base + "/printers/office")
	held, err := office.CreateJob(ctx, "held", nil)
	if err != nil {
		t.Fatal(err)
	}
	m := server.NewRequest(ipp.OpCupsMoveJob)
	m.Add(ipp.TagOperation, ipp.Attr("job-id", ipp.Integer(int(held.ID))))
	m.Add(ipp.TagJob, ipp.Attr("job-printer-uri", ipp.URI(base+"/printers/ps")))
	if _, err := server.Do(ctx, m, nil); err != nil {
		t.Fatal(err)
	}
	if c.Printer("ps").Job(held.ID) == nil || c.Printer("office").Job(held.ID) != nil {
		t.Error("the job did not move")
	}
	// CUPS-Get-Default, CUPS-Reject-Jobs
	m = server.NewRequest(ipp.OpCupsGetDefault)
	res, err := server.Do(ctx, m, nil)
	if err != nil || res.Group(ipp.TagPrinter).Get("printer-name").Text() != "office" {
		t.Errorf("default: %v", err)
	}
	m = client(base + "/printers/laser").NewRequest(ipp.OpCupsRejectJobs)
	if _, err := server.Do(ctx, m, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := client(base+"/printers/laser").PrintJob(ctx, "x", "image/pwg-raster", nil, strings.NewReader("RaS2")); !ipp.IsStatus(err, ipp.StatusServerErrorNotAcceptingJobs) {
		t.Errorf("a rejecting printer: %v", err)
	}
}

func TestSniff(t *testing.T) {
	for in, want := range map[string]string{
		"%PDF-1.7": "application/pdf", "%!PS-Adobe-3.0": "application/postscript", "RaS2": "image/pwg-raster", "UNIRAST\x00": "image/urf",
		"^XA^PW": "zpl", "\x02\x1bA\x1bA1": "sbpl", "{D0791,": "tpcl", "\x1b@\x1dP\xcb": "escpos", "\x1b@\x1b3\x18": "escp", "\x1bE": "pcl",
		"\x1b%-12345X@PJL\r\n@PJL ENTER LANGUAGE=PDF\r\n%PDF-1.7": "application/pdf", "hello": "",
	} {
		if got := fakeprinter.Sniff([]byte(in)); got != want {
			t.Errorf("%q: %q, want %q", in, got, want)
		}
	}
}
