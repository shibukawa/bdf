package bdf

import (
	"bytes"
	"encoding/json"
	"slices"
	"testing"
)

func TestDCValuesJSON(t *testing.T) {
	dc := DublinCore{Title: DCValues{"Report"}, Creator: DCValues{"A", "B"}, Subject: DCValues{}}
	b, err := json.Marshal(dc)
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"title":"Report","creator":["A","B"]}`; string(b) != want {
		t.Fatalf("marshal = %s, want %s", b, want)
	}
	var got DublinCore
	if err := json.Unmarshal([]byte(`{"title":"Report","creator":["A","B"],"subject":null,"rights":[],"unknown":"x"}`), &got); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got.Title, DCValues{"Report"}) || !slices.Equal(got.Creator, DCValues{"A", "B"}) || got.Subject != nil || len(got.Rights) != 0 {
		t.Fatalf("unmarshal = %+v", got)
	}
	for _, bad := range []string{`{"title":1}`, `{"title":["a",2]}`, `{"title":{"@value":"a"}}`} {
		if err := json.Unmarshal([]byte(bad), &got); err == nil {
			t.Errorf("%s accepted", bad)
		}
	}
	if b, _ := json.Marshal(Meta{DC: DublinCore{Rights: DCValues{}}, Generator: "g"}); string(b) != `{"generator":"g"}` {
		t.Errorf("empty dc not omitted: %s", b)
	}
}

func TestDCField(t *testing.T) {
	// Field(name) addresses the element whose JSON key is name.
	var dc DublinCore
	for _, name := range DCTerms {
		f := dc.Field(name)
		if f == nil {
			t.Fatalf("Field(%q) = nil", name)
		}
		*f = DCValues{name}
	}
	b, _ := json.Marshal(dc)
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if len(m) != len(DCTerms) {
		t.Fatalf("%d elements set, want %d: %s", len(m), len(DCTerms), b)
	}
	for k, v := range m {
		if k != v {
			t.Errorf("Field(%q) sets %q", v, k)
		}
	}
	if dc.Field("Title") != nil || dc.Field("alternative") != nil {
		t.Error("unknown names must give nil")
	}
}

func TestSplitKeywords(t *testing.T) {
	got := SplitKeywords(" pdf, bdf ;canvas；日本語、検索，, ")
	if want := (DCValues{"pdf", "bdf", "canvas", "日本語", "検索"}); !slices.Equal(got, want) {
		t.Fatalf("SplitKeywords = %q, want %q", got, want)
	}
	if SplitKeywords("  ") != nil {
		t.Fatal("blank keywords must give nil")
	}
}

func TestMagic(t *testing.T) {
	var buf bytes.Buffer
	if err := NewDocument().WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(buf.Bytes(), []byte("bdf\x00")) {
		t.Fatalf("header starts with %q", buf.Bytes()[:4])
	}
	b := buf.Bytes()
	copy(b, "BDF1")
	if _, err := OpenSingle(bytes.NewReader(b), int64(len(b))); err == nil {
		t.Fatal("bad magic accepted")
	}
}
