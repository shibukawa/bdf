// Package converter is the registry of the input formats that convert into
// BDF documents, with what their converters have in common: options,
// format detection and page selection.
//
// The converters themselves are its subpackages (converter/pdf,
// converter/pptx, converter/xlsx, converter/csv, converter/visio,
// converter/emf). Each registers its format when it is imported, so a
// program supports the formats whose packages it links in:
//
//	import _ "github.com/shibukawa/bdf/converter/pdf"  // PDF only
//	import _ "github.com/shibukawa/bdf/converter/all"  // every format
//
// What the Office converters share is in converter/internal: ooxml (OPC
// packages and their XML) with ooxml/drawingml (shapes, text, tables,
// charts), fontset (fonts for text layout and their embedding), canvas
// (objects under construction) and metafile (EMF/WMF pictures).
//
// Password-protected inputs open with Options.Password. Encrypted Office
// documents are decrypted here (converter/internal/offcrypto), before their
// format is detected; PDF passwords are handled by converter/pdf.
package converter

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/offcrypto"
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
	// CheckPassword, for formats with their own encryption, reports
	// whether an input needs a password to open and checks password
	// against it (see the package's CheckPassword).
	CheckPassword func(r io.ReaderAt, size int64, password string) (protected bool, err error)
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
	// Password opens an encrypted input: the open password of an Office
	// document, the user (or owner) password of a PDF. Inputs that open
	// without one ignore it.
	Password string

	// FileName is the input's file name, when it has one (ConvertFile sets
	// it). Its extension tells the format of inputs whose content does not
	// (a CSV file of one line or one column), and formats that name what
	// they convert after the file use it: a CSV file's sheet.
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
	// Protected reports that the input opened only with Options.Password.
	// Whoever needed the password to read the input should need it to read
	// the document too: encrypt it with the same password (bdf.Lock).
	Protected bool
}

// Errors of Convert and CheckPassword; test for them with errors.Is.
var (
	ErrUnknownFormat    = errors.New("unknown input format")
	ErrPasswordRequired = errors.New("the input is encrypted and needs a password")
	ErrWrongPassword    = errors.New("the password does not open the input")
)

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
	o := Options{}
	if opts != nil {
		o = *opts
	}
	if o.FileName == "" {
		o.FileName = filepath.Base(path)
	}
	res, err := Convert(f, st.Size(), name, &o)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return res, nil
}

// Convert converts an input in the named format (detected when name is
// ""). An encrypted Office document is decrypted with opts.Password first.
func Convert(r io.ReaderAt, size int64, name string, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	var warnings []string
	protected := false
	if offcrypto.IsEncrypted(r, size) {
		if opts.Password == "" {
			return nil, ErrPasswordRequired
		}
		b, err := offcrypto.Decrypt(r, size, opts.Password)
		switch {
		case errors.Is(err, offcrypto.ErrWrongPassword):
			return nil, ErrWrongPassword
		case errors.Is(err, offcrypto.ErrIntegrity):
			// The package decrypted and is still a ZIP file with its own
			// checksums; say so, and convert it.
			warnings = append(warnings, err.Error())
		case err != nil:
			return nil, err
		}
		r, size, protected = bytes.NewReader(b), int64(len(b)), true
	}
	var format *Format
	if name == "" {
		if format = Detect(r, size); format == nil {
			if format = byExtension(opts.FileName); format == nil {
				return nil, ErrUnknownFormat
			}
		}
	} else if format = Lookup(name); format == nil {
		return nil, fmt.Errorf("unknown format %q", name)
	}
	if opts.Warn != nil {
		for _, w := range warnings {
			opts.Warn(w)
		}
		warnings = nil
	}
	res, err := format.Convert(r, size, opts)
	if err != nil {
		return nil, err
	}
	res.Warnings = append(warnings, res.Warnings...)
	res.Protected = res.Protected || protected
	return res, nil
}

// byExtension returns the format whose usual extension a file name has,
// or nil.
func byExtension(fileName string) *Format {
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == "" {
		return nil
	}
	for _, f := range Formats() {
		if slices.Contains(f.Extensions, ext) {
			return f
		}
	}
	return nil
}

// CheckPassword reports whether an input needs a password to open, and
// checks password against it without converting it: nil when the input
// opens with password (or needs none), ErrPasswordRequired when it needs
// one and password is empty, ErrWrongPassword when password does not open
// it. A server can call it when the file is uploaded and answer at once.
func CheckPassword(r io.ReaderAt, size int64, password string) (protected bool, err error) {
	if offcrypto.IsEncrypted(r, size) {
		if password == "" {
			return true, ErrPasswordRequired
		}
		err := offcrypto.Verify(r, size, password)
		if errors.Is(err, offcrypto.ErrWrongPassword) {
			return true, ErrWrongPassword
		}
		return true, err
	}
	if f := Detect(r, size); f != nil && f.CheckPassword != nil {
		return f.CheckPassword(r, size, password)
	}
	return false, nil
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
