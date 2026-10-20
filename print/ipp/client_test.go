package ipp_test

import (
	"bytes"
	"context"
	"errors"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/print/fakeprinter"
	"github.com/shibukawa/bdf/print/ipp"
)

func serve(t *testing.T, p *fakeprinter.Printer) *ipp.Client {
	t.Helper()
	return &ipp.Client{URL: fakeprinter.Start(t, p), User: "tester"}
}

func TestParseURL(t *testing.T) {
	for in, want := range map[string][2]string{
		"ipp://printer.local/ipp/print":        {"http://printer.local:631/ipp/print", "ipp://printer.local:631/ipp/print"},
		"ipps://192.0.2.1/ipp/print":           {"https://192.0.2.1:631/ipp/print", "ipps://192.0.2.1:631/ipp/print"},
		"ipp://localhost:8631/printers/office": {"http://localhost:8631/printers/office", "ipp://localhost:8631/printers/office"},
		"http://[fe80::1%25en0]:631/ipp/print": {"http://[fe80::1%25en0]:631/ipp/print", "ipp://[fe80::1%25en0]:631/ipp/print"},
		"ipp://host":                           {"http://host:631/", "ipp://host:631/"},
	} {
		h, p, err := ipp.ParseURL(in)
		if err != nil || h != want[0] || p != want[1] {
			t.Errorf("ParseURL(%q) = %q, %q, %v; want %q, %q", in, h, p, err, want[0], want[1])
		}
	}
	for _, bad := range []string{"lpd://host/queue", "ipp:///nohost", "://"} {
		if _, _, err := ipp.ParseURL(bad); err == nil {
			t.Errorf("ParseURL(%q) succeeded", bad)
		}
	}
}

func TestPrinterAttributes(t *testing.T) {
	p := fakeprinter.New()
	c := serve(t, p)
	pr, err := c.GetPrinterAttributes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if pr.Name != "fakeprinter" || pr.MakeModel != "BDF Fake Printer" || pr.State != ipp.PrinterIdle || !pr.AcceptingJobs {
		t.Fatalf("identity: %+v", pr)
	}
	if !pr.Supports("application/pdf") || !pr.Supports("image/pwg-raster") || pr.Supports("image/urf") {
		t.Fatalf("formats: %v", pr.Formats)
	}
	if pr.Copies != (ipp.Range{Lo: 1, Hi: 99}) || !pr.PageRanges || pr.PWGRasterSheetBack != "normal" {
		t.Fatalf("copies %v, page ranges %v, sheet back %q", pr.Copies, pr.PageRanges, pr.PWGRasterSheetBack)
	}
	if len(pr.Resolutions) != 2 || pr.Resolutions[1] != (ipp.Resolution{X: 600, Y: 600, Units: ipp.UnitsDPI}) {
		t.Fatalf("resolutions: %v", pr.Resolutions)
	}
	if len(pr.PWGRasterTypes) != 3 || pr.Sides[1] != "two-sided-long-edge" || pr.Media[2] != "jis_b5_182x257mm" {
		t.Fatalf("types %v, sides %v, media %v", pr.PWGRasterTypes, pr.Sides, pr.Media)
	}
	if !pr.SupportsOp(ipp.OpCreateJob) || !pr.SupportsOp(ipp.OpHoldJob) || pr.SupportsOp(ipp.OpPrintURI) {
		t.Fatalf("operations: %v", pr.Operations)
	}
	if !strings.Contains(strings.Join(pr.Summary(), "\n"), "formats: application/pdf, image/pwg-raster") {
		t.Fatalf("summary: %v", pr.Summary())
	}
	if p.Requests["Get-Printer-Attributes"] != 1 {
		t.Fatalf("requests: %v", p.Requests)
	}
}

func TestCreateSendWait(t *testing.T) {
	p := fakeprinter.New()
	p.Processing = 50 * time.Millisecond
	c := serve(t, p)
	ctx := context.Background()
	job := []ipp.Attribute{ipp.Attr("copies", ipp.Integer(2)), ipp.Attr("sides", ipp.Keyword("two-sided-long-edge"))}
	if err := c.ValidateJob(ctx, job, "application/pdf"); err != nil {
		t.Fatal(err)
	}
	js, err := c.CreateJob(ctx, "invoice", job)
	if err != nil {
		t.Fatal(err)
	}
	if js.ID != 1 || js.State != ipp.JobPending {
		t.Fatalf("created: %v", js)
	}
	doc := strings.Repeat("%PDF-1.7 a page\n", 10000) // larger than a buffer: streamed
	js, err = c.SendDocument(ctx, js.ID, "invoice.pdf", "application/pdf", strings.NewReader(doc), true)
	if err != nil {
		t.Fatal(err)
	}
	if js.State != ipp.JobProcessing {
		t.Fatalf("after the document: %v", js)
	}
	time.Sleep(80 * time.Millisecond)
	js, err = c.GetJobAttributes(ctx, js.ID)
	if err != nil {
		t.Fatal(err)
	}
	if js.State != ipp.JobCompleted || js.Impressions != 2 || len(js.Reasons) != 1 || js.Reasons[0] != "job-completed-successfully" {
		t.Fatalf("done: %v", js)
	}
	j := p.Job(js.ID)
	if j.Name != "invoice" || j.User != "tester" || j.Format != "application/pdf" || string(j.Document(0)) != doc {
		t.Fatalf("the printer's job: %q %q %q %d bytes", j.Name, j.User, j.Format, len(j.Document(0)))
	}
	jobs, err := c.GetJobs(ctx, "completed", true, 0)
	if err != nil || len(jobs) != 1 || jobs[0].ID != js.ID {
		t.Fatalf("Get-Jobs: %v %v", jobs, err)
	}
	if jobs, _ := c.GetJobs(ctx, "", false, 0); len(jobs) != 0 {
		t.Fatalf("not completed: %v", jobs)
	}
}

func TestPrintJobAndCancel(t *testing.T) {
	p := fakeprinter.New()
	p.Dir = t.TempDir()
	p.Processing = time.Hour
	c := serve(t, p)
	ctx := context.Background()
	js, err := c.PrintJob(ctx, "slip", "image/pwg-raster", nil, bytes.NewReader([]byte("RaS2")))
	if err != nil {
		t.Fatal(err)
	}
	if js.State != ipp.JobProcessing {
		t.Fatalf("Print-Job: %v", js)
	}
	if matches, _ := filepath.Glob(filepath.Join(p.Dir, "1-1.pwg")); len(matches) != 1 {
		t.Fatalf("the document was not saved in %s", p.Dir)
	}
	if err := c.CancelJob(ctx, js.ID); err != nil {
		t.Fatal(err)
	}
	if js, _ := c.GetJobAttributes(ctx, js.ID); js.State != ipp.JobCanceled {
		t.Fatalf("after Cancel-Job: %v", js)
	}
	if err := c.CancelJob(ctx, js.ID); !ipp.IsStatus(err, ipp.StatusClientErrorNotPossible) {
		t.Fatalf("cancelling twice: %v", err)
	}
	if _, err := c.GetJobAttributes(ctx, 99); !ipp.IsStatus(err, ipp.StatusClientErrorNotFound) {
		t.Fatalf("no such job: %v", err)
	}
}

func TestUnsupported(t *testing.T) {
	p := fakeprinter.New()
	c := serve(t, p)
	ctx := context.Background()
	err := c.ValidateJob(ctx, []ipp.Attribute{ipp.Attr("copies", ipp.Integer(500)), ipp.Attr("media", ipp.Keyword("jpn_hagaki_100x148mm")),
		ipp.Attr("sides", ipp.Keyword("one-sided"))}, "application/pdf")
	var se *ipp.StatusError
	if !errors.As(err, &se) || se.Status != ipp.StatusClientErrorAttributesOrValuesNotSupported || len(se.Unsupported) != 2 {
		t.Fatalf("unsupported attributes: %v", err)
	}
	if se.Unsupported[0].Name != "copies" || se.Unsupported[1].Name != "media" || !strings.Contains(err.Error(), "copies, media") {
		t.Fatalf("unsupported: %v", se)
	}
	if err := c.ValidateJob(ctx, nil, "image/urf"); !ipp.IsStatus(err, ipp.StatusClientErrorDocumentFormatNotSupported) {
		t.Fatalf("format: %v", err)
	}
	p.Reject = []string{"application/pdf"}
	js, err := c.CreateJob(ctx, "x", nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.SendDocument(ctx, js.ID, "", "application/pdf", strings.NewReader("pdf"), true); !ipp.IsStatus(err, ipp.StatusClientErrorDocumentFormatNotSupported) {
		t.Fatalf("a listed format refused: %v", err)
	}
	p.Accepting = false
	if _, err := c.CreateJob(ctx, "x", nil); !ipp.IsStatus(err, ipp.StatusServerErrorNotAcceptingJobs) {
		t.Fatalf("not accepting: %v", err)
	}
}

func TestAuthAndVersion(t *testing.T) {
	p := fakeprinter.New()
	p.Username, p.Password = "alice", "secret"
	c := serve(t, p)
	ctx := context.Background()
	if _, err := c.GetPrinterAttributes(ctx); !errors.Is(err, ipp.ErrAuth) {
		t.Fatalf("without a password: %v", err)
	}
	c.Username, c.Password = "alice", "wrong"
	if _, err := c.GetPrinterAttributes(ctx); !errors.Is(err, ipp.ErrAuth) {
		t.Fatalf("a wrong password: %v", err)
	}
	c.Password = "secret"
	if _, err := c.GetPrinterAttributes(ctx); err != nil {
		t.Fatalf("with the password: %v", err)
	}
	q := fakeprinter.New()
	q.Only11 = true
	c11 := serve(t, q)
	if _, err := c11.GetPrinterAttributes(ctx); err != nil {
		t.Fatalf("falling back to 1.1: %v", err)
	}
	if q.Requests["Get-Printer-Attributes"] != 2 {
		t.Fatalf("requests: %v", q.Requests)
	}
	if _, err := c11.GetJobs(ctx, "", false, 0); err != nil {
		t.Fatalf("1.1 kept: %v", err)
	}
	if q.Requests["Get-Jobs"] != 1 {
		t.Fatalf("the version was asked again: %v", q.Requests)
	}
}

func TestCupsGetPrinters(t *testing.T) {
	c := serve(t, fakeprinter.New())
	ps, err := c.CupsGetPrinters(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 1 || ps[0].Name != "fakeprinter" || !ps[0].Accepting || !strings.HasSuffix(ps[0].URI, "/ipp/print") {
		t.Fatalf("queues: %+v", ps)
	}
}

func TestNotIPP(t *testing.T) {
	srv := httptest.NewServer(nil) // 404 for everything
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	c := &ipp.Client{URL: "ipp://" + u.Host + "/ipp/print"}
	if _, err := c.GetPrinterAttributes(context.Background()); !errors.Is(err, ipp.ErrProtocol) {
		t.Fatalf("a web server: %v", err)
	}
}

func TestTrustOnFirstUse(t *testing.T) {
	ctx := context.Background()
	store := filepath.Join(t.TempDir(), "printers.json")
	p := fakeprinter.New()
	url1 := fakeprinter.StartTLS(t, p)
	if _, err := (&ipp.Client{URL: url1}).GetPrinterAttributes(ctx); err == nil {
		t.Fatal("a self-signed certificate was taken with the roots of the system")
	}
	if _, err := (&ipp.Client{URL: url1, TLS: ipp.TLSPolicy{Store: &ipp.TrustStore{Path: store, ReadOnly: true}}}).GetPrinterAttributes(ctx); err == nil {
		t.Fatal("a read-only store took a certificate it has not recorded")
	}
	if _, err := (&ipp.Client{URL: url1, TLS: ipp.TLSPolicy{Store: &ipp.TrustStore{Path: store}}}).GetPrinterAttributes(ctx); err != nil {
		t.Fatalf("first use: %v", err)
	}
	if _, err := (&ipp.Client{URL: url1, TLS: ipp.TLSPolicy{Store: &ipp.TrustStore{Path: store, ReadOnly: true}}}).GetPrinterAttributes(ctx); err != nil {
		t.Fatalf("the recorded certificate: %v", err)
	}
	// The same printer with another certificate: a printer replaced, or
	// someone in between.
	p.Inject(fakeprinter.Fault{NewCertificate: true})
	url2 := url1
	_, err := (&ipp.Client{URL: url2, TLS: ipp.TLSPolicy{Store: &ipp.TrustStore{Path: store}}}).GetPrinterAttributes(ctx)
	if !errors.Is(err, ipp.ErrCertificateChanged) {
		t.Fatalf("a changed certificate: %v", err)
	}
	st := &ipp.TrustStore{Path: store}
	if err := st.Forget("127.0.0.1"); err != nil {
		t.Fatal(err)
	}
	if _, err := (&ipp.Client{URL: url2, TLS: ipp.TLSPolicy{Store: &ipp.TrustStore{Path: store}}}).GetPrinterAttributes(ctx); err != nil {
		t.Fatalf("after Forget: %v", err)
	}
	if _, err := (&ipp.Client{URL: url1, TLS: ipp.TLSPolicy{Insecure: true}}).GetPrinterAttributes(ctx); err != nil {
		t.Fatalf("insecure: %v", err)
	}
}
