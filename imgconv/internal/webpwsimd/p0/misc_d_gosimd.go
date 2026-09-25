//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_DC16NoLeft_SSE2(m *base.Module, l0 int32) {
	var v2 base.V128
	_ = v2
	var v7 int32
	_ = v7
	var v8 base.V128
	_ = v8
	var v11 base.V128
	_ = v11
	var v12 int32
	_ = v12
	var v16 base.V128
	_ = v16
	var v19 base.V128
	_ = v19
	var v20 int32
	_ = v20
	var v24 base.V128
	_ = v24
	var v34 base.V128
	_ = v34
	v2 = base.Simd_g_const(&F_DC16NoLeft_SSE2__k0)
	v7 = int32(0)
	v8 = base.Simd_g_v128_load(m, l0+int32(-32), v7)
	v11 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v2, v8), base.Simd_g_i8x16_sub_sat_u(v8, v2))
	v12 = int32(8)
	v16 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_u(v11, v12), base.Simd_g_v128_and(v11, base.Simd_g_const(&F_DC16NoLeft_SSE2__k1)))
	v19 = base.Simd_g_i16x8_add(v16, base.Simd_g_i32x4_shl(v16, int32(16)))
	v20 = int32(32)
	v24 = base.Simd_g_i64x2_shr_u(base.Simd_g_i16x8_add(v19, base.Simd_g_i64x2_shl(v19, v20)), int32(48))
	v34 = base.Simd_g_i8x16_splat(int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i16x8_add(base.Simd_g_i8x16_swizzle_c(v24, base.Simd_g_const(&F_DC16NoLeft_SSE2__k2)), v24))+v12) >> (uint(int32(4)) % 32)))
	base.Simd_g_v128_store(m, l0, int32(480), v34)
	base.Simd_g_v128_store(m, l0, int32(448), v34)
	base.Simd_g_v128_store(m, l0, int32(416), v34)
	base.Simd_g_v128_store(m, l0, int32(384), v34)
	base.Simd_g_v128_store(m, l0, int32(352), v34)
	base.Simd_g_v128_store(m, l0, int32(320), v34)
	base.Simd_g_v128_store(m, l0, int32(288), v34)
	base.Simd_g_v128_store(m, l0, int32(256), v34)
	base.Simd_g_v128_store(m, l0, int32(224), v34)
	base.Simd_g_v128_store(m, l0, int32(192), v34)
	base.Simd_g_v128_store(m, l0, int32(160), v34)
	base.Simd_g_v128_store(m, l0, int32(128), v34)
	base.Simd_g_v128_store(m, l0, int32(96), v34)
	base.Simd_g_v128_store(m, l0, int32(64), v34)
	base.Simd_g_v128_store(m, l0, v20, v34)
	base.Simd_g_v128_store(m, l0, v7, v34)
	return
}

var F_DC16NoLeft_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_DC16NoLeft_SSE2__k1 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_DC16NoLeft_SSE2__k2 = [2]uint64{0x30201000b0a0908, 0x302010003020100}

func F_DC16NoTopLeft_SSE2(m *base.Module, l0 int32) {
	var v3 base.V128
	_ = v3
	v3 = base.Simd_g_const(&F_DC16NoTopLeft_SSE2__k0)
	base.Simd_g_v128_store(m, l0, int32(480), v3)
	base.Simd_g_v128_store(m, l0, int32(448), v3)
	base.Simd_g_v128_store(m, l0, int32(416), v3)
	base.Simd_g_v128_store(m, l0, int32(384), v3)
	base.Simd_g_v128_store(m, l0, int32(352), v3)
	base.Simd_g_v128_store(m, l0, int32(320), v3)
	base.Simd_g_v128_store(m, l0, int32(288), v3)
	base.Simd_g_v128_store(m, l0, int32(256), v3)
	base.Simd_g_v128_store(m, l0, int32(224), v3)
	base.Simd_g_v128_store(m, l0, int32(192), v3)
	base.Simd_g_v128_store(m, l0, int32(160), v3)
	base.Simd_g_v128_store(m, l0, int32(128), v3)
	base.Simd_g_v128_store(m, l0, int32(96), v3)
	base.Simd_g_v128_store(m, l0, int32(64), v3)
	base.Simd_g_v128_store(m, l0, int32(32), v3)
	base.Simd_g_v128_store(m, l0, int32(0), v3)
	return
}

var F_DC16NoTopLeft_SSE2__k0 = [2]uint64{0x8080808080808080, 0x8080808080808080}

func F_DC16NoTop_SSE2(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v8 int32
	_ = v8
	var v12 int32
	_ = v12
	var v16 int32
	_ = v16
	var v20 int32
	_ = v20
	var v24 int32
	_ = v24
	var v28 int32
	_ = v28
	var v32 int32
	_ = v32
	var v36 int32
	_ = v36
	var v40 int32
	_ = v40
	var v44 int32
	_ = v44
	var v48 int32
	_ = v48
	var v52 int32
	_ = v52
	var v56 int32
	_ = v56
	var v60 int32
	_ = v60
	var v64 int32
	_ = v64
	var v70 base.V128
	_ = v70
	v5 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-1)))))
	v8 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(31)))))
	v12 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(63)))))
	v16 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(95)))))
	v20 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(127)))))
	v24 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(159)))))
	v28 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(191)))))
	v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(223)))))
	v36 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(255)))))
	v40 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(287)))))
	v44 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(319)))))
	v48 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(351)))))
	v52 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(383)))))
	v56 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(415)))))
	v60 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(447)))))
	v64 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(479)))))
	v70 = base.Simd_g_i8x16_splat(int32(base.Ui32(v5+v8+v12+v16+v20+v24+v28+v32+v36+v40+v44+v48+v52+v56+v60+v64+int32(8)) >> (uint(int32(4)) % 32)))
	base.Simd_g_v128_store(m, l0, int32(480), v70)
	base.Simd_g_v128_store(m, l0, int32(448), v70)
	base.Simd_g_v128_store(m, l0, int32(416), v70)
	base.Simd_g_v128_store(m, l0, int32(384), v70)
	base.Simd_g_v128_store(m, l0, int32(352), v70)
	base.Simd_g_v128_store(m, l0, int32(320), v70)
	base.Simd_g_v128_store(m, l0, int32(288), v70)
	base.Simd_g_v128_store(m, l0, int32(256), v70)
	base.Simd_g_v128_store(m, l0, int32(224), v70)
	base.Simd_g_v128_store(m, l0, int32(192), v70)
	base.Simd_g_v128_store(m, l0, int32(160), v70)
	base.Simd_g_v128_store(m, l0, int32(128), v70)
	base.Simd_g_v128_store(m, l0, int32(96), v70)
	base.Simd_g_v128_store(m, l0, int32(64), v70)
	base.Simd_g_v128_store(m, l0, int32(32), v70)
	base.Simd_g_v128_store(m, l0, int32(0), v70)
	return
}
func F_DC16_SSE2(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 base.V128
	_ = v2
	var v7 int32
	_ = v7
	var v8 base.V128
	_ = v8
	var v11 base.V128
	_ = v11
	var v16 base.V128
	_ = v16
	var v17 int32
	_ = v17
	var v19 base.V128
	_ = v19
	var v20 int32
	_ = v20
	var v24 base.V128
	_ = v24
	var v32 int32
	_ = v32
	var v35 int32
	_ = v35
	var v39 int32
	_ = v39
	var v43 int32
	_ = v43
	var v47 int32
	_ = v47
	var v51 int32
	_ = v51
	var v55 int32
	_ = v55
	var v59 int32
	_ = v59
	var v63 int32
	_ = v63
	var v67 int32
	_ = v67
	var v71 int32
	_ = v71
	var v75 int32
	_ = v75
	var v79 int32
	_ = v79
	var v83 int32
	_ = v83
	var v87 int32
	_ = v87
	var v91 int32
	_ = v91
	var v98 base.V128
	_ = v98
	v2 = base.Simd_g_const(&F_DC16_SSE2__k0)
	v7 = int32(0)
	v8 = base.Simd_g_v128_load(m, l0+int32(-32), v7)
	v11 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v2, v8), base.Simd_g_i8x16_sub_sat_u(v8, v2))
	v16 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_u(v11, int32(8)), base.Simd_g_v128_and(v11, base.Simd_g_const(&F_DC16_SSE2__k1)))
	v17 = int32(16)
	v19 = base.Simd_g_i16x8_add(v16, base.Simd_g_i32x4_shl(v16, v17))
	v20 = int32(32)
	v24 = base.Simd_g_i64x2_shr_u(base.Simd_g_i16x8_add(v19, base.Simd_g_i64x2_shl(v19, v20)), int32(48))
	v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-1)))))
	v35 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(31)))))
	v39 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(63)))))
	v43 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(95)))))
	v47 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(127)))))
	v51 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(159)))))
	v55 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(191)))))
	v59 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(223)))))
	v63 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(255)))))
	v67 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(287)))))
	v71 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(319)))))
	v75 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(351)))))
	v79 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(383)))))
	v83 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(415)))))
	v87 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(447)))))
	v91 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(479)))))
	v98 = base.Simd_g_i8x16_splat(int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i16x8_add(base.Simd_g_i8x16_swizzle_c(v24, base.Simd_g_const(&F_DC16_SSE2__k2)), v24))+(v32+v35+v39+v43+v47+v51+v55+v59+v63+v67+v71+v75+v79+v83+v87+v91)+v17) >> (uint(int32(5)) % 32)))
	base.Simd_g_v128_store(m, l0, int32(480), v98)
	base.Simd_g_v128_store(m, l0, int32(448), v98)
	base.Simd_g_v128_store(m, l0, int32(416), v98)
	base.Simd_g_v128_store(m, l0, int32(384), v98)
	base.Simd_g_v128_store(m, l0, int32(352), v98)
	base.Simd_g_v128_store(m, l0, int32(320), v98)
	base.Simd_g_v128_store(m, l0, int32(288), v98)
	base.Simd_g_v128_store(m, l0, int32(256), v98)
	base.Simd_g_v128_store(m, l0, int32(224), v98)
	base.Simd_g_v128_store(m, l0, int32(192), v98)
	base.Simd_g_v128_store(m, l0, int32(160), v98)
	base.Simd_g_v128_store(m, l0, int32(128), v98)
	base.Simd_g_v128_store(m, l0, int32(96), v98)
	base.Simd_g_v128_store(m, l0, int32(64), v98)
	base.Simd_g_v128_store(m, l0, v20, v98)
	base.Simd_g_v128_store(m, l0, v7, v98)
	return
}

var F_DC16_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_DC16_SSE2__k1 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_DC16_SSE2__k2 = [2]uint64{0x30201000b0a0908, 0x302010003020100}

func F_DC8uvNoLeft_SSE2(m *base.Module, l0 int32) {
	var v2 base.V128
	_ = v2
	var v7 int32
	_ = v7
	var v8 base.V128
	_ = v8
	var v11 base.V128
	_ = v11
	var v16 base.V128
	_ = v16
	var v19 base.V128
	_ = v19
	var v20 int32
	_ = v20
	var v31 base.V128
	_ = v31
	v2 = base.Simd_g_const(&F_DC8uvNoLeft_SSE2__k0)
	v7 = int32(0)
	v8 = base.Simd_g_v128_load64_zero(m, l0+int32(-32), v7)
	v11 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v2, v8), base.Simd_g_i8x16_sub_sat_u(v8, v2))
	v16 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_u(v11, int32(8)), base.Simd_g_v128_and(v11, base.Simd_g_const(&F_DC8uvNoLeft_SSE2__k1)))
	v19 = base.Simd_g_i16x8_add(v16, base.Simd_g_i32x4_shl(v16, int32(16)))
	v20 = int32(32)
	v31 = base.Simd_g_i8x16_splat(int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i64x2_shr_u(base.Simd_g_i16x8_add(v19, base.Simd_g_i64x2_shl(v19, v20)), int32(48)))+int32(4)) >> (uint(int32(3)) % 32)))
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(224), v31)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(192), v31)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(160), v31)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(128), v31)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(96), v31)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(64), v31)
	base.Simd_g_v128_store64_lane_l0(m, l0, v20, v31)
	base.Simd_g_v128_store64_lane_l0(m, l0, v7, v31)
	return
}

var F_DC8uvNoLeft_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_DC8uvNoLeft_SSE2__k1 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}

func F_DC8uvNoTop_SSE2(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v8 int32
	_ = v8
	var v12 int32
	_ = v12
	var v16 int32
	_ = v16
	var v20 int32
	_ = v20
	var v24 int32
	_ = v24
	var v28 int32
	_ = v28
	var v32 int32
	_ = v32
	var v38 base.V128
	_ = v38
	var v40 int32
	_ = v40
	v5 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-1)))))
	v8 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(31)))))
	v12 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(63)))))
	v16 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(95)))))
	v20 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(127)))))
	v24 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(159)))))
	v28 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(191)))))
	v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(223)))))
	v38 = base.Simd_g_i8x16_splat(int32(base.Ui32(v5+v8+v12+v16+v20+v24+v28+v32+int32(4)) >> (uint(int32(3)) % 32)))
	v40 = int32(0)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(224), v38)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(192), v38)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(160), v38)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(128), v38)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(96), v38)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(64), v38)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(32), v38)
	base.Simd_g_v128_store64_lane_l0(m, l0, v40, v38)
	return
}
func F_DC8uv_SSE2(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 base.V128
	_ = v2
	var v7 int32
	_ = v7
	var v8 base.V128
	_ = v8
	var v11 base.V128
	_ = v11
	var v12 int32
	_ = v12
	var v16 base.V128
	_ = v16
	var v19 base.V128
	_ = v19
	var v20 int32
	_ = v20
	var v29 int32
	_ = v29
	var v32 int32
	_ = v32
	var v36 int32
	_ = v36
	var v40 int32
	_ = v40
	var v44 int32
	_ = v44
	var v48 int32
	_ = v48
	var v52 int32
	_ = v52
	var v56 int32
	_ = v56
	var v63 base.V128
	_ = v63
	v2 = base.Simd_g_const(&F_DC8uv_SSE2__k0)
	v7 = int32(0)
	v8 = base.Simd_g_v128_load64_zero(m, l0+int32(-32), v7)
	v11 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v2, v8), base.Simd_g_i8x16_sub_sat_u(v8, v2))
	v12 = int32(8)
	v16 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_u(v11, v12), base.Simd_g_v128_and(v11, base.Simd_g_const(&F_DC8uv_SSE2__k1)))
	v19 = base.Simd_g_i16x8_add(v16, base.Simd_g_i32x4_shl(v16, int32(16)))
	v20 = int32(32)
	v29 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-1)))))
	v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(31)))))
	v36 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(63)))))
	v40 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(95)))))
	v44 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(127)))))
	v48 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(159)))))
	v52 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(191)))))
	v56 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(223)))))
	v63 = base.Simd_g_i8x16_splat(int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i64x2_shr_u(base.Simd_g_i16x8_add(v19, base.Simd_g_i64x2_shl(v19, v20)), int32(48)))+(v29+v32+v36+v40+v44+v48+v52+v56)+v12) >> (uint(int32(4)) % 32)))
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(224), v63)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(192), v63)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(160), v63)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(128), v63)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(96), v63)
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(64), v63)
	base.Simd_g_v128_store64_lane_l0(m, l0, v20, v63)
	base.Simd_g_v128_store64_lane_l0(m, l0, v7, v63)
	return
}

var F_DC8uv_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_DC8uv_SSE2__k1 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}

func F_DispatchAlphaToGreen_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v19 int32
	_ = v19
	var v23 int32
	_ = v23
	var v27 int32
	_ = v27
	var v31 int32
	_ = v31
	var v41 int32
	_ = v41
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v49 base.V128
	_ = v49
	var v53 base.V128
	_ = v53
	var v59 int32
	_ = v59
	var v67 int32
	_ = v67
	var v81 int32
	_ = v81
	var v86 int32
	_ = v86
	var v88 int32
	_ = v88
	var v95 int32
	_ = v95
	var v111 int32
	_ = v111
	if l3 < int32(1) {
	} else {
		if l2 < int32(1) {
		} else {
			v19 = l2 & int32(2147483644)
			v23 = l0
			v27 = l4
			v31 = int32(0)
			for {
				if base.Ui32(l2) < base.Ui32(int32(4)) {
					v67 = int32(0)
					v81 = v67
					v86 = v27 + v67<<(uint(int32(2))%32)
					for {
						v88 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v23+v81))))
						*(*int32)(unsafe.Add(mBase, uint32(v86))) = v88 << (uint(int32(8)) % 32)
						v95 = v81 + int32(1)
						if l2 != v95 {
							v81 = v95
							v86 = v86 + int32(4)
							continue
						} else {
							break
						}
						break
					}
				} else {
					v41 = int32(0)
					v46 = v27
					for {
						v48 = int32(0)
						v49 = base.Simd_g_v128_load32_zero(m, v23+v41, v48)
						v53 = base.Simd_g_i32x4_shl(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v49)), int32(8))
						base.Simd_g_v128_store(m, v46, v48, v53)
						v59 = v41 + int32(4)
						if v19 != v59 {
							v41 = v59
							v46 = v46 + int32(16)
							continue
						} else {
							break
						}
						break
					}
					if v19 == l2 {
					} else {
						v67 = v19
						v81 = v67
						v86 = v27 + v67<<(uint(int32(2))%32)
						for {
							v88 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v23+v81))))
							*(*int32)(unsafe.Add(mBase, uint32(v86))) = v88 << (uint(int32(8)) % 32)
							v95 = v81 + int32(1)
							if l2 != v95 {
								v81 = v95
								v86 = v86 + int32(4)
								continue
							} else {
								break
							}
							break
						}
					}
				}
				v111 = v31 + int32(1)
				if v111 != l3 {
					v23 = v23 + l1
					v27 = v27 + l5<<(uint(int32(2))%32)
					v31 = v111
					continue
				} else {
					break
				}
				break
			}
		}
	}
	return
}
func F_DispatchAlphaToGreen_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v25 int32
	_ = v25
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v38 int32
	_ = v38
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v43 int32
	_ = v43
	var v49 int32
	_ = v49
	var v53 int32
	_ = v53
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v64 int32
	_ = v64
	var v76 int32
	_ = v76
	var v79 int32
	_ = v79
	var v87 int32
	_ = v87
	var v94 int32
	_ = v94
	var v95 base.V128
	_ = v95
	var v97 int32
	_ = v97
	var v98 base.V128
	_ = v98
	var v100 base.V128
	_ = v100
	var v101 base.V128
	_ = v101
	var v102 base.V128
	_ = v102
	var v105 int32
	_ = v105
	var v107 base.V128
	_ = v107
	var v108 base.V128
	_ = v108
	var v114 base.V128
	_ = v114
	var v116 base.V128
	_ = v116
	var v122 base.V128
	_ = v122
	var v130 int32
	_ = v130
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v156 int32
	_ = v156
	var v157 int32
	_ = v157
	var v158 base.V128
	_ = v158
	var v162 base.V128
	_ = v162
	var v170 int32
	_ = v170
	var v178 int32
	_ = v178
	var v202 int32
	_ = v202
	var v204 int32
	_ = v204
	var v219 int32
	_ = v219
	var v226 int32
	_ = v226
	var v254 int32
	_ = v254
	var v261 int32
	_ = v261
	var v265 int32
	_ = v265
	var v269 int32
	_ = v269
	var v281 int32
	_ = v281
	var v293 int32
	_ = v293
	var v295 int32
	_ = v295
	var v310 int32
	_ = v310
	var v311 base.V128
	_ = v311
	var v315 base.V128
	_ = v315
	var v321 int32
	_ = v321
	var v329 int32
	_ = v329
	var v353 int32
	_ = v353
	var v355 int32
	_ = v355
	var v370 int32
	_ = v370
	var v377 int32
	_ = v377
	var v403 int32
	_ = v403
	if l3 < int32(1) {
	} else {
		v25 = l2 & int32(-16)
		if v25 < int32(1) {
			if l2 < int32(1) {
			} else {
				v261 = l2 & int32(2147483644)
				v265 = l0
				v269 = l4
				v281 = int32(0)
				for {
					if base.Ui32(l2) < base.Ui32(int32(4)) {
						v329 = int32(0)
						v353 = v329
						v355 = v269 + v329<<(uint(int32(2))%32)
						for {
							v370 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v265+v353))))
							*(*int32)(unsafe.Add(mBase, uint32(v355))) = v370 << (uint(int32(8)) % 32)
							v377 = v353 + int32(1)
							if l2 != v377 {
								v353 = v377
								v355 = v355 + int32(4)
								continue
							} else {
								break
							}
							break
						}
					} else {
						v293 = int32(0)
						v295 = v269
						for {
							v310 = int32(0)
							v311 = base.Simd_g_v128_load32_zero(m, v265+v293, v310)
							v315 = base.Simd_g_i32x4_shl(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v311)), int32(8))
							base.Simd_g_v128_store(m, v295, v310, v315)
							v321 = v293 + int32(4)
							if v261 != v321 {
								v293 = v321
								v295 = v295 + int32(16)
								continue
							} else {
								break
							}
							break
						}
						if v261 == l2 {
						} else {
							v329 = v261
							v353 = v329
							v355 = v269 + v329<<(uint(int32(2))%32)
							for {
								v370 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v265+v353))))
								*(*int32)(unsafe.Add(mBase, uint32(v355))) = v370 << (uint(int32(8)) % 32)
								v377 = v353 + int32(1)
								if l2 != v377 {
									v353 = v377
									v355 = v355 + int32(4)
									continue
								} else {
									break
								}
								break
							}
						}
					}
					v403 = v281 + int32(1)
					if v403 != l3 {
						v265 = v265 + l1
						v269 = v269 + l5<<(uint(int32(2))%32)
						v281 = v403
						continue
					} else {
						break
					}
					break
				}
			}
		} else {
			v31 = (v25 + int32(-1)) & int32(-16)
			v33 = v31 + int32(17)
			if v33 < l2 {
				v35 = l2
			} else {
				v35 = v33
			}
			v38 = v35 - v31 + int32(-16)
			v40 = v35 & int32(3)
			v41 = v38 - v40
			v43 = l5 << (uint(int32(2)) % 32)
			v49 = l0
			v53 = l4
			v62 = l4 + int32(64)
			v63 = l0 + int32(16)
			v64 = int32(0)
			for {
				v76 = int32(0)
				v79 = v62
				v87 = v63
				for {
					v94 = v53 + v76<<(uint(int32(2))%32)
					v95 = base.Simd_g_const(&F_DispatchAlphaToGreen_SSE2__k0)
					v97 = int32(0)
					v98 = base.Simd_g_v128_load(m, v49+v76, v97)
					v100 = base.Simd_g_i8x16_shuffle2(v95, v98, base.Simd_g_const(&F_DispatchAlphaToGreen_SSE2__k1), base.Simd_g_const(&F_DispatchAlphaToGreen_SSE2__k2))
					v101 = base.Simd_g_const(&F_DispatchAlphaToGreen_SSE2__k3)
					v102 = base.Simd_g_i8x16_shuffle2(v100, v95, base.Simd_g_const(&F_DispatchAlphaToGreen_SSE2__k4), base.Simd_g_const(&F_DispatchAlphaToGreen_SSE2__k5))
					base.Simd_g_v128_store(m, v94, v97, v102)
					v105 = int32(16)
					v107 = base.Simd_g_const(&F_DispatchAlphaToGreen_SSE2__k6)
					v108 = base.Simd_g_i8x16_shuffle2(v100, v95, base.Simd_g_const(&F_DispatchAlphaToGreen_SSE2__k7), base.Simd_g_const(&F_DispatchAlphaToGreen_SSE2__k8))
					base.Simd_g_v128_store(m, v94+v105, v97, v108)
					v114 = base.Simd_g_i8x16_shuffle2(v95, v98, base.Simd_g_const(&F_DispatchAlphaToGreen_SSE2__k9), base.Simd_g_const(&F_DispatchAlphaToGreen_SSE2__k10))
					v116 = base.Simd_g_i8x16_shuffle2(v114, v95, base.Simd_g_const(&F_DispatchAlphaToGreen_SSE2__k4), base.Simd_g_const(&F_DispatchAlphaToGreen_SSE2__k5))
					base.Simd_g_v128_store(m, v94+int32(32), v97, v116)
					v122 = base.Simd_g_i8x16_shuffle2(v114, v95, base.Simd_g_const(&F_DispatchAlphaToGreen_SSE2__k7), base.Simd_g_const(&F_DispatchAlphaToGreen_SSE2__k8))
					base.Simd_g_v128_store(m, v94+int32(48), v97, v122)
					v130 = v76 + v105
					if v130 < v25 {
						v76 = v130
						v79 = v79 + int32(64)
						v87 = v87 + v105
						continue
					} else {
						break
					}
					break
				}
				if l2 <= v130 {
				} else {
					if base.Ui32(v38) <= base.Ui32(int32(3)) {
						v178 = v130
						v202 = v178
						v204 = v53 + v178<<(uint(int32(2))%32)
						for {
							v219 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v49+v202))))
							*(*int32)(unsafe.Add(mBase, uint32(v204))) = v219 << (uint(int32(8)) % 32)
							v226 = v202 + int32(1)
							if v226 < l2 {
								v202 = v226
								v204 = v204 + int32(4)
								continue
							} else {
								break
							}
							break
						}
					} else {
						v143 = v79
						v144 = v41
						v156 = v87
						for {
							v157 = int32(0)
							v158 = base.Simd_g_v128_load32_zero(m, v156, v157)
							v162 = base.Simd_g_i32x4_shl(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v158)), int32(8))
							base.Simd_g_v128_store(m, v143, v157, v162)
							v170 = v144 + int32(-4)
							if v170 != 0 {
								v143 = v143 + int32(16)
								v144 = v170
								v156 = v156 + int32(4)
								continue
							} else {
								break
							}
							break
						}
						if v40 == int32(0) {
						} else {
							v178 = v130 + v41
							v202 = v178
							v204 = v53 + v178<<(uint(int32(2))%32)
							for {
								v219 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v49+v202))))
								*(*int32)(unsafe.Add(mBase, uint32(v204))) = v219 << (uint(int32(8)) % 32)
								v226 = v202 + int32(1)
								if v226 < l2 {
									v202 = v226
									v204 = v204 + int32(4)
									continue
								} else {
									break
								}
								break
							}
						}
					}
				}
				v254 = v64 + int32(1)
				if v254 != l3 {
					v49 = v49 + l1
					v53 = v53 + v43
					v62 = v62 + v43
					v63 = v63 + l1
					v64 = v254
					continue
				} else {
					break
				}
				break
			}
		}
	}
	return
}

var F_DispatchAlphaToGreen_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_DispatchAlphaToGreen_SSE2__k1 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_DispatchAlphaToGreen_SSE2__k2 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_DispatchAlphaToGreen_SSE2__k3 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_DispatchAlphaToGreen_SSE2__k4 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_DispatchAlphaToGreen_SSE2__k5 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_DispatchAlphaToGreen_SSE2__k6 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_DispatchAlphaToGreen_SSE2__k7 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_DispatchAlphaToGreen_SSE2__k8 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_DispatchAlphaToGreen_SSE2__k9 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_DispatchAlphaToGreen_SSE2__k10 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}

func F_DispatchAlpha_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) int32 {
	mBase := m.M
	_ = mBase
	var v7 int32
	_ = v7
	var v24 int32
	_ = v24
	var v29 int32
	_ = v29
	var v33 int32
	_ = v33
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v50 int32
	_ = v50
	var v60 int32
	_ = v60
	var v65 base.V128
	_ = v65
	var v66 base.V128
	_ = v66
	var v71 int32
	_ = v71
	var v72 base.V128
	_ = v72
	var v73 int32
	_ = v73
	var v78 base.V128
	_ = v78
	var v82 int32
	_ = v82
	var v94 int32
	_ = v94
	var v104 base.V128
	_ = v104
	var v106 int32
	_ = v106
	var v110 base.V128
	_ = v110
	var v115 int32
	_ = v115
	var v123 int32
	_ = v123
	var v125 int32
	_ = v125
	var v143 int32
	_ = v143
	var v145 int32
	_ = v145
	var v152 int32
	_ = v152
	var v155 int32
	_ = v155
	var v159 int32
	_ = v159
	var v161 int32
	_ = v161
	var v171 int32
	_ = v171
	var v183 int32
	_ = v183
	var v193 int32
	_ = v193
	v7 = int32(0)
	if l3 < int32(1) {
		v193 = v7
	} else {
		if l2 < int32(1) {
			v193 = v7
		} else {
			v24 = l2 & int32(2147483644)
			v29 = l0
			v33 = l4
			v37 = int32(255)
			v38 = int32(0)
			for {
				if base.B2i32(base.Ui32(l2) < base.Ui32(int32(4))) == int32(0) {
					v50 = int32(0)
					v60 = v50
					v65 = base.Simd_g_i32x4_replace_lane_l0(base.Simd_g_const(&F_DispatchAlpha_C__k0), v37)
					v66 = base.Simd_g_const(&F_DispatchAlpha_C__k1)
					for {
						v71 = int32(2)
						v72 = base.Simd_g_i32x4_shl(v66, v71)
						v73 = int32(0)
						v78 = base.Simd_g_v128_load32_zero(m, v29+v60, v73)
						base.Simd_g_v128_store8_lane_l0(m, v33+base.Simd_g_i32x4_extract_lane_l0(v72), v73, v78)
						v82 = int32(1)
						base.Simd_g_v128_store8_lane_l1(m, v33+base.Simd_g_i32x4_extract_lane_l1(v72), v73, v78)
						base.Simd_g_v128_store8_lane_l2(m, v33+base.Simd_g_i32x4_extract_lane_l2(v72), v73, v78)
						v94 = int32(3)
						base.Simd_g_v128_store8_lane_l3(m, v33+base.Simd_g_i32x4_extract_lane_l3(v72), v73, v78)
						v104 = base.Simd_g_v128_and(v65, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v78)))
						v106 = v60 + int32(4)
						if v24 != v106 {
							v60 = v106
							v65 = v104
							v66 = base.Simd_g_i32x4_add(v66, base.Simd_g_const(&F_DispatchAlpha_C__k2))
							continue
						} else {
							break
						}
						break
					}
					v110 = base.Simd_g_v128_and(v104, base.Simd_g_i8x16_shuffle2(v104, v78, base.Simd_g_const(&F_DispatchAlpha_C__k3), base.Simd_g_const(&F_DispatchAlpha_C__k4)))
					v115 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_v128_and(v110, base.Simd_g_i8x16_shuffle2(v110, v110, base.Simd_g_const(&F_DispatchAlpha_C__k5), base.Simd_g_const(&F_DispatchAlpha_C__k4))))
					if v24 == l2 {
						v171 = v115
					} else {
						v123 = v24
						v125 = v115
						v143 = v123
						v145 = v125
						v152 = v33 + v123<<(uint(int32(2))%32)
						for {
							v155 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v29+v143))))
							*(*uint8)(unsafe.Add(mBase, uint32(v152))) = uint8(v155)
							v159 = v145 & v155
							v161 = v143 + int32(1)
							if l2 != v161 {
								v143 = v161
								v145 = v159
								v152 = v152 + int32(4)
								continue
							} else {
								break
							}
							break
						}
						v171 = v159
					}
				} else {
					v123 = int32(0)
					v125 = v37
					v143 = v123
					v145 = v125
					v152 = v33 + v123<<(uint(int32(2))%32)
					for {
						v155 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v29+v143))))
						*(*uint8)(unsafe.Add(mBase, uint32(v152))) = uint8(v155)
						v159 = v145 & v155
						v161 = v143 + int32(1)
						if l2 != v161 {
							v143 = v161
							v145 = v159
							v152 = v152 + int32(4)
							continue
						} else {
							break
						}
						break
					}
					v171 = v159
				}
				v183 = v38 + int32(1)
				if v183 != l3 {
					v29 = v29 + l1
					v33 = v33 + l5
					v37 = v171
					v38 = v183
					continue
				} else {
					break
				}
				break
			}
			v193 = base.B2i32(v171 != int32(255))
		}
	}
	return v193
}

var F_DispatchAlpha_C__k0 = [2]uint64{0xffffffff00000000, 0xffffffffffffffff}
var F_DispatchAlpha_C__k1 = [2]uint64{0x100000000, 0x300000002}
var F_DispatchAlpha_C__k2 = [2]uint64{0x400000004, 0x400000004}
var F_DispatchAlpha_C__k3 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_DispatchAlpha_C__k4 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_DispatchAlpha_C__k5 = [2]uint64{0x302010007060504, 0x302010003020100}

func F_DispatchAlpha_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) int32 {
	mBase := m.M
	_ = mBase
	var v23 base.V128
	_ = v23
	var v29 base.V128
	_ = v29
	var v30 int32
	_ = v30
	var v34 int32
	_ = v34
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v42 base.V128
	_ = v42
	var v43 base.V128
	_ = v43
	var v62 int32
	_ = v62
	var v66 base.V128
	_ = v66
	var v69 int32
	_ = v69
	var v75 int32
	_ = v75
	var v76 base.V128
	_ = v76
	var v77 base.V128
	_ = v77
	var v79 base.V128
	_ = v79
	var v80 base.V128
	_ = v80
	var v81 base.V128
	_ = v81
	var v83 int32
	_ = v83
	var v86 int32
	_ = v86
	var v89 int32
	_ = v89
	var v94 base.V128
	_ = v94
	var v95 base.V128
	_ = v95
	var v105 int32
	_ = v105
	var v109 base.V128
	_ = v109
	var v111 base.V128
	_ = v111
	var v121 int32
	_ = v121
	var v125 base.V128
	_ = v125
	var v139 int32
	_ = v139
	var v140 base.V128
	_ = v140
	var v144 int32
	_ = v144
	var v154 int32
	_ = v154
	var v158 base.V128
	_ = v158
	var v160 int32
	_ = v160
	var v167 int32
	_ = v167
	var v170 int32
	_ = v170
	var v171 base.V128
	_ = v171
	var v172 base.V128
	_ = v172
	var v174 base.V128
	_ = v174
	var v176 base.V128
	_ = v176
	var v178 int32
	_ = v178
	var v181 int32
	_ = v181
	var v184 int32
	_ = v184
	var v190 base.V128
	_ = v190
	var v206 base.V128
	_ = v206
	var v207 int32
	_ = v207
	var v211 int32
	_ = v211
	var v222 int32
	_ = v222
	var v232 int32
	_ = v232
	var v240 base.V128
	_ = v240
	var v241 base.V128
	_ = v241
	var v242 int32
	_ = v242
	var v244 int32
	_ = v244
	var v245 base.V128
	_ = v245
	var v246 int32
	_ = v246
	var v250 base.V128
	_ = v250
	var v254 int32
	_ = v254
	var v266 int32
	_ = v266
	var v278 base.V128
	_ = v278
	var v280 int32
	_ = v280
	var v283 base.V128
	_ = v283
	var v288 int32
	_ = v288
	var v300 int32
	_ = v300
	var v305 int32
	_ = v305
	var v321 int32
	_ = v321
	var v323 int32
	_ = v323
	var v328 int32
	_ = v328
	var v334 int32
	_ = v334
	var v338 int32
	_ = v338
	var v340 int32
	_ = v340
	var v352 int32
	_ = v352
	var v365 int32
	_ = v365
	var v367 base.V128
	_ = v367
	var v378 base.V128
	_ = v378
	var v379 base.V128
	_ = v379
	var v380 int32
	_ = v380
	if int32(1) <= l3 {
		v29 = base.Simd_g_const(&F_DispatchAlpha_SSE2__k0)
		v30 = l0
		v34 = l4
		v39 = int32(0)
		v40 = int32(255)
		v42 = v29
		v43 = v29
		for {
			if base.B2i32(l2 < int32(17)) == int32(0) {
				v62 = v34
				v66 = v42
				v69 = int32(0)
				for {
					v75 = int32(0)
					v76 = base.Simd_g_v128_load(m, v30+v69, v75)
					v77 = base.Simd_g_const(&F_DispatchAlpha_SSE2__k1)
					v79 = base.Simd_g_i8x16_shuffle2(v76, v77, base.Simd_g_const(&F_DispatchAlpha_SSE2__k2), base.Simd_g_const(&F_DispatchAlpha_SSE2__k3))
					v80 = base.Simd_g_const(&F_DispatchAlpha_SSE2__k4)
					v81 = base.Simd_g_i8x16_shuffle2(v79, v77, base.Simd_g_const(&F_DispatchAlpha_SSE2__k5), base.Simd_g_const(&F_DispatchAlpha_SSE2__k6))
					v83 = int32(12)
					base.Simd_g_v128_store8_lane_l12(m, v62, int32(60), v81)
					v86 = int32(8)
					base.Simd_g_v128_store8_lane_l8(m, v62, int32(56), v81)
					v89 = int32(4)
					base.Simd_g_v128_store8_lane_l4(m, v62, int32(52), v81)
					base.Simd_g_v128_store8_lane_l0(m, v62, int32(48), v81)
					v94 = base.Simd_g_const(&F_DispatchAlpha_SSE2__k7)
					v95 = base.Simd_g_i8x16_shuffle2(v79, v77, base.Simd_g_const(&F_DispatchAlpha_SSE2__k8), base.Simd_g_const(&F_DispatchAlpha_SSE2__k9))
					base.Simd_g_v128_store8_lane_l12(m, v62, int32(44), v95)
					base.Simd_g_v128_store8_lane_l8(m, v62, int32(40), v95)
					base.Simd_g_v128_store8_lane_l4(m, v62, int32(36), v95)
					v105 = int32(32)
					base.Simd_g_v128_store8_lane_l0(m, v62, v105, v95)
					v109 = base.Simd_g_i8x16_shuffle2(v76, v77, base.Simd_g_const(&F_DispatchAlpha_SSE2__k10), base.Simd_g_const(&F_DispatchAlpha_SSE2__k11))
					v111 = base.Simd_g_i8x16_shuffle2(v109, v77, base.Simd_g_const(&F_DispatchAlpha_SSE2__k5), base.Simd_g_const(&F_DispatchAlpha_SSE2__k6))
					base.Simd_g_v128_store8_lane_l12(m, v62, int32(28), v111)
					base.Simd_g_v128_store8_lane_l8(m, v62, int32(24), v111)
					base.Simd_g_v128_store8_lane_l4(m, v62, int32(20), v111)
					v121 = int32(16)
					base.Simd_g_v128_store8_lane_l0(m, v62, v121, v111)
					v125 = base.Simd_g_i8x16_shuffle2(v109, v77, base.Simd_g_const(&F_DispatchAlpha_SSE2__k8), base.Simd_g_const(&F_DispatchAlpha_SSE2__k9))
					base.Simd_g_v128_store8_lane_l12(m, v62, v83, v125)
					base.Simd_g_v128_store8_lane_l8(m, v62, v86, v125)
					base.Simd_g_v128_store8_lane_l4(m, v62, v89, v125)
					base.Simd_g_v128_store8_lane_l0(m, v62, v75, v125)
					v139 = v62 + int32(64)
					v140 = base.Simd_g_v128_and(v76, v66)
					v144 = v69 + v121
					if v69+v105 < l2 {
						v62 = v139
						v66 = v140
						v69 = v144
						continue
					} else {
						break
					}
					break
				}
				v154 = v139
				v158 = v140
				v160 = v144
			} else {
				v154 = v34
				v158 = v42
				v160 = int32(0)
			}
			v167 = v160 | int32(8)
			if v167 < l2 {
				v170 = int32(0)
				v171 = base.Simd_g_v128_load64_zero(m, v30+v160, v170)
				v172 = base.Simd_g_const(&F_DispatchAlpha_SSE2__k1)
				v174 = base.Simd_g_i8x16_shuffle2(v171, v172, base.Simd_g_const(&F_DispatchAlpha_SSE2__k10), base.Simd_g_const(&F_DispatchAlpha_SSE2__k11))
				v176 = base.Simd_g_i8x16_shuffle2(v174, v172, base.Simd_g_const(&F_DispatchAlpha_SSE2__k5), base.Simd_g_const(&F_DispatchAlpha_SSE2__k6))
				v178 = int32(12)
				base.Simd_g_v128_store8_lane_l12(m, v154, int32(28), v176)
				v181 = int32(8)
				base.Simd_g_v128_store8_lane_l8(m, v154, int32(24), v176)
				v184 = int32(4)
				base.Simd_g_v128_store8_lane_l4(m, v154, int32(20), v176)
				base.Simd_g_v128_store8_lane_l0(m, v154, int32(16), v176)
				v190 = base.Simd_g_i8x16_shuffle2(v174, v172, base.Simd_g_const(&F_DispatchAlpha_SSE2__k8), base.Simd_g_const(&F_DispatchAlpha_SSE2__k9))
				base.Simd_g_v128_store8_lane_l12(m, v154, v178, v190)
				base.Simd_g_v128_store8_lane_l8(m, v154, v181, v190)
				base.Simd_g_v128_store8_lane_l4(m, v154, v184, v190)
				base.Simd_g_v128_store8_lane_l0(m, v154, v170, v190)
				v206 = base.Simd_g_v128_and(v43, v171)
				v207 = v167
			} else {
				v206 = v43
				v207 = v160
			}
			if l2 <= v207 {
				v352 = v40
			} else {
				v211 = l2 - v207
				if base.Ui32(v211) < base.Ui32(int32(4)) {
					v300 = v40
					v305 = v207
					v321 = v34 + v305<<(uint(int32(2))%32)
					v323 = v300
					v328 = v305
					for {
						v334 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v30+v328))))
						*(*uint8)(unsafe.Add(mBase, uint32(v321))) = uint8(v334)
						v338 = v323 & v334
						v340 = v328 + int32(1)
						if l2 != v340 {
							v321 = v321 + int32(4)
							v323 = v338
							v328 = v340
							continue
						} else {
							break
						}
						break
					}
					v352 = v338
				} else {
					v222 = v211 & int32(-4)
					v232 = v30 + v207
					v240 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_splat(v207), base.Simd_g_const(&F_DispatchAlpha_SSE2__k12))
					v241 = base.Simd_g_i32x4_replace_lane_l0(base.Simd_g_const(&F_DispatchAlpha_SSE2__k13), v40)
					v242 = v222
					for {
						v244 = int32(2)
						v245 = base.Simd_g_i32x4_shl(v240, v244)
						v246 = int32(0)
						v250 = base.Simd_g_v128_load32_zero(m, v232, v246)
						base.Simd_g_v128_store8_lane_l0(m, v34+base.Simd_g_i32x4_extract_lane_l0(v245), v246, v250)
						v254 = int32(1)
						base.Simd_g_v128_store8_lane_l1(m, v34+base.Simd_g_i32x4_extract_lane_l1(v245), v246, v250)
						base.Simd_g_v128_store8_lane_l2(m, v34+base.Simd_g_i32x4_extract_lane_l2(v245), v246, v250)
						v266 = int32(3)
						base.Simd_g_v128_store8_lane_l3(m, v34+base.Simd_g_i32x4_extract_lane_l3(v245), v246, v250)
						v278 = base.Simd_g_v128_and(v241, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v250)))
						v280 = v242 + int32(-4)
						if v280 != 0 {
							v232 = v232 + int32(4)
							v240 = base.Simd_g_i32x4_add(v240, base.Simd_g_const(&F_DispatchAlpha_SSE2__k14))
							v241 = v278
							v242 = v280
							continue
						} else {
							break
						}
						break
					}
					v283 = base.Simd_g_v128_and(v278, base.Simd_g_i8x16_shuffle2(v278, v250, base.Simd_g_const(&F_DispatchAlpha_SSE2__k15), base.Simd_g_const(&F_DispatchAlpha_SSE2__k16)))
					v288 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_v128_and(v283, base.Simd_g_i8x16_shuffle2(v283, v283, base.Simd_g_const(&F_DispatchAlpha_SSE2__k17), base.Simd_g_const(&F_DispatchAlpha_SSE2__k16))))
					if v211 == v222 {
						v352 = v288
					} else {
						v300 = v288
						v305 = v207 + v222
						v321 = v34 + v305<<(uint(int32(2))%32)
						v323 = v300
						v328 = v305
						for {
							v334 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v30+v328))))
							*(*uint8)(unsafe.Add(mBase, uint32(v321))) = uint8(v334)
							v338 = v323 & v334
							v340 = v328 + int32(1)
							if l2 != v340 {
								v321 = v321 + int32(4)
								v323 = v338
								v328 = v340
								continue
							} else {
								break
							}
							break
						}
						v352 = v338
					}
				}
			}
			v365 = v39 + int32(1)
			if v365 != l3 {
				v30 = v30 + l1
				v34 = v34 + l5
				v39 = v365
				v40 = v352
				v42 = v158
				v43 = v206
				continue
			} else {
				break
			}
			break
		}
		v367 = base.Simd_g_const(&F_DispatchAlpha_SSE2__k0)
		v378 = base.Simd_g_i8x16_eq(v206, v367)
		v379 = base.Simd_g_i8x16_eq(v158, v367)
		v380 = v352 & int32(255)
	} else {
		v23 = base.Simd_g_const(&F_DispatchAlpha_SSE2__k0)
		v378 = v23
		v379 = v23
		v380 = int32(255)
	}
	return base.B2i32(v380&base.Simd_g_i8x16_bitmask(v378) != int32(255)) | base.B2i32(base.Simd_g_i8x16_bitmask(v379) != int32(_a_F_DispatchAlpha_SSE2_0))
}

var F_DispatchAlpha_SSE2__k0 = [2]uint64{0xffffffffffffffff, 0xffffffffffffffff}
var F_DispatchAlpha_SSE2__k1 = [2]uint64{0x0, 0x0}
var F_DispatchAlpha_SSE2__k2 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_DispatchAlpha_SSE2__k3 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_DispatchAlpha_SSE2__k4 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_DispatchAlpha_SSE2__k5 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_DispatchAlpha_SSE2__k6 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_DispatchAlpha_SSE2__k7 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_DispatchAlpha_SSE2__k8 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_DispatchAlpha_SSE2__k9 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_DispatchAlpha_SSE2__k10 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_DispatchAlpha_SSE2__k11 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_DispatchAlpha_SSE2__k12 = [2]uint64{0x100000000, 0x300000002}
var F_DispatchAlpha_SSE2__k13 = [2]uint64{0xffffffff00000000, 0xffffffffffffffff}
var F_DispatchAlpha_SSE2__k14 = [2]uint64{0x400000004, 0x400000004}
var F_DispatchAlpha_SSE2__k15 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_DispatchAlpha_SSE2__k16 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_DispatchAlpha_SSE2__k17 = [2]uint64{0x302010007060504, 0x302010003020100}

func F_Disto16x16_C(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v82 int32
	_ = v82
	var v83 base.V128
	_ = v83
	var v84 int32
	_ = v84
	var v85 base.V128
	_ = v85
	var v86 int32
	_ = v86
	var v87 base.V128
	_ = v87
	var v88 int32
	_ = v88
	var v89 base.V128
	_ = v89
	var v90 int32
	_ = v90
	var v91 base.V128
	_ = v91
	var v92 int32
	_ = v92
	var v93 base.V128
	_ = v93
	var v94 int32
	_ = v94
	var v95 base.V128
	_ = v95
	var v96 int32
	_ = v96
	var v97 base.V128
	_ = v97
	var v98 int32
	_ = v98
	var v99 base.V128
	_ = v99
	var v100 int32
	_ = v100
	var v101 base.V128
	_ = v101
	var v102 int32
	_ = v102
	var v103 base.V128
	_ = v103
	var v104 int32
	_ = v104
	var v105 base.V128
	_ = v105
	var v106 int32
	_ = v106
	var v107 base.V128
	_ = v107
	var v108 int32
	_ = v108
	var v109 base.V128
	_ = v109
	var v110 int32
	_ = v110
	var v111 base.V128
	_ = v111
	var v112 int32
	_ = v112
	var v113 base.V128
	_ = v113
	var v135 int32
	_ = v135
	var v136 int32
	_ = v136
	var v197 int32
	_ = v197
	var v198 int32
	_ = v198
	var v200 int32
	_ = v200
	var v202 int32
	_ = v202
	var v204 int32
	_ = v204
	var v206 int32
	_ = v206
	var v207 base.V128
	_ = v207
	var v209 int32
	_ = v209
	var v210 base.V128
	_ = v210
	var v212 int32
	_ = v212
	var v213 base.V128
	_ = v213
	var v215 int32
	_ = v215
	var v216 base.V128
	_ = v216
	var v218 base.V128
	_ = v218
	var v219 int32
	_ = v219
	var v221 int32
	_ = v221
	var v223 int32
	_ = v223
	var v225 int32
	_ = v225
	var v228 base.V128
	_ = v228
	var v231 base.V128
	_ = v231
	var v234 base.V128
	_ = v234
	var v237 base.V128
	_ = v237
	var v239 base.V128
	_ = v239
	var v240 base.V128
	_ = v240
	var v241 int32
	_ = v241
	var v243 int32
	_ = v243
	var v245 int32
	_ = v245
	var v247 int32
	_ = v247
	var v250 base.V128
	_ = v250
	var v253 base.V128
	_ = v253
	var v256 base.V128
	_ = v256
	var v259 base.V128
	_ = v259
	var v261 base.V128
	_ = v261
	var v262 int32
	_ = v262
	var v264 int32
	_ = v264
	var v266 int32
	_ = v266
	var v268 int32
	_ = v268
	var v271 base.V128
	_ = v271
	var v274 base.V128
	_ = v274
	var v277 base.V128
	_ = v277
	var v280 base.V128
	_ = v280
	var v282 base.V128
	_ = v282
	var v283 base.V128
	_ = v283
	var v284 base.V128
	_ = v284
	var v285 int32
	_ = v285
	var v287 int32
	_ = v287
	var v289 int32
	_ = v289
	var v291 int32
	_ = v291
	var v294 base.V128
	_ = v294
	var v297 base.V128
	_ = v297
	var v300 base.V128
	_ = v300
	var v303 base.V128
	_ = v303
	var v305 base.V128
	_ = v305
	var v306 int32
	_ = v306
	var v308 int32
	_ = v308
	var v310 int32
	_ = v310
	var v312 int32
	_ = v312
	var v315 base.V128
	_ = v315
	var v318 base.V128
	_ = v318
	var v321 base.V128
	_ = v321
	var v324 base.V128
	_ = v324
	var v326 base.V128
	_ = v326
	var v327 base.V128
	_ = v327
	var v328 int32
	_ = v328
	var v330 int32
	_ = v330
	var v332 int32
	_ = v332
	var v334 int32
	_ = v334
	var v337 base.V128
	_ = v337
	var v340 base.V128
	_ = v340
	var v343 base.V128
	_ = v343
	var v346 base.V128
	_ = v346
	var v348 base.V128
	_ = v348
	var v349 int32
	_ = v349
	var v351 int32
	_ = v351
	var v353 int32
	_ = v353
	var v355 int32
	_ = v355
	var v358 base.V128
	_ = v358
	var v361 base.V128
	_ = v361
	var v364 base.V128
	_ = v364
	var v367 base.V128
	_ = v367
	var v369 base.V128
	_ = v369
	var v370 base.V128
	_ = v370
	var v371 base.V128
	_ = v371
	var v372 base.V128
	_ = v372
	var v373 int32
	_ = v373
	var v375 int32
	_ = v375
	var v377 int32
	_ = v377
	var v379 int32
	_ = v379
	var v382 base.V128
	_ = v382
	var v385 base.V128
	_ = v385
	var v388 base.V128
	_ = v388
	var v391 base.V128
	_ = v391
	var v393 base.V128
	_ = v393
	var v394 int32
	_ = v394
	var v396 int32
	_ = v396
	var v398 int32
	_ = v398
	var v400 int32
	_ = v400
	var v403 base.V128
	_ = v403
	var v406 base.V128
	_ = v406
	var v409 base.V128
	_ = v409
	var v412 base.V128
	_ = v412
	var v414 base.V128
	_ = v414
	var v415 base.V128
	_ = v415
	var v416 int32
	_ = v416
	var v418 int32
	_ = v418
	var v420 int32
	_ = v420
	var v422 int32
	_ = v422
	var v425 base.V128
	_ = v425
	var v428 base.V128
	_ = v428
	var v431 base.V128
	_ = v431
	var v434 base.V128
	_ = v434
	var v436 base.V128
	_ = v436
	var v437 int32
	_ = v437
	var v439 int32
	_ = v439
	var v441 int32
	_ = v441
	var v443 int32
	_ = v443
	var v446 base.V128
	_ = v446
	var v449 base.V128
	_ = v449
	var v452 base.V128
	_ = v452
	var v455 base.V128
	_ = v455
	var v457 base.V128
	_ = v457
	var v458 base.V128
	_ = v458
	var v459 base.V128
	_ = v459
	var v460 int32
	_ = v460
	var v462 int32
	_ = v462
	var v464 int32
	_ = v464
	var v466 int32
	_ = v466
	var v469 base.V128
	_ = v469
	var v472 base.V128
	_ = v472
	var v475 base.V128
	_ = v475
	var v478 base.V128
	_ = v478
	var v480 base.V128
	_ = v480
	var v481 int32
	_ = v481
	var v483 int32
	_ = v483
	var v485 int32
	_ = v485
	var v487 int32
	_ = v487
	var v490 base.V128
	_ = v490
	var v493 base.V128
	_ = v493
	var v496 base.V128
	_ = v496
	var v499 base.V128
	_ = v499
	var v501 base.V128
	_ = v501
	var v502 base.V128
	_ = v502
	var v503 int32
	_ = v503
	var v505 int32
	_ = v505
	var v507 int32
	_ = v507
	var v509 int32
	_ = v509
	var v512 base.V128
	_ = v512
	var v515 base.V128
	_ = v515
	var v518 base.V128
	_ = v518
	var v521 base.V128
	_ = v521
	var v523 base.V128
	_ = v523
	var v524 int32
	_ = v524
	var v526 int32
	_ = v526
	var v528 int32
	_ = v528
	var v530 int32
	_ = v530
	var v533 base.V128
	_ = v533
	var v536 base.V128
	_ = v536
	var v539 base.V128
	_ = v539
	var v542 base.V128
	_ = v542
	var v544 base.V128
	_ = v544
	var v545 base.V128
	_ = v545
	var v546 base.V128
	_ = v546
	var v547 base.V128
	_ = v547
	var v550 int32
	_ = v550
	var v560 base.V128
	_ = v560
	var v563 base.V128
	_ = v563
	var v566 base.V128
	_ = v566
	var v569 base.V128
	_ = v569
	var v571 base.V128
	_ = v571
	var v581 base.V128
	_ = v581
	var v584 base.V128
	_ = v584
	var v587 base.V128
	_ = v587
	var v590 base.V128
	_ = v590
	var v592 base.V128
	_ = v592
	var v593 base.V128
	_ = v593
	var v603 base.V128
	_ = v603
	var v606 base.V128
	_ = v606
	var v609 base.V128
	_ = v609
	var v612 base.V128
	_ = v612
	var v614 base.V128
	_ = v614
	var v624 base.V128
	_ = v624
	var v627 base.V128
	_ = v627
	var v630 base.V128
	_ = v630
	var v633 base.V128
	_ = v633
	var v635 base.V128
	_ = v635
	var v636 base.V128
	_ = v636
	var v637 base.V128
	_ = v637
	var v647 base.V128
	_ = v647
	var v650 base.V128
	_ = v650
	var v653 base.V128
	_ = v653
	var v656 base.V128
	_ = v656
	var v658 base.V128
	_ = v658
	var v668 base.V128
	_ = v668
	var v671 base.V128
	_ = v671
	var v674 base.V128
	_ = v674
	var v677 base.V128
	_ = v677
	var v679 base.V128
	_ = v679
	var v680 base.V128
	_ = v680
	var v690 base.V128
	_ = v690
	var v693 base.V128
	_ = v693
	var v696 base.V128
	_ = v696
	var v699 base.V128
	_ = v699
	var v701 base.V128
	_ = v701
	var v711 base.V128
	_ = v711
	var v714 base.V128
	_ = v714
	var v717 base.V128
	_ = v717
	var v720 base.V128
	_ = v720
	var v722 base.V128
	_ = v722
	var v723 base.V128
	_ = v723
	var v724 base.V128
	_ = v724
	var v725 base.V128
	_ = v725
	var v735 base.V128
	_ = v735
	var v738 base.V128
	_ = v738
	var v741 base.V128
	_ = v741
	var v744 base.V128
	_ = v744
	var v746 base.V128
	_ = v746
	var v756 base.V128
	_ = v756
	var v759 base.V128
	_ = v759
	var v762 base.V128
	_ = v762
	var v765 base.V128
	_ = v765
	var v767 base.V128
	_ = v767
	var v768 base.V128
	_ = v768
	var v778 base.V128
	_ = v778
	var v781 base.V128
	_ = v781
	var v784 base.V128
	_ = v784
	var v787 base.V128
	_ = v787
	var v789 base.V128
	_ = v789
	var v799 base.V128
	_ = v799
	var v802 base.V128
	_ = v802
	var v805 base.V128
	_ = v805
	var v808 base.V128
	_ = v808
	var v810 base.V128
	_ = v810
	var v811 base.V128
	_ = v811
	var v812 base.V128
	_ = v812
	var v822 base.V128
	_ = v822
	var v825 base.V128
	_ = v825
	var v828 base.V128
	_ = v828
	var v831 base.V128
	_ = v831
	var v833 base.V128
	_ = v833
	var v843 base.V128
	_ = v843
	var v846 base.V128
	_ = v846
	var v849 base.V128
	_ = v849
	var v852 base.V128
	_ = v852
	var v854 base.V128
	_ = v854
	var v855 base.V128
	_ = v855
	var v865 base.V128
	_ = v865
	var v868 base.V128
	_ = v868
	var v871 base.V128
	_ = v871
	var v874 base.V128
	_ = v874
	var v876 base.V128
	_ = v876
	var v886 base.V128
	_ = v886
	var v889 base.V128
	_ = v889
	var v892 base.V128
	_ = v892
	var v895 base.V128
	_ = v895
	var v897 base.V128
	_ = v897
	var v898 base.V128
	_ = v898
	var v899 base.V128
	_ = v899
	var v900 base.V128
	_ = v900
	var v904 base.V128
	_ = v904
	var v905 base.V128
	_ = v905
	var v906 base.V128
	_ = v906
	var v907 base.V128
	_ = v907
	var v908 base.V128
	_ = v908
	var v909 base.V128
	_ = v909
	var v910 base.V128
	_ = v910
	var v911 base.V128
	_ = v911
	var v912 base.V128
	_ = v912
	var v913 base.V128
	_ = v913
	var v914 base.V128
	_ = v914
	var v915 base.V128
	_ = v915
	var v916 base.V128
	_ = v916
	var v917 base.V128
	_ = v917
	var v925 base.V128
	_ = v925
	var v926 base.V128
	_ = v926
	var v935 base.V128
	_ = v935
	var v936 base.V128
	_ = v936
	var v949 base.V128
	_ = v949
	var v950 base.V128
	_ = v950
	var v951 base.V128
	_ = v951
	var v952 base.V128
	_ = v952
	var v953 base.V128
	_ = v953
	var v954 base.V128
	_ = v954
	var v959 base.V128
	_ = v959
	var v960 base.V128
	_ = v960
	var v973 base.V128
	_ = v973
	var v974 base.V128
	_ = v974
	var v975 base.V128
	_ = v975
	var v976 base.V128
	_ = v976
	var v977 base.V128
	_ = v977
	var v978 base.V128
	_ = v978
	var v983 base.V128
	_ = v983
	var v984 base.V128
	_ = v984
	var v998 base.V128
	_ = v998
	var v999 base.V128
	_ = v999
	var v1012 base.V128
	_ = v1012
	var v1013 base.V128
	_ = v1013
	var v1014 base.V128
	_ = v1014
	var v1015 base.V128
	_ = v1015
	var v1016 base.V128
	_ = v1016
	var v1017 base.V128
	_ = v1017
	var v1018 base.V128
	_ = v1018
	var v1019 base.V128
	_ = v1019
	var v1020 base.V128
	_ = v1020
	var v1021 base.V128
	_ = v1021
	var v1022 base.V128
	_ = v1022
	var v1023 base.V128
	_ = v1023
	var v1024 base.V128
	_ = v1024
	var v1025 base.V128
	_ = v1025
	var v1030 base.V128
	_ = v1030
	var v1031 base.V128
	_ = v1031
	var v1044 base.V128
	_ = v1044
	var v1045 base.V128
	_ = v1045
	var v1046 base.V128
	_ = v1046
	var v1047 base.V128
	_ = v1047
	var v1048 base.V128
	_ = v1048
	var v1049 base.V128
	_ = v1049
	var v1054 base.V128
	_ = v1054
	var v1055 base.V128
	_ = v1055
	var v1068 base.V128
	_ = v1068
	var v1069 base.V128
	_ = v1069
	var v1070 base.V128
	_ = v1070
	var v1071 base.V128
	_ = v1071
	var v1072 base.V128
	_ = v1072
	var v1073 base.V128
	_ = v1073
	var v1078 base.V128
	_ = v1078
	var v1079 base.V128
	_ = v1079
	var v1098 base.V128
	_ = v1098
	var v1101 base.V128
	_ = v1101
	var v1106 int32
	_ = v1106
	var v1108 int32
	_ = v1108
	v82 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2)+30)))
	v83 = base.Simd_g_i32x4_splat(v82)
	v84 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2)+22)))
	v85 = base.Simd_g_i32x4_splat(v84)
	v86 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2)+14)))
	v87 = base.Simd_g_i32x4_splat(v86)
	v88 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2)+6)))
	v89 = base.Simd_g_i32x4_splat(v88)
	v90 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2)+28)))
	v91 = base.Simd_g_i32x4_splat(v90)
	v92 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2)+20)))
	v93 = base.Simd_g_i32x4_splat(v92)
	v94 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2)+12)))
	v95 = base.Simd_g_i32x4_splat(v94)
	v96 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2)+4)))
	v97 = base.Simd_g_i32x4_splat(v96)
	v98 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2)+26)))
	v99 = base.Simd_g_i32x4_splat(v98)
	v100 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2)+18)))
	v101 = base.Simd_g_i32x4_splat(v100)
	v102 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2)+10)))
	v103 = base.Simd_g_i32x4_splat(v102)
	v104 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2)+2)))
	v105 = base.Simd_g_i32x4_splat(v104)
	v106 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2)+24)))
	v107 = base.Simd_g_i32x4_splat(v106)
	v108 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2)+16)))
	v109 = base.Simd_g_i32x4_splat(v108)
	v110 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2)+8)))
	v111 = base.Simd_g_i32x4_splat(v110)
	v112 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2))))
	v113 = base.Simd_g_i32x4_splat(v112)
	v135 = int32(0)
	v136 = int32(-128)
	for {
		v197 = l1 + v136
		v198 = int32(239)
		v200 = int32(235)
		v202 = int32(231)
		v204 = int32(227)
		v206 = int32(0)
		v207 = base.Simd_g_v128_load8_splat(m, v197+v204, v206)
		v209 = int32(1)
		v210 = base.Simd_g_v128_load8_lane_l1(m, v197+v202, v206, v207)
		v212 = int32(2)
		v213 = base.Simd_g_v128_load8_lane_l2(m, v197+v200, v206, v210)
		v215 = int32(3)
		v216 = base.Simd_g_v128_load8_lane_l3(m, v197+v198, v206, v213)
		v218 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v216))
		v219 = int32(237)
		v221 = int32(233)
		v223 = int32(229)
		v225 = int32(225)
		v228 = base.Simd_g_v128_load8_splat(m, v197+v225, v206)
		v231 = base.Simd_g_v128_load8_lane_l1(m, v197+v223, v206, v228)
		v234 = base.Simd_g_v128_load8_lane_l2(m, v197+v221, v206, v231)
		v237 = base.Simd_g_v128_load8_lane_l3(m, v197+v219, v206, v234)
		v239 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v237))
		v240 = base.Simd_g_i32x4_add(v218, v239)
		v241 = int32(238)
		v243 = int32(234)
		v245 = int32(230)
		v247 = int32(226)
		v250 = base.Simd_g_v128_load8_splat(m, v197+v247, v206)
		v253 = base.Simd_g_v128_load8_lane_l1(m, v197+v245, v206, v250)
		v256 = base.Simd_g_v128_load8_lane_l2(m, v197+v243, v206, v253)
		v259 = base.Simd_g_v128_load8_lane_l3(m, v197+v241, v206, v256)
		v261 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v259))
		v262 = int32(236)
		v264 = int32(232)
		v266 = int32(228)
		v268 = int32(224)
		v271 = base.Simd_g_v128_load8_splat(m, v197+v268, v206)
		v274 = base.Simd_g_v128_load8_lane_l1(m, v197+v266, v206, v271)
		v277 = base.Simd_g_v128_load8_lane_l2(m, v197+v264, v206, v274)
		v280 = base.Simd_g_v128_load8_lane_l3(m, v197+v262, v206, v277)
		v282 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v280))
		v283 = base.Simd_g_i32x4_add(v261, v282)
		v284 = base.Simd_g_i32x4_add(v240, v283)
		v285 = int32(175)
		v287 = int32(171)
		v289 = int32(167)
		v291 = int32(163)
		v294 = base.Simd_g_v128_load8_splat(m, v197+v291, v206)
		v297 = base.Simd_g_v128_load8_lane_l1(m, v197+v289, v206, v294)
		v300 = base.Simd_g_v128_load8_lane_l2(m, v197+v287, v206, v297)
		v303 = base.Simd_g_v128_load8_lane_l3(m, v197+v285, v206, v300)
		v305 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v303))
		v306 = int32(173)
		v308 = int32(169)
		v310 = int32(165)
		v312 = int32(161)
		v315 = base.Simd_g_v128_load8_splat(m, v197+v312, v206)
		v318 = base.Simd_g_v128_load8_lane_l1(m, v197+v310, v206, v315)
		v321 = base.Simd_g_v128_load8_lane_l2(m, v197+v308, v206, v318)
		v324 = base.Simd_g_v128_load8_lane_l3(m, v197+v306, v206, v321)
		v326 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v324))
		v327 = base.Simd_g_i32x4_add(v305, v326)
		v328 = int32(174)
		v330 = int32(170)
		v332 = int32(166)
		v334 = int32(162)
		v337 = base.Simd_g_v128_load8_splat(m, v197+v334, v206)
		v340 = base.Simd_g_v128_load8_lane_l1(m, v197+v332, v206, v337)
		v343 = base.Simd_g_v128_load8_lane_l2(m, v197+v330, v206, v340)
		v346 = base.Simd_g_v128_load8_lane_l3(m, v197+v328, v206, v343)
		v348 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v346))
		v349 = int32(172)
		v351 = int32(168)
		v353 = int32(164)
		v355 = int32(160)
		v358 = base.Simd_g_v128_load8_splat(m, v197+v355, v206)
		v361 = base.Simd_g_v128_load8_lane_l1(m, v197+v353, v206, v358)
		v364 = base.Simd_g_v128_load8_lane_l2(m, v197+v351, v206, v361)
		v367 = base.Simd_g_v128_load8_lane_l3(m, v197+v349, v206, v364)
		v369 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v367))
		v370 = base.Simd_g_i32x4_add(v348, v369)
		v371 = base.Simd_g_i32x4_add(v327, v370)
		v372 = base.Simd_g_i32x4_add(v284, v371)
		v373 = int32(207)
		v375 = int32(203)
		v377 = int32(199)
		v379 = int32(195)
		v382 = base.Simd_g_v128_load8_splat(m, v197+v379, v206)
		v385 = base.Simd_g_v128_load8_lane_l1(m, v197+v377, v206, v382)
		v388 = base.Simd_g_v128_load8_lane_l2(m, v197+v375, v206, v385)
		v391 = base.Simd_g_v128_load8_lane_l3(m, v197+v373, v206, v388)
		v393 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v391))
		v394 = int32(205)
		v396 = int32(201)
		v398 = int32(197)
		v400 = int32(193)
		v403 = base.Simd_g_v128_load8_splat(m, v197+v400, v206)
		v406 = base.Simd_g_v128_load8_lane_l1(m, v197+v398, v206, v403)
		v409 = base.Simd_g_v128_load8_lane_l2(m, v197+v396, v206, v406)
		v412 = base.Simd_g_v128_load8_lane_l3(m, v197+v394, v206, v409)
		v414 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v412))
		v415 = base.Simd_g_i32x4_add(v393, v414)
		v416 = int32(206)
		v418 = int32(202)
		v420 = int32(198)
		v422 = int32(194)
		v425 = base.Simd_g_v128_load8_splat(m, v197+v422, v206)
		v428 = base.Simd_g_v128_load8_lane_l1(m, v197+v420, v206, v425)
		v431 = base.Simd_g_v128_load8_lane_l2(m, v197+v418, v206, v428)
		v434 = base.Simd_g_v128_load8_lane_l3(m, v197+v416, v206, v431)
		v436 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v434))
		v437 = int32(204)
		v439 = int32(200)
		v441 = int32(196)
		v443 = int32(192)
		v446 = base.Simd_g_v128_load8_splat(m, v197+v443, v206)
		v449 = base.Simd_g_v128_load8_lane_l1(m, v197+v441, v206, v446)
		v452 = base.Simd_g_v128_load8_lane_l2(m, v197+v439, v206, v449)
		v455 = base.Simd_g_v128_load8_lane_l3(m, v197+v437, v206, v452)
		v457 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v455))
		v458 = base.Simd_g_i32x4_add(v436, v457)
		v459 = base.Simd_g_i32x4_add(v415, v458)
		v460 = int32(143)
		v462 = int32(139)
		v464 = int32(135)
		v466 = int32(131)
		v469 = base.Simd_g_v128_load8_splat(m, v197+v466, v206)
		v472 = base.Simd_g_v128_load8_lane_l1(m, v197+v464, v206, v469)
		v475 = base.Simd_g_v128_load8_lane_l2(m, v197+v462, v206, v472)
		v478 = base.Simd_g_v128_load8_lane_l3(m, v197+v460, v206, v475)
		v480 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v478))
		v481 = int32(141)
		v483 = int32(137)
		v485 = int32(133)
		v487 = int32(129)
		v490 = base.Simd_g_v128_load8_splat(m, v197+v487, v206)
		v493 = base.Simd_g_v128_load8_lane_l1(m, v197+v485, v206, v490)
		v496 = base.Simd_g_v128_load8_lane_l2(m, v197+v483, v206, v493)
		v499 = base.Simd_g_v128_load8_lane_l3(m, v197+v481, v206, v496)
		v501 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v499))
		v502 = base.Simd_g_i32x4_add(v480, v501)
		v503 = int32(142)
		v505 = int32(138)
		v507 = int32(134)
		v509 = int32(130)
		v512 = base.Simd_g_v128_load8_splat(m, v197+v509, v206)
		v515 = base.Simd_g_v128_load8_lane_l1(m, v197+v507, v206, v512)
		v518 = base.Simd_g_v128_load8_lane_l2(m, v197+v505, v206, v515)
		v521 = base.Simd_g_v128_load8_lane_l3(m, v197+v503, v206, v518)
		v523 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v521))
		v524 = int32(140)
		v526 = int32(136)
		v528 = int32(132)
		v530 = int32(128)
		v533 = base.Simd_g_v128_load8_splat(m, v197+v530, v206)
		v536 = base.Simd_g_v128_load8_lane_l1(m, v197+v528, v206, v533)
		v539 = base.Simd_g_v128_load8_lane_l2(m, v197+v526, v206, v536)
		v542 = base.Simd_g_v128_load8_lane_l3(m, v197+v524, v206, v539)
		v544 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v542))
		v545 = base.Simd_g_i32x4_add(v523, v544)
		v546 = base.Simd_g_i32x4_add(v502, v545)
		v547 = base.Simd_g_i32x4_add(v459, v546)
		v550 = l0 + v136
		v560 = base.Simd_g_v128_load8_splat(m, v550+v225, v206)
		v563 = base.Simd_g_v128_load8_lane_l1(m, v550+v223, v206, v560)
		v566 = base.Simd_g_v128_load8_lane_l2(m, v550+v221, v206, v563)
		v569 = base.Simd_g_v128_load8_lane_l3(m, v550+v219, v206, v566)
		v571 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v569))
		v581 = base.Simd_g_v128_load8_splat(m, v550+v204, v206)
		v584 = base.Simd_g_v128_load8_lane_l1(m, v550+v202, v206, v581)
		v587 = base.Simd_g_v128_load8_lane_l2(m, v550+v200, v206, v584)
		v590 = base.Simd_g_v128_load8_lane_l3(m, v550+v198, v206, v587)
		v592 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v590))
		v593 = base.Simd_g_i32x4_sub(v571, v592)
		v603 = base.Simd_g_v128_load8_splat(m, v550+v268, v206)
		v606 = base.Simd_g_v128_load8_lane_l1(m, v550+v266, v206, v603)
		v609 = base.Simd_g_v128_load8_lane_l2(m, v550+v264, v206, v606)
		v612 = base.Simd_g_v128_load8_lane_l3(m, v550+v262, v206, v609)
		v614 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v612))
		v624 = base.Simd_g_v128_load8_splat(m, v550+v247, v206)
		v627 = base.Simd_g_v128_load8_lane_l1(m, v550+v245, v206, v624)
		v630 = base.Simd_g_v128_load8_lane_l2(m, v550+v243, v206, v627)
		v633 = base.Simd_g_v128_load8_lane_l3(m, v550+v241, v206, v630)
		v635 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v633))
		v636 = base.Simd_g_i32x4_sub(v614, v635)
		v637 = base.Simd_g_i32x4_add(v593, v636)
		v647 = base.Simd_g_v128_load8_splat(m, v550+v312, v206)
		v650 = base.Simd_g_v128_load8_lane_l1(m, v550+v310, v206, v647)
		v653 = base.Simd_g_v128_load8_lane_l2(m, v550+v308, v206, v650)
		v656 = base.Simd_g_v128_load8_lane_l3(m, v550+v306, v206, v653)
		v658 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v656))
		v668 = base.Simd_g_v128_load8_splat(m, v550+v291, v206)
		v671 = base.Simd_g_v128_load8_lane_l1(m, v550+v289, v206, v668)
		v674 = base.Simd_g_v128_load8_lane_l2(m, v550+v287, v206, v671)
		v677 = base.Simd_g_v128_load8_lane_l3(m, v550+v285, v206, v674)
		v679 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v677))
		v680 = base.Simd_g_i32x4_sub(v658, v679)
		v690 = base.Simd_g_v128_load8_splat(m, v550+v355, v206)
		v693 = base.Simd_g_v128_load8_lane_l1(m, v550+v353, v206, v690)
		v696 = base.Simd_g_v128_load8_lane_l2(m, v550+v351, v206, v693)
		v699 = base.Simd_g_v128_load8_lane_l3(m, v550+v349, v206, v696)
		v701 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v699))
		v711 = base.Simd_g_v128_load8_splat(m, v550+v334, v206)
		v714 = base.Simd_g_v128_load8_lane_l1(m, v550+v332, v206, v711)
		v717 = base.Simd_g_v128_load8_lane_l2(m, v550+v330, v206, v714)
		v720 = base.Simd_g_v128_load8_lane_l3(m, v550+v328, v206, v717)
		v722 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v720))
		v723 = base.Simd_g_i32x4_sub(v701, v722)
		v724 = base.Simd_g_i32x4_add(v680, v723)
		v725 = base.Simd_g_i32x4_add(v637, v724)
		v735 = base.Simd_g_v128_load8_splat(m, v550+v400, v206)
		v738 = base.Simd_g_v128_load8_lane_l1(m, v550+v398, v206, v735)
		v741 = base.Simd_g_v128_load8_lane_l2(m, v550+v396, v206, v738)
		v744 = base.Simd_g_v128_load8_lane_l3(m, v550+v394, v206, v741)
		v746 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v744))
		v756 = base.Simd_g_v128_load8_splat(m, v550+v379, v206)
		v759 = base.Simd_g_v128_load8_lane_l1(m, v550+v377, v206, v756)
		v762 = base.Simd_g_v128_load8_lane_l2(m, v550+v375, v206, v759)
		v765 = base.Simd_g_v128_load8_lane_l3(m, v550+v373, v206, v762)
		v767 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v765))
		v768 = base.Simd_g_i32x4_sub(v746, v767)
		v778 = base.Simd_g_v128_load8_splat(m, v550+v443, v206)
		v781 = base.Simd_g_v128_load8_lane_l1(m, v550+v441, v206, v778)
		v784 = base.Simd_g_v128_load8_lane_l2(m, v550+v439, v206, v781)
		v787 = base.Simd_g_v128_load8_lane_l3(m, v550+v437, v206, v784)
		v789 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v787))
		v799 = base.Simd_g_v128_load8_splat(m, v550+v422, v206)
		v802 = base.Simd_g_v128_load8_lane_l1(m, v550+v420, v206, v799)
		v805 = base.Simd_g_v128_load8_lane_l2(m, v550+v418, v206, v802)
		v808 = base.Simd_g_v128_load8_lane_l3(m, v550+v416, v206, v805)
		v810 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v808))
		v811 = base.Simd_g_i32x4_sub(v789, v810)
		v812 = base.Simd_g_i32x4_add(v768, v811)
		v822 = base.Simd_g_v128_load8_splat(m, v550+v487, v206)
		v825 = base.Simd_g_v128_load8_lane_l1(m, v550+v485, v206, v822)
		v828 = base.Simd_g_v128_load8_lane_l2(m, v550+v483, v206, v825)
		v831 = base.Simd_g_v128_load8_lane_l3(m, v550+v481, v206, v828)
		v833 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v831))
		v843 = base.Simd_g_v128_load8_splat(m, v550+v466, v206)
		v846 = base.Simd_g_v128_load8_lane_l1(m, v550+v464, v206, v843)
		v849 = base.Simd_g_v128_load8_lane_l2(m, v550+v462, v206, v846)
		v852 = base.Simd_g_v128_load8_lane_l3(m, v550+v460, v206, v849)
		v854 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v852))
		v855 = base.Simd_g_i32x4_sub(v833, v854)
		v865 = base.Simd_g_v128_load8_splat(m, v550+v530, v206)
		v868 = base.Simd_g_v128_load8_lane_l1(m, v550+v528, v206, v865)
		v871 = base.Simd_g_v128_load8_lane_l2(m, v550+v526, v206, v868)
		v874 = base.Simd_g_v128_load8_lane_l3(m, v550+v524, v206, v871)
		v876 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v874))
		v886 = base.Simd_g_v128_load8_splat(m, v550+v509, v206)
		v889 = base.Simd_g_v128_load8_lane_l1(m, v550+v507, v206, v886)
		v892 = base.Simd_g_v128_load8_lane_l2(m, v550+v505, v206, v889)
		v895 = base.Simd_g_v128_load8_lane_l3(m, v550+v503, v206, v892)
		v897 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v895))
		v898 = base.Simd_g_i32x4_sub(v876, v897)
		v899 = base.Simd_g_i32x4_add(v855, v898)
		v900 = base.Simd_g_i32x4_add(v812, v899)
		v904 = base.Simd_g_i32x4_add(v592, v571)
		v905 = base.Simd_g_i32x4_add(v635, v614)
		v906 = base.Simd_g_i32x4_add(v904, v905)
		v907 = base.Simd_g_i32x4_add(v679, v658)
		v908 = base.Simd_g_i32x4_add(v722, v701)
		v909 = base.Simd_g_i32x4_add(v907, v908)
		v910 = base.Simd_g_i32x4_add(v906, v909)
		v911 = base.Simd_g_i32x4_add(v767, v746)
		v912 = base.Simd_g_i32x4_add(v810, v789)
		v913 = base.Simd_g_i32x4_add(v911, v912)
		v914 = base.Simd_g_i32x4_add(v854, v833)
		v915 = base.Simd_g_i32x4_add(v897, v876)
		v916 = base.Simd_g_i32x4_add(v914, v915)
		v917 = base.Simd_g_i32x4_add(v913, v916)
		v925 = base.Simd_g_i32x4_sub(v916, v913)
		v926 = base.Simd_g_i32x4_sub(v909, v906)
		v935 = base.Simd_g_i32x4_sub(v724, v637)
		v936 = base.Simd_g_i32x4_sub(v899, v812)
		v949 = base.Simd_g_i32x4_sub(v636, v593)
		v950 = base.Simd_g_i32x4_sub(v723, v680)
		v951 = base.Simd_g_i32x4_add(v949, v950)
		v952 = base.Simd_g_i32x4_sub(v811, v768)
		v953 = base.Simd_g_i32x4_sub(v898, v855)
		v954 = base.Simd_g_i32x4_add(v952, v953)
		v959 = base.Simd_g_i32x4_sub(v950, v949)
		v960 = base.Simd_g_i32x4_sub(v953, v952)
		v973 = base.Simd_g_i32x4_sub(v905, v904)
		v974 = base.Simd_g_i32x4_sub(v908, v907)
		v975 = base.Simd_g_i32x4_add(v973, v974)
		v976 = base.Simd_g_i32x4_sub(v912, v911)
		v977 = base.Simd_g_i32x4_sub(v915, v914)
		v978 = base.Simd_g_i32x4_add(v976, v977)
		v983 = base.Simd_g_i32x4_sub(v974, v973)
		v984 = base.Simd_g_i32x4_sub(v977, v976)
		v998 = base.Simd_g_i32x4_sub(v371, v284)
		v999 = base.Simd_g_i32x4_sub(v546, v459)
		v1012 = base.Simd_g_i32x4_sub(v239, v218)
		v1013 = base.Simd_g_i32x4_sub(v282, v261)
		v1014 = base.Simd_g_i32x4_add(v1012, v1013)
		v1015 = base.Simd_g_i32x4_sub(v326, v305)
		v1016 = base.Simd_g_i32x4_sub(v369, v348)
		v1017 = base.Simd_g_i32x4_add(v1015, v1016)
		v1018 = base.Simd_g_i32x4_add(v1014, v1017)
		v1019 = base.Simd_g_i32x4_sub(v414, v393)
		v1020 = base.Simd_g_i32x4_sub(v457, v436)
		v1021 = base.Simd_g_i32x4_add(v1019, v1020)
		v1022 = base.Simd_g_i32x4_sub(v501, v480)
		v1023 = base.Simd_g_i32x4_sub(v544, v523)
		v1024 = base.Simd_g_i32x4_add(v1022, v1023)
		v1025 = base.Simd_g_i32x4_add(v1021, v1024)
		v1030 = base.Simd_g_i32x4_sub(v1017, v1014)
		v1031 = base.Simd_g_i32x4_sub(v1024, v1021)
		v1044 = base.Simd_g_i32x4_sub(v1013, v1012)
		v1045 = base.Simd_g_i32x4_sub(v1016, v1015)
		v1046 = base.Simd_g_i32x4_add(v1044, v1045)
		v1047 = base.Simd_g_i32x4_sub(v1020, v1019)
		v1048 = base.Simd_g_i32x4_sub(v1023, v1022)
		v1049 = base.Simd_g_i32x4_add(v1047, v1048)
		v1054 = base.Simd_g_i32x4_sub(v1045, v1044)
		v1055 = base.Simd_g_i32x4_sub(v1048, v1047)
		v1068 = base.Simd_g_i32x4_sub(v283, v240)
		v1069 = base.Simd_g_i32x4_sub(v370, v327)
		v1070 = base.Simd_g_i32x4_add(v1068, v1069)
		v1071 = base.Simd_g_i32x4_sub(v458, v415)
		v1072 = base.Simd_g_i32x4_sub(v545, v502)
		v1073 = base.Simd_g_i32x4_add(v1071, v1072)
		v1078 = base.Simd_g_i32x4_sub(v1069, v1068)
		v1079 = base.Simd_g_i32x4_sub(v1072, v1071)
		v1098 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_add(v372, v547), v113), base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_add(v725, v900)), v105), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_add(v910, v917), v113)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_sub(v917, v910)), v107)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_sub(v925, v926)), v109)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_add(v926, v925)), v111)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_add(v935, v936)), v103)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_sub(v936, v935)), v101)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_sub(v900, v725)), v99)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_add(v951, v954)), v97)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_add(v959, v960)), v95)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_sub(v960, v959)), v93)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_sub(v954, v951)), v91)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_add(v975, v978)), v89)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_add(v983, v984)), v87)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_sub(v984, v983)), v85)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_sub(v978, v975)), v83))), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_add(v998, v999)), v111)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_sub(v999, v998)), v109)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_sub(v547, v372)), v107)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_add(v1018, v1025)), v105)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_add(v1030, v1031)), v103)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_sub(v1031, v1030)), v101)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_sub(v1025, v1018)), v99)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_add(v1046, v1049)), v97)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_add(v1054, v1055)), v95)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_sub(v1055, v1054)), v93)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_sub(v1049, v1046)), v91)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_add(v1070, v1073)), v89)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_add(v1078, v1079)), v87)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_sub(v1079, v1078)), v85)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_abs(base.Simd_g_i32x4_sub(v1073, v1070)), v83))), int32(5)), base.Simd_g_i32x4_replace_lane_l0(base.Simd_g_const(&F_Disto16x16_C__k0), v135))
		v1101 = base.Simd_g_i32x4_add(v1098, base.Simd_g_i8x16_shuffle2(v1098, v83, base.Simd_g_const(&F_Disto16x16_C__k1), base.Simd_g_const(&F_Disto16x16_C__k2)))
		v1106 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v1101, base.Simd_g_i8x16_shuffle2(v1101, v83, base.Simd_g_const(&F_Disto16x16_C__k3), base.Simd_g_const(&F_Disto16x16_C__k2))))
		v1108 = v136 + v530
		if base.Ui32(v1108) < base.Ui32(int32(384)) {
			v135 = v1106
			v136 = v1108
			continue
		} else {
			break
		}
		break
	}
	return v1106
}

var F_Disto16x16_C__k0 = [2]uint64{0x0, 0x0}
var F_Disto16x16_C__k1 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_Disto16x16_C__k2 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_Disto16x16_C__k3 = [2]uint64{0x302010007060504, 0x302010003020100}

func F_Disto16x16_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	var v6 int32
	_ = v6
	var v20 base.V128
	_ = v20
	var v22 base.V128
	_ = v22
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v44 int32
	_ = v44
	var v48 int32
	_ = v48
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	var v62 base.V128
	_ = v62
	var v63 int32
	_ = v63
	var v67 base.V128
	_ = v67
	var v68 base.V128
	_ = v68
	var v70 base.V128
	_ = v70
	var v71 base.V128
	_ = v71
	var v72 base.V128
	_ = v72
	var v73 int32
	_ = v73
	var v76 base.V128
	_ = v76
	var v80 base.V128
	_ = v80
	var v84 base.V128
	_ = v84
	var v85 base.V128
	_ = v85
	var v86 int32
	_ = v86
	var v89 base.V128
	_ = v89
	var v93 base.V128
	_ = v93
	var v97 base.V128
	_ = v97
	var v98 int32
	_ = v98
	var v101 base.V128
	_ = v101
	var v105 base.V128
	_ = v105
	var v109 base.V128
	_ = v109
	var v110 base.V128
	_ = v110
	var v111 base.V128
	_ = v111
	var v112 base.V128
	_ = v112
	var v113 base.V128
	_ = v113
	var v114 base.V128
	_ = v114
	var v115 base.V128
	_ = v115
	var v116 base.V128
	_ = v116
	var v117 base.V128
	_ = v117
	var v118 base.V128
	_ = v118
	var v120 base.V128
	_ = v120
	var v122 base.V128
	_ = v122
	var v123 base.V128
	_ = v123
	var v124 base.V128
	_ = v124
	var v126 base.V128
	_ = v126
	var v128 base.V128
	_ = v128
	var v129 base.V128
	_ = v129
	var v130 base.V128
	_ = v130
	var v131 base.V128
	_ = v131
	var v132 base.V128
	_ = v132
	var v134 base.V128
	_ = v134
	var v136 base.V128
	_ = v136
	var v137 base.V128
	_ = v137
	var v138 base.V128
	_ = v138
	var v139 base.V128
	_ = v139
	var v141 base.V128
	_ = v141
	var v142 base.V128
	_ = v142
	var v143 base.V128
	_ = v143
	var v144 base.V128
	_ = v144
	var v145 base.V128
	_ = v145
	var v146 base.V128
	_ = v146
	var v151 base.V128
	_ = v151
	var v152 base.V128
	_ = v152
	var v167 base.V128
	_ = v167
	var v178 int32
	_ = v178
	var v180 int32
	_ = v180
	var v185 int32
	_ = v185
	var v187 int32
	_ = v187
	var v190 int32
	_ = v190
	v6 = int32(0)
	v20 = base.Simd_g_v128_load_rng(m, l2+int32(16), v6, int32(-16), int32(32))
	v22 = base.Simd_g_v128_load_nc(m, l2, v6)
	v25 = l0
	v26 = l1
	v30 = v6
	v31 = v6
	for {
		v44 = int32(-4)
		v48 = v31
		for {
			v58 = v25 + v44
			v59 = int32(100)
			v61 = int32(0)
			v62 = base.Simd_g_v128_load64_zero(m, v58+v59, v61)
			v63 = v26 + v44
			v67 = base.Simd_g_v128_load64_zero(m, v63+v59, v61)
			v68 = base.Simd_g_const(&F_Disto16x16_SSE2__k0)
			v70 = base.Simd_g_const(&F_Disto16x16_SSE2__k1)
			v71 = base.Simd_g_const(&F_Disto16x16_SSE2__k2)
			v72 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v62, v67, base.Simd_g_const(&F_Disto16x16_SSE2__k3), base.Simd_g_const(&F_Disto16x16_SSE2__k4)), v70, base.Simd_g_const(&F_Disto16x16_SSE2__k5), base.Simd_g_const(&F_Disto16x16_SSE2__k6))
			v73 = int32(36)
			v76 = base.Simd_g_v128_load64_zero(m, v58+v73, v61)
			v80 = base.Simd_g_v128_load64_zero(m, v63+v73, v61)
			v84 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v76, v80, base.Simd_g_const(&F_Disto16x16_SSE2__k3), base.Simd_g_const(&F_Disto16x16_SSE2__k4)), v70, base.Simd_g_const(&F_Disto16x16_SSE2__k5), base.Simd_g_const(&F_Disto16x16_SSE2__k6))
			v85 = base.Simd_g_i16x8_add(v72, v84)
			v86 = int32(68)
			v89 = base.Simd_g_v128_load64_zero(m, v58+v86, v61)
			v93 = base.Simd_g_v128_load64_zero(m, v63+v86, v61)
			v97 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v89, v93, base.Simd_g_const(&F_Disto16x16_SSE2__k3), base.Simd_g_const(&F_Disto16x16_SSE2__k4)), v70, base.Simd_g_const(&F_Disto16x16_SSE2__k5), base.Simd_g_const(&F_Disto16x16_SSE2__k6))
			v98 = int32(4)
			v101 = base.Simd_g_v128_load64_zero(m, v58+v98, v61)
			v105 = base.Simd_g_v128_load64_zero(m, v63+v98, v61)
			v109 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v101, v105, base.Simd_g_const(&F_Disto16x16_SSE2__k3), base.Simd_g_const(&F_Disto16x16_SSE2__k4)), v70, base.Simd_g_const(&F_Disto16x16_SSE2__k5), base.Simd_g_const(&F_Disto16x16_SSE2__k6))
			v110 = base.Simd_g_i16x8_add(v97, v109)
			v111 = base.Simd_g_i16x8_add(v85, v110)
			v112 = base.Simd_g_i16x8_sub(v84, v72)
			v113 = base.Simd_g_i16x8_sub(v109, v97)
			v114 = base.Simd_g_i16x8_add(v112, v113)
			v115 = base.Simd_g_const(&F_Disto16x16_SSE2__k7)
			v116 = base.Simd_g_i8x16_shuffle2(v111, v114, base.Simd_g_const(&F_Disto16x16_SSE2__k8), base.Simd_g_const(&F_Disto16x16_SSE2__k9))
			v117 = base.Simd_g_i16x8_sub(v113, v112)
			v118 = base.Simd_g_i16x8_sub(v110, v85)
			v120 = base.Simd_g_i8x16_shuffle2(v117, v118, base.Simd_g_const(&F_Disto16x16_SSE2__k8), base.Simd_g_const(&F_Disto16x16_SSE2__k9))
			v122 = base.Simd_g_i8x16_shuffle2(v116, v120, base.Simd_g_const(&F_Disto16x16_SSE2__k3), base.Simd_g_const(&F_Disto16x16_SSE2__k4))
			v123 = base.Simd_g_const(&F_Disto16x16_SSE2__k10)
			v124 = base.Simd_g_i8x16_shuffle2(v111, v114, base.Simd_g_const(&F_Disto16x16_SSE2__k11), base.Simd_g_const(&F_Disto16x16_SSE2__k12))
			v126 = base.Simd_g_i8x16_shuffle2(v117, v118, base.Simd_g_const(&F_Disto16x16_SSE2__k11), base.Simd_g_const(&F_Disto16x16_SSE2__k12))
			v128 = base.Simd_g_i8x16_shuffle2(v124, v126, base.Simd_g_const(&F_Disto16x16_SSE2__k3), base.Simd_g_const(&F_Disto16x16_SSE2__k4))
			v129 = base.Simd_g_const(&F_Disto16x16_SSE2__k13)
			v130 = base.Simd_g_i8x16_shuffle2(v122, v128, base.Simd_g_const(&F_Disto16x16_SSE2__k14), base.Simd_g_const(&F_Disto16x16_SSE2__k15))
			v131 = base.Simd_g_const(&F_Disto16x16_SSE2__k16)
			v132 = base.Simd_g_i8x16_shuffle2(v116, v120, base.Simd_g_const(&F_Disto16x16_SSE2__k17), base.Simd_g_const(&F_Disto16x16_SSE2__k18))
			v134 = base.Simd_g_i8x16_shuffle2(v124, v126, base.Simd_g_const(&F_Disto16x16_SSE2__k17), base.Simd_g_const(&F_Disto16x16_SSE2__k18))
			v136 = base.Simd_g_i8x16_shuffle2(v132, v134, base.Simd_g_const(&F_Disto16x16_SSE2__k14), base.Simd_g_const(&F_Disto16x16_SSE2__k15))
			v137 = base.Simd_g_i16x8_add(v130, v136)
			v138 = base.Simd_g_const(&F_Disto16x16_SSE2__k19)
			v139 = base.Simd_g_i8x16_shuffle2(v122, v128, base.Simd_g_const(&F_Disto16x16_SSE2__k20), base.Simd_g_const(&F_Disto16x16_SSE2__k21))
			v141 = base.Simd_g_i8x16_shuffle2(v132, v134, base.Simd_g_const(&F_Disto16x16_SSE2__k20), base.Simd_g_const(&F_Disto16x16_SSE2__k21))
			v142 = base.Simd_g_i16x8_add(v139, v141)
			v143 = base.Simd_g_i16x8_add(v137, v142)
			v144 = base.Simd_g_i16x8_sub(v130, v136)
			v145 = base.Simd_g_i16x8_sub(v139, v141)
			v146 = base.Simd_g_i16x8_add(v144, v145)
			v151 = base.Simd_g_i16x8_sub(v144, v145)
			v152 = base.Simd_g_i16x8_sub(v137, v142)
			v167 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(v143, v146, base.Simd_g_const(&F_Disto16x16_SSE2__k14), base.Simd_g_const(&F_Disto16x16_SSE2__k15))), v22), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(v151, v152, base.Simd_g_const(&F_Disto16x16_SSE2__k14), base.Simd_g_const(&F_Disto16x16_SSE2__k15))), v20)), base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(v151, v152, base.Simd_g_const(&F_Disto16x16_SSE2__k20), base.Simd_g_const(&F_Disto16x16_SSE2__k21))), v20), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(v143, v146, base.Simd_g_const(&F_Disto16x16_SSE2__k20), base.Simd_g_const(&F_Disto16x16_SSE2__k21))), v22)))
			v178 = base.Simd_g_i32x4_extract_lane_l0(v167) + base.Simd_g_i32x4_extract_lane_l1(v167) + base.Simd_g_i32x4_extract_lane_l2(v167) + base.Simd_g_i32x4_extract_lane_l3(v167)
			v180 = v178 >> (uint(int32(31)) % 32)
			v185 = int32(base.Ui32(v178^v180-v180)>>(uint(int32(5))%32)) + v48
			v187 = v44 + v98
			if base.Ui32(v187) < base.Ui32(int32(12)) {
				v44 = v187
				v48 = v185
				continue
			} else {
				break
			}
			break
		}
		v190 = int32(128)
		if base.Ui32(v30) < base.Ui32(int32(384)) {
			v25 = v25 + v190
			v26 = v26 + v190
			v30 = v30 + v190
			v31 = v185
			continue
		} else {
			break
		}
		break
	}
	return v185
}

var F_Disto16x16_SSE2__k0 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_Disto16x16_SSE2__k1 = [2]uint64{0x0, 0x0}
var F_Disto16x16_SSE2__k2 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_Disto16x16_SSE2__k3 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_Disto16x16_SSE2__k4 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_Disto16x16_SSE2__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_Disto16x16_SSE2__k6 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_Disto16x16_SSE2__k7 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_Disto16x16_SSE2__k8 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_Disto16x16_SSE2__k9 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_Disto16x16_SSE2__k10 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_Disto16x16_SSE2__k11 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_Disto16x16_SSE2__k12 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_Disto16x16_SSE2__k13 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_Disto16x16_SSE2__k14 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_Disto16x16_SSE2__k15 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_Disto16x16_SSE2__k16 = [2]uint64{0x1b1a19180b0a0908, 0x1f1e1d1c0f0e0d0c}
var F_Disto16x16_SSE2__k17 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_Disto16x16_SSE2__k18 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_Disto16x16_SSE2__k19 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_Disto16x16_SSE2__k20 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_Disto16x16_SSE2__k21 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}

func F_Disto16x16_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	var v6 int32
	_ = v6
	var v20 base.V128
	_ = v20
	var v22 base.V128
	_ = v22
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v44 int32
	_ = v44
	var v48 int32
	_ = v48
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	var v62 base.V128
	_ = v62
	var v63 int32
	_ = v63
	var v67 base.V128
	_ = v67
	var v68 base.V128
	_ = v68
	var v70 base.V128
	_ = v70
	var v71 int32
	_ = v71
	var v74 base.V128
	_ = v74
	var v78 base.V128
	_ = v78
	var v81 base.V128
	_ = v81
	var v82 base.V128
	_ = v82
	var v83 int32
	_ = v83
	var v86 base.V128
	_ = v86
	var v90 base.V128
	_ = v90
	var v93 base.V128
	_ = v93
	var v94 int32
	_ = v94
	var v97 base.V128
	_ = v97
	var v101 base.V128
	_ = v101
	var v104 base.V128
	_ = v104
	var v105 base.V128
	_ = v105
	var v106 base.V128
	_ = v106
	var v107 base.V128
	_ = v107
	var v108 base.V128
	_ = v108
	var v109 base.V128
	_ = v109
	var v110 base.V128
	_ = v110
	var v111 base.V128
	_ = v111
	var v112 base.V128
	_ = v112
	var v113 base.V128
	_ = v113
	var v115 base.V128
	_ = v115
	var v117 base.V128
	_ = v117
	var v118 base.V128
	_ = v118
	var v119 base.V128
	_ = v119
	var v121 base.V128
	_ = v121
	var v123 base.V128
	_ = v123
	var v124 base.V128
	_ = v124
	var v125 base.V128
	_ = v125
	var v126 base.V128
	_ = v126
	var v127 base.V128
	_ = v127
	var v129 base.V128
	_ = v129
	var v131 base.V128
	_ = v131
	var v132 base.V128
	_ = v132
	var v133 base.V128
	_ = v133
	var v134 base.V128
	_ = v134
	var v136 base.V128
	_ = v136
	var v137 base.V128
	_ = v137
	var v138 base.V128
	_ = v138
	var v139 base.V128
	_ = v139
	var v140 base.V128
	_ = v140
	var v141 base.V128
	_ = v141
	var v146 base.V128
	_ = v146
	var v147 base.V128
	_ = v147
	var v162 base.V128
	_ = v162
	var v173 int32
	_ = v173
	var v175 int32
	_ = v175
	var v180 int32
	_ = v180
	var v182 int32
	_ = v182
	var v185 int32
	_ = v185
	v6 = int32(0)
	v20 = base.Simd_g_v128_load_rng(m, l2+int32(16), v6, int32(-16), int32(32))
	v22 = base.Simd_g_v128_load_nc(m, l2, v6)
	v25 = l0
	v26 = l1
	v30 = v6
	v31 = v6
	for {
		v44 = int32(-4)
		v48 = v31
		for {
			v58 = v25 + v44
			v59 = int32(100)
			v61 = int32(0)
			v62 = base.Simd_g_v128_load64_zero(m, v58+v59, v61)
			v63 = v26 + v44
			v67 = base.Simd_g_v128_load64_zero(m, v63+v59, v61)
			v68 = base.Simd_g_const(&F_Disto16x16_SSE41__k0)
			v70 = base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v62, v67, base.Simd_g_const(&F_Disto16x16_SSE41__k1), base.Simd_g_const(&F_Disto16x16_SSE41__k2)))
			v71 = int32(36)
			v74 = base.Simd_g_v128_load_rng(m, v58+v71, v61, int32(-32), int32(80))
			v78 = base.Simd_g_v128_load_rng(m, v63+v71, v61, int32(-32), int32(80))
			v81 = base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v74, v78, base.Simd_g_const(&F_Disto16x16_SSE41__k1), base.Simd_g_const(&F_Disto16x16_SSE41__k2)))
			v82 = base.Simd_g_i16x8_add(v70, v81)
			v83 = int32(68)
			v86 = base.Simd_g_v128_load_nc(m, v58+v83, v61)
			v90 = base.Simd_g_v128_load_nc(m, v63+v83, v61)
			v93 = base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v86, v90, base.Simd_g_const(&F_Disto16x16_SSE41__k1), base.Simd_g_const(&F_Disto16x16_SSE41__k2)))
			v94 = int32(4)
			v97 = base.Simd_g_v128_load_nc(m, v58+v94, v61)
			v101 = base.Simd_g_v128_load_nc(m, v63+v94, v61)
			v104 = base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v97, v101, base.Simd_g_const(&F_Disto16x16_SSE41__k1), base.Simd_g_const(&F_Disto16x16_SSE41__k2)))
			v105 = base.Simd_g_i16x8_add(v93, v104)
			v106 = base.Simd_g_i16x8_add(v82, v105)
			v107 = base.Simd_g_i16x8_sub(v81, v70)
			v108 = base.Simd_g_i16x8_sub(v104, v93)
			v109 = base.Simd_g_i16x8_add(v107, v108)
			v110 = base.Simd_g_const(&F_Disto16x16_SSE41__k3)
			v111 = base.Simd_g_i8x16_shuffle2(v106, v109, base.Simd_g_const(&F_Disto16x16_SSE41__k4), base.Simd_g_const(&F_Disto16x16_SSE41__k5))
			v112 = base.Simd_g_i16x8_sub(v108, v107)
			v113 = base.Simd_g_i16x8_sub(v105, v82)
			v115 = base.Simd_g_i8x16_shuffle2(v112, v113, base.Simd_g_const(&F_Disto16x16_SSE41__k4), base.Simd_g_const(&F_Disto16x16_SSE41__k5))
			v117 = base.Simd_g_i8x16_shuffle2(v111, v115, base.Simd_g_const(&F_Disto16x16_SSE41__k1), base.Simd_g_const(&F_Disto16x16_SSE41__k2))
			v118 = base.Simd_g_const(&F_Disto16x16_SSE41__k6)
			v119 = base.Simd_g_i8x16_shuffle2(v106, v109, base.Simd_g_const(&F_Disto16x16_SSE41__k7), base.Simd_g_const(&F_Disto16x16_SSE41__k8))
			v121 = base.Simd_g_i8x16_shuffle2(v112, v113, base.Simd_g_const(&F_Disto16x16_SSE41__k7), base.Simd_g_const(&F_Disto16x16_SSE41__k8))
			v123 = base.Simd_g_i8x16_shuffle2(v119, v121, base.Simd_g_const(&F_Disto16x16_SSE41__k1), base.Simd_g_const(&F_Disto16x16_SSE41__k2))
			v124 = base.Simd_g_const(&F_Disto16x16_SSE41__k9)
			v125 = base.Simd_g_i8x16_shuffle2(v117, v123, base.Simd_g_const(&F_Disto16x16_SSE41__k10), base.Simd_g_const(&F_Disto16x16_SSE41__k11))
			v126 = base.Simd_g_const(&F_Disto16x16_SSE41__k12)
			v127 = base.Simd_g_i8x16_shuffle2(v111, v115, base.Simd_g_const(&F_Disto16x16_SSE41__k13), base.Simd_g_const(&F_Disto16x16_SSE41__k14))
			v129 = base.Simd_g_i8x16_shuffle2(v119, v121, base.Simd_g_const(&F_Disto16x16_SSE41__k13), base.Simd_g_const(&F_Disto16x16_SSE41__k14))
			v131 = base.Simd_g_i8x16_shuffle2(v127, v129, base.Simd_g_const(&F_Disto16x16_SSE41__k10), base.Simd_g_const(&F_Disto16x16_SSE41__k11))
			v132 = base.Simd_g_i16x8_add(v125, v131)
			v133 = base.Simd_g_const(&F_Disto16x16_SSE41__k15)
			v134 = base.Simd_g_i8x16_shuffle2(v117, v123, base.Simd_g_const(&F_Disto16x16_SSE41__k16), base.Simd_g_const(&F_Disto16x16_SSE41__k17))
			v136 = base.Simd_g_i8x16_shuffle2(v127, v129, base.Simd_g_const(&F_Disto16x16_SSE41__k16), base.Simd_g_const(&F_Disto16x16_SSE41__k17))
			v137 = base.Simd_g_i16x8_add(v134, v136)
			v138 = base.Simd_g_i16x8_add(v132, v137)
			v139 = base.Simd_g_i16x8_sub(v125, v131)
			v140 = base.Simd_g_i16x8_sub(v134, v136)
			v141 = base.Simd_g_i16x8_add(v139, v140)
			v146 = base.Simd_g_i16x8_sub(v139, v140)
			v147 = base.Simd_g_i16x8_sub(v132, v137)
			v162 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(v138, v141, base.Simd_g_const(&F_Disto16x16_SSE41__k10), base.Simd_g_const(&F_Disto16x16_SSE41__k11))), v22), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(v146, v147, base.Simd_g_const(&F_Disto16x16_SSE41__k10), base.Simd_g_const(&F_Disto16x16_SSE41__k11))), v20)), base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(v146, v147, base.Simd_g_const(&F_Disto16x16_SSE41__k16), base.Simd_g_const(&F_Disto16x16_SSE41__k17))), v20), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(v138, v141, base.Simd_g_const(&F_Disto16x16_SSE41__k16), base.Simd_g_const(&F_Disto16x16_SSE41__k17))), v22)))
			v173 = base.Simd_g_i32x4_extract_lane_l0(v162) + base.Simd_g_i32x4_extract_lane_l1(v162) + base.Simd_g_i32x4_extract_lane_l2(v162) + base.Simd_g_i32x4_extract_lane_l3(v162)
			v175 = v173 >> (uint(int32(31)) % 32)
			v180 = int32(base.Ui32(v173^v175-v175)>>(uint(int32(5))%32)) + v48
			v182 = v44 + v94
			if base.Ui32(v182) < base.Ui32(int32(12)) {
				v44 = v182
				v48 = v180
				continue
			} else {
				break
			}
			break
		}
		v185 = int32(128)
		if base.Ui32(v30) < base.Ui32(int32(384)) {
			v25 = v25 + v185
			v26 = v26 + v185
			v30 = v30 + v185
			v31 = v180
			continue
		} else {
			break
		}
		break
	}
	return v180
}

var F_Disto16x16_SSE41__k0 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_Disto16x16_SSE41__k1 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_Disto16x16_SSE41__k2 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_Disto16x16_SSE41__k3 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_Disto16x16_SSE41__k4 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_Disto16x16_SSE41__k5 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_Disto16x16_SSE41__k6 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_Disto16x16_SSE41__k7 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_Disto16x16_SSE41__k8 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_Disto16x16_SSE41__k9 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_Disto16x16_SSE41__k10 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_Disto16x16_SSE41__k11 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_Disto16x16_SSE41__k12 = [2]uint64{0x1b1a19180b0a0908, 0x1f1e1d1c0f0e0d0c}
var F_Disto16x16_SSE41__k13 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_Disto16x16_SSE41__k14 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_Disto16x16_SSE41__k15 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_Disto16x16_SSE41__k16 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_Disto16x16_SSE41__k17 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}

func F_Disto4x4_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	var v4 base.V128
	_ = v4
	var v11 int32
	_ = v11
	var v12 base.V128
	_ = v12
	var v14 base.V128
	_ = v14
	var v15 base.V128
	_ = v15
	var v18 base.V128
	_ = v18
	var v19 base.V128
	_ = v19
	var v20 int32
	_ = v20
	var v21 base.V128
	_ = v21
	var v23 base.V128
	_ = v23
	var v27 base.V128
	_ = v27
	var v28 base.V128
	_ = v28
	var v29 int32
	_ = v29
	var v30 base.V128
	_ = v30
	var v32 base.V128
	_ = v32
	var v36 base.V128
	_ = v36
	var v37 int32
	_ = v37
	var v38 base.V128
	_ = v38
	var v40 base.V128
	_ = v40
	var v44 base.V128
	_ = v44
	var v45 base.V128
	_ = v45
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v48 base.V128
	_ = v48
	var v49 base.V128
	_ = v49
	var v50 base.V128
	_ = v50
	var v51 base.V128
	_ = v51
	var v52 base.V128
	_ = v52
	var v53 base.V128
	_ = v53
	var v55 base.V128
	_ = v55
	var v57 base.V128
	_ = v57
	var v58 base.V128
	_ = v58
	var v59 base.V128
	_ = v59
	var v61 base.V128
	_ = v61
	var v63 base.V128
	_ = v63
	var v64 base.V128
	_ = v64
	var v65 base.V128
	_ = v65
	var v66 base.V128
	_ = v66
	var v67 base.V128
	_ = v67
	var v69 base.V128
	_ = v69
	var v71 base.V128
	_ = v71
	var v72 base.V128
	_ = v72
	var v73 base.V128
	_ = v73
	var v74 base.V128
	_ = v74
	var v76 base.V128
	_ = v76
	var v77 base.V128
	_ = v77
	var v78 base.V128
	_ = v78
	var v79 base.V128
	_ = v79
	var v80 base.V128
	_ = v80
	var v81 base.V128
	_ = v81
	var v86 base.V128
	_ = v86
	var v88 base.V128
	_ = v88
	var v89 base.V128
	_ = v89
	var v96 base.V128
	_ = v96
	var v108 base.V128
	_ = v108
	var v119 int32
	_ = v119
	var v121 int32
	_ = v121
	v4 = base.Simd_g_const(&F_Disto4x4_SSE2__k0)
	v11 = int32(96)
	v12 = base.Simd_g_v128_load64_zero(m, l0, v11)
	v14 = base.Simd_g_v128_load64_zero(m, l1, v11)
	v15 = base.Simd_g_const(&F_Disto4x4_SSE2__k1)
	v18 = base.Simd_g_const(&F_Disto4x4_SSE2__k2)
	v19 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v12, v14, base.Simd_g_const(&F_Disto4x4_SSE2__k3), base.Simd_g_const(&F_Disto4x4_SSE2__k4)), v4, base.Simd_g_const(&F_Disto4x4_SSE2__k5), base.Simd_g_const(&F_Disto4x4_SSE2__k6))
	v20 = int32(32)
	v21 = base.Simd_g_v128_load64_zero(m, l0, v20)
	v23 = base.Simd_g_v128_load64_zero(m, l1, v20)
	v27 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v21, v23, base.Simd_g_const(&F_Disto4x4_SSE2__k3), base.Simd_g_const(&F_Disto4x4_SSE2__k4)), v4, base.Simd_g_const(&F_Disto4x4_SSE2__k5), base.Simd_g_const(&F_Disto4x4_SSE2__k6))
	v28 = base.Simd_g_i16x8_add(v19, v27)
	v29 = int32(64)
	v30 = base.Simd_g_v128_load64_zero(m, l0, v29)
	v32 = base.Simd_g_v128_load64_zero(m, l1, v29)
	v36 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v30, v32, base.Simd_g_const(&F_Disto4x4_SSE2__k3), base.Simd_g_const(&F_Disto4x4_SSE2__k4)), v4, base.Simd_g_const(&F_Disto4x4_SSE2__k5), base.Simd_g_const(&F_Disto4x4_SSE2__k6))
	v37 = int32(0)
	v38 = base.Simd_g_v128_load64_zero(m, l0, v37)
	v40 = base.Simd_g_v128_load64_zero(m, l1, v37)
	v44 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v38, v40, base.Simd_g_const(&F_Disto4x4_SSE2__k3), base.Simd_g_const(&F_Disto4x4_SSE2__k4)), v4, base.Simd_g_const(&F_Disto4x4_SSE2__k5), base.Simd_g_const(&F_Disto4x4_SSE2__k6))
	v45 = base.Simd_g_i16x8_add(v36, v44)
	v46 = base.Simd_g_i16x8_add(v28, v45)
	v47 = base.Simd_g_i16x8_sub(v27, v19)
	v48 = base.Simd_g_i16x8_sub(v44, v36)
	v49 = base.Simd_g_i16x8_add(v47, v48)
	v50 = base.Simd_g_const(&F_Disto4x4_SSE2__k7)
	v51 = base.Simd_g_i8x16_shuffle2(v46, v49, base.Simd_g_const(&F_Disto4x4_SSE2__k8), base.Simd_g_const(&F_Disto4x4_SSE2__k9))
	v52 = base.Simd_g_i16x8_sub(v48, v47)
	v53 = base.Simd_g_i16x8_sub(v45, v28)
	v55 = base.Simd_g_i8x16_shuffle2(v52, v53, base.Simd_g_const(&F_Disto4x4_SSE2__k8), base.Simd_g_const(&F_Disto4x4_SSE2__k9))
	v57 = base.Simd_g_i8x16_shuffle2(v51, v55, base.Simd_g_const(&F_Disto4x4_SSE2__k3), base.Simd_g_const(&F_Disto4x4_SSE2__k4))
	v58 = base.Simd_g_const(&F_Disto4x4_SSE2__k10)
	v59 = base.Simd_g_i8x16_shuffle2(v46, v49, base.Simd_g_const(&F_Disto4x4_SSE2__k11), base.Simd_g_const(&F_Disto4x4_SSE2__k12))
	v61 = base.Simd_g_i8x16_shuffle2(v52, v53, base.Simd_g_const(&F_Disto4x4_SSE2__k11), base.Simd_g_const(&F_Disto4x4_SSE2__k12))
	v63 = base.Simd_g_i8x16_shuffle2(v59, v61, base.Simd_g_const(&F_Disto4x4_SSE2__k3), base.Simd_g_const(&F_Disto4x4_SSE2__k4))
	v64 = base.Simd_g_const(&F_Disto4x4_SSE2__k13)
	v65 = base.Simd_g_i8x16_shuffle2(v57, v63, base.Simd_g_const(&F_Disto4x4_SSE2__k14), base.Simd_g_const(&F_Disto4x4_SSE2__k15))
	v66 = base.Simd_g_const(&F_Disto4x4_SSE2__k16)
	v67 = base.Simd_g_i8x16_shuffle2(v51, v55, base.Simd_g_const(&F_Disto4x4_SSE2__k17), base.Simd_g_const(&F_Disto4x4_SSE2__k18))
	v69 = base.Simd_g_i8x16_shuffle2(v59, v61, base.Simd_g_const(&F_Disto4x4_SSE2__k17), base.Simd_g_const(&F_Disto4x4_SSE2__k18))
	v71 = base.Simd_g_i8x16_shuffle2(v67, v69, base.Simd_g_const(&F_Disto4x4_SSE2__k14), base.Simd_g_const(&F_Disto4x4_SSE2__k15))
	v72 = base.Simd_g_i16x8_add(v65, v71)
	v73 = base.Simd_g_const(&F_Disto4x4_SSE2__k19)
	v74 = base.Simd_g_i8x16_shuffle2(v57, v63, base.Simd_g_const(&F_Disto4x4_SSE2__k20), base.Simd_g_const(&F_Disto4x4_SSE2__k21))
	v76 = base.Simd_g_i8x16_shuffle2(v67, v69, base.Simd_g_const(&F_Disto4x4_SSE2__k20), base.Simd_g_const(&F_Disto4x4_SSE2__k21))
	v77 = base.Simd_g_i16x8_add(v74, v76)
	v78 = base.Simd_g_i16x8_add(v72, v77)
	v79 = base.Simd_g_i16x8_sub(v65, v71)
	v80 = base.Simd_g_i16x8_sub(v74, v76)
	v81 = base.Simd_g_i16x8_add(v79, v80)
	v86 = base.Simd_g_v128_load_rng(m, l2, v37, int32(0), int32(32))
	v88 = base.Simd_g_i16x8_sub(v79, v80)
	v89 = base.Simd_g_i16x8_sub(v72, v77)
	v96 = base.Simd_g_v128_load_nc(m, l2+int32(16), v37)
	v108 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(v78, v81, base.Simd_g_const(&F_Disto4x4_SSE2__k14), base.Simd_g_const(&F_Disto4x4_SSE2__k15))), v86), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(v88, v89, base.Simd_g_const(&F_Disto4x4_SSE2__k14), base.Simd_g_const(&F_Disto4x4_SSE2__k15))), v96)), base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(v78, v81, base.Simd_g_const(&F_Disto4x4_SSE2__k20), base.Simd_g_const(&F_Disto4x4_SSE2__k21))), v86), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(v88, v89, base.Simd_g_const(&F_Disto4x4_SSE2__k20), base.Simd_g_const(&F_Disto4x4_SSE2__k21))), v96)))
	v119 = base.Simd_g_i32x4_extract_lane_l0(v108) + base.Simd_g_i32x4_extract_lane_l1(v108) + base.Simd_g_i32x4_extract_lane_l2(v108) + base.Simd_g_i32x4_extract_lane_l3(v108)
	v121 = v119 >> (uint(int32(31)) % 32)
	return int32(base.Ui32(v119^v121-v121) >> (uint(int32(5)) % 32))
}

var F_Disto4x4_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_Disto4x4_SSE2__k1 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_Disto4x4_SSE2__k2 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_Disto4x4_SSE2__k3 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_Disto4x4_SSE2__k4 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_Disto4x4_SSE2__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_Disto4x4_SSE2__k6 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_Disto4x4_SSE2__k7 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_Disto4x4_SSE2__k8 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_Disto4x4_SSE2__k9 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_Disto4x4_SSE2__k10 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_Disto4x4_SSE2__k11 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_Disto4x4_SSE2__k12 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_Disto4x4_SSE2__k13 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_Disto4x4_SSE2__k14 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_Disto4x4_SSE2__k15 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_Disto4x4_SSE2__k16 = [2]uint64{0x1b1a19180b0a0908, 0x1f1e1d1c0f0e0d0c}
var F_Disto4x4_SSE2__k17 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_Disto4x4_SSE2__k18 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_Disto4x4_SSE2__k19 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_Disto4x4_SSE2__k20 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_Disto4x4_SSE2__k21 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}

func F_Disto4x4_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	var v11 int32
	_ = v11
	var v12 base.V128
	_ = v12
	var v14 base.V128
	_ = v14
	var v15 base.V128
	_ = v15
	var v17 base.V128
	_ = v17
	var v18 int32
	_ = v18
	var v19 base.V128
	_ = v19
	var v21 base.V128
	_ = v21
	var v24 base.V128
	_ = v24
	var v25 base.V128
	_ = v25
	var v26 int32
	_ = v26
	var v27 base.V128
	_ = v27
	var v29 base.V128
	_ = v29
	var v32 base.V128
	_ = v32
	var v33 int32
	_ = v33
	var v34 base.V128
	_ = v34
	var v36 base.V128
	_ = v36
	var v39 base.V128
	_ = v39
	var v40 base.V128
	_ = v40
	var v41 base.V128
	_ = v41
	var v42 base.V128
	_ = v42
	var v43 base.V128
	_ = v43
	var v44 base.V128
	_ = v44
	var v45 base.V128
	_ = v45
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v48 base.V128
	_ = v48
	var v50 base.V128
	_ = v50
	var v52 base.V128
	_ = v52
	var v53 base.V128
	_ = v53
	var v54 base.V128
	_ = v54
	var v56 base.V128
	_ = v56
	var v58 base.V128
	_ = v58
	var v59 base.V128
	_ = v59
	var v60 base.V128
	_ = v60
	var v61 base.V128
	_ = v61
	var v62 base.V128
	_ = v62
	var v64 base.V128
	_ = v64
	var v66 base.V128
	_ = v66
	var v67 base.V128
	_ = v67
	var v68 base.V128
	_ = v68
	var v69 base.V128
	_ = v69
	var v71 base.V128
	_ = v71
	var v72 base.V128
	_ = v72
	var v73 base.V128
	_ = v73
	var v74 base.V128
	_ = v74
	var v75 base.V128
	_ = v75
	var v76 base.V128
	_ = v76
	var v81 base.V128
	_ = v81
	var v83 base.V128
	_ = v83
	var v84 base.V128
	_ = v84
	var v91 base.V128
	_ = v91
	var v103 base.V128
	_ = v103
	var v114 int32
	_ = v114
	var v116 int32
	_ = v116
	v11 = int32(96)
	v12 = base.Simd_g_v128_load64_zero(m, l0, v11)
	v14 = base.Simd_g_v128_load64_zero(m, l1, v11)
	v15 = base.Simd_g_const(&F_Disto4x4_SSE41__k0)
	v17 = base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v12, v14, base.Simd_g_const(&F_Disto4x4_SSE41__k1), base.Simd_g_const(&F_Disto4x4_SSE41__k2)))
	v18 = int32(32)
	v19 = base.Simd_g_v128_load_rng(m, l0, v18, int32(0), int32(80))
	v21 = base.Simd_g_v128_load_rng(m, l1, v18, int32(0), int32(80))
	v24 = base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v19, v21, base.Simd_g_const(&F_Disto4x4_SSE41__k1), base.Simd_g_const(&F_Disto4x4_SSE41__k2)))
	v25 = base.Simd_g_i16x8_add(v17, v24)
	v26 = int32(64)
	v27 = base.Simd_g_v128_load_nc(m, l0, v26)
	v29 = base.Simd_g_v128_load_nc(m, l1, v26)
	v32 = base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v27, v29, base.Simd_g_const(&F_Disto4x4_SSE41__k1), base.Simd_g_const(&F_Disto4x4_SSE41__k2)))
	v33 = int32(0)
	v34 = base.Simd_g_v128_load_nc(m, l0, v33)
	v36 = base.Simd_g_v128_load_nc(m, l1, v33)
	v39 = base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v34, v36, base.Simd_g_const(&F_Disto4x4_SSE41__k1), base.Simd_g_const(&F_Disto4x4_SSE41__k2)))
	v40 = base.Simd_g_i16x8_add(v32, v39)
	v41 = base.Simd_g_i16x8_add(v25, v40)
	v42 = base.Simd_g_i16x8_sub(v24, v17)
	v43 = base.Simd_g_i16x8_sub(v39, v32)
	v44 = base.Simd_g_i16x8_add(v42, v43)
	v45 = base.Simd_g_const(&F_Disto4x4_SSE41__k3)
	v46 = base.Simd_g_i8x16_shuffle2(v41, v44, base.Simd_g_const(&F_Disto4x4_SSE41__k4), base.Simd_g_const(&F_Disto4x4_SSE41__k5))
	v47 = base.Simd_g_i16x8_sub(v43, v42)
	v48 = base.Simd_g_i16x8_sub(v40, v25)
	v50 = base.Simd_g_i8x16_shuffle2(v47, v48, base.Simd_g_const(&F_Disto4x4_SSE41__k4), base.Simd_g_const(&F_Disto4x4_SSE41__k5))
	v52 = base.Simd_g_i8x16_shuffle2(v46, v50, base.Simd_g_const(&F_Disto4x4_SSE41__k1), base.Simd_g_const(&F_Disto4x4_SSE41__k2))
	v53 = base.Simd_g_const(&F_Disto4x4_SSE41__k6)
	v54 = base.Simd_g_i8x16_shuffle2(v41, v44, base.Simd_g_const(&F_Disto4x4_SSE41__k7), base.Simd_g_const(&F_Disto4x4_SSE41__k8))
	v56 = base.Simd_g_i8x16_shuffle2(v47, v48, base.Simd_g_const(&F_Disto4x4_SSE41__k7), base.Simd_g_const(&F_Disto4x4_SSE41__k8))
	v58 = base.Simd_g_i8x16_shuffle2(v54, v56, base.Simd_g_const(&F_Disto4x4_SSE41__k1), base.Simd_g_const(&F_Disto4x4_SSE41__k2))
	v59 = base.Simd_g_const(&F_Disto4x4_SSE41__k9)
	v60 = base.Simd_g_i8x16_shuffle2(v52, v58, base.Simd_g_const(&F_Disto4x4_SSE41__k10), base.Simd_g_const(&F_Disto4x4_SSE41__k11))
	v61 = base.Simd_g_const(&F_Disto4x4_SSE41__k12)
	v62 = base.Simd_g_i8x16_shuffle2(v46, v50, base.Simd_g_const(&F_Disto4x4_SSE41__k13), base.Simd_g_const(&F_Disto4x4_SSE41__k14))
	v64 = base.Simd_g_i8x16_shuffle2(v54, v56, base.Simd_g_const(&F_Disto4x4_SSE41__k13), base.Simd_g_const(&F_Disto4x4_SSE41__k14))
	v66 = base.Simd_g_i8x16_shuffle2(v62, v64, base.Simd_g_const(&F_Disto4x4_SSE41__k10), base.Simd_g_const(&F_Disto4x4_SSE41__k11))
	v67 = base.Simd_g_i16x8_add(v60, v66)
	v68 = base.Simd_g_const(&F_Disto4x4_SSE41__k15)
	v69 = base.Simd_g_i8x16_shuffle2(v52, v58, base.Simd_g_const(&F_Disto4x4_SSE41__k16), base.Simd_g_const(&F_Disto4x4_SSE41__k17))
	v71 = base.Simd_g_i8x16_shuffle2(v62, v64, base.Simd_g_const(&F_Disto4x4_SSE41__k16), base.Simd_g_const(&F_Disto4x4_SSE41__k17))
	v72 = base.Simd_g_i16x8_add(v69, v71)
	v73 = base.Simd_g_i16x8_add(v67, v72)
	v74 = base.Simd_g_i16x8_sub(v60, v66)
	v75 = base.Simd_g_i16x8_sub(v69, v71)
	v76 = base.Simd_g_i16x8_add(v74, v75)
	v81 = base.Simd_g_v128_load_rng(m, l2, v33, int32(0), int32(32))
	v83 = base.Simd_g_i16x8_sub(v74, v75)
	v84 = base.Simd_g_i16x8_sub(v67, v72)
	v91 = base.Simd_g_v128_load_nc(m, l2+int32(16), v33)
	v103 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(v73, v76, base.Simd_g_const(&F_Disto4x4_SSE41__k10), base.Simd_g_const(&F_Disto4x4_SSE41__k11))), v81), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(v83, v84, base.Simd_g_const(&F_Disto4x4_SSE41__k10), base.Simd_g_const(&F_Disto4x4_SSE41__k11))), v91)), base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(v73, v76, base.Simd_g_const(&F_Disto4x4_SSE41__k16), base.Simd_g_const(&F_Disto4x4_SSE41__k17))), v81), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(v83, v84, base.Simd_g_const(&F_Disto4x4_SSE41__k16), base.Simd_g_const(&F_Disto4x4_SSE41__k17))), v91)))
	v114 = base.Simd_g_i32x4_extract_lane_l0(v103) + base.Simd_g_i32x4_extract_lane_l1(v103) + base.Simd_g_i32x4_extract_lane_l2(v103) + base.Simd_g_i32x4_extract_lane_l3(v103)
	v116 = v114 >> (uint(int32(31)) % 32)
	return int32(base.Ui32(v114^v116-v116) >> (uint(int32(5)) % 32))
}

var F_Disto4x4_SSE41__k0 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_Disto4x4_SSE41__k1 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_Disto4x4_SSE41__k2 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_Disto4x4_SSE41__k3 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_Disto4x4_SSE41__k4 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_Disto4x4_SSE41__k5 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_Disto4x4_SSE41__k6 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_Disto4x4_SSE41__k7 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_Disto4x4_SSE41__k8 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_Disto4x4_SSE41__k9 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_Disto4x4_SSE41__k10 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_Disto4x4_SSE41__k11 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_Disto4x4_SSE41__k12 = [2]uint64{0x1b1a19180b0a0908, 0x1f1e1d1c0f0e0d0c}
var F_Disto4x4_SSE41__k13 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_Disto4x4_SSE41__k14 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_Disto4x4_SSE41__k15 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_Disto4x4_SSE41__k16 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_Disto4x4_SSE41__k17 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}

func F_DoSegmentsJob(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v93 int32
	_ = v93
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v101 int32
	_ = v101
	var v103 int32
	_ = v103
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v111 int32
	_ = v111
	var v112 int32
	_ = v112
	var v116 int32
	_ = v116
	var v119 int32
	_ = v119
	var v120 int32
	_ = v120
	var v127 int32
	_ = v127
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v138 int32
	_ = v138
	var v140 int32
	_ = v140
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v145 float32
	_ = v145
	var v146 int32
	_ = v146
	var v147 int32
	_ = v147
	var v150 int32
	_ = v150
	var v152 int32
	_ = v152
	var v155 int32
	_ = v155
	var v157 int32
	_ = v157
	var v160 int32
	_ = v160
	var v162 int32
	_ = v162
	var v165 int32
	_ = v165
	var v172 int32
	_ = v172
	var v174 int32
	_ = v174
	var v175 int32
	_ = v175
	var v177 int32
	_ = v177
	var v180 int32
	_ = v180
	var v183 int32
	_ = v183
	var v186 int32
	_ = v186
	var v189 int32
	_ = v189
	var v192 int32
	_ = v192
	var v195 int32
	_ = v195
	var v198 int32
	_ = v198
	var v201 int32
	_ = v201
	var v204 int32
	_ = v204
	var v207 int32
	_ = v207
	var v210 int32
	_ = v210
	var v213 int32
	_ = v213
	var v216 int32
	_ = v216
	var v219 int32
	_ = v219
	var v225 int32
	_ = v225
	var v243 int32
	_ = v243
	var v246 int32
	_ = v246
	var v248 int32
	_ = v248
	var v254 int32
	_ = v254
	var v255 int32
	_ = v255
	var v256 int32
	_ = v256
	var v258 int32
	_ = v258
	var v259 int32
	_ = v259
	var v260 int32
	_ = v260
	var v262 int32
	_ = v262
	var v263 int32
	_ = v263
	var v266 int32
	_ = v266
	var v267 int32
	_ = v267
	var v271 int32
	_ = v271
	var v274 base.V128
	_ = v274
	var v275 int32
	_ = v275
	var v278 int32
	_ = v278
	var v280 int32
	_ = v280
	var v281 int32
	_ = v281
	var v283 int32
	_ = v283
	var v284 int32
	_ = v284
	var v285 int32
	_ = v285
	var v286 int32
	_ = v286
	var v288 int32
	_ = v288
	var v289 int32
	_ = v289
	var v290 int32
	_ = v290
	var v291 int32
	_ = v291
	var v293 int32
	_ = v293
	var v294 int32
	_ = v294
	var v296 int32
	_ = v296
	var v298 int32
	_ = v298
	var v299 int32
	_ = v299
	var v301 int32
	_ = v301
	var v307 int32
	_ = v307
	var v308 int32
	_ = v308
	var v309 int32
	_ = v309
	var v310 int32
	_ = v310
	var v316 int32
	_ = v316
	var v319 int32
	_ = v319
	var v322 int32
	_ = v322
	var v325 int32
	_ = v325
	var v326 int32
	_ = v326
	var v329 int32
	_ = v329
	var v330 int32
	_ = v330
	var v331 int32
	_ = v331
	var v337 int32
	_ = v337
	var v339 int32
	_ = v339
	var v342 int32
	_ = v342
	var v345 int32
	_ = v345
	var v346 int32
	_ = v346
	var v347 int32
	_ = v347
	var v350 int32
	_ = v350
	var v351 int32
	_ = v351
	var v353 int32
	_ = v353
	var v357 int32
	_ = v357
	var v359 int32
	_ = v359
	var v360 int32
	_ = v360
	var v361 int32
	_ = v361
	var v363 int32
	_ = v363
	var v364 int32
	_ = v364
	var v365 int32
	_ = v365
	var v367 int32
	_ = v367
	var v368 int32
	_ = v368
	var v371 int32
	_ = v371
	var v372 int32
	_ = v372
	var v376 int32
	_ = v376
	var v378 int32
	_ = v378
	var v387 int32
	_ = v387
	var v404 int32
	_ = v404
	var v405 int32
	_ = v405
	var v407 int32
	_ = v407
	var v413 int32
	_ = v413
	var v414 int32
	_ = v414
	var v416 int32
	_ = v416
	var v418 int32
	_ = v418
	var v421 int32
	_ = v421
	var v424 int32
	_ = v424
	var v425 int32
	_ = v425
	var v428 int32
	_ = v428
	var v429 int32
	_ = v429
	var v431 int32
	_ = v431
	var v437 int32
	_ = v437
	var v438 int32
	_ = v438
	var v440 int32
	_ = v440
	var v443 int32
	_ = v443
	var v446 int32
	_ = v446
	var v447 int32
	_ = v447
	var v449 int32
	_ = v449
	var v450 int32
	_ = v450
	var v457 int32
	_ = v457
	var v459 int32
	_ = v459
	var v461 int32
	_ = v461
	var v462 int32
	_ = v462
	var v465 int32
	_ = v465
	var v468 int32
	_ = v468
	var v470 int32
	_ = v470
	var v473 int32
	_ = v473
	var v474 int32
	_ = v474
	var v478 int32
	_ = v478
	var v481 int32
	_ = v481
	var v483 int32
	_ = v483
	var v488 int32
	_ = v488
	var v489 int32
	_ = v489
	var v490 int32
	_ = v490
	var v493 int32
	_ = v493
	var v500 int32
	_ = v500
	var v501 int32
	_ = v501
	var v502 int32
	_ = v502
	var v505 int32
	_ = v505
	var v508 int32
	_ = v508
	var v509 int32
	_ = v509
	var v512 int32
	_ = v512
	var v513 int32
	_ = v513
	var v515 int32
	_ = v515
	var v518 int32
	_ = v518
	var v520 int32
	_ = v520
	var v529 int32
	_ = v529
	var v531 int32
	_ = v531
	var v533 int32
	_ = v533
	var v534 int32
	_ = v534
	var v538 int32
	_ = v538
	var v540 int64
	_ = v540
	var v542 int32
	_ = v542
	var v544 int32
	_ = v544
	var v546 int32
	_ = v546
	var v547 int32
	_ = v547
	var v549 int32
	_ = v549
	var v553 int32
	_ = v553
	var v559 int32
	_ = v559
	var v560 int32
	_ = v560
	var v569 int32
	_ = v569
	var v576 int32
	_ = v576
	var v578 int32
	_ = v578
	var v579 int32
	_ = v579
	var v582 int32
	_ = v582
	var v586 int32
	_ = v586
	var v587 int64
	_ = v587
	var v593 int32
	_ = v593
	var v596 int32
	_ = v596
	var v599 int32
	_ = v599
	var v601 int32
	_ = v601
	var v606 int32
	_ = v606
	var v607 int32
	_ = v607
	var v610 int32
	_ = v610
	var v614 int32
	_ = v614
	var v618 int32
	_ = v618
	var v619 int32
	_ = v619
	var v622 int32
	_ = v622
	var v629 int32
	_ = v629
	var v638 int32
	_ = v638
	v31 = m.G0
	v33 = v31 - int32(144)
	m.G0 = v33
	v35 = int32(1)
	v36 = *(*int32)(unsafe.Add(mBase, uint32(l1)+288))
	goto L2
L1:
	;
	m.G0 = v33 + int32(144)
	return v638
L2:
	;
	if v36 < v35 {
		v638 = v35
		goto L1
	} else {
		goto L3
	}
L3:
	;
	v45 = m.G89
	v46 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v45)+2)))
	v47 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v45))))
	v56 = m.G90
	v57 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v56)+2)))
	v58 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v56))))
	goto L4
L4:
	;
	F_VP8IteratorImport(m, l1, (v33+int32(31))&int32(-32))
	mBase = m.M
	v90 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v91 = int32(0)
	v93 = *(*int32)(unsafe.Add(mBase, uint32(l1)+36))
	*(*int32)(unsafe.Add(mBase, uint32(v93))) = v91
	v99 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v100 = *(*int32)(unsafe.Add(mBase, uint32(v99)+48))
	v101 = v93 + v100
	*(*int32)(unsafe.Add(mBase, uint32(v101))) = v91
	v103 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v104 = *(*int32)(unsafe.Add(mBase, uint32(v103)+48))
	v105 = v101 + v104
	*(*int32)(unsafe.Add(mBase, uint32(v105))) = v91
	v107 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v108 = *(*int32)(unsafe.Add(mBase, uint32(v107)+48))
	*(*int32)(unsafe.Add(mBase, uint32(v105+v108))) = v91
	v111 = *(*int32)(unsafe.Add(mBase, uint32(l1)+28))
	v112 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v111))))
	v116 = v112&int32(252) | int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v111))) = uint8(v116)
	goto L6
L5:
	;
	v638 = v520
	goto L1
L6:
	;
	v119 = *(*int32)(unsafe.Add(mBase, uint32(l1)+28))
	v120 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v119))))
	v127 = v120&int32(239) | int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v119))) = uint8(v127)
	goto L7
L7:
	;
	v130 = *(*int32)(unsafe.Add(mBase, uint32(l1)+28))
	v131 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v130))))
	v138 = v131&int32(159) | int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v130))) = uint8(v138)
	goto L8
L8:
	;
	v140 = *(*int32)(unsafe.Add(mBase, uint32(v90)+uint32(_c_F_DoSegmentsJob[0])))
	if int32(1) < v140 {
		goto L10
	} else {
		goto L11
	}
L9:
	;
	F_VP8MakeChroma8Preds(m, l1)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v33)+80)) = int64(4294967296)
	v404 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v405 = int32(16)
	v407 = *(*int32)(unsafe.Add(mBase, uint32(l1)+20))
	v413 = m.G59
	v414 = *(*int32)(unsafe.Add(mBase, uint32(v413)))
	m.T0[v414].(func(*base.Module, int32, int32, int32, int32, int32))(m, v404+v405, v407+v47, v405, int32(24), v33+int32(80))
	mBase = m.M
	v416 = int32(0)
	v418 = *(*int32)(unsafe.Add(mBase, uint32(v33)+80))
	if v418 < int32(2) {
		v425 = v416
		goto L30
	} else {
		goto L31
	}
L10:
	;
	F_VP8MakeLuma16Preds(m, l1)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v33)+80)) = int64(4294967296)
	v307 = *(*int32)(unsafe.Add(mBase, uint32(l1)+20))
	v308 = m.G59
	v309 = int32(0)
	v310 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v316 = *(*int32)(unsafe.Add(mBase, uint32(v308)))
	m.T0[v316].(func(*base.Module, int32, int32, int32, int32, int32))(m, v310, v307+v58, v309, int32(16), v33+int32(80))
	mBase = m.M
	v319 = *(*int32)(unsafe.Add(mBase, uint32(v33)+80))
	if v319 < int32(2) {
		v326 = v309
		goto L19
	} else {
		goto L20
	}
L11:
	;
	v143 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v144 = *(*int32)(unsafe.Add(mBase, uint32(v143)))
	v145 = *(*float32)(unsafe.Add(mBase, uint32(v144)+4))
	v146 = m.G75
	v147 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v150 = *(*int32)(unsafe.Add(mBase, uint32(v146)))
	m.T0[v150].(func(*base.Module, int32, int32))(m, v147, v33+int32(80))
	mBase = m.M
	v152 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v155 = *(*int32)(unsafe.Add(mBase, uint32(v146)))
	m.T0[v155].(func(*base.Module, int32, int32))(m, v152+int32(128), v33+int32(96))
	mBase = m.M
	v157 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v160 = *(*int32)(unsafe.Add(mBase, uint32(v146)))
	m.T0[v160].(func(*base.Module, int32, int32))(m, v157+int32(256), v33+int32(112))
	mBase = m.M
	v162 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v165 = *(*int32)(unsafe.Add(mBase, uint32(v146)))
	m.T0[v165].(func(*base.Module, int32, int32))(m, v162+int32(384), v33+int32(128))
	mBase = m.M
	if base.F32_lt(base.F32_abs(v145), float32(2.1474836e+09)) == int32(0) {
		goto L13
	} else {
		goto L14
	}
L12:
	;
	v175 = *(*int32)(unsafe.Add(mBase, uint32(v33)+84))
	v177 = *(*int32)(unsafe.Add(mBase, uint32(v33)+80))
	v180 = *(*int32)(unsafe.Add(mBase, uint32(v33)+88))
	v183 = *(*int32)(unsafe.Add(mBase, uint32(v33)+92))
	v186 = *(*int32)(unsafe.Add(mBase, uint32(v33)+96))
	v189 = *(*int32)(unsafe.Add(mBase, uint32(v33)+100))
	v192 = *(*int32)(unsafe.Add(mBase, uint32(v33)+104))
	v195 = *(*int32)(unsafe.Add(mBase, uint32(v33)+108))
	v198 = *(*int32)(unsafe.Add(mBase, uint32(v33)+112))
	v201 = *(*int32)(unsafe.Add(mBase, uint32(v33)+116))
	v204 = *(*int32)(unsafe.Add(mBase, uint32(v33)+120))
	v207 = *(*int32)(unsafe.Add(mBase, uint32(v33)+124))
	v210 = *(*int32)(unsafe.Add(mBase, uint32(v33)+128))
	v213 = *(*int32)(unsafe.Add(mBase, uint32(v33)+132))
	v216 = *(*int32)(unsafe.Add(mBase, uint32(v33)+136))
	v219 = *(*int32)(unsafe.Add(mBase, uint32(v33)+140))
	v225 = base.I32_div_s(v174*int32(9), int32(100))
	v243 = v175 + v177 + v180 + v183 + v186 + v189 + v192 + v195 + v198 + v201 + v204 + v207 + v210 + v213 + v216 + v219
	if base.Ui32(v243*v243) <= base.Ui32((v175*v175+v177*v177+v180*v180+v183*v183+v186*v186+v189*v189+v192*v192+v195*v195+v198*v198+v201*v201+v204*v204+v207*v207+v210*v210+v213*v213+v216*v216+v219*v219)*(v225+int32(8))) {
		goto L15
	} else {
		goto L16
	}
L13:
	;
	v174 = int32(-2147483648)
	goto L12
L14:
	;
	v172 = base.I32_trunc_f32_s(v145)
	v174 = v172
	goto L12
L15:
	;
	v274 = base.Simd_g_const(&F_DoSegmentsJob__k0)
	v275 = int32(64)
	base.Simd_g_v128_store(m, v33, v275, v274)
	v278 = v33 + v275
	v280 = *(*int32)(unsafe.Add(mBase, uint32(l1)+36))
	v281 = *(*int32)(unsafe.Add(mBase, uint32(v278)))
	*(*int32)(unsafe.Add(mBase, uint32(v280))) = v281
	v283 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v284 = *(*int32)(unsafe.Add(mBase, uint32(v283)+48))
	v285 = v280 + v284
	v286 = *(*int32)(unsafe.Add(mBase, uint32(v278)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v285))) = v286
	v288 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v289 = *(*int32)(unsafe.Add(mBase, uint32(v288)+48))
	v290 = v285 + v289
	v291 = *(*int32)(unsafe.Add(mBase, uint32(v278)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v290))) = v291
	v293 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v294 = *(*int32)(unsafe.Add(mBase, uint32(v293)+48))
	v296 = *(*int32)(unsafe.Add(mBase, uint32(v278)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v290+v294))) = v296
	v298 = *(*int32)(unsafe.Add(mBase, uint32(l1)+28))
	v299 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v298))))
	v301 = v299 & int32(252)
	*(*uint8)(unsafe.Add(mBase, uint32(v298))) = uint8(v301)
	goto L18
L16:
	;
	v246 = int32(0)
	v248 = *(*int32)(unsafe.Add(mBase, uint32(l1)+36))
	*(*int32)(unsafe.Add(mBase, uint32(v248))) = v246
	v254 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v255 = *(*int32)(unsafe.Add(mBase, uint32(v254)+48))
	v256 = v248 + v255
	*(*int32)(unsafe.Add(mBase, uint32(v256))) = v246
	v258 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v259 = *(*int32)(unsafe.Add(mBase, uint32(v258)+48))
	v260 = v256 + v259
	*(*int32)(unsafe.Add(mBase, uint32(v260))) = v246
	v262 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v263 = *(*int32)(unsafe.Add(mBase, uint32(v262)+48))
	*(*int32)(unsafe.Add(mBase, uint32(v260+v263))) = v246
	v266 = *(*int32)(unsafe.Add(mBase, uint32(l1)+28))
	v267 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v266))))
	v271 = v267&int32(252) | int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v266))) = uint8(v271)
	goto L17
L17:
	;
	v387 = int32(2)
	goto L9
L18:
	;
	v387 = int32(2)
	goto L9
L19:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v33)+80)) = int64(4294967296)
	v329 = *(*int32)(unsafe.Add(mBase, uint32(l1)+20))
	v330 = m.G59
	v331 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v337 = *(*int32)(unsafe.Add(mBase, uint32(v330)))
	m.T0[v337].(func(*base.Module, int32, int32, int32, int32, int32))(m, v331, v329+v57, int32(0), int32(16), v33+int32(80))
	mBase = m.M
	v339 = *(*int32)(unsafe.Add(mBase, uint32(v33)+80))
	if v339 < int32(2) {
		v346 = v309
		goto L21
	} else {
		goto L22
	}
L20:
	;
	v322 = *(*int32)(unsafe.Add(mBase, uint32(v33)+84))
	v325 = base.I32_div_s(v322*int32(510), v319)
	v326 = v325
	goto L19
L21:
	;
	v347 = int32(-1)
	if v347 < v326 {
		goto L23
	} else {
		goto L24
	}
L22:
	;
	v342 = *(*int32)(unsafe.Add(mBase, uint32(v33)+84))
	v345 = base.I32_div_s(v342*int32(510), v339)
	v346 = v345
	goto L21
L23:
	;
	v350 = v326
	goto L25
L24:
	;
	v350 = v347
	goto L25
L25:
	;
	v351 = base.B2i32(v350 < v346)
	v353 = *(*int32)(unsafe.Add(mBase, uint32(l1)+36))
	v357 = v351 & int32(255) * int32(16843009)
	*(*int32)(unsafe.Add(mBase, uint32(v353))) = v357
	v359 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v360 = *(*int32)(unsafe.Add(mBase, uint32(v359)+48))
	v361 = v353 + v360
	*(*int32)(unsafe.Add(mBase, uint32(v361))) = v357
	v363 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v364 = *(*int32)(unsafe.Add(mBase, uint32(v363)+48))
	v365 = v361 + v364
	*(*int32)(unsafe.Add(mBase, uint32(v365))) = v357
	v367 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v368 = *(*int32)(unsafe.Add(mBase, uint32(v367)+48))
	*(*int32)(unsafe.Add(mBase, uint32(v365+v368))) = v357
	v371 = *(*int32)(unsafe.Add(mBase, uint32(l1)+28))
	v372 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371))))
	v376 = v372&int32(252) | int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v371))) = uint8(v376)
	goto L26
L26:
	;
	if v350 < v346 {
		goto L27
	} else {
		goto L28
	}
L27:
	;
	v378 = v346
	goto L29
L28:
	;
	v378 = v350
	goto L29
L29:
	;
	v387 = v378*int32(3) + int32(2)
	goto L9
L30:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v33)+80)) = int64(4294967296)
	v428 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v429 = int32(16)
	v431 = *(*int32)(unsafe.Add(mBase, uint32(l1)+20))
	v437 = m.G59
	v438 = *(*int32)(unsafe.Add(mBase, uint32(v437)))
	m.T0[v438].(func(*base.Module, int32, int32, int32, int32, int32))(m, v428+v429, v431+v46, v429, int32(24), v33+int32(80))
	mBase = m.M
	v440 = *(*int32)(unsafe.Add(mBase, uint32(v33)+80))
	if v440 < int32(2) {
		v447 = v416
		goto L32
	} else {
		goto L33
	}
L31:
	;
	v421 = *(*int32)(unsafe.Add(mBase, uint32(v33)+84))
	v424 = base.I32_div_s(v421*int32(510), v418)
	v425 = v424
	goto L30
L32:
	;
	v449 = *(*int32)(unsafe.Add(mBase, uint32(l1)+28))
	v450 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v449))))
	v457 = v450&int32(243) | base.B2i32(v447 < v425)<<(uint(int32(2))%32)&int32(12)
	*(*uint8)(unsafe.Add(mBase, uint32(v449))) = uint8(v457)
	goto L34
L33:
	;
	v443 = *(*int32)(unsafe.Add(mBase, uint32(v33)+84))
	v446 = base.I32_div_s(v443*int32(510), v440)
	v447 = v446
	goto L32
L34:
	;
	v459 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1052))
	if v447 < v425 {
		goto L35
	} else {
		goto L36
	}
L35:
	;
	v461 = v425
	goto L37
L36:
	;
	v461 = v447
	goto L37
L37:
	;
	v462 = int32(-1)
	if v462 < v461 {
		goto L38
	} else {
		goto L39
	}
L38:
	;
	v465 = v461
	goto L40
L39:
	;
	v465 = v462
	goto L40
L40:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+1052)) = v459 + v465
	v468 = *(*int32)(unsafe.Add(mBase, uint32(l1)+28))
	v470 = int32(255)
	v473 = (v465 + v387) >> (uint(int32(2)) % 32)
	v474 = v470 - v473
	if base.Ui32(v474) < base.Ui32(v470) {
		goto L41
	} else {
		goto L42
	}
L41:
	;
	v478 = v474
	goto L43
L42:
	;
	v478 = v470
	goto L43
L43:
	;
	if int32(255) < v473 {
		goto L44
	} else {
		goto L45
	}
L44:
	;
	v481 = int32(0)
	goto L46
L45:
	;
	v481 = v478
	goto L46
L46:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v468)+1)) = uint8(v481)
	v483 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1048))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+1048)) = v483 + v481
	v488 = l0 + int32(24) + v481<<(uint(int32(2))%32)
	v489 = *(*int32)(unsafe.Add(mBase, uint32(v488)))
	v490 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v488))) = v489 + v490
	v493 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_DoSegmentsJob[1])))
	if v493 == int32(0) {
		v520 = v490
		goto L48
	} else {
		goto L49
	}
L47:
	;
	if v520 == int32(0) {
		v638 = v520
		goto L1
	} else {
		goto L54
	}
L48:
	;
	goto L47
L49:
	;
	v500 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v501 = *(*int32)(unsafe.Add(mBase, uint32(v500)+4))
	v502 = *(*int32)(unsafe.Add(mBase, uint32(v501)+96))
	if v502 == int32(0) {
		v520 = v490
		goto L48
	} else {
		goto L50
	}
L50:
	;
	v505 = *(*int32)(unsafe.Add(mBase, uint32(l1)+292))
	if int32(0) < v505 {
		goto L52
	} else {
		goto L53
	}
L51:
	;
	v518 = F_WebPReportProgress(m, v501, v515, v500+int32(368))
	mBase = m.M
	v520 = v518
	goto L48
L52:
	;
	v509 = *(*int32)(unsafe.Add(mBase, uint32(l1)+288))
	v512 = base.I32_div_s((v505-v509)*v493, v505)
	v513 = *(*int32)(unsafe.Add(mBase, uint32(l1)+296))
	v515 = v512 + v513
	goto L51
L53:
	;
	v508 = *(*int32)(unsafe.Add(mBase, uint32(l1)+296))
	v515 = v508
	goto L51
L54:
	;
	v529 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	v531 = v529 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l1))) = v531
	v533 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v534 = *(*int32)(unsafe.Add(mBase, uint32(v533)+40))
	if v531 != v534 {
		goto L57
	} else {
		goto L58
	}
L55:
	;
	if int32(1) < v629 {
		goto L4
	} else {
		goto L63
	}
L56:
	;
	v629 = *(*int32)(unsafe.Add(mBase, uint32(l1)+288))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+288)) = v629 + int32(-1)
	goto L55
L57:
	;
	v606 = *(*int32)(unsafe.Add(mBase, uint32(l1)+36))
	v607 = int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l1)+36)) = v606 + v607
	v610 = *(*int32)(unsafe.Add(mBase, uint32(l1)+28))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+28)) = v610 + v607
	v614 = *(*int32)(unsafe.Add(mBase, uint32(l1)+40))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+40)) = v614 + v607
	v618 = *(*int32)(unsafe.Add(mBase, uint32(l1)+320))
	v619 = int32(16)
	*(*int32)(unsafe.Add(mBase, uint32(l1)+320)) = v618 + v619
	v622 = *(*int32)(unsafe.Add(mBase, uint32(l1)+324))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+324)) = v622 + v619
	goto L56
L58:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1))) = int32(0)
	v538 = *(*int32)(unsafe.Add(mBase, uint32(v533)+uint32(_c_F_DoSegmentsJob[2])))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+40)) = v538
	v540 = *(*int64)(unsafe.Add(mBase, uint32(v533)+uint32(_c_F_DoSegmentsJob[3])))
	*(*int64)(unsafe.Add(mBase, uint32(l1)+320)) = v540
	v542 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
	v544 = v542 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v544
	v546 = *(*int32)(unsafe.Add(mBase, uint32(v533)+uint32(_c_F_DoSegmentsJob[4])))
	v547 = *(*int32)(unsafe.Add(mBase, uint32(v533)+48))
	v549 = int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l1)+36)) = v546 + v544*v547<<(uint(v549)%32)
	v553 = *(*int32)(unsafe.Add(mBase, uint32(v533)+uint32(_c_F_DoSegmentsJob[5])))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+28)) = v553 + v544*v531<<(uint(v549)%32)
	v559 = *(*int32)(unsafe.Add(mBase, uint32(v533)+52))
	v560 = int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l1)+32)) = v533 + (v559+v560)&v544<<(uint(int32(5))%32) + int32(88)
	v569 = *(*int32)(unsafe.Add(mBase, uint32(l1)+316))
	if v560 < v542 {
		goto L59
	} else {
		goto L60
	}
L59:
	;
	v576 = int32(-127)
	goto L61
L60:
	;
	v576 = int32(127)
	goto L61
L61:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v569+v560))) = uint8(v576)
	v578 = *(*int32)(unsafe.Add(mBase, uint32(l1)+312))
	v579 = int32(-1)
	*(*uint8)(unsafe.Add(mBase, uint32(v578+v579))) = uint8(v576)
	v582 = *(*int32)(unsafe.Add(mBase, uint32(l1)+308))
	*(*uint8)(unsafe.Add(mBase, uint32(v582+v579))) = uint8(v576)
	v586 = *(*int32)(unsafe.Add(mBase, uint32(l1)+308))
	v587 = int64(-9114861777597660799)
	*(*int64)(unsafe.Add(mBase, uint32(v586))) = v587
	*(*int64)(unsafe.Add(mBase, uint32(v586+int32(8)))) = v587
	v593 = *(*int32)(unsafe.Add(mBase, uint32(l1)+312))
	*(*int64)(unsafe.Add(mBase, uint32(v593))) = v587
	v596 = *(*int32)(unsafe.Add(mBase, uint32(l1)+316))
	*(*int64)(unsafe.Add(mBase, uint32(v596))) = v587
	v599 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l1)+160)) = v599
	v601 = *(*int32)(unsafe.Add(mBase, uint32(l1)+304))
	if v601 == v599 {
		goto L56
	} else {
		goto L62
	}
L62:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+300)) = int32(0)
	goto L56
L63:
	;
	goto L5
}

var F_DoSegmentsJob__k0 = [2]uint64{0x0, 0x0}
