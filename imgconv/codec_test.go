package imgconv

import "testing"

// TestCodecBinding records which codec this build compiled in, so a
// GOEXPERIMENT=simd run shows in the log that it exercised webpwsimd.
func TestCodecBinding(t *testing.T) {
	t.Logf("Available=%v SIMD=%v", Available(), SIMD())
	if SIMD() && !Available() {
		t.Fatal("SIMD codec without Available")
	}
}
