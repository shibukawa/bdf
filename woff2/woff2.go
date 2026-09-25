// Package woff2 packs TrueType and CFF-flavored OpenType fonts into WOFF2
// (W3C WOFF File Format 2.0), the web font format browsers load with FontFace.
//
// glyf and loca are stored with the WOFF2 glyf transform (the decoder rebuilds
// both from separated point, flag and instruction streams), every other table
// passes through, and the whole table stream is Brotli-compressed. The Brotli
// encoder is left out of builds tagged bdf_noconv; Encode then returns
// ErrNotAvailable and callers keep the plain sfnt.
package woff2

import (
	"encoding/binary"
	"errors"
	"sort"
)

// ErrNotAvailable is returned by Encode when the Brotli encoder is not compiled in.
var ErrNotAvailable = errors.New("woff2: encoder not available in this build")

// IsWOFF reports whether data is a WOFF or WOFF2 file (already compressed).
func IsWOFF(data []byte) bool {
	return len(data) >= 4 && (string(data[:4]) == "wOF2" || string(data[:4]) == "wOFF")
}

// Encode converts an sfnt font (TrueType outlines or CFF) into WOFF2.
func Encode(font []byte) ([]byte, error) {
	if !Available() {
		return nil, ErrNotAvailable
	}
	if len(font) < 12 {
		return nil, errors.New("woff2: font too short")
	}
	flavor := binary.BigEndian.Uint32(font)
	switch flavor {
	case 0x00010000, 0x74727565, 0x4f54544f: // 1.0, 'true', 'OTTO'
	default:
		return nil, errors.New("woff2: not a TrueType or OpenType font")
	}
	n := int(binary.BigEndian.Uint16(font[4:]))
	tables := map[string][]byte{}
	for i := 0; i < n; i++ {
		rec := 12 + i*16
		if rec+16 > len(font) {
			return nil, errors.New("woff2: truncated table directory")
		}
		tag := string(font[rec : rec+4])
		off, length := int(binary.BigEndian.Uint32(font[rec+8:])), int(binary.BigEndian.Uint32(font[rec+12:]))
		if off < 0 || length < 0 || off+length > len(font) || off+length < off {
			return nil, errors.New("woff2: table " + tag + " out of range")
		}
		tables[tag] = font[off : off+length]
	}
	tags := make([]string, 0, len(tables))
	for t := range tables {
		tags = append(tags, t)
	}
	// Tables are listed by tag, which also puts loca after glyf as the format requires.
	sort.Strings(tags)

	// The glyf transform is lossless for well-formed glyphs; anything it cannot
	// represent (cubic or malformed outlines) keeps both tables as they are.
	var glyfT []byte
	transformed, locaLen := false, 0
	if head, maxp := tables["head"], tables["maxp"]; flavor != 0x4f54544f && tables["glyf"] != nil && tables["loca"] != nil && len(head) >= 54 && len(maxp) >= 6 {
		numGlyphs := int(binary.BigEndian.Uint16(maxp[4:]))
		long := binary.BigEndian.Uint16(head[50:]) != 0
		if t, ok := transformGlyf(tables["glyf"], tables["loca"], numGlyphs, long); ok {
			glyfT, transformed = t, true
			// The decoder rebuilds loca in the long format (it costs nothing on the
			// wire) and head must agree with it; bit 11 marks the transform.
			h := append([]byte(nil), head...)
			binary.BigEndian.PutUint16(h[50:], 1)
			binary.BigEndian.PutUint16(h[16:], binary.BigEndian.Uint16(h[16:])|1<<11)
			tables["head"] = h
			locaLen = (numGlyphs + 1) * 4
		}
	}

	var dir, stream []byte
	totalSfnt := 12 + 16*len(tags)
	for _, tag := range tags {
		data := tables[tag]
		origLen := len(data)
		if transformed && tag == "loca" {
			origLen = locaLen
		}
		totalSfnt += (origLen + 3) &^ 3
		flags := byte(knownTag(tag))
		if (tag == "glyf" || tag == "loca") && !transformed {
			flags |= 3 << 6 // null transform
		}
		dir = append(dir, flags)
		if flags&0x3f == 0x3f {
			dir = append(dir, tag...)
		}
		dir = appendBase128(dir, uint32(origLen))
		switch {
		case transformed && tag == "glyf":
			dir = appendBase128(dir, uint32(len(glyfT)))
			stream = append(stream, glyfT...)
		case transformed && tag == "loca": // rebuilt by the decoder, nothing stored
			dir = appendBase128(dir, 0)
		default:
			stream = append(stream, data...)
		}
	}
	compressed, err := compress(stream)
	if err != nil {
		return nil, err
	}
	out := make([]byte, 48, 48+len(dir)+len(compressed)+3)
	out = append(out, dir...)
	out = append(out, compressed...)
	for len(out)%4 != 0 {
		out = append(out, 0)
	}
	if totalSfnt < len(out) {
		totalSfnt = len(out)
	}
	copy(out, "wOF2")
	binary.BigEndian.PutUint32(out[4:], flavor)
	binary.BigEndian.PutUint32(out[8:], uint32(len(out)))
	binary.BigEndian.PutUint16(out[12:], uint16(len(tags)))
	binary.BigEndian.PutUint32(out[16:], uint32(totalSfnt))
	binary.BigEndian.PutUint32(out[20:], uint32(len(compressed)))
	binary.BigEndian.PutUint16(out[24:], 1) // font version 1.0; no metadata or private block
	return out, nil
}

// knownTags is the WOFF2 known-table-tag list; a tag's index is its flag value.
var knownTags = [...]string{
	"cmap", "head", "hhea", "hmtx", "maxp", "name", "OS/2", "post", "cvt ", "fpgm", "glyf", "loca", "prep", "CFF ", "VORG", "EBDT",
	"EBLC", "gasp", "hdmx", "kern", "LTSH", "PCLT", "VDMX", "vhea", "vmtx", "BASE", "GDEF", "GPOS", "GSUB", "EBSC", "JSTF", "MATH",
	"CBDT", "CBLC", "COLR", "CPAL", "SVG ", "sbix", "acnt", "avar", "bdat", "bloc", "bsln", "cvar", "fdsc", "feat", "fmtx", "fvar",
	"gvar", "hsty", "just", "lcar", "mort", "morx", "opbd", "prop", "trak", "Zapf", "Silf", "Glat", "Gloc", "Feat", "Sill",
}

func knownTag(tag string) int {
	for i, t := range knownTags {
		if t == tag {
			return i
		}
	}
	return 0x3f
}

// appendBase128 writes a UIntBase128: big-endian 7-bit groups, high bit set on all but the last.
func appendBase128(b []byte, v uint32) []byte {
	var tmp [5]byte
	i := len(tmp) - 1
	tmp[i] = byte(v & 0x7f)
	for v >>= 7; v != 0; v >>= 7 {
		i--
		tmp[i] = byte(v&0x7f) | 0x80
	}
	return append(b, tmp[i:]...)
}

// append255 writes a 255UInt16.
func append255(b []byte, v int) []byte {
	switch {
	case v < 253:
		return append(b, byte(v))
	case v < 506:
		return append(b, 255, byte(v-253))
	case v < 762:
		return append(b, 254, byte(v-506))
	default:
		return append(b, 253, byte(v>>8), byte(v))
	}
}
