package bdf

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
)

// Part is a raw (uncompressed) part.
type Part struct {
	Hash Hash
	Type string
	Data []byte
}

// defaultMinCompress is the default of Document.MinCompress.
const defaultMinCompress = 512

// Document accumulates parts and views and writes them out.
type Document struct {
	Meta  Meta
	Views []*View

	// CompressionLevel is the flate level for deflate-raw parts (default: flate.BestCompression).
	CompressionLevel int
	// MinCompress is the smallest decoded size that gets compressed (default 512).
	MinCompress int
	// Lock, when set, makes WriteSingle and WriteSplit encrypt the document
	// (docs/spec.md §3.5).
	Lock *Lock

	parts map[Hash]*Part
	order []Hash
}

// NewDocument creates an empty document.
func NewDocument() *Document {
	return &Document{
		CompressionLevel: flate.BestCompression,
		MinCompress:      defaultMinCompress,
		parts:            map[Hash]*Part{},
	}
}

// AddPart stores data as a part of the given type and returns its hash.
// Identical content is stored once.
func (d *Document) AddPart(typ string, data []byte) Hash {
	h := HashOf(data)
	if _, ok := d.parts[h]; !ok {
		d.parts[h] = &Part{Hash: h, Type: typ, Data: data}
		d.order = append(d.order, h)
	}
	return h
}

// AddFont stores a font file (WOFF2/TTF/OTF).
func (d *Document) AddFont(data []byte) Hash { return d.AddPart(PartFont, data) }

// AddImage stores an image file (PNG/JPEG/WebP/AVIF/SVG).
func (d *Document) AddImage(data []byte) Hash { return d.AddPart(PartImage, data) }

// AddPaths stores a Path collection part.
func (d *Document) AddPaths(paths []*Path) Hash {
	return d.AddPart(PartPath, EncodePathCollection(paths))
}

// AddObject encodes and stores an object, returning its hash and bbox.
func (d *Document) AddObject(o *Object) (Hash, Rect) {
	for _, dep := range o.Deps() {
		if _, ok := d.parts[dep]; !ok {
			panic(fmt.Sprintf("bdf: object references unknown part %s", dep))
		}
	}
	return d.AddPart(PartObject, o.Encode()), o.BBox
}

// NewView appends a view.
func (d *Document) NewView(id, kind, title string) *View {
	v := &View{ID: id, Kind: kind, Title: title}
	d.Views = append(d.Views, v)
	return v
}

// Part returns a stored part.
func (d *Document) Part(h Hash) *Part { return d.parts[h] }

// Parts returns parts in insertion order.
func (d *Document) Parts() []*Part {
	out := make([]*Part, len(d.order))
	for i, h := range d.order {
		out[i] = d.parts[h]
	}
	return out
}

func (d *Document) shouldCompress(p *Part) bool {
	switch p.Type {
	case PartImage:
		// Encoded images are compressed already, but not SVG (text) or the
		// pixels of BMP and of the BMP images in icons.
		if !uncompressedImage(p.Data) {
			return false
		}
	case PartFont:
		// WOFF and WOFF2 are compressed already; TTF/OTF are not.
		if len(p.Data) >= 4 && (string(p.Data[:4]) == "wOF2" || string(p.Data[:4]) == "wOFF") {
			return false
		}
	}
	return len(p.Data) >= d.MinCompress
}

// uncompressedImage reports whether an image part is SVG, BMP or an icon,
// which compress well (spec §3.2).
func uncompressedImage(b []byte) bool {
	switch {
	case len(b) >= 2 && b[0] == 'B' && b[1] == 'M':
		return true
	case len(b) >= 4 && b[0] == 0 && b[1] == 0 && b[2] == 1 && b[3] == 0:
		return true
	}
	// SVG: markup, possibly after a byte order mark
	t := bytes.TrimLeft(bytes.TrimPrefix(b, []byte("\xef\xbb\xbf")), " \t\r\n")
	return len(t) > 0 && t[0] == '<'
}

func (d *Document) compress(data []byte) ([]byte, error) { return compress(data, d.CompressionLevel) }

func compress(data []byte, level int) ([]byte, error) {
	var b bytes.Buffer
	w, err := flate.NewWriter(&b, level)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

type encodedPart struct {
	entry PartEntry
	data  []byte
}

// encodeParts compresses parts and builds the parts table with offsets.
func (d *Document) encodeParts() ([]encodedPart, error) {
	out := make([]encodedPart, 0, len(d.order))
	var off int64
	for _, p := range d.Parts() {
		ep := encodedPart{entry: PartEntry{H: p.Hash, T: p.Type, Enc: EncIdentity, Size: len(p.Data), Off: off}, data: p.Data}
		if d.shouldCompress(p) {
			c, err := d.compress(p.Data)
			if err != nil {
				return nil, err
			}
			if len(c) < len(p.Data) {
				ep.entry.Enc = EncDeflateRaw
				ep.data = c
			}
		}
		ep.entry.Len = len(ep.data)
		off += int64(ep.entry.Len)
		out = append(out, ep)
	}
	return out, nil
}

func (d *Document) manifest(parts []encodedPart) *Manifest {
	m := &Manifest{BDF: FormatVersion, Opset: OpsetVersion, Unit: "pt", Meta: d.Meta, Views: d.Views}
	if m.Meta.Generator == "" {
		m.Meta.Generator = "bdf-go/0.1"
	}
	m.Parts = make([]PartEntry, len(parts))
	for i, p := range parts {
		m.Parts[i] = p.entry
	}
	return m
}

// Magic is the single-file magic number: "bdf" and a NUL byte.
var Magic = [4]byte{'b', 'd', 'f', 0}

// HeaderSize is the fixed header length of the single-file form.
const HeaderSize = 32

// stored builds what is written: the manifest (with offsets), the stored
// bytes of each part in manifest order, and the header flags. With a Lock
// the parts and the manifest are sealed and the outer manifest is returned.
func (d *Document) stored() (*Manifest, [][]byte, uint16, error) {
	parts, err := d.encodeParts()
	if err != nil {
		return nil, nil, 0, err
	}
	m := d.manifest(parts)
	data := make([][]byte, len(parts))
	for i, p := range parts {
		data[i] = p.data
	}
	if d.Lock == nil {
		return m, data, 0, nil
	}
	a, err := d.Lock.aead()
	if err != nil {
		return nil, nil, 0, err
	}
	outer := &Manifest{BDF: FormatVersion, Parts: make([]PartEntry, 0, len(parts)+1)}
	sealed := make([][]byte, 0, len(parts)+1)
	add := func(b []byte) Hash {
		h := HashOf(b)
		outer.Parts = append(outer.Parts, PartEntry{H: h, T: PartSealed, Enc: EncIdentity, Len: len(b), Size: len(b)})
		sealed = append(sealed, b)
		return h
	}
	add(nil) // the manifest goes first; sealed below, once the part names are known
	for i := range m.Parts {
		e := &m.Parts[i]
		e.Off = 0
		e.Sealed = add(seal(a, data[i], e.H[:]))
	}
	mj, err := json.Marshal(m)
	if err != nil {
		return nil, nil, 0, err
	}
	menc := EncIdentity
	if c, err := d.compress(mj); err == nil && len(c) < len(mj) {
		mj, menc = c, EncDeflateRaw
	}
	sm := seal(a, mj, manifestAAD)
	h := HashOf(sm)
	outer.Parts[0] = PartEntry{H: h, T: PartSealed, Enc: EncIdentity, Len: len(sm), Size: len(sm)}
	sealed[0] = sm
	outer.Encryption = &Encryption{Cipher: CipherA256GCM, Keys: d.Lock.slots, Manifest: SealedManifest{Part: h, Enc: menc}}
	var off int64
	for i := range outer.Parts {
		outer.Parts[i].Off = off
		off += int64(outer.Parts[i].Len)
	}
	return outer, sealed, FlagEncrypted, nil
}

// WriteSingle writes the single-file form.
func (d *Document) WriteSingle(w io.Writer) error {
	m, data, flags, err := d.stored()
	if err != nil {
		return err
	}
	return writeSingle(w, m, data, flags, d.CompressionLevel, d.MinCompress)
}

// WriteSplit writes the split form into dir (manifest.json and parts/<hash>).
func (d *Document) WriteSplit(dir string) error {
	m, data, _, err := d.stored()
	if err != nil {
		return err
	}
	return writeSplit(dir, m, data)
}

// writeSingle writes a manifest whose part offsets follow data, and data.
func writeSingle(w io.Writer, m *Manifest, data [][]byte, flags uint16, level, minCompress int) error {
	mj, err := json.Marshal(m)
	if err != nil {
		return err
	}
	menc := byte(0)
	if len(mj) >= minCompress {
		if c, err := compress(mj, level); err == nil && len(c) < len(mj) {
			mj, menc = c, 1
		}
	}
	var h buf
	h.bytes(Magic[:])
	h.u16(FormatVersion)
	h.u16(flags)
	h.u64(HeaderSize)
	h.u64(uint64(len(mj)))
	h.u8(menc)
	h.bytes(make([]byte, 7))
	if _, err := w.Write(h.b); err != nil {
		return err
	}
	if _, err := w.Write(mj); err != nil {
		return err
	}
	for _, b := range data {
		if _, err := w.Write(b); err != nil {
			return err
		}
	}
	return nil
}

// writeSplit writes manifest.json (without offsets) and a file per part.
func writeSplit(dir string, m *Manifest, data [][]byte) error {
	sm := *m
	sm.Parts = slices.Clone(m.Parts)
	for i := range sm.Parts {
		sm.Parts[i].Off = 0
	}
	mj, err := json.MarshalIndent(&sm, "", " ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(dir, "parts"), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), mj, 0o644); err != nil {
		return err
	}
	for i, e := range sm.Parts {
		if err := os.WriteFile(filepath.Join(dir, "parts", e.H.String()), data[i], 0o644); err != nil {
			return err
		}
	}
	return nil
}
