package cff

import (
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
