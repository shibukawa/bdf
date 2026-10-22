package docflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
)

// A sink keeps the latest documents, lets the oldest go, and shows their
// pages to those who see the jobs.
func TestSink(t *testing.T) {
	e := newEnv(t, func(o *Options) { o.SubmitTokens = map[string]string{"wms": "s3cret", "erp": "pw"} })
	e.node("latest", Node{Kind: KindSink, Latest: 2000}, http.StatusBadRequest)
	e.node("latest", Node{Kind: KindSink, Title: "Latest", Latest: 2}, http.StatusOK)
	e.node("in", Node{Kind: KindQueue, To: []string{"latest"}}, http.StatusOK)
	doc := testDoc(t, 2)
	for i, who := range []string{"wms:s3cret", "erp:pw", "wms:s3cret"} {
		if res, out := e.do(nil, "POST", fmt.Sprintf("/api/destinations/in/jobs?name=doc-%d", i+1), who, doc); res.StatusCode != http.StatusCreated {
			t.Fatalf("job %d: %d %s", i+1, res.StatusCode, out)
		}
	}
	j := e.job(3)
	if legs(j) != "in:routed latest:kept" || j.State != JobCompleted || j.Live != 1 {
		t.Fatalf("job %+v: %s", j, legs(j))
	}
	// the sink keeps the last two; the first was let go
	var kept []Job
	json.Unmarshal(e.admin("GET", "/api/nodes/latest/latest", nil, http.StatusOK), &kept)
	if len(kept) != 2 || kept[0].Number != 3 || kept[1].Number != 2 {
		t.Fatalf("kept %+v", kept)
	}
	if j := e.job(1); j.Live != 0 {
		t.Fatalf("not let go: %+v", j)
	}
	// a user sees their own there
	_, out := e.do(nil, "GET", "/api/nodes/latest/latest", "erp:pw", nil)
	json.Unmarshal(out, &kept)
	if len(kept) != 1 || kept[0].Number != 2 {
		t.Fatalf("erp's kept %+v", kept)
	}
	e.admin("GET", "/api/nodes/in/latest", nil, http.StatusNotFound)
	// the pages of a job, for those who see it
	res, out := e.do(nil, "GET", "/api/jobs/3/pages/2.png?w=300", "wms:s3cret", nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("page: %d %s", res.StatusCode, out)
	}
	if img, err := png.Decode(bytes.NewReader(out)); err != nil || img.Bounds().Dx() != 300 {
		t.Fatalf("page image %v", err)
	}
	if res, _ := e.do(nil, "GET", "/api/jobs/3/pages/3.png", "wms:s3cret", nil); res.StatusCode != http.StatusNotFound {
		t.Fatalf("page 3 of 2: %d", res.StatusCode)
	}
	if res, _ := e.do(nil, "GET", "/api/jobs/3/pages/1.png", "erp:pw", nil); res.StatusCode != http.StatusNotFound {
		t.Fatalf("another's page: %d", res.StatusCode)
	}
	// the topology and the activity tell what is kept
	var topo []TopoNode
	json.Unmarshal(e.admin("GET", "/api/topology", nil, http.StatusOK), &topo)
	if topo[1].Name != "latest" || topo[1].Kept != 2 || topo[1].Latest != jobID(3) {
		t.Fatalf("topology %+v", topo[1])
	}
	var l Live
	json.Unmarshal(e.admin("GET", "/api/topology/live", nil, http.StatusOK), &l)
	if l.Nodes["latest"].Kept != 2 || l.Nodes["latest"].Latest != jobID(3) {
		t.Fatalf("live %+v", l.Nodes)
	}
	// after the keep, the jobs the sink keeps stay; the one let go goes
	e.clock.add(15 * 24 * time.Hour)
	if rep := e.tick(); rep.Jobs != 1 {
		t.Fatalf("tick %+v", rep)
	}
	e.job(2)
	if res, _ := e.do(nil, "GET", "/api/jobs/1", adminToken, nil); res.StatusCode != http.StatusNotFound {
		t.Fatalf("the job let go stayed: %d", res.StatusCode)
	}
	// deleting the sink lets its jobs go
	e.admin("DELETE", "/api/nodes/latest", nil, http.StatusConflict) // in sends to it
	e.admin("DELETE", "/api/nodes/in", nil, http.StatusNoContent)
	e.admin("DELETE", "/api/nodes/latest", nil, http.StatusNoContent)
	if j := e.job(3); j.Live != 0 {
		t.Fatalf("kept after the sink went: %+v", j)
	}
}

// An entry node's own port: a CUPS server's IPP, an endpoint's REST. A
// CUPS server forwards what it receives along its edges, as an endpoint
// does; its port is a CUPS server with one printer, the node.
func TestPortHandler(t *testing.T) {
	e := newEnv(t, nil)
	ctx := context.Background()
	e.display("wh-a", Node{Title: "Warehouse A"})
	e.printer("laser")
	e.node("cups1", Node{Kind: KindCups, Port: 6310}, http.StatusBadRequest)
	e.node("cups1", Node{Kind: KindCups, To: []string{"wh-a"}}, http.StatusBadRequest)
	e.node("cups1", Node{Kind: KindCups, To: []string{"wh-a"}, Port: 80}, http.StatusBadRequest)
	e.node("cups1", Node{Kind: KindCups, To: []string{"wh-a"}, Port: 6310, TLSCert: "cert.pem"}, http.StatusBadRequest)
	e.node("cups1", Node{Kind: KindCups, To: []string{"nowhere"}, Port: 6310}, http.StatusBadRequest)
	e.node("cups1", Node{Kind: KindCups, To: []string{"cups1"}, Port: 6310}, http.StatusBadRequest)
	e.node("cups1", Node{Kind: KindCups, Title: "CUPS", To: []string{"wh-a", "laser"}, Port: 6310}, http.StatusOK)
	e.node("api", Node{Kind: KindEndpoint, To: []string{"wh-a"}, Port: 6310}, http.StatusConflict)
	var made struct {
		Key string `json:"key"`
	}
	json.Unmarshal(e.node("api", Node{Kind: KindEndpoint, To: []string{"wh-a"}, Port: 6311}, http.StatusOK), &made)
	// its edges are edges: a node it sends to stays, and no path leads back
	e.admin("DELETE", "/api/nodes/laser", nil, http.StatusConflict)
	if out := e.node("relay", Node{Kind: KindQueue, To: []string{"cups1"}}, http.StatusOK); len(out) == 0 {
		t.Fatal("a queue may send to a CUPS server")
	}
	if out := e.node("cups1", Node{Kind: KindCups, To: []string{"relay"}, Port: 6310}, http.StatusBadRequest); !strings.Contains(string(out), "cycle") {
		t.Fatalf("cycle: %s", out)
	}
	e.admin("DELETE", "/api/nodes/relay", nil, http.StatusNoContent)

	// the CUPS server's port: the server at /, with one printer, the node
	// at /printers/<node>; what is printed there follows the edges
	cups := httptest.NewServer(PortHandler(e.m(), "cups1"))
	defer cups.Close()
	ps, err := (&ipp.Client{URL: cups.URL + "/", User: "lp"}).CupsGetPrinters(ctx)
	if err != nil || len(ps) != 1 || ps[0].Name != "cups1" || ps[0].Info != "CUPS" || ps[0].URI != strings.Replace(cups.URL, "http://", "ipp://", 1)+"/printers/cups1" {
		t.Fatalf("CUPS-Get-Printers on the port: %+v %v", ps, err)
	}
	pa, err := (&ipp.Client{URL: cups.URL + "/printers/cups1", User: "lp"}).GetPrinterAttributes(ctx)
	if err != nil || pa.Name != "cups1" || pa.URI != strings.Replace(cups.URL, "http://", "ipp://", 1)+"/printers/cups1" {
		t.Fatalf("the printer's attributes on the port: %+v %v", pa, err)
	}
	// the delivery to the printer waits for the tick: the job is processing
	js, err := (&ipp.Client{URL: cups.URL + "/printers/cups1", User: "lp"}).PrintJob(ctx, "from cups", "application/x-bdf", nil, bytes.NewReader(testDoc(t, 1)))
	if err != nil || js.State != ipp.JobProcessing {
		t.Fatalf("print through the port: %+v %v", js, err)
	}
	if j := e.job(1); j.Dest != "cups1" || j.Source != "ipp" || legs(j) != "cups1:routed wh-a:shown laser:pending" {
		t.Fatalf("job %+v", j)
	}
	// a plain IPP client prints to /ipp/print of the port: the node too
	if js, err := (&ipp.Client{URL: cups.URL + "/ipp/print", User: "lp"}).PrintJob(ctx, "plain ipp", "application/x-bdf", nil, bytes.NewReader(testDoc(t, 1))); err != nil || js.State != ipp.JobProcessing {
		t.Fatalf("print to /ipp/print of the port: %+v %v", js, err)
	}
	if j := e.job(2); j.Dest != "cups1" || legs(j) != "cups1:routed wh-a:shown laser:pending" {
		t.Fatalf("job %+v", j)
	}
	// another node's name is no printer on the port: it is reached along the edges
	if _, err := (&ipp.Client{URL: cups.URL + "/printers/wh-a", User: "lp"}).PrintJob(ctx, "to a node", "application/x-bdf", nil, bytes.NewReader(testDoc(t, 1))); err == nil {
		t.Fatal("printed to a node behind the CUPS server")
	}
	if res, _ := http.Post(cups.URL+"/nothing/here", "application/ipp", nil); res.StatusCode != http.StatusNotFound {
		t.Fatalf("a path that is no printer: %d", res.StatusCode)
	}
	if res, err := http.Get(cups.URL + "/"); err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("GET on the port: %v", err)
	}
	// on the manager's own server the node is a destination like any
	if js := e.printJob("cups1", "own server", 1); js.State != ipp.JobProcessing {
		t.Fatalf("print to the node on the manager's server: %+v", js)
	}

	if j := e.job(3); j.Dest != "cups1" || legs(j) != "cups1:routed wh-a:shown laser:pending" {
		t.Fatalf("job %+v", j)
	}
	// the endpoint's port: every POST is to it, with its key
	api := httptest.NewServer(PortHandler(e.m(), "api"))
	defer api.Close()
	post := func(key string) int {
		req, _ := http.NewRequest("POST", api.URL+"/anything?name=x", bytes.NewReader(testDoc(t, 1)))
		if key != "" {
			req.Header.Set("Authorization", "Bearer "+key)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		return res.StatusCode
	}
	if post(made.Key) != http.StatusCreated || post("") != http.StatusUnauthorized {
		t.Fatal("the endpoint's port")
	}
	if j := e.job(4); j.Dest != "api" || j.Source != "endpoint" {

		t.Fatalf("job %+v", j)
	}
	// the node gone, or without its port: 404
	e.node("api", Node{Kind: KindEndpoint, To: []string{"wh-a"}}, http.StatusOK)
	if post(made.Key) != http.StatusNotFound {
		t.Fatal("a port the node lost")
	}
	// where there are no ports (serverless), none is given
	s := newEnv(t, func(o *Options) { o.NoPorts = true })
	s.display("wh-a", Node{})
	if out := s.node("cups1", Node{Kind: KindCups, To: []string{"wh-a"}, Port: 6310}, http.StatusBadRequest); !strings.Contains(string(out), "serverless") {
		t.Fatalf("no ports: %s", out)
	}
	// cupsDest: the server, and the node under CUPS's and IPP's paths
	for path, want := range map[string]string{"/": "", "/admin": "", "/ipp/print": "cups1", "/printers/cups1": "cups1", "/classes/cups1": "cups1", "/ipp/print/cups1/": "cups1"} {
		if got, ok := cupsDest(path, "cups1"); !ok || got != want {
			t.Errorf("%s: %q %v", path, got, ok)
		}
	}
	for _, path := range []string{"/jobs/1", "/printers/wh-a", "/classes/tokyo", "/ipp/print/x", "/printers/cups1x", "/printers/"} {
		if _, ok := cupsDest(path, "cups1"); ok {
			t.Errorf("%s is a destination", path)
		}
	}
}

// A printer's preview shows what its fake printer received, from the
// farm's API.
func TestPreview(t *testing.T) {
	e := newEnv(t, nil)
	var pic bytes.Buffer
	png.Encode(&pic, testImage())
	farm := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/jobs" && r.URL.Query().Get("printer") == "eve":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `[{"printer":"eve","id":7,"name":"invoice","state":"completed","created":"2026-10-21T00:00:00Z"},{"printer":"eve","id":6,"name":"older","state":"completed","created":"2026-10-20T00:00:00Z"}]`)
		case r.URL.Path == "/api/printers/eve/jobs/7/pages":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"format":"pdf","pages":[{"Width":100,"Height":100,"DPI":72},{"Width":100,"Height":100,"DPI":72}]}`)
		case r.URL.Path == "/api/printers/eve/jobs/7/pages/1":
			w.Header().Set("Content-Type", "image/png")
			w.Write(pic.Bytes())
		default:
			http.NotFound(w, r)
		}
	}))
	defer farm.Close()
	url := "ipp://127.0.0.1:1/ipp/print"
	e.admin("PUT", "/api/nodes/laser", Node{Kind: KindPrinter, URL: url, Preview: "ftp://x"}, http.StatusBadRequest)
	e.admin("PUT", "/api/nodes/laser", Node{Kind: KindPrinter, URL: url, Preview: farm.URL + "/printers/eve"}, http.StatusBadRequest)
	e.admin("PUT", "/api/nodes/laser", Node{Kind: KindPrinter, URL: url}, http.StatusOK)
	e.admin("GET", "/api/nodes/laser/preview", nil, http.StatusNotFound)
	e.admin("PUT", "/api/nodes/laser", Node{Kind: KindPrinter, URL: url, Preview: farm.URL + "/api/printers/eve"}, http.StatusOK)
	var jobs []PreviewJob
	json.Unmarshal(e.admin("GET", "/api/nodes/laser/preview", nil, http.StatusOK), &jobs)
	if len(jobs) != 2 || jobs[0].ID != 7 || jobs[0].Pages != 2 || jobs[0].Name != "invoice" || jobs[1].Pages != 0 {
		t.Fatalf("preview %+v", jobs)
	}
	res, out := e.do(nil, "GET", "/api/nodes/laser/preview/7/pages/1.png", adminToken, nil)
	if res.StatusCode != http.StatusOK || res.Header.Get("Content-Type") != "image/png" || !bytes.Equal(out, pic.Bytes()) {
		t.Fatalf("page: %d %s", res.StatusCode, res.Header.Get("Content-Type"))
	}
	if res, _ := e.do(nil, "GET", "/api/nodes/laser/preview/7/pages/2.png", adminToken, nil); res.StatusCode != http.StatusNotFound {
		t.Fatalf("a page the farm has not: %d", res.StatusCode)
	}
	e.display("wh-a", Node{})
	e.admin("GET", "/api/nodes/wh-a/preview", nil, http.StatusNotFound)
}

// testImage is a small image for a farm's page.
func testImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for i := range img.Pix {
		img.Pix[i] = 0x80
	}
	return img
}
