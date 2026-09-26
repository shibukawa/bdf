package pdf

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/sfnt"
)

// streamReader is what a viewer holding the page documents of a Stream has.
type streamReader struct {
	t     *testing.T
	parts map[bdf.Hash]*bdf.Part
	objs  map[bdf.Hash]*bdf.ObjectPart
	cmaps map[bdf.Hash]map[uint32]uint16
}

func (r *streamReader) add(doc *bdf.Document) {
	for _, p := range doc.Parts() {
		if r.parts[p.Hash] != nil {
			r.t.Errorf("part %s returned twice", p.Hash)
		}
		r.parts[p.Hash] = p
	}
}

func (r *streamReader) object(h bdf.Hash) *bdf.ObjectPart {
	if o := r.objs[h]; o != nil {
		return o
	}
	p := r.parts[h]
	if p == nil {
		r.t.Fatalf("object %s was not returned", h)
	}
	o, err := bdf.DecodeObject(p.Data)
	if err != nil {
		r.t.Fatal(err)
	}
	r.objs[h] = o
	return o
}

func (r *streamReader) cmap(h bdf.Hash) map[uint32]uint16 {
	if m, ok := r.cmaps[h]; ok {
		return m
	}
	f, err := sfnt.Parse(r.parts[h].Data)
	if err != nil {
		r.t.Fatalf("font %s: %v", h, err)
	}
	r.cmaps[h] = f.Cmap
	return f.Cmap
}

// check verifies that an object and those it uses can be drawn with the
// parts returned: they are all there, and the fonts map every character
// the text is drawn with.
func (r *streamReader) check(h bdf.Hash) {
	o := r.object(h)
	for _, img := range o.Images {
		if r.parts[img] == nil {
			r.t.Errorf("image %s was not returned", img)
		}
	}
	for _, f := range o.Fonts {
		if f.Kind == bdf.FontEmbedded && r.parts[f.Hash] == nil {
			r.t.Errorf("font %s was not returned", f.Hash)
		}
	}
	var font *bdf.Font
	err := o.Walk(func(in bdf.Instr) {
		switch in.Op {
		case bdf.OpFont:
			font = &o.Fonts[in.Args[0].(uint64)]
		case bdf.OpFillText, bdf.OpStrokeText:
			if font == nil || font.Kind != bdf.FontEmbedded {
				return
			}
			cm := r.cmap(font.Hash)
			for _, c := range in.Args[0].(string) {
				if _, ok := cm[uint32(c)]; !ok {
					r.t.Errorf("font %s does not map %U (text %q)", font.Hash, c, in.Args[0])
				}
			}
		}
	})
	if err != nil {
		r.t.Fatal(err)
	}
	for _, child := range o.Objects {
		r.check(child)
	}
}

func layersText(t *testing.T, obj func(bdf.Hash) *bdf.ObjectPart, page *bdf.Page) string {
	var b strings.Builder
	for _, l := range page.Layers {
		runs, err := bdf.ExtractText(obj(l.Obj), obj)
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range runs {
			b.WriteString(r.Text)
		}
	}
	return b.String()
}

func singleBytes(t *testing.T, doc *bdf.Document) []byte {
	var b bytes.Buffer
	if err := doc.WriteSingle(&b); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// TestStream converts the test PDFs a page at a time in several orders:
// each page document holds what its page needs beyond the pages before,
// and Finish makes the document Convert makes (the same bytes when the
// pages were converted in order).
func TestStream(t *testing.T) {
	files, _ := filepath.Glob("testdata/*.pdf")
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		opts := func() *Options { return &Options{NoWOFF2: true} } // TrueType/OpenType fonts, whose cmaps the test reads
		want, err := Convert(bytes.NewReader(data), opts())
		if err != nil {
			if strings.Contains(err.Error(), "password") {
				continue
			}
			t.Fatalf("%s: %v", file, err)
		}
		wantPages := want.Doc.Views[0].Pages
		n := len(wantPages)
		orders := map[string][]int{"in order": nil, "reversed": nil, "middle out": nil}
		for i := range n {
			orders["in order"] = append(orders["in order"], i)
			orders["reversed"] = append(orders["reversed"], n-1-i)
			orders["middle out"] = append(orders["middle out"], (n/2+i)%n)
		}
		for name, order := range orders {
			t.Run(filepath.Base(file)+"/"+name, func(t *testing.T) {
				s, err := NewStream(bytes.NewReader(data), opts())
				if err != nil {
					t.Fatal(err)
				}
				if s.Pages() != n {
					t.Fatalf("Pages() = %d, want %d", s.Pages(), n)
				}
				outline := s.Outline()
				if len(outline.Parts()) != 0 {
					t.Errorf("the outline has %d parts", len(outline.Parts()))
				}
				pages := outline.Views[0].Pages
				for i, p := range pages {
					if p.W != wantPages[i].W || p.H != wantPages[i].H || p.Layers == nil || len(p.Layers) != 0 {
						t.Errorf("outline page %d: %gx%g, %d layers", i, p.W, p.H, len(p.Layers))
					}
				}
				r := &streamReader{t: t, parts: map[bdf.Hash]*bdf.Part{}, objs: map[bdf.Hash]*bdf.ObjectPart{}, cmaps: map[bdf.Hash]map[uint32]uint16{}}
				for _, i := range order {
					doc, err := s.Page(i)
					if err != nil {
						t.Fatal(err)
					}
					if len(doc.Views) != 1 || len(doc.Views[0].Pages) != 1 {
						t.Fatalf("page %d: the page document has %d views", i, len(doc.Views))
					}
					page := doc.Views[0].Pages[0]
					if page.W != pages[i].W || page.H != pages[i].H || len(page.Layers) == 0 {
						t.Fatalf("page %d: %gx%g, %d layers", i, page.W, page.H, len(page.Layers))
					}
					r.add(doc)
					for _, l := range page.Layers {
						r.check(l.Obj)
					}
					pages[i] = page
				}
				// asked again: nothing new
				again, err := s.Page(order[0])
				if err != nil || len(again.Parts()) != 0 || again.Views[0].Pages[0].Layers[0] != pages[order[0]].Layers[0] {
					t.Errorf("page %d again: %v, %d parts", order[0], err, len(again.Parts()))
				}
				res, err := s.Finish()
				if err != nil {
					t.Fatal(err)
				}
				if _, err := s.Page(0); err == nil {
					t.Error("Page after Finish")
				}
				if name == "in order" {
					if !bytes.Equal(singleBytes(t, res.Doc), singleBytes(t, want.Doc)) {
						t.Error("Finish made other bytes than Convert")
					}
					return
				}
				if len(res.Doc.Views[0].Pages) != n || res.Doc.Views[0].TextIndex == "" {
					t.Fatalf("the finished document has %d pages, text index %q", len(res.Doc.Views[0].Pages), res.Doc.Views[0].TextIndex)
				}
				finished := func(h bdf.Hash) *bdf.ObjectPart {
					o, err := bdf.DecodeObject(res.Doc.Part(h).Data)
					if err != nil {
						t.Fatal(err)
					}
					return o
				}
				expected := func(h bdf.Hash) *bdf.ObjectPart {
					o, err := bdf.DecodeObject(want.Doc.Part(h).Data)
					if err != nil {
						t.Fatal(err)
					}
					return o
				}
				for i, p := range res.Doc.Views[0].Pages {
					got, streamed, want := layersText(t, finished, p), layersText(t, r.object, pages[i]), layersText(t, expected, wantPages[i])
					if got != want || streamed != want {
						t.Errorf("page %d text:\nfinished %q\nstreamed %q\nconvert  %q", i, got, streamed, want)
					}
				}
			})
		}
	}
}

// TestStreamFontVersions checks that the font versions of a streamed
// document grow into one that the later pages share, when their text keeps
// to a few glyphs.
func TestStreamFontVersions(t *testing.T) {
	data, err := os.ReadFile("testdata/chrome-doc.pdf")
	if err != nil {
		t.Fatal(err)
	}
	s, err := NewStream(bytes.NewReader(data), &Options{NoWOFF2: true, Pages: conv.PageList(slices.Repeat([]int{1, 2}, 4)...)})
	if err != nil {
		t.Fatal(err)
	}
	fonts := 0
	for i := range s.Pages() {
		doc, err := s.Page(i)
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range doc.Parts() {
			if p.Type == bdf.PartFont {
				fonts++
			}
		}
		if i >= 2 && len(doc.Parts()) > 1 {
			// the same pages again: only the page's own object, if that
			t.Errorf("page %d repeats a page and still brings %d parts", i, len(doc.Parts()))
		}
	}
	if fonts == 0 {
		t.Fatal("no fonts returned")
	}
}
