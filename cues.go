package bdf

import "fmt"

// CueSystem is a rectangle of a page that the playing position moves
// through: a system of a score (docs/spec.md §4.4).
type CueSystem struct {
	Page       uint32 // page index in the view (0-based)
	X, Y, W, H float32
}

// Cue ties a time of a view's music to a place on its pages: at Tick (in
// the ticks of the Standard MIDI File) the playing position is at X on
// System.
type Cue struct {
	Tick   uint32
	System uint32
	X      float32
}

// Cues is the content of a cue index part.
type Cues struct {
	Systems []CueSystem
	Cues    []Cue // in the order of Tick
}

// EncodeCues serializes cues as a cue index part (type "idx").
func EncodeCues(c *Cues) []byte {
	var w buf
	w.bytes([]byte("BCUE"))
	w.u16(1)
	w.varuint(uint64(len(c.Systems)))
	for _, s := range c.Systems {
		w.varuint(uint64(s.Page))
		w.f32s(s.X, s.Y, s.W, s.H)
	}
	w.varuint(uint64(len(c.Cues)))
	var t uint32
	for _, q := range c.Cues {
		w.varuint(uint64(q.Tick - t))
		t = q.Tick
		w.varuint(uint64(q.System))
		w.f32(q.X)
	}
	return w.b
}

// DecodeCues parses a cue index part.
func DecodeCues(data []byte) (*Cues, error) {
	r := &reader{b: data}
	if string(r.bytes(4)) != "BCUE" {
		return nil, &FormatError{Msg: "not a cue index part"}
	}
	if v := r.u16(); v > 1 {
		return nil, &FormatError{Msg: fmt.Sprintf("unsupported cue index version %d", v)}
	}
	c := &Cues{}
	n := r.count(17, "bad system count")
	for i := 0; i < n && r.err == nil; i++ {
		c.Systems = append(c.Systems, CueSystem{Page: uint32(r.varuint()), X: r.f32(), Y: r.f32(), W: r.f32(), H: r.f32()})
	}
	n = r.count(6, "bad cue count")
	var t uint32
	for i := 0; i < n && r.err == nil; i++ {
		t += uint32(r.varuint())
		q := Cue{Tick: t, System: uint32(r.varuint()), X: r.f32()}
		if int(q.System) >= len(c.Systems) && r.err == nil {
			return nil, &FormatError{Msg: "cue refers to a missing system"}
		}
		c.Cues = append(c.Cues, q)
	}
	return c, r.err
}
