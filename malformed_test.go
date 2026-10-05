package bdf

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math"
	"testing"
	"time"
)

// header writes a single-file header with the manifest at offset 32.
func header(flags uint16, moff, mlen uint64, menc byte) []byte {
	var h buf
	h.bytes(Magic[:])
	h.u16(FormatVersion)
	h.u16(flags)
	h.u64(moff)
	h.u64(mlen)
	h.u8(menc)
	h.bytes(make([]byte, 7))
	return h.b
}

func openBytes(b []byte) (*Reader, error) {
	return OpenSingle(bytes.NewReader(b), int64(len(b)))
}

func TestOpenSingleManifestRange(t *testing.T) {
	for _, c := range []struct {
		name       string
		moff, mlen uint64
	}{
		{"negative length", HeaderSize, 1 << 63},
		{"negative offset", 1 << 63, 8},
		{"sum overflows", math.MaxInt64, math.MaxInt64},
		{"offset and length cancel", 1<<64 - 8, 16},
		{"past the end", HeaderSize, 1 << 20},
	} {
		t.Run(c.name, func(t *testing.T) {
			b := append(header(0, c.moff, c.mlen, 0), `{"bdf":1,"parts":[]}`...)
			if _, err := openBytes(b); err == nil {
				t.Fatal("no error")
			}
		})
	}
}

// single writes a document with a manifest and the stored bytes of its parts.
func single(t *testing.T, m *Manifest, data []byte) []byte {
	t.Helper()
	mj, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	b := append(header(0, HeaderSize, uint64(len(mj)), 0), mj...)
	return append(b, data...)
}

func TestPartRange(t *testing.T) {
	data := []byte("0123456789")
	h := HashOf(data)
	for _, c := range []struct {
		name string
		e    PartEntry
	}{
		{"negative length", PartEntry{Len: -1, Size: 10}},
		{"negative offset", PartEntry{Len: 10, Size: 10, Off: -4}},
		{"offset overflows", PartEntry{Len: 10, Size: 10, Off: math.MaxInt64}},
		{"past the end", PartEntry{Len: 11, Size: 11}},
	} {
		t.Run(c.name, func(t *testing.T) {
			c.e.H, c.e.T, c.e.Enc = h, PartImage, EncIdentity
			r, err := openBytes(single(t, &Manifest{BDF: 1, Parts: []PartEntry{c.e}}, data))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := r.Part(h); err == nil {
				t.Fatal("no error")
			}
		})
	}
}

func TestInflateMaximumIntegerLimit(t *testing.T) {
	plain := []byte("a valid compressed part")
	packed, err := compress(plain, 9)
	if err != nil {
		t.Fatal(err)
	}
	for _, limit := range []int{len(plain), math.MaxInt} {
		got, err := inflate(packed, limit, limit)
		if err != nil || !bytes.Equal(got, plain) {
			t.Errorf("limit %d: %q, %v", limit, got, err)
		}
	}
	if _, err := inflate(packed, len(plain)-1, 0); err == nil {
		t.Error("data larger than the limit was accepted")
	}
}

func TestInflateLimits(t *testing.T) {
	// 64 MiB of zeros deflate to a few dozen kilobytes
	plain := make([]byte, 64<<20)
	packed, err := compress(plain, 9)
	if err != nil {
		t.Fatal(err)
	}
	h := HashOf(plain)

	t.Run("part larger than stated", func(t *testing.T) {
		e := PartEntry{H: h, T: PartImage, Enc: EncDeflateRaw, Len: len(packed), Size: 1024}
		r, err := openBytes(single(t, &Manifest{BDF: 1, Parts: []PartEntry{e}}, packed))
		if err != nil {
			t.Fatal(err)
		}
		var fe *FormatError
		if _, err := r.Part(h); !errors.As(err, &fe) {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("part as stated", func(t *testing.T) {
		e := PartEntry{H: h, T: PartImage, Enc: EncDeflateRaw, Len: len(packed), Size: len(plain)}
		r, err := openBytes(single(t, &Manifest{BDF: 1, Parts: []PartEntry{e}}, packed))
		if err != nil {
			t.Fatal(err)
		}
		if b, err := r.Part(h); err != nil || len(b) != len(plain) {
			t.Fatalf("%d bytes, err = %v", len(b), err)
		}
	})
	t.Run("manifest", func(t *testing.T) {
		big := make([]byte, MaxManifestSize+1)
		for i := range big {
			big[i] = ' '
		}
		copy(big, `{"bdf":1,"parts":[]}`)
		packed, err := compress(big, 1)
		if err != nil {
			t.Fatal(err)
		}
		var fe *FormatError
		if _, err := openBytes(append(header(0, HeaderSize, uint64(len(packed)), 1), packed...)); !errors.As(err, &fe) {
			t.Fatalf("err = %v", err)
		}
	})
}

// object encodes an object part with the given tables and op stream; the
// tables are the encoded bytes after their counts.
func object(strings []string, paths, paints [][]byte, nObjects int, ops []byte) []byte {
	var w buf
	w.bytes([]byte("BOBJ"))
	w.u16(OpsetVersion)
	w.u16(0)
	w.f32s(0, 0, 10, 10)
	w.varuint(uint64(len(strings)))
	for _, s := range strings {
		w.str(s)
	}
	w.varuint(uint64(len(paths)))
	for _, p := range paths {
		w.bytes(p)
	}
	w.varuint(uint64(len(paints)))
	for _, p := range paints {
		w.bytes(p)
	}
	w.varuint(0) // fonts
	w.varuint(0) // images
	w.varuint(uint64(nObjects))
	for i := range nObjects {
		w.hash(Hash{byte(i + 1)})
	}
	w.varuint(uint64(len(ops)))
	w.bytes(ops)
	return w.b
}

// huge is a count no data can hold; as an int it is negative.
var huge = binary.AppendUvarint(nil, math.MaxUint64)

func TestDecodeHugeCounts(t *testing.T) {
	cat := func(parts ...[]byte) []byte { return bytes.Join(parts, nil) }
	pad := make([]byte, 64)
	for name, data := range map[string][]byte{
		"path verbs":  object(nil, [][]byte{cat([]byte{0}, huge, pad)}, nil, 0, nil),
		"paint stops": object(nil, nil, [][]byte{cat([]byte{PaintLinear}, make([]byte, 16), huge, pad)}, 0, nil),
		"dash":        object(nil, nil, nil, 0, cat([]byte{OpDash}, huge, pad)),
		"glyph run":   object(nil, nil, nil, 0, cat([]byte{OpFillPathRun, NonZero}, huge, pad)),
		"mask bytes":  object(nil, nil, nil, 0, cat([]byte{OpMaskBegin, 0, 0, 0, 0, 0}, huge, pad)),
	} {
		t.Run(name, func(t *testing.T) {
			o, err := DecodeObject(data)
			if err == nil {
				_, err = o.Instructions()
			}
			var fe *FormatError
			if !errors.As(err, &fe) {
				t.Fatalf("err = %v", err)
			}
		})
	}
	t.Run("path collection", func(t *testing.T) {
		if _, err := DecodePathCollection(cat(huge, pad)); err == nil {
			t.Fatal("no error")
		}
	})
}

func TestExtractTextBadReferences(t *testing.T) {
	var ops buf
	ops.u8(OpFont)
	ops.bytes(huge) // a font that is not in the table
	ops.f32(12)
	ops.u8(OpFillText)
	ops.varuint(0)
	ops.f32s(0, 0, 0)
	ops.u8(OpUse)
	ops.varuint(7) // an object that is not in the table
	ops.u8(OpMark)
	ops.u8(MarkAltText)
	ops.varuint(0)
	ops.u8(OpUseAt)
	ops.bytes(huge)
	ops.f32s(0, 0)
	o, err := DecodeObject(object([]string{"text"}, nil, nil, 1, ops.b))
	if err != nil {
		t.Fatal(err)
	}
	runs, err := ExtractText(o, func(Hash) *ObjectPart { return nil })
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) != 2 || runs[0].Text != "text" || !runs[1].AltText {
		t.Fatalf("runs = %+v", runs)
	}
}

// nested builds objects that each draw the next one fan times, depth deep,
// the last one drawing a text; it returns the first and a resolver.
func nested(t *testing.T, depth, fan int) (*ObjectPart, func(Hash) *ObjectPart) {
	t.Helper()
	parts := map[Hash]*ObjectPart{}
	decode := func(o *Object) (*ObjectPart, Hash) {
		b := o.Encode()
		p, err := DecodeObject(b)
		if err != nil {
			t.Fatal(err)
		}
		h := HashOf(b)
		parts[h] = p
		return p, h
	}
	leaf := NewObject()
	leaf.FillText("x", 0, 0, 5)
	top, h := decode(leaf)
	for range depth {
		o := NewObject()
		ref := o.AddObject(h, Rect{})
		for range fan {
			o.Use(ref)
		}
		top, h = decode(o)
	}
	return top, func(h Hash) *ObjectPart { return parts[h] }
}

func TestExtractTextNesting(t *testing.T) {
	t.Run("shallow", func(t *testing.T) {
		top, resolve := nested(t, 8, 2)
		runs, err := ExtractText(top, resolve)
		if err != nil || len(runs) != 256 {
			t.Fatalf("%d runs, err = %v", len(runs), err)
		}
	})
	t.Run("deep", func(t *testing.T) {
		top, resolve := nested(t, 100_000, 1)
		var fe *FormatError
		if _, err := ExtractText(top, resolve); !errors.As(err, &fe) {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("fan out", func(t *testing.T) {
		// 4^30 texts, were they all drawn
		top, resolve := nested(t, 30, 4)
		start := time.Now()
		var fe *FormatError
		if _, err := ExtractText(top, resolve); !errors.As(err, &fe) {
			t.Fatalf("err = %v", err)
		}
		if d := time.Since(start); d > 30*time.Second {
			t.Fatalf("took %v", d)
		}
	})
}

func TestUnlockIterationBudget(t *testing.T) {
	// the second slot would take the iterations past what a reader accepts
	enc := &Encryption{Cipher: CipherA256GCM, Keys: []KeySlot{
		{Type: KeyPassword, KDF: KDFPBKDF2SHA256, Iter: 1, Salt: make([]byte, saltSize), Key: make([]byte, 40)},
		{Type: KeyPassword, KDF: KDFPBKDF2SHA256, Iter: MaxIterations, Salt: make([]byte, saltSize), Key: make([]byte, 40)},
	}}
	start := time.Now()
	_, err := unlockKey(enc, "password")
	var fe *FormatError
	if !errors.As(err, &fe) {
		t.Fatalf("err = %v", err)
	}
	if d := time.Since(start); d > time.Second {
		t.Fatalf("took %v", d)
	}
}

func TestManifestNumbers(t *testing.T) {
	obj := `"layers":[]`
	for name, views := range map[string]string{
		"null view":         `[null]`,
		"null page":         `[{"id":"v","kind":"fixed","pages":[null]}]`,
		"negative page":     `[{"id":"v","kind":"fixed","pages":[{"w":-5,"h":80000000,` + obj + `}]}]`,
		"huge page":         `[{"id":"v","kind":"fixed","pages":[{"w":100,"h":1e9,` + obj + `}]}]`,
		"page body":         `[{"id":"v","kind":"flow","pages":[{"w":100,"h":100,"body":{"x":0,"y":-1e30,"w":10,"h":10},` + obj + `}]}]`,
		"gap":               `[{"id":"v","kind":"flow","continuous":{"gap":1e30}}]`,
		"tiny tile":         `[{"id":"v","kind":"sheet","tile":0.000001}]`,
		"many rows":         `[{"id":"v","kind":"sheet","rows":[[16777216,10],[1,10]]}]`,
		"rows beyond f32":   `[{"id":"v","kind":"sheet","rows":[[1e11,0]]}]`,
		"half a column":     `[{"id":"v","kind":"sheet","cols":[[1.5,10]]}]`,
		"negative columns":  `[{"id":"v","kind":"sheet","cols":[[-1,10]]}]`,
		"thin columns":      `[{"id":"v","kind":"sheet","cols":[[400000,0.0000001]]}]`,
		"negative row size": `[{"id":"v","kind":"sheet","rows":[[4,-10]]}]`,
		"negative freeze":   `[{"id":"v","kind":"sheet","freeze":{"rows":-1}}]`,
	} {
		t.Run(name, func(t *testing.T) {
			mj := `{"bdf":1,"views":` + views + `,"parts":[]}`
			_, err := openBytes(append(header(0, HeaderSize, uint64(len(mj)), 0), mj...))
			var fe *FormatError
			if !errors.As(err, &fe) {
				t.Fatalf("err = %v", err)
			}
		})
	}
	// what writers write: pages, hidden rows, the tile of the default size
	mj := `{"bdf":1,"views":[{"id":"p","kind":"flow","continuous":{"gap":12},"pages":[{"w":595.28,"h":841.89,"body":{"x":72,"y":72,"w":451,"h":697},` + obj + `}]},` +
		`{"id":"s","kind":"sheet","tile":2048,"cols":[[3,64],[16381,48]],"rows":[[1,0],[1048575,15]],"freeze":{"rows":1}},{"id":"t","kind":"sheet"}],"parts":[]}`
	r, err := openBytes(append(header(0, HeaderSize, uint64(len(mj)), 0), mj...))
	if err != nil {
		t.Fatal(err)
	}
	if s := r.Manifest.Views[2].TileSize(); s != DefaultTile {
		t.Errorf("tile size %g", s)
	}
}
