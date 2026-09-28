package metafile

import (
	"bytes"
	"encoding/binary"
	"slices"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// Metafiles that state more than they hold, or make the replay keep and
// write more than a picture needs: the replay takes what is there, leaves
// out what goes beyond its limits and says so.

// TestPointCount reads a polyline that states a million points and holds
// two: the points come from the record, and so does the room for them.
func TestPointCount(t *testing.T) {
	r := make([]byte, 28+8)
	binary.LittleEndian.PutUint32(r[24:], 1<<20)
	copy(r[28:], []byte{10, 0, 20, 0, 30, 0, 40, 0})
	pts := emfPoints(r, 28, 1<<20, true)
	if want := [][2]float64{{10, 20}, {30, 40}}; !slices.Equal(pts, want) || cap(pts) != 2 {
		t.Errorf("points %v with room for %d, want %v", pts, cap(pts), want)
	}
	if pts := emfPoints(r, 28, 1<<20, false); len(pts) != 1 || pts[0] != [2]float64{10 | 20<<16, 30 | 40<<16} {
		t.Errorf("32-bit points %v", pts)
	}
	for _, off := range []int{-4, 36, 40} {
		if pts := emfPoints(r, off, 1<<20, true); len(pts) != 0 {
			t.Errorf("points at %d: %v", off, pts)
		}
	}

	b := emfHeader()
	emfRecord(b, 87, []int32{0, 0, 99, 99, 1 << 20}, r[28:]) // POLYLINE16
	emfRecord(b, 14, []int32{0, 16, 20}, nil)
	if d := replay(t, b.Bytes(), nil); d.ops[bdf.OpStrokePath] != 1 || len(d.warnings) != 0 {
		t.Errorf("%d lines drawn, warnings %v", d.ops[bdf.OpStrokePath], d.warnings)
	}
}

// TestDIBDepth refuses a bitmap of a depth that has no rows before it makes
// room for its pixels.
func TestDIBDepth(t *testing.T) {
	for _, bpp := range []uint16{0, 3, 12, 64} {
		var bmi bytes.Buffer
		binary.Write(&bmi, binary.LittleEndian, []int32{40, 64, 64})
		binary.Write(&bmi, binary.LittleEndian, []uint16{1, bpp})
		binary.Write(&bmi, binary.LittleEndian, []int32{0, 0, 0, 0, 0, 0})
		bits := make([]byte, 64*64*8)
		var err error
		allocs := testing.AllocsPerRun(5, func() { _, err = decodeDIB(bmi.Bytes(), bits) })
		if err == nil || allocs != 0 {
			t.Errorf("%d bits per pixel: error %v after %v allocations", bpp, err, allocs)
		}
	}
}

// limit sets a limit of the replay for a test.
func limit(t *testing.T, v *int, n int) {
	t.Helper()
	old := *v
	*v = n
	t.Cleanup(func() { *v = old })
}

func TestClipLimits(t *testing.T) {
	// Every clip opens a drawing group that writes all the clips in effect:
	// a hundred clips are 5050 polygons without a limit.
	b := emfHeader()
	for range 100 {
		emfRecord(b, 30, []int32{0, 0, 90, 90}, nil)   // INTERSECTCLIPRECT
		emfRecord(b, 43, []int32{10, 10, 50, 50}, nil) // RECTANGLE
	}
	emfRecord(b, 14, []int32{0, 16, 20}, nil)
	d := replay(t, b.Bytes(), nil)
	if want := maxClips * (maxClips + 1) / 2; d.ops[bdf.OpClipPath] != want || d.ops[bdf.OpFillPath] != 100 {
		t.Errorf("%d clips written for %d rectangles, want %d", d.ops[bdf.OpClipPath], d.ops[bdf.OpFillPath], want)
	}
	if want := []string{"a metafile has too many clips in effect; the further ones are not applied"}; !slices.Equal(d.warnings, want) {
		t.Errorf("warnings = %q", d.warnings)
	}

	// One clip, written by every group: the groups after the limit of
	// points are drawn without it.
	limit(t, &maxClipPoints, 100)
	b = emfHeader()
	emfRecord(b, 30, []int32{0, 0, 90, 90}, nil)
	for i := range 40 {
		emfRecord(b, 10, []int32{int32(i), 0}, nil) // SETWINDOWORGEX: a new group
		emfRecord(b, 43, []int32{10, 10, 50, 50}, nil)
	}
	emfRecord(b, 14, []int32{0, 16, 20}, nil)
	d = replay(t, b.Bytes(), nil)
	if d.ops[bdf.OpClipPath] != 25 || d.ops[bdf.OpFillPath] != 40 {
		t.Errorf("%d clips written for %d rectangles, want 25", d.ops[bdf.OpClipPath], d.ops[bdf.OpFillPath])
	}
	if want := []string{"the clips of a metafile are too complex; the rest of it is drawn without them"}; !slices.Equal(d.warnings, want) {
		t.Errorf("warnings = %q", d.warnings)
	}
}

func TestSavedStates(t *testing.T) {
	// A saved state shares the clips in effect instead of copying them.
	g := newGDI(canvas.NewBuilder(bdf.NewDocument(), nil).New(), &Options{})
	g.out = canvas.Identity
	g.clipRect(0, 0, 10, 10)
	g.clipRect(0, 0, 5, 5)
	g.save()
	if len(g.stack) != 1 || len(g.stack[0].clips) != 2 || &g.stack[0].clips[0] != &g.st.clips[0] {
		t.Error("the saved state has clips of its own")
	}
	// A clip after it leaves the saved ones alone.
	g.clipRect(0, 0, 2, 2)
	if g.restore(-1); len(g.st.clips) != 2 || g.st.clips[1][2] != [2]float64{5, 5} {
		t.Errorf("restored clips = %v", g.st.clips)
	}

	limit(t, &maxSaved, 4)
	b := emfHeader()
	for range 6 {
		emfRecord(b, 33, nil, nil) // SAVEDC
	}
	emfRecord(b, 43, []int32{10, 10, 50, 50}, nil)
	emfRecord(b, 14, []int32{0, 16, 20}, nil)
	d := replay(t, b.Bytes(), nil)
	if want := []string{"a metafile saves too many states; the further ones are not saved"}; !slices.Equal(d.warnings, want) || !d.drawn {
		t.Errorf("drawn %v, warnings = %q", d.drawn, d.warnings)
	}
}

// wmfOf makes a Windows metafile of records: function numbers and their
// parameters.
func wmfOf(recs ...[]uint16) []byte {
	var b bytes.Buffer
	binary.Write(&b, binary.LittleEndian, []uint16{1, 9, 0x300, 0, 0, 0, 0, 0, 0})
	for _, r := range append(recs, []uint16{0}) {
		binary.Write(&b, binary.LittleEndian, uint32(2+len(r)))
		binary.Write(&b, binary.LittleEndian, r)
	}
	return b.Bytes()
}

func TestObjectTable(t *testing.T) {
	// A new object takes the lowest free entry of the table.
	brush := func(rgb uint16) []uint16 { return []uint16{0x02FC, 0, rgb, 0, 0} } // CREATEBRUSHINDIRECT: solid, a shade of red
	d := replay(t, wmfOf(
		brush(1), brush(2), brush(3), // entries 0, 1, 2
		[]uint16{0x01F0, 1}, []uint16{0x01F0, 0}, []uint16{0x01F0, 0}, // DELETEOBJECT
		brush(4), brush(5), brush(6), // entries 0, 1, 3
		[]uint16{0x012D, 0}, []uint16{0x041B, 50, 50, 10, 10}, // SELECTOBJECT, RECTANGLE
		[]uint16{0x012D, 1}, []uint16{0x041B, 50, 50, 10, 10},
		[]uint16{0x012D, 2}, []uint16{0x041B, 50, 50, 10, 10},
		[]uint16{0x012D, 3}, []uint16{0x041B, 50, 50, 10, 10},
	), nil)
	want := []bdf.Color{bdf.RGBA(4, 0, 0, 255), bdf.RGBA(5, 0, 0, 255), bdf.RGBA(3, 0, 0, 255), bdf.RGBA(6, 0, 0, 255)}
	if !slices.Equal(d.colors, want) {
		t.Errorf("brushes = %08X, want %08X", d.colors, want)
	}

	// A table that only grows: finding the free entry must not walk over it
	// for every object (that took seconds for this many).
	recs := make([][]uint16, 100000)
	for i := range recs {
		recs[i] = []uint16{0x06FF} // CREATEREGION
	}
	data := wmfOf(recs...)
	start := time.Now()
	replay(t, data, nil)
	if d := time.Since(start); d > time.Second {
		t.Errorf("%d objects took %v", len(recs), d)
	}
}
