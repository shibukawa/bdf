// Package emf converts Windows metafiles (.emf and .wmf) into BDF
// documents: one page the size of the picture, drawn by replaying the
// metafile's records (the same replay draws the metafiles pictures of
// Office documents; see docs/design.md §3.4). Text is laid out with the
// available fonts, which are embedded as subsets.
package emf

import (
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/metafile"
	"github.com/shibukawa/bdf/imgconv"
)

// Options controls the conversion.
type Options struct {
	// Title overrides the document title (by default the picture name an
	// EMF header carries).
	Title string
	// Images controls whether the bitmaps of the metafile are re-encoded
	// (see imgconv). The zero value stores them as PNG.
	Images imgconv.Options
	// FontFS holds fonts that are not in the local file system; it is
	// searched before FontDirs (see converter.Options.FontFS).
	FontFS fs.FS
	// FontDirs are searched for fonts before the system font directories.
	FontDirs []string
	// NoSystemFonts restricts font lookup to FontFS and FontDirs.
	NoSystemFonts bool
	// SystemFonts refers to fonts by family name instead of embedding the
	// fonts used for layout.
	SystemFonts bool
	// NoSubset embeds whole fonts instead of the glyphs in use.
	NoSubset bool
	// NoWOFF2 stores embedded fonts as TrueType/OpenType instead of WOFF2.
	NoWOFF2 bool
	// IgnoreFSType embeds fonts whose OS/2 fsType forbids embedding or
	// subsetting. Set it only when you hold the rights to embed the fonts.
	IgnoreFSType bool
	// NoTextIndex skips building the text index part.
	NoTextIndex bool
	// Warn receives non-fatal problems; when nil they are collected in
	// Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc      *bdf.Document
	Warnings []string
	// W and H are the page size in points.
	W, H          float64
	EmbeddedFonts int
}

// maxSize bounds the metafiles read.
const maxSize = 1 << 30

// ConvertFile converts a .emf or .wmf file.
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

// Convert converts a metafile read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if size > maxSize {
		return nil, fmt.Errorf("emf: the metafile is larger than %d bytes", maxSize)
	}
	data := make([]byte, size)
	if _, err := r.ReadAt(data, 0); err != nil && err != io.EOF {
		return nil, fmt.Errorf("emf: %w", err)
	}
	kind := metafile.Kind(data)
	if kind == "" {
		return nil, fmt.Errorf("emf: not a Windows metafile")
	}
	w, h, ok := metafile.Size(data)
	if !ok {
		return nil, fmt.Errorf("emf: the metafile has no picture size")
	}
	res := &Result{Doc: bdf.NewDocument(), W: w, H: h}
	warned := map[string]bool{}
	warn := func(msg string) {
		if warned[msg] {
			return
		}
		warned[msg] = true
		if opts.Warn != nil {
			opts.Warn(msg)
		} else {
			res.Warnings = append(res.Warnings, msg)
		}
	}
	doc := res.Doc
	doc.Meta.Source = kind
	if _, name := metafile.Description(data); name != "" {
		doc.Meta.DC.Title = bdf.DCValues{name}
	}
	if opts.Title != "" {
		doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	db := fontdb.New(opts.FontFS, opts.FontDirs, !opts.NoSystemFonts)
	fonts := fontset.New(db, warn)
	b := canvas.NewBuilder(doc, fonts)
	cv := b.New()
	metafile.Draw(cv, data, 0, 0, w, h, &metafile.Options{Doc: doc, Fonts: fonts, Images: opts.Images, Warn: warn})
	if !cv.Drawn {
		cv.Obj.SetBBox(0, 0, float32(w), float32(h))
	}
	if !opts.SystemFonts {
		res.EmbeddedFonts = fonts.Embed(doc, fontset.EmbedOptions{NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType})
	}
	b.Encode()
	fonts.ReportMissing()
	view := doc.NewView("pages", bdf.ViewFixed, doc.Meta.DC.Title.First())
	page := view.AddPage(float32(w), float32(h))
	page.Layers = append(page.Layers, bdf.Layer{Role: bdf.RoleBody, Obj: cv.Hash()})
	if !opts.NoTextIndex {
		if _, err := doc.BuildTextIndex(view); err != nil {
			warn(fmt.Sprintf("text index: %v", err))
		}
	}
	return res, nil
}
