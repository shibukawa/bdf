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
	if c := by[KindCups]; c.Role != RoleEntry || c.Title != "CUPS server (source)" || !c.Port || !c.PassesOn || c.SendsOut {
		t.Errorf("cups: %+v", c)
	}
	if c := by[KindCupsOut]; c.Role != RoleExit || c.Title != "CUPS server (sink)" || !c.SendsOut || !c.Fixed || !c.Discovers || !slices.Equal(c.Takes, []string{PayloadBDF}) {
		t.Errorf("cups-out: %+v", c)
	}
	if x := by[KindUpload]; x.Port || x.SendsOut || !x.Keyed || x.Title != "Uploader" {
		t.Errorf("upload: %+v", x)
	}
	if x := by[KindDownload]; x.Title != "Downloader" || !slices.Equal(x.Takes, []string{PayloadStocked}) {
		t.Errorf("download: %+v", x)
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

// IPP enters the topology at a CUPS server alone: the other nodes are no
// IPP printers at /ipp/print, not listed there, and have no ipp address;
// what the CUPS server receives goes on along its edges.
func TestIPPAtCups(t *testing.T) {
	e := newEnv(t, nil)
	e.display("wh-a", Node{})
	e.printer("laser")
	e.node("front", Node{Kind: KindQueue, To: []string{"wh-a"}}, http.StatusOK)
	e.node("tokyo", Node{Kind: KindClass, Members: []string{"laser"}}, http.StatusOK)
	c := e.cups("cups1", "front")
	ctx := context.Background()

	for _, name := range []string{"front", "wh-a", "laser", "tokyo"} {
		if _, err := e.ippClient(name).PrintJob(ctx, "hidden", "application/x-bdf", nil, bytes.NewReader(testDoc(t, 1))); err == nil || !strings.Contains(err.Error(), "404") {
			t.Errorf("print to %s over IPP: %v", name, err)
		}
	}
	ps, err := (&ipp.Client{URL: e.srv.URL + IPPPath, User: "wms"}).CupsGetPrinters(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, p := range ps {
		names = append(names, p.Name)
	}
	if !slices.Equal(names, []string{"cups1"}) {
		t.Errorf("CUPS-Get-Printers: %v", names)
	}
	var topo []TopoNode
	json.Unmarshal(e.admin("GET", "/api/topology", nil, http.StatusOK), &topo)
	for _, n := range topo {
		want := ""
		if n.Name == "cups1" {
			want = IPPPath + "/cups1"
		}
		if n.Addresses["ipp"] != want {
			t.Errorf("%s: addresses %v", n.Name, n.Addresses)
		}
	}

	// through the CUPS server, the job goes on through the queue
	js, err := c.PrintJob(ctx, "through", "application/x-bdf", nil, bytes.NewReader(testDoc(t, 1)))
	if err != nil || js.State != ipp.JobCompleted {
		t.Fatalf("print to cups1: %+v %v", js, err)
	}
	if j := e.job(js.ID); j.Dest != "cups1" || j.Source != "ipp" || legs(j) != "cups1:routed front:routed wh-a:shown" {
		t.Errorf("job %+v: %s", j, legs(j))
	}
	// what the API sends the queue goes there as before
	if j := e.printJob("front", "by the API", 1); legs(j) != "front:routed wh-a:shown" {
		t.Errorf("job %+v: %s", j, legs(j))
	}
}
