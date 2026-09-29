package music

import (
	"fmt"
	"strconv"
	"strings"

	conv "github.com/shibukawa/bdf/converter"
)

// ConverterOptions returns the engraving options for a conversion with
// the converter options o, from the input format source; warn receives
// the warnings.
func ConverterOptions(source string, o *conv.Options, warn func(string)) Options {
	return Options{Source: source, Title: o.Title, Pages: o.Pages,
		FontFS: o.FontFS, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts,
		NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: warn}
}

// TabEnabled reads the optional TAB switch. Scores include TAB by default.
func TabEnabled(o *conv.Options) (bool, error) {
	if o.Param("tab") == "" {
		return true, nil
	}
	return o.BoolParam("tab")
}

// NotateParams reads the time and key parameters of a conversion (an
// input's own time and key signatures are replaced by them).
func NotateParams(o *conv.Options) (NotateOptions, error) {
	var no NotateOptions
	if v := o.Param("time"); v != "" {
		t, err := ParseTime(v)
		if err != nil {
			return no, err
		}
		no.Time = &t
	}
	if v := o.Param("key"); v != "" {
		k, err := ParseKey(v)
		if err != nil {
			return no, err
		}
		no.Key = &k
	}
	return no, nil
}

// ParseTime reads a time signature: 3/4, C (common time) or cut (alla
// breve).
func ParseTime(s string) (TimeSig, error) {
	switch v := strings.TrimSpace(s); strings.ToLower(v) {
	case "c", "common":
		return TimeSig{Beats: 4, BeatType: 4, Symbol: "common"}, nil
	case "c|", "cut", "alla breve":
		return TimeSig{Beats: 2, BeatType: 2, Symbol: "cut"}, nil
	default:
		n, d, ok := strings.Cut(v, "/")
		beats, err1 := strconv.Atoi(strings.TrimSpace(n))
		beatType, err2 := strconv.Atoi(strings.TrimSpace(d))
		if ok && err1 == nil && err2 == nil && beats >= 1 && beats <= 99 && beatType >= 1 && beatType <= 64 && beatType&(beatType-1) == 0 {
			return TimeSig{Beats: beats, BeatType: beatType}, nil
		}
	}
	return TimeSig{}, fmt.Errorf("parameter time: %q is not a time signature such as 3/4", s)
}

// ParseKey reads a key: a note name with # or b and m for minor (G, Bb,
// F#m, E♭ minor), or a number of sharps (positive) or flats (negative) with
// m for minor.
func ParseKey(s string) (KeySig, error) {
	fail := fmt.Errorf("parameter key: %q is not a key such as G, Bb or F#m", s)
	v := strings.TrimSpace(s)
	minor := false
	for _, suf := range []struct {
		s     string
		minor bool
	}{{"minor", true}, {"Minor", true}, {"min", true}, {"major", false}, {"Major", false}, {"maj", false}, {"m", true}, {"M", false}} {
		if rest, ok := strings.CutSuffix(v, suf.s); ok && rest != "" {
			v, minor = strings.TrimSpace(rest), suf.minor
			break
		}
	}
	fifths, err := strconv.Atoi(v)
	if err != nil {
		if v == "" {
			return KeySig{}, fail
		}
		fifths = strings.IndexByte("FCGDAEB", v[0]&^0x20) - 1 // F is -1, B is 5
		if fifths < -1 {
			return KeySig{}, fail
		}
		switch v[1:] {
		case "":
		case "#", "♯":
			fifths += 7
		case "b", "♭":
			fifths -= 7
		default:
			return KeySig{}, fail
		}
		if minor {
			fifths -= 3
		}
	}
	if fifths < -7 || fifths > 7 {
		return KeySig{}, fmt.Errorf("parameter key: %q has more than 7 sharps or flats", s)
	}
	return KeySig{Fifths: fifths, Minor: minor}, nil
}
