// Package equation lays out mathematical formulas and draws them into BDF
// objects, for every converter whose documents hold formulas: Office Math
// (OMML) in Word documents, MathML in HTML pages, and LaTeX in Markdown.
//
// The three notations are read into one tree (the shape of presentation
// MathML: rows of tokens, fractions, radicals, scripts, limits, tables),
// which is laid out the way TeX and MathML Core lay formulas out, with the
// parameters and glyph variants of the OpenType MATH table of a formula
// font (STIX Two Math, Cambria Math, Latin Modern Math …): script sizes and
// shifts, fraction and radical gaps, the spacing between atoms, large
// operators, and delimiters, radicals and accents that grow from larger
// variants and assemblies of parts. Glyphs no character maps to are drawn
// through private use characters of the embedded font subset
// (fontset.Set.GlyphRune). Without a formula font, formulas are laid out
// with the text font and default parameters, and delimiters are scaled.
//
// A laid out formula is drawn into a child object that its host places with
// an ALT_TEXT carrying the formula in a linear notation, such as
// "x=(−b±√(b^2−4ac))/(2a)", which text extraction, search and copy use.
package equation

import "github.com/shibukawa/bdf"

// Node is a node of a formula.
type Node interface{ isNode() }

func (*Row) isNode()       {}
func (*Atom) isNode()      {}
func (*Frac) isNode()      {}
func (*Radical) isNode()   {}
func (*Scripts) isNode()   {}
func (*UnderOver) isNode() {}
func (*Table) isNode()     {}
func (*Space) isNode()     {}
func (*Styled) isNode()    {}
func (*Enclose) isNode()   {}
func (*Phantom) isNode()   {}
func (*Bar) isNode()       {}

// Row is a sequence of nodes laid out side by side (mrow, a TeX group).
// Its vertically stretchy operators grow to the height of the rest.
type Row struct {
	Kids []Node
	// Class is how the row spaces from its neighbors; ClassAuto makes it
	// an ordinary atom.
	Class Class
}

// Kind is the kind of a token.
type Kind uint8

const (
	Ident  Kind = iota // a variable or function name (mi)
	Number             // mn
	Op                 // an operator, relation, delimiter or punctuation (mo)
	Text               // text (mtext, \text{}, Office's normal text)
)

// Class is a TeX atom class, which sets the spacing between atoms.
type Class uint8

const (
	ClassAuto Class = iota // from the token (operators: the operator dictionary)
	Ord
	LargeOp // a large operator or a function name
	Bin
	Rel
	Open
	Close
	Punct
	Inner
	None // no spacing around it
)

// Variant is a math variant: the style of the letters and digits of a
// token, drawn with the Mathematical Alphanumeric Symbols.
type Variant uint8

const (
	VarAuto Variant = iota // italic for single-letter identifiers, else normal
	VarNormal
	VarItalic
	VarBold
	VarBoldItalic
	VarDoubleStruck
	VarScript
	VarBoldScript
	VarFraktur
	VarBoldFraktur
	VarSansSerif
	VarSansSerifBold
	VarSansSerifItalic
	VarSansSerifBoldItalic
	VarMonospace
)

// Atom is a token: identifiers, numbers, operators and text.
type Atom struct {
	Kind    Kind
	Text    string
	Variant Variant
	Class   Class
	// Font is the family of a Text atom drawn in a text font instead of
	// the formula font ("" draws it in the formula font when that has its
	// characters).
	Font string

	// Operator properties. Stretchy operators grow to the height of their
	// row (vertical ones) or the width of what they are over or under
	// (horizontal ones); Symmetric ones grow evenly around the math axis.
	Stretchy, Symmetric bool
	LargeOp             bool // drawn larger in display style
	MovableLimits       bool // limits go to the side outside display style
	Accent              bool // an accent over or under its base
	// Shortfall grows delimiters as TeX's \left and \right do (not quite
	// to the full height of what they enclose).
	Shortfall bool
	// MinSize is a size (em) the operator grows to at least, as \big and
	// \Bigg ask for.
	MinSize float64
	// Spaced sets the space before and after an operator to LSpace and
	// RSpace (em) instead of the space of its class.
	Spaced         bool
	LSpace, RSpace float64
}

// NewOp returns an operator token with the properties the operator
// dictionary gives it (stretchy delimiters, large operators, accents).
func NewOp(s string) *Atom {
	p := lookupOp(normalizeOp(s))
	return &Atom{Kind: Op, Text: s, Stretchy: p.stretchy || p.horizontal, Symmetric: p.symmetric,
		LargeOp: p.largeOp, MovableLimits: p.movable, Accent: p.accent}
}

// FracKind is how a fraction is drawn.
type FracKind uint8

const (
	FracStacked FracKind = iota // numerator over denominator
	FracSkewed                  // numerator raised, a slash, denominator lowered
	FracLinear                  // numerator, slash, denominator on the baseline
)

// Align is a horizontal alignment.
type Align uint8

const (
	Center Align = iota
	Left
	Right
)

// Frac is a fraction (mfrac, \frac, \binom's stack, OMML f).
type Frac struct {
	Num, Den           Node
	Kind               FracKind
	NoBar              bool    // a stack without a rule
	Thickness          float64 // the rule in em; 0 for the font's
	NumAlign, DenAlign Align
}

// Radical is a square root (Degree nil) or a root of another degree.
type Radical struct {
	Base, Degree Node
}

// Scripts is a base with subscripts and superscripts after it and before
// it (msub, msup, msubsup, mmultiscripts; nil for none).
type Scripts struct {
	Base           Node
	Sub, Sup       Node
	PreSub, PreSup Node
}

// UnderOver is a base with something under and over it (limits, accents,
// braces).
type UnderOver struct {
	Base                    Node
	Under, Over             Node
	AccentUnder, AccentOver bool
}

// Table is a matrix, the rows of an array of equations, or cases.
type Table struct {
	Rows     [][]Node
	ColAlign []Align // per column, the last repeated; nil centers
	// Aligned lays the columns out in pairs, the first of a pair aligned
	// right and the second left, without space between them (LaTeX's
	// aligned, Office's equation arrays).
	Aligned bool
	RowGap  float64 // em between rows; 0 for the default
	ColGap  float64 // em between columns; 0 for the default
	// Display sets the cells in display style (equation arrays).
	Display bool
}

// Space is horizontal space (negative: backwards), in em.
type Space struct {
	Width float64
}

// Styled changes the style of what it holds; zero fields leave it as it
// is.
type Styled struct {
	Kid Node
	// Display sets display (1) or inline (2) style.
	Display int8
	// Level sets the script level (1+level: 1 is display/text size, 2
	// script, 3 script-script); LevelUp makes it one smaller.
	Level   int8
	LevelUp bool
	Color   *bdf.Color
	Size    float64 // font size in points
	// Base is the font size in points of text outside scripts (Office's
	// run size, which scripts are reduced from).
	Base    float64
	Scale   float64 // font size factor
	Variant Variant
	// Bold draws the text and operators bold where the font has bold
	// forms (Office's bold runs).
	Bold bool
}

// Enclose draws lines around or across what it holds (menclose, OMML
// borderBox, \boxed, \cancel).
type Enclose struct {
	Kid                      Node
	Top, Bottom, Left, Right bool
	Round, Circle            bool
	StrikeH, StrikeV         bool
	StrikeUp, StrikeDown     bool // lower left to upper right, upper left to lower right
}

// Phantom lays out what it holds without drawing it (Show false), or
// draws it without some of its dimensions (\smash, OMML phant).
type Phantom struct {
	Kid                          Node
	Show                         bool
	ZeroWidth, ZeroAsc, ZeroDesc bool
}

// Bar draws a rule over (or under) what it holds (\overline, OMML bar).
type Bar struct {
	Kid   Node
	Under bool
}
