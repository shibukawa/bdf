package audio

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

// mpegFrame is the header of an MPEG audio frame (MPEG-1, MPEG-2 and
// MPEG-2.5, layers I to III).
type mpegFrame struct {
	version    int // 1, 2, or 25 for MPEG-2.5
	layer      int
	bitrate    int // bits per second; 0 for the free format
	sampleRate int
	channels   int
	length     int // bytes, header included
	samples    int // per frame
}

var (
	mpegBitrates = [5][16]int{ // kbit/s by version and layer
		{0, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448, 0}, // MPEG-1 layer I
		{0, 32, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384, 0},    // MPEG-1 layer II
		{0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 0},     // MPEG-1 layer III
		{0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256, 0},    // MPEG-2, 2.5 layer I
		{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160, 0},         // MPEG-2, 2.5 layers II, III
	}
	mpegSampleRates = map[int][3]int{1: {44100, 48000, 32000}, 2: {22050, 24000, 16000}, 25: {11025, 12000, 8000}}
)

// parseMPEGHeader reads a frame header; ok is false for bytes that are
// not one.
func parseMPEGHeader(b []byte) (f mpegFrame, ok bool) {
	if len(b) < 4 || b[0] != 0xFF || b[1]&0xE0 != 0xE0 {
		return f, false
	}
	switch b[1] >> 3 & 3 {
	case 0:
		f.version = 25
	case 2:
		f.version = 2
	case 3:
		f.version = 1
	default:
		return f, false
	}
	f.layer = 4 - int(b[1]>>1&3)
	if f.layer == 4 {
		return f, false
	}
	br, sr := int(b[2]>>4), int(b[2]>>2&3)
	if br == 15 || sr == 3 {
		return f, false
	}
	if b[3]>>0&3 == 2 { // emphasis
		return f, false
	}
	padding := int(b[2] >> 1 & 1)
	f.sampleRate = mpegSampleRates[f.version][sr]
	f.channels = 2
	if b[3]>>6 == 3 {
		f.channels = 1
	}
	table := f.layer - 1
	if f.version != 1 {
		table = 3
		if f.layer > 1 {
			table = 4
		}
	}
	f.bitrate = mpegBitrates[table][br] * 1000
	switch {
	case f.layer == 1:
		f.samples = 384
	case f.layer == 2 || f.version == 1:
		f.samples = 1152
	default:
		f.samples = 576
	}
	if f.bitrate == 0 {
		return f, false // free format: the frame length is not known
	}
	switch {
	case f.layer == 1:
		f.length = (12*f.bitrate/f.sampleRate + padding) * 4
	case f.layer == 3 && f.version != 1:
		f.length = 72*f.bitrate/f.sampleRate + padding
	default:
		f.length = 144*f.bitrate/f.sampleRate + padding
	}
	return f, f.length >= 4
}

// streamInfo is what a scan of the frames of a stream finds.
type streamInfo struct {
	codec      string
	duration   float64
	bitrate    int
	sampleRate int
	channels   int
	frames     int
	vbr        bool
}

// maxJunk is how far a scan looks for the first frame, and for the next
// one after bytes that are not a frame.
const maxJunk = 64 << 10

// scanner reads a stream frame by frame.
type scanner struct {
	br   *bufio.Reader
	pos  int64
	size int64
}

func newScanner(r io.ReaderAt, start, end int64) *scanner {
	return &scanner{br: bufio.NewReaderSize(io.NewSectionReader(r, start, end-start), 64<<10), size: end - start}
}

// peek returns the next n bytes without reading them, or what is left.
func (s *scanner) peek(n int) []byte {
	b, _ := s.br.Peek(n)
	return b
}

func (s *scanner) skip(n int) {
	m, _ := s.br.Discard(n)
	s.pos += int64(m)
}

// isMPEG reports whether an MPEG audio frame starts at off of r, followed
// by another (or by the end of the file).
func isMPEG(r io.ReaderAt, off, size int64) bool {
	var b [8]byte
	n, _ := r.ReadAt(b[:], off)
	f, ok := parseMPEGHeader(b[:n])
	if !ok {
		return false
	}
	if off+int64(f.length)+4 > size {
		return off+int64(f.length) <= size+1 // the only frame, maybe cut short
	}
	n, _ = r.ReadAt(b[:4], off+int64(f.length))
	g, ok := parseMPEGHeader(b[:n])
	return ok && g.version == f.version && g.layer == f.layer && g.sampleRate == f.sampleRate
}

// scanMPEG reads the MPEG audio frames between start and end: the first
// frame's parameters, and the duration from the Xing/Info or VBRI header
// when there is one, else by counting the frames.
func scanMPEG(r io.ReaderAt, start, end int64) (*streamInfo, bool) {
	s := newScanner(r, start, end)
	info := &streamInfo{}
	var first *mpegFrame
	var samples float64 // seconds, summed per frame
	var bytes int64
	junk := 0
	for s.pos < s.size {
		h := s.peek(4)
		f, ok := parseMPEGHeader(h)
		if ok && first != nil && (f.version != first.version || f.layer != first.layer || f.sampleRate != first.sampleRate) {
			ok = false
		}
		if ok && first == nil {
			// the first frame must be followed by another, unless it is
			// the last bytes of the file
			if s.pos+int64(f.length)+4 <= s.size {
				g, next := parseMPEGHeader(s.peek(f.length + 4)[f.length:])
				ok = next && g.version == f.version && g.layer == f.layer && g.sampleRate == f.sampleRate
			}
		}
		if !ok {
			junk++
			if junk > maxJunk || first == nil && s.pos > maxJunk {
				break
			}
			s.skip(1)
			continue
		}
		junk = 0
		if first == nil {
			fc := f
			first = &fc
			info.sampleRate, info.channels = f.sampleRate, f.channels
			switch f.layer {
			case 1:
				info.codec = "MPEG layer I"
			case 2:
				info.codec = "MPEG layer II"
			default:
				info.codec = "MP3"
			}
			if frames, total, cbr, ok := vbrHeader(s.peek(f.length), f); ok {
				info.frames = frames
				info.duration = float64(frames) * float64(f.samples) / float64(f.sampleRate)
				if total <= 0 {
					total = s.size - s.pos
				}
				if info.duration > 0 {
					info.bitrate = int(float64(total)*8/info.duration + 0.5)
				}
				info.vbr = !cbr
				return info, true
			}
		}
		if f.bitrate != first.bitrate {
			info.vbr = true
		}
		info.frames++
		samples += float64(f.samples) / float64(f.sampleRate)
		bytes += int64(f.length)
		s.skip(f.length)
	}
	if first == nil {
		return nil, false
	}
	info.duration = samples
	switch {
	case !info.vbr:
		info.bitrate = first.bitrate // the nominal rate: frames are cut to whole bytes
	case samples > 0:
		info.bitrate = int(float64(bytes)*8/samples + 0.5)
	}
	return info, true
}

// vbrHeader reads the Xing/Info or VBRI header of the first frame: the
// number of frames and of bytes of the stream, and whether the header says
// the bitrate is constant (Info).
func vbrHeader(frame []byte, f mpegFrame) (frames int, total int64, cbr, ok bool) {
	side := 17
	if f.version == 1 {
		side = 32
	}
	if f.channels == 1 {
		side = (side + 1) / 2 // 9 or 17
	}
	if at := 4 + side; at+8 <= len(frame) && (string(frame[at:at+4]) == "Xing" || string(frame[at:at+4]) == "Info") {
		flags := binary.BigEndian.Uint32(frame[at+4:])
		p := at + 8
		if flags&1 != 0 && p+4 <= len(frame) {
			frames = int(binary.BigEndian.Uint32(frame[p:]))
			p += 4
		}
		if flags&2 != 0 && p+4 <= len(frame) {
			total = int64(binary.BigEndian.Uint32(frame[p:]))
		}
		return frames, total, string(frame[at:at+4]) == "Info", frames > 0
	}
	if at := 4 + 32; at+26 <= len(frame) && string(frame[at:at+4]) == "VBRI" {
		total = int64(binary.BigEndian.Uint32(frame[at+10:]))
		frames = int(binary.BigEndian.Uint32(frame[at+14:]))
		return frames, total, false, frames > 0
	}
	return 0, 0, false, false
}

// adtsSampleRates are the sampling frequencies of AAC by index.
var adtsSampleRates = [16]int{96000, 88200, 64000, 48000, 44100, 32000, 24000, 22050, 16000, 12000, 11025, 8000, 7350, 0, 0, 0}

// adtsChannels are the channels of AAC by channel configuration.
var adtsChannels = [8]int{0, 1, 2, 3, 4, 5, 6, 8}

// adtsFrame is the header of an ADTS frame.
type adtsFrame struct {
	profile    int // audio object type: 1 Main, 2 LC, 3 SSR, 4 LTP
	sampleRate int
	channels   int
	length     int
	blocks     int // raw data blocks, of 1024 samples each
}

func parseADTSHeader(b []byte) (f adtsFrame, ok bool) {
	if len(b) < 7 || b[0] != 0xFF || b[1]&0xF6 != 0xF0 {
		return f, false
	}
	f.profile = int(b[2]>>6) + 1
	f.sampleRate = adtsSampleRates[b[2]>>2&15]
	f.channels = adtsChannels[(b[2]&1)<<2|b[3]>>6]
	f.length = int(b[3]&3)<<11 | int(b[4])<<3 | int(b[5]>>5)
	f.blocks = int(b[6]&3) + 1
	return f, f.sampleRate > 0 && f.length >= 7
}

// isADTS reports whether an ADTS frame starts at off of r, followed by
// another (or by the end of the file).
func isADTS(r io.ReaderAt, off, size int64) bool {
	var b [8]byte
	n, _ := r.ReadAt(b[:], off)
	f, ok := parseADTSHeader(b[:n])
	if !ok {
		return false
	}
	if off+int64(f.length)+7 > size {
		return off+int64(f.length) <= size+1
	}
	n, _ = r.ReadAt(b[:7], off+int64(f.length))
	g, ok := parseADTSHeader(b[:n])
	return ok && g.sampleRate == f.sampleRate
}

// scanADTS reads the ADTS frames of a raw AAC stream.
func scanADTS(r io.ReaderAt, start, end int64) (*streamInfo, bool) {
	s := newScanner(r, start, end)
	info := &streamInfo{}
	var first *adtsFrame
	var seconds float64
	var bytes int64
	junk := 0
	for s.pos < s.size {
		f, ok := parseADTSHeader(s.peek(7))
		if ok && first != nil && f.sampleRate != first.sampleRate {
			ok = false
		}
		if ok && first == nil && s.pos+int64(f.length)+7 <= s.size {
			g, next := parseADTSHeader(s.peek(f.length + 7)[f.length:])
			ok = next && g.sampleRate == f.sampleRate
		}
		if !ok {
			junk++
			if junk > maxJunk || first == nil && s.pos > maxJunk {
				break
			}
			s.skip(1)
			continue
		}
		junk = 0
		if first == nil {
			fc := f
			first = &fc
			info.sampleRate, info.channels = f.sampleRate, f.channels
			info.codec = "AAC"
			switch f.profile {
			case 1:
				info.codec = "AAC Main"
			case 2:
				info.codec = "AAC LC"
			case 3:
				info.codec = "AAC SSR"
			case 4:
				info.codec = "AAC LTP"
			}
		}
		info.frames++
		seconds += float64(1024*f.blocks) / float64(f.sampleRate)
		bytes += int64(f.length)
		s.skip(f.length)
	}
	if first == nil {
		return nil, false
	}
	info.duration = seconds
	if seconds > 0 {
		info.bitrate = int(float64(bytes)*8/seconds + 0.5)
	}
	return info, true
}

// apeTagSize returns the length of an APEv2 tag at the end of a file (of
// which tail is the last 32 bytes, the tag's footer), or 0.
func apeTagSize(tail []byte) int64 {
	if len(tail) < 32 || !bytes.HasPrefix(tail, []byte("APETAGEX")) {
		return 0
	}
	size := int64(binary.LittleEndian.Uint32(tail[12:]))
	if binary.LittleEndian.Uint32(tail[20:])&1<<31 != 0 {
		size += 32 // a header
	}
	return size
}
