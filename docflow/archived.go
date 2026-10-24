package docflow

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/yuruidb"
)

// What an archive wrote lately, to look at on the topology. Its store is
// listed by key alone, oldest first, so an archive keeps the latest
// entries it wrote in a record of its own (archive-states, by its name),
// newest last: the last recentArchived. An older one rolls off the
// record and stays in the store, where a downloader after the archive
// still reads it. Who sees an entry is who sees the job: its sender, the
// members of its group, the auditors and the admins.
//
//	GET /api/nodes/<archive>/archived                    the entries the viewer sees, newest first (?limit=)
//	GET /api/nodes/<archive>/archived/<key>/pages/<n>    a page of one (PNG; ?w=), drawn from the store
//	GET /api/nodes/<archive>/archived/<key>/document     the document as written (in the audit log)

const (
	recentArchived = 50
	defaultShown   = 20
)

// ArchivedEntry is a payload an archive wrote: the job as it was then,
// and the key of its objects in the archive's store (Object the payload,
// Key + ".json" the record).
type ArchivedEntry struct {
	Job      int32     `json:"job"`
	Name     string    `json:"name,omitempty"`
	Kind     string    `json:"kind"`
	Format   string    `json:"format,omitempty"`
	Pages    int       `json:"pages,omitempty"`
	Size     int64     `json:"size"`
	Received time.Time `json:"received"`
	Archived time.Time `json:"archived"`
	Key      string    `json:"key"`
	Object   string    `json:"object"`
	// SubmittedBy, ClaimedUser and Group are the job's: who sees it.
	SubmittedBy *Principal `json:"submitted_by,omitempty"`
	ClaimedUser string     `json:"claimed_user,omitempty"`
	Group       string     `json:"group,omitempty"`
}

// archiveState is the latest entries an archive wrote, oldest first.
type archiveState struct {
	Seq     int64           `json:"seq"`
	Entries []ArchivedEntry `json:"entries,omitempty"`
}

// noteArchived adds what an archive wrote to its latest entries; one
// written again (a retry) takes the place it had.
func (m *Manager) noteArchived(ctx context.Context, name string, e ArchivedEntry) error {
	_, err := change(ctx, m.st.db, name, true, func(st *archiveState) error {
		st.Entries = slices.DeleteFunc(st.Entries, func(x ArchivedEntry) bool { return x.Key == e.Key })
		st.Entries = append(st.Entries, e)
		if extra := len(st.Entries) - recentArchived; extra > 0 {
			st.Entries = slices.Delete(st.Entries, 0, extra)
		}
		st.Seq++
		return nil
	})
	return err
}

// archivedState is an archive's latest entries; none before it wrote.
func (m *Manager) archivedState(ctx context.Context, name string) archiveState {
	rec, err := m.st.db.Get[archiveState](ctx, name)
	if err != nil {
		return archiveState{}
	}
	return rec.Data
}

// dropArchived forgets an archive's latest entries (the archive is
// deleted); its store keeps what it wrote.
func (m *Manager) dropArchived(ctx context.Context, name string) error {
	rec, err := m.st.db.Get[archiveState](ctx, name)
	if errors.Is(err, yuruidb.ErrNotFound) {
		return nil
	}
	if err != nil {
		return mapErr(err)
	}
	if err := m.st.db.Delete[archiveState](ctx, name, rec.Revision()); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
		return mapErr(err)
	}
	return nil
}

// Archived returns the latest entries an archive wrote that the viewer
// may see, newest first, at most limit (defaultShown when 0).
func (m *Manager) Archived(ctx context.Context, name string, v *Viewer, limit int) ([]ArchivedEntry, error) {
	if _, err := m.nodeOfKind(ctx, name, KindArchive); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = defaultShown
	}
	st := m.archivedState(ctx, name)
	out := []ArchivedEntry{}
	for i := len(st.Entries) - 1; i >= 0 && len(out) < limit; i-- {
		if e := st.Entries[i]; v.seesSent(e.SubmittedBy, e.Group) {
			out = append(out, e)
		}
	}
	return out, nil
}

// archivedEntry finds an entry of an archive's latest by its key, of
// those the viewer may see; any other is not found.
func (m *Manager) archivedEntry(ctx context.Context, name, key string, v *Viewer) (*Node, *ArchivedEntry, error) {
	n, err := m.nodeOfKind(ctx, name, KindArchive)
	if err != nil {
		return nil, nil, err
	}
	st := m.archivedState(ctx, name)
	for i := range st.Entries {
		if e := &st.Entries[i]; e.Key == key && v.seesSent(e.SubmittedBy, e.Group) {
			return n, e, nil
		}
	}
	return nil, nil, fmt.Errorf("%w: %s lists no %q", ErrNotFound, name, key)
}

func (h *handler) listArchived(w http.ResponseWriter, r *http.Request, v *Viewer) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	es, err := h.m.Archived(r.Context(), r.PathValue("name"), v, min(limit, recentArchived))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, es)
}

// archivedPage is …/archived/<key>/pages/<n>.png?w=: a page of an
// entry's .bdf document, drawn from the archive's store, to look at
// before the download (which alone is in the audit log).
func (h *handler) archivedPage(w http.ResponseWriter, r *http.Request, v *Viewer) {
	a, e, err := h.m.archivedEntry(r.Context(), r.PathValue("name"), r.PathValue("key"), v)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	n, err := strconv.Atoi(strings.TrimSuffix(r.PathValue("page"), ".png"))
	if err != nil {
		h.writeErr(w, ErrNotFound)
		return
	}
	if e.Kind != PayloadBDF {
		h.writeErr(w, fmt.Errorf("%w: job %d is %s, not a .bdf document", ErrNotFound, e.Job, e.Kind))
		return
	}
	st, err := h.m.st.archive(r.Context(), a.Store, a.collection())
	if err != nil {
		h.writeErr(w, err)
		return
	}
	width, _ := strconv.Atoi(r.URL.Query().Get("w"))
	if width == 0 {
		width = 800
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "private, max-age=3600")
	if err := h.m.writePageIn(r.Context(), st, e.Object, e.Pages, n-1, width, w); err != nil {
		h.writeErr(w, err)
	}
}

// archivedDocument is …/archived/<key>/document: the document the
// archive wrote, as it was. The download is in the audit log.
func (h *handler) archivedDocument(w http.ResponseWriter, r *http.Request, v *Viewer) {
	a, e, err := h.m.archivedEntry(r.Context(), r.PathValue("name"), r.PathValue("key"), v)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	st, err := h.m.st.archive(r.Context(), a.Store, a.collection())
	if err != nil {
		h.writeErr(w, err)
		return
	}
	rd, err := st.Open(r.Context(), e.Object)
	if err != nil {
		h.writeErr(w, mapErr(err))
		return
	}
	defer rd.Close()
	ext := extensionOf(&Payload{Name: e.Name, Format: e.Format})
	name := e.Name
	if name == "" || strings.ContainsAny(name, `"\/`) {
		name = "job-" + jobID(e.Job) + ext
	} else if path.Ext(name) == "" {
		name += ext
	}
	h.m.audit(r.Context(), v.Label(), "archive-document", a.Name, fmt.Sprintf("job %d %s", e.Job, clip(e.Name, 100)))
	w.Header().Set("Content-Type", contentTypeOf(e.Format, e.Name))
	w.Header().Set("Content-Length", strconv.FormatInt(rd.Size, 10))
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("Content-Disposition", `attachment; filename="`+name+`"`)
	if _, err := io.Copy(w, rd); err != nil {
		panic(http.ErrAbortHandler)
	}
}
