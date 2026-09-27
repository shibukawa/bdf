package music

import "testing"

func TestParams(t *testing.T) {
	for s, want := range map[string]TimeSig{
		"3/4":     {Beats: 3, BeatType: 4},
		" 6/8 ":   {Beats: 6, BeatType: 8},
		"12/16":   {Beats: 12, BeatType: 16},
		"C":       {Beats: 4, BeatType: 4, Symbol: "common"},
		"cut":     {Beats: 2, BeatType: 2, Symbol: "cut"},
		"3/5":     {},
		"0/4":     {},
		"3":       {},
		"waltz":   {},
		"3/4/4":   {},
		"-3/4":    {},
		"100/4":   {},
		"4/128":   {},
		"4/0":     {},
		"":        {},
		"c|":      {Beats: 2, BeatType: 2, Symbol: "cut"},
		"common":  {Beats: 4, BeatType: 4, Symbol: "common"},
		" 6 / 8 ": {Beats: 6, BeatType: 8},
		"4/3":     {},
		"a/b":     {},
	} {
		got, err := ParseTime(s)
		if (err != nil) != (want == TimeSig{}) || got != want {
			t.Errorf("ParseTime(%q) = %v, %v", s, got, err)
		}
	}
	for s, want := range map[string]*KeySig{
		"C":        {},
		"G":        {Fifths: 1},
		"Bb":       {Fifths: -2},
		"b":        {Fifths: 5},
		"Bm":       {Fifths: 2, Minor: true},
		"bbm":      {Fifths: -5, Minor: true},
		"F#m":      {Fifths: 3, Minor: true},
		"E♭ minor": {Fifths: -6, Minor: true},
		"C#":       {Fifths: 7},
		"Cb":       {Fifths: -7},
		"A minor":  {Minor: true},
		"Am":       {Minor: true},
		"-3":       {Fifths: -3},
		"2m":       {Fifths: 2, Minor: true},
		"Gmaj":     {Fifths: 1},
		"Dmaj":     {Fifths: 2},
		"C minor":  {Fifths: -3, Minor: true},
		"E♭ major": {Fifths: -3},
		"Em":       {Fifths: 1, Minor: true},
		"Fbm":      nil,
		"C##":      nil,
		"8":        nil,
		"m":        nil,
		"G#":       nil,
		"Fb":       nil,
		"H":        nil,
		"":         nil,
		"Cx":       nil,
		"#":        nil,
	} {
		got, err := ParseKey(s)
		if want == nil {
			if err == nil {
				t.Errorf("ParseKey(%q) = %v, no error", s, got)
			}
			continue
		}
		if err != nil || got != *want {
			t.Errorf("ParseKey(%q) = %v, %v", s, got, err)
		}
	}
}
