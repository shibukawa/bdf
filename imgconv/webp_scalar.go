//go:build !bdf_noconv && !(goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64))

package imgconv

import (
	"github.com/shibukawa/bdf/imgconv/internal/webpw"
	"github.com/shibukawa/bdf/imgconv/internal/webpw/base"
)

// The scalar codec: libwebp built without SIMD, translated over [2]uint64
// pairs. Every build that is not the archsimd one below uses it.

// SIMD reports whether the codec compiled in runs libwebp's SIMD kernels
// over simd/archsimd vector registers (a GOEXPERIMENT=simd build on Go
// 1.27, amd64 or arm64).
func SIMD() bool { return false }

type codecModule = base.Module

func codecNew() *codecModule {
	m := webpw.New()
	webpw.Initialize(m)
	return m
}

func codecMalloc(m *codecModule, n int32) int32 { return webpw.Malloc(m, n) }
func codecMemory(m *codecModule) []byte         { return webpw.Memory(m) }
func codecFree(m *codecModule, p int32)         { webpw.Free(m, p) }

func codecEncode(m *codecModule, in, w, h, sizePtr, quality, method, lossless, exact int32) int32 {
	return webpw.Encode(m, in, w, h, sizePtr, quality, method, lossless, exact)
}
