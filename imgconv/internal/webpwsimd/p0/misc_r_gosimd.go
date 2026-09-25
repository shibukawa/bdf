//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_RD4_SSE2(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 base.V128
	_ = v2
	var v8 int32
	_ = v8
	var v9 base.V128
	_ = v9
	var v13 int32
	_ = v13
	var v18 int32
	_ = v18
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	var v30 base.V128
	_ = v30
	var v31 base.V128
	_ = v31
	var v32 base.V128
	_ = v32
	var v38 base.V128
	_ = v38
	var v40 base.V128
	_ = v40
	v2 = base.Simd_g_const(&F_RD4_SSE2__k0)
	v8 = int32(0)
	v9 = base.Simd_g_v128_load64_zero(m, l0+int32(-33), v8)
	v13 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+31)))
	v18 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-1)))))
	v22 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+95)))
	v24 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+63)))
	v30 = base.Simd_g_v128_or(base.Simd_g_i8x16_shuffle2(v2, v9, base.Simd_g_const(&F_RD4_SSE2__k1), base.Simd_g_const(&F_RD4_SSE2__k2)), base.Simd_g_i32x4_replace_lane_l0(v2, v13<<(uint(int32(16))%32)|v18<<(uint(int32(24))%32)|v22|v24<<(uint(int32(8))%32)))
	v31 = base.Simd_g_const(&F_RD4_SSE2__k3)
	v32 = base.Simd_g_i8x16_shuffle2(v30, v2, base.Simd_g_const(&F_RD4_SSE2__k4), base.Simd_g_const(&F_RD4_SSE2__k5))
	v38 = base.Simd_g_const(&F_RD4_SSE2__k6)
	v40 = base.Simd_g_i8x16_avgr_u(base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_avgr_u(v32, v30), base.Simd_g_v128_and(base.Simd_g_v128_xor(v32, v30), base.Simd_g_const(&F_RD4_SSE2__k7))), base.Simd_g_i8x16_shuffle2(v30, v2, base.Simd_g_const(&F_RD4_SSE2__k8), base.Simd_g_const(&F_RD4_SSE2__k9)))
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(96), v40)
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(64), base.Simd_g_i8x16_shuffle2(v40, v2, base.Simd_g_const(&F_RD4_SSE2__k8), base.Simd_g_const(&F_RD4_SSE2__k9)))
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(32), base.Simd_g_i8x16_shuffle2(v40, v2, base.Simd_g_const(&F_RD4_SSE2__k4), base.Simd_g_const(&F_RD4_SSE2__k5)))
	base.Simd_g_v128_store32_lane_l0(m, l0, v8, base.Simd_g_i8x16_shuffle2(v40, v2, base.Simd_g_const(&F_RD4_SSE2__k10), base.Simd_g_const(&F_RD4_SSE2__k11)))
	return
}

var F_RD4_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_RD4_SSE2__k1 = [2]uint64{0x808080800f0e0d0c, 0x8080808080808080}
var F_RD4_SSE2__k2 = [2]uint64{0x302010080808080, 0xb0a090807060504}
var F_RD4_SSE2__k3 = [2]uint64{0x908070605040302, 0x11100f0e0d0c0b0a}
var F_RD4_SSE2__k4 = [2]uint64{0x908070605040302, 0x80800f0e0d0c0b0a}
var F_RD4_SSE2__k5 = [2]uint64{0x8080808080808080, 0x100808080808080}
var F_RD4_SSE2__k6 = [2]uint64{0x807060504030201, 0x100f0e0d0c0b0a09}
var F_RD4_SSE2__k7 = [2]uint64{0x101010101010101, 0x101010101010101}
var F_RD4_SSE2__k8 = [2]uint64{0x807060504030201, 0x800f0e0d0c0b0a09}
var F_RD4_SSE2__k9 = [2]uint64{0x8080808080808080, 0x80808080808080}
var F_RD4_SSE2__k10 = [2]uint64{0xa09080706050403, 0x8080800f0e0d0c0b}
var F_RD4_SSE2__k11 = [2]uint64{0x8080808080808080, 0x201008080808080}
