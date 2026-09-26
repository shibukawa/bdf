package csv

// Records are read leniently, as RFC 4180 describes them and as
// spreadsheets write them: fields separated by a delimiter, records ended
// by CRLF, LF or CR, and fields that contain delimiters, quotes or line
// breaks enclosed in quotes, a quote in them doubled (or, in some exports,
// escaped with a backslash). Spaces before an opening quote and after a
// closing one are dropped. A quote that breaks these rules is text.
//
// The dialect of a file (delimiter, quote, escape) is guessed from its
// first 64 KiB: each candidate reads the records there, and the one whose
// records most consistently have the same number of fields, with the
// fewest quotes out of place, wins; between equally consistent ones, the
// one with more fields (a tab-separated file whose values contain commas
// splits into more fields at tabs), then the one with more numbers among
// its fields (semicolons between numbers with decimal commas), then the
// more usual one.

// Dialect is how a file writes its records.
type Dialect struct {
	// Delimiter separates the fields: ',', '\t', ';', '|' or another
	// ASCII character.
	Delimiter rune
	// Quote encloses fields ('"' or '\''), or is 0 when fields are never
	// quoted (quotes are text).
	Quote rune
	// Escape is '\\' when a quote in a quoted field is escaped with a
	// backslash rather than doubled, else 0.
	Escape rune
}

// reader splits text into records.
type reader struct {
	s          string
	pos        int
	delim      byte
	quote      byte
	escape     byte
	bad        int  // quotes out of place
	quoted     int  // fields in quotes
	unfinished bool // the text ended in a quoted field
}

func newReader(s string, d Dialect) *reader {
	return &reader{s: s, delim: byte(d.Delimiter), quote: byte(d.Quote), escape: byte(d.Escape)}
}

// record returns the next record's fields (none for a blank line), or
// false at the end of the text. fields is reused for the result.
func (r *reader) record(fields []string) ([]string, bool) {
	s := r.s
	if r.pos >= len(s) {
		return nil, false
	}
	fields = fields[:0]
	if c := s[r.pos]; c == '\n' || c == '\r' {
		r.pos = skipEOL(s, r.pos)
		return fields, true
	}
	for {
		f, end := r.field()
		fields = append(fields, f)
		if end >= len(s) {
			r.pos = len(s)
			return fields, true
		}
		if s[end] == r.delim {
			r.pos = end + 1
			continue
		}
		r.pos = skipEOL(s, end)
		return fields, true
	}
}

func skipEOL(s string, i int) int {
	if s[i] == '\r' && i+1 < len(s) && s[i+1] == '\n' {
		return i + 2
	}
	return i + 1
}

// field reads a field at r.pos; it returns the field and the position of
// the delimiter or line break after it (or the end of the text).
func (r *reader) field() (string, int) {
	s, i := r.s, r.pos
	if r.quote != 0 {
		j := i
		for j < len(s) && s[j] == ' ' && r.delim != ' ' {
			j++
		}
		if j < len(s) && s[j] == r.quote {
			return r.quotedField(j + 1)
		}
	}
	j := i
	for ; j < len(s); j++ {
		c := s[j]
		if c == r.delim || c == '\n' || c == '\r' {
			break
		}
		if c == r.quote && r.quote != 0 {
			r.bad++
		}
	}
	return s[i:j], j
}

// quotedField reads a quoted field whose text starts at i.
func (r *reader) quotedField(i int) (string, int) {
	s := r.s
	r.quoted++
	var b []byte // the unescaped text once an escape is met
	from := i
	text := func(end int) string {
		if b == nil {
			return s[from:end]
		}
		return string(append(b, s[from:end]...))
	}
	for i < len(s) {
		c := s[i]
		switch {
		case c == r.escape && r.escape != 0 && i+1 < len(s) && (s[i+1] == r.quote || s[i+1] == r.escape):
			b = append(b, s[from:i]...)
			b = append(b, s[i+1])
			i += 2
			from = i
		case c != r.quote:
			i++
		case i+1 < len(s) && s[i+1] == r.quote: // a doubled quote
			b = append(b, s[from:i+1]...)
			i += 2
			from = i
		default:
			// the closing quote, if the field ends after it
			k := i + 1
			for k < len(s) && s[k] == ' ' && r.delim != ' ' {
				k++
			}
			if k >= len(s) || s[k] == r.delim || s[k] == '\n' || s[k] == '\r' {
				return text(i), k
			}
			r.bad++ // a quote inside the field: text
			i++
		}
	}
	r.unfinished = true
	return text(len(s)), len(s)
}

// candidates of dialect guessing, the more usual first
var (
	delimiters = []rune{',', '\t', ';', '|'}
	quotes     = []Dialect{{Quote: '"'}, {Quote: '"', Escape: '\\'}, {}, {Quote: '\''}, {Quote: '\'', Escape: '\\'}}
)

// maxSniffRecords bounds the records dialect guessing reads.
const maxSniffRecords = 2000

// trial is how a dialect reads a sample.
type trial struct {
	d       Dialect
	records int     // records with fields
	fields  int     // the usual number of fields
	score   float64 // consistency less the quotes out of place
	numbers int     // fields that are numbers in the first records
}

// try reads a sample in a dialect. cut tells that the sample ends before
// the file does: its last record may be incomplete and is left out.
func try(sample string, d Dialect, cut bool) trial {
	r := newReader(sample, d)
	counts := map[int]int{}
	n, numbers := 0, 0
	var fields []string
	for n < maxSniffRecords {
		f, ok := r.record(fields)
		if !ok || cut && r.pos >= len(sample) {
			break
		}
		fields = f
		if len(f) > 0 {
			counts[len(f)]++
			n++
		}
		if n <= 100 {
			for _, v := range f {
				if classify(v).numeric() {
					numbers++
				}
			}
		}
	}
	bad := r.bad
	if r.unfinished && !cut {
		bad += n // a quote left open swallowed the rest
	}
	t := trial{d: d, records: n, numbers: numbers}
	if n == 0 {
		return t
	}
	best := 0
	for k, v := range counts {
		if v > best || v == best && k > t.fields {
			t.fields, best = k, v
		}
	}
	t.score = float64(best)/float64(n) - 0.5*float64(bad)/float64(n)
	return t
}

// sniff guesses the dialect of a sample; fixed delimiter or quote (0 for
// none given, noQuote for none) narrow the candidates. The trial of the
// chosen dialect is returned with it.
func sniff(sample string, cut bool, delim, quote rune) trial {
	delims := delimiters
	if delim != 0 {
		delims = []rune{delim}
	}
	qs := quotes
	switch quote {
	case 0:
	case noQuote:
		qs = []Dialect{{}}
	default:
		qs = []Dialect{{Quote: quote}, {Quote: quote, Escape: '\\'}}
	}
	var best trial
	found := false
	for _, dl := range delims {
		// the best quoting for this delimiter
		var bq trial
		tried := false
		for _, q := range qs {
			if q.Quote == dl {
				continue
			}
			q.Delimiter = dl
			if t := try(sample, q, cut); !tried || t.score > bq.score+1e-9 {
				bq, tried = t, true
			}
		}
		if tried && (!found || better(bq, best)) {
			best, found = bq, true
		}
	}
	return best
}

// better reports whether trial a is a better guess than b (a is the less
// usual delimiter).
func better(a, b trial) bool {
	if (a.fields >= 2) != (b.fields >= 2) {
		return a.fields >= 2
	}
	if d := a.score - b.score; d > 0.02 || d < -0.02 {
		return d > 0
	}
	if a.fields != b.fields {
		return a.fields > b.fields
	}
	return a.numbers > b.numbers
}

// noQuote is the Quote option for files whose fields are never quoted.
const noQuote rune = -1

// sampleOf returns the start of a text that dialect guessing reads, and
// whether it is cut short.
func sampleOf(text string) (string, bool) {
	if len(text) <= sniffLen {
		return text, false
	}
	return text[:sniffLen], true
}
