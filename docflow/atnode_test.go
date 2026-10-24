package docflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
)

// What is at a display, an inbox and an archive is seen on the topology,
// as a sink's is: the slips on a display's board, the latest documents an
// archive wrote, and the first page of each on the node (Live.Thumb), each
// by those who see the job (a display's every one by the role topology).
func TestAtNodes(t *testing.T) {
	topo := &Viewer{Kind: PrincipalUser, Via: ViaOIDC, ID: "u-tom", Name: "Tom", Email: "tom@example.com", Roles: []string{RoleTopology}}
	e := newEnv(t, func(o *Options) {
		withViewers(o)
		o.Viewer = func(r *http.Request) *Viewer {
			if r.Header.Get("X-Test-Viewer") == "tom" {
				return topo
			}
			return viewers[r.Header.Get("X-Test-Viewer")]
		}
	})
	ctx := context.Background()
	store := "bbolt:" + filepath.ToSlash(filepath.Join(t.TempDir(), "vault.db"))
	e.display("dock", Node{Title: "Dock"})
	e.node("vault", Node{Kind: KindArchive, Store: store}, http.StatusOK)
	e.node("hold", Node{Kind: KindInbox, To: []string{"dock"}}, http.StatusOK)
	e.node("in", Node{Kind: KindQueue, To: []string{"dock", "vault", "hold"}}, http.StatusOK)
	doc := testDoc(t, 2)
	for _, s := range []struct{ who, name string }{{"alice", "alice-1"}, {"wms", "wms-2"}, {"alice", "alice-3"}} {
		path := "/api/destinations/in/jobs?name=" + s.name
		var res *http.Response
		var out []byte
		if s.who == "wms" {
			res, out = e.do(nil, "POST", path, "wms:s3cret", doc)
		} else {
			res, out = e.as(s.who, "POST", path, doc)
		}
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("%s: %d %s", s.name, res.StatusCode, out)
		}
	}
	if rep := e.tick(); rep.Delivered != 3 {
		t.Fatalf("tick %+v", rep)
	}
	get := func(who, path string, status int) []byte {
		t.Helper()
		res, out := e.as(who, "GET", path, nil)
		if res.StatusCode != status {
			t.Fatalf("%s GET %s: %d %s; want %d", who, path, res.StatusCode, out, status)
		}
		return out
	}

	// the board: the role topology sees every slip, a person those of the
	// jobs they see
	board := func(who string) []string {
		t.Helper()
		var b Board
		json.Unmarshal(get(who, "/api/nodes/dock/board", http.StatusOK), &b)
		var names []string
		for _, s := range b.Open {
			names = append(names, s.Name)
		}
		return names
	}
	for who, want := range map[string]string{"tom": "alice-1 wms-2 alice-3", "root": "alice-1 wms-2 alice-3", "alice": "alice-1 alice-3", "bob": ""} {
		if got := strings.Join(board(who), " "); got != want {
			t.Fatalf("%s's board: %q; want %q", who, got, want)
		}
	}
	get("alice", "/api/nodes/vault/board", http.StatusNotFound)
	if res, _ := e.do(nil, "GET", "/api/nodes/dock/board", "", nil); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("the board signed out: %d", res.StatusCode)
	}
	// a slip's pages and its document, to those who see it
	mine, theirs := slipID(1, "dock"), slipID(2, "dock")
	out := get("alice", "/api/slips/"+mine+"/pages/2.png?w=200", http.StatusOK)
	if img, err := png.Decode(bytes.NewReader(out)); err != nil || img.Bounds().Dx() != 200 {
		t.Fatalf("slip page: %v", err)
	}
	get("alice", "/api/slips/"+mine+"/pages/3.png", http.StatusNotFound)
	get("alice", "/api/slips/"+theirs+"/pages/1.png", http.StatusNotFound)
	get("tom", "/api/slips/"+theirs+"/pages/1.png", http.StatusOK)
	if out := get("alice", "/api/slips/"+mine+"/document", http.StatusOK); !bytes.Equal(out, doc) {
		t.Fatal("the slip's document is not what was sent")
	}
	get("bob", "/api/slips/"+mine+"/document", http.StatusNotFound)
	get("alice", "/api/slips/"+mine, http.StatusForbidden) // the slip itself is the role topology's

	// the archive: what it wrote, newest first, of the jobs the viewer sees
	archived := func(who, query string) []ArchivedEntry {
		t.Helper()
		var es []ArchivedEntry
		json.Unmarshal(get(who, "/api/nodes/vault/archived"+query, http.StatusOK), &es)
		return es
	}
	es := archived("root", "")
	if len(es) != 3 || es[0].Job != 3 || es[2].Job != 1 || es[0].Pages != 2 || es[0].Kind != PayloadBDF || es[0].Object != es[0].Key+".bdf" || es[0].Archived.IsZero() {
		t.Fatalf("root's archived %+v", es)
	}
	if es := archived("alice", ""); len(es) != 2 || es[0].Job != 3 || es[1].Job != 1 {
		t.Fatalf("alice's archived %+v", es)
	}
	if es := archived("alice", "?limit=1"); len(es) != 1 || es[0].Job != 3 {
		t.Fatalf("limit %+v", es)
	}
	if es := archived("bob", ""); len(es) != 0 {
		t.Fatalf("bob's archived %+v", es)
	}
	if es := archived("tom", ""); len(es) != 0 {
		t.Fatalf("the role topology sees the jobs it sees: %+v", es)
	}
	get("alice", "/api/nodes/dock/archived", http.StatusNotFound)
	key, wms := es[0].Key, es[1].Key
	out = get("alice", "/api/nodes/vault/archived/"+key+"/pages/1.png?w=120", http.StatusOK)
	if img, err := png.Decode(bytes.NewReader(out)); err != nil || img.Bounds().Dx() != 120 {
		t.Fatalf("archived page: %v", err)
	}
	get("alice", "/api/nodes/vault/archived/"+wms+"/pages/1.png", http.StatusNotFound)
	get("alice", "/api/nodes/vault/archived/nothing/pages/1.png", http.StatusNotFound)
	if out := get("alice", "/api/nodes/vault/archived/"+key+"/document", http.StatusOK); !bytes.Equal(out, doc) {
		t.Fatal("the archived document is not what was sent")
	}
	var audit []AuditEntry
	json.Unmarshal(e.admin("GET", "/api/audit", nil, http.StatusOK), &audit)
	if len(audit) == 0 || audit[0].Action != "archive-document" || audit[0].Subject != "vault" || !strings.HasPrefix(audit[0].Detail, "job 3 ") {
		t.Fatalf("audit %+v", audit[:min(len(audit), 1)])
	}

	// the first page of each on the node, for the viewer: the newest of
	// theirs
	live := func(who string) map[string]LiveNode {
		t.Helper()
		var l Live
		json.Unmarshal(get(who, "/api/topology/live", http.StatusOK), &l)
		return l.Nodes
	}
	l := live("alice")
	for name, want := range map[string]string{
		"dock":  "/api/slips/" + slipID(3, "dock") + "/pages/1.png",
		"hold":  "/api/jobs/" + jobID(3) + "/pages/1.png?payload=",
		"vault": "/api/nodes/vault/archived/" + key + "/pages/1.png",
	} {
		if !strings.HasPrefix(l[name].Thumb, want) {
			t.Fatalf("%s's thumb %q; want %q", name, l[name].Thumb, want)
		}
		get("alice", l[name].Thumb, http.StatusOK)
	}
	if l := live("bob"); l["dock"].Thumb != "" || l["vault"].Thumb != "" || l["hold"].Thumb != "" {
		t.Fatalf("bob's thumbs %+v", l)
	}
	if l := live("root"); !strings.Contains(l["vault"].Thumb, key) {
		t.Fatalf("root's archive thumb %q", l["vault"].Thumb)
	}

	// a write again (a retry) takes its place; past the latest, the oldest
	// rolls off
	m := e.m()
	for i := range recentArchived + 2 {
		if err := m.noteArchived(ctx, "vault", ArchivedEntry{Job: int32(100 + i), Kind: PayloadBDF, Key: fmt.Sprintf("k%03d", i)}); err != nil {
			t.Fatal(err)
		}
	}
	m.noteArchived(ctx, "vault", ArchivedEntry{Job: 100 + recentArchived, Kind: PayloadBDF, Key: fmt.Sprintf("k%03d", recentArchived)})
	st := m.archivedState(ctx, "vault")
	if len(st.Entries) != recentArchived || st.Entries[0].Key != "k002" || st.Entries[len(st.Entries)-1].Key != fmt.Sprintf("k%03d", recentArchived) {
		t.Fatalf("latest %d: %s … %s", len(st.Entries), st.Entries[0].Key, st.Entries[len(st.Entries)-1].Key)
	}
	if es := archived("root", "?limit=500"); len(es) != recentArchived {
		t.Fatalf("limit past the latest: %d", len(es))
	}

	// renamed, the archive's latest go with it; deleted, they go
	e.admin("PUT", "/api/nodes/in", Node{Kind: KindQueue, To: []string{"dock", "hold"}}, http.StatusOK)
	e.admin("POST", "/api/nodes/vault/rename", map[string]string{"name": "vault2"}, http.StatusOK)
	if st := m.archivedState(ctx, "vault2"); len(st.Entries) != recentArchived {
		t.Fatalf("renamed: %d", len(st.Entries))
	}
	if st := m.archivedState(ctx, "vault"); len(st.Entries) != 0 {
		t.Fatalf("left under the old name: %d", len(st.Entries))
	}
	e.admin("DELETE", "/api/nodes/vault2", nil, http.StatusNoContent)
	if st := m.archivedState(ctx, "vault2"); len(st.Entries) != 0 {
		t.Fatalf("kept after the archive went: %d", len(st.Entries))
	}
}
