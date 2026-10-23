package docflow

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/shibukawa/yuruidb"
)

// RenameNode gives a node another name. What is kept under its name
// moves with it: where it is drawn, the edges of the nodes that send to
// it, a template node's templates, a display's devices, enrollment codes
// and the slips it still shows, an inbox's items it still shows, a sink's
// documents, an export's entries, and the list of the jobs sent to it.
// An archive keeps writing to the collection it wrote to. Where the node
// is reached changes with it (/ipp/print/<name>, /in/<name>,
// /out/<name>, a CUPS server's /printers/<name>); the jobs that reached
// it keep the name they had. A node with work on its way (deliveries
// waiting or failed, items held, slips open) is refused, and so is one
// the manager's options name (on_error).
func (m *Manager) RenameNode(ctx context.Context, from, to string) error {
	if !validName(to) {
		return fmt.Errorf("%w: node name %q: letters, digits, . _ and -, up to 64, starting with a letter or a digit", ErrInvalid, to)
	}
	cur, err := m.Node(ctx, from)
	if err != nil {
		return err
	}
	if from == to {
		return nil
	}
	if _, err := m.Node(ctx, to); err == nil {
		return fmt.Errorf("%w: %q is a node already", ErrConflict, to)
	} else if !errors.Is(err, ErrNotFound) {
		return err
	}
	if slices.Contains(m.o.OnError, from) {
		return fmt.Errorf("%w: the manager's on_error names %s: change docflow.on_error first", ErrConflict, from)
	}
	if err := m.idle(ctx, cur); err != nil {
		return err
	}
	n := *cur
	n.Name = to
	if n.Kind == KindArchive {
		n.Collection = cur.collection() // the objects stay where they are
	}
	if _, err := m.st.db.Insert(ctx, to, nodeRecord(n)); err != nil {
		if errors.Is(err, yuruidb.ErrAlreadyExists) {
			return fmt.Errorf("%w: %q is a node already", ErrConflict, to)
		}
		return mapErr(err)
	}
	if err := m.moveState(ctx, cur, to); err != nil {
		return err
	}
	nodes, err := m.Nodes(ctx)
	if err != nil {
		return err
	}
	for _, x := range nodes {
		if x.Name == from || x.Name == to || !slices.Contains(x.edges(), from) {
			continue
		}
		if _, err := change(ctx, m.st.db, x.Name, false, func(r *nodeRecord) error {
			if !(*Node)(r).renameEdge(from, to) {
				return errSkip
			}
			return nil
		}); err != nil {
			return err
		}
	}
	rec, err := m.st.db.Get[nodeRecord](ctx, from)
	if err == nil {
		err = m.st.db.Delete[nodeRecord](ctx, from, rec.Revision())
	}
	if err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
		return mapErr(err)
	}
	m.o.Log.Info("node renamed", "from", from, "to", to, "kind", n.Kind)
	m.changed()
	return nil
}

// idle refuses a node with work on its way under its name.
func (m *Manager) idle(ctx context.Context, n *Node) error {
	ds, err := m.Pending(ctx)
	if err != nil {
		return err
	}
	waiting := 0
	for _, d := range ds {
		if d.Node == n.Name || d.Printer == n.Name || d.Class == n.Name {
			waiting++
		}
	}
	if waiting > 0 {
		return fmt.Errorf("%w: %d deliveries are waiting or failed at %s: rename it once they are done (Admin → Deliveries)", ErrConflict, waiting, n.Name)
	}
	switch n.Kind {
	case KindInbox:
		if st, err := m.st.db.Get[inboxState](ctx, n.Name); err == nil && len(st.Data.Held) > 0 {
			return fmt.Errorf("%w: inbox %s holds %d items: release or discard them first", ErrConflict, n.Name, len(st.Data.Held))
		}
	case KindDisplay:
		if st, err := m.st.db.Get[locationState](ctx, n.Name); err == nil && len(st.Data.Open) > 0 {
			return fmt.Errorf("%w: display %s shows %d open slips: rename it once they are done or canceled", ErrConflict, n.Name, len(st.Data.Open))
		}
	}
	return nil
}

// moveState moves what is kept under a node's name to another name.
func (m *Manager) moveState(ctx context.Context, n *Node, to string) error {
	from := n.Name
	if _, err := change(ctx, m.st.db, layoutID, false, func(r *layoutRecord) error {
		p, ok := r.Positions[from]
		if !ok {
			return errSkip
		}
		delete(r.Positions, from)
		r.Positions[to] = p
		return nil
	}); err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}
	if err := moveRecord(ctx, m.st.db, from, to, func(d *destState, s destState) { d.Jobs = append(s.Jobs, d.Jobs...) }); err != nil {
		return err
	}
	switch n.Kind {
	case KindTemplate:
		return m.moveTemplates(ctx, from, to)
	case KindClass:
		return moveRecord(ctx, m.st.db, "rr."+from, "rr."+to, func(d *counterRecord, s counterRecord) { d.Next = s.Next })
	case KindSink:
		return moveRecord(ctx, m.st.db, from, to, func(d *sinkState, s sinkState) { d.Seq, d.Jobs = d.Seq+s.Seq+1, append(s.Jobs, d.Jobs...) })
	case KindInbox:
		var ended []string
		if err := moveRecord(ctx, m.st.db, from, to, func(d *inboxState, s inboxState) {
			ended = s.Ended
			d.Seq, d.Held, d.Ended = d.Seq+s.Seq+1, append(s.Held, d.Held...), append(s.Ended, d.Ended...)
		}); err != nil {
			return err
		}
		for _, id := range ended {
			if _, err := change(ctx, m.st.db, id, false, func(r *itemRecord) error { r.Inbox = to; return nil }); err != nil && !errors.Is(err, ErrNotFound) {
				return err
			}
		}
	case KindDisplay:
		var st locationState
		if err := moveRecord(ctx, m.st.db, from, to, func(d *locationState, s locationState) {
			st = s
			d.Seq, d.Open, d.Ended, d.Devices = d.Seq+s.Seq+1, append(s.Open, d.Open...), append(s.Ended, d.Ended...), append(s.Devices, d.Devices...)
		}); err != nil {
			return err
		}
		for _, id := range st.Devices {
			if _, err := change(ctx, m.st.db, id, false, func(r *deviceRecord) error { r.Location = to; return nil }); err != nil && !errors.Is(err, ErrNotFound) {
				return err
			}
		}
		for _, id := range append(st.Open, st.Ended...) {
			if _, err := change(ctx, m.st.db, id, false, func(r *slipRecord) error { r.Location = to; return nil }); err != nil && !errors.Is(err, ErrNotFound) {
				return err
			}
		}
		codes, err := list[enrollRecord](ctx, m.st.db)
		if err != nil {
			return err
		}
		for _, c := range codes {
			if c.Data.Location != from {
				continue
			}
			if _, err := change(ctx, m.st.db, c.ID(), false, func(r *enrollRecord) error { r.Location = to; return nil }); err != nil && !errors.Is(err, ErrNotFound) {
				return err
			}
		}
	case KindExport:
		if err := m.moveExport(ctx, from, to); err != nil {
			return err
		}
		return m.moveListedAt(ctx, from, to)
	}
	return nil
}

// moveRecord moves the record of a collection under one id to another,
// merged into one there already (a job that reached the new name
// meanwhile).
func moveRecord[T any](ctx context.Context, db *yuruidb.Store, from, to string, merge func(dst *T, src T)) error {
	rec, err := db.Get[T](ctx, from)
	if errors.Is(err, yuruidb.ErrNotFound) {
		return nil
	}
	if err != nil {
		return mapErr(err)
	}
	if _, err := change(ctx, db, to, true, func(d *T) error {
		merge(d, rec.Data)
		return nil
	}); err != nil {
		return err
	}
	if err := db.Delete[T](ctx, from, rec.Revision()); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
		return mapErr(err)
	}
	return nil
}

// moveTemplates moves a template node's templates, records and files.
func (m *Manager) moveTemplates(ctx context.Context, from, to string) error {
	ts, err := m.templatesOf(ctx, from)
	if err != nil {
		return err
	}
	for _, t := range ts {
		old, rid := templateID(from, t.ID), templateID(to, t.ID)
		b, err := readBlob(ctx, m.st.tpls, old, MaxTemplate)
		if err != nil {
			return err
		}
		if err := m.st.tpls.Remove(ctx, rid); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
			return mapErr(err)
		}
		if _, err := m.st.tpls.Insert(ctx, rid, bytes.NewReader(b)); err != nil {
			return mapErr(err)
		}
		t.Node = to
		if _, err := change(ctx, m.st.db, rid, true, func(r *templateRecord) error {
			*r = templateRecord(t)
			return nil
		}); err != nil {
			return err
		}
		if err := m.DeleteTemplate(ctx, from, t.ID); err != nil && !errors.Is(err, ErrNotFound) {
			return err
		}
	}
	return nil
}

// moveExport moves an export's entries (<name>.<job>) and lists (<name>,
// <name>.<owner's hash>).
func (m *Manager) moveExport(ctx context.Context, from, to string) error {
	entries, err := list[exportEntryRecord](ctx, m.st.db)
	if err != nil {
		return err
	}
	for _, r := range entries {
		if rest, ok := strings.CutPrefix(r.ID(), from+"."); ok && len(rest) == 10 && isDigits(rest) {
			if err := moveRecord(ctx, m.st.db, r.ID(), exportEntryID(to, rest), func(d *exportEntryRecord, s exportEntryRecord) { *d = s }); err != nil {
				return err
			}
		}
	}
	lists, err := list[exportList](ctx, m.st.db)
	if err != nil {
		return err
	}
	for _, r := range lists {
		id := ""
		if r.ID() == from {
			id = to
		} else if rest, ok := strings.CutPrefix(r.ID(), from+"."); ok && len(rest) == 16 && isHex(rest) {
			id = to + "." + rest
		}
		if id == "" {
			continue
		}
		if err := moveRecord(ctx, m.st.db, r.ID(), id, func(d *exportList, s exportList) { d.Seq, d.Jobs = d.Seq+s.Seq+1, append(s.Jobs, d.Jobs...) }); err != nil {
			return err
		}
	}
	return nil
}

// renameEdge points the node's edges to a node at its new name, telling
// whether there were any.
func (n *Node) renameEdge(from, to string) bool {
	changed := false
	swap := func(names []string) {
		for i, x := range names {
			if x == from {
				names[i] = to
				changed = true
			}
		}
	}
	swap(n.To)
	swap(n.Original)
	swap(n.Members)
	swap(n.Otherwise)
	swap(n.OnError)
	for i := range n.Rules {
		swap(n.Rules[i].To)
	}
	if n.Fallback != nil && n.Fallback.Printer == from {
		n.Fallback.Printer = to
		changed = true
	}
	return changed
}
