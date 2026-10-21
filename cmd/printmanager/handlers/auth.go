package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/shibukawa/bdf/print/manager"
	"github.com/shibukawa/popcornweb/authstate"
	"github.com/shibukawa/popcornweb/plugin/auth"
	"github.com/shibukawa/popcornweb/pw"
	"github.com/shibukawa/popcornweb/sessionconfig"
)

// People sign in to the pages with OIDC (authorization code and PKCE) at
// an external IdP, through popcornweb's plugin/auth: anyone it signs in
// gets in (OIDCConfig.Allowed narrows that), and sees what their scope
// lets them (print/manager's Viewer). Nothing is kept in an instance: the
// state of a login on its way through the IdP is in the manager's store
// (any instance finishes it), and the session is a sealed cookie. The
// groups the IdP gives are copied at each login to the person's record;
// who is an admin is decided on each request from that record and the
// configuration.

// authBackend is plugin/auth's storage of the state of logins: the
// manager's store.
const authBackend = "yuruidb"

// authStore is the store of the open App, for the backend that plugin/auth
// opens when the server starts.
var authStore atomic.Pointer[manager.Store]

func init() {
	auth.RegisterBackend(authBackend, func(context.Context, auth.Config, auth.Resources) (auth.Backend, error) {
		st := authStore.Load()
		if st == nil {
			return auth.Backend{}, errors.New("printmanager: the store is not open")
		}
		return auth.Backend{OpenState: func(_ context.Context, namespace string) (authstate.RawStore, error) {
			return loginStates{st: st, ns: namespace}, nil
		}}, nil
	})
}

// loginStates keeps plugin/auth's state of logins (state, nonce, PKCE
// verifier) in the manager's store.
type loginStates struct {
	st *manager.Store
	ns string
}

func (s loginStates) Put(ctx context.Context, key string, payload []byte, expires time.Time) error {
	switch err := s.st.PutState(ctx, s.ns+"\x00"+key, payload, expires); {
	case errors.Is(err, manager.ErrConflict):
		return authstate.ErrAlreadyExists
	case err != nil:
		return fmt.Errorf("%w: %v", authstate.ErrUnavailable, err)
	}
	return nil
}

func (s loginStates) Take(ctx context.Context, key string) ([]byte, error) {
	v, err := s.st.TakeState(ctx, s.ns+"\x00"+key, time.Now())
	switch {
	case errors.Is(err, manager.ErrNotFound):
		return nil, authstate.ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("%w: %v", authstate.ErrUnavailable, err)
	}
	return v, nil
}

// Prune removes the states of logins nobody finished.
func (s loginStates) Prune(ctx context.Context, before time.Time, limit int) (int64, error) {
	n, err := s.st.PruneStates(ctx, before, limit)
	return int64(n), err
}

// The paths that are not the pages: the top page is the sign-in (a
// button to the IdP), the devices' page and endpoints keep
// their own sign-in (the device's key), the API and the printers take
// their tokens (and the API the session too), the endpoints their keys.
var publicPaths = []string{"/", "/tablet", "/public/**", "/api/**", manager.IPPPath + "/**", manager.InPath + "/**", "/device/**"}

// ConfigureAuth sets popcornweb's [auth] and [session] for the pages:
// what is the Print Manager's design (the OIDC mode, the paths a person
// signs in for, where the state of logins is kept, a sealed cookie as the
// session) rather than a deployment's choice. Outside dev it refuses to
// go on without the IdP or an admin rule. secret is MANAGER_SECRET: the
// session's keys are made of it when session.keyring.secret is not set.
//
// Without an IdP in dev, nobody signs in: the pages say so.
func ConfigureAuth(a *auth.Config, s *sessionconfig.SessionConfig, secret string, secureCookies bool, o OIDCConfig) error {
	o2 := a.OIDC
	if o2.Issuer == "" && pw.Development() {
		a.Enabled = false
		return nil
	}
	var missing []string
	for _, f := range []struct{ v, name string }{{o2.Issuer, "auth.oidc.issuer (AUTH_OIDC_ISSUER)"}, {o2.ClientID, "auth.oidc.client_id (AUTH_OIDC_CLIENT_ID)"},
		{o2.ClientSecret, "auth.oidc.client_secret (AUTH_OIDC_CLIENT_SECRET)"}} {
		if f.v == "" {
			missing = append(missing, f.name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("printmanager: people sign in at an OIDC IdP: %s required", strings.Join(missing, ", "))
	}
	if len(o.Admins) == 0 && len(o.AdminGroups) == 0 && !pw.Development() {
		return errors.New("printmanager: oidc.admins (OIDC_ADMINS) or oidc.admin_groups (OIDC_ADMIN_GROUPS) is required: nobody would be an admin")
	}
	for _, e := range o.Admins {
		if iss, sub, ok := strings.Cut(e, "#"); ok && (iss == "" || sub == "") || !ok && !strings.Contains(e, "@") {
			return fmt.Errorf("printmanager: oidc.admins %q: issuer#subject or an email", e)
		}
	}
	a.Enabled = true
	a.Mode = auth.ModeOIDCOnly
	a.Backend = authBackend
	a.Protection.Include = []string{"/**"}
	for _, p := range publicPaths {
		if !slices.Contains(a.Protection.Exclude, p) {
			a.Protection.Exclude = append(a.Protection.Exclude, p)
		}
	}
	a.Protection.Unauthenticated = auth.UnauthenticatedRedirect
	if len(a.OIDC.Scopes) == 0 {
		a.OIDC.Scopes = []string{"email", "profile"}
	}
	if u, err := url.Parse(o2.Issuer); err == nil && u.Scheme == "http" && pw.Development() && loopback(u.Hostname()) {
		a.OIDC.AllowLoopbackHTTP = true
	}

	s.Enabled = true
	if s.Backend == "" || s.Backend == sessionconfig.SessionBackendRDB {
		// the default: the manager has no database of sessions
		s.Backend = sessionconfig.SessionBackendCookie
	}
	if s.Keyring.Secret == "" {
		key := make([]byte, 32)
		if secret == "" {
			rand.Read(key) // dev: the sessions end with the process
		} else {
			mac := hmac.New(sha256.New, []byte(secret))
			mac.Write([]byte("bdf print manager: session keyring"))
			key = mac.Sum(nil)
		}
		s.Keyring.Secret = base64.StdEncoding.EncodeToString(key)
	}
	if !secureCookies {
		s.Cookie.Secure = false
	}
	return nil
}

func loopback(host string) bool {
	return host == "localhost" || host == "127.0.0.1" || host == "::1" || strings.HasSuffix(host, ".localhost")
}

// policy is who gets in, their groups, and who is an admin.
type policy struct{ OIDCConfig }

// groups reads the groups of a person from the claims (claim gives a
// claim's JSON): a string or an array of strings at the claim's dotted
// path.
func (p policy) groups(claim func(name string) (json.RawMessage, bool)) []string {
	path := strings.Split(p.GroupsClaim, ".")
	raw, ok := claim(path[0])
	if !ok {
		return nil
	}
	var v any
	if json.Unmarshal(raw, &v) != nil {
		return nil
	}
	for _, k := range path[1:] {
		obj, ok := v.(map[string]any)
		if !ok {
			return nil
		}
		v = obj[k]
	}
	var out []string
	add := func(g any) {
		if s, ok := g.(string); ok && s != "" && strings.HasPrefix(s, p.GroupsPrefix) && !slices.Contains(out, s) {
			out = append(out, s)
		}
	}
	switch v := v.(type) {
	case string:
		add(v)
	case []any:
		for _, g := range v {
			add(g)
		}
	}
	return out
}

// email is the person's email when the IdP does not say it is
// unverified.
func email(c auth.Claims) string {
	e, _ := c.String("email")
	if raw, ok := c.Raw("email_verified"); ok && string(raw) == "false" {
		return ""
	}
	return e
}

func domainOf(email string) string {
	_, d, _ := strings.Cut(email, "@")
	return strings.ToLower(d)
}

// admin tells an admin: by issuer#subject, email, or a group.
func (p policy) admin(iss, sub, email string, groups []string) bool {
	for _, e := range p.Admins {
		if i, s, ok := strings.Cut(e, "#"); ok {
			if i == iss && s == sub {
				return true
			}
		} else if email != "" && strings.EqualFold(e, email) {
			return true
		}
	}
	for _, g := range p.AdminGroups {
		if slices.Contains(groups, g) {
			return true
		}
	}
	return false
}

// allowed tells whether a person gets in.
func (p policy) allowed(iss, sub, email, hd string, groups []string) bool {
	if p.admin(iss, sub, email, groups) {
		return true
	}
	a := p.Allowed
	if len(a.Issuers) > 0 && !slices.Contains(a.Issuers, iss) {
		return false
	}
	if len(a.Domains) > 0 && !slices.ContainsFunc(a.Domains, func(d string) bool {
		d = strings.ToLower(d)
		return d == domainOf(email) && email != "" || hd != "" && d == strings.ToLower(hd)
	}) {
		return false
	}
	if len(a.Groups) > 0 && !slices.ContainsFunc(a.Groups, func(g string) bool { return slices.Contains(groups, g) }) {
		return false
	}
	return true
}

// personID is the ID of a person: the same for the same issuer and
// subject, and nothing else of them.
func personID(iss, key string) string {
	sum := sha256.Sum256([]byte(iss + "\x00" + key))
	return "u" + hex.EncodeToString(sum[:16])
}

// resolveAccount is plugin/auth's account resolver: a person the IdP
// signed in gets in when the policy lets them, and their record is
// written with the groups of this login.
func (a *App) resolveAccount(ctx context.Context, id auth.Identity, _ bool) (auth.Account, error) {
	c := id.Claims
	groups := a.policy.groups(c.Raw)
	mail := email(c)
	hd, _ := c.String("hd")
	if !a.policy.allowed(id.Issuer, id.Subject, mail, hd, groups) {
		a.log.Info("sign-in refused", "iss", id.Issuer, "sub", id.Subject)
		a.m.Audit(ctx, firstOr(mail, id.Subject), "login-refused", "", id.Issuer)
		return auth.Account{}, auth.ErrAccessDenied
	}
	name, _ := c.String("name")
	if name == "" {
		name, _ = c.String("preferred_username")
	}
	pid := personID(id.Issuer, id.Key)
	if err := a.m.PutPerson(ctx, pid, manager.Person{Issuer: id.Issuer, Subject: id.Subject, Name: name, Email: mail,
		Groups: groups, Login: time.Now().UTC()}); err != nil {
		return auth.Account{}, err
	}
	a.log.Info("sign-in", "person", pid, "iss", id.Issuer, "sub", id.Subject, "groups", len(groups))
	a.m.Audit(ctx, firstOr(name, firstOr(mail, pid)), "login", pid, id.Issuer)
	return auth.Account{ID: pid, DisplayName: name, Email: mail}, nil
}

// viewer is the person signed in to the request's session, or nil. A
// session whose person's record is gone (a store emptied, or replaced)
// is still the person: the session's name and email, without groups
// until the next login.
func (a *App) viewer(r *http.Request) *manager.Viewer {
	s, ok := auth.User(r.Context())
	if !ok {
		return nil
	}
	p, err := a.m.Person(r.Context(), s.AccountID)
	if err != nil {
		if !errors.Is(err, manager.ErrNotFound) {
			a.log.Warn("person", "id", s.AccountID, "err", err)
		}
		p = &manager.Person{Issuer: s.Issuer, Subject: s.Subject, Name: s.DisplayName, Email: s.Email}
	}
	return &manager.Viewer{Kind: manager.PrincipalUser, Via: manager.ViaOIDC, ID: s.AccountID, Name: p.Name, Email: p.Email, Groups: p.Groups,
		Admin: a.policy.admin(p.Issuer, p.Subject, p.Email, p.Groups)}
}

// statusWriter answers what a page writes with another status.
type statusWriter struct {
	http.ResponseWriter
	status int
	wrote  bool
}

func (w *statusWriter) WriteHeader(int) {
	if !w.wrote {
		w.wrote = true
		w.ResponseWriter.WriteHeader(w.status)
	}
}

func (w *statusWriter) Write(b []byte) (int, error) {
	w.WriteHeader(0)
	return w.ResponseWriter.Write(b)
}

func (w *statusWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }
