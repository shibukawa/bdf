package bdf

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"errors"
	"testing"
)

// segmentSample is a document of 7 pages that all use one shared object
// (standing for a font or a master), each with an object of its own.
func segmentSample(t *testing.T) (*Reader, Hash) {
	t.Helper()
	d := NewDocument()
	d.Meta.DC.Title = DCValues{"secret title"}
	shared := NewObject()
	for i := range 100 {
		shared.FillRect(float32(i), 0, 1, 1)
	}
	sh, bbox := d.AddObject(shared)
	v := d.NewView("pages", ViewFixed, "")
	for i := range 7 {
		o := NewObject()
		o.Use(o.AddObject(sh, bbox))
		o.FillRect(float32(i), float32(i), 2, 2)
		h, _ := d.AddObject(o)
		v.AddPage(100, 100, Layer{Role: RoleBody, Obj: h})
	}
	var b bytes.Buffer
	if err := d.WriteSingle(&b); err != nil {
		t.Fatal(err)
	}
	r, err := OpenSingle(bytes.NewReader(b.Bytes()), int64(b.Len()))
	if err != nil {
		t.Fatal(err)
	}
	return r, sh
}

func TestSegmentSealedForKey(t *testing.T) {
	r, shared := segmentSample(t)
	priv, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	lock, err := NewECDHLock(priv.PublicKey())
	if err != nil {
		t.Fatal(err)
	}
	s := SegmentAt(r.Manifest.Views[0], 4, 3)
	if s != (Segment{View: "pages", From: 3, To: 6}) {
		t.Fatalf("SegmentAt = %+v", s)
	}
	var b bytes.Buffer
	if err := r.WriteSegment(&b, SegmentOptions{Segment: s, Have: []Segment{{View: "pages", From: 0, To: 3}}, Lock: lock}); err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(b.Bytes(), []byte("secret title")) {
		t.Fatal("the title is in the clear")
	}
	seg, err := OpenSingle(bytes.NewReader(b.Bytes()), int64(b.Len()))
	if err != nil {
		t.Fatal(err)
	}
	other, _ := ecdh.P256().GenerateKey(rand.Reader)
	if err := seg.UnlockECDH(other); !errors.Is(err, ErrWrongKey) {
		t.Fatalf("another key: %v", err)
	}
	if err := seg.Unlock("password"); !errors.Is(err, ErrWrongPassword) {
		t.Fatalf("a password: %v", err)
	}
	if err := seg.UnlockECDH(priv); err != nil {
		t.Fatal(err)
	}
	m := seg.Manifest
	if m.Segment == nil || *m.Segment != s {
		t.Fatalf("segment = %+v", m.Segment)
	}
	pages := m.Views[0].Pages
	if len(pages) != 7 {
		t.Fatalf("%d pages", len(pages))
	}
	for i, p := range pages {
		if carried := i >= 3 && i < 6; carried != (len(p.Layers) > 0) {
			t.Errorf("page %d: %d layers", i, len(p.Layers))
		}
	}
	// the shared object came with pages 0-2, which the reader holds
	if _, ok := seg.Entry(shared); ok {
		t.Error("the shared object was sent again")
	}
	for i := 3; i < 6; i++ {
		obj := pages[i].Layers[0].Obj
		got, err := seg.Part(obj)
		if err != nil {
			t.Fatal(err)
		}
		want, _ := r.Part(obj)
		if !bytes.Equal(got, want) {
			t.Errorf("page %d differs", i)
		}
	}
	// without Have the segment brings the shared object
	b.Reset()
	if err := r.WriteSegment(&b, SegmentOptions{Segment: s}); err != nil {
		t.Fatal(err)
	}
	plain, err := OpenSingle(bytes.NewReader(b.Bytes()), int64(b.Len()))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := plain.Entry(shared); !ok || plain.Encrypted() {
		t.Error("a segment without Have or Lock")
	}
}

func TestSegmentRefusesSheets(t *testing.T) {
	d := NewDocument()
	d.NewView("s", ViewSheet, "")
	var b bytes.Buffer
	if err := d.WriteSingle(&b); err != nil {
		t.Fatal(err)
	}
	r, _ := OpenSingle(bytes.NewReader(b.Bytes()), int64(b.Len()))
	if err := r.WriteSegment(&bytes.Buffer{}, SegmentOptions{Segment: Segment{View: "s", From: 0, To: 1}}); err == nil {
		t.Fatal("a sheet was cut into segments")
	}
}
