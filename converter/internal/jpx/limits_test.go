package jpx

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
)

// emptyTiles makes a codestream of w × h tiles of one pixel, each with a
// tile-part that holds no packets, for comps components with levels
// decomposition levels.
func emptyTiles(w, h, comps, levels int) []byte {
	var b bytes.Buffer
	put := func(v ...any) {
		for _, x := range v {
			binary.Write(&b, binary.BigEndian, x)
		}
	}
	put(uint16(mSOC), uint16(mSIZ), uint16(38+3*comps), uint16(0))
	put(uint32(w), uint32(h), uint32(0), uint32(0)) // image size and offset
	put(uint32(1), uint32(1), uint32(0), uint32(0)) // tile size and offset
	put(uint16(comps))
	for range comps {
		put([]byte{7, 1, 1}) // 8 bits, not subsampled
	}
	// LRCP, one layer, no component transform; code-blocks of 64 × 64, 5/3 filter
	put(uint16(mCOD), uint16(12), []byte{0, progLRCP, 0, 1, 0, byte(levels), 4, 4, 0, 1})
	put(uint16(mQCD), uint16(4), []byte{0x40, 0x40}) // no quantization, two guard bits
	for i := range w * h {
		put(uint16(mSOT), uint16(10), uint16(i), uint32(14), []byte{0, 1}, uint16(mSOD))
	}
	put(uint16(mEOC))
	return b.Bytes()
}

// TestManyTiles decodes codestreams whose tiles hold no data: laying out
// the resolution levels of every component of every tile is work that the
// size of the image does not bound, so it has a limit of its own.
func TestManyTiles(t *testing.T) {
	img, err := Decode(emptyTiles(4, 4, 3, 5))
	if err != nil || img.Width != 4 || len(img.Components) != 3 || img.Components[0].Data[15] != 128 {
		t.Fatalf("16 tiles of 3 components: %v", err)
	}

	old := maxLevels
	defer func() { maxLevels = old }()
	maxLevels = 16 * 8 * 33
	if _, err := Decode(emptyTiles(4, 4, 8, 32)); err != nil {
		t.Errorf("%d resolution levels: %v", maxLevels, err)
	}
	if _, err := Decode(emptyTiles(4, 4, 8, 32)); err != nil {
		t.Errorf("the same once more: %v", err) // the count is the codestream's
	}
	maxLevels--
	if _, err := Decode(emptyTiles(4, 4, 8, 32)); err == nil || !strings.Contains(err.Error(), "too many tiles and components") {
		t.Errorf("one more than %d resolution levels: %v", maxLevels, err)
	}
}
