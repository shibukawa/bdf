package fakeidp

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestCodeFlow(t *testing.T) {
	idp := New(User{Subject: "alice", Email: "alice@example.com", Groups: []string{"a"}})
	idp.GroupsClaim = "realm_access.roles"
	iss := Start(t, idp)
	var disc map[string]any
	res, err := http.Get(iss + "/.well-known/openid-configuration")
	if err != nil {
		t.Fatal(err)
	}
	json.NewDecoder(res.Body).Decode(&disc)
	res.Body.Close()
	if disc["issuer"] != iss || disc["token_endpoint"] != iss+"/token" {
		t.Fatalf("discovery %v", disc)
	}

	noRedirect := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	verifier := "a-verifier-of-forty-three-characters-or-more-x"
	sum := sha256.Sum256([]byte(verifier))
	authorize := func(challenge string) url.Values {
		t.Helper()
		q := url.Values{"client_id": {ClientID}, "redirect_uri": {"http://rp/cb"}, "response_type": {"code"}, "scope": {"openid email"},
			"state": {"st"}, "nonce": {"n"}, "code_challenge": {challenge}, "code_challenge_method": {"S256"}}
		res, err := noRedirect.Get(iss + "/authorize?" + q.Encode())
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		u, _ := url.Parse(res.Header.Get("Location"))
		return u.Query()
	}
	idp.LoginAs("alice")
	back := authorize(b64(sum[:]))
	if back.Get("state") != "st" || back.Get("code") == "" {
		t.Fatalf("back %v", back)
	}
	token := func(code, verifier string) (int, map[string]any) {
		form := url.Values{"grant_type": {"authorization_code"}, "code": {code}, "redirect_uri": {"http://rp/cb"}, "code_verifier": {verifier}}
		req, _ := http.NewRequest("POST", iss+"/token", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth(ClientID, ClientSecret)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		var out map[string]any
		json.NewDecoder(res.Body).Decode(&out)
		return res.StatusCode, out
	}
	if code, _ := token(back.Get("code"), "wrong-verifier"); code != http.StatusBadRequest {
		t.Fatalf("a wrong verifier: %d", code)
	}
	back = authorize(b64(sum[:]))
	code, out := token(back.Get("code"), verifier)
	if code != http.StatusOK {
		t.Fatalf("token: %d %v", code, out)
	}
	if code, _ := token(back.Get("code"), verifier); code != http.StatusBadRequest {
		t.Fatalf("a code twice: %d", code)
	}
	parts := strings.Split(out["id_token"].(string), ".")
	var claims struct {
		Iss, Sub, Aud, Nonce string
		Realm                struct{ Roles []string } `json:"realm_access"`
	}
	b, _ := b64dec(parts[1])
	json.Unmarshal(b, &claims)
	if claims.Iss != iss || claims.Sub != "alice" || claims.Aud != ClientID || claims.Nonce != "n" || len(claims.Realm.Roles) != 1 {
		t.Fatalf("claims %s", b)
	}
	// without PKCE, the IdP refuses
	if back := authorize(""); back.Get("error") != "invalid_request" {
		t.Fatalf("no PKCE: %v", back)
	}
}

func b64dec(s string) ([]byte, error) { return base64.RawURLEncoding.DecodeString(s) }
