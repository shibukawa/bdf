package image

import (
	"bytes"
	"encoding/binary"

	"github.com/shibukawa/bdf"
)

// exifInfo is what an EXIF block (a TIFF structure) tells about an image.
type exifInfo struct {
	dc bdf.DublinCore
	// orientation is the Orientation tag (1–8), 0 when absent.
	orientation int
}

// EXIF tags read (EXIF 2.32; the XP tags are Windows').
const (
	tagImageDescription   = 0x010e
	tagOrientation        = 0x0112
	tagDateTime           = 0x0132
	tagArtist             = 0x013b
	tagCopyright          = 0x8298
	tagExifIFD            = 0x8769
	tagDateTimeOriginal   = 0x9003
	tagDateTimeDigitized  = 0x9004
	tagOffsetTime         = 0x9010
	tagOffsetTimeOriginal = 0x9011
	tagOffsetTimeDigital  = 0x9012
	tagUserComment        = 0x9286
	tagXPTitle            = 0x9c9b
	tagXPComment          = 0x9c9c
	tagXPAuthor           = 0x9c9d
	tagXPKeywords         = 0x9c9e
	tagXPSubject          = 0x9c9f
)

// tiff reads the IFDs of a TIFF structure.
type tiff struct {
	b     []byte
	order binary.ByteOrder
}

// ifd reads the entries of the IFD at off: the raw value of each tag.
func (t *tiff) ifd(off uint32) map[uint16][]byte {
	b := t.b
	if int64(off)+2 > int64(len(b)) {
		return nil
	}
	n := int(t.order.Uint16(b[off:]))
	out := map[uint16][]byte{}
	for i := range n {
		e := int64(off) + 2 + 12*int64(i)
		if e+12 > int64(len(b)) {
			break
		}
		tag := t.order.Uint16(b[e:])
		typ := t.order.Uint16(b[e+2:])
		count := int64(t.order.Uint32(b[e+4:]))
		var size int64
		switch typ {
		case 1, 2, 6, 7: // BYTE, ASCII, SBYTE, UNDEFINED
			size = 1
		case 3, 8: // SHORT, SSHORT
			size = 2
		case 4, 9, 11: // LONG, SLONG, FLOAT
			size = 4
		case 5, 10, 12: // RATIONAL, SRATIONAL, DOUBLE
			size = 8
		default:
			continue
		}
		total := count * size
		if total <= 4 {
			out[tag] = b[e+8 : e+8+total]
			continue
		}
		at := int64(t.order.Uint32(b[e+8:]))
		if total > int64(len(b)) || at+total > int64(len(b)) {
			continue
		}
		out[tag] = b[at : at+total]
	}
	return out
}

func (t *tiff) short(v []byte) int {
	if len(v) < 2 {
		return 0
	}
	return int(t.order.Uint16(v))
}

func (t *tiff) long(v []byte) uint32 {
	if len(v) < 4 {
		return 0
	}
	return t.order.Uint32(v)
}

// parseEXIF reads an EXIF block: a TIFF header and its IFDs. A leading
// "Exif\0\0" (the JPEG APP1 identifier, which some writers keep in PNG and
// WebP) is skipped.
func parseEXIF(b []byte) (exifInfo, bool) {
	b = bytes.TrimPrefix(b, []byte("Exif\x00\x00"))
	var info exifInfo
	if len(b) < 8 {
		return info, false
	}
	t := &tiff{b: b}
	switch string(b[:4]) {
	case "II*\x00":
		t.order = binary.LittleEndian
	case "MM\x00*":
		t.order = binary.BigEndian
	default:
		return info, false
	}
	ifd0 := t.ifd(t.order.Uint32(b[4:]))
	if ifd0 == nil {
		return info, false
	}
	var ex map[uint16][]byte
	if v, ok := ifd0[tagExifIFD]; ok {
		ex = t.ifd(t.long(v))
	}
	if v, ok := ifd0[tagOrientation]; ok {
		if o := t.short(v); o >= 1 && o <= 8 {
			info.orientation = o
		}
	}
	dc := &info.dc
	xp := func(tag uint16) string { return utf16Text(ifd0[tag], false) }
	add(&dc.Title, xp(tagXPTitle))
	if dc.Creator = splitList(text(ifd0[tagArtist])); len(dc.Creator) == 0 {
		dc.Creator = splitList(xp(tagXPAuthor))
	}
	dc.Subject = splitList(xp(tagXPKeywords))
	for _, d := range []string{text(ifd0[tagImageDescription]), xp(tagXPSubject), xp(tagXPComment), userComment(ex[tagUserComment], t.order)} {
		if !placeholder(d) {
			add(&dc.Description, d)
			break
		}
	}
	// the photographer's copyright, a NUL, then the editor's
	for _, r := range bytes.Split(ifd0[tagCopyright], []byte{0}) {
		add(&dc.Rights, text(r))
	}
	if d := exifDate(text(ex[tagDateTimeOriginal]), text(ex[tagOffsetTimeOriginal])); d != "" {
		add(&dc.Created, d)
	} else {
		add(&dc.Created, exifDate(text(ex[tagDateTimeDigitized]), text(ex[tagOffsetTimeDigital])))
	}
	add(&dc.Modified, exifDate(text(ifd0[tagDateTime]), text(ex[tagOffsetTime])))
	return info, true
}

// userComment decodes EXIF's UserComment: an 8-byte character code, then
// the text. Only ASCII and Unicode (UCS-2 in the byte order of the TIFF)
// comments are read.
func userComment(v []byte, order binary.ByteOrder) string {
	if len(v) < 8 {
		return ""
	}
	switch string(v[:8]) {
	case "ASCII\x00\x00\x00":
		return text(v[8:])
	case "UNICODE\x00":
		return utf16Text(v[8:], order == binary.BigEndian)
	}
	return ""
}
