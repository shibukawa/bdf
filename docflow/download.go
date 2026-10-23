package docflow

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/yuruidb"
)

// A downloader is an exit that applications and people fetch documents
// from, at GET /api/download/<name>. It takes a stocked file
// (PayloadStocked) alone: what an inbox or an archive right before it
// keeps, which says where the document comes from and how long it is
// there. A job goes on from an archive once it wrote (the archive's To)
// and stays there for good; an inbox lists the job it holds at the
// downloaders of its To as it holds it, and unlists it when a person
// releases or discards it. The downloader lists the job as it arrives: an
// entry with where it is kept and whose it is, the value of the metadata
// key By (a customer's id the sender gave, say). The documents stay where
// they are kept and the downloader reads them from there: an archive's do
// not go when the job's record does, and an inbox's keep the job while
// held.
//
// Who may read what: an app the downloader lists among its Apps, with the
// scope download, reads every entry (a backend that knows its own users;
// ?owner= narrows to one owner's); so do the administrators and the
// auditors. A person reads the entries whose owner is them: what By names
// is compared with the person's email (Subject, the default), subject at
// the IdP, ID or groups. A downloader with no By shows a person what they
// could see of the job (their own, their groups'). A document that is not
// theirs is not there (404), not forbidden: the downloader does not say what
// it holds. A person reads with the pages' session, or with a token a
// frontend got for them (oauth.go).
//
// The entries are records of their own (download-entries, <node>.<job>),
// and the lists (download-lists) are the ids of the latest Entries of
// them: one list of every entry, and one per owner. An entry that rolled
// off the lists is still read by its job's number.

const (
	defaultEntries = 1000
	maxEntries     = 5000
	maxFields      = 16
)

func (n *Node) entries() int { return or(n.Entries, defaultEntries) }

// stockRef is where an inbox or an archive keeps a job's payload: the
// node, and the key of the objects an archive wrote (archiveKey) or the
// item an inbox holds.
type stockRef struct {
	Node string
	Key  string
	Item string
}

// DownloadEntry is a document a downloader lists: the job as it reached
// the downloader, where it is kept (an archive, or an inbox), and whose it
// is.
type DownloadEntry struct {
	Job    int32  `json:"job"`
	Name   string `json:"name,omitempty"`
	Kind   string `json:"kind"`
	Format string `json:"format,omitempty"`
	Pages  int    `json:"pages,omitempty"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
	// Received is when the job was received, Listed when it reached the
	// downloader.
	Received time.Time `json:"received"`
	Listed   time.Time `json:"listed"`
	// Archive is the archive that wrote it, Key the key of its objects
	// there and Object the payload's object (the record is Key +
	// ".json"); or Inbox is the inbox that holds it, and Item the item.
	// Payload is the payload's id in the job.
	Archive string `json:"archive,omitempty"`
	Key     string `json:"key,omitempty"`
	Object  string `json:"object,omitempty"`
	Inbox   string `json:"inbox,omitempty"`
	Item    string `json:"item,omitempty"`
	Payload string `json:"payload,omitempty"`
	// Trace is the step the document came out of or went into, which
	// names the other side's objects (trace.go).
	Trace *Trace `json:"trace,omitempty"`
	// Owner is whose the document is: the value of the downloader's By in
	// the job's metadata; empty when the job had none.
	Owner string `json:"owner,omitempty"`
	// Fields are the metadata the downloader copies (Node.Fields), as text.
	Fields map[string]string `json:"fields,omitempty"`
	// SubmittedBy, ClaimedUser and Group are the job's: who sent it, and
	// whose it is when the downloader has no By.
	SubmittedBy *Principal `json:"submitted_by,omitempty"`
	ClaimedUser string     `json:"claimed_user,omitempty"`
	Group       string     `json:"group,omitempty"`
}

// downloadList is a list of a downloader's entries, newest last: every one's
// (the downloader's name), or an owner's (<name>.<owner's hash>).
type downloadList struct {
	Seq  int64    `json:"seq"`
	Jobs []string `json:"jobs,omitempty"`
}

// metaKeyOrDollarRE is the shape of a downloader's By and Fields: a sender's
// key, a $ key of the manager's, or name and user.
var metaKeyOrDollarRE = regexp.MustCompile(`^\$?[A-Za-z0-9_][A-Za-z0-9_./-]{0,63}$`)

// checkDownload cleans a downloader's settings.
func (m *Manager) checkDownload(ctx context.Context, in Node, n *Node) error {
	n.By = strings.TrimSpace(in.By)
	if n.By != "" && !metaKeyOrDollarRE.MatchString(n.By) {
		return fmt.Errorf("%w: by %q: a metadata key", ErrInvalid, in.By)
	}
	n.Subject = or(strings.TrimSpace(in.Subject), SubjectEmail)
	if !slices.Contains(ownerSubjects, n.Subject) {
		return fmt.Errorf("%w: subject %q: %s", ErrInvalid, in.Subject, strings.Join(ownerSubjects, ", "))
	}
	n.Fields = cleanNames(in.Fields)
	if len(n.Fields) > maxFields {
		return fmt.Errorf("%w: %d fields; at most %d", ErrInvalid, len(n.Fields), maxFields)
	}
	for _, f := range n.Fields {
		if !metaKeyOrDollarRE.MatchString(f) {
			return fmt.Errorf("%w: field %q: a metadata key", ErrInvalid, f)
		}
	}
	n.Entries = in.Entries
	if n.Entries < 0 || n.Entries > maxEntries {
		return fmt.Errorf("%w: entries %d: a downloader lists 1 to %d entries", ErrInvalid, in.Entries, maxEntries)
	}
	n.Apps = cleanNames(in.Apps)
	for _, id := range n.Apps {
		if _, err := m.App(ctx, id); err != nil {
			return fmt.Errorf("%w: no app %q", ErrInvalid, id)
		}
	}
	return nil
}

// fromStock refuses a new edge of n to a node that takes a stocked file
// unless n keeps what it sends there: an inbox or an
// archive, along its To. Typing refuses most such edges (what reaches n
// is not a stocked file), but not those of a node nothing reaches yet,
// to which a job may be sent straight. An edge n had already is left, as
// typing leaves a broken edge: a job gets a leg of the error there.
func (m *Manager) fromStock(ctx context.Context, n *Node) error {
	had := map[string]bool{}
	if cur, err := m.Node(ctx, n.Name); err == nil {
		for _, to := range cur.edges() {
			had[to] = true
		}
	}
	for _, to := range n.edges() {
		if had[to] {
			continue
		}
		x, err := m.Node(ctx, to)
		if err != nil || !takes(x.Kind).has(PayloadStocked) || n.stockedTo(x) {
			continue
		}
		return fmt.Errorf("%w: %s takes a stocked file, which only an inbox or an archive sends, along its to: not %s (%s)", ErrInvalid, to, n.Name, n.Kind)
	}
	return nil
}

func downloadEntryID(node, job string) string { return node + "." + job }

// ownerHash names an owner's list: the owner's value is anything.
func ownerHash(owner string) string {
	sum := sha256.Sum256([]byte(owner))
	return hex.EncodeToString(sum[:8])
}

// listKey is the owner's value as the downloader compares it (sameOwner):
// an email without case, so that Alice@Example.com and alice@example.com
// are one owner's list.
func (n *Node) listKey(owner string) string {
	if n.Subject == "" || n.Subject == SubjectEmail {
		return strings.ToLower(owner)
	}
	return owner
}

func downloadListID(node, owner string) string {
	if owner == "" {
		return node
	}
	return node + "." + ownerHash(owner)
}

// routeDownload lists the job at a downloader: what the inbox or the archive
// right before it keeps, under the owner its metadata names.
func (m *Manager) routeDownload(ctx context.Context, r *routing, name, from string) {
	n, err := m.nodeOfKind(ctx, name, KindDownload)
	if err != nil {
		r.add(Leg{Node: name, Kind: KindDownload, From: from, State: LegError, Error: err.Error()})
		return
	}
	if r.stock == nil || r.stock.Node != from || r.payload == nil {
		r.add(Leg{Node: name, Kind: KindDownload, From: from, State: LegError, Error: "not from an inbox or an archive right before it: a downloader lists what one of them keeps"})
		return
	}
	j, p := r.job, r.payload
	view := metaView(j, p, r.trace)
	e := DownloadEntry{Job: j.Number, Name: j.Name, Kind: p.Kind, Format: p.Format, Pages: p.Pages, Size: p.Size, SHA256: p.SHA256, Received: j.Received, Listed: r.now,
		Payload: p.ID, Trace: r.trace, SubmittedBy: j.SubmittedBy, ClaimedUser: j.ClaimedUser, Group: j.Group}
	if r.stock.Item != "" {
		e.Inbox, e.Item = r.stock.Node, r.stock.Item
	} else {
		e.Archive, e.Key, e.Object = r.stock.Node, r.stock.Key, r.stock.Key+extensionOf(p)
	}
	if n.By != "" {
		e.Owner = clip(view[n.By], 255)
	}
	for _, f := range n.Fields {
		if v, ok := view[f]; ok {
			if e.Fields == nil {
				e.Fields = map[string]string{}
			}
			e.Fields[f] = clip(v, 255)
		}
	}
	if _, err := m.st.db.Insert(ctx, downloadEntryID(name, j.ID), downloadEntryRecord(e)); err != nil && !errors.Is(err, yuruidb.ErrAlreadyExists) {
		r.add(Leg{Node: name, Kind: KindDownload, From: from, State: LegError, Error: mapErr(err).Error()})
		return
	}
	// every one's list, and the owner's when the job has one
	lists := []string{downloadListID(name, "")}
	if e.Owner != "" {
		lists = append(lists, downloadListID(name, n.listKey(e.Owner)))
	}
	for _, id := range lists {
		if err := m.listAt(ctx, id, j.ID, n.entries()); err != nil {
			r.add(Leg{Node: name, Kind: KindDownload, From: from, State: LegError, Error: err.Error()})
			return
		}
	}
	r.add(Leg{Node: name, Kind: KindDownload, From: from, State: LegListed, Member: e.Owner})
}

// listAt adds a job to a list of a downloader, newest last, keeping the
// latest limit of them; a job there already stays where it is.
func (m *Manager) listAt(ctx context.Context, id, job string, limit int) error {
	_, err := change(ctx, m.st.db, id, true, func(l *downloadList) error {
		if slices.Contains(l.Jobs, job) {
			return errSkip
		}
		l.Jobs = append(l.Jobs, job)
		if extra := len(l.Jobs) - limit; extra > 0 {
			l.Jobs = l.Jobs[extra:]
		}
		l.Seq++
		return nil
	})
	return err
}

// unlist takes a job off a downloader: its entry, and its place in the
// lists (every one's, and its owner's). A job not listed there is left.
func (m *Manager) unlist(ctx context.Context, name string, job int32) error {
	id := jobID(job)
	rec, err := m.st.db.Get[downloadEntryRecord](ctx, downloadEntryID(name, id))
	if errors.Is(err, yuruidb.ErrNotFound) {
		return nil
	}
	if err != nil {
		return mapErr(err)
	}
	lists := []string{downloadListID(name, "")}
	if n, err := m.Node(ctx, name); err == nil && rec.Data.Owner != "" {
		lists = append(lists, downloadListID(name, n.listKey(rec.Data.Owner)))
	}
	for _, list := range lists {
		if _, err := change(ctx, m.st.db, list, false, func(l *downloadList) error {
			if !slices.Contains(l.Jobs, id) {
				return errSkip
			}
			l.Jobs = slices.DeleteFunc(l.Jobs, func(x string) bool { return x == id })
			l.Seq++
			return nil
		}); err != nil && !errors.Is(err, ErrNotFound) {
			return err
		}
	}
	if err := m.st.db.Delete[downloadEntryRecord](ctx, rec.ID(), rec.Revision()); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
		return mapErr(err)
	}
	return nil
}

// kept tells whether what an entry lists is still where it is kept: an
// inbox's item while it is held (an archive keeps it for good).
func (m *Manager) kept(ctx context.Context, e *DownloadEntry) bool {
	if e.Item == "" {
		return true
	}
	it, err := m.Item(ctx, e.Item)
	return err == nil && it.State == ItemHeld
}

// downloadJobs are the ids of the jobs a list of a downloader holds, newest
// first.
func (m *Manager) downloadJobs(ctx context.Context, id string) []string {
	rec, err := m.st.db.Get[downloadList](ctx, id)
	if err != nil {
		return nil
	}
	ids := slices.Clone(rec.Data.Jobs)
	slices.Reverse(ids)
	return ids
}

// Entry returns an entry of a downloader, by the job's number.
func (m *Manager) Entry(ctx context.Context, node string, job int32) (*DownloadEntry, error) {
	rec, err := m.st.db.Get[downloadEntryRecord](ctx, downloadEntryID(node, jobID(job)))
	if err != nil {
		return nil, mapErr(err)
	}
	e := DownloadEntry(rec.Data)
	if !m.kept(ctx, &e) {
		return nil, fmt.Errorf("%w: job %d is no longer held at %s", ErrNotFound, job, e.Inbox)
	}
	return &e, nil
}

// Entries lists the entries of the lists named, newest first, up to
// limit: every one's (""), or the owners' (an owner for each of a
// person's groups).
func (m *Manager) Entries(ctx context.Context, n *Node, owners []string, limit int) ([]DownloadEntry, error) {
	if limit <= 0 || limit > maxEntries {
		limit = 100
	}
	node := n.Name
	seen := map[string]bool{}
	var out []DownloadEntry
	for _, owner := range owners {
		if owner != "" {
			owner = n.listKey(owner)
		}
		for _, id := range m.downloadJobs(ctx, downloadListID(node, owner)) {
			if seen[id] {
				continue
			}
			seen[id] = true
			rec, err := m.st.db.Get[downloadEntryRecord](ctx, downloadEntryID(node, id))
			if errors.Is(err, yuruidb.ErrNotFound) {
				continue
			}
			if err != nil {
				return nil, mapErr(err)
			}
			if e := DownloadEntry(rec.Data); m.kept(ctx, &e) {
				out = append(out, e)
			}
		}
	}
	slices.SortFunc(out, func(a, b DownloadEntry) int { return int(b.Job - a.Job) })
	if len(out) > limit {
		out = out[:limit]
	}
	if out == nil {
		out = []DownloadEntry{}
	}
	return out, nil
}

// dropDownload forgets a downloader's entries and lists (the downloader is
// deleted); the archives and the inboxes keep the documents.
func (m *Manager) dropDownload(ctx context.Context, name string) error {
	entries, err := list[downloadEntryRecord](ctx, m.st.db)
	if err != nil {
		return err
	}
	for _, r := range entries {
		if rest, ok := strings.CutPrefix(r.ID(), name+"."); ok && len(rest) == 10 && isDigits(rest) {
			if err := m.st.db.Delete[downloadEntryRecord](ctx, r.ID(), r.Revision()); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
				return mapErr(err)
			}
		}
	}
	lists, err := list[downloadList](ctx, m.st.db)
	if err != nil {
		return err
	}
	for _, r := range lists {
		rest, ok := strings.CutPrefix(r.ID(), name+".")
		if r.ID() == name || ok && len(rest) == 16 && isHex(rest) {
			if err := m.st.db.Delete[downloadList](ctx, r.ID(), r.Revision()); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
				return mapErr(err)
			}
		}
	}
	return nil
}

func isDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return s != ""
}

func isHex(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil && s != ""
}

// The downloader's API: GET /api/download/<name> lists the entries the
// viewer may see (?owner= one owner's, ?limit=), GET
// /api/download/<name>/<job> is one entry, …/record the archive's record
// of it (of an inbox's, the job as it is), and …/document the document
// itself, read from where it is kept. What the viewer may not see is not
// found.

// downloadViewer finds who reads a downloader, and what they may read of
// it: every entry (an administrator, an auditor, an app the downloader
// lists with the scope download), or the entries of the owners returned
// (a person: their identities under the downloader's subject). all false
// with no owners is a person who sees the entries of the jobs they could
// see (a downloader with no By).
func (h *handler) downloadViewer(w http.ResponseWriter, r *http.Request, n *Node) (v *Viewer, all bool, owners []string, ok bool) {
	v, signedIn := h.viewerOf(r)
	if !signedIn {
		w.Header().Set("WWW-Authenticate", `Bearer realm="bdf docflow download"`)
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "a token or a sign-in is required"})
		return nil, false, nil, false
	}
	switch {
	case v.Admin || v.Can(RoleAudit):
		return v, true, nil, true
	case v.Kind == PrincipalApp:
		if !slices.Contains(n.Apps, v.ID) {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": n.Name + " does not take the app " + v.ID})
			return nil, false, nil, false
		}
		if !v.Has(ScopeDownload) {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "the app " + v.ID + " has no scope " + ScopeDownload})
			return nil, false, nil, false
		}
		return v, true, nil, true
	case v.Kind == PrincipalUser:
		if v.Token() && !v.Has(ScopeDownload) {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "the token has no scope " + ScopeDownload})
			return nil, false, nil, false
		}
		if !v.CanSendTo(n.Scope, n.Owner) || !sameOrigin(r) {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "you may not read " + n.Name})
			return nil, false, nil, false
		}
		if n.By == "" {
			return v, false, nil, true
		}
		return v, false, v.identities(n.Subject), true
	}
	writeJSON(w, http.StatusForbidden, map[string]string{"error": "the downloader is read by people and apps"})
	return nil, false, nil, false
}

// canSeeEntry tells whether a person sees an entry of a downloader with no
// By: as they would the job (their own, their groups').
func (v *Viewer) canSeeEntry(e *DownloadEntry) bool {
	if v == nil {
		return false
	}
	by := e.SubmittedBy
	return v.Admin || v.Can(RoleAudit) || v.Member(e.Group) || by != nil && by.Kind == v.Kind && by.ID != "" && by.ID == v.ID
}

// sees tells whether the viewer may read an entry, as downloadViewer told
// it.
func (n *Node) sees(v *Viewer, all bool, owners []string, e *DownloadEntry, owner string) bool {
	if owner != "" && !n.sameOwner(owner, e.Owner) {
		return false
	}
	switch {
	case all:
		return true
	case n.By == "":
		return v.canSeeEntry(e)
	}
	return e.Owner != "" && slices.ContainsFunc(owners, func(id string) bool { return n.sameOwner(id, e.Owner) })
}

// sameOwner compares an identity with an owner under the downloader's
// subject: an email without case.
func (n *Node) sameOwner(id, owner string) bool {
	if id == "" || owner == "" {
		return false
	}
	if n.Subject == "" || n.Subject == SubjectEmail {
		return strings.EqualFold(id, owner)
	}
	return id == owner
}

func (h *handler) download(w http.ResponseWriter, r *http.Request) (*Node, bool) {
	n, err := h.m.nodeOfKind(r.Context(), r.PathValue("node"), KindDownload)
	if err != nil {
		h.writeErr(w, err)
		return nil, false
	}
	return n, true
}

// listDownload is GET /api/download/<name>: the entries the viewer may see, newest
// first; ?owner= one owner's (of a person, one of theirs), ?limit=.
func (h *handler) listDownload(w http.ResponseWriter, r *http.Request) {
	n, ok := h.download(w, r)
	if !ok {
		return
	}
	v, all, owners, ok := h.downloadViewer(w, r, n)
	if !ok {
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	owner := r.URL.Query().Get("owner")
	lists := []string{""}
	switch {
	case all && owner != "":
		lists = []string{owner}
	case all:
	case n.By != "":
		lists = owners
		if owner != "" {
			lists = nil
			for _, id := range owners {
				if n.sameOwner(id, owner) {
					lists = []string{id}
				}
			}
		}
	}
	es, err := h.m.Entries(r.Context(), n, lists, limit)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	if !all && n.By == "" {
		es = slices.DeleteFunc(es, func(e DownloadEntry) bool { return !v.canSeeEntry(&e) || owner != "" && !n.sameOwner(owner, e.Owner) })
	}
	writeJSON(w, http.StatusOK, es)
}

// downloadEntry is GET /api/download/<name>/<job>, and its record and
// document: an entry the viewer may see, or not found.
func (h *handler) downloadEntry(w http.ResponseWriter, r *http.Request) (*Node, *Viewer, *DownloadEntry, bool) {
	n, ok := h.download(w, r)
	if !ok {
		return nil, nil, nil, false
	}
	v, all, owners, ok := h.downloadViewer(w, r, n)
	if !ok {
		return nil, nil, nil, false
	}
	job, err := strconv.ParseInt(r.PathValue("job"), 10, 32)
	if err != nil || job <= 0 {
		h.writeErr(w, fmt.Errorf("%w: job %q", ErrNotFound, r.PathValue("job")))
		return nil, nil, nil, false
	}
	e, err := h.m.Entry(r.Context(), n.Name, int32(job))
	if err == nil && !n.sees(v, all, owners, e, r.URL.Query().Get("owner")) {
		err = fmt.Errorf("%w: job %d is not at %s", ErrNotFound, job, n.Name)
	}
	if err != nil {
		h.writeErr(w, err)
		return nil, nil, nil, false
	}
	return n, v, e, true
}

func (h *handler) getDownloadEntry(w http.ResponseWriter, r *http.Request) {
	if _, _, e, ok := h.downloadEntry(w, r); ok {
		writeJSON(w, http.StatusOK, e)
	}
}

// archiveOf opens the store of the archive an entry names.
func (m *Manager) archiveOf(ctx context.Context, e *DownloadEntry) (*Node, error) {
	a, err := m.nodeOfKind(ctx, e.Archive, KindArchive)
	if err != nil {
		return nil, fmt.Errorf("%w: the archive %s of job %d is gone", ErrNotFound, e.Archive, e.Job)
	}
	return a, nil
}

// downloadObject streams an object of the archive that wrote an entry, or
// of an inbox's the job's payload it holds.
func (h *handler) downloadObject(w http.ResponseWriter, r *http.Request, e *DownloadEntry, object, contentType, disposition string) {
	st := h.m.st.docs
	if e.Inbox != "" {
		object = e.Payload
	} else {
		a, err := h.m.archiveOf(r.Context(), e)
		if err != nil {
			h.writeErr(w, err)
			return
		}
		if st, err = h.m.st.archive(r.Context(), a.Store, a.collection()); err != nil {
			h.writeErr(w, err)
			return
		}
	}
	rd, err := st.Open(r.Context(), object)
	if err != nil {
		h.writeErr(w, mapErr(err))
		return
	}
	defer rd.Close()
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.FormatInt(rd.Size, 10))
	w.Header().Set("Cache-Control", "private, no-store")
	if disposition != "" {
		w.Header().Set("Content-Disposition", disposition)
	}
	if _, err := io.Copy(w, rd); err != nil {
		panic(http.ErrAbortHandler)
	}
}

// downloadRecord is …/record: the archive's record of the job, as it
// wrote it (ArchiveRecord: its metadata, payloads and legs then); of an
// inbox's entry, the job as it is now, which nothing wrote down.
func (h *handler) downloadRecord(w http.ResponseWriter, r *http.Request) {
	_, _, e, ok := h.downloadEntry(w, r)
	if !ok {
		return
	}
	if e.Inbox != "" {
		j, err := h.m.JobOfNumber(r.Context(), e.Job)
		if err != nil {
			h.writeErr(w, err)
			return
		}
		w.Header().Set("Cache-Control", "private, no-store")
		writeJSON(w, http.StatusOK, j)
		return
	}
	h.downloadObject(w, r, e, e.Key+".json", "application/json; charset=utf-8", "")
}

// downloadDocument is …/document: the document the archive wrote, as it
// was. The download is in the audit log.
func (h *handler) downloadDocument(w http.ResponseWriter, r *http.Request) {
	n, v, e, ok := h.downloadEntry(w, r)
	if !ok {
		return
	}
	ext := extensionOf(&Payload{Name: e.Name, Format: e.Format})
	name := e.Name
	if name == "" || strings.ContainsAny(name, `"\/`) {
		name = "job-" + jobID(e.Job) + ext
	} else if path.Ext(name) == "" {
		name += ext
	}
	h.m.audit(r.Context(), v.Label(), "download-document", n.Name, fmt.Sprintf("job %d %s", e.Job, clip(e.Name, 100)))
	h.downloadObject(w, r, e, e.Object, contentTypeOf(e.Format, e.Name), `attachment; filename="`+name+`"`)
}
