package xlsx

import (
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Number formats (ECMA-376 Part 1 §18.8.31): a format code has up to four
// sections separated by ";" (positive; negative; zero; text) or sections
// chosen by conditions ("[>=100]"). A section is a sequence of tokens: digit
// placeholders (0 # ?), the decimal point, thousands separators and scaling
// commas, percent, exponents, fractions, date and time codes, literal text,
// gaps the width of a character ("_x"), a character repeated to fill the
// cell ("*x"), and the text placeholder "@".

// piece is part of a formatted value.
type piece struct {
	kind byte // pieceText, pieceGap (the width of s, drawn as nothing), pieceFill (s repeated to fill the cell)
	s    string
}

const (
	pieceText = iota
	pieceGap
	pieceFill
)

// formatted is a value formatted for display.
type formatted struct {
	pieces []piece
	color  *rgb // from [Red], [Color10] …; nil for the cell's font color
	// general is set when the value was formatted by "General", which a
	// narrow column shortens (see formatGeneral).
	general bool
	// number is set for numbers and dates, which do not overflow into the
	// next cells: they show "#" when they do not fit.
	number bool
}

func (f formatted) text() string {
	var b strings.Builder
	for _, p := range f.pieces {
		if p.kind == pieceText {
			b.WriteString(p.s)
		}
	}
	return b.String()
}

type tokKind byte

const (
	tLit     tokKind = iota // literal text
	tGap                    // _x
	tFill                   // *x
	tText                   // @
	tDigit                  // 0 # ?
	tPoint                  // .
	tComma                  // ,
	tPercent                // %
	tExp                    // E+ E- e+ e-
	tSlash                  // / of a fraction
	tDate                   // y m d h s AM/PM g e aaa b …
	tElapsed                // [h] [m] [s]
	tGeneral                // General
	tSubsec                 // .0 .00 .000 after seconds
)

type token struct {
	kind tokKind
	s    string // literal, placeholder character, date code ("yyyy", "mm", "AM/PM" …)
}

type condition struct {
	op  string // < > = <= >= <>
	val float64
}

type section struct {
	toks   []token
	color  *rgb
	cond   *condition
	date   bool
	text   bool // has @
	locale string
	dbnum  int
}

// numFormat is a parsed format code.
type numFormat struct {
	code     string
	sections []*section
}

var fmtColors = map[string]rgb{
	"black": {0, 0, 0}, "blue": {0, 0, 1}, "cyan": {0, 1, 1}, "green": {0, 1, 0},
	"magenta": {1, 0, 1}, "red": {1, 0, 0}, "white": {1, 1, 1}, "yellow": {1, 1, 0},
}

// parseNumFormat parses a format code; indexed colors ([Color10]) come from palette.
func parseNumFormat(code string, palette []rgb) *numFormat {
	nf := &numFormat{code: code}
	for _, s := range splitSections(code) {
		nf.sections = append(nf.sections, parseSection(s, palette))
		if len(nf.sections) == 4 {
			break
		}
	}
	if len(nf.sections) == 0 {
		nf.sections = []*section{{toks: []token{{kind: tGeneral}}}}
	}
	return nf
}

// splitSections splits a code at the semicolons outside quotes and brackets.
func splitSections(code string) []string {
	var out []string
	start, quoted, bracket := 0, false, false
	for i := 0; i < len(code); i++ {
		switch c := code[i]; {
		case quoted:
			if c == '"' {
				quoted = false
			}
		case c == '"':
			quoted = true
		case c == '\\' || c == '_' || c == '*':
			if i+1 < len(code) {
				_, n := utf8.DecodeRuneInString(code[i+1:])
				i += n
			}
		case c == '[':
			bracket = true
		case c == ']':
			bracket = false
		case c == ';' && !bracket:
			out = append(out, code[start:i])
			start = i + 1
		}
	}
	return append(out, code[start:])
}

func parseSection(s string, palette []rgb) *section {
	sec := &section{}
	lower := strings.ToLower(s)
	add := func(k tokKind, v string) {
		// merge adjacent literals
		if k == tLit && len(sec.toks) > 0 && sec.toks[len(sec.toks)-1].kind == tLit {
			sec.toks[len(sec.toks)-1].s += v
			return
		}
		sec.toks = append(sec.toks, token{k, v})
	}
	for i := 0; i < len(s); {
		c := s[i]
		r, n := utf8.DecodeRuneInString(s[i:])
		switch {
		case c == '"':
			j := strings.IndexByte(s[i+1:], '"')
			if j < 0 {
				add(tLit, s[i+1:])
				i = len(s)
				continue
			}
			add(tLit, s[i+1:i+1+j])
			i += j + 2
			continue
		case c == '\\' || c == '_' || c == '*':
			if i+1 >= len(s) {
				i++
				continue
			}
			nr, nn := utf8.DecodeRuneInString(s[i+1:])
			switch c {
			case '\\':
				add(tLit, string(nr))
			case '_':
				add(tGap, string(nr))
			case '*':
				add(tFill, string(nr))
			}
			i += 1 + nn
			continue
		case c == '[':
			j := strings.IndexByte(s[i:], ']')
			if j < 0 {
				i = len(s)
				continue
			}
			sec.bracket(s[i+1:i+j], palette, add)
			i += j + 1
			continue
		case c == '@':
			add(tText, "@")
			sec.text = true
		case c == '0' || c == '#' || c == '?':
			add(tDigit, string(c))
		case c == '.':
			// fractions of a second follow a seconds code
			if sec.lastDate() == 's' && i+1 < len(s) && s[i+1] == '0' {
				j := i + 1
				for j < len(s) && s[j] == '0' {
					j++
				}
				add(tSubsec, s[i+1:j])
				i = j
				continue
			}
			add(tPoint, ".")
		case c == ',':
			add(tComma, ",")
		case c == '%':
			add(tPercent, "%")
		case (c == 'E' || c == 'e') && i+1 < len(s) && (s[i+1] == '+' || s[i+1] == '-') && !sec.hasDate():
			add(tExp, s[i:i+2])
			i += 2
			continue
		case c == '/' && sec.prevDigit():
			add(tSlash, "/")
		case strings.HasPrefix(lower[i:], "general"):
			add(tGeneral, "General")
			i += 7
			continue
		case strings.HasPrefix(lower[i:], "am/pm"):
			add(tDate, "AM/PM")
			sec.date = true
			i += 5
			continue
		case strings.HasPrefix(lower[i:], "a/p"):
			add(tDate, "A/P")
			sec.date = true
			i += 3
			continue
		case strings.ContainsRune("ymdhsgebYMDHSGEB", r) || (c == 'a' || c == 'A') && strings.HasPrefix(lower[i:], "aaa"):
			j := i + 1
			lc := lower[i]
			for j < len(s) && lower[j] == lc {
				j++
			}
			add(tDate, lower[i:j])
			sec.date = true
			i = j
			continue
		default:
			add(tLit, string(r))
		}
		i += n
	}
	return sec
}

func (sec *section) lastDate() byte {
	for i := len(sec.toks) - 1; i >= 0; i-- {
		t := sec.toks[i]
		switch t.kind {
		case tDate:
			return t.s[0]
		case tElapsed:
			return t.s[0]
		case tLit, tGap:
			continue
		}
		return 0
	}
	return 0
}

func (sec *section) hasDate() bool { return sec.date }

// prevDigit reports whether a digit placeholder comes before the current
// position with only digits or literals between (the "/" of a fraction).
func (sec *section) prevDigit() bool {
	for i := len(sec.toks) - 1; i >= 0; i-- {
		switch sec.toks[i].kind {
		case tDigit:
			return true
		case tLit:
			continue
		default:
			return false
		}
	}
	return false
}

// bracket interprets a [...] code: a color, a condition, a locale or
// currency, an elapsed time unit or a numeral system.
func (sec *section) bracket(b string, palette []rgb, add func(tokKind, string)) {
	l := strings.ToLower(strings.TrimSpace(b))
	switch {
	case l == "":
	case l[0] == '$':
		// [$€-407]: a currency symbol and a locale id
		sym, loc, _ := strings.Cut(b[1:], "-")
		if sym != "" {
			add(tLit, sym)
		}
		sec.locale = strings.ToLower(loc)
	case strings.Trim(l, "h") == "" || strings.Trim(l, "m") == "" || strings.Trim(l, "s") == "":
		add(tElapsed, l)
		sec.date = true
	case l[0] == '<' || l[0] == '>' || l[0] == '=':
		op := l[:1]
		rest := l[1:]
		if len(rest) > 0 && (rest[0] == '=' || rest[0] == '>') {
			op += rest[:1]
			rest = rest[1:]
		}
		if v, err := strconv.ParseFloat(strings.TrimSpace(rest), 64); err == nil {
			sec.cond = &condition{op, v}
		}
	case strings.HasPrefix(l, "color"):
		if n, err := strconv.Atoi(strings.TrimSpace(l[5:])); err == nil && n >= 1 && n <= 56 && n+7 < len(palette) {
			c := palette[n+7]
			sec.color = &c
		}
	case strings.HasPrefix(l, "dbnum"):
		sec.dbnum, _ = strconv.Atoi(l[5:])
	default:
		if c, ok := fmtColors[l]; ok {
			sec.color = &c
		}
	}
}

func (c *condition) match(v float64) bool {
	switch c.op {
	case "<":
		return v < c.val
	case ">":
		return v > c.val
	case "=":
		return v == c.val
	case "<=":
		return v <= c.val
	case ">=":
		return v >= c.val
	case "<>":
		return v != c.val
	}
	return false
}

// isDate reports whether the format shows numbers as dates or times.
func (nf *numFormat) isDate() bool {
	return len(nf.sections) > 0 && nf.sections[0].date
}

// pick returns the section for a number and whether the number's sign is
// shown by the section itself (the value is then formatted without it).
func (nf *numFormat) pick(v float64) (*section, bool) {
	secs := nf.sections
	// a trailing text section does not take numbers
	if n := len(secs); n > 1 && secs[n-1].text && !secs[n-1].hasDigits() {
		secs = secs[:n-1]
	}
	if secs[0].cond != nil || len(secs) > 1 && secs[1].cond != nil {
		// conditional sections: the first two are tried in order, the
		// third (or the second without a condition) takes the rest
		for i, s := range secs {
			if i == 2 || s.cond == nil || s.cond.match(v) {
				// Sections chosen by a condition show the value without its sign
				// when a later section exists for the other values.
				return s, i > 0 || len(secs) > 1 && v < 0 && s.cond != nil && s.cond.val <= 0
			}
		}
		return secs[len(secs)-1], true
	}
	switch {
	case v > 0 || len(secs) == 1:
		return secs[0], false
	case v < 0:
		return secs[1], true
	case len(secs) >= 3:
		return secs[2], false
	}
	return secs[0], false
}

func (sec *section) hasDigits() bool {
	for _, t := range sec.toks {
		switch t.kind {
		case tDigit, tGeneral, tDate, tElapsed:
			return true
		}
	}
	return false
}

// formatNumber formats a number.
func (nf *numFormat) formatNumber(v float64, date1904 bool) formatted {
	sec, unsigned := nf.pick(v)
	out := formatted{color: sec.color, number: true}
	if sec.date {
		s, ok := formatDate(sec, v, date1904)
		if !ok {
			out.pieces = []piece{{pieceText, "#"}}
			out.general = false
			return out
		}
		out.pieces = s
		return out
	}
	if unsigned {
		v = math.Abs(v)
	}
	if len(sec.toks) == 0 {
		return out // an empty section hides the value
	}
	if sec.onlyText() {
		// a text-only section shows numbers in General
		out.pieces = []piece{{pieceText, formatGeneral(v, 11)}}
		out.general = true
		return out
	}
	out.pieces, out.general = formatNumeric(sec, v)
	if sec.dbnum > 0 {
		for i := range out.pieces {
			if out.pieces[i].kind == pieceText {
				out.pieces[i].s = dbNum(out.pieces[i].s, sec.dbnum)
			}
		}
	}
	return out
}

func (sec *section) onlyText() bool {
	for _, t := range sec.toks {
		switch t.kind {
		case tDigit, tGeneral, tPoint, tPercent, tExp, tDate, tElapsed:
			return false
		}
	}
	return sec.text
}

// formatText formats a string: with the text section (the fourth, or a
// single section with @), or as is.
func (nf *numFormat) formatText(s string) formatted {
	var sec *section
	for i, sc := range nf.sections {
		if sc.text || i == 3 {
			sec = sc
		}
	}
	if sec == nil {
		return formatted{pieces: []piece{{pieceText, s}}}
	}
	out := formatted{color: sec.color}
	for _, t := range sec.toks {
		switch t.kind {
		case tText:
			out.pieces = append(out.pieces, piece{pieceText, s})
		case tLit:
			out.pieces = append(out.pieces, piece{pieceText, t.s})
		case tGap:
			out.pieces = append(out.pieces, piece{pieceGap, t.s})
		case tFill:
			out.pieces = append(out.pieces, piece{pieceFill, t.s})
		}
	}
	return out
}

// formatGeneral formats a number the way General does in a cell that holds
// at most maxChars characters (11 in a standard column; not counting the
// sign): integers as they are, other numbers rounded to the digits that
// fit, and very large or small ones in scientific notation. It returns ""
// when nothing fits.
func formatGeneral(v float64, maxChars int) string {
	if v == 0 || math.IsNaN(v) {
		return "0"
	}
	if math.IsInf(v, 0) {
		return "#NUM!"
	}
	sign := ""
	if v < 0 {
		sign = "-"
		v = -v
	}
	v = round15(v)
	exp := int(math.Floor(math.Log10(v)))
	if exp < 11 && exp >= -9 {
		intDigits := max(exp+1, 1)
		if intDigits <= maxChars {
			decimals := maxChars - intDigits - 1
			if decimals < 0 {
				decimals = 0
			}
			s := strconv.FormatFloat(v, 'f', min(decimals, 30), 64)
			if strings.Contains(s, ".") {
				s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
			}
			if s != "0" && len(s) <= maxChars {
				return sign + s
			}
		}
	}
	// scientific: mantissa digits that fit next to "E+nn"
	es := "E+"
	if exp < 0 {
		es = "E-"
	}
	ex := exp
	if ex < 0 {
		ex = -ex
	}
	expStr := es + pad2(ex)
	room := maxChars - len(expStr)
	if room < 1 {
		return ""
	}
	m := v / math.Pow(10, float64(exp))
	decimals := room - 2
	if decimals < 0 {
		decimals = 0
	}
	ms := strconv.FormatFloat(m, 'f', decimals, 64)
	if strings.HasPrefix(ms, "10") { // rounded up to 10
		exp++
		return formatGeneral(sign2(sign)*math.Pow(10, float64(exp)), maxChars)
	}
	if strings.Contains(ms, ".") {
		ms = strings.TrimRight(strings.TrimRight(ms, "0"), ".")
	}
	return sign + ms + expStr
}

func sign2(s string) float64 {
	if s == "-" {
		return -1
	}
	return 1
}

func pad2(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}

// round15 rounds to the 15 significant digits Excel keeps.
func round15(v float64) float64 {
	if v == 0 || math.IsInf(v, 0) || math.IsNaN(v) {
		return v
	}
	f, err := strconv.ParseFloat(strconv.FormatFloat(v, 'g', 15, 64), 64)
	if err != nil {
		return v
	}
	return f
}

// numeric parts of a number section
type numLayout struct {
	intPh, fracPh []byte // placeholder characters of the integer and fraction parts
	point         bool
	thousands     bool
	scale         int // scaling commas (each divides by 1000)
	percent       int
	exp           string
	expPh         []byte
	numPh, denPh  []byte // fraction placeholders
	denFixed      int    // fraction with a fixed denominator (/8)
	fraction      bool
	general       bool
}

// analyze finds the numeric structure of a section.
func analyze(toks []token) numLayout {
	var l numLayout
	part := 0 // 0 integer, 1 fraction digits, 2 exponent, 3 numerator/denominator handled separately
	lastDigit := -1
	for i, t := range toks {
		switch t.kind {
		case tGeneral:
			l.general = true
		case tSlash:
			l.fraction = true
		case tDigit:
			lastDigit = i
		}
	}
	if l.fraction {
		// integer part, numerator and denominator: "# ?/?" or "?/8"
		slash := -1
		for i, t := range toks {
			if t.kind == tSlash {
				slash = i
			}
		}
		// denominator: digits (or a literal number) after the slash
		for i := slash + 1; i < len(toks); i++ {
			t := toks[i]
			if t.kind == tDigit {
				l.denPh = append(l.denPh, t.s[0])
				continue
			}
			if t.kind == tLit && len(l.denPh) == 0 {
				if n, err := strconv.Atoi(strings.TrimSpace(t.s)); err == nil && n > 0 {
					l.denFixed = n
				}
			}
			break
		}
		// numerator: the digit group right before the slash
		j := slash - 1
		for j >= 0 && toks[j].kind == tLit {
			j--
		}
		for j >= 0 && toks[j].kind == tDigit {
			l.numPh = append([]byte{toks[j].s[0]}, l.numPh...)
			j--
		}
		// integer part: digits before the numerator group
		for i := 0; i <= j; i++ {
			if toks[i].kind == tDigit {
				l.intPh = append(l.intPh, toks[i].s[0])
			}
			if toks[i].kind == tComma {
				l.thousands = true
			}
		}
		for _, t := range toks {
			if t.kind == tPercent {
				l.percent++
			}
		}
		return l
	}
	for i, t := range toks {
		switch t.kind {
		case tDigit:
			switch part {
			case 0:
				l.intPh = append(l.intPh, t.s[0])
			case 1:
				l.fracPh = append(l.fracPh, t.s[0])
			case 2:
				l.expPh = append(l.expPh, t.s[0])
			}
		case tPoint:
			if part == 0 {
				l.point = true
				part = 1
			}
		case tExp:
			if part < 2 {
				l.exp = t.s
				part = 2
			}
		case tPercent:
			l.percent++
		case tComma:
			if part != 0 && part != 1 {
				continue
			}
			// a comma between digit placeholders groups thousands; commas
			// right after the last digit of the number scale it
			next := nextKind(toks, i)
			if part == 0 && next == tDigit && len(l.intPh) > 0 {
				l.thousands = true
			} else if i > lastDigit || next != tDigit {
				if prevDigitOrComma(toks, i) {
					l.scale++
				}
			}
		}
	}
	return l
}

func nextKind(toks []token, i int) tokKind {
	for j := i + 1; j < len(toks); j++ {
		if toks[j].kind != tComma {
			return toks[j].kind
		}
	}
	return tLit
}

func prevDigitOrComma(toks []token, i int) bool {
	for j := i - 1; j >= 0; j-- {
		switch toks[j].kind {
		case tComma:
			continue
		case tDigit:
			return true
		}
		return false
	}
	return false
}

// formatNumeric formats a number (without its sign when the section shows
// it) with a number section.
func formatNumeric(sec *section, v float64) ([]piece, bool) {
	l := analyze(sec.toks)
	neg := v < 0
	v = math.Abs(v)
	for i := 0; i < l.percent; i++ {
		v *= 100
	}
	for i := 0; i < l.scale; i++ {
		v /= 1000
	}
	var intStr, fracStr, expStr, numStr, denStr string
	expNeg := false
	general := false
	switch {
	case l.general && len(l.intPh) == 0:
		intStr = formatGeneral(v, 11)
		general = true
	case l.fraction:
		intStr, numStr, denStr = fractionParts(v, l)
	case l.exp != "":
		intStr, fracStr, expStr, expNeg = sciParts(v, l)
	default:
		rv := roundHalfAway(round15(v), len(l.fracPh))
		s := strconv.FormatFloat(rv, 'f', len(l.fracPh), 64)
		intStr, fracStr, _ = strings.Cut(s, ".")
		if intStr == "0" {
			intStr = ""
		}
	}
	// Negative numbers that round to zero lose their sign.
	if neg && !isZeroDigits(intStr+fracStr+numStr) || neg && general {
		intStr = "-" + intStr
		if general {
			intStr = strings.Replace(intStr, "--", "-", 1)
		}
	}
	if l.thousands {
		intStr = groupThousands(intStr)
	}
	var out []piece
	text := func(s string) {
		if s == "" {
			return
		}
		if n := len(out); n > 0 && out[n-1].kind == pieceText {
			out[n-1].s += s
			return
		}
		out = append(out, piece{pieceText, s})
	}
	intDigits := fillPlaceholders(intStr, l.intPh, true)
	fracDigits := fillFraction(fracStr, l.fracPh)
	expDigits := fillPlaceholders(expStr, l.expPh, true)
	numDigits := fillPlaceholders(numStr, l.numPh, true)
	denDigits := fillFraction(denStr, l.denPh)
	if l.fraction && numStr == "" && denStr == "" {
		// a whole number: the fraction's place stays blank
		numDigits = fillPlaceholders("", bytesOf('?', len(l.numPh)), true)
		denDigits = fillPlaceholders("", bytesOf('?', len(l.denPh)), true)
	}
	part := 0
	intIdx, fracIdx, expIdx, numIdx, denIdx := 0, 0, 0, 0, 0
	slashSeen := false
	numStart := len(sec.toks)
	if l.fraction {
		// the numerator's first placeholder
		cnt := 0
		for i := len(sec.toks) - 1; i >= 0; i-- {
			if sec.toks[i].kind == tSlash {
				for j := i - 1; j >= 0; j-- {
					if sec.toks[j].kind == tDigit {
						cnt++
						if cnt == len(l.numPh) {
							numStart = j
							break
						}
					} else if sec.toks[j].kind != tLit {
						break
					}
				}
				break
			}
		}
	}
	for i, t := range sec.toks {
		switch t.kind {
		case tLit:
			// the literal between a whole number and a fraction disappears
			// with the whole number
			if l.fraction && i < numStart && intStr == "" && len(l.intPh) > 0 && !hasZero(l.intPh) {
				if onlySpaces(t.s) {
					out = append(out, piece{pieceGap, t.s})
					continue
				}
			}
			text(t.s)
		case tGap:
			out = append(out, piece{pieceGap, t.s})
		case tFill:
			out = append(out, piece{pieceFill, t.s})
		case tGeneral:
			text(intStr)
		case tText:
		case tPercent:
			text("%")
		case tPoint:
			if part == 0 {
				text(".")
				part = 1
			} else {
				text(".")
			}
		case tExp:
			if expNeg {
				text(t.s[:1] + "-")
			} else if t.s[1] == '+' {
				text(t.s[:1] + "+")
			} else {
				text(t.s[:1])
			}
			part = 2
		case tSlash:
			if l.fraction && numStr == "" && denStr == "" {
				out = append(out, piece{pieceGap, "/"})
			} else {
				text("/")
			}
			slashSeen = true
		case tDigit:
			switch {
			case l.fraction && slashSeen:
				if denIdx < len(denDigits) {
					out = appendDigits(out, denDigits[denIdx])
				}
				denIdx++
			case l.fraction && i >= numStart:
				if numIdx < len(numDigits) {
					out = appendDigits(out, numDigits[numIdx])
				}
				numIdx++
			case part == 0:
				if intIdx < len(intDigits) {
					out = appendDigits(out, intDigits[intIdx])
				}
				intIdx++
			case part == 1:
				if fracIdx < len(fracDigits) {
					out = appendDigits(out, fracDigits[fracIdx])
				}
				fracIdx++
			default:
				if expIdx < len(expDigits) {
					out = appendDigits(out, expDigits[expIdx])
				}
				expIdx++
			}
		}
	}
	if len(l.intPh) == 0 && !l.general && intStr != "" && intStr != "-" && !l.fraction {
		// no integer placeholders: the integer digits still show before the point
		for i, p := range out {
			if p.kind == pieceText && strings.HasPrefix(p.s, ".") {
				out = append(out[:i], append([]piece{{pieceText, intStr}}, out[i:]...)...)
				break
			}
		}
	}
	return out, general
}

func bytesOf(c byte, n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = c
	}
	return b
}

func hasZero(ph []byte) bool {
	for _, c := range ph {
		if c == '0' {
			return true
		}
	}
	return false
}

func onlySpaces(s string) bool { return strings.Trim(s, " ") == "" }

func isZeroDigits(s string) bool {
	for _, c := range s {
		if c >= '1' && c <= '9' {
			return false
		}
	}
	return true
}

// appendDigits adds the output of one placeholder: digits, or for "?"
// a gap the width of a digit.
func appendDigits(out []piece, s string) []piece {
	for _, part := range splitGaps(s) {
		if part.kind == pieceText {
			if n := len(out); n > 0 && out[n-1].kind == pieceText {
				out[n-1].s += part.s
				continue
			}
		}
		out = append(out, part)
	}
	return out
}

// splitGaps turns the "\x00" markers of "?" placeholders into gaps.
func splitGaps(s string) []piece {
	var out []piece
	for s != "" {
		i := strings.IndexByte(s, 0)
		if i < 0 {
			out = append(out, piece{pieceText, s})
			break
		}
		if i > 0 {
			out = append(out, piece{pieceText, s[:i]})
		}
		out = append(out, piece{pieceGap, "0"})
		s = s[i+1:]
	}
	return out
}

// fillPlaceholders distributes the digits of an integer (with its sign and
// separators) over placeholders from the right; the first placeholder takes
// the digits left over. "0" pads with zeros, "?" with gaps, "#" with nothing.
func fillPlaceholders(digits string, ph []byte, leading bool) []string {
	out := make([]string, len(ph))
	if len(ph) == 0 {
		return out
	}
	// split the digit string into units: a digit with its grouping comma
	var units []string
	for i := 0; i < len(digits); {
		c := digits[i]
		if c == '-' || c == ',' {
			// keep with the next digit
			j := i + 1
			for j < len(digits) && (digits[j] == ',') {
				j++
			}
			if j < len(digits) {
				units = append(units, digits[i:j+1])
				i = j + 1
				continue
			}
			units = append(units, digits[i:])
			break
		}
		units = append(units, digits[i:i+1])
		i++
	}
	j := len(units) - 1
	for i := len(ph) - 1; i >= 0; i-- {
		if j >= 0 {
			if i == 0 {
				out[i] = strings.Join(units[:j+1], "")
			} else {
				out[i] = units[j]
			}
			j--
			continue
		}
		switch ph[i] {
		case '0':
			out[i] = "0"
		case '?':
			out[i] = "\x00"
		}
	}
	return out
}

// fillFraction distributes decimal digits over placeholders from the left:
// trailing zeros show for "0", become gaps for "?" and disappear for "#".
func fillFraction(digits string, ph []byte) []string {
	out := make([]string, len(ph))
	last := len(digits) - 1
	for last >= 0 && digits[last] == '0' {
		last--
	}
	for i := range ph {
		d := byte('0')
		if i < len(digits) {
			d = digits[i]
		}
		if i <= last {
			out[i] = string(d)
			continue
		}
		switch ph[i] {
		case '0':
			out[i] = string(d)
		case '?':
			out[i] = "\x00"
		}
	}
	return out
}

func groupThousands(s string) string {
	neg := strings.HasPrefix(s, "-")
	if neg {
		s = s[1:]
	}
	var b strings.Builder
	for i := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteByte(s[i])
	}
	if neg {
		return "-" + b.String()
	}
	return b.String()
}

func roundHalfAway(v float64, decimals int) float64 {
	p := math.Pow(10, float64(decimals))
	r := math.Round(v*p) / p
	if math.IsInf(r, 0) || math.IsNaN(r) {
		return v
	}
	return r
}

// sciParts formats the mantissa and exponent of scientific notation;
// "##0.0E+0" (integer placeholders beyond one) uses exponents that are
// multiples of their number.
func sciParts(v float64, l numLayout) (intStr, fracStr, expStr string, expNeg bool) {
	if v == 0 {
		return "0", strings.Repeat("0", len(l.fracPh)), "0", false
	}
	exp := int(math.Floor(math.Log10(v)))
	step := 1
	if n := len(l.intPh); n > 1 && hasHash(l.intPh) {
		step = n
	}
	e := exp
	if step > 1 {
		e = int(math.Floor(float64(exp)/float64(step))) * step
	} else if n := len(l.intPh); n > 1 {
		e = exp - (n - 1)
	}
	m := v / math.Pow(10, float64(e))
	m = roundHalfAway(m, len(l.fracPh))
	if step == 1 && len(l.intPh) <= 1 && m >= 10 {
		m /= 10
		e++
	}
	s := strconv.FormatFloat(m, 'f', len(l.fracPh), 64)
	intStr, fracStr, _ = strings.Cut(s, ".")
	if e < 0 {
		return intStr, fracStr, strconv.Itoa(-e), true
	}
	return intStr, fracStr, strconv.Itoa(e), false
}

func hasHash(ph []byte) bool {
	for _, c := range ph {
		if c == '#' {
			return true
		}
	}
	return false
}

// fractionParts finds the closest fraction with as many denominator digits
// as the format has (or its fixed denominator).
func fractionParts(v float64, l numLayout) (intStr, numStr, denStr string) {
	whole := 0.0
	frac := v
	if len(l.intPh) > 0 {
		whole = math.Floor(v)
		frac = v - whole
	}
	var num, den int
	if l.denFixed > 0 {
		den = l.denFixed
		num = int(math.Round(frac * float64(den)))
	} else {
		maxDen := int(math.Pow(10, float64(max(len(l.denPh), 1)))) - 1
		num, den = bestFraction(frac, maxDen)
	}
	if den > 0 && num == den && len(l.intPh) > 0 {
		whole++
		num = 0
	}
	if whole != 0 {
		intStr = strconv.FormatFloat(whole, 'f', 0, 64)
	}
	if num == 0 {
		if len(l.intPh) > 0 && whole != 0 {
			// "1" rather than "1 0/1": the fraction part becomes blank
			return intStr, "", ""
		}
		if len(l.intPh) > 0 {
			return "0", "", ""
		}
		return "", "0", strconv.Itoa(max(den, 1))
	}
	return intStr, strconv.Itoa(num), strconv.Itoa(den)
}

func bestFraction(x float64, maxDen int) (int, int) {
	bestN, bestD, bestErr := 0, 1, math.Abs(x)
	for d := 1; d <= maxDen; d++ {
		n := int(math.Round(x * float64(d)))
		if e := math.Abs(x - float64(n)/float64(d)); e < bestErr-1e-12 {
			bestN, bestD, bestErr = n, d, e
			if e == 0 {
				break
			}
		}
	}
	return bestN, bestD
}

// Dates are serial numbers: days since 1900-01-00 (with the fictitious
// 1900-02-29 of Lotus 1-2-3 that Excel keeps) or, in the 1904 date system,
// since 1904-01-01. The fraction is the time of day.

// serialTime converts a serial date to a time (UTC); false for dates Excel
// cannot show (negative, after 9999-12-31).
func serialTime(v float64, date1904 bool) (t time.Time, leapBug bool, ok bool) {
	if v < 0 || v > 2958465.99999999 {
		return time.Time{}, false, false
	}
	days := math.Floor(v)
	ms := math.Round((v - days) * 86400000)
	if ms >= 86400000 {
		days++
		ms -= 86400000
	}
	var base time.Time
	if date1904 {
		base = time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
	} else {
		switch {
		case days == 60:
			// 1900-02-29
			return time.Date(1900, 2, 28, 0, 0, 0, 0, time.UTC).Add(time.Duration(ms) * time.Millisecond), true, true
		case days < 60:
			base = time.Date(1899, 12, 31, 0, 0, 0, 0, time.UTC)
		default:
			base = time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
		}
	}
	t = base.AddDate(0, 0, int(days)).Add(time.Duration(ms) * time.Millisecond)
	return t, false, true
}

var (
	monthNames = [...]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	dayNames   = [...]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	jaDayNames = [...]string{"日", "月", "火", "水", "木", "金", "土"}
)

type era struct {
	start          time.Time
	name, abbr, ch string // 平成, 平, H
}

var jaEras = []era{
	{time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC), "令和", "令", "R"},
	{time.Date(1989, 1, 8, 0, 0, 0, 0, time.UTC), "平成", "平", "H"},
	{time.Date(1926, 12, 25, 0, 0, 0, 0, time.UTC), "昭和", "昭", "S"},
	{time.Date(1912, 7, 30, 0, 0, 0, 0, time.UTC), "大正", "大", "T"},
	{time.Date(1868, 1, 1, 0, 0, 0, 0, time.UTC), "明治", "明", "M"},
}

func eraOf(t time.Time) (era, int) {
	for _, e := range jaEras {
		if !t.Before(e.start) {
			return e, t.Year() - e.start.Year() + 1
		}
	}
	return era{name: "", abbr: "", ch: ""}, t.Year()
}

// formatDate formats a serial date or time with a date section.
func formatDate(sec *section, v float64, date1904 bool) ([]piece, bool) {
	// round to the smallest unit shown
	sub := 0
	for _, t := range sec.toks {
		if t.kind == tSubsec {
			sub = max(sub, len(t.s))
		}
	}
	unit := 86400.0 * math.Pow(10, float64(sub))
	v = math.Round(v*unit) / unit
	t, leap, ok := serialTime(v, date1904)
	if !ok {
		return nil, false
	}
	hour12 := false
	for _, tk := range sec.toks {
		if tk.kind == tDate && (tk.s == "AM/PM" || tk.s == "A/P") {
			hour12 = true
		}
	}
	// month or minute: "m" after an hour or before a second is minutes
	isMinute := make([]bool, len(sec.toks))
	for i, tk := range sec.toks {
		if tk.kind != tDate || tk.s[0] != 'm' || len(tk.s) > 2 {
			continue
		}
		for j := i - 1; j >= 0; j-- {
			k := sec.toks[j]
			if k.kind == tDate || k.kind == tElapsed {
				isMinute[i] = k.s[0] == 'h' || k.s == "[h]" || strings.HasPrefix(k.s, "[h")
				break
			}
		}
		if !isMinute[i] {
			for j := i + 1; j < len(sec.toks); j++ {
				k := sec.toks[j]
				if k.kind == tDate || k.kind == tElapsed {
					isMinute[i] = k.s[0] == 's' || strings.HasPrefix(k.s, "[s")
					break
				}
			}
		}
	}
	day, month, year := t.Day(), int(t.Month()), t.Year()
	if leap {
		day, month, year = 29, 2, 1900
	}
	weekday := int(t.Weekday())
	if !date1904 && math.Floor(v) < 61 {
		// Excel's 1900 calendar is a day off before March 1900
		weekday = (int(math.Floor(v)) + 6) % 7
	}
	ja := sec.locale == "411" || strings.HasSuffix(sec.locale, "0411") || strings.HasPrefix(sec.locale, "ja")
	var out []piece
	text := func(s string) {
		if n := len(out); n > 0 && out[n-1].kind == pieceText {
			out[n-1].s += s
			return
		}
		out = append(out, piece{pieceText, s})
	}
	totalSec := v * 86400
	for i, tk := range sec.toks {
		switch tk.kind {
		case tLit, tDigit, tPoint, tComma, tPercent, tSlash:
			text(tk.s)
		case tGap:
			out = append(out, piece{pieceGap, tk.s})
		case tFill:
			out = append(out, piece{pieceFill, tk.s})
		case tSubsec:
			frac := t.Nanosecond() / 1e6
			s := pad3(frac)
			text("." + s[:min(len(tk.s), 3)])
		case tElapsed:
			switch tk.s[0] {
			case 'h':
				text(padN(int(math.Floor(totalSec/3600+1e-9)), len(tk.s)))
			case 'm':
				text(padN(int(math.Floor(totalSec/60+1e-9)), len(tk.s)))
			case 's':
				text(padN(int(math.Floor(totalSec+1e-9)), len(tk.s)))
			}
		case tDate:
			s := tk.s
			switch {
			case s == "AM/PM":
				if t.Hour() < 12 {
					text("AM")
				} else {
					text("PM")
				}
			case s == "A/P":
				if t.Hour() < 12 {
					text("A")
				} else {
					text("P")
				}
			case s[0] == 'y':
				if len(s) <= 2 {
					text(pad2(year % 100))
				} else {
					text(strconv.Itoa(year))
				}
			case s[0] == 'b':
				// Buddhist era
				if len(s) <= 2 {
					text(pad2((year + 543) % 100))
				} else {
					text(strconv.Itoa(year + 543))
				}
			case s[0] == 'm' && isMinute[i]:
				text(padN(t.Minute(), len(s)))
			case s[0] == 'm':
				switch len(s) {
				case 1, 2:
					text(padN(month, len(s)))
				case 3:
					if ja {
						text(strconv.Itoa(month) + "月")
					} else {
						text(monthNames[month-1][:3])
					}
				case 4:
					if ja {
						text(strconv.Itoa(month) + "月")
					} else {
						text(monthNames[month-1])
					}
				default:
					text(monthNames[month-1][:1])
				}
			case s[0] == 'd':
				switch len(s) {
				case 1, 2:
					text(padN(day, len(s)))
				case 3:
					if ja {
						text(jaDayNames[weekday])
					} else {
						text(dayNames[weekday][:3])
					}
				default:
					if ja {
						text(jaDayNames[weekday] + "曜日")
					} else {
						text(dayNames[weekday])
					}
				}
			case s[0] == 'a':
				if len(s) >= 4 {
					text(jaDayNames[weekday] + "曜日")
				} else {
					text(jaDayNames[weekday])
				}
			case s[0] == 'h':
				h := t.Hour()
				if hour12 {
					h %= 12
					if h == 0 {
						h = 12
					}
				}
				text(padN(h, len(s)))
			case s[0] == 's':
				text(padN(t.Second(), len(s)))
			case s[0] == 'g':
				e, _ := eraOf(t)
				switch len(s) {
				case 1:
					text(e.ch)
				case 2:
					text(e.abbr)
				default:
					text(e.name)
				}
			case s[0] == 'e':
				if _, y := eraOf(t); true {
					if len(s) >= 2 {
						text(pad2(y))
					} else {
						text(strconv.Itoa(y))
					}
				}
			default:
				text(s)
			}
		}
	}
	return out, true
}

func padN(n, width int) string {
	if width >= 2 {
		return pad2(n)
	}
	return strconv.Itoa(n)
}

func pad3(n int) string {
	s := strconv.Itoa(n)
	for len(s) < 3 {
		s = "0" + s
	}
	return s
}

// dbNum rewrites the digits of formatted text in East Asian numerals:
// [DBNum1] kanji digits (一二三), [DBNum2] formal ones (壱弐参), [DBNum3]
// full-width digits.
func dbNum(s string, kind int) string {
	var digits [10]string
	switch kind {
	case 1:
		digits = [10]string{"〇", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
	case 2:
		digits = [10]string{"〇", "壱", "弐", "参", "四", "伍", "六", "七", "八", "九"}
	case 3:
		digits = [10]string{"０", "１", "２", "３", "４", "５", "６", "７", "８", "９"}
	default:
		return s
	}
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteString(digits[r-'0'])
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// builtinFormat returns the code of a built-in number format (ECMA-376
// Part 1 §18.8.30) for the workbook's locale ("ja" for Japanese, where the
// currency, short dates and IDs 27–58 differ).
func builtinFormat(id int, locale string) string {
	if locale == "ja" {
		if s, ok := builtinJa[id]; ok {
			return s
		}
	}
	if s, ok := builtinFormats[id]; ok {
		return s
	}
	if id >= 27 && id <= 81 {
		if s, ok := builtinJa[id]; ok {
			return s
		}
	}
	return "General"
}

var builtinFormats = map[int]string{
	0: "General", 1: "0", 2: "0.00", 3: "#,##0", 4: "#,##0.00",
	5: `"$"#,##0_);("$"#,##0)`, 6: `"$"#,##0_);[Red]("$"#,##0)`,
	7: `"$"#,##0.00_);("$"#,##0.00)`, 8: `"$"#,##0.00_);[Red]("$"#,##0.00)`,
	9: "0%", 10: "0.00%", 11: "0.00E+00", 12: "# ?/?", 13: "# ??/??",
	14: "m/d/yyyy", 15: "d-mmm-yy", 16: "d-mmm", 17: "mmm-yy",
	18: "h:mm AM/PM", 19: "h:mm:ss AM/PM", 20: "h:mm", 21: "h:mm:ss", 22: "m/d/yyyy h:mm",
	37: "#,##0 ;(#,##0)", 38: "#,##0 ;[Red](#,##0)", 39: "#,##0.00;(#,##0.00)", 40: "#,##0.00;[Red](#,##0.00)",
	41: `_(* #,##0_);_(* \(#,##0\);_(* "-"_);_(@_)`,
	42: `_("$"* #,##0_);_("$"* \(#,##0\);_("$"* "-"_);_(@_)`,
	43: `_(* #,##0.00_);_(* \(#,##0.00\);_(* "-"??_);_(@_)`,
	44: `_("$"* #,##0.00_);_("$"* \(#,##0.00\);_("$"* "-"??_);_(@_)`,
	45: "mm:ss", 46: "[h]:mm:ss", 47: "mmss.0", 48: "##0.0E+0", 49: "@",
}

var builtinJa = map[int]string{
	5: `"¥"#,##0;"¥"\-#,##0`, 6: `"¥"#,##0;[Red]"¥"\-#,##0`,
	7: `"¥"#,##0.00;"¥"\-#,##0.00`, 8: `"¥"#,##0.00;[Red]"¥"\-#,##0.00`,
	14: "yyyy/m/d", 22: "yyyy/m/d h:mm",
	27: `[$-411]ge.m.d`, 28: `[$-411]ggge"年"m"月"d"日"`, 29: `[$-411]ggge"年"m"月"d"日"`,
	30: "m/d/yy", 31: `yyyy"年"m"月"d"日"`, 32: `h"時"mm"分"`, 33: `h"時"mm"分"ss"秒"`,
	34: `yyyy"年"m"月"`, 35: `m"月"d"日"`, 36: `[$-411]ge.m.d`,
	42: `_ "¥"* #,##0_ ;_ "¥"* \-#,##0_ ;_ "¥"* "-"_ ;_ @_ `,
	44: `_ "¥"* #,##0.00_ ;_ "¥"* \-#,##0.00_ ;_ "¥"* "-"??_ ;_ @_ `,
	50: `[$-411]ge.m.d`, 51: `[$-411]ggge"年"m"月"d"日"`, 52: `yyyy"年"m"月"`, 53: `m"月"d"日"`,
	54: `[$-411]ggge"年"m"月"d"日"`, 55: `yyyy"年"m"月"`, 56: `m"月"d"日"`,
	57: `[$-411]ge.m.d`, 58: `[$-411]ggge"年"m"月"d"日"`,
}
