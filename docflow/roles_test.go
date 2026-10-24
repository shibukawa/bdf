package docflow

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"testing"
)

// People with a role: topology changes the topology and sends to any
// node, audit reads the log and every job; a reprint sends a job's
// document again as the person.
func TestRolesAndReprint(t *testing.T) {
	e := newEnv(t, func(o *Options) {
		withViewers(o)
		people := map[string]*Viewer{
			"tess":   {Kind: PrincipalUser, Via: ViaOIDC, ID: "u-tess", Name: "Tess", Roles: []string{RoleTopology}},
			"audrey": {Kind: PrincipalUser, Via: ViaOIDC, ID: "u-audrey", Name: "Audrey", Roles: []string{RoleAudit}},
		}
		base := o.Viewer
		o.Viewer = func(r *http.Request) *Viewer {
			if v := people[r.Header.Get("X-Test-Viewer")]; v != nil {
				return v
			}
			return base(r)
		}
	})
	e.display("office", Node{})
	e.admin("PUT", "/api/nodes/wh-a", Node{Kind: KindDisplay, Scope: ScopeGroup, Owner: "warehouse-a"}, http.StatusOK)
	e.admin("PUT", "/api/nodes/wh-b", Node{Kind: KindDisplay, Scope: ScopeGroup, Owner: "warehouse-b"}, http.StatusOK)
	e.admin("PUT", "/api/nodes/a-in", Node{Kind: KindUpload, To: []string{"wh-a"}, Scope: ScopeGroup, Owner: "warehouse-a"}, http.StatusOK)
	doc := testDoc(t, 1)
	if res, out := e.as("alice", "POST", "/api/destinations/a-in/jobs?name=alice-a", doc); res.StatusCode != http.StatusCreated {
		t.Fatalf("alice: %d %s", res.StatusCode, out)
	}
	// the role topology changes the topology (through the API too), the
	// role audit reads the log and every job; neither has the other's
	if res, _ := e.as("tess", "PUT", "/api/nodes/q", Node{Kind: KindQueue, To: []string{"office"}}); res.StatusCode != http.StatusOK {
		t.Fatalf("tess puts a node: %d", res.StatusCode)
	}
	if res, _ := e.as("tess", "GET", "/api/audit", nil); res.StatusCode != http.StatusForbidden {
		t.Fatalf("tess reads the audit log: %d", res.StatusCode)
	}
	if res, out := e.as("tess", "GET", "/api/topology", nil); res.StatusCode != http.StatusOK || !bytes.Contains(out, []byte(`"scope":"group"`)) {
		t.Fatalf("tess's topology has no settings: %d", res.StatusCode)
	}
	if res, _ := e.as("audrey", "GET", "/api/audit", nil); res.StatusCode != http.StatusOK {
		t.Fatalf("audrey reads the audit log: %d", res.StatusCode)
	}
	if res, _ := e.as("audrey", "PUT", "/api/nodes/q2", Node{Kind: KindQueue, To: []string{"office"}}); res.StatusCode != http.StatusForbidden {
		t.Fatalf("audrey puts a node: %d", res.StatusCode)
	}
	var jobs []Job
	_, out := e.as("audrey", "GET", "/api/my/jobs", nil)
	json.Unmarshal(out, &jobs)
	if len(jobs) != 1 || jobs[0].Name != "alice-a" {
		t.Fatalf("audrey sees %+v", jobs)
	}
	if res, _ := e.as("audrey", "POST", "/api/jobs/"+jobs[0].ID+"/cancel", nil); res.StatusCode != http.StatusForbidden {
		t.Fatalf("audrey cancels: %d", res.StatusCode)
	}
	if res, out := e.as("audrey", "GET", "/api/topology", nil); res.StatusCode != http.StatusOK || bytes.Contains(out, []byte(`"scope"`)) {
		t.Fatalf("audrey's topology has settings: %d %s", res.StatusCode, out)
	}

	// a person sends to a source alone; the role topology to any node, to
	// try it out, and the role audit as everyone
	if res, out := e.as("audrey", "POST", "/api/destinations/office/jobs?name=audrey", doc); res.StatusCode != http.StatusForbidden || !bytes.Contains(out, []byte("not a source")) {
		t.Fatalf("audrey to a display: %d %s", res.StatusCode, out)
	}
	if res, out := e.as("tess", "POST", "/api/destinations/office/jobs?name=tess", doc); res.StatusCode != http.StatusCreated {
		t.Fatalf("tess to a display: %d %s", res.StatusCode, out)
	}
	// the scopes still hold: the groups' nodes are their members'
	var ds []Destination
	_, out = e.as("tess", "GET", "/api/destinations", nil)
	if json.Unmarshal(out, &ds); len(ds) != 2 || ds[0].Name != "office" || ds[1].Name != "q" {
		t.Fatalf("tess sends to %+v", ds)
	}

	// a reprint: the person's own job again, as theirs, to a node they may
	// send to; a member of the group reprints the group's, to a source, or
	// to the node it was sent to
	var again Job
	res, out := e.as("alice", "POST", "/api/jobs/"+jobs[0].ID+"/reprint", nil)
	json.Unmarshal(out, &again)
	if res.StatusCode != http.StatusCreated || again.Dest != "a-in" || again.Source != "reprint" || again.Name != "alice-a" || again.SubmittedBy.ID != "u-alice" || again.Number == jobs[0].Number {
		t.Fatalf("reprint: %d %+v", res.StatusCode, again)
	}
	if res, _ := e.as("alice", "POST", "/api/jobs/"+jobs[0].ID+"/reprint", map[string]string{"to": "wh-b"}); res.StatusCode != http.StatusForbidden {
		t.Fatalf("reprint to a node not hers: %d", res.StatusCode)
	}
	if res, _ := e.as("bob", "POST", "/api/jobs/"+jobs[0].ID+"/reprint", nil); res.StatusCode != http.StatusNotFound {
		t.Fatalf("bob reprints alice's: %d", res.StatusCode)
	}
	res, out = e.do(nil, "POST", "/api/destinations/wh-b/jobs?name=wms-b", "wms:s3cret", doc)
	var wmsB Job
	json.Unmarshal(out, &wmsB)
	if res, out := e.as("bob", "POST", "/api/jobs/"+wmsB.ID+"/reprint", nil); res.StatusCode != http.StatusCreated {
		t.Fatalf("bob reprints his group's where it went: %d %s", res.StatusCode, out)
	}
	if res, out := e.as("bob", "POST", "/api/jobs/"+wmsB.ID+"/reprint", map[string]string{"to": "office"}); res.StatusCode != http.StatusForbidden {
		t.Fatalf("bob reprints his group's to a display: %d %s", res.StatusCode, out)
	}
	if res, out := e.as("alice", "POST", "/api/jobs/"+jobs[0].ID+"/reprint", map[string]string{"to": "office"}); res.StatusCode != http.StatusForbidden {
		t.Fatalf("alice reprints to a display: %d %s", res.StatusCode, out)
	}
	var log []AuditEntry
	json.Unmarshal(e.admin("GET", "/api/audit?limit=20", nil, http.StatusOK), &log)
	if !slices.ContainsFunc(log, func(en AuditEntry) bool { return en.Action == "reprint" && en.Who == "Alice" }) {
		t.Fatalf("audit %+v", log)
	}
}

// A retired secret still opens what was signed with it, while a rotation
// goes on; nothing is signed with it.
func TestKeyRing(t *testing.T) {
	ctx := context.Background()
	old := bytes.Repeat([]byte{1}, 32)
	e := newEnv(t, func(o *Options) { o.Secret = old })
	e.display("wh-a", Node{})
	d := e.enroll("wh-a", "tablet")
	u, _ := http.NewRequest("GET", e.srv.URL, nil)
	cookie := d.c.Jar.Cookies(u.URL)[0].Value
	if _, err := e.m().Session(ctx, cookie); err != nil {
		t.Fatal(err)
	}
	// the secret rotated, the old one kept: the session goes on
	e.o.Secret = bytes.Repeat([]byte{2}, 32)
	e.o.PreviousSecrets = [][]byte{old}
	if _, err := e.m().Session(ctx, cookie); err != nil {
		t.Fatalf("with the old secret kept: %v", err)
	}
	// a new session is signed with the new secret alone
	d.login()
	fresh := d.c.Jar.Cookies(u.URL)[0].Value
	e.o.PreviousSecrets = nil
	if _, err := e.m().Session(ctx, fresh); err != nil {
		t.Fatalf("a session of the new secret: %v", err)
	}
	if _, err := e.m().Session(ctx, cookie); err == nil {
		t.Fatal("the old secret, dropped, still opens")
	}
}
