package docflow

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// A downloader lists what an archive before it wrote, under the owner the
// metadata names, and serves each entry to its owner, to the apps it
// lists and to the administrators: anything else is not found.
func TestDownload(t *testing.T) {
	e := newEnv(t, withViewers)
	ctx := context.Background()
	store := "bbolt:" + filepath.ToSlash(filepath.Join(t.TempDir(), "vault.db"))
	// the apps: a backend with the scope, one without, one not listed
	var backend, nope, other struct {
		Secret string `json:"secret"`
	}
	json.Unmarshal(e.admin("PUT", "/api/apps/backend", App{Scopes: []string{ScopeDownload}}, http.StatusOK), &backend)
	json.Unmarshal(e.admin("PUT", "/api/apps/nope", App{Scopes: []string{ScopePrint}}, http.StatusOK), &nope)
	json.Unmarshal(e.admin("PUT", "/api/apps/other", App{Scopes: []string{ScopeDownload}}, http.StatusOK), &other)
	// the downloader's settings are checked
	e.node("out", Node{Kind: KindDownload, By: "bad key!"}, http.StatusBadRequest)
	e.node("out", Node{Kind: KindDownload, By: "customer", Subject: "phone"}, http.StatusBadRequest)
	e.node("out", Node{Kind: KindDownload, By: "customer", Entries: 9999}, http.StatusBadRequest)
	e.node("out", Node{Kind: KindDownload, By: "customer", Apps: []string{"nobody"}}, http.StatusBadRequest)
	e.node("out", Node{Kind: KindDownload, By: "customer", Fields: []string{"invoice"}, Apps: []string{"backend", "nope"}}, http.StatusOK)
	var got Node
	json.Unmarshal(e.admin("GET", "/api/nodes/out", nil, http.StatusOK), &got)
	if got.By != "customer" || got.Subject != SubjectEmail || len(got.Fields) != 1 || len(got.Apps) != 2 {
		t.Fatalf("export %+v", got)
	}
	// a downloader takes a stocked file, which an archive sends along its to;
	// any other node, or another of its edges, is refused, with the edge
	e.node("vault", Node{Kind: KindArchive, Store: store, To: []string{"out"}}, http.StatusOK)
	e.node("in", Node{Kind: KindQueue, To: []string{"vault"}}, http.StatusOK)
	e.node("seen", Node{Kind: KindSink}, http.StatusOK)
	for name, n := range map[string]Node{
		"direct": {Kind: KindQueue, To: []string{"out"}},
		"api":    {Kind: KindUpload, To: []string{"in", "out"}},
		"conv":   {Kind: KindConvert, To: []string{"out"}},
		"pick":   {Kind: KindFilter, Otherwise: []string{"out"}},
		"hold":   {Kind: KindInbox, To: []string{"seen"}, OnError: []string{"out"}},
	} {
		if out := e.node(name, n, http.StatusBadRequest); !strings.Contains(string(out), "out takes a stocked file") || !strings.Contains(string(out), "not "+name) {
			t.Fatalf("%s to the export: %s", name, out)
		}
	}
	e.node("fan", Node{Kind: KindQueue, To: []string{"vault"}}, http.StatusOK)
	e.node("api", Node{Kind: KindUpload, To: []string{"fan"}}, http.StatusOK)
	// a queue after the archive gets what it wrote, not a stocked file;
	// the archive still sends that on to the others
	e.node("after", Node{Kind: KindQueue, To: []string{"seen"}}, http.StatusOK)
	e.node("after", Node{Kind: KindQueue, To: []string{"seen", "out"}}, http.StatusBadRequest)
	e.node("vault", Node{Kind: KindArchive, Store: store, To: []string{"out", "after"}}, http.StatusOK)
	var typed []TopoNode
	json.Unmarshal(e.admin("GET", "/api/topology", nil, http.StatusOK), &typed)
	for _, n := range typed {
		switch {
		case n.Name == "out" && (n.Takes != PayloadStocked || len(n.Broken) > 0),
			n.Name == "vault" && (len(n.Accepts) != 3 || len(n.Broken) > 0 || strings.Join(n.Emits, ",") != "bdf"),
			n.Name == "after" && strings.Join(n.Emits, ",") != "bdf":
			t.Fatalf("typing of %s: %+v", n.Name, n)
		}
	}
	e.node("vault", Node{Kind: KindArchive, Store: store, To: []string{"out"}}, http.StatusOK)
	e.admin("DELETE", "/api/nodes/after", nil, http.StatusNoContent)
	e.admin("DELETE", "/api/nodes/seen", nil, http.StatusNoContent)

	doc := testDoc(t, 2)
	send := func(name, meta string) {
		t.Helper()
		if res, out := e.do(nil, "POST", "/api/destinations/in/jobs?name="+name+meta, "wms:s3cret", doc); res.StatusCode != http.StatusCreated {
			t.Fatalf("%s: %d %s", name, res.StatusCode, out)
		}
	}
	send("inv-1", "&meta.customer=alice@example.com&meta.invoice=INV-1")
	send("inv-2", "&meta.customer=bob@example.com&meta.invoice=INV-2")
	send("inv-3", "&meta.customer=Alice@Example.com&meta.invoice=INV-3&meta.secret=hidden")
	send("memo", "")
	if rep := e.tick(); rep.Delivered != 4 {
		t.Fatalf("tick %+v", rep)
	}
	j := e.job(3)
	if legs(j) != "in:routed vault:done out:listed(Alice@Example.com)" || j.State != JobCompleted || j.Live != 0 {
		t.Fatalf("job %+v: %s", j, legs(j))
	}
	if legs(e.job(4)) != "in:routed vault:done out:listed" {
		t.Fatalf("a job with no owner: %s", legs(e.job(4)))
	}

	// a person reads their own: by the email the metadata names, without case
	entries := func(res *http.Response, out []byte) []DownloadEntry {
		t.Helper()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("list: %d %s", res.StatusCode, out)
		}
		var es []DownloadEntry
		json.Unmarshal(out, &es)
		return es
	}
	es := entries(e.as("alice", "GET", "/api/download/out", nil))
	if len(es) != 2 || es[0].Job != 3 || es[1].Job != 1 || es[0].Owner != "Alice@Example.com" || es[0].Fields["invoice"] != "INV-3" || es[0].Fields["secret"] != "" ||
		es[0].Archive != "vault" || es[0].Object != es[0].Key+".bdf" || es[0].Pages != 2 || es[0].SubmittedBy == nil || es[0].SubmittedBy.ID != "wms" {
		t.Fatalf("alice's entries %+v", es)
	}
	if es := entries(e.as("alice", "GET", "/api/download/out?limit=1", nil)); len(es) != 1 || es[0].Job != 3 {
		t.Fatalf("limit %+v", es)
	}
	if es := entries(e.as("alice", "GET", "/api/download/out?owner=bob@example.com", nil)); len(es) != 0 {
		t.Fatalf("another's owner %+v", es)
	}
	if es := entries(e.as("bob", "GET", "/api/download/out", nil)); len(es) != 1 || es[0].Job != 2 {
		t.Fatalf("bob's entries %+v", es)
	}
	// an administrator reads every one, and one owner's
	if es := entries(e.as("root", "GET", "/api/download/out", nil)); len(es) != 4 || es[0].Job != 4 || es[3].Owner != "alice@example.com" {
		t.Fatalf("root's entries %+v", es)
	}
	if es := entries(e.as("root", "GET", "/api/download/out?owner=alice@example.com", nil)); len(es) != 2 || es[0].Job != 3 || es[1].Job != 1 {
		t.Fatalf("root, one owner's (an email without case): %+v", es)
	}
	// an app the downloader lists reads every one with the scope; without
	// it, or not listed, it is forbidden; a submit token is nobody here
	token := func(id, secret string) string {
		t.Helper()
		res, out := e.doWith(func(r *http.Request) {
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			r.SetBasicAuth(id, secret)
		}, "POST", OAuthTokenPath, []byte("grant_type=client_credentials"))
		var body map[string]any
		json.Unmarshal(out, &body)
		if res.StatusCode != http.StatusOK {
			t.Fatalf("token of %s: %d %+v", id, res.StatusCode, body)
		}
		return body["access_token"].(string)
	}
	bt := token("backend", backend.Secret)
	if es := entries(e.do(nil, "GET", "/api/download/out", bt, nil)); len(es) != 4 {
		t.Fatalf("the backend's entries %+v", es)
	}
	if es := entries(e.do(nil, "GET", "/api/download/out?owner=ALICE@example.com", bt, nil)); len(es) != 2 || es[0].Job != 3 {
		t.Fatalf("the backend, one owner's: %+v", es)
	}
	for _, c := range []struct {
		who  string
		tok  string
		want int
	}{{"nope", token("nope", nope.Secret), http.StatusForbidden}, {"other", token("other", other.Secret), http.StatusForbidden}, {"wms", "wms:s3cret", http.StatusForbidden}, {"nobody", "", http.StatusUnauthorized}} {
		if res, out := e.do(nil, "GET", "/api/download/out", c.tok, nil); res.StatusCode != c.want {
			t.Fatalf("%s: %d %s; want %d", c.who, res.StatusCode, out, c.want)
		}
	}
	// a person may not read the downloader of a group they are not in
	e.node("out", Node{Kind: KindDownload, By: "customer", Fields: []string{"invoice"}, Apps: []string{"backend", "nope"}, Scope: ScopeGroup, Owner: "warehouse-b"}, http.StatusOK)
	if res, _ := e.as("alice", "GET", "/api/download/out", nil); res.StatusCode != http.StatusForbidden {
		t.Fatalf("alice at warehouse-b's export: %d", res.StatusCode)
	}
	if es := entries(e.as("bob", "GET", "/api/download/out", nil)); len(es) != 1 {
		t.Fatalf("bob at his group's export: %+v", es)
	}
	e.node("out", Node{Kind: KindDownload, By: "customer", Fields: []string{"invoice"}, Apps: []string{"backend", "nope"}}, http.StatusOK)

	// one entry: the owner's, or not found; its record and document come
	// from the archive
	for _, c := range []struct {
		who, path string
		want      int
	}{{"alice", "/api/download/out/1", 200}, {"alice", "/api/download/out/3", 200}, {"alice", "/api/download/out/2", 404}, {"alice", "/api/download/out/4", 404}, {"root", "/api/download/out/4", 200}, {"alice", "/api/download/out/99", 404},
		{"alice", "/api/download/nowhere/1", 404}, {"alice", "/api/download/in/1", 404}, {"alice", "/api/download/out/x", 404}, {"bob", "/api/download/out/1/document", 404}, {"alice", "/api/download/out/1?owner=bob@example.com", 404}} {
		if res, out := e.as(c.who, "GET", c.path, nil); res.StatusCode != c.want {
			t.Fatalf("%s %s: %d %s; want %d", c.who, c.path, res.StatusCode, out, c.want)
		}
	}
	res, out := e.as("alice", "GET", "/api/download/out/1/document", nil)
	if res.StatusCode != http.StatusOK || !bytes.Equal(out, doc) || res.Header.Get("Content-Disposition") != `attachment; filename="inv-1.bdf"` || res.Header.Get("Content-Type") != "application/x-bdf" {
		t.Fatalf("document: %d %s %s", res.StatusCode, res.Header.Get("Content-Disposition"), res.Header.Get("Content-Type"))
	}
	res, out = e.do(nil, "GET", "/api/download/out/3/record", bt, nil)
	var rec ArchiveRecord
	if res.StatusCode != http.StatusOK || json.Unmarshal(out, &rec) != nil || rec.Job != 3 || rec.Meta["customer"] != "Alice@Example.com" || rec.Meta["secret"] != "hidden" || rec.Archive != "vault" {
		t.Fatalf("record: %d %+v", res.StatusCode, rec)
	}
	// the download is in the audit log
	var log []AuditEntry
	json.Unmarshal(e.admin("GET", "/api/audit?limit=5", nil, http.StatusOK), &log)
	if len(log) == 0 || log[0].Action != "download-document" || log[0].Who != "Alice" || log[0].Subject != "out" || !strings.HasPrefix(log[0].Detail, "job 1 ") {
		t.Fatalf("audit %+v", log)
	}
	// the topology and the activity count what is listed
	var topo []TopoNode
	json.Unmarshal(e.admin("GET", "/api/topology", nil, http.StatusOK), &topo)
	for _, n := range topo {
		if n.Name == "out" && n.Listed != 4 {
			t.Fatalf("topology %+v", n)
		}
		if n.Name == "vault" && (len(n.To) != 1 || len(n.Emits) != 1 || n.Emits[0] != PayloadBDF) {
			t.Fatalf("the archive's edge %+v", n)
		}
	}
	var l Live
	_, out = e.as("alice", "GET", "/api/topology/live", nil)
	json.Unmarshal(out, &l)
	if l.Nodes["out"].Listed != 0 { // alice did not send them: she sees her entries at the export, not the jobs
		t.Fatalf("live %+v", l.Nodes["out"])
	}
	_, out = e.as("root", "GET", "/api/topology/live", nil)
	json.Unmarshal(out, &l)
	if l.Nodes["out"].Listed != 4 {
		t.Fatalf("live %+v", l.Nodes["out"])
	}

	// the archive's step done twice (a retry whose first outcome was
	// lost) lists the job once
	m := e.m()
	var d *Delivery
	for _, leg := range e.job(1).Legs {
		if leg.Node == "vault" {
			d, _ = m.Delivery(ctx, leg.Ref)
		}
	}
	vault, _ := m.Node(ctx, "vault")
	if err := m.continueFrom(ctx, d, vault, e.job(1).received(), false, &stockRef{Node: "vault", Key: d.Result}); err != nil {
		t.Fatal(err)
	}
	if es := entries(e.as("alice", "GET", "/api/download/out", nil)); len(es) != 2 {
		t.Fatalf("listed twice %+v", es)
	}
	// a job sent straight to the downloader is refused: it has no archive
	// before it; nor is the downloader a destination to choose
	if res, out := e.do(nil, "POST", "/api/destinations/out/jobs?name=straight", "wms:s3cret", doc); res.StatusCode != http.StatusBadRequest || !strings.Contains(string(out), "is a downloader") {
		t.Fatalf("straight: %d %s", res.StatusCode, out)
	}
	var ds []Destination
	_, out = e.as("alice", "GET", "/api/destinations", nil)
	json.Unmarshal(out, &ds)
	for _, d := range ds {
		if d.Kind == KindDownload {
			t.Fatalf("an export among the destinations: %+v", ds)
		}
	}
	// a job routed to it by hand, or along a way the graph no longer
	// has, gets an error leg there
	r := &routing{job: e.job(1), now: m.now(), payload: e.job(1).received()}
	m.routeDownload(ctx, r, "out", "hold")
	r.stock = &stockRef{Node: "vault", Key: "x"}
	m.route(ctx, r, "out", "in", 0)
	if len(r.legs) != 2 || r.legs[0].State != LegError || !strings.Contains(r.legs[0].Error, "not from an inbox or an archive") ||
		r.legs[1].State != LegError || !strings.Contains(r.legs[1].Error, "takes stocked; this is bdf from in") {
		t.Fatalf("routed by hand: %+v", r.legs)
	}

	// a downloader with no By shows a person the entries of the jobs they
	// could see; a group subject compares the person's groups
	store2 := "bbolt:" + filepath.ToSlash(filepath.Join(t.TempDir(), "vault2.db"))
	e.node("mine", Node{Kind: KindDownload}, http.StatusOK)
	e.node("dept", Node{Kind: KindDownload, By: "dept", Subject: SubjectGroup}, http.StatusOK)
	e.node("vault2", Node{Kind: KindArchive, Store: store2, To: []string{"mine", "dept"}}, http.StatusOK)
	e.node("in2", Node{Kind: KindQueue, To: []string{"vault2"}}, http.StatusOK)
	if res, out := e.as("alice", "POST", "/api/destinations/in2/jobs?name=alices&meta.dept=warehouse-b", doc); res.StatusCode != http.StatusCreated {
		t.Fatalf("alice sends: %d %s", res.StatusCode, out)
	}
	send2 := func(name, meta string) {
		t.Helper()
		if res, out := e.do(nil, "POST", "/api/destinations/in2/jobs?name="+name+meta, "wms:s3cret", doc); res.StatusCode != http.StatusCreated {
			t.Fatalf("%s: %d %s", name, res.StatusCode, out)
		}
	}
	send2("wms-a", "&meta.dept=warehouse-a")
	send2("wms-b", "&meta.dept=warehouse-b")
	if rep := e.tick(); rep.Delivered != 3 {
		t.Fatalf("tick %+v", rep)
	}
	if es := entries(e.as("alice", "GET", "/api/download/mine", nil)); len(es) != 1 || es[0].Name != "alices" {
		t.Fatalf("alice's own at mine: %+v", es)
	}
	if es := entries(e.as("bob", "GET", "/api/download/mine", nil)); len(es) != 0 {
		t.Fatalf("bob at mine: %+v", es)
	}
	if es := entries(e.as("root", "GET", "/api/download/mine", nil)); len(es) != 3 {
		t.Fatalf("root at mine: %+v", es)
	}
	if es := entries(e.as("alice", "GET", "/api/download/dept", nil)); len(es) != 1 || es[0].Name != "wms-a" {
		t.Fatalf("alice's group at dept: %+v", es)
	}
	if es := entries(e.as("bob", "GET", "/api/download/dept", nil)); len(es) != 2 || es[0].Name != "wms-b" || es[1].Name != "alices" {
		t.Fatalf("bob's group at dept: %+v", es)
	}

	// after the jobs are gone, the entries stay and the documents are
	// read from the archive
	e.clock.add(15 * 24 * time.Hour)
	if rep := e.tick(); rep.Jobs == 0 {
		t.Fatalf("tick %+v", rep)
	}
	if res, _ := e.do(nil, "GET", "/api/jobs/1", adminToken, nil); res.StatusCode != http.StatusNotFound {
		t.Fatalf("the job stayed: %d", res.StatusCode)
	}
	if es := entries(e.as("alice", "GET", "/api/download/out", nil)); len(es) != 2 {
		t.Fatalf("after the jobs went: %+v", es)
	}
	if res, out := e.as("alice", "GET", "/api/download/out/1/document", nil); res.StatusCode != http.StatusOK || !bytes.Equal(out, doc) {
		t.Fatalf("document after the job went: %d", res.StatusCode)
	}
	// deleting the downloader forgets its entries; the archive keeps the documents
	e.admin("DELETE", "/api/nodes/out", nil, http.StatusConflict) // vault sends to it
	e.admin("DELETE", "/api/nodes/vault", nil, http.StatusConflict)
	e.admin("DELETE", "/api/nodes/in", nil, http.StatusNoContent)
	e.admin("DELETE", "/api/nodes/api", nil, http.StatusNoContent)
	e.admin("DELETE", "/api/nodes/fan", nil, http.StatusNoContent)
	e.admin("DELETE", "/api/nodes/vault", nil, http.StatusNoContent)
	e.admin("DELETE", "/api/nodes/out", nil, http.StatusNoContent)
	if _, err := m.Entry(ctx, "out", 1); !errors.Is(err, ErrNotFound) {
		t.Fatalf("the entry stayed: %v", err)
	}
	if ids := m.downloadJobs(ctx, downloadListID("out", "alice@example.com")); len(ids) != 0 {
		t.Fatalf("the list stayed: %v", ids)
	}
	if _, err := m.Entry(ctx, "mine", 6); err != nil {
		t.Fatalf("another export's entry went: %v", err)
	}
}

// A frontend registered as a public client gets a person's token with
// the authorization code grant and PKCE, and reads a downloader with it.
func TestOAuthCode(t *testing.T) {
	e := newEnv(t, withViewers)
	ctx := context.Background()
	m := e.m()
	// the people the tokens are of: their records, as the web
	// application keeps them at each login
	for id, v := range viewers {
		if err := m.PutPerson(ctx, v.ID, Person{Issuer: "https://idp.example", Subject: id, Name: v.Name, Email: v.Email, Groups: v.Groups, Login: m.now()}); err != nil {
			t.Fatal(err)
		}
	}
	store := "bbolt:" + filepath.ToSlash(filepath.Join(t.TempDir(), "vault.db"))
	e.node("out", Node{Kind: KindDownload, By: "customer", Subject: SubjectSub}, http.StatusOK)
	e.node("vault", Node{Kind: KindArchive, Store: store, To: []string{"out"}}, http.StatusOK)
	e.node("in", Node{Kind: KindQueue, To: []string{"vault"}}, http.StatusOK)
	doc := testDoc(t, 1)
	for _, who := range []string{"alice", "bob"} {
		if res, out := e.do(nil, "POST", "/api/destinations/in/jobs?name="+who+"&meta.customer="+who, "wms:s3cret", doc); res.StatusCode != http.StatusCreated {
			t.Fatalf("%s: %d %s", who, res.StatusCode, out)
		}
	}
	e.tick()

	// a public client: redirect URIs are required and checked, and it
	// has no secret
	e.admin("PUT", "/api/apps/portal", App{Public: true, Scopes: []string{ScopeDownload}}, http.StatusBadRequest)
	e.admin("PUT", "/api/apps/portal", App{Public: true, Scopes: []string{ScopeDownload}, RedirectURIs: []string{"http://portal.example/cb"}}, http.StatusBadRequest)
	e.admin("PUT", "/api/apps/portal", App{Public: true, Scopes: []string{ScopeDownload}, RedirectURIs: []string{"https://portal.example/cb#x"}}, http.StatusBadRequest)
	e.admin("PUT", "/api/apps/portal", App{Scopes: []string{ScopeDownload}, RedirectURIs: []string{"https://portal.example/cb"}}, http.StatusBadRequest)
	var made struct {
		App
		Secret string `json:"secret"`
	}
	json.Unmarshal(e.admin("PUT", "/api/apps/portal", App{Title: "Portal", Public: true, Scopes: []string{ScopeDownload, ScopePrint}, RedirectURIs: []string{"https://portal.example/cb", "http://localhost:5173/cb"}}, http.StatusOK), &made)
	if made.Secret != "" || !made.Public || len(made.RedirectURIs) != 2 {
		t.Fatalf("public client %+v", made)
	}
	e.admin("POST", "/api/apps/portal/secret", nil, http.StatusBadRequest)
	res, out := e.doWith(func(r *http.Request) {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.SetBasicAuth("portal", "")
	}, "POST", OAuthTokenPath, []byte("grant_type=client_credentials"))
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("client credentials of a public client: %d %s", res.StatusCode, out)
	}

	// the person is sent back with a code: a registered client and
	// redirect URI, PKCE with S256, scopes of the client's
	noRedirect := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	verifier := "a-verifier-of-at-least-forty-three-characters-0123456789"
	sum := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(sum[:])
	authorize := func(who string, params url.Values) (int, *url.URL) {
		t.Helper()
		q := url.Values{"response_type": {"code"}, "client_id": {"portal"}, "redirect_uri": {"https://portal.example/cb"}, "state": {"xyz"}, "code_challenge": {challenge}, "code_challenge_method": {"S256"}}
		for k, v := range params {
			q[k] = v
		}
		res, out := e.send(noRedirect, func(r *http.Request) {
			if who != "" {
				r.Header.Set("X-Test-Viewer", who)
			}
		}, "GET", OAuthAuthorizePath+"?"+q.Encode(), nil)
		loc, _ := url.Parse(res.Header.Get("Location"))
		if res.StatusCode != http.StatusFound {
			t.Logf("authorize: %d %s", res.StatusCode, out)
		}
		return res.StatusCode, loc
	}
	if code, _ := authorize("alice", url.Values{"client_id": {"nobody"}}); code != http.StatusBadRequest {
		t.Fatalf("an unknown client: %d", code)
	}
	if code, _ := authorize("alice", url.Values{"redirect_uri": {"https://evil.example/cb"}}); code != http.StatusBadRequest {
		t.Fatalf("another redirect URI: %d", code)
	}
	if code, _ := authorize("", nil); code != http.StatusUnauthorized {
		t.Fatalf("nobody signed in: %d", code)
	}
	if code, loc := authorize("alice", url.Values{"code_challenge_method": {"plain"}}); code != http.StatusFound || loc.Query().Get("error") != "invalid_request" || loc.Query().Get("state") != "xyz" {
		t.Fatalf("plain PKCE: %d %s", code, loc)
	}
	if code, loc := authorize("alice", url.Values{"scope": {"admin"}}); code != http.StatusFound || loc.Query().Get("error") != "invalid_scope" {
		t.Fatalf("a scope not the client's: %d %s", code, loc)
	}
	code, loc := authorize("alice", url.Values{"scope": {"download"}})
	if code != http.StatusFound || loc.Host != "portal.example" || loc.Path != "/cb" || loc.Query().Get("code") == "" || loc.Query().Get("state") != "xyz" || loc.Query().Get("error") != "" {
		t.Fatalf("a code: %d %s", code, loc)
	}
	// the code is exchanged once, with the verifier, by the client it was
	// sent to
	exchange := func(form url.Values) (int, map[string]any) {
		t.Helper()
		res, out := e.doWith(func(r *http.Request) {
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Set("Origin", "https://portal.example")
		}, "POST", OAuthTokenPath, []byte(form.Encode()))
		var body map[string]any
		json.Unmarshal(out, &body)
		if res.Header.Get("Access-Control-Allow-Origin") != "https://portal.example" {
			t.Fatalf("CORS of the token endpoint: %q", res.Header.Get("Access-Control-Allow-Origin"))
		}
		return res.StatusCode, body
	}
	form := func(code, verifier string) url.Values {
		return url.Values{"grant_type": {"authorization_code"}, "client_id": {"portal"}, "code": {code}, "code_verifier": {verifier}, "redirect_uri": {"https://portal.example/cb"}}
	}
	authCode := loc.Query().Get("code")
	if code, body := exchange(form(authCode, "wrong-verifier-wrong-verifier-wrong-verifier-wrong")); code != http.StatusBadRequest || body["error"] != "invalid_grant" {
		t.Fatalf("a wrong verifier: %d %+v", code, body)
	}
	// the wrong try spent the code: another is needed
	_, loc = authorize("alice", nil)
	authCode = loc.Query().Get("code")
	f := form(authCode, verifier)
	f.Set("redirect_uri", "http://localhost:5173/cb")
	if code, body := exchange(f); code != http.StatusBadRequest || body["error"] != "invalid_grant" {
		t.Fatalf("another redirect URI: %d %+v", code, body)
	}
	_, loc = authorize("alice", nil)
	authCode = loc.Query().Get("code")
	code, body := exchange(form(authCode, verifier))
	tok, _ := body["access_token"].(string)
	if code != http.StatusOK || tok == "" || body["token_type"] != "Bearer" || body["scope"] != "download print" || body["expires_in"].(float64) != 3600 {
		t.Fatalf("the token: %d %+v", code, body)
	}
	if code, body := exchange(form(authCode, verifier)); code != http.StatusBadRequest || body["error"] != "invalid_grant" {
		t.Fatalf("the code twice: %d %+v", code, body)
	}
	if code, body := exchange(url.Values{"grant_type": {"password"}}); code != http.StatusBadRequest || body["error"] != "unsupported_grant_type" {
		t.Fatalf("another grant: %d %+v", code, body)
	}

	// the token is the person, at the downloaders, the inboxes and /api/me
	// alone
	var me struct {
		Kind, ID, Email string
		Groups          []string
	}
	res, out = e.do(nil, "GET", "/api/me", tok, nil)
	json.Unmarshal(out, &me)
	if res.StatusCode != http.StatusOK || me.Kind != PrincipalUser || me.ID != "u-alice" || me.Email != "alice@example.com" || len(me.Groups) != 1 {
		t.Fatalf("me: %d %+v", res.StatusCode, me)
	}
	var es []DownloadEntry
	res, out = e.doWith(func(r *http.Request) {
		r.Header.Set("Authorization", "Bearer "+tok)
		r.Header.Set("Origin", "https://portal.example")
	}, "GET", "/api/download/out", nil)
	json.Unmarshal(out, &es)
	if res.StatusCode != http.StatusOK || len(es) != 1 || es[0].Owner != "alice" || res.Header.Get("Access-Control-Allow-Origin") != "https://portal.example" {
		t.Fatalf("alice's entries by subject: %d %+v %q", res.StatusCode, es, res.Header.Get("Access-Control-Allow-Origin"))
	}
	if res, out := e.do(nil, "GET", "/api/download/out/1/document", tok, nil); res.StatusCode != http.StatusOK || !bytes.Equal(out, doc) {
		t.Fatalf("the document with the token: %d", res.StatusCode)
	}
	if res, _ := e.do(nil, "GET", "/api/download/out/2", tok, nil); res.StatusCode != http.StatusNotFound {
		t.Fatalf("bob's with the token: %d", res.StatusCode)
	}
	// it opens the person's inboxes too, with the scope inbox alone
	e.node("desk", Node{Kind: KindInbox, Scope: ScopeGroup, Owner: "warehouse-a"}, http.StatusOK)
	if res, out := e.do(nil, "GET", "/api/inbox/desk", tok, nil); res.StatusCode != http.StatusForbidden || !strings.Contains(string(out), "the token has no scope inbox") {
		t.Fatalf("an inbox with a token of no scope inbox: %d %s", res.StatusCode, out)
	}
	for _, c := range []struct{ method, path string }{{"GET", "/api/my/jobs"}, {"GET", "/api/jobs/1"}, {"POST", "/api/destinations/in/jobs?name=x"}, {"GET", "/api/nodes"}, {"GET", "/api/topology"}} {
		if res, _ := e.do(nil, c.method, c.path, tok, doc); res.StatusCode != http.StatusUnauthorized {
			t.Fatalf("%s %s with a person's token: %d", c.method, c.path, res.StatusCode)
		}
	}
	// a page of the client's origin may ask; another origin gets no
	// permission
	res, _ = e.doWith(func(r *http.Request) {
		r.Header.Set("Origin", "https://portal.example")
		r.Header.Set("Access-Control-Request-Method", "GET")
		r.Header.Set("Access-Control-Request-Headers", "authorization")
	}, "OPTIONS", "/api/download/out", nil)
	if res.StatusCode != http.StatusNoContent || res.Header.Get("Access-Control-Allow-Origin") != "https://portal.example" || !strings.Contains(res.Header.Get("Access-Control-Allow-Headers"), "Authorization") {
		t.Fatalf("preflight: %d %v", res.StatusCode, res.Header)
	}
	res, _ = e.doWith(func(r *http.Request) { r.Header.Set("Origin", "https://evil.example") }, "OPTIONS", "/api/download/out", nil)
	if res.StatusCode != http.StatusNoContent || res.Header.Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("preflight of another origin: %d %v", res.StatusCode, res.Header)
	}
	// the scopes of the token are the client's still: taking one away
	// takes it from the token; disabling the client ends it; so does time
	e.admin("PUT", "/api/apps/portal", App{Title: "Portal", Public: true, Scopes: []string{ScopePrint}, RedirectURIs: []string{"https://portal.example/cb"}}, http.StatusOK)
	if res, _ := e.do(nil, "GET", "/api/download/out", tok, nil); res.StatusCode != http.StatusForbidden {
		t.Fatalf("the scope taken away: %d", res.StatusCode)
	}
	e.admin("PUT", "/api/apps/portal", App{Title: "Portal", Public: true, Scopes: []string{ScopeDownload}, RedirectURIs: []string{"https://portal.example/cb"}, Disabled: true}, http.StatusOK)
	if res, _ := e.do(nil, "GET", "/api/download/out", tok, nil); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("the client disabled: %d", res.StatusCode)
	}
	e.admin("PUT", "/api/apps/portal", App{Title: "Portal", Public: true, Scopes: []string{ScopeDownload}, RedirectURIs: []string{"https://portal.example/cb"}}, http.StatusOK)
	if res, _ := e.do(nil, "GET", "/api/download/out", tok, nil); res.StatusCode != http.StatusOK {
		t.Fatalf("the client enabled again: %d", res.StatusCode)
	}
	e.clock.add(2 * time.Hour)
	if res, _ := e.do(nil, "GET", "/api/download/out", tok, nil); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("the token expired: %d", res.StatusCode)
	}
	// the code and the token are in the audit log
	var log []AuditEntry
	json.Unmarshal(e.admin("GET", "/api/audit?limit=20", nil, http.StatusOK), &log)
	actions := map[string]int{}
	for _, en := range log {
		actions[en.Action]++
		if (en.Action == "oauth-code" || en.Action == "user-token") && (en.Who != "Alice" || en.Subject != "portal") {
			t.Fatalf("audit %+v", en)
		}
	}
	if actions["oauth-code"] < 3 || actions["user-token"] != 1 || actions["code-refused"] < 2 {
		t.Fatalf("audit %+v", actions)
	}
	// a public client turned confidential gets a secret; back, none
	var turned struct {
		App
		Secret string `json:"secret"`
	}
	json.Unmarshal(e.admin("PUT", "/api/apps/portal", App{Title: "Portal", Scopes: []string{ScopeDownload}}, http.StatusOK), &turned)
	if turned.Secret == "" || turned.Public || len(turned.RedirectURIs) != 0 {
		t.Fatalf("turned confidential %+v", turned)
	}
	if _, _, _, err := m.RedeemCode(ctx, "nothing", verifier, "portal", "https://portal.example/cb"); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("an unknown code: %v", err)
	}
}

// A filter in front of the archives picks the downloader a document is
// listed at: each archive sends what it wrote, a stocked file, to its
// own. A filter after an archive is refused before a downloader: what it
// sends on is the document, not a stocked file.
func TestDownloadAfterFilter(t *testing.T) {
	e := newEnv(t, withViewers)
	dir := t.TempDir()
	store := func(name string) string { return "bbolt:" + filepath.ToSlash(filepath.Join(dir, name+".db")) }
	e.node("invoices", Node{Kind: KindDownload, By: "customer"}, http.StatusOK)
	e.node("receipts", Node{Kind: KindDownload, By: "customer"}, http.StatusOK)
	e.node("seen", Node{Kind: KindSink}, http.StatusOK)
	e.node("vault-i", Node{Kind: KindArchive, Store: store("i"), To: []string{"invoices"}}, http.StatusOK)
	e.node("vault-r", Node{Kind: KindArchive, Store: store("r"), To: []string{"receipts"}}, http.StatusOK)
	e.node("by-type", Node{Kind: KindFilter, Rules: []Rule{{Match: Match{Meta: map[string]string{"type": "^invoice$"}}, To: []string{"vault-i"}},
		{Match: Match{Meta: map[string]string{"type": "^receipt$"}}, To: []string{"vault-r"}}}, Otherwise: []string{"seen"}}, http.StatusOK)
	e.node("in", Node{Kind: KindQueue, To: []string{"by-type"}}, http.StatusOK)
	e.node("late", Node{Kind: KindFilter, Otherwise: []string{"seen"}}, http.StatusOK)
	e.node("vault-i", Node{Kind: KindArchive, Store: store("i"), To: []string{"invoices", "late"}}, http.StatusOK)
	if out := e.node("late", Node{Kind: KindFilter, Rules: []Rule{{Match: Match{Meta: map[string]string{"type": "^invoice$"}}, To: []string{"receipts"}}}}, http.StatusBadRequest); !strings.Contains(string(out), "not late (filter)") {
		t.Fatalf("a filter after an archive: %s", out)
	}
	doc := testDoc(t, 1)
	for _, q := range []string{"name=a&meta.type=invoice&meta.customer=alice@example.com", "name=b&meta.type=receipt&meta.customer=alice@example.com", "name=c&meta.type=memo"} {
		if res, out := e.do(nil, "POST", "/api/destinations/in/jobs?"+q, "wms:s3cret", doc); res.StatusCode != http.StatusCreated {
			t.Fatalf("%s: %d %s", q, res.StatusCode, out)
		}
	}
	if rep := e.tick(); rep.Delivered != 2 {
		t.Fatalf("tick %+v", rep)
	}
	if l := legs(e.job(1)); l != "in:routed by-type:routed(rule 1) vault-i:done invoices:listed(alice@example.com) late:routed(otherwise) seen:kept" {
		t.Fatalf("invoice: %s", l)
	}
	if l := legs(e.job(3)); l != "in:routed by-type:routed(otherwise) seen:kept" {
		t.Fatalf("memo: %s", l)
	}
	for name, want := range map[string]string{"invoices": "a", "receipts": "b"} {
		_, out := e.as("alice", "GET", "/api/download/"+name, nil)
		var es []DownloadEntry
		json.Unmarshal(out, &es)
		if len(es) != 1 || es[0].Name != want || es[0].Archive != "vault-"+name[:1] || es[0].Inbox != "" {
			t.Fatalf("%s: %+v", name, es)
		}
	}
	// the activity counts what a downloader lists now: one deleted and put
	// again under the same name lists nothing, whatever the legs say
	live := func() LiveNode {
		t.Helper()
		var l Live
		_, out := e.as("root", "GET", "/api/topology/live", nil)
		json.Unmarshal(out, &l)
		return l.Nodes["invoices"]
	}
	if n := live(); n.Listed != 1 {
		t.Fatalf("live before %+v", n)
	}
	e.node("vault-i", Node{Kind: KindArchive, Store: store("i")}, http.StatusOK)
	e.admin("DELETE", "/api/nodes/invoices", nil, http.StatusNoContent)
	e.node("invoices", Node{Kind: KindDownload, By: "customer"}, http.StatusOK)
	if n := live(); n.Listed != 0 {
		t.Fatalf("live after the export was put again %+v", n)
	}
}

// A downloader after an inbox lists each job while the inbox holds it,
// read from the job's documents, with the job as its record. The inbox
// does not release to it; a release or a discard takes the job off. The
// inbox takes any kind still: its edge to the downloader carries a stocked
// file.
func TestDownloadInbox(t *testing.T) {
	e := newEnv(t, withViewers)
	e.node("pull", Node{Kind: KindDownload, By: "customer"}, http.StatusOK)
	e.node("done", Node{Kind: KindSink}, http.StatusOK)
	e.node("hold", Node{Kind: KindInbox, To: []string{"pull", "done"}}, http.StatusOK)
	e.node("in", Node{Kind: KindQueue, To: []string{"hold"}}, http.StatusOK)
	doc := testDoc(t, 1)
	for _, q := range []string{"name=a&meta.customer=alice@example.com", "name=b&meta.customer=alice@example.com"} {
		if res, out := e.do(nil, "POST", "/api/destinations/in/jobs?"+q, "wms:s3cret", doc); res.StatusCode != http.StatusCreated {
			t.Fatalf("%s: %d %s", q, res.StatusCode, out)
		}
	}
	if l := legs(e.job(1)); l != "in:routed hold:held pull:listed(alice@example.com)" {
		t.Fatalf("held: %s", l)
	}
	listed := func() []DownloadEntry {
		t.Helper()
		res, out := e.as("alice", "GET", "/api/download/pull", nil)
		if res.StatusCode != http.StatusOK {
			t.Fatalf("list: %d %s", res.StatusCode, out)
		}
		var es []DownloadEntry
		json.Unmarshal(out, &es)
		return es
	}
	if es := listed(); len(es) != 2 || es[1].Inbox != "hold" || es[1].Item == "" || es[1].Archive != "" || es[1].Name != "a" {
		t.Fatalf("listed %+v", es)
	}
	if res, out := e.as("alice", "GET", "/api/download/pull/1/document", nil); res.StatusCode != http.StatusOK || !bytes.Equal(out, doc) {
		t.Fatalf("document: %d %s", res.StatusCode, out)
	}
	var rec Job
	if res, out := e.as("alice", "GET", "/api/download/pull/1/record", nil); res.StatusCode != http.StatusOK || json.Unmarshal(out, &rec) != nil || rec.Number != 1 {
		t.Fatalf("record: %d %s", res.StatusCode, out)
	}
	// the inbox releases to done alone
	var v InboxView
	json.Unmarshal(e.admin("GET", "/api/inbox/hold", nil, http.StatusOK), &v)
	if strings.Join(v.To, ",") != "done" || len(v.Held) != 2 || v.Held[0].Job != 1 || strings.Join(v.Held[0].Listed, ",") != "pull" {
		t.Fatalf("inbox %+v", v)
	}
	if out := e.admin("POST", "/api/inbox/hold/"+fmt.Sprint(v.Held[0].Job)+"/release", map[string]any{"to": []string{"pull"}}, http.StatusBadRequest); !strings.Contains(string(out), "not a node to release to") {
		t.Fatalf("release to the export: %s", out)
	}
	e.admin("POST", "/api/inbox/hold/"+fmt.Sprint(v.Held[0].Job)+"/release", nil, http.StatusOK)
	if l := legs(e.job(1)); l != "in:routed hold:released pull:listed(alice@example.com) done:kept" {
		t.Fatalf("released: %s", l)
	}
	if es := listed(); len(es) != 1 || es[0].Job != 2 {
		t.Fatalf("after the release %+v", es)
	}
	for _, path := range []string{"/api/download/pull/1", "/api/download/pull/1/document", "/api/download/pull/1/record"} {
		if res, _ := e.as("alice", "GET", path, nil); res.StatusCode != http.StatusNotFound {
			t.Fatalf("%s after the release: %d", path, res.StatusCode)
		}
	}
	// a discard takes it off as well, under the downloader's new name
	e.admin("POST", "/api/nodes/pull/rename", map[string]string{"name": "fetch"}, http.StatusOK)
	e.admin("POST", "/api/inbox/hold/"+fmt.Sprint(v.Held[1].Job)+"/discard", nil, http.StatusOK)
	if res, out := e.as("alice", "GET", "/api/download/fetch", nil); res.StatusCode != http.StatusOK || string(bytes.TrimSpace(out)) != "[]" {
		t.Fatalf("after the discard: %d %s", res.StatusCode, out)
	}
	var typed []TopoNode
	json.Unmarshal(e.admin("GET", "/api/topology", nil, http.StatusOK), &typed)
	for _, n := range typed {
		if n.Name == "hold" && (len(n.Accepts) != 3 || len(n.Broken) > 0) || n.Name == "fetch" && (n.Takes != PayloadStocked || n.Listed != 0) {
			t.Fatalf("typing of %s: %+v", n.Name, n)
		}
	}
}
