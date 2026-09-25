// Package converter holds what the format converters have in common:
// input format detection and page selection. The converters themselves are
// its subpackages: converter/pdf (PDF) and converter/pptx (PowerPoint).
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
)

// Detect sniffs the format of a document: PDF by its header, PowerPoint by
// the presentation part of an Office Open XML package (.pptx, .pptm, .ppsx,
// .potx). The legacy binary .ppt format is reported as Unknown.
func Detect(r io.ReaderAt, size int64) Format {
	head := make([]byte, 1024)
	n, _ := r.ReadAt(head, 0)
	head = head[:n]
	if bytes.Contains(head, []byte("%PDF-")) {
		return PDF
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
