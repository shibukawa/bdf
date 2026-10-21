package manager

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/raw"
)

// ErrBranch is the error of a printer that a branch's bdf printd is to
// fetch the jobs of: not done yet.
var ErrBranch = errors.New("manager: printers of a branch (route branch) are not supported yet")

// ErrPrint is the error of a printer that did not take a slip.
var ErrPrint = errors.New("manager: the printer did not print")

// ErrLeased is the error of a slip another manager is sending.
var ErrLeased = errors.New("manager: the slip is being printed")

// leaseFor bounds a send: a lease older than this is taken over (the
// manager that held it stopped).
const leaseFor = 5 * time.Minute

// checkPrinterURL takes the printers a manager reaches over the network:
// IPP, and 9100 and LPD. Serial, USB and Windows queues are a branch
// printd's.
func checkPrinterURL(u string) error {
	if raw.IsURL(u) {
		p, err := raw.Parse(u)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInvalid, err)
		}
		if p.Scheme != "socket" && p.Scheme != "lpd" {
			return fmt.Errorf("%w: %s: a printer of a manager is reached over the network (ipp, ipps, socket, lpd)", ErrInvalid, u)
		}
		return nil
	}
	if _, _, err := ipp.ParseURL(u); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	return nil
}

// PrintSlip prints a slip on a printer now, as a device asked.
func (m *Manager) PrintSlip(ctx context.Context, id, printer string, by Mark) (*Slip, error) {
	return m.printSlip(ctx, id, printer, whyDevice, by.Device)
}

// printSlip sends a slip to a printer under a lease: the lease is written
// first, so that two managers do not both send it; a failed send is
// tried again later (RetryAt).
func (m *Manager) printSlip(ctx context.Context, id, printer, why, device string) (*Slip, error) {
	p, err := m.printer(ctx, printer)
	if err != nil {
		return nil, err
	}
	if p.Route == RouteBranch {
		return nil, ErrBranch
	}
	owner := randomHex(8)
	now := m.now()
	s, err := m.changeSlip(ctx, id, func(s *Slip) error {
		if s.Lease != nil && now.Before(s.Lease.Until) {
			return ErrLeased
		}
		if s.State == SlipCanceled {
			return fmt.Errorf("%w: the slip is canceled", ErrConflict)
		}
		if why == whyFallback && (s.State != SlipNew || s.fallbackPrinted()) {
			return errSkip
		}
		s.Lease = &Lease{Owner: owner, Until: now.Add(leaseFor), Printer: printer, Why: why, Device: device}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if s.Lease == nil || s.Lease.Owner != owner {
		return s, nil // nothing to print (errSkip)
	}
	sctx, cancel := context.WithTimeout(ctx, leaseFor-30*time.Second)
	job, sendErr := m.sendSlip(sctx, p, s)
	cancel()
	done := m.now()
	s, err = m.changeSlip(context.WithoutCancel(ctx), id, func(s *Slip) error {
		if s.Lease != nil && s.Lease.Owner == owner {
			s.Lease = nil
		}
		if sendErr != nil {
			s.Attempts++
			s.LastError = clip(sendErr.Error(), 500)
			s.RetryAt = done.Add(retryDelay(s.Attempts))
			return nil
		}
		s.Printed = append(s.Printed, Printed{Printer: printer, Why: why, Device: device, Job: job, At: done})
		s.LastError, s.RetryAt = "", time.Time{}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if sendErr != nil {
		m.o.Log.Warn("print", "slip", id, "printer", printer, "why", why, "err", sendErr)
	} else {
		m.o.Log.Info("print", "slip", id, "printer", printer, "why", why, "job", job)
	}
	if err := m.touch(context.WithoutCancel(ctx), s.Location, func(*locationState) {}); err != nil {
		return s, err
	}
	if sendErr != nil {
		return s, fmt.Errorf("%w on %s: %v", ErrPrint, printer, sendErr)
	}
	return s, nil
}

// retryDelay is the wait before the next try of a print that failed.
func retryDelay(attempts int) time.Duration {
	d := 30 * time.Second << min(attempts-1, 5)
	return min(d, 15*time.Minute)
}

func (m *Manager) sendSlip(ctx context.Context, p *Node, s *Slip) (int32, error) {
	doc, err := m.Document(ctx, s.Job)
	if err != nil {
		return 0, err
	}
	var meta *print.Meta
	if j, err := m.JobOfNumber(ctx, s.Job); err == nil {
		meta = j.Meta
	}
	return m.send(ctx, p, doc, firstOr(s.Name, "slip-"+s.ID), meta)
}

// send prints a document on a printer, with what the job says about
// itself (meta, which may be nil), and returns the printer's job id (0
// for a raw printer, which has none).
func (m *Manager) send(ctx context.Context, p *Node, doc *bdf.Document, name string, meta *print.Meta) (int32, error) {
	opts := &print.Options{Ticket: print.Ticket{JobName: name}, FontDirs: m.o.FontDirs, NoSystemFonts: m.o.NoSystemFonts}
	if meta != nil {
		opts.Ticket.Meta = *meta
	}
	if raw.IsURL(p.URL) {
		rp, err := raw.Parse(p.URL)
		if err != nil {
			return 0, err
		}
		_, err = print.PrintRaw(ctx, rp, doc, opts)
		return 0, err
	}
	res, err := print.Print(ctx, &ipp.Client{URL: p.URL, User: "bdf-print-manager"}, doc, opts)
	if err != nil {
		return 0, err
	}
	if res.Job != nil {
		return res.Job.ID, nil
	}
	return 0, nil
}

func firstOr(s, or string) string {
	if strings.TrimSpace(s) == "" {
		return or
	}
	return s
}

// TickReport is what a tick did.
type TickReport struct {
	// Delivered and Undelivered count the deliveries to printers and
	// archives sent, and the sends that failed (tried again later).
	Delivered   int `json:"delivered"`
	Undelivered int `json:"undelivered"`
	// Printed and Failed count the slips printed on paper as a fallback,
	// and those that could not be.
	Printed int `json:"printed"`
	Failed  int `json:"failed"`
	// Removed, Items and Jobs count the slips, the items and the jobs
	// removed after their keep.
	Removed  int `json:"removed"`
	Items    int `json:"items"`
	Jobs     int `json:"jobs"`
	Expired  int `json:"expired"`
	Leased   int `json:"leased"`
	Location int `json:"locations"`
}

// Tick does what is due: it sends the deliveries to printers and
// archives that wait (or failed and are to be tried again), prints on
// paper the slips nobody opened in time (or every unopened slip of a
// display none of whose devices is connected), removes the slips and
// the items ended longer ago than their display or inbox keeps them, the
// jobs that ended and nothing shows any more, and the enrollment codes
// expired. Any instance may call it, at any time, as often as it likes: a
// scheduler of the platform (every minute), or a timer of the process.
func (m *Manager) Tick(ctx context.Context) (*TickReport, error) {
	rep := &TickReport{}
	var errs []error
	if err := m.tickDeliveries(ctx, rep); err != nil {
		errs = append(errs, err)
	}
	nodes, err := m.Nodes(ctx)
	if err != nil {
		return nil, err
	}
	for _, n := range nodes {
		switch n.Kind {
		case KindDisplay:
			rep.Location++
			if err := m.tickLocation(ctx, &n, rep); err != nil {
				errs = append(errs, fmt.Errorf("%s: %w", n.Name, err))
			}
		case KindInbox:
			if err := m.tickInbox(ctx, &n, rep); err != nil {
				errs = append(errs, fmt.Errorf("%s: %w", n.Name, err))
			}
		}
	}
	codes, err := list[enrollRecord](ctx, m.st.db)
	if err != nil {
		errs = append(errs, err)
	}
	now := m.now()
	for _, c := range codes {
		if now.After(c.Data.Expires) {
			if err := m.st.db.Delete[enrollRecord](ctx, c.ID(), c.Revision()); err == nil {
				rep.Expired++
			}
		}
	}
	if err := m.tickJobs(ctx, rep); err != nil {
		errs = append(errs, err)
	}
	return rep, errors.Join(errs...)
}

func (m *Manager) tickLocation(ctx context.Context, l *Node, rep *TickReport) error {
	rec, err := m.st.db.Get[locationState](ctx, l.Name)
	if err != nil {
		if errors.Is(mapErr(err), ErrNotFound) {
			return nil
		}
		return mapErr(err)
	}
	st := rec.Data
	now := m.now()
	var errs []error
	if l.Fallback != nil {
		offline := !m.anyOnline(ctx, st.Devices, now.Add(-l.offlineAfter()))
		for _, id := range st.Open {
			s, err := m.Slip(ctx, id)
			if err != nil {
				continue
			}
			due := offline || (!s.FallbackAt.IsZero() && !now.Before(s.FallbackAt))
			if s.State != SlipNew || s.fallbackPrinted() || !due || now.Before(s.RetryAt) {
				continue
			}
			_, err = m.printSlip(ctx, id, l.Fallback.Printer, whyFallback, "")
			switch {
			case errors.Is(err, ErrLeased):
				rep.Leased++
			case err != nil:
				rep.Failed++
				errs = append(errs, err)
			default:
				rep.Printed++
			}
		}
	}
	for _, id := range st.Ended {
		s, err := m.Slip(ctx, id)
		if errors.Is(err, ErrNotFound) {
			s = &Slip{ID: id, Location: l.Name}
		} else if err != nil {
			continue
		} else if at := endedAt(s); at.IsZero() || now.Sub(at) < l.keep() {
			continue
		}
		if err := m.removeSlip(ctx, s); err != nil {
			errs = append(errs, err)
			continue
		}
		rep.Removed++
	}
	return errors.Join(errs...)
}

func endedAt(s *Slip) time.Time {
	switch {
	case s.Done != nil:
		return s.Done.At
	case s.Canceled != nil:
		return s.Canceled.At
	}
	return time.Time{}
}

// anyOnline tells whether a device was seen since a time.
func (m *Manager) anyOnline(ctx context.Context, ids []string, since time.Time) bool {
	for _, id := range ids {
		d, err := m.Device(ctx, id)
		if err == nil && d.LastSeen.After(since) {
			return true
		}
	}
	return false
}
