package pdf

import (
	"bytes"
	"github.com/shibukawa/bdf/converter/internal/sfnt"
	"io"
)

func bytesReader(b []byte) io.Reader { return bytes.NewReader(b) }

func readAll(r io.Reader) ([]byte, error) { return io.ReadAll(r) }

func be16(b []byte, i int) uint16 { return sfnt.BE16(b, i) }

func be32(b []byte, i int) uint32 { return sfnt.BE32(b, i) }
