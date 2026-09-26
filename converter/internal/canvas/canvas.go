// Package canvas builds the objects of a document whose text the converter
// lays out itself. Font references stay placeholders until the whole
// document has been laid out, because embedded font programs are subsets of
// every character used anywhere (see fontset); child objects (vertical text
// columns) are encoded before their parents.
package canvas

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/fontset"
)

// Builder makes the canvases of one document.
type Builder struct {
	doc      *bdf.Document
	fonts    *fontset.Set
	canvases []*Canvas
}

// NewBuilder returns a builder of objects for doc whose text is measured
// with fonts.
func NewBuilder(doc *bdf.Document, fonts *fontset.Set) *Builder {
	return &Builder{doc: doc, fonts: fonts}
}

// New starts an object.
func (b *Builder) New() *Canvas {
	cv := &Canvas{b: b, Obj: bdf.NewObject(), fontIdx: map[fontset.Use]bdf.FontRef{}, images: map[bdf.Hash]bdf.ImageRef{}}
	b.canvases = append(b.canvases, cv)
	return cv
}

// Encode resolves the font references of every canvas with the fonts'
// records and adds the objects to the document. Call it after the fonts
// have been embedded (fontset.Set.Embed).
func (b *Builder) Encode() {
	for _, cv := range b.canvases {
		cv.encode()
	}
}

// Canvas is an object under construction.
type Canvas struct {
	b        *Builder
	Obj      *bdf.Object
	fonts    []fontset.Use
	fontIdx  map[fontset.Use]bdf.FontRef
	images   map[bdf.Hash]bdf.ImageRef
	children []*Canvas
	// Drawn records that something visible was drawn.
	Drawn bool
	// lang is the language in effect at the end of the instructions so far
	// ("" = the document's): LANG marks carry over in drawing order, and
	// each top-level object starts with the document's language.
	lang string

	encoded bool
	hash    bdf.Hash
}

// Font returns the reference of a font in the object.
func (cv *Canvas) Font(u fontset.Use) bdf.FontRef {
	if r, ok := cv.fontIdx[u]; ok {
		return r
	}
	r := cv.Obj.AddFont(bdf.SystemFont("\x00pending", 0, 0))
	cv.fonts = append(cv.fonts, u)
	cv.fontIdx[u] = r
	return r
}

// Image returns the reference of a document image in the object.
func (cv *Canvas) Image(h bdf.Hash) bdf.ImageRef {
	if r, ok := cv.images[h]; ok {
		return r
	}
	r := cv.Obj.AddImage(h)
	cv.images[h] = r
	return r
}

// Lang returns the language in effect after the instructions so far ("" for
// the document's).
func (cv *Canvas) Lang() string { return cv.lang }

// SetLang makes lang (a language tag, "" for the document's) the language
// of the text that follows, with a LANG mark when it changes.
func (cv *Canvas) SetLang(lang string) {
	lang = NormLang(lang)
	if isDocLang(lang, cv.b.doc.Meta.DC.Language.First()) {
		lang = ""
	}
	if lang != cv.lang {
		cv.Obj.Mark(bdf.MarkLang, lang)
		cv.lang = lang
	}
}

// NormLang cleans up a language tag from the markup; "" means unknown.
func NormLang(l string) string {
	l = strings.TrimSpace(l)
	if strings.EqualFold(l, "x-none") {
		return ""
	}
	return l
}

// isDocLang reports whether lang stands for the document's language doc:
// the same tag, or its language subtag alone (as the script of East Asian
// text implies: "ja" in a "ja-JP" document).
func isDocLang(lang, doc string) bool {
	if strings.EqualFold(lang, doc) {
		return true
	}
	n := len(lang)
	return n > 0 && !strings.Contains(lang, "-") && len(doc) > n && doc[n] == '-' && strings.EqualFold(doc[:n], lang)
}

// Child adds a child object and returns its canvas and reference.
func (cv *Canvas) Child(bbox bdf.Rect) (*Canvas, bdf.ObjRef) {
	ch := cv.b.New()
	ref := cv.Obj.AddObject(bdf.Hash{}, bbox)
	cv.children = append(cv.children, ch)
	return ch, ref
}

// Share makes ch, a canvas of the same builder, a child of cv as well and
// returns its reference: an object that several parents draw (a picture
// over several tiles of a sheet) is stored once. bbox is ch's bounding box.
func (cv *Canvas) Share(ch *Canvas, bbox bdf.Rect) bdf.ObjRef {
	ref := cv.Obj.AddObject(bdf.Hash{}, bbox)
	cv.children = append(cv.children, ch)
	return ref
}

// Used records that cv has just used ch: structure and languages carry on
// through a USE (spec §7.8), so the language in effect is now the one ch
// ended with.
func (cv *Canvas) Used(ch *Canvas) { cv.lang = ch.lang }

// Transform emits m (as TRANSLATE when it is one).
func (cv *Canvas) Transform(m Matrix) {
	if m[0] == 1 && m[1] == 0 && m[2] == 0 && m[3] == 1 {
		if m[4] != 0 || m[5] != 0 {
			cv.Obj.Translate(f32(m[4]), f32(m[5]))
		}
		return
	}
	cv.Obj.Transform(f32(m[0]), f32(m[1]), f32(m[2]), f32(m[3]), f32(m[4]), f32(m[5]))
}

// Hash returns the hash of the encoded object (after Builder.Encode).
func (cv *Canvas) Hash() bdf.Hash { return cv.hash }

func (cv *Canvas) encode() bdf.Hash {
	if cv.encoded {
		return cv.hash
	}
	cv.encoded = true
	for i, u := range cv.fonts {
		cv.Obj.UpdateFont(bdf.FontRef(i), cv.b.fonts.Font(u))
	}
	for i, ch := range cv.children {
		cv.Obj.UpdateObject(bdf.ObjRef(i), ch.encode())
	}
	cv.hash, _ = cv.b.doc.AddObject(cv.Obj)
	return cv.hash
}

func f32(v float64) float32 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return float32(v)
}
