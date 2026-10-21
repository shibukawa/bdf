package manager

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

// A filter sends a job on by the first rule that fits it, or otherwise,
// or ends it.
func TestFilter(t *testing.T) {
	e := newEnv(t, func(o *Options) { o.SubmitTokens = map[string]string{"wms": "s3cret", "erp": "pw"} })
	for _, d := range []string{"wh-a", "wh-b", "office"} {
		e.display(d, Node{})
	}
	e.node("sort", Node{Kind: KindFilter}, http.StatusBadRequest)
	e.node("sort", Node{Kind: KindFilter, Rules: []Rule{{Match: Match{Name: "("}, To: []string{"wh-a"}}}}, http.StatusBadRequest)
	e.node("sort", Node{Kind: KindFilter, Rules: []Rule{{Match: Match{}, To: []string{"wh-a"}}}}, http.StatusBadRequest)
	e.node("sort", Node{Kind: KindFilter, Rules: []Rule{{Match: Match{Name: "x"}}}}, http.StatusBadRequest)
	e.node("sort", Node{Kind: KindFilter, Rules: []Rule{{Match: Match{MinPages: 3, MaxPages: 2}, To: []string{"wh-a"}}}}, http.StatusBadRequest)
	e.node("sort", Node{Kind: KindFilter, Rules: []Rule{{Match: Match{Name: "x"}, To: []string{"nowhere"}}}}, http.StatusBadRequest)
	e.node("sort", Node{Kind: KindFilter, Title: "Sorter", Rules: []Rule{
		{Match: Match{Name: "(?i)^urgent"}, To: []string{"wh-a", "wh-b"}},
		{Match: Match{MinPages: 3}, To: []string{"wh-b"}},
		{Match: Match{User: "^erp$", Format: "bdf"}, To: []string{"office"}},
	}, Otherwise: []string{"office"}}, http.StatusOK)
	e.node("in", Node{Kind: KindQueue, To: []string{"sort"}}, http.StatusOK)
	// a cycle through a rule
	e.node("loop", Node{Kind: KindFilter, Otherwise: []string{"in"}}, http.StatusOK)
	if out := e.node("in", Node{Kind: KindQueue, To: []string{"sort", "loop"}}, http.StatusBadRequest); !strings.Contains(string(out), "cycle") {
		t.Fatalf("cycle: %s", out)
	}
	e.admin("DELETE", "/api/nodes/loop", nil, http.StatusNoContent)

	send := func(who, name string, pages int) *Job {
		t.Helper()
		res, out := e.do(nil, "POST", "/api/destinations/in/jobs?name="+name, who, testDoc(t, pages))
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("%s: %d %s", name, res.StatusCode, out)
		}
		var j Job
		json.Unmarshal(out, &j)
		return &j
	}
	for _, c := range []struct {
		who, name string
		pages     int
		want      string
	}{
		{"wms:s3cret", "Urgent-pick", 1, "in:routed sort:routed(rule 1) wh-a:shown wh-b:shown"},
		{"wms:s3cret", "long", 4, "in:routed sort:routed(rule 2) wh-b:shown"},
		{"erp:pw", "invoice", 1, "in:routed sort:routed(rule 3) office:shown"},
		{"wms:s3cret", "plain", 1, "in:routed sort:routed(otherwise) office:shown"},
	} {
		if j := send(c.who, c.name, c.pages); legs(j) != c.want {
			t.Errorf("%s: %s; want %s", c.name, legs(j), c.want)
		}
	}
	// no otherwise: a job no rule fits ends at the filter
	e.node("sort", Node{Kind: KindFilter, Rules: []Rule{{Match: Match{Name: "^never"}, To: []string{"wh-a"}}}}, http.StatusOK)
	j := send("wms:s3cret", "dropped", 1)
	if legs(j) != "in:routed sort:dropped" || j.State != JobCompleted || !strings.Contains(j.Legs[1].Error, "no rule fits") {
		t.Fatalf("dropped %+v: %s", j, legs(j))
	}
	// the topology shows the rules; the filter is sent to, so it stays
	var topo []TopoNode
	json.Unmarshal(e.admin("GET", "/api/topology", nil, http.StatusOK), &topo)
	for _, n := range topo {
		if n.Name == "sort" && (len(n.Rules) != 1 || n.Rules[0].Match.Name != "^never") {
			t.Fatalf("topology %+v", n)
		}
	}
	e.admin("DELETE", "/api/nodes/wh-a", nil, http.StatusConflict)
	e.admin("DELETE", "/api/nodes/sort", nil, http.StatusConflict)
}
