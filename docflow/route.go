package docflow

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/shibukawa/yuruidb"
)

// routing is a job on its way through the topology with a payload (the
// one received, or one a node made): the legs it gets, how many records
// that need its payloads were made, and the deliveries made, to send in
// the background.
type routing struct {
	job     *Job
	payload *Payload
	// archived is what the last archive on the way wrote (the job goes on
	// from an archive once it wrote): what an export lists.
	archived *archiveRef
	// trace is the step the payload came out of, or went into (trace.go).
	trace *Trace
	// cause is the failure the job left a node for along an on_error
	// edge: the deliveries made on the way carry it.
	cause      *Failure
	now        time.Time
	legs       []Leg
	live       int
	deliveries []string
}

// next is the index of the next leg, which names a delivery.
func (r *routing) next() int { return len(r.job.Legs) + len(r.legs) }

func (r *routing) add(l Leg) {
	l.At = r.now
	if l.Payload == "" && r.payload != nil {
		l.Payload = r.payload.ID
	}
	r.legs = append(r.legs, l)
}

// route takes a job to a node, and on along the edges: a slip at a
// location, a delivery to a printer or an archive, a step of a template or
// a convert node, an item in an inbox, and the next nodes of a queue, an
// entry (an endpoint, a CUPS server) or a class. A node that cannot take
// it (one that does not take the payload's kind: the graph is typed when
// a node is put, but a job may be on its way while it changes) gets a leg
// of its error; the job goes on elsewhere.
func (m *Manager) route(ctx context.Context, r *routing, name, from string, depth int) {
	if depth > maxDepth {
		r.add(Leg{Node: name, From: from, State: LegError, Error: "too deep: a cycle?"})
		return
	}
	kind, err := m.kindOf(ctx, name)
	if err != nil {
		r.add(Leg{Node: name, From: from, State: LegError, Error: "no such node"})
		return
	}
	if r.payload != nil && !takes(kind).has(r.payload.Kind) {
		r.add(Leg{Node: name, Kind: kind, From: from, State: LegError, Error: fmt.Sprintf("takes %s; this is %s", takes(kind), r.payload.Kind)})
		return
	}
	k := kindNamed(kind)
	switch {
	case k == nil:
		r.add(Leg{Node: name, Kind: kind, From: from, State: LegError, Error: "no such kind of node here: " + kind})
	case k.route != nil:
		k.route(m, ctx, r, name, from, depth)
	case k.Step != nil:
		m.routeDelivery(ctx, r, Leg{Node: name, Kind: kind, From: from}, "", "", kind)
	case k.Send != nil:
		m.routeDelivery(ctx, r, Leg{Node: name, Kind: kind, From: from}, name, "", kind)
	case k.PassesOn:
		n, err := m.Node(ctx, name)
		if err != nil {
			r.add(Leg{Node: name, Kind: kind, From: from, State: LegError, Error: err.Error()})
			return
		}
		r.add(Leg{Node: name, Kind: kind, From: from, State: LegRouted})
		for _, to := range n.To {
			m.route(ctx, r, to, name, depth+1)
		}
	default:
		r.add(Leg{Node: name, Kind: kind, From: from, State: LegError, Error: "cannot take a job"})
	}
}

// routeClass sends a job to a member of a class. A member that sends out
// (a printer) is the class's delivery: one that fails is skipped for the
// next (deliver).
func (m *Manager) routeClass(ctx context.Context, r *routing, name, from string, depth int) {
	n, err := m.Node(ctx, name)
	if err != nil {
		r.add(Leg{Node: name, Kind: KindClass, From: from, State: LegError, Error: err.Error()})
		return
	}
	member, err := m.pick(ctx, n)
	if err != nil {
		r.add(Leg{Node: name, Kind: KindClass, From: from, State: LegError, Error: err.Error()})
		return
	}
	if mk, err := m.kindOf(ctx, member); err == nil && sendsOut(mk) {
		m.routeDelivery(ctx, r, Leg{Node: name, Kind: KindClass, From: from, Member: member}, member, name, mk)
		return
	}
	r.add(Leg{Node: name, Kind: KindClass, From: from, State: LegRouted, Member: member})
	m.route(ctx, r, member, name, depth+1)
}

// pick chooses the member of a class a job goes to: the next in turn
// (round-robin: a count kept with the class), or the first (failover).
func (m *Manager) pick(ctx context.Context, n *Node) (string, error) {
	if len(n.Members) == 0 {
		return "", errors.New("the class has no members")
	}
	if n.Policy == PolicyFailover {
		return n.Members[0], nil
	}
	var i int32
	if _, err := change(ctx, m.st.db, "rr."+n.Name, true, func(c *counterRecord) error {
		i = c.Next
		c.Next++
		if c.Next < 0 {
			c.Next = 0
		}
		return nil
	}); err != nil {
		return "", err
	}
	return n.Members[int(i)%len(n.Members)], nil
}

// routeDisplay makes the job a slip at a display.
func (m *Manager) routeDisplay(ctx context.Context, r *routing, loc, from string) {
	l, err := m.display(ctx, loc)
	if err != nil {
		r.add(Leg{Node: loc, Kind: KindDisplay, From: from, State: LegError, Error: err.Error()})
		return
	}
	j, p := r.job, r.payload
	s := Slip{ID: slipID(j.Number, loc), Job: j.Number, Location: loc, Name: j.Name, SubmittedBy: j.SubmittedBy, ClaimedUser: j.ClaimedUser, Source: j.Source, Format: j.Format,
		Payload: p.ID, Pages: p.Pages, Size: p.Size, SHA256: p.SHA256, Received: r.now, State: SlipNew}
	if l.Fallback != nil && l.Fallback.After > 0 {
		s.FallbackAt = r.now.Add(l.Fallback.After.d())
	}
	_, err = m.st.db.Insert(ctx, s.ID, slipRecord(s))
	if errors.Is(err, yuruidb.ErrAlreadyExists) {
		// reached twice, by two paths: shown once
		r.add(Leg{Node: loc, Kind: KindDisplay, From: from, State: LegShown, Ref: s.ID})
		return
	}
	if err != nil {
		r.add(Leg{Node: loc, Kind: KindDisplay, From: from, State: LegError, Error: mapErr(err).Error()})
		return
	}
	r.live++
	if err := m.touch(ctx, loc, func(st *locationState) { st.Open = append(st.Open, s.ID) }); err != nil {
		r.add(Leg{Node: loc, Kind: KindDisplay, From: from, State: LegShown, Ref: s.ID, Error: err.Error()})
		return
	}
	r.add(Leg{Node: loc, Kind: KindDisplay, From: from, State: LegShown, Ref: s.ID})
}

// Job returns a job by its id, or its number.
func (m *Manager) Job(ctx context.Context, id string) (*Job, error) {
	if n, err := strconv.ParseInt(id, 10, 32); err == nil && n > 0 && len(id) <= 10 {
		id = jobID(int32(n))
	}
	rec, err := m.st.db.Get[jobRecord](ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	j := Job(rec.Data)
	return &j, nil
}

// JobOfNumber returns the job of a number.
func (m *Manager) JobOfNumber(ctx context.Context, n int32) (*Job, error) {
	if n <= 0 {
		return nil, ErrNotFound
	}
	return m.Job(ctx, jobID(n))
}

// changeJob changes a job and derives its state; a job that ends is
// listed for its removal, under the day it ended.
func (m *Manager) changeJob(ctx context.Context, id string, fn func(*Job) error) (*Job, error) {
	now := m.now()
	ended := false
	rec, err := change(ctx, m.st.db, id, false, func(r *jobRecord) error {
		j := (*Job)(r)
		if err := fn(j); err != nil {
			return err
		}
		was := j.Ended
		j.derive(now)
		ended = was.IsZero() && !j.Ended.IsZero()
		return nil
	})
	if err != nil {
		return nil, err
	}
	j := Job(rec)
	if ended {
		if err := m.endedOn(ctx, &j); err != nil {
			return &j, err
		}
	}
	return &j, nil
}

// endedOn lists a job under the day it ended.
func (m *Manager) endedOn(ctx context.Context, j *Job) error {
	day := j.Ended.UTC().Format("2006-01-02")
	_, err := change(ctx, m.st.db, day, true, func(d *dayRecord) error {
		if !slices.Contains(d.Jobs, j.ID) {
			d.Jobs = append(d.Jobs, j.ID)
		}
		return nil
	})
	return err
}

// release tells that a record that needed the job's document is gone.
func (m *Manager) release(ctx context.Context, job int32) {
	if _, err := m.changeJob(ctx, jobID(job), func(j *Job) error {
		j.Live = max(j.Live-1, 0)
		return nil
	}); err != nil && !errors.Is(err, ErrNotFound) {
		m.o.Log.Warn("release a job's document", "job", job, "err", err)
	}
}

// touchDest lists a job under its destination.
func (m *Manager) touchDest(ctx context.Context, dest, id string) error {
	_, err := change(ctx, m.st.db, dest, true, func(d *destState) error {
		d.Jobs = append(d.Jobs, id)
		if len(d.Jobs) > recentJobs {
			d.Jobs = d.Jobs[len(d.Jobs)-recentJobs:]
		}
		return nil
	})
	return err
}

// JobsOf returns the last jobs sent to a destination, latest first.
func (m *Manager) JobsOf(ctx context.Context, dest string, limit int) ([]Job, error) {
	if _, err := m.kindOf(ctx, dest); err != nil {
		return nil, err
	}
	return m.jobsOfState(ctx, dest, limit)
}

// RecentJobs returns the last jobs of every destination, latest first.
func (m *Manager) RecentJobs(ctx context.Context, limit int) ([]Job, error) {
	return m.jobsOfState(ctx, allJobs, limit)
}

// jobsOfState reads the jobs a destination state lists.
func (m *Manager) jobsOfState(ctx context.Context, dest string, limit int) ([]Job, error) {
	rec, err := m.st.db.Get[destState](ctx, dest)
	if errors.Is(err, yuruidb.ErrNotFound) {
		return []Job{}, nil
	}
	if err != nil {
		return nil, mapErr(err)
	}
	if limit <= 0 {
		limit = recentJobs
	}
	out := make([]Job, 0, min(limit, len(rec.Data.Jobs)))
	for i := len(rec.Data.Jobs) - 1; i >= 0 && len(out) < limit; i-- {
		j, err := m.Job(ctx, rec.Data.Jobs[i])
		if errors.Is(err, ErrNotFound) {
			continue // removed
		}
		if err != nil {
			return nil, err
		}
		out = append(out, *j)
	}
	return out, nil
}

// Cancel cancels a job: its slips not done leave the boards, its items
// leave their inboxes, and its deliveries not sent are not. A job that
// ended is canceled while a slip of it is still open (a slip is done by
// its reader, after the job completed).
func (m *Manager) Cancel(ctx context.Context, id string, by Mark) (*Job, error) {
	by.At = m.now()
	cur, err := m.Job(ctx, id)
	if err != nil {
		return nil, err
	}
	openSlip := false
	if cur.State.ended() {
		for _, l := range cur.Legs {
			if l.Kind == KindDisplay && l.State == LegShown {
				if s, err := m.Slip(ctx, l.Ref); err == nil && !s.State.ended() {
					openSlip = true
				}
			}
		}
	}
	var slips, items, deliveries []string
	j, err := m.changeJob(ctx, id, func(j *Job) error {
		if j.State == JobCanceled || (j.State.ended() && !openSlip) {
			return fmt.Errorf("%w: the job is %s", ErrConflict, j.State)
		}
		j.Canceled = &by
		for i := range j.Legs {
			l := &j.Legs[i]
			switch {
			case l.Kind == KindDisplay && l.State == LegShown:
				slips = append(slips, l.Ref)
			case l.Kind == KindInbox && l.State == LegHeld:
				items = append(items, l.Ref)
				l.State = LegDiscarded
			case l.State == LegPending:
				deliveries = append(deliveries, l.Ref)
				l.State = LegCanceled
				j.Live = max(j.Live-1, 0)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	ctx = context.WithoutCancel(ctx)
	var errs []error
	for _, id := range slips {
		if _, err := m.CancelSlip(ctx, id, by); err != nil && !errors.Is(err, ErrConflict) && !errors.Is(err, ErrNotFound) {
			errs = append(errs, err)
		}
	}
	for _, id := range items {
		if err := m.discardItem(ctx, id, by); err != nil && !errors.Is(err, ErrConflict) && !errors.Is(err, ErrNotFound) {
			errs = append(errs, err)
		}
	}
	for _, id := range deliveries {
		if err := m.cancelDelivery(ctx, id); err != nil && !errors.Is(err, ErrNotFound) {
			errs = append(errs, err)
		}
	}
	m.o.Log.Info("job canceled", "job", j.Number, "by", by.Worker)
	return j, errors.Join(errs...)
}

// tickJobs removes the jobs that ended longer ago than they are kept, and
// that nothing shows any more, with their documents and deliveries.
func (m *Manager) tickJobs(ctx context.Context, rep *TickReport) error {
	days, err := list[dayRecord](ctx, m.st.db)
	if err != nil {
		return err
	}
	now := m.now()
	cut := now.Add(-m.o.JobKeep).UTC().Format("2006-01-02")
	var errs []error
	for _, d := range days {
		if d.ID() > cut {
			break // the days are listed in order
		}
		var removed []string
		for _, id := range d.Data.Jobs {
			j, err := m.Job(ctx, id)
			if errors.Is(err, ErrNotFound) {
				removed = append(removed, id)
				continue
			}
			if err != nil {
				errs = append(errs, err)
				continue
			}
			if j.Live > 0 || !j.State.ended() || now.Sub(j.Ended) < m.o.JobKeep {
				continue // a slip or an item still shows it, or it is kept yet
			}
			if err := m.removeJob(ctx, j); err != nil {
				errs = append(errs, err)
				continue
			}
			removed = append(removed, id)
			rep.Jobs++
		}
		if len(removed) == 0 {
			continue
		}
		rec, err := change(ctx, m.st.db, d.ID(), false, func(r *dayRecord) error {
			r.Jobs = slices.DeleteFunc(r.Jobs, func(id string) bool { return slices.Contains(removed, id) })
			return nil
		})
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if len(rec.Jobs) == 0 {
			if cur, err := m.st.db.Get[dayRecord](ctx, d.ID()); err == nil && len(cur.Data.Jobs) == 0 {
				m.st.db.Delete[dayRecord](ctx, d.ID(), cur.Revision())
			}
		}
	}
	return errors.Join(errs...)
}

// namesDelivery tells a kind of node whose leg names a delivery: an
// exit that sends out, a class, an archive, a step.
func namesDelivery(kind string) bool {
	switch kind {
	case KindArchive, KindClass, KindTemplate, KindConvert, KindSign:
		return true
	}
	return sendsOut(kind) || steps(kind)
}

// removeJob deletes a job, its payloads and its deliveries.
func (m *Manager) removeJob(ctx context.Context, j *Job) error {
	for _, l := range j.Legs {
		if l.Ref == "" || !namesDelivery(l.Kind) {
			continue
		}
		if rec, err := m.st.db.Get[deliveryRecord](ctx, l.Ref); err == nil {
			if err := m.st.db.Delete[deliveryRecord](ctx, l.Ref, rec.Revision()); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
				return mapErr(err)
			}
		}
	}
	ids := []string{j.ID}
	for _, p := range j.Payloads {
		if !slices.Contains(ids, p.ID) {
			ids = append(ids, p.ID)
		}
	}
	for _, id := range ids {
		if err := m.st.docs.Remove(ctx, id); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
			return err
		}
	}
	rec, err := m.st.db.Get[jobRecord](ctx, j.ID)
	if err == nil {
		err = m.st.db.Delete[jobRecord](ctx, j.ID, rec.Revision())
	}
	if err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
		return mapErr(err)
	}
	return nil
}
