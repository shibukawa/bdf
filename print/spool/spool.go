// Package spool keeps print jobs and sends them to printers: it takes a
// job and its document, keeps both across restarts, queues the job, sends
// it to a printer of its queue as print.Print does, follows it until the
// printer is done, retries what failed for want of the printer, and never
// prints a job twice after a crash. The states of the jobs are IPP's.
//
// One process owns a spool. Its records and documents are kept in one
// file through yuruidb (store.go).
package spool

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/ipp"
)

// Errors of the spool, which a server turns into HTTP statuses.
var (
	ErrInvalid       = errors.New("spool: invalid request")
	ErrNotFound      = errors.New("spool: no such job")
	ErrNoQueue       = errors.New("spool: no such queue")
	ErrNotAccepting  = errors.New("spool: the queue is not accepting jobs")
	ErrFull          = errors.New("spool: the spool is full")
	ErrTooLarge      = errors.New("spool: the document is too large")
	ErrBadDocument   = errors.New("spool: the document cannot be read")
	ErrUnknownFormat = errors.New("spool: the document is neither a bdf document nor a format the spool converts")
	ErrNotPossible   = errors.New("spool: not possible in the state of the job")
	ErrWrongPassword = errors.New("spool: the password does not open the document")
	// ErrPasswordRequired is returned by a Converter for an input that is
	// encrypted when no password was given.
	ErrPasswordRequired = errors.New("spool: the document is encrypted")
	ErrClosed           = errors.New("spool: closed")
)

// Converter turns documents that are not .bdf into documents: the spool
// does not depend on the converters, the program that runs it does.
type Converter interface {
	// Detect names the format of a document, or says it is unknown; name
	// is its file name, whose extension tells the formats whose content
	// does not (CSV, Markdown).
	Detect(r io.ReaderAt, size int64, name string) (format string, ok bool)
	// Convert makes a document of one in a format Detect named; a
	// password opens an encrypted one (ErrPasswordRequired and
	// ErrWrongPassword tell why it does not). The document of an input
	// that needed the password should be locked with it: the spool keeps
	// it so.
	Convert(ctx context.Context, r io.ReaderAt, size int64, format, name, password string) (*bdf.Document, []string, error)
}

// Queue is a queue of jobs and the printers its jobs are sent to.
type Queue struct {
	Name string
	// Printers are the names of the printers, in order; more than one
	// makes a class.
	Printers []string
	// Policy is "any" (the first printer free) or "failover" (the first
	// printer that is up, in order); "" is "failover".
	Policy string
	// Defaults are the ticket's defaults: a job's ticket asks over them.
	Defaults print.Ticket
	// Format and RasterIfLossy are print.Options'.
	Format        string
	RasterIfLossy bool
	// MaxInflight is how many jobs of the queue a printer is given before
	// the first is done: 2 by default, so that the next is there when
	// the first is printed.
	MaxInflight int
}

// Printer is a printer jobs are sent to.
type Printer struct {
	Name   string
	Client *ipp.Client
}

// Limits bound what the spool takes.
type Limits struct {
	// MaxDocument is the size of a document: 256 MiB by default.
	MaxDocument int64
	// MaxJobs and MaxBytes bound the jobs not done and their documents:
	// 10,000 and 10 GiB by default. Beyond them a job is refused (ErrFull).
	MaxJobs  int
	MaxBytes int64
	// RetryFor is how long a job is sent again to a printer that is not
	// reached before it is aborted: 24 hours by default.
	RetryFor time.Duration
	// Renders is how many jobs are converted or drawn at once: the
	// number of CPUs by default.
	Renders int
	// MaxImpressions is print.Options'.
	MaxImpressions int
}

// Retention is how long what is done is kept.
type Retention struct {
	// Documents of the jobs done: 7 days by default.
	Documents time.Duration
	// Records of the jobs done, which are the accounting: 400 days by
	// default.
	Records time.Duration
}

// Options are a spool's.
type Options struct {
	// Path is the file of the spool.
	Path string
	// Timeout is how long Open waits for the file's lock: 1 second by
	// default. Another process holding it is an error.
	Timeout time.Duration
	Queues  []Queue
	// Printers are the printers the queues name.
	Printers []Printer
	// Converter converts documents that are not .bdf; nil takes .bdf
	// only.
	Converter Converter
	Limits    Limits
	Retention Retention
	// TempDir holds documents while they arrive: os.TempDir by default.
	TempDir string
	// FontFS, FontDirs and NoSystemFonts are print.Options' and the
	// converter's.
	FontFS        fs.FS
	FontDirs      []string
	NoSystemFonts bool
	// Now is the clock; time.Now by default.
	Now func() time.Time
	// Log receives what happens to jobs and printers; nothing by default.
	Log *slog.Logger
}

// Spool is a spool, open.
type Spool struct {
	o    Options
	st   *store
	path string // the file, absolute

	mu        sync.Mutex
	jobs      map[int32]*entry
	recent    []Job // the jobs done within the retention of documents, oldest first
	queues    map[string]*queueState
	printers  map[string]*printerState
	passwords map[int32]string
	uploads   map[string]bool // keys of documents being stored
	subs      map[chan Event]struct{}
	stats     Stats
	closed    bool

	started bool
	kick    chan struct{}  // wakes the runner
	stop    chan struct{}  // closed by Close
	running sync.WaitGroup // the runner and the steps
	trust   *ipp.TrustStore
	renders chan struct{} // the slots of conversion and drawing
	// stepCtx ends the steps when the spool closes; sendCtx ends the
	// sends, which Close lets finish for a while first.
	stepCtx  context.Context
	stepStop context.CancelFunc
	sendCtx  context.Context
	sendStop context.CancelFunc
}

// entry is a job not done, as the spool holds it.
type entry struct {
	job Job
	rev string // the store's revision of job
	// busy is set while a step of the job runs; cancel stops it.
	busy   bool
	cancel context.CancelFunc
	// canceling asks the step to cancel the job (at the printer too).
	canceling bool
	// recover marks a job that was at a printer when the spool stopped:
	// it is checked against the printer before anything else.
	recover bool
}

type queueState struct {
	Queue
	rec queueRecord
	rev string
}

// Event tells a change: of a job, of a queue, or of a printer.
type Event struct {
	Job     *Job
	Queue   *QueueInfo
	Printer *PrinterInfo
}

// Open opens a spool, recovering what a crash left: documents that were
// arriving, jobs recorded as done but not removed, and jobs that were at a
// printer, which are checked against the printer before anything is sent
// again (Start sends them).
func Open(ctx context.Context, o Options) (*Spool, error) {
	if o.Path == "" {
		return nil, fmt.Errorf("%w: no path", ErrInvalid)
	}
	if o.Timeout == 0 {
		o.Timeout = time.Second
	}
	if o.Now == nil {
		o.Now = time.Now
	}
	if o.Log == nil {
		o.Log = slog.New(slog.DiscardHandler)
	}
	l := &o.Limits
	if l.MaxDocument == 0 {
		l.MaxDocument = 256 << 20
	}
	if l.MaxJobs == 0 {
		l.MaxJobs = 10000
	}
	if l.MaxBytes == 0 {
		l.MaxBytes = 10 << 30
	}
	if l.RetryFor == 0 {
		l.RetryFor = 24 * time.Hour
	}
	if l.Renders <= 0 {
		l.Renders = max(numCPU(), 1)
	}
	if o.Retention.Documents == 0 {
		o.Retention.Documents = 7 * 24 * time.Hour
	}
	if o.Retention.Records == 0 {
		o.Retention.Records = 400 * 24 * time.Hour
	}
	s := &Spool{o: o, jobs: map[int32]*entry{}, queues: map[string]*queueState{}, printers: map[string]*printerState{},
		passwords: map[int32]string{}, uploads: map[string]bool{}, subs: map[chan Event]struct{}{}, stats: newStats(),
		kick: make(chan struct{}, 1), stop: make(chan struct{}), renders: make(chan struct{}, l.Renders)}
	s.stepCtx, s.stepStop = context.WithCancel(context.Background())
	s.sendCtx, s.sendStop = context.WithCancel(context.Background())
	for _, p := range o.Printers {
		if p.Name == "" || p.Client == nil {
			return nil, fmt.Errorf("%w: a printer without name or client", ErrInvalid)
		}
		if _, dup := s.printers[p.Name]; dup {
			return nil, fmt.Errorf("%w: printer %q twice", ErrInvalid, p.Name)
		}
		s.printers[p.Name] = newPrinterState(p)
	}
	for _, q := range o.Queues {
		if q.Name == "" {
			return nil, fmt.Errorf("%w: a queue without name", ErrInvalid)
		}
		if _, dup := s.queues[q.Name]; dup {
			return nil, fmt.Errorf("%w: queue %q twice", ErrInvalid, q.Name)
		}
		if len(q.Printers) == 0 {
			return nil, fmt.Errorf("%w: queue %q has no printer", ErrInvalid, q.Name)
		}
		for _, p := range q.Printers {
			if s.printers[p] == nil {
				return nil, fmt.Errorf("%w: queue %q names printer %q, which is not configured", ErrInvalid, q.Name, p)
			}
		}
		switch q.Policy {
		case "":
			q.Policy = "failover"
		case "any", "failover":
		default:
			return nil, fmt.Errorf("%w: queue %q: policy %q: any or failover", ErrInvalid, q.Name, q.Policy)
		}
		if q.MaxInflight <= 0 {
			q.MaxInflight = 2
		}
		if err := q.Defaults.Check(); err != nil {
			return nil, fmt.Errorf("%w: queue %q: %v", ErrInvalid, q.Name, err)
		}
		s.queues[q.Name] = &queueState{Queue: q}
	}
	path, err := claim(o.Path)
	if err != nil {
		return nil, err
	}
	s.path = path
	st, err := openStore(ctx, o.Path, o.Timeout)
	if err != nil {
		unclaim(path)
		return nil, err
	}
	s.st = st
	s.trust = &ipp.TrustStore{Backend: trustBackend{st}}
	if err := s.load(ctx); err != nil {
		st.Close()
		unclaim(path)
		return nil, err
	}
	return s, nil
}

// open are the files of the spools open in this process: a file is
// shared within a process, so a second spool of it is refused here (the
// file's lock refuses another process).
var open = struct {
	sync.Mutex
	paths map[string]bool
}{paths: map[string]bool{}}

func claim(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	// the links are resolved, of the directory when the file is not made
	// yet, so that a path and the path of its target are one spool
	if r, err := filepath.EvalSymlinks(abs); err == nil {
		abs = r
	} else if dir, err := filepath.EvalSymlinks(filepath.Dir(abs)); err == nil {
		abs = filepath.Join(dir, filepath.Base(abs))
	}
	open.Lock()
	defer open.Unlock()
	if open.paths[abs] {
		return "", fmt.Errorf("spool: %s is open already", path)
	}
	open.paths[abs] = true
	return abs, nil
}

func unclaim(path string) {
	open.Lock()
	defer open.Unlock()
	delete(open.paths, path)
}

// load reads the jobs and the queues' states, and puts right what a crash
// between two writes left.
func (s *Spool) load(ctx context.Context) error {
	for name, q := range s.queues {
		rec, rev, err := s.st.queue(ctx, name)
		if err != nil {
			return err
		}
		q.rec, q.rev = rec, rev
	}
	live, err := s.st.liveJobs(ctx)
	if err != nil {
		return err
	}
	now := s.o.Now()
	for _, r := range live {
		j := r.job
		if _, _, err := s.st.getDone(ctx, j.ID); err == nil {
			// recorded as done, not yet removed from the jobs
			if err := s.st.deleteLive(ctx, j.ID, r.rev); err != nil {
				return err
			}
			continue
		} else if !isNotFound(err) {
			return err
		}
		if s.queues[j.Queue] == nil {
			// a queue no longer configured: its jobs wait for one, held
			if j.State != ipp.JobPendingHeld && !j.State.Done() {
				s.o.Log.Warn("job of a queue that is not configured", "job", j.ID, "queue", j.Queue)
			}
		}
		s.jobs[j.ID] = &entry{job: j, rev: r.rev}
	}
	// the jobs done recently, for the lists
	cut := now.Add(-s.o.Retention.Documents)
	err = s.st.eachDone(ctx, func(j Job, _ string) bool {
		if j.Finished.After(cut) {
			s.recent = append(s.recent, j)
		}
		return true
	})
	if err != nil {
		return err
	}
	sort.SliceStable(s.recent, func(a, b int) bool { return s.recent[a].Finished.Before(s.recent[b].Finished) })
	// documents of no job: a crash while one arrived
	return s.dropOrphans(ctx)
}

// Close stops the spool as Shutdown does, giving sends 30 seconds.
func (s *Spool) Close() error { return s.Shutdown(30 * time.Second) }

// Shutdown stops the spool: no step starts, the steps that only wait stop
// at once, and the documents being sent to printers are given grace to
// finish (a send cut is put right when the spool opens again). Then the
// file is closed.
func (s *Spool) Shutdown(grace time.Duration) error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil
	}
	s.closed = true
	close(s.stop)
	for ch := range s.subs {
		close(ch)
	}
	s.subs = nil
	s.mu.Unlock()
	s.stepStop()
	t := time.AfterFunc(grace, s.sendStop)
	s.running.Wait()
	t.Stop()
	s.sendStop()
	err := s.st.Close()
	unclaim(s.path)
	return err
}

// TrustStore keeps the fingerprints of the printers' certificates in the
// spool, for the printers' clients.
func (s *Spool) TrustStore() *ipp.TrustStore { return s.trust }

// Submission is a job to submit.
type Submission struct {
	Queue string
	// Name is the document's file name; User who submits it.
	Name, User string
	// Document is read to its end by Submit; or Upload, a document
	// received before (Receive), which Submit takes.
	Document io.Reader
	Upload   *Upload
	Ticket   print.Ticket
	Options  JobOptions
	// Priority is 1 to 100; 0 is 50.
	Priority int
	// HoldUntil and HoldUntilTime hold the job (job-hold-until).
	HoldUntil     string
	HoldUntilTime time.Time
	// Password opens an encrypted document. It is kept in memory only:
	// after a restart such a job waits for it again (job-password-wait).
	Password string
}

// Submit takes a job: it stores its document and records the job, which
// is then queued (or held). The document is read to its end before the
// job is recorded.
func (s *Spool) Submit(ctx context.Context, sub Submission) (Job, error) {
	if sub.Upload != nil {
		defer sub.Upload.Discard()
	} else if sub.Document == nil {
		return Job{}, fmt.Errorf("%w: no document", ErrInvalid)
	}
	if sub.Priority == 0 {
		sub.Priority = 50
	}
	if sub.Priority < 1 || sub.Priority > 100 {
		return Job{}, fmt.Errorf("%w: priority %d: 1 to 100", ErrInvalid, sub.Priority)
	}
	if err := sub.Ticket.Check(); err != nil {
		return Job{}, fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	if err := checkHold(sub.HoldUntil); err != nil {
		return Job{}, err
	}
	if err := checkOptions(sub.Options); err != nil {
		return Job{}, err
	}
	if err := s.accepts(sub.Queue); err != nil {
		return Job{}, err
	}
	if sub.Name == "" {
		sub.Name = "document"
	}
	up := sub.Upload
	if up == nil {
		var err error
		if up, err = s.Receive(ctx, sub.Document); err != nil {
			return Job{}, err
		}
		defer up.Discard()
	}
	f, size, sum := up.f, up.size, up.sum
	if err := s.room(size); err != nil {
		return Job{}, err
	}
	doc := Document{Name: sub.Name, Size: size, SHA256: sum}
	head := make([]byte, 4)
	if n, _ := f.ReadAt(head, 0); n == 4 && bytes.Equal(head, bdf.Magic[:]) {
		doc.Format = "bdf"
		r, err := bdf.OpenSingle(f, size)
		if err != nil {
			return Job{}, fmt.Errorf("%w: %v", ErrBadDocument, err)
		}
		if doc.Encrypted = r.Encrypted(); doc.Encrypted && sub.Password != "" {
			if err := r.Unlock(sub.Password); errors.Is(err, bdf.ErrWrongPassword) {
				return Job{}, ErrWrongPassword
			} else if err != nil {
				return Job{}, fmt.Errorf("%w: %v", ErrBadDocument, err)
			}
		}
		r.Close()
	} else {
		if s.o.Converter == nil {
			return Job{}, ErrUnknownFormat
		}
		format, ok := s.o.Converter.Detect(f, size, sub.Name)
		if !ok {
			return Job{}, ErrUnknownFormat
		}
		doc.Format = format
	}
	ctx = context.WithoutCancel(ctx) // a client that goes away does not leave half a job
	id, err := s.st.nextID(ctx)
	if err != nil {
		return Job{}, err
	}
	key := originalKey(id)
	s.mu.Lock()
	s.uploads[key] = true
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		delete(s.uploads, key)
		s.mu.Unlock()
	}()
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return Job{}, err
	}
	if _, err := s.st.putDoc(ctx, key, f); err != nil {
		return Job{}, err
	}
	now := s.o.Now()
	if sub.Ticket.JobName == "" {
		sub.Ticket.JobName = sub.Name
	}
	j := Job{ID: id, UUID: newUUID(), Queue: sub.Queue, Name: sub.Name, User: sub.User, Document: doc, Ticket: sub.Ticket,
		Options: sub.Options, State: ipp.JobPending, Priority: sub.Priority, HoldUntil: sub.HoldUntil, HoldUntilTime: sub.HoldUntilTime,
		Step: StepSend, Submitted: now}
	if doc.Format != "bdf" {
		j.Step = StepConvert
	}
	switch end := holdEnd(sub.HoldUntil, sub.HoldUntilTime, now); {
	case end.IsZero() || end.After(now):
		j.State, j.Reasons, j.Due = ipp.JobPendingHeld, []string{"job-hold-until-specified"}, end
	case doc.Encrypted && sub.Password == "":
		j.State, j.Reasons = ipp.JobPendingHeld, []string{"job-password-wait"}
	default:
		j.Reasons, j.Due = []string{"job-queued"}, now
	}
	rev, err := s.st.insertLive(ctx, &j)
	if err != nil {
		s.st.removeDoc(ctx, key)
		return Job{}, err
	}
	s.mu.Lock()
	s.jobs[id] = &entry{job: j, rev: rev}
	if sub.Password != "" {
		s.passwords[id] = sub.Password
	}
	s.mu.Unlock()
	s.o.Log.Info("job submitted", "job", id, "queue", j.Queue, "user", j.User, "name", j.Name, "format", doc.Format, "size", size)
	s.notify(&j)
	s.wake()
	return j, nil
}

func checkOptions(o JobOptions) error {
	switch o.Format {
	case "", "auto", print.FormatPDF, print.FormatPWGRaster, "pdf", "pwg-raster":
	default:
		return fmt.Errorf("%w: format %q", ErrInvalid, o.Format)
	}
	return nil
}

// accepts refuses a queue that does not exist or does not accept jobs.
func (s *Spool) accepts(queue string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	q := s.queues[queue]
	if q == nil {
		return fmt.Errorf("%w: %q", ErrNoQueue, queue)
	}
	if q.rec.Rejecting {
		return fmt.Errorf("%w: %s", ErrNotAccepting, queue)
	}
	return nil
}

// room refuses a document beyond the spool's limits.
func (s *Spool) room(size int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.jobs) >= s.o.Limits.MaxJobs {
		return fmt.Errorf("%w: %d jobs", ErrFull, len(s.jobs))
	}
	var total int64
	for _, e := range s.jobs {
		total += e.job.Document.Size
		if e.job.Document.BDF != nil {
			total += e.job.Document.BDF.Size
		}
	}
	if total+size > s.o.Limits.MaxBytes {
		return fmt.Errorf("%w: %d bytes of documents", ErrFull, total)
	}
	return nil
}

// Upload is a document received before its job is submitted: a server
// reads the document of a request before the fields that follow it.
type Upload struct {
	f    *os.File
	size int64
	sum  string
}

// Size is the document's size.
func (u *Upload) Size() int64 { return u.size }

// Discard drops a document no job took.
func (u *Upload) Discard() {
	if u.f != nil {
		u.f.Close()
		os.Remove(u.f.Name())
		u.f = nil
	}
}

// Receive reads a document to a temporary file (Options.TempDir),
// measuring it, for Submission.Upload.
func (s *Spool) Receive(ctx context.Context, r io.Reader) (*Upload, error) {
	f, size, sum, err := s.receive(ctx, r)
	if err != nil {
		return nil, err
	}
	return &Upload{f, size, sum}, nil
}

// receive reads a document to a temporary file, measuring it.
func (s *Spool) receive(ctx context.Context, r io.Reader) (*os.File, int64, string, error) {
	f, err := os.CreateTemp(s.o.TempDir, "bdf-spool-*")
	if err != nil {
		return nil, 0, "", err
	}
	h := sha256.New()
	n, err := io.Copy(io.MultiWriter(f, h), io.LimitReader(ctxReader{ctx, r}, s.o.Limits.MaxDocument+1))
	if err == nil && n > s.o.Limits.MaxDocument {
		err = fmt.Errorf("%w: more than %d bytes", ErrTooLarge, s.o.Limits.MaxDocument)
	}
	if err == nil && n == 0 {
		err = fmt.Errorf("%w: the document is empty", ErrBadDocument)
	}
	if err != nil {
		f.Close()
		os.Remove(f.Name())
		return nil, 0, "", err
	}
	return f, n, hex.EncodeToString(h.Sum(nil)), nil
}

type ctxReader struct {
	ctx context.Context
	r   io.Reader
}

func (r ctxReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.r.Read(p)
}

// Get returns a job, done or not.
func (s *Spool) Get(ctx context.Context, id int32) (Job, error) {
	s.mu.Lock()
	if e := s.jobs[id]; e != nil {
		j := e.job.clone()
		s.mu.Unlock()
		return j, nil
	}
	s.mu.Unlock()
	j, _, err := s.st.getDone(ctx, id)
	if isNotFound(err) {
		return Job{}, fmt.Errorf("%w: %d", ErrNotFound, id)
	}
	return j, err
}

// Filter selects jobs to list.
type Filter struct {
	// Queue and User select the jobs of a queue, of a user.
	Queue, User string
	// Done lists the jobs done within the retention of documents too.
	Done bool
	// After lists the jobs of greater ids; Limit how many (100 by
	// default, at most 1000).
	After int32
	Limit int
}

// List returns jobs in the order of their ids.
func (s *Spool) List(f Filter) []Job {
	if f.Limit <= 0 {
		f.Limit = 100
	}
	f.Limit = min(f.Limit, 1000)
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Job
	keep := func(j *Job) bool {
		return j.ID > f.After && (f.Queue == "" || j.Queue == f.Queue) && (f.User == "" || j.User == f.User)
	}
	for _, e := range s.jobs {
		if keep(&e.job) {
			out = append(out, e.job.clone())
		}
	}
	if f.Done {
		for i := range s.recent {
			if keep(&s.recent[i]) {
				out = append(out, s.recent[i].clone())
			}
		}
	}
	slices.SortFunc(out, func(a, b Job) int { return int(a.ID) - int(b.ID) })
	if len(out) > f.Limit {
		out = out[:f.Limit]
	}
	return out
}

// Cancel cancels a job. A job at a printer is canceled there too.
func (s *Spool) Cancel(ctx context.Context, id int32, reason string) (Job, error) {
	if reason == "" {
		reason = "job-canceled-by-user"
	}
	s.mu.Lock()
	e := s.jobs[id]
	if e == nil {
		s.mu.Unlock()
		if _, err := s.Get(ctx, id); err != nil {
			return Job{}, err
		}
		return Job{}, fmt.Errorf("%w: job %d is done", ErrNotPossible, id)
	}
	if e.busy {
		// the step running cancels the job at the printer and records it
		e.canceling = true
		if e.cancel != nil {
			e.cancel()
		}
		j := e.job.clone()
		s.mu.Unlock()
		return j, nil
	}
	e.busy = true // no step starts while the job is canceled
	j := e.job.clone()
	s.mu.Unlock()
	if j.Printer != nil && j.Printer.JobID != 0 {
		s.cancelAtPrinter(ctx, &j)
	}
	if err := j.to(s.o.Now(), ipp.JobCanceled, reason); err != nil {
		s.release(e)
		return Job{}, err
	}
	if err := s.finish(ctx, e, j); err != nil {
		s.release(e)
		return Job{}, err
	}
	s.o.Log.Info("job canceled", "job", id, "reason", reason)
	return j, nil
}

// cancelAtPrinter asks the printer of a job to cancel it, as best it can:
// a printer that is not reached is told nothing.
func (s *Spool) cancelAtPrinter(ctx context.Context, j *Job) {
	p := s.printer(j.Printer.Name)
	if p == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer cancel()
	if err := p.client.CancelJob(ctx, j.Printer.JobID); err != nil {
		s.o.Log.Warn("cancel at the printer", "job", j.ID, "printer", p.Name, "err", err)
	}
}

// change applies f to a job that no step runs, and records the result.
func (s *Spool) change(ctx context.Context, id int32, f func(j *Job, now time.Time) error) (Job, error) {
	s.mu.Lock()
	e := s.jobs[id]
	if e == nil {
		s.mu.Unlock()
		if _, err := s.Get(ctx, id); err != nil {
			return Job{}, err
		}
		return Job{}, fmt.Errorf("%w: job %d is done", ErrNotPossible, id)
	}
	if e.busy {
		s.mu.Unlock()
		return Job{}, fmt.Errorf("%w: job %d is being processed", ErrNotPossible, id)
	}
	e.busy = true
	j := e.job.clone()
	s.mu.Unlock()
	defer s.release(e)
	if err := f(&j, s.o.Now()); err != nil {
		return Job{}, err
	}
	if err := s.save(ctx, e, j); err != nil {
		return Job{}, err
	}
	s.wake()
	return j, nil
}

// Hold holds a job that is pending until it is released.
func (s *Spool) Hold(ctx context.Context, id int32) (Job, error) {
	return s.change(ctx, id, func(j *Job, now time.Time) error {
		if j.State != ipp.JobPending && j.State != ipp.JobPendingHeld {
			return fmt.Errorf("%w: job %d is %s", ErrNotPossible, id, j.State)
		}
		j.HoldUntil, j.HoldUntilTime, j.Due = "indefinite", time.Time{}, time.Time{}
		return j.to(now, ipp.JobPendingHeld, "job-hold-until-specified")
	})
}

// Release releases a held job; a password opens its document when it
// waits for one.
func (s *Spool) Release(ctx context.Context, id int32, password string) (Job, error) {
	if password != "" {
		if err := s.checkPassword(ctx, id, password); err != nil {
			return Job{}, err
		}
	}
	return s.change(ctx, id, func(j *Job, now time.Time) error {
		if j.State != ipp.JobPendingHeld {
			return fmt.Errorf("%w: job %d is %s", ErrNotPossible, id, j.State)
		}
		if slices.Contains(j.Reasons, "job-password-wait") && !s.hasPassword(id) {
			return fmt.Errorf("%w: job %d waits for its password", ErrNotPossible, id)
		}
		j.HoldUntil, j.HoldUntilTime, j.Due = "", time.Time{}, now
		return j.to(now, ipp.JobPending, "job-queued")
	})
}

// SetPassword gives the password of an encrypted document: a job waiting
// for it is queued.
func (s *Spool) SetPassword(ctx context.Context, id int32, password string) (Job, error) {
	if err := s.checkPassword(ctx, id, password); err != nil {
		return Job{}, err
	}
	return s.change(ctx, id, func(j *Job, now time.Time) error {
		if j.State == ipp.JobPendingHeld && slices.Contains(j.Reasons, "job-password-wait") {
			j.Due = now
			return j.to(now, ipp.JobPending, "job-queued")
		}
		return nil
	})
}

func (s *Spool) hasPassword(id int32) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.passwords[id] != ""
}

func (s *Spool) password(id int32) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.passwords[id]
}

// checkPassword keeps a password that opens a job's .bdf document; the
// password of a document converted is checked when it is converted.
func (s *Spool) checkPassword(ctx context.Context, id int32, password string) error {
	if password == "" {
		return fmt.Errorf("%w: no password", ErrInvalid)
	}
	j, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if j.State.Done() {
		return fmt.Errorf("%w: job %d is done", ErrNotPossible, id)
	}
	if j.Document.Format == "bdf" {
		if !j.Document.Encrypted {
			return fmt.Errorf("%w: the document of job %d is not encrypted", ErrInvalid, id)
		}
		r, err := s.st.getDoc(ctx, originalKey(id))
		if err != nil {
			return err
		}
		defer r.Close()
		br, err := bdf.OpenSingle(r, r.Size)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrBadDocument, err)
		}
		if err := br.Unlock(password); errors.Is(err, bdf.ErrWrongPassword) {
			return ErrWrongPassword
		} else if err != nil {
			return err
		}
	}
	s.mu.Lock()
	s.passwords[id] = password
	s.mu.Unlock()
	return nil
}

// Retry queues again a job that was canceled or aborted, as a new attempt.
func (s *Spool) Retry(ctx context.Context, id int32) (Job, error) {
	s.mu.Lock()
	if e := s.jobs[id]; e != nil {
		st := e.job.State
		s.mu.Unlock()
		return Job{}, fmt.Errorf("%w: job %d is %s", ErrNotPossible, id, st)
	}
	s.mu.Unlock()
	ctx = context.WithoutCancel(ctx)
	j, rev, err := s.st.getDone(ctx, id)
	if isNotFound(err) {
		return Job{}, fmt.Errorf("%w: %d", ErrNotFound, id)
	} else if err != nil {
		return Job{}, err
	}
	if j.State != ipp.JobCanceled && j.State != ipp.JobAborted {
		return Job{}, fmt.Errorf("%w: job %d is %s", ErrNotPossible, id, j.State)
	}
	if j.DocumentsRemoved {
		return Job{}, fmt.Errorf("%w: the documents of job %d are no longer kept", ErrNotPossible, id)
	}
	if s.queues[j.Queue] == nil {
		return Job{}, fmt.Errorf("%w: %q", ErrNoQueue, j.Queue)
	}
	now := s.o.Now()
	if err := j.to(now, ipp.JobPending, "job-queued"); err != nil {
		return Job{}, err
	}
	j.Printer, j.RetrySince, j.Due, j.Message = nil, time.Time{}, now, ""
	j.Step = StepSend
	if j.Document.Format != "bdf" && j.Document.BDF == nil {
		j.Step = StepConvert
	}
	if j.Document.Encrypted && !s.hasPassword(id) {
		j.State, j.Reasons = ipp.JobPendingHeld, []string{"job-password-wait"}
	}
	nrev, err := s.st.insertLive(ctx, &j)
	if err != nil {
		return Job{}, err
	}
	if err := s.st.deleteDone(ctx, id, rev); err != nil && !isNotFound(err) {
		return Job{}, err
	}
	s.mu.Lock()
	s.jobs[id] = &entry{job: j, rev: nrev}
	s.recent = slices.DeleteFunc(s.recent, func(r Job) bool { return r.ID == id })
	s.mu.Unlock()
	s.o.Log.Info("job retried", "job", id)
	s.notify(&j)
	s.wake()
	return j, nil
}

// Move moves a job that is not processed yet to another queue.
func (s *Spool) Move(ctx context.Context, id int32, queue string) (Job, error) {
	s.mu.Lock()
	ok := s.queues[queue] != nil
	s.mu.Unlock()
	if !ok {
		return Job{}, fmt.Errorf("%w: %q", ErrNoQueue, queue)
	}
	return s.change(ctx, id, func(j *Job, _ time.Time) error {
		if j.State != ipp.JobPending && j.State != ipp.JobPendingHeld {
			return fmt.Errorf("%w: job %d is %s", ErrNotPossible, id, j.State)
		}
		j.Queue = queue
		return nil
	})
}

// SetPriority changes the priority of a job not processed yet.
func (s *Spool) SetPriority(ctx context.Context, id int32, priority int) (Job, error) {
	if priority < 1 || priority > 100 {
		return Job{}, fmt.Errorf("%w: priority %d: 1 to 100", ErrInvalid, priority)
	}
	return s.change(ctx, id, func(j *Job, _ time.Time) error {
		if j.State != ipp.JobPending && j.State != ipp.JobPendingHeld {
			return fmt.Errorf("%w: job %d is %s", ErrNotPossible, id, j.State)
		}
		j.Priority = priority
		return nil
	})
}

// save records a job that a step or an operation changed. The spool is
// the only writer: a conflict is a bug, or another process.
func (s *Spool) save(ctx context.Context, e *entry, j Job) error {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
	defer cancel()
	s.mu.Lock()
	rev := e.rev
	s.mu.Unlock()
	nrev, err := s.st.updateLive(ctx, &j, rev)
	if err != nil {
		if isConflict(err) {
			err = fmt.Errorf("spool: job %d changed in the store under the spool: %w", j.ID, err)
		}
		return err
	}
	s.mu.Lock()
	e.job, e.rev = j.clone(), nrev // the step goes on changing its own copy
	s.mu.Unlock()
	s.notify(&j)
	return nil
}

// finish records a job done: it is added to the jobs done, then removed
// from the jobs not done (Open puts right a crash between the two).
func (s *Spool) finish(ctx context.Context, e *entry, j Job) error {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
	defer cancel()
	if err := s.st.insertDone(ctx, &j); err != nil {
		return err
	}
	s.mu.Lock()
	rev := e.rev
	s.mu.Unlock()
	if err := s.st.deleteLive(ctx, j.ID, rev); err != nil && !isNotFound(err) {
		return err
	}
	s.mu.Lock()
	delete(s.jobs, j.ID)
	delete(s.passwords, j.ID)
	s.recent = append(s.recent, j.clone())
	s.stats.Done[QueueState{j.Queue, j.State.String()}]++
	s.stats.Sheets[j.Queue] += int64(j.Sheets)
	e.job = j.clone()
	s.mu.Unlock()
	s.notify(&j)
	s.wake()
	return nil
}

// release marks a job's step over.
func (s *Spool) release(e *entry) {
	s.mu.Lock()
	e.busy, e.cancel, e.canceling = false, nil, false
	s.mu.Unlock()
}

// Subscribe returns a channel of the changes of jobs and queues, and a
// function that ends the subscription. Changes are dropped when the
// channel is full.
func (s *Spool) Subscribe(buffer int) (<-chan Event, func()) {
	ch := make(chan Event, max(buffer, 1))
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		close(ch)
		return ch, func() {}
	}
	s.subs[ch] = struct{}{}
	return ch, func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if _, ok := s.subs[ch]; ok {
			delete(s.subs, ch)
			close(ch)
		}
	}
}

func (s *Spool) notify(j *Job) {
	c := j.clone()
	s.publish(Event{Job: &c})
}

func (s *Spool) publish(ev Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for ch := range s.subs {
		select {
		case ch <- ev:
		default:
		}
	}
}

// wake tells the runner that something changed.
func (s *Spool) wake() {
	select {
	case s.kick <- struct{}{}:
	default:
	}
}

// Backup writes a consistent copy of the spool's file.
func (s *Spool) Backup(ctx context.Context, w io.Writer) error {
	return s.st.backup(ctx, w)
}

func newUUID() string {
	var b [16]byte
	rand.Read(b[:])
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	h := hex.EncodeToString(b[:])
	return "urn:uuid:" + h[:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:]
}
