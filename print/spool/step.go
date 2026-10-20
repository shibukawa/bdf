package spool

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/ipp"
)

// convert turns the document of a job into a .bdf and keeps it, so that
// it is not converted again when the job is sent again.
func (s *Spool) convert(ctx context.Context, e *entry) {
	if !s.slot(ctx) {
		return
	}
	defer s.unslot()
	j := s.snapshot(e)
	now := s.o.Now()
	if j.to(now, ipp.JobProcessing, "job-transforming") != nil || !s.commit(ctx, e, j) {
		return
	}
	doc, warnings, err := s.convertDocument(ctx, &j)
	now = s.o.Now()
	switch {
	case s.isCanceling(e):
		s.canceled(ctx, e, j)
		return
	case ctx.Err() != nil:
		// the spool stops: the job is converted again when it starts
		j.State, j.Reasons = ipp.JobPending, []string{"job-queued"}
		s.commit(context.WithoutCancel(ctx), e, j)
		return
	case errors.Is(err, ErrPasswordRequired), errors.Is(err, ErrWrongPassword):
		j.Message = err.Error()
		j.to(now, ipp.JobPendingHeld, "job-password-wait")
		s.forgetPassword(j.ID)
		s.commit(ctx, e, j)
		return
	case err != nil:
		j.Message = err.Error()
		j.to(now, ipp.JobAborted, "document-format-error")
		s.o.Log.Warn("conversion failed", "job", j.ID, "err", err)
		s.commit(ctx, e, j)
		return
	}
	size, sum, err := s.keepConverted(ctx, j.ID, doc)
	if err != nil {
		s.o.Log.Error("keep a converted document", "job", j.ID, "err", err)
		j.State, j.Reasons, j.Due = ipp.JobPending, []string{"job-queued"}, now.Add(retryDelay(0))
		s.commit(ctx, e, j)
		return
	}
	j.Document.BDF = &Converted{Size: size, SHA256: sum}
	j.Warnings = append(j.Warnings, warnings...)
	j.Step, j.Due, j.Message = StepSend, now, ""
	j.to(now, ipp.JobPending, "job-queued")
	s.commit(ctx, e, j)
}

func (s *Spool) convertDocument(ctx context.Context, j *Job) (*bdf.Document, []string, error) {
	r, err := s.st.getDoc(ctx, originalKey(j.ID))
	if err != nil {
		return nil, nil, err
	}
	defer r.Close()
	return s.o.Converter.Convert(ctx, r, r.Size, j.Document.Format, j.Document.Name, s.password(j.ID))
}

// keepConverted stores a converted document, replacing one a crash left.
func (s *Spool) keepConverted(ctx context.Context, id int32, doc *bdf.Document) (int64, string, error) {
	f, err := os.CreateTemp(s.o.TempDir, "bdf-spool-*.bdf")
	if err != nil {
		return 0, "", err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	h := sha256.New()
	if err := doc.WriteSingle(io.MultiWriter(f, h)); err != nil {
		return 0, "", err
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return 0, "", err
	}
	ctx = context.WithoutCancel(ctx)
	if err := s.st.removeDoc(ctx, bdfKey(id)); err != nil {
		return 0, "", err
	}
	size, err := s.st.putDoc(ctx, bdfKey(id), f)
	return size, hex.EncodeToString(h.Sum(nil)), err
}

// send sends a job to a printer: the attempt is recorded before
// Create-Job, the printer's job id as soon as it is known, and that the
// printer has the document once it has it, so that a crash at any point
// is put right by recoverJob without printing the job twice.
func (s *Spool) send(ctx context.Context, e *entry, p *printerState) {
	defer func() {
		s.mu.Lock()
		p.sending = false
		s.mu.Unlock()
	}()
	if !s.slot(ctx) {
		return
	}
	defer s.unslot()
	j := s.snapshot(e)
	s.mu.Lock()
	q := s.queues[j.Queue].Queue
	caps := p.caps
	s.mu.Unlock()
	doc, err := s.document(ctx, &j)
	now := s.o.Now()
	if err != nil {
		switch {
		case errors.Is(err, ErrWrongPassword), errors.Is(err, ErrPasswordRequired):
			j.Message = err.Error()
			j.to(now, ipp.JobPendingHeld, "job-password-wait")
			s.forgetPassword(j.ID)
		case ctx.Err() != nil:
			return
		default:
			j.Message = err.Error()
			j.to(now, ipp.JobAborted, "document-format-error")
		}
		s.commit(ctx, e, j)
		return
	}
	j.Attempts++
	j.Printer = &PrinterJob{Name: p.Name, URL: p.client.URL, AttemptStarted: now}
	if j.Started.IsZero() {
		j.Started = now
	}
	j.Message = ""
	if j.to(now, ipp.JobProcessing, "job-printing") != nil || !s.commit(ctx, e, j) {
		return
	}
	crashpoint("attempt-recorded")
	ticket := j.Ticket.Merge(q.Defaults)
	ticket.JobName = printerJobName(&j)
	format := j.Options.Format
	if format == "" || format == "auto" {
		format = q.Format
	}
	lossy := q.RasterIfLossy
	if j.Options.RasterIfLossy != nil {
		lossy = *j.Options.RasterIfLossy
	}
	// the context of the send: canceled by Cancel, and by a stop of the
	// spool only after the grace of Close
	sendCtx, stop := context.WithCancel(s.sendCtx)
	defer stop()
	go func() {
		select {
		case <-ctx.Done():
			if s.isCanceling(e) {
				stop()
			}
		case <-sendCtx.Done():
		}
	}()
	res, err := print.Print(sendCtx, p.client, doc, &print.Options{
		Ticket: ticket, Format: format, RasterIfLossy: lossy, View: j.Options.View, SkipRoles: j.Options.SkipRoles, Strict: j.Options.Strict,
		FontFS: s.o.FontFS, FontDirs: s.o.FontDirs, NoSystemFonts: s.o.NoSystemFonts, MaxImpressions: s.o.Limits.MaxImpressions,
		Printer: caps,
		OnCreated: func(js *ipp.JobStatus) error {
			crashpoint("created")
			j.Printer.JobID, j.Printer.JobURI = js.ID, js.URI
			if !s.commit(sendCtx, e, j) {
				return errors.New("spool: the printer's job could not be recorded")
			}
			crashpoint("created-recorded")
			return nil
		},
		OnSent: func(js *ipp.JobStatus) error {
			j.Printer.Sent = true
			if !s.commit(sendCtx, e, j) {
				return errors.New("spool: the sent document could not be recorded")
			}
			crashpoint("sent-recorded")
			return nil
		},
	})
	now = s.o.Now()
	if res != nil {
		j.Format, j.Media, j.Sheets, j.Bytes = res.Format, res.Media, res.Sheets, res.Bytes
		j.Unsupported = res.Unsupported
		j.Warnings = append(j.Warnings, res.Warnings...)
	}
	if j.Printer.Sent {
		s.count(func(st *Stats) { st.SentBytes[p.Name] += j.Bytes })
	}
	switch {
	case err == nil:
		j.Step, j.RetrySince, j.Due = StepWatch, time.Time{}, time.Time{}
		j.Retries = 0
		s.mirror(&j, res.Job)
		s.o.Log.Info("job sent", "job", j.ID, "printer", p.Name, "printer-job", j.Printer.JobID, "format", j.Format, "bytes", j.Bytes)
		s.commit(ctx, e, j)
	case s.isCanceling(e):
		s.canceled(ctx, e, j)
	case j.Printer.Sent:
		// the printer has the job; it is followed
		j.Step = StepWatch
		s.commit(ctx, e, j)
	case sendCtx.Err() != nil:
		// the spool stopped before the printer had the document: the job
		// is put right when the spool starts again
	case retryable(err):
		s.retry(ctx, e, j, p, err)
	default:
		j.Message = err.Error()
		reason := "aborted-by-system"
		switch {
		case errors.Is(err, print.ErrUnsupported):
			reason = ReasonUnsupported
		case errors.Is(err, print.ErrNoPages):
			reason = "document-format-error"
		}
		s.o.Log.Warn("job aborted", "job", j.ID, "printer", p.Name, "err", err)
		j.to(now, ipp.JobAborted, reason)
		s.commit(ctx, e, j)
	}
}

// retry puts back a job a printer did not take, to send it again later,
// unless it has been tried for longer than Limits.RetryFor.
func (s *Spool) retry(ctx context.Context, e *entry, j Job, p *printerState, cause error) {
	now := s.o.Now()
	if j.RetrySince.IsZero() {
		j.RetrySince = now
	}
	j.Message = cause.Error()
	if now.Sub(j.RetrySince) >= s.o.Limits.RetryFor {
		s.o.Log.Warn("job aborted after retrying", "job", j.ID, "err", cause)
		j.to(now, ipp.JobAborted, ReasonRetryExhausted)
		s.commit(ctx, e, j)
		return
	}
	j.Due = now.Add(retryDelay(j.Retries))
	j.Retries++
	s.count(func(st *Stats) { st.Retries[p.Name]++ })
	if j.Printer.JobID == 0 {
		// the printer may have made a job whose answer was lost
		s.dropAttempt(ctx, p, &j)
	}
	j.Step, j.Printer = StepSend, nil
	j.to(now, ipp.JobPending, ReasonRetryWait)
	s.o.Log.Info("job will be sent again", "job", j.ID, "printer", p.Name, "at", j.Due, "err", cause)
	s.commit(ctx, e, j)
	// the printer is checked again before the next try
	s.mu.Lock()
	if errors.Is(cause, print.ErrPrinterStopped) && p.caps != nil {
		p.caps.State = ipp.PrinterStopped
	} else if !errors.Is(cause, print.ErrPrinterStopped) {
		p.up = false
	}
	s.mu.Unlock()
	p.poke()
}

// dropAttempt cancels, as best it can, the jobs of a failed attempt that
// the printer made although its answer was lost: they wait for a document
// that does not come.
func (s *Spool) dropAttempt(ctx context.Context, p *printerState, j *Job) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer cancel()
	jobs, err := p.client.GetJobs(ctx, "not-completed", true, 0)
	if err != nil {
		return
	}
	prefix := attemptPrefix(j)
	for _, js := range jobs {
		if strings.HasPrefix(js.Name, prefix) {
			p.client.CancelJob(ctx, js.ID)
		}
	}
}

// crashHook, set by tests, is called at the points of a send where a
// crash is put right differently.
var crashHook func(point string)

func crashpoint(point string) {
	if crashHook != nil {
		crashHook(point)
	}
}

// retryable tells a failure for want of the printer: it was not reached,
// is stopped, busy or not taking jobs. Other failures are the job's.
func retryable(err error) bool {
	if errors.Is(err, print.ErrPrinterStopped) || errors.Is(err, ipp.ErrAuth) || errors.Is(err, ipp.ErrCertificateChanged) {
		return true
	}
	var se *ipp.StatusError
	if errors.As(err, &se) {
		switch se.Status {
		case ipp.StatusServerErrorServiceUnavailable, ipp.StatusServerErrorTemporaryError, ipp.StatusServerErrorNotAcceptingJobs,
			ipp.StatusServerErrorBusy, ipp.StatusServerErrorDeviceError, ipp.StatusServerErrorTooManyJobs, ipp.StatusServerErrorInternal,
			ipp.StatusClientErrorNotAuthenticated, ipp.StatusClientErrorNotAuthorized, ipp.StatusClientErrorTimeout:
			return true
		}
		return false
	}
	if errors.Is(err, ipp.ErrProtocol) {
		return true
	}
	var ne net.Error
	if errors.As(err, &ne) {
		return true
	}
	// errors of the connection that are not net.Error (a body cut, an
	// EOF of the answer)
	return errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) || strings.Contains(err.Error(), "connection")
}

// document opens the document of a job, as a .bdf: the one received, or
// the one converted.
func (s *Spool) document(ctx context.Context, j *Job) (*bdf.Document, error) {
	key := originalKey(j.ID)
	if j.Document.Format != "bdf" {
		if j.Document.BDF == nil {
			return nil, errors.New("spool: the document is not converted")
		}
		key = bdfKey(j.ID)
	}
	r, err := s.st.getDoc(ctx, key)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	br, err := bdf.OpenSingle(r, r.Size)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBadDocument, err)
	}
	if br.Locked() {
		pw := s.password(j.ID)
		if pw == "" {
			return nil, ErrPasswordRequired
		}
		if err := br.Unlock(pw); errors.Is(err, bdf.ErrWrongPassword) {
			return nil, ErrWrongPassword
		} else if err != nil {
			return nil, err
		}
	}
	return br.ToDocument()
}

// watch follows a job at its printer until it is done.
func (s *Spool) watch(ctx context.Context, e *entry) {
	j := s.snapshot(e)
	p := s.printer(j.Printer.Name)
	if p == nil {
		s.o.Log.Error("job at a printer that is not configured", "job", j.ID, "printer", j.Printer.Name)
		return
	}
	last := &ipp.JobStatus{ID: j.Printer.JobID, URI: j.Printer.JobURI, State: ipp.JobPending}
	delay := watchFirst
	crashpoint("watching")
	for {
		select {
		case <-ctx.Done():
			if s.isCanceling(e) {
				s.cancelAtPrinter(context.WithoutCancel(ctx), &j)
				s.canceled(context.WithoutCancel(ctx), e, j)
			}
			return
		case <-time.After(delay):
		}
		cctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		js, err := print.Check(cctx, p.client, last)
		cancel()
		switch {
		case js == nil && ipp.IsStatus(err, ipp.StatusClientErrorNotFound):
			// the printer keeps no history of the job: it had the document
			j.Warnings = append(j.Warnings, "the printer keeps no history of the job: it is taken as printed")
			s.done(ctx, e, j, ipp.JobCompleted, "job-completed-with-warnings")
			return
		case js == nil:
			if ctx.Err() == nil {
				delay = min(delay*2, time.Minute)
			}
			continue
		}
		last = js
		if s.mirror(&j, js) {
			if !s.commit(ctx, e, j) {
				return
			}
		}
		if js.State.Done() {
			switch js.State {
			case ipp.JobCompleted:
				reason := "job-completed-successfully"
				if len(j.Warnings) > 0 || len(j.Unsupported) > 0 {
					reason = "job-completed-with-warnings"
				}
				s.done(ctx, e, j, ipp.JobCompleted, reason)
			case ipp.JobCanceled:
				s.done(ctx, e, j, ipp.JobCanceled, "job-canceled-at-device")
			default:
				s.done(ctx, e, j, ipp.JobAborted, "aborted-by-system")
			}
			return
		}
		delay = min(delay*3/2, watchMax)
		if js.State == ipp.JobProcessingStopped {
			delay = watchMax
		}
	}
}

// mirror copies the state the printer reports of a job into it, and tells
// whether it changed.
func (s *Spool) mirror(j *Job, js *ipp.JobStatus) bool {
	if js == nil {
		return false
	}
	st := ipp.JobProcessing
	if js.State == ipp.JobProcessingStopped {
		st = ipp.JobProcessingStopped
	}
	reasons := slices.DeleteFunc(slices.Clone(js.Reasons), func(r string) bool { return r == "none" || r == "" })
	if len(reasons) == 0 {
		reasons = []string{"job-printing"}
	}
	changed := j.State != st || !slices.Equal(j.Reasons, reasons)
	if js.Impressions > 0 && int(js.Impressions) != j.Impressions {
		j.Impressions = int(js.Impressions)
		changed = true
	}
	if js.State.Done() {
		return changed
	}
	if changed {
		if err := j.to(s.o.Now(), st, reasons...); err != nil {
			return false
		}
	}
	return changed
}

// done records the end of a job at its printer.
func (s *Spool) done(ctx context.Context, e *entry, j Job, st ipp.JobState, reason string) {
	if j.State == ipp.JobPending || j.State == ipp.JobPendingHeld {
		j.State = ipp.JobProcessing
	}
	if err := j.to(s.o.Now(), st, reason); err != nil {
		s.o.Log.Error("end a job", "job", j.ID, "err", err)
		return
	}
	s.o.Log.Info("job done", "job", j.ID, "state", st, "reason", reason, "printer", j.Printer.Name)
	s.commit(ctx, e, j)
}

// canceled records a job canceled while a step ran.
func (s *Spool) canceled(ctx context.Context, e *entry, j Job) {
	ctx = context.WithoutCancel(ctx)
	if j.Printer != nil && j.Printer.JobID != 0 {
		s.cancelAtPrinter(ctx, &j)
	}
	if j.State == ipp.JobPendingHeld {
		j.State = ipp.JobPending
	}
	if err := j.to(s.o.Now(), ipp.JobCanceled, "job-canceled-by-user"); err != nil {
		s.o.Log.Error("cancel", "job", j.ID, "err", err)
		return
	}
	s.o.Log.Info("job canceled", "job", j.ID)
	s.commit(ctx, e, j)
}

func (s *Spool) forgetPassword(id int32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.passwords, id)
}

// slot takes one of the slots of conversion and drawing, or returns false
// when the step is canceled first.
func (s *Spool) slot(ctx context.Context) bool {
	select {
	case s.renders <- struct{}{}:
		return true
	case <-ctx.Done():
		return false
	}
}

func (s *Spool) unslot() { <-s.renders }
