package spool

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/ipp/ipptest"
	bolt "go.etcd.io/bbolt"
)

const a4w, a4h = 595.276, 841.89

// testDoc is a document of n A4 pages.
func testDoc(n int) *bdf.Document {
	d := bdf.NewDocument()
	v := d.NewView("pages", bdf.ViewFixed, "")
	for i := 0; i < n; i++ {
		o := bdf.NewObject()
		o.FillColor(0x000000ff).FillRect(50, 100, float32(50*(i+1)), 40)
		h, _ := d.AddObject(o)
		v.AddPage(a4w, a4h, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	}
	return d
}

func bdfBytes(t *testing.T, d *bdf.Document) []byte {
	t.Helper()
	var b bytes.Buffer
	if err := d.WriteSingle(&b); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func encrypted(t *testing.T, password string) []byte {
	t.Helper()
	d := testDoc(1)
	lock, err := bdf.NewPasswordLock(password, 1)
	if err != nil {
		t.Fatal(err)
	}
	d.Lock = lock
	return bdfBytes(t, d)
}

// fakeConverter takes "%PDF" documents: a page each, or the error of the
// password asked.
type fakeConverter struct {
	mu       sync.Mutex
	converts int
	password string // the password of encrypted inputs ("%PDF-encrypted")
}

func (c *fakeConverter) Detect(r io.ReaderAt, size int64, name string) (string, bool) {
	b := make([]byte, 4)
	if n, _ := r.ReadAt(b, 0); n == 4 && string(b) == "%PDF" {
		return "pdf", true
	}
	return "", false
}

func (c *fakeConverter) Convert(ctx context.Context, r io.ReaderAt, size int64, format, name, password string) (*bdf.Document, []string, error) {
	c.mu.Lock()
	c.converts++
	c.mu.Unlock()
	b := make([]byte, size)
	r.ReadAt(b, 0)
	if bytes.HasPrefix(b, []byte("%PDF-encrypted")) {
		switch password {
		case "":
			return nil, nil, ErrPasswordRequired
		case c.password:
		default:
			return nil, nil, ErrWrongPassword
		}
	}
	if bytes.HasPrefix(b, []byte("%PDF-broken")) {
		return nil, nil, errors.New("broken pdf")
	}
	return testDoc(1), []string{"converted"}, nil
}

// clock is a clock a test moves.
type clock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *clock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *clock) Add(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// fixture is a spool over fake printers.
type fixture struct {
	t        *testing.T
	path     string
	clock    *clock
	printers map[string]*ipptest.Printer
	clients  map[string]*ipp.Client
	conv     *fakeConverter
	o        Options
}

func newFixture(t *testing.T, printers ...string) *fixture {
	t.Helper()
	if len(printers) == 0 {
		printers = []string{"mfp1"}
	}
	f := &fixture{t: t, path: filepath.Join(t.TempDir(), "spool.db"), clock: &clock{now: time.Date(2026, 10, 20, 10, 0, 0, 0, time.UTC)},
		printers: map[string]*ipptest.Printer{}, clients: map[string]*ipp.Client{}, conv: &fakeConverter{password: "secret"}}
	var ps []Printer
	for _, name := range printers {
		p := ipptest.New()
		srv := httptest.NewServer(p)
		t.Cleanup(srv.Close)
		u, _ := url.Parse(srv.URL)
		c := &ipp.Client{URL: "ipp://" + u.Host + "/ipp/print", User: "bdf"}
		f.printers[name], f.clients[name] = p, c
		ps = append(ps, Printer{Name: name, Client: c})
	}
	f.o = Options{Path: f.path, Queues: []Queue{{Name: "office", Printers: printers}, {Name: "spare", Printers: printers[:1]}},
		Printers: ps, Converter: f.conv, Now: f.clock.Now, TempDir: t.TempDir()}
	return f
}

func (f *fixture) open() *Spool {
	f.t.Helper()
	s, err := Open(context.Background(), f.o)
	if err != nil {
		f.t.Fatal(err)
	}
	f.t.Cleanup(func() { s.Close() })
	return s
}

func submit(t *testing.T, s *Spool, sub Submission) Job {
	t.Helper()
	if sub.Queue == "" {
		sub.Queue = "office"
	}
	j, err := s.Submit(context.Background(), sub)
	if err != nil {
		t.Fatal(err)
	}
	return j
}

func TestSubmitAndReopen(t *testing.T) {
	f := newFixture(t)
	s := f.open()
	ctx := context.Background()
	doc := bdfBytes(t, testDoc(2))
	a := submit(t, s, Submission{Name: "invoice.bdf", User: "erp", Document: bytes.NewReader(doc), Ticket: print.Ticket{Sides: "two-sided-long-edge"}})
	b := submit(t, s, Submission{Name: "report.pdf", User: "erp", Document: strings.NewReader("%PDF-1.7 …"), Priority: 80})
	if a.ID != 1 || b.ID != 2 || a.State != ipp.JobPending || a.Step != StepSend || b.Step != StepConvert || a.Document.Format != "bdf" || b.Document.Format != "pdf" {
		t.Fatalf("jobs %+v\n%+v", a, b)
	}
	if a.Ticket.JobName != "invoice.bdf" || a.Document.Size != int64(len(doc)) || len(a.Document.SHA256) != 64 || a.UUID == "" {
		t.Fatalf("job a %+v", a)
	}
	r, err := s.st.getDoc(ctx, originalKey(a.ID))
	if err != nil {
		t.Fatal(err)
	}
	got, _ := io.ReadAll(r)
	r.Close()
	if !bytes.Equal(got, doc) {
		t.Fatal("the document kept is not the one received")
	}
	// unknown, empty and too large documents
	if _, err := s.Submit(ctx, Submission{Queue: "office", Document: strings.NewReader("hello")}); !errors.Is(err, ErrUnknownFormat) {
		t.Fatalf("unknown format: %v", err)
	}
	if _, err := s.Submit(ctx, Submission{Queue: "office", Document: strings.NewReader("")}); !errors.Is(err, ErrBadDocument) {
		t.Fatalf("empty: %v", err)
	}
	if _, err := s.Submit(ctx, Submission{Queue: "nowhere", Document: bytes.NewReader(doc)}); !errors.Is(err, ErrNoQueue) {
		t.Fatalf("no queue: %v", err)
	}
	s.Close()
	f.o.Limits.MaxDocument = 10
	s = f.open()
	if _, err := s.Submit(ctx, Submission{Queue: "office", Document: bytes.NewReader(doc)}); !errors.Is(err, ErrTooLarge) {
		t.Fatalf("too large: %v", err)
	}
	// reopened: the same jobs, and ids that go on
	list := s.List(Filter{})
	if len(list) != 2 || list[0].ID != 1 || list[1].ID != 2 || list[1].Priority != 80 {
		t.Fatalf("after reopening: %+v", list)
	}
	f.o.Limits.MaxDocument = 0
	s.Close()
	s = f.open()
	c := submit(t, s, Submission{Document: bytes.NewReader(doc)})
	if c.ID <= 2 {
		t.Fatalf("an id given again: %d", c.ID)
	}
	// a second process cannot open the spool
	f.o.Timeout = 50 * time.Millisecond
	other := exec2(t, f.o)
	if other == nil {
		t.Fatal("a second spool opened the same file")
	}
}

// exec2 opens the spool a second time in this process: the file is
// shared by yuruidb within a process, so the lock is checked by opening it
// with bbolt directly as another process would.
func exec2(t *testing.T, o Options) error {
	t.Helper()
	return openLocked(o.Path, o.Timeout)
}

func TestOperations(t *testing.T) {
	f := newFixture(t)
	s := f.open()
	ctx := context.Background()
	doc := bdfBytes(t, testDoc(1))
	a := submit(t, s, Submission{Document: bytes.NewReader(doc)})
	b := submit(t, s, Submission{Document: bytes.NewReader(doc), HoldUntil: "indefinite"})
	if b.State != ipp.JobPendingHeld || b.Reasons[0] != "job-hold-until-specified" {
		t.Fatalf("held job %+v", b)
	}
	if j, err := s.Hold(ctx, a.ID); err != nil || j.State != ipp.JobPendingHeld {
		t.Fatalf("hold: %+v %v", j, err)
	}
	if j, err := s.Release(ctx, a.ID, ""); err != nil || j.State != ipp.JobPending {
		t.Fatalf("release: %+v %v", j, err)
	}
	if j, err := s.Move(ctx, a.ID, "spare"); err != nil || j.Queue != "spare" {
		t.Fatalf("move: %+v %v", j, err)
	}
	if _, err := s.Move(ctx, a.ID, "nowhere"); !errors.Is(err, ErrNoQueue) {
		t.Fatalf("move nowhere: %v", err)
	}
	if j, err := s.SetPriority(ctx, a.ID, 90); err != nil || j.Priority != 90 {
		t.Fatalf("priority: %+v %v", j, err)
	}
	if j, err := s.Cancel(ctx, a.ID, ""); err != nil || j.State != ipp.JobCanceled || j.Finished.IsZero() {
		t.Fatalf("cancel: %+v %v", j, err)
	}
	if _, err := s.Cancel(ctx, a.ID, ""); !errors.Is(err, ErrNotPossible) {
		t.Fatalf("cancel twice: %v", err)
	}
	if _, err := s.Hold(ctx, 999); !errors.Is(err, ErrNotFound) {
		t.Fatalf("no such job: %v", err)
	}
	if got := s.List(Filter{}); len(got) != 1 || got[0].ID != b.ID {
		t.Fatalf("jobs not done: %+v", got)
	}
	if got := s.List(Filter{Done: true}); len(got) != 2 {
		t.Fatalf("with the jobs done: %+v", got)
	}
	done, err := s.Get(ctx, a.ID)
	if err != nil || done.State != ipp.JobCanceled {
		t.Fatalf("get a job done: %+v %v", done, err)
	}
	r, err := s.Retry(ctx, a.ID)
	if err != nil || r.State != ipp.JobPending || r.Step != StepSend || r.Queue != "spare" {
		t.Fatalf("retry: %+v %v", r, err)
	}
	if _, err := s.Retry(ctx, a.ID); !errors.Is(err, ErrNotPossible) {
		t.Fatalf("retry a job not done: %v", err)
	}
	// queues
	if q, err := s.RejectQueue(ctx, "office"); err != nil || !q.Rejecting {
		t.Fatalf("reject: %+v %v", q, err)
	}
	if _, err := s.Submit(ctx, Submission{Queue: "office", Document: bytes.NewReader(doc)}); !errors.Is(err, ErrNotAccepting) {
		t.Fatalf("a queue rejecting: %v", err)
	}
	if _, err := s.PauseQueue(ctx, "spare"); err != nil {
		t.Fatal(err)
	}
	s.Close()
	s = f.open()
	qs := s.Queues()
	if len(qs) != 2 || !qs[0].Rejecting || !qs[1].Paused || qs[1].Jobs["pending"] != 1 {
		t.Fatalf("queues after reopening: %+v", qs)
	}
}

func TestEncryptedSubmit(t *testing.T) {
	f := newFixture(t)
	s := f.open()
	ctx := context.Background()
	doc := encrypted(t, "secret")
	if _, err := s.Submit(ctx, Submission{Queue: "office", Document: bytes.NewReader(doc), Password: "wrong"}); !errors.Is(err, ErrWrongPassword) {
		t.Fatalf("wrong password: %v", err)
	}
	a := submit(t, s, Submission{Document: bytes.NewReader(doc)})
	if a.State != ipp.JobPendingHeld || a.Reasons[0] != "job-password-wait" || !a.Document.Encrypted {
		t.Fatalf("encrypted without password: %+v", a)
	}
	if _, err := s.Release(ctx, a.ID, ""); !errors.Is(err, ErrNotPossible) {
		t.Fatalf("release without password: %v", err)
	}
	if _, err := s.SetPassword(ctx, a.ID, "wrong"); !errors.Is(err, ErrWrongPassword) {
		t.Fatalf("set a wrong password: %v", err)
	}
	if j, err := s.SetPassword(ctx, a.ID, "secret"); err != nil || j.State != ipp.JobPending {
		t.Fatalf("set password: %+v %v", j, err)
	}
	b := submit(t, s, Submission{Document: bytes.NewReader(doc), Password: "secret"})
	if b.State != ipp.JobPending {
		t.Fatalf("encrypted with password: %+v", b)
	}
	// the password is not kept: after a restart the job waits for it
	if s.password(b.ID) != "secret" {
		t.Fatal("the password is not in memory")
	}
	s.Close()
	s = f.open()
	if s.password(b.ID) != "" {
		t.Fatal("a password survived a restart")
	}
	raw, _ := json.Marshal(s.List(Filter{}))
	if bytes.Contains(raw, []byte("secret")) {
		t.Fatal("a password is in a job")
	}
}

func TestCrashLeftovers(t *testing.T) {
	f := newFixture(t)
	s := f.open()
	ctx := context.Background()
	doc := bdfBytes(t, testDoc(1))
	a := submit(t, s, Submission{Document: bytes.NewReader(doc)})
	b := submit(t, s, Submission{Document: bytes.NewReader(doc)})
	// a crash between recording a job done and removing it from the jobs
	done := a
	done.to(f.clock.Now(), ipp.JobCanceled, "job-canceled-by-user")
	if err := s.st.insertDone(ctx, &done); err != nil {
		t.Fatal(err)
	}
	// a crash while a document arrived
	if _, err := s.st.putDoc(ctx, originalKey(77), strings.NewReader("half")); err != nil {
		t.Fatal(err)
	}
	s.Close()
	s = f.open()
	if got := s.List(Filter{}); len(got) != 1 || got[0].ID != b.ID {
		t.Fatalf("jobs after the crash: %+v", got)
	}
	if j, err := s.Get(ctx, a.ID); err != nil || j.State != ipp.JobCanceled {
		t.Fatalf("job done: %+v %v", j, err)
	}
	if _, err := s.st.getDoc(ctx, originalKey(77)); !isNotFound(err) {
		t.Fatalf("an orphan document is kept: %v", err)
	}
	if _, err := s.st.getDoc(ctx, originalKey(a.ID)); err != nil {
		t.Fatalf("the document of a job done is removed before its retention: %v", err)
	}
}

func TestCollectAndAccounting(t *testing.T) {
	f := newFixture(t)
	s := f.open()
	ctx := context.Background()
	doc := bdfBytes(t, testDoc(1))
	a := submit(t, s, Submission{Name: "a.bdf", User: "erp", Document: bytes.NewReader(doc), Ticket: print.Ticket{Sides: "two-sided-long-edge", Copies: 2}})
	if _, err := s.Cancel(ctx, a.ID, "job-canceled-by-operator"); err != nil {
		t.Fatal(err)
	}
	f.clock.Add(8 * 24 * time.Hour)
	b := submit(t, s, Submission{Name: "b.bdf", User: "web", Document: bytes.NewReader(doc)})
	if _, err := s.Cancel(ctx, b.ID, ""); err != nil {
		t.Fatal(err)
	}
	if err := s.Collect(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := s.st.getDoc(ctx, originalKey(a.ID)); !isNotFound(err) {
		t.Fatalf("a document past its retention: %v", err)
	}
	if _, err := s.st.getDoc(ctx, originalKey(b.ID)); err != nil {
		t.Fatalf("a document within its retention: %v", err)
	}
	if j, _ := s.Get(ctx, a.ID); !j.DocumentsRemoved {
		t.Fatal("the job does not say its documents are removed")
	}
	if _, err := s.Retry(ctx, a.ID); !errors.Is(err, ErrNotPossible) {
		t.Fatalf("retry without documents: %v", err)
	}
	if got := s.List(Filter{Done: true}); len(got) != 1 || got[0].ID != b.ID {
		t.Fatalf("recent jobs: %+v", got)
	}
	var out bytes.Buffer
	if err := s.Accounting(ctx, time.Time{}, time.Time{}, "csv", &out); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 3 || !strings.HasPrefix(lines[1], "1,") || !strings.Contains(lines[1], ",erp,a.bdf,canceled,job-canceled-by-operator,") || !strings.Contains(lines[1], ",two-sided-long-edge,,2,") {
		t.Fatalf("accounting:\n%s", out.String())
	}
	out.Reset()
	if err := s.Accounting(ctx, f.clock.Now().Add(-time.Hour), time.Time{}, "jsonl", &out); err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(out.String(), "\n"); n != 1 || !strings.Contains(out.String(), `"user":"web"`) {
		t.Fatalf("accounting from an hour ago:\n%s", out.String())
	}
	// the records past their retention
	f.clock.Add(400 * 24 * time.Hour)
	if err := s.Collect(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get(ctx, a.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("a record past its retention: %v", err)
	}
}

func TestBackupAndTrust(t *testing.T) {
	f := newFixture(t)
	s := f.open()
	ctx := context.Background()
	submit(t, s, Submission{Document: bytes.NewReader(bdfBytes(t, testDoc(1)))})
	be := trustBackend{s.st}
	if err := be.Save(map[string]string{"mfp1.example.local": "sha256:aa", "[fe80::1]": "sha256:bb"}); err != nil {
		t.Fatal(err)
	}
	if err := be.Save(map[string]string{"[fe80::1]": "sha256:cc"}); err != nil {
		t.Fatal(err)
	}
	m, err := be.Load()
	if err != nil || len(m) != 1 || m["[fe80::1]"] != "sha256:cc" {
		t.Fatalf("trust: %v %v", m, err)
	}
	var b bytes.Buffer
	if err := s.Backup(ctx, &b); err != nil || b.Len() == 0 {
		t.Fatalf("backup: %d bytes, %v", b.Len(), err)
	}
	copyPath := filepath.Join(t.TempDir(), "copy.db")
	if err := writeFile(copyPath, b.Bytes()); err != nil {
		t.Fatal(err)
	}
	f.o.Path = copyPath
	c := f.open()
	if got := c.List(Filter{}); len(got) != 1 {
		t.Fatalf("the backup has %d jobs", len(got))
	}
}

func TestHoldEnd(t *testing.T) {
	jst := time.FixedZone("JST", 9*3600)
	at := func(h int) time.Time { return time.Date(2026, 10, 20, h, 30, 0, 0, jst) } // a Tuesday
	for _, c := range []struct {
		hold string
		now  time.Time
		want time.Time
	}{
		{"", at(10), at(10)},
		{"indefinite", at(10), time.Time{}},
		{"night", at(10), time.Date(2026, 10, 20, 18, 0, 0, 0, jst)},
		{"night", at(23), at(23)},
		{"night", at(3), at(3)},
		{"day-time", at(20), time.Date(2026, 10, 21, 6, 0, 0, 0, jst)},
		{"third-shift", at(9), time.Date(2026, 10, 21, 0, 0, 0, 0, jst)},
		{"second-shift", at(17), at(17)},
		{"weekend", at(10), time.Date(2026, 10, 24, 0, 0, 0, 0, jst)},
	} {
		if got := holdEnd(c.hold, time.Time{}, c.now); !got.Equal(c.want) {
			t.Errorf("%s at %v: %v, want %v", c.hold, c.now, got, c.want)
		}
	}
	if got := holdEnd("night", at(12), at(10)); !got.Equal(at(12)) {
		t.Errorf("job-hold-until-time: %v", got)
	}
	if checkHold("lunch") == nil {
		t.Error("an unknown hold is taken")
	}
}

func TestTransitions(t *testing.T) {
	now := time.Now()
	j := Job{ID: 1, State: ipp.JobPending}
	for _, st := range []ipp.JobState{ipp.JobProcessing, ipp.JobProcessingStopped, ipp.JobProcessing, ipp.JobCompleted} {
		if err := j.to(now, st); err != nil {
			t.Fatal(err)
		}
	}
	if j.Finished.IsZero() {
		t.Fatal("a job completed has no end")
	}
	if err := j.to(now, ipp.JobPending); err == nil {
		t.Fatal("a job completed went back to pending")
	}
	h := Job{ID: 2, State: ipp.JobPendingHeld}
	if err := h.to(now, ipp.JobCompleted); err == nil {
		t.Fatal("a held job completed")
	}
	b, _ := json.Marshal(Job{ID: 3, State: ipp.JobProcessingStopped})
	if !bytes.Contains(b, []byte(`"state":"processing-stopped"`)) {
		t.Fatalf("JSON of a job: %s", b)
	}
}

func TestNameKey(t *testing.T) {
	for in, want := range map[string]string{"office": "office", "mfp1.example.local": "mfp1.example.local", "[fe80::1]": "_x5b666538303a3a315d",
		"受付": "_xe58f97e4bb98", "_xab": "_x5f786162", "__x__": "_x5f5f785f5f"} {
		if got := nameKey(in); got != want {
			t.Errorf("nameKey(%q) = %q, want %q", in, got, want)
		}
	}
}

// openLocked opens a spool's file as another process would.
func openLocked(path string, timeout time.Duration) error {
	db, err := bolt.Open(path, 0o600, &bolt.Options{Timeout: timeout})
	if err == nil {
		db.Close()
		return nil
	}
	return err
}

func writeFile(path string, b []byte) error { return os.WriteFile(path, b, 0o600) }

func TestOpenTwice(t *testing.T) {
	f := newFixture(t)
	f.open()
	if _, err := Open(context.Background(), f.o); err == nil {
		t.Fatal("the same spool opened twice in a process")
	}
}

func TestCompact(t *testing.T) {
	f := newFixture(t)
	s, err := Open(context.Background(), f.o)
	if err != nil {
		t.Fatal(err)
	}
	doc := bdfBytes(t, testDoc(1))
	var ids []int32
	for i := 0; i < 20; i++ {
		ids = append(ids, submit(t, s, Submission{Document: bytes.NewReader(doc)}).ID)
	}
	for _, id := range ids[1:] {
		if _, err := s.Cancel(context.Background(), id, ""); err != nil {
			t.Fatal(err)
		}
		if err := s.removeDocuments(context.Background(), id); err != nil {
			t.Fatal(err)
		}
	}
	if _, _, err := Compact(f.path); err == nil {
		t.Fatal("compacted a spool that is open")
	}
	s.Close()
	before, after, err := Compact(f.path)
	if err != nil {
		t.Fatal(err)
	}
	if after > before {
		t.Fatalf("%d → %d bytes", before, after)
	}
	s = f.open()
	if j, err := s.Get(context.Background(), ids[0]); err != nil || j.State != ipp.JobPending {
		t.Fatalf("after compaction: %v %+v", err, j)
	}
	if j, err := s.Get(context.Background(), ids[5]); err != nil || j.State != ipp.JobCanceled {
		t.Fatalf("after compaction: %v %+v", err, j)
	}
}
