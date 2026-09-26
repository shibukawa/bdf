// Package converter is the registry of the input formats that convert into
// BDF documents, with what their converters have in common: options,
// format detection and page selection.
//
// The converters themselves are its subpackages (converter/pdf,
// converter/pptx, converter/xlsx, converter/csv, converter/emf). Each registers its format
// when it is imported, so a program supports the formats whose packages it
// links in:
//
//	import _ "github.com/shibukawa/bdf/converter/pdf"  // PDF only
//	import _ "github.com/shibukawa/bdf/converter/all"  // every format
//
// What the Office converters share is in converter/internal: ooxml (OPC
// packages and their XML) with ooxml/drawingml (shapes, text, tables,
// charts), fontset (fonts for text layout and their embedding), canvas
// (objects under construction) and metafile (EMF/WMF pictures).
package converter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/imgconv"
)

// Format is an input format. Converter packages register theirs in init.
type Format struct {
	// Name identifies the format ("pdf", "pptx", "emf").
	Name string
	// Description names the format for people ("PowerPoint presentation").
	Description string
	// Extensions are the usual file name extensions, with the dot.
	Extensions []string
	// Params are the format-specific options it reads from Options.Params.
	Params []Param
	// Detect reports whether an input is in the format; head holds its
	// first bytes (up to 1 KiB).
	Detect func(head []byte, r io.ReaderAt, size int64) bool
	// Convert converts an input.
	Convert func(r io.ReaderAt, size int64, opts *Options) (*Result, error)
}

// Param is a format-specific option.
type Param struct {
	Name  string
	Usage string
}

// Options are the conversion settings; each format uses those that apply
// to it and ignores the others.
type Options struct {
	// Title overrides the document title.
	Title string
	// Pages selects 1-based pages (or slides, or sheets); nil converts all
	// of them.
	Pages []int
	// Images controls whether raster images are re-encoded (see imgconv).
	// The zero value keeps images as they are.
	Images imgconv.Options

	// The fonts of formats whose text the converter lays out (Office
	// documents, metafiles):

	// FontDirs are searched for fonts before the system font directories.
	FontDirs []string
	// NoSystemFonts restricts font lookup to FontDirs.
	NoSystemFonts bool
	// SystemFonts refers to fonts by family name instead of embedding them.
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
	// Params holds format-specific options by name (see Format.Params).
	Params map[string]string
	// FileName is the input's file name, when it has one (ConvertFile sets
	// it). Formats that name what they convert after it use it: a CSV
	// file's sheet.
	FileName string
	// Warn receives non-fatal problems; when nil they are collected in
	// Result.Warnings.
	Warn func(msg string)
}

// Param returns a format-specific option ("" when it is not set).
func (o *Options) Param(name string) string { return o.Params[name] }

// BoolParam reads a format-specific option that is a boolean ("" is false).
func (o *Options) BoolParam(name string) (bool, error) {
	v := o.Params[name]
	if v == "" {
		return false, nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("parameter %s: %q is not a boolean", name, v)
	}
	return b, nil
}

// Result is the outcome of a conversion.
type Result struct {
	Doc      *bdf.Document
	Warnings []string
	// Summary describes the result in a line ("5 slide(s), 2 embedded
	// font(s)").
	Summary string
}

var (
	mu      sync.RWMutex
	formats = map[string]*Format{}
)

// Register makes a format available. It panics when the name is taken or
// the format cannot detect or convert.
func Register(f *Format) {
	mu.Lock()
	defer mu.Unlock()
	if f.Name == "" || f.Detect == nil || f.Convert == nil {
		panic("converter: incomplete format " + strconv.Quote(f.Name))
	}
	if _, ok := formats[f.Name]; ok {
		panic("converter: format " + f.Name + " registered twice")
	}
	formats[f.Name] = f
}

// Formats returns the registered formats by name.
func Formats() []*Format {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]*Format, 0, len(formats))
	for _, f := range formats {
		out = append(out, f)
	}
	slices.SortFunc(out, func(a, b *Format) int { return strings.Compare(a.Name, b.Name) })
	return out
}

// Lookup returns the registered format with the name, or nil.
func Lookup(name string) *Format {
	mu.RLock()
	defer mu.RUnlock()
	return formats[name]
}

// Detect returns the registered format of an input, or nil when no format
// recognizes it.
func Detect(r io.ReaderAt, size int64) *Format {
	head := make([]byte, 1024)
	n, _ := r.ReadAt(head, 0)
	head = head[:n]
	for _, f := range Formats() {
		if f.Detect(head, r, size) {
			return f
		}
	}
	return nil
}

// DetectFile is Detect for a file path.
func DetectFile(path string) (*Format, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return Detect(f, st.Size()), nil
}

// ConvertFile converts a file in the named format (detected when name is
// "").
func ConvertFile(path, name string, opts *Options) (*Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	var format *Format
	if name == "" {
		if format = Detect(f, st.Size()); format == nil {
			return nil, fmt.Errorf("%s: unknown input format", path)
		}
	} else if format = Lookup(name); format == nil {
		return nil, fmt.Errorf("unknown format %q", name)
	}
	o := Options{}
	if opts != nil {
		o = *opts
	}
	if o.FileName == "" {
		o.FileName = filepath.Base(path)
	}
	return format.Convert(f, st.Size(), &o)
}

// PageRange parses "1-3,5,8-" style selections of 1-based page (or slide)
// numbers, clamped to 1..count.
func PageRange(spec string, count int) ([]int, error) {
	var out []int
	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		lo, hi := 1, count
		if i := strings.IndexByte(part, '-'); i >= 0 {
			if i > 0 {
				if _, err := fmt.Sscanf(part[:i], "%d", &lo); err != nil {
					return nil, fmt.Errorf("bad page range %q", part)
				}
			}
			if i+1 < len(part) {
				if _, err := fmt.Sscanf(part[i+1:], "%d", &hi); err != nil {
					return nil, fmt.Errorf("bad page range %q", part)
				}
			}
		} else {
			if _, err := fmt.Sscanf(part, "%d", &lo); err != nil {
				return nil, fmt.Errorf("bad page %q", part)
			}
			hi = lo
		}
		lo, hi = max(1, lo), min(count, hi)
		for n := lo; n <= hi; n++ {
			out = append(out, n)
		}
	}
	return out, nil
}
