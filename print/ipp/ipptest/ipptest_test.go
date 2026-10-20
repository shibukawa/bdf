package ipptest_test

import (
	"bytes"
	"context"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/ipp/ipptest"
)

func client(t *testing.T, p *ipptest.Printer) *ipp.Client {
	t.Helper()
	srv := httptest.NewServer(p)
	t.Cleanup(srv.Close)
	u, _ := url.Parse(srv.URL)
	return &ipp.Client{URL: "ipp://" + u.Host + "/ipp/print", User: "spool"}
}

func TestStopForget(t *testing.T) {
	p := ipptest.New()
	c := client(t, p)
	ctx := context.Background()
	p.Stop("media-empty-error")
	pr, err := c.GetPrinterAttributes(ctx)
	if err != nil || pr.State != ipp.PrinterStopped || strings.Join(pr.Reasons, ",") != "media-empty-error" {
		t.Fatalf("stopped: %v %v %v", pr.State, pr.Reasons, err)
	}
	js, err := c.PrintJob(ctx, "#1.1 a.pdf", "application/pdf", nil, bytes.NewReader([]byte("%PDF")))
	if err != nil {
		t.Fatal(err)
	}
	if st, _ := c.GetJobAttributes(ctx, js.ID); st.State != ipp.JobProcessing {
		t.Fatalf("a job of a stopped printer is %v", st.State)
	}
	p.Resume()
	if st, _ := c.GetJobAttributes(ctx, js.ID); st.State != ipp.JobCompleted {
		t.Fatalf("after Resume the job is %v", st.State)
	}
	all, err := c.GetJobs(ctx, "all", true, 0)
	if err != nil || len(all) != 1 || all[0].Name != "#1.1 a.pdf" {
		t.Fatalf("which-jobs=all: %v %v", all, err)
	}
	p.Forget()
	if _, err := c.GetJobAttributes(ctx, js.ID); !ipp.IsStatus(err, ipp.StatusClientErrorNotFound) {
		t.Fatalf("forgotten job: %v", err)
	}
}

func TestFailures(t *testing.T) {
	p := ipptest.New()
	c := client(t, p)
	ctx := context.Background()
	p.DropCreateResponse()
	if _, err := c.CreateJob(ctx, "#2.1 b.pdf", nil); err == nil {
		t.Fatal("Create-Job answered although dropped")
	}
	if jobs := p.Jobs(); len(jobs) != 1 || jobs[0].Name != "#2.1 b.pdf" {
		t.Fatalf("the dropped Create-Job made %v", jobs)
	}
	js, err := c.CreateJob(ctx, "#3.1 c.pdf", nil)
	if err != nil {
		t.Fatal(err)
	}
	p.FailAfter(10)
	if _, err := c.SendDocument(ctx, js.ID, "c.pdf", "application/pdf", bytes.NewReader(bytes.Repeat([]byte("x"), 1<<16)), true); err == nil {
		t.Fatal("the document was taken although cut")
	}
	j := p.Job(js.ID)
	if len(j.Documents) != 0 || p.JobState(j) != ipp.JobPending {
		t.Fatalf("a cut document: %d documents, %v", len(j.Documents), p.JobState(j))
	}
	// the next document goes through
	if _, err := c.SendDocument(ctx, js.ID, "c.pdf", "application/pdf", bytes.NewReader([]byte("%PDF")), true); err != nil {
		t.Fatal(err)
	}
}
