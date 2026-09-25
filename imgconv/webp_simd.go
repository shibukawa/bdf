//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package imgconv

import (
	"github.com/shibukawa/bdf/imgconv/internal/webpwsimd"
	"github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
)

// The SIMD codec: libwebp's SSE2/SSE4.1 kernels through wasm SIMD,
// translated over simd/archsimd vector registers (wasm2go -simd=go127).
// Selected by a GOEXPERIMENT=simd build on Go 1.27; needs AVX2 on amd64
// at run time (checked at init).

// SIMD reports whether the codec compiled in runs libwebp's SIMD kernels
// over simd/archsimd vector registers (a GOEXPERIMENT=simd build on Go
// 1.27, amd64 or arm64).
func SIMD() bool { return true }

type codecModule = base.Module

func codecNew() *codecModule {
	m := webpwsimd.New()
	webpwsimd.Initialize(m)
	return m
}

func codecMalloc(m *codecModule, n int32) int32 { return webpwsimd.Malloc(m, n) }
func codecMemory(m *codecModule) []byte         { return webpwsimd.Memory(m) }
func codecFree(m *codecModule, p int32)         { webpwsimd.Free(m, p) }

func codecEncode(m *codecModule, in, w, h, sizePtr, quality, method, lossless, exact int32) int32 {
	return webpwsimd.Encode(m, in, w, h, sizePtr, quality, method, lossless, exact)
}
