// Package musicxml reads MusicXML scores (.musicxml, .xml, and compressed
// .mxl) into the music engine's Score, which converter/internal/music
// engraves as staff notation on pages with the music played by the viewer.
//
// It reads partwise and timewise scores (timewise ones are turned
// partwise), in UTF-8, UTF-16 or a declared legacy encoding, with the
// DOCTYPE and entities of any exporter. From the header it takes the
// title, subtitle and credits (typed credits first, then work and movement
// titles and the identification), and the parts with their names and MIDI
// instruments. From the music it takes what the score shows: time and key
// signatures, clefs, staves, notes, chords, rests, voices, grace and cue
// notes, tuplets, beams, stems, accidentals, ties, slurs, articulations,
// fermatas, ornaments, single-note tremolos, lyrics, dynamics, words,
// hairpins, pedal marks, octave lines, segno and coda signs, chord
// symbols, tempo marks, rehearsal marks, bar lines, repeats, voltas and
// system and page breaks. From the sound it takes the tempo, the dynamics
// and the jumps (D.C., D.S., Fine, To Coda) for the performance, played in
// the order the repeats and jumps give.
//
// It does not read: layout and formatting (positions, fonts, colours,
// system and staff distances are the engine's), figured bass, fingerings
// and other technical marks, arpeggios, glissandos, multi-measure rest
// styles, non-traditional key signatures, nested tuplets (the outer one
// is kept), cross-staff chords (a chord keeps the staff of its first
// note), mid-part instrument changes and MusicXML opus files.
//
// Decisions the format leaves open:
//
//   - Measure numbers are those of the file (the text attribute, else the
//     number attribute), always.
//   - A measure is as long as its time signature, unless what its parts
//     hold is shorter (a pickup, an incomplete last measure) or longer (a
//     cadenza, senza misura): then it is as long as that.
//   - Accidentals are drawn where the file has an accidental element; a
//     note without one has none. A file with altered notes and no
//     accidental elements at all (some exporters leave them out) has its
//     accidentals chosen by the engine from the key signature instead.
//   - Beams are as the file has them. When the file has no beam elements
//     at all, the engine beams the notes; when it has some, a note
//     without any is not beamed.
//   - Octave lines follow MusicXML: the pitch of a note is its sounding
//     pitch, and an octave-shift of type "down" (an 8va line above the
//     staff, Direction.Size 8) shows the notes an octave lower than they
//     sound; "up" (8vb below, Size -8) an octave higher. The notes are
//     written shifted, and played at their pitch.
//   - Grace notes are played before their principal note, a sixteenth
//     each at most, in time taken from it; cue notes are not played.
//   - Chord symbols are written with "#" and "b" for sharps and flats.
package musicxml

import (
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/shibukawa/tinygodriver/encoding/xmlro"

	"github.com/shibukawa/bdf/converter/internal/music"
)

// Options controls how a score is read.
type Options struct {
	// NoPlay leaves the performance out: Score.Play and Score.PlayOrder
	// stay nil.
	NoPlay bool
}

// Budgets against hostile inputs.
const (
	maxParts     = 256
	maxMeasures  = 100000
	maxNotes     = 2000000
	maxStaves    = 16
	maxQuarters  = 4096 // the furthest position in a measure
	maxPlays     = 10000
	maxPlayNotes = 4000000

	maxInMeasure    = 10000   // notes, slurs, syllables, directions and clef changes of a part in a measure
	maxTempoMarks   = 64      // tempo marks of a measure
	maxVerses       = 99      // verses of lyrics
	maxEndings      = 100     // numbers of a volta
	maxPlayQuarters = 1000000 // the length of the performance in quarter notes
)

// Budgets that are variables for the tests.
var (
	// maxStaffMeasures bounds the staves times the measures of the score,
	// as the engine does, and the measures of the parts read.
	maxStaffMeasures = 1_000_000
	// maxPlayChanges bounds the tempo changes of the performance, and the
	// syllables sung: repeats play those of a measure many times.
	maxPlayChanges = 1_000_000
)

// Parse reads a MusicXML score, uncompressed or compressed. It returns
// the score with its performance (unless o.NoPlay) and the problems it
// read past.
func Parse(r io.ReaderAt, size int64, o *Options) (*music.Score, []string, error) {
	if o == nil {
		o = &Options{}
	}
	data, err := readInput(r, size)
	if err != nil {
		return nil, nil, err
	}
	rd := newReader()
	data, note := toUTF8(data)
	if note != "" {
		rd.warn("%s", note)
	}
	if err := rd.read(newXMLReader(data)); err != nil {
		return nil, rd.warnings, err
	}
	rd.finish()
	if !o.NoPlay {
		rd.perform()
	}
	return rd.score, rd.warnings, nil
}

// reader is the state of reading a score.
type reader struct {
	score    *music.Score
	warnings []string
	warned   map[string]bool

	parts []*partState // in the order of the part list
	byID  map[string]*partState
	ms    []*measureInfo // what the parts share, one per measure

	workTitle, movementTitle string
	credits                  map[string]string // by credit type
	untyped                  []credit
	creators                 map[string]string
	rights                   string

	notes                                int
	partMeasures                         int // the measures of the parts read, against maxStaffMeasures
	hasBeams, hasAccidentals, hasAltered bool
	verses                               *numbers
	gotMusic                             bool
	timewiseCount                        int // the measures of a timewise score read
}

func newReader() *reader {
	return &reader{score: &music.Score{}, warned: map[string]bool{}, byID: map[string]*partState{},
		credits: map[string]string{}, creators: map[string]string{}, verses: newNumbers(maxVerses)}
}

// warn reports a problem once.
func (r *reader) warn(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if r.warned[msg] {
		return
	}
	r.warned[msg] = true
	r.warnings = append(r.warnings, "musicxml: "+msg)
}

// read reads the document: the header whole, the music a measure at a
// time.
func (r *reader) read(d *xmlro.Reader) error {
	for {
		k, err := d.Next()
		if err != nil {
			return fmt.Errorf("musicxml: %w", err)
		}
		if k == xmlro.EOF {
			return errors.New("musicxml: no score in the file")
		}
		if k != xmlro.StartElement {
			continue
		}
		switch name := string(d.LocalName()); name {
		case "score-partwise":
			return r.body(d, false)
		case "score-timewise":
			return r.body(d, true)
		case "opus":
			return errors.New("musicxml: the file is an opus (a list of scores), not a score")
		default:
			return fmt.Errorf("musicxml: the document is a <%s>, not a score", name)
		}
	}
}

// body reads the children of the root.
func (r *reader) body(d *xmlro.Reader, timewise bool) error {
	for {
		k, err := next(d)
		if err != nil {
			return r.early(err)
		}
		switch k {
		case xmlro.EndElement:
			if len(r.parts) == 0 {
				return errors.New("musicxml: the score has no parts")
			}
			return nil
		case xmlro.StartElement:
			part, measure := xmlro.Equal(d.LocalName(), "part"), xmlro.Equal(d.LocalName(), "measure")
			switch {
			case part && !timewise:
				if err := r.partwise(d); err != nil {
					return r.early(err)
				}
			case measure && timewise:
				n, err := readNode(d)
				if err != nil {
					return r.early(err)
				}
				r.timewise(n)
			case part || measure:
				if err := d.Skip(); err != nil {
					return r.early(err)
				}
			default:
				n, err := readNode(d)
				if err != nil {
					return r.early(err)
				}
				r.header(n)
			}
		}
	}
}

// next reads a token inside the root element, where the end of the
// document is an error: io.EOF, as for a document that ends in the middle
// of an element.
func next(d *xmlro.Reader) (xmlro.Kind, error) {
	k, err := d.Next()
	if err == nil && k == xmlro.EOF || errors.Is(err, xmlro.ErrTruncated) {
		return k, io.EOF
	}
	return k, err
}

// early handles the end of a document that stops short: what was read of
// the music is kept.
func (r *reader) early(err error) error {
	if errors.Is(err, errTooLarge) {
		return err
	}
	if r.gotMusic && len(r.parts) > 0 {
		if err == io.EOF {
			r.warn("the document ends early")
		} else {
			r.warn("the document ends early: %v", err)
		}
		return nil
	}
	if err == io.EOF {
		return errors.New("musicxml: the document ends early")
	}
	return fmt.Errorf("musicxml: %w", err)
}

// partwise reads a <part> of a partwise score a measure at a time.
func (r *reader) partwise(d *xmlro.Reader) error {
	var id string
	for _, a := range attrsOf(d) {
		if a.name == "id" {
			id = a.value
		}
	}
	p := r.part(id)
	for {
		k, err := next(d)
		if err != nil {
			return err
		}
		switch k {
		case xmlro.EndElement:
			return nil
		case xmlro.StartElement:
			if !xmlro.Equal(d.LocalName(), "measure") || p == nil {
				if err := d.Skip(); err != nil {
					return err
				}
				continue
			}
			n, err := readNode(d)
			if err != nil {
				return err
			}
			r.measure(p, p.count, n, n)
			p.count++
		}
	}
}

// timewise reads a <measure> of a timewise score: the parts in it.
func (r *reader) timewise(n *node) {
	idx := r.timewiseCount
	r.timewiseCount++
	for _, k := range n.kids {
		if k.name != "part" {
			continue
		}
		if p := r.part(k.attr("id")); p != nil {
			r.measure(p, idx, n, k)
			p.count = idx + 1
		}
	}
}

// part returns the part with an id, adding it when the part list does
// not name it; nil past the budget.
func (r *reader) part(id string) *partState {
	if p := r.byID[id]; p != nil {
		p.seen = true
		return p
	}
	if len(r.parts) >= maxParts {
		r.warn("more than %d parts; the rest are left out", maxParts)
		return nil
	}
	p := newPart(id, id)
	p.seen = true
	r.byID[id] = p
	r.parts = append(r.parts, p)
	return p
}

// header reads an element before the music.
func (r *reader) header(n *node) {
	switch n.name {
	case "work":
		if t := collapse(n.textOf("work-title")); t != "" {
			r.workTitle = t
		}
	case "movement-title":
		r.movementTitle = collapse(n.text)
	case "identification":
		for _, k := range n.kids {
			switch k.name {
			case "creator":
				typ := k.attr("type")
				if typ == "poet" {
					typ = "lyricist"
				}
				if t := collapse(k.text); t != "" && r.creators[typ] == "" {
					r.creators[typ] = t
				}
			case "rights":
				if t := collapse(k.text); t != "" && r.rights == "" {
					r.rights = t
				}
			}
		}
	case "credit":
		r.credit(n)
	case "part-list":
		r.partList(n)
	}
}

// credit is an untyped credit: text on the first page.
type credit struct {
	text     string
	size     float64
	centered bool
	right    bool
}

// credit reads the words of a credit: typed credits fill the credits
// they name; untyped ones on the first page are kept for a guess.
func (r *reader) credit(n *node) {
	var words []string
	var first *node
	for _, k := range n.kids {
		if k.name == "credit-words" {
			if first == nil {
				first = k
			}
			if t := collapse(k.text); t != "" {
				words = append(words, t)
			}
		}
	}
	text := strings.Join(words, " ")
	if text == "" {
		return
	}
	typed := false
	for _, k := range n.kids {
		if k.name != "credit-type" {
			continue
		}
		typed = true
		switch typ := strings.ToLower(k.trim()); typ {
		case "title", "subtitle", "composer", "lyricist", "arranger", "rights":
			if r.credits[typ] == "" {
				r.credits[typ] = text
			}
		case "poet":
			if r.credits["lyricist"] == "" {
				r.credits["lyricist"] = text
			}
		}
	}
	if typed || (n.attr("page") != "" && n.attr("page") != "1") {
		return
	}
	c := credit{text: text}
	c.size, _ = strconv.ParseFloat(first.attr("font-size"), 64)
	j := first.attr("justify")
	if j == "" {
		j = first.attr("halign")
	}
	c.centered, c.right = j == "center", j == "right"
	r.untyped = append(r.untyped, c)
}

// partList reads the parts, their names and MIDI instruments.
func (r *reader) partList(n *node) {
	for _, sp := range n.kids {
		if sp.name != "score-part" {
			continue
		}
		id := sp.attr("id")
		if r.byID[id] != nil {
			continue
		}
		if len(r.parts) >= maxParts {
			r.warn("more than %d parts; the rest are left out", maxParts)
			return
		}
		raw := collapse(sp.textOf("part-name"))
		p := newPart(id, raw)
		p.part.Name = displayName(sp, "part-name")
		p.part.Abbrev = displayName(sp, "part-abbreviation")
		if p.trackName == "" {
			p.trackName = p.part.Name
		}
		for _, k := range sp.kids {
			if k.name != "midi-instrument" {
				continue
			}
			in := &instrument{id: k.attr("id"), channel: -1, program: -1, unpitched: -1, volume: -1, pan: -1}
			if v, ok := atoi(k.textOf("midi-channel")); ok && v >= 1 && v <= 16 {
				in.channel = v - 1
			}
			if v, ok := atoi(k.textOf("midi-program")); ok && v >= 1 && v <= 128 {
				in.program = v - 1
			}
			if v, ok := atoi(k.textOf("midi-unpitched")); ok && v >= 1 && v <= 128 {
				in.unpitched = v - 1
			}
			if v, ok := parseFloat(k.textOf("volume")); ok {
				in.volume = min(max(int(math.Round(v*127/100)), 0), 127)
			}
			if v, ok := parseFloat(k.textOf("pan")); ok {
				v = min(max(v, -90), 90)
				in.pan = min(max(int(math.Round((v+90)*127/180)), 0), 127)
			}
			p.instruments = append(p.instruments, in)
			if in.id != "" {
				p.instByID[in.id] = in
			}
		}
		r.byID[id] = p
		r.parts = append(r.parts, p)
	}
}

// displayName returns the name a part shows: its display text when it
// has one, its name, or "" when it is hidden.
func displayName(sp *node, name string) string {
	n := sp.child(name)
	if n == nil || n.attr("print-object") == "no" {
		return ""
	}
	if disp := sp.child(name + "-display"); disp != nil {
		if disp.attr("print-object") == "no" {
			return ""
		}
		var b strings.Builder
		for _, k := range disp.kids {
			switch k.name {
			case "display-text":
				b.WriteString(k.text)
			case "accidental-text":
				b.WriteString(accidentalText(k.trim()))
			}
		}
		if t := collapse(b.String()); t != "" {
			return t
		}
	}
	return collapse(n.text)
}

// accidentalText writes an accidental as text.
func accidentalText(s string) string {
	switch s {
	case "sharp":
		return "#"
	case "flat":
		return "b"
	case "natural":
		return ""
	case "double-sharp", "sharp-sharp":
		return "##"
	case "flat-flat", "double-flat":
		return "bb"
	}
	return ""
}

// credits fills the credits of the score.
func (r *reader) setCredits() {
	s := r.score
	s.Title = r.credits["title"]
	s.Subtitle = r.credits["subtitle"]
	switch {
	case s.Title == "" && r.workTitle != "":
		s.Title = r.workTitle
		if s.Subtitle == "" && r.movementTitle != "" && r.movementTitle != r.workTitle {
			s.Subtitle = r.movementTitle
		}
	case s.Title == "" && r.movementTitle != "":
		s.Title = r.movementTitle
	}
	s.Composer = firstOf(r.credits["composer"], r.creators["composer"])
	s.Lyricist = firstOf(r.credits["lyricist"], r.creators["lyricist"])
	s.Arranger = firstOf(r.credits["arranger"], r.creators["arranger"])
	s.Rights = firstOf(r.credits["rights"], r.rights)
	// untyped credits: the largest centered text is the title, the first
	// on the right the composer
	var title *credit
	for i := range r.untyped {
		c := &r.untyped[i]
		if c.centered && (title == nil || c.size > title.size) {
			title = c
		}
	}
	if s.Title == "" && title != nil {
		s.Title = title.text
	}
	if s.Composer == "" {
		for _, c := range r.untyped {
			if c.right {
				s.Composer = c.text
				break
			}
		}
	}
}

func firstOf(s ...string) string {
	for _, v := range s {
		if v != "" {
			return v
		}
	}
	return ""
}

// collapse joins the white space of a text into single spaces.
func collapse(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func atoi(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		f, ok := parseFloat(s)
		if !ok || math.Abs(f) > 1e9 {
			return 0, false
		}
		return int(math.Round(f)), true
	}
	return v, true
}

func atoiDef(s string, def int) int {
	if v, ok := atoi(s); ok {
		return v
	}
	return def
}

// parseFloat reads a finite number.
func parseFloat(s string) (float64, bool) {
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, false
	}
	return v, true
}

// tickOf converts a position in quarter notes into ticks.
func tickOf(q float64) int {
	return int(math.Round(min(max(q, -maxQuarters), maxQuarters) * music.PPQ))
}

// finish completes the score once every part is read: the measure
// lengths and time signatures, the voltas, what each part lacks, and the
// policies that need the whole file (beams, accidentals).
func (r *reader) finish() {
	r.setCredits()
	s := r.score
	// every staff is laid out in every measure
	staves := 0
	for _, p := range r.parts {
		if p.seen {
			staves += p.part.Staves
		}
	}
	if n := maxStaffMeasures / max(staves, 1); len(r.ms) > n {
		r.warn("more than %d measures on the staves (%d staves of %d measures); the measures after %d are left out",
			maxStaffMeasures, staves, len(r.ms), n)
		r.ms, s.Measures = r.ms[:n], s.Measures[:n]
	}
	cur, free := music.TimeSig{Beats: 4, BeatType: 4}, false
	for i, mi := range r.ms {
		m := mi.m
		switch {
		case mi.time != nil:
			if i == 0 || *mi.time != cur || free {
				t := *mi.time
				m.Time = &t
			}
			cur, free = *mi.time, false
		case mi.free:
			free = true
		}
		tsLen := cur.Length()
		switch {
		case free && mi.content > 0:
			m.Length = mi.content
		case mi.content <= 0 || abs(mi.content-tsLen) <= 2:
			m.Length = tsLen
		default:
			m.Length = mi.content
		}
		sort.SliceStable(m.Tempo, func(a, b int) bool { return m.Tempo[a].Offset < m.Tempo[b].Offset })
		for j := range m.Tempo {
			m.Tempo[j].Offset = min(max(m.Tempo[j].Offset, 0), m.Length)
		}
	}
	r.endings()

	for _, p := range r.parts {
		if !p.seen {
			r.warn("part %q has no music; it is left out", p.id)
			continue
		}
		p.finish(r)
		s.Parts = append(s.Parts, p.part)
	}
	if len(s.Parts) == 0 {
		// a score whose parts are all empty keeps the first as a staff
		p := r.parts[0]
		p.finish(r)
		s.Parts = append(s.Parts, p.part)
	}
	for _, p := range s.Parts {
		for _, pm := range p.Measures {
			for _, ev := range pm.Events {
				if !r.hasBeams {
					ev.Beams = nil
				}
				if !r.hasAccidentals && r.hasAltered {
					for _, n := range ev.Notes {
						if n.Accidental == music.AccNone && !p.Percussion {
							n.Accidental = music.AccAuto
						}
					}
				}
			}
		}
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
