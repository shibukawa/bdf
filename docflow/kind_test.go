package docflow

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/print/ipp"
)

// The kinds are a registry: the core's and the IPP adapter's
// (adapters_test.go), entries first, as the API tells them.
func TestKindRegistry(t *testing.T) {
	e := newEnv(t, nil)
	var ks []KindInfo
	if err := json.Unmarshal(e.admin("GET", "/api/kinds", nil, http.StatusOK), &ks); err != nil {
		t.Fatal(err)
	}
	by := map[string]KindInfo{}
	var order []Role
	for _, k := range ks {
		by[k.Name] = k
		if len(order) == 0 || order[len(order)-1] != k.Role {
			order = append(order, k.Role)
		}
	}
	if !slices.Equal(order, []Role{RoleEntry, RoleRoute, RoleStep, RoleExit}) {
		t.Errorf("roles in order: %v", order)
	}
	if p := by[KindPrinter]; p.Role != RoleExit || !p.SendsOut || !p.Fixed || !slices.Equal(p.Takes, []string{PayloadBDF}) {
		t.Errorf("printer: %+v", p)
	}
	if c := by[KindCups]; c.Role != RoleEntry || c.Title != "CUPS server" || !c.Port || !c.PassesOn || c.SendsOut {
		t.Errorf("cups: %+v", c)
	}
	if x := by[KindEndpoint]; !x.Port || x.SendsOut {
		t.Errorf("endpoint: %+v", x)
	}
	if x := by[KindTemplate]; x.Makes != PayloadBDF || !slices.Equal(x.Takes, []string{PayloadData}) {
		t.Errorf("template: %+v", x)
	}
	// a kind nobody registered is refused, naming those there are
	if out := e.node("x", Node{Kind: "smtp", To: []string{"y"}}, http.StatusBadRequest); !strings.Contains(string(out), "printer") {
		t.Errorf("unknown kind: %s", out)
	}
	defer func() {
		if recover() == nil {
			t.Error("a kind registered twice does not panic")
		}
	}()
	RegisterKind(Kind{Name: KindPrinter, Role: RoleExit})
}

// A node with no_ipp is no IPP printer at /ipp/print: not there, not
// listed, no address; what other nodes send it still reaches it.
func TestNoIPP(t *testing.T) {
	e := newEnv(t, nil)
	e.display("wh-a", Node{})
	e.node("quiet", Node{Kind: KindQueue, To: []string{"wh-a"}, NoIPP: true}, http.StatusOK)
	e.node("front", Node{Kind: KindQueue, To: []string{"quiet"}}, http.StatusOK)
	ctx := context.Background()

	if _, err := e.ippClient("quiet").PrintJob(ctx, "hidden", "application/x-bdf", nil, bytes.NewReader(testDoc(t, 1))); err == nil || !strings.Contains(err.Error(), "404") {
		t.Fatalf("print to a node with no_ipp: %v", err)
	}
	ps, err := (&ipp.Client{URL: e.srv.URL + IPPPath, User: "wms"}).CupsGetPrinters(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, p := range ps {
		names = append(names, p.Name)
	}
	if slices.Contains(names, "quiet") || !slices.Contains(names, "front") || !slices.Contains(names, "wh-a") {
		t.Errorf("CUPS-Get-Printers: %v", names)
	}
	var topo []TopoNode
	json.Unmarshal(e.admin("GET", "/api/topology", nil, http.StatusOK), &topo)
	for _, n := range topo {
		want := IPPPath + "/" + n.Name
		if n.Name == "quiet" {
			want = ""
		}
		if n.Addresses["ipp"] != want {
			t.Errorf("%s: addresses %v", n.Name, n.Addresses)
		}
	}

	// through the node before it, the job goes on through it
	js := e.printJob("front", "through", 1)
	j := e.job(js.ID)
	if !slices.ContainsFunc(j.Legs, func(l Leg) bool { return l.Node == "quiet" && l.State == LegRouted }) {
		t.Errorf("legs: %+v", j.Legs)
	}

	// IPP on again
	e.node("quiet", Node{Kind: KindQueue, To: []string{"wh-a"}}, http.StatusOK)
	e.printJob("quiet", "shown", 1)
}
