package pdf2bdf

import (
	"bytes"
	"io"
)

func bytesReader(b []byte) io.Reader { return bytes.NewReader(b) }

func readAll(r io.Reader) ([]byte, error) { return io.ReadAll(r) }
