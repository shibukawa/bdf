package audio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/shibukawa/bdf/converter/internal/isobmff"
)

// Bounds on an MP4 file: the movie box (which holds the sample tables of
// every track and the metadata), and an XMP box.
const (
	maxMoov = 64 << 20
	maxXMP  = 16 << 20
)

// xmpUUID identifies the uuid box that holds an XMP packet.
var xmpUUID = []byte{0xBE, 0x7A, 0xCF, 0xCB, 0x97, 0xA9, 0x42, 0xE8, 0x9C, 0x71, 0x99, 0x94, 0x91, 0xE3, 0xAF, 0xAC}

var be = binary.BigEndian

// imageBrands are the brands of HEIF image files (AVIF, HEIC), which are
// not audio.
var imageBrands = map[string]bool{"avif": true, "avis": true, "heic": true, "heix": true, "hevc": true, "hevx": true, "mif1": true, "msf1": true, "miaf": true}

// mp4Kind tells what an ISOBMFF file holds by its tracks: "audio" for
// audio tracks only, "video" when it has a video track, "image" for a
// HEIF image, "" for a file without tracks (or one that cannot be read).
func mp4Kind(head []byte, r io.ReaderAt, size int64) string {
	if len(head) < 12 || string(head[4:8]) != "ftyp" {
		return ""
	}
	n := min(len(head), int(be.Uint32(head)))
	for i := 8; i+4 <= n; i += 4 {
		if i != 12 && imageBrands[string(head[i:i+4])] {
			return "image"
		}
	}
	for _, top := range isobmff.Walk(r, size) {
		if top.Type != "moov" || top.Size-top.Header > maxMoov {
			continue
		}
		moov, err := top.Contents(r)
		if err != nil {
			return ""
		}
		audio, video := false, false
		for _, trak := range isobmff.Boxes(moov) {
			if trak.Type != "trak" {
				continue
			}
			switch handler(trak.Data) {
			case "soun":
				audio = true
			case "vide":
				video = true
			}
		}
		switch {
		case video:
			return "video"
		case audio:
			return "audio"
		}
		return ""
	}
	return ""
}

// handler returns the handler type of a track ("soun", "vide", "text").
func handler(trak []byte) string {
	mdia := isobmff.Boxes(isobmff.Find(isobmff.Boxes(trak), "mdia"))
	if hdlr := isobmff.Find(mdia, "hdlr"); len(hdlr) >= 12 {
		return string(hdlr[8:12])
	}
	return ""
}

// readMP4 reads an MP4 (M4A) file: its first audio track, its duration,
// and the iTunes-style metadata, chapters and XMP of the movie.
func readMP4(r io.ReaderAt, size int64) (*track, error) {
	t := newTrack("m4a")
	var moov []byte
	var mdat int64
	for _, top := range isobmff.Walk(r, size) {
		switch top.Type {
		case "moov":
			if moov != nil {
				continue
			}
			if top.Size-top.Header > maxMoov {
				return nil, fmt.Errorf("the movie box is larger than %d bytes", maxMoov)
			}
			var err error
			if moov, err = top.Contents(r); err != nil {
				return nil, err
			}
		case "mdat":
			mdat += top.Size - top.Header
		case "uuid":
			if top.Size-top.Header < 16 || top.Size-top.Header > maxXMP {
				continue
			}
			b, err := top.Contents(r)
			if err == nil && bytes.Equal(b[:16], xmpUUID) {
				t.xmp = append(t.xmp, b[16:])
			}
		}
	}
	if moov == nil {
		return nil, errors.New("the file has no movie box")
	}
	boxes := isobmff.Boxes(moov)
	var timescale, duration uint64
	if mvhd := isobmff.Find(boxes, "mvhd"); len(mvhd) >= 24 {
		if mvhd[0] == 1 && len(mvhd) >= 32 {
			timescale, duration = uint64(be.Uint32(mvhd[20:])), be.Uint64(mvhd[24:])
		} else {
			timescale, duration = uint64(be.Uint32(mvhd[12:])), uint64(be.Uint32(mvhd[16:]))
			if duration == math.MaxUint32 {
				duration = 0
			}
		}
	}
	found := false
	for _, trak := range boxes {
		if trak.Type != "trak" || handler(trak.Data) != "soun" {
			continue
		}
		if found {
			t.warn("the file has more than one audio track: the first is described")
			break
		}
		found = true
		t.readSoundTrack(trak.Data, timescale, duration)
	}
	if !found {
		return nil, errors.New("the file has no audio track")
	}
	if timescale > 0 && duration > 0 {
		t.duration = float64(duration) / float64(timescale)
	} else if mehd := isobmff.Find(isobmff.Boxes(isobmff.Find(boxes, "mvex")), "mehd"); len(mehd) >= 8 && timescale > 0 {
		if mehd[0] == 1 && len(mehd) >= 12 {
			t.duration = float64(be.Uint64(mehd[4:])) / float64(timescale)
		} else {
			t.duration = float64(be.Uint32(mehd[4:])) / float64(timescale)
		}
	}
	if t.bitrate == 0 && t.duration > 0 && mdat > 0 {
		t.bitrate = int(float64(mdat)*8/t.duration + 0.5)
	}
	// metadata: the movie's user data, and iTunes' metadata under it (or
	// under the movie itself)
	for _, parent := range [][]byte{isobmff.Find(boxes, "udta"), moov} {
		if parent == nil {
			continue
		}
		children := isobmff.Boxes(parent)
		if meta := isobmff.Find(children, "meta"); meta != nil {
			t.readMeta(meta)
		}
		if parent != nil && &parent[0] != &moov[0] {
			t.readUserData(children)
		}
	}
	return t, nil
}

// readSoundTrack reads the sample description of a sound track: the
// coding, sample rate and channels; and its language and duration.
func (t *track) readSoundTrack(trak []byte, mvTimescale, mvDuration uint64) {
	mdia := isobmff.Boxes(isobmff.Find(isobmff.Boxes(trak), "mdia"))
	var timescale, duration uint64
	if mdhd := isobmff.Find(mdia, "mdhd"); len(mdhd) >= 22 {
		lang := 20
		if mdhd[0] == 1 && len(mdhd) >= 34 {
			timescale, duration, lang = uint64(be.Uint32(mdhd[20:])), be.Uint64(mdhd[24:]), 32
		} else {
			timescale, duration = uint64(be.Uint32(mdhd[12:])), uint64(be.Uint32(mdhd[16:]))
		}
		if code := be.Uint16(mdhd[lang:]); code != 0 && code != 0x7FFF {
			t.tags.add(keyLanguage, string([]byte{byte(code>>10&31) + 0x60, byte(code>>5&31) + 0x60, byte(code&31) + 0x60}))
		}
	}
	if mvTimescale == 0 || mvDuration == 0 {
		if timescale > 0 && duration > 0 && duration != math.MaxUint32 {
			t.duration = float64(duration) / float64(timescale)
		}
	}
	stbl := isobmff.Boxes(isobmff.Find(isobmff.Boxes(isobmff.Find(mdia, "minf")), "stbl"))
	stsd := isobmff.Find(stbl, "stsd")
	if len(stsd) < 8 {
		return
	}
	entries := isobmff.Boxes(stsd[8:])
	if len(entries) == 0 {
		return
	}
	e := entries[0]
	t.codec = sampleEntryCodec(e.Type)
	switch e.Type {
	case "drms", "drmi":
		t.protected = true
	}
	b := e.Data
	if len(b) < 28 {
		return
	}
	version := be.Uint16(b[8:])
	children := 36 // version 0: 8 reserved, 20 of fields, then the boxes
	switch version {
	case 1:
		children = 52
	case 2:
		children = 72
	}
	if version == 2 {
		if len(b) >= 44 {
			t.sampleRate = int(math.Float64frombits(be.Uint64(b[28:])) + 0.5)
			t.channels = int(be.Uint32(b[36:]))
		}
	} else {
		t.channels = int(be.Uint16(b[16:]))
		t.bits = int(be.Uint16(b[18:]))
		t.sampleRate = int(be.Uint32(b[24:]) >> 16)
	}
	if t.sampleRate == 0 && timescale > 0 && timescale < 1<<20 {
		t.sampleRate = int(timescale)
	}
	if children > len(b) {
		return
	}
	for _, c := range isobmff.Boxes(b[children:]) {
		switch c.Type {
		case "esds":
			t.readESDS(c.Data)
		case "alac":
			if len(c.Data) >= 36 {
				t.bits = int(c.Data[9])
				t.channels = int(c.Data[13])
				t.sampleRate = int(be.Uint32(c.Data[32:]))
			}
		case "dfLa":
			t.bits = 0
			if len(c.Data) >= 4+4+18 {
				si := c.Data[8:]
				t.sampleRate = int(be.Uint32(si[10:]) >> 12)
				t.channels = int(si[12]>>1&7) + 1
				t.bits = int(si[12]&1)<<4 | int(si[13]>>4) + 1
			}
		case "dOps":
			if len(c.Data) >= 2 {
				t.channels = int(c.Data[1])
			}
			t.sampleRate = 48000
		}
	}
	switch t.codec {
	case "PCM", "Apple Lossless", "FLAC":
	default:
		t.bits = 0
	}
}

// sampleEntryCodec names the coding of a sample entry type.
func sampleEntryCodec(typ string) string {
	switch typ {
	case "mp4a", "drms":
		return "AAC"
	case "alac":
		return "Apple Lossless"
	case "fLaC":
		return "FLAC"
	case "Opus":
		return "Opus"
	case "ac-3", "sac3":
		return "AC-3"
	case "ec-3":
		return "E-AC-3"
	case ".mp3", "mp3 ":
		return "MP3"
	case "samr":
		return "AMR"
	case "sawb":
		return "AMR-WB"
	case "lpcm", "sowt", "twos", "raw ", "in24", "in32", "fl32", "fl64", "NONE":
		return "PCM"
	case "ulaw":
		return "μ-law"
	case "alaw":
		return "A-law"
	case "ima4":
		return "IMA ADPCM"
	case "mp4v", "avc1", "hvc1", "hev1", "vp09", "av01":
		return ""
	}
	return strings.TrimSpace(typ)
}

// readESDS reads the elementary stream descriptor of an MPEG-4 audio
// entry: the object type (AAC, MP3) and the average bitrate.
func (t *track) readESDS(b []byte) {
	if len(b) < 4 {
		return
	}
	d := descriptor(b[4:], 0x03) // ES_Descriptor
	if len(d) < 3 {
		return
	}
	p := 3
	flags := d[2]
	if flags&0x80 != 0 {
		p += 2
	}
	if flags&0x40 != 0 && p < len(d) {
		p += 1 + int(d[p])
	}
	if flags&0x20 != 0 {
		p += 2
	}
	if p > len(d) {
		return
	}
	dc := descriptor(d[p:], 0x04) // DecoderConfigDescriptor
	if len(dc) < 13 {
		return
	}
	oti := dc[0]
	if avg := be.Uint32(dc[9:]); avg > 0 {
		t.bitrate = int(avg)
	}
	switch oti {
	case 0x40: // MPEG-4 audio: the object type says which
		t.codec = "AAC"
		if dsi := descriptor(dc[13:], 0x05); len(dsi) >= 2 {
			aot := int(dsi[0] >> 3)
			if aot == 31 {
				aot = 32 + int(dsi[0]&7)<<3 | int(dsi[1]>>5)
			}
			switch aot {
			case 5:
				t.codec = "HE-AAC"
			case 29:
				t.codec = "HE-AAC v2"
			case 23, 39:
				t.codec = "AAC LD"
			case 42:
				t.codec = "xHE-AAC"
			}
		}
	case 0x66, 0x67, 0x68:
		t.codec = "AAC"
	case 0x69, 0x6B:
		t.codec = "MP3"
	case 0xA5:
		t.codec = "AC-3"
	case 0xA6:
		t.codec = "E-AC-3"
	case 0xA9, 0xAC:
		t.codec = "DTS"
	case 0xAD:
		t.codec = "Opus"
	case 0xE1:
		t.codec = "QCELP"
	}
}

// descriptor returns the payload of the first descriptor of a tag in b
// (MPEG-4 descriptors: a tag byte and a length in up to four bytes of
// seven bits), or nil.
func descriptor(b []byte, tag byte) []byte {
	for len(b) >= 2 {
		typ := b[0]
		p := 1
		size := 0
		for i := 0; i < 4 && p < len(b); i++ {
			c := b[p]
			p++
			size = size<<7 | int(c&0x7F)
			if c&0x80 == 0 {
				break
			}
		}
		if size > len(b)-p {
			size = len(b) - p
		}
		if typ == tag {
			return b[p : p+size]
		}
		b = b[p+size:]
	}
	return nil
}

// readMeta reads a meta box: iTunes' metadata list, with the keys box of
// QuickTime's metadata when it has one.
func (t *track) readMeta(meta []byte) {
	// a full box in MP4, a plain box in QuickTime: the first child tells
	children := isobmff.Boxes(meta)
	if len(children) == 0 || children[0].Type != "hdlr" && children[0].Type != "keys" && children[0].Type != "ilst" {
		if len(meta) < 4 {
			return
		}
		children = isobmff.Boxes(meta[4:])
	}
	var keys []string
	if kb := isobmff.Find(children, "keys"); len(kb) >= 8 {
		for _, k := range isobmff.Boxes(kb[8:]) {
			if k.Type == "mdta" {
				keys = append(keys, string(k.Data))
			} else {
				keys = append(keys, "")
			}
		}
	}
	for _, item := range isobmff.Boxes(isobmff.Find(children, "ilst")) {
		name := itemName(item.Type)
		parts := isobmff.Boxes(item.Data)
		if keys != nil {
			if i := int(be.Uint32([]byte(name))); i >= 1 && i <= len(keys) {
				name = keys[i-1]
			}
		}
		if name == "----" {
			name = "----:" + string(isobmff.Find(parts, "name"))
			if len(name) > 4+4 {
				name = "----:" + string(isobmff.Find(parts, "name")[4:])
			}
		}
		for _, data := range isobmff.FindAll(parts, "data") {
			if len(data) < 8 {
				continue
			}
			t.addItem(name, be.Uint32(data)&0xFFFFFF, data[8:])
		}
	}
}

// itemName normalizes the name of a metadata item: the copyright sign
// that starts many of them is the byte 0xA9 in the file.
func itemName(typ string) string {
	if len(typ) == 4 && typ[0] == 0xA9 {
		return "©" + typ[1:]
	}
	return typ
}

// addItem reads one data atom of an iTunes metadata item: text (type 1
// UTF-8, 2 UTF-16), a picture (13 JPEG, 14 PNG, 27 BMP) or a number.
func (t *track) addItem(name string, typ uint32, b []byte) {
	tg := t.tags
	text := func() string {
		switch typ {
		case 2:
			return decodeUTF16(b, true)
		case 1, 0:
			return string(bytes.ToValidUTF8(b, []byte("�")))
		}
		return ""
	}
	number := func() int {
		switch len(b) {
		case 1:
			return int(b[0])
		case 2:
			return int(be.Uint16(b))
		case 4:
			return int(be.Uint32(b))
		case 8:
			return int(be.Uint64(b))
		}
		return 0
	}
	switch name {
	case "©nam", "com.apple.quicktime.title":
		tg.add(keyTitle, text())
	case "©st3":
		tg.add(keySubtitle, text())
	case "©ART", "com.apple.quicktime.artist":
		tg.add(keyArtist, text())
	case "aART":
		tg.add(keyAlbumArtist, text())
	case "©alb", "com.apple.quicktime.album":
		tg.add(keyAlbum, text())
	case "©wrt", "com.apple.quicktime.composer":
		tg.add(keyComposer, text())
	case "©day", "com.apple.quicktime.year", "com.apple.quicktime.creationdate":
		if !tg.has(keyDate) {
			tg.add(keyDate, text())
		}
	case "©gen", "com.apple.quicktime.genre":
		tg.add(keyGenre, text())
	case "gnre":
		if n := number(); n >= 1 && n <= len(id3Genres) {
			tg.add(keyGenre, id3Genres[n-1])
		}
	case "trkn", "disk":
		if len(b) >= 6 {
			key, total := keyTrack, keyTrackTotal
			if name == "disk" {
				key, total = keyDisc, keyDiscTotal
			}
			if n := be.Uint16(b[2:]); n > 0 {
				tg.add(key, fmt.Sprint(n))
			}
			if n := be.Uint16(b[4:]); n > 0 {
				tg.add(total, fmt.Sprint(n))
			}
		}
	case "©cmt", "com.apple.quicktime.comment":
		tg.add(keyComment, text())
	case "desc", "ldes", "com.apple.quicktime.description":
		if !tg.has(keyComment) {
			tg.add(keyComment, text())
		}
	case "©lyr":
		tg.add(keyLyrics, text())
	case "cprt", "©cpy", "com.apple.quicktime.copyright":
		tg.add(keyCopyright, text())
	case "©grp":
		tg.add(keyGrouping, text())
	case "©pub":
		tg.add(keyPublisher, text())
	case "tmpo":
		if n := number(); n > 0 {
			tg.add(keyBPM, fmt.Sprint(n))
		}
	case "covr":
		mime := ""
		switch typ {
		case 13:
			mime = "image/jpeg"
		case 14:
			mime = "image/png"
		case 27:
			mime = "image/bmp"
		}
		if len(b) > 0 && len(t.pictures) < maxValues {
			t.pictures = append(t.pictures, picture{typ: 3, mime: mime, data: b})
		}
	default:
		if strings.HasPrefix(name, "----:") {
			if key := userKey(name[5:]); key != "" && !tg.has(key) {
				tg.add(key, text())
			}
		}
	}
}

// readUserData reads what a udta box holds outside meta: QuickTime's
// text atoms (a size, a language and the text) and Nero's chapter list.
func (t *track) readUserData(children []isobmff.Box) {
	for _, c := range children {
		switch name := itemName(c.Type); name {
		case "chpl":
			t.readChapterList(c.Data)
		case "©nam", "©ART", "©alb", "©day", "©cmt", "©gen", "©wrt", "©cpy", "©lyr", "©grp":
			if len(c.Data) >= 4 && int(be.Uint16(c.Data)) == len(c.Data)-4 {
				t.addItem(name, 1, c.Data[4:])
			}
		}
	}
}

// readChapterList reads a chpl box: chapters as start times in units of
// 100 ns and titles.
func (t *track) readChapterList(b []byte) {
	if len(b) < 5 {
		return
	}
	p := 4
	if b[0] == 1 {
		p += 4
	}
	if p >= len(b) {
		return
	}
	n := int(b[p])
	p++
	for i := 0; i < n && p+9 <= len(b); i++ {
		start := be.Uint64(b[p:])
		size := int(b[p+8])
		p += 9
		if p+size > len(b) {
			break
		}
		t.addChapter(float64(start)/1e7, string(bytes.ToValidUTF8(b[p:p+size], []byte("�"))))
		p += size
	}
}
