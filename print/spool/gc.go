package spool

import (
	"context"
	"strings"
	"time"
)

// Collect removes what the retention no longer keeps: the documents of
// the jobs done after Retention.Documents, the records of the jobs done
// after Retention.Records, and documents of no job (a crash while one
// arrived). The runner calls it every hour.
func (s *Spool) Collect(ctx context.Context) error {
	now := s.o.Now()
	docCut := now.Add(-s.o.Retention.Documents)
	recCut := now.Add(-s.o.Retention.Records)
	// the documents of the jobs done recently
	s.mu.Lock()
	var expired []int32
	keep := s.recent[:0]
	for _, j := range s.recent {
		if j.Finished.Before(docCut) {
			expired = append(expired, j.ID)
		} else {
			keep = append(keep, j)
		}
	}
	s.recent = keep
	s.mu.Unlock()
	for _, id := range expired {
		if err := s.removeDocuments(ctx, id); err != nil {
			return err
		}
	}
	// the records past their retention, and documents a crash left
	var old []struct {
		id  int32
		rev string
	}
	var stray []int32
	err := s.st.eachDone(ctx, func(j Job, rev string) bool {
		switch {
		case j.Finished.Before(recCut):
			old = append(old, struct {
				id  int32
				rev string
			}{j.ID, rev})
		case j.Finished.Before(docCut) && !j.DocumentsRemoved:
			stray = append(stray, j.ID)
		}
		return true
	})
	if err != nil {
		return err
	}
	for _, id := range stray {
		if err := s.removeDocuments(ctx, id); err != nil {
			return err
		}
	}
	for _, r := range old {
		s.st.removeDoc(ctx, originalKey(r.id))
		s.st.removeDoc(ctx, bdfKey(r.id))
		if err := s.st.deleteDone(ctx, r.id, r.rev); err != nil && !isNotFound(err) {
			return err
		}
	}
	return s.dropOrphans(ctx)
}

// removeDocuments removes the documents of a job done and records it.
func (s *Spool) removeDocuments(ctx context.Context, id int32) error {
	if err := s.st.removeDoc(ctx, originalKey(id)); err != nil {
		return err
	}
	if err := s.st.removeDoc(ctx, bdfKey(id)); err != nil {
		return err
	}
	j, rev, err := s.st.getDone(ctx, id)
	if isNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	j.DocumentsRemoved = true
	return s.st.updateDone(ctx, &j, rev)
}

// dropOrphans removes the documents of no job: a job not done has its
// documents, and a job done has them until they are removed.
func (s *Spool) dropOrphans(ctx context.Context) error {
	var orphans []string
	err := s.st.eachDoc(ctx, func(key string, _ int64) {
		s.mu.Lock()
		uploading := s.uploads[key]
		id := parseKey(strings.SplitN(key, ".", 2)[0])
		_, live := s.jobs[id]
		s.mu.Unlock()
		if uploading || live {
			return
		}
		if id != 0 {
			if j, _, err := s.st.getDone(ctx, id); err == nil && !j.DocumentsRemoved {
				return
			}
		}
		orphans = append(orphans, key)
	})
	if err != nil {
		return err
	}
	for _, key := range orphans {
		s.o.Log.Info("document of no job removed", "document", key)
		if err := s.st.removeDoc(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

// collectEvery is how often the runner collects.
const collectEvery = time.Hour
