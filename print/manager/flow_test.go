package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	_ "github.com/shibukawa/bdf/converter/csv"
	_ "github.com/shibukawa/bdf/converter/image"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/printconv"
	"github.com/shibukawa/yuruidb"
)

// testPNG is a small picture: a file the converter makes a page of.
func testPNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for i := range img.Pix {
		img.Pix[i] = 0x80
	}
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// The graph is typed: what an entry sends on (a document made .bdf, data,
// a file as received) must be taken by every node it reaches; a convert
// node makes a .bdf of a file on the way; an archive takes anything.
func TestKinds(t *testing.T) {
	e := newEnv(t, func(o *Options) { o.Converter = printconv.Converter{NoSystemFonts: true} })
	ctx := context.Background()
	laser := e.printer("laser")
	store := "bbolt:" + filepath.ToSlash(filepath.Join(t.TempDir(), "audit.db"))
	e.node("audit", Node{Kind: KindArchive, Store: store}, http.StatusOK)
	// a file as received does not go to a printer, nor does data; a
	// convert node in between takes anything and sends a .bdf on
	if out := e.node("raw", Node{Kind: KindEndpoint, Output: OutputFile, To: []string{"laser"}}, http.StatusBadRequest); !strings.Contains(string(out), "takes bdf") || !strings.Contains(string(out), "any reaches it: raw → laser") {
		t.Fatalf("typed: %s", out)
	}
	e.node("raw", Node{Kind: KindEndpoint, Output: "maybe", To: []string{"audit"}}, http.StatusBadRequest)
	e.node("conv", Node{Kind: KindConvert}, http.StatusBadRequest)
	e.node("conv", Node{Kind: KindConvert, To: []string{"laser"}, Formats: []string{"PDF/A"}}, http.StatusBadRequest)
	e.node("conv", Node{Kind: KindConvert, To: []string{"laser"}}, http.StatusOK)
	e.node("fan", Node{Kind: KindQueue, To: []string{"audit", "conv"}}, http.StatusOK)
	var raw struct {
		Key string `json:"key"`
	}
	json.Unmarshal(e.node("raw", Node{Kind: KindEndpoint, Output: OutputFile, To: []string{"fan"}}, http.StatusOK), &raw)
	// a change deeper in the graph is checked against every entry: the
	// queue cannot send to the printer while a file as received reaches it
	if out := e.node("fan", Node{Kind: KindQueue, To: []string{"audit", "conv", "laser"}}, http.StatusBadRequest); !strings.Contains(string(out), "raw → fan → laser") {
		t.Fatalf("typed through a queue: %s", out)
	}
	e.node("din", Node{Kind: KindEndpoint, Accept: AcceptData, To: []string{"laser"}}, http.StatusBadRequest)
	e.node("forms", Node{Kind: KindTemplate, Template: "x", To: []string{"laser"}}, http.StatusOK)
	if out := e.node("docs", Node{Kind: KindEndpoint, To: []string{"forms"}}, http.StatusBadRequest); !strings.Contains(string(out), "takes data, but bdf reaches it") {
		t.Fatalf("typed template: %s", out)
	}
	var topo []TopoNode
	json.Unmarshal(e.admin("GET", "/api/topology", nil, http.StatusOK), &topo)
	kinds := map[string]string{}
	for _, n := range topo {
		kinds[n.Name] = n.Takes + "<" + strings.Join(n.Emits, "+")
	}
	// an archive sends on what it wrote (an export after it), as a queue does
	want := map[string]string{"raw": "any<any", "fan": "any<any", "conv": "any<bdf", "laser": "bdf<", "audit": "any<any", "forms": "data<bdf"}
	for name, w := range want {
		if kinds[name] != w {
			t.Errorf("%s: %s; want %s", name, kinds[name], w)
		}
	}

	// what a node may take now is what every node after it takes: a queue
	// in front of the printer takes .bdf alone, one in front of the
	// printer and the template node nothing (its first edge types it, as
	// a wildcard pin's does)
	e.node("q", Node{Kind: KindQueue, To: []string{"laser"}}, http.StatusOK)
	accepts := func() map[string]string {
		var topo []TopoNode
		json.Unmarshal(e.admin("GET", "/api/topology", nil, http.StatusOK), &topo)
		out := map[string]string{}
		for _, n := range topo {
			out[n.Name] = strings.Join(n.Accepts, "+")
			for to, bad := range n.Broken {
				out[n.Name+">"+to] = strings.Join(bad, "+")
			}
		}
		return out
	}
	if got := accepts(); got["q"] != "bdf" || got["fan"] != "bdf+data+any" || got["conv"] != "bdf+data+any" || got["forms"] != "data" || got["laser"] != "bdf" || got["raw"] != "bdf+data+any" {
		t.Fatalf("accepts %v", got)
	}
	e.node("q", Node{Kind: KindQueue, To: []string{"laser", "forms"}}, http.StatusOK)
	if got := accepts(); got["q"] != "" || !strings.Contains(string(e.admin("GET", "/api/topology", nil, http.StatusOK)), `"accepts":[]`) {
		t.Fatalf("a dead end accepts nothing: %v", got)
	}
	// a change of settings that leaves an edge carrying what the node
	// after it does not take is not refused: the edge is shown broken,
	// and a job gets a leg of the error there
	var din struct {
		Key string `json:"key"`
	}
	json.Unmarshal(e.node("din", Node{Kind: KindEndpoint, Accept: AcceptData, To: []string{"forms"}}, http.StatusOK), &din)
	e.node("din", Node{Kind: KindEndpoint, To: []string{"forms"}}, http.StatusOK)
	if got := accepts(); got["din>forms"] != "bdf" || got["din"] != "data" || len(got) != 8+1 {
		t.Fatalf("broken %v", got)
	}
	bres, bout := e.do(nil, "POST", "/in/din?name=doc", din.Key, testDoc(t, 1))
	var bj Job
	json.Unmarshal(bout, &bj)
	if bres.StatusCode != http.StatusCreated || legs(&bj) != "din:routed forms:error" || !strings.Contains(bj.Legs[1].Error, "takes data") {
		t.Fatalf("along a broken edge: %d %s %+v", bres.StatusCode, legs(&bj), bj.Legs)
	}
	// a new edge along which a kind reaches a node that does not take it
	// is refused still, broken edges or not
	e.node("q", Node{Kind: KindQueue, To: []string{"forms"}}, http.StatusOK)
	if out := e.node("din", Node{Kind: KindEndpoint, To: []string{"forms", "q"}}, http.StatusBadRequest); !strings.Contains(string(out), "forms (template) takes data, but bdf reaches it: din → q → forms") {
		t.Fatalf("a new edge: %s", out)
	}
	// and the same edges put again are not refused for the broken one
	e.node("din", Node{Kind: KindEndpoint, To: []string{"forms"}, Title: "documents"}, http.StatusOK)
	// mended: the entry takes data again
	e.node("din", Node{Kind: KindEndpoint, Accept: AcceptData, To: []string{"forms"}}, http.StatusOK)
	if got := accepts(); got["din>forms"] != "" || len(got) != 8 {
		t.Fatalf("mended %v", got)
	}
	e.admin("DELETE", "/api/nodes/din", nil, http.StatusNoContent)
	e.admin("DELETE", "/api/nodes/q", nil, http.StatusNoContent)

	// a picture received as it is (any file): kept at the archive as it
	// came, made a document by the convert node and printed
	pic := testPNG(t)
	res, out := e.do(nil, "POST", "/in/raw?name=photo.png", raw.Key, pic)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("png: %d %s", res.StatusCode, out)
	}
	var j Job
	json.Unmarshal(out, &j)
	if j.Kind != PayloadAny || j.Format != "image" || j.Size != int64(len(pic)) || legs(&j) != "raw:routed fan:routed audit:pending conv:pending" {
		t.Fatalf("png job %+v: %s", j, legs(&j))
	}
	if rep := e.tick(); rep.Delivered != 2 {
		t.Fatalf("tick %+v: %+v", rep, e.job(j.Number).Legs)
	}
	if rep := e.tick(); rep.Delivered != 1 { // the printer's delivery, made by the convert node
		t.Fatalf("tick %+v", rep)
	}
	got := e.job(j.Number)
	if got.State != JobCompleted || legs(got) != "raw:routed fan:routed audit:done conv:done laser:done" || len(got.Payloads) != 2 || len(laser.Jobs()) != 1 {
		t.Fatalf("job %+v: %s", got, legs(got))
	}
	made := got.Payloads[1]
	if made.Kind != PayloadBDF || made.Made != "conv" || made.Name != "photo.bdf" || made.Pages != 1 || got.Legs[3].Made != made.ID || got.Legs[4].Payload != made.ID || got.Legs[2].Payload != got.ID {
		t.Fatalf("payloads %+v legs %+v", got.Payloads, got.Legs)
	}
	st, err := e.st.archive(ctx, store, "archive-audit")
	if err != nil {
		t.Fatal(err)
	}
	page, _ := st.List(ctx, yuruidb.ListOptions{Limit: 10})
	key := archiveKey(got, got.received())
	if len(page.Blobs) != 2 || page.Blobs[0].ID != key+".json" || page.Blobs[1].ID != key+".png" {
		t.Fatalf("archived %+v", page.Blobs)
	}
	r, _ := st.Open(ctx, key+".json")
	var rec ArchiveRecord
	json.NewDecoder(r).Decode(&rec)
	r.Close()
	if rec.Kind != PayloadAny || rec.Format != "image" || rec.Object != key+".png" || rec.Size != int64(len(pic)) || rec.Meta["$file"] != "photo.png" {
		t.Fatalf("record %+v", rec)
	}
	// the document to show is the one made; the file is downloaded as it came
	if res, out := e.do(nil, "GET", fmt.Sprintf("/api/jobs/%d/document", j.Number), adminToken, nil); res.StatusCode != http.StatusOK || !bytes.HasPrefix(out, bdf.Magic[:]) || !strings.Contains(res.Header.Get("Content-Disposition"), "photo.bdf") {
		t.Fatalf("document: %d %s", res.StatusCode, res.Header.Get("Content-Disposition"))
	}
	if res, out := e.do(nil, "GET", fmt.Sprintf("/api/jobs/%d/document?payload=%s", j.Number, j.ID), adminToken, nil); res.StatusCode != http.StatusOK || !bytes.Equal(out, pic) || res.Header.Get("Content-Type") != "image/png" {
		t.Fatalf("png: %d %s", res.StatusCode, res.Header.Get("Content-Type"))
	}
	if res, _ := e.do(nil, "GET", fmt.Sprintf("/api/jobs/%d/pages/1.png", j.Number), adminToken, nil); res.StatusCode != http.StatusOK {
		t.Fatalf("a page of the document made: %d", res.StatusCode)
	}
	if res, _ := e.do(nil, "GET", fmt.Sprintf("/api/jobs/%d/document?payload=nowhere", j.Number), adminToken, nil); res.StatusCode != http.StatusNotFound {
		t.Fatalf("no payload: %d", res.StatusCode)
	}
	// a CSV file is data: the archive keeps it as it came; the convert
	// node has no pages of it yet (a sheet has no print view), so that
	// step fails for good
	csv := []byte("id,name\n1,Ann\n2,Bob\n")
	res, out = e.do(nil, "POST", "/in/raw?name=people.csv", raw.Key, csv)
	json.Unmarshal(out, &j)
	if res.StatusCode != http.StatusCreated || j.Kind != PayloadData || j.Format != "csv" {
		t.Fatalf("csv: %d %+v", res.StatusCode, j)
	}
	if rep := e.tick(); rep.Delivered != 1 || rep.Undelivered != 1 {
		t.Fatalf("tick %+v", rep)
	}
	if got := e.job(j.Number); got.State != JobAborted || !strings.Contains(got.Legs[3].Error, "no pages") {
		t.Fatalf("csv through convert %+v: %s", got.Legs, legs(got))
	}
	if page, _ := st.List(ctx, yuruidb.ListOptions{Limit: 10}); len(page.Blobs) != 4 || page.Blobs[2].ID != archiveKey(e.job(j.Number), e.job(j.Number).received())+".csv" {
		t.Fatalf("archived csv %+v", page.Blobs)
	}

	// a .bdf document received as it is stays one: the convert node sends
	// it on as it is, and the printer prints it
	res, out = e.do(nil, "POST", "/in/raw?name=plain", raw.Key, testDoc(t, 2))
	json.Unmarshal(out, &j)
	if res.StatusCode != http.StatusCreated || j.Kind != PayloadBDF || j.Pages != 2 {
		t.Fatalf("bdf as received: %d %+v", res.StatusCode, j)
	}
	e.tick()
	e.tick()
	got = e.job(j.Number)
	if got.State != JobCompleted || len(got.Payloads) != 1 || got.Legs[3].Made != got.ID || len(laser.Jobs()) != 2 {
		t.Fatalf("bdf through convert %+v: %s", got, legs(got))
	}
	// data with no document of it: the convert node fails for good
	res, out = e.do(nil, "POST", "/in/raw?name=x.json", raw.Key, map[string]any{"a": 1})
	json.Unmarshal(out, &j)
	if res.StatusCode != http.StatusCreated || j.Kind != PayloadData || j.Format != "json" {
		t.Fatalf("json: %d %+v", res.StatusCode, j)
	}
	if rep := e.tick(); rep.Delivered != 1 || rep.Undelivered != 1 {
		t.Fatalf("tick %+v", rep)
	}
	if got := e.job(j.Number); got.State != JobAborted || !strings.Contains(got.Legs[3].Error, "no document of json") {
		t.Fatalf("json through convert %+v", got)
	}
	// what is not known at all is refused at the entry
	if res, _ := e.do(nil, "POST", "/in/raw?name=x.bin", raw.Key, []byte("\x00\x01\x02 nothing")); res.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("unknown: %d", res.StatusCode)
	}
	// a direct submission to the convert node takes the file as it is
	res, out = e.do(nil, "POST", "/api/destinations/conv/jobs?name=direct.csv", adminToken, csv)
	json.Unmarshal(out, &j)
	if res.StatusCode != http.StatusCreated || j.Kind != PayloadData || legs(&j) != "conv:pending" {
		t.Fatalf("direct to convert: %d %+v %s", res.StatusCode, j, legs(&j))
	}
	// a reprint sends what was received again, with its metadata
	res, out = e.do(nil, "POST", fmt.Sprintf("/api/jobs/%d/reprint", got.Number), adminToken, map[string]string{"to": "raw"})
	var again Job
	json.Unmarshal(out, &again)
	if res.StatusCode != http.StatusCreated || again.Kind != PayloadBDF || again.Pages != 2 || again.Source != "reprint" || again.Meta["$entry"] != "raw" {
		t.Fatalf("reprint: %d %+v", res.StatusCode, again)
	}
	// the nodes the deliveries of the jobs refer to stay
	e.admin("DELETE", "/api/nodes/conv", nil, http.StatusConflict)
	if _, err := e.ippClient("forms").PrintJob(ctx, "x", "application/x-bdf", nil, bytes.NewReader(testDoc(t, 1))); !ipp.IsStatus(err, ipp.StatusClientErrorDocumentFormatError) {
		t.Fatalf("a document to a template node: %v", err)
	}
}

// Metadata goes with a job: the sender's (headers, parameters, IPP
// attributes), the entry's, and the manager's; a filter routes by it.
func TestMeta(t *testing.T) {
	e := newEnv(t, func(o *Options) {
		o.Converter = printconv.Converter{NoSystemFonts: true}
		o.SubmitTokens = map[string]string{"wms": "s3cret"}
	})
	ctx := context.Background()
	e.printer("laser")
	store := "bbolt:" + filepath.ToSlash(filepath.Join(t.TempDir(), "audit.db"))
	e.node("audit", Node{Kind: KindArchive, Store: store}, http.StatusOK)
	e.node("keep", Node{Kind: KindSink}, http.StatusOK)
	e.node("conv", Node{Kind: KindConvert, To: []string{"laser"}}, http.StatusOK)
	e.node("sort", Node{Kind: KindFilter, Rules: []Rule{{Match: Match{Meta: map[string]string{"dept": "("}}, To: []string{"conv"}}}}, http.StatusBadRequest)
	e.node("sort", Node{Kind: KindFilter, Rules: []Rule{{Match: Match{Kind: "text"}, To: []string{"conv"}}}}, http.StatusBadRequest)
	e.node("sort", Node{Kind: KindFilter, Rules: []Rule{{Match: Match{Meta: map[string]string{"": "x"}}, To: []string{"conv"}}}}, http.StatusBadRequest)
	e.node("sort", Node{Kind: KindFilter, Rules: []Rule{
		{Match: Match{Meta: map[string]string{"dept": "^sales$", "$sender_kind": "^endpoint$"}}, To: []string{"conv"}},
		{Match: Match{Kind: PayloadData}, To: []string{"keep"}},
	}, Otherwise: []string{"audit"}}, http.StatusOK)
	e.node("tagged", Node{Kind: KindEndpoint, Output: OutputFile, To: []string{"sort"}, Meta: map[string]any{"$site": "tokyo"}}, http.StatusBadRequest)
	var tagged struct {
		Key string `json:"key"`
	}
	json.Unmarshal(e.node("tagged", Node{Kind: KindEndpoint, Output: OutputFile, To: []string{"sort"}, Meta: map[string]any{"site": "tokyo", "floor": 3}}, http.StatusOK), &tagged)
	csv := []byte("id,name\n1,Ann\n")
	send := func(name string, body any, edit func(*http.Request)) *Job {
		t.Helper()
		res, out := e.doWith(func(r *http.Request) {
			r.Header.Set("Authorization", "Bearer "+tagged.Key)
			if edit != nil {
				edit(r)
			}
		}, "POST", "/in/tagged?name="+name, body)
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("%s: %d %s", name, res.StatusCode, out)
		}
		var j Job
		json.Unmarshal(out, &j)
		return &j
	}
	// a header: the sender's key; the entry's keys over the sender's
	j := send("a.csv", csv, func(r *http.Request) {
		r.Header.Set("X-Meta-Dept", "sales")
		r.Header.Set("X-Meta-Site", "osaka")
	})
	if legs(j) != "tagged:routed sort:routed(rule 1) conv:pending" || j.Meta["dept"] != "sales" || j.Meta["site"] != "tokyo" || j.Meta["floor"] != float64(3) ||
		j.Meta["$entry"] != "tagged" || j.Meta["$file"] != "a.csv" || j.Meta["$job"] != float64(j.Number) || j.Meta["$received"] == "" {
		t.Fatalf("by header %+v: %s", j.Meta, legs(j))
	}
	// a parameter; the kind of the payload
	if j := send("b.csv&meta.dept=ops", csv, nil); legs(j) != "tagged:routed sort:routed(rule 2) keep:kept" || j.Meta["dept"] != "ops" {
		t.Fatalf("by parameter %+v: %s", j.Meta, legs(j))
	}
	if j := send("c", testDoc(t, 1), nil); legs(j) != "tagged:routed sort:routed(otherwise) audit:pending" {
		t.Fatalf("otherwise: %s", legs(j))
	}
	// too much, and the manager's keys, are refused
	many := map[string]any{}
	for i := range 65 {
		many[fmt.Sprint("k", i)] = i
	}
	if res, _ := e.do(nil, "POST", "/in/tagged?name=many.json", tagged.Key, many); res.StatusCode != http.StatusCreated {
		t.Fatalf("a JSON body is data, not metadata: %d", res.StatusCode)
	}
	if res, _ := e.doWith(func(r *http.Request) {
		r.Header.Set("Authorization", "Bearer "+tagged.Key)
		r.Header.Set("X-Meta-$job", "9")
	}, "POST", "/in/tagged?name=d", testDoc(t, 1)); res.StatusCode != http.StatusBadRequest {
		t.Fatalf("a $ key: %d", res.StatusCode)
	}
	if res, _ := e.do(nil, "POST", "/in/tagged?name=e&meta.a="+strings.Repeat("x", 17<<10), tagged.Key, testDoc(t, 1)); res.StatusCode != http.StatusBadRequest {
		t.Fatalf("too much: %d", res.StatusCode)
	}
	// IPP: the job attributes, and who printed
	c := &ipp.Client{URL: e.srv.URL + IPPPath + "/tagged", User: "wms", Username: "wms", Password: "s3cret"}
	js, err := c.PrintJob(ctx, "ipp", "application/x-bdf", []ipp.Attribute{ipp.Attr("copies", ipp.Integer(2)), ipp.Attr("media", ipp.Keyword("iso_a4_210x297mm"))}, bytes.NewReader(testDoc(t, 1)))
	if err != nil {
		t.Fatal(err)
	}
	got := e.job(js.ID)
	if got.Meta["ipp.copies"] != float64(2) || got.Meta["ipp.media"] != "iso_a4_210x297mm" || got.Meta["ipp.document-format"] != "application/x-bdf" || got.Meta["$sender"] != "wms" || got.Meta["$source"] != "ipp" || got.Meta["site"] != "tokyo" {
		t.Fatalf("ipp meta %+v", got.Meta)
	}
	// the API's destinations take metadata too
	res, out := e.do(nil, "POST", "/api/destinations/keep/jobs?name=f&meta.dept=ops", "wms:s3cret", testDoc(t, 1))
	json.Unmarshal(out, &j)
	if res.StatusCode != http.StatusCreated || j.Meta["dept"] != "ops" || j.Meta["$entry"] != "keep" || j.Meta["$sender"] != "wms" {
		t.Fatalf("api meta: %d %+v", res.StatusCode, j.Meta)
	}
	// the topology shows the rules with their metadata
	var topo []TopoNode
	json.Unmarshal(e.admin("GET", "/api/topology", nil, http.StatusOK), &topo)
	for _, n := range topo {
		if n.Name == "sort" && (n.Rules[0].Match.Meta["dept"] != "^sales$" || n.Rules[1].Match.Kind != PayloadData) {
			t.Fatalf("topology %+v", n)
		}
	}
}

// A convert node sends the document it makes along To and the file it
// took, as it was, along Original: to archives of their own, each record
// with the same trace, which names the other's objects. A filter sees the
// trace; a printer along Original, which a file does not reach, is
// refused.
func TestOriginalTrace(t *testing.T) {
	e := newEnv(t, func(o *Options) { o.Converter = printconv.Converter{NoSystemFonts: true} })
	ctx := context.Background()
	e.printer("laser")
	dir := t.TempDir()
	madeStore := "bbolt:" + filepath.ToSlash(filepath.Join(dir, "made.db"))
	origStore := "bbolt:" + filepath.ToSlash(filepath.Join(dir, "orig.db"))
	e.node("made", Node{Kind: KindArchive, Store: madeStore}, http.StatusOK)
	e.node("orig", Node{Kind: KindArchive, Store: origStore}, http.StatusOK)
	e.node("conv", Node{Kind: KindConvert, To: []string{"made"}}, http.StatusOK)
	var raw struct {
		Key string `json:"key"`
	}
	json.Unmarshal(e.node("raw", Node{Kind: KindEndpoint, Output: OutputFile, To: []string{"conv"}}, http.StatusOK), &raw)
	// the file goes on along Original as it came: not to a printer
	if out := e.node("conv", Node{Kind: KindConvert, To: []string{"made"}, Original: []string{"laser"}}, http.StatusBadRequest); !strings.Contains(string(out), "laser (printer) takes bdf, but any reaches it: raw → conv → laser") {
		t.Fatalf("typed original: %s", out)
	}
	// along Original through a filter that routes by the trace's side
	e.node("side", Node{Kind: KindFilter, Rules: []Rule{{Match: Match{Meta: map[string]string{"$trace_role": "^source$"}}, To: []string{"orig"}}}}, http.StatusOK)
	e.node("conv", Node{Kind: KindConvert, To: []string{"made"}, Original: []string{"side"}}, http.StatusOK)
	var topo []TopoNode
	json.Unmarshal(e.admin("GET", "/api/topology", nil, http.StatusOK), &topo)
	for _, n := range topo {
		if n.Name == "conv" && (strings.Join(n.Reaches, "+") != "any" || strings.Join(n.Emits, "+") != "bdf+any" || len(n.Broken) != 0) {
			t.Fatalf("conv %+v", n)
		}
		if n.Name == "side" && strings.Join(n.Accepts, "+") != "bdf+data+any" {
			t.Fatalf("side %+v", n)
		}
	}

	pic := testPNG(t)
	res, out := e.do(nil, "POST", "/in/raw?name=photo.png", raw.Key, pic)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("png: %d %s", res.StatusCode, out)
	}
	var j Job
	json.Unmarshal(out, &j)
	e.tick()
	e.tick()
	got := e.job(j.Number)
	if got.State != JobCompleted || legs(got) != "raw:routed conv:done made:done side:routed(rule 1) orig:done" || len(got.Payloads) != 2 {
		t.Fatalf("job %s: %+v", legs(got), got.Legs)
	}
	record := func(store, collection string, p *Payload) ArchiveRecord {
		st, err := e.st.archive(ctx, store, collection)
		if err != nil {
			t.Fatal(err)
		}
		r, err := st.Open(ctx, archiveKey(got, p)+".json")
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()
		var rec ArchiveRecord
		json.NewDecoder(r).Decode(&rec)
		return rec
	}
	src, made := got.received(), &got.Payloads[1]
	a, b := record(madeStore, "archive-made", made), record(origStore, "archive-orig", src)
	if a.Trace == nil || b.Trace == nil || a.Trace.ID != got.Legs[1].Ref || b.Trace.ID != a.Trace.ID || a.Trace.Node != "conv" {
		t.Fatalf("traces %+v %+v", a.Trace, b.Trace)
	}
	if a.Trace.Role != TraceOutput || b.Trace.Role != TraceSource || a.Format != "bdf" || b.Format != "image" {
		t.Fatalf("sides %+v %+v", a, b)
	}
	// each names the other's objects in its own archive
	if a.Trace.peer().Object != b.Object || b.Trace.peer().Object != a.Object || a.Trace.Source.Key != archiveKey(got, src) || b.Trace.Output.SHA256 != made.SHA256 {
		t.Fatalf("peers %+v / %s, %+v / %s", a.Trace, a.Object, b.Trace, b.Object)
	}
	fps := func(tr *Trace) string {
		var out []string
		for _, f := range tr.Fingerprints {
			out = append(out, f.Node+":"+f.Role)
		}
		return strings.Join(out, " ")
	}
	if fps(a.Trace) != "conv:output" || fps(b.Trace) != "conv:source" {
		t.Fatalf("fingerprints %s / %s", fps(a.Trace), fps(b.Trace))
	}

	// two steps: what conv made goes through stamp (a .bdf goes on as it
	// is) before the archive; its record keeps conv's fingerprint, which
	// the original's shares
	e.node("stamp", Node{Kind: KindConvert, To: []string{"made"}}, http.StatusOK)
	e.node("conv", Node{Kind: KindConvert, To: []string{"stamp"}, Original: []string{"side"}}, http.StatusOK)
	res, out = e.do(nil, "POST", "/in/raw?name=photo2.png", raw.Key, append(pic, 0))
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("png: %d %s", res.StatusCode, out)
	}
	json.Unmarshal(out, &j)
	e.tick()
	e.tick()
	e.tick()
	got = e.job(j.Number)
	if got.State != JobCompleted || legs(got) != "raw:routed conv:done stamp:done side:routed(rule 1) orig:done made:done" {
		t.Fatalf("job %s: %+v", legs(got), got.Legs)
	}
	a, b = record(madeStore, "archive-made", &got.Payloads[1]), record(origStore, "archive-orig", got.received())
	if fps(a.Trace) != "conv:output stamp:output" || fps(b.Trace) != "conv:source" || a.Trace.Fingerprints[0].ID != b.Trace.ID || a.Trace.ID == b.Trace.ID {
		t.Fatalf("fingerprints %+v / %+v", a.Trace, b.Trace)
	}
}
