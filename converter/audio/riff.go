package audio

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
)

// Bounds on the chunks of WAV and AIFF files that are read whole.
const (
	maxInfoChunk = 16 << 20
	maxTextChunk = 1 << 20
)

// chunkReader walks the chunks of a RIFF (little-endian) or IFF
// (big-endian) file without reading their contents.
type chunkReader struct {
	r    io.ReaderAt
	size int64
	pos  int64
	be   bool
	// declared is the size the last chunk claimed, before it was cut to
	// the file.
	declared int64
}

// next returns the next chunk's identifier, the offset and size of its
// contents, and false at the end.
func (c *chunkReader) next() (id string, off, size int64, ok bool) {
	var hdr [8]byte
	if c.pos+8 > c.size {
		return "", 0, 0, false
	}
	if n, _ := c.r.ReadAt(hdr[:], c.pos); n < 8 {
		return "", 0, 0, false
	}
	id = string(hdr[:4])
	if c.be {
		size = int64(binary.BigEndian.Uint32(hdr[4:]))
	} else {
		size = int64(binary.LittleEndian.Uint32(hdr[4:]))
	}
	off = c.pos + 8
	c.declared = size
	if size > c.size-off {
		size = c.size - off // a truncated last chunk
	}
	c.pos = off + size + size&1
	return id, off, size, true
}

// read reads a chunk's contents, up to limit bytes.
func (c *chunkReader) read(off, size, limit int64) []byte {
	if size > limit {
		return nil
	}
	b := make([]byte, size)
	n, _ := c.r.ReadAt(b, off)
	return b[:n]
}

// readWAV reads a WAVE (or RF64) file: its format chunk, the size of its
// data, the INFO list, the broadcast extension and an ID3 chunk.
func readWAV(r io.ReaderAt, size int64) (*track, error) {
	t := newTrack("wav")
	c := &chunkReader{r: r, size: size, pos: 12}
	le := binary.LittleEndian
	var dataSize, ds64Data int64
	dataDeclared := int64(-1)
	var byteRate int
	var info [][]byte
	var bext []byte
	for {
		id, off, n, ok := c.next()
		if !ok {
			break
		}
		switch id {
		case "fmt ":
			b := c.read(off, n, 64<<10)
			if len(b) < 16 {
				return nil, errors.New("the format chunk is truncated")
			}
			tag := int(le.Uint16(b))
			t.channels = int(le.Uint16(b[2:]))
			t.sampleRate = int(le.Uint32(b[4:]))
			byteRate = int(le.Uint32(b[8:]))
			t.bits = int(le.Uint16(b[14:]))
			if tag == 0xFFFE && len(b) >= 26 { // extensible: the subformat's first bytes
				tag = int(le.Uint16(b[24:]))
			}
			t.codec = waveCodec(tag)
			if t.codec != "PCM" {
				t.bits = 0
			}
		case "data":
			dataSize, dataDeclared = n, c.declared
		case "ds64":
			if b := c.read(off, n, 64<<10); len(b) >= 16 {
				ds64Data = int64(le.Uint64(b[8:]))
			}
		case "LIST":
			if b := c.read(off, n, maxInfoChunk); len(b) >= 4 && string(b[:4]) == "INFO" {
				info = append(info, b[4:])
			}
		case "bext":
			bext = c.read(off, n, maxTextChunk)
		case "id3 ", "ID3 ":
			if b := c.read(off, n, maxID3); b != nil {
				if tag, err := parseID3v2(b); err == nil {
					t.applyID3v2(tag)
				}
			}
		}
	}
	if t.sampleRate == 0 {
		return nil, errors.New("the file has no format chunk")
	}
	if ds64Data > 0 && (dataDeclared == math.MaxUint32 || dataSize == 0) {
		dataSize = ds64Data
	}
	if byteRate > 0 && dataSize > 0 {
		t.duration = float64(dataSize) / float64(byteRate)
		t.bitrate = byteRate * 8
	}
	// the INFO list, then the broadcast extension, after an ID3 chunk
	native := tags{}
	var samples [][]byte
	type field struct{ key, value []byte }
	var fields []field
	for _, list := range info {
		for p := 0; p+8 <= len(list); {
			id := string(list[p : p+4])
			n := int(le.Uint32(list[p+4:]))
			p += 8
			if n < 0 || n > len(list)-p {
				break
			}
			fields = append(fields, field{[]byte(id), list[p : p+n]})
			samples = append(samples, list[p:p+n])
			p += n + n&1
		}
	}
	if len(bext) >= 256+32+32+10+8 {
		samples = append(samples, bext[:256], bext[256:288])
	}
	l := guessLegacy(samples)
	for _, f := range fields {
		v := l.decode(f.value)
		switch string(f.key) {
		case "INAM":
			native.add(keyTitle, v)
		case "ISBJ":
			native.add(keySubtitle, v)
		case "IART":
			native.add(keyArtist, splitPeople([]string{v})...)
		case "IPRD":
			native.add(keyAlbum, v)
		case "ICRD":
			native.add(keyDate, v)
		case "IGNR":
			native.add(keyGenre, v)
		case "ICMT":
			native.add(keyComment, v)
		case "ICOP":
			native.add(keyCopyright, v)
		case "ITRK", "IPRT":
			native.addNumber(keyTrack, keyTrackTotal, v)
		case "IMUS":
			native.add(keyComposer, splitPeople([]string{v})...)
		case "IWRI":
			native.add(keyLyricist, splitPeople([]string{v})...)
		case "ILNG":
			native.add(keyLanguage, v)
		}
	}
	if len(bext) >= 256+32+32+10+8 {
		if !native.has(keyComment) {
			native.add(keyComment, l.decode(bext[:256]))
		}
		if !native.has(keyDate) {
			date := strings.ReplaceAll(strings.Trim(string(bext[320:330]), "\x00 "), ":", "-")
			if tm := strings.Trim(string(bext[330:338]), "\x00 "); len(tm) == 8 {
				date += "T" + strings.ReplaceAll(tm, "-", ":")
			}
			native.add(keyDate, date)
		}
	}
	t.tags.fill(native)
	return t, nil
}

// waveCodec names a WAVE format tag.
func waveCodec(tag int) string {
	switch tag {
	case 1:
		return "PCM"
	case 3:
		return "PCM (float)"
	case 6:
		return "A-law"
	case 7:
		return "μ-law"
	case 2:
		return "MS ADPCM"
	case 0x11:
		return "IMA ADPCM"
	case 0x50:
		return "MPEG"
	case 0x55:
		return "MP3"
	case 0x2000:
		return "AC-3"
	case 0xFF, 0xA106:
		return "AAC"
	}
	return fmt.Sprintf("WAVE format 0x%X", tag)
}

// readAIFF reads an AIFF or AIFF-C file: its common chunk, the size of
// its sound data, its text chunks and an ID3 chunk.
func readAIFF(r io.ReaderAt, size int64) (*track, error) {
	t := newTrack("aiff")
	var form [12]byte
	if n, _ := r.ReadAt(form[:], 0); n < 12 {
		return nil, errors.New("the file is truncated")
	}
	aifc := string(form[8:12]) == "AIFC"
	c := &chunkReader{r: r, size: size, pos: 12, be: true}
	be := binary.BigEndian
	var frames uint32
	type text struct {
		key   string
		value []byte
	}
	var texts []text
	var samples [][]byte
	for {
		id, off, n, ok := c.next()
		if !ok {
			break
		}
		switch id {
		case "COMM":
			b := c.read(off, n, 64<<10)
			if len(b) < 18 {
				return nil, errors.New("the common chunk is truncated")
			}
			t.channels = int(be.Uint16(b))
			frames = be.Uint32(b[2:])
			t.bits = int(be.Uint16(b[6:]))
			t.sampleRate = int(extended(b[8:18]) + 0.5)
			t.codec = "PCM"
			if aifc && len(b) >= 22 {
				t.codec = aifcCodec(string(b[18:22]))
				if t.codec != "PCM" && t.codec != "PCM (float)" {
					t.bits = 0
				}
			}
		case "NAME", "AUTH", "(c) ", "ANNO":
			if b := c.read(off, n, maxTextChunk); b != nil {
				texts = append(texts, text{id, b})
				samples = append(samples, b)
			}
		case "ID3 ", "id3 ":
			if b := c.read(off, n, maxID3); b != nil {
				if tag, err := parseID3v2(b); err == nil {
					t.applyID3v2(tag)
				}
			}
		}
	}
	if t.sampleRate == 0 {
		return nil, errors.New("the file has no common chunk")
	}
	t.duration = float64(frames) / float64(t.sampleRate)
	if t.bits > 0 {
		t.bitrate = t.sampleRate * t.channels * t.bits
	}
	native := tags{}
	l := guessLegacy(samples)
	for _, x := range texts {
		v := l.decode(x.value)
		switch x.key {
		case "NAME":
			native.add(keyTitle, v)
		case "AUTH":
			native.add(keyArtist, splitPeople([]string{v})...)
		case "(c) ":
			native.add(keyCopyright, v)
		case "ANNO":
			native.add(keyComment, v)
		}
	}
	t.tags.fill(native)
	return t, nil
}

// aifcCodec names an AIFF-C compression type.
func aifcCodec(typ string) string {
	switch strings.ToUpper(typ) {
	case "NONE", "SOWT", "TWOS", "IN24", "IN32":
		return "PCM"
	case "FL32", "FL64":
		return "PCM (float)"
	case "ULAW":
		return "μ-law"
	case "ALAW":
		return "A-law"
	case "IMA4":
		return "IMA ADPCM"
	}
	return strings.TrimSpace(typ)
}

// extended reads an 80-bit IEEE 754 extended-precision number (the sample
// rate of AIFF).
func extended(b []byte) float64 {
	sign := 1.0
	if b[0]&0x80 != 0 {
		sign = -1
	}
	exp := int(binary.BigEndian.Uint16(b)&0x7FFF) - 16383
	mantissa := binary.BigEndian.Uint64(b[2:])
	if exp == -16383 && mantissa == 0 {
		return 0
	}
	return sign * math.Ldexp(float64(mantissa), exp-63)
}
