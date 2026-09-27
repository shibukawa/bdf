// Package kicad converts KiCad schematics (.kicad_sch) and boards
// (.kicad_pcb) of KiCad 6 and later, and zip archives of their projects,
// into BDF documents.
//
// A schematic becomes a view whose pages are its sheets, in the order of
// their page numbers: a hierarchical sheet used twice is two pages, with
// the reference designators of each instance. The sheets a schematic
// refers to are read from the files given with it (Options.Files, the
// other files of a zip archive) or from its directory. A board becomes a
// view of its front and one of its back, drawn as KiCad's board editor
// shows them, and a view of each layer. Colors are those of KiCad's
// default theme; text is drawn with NewStroke, the stroke font KiCad draws
// text with (CJK characters, which it lacks, with TrueType fonts), or with
// the TrueType font a text names. See docs/design.md §3.29.
package kicad

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/internal/fontdb"
)

// Options controls the conversion.
type Options struct {
	// Pages selects 1-based sheets of the schematic in page order (see
	// converter.Pages); nil converts all of them.
	Pages converter.Pages
	// Title overrides the document title.
	Title string
	// Views selects what a project converts: "all" (the default: the
	// schematic and the board), "schematic" or "board". A board converts
	// into its front, back and layer views unless Layers is "board" (the
	// front and back only) or "layers" (the layer views only).
	Views  string
	Layers string
	// FileName is the input's file name: it names the project of a
	// schematic without a project file.
	FileName string
	// Refs reads the files an input refers to (the sheets of a
	// hierarchical schematic, the project file): converter.Options.ReadRef
	// with the input's Files and Dir. nil reads none.
	Refs func(from, rel string) ([]byte, string, error)
	// FontFS, FontDirs and NoSystemFonts find the TrueType fonts of CJK
	// text and of texts that name a font (see converter.Options).
	FontFS        fs.FS
	FontDirs      []string
	NoSystemFonts bool
	// SystemFonts refers to fonts by family name instead of embedding them.
	SystemFonts bool
	// NoSubset embeds whole fonts instead of the glyphs in use.
	NoSubset bool
	// NoWOFF2 stores embedded fonts as TrueType/OpenType instead of WOFF2.
	NoWOFF2 bool
	// IgnoreFSType embeds fonts whose OS/2 fsType forbids it.
	IgnoreFSType bool
	// NoTextIndex skips building the text index parts.
	NoTextIndex bool
	// Warn receives non-fatal problems; when nil they are collected in
	// Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc           *bdf.Document
	Warnings      []string
	Sheets        int // pages of the schematic
	Boards        int // boards converted
	EmbeddedFonts int
}

// maxSize bounds the files read.
const maxSize = 256 << 20

// conv is the state of a conversion.
type conv struct {
	opts       *Options
	doc        *bdf.Document
	fonts      *cad.Fonts
	cvs        *canvas.Builder
	plotter    *cad.Plotter
	warnings   []string
	warned     map[string]bool
	iu         float64 // internal units per mm of the file being drawn
	defaultPen float64 // the pen of text and lines without a width (mm)
	refs       func(from, rel string) ([]byte, string, error)
	textVars   map[string]string // the project's text variables
	project    string            // the project name
	// schSheetFile is the drawing sheet the project names for its
	// schematic, and schSheet the one read
	schSheetFile string
	schSheet     *worksheet
	sch          schSettings            // the project's schematic drawing settings
	read         int                    // bytes read from the files the input refers to
	drawn        int                    // the work of the pages drawn (see maxDrawn)
	refCache     map[string]refResult   // the files read, by reference
	bitmaps      map[*node]*sheetBitmap // the drawing sheet's images, decoded
}

// refResult is a file read by reference.
type refResult struct {
	b    []byte
	path string
	err  error
}

// maxRead bounds the bytes read from the files an input refers to, in all.
const maxRead = 1 << 30

// maxDrawn bounds the work of drawing a schematic's pages, counted as the
// bytes of the files each page draws plus 100 an item drawn: a sheet used
// in many places, or a drawing sheet that repeats much, is drawn on each
// page. The pages past it are left blank.
const maxDrawn = 1 << 30

// readRef reads a file the input refers to (see Options.Refs), once, and
// within maxRead in all.
func (c *conv) readRef(from, rel string) ([]byte, string, error) {
	if c.refs == nil {
		return nil, "", fmt.Errorf("%s: %w", rel, fs.ErrNotExist)
	}
	key := path.Clean(path.Join(path.Dir(from), strings.ReplaceAll(rel, `\`, "/")))
	if r, ok := c.refCache[key]; ok {
		return r.b, r.path, r.err
	}
	b, p, err := c.refs(from, rel)
	if err == nil {
		c.read += len(b)
		if c.read > maxRead {
			b, p, err = nil, "", fmt.Errorf("%s: the files read are larger than %d bytes in all", rel, maxRead)
		}
	}
	if c.refCache == nil {
		c.refCache = map[string]refResult{}
	}
	c.refCache[key] = refResult{b, p, err}
	return b, p, err
}

func (c *conv) warnf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if c.opts.Warn != nil {
		c.opts.Warn(msg)
		return
	}
	c.warnings = append(c.warnings, msg)
}

// warnOnce warns about a kind of problem once per conversion.
func (c *conv) warnOnce(key, format string, args ...any) {
	if c.warned[key] {
		return
	}
	c.warned[key] = true
	c.warnf(format, args...)
}

// input is what an input holds: a schematic, a board or both, by their
// paths (relative to the input's directory, or in the zip archive).
type input struct {
	sch, pcb, pro string
	schData       []byte
	pcbData       []byte
	proData       []byte
}

// Convert converts a schematic, a board, a project file or a zip archive of
// a project read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if size > maxSize {
		return nil, fmt.Errorf("kicad: the file is larger than %d bytes", maxSize)
	}
	data := make([]byte, size)
	if _, err := r.ReadAt(data, 0); err != nil && err != io.EOF {
		return nil, fmt.Errorf("kicad: %w", err)
	}
	c := &conv{opts: opts, doc: bdf.NewDocument(), warned: map[string]bool{}, refs: opts.Refs, sch: defaultSchSettings}
	in, err := c.readInput(data)
	if err != nil {
		return nil, err
	}
	c.doc.Meta.Source = "kicad"
	db := fontdb.New(opts.FontFS, opts.FontDirs, !opts.NoSystemFonts)
	set := fontset.New(db, func(msg string) { c.warnf("%s", msg) })
	c.fonts = &cad.Fonts{Set: set}
	c.plotter = &cad.Plotter{Fonts: c.fonts, Thin: 0.1 * 72 / 25.4}
	c.cvs = canvas.NewBuilder(c.doc, set)
	c.readProject(in)

	views := strings.ToLower(opts.Views)
	switch views {
	case "", "all", "schematic", "board":
	default:
		return nil, fmt.Errorf("kicad: views %q is not all, schematic or board", opts.Views)
	}
	res := &Result{Doc: c.doc}
	var pages []pageOut
	if in.schData != nil && views != "board" {
		ps, err := c.schematic(in)
		if err != nil {
			if in.pcbData == nil {
				return nil, err
			}
			c.warnf("%v", err)
		}
		pages = append(pages, ps...)
		res.Sheets = len(ps)
	}
	if in.pcbData != nil && views != "schematic" {
		ps, err := c.board(in)
		if err != nil {
			if len(pages) == 0 {
				return nil, err
			}
			c.warnf("%v", err)
		} else {
			res.Boards = 1
		}
		pages = append(pages, ps...)
	}
	if len(pages) == 0 {
		return nil, fmt.Errorf("kicad: nothing to convert")
	}
	if opts.Title != "" {
		c.doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	if !opts.SystemFonts {
		res.EmbeddedFonts = set.Embed(c.doc, fontset.EmbedOptions{NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType})
	}
	c.cvs.Encode()
	set.ReportMissing()
	indexed := map[*bdf.View]bool{}
	for _, p := range pages {
		p.page.Layers = []bdf.Layer{{Role: bdf.RoleBody, Obj: p.cv.Hash()}}
		if p.background != nil {
			p.page.Layers = append([]bdf.Layer{{Role: bdf.RoleBackground, Obj: p.background.Hash()}}, p.page.Layers...)
		}
	}
	for _, p := range pages {
		if !opts.NoTextIndex && !indexed[p.view] {
			indexed[p.view] = true
			if _, err := c.doc.BuildTextIndex(p.view); err != nil {
				c.warnf("text index: %v", err)
			}
		}
	}
	res.Warnings = c.warnings
	return res, nil
}

// pageOut is a page made, whose objects are encoded at the end.
type pageOut struct {
	view       *bdf.View
	page       *bdf.Page
	cv         *canvas.Canvas
	background *canvas.Canvas // the drawing sheet, or nil
}

// readInput finds the schematic, board and project of an input.
func (c *conv) readInput(data []byte) (*input, error) {
	head := bytes.TrimLeft(bytes.TrimPrefix(data[:min(len(data), 512)], []byte("\xef\xbb\xbf")), " \t\r\n")
	name := path.Base(strings.ReplaceAll(c.opts.FileName, `\`, "/"))
	if name == "." || name == "/" {
		name = ""
	}
	switch {
	case bytes.HasPrefix(head, []byte("(kicad_sch")):
		in := &input{sch: nameOr(name, "schematic.kicad_sch"), schData: data}
		c.siblings(in, in.sch)
		return in, nil
	case bytes.HasPrefix(head, []byte("(kicad_pcb")):
		in := &input{pcb: nameOr(name, "board.kicad_pcb"), pcbData: data}
		c.siblings(in, in.pcb)
		return in, nil
	case bytes.HasPrefix(head, []byte("PK\x03\x04")):
		return c.readZip(data)
	case bytes.HasPrefix(head, []byte("{")):
		in := &input{pro: nameOr(name, "project.kicad_pro"), proData: data}
		base := strings.TrimSuffix(in.pro, path.Ext(in.pro))
		if c.refs != nil {
			if b, p, err := c.readRef(in.pro, path.Base(base)+".kicad_sch"); err == nil {
				in.sch, in.schData = p, b
			}
			if b, p, err := c.readRef(in.pro, path.Base(base)+".kicad_pcb"); err == nil {
				in.pcb, in.pcbData = p, b
			}
		}
		if in.schData == nil && in.pcbData == nil {
			return nil, fmt.Errorf("kicad: the project's schematic and board (%s.kicad_sch, %s.kicad_pcb) are not among the files given", base, base)
		}
		return in, nil
	}
	return nil, fmt.Errorf("kicad: not a KiCad 6 or later schematic, board, project or zip archive")
}

func nameOr(name, def string) string {
	if name == "" {
		return def
	}
	return name
}

// siblings reads the project file of a schematic or board, when the files
// given with it hold one.
func (c *conv) siblings(in *input, main string) {
	if c.refs == nil {
		return
	}
	base := strings.TrimSuffix(path.Base(main), path.Ext(main))
	if b, p, err := c.readRef(main, base+".kicad_pro"); err == nil {
		in.pro, in.proData = p, b
	}
}

// readZip finds the project in a zip archive: the schematic and board that
// the project file names, or the root schematic (the one no other sheet
// refers to) and the board.
func (c *conv) readZip(data []byte) (*input, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("kicad: %w", err)
	}
	var schs, pcbs, pros []string
	for _, f := range zr.File {
		n := f.Name
		if f.FileInfo().IsDir() || strings.HasPrefix(n, "__MACOSX/") || strings.Contains(n, "/.") || strings.HasPrefix(n, ".") {
			continue
		}
		switch strings.ToLower(path.Ext(n)) {
		case ".kicad_sch":
			schs = append(schs, n)
		case ".kicad_pcb":
			pcbs = append(pcbs, n)
		case ".kicad_pro":
			pros = append(pros, n)
		}
	}
	zopts := converter.Options{Files: zr}
	c.refs = zopts.ReadRef
	read := func(name string) []byte {
		b, _, err := c.readRef("", name)
		if err != nil {
			c.warnf("%v", err)
			return nil
		}
		return b
	}
	in := &input{}
	if len(pros) > 0 {
		// the project nearest the top of the archive
		best := pros[0]
		for _, p := range pros[1:] {
			if strings.Count(p, "/") < strings.Count(best, "/") {
				best = p
			}
		}
		in.pro, in.proData = best, read(best)
		base := strings.TrimSuffix(best, path.Ext(best))
		for _, s := range schs {
			if strings.EqualFold(s, base+".kicad_sch") {
				in.sch = s
			}
		}
		for _, p := range pcbs {
			if strings.EqualFold(p, base+".kicad_pcb") {
				in.pcb = p
			}
		}
	}
	if in.sch == "" && len(schs) > 0 {
		in.sch = rootSchematic(schs, read)
	}
	if in.pcb == "" && len(pcbs) > 0 {
		in.pcb = pcbs[0]
		for _, p := range pcbs[1:] {
			if strings.Count(p, "/") < strings.Count(in.pcb, "/") {
				in.pcb = p
			}
		}
	}
	if in.sch != "" {
		in.schData = read(in.sch)
	}
	if in.pcb != "" {
		in.pcbData = read(in.pcb)
	}
	if in.schData == nil && in.pcbData == nil {
		return nil, fmt.Errorf("kicad: the zip archive holds no KiCad schematic or board")
	}
	return in, nil
}

// rootSchematic returns the schematic that no other schematic uses as a
// sheet (the first such, nearest the top of the archive).
func rootSchematic(schs []string, read func(string) []byte) string {
	used := map[string]bool{}
	for _, s := range schs {
		b := read(s)
		n, err := parse(b)
		if err != nil {
			continue
		}
		for _, sh := range n.children("sheet") {
			if f := sheetFile(sh); f != "" {
				used[path.Clean(path.Join(path.Dir(s), strings.ReplaceAll(f, `\`, "/")))] = true
			}
		}
	}
	best := ""
	for _, s := range schs {
		if used[s] {
			continue
		}
		if best == "" || strings.Count(s, "/") < strings.Count(best, "/") {
			best = s
		}
	}
	if best == "" {
		best = schs[0]
	}
	return best
}
