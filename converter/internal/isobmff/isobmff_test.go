package isobmff

import (
	"bytes"
	"encoding/binary"
	"slices"
	"testing"
)

func box(typ string, content ...[]byte) []byte {
	b := slices.Concat(content...)
	return slices.Concat(binary.BigEndian.AppendUint32(nil, uint32(8+len(b))), []byte(typ), b)
}

func TestBoxes(t *testing.T) {
	large := slices.Concat([]byte{0, 0, 0, 1}, []byte("big "), binary.BigEndian.AppendUint64(nil, 16+3), []byte("abc"))
	toEnd := slices.Concat([]byte{0, 0, 0, 0}, []byte("rest"), []byte("xyz"))
	data := slices.Concat(box("ftyp", []byte("M4A ")), box("moov", box("mvhd", []byte{1, 2})), large, toEnd)
	bs := Boxes(data)
	if len(bs) != 4 || bs[0].Type != "ftyp" || string(bs[0].Data) != "M4A " || bs[2].Type != "big " || string(bs[2].Data) != "abc" ||
		bs[3].Type != "rest" || string(bs[3].Data) != "xyz" {
		t.Fatalf("Boxes = %v", bs)
	}
	if got := Find(Boxes(Find(bs, "moov")), "mvhd"); !bytes.Equal(got, []byte{1, 2}) {
		t.Errorf("mvhd = %v", got)
	}
	if Find(bs, "none") != nil {
		t.Error("found a box that is not there")
	}
	// a box that claims more than the data holds ends the list
	bad := slices.Concat(box("a   "), []byte{0, 0, 1, 0}, []byte("b   "), []byte{1, 2, 3})
	if bs := Boxes(bad); len(bs) != 1 {
		t.Errorf("Boxes of a truncated file = %v", bs)
	}
	tops := Walk(bytes.NewReader(data), int64(len(data)))
	if len(tops) != 4 || tops[1].Type != "moov" || tops[1].Off != 12 || tops[2].Header != 16 || tops[3].Size != 11 {
		t.Fatalf("Walk = %+v", tops)
	}
	if c, err := tops[1].Contents(bytes.NewReader(data)); err != nil || !bytes.Equal(c, box("mvhd", []byte{1, 2})) {
		t.Errorf("Contents = %v, %v", c, err)
	}
	if tops := Walk(bytes.NewReader(bad), int64(len(bad))); len(tops) != 1 {
		t.Errorf("Walk of a truncated file = %+v", tops)
	}
	if got := FindAll(Boxes(slices.Concat(box("data", []byte("1")), box("data", []byte("2")))), "data"); len(got) != 2 {
		t.Errorf("FindAll = %q", got)
	}
}
