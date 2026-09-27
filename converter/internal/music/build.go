package music

import (
	"io/fs"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
)

// Options controls how a score is engraved.
type Options struct {
	// Source is the input format recorded in the document ("mml", "midi",
	// "musicxml").
	Source string
	// Title overrides the document title (and the title printed above the
	// music).
	Title string
	// Pages selects 1-based pages (see converter.Pages); nil keeps all.
	Pages conv.Pages

	// PageWidth and PageHeight are the paper size in pt (default A4
	// portrait); StaffSpace is the distance between staff lines in pt
	// (default 1.75 mm, a 7 mm staff).
	PageWidth, PageHeight float64
	StaffSpace            float64

	// NoPlay leaves the sound out: no seq part and no cues.
	NoPlay bool

	// The fonts of the text (title, lyrics, words), as in
	// converter.Options.
	FontFS        fs.FS
	FontDirs      []string
	NoSystemFonts bool
	SystemFonts   bool
	NoSubset      bool
	NoWOFF2       bool
	IgnoreFSType  bool
	NoTextIndex   bool

	// Warn receives non-fatal problems; when nil they are collected in
	// Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of engraving.
type Result struct {
	Doc      *bdf.Document
	Warnings []string
	// Pages, Systems and Measures count what was engraved; Parts is the
	// number of parts and Staves the staves of a system.
	Pages, Systems, Measures, Parts, Staves int
	// Notes is the number of notes drawn.
	Notes int
	// Seconds is the length of the music played.
	Seconds       float64
	EmbeddedFonts int
}

// Summary describes the result for the command line.
func (r *Result) Summary() string {
	return summary(r)
}

// Build engraves a score into a document with one fixed view of pages.
func Build(s *Score, o Options) (*Result, error) {
	return build(s, &o)
}
