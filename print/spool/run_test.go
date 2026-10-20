package spool

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/fakeprinter"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/raw"
)

// fast shortens the spool's waits for a test.
func fast(t *testing.T) {
	t.Helper()
	li, lb, wf, wm, rd := lookIdle, lookBusy, watchFirst, watchMax, retryDelays
	lookIdle, lookBusy, watchFirst, watchMax = 100*time.Millisecond, 30*time.Millisecond, 10*time.Millisecond, 50*time.Millisecond
	retryDelays = []time.Duration{50 * time.Millisecond, 100 * time.Millisecond}
	t.Cleanup(func() { lookIdle, lookBusy, watchFirst, watchMax, retryDelays = li, lb, wf, wm, rd })
}

// running opens a spool on the real clock and starts it.
func (f *fixture) running() *Spool {
	f.t.Helper()
	f.o.Now = nil
	s := f.open()
	s.Start()
	return s
}

func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func jobState(s *Spool, id int32) Job {
	j, _ := s.Get(context.Background(), id)
	return j
}

func documents(p *fakeprinter.Printer) int {
	n := 0
	for _, j := range p.Jobs() {
		n += len(j.Documents)
	}
	return n
}

func TestRunPrints(t *testing.T) {
	fast(t)
	f := newFixture(t)
	s := f.running()
	doc := bdfBytes(t, testDoc(2))
	var ids []int32
	for i := 0; i < 3; i++ {
		ids = append(ids, submit(t, s, Submission{Name: fmt.Sprintf("d%d.bdf", i), Document: bytes.NewReader(doc), Ticket: print.Ticket{Sides: "two-sided-long-edge"}}).ID)
	}
	for _, id := range ids {
		waitFor(t, fmt.Sprintf("job %d", id), func() bool { return jobState(s, id).State == ipp.JobCompleted })
	}
	p := f.printers["mfp1"]
	if len(p.Jobs()) != 3 || documents(p) != 3 {
		t.Fatalf("the printer has %d jobs and %d documents", len(p.Jobs()), documents(p))
	}
	j := jobState(s, ids[0])
	if j.Reasons[0] != "job-completed-successfully" || j.Attempts != 1 || j.Printer.Name != "mfp1" || !j.Printer.Sent || j.Format != "application/pdf" || j.Sheets != 2 {
		t.Fatalf("job done: %+v", j)
	}
	if got := p.Jobs()[0].Name; got != "#1.1 d0.bdf" {
		t.Fatalf("job-name at the printer: %q", got)
	}
	if len(s.List(Filter{})) != 0 {
		t.Fatal("jobs left not done")
	}
}

func TestRunInflightAndPriority(t *testing.T) {
	fast(t)
	f := newFixture(t)
	p := f.printers["mfp1"]
	p.Processing = 300 * time.Millisecond
	s := f.running()
	doc := bdfBytes(t, testDoc(1))
	s.PauseQueue(context.Background(), "office")
	var ids []int32
	for i := 0; i < 4; i++ {
		ids = append(ids, submit(t, s, Submission{Document: bytes.NewReader(doc), Priority: 10 + 10*i}).ID)
	}
	s.ResumeQueue(context.Background(), "office")
	most := 0
	waitFor(t, "all the jobs", func() bool {
		inside := 0
		for _, j := range p.Jobs() {
			if st := p.JobState(j); !st.Done() {
				inside++
			}
		}
		most = max(most, inside)
		for _, id := range ids {
			if jobState(s, id).State != ipp.JobCompleted {
				return false
			}
		}
		return true
	})
	if most > 2 {
		t.Fatalf("the printer had %d jobs at once", most)
	}
	// the highest priority first
	if got := p.Jobs()[0].Name; !strings.HasPrefix(got, fmt.Sprintf("#%d.", ids[3])) {
		t.Fatalf("first sent: %q", got)
	}
}

func TestRunFailoverAndAny(t *testing.T) {
	fast(t)
	f := newFixture(t, "mfp1", "mfp2")
	f.printers["mfp1"].Stop("media-empty-error")
	s := f.running()
	doc := bdfBytes(t, testDoc(1))
	a := submit(t, s, Submission{Document: bytes.NewReader(doc)})
	waitFor(t, "failover", func() bool { return jobState(s, a.ID).State == ipp.JobCompleted })
	if jobState(s, a.ID).Printer.Name != "mfp2" || len(f.printers["mfp1"].Jobs()) != 0 {
		t.Fatalf("failover: %+v", jobState(s, a.ID).Printer)
	}
	s.Close()
	// any: both printers take jobs
	f2 := newFixture(t, "mfp1", "mfp2")
	f2.o.Queues[0].Policy = "any"
	for _, p := range f2.printers {
		p.Processing = 200 * time.Millisecond
	}
	s2 := f2.running()
	var ids []int32
	for i := 0; i < 6; i++ {
		ids = append(ids, submit(t, s2, Submission{Document: bytes.NewReader(doc)}).ID)
	}
	for _, id := range ids {
		waitFor(t, "any", func() bool { return jobState(s2, id).State == ipp.JobCompleted })
	}
	if len(f2.printers["mfp1"].Jobs()) == 0 || len(f2.printers["mfp2"].Jobs()) == 0 {
		t.Fatalf("any: %d and %d jobs", len(f2.printers["mfp1"].Jobs()), len(f2.printers["mfp2"].Jobs()))
	}
}

func TestRunStoppedPrinterWaits(t *testing.T) {
	fast(t)
	f := newFixture(t)
	p := f.printers["mfp1"]
	p.Stop("media-empty-error")
	s := f.running()
	a := submit(t, s, Submission{Document: bytes.NewReader(bdfBytes(t, testDoc(1)))})
	time.Sleep(200 * time.Millisecond)
	if len(p.Jobs()) != 0 || jobState(s, a.ID).State != ipp.JobPending {
		t.Fatalf("a job was sent to a stopped printer: %+v", jobState(s, a.ID))
	}
	if info := s.Printers(); info[0].State != "stopped" || info[0].Reasons[0] != "media-empty-error" {
		t.Fatalf("printer: %+v", info)
	}
	p.Resume()
	waitFor(t, "the job after the printer is back", func() bool { return jobState(s, a.ID).State == ipp.JobCompleted })
}

func TestRunRetry(t *testing.T) {
	fast(t)
	f := newFixture(t)
	p := f.printers["mfp1"]
	p.Busy = 2 // refuses the job twice: it is sent again
	s := f.running()
	a := submit(t, s, Submission{Document: bytes.NewReader(bdfBytes(t, testDoc(1)))})
	waitFor(t, "the job retried", func() bool { return jobState(s, a.ID).State == ipp.JobCompleted })
	if j := jobState(s, a.ID); j.Attempts != 3 || j.Retries != 0 || len(p.Jobs()) != 1 || p.Requests["Create-Job"] != 3 {
		t.Fatalf("retried: %d attempts, %d retries, %d jobs, %v", j.Attempts, j.Retries, len(p.Jobs()), p.Requests)
	}
	// a printer not accepting jobs is not sent to (once the watcher saw
	// it: an attempt before is refused and retried)
	p.Accepting = false
	b := submit(t, s, Submission{Document: bytes.NewReader(bdfBytes(t, testDoc(1)))})
	time.Sleep(300 * time.Millisecond)
	if j := jobState(s, b.ID); j.State != ipp.JobPending || j.Attempts > 1 {
		t.Fatalf("sent to a printer not accepting jobs: %+v", j)
	}
	p.Accepting = true
	waitFor(t, "the job once the printer accepts", func() bool { return jobState(s, b.ID).State == ipp.JobCompleted })
	s.Close()
	// retried for too long: aborted
	f2 := newFixture(t)
	f2.printers["mfp1"].Busy = 1000
	f2.o.Limits.RetryFor = 50 * time.Millisecond
	s2 := f2.running()
	c := submit(t, s2, Submission{Document: bytes.NewReader(bdfBytes(t, testDoc(1)))})
	waitFor(t, "the retries to end", func() bool { return jobState(s2, c.ID).State == ipp.JobAborted })
	if j := jobState(s2, c.ID); j.Reasons[0] != ReasonRetryExhausted {
		t.Fatalf("aborted for %v", j.Reasons)
	}
}

func TestRunAbortedNotResent(t *testing.T) {
	fast(t)
	f := newFixture(t)
	p := f.printers["mfp1"]
	p.Processing = time.Hour
	s := f.running()
	a := submit(t, s, Submission{Document: bytes.NewReader(bdfBytes(t, testDoc(1)))})
	waitFor(t, "the job at the printer", func() bool { return len(p.Jobs()) == 1 && jobState(s, a.ID).Step == StepWatch })
	p.Abort(p.Jobs()[0].ID)
	waitFor(t, "the job aborted", func() bool { return jobState(s, a.ID).State == ipp.JobAborted })
	time.Sleep(150 * time.Millisecond)
	if len(p.Jobs()) != 1 {
		t.Fatalf("an aborted job was sent again: %d jobs", len(p.Jobs()))
	}
	// Retry sends it again, as a new attempt
	if _, err := s.Retry(context.Background(), a.ID); err != nil {
		t.Fatal(err)
	}
	p.Processing = 0
	waitFor(t, "the job retried", func() bool { return jobState(s, a.ID).State == ipp.JobCompleted })
	if j := jobState(s, a.ID); j.Attempts != 2 || p.Jobs()[1].Name != fmt.Sprintf("#%d.2 document", a.ID) {
		t.Fatalf("retried: %d attempts, %q", j.Attempts, p.Jobs()[1].Name)
	}
}

func TestRunCancelAtPrinter(t *testing.T) {
	fast(t)
	f := newFixture(t)
	p := f.printers["mfp1"]
	p.Processing = time.Hour
	s := f.running()
	a := submit(t, s, Submission{Document: bytes.NewReader(bdfBytes(t, testDoc(1)))})
	waitFor(t, "the job at the printer", func() bool { return jobState(s, a.ID).Step == StepWatch })
	if _, err := s.Cancel(context.Background(), a.ID, ""); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "the job canceled", func() bool { return jobState(s, a.ID).State == ipp.JobCanceled })
	if st := p.JobState(p.Jobs()[0]); st != ipp.JobCanceled {
		t.Fatalf("the job at the printer is %v", st)
	}
}

func TestRunConvert(t *testing.T) {
	fast(t)
	f := newFixture(t)
	s := f.running()
	ctx := context.Background()
	a := submit(t, s, Submission{Name: "a.pdf", Document: strings.NewReader("%PDF-1.7")})
	waitFor(t, "a converted job", func() bool { return jobState(s, a.ID).State == ipp.JobCompleted })
	j := jobState(s, a.ID)
	if j.Document.BDF == nil || j.Document.BDF.Size == 0 || !slices.Contains(j.Warnings, "converted") {
		t.Fatalf("converted: %+v", j)
	}
	b := submit(t, s, Submission{Name: "b.pdf", Document: strings.NewReader("%PDF-encrypted")})
	waitFor(t, "a job waiting for its password", func() bool {
		j := jobState(s, b.ID)
		return j.State == ipp.JobPendingHeld && slices.Contains(j.Reasons, "job-password-wait")
	})
	if _, err := s.SetPassword(ctx, b.ID, "secret"); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "the job with its password", func() bool { return jobState(s, b.ID).State == ipp.JobCompleted })
	c := submit(t, s, Submission{Name: "c.pdf", Document: strings.NewReader("%PDF-broken")})
	waitFor(t, "a broken document", func() bool { return jobState(s, c.ID).State == ipp.JobAborted })
	if j := jobState(s, c.ID); j.Reasons[0] != "document-format-error" {
		t.Fatalf("broken: %+v", j)
	}
	if f.conv.converts != 4 {
		t.Fatalf("%d conversions", f.conv.converts)
	}
}

func TestRunEncryptedAfterRestart(t *testing.T) {
	fast(t)
	f := newFixture(t)
	f.printers["mfp1"].Stop()
	s := f.running()
	a := submit(t, s, Submission{Document: bytes.NewReader(encrypted(t, "secret")), Password: "secret"})
	s.Close()
	f.printers["mfp1"].Resume()
	s = f.running()
	waitFor(t, "the job without its password", func() bool {
		j := jobState(s, a.ID)
		return j.State == ipp.JobPendingHeld && slices.Contains(j.Reasons, "job-password-wait")
	})
	if _, err := s.SetPassword(context.Background(), a.ID, "secret"); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "the encrypted job", func() bool { return jobState(s, a.ID).State == ipp.JobCompleted })
}

func TestRunHoldWindow(t *testing.T) {
	fast(t)
	f := newFixture(t)
	s := f.running()
	a := submit(t, s, Submission{Document: bytes.NewReader(bdfBytes(t, testDoc(1))), HoldUntilTime: time.Now().Add(200 * time.Millisecond)})
	if a.State != ipp.JobPendingHeld {
		t.Fatalf("held: %+v", a)
	}
	waitFor(t, "the hold to end", func() bool { return jobState(s, a.ID).State == ipp.JobCompleted })
}

// The crash tests run a spool in a child process that exits at a point of
// a send, then open the spool again here and count what the printer
// printed.
func TestCrashChild(t *testing.T) {
	point := os.Getenv("SPOOL_CRASH_POINT")
	if point == "" {
		t.Skip("the child of TestCrash")
	}
	crashHook = func(p string) {
		if p == point {
			os.Exit(3)
		}
	}
	watchFirst, lookBusy = 10*time.Millisecond, 30*time.Millisecond
	s, err := Open(context.Background(), Options{Path: os.Getenv("SPOOL_CRASH_PATH"), Queues: []Queue{{Name: "office", Printers: []string{"mfp1"}}},
		Printers: []Printer{crashPrinter(t, os.Getenv("SPOOL_CRASH_PRINTER"))}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Submit(context.Background(), Submission{Queue: "office", Name: "doc.bdf", Document: bytes.NewReader(bdfBytes(t, testDoc(1)))}); err != nil {
		t.Fatal(err)
	}
	s.Start()
	time.Sleep(10 * time.Second)
	os.Exit(4) // the point was not reached
}

// crashPrinter is the printer of a URL: of IPP, or raw.
func crashPrinter(t *testing.T, u string) Printer {
	if raw.IsURL(u) {
		p, err := raw.Parse(u)
		if err != nil {
			t.Fatal(err)
		}
		return Printer{Name: "mfp1", Raw: p}
	}
	return Printer{Name: "mfp1", Client: &ipp.Client{URL: u, User: "bdf"}}
}

func TestCrash(t *testing.T) {
	fast(t)
	for _, c := range []struct {
		point string
		// before changes the printer before the spool opens again
		before func(p *fakeprinter.Printer, c *ipp.Client)
		// state of the job, and jobs and documents of the printer, after
		state            ipp.JobState
		reason           string
		jobs, docs, canc int
	}{
		{point: "attempt-recorded", state: ipp.JobCompleted, reason: "job-completed-successfully", jobs: 1, docs: 1},
		{point: "created", state: ipp.JobCompleted, reason: "job-completed-successfully", jobs: 2, docs: 1, canc: 1},
		{point: "created-recorded", state: ipp.JobCompleted, reason: "job-completed-successfully", jobs: 2, docs: 1, canc: 1},
		{point: "sent-recorded", state: ipp.JobCompleted, reason: "job-completed-successfully", jobs: 1, docs: 1},
		{point: "watching", state: ipp.JobCompleted, reason: "job-completed-successfully", jobs: 1, docs: 1},
		{point: "sent-recorded", before: func(p *fakeprinter.Printer, _ *ipp.Client) { p.Forget() },
			state: ipp.JobCompleted, reason: "job-completed-with-warnings", jobs: 0, docs: 0},
		{point: "attempt-recorded", before: func(p *fakeprinter.Printer, c *ipp.Client) {
			// two jobs that may be this one
			for i := 0; i < 2; i++ {
				c.PrintJob(context.Background(), "#1.1 doc.bdf", "application/pdf", nil, strings.NewReader("%PDF"))
			}
		}, state: ipp.JobPendingHeld, reason: ReasonDuplicateCheck, jobs: 2, docs: 2},
		// the printer was switched off and on: its job ids count from 1
		// again, and another's job has the id the spool kept
		{point: "created-recorded", before: func(p *fakeprinter.Printer, c *ipp.Client) {
			p.Inject(fakeprinter.Fault{Restart: true})
			p.Stop("media-empty-error")
			c.PrintJob(context.Background(), "someone else's", "application/pdf", nil, strings.NewReader("%PDF"))
			p.Resume()
		}, state: ipp.JobCompleted, reason: "job-completed-successfully", jobs: 2, docs: 2},
		{point: "sent-recorded", before: func(p *fakeprinter.Printer, c *ipp.Client) {
			p.Inject(fakeprinter.Fault{Restart: true})
			p.Inject(fakeprinter.Fault{NeverComplete: ptr(true)})
			c.PrintJob(context.Background(), "someone else's", "application/pdf", nil, strings.NewReader("%PDF"))
		}, state: ipp.JobCompleted, reason: "job-completed-with-warnings", jobs: 1, docs: 1},
	} {
		name := c.point
		if c.before != nil {
			name += "+"
		}
		t.Run(name, func(t *testing.T) {
			p := fakeprinter.New()
			printerURL := fakeprinter.Start(t, p)
			path := filepath.Join(t.TempDir(), "spool.db")
			cmd := exec.Command(os.Args[0], "-test.run=^TestCrashChild$")
			cmd.Env = append(os.Environ(), "SPOOL_CRASH_POINT="+c.point, "SPOOL_CRASH_PATH="+path, "SPOOL_CRASH_PRINTER="+printerURL)
			out, err := cmd.CombinedOutput()
			if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 3 {
				t.Fatalf("the child did not crash at %s: %v\n%s", c.point, err, out)
			}
			client := &ipp.Client{URL: printerURL, User: "bdf"}
			if c.before != nil {
				c.before(p, client)
			}
			s, err := Open(context.Background(), Options{Path: path, Queues: []Queue{{Name: "office", Printers: []string{"mfp1"}}},
				Printers: []Printer{{Name: "mfp1", Client: client}}})
			if err != nil {
				t.Fatal(err)
			}
			defer s.Close()
			s.Start()
			waitFor(t, "the job after the crash", func() bool { return jobState(s, 1).State == c.state })
			time.Sleep(100 * time.Millisecond)
			j := jobState(s, 1)
			canceled := 0
			for _, pj := range p.Jobs() {
				if p.JobState(pj) == ipp.JobCanceled {
					canceled++
				}
			}
			if j.State != c.state || j.Reasons[0] != c.reason || len(p.Jobs()) != c.jobs || documents(p) != c.docs || canceled != c.canc {
				t.Fatalf("after %s: job %s %v; printer: %d jobs, %d documents, %d canceled\n%+v", c.point, j.State, j.Reasons, len(p.Jobs()), documents(p), canceled, j)
			}
		})
	}
}

func TestNextLookOfPrinterNotLookedAt(t *testing.T) {
	f := newFixture(t)
	s := f.open()
	// a send that did not reach the printer marks it down before a look
	// failed
	p := s.printer("mfp1")
	s.mu.Lock()
	p.up, p.failing = false, 0
	s.mu.Unlock()
	if d := s.nextLook(p); d != lookBusy {
		t.Fatalf("next look in %v", d)
	}
}

func ptr[T any](v T) *T { return &v }

// TestPrinterRestartedWhileWatched restarts the printer while the spool
// follows a job there: the printer counts its job ids from 1 again and
// another's job takes the id. The spool does not follow the other's job;
// its own is gone with the printer's memory, and taken as printed.
func TestPrinterRestartedWhileWatched(t *testing.T) {
	fast(t)
	f := newFixture(t)
	p, c := f.printers["mfp1"], f.clients["mfp1"]
	p.Processing = time.Hour
	s := f.running()
	a := submit(t, s, Submission{Document: bytes.NewReader(bdfBytes(t, testDoc(1)))})
	waitFor(t, "the job sent to the printer", func() bool {
		j := jobState(s, a.ID)
		return j.Printer != nil && j.Printer.Sent && j.Step == StepWatch
	})
	p.Inject(fakeprinter.Fault{Restart: true})
	if _, err := c.PrintJob(context.Background(), "someone else's", "application/pdf", nil, strings.NewReader("%PDF")); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "the job done", func() bool { return jobState(s, a.ID).State.Done() })
	j := jobState(s, a.ID)
	if j.State != ipp.JobCompleted || j.Reasons[0] != "job-completed-with-warnings" {
		t.Fatalf("after the restart: %s %v %q %v", j.State, j.Reasons, j.Message, j.Warnings)
	}
	if pj := p.Jobs(); len(pj) != 1 || pj[0].Name != "someone else's" || p.JobState(pj[0]) == ipp.JobCanceled {
		t.Fatalf("the other's job: %+v", pj)
	}
}
