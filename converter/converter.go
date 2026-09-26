// Package converter holds what the format converters have in common:
// input format detection and page selection. The converters themselves are
// its subpackages: converter/pdf (PDF), converter/pptx (PowerPoint) and
// converter/drawio (draw.io diagrams).
package converter

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// Format is an input document format.
type Format string

const (
	Unknown Format = ""
	PDF     Format = "pdf"
	PPTX    Format = "pptx"
	Drawio  Format = "drawio"
)

// Detect sniffs the format of a document: PDF by its header, PowerPoint by
// the presentation part of an Office Open XML package (.pptx, .pptm, .ppsx,
// .potx), draw.io by its mxfile or mxGraphModel XML (also embedded in SVG
// and PNG exports). The legacy binary .ppt format is reported as Unknown.
func Detect(r io.ReaderAt, size int64) Format {
	head := make([]byte, 1024)
	n, _ := r.ReadAt(head, 0)
	head = head[:n]
	if bytes.Contains(head, []byte("%PDF-")) {
		return PDF
	}
	if isDrawio(r, size, head) {
		return Drawio
	}
	if !bytes.HasPrefix(head, []byte("PK\x03\x04")) {
		return Unknown
	}
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return Unknown
	}
	for _, f := range zr.File {
		if strings.EqualFold(f.Name, "ppt/presentation.xml") {
			return PPTX
		}
	}
	return Unknown
}

// maxDrawioSniff is the most read to find a diagram in an SVG or PNG export.
const maxDrawioSniff = 64 << 20

// isDrawio reports whether the input is a draw.io file: XML with an mxfile
// or mxGraphModel root, an SVG export whose root element carries the
// diagram in its content attribute, or a PNG export with an mxfile text
// chunk.
func isDrawio(r io.ReaderAt, size int64, head []byte) bool {
	t := bytes.TrimLeft(head, "\xef\xbb\xbf \t\r\n")
	switch {
	case bytes.HasPrefix(t, []byte("<")):
		if bytes.Contains(t, []byte("<mxfile")) || bytes.Contains(t, []byte("<mxGraphModel")) {
			return true
		}
		i := bytes.Index(t, []byte("<svg"))
		if i < 0 {
			return false
		}
		// the root tag may be longer than the head: read up to its end
		tag := t[i:]
		if !bytes.Contains(tag, []byte(">")) && size <= maxDrawioSniff {
			all := make([]byte, size)
			n, _ := r.ReadAt(all, 0)
			all = all[:n]
			if j := bytes.Index(all, []byte("<svg")); j >= 0 {
				tag = all[j:]
			}
		}
		if end := bytes.IndexByte(tag, '>'); end >= 0 {
			tag = tag[:end]
		}
		return bytes.Contains(tag, []byte(" content=\"")) &&
			(bytes.Contains(tag, []byte("&lt;mxfile")) || bytes.Contains(tag, []byte("&lt;mxGraphModel")))
	case bytes.HasPrefix(head, []byte("\x89PNG\r\n\x1a\n")):
		// text chunks come before the image data; look through the chunk list
		off := int64(8)
		var hdr [8]byte
		for off+8 <= size {
			if _, err := r.ReadAt(hdr[:], off); err != nil {
				return false
			}
			n := int64(uint32(hdr[0])<<24 | uint32(hdr[1])<<16 | uint32(hdr[2])<<8 | uint32(hdr[3]))
			typ := string(hdr[4:8])
			if typ == "tEXt" || typ == "zTXt" {
				key := make([]byte, 13)
				m, _ := r.ReadAt(key, off+8)
				key = key[:m]
				if bytes.HasPrefix(key, []byte("mxfile\x00")) || bytes.HasPrefix(key, []byte("mxGraphModel\x00")) {
					return true
				}
			}
			if typ == "IEND" || typ == "IDAT" {
				return false
			}
			off += 12 + n
		}
	}
	return false
}

// DetectFile is Detect for a file path.
func DetectFile(path string) (Format, error) {
	f, err := os.Open(path)
	if err != nil {
		return Unknown, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return Unknown, err
	}
	return Detect(f, st.Size()), nil
}

// PageRange parses "1-3,5,8-" style selections of 1-based page (or slide)
// numbers, clamped to 1..count.
func PageRange(spec string, count int) ([]int, error) {
	var out []int
	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		lo, hi := 1, count
		if i := strings.IndexByte(part, '-'); i >= 0 {
			if i > 0 {
				if _, err := fmt.Sscanf(part[:i], "%d", &lo); err != nil {
					return nil, fmt.Errorf("bad page range %q", part)
				}
			}
			if i+1 < len(part) {
				if _, err := fmt.Sscanf(part[i+1:], "%d", &hi); err != nil {
					return nil, fmt.Errorf("bad page range %q", part)
				}
			}
		} else {
			if _, err := fmt.Sscanf(part, "%d", &lo); err != nil {
				return nil, fmt.Errorf("bad page %q", part)
			}
			hi = lo
		}
		lo, hi = max(1, lo), min(count, hi)
		for n := lo; n <= hi; n++ {
			out = append(out, n)
		}
	}
	return out, nil
}
