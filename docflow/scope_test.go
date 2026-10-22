package docflow

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"
)

// viewers are the people of a test, signed in by a header (the web
// application's OIDC session in cmd/docflow).
var viewers = map[string]*Viewer{
	"alice": {Kind: PrincipalUser, Via: ViaOIDC, ID: "u-alice", Name: "Alice", Email: "alice@example.com", Groups: []string{"warehouse-a"}},
	"bob":   {Kind: PrincipalUser, Via: ViaOIDC, ID: "u-bob", Name: "Bob", Email: "bob@example.com", Groups: []string{"warehouse-b"}},
	"root":  {Kind: PrincipalUser, Via: ViaOIDC, ID: "u-root", Name: "Root", Email: "root@example.com", Admin: true},
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
	e.admin("PUT", "/api/nodes/wh-b-laser", Node{Kind: KindPrinter, URL: "ipp://127.0.0.1:1/ipp/print", Scope: ScopeGroup, Owner: "warehouse-b"}, http.StatusOK)
	e.admin("PUT", "/api/nodes/wh-a", Node{Kind: KindDisplay, Scope: ScopeGroup, Owner: "warehouse-a"}, http.StatusOK)
	// a location of no scope: shared, and its jobs are the group's of its
	// fallback printer
	e.admin("PUT", "/api/nodes/wh-b", Node{Kind: KindDisplay, Fallback: &Fallback{Printer: "wh-b-laser"}}, http.StatusOK)
	e.admin("PUT", "/api/nodes/office", Node{Kind: KindDisplay, Scope: ScopeShared}, http.StatusOK)
	e.admin("PUT", "/api/nodes/alice-desk", Node{Kind: KindDisplay, Scope: ScopePersonal, Owner: "alice@example.com"}, http.StatusOK)
	e.admin("PUT", "/api/nodes/team-a", Node{Kind: KindQueue, To: []string{"wh-a"}, Scope: ScopeGroup, Owner: "warehouse-a"}, http.StatusOK)
	e.admin("PUT", "/api/nodes/bad", Node{Kind: KindDisplay, Scope: ScopeGroup}, http.StatusBadRequest)
	e.admin("PUT", "/api/nodes/bad", Node{Kind: KindDisplay, Scope: "team", Owner: "x"}, http.StatusBadRequest)
	e.admin("PUT", "/api/nodes/bad", Node{Kind: KindQueue, To: []string{"wh-a"}, Scope: ScopeShared, Owner: "x"}, http.StatusBadRequest)

	doc := testDoc(t, 1)
	submit := func(who, dest, name string, status int) *Job {
		t.Helper()
		var res *http.Response
		var out []byte
		path := "/api/destinations/" + dest + "/jobs?name=" + name
		if who == "wms" {
			res, out = e.do(nil, "POST", path+"&user=sato", "wms:s3cret", doc)
		} else {
			res, out = e.as(who, "POST", path+"&user=someone-else", doc)
		}
		if res.StatusCode != status {
			t.Fatalf("%s to %s: %d %s; want %d", who, dest, res.StatusCode, out, status)
		}
		var j Job
		json.Unmarshal(out, &j)
		return &j
	}
	mine := submit("alice", "office", "alice-office", http.StatusCreated)
	if by := mine.SubmittedBy; by == nil || *by != (Principal{Kind: PrincipalUser, ID: "u-alice", Name: "Alice", Via: ViaOIDC}) || mine.ClaimedUser != "someone-else" || mine.Group != "" {
		t.Fatalf("alice's job %+v by %+v", mine, mine.SubmittedBy)
	}
	submit("alice", "alice-desk", "alice-desk", http.StatusCreated)
	if j := submit("alice", "wh-b", "alice-wh-b", http.StatusCreated); j.Group != "warehouse-b" {
		t.Fatalf("the fallback printer's group: %+v", j)
	}
	if j := submit("alice", "team-a", "alice-team", http.StatusCreated); j.Group != "warehouse-a" {
		t.Fatalf("a group's queue: %+v", j)
	}
	submit("bob", "wh-a", "bob-wh-a", http.StatusForbidden)
	submit("bob", "alice-desk", "bob-desk", http.StatusForbidden)
	if j := submit("wms", "wh-a", "wms-wh-a", http.StatusCreated); j.Group != "warehouse-a" || j.SubmittedBy.Kind != PrincipalToken || j.ClaimedUser != "sato" {
		t.Fatalf("wms's job %+v", j)
	}
	submit("wms", "wh-b", "wms-wh-b", http.StatusCreated)
	bobs := submit("bob", "office", "bob-office", http.StatusCreated)
	// a page of another site cannot send as the person
	res, _ := e.doWith(func(r *http.Request) {
		r.Header.Set("X-Test-Viewer", "alice")
		r.Header.Set("Sec-Fetch-Site", "cross-site")
	}, "POST", "/api/destinations/office/jobs", doc)
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("cross-site: %d", res.StatusCode)
	}

	sees := func(who string) []string {
		var js []Job
		_, out := e.as(who, "GET", "/api/my/jobs", nil)
		json.Unmarshal(out, &js)
		var names []string
		for _, j := range js {
			names = append(names, j.Name)
		}
		return names
	}
	for who, want := range map[string][]string{
		"alice": {"wms-wh-a", "alice-team", "alice-wh-b", "alice-desk", "alice-office"},
		"bob":   {"bob-office", "wms-wh-b", "alice-wh-b"},
		"root":  {"bob-office", "wms-wh-b", "wms-wh-a", "alice-team", "alice-wh-b", "alice-desk", "alice-office"},
	} {
		if got := sees(who); !slices.Equal(got, want) {
			t.Errorf("%s sees %q; want %q", who, got, want)
		}
	}
	if (&Viewer{Kind: PrincipalUser, ID: "u-carol"}).CanSee(mine) || (*Viewer)(nil).CanSee(mine) {
		t.Error("an outsider sees alice's job")
	}
	// a member of the group sees a job, but only its sender cancels it
	if res, _ := e.as("bob", "GET", "/api/jobs/"+jobID(3), nil); res.StatusCode != http.StatusOK {
		t.Errorf("bob on alice's warehouse-b job: %d", res.StatusCode)
	}
	if res, _ := e.as("bob", "POST", "/api/jobs/"+jobID(3)+"/cancel", nil); res.StatusCode != http.StatusForbidden {
		t.Errorf("bob cancels alice's job: %d", res.StatusCode)
	}
	if res, _ := e.as("alice", "GET", "/api/jobs/"+bobs.ID, nil); res.StatusCode != http.StatusNotFound {
		t.Errorf("alice on bob's job: %d", res.StatusCode)
	}
	if res, _ := e.as("alice", "POST", "/api/jobs/"+jobID(3)+"/cancel", nil); res.StatusCode != http.StatusOK {
		t.Errorf("alice cancels her job: %d", res.StatusCode)
	}

	ds, err := e.m().Destinations(ctx, viewers["alice"])
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, d := range ds {
		names = append(names, d.Name)
	}
	if !slices.Equal(names, []string{"alice-desk", "office", "team-a", "wh-a", "wh-b"}) {
		t.Errorf("alice sends to %q", names)
	}

	// the API: an admin's session, as the admin token; others are refused
	if res, out := e.as("root", "GET", "/api/topology", nil); res.StatusCode != http.StatusOK {
		t.Fatalf("admin: %d %s", res.StatusCode, out)
	}
	if res, _ := e.as("alice", "GET", "/api/nodes", nil); res.StatusCode != http.StatusForbidden {
		t.Fatalf("alice on the API: %d", res.StatusCode)
	}
	if res, _ := e.as("", "GET", "/api/topology", nil); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("nobody on the API: %d", res.StatusCode)
	}
	if res, _ := e.as("root", "POST", "/api/tick", nil); res.StatusCode != http.StatusOK {
		t.Fatalf("admin's tick: %d", res.StatusCode)
	}
	res, _ = e.doWith(func(r *http.Request) {
		r.Header.Set("X-Test-Viewer", "root")
		r.Header.Set("Origin", "https://evil.example")
	}, "DELETE", "/api/nodes/office", nil)
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("admin's session from another origin: %d", res.StatusCode)
	}

	// everyone signed in sees the whole topology (its nodes and edges, not
	// its settings) and where it is drawn; an administrator the settings
	// too; nobody else at all
	res, out := e.as("alice", "GET", "/api/topology", nil)
	if res.StatusCode != http.StatusOK || strings.Contains(string(out), `"scope"`) || strings.Contains(string(out), `"owner"`) || strings.Contains(string(out), `"url"`) {
		t.Fatalf("alice's topology: %d %s", res.StatusCode, out)
	}
	var topo []TopoNode
	json.Unmarshal(out, &topo)
	if len(topo) != 6 || topo[4].Name != "wh-b" || topo[4].Fallback == nil || topo[4].Fallback.Printer != "wh-b-laser" || topo[2].To[0] != "wh-a" {
		t.Fatalf("alice's topology %+v", topo)
	}
	if _, out := e.as("root", "GET", "/api/topology", nil); !strings.Contains(string(out), `"scope":"group"`) {
		t.Fatalf("root's topology without the settings: %s", out)
	}
	if res, _ := e.as("", "GET", "/api/topology", nil); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("nobody's topology: %d", res.StatusCode)
	}
	if res, _ := e.as("alice", "GET", "/api/layout", nil); res.StatusCode != http.StatusOK {
		t.Fatalf("alice's layout: %d", res.StatusCode)
	}
	// the activity is that of the jobs each may see
	liveOf := func(who string) *Live {
		t.Helper()
		res, out := e.as(who, "GET", "/api/topology/live", nil)
		if res.StatusCode != http.StatusOK {
			t.Fatalf("%s's live: %d %s", who, res.StatusCode, out)
		}
		var l Live
		json.Unmarshal(out, &l)
		return &l
	}
	a := liveOf("alice")
	if a.Nodes["wh-a"].Open != 2 || a.Nodes["office"].Open != 1 || a.Nodes["alice-desk"].Open != 1 || a.Nodes["wh-b"].Open != 0 || a.Nodes["team-a"].Recent != 1 || a.Nodes["wh-a"].Recent != 2 {
		t.Fatalf("alice's live %+v", a.Nodes)
	}
	if !slices.ContainsFunc(a.Flows, func(f Flow) bool { return f.From == "team-a" && f.To == "wh-a" && f.Name == "alice-team" }) {
		t.Fatalf("alice's flows %+v", a.Flows)
	}
	if b := liveOf("bob"); b.Nodes["wh-b"].Open != 1 || b.Nodes["office"].Open != 1 || b.Nodes["wh-a"].Open != 0 || b.Nodes["alice-desk"].Open != 0 {
		t.Fatalf("bob's live %+v", b.Nodes)
	}
	if r := liveOf("root"); r.Nodes["office"].Open != 2 || r.Nodes["wh-a"].Open != 2 || r.Nodes["wh-b"].Open != 1 {
		t.Fatalf("root's live %+v", r.Nodes)
	}
	// a client polls with the revision: 304 while nothing changed
	res, _ = e.doWith(func(r *http.Request) {
		r.Header.Set("X-Test-Viewer", "alice")
		r.Header.Set("If-None-Match", `"`+a.Rev+`"`)
	}, "GET", "/api/topology/live", nil)
	if res.StatusCode != http.StatusNotModified {
		t.Fatalf("live with its revision: %d", res.StatusCode)
	}
	submit("alice", "office", "alice-again", http.StatusCreated)
	if again := liveOf("alice"); again.Rev == a.Rev || again.Nodes["office"].Open != 2 {
		t.Fatalf("live after a job: %+v", again)
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
