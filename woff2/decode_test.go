package woff2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
)

func TestDecodeRoundTrip(t *testing.T) {
	if !Available() {
		t.Skip("built without Brotli")
	}
	for name, font := range testFonts(t) {
		t.Run(name, func(t *testing.T) {
			w, err := Encode(font)
			if err != nil {
				t.Fatal(err)
			}
			dec, err := Decode(w)
			if err != nil {
				t.Fatal(err)
			}
			if binary.BigEndian.Uint32(dec) != binary.BigEndian.Uint32(font) {
				t.Fatal("flavor changed")
			}
			orig, got := sfntTables(font), sfntTables(dec)
			if len(orig) != len(got) {
				t.Fatalf("%d tables, want %d", len(got), len(orig))
			}
			for tag, data := range orig {
				switch tag {
				case "glyf", "loca":
				case "head":
					a, b := bytes.Clone(data), bytes.Clone(got[tag])
					for _, r := range [][2]int{{8, 12}, {16, 18}, {50, 52}} { // checksum adjustment, flags, loca format
						clear(a[r[0]:r[1]])
						clear(b[r[0]:r[1]])
					}
					if !bytes.Equal(a, b) {
						t.Error("head changed")
					}
				default:
					if !bytes.Equal(got[tag], data) {
						t.Errorf("table %s changed", tag)
					}
				}
			}
			n := int(binary.BigEndian.Uint16(orig["maxp"][4:]))
			want := parseGlyf(t, orig["glyf"], orig["loca"], n, binary.BigEndian.Uint16(orig["head"][50:]) != 0)
			have := parseGlyf(t, got["glyf"], got["loca"], n, true)
			for g := range want {
				if !reflect.DeepEqual(fmt.Sprint(want[g]), fmt.Sprint(have[g])) {
					t.Fatalf("glyph %d differs:\n want %v\n got  %v", g, want[g], have[g])
				}
			}
		})
	}
}

func TestDecodeRejectsGarbage(t *testing.T) {
	for _, b := range [][]byte{nil, []byte("wOF2"), append([]byte("wOF2"), make([]byte, 60)...), []byte("OTTO0000")} {
		if _, err := Decode(b); err == nil {
			t.Errorf("Decode(%q) succeeded", b)
		}
	}
}
