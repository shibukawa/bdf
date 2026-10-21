// Package manager is the BDF Print Manager's core: a topology of nodes
// a job flows through, as an administrator composes it. Endpoints take
// documents and data from applications; queues forward to several nodes
// at once, classes to one of their members in turn, inboxes hold jobs
// until a person sends them on; displays show jobs as slips on the
// tablets enrolled there (a warehouse) instead of printing them, and
// print on paper what nobody opened; printers print; archives keep a copy
// that is never changed. Every node is an IPP printer at
// /ipp/print/<name>.
//
// A Manager holds no state of its own: everything is in its Store, so
// that any request may reach any instance (a container, several, or a
// function of AWS Lambda). Nothing runs in the background either: what is
// due is done by Tick, which a scheduler calls (or a timer of the
// process). The web application is cmd/printmanager; Handler is its API,
// its IPP printers and the devices' endpoints.
package manager

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/shibukawa/bdf/print/spool"
	"github.com/shibukawa/yuruidb"
)

// Options configures a manager.
type Options struct {
	// Converter makes .bdf documents of what senders send (PDF, Office,
	// rasters …); nil takes .bdf only.
	Converter spool.Converter
	// Secret signs the challenges and sessions of devices: at least 32
	// random bytes, the same for every instance.
	Secret []byte
	// FontDirs and NoSystemFonts find the fonts of documents, to draw
	// their pages and print them.
	FontDirs      []string
	NoSystemFonts bool
	// AdminToken authorizes the API (bearer); SubmitTokens are the names
	// and secrets that may send jobs (IPP Basic, or bearer name:secret).
	// With no SubmitTokens anyone may send.
	AdminToken   string
	SubmitTokens map[string]string
	// TickToken authorizes POST /api/tick alone (a scheduler's).
	TickToken string
	// Viewer returns the person signed in to the pages a request comes
	// from (the web application's session), or nil. An admin uses the
	// API as the AdminToken does; anyone sees the jobs of their scope and
	// sends jobs to the nodes they may send to.
	Viewer func(*http.Request) *Viewer
	// MaxDocument bounds a document received (default 64 MiB).
	MaxDocument int64
	// SessionFor is how long a device's session lasts before it signs a
	// challenge again (default 10 minutes); SecureCookies sets Secure.
	SessionFor    time.Duration
	SecureCookies bool
	// TabletPath is the devices' page, which an enrollment's URL opens
	// with ?code= (default /tablet).
	TabletPath string
	// JobKeep is how long a job that ended stays (its record and, once
	// nothing shows it any more, its document): 14 days by default.
	JobKeep time.Duration
	// Background sends the deliveries of a job (to printers and archives)
	// right after it is routed, in a goroutine of the process; without
	// it the next Tick does (a function of a serverless platform).
	Background bool
	// TempDir holds documents while they are converted ("" the system's).
	TempDir string
	Log     *slog.Logger
	// Now is the clock (default time.Now).
	Now func() time.Time
}

// Manager is a view of a Store. It is cheap to make: an instance may make
// one per request.
type Manager struct {
	st *Store
	o  Options
}

// New makes a manager of a store.
func New(st *Store, o Options) (*Manager, error) {
	if len(o.Secret) < 32 {
		return nil, fmt.Errorf("%w: the secret has %d bytes; it needs 32", ErrInvalid, len(o.Secret))
	}
	if o.MaxDocument <= 0 {
		o.MaxDocument = 64 << 20
	}
	if o.SessionFor <= 0 {
		o.SessionFor = 10 * time.Minute
	}
	if o.TabletPath == "" {
		o.TabletPath = "/tablet"
	}
	if o.JobKeep <= 0 {
		o.JobKeep = 14 * 24 * time.Hour
	}
	if o.Log == nil {
		o.Log = slog.New(slog.DiscardHandler)
	}
	if o.Now == nil {
		o.Now = time.Now
	}
	return &Manager{st: st, o: o}, nil
}

func (m *Manager) now() time.Time { return m.o.Now().UTC().Truncate(time.Millisecond) }

// mapErr turns yuruidb's errors into the manager's.
func mapErr(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, yuruidb.ErrNotFound):
		return fmt.Errorf("%w: %w", ErrNotFound, err)
	case errors.Is(err, yuruidb.ErrConflict), errors.Is(err, yuruidb.ErrAlreadyExists):
		return fmt.Errorf("%w: %w", ErrConflict, err)
	}
	return err
}

// tries is how many times a read-modify-write is tried against others
// writing the same record.
const tries = 16

// change reads a record, changes it and writes it back, again while
// another write comes between; a record not there is made when create is
// set (fn sees the zero value). fn returning errSkip leaves the record as
// it is.
func change[T any](ctx context.Context, db *yuruidb.Store, id string, create bool, fn func(*T) error) (T, error) {
	var zero T
	for range tries {
		rec, err := db.Get[T](ctx, id)
		missing := errors.Is(err, yuruidb.ErrNotFound)
		if err != nil && !(missing && create) {
			return zero, mapErr(err)
		}
		v := rec.Data
		if missing {
			v = zero
		}
		if err := fn(&v); err != nil {
			if errors.Is(err, errSkip) {
				return v, nil
			}
			return zero, err
		}
		if missing {
			_, err = db.Insert(ctx, id, v)
		} else {
			_, err = db.Update(ctx, id, v, rec.Revision())
		}
		if errors.Is(err, yuruidb.ErrConflict) || errors.Is(err, yuruidb.ErrAlreadyExists) {
			continue
		}
		return v, mapErr(err)
	}
	return zero, ErrConflict
}

var errSkip = errors.New("manager: no change")

// yuruidbListOptions is a page of a listing.
func yuruidbListOptions(limit int, cursor string) yuruidb.ListOptions {
	return yuruidb.ListOptions{Limit: limit, Cursor: cursor}
}

// list reads every record of a collection.
func list[T any](ctx context.Context, db *yuruidb.Store) ([]yuruidb.Record[T], error) {
	var all []yuruidb.Record[T]
	cursor := ""
	for {
		page, err := db.List[T](ctx, yuruidb.ListOptions{Limit: 200, Cursor: cursor})
		if err != nil {
			return nil, mapErr(err)
		}
		all = append(all, page.Records...)
		if page.Cursor == "" {
			return all, nil
		}
		cursor = page.Cursor
	}
}

func randomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// display returns a display node.
func (m *Manager) display(ctx context.Context, name string) (*Node, error) {
	return m.nodeOfKind(ctx, name, KindDisplay)
}

// printer returns a printer node.
func (m *Manager) printer(ctx context.Context, name string) (*Node, error) {
	return m.nodeOfKind(ctx, name, KindPrinter)
}

// nodeOfKind returns a node of a kind; one of another kind is not found.
func (m *Manager) nodeOfKind(ctx context.Context, name, kind string) (*Node, error) {
	n, err := m.Node(ctx, name)
	if err != nil {
		return nil, err
	}
	if n.Kind != kind {
		return nil, fmt.Errorf("%w: %s is a %s, not a %s", ErrNotFound, name, n.Kind, kind)
	}
	return n, nil
}

// nodesOfKind returns the nodes of a kind, by name.
func (m *Manager) nodesOfKind(ctx context.Context, kind string) ([]Node, error) {
	nodes, err := m.Nodes(ctx)
	if err != nil {
		return nil, err
	}
	out := nodes[:0]
	for _, n := range nodes {
		if n.Kind == kind {
			out = append(out, n)
		}
	}
	return out, nil
}

// touch changes a location's state, and its Seq so that its revision
// changes.
func (m *Manager) touch(ctx context.Context, loc string, fn func(*locationState)) error {
	_, err := change(ctx, m.st.db, loc, true, func(s *locationState) error {
		fn(s)
		s.Seq++
		return nil
	})
	return err
}

// nextJob takes the next job number, from 1.
func (m *Manager) nextJob(ctx context.Context) (int32, error) {
	var job int32
	_, err := change(ctx, m.st.db, "next-job", true, func(c *counterRecord) error {
		if c.Next <= 0 || c.Next == 1<<31-1 {
			c.Next = 1
		}
		job = c.Next
		c.Next++
		return nil
	})
	return job, err
}

// Slip returns a slip.
func (m *Manager) Slip(ctx context.Context, id string) (*Slip, error) {
	rec, err := m.st.db.Get[slipRecord](ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	s := Slip(rec.Data)
	return &s, nil
}

// SlipOfJob returns the first slip of a job number: the one at the first
// location the job reached.
func (m *Manager) SlipOfJob(ctx context.Context, job int32) (*Slip, error) {
	j, err := m.JobOfNumber(ctx, job)
	if err != nil {
		return nil, err
	}
	for _, l := range j.Legs {
		if l.Kind == KindDisplay && l.Ref != "" {
			return m.Slip(ctx, l.Ref)
		}
	}
	return nil, fmt.Errorf("%w: job %d has no slip", ErrNotFound, job)
}

// changeSlip changes a slip; a change of its state touches its location.
func (m *Manager) changeSlip(ctx context.Context, id string, fn func(*Slip) error) (*Slip, error) {
	r, err := change(ctx, m.st.db, id, false, func(r *slipRecord) error {
		return fn((*Slip)(r))
	})
	if err != nil {
		return nil, err
	}
	s := Slip(r)
	return &s, nil
}

// Board is what a location's devices show: the open slips, oldest first,
// and the ended ones still shown, latest first. Rev changes whenever any
// of them does.
type Board struct {
	Location string `json:"location"`
	Title    string `json:"title"`
	Rev      string `json:"rev"`
	Open     []Slip `json:"open"`
	Ended    []Slip `json:"ended"`
}

// BoardRev returns the revision of a location's board, which a device
// compares with the one it has before reading the board again.
func (m *Manager) BoardRev(ctx context.Context, loc string) (string, error) {
	rec, err := m.st.db.Get[locationState](ctx, loc)
	if errors.Is(err, yuruidb.ErrNotFound) {
		return "0", nil
	}
	if err != nil {
		return "", mapErr(err)
	}
	return fmt.Sprintf("%d", rec.Data.Seq), nil
}

// Board returns a display's board.
func (m *Manager) Board(ctx context.Context, loc string) (*Board, error) {
	l, err := m.display(ctx, loc)
	if err != nil {
		return nil, err
	}
	b := &Board{Location: loc, Title: or(l.Title, l.Name), Rev: "0", Open: []Slip{}, Ended: []Slip{}}
	rec, err := m.st.db.Get[locationState](ctx, loc)
	if errors.Is(err, yuruidb.ErrNotFound) {
		return b, nil
	}
	if err != nil {
		return nil, mapErr(err)
	}
	b.Rev = fmt.Sprintf("%d", rec.Data.Seq)
	read := func(ids []string) ([]Slip, error) {
		out := make([]Slip, 0, len(ids))
		for _, id := range ids {
			s, err := m.Slip(ctx, id)
			if errors.Is(err, ErrNotFound) {
				continue // removed meanwhile
			}
			if err != nil {
				return nil, err
			}
			out = append(out, *s)
		}
		return out, nil
	}
	if b.Open, err = read(rec.Data.Open); err != nil {
		return nil, err
	}
	ended := slices.Clone(rec.Data.Ended)
	slices.Reverse(ended)
	if b.Ended, err = read(ended); err != nil {
		return nil, err
	}
	return b, nil
}

// end moves a slip from a location's open slips to its ended ones.
func (m *Manager) end(ctx context.Context, s *Slip) error {
	return m.touch(ctx, s.Location, func(st *locationState) {
		st.Open = slices.DeleteFunc(st.Open, func(id string) bool { return id == s.ID })
		if !slices.Contains(st.Ended, s.ID) {
			st.Ended = append(st.Ended, s.ID)
		}
	})
}

// Open marks a slip opened by a device; a slip opened before stays as it
// was.
func (m *Manager) Open(ctx context.Context, id string, by Mark) (*Slip, error) {
	by.At = m.now()
	changed := false
	s, err := m.changeSlip(ctx, id, func(s *Slip) error {
		if s.State != SlipNew {
			return errSkip
		}
		s.State, s.Opened, changed = SlipOpened, &by, true
		return nil
	})
	if err != nil || !changed {
		return s, err
	}
	return s, m.touch(ctx, s.Location, func(*locationState) {})
}

// Done marks a slip done.
func (m *Manager) Done(ctx context.Context, id string, by Mark) (*Slip, error) {
	by.At = m.now()
	s, err := m.changeSlip(ctx, id, func(s *Slip) error {
		if s.State.ended() {
			return fmt.Errorf("%w: the slip is %s", ErrConflict, s.State)
		}
		if s.Opened == nil {
			s.Opened = &by
		}
		s.State, s.Done = SlipDone, &by
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s, m.end(ctx, s)
}

// CancelSlip cancels a slip that is not done: it leaves the devices, with
// a mark that it was canceled. The job goes on elsewhere; Cancel cancels
// the job.
func (m *Manager) CancelSlip(ctx context.Context, id string, by Mark) (*Slip, error) {
	by.At = m.now()
	s, err := m.changeSlip(ctx, id, func(s *Slip) error {
		if s.State.ended() {
			return fmt.Errorf("%w: the slip is %s", ErrConflict, s.State)
		}
		s.State, s.Canceled = SlipCanceled, &by
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s, m.end(ctx, s)
}

// removeSlip deletes an ended slip; the job's document stays while
// something else needs it (release).
func (m *Manager) removeSlip(ctx context.Context, s *Slip) error {
	rec, err := m.st.db.Get[slipRecord](ctx, s.ID)
	switch {
	case err == nil:
		if err := m.st.db.Delete[slipRecord](ctx, s.ID, rec.Revision()); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
			return mapErr(err)
		}
		m.release(ctx, s.Job)
	case !errors.Is(err, yuruidb.ErrNotFound):
		return mapErr(err)
	}
	return m.touch(ctx, s.Location, func(st *locationState) {
		st.Ended = slices.DeleteFunc(st.Ended, func(id string) bool { return id == s.ID })
		st.Open = slices.DeleteFunc(st.Open, func(id string) bool { return id == s.ID })
	})
}

// discard drains and closes what a sender sent past what is read.
func discard(r io.Reader) { io.Copy(io.Discard, io.LimitReader(r, 1<<20)) }
