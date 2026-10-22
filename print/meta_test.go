package print

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/print/fakeprinter"
	"github.com/shibukawa/bdf/print/ipp"
)

func TestMetaCheck(t *testing.T) {
	long := strings.Repeat("a", 256)
	bad := []Meta{
		{AccountID: long},
		{AccountingUserID: "a\nb"},
		{Metadata: []string{"title"}},
		{Metadata: []string{"foo=bar"}},
		{Metadata: []string{"X-a=1"}},
		{Metadata: []string{"title=" + strings.Repeat("a", 1020)}},
		{Metadata: make([]string, 65)},
		{Custom: map[string]Strings{"order": {"1"}}},
		{Custom: map[string]Strings{"x-Order": {"1"}}},
		{Custom: map[string]Strings{"x-": {"1"}}},
		{Custom: map[string]Strings{"x-order": {}}},
		{Custom: map[string]Strings{"x-order": {"a\tb"}}},
	}
	for i, m := range bad {
		if err := m.Check(); err == nil {
			t.Errorf("%d: %+v passes", i, m)
		}
	}
	good := Meta{AccountID: "dept-7", AccountingUserID: "alice", Metadata: []string{"title=Invoice 42", "accessRights=all", "x-Order.Id=1", "date=2026-10-21"},
		Custom: map[string]Strings{"x-order": {"A-1"}, "x-tags": {"a", "b"}, "x-a.b_c-d": {"1"}}}
	if err := good.Check(); err != nil {
		t.Error(err)
	}
	if !(&Meta{}).Empty() || good.Empty() {
		t.Error("empty")
	}
}

func TestMetaJSON(t *testing.T) {
	tk := Ticket{Copies: 2, Meta: Meta{AccountID: "dept-7", Metadata: []string{"title=Invoice"}, Custom: map[string]Strings{"x-order": {"A-1"}, "x-tags": {"a", "b"}}}}
	b, err := json.Marshal(tk)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"copies":2,"job-account-id":"dept-7","document-metadata":["title=Invoice"],"bdf-attributes":{"x-order":"A-1","x-tags":["a","b"]}}`
	if string(b) != want {
		t.Errorf("json %s", b)
	}
	var back Ticket
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&back); err != nil || !reflect.DeepEqual(back, tk) {
		t.Errorf("back %+v %v", back, err)
	}
	if err := json.Unmarshal([]byte(`{"bdf-attributes":{"x-a":1}}`), &back); err == nil {
		t.Error("a number as a value")
	}
}

func TestMetaMerge(t *testing.T) {
	tk := Ticket{Meta: Meta{AccountID: "a", Custom: map[string]Strings{"x-a": {"1"}}}}
	d := Ticket{Meta: Meta{AccountID: "b", AccountingUserID: "u", Metadata: []string{"title=T"}, Custom: map[string]Strings{"x-a": {"9"}, "x-b": {"2"}}}}
	m := tk.Merge(d).Meta
	want := Meta{AccountID: "a", AccountingUserID: "u", Metadata: []string{"title=T"}, Custom: map[string]Strings{"x-a": {"1"}, "x-b": {"2"}}}
	if !reflect.DeepEqual(m, want) {
		t.Errorf("%+v", m)
	}
	if tk.Custom["x-b"] != nil || d.Custom["x-a"][0] != "9" {
		t.Error("the tickets merged were changed")
	}
}

func TestMetaTake(t *testing.T) {
	var m Meta
	for _, c := range []struct {
		a         ipp.Attribute
		known, ok bool
	}{
		{ipp.Attr("job-account-id", ipp.Name("dept-7")), true, true},
		{ipp.Attr("job-accounting-user-id", ipp.Integer(1)), true, false},
		{ipp.Attr("job-account-id", ipp.Name("a"), ipp.Name("b")), true, false},
		{ipp.Attr("document-metadata", ipp.OctetString([]byte("title=T")), ipp.OctetString([]byte("x-k=v"))), true, true},
		{ipp.Attr("document-metadata", ipp.OctetString([]byte("nokey"))), true, false},
		{ipp.Attr("x-n", ipp.Integer(3), ipp.Boolean(true), ipp.Keyword("k")), true, true},
		{ipp.Attr("x-c", ipp.Collection(ipp.Attr("a", ipp.Integer(1)))), true, false},
		{ipp.Attr("x-Bad", ipp.Name("1")), true, false},
		{ipp.Attr("copies", ipp.Integer(1)), false, false},
	} {
		known, ok := m.Take(c.a)
		if known != c.known || ok != c.ok {
			t.Errorf("%s: known %v ok %v", c.a.Name, known, ok)
		}
	}
	want := Meta{AccountID: "dept-7", Metadata: []string{"title=T", "x-k=v"}, Custom: map[string]Strings{"x-n": {"3", "true", "k"}}}
	if !reflect.DeepEqual(m, want) {
		t.Errorf("%+v", m)
	}
	attrs := m.Attrs()
	if len(attrs) != 3 || attrs[0].Name != "job-account-id" || attrs[1].Name != "x-n" || attrs[1].Values[0].Tag != ipp.TagName ||
		attrs[2].Name != "document-metadata" || attrs[2].Values[0].Tag != ipp.TagOctetString {
		t.Errorf("attrs %v", attrs)
	}
}

// TestPrintMeta sends what a job says about itself: the account to a
// printer that takes one, the attributes of one's own, and
// document-metadata with the document; a printer that knows no attribute
// of one's own says so, and the job is printed without it.
func TestPrintMeta(t *testing.T) {
	p := fakeprinter.New()
	c := serve(t, p)
	tk := Ticket{Meta: Meta{AccountID: "dept-7", AccountingUserID: "alice", Metadata: []string{"title=Invoice 42", "x-order=A-1"},
		Custom: map[string]Strings{"x-cost-center": {"cc-42"}, "x-tags": {"a", "b"}}}}
	res, err := Print(context.Background(), c, testDoc(1), &Options{Ticket: tk})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Unsupported) > 0 {
		t.Errorf("unsupported: %v", res.Unsupported)
	}
	j := p.Jobs()[0]
	if attr(j, "job-account-id") != "dept-7" || attr(j, "job-accounting-user-id") != "alice" || attr(j, "x-cost-center") != "cc-42" ||
		strings.Join(ipp.Get(j.Attrs, "x-tags").Strings(), ",") != "a,b" ||
		strings.Join(ipp.Get(j.Attrs, "document-metadata").Strings(), ";") != "title=Invoice 42;x-order=A-1" {
		t.Errorf("attributes: %v", j.Attrs)
	}
	// a printer that knows no attribute of one's own
	p3 := fakeprinter.New()
	p3.NoCustom = true
	c3 := serve(t, p3)
	res, err = Print(context.Background(), c3, testDoc(1), &Options{Ticket: tk})
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(res.Unsupported, ";"); got != "x-cost-center=cc-42;x-tags=a,b" {
		t.Errorf("unsupported: %q", got)
	}
	if j := p3.Jobs()[0]; attr(j, "x-cost-center") != "" || attr(j, "job-account-id") != "dept-7" {
		t.Errorf("attributes: %v", j.Attrs)
	}
	if _, err := Print(context.Background(), c3, testDoc(1), &Options{Ticket: tk, Strict: true}); !errors.Is(err, ErrUnsupported) {
		t.Errorf("strict: %v", err)
	}
	// one that says so only once the job is made (no Validate-Job)
	p4 := fakeprinter.New()
	p4.NoCustom, p4.IgnoreUnsupported = true, true
	res, err = Print(context.Background(), serve(t, p4), testDoc(1), &Options{Ticket: tk})
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(res.Unsupported, ";"); got != "x-cost-center=cc-42;x-tags=a,b" {
		t.Errorf("unsupported after the fact: %q", got)
	}
	// a ticket that does not pass
	if _, err := Print(context.Background(), c, testDoc(1), &Options{Ticket: Ticket{Meta: Meta{Custom: map[string]Strings{"copies": {"2"}}}}}); err == nil {
		t.Error("an attribute of one's own named as a standard one")
	}
}
