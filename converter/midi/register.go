package midi

import (
	"encoding/binary"
	"fmt"
	"io"

	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/music"
)

// maxSize bounds the input read.
const maxSize = 64 << 20

func init() {
	conv.Register(&conv.Format{
		Name:        "midi",
		Description: "Standard MIDI File",
		Extensions:  []string{".mid", ".midi", ".smf", ".kar", ".rmi"},
		Params: []conv.Param{
			{Name: "time", Usage: "time signature instead of the file's (3/4, 6/8, C for common time, cut)"},
			{Name: "key", Usage: "key signature instead of the file's or the estimated one (G, Bb, F#m, or -2 for two flats)"},
		},
		Detect: func(head []byte, r io.ReaderAt, size int64) bool { return detect(head) },
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			no, err := music.NotateParams(o)
			if err != nil {
				return nil, err
			}
			if size > maxSize {
				return nil, fmt.Errorf("midi: the file is larger than %d bytes", maxSize)
			}
			data := make([]byte, size)
			if n, err := r.ReadAt(data, 0); n < len(data) {
				if err == nil {
					err = io.ErrUnexpectedEOF
				}
				return nil, fmt.Errorf("midi: %w", err)
			}
			var warnings []string
			warn := o.Warn
			if warn == nil {
				warn = func(msg string) { warnings = append(warnings, msg) }
			}
			parsed, err := Parse(data, &Options{Warn: warn})
			if err != nil {
				return nil, err
			}
			score := music.Notate(parsed.Perf, no)
			score.SMF = parsed.SMF
			res, err := music.Build(score, music.ConverterOptions("midi", o, warn))
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: append(warnings, res.Warnings...), Summary: res.Summary()}, nil
		},
	})
}

// detect tells a Standard MIDI File (or a RIFF MIDI file) by its header.
func detect(head []byte) bool {
	switch {
	case len(head) >= 8 && string(head[:4]) == "MThd":
		n := binary.BigEndian.Uint32(head[4:])
		return n >= 6 && n < 1<<16
	case len(head) >= 12 && string(head[:4]) == "RIFF" && string(head[8:12]) == "RMID":
		return true
	}
	return false
}
