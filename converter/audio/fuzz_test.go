package audio

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// FuzzRead checks that no input makes the readers panic, or the converter.
func FuzzRead(f *testing.F) {
	files, _ := filepath.Glob(filepath.Join("testdata", "*"))
	for _, name := range files {
		if b, err := os.ReadFile(name); err == nil {
			f.Add(b)
		}
	}
	f.Add([]byte("ID3\x04\x00\x10\x00\x00\x00\x10TIT2\x00\x00\x00\x06\x00\x00\x03hello3DI\x04\x00\x10\x00\x00\x00\x10"))
	f.Add(oggPages(1, 0, vorbisIdent(44100, 2, 0), []byte("\x03vorbis")))
	f.Add(box("ftyp", []byte("M4A ")))
	f.Fuzz(func(t *testing.T, data []byte) {
		r := bytes.NewReader(data)
		Detect(data[:min(len(data), 1024)], r, int64(len(data)))
		if _, err := Convert(r, int64(len(data)), &Options{FileName: "fuzz.mp3"}); err != nil {
			return
		}
	})
}
