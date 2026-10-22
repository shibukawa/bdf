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

// An export is an exit that applications and people fetch documents
// from, at GET /out/<name>: the documents an archive before it wrote.
// The job goes on from an archive once it wrote (the archive's To), and
// the export lists it as it arrives: an entry with where the archive
// wrote it and whose it is, the value of the metadata key By (a
// customer's id the sender gave, say). The documents themselves stay in
// the archive, which keeps them for good; the export reads them from
// there, so what it lists does not go when the job's record does.
//
// Who may read what: an app the export lists among its Apps, with the
// scope export, reads every entry (a backend that knows its own users;
// ?owner= narrows to one owner's); so do the administrators and the
// auditors. A person reads the entries whose owner is them: what By names
// is compared with the person's email (Subject, the default), subject at
// the IdP, ID or groups. An export with no By shows a person what they
// could see of the job (their own, their groups'). A document that is not
// theirs is not there (404), not forbidden: the export does not say what
// it holds. A person reads with the pages' session, or with a token a
// frontend got for them (oauth.go).
//
// The entries are records of their own (export-entries, <node>.<job>),
// and the lists (export-lists) are the ids of the latest Entries of
// them: one list of every entry, and one per owner. An entry that rolled
// off the lists is still read by its job's number.

// OutPath is where the exports are.
const OutPath = "/out"

const (
	defaultEntries = 1000
	maxEntries     = 5000
	maxFields      = 16
)

func (n *Node) entries() int { return or(n.Entries, defaultEntries) }

// archiveRef is what an archive wrote of a job: the archive, and the key
// of its objects there (archiveKey).
type archiveRef struct {
	Node string
	Key  string
}

// ExportEntry is a document an export lists: the job as it reached the
// export, where the archive wrote it, and whose it is.
type ExportEntry struct {
	Job    int32  `json:"job"`
	Name   string `json:"name,omitempty"`
	Kind   string `json:"kind"`
	Format string `json:"format,omitempty"`
	Pages  int    `json:"pages,omitempty"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
	// Received is when the job was received, Listed when it reached the
	// export.
	Received time.Time `json:"received"`
	Listed   time.Time `json:"listed"`
	// Archive is the archive that wrote it, Key the key of its objects
	// there, Object the payload's object (the record is Key + ".json"),
	// and Payload the payload's id in the job.
	Archive string `json:"archive"`
	Key     string `json:"key"`
	Object  string `json:"object"`
	Payload string `json:"payload,omitempty"`
	// Trace is the step the document came out of or went into, which
	// names the other side's objects (trace.go).
	Trace *Trace `json:"trace,omitempty"`
	// Owner is whose the document is: the value of the export's By in
	// the job's metadata; empty when the job had none.
	Owner string `json:"owner,omitempty"`
	// Fields are the metadata the export copies (Node.Fields), as text.
	Fields map[string]string `json:"fields,omitempty"`
	// SubmittedBy, ClaimedUser and Group are the job's: who sent it, and
	// whose it is when the export has no By.
	SubmittedBy *Principal `json:"submitted_by,omitempty"`
	ClaimedUser string     `json:"claimed_user,omitempty"`
	Group       string     `json:"group,omitempty"`
}

// exportList is a list of an export's entries, newest last: every one's
// (the export's name), or an owner's (<name>.<owner's hash>).
type exportList struct {
	Seq  int64    `json:"seq"`
	Jobs []string `json:"jobs,omitempty"`
}

// metaKeyOrDollarRE is the shape of an export's By and Fields: a sender's
// key, a $ key of the manager's, or name and user.
var metaKeyOrDollarRE = regexp.MustCompile(`^\$?[A-Za-z0-9_][A-Za-z0-9_./-]{0,63}$`)

// checkExport cleans an export's settings.
func (m *Manager) checkExport(ctx context.Context, in Node, n *Node) error {
	n.By = strings.TrimSpace(in.By)
	if n.By != "" && !metaKeyOrDollarRE.MatchString(n.By) {
		return fmt.Errorf("%w: by %q: a metadata key", ErrInvalid, in.By)
	}
	n.Subject = or(strings.TrimSpace(in.Subject), SubjectEmail)
	if !slices.Contains(exportSubjects, n.Subject) {
		return fmt.Errorf("%w: subject %q: %s", ErrInvalid, in.Subject, strings.Join(exportSubjects, ", "))
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
		return fmt.Errorf("%w: entries %d: an export lists 1 to %d entries", ErrInvalid, in.Entries, maxEntries)
	}
	n.Apps = cleanNames(in.Apps)
	for _, id := range n.Apps {
		if _, err := m.App(ctx, id); err != nil {
			return fmt.Errorf("%w: no app %q", ErrInvalid, id)
		}
	}
	return nil
}

// exportsArchived refuses a graph in which a job may reach an export
// without an archive before it, with n put as it is. The walk starts at
// every entry, every node nothing sends to (a destination of its own),
// and every node whose output is a payload no archive wrote: a template
// or a convert node (what it makes is new) and an inbox (a release is a
// routing of its own); an archive ends the walk (what goes on from it
// was written). A job sent straight to a node in between is checked when
// it is routed.
func (m *Manager) exportsArchived(ctx context.Context, n *Node) error {
	list, err := m.Nodes(ctx)
	if err != nil {
		return err
	}
	nodes := map[string]*Node{}
	for i := range list {
		nodes[list[i].Name] = &list[i]
	}
	nodes[n.Name] = n
	var exports []string
	incoming := map[string]bool{}
	for _, x := range nodes {
		if x.Kind == KindExport {
			exports = append(exports, x.Name)
		}
		for _, to := range x.edges() {
			incoming[to] = true
		}
	}
	if len(exports) == 0 {
		return nil
	}
	from := map[string]string{} // the node a job came from, on a way without an archive
	var walk func(name string)
	walk = func(name string) {
		for _, to := range nodes[name].edges() {
			x := nodes[to]
			if x == nil {
				continue
			}
			if _, seen := from[to]; seen {
				continue
			}
			from[to] = name
			if x.Kind != KindArchive {
				walk(to)
			}
		}
	}
	roots := make([]string, 0, len(nodes))
	for name, x := range nodes {
		switch {
		case x.Kind == KindArchive:
		case isEntry(x.Kind), x.Kind == KindTemplate, x.Kind == KindConvert, x.Kind == KindInbox, !incoming[name]:
			roots = append(roots, name)
		}
	}
	slices.Sort(roots)
	for _, root := range roots {
		if _, seen := from[root]; !seen {
			from[root] = ""
			walk(root)
		}
	}
	slices.Sort(exports)
	for _, name := range exports {
		if _, reached := from[name]; !reached {
			continue
		}
		path := []string{name}
		for prev := from[name]; prev != ""; prev = from[prev] {
			path = append([]string{prev}, path...)
		}
		if len(path) == 1 {
			continue // a root itself: jobs sent straight to it are refused when routed
		}
		return fmt.Errorf("%w: %s (export) lists what an archive wrote, but a job reaches it without one: %s", ErrInvalid, name, strings.Join(path, " → "))
	}
	return nil
}

func exportEntryID(node, job string) string { return node + "." + job }

// ownerHash names an owner's list: the owner's value is anything.
func ownerHash(owner string) string {
	sum := sha256.Sum256([]byte(owner))
	return hex.EncodeToString(sum[:8])
}

// listKey is the owner's value as the export compares it (sameOwner):
// an email without case, so that Alice@Example.com and alice@example.com
// are one owner's list.
func (n *Node) listKey(owner string) string {
	if n.Subject == "" || n.Subject == SubjectEmail {
		return strings.ToLower(owner)
	}
	return owner
}

func exportListID(node, owner string) string {
	if owner == "" {
		return node
	}
	return node + "." + ownerHash(owner)
}

// routeExport lists the job at an export: what the archive before it
// wrote, under the owner its metadata names.
func (m *Manager) routeExport(ctx context.Context, r *routing, name, from string) {
	n, err := m.nodeOfKind(ctx, name, KindExport)
	if err != nil {
		r.add(Leg{Node: name, Kind: KindExport, From: from, State: LegError, Error: err.Error()})
		return
	}
	if r.archived == nil || r.payload == nil {
		r.add(Leg{Node: name, Kind: KindExport, From: from, State: LegError, Error: "reached without an archive before it: an export lists what an archive wrote"})
		return
	}
	j, p := r.job, r.payload
	view := metaView(j, p, r.trace)
	e := ExportEntry{Job: j.Number, Name: j.Name, Kind: p.Kind, Format: p.Format, Pages: p.Pages, Size: p.Size, SHA256: p.SHA256, Received: j.Received, Listed: r.now,
		Archive: r.archived.Node, Key: r.archived.Key, Object: r.archived.Key + extensionOf(p), Payload: p.ID, Trace: r.trace, SubmittedBy: j.SubmittedBy, ClaimedUser: j.ClaimedUser, Group: j.Group}
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
	if _, err := m.st.db.Insert(ctx, exportEntryID(name, j.ID), exportEntryRecord(e)); err != nil && !errors.Is(err, yuruidb.ErrAlreadyExists) {
		r.add(Leg{Node: name, Kind: KindExport, From: from, State: LegError, Error: mapErr(err).Error()})
		return
	}
	// every one's list, and the owner's when the job has one
	lists := []string{exportListID(name, "")}
	if e.Owner != "" {
		lists = append(lists, exportListID(name, n.listKey(e.Owner)))
	}
	for _, id := range lists {
		if err := m.listAt(ctx, id, j.ID, n.entries()); err != nil {
			r.add(Leg{Node: name, Kind: KindExport, From: from, State: LegError, Error: err.Error()})
			return
		}
	}
	r.add(Leg{Node: name, Kind: KindExport, From: from, State: LegListed, Member: e.Owner})
}

// listAt adds a job to a list of an export, newest last, keeping the
// latest limit of them; a job there already stays where it is.
func (m *Manager) listAt(ctx context.Context, id, job string, limit int) error {
	_, err := change(ctx, m.st.db, id, true, func(l *exportList) error {
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

// exportJobs are the ids of the jobs a list of an export holds, newest
// first.
func (m *Manager) exportJobs(ctx context.Context, id string) []string {
	rec, err := m.st.db.Get[exportList](ctx, id)
	if err != nil {
		return nil
	}
	ids := slices.Clone(rec.Data.Jobs)
	slices.Reverse(ids)
	return ids
}

// Entry returns an entry of an export, by the job's number.
func (m *Manager) Entry(ctx context.Context, node string, job int32) (*ExportEntry, error) {
	rec, err := m.st.db.Get[exportEntryRecord](ctx, exportEntryID(node, jobID(job)))
	if err != nil {
		return nil, mapErr(err)
	}
	e := ExportEntry(rec.Data)
	return &e, nil
}

// Entries lists the entries of the lists named, newest first, up to
// limit: every one's (""), or the owners' (an owner for each of a
// person's groups).
func (m *Manager) Entries(ctx context.Context, n *Node, owners []string, limit int) ([]ExportEntry, error) {
	if limit <= 0 || limit > maxEntries {
		limit = 100
	}
	node := n.Name
	seen := map[string]bool{}
	var out []ExportEntry
	for _, owner := range owners {
		if owner != "" {
			owner = n.listKey(owner)
		}
		for _, id := range m.exportJobs(ctx, exportListID(node, owner)) {
			if seen[id] {
				continue
			}
			seen[id] = true
			rec, err := m.st.db.Get[exportEntryRecord](ctx, exportEntryID(node, id))
			if errors.Is(err, yuruidb.ErrNotFound) {
				continue
			}
			if err != nil {
				return nil, mapErr(err)
			}
			out = append(out, ExportEntry(rec.Data))
		}
	}
	slices.SortFunc(out, func(a, b ExportEntry) int { return int(b.Job - a.Job) })
	if len(out) > limit {
		out = out[:limit]
	}
	if out == nil {
		out = []ExportEntry{}
	}
	return out, nil
}

// dropExport forgets an export's entries and lists (the export is
// deleted); the archive keeps the documents.
func (m *Manager) dropExport(ctx context.Context, name string) error {
	entries, err := list[exportEntryRecord](ctx, m.st.db)
	if err != nil {
		return err
	}
	for _, r := range entries {
		if rest, ok := strings.CutPrefix(r.ID(), name+"."); ok && len(rest) == 10 && isDigits(rest) {
			if err := m.st.db.Delete[exportEntryRecord](ctx, r.ID(), r.Revision()); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
				return mapErr(err)
			}
		}
	}
	lists, err := list[exportList](ctx, m.st.db)
	if err != nil {
		return err
	}
	for _, r := range lists {
		rest, ok := strings.CutPrefix(r.ID(), name+".")
		if r.ID() == name || ok && len(rest) == 16 && isHex(rest) {
			if err := m.st.db.Delete[exportList](ctx, r.ID(), r.Revision()); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
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

// The export's endpoints: GET /out/<name> lists the entries the viewer
// may see (?owner= one owner's, ?limit=), GET /out/<name>/<job> is one
// entry, …/record the archive's record of it, and …/document the
// document itself, read from the archive. What the viewer may not see is
// not found.

// exportViewer finds who reads an export, and what they may read of it:
// every entry (an administrator, an auditor, an app the export lists with
// the scope export), or the entries of the owners returned (a person:
// their identities under the export's subject). all false with no owners
// is a person who sees the entries of the jobs they could see (an export
// with no By).
func (h *handler) exportViewer(w http.ResponseWriter, r *http.Request, n *Node) (v *Viewer, all bool, owners []string, ok bool) {
	v, signedIn := h.viewerOf(r)
	if !signedIn {
		w.Header().Set("WWW-Authenticate", `Bearer realm="bdf docflow export"`)
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
		if !v.Has(ScopeExport) {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "the app " + v.ID + " has no scope " + ScopeExport})
			return nil, false, nil, false
		}
		return v, true, nil, true
	case v.Kind == PrincipalUser:
		if v.Token() && !v.Has(ScopeExport) {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "the token has no scope " + ScopeExport})
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
	writeJSON(w, http.StatusForbidden, map[string]string{"error": "the export is read by people and apps"})
	return nil, false, nil, false
}

// canSeeEntry tells whether a person sees an entry of an export with no
// By: as they would the job (their own, their groups').
func (v *Viewer) canSeeEntry(e *ExportEntry) bool {
	if v == nil {
		return false
	}
	by := e.SubmittedBy
	return v.Admin || v.Can(RoleAudit) || v.Member(e.Group) || by != nil && by.Kind == v.Kind && by.ID != "" && by.ID == v.ID
}

// sees tells whether the viewer may read an entry, as exportViewer told
// it.
func (n *Node) sees(v *Viewer, all bool, owners []string, e *ExportEntry, owner string) bool {
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

// sameOwner compares an identity with an owner under the export's
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

func (h *handler) export(w http.ResponseWriter, r *http.Request) (*Node, bool) {
	n, err := h.m.nodeOfKind(r.Context(), r.PathValue("node"), KindExport)
	if err != nil {
		h.writeErr(w, err)
		return nil, false
	}
	return n, true
}

// listExport is GET /out/<name>: the entries the viewer may see, newest
// first; ?owner= one owner's (of a person, one of theirs), ?limit=.
func (h *handler) listExport(w http.ResponseWriter, r *http.Request) {
	n, ok := h.export(w, r)
	if !ok {
		return
	}
	v, all, owners, ok := h.exportViewer(w, r, n)
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
		es = slices.DeleteFunc(es, func(e ExportEntry) bool { return !v.canSeeEntry(&e) || owner != "" && !n.sameOwner(owner, e.Owner) })
	}
	writeJSON(w, http.StatusOK, es)
}

// exportEntry is GET /out/<name>/<job>, and its record and document: an
// entry the viewer may see, or not found.
func (h *handler) exportEntry(w http.ResponseWriter, r *http.Request) (*Node, *Viewer, *ExportEntry, bool) {
	n, ok := h.export(w, r)
	if !ok {
		return nil, nil, nil, false
	}
	v, all, owners, ok := h.exportViewer(w, r, n)
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

func (h *handler) getExportEntry(w http.ResponseWriter, r *http.Request) {
	if _, _, e, ok := h.exportEntry(w, r); ok {
		writeJSON(w, http.StatusOK, e)
	}
}

// archiveOf opens the store of the archive an entry names.
func (m *Manager) archiveOf(ctx context.Context, e *ExportEntry) (*Node, error) {
	a, err := m.nodeOfKind(ctx, e.Archive, KindArchive)
	if err != nil {
		return nil, fmt.Errorf("%w: the archive %s of job %d is gone", ErrNotFound, e.Archive, e.Job)
	}
	return a, nil
}

// exportObject streams an object of the archive that wrote an entry.
func (h *handler) exportObject(w http.ResponseWriter, r *http.Request, e *ExportEntry, object, contentType, disposition string) {
	a, err := h.m.archiveOf(r.Context(), e)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	st, err := h.m.st.archive(r.Context(), a.Store, a.collection())
	if err != nil {
		h.writeErr(w, err)
		return
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

// exportRecord is …/record: the archive's record of the job, as it
// wrote it (ArchiveRecord: its metadata, payloads and legs then).
func (h *handler) exportRecord(w http.ResponseWriter, r *http.Request) {
	_, _, e, ok := h.exportEntry(w, r)
	if !ok {
		return
	}
	h.exportObject(w, r, e, e.Key+".json", "application/json; charset=utf-8", "")
}

// exportDocument is …/document: the document the archive wrote, as it
// was. The download is in the audit log.
func (h *handler) exportDocument(w http.ResponseWriter, r *http.Request) {
	n, v, e, ok := h.exportEntry(w, r)
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
	h.m.audit(r.Context(), v.Label(), "export-document", n.Name, fmt.Sprintf("job %d %s", e.Job, clip(e.Name, 100)))
	h.exportObject(w, r, e, e.Object, contentTypeOf(e.Format, e.Name), `attachment; filename="`+name+`"`)
}
