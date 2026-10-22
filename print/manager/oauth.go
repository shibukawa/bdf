package manager

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Frontends: a web application of the organization's own, running in
// people's browsers on another origin, reads an export for the person
// using it with a token of theirs. The frontend is registered as a
// public client (an App with Public and its RedirectURIs: no secret, as
// a browser keeps none), and gets the token with the authorization code
// grant and PKCE (RFC 6749 §4.1, RFC 7636) at the manager, which is the
// authorization server of its machines already (apps.go): the person
// signs in at the IdP as for the pages, GET /oauth/authorize sends a
// code back to the frontend, and POST /oauth/token exchanges it, with
// the verifier, for a token of UserTokenFor. The administrator's
// registration of the client is the consent: no screen asks the person.
// The token carries the person (their record is read on each request:
// the groups of their last login, the subject at the IdP) and the scopes
// granted, no role: it opens the export endpoints and /api/me, nothing
// else of the API. The frontend asks for a code again when it expires
// (the person's session at the manager makes that silent).

// OAuthAuthorizePath is where a frontend sends the person for a code.
const OAuthAuthorizePath = "/oauth/authorize"

const (
	// UserTokenFor is how long a person's token lasts.
	UserTokenFor = time.Hour
	// codeFor is how long a code may be exchanged.
	codeFor = 5 * time.Minute
	// maxRedirectURIs bounds a public client's redirect URIs.
	maxRedirectURIs = 8
)

// authCode is a code on its way back to the frontend, kept until it is
// exchanged (once) or expires.
type authCode struct {
	Client    string   `json:"client"`
	Redirect  string   `json:"redirect"`
	Challenge string   `json:"challenge"`
	Person    string   `json:"person"`
	Scopes    []string `json:"scopes"`
}

func codeKey(code string) string { return "oauth-code\x00" + code }

// checkRedirectURIs cleans a public client's redirect URIs: absolute,
// https (http on the loopback), without a fragment.
func checkRedirectURIs(in []string) ([]string, error) {
	uris := cleanNames(in)
	if len(uris) > maxRedirectURIs {
		return nil, fmt.Errorf("%w: %d redirect URIs; at most %d", ErrInvalid, len(uris), maxRedirectURIs)
	}
	for _, s := range uris {
		u, err := url.Parse(s)
		if err != nil || !u.IsAbs() || u.Host == "" || u.Fragment != "" || len(s) > 255 {
			return nil, fmt.Errorf("%w: redirect URI %q: an absolute https URL without a fragment", ErrInvalid, s)
		}
		switch u.Scheme {
		case "https":
		case "http":
			if !loopbackHost(u.Hostname()) {
				return nil, fmt.Errorf("%w: redirect URI %q: http only on the loopback (localhost)", ErrInvalid, s)
			}
		default:
			return nil, fmt.Errorf("%w: redirect URI %q: https", ErrInvalid, s)
		}
	}
	return uris, nil
}

func loopbackHost(host string) bool {
	return host == "localhost" || host == "127.0.0.1" || host == "::1" || strings.HasSuffix(host, ".localhost")
}

// clientOrigin tells an origin of a public client: of one of its
// redirect URIs. A page there may call the token endpoint and the
// exports (CORS).
func (m *Manager) clientOrigin(ctx context.Context, origin string) bool {
	o, err := url.Parse(origin)
	if err != nil || o.Host == "" {
		return false
	}
	apps, err := m.Apps(ctx)
	if err != nil {
		return false
	}
	for _, a := range apps {
		if !a.Public || a.Disabled {
			continue
		}
		for _, s := range a.RedirectURIs {
			if u, err := url.Parse(s); err == nil && strings.EqualFold(u.Scheme, o.Scheme) && strings.EqualFold(u.Host, o.Host) {
				return true
			}
		}
	}
	return false
}

// cors lets a page of a public client's origin call: the token endpoint
// and the exports answer its preflight, and name the origin.
func (h *handler) cors(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" && h.m.clientOrigin(r.Context(), origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition, Content-Length")
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
				w.Header().Set("Access-Control-Max-Age", "600")
			}
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		fn(w, r)
	}
}

// oauthAuthorize is GET /oauth/authorize: a person signed in to the
// pages is sent back to the public client with a code. A client or a
// redirect URI that is not registered is refused here; the other errors
// go back to the client (RFC 6749 §4.1.2.1).
func (h *handler) oauthAuthorize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()
	a, err := h.m.App(ctx, q.Get("client_id"))
	if err != nil || !a.Public || a.Disabled {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_request", "error_description": "the client is not a public client registered here"})
		return
	}
	redirect := q.Get("redirect_uri")
	if !slices.Contains(a.RedirectURIs, redirect) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_request", "error_description": "the redirect URI is not the client's"})
		return
	}
	back := func(params map[string]string) {
		u, _ := url.Parse(redirect)
		v := u.Query()
		for k, val := range params {
			v.Set(k, val)
		}
		if state := q.Get("state"); state != "" {
			v.Set("state", state)
		}
		u.RawQuery = v.Encode()
		w.Header().Set("Cache-Control", "no-store")
		http.Redirect(w, r, u.String(), http.StatusFound)
	}
	fail := func(code, desc string) { back(map[string]string{"error": code, "error_description": desc}) }
	if q.Get("response_type") != "code" {
		fail("unsupported_response_type", "code only")
		return
	}
	challenge := q.Get("code_challenge")
	if challenge == "" || q.Get("code_challenge_method") != "S256" {
		fail("invalid_request", "PKCE with S256 is required (code_challenge, code_challenge_method)")
		return
	}
	scopes := a.Scopes
	if asked := strings.Fields(q.Get("scope")); len(asked) > 0 {
		for _, s := range asked {
			if !slices.Contains(a.Scopes, s) {
				fail("invalid_scope", "the scope "+s+" is not the client's")
				return
			}
		}
		scopes = cleanNames(asked)
	}
	v, ok := h.viewerOf(r)
	if !ok || v.Kind != PrincipalUser || v.Token() {
		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "sign in first"})
		return
	}
	code := make([]byte, 32)
	rand.Read(code)
	c := authCode{Client: a.ID, Redirect: redirect, Challenge: challenge, Person: v.ID, Scopes: scopes}
	b, _ := json.Marshal(c)
	s := base64.RawURLEncoding.EncodeToString(code)
	if err := h.m.st.PutState(ctx, codeKey(s), b, h.m.now().Add(codeFor)); err != nil {
		h.writeErr(w, err)
		return
	}
	h.m.audit(ctx, v.Label(), "oauth-code", a.ID, strings.Join(scopes, " "))
	back(map[string]string{"code": s})
}

// RedeemCode exchanges a code for a person's token: the code is taken
// once, and is the client's, for the redirect URI it was sent to, with
// the verifier of its challenge.
func (m *Manager) RedeemCode(ctx context.Context, code, verifier, client, redirect string) (token string, expires time.Time, scopes []string, err error) {
	refused := fmt.Errorf("%w: the code", ErrUnauthorized)
	raw, err := m.st.TakeState(ctx, codeKey(code), m.now())
	if errors.Is(err, ErrNotFound) {
		return "", time.Time{}, nil, refused
	}
	if err != nil {
		return "", time.Time{}, nil, err
	}
	var c authCode
	if json.Unmarshal(raw, &c) != nil || c.Client != client || c.Redirect != redirect {
		return "", time.Time{}, nil, refused
	}
	sum := sha256.Sum256([]byte(verifier))
	if verifier == "" || base64.RawURLEncoding.EncodeToString(sum[:]) != c.Challenge {
		return "", time.Time{}, nil, refused
	}
	a, err := m.App(ctx, client)
	if err != nil || !a.Public || a.Disabled {
		return "", time.Time{}, nil, refused
	}
	for _, s := range c.Scopes {
		if slices.Contains(a.Scopes, s) {
			scopes = append(scopes, s)
		}
	}
	expires = m.now().Add(UserTokenFor)
	token = m.token("user", []byte(c.Person+"|"+strconv.FormatInt(expires.Unix(), 10)+"|"+strings.Join(scopes, " ")+"|"+a.ID))
	m.touchApp(ctx, a)
	who := c.Person
	if p, err := m.Person(ctx, c.Person); err == nil {
		who = firstOr(p.Name, firstOr(p.Email, c.Person))
	}
	m.audit(ctx, who, "user-token", a.ID, strings.Join(scopes, " "))
	return token, expires, scopes, nil
}

// userTokenViewer is the person a token is of: one not expired, of a
// person still kept, through a client still there, with the scopes
// granted that are still the client's.
func (m *Manager) userTokenViewer(ctx context.Context, tok string) (*Viewer, bool) {
	payload, ok := m.open("user", tok)
	if !ok {
		return nil, false
	}
	parts := strings.SplitN(string(payload), "|", 4)
	if len(parts) != 4 {
		return nil, false
	}
	exp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || m.now().Unix() > exp {
		return nil, false
	}
	p, err := m.Person(ctx, parts[0])
	if err != nil {
		return nil, false
	}
	a, err := m.App(ctx, parts[3])
	if err != nil || a.Disabled || !a.Public {
		return nil, false
	}
	var scopes []string
	for _, s := range strings.Fields(parts[2]) {
		if slices.Contains(a.Scopes, s) {
			scopes = append(scopes, s)
		}
	}
	return &Viewer{Kind: PrincipalUser, Via: ViaOAuth, ID: parts[0], Name: p.Name, Email: p.Email, Groups: p.Groups, Subject: p.Subject, Scopes: scopes}, true
}

// oauthCodeGrant is POST /oauth/token with grant_type=authorization_code:
// the public client, named in the form (or as the Basic user, with no
// password), sends the code, the verifier and the redirect URI.
func (h *handler) oauthCodeGrant(w http.ResponseWriter, r *http.Request, oauthErr func(status int, code, desc string)) {
	client := r.PostForm.Get("client_id")
	if id, _, ok := r.BasicAuth(); ok && client == "" {
		client = id
	}
	code, verifier, redirect := r.PostForm.Get("code"), r.PostForm.Get("code_verifier"), r.PostForm.Get("redirect_uri")
	if client == "" || code == "" || verifier == "" || redirect == "" {
		oauthErr(http.StatusBadRequest, "invalid_request", "client_id, code, code_verifier and redirect_uri are required")
		return
	}
	tok, exp, granted, err := h.m.RedeemCode(r.Context(), code, verifier, client, redirect)
	switch {
	case errors.Is(err, ErrUnauthorized):
		h.m.audit(r.Context(), "app:"+clip(client, 60), "code-refused", client, "")
		oauthErr(http.StatusBadRequest, "invalid_grant", "the code, the verifier or the redirect URI")
		return
	case err != nil:
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"access_token": tok, "token_type": "Bearer", "expires_in": int(exp.Sub(h.m.now()).Seconds()), "scope": strings.Join(granted, " ")})
}
