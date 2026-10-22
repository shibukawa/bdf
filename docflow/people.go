package docflow

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// People sign in to the pages with OIDC: the web application signs them
// in and Options.Viewer tells the manager who they are, with the groups
// the IdP gave them. Scripts and senders use the tokens of the
// configuration: the admin token (an administrator), a submit token (a
// user of its name). What each may see and send to is scope.go's.

// viewerOf finds who a request is from: the admin token (bearer), an
// app's token, a submit token (bearer name:secret, or Basic), a person's
// token (a frontend's, at the exports and /api/me alone), or a person
// signed in to the pages (their session). A request with a credential
// that is not known is nobody's.
func (h *handler) viewerOf(r *http.Request) (*Viewer, bool) {
	tok := bearer(r)
	if sameSecret(h.m.o.AdminToken, tok) {
		return &Viewer{Kind: PrincipalToken, Via: ViaAdminToken, ID: "admin", Name: "admin", Admin: true}, true
	}
	if tok != "" && !strings.Contains(tok, ":") {
		// an app's token (apps.go)
		if v, ok := h.m.appTokenViewer(r.Context(), tok); ok {
			return v, true
		}
		// a person's token (oauth.go): it opens the exports and /api/me,
		// not the pages' API
		if v, ok := h.m.userTokenViewer(r.Context(), tok); ok {
			if strings.HasPrefix(r.URL.Path, OutPath+"/") || r.URL.Path == "/api/me" {
				return v, true
			}
			return nil, false
		}
	}
	name, secret, ok := r.BasicAuth()
	if !ok && tok != "" {
		name, secret, ok = strings.Cut(tok, ":")
	}
	if ok {
		if want, known := h.m.o.SubmitTokens[name]; known && sameSecret(want, secret) {
			return &Viewer{Kind: PrincipalToken, Via: ViaSubmitToken, ID: name, Name: name}, true
		}
		if v, ok := h.m.appBasicViewer(r.Context(), name, secret); ok {
			return v, true
		}
		return nil, false
	}
	if r.Header.Get("Authorization") != "" || h.m.o.Viewer == nil {
		return nil, false
	}
	v := h.m.o.Viewer(r)
	return v, v != nil
}

// sameOrigin tells a request that a page of the manager's own origin
// made, or one that changes nothing. The session's cookie rides on every
// request to the origin, so what it authorizes to change must come from
// the origin's pages.
func sameOrigin(r *http.Request) bool {
	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	}
	if r.Header.Get("Authorization") != "" {
		return true // not the cookie's
	}
	if site := r.Header.Get("Sec-Fetch-Site"); site != "" {
		return site == "same-origin"
	}
	o, err := url.Parse(r.Header.Get("Origin"))
	return err == nil && o.Host != "" && o.Host == r.Host
}

// whoami answers who a request is from.
func (h *handler) whoami(w http.ResponseWriter, r *http.Request) {
	v, ok := h.viewerOf(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "not signed in"})
		return
	}
	me := map[string]any{"kind": v.Kind, "id": v.ID, "name": v.Label(), "email": v.Email, "groups": v.Groups, "admin": v.Admin,
		"topology": v.Can(RoleTopology), "audit": v.Can(RoleAudit)}
	if v.Admin {
		// the pages: whether entry nodes may have ports of their own here
		me["ports"] = !h.m.o.NoPorts
	}
	writeJSON(w, http.StatusOK, me)
}

type viewerHandler func(http.ResponseWriter, *http.Request, *Viewer)

// userOnly serves a request of anyone signed in (or with a token).
func (h *handler) userOnly(fn viewerHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, ok := h.viewerOf(r)
		switch {
		case !ok:
			w.Header().Set("WWW-Authenticate", `Bearer realm="bdf docflow"`)
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "not signed in"})
		case !sameOrigin(r):
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "a request of another origin"})
		default:
			fn(w, r, v)
		}
	}
}

// ownJob serves a request about a job the viewer may see.
func (h *handler) ownJob(fn func(http.ResponseWriter, *http.Request, *Viewer, *Job)) http.HandlerFunc {
	return h.userOnly(func(w http.ResponseWriter, r *http.Request, v *Viewer) {
		j, err := h.m.Job(r.Context(), r.PathValue("id"))
		if err == nil && !v.CanSee(j) {
			err = ErrNotFound
		}
		if err != nil {
			h.writeErr(w, err)
			return
		}
		fn(w, r, v, j)
	})
}

// allJobs is the destination state that lists every job.
const allJobs = "-all"

// MyJobs returns the last jobs a viewer may see, latest first: their
// own, their groups', every job for an administrator.
func (m *Manager) MyJobs(ctx context.Context, v *Viewer, limit int) ([]Job, error) {
	jobs, err := m.jobsOfState(ctx, allJobs, 0)
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = recentJobs
	}
	out := make([]Job, 0, min(limit, len(jobs)))
	for _, j := range jobs {
		if len(out) >= limit {
			break
		}
		if v.CanSee(&j) {
			out = append(out, j)
		}
	}
	return out, nil
}

func (h *handler) myJobs(w http.ResponseWriter, r *http.Request, v *Viewer) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	js, err := h.m.MyJobs(r.Context(), v, limit)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, js)
}

// Destination is a node as a person choosing where to print sees it.
type Destination struct {
	Name  string `json:"name"`
	Kind  string `json:"kind"`
	Title string `json:"title,omitempty"`
	Scope string `json:"scope,omitempty"`
	Owner string `json:"owner,omitempty"`
}

// Destinations returns the nodes a viewer may send to: an app, the
// endpoints that list it.
func (m *Manager) Destinations(ctx context.Context, v *Viewer) ([]Destination, error) {
	topo, err := m.Topology(ctx)
	if err != nil {
		return nil, err
	}
	out := []Destination{}
	for _, n := range topo {
		if n.Kind == KindExport {
			continue // read from, never sent to
		}
		allowed := v.CanSendTo(n.Scope, n.Owner)
		if v.Kind == PrincipalApp && !v.Admin {
			allowed = n.Kind == KindEndpoint && slices.Contains(n.Apps, v.ID)
		}
		if allowed {
			out = append(out, Destination{Name: n.Name, Kind: n.Kind, Title: n.Title, Scope: n.Scope, Owner: n.Owner})
		}
	}
	return out, nil
}

func (h *handler) destinations(w http.ResponseWriter, r *http.Request, v *Viewer) {
	ds, err := h.m.Destinations(r.Context(), v)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ds)
}

// The audit log: who did what, kept newest first.

// AuditEntry is an action a person, a device or a sender took.
type AuditEntry struct {
	ID      string    `json:"id"`
	At      time.Time `json:"at"`
	Who     string    `json:"who"`
	Action  string    `json:"action"`
	Subject string    `json:"subject,omitempty"`
	Detail  string    `json:"detail,omitempty"`
}

// auditSeq orders the entries of one process written at the same instant.
var auditSeq atomic.Uint32

// auditID sorts newest first in the store's order of ids: the time, then
// the sequence of the process, both reversed.
func (m *Manager) auditID(at time.Time) string {
	return fmt.Sprintf("%019d.%010d.%s", math.MaxInt64-at.UnixNano(), math.MaxUint32-auditSeq.Add(1), randomHex(2))
}

// Audit records an action of the web application's (a sign-in).
func (m *Manager) Audit(ctx context.Context, who, action, subject, detail string) {
	m.audit(ctx, who, action, subject, detail)
}

// audit records an action; a record that cannot be written is logged.
func (m *Manager) audit(ctx context.Context, who, action, subject, detail string) {
	at := m.now()
	e := AuditEntry{At: at, Who: clip(who, 80), Action: action, Subject: clip(subject, 255), Detail: clip(detail, 500)}
	e.ID = m.auditID(at)
	if _, err := m.st.db.Insert(context.WithoutCancel(ctx), e.ID, auditRecord(e)); err != nil {
		m.o.Log.Warn("audit", "action", action, "subject", subject, "err", err)
	}
}

// AuditLog returns the last entries, newest first.
func (m *Manager) AuditLog(ctx context.Context, limit int) ([]AuditEntry, error) {
	if limit <= 0 || limit > 1000 {
		limit = 200
	}
	var out []AuditEntry
	cursor := ""
	for len(out) < limit {
		page, err := m.st.db.List[auditRecord](ctx, yuruidbListOptions(min(200, limit-len(out)), cursor))
		if err != nil {
			return nil, mapErr(err)
		}
		for _, r := range page.Records {
			e := AuditEntry(r.Data)
			e.ID = r.ID()
			out = append(out, e)
		}
		if page.Cursor == "" || len(page.Records) == 0 {
			break
		}
		cursor = page.Cursor
	}
	if out == nil {
		out = []AuditEntry{}
	}
	return out, nil
}

func (h *handler) auditLog(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	es, err := h.m.AuditLog(r.Context(), limit)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, es)
}

// who names the principal of a request for the audit log, or the
// sender a request says it is.
func (h *handler) who(r *http.Request) string {
	if v, ok := h.viewerOf(r); ok {
		return v.Label()
	}
	return "?"
}
