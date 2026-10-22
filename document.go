package bdf

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shibukawa/bdf/encoding/parquet"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/klauspost/compress/flate"
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
	// Tables are the tables the document holds beside its drawing
	// (docs/spec.md §4.6): see AddTable.
	Tables []*Table

	// CompressionLevel is the flate level of the deflate-raw parts: 1
	// (BestSpeed) to 9 (BestCompression). The default, 6, is what
	// DefaultCompression names too; 9 writes files a few percent smaller
	// in four times the time of the compression.
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
		CompressionLevel: defaultLevel,
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
	case PartAudio:
		// An audio file is stored as it is, whatever its coding (spec
		// §3.2): a reader hands its bytes to the browser's player.
		return false
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

type deflateWriter struct {
	writer *flate.Writer
	output bytes.Buffer
}

// The compressor's working memory is much larger than most document
// parts. Reuse it at the same level; returned bytes are never pooled.
var deflateWriters [flate.BestCompression - flate.HuffmanOnly + 1]sync.Pool

// defaultLevel is used when Document.CompressionLevel is
// flate.DefaultCompression.
const defaultLevel = 6

func compress(data []byte, level int) ([]byte, error) {
	if level < flate.HuffmanOnly || level > flate.BestCompression {
		_, err := flate.NewWriter(io.Discard, level)
		return nil, err
	}
	if level == flate.DefaultCompression {
		level = defaultLevel
	}
	pool := &deflateWriters[level-flate.HuffmanOnly]
	var enc *deflateWriter
	if reused := pool.Get(); reused != nil {
		enc = reused.(*deflateWriter)
		enc.writer.Reset(&enc.output)
	} else {
		enc = &deflateWriter{}
		var err error
		enc.writer, err = flate.NewWriter(&enc.output, level)
		if err != nil {
			return nil, err
		}
	}
	defer func() {
		// The caller owns the output. The pooled writer must not keep it
		// alive or overwrite it on the next compression.
		enc.output = bytes.Buffer{}
		pool.Put(enc)
	}()
	if _, err := enc.writer.Write(data); err != nil {
		return nil, err
	}
	if err := enc.writer.Close(); err != nil {
		return nil, err
	}
	return enc.output.Bytes(), nil
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
	m := &Manifest{BDF: FormatVersion, Opset: OpsetVersion, Unit: "pt", Meta: d.Meta, Views: d.Views, Tables: d.Tables}
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
	outer, sealed, err := sealStored(m, data, d.Lock, d.CompressionLevel)
	if err != nil {
		return nil, nil, 0, err
	}
	return outer, sealed, FlagEncrypted, nil
}

// sealStored seals the stored parts of m (data, in manifest order) and m
// itself with lock, and returns the outer manifest, with offsets, and the
// sealed parts, the manifest first (docs/spec.md §3.5). It sets the Sealed
// names of m's entries. level is the flate level the manifest is
// compressed with.
func sealStored(m *Manifest, data [][]byte, lock *Lock, level int) (*Manifest, [][]byte, error) {
	a, err := lock.aead()
	if err != nil {
		return nil, nil, err
	}
	outer := &Manifest{BDF: FormatVersion, Parts: make([]PartEntry, 0, len(m.Parts)+1)}
	sealed := make([][]byte, 0, len(m.Parts)+1)
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
		return nil, nil, err
	}
	menc := EncIdentity
	if c, err := compress(mj, level); err == nil && len(c) < len(mj) {
		mj, menc = c, EncDeflateRaw
	}
	sm := seal(a, mj, manifestAAD)
	h := HashOf(sm)
	outer.Parts[0] = PartEntry{H: h, T: PartSealed, Enc: EncIdentity, Len: len(sm), Size: len(sm)}
	sealed[0] = sm
	outer.Encryption = &Encryption{Cipher: CipherA256GCM, Keys: lock.slots, Manifest: SealedManifest{Part: h, Enc: menc}}
	var off int64
	for i := range outer.Parts {
		outer.Parts[i].Off = off
		off += int64(outer.Parts[i].Len)
	}
	return outer, sealed, nil
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

// AddTable adds a table the document holds beside its drawing (docs/spec.md
// §4.6): t names it, and data is the Apache Parquet file of its rows,
// which becomes a PartTable part (t.Part is set). The id must be new to
// the document; the rows and the columns of t should be those of the
// file (AddTableRows makes both of a table in memory).
func (d *Document) AddTable(t *Table, data []byte) error {
	switch {
	case t == nil || t.ID == "":
		return errors.New("bdf: a table needs an id")
	case d.Table(t.ID) != nil:
		return fmt.Errorf("bdf: a table %q already", t.ID)
	case len(d.Tables) >= MaxTables:
		return fmt.Errorf("bdf: more than %d tables", MaxTables)
	case len(t.Columns) > MaxTableColumns:
		return fmt.Errorf("bdf: table %s: %d columns, more than %d", t.ID, len(t.Columns), MaxTableColumns)
	case len(data) < 12 || string(data[:4]) != "PAR1":
		return fmt.Errorf("bdf: table %s: not a Parquet file", t.ID)
	case t.Rows < 0 || t.Page < 0:
		return fmt.Errorf("bdf: table %s: numbers out of range", t.ID)
	}
	t.Part = d.AddPart(PartTable, data).String()
	d.Tables = append(d.Tables, t)
	return nil
}

// AddTableRows adds a table in memory (encoding/parquet.Table) as a table
// of the document, written as a Parquet file: id and title name it, and
// view and page say where the figure it is the data of is (page from 1;
// 0 for a sheet, "" for none).
func (d *Document) AddTableRows(id, title string, rows *parquet.Table, view string, page int) (*Table, error) {
	data, err := parquet.Encode(rows)
	if err != nil {
		return nil, fmt.Errorf("bdf: table %s: %w", id, err)
	}
	t := &Table{ID: id, Title: title, Rows: int64(len(rows.Rows)), Columns: tableColumns(rows.Columns), View: view, Page: page}
	if err := d.AddTable(t, data); err != nil {
		return nil, err
	}
	return t, nil
}

// tableColumns lists columns as the manifest has them.
func tableColumns(cols []parquet.Column) []TableColumn {
	out := make([]TableColumn, len(cols))
	for i, c := range cols {
		out[i] = TableColumn{Name: c.Name, Type: c.Type.String()}
	}
	return out
}

// Table returns the table with an id, or nil.
func (d *Document) Table(id string) *Table {
	for _, t := range d.Tables {
		if t.ID == id {
			return t
		}
	}
	return nil
}
