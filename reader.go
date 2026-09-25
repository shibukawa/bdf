package bdf

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Reader reads a document from either form.
type Reader struct {
	Manifest *Manifest
	entries  map[Hash]PartEntry
	load     func(PartEntry) ([]byte, error) // stored bytes
}

// OpenSingle parses the single-file form.
func OpenSingle(r io.ReaderAt, size int64) (*Reader, error) {
	hdr := make([]byte, HeaderSize)
	if _, err := r.ReadAt(hdr, 0); err != nil {
		return nil, err
	}
	rd := &reader{b: hdr}
	if !bytes.Equal(rd.bytes(4), Magic[:]) {
		return nil, &FormatError{Msg: "bad magic"}
	}
	if v := rd.u16(); v > FormatVersion {
		return nil, &FormatError{Msg: fmt.Sprintf("unsupported format version %d", v)}
	}
	rd.u16()
	moff := int64(rd.u64())
	mlen := int64(rd.u64())
	menc := rd.u8()
	if moff+mlen > size {
		return nil, &FormatError{Msg: "manifest out of range"}
	}
	mj := make([]byte, mlen)
	if _, err := r.ReadAt(mj, moff); err != nil {
		return nil, err
	}
	if menc == 1 {
		var err error
		if mj, err = inflate(mj); err != nil {
			return nil, err
		}
	}
	var m Manifest
	if err := json.Unmarshal(mj, &m); err != nil {
		return nil, err
	}
	base := moff + mlen
	return newReader(&m, func(e PartEntry) ([]byte, error) {
		if base+e.Off+int64(e.Len) > size {
			return nil, &FormatError{Msg: "part out of range"}
		}
		b := make([]byte, e.Len)
		_, err := r.ReadAt(b, base+e.Off)
		return b, err
	}), nil
}

// OpenSingleFile opens a single-form file by path.
func OpenSingleFile(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	st, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	return OpenSingle(f, st.Size())
}

// OpenSplit opens the split form from a directory.
func OpenSplit(dir string) (*Reader, error) {
	mj, err := os.ReadFile(filepath.Join(dir, "manifest.json"))
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(mj, &m); err != nil {
		return nil, err
	}
	return newReader(&m, func(e PartEntry) ([]byte, error) {
		return os.ReadFile(filepath.Join(dir, "parts", e.H.String()))
	}), nil
}

func newReader(m *Manifest, load func(PartEntry) ([]byte, error)) *Reader {
	r := &Reader{Manifest: m, entries: map[Hash]PartEntry{}, load: load}
	for _, e := range m.Parts {
		r.entries[e.H] = e
	}
	return r
}

// Entry returns the manifest entry for a part.
func (r *Reader) Entry(h Hash) (PartEntry, bool) {
	e, ok := r.entries[h]
	return e, ok
}

// Part returns the decoded bytes of a part.
func (r *Reader) Part(h Hash) ([]byte, error) {
	e, ok := r.entries[h]
	if !ok {
		return nil, fmt.Errorf("bdf: unknown part %s", h)
	}
	b, err := r.load(e)
	if err != nil {
		return nil, err
	}
	switch e.Enc {
	case EncIdentity, "":
		return b, nil
	case EncDeflateRaw:
		return inflate(b)
	default:
		return nil, fmt.Errorf("bdf: unknown encoding %q", e.Enc)
	}
}

// Object returns a decoded object part.
func (r *Reader) Object(h Hash) (*ObjectPart, error) {
	b, err := r.Part(h)
	if err != nil {
		return nil, err
	}
	return DecodeObject(b)
}

// ToDocument rebuilds a Document (for split/join round trips).
func (r *Reader) ToDocument() (*Document, error) {
	d := NewDocument()
	d.Meta = r.Manifest.Meta
	d.Views = r.Manifest.Views
	for _, e := range r.Manifest.Parts {
		b, err := r.Part(e.H)
		if err != nil {
			return nil, err
		}
		if got := HashOf(b); got != e.H {
			return nil, fmt.Errorf("bdf: part %s hash mismatch (%s)", e.H, got)
		}
		d.AddPart(e.T, b)
	}
	return d, nil
}

func inflate(b []byte) ([]byte, error) {
	return io.ReadAll(flate.NewReader(bytes.NewReader(b)))
}
