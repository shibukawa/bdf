package pdf

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"

	"github.com/shibukawa/bdf"
)

// Stream converts a PDF a page at a time, in the order the pages are asked
// for, so that a viewer can show the pages as they are converted and convert
// the ones the reader looks at first. NewStream reads the document and sizes
// its pages (Outline), Page converts one page and returns what a reader
// holding the outline and the pages returned before needs to draw it, and
// Finish converts the pages left and makes the document Convert makes: the
// same bytes when the pages were converted in order.
//
// The page documents cannot share fonts the way the finished document does:
// a font is subset to the glyphs its pages use, and later pages use more.
// Each page document gets a version of the font rebuilt for the codes of its
// objects, unless a version returned before holds them all; a new version
// also takes the codes of the last one while it is small, so that the
// versions of a font whose pages keep to a small set of glyphs (Latin text)
// grow into one that serves the pages after.
//
// A Stream is not safe for concurrent use.
type Stream struct {
	c         *converter
	protected bool
	finished  bool

	sent     map[bdf.Hash]bool              // parts returned in the page documents
	encoded  map[*pending]bdf.Hash          // objects returned, with the hash they were returned under
	created  map[*pending]int               // the page whose conversion made each object
	codes    []map[*pdfFont]map[uint32]bool // by page: the codes each font drew
	versions map[*pdfFont][]*fontVersion    // font versions returned
}

// fontVersion is a font rebuilt for some of its codes.
type fontVersion struct {
	codes map[uint32]bool
	font  bdf.Font
}

// streamFontGrowth is the number of codes up to which a new font version
// also holds the codes of the last one.
const streamFontGrowth = 512

// NewStream reads a PDF for a conversion done a page at a time.
func NewStream(rs io.ReadSeeker, opts *Options) (*Stream, error) {
	c, protected, err := newConverter(rs, opts)
	if err != nil {
		return nil, err
	}
	return &Stream{c: c, protected: protected, sent: map[bdf.Hash]bool{}, encoded: map[*pending]bdf.Hash{},
		created: map[*pending]int{}, codes: make([]map[*pdfFont]map[uint32]bool, len(c.pageBodies)),
		versions: map[*pdfFont][]*fontVersion{}}, nil
}

// Pages returns the number of pages of the document's view.
func (s *Stream) Pages() int { return len(s.c.pageBodies) }

// Protected reports that the PDF opened only with Options.Password.
func (s *Stream) Protected() bool { return s.protected }

// Outline returns the document without page content: its metadata and its
// view with every page sized, but without layers until Page converts it.
func (s *Stream) Outline() *bdf.Document {
	doc := bdf.NewDocument()
	doc.Meta = s.c.doc.Meta
	view := *s.c.doc.Views[0]
	view.Pages = make([]*bdf.Page, len(s.c.pageBodies))
	for i, pr := range s.c.pageBodies {
		view.Pages[i] = &bdf.Page{W: pr.page.W, H: pr.page.H, Body: pr.page.Body, Layers: []bdf.Layer{}}
	}
	doc.Views = []*bdf.View{&view}
	return doc
}

// Page converts page i of the view (0-based), unless an earlier call did,
// and returns a document whose only view holds that page, with the parts it
// needs that no earlier call returned: a reader puts the page in its outline
// and adds the parts to those it has.
func (s *Stream) Page(i int) (*bdf.Document, error) {
	c := s.c
	if s.finished {
		return nil, errors.New("pdf: the stream is finished")
	}
	if i < 0 || i >= len(c.pageBodies) {
		return nil, fmt.Errorf("pdf: no page %d", i)
	}
	pr := c.pageBodies[i]
	if pr.body == nil {
		start := len(c.pendings)
		c.pageCodes = map[*pdfFont]map[uint32]bool{}
		err := c.convertPage(pr)
		s.codes[i], c.pageCodes = c.pageCodes, nil
		for _, p := range c.pendings[start:] {
			s.created[p] = i
		}
		if err != nil {
			return nil, fmt.Errorf("pdf: page %d: %w", pr.nr, err)
		}
	}
	return s.pageDocument(pr), nil
}

// pageDocument encodes the objects of a converted page that were not
// returned before, with the font versions they need.
func (s *Stream) pageDocument(pr *pageRef) *bdf.Document {
	c := s.c
	var objs []*pending // children first
	seen := map[*pending]bool{}
	var visit func(p *pending)
	visit = func(p *pending) {
		if seen[p] {
			return
		}
		seen[p] = true
		if _, ok := s.encoded[p]; ok {
			return
		}
		for _, r := range slices.Sorted(maps.Keys(p.children)) {
			visit(p.children[r])
		}
		objs = append(objs, p)
	}
	visit(pr.body)

	// The codes each font draws in these objects: those of the pages that
	// made them.
	var fonts []*pdfFont
	need := map[*pdfFont]map[uint32]bool{}
	pages := map[int]bool{}
	for _, p := range objs {
		pages[s.created[p]] = true
		for _, r := range slices.Sorted(maps.Keys(p.fonts)) {
			if f := p.fonts[r]; need[f] == nil {
				need[f] = map[uint32]bool{}
				fonts = append(fonts, f)
			}
		}
	}
	for pg := range pages {
		for f, codes := range s.codes[pg] {
			if n := need[f]; n != nil {
				maps.Copy(n, codes)
			}
		}
	}

	doc := bdf.NewDocument()
	add := func(typ string, data []byte) bdf.Hash {
		h := bdf.HashOf(data)
		if !s.sent[h] {
			s.sent[h] = true
			doc.AddPart(typ, data)
		}
		return h
	}
	chosen := map[*pdfFont]bdf.Font{}
	for _, f := range fonts {
		chosen[f] = s.fontFor(f, need[f], add)
	}
	for _, p := range objs {
		for r, f := range p.fonts {
			p.obj.UpdateFont(r, chosen[f])
		}
		for r, child := range p.children {
			p.obj.UpdateObject(r, s.encoded[child])
		}
		for _, dep := range p.obj.Deps() {
			// images and path collections are in the document already
			if part := c.doc.Part(dep); part != nil {
				add(part.Type, part.Data)
			}
		}
		s.encoded[p] = add(bdf.PartObject, p.obj.Encode())
	}
	view := doc.NewView(c.doc.Views[0].ID, c.doc.Views[0].Kind, "")
	page := view.AddPage(pr.page.W, pr.page.H, bdf.Layer{Role: bdf.RoleBody, Obj: s.encoded[pr.body]})
	page.Body = pr.page.Body
	return doc
}

// fontFor returns a version of f that maps the codes in need, rebuilding one
// when no version returned before does.
func (s *Stream) fontFor(f *pdfFont, need map[uint32]bool, add func(string, []byte) bdf.Hash) bdf.Font {
	vs := s.versions[f]
	for i := len(vs) - 1; i >= 0; i-- {
		if covers(vs[i].codes, need) {
			return vs[i].font
		}
	}
	codes := need
	if n := len(vs); n > 0 && len(vs[n-1].codes) < streamFontGrowth {
		codes = maps.Clone(vs[n-1].codes)
		maps.Copy(codes, need)
	}
	used := make(map[uint32]*usedGlyph, len(codes))
	for code := range codes {
		if u := f.used[code]; u != nil {
			used[code] = u
		}
	}
	font := s.c.embedFont(f, used, func(data []byte) bdf.Hash { return add(bdf.PartFont, data) })
	s.versions[f] = append(vs, &fontVersion{codes: codes, font: font})
	return font
}

func covers(have, need map[uint32]bool) bool {
	for code := range need {
		if !have[code] {
			return false
		}
	}
	return true
}

// Finish converts the pages no Page call converted and returns the whole
// document. The stream is done with.
func (s *Stream) Finish() (*Result, error) {
	c := s.c
	if s.finished {
		return nil, errors.New("pdf: the stream is finished")
	}
	s.finished = true
	s.sent, s.encoded, s.created, s.codes, s.versions = nil, nil, nil, nil, nil
	for _, pr := range c.pageBodies {
		if pr.body != nil {
			continue
		}
		if err := c.convertPage(pr); err != nil {
			return nil, fmt.Errorf("pdf: page %d: %w", pr.nr, err)
		}
	}
	return c.finish(s.protected), nil
}

// use is f.use, recording the code for the page a Stream converts.
func (in *interp) use(f *pdfFont, g glyphCode) (draw, text string) {
	if codes := in.c.pageCodes; codes != nil {
		set := codes[f]
		if set == nil {
			set = map[uint32]bool{}
			codes[f] = set
		}
		set[g.code] = true
	}
	return f.use(g)
}
