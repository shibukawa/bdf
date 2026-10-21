package manager

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image/png"
	"io"
	"os"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

// Input is a document sent to a location.
type Input struct {
	// Name is the job's name (job-name, or the file's).
	Name string
	// User is who the sender says sent it (requesting-user-name).
	User string
	// Source is "ipp" or "api".
	Source string
	// Doc is the document, of Size bytes: .bdf, or what the converter
	// reads.
	Doc  io.ReaderAt
	Size int64
}

// ErrFormat is the error of a document neither .bdf nor of a format the
// converter reads.
var ErrFormat = errors.New("manager: the document's format is not known")

// Submit takes a document for a location's devices: it is made a .bdf
// document, kept, and put on the location's board.
func (m *Manager) Submit(ctx context.Context, loc string, in Input) (*Slip, error) {
	l, err := m.Location(ctx, loc)
	if err != nil {
		return nil, err
	}
	doc, format, err := m.document(ctx, in)
	if err != nil {
		return nil, err
	}
	v, err := print.PrintView(doc, "")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFormat, err)
	}
	f, err := os.CreateTemp(m.o.TempDir, "bdf-pm-*.bdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	h := sha256.New()
	if err := doc.WriteSingle(io.MultiWriter(f, h)); err != nil {
		return nil, err
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	job, err := m.nextJob(ctx)
	if err != nil {
		return nil, err
	}
	now := m.now()
	s := Slip{ID: slipID(job), Job: job, Location: loc, Name: clip(in.Name, 255), User: clip(in.User, 255), Source: in.Source,
		Format: format, Pages: len(v.Pages), SHA256: hex.EncodeToString(h.Sum(nil)), Received: now, State: SlipNew}
	if l.Fallback != nil && l.Fallback.After > 0 {
		s.FallbackAt = now.Add(l.Fallback.After.d())
	}
	meta, err := m.st.docs.Insert(ctx, s.ID, f)
	if err != nil {
		return nil, mapErr(err)
	}
	s.Size = meta.Size
	if _, err := m.st.db.Insert(ctx, s.ID, slipRecord(s)); err != nil {
		m.st.docs.Remove(context.WithoutCancel(ctx), s.ID)
		return nil, mapErr(err)
	}
	if err := m.touch(ctx, loc, func(st *locationState) { st.Open = append(st.Open, s.ID) }); err != nil {
		return nil, err
	}
	m.o.Log.Info("slip", "location", loc, "job", job, "name", s.Name, "pages", s.Pages, "source", s.Source)
	return &s, nil
}

// document reads or converts what a sender sent.
func (m *Manager) document(ctx context.Context, in Input) (*bdf.Document, string, error) {
	if in.Size > m.o.MaxDocument {
		return nil, "", fmt.Errorf("%w: %d bytes; at most %d", ErrInvalid, in.Size, m.o.MaxDocument)
	}
	head := make([]byte, len(bdf.Magic))
	if n, _ := in.Doc.ReadAt(head, 0); n == len(head) && bytes.Equal(head, bdf.Magic[:]) {
		r, err := bdf.OpenSingle(in.Doc, in.Size)
		if err != nil {
			return nil, "", fmt.Errorf("%w: %v", ErrFormat, err)
		}
		defer r.Close()
		if r.Locked() {
			return nil, "", fmt.Errorf("%w: the document is locked with a password", ErrFormat)
		}
		doc, err := r.ToDocument()
		if err != nil {
			return nil, "", fmt.Errorf("%w: %v", ErrFormat, err)
		}
		return doc, "bdf", nil
	}
	if m.o.Converter == nil {
		return nil, "", ErrFormat
	}
	format, ok := m.o.Converter.Detect(in.Doc, in.Size, in.Name)
	if !ok {
		return nil, "", ErrFormat
	}
	doc, _, err := m.o.Converter.Convert(ctx, in.Doc, in.Size, format, in.Name, "")
	if err != nil {
		return nil, "", fmt.Errorf("%w: %s: %v", ErrFormat, format, err)
	}
	if doc.Lock != nil {
		return nil, "", fmt.Errorf("%w: the document needs a password", ErrFormat)
	}
	return doc, format, nil
}

// Document returns the document of a slip.
func (m *Manager) Document(ctx context.Context, s *Slip) (*bdf.Document, error) {
	r, err := m.st.docs.Open(ctx, s.ID)
	if err != nil {
		return nil, mapErr(err)
	}
	defer r.Close()
	br, err := bdf.OpenSingle(r, r.Size)
	if err != nil {
		return nil, err
	}
	defer br.Close()
	return br.ToDocument()
}

// WriteDocument writes the .bdf file of a slip.
func (m *Manager) WriteDocument(ctx context.Context, s *Slip, w io.Writer) error {
	r, err := m.st.docs.Open(ctx, s.ID)
	if err != nil {
		return mapErr(err)
	}
	defer r.Close()
	_, err = io.Copy(w, r)
	return err
}

// MaxPageWidth bounds the width of a page drawn for a device.
const MaxPageWidth = 2400

// WritePage draws a page (from 0) of a slip as PNG, width pixels wide.
func (m *Manager) WritePage(ctx context.Context, s *Slip, page, width int, w io.Writer) error {
	if page < 0 || page >= s.Pages {
		return fmt.Errorf("%w: page %d of %d", ErrNotFound, page+1, s.Pages)
	}
	width = min(max(width, 200), MaxPageWidth)
	doc, err := m.Document(ctx, s)
	if err != nil {
		return err
	}
	v, err := print.PrintView(doc, "")
	if err != nil {
		return err
	}
	p := v.Pages[page]
	if p == nil || p.W <= 0 {
		return fmt.Errorf("%w: page %d has no size", ErrInvalid, page+1)
	}
	r := imagebdf.New(doc, &imagebdf.Options{FontDirs: m.o.FontDirs, NoSystemFonts: m.o.NoSystemFonts})
	img, err := r.Page(v, page, float64(width)/float64(p.W))
	if err != nil {
		return err
	}
	return (&png.Encoder{CompressionLevel: png.BestSpeed}).Encode(w, img)
}

// clip cuts a string to n bytes, at a rune.
func clip(s string, n int) string {
	if len(s) <= n {
		return s
	}
	for n > 0 && s[n]&0xc0 == 0x80 {
		n--
	}
	return s[:n]
}
