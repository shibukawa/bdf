// Package equation is the formula engine (internal/mathlayout) for the
// converters: it lays formulas out with the fonts of a document
// (fontset.Set), which keeps the characters they draw for the embedded
// subsets, and draws them into BDF objects as text.
//
// A laid out formula is drawn into a child object that its host places with
// an ALT_TEXT carrying the formula in a linear notation, such as
// "x=(−b±√(b^2−4ac))/(2a)", which text extraction, search and copy use.
// Glyphs no character maps to are drawn through private use characters of
// the embedded font subset (fontset.Set.GlyphRune).
package equation

import (
	"golang.org/x/net/html"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/internal/mathlayout"
)

// The formula tree, its readers and the layout are those of mathlayout.
type (
	Engine    = mathlayout.Engine
	Box       = mathlayout.Box
	Style     = mathlayout.Style
	OMML      = mathlayout.OMML
	Node      = mathlayout.Node
	Row       = mathlayout.Row
	Atom      = mathlayout.Atom
	Frac      = mathlayout.Frac
	Radical   = mathlayout.Radical
	Scripts   = mathlayout.Scripts
	UnderOver = mathlayout.UnderOver
	Table     = mathlayout.Table
	Space     = mathlayout.Space
	Styled    = mathlayout.Styled
	Enclose   = mathlayout.Enclose
	Phantom   = mathlayout.Phantom
	Bar       = mathlayout.Bar
)

// ParseTeX reads a formula written in LaTeX's math mode.
func ParseTeX(src string) Node { return mathlayout.ParseTeX(src) }

// ParseMathML reads a math element.
func ParseMathML(n *html.Node) (formula Node, display bool) { return mathlayout.ParseMathML(n) }

// Linear returns a formula in a linear notation.
func Linear(n Node) string { return mathlayout.Linear(n) }

// ParseColor reads a color of a formula (a name or #rrggbb).
func ParseColor(s string) (bdf.Color, bool) { return mathlayout.ParseColor(s) }

// IsFunctionName reports whether s is the name of a function (sin, log …).
func IsFunctionName(s string) bool { return mathlayout.IsFunctionName(s) }

// Fonts are what formulas are laid out with.
type Fonts struct {
	Set *fontset.Set
	// Math is the formula font asked for ("" for the default): a font with
	// a MATH table resolves it (see fontdb.DB.ResolveMath).
	Math string
	// Text returns the face that draws r in a text font: for characters
	// the formula font lacks (East Asian text) and for text asked in a
	// family (family "" is the text around the formula). It is also what
	// formulas are drawn with when no formula font is available.
	Text func(family string, r rune, bold, italic bool) *fontset.Choice
	// Warn receives problems, such as a missing formula font.
	Warn func(msg string)
}

// New returns an engine that lays formulas out with f. It is not safe for
// concurrent use.
func New(f Fonts) *Engine {
	// one face of the engine for a face of the set: the engine tells its
	// faces apart by their pointers
	faces := map[*fontset.Choice]*mathlayout.Face{}
	face := func(fc *fontset.Choice) *mathlayout.Face {
		if fc == nil {
			return nil
		}
		mf := faces[fc]
		if mf == nil {
			mf = &mathlayout.Face{Loaded: fc.Loaded, Host: fc}
			faces[fc] = mf
		}
		return mf
	}
	mf := mathlayout.Fonts{Warn: f.Warn}
	if f.Set != nil {
		if fc := f.Set.ChooseMath(f.Math); fc != nil {
			f.Set.Pin(fc)
			mf.Math = face(fc)
		}
		mf.Advance = func(fa *mathlayout.Face, r rune) float64 { return f.Set.Advance(choice(fa), r) }
		mf.GlyphRune = func(fa *mathlayout.Face, g uint16) (rune, bool) { return f.Set.GlyphRune(choice(fa), g) }
	} else {
		mf.Advance = func(*mathlayout.Face, rune) float64 { return 0.5 }
	}
	if f.Text != nil {
		mf.Text = func(family string, r rune, bold, italic bool) *mathlayout.Face {
			return face(f.Text(family, r, bold, italic))
		}
	}
	return mathlayout.New(mf)
}

// choice returns the face of the set that a face of the engine stands for.
func choice(f *mathlayout.Face) *fontset.Choice {
	fc, _ := f.Host.(*fontset.Choice)
	return fc
}

// XHeight returns the height of the x of a face in em (0 when it has
// none).
func XHeight(fc *fontset.Choice) float64 {
	if fc == nil {
		return 0
	}
	return mathlayout.XHeight(&mathlayout.Face{Loaded: fc.Loaded})
}
