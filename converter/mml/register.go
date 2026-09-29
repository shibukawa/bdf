package mml

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/music"
)

func init() {
	conv.Register(&conv.Format{
		Name:        "mml",
		Description: "Music Macro Language (MML)",
		Extensions:  []string{".mml"},
		Params: []conv.Param{
			{Name: "dialect", Usage: "auto (default), generic (FlMML, MSX BASIC), mabinogi (MML@…;) or ppmck (NES MCK/PPMCK)"},
			{Name: "octave", Usage: "auto (default: the dialect's), normal (> raises the octave) or reverse (< raises it)"},
			{Name: "time", Usage: "the time signature of the score, e.g. 3/4 or 6/8 (default 4/4)"},
			{Name: "key", Usage: "the key signature, e.g. G, Bb, F#m or Em (default: estimated from the notes)"},
			{Name: "program", Usage: "the General MIDI program (0–127) of the tracks that set no tone (default: the dialect's)"},
			{Name: "tab", Usage: "add guitar TAB views (true by default; false to omit)"},
		},
		Detect: Detect,
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			tab, err := music.TabEnabled(o)
			if err != nil {
				return nil, err
			}
			opts, no, err := options(o)
			if err != nil {
				return nil, err
			}
			if size > maxFile {
				return nil, fmt.Errorf("mml: the file is larger than %d MiB", maxFile>>20)
			}
			b := make([]byte, size)
			if n, err := r.ReadAt(b, 0); n < len(b) {
				return nil, err
			}
			perf, warnings, err := Parse(b, opts)
			if err != nil {
				return nil, err
			}
			buildOptions := music.ConverterOptions("mml", o, o.Warn)
			buildOptions.GuitarTAB = tab
			res, err := music.Build(music.Notate(perf, no), buildOptions)
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: append(warnings, res.Warnings...), Summary: res.Summary()}, nil
		},
	})
}

// maxFile bounds the input read.
const maxFile = 16 << 20

// options reads the format's parameters.
func options(o *conv.Options) (*Options, music.NotateOptions, error) {
	opts := &Options{Warn: o.Warn}
	switch v := strings.ToLower(o.Param("dialect")); v {
	case "", "auto":
	case Generic, Mabinogi, PPMCK:
		opts.Dialect = v
	default:
		return nil, music.NotateOptions{}, fmt.Errorf("mml: parameter dialect: %q is not auto, generic, mabinogi or ppmck", v)
	}
	switch v := strings.ToLower(o.Param("octave")); v {
	case "", "auto":
	case "normal":
		opts.Octave = OctaveNormal
	case "reverse":
		opts.Octave = OctaveReverse
	default:
		return nil, music.NotateOptions{}, fmt.Errorf("mml: parameter octave: %q is not auto, normal or reverse", v)
	}
	if v := o.Param("program"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 || n > 127 {
			return nil, music.NotateOptions{}, fmt.Errorf("mml: parameter program: %q is not a number 0–127", v)
		}
		opts.Program = &n
	}
	no, err := music.NotateParams(o)
	if err != nil {
		return nil, no, fmt.Errorf("mml: %w", err)
	}
	return opts, no, nil
}

// Detect reports whether the first bytes of an input are MML: a line that
// starts with Mabinogi's MML@, header lines of MML in upper case (#TITLE,
// #COMPOSER, #OCTAVE-REV, #WAV9 …) with lines of MML or PPMCK definitions
// (@v0 = { … }), or two PPMCK track lines ("A t150 o4 l8 cdef"). Text that
// is only notes is not told from other text: it is MML by its extension.
func Detect(head []byte, r io.ReaderAt, size int64) bool {
	head = bytes.TrimPrefix(head, []byte{0xEF, 0xBB, 0xBF})
	if len(head) == 0 || bytes.IndexByte(head, 0) >= 0 {
		return false
	}
	var mabinogi, fence bool
	var headers, defs, tracks, lines int
	for _, l := range strings.Split(strings.ReplaceAll(string(head), "\r", "\n"), "\n") {
		t := strings.TrimSpace(l)
		switch {
		case t == "":
		case strings.HasPrefix(t, "```"), strings.HasPrefix(t, "~~~"): // a Markdown code block
			fence = true
		case len(t) > 4 && strings.EqualFold(t[:4], "MML@"):
			if ok, tokens, _ := mmlLike(t[4:min(len(t), 68)]); ok && tokens >= 1 {
				mabinogi = true
			}
		case t[0] == '#':
			if name, _, _ := strings.Cut(strings.ReplaceAll(t, "\t", " "), " "); name == strings.ToUpper(name) &&
				(name == "#TITLE" || headerDialect(name) != "") {
				headers++
			}
		case isDefinition(t):
			defs++
		case isTrackLine(l):
			tracks++
		default:
			if k := strings.Index(t, "//"); k >= 0 {
				t = t[:k]
			}
			if ok, tokens, notes := mmlLike(t); ok && tokens >= 1 && tokens+notes >= 3 {
				lines++
			}
		}
	}
	return !fence && (mabinogi || tracks >= 2 || headers > 0 && tracks+defs+lines > 0)
}
