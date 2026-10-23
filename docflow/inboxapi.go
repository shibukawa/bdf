package docflow

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/raster/imagebdf"
	"github.com/shibukawa/bdf/thumbnail"
)

// The inboxes' API: what an inbox holds is worked at /api/inbox/<name>,
// as on the Inboxes page, by applications as well as people:
//
//	GET  /api/inbox                          the inboxes the viewer works, with what each holds
//	GET  /api/inbox/<name>                   the items held (oldest first), the ended ones, and the nodes it releases to
//	GET  /api/inbox/<name>/<job>             an item, by its job's number
//	GET  /api/inbox/<name>/<job>/thumbnail   its first page, small (PNG; ?size=, ?mode=)
//	GET  /api/inbox/<name>/<job>/pages/<n>   a page of it (PNG; ?w=)
//	GET  /api/inbox/<name>/<job>/document    the document it holds
//	POST /api/inbox/<name>/<job>/release     sends it on: {"to": [nodes], "worker": name}
//	POST /api/inbox/<name>/<job>/discard     drops it: {"worker": name}
//
// Who works an inbox: the role topology, every one; an app the inbox
// lists among its Apps, with the scope inbox; and the inbox's people, by
// its scope: the members of a group's, the owner of a personal one (a
// shared inbox has none: the role topology works it). A person works it
// with the pages' session, or with a token a frontend got for them, of
// the scope inbox. All but the role topology release to the inbox's To
// alone, so that an inbox that names none is released by the role
// topology. The thumbnail and the pages are of a .bdf document: an item of
// a file as received (an uploader whose output is file) or of data has
// none, only its document.

// mayWork tells why the viewer may not work an inbox (ErrForbidden), or
// nil when they may.
func (v *Viewer) mayWork(n *Node) error {
	switch {
	case v == nil:
		return fmt.Errorf("%w: a token or a sign-in is required", ErrUnauthorized)
	case v.Can(RoleTopology):
		return nil
	case v.Kind == PrincipalApp:
		if !slices.Contains(n.Apps, v.ID) {
			return fmt.Errorf("%w: %s does not take the app %s", ErrForbidden, n.Name, v.ID)
		}
		if !v.Has(ScopeInbox) {
			return fmt.Errorf("%w: the app %s has no scope %s", ErrForbidden, v.ID, ScopeInbox)
		}
		return nil
	case v.Kind == PrincipalUser:
		if v.Token() && !v.Has(ScopeInbox) {
			return fmt.Errorf("%w: the token has no scope %s", ErrForbidden, ScopeInbox)
		}
		if n.Scope == ScopeGroup && v.Member(n.Owner) || n.Scope == ScopePersonal && v.is(n.Owner) {
			return nil
		}
	}
	return fmt.Errorf("%w: you do not work the inbox %s", ErrForbidden, n.Name)
}

// InboxSummary is an inbox as GET /api/inbox lists it: the nodes it
// releases to, and how many items it holds.
type InboxSummary struct {
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	To          []string `json:"to"`
	Held        int      `json:"held"`
	Rev         string   `json:"rev"`
}

// Inboxes lists the inboxes a viewer works, by name.
func (m *Manager) Inboxes(ctx context.Context, v *Viewer) ([]InboxSummary, error) {
	nodes, err := m.nodesOfKind(ctx, KindInbox)
	if err != nil {
		return nil, err
	}
	out := []InboxSummary{}
	for i := range nodes {
		n := &nodes[i]
		if v.mayWork(n) != nil {
			continue
		}
		_, release := m.pullsOf(ctx, n)
		s := InboxSummary{Name: n.Name, Title: or(n.Title, n.Name), Description: n.Description, To: release, Rev: "0"}
		if s.To == nil {
			s.To = []string{}
		}
		if st, err := m.st.db.Get[inboxState](ctx, n.Name); err == nil {
			s.Held, s.Rev = len(st.Data.Held), strconv.FormatInt(st.Data.Seq, 10)
		}
		out = append(out, s)
	}
	slices.SortFunc(out, func(a, b InboxSummary) int { return strings.Compare(a.Name, b.Name) })
	return out, nil
}

func (h *handler) listInboxes(w http.ResponseWriter, r *http.Request) {
	v, ok := h.viewerOf(r)
	if !ok {
		w.Header().Set("WWW-Authenticate", `Bearer realm="bdf docflow inbox"`)
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "a token or a sign-in is required"})
		return
	}
	out, err := h.m.Inboxes(r.Context(), v)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

// workedInbox finds the inbox a request names and who works it; a viewer
// who does not is refused.
func (h *handler) workedInbox(w http.ResponseWriter, r *http.Request) (*Node, *Viewer, bool) {
	v, ok := h.viewerOf(r)
	switch {
	case !ok:
		w.Header().Set("WWW-Authenticate", `Bearer realm="bdf docflow inbox"`)
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "a token or a sign-in is required"})
		return nil, nil, false
	case !sameOrigin(r):
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "a request of another origin"})
		return nil, nil, false
	}
	n, err := h.m.nodeOfKind(r.Context(), r.PathValue("node"), KindInbox)
	if err == nil {
		err = v.mayWork(n)
	}
	if err != nil {
		h.writeErr(w, err)
		return nil, nil, false
	}
	return n, v, true
}

// workedItem finds the item of the job a request names at its inbox.
func (h *handler) workedItem(w http.ResponseWriter, r *http.Request) (*Node, *Viewer, *Item, bool) {
	n, v, ok := h.workedInbox(w, r)
	if !ok {
		return nil, nil, nil, false
	}
	job, err := strconv.ParseInt(r.PathValue("job"), 10, 32)
	if err != nil || job <= 0 {
		h.writeErr(w, fmt.Errorf("%w: job %q", ErrNotFound, r.PathValue("job")))
		return nil, nil, nil, false
	}
	it, err := h.m.Item(r.Context(), itemID(int32(job), n.Name))
	if err != nil {
		h.writeErr(w, err)
		return nil, nil, nil, false
	}
	return n, v, it, true
}

// inbox answers an inbox, or 304 while it is what the client has (the
// ETag is of the answer).
func (h *handler) inbox(w http.ResponseWriter, r *http.Request) {
	n, _, ok := h.workedInbox(w, r)
	if !ok {
		return
	}
	view, err := h.m.Inbox(r.Context(), n.Name)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	b, err := json.Marshal(view)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	sum := sha256.Sum256(b)
	etag := `"` + hex.EncodeToString(sum[:8]) + `"`
	w.Header().Set("ETag", etag)
	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	writeJSON(w, http.StatusOK, view)
}

func (h *handler) getItem(w http.ResponseWriter, r *http.Request) {
	if _, _, it, ok := h.workedItem(w, r); ok {
		writeJSON(w, http.StatusOK, it)
	}
}

// itemPayload is the payload an item holds; with pages, a .bdf document's
// alone.
func (h *handler) itemPayload(w http.ResponseWriter, r *http.Request, it *Item, pages bool) (*Payload, bool) {
	j, err := h.m.JobOfNumber(r.Context(), it.Job)
	if err != nil {
		h.writeErr(w, err)
		return nil, false
	}
	p := j.payload(it.Payload)
	switch {
	case p == nil:
		err = fmt.Errorf("%w: job %d has no payload %q", ErrNotFound, it.Job, it.Payload)
	case pages && p.Kind != PayloadBDF:
		err = fmt.Errorf("%w: the item holds %s (%s), which has no pages: its document is …/document", ErrNotFound, or(p.Format, p.Kind), p.Kind)
	}
	if err != nil {
		h.writeErr(w, err)
		return nil, false
	}
	return p, true
}

// itemThumbnail draws the first page of an item's document, the whole of
// it in ?size= pixels (256 by default) on its longer side; ?mode=crop
// draws the square at its top instead, auto as the kind of document asks
// (package thumbnail).
func (h *handler) itemThumbnail(w http.ResponseWriter, r *http.Request) {
	_, _, it, ok := h.workedItem(w, r)
	if !ok {
		return
	}
	p, ok := h.itemPayload(w, r, it, true)
	if !ok {
		return
	}
	q := r.URL.Query()
	size, _ := strconv.Atoi(q.Get("size"))
	if size == 0 {
		size = thumbnail.DefaultSize
	}
	size = min(max(size, 32), 1024)
	mode, err := thumbnail.ParseMode(strings.ToLower(firstOr(q.Get("mode"), "fit")))
	if err != nil {
		h.writeErr(w, fmt.Errorf("%w: mode %q: fit, crop or auto", ErrInvalid, q.Get("mode")))
		return
	}
	doc, err := h.m.Document(r.Context(), p.ID)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	view, err := pagesView(doc)
	if err != nil {
		h.writeErr(w, fmt.Errorf("%w: %v", ErrNotFound, err))
		return
	}
	res, err := thumbnail.Make(doc, &thumbnail.Options{Size: size, Mode: mode, View: view.ID,
		Raster: imagebdf.Options{FontDirs: h.m.o.FontDirs, NoSystemFonts: h.m.o.NoSystemFonts}})
	if err != nil {
		h.writeErr(w, err)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "private, max-age=3600")
	(&png.Encoder{CompressionLevel: png.BestSpeed}).Encode(w, res.Image)
}

// itemPage draws a page of an item's document: …/pages/{n}.png?w= (800
// pixels wide by default).
func (h *handler) itemPage(w http.ResponseWriter, r *http.Request) {
	_, _, it, ok := h.workedItem(w, r)
	if !ok {
		return
	}
	n, err := strconv.Atoi(strings.TrimSuffix(r.PathValue("page"), ".png"))
	if err != nil {
		h.writeErr(w, ErrNotFound)
		return
	}
	p, ok := h.itemPayload(w, r, it, true)
	if !ok {
		return
	}
	width, _ := strconv.Atoi(r.URL.Query().Get("w"))
	if width == 0 {
		width = 800
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "private, max-age=3600")
	if err := h.m.WritePage(r.Context(), p.ID, p.Pages, n-1, width, w); err != nil {
		h.writeErr(w, err)
	}
}

// itemDocument downloads the document an item holds, as it is held. The
// download is in the audit log.
func (h *handler) itemDocument(w http.ResponseWriter, r *http.Request) {
	n, v, it, ok := h.workedItem(w, r)
	if !ok {
		return
	}
	p, ok := h.itemPayload(w, r, it, false)
	if !ok {
		return
	}
	ext := extensionOf(p)
	name := firstOr(it.Name, p.Name)
	if name == "" || strings.ContainsAny(name, `"\/`) {
		name = "job-" + jobID(it.Job) + ext
	} else if path.Ext(name) == "" {
		name += ext
	}
	h.m.audit(r.Context(), v.Label(), "inbox-document", n.Name, fmt.Sprintf("job %d %s", it.Job, clip(it.Name, 100)))
	w.Header().Set("Content-Type", contentTypeOf(p.Format, p.Name))
	w.Header().Set("Content-Length", strconv.FormatInt(p.Size, 10))
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("Content-Disposition", `attachment; filename="`+name+`"`)
	if err := h.m.WritePayload(r.Context(), p.ID, w); err != nil {
		panic(http.ErrAbortHandler)
	}
}

// worker is the name a release or a discard is recorded under: the one
// the body gives (a person at a shared screen), else the viewer's.
func worker(v *Viewer, given string) Mark {
	return Mark{Worker: clip(strings.TrimSpace(firstOr(given, v.Label())), 80)}
}

// release sends an item on: {"to": [nodes], "worker": name}; no nodes,
// every node the inbox releases to.
func (h *handler) release(w http.ResponseWriter, r *http.Request) {
	n, v, it, ok := h.workedItem(w, r)
	if !ok {
		return
	}
	var in struct {
		To     []string `json:"to"`
		Worker string   `json:"worker"`
	}
	if r.ContentLength != 0 {
		if err := readJSON(r, &in); err != nil {
			h.writeErr(w, err)
			return
		}
	}
	if _, to := h.m.pullsOf(r.Context(), n); len(to) == 0 && !v.Can(RoleTopology) {
		h.writeErr(w, fmt.Errorf("%w: %s names no nodes to release to: the role %s releases what it holds", ErrForbidden, n.Name, RoleTopology))
		return
	}
	it, err := h.m.Release(r.Context(), it.ID, in.To, worker(v, in.Worker))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	detail := strings.Join(it.Released.To, ", ")
	if in.Worker != "" {
		detail += " by " + in.Worker
	}
	h.m.audit(r.Context(), v.Label(), "release-item", it.ID, detail)
	writeJSON(w, http.StatusOK, it)
}

// discardItem drops an item: {"worker": name}.
func (h *handler) discardItem(w http.ResponseWriter, r *http.Request) {
	_, v, it, ok := h.workedItem(w, r)
	if !ok {
		return
	}
	var in struct {
		Worker string `json:"worker"`
	}
	if r.ContentLength != 0 {
		if err := readJSON(r, &in); err != nil {
			h.writeErr(w, err)
			return
		}
	}
	it, err := h.m.Discard(r.Context(), it.ID, worker(v, in.Worker))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	h.m.audit(r.Context(), v.Label(), "discard-item", it.ID, in.Worker)
	writeJSON(w, http.StatusOK, it)
}
