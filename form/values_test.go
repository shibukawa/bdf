package form

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/cms/cmstest"
	"github.com/shibukawa/bdf/sign"
)

// TestValueSignatures: party A fills in theirs and signs the values, party
// B completes the contract and signs theirs; the completed contract holds
// both signatures, which hold for its values.
func TestValueSignatures(t *testing.T) {
	tpl, err := os.ReadFile("../converter/pdfform/testdata/contract.pdf")
	if err != nil {
		t.Fatal(err)
	}
	ca := cmstest.NewAuthority("BDF test root")
	alice, bob := ca.Issue("Alice", false), ca.Issue("Bob", true)
	tsa := cmstest.NewTSA(ca)
	defer tsa.Close()
	a := map[string]any{"party_a": map[string]any{"name": "Alice", "sign": nil}}
	res, err := Fill(bytes.NewReader(tpl), int64(len(tpl)), &Options{Data: a, Open: []string{"party_b"}, Convert: testConvert(),
		Sign: &ValueSigner{Key: alice.Key, Chain: alice.Chain, Claims: map[string]any{"role": "author"}, TSA: tsa.URL}})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Partial.Signed) != 1 {
		t.Fatalf("%d signatures in the partly filled form", len(res.Partial.Signed))
	}
	bdff := write(t, res)
	check := func(b []byte, n int) []*ValueSignature {
		t.Helper()
		sigs, err := VerifyValues(bytes.NewReader(b), int64(len(b)), &sign.VerifyOptions{Roots: ca.Pool()})
		if err != nil {
			t.Fatal(err)
		}
		if len(sigs) != n {
			t.Fatalf("%d signatures, not %d", len(sigs), n)
		}
		return sigs
	}
	if s := check(bdff, 1)[0]; !s.Valid() || s.Values.Claims["role"] != "author" || s.TimeStamp.IsZero() || s.Certificate.Subject.CommonName != "Alice" {
		t.Errorf("the signature of A in the partly filled form: %v %v %+v", s.Err, s.TrustErr, s.Values)
	}

	// B completes it, and signs theirs
	b := map[string]any{"party_b": map[string]any{"name": "ACME", "sign": nil}}
	done, err := Fill(bytes.NewReader(bdff), int64(len(bdff)), &Options{Data: b, Convert: testConvert(), Sign: &ValueSigner{Key: bob.Key, Chain: bob.Chain}})
	if err != nil {
		t.Fatal(err)
	}
	final := write(t, done)
	sigs := check(final, 2)
	for _, s := range sigs {
		if !s.Valid() {
			t.Errorf("signature %d of the completed form: %v %v", s.Index, s.Err, s.TrustErr)
		}
	}
	if len(sigs[0].Values.Open) != 1 || len(sigs[1].Values.Open) != 0 {
		t.Errorf("what was open: %v / %v", sigs[0].Values.Open, sigs[1].Values.Open)
	}
	// the record is a part that no view refers to: the document reads as it did
	rd, _ := bdf.OpenSingle(bytes.NewReader(final), int64(len(final)))
	n := 0
	for _, e := range rd.Manifest.Parts {
		if e.T == bdf.PartFill {
			n++
		}
	}
	if n != 1 || !strings.Contains(text(t, done.Doc), "ACME") {
		t.Errorf("%d records", n)
	}

	// someone changes the name of A in the partly filled form before B
	// completes it: the signature of A no longer holds
	tampered := retemplate(t, bdff, func(def *Template) { def.Filled["party_a"].(map[string]any)["name"] = "Mallory" })
	done, err = Fill(bytes.NewReader(tampered), int64(len(tampered)), &Options{Data: b, Convert: testConvert()})
	if err != nil {
		t.Fatal(err)
	}
	s := check(write(t, done), 1)[0]
	if s.Err == nil || !strings.Contains(s.Err.Error(), `"party_a.name" is not what was signed`) {
		t.Errorf("a changed value: %v", s.Err)
	}
	// or opens a name that A filled
	reopened := retemplate(t, bdff, func(def *Template) { def.Open = append(def.Open, "party_a.name") })
	if s := check(reopened, 1)[0]; s.Err == nil || !strings.Contains(s.Err.Error(), `"party_a.name" was filled`) {
		t.Errorf("a name opened again: %v", s.Err)
	}
	// a form without signed values has none
	plain, _ := Fill(bytes.NewReader(tpl), int64(len(tpl)), &Options{Data: a, Open: []string{"party_b"}, Convert: testConvert()})
	check(write(t, plain), 0)
}

// retemplate writes a partly filled form again with its description
// changed.
func retemplate(t *testing.T, b []byte, change func(*Template)) []byte {
	t.Helper()
	def, _, err := ReadTemplate(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		t.Fatal(err)
	}
	change(def)
	rd, _ := bdf.OpenSingle(bytes.NewReader(b), int64(len(b)))
	old, _ := rd.ToDocument()
	doc := bdf.NewDocument()
	doc.Meta, doc.Views = old.Meta, old.Views
	for _, p := range old.Parts() {
		if p.Type != bdf.PartTemplate {
			doc.AddPart(p.Type, p.Data)
		}
	}
	j, _ := json.Marshal(def)
	doc.AddPart(bdf.PartTemplate, j)
	var out bytes.Buffer
	if err := doc.WriteSingle(&out); err != nil {
		t.Fatal(err)
	}
	return out.Bytes()
}
