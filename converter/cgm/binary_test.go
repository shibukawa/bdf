package cgm

import (
	"bytes"
	"encoding/binary"
	"runtime"
	"testing"
)

// TestManyPartitions reads an element whose parameters come in many small
// partitions: they are joined without copying what was joined before for
// every partition.
func TestManyPartitions(t *testing.T) {
	const n = 50_000
	var data, want []byte
	data = binary.BigEndian.AppendUint16(data, 4<<12|1<<5|31) // POLYLINE, a long form
	for i := range n {
		w := uint16(1)
		if i < n-1 {
			w |= 0x8000 // more partitions follow
		}
		data = binary.BigEndian.AppendUint16(data, w)
		data = append(data, byte(i), 0) // a byte and its padding
		want = append(want, byte(i))
	}
	file := bytes.Clone(data)
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	r := &binReader{data: data}
	e, ok := r.next()
	runtime.ReadMemStats(&after)
	if !ok || r.truncated || !bytes.Equal(e.data, want) {
		t.Fatalf("ok %v, truncated %v, %d bytes of parameters, want %d", ok, r.truncated, len(e.data), len(want))
	}
	if !bytes.Equal(data, file) {
		t.Error("the file was written to")
	}
	if got := after.TotalAlloc - before.TotalAlloc; got > 64*n {
		t.Errorf("%d bytes allocated for %d partitions of a byte", got, n)
	}
	if _, ok := r.next(); ok {
		t.Error("an element after the last one")
	}
}
