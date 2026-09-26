// Package csv converts CSV and TSV files into BDF documents: one sheet
// view that shows the values the way Excel shows a CSV file it opens (the
// sheet is laid out and drawn by converter/xlsx; see xlsx.Grid).
//
// What a file does not say is guessed, and options override each guess:
// its character encoding (a byte order mark, UTF-8, UTF-16, Shift_JIS,
// EUC-JP, ISO-2022-JP or Windows-1252), its dialect (comma, tab,
// semicolon or vertical bar between fields; double, single or no quotes
// around them; quotes doubled or escaped with a backslash), which values
// are numbers or dates (aligned right, as they are written), and whether
// the first record is a header row (then bold, frozen and marked as the
// column headers). See docs/design.md §3.7.
package csv

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/xlsx"
)

// Header tells whether the first record is a header row.
type Header int

const (
	// HeaderAuto guesses it from the values.
	HeaderAuto Header = iota
	// HeaderYes takes the first record as the header row.
	HeaderYes
	// HeaderNo takes every record as data.
	HeaderNo
)

// NoQuote is the Quote option for files whose fields are never quoted.
const NoQuote = noQuote

// Options controls the conversion.
type Options struct {
	// Name is the sheet name; ConvertFile names the sheet after the file,
	// as Excel does, and the default is "Sheet1".
	Name string
	// Charset is the character encoding as a WHATWG label ("utf-8",
	// "shift_jis", "euc-kr", "gb18030", "windows-1252" …); "" detects it.
	// A byte order mark overrides it.
	Charset string
	// Delimiter separates the fields (an ASCII character); 0 detects it.
	Delimiter rune
	// Quote encloses fields ('"' or '\''); 0 detects it, NoQuote reads
	// quotes as text.
	Quote rune
	// Header tells whether the first record is a header row.
	Header Header
	// TableStyle formats the values as an Excel table with a built-in
	// table style ("TableStyleMedium2" …); "" leaves them plain.
	TableStyle string
	// Title sets the document title.
	Title string
	// FontDirs are searched for fonts before the system font directories.
	FontDirs []string
	// NoSystemFonts restricts font lookup to FontDirs.
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

// Result is the outcome of a conversion, with what was guessed.
type Result struct {
	Doc           *bdf.Document
	Warnings      []string
	EmbeddedFonts int
	// Charset is the encoding the file was read in ("UTF-8", "Shift_JIS",
	// "UTF-16LE" …), and BOM whether a byte order mark told it.
	Charset string
	BOM     bool
	Dialect Dialect
	// Header reports whether the first record is the header row.
	Header bool
	// Rows and Cols are the number of records and the most fields in one.
	Rows, Cols int
}

// ConvertFile converts a CSV or TSV file; the sheet is named after it.
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
	o := Options{}
	if opts != nil {
		o = *opts
	}
	if o.Name == "" {
		o.Name = sheetName(path)
	}
	return Convert(f, st.Size(), &o)
}

// sheetName is the name Excel gives the sheet of a CSV file: the file
// name without its extension.
func sheetName(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

// Convert converts a CSV or TSV file read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if err := checkRune("delimiter", opts.Delimiter); err != nil {
		return nil, err
	}
	if opts.Quote != NoQuote {
		if err := checkRune("quote", opts.Quote); err != nil {
			return nil, err
		}
	}
	if opts.Quote > 0 && opts.Quote == opts.Delimiter {
		return nil, fmt.Errorf("csv: the delimiter and the quote are both %q", opts.Delimiter)
	}
	data := make([]byte, size)
	if n, err := r.ReadAt(data, 0); n < len(data) {
		return nil, fmt.Errorf("csv: %w", err)
	}
	res := &Result{}
	cs, bom := bomCharset(data)
	switch {
	case bom:
	case opts.Charset != "" && opts.Charset != "auto":
		var err error
		if cs, err = namedCharset(opts.Charset); err != nil {
			return nil, fmt.Errorf("csv: %w", err)
		}
	default:
		cs = detectCharset(data)
	}
	res.Charset, res.BOM = cs.name, bom
	text := decode(data, cs)
	data = nil

	sample, cut := sampleOf(text)
	res.Dialect = sniff(sample, cut, opts.Delimiter, opts.Quote).d

	var warnings []string
	warn := func(msg string) {
		if opts.Warn != nil {
			opts.Warn(msg)
		} else {
			warnings = append(warnings, msg)
		}
	}
	rd := newReader(text, res.Dialect)
	var rows [][]xlsx.GridCell
	var head [][]string
	var fields []string
	var first, rest []digitStat // of the first record and of the others
	for {
		f, ok := rd.record(fields)
		if !ok {
			break
		}
		fields = f
		if rd.unfinished {
			warn(fmt.Sprintf("record %d: a quote is never closed; the rest of the file is one value", len(rows)+1))
		}
		if len(head) < maxHeaderSample {
			head = append(head, append([]string(nil), f...))
		}
		cells := make([]xlsx.GridCell, len(f))
		stats := &rest
		if len(rows) == 0 {
			stats = &first
		}
		for len(*stats) < len(f) {
			*stats = append(*stats, digitStat{})
		}
		for i, v := range f {
			cells[i] = xlsx.GridCell{Text: v, Number: v != "" && classify(v).numeric()}
			(*stats)[i].add(v)
		}
		rows = append(rows, cells)
		res.Cols = max(res.Cols, len(f))
	}
	res.Rows = len(rows)
	switch opts.Header {
	case HeaderYes:
		res.Header = len(rows) > 0
	case HeaderAuto:
		res.Header = hasHeader(head)
	}
	codeColumns(rows, first, rest, res.Header)
	g := &xlsx.Grid{Name: opts.Name, Rows: rows, Lang: language(text, cs), TableStyle: opts.TableStyle}
	if res.Header {
		g.HeaderRows = 1
	}
	xr, err := xlsx.ConvertGrid(g, &xlsx.Options{Title: opts.Title, FontDirs: opts.FontDirs, NoSystemFonts: opts.NoSystemFonts,
		SystemFonts: opts.SystemFonts, NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType,
		NoTextIndex: opts.NoTextIndex, Warn: opts.Warn})
	if err != nil {
		return nil, fmt.Errorf("csv: %s", strings.TrimPrefix(err.Error(), "xlsx: "))
	}
	xr.Doc.Meta.Source = "csv"
	res.Doc, res.EmbeddedFonts = xr.Doc, xr.EmbeddedFonts
	res.Warnings = append(warnings, xr.Warnings...)
	return res, nil
}

func checkRune(what string, r rune) error {
	if r == 0 {
		return nil
	}
	if r < 0 || r >= 0x80 || r == '\n' || r == '\r' {
		return fmt.Errorf("csv: the %s must be an ASCII character other than a line break, not %q", what, r)
	}
	return nil
}

// language guesses the language of a file's text for its fonts and its
// Han characters: that of a Japanese, Korean or Chinese encoding, or of
// the kana or hangul in it.
func language(text string, cs charset) string {
	switch cs.name {
	case "Shift_JIS", "EUC-JP", "ISO-2022-JP":
		return "ja"
	case "EUC-KR":
		return "ko"
	case "GBK", "GB18030":
		return "zh-CN"
	case "BIG5":
		return "zh-TW"
	}
	kana, hangul := 0, 0
	for i, r := range text {
		if i > 1<<20 {
			break
		}
		switch {
		case unicode.In(r, unicode.Hiragana, unicode.Katakana):
			kana++
		case unicode.Is(unicode.Hangul, r):
			hangul++
		}
	}
	switch {
	case kana > 0 && kana >= hangul:
		return "ja"
	case hangul > 0:
		return "ko"
	}
	return ""
}

// digitStat tells whether a column holds codes: a value of digits with
// leading zeros, and no value other than digits.
type digitStat struct{ code, other bool }

func (d *digitStat) add(v string) {
	v = strings.TrimSpace(v)
	switch {
	case v == "":
	case !allDigits(v):
		d.other = true
	case len(v) > 1 && v[0] == '0':
		d.code = true
	}
}

// codeColumns makes the integers of columns of codes (0012, 0107, 1003)
// text, as the codes with leading zeros are, so that the column is aligned
// alike.
func codeColumns(rows [][]xlsx.GridCell, first, rest []digitStat, header bool) {
	for j := range rest {
		st := rest[j]
		if !header && j < len(first) {
			st.code, st.other = st.code || first[j].code, st.other || first[j].other
		}
		if !st.code || st.other {
			continue
		}
		for i, cells := range rows {
			if j < len(cells) && (i > 0 || !header) {
				cells[j].Number = false
			}
		}
	}
}
