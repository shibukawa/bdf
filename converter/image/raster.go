package image

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/shibukawa/bdf"
)

// picture is what the reader of a format found in an image.
type picture struct {
	format string
	// w and h are the size in CSS pixels as browsers draw the image, before
	// its EXIF orientation (which Convert applies).
	w, h     float64
	animated bool
	// Metadata blocks, and what the format itself says (PNG text chunks,
	// SVG title and description).
	xmp, exif [][]byte
	iptc      []byte
	native    bdf.DublinCore
	warnings  []string
}

var (
	be = binary.BigEndian
	le = binary.LittleEndian
)

// errTruncated is returned for images that end before their header does.
var errTruncated = errors.New("the image is truncated")

// maxText bounds the decompressed size of the compressed text chunks of an
// image.
const maxText = 16 << 20

// pngKeys are the keywords of the PNG text chunks read.
var pngKeys = map[string]bool{
	"XML:com.adobe.xmp": true, "Raw profile type exif": true, "Raw profile type APP1": true, "Raw profile type xmp": true,
	"Title": true, "Author": true, "Description": true, "Comment": true, "Copyright": true, "Creation Time": true,
}

// readPNG reads the header, text chunks (tEXt, zTXt, iTXt), eXIf and tIME
// of a PNG image, and whether it is animated (APNG).
func readPNG(b []byte) (*picture, error) {
	if len(b) < 33 || string(b[12:16]) != "IHDR" {
		return nil, errTruncated
	}
	p := &picture{format: "png", w: float64(be.Uint32(b[16:])), h: float64(be.Uint32(b[20:]))}
	var comment string
	budget := maxText // for decompressing text
	for off := 8; off+12 <= len(b); {
		n := int64(be.Uint32(b[off:]))
		typ := string(b[off+4 : off+8])
		if int64(off)+12+n > int64(len(b)) {
			break
		}
		data := b[off+8 : off+8+int(n)]
		off += 12 + int(n)
		switch typ {
		case "tEXt", "zTXt", "iTXt":
			key, value, ok := pngText(typ, data, &budget)
			if !ok {
				continue
			}
			switch key {
			case "XML:com.adobe.xmp":
				p.xmp = append(p.xmp, []byte(value))
			case "Raw profile type exif", "Raw profile type APP1":
				// ImageMagick's hex dump: "\nexif\n    1234\n4578…"
				if raw := rawProfile(value); bytes.HasPrefix(raw, []byte("Exif\x00\x00")) || bytes.HasPrefix(raw, []byte("II*\x00")) || bytes.HasPrefix(raw, []byte("MM\x00*")) {
					p.exif = append(p.exif, raw)
				}
			case "Raw profile type xmp":
				p.xmp = append(p.xmp, rawProfile(value))
			case "Title":
				p.setText(&p.native.Title, value)
			case "Author":
				p.setText(&p.native.Creator, value)
			case "Description":
				p.setText(&p.native.Description, value)
			case "Comment":
				if comment == "" {
					comment = clean(value)
				}
			case "Copyright":
				p.setText(&p.native.Rights, value)
			case "Creation Time":
				p.setText(&p.native.Created, textDate(value))
			}
		case "eXIf":
			p.exif = append(p.exif, data)
		case "tIME":
			if len(data) == 7 {
				f := [6]int{int(be.Uint16(data)), int(data[2]), int(data[3]), int(data[4]), int(data[5]), int(data[6])}
				if validDate(f, 6) {
					p.setText(&p.native.Modified, w3cdtf(f, 6, "Z"))
				}
			}
		case "acTL":
			p.animated = len(data) >= 4 && be.Uint32(data) > 1
		case "IEND":
			off = len(b)
		}
	}
	if len(p.native.Description) == 0 && !placeholder(comment) {
		p.setText(&p.native.Description, comment)
	}
	return p, nil
}

// setText sets an element that has no value yet.
func (p *picture) setText(f *bdf.DCValues, v string) {
	if len(*f) == 0 {
		add(f, v)
	}
}

// pngText decodes a text chunk whose keyword is read (pngKeys): its
// keyword and text. Compressed text takes from budget.
func pngText(typ string, data []byte, budget *int) (key, value string, ok bool) {
	i := bytes.IndexByte(data, 0)
	if i < 1 || !pngKeys[string(data[:i])] {
		return "", "", false
	}
	key, rest := latin1(data[:i]), data[i+1:]
	switch typ {
	case "tEXt":
		return key, latin1(rest), true
	case "zTXt":
		if len(rest) < 1 || rest[0] != 0 {
			return "", "", false
		}
		t, err := inflate(rest[1:], budget)
		return key, latin1(t), err == nil
	}
	// iTXt: compression flag and method, language tag, translated keyword, UTF-8 text
	if len(rest) < 2 {
		return "", "", false
	}
	compressed := rest[0] == 1
	rest = rest[2:]
	for range 2 {
		j := bytes.IndexByte(rest, 0)
		if j < 0 {
			return "", "", false
		}
		rest = rest[j+1:]
	}
	if compressed {
		t, err := inflate(rest, budget)
		if err != nil {
			return "", "", false
		}
		rest = t
	}
	return key, string(bytes.ToValidUTF8(rest, nil)), true
}

// inflate decompresses zlib data, up to the bytes left in budget.
func inflate(b []byte, budget *int) ([]byte, error) {
	if *budget <= 0 {
		return nil, errors.New("too much compressed text")
	}
	r, err := zlib.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	t, err := io.ReadAll(io.LimitReader(r, int64(*budget)))
	*budget -= len(t)
	return t, err
}

// rawProfile decodes ImageMagick's "Raw profile type" text: a name line, a
// length line, then the bytes in hex.
func rawProfile(s string) []byte {
	lines := strings.SplitN(strings.TrimLeft(s, "\n"), "\n", 3)
	if len(lines) < 3 {
		return nil
	}
	b, _ := hex.DecodeString(strings.Join(strings.Fields(lines[2]), ""))
	return b
}

// readJPEG reads the frame header and the metadata segments of a JPEG
// image: EXIF and XMP (APP1) and Photoshop's IPTC (APP13).
func readJPEG(b []byte) (*picture, error) {
	p := &picture{format: "jpeg"}
	for off := 2; off+4 <= len(b); {
		if b[off] != 0xff {
			break
		}
		m := b[off+1]
		switch {
		case m == 0xff: // fill byte
			off++
			continue
		case m == 0x01 || m >= 0xd0 && m <= 0xd8: // markers without a segment
			off += 2
			continue
		case m == 0xd9 || m == 0xda: // metadata and the frame header come before the first scan
			off = len(b)
			continue
		}
		n := int(be.Uint16(b[off+2:]))
		if n < 2 || off+2+n > len(b) {
			break
		}
		seg := b[off+4 : off+2+n]
		off += 2 + n
		switch {
		case m == 0xe1 && bytes.HasPrefix(seg, []byte("Exif\x00\x00")):
			p.exif = append(p.exif, seg[6:])
		case m == 0xe1 && bytes.HasPrefix(seg, []byte("http://ns.adobe.com/xap/1.0/\x00")):
			p.xmp = append(p.xmp, seg[29:])
		case m == 0xed && bytes.HasPrefix(seg, []byte("Photoshop 3.0\x00")) && p.iptc == nil:
			p.iptc = photoshopIPTC(seg[14:])
		case m >= 0xc0 && m <= 0xcf && m != 0xc4 && m != 0xc8 && m != 0xcc && p.w == 0: // SOFn
			if len(seg) >= 5 {
				p.h, p.w = float64(be.Uint16(seg[1:])), float64(be.Uint16(seg[3:]))
			}
		}
	}
	if p.w == 0 || p.h == 0 {
		return nil, errors.New("the JPEG image has no frame header")
	}
	return p, nil
}

// readGIF reads the logical screen size of a GIF image, its first comment
// and XMP, and whether it has several frames.
func readGIF(b []byte) (*picture, error) {
	if len(b) < 13 {
		return nil, errTruncated
	}
	p := &picture{format: "gif", w: float64(le.Uint16(b[6:])), h: float64(le.Uint16(b[8:]))}
	off := 13
	if b[10]&0x80 != 0 { // global color table
		off += 3 << (b[10]&7 + 1)
	}
	frames := 0
	var comment string
loop:
	for off < len(b) {
		switch b[off] {
		case 0x21: // extension
			if off+2 > len(b) {
				break loop
			}
			label := b[off+1]
			off += 2
			if label == 0xff && off+12 <= len(b) && string(b[off:off+12]) == "\x0bXMP DataXMP" {
				// the packet is stored as it is, followed by a trailer that
				// makes it read as sub-blocks
				x := b[off+12:]
				if i := bytes.Index(x, []byte("<?xpacket end=")); i >= 0 {
					if j := bytes.Index(x[i:], []byte("?>")); j >= 0 {
						p.xmp = append(p.xmp, x[:i+j+2])
					}
				}
			}
			var data []byte
			data, off = subBlocks(b, off, label == 0xfe && comment == "")
			if data != nil {
				comment = clean(latin1(data))
			}
		case 0x2c: // image descriptor
			frames++
			if off+10 > len(b) {
				break loop
			}
			f := b[off+9]
			off += 10
			if f&0x80 != 0 { // local color table
				off += 3 << (f&7 + 1)
			}
			_, off = subBlocks(b, off+1, false) // after the LZW minimum code size
		default: // trailer
			break loop
		}
	}
	p.animated = frames > 1
	if !placeholder(comment) {
		add(&p.native.Description, comment)
	}
	return p, nil
}

// subBlocks reads the data sub-blocks at off, returning their contents
// (when keep is set) and the offset after the terminator.
func subBlocks(b []byte, off int, keep bool) ([]byte, int) {
	var data []byte
	for off < len(b) {
		n := int(b[off])
		off++
		if n == 0 {
			break
		}
		if keep {
			data = append(data, b[off:min(off+n, len(b))]...)
		}
		off += n
	}
	return data, off
}

// readWebP reads the canvas or frame size of a WebP image, its EXIF and
// XMP chunks, and whether it has several frames.
func readWebP(b []byte) (*picture, error) {
	if len(b) < 20 {
		return nil, errTruncated
	}
	p := &picture{format: "webp"}
	end := min(int64(len(b)), 8+int64(le.Uint32(b[4:])))
	frames := 0
	extended := false
	for off := int64(12); off+8 <= end; {
		typ := string(b[off : off+4])
		n := int64(le.Uint32(b[off+4:]))
		if off+8+n > end {
			break
		}
		data := b[off+8 : off+8+n]
		off += 8 + n + n&1
		switch typ {
		case "VP8X":
			if len(data) >= 10 {
				extended = true
				p.w = float64(1 + (uint32(data[4]) | uint32(data[5])<<8 | uint32(data[6])<<16))
				p.h = float64(1 + (uint32(data[7]) | uint32(data[8])<<8 | uint32(data[9])<<16))
			}
		case "VP8 ":
			if !extended && len(data) >= 10 && data[3] == 0x9d && data[4] == 0x01 && data[5] == 0x2a {
				p.w, p.h = float64(le.Uint16(data[6:])&0x3fff), float64(le.Uint16(data[8:])&0x3fff)
			}
		case "VP8L":
			if !extended && len(data) >= 5 && data[0] == 0x2f {
				bits := le.Uint32(data[1:])
				p.w, p.h = float64(1+bits&0x3fff), float64(1+bits>>14&0x3fff)
			}
		case "ANMF":
			frames++
		case "EXIF":
			p.exif = append(p.exif, data)
		case "XMP ":
			p.xmp = append(p.xmp, data)
		}
	}
	if p.w == 0 || p.h == 0 {
		return nil, errors.New("the WebP image has no size")
	}
	p.animated = frames > 1
	return p, nil
}

// bmpHeader reports whether b starts like a BMP file: "BM" and the size of
// a known DIB header.
func bmpHeader(b []byte) bool {
	if len(b) < 18 || b[0] != 'B' || b[1] != 'M' {
		return false
	}
	switch le.Uint32(b[14:]) {
	case 12, 16, 40, 52, 56, 64, 108, 124:
		return true
	}
	return false
}

// readBMP reads the size of a BMP image (a negative height is a top-down
// image).
func readBMP(b []byte) (*picture, error) {
	if !bmpHeader(b) || len(b) < 26 {
		return nil, errTruncated
	}
	p := &picture{format: "bmp"}
	if le.Uint32(b[14:]) == 12 { // BITMAPCOREHEADER
		p.w, p.h = float64(le.Uint16(b[18:])), float64(le.Uint16(b[20:]))
	} else {
		w, h := int32(le.Uint32(b[18:])), int32(le.Uint32(b[22:]))
		p.w, p.h = float64(w), float64(h)
		if h < 0 {
			p.h = -p.h
		}
	}
	if p.w <= 0 || p.h <= 0 {
		return nil, errors.New("the BMP image has no size")
	}
	return p, nil
}

// icoHeader reports whether b starts like an icon file: the header and a
// first directory entry whose image lies inside the file.
func icoHeader(b []byte, size int64) bool {
	if len(b) < 22 || le.Uint16(b) != 0 || le.Uint16(b[2:]) != 1 {
		return false
	}
	n := int64(le.Uint16(b[4:]))
	if n == 0 || b[9] != 0 {
		return false
	}
	dataSize, off := int64(le.Uint32(b[14:])), int64(le.Uint32(b[18:]))
	return dataSize > 0 && off >= 6+16*n && off+dataSize <= size
}

// readICO reads the size of the largest image of an icon file, the one
// browsers draw.
func readICO(b []byte) (*picture, error) {
	if !icoHeader(b, int64(len(b))) {
		return nil, errors.New("not an icon file")
	}
	p := &picture{format: "ico"}
	n := int(le.Uint16(b[4:]))
	bestBits := 0
	for i := range n {
		e := 6 + 16*i
		if e+16 > len(b) {
			break
		}
		w, h, bits := float64(b[e]), float64(b[e+1]), int(le.Uint16(b[e+6:]))
		if w == 0 {
			w = 256
		}
		if h == 0 {
			h = 256
		}
		if w*h > p.w*p.h || w*h == p.w*p.h && bits > bestBits {
			p.w, p.h, bestBits = w, h, bits
		}
	}
	if p.w == 0 {
		return nil, errTruncated
	}
	return p, nil
}

// readRaster reads an image in one of the raster formats.
func readRaster(format string, b []byte) (*picture, error) {
	switch format {
	case "png":
		return readPNG(b)
	case "jpeg":
		return readJPEG(b)
	case "gif":
		return readGIF(b)
	case "webp":
		return readWebP(b)
	case "avif":
		return readAVIF(b)
	case "bmp":
		return readBMP(b)
	case "ico":
		return readICO(b)
	}
	return nil, fmt.Errorf("unknown image format %q", format)
}
