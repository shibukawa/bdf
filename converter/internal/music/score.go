package music

// Score is music as written (MusicXML, or a Performance after Notate).
type Score struct {
	// Title and the other credits are printed above the first system.
	Title, Subtitle, Composer, Lyricist, Arranger, Rights string

	// Measures are the measures every part shares: their lengths, time
	// signatures, bar lines, repeats and tempo marks.
	Measures []*Measure
	Parts    []*Part

	// Play is what the viewer plays, or nil for a score without sound.
	Play *Performance
	// PlayOrder is the order the measures are played in, when repeats make
	// it differ from their order in the score: the measure (an index into
	// Measures) played from each Tick of Play. nil plays every measure
	// once, in order, from tick 0.
	PlayOrder []PlayedMeasure
	// SMF is the Standard MIDI File stored for Play as it is (a MIDI
	// input). nil writes one from Play.
	SMF []byte

	// err is why Notate did not write the performance (it is too large);
	// Build returns it.
	err error
	// warnings are what Notate left out; Build reports them.
	warnings []string
}

// PlayedMeasure is a measure played from a tick of the performance.
type PlayedMeasure struct {
	Tick    int
	Measure int
}

// Measure is what the parts of a score share about one measure.
type Measure struct {
	// Number is the measure number shown; "" numbers the measure from the
	// previous one.
	Number string
	// Length is the length of the measure in ticks: that of its time
	// signature, or shorter in an incomplete (pickup) measure.
	Length int
	// Time is the time signature that starts at this measure, nil when it
	// does not change. The first measure's nil means 4/4.
	Time *TimeSig
	// Implicit is set on a measure that does not count in the numbering
	// (a pickup).
	Implicit bool

	// Left and Right are the bar lines at the start and end.
	Left, Right Barline
	// Times is how many times a repeat that ends here is played (0 for 2).
	Times int
	// Ending is a volta bracket over the measure.
	Ending *Ending

	// Tempo marks are printed above the top staff.
	Tempo []TempoMark
	// Rehearsal is a rehearsal mark ("A") printed in a box.
	Rehearsal string

	// Break asks for a new system (or page) to start at this measure.
	Break Break
}

// TimeSig is a time signature.
type TimeSig struct {
	Beats, BeatType int
	// Symbol is "common" (C) or "cut" (alla breve) to draw as a symbol, ""
	// for the numbers.
	Symbol string
}

// Length returns the length of a measure of the time signature in ticks.
func (t TimeSig) Length() int {
	if t.BeatType <= 0 || t.Beats <= 0 {
		return 4 * PPQ
	}
	return t.Beats * 4 * PPQ / t.BeatType
}

// KeySig is a key signature: the number of sharps (positive) or flats
// (negative), and whether the key is minor.
type KeySig struct {
	Fifths int
	Minor  bool
}

// Barline is the style of a bar line.
type Barline int

const (
	BarRegular Barline = iota
	BarNone
	BarDouble      // two thin lines
	BarFinal       // thin and thick (light-heavy)
	BarHeavyLight  // thick and thin
	BarRepeatStart // heavy-light with dots after
	BarRepeatEnd   // light-heavy with dots before
	BarRepeatBoth  // dots on both sides
	BarDashed
)

// Ending is a volta bracket.
type Ending struct {
	Text  string // "1.", "1, 2."
	Start bool   // the bracket starts at this measure
	Stop  bool   // the bracket ends at this measure
	Open  bool   // with Stop: no hook at the end (discontinue)
}

// TempoMark is a tempo indication: a text ("Allegro"), a metronome mark
// (a note value = BPM), or both.
type TempoMark struct {
	Offset int // ticks into the measure
	Text   string
	// Beat and Dots are the note value of the metronome mark and PerMinute
	// its count; PerMinute 0 has no metronome mark.
	Beat      NoteType
	Dots      int
	PerMinute float64
}

// Break is a break before a measure.
type Break int

const (
	BreakNone Break = iota
	BreakSystem
	BreakPage
)

// Part is one instrument of a score, on one or more staves (two for a
// piano's grand staff).
type Part struct {
	Name, Abbrev string
	Staves       int // 1 or more
	// Percussion parts are drawn on percussion clefs.
	Percussion bool
	// Measures has one entry for each measure of the score.
	Measures []*PartMeasure
}

// PartMeasure is what one part has in a measure.
type PartMeasure struct {
	// Key is the key signature that starts at this measure, nil when it
	// does not change. The first measure's nil means C major.
	Key *KeySig
	// Clefs are the clef changes in the measure, in the order of Offset;
	// the first measure has one at offset 0 for each staff.
	Clefs []ClefChange
	// Events are the notes, chords and rests of all voices and staves.
	Events []*Event
	// Directions are dynamics, words, hairpins, pedal marks and the like.
	Directions []*Direction
}

// ClefChange is a clef that starts at an offset of a measure on a staff.
type ClefChange struct {
	Offset int
	Staff  int // 0-based
	Clef   Clef
}

// Clef is a clef.
type Clef struct {
	Sign ClefSign
	// Line is the staff line the clef sits on, 1 (bottom) to 5; 0 is the
	// usual line of the sign (G 2, F 4, C 3).
	Line int
	// Octave is the octave the music sounds in relative to the clef: -1
	// for a treble clef with an 8 below (a tenor or guitar part).
	Octave int
}

// ClefSign is the sign of a clef.
type ClefSign int

const (
	ClefG ClefSign = iota
	ClefF
	ClefC
	ClefPercussion
	ClefTab
	ClefNone
)

// NoteType is a note value.
type NoteType int

const (
	Breve NoteType = iota // double whole
	Whole
	Half
	Quarter
	Eighth
	N16th
	N32nd
	N64th
	N128th
)

// Ticks returns the length of the note value with dots, in ticks.
func (t NoteType) Ticks(dots int) int {
	d := 8 * PPQ >> t
	n := d
	for i := 0; i < dots; i++ {
		d /= 2
		n += d
	}
	return n
}

// Event is a note, a chord or a rest in one voice of one staff.
type Event struct {
	// Offset is where the event starts in the measure and Duration how long
	// it lasts (as notated: tuplets scale it), in ticks. Grace notes have
	// Duration 0 and the Offset of the note they lead to.
	Offset, Duration int
	// Staff (0-based) and Voice (1-based) place the event.
	Staff, Voice int

	// Notes are the notes of a chord; a rest has none.
	Notes []*Note
	Rest  bool
	// MeasureRest is a rest that fills the measure (a whole rest centred
	// in it whatever the time signature).
	MeasureRest bool
	// Hidden events take time without being drawn.
	Hidden bool

	Type NoteType
	Dots int
	// Tuplet is the tuplet the event belongs to, nil for none.
	Tuplet *Tuplet
	Grace  bool
	// Slash is a grace note with a slash (an acciaccatura).
	Slash bool
	// Cue notes are drawn small.
	Cue bool

	// Beams are the beam states of the event per level (the eighth's beam
	// first). nil beams it automatically; an empty non-nil slice leaves it
	// unbeamed.
	Beams []Beam
	Stem  Stem

	// Artics are articulations and ornaments by SMuFL name without
	// "Above"/"Below" (articStaccato, articAccent, fermata, ornamentTrill).
	Artics []string
	// Lyrics are the syllables sung on the event.
	Lyrics []*LyricSyllable
	// Slurs start or end on the event.
	Slurs []SlurMark
	// Tremolo is the number of tremolo strokes on the stem.
	Tremolo int
}

// Tuplet is a group of events played in the time of fewer (3 in the time
// of 2). Events of one group share the pointer.
type Tuplet struct {
	Actual, Normal int
	// Bracket draws a bracket (default: when the group is not beamed as
	// one).
	Bracket *bool
	// ShowNumber is false to hide the number.
	ShowNumber bool
}

// Beam is the state of a beam level at an event.
type Beam int

const (
	BeamBegin Beam = iota + 1
	BeamContinue
	BeamEnd
	BeamForwardHook
	BeamBackwardHook
)

// Stem is the direction of a stem.
type Stem int

const (
	StemAuto Stem = iota
	StemUp
	StemDown
	StemNone
)

// Note is a note of an event.
type Note struct {
	Pitch Pitch
	// Accidental is the accidental drawn before the note. AccAuto chooses
	// it from the key signature and the notes before in the measure.
	Accidental Accidental
	// Cautionary accidentals are drawn in parentheses.
	Cautionary bool
	// TieStart and TieStop tie the note to the next and previous one.
	TieStart, TieStop bool
	Head              Notehead
}

// Pitch is a written pitch.
type Pitch struct {
	Step   int // 0 C … 6 B
	Alter  int // -2 … 2 semitones
	Octave int // 4 holds middle C
}

// MIDI returns the MIDI note number of the pitch.
func (p Pitch) MIDI() int {
	return (p.Octave+1)*12 + stepSemitones[(p.Step%7+7)%7] + p.Alter
}

var stepSemitones = [7]int{0, 2, 4, 5, 7, 9, 11}

// Accidental is an accidental sign.
type Accidental int

const (
	AccAuto Accidental = iota
	AccNone
	AccSharp
	AccFlat
	AccNatural
	AccDoubleSharp
	AccDoubleFlat
)

// Notehead is the shape of a notehead.
type Notehead int

const (
	HeadNormal Notehead = iota
	HeadX
)

// LyricSyllable is one syllable of a verse of lyrics.
type LyricSyllable struct {
	Verse int // 1-based
	Text  string
	// Syllabic places the syllable in its word: hyphens are drawn to the
	// next syllable after SyllableBegin and SyllableMiddle.
	Syllabic Syllabic
	// Extend draws an extender line to the end of the melisma.
	Extend bool
}

// Syllabic is the place of a syllable in its word.
type Syllabic int

const (
	SyllableSingle Syllabic = iota
	SyllableBegin
	SyllableMiddle
	SyllableEnd
)

// SlurMark starts or stops a slur on an event. Number tells overlapping
// slurs apart.
type SlurMark struct {
	Number int
	Start  bool // false: stop
	Above  *bool
}

// Direction is a mark at an offset of a measure: a dynamic, words, a
// hairpin, a pedal mark, an octave line or a chord symbol.
type Direction struct {
	Offset int
	Staff  int
	Kind   DirectionKind
	// Text is the dynamic ("mf", "sfz"), the words, or the chord symbol
	// ("Cmaj7").
	Text string
	// Above places the direction above the staff; the default depends on
	// the kind (dynamics and hairpins below, words and chord symbols above).
	Above *bool
	// Italic words ("dolce", "cresc.").
	Italic bool
	// Stop ends a hairpin, pedal or octave line that an earlier direction
	// with the same Number started.
	Stop   bool
	Number int
	// Size is the octave line's shift (8 or 15; negative below).
	Size int
}

// DirectionKind is the kind of a Direction.
type DirectionKind int

const (
	DirDynamic DirectionKind = iota
	DirWords
	DirCrescendo  // hairpin opening
	DirDiminuendo // hairpin closing
	DirPedal
	DirOctave
	DirChord
	DirSegno
	DirCoda
)
