// Package midi reads Standard MIDI Files: the .mid (.midi, .smf) files of
// sequencers and notation programs, karaoke files (.kar) and RIFF MIDI
// files (.rmi). The format registered as "midi" writes their music as staff
// notation with converter/internal/music and stores the file for the viewer
// to play (docs/spec.md §4.4).
//
// Files of formats 0, 1 and 2 are read; the patterns of a format 2 file are
// taken to start together, as the tracks of format 1 do. Each channel of
// each track chunk that plays notes becomes a track of the performance,
// named after the chunk's track name, the channel's name (a name after a
// channel prefix) or its General MIDI program. A note-off (or a note-on of
// velocity 0) ends the earliest note of its key and channel that is still
// sounding; "all notes off" ends them all, and notes left sounding end with
// their track. Running status is followed, also across meta and system
// exclusive events. Tempo, time signature and key signature changes, lyrics
// (lyric events, and the text events of karaoke files, given to the track
// they are sung on), the title (the name of a format 0 track or of a
// conductor track, @T of a karaoke file, INAM of a RIFF MIDI file), the
// copyright and credits written as "Composer: …" are read; program changes,
// volume, pan, expression, sustain, CC 70 and pitch bend are kept as the
// controls of the tracks. Text in UTF-8, Shift_JIS or Windows-1252 is
// recognized.
//
// A file timed in SMPTE frames is notated at 120 BPM and played from a file
// written from its notes. A damaged file (truncated, with a wrong track
// count, a format 0 file of several tracks, bytes that are not chunks) is
// read as far as it goes and stored repaired. System exclusive messages,
// markers, cue points and the other controllers stay in the stored file but
// are not notated.
package midi

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"slices"
	"sort"

	"github.com/shibukawa/bdf/converter/internal/music"
)

// Options controls how a file is read.
type Options struct {
	// Warn receives non-fatal problems; when nil they are collected in
	// Parsed.Warnings.
	Warn func(msg string)
}

// Parsed is a Standard MIDI File read by Parse.
type Parsed struct {
	// Perf is the music: one track for each channel of each track chunk
	// that plays notes, in music.PPQ ticks. Its Division is the ticks per
	// quarter note of SMF (0 when SMF is nil).
	Perf *music.Performance
	// SMF is the Standard MIDI File to store for playback: the input (the
	// file inside a RIFF MIDI wrapper; it shares memory with the input), or
	// a repaired copy of a damaged file, or of a format 2 file as format 1.
	// It is nil for a file timed in SMPTE frames: a file written from Perf
	// is played instead.
	SMF []byte
	// Format is the format of the file (0, 1 or 2).
	Format int
	// Warnings are the non-fatal problems when Options.Warn is nil.
	Warnings []string
}

// Budgets that keep hostile files from running away.
const (
	maxChunks   = 1024 // track chunks
	maxNotes    = 2_000_000
	maxControls = 2_000_000 // controls kept; more are left out
	maxEvents   = 1_000_000 // tempo, time and key changes and texts kept
	maxQuarters = 1_000_000 // the length of a track in quarter notes
)

// Parse reads a Standard MIDI File (or a RIFF MIDI file).
func Parse(data []byte, o *Options) (*Parsed, error) {
	if o == nil {
		o = &Options{}
	}
	p := &parser{o: o, warned: map[string]bool{}}
	res, err := p.parse(data)
	if err != nil {
		return nil, err
	}
	res.Warnings = p.warnings
	return res, nil
}

type parser struct {
	o        *Options
	warnings []string
	warned   map[string]bool

	tm      timing
	maxTick int64
	// counts of what was kept, against the budgets
	notes, controls, events int
	// pending holds the notes sounding in the chunk being read, by channel
	// and key, earliest first.
	pending [16][128][]pendingNote
}

func (p *parser) warnOnce(key, format string, args ...any) {
	if p.warned[key] {
		return
	}
	p.warned[key] = true
	msg := fmt.Sprintf(format, args...)
	if p.o.Warn != nil {
		p.o.Warn(msg)
		return
	}
	p.warnings = append(p.warnings, msg)
}

// timing converts the ticks of a file (Num/Den of them in a quarter note)
// to PPQ ticks.
type timing struct{ num, den int64 }

// ppq rounds an absolute tick of the file to PPQ ticks; ticks are bounded
// by maxTick, so the products do not overflow.
func (t timing) ppq(tick int64) int {
	return int((2*tick*music.PPQ*t.den + t.num) / (2 * t.num))
}

type pendingNote struct {
	tick int64
	vel  int
}

type rawNote struct {
	on, off  int64
	key, vel int
}

type rawControl struct {
	tick  int64
	kind  music.ControlKind
	num   int
	value int
}

type rawText struct {
	tick int64
	typ  byte // meta event type: 1 text, 5 lyric
	data []byte
}

type change[T comparable] struct {
	tick int64
	v    T
}

// chunk is what a track chunk holds.
type chunk struct {
	index int
	body  []byte
	// parsed is the length of the body read as whole events, up to the end
	// of track; sawEnd is set when it ends with one.
	parsed int
	sawEnd bool
	// damaged is set when the body could not be read to its end.
	damaged bool
	// end is the tick the track ends at.
	end int64

	notes    [16][]rawNote
	controls [16][]rawControl

	name, instrument, copyright []byte
	chanNames                   [16][]byte // names given after a channel prefix
	texts                       []rawText
	tempo                       []change[float64]
	times                       []change[music.TimeSig]
	keys                        []change[music.KeySig]
}

func (c *chunk) hasNotes() bool {
	for _, n := range c.notes {
		if len(n) > 0 {
			return true
		}
	}
	return false
}

// riffInfo is what the INFO list of a RIFF MIDI file says.
type riffInfo struct{ title, copyright []byte }

func (p *parser) parse(data []byte) (*Parsed, error) {
	smf, info, err := unwrapRMID(data)
	if err != nil {
		return nil, err
	}
	if len(smf) < 8 || string(smf[:4]) != "MThd" {
		return nil, errors.New("midi: not a Standard MIDI File")
	}
	hl := int64(binary.BigEndian.Uint32(smf[4:]))
	if hl < 6 {
		return nil, fmt.Errorf("midi: the header chunk is too short (%d bytes)", hl)
	}
	if hl > int64(len(smf)-8) {
		return nil, errors.New("midi: the file is truncated in its header")
	}
	format := int(binary.BigEndian.Uint16(smf[8:]))
	declared := int(binary.BigEndian.Uint16(smf[10:]))
	division := binary.BigEndian.Uint16(smf[12:])

	smpte := division&0x8000 != 0
	switch {
	case smpte:
		fps, tpf := -int64(int8(division>>8)), int64(division&0xFF)
		if fps != 24 && fps != 25 && fps != 29 && fps != 30 || tpf == 0 {
			return nil, fmt.Errorf("midi: the SMPTE division %#04x is invalid", division)
		}
		// ticks per second, taken at 120 BPM: two quarter notes a second
		p.tm = timing{fps * tpf, 2}
		if fps == 29 { // 30 drop-frame, 29.97 frames a second
			p.tm = timing{2997 * tpf, 200}
		}
		p.warnOnce("smpte", "the file is timed in SMPTE frames (%d fps × %d): it is notated at 120 BPM and played from a file written from its notes", fps, tpf)
	case division == 0:
		return nil, errors.New("midi: the division of the header is 0")
	default:
		p.tm = timing{int64(division), 1}
	}
	p.maxTick = maxQuarters * p.tm.num / p.tm.den

	outFormat := format
	clean := !smpte
	switch {
	case format == 2:
		p.warnOnce("format2", "the patterns of a format 2 file are read as tracks that start together")
		outFormat, clean = 1, false
	case format > 2:
		p.warnOnce("format", "format %d is read as format 1", format)
		outFormat, clean = 1, false
	}

	var chunks []*chunk
	pos := int64(8 + hl)
	for pos < int64(len(smf)) {
		rest := smf[pos:]
		if len(rest) < 8 || !isChunkID(rest[:4]) {
			// not a chunk: look for the next track chunk
			next := bytes.Index(rest[min(len(rest), 1):], []byte("MTrk"))
			if next < 0 {
				if slices.ContainsFunc(rest, func(b byte) bool { return b != 0 }) {
					p.warnOnce("trailing", "%d bytes after the last track are ignored", len(rest))
				}
				clean = false
				break
			}
			p.warnOnce("skipped", "%d bytes that are not a chunk are skipped", next+1)
			pos += int64(next + 1)
			clean = false
			continue
		}
		n := int64(binary.BigEndian.Uint32(rest[4:]))
		end := pos + 8 + n
		if end > int64(len(smf)) {
			end = int64(len(smf))
			clean = false
		}
		if string(rest[:4]) == "MTrk" {
			if len(chunks) == maxChunks {
				return nil, fmt.Errorf("midi: the file has more than %d tracks", maxChunks)
			}
			c := &chunk{index: len(chunks), body: smf[pos+8 : end]}
			if err := p.readTrack(c); err != nil {
				return nil, err
			}
			if c.damaged || !c.sawEnd {
				clean = false
			}
			chunks = append(chunks, c)
		}
		pos = end
	}
	if len(chunks) == 0 {
		return nil, errors.New("midi: the file has no tracks")
	}
	if declared != len(chunks) {
		p.warnOnce("count", "the header counts %d tracks, the file has %d", declared, len(chunks))
		clean = false
	}
	if outFormat == 0 && len(chunks) > 1 {
		p.warnOnce("format0", "a format 0 file of %d tracks is read as format 1", len(chunks))
		outFormat, clean = 1, false
	}

	perf := p.perform(chunks, format, info)
	res := &Parsed{Perf: perf, Format: format}
	switch {
	case smpte:
		perf.Tempo = []music.Tempo{{Tick: 0, BPM: 120}}
	case clean:
		res.SMF = smf
		perf.Division = int(division)
	default:
		res.SMF = rebuild(outFormat, division, chunks)
		perf.Division = int(division)
	}
	return res, nil
}

// unwrapRMID returns the Standard MIDI File in a RIFF MIDI file (and what
// its INFO list says), or data itself.
func unwrapRMID(data []byte) ([]byte, riffInfo, error) {
	var info riffInfo
	if len(data) < 12 || string(data[:4]) != "RIFF" || string(data[8:12]) != "RMID" {
		return data, info, nil
	}
	var smf []byte
	for pos := int64(12); pos+8 <= int64(len(data)); {
		id := string(data[pos : pos+4])
		n := int64(binary.LittleEndian.Uint32(data[pos+4:]))
		body := data[pos+8 : min(pos+8+n, int64(len(data)))]
		switch id {
		case "data":
			if smf == nil {
				smf = body
			}
		case "LIST":
			if len(body) >= 4 && string(body[:4]) == "INFO" {
				for q := int64(4); q+8 <= int64(len(body)); {
					m := int64(binary.LittleEndian.Uint32(body[q+4:]))
					s := body[q+8 : min(q+8+m, int64(len(body)))]
					switch string(body[q : q+4]) {
					case "INAM":
						info.title = s
					case "ICOP":
						info.copyright = s
					}
					q += 8 + m + m&1
				}
			}
		}
		pos += 8 + n + n&1
	}
	if smf == nil {
		return nil, info, errors.New("midi: the RIFF MIDI file has no data chunk")
	}
	return smf, info, nil
}

// isChunkID tells whether b looks like the type of a chunk.
func isChunkID(b []byte) bool {
	for _, c := range b {
		if c < 0x20 || c > 0x7E {
			return false
		}
	}
	return true
}

// readVarlen reads a variable-length quantity of at most 4 bytes. n is 0
// when b ends inside it and -1 when it is longer.
func readVarlen(b []byte) (v, n int) {
	for i := 0; i < 4; i++ {
		if i >= len(b) {
			return 0, 0
		}
		v = v<<7 | int(b[i]&0x7F)
		if b[i] < 0x80 {
			return v, i + 1
		}
	}
	return 0, -1
}

// readTrack reads the events of a track chunk.
func (p *parser) readTrack(c *chunk) error {
	b := c.body
	var (
		tick    int64
		running byte
		prefix  = -1 // the channel of a channel prefix in effect
	)
	cut := func(why string) {
		c.damaged = true
		if why == "" {
			p.warnOnce("truncated", "track %d is cut short; it is read up to the cut", c.index+1)
		} else {
			p.warnOnce("damaged", "track %d is damaged (%s); it is read up to the damage", c.index+1, why)
		}
	}
loop:
	for i := 0; i < len(b); {
		delta, n := readVarlen(b[i:])
		if n <= 0 {
			if n == 0 {
				cut("")
			} else {
				cut("a delta time longer than 4 bytes")
			}
			break
		}
		i += n
		tick += int64(delta)
		if tick > p.maxTick {
			c.damaged = true
			p.warnOnce("long", "tracks longer than %d quarter notes are cut", maxQuarters)
			break
		}
		if i >= len(b) {
			cut("")
			break
		}
		st := b[i]
		switch {
		case st >= 0x80:
			i++
		case running != 0:
			st = running
		default:
			cut("a data byte without a status")
			break loop
		}
		switch {
		case st < 0xF0:
			n := 2
			if st&0xE0 == 0xC0 { // program change, channel pressure
				n = 1
			}
			if len(b)-i < n {
				cut("")
				break loop
			}
			running = st
			d1, d2 := b[i], byte(0)
			if n == 2 {
				d2 = b[i+1]
			}
			i += n
			if (d1|d2)&0x80 != 0 {
				p.warnOnce("data", "data bytes out of range are read without their top bit")
				d1, d2 = d1&0x7F, d2&0x7F
			}
			prefix = -1
			if err := p.channelEvent(c, tick, st, int(d1), int(d2)); err != nil {
				return err
			}
		case st == 0xFF:
			if i >= len(b) {
				cut("")
				break loop
			}
			typ := b[i]
			l, n := readVarlen(b[i+1:])
			if n <= 0 || l > len(b)-(i+1+n) {
				cut("")
				break loop
			}
			data := b[i+1+n : i+1+n+l]
			i += 1 + n + l
			if typ == 0x2F {
				c.end, c.sawEnd, c.parsed = tick, true, i
				break loop
			}
			p.meta(c, tick, typ, data, &prefix)
		case st == 0xF0 || st == 0xF7: // system exclusive
			l, n := readVarlen(b[i:])
			if n <= 0 || l > len(b)-(i+n) {
				cut("")
				break loop
			}
			i += n + l
		default:
			// system common and real-time messages have no place in a file
			n := 0
			switch st {
			case 0xF1, 0xF3:
				n = 1
			case 0xF2:
				n = 2
			}
			if len(b)-i < n {
				cut("")
				break loop
			}
			i += n
			p.warnOnce("system", "system messages in tracks are ignored")
		}
		c.end, c.parsed = tick, i
	}
	// notes left sounding end with the track
	for ch := range p.pending {
		for key, q := range p.pending[ch] {
			if len(q) > 0 {
				p.warnOnce("unended", "notes without a note-off end with their track")
				for _, pn := range q {
					c.notes[ch] = append(c.notes[ch], rawNote{on: pn.tick, off: c.end, key: key, vel: pn.vel})
				}
				p.pending[ch][key] = nil
			}
		}
	}
	return nil
}

// channelEvent reads a channel message.
func (p *parser) channelEvent(c *chunk, tick int64, st byte, d1, d2 int) error {
	ch := int(st & 15)
	switch st & 0xF0 {
	case 0x90:
		if d2 > 0 {
			if p.notes++; p.notes > maxNotes {
				return fmt.Errorf("midi: the file has more than %d notes", maxNotes)
			}
			p.pending[ch][d1] = append(p.pending[ch][d1], pendingNote{tick, d2})
			break
		}
		fallthrough // a note-on of velocity 0 is a note-off
	case 0x80:
		if q := p.pending[ch][d1]; len(q) > 0 {
			c.notes[ch] = append(c.notes[ch], rawNote{on: q[0].tick, off: tick, key: d1, vel: q[0].vel})
			p.pending[ch][d1] = q[1:]
		}
	case 0xB0:
		switch d1 {
		case 120, 123: // all sound off, all notes off
			for key, q := range p.pending[ch] {
				for _, pn := range q {
					c.notes[ch] = append(c.notes[ch], rawNote{on: pn.tick, off: tick, key: key, vel: pn.vel})
				}
				p.pending[ch][key] = nil
			}
		case 7, 10, 11, 64, 70: // volume, pan, expression, sustain, sound variation
			p.control(c, ch, rawControl{tick, music.ControlChange, d1, d2})
		}
	case 0xC0:
		p.control(c, ch, rawControl{tick, music.ControlProgram, 0, d1})
	case 0xE0:
		p.control(c, ch, rawControl{tick, music.ControlPitchBend, 0, d2<<7 | d1 - 8192})
	}
	return nil
}

func (p *parser) control(c *chunk, ch int, e rawControl) {
	if p.controls++; p.controls > maxControls {
		p.warnOnce("controls", "controls past the first %d are left out", maxControls)
		return
	}
	c.controls[ch] = append(c.controls[ch], e)
}

// meta reads a meta event.
func (p *parser) meta(c *chunk, tick int64, typ byte, data []byte, prefix *int) {
	switch typ {
	case 0x01, 0x03, 0x04, 0x05, 0x51, 0x58, 0x59:
		if p.events++; p.events > maxEvents {
			p.warnOnce("events", "meta events past the first %d are left out", maxEvents)
			return
		}
	}
	switch typ {
	case 0x01, 0x05: // text, lyric
		c.texts = append(c.texts, rawText{tick, typ, data})
	case 0x02: // copyright
		if len(c.copyright) == 0 {
			c.copyright = data
		}
	case 0x03, 0x04: // track name, instrument name
		switch {
		case *prefix >= 0:
			if len(c.chanNames[*prefix]) == 0 {
				c.chanNames[*prefix] = data
			}
		case typ == 0x03 && len(c.name) == 0:
			c.name = data
		case typ == 0x04 && len(c.instrument) == 0:
			c.instrument = data
		}
	case 0x20: // channel prefix
		if len(data) >= 1 {
			*prefix = int(data[0] & 15)
		}
	case 0x51: // tempo
		if len(data) >= 3 {
			if us := float64(int(data[0])<<16 | int(data[1])<<8 | int(data[2])); us > 0 {
				// a whole BPM is written rounded to microseconds
				bpm := 60e6 / us
				if r := math.Round(bpm); math.Round(60e6/r) == us {
					bpm = r
				}
				c.tempo = append(c.tempo, change[float64]{tick, bpm})
			}
		}
	case 0x58: // time signature
		if len(data) >= 2 && data[0] > 0 && data[1] <= 6 {
			c.times = append(c.times, change[music.TimeSig]{tick, music.TimeSig{Beats: int(data[0]), BeatType: 1 << data[1]}})
		} else {
			p.warnOnce("time", "invalid time signatures are ignored")
		}
	case 0x59: // key signature
		if len(data) >= 2 && int8(data[0]) >= -7 && int8(data[0]) <= 7 && data[1] <= 1 {
			c.keys = append(c.keys, change[music.KeySig]{tick, music.KeySig{Fifths: int(int8(data[0])), Minor: data[1] == 1}})
		} else {
			p.warnOnce("key", "invalid key signatures are ignored")
		}
	}
}

// rebuild writes the track chunks read as a Standard MIDI File of a
// format, each up to the damage or its end of track, ended when it has
// none.
func rebuild(format int, division uint16, chunks []*chunk) []byte {
	out := []byte("MThd\x00\x00\x00\x06")
	out = binary.BigEndian.AppendUint16(out, uint16(format))
	out = binary.BigEndian.AppendUint16(out, uint16(len(chunks)))
	out = binary.BigEndian.AppendUint16(out, division)
	for _, c := range chunks {
		body := c.body[:c.parsed]
		n := len(body)
		if !c.sawEnd {
			n += 4
		}
		out = append(out, "MTrk"...)
		out = binary.BigEndian.AppendUint32(out, uint32(n))
		out = append(out, body...)
		if !c.sawEnd {
			out = append(out, 0, 0xFF, 0x2F, 0)
		}
	}
	return out
}

// perform turns the chunks into a performance.
func (p *parser) perform(chunks []*chunk, format int, info riffInfo) *music.Performance {
	dec := newDecoder(chunks, info)
	perf := &music.Performance{}

	// the program, volume and pan each channel has, by tick, over all
	// chunks (a later chunk wins at a tick): for tracks that set none
	// before their first note
	var global [16][3][]change[int]
	for _, c := range chunks {
		for ch := range c.controls {
			for _, e := range c.controls[ch] {
				if k := initialKind(e); k >= 0 {
					global[ch][k] = append(global[ch][k], change[int]{e.tick, e.value})
				}
			}
		}
	}
	for ch := range global {
		for k := range global[ch] {
			sort.SliceStable(global[ch][k], func(i, j int) bool { return global[ch][k][i].tick < global[ch][k][j].tick })
		}
	}

	var owners []int // the chunk of each track
	zero := false
	for _, c := range chunks {
		var chans []int
		for ch := range c.notes {
			if len(c.notes[ch]) > 0 {
				chans = append(chans, ch)
			}
		}
		for _, ch := range chans {
			raw := c.notes[ch]
			slices.SortStableFunc(raw, func(a, b rawNote) int {
				if a.on != b.on {
					return cmpInt64(a.on, b.on)
				}
				if a.key != b.key {
					return a.key - b.key
				}
				return cmpInt64(a.off, b.off)
			})
			t := &music.Track{Channel: ch, Volume: -1, Pan: -1, Notes: make([]music.PlayNote, 0, len(raw))}
			for _, n := range raw {
				tick := p.tm.ppq(n.on)
				dur := p.tm.ppq(n.off) - tick
				if n.off == n.on && ch != 9 {
					zero = true
					continue
				}
				t.Notes = append(t.Notes, music.PlayNote{Tick: tick, Dur: max(dur, 1), Key: n.key, Vel: n.vel})
			}
			if len(t.Notes) == 0 {
				continue
			}
			// the initial values: those set at or before the first note,
			// in the chunk or else in any
			var init [3]int
			var has [3]bool
			first := raw[0].on
			for _, e := range c.controls[ch] {
				if e.tick > first {
					break
				}
				if k := initialKind(e); k >= 0 {
					init[k], has[k] = e.value, true
				}
			}
			for k := range init {
				if has[k] {
					continue
				}
				l := global[ch][k]
				if i := sort.Search(len(l), func(i int) bool { return l[i].tick > first }) - 1; i >= 0 {
					init[k], has[k] = l[i].v, true
				}
			}
			if has[0] && ch != 9 {
				t.Program = init[0]
			}
			if has[1] {
				t.Volume = init[1]
			}
			if has[2] {
				t.Pan = init[2]
			}
			if len(c.controls[ch]) > 0 {
				t.Controls = make([]music.Control, len(c.controls[ch]))
				for i, e := range c.controls[ch] {
					t.Controls[i] = music.Control{Tick: p.tm.ppq(e.tick), Kind: e.kind, Num: e.num, Value: e.value}
				}
			}
			t.Name = dec.trackName(c, ch, len(chans), has[0], t.Program, format == 0 && c.index == 0)
			perf.Tracks = append(perf.Tracks, t)
			owners = append(owners, c.index)
		}
	}
	if zero {
		p.warnOnce("zero", "notes of no length are left out")
	}

	// the end: that of the longest track, at least that of the last note
	for _, c := range chunks {
		perf.End = max(perf.End, p.tm.ppq(c.end))
	}
	for _, t := range perf.Tracks {
		for _, n := range t.Notes {
			perf.End = max(perf.End, n.Tick+n.Dur)
		}
	}

	var tempo [][]change[float64]
	var times [][]change[music.TimeSig]
	var keys [][]change[music.KeySig]
	for _, c := range chunks {
		tempo = append(tempo, c.tempo)
		times = append(times, c.times)
		keys = append(keys, c.keys)
	}
	for _, e := range merge(p, tempo) {
		perf.Tempo = append(perf.Tempo, music.Tempo{Tick: int(e.tick), BPM: e.v})
	}
	for _, e := range merge(p, times) {
		perf.Time = append(perf.Time, music.TimeChange{Tick: int(e.tick), Time: e.v})
	}
	for _, e := range merge(p, keys) {
		perf.Key = append(perf.Key, music.KeyChange{Tick: int(e.tick), Key: e.v})
	}

	p.credits(perf, dec, chunks, format, info)
	p.lyrics(perf, dec, chunks, owners)
	return perf
}

// initialKind is the index of a control among the initial values of a
// track (program, volume, pan), or -1.
func initialKind(e rawControl) int {
	switch {
	case e.kind == music.ControlProgram:
		return 0
	case e.kind == music.ControlChange && e.num == 7:
		return 1
	case e.kind == music.ControlChange && e.num == 10:
		return 2
	}
	return -1
}

func cmpInt64(a, b int64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

// merge merges the changes of the chunks in PPQ ticks: one at a tick (that
// of the later chunk, the later event), without repeating the value in
// effect.
func merge[T comparable](p *parser, lists [][]change[T]) []change[T] {
	at := map[int64]T{}
	for _, l := range lists {
		for _, e := range l {
			at[int64(p.tm.ppq(e.tick))] = e.v
		}
	}
	ticks := make([]int64, 0, len(at))
	for t := range at {
		ticks = append(ticks, t)
	}
	slices.Sort(ticks)
	var out []change[T]
	for _, t := range ticks {
		if len(out) > 0 && out[len(out)-1].v == at[t] {
			continue
		}
		out = append(out, change[T]{t, at[t]})
	}
	return out
}
