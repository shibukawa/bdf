package image

import (
	"bytes"
	"encoding/binary"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
)

// photoshopIPTC returns the IPTC-IIM block of Photoshop image resources
// (the body of a JPEG APP13 segment after "Photoshop 3.0\0"): the resource
// 0x0404.
func photoshopIPTC(b []byte) []byte {
	for len(b) >= 12 && string(b[:4]) == "8BIM" {
		id := binary.BigEndian.Uint16(b[4:])
		// a Pascal string name, padded to an even length
		nameLen := 1 + int(b[6])
		nameLen += nameLen & 1
		at := 6 + nameLen
		if at+4 > len(b) {
			return nil
		}
		size := int(binary.BigEndian.Uint32(b[at:]))
		at += 4
		if size < 0 || at+size > len(b) {
			return nil
		}
		if id == 0x0404 {
			return b[at : at+size]
		}
		b = b[min(len(b), at+size+size&1):] // data padded to an even length
	}
	return nil
}

// parseIPTC reads the application record of IPTC-IIM data (IPTC IIM 4.2):
// the object name, keywords, creation date and time, by-lines, copyright
// notice and caption. Text is UTF-8 when the envelope says so (or when it
// is valid UTF-8), Latin-1 otherwise.
func parseIPTC(b []byte) bdf.DublinCore {
	type dataset struct {
		record, id byte
		data       []byte
	}
	var sets []dataset
	for len(b) >= 5 && b[0] == 0x1c {
		size := int(binary.BigEndian.Uint16(b[3:]))
		if size&0x8000 != 0 || 5+size > len(b) { // extended lengths are not used for text
			break
		}
		sets = append(sets, dataset{b[1], b[2], b[5 : 5+size]})
		b = b[5+size:]
	}
	isUTF8 := false
	for _, s := range sets {
		if s.record == 1 && s.id == 90 && bytes.Equal(s.data, []byte("\x1b%G")) {
			isUTF8 = true
		}
	}
	str := func(v []byte) string {
		if isUTF8 || utf8.Valid(v) {
			return clean(string(bytes.ToValidUTF8(v, nil)))
		}
		return clean(latin1(v))
	}
	var dc bdf.DublinCore
	var date, tm, headline string
	for _, s := range sets {
		if s.record != 2 {
			continue
		}
		v := str(s.data)
		switch s.id {
		case 5: // object name
			if len(dc.Title) == 0 {
				add(&dc.Title, v)
			}
		case 25: // keywords
			add(&dc.Subject, v)
		case 55: // date created, CCYYMMDD
			date = v
		case 60: // time created, HHMMSS±HHMM
			tm = v
		case 80: // by-line
			add(&dc.Creator, v)
		case 105: // headline
			headline = v
		case 116: // copyright notice
			if len(dc.Rights) == 0 {
				add(&dc.Rights, v)
			}
		case 120: // caption / abstract
			if len(dc.Description) == 0 && !placeholder(v) {
				add(&dc.Description, v)
			}
		}
	}
	if len(dc.Title) == 0 {
		add(&dc.Title, headline)
	}
	add(&dc.Created, iptcDate(date, tm))
	return dc
}

// iptcDate converts IPTC's date (CCYYMMDD) and time (HHMMSS±HHMM) to the
// W3C format.
func iptcDate(date, tm string) string {
	if len(date) != 8 {
		return ""
	}
	d := date[:4] + "-" + date[4:6] + "-" + date[6:]
	if len(tm) == 11 {
		d += "T" + tm[:2] + ":" + tm[2:4] + ":" + tm[4:6] + tm[6:9] + ":" + tm[9:]
	} else if len(tm) == 6 {
		d += "T" + tm[:2] + ":" + tm[2:4] + ":" + tm[4:6]
	}
	return isoDate(d)
}
