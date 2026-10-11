package pdf

import (
	"bytes"
	"io"
	"regexp"
	"sort"
	"strings"
)

// An AES stream is an initialization vector of 16 bytes and the data, padded
// to whole blocks: an empty one takes 32 bytes. Some writers leave the block
// of padding out and write an empty stream as its vector alone (seen in a
// file whose producer is PDF-XChange Standard 10.5); pdfcpu refuses such a
// file ("ciphertext too short"). A file refused for that is read
// again through emptyStreams, which shows pdfcpu those streams empty and
// in the clear: a Length of 0, which makes pdfcpu look for endstream, and
// line ends for data. Nothing moves in the file, so the cross-reference
// table holds.

// ivOnly is the size of an AES stream that is its initialization vector alone.
const ivOnly = 16

// shortCiphertext reports that pdfcpu refused an AES stream or string
// shorter than a vector and a block. pdfcpu does not export the error.
func shortCiphertext(err error) bool {
	return err != nil && strings.Contains(err.Error(), "ciphertext too short")
}

// patch replaces the bytes of a file from an offset on.
type patch struct {
	at   int64
	with []byte
}

// patched reads a file with some of its bytes replaced.
type patched struct {
	rs      io.ReadSeeker
	at      int64
	patches []patch // in the order of the file
}

func (p *patched) Read(b []byte) (int, error) {
	n, err := p.rs.Read(b)
	end := p.at + int64(n)
	// the first patch that ends after the bytes read start
	i := sort.Search(len(p.patches), func(i int) bool { return p.patches[i].at+int64(len(p.patches[i].with)) > p.at })
	for ; i < len(p.patches) && p.patches[i].at < end; i++ {
		pt := p.patches[i]
		if pt.at >= p.at {
			copy(b[pt.at-p.at:n], pt.with)
		} else {
			copy(b[:n], pt.with[p.at-pt.at:])
		}
	}
	p.at = end
	return n, err
}

func (p *patched) Seek(offset int64, whence int) (int64, error) {
	at, err := p.rs.Seek(offset, whence)
	if err == nil {
		p.at = at
	}
	return at, err
}

const (
	// scanChunk is how much of the file the scan reads at a time, and
	// scanBack how far before its endstream an object may start.
	scanChunk = 1 << 20
	scanBack  = 64 << 10
)

const pdfSpace = `[\x00\t\n\f\r ]`

var (
	// the header of the object a stream dictionary is in: the last one before the stream
	objHeader = regexp.MustCompile(`[0-9]+` + pdfSpace + `+[0-9]+` + pdfSpace + `+obj`)
	// the Length of a stream, a number or a reference to one (not Length1 of a font file)
	streamLength = regexp.MustCompile(`/Length` + pdfSpace + `+(([0-9]+)(?:` + pdfSpace + `+[0-9]+` + pdfSpace + `+R)?)`)
	// streams that are never encrypted, or not with the document's filter
	clearStream = regexp.MustCompile(`/Type` + pdfSpace + `*/XRef|/Crypt`)
)

// emptyStreams returns rs with every stream that is 16 bytes long shown as
// an empty one, or nil when the file has none. The streams are found by
// their shape, without the cross-reference table: a stream dictionary of an
// object, then stream, 16 bytes and endstream.
func emptyStreams(rs io.ReadSeeker) (io.ReadSeeker, error) {
	if _, err := rs.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	var patches []patch
	buf := make([]byte, 0, scanChunk+scanBack)
	base := int64(0) // the offset of buf[0] in the file
	from := 0        // where the search goes on in buf
	for {
		n, err := io.ReadFull(rs, buf[len(buf):min(cap(buf), len(buf)+scanChunk)])
		buf = buf[:len(buf)+n]
		for {
			i := bytes.Index(buf[from:], []byte("endstream"))
			if i < 0 {
				break
			}
			i += from
			from = i + len("endstream")
			for _, p := range emptyStream(buf, i) {
				patches = append(patches, patch{at: base + int64(p.at), with: p.with})
			}
		}
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			return nil, err
		}
		// an endstream cut by the end of the chunk is found with the next one
		from = max(from, len(buf)-len("endstream")+1)
		if drop := len(buf) - scanBack; drop > 0 {
			copy(buf, buf[drop:])
			buf = buf[:scanBack]
			base += int64(drop)
			from -= drop
		}
	}
	if len(patches) == 0 {
		return nil, nil
	}
	return &patched{rs: rs, patches: patches}, nil
}

// emptyStream returns the patches (with offsets in buf) that empty the
// stream ending with the endstream at end, when it is 16 bytes long and
// can be encrypted.
func emptyStream(buf []byte, end int) []patch {
	// the data, and an optional line end before endstream
	for _, tail := range []string{"", "\n", "\r", "\r\n"} {
		data := end - len(tail) - ivOnly
		if data < 0 || string(buf[data+ivOnly:end]) != tail {
			continue
		}
		// the line end after stream is LF or CR LF
		kw := data - 1
		if kw < 1 || buf[kw] != '\n' {
			continue
		}
		if buf[kw-1] == '\r' {
			kw--
		}
		kw -= len("stream")
		if kw < 0 || !bytes.HasPrefix(buf[kw:], []byte("stream")) || bytes.HasSuffix(buf[:kw], []byte("end")) {
			continue
		}
		// its dictionary: all there is between the header of the object and the stream
		start := max(0, kw-scanBack)
		headers := objHeader.FindAllIndex(buf[start:kw], -1)
		if headers == nil {
			return nil
		}
		dictAt := start + headers[len(headers)-1][1]
		dict := buf[dictAt:kw]
		if t := bytes.TrimSpace(dict); !bytes.HasPrefix(t, []byte("<<")) || !bytes.HasSuffix(t, []byte(">>")) || clearStream.Match(dict) {
			return nil
		}
		m := streamLength.FindSubmatchIndex(dict)
		if m == nil {
			return nil
		}
		value := dict[m[2]:m[3]]
		if m[3] == m[5] && string(value) != "16" {
			return nil // another length said, directly: not a stream of this shape
		}
		zero := bytes.Repeat([]byte{' '}, len(value))
		zero[0] = '0'
		return []patch{{at: int64(dictAt + m[2]), with: zero}, {at: int64(data), with: bytes.Repeat([]byte{'\n'}, ivOnly)}}
	}
	return nil
}
