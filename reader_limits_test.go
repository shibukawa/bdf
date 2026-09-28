package bdf

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// manifestHeaderReader exposes only a header. Its declared file size is much
// larger, so the reader must reject the header before requesting the payload.
type manifestHeaderReader struct {
	header []byte
	reads  int
}

func (r *manifestHeaderReader) ReadAt(b []byte, off int64) (int, error) {
	r.reads++
	if off == 0 && len(b) == HeaderSize {
		return copy(b, r.header), nil
	}
	return 0, io.ErrUnexpectedEOF
}

func TestOpenSingleStoredManifestLimit(t *testing.T) {
	for _, encoding := range []byte{0, 1} {
		r := &manifestHeaderReader{header: header(0, HeaderSize, MaxManifestSize+1, encoding)}
		_, err := OpenSingle(r, HeaderSize+MaxManifestSize+1)
		var fe *FormatError
		if !errors.As(err, &fe) || r.reads != 1 {
			t.Errorf("encoding %d: error = %v, reads = %d; want format error before payload read", encoding, err, r.reads)
		}
	}
}

func TestOpenSplitStoredManifestLimit(t *testing.T) {
	dir := t.TempDir()
	f, err := os.Create(filepath.Join(dir, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	// A sparse file exercises the size check without allocating its contents.
	err = f.Truncate(MaxManifestSize + 1)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	var fe *FormatError
	if _, err := OpenSplit(dir); !errors.As(err, &fe) {
		t.Fatalf("error = %v; want format error", err)
	}
}

func TestOpenSplitStoredPartLength(t *testing.T) {
	for _, tc := range []struct {
		name   string
		stored []byte
		length int
	}{
		{"long", []byte("contents"), 1},
		{"short", []byte("contents"), 20},
		{"negative", []byte("contents"), -1},
		{"empty", nil, 0},
		{"exact", []byte("contents"), 8},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := NewDocument()
			h := d.AddPart(PartImage, tc.stored)
			dir := t.TempDir()
			if err := d.WriteSplit(dir); err != nil {
				t.Fatal(err)
			}
			r, err := OpenSplit(dir)
			if err != nil {
				t.Fatal(err)
			}
			e := r.entries[h]
			e.Len = tc.length
			r.entries[h] = e
			b, err := r.Part(h)
			if tc.length == len(tc.stored) {
				if err != nil || !bytes.Equal(b, tc.stored) {
					t.Fatalf("part = %q, error = %v", b, err)
				}
			} else {
				var fe *FormatError
				if !errors.As(err, &fe) {
					t.Fatalf("error = %v; want format error", err)
				}
			}
		})
	}
}

func TestUnlockStoredManifestLimit(t *testing.T) {
	d := NewDocument()
	var err error
	d.Lock, err = NewPasswordLock("password", 1)
	if err != nil {
		t.Fatal(err)
	}
	m, _, _, err := d.stored()
	if err != nil {
		t.Fatal(err)
	}
	for i := range m.Parts {
		if m.Parts[i].H == m.Encryption.Manifest.Part {
			m.Parts[i].Len = MaxManifestSize + nonceSize + 16 + 1
		}
	}
	loaded := false
	r := newReader(m, func(PartEntry) ([]byte, error) {
		loaded = true
		return nil, io.ErrUnexpectedEOF
	})
	var fe *FormatError
	if err := r.Unlock("password"); !errors.As(err, &fe) || loaded {
		t.Fatalf("error = %v, loaded = %v; want format error before load", err, loaded)
	}
}
