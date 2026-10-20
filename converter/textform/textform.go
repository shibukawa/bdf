// Package textform fills templates of plain text with data: forms whose
// characters are on a grid of columns and lines, as dot matrix printers
// print continuous forms and slips in their text mode, and as receipt
// printers print receipts.
//
// A template is a text file of UTF-8: optional front matter (YAML between
// lines "---") that gives the grid, and the body, which is the print image
// of a page. A line of the body is a line of the page, and a character
// takes a column (a full-width one two):
//
//	---
//	cpi: 10            # columns to the inch (10, 12, 15, 17.1, 20)
//	lpi: 6             # lines to the inch (6, 8)
//	columns: 136       # the width of the paper in columns (80)
//	lines: 66          # the length of a page in lines (66; 0: a roll, a receipt)
//	margin: {left: 2, top: 1}
//	skip: 3            # lines skipped over the perforation at the bottom
//	title: 4           # the first lines repeat at the top of each page
//	font: BIZ UDGothic # the monospace font of the pages (a family, or a font file)
//	ambiguous: narrow  # the columns of box drawing, ○, ※ …: narrow (1) or wide (2)
//	---
//	{{#_big}}納品書{{/_big}}                      No. {{no | right:8}}
//	{{customer | width:30}} 御中
//	{{#preprinted}}
//	┌────────────────────┬────┐
//	{{/preprinted}}
//	{{#items}}
//	│{{name | width:20 | cut}}│{{qty | right:4}}│
//	│                    │    │
//	{{/items}}
//
// The tags are those of the other templates (package form/tag): values,
// sections, number formats. A value takes the columns of its text, or
// those of its tag: {{name | width:20}} on the left of 20 columns,
// right:n and center:n; a longer value is written whole with a warning,
// cut (| cut), or goes on to the lines under it in the same columns
// (| wrap). A line with only the tags of sections of lines takes no line.
//
// Sections repeat as in a workbook: one whose tags are in the same line
// repeats characters; one over lines fills the lines between its tags, an
// item to a line (slots; a line of them without tags is the first again);
// {{#items | rows}} repeats the lines for each item, and the pages break
// when they are full, with the title lines on each; {{#pages | pages}}
// repeats them on a page for each item. A form feed in the body breaks the
// page. A form without a section of rows must fit its pages.
//
// Three sections are not of the data: {{#preprinted}} marks the rules and
// names of paper that has them printed already (they are drawn on a layer
// of the role "background", which printing on that paper leaves out), and
// {{#_wide}}, {{#_tall}} and {{#_big}} characters of double width, height
// or both. A template has no pictures and no codes.
//
// The form is a fixed view of pages of the size of the paper, its
// characters drawn in their columns with a monospace font (the advance of
// each run is that of its columns, so they keep to the grid whatever the
// font), and the print image itself, as the part of the type
// bdf.PartText (package form/textpart): what the text modes of printers
// and plain text are written from.
package textform

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
	"time"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/form/tag"
	"github.com/shibukawa/bdf/form/textpart"
	"github.com/shibukawa/bdf/internal/fontdb"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
)

// TemplateData is what a template is filled with.
type TemplateData struct {
	// Data is the data: an object of the values of JSON.
	Data map[string]any
	// Strict makes a tag whose name the data does not have an error;
	// without it the tag is replaced by nothing, with a warning.
	Strict bool
	// Files holds the font file that the front matter names by its path
	// (font: fonts/slip.ttf), when it names one.
	Files fs.FS
	// Format writes a value in the number format of a tag (nil: values
	// are written as the General format of a spreadsheet writes them).
	Format tag.Formatter
}

// Options controls the fonts of the pages.
type Options struct {
	// Pages are the pages of the form to keep (nil: all).
	Pages         conv.Pages
	Title         string
	FontFS        fs.FS
	FontDirs      []string
	NoSystemFonts bool
	SystemFonts   bool
	NoSubset      bool
	NoWOFF2       bool
	IgnoreFSType  bool
	NoTextIndex   bool
}

// OptionsOf returns the options that those of a conversion give.
func OptionsOf(o *conv.Options) *Options {
	if o == nil {
		return &Options{}
	}
	return &Options{Pages: o.Pages, Title: o.Title, FontFS: o.FontFS, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts,
		NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex}
}

// Result is a filled template.
type Result struct {
	Doc      *bdf.Document
	Warnings []string
	Fields   []tag.Field
	// Image is the print image of the form: the part of the document of
	// the type bdf.PartText.
	Image *textpart.Image
}

// Inspect reads the tags of a template: the names they look up, or the
// first mistake in them.
func Inspect(r io.ReaderAt, size int64) ([]tag.Field, error) {
	data, err := readAll(r, size)
	if err != nil {
		return nil, err
	}
	t, err := parse(data)
	if err != nil {
		return nil, err
	}
	return t.fields(), nil
}

// Fill fills a template with data.
func Fill(r io.ReaderAt, size int64, td *TemplateData, opts *Options) (*Result, error) {
	return fill(r, size, td, opts, false)
}

// Show lays a template out as it is, with its tags: its lines on the pages
// of its paper (as a compiled template shows a template without sample
// data).
func Show(r io.ReaderAt, size int64, td *TemplateData, opts *Options) (*Result, error) {
	return fill(r, size, td, opts, true)
}

func fill(r io.ReaderAt, size int64, td *TemplateData, opts *Options, raw bool) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if td == nil {
		td = &TemplateData{}
	}
	data, err := readAll(r, size)
	if err != nil {
		return nil, err
	}
	t, err := parse(data)
	if err != nil {
		return nil, err
	}
	f := &filler{t: t, td: td, format: td.Format, warned: map[string]bool{}, raw: raw}
	f.warnings = append(f.warnings, t.warnings...)
	if raw {
		for _, l := range t.lines {
			if err := f.emit(l, 0, nil); err != nil {
				return nil, err
			}
		}
	} else {
		root, _ := tag.Normalize(td.Data).(map[string]any)
		if root == nil {
			root = map[string]any{}
		}
		if err := f.lines(t.root, 0, len(t.lines)-1, tag.NewScope(root)); err != nil {
			return nil, err
		}
	}
	pages, err := f.paginate()
	if err != nil {
		return nil, err
	}
	if sel := opts.Pages.Numbers(len(pages)); sel != nil {
		var kept []page
		for _, n := range sel {
			if n >= 1 && n <= len(pages) {
				kept = append(kept, pages[n-1])
			}
		}
		pages = kept
	}
	if len(pages) > textpart.MaxPages {
		return nil, fmt.Errorf("text: the form has more than %d pages", textpart.MaxPages)
	}
	im := f.image(pages)
	tables := make([][]int, len(pages))
	for i, p := range pages {
		for k, ol := range p.lines {
			tables[i] = append(tables[i], p.tables[k])
			if ol.tall() {
				tables[i] = append(tables[i], p.tables[k])
			}
		}
	}

	doc := bdf.NewDocument()
	doc.Meta.Source = "text"
	if opts.Title != "" {
		doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	view := doc.NewView("pages", bdf.ViewFixed, opts.Title)
	warn := func(msg string) { f.warnOnce("font:"+msg, "%s", msg) }
	fontFS := opts.FontFS
	named := t.grid.font
	if isFontPath(named) {
		file, err := readFont(td.Files, named)
		if err != nil {
			return nil, err
		}
		fontFS = &oneFontFS{name: path.Base(named), data: file, rest: opts.FontFS}
	}
	db := fontdb.New(fontFS, opts.FontDirs, !opts.NoSystemFonts)
	set := fontset.New(db, warn)
	fams := families(db, "")
	if named != "" {
		if isFontPath(named) {
			// (the family of the font file: of its first face)
			for _, fc := range db.Faces {
				if fc.Path == path.Base(named) && fc.Family != "" {
					fams = families(db, fc.Family)
					break
				}
			}
		} else {
			fams = families(db, named)
		}
	}
	im.Font = fams[0]
	draw(doc, view, im, tables, fams, set, func() {
		if !opts.SystemFonts {
			set.Embed(doc, fontset.EmbedOptions{NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType})
		}
	})
	set.ReportMissing()
	if _, err := im.AddTo(doc); err != nil {
		return nil, err
	}
	if !opts.NoTextIndex {
		if _, err := doc.BuildTextIndex(view); err != nil {
			f.warnf("text index: %v", err)
		}
	}
	return &Result{Doc: doc, Warnings: f.warnings, Fields: t.fields(), Image: im}, nil
}

// isFontPath reports whether the font of the front matter is a font file
// rather than a family.
func isFontPath(s string) bool {
	switch strings.ToLower(path.Ext(s)) {
	case ".ttf", ".otf", ".ttc", ".otc":
		return true
	}
	return false
}

// maxFontFile bounds the font file that a template names.
const maxFontFile = 64 << 20

// readFont reads the font file that the front matter names, from the files
// of the template.
func readFont(files fs.FS, name string) ([]byte, error) {
	if files == nil {
		return nil, fmt.Errorf("text: the front matter names the font file %q, and the template has no files beside it to read it from", name)
	}
	p := path.Clean(strings.ReplaceAll(name, "\\", "/"))
	if !fs.ValidPath(p) {
		return nil, fmt.Errorf("text: the font file %q is not a path in the directory of the template", name)
	}
	fl, err := files.Open(p)
	if err != nil {
		return nil, fmt.Errorf("text: the font file of the front matter: %w", err)
	}
	defer fl.Close()
	data, err := io.ReadAll(io.LimitReader(fl, maxFontFile+1))
	if err != nil {
		return nil, fmt.Errorf("text: the font file of the front matter: %w", err)
	}
	if len(data) > maxFontFile {
		return nil, fmt.Errorf("text: the font file %q is larger than %d MiB", name, maxFontFile>>20)
	}
	return data, nil
}

// oneFontFS is a file system of a font file, at its top, and the files of
// another file system (when there is one).
type oneFontFS struct {
	name string
	data []byte
	rest fs.FS
}

func (o *oneFontFS) Open(name string) (fs.File, error) {
	switch {
	case name == o.name:
		return &memFile{Reader: bytes.NewReader(o.data), info: fileInfo{name: o.name, size: int64(len(o.data))}}, nil
	case name == ".":
		return &memFile{Reader: bytes.NewReader(nil), info: fileInfo{name: ".", dir: true}}, nil
	case o.rest != nil:
		return o.rest.Open(name)
	}
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}

func (o *oneFontFS) ReadDir(name string) ([]fs.DirEntry, error) {
	var out []fs.DirEntry
	if o.rest != nil {
		es, err := fs.ReadDir(o.rest, name)
		if err != nil && !(name == "." && errors.Is(err, fs.ErrNotExist)) {
			return nil, err
		}
		for _, e := range es {
			if !(name == "." && e.Name() == o.name) {
				out = append(out, e)
			}
		}
	} else if name != "." {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrNotExist}
	}
	if name == "." {
		out = append(out, fs.FileInfoToDirEntry(fileInfo{name: o.name, size: int64(len(o.data))}))
	}
	return out, nil
}

type memFile struct {
	*bytes.Reader
	info fileInfo
}

func (m *memFile) Stat() (fs.FileInfo, error) { return m.info, nil }
func (m *memFile) Close() error               { return nil }

type fileInfo struct {
	name string
	size int64
	dir  bool
}

func (i fileInfo) Name() string { return i.name }
func (i fileInfo) Size() int64  { return i.size }
func (i fileInfo) Mode() fs.FileMode {
	if i.dir {
		return fs.ModeDir | 0o555
	}
	return 0o444
}
func (i fileInfo) ModTime() time.Time { return time.Time{} }
func (i fileInfo) IsDir() bool        { return i.dir }
func (i fileInfo) Sys() any           { return nil }

// Params are the options of writing a form as plain text (Plain).
var Params = []conv.Param{
	{Name: "encoding", Usage: "the encoding of the text: utf-8 (the default), shift_jis, euc-jp, iso-2022-jp, or another name of the WHATWG Encoding Standard"},
	{Name: "newline", Usage: "the end of a line: lf (the default) or crlf"},
	{Name: "paper", Usage: "with preprinted, the ranges of the preprinted paper are written as spaces"},
}

// Plain writes the print image of a form as plain text: its lines, the
// pages separated by form feeds, with the margins of the paper, in an
// encoding (the parameters of Params). A character that the encoding does
// not have is written as 〓 (or ? when it has none either), with a warning.
func Plain(im *textpart.Image, params map[string]string) ([]byte, []string, error) {
	o := textpart.TextOptions{Margins: true}
	switch strings.ToLower(params["newline"]) {
	case "", "lf":
	case "crlf":
		o.CRLF = true
	default:
		return nil, nil, fmt.Errorf("text: newline is %q; it is lf or crlf", params["newline"])
	}
	switch strings.ToLower(params["paper"]) {
	case "", "plain":
	case "preprinted":
		o.NoPaper = true
	default:
		return nil, nil, fmt.Errorf("text: paper is %q; it is plain or preprinted", params["paper"])
	}
	s := im.Text(o)
	name := strings.ToLower(params["encoding"])
	if name == "" || name == "utf-8" || name == "utf8" {
		return []byte(s), nil, nil
	}
	enc, err := htmlindex.Get(name)
	if err != nil {
		return nil, nil, fmt.Errorf("text: the encoding %q is not known", params["encoding"])
	}
	return encode(s, enc)
}

// encode writes text in an encoding, with the characters it does not
// have replaced.
func encode(s string, enc encoding.Encoding) ([]byte, []string, error) {
	if b, err := enc.NewEncoder().Bytes([]byte(s)); err == nil {
		return b, nil, nil
	}
	geta := "〓"
	if _, err := enc.NewEncoder().Bytes([]byte(geta)); err != nil {
		geta = "?"
	}
	var warnings []string
	missing := map[rune]bool{}
	var b strings.Builder
	for _, r := range s {
		if _, ok := missing[r]; !ok {
			_, err := enc.NewEncoder().Bytes([]byte(string(r)))
			missing[r] = err != nil
			if err != nil && len(warnings) < 32 {
				warnings = append(warnings, fmt.Sprintf("the character %c (U+%04X) is not in the encoding; it is written as %s", r, r, geta))
			}
		}
		if missing[r] {
			b.WriteString(geta)
			continue
		}
		b.WriteRune(r)
	}
	out, err := enc.NewEncoder().Bytes([]byte(b.String()))
	if err != nil {
		return nil, nil, fmt.Errorf("text: %w", err)
	}
	return out, warnings, nil
}
