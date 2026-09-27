package epub

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"

	conv "github.com/shibukawa/bdf/converter"
)

// Params are the format-specific options of EPUB publications.
var Params = []conv.Param{
	{Name: "views", Usage: "pages (default): book pages; scroll: one long column, horizontal; both"},
	{Name: "paper", Usage: "page size: a5 (default), a4, a6, b5, b6, letter, or WIDTHxHEIGHT in millimeters"},
	{Name: "size", Usage: "text size in points (default 12)"},
	{Name: "font", Usage: "font family of the text (default: an installed sans-serif family)"},
	{Name: "mono", Usage: "font family of code (default: an installed monospaced family)"},
	{Name: "width", Usage: "text width of the scroll view in points (default: 36 ems of the text)"},
}

func init() {
	conv.Register(&conv.Format{
		Name:        "epub",
		Description: "EPUB publication (reflowable books in reader mode, fixed-layout books of pictures)",
		Extensions:  []string{".epub"},
		Params:      Params,
		Detect:      Detect,
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			opts, err := FromConverter(o)
			if err != nil {
				return nil, err
			}
			res, err := Convert(r, size, opts)
			if err != nil {
				return nil, err
			}
			return Summary(res), nil
		},
	})
}

// FromConverter maps the registry's options (and Params) onto Options.
func FromConverter(o *conv.Options) (*Options, error) {
	opts := &Options{Views: o.Param("views"), Pages: o.Pages, Title: o.Title, Font: o.Param("font"), MonoFont: o.Param("mono"),
		Images: o.Images, FontFS: o.FontFS, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts,
		EmbedFonts: o.EmbedFonts && !o.SystemFonts, NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType,
		NoTextIndex: o.NoTextIndex, Warn: o.Warn}
	for _, p := range []struct {
		name string
		dst  *float64
	}{{"width", &opts.Width}, {"size", &opts.FontSize}} {
		if v := o.Param(p.name); v != "" {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil || f <= 0 {
				return nil, fmt.Errorf("parameter %s: %q is not a positive number", p.name, v)
			}
			*p.dst = f
		}
	}
	if v := o.Param("paper"); v != "" {
		p, err := ParsePaper(v)
		if err != nil {
			return nil, fmt.Errorf("parameter paper: %w", err)
		}
		opts.Paper = p
	}
	return opts, nil
}

// Summary makes a registry result of a result.
func Summary(res *Result) *conv.Result {
	var s string
	switch {
	case res.FixedLayout:
		s = fmt.Sprintf("%d fixed-layout page(s)", res.Pages)
	case res.Pages > 0 && res.Strips > 0:
		s = fmt.Sprintf("%d page(s), %d scroll strip(s)", res.Pages, res.Strips)
	case res.Pages > 0:
		s = fmt.Sprintf("%d page(s)", res.Pages)
	default:
		s = fmt.Sprintf("%d scroll strip(s)", res.Strips)
	}
	s += fmt.Sprintf(" from %d content document(s), %d image(s)", res.Chapters, res.Images)
	if !res.FixedLayout {
		s += fmt.Sprintf(", %d embedded font(s)", res.EmbeddedFonts)
	}
	if res.Vertical {
		s += ", vertical text"
	}
	return &conv.Result{Doc: res.Doc, Warnings: res.Warnings, Summary: s}
}

const mimetype = "application/epub+zip"

// Detect reports whether an input is an EPUB publication: a ZIP file whose
// first entry is the mimetype file saying application/epub+zip, as the
// specification requires, or (for files made carelessly) one that has
// such a mimetype file elsewhere, or a container.xml naming a package
// document.
func Detect(head []byte, r io.ReaderAt, size int64) bool {
	if !bytes.HasPrefix(head, []byte("PK\x03\x04")) {
		return false
	}
	if len(head) >= 30 {
		n, extra := int(binary.LittleEndian.Uint16(head[26:])), int(binary.LittleEndian.Uint16(head[28:]))
		if 30+n <= len(head) && string(head[30:30+n]) == "mimetype" {
			if data := head[min(30+n+extra, len(head)):]; bytes.HasPrefix(data, []byte(mimetype)) {
				return true
			}
		}
	}
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return false
	}
	for _, f := range zr.File {
		switch f.Name {
		case "mimetype":
			if f.UncompressedSize64 < 64 {
				if rc, err := f.Open(); err == nil {
					b, _ := io.ReadAll(rc)
					rc.Close()
					if strings.TrimSpace(string(b)) == mimetype {
						return true
					}
				}
			}
		case "META-INF/container.xml":
			if f.UncompressedSize64 < 1<<20 {
				if rc, err := f.Open(); err == nil {
					b, _ := io.ReadAll(rc)
					rc.Close()
					if bytes.Contains(b, []byte("application/oebps-package+xml")) {
						return true
					}
				}
			}
		}
	}
	return false
}
