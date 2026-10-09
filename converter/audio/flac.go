package audio

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Bounds on FLAC and Vorbis metadata: a comment block, a picture block,
// and all the metadata of a file together.
const (
	maxComment  = 8 << 20
	maxPicture  = 16 << 20
	maxMetadata = 256 << 20
)

// readFLAC reads the metadata blocks of a FLAC stream that starts at
// start: STREAMINFO, VORBIS_COMMENT and PICTURE.
func readFLAC(r io.ReaderAt, start, size int64, t *track) error {
	var hdr [4]byte
	pos := start + 4
	var total int64
	for {
		if n, _ := r.ReadAt(hdr[:], pos); n < 4 {
			return errors.New("the FLAC metadata is truncated")
		}
		last, typ := hdr[0]&0x80 != 0, hdr[0]&0x7F
		length := int64(hdr[1])<<16 | int64(hdr[2])<<8 | int64(hdr[3])
		pos += 4
		if total += length; total > maxMetadata || pos+length > size {
			return fmt.Errorf("the FLAC metadata is larger than the file or than %d bytes", maxMetadata)
		}
		read := func(limit int64) []byte {
			if length > limit {
				t.warn(fmt.Sprintf("a FLAC metadata block of %d bytes is skipped", length))
				return nil
			}
			b := make([]byte, length)
			n, _ := r.ReadAt(b, pos)
			return b[:n]
		}
		switch typ {
		case 0:
			if b := read(64); len(b) >= 18 {
				t.readStreamInfo(b)
			}
		case 4:
			if b := read(maxComment); b != nil {
				parseVorbisComment(b, t)
			}
		case 6:
			if b := read(maxPicture); b != nil {
				if p, ok := parsePictureBlock(b); ok && len(t.pictures) < maxValues {
					t.pictures = append(t.pictures, p)
				}
			}
		}
		pos += length
		if last {
			break
		}
	}
	if t.duration > 0 && pos < size {
		t.bitrate = int(float64(size-pos)*8/t.duration + 0.5)
	}
	return nil
}

// readStreamInfo reads a STREAMINFO block: sample rate, channels, bits
// per sample and the number of samples.
func (t *track) readStreamInfo(b []byte) {
	t.codec = "FLAC"
	t.sampleRate = int(binary.BigEndian.Uint32(b[10:]) >> 12)
	t.channels = int(b[12]>>1&7) + 1
	t.bits = int(b[12]&1)<<4 | int(b[13]>>4) + 1
	samples := uint64(b[13]&0x0F)<<32 | uint64(binary.BigEndian.Uint32(b[14:]))
	if t.sampleRate > 0 && samples > 0 {
		t.duration = float64(samples) / float64(t.sampleRate)
	}
}

// parseVorbisComment reads a Vorbis comment (little-endian: the vendor
// string, then fields "KEY=value") into the track: the usual fields,
// pictures (METADATA_BLOCK_PICTURE, and the older COVERART), and
// chapters (CHAPTER001 and CHAPTER001NAME).
func parseVorbisComment(b []byte, t *track) {
	le := binary.LittleEndian
	if len(b) < 8 {
		return
	}
	vendor := int(le.Uint32(b))
	if vendor < 0 || vendor > len(b)-8 {
		return
	}
	b = b[4+vendor:]
	n := int(le.Uint32(b))
	b = b[4:]
	chapters := map[string]*chapter{}
	var coverArt, coverMIME string
	for i := 0; i < n && len(b) >= 4; i++ {
		size := int(le.Uint32(b))
		b = b[4:]
		if size < 0 || size > len(b) {
			break
		}
		field := b[:size]
		b = b[size:]
		key, value, ok := strings.Cut(string(field), "=")
		if !ok {
			continue
		}
		key = strings.ToUpper(key)
		switch {
		case key == "METADATA_BLOCK_PICTURE":
			if raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(value)); err == nil {
				if p, ok := parsePictureBlock(raw); ok && len(t.pictures) < maxValues {
					t.pictures = append(t.pictures, p)
				}
			}
		case key == "COVERART":
			coverArt = value
		case key == "COVERARTMIME":
			coverMIME = value
		case strings.HasPrefix(key, "CHAPTER") && len(key) >= 10:
			name := strings.HasSuffix(key, "NAME")
			num := strings.TrimSuffix(key[7:], "NAME")
			if _, err := strconv.Atoi(num); err != nil || len(chapters) >= maxChapters && chapters[num] == nil {
				continue
			}
			c := chapters[num]
			if c == nil {
				c = &chapter{start: -1}
				chapters[num] = c
			}
			if name {
				c.title = value
			} else {
				c.start = chapterTime(value)
			}
		case key == "TRACKNUMBER" || key == "DISCNUMBER":
			t.tags.addNumber(userKey(key), totalKey(userKey(key)), value)
		default:
			if k := userKey(key); k != "" {
				t.tags.add(k, value)
			}
		}
	}
	if coverArt != "" && len(t.pictures) < maxValues {
		if raw, err := base64.StdEncoding.DecodeString(strings.Join(strings.Fields(coverArt), "")); err == nil {
			t.pictures = append(t.pictures, picture{typ: 3, mime: strings.ToLower(coverMIME), data: raw})
		}
	}
	if len(chapters) > 0 {
		nums := make([]string, 0, len(chapters))
		for k := range chapters {
			nums = append(nums, k)
		}
		sortNumbers(nums)
		for _, k := range nums {
			if c := chapters[k]; c.start >= 0 {
				title := c.title
				if title == "" {
					title = "Chapter " + strings.TrimLeft(k, "0")
				}
				t.addChapter(c.start, title)
			}
		}
	}
}

// sortNumbers sorts the numbers of chapters as numbers ("002" before "10").
func sortNumbers(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0; j-- {
			a, _ := strconv.Atoi(s[j-1])
			b, _ := strconv.Atoi(s[j])
			if a <= b {
				break
			}
			s[j-1], s[j] = s[j], s[j-1]
		}
	}
}

var chapterTimeRE = regexp.MustCompile(`^(\d+):(\d+):(\d+)(?:\.(\d+))?$`)

// chapterTime reads a chapter's start time "HH:MM:SS.mmm"; -1 when it is
// not one.
func chapterTime(s string) float64 {
	m := chapterTimeRE.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return -1
	}
	h, _ := strconv.Atoi(m[1])
	mi, _ := strconv.Atoi(m[2])
	sec, _ := strconv.Atoi(m[3])
	t := float64(h*3600 + mi*60 + sec)
	if m[4] != "" {
		f, _ := strconv.ParseFloat("0."+m[4], 64)
		t += f
	}
	return t
}

// parsePictureBlock reads a FLAC picture block (big-endian: the type, the
// media type, a description, the size, and the picture).
func parsePictureBlock(b []byte) (picture, bool) {
	be := binary.BigEndian
	if len(b) < 32 {
		return picture{}, false
	}
	typ := int(be.Uint32(b))
	n := int(be.Uint32(b[4:]))
	if n < 0 || n > len(b)-8 {
		return picture{}, false
	}
	mime := strings.ToLower(string(b[8 : 8+n]))
	b = b[8+n:]
	if len(b) < 4 {
		return picture{}, false
	}
	n = int(be.Uint32(b))
	if n < 0 || n > len(b)-4 {
		return picture{}, false
	}
	b = b[4+n:] // the description
	if len(b) < 20 {
		return picture{}, false
	}
	n = int(be.Uint32(b[16:]))
	b = b[20:]
	if n <= 0 || n > len(b) {
		return picture{}, false
	}
	return picture{typ: typ, mime: mime, data: b[:n]}, true
}
