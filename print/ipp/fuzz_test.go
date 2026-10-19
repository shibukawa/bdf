package ipp

import (
	"bytes"
	"errors"
	"testing"
)

// FuzzDecode checks that any bytes decode to a message that encodes
// back to what decodes the same, or are refused as malformed or past a
// limit, without a panic.
func FuzzDecode(f *testing.F) {
	s, _ := sample().Marshal()
	f.Add(s)
	b, _ := NewRequest(OpPrintJob, 1, "", "").Marshal()
	f.Add(b)
	f.Add([]byte{2, 0, 0, 1, 0, 0, 0, 1, 0x01, 0x34, 0, 1, 'c', 0, 0, 0x4a, 0, 0, 0, 1, 'm', 0x21, 0, 0, 0, 4, 0, 0, 0, 1, 0x37, 0, 0, 0, 0, 3})
	f.Add([]byte{})
	lim := &Limits{MaxBytes: 1 << 16, MaxAttributes: 200, MaxValueBytes: 4096, MaxValues: 50, MaxDepth: 4}
	f.Fuzz(func(t *testing.T, in []byte) {
		m, err := Unmarshal(in, lim)
		if err != nil {
			if !errors.Is(err, ErrMalformed) && !errors.Is(err, ErrLimit) {
				t.Fatalf("Unmarshal: an error that is not malformed nor a limit: %v", err)
			}
			return
		}
		out, err := m.Marshal()
		if err != nil {
			t.Fatalf("Marshal of a decoded message: %v", err)
		}
		again, err := Unmarshal(out, lim)
		if err != nil {
			t.Fatalf("Unmarshal of an encoding: %v (%x)", err, out)
		}
		out2, err := again.Marshal()
		if err != nil {
			t.Fatalf("Marshal again: %v", err)
		}
		if !bytes.Equal(out, out2) {
			t.Fatalf("the encoding is not stable:\n%x\n%x", out, out2)
		}
	})
}
