package bdf

import (
	"bytes"
	"compress/flate"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
)

// MaxManifestSize bounds the manifest JSON a reader inflates: a manifest
// states no size of its own, so a small file could otherwise ask for any
// amount of memory.
const MaxManifestSize = 256 << 20

// Reader reads a document from either form.
//
// An encrypted document opens locked: Manifest is then the outer manifest,
// whose parts are sealed, until Unlock replaces it with the document's own.
type Reader struct {
	Manifest *Manifest
	stored   *Manifest // as stored: the outer manifest of an encrypted document
	entries  map[Hash]PartEntry
	outer    map[Hash]PartEntry              // stored parts of an encrypted document
	load     func(PartEntry) ([]byte, error) // stored bytes
	aead     cipher.AEAD                     // set by Unlock
	file     io.Closer                       // the file OpenSingleFile opened
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
	flags := rd.u16()
	moff := int64(rd.u64())
	mlen := int64(rd.u64())
	menc := rd.u8()
	if !inRange(moff, mlen, size) {
		return nil, &FormatError{Msg: "manifest out of range"}
	}
	mj := make([]byte, mlen)
	if _, err := r.ReadAt(mj, moff); err != nil {
		return nil, err
	}
	if menc == 1 {
		var err error
		if mj, err = inflate(mj, MaxManifestSize, 0); err != nil {
			return nil, err
		}
	}
	var m Manifest
	if err := readManifest(mj, &m); err != nil {
		return nil, err
	}
	if (flags&FlagEncrypted != 0) != (m.Encryption != nil) {
		return nil, &FormatError{Msg: "encryption flag does not match the manifest"}
	}
	base := moff + mlen
	return newReader(&m, func(e PartEntry) ([]byte, error) {
		if !inRange(e.Off, int64(e.Len), size-base) {
			return nil, &FormatError{Msg: "part out of range"}
		}
		b := make([]byte, e.Len)
		_, err := r.ReadAt(b, base+e.Off)
		return b, err
	}), nil
}

// OpenSingleFile opens a single-form file by path. The parts are read from
// the file as they are asked for, so it stays open until Close.
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
	r, err := OpenSingle(f, st.Size())
	if err != nil {
		f.Close()
		return nil, err
	}
	r.file = f
	return r, nil
}

// Close closes the file of a Reader from OpenSingleFile; the parts cannot be
// read after it. Readers opened in the other ways hold nothing to close.
func (r *Reader) Close() error {
	if r.file == nil {
		return nil
	}
	f := r.file
	r.file = nil
	return f.Close()
}

// OpenSplit opens the split form from a directory.
func OpenSplit(dir string) (*Reader, error) {
	mj, err := os.ReadFile(filepath.Join(dir, "manifest.json"))
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := readManifest(mj, &m); err != nil {
		return nil, err
	}
	return newReader(&m, func(e PartEntry) ([]byte, error) {
		return os.ReadFile(filepath.Join(dir, "parts", e.H.String()))
	}), nil
}

// readManifest parses the JSON of a manifest and checks its numbers.
func readManifest(b []byte, m *Manifest) error {
	if err := json.Unmarshal(b, m); err != nil {
		return err
	}
	return m.check()
}

// inRange reports whether n bytes at off lie within size bytes; the values
// come from the file, so they may be negative or overflow when added.
func inRange(off, n, size int64) bool {
	return off >= 0 && n >= 0 && off <= size && n <= size-off
}

func newReader(m *Manifest, load func(PartEntry) ([]byte, error)) *Reader {
	r := &Reader{Manifest: m, stored: m, load: load}
	r.entries = entries(m)
	return r
}

func entries(m *Manifest) map[Hash]PartEntry {
	out := make(map[Hash]PartEntry, len(m.Parts))
	for _, e := range m.Parts {
		out[e.H] = e
	}
	return out
}

// Encrypted reports whether the document is encrypted.
func (r *Reader) Encrypted() bool { return r.stored.Encryption != nil }

// Locked reports whether the document is encrypted and not unlocked yet.
func (r *Reader) Locked() bool { return r.Encrypted() && r.aead == nil }

// Unlock opens an encrypted document with its password: it decrypts the
// manifest, which becomes Manifest, and lets Part decrypt parts. It returns
// ErrWrongPassword when no key slot opens with the password.
func (r *Reader) Unlock(password string) error {
	enc := r.stored.Encryption
	if enc == nil {
		return errors.New("bdf: the document is not encrypted")
	}
	key, err := unlockKey(enc, password)
	if err != nil {
		return err
	}
	a, err := newAEAD(key)
	if err != nil {
		return err
	}
	outer := entries(r.stored)
	e, ok := outer[enc.Manifest.Part]
	if !ok {
		return &FormatError{Msg: "sealed manifest part missing"}
	}
	b, err := r.load(e)
	if err != nil {
		return err
	}
	if b, err = open(a, b, manifestAAD); err != nil {
		return err
	}
	if b, err = decodeStored(b, enc.Manifest.Enc, MaxManifestSize, 0); err != nil {
		return err
	}
	var m Manifest
	if err := readManifest(b, &m); err != nil {
		return err
	}
	if m.Encryption != nil {
		return &FormatError{Msg: "sealed manifest is encrypted again"}
	}
	for _, e := range m.Parts {
		if _, ok := outer[e.Sealed]; !ok {
			return &FormatError{Msg: fmt.Sprintf("part %s: sealed part %s missing", e.H, e.Sealed)}
		}
	}
	r.Manifest, r.entries, r.outer, r.aead = &m, entries(&m), outer, a
	return nil
}

// Entry returns the manifest entry for a part.
func (r *Reader) Entry(h Hash) (PartEntry, bool) {
	e, ok := r.entries[h]
	return e, ok
}

// Part returns the decoded bytes of a part. The parts of a locked document
// are its sealed parts, returned as they are stored.
func (r *Reader) Part(h Hash) ([]byte, error) {
	e, ok := r.entries[h]
	if !ok {
		return nil, fmt.Errorf("bdf: unknown part %s", h)
	}
	var b []byte
	var err error
	if r.aead != nil {
		if b, err = r.load(r.outer[e.Sealed]); err != nil {
			return nil, err
		}
		if b, err = open(r.aead, b, e.H[:]); err != nil {
			return nil, fmt.Errorf("part %s: %w", e.H, err)
		}
	} else if b, err = r.load(e); err != nil {
		return nil, err
	}
	return decodeStored(b, e.Enc, e.Size, e.Size)
}

// decodeStored decodes stored bytes, which decode to limit bytes at most
// and are expected to decode to size bytes (0: not known).
func decodeStored(b []byte, enc string, limit, size int) ([]byte, error) {
	switch enc {
	case EncIdentity, "":
		return b, nil
	case EncDeflateRaw:
		return inflate(b, limit, size)
	default:
		return nil, fmt.Errorf("bdf: unknown encoding %q", enc)
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

// ToDocument rebuilds a Document (for split/join round trips). An encrypted
// document must be unlocked first, and the Document it gives is not
// encrypted until a Lock is set; WriteSingle and WriteSplit of the Reader
// copy an encrypted document without its password.
func (r *Reader) ToDocument() (*Document, error) {
	if r.Locked() {
		return nil, ErrLocked
	}
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

// WriteSingle writes the document in the single-file form. The stored parts
// are copied as they are (after checking their hashes), so an encrypted
// document needs no password.
func (r *Reader) WriteSingle(w io.Writer) error {
	m, data, err := r.storedParts()
	if err != nil {
		return err
	}
	var off int64
	for i := range m.Parts {
		m.Parts[i].Off = off
		off += int64(m.Parts[i].Len)
	}
	flags := uint16(0)
	if m.Encryption != nil {
		flags = FlagEncrypted
	}
	return writeSingle(w, m, data, flags, flate.BestCompression, defaultMinCompress)
}

// WriteSplit writes the document in the split form into dir, copying the
// stored parts as WriteSingle does.
func (r *Reader) WriteSplit(dir string) error {
	m, data, err := r.storedParts()
	if err != nil {
		return err
	}
	return writeSplit(dir, m, data)
}

// storedParts reads every stored part and checks it against its name: the
// hash of the decoded bytes, or of the stored bytes for a sealed part.
func (r *Reader) storedParts() (*Manifest, [][]byte, error) {
	m := *r.stored
	m.Parts = slices.Clone(m.Parts)
	data := make([][]byte, len(m.Parts))
	for i, e := range m.Parts {
		b, err := r.load(e)
		if err != nil {
			return nil, nil, err
		}
		if len(b) != e.Len {
			return nil, nil, fmt.Errorf("bdf: part %s is %d bytes, not %d", e.H, len(b), e.Len)
		}
		dec := b
		if e.T != PartSealed {
			if dec, err = decodeStored(b, e.Enc, e.Size, e.Size); err != nil {
				return nil, nil, err
			}
		}
		if got := HashOf(dec); got != e.H {
			return nil, nil, fmt.Errorf("bdf: part %s hash mismatch (%s)", e.H, got)
		}
		data[i] = b
	}
	return &m, data, nil
}

// inflate decodes deflate-raw data that holds limit bytes at most (for a
// part, the size the manifest states). Data that goes on past the limit is an
// error, so a small file cannot ask for more memory than its manifest says it
// needs. size is the size to expect, allocated at once (0: not known).
func inflate(b []byte, limit, size int) ([]byte, error) {
	if limit < 0 || size < 0 {
		return nil, &FormatError{Msg: "negative part size"}
	}
	var out bytes.Buffer
	// deflate packs 1032 bytes into one at best, so no more than that is
	// allocated on the word of the manifest
	out.Grow(int(min(int64(size), int64(limit), 1032*int64(len(b)))))
	n, err := io.Copy(&out, io.LimitReader(flate.NewReader(bytes.NewReader(b)), int64(limit)+1))
	if err != nil {
		return nil, err
	}
	if n > int64(limit) {
		return nil, &FormatError{Msg: "data inflates to more than its stated size"}
	}
	return out.Bytes(), nil
}
