package docflow

import (
	"bytes"
	"encoding/json"
	"image/png"
	"net/http"
	"slices"
	"strings"
	"testing"
)

// The inboxes' API: an app the inbox lists with the scope inbox, the
// people of a group's inbox and the role topology list what it holds,
// look at its pages and release or discard it; anyone else is refused.
func TestInboxAPI(t *testing.T) {
	e := newEnv(t, withViewers)
	secrets := map[string]string{}
	for id, scopes := range map[string][]string{"clerk": {ScopeInbox}, "reader": {ScopeDownload}, "outsider": {ScopeInbox}} {
		var made struct {
			Secret string `json:"secret"`
		}
		json.Unmarshal(e.admin("PUT", "/api/apps/"+id, App{Scopes: scopes}, http.StatusOK), &made)
		secrets[id] = made.Secret
	}
	token := func(id string) string {
		t.Helper()
		res, out := e.doWith(func(r *http.Request) {
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			r.SetBasicAuth(id, secrets[id])
		}, "POST", OAuthTokenPath, []byte("grant_type=client_credentials"))
		var body struct {
			Token string `json:"access_token"`
		}
		if json.Unmarshal(out, &body); res.StatusCode != http.StatusOK || body.Token == "" {
			t.Fatalf("token of %s: %d %s", id, res.StatusCode, out)
		}
		return body.Token
	}
	clerk, reader, outsider := token("clerk"), token("reader"), token("outsider")

	e.display("wh-a", Node{})
	e.printer("laser")
	// a group's inbox that releases to two nodes, one of the same group
	// that names none, and a shared one; the apps an inbox lists are apps
	e.node("desk", Node{Kind: KindInbox, Title: "Approvals", To: []string{"wh-a", "laser"}, Scope: ScopeGroup, Owner: "warehouse-a", Apps: []string{"clerk", "reader"}}, http.StatusOK)
	e.node("free", Node{Kind: KindInbox, Scope: ScopeGroup, Owner: "warehouse-a"}, http.StatusOK)
	e.node("shared", Node{Kind: KindInbox, To: []string{"wh-a"}}, http.StatusOK)
	if out := e.node("bad", Node{Kind: KindInbox, Apps: []string{"nobody"}}, http.StatusBadRequest); !strings.Contains(string(out), `no app \"nobody\"`) {
		t.Fatalf("an app that is not there: %s", out)
	}
	doc := testDoc(t, 2)
	for _, dest := range []string{"desk", "desk", "free", "shared", "desk"} {
		if res, out := e.do(nil, "POST", "/api/destinations/"+dest+"/jobs?name=to-"+dest, "wms:s3cret", doc); res.StatusCode != http.StatusCreated {
			t.Fatalf("send to %s: %d %s", dest, res.StatusCode, out)
		}
	}

	// the inboxes each works
	names := func(res *http.Response, out []byte) string {
		t.Helper()
		var ins []InboxSummary
		if res.StatusCode != http.StatusOK || json.Unmarshal(out, &ins) != nil {
			t.Fatalf("inboxes: %d %s", res.StatusCode, out)
		}
		var ns []string
		for _, in := range ins {
			ns = append(ns, in.Name)
		}
		return strings.Join(ns, " ")
	}
	var ins []InboxSummary
	_, out := e.as("alice", "GET", "/api/inbox", nil)
	if json.Unmarshal(out, &ins); len(ins) != 2 || ins[0].Name != "desk" || ins[0].Title != "Approvals" || ins[0].Held != 3 || !slices.Equal(ins[0].To, []string{"wh-a", "laser"}) ||
		ins[1].Name != "free" || ins[1].Held != 1 || len(ins[1].To) != 0 {
		t.Fatalf("alice's inboxes %+v", ins)
	}
	for who, want := range map[string]string{"bob": "", "root": "desk free shared"} {
		if got := names(e.as(who, "GET", "/api/inbox", nil)); got != want {
			t.Fatalf("%s's inboxes %q; want %q", who, got, want)
		}
	}
	for tok, want := range map[string]string{clerk: "desk", reader: "", outsider: ""} {
		if got := names(e.do(nil, "GET", "/api/inbox", tok, nil)); got != want {
			t.Fatalf("an app's inboxes %q; want %q", got, want)
		}
	}
	if res, _ := e.do(nil, "GET", "/api/inbox", "", nil); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("nobody's inboxes: %d", res.StatusCode)
	}

	// an inbox: who works it reads it, the others are told why not
	res, out := e.as("alice", "GET", "/api/inbox/desk", nil)
	var v InboxView
	if json.Unmarshal(out, &v); res.StatusCode != http.StatusOK || len(v.Held) != 3 || v.Held[0].Job != 1 || v.Held[0].Name != "to-desk" || v.Held[0].Pages != 2 || res.Header.Get("ETag") == "" {
		t.Fatalf("the inbox: %d %+v", res.StatusCode, v)
	}
	etag := res.Header.Get("ETag")
	if res, _ := e.doWith(func(r *http.Request) {
		r.Header.Set("Authorization", "Bearer "+clerk)
		r.Header.Set("If-None-Match", etag)
	}, "GET", "/api/inbox/desk", nil); res.StatusCode != http.StatusNotModified {
		t.Fatalf("the inbox unchanged: %d", res.StatusCode)
	}
	for _, c := range []struct {
		name string
		res  *http.Response
		out  []byte
		want int
		why  string
	}{
		{"bob", nil, nil, http.StatusForbidden, "you do not work the inbox desk"},
		{"reader", nil, nil, http.StatusForbidden, "has no scope inbox"},
		{"outsider", nil, nil, http.StatusForbidden, "does not take the app outsider"},
		{"wms", nil, nil, http.StatusForbidden, "you do not work"},
	} {
		switch c.name {
		case "bob":
			c.res, c.out = e.as("bob", "GET", "/api/inbox/desk/1", nil)
		case "reader":
			c.res, c.out = e.do(nil, "GET", "/api/inbox/desk", reader, nil)
		case "outsider":
			c.res, c.out = e.do(nil, "GET", "/api/inbox/desk/1/thumbnail", outsider, nil)
		case "wms":
			c.res, c.out = e.do(nil, "GET", "/api/inbox/desk", "wms:s3cret", nil)
		}
		if c.res.StatusCode != c.want || !strings.Contains(string(c.out), c.why) {
			t.Fatalf("%s: %d %s", c.name, c.res.StatusCode, c.out)
		}
	}
	if res, _ := e.as("alice", "GET", "/api/inbox/shared", nil); res.StatusCode != http.StatusForbidden {
		t.Fatalf("a shared inbox is the role topology's: %d", res.StatusCode)
	}
	for _, path := range []string{"/api/inbox/wh-a", "/api/inbox/desk/3", "/api/inbox/desk/x", "/api/inbox/desk/99/thumbnail"} {
		if res, _ := e.do(nil, "GET", path, clerk, nil); res.StatusCode != http.StatusNotFound {
			t.Fatalf("%s: %d", path, res.StatusCode)
		}
	}

	// an item, its thumbnail, its pages and its document
	var it Item
	if _, out := e.do(nil, "GET", "/api/inbox/desk/2", clerk, nil); json.Unmarshal(out, &it) != nil || it.Job != 2 || it.Inbox != "desk" || it.State != ItemHeld {
		t.Fatalf("item %s", out)
	}
	image := func(path string) (w, h int) {
		t.Helper()
		res, out := e.do(nil, "GET", path, clerk, nil)
		cfg, err := png.DecodeConfig(bytes.NewReader(out))
		if res.StatusCode != http.StatusOK || res.Header.Get("Content-Type") != "image/png" || err != nil {
			t.Fatalf("%s: %d %v %.200s", path, res.StatusCode, err, out)
		}
		return cfg.Width, cfg.Height
	}
	// a page of A4: the whole of it, by default, its longer side the size
	if w, h := image("/api/inbox/desk/1/thumbnail"); h != 256 || w >= h {
		t.Fatalf("thumbnail %dx%d", w, h)
	}
	if w, h := image("/api/inbox/desk/1/thumbnail?size=64&mode=crop"); w != 64 || h != 64 {
		t.Fatalf("cropped thumbnail %dx%d", w, h)
	}
	if w, _ := image("/api/inbox/desk/1/pages/2.png?w=200"); w != 200 {
		t.Fatalf("page 2 is %d wide", w)
	}
	for path, want := range map[string]int{"/api/inbox/desk/1/thumbnail?mode=tiles": http.StatusBadRequest, "/api/inbox/desk/1/pages/3": http.StatusNotFound} {
		if res, _ := e.do(nil, "GET", path, clerk, nil); res.StatusCode != want {
			t.Fatalf("%s: %d; want %d", path, res.StatusCode, want)
		}
	}
	if res, out := e.do(nil, "GET", "/api/inbox/desk/1/document", clerk, nil); res.StatusCode != http.StatusOK || !bytes.Equal(out, doc) ||
		res.Header.Get("Content-Disposition") != `attachment; filename="to-desk.bdf"` {
		t.Fatalf("document: %d %q", res.StatusCode, res.Header.Get("Content-Disposition"))
	}

	// the decision: the inbox's nodes, by whoever works it
	if res, _ := e.as("bob", "POST", "/api/inbox/desk/1/release", nil); res.StatusCode != http.StatusForbidden {
		t.Fatalf("bob released: %d", res.StatusCode)
	}
	if res, out := e.as("alice", "POST", "/api/inbox/desk/1/release", map[string]any{"to": []string{"shared"}}); res.StatusCode != http.StatusBadRequest || !strings.Contains(string(out), "releases to wh-a, laser only") {
		t.Fatalf("to another node: %d %s", res.StatusCode, out)
	}
	// a page of another origin may not decide with the session
	if res, _ := e.doWith(func(r *http.Request) {
		r.Header.Set("X-Test-Viewer", "alice")
		r.Header.Set("Sec-Fetch-Site", "cross-site")
	}, "POST", "/api/inbox/desk/1/release", nil); res.StatusCode != http.StatusForbidden {
		t.Fatalf("released from another origin: %d", res.StatusCode)
	}
	res, out = e.as("alice", "POST", "/api/inbox/desk/1/release", map[string]any{"to": []string{"wh-a"}})
	if json.Unmarshal(out, &it); res.StatusCode != http.StatusOK || it.State != ItemReleased || it.Released.Worker != "Alice" || !slices.Equal(it.Released.To, []string{"wh-a"}) {
		t.Fatalf("released: %d %s", res.StatusCode, out)
	}
	if res, _ := e.as("alice", "POST", "/api/inbox/desk/1/release", nil); res.StatusCode != http.StatusConflict {
		t.Fatalf("released twice: %d", res.StatusCode)
	}
	res, out = e.do(nil, "POST", "/api/inbox/desk/2/release", clerk, map[string]any{"worker": "Sato"})
	if json.Unmarshal(out, &it); res.StatusCode != http.StatusOK || it.Released.Worker != "Sato" || !slices.Equal(it.Released.To, []string{"wh-a", "laser"}) {
		t.Fatalf("released by the app: %d %s", res.StatusCode, out)
	}
	if j := e.job(2); legs(j) != "desk:released wh-a:shown laser:pending" {
		t.Fatalf("job 2: %s", legs(j))
	}
	res, out = e.as("alice", "POST", "/api/inbox/desk/5/discard", map[string]any{"worker": "Mina"})
	if json.Unmarshal(out, &it); res.StatusCode != http.StatusOK || it.State != ItemDiscarded || it.Discarded.Worker != "Mina" {
		t.Fatalf("discarded: %d %s", res.StatusCode, out)
	}
	// an inbox that names no nodes releases where the role topology says
	if res, out := e.as("alice", "POST", "/api/inbox/free/3/release", map[string]any{"to": []string{"wh-a"}}); res.StatusCode != http.StatusForbidden || !strings.Contains(string(out), "names no nodes to release to") {
		t.Fatalf("a member released anywhere: %d %s", res.StatusCode, out)
	}
	if res, out := e.as("root", "POST", "/api/inbox/free/3/release", map[string]any{"to": []string{"wh-a"}}); res.StatusCode != http.StatusOK {
		t.Fatalf("root released: %d %s", res.StatusCode, out)
	}
	if _, out := e.as("alice", "GET", "/api/inbox/desk", nil); json.Unmarshal(out, &v) != nil || len(v.Held) != 0 || len(v.Ended) != 3 {
		t.Fatalf("the inbox after: %s", out)
	}
	// who decided, and the document read, are in the audit log
	var log []AuditEntry
	json.Unmarshal(e.admin("GET", "/api/audit?limit=10", nil, http.StatusOK), &log)
	got := map[string]string{}
	for _, en := range log {
		got[en.Action+" "+en.Subject] = en.Who + ": " + en.Detail
	}
	if got["release-item "+itemID(2, "desk")] != "clerk: wh-a, laser by Sato" || got["discard-item "+itemID(5, "desk")] != "Alice: Mina" ||
		!strings.HasPrefix(got["inbox-document desk"], "clerk: job 1 ") {
		t.Fatalf("audit %+v", got)
	}

	// the OpenAPI document shows each the inboxes they work
	paths := func(tok string) []string {
		t.Helper()
		_, out := e.do(nil, "GET", OpenAPIPath, tok, nil)
		var doc map[string]any
		json.Unmarshal(out, &doc)
		return keysOf(doc["paths"])
	}
	want := []string{"/api/inbox", "/api/inbox/desk", "/api/inbox/desk/{job}", "/api/inbox/desk/{job}/discard", "/api/inbox/desk/{job}/document",
		"/api/inbox/desk/{job}/pages/{page}", "/api/inbox/desk/{job}/release", "/api/inbox/desk/{job}/thumbnail"}
	if got := paths(clerk); !slices.Equal(got, want) {
		t.Fatalf("the clerk's paths %v", got)
	}
	if got := paths(reader); len(got) != 0 {
		t.Fatalf("the reader's paths %v", got)
	}
	_, out = e.do(nil, "GET", OpenAPIPath+"?node=desk", clerk, nil)
	var api map[string]any
	json.Unmarshal(out, &api)
	if enum := dig(api, "paths", "/api/inbox/desk/{job}/release", "post", "requestBody", "content", "application/json", "schema", "properties", "to", "items", "enum"); !slices.Equal(anyStrings(enum), []string{"wh-a", "laser"}) {
		t.Fatalf("release to %v", enum)
	}
	// the topology says where the API reaches a node
	var topo []TopoNode
	json.Unmarshal(e.admin("GET", "/api/topology", nil, http.StatusOK), &topo)
	for _, n := range topo {
		if n.Name == "desk" && n.API != "/api/inbox/desk" || n.Name == "wh-a" && n.API != "" {
			t.Fatalf("%s at %q", n.Name, n.API)
		}
	}
}

// anyStrings is a JSON array of strings.
func anyStrings(v any) []string {
	var out []string
	if xs, ok := v.([]any); ok {
		for _, x := range xs {
			s, _ := x.(string)
			out = append(out, s)
		}
	}
	return out
}
