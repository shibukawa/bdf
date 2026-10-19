// Package pdfbdf writes BDF documents as PDF files, for what a PDF is
// there for: printing, filing, sending to those who have no viewer of BDF.
//
// The instructions of the objects (docs/spec.md §7) become the operators
// of content streams, in the same order: paths, transforms, clips, lines
// and dashes, alpha and blend modes (the graphics state parameters ca, CA
// and BM), linear and radial gradients (axial and radial shadings; conic
// ones as a mesh of triangles), image patterns (tiling patterns), groups
// (transparency group forms) and soft masks (the SMask of the state a
// group is drawn with). Text stays text: the fonts of the document are
// embedded with the glyphs in use (the fonts of the system, for fonts
// referred to by name), as CID fonts with the characters of their codes,
// so that the text of the file is searched and copied; a font that may not
// be embedded, or that a PDF cannot hold, is drawn as the outlines of its
// glyphs under the text they stand for. JPEG images are stored as they
// are, other bitmaps as their pixels, SVG images as pictures of them.
// Links are link annotations, and the metadata that of the file.
//
// The pages of fixed and flow views are the pages of the file. A view of
// scroll strips is one long page, and a sheet one page of its tiles; they
// are written when the document has no pages. The layers that pages share
// (the masters of slides) are written once.
//
// A PDF has no shadows: the shadow of a shape is a picture of it (the
// shape drawn by raster/imagebdf at 300 dpi and blurred as that renderer
// blurs it), put under the shape, which stays vectors.
//
// What a PDF has no counterpart for is left out with a warning: FILTER,
// the composite operations of Canvas that are not blend modes (drawn as
// source-over), CLEAR_RECT in a group, AVIF images. Text is laid
// out as the image renderer lays it out (raster/imagebdf): not shaped, and
// with the fonts of the machine for the text of fonts referred to by name.
//
// Options.Tagged writes a tagged PDF: the structure that the MARK
// instructions of the document tell (headings, paragraphs, lists, tables,
// figures with their alternative texts, links, the language of the text)
// as the structure tree of the file, with the content of the pages marked
// by what it is of, and what is no content of the document (page headers
// and footers, masters, backgrounds, shadows) as artifacts. Options.PDFA
// makes the file a PDF/A file of part 2 or 3, of the level b, u or a: with
// an output intent of sRGB, what PDF/A does not have written another way,
// and, at level a, tagged. The files are checked with veraPDF.
package pdfbdf

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/fontdb"
)

// Options controls the writing of a file.
type Options struct {
	// Views are the ids of the views to write, in the order of the
	// document. nil writes the views of pages (fixed and flow); when the
	// document has none, all its views.
	Views []string
	// SkipRoles are the roles of the layers that are left out (bdf.RoleMaster,
	// bdf.RoleBackground …): for printing on paper that has them already.
	SkipRoles []string
	// FontFS holds fonts searched before FontDirs, for the text of fonts
	// referred to by name and for characters an embedded font lacks.
	FontFS fs.FS
	// FontDirs are searched for fonts before the system font directories.
	FontDirs []string
	// NoSystemFonts restricts font lookup to FontFS and FontDirs.
	NoSystemFonts bool
	// IgnoreFSType embeds fonts of the system whose OS/2 fsType forbids
	// embedding, and leaves glyphs out of fonts whose fsType forbids
	// subsetting. Without it the text of such fonts is drawn as outlines,
	// or the fonts are embedded whole. Set it only when you hold the rights.
	IgnoreFSType bool
	// NoCompression writes the streams of the file as they are (for
	// reading them; the file is several times larger).
	NoCompression bool
	// CompressionLevel is the level of the deflate of the streams, 1 to 9
	// (0: 6).
	CompressionLevel int
	// Tagged writes a tagged PDF: the structure of the document (its
	// headings, paragraphs, lists, tables and figures with their
	// alternative texts, as the MARKs of its objects tell them) as the
	// structure tree of the file, with the content of the pages marked as
	// that of its elements, and the rest as artifacts. Screen readers and
	// the tools that take the text of a file out go by it.
	Tagged bool
	// Pages selects the pages written, by their zero-based number among
	// the pages the file has without it, in that order; numbers past
	// them are left out. nil writes all of them. For a printer that
	// takes no page-ranges.
	Pages []int
	// PDFA makes the file a PDF/A file, of the part and the level named:
	// "2b" (it shows the same in any reader), "2u" (and its text is
	// Unicode), "2a" (and it is tagged: Tagged is set), or "3b", "3u" and
	// "3a", the same of part 3. "" writes a PDF that names no standard.
	PDFA string
}

// Result tells what was written.
type Result struct {
	Pages int
	// Fonts is the number of fonts embedded.
	Fonts    int
	Warnings []string
}

// maxPageSide is the largest side of a page in the units of a PDF without
// a larger unit (200 inches).
const maxPageSide = 14400

// ErrNoPages is returned for a document without anything to write.
var ErrNoPages = errors.New("pdfbdf: the document has no pages")

type exporter struct {
	doc  *bdf.Document
	opts Options
	pw   *writer

	// tags makes the structure tree of a tagged file, and pdfa is the
	// PDF/A that the file is (nil: none)
	tags *tagger
	pdfa *pdfaLevel

	objects  map[bdf.Hash]*bdf.ObjectPart
	coll     map[bdf.Hash][]*bdf.Path
	images   map[bdf.Hash]*pdfImage
	shadings map[*bdf.Paint]*shading
	patterns map[patKey]named
	gstates  map[string]named
	forms    map[layerKey]*layerForm
	fonts    *fonts
	shadows  *shadows
	// pageW and pageH are the size of the page being written (of the
	// first page of a layer that pages share), in units of the document
	pageW, pageH float64
	names        int
	warnings     []string
	warned       map[string]bool

	// the stream being written, the object whose instructions are run,
	// and the groups and soft masks that are open
	s      *stream
	obj    *bdf.ObjectPart
	groups []*group
	masks  []*openMask
	depth  int
	alt    *string
	// floor is the number of states of the stream that RESTORE does not
	// take (those from before a USE), unsaved the SAVEs that kept no
	// state, and floors what they were around the open groups and masks
	floor, unsaved int
	floors         []int
	skipped        int
	// seen holds the objects written before, again says that the object
	// being written is one, and reused counts the instructions of such
	seen   map[bdf.Hash]struct{}
	again  bool
	reused int
}

// layerKey is a layer on pages of a size: the object of the layer, and
// the transform of such pages (see begin).
type layerKey struct {
	obj  bdf.Hash
	flip matrix
}

// layerForm is a layer that pages share, as a form XObject.
type layerForm struct {
	name  string
	ref   ref
	links []link
}

// pageOut is a page of the file.
type pageOut struct {
	// key is the key of the content of the page in the parent tree of a
	// tagged file (-1 without marked content)
	key     int
	ref     ref
	w, h    float64 // in units of the document
	unit    float64 // the units of the document in a unit of the page
	flip    matrix  // from the space of the document to that of the page
	content ref
	res     string
	links   []link
	view    *bdf.View
	index   int // in its view
}

func (ex *exporter) warnf(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	if !ex.warned[s] {
		ex.warned[s] = true
		ex.warnings = append(ex.warnings, s)
	}
}

func (ex *exporter) object(h bdf.Hash) *bdf.ObjectPart {
	if o, ok := ex.objects[h]; ok {
		return o
	}
	ex.objects[h] = nil
	p := ex.doc.Part(h)
	if p == nil {
		ex.warnf("object %s is missing", h)
		return nil
	}
	o, err := bdf.DecodeObject(p.Data)
	if err != nil {
		ex.warnf("object %s: %v", h, err)
		return nil
	}
	ex.objects[h] = o
	return o
}

// Write writes a document as a PDF file.
func Write(w io.Writer, doc *bdf.Document, opts *Options) (res *Result, err error) {
	defer func() {
		if p := recover(); p != nil {
			res, err = nil, fmt.Errorf("pdfbdf: the writing failed: %v", p)
		}
	}()
	if doc == nil {
		return nil, errors.New("pdfbdf: no document")
	}
	ex := &exporter{doc: doc, objects: map[bdf.Hash]*bdf.ObjectPart{}, coll: map[bdf.Hash][]*bdf.Path{}, images: map[bdf.Hash]*pdfImage{},
		shadings: map[*bdf.Paint]*shading{}, patterns: map[patKey]named{}, gstates: map[string]named{}, forms: map[layerKey]*layerForm{},
		warned: map[string]bool{}, seen: map[bdf.Hash]struct{}{}}
	if opts != nil {
		ex.opts = *opts
	}
	ex.fonts = &fonts{ex: ex, parts: map[bdf.Hash]*face{}, files: map[*fontdb.Face]*face{}, chains: map[bdf.Font]*fontChain{}}
	if ex.pdfa, err = parsePDFA(ex.opts.PDFA); err != nil {
		return nil, err
	}
	if ex.opts.Tagged || ex.pdfa != nil && ex.pdfa.conf == "A" {
		ex.tags = newTagger(ex)
	}
	views, err := ex.views()
	if err != nil {
		return nil, err
	}
	level := ex.opts.CompressionLevel
	if level < 1 || level > 9 {
		level = 6
	}
	ex.pw = newWriter(w, !ex.opts.NoCompression, level)
	pw := ex.pw
	catalog, tree := pw.alloc(), pw.alloc()

	// the pages: their objects first, which links name
	var pages []*pageOut
	for _, v := range views {
		switch v.Kind {
		case bdf.ViewFixed, bdf.ViewFlow:
			for i, p := range v.Pages {
				if p == nil || !(p.W > 0) || !(p.H > 0) || math.IsInf(float64(p.W), 0) || math.IsInf(float64(p.H), 0) {
					continue
				}
				pages = append(pages, &pageOut{key: -1, w: float64(p.W), h: float64(p.H), view: v, index: i})
			}
		case bdf.ViewScroll:
			if w, h := scrollSize(v); w > 0 && h > 0 {
				pages = append(pages, &pageOut{key: -1, w: w, h: h, view: v, index: -1})
			}
		case bdf.ViewSheet:
			if w, h := ex.sheetSize(v); w > 0 && h > 0 {
				pages = append(pages, &pageOut{key: -1, w: w, h: h, view: v, index: -1})
			}
		}
	}
	if ex.opts.Pages != nil {
		var chosen []*pageOut
		for _, n := range ex.opts.Pages {
			if n >= 0 && n < len(pages) {
				chosen = append(chosen, pages[n])
			}
		}
		pages = chosen
	}
	if len(pages) == 0 {
		return nil, ErrNoPages
	}
	for _, pg := range pages {
		pg.ref = pw.alloc()
	}
	ex.shared(pages)
	for _, pg := range pages {
		ex.page(pg)
	}
	// where the links to pages go
	place := map[*bdf.View][]*pageOut{}
	for _, pg := range pages {
		place[pg.view] = append(place[pg.view], pg)
	}
	var kids strings.Builder
	for _, pg := range pages {
		kids.WriteString(pg.ref.String() + " ")
		k := 1 / pg.unit
		var annots strings.Builder
		for _, l := range pg.links {
			action := ex.linkAction(l.url, pg, place)
			if action == "" {
				continue
			}
			// (the page has its y axis up)
			rect := nums(l.x0*k, (pg.h-l.y1)*k, l.x1*k, (pg.h-l.y0)*k)
			tagged := ""
			annot := pw.alloc()
			if t := ex.tags; t != nil {
				// what the link is, for who does not see where it is
				tagged = " /Contents " + text(l.url)
				if l.elem != nil {
					tagged += fmt.Sprintf(" /StructParent %d", t.next)
					t.annots = append(t.annots, annotTag{key: t.next, elem: l.elem})
					t.next++
					if l.elem.page == nil {
						l.elem.page = pg
					}
					l.elem.kids = append(l.elem.kids, structKid{annot: &annot, page: pg})
				}
			}
			pw.put(annot, fmt.Sprintf("<< /Type /Annot /Subtype /Link /Rect [%s] /Border [0 0 0] /F 4%s %s >>", rect, tagged, action))
			annots.WriteString(annot.String() + " ")
		}
		dict := fmt.Sprintf("<< /Type /Page /Parent %s /MediaBox [0 0 %s] /Resources %s /Contents %s", tree, nums(pg.w*k, pg.h*k), pg.res, pg.content)
		if pg.unit != 1 {
			dict += " /UserUnit " + num(pg.unit)
		}
		if ex.tags != nil {
			// (the keys go through the page in the order of its structure)
			dict += " /Tabs /S"
			if pg.key >= 0 {
				dict += fmt.Sprintf(" /StructParents %d", pg.key)
			}
		}
		if annots.Len() > 0 {
			dict += " /Annots [" + annots.String() + "]"
		}
		pw.put(pg.ref, dict+" >>")
	}
	pw.put(tree, fmt.Sprintf("<< /Type /Pages /Count %d /Kids [%s] >>", len(pages), kids.String()))
	nfonts := ex.writeFonts()

	meta := ex.metadata()
	cat := "<< /Type /Catalog /Pages " + tree.String() + " /Metadata " + meta.String()
	if ex.tags != nil {
		cat += ex.tags.write()
	}
	if ex.pdfa != nil {
		cat += ex.outputIntent()
	}
	if lang := doc.Meta.DC.Language.First(); lang != "" {
		cat += " /Lang " + text(lang)
	}
	prefs := ""
	if doc.Meta.DC.Title.First() != "" {
		prefs += " /DisplayDocTitle true"
	}
	if views[0].Direction == bdf.DirectionRTL {
		prefs += " /Direction /R2L"
	}
	if prefs != "" {
		cat += " /ViewerPreferences <<" + prefs + " >>"
	}
	pw.put(catalog, cat+" >>")
	if err := pw.finish(catalog, ex.info()); err != nil {
		return nil, fmt.Errorf("pdfbdf: %w", err)
	}
	return &Result{Pages: len(pages), Fonts: nfonts, Warnings: ex.warnings}, nil
}

// views returns the views to write.
func (ex *exporter) views() ([]*bdf.View, error) {
	doc := ex.doc
	if len(ex.opts.Views) > 0 {
		var out []*bdf.View
		for _, id := range ex.opts.Views {
			found := false
			for _, v := range doc.Views {
				if v != nil && v.ID == id {
					out, found = append(out, v), true
				}
			}
			if !found {
				return nil, fmt.Errorf("pdfbdf: no view %q", id)
			}
		}
		return out, nil
	}
	var paged, all []*bdf.View
	for _, v := range doc.Views {
		if v == nil {
			continue
		}
		all = append(all, v)
		if v.Kind == bdf.ViewFixed || v.Kind == bdf.ViewFlow {
			paged = append(paged, v)
		}
	}
	if len(paged) > 0 {
		return paged, nil
	}
	if len(all) == 0 {
		return nil, ErrNoPages
	}
	return all, nil
}

func (ex *exporter) skip(role string) bool {
	for _, r := range ex.opts.SkipRoles {
		if r == role {
			return true
		}
	}
	return false
}

// shared makes forms of the layers that more than one page has. The form
// of a layer is drawn in the space of the page of the file (see
// stream.base): it holds the transform of the pages of its size.
func (ex *exporter) shared(pages []*pageOut) {
	count := map[layerKey]int{}
	for _, pg := range pages {
		ex.size(pg)
		if pg.index < 0 {
			continue
		}
		for _, l := range pg.view.Pages[pg.index].Layers {
			// (in a tagged file the content of a page is marked as that
			// of the page: a layer of content is written for each page)
			if !ex.skip(l.Role) && (ex.tags == nil || artifactRole(l.Role)) {
				count[layerKey{l.Obj, pg.flip}]++
			}
		}
	}
	for _, pg := range pages {
		if pg.index < 0 {
			continue
		}
		for _, l := range pg.view.Pages[pg.index].Layers {
			key := layerKey{l.Obj, pg.flip}
			if count[key] < 2 || ex.forms[key] != nil || ex.skip(l.Role) {
				continue
			}
			s := newStream(identity, pg.flip)
			s.op("%s cm", pg.flip.exact())
			ex.s, ex.pageW, ex.pageH = s, pg.w, pg.h
			if ex.tags != nil {
				ex.tags.off++ // (the form is an artifact as a whole, where it is drawn)
			}
			ex.top(l.Obj)
			if ex.tags != nil {
				ex.tags.off--
			}
			ex.forms[key] = &layerForm{name: ex.newName("Fm"), ref: ex.form(s, [4]float64{-huge, -huge, huge, huge}, formPlain), links: s.links}
		}
	}
}

// artifactRole reports whether the layers of a role are not the content of
// a page: what is under it and around it on every page.
func artifactRole(role string) bool {
	switch role {
	case bdf.RoleBackground, bdf.RoleMaster, bdf.RoleHeader, bdf.RoleFooter:
		return true
	}
	return false
}

// pageLayer writes a layer of a page: as an artifact, when it is not content.
func (ex *exporter) pageLayer(role string, h bdf.Hash) {
	t := ex.tags
	if t == nil || !artifactRole(role) {
		ex.top(h)
		return
	}
	t.artifact++
	ex.top(h)
	t.artifact--
}

// size gives a page its unit, when it is too large for the units of a
// PDF, and the transform from the space of the document, whose y axis
// goes down.
func (ex *exporter) size(pg *pageOut) {
	pg.unit = 1
	if side := math.Max(pg.w, pg.h); side > maxPageSide {
		pg.unit = math.Ceil(side / maxPageSide)
	}
	k := 1 / pg.unit
	pg.flip = matrix{k, 0, 0, -k, 0, pg.h * k}
}

// page writes the content of a page.
func (ex *exporter) page(pg *pageOut) {
	s := newStream(identity, pg.flip)
	s.page = true
	s.op("%s cm", pg.flip.exact())
	ex.s, ex.pageW, ex.pageH = s, pg.w, pg.h
	if ex.tags != nil {
		ex.tags.page = pg
	}
	back, _ := pg.flip.invert()
	v := pg.view
	switch {
	case pg.index >= 0:
		for _, l := range v.Pages[pg.index].Layers {
			if ex.skip(l.Role) {
				continue
			}
			if f := ex.forms[layerKey{l.Obj, pg.flip}]; f != nil {
				ex.tagArtifact()
				s.op("q %s cm %s Do Q", back.exact(), s.use("XObject", f.name, f.ref))
				s.links = append(s.links, f.links...)
				continue
			}
			ex.pageLayer(l.Role, l.Obj)
		}
	case v.Kind == bdf.ViewScroll:
		ex.scroll(v)
	default:
		ex.sheet(v, pg)
	}
	s.finish()
	pg.content = ex.pw.addStream("", s.b.Bytes(), true)
	pg.res, pg.links, pg.key = s.resources(), s.links, s.key
}

// scrollSize is the size of the strips of a scroll view stacked.
func scrollSize(v *bdf.View) (w, h float64) {
	for _, p := range v.Pages {
		if p == nil {
			continue
		}
		b := bdf.RectDef{W: p.W, H: p.H}
		if p.Body != nil {
			b = *p.Body
		}
		if !(b.W > 0) || !(b.H > 0) {
			continue
		}
		w = math.Max(w, float64(b.W))
		h += float64(b.H)
	}
	if math.IsInf(w, 0) || math.IsInf(h, 0) {
		return 0, 0
	}
	return w, h
}

// scroll writes the strips of a scroll view stacked without gaps, each cut
// to its body, as the viewer scrolls them.
func (ex *exporter) scroll(v *bdf.View) {
	s := ex.s
	y := 0.0
	for _, p := range v.Pages {
		if p == nil {
			continue
		}
		b := bdf.RectDef{W: p.W, H: p.H}
		if p.Body != nil {
			b = *p.Body
		}
		if !(b.W > 0) || !(b.H > 0) {
			continue
		}
		ex.scoped(func() {
			saved := s.st.ctm
			s.op("%s re W n", nums(0, y, float64(b.W), float64(b.H)))
			ex.concat(matrix{1, 0, 0, 1, -float64(b.X), y - float64(b.Y)})
			for _, l := range p.Layers {
				if l.Role == bdf.RoleBody || l.Role == bdf.RoleAnnotation {
					ex.pageLayer(l.Role, l.Obj)
				}
			}
			s.st.ctm = saved
		})
		y += float64(b.H)
	}
}

// maxTiles bounds the tiles of a sheet that are written.
const maxTiles = 1 << 16

// tilesOf returns the tiles of a sheet with their places, in the order of
// the rows.
func tilesOf(v *bdf.View) (out []tilePlace) {
	for k, hs := range v.Tiles {
		x, y, _ := strings.Cut(k, ",")
		tx, errx := strconv.Atoi(x)
		ty, erry := strconv.Atoi(y)
		if errx != nil || erry != nil || tx < 0 || ty < 0 || tx > 1<<20 || ty > 1<<20 || strconv.Itoa(tx)+","+strconv.Itoa(ty) != k {
			continue
		}
		h, err := bdf.ParseHash(hs)
		if err != nil {
			continue
		}
		out = append(out, tilePlace{tx, ty, h})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ty < out[j].ty || out[i].ty == out[j].ty && out[i].tx < out[j].tx })
	if len(out) > maxTiles {
		out = out[:maxTiles]
	}
	return out
}

type tilePlace struct {
	tx, ty int
	hash   bdf.Hash
}

func tileSize(v *bdf.View) float64 {
	if v.Tile > 0 && !math.IsInf(float64(v.Tile), 0) {
		return float64(v.Tile)
	}
	return 2048
}

// sheetSize is the size of the page of a sheet: what its tiles reach,
// within its rows and columns.
func (ex *exporter) sheetSize(v *bdf.View) (w, h float64) {
	tile := tileSize(v)
	if tile < 1 {
		return 0, 0
	}
	for _, t := range tilesOf(v) {
		w, h = math.Max(w, float64(t.tx+1)*tile), math.Max(h, float64(t.ty+1)*tile)
	}
	sw, sh := 0.0, 0.0
	for _, c := range v.Cols {
		sw += float64(c[0]) * float64(c[1])
	}
	for _, r := range v.Rows {
		sh += float64(r[0]) * float64(r[1])
	}
	if sw > 0 && !math.IsInf(sw, 0) {
		w = math.Min(w, sw)
	}
	if sh > 0 && !math.IsInf(sh, 0) {
		h = math.Min(h, sh)
	}
	return w, h
}

// maxGridLines bounds the gridlines of a sheet that are written.
const maxGridLines = 1 << 14

// sheet writes a sheet: its gridlines, when the view asks for them, and
// its tiles.
func (ex *exporter) sheet(v *bdf.View, pg *pageOut) {
	s := ex.s
	tile := tileSize(v)
	if v.Gridlines {
		// the lines between the cells, as the viewer draws them
		ex.tagArtifact()
		s.op("q 0.851 0.851 0.851 RG %s w", num(0.5*pg.unit))
		lines := 0
		walk := func(runs []bdf.Run, limit float64, draw func(pos float64)) {
			pos := 0.0
			for _, run := range runs {
				n, size := float64(run[0]), float64(run[1])
				if !(n >= 1) || !(size > 0) {
					continue
				}
				for i := 0; float64(i) < math.Floor(n); i++ {
					if pos += size; pos > limit || lines >= maxGridLines {
						return
					}
					lines++
					draw(pos)
				}
			}
		}
		walk(v.Cols, pg.w, func(x float64) { s.op("%s m %s l", nums(x, 0), nums(x, pg.h)) })
		walk(v.Rows, pg.h, func(y float64) { s.op("%s m %s l", nums(0, y), nums(pg.w, y)) })
		s.op("S Q")
		if lines >= maxGridLines {
			ex.warnf("a sheet has more than %d gridlines: the rest are not written", maxGridLines)
		}
	}
	for _, t := range tilesOf(v) {
		ex.scoped(func() {
			saved := s.st.ctm
			ex.concat(matrix{1, 0, 0, 1, float64(t.tx) * tile, float64(t.ty) * tile})
			s.op("0 0 %s re W n", nums(tile, tile))
			ex.top(t.hash)
			s.st.ctm = saved
		})
	}
}

// linkAction writes where a link goes: the entry of its annotation ("" for
// a link that goes nowhere in the file).
func (ex *exporter) linkAction(target string, from *pageOut, place map[*bdf.View][]*pageOut) string {
	dest := func(v *bdf.View, n int) string {
		pages := place[v]
		if len(pages) == 0 {
			return ""
		}
		if pages[0].index < 0 {
			// (a view on one page: the page)
			return fmt.Sprintf("/Dest [%s /Fit]", pages[0].ref)
		}
		for _, pg := range pages {
			if pg.index == n-1 {
				return fmt.Sprintf("/Dest [%s /Fit]", pg.ref)
			}
		}
		return ""
	}
	if frag, ok := strings.CutPrefix(target, "#"); ok {
		q, err := url.ParseQuery(frag)
		if err != nil {
			return ""
		}
		n := 1
		if p := q.Get("page"); p != "" {
			if n, err = strconv.Atoi(p); err != nil {
				return ""
			}
		}
		if id := q.Get("view"); id != "" {
			for v := range place {
				if v.ID == id {
					return dest(v, n)
				}
			}
			return ""
		}
		if q.Has("page") {
			return dest(from.view, n)
		}
		return ""
	}
	u, err := url.Parse(target)
	if err != nil || u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "mailto" {
		return ""
	}
	// (a URI is ASCII: what is not is escaped)
	var b strings.Builder
	for i := 0; i < len(target); i++ {
		if c := target[i]; c <= 0x20 || c >= 0x7f {
			fmt.Fprintf(&b, "%%%02X", c)
		} else {
			b.WriteByte(c)
		}
	}
	return "/A << /Type /Action /S /URI /URI " + text(b.String()) + " >>"
}

// pdfDate writes a date of the metadata (ISO 8601) as a date of a PDF
// ("" for one that is not a date).
func pdfDate(s string) (pdf, iso string) {
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05", "2006-01-02T15:04Z07:00", "2006-01-02", "2006-01", "2006"} {
		t, err := time.Parse(layout, strings.TrimSpace(s))
		if err != nil {
			continue
		}
		_, off := t.Zone()
		zone := "Z"
		if off != 0 {
			sign := '+'
			if off < 0 {
				sign, off = '-', -off
			}
			zone = fmt.Sprintf("%c%02d'%02d'", sign, off/3600, off/60%60)
		}
		return "D:" + t.Format("20060102150405") + zone, t.Format(time.RFC3339)
	}
	return "", ""
}

// producer names this writer in the metadata of the file.
const producer = "BDF (github.com/shibukawa/bdf)"

// info writes the document information dictionary.
func (ex *exporter) info() ref {
	dc := ex.doc.Meta.DC
	var b strings.Builder
	b.WriteString("<<")
	add := func(key, v string) {
		if v != "" {
			b.WriteString(" /" + key + " " + text(v))
		}
	}
	add("Title", dc.Title.First())
	add("Author", strings.Join(dc.Creator, ", "))
	add("Subject", dc.Description.First())
	add("Keywords", strings.Join(dc.Subject, ", "))
	add("Producer", producer)
	if d, _ := pdfDate(dc.Created.First()); d != "" {
		add("CreationDate", d)
	}
	if d, _ := pdfDate(dc.Modified.First()); d != "" {
		add("ModDate", d)
	}
	b.WriteString(" >>")
	return ex.pw.add(b.String())
}

// metadata writes the metadata of the document as an XMP packet.
func (ex *exporter) metadata() ref {
	dc := ex.doc.Meta.DC
	esc := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\"", "&quot;")
	clean := func(s string) string {
		// (what XML has no place for)
		return esc.Replace(strings.Map(func(r rune) rune {
			if r < 0x20 && r != '\t' && r != '\n' && r != '\r' || r == 0xfffe || r == 0xffff {
				return -1
			}
			return r
		}, strings.ToValidUTF8(s, "")))
	}
	var b strings.Builder
	b.WriteString("<?xpacket begin=\"\ufeff\" id=\"W5M0MpCehiHzreSzNTczkc9d\"?>\n<x:xmpmeta xmlns:x=\"adobe:ns:meta/\">\n" +
		"<rdf:RDF xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\">\n<rdf:Description rdf:about=\"\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" " +
		"xmlns:xmp=\"http://ns.adobe.com/xap/1.0/\" xmlns:pdf=\"http://ns.adobe.com/pdf/1.3/\">\n")
	list := func(tag, kind string, vs []string, lang bool) {
		if len(vs) == 0 {
			return
		}
		b.WriteString("<" + tag + "><rdf:" + kind + ">")
		for _, v := range vs {
			if lang {
				b.WriteString("<rdf:li xml:lang=\"x-default\">" + clean(v) + "</rdf:li>")
			} else {
				b.WriteString("<rdf:li>" + clean(v) + "</rdf:li>")
			}
		}
		b.WriteString("</rdf:" + kind + "></" + tag + ">\n")
	}
	list("dc:title", "Alt", dc.Title[:min(len(dc.Title), 1)], true)
	list("dc:creator", "Seq", dc.Creator, false)
	list("dc:description", "Alt", dc.Description[:min(len(dc.Description), 1)], true)
	list("dc:subject", "Bag", dc.Subject, false)
	list("dc:publisher", "Bag", dc.Publisher, false)
	list("dc:language", "Bag", dc.Language, false)
	list("dc:rights", "Alt", dc.Rights[:min(len(dc.Rights), 1)], true)
	b.WriteString("<dc:format>application/pdf</dc:format>\n")
	if _, d := pdfDate(dc.Created.First()); d != "" {
		b.WriteString("<xmp:CreateDate>" + d + "</xmp:CreateDate>\n")
	}
	if _, d := pdfDate(dc.Modified.First()); d != "" {
		b.WriteString("<xmp:ModifyDate>" + d + "</xmp:ModifyDate>\n")
	}
	b.WriteString("<pdf:Producer>" + clean(producer) + "</pdf:Producer>\n</rdf:Description>\n")
	if a := ex.pdfa; a != nil {
		// which PDF/A the file is
		fmt.Fprintf(&b, "<rdf:Description rdf:about=\"\" xmlns:pdfaid=\"http://www.aiim.org/pdfa/ns/id/\">\n<pdfaid:part>%d</pdfaid:part>\n<pdfaid:conformance>%s</pdfaid:conformance>\n</rdf:Description>\n", a.part, a.conf)
	}
	b.WriteString("</rdf:RDF>\n</x:xmpmeta>\n<?xpacket end=\"w\"?>")
	// (not deflated: tools read the metadata of a file without its streams)
	return ex.pw.addStream(" /Type /Metadata /Subtype /XML", []byte(b.String()), false)
}
