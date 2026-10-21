package manager

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"testing"
	"time"
)

// viewers are the people of a test, signed in by a header (the web
// application's session in cmd/printmanager).
var viewers = map[string]*Viewer{
	"alice": {ID: "u-alice", Name: "Alice", Email: "alice@example.com", Groups: []string{"warehouse-a"}},
	"bob":   {ID: "u-bob", Name: "Bob", Email: "bob@example.com", Groups: []string{"warehouse-b"}},
	"root":  {ID: "u-root", Name: "Root", Email: "root@example.com", Admin: true},
}

func withViewers(o *Options) {
	o.Viewer = func(r *http.Request) *Viewer { return viewers[r.Header.Get("X-Test-Viewer")] }
	o.SubmitTokens = map[string]string{"wms": "s3cret"}
}

// as sends a request of a signed-in person, from the manager's own page.
func (e *env) as(who, method, path string, body any) (*http.Response, []byte) {
	e.t.Helper()
	return e.doWith(func(r *http.Request) {
		r.Header.Set("X-Test-Viewer", who)
		r.Header.Set("Sec-Fetch-Site", "same-origin")
	}, method, path, body)
}

func TestScope(t *testing.T) {
	e := newEnv(t, withViewers)
	ctx := context.Background()
	e.admin("PUT", "/api/printers/wh-b-laser", Printer{URL: "ipp://127.0.0.1:1/ipp/print", Kind: KindGroup, Owner: "warehouse-b"}, http.StatusOK)
	e.admin("PUT", "/api/locations/wh-a", Location{Kind: KindGroup, Owner: "warehouse-a"}, http.StatusOK)
	// a location of no group: its slips are the group's of its fallback printer
	e.admin("PUT", "/api/locations/wh-b", Location{Fallback: &Fallback{Printer: "wh-b-laser"}}, http.StatusOK)
	e.admin("PUT", "/api/locations/office", Location{Kind: KindShared}, http.StatusOK)
	e.admin("PUT", "/api/locations/alice-desk", Location{Kind: KindPersonal, Owner: "alice@example.com"}, http.StatusOK)
	e.admin("PUT", "/api/locations/bad", Location{Kind: KindGroup}, http.StatusBadRequest)
	e.admin("PUT", "/api/locations/bad", Location{Kind: "team", Owner: "x"}, http.StatusBadRequest)

	doc := testDoc(t, 1)
	submit := func(who, loc, name string, status int) *Slip {
		t.Helper()
		var res *http.Response
		var out []byte
		if who == "wms" {
			res, out = e.do(nil, "POST", "/api/locations/"+loc+"/slips?name="+name+"&user=sato", "wms:s3cret", doc)
		} else {
			res, out = e.as(who, "POST", "/api/locations/"+loc+"/slips?name="+name+"&user=someone-else", doc)
		}
		if res.StatusCode != status {
			t.Fatalf("%s to %s: %d %s; want %d", who, loc, res.StatusCode, out, status)
		}
		var s Slip
		json.Unmarshal(out, &s)
		return &s
	}
	mine := submit("alice", "office", "alice-office", http.StatusCreated)
	if by := mine.SubmittedBy; by == nil || *by != (Principal{Kind: PrincipalUser, ID: "u-alice", Name: "Alice", Via: ViaOIDC}) || mine.ClaimedUser != "someone-else" || mine.Group != "" {
		t.Fatalf("alice's slip %+v by %+v", mine, mine.SubmittedBy)
	}
	submit("alice", "alice-desk", "alice-desk", http.StatusCreated)
	// wh-b is shared (no kind), and its slips are warehouse-b's
	if s := submit("alice", "wh-b", "alice-wh-b", http.StatusCreated); s.Group != "warehouse-b" {
		t.Fatalf("the fallback printer's group: %+v", s)
	}
	submit("bob", "wh-a", "bob-wh-a", http.StatusForbidden)
	submit("bob", "alice-desk", "bob-desk", http.StatusForbidden)
	if s := submit("wms", "wh-a", "wms-wh-a", http.StatusCreated); s.Group != "warehouse-a" || s.SubmittedBy.Kind != PrincipalToken || s.ClaimedUser != "sato" {
		t.Fatalf("wms's slip %+v", s)
	}
	submit("wms", "wh-b", "wms-wh-b", http.StatusCreated)
	submit("bob", "office", "bob-office", http.StatusCreated)
	// a page of another site cannot send as the person
	res, _ := e.doWith(func(r *http.Request) {
		r.Header.Set("X-Test-Viewer", "alice")
		r.Header.Set("Sec-Fetch-Site", "cross-site")
	}, "POST", "/api/locations/office/slips", doc)
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("cross-site: %d", res.StatusCode)
	}

	names := func(who string) []string {
		ss, err := e.m().SlipsOf(ctx, viewers[who])
		if err != nil {
			t.Fatal(err)
		}
		var out []string
		for _, s := range ss {
			out = append(out, s.Name)
		}
		return out
	}
	for who, want := range map[string][]string{
		"alice": {"wms-wh-a", "alice-wh-b", "alice-desk", "alice-office"},
		"bob":   {"bob-office", "wms-wh-b", "alice-wh-b"},
		"root":  {"bob-office", "wms-wh-b", "wms-wh-a", "alice-wh-b", "alice-desk", "alice-office"},
	} {
		if got := names(who); !slices.Equal(got, want) {
			t.Errorf("%s sees %q; want %q", who, got, want)
		}
	}
	if (&Viewer{ID: "u-carol"}).CanSee(mine) || (*Viewer)(nil).CanSee(mine) {
		t.Error("an outsider sees alice's slip")
	}

	p, err := e.m().PlacesOf(ctx, viewers["alice"])
	if err != nil {
		t.Fatal(err)
	}
	var places []string
	for _, l := range p.Locations {
		places = append(places, l.Name)
	}
	if !slices.Equal(places, []string{"alice-desk", "office", "wh-a", "wh-b"}) {
		t.Errorf("alice prints to %q", places)
	}
	if len(p.Printers) != 0 {
		t.Errorf("alice prints on %+v", p.Printers)
	}

	// the API: an admin's session, as the admin token; others are refused
	if res, out := e.as("root", "GET", "/api/locations", nil); res.StatusCode != http.StatusOK {
		t.Fatalf("admin: %d %s", res.StatusCode, out)
	}
	if res, _ := e.as("alice", "GET", "/api/locations", nil); res.StatusCode != http.StatusForbidden {
		t.Fatalf("alice on the API: %d", res.StatusCode)
	}
	if res, _ := e.as("", "GET", "/api/locations", nil); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("nobody on the API: %d", res.StatusCode)
	}
	if res, _ := e.as("root", "POST", "/api/tick", nil); res.StatusCode != http.StatusOK {
		t.Fatalf("admin's tick: %d", res.StatusCode)
	}
	res, _ = e.doWith(func(r *http.Request) {
		r.Header.Set("X-Test-Viewer", "root")
		r.Header.Set("Origin", "https://evil.example")
	}, "DELETE", "/api/locations/office", nil)
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("admin's session from another origin: %d", res.StatusCode)
	}
}

func TestStates(t *testing.T) {
	e := newEnv(t, nil)
	ctx := context.Background()
	now := e.clock.Now()
	if err := e.st.PutState(ctx, "auth-oidc\x00k1", []byte("v1"), now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := e.st.PutState(ctx, "auth-oidc\x00k1", []byte("again"), now.Add(time.Minute)); !errors.Is(err, ErrConflict) {
		t.Fatalf("a key twice: %v", err)
	}
	e.st.PutState(ctx, "auth-oidc\x00old", []byte("v2"), now.Add(-time.Second))
	e.st.PutState(ctx, "auth-oidc\x00later", []byte("v3"), now.Add(-time.Second))
	if v, err := e.st.TakeState(ctx, "auth-oidc\x00k1", now); err != nil || string(v) != "v1" {
		t.Fatalf("take: %q %v", v, err)
	}
	if _, err := e.st.TakeState(ctx, "auth-oidc\x00k1", now); !errors.Is(err, ErrNotFound) {
		t.Fatalf("taken twice: %v", err)
	}
	if _, err := e.st.TakeState(ctx, "auth-oidc\x00old", now); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expired: %v", err)
	}
	if n, err := e.st.PruneStates(ctx, now, 10); err != nil || n != 1 {
		t.Fatalf("pruned %d %v", n, err)
	}
}
