package pdfbdf

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/klauspost/compress/zlib"
)

// ref is the number of an indirect object of the file (0: none).
type ref int

func (r ref) String() string { return strconv.Itoa(int(r)) + " 0 R" }

// writer writes the objects of a PDF file one after the other, and the
// table of where they are at the end.
type writer struct {
	w        *bufio.Writer
	off      int64
	offs     []int64 // by object number; 0 for one that is not written yet
	sum      hash.Hash
	err      error
	compress bool
	level    int
	zbuf     bytes.Buffer
}

func newWriter(w io.Writer, compress bool, level int) *writer {
	pw := &writer{w: bufio.NewWriterSize(w, 1<<16), offs: make([]int64, 1), sum: sha256.New(), compress: compress, level: level}
	// (the bytes of the second line tell tools that the file is binary)
	pw.write([]byte("%PDF-1.7\n%\xe2\xe3\xcf\xd3\n"))
	return pw
}

func (pw *writer) write(b []byte) {
	if pw.err != nil {
		return
	}
	n, err := pw.w.Write(b)
	pw.off += int64(n)
	pw.sum.Write(b)
	pw.err = err
}

// alloc reserves the number of an object that is written later.
func (pw *writer) alloc() ref {
	pw.offs = append(pw.offs, 0)
	return ref(len(pw.offs) - 1)
}

// put writes an object: a dictionary, an array or another value.
func (pw *writer) put(r ref, body string) {
	pw.offs[r] = pw.off
	pw.write([]byte(strconv.Itoa(int(r)) + " 0 obj\n" + body + "\nendobj\n"))
}

// add writes an object under a new number.
func (pw *writer) add(body string) ref {
	r := pw.alloc()
	pw.put(r, body)
	return r
}

// stream writes a stream object: dict is what its dictionary holds beside
// its length and its filter. The data is deflated when deflate is set and
// the file is compressed; other filters are those of the data as it is.
func (pw *writer) stream(r ref, dict string, data []byte, deflate bool) {
	if deflate && pw.compress {
		pw.zbuf.Reset()
		zw, _ := zlib.NewWriterLevel(&pw.zbuf, pw.level)
		zw.Write(data)
		zw.Close()
		data = pw.zbuf.Bytes()
		dict += " /Filter /FlateDecode"
	}
	pw.offs[r] = pw.off
	pw.write([]byte(strconv.Itoa(int(r)) + " 0 obj\n<<" + dict + " /Length " + strconv.Itoa(len(data)) + " >>\nstream\n"))
	pw.write(data)
	pw.write([]byte("\nendstream\nendobj\n"))
}

func (pw *writer) addStream(dict string, data []byte, deflate bool) ref {
	r := pw.alloc()
	pw.stream(r, dict, data, deflate)
	return r
}

// finish writes the table of the objects and the trailer.
func (pw *writer) finish(root, info ref) error {
	id := fmt.Sprintf("%x", pw.sum.Sum(nil)[:16])
	start := pw.off
	var b strings.Builder
	fmt.Fprintf(&b, "xref\n0 %d\n0000000000 65535 f \n", len(pw.offs))
	for _, off := range pw.offs[1:] {
		if off == 0 {
			// a number that no object took
			b.WriteString("0000000000 65535 f \n")
			continue
		}
		fmt.Fprintf(&b, "%010d 00000 n \n", off)
	}
	fmt.Fprintf(&b, "trailer\n<< /Size %d /Root %s", len(pw.offs), root)
	if info != 0 {
		fmt.Fprintf(&b, " /Info %s", info)
	}
	fmt.Fprintf(&b, " /ID [<%s> <%s>] >>\nstartxref\n%d\n%%%%EOF\n", id, id, start)
	pw.write([]byte(b.String()))
	if pw.err != nil {
		return pw.err
	}
	return pw.w.Flush()
}

// text writes a text string of a dictionary: as it is in parentheses when
// it is ASCII, else as UTF-16 in hexadecimal.
func text(s string) string {
	ascii := true
	for i := 0; i < len(s); i++ {
		if s[i] < 0x20 || s[i] > 0x7e {
			ascii = false
			break
		}
	}
	if ascii {
		var b strings.Builder
		b.WriteByte('(')
		for i := 0; i < len(s); i++ {
			switch c := s[i]; c {
			case '(', ')', '\\':
				b.WriteByte('\\')
				b.WriteByte(c)
			default:
				b.WriteByte(c)
			}
		}
		b.WriteByte(')')
		return b.String()
	}
	var b strings.Builder
	b.WriteString("<FEFF")
	for _, u := range utf16.Encode([]rune(s)) {
		fmt.Fprintf(&b, "%04X", u)
	}
	b.WriteByte('>')
	return b.String()
}

// name writes a name object.
func name(s string) string {
	var b strings.Builder
	b.WriteByte('/')
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c <= 0x20 || c >= 0x7f || strings.IndexByte("#()<>[]{}/%", c) >= 0 {
			fmt.Fprintf(&b, "#%02X", c)
			continue
		}
		b.WriteByte(c)
	}
	return b.String()
}
