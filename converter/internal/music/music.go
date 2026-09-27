// Package music engraves scores: it lays out music as staff notation on
// pages of a BDF document and stores what the viewer plays (a Standard
// MIDI File and the cues that tie its time to places on the pages,
// docs/spec.md §4.4).
//
// Music comes in two shapes. A Performance is music as played: the timed
// notes of each track, as MML and MIDI files hold it; Notate turns it
// into a Score, choosing measures, note values, ties, voices, clefs and
// spelling. A Score is music as written, as MusicXML holds it: parts of
// measures of notes with their note values, beams and marks. Build
// engraves a Score into a document.
//
// Positions and durations are in ticks, PPQ to a quarter note.
package music

// PPQ is the number of ticks in a quarter note.
const PPQ = 960

// Performance is music as played (MML, MIDI).
type Performance struct {
	Title, Subtitle, Composer, Lyricist, Arranger, Copyright string

	// Tempo, Time and Key are the changes of tempo, time signature and key
	// signature, in the order of Tick. Without Time the music is in 4/4;
	// without Key, Notate estimates the key from the notes.
	Tempo []Tempo
	Time  []TimeChange
	Key   []KeyChange

	Tracks []*Track

	// End is where the music ends (the end of the longest track), at least
	// the end of the last note.
	End int

	// Division is the ticks per quarter note of the Standard MIDI File
	// that is played instead of one written from the performance (a MIDI
	// file stored as it is): cue times are converted to it. 0 is PPQ.
	Division int
}

// Tempo is a tempo change.
type Tempo struct {
	Tick int
	BPM  float64 // quarter notes per minute
}

// TimeChange is a time signature change.
type TimeChange struct {
	Tick int
	Time TimeSig
}

// KeyChange is a key signature change.
type KeyChange struct {
	Tick int
	Key  KeySig
}

// Track is one line of a performance: an MML track or part, the notes of
// one channel of a MIDI track.
type Track struct {
	Name    string
	Channel int // 0–15; 9 is percussion
	Program int // General MIDI program at the start (0–127)
	// Volume is the initial channel volume (CC 7), or -1 for none.
	Volume int
	// Pan is the initial pan (CC 10), or -1 for none.
	Pan int

	Notes    []PlayNote // in the order of Tick
	Controls []Control  // in the order of Tick
	Lyrics   []Lyric    // in the order of Tick: syllables sung from the note starting at Tick
}

// PlayNote is a note as played.
type PlayNote struct {
	Tick, Dur int
	Key       int // MIDI note number (60 is middle C)
	Vel       int // velocity 1–127
	// Spell is how the note was written when the input says so (MML's
	// c+ and d-): 1 sharp, -1 flat, 0 unknown or natural.
	Spell int
}

// ControlKind is the kind of a Control.
type ControlKind int

const (
	// ControlProgram changes the program to Value.
	ControlProgram ControlKind = iota
	// ControlChange sets controller Num to Value.
	ControlChange
	// ControlPitchBend bends by Value (-8192–8191).
	ControlPitchBend
)

// Control is a program change or a controller of a track.
type Control struct {
	Tick  int
	Kind  ControlKind
	Num   int
	Value int
}

// Lyric is a syllable of a track's lyrics.
type Lyric struct {
	Tick int
	Text string
}
