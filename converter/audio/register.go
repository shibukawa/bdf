package audio

import (
	"io"

	"github.com/shibukawa/bdf/converter"
)

func init() {
	converter.Register(&converter.Format{
		Name:        "audio",
		Description: "Audio file: MP3, AAC (M4A or ADTS), FLAC, Ogg Vorbis or Opus, WAV or AIFF (the cover art and the tags)",
		Extensions:  []string{".mp3", ".m4a", ".m4b", ".m4p", ".aac", ".flac", ".ogg", ".oga", ".opus", ".spx", ".wav", ".wave", ".aif", ".aiff", ".aifc"},
		Detect:      Detect,
		Convert: func(r io.ReaderAt, size int64, o *converter.Options) (*converter.Result, error) {
			res, err := Convert(r, size, &Options{Title: o.Title, FileName: o.FileName, Images: o.Images, NoTextIndex: o.NoTextIndex, Warn: o.Warn})
			if err != nil {
				return nil, err
			}
			return &converter.Result{Doc: res.Doc, Warnings: res.Warnings, Summary: res.Summary()}, nil
		},
	})
}

// Detect reports whether an input is an audio file: by the ID3 tag or
// the signature it starts with, by its first frames when it starts with
// MPEG audio or ADTS frames, and for an MP4 file by its having audio
// tracks only.
func Detect(head []byte, r io.ReaderAt, size int64) bool {
	switch {
	case id3v2Size(head) > 0:
		return true
	case len(head) >= 4 && string(head[:4]) == "fLaC":
		return true
	case len(head) >= 4 && string(head[:4]) == "OggS":
		return isOggAudio(head)
	case len(head) >= 12 && (string(head[:4]) == "RIFF" || string(head[:4]) == "RF64") && string(head[8:12]) == "WAVE":
		return true
	case len(head) >= 12 && string(head[:4]) == "FORM" && (string(head[8:12]) == "AIFF" || string(head[8:12]) == "AIFC"):
		return true
	case len(head) >= 12 && string(head[4:8]) == "ftyp":
		return mp4Kind(head, r, size) == "audio"
	}
	return isMPEG(r, 0, size) || isADTS(r, 0, size)
}
