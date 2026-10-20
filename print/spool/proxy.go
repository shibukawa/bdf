package spool

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print/ipp"
)

// A printer of Proxy is one that proxies fetch jobs from, as IPP's
// Shared Infrastructure Extensions (PWG 5100.18) have it: a branch's print
// server, behind a NAT, fetches the jobs of its queue from the server of
// the head office and reports them back. The spool sends nothing to such
// a printer: a job for it waits to be fetched (StepFetch), is bound to the
// proxy that acknowledges it, and once the proxy has its document
// (AcknowledgeDocument) waits for its reports (StepRemote). A job bound to
// a proxy is never offered to another, so that it is not printed twice;
// Retry offers it again.

// Steps of the jobs of proxies.
const (
	// StepFetch waits for a proxy to fetch the job.
	StepFetch Step = "fetch"
	// StepRemote waits for the proxy that has the job to report it.
	StepRemote Step = "remote"
)

// proxyStale is how long after its last report a proxy is taken as not
// there.
var proxyStale = 2 * time.Minute

// Origin is where a job a proxy fetched comes from: the upstream (the
// name of its configuration), the job's id and uuid there.
type Origin struct {
	Upstream string `json:"upstream"`
	JobID    int32  `json:"job-id"`
	UUID     string `json:"uuid"`
}

// offer makes a job one that the proxies of a printer may fetch.
func (s *Spool) offer(ctx context.Context, e *entry, p *printerState) {
	j := s.snapshot(e)
	now := s.o.Now()
	j.Printer = &PrinterJob{Name: p.Name, URL: "proxy:" + p.Name, AttemptStarted: now}
	j.Step, j.Message = StepFetch, ""
	j.Reasons = []string{"job-fetchable"}
	s.commit(ctx, e, j)
	s.o.Log.Info("job fetchable", "job", j.ID, "printer", p.Name)
}

// Fetchable are the jobs the proxy of device may fetch for a printer: the
// pending jobs offered to it and bound to no proxy, and those bound to
// this one whose document it has not acknowledged (it fetches them again
// after it stopped).
func (s *Spool) Fetchable(printer, device string) []Job {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Job
	for _, e := range s.jobs {
		j := &e.job
		if q := s.queues[j.Queue]; q == nil || q.rec.Paused {
			continue
		}
		if j.Step == StepFetch && j.State == ipp.JobPending && j.Printer != nil && j.Printer.Name == printer &&
			(j.Printer.Device == "" || j.Printer.Device == device) {
			out = append(out, j.clone())
		}
	}
	slices.SortFunc(out, func(a, b Job) int {
		if a.Priority != b.Priority {
			return b.Priority - a.Priority
		}
		return int(a.ID - b.ID)
	})
	return out
}

// proxyJob runs f on a job of a printer of proxies that no step runs, and
// records it.
func (s *Spool) proxyJob(ctx context.Context, id int32, printer string, f func(j *Job, now time.Time) error) (Job, error) {
	s.mu.Lock()
	e := s.jobs[id]
	if e == nil {
		s.mu.Unlock()
		j, err := s.Get(ctx, id)
		if err != nil {
			return Job{}, err
		}
		return j, fmt.Errorf("%w: job %d is %s", ErrNotPossible, id, j.State)
	}
	if e.job.Printer == nil || e.job.Printer.Name != printer {
		s.mu.Unlock()
		return Job{}, fmt.Errorf("%w: job %d is not for %s", ErrNotFound, id, printer)
	}
	if e.busy {
		s.mu.Unlock()
		return Job{}, fmt.Errorf("%w: job %d is being processed", ErrNotPossible, id)
	}
	e.busy = true
	j := e.job.clone()
	s.mu.Unlock()
	// a Cancel while f ran waits for it: done after
	defer func() {
		s.mu.Lock()
		canceling := e.canceling
		cur := e.job.clone()
		s.mu.Unlock()
		if canceling && !cur.State.Done() {
			if cur.State == ipp.JobPendingHeld {
				cur.State = ipp.JobPending
			}
			if cur.to(s.o.Now(), ipp.JobCanceled, "job-canceled-by-user") == nil {
				s.commitErr(ctx, e, cur)
				s.o.Log.Info("job canceled", "job", cur.ID)
			}
		}
		s.release(e)
	}()
	if err := f(&j, s.o.Now()); err != nil {
		return j, err
	}
	if err := s.commitErr(ctx, e, j); err != nil {
		return Job{}, err
	}
	s.wake()
	return j, nil
}

func (s *Spool) commitErr(ctx context.Context, e *entry, j Job) error {
	ctx = context.WithoutCancel(ctx)
	if j.State.Done() {
		return s.finish(ctx, e, j)
	}
	return s.save(ctx, e, j)
}

// AcknowledgeJob binds a fetchable job to the proxy of device: no other
// proxy is given it.
func (s *Spool) AcknowledgeJob(ctx context.Context, id int32, printer, device string) (Job, error) {
	return s.proxyJob(ctx, id, printer, func(j *Job, now time.Time) error {
		if j.Step != StepFetch || j.State != ipp.JobPending || j.Printer.Device != "" && j.Printer.Device != device {
			return fmt.Errorf("%w: job %d is not fetchable", ErrNotPossible, id)
		}
		j.Printer.Device = device
		if j.Started.IsZero() {
			j.Started = now
		}
		j.Attempts++
		return j.to(now, ipp.JobPending, "job-queued")
	})
}

// AcknowledgeDocument records that the proxy of device has the document
// of a job it acknowledged: the job waits for its reports.
func (s *Spool) AcknowledgeDocument(ctx context.Context, id int32, printer, device string) (Job, error) {
	return s.proxyJob(ctx, id, printer, func(j *Job, now time.Time) error {
		if j.Printer.Device != device {
			return fmt.Errorf("%w: job %d is another proxy's", ErrNotPossible, id)
		}
		if j.Step == StepRemote {
			return nil
		}
		if j.Step != StepFetch {
			return fmt.Errorf("%w: job %d is not fetched", ErrNotPossible, id)
		}
		j.Printer.Sent, j.Step = true, StepRemote
		return j.to(now, ipp.JobProcessing, "job-printing")
	})
}

// ProxyReport records what the proxy of device reports of a job: its
// state at the proxy, the reasons, the impressions, a message. The job
// ends when the proxy's does. A job already ended here (canceled) is
// returned as it is, with ErrNotPossible, for the proxy to end its own.
func (s *Spool) ProxyReport(ctx context.Context, id int32, printer, device string, st ipp.JobState, reasons []string, impressions int, message string) (Job, error) {
	return s.proxyJob(ctx, id, printer, func(j *Job, now time.Time) error {
		if j.Printer.Device != device {
			return fmt.Errorf("%w: job %d is another proxy's", ErrNotPossible, id)
		}
		if j.Step == StepFetch {
			// the proxy has the document though it did not say so
			j.Printer.Sent, j.Step = true, StepRemote
		}
		if impressions > 0 {
			j.Impressions = impressions
		}
		if message != "" {
			j.Message = message
		}
		reasons = slices.DeleteFunc(slices.Clone(reasons), func(r string) bool { return r == "" || r == "none" })
		if j.State == ipp.JobPending {
			j.State = ipp.JobProcessing
		}
		switch st {
		case ipp.JobCompleted:
			j.Step = ""
			return j.to(now, ipp.JobCompleted, cmpOr(first(reasons), "job-completed-successfully"))
		case ipp.JobCanceled:
			j.Step = ""
			return j.to(now, ipp.JobCanceled, "job-canceled-at-device")
		case ipp.JobAborted:
			j.Step = ""
			return j.to(now, ipp.JobAborted, cmpOr(first(reasons), "aborted-by-system"))
		case ipp.JobProcessingStopped:
			return j.to(now, ipp.JobProcessingStopped, cmpOr(first(reasons), "printer-stopped"))
		}
		return j.to(now, ipp.JobProcessing, cmpOr(first(reasons), "job-printing"))
	})
}

func first(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return s[0]
}

func cmpOr(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// ProxyJobs are the jobs bound to the proxy of device that are not done
// here: those it should have.
func (s *Spool) ProxyJobs(printer, device string) []Job {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Job
	for _, e := range s.jobs {
		j := &e.job
		if j.Printer != nil && j.Printer.Name == printer && j.Printer.Device == device && (j.Step == StepRemote || j.Step == StepFetch) {
			out = append(out, j.clone())
		}
	}
	return out
}

// ProxyLost ends a job the proxy of device had the document of and no
// longer knows (it lost its spool): aborted, for a person to retry.
func (s *Spool) ProxyLost(ctx context.Context, id int32, printer, device string) (Job, error) {
	return s.proxyJob(ctx, id, printer, func(j *Job, now time.Time) error {
		if j.Printer.Device != device || j.Step != StepRemote {
			return fmt.Errorf("%w: job %d", ErrNotPossible, id)
		}
		j.Step, j.Message = "", "the proxy no longer has the job"
		return j.to(now, ipp.JobAborted, "aborted-by-system")
	})
}

// ProxySeen records what a proxy reports of the printer behind it (its
// state and capabilities): the printer of proxies is up while the proxies
// report, for proxyStale.
func (s *Spool) ProxySeen(printer, device string, caps *ipp.Printer) error {
	s.mu.Lock()
	p := s.printers[printer]
	if p == nil || !p.Proxy {
		s.mu.Unlock()
		return fmt.Errorf("%w: printer %q is not one of proxies", ErrNotFound, printer)
	}
	was := s.printerInfo(p)
	p.up, p.err, p.checked, p.seenBy = true, "", s.o.Now(), device
	if caps != nil {
		p.caps = caps
	}
	now := s.printerInfo(p)
	s.mu.Unlock()
	if was.Up != now.Up || was.State != now.State || !slices.Equal(was.Reasons, now.Reasons) || was.Accepting != now.Accepting {
		s.o.Log.Info("printer of proxies", "printer", printer, "device", device, "state", now.State)
		s.publish(Event{Printer: &now})
	}
	return nil
}

// ProxyGone records that a proxy stopped (Deregister-Output-Device).
func (s *Spool) ProxyGone(printer, device string) {
	s.mu.Lock()
	p := s.printers[printer]
	if p == nil || !p.Proxy || p.seenBy != device {
		s.mu.Unlock()
		return
	}
	p.up, p.err = false, "the proxy stopped"
	info := s.printerInfo(p)
	s.mu.Unlock()
	s.publish(Event{Printer: &info})
}

// lookProxy takes a printer of proxies as down when no proxy reported for
// proxyStale.
func (s *Spool) lookProxy(p *printerState) {
	s.mu.Lock()
	stale := p.up && s.o.Now().Sub(p.checked) > proxyStale
	if stale {
		p.up, p.err = false, "no proxy reported for "+proxyStale.String()
	}
	info := s.printerInfo(p)
	s.mu.Unlock()
	if stale {
		s.o.Log.Warn("no proxy reports", "printer", p.Name)
		s.publish(Event{Printer: &info})
	}
}

// IsProxy tells a printer of proxies.
func (s *Spool) IsProxy(printer string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.printers[printer]
	return p != nil && p.Proxy
}

// Document opens the document of a job as a .bdf (the one received, or
// the one converted), unlocked with the password the spool holds.
func (s *Spool) Document(ctx context.Context, id int32) (*bdf.Document, error) {
	j, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.document(ctx, &j)
}

// OpenBDF opens the .bdf of a job as it is kept (locked when it was):
// the one received, or the one converted. The caller closes it.
func (s *Spool) OpenBDF(ctx context.Context, id int32) (io.ReadCloser, int64, error) {
	j, err := s.Get(ctx, id)
	if err != nil {
		return nil, 0, err
	}
	key := originalKey(id)
	if j.Document.Format != "bdf" {
		if j.Document.BDF == nil {
			return nil, 0, errors.New("spool: the document is not converted")
		}
		key = bdfKey(id)
	}
	return s.st.openDoc(ctx, key)
}

// Password is the password the spool holds of a job's document, "" when
// none.
func (s *Spool) Password(id int32) string { return s.password(id) }

// FindOrigin finds the job, done or not, that came from an upstream's job
// of a uuid: a proxy does not take a job twice.
func (s *Spool) FindOrigin(ctx context.Context, upstream, uuid string) (Job, bool) {
	s.mu.Lock()
	for _, e := range s.jobs {
		if o := e.job.Origin; o != nil && o.Upstream == upstream && o.UUID == uuid {
			j := e.job.clone()
			s.mu.Unlock()
			return j, true
		}
	}
	for i := range s.recent {
		if o := s.recent[i].Origin; o != nil && o.Upstream == upstream && o.UUID == uuid {
			j := s.recent[i].clone()
			s.mu.Unlock()
			return j, true
		}
	}
	s.mu.Unlock()
	var found Job
	ok := false
	s.st.eachDone(context.WithoutCancel(ctx), func(j Job, _ string) bool {
		if o := j.Origin; o != nil && o.Upstream == upstream && o.UUID == uuid {
			found, ok = j, true
			return false
		}
		return true
	})
	return found, ok
}

// Origins are the jobs that came from an upstream, done within the
// retention of documents or not.
func (s *Spool) Origins(upstream string) []Job {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Job
	for _, e := range s.jobs {
		if o := e.job.Origin; o != nil && o.Upstream == upstream {
			out = append(out, e.job.clone())
		}
	}
	for i := range s.recent {
		if o := s.recent[i].Origin; o != nil && o.Upstream == upstream {
			out = append(out, s.recent[i].clone())
		}
	}
	return out
}
