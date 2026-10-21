package manager

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// People sign in on the pages: an administrator with the admin token, a
// user with a submit token's name and secret. The session is a cookie
// the manager signs (as the devices' is), so that any instance reads it.
// This is the provisional scheme of pm-1's tokens; the authentication of
// people over OIDC, and the roles of administrators (the topology, the
// audit) are requirement:print-manager-auth.

// UserCookie is the cookie of a person's session.
const UserCookie = "bdf_pm_user"

// Principal is who a request is from: an administrator, or a user of a
// name (a submit token's).
type Principal struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
}

// Login checks a token: the admin token, or a submit token as
// name:secret.
func (m *Manager) Login(token string) (Principal, bool) {
	token = strings.TrimSpace(token)
	if sameSecret(m.o.AdminToken, token) {
		return Principal{Name: "admin", Admin: true}, true
	}
	name, secret, ok := strings.Cut(token, ":")
	if !ok {
		return Principal{}, false
	}
	want, known := m.o.SubmitTokens[name]
	if !known || !sameSecret(want, secret) {
		return Principal{}, false
	}
	return Principal{Name: name}, true
}

// UserSessionFor is how long a person's session lasts.
const UserSessionFor = 12 * time.Hour

// userSession makes the session of a principal.
func (m *Manager) userSession(p Principal) (string, time.Time) {
	until := m.now().Add(UserSessionFor)
	admin := "0"
	if p.Admin {
		admin = "1"
	}
	return m.token("user", []byte(p.Name+"|"+admin+"|"+strconv.FormatInt(until.Unix(), 10))), until
}

// UserSession reads a session: one not expired, of a user still known
// (the admin, or a submit token's name).
func (m *Manager) UserSession(tok string) (Principal, bool) {
	payload, ok := m.open("user", tok)
	if !ok {
		return Principal{}, false
	}
	parts := strings.Split(string(payload), "|")
	if len(parts) != 3 {
		return Principal{}, false
	}
	until, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || m.now().Unix() > until {
		return Principal{}, false
	}
	p := Principal{Name: parts[0], Admin: parts[1] == "1"}
	if p.Admin {
		return p, m.o.AdminToken != ""
	}
	_, known := m.o.SubmitTokens[p.Name]
	return p, known
}

// PrincipalOf finds who a request is from: the admin token (bearer), a
// submit token (bearer name:secret, or Basic), or a person's session
// (cookie).
func (m *Manager) PrincipalOf(r *http.Request) (Principal, bool) {
	if tok := bearer(r); tok != "" {
		if p, ok := m.Login(tok); ok {
			return p, true
		}
	}
	if name, secret, ok := r.BasicAuth(); ok {
		if p, ok := m.Login(name + ":" + secret); ok {
			return p, true
		}
	}
	if c, err := r.Cookie(UserCookie); err == nil {
		return m.UserSession(c.Value)
	}
	return Principal{}, false
}

// login signs a person in: a form (token=…, then a redirect to /jobs,
// or back to / with ?error=1) or JSON ({"token": …}, then the
// principal).
func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var token string
	isJSON := strings.HasPrefix(r.Header.Get("Content-Type"), "application/json")
	if isJSON {
		var in struct {
			Token string `json:"token"`
		}
		if err := readJSON(r, &in); err != nil {
			h.writeErr(w, err)
			return
		}
		token = in.Token
	} else {
		if err := r.ParseForm(); err != nil {
			h.writeErr(w, errors.Join(ErrInvalid, err))
			return
		}
		token = r.PostForm.Get("token")
	}
	p, ok := h.m.Login(token)
	if !ok {
		name, _, _ := strings.Cut(token, ":")
		h.m.audit(r.Context(), clip(firstOr(name, "?"), 80), "login-failed", "", "")
		if isJSON {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "the token is not known"})
			return
		}
		http.Redirect(w, r, "/?error=1", http.StatusSeeOther)
		return
	}
	tok, until := h.m.userSession(p)
	http.SetCookie(w, &http.Cookie{Name: UserCookie, Value: tok, Path: "/", MaxAge: int(UserSessionFor / time.Second), HttpOnly: true,
		SameSite: http.SameSiteLaxMode, Secure: h.m.o.SecureCookies})
	h.m.audit(r.Context(), p.Name, "login", "", "")
	if isJSON {
		writeJSON(w, http.StatusOK, struct {
			Principal
			Until time.Time `json:"until"`
		}{p, until})
		return
	}
	http.Redirect(w, r, "/jobs", http.StatusSeeOther)
}

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	if p, ok := h.m.PrincipalOf(r); ok {
		h.m.audit(r.Context(), p.Name, "logout", "", "")
	}
	http.SetCookie(w, &http.Cookie{Name: UserCookie, Value: "", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: h.m.o.SecureCookies})
	if strings.HasPrefix(r.Header.Get("Accept"), "application/json") {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// whoami answers the principal of a request.
func (h *handler) whoami(w http.ResponseWriter, r *http.Request) {
	p, ok := h.m.PrincipalOf(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "not signed in"})
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// userOnly serves a request of a signed-in person.
func (h *handler) userOnly(fn func(http.ResponseWriter, *http.Request, Principal)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := h.m.PrincipalOf(r)
		if !ok {
			w.Header().Set("WWW-Authenticate", `Bearer realm="bdf print manager"`)
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "not signed in"})
			return
		}
		fn(w, r, p)
	}
}

// ownJob serves a request about a job of the person's: any job for an
// administrator, the jobs sent under their name for a user.
func (h *handler) ownJob(fn func(http.ResponseWriter, *http.Request, Principal, *Job)) http.HandlerFunc {
	return h.userOnly(func(w http.ResponseWriter, r *http.Request, p Principal) {
		j, err := h.m.Job(r.Context(), r.PathValue("id"))
		if err == nil && !p.Admin && j.User != p.Name {
			err = ErrNotFound
		}
		if err != nil {
			h.writeErr(w, err)
			return
		}
		fn(w, r, p, j)
	})
}

// allJobs is the destination state that lists every job.
const allJobs = "-all"

// MyJobs returns the last jobs of a person, latest first: every job for
// an administrator, those sent under their name for a user.
func (m *Manager) MyJobs(ctx context.Context, p Principal, limit int) ([]Job, error) {
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
		if p.Admin || j.User == p.Name {
			out = append(out, j)
		}
	}
	return out, nil
}

func (h *handler) myJobs(w http.ResponseWriter, r *http.Request, p Principal) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	js, err := h.m.MyJobs(r.Context(), p, limit)
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
}

// destinations lists the nodes a person may print to: every node, for
// now (the scopes of requirement:print-manager-auth narrow it later).
func (h *handler) destinations(w http.ResponseWriter, r *http.Request, p Principal) {
	topo, err := h.m.Topology(r.Context())
	if err != nil {
		h.writeErr(w, err)
		return
	}
	out := make([]Destination, len(topo))
	for i, n := range topo {
		out[i] = Destination{Name: n.Name, Kind: n.Kind, Title: n.Title}
	}
	writeJSON(w, http.StatusOK, out)
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
	if p, ok := h.m.PrincipalOf(r); ok {
		return p.Name
	}
	return "?"
}
