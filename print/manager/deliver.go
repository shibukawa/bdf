package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/shibukawa/yuruidb"
)

// A delivery takes a job to a printer or an archive. It is made when the
// job is routed and done later: by the next Tick, or at once in the
// background of a process (Options.Background), under a lease so that two
// instances do not both send it. A send that fails is tried again after
// a wait that grows, on the next member of a class, and given up after
// maxAttempts (the job is then aborted); an administrator may retry it.

// maxAttempts is how many sends are tried before a delivery fails for
// good: with retryDelay, about four hours.
const maxAttempts = 20

// errNotDue is the error of a delivery whose RetryAt has not come.
var errNotDue = errors.New("manager: not due")

// routeDelivery makes a delivery for a leg: to printer (a member, when
// the leg is a class's), or to the archive the leg names.
func (m *Manager) routeDelivery(ctx context.Context, r *routing, leg Leg, printer, class string) {
	j := r.job
	d := Delivery{ID: fmt.Sprintf("%s.%02d", j.ID, r.next()), Job: j.Number, Name: j.Name, Node: leg.Node, Kind: leg.Kind, Printer: printer, Class: class,
		State: DeliveryPending, Created: r.now}
	if d.Kind == KindClass {
		d.Kind = KindPrinter
	}
	if _, err := m.st.db.Insert(ctx, d.ID, deliveryRecord(d)); err != nil {
		leg.State, leg.Error = LegError, mapErr(err).Error()
		r.add(leg)
		return
	}
	if _, err := m.st.db.Insert(ctx, d.ID, pendingRecord{Job: j.Number}); err != nil && !errors.Is(err, yuruidb.ErrAlreadyExists) {
		leg.State, leg.Error = LegError, mapErr(err).Error()
		r.add(leg)
		return
	}
	r.live++
	r.deliveries = append(r.deliveries, d.ID)
	leg.State, leg.Ref = LegPending, d.ID
	r.add(leg)
}

// kick sends deliveries in the background of the process, when the
// manager does that.
func (m *Manager) kick(ids []string) {
	if !m.o.Background || len(ids) == 0 {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), leaseFor)
		defer cancel()
		for _, id := range ids {
			if _, err := m.deliver(ctx, id); err != nil && !errors.Is(err, ErrLeased) && !errors.Is(err, errNotDue) {
				m.o.Log.Warn("deliver", "delivery", id, "err", err)
			}
		}
	}()
}

// Delivery returns a delivery.
func (m *Manager) Delivery(ctx context.Context, id string) (*Delivery, error) {
	rec, err := m.st.db.Get[deliveryRecord](ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	d := Delivery(rec.Data)
	return &d, nil
}

// Pending returns the deliveries not done, oldest first.
func (m *Manager) Pending(ctx context.Context) ([]Delivery, error) {
	recs, err := list[pendingRecord](ctx, m.st.db)
	if err != nil {
		return nil, err
	}
	out := make([]Delivery, 0, len(recs))
	for _, r := range recs {
		d, err := m.Delivery(ctx, r.ID())
		if errors.Is(err, ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}
		out = append(out, *d)
	}
	return out, nil
}

func (m *Manager) changeDelivery(ctx context.Context, id string, fn func(*Delivery) error) (*Delivery, error) {
	r, err := change(ctx, m.st.db, id, false, func(r *deliveryRecord) error {
		return fn((*Delivery)(r))
	})
	if err != nil {
		return nil, err
	}
	d := Delivery(r)
	return &d, nil
}

// unpend takes a delivery off the pending ones.
func (m *Manager) unpend(ctx context.Context, id string) {
	rec, err := m.st.db.Get[pendingRecord](ctx, id)
	if err != nil {
		return
	}
	m.st.db.Delete[pendingRecord](ctx, id, rec.Revision())
}

// Retry sends a delivery again now: one waiting for its time, or one
// that failed for good.
func (m *Manager) Retry(ctx context.Context, id string) (*Delivery, error) {
	revived := false
	d, err := m.changeDelivery(ctx, id, func(d *Delivery) error {
		switch d.State {
		case DeliveryPending:
			d.RetryAt = time.Time{}
		case DeliveryFailed:
			d.State, d.RetryAt, d.Attempts, d.Tried, d.Done = DeliveryPending, time.Time{}, 0, nil, time.Time{}
			revived = true
		default:
			return fmt.Errorf("%w: the delivery is %s", ErrConflict, d.State)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if revived {
		if _, err := m.st.db.Insert(ctx, id, pendingRecord{Job: d.Job}); err != nil && !errors.Is(err, yuruidb.ErrAlreadyExists) {
			return nil, mapErr(err)
		}
		if _, err := m.changeJob(ctx, jobID(d.Job), func(j *Job) error {
			if l := j.leg(id); l != nil {
				l.State, l.Error = LegPending, ""
			}
			j.Live++
			return nil
		}); err != nil {
			return nil, err
		}
	}
	m.kick([]string{id})
	return d, nil
}

// deliver sends a delivery that is due, under a lease, and records the
// outcome. It returns the delivery as it is after.
func (m *Manager) deliver(ctx context.Context, id string) (*Delivery, error) {
	owner := randomHex(8)
	now := m.now()
	d, err := m.changeDelivery(ctx, id, func(d *Delivery) error {
		if d.State != DeliveryPending {
			return errSkip
		}
		if d.Lease != nil && now.Before(d.Lease.Until) {
			return ErrLeased
		}
		if now.Before(d.RetryAt) {
			return errNotDue
		}
		d.Lease = &Lease{Owner: owner, Until: now.Add(leaseFor), Why: "delivery"}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if d.Lease == nil || d.Lease.Owner != owner {
		m.unpend(ctx, id) // not pending: left over
		return d, nil
	}
	j, err := m.JobOfNumber(ctx, d.Job)
	if err != nil {
		return m.outcome(ctx, d, owner, "", "", err, true)
	}
	sctx, cancel := context.WithTimeout(ctx, leaseFor-30*time.Second)
	defer cancel()
	switch d.Kind {
	case KindArchive:
		n, err := m.Node(sctx, d.Node)
		if err == nil && n.Kind != KindArchive {
			err = fmt.Errorf("%s is a %s, not an archive", d.Node, n.Kind)
		}
		if err != nil {
			return m.outcome(ctx, d, owner, "", "", err, false)
		}
		key, err := m.archive(sctx, n, j)
		return m.outcome(ctx, d, owner, "", key, err, false)
	default:
		printer, final, err := m.member(sctx, d)
		if err != nil {
			return m.outcome(ctx, d, owner, printer, "", err, final)
		}
		p, err := m.Printer(sctx, printer)
		if err != nil {
			return m.outcome(ctx, d, owner, printer, "", err, false)
		}
		if p.Route == RouteBranch {
			return m.outcome(ctx, d, owner, printer, "", ErrBranch, true)
		}
		doc, err := m.Document(sctx, j.Number)
		if err != nil {
			return m.outcome(ctx, d, owner, printer, "", err, false)
		}
		jobID, err := m.send(sctx, p, doc, firstOr(d.Name, "job-"+j.ID))
		result := ""
		if err == nil && jobID != 0 {
			result = fmt.Sprint(jobID)
		}
		return m.outcome(ctx, d, owner, printer, result, err, false)
	}
}

// member chooses the printer a delivery goes to now: the one it was
// routed to, or, after that one failed, the next member of its class not
// tried yet (all of them tried, they are tried over). final tells an
// error that no retry mends.
func (m *Manager) member(ctx context.Context, d *Delivery) (printer string, final bool, err error) {
	if d.Class == "" || d.Attempts == 0 {
		return d.Printer, false, nil
	}
	n, err := m.Node(ctx, d.Class)
	if err != nil || n.Kind != KindClass {
		return d.Printer, false, nil // the class is gone: the printer it chose
	}
	var members []string
	for _, name := range n.Members {
		if k, err := m.kindOf(ctx, name); err == nil && k == KindPrinter {
			members = append(members, name)
		}
	}
	if len(members) == 0 {
		return d.Printer, true, fmt.Errorf("class %s has no printer", d.Class)
	}
	untried := slices.DeleteFunc(slices.Clone(members), func(name string) bool { return slices.Contains(d.Tried, name) })
	if len(untried) == 0 {
		untried = members
	}
	// the next after the last one, in the class's order
	last := slices.Index(members, d.Printer)
	for i := range members {
		name := members[(last+1+i)%len(members)]
		if slices.Contains(untried, name) {
			return name, false, nil
		}
	}
	return untried[0], false, nil
}

// outcome records what a send did and releases the lease; a delivery
// canceled meanwhile stays canceled.
func (m *Manager) outcome(ctx context.Context, d *Delivery, owner, printer, result string, sendErr error, final bool) (*Delivery, error) {
	ctx = context.WithoutCancel(ctx)
	done := m.now()
	ended := false
	d, err := m.changeDelivery(ctx, d.ID, func(d *Delivery) error {
		if d.Lease != nil && d.Lease.Owner == owner {
			d.Lease = nil
		}
		if printer != "" {
			d.Printer = printer
		}
		if d.State != DeliveryPending {
			return nil
		}
		if sendErr == nil {
			d.State, d.Result, d.Done, d.LastError, d.RetryAt = DeliveryDone, result, done, "", time.Time{}
			ended = true
			return nil
		}
		d.Attempts++
		d.LastError = clip(sendErr.Error(), 500)
		if d.Class != "" && printer != "" {
			d.Tried = append(d.Tried, printer)
			if len(d.Tried) > len(d.Tried)/2+16 {
				d.Tried = d.Tried[len(d.Tried)-16:]
			}
		}
		if final || d.Attempts >= maxAttempts {
			d.State, d.Done, d.RetryAt = DeliveryFailed, done, time.Time{}
			ended = true
			return nil
		}
		d.RetryAt = done.Add(retryDelay(d.Attempts))
		return nil
	})
	if err != nil {
		return nil, err
	}
	if sendErr != nil {
		m.o.Log.Warn("deliver", "delivery", d.ID, "job", d.Job, "node", d.Node, "printer", printer, "attempt", d.Attempts, "state", d.State, "err", sendErr)
	} else {
		m.o.Log.Info("delivered", "delivery", d.ID, "job", d.Job, "node", d.Node, "printer", printer, "result", result)
	}
	if ended {
		m.unpend(ctx, d.ID)
		if _, err := m.changeJob(ctx, jobID(d.Job), func(j *Job) error {
			if l := j.leg(d.ID); l != nil {
				l.State, l.Error = LegDone, ""
				if d.State == DeliveryFailed {
					l.State, l.Error = LegFailed, d.LastError
				}
				if printer != "" && l.Kind == KindClass {
					l.Member = printer
				}
			}
			j.Live = max(j.Live-1, 0)
			return nil
		}); err != nil && !errors.Is(err, ErrNotFound) {
			return d, err
		}
	}
	if sendErr != nil {
		return d, fmt.Errorf("%w: %s: %v", ErrPrint, d.Node, sendErr)
	}
	return d, nil
}

// cancelDelivery stops a pending delivery; one being sent finishes, and
// stays canceled.
func (m *Manager) cancelDelivery(ctx context.Context, id string) error {
	_, err := m.changeDelivery(ctx, id, func(d *Delivery) error {
		if d.State != DeliveryPending {
			return errSkip
		}
		d.State, d.Done = DeliveryCanceled, m.now()
		return nil
	})
	if err != nil {
		return err
	}
	m.unpend(ctx, id)
	return nil
}

// ArchiveRecord is the record an archive keeps beside a job's document.
type ArchiveRecord struct {
	Job         int32      `json:"job"`
	Dest        string     `json:"dest"`
	Name        string     `json:"name,omitempty"`
	SubmittedBy *Principal `json:"submitted_by,omitempty"`
	ClaimedUser string     `json:"claimed_user,omitempty"`
	Source      string     `json:"source"`
	Format      string     `json:"format,omitempty"`
	Pages       int        `json:"pages"`
	Size        int64      `json:"size"`
	SHA256      string     `json:"sha256"`
	Received    time.Time  `json:"received"`
	Archived    time.Time  `json:"archived"`
	Archive     string     `json:"archive"`
	Legs        []Leg      `json:"legs"`
}

// archiveKey names a job's objects in an archive: the day it was
// received, its id and the start of its hash, so that the objects of a
// day are listed together and a job is written once.
func archiveKey(j *Job) string {
	return j.Received.UTC().Format("2006-01-02") + "." + j.ID + "." + j.SHA256[:min(12, len(j.SHA256))]
}

// archive writes a job's document and its record to an archive's store,
// each once: an object there already is taken as written (a retry after
// the outcome of a write was not known).
func (m *Manager) archive(ctx context.Context, n *Node, j *Job) (string, error) {
	st, err := m.st.archive(ctx, n.Store, n.collection())
	if err != nil {
		return "", err
	}
	key := archiveKey(j)
	r, err := m.st.docs.Open(ctx, j.ID)
	if err != nil {
		return "", mapErr(err)
	}
	defer r.Close()
	if _, err := st.Insert(ctx, key+".bdf", r); err != nil && !errors.Is(err, yuruidb.ErrAlreadyExists) {
		return "", err
	}
	rec := ArchiveRecord{Job: j.Number, Dest: j.Dest, Name: j.Name, SubmittedBy: j.SubmittedBy, ClaimedUser: j.ClaimedUser, Source: j.Source, Format: j.Format, Pages: j.Pages, Size: j.Size,
		SHA256: j.SHA256, Received: j.Received, Archived: m.now(), Archive: n.Name, Legs: j.Legs}
	b, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return "", err
	}
	if _, err := st.Insert(ctx, key+".json", bytes.NewReader(b)); err != nil && !errors.Is(err, yuruidb.ErrAlreadyExists) {
		return "", err
	}
	return key, nil
}

// tickDeliveries sends the pending deliveries that are due.
func (m *Manager) tickDeliveries(ctx context.Context, rep *TickReport) error {
	recs, err := list[pendingRecord](ctx, m.st.db)
	if err != nil {
		return err
	}
	var errs []error
	for _, r := range recs {
		d, err := m.deliver(ctx, r.ID())
		switch {
		case errors.Is(err, ErrLeased):
			rep.Leased++
		case errors.Is(err, errNotDue):
		case err != nil:
			rep.Undelivered++
			errs = append(errs, err)
		case d != nil && d.State == DeliveryDone:
			rep.Delivered++
		}
	}
	return errors.Join(errs...)
}
