//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_TM16_SSE2(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v3 base.V128
	_ = v3
	var v13 base.V128
	_ = v13
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v33 base.V128
	_ = v33
	var v36 base.V128
	_ = v36
	var v40 int32
	_ = v40
	v2 = int32(0)
	v3 = base.Simd_g_const(&F_TM16_SSE2__k0)
	v13 = base.Simd_g_v128_load(m, l0+int32(-32), v2)
	v25 = v2
	for {
		v27 = l0 + v25
		v30 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v27+int32(-1)))))
		v31 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-33)))))
		v33 = base.Simd_g_i16x8_splat(v30 - v31)
		v36 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v33, base.Simd_g_i8x16_shuffle2(v13, v3, base.Simd_g_const(&F_TM16_SSE2__k1), base.Simd_g_const(&F_TM16_SSE2__k2))), base.Simd_g_i16x8_add(v33, base.Simd_g_i8x16_shuffle2(v13, v3, base.Simd_g_const(&F_TM16_SSE2__k3), base.Simd_g_const(&F_TM16_SSE2__k4))))
		base.Simd_g_v128_store(m, v27, int32(0), v36)
		v40 = v25 + int32(32)
		if v40 != int32(512) {
			v25 = v40
			continue
		} else {
			break
		}
		break
	}
	return
}

var F_TM16_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_TM16_SSE2__k1 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_TM16_SSE2__k2 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_TM16_SSE2__k3 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_TM16_SSE2__k4 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}

func F_TM4_SSE2(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v3 base.V128
	_ = v3
	var v7 int32
	_ = v7
	var v10 int32
	_ = v10
	var v16 base.V128
	_ = v16
	var v19 base.V128
	_ = v19
	var v26 int32
	_ = v26
	var v34 int32
	_ = v34
	var v42 int32
	_ = v42
	v2 = int32(0)
	v3 = base.Simd_g_const(&F_TM4_SSE2__k0)
	v7 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-1)))))
	v10 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-33)))))
	v16 = base.Simd_g_v128_load32_zero(m, l0+int32(-32), v2)
	v19 = base.Simd_g_i8x16_shuffle2(v16, v3, base.Simd_g_const(&F_TM4_SSE2__k1), base.Simd_g_const(&F_TM4_SSE2__k2))
	base.Simd_g_v128_store32_lane_l0(m, l0, v2, base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v7-v10), v19), v3))
	v26 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+31)))
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(32), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v26-v10), v19), v3))
	v34 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+63)))
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(64), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v34-v10), v19), v3))
	v42 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+95)))
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(96), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v42-v10), v19), v3))
	return
}

var F_TM4_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_TM4_SSE2__k1 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_TM4_SSE2__k2 = [2]uint64{0x380028001800080, 0x780068005800480}

func F_TM8uv_SSE2(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v3 base.V128
	_ = v3
	var v7 int32
	_ = v7
	var v10 int32
	_ = v10
	var v16 base.V128
	_ = v16
	var v19 base.V128
	_ = v19
	var v26 int32
	_ = v26
	var v34 int32
	_ = v34
	var v42 int32
	_ = v42
	var v50 int32
	_ = v50
	var v58 int32
	_ = v58
	var v66 int32
	_ = v66
	var v74 int32
	_ = v74
	v2 = int32(0)
	v3 = base.Simd_g_const(&F_TM8uv_SSE2__k0)
	v7 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-1)))))
	v10 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-33)))))
	v16 = base.Simd_g_v128_load64_zero(m, l0+int32(-32), v2)
	v19 = base.Simd_g_i8x16_shuffle2(v16, v3, base.Simd_g_const(&F_TM8uv_SSE2__k1), base.Simd_g_const(&F_TM8uv_SSE2__k2))
	base.Simd_g_v128_store64_lane_l0(m, l0, v2, base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v7-v10), v19), v3))
	v26 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+31)))
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(32), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v26-v10), v19), v3))
	v34 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+63)))
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(64), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v34-v10), v19), v3))
	v42 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+95)))
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(96), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v42-v10), v19), v3))
	v50 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+127)))
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(128), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v50-v10), v19), v3))
	v58 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+159)))
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(160), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v58-v10), v19), v3))
	v66 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+191)))
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(192), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v66-v10), v19), v3))
	v74 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+223)))
	base.Simd_g_v128_store64_lane_l0(m, l0, int32(224), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v74-v10), v19), v3))
	return
}

var F_TM8uv_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_TM8uv_SSE2__k1 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_TM8uv_SSE2__k2 = [2]uint64{0x380028001800080, 0x780068005800480}

func F_TransformColorInverse_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v24 int32
	_ = v24
	var v26 int32
	_ = v26
	var v35 base.V128
	_ = v35
	var v38 int32
	_ = v38
	var v44 base.V128
	_ = v44
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v64 int32
	_ = v64
	var v65 base.V128
	_ = v65
	var v67 base.V128
	_ = v67
	var v68 base.V128
	_ = v68
	var v72 base.V128
	_ = v72
	var v77 base.V128
	_ = v77
	var v80 int32
	_ = v80
	var v81 base.V128
	_ = v81
	var v93 base.V128
	_ = v93
	var v96 int32
	_ = v96
	var v103 int32
	_ = v103
	var v109 int32
	_ = v109
	var v123 int32
	_ = v123
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v126 int32
	_ = v126
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v154 int32
	_ = v154
	var v155 int32
	_ = v155
	var v166 int32
	_ = v166
	var v167 int32
	_ = v167
	var v171 int32
	_ = v171
	var v172 base.V128
	_ = v172
	var v173 int32
	_ = v173
	var v175 int32
	_ = v175
	var v176 base.V128
	_ = v176
	var v178 int32
	_ = v178
	var v182 base.V128
	_ = v182
	var v204 base.V128
	_ = v204
	var v212 int32
	_ = v212
	var v221 int32
	_ = v221
	var v232 int32
	_ = v232
	var v237 int32
	_ = v237
	var v246 int32
	_ = v246
	var v247 int32
	_ = v247
	var v251 int32
	_ = v251
	var v252 int32
	_ = v252
	var v255 int32
	_ = v255
	var v257 int32
	_ = v257
	var v261 int32
	_ = v261
	var v282 int32
	_ = v282
	var v287 int32
	_ = v287
	if l2 < int32(4) {
		v109 = int32(0)
	} else {
		v20 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
		v21 = int32(8)
		v24 = int32(11)
		v26 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)))
		v35 = base.Simd_g_i32x4_splat(base.I32_extend16_s(v20<<(uint(v21)%32))<<(uint(v24)%32) | int32(base.Ui32(base.I32_extend16_s(v26<<(uint(v21)%32)))>>(uint(int32(5))%32))&int32(_a_F_TransformColorInverse_SSE2_0))
		v38 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)))
		v44 = base.Simd_g_i32x4_splat(base.I32_extend16_s(v38<<(uint(v21)%32)) << (uint(v24) % 32))
		v58 = int32(0)
		v59 = l3
		v60 = l1
		for {
			v64 = int32(0)
			v65 = base.Simd_g_v128_load(m, v60, v64)
			v67 = base.Simd_g_v128_and(v65, base.Simd_g_const(&F_TransformColorInverse_SSE2__k0))
			v68 = base.Simd_g_const(&F_TransformColorInverse_SSE2__k1)
			v72 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(v67, base.Simd_g_const(&F_TransformColorInverse_SSE2__k2)), base.Simd_g_const(&F_TransformColorInverse_SSE2__k3))
			v77 = base.Simd_g_const(&F_TransformColorInverse_SSE2__k4)
			v80 = int32(8)
			v81 = base.Simd_g_i16x8_shl(base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v72), base.Simd_g_i32x4_extend_low_i16x8_s(v35)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v72), base.Simd_g_i32x4_extend_high_i16x8_s(v35)), base.Simd_g_const(&F_TransformColorInverse_SSE2__k5), base.Simd_g_const(&F_TransformColorInverse_SSE2__k6)), v65), v80)
			v93 = base.Simd_g_v128_or(v67, base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_add(base.Simd_g_i32x4_shr_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v81), base.Simd_g_i32x4_extend_low_i16x8_s(v44)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v81), base.Simd_g_i32x4_extend_high_i16x8_s(v44)), base.Simd_g_const(&F_TransformColorInverse_SSE2__k5), base.Simd_g_const(&F_TransformColorInverse_SSE2__k6)), v80), v81), v80))
			base.Simd_g_v128_store(m, v59, v64, v93)
			v96 = int32(16)
			v103 = v58 + int32(4)
			if v58+v80 <= l2 {
				v58 = v103
				v59 = v59 + v96
				v60 = v60 + v96
				continue
			} else {
				break
			}
			break
		}
		v109 = v103
	}
	if l2 == v109 {
	} else {
		v123 = v109 << (uint(int32(2)) % 32)
		v124 = l1 + v123
		v125 = l2 - v109
		v126 = l3 + v123
		if v125 < int32(1) {
		} else {
			v141 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+2)))
			v142 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+1)))
			v143 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0))))
			v144 = int32(0)
			if base.Ui32(v125) < base.Ui32(int32(4)) {
				v221 = v144
				v232 = v221 << (uint(int32(2)) % 32)
				v237 = v125 - v221
				v246 = v126 + v232
				v247 = v124 + v232
				for {
					v251 = *(*int32)(unsafe.Add(mBase, uint32(v247)))
					v252 = int32(16)
					v255 = v251 << (uint(v252) % 32) >> (uint(int32(24)) % 32)
					v257 = int32(5)
					v261 = v255*v143>>(uint(v257)%32) + int32(base.Ui32(v251)>>(uint(v252)%32))
					*(*int32)(unsafe.Add(mBase, uint32(v246))) = v261<<(uint(v252)%32)&int32(16711680) | v251&int32(-16711936) | (int32(base.Ui32(v255*v142)>>(uint(v257)%32))+v251+int32(base.Ui32(base.I32_extend8_s(v261)*v141)>>(uint(v257)%32)))&int32(255)
					v282 = int32(4)
					v287 = v237 + int32(-1)
					if v287 != 0 {
						v237 = v287
						v246 = v246 + v282
						v247 = v247 + v282
						continue
					} else {
						break
					}
					break
				}
			} else {
				if base.Ui32(v126-v124) < base.Ui32(int32(16)) {
					v221 = v144
					v232 = v221 << (uint(int32(2)) % 32)
					v237 = v125 - v221
					v246 = v126 + v232
					v247 = v124 + v232
					for {
						v251 = *(*int32)(unsafe.Add(mBase, uint32(v247)))
						v252 = int32(16)
						v255 = v251 << (uint(v252) % 32) >> (uint(int32(24)) % 32)
						v257 = int32(5)
						v261 = v255*v143>>(uint(v257)%32) + int32(base.Ui32(v251)>>(uint(v252)%32))
						*(*int32)(unsafe.Add(mBase, uint32(v246))) = v261<<(uint(v252)%32)&int32(16711680) | v251&int32(-16711936) | (int32(base.Ui32(v255*v142)>>(uint(v257)%32))+v251+int32(base.Ui32(base.I32_extend8_s(v261)*v141)>>(uint(v257)%32)))&int32(255)
						v282 = int32(4)
						v287 = v237 + int32(-1)
						if v287 != 0 {
							v237 = v287
							v246 = v246 + v282
							v247 = v247 + v282
							continue
						} else {
							break
						}
						break
					}
				} else {
					v154 = v125 & int32(2147483644)
					v155 = v126
					v166 = v154
					v167 = v124
					for {
						v171 = int32(0)
						v172 = base.Simd_g_v128_load(m, v167, v171)
						v173 = int32(16)
						v175 = int32(24)
						v176 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_shl(v172, v173), v175)
						v178 = int32(5)
						v182 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_mul(v176, base.Simd_g_i32x4_splat(v143)), v178), base.Simd_g_i32x4_shr_u(v172, v173))
						v204 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(v182, v173), base.Simd_g_const(&F_TransformColorInverse_SSE2__k7)), base.Simd_g_v128_and(v172, base.Simd_g_const(&F_TransformColorInverse_SSE2__k0))), base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v176, base.Simd_g_i32x4_splat(v142)), v178), v172), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_shl(v182, v175), v175), base.Simd_g_i32x4_splat(v141)), v178)), base.Simd_g_const(&F_TransformColorInverse_SSE2__k8)))
						base.Simd_g_v128_store(m, v155, v171, v204)
						v212 = v166 + int32(-4)
						if v212 != 0 {
							v155 = v155 + v173
							v166 = v212
							v167 = v167 + v173
							continue
						} else {
							break
						}
						break
					}
					if v154 == v125 {
					} else {
						v221 = v154
						v232 = v221 << (uint(int32(2)) % 32)
						v237 = v125 - v221
						v246 = v126 + v232
						v247 = v124 + v232
						for {
							v251 = *(*int32)(unsafe.Add(mBase, uint32(v247)))
							v252 = int32(16)
							v255 = v251 << (uint(v252) % 32) >> (uint(int32(24)) % 32)
							v257 = int32(5)
							v261 = v255*v143>>(uint(v257)%32) + int32(base.Ui32(v251)>>(uint(v252)%32))
							*(*int32)(unsafe.Add(mBase, uint32(v246))) = v261<<(uint(v252)%32)&int32(16711680) | v251&int32(-16711936) | (int32(base.Ui32(v255*v142)>>(uint(v257)%32))+v251+int32(base.Ui32(base.I32_extend8_s(v261)*v141)>>(uint(v257)%32)))&int32(255)
							v282 = int32(4)
							v287 = v237 + int32(-1)
							if v287 != 0 {
								v237 = v287
								v246 = v246 + v282
								v247 = v247 + v282
								continue
							} else {
								break
							}
							break
						}
					}
				}
			}
		}
	}
	return
}

var F_TransformColorInverse_SSE2__k0 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_TransformColorInverse_SSE2__k1 = [2]uint64{0x0, 0x0}
var F_TransformColorInverse_SSE2__k2 = [2]uint64{0x504050401000100, 0xf0e0d0c0b0a0908}
var F_TransformColorInverse_SSE2__k3 = [2]uint64{0x706050403020100, 0xd0c0d0c09080908}
var F_TransformColorInverse_SSE2__k4 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_TransformColorInverse_SSE2__k5 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_TransformColorInverse_SSE2__k6 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_TransformColorInverse_SSE2__k7 = [2]uint64{0xff000000ff0000, 0xff000000ff0000}
var F_TransformColorInverse_SSE2__k8 = [2]uint64{0xff000000ff, 0xff000000ff}

func F_TransformColorInverse_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v25 int32
	_ = v25
	var v29 int32
	_ = v29
	var v34 base.V128
	_ = v34
	var v37 int32
	_ = v37
	var v43 base.V128
	_ = v43
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v62 int32
	_ = v62
	var v63 base.V128
	_ = v63
	var v65 base.V128
	_ = v65
	var v70 base.V128
	_ = v70
	var v72 base.V128
	_ = v72
	var v74 base.V128
	_ = v74
	var v83 base.V128
	_ = v83
	var v86 int32
	_ = v86
	var v93 int32
	_ = v93
	var v99 int32
	_ = v99
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v132 int32
	_ = v132
	var v133 int32
	_ = v133
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v155 int32
	_ = v155
	var v156 int32
	_ = v156
	var v160 int32
	_ = v160
	var v161 base.V128
	_ = v161
	var v162 int32
	_ = v162
	var v164 int32
	_ = v164
	var v165 base.V128
	_ = v165
	var v167 int32
	_ = v167
	var v171 base.V128
	_ = v171
	var v193 base.V128
	_ = v193
	var v201 int32
	_ = v201
	var v210 int32
	_ = v210
	var v221 int32
	_ = v221
	var v226 int32
	_ = v226
	var v235 int32
	_ = v235
	var v236 int32
	_ = v236
	var v240 int32
	_ = v240
	var v241 int32
	_ = v241
	var v244 int32
	_ = v244
	var v246 int32
	_ = v246
	var v250 int32
	_ = v250
	var v271 int32
	_ = v271
	var v276 int32
	_ = v276
	if l2 < int32(4) {
		v99 = int32(0)
	} else {
		v19 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
		v20 = int32(8)
		v25 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)))
		v29 = int32(5)
		v34 = base.Simd_g_i32x4_splat(base.I32_extend16_s(v19<<(uint(v20)%32))<<(uint(int32(11))%32) | int32(base.Ui32(base.I32_extend16_s(v25<<(uint(v20)%32)))>>(uint(v29)%32))&int32(_a_F_TransformColorInverse_SSE41_0))
		v37 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)))
		v43 = base.Simd_g_i32x4_splat(base.I32_extend16_s(v37<<(uint(v20)%32)) >> (uint(v29) % 32))
		v57 = int32(0)
		v58 = l3
		v59 = l1
		for {
			v62 = int32(0)
			v63 = base.Simd_g_v128_load(m, v59, v62)
			v65 = base.Simd_g_i8x16_swizzle(v63, base.Simd_g_const(&F_TransformColorInverse_SSE41__k0))
			v70 = base.Simd_g_const(&F_TransformColorInverse_SSE41__k1)
			v72 = base.Simd_g_i8x16_add(v63, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v65), base.Simd_g_i32x4_extend_low_i16x8_s(v34)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v65), base.Simd_g_i32x4_extend_high_i16x8_s(v34)), base.Simd_g_const(&F_TransformColorInverse_SSE41__k2), base.Simd_g_const(&F_TransformColorInverse_SSE41__k3)))
			v74 = base.Simd_g_i8x16_swizzle(v72, base.Simd_g_const(&F_TransformColorInverse_SSE41__k4))
			v83 = base.Simd_g_v128_bitselect(v63, base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v74), base.Simd_g_i32x4_extend_low_i16x8_s(v43)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v74), base.Simd_g_i32x4_extend_high_i16x8_s(v43)), base.Simd_g_const(&F_TransformColorInverse_SSE41__k2), base.Simd_g_const(&F_TransformColorInverse_SSE41__k3)), v72), base.Simd_g_const(&F_TransformColorInverse_SSE41__k5))
			base.Simd_g_v128_store(m, v58, v62, v83)
			v86 = int32(16)
			v93 = v57 + int32(4)
			if v57+int32(8) <= l2 {
				v57 = v93
				v58 = v58 + v86
				v59 = v59 + v86
				continue
			} else {
				break
			}
			break
		}
		v99 = v93
	}
	if l2 == v99 {
	} else {
		v112 = v99 << (uint(int32(2)) % 32)
		v113 = l1 + v112
		v114 = l2 - v99
		v115 = l3 + v112
		if v114 < int32(1) {
		} else {
			v130 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+2)))
			v131 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+1)))
			v132 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0))))
			v133 = int32(0)
			if base.Ui32(v114) < base.Ui32(int32(4)) {
				v210 = v133
				v221 = v210 << (uint(int32(2)) % 32)
				v226 = v114 - v210
				v235 = v115 + v221
				v236 = v113 + v221
				for {
					v240 = *(*int32)(unsafe.Add(mBase, uint32(v236)))
					v241 = int32(16)
					v244 = v240 << (uint(v241) % 32) >> (uint(int32(24)) % 32)
					v246 = int32(5)
					v250 = v244*v132>>(uint(v246)%32) + int32(base.Ui32(v240)>>(uint(v241)%32))
					*(*int32)(unsafe.Add(mBase, uint32(v235))) = v250<<(uint(v241)%32)&int32(16711680) | v240&int32(-16711936) | (int32(base.Ui32(v244*v131)>>(uint(v246)%32))+v240+int32(base.Ui32(base.I32_extend8_s(v250)*v130)>>(uint(v246)%32)))&int32(255)
					v271 = int32(4)
					v276 = v226 + int32(-1)
					if v276 != 0 {
						v226 = v276
						v235 = v235 + v271
						v236 = v236 + v271
						continue
					} else {
						break
					}
					break
				}
			} else {
				if base.Ui32(v115-v113) < base.Ui32(int32(16)) {
					v210 = v133
					v221 = v210 << (uint(int32(2)) % 32)
					v226 = v114 - v210
					v235 = v115 + v221
					v236 = v113 + v221
					for {
						v240 = *(*int32)(unsafe.Add(mBase, uint32(v236)))
						v241 = int32(16)
						v244 = v240 << (uint(v241) % 32) >> (uint(int32(24)) % 32)
						v246 = int32(5)
						v250 = v244*v132>>(uint(v246)%32) + int32(base.Ui32(v240)>>(uint(v241)%32))
						*(*int32)(unsafe.Add(mBase, uint32(v235))) = v250<<(uint(v241)%32)&int32(16711680) | v240&int32(-16711936) | (int32(base.Ui32(v244*v131)>>(uint(v246)%32))+v240+int32(base.Ui32(base.I32_extend8_s(v250)*v130)>>(uint(v246)%32)))&int32(255)
						v271 = int32(4)
						v276 = v226 + int32(-1)
						if v276 != 0 {
							v226 = v276
							v235 = v235 + v271
							v236 = v236 + v271
							continue
						} else {
							break
						}
						break
					}
				} else {
					v143 = v114 & int32(2147483644)
					v144 = v115
					v155 = v143
					v156 = v113
					for {
						v160 = int32(0)
						v161 = base.Simd_g_v128_load(m, v156, v160)
						v162 = int32(16)
						v164 = int32(24)
						v165 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_shl(v161, v162), v164)
						v167 = int32(5)
						v171 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_mul(v165, base.Simd_g_i32x4_splat(v132)), v167), base.Simd_g_i32x4_shr_u(v161, v162))
						v193 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(v171, v162), base.Simd_g_const(&F_TransformColorInverse_SSE41__k6)), base.Simd_g_v128_and(v161, base.Simd_g_const(&F_TransformColorInverse_SSE41__k5))), base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v165, base.Simd_g_i32x4_splat(v131)), v167), v161), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_shl(v171, v164), v164), base.Simd_g_i32x4_splat(v130)), v167)), base.Simd_g_const(&F_TransformColorInverse_SSE41__k7)))
						base.Simd_g_v128_store(m, v144, v160, v193)
						v201 = v155 + int32(-4)
						if v201 != 0 {
							v144 = v144 + v162
							v155 = v201
							v156 = v156 + v162
							continue
						} else {
							break
						}
						break
					}
					if v143 == v114 {
					} else {
						v210 = v143
						v221 = v210 << (uint(int32(2)) % 32)
						v226 = v114 - v210
						v235 = v115 + v221
						v236 = v113 + v221
						for {
							v240 = *(*int32)(unsafe.Add(mBase, uint32(v236)))
							v241 = int32(16)
							v244 = v240 << (uint(v241) % 32) >> (uint(int32(24)) % 32)
							v246 = int32(5)
							v250 = v244*v132>>(uint(v246)%32) + int32(base.Ui32(v240)>>(uint(v241)%32))
							*(*int32)(unsafe.Add(mBase, uint32(v235))) = v250<<(uint(v241)%32)&int32(16711680) | v240&int32(-16711936) | (int32(base.Ui32(v244*v131)>>(uint(v246)%32))+v240+int32(base.Ui32(base.I32_extend8_s(v250)*v130)>>(uint(v246)%32)))&int32(255)
							v271 = int32(4)
							v276 = v226 + int32(-1)
							if v276 != 0 {
								v226 = v276
								v235 = v235 + v271
								v236 = v236 + v271
								continue
							} else {
								break
							}
							break
						}
					}
				}
			}
		}
	}
	return
}

var F_TransformColorInverse_SSE41__k0 = [2]uint64{0x58f058f018f018f, 0xd8f0d8f098f098f}
var F_TransformColorInverse_SSE41__k1 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_TransformColorInverse_SSE41__k2 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_TransformColorInverse_SSE41__k3 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_TransformColorInverse_SSE41__k4 = [2]uint64{0x8f8f068f8f8f028f, 0x8f8f0e8f8f8f0a8f}
var F_TransformColorInverse_SSE41__k5 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_TransformColorInverse_SSE41__k6 = [2]uint64{0xff000000ff0000, 0xff000000ff0000}
var F_TransformColorInverse_SSE41__k7 = [2]uint64{0xff000000ff, 0xff000000ff}

func F_TransformColor_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	var v32 base.V128
	_ = v32
	var v35 int32
	_ = v35
	var v41 base.V128
	_ = v41
	var v54 int32
	_ = v54
	var v55 int32
	_ = v55
	var v58 int32
	_ = v58
	var v59 base.V128
	_ = v59
	var v62 base.V128
	_ = v62
	var v66 base.V128
	_ = v66
	var v71 base.V128
	_ = v71
	var v73 int32
	_ = v73
	var v74 base.V128
	_ = v74
	var v81 int32
	_ = v81
	var v86 base.V128
	_ = v86
	var v94 int32
	_ = v94
	var v99 int32
	_ = v99
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v127 int32
	_ = v127
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v137 int32
	_ = v137
	var v138 int32
	_ = v138
	var v148 int32
	_ = v148
	var v152 int32
	_ = v152
	var v153 base.V128
	_ = v153
	var v154 int32
	_ = v154
	var v155 base.V128
	_ = v155
	var v158 int32
	_ = v158
	var v159 base.V128
	_ = v159
	var v161 int32
	_ = v161
	var v185 base.V128
	_ = v185
	var v191 int32
	_ = v191
	var v199 int32
	_ = v199
	var v211 int32
	_ = v211
	var v213 int32
	_ = v213
	var v225 int32
	_ = v225
	var v226 int32
	_ = v226
	var v227 int32
	_ = v227
	var v231 int32
	_ = v231
	var v233 int32
	_ = v233
	var v259 int32
	_ = v259
	if l2 < int32(4) {
		v99 = int32(0)
	} else {
		v17 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
		v18 = int32(8)
		v21 = int32(11)
		v23 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)))
		v32 = base.Simd_g_i32x4_splat(base.I32_extend16_s(v17<<(uint(v18)%32))<<(uint(v21)%32) | int32(base.Ui32(base.I32_extend16_s(v23<<(uint(v18)%32)))>>(uint(int32(5))%32))&int32(_a_F_TransformColor_SSE2_0))
		v35 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)))
		v41 = base.Simd_g_i32x4_splat(base.I32_extend16_s(v35<<(uint(v18)%32)) << (uint(v21) % 32))
		v54 = int32(0)
		v55 = l1
		for {
			v58 = int32(0)
			v59 = base.Simd_g_v128_load(m, v55, v58)
			v62 = base.Simd_g_const(&F_TransformColor_SSE2__k0)
			v66 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_v128_and(v59, base.Simd_g_const(&F_TransformColor_SSE2__k1)), base.Simd_g_const(&F_TransformColor_SSE2__k2)), base.Simd_g_const(&F_TransformColor_SSE2__k3))
			v71 = base.Simd_g_const(&F_TransformColor_SSE2__k4)
			v73 = int32(8)
			v74 = base.Simd_g_i16x8_shl(v59, v73)
			v81 = int32(16)
			v86 = base.Simd_g_i8x16_sub(v59, base.Simd_g_v128_and(base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v66), base.Simd_g_i32x4_extend_low_i16x8_s(v32)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v66), base.Simd_g_i32x4_extend_high_i16x8_s(v32)), base.Simd_g_const(&F_TransformColor_SSE2__k5), base.Simd_g_const(&F_TransformColor_SSE2__k6)), base.Simd_g_i32x4_shr_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v74), base.Simd_g_i32x4_extend_low_i16x8_s(v41)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v74), base.Simd_g_i32x4_extend_high_i16x8_s(v41)), base.Simd_g_const(&F_TransformColor_SSE2__k5), base.Simd_g_const(&F_TransformColor_SSE2__k6)), v81)), base.Simd_g_const(&F_TransformColor_SSE2__k7)))
			base.Simd_g_v128_store(m, v55, v58, v86)
			v94 = v54 + int32(4)
			if v54+v73 <= l2 {
				v54 = v94
				v55 = v55 + v81
				continue
			} else {
				break
			}
			break
		}
		v99 = v94
	}
	if l2 == v99 {
	} else {
		v112 = l1 + v99<<(uint(int32(2))%32)
		v113 = l2 - v99
		if v113 < int32(1) {
		} else {
			v127 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+2)))
			v128 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+1)))
			v129 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0))))
			if base.Ui32(v113) < base.Ui32(int32(4)) {
				v199 = int32(0)
				v211 = v112 + v199<<(uint(int32(2))%32)
				v213 = v113 - v199
				for {
					v225 = *(*int32)(unsafe.Add(mBase, uint32(v211)))
					v226 = int32(16)
					v227 = int32(base.Ui32(v225) >> (uint(v226) % 32))
					v231 = v225 << (uint(v226) % 32) >> (uint(int32(24)) % 32)
					v233 = int32(5)
					*(*int32)(unsafe.Add(mBase, uint32(v211))) = (v227-int32(base.Ui32(v231*v129)>>(uint(v233)%32)))<<(uint(v226)%32)&int32(16711680) | v225&int32(-16711936) | (v225-(int32(base.Ui32(v231*v128)>>(uint(v233)%32))+int32(base.Ui32(base.I32_extend8_s(v227)*v127)>>(uint(v233)%32))))&int32(255)
					v259 = v213 + int32(-1)
					if v259 != 0 {
						v211 = v211 + int32(4)
						v213 = v259
						continue
					} else {
						break
					}
					break
				}
			} else {
				v137 = v113 & int32(2147483644)
				v138 = v112
				v148 = v137
				for {
					v152 = int32(0)
					v153 = base.Simd_g_v128_load(m, v138, v152)
					v154 = int32(16)
					v155 = base.Simd_g_i32x4_shr_u(v153, v154)
					v158 = int32(24)
					v159 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_shl(v153, v154), v158)
					v161 = int32(5)
					v185 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(base.Simd_g_i32x4_sub(v155, base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v159, base.Simd_g_i32x4_splat(v129)), v161)), v154), base.Simd_g_const(&F_TransformColor_SSE2__k8)), base.Simd_g_v128_and(v153, base.Simd_g_const(&F_TransformColor_SSE2__k1))), base.Simd_g_v128_and(base.Simd_g_i32x4_sub(v153, base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v159, base.Simd_g_i32x4_splat(v128)), v161), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_shl(v155, v158), v158), base.Simd_g_i32x4_splat(v127)), v161))), base.Simd_g_const(&F_TransformColor_SSE2__k9)))
					base.Simd_g_v128_store(m, v138, v152, v185)
					v191 = v148 + int32(-4)
					if v191 != 0 {
						v138 = v138 + v154
						v148 = v191
						continue
					} else {
						break
					}
					break
				}
				if v137 == v113 {
				} else {
					v199 = v137
					v211 = v112 + v199<<(uint(int32(2))%32)
					v213 = v113 - v199
					for {
						v225 = *(*int32)(unsafe.Add(mBase, uint32(v211)))
						v226 = int32(16)
						v227 = int32(base.Ui32(v225) >> (uint(v226) % 32))
						v231 = v225 << (uint(v226) % 32) >> (uint(int32(24)) % 32)
						v233 = int32(5)
						*(*int32)(unsafe.Add(mBase, uint32(v211))) = (v227-int32(base.Ui32(v231*v129)>>(uint(v233)%32)))<<(uint(v226)%32)&int32(16711680) | v225&int32(-16711936) | (v225-(int32(base.Ui32(v231*v128)>>(uint(v233)%32))+int32(base.Ui32(base.I32_extend8_s(v227)*v127)>>(uint(v233)%32))))&int32(255)
						v259 = v213 + int32(-1)
						if v259 != 0 {
							v211 = v211 + int32(4)
							v213 = v259
							continue
						} else {
							break
						}
						break
					}
				}
			}
		}
	}
	return
}

var F_TransformColor_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_TransformColor_SSE2__k1 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_TransformColor_SSE2__k2 = [2]uint64{0x504050401000100, 0xf0e0d0c0b0a0908}
var F_TransformColor_SSE2__k3 = [2]uint64{0x706050403020100, 0xd0c0d0c09080908}
var F_TransformColor_SSE2__k4 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_TransformColor_SSE2__k5 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_TransformColor_SSE2__k6 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_TransformColor_SSE2__k7 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_TransformColor_SSE2__k8 = [2]uint64{0xff000000ff0000, 0xff000000ff0000}
var F_TransformColor_SSE2__k9 = [2]uint64{0xff000000ff, 0xff000000ff}

func F_Transform_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	var v18 base.V128
	_ = v18
	var v20 base.V128
	_ = v20
	var v22 base.V128
	_ = v22
	var v23 int32
	_ = v23
	var v24 base.V128
	_ = v24
	var v28 base.V128
	_ = v28
	var v29 base.V128
	_ = v29
	var v32 base.V128
	_ = v32
	var v36 base.V128
	_ = v36
	var v40 base.V128
	_ = v40
	var v43 base.V128
	_ = v43
	var v44 base.V128
	_ = v44
	var v45 base.V128
	_ = v45
	var v46 base.V128
	_ = v46
	var v48 base.V128
	_ = v48
	var v49 base.V128
	_ = v49
	var v51 base.V128
	_ = v51
	var v53 base.V128
	_ = v53
	var v56 base.V128
	_ = v56
	var v57 base.V128
	_ = v57
	var v59 base.V128
	_ = v59
	var v63 base.V128
	_ = v63
	var v64 base.V128
	_ = v64
	var v65 base.V128
	_ = v65
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
	var v82 base.V128
	_ = v82
	var v84 base.V128
	_ = v84
	var v85 base.V128
	_ = v85
	var v86 base.V128
	_ = v86
	var v87 base.V128
	_ = v87
	var v88 base.V128
	_ = v88
	var v90 base.V128
	_ = v90
	var v92 base.V128
	_ = v92
	var v93 base.V128
	_ = v93
	var v94 base.V128
	_ = v94
	var v95 base.V128
	_ = v95
	var v96 base.V128
	_ = v96
	var v98 base.V128
	_ = v98
	var v100 base.V128
	_ = v100
	var v102 base.V128
	_ = v102
	var v104 base.V128
	_ = v104
	var v109 base.V128
	_ = v109
	var v111 base.V128
	_ = v111
	var v115 base.V128
	_ = v115
	var v116 base.V128
	_ = v116
	var v119 base.V128
	_ = v119
	var v121 base.V128
	_ = v121
	var v122 base.V128
	_ = v122
	var v124 int32
	_ = v124
	var v125 base.V128
	_ = v125
	var v136 base.V128
	_ = v136
	var v137 base.V128
	_ = v137
	var v140 base.V128
	_ = v140
	var v142 base.V128
	_ = v142
	var v145 base.V128
	_ = v145
	var v148 base.V128
	_ = v148
	var v150 base.V128
	_ = v150
	var v152 base.V128
	_ = v152
	var v154 base.V128
	_ = v154
	var v156 base.V128
	_ = v156
	var v158 base.V128
	_ = v158
	var v164 base.V128
	_ = v164
	var v166 base.V128
	_ = v166
	var v174 base.V128
	_ = v174
	var v176 base.V128
	_ = v176
	var v178 base.V128
	_ = v178
	var v180 base.V128
	_ = v180
	var v182 base.V128
	_ = v182
	var v184 base.V128
	_ = v184
	var v186 base.V128
	_ = v186
	var v188 base.V128
	_ = v188
	var v189 base.V128
	_ = v189
	var v190 base.V128
	_ = v190
	var v191 base.V128
	_ = v191
	var v192 base.V128
	_ = v192
	var v193 base.V128
	_ = v193
	var v194 base.V128
	_ = v194
	var v196 base.V128
	_ = v196
	var v197 base.V128
	_ = v197
	var v200 base.V128
	_ = v200
	var v201 base.V128
	_ = v201
	var v204 base.V128
	_ = v204
	var v205 base.V128
	_ = v205
	var v208 base.V128
	_ = v208
	var v209 base.V128
	_ = v209
	var v213 int32
	_ = v213
	var v225 int32
	_ = v225
	v18 = base.Simd_g_v128_load64_zero(m, l0, int32(24))
	v20 = base.Simd_g_v128_load64_zero(m, l0, int32(16))
	v22 = base.Simd_g_v128_load64_zero(m, l0, int32(8))
	v23 = int32(0)
	v24 = base.Simd_g_v128_load64_zero(m, l0, v23)
	if l2 == v23 {
		v43 = v18
		v44 = v20
		v45 = v22
		v46 = v24
	} else {
		v28 = base.Simd_g_v128_load64_zero(m, l0, int32(56))
		v29 = base.Simd_g_const(&F_Transform_SSE2__k0)
		v32 = base.Simd_g_v128_load64_zero(m, l0, int32(48))
		v36 = base.Simd_g_v128_load64_zero(m, l0, int32(40))
		v40 = base.Simd_g_v128_load64_zero(m, l0, int32(32))
		v43 = base.Simd_g_i8x16_shuffle2(v18, v28, base.Simd_g_const(&F_Transform_SSE2__k1), base.Simd_g_const(&F_Transform_SSE2__k2))
		v44 = base.Simd_g_i8x16_shuffle2(v20, v32, base.Simd_g_const(&F_Transform_SSE2__k1), base.Simd_g_const(&F_Transform_SSE2__k2))
		v45 = base.Simd_g_i8x16_shuffle2(v22, v36, base.Simd_g_const(&F_Transform_SSE2__k1), base.Simd_g_const(&F_Transform_SSE2__k2))
		v46 = base.Simd_g_i8x16_shuffle2(v24, v40, base.Simd_g_const(&F_Transform_SSE2__k1), base.Simd_g_const(&F_Transform_SSE2__k2))
	}
	v48 = base.Simd_g_i32x4_extend_low_i16x8_s(v43)
	v49 = base.Simd_g_const(&F_Transform_SSE2__k3)
	v51 = base.Simd_g_i32x4_extend_high_i16x8_s(v43)
	v53 = base.Simd_g_const(&F_Transform_SSE2__k4)
	v56 = base.Simd_g_i32x4_extend_low_i16x8_s(v45)
	v57 = base.Simd_g_const(&F_Transform_SSE2__k5)
	v59 = base.Simd_g_i32x4_extend_high_i16x8_s(v45)
	v63 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v45, v43), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v48, v49), base.Simd_g_i32x4_mul(v51, v49), base.Simd_g_const(&F_Transform_SSE2__k6), base.Simd_g_const(&F_Transform_SSE2__k7))), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v56, v57), base.Simd_g_i32x4_mul(v59, v57), base.Simd_g_const(&F_Transform_SSE2__k6), base.Simd_g_const(&F_Transform_SSE2__k7)))
	v64 = base.Simd_g_i16x8_add(v46, v44)
	v65 = base.Simd_g_i16x8_add(v63, v64)
	v76 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v56, v49), base.Simd_g_i32x4_mul(v59, v49), base.Simd_g_const(&F_Transform_SSE2__k6), base.Simd_g_const(&F_Transform_SSE2__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v48, v57), base.Simd_g_i32x4_mul(v51, v57), base.Simd_g_const(&F_Transform_SSE2__k6), base.Simd_g_const(&F_Transform_SSE2__k7))), base.Simd_g_i16x8_sub(v45, v43))
	v77 = base.Simd_g_i16x8_sub(v46, v44)
	v78 = base.Simd_g_i16x8_add(v76, v77)
	v79 = base.Simd_g_const(&F_Transform_SSE2__k8)
	v80 = base.Simd_g_i8x16_shuffle2(v65, v78, base.Simd_g_const(&F_Transform_SSE2__k9), base.Simd_g_const(&F_Transform_SSE2__k10))
	v81 = base.Simd_g_i16x8_sub(v77, v76)
	v82 = base.Simd_g_i16x8_sub(v64, v63)
	v84 = base.Simd_g_i8x16_shuffle2(v81, v82, base.Simd_g_const(&F_Transform_SSE2__k9), base.Simd_g_const(&F_Transform_SSE2__k10))
	v85 = base.Simd_g_const(&F_Transform_SSE2__k11)
	v86 = base.Simd_g_i8x16_shuffle2(v80, v84, base.Simd_g_const(&F_Transform_SSE2__k12), base.Simd_g_const(&F_Transform_SSE2__k13))
	v87 = base.Simd_g_const(&F_Transform_SSE2__k14)
	v88 = base.Simd_g_i8x16_shuffle2(v65, v78, base.Simd_g_const(&F_Transform_SSE2__k15), base.Simd_g_const(&F_Transform_SSE2__k16))
	v90 = base.Simd_g_i8x16_shuffle2(v81, v82, base.Simd_g_const(&F_Transform_SSE2__k15), base.Simd_g_const(&F_Transform_SSE2__k16))
	v92 = base.Simd_g_i8x16_shuffle2(v88, v90, base.Simd_g_const(&F_Transform_SSE2__k12), base.Simd_g_const(&F_Transform_SSE2__k13))
	v93 = base.Simd_g_const(&F_Transform_SSE2__k17)
	v94 = base.Simd_g_i8x16_shuffle2(v86, v92, base.Simd_g_const(&F_Transform_SSE2__k18), base.Simd_g_const(&F_Transform_SSE2__k19))
	v95 = base.Simd_g_const(&F_Transform_SSE2__k20)
	v96 = base.Simd_g_i8x16_shuffle2(v80, v84, base.Simd_g_const(&F_Transform_SSE2__k21), base.Simd_g_const(&F_Transform_SSE2__k22))
	v98 = base.Simd_g_i8x16_shuffle2(v88, v90, base.Simd_g_const(&F_Transform_SSE2__k21), base.Simd_g_const(&F_Transform_SSE2__k22))
	v100 = base.Simd_g_i8x16_shuffle2(v96, v98, base.Simd_g_const(&F_Transform_SSE2__k18), base.Simd_g_const(&F_Transform_SSE2__k19))
	v102 = base.Simd_g_i32x4_extend_low_i16x8_s(v100)
	v104 = base.Simd_g_i32x4_extend_high_i16x8_s(v100)
	v109 = base.Simd_g_i32x4_extend_low_i16x8_s(v94)
	v111 = base.Simd_g_i32x4_extend_high_i16x8_s(v94)
	v115 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v94, v100), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v102, v49), base.Simd_g_i32x4_mul(v104, v49), base.Simd_g_const(&F_Transform_SSE2__k6), base.Simd_g_const(&F_Transform_SSE2__k7))), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v109, v57), base.Simd_g_i32x4_mul(v111, v57), base.Simd_g_const(&F_Transform_SSE2__k6), base.Simd_g_const(&F_Transform_SSE2__k7)))
	v116 = base.Simd_g_const(&F_Transform_SSE2__k0)
	v119 = base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v86, v92, base.Simd_g_const(&F_Transform_SSE2__k1), base.Simd_g_const(&F_Transform_SSE2__k2)), base.Simd_g_const(&F_Transform_SSE2__k23))
	v121 = base.Simd_g_i8x16_shuffle2(v96, v98, base.Simd_g_const(&F_Transform_SSE2__k1), base.Simd_g_const(&F_Transform_SSE2__k2))
	v122 = base.Simd_g_i16x8_add(v119, v121)
	v124 = int32(3)
	v125 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(v115, v122), v124)
	v136 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v109, v49), base.Simd_g_i32x4_mul(v111, v49), base.Simd_g_const(&F_Transform_SSE2__k6), base.Simd_g_const(&F_Transform_SSE2__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v102, v57), base.Simd_g_i32x4_mul(v104, v57), base.Simd_g_const(&F_Transform_SSE2__k6), base.Simd_g_const(&F_Transform_SSE2__k7))), base.Simd_g_i16x8_sub(v94, v100))
	v137 = base.Simd_g_i16x8_sub(v119, v121)
	v140 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(v136, v137), v124)
	v142 = base.Simd_g_i8x16_shuffle2(v125, v140, base.Simd_g_const(&F_Transform_SSE2__k9), base.Simd_g_const(&F_Transform_SSE2__k10))
	v145 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_sub(v137, v136), v124)
	v148 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_sub(v122, v115), v124)
	v150 = base.Simd_g_i8x16_shuffle2(v145, v148, base.Simd_g_const(&F_Transform_SSE2__k9), base.Simd_g_const(&F_Transform_SSE2__k10))
	v152 = base.Simd_g_i8x16_shuffle2(v142, v150, base.Simd_g_const(&F_Transform_SSE2__k21), base.Simd_g_const(&F_Transform_SSE2__k22))
	v154 = base.Simd_g_i8x16_shuffle2(v125, v140, base.Simd_g_const(&F_Transform_SSE2__k15), base.Simd_g_const(&F_Transform_SSE2__k16))
	v156 = base.Simd_g_i8x16_shuffle2(v145, v148, base.Simd_g_const(&F_Transform_SSE2__k15), base.Simd_g_const(&F_Transform_SSE2__k16))
	v158 = base.Simd_g_i8x16_shuffle2(v154, v156, base.Simd_g_const(&F_Transform_SSE2__k21), base.Simd_g_const(&F_Transform_SSE2__k22))
	v164 = base.Simd_g_i8x16_shuffle2(v142, v150, base.Simd_g_const(&F_Transform_SSE2__k12), base.Simd_g_const(&F_Transform_SSE2__k13))
	v166 = base.Simd_g_i8x16_shuffle2(v154, v156, base.Simd_g_const(&F_Transform_SSE2__k12), base.Simd_g_const(&F_Transform_SSE2__k13))
	if l2 == int32(0) {
		v182 = base.Simd_g_v128_load32_zero(m, l1, int32(96))
		v184 = base.Simd_g_v128_load32_zero(m, l1, int32(64))
		v186 = base.Simd_g_v128_load32_zero(m, l1, int32(32))
		v188 = base.Simd_g_v128_load32_zero(m, l1, int32(0))
		v189 = v182
		v190 = v184
		v191 = v186
		v192 = v188
	} else {
		v174 = base.Simd_g_v128_load64_zero(m, l1, int32(96))
		v176 = base.Simd_g_v128_load64_zero(m, l1, int32(64))
		v178 = base.Simd_g_v128_load64_zero(m, l1, int32(32))
		v180 = base.Simd_g_v128_load64_zero(m, l1, int32(0))
		v189 = v174
		v190 = v176
		v191 = v178
		v192 = v180
	}
	v193 = base.Simd_g_const(&F_Transform_SSE2__k24)
	v194 = base.Simd_g_const(&F_Transform_SSE2__k25)
	v196 = base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v189, v193, base.Simd_g_const(&F_Transform_SSE2__k26), base.Simd_g_const(&F_Transform_SSE2__k27)), base.Simd_g_i8x16_shuffle2(v152, v158, base.Simd_g_const(&F_Transform_SSE2__k18), base.Simd_g_const(&F_Transform_SSE2__k19)))
	v197 = base.Simd_g_i8x16_narrow_i16x8_u(v196, v196)
	v200 = base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v190, v193, base.Simd_g_const(&F_Transform_SSE2__k26), base.Simd_g_const(&F_Transform_SSE2__k27)), base.Simd_g_i8x16_shuffle2(v152, v158, base.Simd_g_const(&F_Transform_SSE2__k1), base.Simd_g_const(&F_Transform_SSE2__k2)))
	v201 = base.Simd_g_i8x16_narrow_i16x8_u(v200, v200)
	v204 = base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v191, v193, base.Simd_g_const(&F_Transform_SSE2__k26), base.Simd_g_const(&F_Transform_SSE2__k27)), base.Simd_g_i8x16_shuffle2(v164, v166, base.Simd_g_const(&F_Transform_SSE2__k18), base.Simd_g_const(&F_Transform_SSE2__k19)))
	v205 = base.Simd_g_i8x16_narrow_i16x8_u(v204, v204)
	v208 = base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v192, v193, base.Simd_g_const(&F_Transform_SSE2__k26), base.Simd_g_const(&F_Transform_SSE2__k27)), base.Simd_g_i8x16_shuffle2(v164, v166, base.Simd_g_const(&F_Transform_SSE2__k1), base.Simd_g_const(&F_Transform_SSE2__k2)))
	v209 = base.Simd_g_i8x16_narrow_i16x8_u(v208, v208)
	if l2 == int32(0) {
		v225 = int32(0)
		base.Simd_g_v128_store32_lane_l0(m, l1, int32(96), v197)
		base.Simd_g_v128_store32_lane_l0(m, l1, int32(64), v201)
		base.Simd_g_v128_store32_lane_l0(m, l1, int32(32), v205)
		base.Simd_g_v128_store32_lane_l0(m, l1, v225, v209)
		return
	} else {
		v213 = int32(0)
		base.Simd_g_v128_store64_lane_l0(m, l1, int32(96), v197)
		base.Simd_g_v128_store64_lane_l0(m, l1, int32(64), v201)
		base.Simd_g_v128_store64_lane_l0(m, l1, int32(32), v205)
		base.Simd_g_v128_store64_lane_l0(m, l1, v213, v209)
		return
	}
}

var F_Transform_SSE2__k0 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_Transform_SSE2__k1 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_Transform_SSE2__k2 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_Transform_SSE2__k3 = [2]uint64{0xffff8a8cffff8a8c, 0xffff8a8cffff8a8c}
var F_Transform_SSE2__k4 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_Transform_SSE2__k5 = [2]uint64{0x4e7b00004e7b, 0x4e7b00004e7b}
var F_Transform_SSE2__k6 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_Transform_SSE2__k7 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_Transform_SSE2__k8 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_Transform_SSE2__k9 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_Transform_SSE2__k10 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_Transform_SSE2__k11 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_Transform_SSE2__k12 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_Transform_SSE2__k13 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_Transform_SSE2__k14 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_Transform_SSE2__k15 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_Transform_SSE2__k16 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_Transform_SSE2__k17 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_Transform_SSE2__k18 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_Transform_SSE2__k19 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_Transform_SSE2__k20 = [2]uint64{0x1b1a19180b0a0908, 0x1f1e1d1c0f0e0d0c}
var F_Transform_SSE2__k21 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_Transform_SSE2__k22 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_Transform_SSE2__k23 = [2]uint64{0x4000400040004, 0x4000400040004}
var F_Transform_SSE2__k24 = [2]uint64{0x0, 0x0}
var F_Transform_SSE2__k25 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_Transform_SSE2__k26 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_Transform_SSE2__k27 = [2]uint64{0x380028001800080, 0x780068005800480}

func F_TrellisQuantizeBlock(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32) int32 {
	mBase := m.M
	_ = mBase
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v43 int32
	_ = v43
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v60 int32
	_ = v60
	var v66 int32
	_ = v66
	var v86 int32
	_ = v86
	var v109 int32
	_ = v109
	var v113 int32
	_ = v113
	var v117 int32
	_ = v117
	var v123 int32
	_ = v123
	var v131 int32
	_ = v131
	var v134 int32
	_ = v134
	var v135 int32
	_ = v135
	var v139 int64
	_ = v139
	var v140 int64
	_ = v140
	var v142 int32
	_ = v142
	var v148 int64
	_ = v148
	var v150 int64
	_ = v150
	var v157 int32
	_ = v157
	var v160 int32
	_ = v160
	var v179 int32
	_ = v179
	var v182 int32
	_ = v182
	var __phi182 int32
	_ = __phi182
	var v190 int32
	_ = v190
	var __phi190 int32
	_ = __phi190
	var v193 int32
	_ = v193
	var __phi193 int32
	_ = __phi193
	var v196 int32
	_ = v196
	var __phi196 int32
	_ = __phi196
	var v197 int32
	_ = v197
	var __phi197 int32
	_ = __phi197
	var v198 int32
	_ = v198
	var __phi198 int32
	_ = __phi198
	var v199 int64
	_ = v199
	var __phi199 int64
	_ = __phi199
	var v201 int64
	_ = v201
	var __phi201 int64
	_ = __phi201
	var v202 int32
	_ = v202
	var __phi202 int32
	_ = __phi202
	var v203 int32
	_ = v203
	var __phi203 int32
	_ = __phi203
	var v218 int32
	_ = v218
	var v222 int32
	_ = v222
	var v224 int32
	_ = v224
	var v226 int32
	_ = v226
	var v228 int32
	_ = v228
	var v232 int32
	_ = v232
	var v233 int32
	_ = v233
	var v235 int32
	_ = v235
	var v236 int32
	_ = v236
	var v238 int32
	_ = v238
	var v239 int32
	_ = v239
	var v242 int32
	_ = v242
	var v246 int32
	_ = v246
	var v251 int32
	_ = v251
	var v252 int32
	_ = v252
	var v255 int32
	_ = v255
	var v258 int32
	_ = v258
	var v259 int32
	_ = v259
	var v260 int32
	_ = v260
	var v264 int32
	_ = v264
	var v265 int32
	_ = v265
	var v269 int32
	_ = v269
	var v271 int32
	_ = v271
	var v272 int32
	_ = v272
	var v275 int32
	_ = v275
	var v277 int64
	_ = v277
	var v286 int32
	_ = v286
	var v287 int32
	_ = v287
	var v288 int32
	_ = v288
	var v291 int32
	_ = v291
	var v292 int32
	_ = v292
	var v293 int32
	_ = v293
	var v295 int32
	_ = v295
	var v299 int32
	_ = v299
	var v303 int64
	_ = v303
	var v304 int64
	_ = v304
	var v305 int32
	_ = v305
	var v307 int32
	_ = v307
	var v311 int64
	_ = v311
	var v312 int32
	_ = v312
	var v314 int64
	_ = v314
	var v315 int32
	_ = v315
	var v318 int32
	_ = v318
	var v323 int64
	_ = v323
	var v337 int32
	_ = v337
	var v338 int32
	_ = v338
	var v342 int64
	_ = v342
	var v344 int64
	_ = v344
	var v346 int64
	_ = v346
	var v349 int32
	_ = v349
	var v351 int64
	_ = v351
	var v352 int64
	_ = v352
	var v353 int32
	_ = v353
	var v354 int32
	_ = v354
	var v359 int32
	_ = v359
	var v360 int32
	_ = v360
	var v363 int32
	_ = v363
	var v367 int32
	_ = v367
	var v376 int32
	_ = v376
	var v379 int32
	_ = v379
	var v380 int32
	_ = v380
	var v383 int32
	_ = v383
	var v384 int32
	_ = v384
	var v385 int32
	_ = v385
	var v387 int32
	_ = v387
	var v391 int32
	_ = v391
	var v395 int64
	_ = v395
	var v396 int64
	_ = v396
	var v397 int32
	_ = v397
	var v399 int32
	_ = v399
	var v403 int64
	_ = v403
	var v404 int64
	_ = v404
	var v405 int32
	_ = v405
	var v407 int64
	_ = v407
	var v408 int32
	_ = v408
	var v411 int32
	_ = v411
	var v416 int64
	_ = v416
	var v428 int32
	_ = v428
	var v429 int32
	_ = v429
	var v433 int64
	_ = v433
	var v435 int64
	_ = v435
	var v437 int64
	_ = v437
	var v444 int32
	_ = v444
	var v446 int64
	_ = v446
	var v447 int32
	_ = v447
	var v448 int32
	_ = v448
	var v456 int32
	_ = v456
	var v466 int32
	_ = v466
	var v478 int32
	_ = v478
	var v479 int32
	_ = v479
	var v494 base.V128
	_ = v494
	var v495 int32
	_ = v495
	var v497 int32
	_ = v497
	var v499 int64
	_ = v499
	var v501 int32
	_ = v501
	var v515 base.V128
	_ = v515
	var v516 int32
	_ = v516
	var v518 int32
	_ = v518
	var v529 int32
	_ = v529
	var v536 int32
	_ = v536
	var v549 int32
	_ = v549
	var v554 int32
	_ = v554
	var v559 int32
	_ = v559
	var v562 int32
	_ = v562
	var v567 int32
	_ = v567
	var v583 int32
	_ = v583
	var v585 int32
	_ = v585
	var v588 int32
	_ = v588
	var v592 int32
	_ = v592
	var v593 int32
	_ = v593
	var v595 int32
	_ = v595
	var v599 int32
	_ = v599
	var v601 int32
	_ = v601
	var v604 int32
	_ = v604
	var v605 int32
	_ = v605
	var v611 int32
	_ = v611
	var v613 int32
	_ = v613
	var v634 int32
	_ = v634
	v37 = m.G0
	v39 = v37 - int32(192)
	m.G0 = v39
	if l4 != 0 {
		goto L1
	} else {
		goto L2
	}
L1:
	;
	v43 = int32(-1)
	goto L3
L2:
	;
	v43 = int32(0)
	goto L3
L3:
	;
	v48 = l0 + l4*int32(264) + int32(3420)
	v49 = m.G81
	v51 = base.B2i32(l4 == int32(0))
	v53 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v49+v51))))
	v60 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48+v53*int32(33)+l3*int32(11)))))
	v66 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l5)+2)))
	v86 = int32(15)
	goto L5
L4:
	;
	v131 = *(*int32)(unsafe.Add(mBase, uint32(l0+l4*int32(192)+int32(_a_F_TrellisQuantizeBlock_0)+v51*int32(12)+l3<<(uint(int32(2))%32))))
	v134 = v123 + base.B2i32(v123 < int32(15))
	v135 = m.G79
	v139 = int64(*(*uint16)(unsafe.Add(mBase, uint32(v135+v60<<(uint(int32(1))%32)))))
	v140 = base.I64_extend_i32_s(l6)
	if l3 != 0 {
		v150 = int64(0)
		goto L9
	} else {
		goto L10
	}
L5:
	;
	v109 = m.G1
	v113 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v109+int32(_a_F_TrellisQuantizeBlock_1)+v86))))
	v117 = int32(*(*int16)(unsafe.Add(mBase, uint32(l1+v113<<(uint(int32(1))%32)))))
	if base.Ui32(int32(base.Ui32(v66*v66)>>(uint(int32(2))%32))) < base.Ui32(v117*v117) {
		v123 = v86
		goto L4
	} else {
		goto L7
	}
L6:
	;
	v123 = v43
	goto L4
L7:
	;
	if base.Ui32(v51) < base.Ui32(v86) {
		v86 = v86 + int32(-1)
		goto L5
	} else {
		goto L8
	}
L8:
	;
	goto L6
L9:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v39)+24)) = v131
	*(*int64)(unsafe.Add(mBase, uint32(v39)+16)) = v150
	*(*int32)(unsafe.Add(mBase, uint32(v39)+8)) = v131
	*(*int64)(unsafe.Add(mBase, uint32(v39))) = v150
	if v51 <= v134 {
		goto L12
	} else {
		goto L13
	}
L10:
	;
	v142 = m.G79
	v148 = int64(*(*uint16)(unsafe.Add(mBase, uint32(v142+(v60^int32(255))<<(uint(int32(1))%32)))))
	v150 = v148 * v140
	goto L9
L11:
	;
	if l4 != 0 {
		goto L57
	} else {
		goto L58
	}
L12:
	;
	v160 = int32(192)
	v179 = int32(-1)
	__phi182 = v51
	__phi190 = v179
	__phi193 = l4*v160 | v51*int32(12) + l0 + int32(_a_F_TrellisQuantizeBlock_2)
	__phi196 = v39 + int32(64) | v51<<(uint(int32(3))%32)
	__phi197 = v39 + int32(32)
	__phi198 = v39
	__phi199 = v139 * v140
	__phi201 = v150
	__phi202 = v179
	__phi203 = v179
	v182 = __phi182
	v190 = __phi190
	v193 = __phi193
	v196 = __phi196
	v197 = __phi197
	v198 = __phi198
	v199 = __phi199
	v201 = __phi201
	v202 = __phi202
	v203 = __phi203
	goto L14
L13:
	;
	v157 = int32(-1)
	v466 = v157
	v478 = int32(255)
	v479 = v157
	goto L11
L14:
	;
	v218 = m.G1
	v222 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v218+int32(_a_F_TrellisQuantizeBlock_1)+v182))))
	v224 = v222 << (uint(int32(1)) % 32)
	v226 = int32(*(*int16)(unsafe.Add(mBase, uint32(l1+v224))))
	v228 = v226 >> (uint(int32(31)) % 32)
	v232 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l5+v160+v224))))
	v233 = v226 ^ v228 - v228 + v232
	v235 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l5+int32(32)+v224))))
	v236 = v233 * v235
	v238 = int32(base.Ui32(v236) >> (uint(int32(17)) % 32))
	v239 = int32(2)
	if base.Ui32(v238) < base.Ui32(v239) {
		goto L16
	} else {
		goto L17
	}
L15:
	;
	v466 = v444
	v478 = v447
	v479 = v448
	goto L11
L16:
	;
	v242 = v238
	goto L18
L17:
	;
	v242 = v239
	goto L18
L18:
	;
	v246 = *(*int32)(unsafe.Add(mBase, uint32(v193+v242<<(uint(int32(2))%32))))
	*(*int32)(unsafe.Add(mBase, uint32(v197)+8)) = v246
	v251 = int32(base.Ui32(v236+int32(_a_F_TrellisQuantizeBlock_3)) >> (uint(int32(17)) % 32))
	v252 = int32(2047)
	if base.Ui32(v251) < base.Ui32(v252) {
		goto L19
	} else {
		goto L20
	}
L19:
	;
	v255 = v251
	goto L21
L20:
	;
	v255 = v252
	goto L21
L21:
	;
	v258 = v218 + int32(_a_F_TrellisQuantizeBlock_4) + v224
	v259 = int32(1)
	v260 = v233 << (uint(v259) % 32)
	v264 = int32(base.Ui32(v226&int32(_a_F_TrellisQuantizeBlock_5)) >> (uint(int32(15)) % 32))
	v265 = m.G81
	v269 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v265+v182+v259))))
	v271 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l5+v224))))
	v272 = int32(2047)
	if base.Ui32(v238) < base.Ui32(v272) {
		goto L24
	} else {
		goto L25
	}
L22:
	;
	v359 = v275 + int32(1)
	v360 = int32(2)
	if base.Ui32(v359) < base.Ui32(v360) {
		goto L39
	} else {
		goto L40
	}
L23:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v196+int32(2)))) = uint16(v275)
	*(*uint8)(unsafe.Add(mBase, uint32(v196+int32(1)))) = uint8(v264)
	v286 = m.G80
	v287 = *(*int32)(unsafe.Add(mBase, uint32(v198)+24))
	v288 = int32(67)
	if base.Ui32(v238) < base.Ui32(v288) {
		goto L28
	} else {
		goto L29
	}
L24:
	;
	v275 = v238
	goto L26
L25:
	;
	v275 = v272
	goto L26
L26:
	;
	if base.Ui32(v275) <= base.Ui32(v251) {
		goto L23
	} else {
		goto L27
	}
L27:
	;
	v277 = int64(36028797018963967)
	*(*int64)(unsafe.Add(mBase, uint32(v197))) = v277
	v349 = v190
	v351 = v199
	v352 = v277
	v353 = v202
	v354 = v203
	goto L22
L28:
	;
	v291 = v238
	goto L30
L29:
	;
	v291 = v288
	goto L30
L30:
	;
	v292 = int32(1)
	v293 = v291 << (uint(v292) % 32)
	v295 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v287+v293))))
	v299 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v286+v275<<(uint(v292)%32)))))
	v303 = *(*int64)(unsafe.Add(mBase, uint32(v198)+16))
	v304 = base.I64_extend_i32_u(v295+v299)*v140 + v303
	v305 = *(*int32)(unsafe.Add(mBase, uint32(v198)+8))
	v307 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v305+v293))))
	v311 = base.I64_extend_i32_u(v307+v299)*v140 + v201
	v312 = base.B2i32(v304 < v311)
	*(*uint8)(unsafe.Add(mBase, uint32(v196))) = uint8(v312)
	if v304 < v311 {
		goto L31
	} else {
		goto L32
	}
L31:
	;
	v314 = v304
	goto L33
L32:
	;
	v314 = v311
	goto L33
L33:
	;
	v315 = v275 * v271
	v318 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v258))))
	v323 = v314 + base.I64_extend_i32_s((v315-v260)*v315*v318)<<(uint(int64(8))%64)
	*(*int64)(unsafe.Add(mBase, uint32(v197))) = v323
	if base.Ui32(v236) < base.Ui32(int32(131072)) {
		v349 = v190
		v351 = v199
		v352 = v323
		v353 = v202
		v354 = v203
		goto L22
	} else {
		goto L34
	}
L34:
	;
	if v199 <= v323 {
		v349 = v190
		v351 = v199
		v352 = v323
		v353 = v202
		v354 = v203
		goto L22
	} else {
		goto L35
	}
L35:
	;
	if base.Ui32(int32(14)) < base.Ui32(v182) {
		v344 = int64(0)
		goto L36
	} else {
		goto L37
	}
L36:
	;
	v346 = v344*v140 + v323
	if v199 <= v346 {
		v349 = v190
		v351 = v199
		v352 = v323
		v353 = v202
		v354 = v203
		goto L22
	} else {
		goto L38
	}
L37:
	;
	v337 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48+v269*int32(33)+v242*int32(11)))))
	v338 = m.G79
	v342 = int64(*(*uint16)(unsafe.Add(mBase, uint32(v338+v337<<(uint(int32(1))%32)))))
	v344 = v342
	goto L36
L38:
	;
	v349 = v182
	v351 = v346
	v352 = v323
	v353 = v312
	v354 = int32(0)
	goto L22
L39:
	;
	v363 = v359
	goto L41
L40:
	;
	v363 = v360
	goto L41
L41:
	;
	v367 = *(*int32)(unsafe.Add(mBase, uint32(v193+v363<<(uint(int32(2))%32))))
	*(*int32)(unsafe.Add(mBase, uint32(v197)+24)) = v367
	if base.Ui32(v255) <= base.Ui32(v238) {
		goto L43
	} else {
		goto L44
	}
L42:
	;
	v456 = v182 + int32(1)
	if v134+int32(1) != v456 {
		__phi182 = v456
		__phi190 = v444
		__phi193 = v193 + int32(12)
		__phi196 = v196 + int32(8)
		__phi197 = v198
		__phi198 = v197
		__phi199 = v446
		__phi201 = v352
		__phi202 = v447
		__phi203 = v448
		v182 = __phi182
		v190 = __phi190
		v193 = __phi193
		v196 = __phi196
		v197 = __phi197
		v198 = __phi198
		v199 = __phi199
		v201 = __phi201
		v202 = __phi202
		v203 = __phi203
		goto L14
	} else {
		goto L55
	}
L43:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v197)+16)) = int64(36028797018963967)
	v444 = v349
	v446 = v351
	v447 = v353
	v448 = v354
	goto L42
L44:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v196+int32(6)))) = uint16(v359)
	*(*uint8)(unsafe.Add(mBase, uint32(v196+int32(5)))) = uint8(v264)
	v376 = m.G80
	v379 = *(*int32)(unsafe.Add(mBase, uint32(v198)+24))
	v380 = int32(67)
	if base.Ui32(v359) < base.Ui32(v380) {
		goto L45
	} else {
		goto L46
	}
L45:
	;
	v383 = v359
	goto L47
L46:
	;
	v383 = v380
	goto L47
L47:
	;
	v384 = int32(1)
	v385 = v383 << (uint(v384) % 32)
	v387 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v379+v385))))
	v391 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v376+v359<<(uint(v384)%32)))))
	v395 = *(*int64)(unsafe.Add(mBase, uint32(v198)+16))
	v396 = base.I64_extend_i32_u(v387+v391)*v140 + v395
	v397 = *(*int32)(unsafe.Add(mBase, uint32(v198)+8))
	v399 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v397+v385))))
	v403 = *(*int64)(unsafe.Add(mBase, uint32(v198)))
	v404 = base.I64_extend_i32_u(v399+v391)*v140 + v403
	v405 = base.B2i32(v396 < v404)
	*(*uint8)(unsafe.Add(mBase, uint32(v196+int32(4)))) = uint8(v405)
	if v396 < v404 {
		goto L48
	} else {
		goto L49
	}
L48:
	;
	v407 = v396
	goto L50
L49:
	;
	v407 = v404
	goto L50
L50:
	;
	v408 = v359 * v271
	v411 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v258))))
	v416 = v407 + base.I64_extend_i32_s((v408-v260)*v408*v411)<<(uint(int64(8))%64)
	*(*int64)(unsafe.Add(mBase, uint32(v197)+16)) = v416
	if v351 <= v416 {
		v444 = v349
		v446 = v351
		v447 = v353
		v448 = v354
		goto L42
	} else {
		goto L51
	}
L51:
	;
	if base.Ui32(int32(14)) < base.Ui32(v182) {
		v435 = int64(0)
		goto L52
	} else {
		goto L53
	}
L52:
	;
	v437 = v435*v140 + v416
	if v351 <= v437 {
		v444 = v349
		v446 = v351
		v447 = v353
		v448 = v354
		goto L42
	} else {
		goto L54
	}
L53:
	;
	v428 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48+v269*int32(33)+v363*int32(11)))))
	v429 = m.G79
	v433 = int64(*(*uint16)(unsafe.Add(mBase, uint32(v429+v428<<(uint(int32(1))%32)))))
	v435 = v433
	goto L52
L54:
	;
	v444 = v182
	v446 = v437
	v447 = v405
	v448 = int32(1)
	goto L42
L55:
	;
	goto L15
L56:
	;
	v529 = int32(0)
	if v466 == int32(-1) {
		v634 = v529
		goto L59
	} else {
		goto L60
	}
L57:
	;
	v515 = base.Simd_g_const(&F_TrellisQuantizeBlock__k0)
	v516 = int32(0)
	base.Simd_g_v128_store(m, l1, v516, v515)
	v518 = int32(16)
	base.Simd_g_v128_store(m, l1+v518, v516, v515)
	base.Simd_g_v128_store(m, l2+v518, v516, v515)
	base.Simd_g_v128_store(m, l2, v516, v515)
	goto L56
L58:
	;
	v494 = base.Simd_g_const(&F_TrellisQuantizeBlock__k0)
	v495 = int32(2)
	base.Simd_g_v128_store(m, l1, v495, v494)
	v497 = int32(24)
	v499 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(l1+v497))) = v499
	v501 = int32(18)
	*(*int64)(unsafe.Add(mBase, uint32(l1+v501))) = v499
	*(*int64)(unsafe.Add(mBase, uint32(l2+v497))) = v499
	*(*int64)(unsafe.Add(mBase, uint32(l2+v501))) = v499
	base.Simd_g_v128_store(m, l2, v495, v494)
	goto L56
L59:
	;
	m.G0 = v39 + int32(192)
	return v634
L60:
	;
	v536 = v39 + int32(64) + v466<<(uint(int32(3))%32)
	*(*uint8)(unsafe.Add(mBase, uint32(v536+v479<<(uint(int32(2))%32)))) = uint8(v478)
	if v466 < v51 {
		v634 = v529
		goto L59
	} else {
		goto L61
	}
L61:
	;
	v549 = int32(0)
	v554 = v466
	v559 = v536
	v562 = l2 + v466<<(uint(int32(1))%32)
	v567 = v479
	goto L62
L62:
	;
	v583 = int32(2)
	v585 = v559 + v567<<(uint(v583)%32)
	v588 = int32(*(*int16)(unsafe.Add(mBase, uint32(v585+v583))))
	v592 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v585+int32(1)))))
	if v592 != 0 {
		goto L64
	} else {
		goto L65
	}
L63:
	;
	v634 = base.B2i32(v611 != int32(0))
	goto L59
L64:
	;
	v593 = int32(0) - v588
	goto L66
L65:
	;
	v593 = v588
	goto L66
L66:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v562))) = uint16(v593)
	v595 = m.G1
	v599 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v595+int32(_a_F_TrellisQuantizeBlock_1)+v554))))
	v601 = v599 << (uint(int32(1)) % 32)
	v604 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l5+v601))))
	v605 = v604 * v593
	*(*uint16)(unsafe.Add(mBase, uint32(l1+v601))) = uint16(v605)
	v611 = v549 | v588
	v613 = int32(*(*int8)(unsafe.Add(mBase, uint32(v585))))
	if v51 < v554 {
		v549 = v611
		v554 = v554 + int32(-1)
		v559 = v559 + int32(-8)
		v562 = v562 + int32(-2)
		v567 = v613
		goto L62
	} else {
		goto L67
	}
L67:
	;
	goto L63
}

var F_TrellisQuantizeBlock__k0 = [2]uint64{0x0, 0x0}

func F_TrueMotion_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v20 int32
	_ = v20
	var v23 int32
	_ = v23
	var v24 base.V128
	_ = v24
	var v25 base.V128
	_ = v25
	var v27 base.V128
	_ = v27
	var v34 int32
	_ = v34
	var v42 int32
	_ = v42
	var v50 int32
	_ = v50
	var v58 int32
	_ = v58
	var v66 int32
	_ = v66
	var v74 int32
	_ = v74
	var v82 int32
	_ = v82
	var v92 int32
	_ = v92
	var v94 base.V128
	_ = v94
	var v95 int32
	_ = v95
	var v96 base.V128
	_ = v96
	var v97 base.V128
	_ = v97
	var v99 base.V128
	_ = v99
	var v102 base.V128
	_ = v102
	var v104 base.V128
	_ = v104
	var v107 int32
	_ = v107
	var v109 base.V128
	_ = v109
	var v112 base.V128
	_ = v112
	var v115 int32
	_ = v115
	var v117 base.V128
	_ = v117
	var v120 base.V128
	_ = v120
	var v123 int32
	_ = v123
	var v125 base.V128
	_ = v125
	var v128 base.V128
	_ = v128
	var v131 int32
	_ = v131
	var v133 base.V128
	_ = v133
	var v136 base.V128
	_ = v136
	var v139 int32
	_ = v139
	var v141 base.V128
	_ = v141
	var v144 base.V128
	_ = v144
	var v147 int32
	_ = v147
	var v149 base.V128
	_ = v149
	var v152 base.V128
	_ = v152
	var v155 int32
	_ = v155
	var v157 base.V128
	_ = v157
	var v160 base.V128
	_ = v160
	var v163 int32
	_ = v163
	var v165 base.V128
	_ = v165
	var v168 base.V128
	_ = v168
	var v171 int32
	_ = v171
	var v173 base.V128
	_ = v173
	var v176 base.V128
	_ = v176
	var v179 int32
	_ = v179
	var v181 base.V128
	_ = v181
	var v184 base.V128
	_ = v184
	var v187 int32
	_ = v187
	var v189 base.V128
	_ = v189
	var v192 base.V128
	_ = v192
	var v195 int32
	_ = v195
	var v197 base.V128
	_ = v197
	var v200 base.V128
	_ = v200
	var v203 int32
	_ = v203
	var v205 base.V128
	_ = v205
	var v208 base.V128
	_ = v208
	var v211 int32
	_ = v211
	var v213 base.V128
	_ = v213
	var v216 base.V128
	_ = v216
	var v219 int32
	_ = v219
	var v221 base.V128
	_ = v221
	var v224 base.V128
	_ = v224
	var v229 int32
	_ = v229
	var v230 base.V128
	_ = v230
	var v235 base.V128
	_ = v235
	var v240 base.V128
	_ = v240
	var v245 base.V128
	_ = v245
	var v250 base.V128
	_ = v250
	var v255 base.V128
	_ = v255
	var v260 base.V128
	_ = v260
	var v265 base.V128
	_ = v265
	var v270 int32
	_ = v270
	var v273 int32
	_ = v273
	var v281 int32
	_ = v281
	var v282 base.V128
	_ = v282
	var v288 int32
	_ = v288
	var v295 int64
	_ = v295
	var v304 int32
	_ = v304
	var v305 base.V128
	_ = v305
	var v340 int32
	_ = v340
	var v348 int64
	_ = v348
	var v364 base.V128
	_ = v364
	if l1 == int32(0) {
		if l2 == int32(0) {
			switch l3 + int32(-4) {
			case 0:
				v340 = int32(-2122219135)
				*(*int32)(unsafe.Add(mBase, uint32(l0)+96)) = v340
				*(*int32)(unsafe.Add(mBase, uint32(l0)+64)) = v340
				*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = v340
				*(*int32)(unsafe.Add(mBase, uint32(l0))) = v340
				return
			default:
				v364 = base.Simd_g_const(&F_TrueMotion_SSE2__k0)
				base.Simd_g_v128_store(m, l0, int32(480), v364)
				base.Simd_g_v128_store(m, l0, int32(448), v364)
				base.Simd_g_v128_store(m, l0, int32(416), v364)
				base.Simd_g_v128_store(m, l0, int32(384), v364)
				base.Simd_g_v128_store(m, l0, int32(352), v364)
				base.Simd_g_v128_store(m, l0, int32(320), v364)
				base.Simd_g_v128_store(m, l0, int32(288), v364)
				base.Simd_g_v128_store(m, l0, int32(256), v364)
				base.Simd_g_v128_store(m, l0, int32(224), v364)
				base.Simd_g_v128_store(m, l0, int32(192), v364)
				base.Simd_g_v128_store(m, l0, int32(160), v364)
				base.Simd_g_v128_store(m, l0, int32(128), v364)
				base.Simd_g_v128_store(m, l0, int32(96), v364)
				base.Simd_g_v128_store(m, l0, int32(64), v364)
				base.Simd_g_v128_store(m, l0, int32(32), v364)
				base.Simd_g_v128_store(m, l0, int32(0), v364)
				return
			case 4:
				v348 = int64(-9114861777597660799)
				*(*int64)(unsafe.Add(mBase, uint32(l0)+224)) = v348
				*(*int64)(unsafe.Add(mBase, uint32(l0)+192)) = v348
				*(*int64)(unsafe.Add(mBase, uint32(l0)+160)) = v348
				*(*int64)(unsafe.Add(mBase, uint32(l0)+128)) = v348
				*(*int64)(unsafe.Add(mBase, uint32(l0)+96)) = v348
				*(*int64)(unsafe.Add(mBase, uint32(l0)+64)) = v348
				*(*int64)(unsafe.Add(mBase, uint32(l0)+32)) = v348
				*(*int64)(unsafe.Add(mBase, uint32(l0))) = v348
				return
			}
		} else {
			if l3 != int32(8) {
				v304 = int32(0)
				v305 = base.Simd_g_v128_load(m, l2, v304)
				base.Simd_g_v128_store(m, l0, int32(480), v305)
				base.Simd_g_v128_store(m, l0, int32(448), v305)
				base.Simd_g_v128_store(m, l0, int32(416), v305)
				base.Simd_g_v128_store(m, l0, int32(384), v305)
				base.Simd_g_v128_store(m, l0, int32(352), v305)
				base.Simd_g_v128_store(m, l0, int32(320), v305)
				base.Simd_g_v128_store(m, l0, int32(288), v305)
				base.Simd_g_v128_store(m, l0, int32(256), v305)
				base.Simd_g_v128_store(m, l0, int32(224), v305)
				base.Simd_g_v128_store(m, l0, int32(192), v305)
				base.Simd_g_v128_store(m, l0, int32(160), v305)
				base.Simd_g_v128_store(m, l0, int32(128), v305)
				base.Simd_g_v128_store(m, l0, int32(96), v305)
				base.Simd_g_v128_store(m, l0, int32(64), v305)
				base.Simd_g_v128_store(m, l0, int32(32), v305)
				base.Simd_g_v128_store(m, l0, v304, v305)
				return
			} else {
				v295 = *(*int64)(unsafe.Add(mBase, uint32(l2)))
				*(*int64)(unsafe.Add(mBase, uint32(l0)+224)) = v295
				*(*int64)(unsafe.Add(mBase, uint32(l0)+192)) = v295
				*(*int64)(unsafe.Add(mBase, uint32(l0)+160)) = v295
				*(*int64)(unsafe.Add(mBase, uint32(l0)+128)) = v295
				*(*int64)(unsafe.Add(mBase, uint32(l0)+96)) = v295
				*(*int64)(unsafe.Add(mBase, uint32(l0)+64)) = v295
				*(*int64)(unsafe.Add(mBase, uint32(l0)+32)) = v295
				*(*int64)(unsafe.Add(mBase, uint32(l0))) = v295
				return
			}
		}
	} else {
		if l2 == int32(0) {
			if l3 != int32(8) {
				v270 = l0
				v273 = int32(0)
				for {
					v281 = int32(0)
					v282 = base.Simd_g_v128_load8_splat(m, l1+v273, v281)
					base.Simd_g_v128_store(m, v270, v281, v282)
					v288 = v273 + int32(1)
					if v288 != int32(16) {
						v270 = v270 + int32(32)
						v273 = v288
						continue
					} else {
						break
					}
					break
				}
				return
			} else {
				v229 = int32(0)
				v230 = base.Simd_g_v128_load8_splat(m, l1, v229)
				base.Simd_g_v128_store64_lane_l0(m, l0, v229, v230)
				v235 = base.Simd_g_v128_load8_splat(m, l1, int32(1))
				base.Simd_g_v128_store64_lane_l0(m, l0, int32(32), v235)
				v240 = base.Simd_g_v128_load8_splat(m, l1, int32(2))
				base.Simd_g_v128_store64_lane_l0(m, l0, int32(64), v240)
				v245 = base.Simd_g_v128_load8_splat(m, l1, int32(3))
				base.Simd_g_v128_store64_lane_l0(m, l0, int32(96), v245)
				v250 = base.Simd_g_v128_load8_splat(m, l1, int32(4))
				base.Simd_g_v128_store64_lane_l0(m, l0, int32(128), v250)
				v255 = base.Simd_g_v128_load8_splat(m, l1, int32(5))
				base.Simd_g_v128_store64_lane_l0(m, l0, int32(160), v255)
				v260 = base.Simd_g_v128_load8_splat(m, l1, int32(6))
				base.Simd_g_v128_store64_lane_l0(m, l0, int32(192), v260)
				v265 = base.Simd_g_v128_load8_splat(m, l1, int32(7))
				base.Simd_g_v128_store64_lane_l0(m, l0, int32(224), v265)
				return
			}
		} else {
			v15 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
			if l3 != int32(8) {
				v92 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+int32(-1)))))
				v94 = base.Simd_g_i16x8_splat(v15 - v92)
				v95 = int32(0)
				v96 = base.Simd_g_v128_load(m, l2, v95)
				v97 = base.Simd_g_const(&F_TrueMotion_SSE2__k1)
				v99 = base.Simd_g_i8x16_shuffle2(v96, v97, base.Simd_g_const(&F_TrueMotion_SSE2__k2), base.Simd_g_const(&F_TrueMotion_SSE2__k3))
				v102 = base.Simd_g_i8x16_shuffle2(v96, v97, base.Simd_g_const(&F_TrueMotion_SSE2__k4), base.Simd_g_const(&F_TrueMotion_SSE2__k5))
				v104 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v94, v99), base.Simd_g_i16x8_add(v94, v102))
				base.Simd_g_v128_store(m, l0, v95, v104)
				v107 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+1)))
				v109 = base.Simd_g_i16x8_splat(v107 - v92)
				v112 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v109, v99), base.Simd_g_i16x8_add(v109, v102))
				base.Simd_g_v128_store(m, l0, int32(32), v112)
				v115 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+2)))
				v117 = base.Simd_g_i16x8_splat(v115 - v92)
				v120 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v117, v99), base.Simd_g_i16x8_add(v117, v102))
				base.Simd_g_v128_store(m, l0, int32(64), v120)
				v123 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3)))
				v125 = base.Simd_g_i16x8_splat(v123 - v92)
				v128 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v125, v99), base.Simd_g_i16x8_add(v125, v102))
				base.Simd_g_v128_store(m, l0, int32(96), v128)
				v131 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+4)))
				v133 = base.Simd_g_i16x8_splat(v131 - v92)
				v136 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v133, v99), base.Simd_g_i16x8_add(v133, v102))
				base.Simd_g_v128_store(m, l0, int32(128), v136)
				v139 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+5)))
				v141 = base.Simd_g_i16x8_splat(v139 - v92)
				v144 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v141, v99), base.Simd_g_i16x8_add(v141, v102))
				base.Simd_g_v128_store(m, l0, int32(160), v144)
				v147 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+6)))
				v149 = base.Simd_g_i16x8_splat(v147 - v92)
				v152 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v149, v99), base.Simd_g_i16x8_add(v149, v102))
				base.Simd_g_v128_store(m, l0, int32(192), v152)
				v155 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+7)))
				v157 = base.Simd_g_i16x8_splat(v155 - v92)
				v160 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v157, v99), base.Simd_g_i16x8_add(v157, v102))
				base.Simd_g_v128_store(m, l0, int32(224), v160)
				v163 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+8)))
				v165 = base.Simd_g_i16x8_splat(v163 - v92)
				v168 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v165, v99), base.Simd_g_i16x8_add(v165, v102))
				base.Simd_g_v128_store(m, l0, int32(256), v168)
				v171 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+9)))
				v173 = base.Simd_g_i16x8_splat(v171 - v92)
				v176 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v173, v99), base.Simd_g_i16x8_add(v173, v102))
				base.Simd_g_v128_store(m, l0, int32(288), v176)
				v179 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+10)))
				v181 = base.Simd_g_i16x8_splat(v179 - v92)
				v184 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v181, v99), base.Simd_g_i16x8_add(v181, v102))
				base.Simd_g_v128_store(m, l0, int32(320), v184)
				v187 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+11)))
				v189 = base.Simd_g_i16x8_splat(v187 - v92)
				v192 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v189, v99), base.Simd_g_i16x8_add(v189, v102))
				base.Simd_g_v128_store(m, l0, int32(352), v192)
				v195 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+12)))
				v197 = base.Simd_g_i16x8_splat(v195 - v92)
				v200 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v197, v99), base.Simd_g_i16x8_add(v197, v102))
				base.Simd_g_v128_store(m, l0, int32(384), v200)
				v203 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+13)))
				v205 = base.Simd_g_i16x8_splat(v203 - v92)
				v208 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v205, v99), base.Simd_g_i16x8_add(v205, v102))
				base.Simd_g_v128_store(m, l0, int32(416), v208)
				v211 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+14)))
				v213 = base.Simd_g_i16x8_splat(v211 - v92)
				v216 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v213, v99), base.Simd_g_i16x8_add(v213, v102))
				base.Simd_g_v128_store(m, l0, int32(448), v216)
				v219 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+15)))
				v221 = base.Simd_g_i16x8_splat(v219 - v92)
				v224 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v221, v99), base.Simd_g_i16x8_add(v221, v102))
				base.Simd_g_v128_store(m, l0, int32(480), v224)
				return
			} else {
				v20 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+int32(-1)))))
				v23 = int32(0)
				v24 = base.Simd_g_v128_load64_zero(m, l2, v23)
				v25 = base.Simd_g_const(&F_TrueMotion_SSE2__k1)
				v27 = base.Simd_g_i8x16_shuffle2(v24, v25, base.Simd_g_const(&F_TrueMotion_SSE2__k2), base.Simd_g_const(&F_TrueMotion_SSE2__k3))
				base.Simd_g_v128_store64_lane_l0(m, l0, v23, base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v15-v20), v27), v25))
				v34 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+1)))
				base.Simd_g_v128_store64_lane_l0(m, l0, int32(32), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v34-v20), v27), v25))
				v42 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+2)))
				base.Simd_g_v128_store64_lane_l0(m, l0, int32(64), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v42-v20), v27), v25))
				v50 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3)))
				base.Simd_g_v128_store64_lane_l0(m, l0, int32(96), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v50-v20), v27), v25))
				v58 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+4)))
				base.Simd_g_v128_store64_lane_l0(m, l0, int32(128), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v58-v20), v27), v25))
				v66 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+5)))
				base.Simd_g_v128_store64_lane_l0(m, l0, int32(160), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v66-v20), v27), v25))
				v74 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+6)))
				base.Simd_g_v128_store64_lane_l0(m, l0, int32(192), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v74-v20), v27), v25))
				v82 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+7)))
				base.Simd_g_v128_store64_lane_l0(m, l0, int32(224), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_splat(v82-v20), v27), v25))
				return
			}
		}
	}
}

var F_TrueMotion_SSE2__k0 = [2]uint64{0x8181818181818181, 0x8181818181818181}
var F_TrueMotion_SSE2__k1 = [2]uint64{0x0, 0x0}
var F_TrueMotion_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_TrueMotion_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_TrueMotion_SSE2__k4 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_TrueMotion_SSE2__k5 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
