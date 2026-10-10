//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64) && arm64 && !darwin && !linux

package base

// CPUDotProd / CPUI8MM / CPUFHM: no detection story on this OS — run
// the portable bodies.
var CPUDotProd = false
var CPUI8MM = false
var CPUFHM = false
