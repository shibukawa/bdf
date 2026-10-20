// Package textpart reads and writes the print image of a form filled from
// a template of plain text: the part of the type bdf.PartText ("text") of
// a document, as JSON (docs/spec.md §4).
//
// A template of text (converter/textform) puts its characters on a grid
// of columns and lines, as a dot matrix printer prints in its text mode.
// The pages of the document draw that grid with a monospace font; this
// part keeps the grid itself, so that the writers of the text modes of
// printers (ESC/P, ESC/POS) and of plain text print the characters and
// not the pictures of the pages. It is small and imports neither package
// form nor the converters, for the package print to read.
//
// The JSON of version 1:
//
//	{
//	  "version": 1,
//	  "cpi": 10, "lpi": 6,            // columns and lines to the inch
//	  "columns": 136, "lines": 66,    // the paper, in columns and lines (lines 0: a roll)
//	  "margin": {"left": 2, "top": 1},
//	  "skip": 3,                      // lines skipped over the perforation
//	  "pages": [
//	    {"lines": [
//	      {"text": "納品書", "size": [{"col": 30, "cols": 12, "size": "big"}]},
//	      {"text": ""},
//	      {"text": "┌────┬──┐", "paper": [{"col": 0, "cols": 9}]},
//	      {"text": "│ ペン │ 2 │", "paper": [{"col": 0, "cols": 1}, {"col": 7, "cols": 1}, {"col": 11, "cols": 1}]}
//	    ]}
//	  ]
//	}
//
// Columns are counted from the left margin: column 0 of a line is the
// first column after it, and line 0 of a page the first line under the top
// margin. A character takes the columns that Width gives it, twice as
// many in a range of the size "wide" or "big"; Cells walks a line with
// the columns of its characters.
package textpart

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
	"golang.org/x/text/width"
)

// Version is the version of the JSON this package writes and reads.
const Version = 1

// Sizes of characters (Range.Size).
const (
	// SizeWide is double width: ESC W 1 of ESC/P, GS ! 0x10 of ESC/POS.
	SizeWide = "wide"
	// SizeTall is double height: ESC w 1, GS ! 0x01.
	SizeTall = "tall"
	// SizeBig is double width and height.
	SizeBig = "big"
)

// Limits of a part (it is read from documents that are not trusted).
const (
	// MaxSize bounds the bytes of the JSON.
	MaxSize = 64 << 20
	// MaxColumns and MaxLines bound the paper.
	MaxColumns = 4096
	MaxLines   = 4096
	// MaxPages bounds the pages, and MaxPageLines the lines of a page (a
	// roll of paper has as many as it has).
	MaxPages     = 100000
	MaxPageLines = 1 << 20
)

// Image is the print image of a form: the grid and the lines of its pages.
type Image struct {
	// Version is Version.
	Version int `json:"version"`
	// CPI is the width of a column: columns to the inch (10, 12, 15, 17.1
	// …). LPI is the height of a line: lines to the inch (6, 8 …).
	CPI float64 `json:"cpi"`
	LPI float64 `json:"lpi"`
	// Columns is the width of the paper in columns, the left margin
	// included (80; 136 for paper of 15 inches).
	Columns int `json:"columns"`
	// Lines is the length of a page in lines, the top margin and the skip
	// included (66 for 11 inches at 6 LPI). 0 is a roll of paper (a
	// receipt): each page is as long as its lines.
	Lines int `json:"lines"`
	// Margin is the room left of the lines (in columns) and above them (in
	// lines) on each page.
	Margin Margin `json:"margin"`
	// Skip is the lines left empty at the bottom of each page, over the
	// perforation of fanfold paper (ESC N of ESC/P). Lines - Margin.Top -
	// Skip lines are printed on a page.
	Skip int `json:"skip,omitempty"`
	// AmbiguousWide says that the characters of East Asian width
	// "ambiguous" (box drawing, ○, ※, Greek, Cyrillic …) take two
	// columns; without it they take one (see Width).
	AmbiguousWide bool `json:"ambiguousWide,omitempty"`
	// Font is the family of the monospace font that the pages are drawn
	// with ("" for none named).
	Font string `json:"font,omitempty"`
	// Pages are the pages of the form, in the order of the pages of the
	// first view of the document.
	Pages []Page `json:"pages"`
}

// Margin is the room around the lines of a page.
type Margin struct {
	// Left is in columns, Top in lines.
	Left int `json:"left"`
	Top  int `json:"top"`
}

// Page is a page of the form.
type Page struct {
	// Lines are the lines of the page from the top margin down, a line
	// for each line of the grid: Lines[i] is line i. Lines past the last
	// one are empty.
	Lines []Line `json:"lines"`
}

// Line is a line of a page.
type Line struct {
	// Text is the characters of the line from its first column (the left
	// margin is not in it), without spaces at its end. A character of a
	// double width range is in it once.
	Text string `json:"text"`
	// Paper are the ranges of the line that are of the preprinted paper
	// (the rules and the names of a printed form): they are drawn on the
	// page, as a layer of the role "background", and are not printed on
	// paper that has them already. They do not overlap and are in order.
	Paper []Range `json:"paper,omitempty"`
	// Size are the ranges of the line whose characters are printed double
	// in size (Range.Size). They do not overlap and are in order. A line
	// with a range of the size "tall" or "big" takes two lines of the
	// page: its own and the one under it, whose Text is empty. All of its
	// characters stand on the base line of the lower one (as ESC/POS
	// prints a line with characters of double height).
	Size []Range `json:"size,omitempty"`
}

// Range is the columns of a part of a line.
type Range struct {
	// Col is the first column of the range, and Cols the number of its
	// columns as printed (a character of double width takes twice its
	// columns).
	Col  int `json:"col"`
	Cols int `json:"cols"`
	// Size is the size of the characters of a range of Line.Size:
	// SizeWide, SizeTall or SizeBig. A range of Line.Paper has none.
	Size string `json:"size,omitempty"`
}

// Tall reports whether a line takes two lines of the page.
func (l Line) Tall() bool {
	for _, r := range l.Size {
		if r.Size == SizeTall || r.Size == SizeBig {
			return true
		}
	}
	return false
}

// Width is the columns of a character: 2 for a wide or full-width one
// (East Asian width W or F: kanji, kana, full-width forms, most emoji), and
// for an ambiguous one (A) when ambiguousWide is set; 0 for a character
// that takes no column of its own (a combining mark, a variation selector,
// a zero-width space or joiner, a control character); 1 for the others.
func Width(r rune, ambiguousWide bool) int {
	switch {
	case r < 0x20 || r >= 0x7f && r < 0xa0:
		return 0
	case r < 0x7f:
		return 1
	case unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || unicode.Is(unicode.Cf, r) ||
		r >= 0xfe00 && r <= 0xfe0f || r >= 0xe0100 && r <= 0xe01ef || r == 0x200b:
		return 0
	}
	switch width.LookupRune(r).Kind() {
	case width.EastAsianWide, width.EastAsianFullwidth:
		return 2
	case width.EastAsianAmbiguous:
		if ambiguousWide {
			return 2
		}
	}
	return 1
}

// Cell is a character of a line with its place: a character and the
// characters of no width that follow it (combining marks).
type Cell struct {
	// Text is the character.
	Text string
	// Col is its first column and Cols its columns as printed: Width of
	// the character, twice that in a range of the size "wide" or "big".
	Col, Cols int
	// Size is the size of the range the character is in ("" for none),
	// and Paper whether it is in a range of the paper.
	Size  string
	Paper bool
}

// Cells returns the characters of a line of the image with their columns.
// Spaces are characters too.
func (im *Image) Cells(l Line) []Cell {
	var out []Cell
	col := 0
	for i := 0; i < len(l.Text); {
		r, n := utf8.DecodeRuneInString(l.Text[i:])
		w := Width(r, im.AmbiguousWide)
		if w == 0 && len(out) > 0 {
			out[len(out)-1].Text += l.Text[i : i+n]
			i += n
			continue
		}
		c := Cell{Text: l.Text[i : i+n], Col: col, Cols: w}
		for _, rg := range l.Size {
			if col >= rg.Col && col < rg.Col+rg.Cols {
				c.Size = rg.Size
			}
		}
		if c.Size == SizeWide || c.Size == SizeBig {
			c.Cols *= 2
		}
		for _, rg := range l.Paper {
			if col >= rg.Col && col < rg.Col+rg.Cols {
				c.Paper = true
			}
		}
		out = append(out, c)
		col += c.Cols
		i += n
	}
	return out
}

// Rows is the lines that a page prints: those of the paper without the top
// margin and the skip, or, on a roll (Lines 0), those of the page.
func (im *Image) Rows(p Page) int {
	if im.Lines == 0 {
		return len(p.Lines)
	}
	return im.Lines - im.Margin.Top - im.Skip
}

// TextOptions controls Image.Text.
type TextOptions struct {
	// Margins writes the margins: spaces before each line, and empty
	// lines above the first.
	Margins bool
	// NoPaper writes the ranges of the paper as spaces, as they are
	// printed on paper that has them already.
	NoPaper bool
	// CRLF ends the lines with CR LF instead of LF.
	CRLF bool
}

// Text writes the image as plain text: its lines, ended by LF, and its
// pages separated by a form feed (FF). A character of double width is
// followed by the spaces of its other columns, so that the columns of the
// line stay those of the grid; the second line of a line of double height
// is an empty line.
func (im *Image) Text(o TextOptions) string {
	nl := "\n"
	if o.CRLF {
		nl = "\r\n"
	}
	var b strings.Builder
	for pi, p := range im.Pages {
		if pi > 0 {
			b.WriteByte('\f')
		}
		if o.Margins {
			for range im.Margin.Top {
				b.WriteString(nl)
			}
		}
		for _, l := range p.Lines {
			var line strings.Builder
			if o.Margins {
				line.WriteString(strings.Repeat(" ", im.Margin.Left))
			}
			for _, c := range im.Cells(l) {
				w := Width(firstRune(c.Text), im.AmbiguousWide)
				switch {
				case o.NoPaper && c.Paper:
					line.WriteString(strings.Repeat(" ", c.Cols))
				default:
					line.WriteString(c.Text)
					line.WriteString(strings.Repeat(" ", c.Cols-w))
				}
			}
			b.WriteString(strings.TrimRight(line.String(), " "))
			b.WriteString(nl)
		}
	}
	return b.String()
}

func firstRune(s string) rune {
	r, _ := utf8.DecodeRuneInString(s)
	return r
}

// Encode writes the image as the JSON of the part.
func (im *Image) Encode() ([]byte, error) {
	if im.Version == 0 {
		im.Version = Version
	}
	if err := im.check(); err != nil {
		return nil, err
	}
	return json.Marshal(im)
}

// Decode reads the JSON of a part, and checks it.
func Decode(b []byte) (*Image, error) {
	if len(b) > MaxSize {
		return nil, fmt.Errorf("textpart: the part is larger than %d MiB", MaxSize>>20)
	}
	im := &Image{}
	if err := json.Unmarshal(b, im); err != nil {
		return nil, fmt.Errorf("textpart: %w", err)
	}
	if im.Version != Version {
		return nil, fmt.Errorf("textpart: the part is of version %d; this reads version %d", im.Version, Version)
	}
	if err := im.check(); err != nil {
		return nil, err
	}
	return im, nil
}

// check refuses an image whose numbers are out of range.
func (im *Image) check() error {
	bad := func(format string, args ...any) error {
		return fmt.Errorf("textpart: "+format, args...)
	}
	switch {
	case !(im.CPI > 0 && im.CPI <= 100) || !(im.LPI > 0 && im.LPI <= 100):
		return bad("the pitch %g CPI, %g LPI is out of range", im.CPI, im.LPI)
	case im.Columns < 1 || im.Columns > MaxColumns:
		return bad("%d columns are out of range", im.Columns)
	case im.Lines < 0 || im.Lines > MaxLines:
		return bad("%d lines are out of range", im.Lines)
	case im.Margin.Left < 0 || im.Margin.Left >= im.Columns || im.Margin.Top < 0 || im.Skip < 0:
		return bad("the margins are out of range")
	case im.Lines > 0 && im.Margin.Top+im.Skip >= im.Lines:
		return bad("the margins leave no line of the page")
	case len(im.Pages) > MaxPages:
		return bad("more than %d pages", MaxPages)
	}
	for pi, p := range im.Pages {
		if len(p.Lines) > MaxPageLines || im.Lines > 0 && len(p.Lines) > im.Rows(p) {
			return bad("page %d has %d lines", pi+1, len(p.Lines))
		}
		for li, l := range p.Lines {
			for _, rs := range [][]Range{l.Paper, l.Size} {
				end := 0
				for _, r := range rs {
					if r.Col < end || r.Cols < 1 || r.Cols > MaxColumns*2 || r.Col > MaxColumns*2 {
						return bad("page %d, line %d: the ranges overlap or are out of range", pi+1, li+1)
					}
					end = r.Col + r.Cols
				}
			}
			for _, r := range l.Size {
				if r.Size != SizeWide && r.Size != SizeTall && r.Size != SizeBig {
					return bad("page %d, line %d: the size %q", pi+1, li+1, r.Size)
				}
			}
		}
	}
	return nil
}

// ErrNone tells that a document has no part of the print image.
var ErrNone = errors.New("textpart: the document has no print image of text")

// Read reads the print image of an open document: its first part of the
// type bdf.PartText. It returns ErrNone when it has none.
func Read(rd *bdf.Reader) (*Image, error) {
	if rd.Locked() {
		return nil, fmt.Errorf("textpart: the document is encrypted")
	}
	for _, e := range rd.Manifest.Parts {
		if e.T != bdf.PartText {
			continue
		}
		if e.Size < 0 || e.Size > MaxSize || e.Len < 0 || e.Len > MaxSize {
			return nil, fmt.Errorf("textpart: the part is larger than %d MiB", MaxSize>>20)
		}
		b, err := rd.Part(e.H)
		if err != nil {
			return nil, fmt.Errorf("textpart: %w", err)
		}
		return Decode(b)
	}
	return nil, ErrNone
}

// FromDocument reads the print image of a document being written. It
// returns ErrNone when it has none.
func FromDocument(doc *bdf.Document) (*Image, error) {
	for _, p := range doc.Parts() {
		if p.Type == bdf.PartText {
			return Decode(p.Data)
		}
	}
	return nil, ErrNone
}

// AddTo adds the image to a document as its part of the print image.
func (im *Image) AddTo(doc *bdf.Document) (bdf.Hash, error) {
	b, err := im.Encode()
	if err != nil {
		return bdf.Hash{}, err
	}
	return doc.AddPart(bdf.PartText, b), nil
}
