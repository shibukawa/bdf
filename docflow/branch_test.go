package docflow

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

// A branch sends a job along the first case it passes, else its default,
// else ends it; its arms are edges as the others are.
func TestBranch(t *testing.T) {
	e := newEnv(t, nil)
	for _, s := range []string{"invoices", "big", "rest", "other"} {
		e.node(s, Node{Kind: KindSink}, http.StatusOK)
	}
	for _, bad := range []Node{
		{Cases: []Case{{Name: "x", To: []string{"rest"}}}},                                              // checks nothing
		{Cases: []Case{{Name: "default", When: []Check{{Key: "a"}}}}},                                   // the default's name
		{Cases: []Case{{Name: "x", When: []Check{{Key: "a"}}}, {Name: "x", When: []Check{{Key: "b"}}}}}, // twice
		{Cases: []Case{{When: []Check{{Key: "a", Pattern: "["}}}}},
		{Cases: []Case{{When: []Check{{Key: "a"}}, To: []string{"nowhere"}}}},
	} {
		bad.Kind = KindBranch
		e.node("route", bad, http.StatusBadRequest)
	}
	// put before the nodes after it
	e.node("route", Node{Kind: KindBranch}, http.StatusOK)
	e.node("route", Node{Kind: KindBranch, Cases: []Case{
		{Name: "invoice", When: []Check{{Key: "doc-type", Enum: []string{"invoice"}}}, To: []string{"invoices"}},
		{When: []Check{{Key: "amount", Min: ptr(10000)}, {Key: "currency", Enum: []string{"JPY"}}}, To: []string{"big", "invoices"}},
		{Name: "pending", When: []Check{{Key: "doc-type", Enum: []string{"draft"}}}},
	}, Default: []string{"rest"}}, http.StatusOK)
	e.node("in", Node{Kind: KindQueue, To: []string{"route"}}, http.StatusOK)
	send := func(query string) *Job {
		t.Helper()
		res, out := e.do(nil, "POST", "/api/destinations/in/jobs?name=x&"+query, "", testDoc(t, 1))
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("%s: %d %s", query, res.StatusCode, out)
		}
		var j Job
		json.Unmarshal(out, &j)
		return &j
	}
	for _, c := range []struct{ query, want string }{
		{"meta.doc-type=invoice&meta.amount=99999&meta.currency=JPY", "in:routed route:routed(invoice) invoices:kept"},
		{"meta.doc-type=order&meta.amount=99999&meta.currency=JPY", "in:routed route:routed(case 2) big:kept invoices:kept"},
		{"meta.doc-type=order&meta.amount=99999&meta.currency=USD", "in:routed route:routed(default) rest:kept"},
		{"meta.amount=lots&meta.currency=JPY", "in:routed route:routed(default) rest:kept"},
		{"meta.doc-type=draft", "in:routed route:error(pending)"},
	} {
		if j := send(c.query); legs(j) != c.want {
			t.Errorf("%s: %s; want %s", c.query, legs(j), c.want)
		}
	}
	// no default: a job no case fits ends there
	e.node("route", Node{Kind: KindBranch, Cases: []Case{{When: []Check{{Key: "doc-type"}}, To: []string{"invoices"}}}}, http.StatusOK)
	if j := send(""); legs(j) != "in:routed route:dropped" || j.State != JobCompleted || !strings.HasPrefix(j.Legs[1].Error, "no case fits") {
		t.Errorf("dropped %s: %s", j.State, legs(j))
	}

	// the arms are edges: a node they name is not removed, a rename
	// follows them, a cycle through them and a kind the node after does
	// not take are refused
	e.node("route", Node{Kind: KindBranch, Cases: []Case{{Name: "a", When: []Check{{Key: "x"}}, To: []string{"invoices"}}}, Default: []string{"rest"}}, http.StatusOK)
	e.admin("DELETE", "/api/nodes/invoices", nil, http.StatusConflict)
	e.admin("DELETE", "/api/nodes/rest", nil, http.StatusConflict)
	e.admin("POST", "/api/nodes/invoices/rename", map[string]string{"name": "bills"}, http.StatusOK)
	var n Node
	json.Unmarshal(e.admin("GET", "/api/nodes/route", nil, http.StatusOK), &n)
	if n.Cases[0].To[0] != "bills" {
		t.Errorf("renamed %+v", n.Cases)
	}
	e.node("route", Node{Kind: KindBranch, Default: []string{"in"}}, http.StatusBadRequest)
	e.printer("laser")
	e.node("forms", Node{Kind: KindTemplate, To: []string{"laser"}}, http.StatusOK)
	e.node("docs", Node{Kind: KindEndpoint, To: []string{"route"}}, http.StatusOK)
	if out := e.node("route", Node{Kind: KindBranch, Cases: []Case{{When: []Check{{Key: "x"}}, To: []string{"forms"}}}}, http.StatusBadRequest); !strings.Contains(string(out), "takes data, but bdf reaches it") {
		t.Errorf("typed: %s", out)
	}

	// the public topology shows the arms, not what they check
	e.node("route", Node{Kind: KindBranch, Cases: []Case{{Name: "a", When: []Check{{Key: "secret"}}, To: []string{"bills"}}}, Default: []string{"rest"}}, http.StatusOK)
	var topo []TopoNode
	json.Unmarshal(e.admin("GET", "/api/topology", nil, http.StatusOK), &topo)
	for _, x := range publicTopology(topo) {
		if x.Name == "route" && (len(x.Cases) != 1 || x.Cases[0].Name != "a" || x.Cases[0].To[0] != "bills" || x.Cases[0].When != nil || x.Default[0] != "rest") {
			t.Errorf("public %+v", x.Node)
		}
	}
}
