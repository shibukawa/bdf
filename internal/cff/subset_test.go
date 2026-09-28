package cff

import (
	"errors"
	"slices"
	"testing"
)

func TestT2Exec(t *testing.T) {
	num := func(v int) byte { return byte(v + 139) } // -107..107
	local := newSubrSet([][]byte{{11}, {num(5), 11}, {11}})
	global := newSubrSet([][]byte{{11}})
	// Two stems, then a hintmask whose mask byte 0xff must not be read as an
	// operand; subr 1 (biased -106) is called, then seac 'A' + grave (code 193).
	cs := []byte{num(10), num(20), num(30), num(40), 1, 19, 0xff, num(-106), 10, num(0), num(0), num(65), 247, 193 - 108, 14}
	st := &t2State{seac: [2]int{-1, -1}}
	done, err := st.exec(cs, local, global, 0)
	if err != nil || !done {
		t.Fatalf("exec = %v, %v", done, err)
	}
	if !slices.Equal(local.used, []bool{false, true, false}) || global.used[0] {
		t.Fatalf("used = %v %v", local.used, global.used)
	}
	if st.nStems != 2 || st.seac != [2]int{65, 193} {
		t.Fatalf("stems %d, seac %v", st.nStems, st.seac)
	}
	// Arithmetic can compute subroutine numbers: give up rather than guess.
	st = &t2State{seac: [2]int{-1, -1}}
	if _, err := st.exec([]byte{num(1), num(2), 12, 10, 10, 14}, local, global, 0); err == nil {
		t.Fatal("arithmetic operator was followed")
	}
	// A subroutine number out of range is an error too.
	if _, err := (&t2State{}).exec([]byte{num(50), 10, 14}, local, global, 0); err == nil {
		t.Fatal("out-of-range subroutine was followed")
	}
}

// program assembles a name-keyed font of global subroutines and
// charstrings; top makes the Top DICT of the offset of the CharStrings.
func program(gsubrs, charStrings [][]byte, top func(charStrings int) []byte) []byte {
	header, name, strs := []byte{1, 0, 4, 1}, Index([][]byte{[]byte("A")}), Index(nil)
	at := len(header) + len(name) + len(Index([][]byte{top(0)})) + len(strs) + len(Index(gsubrs))
	var out []byte
	for _, b := range [][]byte{header, name, Index([][]byte{top(at)}), strs, Index(gsubrs), Index(charStrings)} {
		out = append(out, b...)
	}
	return out
}

// number writes an integer operand of a DICT in five bytes.
func number(v int) []byte { return []byte{29, byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)} }

// TestSubsetBounded subsets fonts whose charstrings run long: subroutines
// that call the next one eight times, ten levels deep, and a subroutine
// that pushes operands called again and again.
func TestSubsetBounded(t *testing.T) {
	call := func(subr int) []byte { return []byte{byte(subr - 107 + 139), 29} } // callgsubr
	top := func(at int) []byte { return append(number(at), 17) }
	var gsubrs [][]byte
	for level := range 10 {
		var s []byte
		for range 8 {
			if level < 9 {
				s = append(s, call(level+1)...)
			}
		}
		gsubrs = append(gsubrs, append(s, 11))
	}
	glyph := append(call(0), 14)
	if _, _, err := Subset(program(gsubrs, [][]byte{glyph}, top), map[int]bool{0: true}); !errors.Is(err, ErrTooLong) {
		t.Errorf("subroutines that call subroutines: error %v, want ErrTooLong", err)
	}

	push := make([]byte, 0, 1001)
	for range 1000 {
		push = append(push, 139)
	}
	var calls []byte
	for range 20 {
		calls = append(calls, call(1)...)
	}
	gsubrs = [][]byte{append(calls, 11), append(push, 11)}
	if _, _, err := Subset(program(gsubrs, [][]byte{glyph}, top), map[int]bool{0: true}); !errors.Is(err, ErrTooLong) {
		t.Errorf("operands left on the stack: error %v, want ErrTooLong", err)
	}

	// subroutines called a few times are followed
	gsubrs = [][]byte{append(append(call(1), call(1)...), 11), {11}}
	out, order, err := Subset(program(gsubrs, [][]byte{glyph, glyph}, top), map[int]bool{1: true})
	if err != nil || len(order) != 2 || len(out) == 0 {
		t.Errorf("subroutines called twice: %d glyphs, error %v", len(order), err)
	}
}

// TestOperands reads DICTs whose offsets are real numbers no file is as
// large as: the font is not read, and nothing is read out of range.
func TestOperands(t *testing.T) {
	huge := []byte{30, 0x1b, 0x30, 0x0f} // 1E300
	// the largest number below 2^63: added to an offset, it would wrap
	wraps := []byte{30, 0x9a, 0x22, 0x33, 0x72, 0x03, 0x68, 0x54, 0x77, 0x47, 0x84, 0xb1, 0x8f}
	glyph := make([]byte, 2000)
	for i := range glyph {
		glyph[i] = 14
	}
	for name, top := range map[string]func(at int) []byte{
		"CharStrings": func(at int) []byte { return append(huge[:len(huge):len(huge)], 17) },
		"Encoding": func(at int) []byte {
			return append(append(append(number(at), 17), huge...), 16)
		},
		"charset": func(at int) []byte {
			return append(append(append(number(at), 17), huge...), 15)
		},
		"Private": func(at int) []byte {
			return append(append(append(append(number(at), 17), wraps...), number(1500)...), 18)
		},
	} {
		data := program(nil, [][]byte{glyph}, top)
		if f, err := Parse(data); err == nil {
			if f.NumGlyphs != 1 {
				t.Errorf("%s: %d glyphs", name, f.NumGlyphs)
			}
			if _, _, err := Subset(data, map[int]bool{0: true}); err == nil && name == "Private" {
				t.Errorf("%s: the font is subset", name)
			}
		} else if name != "CharStrings" {
			t.Errorf("%s: %v", name, err)
		}
	}
	for v, want := range map[float64]int{0: 0, 1500: 1500, -3: -3, 1e300: -1 << 40, -1e300: -1 << 40, 1e15: 1 << 40, 9223372036854774784: 1 << 40} {
		if got := operand(v); got != want {
			t.Errorf("operand(%g) = %d, want %d", v, got, want)
		}
	}
}
