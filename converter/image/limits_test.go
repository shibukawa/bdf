package image

import (
	"bytes"
	"encoding/binary"
	"slices"
	"strings"
	"testing"
)

// Images whose metadata is stated over and over, or stands for more than
// the file holds: what the converter keeps of it is no larger than the file.

// isoBox makes an ISOBMFF box.
func isoBox(typ string, content ...[]byte) []byte {
	b := slices.Concat(content...)
	return slices.Concat(binary.BigEndian.AppendUint32(nil, uint32(8+len(b))), []byte(typ), b)
}

// TestAVIFItems reads an AVIF image whose item locations list the Exif item
// a hundred times, each time as the whole file.
func TestAVIFItems(t *testing.T) {
	u16 := func(v int) []byte { return binary.BigEndian.AppendUint16(nil, uint16(v)) }
	u32 := func(v int) []byte { return binary.BigEndian.AppendUint32(nil, uint32(v)) }
	full := []byte{0, 0, 0, 0} // version and flags
	var items []byte
	for range 100 {
		// the item, its construction method and data reference, then an
		// extent without offset and length: all of the file
		items = slices.Concat(items, u32(2), u16(0), u16(0), u16(1))
	}
	data := slices.Concat(
		isoBox("ftyp", []byte("avif"), u32(0), []byte("avifmif1")),
		isoBox("meta", full,
			isoBox("pitm", full, u16(1)),
			isoBox("iinf", full, u16(1), isoBox("infe", []byte{2, 0, 0, 0}, u16(2), u16(0), []byte("Exif\x00"))),
			isoBox("iprp", isoBox("ipco", isoBox("ispe", full, u32(64), u32(48))), isoBox("ipma", full, u32(1), u16(1), []byte{1, 1})),
			isoBox("iloc", []byte{2, 0, 0, 0}, []byte{0, 0}, u32(100), items)))
	p, err := readAVIF(data)
	if err != nil || p.w != 64 || p.h != 48 {
		t.Fatalf("%v, %+v", err, p)
	}
	size := 0
	for _, e := range p.exif {
		size += len(e)
	}
	if len(p.exif) != 1 || size > len(data) {
		t.Errorf("%d Exif blocks of %d bytes from a file of %d", len(p.exif), size, len(data))
	}
}

// TestGIFPackets reads the XMP packet of a GIF image from its application
// extension, not from what follows it.
func TestGIFPackets(t *testing.T) {
	head := "GIF89a\x01\x00\x01\x00\x00\x00\x00"
	packet := `<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?><x:xmpmeta xmlns:x="adobe:ns:meta/"/><?xpacket end="w"?>`
	// the trailer that ends the packet for a reader of sub-blocks, wherever
	// in the packet it stands: 1, 255 down to 0, and the terminator
	trailer := []byte{1}
	for i := 255; i >= 0; i-- {
		trailer = append(trailer, byte(i))
	}
	trailer = append(trailer, 0)
	ext := "\x21\xff\x0bXMP DataXMP"
	p, err := readGIF([]byte(head + ext + packet + string(trailer) + "\x3b"))
	if err != nil || len(p.xmp) != 1 || string(p.xmp[0]) != packet {
		t.Errorf("%v, packets %q", err, p.xmp)
	}
	// A hundred extensions without a packet, and the end of one after the
	// image: none of them has it.
	p, err = readGIF([]byte(head + strings.Repeat(ext+"\x00", 100) + "\x3b" + packet))
	if err != nil || len(p.xmp) != 0 {
		t.Errorf("%v, %d packets", err, len(p.xmp))
	}
}

// TestSVGEntities reads an SVG document whose title refers a hundred times
// to an entity of 50 KiB.
func TestSVGEntities(t *testing.T) {
	doc := `<!DOCTYPE svg [<!ENTITY a "` + strings.Repeat("A", 50<<10) + `"><!ENTITY ns "http://www.w3.org/2000/svg">]>` +
		`<svg xmlns="&ns;" width="8" height="4"><title>` + strings.Repeat("&a;", 100) + `</title></svg>`
	res, err := Convert(bytes.NewReader([]byte(doc)), int64(len(doc)), nil)
	if err != nil {
		t.Fatal(err)
	}
	want := "the entities of the SVG document stand for too much text: they are not replaced"
	if title := res.Doc.Meta.DC.Title.First(); len(title) > 300 || res.Width != 8 || !slices.Contains(res.Warnings, want) {
		t.Errorf("title of %d bytes, warnings %q", len(title), res.Warnings)
	}
	// Twenty references are within the limit.
	doc = strings.Replace(doc, strings.Repeat("&a;", 100), strings.Repeat("&a;", 20), 1)
	res, err = Convert(bytes.NewReader([]byte(doc)), int64(len(doc)), nil)
	if err != nil || len(res.Doc.Meta.DC.Title.First()) != 20*50<<10 || len(res.Warnings) != 0 {
		t.Errorf("%v, warnings %q", err, res.Warnings)
	}
}
