package manager

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

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
	locationRecord Location
	printerRecord  Printer
	slipRecord     Slip
	deviceRecord   Device
	enrollRecord   enrollment
	counterRecord  struct {
		Next int32 `json:"next"`
	}
)

// Store keeps the records of a manager (locations, printers, slips,
// devices) and the documents of its slips. Every manager of the same
// store sees the same state: a manager holds none of its own.
type Store struct {
	db   *yuruidb.Store
	docs *blob.Store
	// Shared tells a store that several processes can use at once (a
	// database or an object storage, not a file).
	Shared bool
	closed bool
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
		db.Bind[locationRecord]("locations"),
		db.Bind[printerRecord]("printers"),
		db.Bind[locationState]("location-states"),
		db.Bind[slipRecord]("slips"),
		db.Bind[deviceRecord]("devices"),
		db.Bind[enrollRecord]("enrollments"),
		db.Bind[counterRecord]("meta"),
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
	return &Store{db: db, docs: docs, Shared: caps.Distributed}, nil
}

// Close closes the store.
func (s *Store) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	return errors.Join(s.docs.Close(), s.db.Close())
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

const docsCollection = "slip-documents"

func openDocuments(ctx context.Context, raw string) (*blob.Store, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStore, err)
	}
	q := u.Query()
	var docs *blob.Store
	switch u.Scheme {
	case "bbolt":
		docs, err = ybolt.OpenBlobs(ctx, ybolt.Config{Path: boltPath(u), Collection: docsCollection})
	case "s3":
		docs, err = s3.OpenBlobs(ctx, s3.Config{Bucket: u.Host, Region: q.Get("region"), Endpoint: q.Get("endpoint"),
			AddressStyle: q.Get("address-style"), Collection: docsCollection})
	case "r2":
		docs, err = r2.OpenBlobs(ctx, r2.Config{Bucket: u.Host, Endpoint: q.Get("endpoint"), Collection: docsCollection})
	case "gcs":
		docs, err = gcs.OpenBlobs(ctx, gcs.Config{Bucket: u.Host, Endpoint: q.Get("endpoint"), Collection: docsCollection})
	default:
		return nil, fmt.Errorf("%w: documents in %q: bbolt, s3, r2 or gcs", ErrStore, u.Scheme)
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %w", ErrStore, u.Scheme, err)
	}
	return docs, nil
}
