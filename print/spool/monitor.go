package spool

import (
	"context"
	"slices"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
)

// statusAttributes are what the watcher asks a printer at each look; all
// the capabilities are asked again when printer-config-change-time moves.
var statusAttributes = []string{"printer-state", "printer-state-reasons", "printer-state-message", "printer-is-accepting-jobs",
	"printer-config-change-time"}

// monitor watches a printer: every minute when nothing waits for it, every
// ten seconds when jobs wait or it is stopped, and less often while it is
// not reached (up to five minutes).
func (s *Spool) monitor(p *printerState) {
	defer s.running.Done()
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		select {
		case <-s.stop:
			return
		case <-timer.C:
		case <-p.pokes:
		}
		s.look(p)
		timer.Reset(s.nextLook(p))
	}
}

// look asks a printer its state, and its capabilities when they changed
// or are not known.
func (s *Spool) look(p *printerState) {
	ctx, cancel := context.WithTimeout(s.stepCtx, 30*time.Second)
	defer cancel()
	s.mu.Lock()
	caps := p.caps
	s.mu.Unlock()
	var full *ipp.Printer
	st, err := p.client.GetPrinterAttributes(ctx, statusAttributes...)
	if err == nil && (caps == nil || st.ConfigChangeTime != caps.ConfigChangeTime) {
		full, err = p.client.GetPrinterAttributes(ctx)
	}
	if s.stepCtx.Err() != nil {
		return
	}
	s.mu.Lock()
	was := s.printerInfo(p)
	p.checked = s.o.Now()
	if err != nil {
		p.up, p.err = false, err.Error()
		p.failing++
	} else {
		p.up, p.err, p.failing = true, "", 0
		if full != nil {
			p.caps = full
		} else {
			c := *p.caps
			c.State, c.Reasons, c.StateMessage, c.AcceptingJobs = st.State, st.Reasons, st.StateMessage, st.AcceptingJobs
			p.caps = &c
		}
	}
	now := s.printerInfo(p)
	s.mu.Unlock()
	if was.Up != now.Up || was.State != now.State || !slices.Equal(was.Reasons, now.Reasons) || was.Accepting != now.Accepting {
		if err != nil {
			s.o.Log.Warn("printer not reached", "printer", p.Name, "err", err)
		} else {
			s.o.Log.Info("printer", "printer", p.Name, "state", now.State, "reasons", now.Reasons, "accepting", now.Accepting)
		}
		s.publish(Event{Printer: &now})
		s.wake()
	}
}

// nextLook is when to look at a printer again.
func (s *Spool) nextLook(p *printerState) time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !p.up {
		return min(lookBusy<<min(p.failing-1, 5), 5*time.Minute)
	}
	if p.caps != nil && p.caps.State == ipp.PrinterStopped {
		return lookBusy
	}
	for _, e := range s.jobs {
		j := &e.job
		if j.Printer != nil && j.Printer.Name == p.Name {
			return lookBusy
		}
		if j.State == ipp.JobPending {
			if q := s.queues[j.Queue]; q != nil && slices.Contains(q.Printers, p.Name) {
				return lookBusy
			}
		}
	}
	return lookIdle
}

// How often printers are looked at and jobs at printers asked about (tests
// shorten them).
var (
	lookIdle   = time.Minute
	lookBusy   = 10 * time.Second
	watchFirst = 2 * time.Second
	watchMax   = 10 * time.Second
)

// poke makes the watcher look at the printer now.
func (p *printerState) poke() {
	select {
	case p.pokes <- struct{}{}:
	default:
	}
}
