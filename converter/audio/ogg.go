package audio

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// Bounds on an Ogg file: how much is read for the headers of its stream,
// and how far from the end its last page is looked for (a page holds at
// most 65,307 bytes).
const (
	maxOggHeaders = 64 << 20
	maxOggPages   = 10000
	oggTail       = 256 << 10
)

// oggPage is a page of an Ogg file.
type oggPage struct {
	typ     byte
	granule int64
	serial  uint32
	lacing  []byte
	body    []byte
}

// readOggPage reads the next page.
func readOggPage(br *bufio.Reader) (*oggPage, error) {
	var hdr [27]byte
	if _, err := io.ReadFull(br, hdr[:]); err != nil {
		return nil, err
	}
	if string(hdr[:4]) != "OggS" || hdr[4] != 0 {
		return nil, errors.New("not an Ogg page")
	}
	p := &oggPage{typ: hdr[5], granule: int64(binary.LittleEndian.Uint64(hdr[6:])), serial: binary.LittleEndian.Uint32(hdr[14:])}
	p.lacing = make([]byte, hdr[26])
	if _, err := io.ReadFull(br, p.lacing); err != nil {
		return nil, err
	}
	n := 0
	for _, l := range p.lacing {
		n += int(l)
	}
	p.body = make([]byte, n)
	if _, err := io.ReadFull(br, p.body); err != nil {
		return nil, err
	}
	return p, nil
}

// oggAudio tells the coding of a stream by its first packet: "vorbis",
// "opus", "flac", "speex", "" for others (Theora video, Skeleton, text).
func oggAudio(packet []byte) string {
	switch {
	case bytes.HasPrefix(packet, []byte("\x01vorbis")):
		return "vorbis"
	case bytes.HasPrefix(packet, []byte("OpusHead")):
		return "opus"
	case bytes.HasPrefix(packet, []byte("\x7fFLAC")):
		return "flac"
	case bytes.HasPrefix(packet, []byte("Speex   ")):
		return "speex"
	}
	return ""
}

// isOggAudio reports whether the first page of an Ogg file (head: its
// first bytes) begins an audio stream.
func isOggAudio(head []byte) bool {
	if len(head) < 28 || string(head[:4]) != "OggS" {
		return false
	}
	n := int(head[26])
	if 27+n+8 > len(head) {
		return false
	}
	return oggAudio(head[27+n:]) != ""
}

// readOgg reads the headers of the first audio stream of an Ogg file
// (identification and comments) and the time of its last page.
func readOgg(r io.ReaderAt, size int64, t *track) error {
	br := bufio.NewReaderSize(io.NewSectionReader(r, 0, min(size, maxOggHeaders)), 64<<10)
	var serial uint32
	var codec string
	var packets [][]byte
	var packet []byte
	done := false
	for pages := 0; pages < maxOggPages && !done; pages++ {
		p, err := readOggPage(br)
		if err != nil {
			break
		}
		if codec == "" {
			if p.typ&2 == 0 || len(p.body) == 0 {
				continue
			}
			if c := oggAudio(p.body); c != "" {
				codec, serial = c, p.serial
			} else {
				if bytes.HasPrefix(p.body, []byte("\x80theora")) || bytes.HasPrefix(p.body, []byte("BBCD")) {
					t.warn("the file has a video stream, which is not described")
				}
				continue
			}
		} else if p.serial != serial {
			continue
		}
		off := 0
		for _, l := range p.lacing {
			packet = append(packet, p.body[off:off+int(l)]...)
			off += int(l)
			if l < 255 {
				packets = append(packets, packet)
				packet = nil
				if done = headersComplete(codec, packets); done {
					break
				}
			}
		}
	}
	if codec == "" {
		return errors.New("the file has no audio stream")
	}
	t.readOggHeaders(codec, packets)
	// the last page of the stream tells how long it is
	tail := min(size, oggTail)
	b := make([]byte, tail)
	n, _ := r.ReadAt(b, size-tail)
	b = b[:n]
	for at := bytes.LastIndex(b, []byte("OggS")); at >= 0; at = bytes.LastIndex(b[:at], []byte("OggS")) {
		if at+27 > len(b) || b[at+4] != 0 || binary.LittleEndian.Uint32(b[at+14:]) != serial {
			continue
		}
		granule := int64(binary.LittleEndian.Uint64(b[at+6:]))
		if granule < 0 {
			continue
		}
		t.oggDuration(codec, granule)
		break
	}
	if t.duration > 0 && t.bitrate == 0 {
		t.bitrate = int(float64(size)*8/t.duration + 0.5)
	}
	return nil
}

// headersComplete reports whether the packets read hold the headers the
// reader needs: identification and comments.
func headersComplete(codec string, packets [][]byte) bool {
	switch codec {
	case "flac":
		// the first packet holds STREAMINFO; the metadata blocks that
		// follow end with one flagged last
		if len(packets) < 2 {
			return false
		}
		last := packets[len(packets)-1]
		return len(last) > 0 && last[0]&0x80 != 0
	}
	return len(packets) >= 2
}

// readOggHeaders reads the identification header and the comments of a
// stream.
func (t *track) readOggHeaders(codec string, packets [][]byte) {
	le := binary.LittleEndian
	ident := packets[0]
	switch codec {
	case "vorbis":
		t.codec = "Vorbis"
		if len(ident) >= 28 {
			t.channels = int(ident[11])
			t.sampleRate = int(le.Uint32(ident[12:]))
			if nominal := int32(le.Uint32(ident[20:])); nominal > 0 {
				t.bitrate = int(nominal)
			}
		}
		if len(packets) > 1 && bytes.HasPrefix(packets[1], []byte("\x03vorbis")) {
			parseVorbisComment(packets[1][7:], t)
		}
	case "opus":
		t.codec = "Opus"
		t.sampleRate = 48000
		if len(ident) >= 12 {
			t.channels = int(ident[9])
			t.preskip = int64(le.Uint16(ident[10:]))
		}
		if len(packets) > 1 && bytes.HasPrefix(packets[1], []byte("OpusTags")) {
			parseVorbisComment(packets[1][8:], t)
		}
	case "speex":
		t.codec = "Speex"
		if len(ident) >= 56 {
			t.sampleRate = int(le.Uint32(ident[36:]))
			t.channels = int(le.Uint32(ident[48:]))
			if br := int32(le.Uint32(ident[52:])); br > 0 {
				t.bitrate = int(br)
			}
		}
		if len(packets) > 1 {
			parseVorbisComment(packets[1], t)
		}
	case "flac":
		t.codec = "FLAC"
		if len(ident) >= 13+4+34 && string(ident[9:13]) == "fLaC" {
			t.readStreamInfo(ident[17:])
		}
		for _, p := range packets[1:] {
			if len(p) < 4 {
				continue
			}
			switch p[0] & 0x7F {
			case 4:
				parseVorbisComment(p[4:], t)
			case 6:
				if pic, ok := parsePictureBlock(p[4:]); ok && len(t.pictures) < maxValues {
					t.pictures = append(t.pictures, pic)
				}
			}
		}
	}
}

// oggDuration sets the duration from the granule position of the last
// page: samples for Vorbis, Speex and FLAC, 48 kHz samples less the
// pre-skip for Opus.
func (t *track) oggDuration(codec string, granule int64) {
	switch codec {
	case "opus":
		t.duration = float64(granule-t.preskip) / 48000
	default:
		if t.sampleRate > 0 {
			t.duration = float64(granule) / float64(t.sampleRate)
		}
	}
	if t.duration < 0 {
		t.duration = 0
	}
}
