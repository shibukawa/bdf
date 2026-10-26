package ipp_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/flowtest"
	dipp "github.com/shibukawa/bdf/docflow/io/ipp"
	"github.com/shibukawa/bdf/print/printconv"
)

// A CUPS server (sink) prints to a queue of a CUPS server over IPP: here
// the port of a CUPS server (source) of the same manager, so that what
// it prints comes back as a job of its own, kept at a sink. Its queues
// are listed before it is saved.
func TestCupsOut(t *testing.T) {
	// the CUPS server (source) makes .bdf of the PDF printed to it
	e := flowtest.New(t, func(o *docflow.Options) { o.Converter = printconv.Converter{NoSystemFonts: true} })
	ctx := context.Background()
	e.Put(docflow.Node{Name: "keep", Kind: docflow.KindSink},
		docflow.Node{Name: "office", Kind: dipp.KindCups, Title: "Office", To: []string{"keep"}})
	// the CUPS server (source)'s own port, as cmd/docflow serves it
	cups := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		docflow.PortHandler(e.M(), "office").ServeHTTP(w, r)
	}))
	defer cups.Close()
	server := strings.Replace(cups.URL, "http://", "ipp://", 1)

	k := kind(t, dipp.KindCupsOut)
	got, err := k.Discover(ctx, e.M(), docflow.Node{URL: server})
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(got)
	var qs []dipp.Queue
	json.Unmarshal(b, &qs)
	if len(qs) != 1 || qs[0].Name != "office" || qs[0].Path != "/printers/office" || qs[0].Info != "Office" || !qs[0].Accepting || qs[0].Class {
		t.Fatalf("queues %+v", qs)
	}
	if _, err := k.Discover(ctx, e.M(), docflow.Node{URL: "http://" + cups.Listener.Addr().String()}); !errors.Is(err, docflow.ErrInvalid) {
		t.Fatalf("an http URL: %v", err)
	}
	// the editor asks through the API; a kind with no far side is not found
	for path, want := range map[string]int{"/api/kinds/cups-out/discover": http.StatusOK, "/api/kinds/queue/discover": http.StatusNotFound} {
		req, _ := http.NewRequest("POST", e.Srv.URL+path, strings.NewReader(`{"url": "`+server+`"}`))
		req.Header.Set("Authorization", "Bearer "+flowtest.AdminToken)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		var body []dipp.Queue
		json.NewDecoder(res.Body).Decode(&body)
		res.Body.Close()
		if res.StatusCode != want || want == http.StatusOK && (len(body) != 1 || body[0].Name != "office") {
			t.Errorf("%s: %d %+v", path, res.StatusCode, body)
		}
	}

	for _, bad := range []docflow.Node{
		{Name: "out", Kind: dipp.KindCupsOut, URL: server},
		{Name: "out", Kind: dipp.KindCupsOut, URL: server + "/jobs/1"},
		{Name: "out", Kind: dipp.KindCupsOut, URL: strings.Replace(server, "ipp:", "http:", 1) + "/printers/office"},
		{Name: "out", Kind: dipp.KindCupsOut, URL: server + "/printers/office", Credential: "../secret"},
	} {
		if _, err := e.M().PutNode(ctx, bad); !errors.Is(err, docflow.ErrInvalid) {
			t.Errorf("%+v: %v", bad, err)
		}
	}
	e.Put(docflow.Node{Name: "out", Kind: dipp.KindCupsOut, URL: server + qs[0].Path})
	j := e.Submit("out", "statement", flowtest.Doc(t, 2), nil)
	e.Tick()
	if d := e.Delivery(j, "out"); d.State != docflow.DeliveryDone || d.Result == "" {
		t.Fatalf("delivery %+v", d)
	}
	// it came back through the CUPS server (source), as a job of its own
	back := e.Job(2)
	if back.Dest != "office" || back.Source != "ipp" || back.Name != "statement" || back.Pages != 2 {
		t.Fatalf("the job printed %+v", back)
	}
}

func kind(t *testing.T, name string) docflow.Kind {
	for _, k := range docflow.Kinds() {
		if k.Name == name {
			return k
		}
	}
	t.Fatalf("no kind %s", name)
	return docflow.Kind{}
}
