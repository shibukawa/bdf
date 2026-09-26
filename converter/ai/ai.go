// Package ai converts Adobe Illustrator artwork (.ai) into BDF documents,
// one page per artboard.
//
// An .ai file of Illustrator 9 or later is a PDF with Illustrator's own
// data beside it: saved with "Create PDF Compatible File" (the default), its
// pages are the artboards, drawn by the PDF converter (converter/pdf). A
// page's media box includes the document's bleed; its trim box is the
// artboard, which becomes the page. Illustrator's layers are PDF optional
// content, so hidden layers stay hidden.
//
// Files saved without PDF content hold only a placeholder page, and the
// PostScript-based format of Illustrator 8 and earlier has no PDF at all:
// both are reported as errors rather than converted.
package ai

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/pdf"
	"github.com/shibukawa/bdf/imgconv"
)

// Options controls the conversion.
type Options struct {
	// Pages selects 1-based artboards; nil converts every artboard.
	Pages []int
	// Title overrides the document title.
	Title string
	// Box is the page boundary that becomes the page (see pdf.Options.Box):
	// "trim", the artboard, by default; "bleed" or "media" include the
	// document's bleed.
	Box string
	// NoTextIndex skips building the text index part.
	NoTextIndex bool
	// NoSubset keeps every glyph of embedded fonts.
	NoSubset bool
	// NoWOFF2 stores rebuilt fonts as TrueType/OpenType instead of WOFF2.
	NoWOFF2 bool
	// IgnoreFSType embeds fonts whose OS/2 fsType forbids it. Set it only
	// when you hold the rights to embed the fonts of the document.
	IgnoreFSType bool
	// Images controls whether raster images are re-encoded (see imgconv).
	Images imgconv.Options
	// Password opens an encrypted file.
	Password string
	// Warn receives non-fatal problems; when nil they are collected in
	// Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc       *bdf.Document
	Warnings  []string
	Artboards int
	// Protected reports that the file opened only with Options.Password.
	Protected bool
}

// Errors of Convert for the files it cannot convert; test for them with
// errors.Is.
var (
	// ErrNoPDFContent is a file saved without "Create PDF Compatible File".
	ErrNoPDFContent = errors.New(`the file was saved without PDF content; save it again in Illustrator with "Create PDF Compatible File" on`)
	// ErrPostScript is the PostScript-based format of Illustrator 8 and
	// earlier (which later versions also save as "Illustrator 8").
	ErrPostScript = errors.New("PostScript-based Illustrator files (Illustrator 8 and earlier) are not supported; save it again in Illustrator 9 or later")
)

// Kind tells an Illustrator file from its first bytes: "pdf" for the PDF
// based format (Illustrator 9 and later), "ps" for the PostScript based one,
// "" for other files. The PDF based format is known by its XMP metadata,
// which Illustrator writes near the start of the file.
func Kind(r io.ReaderAt, size int64) string {
	head := make([]byte, min(size, 1024))
	n, _ := r.ReadAt(head, 0)
	head = head[:n]
	switch {
	case bytes.HasPrefix(head, []byte("%!PS-Adobe")):
		if bytes.Contains(head, []byte("%%Creator: Adobe Illustrator")) || bytes.Contains(head, []byte("%%Creator:Adobe Illustrator")) {
			return "ps"
		}
	case bytes.Contains(head, []byte("%PDF-")):
		// Illustrator writes the catalog's metadata first; its thumbnail
		// comes before the property that says the file is a document.
		buf := make([]byte, min(size, 256<<10))
		n, _ := r.ReadAt(buf, 0)
		buf = buf[:n]
		if bytes.Contains(buf, []byte("<illustrator:Type>Document</illustrator:Type>")) || bytes.Contains(buf, []byte(`illustrator:Type="Document"`)) {
			return "pdf"
		}
	}
	return ""
}

// ConvertFile converts an .ai file.
func ConvertFile(path string, opts *Options) (*Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return Convert(f, st.Size(), opts)
}

// Convert converts an Illustrator file read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if Kind(r, size) == "ps" {
		return nil, fmt.Errorf("ai: %w", ErrPostScript)
	}
	box := opts.Box
	if box == "" {
		box = "trim"
	}
	res, err := pdf.Convert(io.NewSectionReader(r, 0, size), &pdf.Options{Pages: opts.Pages, Title: opts.Title, Box: box,
		NoTextIndex: opts.NoTextIndex, NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType,
		Images: opts.Images, Password: opts.Password, Warn: opts.Warn})
	if err != nil {
		return nil, fmt.Errorf("ai: %w", err)
	}
	if placeholder(res.Doc) {
		return nil, fmt.Errorf("ai: %w", ErrNoPDFContent)
	}
	res.Doc.Meta.Source = "ai"
	return &Result{Doc: res.Doc, Warnings: res.Warnings, Artboards: res.Pages, Protected: res.Protected}, nil
}

// placeholder reports whether a document is the one page Illustrator saves
// in place of the artwork when the PDF content is left out. The page says
// so, in the language of the Illustrator that saved it.
func placeholder(doc *bdf.Document) bool {
	if len(doc.Views) != 1 || len(doc.Views[0].Pages) != 1 {
		return false
	}
	resolve := func(h bdf.Hash) *bdf.ObjectPart {
		p := doc.Part(h)
		if p == nil {
			return nil
		}
		o, err := bdf.DecodeObject(p.Data)
		if err != nil {
			return nil
		}
		return o
	}
	var text strings.Builder
	for _, l := range doc.Views[0].Pages[0].Layers {
		o := resolve(l.Obj)
		if o == nil {
			continue
		}
		runs, _ := bdf.ExtractText(o, resolve)
		for _, run := range runs {
			for _, c := range run.Text {
				if !unicode.IsSpace(c) {
					text.WriteRune(unicode.ToLower(c))
				}
			}
		}
	}
	s := text.String()
	return strings.Contains(s, "savedwithoutpdfcontent") || strings.Contains(s, "createpdfcompatiblefile") ||
		strings.Contains(s, "pdf互換ファイルを作成")
}
