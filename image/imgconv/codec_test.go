package imgconv

import "testing"

func TestCodecAvailability(t *testing.T) {
	if SIMD() && !Available() {
		t.Fatal("SIMD codec without Available")
	}
}
