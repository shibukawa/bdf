package manager

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/shibukawa/yuruidb"
)

// An inbox holds the jobs that reach it as items, until a person
// releases one to a node (one of the inbox's To, when it names any) or
// discards it: a filter worked by hand. An item released or discarded
// stays, shown as ended, for the inbox's keep.

// routeInbox holds the job in an inbox.
func (m *Manager) routeInbox(ctx context.Context, r *routing, name, from string) {
	j, p := r.job, r.payload
	it := Item{ID: itemID(j.Number, name), Job: j.Number, Inbox: name, Name: j.Name, SubmittedBy: j.SubmittedBy, ClaimedUser: j.ClaimedUser, Dest: j.Dest, Payload: p.ID, Trace: r.trace, Pages: p.Pages, Received: r.now, State: ItemHeld}
	_, err := m.st.db.Insert(ctx, it.ID, itemRecord(it))
	if errors.Is(err, yuruidb.ErrAlreadyExists) {
		r.add(Leg{Node: name, Kind: KindInbox, From: from, State: LegError, Ref: it.ID, Error: "the job passed this inbox already"})
		return
	}
	if err != nil {
		r.add(Leg{Node: name, Kind: KindInbox, From: from, State: LegError, Error: mapErr(err).Error()})
		return
	}
	r.live++
	if err := m.touchInbox(ctx, name, func(st *inboxState) { st.Held = append(st.Held, it.ID) }); err != nil {
		r.add(Leg{Node: name, Kind: KindInbox, From: from, State: LegHeld, Ref: it.ID, Error: err.Error()})
		return
	}
	r.add(Leg{Node: name, Kind: KindInbox, From: from, State: LegHeld, Ref: it.ID})
}

// touchInbox changes an inbox's state.
func (m *Manager) touchInbox(ctx context.Context, name string, fn func(*inboxState)) error {
	_, err := change(ctx, m.st.db, name, true, func(s *inboxState) error {
		fn(s)
		s.Seq++
		return nil
	})
	return err
}

// Item returns an item of an inbox.
func (m *Manager) Item(ctx context.Context, id string) (*Item, error) {
	rec, err := m.st.db.Get[itemRecord](ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	it := Item(rec.Data)
	return &it, nil
}

// InboxView is an inbox as a person works it: the items held, oldest
// first, and the ended ones still shown, latest first.
type InboxView struct {
	Name  string   `json:"name"`
	Title string   `json:"title"`
	To    []string `json:"to"`
	Rev   string   `json:"rev"`
	Held  []Item   `json:"held"`
	Ended []Item   `json:"ended"`
}

// Inbox returns an inbox's items.
func (m *Manager) Inbox(ctx context.Context, name string) (*InboxView, error) {
	n, err := m.Node(ctx, name)
	if err != nil {
		return nil, err
	}
	if n.Kind != KindInbox {
		return nil, fmt.Errorf("%w: %s is a %s", ErrNotFound, name, n.Kind)
	}
	v := &InboxView{Name: name, Title: or(n.Title, name), To: n.To, Rev: "0", Held: []Item{}, Ended: []Item{}}
	if v.To == nil {
		v.To = []string{}
	}
	rec, err := m.st.db.Get[inboxState](ctx, name)
	if errors.Is(err, yuruidb.ErrNotFound) {
		return v, nil
	}
	if err != nil {
		return nil, mapErr(err)
	}
	v.Rev = fmt.Sprintf("%d", rec.Data.Seq)
	read := func(ids []string) ([]Item, error) {
		out := make([]Item, 0, len(ids))
		for _, id := range ids {
			it, err := m.Item(ctx, id)
			if errors.Is(err, ErrNotFound) {
				continue
			}
			if err != nil {
				return nil, err
			}
			out = append(out, *it)
		}
		return out, nil
	}
	if v.Held, err = read(rec.Data.Held); err != nil {
		return nil, err
	}
	ended := slices.Clone(rec.Data.Ended)
	slices.Reverse(ended)
	if v.Ended, err = read(ended); err != nil {
		return nil, err
	}
	return v, nil
}

// Release sends a held item on to nodes: those named, or all of the
// inbox's To when none is. When the inbox names nodes, only those may be
// chosen.
func (m *Manager) Release(ctx context.Context, id string, to []string, by Mark) (*Item, error) {
	by.At = m.now()
	it, err := m.Item(ctx, id)
	if err != nil {
		return nil, err
	}
	if it.State != ItemHeld {
		return nil, fmt.Errorf("%w: the item is %s", ErrConflict, it.State)
	}
	var allowed []string
	if n, err := m.Node(ctx, it.Inbox); err == nil && n.Kind == KindInbox {
		allowed = n.To
	}
	to = cleanNames(to)
	if len(to) == 0 {
		to = allowed
	}
	if len(to) == 0 {
		return nil, fmt.Errorf("%w: name the nodes to release to", ErrInvalid)
	}
	for _, name := range to {
		if len(allowed) > 0 && !slices.Contains(allowed, name) {
			return nil, fmt.Errorf("%w: %s releases to %s only", ErrInvalid, it.Inbox, strings.Join(allowed, ", "))
		}
		if _, err := m.kindOf(ctx, name); err != nil {
			return nil, fmt.Errorf("%w: no node %q", ErrInvalid, name)
		}
	}
	j, err := m.JobOfNumber(ctx, it.Job)
	if err != nil {
		return nil, err
	}
	rel := &Release{To: to, Worker: by.Worker, At: by.At}
	it, err = m.changeItem(ctx, id, func(it *Item) error {
		if it.State != ItemHeld {
			return fmt.Errorf("%w: the item is %s", ErrConflict, it.State)
		}
		it.State, it.Released = ItemReleased, rel
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := m.endItem(ctx, it.Inbox, id); err != nil {
		return it, err
	}
	r := &routing{job: j, now: by.At, payload: j.payload(it.Payload), trace: it.Trace}
	for _, name := range to {
		m.route(ctx, r, name, it.Inbox, 0)
	}
	if _, err := m.changeJob(ctx, j.ID, func(j *Job) error {
		if l := j.leg(id); l != nil {
			l.State = LegReleased
		}
		j.Legs = append(j.Legs, r.legs...)
		j.Live += r.live
		return nil
	}); err != nil {
		return it, err
	}
	m.o.Log.Info("item released", "inbox", it.Inbox, "job", it.Job, "to", to, "by", by.Worker)
	m.kick(r.deliveries)
	return it, nil
}

// Discard drops a held item: the job goes no further from here.
func (m *Manager) Discard(ctx context.Context, id string, by Mark) (*Item, error) {
	by.At = m.now()
	if err := m.discardItem(ctx, id, by); err != nil {
		return nil, err
	}
	it, err := m.Item(ctx, id)
	if err != nil {
		return nil, err
	}
	if _, err := m.changeJob(ctx, jobID(it.Job), func(j *Job) error {
		if l := j.leg(id); l != nil {
			l.State = LegDiscarded
		}
		return nil
	}); err != nil && !errors.Is(err, ErrNotFound) {
		return it, err
	}
	m.o.Log.Info("item discarded", "inbox", it.Inbox, "job", it.Job, "by", by.Worker)
	return it, nil
}

// discardItem marks an item discarded and moves it to the ended ones.
func (m *Manager) discardItem(ctx context.Context, id string, by Mark) error {
	it, err := m.changeItem(ctx, id, func(it *Item) error {
		if it.State != ItemHeld {
			return fmt.Errorf("%w: the item is %s", ErrConflict, it.State)
		}
		it.State, it.Discarded = ItemDiscarded, &by
		return nil
	})
	if err != nil {
		return err
	}
	return m.endItem(ctx, it.Inbox, id)
}

// end moves an item from the held ones to the ended ones.
func (m *Manager) endItem(ctx context.Context, inbox, id string) error {
	return m.touchInbox(ctx, inbox, func(st *inboxState) {
		st.Held = slices.DeleteFunc(st.Held, func(x string) bool { return x == id })
		if !slices.Contains(st.Ended, id) {
			st.Ended = append(st.Ended, id)
		}
	})
}

func (m *Manager) changeItem(ctx context.Context, id string, fn func(*Item) error) (*Item, error) {
	r, err := change(ctx, m.st.db, id, false, func(r *itemRecord) error {
		return fn((*Item)(r))
	})
	if err != nil {
		return nil, err
	}
	it := Item(r)
	return &it, nil
}

// tickInbox removes the ended items an inbox showed for its keep.
func (m *Manager) tickInbox(ctx context.Context, n *Node, rep *TickReport) error {
	rec, err := m.st.db.Get[inboxState](ctx, n.Name)
	if err != nil {
		if errors.Is(err, yuruidb.ErrNotFound) {
			return nil
		}
		return mapErr(err)
	}
	now := m.now()
	var errs []error
	for _, id := range rec.Data.Ended {
		it, err := m.Item(ctx, id)
		if errors.Is(err, ErrNotFound) {
			it = &Item{ID: id, Inbox: n.Name}
		} else if err != nil {
			continue
		} else if at := it.endedAt(); at.IsZero() || now.Sub(at) < n.keep() {
			continue
		}
		if err := m.removeItem(ctx, it); err != nil {
			errs = append(errs, err)
			continue
		}
		rep.Items++
	}
	return errors.Join(errs...)
}

// removeItem deletes an ended item; the job's document is released.
func (m *Manager) removeItem(ctx context.Context, it *Item) error {
	rec, err := m.st.db.Get[itemRecord](ctx, it.ID)
	switch {
	case err == nil:
		if err := m.st.db.Delete[itemRecord](ctx, it.ID, rec.Revision()); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
			return mapErr(err)
		}
		m.release(ctx, it.Job)
	case !errors.Is(err, yuruidb.ErrNotFound):
		return mapErr(err)
	}
	return m.touchInbox(ctx, it.Inbox, func(st *inboxState) {
		st.Ended = slices.DeleteFunc(st.Ended, func(id string) bool { return id == it.ID })
		st.Held = slices.DeleteFunc(st.Held, func(id string) bool { return id == it.ID })
	})
}
