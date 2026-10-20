package spool

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
)

// recoverJob puts right a job that was at a printer when the spool
// stopped, from what was recorded and what the printer says, so that no
// job is printed twice:
//
//   - the printer's job id recorded and the document sent: the job is
//     followed (or taken as printed when the printer forgot it);
//   - the id recorded, the document not sent: a job the printer still
//     waits a document for is canceled, and the job is sent again;
//   - no id (a crash before Create-Job's answer was recorded): the
//     printer's jobs are searched for the attempt's name; none, and the
//     job is sent again; one is the job; more are held for a person
//     (bdf-duplicate-check).
//
// A printer not reached is asked again later.
func (s *Spool) recoverJob(ctx context.Context, e *entry) {
	j := s.snapshot(e)
	if j.Printer == nil {
		s.recovered(e)
		return
	}
	p := s.printer(j.Printer.Name)
	if p == nil {
		s.o.Log.Error("job at a printer that is not configured", "job", j.ID, "printer", j.Printer.Name)
		return
	}
	actx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	now := s.o.Now()
	if j.Printer.JobID == 0 {
		jobs, err := p.client.GetJobs(actx, "all", true, 0)
		if err != nil {
			s.later(ctx, e, j, err)
			return
		}
		prefix := attemptPrefix(&j)
		var found []*ipp.JobStatus
		for _, js := range jobs {
			if strings.HasPrefix(js.Name, prefix) {
				found = append(found, js)
			}
		}
		switch len(found) {
		case 0:
			s.o.Log.Info("recovery: the printer has not the job", "job", j.ID, "printer", p.Name)
			s.resend(ctx, e, j)
			return
		case 1:
			j.Printer.JobID, j.Printer.JobURI = found[0].ID, found[0].URI
			// whether the printer had the document: a job it still waits
			// a document for has not
			j.Printer.Sent = found[0].State != ipp.JobPending || !waitsDocument(found[0])
		default:
			ids := make([]string, len(found))
			for i, js := range found {
				ids[i] = fmt.Sprint(js.ID)
			}
			// a person tells: Cancel when one of them is the job, Release
			// to send it again
			j.Message = fmt.Sprintf("the printer has %d jobs that may be this one: %s", len(found), strings.Join(ids, ", "))
			j.Step = StepSend
			j.to(now, ipp.JobPendingHeld, ReasonDuplicateCheck)
			s.o.Log.Warn("recovery: a person decides", "job", j.ID, "printer", p.Name, "printer-jobs", ids)
			s.recovered(e)
			s.commit(ctx, e, j)
			return
		}
	}
	js, err := p.client.GetJobAttributes(actx, j.Printer.JobID)
	switch {
	case ipp.IsStatus(err, ipp.StatusClientErrorNotFound):
		if j.Printer.Sent {
			j.Warnings = append(j.Warnings, "the printer keeps no history of the job: it is taken as printed")
			s.recovered(e)
			s.done(ctx, e, j, ipp.JobCompleted, "job-completed-with-warnings")
			return
		}
		s.resend(ctx, e, j)
		return
	case err != nil:
		s.later(ctx, e, j, err)
		return
	}
	if !j.Printer.Sent {
		if !js.State.Done() {
			if waitsDocument(js) || js.State == ipp.JobPending {
				p.client.CancelJob(actx, js.ID)
			} else {
				// the printer prints it: it had the document
				j.Printer.Sent = true
			}
		}
		if !j.Printer.Sent {
			s.resend(ctx, e, j)
			return
		}
	}
	s.o.Log.Info("recovery: the job is followed at the printer", "job", j.ID, "printer", p.Name, "printer-job", js.ID)
	j.Step = StepWatch
	s.mirror(&j, js)
	s.recovered(e)
	s.commit(ctx, e, j)
}

// waitsDocument tells a job the printer made and has no document of.
func waitsDocument(js *ipp.JobStatus) bool {
	for _, r := range js.Reasons {
		if r == "job-incoming" || r == "job-data-insufficient" {
			return true
		}
	}
	return false
}

// resend queues a job to be sent again, as a new attempt.
func (s *Spool) resend(ctx context.Context, e *entry, j Job) {
	now := s.o.Now()
	j.Step, j.Printer, j.Due = StepSend, nil, now
	if j.State == ipp.JobPendingHeld {
		j.State = ipp.JobPending
	}
	j.to(now, ipp.JobPending, "job-queued")
	s.recovered(e)
	s.commit(ctx, e, j)
}

// later asks the printer again after a while.
func (s *Spool) later(ctx context.Context, e *entry, j Job, err error) {
	s.o.Log.Info("recovery: the printer is not reached", "job", j.ID, "printer", j.Printer.Name, "err", err)
	s.mu.Lock()
	e.job.Due = s.o.Now().Add(retryDelay(e.job.Retries))
	e.job.Retries++
	s.mu.Unlock()
}

func (s *Spool) recovered(e *entry) {
	s.mu.Lock()
	e.recover = false
	s.mu.Unlock()
}
