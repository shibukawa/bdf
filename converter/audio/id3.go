package audio

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Bounds on ID3 tags: how much of a tag is read (its size field says up to
// 256 MiB), and how large a compressed frame may inflate to.
const (
	maxID3     = 64 << 20
	maxInflate = 16 << 20
)

// id3Frame is a frame of an ID3v2 tag, with its v2.3/v2.4 identifier.
type id3Frame struct {
	id   string
	data []byte
}

// id3v2 is a parsed ID3v2 tag.
type id3v2 struct {
	version int
	frames  []id3Frame
	// truncated reports a tag that claims more than the file holds.
	truncated bool
}

var errNotID3 = errors.New("not an ID3v2 tag")

// id3v2Size returns the length of the ID3v2 tag a file starts with: its
// header, body and footer; 0 when the file has no tag.
func id3v2Size(head []byte) int64 {
	if len(head) < 10 || string(head[:3]) != "ID3" || head[3] < 2 || head[3] > 4 || head[3] == 0xFF || head[4] == 0xFF {
		return 0
	}
	for _, b := range head[6:10] {
		if b >= 0x80 {
			return 0
		}
	}
	n := 10 + int64(syncsafe(head[6:10]))
	if head[3] == 4 && head[5]&0x10 != 0 {
		n += 10 // a footer
	}
	return n
}

func syncsafe(b []byte) uint32 {
	return uint32(b[0]&0x7F)<<21 | uint32(b[1]&0x7F)<<14 | uint32(b[2]&0x7F)<<7 | uint32(b[3]&0x7F)
}

// unsync undoes unsynchronisation: FF 00 is FF.
func unsync(b []byte) []byte {
	if !bytes.Contains(b, []byte{0xFF, 0x00}) {
		return b
	}
	out := make([]byte, 0, len(b))
	for i := 0; i < len(b); i++ {
		out = append(out, b[i])
		if b[i] == 0xFF && i+1 < len(b) && b[i+1] == 0 {
			i++
		}
	}
	return out
}

// parseID3v2 parses the tag that b (the tag, or the file) starts with.
func parseID3v2(b []byte) (*id3v2, error) {
	if id3v2Size(b) == 0 {
		return nil, errNotID3
	}
	tag := &id3v2{version: int(b[3])}
	flags := b[5]
	size := int(syncsafe(b[6:10]))
	body := b[10:]
	if size > len(body) {
		size, tag.truncated = len(body), true
	}
	body = body[:size]
	if tag.version >= 3 && flags&0x40 != 0 && len(body) >= 4 {
		// the extended header
		skip := 0
		if tag.version == 3 {
			skip = 4 + int(binary.BigEndian.Uint32(body))
		} else {
			skip = int(syncsafe(body))
		}
		if skip < 0 || skip > len(body) {
			skip = len(body)
		}
		body = body[skip:]
	}
	if tag.version < 4 && flags&0x80 != 0 {
		body = unsync(body)
	}
	tag.frames = parseFrames(body, tag.version, tag.version == 4 && flags&0x80 != 0)
	return tag, nil
}

// v22IDs maps the three-letter frame identifiers of ID3v2.2 to those of
// v2.3 and v2.4.
var v22IDs = map[string]string{
	"TT1": "TIT1", "TT2": "TIT2", "TT3": "TIT3", "TP1": "TPE1", "TP2": "TPE2", "TP3": "TPE3", "TP4": "TPE4", "TAL": "TALB",
	"TCM": "TCOM", "TXT": "TEXT", "TCO": "TCON", "TYE": "TYER", "TDA": "TDAT", "TIM": "TIME", "TRK": "TRCK", "TPA": "TPOS",
	"TPB": "TPUB", "TCR": "TCOP", "TLA": "TLAN", "TRC": "TSRC", "TBP": "TBPM", "TOR": "TORY", "TXX": "TXXX", "COM": "COMM",
	"ULT": "USLT", "SLT": "SYLT", "PIC": "APIC",
}

// validFrameID reports whether b starts with a frame identifier of n
// letters and digits.
func validFrameID(b []byte, n int) bool {
	if len(b) < n {
		return false
	}
	for _, c := range b[:n] {
		if !(c >= 'A' && c <= 'Z' || c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}

// parseFrames splits the body of a tag (or of a CHAP frame) into frames.
func parseFrames(body []byte, version int, allUnsync bool) []id3Frame {
	var out []id3Frame
	idLen, hdr := 4, 10
	if version == 2 {
		idLen, hdr = 3, 6
	}
	for off := 0; off+hdr <= len(body); {
		if body[off] == 0 {
			break // padding
		}
		if !validFrameID(body[off:], idLen) {
			break
		}
		id := string(body[off : off+idLen])
		var size int
		var flags uint16
		switch version {
		case 2:
			size = int(body[off+3])<<16 | int(body[off+4])<<8 | int(body[off+5])
		case 3:
			size = int(binary.BigEndian.Uint32(body[off+4:]))
			flags = binary.BigEndian.Uint16(body[off+8:])
		default:
			size = int(syncsafe(body[off+4:]))
			flags = binary.BigEndian.Uint16(body[off+8:])
			// Some writers put plain sizes in v2.4 frames: take the one
			// that leads to the next frame (or the end).
			plain := int(binary.BigEndian.Uint32(body[off+4:]))
			if plain != size && plain >= 0 && off+hdr+plain <= len(body) &&
				!frameFollows(body, off+hdr+size, idLen) && frameFollows(body, off+hdr+plain, idLen) {
				size = plain
			}
		}
		if size < 0 || off+hdr+size > len(body) {
			size = len(body) - off - hdr // a truncated last frame
		}
		data := body[off+hdr : off+hdr+size]
		off += hdr + size
		if version == 2 {
			id = v22IDs[id]
			if id == "" {
				continue
			}
		}
		var compressed, encrypted, unsynced bool
		switch version {
		case 3:
			compressed, encrypted = flags&0x80 != 0, flags&0x40 != 0
			if compressed && len(data) >= 4 {
				data = data[4:]
			}
			if encrypted && len(data) >= 1 {
				data = data[1:]
			}
			if flags&0x20 != 0 && len(data) >= 1 {
				data = data[1:]
			}
		case 4:
			if flags&0x40 != 0 && len(data) >= 1 {
				data = data[1:]
			}
			compressed, encrypted, unsynced = flags&0x08 != 0, flags&0x04 != 0, flags&0x02 != 0 || allUnsync
			if encrypted && len(data) >= 1 {
				data = data[1:]
			}
			if flags&0x01 != 0 && len(data) >= 4 {
				data = data[4:]
			}
		}
		if encrypted {
			continue
		}
		if unsynced {
			data = unsync(data)
		}
		if compressed {
			var err error
			if data, err = inflate(data); err != nil {
				continue
			}
		}
		out = append(out, id3Frame{id, data})
	}
	return out
}

// frameFollows reports whether a frame or the padding (or the end) is at
// off.
func frameFollows(body []byte, off, idLen int) bool {
	return off >= len(body) || body[off] == 0 || validFrameID(body[off:], idLen)
}

// inflate decompresses a compressed frame, to at most maxInflate bytes.
func inflate(b []byte) ([]byte, error) {
	z, err := zlib.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer z.Close()
	out, err := io.ReadAll(io.LimitReader(z, maxInflate+1))
	if err != nil && len(out) == 0 {
		return nil, err
	}
	if len(out) > maxInflate {
		return nil, fmt.Errorf("the frame inflates to more than %d bytes", maxInflate)
	}
	return out, nil
}

// id3Strings decodes the strings of a text frame (after its encoding
// byte): NUL separates several values (v2.4), and ends the last.
func id3Strings(enc byte, b []byte, l legacy) []string {
	var out []string
	add := func(s string) {
		if s = clean(s); s != "" {
			out = append(out, s)
		}
	}
	switch enc {
	case 1, 2:
		bigEndian := enc == 2
		for len(b) > 0 {
			if len(b) >= 2 && !bigEndian && (b[0] == 0xFE && b[1] == 0xFF) {
				bigEndian = true
			}
			end := utf16End(b)
			add(decodeUTF16(b[:end], bigEndian))
			if end+2 > len(b) {
				break
			}
			b = b[end+2:]
		}
	case 3:
		for _, s := range bytes.Split(b, []byte{0}) {
			add(string(bytes.ToValidUTF8(s, []byte("�"))))
		}
	default:
		for _, s := range bytes.Split(b, []byte{0}) {
			add(l.decode(s))
		}
	}
	return out
}

// utf16End returns the offset of the NUL character that ends a UTF-16
// string, or the length.
func utf16End(b []byte) int {
	for i := 0; i+1 < len(b); i += 2 {
		if b[i] == 0 && b[i+1] == 0 {
			return i
		}
	}
	return len(b)
}

// splitTerminated splits a string terminated by NUL (in the frame's
// encoding) from what follows it.
func splitTerminated(enc byte, b []byte) (s, rest []byte) {
	if enc == 1 || enc == 2 {
		end := utf16End(b)
		if end+2 <= len(b) {
			return b[:end], b[end+2:]
		}
		return b[:end], nil
	}
	if i := bytes.IndexByte(b, 0); i >= 0 {
		return b[:i], b[i+1:]
	}
	return b, nil
}

// id3String decodes one string of a frame in its encoding.
func id3String(enc byte, b []byte, l legacy) string {
	switch enc {
	case 1:
		return clean(decodeUTF16(b, false))
	case 2:
		return clean(decodeUTF16(b, true))
	case 3:
		return clean(string(bytes.ToValidUTF8(b, []byte("�"))))
	}
	return clean(l.decode(b))
}

// legacySamples gathers the strings a tag stores in the legacy encoding,
// for guessLegacy.
func legacySamples(frames []id3Frame) [][]byte {
	var out [][]byte
	for _, f := range frames {
		if len(f.data) < 2 || f.data[0] != 0 {
			continue
		}
		switch {
		case f.id == "COMM" || f.id == "USLT":
			if len(f.data) > 4 {
				out = append(out, f.data[4:])
			}
		case f.id[0] == 'T':
			out = append(out, f.data[1:])
		}
	}
	return out
}

// apply reads the frames of a tag into the track.
func (t *track) applyID3v2(tag *id3v2) {
	if tag.truncated {
		t.warn("the ID3v2 tag is truncated")
	}
	t.applyFrames(tag.frames, tag.version)
}

// applyFrames reads frames into the track's tags, pictures and chapters.
func (t *track) applyFrames(frames []id3Frame, version int) {
	l := guessLegacy(legacySamples(frames))
	tg := t.tags
	var year, ddmm, hhmm, date, release string
	var comment, lyrics described
	for _, f := range frames {
		if len(f.data) == 0 {
			continue
		}
		enc, rest := f.data[0], f.data[1:]
		switch f.id {
		case "TIT2":
			tg.add(keyTitle, id3Strings(enc, rest, l)...)
		case "TIT3":
			tg.add(keySubtitle, id3Strings(enc, rest, l)...)
		case "TIT1":
			tg.add(keyGrouping, id3Strings(enc, rest, l)...)
		case "TPE1":
			tg.add(keyArtist, splitPeople(id3Strings(enc, rest, l))...)
		case "TPE2":
			tg.add(keyAlbumArtist, splitPeople(id3Strings(enc, rest, l))...)
		case "TPE3":
			tg.add(keyConductor, splitPeople(id3Strings(enc, rest, l))...)
		case "TPE4":
			tg.add(keyRemixer, splitPeople(id3Strings(enc, rest, l))...)
		case "TALB":
			tg.add(keyAlbum, id3Strings(enc, rest, l)...)
		case "TCOM":
			tg.add(keyComposer, splitPeople(id3Strings(enc, rest, l))...)
		case "TEXT":
			tg.add(keyLyricist, splitPeople(id3Strings(enc, rest, l))...)
		case "TCON":
			for _, s := range id3Strings(enc, rest, l) {
				tg.add(keyGenre, genreNames(s)...)
			}
		case "TYER":
			year = first(id3Strings(enc, rest, l))
		case "TDAT":
			ddmm = first(id3Strings(enc, rest, l))
		case "TIME":
			hhmm = first(id3Strings(enc, rest, l))
		case "TDRC":
			date = first(id3Strings(enc, rest, l))
		case "TDRL":
			release = first(id3Strings(enc, rest, l))
		case "TDOR", "TORY":
			tg.add(keyOriginalDate, first(id3Strings(enc, rest, l)))
		case "TRCK":
			tg.addNumber(keyTrack, keyTrackTotal, first(id3Strings(enc, rest, l)))
		case "TPOS":
			tg.addNumber(keyDisc, keyDiscTotal, first(id3Strings(enc, rest, l)))
		case "TPUB":
			tg.add(keyPublisher, id3Strings(enc, rest, l)...)
		case "TCOP":
			tg.add(keyCopyright, id3Strings(enc, rest, l)...)
		case "TLAN":
			tg.add(keyLanguage, id3Strings(enc, rest, l)...)
		case "TSRC":
			tg.add(keyISRC, id3Strings(enc, rest, l)...)
		case "TBPM":
			tg.add(keyBPM, first(id3Strings(enc, rest, l)))
		case "TXXX":
			desc, value := splitTerminated(enc, rest)
			if key := userKey(id3String(enc, desc, l)); key != "" {
				if key == keyTrack || key == keyDisc {
					tg.addNumber(key, totalKey(key), first(id3Strings(enc, value, l)))
				} else if !tg.has(key) {
					tg.add(key, id3Strings(enc, value, l)...)
				}
			}
		case "COMM", "USLT":
			if len(rest) < 3 {
				continue
			}
			desc, text := splitTerminated(enc, rest[3:])
			d := id3String(enc, desc, l)
			if f.id == "COMM" && skipComment(d) {
				continue
			}
			if f.id == "COMM" {
				comment.offer(d, id3String(enc, text, l))
			} else {
				lyrics.offer(d, id3String(enc, text, l))
			}
		case "SYLT":
			if len(t.synced) == 0 && len(rest) > 6 {
				t.synced = syncedLyrics(enc, rest, l)
			}
		case "APIC":
			t.addAPIC(enc, rest, l, version)
		case "CHAP":
			t.addCHAP(f.data, version)
		}
	}
	tg.add(keyComment, comment.text)
	if lyrics.text != "" {
		tg.add(keyLyrics, lyrics.text)
	} else {
		tg.add(keyLyrics, lyricsText(t.synced))
	}
	switch {
	case date != "":
		tg.add(keyDate, date)
	case year != "":
		if len(ddmm) == 4 && len(year) == 4 {
			year += "-" + ddmm[2:] + "-" + ddmm[:2]
			if len(hhmm) == 4 {
				year += "T" + hhmm[:2] + ":" + hhmm[2:]
			}
		}
		tg.add(keyDate, year)
	case release != "":
		tg.add(keyDate, release)
	}
}

// described is a text that several frames may hold, each with a
// description: the one without a description is taken, else the first.
type described struct {
	text, desc string
	set        bool
}

func (d *described) offer(desc, text string) {
	if text != "" && (!d.set || desc == "" && d.desc != "") {
		d.text, d.desc, d.set = text, desc, true
	}
}

func first(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return s[0]
}

// splitPeople splits names joined with semicolons, as tag editors join
// several artists in one v2.3 frame.
func splitPeople(values []string) []string {
	var out []string
	for _, v := range values {
		for _, p := range strings.Split(v, ";") {
			if p = strings.TrimSpace(p); p != "" {
				out = append(out, p)
			}
		}
	}
	return out
}

// skipComment reports a comment that is a program's own data, not a
// comment on the music.
func skipComment(desc string) bool {
	for _, p := range []string{"iTun", "Songs-DB", "MusicMatch", "Media Jukebox", "Encoder", "Encoded by", "Encoding"} {
		if strings.HasPrefix(desc, p) {
			return true
		}
	}
	return false
}

// syncedLyrics reads the lines of a SYLT frame (after its encoding byte):
// language, time stamp format, content type, descriptor, then entries of a
// terminated string and a time stamp. The lines have their times when the
// stamps are in milliseconds; stamps that count MPEG frames are not
// converted, and leave the lines without times.
func syncedLyrics(enc byte, b []byte, l legacy) []lyricLine {
	if len(b) < 6 || b[4] != 1 { // the lyrics, not a transcription or the like
		return nil
	}
	timed := b[3] == 2
	_, b = splitTerminated(enc, b[5:])
	var entries []syltEntry
	for len(b) > 0 && len(entries) < 10000 {
		s, rest := splitTerminated(enc, b)
		if len(rest) < 4 {
			break
		}
		entries = append(entries, syltEntry{id3Text(enc, s, l), int64(binary.BigEndian.Uint32(rest))})
		b = rest[4:]
	}
	return syltLines(entries, timed)
}

// id3Text decodes a string of a frame in its encoding, keeping the spaces
// and line feeds around it (the pieces of synchronized lyrics).
func id3Text(enc byte, b []byte, l legacy) string {
	var s string
	switch enc {
	case 1:
		s = decodeUTF16(b, false)
	case 2:
		s = decodeUTF16(b, true)
	case 3:
		s = string(bytes.ToValidUTF8(b, []byte("\uFFFD")))
	default:
		s = l.decode(b)
	}
	s = strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\r", "\n")
	return strings.Map(func(r rune) rune {
		switch {
		case r == '\n':
			return r
		case r == '\t':
			return ' '
		case r < 0x20 || r == 0x7F || r >= 0x80 && r < 0xA0 || r == 0xFEFF:
			return -1
		}
		return r
	}, s)
}

// addAPIC reads a picture frame (after its encoding byte).
func (t *track) addAPIC(enc byte, b []byte, l legacy, version int) {
	var mime string
	if version == 2 {
		if len(b) < 4 {
			return
		}
		mime, b = "image/"+strings.ToLower(l.decode(b[:3])), b[3:]
		if mime == "image/jpg" {
			mime = "image/jpeg"
		}
	} else {
		var m []byte
		m, b = splitTerminated(0, b)
		mime = strings.ToLower(l.decode(m))
	}
	if len(b) < 1 {
		return
	}
	typ := int(b[0])
	_, data := splitTerminated(enc, b[1:])
	if len(data) == 0 || len(t.pictures) >= maxValues {
		return
	}
	t.pictures = append(t.pictures, picture{typ: typ, mime: mime, data: data})
}

// addCHAP reads a chapter frame: its element identifier, start and end
// times, offsets, and the frames in it (its title).
func (t *track) addCHAP(b []byte, version int) {
	id, rest := splitTerminated(0, b)
	if len(rest) < 16 {
		return
	}
	start := binary.BigEndian.Uint32(rest)
	title := ""
	for _, f := range parseFrames(rest[16:], max(version, 3), false) {
		if f.id == "TIT2" && len(f.data) > 1 {
			title = first(id3Strings(f.data[0], f.data[1:], guessLegacy([][]byte{f.data[1:]})))
			break
		}
	}
	if title == "" {
		title = string(id)
	}
	t.addChapter(float64(start)/1000, title)
}

// readID3v1 reads an ID3v1 tag (the last 128 bytes of a file) into tags.
func readID3v1(b []byte) (tags, bool) {
	if len(b) != 128 || string(b[:3]) != "TAG" {
		return nil, false
	}
	fields := [][]byte{b[3:33], b[33:63], b[63:93], b[97:127]}
	l := guessLegacy(fields)
	tg := tags{}
	tg.add(keyTitle, l.decode(fields[0]))
	tg.add(keyArtist, l.decode(fields[1]))
	tg.add(keyAlbum, l.decode(fields[2]))
	tg.add(keyDate, l.decode(b[93:97]))
	comment := b[97:127]
	if comment[28] == 0 && comment[29] != 0 {
		tg.add(keyTrack, fmt.Sprint(comment[29]))
		comment = comment[:28]
	}
	tg.add(keyComment, l.decode(comment))
	if g := int(b[127]); g < len(id3Genres) {
		tg.add(keyGenre, id3Genres[g])
	}
	return tg, true
}
