// Package fakeidp is an OpenID Provider for tests: discovery, a JWKS, the
// authorization code flow with PKCE (S256) and RS256 ID tokens, with the
// people of a test and their groups in a claim of any path
// ("groups", "roles", Keycloak's "realm_access.roles"). It checks no
// password: a test says who signs in next (LoginAs), or a person picks
// themselves on its page. It is to Docflow's sign-in what
// print/fakeprinter is to its printers.
//
//	idp := fakeidp.New(fakeidp.User{Subject: "alice", Email: "alice@example.com", Groups: []string{"warehouse-a"}})
//	issuer := fakeidp.Start(t, idp)
//	idp.LoginAs("alice")
package fakeidp

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"
)

// Default credentials of the one client the provider knows.
const (
	ClientID     = "docflow"
	ClientSecret = "fake-client-secret"
)

// User is a person of the provider.
type User struct {
	Subject string
	Name    string
	Email   string
	// EmailUnverified sends email_verified false.
	EmailUnverified bool
	// Groups go to the provider's GroupsClaim.
	Groups []string
	// Claims are more claims of the ID token (they win over the others).
	Claims map[string]any
}

// IdP is the provider. Its fields may be changed between logins.
type IdP struct {
	// Issuer is the provider's URL; Start and Listen set it. Empty, it is
	// http:// and the request's host.
	Issuer       string
	ClientID     string
	ClientSecret string
	// GroupsClaim is where the groups go: a claim's name, or a dotted
	// path into objects (realm_access.roles). Default "groups".
	GroupsClaim string
	// TokenTTL is how long an ID token lasts (default 5 minutes).
	TokenTTL time.Duration
	// Tamper changes the claims of the next ID tokens, to test what a
	// relying party refuses.
	Tamper func(claims map[string]any)
	// Now is the clock (default time.Now).
	Now func() time.Time

	mu    sync.Mutex
	users []User
	login string
	codes map[string]*grant
	key   *rsa.PrivateKey
	// Logins counts the codes the provider gave.
	logins int
}

type grant struct {
	user                       User
	clientID, redirect, nonce  string
	challenge, challengeMethod string
	expires                    time.Time
}

// New makes a provider of users, with the default client.
func New(users ...User) *IdP {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return &IdP{ClientID: ClientID, ClientSecret: ClientSecret, users: users, codes: map[string]*grant{}, key: key}
}

// AddUser adds a person, or replaces the one of the same subject.
func (p *IdP) AddUser(u User) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := range p.users {
		if p.users[i].Subject == u.Subject {
			p.users[i] = u
			return
		}
	}
	p.users = append(p.users, u)
}

// LoginAs makes the next authorizations sign in a person without a page;
// "" shows the page where a person picks themselves.
func (p *IdP) LoginAs(subject string) {
	p.mu.Lock()
	p.login = subject
	p.mu.Unlock()
}

// Logins tells how many codes the provider gave.
func (p *IdP) Logins() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.logins
}

// Start serves the provider on the loopback until the test ends, and
// returns its issuer.
func Start(t testing.TB, p *IdP) string {
	t.Helper()
	s, err := p.Listen("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return p.Issuer
}

// Server is the provider served on a port.
type Server struct {
	URL string
	srv *http.Server
}

// Close stops the server.
func (s *Server) Close() error { return s.srv.Close() }

// Listen serves the provider at addr over http, and sets its Issuer.
func (p *IdP) Listen(addr string) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	u := "http://" + ln.Addr().String()
	p.mu.Lock()
	p.Issuer = u
	p.mu.Unlock()
	srv := &http.Server{Handler: p, ReadHeaderTimeout: 10 * time.Second}
	go srv.Serve(ln)
	return &Server{URL: u, srv: srv}, nil
}

func (p *IdP) now() time.Time {
	if p.Now != nil {
		return p.Now()
	}
	return time.Now()
}

func (p *IdP) issuer(r *http.Request) string {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Issuer != "" {
		return p.Issuer
	}
	return "http://" + r.Host
}

// ServeHTTP serves discovery, the keys, authorization, the token
// endpoint, and the end of a session.
func (p *IdP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/.well-known/openid-configuration":
		iss := p.issuer(r)
		writeJSON(w, http.StatusOK, map[string]any{
			"issuer":                                iss,
			"authorization_endpoint":                iss + "/authorize",
			"token_endpoint":                        iss + "/token",
			"jwks_uri":                              iss + "/jwks",
			"end_session_endpoint":                  iss + "/logout",
			"response_types_supported":              []string{"code"},
			"subject_types_supported":               []string{"public"},
			"id_token_signing_alg_values_supported": []string{"RS256"},
			"code_challenge_methods_supported":      []string{"S256"},
			"token_endpoint_auth_methods_supported": []string{"client_secret_basic", "client_secret_post"},
			"scopes_supported":                      []string{"openid", "email", "profile"},
		})
	case "/jwks":
		pub := p.key.PublicKey
		writeJSON(w, http.StatusOK, map[string]any{"keys": []map[string]string{{
			"kty": "RSA", "kid": "fake-1", "alg": "RS256", "use": "sig",
			"n": b64(pub.N.Bytes()), "e": b64(big.NewInt(int64(pub.E)).Bytes()),
		}}})
	case "/authorize":
		p.authorize(w, r)
	case "/token":
		p.token(w, r)
	case "/logout":
		if u := r.URL.Query().Get("post_logout_redirect_uri"); u != "" {
			http.Redirect(w, r, u, http.StatusFound)
			return
		}
		fmt.Fprintln(w, "signed out")
	default:
		http.NotFound(w, r)
	}
}

var picker = template.Must(template.New("").Parse(`<!doctype html>
<html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<title>Fake IdP</title>
<style>body{font:16px system-ui;margin:2rem auto;max-width:28rem;padding:0 1rem}a{display:block;padding:.75rem 1rem;margin:.5rem 0;border:1px solid #ccc;border-radius:.5rem;text-decoration:none;color:inherit}small{color:#666}</style>
<h1>Fake IdP</h1><p>No password is checked. Sign in as:</p>
{{range .}}<a href="{{.URL}}"><b>{{.Name}}</b><br><small>{{.Email}} · {{.Groups}}</small></a>{{end}}
</html>`))

func (p *IdP) authorize(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	redirect := q.Get("redirect_uri")
	if q.Get("client_id") != p.ClientID || redirect == "" {
		http.Error(w, "unknown client or no redirect_uri", http.StatusBadRequest)
		return
	}
	back := func(v url.Values) {
		v.Set("state", q.Get("state"))
		v.Set("iss", p.issuer(r))
		sep := "?"
		if strings.Contains(redirect, "?") {
			sep = "&"
		}
		http.Redirect(w, r, redirect+sep+v.Encode(), http.StatusFound)
	}
	if q.Get("response_type") != "code" || !strings.Contains(" "+q.Get("scope")+" ", " openid ") {
		back(url.Values{"error": {"invalid_request"}})
		return
	}
	if q.Get("code_challenge") == "" || q.Get("code_challenge_method") != "S256" {
		back(url.Values{"error": {"invalid_request"}, "error_description": {"PKCE S256 is required"}})
		return
	}
	p.mu.Lock()
	subject := q.Get("user")
	if subject == "" {
		subject = p.login
	}
	var user *User
	for i := range p.users {
		if p.users[i].Subject == subject {
			user = &p.users[i]
		}
	}
	if user == nil {
		type choice struct{ URL, Name, Email, Groups string }
		var cs []choice
		for _, u := range p.users {
			v := r.URL.Query()
			v.Set("user", u.Subject)
			cs = append(cs, choice{URL: "?" + v.Encode(), Name: firstOr(u.Name, u.Subject), Email: u.Email, Groups: strings.Join(u.Groups, ", ")})
		}
		p.mu.Unlock()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		picker.Execute(w, cs)
		return
	}
	code := b64(random(24))
	p.codes[code] = &grant{user: *user, clientID: p.ClientID, redirect: redirect, nonce: q.Get("nonce"),
		challenge: q.Get("code_challenge"), challengeMethod: "S256", expires: p.now().Add(time.Minute)}
	p.logins++
	p.mu.Unlock()
	back(url.Values{"code": {code}})
}

func (p *IdP) token(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		oauthError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	id, secret, ok := r.BasicAuth()
	if ok {
		id, _ = url.QueryUnescape(id)
		secret, _ = url.QueryUnescape(secret)
	} else {
		id, secret = r.PostForm.Get("client_id"), r.PostForm.Get("client_secret")
	}
	if id != p.ClientID || subtle.ConstantTimeCompare([]byte(secret), []byte(p.ClientSecret)) != 1 {
		oauthError(w, http.StatusUnauthorized, "invalid_client")
		return
	}
	if r.PostForm.Get("grant_type") != "authorization_code" {
		oauthError(w, http.StatusBadRequest, "unsupported_grant_type")
		return
	}
	p.mu.Lock()
	code := r.PostForm.Get("code")
	g := p.codes[code]
	delete(p.codes, code)
	p.mu.Unlock()
	sum := sha256.Sum256([]byte(r.PostForm.Get("code_verifier")))
	if g == nil || !p.now().Before(g.expires) || g.clientID != id || g.redirect != r.PostForm.Get("redirect_uri") || b64(sum[:]) != g.challenge {
		oauthError(w, http.StatusBadRequest, "invalid_grant")
		return
	}
	idToken, err := p.idToken(r, g)
	if err != nil {
		oauthError(w, http.StatusInternalServerError, "server_error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"access_token": b64(random(24)), "token_type": "Bearer", "expires_in": 300, "id_token": idToken})
}

// Claims are the claims of an ID token of a person (without the time,
// the audience and the nonce).
func (p *IdP) Claims(u User) map[string]any {
	c := map[string]any{"sub": u.Subject}
	if u.Name != "" {
		c["name"] = u.Name
	}
	if u.Email != "" {
		c["email"] = u.Email
		c["email_verified"] = !u.EmailUnverified
	}
	if u.Groups != nil {
		setPath(c, firstOr(p.GroupsClaim, "groups"), u.Groups)
	}
	for k, v := range u.Claims {
		c[k] = v
	}
	return c
}

// setPath puts a value at a dotted path of objects.
func setPath(c map[string]any, path string, v any) {
	keys := strings.Split(path, ".")
	for _, k := range keys[:len(keys)-1] {
		next, ok := c[k].(map[string]any)
		if !ok {
			next = map[string]any{}
			c[k] = next
		}
		c = next
	}
	c[keys[len(keys)-1]] = v
}

func (p *IdP) idToken(r *http.Request, g *grant) (string, error) {
	now := p.now()
	ttl := p.TokenTTL
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	c := p.Claims(g.user)
	c["iss"] = p.issuer(r)
	c["aud"] = g.clientID
	c["iat"] = now.Unix()
	c["exp"] = now.Add(ttl).Unix()
	c["auth_time"] = now.Unix()
	if g.nonce != "" {
		c["nonce"] = g.nonce
	}
	if p.Tamper != nil {
		p.Tamper(c)
	}
	return p.Sign(c)
}

// Sign makes a compact JWS (RS256) of claims with the provider's key.
func (p *IdP) Sign(claims map[string]any) (string, error) {
	head, _ := json.Marshal(map[string]string{"alg": "RS256", "typ": "JWT", "kid": "fake-1"})
	body, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	in := b64(head) + "." + b64(body)
	sum := sha256.Sum256([]byte(in))
	sig, err := rsa.SignPKCS1v15(rand.Reader, p.key, crypto.SHA256, sum[:])
	if err != nil {
		return "", err
	}
	return in + "." + b64(sig), nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func oauthError(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, map[string]string{"error": code})
}

func b64(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }

func random(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

func firstOr(s, or string) string {
	if s == "" {
		return or
	}
	return s
}
