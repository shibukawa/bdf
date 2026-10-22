package docflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/internal/cms/cmstest"
	"github.com/shibukawa/bdf/sign"
)

// A sign node signs the document that reaches it with the manager's key,
// as a witness of who sent it, and sends it on.
func TestSignNode(t *testing.T) {
	ca := cmstest.NewAuthority("BDF test root")
	tsa := cmstest.NewTSA(ca)
	defer tsa.Close()
	id := ca.Issue("BDF Docflow", false)
	e := newEnv(t, func(o *Options) { o.Signer = &Signer{Key: id.Key, Chain: id.Chain, TSA: tsa.URL} })
	e.node("keep", Node{Kind: KindSink}, http.StatusOK)
	e.node("witness", Node{Kind: KindSign}, http.StatusBadRequest)
	e.node("witness", Node{Kind: KindSign, To: []string{"keep"}, SignRole: strings.Repeat("x", 65)}, http.StatusBadRequest)
	e.node("witness", Node{Kind: KindSign, To: []string{"keep"}}, http.StatusOK)
	// it signs .bdf: data does not reach it
	if out := e.node("din", Node{Kind: KindEndpoint, Accept: AcceptData, To: []string{"witness"}}, http.StatusBadRequest); !strings.Contains(string(out), "takes bdf") {
		t.Fatalf("data to a sign node: %s", out)
	}
	var in struct {
		Key string `json:"key"`
	}
	json.Unmarshal(e.node("in", Node{Kind: KindEndpoint, To: []string{"witness"}}, http.StatusOK), &in)
	res, out := e.do(nil, "POST", "/in/in?name=contract", in.Key, testDoc(t, 1))
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("%d %s", res.StatusCode, out)
	}
	var j Job
	json.Unmarshal(out, &j)
	if legs(&j) != "in:routed witness:pending" {
		t.Fatalf("legs %s", legs(&j))
	}
	if rep := e.tick(); rep.Delivered != 1 {
		t.Fatalf("tick %+v", rep)
	}
	got := e.job(j.Number)
	if legs(got) != "in:routed witness:done keep:kept" || len(got.Payloads) != 2 || got.Payloads[1].Made != "witness" || got.Payloads[1].Pages != 1 {
		t.Fatalf("job %s %+v", legs(got), got.Payloads)
	}
	res, doc := e.do(nil, "GET", fmt.Sprintf("/api/jobs/%d/document", j.Number), adminToken, nil)
	if res.StatusCode != http.StatusOK {
		t.Fatal(res.StatusCode)
	}
	d, err := sign.Read(bytes.NewReader(doc), int64(len(doc)))
	if err != nil {
		t.Fatal(err)
	}
	sigs, err := d.Verify(&sign.VerifyOptions{Roots: ca.Pool()})
	if err != nil || len(sigs) != 1 {
		t.Fatalf("%d signatures: %v", len(sigs), err)
	}
	s := sigs[0]
	by, _ := s.Claims["submitted_by"].(map[string]any)
	if !s.Valid() || s.Claims["role"] != RoleWitness || s.Claims["node"] != "witness" || s.Claims["job"] != j.ID || by["kind"] == nil || s.TimeStamp.IsZero() {
		t.Errorf("the signature: %v %v %+v", s.Err, s.TrustErr, s.Claims)
	}

	// a manager without a key fails the step, for good
	e2 := newEnv(t, nil)
	e2.node("keep", Node{Kind: KindSink}, http.StatusOK)
	e2.node("witness", Node{Kind: KindSign, To: []string{"keep"}, SignRole: "notary"}, http.StatusOK)
	json.Unmarshal(e2.node("in", Node{Kind: KindEndpoint, To: []string{"witness"}}, http.StatusOK), &in)
	_, out = e2.do(nil, "POST", "/in/in?name=contract", in.Key, testDoc(t, 1))
	json.Unmarshal(out, &j)
	e2.tick()
	if got := e2.job(j.Number); !strings.HasPrefix(legs(got), "in:routed witness:failed") || !strings.Contains(got.Legs[1].Error, "no key to sign with") {
		t.Errorf("without a key: %s %+v", legs(got), got.Legs)
	}
}
