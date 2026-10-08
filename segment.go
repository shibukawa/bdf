package bdf

import (
	"fmt"
	"io"
)

// Segments (docs/spec.md §3.6): a server that keeps a document hands it to
// a reader a few pages at a time, each segment sealed for a key pair the
// reader made for that request (NewECDHLock). A segment is a whole document
// in form, every view and page with the metadata, but only the pages of the
// segment have layers, and its parts are those these pages need that the
// reader does not hold from earlier segments. The stored parts are copied
// as they are, so a segment costs the server reading and sealing, not
// compressing.

// Segment names the pages [From, To) of a view (0-based).
type Segment struct {
	View string `json:"view"`
	From int    `json:"from"`
	To   int    `json:"to"`
}

// SegmentAt returns the segment of size pages that holds page of v: a
// view's segments are pages [0, size), [size, 2·size) and so on, the last
// one shorter.
func SegmentAt(v *View, page, size int) Segment {
	from := page / size * size
	return Segment{View: v.ID, From: from, To: min(from+size, len(v.Pages))}
}

// SegmentOptions say which pages a segment carries and how it is sealed.
type SegmentOptions struct {
	Segment
	// Have lists the pages the reader holds from earlier segments: the
	// parts they need are left out.
	Have []Segment
	// Lock seals the segment (NewECDHLock); nil writes it in the clear.
	Lock *Lock
}

// segmentLevel is the flate level of a segment's manifest, which is made
// for each request: the standard library's default (see defaultLevel).
const segmentLevel = defaultLevel

type pageKey struct {
	view string
	page int
}

// WriteSegment writes a segment of the document in the single-file form
// (docs/spec.md §3.6). The document must be one of pages: its views fixed,
// flow or scroll, without music, and an encrypted one unlocked. The views
// of the segment have no text index (the text of the whole view). It may
// be called from several goroutines at once.
func (r *Reader) WriteSegment(w io.Writer, o SegmentOptions) error {
	if r.Locked() {
		return ErrLocked
	}
	m := r.Manifest
	for _, v := range m.Views {
		if v.Kind != ViewFixed && v.Kind != ViewFlow && v.Kind != ViewScroll || v.Play != nil {
			return fmt.Errorf("bdf: view %q: segments carry views of pages without music, not a %s view", v.ID, v.Kind)
		}
	}
	need, err := r.segmentParts(o.Segment)
	if err != nil {
		return err
	}
	have := map[Hash]bool{}
	for _, s := range o.Have {
		hs, err := r.segmentParts(s)
		if err != nil {
			return err
		}
		for h := range hs {
			have[h] = true
		}
	}
	s := o.Segment
	out := &Manifest{BDF: FormatVersion, Opset: m.Opset, Unit: m.Unit, Meta: m.Meta, Segment: &s, Views: outline(m.Views, s), Parts: []PartEntry{}}
	var data [][]byte
	var off int64
	for _, e := range m.Parts { // in the order of the document: the recommended one
		if !need[e.H] || have[e.H] {
			continue
		}
		b, err := r.storedPlain(e)
		if err != nil {
			return err
		}
		if len(b) != e.Len {
			return fmt.Errorf("bdf: part %s is %d bytes, not %d", e.H, len(b), e.Len)
		}
		e.Off, e.Sealed = off, Hash{}
		off += int64(e.Len)
		out.Parts = append(out.Parts, e)
		data = append(data, b)
	}
	flags := uint16(0)
	if o.Lock != nil {
		if out, data, err = sealStored(out, data, o.Lock, segmentLevel); err != nil {
			return err
		}
		flags = FlagEncrypted
	}
	return writeSingle(w, out, data, flags, segmentLevel, defaultMinCompress)
}

// outline copies views with the layers of the pages of s only, and without
// text indexes.
func outline(views []*View, s Segment) []*View {
	out := make([]*View, len(views))
	for i, v := range views {
		c := *v
		c.TextIndex = ""
		c.Pages = make([]*Page, len(v.Pages))
		for j, p := range v.Pages {
			q := *p
			if v.ID != s.View || j < s.From || j >= s.To {
				q.Layers = []Layer{}
			}
			c.Pages[j] = &q
		}
		out[i] = &c
	}
	return out
}

// segmentParts returns the parts the pages of s need.
func (r *Reader) segmentParts(s Segment) (map[Hash]bool, error) {
	var v *View
	for _, w := range r.Manifest.Views {
		if w.ID == s.View {
			v = w
			break
		}
	}
	if v == nil {
		return nil, fmt.Errorf("bdf: no view %q", s.View)
	}
	if s.From < 0 || s.From >= s.To || s.To > len(v.Pages) {
		return nil, fmt.Errorf("bdf: view %q has no pages %d to %d", s.View, s.From, s.To)
	}
	out := map[Hash]bool{}
	for i := s.From; i < s.To; i++ {
		hs, err := r.pageNeeds(v, i)
		if err != nil {
			return nil, err
		}
		for _, h := range hs {
			out[h] = true
		}
	}
	return out, nil
}

// pageNeeds returns the parts page i of v needs: its layers' objects and
// everything they reference. It keeps what it found for the next segments.
func (r *Reader) pageNeeds(v *View, i int) ([]Hash, error) {
	k := pageKey{v.ID, i}
	r.pageMu.Lock()
	hs, ok := r.pageParts[k]
	r.pageMu.Unlock()
	if ok {
		return hs, nil
	}
	seen := map[Hash]bool{}
	var walk func(h Hash) error
	walk = func(h Hash) error {
		if seen[h] {
			return nil
		}
		seen[h] = true
		e, ok := r.entries[h]
		if !ok {
			return fmt.Errorf("bdf: page %d of view %q: unknown part %s", i+1, v.ID, h)
		}
		hs = append(hs, h)
		if e.T != PartObject {
			return nil
		}
		o, err := r.Object(h)
		if err != nil {
			return err
		}
		for _, d := range o.Deps() {
			if err := walk(d); err != nil {
				return err
			}
		}
		return nil
	}
	for _, l := range v.Pages[i].Layers {
		if err := walk(l.Obj); err != nil {
			return nil, err
		}
	}
	r.pageMu.Lock()
	if r.pageParts == nil {
		r.pageParts = map[pageKey][]Hash{}
	}
	r.pageParts[k] = hs
	r.pageMu.Unlock()
	return hs, nil
}
