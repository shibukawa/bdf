package server

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/ipp/ipptest"
	"github.com/shibukawa/bdf/print/spool"
)

const proxyToken = "branch-0123456789abcdef"

// branchFixture is a head office whose queue is fetched by the proxy of a
// branch, whose printer is a fake one.
type branchFixture struct {
	t       *testing.T
	hq      *Server
	hqIPP   string // the URL of the queue of the branch at the head office
	branch  *Server
	printer *ipptest.Printer
	stop    func()
}

func serveConfig(t *testing.T, yaml string, ippLn net.Listener) (*Server, func()) {
	t.Helper()
	cfg, err := ParseConfig([]byte(yaml))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	cfg.dir, cfg.TempDir = dir, dir
	cfg.Listen = "127.0.0.1:0"
	if err := cfg.resolve(); err != nil {
		t.Fatal(err)
	}
	srv, err := New(context.Background(), cfg, Options{Grace: time.Second, IPPListener: ippLn, Log: testLog(t)})
	if err != nil {
		t.Fatal(err)
	}
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		srv.Serve(ctx, ln)
		close(done)
	}()
	return srv, func() {
		cancel()
		<-done
	}
}

func newBranchFixture(t *testing.T, branchPaused bool) *branchFixture {
	t.Helper()
	t.Setenv("BDF_TEST_PROXY_SECRET", proxyToken)
	hqLn, _ := net.Listen("tcp", "127.0.0.1:0")
	hq, stopHQ := serveConfig(t, `
printers: { branchp: { proxy: true } }
queues: { branch: { printers: [branchp], defaults: { sides: two-sided-long-edge } } }
ipp: { dns-sd: false }
tokens:
  - { name: branch-a, value: `+proxyToken+`, roles: [proxy], queues: [branch] }
`, hqLn)
	hqIPP := "ipp://" + hqLn.Addr().String() + "/ipp/print/branch"
	p := ipptest.New()
	ps := httptest.NewServer(p)
	t.Cleanup(ps.Close)
	u, _ := url.Parse(ps.URL)
	branch, stopBranch := serveConfig(t, fmt.Sprintf(`
printers: { mfp: { url: "ipp://%s/ipp/print" } }
queues: { local: { printers: [mfp] } }
upstreams:
  hq: { url: "%s", queue: local, username: branch-a, password-env: BDF_TEST_PROXY_SECRET, poll: 100ms }
`, u.Host, hqIPP), nil)
	if branchPaused {
		branch.Spool().PauseQueue(context.Background(), "local")
	}
	f := &branchFixture{t: t, hq: hq, hqIPP: hqIPP, branch: branch, printer: p}
	f.stop = func() {
		stopBranch()
		stopHQ()
	}
	t.Cleanup(func() { f.stop() })
	return f
}

func until(t *testing.T, what string, f func() bool) {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for !f() {
		if time.Now().After(deadline) {
			t.Fatalf("%s: not in time", what)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func (f *branchFixture) submitHQ(name string) spool.Job {
	f.t.Helper()
	j, err := f.hq.Spool().Submit(context.Background(), spool.Submission{Queue: "branch", Name: name, User: "alice",
		Document: bytes.NewReader(docBytes(f.t)), Ticket: print.Ticket{Copies: 2}})
	if err != nil {
		f.t.Fatal(err)
	}
	return j
}

// TestProxyPrints sends a job of the head office to the printer of the
// branch: the branch fetches it, prints it with the head office's ticket
// and queue's defaults, and reports it done; the head office sees the
// branch's printer.
func TestProxyPrints(t *testing.T) {
	f := newBranchFixture(t, false)
	ctx := context.Background()
	hj := f.submitHQ("invoice.bdf")
	var local spool.Job
	until(t, "the branch takes the job", func() bool {
		jobs := f.branch.Spool().List(spool.Filter{Done: true})
		if len(jobs) == 1 {
			local = jobs[0]
		}
		return len(jobs) == 1
	})
	if local.Origin == nil || local.Origin.UUID != hj.UUID || local.Origin.JobID != hj.ID || local.User != "alice" || local.Name != "invoice.bdf" ||
		local.Ticket.Copies != 2 || local.Ticket.Sides != "two-sided-long-edge" || local.Document.Type != bdfFormat {
		t.Fatalf("the branch's job %+v", local)
	}
	until(t, "the head office's job is completed", func() bool {
		j, _ := f.hq.Spool().Get(ctx, hj.ID)
		return j.State == ipp.JobCompleted
	})
	if jobs := f.printer.Jobs(); len(jobs) != 1 {
		t.Fatalf("the printer has %d jobs", len(jobs))
	}
	j, _ := f.hq.Spool().Get(ctx, hj.ID)
	if j.Printer == nil || j.Printer.Device == "" || !j.Printer.Sent {
		t.Fatalf("the head office's job %+v", j.Printer)
	}
	var branchP spool.PrinterInfo
	until(t, "the head office sees the branch's printer", func() bool {
		for _, p := range f.hq.Spool().Printers() {
			branchP = p
		}
		return branchP.Up && branchP.MakeModel == "BDF Test Printer"
	})
}

// TestProxyCancel cancels at the head office a job the branch has: the
// branch cancels it before it prints it.
func TestProxyCancel(t *testing.T) {
	f := newBranchFixture(t, true)
	ctx := context.Background()
	hj := f.submitHQ("a.bdf")
	var local spool.Job
	until(t, "the branch takes the job", func() bool {
		jobs := f.branch.Spool().List(spool.Filter{})
		if len(jobs) == 1 {
			local = jobs[0]
		}
		return len(jobs) == 1
	})
	if _, err := f.hq.Spool().Cancel(ctx, hj.ID, ""); err != nil {
		t.Fatal(err)
	}
	until(t, "the branch cancels it", func() bool {
		j, _ := f.branch.Spool().Get(ctx, local.ID)
		return j.State == ipp.JobCanceled
	})
	if n := len(f.printer.Jobs()); n != 0 {
		t.Fatalf("printed %d", n)
	}
}

// TestProxyTwice: a branch that stopped after it took a job and before it
// told the head office it has the document does not take it again; a job
// the head office gave a branch that no longer has it is aborted, not
// given again. The branch's cycles are run by hand.
func TestProxyTwice(t *testing.T) {
	t.Setenv("BDF_TEST_PROXY_SECRET", proxyToken)
	ctx := context.Background()
	hqLn, _ := net.Listen("tcp", "127.0.0.1:0")
	hq, _ := serveConfig(t, `
printers: { branchp: { proxy: true } }
queues: { branch: { printers: [branchp] } }
ipp: { dns-sd: false }
tokens:
  - { name: branch-a, value: `+proxyToken+`, roles: [proxy], queues: [branch] }
`, hqLn)
	cfg, err := ParseConfig([]byte(fmt.Sprintf(`
printers: { mfp: { url: "ipp://192.0.2.1/ipp/print" } }
queues: { local: { printers: [mfp] } }
upstreams:
  hq: { url: "ipp://%s/ipp/print/branch", queue: local, username: branch-a, password-env: BDF_TEST_PROXY_SECRET }
`, hqLn.Addr())))
	if err != nil {
		t.Fatal(err)
	}
	cfg.dir, cfg.TempDir = t.TempDir(), t.TempDir()
	cfg.resolve()
	branch, err := New(ctx, cfg, Options{Grace: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	defer branch.Close()
	u := branch.ups[0]
	submit := func() spool.Job {
		j, err := hq.Spool().Submit(ctx, spool.Submission{Queue: "branch", Name: "x", User: "alice", Document: bytes.NewReader(docBytes(t))})
		if err != nil {
			t.Fatal(err)
		}
		until(t, "offered", func() bool {
			j, _ := hq.Spool().Get(ctx, j.ID)
			return j.Step == spool.StepFetch
		})
		return j
	}
	// taken, and not acknowledged: the branch has it
	a := submit()
	if _, err := hq.Spool().AcknowledgeJob(ctx, a.ID, "branchp", u.device); err != nil {
		t.Fatal(err)
	}
	if _, err := branch.Spool().Submit(ctx, spool.Submission{Queue: "local", Name: "x", Document: bytes.NewReader(docBytes(t)),
		Origin: &spool.Origin{Upstream: "hq", JobID: a.ID, UUID: a.UUID}}); err != nil {
		t.Fatal(err)
	}
	// given, and lost by the branch
	b := submit()
	hq.Spool().AcknowledgeJob(ctx, b.ID, "branchp", u.device)
	if _, err := hq.Spool().AcknowledgeDocument(ctx, b.ID, "branchp", u.device); err != nil {
		t.Fatal(err)
	}
	// not taken yet
	c := submit()
	if err := u.cycle(ctx); err != nil {
		t.Fatal(err)
	}
	if ja, _ := hq.Spool().Get(ctx, a.ID); !ja.Printer.Sent || ja.Step != spool.StepRemote {
		t.Errorf("a: %+v %+v", ja.Step, ja.Printer)
	}
	if jb, _ := hq.Spool().Get(ctx, b.ID); jb.State != ipp.JobAborted {
		t.Errorf("b: %v", jb.State)
	}
	if jc, _ := hq.Spool().Get(ctx, c.ID); !jc.Printer.Sent {
		t.Errorf("c: %+v", jc.Printer)
	}
	jobs := branch.Spool().List(spool.Filter{})
	if len(jobs) != 2 || jobs[0].Origin.JobID != a.ID || jobs[1].Origin.JobID != c.ID {
		t.Fatalf("the branch has %d jobs", len(jobs))
	}
	// another cycle takes nothing more
	if err := u.cycle(ctx); err != nil {
		t.Fatal(err)
	}
	if n := len(branch.Spool().List(spool.Filter{})); n != 2 {
		t.Fatalf("%d jobs after another cycle", n)
	}
	// another proxy may not take the jobs of this one, nor fetch without a
	// token of a proxy
	other := &ipp.Client{URL: "ipp://" + hqLn.Addr().String() + "/ipp/print/branch", User: "x"}
	m := other.NewRequest(ipp.OpGetJobs)
	m.Add(ipp.TagOperation, ipp.Attr("which-jobs", ipp.Keyword("fetchable")), ipp.Attr("output-device-uuid", ipp.URI(u.device)))
	if _, err := other.Do(ctx, m, nil); err == nil {
		t.Error("fetchable jobs without a token")
	}
	other.Username, other.Password = "branch-a", proxyToken
	d := submit()
	hq.Spool().AcknowledgeJob(ctx, d.ID, "branchp", u.device)
	m = other.NewRequest(ipp.OpAcknowledgeJob)
	m.Add(ipp.TagOperation, ipp.Attr("output-device-uuid", ipp.URI("urn:uuid:00000000-0000-4000-8000-000000000000")), ipp.Attr("job-id", ipp.Integer(int(d.ID))))
	if _, err := other.Do(ctx, m, nil); !ipp.IsStatus(err, ipp.StatusClientErrorNotFetchable) {
		t.Errorf("another proxy's job: %v", err)
	}
}

func testLog(t *testing.T) *slog.Logger {
	if !testing.Verbose() {
		return nil
	}
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

// TestProxyRetry retries at the head office a job the branch took and
// canceled: the branch takes it again.
func TestProxyRetry(t *testing.T) {
	f := newBranchFixture(t, true)
	ctx := context.Background()
	hj := f.submitHQ("a.bdf")
	until(t, "the branch takes the job", func() bool { return len(f.branch.Spool().List(spool.Filter{})) == 1 })
	f.hq.Spool().Cancel(ctx, hj.ID, "")
	until(t, "the branch cancels it", func() bool { return len(f.branch.Spool().List(spool.Filter{})) == 0 })
	if _, err := f.hq.Spool().Retry(ctx, hj.ID); err != nil {
		t.Fatal(err)
	}
	until(t, "the branch takes it again", func() bool { return len(f.branch.Spool().List(spool.Filter{})) == 1 })
	f.branch.Spool().ResumeQueue(ctx, "local")
	until(t, "printed once", func() bool {
		j, _ := f.hq.Spool().Get(ctx, hj.ID)
		return j.State == ipp.JobCompleted
	})
	if n := len(f.printer.Jobs()); n != 1 {
		t.Fatalf("printed %d times", n)
	}
	// a held job bound to the branch cannot be held
	hj2 := f.submitHQ("b.bdf")
	until(t, "the branch takes the second", func() bool {
		j, _ := f.hq.Spool().Get(ctx, hj2.ID)
		return j.Step == spool.StepRemote
	})
	if _, err := f.hq.Spool().Hold(ctx, hj2.ID); err == nil {
		t.Error("held a job at the branch")
	}
}

func TestProxyConfig(t *testing.T) {
	for name, y := range map[string]string{
		"two queues": "printers: { b: { proxy: true } }\nqueues: { q1: { printers: [b] }, q2: { printers: [b] } }\nipp: {}",
		"no ipp":     "printers: { b: { proxy: true } }\nqueues: { q1: { printers: [b] } }",
		"with url":   "printers: { b: { proxy: true, url: ipp://h/ } }\nqueues: { q1: { printers: [b] } }\nipp: {}",
		"class":      "printers: { b: { proxy: true }, p: { url: ipp://h/ } }\nqueues: { q1: { printers: [b, p] } }\nipp: {}",
		"upstream":   "printers: { p: { url: ipp://h/ } }\nqueues: { q: { printers: [p] } }\nupstreams: { hq: { url: ipp://hq/ipp/print/x, queue: nowhere } }",
	} {
		if _, err := ParseConfig([]byte(y)); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
}
