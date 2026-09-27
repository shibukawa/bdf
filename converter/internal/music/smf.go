package music

import (
	"encoding/binary"
	"math"
	"sort"
)

// SMF writes the performance as a Standard MIDI File (format 1, PPQ ticks
// per quarter note): a conductor track with the title, tempo, time and key
// signatures, then one track per Track.
func (p *Performance) SMF() []byte {
	end := p.end()
	var out []byte
	out = append(out, "MThd"...)
	out = binary.BigEndian.AppendUint32(out, 6)
	out = binary.BigEndian.AppendUint16(out, 1)
	out = binary.BigEndian.AppendUint16(out, uint16(1+len(p.Tracks)))
	out = binary.BigEndian.AppendUint16(out, PPQ)

	// conductor track
	var c trackWriter
	if p.Title != "" {
		c.meta(0, 0x03, []byte(p.Title))
	}
	if p.Copyright != "" {
		c.meta(0, 0x02, []byte(p.Copyright))
	}
	type conductor struct {
		tick, order int
		data        []byte
	}
	var evs []conductor
	for _, t := range p.Time {
		dd := 0
		for 1<<dd < t.Time.BeatType {
			dd++
		}
		evs = append(evs, conductor{t.Tick, 0, []byte{0x58, 4, byte(t.Time.Beats), byte(dd), 24, 8}})
	}
	for _, k := range p.Key {
		mi := byte(0)
		if k.Key.Minor {
			mi = 1
		}
		evs = append(evs, conductor{k.Tick, 1, []byte{0x59, 2, byte(int8(k.Key.Fifths)), mi}})
	}
	for _, t := range p.Tempo {
		if t.BPM <= 0 {
			continue
		}
		us := int(math.Round(60e6 / t.BPM))
		us = min(max(us, 1), 0xFFFFFF)
		evs = append(evs, conductor{t.Tick, 2, []byte{0x51, 3, byte(us >> 16), byte(us >> 8), byte(us)}})
	}
	sort.SliceStable(evs, func(i, j int) bool {
		if evs[i].tick != evs[j].tick {
			return evs[i].tick < evs[j].tick
		}
		return evs[i].order < evs[j].order
	})
	for _, e := range evs {
		c.event(e.tick, append([]byte{0xFF}, e.data...))
	}
	c.meta(end, 0x2F, nil)
	out = c.chunk(out)

	for _, t := range p.Tracks {
		out = t.smf(out, end)
	}
	return out
}

// end is the tick the performance ends at.
func (p *Performance) end() int {
	end := p.End
	for _, t := range p.Tracks {
		for _, n := range t.Notes {
			end = max(end, n.Tick+n.Dur)
		}
	}
	return end
}

// smf appends the track chunk of a track.
func (t *Track) smf(out []byte, end int) []byte {
	type ev struct {
		tick, order int
		data        []byte
	}
	ch := byte(t.Channel & 15)
	var evs []ev
	if t.Program >= 0 && t.Channel != 9 {
		evs = append(evs, ev{0, 0, []byte{0xC0 | ch, byte(t.Program & 127)}})
	}
	if t.Volume >= 0 {
		evs = append(evs, ev{0, 0, []byte{0xB0 | ch, 7, byte(min(t.Volume, 127))}})
	}
	if t.Pan >= 0 {
		evs = append(evs, ev{0, 0, []byte{0xB0 | ch, 10, byte(min(t.Pan, 127))}})
	}
	for _, c := range t.Controls {
		switch c.Kind {
		case ControlProgram:
			evs = append(evs, ev{c.Tick, 1, []byte{0xC0 | ch, byte(c.Value & 127)}})
		case ControlChange:
			evs = append(evs, ev{c.Tick, 1, []byte{0xB0 | ch, byte(c.Num & 127), byte(min(max(c.Value, 0), 127))}})
		case ControlPitchBend:
			v := min(max(c.Value+8192, 0), 16383)
			evs = append(evs, ev{c.Tick, 1, []byte{0xE0 | ch, byte(v & 127), byte(v >> 7)}})
		}
	}
	for _, l := range t.Lyrics {
		evs = append(evs, ev{l.Tick, 2, append([]byte{0xFF, 0x05}, varlen(nil, len(l.Text))...)})
		evs[len(evs)-1].data = append(evs[len(evs)-1].data, l.Text...)
	}
	for _, n := range t.Notes {
		if n.Dur <= 0 || n.Key < 0 || n.Key > 127 {
			continue
		}
		vel := byte(min(max(n.Vel, 1), 127))
		// note-offs first at a tick, so that a repeated note is not cut
		evs = append(evs, ev{n.Tick + n.Dur, -1, []byte{0x80 | ch, byte(n.Key), 64}})
		evs = append(evs, ev{n.Tick, 3, []byte{0x90 | ch, byte(n.Key), vel}})
	}
	sort.SliceStable(evs, func(i, j int) bool {
		if evs[i].tick != evs[j].tick {
			return evs[i].tick < evs[j].tick
		}
		return evs[i].order < evs[j].order
	})
	var w trackWriter
	if t.Name != "" {
		w.meta(0, 0x03, []byte(t.Name))
	}
	for _, e := range evs {
		w.event(e.tick, e.data)
	}
	last := 0
	if len(evs) > 0 {
		last = evs[len(evs)-1].tick
	}
	w.meta(max(end, last), 0x2F, nil)
	return w.chunk(out)
}

// trackWriter writes the events of a track chunk with delta times and
// running status.
type trackWriter struct {
	b       []byte
	tick    int
	running byte
}

func (w *trackWriter) event(tick int, data []byte) {
	w.b = varlen(w.b, max(tick-w.tick, 0))
	w.tick = max(tick, w.tick)
	st := data[0]
	switch {
	case st >= 0xF0:
		w.running = 0
		w.b = append(w.b, data...)
	case st == w.running:
		w.b = append(w.b, data[1:]...)
	default:
		w.running = st
		w.b = append(w.b, data...)
	}
}

func (w *trackWriter) meta(tick int, typ byte, data []byte) {
	ev := append([]byte{0xFF, typ}, varlen(nil, len(data))...)
	w.event(tick, append(ev, data...))
}

func (w *trackWriter) chunk(out []byte) []byte {
	out = append(out, "MTrk"...)
	out = binary.BigEndian.AppendUint32(out, uint32(len(w.b)))
	return append(out, w.b...)
}

// varlen appends a MIDI variable-length quantity.
func varlen(b []byte, v int) []byte {
	var tmp [5]byte
	i := len(tmp) - 1
	tmp[i] = byte(v & 0x7F)
	for v >>= 7; v > 0; v >>= 7 {
		i--
		tmp[i] = byte(v&0x7F) | 0x80
	}
	return append(b, tmp[i:]...)
}
