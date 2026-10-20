package spool

import (
	"context"
	"slices"
	"sort"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
)

// Start runs the spool: a watcher per printer, and the loop that converts
// documents, sends jobs to their printers and follows them. Jobs that were
// at a printer when the spool last stopped are checked against the
// printer first. Close stops it.
func (s *Spool) Start() {
	s.mu.Lock()
	if s.started || s.closed {
		s.mu.Unlock()
		return
	}
	s.started = true
	for _, e := range s.jobs {
		switch {
		case e.job.Step == StepWatch:
			e.recover = true
		case e.job.Step == StepSend && e.job.Printer != nil:
			e.recover = true
		case e.job.State == ipp.JobProcessing && e.job.Step == StepConvert:
			// stopped while it was converted: it is converted again
			e.job.State, e.job.Reasons = ipp.JobPending, []string{"job-queued"}
		}
	}
	printers := make([]*printerState, 0, len(s.printers))
	for _, p := range s.printers {
		printers = append(printers, p)
	}
	s.mu.Unlock()
	for _, p := range printers {
		s.running.Add(1)
		go s.monitor(p)
	}
	s.running.Add(1)
	go s.loop()
}

// loop starts the steps whose time has come, until the spool closes.
func (s *Spool) loop() {
	defer s.running.Done()
	timer := time.NewTimer(0)
	defer timer.Stop()
	collect := s.o.Now().Add(collectEvery)
	for {
		select {
		case <-s.stop:
			return
		case <-s.kick:
		case <-timer.C:
		}
		now := s.o.Now()
		if !now.Before(collect) {
			collect = now.Add(collectEvery)
			s.running.Add(1)
			go func() {
				defer s.running.Done()
				if err := s.Collect(s.stepCtx); err != nil && s.stepCtx.Err() == nil {
					s.o.Log.Error("collect", "err", err)
				}
			}()
		}
		next := s.schedule(now)
		if next.IsZero() || next.After(collect) {
			next = collect
		}
		timer.Reset(max(next.Sub(s.o.Now()), 10*time.Millisecond))
	}
}

// schedule starts what is due and returns when it should look again (the
// zero time when nothing waits for a time).
func (s *Spool) schedule(now time.Time) time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return time.Time{}
	}
	var next time.Time
	soon := func(t time.Time) {
		if !t.IsZero() && (next.IsZero() || t.Before(next)) {
			next = t
		}
	}
	var sendable []*entry
	for _, e := range s.jobs {
		j := &e.job
		if e.busy {
			continue
		}
		switch {
		case e.recover:
			if j.Due.After(now) {
				soon(j.Due)
				continue
			}
			s.begin(e, s.recoverJob)
		case j.State == ipp.JobPendingHeld:
			// a hold that ends by itself (job-hold-until with a time)
			if j.Due.IsZero() || slices.Contains(j.Reasons, "job-password-wait") || slices.Contains(j.Reasons, ReasonDuplicateCheck) {
				continue
			}
			if j.Due.After(now) {
				soon(j.Due)
				continue
			}
			s.begin(e, s.releaseHold)
		case j.Step == StepWatch:
			s.begin(e, s.watch)
		case j.Step == StepFetch || j.Step == StepRemote:
			// a proxy fetches it, and reports it
		case j.State != ipp.JobPending:
		case j.Due.After(now):
			soon(j.Due)
		case s.queues[j.Queue] == nil || s.queues[j.Queue].rec.Paused:
		case j.Step == StepConvert:
			s.begin(e, s.convert)
		case j.Step == StepSend:
			sendable = append(sendable, e)
		}
	}
	// the jobs to send, in their queues' order: the highest priority,
	// then the oldest
	sort.Slice(sendable, func(a, b int) bool {
		x, y := &sendable[a].job, &sendable[b].job
		if x.Priority != y.Priority {
			return x.Priority > y.Priority
		}
		if !x.Submitted.Equal(y.Submitted) {
			return x.Submitted.Before(y.Submitted)
		}
		return x.ID < y.ID
	})
	for _, e := range sendable {
		if p := s.printers[s.queues[e.job.Queue].Printers[0]]; p.Proxy {
			s.begin(e, func(ctx context.Context, e *entry) { s.offer(ctx, e, p) })
			continue
		}
		if p := s.choose(s.queues[e.job.Queue]); p != nil {
			p.sending = true
			target := p
			s.begin(e, func(ctx context.Context, e *entry) { s.send(ctx, e, target) })
		}
	}
	return next
}

// choose picks the printer of a queue a job is sent to now, or nil when
// none can take one: the first printer that is up and has room ("any"),
// or the first printer that is up, when it has room ("failover").
func (s *Spool) choose(q *queueState) *printerState {
	for _, name := range q.Printers {
		p := s.printers[name]
		if !p.available() {
			continue
		}
		if !p.sending && s.inflight(p) < q.MaxInflight {
			return p
		}
		if q.Policy == "failover" {
			return nil
		}
	}
	return nil
}

// available tells a printer that answered last, is not stopped, and takes
// jobs.
func (p *printerState) available() bool {
	if !p.up || p.caps == nil || p.caps.State == ipp.PrinterStopped {
		return false
	}
	if _, said := p.caps.Has("printer-is-accepting-jobs"); said && !p.caps.AcceptingJobs {
		return false
	}
	return true
}

// inflight counts the jobs a printer has and has not finished.
func (s *Spool) inflight(p *printerState) int {
	n := 0
	for _, e := range s.jobs {
		if e.job.Printer != nil && e.job.Printer.Name == p.Name && e.job.Step == StepWatch {
			n++
		}
	}
	return n
}

// begin runs a step of a job in its own goroutine. s.mu is held.
func (s *Spool) begin(e *entry, step func(context.Context, *entry)) {
	ctx, cancel := context.WithCancel(s.stepCtx)
	e.busy, e.cancel = true, cancel
	s.running.Add(1)
	go func() {
		defer s.running.Done()
		defer cancel()
		step(ctx, e)
		s.mu.Lock()
		e.busy, e.cancel = false, nil
		s.mu.Unlock()
		s.wake()
	}()
}

// snapshot returns a copy of a job to change in a step.
func (s *Spool) snapshot(e *entry) Job {
	s.mu.Lock()
	defer s.mu.Unlock()
	return e.job.clone()
}

// canceling tells whether a job was asked to be canceled while a step ran.
func (s *Spool) isCanceling(e *entry) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return e.canceling
}

// releaseHold queues a job whose hold has ended.
func (s *Spool) releaseHold(ctx context.Context, e *entry) {
	j := s.snapshot(e)
	now := s.o.Now()
	j.HoldUntil, j.HoldUntilTime, j.Due = "", time.Time{}, now
	if err := j.to(now, ipp.JobPending, "job-queued"); err != nil {
		s.o.Log.Error("release", "job", j.ID, "err", err)
		return
	}
	s.commit(ctx, e, j)
}

// commit saves a job a step changed, logging what fails.
func (s *Spool) commit(ctx context.Context, e *entry, j Job) bool {
	var err error
	if j.State.Done() {
		err = s.finish(ctx, e, j)
	} else {
		err = s.save(ctx, e, j)
	}
	if err != nil {
		s.o.Log.Error("record a job", "job", j.ID, "err", err)
		return false
	}
	return true
}

// retryDelays are the waits before a job is sent again to a printer that
// was not reached; the last repeats.
var retryDelays = []time.Duration{10 * time.Second, 30 * time.Second, time.Minute, 2 * time.Minute, 5 * time.Minute, 10 * time.Minute}

func retryDelay(n int) time.Duration {
	return retryDelays[min(max(n, 0), len(retryDelays)-1)]
}
