package bdf

import (
	"bytes"
	"testing"
)

// twoPages builds page objects that share a "master" prefix (background,
// title band, a heading in a font) and differ afterwards.
func twoPages(t *testing.T) (*Document, []*Object, Hash) {
	t.Helper()
	doc := NewDocument()
	img := doc.AddPart(PartImage, []byte("not really a png"))
	build := func(body string, color Color) *Object {
		o := NewObject().SetBBox(0, 0, 400, 300)
		o.Transform(1, 0, 0, -1, 0, 300) // the converter's page transform
		o.FillColor(RGB(240, 240, 240)).FillRect(0, 0, 400, 300)
		o.Save().FillColor(RGB(31, 58, 95)).FillRect(0, 0, 400, 40).Restore()
		f := o.AddFont(SystemFont("serif", 400, 0))
		o.Font(f, 12).FillColor(RGB(10, 10, 10))
		o.Image(o.AddImage(img), 10, 50, 20, 20)
		o.Mark(1, "").FillText("Master heading", 20, 30, 90) // ends the common part
		// page-specific
		o.FillColor(color)
		o.Save().Translate(0, 100).FillText(body, 20, 0, 60).Restore()
		return o
	}
	return doc, []*Object{build("page one", RGB(200, 0, 0)), build("page two", RGB(0, 0, 200))}, img
}

func TestSharePrefixesSplitsAtCommonPrefix(t *testing.T) {
	doc, objs, _ := twoPages(t)
	shares, err := SharePrefixes(objs, 16)
	if err != nil {
		t.Fatal(err)
	}
	if len(shares) != 1 || len(shares[0].Members) != 2 {
		t.Fatalf("shares = %+v", shares)
	}
	sh := shares[0]
	prefixHash, _ := doc.AddObject(sh.Prefix)
	prefixText, err := Disassemble(sh.Prefix.Encode())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains([]byte(prefixText), []byte("Master heading")) {
		t.Fatalf("prefix lacks the heading:\n%s", prefixText)
	}
	// The prefix ends after the heading, not after the MARK before it, and
	// the rest re-establishes transform, font and fill colour before drawing.
	for _, m := range sh.Members {
		m.Rest.UpdateObject(0, prefixHash)
		part, err := DecodeObject(m.Rest.Encode())
		if err != nil {
			t.Fatal(err)
		}
		ins, err := part.Instructions()
		if err != nil {
			t.Fatal(err)
		}
		ops := []byte{}
		for _, in := range ins {
			ops = append(ops, in.Op)
		}
		want := []byte{OpUse, OpTransform, OpFillColor, OpFont, OpFillColor, OpSave, OpTranslate, OpFillText, OpRestore}
		if !bytes.Equal(ops, want) {
			t.Fatalf("rest ops = % x, want % x\n%s", ops, want, must(Disassemble(m.Rest.Encode())))
		}
		if ins[1].Args[3].(float32) != -1 || ins[1].Args[5].(float32) != 300 {
			t.Fatalf("rest transform = %v", ins[1].Args)
		}
		if part.Objects[0] != prefixHash || len(m.RestObjects) != 1 || m.RestObjects[0] != -1 {
			t.Fatalf("rest child = %v, map %v", part.Objects, m.RestObjects)
		}
		if len(m.RestFonts) != 1 || m.RestFonts[0] != 0 {
			t.Fatalf("rest fonts map = %v", m.RestFonts)
		}
		if m.Saved < 16 {
			t.Fatalf("saved = %d", m.Saved)
		}
	}
	if len(sh.PrefixFonts) != 1 || sh.PrefixFonts[0] != 0 {
		t.Fatalf("prefix fonts map = %v", sh.PrefixFonts)
	}
	// Text extraction over the rewritten pages reads the same runs as before.
	resolve := func(h Hash) *ObjectPart {
		p, _ := DecodeObject(doc.Part(h).Data)
		return p
	}
	for i, m := range sh.Members {
		before, err := ExtractText(must(DecodeObject(objs[i].Encode())), resolve)
		if err != nil {
			t.Fatal(err)
		}
		after, err := ExtractText(must(DecodeObject(m.Rest.Encode())), resolve)
		if err != nil {
			t.Fatal(err)
		}
		if len(before) != len(after) {
			t.Fatalf("page %d: %d runs before, %d after", i, len(before), len(after))
		}
		for j := range before {
			if before[j].Text != after[j].Text || before[j].X != after[j].X || before[j].Y != after[j].Y || before[j].Sep != after[j].Sep {
				t.Fatalf("page %d run %d: %+v vs %+v", i, j, before[j], after[j])
			}
		}
	}
}

func TestSharePrefixesRespectsDepthAndClip(t *testing.T) {
	// A prefix that never returns to depth 0 before the pages diverge, and one
	// with a top-level clip, must not be shared.
	mk := func(clip bool, body string) *Object {
		o := NewObject().SetBBox(0, 0, 100, 100)
		if clip {
			o.ClipRect(0, 0, 100, 100)
		}
		o.FillColor(RGB(1, 2, 3)).FillRect(0, 0, 100, 100).FillRect(1, 1, 50, 50).FillRect(2, 2, 30, 30)
		o.Save().FillText(body, 0, 0, 10)
		o.Restore()
		return o
	}
	shares, err := SharePrefixes([]*Object{mk(true, "a"), mk(true, "b")}, 8)
	if err != nil {
		t.Fatal(err)
	}
	if len(shares) != 0 {
		t.Fatalf("clipped pages were shared: %+v", shares)
	}
	shares, err = SharePrefixes([]*Object{mk(false, "a"), mk(false, "b")}, 8)
	if err != nil {
		t.Fatal(err)
	}
	if len(shares) != 1 {
		t.Fatalf("unclipped pages not shared")
	}
	text := must(Disassemble(shares[0].Prefix.Encode()))
	if bytes.Contains([]byte(text), []byte("SAVE")) {
		t.Fatalf("cut inside the save/restore pair:\n%s", text)
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// TestReplay copies instructions of an object that was read into another:
// those of a part of it, with the resources they refer to and no others.
func TestReplay(t *testing.T) {
	doc := NewDocument()
	img := doc.AddPart(PartImage, []byte("not really a png"))
	child, box := doc.AddObject(NewObject().SetBBox(0, 0, 5, 5).FillRect(0, 0, 5, 5))
	o := NewObject().SetBBox(0, 0, 100, 100)
	o.FillPath(o.AddPath((&Path{}).Rect(0, 0, 1, 1)), NonZero) // (not copied)
	o.Font(o.AddFont(SystemFont("serif", 400, 0)), 12)
	o.FillPaint(o.AddPaint(Paint{Kind: PaintPattern, Image: o.AddImage(img)}))
	o.FillPath(o.AddPath((&Path{}).Circle(5, 5, 3)), EvenOdd).FillText("text", 1, 2, 3)
	o.Image(o.AddImage(img), 1, 2, 3, 4).Use(o.AddObject(child, box)).Mark(MarkFigure, "alt")
	part, err := DecodeObject(o.Encode())
	if err != nil {
		t.Fatal(err)
	}
	ins, err := part.Instructions()
	if err != nil {
		t.Fatal(err)
	}
	dst := NewObject().SetBBox(0, 0, 100, 100)
	dst.FillPath(dst.AddPath((&Path{}).Rect(9, 9, 9, 9)), NonZero) // what the destination has of its own
	rp := NewReplay(part, dst)
	for _, in := range ins[1:] {
		if err := rp.Emit(in); err != nil {
			t.Fatal(err)
		}
	}
	got, err := DecodeObject(dst.Encode())
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Paths) != 2 || len(got.Paths[1].Inline.Verbs) != len(part.Paths[1].Inline.Verbs) || len(got.Fonts) != 1 || len(got.Paints) != 1 ||
		len(got.Images) != len(part.Images) || got.Images[0] != img || len(got.Objects) != 1 || got.Objects[0] != child {
		t.Errorf("resources: %d paths, %d fonts, %d paints, %v images, %v objects", len(got.Paths), len(got.Fonts), len(got.Paints), got.Images, got.Objects)
	}
	out, err := got.Instructions()
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != len(ins) {
		t.Fatalf("%d instructions, want %d", len(out), len(ins))
	}
	for i := 1; i < len(ins); i++ {
		if out[i].Op != ins[i].Op {
			t.Errorf("instruction %d is 0x%02x, want 0x%02x", i, out[i].Op, ins[i].Op)
		}
	}
	// the path of the copied FILL_PATH is the second of the destination
	if out[3].Args[0].(uint64) != 1 || out[3].Args[1].(uint64) != uint64(EvenOdd) {
		t.Errorf("the path of the copy: %v", out[3].Args)
	}
	// an instruction that refers to what the part lacks is an error
	if err := rp.Emit(Instr{Op: OpFillPath, Args: []any{uint64(9), uint64(0)}}); err == nil {
		t.Error("no error for a path the object does not have")
	}
}
