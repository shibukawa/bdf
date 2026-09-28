// Package mml reads MML (Music Macro Language), the text that chip music
// and game music are written in, into a music.Performance that
// converter/internal/music writes as staff notation and plays.
//
// MML writes each track as a string of commands: notes (c d e f g a b,
// with + or # for a sharp and - for a flat, a length such as 8 or 4. and
// n for a note number counted from o0c, so that n48 is o4c), rests (r),
// the octave (o, < and >; o4c is middle C), the default length (l), the
// tempo (t, global whichever track sets it), the volume (v, 0–15, played
// as the velocity of the notes; the notes of v0 are silent and left out as
// rests) and the tone (@). Loops ([…]n, 2 passes by default, with | or :
// to leave on the last pass, and FlMML's /:n … / … :/) are played out,
// tuplets ({cde}4) share their length, & joins a note to the next one of
// the same pitch (to another pitch it is a slur, and the notes stay two),
// ^ lengthens the note before it, and % gives a length in the dialect's
// ticks (FlMML 384 to a whole note, PPMCK 192). Three dialects are read:
//
//   - generic: FlMML (the MML of niconico's Pikokakiko) and the PLAY
//     statements of MSX and N88-BASIC. Tracks are separated by ';', header
//     lines (#TITLE, #ARTIST, #CODING) give the credits, FlMML's macros
//     ($a=cde; then $a) are expanded, and its tones are played with the
//     nearest General MIDI programs (@0 sine 79, @1 saw 81, @2 and @6
//     triangle 82, @3 and @5 pulse 80 with @W setting the duty, @4, @7,
//     @10 and @11 noise 122, others 80; 80 without a tone). Commands are
//     read in either case.
//   - mabinogi: the MML@melody,chord1,chord2; blocks of Mabinogi and
//     ArcheAge. Each block is an instrument, a piano by default.
//   - ppmck: MCK and PPMCK, the MML of NES music. Lines start with the
//     names of their tracks ("A", or "ABC" for the three): A and B are
//     pulse waves (80, with the duty of @0–@3 as controller 70; @0 when
//     none is set), C the triangle (82), D the noise (122), E the samples
//     (not played), F and on the expansion sound (80). ';' starts a
//     comment, header lines (#TITLE, #COMPOSER, #MAKER, #PROGRAMER,
//     #OCTAVE-REV) are read, definitions (@v0 = { … }) are skipped, and K
//     transposes.
//
// Comments are /* … */ and // in all but Mabinogi's. In MCK, Mabinogi and
// MSX BASIC > raises the octave (MCK's #OCTAVE-REV turns it round), but in
// FlMML < does (and #OCTAVE REVERSE turns it round): the generic dialect
// takes the direction that keeps the notes on either side of each < and >
// closer, unless it is set (Options.Octave).
//
// What only shapes the sound is left out: gate time (q; notes last their
// written length), envelopes, LFOs, filters, detune, sweeps, and PPMCK's
// loop point (L): the music plays once.
package mml

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/shibukawa/bdf/converter/internal/music"
	"golang.org/x/text/encoding/japanese"
	xunicode "golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// The dialects of MML (Options.Dialect).
const (
	Generic  = "generic"
	Mabinogi = "mabinogi"
	PPMCK    = "ppmck"
)

// OctaveShift says which of < and > raises the octave.
type OctaveShift int

const (
	// OctaveAuto follows the dialect: > raises in Mabinogi and PPMCK (<
	// after #OCTAVE-REV); the generic dialect takes the direction from the
	// music.
	OctaveAuto OctaveShift = iota
	// OctaveNormal: > raises the octave and < lowers it.
	OctaveNormal
	// OctaveReverse: < raises the octave and > lowers it.
	OctaveReverse
)

// Options control how MML is read.
type Options struct {
	// Dialect is Generic, Mabinogi or PPMCK; "" detects it.
	Dialect string
	// Octave says which of < and > raises the octave; it overrides the
	// file's #OCTAVE-REV.
	Octave OctaveShift
	// Program is the General MIDI program (0–127) of the tracks that set
	// no tone (all of PPMCK's, whose @ sets the duty); nil keeps the
	// dialect's.
	Program *int
	// Warn receives non-fatal problems; when nil they are returned.
	Warn func(msg string)
}

// Budgets that keep hostile files from running away.
const (
	maxNotes  = 1_000_000      // notes of the performance
	maxTracks = 1024           // tracks with notes
	maxTicks  = 10_000 * whole // length of a track: 10,000 whole notes
	maxSteps  = 20_000_000     // commands played, loops played out
	maxDepth  = 64             // nesting of loops and tuplets
	maxRepeat = 10_000         // passes of a loop
	maxText   = 16 << 20       // MML text after the macros are expanded
	maxNumber = 1<<16 - 1      // numbers saturate here
)

// whole is the length of a whole note in ticks.
const whole = 4 * music.PPQ

// Parse reads MML into a performance. The warnings are returned when
// o.Warn is nil. It fails when the MML plays no note.
func Parse(src []byte, o *Options) (*music.Performance, []string, error) {
	if o == nil {
		o = &Options{}
	}
	program := -1
	if o.Program != nil {
		if *o.Program < 0 || *o.Program > 127 {
			return nil, nil, fmt.Errorf("mml: program %d is not 0–127", *o.Program)
		}
		program = *o.Program
	}
	text := decode(src)
	dialect := strings.ToLower(o.Dialect)
	switch dialect {
	case "", "auto":
		dialect = detectDialect(text)
	case Generic, Mabinogi, PPMCK:
	default:
		return nil, nil, fmt.Errorf("mml: unknown dialect %q", o.Dialect)
	}
	w := &warner{}
	var s *source
	switch dialect {
	case Mabinogi:
		s = splitMabinogi(text)
	case PPMCK:
		s = splitPPMCK(text, w)
	default:
		s = splitGeneric(text, w)
	}
	up := 1
	switch {
	case o.Octave == OctaveReverse:
		up = -1
	case o.Octave == OctaveNormal:
	case dialect == PPMCK && s.reverse:
		up = -1
	}
	p := s.play(dialect, up, program)
	if dialect == Generic && o.Octave == OctaveAuto && p.cost > 0 && !p.stop {
		// the direction that moves less across the octave changes
		if q := s.play(dialect, -1, program); q.cost < p.cost {
			p = q
		}
	}
	msgs := append(w.msgs, p.w.msgs...)
	if o.Warn != nil {
		for _, m := range msgs {
			o.Warn(m)
		}
		msgs = nil
	}
	if len(p.perf.Tracks) == 0 {
		return nil, msgs, errors.New("mml: no notes")
	}
	return p.perf, msgs, nil
}

// decode returns the text of an MML file: UTF-8 (a byte order mark is
// dropped), UTF-16 with a byte order mark, or Shift_JIS when it is not
// UTF-8, with its line breaks made \n.
func decode(b []byte) string {
	switch {
	case bytes.HasPrefix(b, []byte{0xEF, 0xBB, 0xBF}), bytes.HasPrefix(b, []byte{0xFF, 0xFE}), bytes.HasPrefix(b, []byte{0xFE, 0xFF}):
		if d, _, err := transform.Bytes(xunicode.BOMOverride(xunicode.UTF8.NewDecoder()), b); err == nil {
			b = d
		}
	case !utf8.Valid(b):
		if d, err := japanese.ShiftJIS.NewDecoder().Bytes(b); err == nil {
			b = d
		}
	}
	s := strings.ToValidUTF8(string(b), "\uFFFD")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.ReplaceAll(s, "\r", "\n")
}

// warner collects warnings, each key once.
type warner struct {
	msgs []string
	seen map[string]bool
}

func (w *warner) once(key, format string, args ...any) {
	if w.seen[key] {
		return
	}
	if w.seen == nil {
		w.seen = map[string]bool{}
	}
	w.seen[key] = true
	w.msgs = append(w.msgs, fmt.Sprintf(format, args...))
}

// source is an MML file split into its tracks.
type source struct {
	title, composer, arranger, lyricist, copyright string
	// reverse is set by #OCTAVE-REV: < raises the octave.
	reverse bool
	tracks  []trackSource
}

// trackSource is the MML of a track.
type trackSource struct {
	name string // "" names it "Track n"
	mml  []byte
	kind chip
}

// chip is the sound channel a PPMCK track plays on.
type chip int

const (
	chipNone chip = iota // not PPMCK: the tone comes from @
	chipPulse
	chipTriangle
	chipNoise
	chipDPCM
	chipExpansion
)

// header reads a header line (#TITLE Song).
func (s *source) header(line string, w *warner) {
	name, val, _ := strings.Cut(line, " ")
	if n, v, ok := strings.Cut(name, "\t"); ok {
		name, val = n, v+" "+val
	}
	val = strings.TrimSpace(val)
	if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
		val = val[1 : len(val)-1]
	}
	set := func(dst *string) {
		if *dst == "" {
			*dst = val
		}
	}
	switch strings.ToUpper(name) {
	case "#TITLE":
		set(&s.title)
	case "#COMPOSER", "#ARTIST":
		set(&s.composer)
	case "#ARRANGER", "#CODING", "#PROGRAMER", "#PROGRAMMER":
		set(&s.arranger)
	case "#LYRICIST":
		set(&s.lyricist)
	case "#MAKER", "#COPYRIGHT":
		set(&s.copyright)
	case "#OCTAVE-REV":
		s.reverse = val != "0"
	case "#INCLUDE":
		w.once("include", "#INCLUDE %s is not read", val)
	}
}

// stripComments removes the comments of MML: /* */ and // (not between
// FlMML's :/ and /:), and in PPMCK ; to the end of the line. Header lines
// (# at the start) are kept as they are.
func stripComments(text string, semicolon bool) string {
	var b strings.Builder
	lineStart := true
	for i := 0; i < len(text); {
		if lineStart {
			j := i
			for j < len(text) && (text[j] == ' ' || text[j] == '\t') {
				j++
			}
			if j < len(text) && text[j] == '#' {
				end := strings.IndexByte(text[j:], '\n')
				if end < 0 {
					end = len(text)
				} else {
					end += j
				}
				b.WriteString(text[i:end])
				i = end
				lineStart = false
				continue
			}
		}
		lineStart = false
		c := text[i]
		switch {
		case c == '\n':
			lineStart = true
			b.WriteByte(c)
			i++
		case strings.HasPrefix(text[i:], "/*"):
			end := strings.Index(text[i+2:], "*/")
			if end < 0 {
				i = len(text)
			} else {
				i += end + 4
			}
			b.WriteByte(' ')
		case c == ';' && semicolon,
			strings.HasPrefix(text[i:], "//") && (i == 0 || text[i-1] != ':') && !strings.HasPrefix(text[i:], "//:"):
			end := strings.IndexByte(text[i:], '\n')
			if end < 0 {
				i = len(text)
			} else {
				i += end
			}
		default:
			b.WriteByte(c)
			i++
		}
	}
	return b.String()
}

// splitGeneric splits FlMML and BASIC MML into its tracks at ';',
// expanding the macros.
func splitGeneric(text string, w *warner) *source {
	s := &source{}
	var body strings.Builder
	braces := 0
	for _, l := range strings.Split(stripComments(text, false), "\n") {
		t := strings.TrimSpace(l)
		if braces > 0 { // the rest of a tone definition (#OPM@0 { … })
			braces += strings.Count(t, "{") - strings.Count(t, "}")
			continue
		}
		if strings.HasPrefix(t, "#") {
			s.header(t, w)
			if name, _, _ := strings.Cut(t, " "); strings.Contains(name, "@") {
				braces = max(strings.Count(t, "{")-strings.Count(t, "}"), 0)
			}
			continue
		}
		body.WriteString(l)
		body.WriteByte('\n')
	}
	m := &macros{defs: map[string]string{}, w: w}
	for _, st := range strings.Split(body.String(), ";") {
		t := strings.TrimSpace(st)
		if t == "" {
			continue
		}
		if name, val, ok := macroDef(t, w); ok {
			if name != "" {
				m.defs[name] = m.expand(val)
			}
			continue
		}
		if x := m.expand(st); strings.TrimSpace(x) != "" {
			s.tracks = append(s.tracks, trackSource{mml: []byte(x)})
		}
	}
	return s
}

// macroDef reads a FlMML macro definition ($name=mml); ok with an empty
// name is a definition that is not read.
func macroDef(t string, w *warner) (name, val string, ok bool) {
	if t[0] != '$' {
		return "", "", false
	}
	i := 1
	for i < len(t) && isIdent(t[i]) {
		i++
	}
	name = t[1:i]
	rest := strings.TrimLeft(t[i:], " \t\n")
	if strings.HasPrefix(rest, "{") { // arguments
		end := strings.IndexByte(rest, '}')
		if end < 0 || !strings.HasPrefix(strings.TrimLeft(rest[end+1:], " \t\n"), "=") {
			return "", "", false
		}
		w.once("macro-args", "macros with arguments ($%s{…}) are not expanded", name)
		return "", "", true
	}
	if name == "" || !strings.HasPrefix(rest, "=") {
		return "", "", false
	}
	return name, rest[1:], true
}

func isIdent(c byte) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

// macros are FlMML's macros. A macro is expanded when it is defined, so
// that it holds only the macros defined before it and cannot recurse.
type macros struct {
	defs  map[string]string
	total int // text expanded so far
	w     *warner
}

// expand replaces the macros in MML ($name, the longest defined name).
func (m *macros) expand(s string) string {
	if !strings.Contains(s, "$") {
		m.total += len(s)
		return s
	}
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] != '$' {
			b.WriteByte(s[i])
			i++
			continue
		}
		j := i + 1
		for j < len(s) && isIdent(s[j]) {
			j++
		}
		name := ""
		for k := min(j, i+1+64); k > i+1; k-- {
			if _, ok := m.defs[s[i+1:k]]; ok {
				name = s[i+1 : k]
				break
			}
		}
		if name == "" {
			m.w.once("macro:"+s[i+1:j], "macro $%s is not defined", s[i+1:j])
			i = j
			if i < len(s) && s[i] == '{' { // its arguments
				if end := strings.IndexByte(s[i:], '}'); end >= 0 {
					i += end + 1
				}
			}
			continue
		}
		val := m.defs[name]
		if m.total+b.Len()+len(val) > maxText {
			m.w.once("macro-size", "the macros expand to more than %d MiB of MML; the rest is left out", maxText>>20)
			break
		}
		b.WriteString(val)
		i += 1 + len(name)
	}
	m.total += b.Len()
	return b.String()
}

// splitMabinogi reads the MML@melody,chord1,chord2; blocks of Mabinogi.
func splitMabinogi(text string) *source {
	s := &source{}
	var blocks [][]string
	for i := mabinogiAt(text, 0); i >= 0; {
		i += 4
		next := mabinogiAt(text, i)
		end := next
		if end < 0 {
			end = len(text)
		}
		if k := strings.IndexByte(text[i:end], ';'); k >= 0 {
			end = i + k
		}
		blocks = append(blocks, strings.Split(text[i:end], ","))
		i = next
	}
	for bi, parts := range blocks {
		for pi, part := range parts {
			name := "Melody"
			if pi > 0 {
				name = fmt.Sprintf("Chord %d", pi)
			}
			if len(blocks) > 1 {
				name += fmt.Sprintf(" (%d)", bi+1)
			}
			s.tracks = append(s.tracks, trackSource{name: name, mml: []byte(part)})
		}
	}
	return s
}

// mabinogiAt returns the index of the next MML@ (in either case) from i
// that starts MML, not an e-mail address, or -1.
func mabinogiAt(text string, i int) int {
	for i < len(text) {
		k := indexFold(text[i:], "MML@")
		if k < 0 {
			return -1
		}
		k += i
		i = k + 4
		end := min(len(text), i+32)
		if n := strings.IndexAny(text[i:end], ",;\n"); n >= 0 {
			end = i + n
		}
		if ok, _, _ := mmlLike(text[i:end]); ok && (k == 0 || !isIdent(text[k-1]) && text[k-1] != '.' && text[k-1] != '-') {
			return k
		}
	}
	return -1
}

// indexFold is strings.Index ignoring the case of ASCII letters.
func indexFold(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		j := 0
		for j < len(sub) && (s[i+j] == sub[j] || isLetter(sub[j]) && s[i+j]|0x20 == sub[j]|0x20) {
			j++
		}
		if j == len(sub) {
			return i
		}
	}
	return -1
}

func isLetter(c byte) bool { return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' }

func isDigit(c byte) bool { return '0' <= c && c <= '9' }

// splitPPMCK splits PPMCK's MML into its tracks by the names that start
// its lines, skipping definitions.
func splitPPMCK(text string, w *warner) *source {
	s := &source{}
	bufs := map[byte]*bytes.Buffer{}
	var prev []byte
	inDef, braces, opened := false, 0, false
	total := 0 // the MML of all tracks: a line can go to 26 of them
	for _, l := range strings.Split(stripComments(text, true), "\n") {
		t := strings.TrimSpace(l)
		if inDef {
			if !opened && t != "" && t[0] != '{' {
				inDef = false // a definition without braces
			} else {
				braces += strings.Count(t, "{") - strings.Count(t, "}")
				opened = opened || strings.Contains(t, "{")
				inDef = !opened || braces > 0
				continue
			}
		}
		switch {
		case t == "":
			continue
		case t[0] == '#':
			s.header(t, w)
			continue
		case t[0] == '@': // @v0 = { 15 14 13 }, @DPCM0 = { "a.dmc", 15 }
			braces = strings.Count(t, "{") - strings.Count(t, "}")
			opened = strings.Contains(t, "{")
			inDef = !opened || braces > 0
			continue
		}
		names, rest := trackNames(l)
		switch {
		case names != "":
			prev = []byte(names)
		case l[0] == ' ' || l[0] == '\t':
			rest = l // a line that continues the tracks before
		default:
			w.once("no-track", "lines that start with no track name are skipped")
			continue
		}
		for _, c := range prev {
			if total += len(rest) + 1; total > maxText {
				w.once("ppmck-size", "the tracks hold more than %d MiB of MML; the rest is left out", maxText>>20)
				break
			}
			b := bufs[c]
			if b == nil {
				b = &bytes.Buffer{}
				bufs[c] = b
			}
			b.WriteString(rest)
			b.WriteByte('\n')
		}
	}
	for c := byte('A'); c <= 'Z'; c++ {
		b := bufs[c]
		if b == nil {
			continue
		}
		kind := chipExpansion
		switch c {
		case 'A', 'B':
			kind = chipPulse
		case 'C':
			kind = chipTriangle
		case 'D':
			kind = chipNoise
		case 'E':
			kind = chipDPCM
		}
		s.tracks = append(s.tracks, trackSource{name: string(c), mml: b.Bytes(), kind: kind})
	}
	return s
}

// trackNames reads the track names that start a PPMCK line ("ABC l8
// cde"): upper case letters followed by a space or a tab.
func trackNames(l string) (names, rest string) {
	i := 0
	for i < len(l) && 'A' <= l[i] && l[i] <= 'Z' {
		i++
	}
	if i == 0 || i > 26 || i < len(l) && l[i] != ' ' && l[i] != '\t' {
		return "", ""
	}
	return l[:i], l[i:]
}

// player plays the tracks of a source into a performance.
type player struct {
	dialect string
	up      int // the octave change of >
	program int // the program of tracks without a tone, -1 the dialect's
	w       warner
	perf    *music.Performance
	tempos  []music.Tempo
	notes   int
	steps   int
	stop    bool
	// cost adds up the intervals across octave changes (< >): the
	// direction that makes them smaller is the right one.
	cost int
}

// channels are the MIDI channels given to the tracks in turn: 9 is
// percussion.
var channels = [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 10, 11, 12, 13, 14, 15}

func (s *source) play(dialect string, up, program int) *player {
	p := &player{dialect: dialect, up: up, program: program}
	perf := &music.Performance{Title: s.title, Composer: s.composer, Arranger: s.arranger, Lyricist: s.lyricist, Copyright: s.copyright}
	p.perf = perf
	for i := range s.tracks {
		if len(perf.Tracks) == maxTracks {
			p.w.once("tracks", "the MML has more than %d tracks; the rest are left out", maxTracks)
			break
		}
		ts := &s.tracks[i]
		v := p.voice(ts)
		v.run(ts.mml, 0)
		if ts.kind == chipDPCM {
			if len(v.tr.Notes) > 0 {
				p.w.once("dpcm", "the samples of track %s (DPCM) are not played", ts.name)
			}
			continue
		}
		if len(v.tr.Notes) > 0 {
			if v.tr.Name == "" {
				v.tr.Name = fmt.Sprintf("Track %d", len(perf.Tracks)+1)
			}
			v.tr.Channel = channels[len(perf.Tracks)%len(channels)]
			perf.Tracks = append(perf.Tracks, v.tr)
			last := v.tr.Notes[len(v.tr.Notes)-1]
			perf.End = max(perf.End, v.tick, last.Tick+last.Dur)
		}
		if p.stop {
			break
		}
	}
	perf.Tempo = tempoMap(p.tempos)
	return p
}

// tempoMap orders the tempo changes of all tracks: at one tick the last
// track's wins, and a tempo that does not change is dropped. The music
// starts at 120 when no tempo is set at its start.
func tempoMap(ts []music.Tempo) []music.Tempo {
	slices.SortStableFunc(ts, func(a, b music.Tempo) int { return cmp.Compare(a.Tick, b.Tick) })
	out := []music.Tempo{{Tick: 0, BPM: 120}}
	for _, t := range ts {
		if last := &out[len(out)-1]; last.Tick == t.Tick {
			last.BPM = t.BPM
		} else {
			out = append(out, t)
		}
	}
	res := out[:1]
	for _, t := range out[1:] {
		if t.BPM != res[len(res)-1].BPM {
			res = append(res, t)
		}
	}
	return res
}

// voice is the state of a track being played.
type voice struct {
	p    *player
	tr   *music.Track
	kind chip
	fold bool // commands in either case (not PPMCK)

	tick   int
	oct    int
	len    int // the default length (l)
	vol    int // v
	vel    int // the velocity of the notes; 0 is silent
	shift  int // transposition in semitones
	prog   int // the current program
	ccs    map[int]int
	last   int  // the index of the last note when nothing came after it, or -1
	tied   bool // & after the last note
	done   bool // the track is too long: the rest is left out
	pct    int  // the ticks of a whole note in % lengths
	moved  bool // an octave change since the last note
	before int  // the key of the last note, or -1
}

func (p *player) voice(ts *trackSource) *voice {
	v := &voice{p: p, tr: &music.Track{Name: ts.name, Volume: -1, Pan: -1}, kind: ts.kind, fold: p.dialect != PPMCK,
		oct: 4, len: music.PPQ, ccs: map[int]int{}, last: -1, before: -1, pct: 384}
	prog := 80
	switch p.dialect {
	case Mabinogi:
		v.setVol(8)
		prog = 0
	case PPMCK:
		v.setVol(15)
		v.pct = 192
		switch ts.kind {
		case chipTriangle:
			prog = 82
		case chipNoise:
			prog = 122
		}
	default:
		v.setVol(12) // FlMML's velocity of 100
	}
	if p.program >= 0 {
		prog = p.program
	}
	v.tr.Program, v.prog = prog, prog
	if ts.kind == chipPulse && prog == 80 {
		v.cc(70, dutyCC[0])
	}
	return v
}

// dutyCC are the values of controller 70 for the duties 12.5%, 25%, 50%
// and 75% (docs/spec.md §4.4).
var dutyCC = [4]int{32, 64, 0, 96}

// run plays MML.
func (v *voice) run(s []byte, depth int) {
	p := v.p
	if depth > maxDepth {
		p.w.once("depth", "loops and tuplets nested more than %d deep are left out", maxDepth)
		return
	}
	for i := 0; ; {
		if p.stop || v.done {
			return
		}
		// a step for each command, and one for each pass of a loop
		if p.steps++; p.steps > maxSteps {
			p.w.once("steps", "the MML plays too long; the rest is left out")
			p.stop = true
			return
		}
		if i >= len(s) {
			return
		}
		c := s[i]
		if c >= utf8.RuneSelf {
			r, n := utf8.DecodeRune(s[i:])
			i += n
			if !unicode.IsSpace(r) {
				p.w.once("text", "text that is not MML (%q) is skipped", r)
			}
			continue
		}
		i++
		if 'A' <= c && c <= 'Z' {
			if !v.fold {
				i = v.upper(s, i-1)
				continue
			}
			c += 'a' - 'A'
		}
		switch c {
		case ' ', '\t', '\n', '\v', '\f', ',', '|', ']', '}', '.', '+', '-', '#', '$':
			// separators, and the breaks and ends of loops and tuplets
			// outside them
		case 'c', 'd', 'e', 'f', 'g', 'a', 'b':
			acc := 0
		accidentals:
			for ; i < len(s); i++ {
				switch s[i] {
				case '+', '#':
					acc++
				case '-':
					acc--
				default:
					break accidentals
				}
			}
			key := (v.oct+1)*12 + int(semitones[c-'a']) + acc + v.shift
			v.note(key, v.length(s, &i), cmp.Compare(acc, 0))
		case 'r':
			v.rest(v.length(s, &i))
		case 'w':
			if p.dialect != PPMCK {
				v.unknown(c, s, &i)
				break
			}
			v.rest(v.length(s, &i)) // wait
		case 'n':
			if p.dialect == Generic && i < len(s) && s[i]|0x20 == 's' { // FlMML's note shift
				i++
				if n, ok := signed(s, &i); ok {
					v.shift = min(max(v.shift+n, -127), 127)
				}
				break
			}
			n, ok := number(s, &i)
			if !ok {
				v.unknown(c, s, &i)
				break
			}
			d := v.len
			if i+1 < len(s) && s[i] == ',' && (isDigit(s[i+1]) || s[i+1] == '%') {
				i++
				d = v.length(s, &i)
			} else {
				d = dots(s, &i, d)
			}
			v.moved = false
			v.note(n+12+v.shift, d, 0)
		case 'o':
			if n, ok := signed(s, &i); ok {
				v.oct = min(max(n, -2), 10)
				v.moved = false
			}
		case '<', '>':
			d := p.up
			if c == '<' {
				d = -d
			}
			v.oct = min(max(v.oct+d, -3), 12)
			v.moved = true
		case 'l':
			if d, ok := v.lengthArg(s, &i); ok {
				v.len = d
			}
		case 't':
			if f, ok := decimal(s, &i); ok && f >= 1 {
				p.tempos = append(p.tempos, music.Tempo{Tick: v.tick, BPM: min(f, 999)})
			}
		case 'v':
			if i < len(s) && (s[i] == '+' || s[i] == '-') {
				sign := 1
				if s[i] == '-' {
					sign = -1
				}
				i++
				n, ok := number(s, &i)
				if !ok {
					n = 1
				}
				v.setVol(v.vol + sign*n)
			} else if n, ok := number(s, &i); ok {
				v.setVol(n)
			}
		case '(', ')':
			if p.dialect != Generic {
				v.unknown(c, s, &i)
				break
			}
			n, ok := number(s, &i)
			if !ok {
				n = 1
			}
			if c == ')' { // FlMML: ( raises the volume
				n = -n
			}
			v.setVol(v.vol + n)
		case 'q':
			skipArgs(s, &i) // gate time: notes last their written length
		case 'x', 's', 'y':
			if c != 'x' && p.dialect != PPMCK || c == 'x' && p.dialect == Mabinogi {
				v.unknown(c, s, &i)
				break
			}
			skipArgs(s, &i) // FlMML's volume mode, PPMCK's sweep and memory writes
		case '@':
			v.at(s, &i)
		case '&':
			v.tied = true
		case '^':
			v.extend(v.length(s, &i))
		case '[':
			end := closing(s, i-1, '[', ']')
			p.scanned(s, i, end)
			if end < 0 {
				p.w.once("open-loop", "a loop that does not end ([ without ]) is played once")
				break
			}
			body := s[i:end]
			i = end + 1
			n, ok := number(s, &i)
			if !ok {
				n = 2
			}
			v.repeat(body, n, breakAt(body, '['), depth)
		case '/':
			if i < len(s) && s[i] == ':' { // FlMML's /:n … / … :/
				end := closingSlash(s, i-1)
				p.scanned(s, i, end)
				i++
				n, ok := number(s, &i)
				if !ok {
					n = 2
				}
				if end < 0 || end < i {
					p.w.once("open-loop", "a loop that does not end ([ without ]) is played once")
					break
				}
				body := s[i:end]
				i = end + 2
				v.repeat(body, n, breakAt(body, '/'), depth)
			}
		case ':':
			if i < len(s) && s[i] == '/' {
				i++
			}
		case '{':
			end := closing(s, i-1, '{', '}')
			p.scanned(s, i, end)
			if end < 0 {
				p.w.once("open-tuplet", "a tuplet that does not end ({ without }) is played as it is")
				break
			}
			body := s[i:end]
			i = end + 1
			v.tuplet(body, s, &i, depth)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			number(s, &i) // stray
		default:
			v.unknown(c, s, &i)
		}
	}
}

// scanned charges the search for the end of a loop or tuplet from s[i]
// (to s[end], or the end of s) to the steps: loops and breaks are looked
// for again on each pass of the loops around them.
func (p *player) scanned(s []byte, i, end int) {
	if end < 0 {
		end = len(s)
	}
	p.steps += 2 * max(end-i, 0)
}

// semitones are the semitones above C of the notes a to g.
var semitones = [7]int8{9, 11, 0, 2, 4, 5, 7}

// upper plays a PPMCK command in upper case from s[i] and returns where it
// ends.
func (v *voice) upper(s []byte, i int) int {
	j := i
	for j < len(s) && 'A' <= s[j] && s[j] <= 'Z' {
		j++
	}
	switch name := string(s[i:j]); name {
	case "K": // transposition in semitones
		if n, ok := signed(s, &j); ok {
			v.shift = min(max(n, -127), 127)
		}
	case "L": // the loop point: the music plays once
	case "D", "EN", "ENOF", "EP", "EPOF", "MP", "MPOF", "MH", "MHOF", "NB", "SA", "SD", "OP":
		skipArgs(s, &j) // detune, envelopes, modulation, sweeps
	default:
		v.p.w.once("cmd:"+name, "unknown command %q is skipped", name)
		skipArgs(s, &j)
	}
	return j
}

// at plays the @ commands: tones, FlMML's fine volume, pulse width and
// pan, and PPMCK's duty.
func (v *voice) at(s []byte, i *int) {
	p := v.p
	if p.dialect == Mabinogi {
		v.unknown('@', s, i)
		return
	}
	if *i < len(s) && isDigit(s[*i]) {
		n, _ := number(s, i)
		if p.dialect == PPMCK {
			if v.kind == chipPulse && v.prog == 80 && n <= 3 {
				v.cc(70, dutyCC[n])
			}
			return
		}
		if *i+1 < len(s) && s[*i] == '-' && isDigit(s[*i+1]) { // FlMML's @9-3: the wave of the tone
			*i++
			number(s, i)
		}
		v.setProgram(toneProgram(n))
		return
	}
	j := *i
	for *i < len(s) && (isLetter(s[*i]) || s[*i] == '@') {
		*i++
	}
	name := strings.ToLower(string(s[j:*i]))
	if p.dialect == PPMCK { // @v @q @@ @EN …: envelopes and tones by number
		skipArgs(s, i)
		return
	}
	switch name {
	case "v": // fine volume 0–127
		if n, ok := number(s, i); ok {
			v.vel = min(n, 127)
		}
	case "w": // pulse width
		if n, ok := number(s, i); ok {
			v.cc(70, pulseWidth(n))
		}
	case "p": // pan
		if n, ok := number(s, i); ok {
			v.pan(min(n, 127))
		}
	case "e", "l", "f", "q", "n", "d", "o", "r", "s", "x", "i", "u", "m", "mh":
		skipArgs(s, i) // envelopes, LFOs, filters, detune, outputs
	default:
		p.w.once("cmd:@"+name, "unknown command %q is skipped", "@"+name)
		skipArgs(s, i)
	}
}

// toneProgram is the General MIDI program of a FlMML tone.
func toneProgram(n int) int {
	switch n {
	case 0: // sine
		return 79
	case 1: // saw
		return 81
	case 2, 6: // triangle, NES triangle
		return 82
	case 4, 7, 10, 11: // white noise, NES noise, Game Boy noises
		return 122
	}
	return 80 // pulse, NES pulse, and the rest (FM, wave tables, samples)
}

// pulseWidth is the controller 70 value of FlMML's @W: a NES duty (0–3)
// or a width in percent.
func pulseWidth(n int) int {
	switch {
	case n <= 3:
		return dutyCC[n]
	case n <= 18:
		return dutyCC[0]
	case n <= 37:
		return dutyCC[1]
	case n <= 62:
		return dutyCC[2]
	case n <= 87:
		return dutyCC[3]
	}
	return dutyCC[0]
}

func (v *voice) unknown(c byte, s []byte, i *int) {
	v.p.w.once("cmd:"+string(c), "unknown command %q is skipped", string(c))
	skipArgs(s, i)
}

// note plays a note. A note tied (&) to the one before of the same pitch
// lengthens it.
func (v *voice) note(key, dur, spell int) {
	p := v.p
	if v.moved && v.before >= 0 {
		p.cost += abs(key - v.before)
	}
	v.moved, v.before = false, key
	tied := v.tied
	v.tied = false
	if tied && v.last >= 0 {
		if n := &v.tr.Notes[v.last]; n.Key == key && n.Tick+n.Dur == v.tick {
			n.Dur += dur
			v.advance(dur)
			return
		}
	}
	v.last = -1
	switch {
	case key < 0 || key > 127:
		p.w.once("range", "notes outside the MIDI range are left out")
	case v.vel <= 0: // v0: silent
	default:
		if p.notes++; p.notes > maxNotes {
			p.w.once("notes", "the MML plays more than %d notes; the rest is left out", maxNotes)
			p.stop = true
			return
		}
		v.tr.Notes = append(v.tr.Notes, music.PlayNote{Tick: v.tick, Dur: dur, Key: key, Vel: v.vel, Spell: spell})
		v.last = len(v.tr.Notes) - 1
	}
	v.advance(dur)
}

func (v *voice) rest(dur int) {
	v.last, v.tied = -1, false
	v.advance(dur)
}

// extend lengthens the last note (^), or the rest.
func (v *voice) extend(dur int) {
	if v.last >= 0 {
		if n := &v.tr.Notes[v.last]; n.Tick+n.Dur == v.tick {
			n.Dur += dur
		}
	}
	v.advance(dur)
}

func (v *voice) advance(dur int) {
	v.tick += dur
	if v.tick > maxTicks {
		v.p.w.once("ticks", "a track plays longer than %d whole notes; the rest is left out", maxTicks/whole)
		v.done = true
	}
}

// repeat plays a loop n times, the last pass up to its break.
func (v *voice) repeat(body []byte, n, brk, depth int) {
	n = min(max(n, 1), maxRepeat)
	for k := 0; k < n && !v.p.stop && !v.done; k++ {
		if k == n-1 && brk >= 0 {
			v.run(body[:brk], depth+1)
		} else {
			v.run(body, depth+1)
		}
	}
}

// tuplet plays the notes of a tuplet in the length that follows it (s[*i]),
// sharing it in the proportion of their own lengths.
func (v *voice) tuplet(body, s []byte, i *int, depth int) {
	start, n0, c0, t0 := v.tick, len(v.tr.Notes), len(v.tr.Controls), len(v.p.tempos)
	v.run(body, depth+1)
	total := v.length(s, i)
	inner := v.tick - start
	switch {
	case v.p.stop || v.done:
		return
	case inner <= 0:
		v.advance(total)
		return
	}
	scale := func(t int) int { return start + int(int64(t-start)*int64(total)/int64(inner)) }
	notes := v.tr.Notes
	if n0 > 0 { // a note before tied into the tuplet
		if n := &notes[n0-1]; n.Tick+n.Dur > start {
			n.Dur = scale(n.Tick+n.Dur) - n.Tick
		}
	}
	out := notes[:n0]
	for k, n := range notes[n0:] {
		end := scale(n.Tick + n.Dur)
		n.Tick = scale(n.Tick)
		n.Dur = end - n.Tick
		if n.Dur > 0 {
			out = append(out, n)
		} else if n0+k == v.last {
			v.last = -1
		}
	}
	v.tr.Notes = out
	if v.last >= n0 { // the last note of the tuplet, kept: the last one
		v.last = len(out) - 1
	}
	for k := c0; k < len(v.tr.Controls); k++ {
		v.tr.Controls[k].Tick = scale(v.tr.Controls[k].Tick)
	}
	for k := t0; k < len(v.p.tempos); k++ {
		v.p.tempos[k].Tick = scale(v.p.tempos[k].Tick)
	}
	v.tick = start
	v.advance(total)
}

// setVol sets the volume (v, 0–15) and the velocity it plays.
func (v *voice) setVol(n int) {
	v.vol = min(max(n, 0), 16)
	v.vel = min((v.vol*127+7)/15, 127)
}

// setProgram changes the program: at the start of the track it is the
// track's.
func (v *voice) setProgram(prog int) {
	if v.tick == 0 && len(v.tr.Notes) == 0 {
		v.tr.Program, v.prog = prog, prog
		return
	}
	if prog == v.prog {
		return
	}
	v.prog = prog
	if n := len(v.tr.Controls); n > 0 && v.tr.Controls[n-1].Tick == v.tick && v.tr.Controls[n-1].Kind == music.ControlProgram {
		v.tr.Controls[n-1].Value = prog
		return
	}
	v.tr.Controls = append(v.tr.Controls, music.Control{Tick: v.tick, Kind: music.ControlProgram, Value: prog})
}

func (v *voice) pan(val int) {
	if v.tick == 0 && len(v.tr.Notes) == 0 {
		v.tr.Pan = val
		v.ccs[10] = val
		return
	}
	v.cc(10, val)
}

// cc sets a controller, replacing a change of it at the same tick.
func (v *voice) cc(num, val int) {
	if cur, ok := v.ccs[num]; ok && cur == val {
		return
	}
	v.ccs[num] = val
	ctl := v.tr.Controls
	for k := len(ctl) - 1; k >= 0 && ctl[k].Tick == v.tick; k-- {
		if ctl[k].Kind == music.ControlChange && ctl[k].Num == num {
			ctl[k].Value = val
			return
		}
	}
	v.tr.Controls = append(v.tr.Controls, music.Control{Tick: v.tick, Kind: music.ControlChange, Num: num, Value: val})
}

// length reads the length of a note or rest (4, 8., %96), the default
// length when it has none.
func (v *voice) length(s []byte, i *int) int {
	if d, ok := v.lengthArg(s, i); ok {
		return d
	}
	return dots(s, i, v.len)
}

// lengthArg reads a length with its dots; ok is false when there is no
// number.
func (v *voice) lengthArg(s []byte, i *int) (int, bool) {
	j := *i
	if j < len(s) && s[j] == '%' {
		j++
		n, ok := number(s, &j)
		if !ok {
			return 0, false
		}
		*i = j
		return dots(s, i, max(n*whole/v.pct, 1)), true
	}
	n, ok := number(s, &j)
	if !ok || n == 0 {
		if ok {
			*i = j
		}
		return 0, false
	}
	*i = j
	return dots(s, i, max(whole/n, 1)), true
}

// dots adds the dots after a length.
func dots(s []byte, i *int, d int) int {
	h := d
	for *i < len(s) && s[*i] == '.' {
		*i++
		h /= 2
		d += h
	}
	return d
}

// number reads an unsigned number, saturating at maxNumber.
func number(s []byte, i *int) (int, bool) {
	j, n := *i, 0
	for j < len(s) && isDigit(s[j]) {
		n = min(n*10+int(s[j]-'0'), maxNumber)
		j++
	}
	if j == *i {
		return 0, false
	}
	*i = j
	return n, true
}

// signed reads a number with an optional sign.
func signed(s []byte, i *int) (int, bool) {
	j, sign := *i, 1
	if j < len(s) && (s[j] == '+' || s[j] == '-') {
		if s[j] == '-' {
			sign = -1
		}
		j++
	}
	n, ok := number(s, &j)
	if !ok {
		return 0, false
	}
	*i = j
	return sign * n, true
}

// decimal reads a number that may have a fraction (t120.5).
func decimal(s []byte, i *int) (float64, bool) {
	j := *i
	for j < len(s) && isDigit(s[j]) && j-*i < 9 {
		j++
	}
	if j == *i {
		return 0, false
	}
	if j+1 < len(s) && s[j] == '.' && isDigit(s[j+1]) {
		j++
		for k := 0; j < len(s) && isDigit(s[j]); j++ {
			if k++; k > 6 {
				break
			}
		}
	}
	f, err := strconv.ParseFloat(string(s[*i:j]), 64)
	for j < len(s) && isDigit(s[j]) {
		j++
	}
	*i = j
	return f, err == nil
}

// skipArgs skips the arguments of a command: numbers with signs, commas
// and dots.
func skipArgs(s []byte, i *int) {
	for *i < len(s) && (isDigit(s[*i]) || strings.IndexByte("+-,.", s[*i]) >= 0) {
		*i++
	}
}

// closing returns the index of the bracket that closes the one at s[open],
// or -1.
func closing(s []byte, open int, o, c byte) int {
	depth := 0
	for i := open; i < len(s); i++ {
		switch s[i] {
		case o:
			depth++
		case c:
			if depth--; depth == 0 {
				return i
			}
		}
	}
	return -1
}

// closingSlash returns the index of the :/ that closes the /: at s[open],
// or -1.
func closingSlash(s []byte, open int) int {
	depth := 0
	for i := open; i+1 < len(s); i++ {
		switch {
		case s[i] == '/' && s[i+1] == ':':
			depth++
			i++
		case s[i] == ':' && s[i+1] == '/':
			if depth--; depth == 0 {
				return i
			}
			i++
		}
	}
	return -1
}

// breakAt returns the index of the break of a loop body (| or : in [ ],
// / in /: :/), or -1.
func breakAt(b []byte, kind byte) int {
	depth := 0
	for i := 0; i < len(b); i++ {
		c := b[i]
		switch {
		case c == '/' && i+1 < len(b) && b[i+1] == ':':
			depth++
			i++
		case c == ':' && i+1 < len(b) && b[i+1] == '/':
			depth--
			i++
		case c == '[' || c == '{':
			depth++
		case c == ']' || c == '}':
			depth--
		case depth == 0 && kind == '[' && (c == '|' || c == ':'), depth == 0 && kind == '/' && c == '/':
			return i
		}
	}
	return -1
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// detectDialect tells the dialect of MML text: Mabinogi's MML@, PPMCK's
// headers, definitions and track names, or else the generic dialect.
func detectDialect(text string) string {
	if mabinogiAt(text, 0) >= 0 {
		return Mabinogi
	}
	var ppmck, generic, tracks, others int
	for _, l := range strings.Split(text, "\n") {
		t := strings.TrimSpace(l)
		switch {
		case t == "" || t[0] == ';' || strings.HasPrefix(t, "//") || strings.HasPrefix(t, "/*"):
		case t[0] == '#':
			switch headerDialect(t) {
			case PPMCK:
				ppmck++
			case Generic:
				generic++
			}
		case isDefinition(t):
			ppmck++
		case isTrackLine(l):
			tracks++
		default:
			others++
		}
	}
	if ppmck > generic || generic == 0 && tracks >= 2 && tracks >= others {
		return PPMCK
	}
	return Generic
}

// headerDialect tells the dialect a header line belongs to: PPMCK or
// Generic (FlMML), "" for either (#TITLE) or none.
func headerDialect(t string) string {
	name, _, _ := strings.Cut(strings.ToUpper(t), " ")
	name, _, _ = strings.Cut(name, "\t")
	switch {
	case name == "#TITLE":
		return ""
	case slices.Contains(ppmckHeaders, name), strings.HasPrefix(name, "#EX-"):
		return PPMCK
	case slices.Contains(flmmlHeaders, name), strings.HasPrefix(name, "#WAV"), strings.HasPrefix(name, "#OPM@"), strings.HasPrefix(name, "#OPN@"):
		return Generic
	}
	return ""
}

var (
	ppmckHeaders = []string{"#COMPOSER", "#MAKER", "#PROGRAMER", "#PROGRAMMER", "#OCTAVE-REV", "#INCLUDE", "#BANK-CHANGE",
		"#SETBANK", "#EFFECT-INCLUDE", "#DPCM-RESTSTOP", "#GATE-DENOM", "#PITCH-CORRECTION", "#AUTO-BANKSWITCH", "#NO-BANKSWITCH"}
	flmmlHeaders = []string{"#ARTIST", "#CODING", "#COMMENT", "#OCTAVE", "#VELOCITY", "#FMGAIN", "#USING", "#PRAGMA"}
)

// isDefinition reports whether a line starts a PPMCK definition: @v0 = {,
// @EN1 = {, @DPCM0 = {.
func isDefinition(t string) bool {
	if len(t) < 2 || t[0] != '@' {
		return false
	}
	i := 1
	for i < len(t) && isLetter(t[i]) && i < 8 {
		i++
	}
	j := i
	for j < len(t) && isDigit(t[j]) {
		j++
	}
	if j == i {
		return false
	}
	rest := strings.TrimLeft(t[j:], " \t")
	return strings.HasPrefix(rest, "=")
}

// isTrackLine reports whether a line is a PPMCK track line: track names,
// then MML.
func isTrackLine(l string) bool {
	names, rest := trackNames(l)
	if names == "" {
		return false
	}
	if k := strings.IndexByte(rest, ';'); k >= 0 {
		rest = rest[:k]
	}
	ok, tokens, _ := mmlLike(rest)
	return ok && tokens >= 1
}

// mmlLike reports whether text is made of the letters and signs of MML,
// with the number of commands with a number (c4, l8, o5, @2) and of notes
// in it.
func mmlLike(t string) (ok bool, tokens, notes int) {
	for i := 0; i < len(t); i++ {
		c := t[i]
		switch {
		case 'a' <= c && c <= 'g' || 'A' <= c && c <= 'G':
			notes++
		case isDigit(c) || c == ' ' || c == '\t' || strings.IndexByte("<>[]{}&^.+-#|:/%,@;$=()_!", c) >= 0:
		case strings.IndexByte("rolnqtvkswxypRLOMNPQTVKSWXY", c) >= 0:
		default:
			return false, 0, 0
		}
		if i+1 < len(t) && isDigit(t[i+1]) && strings.IndexByte("abcdefgrlotvnq@ABCDEFGRLOTVNQ", c) >= 0 {
			tokens++
		}
	}
	return true, tokens, notes
}
