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

// Part is a raw (uncompressed) part.
type Part struct {
	Hash Hash
	Type string
	Data []byte
}

// Document accumulates parts and views and writes them out.
type Document struct {
	Meta  Meta
	Views []*View

	// CompressionLevel is the flate level for deflate-raw parts (default: flate.BestCompression).
	CompressionLevel int
	// MinCompress is the smallest decoded size that gets compressed (default 512).
	MinCompress int

	parts map[Hash]*Part
	order []Hash
}

// NewDocument creates an empty document.
func NewDocument() *Document {
	return &Document{
		CompressionLevel: flate.BestCompression,
		MinCompress:      512,
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
		return false
	case PartFont:
		// WOFF and WOFF2 are compressed already; TTF/OTF are not.
		if len(p.Data) >= 4 && (string(p.Data[:4]) == "wOF2" || string(p.Data[:4]) == "wOFF") {
			return false
		}
	}
	return len(p.Data) >= d.MinCompress
}

func (d *Document) compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := flate.NewWriter(&b, d.CompressionLevel)
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

// Magic is the single-file magic number.
var Magic = [4]byte{'B', 'D', 'F', '1'}

// HeaderSize is the fixed header length of the single-file form.
const HeaderSize = 32

// WriteSingle writes the single-file form.
func (d *Document) WriteSingle(w io.Writer) error {
	parts, err := d.encodeParts()
	if err != nil {
		return err
	}
	mj, err := json.Marshal(d.manifest(parts))
	if err != nil {
		return err
	}
	menc := byte(0)
	if len(mj) >= d.MinCompress {
		if c, err := d.compress(mj); err == nil && len(c) < len(mj) {
			mj, menc = c, 1
		}
	}
	var h buf
	h.bytes(Magic[:])
	h.u16(FormatVersion)
	h.u16(0)
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
	for _, p := range parts {
		if _, err := w.Write(p.data); err != nil {
			return err
		}
	}
	return nil
}

// WriteSplit writes the split form into dir (manifest.json and parts/<hash>).
func (d *Document) WriteSplit(dir string) error {
	parts, err := d.encodeParts()
	if err != nil {
		return err
	}
	m := d.manifest(parts)
	for i := range m.Parts {
		m.Parts[i].Off = 0
	}
	mj, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(dir, "parts"), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), mj, 0o644); err != nil {
		return err
	}
	for _, p := range parts {
		if err := os.WriteFile(filepath.Join(dir, "parts", p.entry.H.String()), p.data, 0o644); err != nil {
			return err
		}
	}
	return nil
}
