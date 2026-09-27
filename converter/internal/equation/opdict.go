package equation

import (
	"strings"
	"unicode/utf8"
)

// opProps are what the operator dictionary says about an operator.
type opProps struct {
	class               Class
	stretchy, symmetric bool
	horizontal          bool // stretches horizontally (arrows, accents, braces)
	largeOp, movable    bool
	accent              bool
}

const (
	relChars = "=<>:≤≥≠≈≡≢∼≃≄≅≆≇≉≊≍≎≏≐≑≒≓≔≕≖≗≘≙≚≛≜≝≞≟∝∈∉∊∋∌∍⊂⊃⊄⊅⊆⊇⊈⊉⊊⊋⊏⊐⊑⊒" +
		"⊢⊣⊥⊦⊧⊨⊩⊪⊫⊬⊭⊮⊯≪≫≺≻≼≽≾≿⪯⪰⋘⋙≲≳≶≷⋖⋗≦≧⩽⩾⪅⪆⪋⪌∣∤∥∦≬⋈⌢⌣⊲⊳⊴⊵∽⋍⋐⋑⊸⟂" +
		"≮≯≰≱⊀⊁⋠⋡⋪⋫⋬⋭⋦⋧⋨⋩≨≩⪇⪈⪉⪊⋚⋛⪋⪌∺∻≁≭⊊⊋⫅⫆⫋⫌∷"
	binChars   = "+−±∓×÷·∗∘∙⋅∪∩∧∨⊕⊖⊗⊘⊙⊚⊛⊝⊎⊓⊔∖⋆⋄◇△▽◁▷▹◃≀⨿†‡⋉⋊⋋⋌⊞⊟⊠⊡⋒⋓⊻⊼⊽⋎⋏⊺⋇⨯⨝⩞"
	openChars  = "([{⟨⌈⌊⟦⦃⦅⟮⟬⦇⦉⌜⌞〈"
	closeChars = ")]}⟩⌉⌋⟧⦄⦆⟯⟭⦈⦊⌝⌟〉"
	// fence characters that open or close depending on where they are
	fenceChars = "|‖∣∥"
	punctChars = ",;"
	postChars  = "!?′″‴⁗'"
	largeOps   = "∑∏∐⋀⋁⋂⋃⨀⨁⨂⨃⨄⨅⨆⨉⫿⫼"
	integrals  = "∫∬∭∮∯∰∱∲∳⨋⨌⨍⨎⨏⨐⨑⨒⨓⨔⨕⨖⨗⨘⨙⨚⨛⨜"
	// vertically stretchy operators besides the brackets and fences
	vertStretchy = "/\\∕∖↑↓↕⇑⇓⇕⎰⎱⟅⟆"
	// horizontally stretchy operators: arrows, and what goes over and under
	horizStretchy = "←→↔↚↛↮⇐⇒⇔⇍⇎⇏↦↤↩↪↼↽⇀⇁⇄⇆⇋⇌⟵⟶⟷⟸⟹⟺⟼↞↠" +
		"^ˆ̂~˜̃¯‾̄̅_̲ˇ̌⏞⏟⏜⏝⎴⎵⏠⏡⃗⃖⃡=−"
	// characters that are accents when over or under a base
	accentChars = "^ˆ~˜¯‾ˇ˘˙¨´`˚→←↔⃗⏞⏟⏜⏝⎴⎵⏠⏡_"
)

// functionNames are the names of functions set upright with operator
// spacing, and those whose limits go under them (lim, max …).
var functionNames = map[string]bool{
	"sin": false, "cos": false, "tan": false, "cot": false, "sec": false, "csc": false,
	"arcsin": false, "arccos": false, "arctan": false, "arccot": false, "arcsec": false, "arccsc": false,
	"sinh": false, "cosh": false, "tanh": false, "coth": false, "sech": false, "csch": false,
	"arsinh": false, "arcosh": false, "artanh": false,
	"log": false, "ln": false, "lg": false, "exp": false, "arg": false, "deg": false, "dim": false,
	"hom": false, "ker": false, "det": true, "gcd": true, "Pr": true, "sgn": false, "tr": false, "Tr": false,
	"lim": true, "liminf": true, "limsup": true, "max": true, "min": true, "sup": true, "inf": true,
	"mod": false, "Re": false, "Im": false, "rank": false, "diag": false, "span": false,
}

// IsFunctionName reports whether s names a function (sin, log, lim …)
// that is set upright.
func IsFunctionName(s string) bool {
	_, ok := functionNames[s]
	return ok
}

// lookupOp returns the dictionary properties of an operator.
func lookupOp(s string) opProps {
	r, n := utf8.DecodeRuneInString(s)
	if n != len(s) || n == 0 {
		// several characters: a function name (lim, max …) or a word
		if mov, ok := functionNames[s]; ok {
			return opProps{class: LargeOp, movable: mov}
		}
		return opProps{class: Ord}
	}
	in := func(set string) bool { return strings.ContainsRune(set, r) }
	p := opProps{class: Ord}
	switch {
	case in(relChars) || isArrow(r):
		p.class = Rel
	case in(binChars):
		p.class = Bin
	case in(openChars):
		p.class, p.stretchy, p.symmetric = Open, true, true
	case in(closeChars):
		p.class, p.stretchy, p.symmetric = Close, true, true
	case in(fenceChars):
		p.class, p.stretchy, p.symmetric = Ord, true, true
	case in(punctChars):
		p.class = Punct
	case in(postChars):
		p.class = Close
	case in(largeOps):
		p.class, p.largeOp, p.movable, p.symmetric = LargeOp, true, true, true
	case in(integrals):
		p.class, p.largeOp, p.symmetric = LargeOp, true, true
	case r >= 0x2061 && r <= 0x2064: // invisible function application, times, separator, plus
		p.class = None
	}
	if in(vertStretchy) {
		p.stretchy = true
	}
	if in(horizStretchy) || isArrow(r) && !in(vertStretchy) {
		p.horizontal = true
	}
	if in(accentChars) || isCombining(r) {
		p.accent = true
	}
	return p
}

// isArrow reports whether r is in the arrow blocks.
func isArrow(r rune) bool {
	return r >= 0x2190 && r <= 0x21FF || r >= 0x27F0 && r <= 0x27FF || r >= 0x2900 && r <= 0x297F || r >= 0x2B00 && r <= 0x2B11
}

// isCombining reports whether r is a combining mark (an accent over the
// character before it).
func isCombining(r rune) bool {
	return r >= 0x0300 && r <= 0x036F || r >= 0x20D0 && r <= 0x20FF
}

// combiningAccent returns the combining form of a spacing accent, whose
// glyph in a formula font sits at the height of the accent over a
// lower-case letter.
func combiningAccent(r rune) rune {
	switch r {
	case '^', 'ˆ':
		return 0x302
	case '~', '˜':
		return 0x303
	case '¯', '‾':
		return 0x304
	case '˘':
		return 0x306
	case '˙':
		return 0x307
	case '¨':
		return 0x308
	case '˚':
		return 0x30A
	case 'ˇ':
		return 0x30C
	case '´':
		return 0x301
	case '`':
		return 0x300
	case '→', '⃗':
		return 0x20D7
	case '←':
		return 0x20D6
	case '↔':
		return 0x20E1
	}
	return r
}

// normalizeOp replaces the ASCII stand-ins of operators with the
// characters formulas set: the hyphen-minus is a minus sign, the asterisk
// an asterisk operator, the apostrophe a prime.
func normalizeOp(s string) string {
	switch s {
	case "-":
		return "−"
	case "*":
		return "∗"
	case "'":
		return "′"
	case "''":
		return "″"
	}
	return s
}
