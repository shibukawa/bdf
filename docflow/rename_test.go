package docflow

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// A node is put before the nodes after it: a job that reaches one that
// sends to none yet gets an error there.
func TestNowhere(t *testing.T) {
	e := newEnv(t, nil)
	e.node("empty", Node{Kind: KindQueue}, http.StatusOK)
	e.printJob("empty", "lost", 1)
	j := e.job(1)
	if legs(j) != "empty:error" || !strings.Contains(j.Legs[0].Error, "no node yet") {
		t.Fatalf("job %+v: %s", j, legs(j))
	}
}

// A node renamed takes what is kept under its name with it, and the
// nodes that send to it send to the new name.
func TestRenameNode(t *testing.T) {
	e := newEnv(t, nil)
	ctx := context.Background()
	rename := func(from, to string, status int) []byte {
		t.Helper()
		return e.admin("POST", "/api/nodes/"+from+"/rename", map[string]string{"name": to}, status)
	}
	e.display("wh-a", Node{Title: "A"})
	d := e.enroll("wh-a", "tablet")
	e.node("q", Node{Kind: KindQueue, To: []string{"wh-a"}}, http.StatusOK)
	e.node("f", Node{Kind: KindFilter, Rules: []Rule{{Match: Match{Name: "^x"}, To: []string{"wh-a"}}}, Otherwise: []string{"q"}}, http.StatusOK)
	e.admin("PUT", "/api/layout", Layout{Positions: map[string]Point{"q": {X: 10, Y: 20}, "wh-a": {X: 300}}}, http.StatusOK)

	// a display with a slip open is refused, until it is done
	e.printJob("q", "pick", 1)
	s := d.board().Open[0]
	if out := rename("wh-a", "wh-1", http.StatusConflict); !strings.Contains(string(out), "open slips") {
		t.Fatalf("open slip: %s", out)
	}
	d.action(s.ID, "done", http.StatusOK)
	var en Enrollment
	json.Unmarshal(e.admin("POST", "/api/nodes/wh-a/enrollments", nil, http.StatusCreated), &en)
	rename("wh-a", "wh-1", http.StatusOK)
	e.admin("GET", "/api/nodes/wh-a", nil, http.StatusNotFound)
	if q, _ := e.m().Node(ctx, "q"); !slices.Equal(q.To, []string{"wh-1"}) {
		t.Fatalf("q %+v", q)
	}
	if f, _ := e.m().Node(ctx, "f"); !slices.Equal(f.Rules[0].To, []string{"wh-1"}) {
		t.Fatalf("f %+v", f)
	}
	// the device signs in again (its session named the old name) and
	// still sees the slip it did
	d.login()
	if b := d.board(); b.Location != "wh-1" || len(b.Ended) != 1 || b.Ended[0].ID != s.ID {
		t.Fatalf("board %+v", b)
	}
	if rec, err := e.m().st.db.Get[enrollRecord](ctx, en.Code); err != nil || rec.Data.Location != "wh-1" {
		t.Fatalf("enrollment %+v %v", rec.Data, err)
	}

	rename("q", "queue-1", http.StatusOK)
	if f, _ := e.m().Node(ctx, "f"); !slices.Equal(f.Otherwise, []string{"queue-1"}) {
		t.Fatalf("f %+v", f)
	}
	var l Layout
	json.Unmarshal(e.admin("GET", "/api/layout", nil, http.StatusOK), &l)
	if _, ok := l.Positions["q"]; ok || l.Positions["queue-1"] != (Point{X: 10, Y: 20}) || l.Positions["wh-1"].X != 300 {
		t.Fatalf("layout %+v", l)
	}
	e.printJob("queue-1", "pick 2", 1)
	if b := d.board(); len(b.Open) != 1 || b.Open[0].Name != "pick 2" {
		t.Fatalf("board %+v", b)
	}

	// the name is checked
	rename("queue-1", "f", http.StatusConflict)
	rename("queue-1", "has space", http.StatusBadRequest)
	rename("nowhere", "x", http.StatusNotFound)
	rename("queue-1", "queue-1", http.StatusOK)

	// a template node's templates go with it
	tpl, err := os.ReadFile("../testdata/form/pages.bdft")
	if err != nil {
		t.Fatal(err)
	}
	e.node("forms", Node{Kind: KindTemplate, Template: "slip", To: []string{"queue-1"}}, http.StatusOK)
	e.admin("PUT", "/api/nodes/forms/templates/slip?name=pages.bdft", tpl, http.StatusOK)
	rename("forms", "forms-2", http.StatusOK)
	if got := e.admin("GET", "/api/nodes/forms-2/templates/slip/file", nil, http.StatusOK); !bytes.Equal(got, tpl) {
		t.Fatalf("template file: %d bytes", len(got))
	}
	e.admin("GET", "/api/nodes/forms/templates/slip", nil, http.StatusNotFound)
	if ts, _ := e.m().templatesOf(ctx, "forms"); len(ts) != 0 {
		t.Fatalf("templates left %+v", ts)
	}

	// an archive keeps writing to its collection
	store := "bbolt:" + filepath.ToSlash(filepath.Join(t.TempDir(), "audit.db"))
	e.node("audit", Node{Kind: KindArchive, Store: store}, http.StatusOK)
	rename("audit", "vault", http.StatusOK)
	if n, _ := e.m().Node(ctx, "vault"); n.Collection != "archive-audit" {
		t.Fatalf("collection %q", n.Collection)
	}
}
