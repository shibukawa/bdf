package docflow

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shibukawa/yuruidb"
)

// Apps are the applications that send to the manager with credentials of
// their own (a business system, a WMS): an administrator registers one,
// with the scopes of what it may do, and gets its secret once. The app
// takes a short token at POST /oauth/token (OAuth 2.0 client credentials:
// Basic client_id:secret, or client_id and client_secret in the form) and
// sends with it as a bearer; over IPP it sends Basic client_id:secret.
// What it may send to are the entries that list it among their apps (an
// uploader, a CUPS server), and what it may read and work the downloaders
// and the inboxes that list it.
// The manager keeps the secret's argon2id hash only.

// Scopes of an app.
const (
	// ScopePrint sends documents to the entries that list the app.
	ScopePrint = "print"
	// ScopePrintTemplate sends data to those entries (for the template
	// nodes after them).
	ScopePrintTemplate = "print.template"
	// ScopeJobsCancel cancels the app's own jobs.
	ScopeJobsCancel = "jobs.cancel"
	// ScopeTick calls POST /api/tick (a scheduler).
	ScopeTick = "tick"
	// ScopeAdmin uses the API as an administrator (a script).
	ScopeAdmin = "admin"
	// ScopeDownload reads the downloaders that list the app
	// (download.go); of a person's token, the downloaders they may read.
	ScopeDownload = "download"
	// ScopeInbox works the inboxes that list the app (inboxapi.go): lists
	// what they hold, reads its pages and documents, releases and
	// discards; of a person's token, the inboxes they work.
	ScopeInbox = "inbox"
)

var appScopes = []string{ScopePrint, ScopePrintTemplate, ScopeJobsCancel, ScopeTick, ScopeAdmin, ScopeDownload, ScopeInbox}

// AppTokenFor is how long an app's token lasts; an app takes another
// with its credentials.
const AppTokenFor = 10 * time.Minute

// App is an application registered to send.
type App struct {
	// ID is the client_id: a name.
	ID    string `json:"id"`
	Title string `json:"title,omitempty"`
	// Scopes are what it may do (appScopes); print when none is given.
	Scopes []string `json:"scopes"`
	// Group is the group whose members see the app's jobs, when the
	// destination gives them none.
	Group string `json:"group,omitempty"`
	// Public tells a public client: a frontend in people's browsers,
	// which keeps no secret and gets tokens of the people using it with
	// the authorization code grant (oauth.go), sent back to one of its
	// RedirectURIs.
	Public       bool     `json:"public,omitempty"`
	RedirectURIs []string `json:"redirect_uris,omitempty"`
	// SecretHash is the argon2id hash of its secret, never the secret;
	// the API leaves it out. A public client has none.
	SecretHash string `json:"secret_hash,omitempty"`
	// Disabled refuses its credentials, keeping it.
	Disabled bool      `json:"disabled,omitempty"`
	Created  time.Time `json:"created"`
	LastUsed time.Time `json:"last_used,omitzero"`
}

// newAppSecret makes an app's secret: 32 random bytes, in base64url.
func newAppSecret() string {
	b := make([]byte, 32)
	rand.Read(b)
	return "bdfa_" + base64.RawURLEncoding.EncodeToString(b)
}

// PutApp registers an app or changes it: its scopes, title, group and
// whether it is disabled. A new app gets a secret, returned once; an app
// keeps its (RotateAppSecret makes a new one). A public client gets
// none: its redirect URIs are checked instead.
func (m *Manager) PutApp(ctx context.Context, in App) (secret string, err error) {
	if !validName(in.ID) {
		return "", fmt.Errorf("%w: app id %q", ErrInvalid, in.ID)
	}
	scopes := cleanNames(in.Scopes)
	if len(scopes) == 0 {
		scopes = []string{ScopePrint}
	}
	for _, s := range scopes {
		if !slices.Contains(appScopes, s) {
			return "", fmt.Errorf("%w: scope %q: %s", ErrInvalid, s, strings.Join(appScopes, ", "))
		}
	}
	uris, err := checkRedirectURIs(in.RedirectURIs)
	if err != nil {
		return "", err
	}
	if in.Public && len(uris) == 0 {
		return "", fmt.Errorf("%w: a public client has at least one redirect URI", ErrInvalid)
	}
	if !in.Public && len(uris) > 0 {
		return "", fmt.Errorf("%w: redirect URIs are a public client's", ErrInvalid)
	}
	a := App{ID: in.ID, Title: clip(strings.TrimSpace(in.Title), 255), Scopes: scopes, Group: strings.TrimSpace(in.Group), Public: in.Public, RedirectURIs: uris, Disabled: in.Disabled}
	var hash string
	_, err = change(ctx, m.st.db, a.ID, true, func(r *appRecord) error {
		switch {
		case a.Public:
			hash, secret = "", "" // no secret: a browser keeps none
			if r.Created.IsZero() {
				a.Created = m.now()
			} else {
				a.Created, a.LastUsed = r.Created, r.LastUsed
			}
		case r.SecretHash == "":
			if hash == "" {
				secret = newAppSecret()
				if hash, err = hashKey(secret); err != nil {
					return err
				}
			}
			a.Created = m.now()
			if !r.Created.IsZero() {
				a.Created, a.LastUsed = r.Created, r.LastUsed // was a public client
			}
		default:
			hash, secret = r.SecretHash, ""
			a.Created, a.LastUsed = r.Created, r.LastUsed
		}
		*r = appRecord(a)
		r.SecretHash = hash
		return nil
	})
	if err != nil {
		return "", err
	}
	return secret, nil
}

// App returns an app.
func (m *Manager) App(ctx context.Context, id string) (*App, error) {
	if !validName(id) {
		return nil, fmt.Errorf("%w: app %q", ErrNotFound, id)
	}
	rec, err := m.st.db.Get[appRecord](ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	a := App(rec.Data)
	a.ID = id
	return &a, nil
}

// Apps returns the apps, by id.
func (m *Manager) Apps(ctx context.Context) ([]App, error) {
	recs, err := list[appRecord](ctx, m.st.db)
	if err != nil {
		return nil, err
	}
	out := make([]App, len(recs))
	for i, r := range recs {
		out[i] = App(r.Data)
		out[i].ID = r.ID()
	}
	return out, nil
}

// DeleteApp removes an app; one a node lists (an entry, a downloader, an
// inbox) stays.
func (m *Manager) DeleteApp(ctx context.Context, id string) error {
	rec, err := m.st.db.Get[appRecord](ctx, id)
	if err != nil {
		return mapErr(err)
	}
	nodes, err := m.Nodes(ctx)
	if err != nil {
		return err
	}
	for _, n := range nodes {
		if slices.Contains(n.Apps, id) {
			return fmt.Errorf("%w: %s %s lists the app %s", ErrConflict, n.Kind, n.Name, id)
		}
	}
	return mapErr(m.st.db.Delete[appRecord](ctx, id, rec.Revision()))
}

// RotateAppSecret gives an app a new secret, returned once; the old one
// is refused from now.
func (m *Manager) RotateAppSecret(ctx context.Context, id string) (string, error) {
	secret := newAppSecret()
	hash, err := hashKey(secret)
	if err != nil {
		return "", err
	}
	if _, err := change(ctx, m.st.db, id, false, func(r *appRecord) error {
		if r.Public {
			return fmt.Errorf("%w: %s is a public client, which has no secret", ErrInvalid, id)
		}
		r.SecretHash = hash
		return nil
	}); err != nil {
		return "", err
	}
	m.o.Log.Info("app secret rotated", "app", id)
	return secret, nil
}

// appChecks remembers the credentials checked for a few minutes, in the
// instance: argon2id is slow, and IPP clients send with every request.
// The key is of the id, the secret and the hash, so that a rotation or
// another secret misses it.
var appChecks sync.Map

const appCheckFor = 5 * time.Minute

var (
	dummyHash     string
	dummyHashOnce sync.Once
)

// verifyApp tells whether a secret is the app's.
func (m *Manager) verifyApp(a *App, secret string) bool {
	if a == nil || a.SecretHash == "" || secret == "" {
		return false
	}
	sum := sha256.Sum256([]byte(a.ID + "\x00" + secret + "\x00" + a.SecretHash))
	key := hex.EncodeToString(sum[:])
	if until, ok := appChecks.Load(key); ok && m.now().Before(until.(time.Time)) {
		return true
	}
	if !checkKey(a.SecretHash, secret) {
		return false
	}
	appChecks.Store(key, m.now().Add(appCheckFor))
	return true
}

// burn spends the time a check takes, for an id that is not an app's,
// so that the time does not tell which ids are.
func burn(secret string) {
	dummyHashOnce.Do(func() { dummyHash, _ = hashKey("not an app") })
	checkKey(dummyHash, secret)
}

// appViewer is the viewer an app is, with the scopes granted.
func appViewer(a *App, scopes []string) *Viewer {
	return &Viewer{Kind: PrincipalApp, Via: ViaClientCredentials, ID: a.ID, Name: firstOr(a.Title, a.ID), Scopes: scopes, Admin: slices.Contains(scopes, ScopeAdmin)}
}

// IssueAppToken checks an app's credentials and makes a token of the
// scopes asked (none: every one of the app's), and tells when it ends.
func (m *Manager) IssueAppToken(ctx context.Context, id, secret string, scopes []string) (token string, expires time.Time, granted []string, err error) {
	a, err := m.App(ctx, id)
	if errors.Is(err, ErrNotFound) {
		burn(secret)
		return "", time.Time{}, nil, fmt.Errorf("%w: the client", ErrUnauthorized)
	}
	if err != nil {
		return "", time.Time{}, nil, err
	}
	if a.Disabled || !m.verifyApp(a, secret) {
		return "", time.Time{}, nil, fmt.Errorf("%w: the client", ErrUnauthorized)
	}
	granted = a.Scopes
	if len(scopes) > 0 {
		for _, s := range scopes {
			if !slices.Contains(a.Scopes, s) {
				return "", time.Time{}, nil, fmt.Errorf("%w: scope %q is not the app's", ErrInvalid, s)
			}
		}
		granted = cleanNames(scopes)
	}
	expires = m.now().Add(AppTokenFor)
	token = m.token("app", []byte(a.ID+"|"+strconv.FormatInt(expires.Unix(), 10)+"|"+strings.Join(granted, " ")))
	m.touchApp(ctx, a)
	return token, expires, granted, nil
}

// appTokenViewer is the app a token is of: one not expired, of an app
// that is there and not disabled, with the scopes granted that are still
// the app's.
func (m *Manager) appTokenViewer(ctx context.Context, tok string) (*Viewer, bool) {
	payload, ok := m.open("app", tok)
	if !ok {
		return nil, false
	}
	parts := strings.SplitN(string(payload), "|", 3)
	if len(parts) != 3 {
		return nil, false
	}
	exp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || m.now().Unix() > exp {
		return nil, false
	}
	a, err := m.App(ctx, parts[0])
	if err != nil || a.Disabled {
		return nil, false
	}
	var scopes []string
	for _, s := range strings.Fields(parts[2]) {
		if slices.Contains(a.Scopes, s) {
			scopes = append(scopes, s)
		}
	}
	m.touchApp(ctx, a)
	return appViewer(a, scopes), true
}

// appBasicViewer is the app of a Basic client_id:secret (IPP, and the
// API without a token).
func (m *Manager) appBasicViewer(ctx context.Context, id, secret string) (*Viewer, bool) {
	a, err := m.App(ctx, id)
	if err != nil {
		burn(secret)
		return nil, false
	}
	if a.Disabled || !m.verifyApp(a, secret) {
		return nil, false
	}
	m.touchApp(ctx, a)
	return appViewer(a, a.Scopes), true
}

// touchApp records when the app was last used, once a minute at most.
func (m *Manager) touchApp(ctx context.Context, a *App) {
	now := m.now()
	if now.Sub(a.LastUsed) < time.Minute {
		return
	}
	if _, err := change(ctx, m.st.db, a.ID, false, func(r *appRecord) error {
		if now.Sub(r.LastUsed) < time.Minute {
			return errSkip
		}
		r.LastUsed = now
		return nil
	}); err != nil && !errors.Is(err, ErrNotFound) {
		m.o.Log.Warn("app last used", "app", a.ID, "err", err)
	}
}

// appMaySend tells whether an app may send to a destination: an entry (an
// uploader, a CUPS server) that lists it, with the scope print (a
// document) or print.template (data). nil, another kind of viewer and an administrator may.
func (m *Manager) appMaySend(ctx context.Context, v *Viewer, dest string, data bool) error {
	if v == nil || v.Kind != PrincipalApp || v.Admin {
		return nil
	}
	need := ScopePrint
	if data {
		need = ScopePrintTemplate
	}
	n, err := m.Node(ctx, dest)
	if err != nil {
		return err
	}
	if !n.takesApp(v.ID) {
		return fmt.Errorf("%w: %s does not take the app %s", ErrForbidden, dest, v.ID)
	}
	if !v.Has(need) {
		return fmt.Errorf("%w: the app %s has no scope %s", ErrForbidden, v.ID, need)
	}
	return nil
}

// The apps' API (administrators), and the token endpoint (the apps).

func (h *handler) listApps(w http.ResponseWriter, r *http.Request) {
	as, err := h.m.Apps(r.Context())
	if err != nil {
		h.writeErr(w, err)
		return
	}
	for i := range as {
		as[i].SecretHash = ""
	}
	if as == nil {
		as = []App{}
	}
	writeJSON(w, http.StatusOK, as)
}

func (h *handler) getApp(w http.ResponseWriter, r *http.Request) {
	a, err := h.m.App(r.Context(), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	a.SecretHash = ""
	writeJSON(w, http.StatusOK, a)
}

// appWithSecret is an app and, once, its secret.
type appWithSecret struct {
	*App
	Secret string `json:"secret,omitempty"`
}

func (h *handler) putApp(w http.ResponseWriter, r *http.Request) {
	var a App
	if err := readJSON(r, &a); err != nil {
		h.writeErr(w, err)
		return
	}
	a.ID = r.PathValue("id")
	secret, err := h.m.PutApp(r.Context(), a)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	got, err := h.m.App(r.Context(), a.ID)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	got.SecretHash = ""
	detail := strings.Join(got.Scopes, " ")
	if secret != "" {
		detail += ", a secret made"
	}
	if got.Disabled {
		detail += ", disabled"
	}
	h.m.audit(r.Context(), h.who(r), "put-app", a.ID, detail)
	writeJSON(w, http.StatusOK, appWithSecret{got, secret})
}

func (h *handler) deleteApp(w http.ResponseWriter, r *http.Request) {
	if err := h.m.DeleteApp(r.Context(), r.PathValue("id")); err != nil {
		h.writeErr(w, err)
		return
	}
	h.m.audit(r.Context(), h.who(r), "delete-app", r.PathValue("id"), "")
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) rotateAppSecret(w http.ResponseWriter, r *http.Request) {
	secret, err := h.m.RotateAppSecret(r.Context(), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	h.m.audit(r.Context(), h.who(r), "rotate-app-secret", r.PathValue("id"), "")
	writeJSON(w, http.StatusOK, map[string]string{"id": r.PathValue("id"), "secret": secret})
}

// oauthToken is POST /oauth/token: the client credentials grant (RFC
// 6749 §4.4), the client in Basic or in the form, an optional scope; or
// the authorization code grant of a public client (oauth.go).
func (h *handler) oauthToken(w http.ResponseWriter, r *http.Request) {
	oauthErr := func(status int, code, desc string) {
		if status == http.StatusUnauthorized {
			w.Header().Set("WWW-Authenticate", `Basic realm="bdf docflow apps"`)
		}
		writeJSON(w, status, map[string]string{"error": code, "error_description": desc})
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<16)
	if err := r.ParseForm(); err != nil {
		oauthErr(http.StatusBadRequest, "invalid_request", "the form could not be read")
		return
	}
	switch r.PostForm.Get("grant_type") {
	case "client_credentials":
	case "authorization_code":
		h.oauthCodeGrant(w, r, oauthErr)
		return
	default:
		oauthErr(http.StatusBadRequest, "unsupported_grant_type", "client_credentials or authorization_code")
		return
	}
	id, secret, ok := r.BasicAuth()
	if !ok {
		id, secret = r.PostForm.Get("client_id"), r.PostForm.Get("client_secret")
	}
	if id == "" {
		oauthErr(http.StatusUnauthorized, "invalid_client", "the client's id and secret are required")
		return
	}
	tok, exp, granted, err := h.m.IssueAppToken(r.Context(), id, secret, strings.Fields(r.PostForm.Get("scope")))
	switch {
	case errors.Is(err, ErrUnauthorized):
		h.m.audit(r.Context(), "app:"+clip(id, 60), "token-refused", id, "")
		oauthErr(http.StatusUnauthorized, "invalid_client", "the client's credentials")
		return
	case errors.Is(err, ErrInvalid):
		oauthErr(http.StatusBadRequest, "invalid_scope", err.Error())
		return
	case err != nil:
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"access_token": tok, "token_type": "Bearer", "expires_in": int(exp.Sub(h.m.now()).Seconds()), "scope": strings.Join(granted, " ")})
}

var _ = yuruidb.ErrNotFound
