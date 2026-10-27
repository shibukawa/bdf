package docflow

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

// ErrPrint is the error of a printer that did not take a slip.
var ErrPrint = errors.New("manager: the printer did not print")

// ErrLeased is the error of a slip another manager is sending.
var ErrLeased = errors.New("manager: the slip is being printed")

// leaseFor bounds a send: a lease older than this is taken over (the
// manager that held it stopped).
const leaseFor = 5 * time.Minute

// PrintSlip prints a slip on a printer now, as a device asked.
func (m *Manager) PrintSlip(ctx context.Context, id, printer string, by Mark) (*Slip, error) {
	return m.printSlip(ctx, id, printer, whyDevice, by.Device)
}

// printSlip sends a slip to a printer under a lease: the lease is written
// first, so that two managers do not both send it; a failed send is
// tried again later (RetryAt).
func (m *Manager) printSlip(ctx context.Context, id, printer, why, device string) (*Slip, error) {
	p, k, err := m.exit(ctx, printer)
	if err != nil {
		return nil, err
	}
	if k.Ready != nil {
		if err := k.Ready(p); err != nil {
			return nil, err
		}
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
	job, sendErr := m.sendSlip(sctx, k, p, s)
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

// sendSlip sends a slip's document out of an exit (a printer), with its
// job's metadata, and returns the printer's job id (0 for none).
func (m *Manager) sendSlip(ctx context.Context, k *Kind, p *Node, s *Slip) (int32, error) {
	o := &Out{ID: "slip." + s.ID, Node: p, Payload: &Payload{ID: firstOr(s.Payload, jobID(s.Job)), Kind: PayloadBDF, Format: "bdf"}, Name: firstOr(s.Name, "slip-"+s.ID), m: m}
	if j, err := m.JobOfNumber(ctx, s.Job); err == nil {
		o.Job, o.Meta = j, j.Meta
	}
	res, err := k.Send(ctx, o)
	if err != nil {
		return 0, err
	}
	id, _ := strconv.ParseInt(res, 10, 32)
	return int32(id), nil
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
	// Removed, Items and Jobs count the slips expired, and the items and
	// the jobs removed after their keep.
	Removed  int `json:"removed"`
	Items    int `json:"items"`
	Jobs     int `json:"jobs"`
	Expired  int `json:"expired"`
	Leased   int `json:"leased"`
	Location int `json:"locations"`
	// Ticked counts the nodes of the kinds with a tick of their own
	// (Kind.Tick) ticked; Polled the files their entries took in as jobs.
	Ticked int `json:"ticked"`
	Polled int `json:"polled"`
	// Notes are what a tick has to say besides the counts (a file an
	// entry left on its far side, and why): maxNotes at most.
	Notes []string `json:"notes,omitempty"`
}

// maxNotes bounds the notes of a tick.
const maxNotes = 50

// Note adds a note to the report (fmt.Sprintf's form), up to maxNotes;
// the first ones stay.
func (r *TickReport) Note(format string, args ...any) {
	if len(r.Notes) < maxNotes {
		r.Notes = append(r.Notes, fmt.Sprintf(format, args...))
	}
}

// Tick does what is due: it sends the deliveries to printers and
// archives that wait (or failed and are to be tried again), prints on
// paper the slips nobody opened in time (or every unopened slip of a
// display none of whose devices is connected), removes the slips
// expired, the items ended longer ago than their inbox keeps them, the
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
		default:
			if k := kindNamed(n.Kind); k != nil && k.Tick != nil {
				rep.Ticked++
				if err := k.Tick(ctx, m, &n, rep); err != nil {
					errs = append(errs, fmt.Errorf("%s: %w", n.Name, err))
				}
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
	// a slip stays until it expires, opened, done or not; one being
	// printed waits for its print
	for _, id := range slices.Concat(st.Open, st.Ended) {
		s, err := m.Slip(ctx, id)
		if errors.Is(err, ErrNotFound) {
			s = &Slip{ID: id, Location: l.Name}
		} else if err != nil {
			continue
		} else if now.Before(s.expiresAt(l.keep())) || (s.Lease != nil && now.Before(s.Lease.Until)) {
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
