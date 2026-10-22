package docflow

import (
	"context"
	"errors"
	"slices"
	"strconv"

	"github.com/shibukawa/yuruidb"
)

// A sink keeps the latest documents that reach it, to look at on the
// topology: what an endpoint or a queue sends, seen without a printer or
// a device. It keeps the last Latest jobs (their documents stay while a
// sink keeps them); an older one is let go when a newer comes.

const defaultLatest = 10

func (n *Node) latest() int { return or(n.Latest, defaultLatest) }

// sinkState is the jobs a sink keeps, oldest first.
type sinkState struct {
	Seq  int64    `json:"seq"`
	Jobs []string `json:"jobs,omitempty"`
}

// routeSink keeps the job at a sink, letting the oldest go past Latest.
func (m *Manager) routeSink(ctx context.Context, r *routing, name, from string) {
	n, err := m.nodeOfKind(ctx, name, KindSink)
	if err != nil {
		r.add(Leg{Node: name, Kind: KindSink, From: from, State: LegError, Error: err.Error()})
		return
	}
	j := r.job
	var pushed []string
	kept := false
	err = m.touchSink(ctx, name, func(st *sinkState) {
		pushed, kept = nil, false
		if slices.Contains(st.Jobs, j.ID) {
			kept = true // reached twice, by two paths
			return
		}
		st.Jobs = append(st.Jobs, j.ID)
		if extra := len(st.Jobs) - n.latest(); extra > 0 {
			pushed = slices.Clone(st.Jobs[:extra])
			st.Jobs = st.Jobs[extra:]
		}
	})
	if err != nil {
		r.add(Leg{Node: name, Kind: KindSink, From: from, State: LegError, Error: err.Error()})
		return
	}
	if !kept {
		r.live++
	}
	r.add(Leg{Node: name, Kind: KindSink, From: from, State: LegKept})
	for _, id := range pushed {
		m.releaseID(ctx, id)
	}
}

// releaseID tells that a record of a job's id no longer needs its
// document.
func (m *Manager) releaseID(ctx context.Context, id string) {
	if n, err := strconv.ParseInt(id, 10, 32); err == nil {
		m.release(ctx, int32(n))
	}
}

func (m *Manager) touchSink(ctx context.Context, name string, fn func(*sinkState)) error {
	_, err := change(ctx, m.st.db, name, true, func(s *sinkState) error {
		fn(s)
		s.Seq++
		return nil
	})
	return err
}

// sinkJobs are the ids of the jobs a sink keeps now, oldest first.
func (m *Manager) sinkJobs(ctx context.Context, name string) []string {
	rec, err := m.st.db.Get[sinkState](ctx, name)
	if err != nil {
		return nil
	}
	return rec.Data.Jobs
}

// Latest returns the jobs a sink keeps, newest first: those the viewer
// may see (nil: every one).
func (m *Manager) Latest(ctx context.Context, name string, v *Viewer) ([]Job, error) {
	if _, err := m.nodeOfKind(ctx, name, KindSink); err != nil {
		return nil, err
	}
	ids := m.sinkJobs(ctx, name)
	out := []Job{}
	for i := len(ids) - 1; i >= 0; i-- {
		j, err := m.Job(ctx, ids[i])
		if errors.Is(err, ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if v == nil || v.CanSee(j) {
			out = append(out, *j)
		}
	}
	return out, nil
}

// dropSink lets go of the jobs a sink keeps and forgets it (the sink is
// deleted).
func (m *Manager) dropSink(ctx context.Context, name string) error {
	rec, err := m.st.db.Get[sinkState](ctx, name)
	if errors.Is(err, yuruidb.ErrNotFound) {
		return nil
	}
	if err != nil {
		return mapErr(err)
	}
	for _, id := range rec.Data.Jobs {
		m.releaseID(ctx, id)
	}
	if err := m.st.db.Delete[sinkState](ctx, name, rec.Revision()); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
		return mapErr(err)
	}
	return nil
}
