package parquet

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/bits"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/klauspost/compress/snappy"
	kzstd "github.com/klauspost/compress/zstd"
)

func TestThrift(t *testing.T) {
	e := newTenc()
	e.i32(1, -5)
	e.bin(2, "name")
	e.i64(300, 1<<40) // a long field header
	e.boolean(301, true)
	e.begin(302) // skipped: a struct holding a list of lists and a map
	e.list(1, tList, 2)
	e.b = append(e.b, 0x25, 2, 4, 0x25, 6, 8)
	e.field(2, tMap)
	e.b = append(e.b, 1, 0x85, 1, 'k', 2)
	e.end()
	e.list(303, tBinary, 16)
	for range 16 {
		e.b = append(e.b, 1, 'x')
	}
	b := e.bytes()

	t1 := &thrift{b: b}
	var got []any
	t1.fields(func(id int16, typ byte) {
		switch id {
		case 1:
			got = append(got, t1.i32())
		case 2:
			got = append(got, t1.string())
		case 300:
			got = append(got, t1.i64())
		case 301:
			got = append(got, typ == tTrue)
		case 303:
			got = append(got, len(t1.strings(typ)))
		default:
			t1.skip(typ)
		}
	})
	if t1.err != nil || t1.pos != len(b) || !slices.Equal(got, []any{int32(-5), "name", int64(1 << 40), true, 16}) {
		t.Errorf("err %v, pos %d of %d, got %v", t1.err, t1.pos, len(b), got)
	}
	for i := range len(b) - 1 {
		t2 := &thrift{b: b[:i]}
		t2.fields(func(_ int16, typ byte) { t2.skip(typ) })
		if t2.err == nil {
			t.Errorf("cut at %d: no error", i)
		}
	}
	// nesting beyond the bound
	deep := bytes.Repeat([]byte{0x1c}, 100)
	t3 := &thrift{b: deep}
	t3.fields(func(_ int16, typ byte) { t3.skip(typ) })
	if t3.err == nil {
		t.Error("deep nesting: no error")
	}
}

func TestPageHeader(t *testing.T) {
	e := newTenc()
	e.i32(1, pageDataV2)
	e.i32(2, 100)
	e.i32(3, 60)
	e.i32(4, 1234) // crc, skipped
	e.begin(8)
	e.i32(1, 10)
	e.i32(2, 1)
	e.i32(3, 10)
	e.i32(4, encRLEDict)
	e.i32(5, 3)
	e.i32(6, 2)
	e.boolean(7, false)
	e.begin(8) // statistics, skipped
	e.bin(5, "max")
	e.end()
	e.end()
	b := append(e.bytes(), "the page"...)
	h, n, err := readPageHeader(b)
	if err != nil || string(b[n:]) != "the page" {
		t.Fatalf("err %v, header of %d bytes", err, n)
	}
	want := pageHeader{typ: pageDataV2, uncompressed: 100, compressed: 60, numValues: 10, encoding: encRLEDict, defLength: 3, repLength: 2, uncompressV2: true}
	if *h != want {
		t.Errorf("header %+v", *h)
	}
}

func collect(t *testing.T, d interface{ next() (uint64, error) }, n int) []uint64 {
	t.Helper()
	var out []uint64
	for range n {
		v, err := d.next()
		if err != nil {
			t.Fatalf("after %v: %v", out, err)
		}
		out = append(out, v)
	}
	return out
}

func TestHybrid(t *testing.T) {
	// the examples of Encodings.md: 0 to 7 bit-packed in 3 bits, and a run
	// of five 3s
	b := []byte{0x03, 0x88, 0xc6, 0xfa, 0x0a, 0x03}
	got := collect(t, newHybrid(b, 3), 13)
	if !slices.Equal(got, []uint64{0, 1, 2, 3, 4, 5, 6, 7, 3, 3, 3, 3, 3}) {
		t.Errorf("got %v", got)
	}
	if _, err := newHybrid(b, 3).next(); err != nil {
		t.Error(err)
	}
	h := newHybrid(b, 3)
	collect(t, h, 13)
	if _, err := h.next(); err == nil {
		t.Error("past the end: no error")
	}
	// BIT_PACKED, most significant bit first
	got = collect(t, &bitPacked{b: []byte{0x05, 0x39, 0x77}, width: 3}, 8)
	if !slices.Equal(got, []uint64{0, 1, 2, 3, 4, 5, 6, 7}) {
		t.Errorf("bit-packed: %v", got)
	}
	// widths up to 64 bits
	for _, w := range []uint{1, 7, 13, 31, 32, 33, 57, 63, 64} {
		vals := make([]uint64, 16)
		var bb []byte
		acc, n := uint64(0), uint(0)
		for i := range vals {
			vals[i] = rand.Uint64() & mask(w)
			for k := range w {
				acc |= (vals[i] >> k & 1) << n
				if n++; n == 64 {
					bb = binary.LittleEndian.AppendUint64(bb, acc)
					acc, n = 0, 0
				}
			}
		}
		for ; n > 0; n -= min(n, 8) {
			bb = append(bb, byte(acc))
			acc >>= 8
		}
		for i, v := range vals {
			if got := bitsLE(bb, uint64(i)*uint64(w), w); got != v {
				t.Errorf("width %d, value %d: %x, want %x", w, i, got, v)
			}
		}
	}
}

// encodeDelta writes DELTA_BINARY_PACKED values (blocks of 128 in four
// miniblocks).
func encodeDelta(vals []int64, i32 bool) []byte {
	out := binary.AppendUvarint(nil, 128)
	out = binary.AppendUvarint(out, 4)
	out = binary.AppendUvarint(out, uint64(len(vals)))
	if len(vals) == 0 {
		return binary.AppendVarint(out, 0)
	}
	out = binary.AppendVarint(out, vals[0])
	for start := 1; start < len(vals); start += 128 {
		block := vals[start:min(start+128, len(vals))]
		deltas := make([]int64, len(block))
		minD := int64(math.MaxInt64)
		for i, v := range block {
			if i32 {
				deltas[i] = int64(int32(v) - int32(vals[start+i-1]))
			} else {
				deltas[i] = v - vals[start+i-1]
			}
			minD = min(minD, deltas[i])
		}
		out = binary.AppendVarint(out, minD)
		widths := make([]byte, 4)
		for m := range 4 {
			for _, d := range deltas[min(m*32, len(deltas)):min(m*32+32, len(deltas))] {
				widths[m] = max(widths[m], byte(bits.Len64(uint64(d-minD))))
			}
		}
		out = append(out, widths...)
		for m := range 4 {
			if m*32 >= len(deltas) {
				break
			}
			w := uint(widths[m])
			mini := make([]byte, 32*w/8)
			for i := range 32 {
				var d uint64
				if m*32+i < len(deltas) {
					d = uint64(deltas[m*32+i] - minD)
				}
				for k := range w {
					if d>>k&1 != 0 {
						p := uint(i)*w + k
						mini[p/8] |= 1 << (p % 8)
					}
				}
			}
			out = append(out, mini...)
		}
	}
	return out
}

func TestDeltaBinary(t *testing.T) {
	// the first example of Encodings.md
	b := []byte{0x80, 0x01, 0x04, 0x05, 0x02, 0x02, 0, 0, 0, 0}
	d, err := newDeltaBinary(b, false)
	if err != nil {
		t.Fatal(err)
	}
	if got := collect(t, d.deltaInts(), 5); !slices.Equal(got, []uint64{1, 2, 3, 4, 5}) {
		t.Errorf("got %v", got)
	}
	if end, err := deltaEnd(b); err != nil || end != len(b) {
		t.Errorf("end %d, %v", end, err)
	}
	for _, c := range []struct {
		vals []int64
		i32  bool
	}{
		{[]int64{7, 5, 3, 1, 2, 3, 4, 5}, false},
		{[]int64{math.MinInt64, math.MaxInt64, 0, -1, math.MaxInt64, math.MinInt64}, false},
		{[]int64{math.MinInt32, math.MaxInt32, 0, -1, math.MaxInt32}, true},
		{func() []int64 {
			v := make([]int64, 1000)
			for i := range v {
				v[i] = rand.Int64N(1<<40) - 1<<39
			}
			return v
		}(), false},
		{[]int64{42}, false},
		{nil, false},
	} {
		b := encodeDelta(c.vals, c.i32)
		b = append(b, "tail"...)
		end, err := deltaEnd(b)
		if err != nil || string(b[end:]) != "tail" {
			t.Errorf("%d values: end %d of %d, %v", len(c.vals), end, len(b), err)
		}
		d, err := newDeltaBinary(b[:end], c.i32)
		if err != nil {
			t.Fatal(err)
		}
		for i, want := range c.vals {
			v, err := d.nextInt()
			if err != nil || int64(v) != want {
				t.Fatalf("%d values, value %d: %d (%v), want %d", len(c.vals), i, int64(v), err, want)
			}
		}
		if _, err := d.nextInt(); err == nil {
			t.Errorf("%d values: a value past the end", len(c.vals))
		}
	}
}

type deltaInts struct{ d *deltaBinary }

func (d deltaInts) next() (uint64, error)   { return d.d.nextInt() }
func (d *deltaBinary) deltaInts() deltaInts { return deltaInts{d} }

func TestDeltaByteArrays(t *testing.T) {
	words := []string{"apple", "applesauce", "apply", "", "banana", "band"}
	lengths := make([]int64, len(words))
	var data []byte
	for i, w := range words {
		lengths[i] = int64(len(w))
		data = append(data, w...)
	}
	d, _, err := newDeltaLength(append(encodeDelta(lengths, true), data...))
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range words {
		if v, err := d.next(); err != nil || string(v.b) != w {
			t.Fatalf("%q (%v), want %q", v.b, err, w)
		}
	}
	// the prefixes shared with the value before, then the rest
	prefixes := []int64{0, 5, 4, 0, 0, 3}
	var suffixes, sufData []byte
	sufLengths := make([]int64, len(words))
	for i, w := range words {
		sufLengths[i] = int64(len(w)) - prefixes[i]
		sufData = append(sufData, w[prefixes[i]:]...)
	}
	suffixes = append(encodeDelta(sufLengths, true), sufData...)
	da, err := newDeltaByteArray(append(encodeDelta(prefixes, true), suffixes...))
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range words {
		if v, err := da.next(); err != nil || string(v.b) != w {
			t.Fatalf("%q (%v), want %q", v.b, err, w)
		}
	}
	// a prefix longer than the value before
	da, _ = newDeltaByteArray(append(encodeDelta([]int64{3}, true), append(encodeDelta([]int64{1}, true), 'x')...))
	if _, err := da.next(); err == nil {
		t.Error("prefix too long: no error")
	}
}

func TestByteStreamSplit(t *testing.T) {
	vals := []float32{1.5, -2.25, 1e30}
	b := make([]byte, 12)
	for i, v := range vals {
		for j := range 4 {
			b[j*3+i] = byte(math.Float32bits(v) >> (8 * j))
		}
	}
	d, err := newByteStreamSplit(b, typeFloat, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range vals {
		if v, err := d.next(); err != nil || math.Float32frombits(uint32(v.u)) != want {
			t.Errorf("%v (%v), want %v", math.Float32frombits(uint32(v.u)), err, want)
		}
	}
	if _, err := newByteStreamSplit(b[:11], typeFloat, 0); err == nil {
		t.Error("a partial value: no error")
	}
}

func TestPlain(t *testing.T) {
	b := binary.LittleEndian.AppendUint32(nil, 3)
	b = append(b, "abc"...)
	b = binary.LittleEndian.AppendUint32(b, 10) // longer than what is left
	d := &plainDecoder{b: b, phys: typeBinary}
	if v, err := d.next(); err != nil || string(v.b) != "abc" {
		t.Errorf("%q, %v", v.b, err)
	}
	if _, err := d.next(); err == nil {
		t.Error("truncated: no error")
	}
	d = &plainDecoder{b: []byte{0b101}, phys: typeBoolean}
	got := collect(t, plainBits{d}, 3)
	if !slices.Equal(got, []uint64{1, 0, 1}) {
		t.Errorf("booleans %v", got)
	}
	i96 := append(binary.LittleEndian.AppendUint64(nil, 5), 0x8c, 0x3d, 0x25, 0)
	if v, err := (&plainDecoder{b: i96, phys: typeInt96}).next(); err != nil || v.u != 5 || v.x != 2440588 {
		t.Errorf("INT96 %+v, %v", v, err)
	}
	if v, _ := (&plainDecoder{b: []byte{0xff, 0xff, 0xff, 0xff}, phys: typeInt32}).next(); int64(v.u) != -1 {
		t.Errorf("INT32 -1: %x", v.u)
	}
}

type plainBits struct{ d *plainDecoder }

func (p plainBits) next() (uint64, error) {
	v, err := p.d.next()
	return v.u, err
}

func TestSnappy(t *testing.T) {
	for _, n := range []int{0, 1, 59, 60, 61, 300, 70000} {
		data := make([]byte, n)
		for i := range data {
			data[i] = "abcabcabd"[i%9] + byte(i/1000)
		}
		enc := snappy.Encode(nil, data)
		got, err := decompress(codecSnappy, enc, n)
		if err != nil || !bytes.Equal(got, data) {
			t.Errorf("%d bytes: %v", n, err)
		}
		if _, err := decompress(codecSnappy, enc, n+1); err == nil {
			t.Errorf("%d bytes: a wrong size decompressed", n)
		}
	}
	// a copy that reaches before the start
	if err := snappyBlock(make([]byte, 5), []byte{5, 0x01 | 1<<2, 9}); err == nil {
		t.Error("bad offset: no error")
	}
}

func TestZstdDecompress(t *testing.T) {
	enc, err := kzstd.NewWriter(nil, kzstd.WithEncoderConcurrency(1))
	if err != nil {
		t.Fatal(err)
	}
	defer enc.Close()

	want := []byte("tinygodriver public Reader")
	frame := enc.EncodeAll(want[:10], nil)
	frame = enc.EncodeAll(want[10:], frame)
	got, err := decompress(codecZstd, frame, len(want))
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("decoded %q, %v; want %q", got, err, want)
	}
	if _, err := decompress(codecZstd, frame, len(want)-1); err == nil {
		t.Error("accepted output larger than the declared page size")
	}
	if _, err := decompress(codecZstd, frame, len(want)+1); err == nil {
		t.Error("accepted output smaller than the declared page size")
	}
}

func TestLZ4(t *testing.T) {
	// "abc", then a copy of 9 bytes 3 back, then "xyz"
	block := []byte{0x35, 'a', 'b', 'c', 3, 0, 0x30, 'x', 'y', 'z'}
	want := "abcabcabcabcxyz"
	got, err := decompress(codecLZ4Raw, block, len(want))
	if err != nil || string(got) != want {
		t.Fatalf("block: %q, %v", got, err)
	}
	// Hadoop's framing, in two blocks
	var hadoop []byte
	for _, part := range [][]byte{block, {0x20, '!', '?'}} {
		n := len(want)
		if part[1] == '!' {
			n = 2
		}
		hadoop = binary.BigEndian.AppendUint32(hadoop, uint32(n))
		hadoop = binary.BigEndian.AppendUint32(hadoop, uint32(len(part)))
		hadoop = append(hadoop, part...)
	}
	if got, err := decompress(codecLZ4, hadoop, len(want)+2); err != nil || string(got) != want+"!?" {
		t.Errorf("Hadoop: %q, %v", got, err)
	}
	// the frame format: a block, a block copying from the one before (4
	// bytes, 6 back), and a block stored raw
	frame := []byte{0x04, 0x22, 0x4d, 0x18, 0x40, 0x40, 0x00}
	for _, b := range [][]byte{block, {0x00, 6, 0}} {
		frame = binary.LittleEndian.AppendUint32(frame, uint32(len(b)))
		frame = append(frame, b...)
	}
	frame = binary.LittleEndian.AppendUint32(frame, 0x80000002)
	frame = append(frame, "ok"...)
	frame = binary.LittleEndian.AppendUint32(frame, 0)
	wantFrame := want + "abcx" + "ok"
	if got, err := decompress(codecLZ4, frame, len(wantFrame)); err != nil || string(got) != wantFrame {
		t.Errorf("frame: %q, %v", got, err)
	}
	// a raw block under the old codec, as old parquet-cpp wrote it
	if got, err := decompress(codecLZ4, block, len(want)); err != nil || string(got) != want {
		t.Errorf("raw block: %q, %v", got, err)
	}
	for _, bad := range [][]byte{{0x35, 'a', 'b', 'c', 4, 0}, {0x35, 'a', 'b', 'c', 0, 0}, {0xf0}} {
		if _, err := decompress(codecLZ4Raw, bad, 12); err == nil {
			t.Errorf("%x: no error", bad)
		}
	}
}

func TestDecompressSizes(t *testing.T) {
	// sizes said that the data cannot hold are refused before anything
	// that large is made
	for _, codec := range []int32{codecSnappy, codecLZ4Raw, codecLZ4, codecGzip, codecBrotli, codecZstd} {
		if _, err := decompress(codec, []byte{0x80, 0x80, 0x80, 0x80, 0x01}, maxPageSize); err == nil {
			t.Errorf("codec %d: no error", codec)
		}
	}
	if _, err := decompress(codecNone, nil, maxPageSize+1); err == nil {
		t.Error("too large: no error")
	}
	if _, err := decompress(codecLZO, nil, 1); err == nil {
		t.Error("LZO: no error")
	}
}
