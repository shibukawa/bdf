//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_LD4_SSE2(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v3 base.V128
	_ = v3
	var v9 int64
	_ = v9
	var v10 int32
	_ = v10
	var v11 base.V128
	_ = v11
	var v13 base.V128
	_ = v13
	var v19 base.V128
	_ = v19
	var v25 base.V128
	_ = v25
	var v27 base.V128
	_ = v27
	v3 = base.Simd_g_const(&F_LD4_SSE2__k0)
	v9 = *(*int64)(unsafe.Add(mBase, uint32(l0+int32(-32))))
	v10 = int32(0)
	v11 = base.Simd_g_i64x2_replace_lane_l0(v3, v9)
	v13 = base.Simd_g_const(&F_LD4_SSE2__k1)
	v19 = base.Simd_g_i16x8_replace_lane_l3(base.Simd_g_i8x16_shuffle2(v11, v3, base.Simd_g_const(&F_LD4_SSE2__k2), base.Simd_g_const(&F_LD4_SSE2__k3)), base.I32_wrap_i64(int64(base.Ui64(v9)>>(uint(int64(56))%64))))
	v25 = base.Simd_g_const(&F_LD4_SSE2__k4)
	v27 = base.Simd_g_i8x16_avgr_u(base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_avgr_u(v11, v19), base.Simd_g_v128_and(base.Simd_g_v128_xor(v19, v11), base.Simd_g_const(&F_LD4_SSE2__k5))), base.Simd_g_i8x16_shuffle2(v11, v3, base.Simd_g_const(&F_LD4_SSE2__k6), base.Simd_g_const(&F_LD4_SSE2__k7)))
	base.Simd_g_v128_store32_lane_l0(m, l0, v10, v27)
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(96), base.Simd_g_i8x16_shuffle2(v27, v3, base.Simd_g_const(&F_LD4_SSE2__k8), base.Simd_g_const(&F_LD4_SSE2__k9)))
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(64), base.Simd_g_i8x16_shuffle2(v27, v3, base.Simd_g_const(&F_LD4_SSE2__k2), base.Simd_g_const(&F_LD4_SSE2__k3)))
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(32), base.Simd_g_i8x16_shuffle2(v27, v3, base.Simd_g_const(&F_LD4_SSE2__k6), base.Simd_g_const(&F_LD4_SSE2__k7)))
	return
}

var F_LD4_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_LD4_SSE2__k1 = [2]uint64{0x908070605040302, 0x11100f0e0d0c0b0a}
var F_LD4_SSE2__k2 = [2]uint64{0x908070605040302, 0x80800f0e0d0c0b0a}
var F_LD4_SSE2__k3 = [2]uint64{0x8080808080808080, 0x100808080808080}
var F_LD4_SSE2__k4 = [2]uint64{0x807060504030201, 0x100f0e0d0c0b0a09}
var F_LD4_SSE2__k5 = [2]uint64{0x101010101010101, 0x101010101010101}
var F_LD4_SSE2__k6 = [2]uint64{0x807060504030201, 0x800f0e0d0c0b0a09}
var F_LD4_SSE2__k7 = [2]uint64{0x8080808080808080, 0x80808080808080}
var F_LD4_SSE2__k8 = [2]uint64{0xa09080706050403, 0x8080800f0e0d0c0b}
var F_LD4_SSE2__k9 = [2]uint64{0x8080808080808080, 0x201008080808080}
