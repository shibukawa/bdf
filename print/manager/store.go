package manager

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/shibukawa/yuruidb"
	"github.com/shibukawa/yuruidb/blob"
	ybolt "github.com/shibukawa/yuruidb/drivers/bbolt"
	_ "github.com/shibukawa/yuruidb/drivers/dynamodb"
	_ "github.com/shibukawa/yuruidb/drivers/firestore"
	"github.com/shibukawa/yuruidb/drivers/gcs"
	"github.com/shibukawa/yuruidb/drivers/r2"
	"github.com/shibukawa/yuruidb/drivers/s3"
)

// This file is the only one that knows which of yuruidb's backends keeps
// the records and the documents: the rest of the manager sees a
// yuruidb.Store and a blob.Store.

// Records of the collections. yuruidb binds a collection to a type, so
// each collection has a type of its own.
type (
	nodeRecord     Node
	templateRecord Template
	layoutRecord   Layout
	jobRecord      Job
	slipRecord     Slip
	itemRecord     Item
	deliveryRecord Delivery
	deviceRecord   Device
	auditRecord    AuditEntry
	enrollRecord   enrollment
	personRecord   Person
	counterRecord  struct {
		Next int32 `json:"next"`
	}
	// stateRecord is a value kept for a while under a key: see PutState.
	stateRecord struct {
		Payload []byte    `json:"payload"`
		Expires time.Time `json:"expires"`
	}
)

// Store keeps the records of a manager (nodes, jobs, slips, devices) and
// the documents of its jobs. Every manager of the same
// store sees the same state: a manager holds none of its own.
type Store struct {
	db   *yuruidb.Store
	docs *blob.Store
	// tpls keeps the templates' files, beside the documents.
	tpls *blob.Store
	// Shared tells a store that several processes can use at once (a
	// database or an object storage, not a file).
	Shared bool
	closed bool
	// archives are the stores of the archive nodes, opened as they are
	// first written to, by URL and collection.
	mu       sync.Mutex
	archives map[string]*blob.Store
}

// ErrStore is the error of a store that cannot keep a manager's records.
var ErrStore = errors.New("manager: store")

// OpenStore opens the store of the records and the store of the
// documents, both given as URLs:
//
//	bbolt:///var/lib/bdf/printmanager.db   one file, one process (records and documents)
//	dynamodb://table?region=ap-northeast-1 records only
//	firestore://project/database           records only
//	s3://bucket?region=…&endpoint=…        records or documents
//	r2://bucket?endpoint=…                 records or documents
//	gcs://bucket                           records or documents
//
// documents "" is records when that keeps documents too (bbolt, s3, r2,
// gcs). The records' store must make writes conditional (yuruidb's
// ConditionalMutations): a store that cannot is refused.
func OpenStore(ctx context.Context, records, documents string) (*Store, error) {
	if records == "" {
		return nil, fmt.Errorf("%w: no store of records", ErrStore)
	}
	if documents == "" {
		documents = records
	}
	db, err := openRecords(ctx, records)
	if err != nil {
		return nil, err
	}
	caps := db.Capabilities()
	if !caps.ConditionalMutations {
		db.Close()
		return nil, fmt.Errorf("%w: %s does not make writes conditional", ErrStore, scheme(records))
	}
	for _, err := range []error{
		db.Bind[nodeRecord]("nodes"),
		db.Bind[templateRecord]("templates"),
		db.Bind[layoutRecord]("layout"),
		db.Bind[locationState]("location-states"),
		db.Bind[inboxState]("inbox-states"),
		db.Bind[sinkState]("sink-states"),
		db.Bind[destState]("dest-states"),
		db.Bind[jobRecord]("jobs"),
		db.Bind[slipRecord]("slips"),
		db.Bind[itemRecord]("items"),
		db.Bind[deliveryRecord]("deliveries"),
		db.Bind[pendingRecord]("pending-deliveries"),
		db.Bind[dayRecord]("ended-days"),
		db.Bind[deviceRecord]("devices"),
		db.Bind[auditRecord]("audit"),
		db.Bind[enrollRecord]("enrollments"),
		db.Bind[counterRecord]("meta"),
		db.Bind[personRecord]("people"),
		db.Bind[stateRecord]("states"),
	} {
		if err != nil {
			db.Close()
			return nil, err
		}
	}
	docs, err := openDocuments(ctx, documents)
	if err != nil {
		db.Close()
		return nil, err
	}
	tpls, err := openBlobs(ctx, documents, templatesCollection)
	if err != nil {
		docs.Close()
		db.Close()
		return nil, err
	}
	return &Store{db: db, docs: docs, tpls: tpls, Shared: caps.Distributed, archives: map[string]*blob.Store{}}, nil
}

// Close closes the store.
func (s *Store) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	errs := []error{s.docs.Close(), s.tpls.Close(), s.db.Close()}
	s.mu.Lock()
	for _, a := range s.archives {
		errs = append(errs, a.Close())
	}
	s.archives = nil
	s.mu.Unlock()
	return errors.Join(errs...)
}

// archive opens the store of an archive node, or returns it opened.
func (s *Store) archive(ctx context.Context, raw, collection string) (*blob.Store, error) {
	key := raw + "|" + collection
	s.mu.Lock()
	defer s.mu.Unlock()
	if a := s.archives[key]; a != nil {
		return a, nil
	}
	if s.archives == nil {
		return nil, fmt.Errorf("%w: closed", ErrStore)
	}
	a, err := openBlobs(ctx, raw, collection)
	if err != nil {
		return nil, err
	}
	s.archives[key] = a
	return a, nil
}

// PutState keeps a value under a key until it is taken or expires: the
// state of a login on its way through the IdP, which any instance may
// finish. A key kept already is ErrConflict.
func (s *Store) PutState(ctx context.Context, key string, payload []byte, expires time.Time) error {
	_, err := s.db.Insert(ctx, stateID(key), stateRecord{Payload: payload, Expires: expires.UTC()})
	return mapErr(err)
}

// TakeState returns the value of a key and removes it, so that it is
// taken once: ErrNotFound when it is not there, was taken meanwhile, or
// expired before now.
func (s *Store) TakeState(ctx context.Context, key string, now time.Time) ([]byte, error) {
	id := stateID(key)
	rec, err := s.db.Get[stateRecord](ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	if err := s.db.Delete[stateRecord](ctx, id, rec.Revision()); err != nil {
		if errors.Is(err, yuruidb.ErrConflict) || errors.Is(err, yuruidb.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, mapErr(err)
	}
	if !now.Before(rec.Data.Expires) {
		return nil, ErrNotFound
	}
	return rec.Data.Payload, nil
}

// PruneStates removes up to limit values that expired before a time, and
// tells how many it removed.
func (s *Store) PruneStates(ctx context.Context, before time.Time, limit int) (int, error) {
	recs, err := list[stateRecord](ctx, s.db)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, r := range recs {
		if n >= limit {
			break
		}
		if r.Data.Expires.Before(before) && s.db.Delete[stateRecord](ctx, r.ID(), r.Revision()) == nil {
			n++
		}
	}
	return n, nil
}

// stateID is the record id of a key: a digest, which every backend takes
// whatever the key's characters.
func stateID(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

func scheme(u string) string {
	s, _, _ := strings.Cut(u, ":")
	return s
}

// boltPath is the file of a bbolt URL: bbolt:///abs/path, bbolt:rel/path
// or bbolt://rel/path.
func boltPath(u *url.URL) string {
	p := u.Opaque
	if p == "" {
		p = u.Host + u.Path
	}
	return filepath.FromSlash(p)
}

func openRecords(ctx context.Context, raw string) (*yuruidb.Store, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStore, err)
	}
	q := u.Query()
	cfg := yuruidb.Config{}
	set := func(key, value string) {
		if value != "" {
			cfg[key] = value
		}
	}
	switch u.Scheme {
	case "bbolt":
		set("Path", boltPath(u))
	case "dynamodb":
		set("Table", u.Host)
		set("Region", q.Get("region"))
		set("Endpoint", q.Get("endpoint"))
	case "firestore":
		set("Project", u.Host)
		set("Database", strings.Trim(u.Path, "/"))
		set("Endpoint", q.Get("endpoint"))
	case "s3", "r2":
		set("Bucket", u.Host)
		set("Region", q.Get("region"))
		set("Endpoint", q.Get("endpoint"))
		set("AddressStyle", q.Get("address-style"))
	case "gcs":
		set("Bucket", u.Host)
		set("Endpoint", q.Get("endpoint"))
	default:
		return nil, fmt.Errorf("%w: records in %q: bbolt, dynamodb, firestore, s3, r2 or gcs", ErrStore, u.Scheme)
	}
	db, err := yuruidb.Open(ctx, u.Scheme, cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %w", ErrStore, u.Scheme, err)
	}
	return db, nil
}

const (
	docsCollection      = "job-documents"
	templatesCollection = "templates"
)

func openDocuments(ctx context.Context, raw string) (*blob.Store, error) {
	return openBlobs(ctx, raw, docsCollection)
}

// openBlobs opens a store of blobs: the documents of the jobs, or those
// of an archive.
func openBlobs(ctx context.Context, raw, collection string) (*blob.Store, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStore, err)
	}
	q := u.Query()
	var docs *blob.Store
	switch u.Scheme {
	case "bbolt":
		docs, err = ybolt.OpenBlobs(ctx, ybolt.Config{Path: boltPath(u), Collection: collection})
	case "s3":
		docs, err = s3.OpenBlobs(ctx, s3.Config{Bucket: u.Host, Region: q.Get("region"), Endpoint: q.Get("endpoint"),
			AddressStyle: q.Get("address-style"), Collection: collection})
	case "r2":
		docs, err = r2.OpenBlobs(ctx, r2.Config{Bucket: u.Host, Endpoint: q.Get("endpoint"), Collection: collection})
	case "gcs":
		docs, err = gcs.OpenBlobs(ctx, gcs.Config{Bucket: u.Host, Endpoint: q.Get("endpoint"), Collection: collection})
	default:
		return nil, fmt.Errorf("%w: documents in %q: bbolt, s3, r2 or gcs", ErrStore, u.Scheme)
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %w", ErrStore, u.Scheme, err)
	}
	return docs, nil
}
