package woff2

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"testing"

	"github.com/andybalholm/brotli"
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

// packed assembles a WOFF2 file of a directory and the Brotli stream of
// its tables.
func packed(t *testing.T, numTables int, dir []byte, write func(w io.Writer)) []byte {
	t.Helper()
	var comp bytes.Buffer
	w := brotli.NewWriterOptions(&comp, brotli.WriterOptions{Quality: 1})
	write(w)
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	out := make([]byte, 48)
	copy(out, "wOF2")
	binary.BigEndian.PutUint32(out[4:], 0x00010000)
	binary.BigEndian.PutUint16(out[12:], uint16(numTables))
	binary.BigEndian.PutUint32(out[20:], uint32(comp.Len()))
	out = append(out, dir...)
	out = append(out, comp.Bytes()...)
	binary.BigEndian.PutUint32(out[8:], uint32(len(out)))
	return out
}

// allocated returns the bytes fn allocates.
func allocated(fn func()) uint64 {
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	fn()
	runtime.ReadMemStats(&after)
	return after.TotalAlloc - before.TotalAlloc
}

// TestDecodeBounded reads fonts that say little and hold much: what
// Decode allocates follows what the directory says, not the stream.
func TestDecodeBounded(t *testing.T) {
	if !Available() {
		t.Skip("built without Brotli")
	}
	// a stream of 32 MiB for a directory of one table of 54 bytes
	head := appendBase128([]byte{1}, 54)
	long := packed(t, 1, head, func(w io.Writer) {
		chunk := make([]byte, 1<<20)
		for range 32 {
			w.Write(chunk)
		}
	})
	var err error
	if n := allocated(func() { _, err = Decode(long) }); err != nil || n > 8<<20 {
		t.Errorf("a stream longer than its tables: %d bytes allocated, error %v", n, err)
	}

	// tables of more than maxSize together
	var dir []byte
	for range 5 {
		dir = appendBase128(append(dir, 0x3f, 'T', 'E', 'S', 'T'), 64<<20)
	}
	large := packed(t, 5, dir, func(w io.Writer) { w.Write(make([]byte, 16)) })
	if _, err := Decode(large); !errors.Is(err, ErrTooLarge) {
		t.Errorf("tables of 320 MiB: error %v, want ErrTooLarge", err)
	}

	// glyphs of 65535 points each, of streams that hold no point
	const glyphs = 64
	var nContour, nPoints []byte
	for range glyphs {
		nContour = append(nContour, 0, 1)
		nPoints = append(nPoints, 253, 0xff, 0xff)
	}
	bbox := make([]byte, ((glyphs+31)>>5)<<2)
	tr := make([]byte, 36)
	binary.BigEndian.PutUint16(tr[4:], glyphs)
	for i, s := range [][]byte{nContour, nPoints, nil, nil, nil, bbox, nil} {
		binary.BigEndian.PutUint32(tr[8+4*i:], uint32(len(s)))
		tr = append(tr, s...)
	}
	dir = appendBase128(appendBase128([]byte{10}, 100), uint32(len(tr))) // glyf
	dir = appendBase128(appendBase128(append(dir, 11), 100), 0)          // loca
	empty := packed(t, 2, dir, func(w io.Writer) { w.Write(tr) })
	if n := allocated(func() { _, err = Decode(empty) }); err == nil || n > 8<<20 {
		t.Errorf("glyphs of streams that ended: %d bytes allocated, error %v", n, err)
	}
}
