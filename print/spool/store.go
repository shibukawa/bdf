package spool

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/shibukawa/yuruidb"
	"github.com/shibukawa/yuruidb/blob"
	ybolt "github.com/shibukawa/yuruidb/drivers/bbolt"
)

// This file is the only one that knows yuruidb: the spool's records are
// collections of one bbolt file and its documents are blobs of the same
// file. Another of yuruidb's backends would change only this file.

// Records of the collections. Each collection has its own Go type, as
// yuruidb binds a collection to a type.
type (
	liveJob     Job // jobs: the jobs not done
	finishedJob Job // done: the jobs done, which are also the accounting
	counter     struct {
		Next int32 `json:"next"`
	}
	queueRecord struct {
		Paused    bool `json:"paused,omitempty"`
		Rejecting bool `json:"rejecting,omitempty"`
	}
	trustRecord struct {
		Host        string `json:"host"`
		Fingerprint string `json:"fingerprint"`
	}
)

const counterKey = "next-id"

// store keeps the spool's records and documents.
type store struct {
	db   *yuruidb.Store
	docs *blob.Store
}

func openStore(ctx context.Context, path string, timeout time.Duration) (*store, error) {
	db, err := ybolt.Open(ctx, ybolt.Config{Path: path, Timeout: timeout})
	if err != nil {
		return nil, fmt.Errorf("spool: %s: %w", path, err)
	}
	for _, err := range []error{
		db.Bind[liveJob]("jobs"),
		db.Bind[finishedJob]("done"),
		db.Bind[counter]("meta"),
		db.Bind[queueRecord]("queues"),
		db.Bind[trustRecord]("trust"),
	} {
		if err != nil {
			db.Close()
			return nil, err
		}
	}
	docs, err := ybolt.OpenBlobs(ctx, ybolt.Config{Path: path, Timeout: timeout, Collection: "docs"})
	if err != nil {
		db.Close()
		return nil, err
	}
	return &store{db: db, docs: docs}, nil
}

func (s *store) Close() error {
	return errors.Join(s.docs.Close(), s.db.Close())
}

// stored is a record and the revision it was read at.
type stored struct {
	job Job
	rev string
}

// liveJobs reads all the jobs not done, in the order of their ids.
func (s *store) liveJobs(ctx context.Context) ([]stored, error) {
	var out []stored
	cursor := ""
	for {
		page, err := s.db.List[liveJob](ctx, yuruidb.ListOptions{Limit: 1000, Cursor: cursor})
		if err != nil {
			return nil, err
		}
		for _, r := range page.Records {
			out = append(out, stored{Job(r.Data), r.Revision()})
		}
		if page.Cursor == "" {
			return out, nil
		}
		cursor = page.Cursor
	}
}

func (s *store) insertLive(ctx context.Context, j *Job) (string, error) {
	r, err := s.db.Insert(ctx, j.key(), liveJob(*j))
	return r.Revision(), err
}

func (s *store) updateLive(ctx context.Context, j *Job, rev string) (string, error) {
	r, err := s.db.Update(ctx, j.key(), liveJob(*j), rev)
	return r.Revision(), err
}

func (s *store) deleteLive(ctx context.Context, id int32, rev string) error {
	return s.db.Delete[liveJob](ctx, jobKey(id), rev)
}

// insertDone records a job done. Recording it twice (after a crash between
// recording it and removing it from the jobs) is not an error.
func (s *store) insertDone(ctx context.Context, j *Job) error {
	_, err := s.db.Insert(ctx, j.key(), finishedJob(*j))
	if errors.Is(err, yuruidb.ErrAlreadyExists) {
		return nil
	}
	return err
}

func (s *store) getDone(ctx context.Context, id int32) (Job, string, error) {
	r, err := s.db.Get[finishedJob](ctx, jobKey(id))
	if err != nil {
		return Job{}, "", err
	}
	return Job(r.Data), r.Revision(), nil
}

func (s *store) updateDone(ctx context.Context, j *Job, rev string) error {
	_, err := s.db.Update(ctx, j.key(), finishedJob(*j), rev)
	return err
}

func (s *store) deleteDone(ctx context.Context, id int32, rev string) error {
	return s.db.Delete[finishedJob](ctx, jobKey(id), rev)
}

// eachDone calls f with the jobs done, in the order of their ids, until f
// returns false.
func (s *store) eachDone(ctx context.Context, f func(Job, string) bool) error {
	cursor := ""
	for {
		page, err := s.db.List[finishedJob](ctx, yuruidb.ListOptions{Limit: 1000, Cursor: cursor})
		if err != nil {
			return err
		}
		for _, r := range page.Records {
			if !f(Job(r.Data), r.Revision()) {
				return nil
			}
		}
		if page.Cursor == "" {
			return nil
		}
		cursor = page.Cursor
	}
}

// nextID takes the next job id. An id is never given twice; one may be
// skipped (when the job is not made after all).
func (s *store) nextID(ctx context.Context) (int32, error) {
	// (submissions at once take ids at once: the one that loses the
	// update tries again)
	for range 100 {
		id, err := s.takeID(ctx)
		if !isConflict(err) {
			return id, err
		}
	}
	return 0, errors.New("spool: the job ids are contended")
}

func (s *store) takeID(ctx context.Context) (int32, error) {
	r, err := s.db.Get[counter](ctx, counterKey)
	if errors.Is(err, yuruidb.ErrNotFound) {
		if _, err := s.db.Insert(ctx, counterKey, counter{Next: 2}); err != nil {
			return 0, err
		}
		return 1, nil
	}
	if err != nil {
		return 0, err
	}
	id := r.Data.Next
	if id <= 0 || id == 1<<31-1 {
		return 0, errors.New("spool: job ids are exhausted")
	}
	if _, err := s.db.Update(ctx, counterKey, counter{Next: id + 1}, r.Revision()); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *store) queue(ctx context.Context, name string) (queueRecord, string, error) {
	r, err := s.db.Get[queueRecord](ctx, nameKey(name))
	if errors.Is(err, yuruidb.ErrNotFound) {
		return queueRecord{}, "", nil
	}
	return r.Data, r.Revision(), err
}

func (s *store) setQueue(ctx context.Context, name string, q queueRecord, rev string) error {
	var err error
	if rev == "" {
		_, err = s.db.Insert(ctx, nameKey(name), q)
	} else {
		_, err = s.db.Update(ctx, nameKey(name), q, rev)
	}
	return err
}

// putDoc stores a document; getDoc opens it for reading at any offset, as
// long as ctx lasts.
func (s *store) putDoc(ctx context.Context, key string, r io.Reader) (int64, error) {
	m, err := s.docs.Insert(ctx, key, r)
	return m.Size, err
}

func (s *store) getDoc(ctx context.Context, key string) (*blob.Reader, error) {
	return s.docs.Open(ctx, key)
}

// removeDoc removes a document; one already gone is not an error.
func (s *store) removeDoc(ctx context.Context, key string) error {
	err := s.docs.Remove(ctx, key)
	if errors.Is(err, yuruidb.ErrNotFound) {
		return nil
	}
	return err
}

// eachDoc calls f with the key and size of each document.
func (s *store) eachDoc(ctx context.Context, f func(key string, size int64)) error {
	cursor := ""
	for {
		page, err := s.docs.List(ctx, yuruidb.ListOptions{Limit: 1000, Cursor: cursor})
		if err != nil {
			return err
		}
		for _, b := range page.Blobs {
			f(b.ID, b.Size)
		}
		if page.Cursor == "" {
			return nil
		}
		cursor = page.Cursor
	}
}

// backup writes a consistent copy of the whole store.
func (s *store) backup(ctx context.Context, w io.Writer) error {
	return s.db.Backup(ctx, w)
}

// trustBackend keeps the fingerprints of the printers' certificates in the
// store, for an ipp.TrustStore.
type trustBackend struct{ s *store }

func (t trustBackend) Load() (map[string]string, error) {
	ctx := context.Background()
	out := map[string]string{}
	cursor := ""
	for {
		page, err := t.s.db.List[trustRecord](ctx, yuruidb.ListOptions{Limit: 1000, Cursor: cursor})
		if err != nil {
			return nil, err
		}
		for _, r := range page.Records {
			out[r.Data.Host] = r.Data.Fingerprint
		}
		if page.Cursor == "" {
			return out, nil
		}
		cursor = page.Cursor
	}
}

// Save writes the fingerprints that changed: hosts are only added or
// forgotten one at a time.
func (t trustBackend) Save(m map[string]string) error {
	ctx := context.Background()
	old, err := t.Load()
	if err != nil {
		return err
	}
	for host, fp := range m {
		if old[host] == fp {
			continue
		}
		r, err := t.s.db.Get[trustRecord](ctx, nameKey(host))
		switch {
		case errors.Is(err, yuruidb.ErrNotFound):
			_, err = t.s.db.Insert(ctx, nameKey(host), trustRecord{host, fp})
		case err == nil:
			_, err = t.s.db.Update(ctx, nameKey(host), trustRecord{host, fp}, r.Revision())
		}
		if err != nil {
			return err
		}
	}
	for host := range old {
		if _, ok := m[host]; ok {
			continue
		}
		r, err := t.s.db.Get[trustRecord](ctx, nameKey(host))
		if err == nil {
			err = t.s.db.Delete[trustRecord](ctx, nameKey(host), r.Revision())
		}
		if err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
			return err
		}
	}
	return nil
}

// nameKey makes a name (a queue's, a host) a key of the store: the name
// itself when it is one of the store's IDs (letters, digits, '.', '_' and
// '-'), else "_x" and the hex of its bytes.
func nameKey(name string) string {
	ok := name != "" && len(name) <= 200 && name != "." && name != ".." && !(len(name) >= 2 && name[0] == '_' && name[1] == 'x')
	for _, c := range name {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '.' || c == '_' || c == '-') {
			ok = false
			break
		}
	}
	if len(name) >= 4 && name[:2] == "__" && name[len(name)-2:] == "__" {
		ok = false
	}
	if ok {
		return name
	}
	return "_x" + hex.EncodeToString([]byte(name))
}

// isConflict tells an error of a record changed since it was read: the
// spool is the store's only writer, so it is a bug or another process.
func isConflict(err error) bool { return errors.Is(err, yuruidb.ErrConflict) }

// isNotFound tells an error of a record or document that is not there.
func isNotFound(err error) bool { return errors.Is(err, yuruidb.ErrNotFound) }
