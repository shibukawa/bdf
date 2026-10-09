package audio

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"slices"
	"testing"
)

// Builders of test files.

var le = binary.LittleEndian

func u16be(v int) []byte   { return be.AppendUint16(nil, uint16(v)) }
func u32be(v int) []byte   { return be.AppendUint32(nil, uint32(v)) }
func u16le(v int) []byte   { return le.AppendUint16(nil, uint16(v)) }
func u32le(v int) []byte   { return le.AppendUint32(nil, uint32(v)) }
func u64le(v int64) []byte { return le.AppendUint64(nil, uint64(v)) }

func syncsafeBytes(n int) []byte {
	return []byte{byte(n >> 21 & 0x7F), byte(n >> 14 & 0x7F), byte(n >> 7 & 0x7F), byte(n & 0x7F)}
}

// id3Tag makes an ID3v2 tag of a version with the tag flags and frames.
func id3Tag(version, flags byte, frames ...[]byte) []byte {
	body := slices.Concat(frames...)
	tag := slices.Concat([]byte("ID3"), []byte{version, 0, flags}, syncsafeBytes(len(body)), body)
	if version == 4 && flags&0x10 != 0 {
		tag = slices.Concat(tag, []byte("3DI"), []byte{version, 0, flags}, syncsafeBytes(len(body)))
	}
	return tag
}

// frameOf makes a frame of a version; v2.4 sizes are syncsafe.
func frameOf(version int, id string, flags uint16, data []byte) []byte {
	switch version {
	case 2:
		return slices.Concat([]byte(id[:3]), []byte{byte(len(data) >> 16), byte(len(data) >> 8), byte(len(data))}, data)
	case 3:
		return slices.Concat([]byte(id), u32be(len(data)), u16be(int(flags)), data)
	}
	return slices.Concat([]byte(id), syncsafeBytes(len(data)), u16be(int(flags)), data)
}

// textFrame makes a text frame in an encoding (0 legacy, 1 UTF-16 with a
// BOM, 2 UTF-16BE, 3 UTF-8).
func textFrame(version int, id string, enc byte, values ...string) []byte {
	var data []byte
	data = append(data, enc)
	for i, v := range values {
		if i > 0 {
			if enc == 1 || enc == 2 {
				data = append(data, 0, 0)
			} else {
				data = append(data, 0)
			}
		}
		data = append(data, encodeText(enc, v)...)
	}
	return frameOf(version, id, 0, data)
}

func encodeText(enc byte, s string) []byte {
	switch enc {
	case 1:
		out := []byte{0xFF, 0xFE}
		for _, r := range s {
			out = le.AppendUint16(out, uint16(r))
		}
		return out
	case 2:
		var out []byte
		for _, r := range s {
			out = be.AppendUint16(out, uint16(r))
		}
		return out
	}
	return []byte(s)
}

func deflate(b []byte) []byte {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write(b)
	w.Close()
	return buf.Bytes()
}

// mpegFrames makes n MPEG-1 layer III frames of 128 kbit/s at 44.1 kHz,
// stereo (417 bytes each), the first with a Xing header saying xing
// frames when xing is not 0.
func mpegFrames(n, xing int) []byte {
	var out []byte
	for i := 0; i < n; i++ {
		f := make([]byte, 417)
		copy(f, []byte{0xFF, 0xFB, 0x90, 0x00})
		if i == 0 && xing != 0 {
			copy(f[4+32:], "Xing")
			copy(f[4+36:], u32be(3))
			copy(f[4+40:], u32be(xing))
			copy(f[4+44:], u32be(xing*417))
		}
		out = append(out, f...)
	}
	return out
}

// adtsFrames makes n ADTS frames of AAC LC at 22.05 kHz, mono, of size
// bytes each.
func adtsFrames(n, size int) []byte {
	var out []byte
	for i := 0; i < n; i++ {
		f := make([]byte, size)
		f[0], f[1] = 0xFF, 0xF1
		f[2] = 0x40 | 7<<2 // profile LC, sampling index 7
		f[3] = 1<<6 | byte(size>>11&3)
		f[4] = byte(size >> 3)
		f[5] = byte(size&7) << 5
		f[6] = 0xFC
		out = append(out, f...)
	}
	return out
}

func id3v1(title, artist, album, year, comment string, track, genre byte) []byte {
	pad := func(s string, n int) []byte {
		b := make([]byte, n)
		copy(b, s)
		return b
	}
	c := pad(comment, 30)
	if track > 0 {
		c[28], c[29] = 0, track
	}
	return slices.Concat([]byte("TAG"), pad(title, 30), pad(artist, 30), pad(album, 30), pad(year, 4), c, []byte{genre})
}

// box makes an ISOBMFF box.
func box(typ string, content ...[]byte) []byte {
	b := slices.Concat(content...)
	return slices.Concat(u32be(8+len(b)), []byte(typ), b)
}

// fullBox makes a box with a version and flags.
func fullBox(typ string, version byte, content ...[]byte) []byte {
	return box(typ, []byte{version, 0, 0, 0}, slices.Concat(content...))
}

// item makes an iTunes metadata item with one data atom of a type.
func item(name string, typ int, payload []byte) []byte {
	return box(name, box("data", u32be(typ), u32be(0), payload))
}

// oggPages writes packets as pages of a stream, each packet on a page of
// its own unless it is longer than the segments of a page can hold; the
// first page is marked as the beginning of the stream, the last as its
// end with the granule position.
func oggPages(serial uint32, granule int64, packets ...[]byte) []byte {
	var out []byte
	seq := 0
	page := func(typ byte, gran int64, lacing, body []byte) {
		hdr := slices.Concat([]byte("OggS"), []byte{0, typ}, le.AppendUint64(nil, uint64(gran)), u32le(int(serial)), u32le(seq), u32le(0), []byte{byte(len(lacing))})
		out = slices.Concat(out, hdr, lacing, body)
		seq++
	}
	for i, p := range packets {
		typ := byte(0)
		if i == 0 {
			typ = 2
		}
		for len(p) > 255*200 {
			// a packet longer than a page: continued on the next
			page(typ, -1, bytes.Repeat([]byte{255}, 200), p[:255*200])
			p = p[255*200:]
			typ = 1
		}
		var lacing []byte
		for n := len(p); n >= 255; n -= 255 {
			lacing = append(lacing, 255)
		}
		lacing = append(lacing, byte(len(p)%255))
		gran := int64(0)
		if i == len(packets)-1 {
			typ |= 4
			gran = granule
		}
		page(typ, gran, lacing, p)
	}
	return out
}

// vorbisComment writes a Vorbis comment with fields "KEY=value".
func vorbisComment(fields ...string) []byte {
	out := slices.Concat(u32le(4), []byte("test"), u32le(len(fields)))
	for _, f := range fields {
		out = slices.Concat(out, u32le(len(f)), []byte(f))
	}
	return out
}

// riffChunk makes a RIFF chunk (little-endian size, padded to even).
func riffChunk(id string, content ...[]byte) []byte {
	b := slices.Concat(content...)
	out := slices.Concat([]byte(id), u32le(len(b)), b)
	if len(b)%2 == 1 {
		out = append(out, 0)
	}
	return out
}

// iffChunk makes an IFF chunk (big-endian size, padded to even).
func iffChunk(id string, content ...[]byte) []byte {
	b := slices.Concat(content...)
	out := slices.Concat([]byte(id), u32be(len(b)), b)
	if len(b)%2 == 1 {
		out = append(out, 0)
	}
	return out
}

// flacBlock makes a FLAC metadata block.
func flacBlock(typ byte, last bool, content ...[]byte) []byte {
	b := slices.Concat(content...)
	h := typ
	if last {
		h |= 0x80
	}
	return slices.Concat([]byte{h, byte(len(b) >> 16), byte(len(b) >> 8), byte(len(b))}, b)
}

// streamInfo makes a STREAMINFO block body.
func streamInfoBlock(rate, channels, bits int, samples uint64) []byte {
	b := make([]byte, 34)
	v := uint64(rate)<<44 | uint64(channels-1)<<41 | uint64(bits-1)<<36 | samples
	be.PutUint64(b[10:], v)
	return b
}

// pictureBlock makes a FLAC picture block body.
func pictureBlock(typ int, mime, desc string, data []byte) []byte {
	return slices.Concat(u32be(typ), u32be(len(mime)), []byte(mime), u32be(len(desc)), []byte(desc), u32be(1), u32be(1), u32be(24), u32be(0), u32be(len(data)), data)
}

// tinyPNG is a 1 × 1 PNG.
var tinyPNG = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00\x00\x1f\x15\xc4\x89\x00\x00\x00\nIDATx\x9cc\x00\x01\x00\x00\x05\x00\x01\r\n-\xb4\x00\x00\x00\x00IEND\xaeB`\x82")

// readTrack reads a file in memory.
func readTrack(t *testing.T, data []byte) *track {
	t.Helper()
	tr, err := read(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	return tr
}

// want checks the values of a key.
func want(t *testing.T, tg tags, key string, values ...string) {
	t.Helper()
	if !slices.Equal(tg[key], values) {
		t.Errorf("%s = %q, want %q", key, tg[key], values)
	}
}

func slicesConcat(b ...[]byte) []byte { return slices.Concat(b...) }

func bytesReader(b []byte) *bytes.Reader { return bytes.NewReader(b) }
