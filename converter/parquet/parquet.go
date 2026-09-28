// Package parquet converts Apache Parquet files into BDF documents: one
// sheet view of the table, the column names bold and frozen at the top
// with a gray row of their types under them, and the values the way data
// tools show them (the sheet is laid out and drawn by converter/xlsx; see
// xlsx.Grid).
//
// The reader is written from the format's specification: the Thrift
// footer and page headers, the PLAIN, dictionary, RLE/bit-packed, DELTA_*
// and BYTE_STREAM_SPLIT encodings, data pages of both versions, the
// Snappy, gzip, Zstandard, Brotli and LZ4 codecs, and nested values
// assembled from their repetition and definition levels (lists and maps,
// with the backward-compatibility rules for the structures older writers
// made; structs; Variants, shredded or not). Lists, maps and structs are
// shown as JSON, and geometries (GEOMETRY and GEOGRAPHY, and the WKB
// columns of GeoParquet) as WKT. The file is read a page at a time, and
// only as far as the rows shown. See docs/design.md §3.26.
package parquet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/xlsx"
)

// DefaultRows is the number of rows shown unless Options.Rows says
// otherwise.
const DefaultRows = 10000

// MaxRows is the most rows a sheet shows (as Excel's sheets hold, less
// the header rows).
const MaxRows = 1048576 - 2

// maxColumns is the most columns a sheet shows, as Excel's sheets hold.
const maxColumns = 16384

// maxCells bounds the cells shown when Options.Rows is not set: a table of
// many columns shows fewer rows.
const maxCells = 1 << 22

// maxSheetCells is a hard cap on the cells of the grid, whatever Options.Rows
// asks for. A file can state, in a few bytes, a row group of a huge row count
// (runs of repeated or null values pack into almost nothing), and the reader
// allocates rows×columns cells before reading a page; this cap lets the
// documented large cases through (a million rows of a few columns) but keeps
// a hostile count from asking for a grid of hundreds of gigabytes.
const maxSheetCells = 1 << 24

// Options controls the conversion.
type Options struct {
	// Name is the sheet name; ConvertFile names the sheet after the file,
	// and the default is "Sheet1".
	Name string
	// Rows is the number of rows shown (from the first); 0 shows
	// DefaultRows (fewer for tables of more than 400 columns), and -1 (or
	// more than MaxRows) as many as a sheet holds.
	Rows int
	// NoTypes leaves out the row of the columns' types.
	NoTypes bool
	// TableStyle formats the values as an Excel table with a built-in
	// table style ("TableStyleMedium2" …); "" leaves them plain.
	TableStyle string
	// Title sets the document title.
	Title string
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
	Doc           *bdf.Document
	Warnings      []string
	EmbeddedFonts int
	// Rows is the number of rows in the file, and Shown of those shown.
	Rows, Shown int64
	// Cols is the number of columns (the fields at the top of the
	// schema), and RowGroups the number of row groups.
	Cols, RowGroups int
	// Codecs names the compression codecs of the column chunks read.
	Codecs []string
	// CreatedBy is the writer the file says wrote it.
	CreatedBy string
}

// ErrEncrypted is returned for files whose footer is encrypted.
var ErrEncrypted = errors.New("parquet: the file is encrypted (Parquet modular encryption is not supported)")

// ConvertFile converts a Parquet file; the sheet is named after it.
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

// sheetName names a file's sheet: the file name without its extension.
func sheetName(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

// file is an opened Parquet file.
type file struct {
	r       io.ReaderAt
	meta    *fileMeta
	root    *node
	leaves  []*node
	dataEnd int64 // where the footer starts
}

var magic = []byte("PAR1")

func open(r io.ReaderAt, size int64) (*file, error) {
	if size < 12 {
		return nil, errors.New("parquet: not a Parquet file")
	}
	var head, tail [8]byte
	if _, err := r.ReadAt(head[:4], 0); err != nil {
		return nil, fmt.Errorf("parquet: %w", err)
	}
	if _, err := r.ReadAt(tail[:], size-8); err != nil && err != io.EOF {
		return nil, fmt.Errorf("parquet: %w", err)
	}
	if string(tail[4:]) == "PARE" {
		return nil, ErrEncrypted
	}
	if !bytes.Equal(head[:4], magic) || !bytes.Equal(tail[4:], magic) {
		return nil, errors.New("parquet: not a Parquet file")
	}
	n := int64(binary.LittleEndian.Uint32(tail[:4]))
	if n > size-12 {
		return nil, errors.New("parquet: the footer is out of the file")
	}
	b := make([]byte, n)
	if _, err := r.ReadAt(b, size-8-n); err != nil && err != io.EOF {
		return nil, fmt.Errorf("parquet: %w", err)
	}
	meta, err := readFileMeta(b)
	if err != nil {
		return nil, fmt.Errorf("parquet: file metadata: %w", err)
	}
	root, leaves, err := buildSchema(meta.schema)
	if err != nil {
		return nil, err
	}
	for _, g := range meta.rowGroups {
		if len(g.columns) != len(leaves) {
			return nil, errors.New("parquet: a row group's columns are not the schema's")
		}
	}
	f := &file{r: r, meta: meta, root: root, leaves: leaves, dataEnd: size - 8 - n}
	f.geoColumns()
	return f, nil
}

// geoColumns marks the columns that GeoParquet metadata says hold WKB
// geometries.
func (f *file) geoColumns() {
	s, ok := f.meta.keyValues["geo"]
	if !ok {
		return
	}
	var geo struct {
		Columns map[string]struct {
			Encoding string `json:"encoding"`
		} `json:"columns"`
	}
	if json.Unmarshal([]byte(s), &geo) != nil {
		return
	}
	for _, c := range f.root.children {
		if g, ok := geo.Columns[c.name]; ok && c.leaf && c.phys == typeBinary && strings.EqualFold(g.Encoding, "WKB") {
			c.geo = true
		}
	}
}

// Convert converts a Parquet file read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	var warnings []string
	warn := func(msg string) {
		if opts.Warn != nil {
			opts.Warn(msg)
		} else {
			warnings = append(warnings, msg)
		}
	}
	g, res, err := table(r, size, opts, warn)
	if err != nil {
		return nil, err
	}
	xr, err := xlsx.ConvertGrid(g, &xlsx.Options{Title: opts.Title, FontFS: opts.FontFS, FontDirs: opts.FontDirs, NoSystemFonts: opts.NoSystemFonts,
		SystemFonts: opts.SystemFonts, NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType,
		NoTextIndex: opts.NoTextIndex, Warn: opts.Warn})
	if err != nil {
		return nil, fmt.Errorf("parquet: %s", strings.TrimPrefix(err.Error(), "xlsx: "))
	}
	xr.Doc.Meta.Source = "parquet"
	res.Doc, res.EmbeddedFonts = xr.Doc, xr.EmbeddedFonts
	res.Warnings = append(warnings, xr.Warnings...)
	return res, nil
}

// table reads the rows shown into a grid: the column names, their types,
// and the values.
func table(r io.ReaderAt, size int64, opts *Options, warn func(string)) (*xlsx.Grid, *Result, error) {
	f, err := open(r, size)
	if err != nil {
		return nil, nil, err
	}
	top := f.root.children
	if len(top) > maxColumns {
		warn(fmt.Sprintf("only the first %d of %d columns are shown (the most a sheet has)", maxColumns, len(top)))
		top = top[:maxColumns]
	}
	limit := rowLimit(opts.Rows, len(top))
	res := &Result{Rows: f.meta.numRows, Cols: len(f.root.children), RowGroups: len(f.meta.rowGroups), CreatedBy: f.meta.createdBy}
	s := newSheet(f, top, limit, warn)
	s.read()
	res.Shown, res.Codecs = int64(len(s.rows)), s.codecs
	if total := f.meta.numRows; total > res.Shown && res.Shown == limit {
		msg := fmt.Sprintf("only the first %d of %d rows are shown", res.Shown, total)
		if limit < MaxRows {
			msg += " (the rows parameter shows more)"
		}
		warn(msg)
	}

	header := [][]xlsx.GridCell{make([]xlsx.GridCell, len(top))}
	for j, c := range top {
		header[0][j] = xlsx.GridCell{Text: c.name}
	}
	if !opts.NoTypes {
		types := make([]xlsx.GridCell, len(top))
		for j, c := range top {
			types[j] = xlsx.GridCell{Text: typeName(c)}
		}
		header = append(header, types)
	}
	g := &xlsx.Grid{Name: opts.Name, Rows: append(header, s.rows...), HeaderRows: len(header), SubHeader: !opts.NoTypes,
		Lang: language(header[0], s.rows), TableStyle: opts.TableStyle}
	return g, res, nil
}

// rowLimit is how many rows to read for a table of cols columns when
// Options.Rows asks for rows (0 the default, negative or too large for all).
// The default shows DefaultRows but no more than maxCells; any request is
// then held to maxSheetCells, since read() allocates rows×cols cells up front
// from a row count the file merely states.
func rowLimit(rows, cols int) int64 {
	limit := int64(rows)
	switch {
	case limit == 0:
		limit = DefaultRows
		if cols > 0 {
			limit = max(min(limit, maxCells/int64(cols)), 1)
		}
	case limit < 0 || limit > MaxRows:
		limit = MaxRows
	}
	if cols > 0 {
		limit = max(min(limit, maxSheetCells/int64(cols)), 1)
	}
	return limit
}

// sheet reads the rows shown into cells.
type sheet struct {
	f      *file
	top    []*node // the columns shown
	limit  int64
	warn   func(string)
	fmts   []format // of the leaves
	rows   [][]xlsx.GridCell
	codecs []string
	failed []bool // of the columns at the top: a warning was given
	// times holds the values of the columns of times and timestamps, which
	// are shown when their fractions of seconds are known
	times map[int]*timeColumn
}

type timeColumn struct {
	cells  []timeCell
	digits int
}

type timeCell struct {
	row int
	u   uint64
	x   uint32
}

func newSheet(f *file, top []*node, limit int64, warn func(string)) *sheet {
	s := &sheet{f: f, top: top, limit: limit, warn: warn, failed: make([]bool, len(top)), times: map[int]*timeColumn{}}
	for _, l := range f.leaves {
		s.fmts = append(s.fmts, formatOf(l))
	}
	return s
}

func (s *sheet) read() {
	cols := len(s.top)
	for gi := range s.f.meta.rowGroups {
		g := &s.f.meta.rowGroups[gi]
		n := min(g.numRows, s.limit-int64(len(s.rows)))
		if n <= 0 {
			if int64(len(s.rows)) >= s.limit {
				break
			}
			continue
		}
		base := len(s.rows)
		cells := make([]xlsx.GridCell, int(n)*cols)
		for i := range int(n) {
			s.rows = append(s.rows, cells[i*cols:(i+1)*cols:(i+1)*cols])
		}
		for j, top := range s.top {
			if top.leaf && top.rep != repRepeated {
				s.flat(g, j, top, base, int(n))
			} else {
				s.nested(g, j, top, base, int(n))
			}
		}
	}
	for j, tc := range s.times {
		f := s.fmts[s.top[j].col]
		for _, c := range tc.cells {
			s.rows[c.row][j] = xlsx.GridCell{Text: f.text(scalar{u: c.u, x: c.x}, tc.digits), Number: true}
		}
	}
}

func (s *sheet) fail(j int, err error) {
	if !s.failed[j] {
		s.failed[j] = true
		s.warn(fmt.Sprintf("column %q: %v", s.top[j].name, err))
	}
}

func (s *sheet) column(g *rowGroup, leaf *node) (*column, error) {
	cc := &g.columns[leaf.col]
	c, err := newColumn(s.f.r, s.f.dataEnd, leaf, cc)
	if err != nil {
		return nil, err
	}
	if name := codecNames[c.codec]; name != "" && !slices.Contains(s.codecs, name) {
		s.codecs = append(s.codecs, name)
	}
	return c, nil
}

// flat reads a column at the top that is not repeated: a value per row.
func (s *sheet) flat(g *rowGroup, j int, leaf *node, base, n int) {
	c, err := s.column(g, leaf)
	if err != nil {
		s.fail(j, err)
		return
	}
	f := s.fmts[leaf.col]
	var tc *timeColumn
	if f.timed() {
		if tc = s.times[j]; tc == nil {
			tc = &timeColumn{}
			s.times[j] = tc
		}
	}
	row := -1
	for {
		rep, def, v, err := c.next()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.fail(j, err)
			return
		}
		if rep == 0 {
			if row++; row >= n {
				return
			}
		}
		if row < 0 || def != leaf.defLevel || f.kind == fNull {
			continue
		}
		if tc != nil {
			tc.cells = append(tc.cells, timeCell{base + row, v.u, v.x})
			tc.digits = max(tc.digits, f.digits(v))
			continue
		}
		s.rows[base+row][j] = xlsx.GridCell{Text: cellText(f.text(v, -1)), Number: f.number()}
	}
	if row+1 < n {
		s.fail(j, errColumnEnd)
	}
}

// nested reads a nested field at the top: the entries of its leaves are
// put together into the values of the rows, which are then written.
func (s *sheet) nested(g *rowGroup, j int, top *node, base, n int) {
	vals := make([]any, n)
	for _, leaf := range leavesOf(top) {
		c, err := s.column(g, leaf)
		if err != nil {
			s.fail(j, err)
			return
		}
		a := newAssembler(top, leaf)
		row := -1
		for {
			rep, def, v, err := c.next()
			if err == io.EOF {
				break
			}
			if err != nil {
				s.fail(j, err)
				return
			}
			if rep == 0 {
				if row++; row >= n {
					break
				}
			}
			if row < 0 || !a.put(&vals[row], rep, def, v) {
				s.fail(j, errEncoding)
				return
			}
		}
		if row+1 < n {
			s.fail(j, errColumnEnd)
			return
		}
	}
	w := &writer{fmts: s.fmts}
	for i, v := range vals {
		if v == nil && top.rep != repRepeated {
			continue
		}
		w.b.Reset()
		w.value(top, v, false)
		s.rows[base+i][j] = xlsx.GridCell{Text: cellText(w.b.String())}
	}
}

// language guesses the language of the text for its fonts and its Han
// characters: that of the kana or hangul in the column names and values.
func language(names []xlsx.GridCell, rows [][]xlsx.GridCell) string {
	kana, hangul, n := 0, 0, 0
	for _, cells := range slices.Concat([][]xlsx.GridCell{names}, rows[:min(len(rows), 1<<16)]) {
		for _, c := range cells {
			if c.Number {
				continue
			}
			for _, r := range c.Text {
				switch {
				case r < 0x3000:
				case unicode.In(r, unicode.Hiragana, unicode.Katakana):
					kana++
				case unicode.Is(unicode.Hangul, r):
					hangul++
				}
			}
			if n += len(c.Text); n > 1<<20 {
				break
			}
		}
		if n > 1<<20 {
			break
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
