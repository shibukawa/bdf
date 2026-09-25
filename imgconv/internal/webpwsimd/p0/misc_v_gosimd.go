//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_VE16_C(m *base.Module, l0 int32) {
	var v5 int32
	_ = v5
	var v6 base.V128
	_ = v6
	v5 = int32(0)
	v6 = base.Simd_g_v128_load(m, l0+int32(-32), v5)
	base.Simd_g_v128_store(m, l0, v5, v6)
	base.Simd_g_v128_store(m, l0, int32(32), v6)
	base.Simd_g_v128_store(m, l0, int32(64), v6)
	base.Simd_g_v128_store(m, l0, int32(96), v6)
	base.Simd_g_v128_store(m, l0, int32(128), v6)
	base.Simd_g_v128_store(m, l0, int32(256), v6)
	base.Simd_g_v128_store(m, l0, int32(224), v6)
	base.Simd_g_v128_store(m, l0, int32(192), v6)
	base.Simd_g_v128_store(m, l0, int32(160), v6)
	base.Simd_g_v128_store(m, l0, int32(384), v6)
	base.Simd_g_v128_store(m, l0, int32(352), v6)
	base.Simd_g_v128_store(m, l0, int32(320), v6)
	base.Simd_g_v128_store(m, l0, int32(288), v6)
	base.Simd_g_v128_store(m, l0, int32(480), v6)
	base.Simd_g_v128_store(m, l0, int32(448), v6)
	base.Simd_g_v128_store(m, l0, int32(416), v6)
	return
}
func F_VE16_SSE2(m *base.Module, l0 int32) {
	var v5 int32
	_ = v5
	var v6 base.V128
	_ = v6
	v5 = int32(0)
	v6 = base.Simd_g_v128_load(m, l0+int32(-32), v5)
	base.Simd_g_v128_store(m, l0, int32(480), v6)
	base.Simd_g_v128_store(m, l0, int32(448), v6)
	base.Simd_g_v128_store(m, l0, int32(416), v6)
	base.Simd_g_v128_store(m, l0, int32(384), v6)
	base.Simd_g_v128_store(m, l0, int32(352), v6)
	base.Simd_g_v128_store(m, l0, int32(320), v6)
	base.Simd_g_v128_store(m, l0, int32(288), v6)
	base.Simd_g_v128_store(m, l0, int32(256), v6)
	base.Simd_g_v128_store(m, l0, int32(224), v6)
	base.Simd_g_v128_store(m, l0, int32(192), v6)
	base.Simd_g_v128_store(m, l0, int32(160), v6)
	base.Simd_g_v128_store(m, l0, int32(128), v6)
	base.Simd_g_v128_store(m, l0, int32(96), v6)
	base.Simd_g_v128_store(m, l0, int32(64), v6)
	base.Simd_g_v128_store(m, l0, int32(32), v6)
	base.Simd_g_v128_store(m, l0, v5, v6)
	return
}
func F_VE4_SSE2(m *base.Module, l0 int32) {
	var v2 base.V128
	_ = v2
	var v7 int32
	_ = v7
	var v8 base.V128
	_ = v8
	var v11 base.V128
	_ = v11
	var v19 base.V128
	_ = v19
	v2 = base.Simd_g_const(&F_VE4_SSE2__k0)
	v7 = int32(0)
	v8 = base.Simd_g_v128_load64_zero(m, l0+int32(-33), v7)
	v11 = base.Simd_g_i8x16_shuffle2(v8, v2, base.Simd_g_const(&F_VE4_SSE2__k1), base.Simd_g_const(&F_VE4_SSE2__k2))
	v19 = base.Simd_g_i8x16_avgr_u(base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_avgr_u(v8, v11), base.Simd_g_v128_and(base.Simd_g_v128_xor(v11, v8), base.Simd_g_const(&F_VE4_SSE2__k3))), base.Simd_g_i8x16_shuffle2(v8, v2, base.Simd_g_const(&F_VE4_SSE2__k4), base.Simd_g_const(&F_VE4_SSE2__k5)))
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(96), v19)
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(64), v19)
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(32), v19)
	base.Simd_g_v128_store32_lane_l0(m, l0, v7, v19)
	return
}

var F_VE4_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_VE4_SSE2__k1 = [2]uint64{0x908070605040302, 0x80800f0e0d0c0b0a}
var F_VE4_SSE2__k2 = [2]uint64{0x8080808080808080, 0x100808080808080}
var F_VE4_SSE2__k3 = [2]uint64{0x101010101010101, 0x101010101010101}
var F_VE4_SSE2__k4 = [2]uint64{0x807060504030201, 0x800f0e0d0c0b0a09}
var F_VE4_SSE2__k5 = [2]uint64{0x8080808080808080, 0x80808080808080}

func F_VFilter16_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	var v6 int32
	_ = v6
	var v7 base.V128
	_ = v7
	var v32 int32
	_ = v32
	var v35 base.V128
	_ = v35
	var v36 base.V128
	_ = v36
	var v40 int32
	_ = v40
	var v43 base.V128
	_ = v43
	var v45 base.V128
	_ = v45
	var v48 base.V128
	_ = v48
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v54 base.V128
	_ = v54
	var v55 int32
	_ = v55
	var v57 base.V128
	_ = v57
	var v70 base.V128
	_ = v70
	var v72 base.V128
	_ = v72
	var v83 base.V128
	_ = v83
	var v87 base.V128
	_ = v87
	var v88 int32
	_ = v88
	var v90 base.V128
	_ = v90
	var v102 base.V128
	_ = v102
	var v103 base.V128
	_ = v103
	var v104 base.V128
	_ = v104
	var v105 base.V128
	_ = v105
	var v106 base.V128
	_ = v106
	var v107 base.V128
	_ = v107
	var v111 base.V128
	_ = v111
	var v116 base.V128
	_ = v116
	var v117 base.V128
	_ = v117
	var v118 base.V128
	_ = v118
	var v119 base.V128
	_ = v119
	var v121 base.V128
	_ = v121
	var v125 base.V128
	_ = v125
	var v126 base.V128
	_ = v126
	var v127 base.V128
	_ = v127
	var v128 base.V128
	_ = v128
	var v129 int32
	_ = v129
	var v131 base.V128
	_ = v131
	var v132 base.V128
	_ = v132
	var v138 base.V128
	_ = v138
	var v139 base.V128
	_ = v139
	var v142 base.V128
	_ = v142
	var v145 base.V128
	_ = v145
	var v149 base.V128
	_ = v149
	var v152 base.V128
	_ = v152
	var v155 base.V128
	_ = v155
	var v157 base.V128
	_ = v157
	var v162 base.V128
	_ = v162
	var v164 base.V128
	_ = v164
	var v167 int32
	_ = v167
	var v181 base.V128
	_ = v181
	var v183 base.V128
	_ = v183
	var v187 base.V128
	_ = v187
	var v199 base.V128
	_ = v199
	var v203 base.V128
	_ = v203
	var v208 base.V128
	_ = v208
	v6 = int32(0)
	v7 = base.Simd_g_const(&F_VFilter16_SSE2__k0)
	v32 = l0 - l1<<(uint(int32(2))%32)
	v35 = base.Simd_g_v128_load(m, v32+l1, v6)
	v36 = base.Simd_g_const(&F_VFilter16_SSE2__k1)
	v40 = l1 * int32(3)
	v43 = base.Simd_g_v128_load(m, v32+v40, v6)
	v45 = base.Simd_g_v128_load(m, l0, v6)
	v48 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v43, v45), base.Simd_g_i8x16_sub_sat_u(v45, v43))
	v50 = int32(1)
	v51 = l1 << (uint(v50) % 32)
	v54 = base.Simd_g_v128_load(m, v32+v51, v6)
	v55 = l0 + l1
	v57 = base.Simd_g_v128_load(m, v55, v6)
	v70 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v54, v43), base.Simd_g_i8x16_sub_sat_u(v43, v54))
	v72 = base.Simd_g_v128_load(m, v32, v6)
	v83 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v57, v45), base.Simd_g_i8x16_sub_sat_u(v45, v57))
	v87 = base.Simd_g_v128_load(m, l0+v40, v6)
	v88 = l0 + v51
	v90 = base.Simd_g_v128_load(m, v88, v6)
	v102 = base.Simd_g_i8x16_eq(v7, base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_add_sat_u(base.Simd_g_i8x16_add_sat_u(v48, v48), base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v54, v57), base.Simd_g_i8x16_sub_sat_u(v57, v54)), v50), base.Simd_g_const(&F_VFilter16_SSE2__k2))), base.Simd_g_i8x16_splat(l2)), base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(v70, base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v72, v35), base.Simd_g_i8x16_sub_sat_u(v35, v72))), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v35, v54), base.Simd_g_i8x16_sub_sat_u(v54, v35))), v83), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v87, v90), base.Simd_g_i8x16_sub_sat_u(v90, v87))), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v90, v57), base.Simd_g_i8x16_sub_sat_u(v57, v90))), base.Simd_g_i8x16_splat(l3))))
	v103 = base.Simd_g_v128_xor(v45, v36)
	v104 = base.Simd_g_v128_xor(v43, v36)
	v105 = base.Simd_g_i8x16_sub_sat_s(v103, v104)
	v106 = base.Simd_g_v128_xor(v54, v36)
	v107 = base.Simd_g_v128_xor(v57, v36)
	v111 = base.Simd_g_i8x16_add_sat_s(v105, base.Simd_g_i8x16_add_sat_s(v105, base.Simd_g_i8x16_add_sat_s(v105, base.Simd_g_i8x16_sub_sat_s(v106, v107))))
	v116 = base.Simd_g_i8x16_eq(v7, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_max_u(v70, v83), base.Simd_g_i8x16_splat(l4)))
	v117 = base.Simd_g_v128_and(base.Simd_g_v128_and(v102, v111), v116)
	v118 = base.Simd_g_const(&F_VFilter16_SSE2__k3)
	v119 = base.Simd_g_i8x16_shuffle2(v7, v117, base.Simd_g_const(&F_VFilter16_SSE2__k4), base.Simd_g_const(&F_VFilter16_SSE2__k5))
	v121 = base.Simd_g_const(&F_VFilter16_SSE2__k6)
	v125 = base.Simd_g_const(&F_VFilter16_SSE2__k7)
	v126 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v119), v121), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v119), v121), base.Simd_g_const(&F_VFilter16_SSE2__k8), base.Simd_g_const(&F_VFilter16_SSE2__k9))
	v127 = base.Simd_g_const(&F_VFilter16_SSE2__k10)
	v128 = base.Simd_g_i16x8_add(v126, v127)
	v129 = int32(7)
	v131 = base.Simd_g_const(&F_VFilter16_SSE2__k11)
	v132 = base.Simd_g_i8x16_shuffle2(v7, v117, base.Simd_g_const(&F_VFilter16_SSE2__k12), base.Simd_g_const(&F_VFilter16_SSE2__k13))
	v138 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v132), v121), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v132), v121), base.Simd_g_const(&F_VFilter16_SSE2__k8), base.Simd_g_const(&F_VFilter16_SSE2__k9))
	v139 = base.Simd_g_i16x8_add(v138, v127)
	v142 = base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(v128, v129), base.Simd_g_i16x8_shr_s(v139, v129))
	v145 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(base.Simd_g_v128_xor(v35, v36), v142), v36)
	base.Simd_g_v128_store(m, l0+l1*int32(-3), v6, v145)
	v149 = base.Simd_g_i16x8_add(v128, v126)
	v152 = base.Simd_g_i16x8_add(v139, v138)
	v155 = base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(v149, v129), base.Simd_g_i16x8_shr_s(v152, v129))
	v157 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v106, v155), v36)
	base.Simd_g_v128_store(m, l0-v51, v6, v157)
	v162 = base.Simd_g_v128_and(base.Simd_g_v128_andnot(v111, v116), v102)
	v164 = base.Simd_g_i8x16_add_sat_s(v162, base.Simd_g_const(&F_VFilter16_SSE2__k14))
	v167 = int32(11)
	v181 = base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(v149, v126), v129), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(v152, v138), v129))
	v183 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(base.Simd_g_i8x16_add_sat_s(v104, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v7, v164, base.Simd_g_const(&F_VFilter16_SSE2__k4), base.Simd_g_const(&F_VFilter16_SSE2__k5)), v167), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v7, v164, base.Simd_g_const(&F_VFilter16_SSE2__k12), base.Simd_g_const(&F_VFilter16_SSE2__k13)), v167))), v181), v36)
	base.Simd_g_v128_store(m, l0-l1, v6, v183)
	v187 = base.Simd_g_i8x16_add_sat_s(v162, base.Simd_g_const(&F_VFilter16_SSE2__k15))
	v199 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(base.Simd_g_i8x16_sub_sat_s(v103, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v7, v187, base.Simd_g_const(&F_VFilter16_SSE2__k4), base.Simd_g_const(&F_VFilter16_SSE2__k5)), v167), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v7, v187, base.Simd_g_const(&F_VFilter16_SSE2__k12), base.Simd_g_const(&F_VFilter16_SSE2__k13)), v167))), v181), v36)
	base.Simd_g_v128_store(m, l0, v6, v199)
	v203 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v107, v155), v36)
	base.Simd_g_v128_store(m, v55, v6, v203)
	v208 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(base.Simd_g_v128_xor(v90, v36), v142), v36)
	base.Simd_g_v128_store(m, v88, v6, v208)
	return
}

var F_VFilter16_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_VFilter16_SSE2__k1 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_VFilter16_SSE2__k2 = [2]uint64{0x7f7f7f7f7f7f7f7f, 0x7f7f7f7f7f7f7f7f}
var F_VFilter16_SSE2__k3 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_VFilter16_SSE2__k4 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_VFilter16_SSE2__k5 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_VFilter16_SSE2__k6 = [2]uint64{0x90000000900, 0x90000000900}
var F_VFilter16_SSE2__k7 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_VFilter16_SSE2__k8 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_VFilter16_SSE2__k9 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_VFilter16_SSE2__k10 = [2]uint64{0x3f003f003f003f, 0x3f003f003f003f}
var F_VFilter16_SSE2__k11 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_VFilter16_SSE2__k12 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_VFilter16_SSE2__k13 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_VFilter16_SSE2__k14 = [2]uint64{0x303030303030303, 0x303030303030303}
var F_VFilter16_SSE2__k15 = [2]uint64{0x404040404040404, 0x404040404040404}

func F_VFilter16i_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	var v6 int32
	_ = v6
	var v40 base.V128
	_ = v40
	var v42 int32
	_ = v42
	var v45 base.V128
	_ = v45
	var v47 int32
	_ = v47
	var v50 base.V128
	_ = v50
	var v55 base.V128
	_ = v55
	var v57 int32
	_ = v57
	var v61 int32
	_ = v61
	var v66 base.V128
	_ = v66
	var v68 base.V128
	_ = v68
	var v70 base.V128
	_ = v70
	var v74 base.V128
	_ = v74
	var v87 base.V128
	_ = v87
	var v88 base.V128
	_ = v88
	var v89 base.V128
	_ = v89
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v92 base.V128
	_ = v92
	var v93 base.V128
	_ = v93
	var v97 base.V128
	_ = v97
	var v98 int32
	_ = v98
	var v100 base.V128
	_ = v100
	var v103 base.V128
	_ = v103
	var v106 base.V128
	_ = v106
	var v108 base.V128
	_ = v108
	var v109 base.V128
	_ = v109
	var v110 base.V128
	_ = v110
	var v124 base.V128
	_ = v124
	var v127 base.V128
	_ = v127
	var v140 base.V128
	_ = v140
	var v145 int32
	_ = v145
	var v153 base.V128
	_ = v153
	var v155 base.V128
	_ = v155
	var v156 base.V128
	_ = v156
	var v158 int32
	_ = v158
	var v160 base.V128
	_ = v160
	var v164 base.V128
	_ = v164
	var v169 base.V128
	_ = v169
	var v172 base.V128
	_ = v172
	var v177 base.V128
	_ = v177
	var v188 base.V128
	_ = v188
	var v192 base.V128
	_ = v192
	var v196 base.V128
	_ = v196
	var v200 int32
	_ = v200
	v6 = int32(0)
	v40 = base.Simd_g_v128_load(m, l0+l1, v6)
	v42 = l1 * int32(3)
	v45 = base.Simd_g_v128_load(m, l0+v42, v6)
	v47 = l1 << (uint(int32(1)) % 32)
	v50 = base.Simd_g_v128_load(m, l0+v47, v6)
	v55 = base.Simd_g_v128_load(m, l0, v6)
	v57 = l0
	v61 = int32(4)
	v66 = v40
	v68 = v45
	v70 = v50
	v74 = v55
	for {
		v87 = base.Simd_g_const(&F_VFilter16i_SSE2__k0)
		v88 = base.Simd_g_v128_xor(v70, v87)
		v89 = base.Simd_g_const(&F_VFilter16i_SSE2__k1)
		v90 = v57 + l1*int32(5)
		v91 = int32(0)
		v92 = base.Simd_g_v128_load(m, v90, v91)
		v93 = base.Simd_g_v128_xor(v92, v87)
		v97 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v70, v68), base.Simd_g_i8x16_sub_sat_u(v68, v70))
		v98 = v57 + l1<<(uint(int32(2))%32)
		v100 = base.Simd_g_v128_load(m, v98, v91)
		v103 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v100, v92), base.Simd_g_i8x16_sub_sat_u(v92, v100))
		v106 = base.Simd_g_i8x16_eq(v89, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_max_u(v97, v103), base.Simd_g_i8x16_splat(l4)))
		v108 = base.Simd_g_v128_xor(v100, v87)
		v109 = base.Simd_g_v128_xor(v68, v87)
		v110 = base.Simd_g_i8x16_sub_sat_s(v108, v109)
		v124 = base.Simd_g_v128_load(m, v57+l1*int32(6), v91)
		v127 = base.Simd_g_v128_load(m, v57+l1*int32(7), v91)
		v140 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v68, v100), base.Simd_g_i8x16_sub_sat_u(v100, v68))
		v145 = int32(1)
		v153 = base.Simd_g_v128_bitselect(base.Simd_g_i8x16_add_sat_s(base.Simd_g_i8x16_add_sat_s(base.Simd_g_i8x16_add_sat_s(base.Simd_g_v128_bitselect(v89, base.Simd_g_i8x16_sub_sat_s(v88, v93), v106), v110), v110), v110), v89, base.Simd_g_i8x16_eq(v89, base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(v97, base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v74, v66), base.Simd_g_i8x16_sub_sat_u(v66, v74))), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v66, v70), base.Simd_g_i8x16_sub_sat_u(v70, v66))), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v124, v127), base.Simd_g_i8x16_sub_sat_u(v127, v124))), v103), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v92, v124), base.Simd_g_i8x16_sub_sat_u(v124, v92))), base.Simd_g_i8x16_splat(l3)), base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_add_sat_u(base.Simd_g_i8x16_add_sat_u(v140, v140), base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v70, v92), base.Simd_g_i8x16_sub_sat_u(v92, v70)), v145), base.Simd_g_const(&F_VFilter16i_SSE2__k2))), base.Simd_g_i8x16_splat(l2)))))
		v155 = base.Simd_g_i8x16_add_sat_s(v153, base.Simd_g_const(&F_VFilter16i_SSE2__k3))
		v156 = base.Simd_g_const(&F_VFilter16i_SSE2__k4)
		v158 = int32(11)
		v160 = base.Simd_g_const(&F_VFilter16i_SSE2__k5)
		v164 = base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v89, v155, base.Simd_g_const(&F_VFilter16i_SSE2__k6), base.Simd_g_const(&F_VFilter16i_SSE2__k7)), v158), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v89, v155, base.Simd_g_const(&F_VFilter16i_SSE2__k8), base.Simd_g_const(&F_VFilter16i_SSE2__k9)), v158))
		v169 = base.Simd_g_v128_bitselect(base.Simd_g_i8x16_add(base.Simd_g_i8x16_avgr_u(base.Simd_g_v128_xor(v164, v87), v89), base.Simd_g_const(&F_VFilter16i_SSE2__k10)), v89, v106)
		v172 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v88, v169), v87)
		base.Simd_g_v128_store(m, v57+v47, v91, v172)
		v177 = base.Simd_g_i8x16_add_sat_s(v153, base.Simd_g_const(&F_VFilter16i_SSE2__k11))
		v188 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v109, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v89, v177, base.Simd_g_const(&F_VFilter16i_SSE2__k6), base.Simd_g_const(&F_VFilter16i_SSE2__k7)), v158), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v89, v177, base.Simd_g_const(&F_VFilter16i_SSE2__k8), base.Simd_g_const(&F_VFilter16i_SSE2__k9)), v158))), v87)
		base.Simd_g_v128_store(m, v57+v42, v91, v188)
		v192 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v108, v164), v87)
		base.Simd_g_v128_store(m, v98, v91, v192)
		v196 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v93, v169), v87)
		base.Simd_g_v128_store(m, v90, v91, v196)
		v200 = v61 + int32(-1)
		if base.Ui32(v145) < base.Ui32(v200) {
			v57 = v98
			v61 = v200
			v66 = v196
			v68 = v127
			v70 = v124
			v74 = v192
			continue
		} else {
			break
		}
		break
	}
	return
}

var F_VFilter16i_SSE2__k0 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_VFilter16i_SSE2__k1 = [2]uint64{0x0, 0x0}
var F_VFilter16i_SSE2__k2 = [2]uint64{0x7f7f7f7f7f7f7f7f, 0x7f7f7f7f7f7f7f7f}
var F_VFilter16i_SSE2__k3 = [2]uint64{0x404040404040404, 0x404040404040404}
var F_VFilter16i_SSE2__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_VFilter16i_SSE2__k5 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_VFilter16i_SSE2__k6 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_VFilter16i_SSE2__k7 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_VFilter16i_SSE2__k8 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_VFilter16i_SSE2__k9 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_VFilter16i_SSE2__k10 = [2]uint64{0xc0c0c0c0c0c0c0c0, 0xc0c0c0c0c0c0c0c0}
var F_VFilter16i_SSE2__k11 = [2]uint64{0x303030303030303, 0x303030303030303}

func F_VFilter8_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	var v7 int32
	_ = v7
	var v11 base.V128
	_ = v11
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v38 base.V128
	_ = v38
	var v39 int32
	_ = v39
	var v42 base.V128
	_ = v42
	var v43 base.V128
	_ = v43
	var v44 base.V128
	_ = v44
	var v45 base.V128
	_ = v45
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v52 base.V128
	_ = v52
	var v55 base.V128
	_ = v55
	var v57 base.V128
	_ = v57
	var v59 int32
	_ = v59
	var v62 base.V128
	_ = v62
	var v65 base.V128
	_ = v65
	var v67 base.V128
	_ = v67
	var v70 base.V128
	_ = v70
	var v72 base.V128
	_ = v72
	var v74 base.V128
	_ = v74
	var v76 base.V128
	_ = v76
	var v85 int32
	_ = v85
	var v87 base.V128
	_ = v87
	var v88 int32
	_ = v88
	var v90 base.V128
	_ = v90
	var v92 base.V128
	_ = v92
	var v94 base.V128
	_ = v94
	var v96 base.V128
	_ = v96
	var v98 base.V128
	_ = v98
	var v101 base.V128
	_ = v101
	var v105 base.V128
	_ = v105
	var v108 base.V128
	_ = v108
	var v110 base.V128
	_ = v110
	var v111 int32
	_ = v111
	var v113 base.V128
	_ = v113
	var v114 int32
	_ = v114
	var v116 base.V128
	_ = v116
	var v118 base.V128
	_ = v118
	var v131 base.V128
	_ = v131
	var v144 base.V128
	_ = v144
	var v145 base.V128
	_ = v145
	var v146 base.V128
	_ = v146
	var v147 base.V128
	_ = v147
	var v148 base.V128
	_ = v148
	var v149 base.V128
	_ = v149
	var v153 base.V128
	_ = v153
	var v158 base.V128
	_ = v158
	var v159 base.V128
	_ = v159
	var v160 base.V128
	_ = v160
	var v161 base.V128
	_ = v161
	var v163 base.V128
	_ = v163
	var v167 base.V128
	_ = v167
	var v168 base.V128
	_ = v168
	var v169 base.V128
	_ = v169
	var v170 base.V128
	_ = v170
	var v171 int32
	_ = v171
	var v173 base.V128
	_ = v173
	var v174 base.V128
	_ = v174
	var v180 base.V128
	_ = v180
	var v181 base.V128
	_ = v181
	var v184 base.V128
	_ = v184
	var v187 base.V128
	_ = v187
	var v192 base.V128
	_ = v192
	var v198 base.V128
	_ = v198
	var v201 base.V128
	_ = v201
	var v204 base.V128
	_ = v204
	var v206 base.V128
	_ = v206
	var v218 base.V128
	_ = v218
	var v220 base.V128
	_ = v220
	var v223 int32
	_ = v223
	var v237 base.V128
	_ = v237
	var v239 base.V128
	_ = v239
	var v250 base.V128
	_ = v250
	var v262 base.V128
	_ = v262
	var v272 base.V128
	_ = v272
	var v283 base.V128
	_ = v283
	v7 = int32(0)
	v11 = base.Simd_g_const(&F_VFilter8_SSE2__k0)
	v31 = l2 * int32(-3)
	v34 = l2 << (uint(int32(2)) % 32)
	v35 = l0 - v34
	v38 = base.Simd_g_v128_load64_zero(m, v35+l2, v7)
	v39 = l1 - v34
	v42 = base.Simd_g_v128_load64_zero(m, v39+l2, v7)
	v43 = base.Simd_g_const(&F_VFilter8_SSE2__k1)
	v44 = base.Simd_g_i8x16_shuffle2(v38, v42, base.Simd_g_const(&F_VFilter8_SSE2__k2), base.Simd_g_const(&F_VFilter8_SSE2__k3))
	v45 = base.Simd_g_const(&F_VFilter8_SSE2__k4)
	v48 = int32(1)
	v49 = l2 << (uint(v48) % 32)
	v52 = base.Simd_g_v128_load64_zero(m, v35+v49, v7)
	v55 = base.Simd_g_v128_load64_zero(m, v39+v49, v7)
	v57 = base.Simd_g_i8x16_shuffle2(v52, v55, base.Simd_g_const(&F_VFilter8_SSE2__k2), base.Simd_g_const(&F_VFilter8_SSE2__k3))
	v59 = l2 * int32(3)
	v62 = base.Simd_g_v128_load64_zero(m, v35+v59, v7)
	v65 = base.Simd_g_v128_load64_zero(m, v39+v59, v7)
	v67 = base.Simd_g_i8x16_shuffle2(v62, v65, base.Simd_g_const(&F_VFilter8_SSE2__k2), base.Simd_g_const(&F_VFilter8_SSE2__k3))
	v70 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v57, v67), base.Simd_g_i8x16_sub_sat_u(v67, v57))
	v72 = base.Simd_g_v128_load64_zero(m, v35, v7)
	v74 = base.Simd_g_v128_load64_zero(m, v39, v7)
	v76 = base.Simd_g_i8x16_shuffle2(v72, v74, base.Simd_g_const(&F_VFilter8_SSE2__k2), base.Simd_g_const(&F_VFilter8_SSE2__k3))
	v85 = l0 + l2
	v87 = base.Simd_g_v128_load64_zero(m, v85, v7)
	v88 = l1 + l2
	v90 = base.Simd_g_v128_load64_zero(m, v88, v7)
	v92 = base.Simd_g_i8x16_shuffle2(v87, v90, base.Simd_g_const(&F_VFilter8_SSE2__k2), base.Simd_g_const(&F_VFilter8_SSE2__k3))
	v94 = base.Simd_g_v128_load64_zero(m, l0, v7)
	v96 = base.Simd_g_v128_load64_zero(m, l1, v7)
	v98 = base.Simd_g_i8x16_shuffle2(v94, v96, base.Simd_g_const(&F_VFilter8_SSE2__k2), base.Simd_g_const(&F_VFilter8_SSE2__k3))
	v101 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v92, v98), base.Simd_g_i8x16_sub_sat_u(v98, v92))
	v105 = base.Simd_g_v128_load64_zero(m, l0+v59, v7)
	v108 = base.Simd_g_v128_load64_zero(m, l1+v59, v7)
	v110 = base.Simd_g_i8x16_shuffle2(v105, v108, base.Simd_g_const(&F_VFilter8_SSE2__k2), base.Simd_g_const(&F_VFilter8_SSE2__k3))
	v111 = l0 + v49
	v113 = base.Simd_g_v128_load64_zero(m, v111, v7)
	v114 = l1 + v49
	v116 = base.Simd_g_v128_load64_zero(m, v114, v7)
	v118 = base.Simd_g_i8x16_shuffle2(v113, v116, base.Simd_g_const(&F_VFilter8_SSE2__k2), base.Simd_g_const(&F_VFilter8_SSE2__k3))
	v131 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v67, v98), base.Simd_g_i8x16_sub_sat_u(v98, v67))
	v144 = base.Simd_g_i8x16_eq(v11, base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(v70, base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v76, v44), base.Simd_g_i8x16_sub_sat_u(v44, v76))), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v44, v57), base.Simd_g_i8x16_sub_sat_u(v57, v44))), v101), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v110, v118), base.Simd_g_i8x16_sub_sat_u(v118, v110))), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v118, v92), base.Simd_g_i8x16_sub_sat_u(v92, v118))), base.Simd_g_i8x16_splat(l4)), base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_add_sat_u(base.Simd_g_i8x16_add_sat_u(v131, v131), base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v57, v92), base.Simd_g_i8x16_sub_sat_u(v92, v57)), v48), base.Simd_g_const(&F_VFilter8_SSE2__k5))), base.Simd_g_i8x16_splat(l3))))
	v145 = base.Simd_g_v128_xor(v98, v45)
	v146 = base.Simd_g_v128_xor(v67, v45)
	v147 = base.Simd_g_i8x16_sub_sat_s(v145, v146)
	v148 = base.Simd_g_v128_xor(v57, v45)
	v149 = base.Simd_g_v128_xor(v92, v45)
	v153 = base.Simd_g_i8x16_add_sat_s(v147, base.Simd_g_i8x16_add_sat_s(v147, base.Simd_g_i8x16_add_sat_s(v147, base.Simd_g_i8x16_sub_sat_s(v148, v149))))
	v158 = base.Simd_g_i8x16_eq(v11, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_max_u(v70, v101), base.Simd_g_i8x16_splat(l5)))
	v159 = base.Simd_g_v128_and(base.Simd_g_v128_and(v144, v153), v158)
	v160 = base.Simd_g_const(&F_VFilter8_SSE2__k6)
	v161 = base.Simd_g_i8x16_shuffle2(v11, v159, base.Simd_g_const(&F_VFilter8_SSE2__k7), base.Simd_g_const(&F_VFilter8_SSE2__k8))
	v163 = base.Simd_g_const(&F_VFilter8_SSE2__k9)
	v167 = base.Simd_g_const(&F_VFilter8_SSE2__k10)
	v168 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v161), v163), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v161), v163), base.Simd_g_const(&F_VFilter8_SSE2__k11), base.Simd_g_const(&F_VFilter8_SSE2__k12))
	v169 = base.Simd_g_const(&F_VFilter8_SSE2__k13)
	v170 = base.Simd_g_i16x8_add(v168, v169)
	v171 = int32(7)
	v173 = base.Simd_g_const(&F_VFilter8_SSE2__k14)
	v174 = base.Simd_g_i8x16_shuffle2(v11, v159, base.Simd_g_const(&F_VFilter8_SSE2__k15), base.Simd_g_const(&F_VFilter8_SSE2__k16))
	v180 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v174), v163), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v174), v163), base.Simd_g_const(&F_VFilter8_SSE2__k11), base.Simd_g_const(&F_VFilter8_SSE2__k12))
	v181 = base.Simd_g_i16x8_add(v180, v169)
	v184 = base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(v170, v171), base.Simd_g_i16x8_shr_s(v181, v171))
	v187 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(base.Simd_g_v128_xor(v44, v45), v184), v45)
	base.Simd_g_v128_store64_lane_l0(m, l0+v31, v7, v187)
	v192 = base.Simd_g_const(&F_VFilter8_SSE2__k17)
	base.Simd_g_v128_store64_lane_l0(m, l1+v31, v7, base.Simd_g_i8x16_shuffle2(v187, v11, base.Simd_g_const(&F_VFilter8_SSE2__k18), base.Simd_g_const(&F_VFilter8_SSE2__k3)))
	v198 = base.Simd_g_i16x8_add(v170, v168)
	v201 = base.Simd_g_i16x8_add(v181, v180)
	v204 = base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(v198, v171), base.Simd_g_i16x8_shr_s(v201, v171))
	v206 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v148, v204), v45)
	base.Simd_g_v128_store64_lane_l0(m, l0-v49, v7, v206)
	base.Simd_g_v128_store64_lane_l0(m, l1-v49, v7, base.Simd_g_i8x16_shuffle2(v206, v11, base.Simd_g_const(&F_VFilter8_SSE2__k18), base.Simd_g_const(&F_VFilter8_SSE2__k3)))
	v218 = base.Simd_g_v128_and(base.Simd_g_v128_andnot(v153, v158), v144)
	v220 = base.Simd_g_i8x16_add_sat_s(v218, base.Simd_g_const(&F_VFilter8_SSE2__k19))
	v223 = int32(11)
	v237 = base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(v198, v168), v171), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(v201, v180), v171))
	v239 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(base.Simd_g_i8x16_add_sat_s(v146, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v11, v220, base.Simd_g_const(&F_VFilter8_SSE2__k7), base.Simd_g_const(&F_VFilter8_SSE2__k8)), v223), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v11, v220, base.Simd_g_const(&F_VFilter8_SSE2__k15), base.Simd_g_const(&F_VFilter8_SSE2__k16)), v223))), v237), v45)
	base.Simd_g_v128_store64_lane_l0(m, l0-l2, v7, v239)
	base.Simd_g_v128_store64_lane_l0(m, l1-l2, v7, base.Simd_g_i8x16_shuffle2(v239, v11, base.Simd_g_const(&F_VFilter8_SSE2__k18), base.Simd_g_const(&F_VFilter8_SSE2__k3)))
	v250 = base.Simd_g_i8x16_add_sat_s(v218, base.Simd_g_const(&F_VFilter8_SSE2__k20))
	v262 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(base.Simd_g_i8x16_sub_sat_s(v145, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v11, v250, base.Simd_g_const(&F_VFilter8_SSE2__k7), base.Simd_g_const(&F_VFilter8_SSE2__k8)), v223), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v11, v250, base.Simd_g_const(&F_VFilter8_SSE2__k15), base.Simd_g_const(&F_VFilter8_SSE2__k16)), v223))), v237), v45)
	base.Simd_g_v128_store64_lane_l0(m, l0, v7, v262)
	base.Simd_g_v128_store64_lane_l0(m, l1, v7, base.Simd_g_i8x16_shuffle2(v262, v11, base.Simd_g_const(&F_VFilter8_SSE2__k18), base.Simd_g_const(&F_VFilter8_SSE2__k3)))
	v272 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v149, v204), v45)
	base.Simd_g_v128_store64_lane_l0(m, v85, v7, v272)
	base.Simd_g_v128_store64_lane_l0(m, v88, v7, base.Simd_g_i8x16_shuffle2(v272, v11, base.Simd_g_const(&F_VFilter8_SSE2__k18), base.Simd_g_const(&F_VFilter8_SSE2__k3)))
	v283 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(base.Simd_g_v128_xor(v118, v45), v184), v45)
	base.Simd_g_v128_store64_lane_l0(m, v111, v7, v283)
	base.Simd_g_v128_store64_lane_l0(m, v114, v7, base.Simd_g_i8x16_shuffle2(v283, v11, base.Simd_g_const(&F_VFilter8_SSE2__k18), base.Simd_g_const(&F_VFilter8_SSE2__k3)))
	return
}

var F_VFilter8_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_VFilter8_SSE2__k1 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_VFilter8_SSE2__k2 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_VFilter8_SSE2__k3 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_VFilter8_SSE2__k4 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_VFilter8_SSE2__k5 = [2]uint64{0x7f7f7f7f7f7f7f7f, 0x7f7f7f7f7f7f7f7f}
var F_VFilter8_SSE2__k6 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_VFilter8_SSE2__k7 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_VFilter8_SSE2__k8 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_VFilter8_SSE2__k9 = [2]uint64{0x90000000900, 0x90000000900}
var F_VFilter8_SSE2__k10 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_VFilter8_SSE2__k11 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_VFilter8_SSE2__k12 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_VFilter8_SSE2__k13 = [2]uint64{0x3f003f003f003f, 0x3f003f003f003f}
var F_VFilter8_SSE2__k14 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_VFilter8_SSE2__k15 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_VFilter8_SSE2__k16 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_VFilter8_SSE2__k17 = [2]uint64{0xf0e0d0c0b0a0908, 0x1716151413121110}
var F_VFilter8_SSE2__k18 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_VFilter8_SSE2__k19 = [2]uint64{0x303030303030303, 0x303030303030303}
var F_VFilter8_SSE2__k20 = [2]uint64{0x404040404040404, 0x404040404040404}

func F_VFilter8i_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	var v7 int32
	_ = v7
	var v10 base.V128
	_ = v10
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v35 base.V128
	_ = v35
	var v38 base.V128
	_ = v38
	var v39 base.V128
	_ = v39
	var v40 base.V128
	_ = v40
	var v41 base.V128
	_ = v41
	var v42 base.V128
	_ = v42
	var v44 int32
	_ = v44
	var v46 base.V128
	_ = v46
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v50 base.V128
	_ = v50
	var v52 base.V128
	_ = v52
	var v53 base.V128
	_ = v53
	var v56 int32
	_ = v56
	var v59 base.V128
	_ = v59
	var v62 base.V128
	_ = v62
	var v64 base.V128
	_ = v64
	var v67 base.V128
	_ = v67
	var v69 base.V128
	_ = v69
	var v71 base.V128
	_ = v71
	var v73 base.V128
	_ = v73
	var v76 base.V128
	_ = v76
	var v80 base.V128
	_ = v80
	var v82 base.V128
	_ = v82
	var v83 base.V128
	_ = v83
	var v84 base.V128
	_ = v84
	var v89 base.V128
	_ = v89
	var v91 base.V128
	_ = v91
	var v93 base.V128
	_ = v93
	var v96 base.V128
	_ = v96
	var v99 base.V128
	_ = v99
	var v101 base.V128
	_ = v101
	var v113 base.V128
	_ = v113
	var v116 base.V128
	_ = v116
	var v118 base.V128
	_ = v118
	var v121 base.V128
	_ = v121
	var v124 base.V128
	_ = v124
	var v126 base.V128
	_ = v126
	var v139 base.V128
	_ = v139
	var v153 base.V128
	_ = v153
	var v155 base.V128
	_ = v155
	var v156 base.V128
	_ = v156
	var v158 int32
	_ = v158
	var v160 base.V128
	_ = v160
	var v164 base.V128
	_ = v164
	var v169 base.V128
	_ = v169
	var v172 base.V128
	_ = v172
	var v177 base.V128
	_ = v177
	var v184 base.V128
	_ = v184
	var v195 base.V128
	_ = v195
	var v206 base.V128
	_ = v206
	var v216 base.V128
	_ = v216
	v7 = int32(0)
	v10 = base.Simd_g_const(&F_VFilter8i_SSE2__k0)
	v28 = l2 << (uint(int32(2)) % 32)
	v29 = l0 + v28
	v30 = int32(1)
	v31 = l2 << (uint(v30) % 32)
	v35 = base.Simd_g_v128_load64_zero(m, l0+v31, v7)
	v38 = base.Simd_g_v128_load64_zero(m, l1+v31, v7)
	v39 = base.Simd_g_const(&F_VFilter8i_SSE2__k1)
	v40 = base.Simd_g_i8x16_shuffle2(v35, v38, base.Simd_g_const(&F_VFilter8i_SSE2__k2), base.Simd_g_const(&F_VFilter8i_SSE2__k3))
	v41 = base.Simd_g_const(&F_VFilter8i_SSE2__k4)
	v42 = base.Simd_g_v128_xor(v40, v41)
	v44 = v29 + l2
	v46 = base.Simd_g_v128_load64_zero(m, v44, v7)
	v47 = l1 + v28
	v48 = v47 + l2
	v50 = base.Simd_g_v128_load64_zero(m, v48, v7)
	v52 = base.Simd_g_i8x16_shuffle2(v46, v50, base.Simd_g_const(&F_VFilter8i_SSE2__k2), base.Simd_g_const(&F_VFilter8i_SSE2__k3))
	v53 = base.Simd_g_v128_xor(v52, v41)
	v56 = l2 * int32(3)
	v59 = base.Simd_g_v128_load64_zero(m, l0+v56, v7)
	v62 = base.Simd_g_v128_load64_zero(m, l1+v56, v7)
	v64 = base.Simd_g_i8x16_shuffle2(v59, v62, base.Simd_g_const(&F_VFilter8i_SSE2__k2), base.Simd_g_const(&F_VFilter8i_SSE2__k3))
	v67 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v40, v64), base.Simd_g_i8x16_sub_sat_u(v64, v40))
	v69 = base.Simd_g_v128_load64_zero(m, v29, v7)
	v71 = base.Simd_g_v128_load64_zero(m, v47, v7)
	v73 = base.Simd_g_i8x16_shuffle2(v69, v71, base.Simd_g_const(&F_VFilter8i_SSE2__k2), base.Simd_g_const(&F_VFilter8i_SSE2__k3))
	v76 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v52, v73), base.Simd_g_i8x16_sub_sat_u(v73, v52))
	v80 = base.Simd_g_i8x16_eq(v10, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_max_u(v67, v76), base.Simd_g_i8x16_splat(l5)))
	v82 = base.Simd_g_v128_xor(v73, v41)
	v83 = base.Simd_g_v128_xor(v64, v41)
	v84 = base.Simd_g_i8x16_sub_sat_s(v82, v83)
	v89 = base.Simd_g_v128_load64_zero(m, l0, v7)
	v91 = base.Simd_g_v128_load64_zero(m, l1, v7)
	v93 = base.Simd_g_i8x16_shuffle2(v89, v91, base.Simd_g_const(&F_VFilter8i_SSE2__k2), base.Simd_g_const(&F_VFilter8i_SSE2__k3))
	v96 = base.Simd_g_v128_load64_zero(m, l0+l2, v7)
	v99 = base.Simd_g_v128_load64_zero(m, l1+l2, v7)
	v101 = base.Simd_g_i8x16_shuffle2(v96, v99, base.Simd_g_const(&F_VFilter8i_SSE2__k2), base.Simd_g_const(&F_VFilter8i_SSE2__k3))
	v113 = base.Simd_g_v128_load64_zero(m, v29+v56, v7)
	v116 = base.Simd_g_v128_load64_zero(m, v47+v56, v7)
	v118 = base.Simd_g_i8x16_shuffle2(v113, v116, base.Simd_g_const(&F_VFilter8i_SSE2__k2), base.Simd_g_const(&F_VFilter8i_SSE2__k3))
	v121 = base.Simd_g_v128_load64_zero(m, v29+v31, v7)
	v124 = base.Simd_g_v128_load64_zero(m, v47+v31, v7)
	v126 = base.Simd_g_i8x16_shuffle2(v121, v124, base.Simd_g_const(&F_VFilter8i_SSE2__k2), base.Simd_g_const(&F_VFilter8i_SSE2__k3))
	v139 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v64, v73), base.Simd_g_i8x16_sub_sat_u(v73, v64))
	v153 = base.Simd_g_v128_bitselect(base.Simd_g_i8x16_add_sat_s(base.Simd_g_i8x16_add_sat_s(base.Simd_g_i8x16_add_sat_s(base.Simd_g_v128_bitselect(v10, base.Simd_g_i8x16_sub_sat_s(v42, v53), v80), v84), v84), v84), v10, base.Simd_g_i8x16_eq(v10, base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(v67, base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v93, v101), base.Simd_g_i8x16_sub_sat_u(v101, v93))), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v101, v40), base.Simd_g_i8x16_sub_sat_u(v40, v101))), v76), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v118, v126), base.Simd_g_i8x16_sub_sat_u(v126, v118))), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v126, v52), base.Simd_g_i8x16_sub_sat_u(v52, v126))), base.Simd_g_i8x16_splat(l4)), base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_add_sat_u(base.Simd_g_i8x16_add_sat_u(v139, v139), base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v40, v52), base.Simd_g_i8x16_sub_sat_u(v52, v40)), v30), base.Simd_g_const(&F_VFilter8i_SSE2__k5))), base.Simd_g_i8x16_splat(l3)))))
	v155 = base.Simd_g_i8x16_add_sat_s(v153, base.Simd_g_const(&F_VFilter8i_SSE2__k6))
	v156 = base.Simd_g_const(&F_VFilter8i_SSE2__k7)
	v158 = int32(11)
	v160 = base.Simd_g_const(&F_VFilter8i_SSE2__k8)
	v164 = base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v10, v155, base.Simd_g_const(&F_VFilter8i_SSE2__k9), base.Simd_g_const(&F_VFilter8i_SSE2__k10)), v158), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v10, v155, base.Simd_g_const(&F_VFilter8i_SSE2__k11), base.Simd_g_const(&F_VFilter8i_SSE2__k12)), v158))
	v169 = base.Simd_g_v128_bitselect(base.Simd_g_i8x16_add(base.Simd_g_i8x16_avgr_u(base.Simd_g_v128_xor(v164, v41), v10), base.Simd_g_const(&F_VFilter8i_SSE2__k13)), v10, v80)
	v172 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v42, v169), v41)
	base.Simd_g_v128_store64_lane_l0(m, v29-v31, v7, v172)
	v177 = base.Simd_g_const(&F_VFilter8i_SSE2__k14)
	base.Simd_g_v128_store64_lane_l0(m, v47-v31, v7, base.Simd_g_i8x16_shuffle2(v172, v10, base.Simd_g_const(&F_VFilter8i_SSE2__k15), base.Simd_g_const(&F_VFilter8i_SSE2__k3)))
	v184 = base.Simd_g_i8x16_add_sat_s(v153, base.Simd_g_const(&F_VFilter8i_SSE2__k16))
	v195 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v83, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v10, v184, base.Simd_g_const(&F_VFilter8i_SSE2__k9), base.Simd_g_const(&F_VFilter8i_SSE2__k10)), v158), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v10, v184, base.Simd_g_const(&F_VFilter8i_SSE2__k11), base.Simd_g_const(&F_VFilter8i_SSE2__k12)), v158))), v41)
	base.Simd_g_v128_store64_lane_l0(m, v29-l2, v7, v195)
	base.Simd_g_v128_store64_lane_l0(m, v47-l2, v7, base.Simd_g_i8x16_shuffle2(v195, v10, base.Simd_g_const(&F_VFilter8i_SSE2__k15), base.Simd_g_const(&F_VFilter8i_SSE2__k3)))
	v206 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v82, v164), v41)
	base.Simd_g_v128_store64_lane_l0(m, v29, v7, v206)
	base.Simd_g_v128_store64_lane_l0(m, v47, v7, base.Simd_g_i8x16_shuffle2(v206, v10, base.Simd_g_const(&F_VFilter8i_SSE2__k15), base.Simd_g_const(&F_VFilter8i_SSE2__k3)))
	v216 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v53, v169), v41)
	base.Simd_g_v128_store64_lane_l0(m, v44, v7, v216)
	base.Simd_g_v128_store64_lane_l0(m, v48, v7, base.Simd_g_i8x16_shuffle2(v216, v10, base.Simd_g_const(&F_VFilter8i_SSE2__k15), base.Simd_g_const(&F_VFilter8i_SSE2__k3)))
	return
}

var F_VFilter8i_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_VFilter8i_SSE2__k1 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_VFilter8i_SSE2__k2 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_VFilter8i_SSE2__k3 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_VFilter8i_SSE2__k4 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_VFilter8i_SSE2__k5 = [2]uint64{0x7f7f7f7f7f7f7f7f, 0x7f7f7f7f7f7f7f7f}
var F_VFilter8i_SSE2__k6 = [2]uint64{0x404040404040404, 0x404040404040404}
var F_VFilter8i_SSE2__k7 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_VFilter8i_SSE2__k8 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_VFilter8i_SSE2__k9 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_VFilter8i_SSE2__k10 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_VFilter8i_SSE2__k11 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_VFilter8i_SSE2__k12 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_VFilter8i_SSE2__k13 = [2]uint64{0xc0c0c0c0c0c0c0c0, 0xc0c0c0c0c0c0c0c0}
var F_VFilter8i_SSE2__k14 = [2]uint64{0xf0e0d0c0b0a0908, 0x1716151413121110}
var F_VFilter8i_SSE2__k15 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_VFilter8i_SSE2__k16 = [2]uint64{0x303030303030303, 0x303030303030303}

func F_VL4_SSE2(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 base.V128
	_ = v2
	var v8 int32
	_ = v8
	var v12 base.V128
	_ = v12
	var v14 base.V128
	_ = v14
	var v15 base.V128
	_ = v15
	var v16 base.V128
	_ = v16
	var v23 base.V128
	_ = v23
	var v24 base.V128
	_ = v24
	var v33 base.V128
	_ = v33
	var v48 int32
	_ = v48
	var v51 int32
	_ = v51
	v2 = base.Simd_g_const(&F_VL4_SSE2__k0)
	v8 = int32(0)
	v12 = base.Simd_g_v128_load64_zero(m, l0+int32(-32), v8)
	v14 = base.Simd_g_const(&F_VL4_SSE2__k1)
	v15 = base.Simd_g_i8x16_shuffle2(v12, v2, base.Simd_g_const(&F_VL4_SSE2__k2), base.Simd_g_const(&F_VL4_SSE2__k3))
	v16 = base.Simd_g_i8x16_avgr_u(v12, v15)
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(64), base.Simd_g_i8x16_shuffle2(v16, v2, base.Simd_g_const(&F_VL4_SSE2__k2), base.Simd_g_const(&F_VL4_SSE2__k3)))
	v23 = base.Simd_g_i8x16_shuffle2(v12, v2, base.Simd_g_const(&F_VL4_SSE2__k4), base.Simd_g_const(&F_VL4_SSE2__k5))
	v24 = base.Simd_g_i8x16_avgr_u(v23, v15)
	v33 = base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_avgr_u(v16, v24), base.Simd_g_v128_and(base.Simd_g_v128_and(base.Simd_g_v128_xor(v24, v16), base.Simd_g_v128_or(base.Simd_g_v128_xor(v15, v12), base.Simd_g_v128_xor(v15, v23))), base.Simd_g_const(&F_VL4_SSE2__k6)))
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(96), base.Simd_g_i8x16_shuffle2(v33, v2, base.Simd_g_const(&F_VL4_SSE2__k2), base.Simd_g_const(&F_VL4_SSE2__k3)))
	base.Simd_g_v128_store32_lane_l0(m, l0, v8, v16)
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(32), v33)
	v48 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i8x16_shuffle2(v33, v2, base.Simd_g_const(&F_VL4_SSE2__k7), base.Simd_g_const(&F_VL4_SSE2__k8)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+67)) = uint8(v48)
	v51 = int32(base.Ui32(v48) >> (uint(int32(8)) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+99)) = uint8(v51)
	return
}

var F_VL4_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_VL4_SSE2__k1 = [2]uint64{0x807060504030201, 0x100f0e0d0c0b0a09}
var F_VL4_SSE2__k2 = [2]uint64{0x807060504030201, 0x800f0e0d0c0b0a09}
var F_VL4_SSE2__k3 = [2]uint64{0x8080808080808080, 0x80808080808080}
var F_VL4_SSE2__k4 = [2]uint64{0x908070605040302, 0x80800f0e0d0c0b0a}
var F_VL4_SSE2__k5 = [2]uint64{0x8080808080808080, 0x100808080808080}
var F_VL4_SSE2__k6 = [2]uint64{0x101010101010101, 0x101010101010101}
var F_VL4_SSE2__k7 = [2]uint64{0xb0a090807060504, 0x808080800f0e0d0c}
var F_VL4_SSE2__k8 = [2]uint64{0x8080808080808080, 0x302010080808080}

func F_VP8LAddGreenToBlueAndRed_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v13 int32
	_ = v13
	var v20 int32
	_ = v20
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v31 int32
	_ = v31
	var v32 base.V128
	_ = v32
	var v34 base.V128
	_ = v34
	var v37 base.V128
	_ = v37
	var v40 int32
	_ = v40
	var v46 base.V128
	_ = v46
	var v54 int32
	_ = v54
	var v59 int32
	_ = v59
	var v71 int32
	_ = v71
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v79 int32
	_ = v79
	var v93 int32
	_ = v93
	var v100 int32
	_ = v100
	var v104 int32
	_ = v104
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v116 int32
	_ = v116
	var v118 int32
	_ = v118
	var v121 int32
	_ = v121
	var v126 int32
	_ = v126
	var v130 int32
	_ = v130
	var v134 int32
	_ = v134
	var v136 int32
	_ = v136
	var v156 int32
	_ = v156
	if l1 < int32(1) {
	} else {
		v13 = int32(0)
		if base.Ui32(l1) < base.Ui32(int32(4)) {
			v59 = v13
			if l1&int32(1) == int32(0) {
				v93 = v59
			} else {
				v71 = v59 << (uint(int32(2)) % 32)
				v74 = *(*int32)(unsafe.Add(mBase, uint32(l0+v71)))
				v76 = int32(base.Ui32(v74) >> (uint(int32(8)) % 32))
				v79 = int32(16711935)
				*(*int32)(unsafe.Add(mBase, uint32(l2+v71))) = (v76&int32(255)+v74&v79+v76<<(uint(int32(16))%32))&v79 | v74&int32(-16711936)
				v93 = v59 | int32(1)
			}
			if v59 == l1+int32(-1) {
			} else {
				v100 = v93 << (uint(int32(2)) % 32)
				v104 = l1 - v93
				v108 = l0 + v100
				v109 = l2 + v100
				for {
					v113 = *(*int32)(unsafe.Add(mBase, uint32(v108)))
					v114 = int32(8)
					v115 = int32(base.Ui32(v113) >> (uint(v114) % 32))
					v116 = int32(255)
					v118 = int32(16711935)
					v121 = int32(16)
					v126 = int32(-16711936)
					*(*int32)(unsafe.Add(mBase, uint32(v109))) = (v115&v116+v113&v118+v115<<(uint(v121)%32))&v118 | v113&v126
					v130 = int32(4)
					v134 = *(*int32)(unsafe.Add(mBase, uint32(v108+v130)))
					v136 = int32(base.Ui32(v134) >> (uint(v114) % 32))
					*(*int32)(unsafe.Add(mBase, uint32(v109+v130))) = (v136&v116+v134&v118+v136<<(uint(v121)%32))&v118 | v134&v126
					v156 = v104 + int32(-2)
					if v156 != 0 {
						v104 = v156
						v108 = v108 + v114
						v109 = v109 + v114
						continue
					} else {
						break
					}
					break
				}
			}
		} else {
			if base.Ui32(l2-l0) < base.Ui32(int32(16)) {
				v59 = v13
				if l1&int32(1) == int32(0) {
					v93 = v59
				} else {
					v71 = v59 << (uint(int32(2)) % 32)
					v74 = *(*int32)(unsafe.Add(mBase, uint32(l0+v71)))
					v76 = int32(base.Ui32(v74) >> (uint(int32(8)) % 32))
					v79 = int32(16711935)
					*(*int32)(unsafe.Add(mBase, uint32(l2+v71))) = (v76&int32(255)+v74&v79+v76<<(uint(int32(16))%32))&v79 | v74&int32(-16711936)
					v93 = v59 | int32(1)
				}
				if v59 == l1+int32(-1) {
				} else {
					v100 = v93 << (uint(int32(2)) % 32)
					v104 = l1 - v93
					v108 = l0 + v100
					v109 = l2 + v100
					for {
						v113 = *(*int32)(unsafe.Add(mBase, uint32(v108)))
						v114 = int32(8)
						v115 = int32(base.Ui32(v113) >> (uint(v114) % 32))
						v116 = int32(255)
						v118 = int32(16711935)
						v121 = int32(16)
						v126 = int32(-16711936)
						*(*int32)(unsafe.Add(mBase, uint32(v109))) = (v115&v116+v113&v118+v115<<(uint(v121)%32))&v118 | v113&v126
						v130 = int32(4)
						v134 = *(*int32)(unsafe.Add(mBase, uint32(v108+v130)))
						v136 = int32(base.Ui32(v134) >> (uint(v114) % 32))
						*(*int32)(unsafe.Add(mBase, uint32(v109+v130))) = (v136&v116+v134&v118+v136<<(uint(v121)%32))&v118 | v134&v126
						v156 = v104 + int32(-2)
						if v156 != 0 {
							v104 = v156
							v108 = v108 + v114
							v109 = v109 + v114
							continue
						} else {
							break
						}
						break
					}
				}
			} else {
				v20 = l1 & int32(2147483644)
				v25 = v20
				v26 = l2
				v27 = l0
				for {
					v31 = int32(0)
					v32 = base.Simd_g_v128_load(m, v27, v31)
					v34 = base.Simd_g_i32x4_shr_u(v32, int32(8))
					v37 = base.Simd_g_const(&F_VP8LAddGreenToBlueAndRed_C__k0)
					v40 = int32(16)
					v46 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v34, base.Simd_g_const(&F_VP8LAddGreenToBlueAndRed_C__k1)), base.Simd_g_v128_and(v32, v37)), base.Simd_g_i32x4_shl(v34, v40)), v37), base.Simd_g_v128_and(v32, base.Simd_g_const(&F_VP8LAddGreenToBlueAndRed_C__k2)))
					base.Simd_g_v128_store(m, v26, v31, v46)
					v54 = v25 + int32(-4)
					if v54 != 0 {
						v25 = v54
						v26 = v26 + v40
						v27 = v27 + v40
						continue
					} else {
						break
					}
					break
				}
				if v20 == l1 {
				} else {
					v59 = v20
					if l1&int32(1) == int32(0) {
						v93 = v59
					} else {
						v71 = v59 << (uint(int32(2)) % 32)
						v74 = *(*int32)(unsafe.Add(mBase, uint32(l0+v71)))
						v76 = int32(base.Ui32(v74) >> (uint(int32(8)) % 32))
						v79 = int32(16711935)
						*(*int32)(unsafe.Add(mBase, uint32(l2+v71))) = (v76&int32(255)+v74&v79+v76<<(uint(int32(16))%32))&v79 | v74&int32(-16711936)
						v93 = v59 | int32(1)
					}
					if v59 == l1+int32(-1) {
					} else {
						v100 = v93 << (uint(int32(2)) % 32)
						v104 = l1 - v93
						v108 = l0 + v100
						v109 = l2 + v100
						for {
							v113 = *(*int32)(unsafe.Add(mBase, uint32(v108)))
							v114 = int32(8)
							v115 = int32(base.Ui32(v113) >> (uint(v114) % 32))
							v116 = int32(255)
							v118 = int32(16711935)
							v121 = int32(16)
							v126 = int32(-16711936)
							*(*int32)(unsafe.Add(mBase, uint32(v109))) = (v115&v116+v113&v118+v115<<(uint(v121)%32))&v118 | v113&v126
							v130 = int32(4)
							v134 = *(*int32)(unsafe.Add(mBase, uint32(v108+v130)))
							v136 = int32(base.Ui32(v134) >> (uint(v114) % 32))
							*(*int32)(unsafe.Add(mBase, uint32(v109+v130))) = (v136&v116+v134&v118+v136<<(uint(v121)%32))&v118 | v134&v126
							v156 = v104 + int32(-2)
							if v156 != 0 {
								v104 = v156
								v108 = v108 + v114
								v109 = v109 + v114
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

var F_VP8LAddGreenToBlueAndRed_C__k0 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_VP8LAddGreenToBlueAndRed_C__k1 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_VP8LAddGreenToBlueAndRed_C__k2 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}

func F_VP8LBundleColorMap_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v20 int32
	_ = v20
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v32 int32
	_ = v32
	var v33 base.V128
	_ = v33
	var v39 base.V128
	_ = v39
	var v47 int32
	_ = v47
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v68 int32
	_ = v68
	var v73 int32
	_ = v73
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v84 int32
	_ = v84
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v93 int32
	_ = v93
	var v94 int32
	_ = v94
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v112 int32
	_ = v112
	var v115 int32
	_ = v115
	var v120 int32
	_ = v120
	var v132 int32
	_ = v132
	var v133 int32
	_ = v133
	var v139 int32
	_ = v139
	var v146 int32
	_ = v146
	var v160 int32
	_ = v160
	var v163 int32
	_ = v163
	var v164 int32
	_ = v164
	var v169 int32
	_ = v169
	var v180 int32
	_ = v180
	if int32(0) < l2 {
		if l1 < int32(1) {
		} else {
			v52 = int32(3) - l2
			v53 = int32(-1)
			v56 = v53<<(uint(l2)%32) ^ v53
			v57 = int32(1)
			if l1 != v57 {
				v68 = int32(-16777216)
				v73 = int32(0)
				for {
					v82 = l0 + v73
					v83 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v82))))
					v84 = v73 & v56
					if v84 != 0 {
						v90 = v68
					} else {
						v90 = int32(-16777216)
					}
					v91 = v83<<(uint(v84<<(uint(v52)%32)+int32(8))%32) | v90
					*(*int32)(unsafe.Add(mBase, uint32(l3+int32(base.Ui32(v73)>>(uint(l2)%32))<<(uint(int32(2))%32)))) = v91
					v93 = int32(1)
					v94 = v73 + v93
					v101 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v82+v93))))
					v102 = v94 & v56
					if v102 != 0 {
						v108 = v91
					} else {
						v108 = int32(-16777216)
					}
					v109 = v101<<(uint(v102<<(uint(v52)%32)+int32(8))%32) | v108
					*(*int32)(unsafe.Add(mBase, uint32(l3+int32(base.Ui32(v94)>>(uint(l2)%32))<<(uint(int32(2))%32)))) = v109
					v112 = v73 + int32(2)
					if v112 != l1&int32(2147483646) {
						v68 = v109
						v73 = v112
						continue
					} else {
						break
					}
					break
				}
				v115 = v109
				v120 = v112
			} else {
				v115 = int32(-16777216)
				v120 = int32(0)
			}
			if l1&v57 == int32(0) {
			} else {
				v132 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v120))))
				v133 = v120 & v56
				if v133 != 0 {
					v139 = v115
				} else {
					v139 = int32(-16777216)
				}
				*(*int32)(unsafe.Add(mBase, uint32(l3+int32(base.Ui32(v120)>>(uint(l2)%32))<<(uint(int32(2))%32)))) = v132<<(uint(v133<<(uint(v52)%32)+int32(8))%32) | v139
			}
		}
	} else {
		if l1 < int32(1) {
		} else {
			if base.Ui32(l1) <= base.Ui32(int32(3)) {
				v146 = int32(0)
				v160 = l3 + v146<<(uint(int32(2))%32)
				v163 = l1 - v146
				v164 = l0 + v146
				for {
					v169 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v164))))
					*(*int32)(unsafe.Add(mBase, uint32(v160))) = v169<<(uint(int32(8))%32) | int32(-16777216)
					v180 = v163 + int32(-1)
					if v180 != 0 {
						v160 = v160 + int32(4)
						v163 = v180
						v164 = v164 + int32(1)
						continue
					} else {
						break
					}
					break
				}
			} else {
				v20 = l1 & int32(2147483644)
				v23 = l0
				v26 = v20
				v27 = l3
				for {
					v32 = int32(0)
					v33 = base.Simd_g_v128_load32_zero(m, v23, v32)
					v39 = base.Simd_g_v128_or(base.Simd_g_i32x4_shl(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v33)), int32(8)), base.Simd_g_const(&F_VP8LBundleColorMap_C__k0))
					base.Simd_g_v128_store(m, v27, v32, v39)
					v47 = v26 + int32(-4)
					if v47 != 0 {
						v23 = v23 + int32(4)
						v26 = v47
						v27 = v27 + int32(16)
						continue
					} else {
						break
					}
					break
				}
				if v20 != l1 {
					v146 = v20
					v160 = l3 + v146<<(uint(int32(2))%32)
					v163 = l1 - v146
					v164 = l0 + v146
					for {
						v169 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v164))))
						*(*int32)(unsafe.Add(mBase, uint32(v160))) = v169<<(uint(int32(8))%32) | int32(-16777216)
						v180 = v163 + int32(-1)
						if v180 != 0 {
							v160 = v160 + int32(4)
							v163 = v180
							v164 = v164 + int32(1)
							continue
						} else {
							break
						}
						break
					}
				} else {
				}
			}
		}
	}
	return
}

var F_VP8LBundleColorMap_C__k0 = [2]uint64{0xff000000ff000000, 0xff000000ff000000}

func F_VP8LConvertBGRAToRGB565_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v17 int32
	_ = v17
	var v19 int32
	_ = v19
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v27 int32
	_ = v27
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v55 int32
	_ = v55
	var v57 int32
	_ = v57
	var v60 int32
	_ = v60
	var v67 int32
	_ = v67
	var v68 base.V128
	_ = v68
	var v69 int32
	_ = v69
	var v77 base.V128
	_ = v77
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v83 int32
	_ = v83
	var v84 int32
	_ = v84
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v92 int32
	_ = v92
	var v96 int32
	_ = v96
	var v104 base.V128
	_ = v104
	var v106 int32
	_ = v106
	var v111 int32
	_ = v111
	var v116 int32
	_ = v116
	var v121 int32
	_ = v121
	var v128 int32
	_ = v128
	var v132 int32
	_ = v132
	var v135 int32
	_ = v135
	var v142 int32
	_ = v142
	var v145 int32
	_ = v145
	var v151 int32
	_ = v151
	var v160 int32
	_ = v160
	var v170 int32
	_ = v170
	var v175 int32
	_ = v175
	if l1 < int32(1) {
	} else {
		v15 = l0 + l1<<(uint(int32(2))%32)
		v17 = l0 + int32(4)
		if base.Ui32(v17) < base.Ui32(v15) {
			v19 = v15
		} else {
			v19 = v17
		}
		v21 = l0 ^ int32(-1)
		v22 = v19 + v21
		if base.Ui32(v22) < base.Ui32(int32(44)) {
			v132 = l2
			v135 = l0
			v142 = v132
			v145 = v135
			for {
				v151 = *(*int32)(unsafe.Add(mBase, uint32(v145)))
				v160 = int32(base.Ui32(v151)>>(uint(int32(5))%32))&int32(224) | int32(base.Ui32(v151)>>(uint(int32(3))%32))&int32(31)
				*(*uint8)(unsafe.Add(mBase, uint32(v142)+1)) = uint8(v160)
				v170 = int32(base.Ui32(v151)>>(uint(int32(16))%32))&int32(248) | int32(base.Ui32(v151)>>(uint(int32(13))%32))&int32(7)
				*(*uint8)(unsafe.Add(mBase, uint32(v142))) = uint8(v170)
				v175 = v145 + int32(4)
				if base.Ui32(v175) < base.Ui32(v15) {
					v142 = v142 + int32(2)
					v145 = v175
					continue
				} else {
					break
				}
				break
			}
		} else {
			v27 = l1<<(uint(int32(2))%32) + l0
			if base.Ui32(v17) < base.Ui32(v27) {
				v29 = v27
			} else {
				v29 = v17
			}
			v30 = v29 + v21
			if base.Ui32(v17+v30&int32(-4)) <= base.Ui32(l2) {
				v43 = int32(2)
				v45 = int32(1)
				v46 = int32(base.Ui32(v22)>>(uint(v43)%32)) + v45
				v48 = v46 & int32(2147483644)
				v55 = l0
				v57 = l2
				v60 = v48
				for {
					v67 = int32(0)
					v68 = base.Simd_g_v128_load(m, v55, v67)
					v69 = int32(16)
					v77 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v68, v69), base.Simd_g_const(&F_VP8LConvertBGRAToRGB565_C__k0)), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v68, int32(13)), base.Simd_g_const(&F_VP8LConvertBGRAToRGB565_C__k1)))
					v78 = int32(3)
					v79 = base.Simd_g_i32x4_extract_lane_l3(v77)
					*(*uint8)(unsafe.Add(mBase, uint32(v57+int32(6)))) = uint8(v79)
					v83 = int32(2)
					v84 = base.Simd_g_i32x4_extract_lane_l2(v77)
					*(*uint8)(unsafe.Add(mBase, uint32(v57+int32(4)))) = uint8(v84)
					v88 = int32(1)
					v89 = base.Simd_g_i32x4_extract_lane_l1(v77)
					*(*uint8)(unsafe.Add(mBase, uint32(v57+v83))) = uint8(v89)
					v92 = base.Simd_g_i32x4_extract_lane_l0(v77)
					*(*uint8)(unsafe.Add(mBase, uint32(v57))) = uint8(v92)
					v96 = int32(5)
					v104 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v68, v96), base.Simd_g_const(&F_VP8LConvertBGRAToRGB565_C__k2)), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v68, v78), base.Simd_g_const(&F_VP8LConvertBGRAToRGB565_C__k3)))
					v106 = base.Simd_g_i32x4_extract_lane_l3(v104)
					*(*uint8)(unsafe.Add(mBase, uint32(v57+int32(7)))) = uint8(v106)
					v111 = base.Simd_g_i32x4_extract_lane_l2(v104)
					*(*uint8)(unsafe.Add(mBase, uint32(v57+v96))) = uint8(v111)
					v116 = base.Simd_g_i32x4_extract_lane_l1(v104)
					*(*uint8)(unsafe.Add(mBase, uint32(v57+v78))) = uint8(v116)
					v121 = base.Simd_g_i32x4_extract_lane_l0(v104)
					*(*uint8)(unsafe.Add(mBase, uint32(v57+v88))) = uint8(v121)
					v128 = v60 + int32(-4)
					if v128 != 0 {
						v55 = v55 + v69
						v57 = v57 + int32(8)
						v60 = v128
						continue
					} else {
						break
					}
					break
				}
				if v46 == v48 {
				} else {
					v132 = l2 + v48<<(uint(v45)%32)
					v135 = l0 + v48<<(uint(v43)%32)
					v142 = v132
					v145 = v135
					for {
						v151 = *(*int32)(unsafe.Add(mBase, uint32(v145)))
						v160 = int32(base.Ui32(v151)>>(uint(int32(5))%32))&int32(224) | int32(base.Ui32(v151)>>(uint(int32(3))%32))&int32(31)
						*(*uint8)(unsafe.Add(mBase, uint32(v142)+1)) = uint8(v160)
						v170 = int32(base.Ui32(v151)>>(uint(int32(16))%32))&int32(248) | int32(base.Ui32(v151)>>(uint(int32(13))%32))&int32(7)
						*(*uint8)(unsafe.Add(mBase, uint32(v142))) = uint8(v170)
						v175 = v145 + int32(4)
						if base.Ui32(v175) < base.Ui32(v15) {
							v142 = v142 + int32(2)
							v145 = v175
							continue
						} else {
							break
						}
						break
					}
				}
			} else {
				if base.Ui32(l0) < base.Ui32(l2+int32(base.Ui32(v30)>>(uint(int32(1))%32))&int32(2147483646)+int32(2)) {
					v132 = l2
					v135 = l0
					v142 = v132
					v145 = v135
					for {
						v151 = *(*int32)(unsafe.Add(mBase, uint32(v145)))
						v160 = int32(base.Ui32(v151)>>(uint(int32(5))%32))&int32(224) | int32(base.Ui32(v151)>>(uint(int32(3))%32))&int32(31)
						*(*uint8)(unsafe.Add(mBase, uint32(v142)+1)) = uint8(v160)
						v170 = int32(base.Ui32(v151)>>(uint(int32(16))%32))&int32(248) | int32(base.Ui32(v151)>>(uint(int32(13))%32))&int32(7)
						*(*uint8)(unsafe.Add(mBase, uint32(v142))) = uint8(v170)
						v175 = v145 + int32(4)
						if base.Ui32(v175) < base.Ui32(v15) {
							v142 = v142 + int32(2)
							v145 = v175
							continue
						} else {
							break
						}
						break
					}
				} else {
					v43 = int32(2)
					v45 = int32(1)
					v46 = int32(base.Ui32(v22)>>(uint(v43)%32)) + v45
					v48 = v46 & int32(2147483644)
					v55 = l0
					v57 = l2
					v60 = v48
					for {
						v67 = int32(0)
						v68 = base.Simd_g_v128_load(m, v55, v67)
						v69 = int32(16)
						v77 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v68, v69), base.Simd_g_const(&F_VP8LConvertBGRAToRGB565_C__k0)), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v68, int32(13)), base.Simd_g_const(&F_VP8LConvertBGRAToRGB565_C__k1)))
						v78 = int32(3)
						v79 = base.Simd_g_i32x4_extract_lane_l3(v77)
						*(*uint8)(unsafe.Add(mBase, uint32(v57+int32(6)))) = uint8(v79)
						v83 = int32(2)
						v84 = base.Simd_g_i32x4_extract_lane_l2(v77)
						*(*uint8)(unsafe.Add(mBase, uint32(v57+int32(4)))) = uint8(v84)
						v88 = int32(1)
						v89 = base.Simd_g_i32x4_extract_lane_l1(v77)
						*(*uint8)(unsafe.Add(mBase, uint32(v57+v83))) = uint8(v89)
						v92 = base.Simd_g_i32x4_extract_lane_l0(v77)
						*(*uint8)(unsafe.Add(mBase, uint32(v57))) = uint8(v92)
						v96 = int32(5)
						v104 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v68, v96), base.Simd_g_const(&F_VP8LConvertBGRAToRGB565_C__k2)), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v68, v78), base.Simd_g_const(&F_VP8LConvertBGRAToRGB565_C__k3)))
						v106 = base.Simd_g_i32x4_extract_lane_l3(v104)
						*(*uint8)(unsafe.Add(mBase, uint32(v57+int32(7)))) = uint8(v106)
						v111 = base.Simd_g_i32x4_extract_lane_l2(v104)
						*(*uint8)(unsafe.Add(mBase, uint32(v57+v96))) = uint8(v111)
						v116 = base.Simd_g_i32x4_extract_lane_l1(v104)
						*(*uint8)(unsafe.Add(mBase, uint32(v57+v78))) = uint8(v116)
						v121 = base.Simd_g_i32x4_extract_lane_l0(v104)
						*(*uint8)(unsafe.Add(mBase, uint32(v57+v88))) = uint8(v121)
						v128 = v60 + int32(-4)
						if v128 != 0 {
							v55 = v55 + v69
							v57 = v57 + int32(8)
							v60 = v128
							continue
						} else {
							break
						}
						break
					}
					if v46 == v48 {
					} else {
						v132 = l2 + v48<<(uint(v45)%32)
						v135 = l0 + v48<<(uint(v43)%32)
						v142 = v132
						v145 = v135
						for {
							v151 = *(*int32)(unsafe.Add(mBase, uint32(v145)))
							v160 = int32(base.Ui32(v151)>>(uint(int32(5))%32))&int32(224) | int32(base.Ui32(v151)>>(uint(int32(3))%32))&int32(31)
							*(*uint8)(unsafe.Add(mBase, uint32(v142)+1)) = uint8(v160)
							v170 = int32(base.Ui32(v151)>>(uint(int32(16))%32))&int32(248) | int32(base.Ui32(v151)>>(uint(int32(13))%32))&int32(7)
							*(*uint8)(unsafe.Add(mBase, uint32(v142))) = uint8(v170)
							v175 = v145 + int32(4)
							if base.Ui32(v175) < base.Ui32(v15) {
								v142 = v142 + int32(2)
								v145 = v175
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

var F_VP8LConvertBGRAToRGB565_C__k0 = [2]uint64{0xf8000000f8, 0xf8000000f8}
var F_VP8LConvertBGRAToRGB565_C__k1 = [2]uint64{0x700000007, 0x700000007}
var F_VP8LConvertBGRAToRGB565_C__k2 = [2]uint64{0xe0000000e0, 0xe0000000e0}
var F_VP8LConvertBGRAToRGB565_C__k3 = [2]uint64{0x1f0000001f, 0x1f0000001f}

func F_VP8LConvertBGRAToRGBA4444_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v16 int32
	_ = v16
	var v18 int32
	_ = v18
	var v20 int32
	_ = v20
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v44 int32
	_ = v44
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v49 int32
	_ = v49
	var v56 int32
	_ = v56
	var v58 int32
	_ = v58
	var v61 int32
	_ = v61
	var v69 int32
	_ = v69
	var v70 base.V128
	_ = v70
	var v71 base.V128
	_ = v71
	var v75 base.V128
	_ = v75
	var v76 int32
	_ = v76
	var v77 int32
	_ = v77
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v92 int32
	_ = v92
	var v96 int32
	_ = v96
	var v103 base.V128
	_ = v103
	var v105 int32
	_ = v105
	var v110 int32
	_ = v110
	var v115 int32
	_ = v115
	var v118 int32
	_ = v118
	var v125 int32
	_ = v125
	var v129 int32
	_ = v129
	var v132 int32
	_ = v132
	var v140 int32
	_ = v140
	var v143 int32
	_ = v143
	var v150 int32
	_ = v150
	var v151 int32
	_ = v151
	var v155 int32
	_ = v155
	var v165 int32
	_ = v165
	var v170 int32
	_ = v170
	if l1 < int32(1) {
	} else {
		v16 = l0 + l1<<(uint(int32(2))%32)
		v18 = l0 + int32(4)
		if base.Ui32(v18) < base.Ui32(v16) {
			v20 = v16
		} else {
			v20 = v18
		}
		v22 = l0 ^ int32(-1)
		v23 = v20 + v22
		if base.Ui32(v23) < base.Ui32(int32(44)) {
			v129 = l2
			v132 = l0
			v140 = v129
			v143 = v132
			for {
				v150 = *(*int32)(unsafe.Add(mBase, uint32(v143)))
				v151 = int32(240)
				v155 = v150&v151 | int32(base.Ui32(v150)>>(uint(int32(28))%32))
				*(*uint8)(unsafe.Add(mBase, uint32(v140)+1)) = uint8(v155)
				v165 = int32(base.Ui32(v150)>>(uint(int32(16))%32))&v151 | int32(base.Ui32(v150)>>(uint(int32(12))%32))&int32(15)
				*(*uint8)(unsafe.Add(mBase, uint32(v140))) = uint8(v165)
				v170 = v143 + int32(4)
				if base.Ui32(v170) < base.Ui32(v16) {
					v140 = v140 + int32(2)
					v143 = v170
					continue
				} else {
					break
				}
				break
			}
		} else {
			v28 = l1<<(uint(int32(2))%32) + l0
			if base.Ui32(v18) < base.Ui32(v28) {
				v30 = v28
			} else {
				v30 = v18
			}
			v31 = v30 + v22
			if base.Ui32(v18+v31&int32(-4)) <= base.Ui32(l2) {
				v44 = int32(2)
				v46 = int32(1)
				v47 = int32(base.Ui32(v23)>>(uint(v44)%32)) + v46
				v49 = v47 & int32(2147483644)
				v56 = l0
				v58 = l2
				v61 = v49
				for {
					v69 = int32(0)
					v70 = base.Simd_g_v128_load(m, v56, v69)
					v71 = base.Simd_g_const(&F_VP8LConvertBGRAToRGBA4444_C__k0)
					v75 = base.Simd_g_v128_or(base.Simd_g_v128_and(v70, v71), base.Simd_g_i32x4_shr_u(v70, int32(28)))
					v76 = int32(3)
					v77 = base.Simd_g_i32x4_extract_lane_l3(v75)
					*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(7)))) = uint8(v77)
					v81 = int32(2)
					v82 = base.Simd_g_i32x4_extract_lane_l2(v75)
					*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(5)))) = uint8(v82)
					v86 = int32(1)
					v87 = base.Simd_g_i32x4_extract_lane_l1(v75)
					*(*uint8)(unsafe.Add(mBase, uint32(v58+v76))) = uint8(v87)
					v92 = base.Simd_g_i32x4_extract_lane_l0(v75)
					*(*uint8)(unsafe.Add(mBase, uint32(v58+v86))) = uint8(v92)
					v96 = int32(16)
					v103 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v70, v96), v71), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v70, int32(12)), base.Simd_g_const(&F_VP8LConvertBGRAToRGBA4444_C__k1)))
					v105 = base.Simd_g_i32x4_extract_lane_l3(v103)
					*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(6)))) = uint8(v105)
					v110 = base.Simd_g_i32x4_extract_lane_l2(v103)
					*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(4)))) = uint8(v110)
					v115 = base.Simd_g_i32x4_extract_lane_l1(v103)
					*(*uint8)(unsafe.Add(mBase, uint32(v58+v81))) = uint8(v115)
					v118 = base.Simd_g_i32x4_extract_lane_l0(v103)
					*(*uint8)(unsafe.Add(mBase, uint32(v58))) = uint8(v118)
					v125 = v61 + int32(-4)
					if v125 != 0 {
						v56 = v56 + v96
						v58 = v58 + int32(8)
						v61 = v125
						continue
					} else {
						break
					}
					break
				}
				if v47 == v49 {
				} else {
					v129 = l2 + v49<<(uint(v46)%32)
					v132 = l0 + v49<<(uint(v44)%32)
					v140 = v129
					v143 = v132
					for {
						v150 = *(*int32)(unsafe.Add(mBase, uint32(v143)))
						v151 = int32(240)
						v155 = v150&v151 | int32(base.Ui32(v150)>>(uint(int32(28))%32))
						*(*uint8)(unsafe.Add(mBase, uint32(v140)+1)) = uint8(v155)
						v165 = int32(base.Ui32(v150)>>(uint(int32(16))%32))&v151 | int32(base.Ui32(v150)>>(uint(int32(12))%32))&int32(15)
						*(*uint8)(unsafe.Add(mBase, uint32(v140))) = uint8(v165)
						v170 = v143 + int32(4)
						if base.Ui32(v170) < base.Ui32(v16) {
							v140 = v140 + int32(2)
							v143 = v170
							continue
						} else {
							break
						}
						break
					}
				}
			} else {
				if base.Ui32(l0) < base.Ui32(l2+int32(base.Ui32(v31)>>(uint(int32(1))%32))&int32(2147483646)+int32(2)) {
					v129 = l2
					v132 = l0
					v140 = v129
					v143 = v132
					for {
						v150 = *(*int32)(unsafe.Add(mBase, uint32(v143)))
						v151 = int32(240)
						v155 = v150&v151 | int32(base.Ui32(v150)>>(uint(int32(28))%32))
						*(*uint8)(unsafe.Add(mBase, uint32(v140)+1)) = uint8(v155)
						v165 = int32(base.Ui32(v150)>>(uint(int32(16))%32))&v151 | int32(base.Ui32(v150)>>(uint(int32(12))%32))&int32(15)
						*(*uint8)(unsafe.Add(mBase, uint32(v140))) = uint8(v165)
						v170 = v143 + int32(4)
						if base.Ui32(v170) < base.Ui32(v16) {
							v140 = v140 + int32(2)
							v143 = v170
							continue
						} else {
							break
						}
						break
					}
				} else {
					v44 = int32(2)
					v46 = int32(1)
					v47 = int32(base.Ui32(v23)>>(uint(v44)%32)) + v46
					v49 = v47 & int32(2147483644)
					v56 = l0
					v58 = l2
					v61 = v49
					for {
						v69 = int32(0)
						v70 = base.Simd_g_v128_load(m, v56, v69)
						v71 = base.Simd_g_const(&F_VP8LConvertBGRAToRGBA4444_C__k0)
						v75 = base.Simd_g_v128_or(base.Simd_g_v128_and(v70, v71), base.Simd_g_i32x4_shr_u(v70, int32(28)))
						v76 = int32(3)
						v77 = base.Simd_g_i32x4_extract_lane_l3(v75)
						*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(7)))) = uint8(v77)
						v81 = int32(2)
						v82 = base.Simd_g_i32x4_extract_lane_l2(v75)
						*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(5)))) = uint8(v82)
						v86 = int32(1)
						v87 = base.Simd_g_i32x4_extract_lane_l1(v75)
						*(*uint8)(unsafe.Add(mBase, uint32(v58+v76))) = uint8(v87)
						v92 = base.Simd_g_i32x4_extract_lane_l0(v75)
						*(*uint8)(unsafe.Add(mBase, uint32(v58+v86))) = uint8(v92)
						v96 = int32(16)
						v103 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v70, v96), v71), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v70, int32(12)), base.Simd_g_const(&F_VP8LConvertBGRAToRGBA4444_C__k1)))
						v105 = base.Simd_g_i32x4_extract_lane_l3(v103)
						*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(6)))) = uint8(v105)
						v110 = base.Simd_g_i32x4_extract_lane_l2(v103)
						*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(4)))) = uint8(v110)
						v115 = base.Simd_g_i32x4_extract_lane_l1(v103)
						*(*uint8)(unsafe.Add(mBase, uint32(v58+v81))) = uint8(v115)
						v118 = base.Simd_g_i32x4_extract_lane_l0(v103)
						*(*uint8)(unsafe.Add(mBase, uint32(v58))) = uint8(v118)
						v125 = v61 + int32(-4)
						if v125 != 0 {
							v56 = v56 + v96
							v58 = v58 + int32(8)
							v61 = v125
							continue
						} else {
							break
						}
						break
					}
					if v47 == v49 {
					} else {
						v129 = l2 + v49<<(uint(v46)%32)
						v132 = l0 + v49<<(uint(v44)%32)
						v140 = v129
						v143 = v132
						for {
							v150 = *(*int32)(unsafe.Add(mBase, uint32(v143)))
							v151 = int32(240)
							v155 = v150&v151 | int32(base.Ui32(v150)>>(uint(int32(28))%32))
							*(*uint8)(unsafe.Add(mBase, uint32(v140)+1)) = uint8(v155)
							v165 = int32(base.Ui32(v150)>>(uint(int32(16))%32))&v151 | int32(base.Ui32(v150)>>(uint(int32(12))%32))&int32(15)
							*(*uint8)(unsafe.Add(mBase, uint32(v140))) = uint8(v165)
							v170 = v143 + int32(4)
							if base.Ui32(v170) < base.Ui32(v16) {
								v140 = v140 + int32(2)
								v143 = v170
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

var F_VP8LConvertBGRAToRGBA4444_C__k0 = [2]uint64{0xf0000000f0, 0xf0000000f0}
var F_VP8LConvertBGRAToRGBA4444_C__k1 = [2]uint64{0xf0000000f, 0xf0000000f}

func F_VP8LSubtractGreenFromBlueAndRed_C(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v22 int32
	_ = v22
	var v23 base.V128
	_ = v23
	var v25 base.V128
	_ = v25
	var v32 int32
	_ = v32
	var v37 base.V128
	_ = v37
	var v43 int32
	_ = v43
	var v47 int32
	_ = v47
	var v57 int32
	_ = v57
	var v60 int32
	_ = v60
	var v63 int32
	_ = v63
	var v65 int32
	_ = v65
	var v82 int32
	_ = v82
	if l1 < int32(1) {
	} else {
		if base.Ui32(l1) < base.Ui32(int32(4)) {
			v47 = int32(0)
			v57 = l1 - v47
			v60 = l0 + v47<<(uint(int32(2))%32)
			for {
				v63 = *(*int32)(unsafe.Add(mBase, uint32(v60)))
				v65 = int32(base.Ui32(v63) >> (uint(int32(8)) % 32))
				*(*int32)(unsafe.Add(mBase, uint32(v60))) = (v63-v65)&int32(255) | v63&int32(-16711936) | (v63-v65<<(uint(int32(16))%32))&int32(16711680)
				v82 = v57 + int32(-1)
				if v82 != 0 {
					v57 = v82
					v60 = v60 + int32(4)
					continue
				} else {
					break
				}
				break
			}
		} else {
			v14 = l1 & int32(2147483644)
			v18 = v14
			v19 = l0
			for {
				v22 = int32(0)
				v23 = base.Simd_g_v128_load(m, v19, v22)
				v25 = base.Simd_g_i32x4_shr_u(v23, int32(8))
				v32 = int32(16)
				v37 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_sub(v23, v25), base.Simd_g_const(&F_VP8LSubtractGreenFromBlueAndRed_C__k0)), base.Simd_g_v128_and(v23, base.Simd_g_const(&F_VP8LSubtractGreenFromBlueAndRed_C__k1))), base.Simd_g_v128_and(base.Simd_g_i32x4_sub(v23, base.Simd_g_i32x4_shl(v25, v32)), base.Simd_g_const(&F_VP8LSubtractGreenFromBlueAndRed_C__k2)))
				base.Simd_g_v128_store(m, v19, v22, v37)
				v43 = v18 + int32(-4)
				if v43 != 0 {
					v18 = v43
					v19 = v19 + v32
					continue
				} else {
					break
				}
				break
			}
			if v14 == l1 {
			} else {
				v47 = v14
				v57 = l1 - v47
				v60 = l0 + v47<<(uint(int32(2))%32)
				for {
					v63 = *(*int32)(unsafe.Add(mBase, uint32(v60)))
					v65 = int32(base.Ui32(v63) >> (uint(int32(8)) % 32))
					*(*int32)(unsafe.Add(mBase, uint32(v60))) = (v63-v65)&int32(255) | v63&int32(-16711936) | (v63-v65<<(uint(int32(16))%32))&int32(16711680)
					v82 = v57 + int32(-1)
					if v82 != 0 {
						v57 = v82
						v60 = v60 + int32(4)
						continue
					} else {
						break
					}
					break
				}
			}
		}
	}
	return
}

var F_VP8LSubtractGreenFromBlueAndRed_C__k0 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_VP8LSubtractGreenFromBlueAndRed_C__k1 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_VP8LSubtractGreenFromBlueAndRed_C__k2 = [2]uint64{0xff000000ff0000, 0xff000000ff0000}

func F_VP8LTransformColorInverse_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v44 int32
	_ = v44
	var v45 int32
	_ = v45
	var v49 int32
	_ = v49
	var v50 base.V128
	_ = v50
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v54 base.V128
	_ = v54
	var v56 int32
	_ = v56
	var v60 base.V128
	_ = v60
	var v82 base.V128
	_ = v82
	var v90 int32
	_ = v90
	var v99 int32
	_ = v99
	var v110 int32
	_ = v110
	var v115 int32
	_ = v115
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v133 int32
	_ = v133
	var v135 int32
	_ = v135
	var v139 int32
	_ = v139
	var v160 int32
	_ = v160
	var v165 int32
	_ = v165
	if l2 < int32(1) {
	} else {
		v19 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+2)))
		v20 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+1)))
		v21 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0))))
		v22 = int32(0)
		if base.Ui32(l2) < base.Ui32(int32(4)) {
			v99 = v22
			v110 = v99 << (uint(int32(2)) % 32)
			v115 = l2 - v99
			v124 = l3 + v110
			v125 = l1 + v110
			for {
				v129 = *(*int32)(unsafe.Add(mBase, uint32(v125)))
				v130 = int32(16)
				v133 = v129 << (uint(v130) % 32) >> (uint(int32(24)) % 32)
				v135 = int32(5)
				v139 = v133*v21>>(uint(v135)%32) + int32(base.Ui32(v129)>>(uint(v130)%32))
				*(*int32)(unsafe.Add(mBase, uint32(v124))) = v139<<(uint(v130)%32)&int32(16711680) | v129&int32(-16711936) | (int32(base.Ui32(v133*v20)>>(uint(v135)%32))+v129+int32(base.Ui32(base.I32_extend8_s(v139)*v19)>>(uint(v135)%32)))&int32(255)
				v160 = int32(4)
				v165 = v115 + int32(-1)
				if v165 != 0 {
					v115 = v165
					v124 = v124 + v160
					v125 = v125 + v160
					continue
				} else {
					break
				}
				break
			}
		} else {
			if base.Ui32(l3-l1) < base.Ui32(int32(16)) {
				v99 = v22
				v110 = v99 << (uint(int32(2)) % 32)
				v115 = l2 - v99
				v124 = l3 + v110
				v125 = l1 + v110
				for {
					v129 = *(*int32)(unsafe.Add(mBase, uint32(v125)))
					v130 = int32(16)
					v133 = v129 << (uint(v130) % 32) >> (uint(int32(24)) % 32)
					v135 = int32(5)
					v139 = v133*v21>>(uint(v135)%32) + int32(base.Ui32(v129)>>(uint(v130)%32))
					*(*int32)(unsafe.Add(mBase, uint32(v124))) = v139<<(uint(v130)%32)&int32(16711680) | v129&int32(-16711936) | (int32(base.Ui32(v133*v20)>>(uint(v135)%32))+v129+int32(base.Ui32(base.I32_extend8_s(v139)*v19)>>(uint(v135)%32)))&int32(255)
					v160 = int32(4)
					v165 = v115 + int32(-1)
					if v165 != 0 {
						v115 = v165
						v124 = v124 + v160
						v125 = v125 + v160
						continue
					} else {
						break
					}
					break
				}
			} else {
				v32 = l2 & int32(2147483644)
				v33 = l3
				v44 = v32
				v45 = l1
				for {
					v49 = int32(0)
					v50 = base.Simd_g_v128_load(m, v45, v49)
					v51 = int32(16)
					v53 = int32(24)
					v54 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_shl(v50, v51), v53)
					v56 = int32(5)
					v60 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_mul(v54, base.Simd_g_i32x4_splat(v21)), v56), base.Simd_g_i32x4_shr_u(v50, v51))
					v82 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(v60, v51), base.Simd_g_const(&F_VP8LTransformColorInverse_C__k0)), base.Simd_g_v128_and(v50, base.Simd_g_const(&F_VP8LTransformColorInverse_C__k1))), base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v54, base.Simd_g_i32x4_splat(v20)), v56), v50), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_shl(v60, v53), v53), base.Simd_g_i32x4_splat(v19)), v56)), base.Simd_g_const(&F_VP8LTransformColorInverse_C__k2)))
					base.Simd_g_v128_store(m, v33, v49, v82)
					v90 = v44 + int32(-4)
					if v90 != 0 {
						v33 = v33 + v51
						v44 = v90
						v45 = v45 + v51
						continue
					} else {
						break
					}
					break
				}
				if v32 == l2 {
				} else {
					v99 = v32
					v110 = v99 << (uint(int32(2)) % 32)
					v115 = l2 - v99
					v124 = l3 + v110
					v125 = l1 + v110
					for {
						v129 = *(*int32)(unsafe.Add(mBase, uint32(v125)))
						v130 = int32(16)
						v133 = v129 << (uint(v130) % 32) >> (uint(int32(24)) % 32)
						v135 = int32(5)
						v139 = v133*v21>>(uint(v135)%32) + int32(base.Ui32(v129)>>(uint(v130)%32))
						*(*int32)(unsafe.Add(mBase, uint32(v124))) = v139<<(uint(v130)%32)&int32(16711680) | v129&int32(-16711936) | (int32(base.Ui32(v133*v20)>>(uint(v135)%32))+v129+int32(base.Ui32(base.I32_extend8_s(v139)*v19)>>(uint(v135)%32)))&int32(255)
						v160 = int32(4)
						v165 = v115 + int32(-1)
						if v165 != 0 {
							v115 = v165
							v124 = v124 + v160
							v125 = v125 + v160
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

var F_VP8LTransformColorInverse_C__k0 = [2]uint64{0xff000000ff0000, 0xff000000ff0000}
var F_VP8LTransformColorInverse_C__k1 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_VP8LTransformColorInverse_C__k2 = [2]uint64{0xff000000ff, 0xff000000ff}

func F_VP8LTransformColor_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v38 int32
	_ = v38
	var v42 int32
	_ = v42
	var v43 base.V128
	_ = v43
	var v44 int32
	_ = v44
	var v45 base.V128
	_ = v45
	var v48 int32
	_ = v48
	var v49 base.V128
	_ = v49
	var v51 int32
	_ = v51
	var v75 base.V128
	_ = v75
	var v81 int32
	_ = v81
	var v89 int32
	_ = v89
	var v101 int32
	_ = v101
	var v103 int32
	_ = v103
	var v115 int32
	_ = v115
	var v116 int32
	_ = v116
	var v117 int32
	_ = v117
	var v121 int32
	_ = v121
	var v123 int32
	_ = v123
	var v149 int32
	_ = v149
	if l2 < int32(1) {
	} else {
		v17 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+2)))
		v18 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+1)))
		v19 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0))))
		if base.Ui32(l2) < base.Ui32(int32(4)) {
			v89 = int32(0)
			v101 = l1 + v89<<(uint(int32(2))%32)
			v103 = l2 - v89
			for {
				v115 = *(*int32)(unsafe.Add(mBase, uint32(v101)))
				v116 = int32(16)
				v117 = int32(base.Ui32(v115) >> (uint(v116) % 32))
				v121 = v115 << (uint(v116) % 32) >> (uint(int32(24)) % 32)
				v123 = int32(5)
				*(*int32)(unsafe.Add(mBase, uint32(v101))) = (v117-int32(base.Ui32(v121*v19)>>(uint(v123)%32)))<<(uint(v116)%32)&int32(16711680) | v115&int32(-16711936) | (v115-(int32(base.Ui32(v121*v18)>>(uint(v123)%32))+int32(base.Ui32(base.I32_extend8_s(v117)*v17)>>(uint(v123)%32))))&int32(255)
				v149 = v103 + int32(-1)
				if v149 != 0 {
					v101 = v101 + int32(4)
					v103 = v149
					continue
				} else {
					break
				}
				break
			}
		} else {
			v27 = l2 & int32(2147483644)
			v28 = l1
			v38 = v27
			for {
				v42 = int32(0)
				v43 = base.Simd_g_v128_load(m, v28, v42)
				v44 = int32(16)
				v45 = base.Simd_g_i32x4_shr_u(v43, v44)
				v48 = int32(24)
				v49 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_shl(v43, v44), v48)
				v51 = int32(5)
				v75 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(base.Simd_g_i32x4_sub(v45, base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v49, base.Simd_g_i32x4_splat(v19)), v51)), v44), base.Simd_g_const(&F_VP8LTransformColor_C__k0)), base.Simd_g_v128_and(v43, base.Simd_g_const(&F_VP8LTransformColor_C__k1))), base.Simd_g_v128_and(base.Simd_g_i32x4_sub(v43, base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v49, base.Simd_g_i32x4_splat(v18)), v51), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_shl(v45, v48), v48), base.Simd_g_i32x4_splat(v17)), v51))), base.Simd_g_const(&F_VP8LTransformColor_C__k2)))
				base.Simd_g_v128_store(m, v28, v42, v75)
				v81 = v38 + int32(-4)
				if v81 != 0 {
					v28 = v28 + v44
					v38 = v81
					continue
				} else {
					break
				}
				break
			}
			if v27 == l2 {
			} else {
				v89 = v27
				v101 = l1 + v89<<(uint(int32(2))%32)
				v103 = l2 - v89
				for {
					v115 = *(*int32)(unsafe.Add(mBase, uint32(v101)))
					v116 = int32(16)
					v117 = int32(base.Ui32(v115) >> (uint(v116) % 32))
					v121 = v115 << (uint(v116) % 32) >> (uint(int32(24)) % 32)
					v123 = int32(5)
					*(*int32)(unsafe.Add(mBase, uint32(v101))) = (v117-int32(base.Ui32(v121*v19)>>(uint(v123)%32)))<<(uint(v116)%32)&int32(16711680) | v115&int32(-16711936) | (v115-(int32(base.Ui32(v121*v18)>>(uint(v123)%32))+int32(base.Ui32(base.I32_extend8_s(v117)*v17)>>(uint(v123)%32))))&int32(255)
					v149 = v103 + int32(-1)
					if v149 != 0 {
						v101 = v101 + int32(4)
						v103 = v149
						continue
					} else {
						break
					}
					break
				}
			}
		}
	}
	return
}

var F_VP8LTransformColor_C__k0 = [2]uint64{0xff000000ff0000, 0xff000000ff0000}
var F_VP8LTransformColor_C__k1 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_VP8LTransformColor_C__k2 = [2]uint64{0xff000000ff, 0xff000000ff}

func F_VP8YuvToArgb32_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v29 base.V128
	_ = v29
	var v31 int32
	_ = v31
	var v32 base.V128
	_ = v32
	var v33 base.V128
	_ = v33
	var v34 base.V128
	_ = v34
	var v36 base.V128
	_ = v36
	var v40 base.V128
	_ = v40
	var v41 base.V128
	_ = v41
	var v44 base.V128
	_ = v44
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v48 base.V128
	_ = v48
	var v50 base.V128
	_ = v50
	var v56 base.V128
	_ = v56
	var v58 base.V128
	_ = v58
	var v59 base.V128
	_ = v59
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v70 int32
	_ = v70
	var v72 base.V128
	_ = v72
	var v73 base.V128
	_ = v73
	var v83 base.V128
	_ = v83
	var v93 base.V128
	_ = v93
	var v95 base.V128
	_ = v95
	var v97 base.V128
	_ = v97
	var v99 base.V128
	_ = v99
	var v103 base.V128
	_ = v103
	v18 = l3
	v19 = int32(0)
	for {
		v29 = base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k0)
		v31 = int32(0)
		v32 = base.Simd_g_v128_load64_zero(m, l0+v19, v31)
		v33 = base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k1)
		v34 = base.Simd_g_i8x16_shuffle2(v29, v32, base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k3))
		v36 = base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k4)
		v40 = base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k5)
		v41 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v34), v36), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v34), v36), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k7))
		v44 = base.Simd_g_v128_load64_zero(m, l1+v19, v31)
		v46 = base.Simd_g_i8x16_shuffle2(v29, v44, base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k3))
		v47 = base.Simd_g_i32x4_extend_low_i16x8_u(v46)
		v48 = base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k8)
		v50 = base.Simd_g_i32x4_extend_high_i16x8_u(v46)
		v56 = base.Simd_g_v128_load64_zero(m, l2+v19, v31)
		v58 = base.Simd_g_i8x16_shuffle2(v29, v56, base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k3))
		v59 = base.Simd_g_i32x4_extend_low_i16x8_u(v58)
		v60 = base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k9)
		v62 = base.Simd_g_i32x4_extend_high_i16x8_u(v58)
		v70 = int32(6)
		v72 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k10), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v41, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v47, v48), base.Simd_g_i32x4_mul(v50, v48), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v59, v60), base.Simd_g_i32x4_mul(v62, v60), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k7)))), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k11)), v70))
		v73 = base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k12)
		v83 = base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k13)
		v93 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v41, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v59, v73), base.Simd_g_i32x4_mul(v62, v73), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k7))), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k14)), v70), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v47, v83), base.Simd_g_i32x4_mul(v50, v83), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k7)), v41), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k15)), v70))
		v95 = base.Simd_g_i8x16_shuffle2(v72, v93, base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k3))
		v97 = base.Simd_g_i8x16_shuffle2(v72, v93, base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k16), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k17))
		v99 = base.Simd_g_i8x16_shuffle2(v95, v97, base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k18), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k19))
		base.Simd_g_v128_store(m, v18, int32(16), v99)
		v103 = base.Simd_g_i8x16_shuffle2(v95, v97, base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k20), base.Simd_g_const(&F_VP8YuvToArgb32_SSE2__k21))
		base.Simd_g_v128_store(m, v18, v31, v103)
		if base.Ui32(v19) < base.Ui32(int32(24)) {
			v18 = v18 + int32(32)
			v19 = v19 + int32(8)
			continue
		} else {
			break
		}
		break
	}
	return
}

var F_VP8YuvToArgb32_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_VP8YuvToArgb32_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_VP8YuvToArgb32_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_VP8YuvToArgb32_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_VP8YuvToArgb32_SSE2__k4 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_VP8YuvToArgb32_SSE2__k5 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_VP8YuvToArgb32_SSE2__k6 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_VP8YuvToArgb32_SSE2__k7 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_VP8YuvToArgb32_SSE2__k8 = [2]uint64{0x191300001913, 0x191300001913}
var F_VP8YuvToArgb32_SSE2__k9 = [2]uint64{0x340800003408, 0x340800003408}
var F_VP8YuvToArgb32_SSE2__k10 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_VP8YuvToArgb32_SSE2__k11 = [2]uint64{0x2204220422042204, 0x2204220422042204}
var F_VP8YuvToArgb32_SSE2__k12 = [2]uint64{0x662500006625, 0x662500006625}
var F_VP8YuvToArgb32_SSE2__k13 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_VP8YuvToArgb32_SSE2__k14 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_VP8YuvToArgb32_SSE2__k15 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_VP8YuvToArgb32_SSE2__k16 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_VP8YuvToArgb32_SSE2__k17 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_VP8YuvToArgb32_SSE2__k18 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_VP8YuvToArgb32_SSE2__k19 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_VP8YuvToArgb32_SSE2__k20 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_VP8YuvToArgb32_SSE2__k21 = [2]uint64{0x302808001008080, 0x706808005048080}

func F_VP8YuvToBgr32_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	var v5 base.V128
	_ = v5
	var v37 int32
	_ = v37
	var v38 base.V128
	_ = v38
	var v39 base.V128
	_ = v39
	var v40 base.V128
	_ = v40
	var v42 base.V128
	_ = v42
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v49 base.V128
	_ = v49
	var v51 base.V128
	_ = v51
	var v52 base.V128
	_ = v52
	var v53 base.V128
	_ = v53
	var v55 base.V128
	_ = v55
	var v60 base.V128
	_ = v60
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v65 base.V128
	_ = v65
	var v67 base.V128
	_ = v67
	var v73 base.V128
	_ = v73
	var v75 base.V128
	_ = v75
	var v77 base.V128
	_ = v77
	var v78 base.V128
	_ = v78
	var v80 base.V128
	_ = v80
	var v88 base.V128
	_ = v88
	var v89 base.V128
	_ = v89
	var v91 int32
	_ = v91
	var v92 base.V128
	_ = v92
	var v94 base.V128
	_ = v94
	var v100 base.V128
	_ = v100
	var v102 base.V128
	_ = v102
	var v104 base.V128
	_ = v104
	var v105 base.V128
	_ = v105
	var v107 base.V128
	_ = v107
	var v115 int32
	_ = v115
	var v116 base.V128
	_ = v116
	var v118 base.V128
	_ = v118
	var v124 base.V128
	_ = v124
	var v126 base.V128
	_ = v126
	var v128 base.V128
	_ = v128
	var v129 base.V128
	_ = v129
	var v131 base.V128
	_ = v131
	var v139 base.V128
	_ = v139
	var v141 base.V128
	_ = v141
	var v145 base.V128
	_ = v145
	var v147 base.V128
	_ = v147
	var v148 base.V128
	_ = v148
	var v149 base.V128
	_ = v149
	var v151 base.V128
	_ = v151
	var v156 base.V128
	_ = v156
	var v161 base.V128
	_ = v161
	var v163 base.V128
	_ = v163
	var v164 base.V128
	_ = v164
	var v166 base.V128
	_ = v166
	var v174 base.V128
	_ = v174
	var v178 base.V128
	_ = v178
	var v180 base.V128
	_ = v180
	var v181 base.V128
	_ = v181
	var v183 base.V128
	_ = v183
	var v192 base.V128
	_ = v192
	var v194 base.V128
	_ = v194
	var v195 base.V128
	_ = v195
	var v197 base.V128
	_ = v197
	var v205 base.V128
	_ = v205
	var v208 base.V128
	_ = v208
	var v211 base.V128
	_ = v211
	var v213 base.V128
	_ = v213
	var v218 base.V128
	_ = v218
	var v225 base.V128
	_ = v225
	var v242 base.V128
	_ = v242
	var v271 base.V128
	_ = v271
	var v274 base.V128
	_ = v274
	var v281 base.V128
	_ = v281
	var v284 base.V128
	_ = v284
	var v286 base.V128
	_ = v286
	var v291 base.V128
	_ = v291
	var v295 base.V128
	_ = v295
	var v297 base.V128
	_ = v297
	var v302 base.V128
	_ = v302
	var v305 base.V128
	_ = v305
	var v308 base.V128
	_ = v308
	var v313 base.V128
	_ = v313
	var v320 base.V128
	_ = v320
	var v323 base.V128
	_ = v323
	var v330 base.V128
	_ = v330
	var v333 base.V128
	_ = v333
	var v336 base.V128
	_ = v336
	var v341 base.V128
	_ = v341
	var v346 base.V128
	_ = v346
	var v351 base.V128
	_ = v351
	var v354 base.V128
	_ = v354
	var v357 base.V128
	_ = v357
	var v362 base.V128
	_ = v362
	var v367 base.V128
	_ = v367
	var v370 base.V128
	_ = v370
	var v375 base.V128
	_ = v375
	var v380 base.V128
	_ = v380
	var v385 base.V128
	_ = v385
	v5 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k0)
	v37 = int32(0)
	v38 = base.Simd_g_v128_load64_zero(m, l0, v37)
	v39 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k1)
	v40 = base.Simd_g_i8x16_shuffle2(v5, v38, base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k3))
	v42 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k4)
	v46 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k5)
	v47 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v40), v42), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v40), v42), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7))
	v49 = base.Simd_g_v128_load64_zero(m, l2, v37)
	v51 = base.Simd_g_i8x16_shuffle2(v5, v49, base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k3))
	v52 = base.Simd_g_i32x4_extend_low_i16x8_u(v51)
	v53 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k8)
	v55 = base.Simd_g_i32x4_extend_high_i16x8_u(v51)
	v60 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k9)
	v62 = int32(6)
	v64 = int32(8)
	v65 = base.Simd_g_v128_load64_zero(m, l0, v64)
	v67 = base.Simd_g_i8x16_shuffle2(v5, v65, base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k3))
	v73 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v67), v42), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v67), v42), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7))
	v75 = base.Simd_g_v128_load64_zero(m, l2, v64)
	v77 = base.Simd_g_i8x16_shuffle2(v5, v75, base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k3))
	v78 = base.Simd_g_i32x4_extend_low_i16x8_u(v77)
	v80 = base.Simd_g_i32x4_extend_high_i16x8_u(v77)
	v88 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v47, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v52, v53), base.Simd_g_i32x4_mul(v55, v53), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7))), v60), v62), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v73, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v78, v53), base.Simd_g_i32x4_mul(v80, v53), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7))), v60), v62))
	v89 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k10)
	v91 = int32(16)
	v92 = base.Simd_g_v128_load64_zero(m, l0, v91)
	v94 = base.Simd_g_i8x16_shuffle2(v5, v92, base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k3))
	v100 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v94), v42), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v94), v42), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7))
	v102 = base.Simd_g_v128_load64_zero(m, l2, v91)
	v104 = base.Simd_g_i8x16_shuffle2(v5, v102, base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k3))
	v105 = base.Simd_g_i32x4_extend_low_i16x8_u(v104)
	v107 = base.Simd_g_i32x4_extend_high_i16x8_u(v104)
	v115 = int32(24)
	v116 = base.Simd_g_v128_load64_zero(m, l0, v115)
	v118 = base.Simd_g_i8x16_shuffle2(v5, v116, base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k3))
	v124 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v118), v42), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v118), v42), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7))
	v126 = base.Simd_g_v128_load64_zero(m, l2, v115)
	v128 = base.Simd_g_i8x16_shuffle2(v5, v126, base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k3))
	v129 = base.Simd_g_i32x4_extend_low_i16x8_u(v128)
	v131 = base.Simd_g_i32x4_extend_high_i16x8_u(v128)
	v139 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v100, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v105, v53), base.Simd_g_i32x4_mul(v107, v53), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7))), v60), v62), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v124, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v129, v53), base.Simd_g_i32x4_mul(v131, v53), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7))), v60), v62))
	v141 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v88, v89), base.Simd_g_v128_and(v139, v89))
	v145 = base.Simd_g_v128_load64_zero(m, l1, v37)
	v147 = base.Simd_g_i8x16_shuffle2(v5, v145, base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k3))
	v148 = base.Simd_g_i32x4_extend_low_i16x8_u(v147)
	v149 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k11)
	v151 = base.Simd_g_i32x4_extend_high_i16x8_u(v147)
	v156 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k12)
	v161 = base.Simd_g_v128_load64_zero(m, l1, v64)
	v163 = base.Simd_g_i8x16_shuffle2(v5, v161, base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k3))
	v164 = base.Simd_g_i32x4_extend_low_i16x8_u(v163)
	v166 = base.Simd_g_i32x4_extend_high_i16x8_u(v163)
	v174 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v148, v149), base.Simd_g_i32x4_mul(v151, v149), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7)), v47), v156), v62), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v164, v149), base.Simd_g_i32x4_mul(v166, v149), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7)), v73), v156), v62))
	v178 = base.Simd_g_v128_load64_zero(m, l1, v91)
	v180 = base.Simd_g_i8x16_shuffle2(v5, v178, base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k3))
	v181 = base.Simd_g_i32x4_extend_low_i16x8_u(v180)
	v183 = base.Simd_g_i32x4_extend_high_i16x8_u(v180)
	v192 = base.Simd_g_v128_load64_zero(m, l1, v115)
	v194 = base.Simd_g_i8x16_shuffle2(v5, v192, base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k3))
	v195 = base.Simd_g_i32x4_extend_low_i16x8_u(v194)
	v197 = base.Simd_g_i32x4_extend_high_i16x8_u(v194)
	v205 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v181, v149), base.Simd_g_i32x4_mul(v183, v149), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7)), v100), v156), v62), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v195, v149), base.Simd_g_i32x4_mul(v197, v149), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7)), v124), v156), v62))
	v208 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v174, v64), base.Simd_g_i16x8_shr_u(v205, v64))
	v211 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v141, v64), base.Simd_g_i16x8_shr_u(v208, v64))
	v213 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k13)
	v218 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k14)
	v225 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k15)
	v242 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v47, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v148, v213), base.Simd_g_i32x4_mul(v151, v213), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v52, v218), base.Simd_g_i32x4_mul(v55, v218), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7)))), v225), v62), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v73, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v164, v213), base.Simd_g_i32x4_mul(v166, v213), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v78, v218), base.Simd_g_i32x4_mul(v80, v218), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7)))), v225), v62))
	v271 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v100, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v181, v213), base.Simd_g_i32x4_mul(v183, v213), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v105, v218), base.Simd_g_i32x4_mul(v107, v218), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7)))), v225), v62), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v124, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v195, v213), base.Simd_g_i32x4_mul(v197, v213), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v129, v218), base.Simd_g_i32x4_mul(v131, v218), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE2__k7)))), v225), v62))
	v274 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v242, v64), base.Simd_g_i16x8_shr_u(v271, v64))
	v281 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v88, v64), base.Simd_g_i16x8_shr_u(v139, v64))
	v284 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v274, v64), base.Simd_g_i16x8_shr_u(v281, v64))
	v286 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v211, v89), base.Simd_g_v128_and(v284, v89))
	v291 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v174, v89), base.Simd_g_v128_and(v205, v89))
	v295 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v242, v89), base.Simd_g_v128_and(v271, v89))
	v297 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v291, v89), base.Simd_g_v128_and(v295, v89))
	v302 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v141, v89), base.Simd_g_v128_and(v208, v89))
	v305 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v297, v64), base.Simd_g_i16x8_shr_u(v302, v64))
	v308 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v286, v64), base.Simd_g_i16x8_shr_u(v305, v64))
	v313 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v274, v89), base.Simd_g_v128_and(v281, v89))
	v320 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v291, v64), base.Simd_g_i16x8_shr_u(v295, v64))
	v323 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v313, v64), base.Simd_g_i16x8_shr_u(v320, v64))
	v330 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v211, v64), base.Simd_g_i16x8_shr_u(v284, v64))
	v333 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v323, v64), base.Simd_g_i16x8_shr_u(v330, v64))
	v336 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v308, v64), base.Simd_g_i16x8_shr_u(v333, v64))
	base.Simd_g_v128_store(m, l3, int32(80), v336)
	v341 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v323, v89), base.Simd_g_v128_and(v330, v89))
	v346 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v297, v89), base.Simd_g_v128_and(v302, v89))
	v351 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v313, v89), base.Simd_g_v128_and(v320, v89))
	v354 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v346, v64), base.Simd_g_i16x8_shr_u(v351, v64))
	v357 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v341, v64), base.Simd_g_i16x8_shr_u(v354, v64))
	base.Simd_g_v128_store(m, l3, int32(64), v357)
	v362 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v346, v89), base.Simd_g_v128_and(v351, v89))
	v367 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v286, v89), base.Simd_g_v128_and(v305, v89))
	v370 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v362, v64), base.Simd_g_i16x8_shr_u(v367, v64))
	base.Simd_g_v128_store(m, l3, int32(48), v370)
	v375 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v308, v89), base.Simd_g_v128_and(v333, v89))
	base.Simd_g_v128_store(m, l3, int32(32), v375)
	v380 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v341, v89), base.Simd_g_v128_and(v354, v89))
	base.Simd_g_v128_store(m, l3, v91, v380)
	v385 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v362, v89), base.Simd_g_v128_and(v367, v89))
	base.Simd_g_v128_store(m, l3, v37, v385)
	return
}

var F_VP8YuvToBgr32_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_VP8YuvToBgr32_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_VP8YuvToBgr32_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_VP8YuvToBgr32_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_VP8YuvToBgr32_SSE2__k4 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_VP8YuvToBgr32_SSE2__k5 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_VP8YuvToBgr32_SSE2__k6 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_VP8YuvToBgr32_SSE2__k7 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_VP8YuvToBgr32_SSE2__k8 = [2]uint64{0x662500006625, 0x662500006625}
var F_VP8YuvToBgr32_SSE2__k9 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_VP8YuvToBgr32_SSE2__k10 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_VP8YuvToBgr32_SSE2__k11 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_VP8YuvToBgr32_SSE2__k12 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_VP8YuvToBgr32_SSE2__k13 = [2]uint64{0x191300001913, 0x191300001913}
var F_VP8YuvToBgr32_SSE2__k14 = [2]uint64{0x340800003408, 0x340800003408}
var F_VP8YuvToBgr32_SSE2__k15 = [2]uint64{0x2204220422042204, 0x2204220422042204}

func F_VP8YuvToBgr32_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	var v5 base.V128
	_ = v5
	var v33 int32
	_ = v33
	var v34 base.V128
	_ = v34
	var v35 base.V128
	_ = v35
	var v36 base.V128
	_ = v36
	var v38 base.V128
	_ = v38
	var v42 base.V128
	_ = v42
	var v43 base.V128
	_ = v43
	var v45 base.V128
	_ = v45
	var v47 base.V128
	_ = v47
	var v48 base.V128
	_ = v48
	var v49 base.V128
	_ = v49
	var v51 base.V128
	_ = v51
	var v56 base.V128
	_ = v56
	var v58 base.V128
	_ = v58
	var v59 base.V128
	_ = v59
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v68 base.V128
	_ = v68
	var v70 int32
	_ = v70
	var v72 int32
	_ = v72
	var v73 base.V128
	_ = v73
	var v75 base.V128
	_ = v75
	var v81 base.V128
	_ = v81
	var v83 base.V128
	_ = v83
	var v85 base.V128
	_ = v85
	var v86 base.V128
	_ = v86
	var v88 base.V128
	_ = v88
	var v93 base.V128
	_ = v93
	var v95 base.V128
	_ = v95
	var v96 base.V128
	_ = v96
	var v98 base.V128
	_ = v98
	var v107 base.V128
	_ = v107
	var v108 base.V128
	_ = v108
	var v110 base.V128
	_ = v110
	var v116 base.V128
	_ = v116
	var v128 base.V128
	_ = v128
	var v129 base.V128
	_ = v129
	var v132 base.V128
	_ = v132
	var v138 base.V128
	_ = v138
	var v150 base.V128
	_ = v150
	var v151 base.V128
	_ = v151
	var v153 base.V128
	_ = v153
	var v156 base.V128
	_ = v156
	var v158 base.V128
	_ = v158
	var v161 base.V128
	_ = v161
	var v163 base.V128
	_ = v163
	var v166 base.V128
	_ = v166
	var v168 base.V128
	_ = v168
	var v171 base.V128
	_ = v171
	var v173 base.V128
	_ = v173
	var v176 int32
	_ = v176
	var v177 base.V128
	_ = v177
	var v179 base.V128
	_ = v179
	var v185 base.V128
	_ = v185
	var v187 base.V128
	_ = v187
	var v189 base.V128
	_ = v189
	var v190 base.V128
	_ = v190
	var v192 base.V128
	_ = v192
	var v197 base.V128
	_ = v197
	var v199 base.V128
	_ = v199
	var v200 base.V128
	_ = v200
	var v202 base.V128
	_ = v202
	var v211 int32
	_ = v211
	var v212 base.V128
	_ = v212
	var v214 base.V128
	_ = v214
	var v220 base.V128
	_ = v220
	var v222 base.V128
	_ = v222
	var v224 base.V128
	_ = v224
	var v225 base.V128
	_ = v225
	var v227 base.V128
	_ = v227
	var v232 base.V128
	_ = v232
	var v234 base.V128
	_ = v234
	var v235 base.V128
	_ = v235
	var v237 base.V128
	_ = v237
	var v246 base.V128
	_ = v246
	var v264 base.V128
	_ = v264
	var v283 base.V128
	_ = v283
	var v285 base.V128
	_ = v285
	var v292 base.V128
	_ = v292
	var v299 base.V128
	_ = v299
	v5 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k0)
	v33 = int32(16)
	v34 = base.Simd_g_v128_load64_zero(m, l0, v33)
	v35 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k1)
	v36 = base.Simd_g_i8x16_shuffle2(v5, v34, base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k3))
	v38 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k4)
	v42 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k5)
	v43 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v36), v38), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v36), v38), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7))
	v45 = base.Simd_g_v128_load64_zero(m, l1, v33)
	v47 = base.Simd_g_i8x16_shuffle2(v5, v45, base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k3))
	v48 = base.Simd_g_i32x4_extend_low_i16x8_u(v47)
	v49 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k8)
	v51 = base.Simd_g_i32x4_extend_high_i16x8_u(v47)
	v56 = base.Simd_g_v128_load64_zero(m, l2, v33)
	v58 = base.Simd_g_i8x16_shuffle2(v5, v56, base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k3))
	v59 = base.Simd_g_i32x4_extend_low_i16x8_u(v58)
	v60 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k9)
	v62 = base.Simd_g_i32x4_extend_high_i16x8_u(v58)
	v68 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k10)
	v70 = int32(6)
	v72 = int32(24)
	v73 = base.Simd_g_v128_load64_zero(m, l0, v72)
	v75 = base.Simd_g_i8x16_shuffle2(v5, v73, base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k3))
	v81 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v75), v38), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v75), v38), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7))
	v83 = base.Simd_g_v128_load64_zero(m, l1, v72)
	v85 = base.Simd_g_i8x16_shuffle2(v5, v83, base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k3))
	v86 = base.Simd_g_i32x4_extend_low_i16x8_u(v85)
	v88 = base.Simd_g_i32x4_extend_high_i16x8_u(v85)
	v93 = base.Simd_g_v128_load64_zero(m, l2, v72)
	v95 = base.Simd_g_i8x16_shuffle2(v5, v93, base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k3))
	v96 = base.Simd_g_i32x4_extend_low_i16x8_u(v95)
	v98 = base.Simd_g_i32x4_extend_high_i16x8_u(v95)
	v107 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v43, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v48, v49), base.Simd_g_i32x4_mul(v51, v49), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v59, v60), base.Simd_g_i32x4_mul(v62, v60), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7)))), v68), v70), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v81, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v86, v49), base.Simd_g_i32x4_mul(v88, v49), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v96, v60), base.Simd_g_i32x4_mul(v98, v60), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7)))), v68), v70))
	v108 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k11)
	v110 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k12)
	v116 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k13)
	v128 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v48, v110), base.Simd_g_i32x4_mul(v51, v110), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7)), v43), v116), v70), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v86, v110), base.Simd_g_i32x4_mul(v88, v110), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7)), v81), v116), v70))
	v129 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k14)
	v132 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k15)
	v138 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k16)
	v150 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v43, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v59, v132), base.Simd_g_i32x4_mul(v62, v132), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7))), v138), v70), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v81, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v96, v132), base.Simd_g_i32x4_mul(v98, v132), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7))), v138), v70))
	v151 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k17)
	v153 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v107, v108), base.Simd_g_i8x16_swizzle(v128, v129)), base.Simd_g_i8x16_swizzle(v150, v151))
	base.Simd_g_v128_store(m, l3, int32(80), v153)
	v156 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k18)
	v158 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k19)
	v161 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k20)
	v163 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v107, v156), base.Simd_g_i8x16_swizzle(v128, v158)), base.Simd_g_i8x16_swizzle(v150, v161))
	base.Simd_g_v128_store(m, l3, int32(64), v163)
	v166 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k21)
	v168 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k22)
	v171 = base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k23)
	v173 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v107, v166), base.Simd_g_i8x16_swizzle(v128, v168)), base.Simd_g_i8x16_swizzle(v150, v171))
	base.Simd_g_v128_store(m, l3, int32(48), v173)
	v176 = int32(0)
	v177 = base.Simd_g_v128_load64_zero(m, l0, v176)
	v179 = base.Simd_g_i8x16_shuffle2(v5, v177, base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k3))
	v185 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v179), v38), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v179), v38), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7))
	v187 = base.Simd_g_v128_load64_zero(m, l1, v176)
	v189 = base.Simd_g_i8x16_shuffle2(v5, v187, base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k3))
	v190 = base.Simd_g_i32x4_extend_low_i16x8_u(v189)
	v192 = base.Simd_g_i32x4_extend_high_i16x8_u(v189)
	v197 = base.Simd_g_v128_load64_zero(m, l2, v176)
	v199 = base.Simd_g_i8x16_shuffle2(v5, v197, base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k3))
	v200 = base.Simd_g_i32x4_extend_low_i16x8_u(v199)
	v202 = base.Simd_g_i32x4_extend_high_i16x8_u(v199)
	v211 = int32(8)
	v212 = base.Simd_g_v128_load64_zero(m, l0, v211)
	v214 = base.Simd_g_i8x16_shuffle2(v5, v212, base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k3))
	v220 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v214), v38), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v214), v38), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7))
	v222 = base.Simd_g_v128_load64_zero(m, l1, v211)
	v224 = base.Simd_g_i8x16_shuffle2(v5, v222, base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k3))
	v225 = base.Simd_g_i32x4_extend_low_i16x8_u(v224)
	v227 = base.Simd_g_i32x4_extend_high_i16x8_u(v224)
	v232 = base.Simd_g_v128_load64_zero(m, l2, v211)
	v234 = base.Simd_g_i8x16_shuffle2(v5, v232, base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k3))
	v235 = base.Simd_g_i32x4_extend_low_i16x8_u(v234)
	v237 = base.Simd_g_i32x4_extend_high_i16x8_u(v234)
	v246 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v185, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v190, v49), base.Simd_g_i32x4_mul(v192, v49), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v200, v60), base.Simd_g_i32x4_mul(v202, v60), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7)))), v68), v70), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v220, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v225, v49), base.Simd_g_i32x4_mul(v227, v49), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v235, v60), base.Simd_g_i32x4_mul(v237, v60), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7)))), v68), v70))
	v264 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v190, v110), base.Simd_g_i32x4_mul(v192, v110), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7)), v185), v116), v70), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v225, v110), base.Simd_g_i32x4_mul(v227, v110), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7)), v220), v116), v70))
	v283 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v185, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v200, v132), base.Simd_g_i32x4_mul(v202, v132), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7))), v138), v70), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v220, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v235, v132), base.Simd_g_i32x4_mul(v237, v132), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToBgr32_SSE41__k7))), v138), v70))
	v285 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v246, v108), base.Simd_g_i8x16_swizzle(v264, v129)), base.Simd_g_i8x16_swizzle(v283, v151))
	base.Simd_g_v128_store(m, l3, int32(32), v285)
	v292 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v246, v156), base.Simd_g_i8x16_swizzle(v264, v158)), base.Simd_g_i8x16_swizzle(v283, v161))
	base.Simd_g_v128_store(m, l3, v33, v292)
	v299 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v246, v166), base.Simd_g_i8x16_swizzle(v264, v168)), base.Simd_g_i8x16_swizzle(v283, v171))
	base.Simd_g_v128_store(m, l3, v176, v299)
	return
}

var F_VP8YuvToBgr32_SSE41__k0 = [2]uint64{0x0, 0x0}
var F_VP8YuvToBgr32_SSE41__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_VP8YuvToBgr32_SSE41__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_VP8YuvToBgr32_SSE41__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_VP8YuvToBgr32_SSE41__k4 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_VP8YuvToBgr32_SSE41__k5 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_VP8YuvToBgr32_SSE41__k6 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_VP8YuvToBgr32_SSE41__k7 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_VP8YuvToBgr32_SSE41__k8 = [2]uint64{0x191300001913, 0x191300001913}
var F_VP8YuvToBgr32_SSE41__k9 = [2]uint64{0x340800003408, 0x340800003408}
var F_VP8YuvToBgr32_SSE41__k10 = [2]uint64{0x2204220422042204, 0x2204220422042204}
var F_VP8YuvToBgr32_SSE41__k11 = [2]uint64{0x8f8f0c8f8f0b8f8f, 0x8f0f8f8f0e8f8f0d}
var F_VP8YuvToBgr32_SSE41__k12 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_VP8YuvToBgr32_SSE41__k13 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_VP8YuvToBgr32_SSE41__k14 = [2]uint64{0xd8f8f0c8f8f0b8f, 0x8f8f0f8f8f0e8f8f}
var F_VP8YuvToBgr32_SSE41__k15 = [2]uint64{0x662500006625, 0x662500006625}
var F_VP8YuvToBgr32_SSE41__k16 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_VP8YuvToBgr32_SSE41__k17 = [2]uint64{0x8f0c8f8f0b8f8f0a, 0xf8f8f0e8f8f0d8f}
var F_VP8YuvToBgr32_SSE41__k18 = [2]uint64{0x8f078f8f068f8f05, 0xa8f8f098f8f088f}
var F_VP8YuvToBgr32_SSE41__k19 = [2]uint64{0x8f8f078f8f068f8f, 0x8f0a8f8f098f8f08}
var F_VP8YuvToBgr32_SSE41__k20 = [2]uint64{0x78f8f068f8f058f, 0x8f8f098f8f088f8f}
var F_VP8YuvToBgr32_SSE41__k21 = [2]uint64{0x28f8f018f8f008f, 0x8f8f048f8f038f8f}
var F_VP8YuvToBgr32_SSE41__k22 = [2]uint64{0x8f028f8f018f8f00, 0x58f8f048f8f038f}
var F_VP8YuvToBgr32_SSE41__k23 = [2]uint64{0x8f8f018f8f008f8f, 0x8f048f8f038f8f02}

func F_VP8YuvToBgra32_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v28 base.V128
	_ = v28
	var v30 int32
	_ = v30
	var v31 base.V128
	_ = v31
	var v32 base.V128
	_ = v32
	var v33 base.V128
	_ = v33
	var v34 base.V128
	_ = v34
	var v35 base.V128
	_ = v35
	var v37 base.V128
	_ = v37
	var v39 base.V128
	_ = v39
	var v43 base.V128
	_ = v43
	var v45 base.V128
	_ = v45
	var v47 base.V128
	_ = v47
	var v52 base.V128
	_ = v52
	var v56 int32
	_ = v56
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v63 base.V128
	_ = v63
	var v64 base.V128
	_ = v64
	var v66 base.V128
	_ = v66
	var v75 base.V128
	_ = v75
	var v76 base.V128
	_ = v76
	var v81 base.V128
	_ = v81
	var v93 base.V128
	_ = v93
	var v95 base.V128
	_ = v95
	var v97 base.V128
	_ = v97
	var v99 base.V128
	_ = v99
	var v103 base.V128
	_ = v103
	v18 = l3
	v19 = int32(0)
	for {
		v28 = base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k0)
		v30 = int32(0)
		v31 = base.Simd_g_v128_load64_zero(m, l1+v19, v30)
		v32 = base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k1)
		v33 = base.Simd_g_i8x16_shuffle2(v28, v31, base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k3))
		v34 = base.Simd_g_i32x4_extend_low_i16x8_u(v33)
		v35 = base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k4)
		v37 = base.Simd_g_i32x4_extend_high_i16x8_u(v33)
		v39 = base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k5)
		v43 = base.Simd_g_v128_load64_zero(m, l0+v19, v30)
		v45 = base.Simd_g_i8x16_shuffle2(v28, v43, base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k3))
		v47 = base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k6)
		v52 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v45), v47), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v45), v47), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k8))
		v56 = int32(6)
		v60 = base.Simd_g_v128_load64_zero(m, l2+v19, v30)
		v62 = base.Simd_g_i8x16_shuffle2(v28, v60, base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k3))
		v63 = base.Simd_g_i32x4_extend_low_i16x8_u(v62)
		v64 = base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k9)
		v66 = base.Simd_g_i32x4_extend_high_i16x8_u(v62)
		v75 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v34, v35), base.Simd_g_i32x4_mul(v37, v35), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k8)), v52), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k10)), v56), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v52, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v63, v64), base.Simd_g_i32x4_mul(v66, v64), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k8))), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k11)), v56))
		v76 = base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k12)
		v81 = base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k13)
		v93 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v52, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v34, v76), base.Simd_g_i32x4_mul(v37, v76), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k8)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v63, v81), base.Simd_g_i32x4_mul(v66, v81), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k8)))), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k14)), v56), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k15))
		v95 = base.Simd_g_i8x16_shuffle2(v75, v93, base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k3))
		v97 = base.Simd_g_i8x16_shuffle2(v75, v93, base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k16), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k17))
		v99 = base.Simd_g_i8x16_shuffle2(v95, v97, base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k18), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k19))
		base.Simd_g_v128_store(m, v18, int32(16), v99)
		v103 = base.Simd_g_i8x16_shuffle2(v95, v97, base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k20), base.Simd_g_const(&F_VP8YuvToBgra32_SSE2__k21))
		base.Simd_g_v128_store(m, v18, v30, v103)
		if base.Ui32(v19) < base.Ui32(int32(24)) {
			v18 = v18 + int32(32)
			v19 = v19 + int32(8)
			continue
		} else {
			break
		}
		break
	}
	return
}

var F_VP8YuvToBgra32_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_VP8YuvToBgra32_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_VP8YuvToBgra32_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_VP8YuvToBgra32_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_VP8YuvToBgra32_SSE2__k4 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_VP8YuvToBgra32_SSE2__k5 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_VP8YuvToBgra32_SSE2__k6 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_VP8YuvToBgra32_SSE2__k7 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_VP8YuvToBgra32_SSE2__k8 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_VP8YuvToBgra32_SSE2__k9 = [2]uint64{0x662500006625, 0x662500006625}
var F_VP8YuvToBgra32_SSE2__k10 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_VP8YuvToBgra32_SSE2__k11 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_VP8YuvToBgra32_SSE2__k12 = [2]uint64{0x191300001913, 0x191300001913}
var F_VP8YuvToBgra32_SSE2__k13 = [2]uint64{0x340800003408, 0x340800003408}
var F_VP8YuvToBgra32_SSE2__k14 = [2]uint64{0x2204220422042204, 0x2204220422042204}
var F_VP8YuvToBgra32_SSE2__k15 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_VP8YuvToBgra32_SSE2__k16 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_VP8YuvToBgra32_SSE2__k17 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_VP8YuvToBgra32_SSE2__k18 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_VP8YuvToBgra32_SSE2__k19 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_VP8YuvToBgra32_SSE2__k20 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_VP8YuvToBgra32_SSE2__k21 = [2]uint64{0x302808001008080, 0x706808005048080}

func F_VP8YuvToRgb32_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	var v5 base.V128
	_ = v5
	var v37 int32
	_ = v37
	var v38 base.V128
	_ = v38
	var v39 base.V128
	_ = v39
	var v40 base.V128
	_ = v40
	var v41 base.V128
	_ = v41
	var v42 base.V128
	_ = v42
	var v44 base.V128
	_ = v44
	var v46 base.V128
	_ = v46
	var v49 base.V128
	_ = v49
	var v51 base.V128
	_ = v51
	var v53 base.V128
	_ = v53
	var v58 base.V128
	_ = v58
	var v60 base.V128
	_ = v60
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v65 base.V128
	_ = v65
	var v67 base.V128
	_ = v67
	var v68 base.V128
	_ = v68
	var v70 base.V128
	_ = v70
	var v75 base.V128
	_ = v75
	var v77 base.V128
	_ = v77
	var v83 base.V128
	_ = v83
	var v88 base.V128
	_ = v88
	var v89 base.V128
	_ = v89
	var v91 int32
	_ = v91
	var v92 base.V128
	_ = v92
	var v94 base.V128
	_ = v94
	var v95 base.V128
	_ = v95
	var v97 base.V128
	_ = v97
	var v102 base.V128
	_ = v102
	var v104 base.V128
	_ = v104
	var v110 base.V128
	_ = v110
	var v115 int32
	_ = v115
	var v116 base.V128
	_ = v116
	var v118 base.V128
	_ = v118
	var v119 base.V128
	_ = v119
	var v121 base.V128
	_ = v121
	var v126 base.V128
	_ = v126
	var v128 base.V128
	_ = v128
	var v134 base.V128
	_ = v134
	var v139 base.V128
	_ = v139
	var v141 base.V128
	_ = v141
	var v145 base.V128
	_ = v145
	var v147 base.V128
	_ = v147
	var v148 base.V128
	_ = v148
	var v149 base.V128
	_ = v149
	var v151 base.V128
	_ = v151
	var v156 base.V128
	_ = v156
	var v161 base.V128
	_ = v161
	var v163 base.V128
	_ = v163
	var v164 base.V128
	_ = v164
	var v166 base.V128
	_ = v166
	var v174 base.V128
	_ = v174
	var v178 base.V128
	_ = v178
	var v180 base.V128
	_ = v180
	var v181 base.V128
	_ = v181
	var v183 base.V128
	_ = v183
	var v192 base.V128
	_ = v192
	var v194 base.V128
	_ = v194
	var v195 base.V128
	_ = v195
	var v197 base.V128
	_ = v197
	var v205 base.V128
	_ = v205
	var v208 base.V128
	_ = v208
	var v211 base.V128
	_ = v211
	var v213 base.V128
	_ = v213
	var v218 base.V128
	_ = v218
	var v225 base.V128
	_ = v225
	var v242 base.V128
	_ = v242
	var v271 base.V128
	_ = v271
	var v274 base.V128
	_ = v274
	var v281 base.V128
	_ = v281
	var v284 base.V128
	_ = v284
	var v286 base.V128
	_ = v286
	var v291 base.V128
	_ = v291
	var v295 base.V128
	_ = v295
	var v297 base.V128
	_ = v297
	var v302 base.V128
	_ = v302
	var v305 base.V128
	_ = v305
	var v308 base.V128
	_ = v308
	var v313 base.V128
	_ = v313
	var v320 base.V128
	_ = v320
	var v323 base.V128
	_ = v323
	var v330 base.V128
	_ = v330
	var v333 base.V128
	_ = v333
	var v336 base.V128
	_ = v336
	var v341 base.V128
	_ = v341
	var v346 base.V128
	_ = v346
	var v351 base.V128
	_ = v351
	var v354 base.V128
	_ = v354
	var v357 base.V128
	_ = v357
	var v362 base.V128
	_ = v362
	var v367 base.V128
	_ = v367
	var v370 base.V128
	_ = v370
	var v375 base.V128
	_ = v375
	var v380 base.V128
	_ = v380
	var v385 base.V128
	_ = v385
	v5 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k0)
	v37 = int32(0)
	v38 = base.Simd_g_v128_load64_zero(m, l1, v37)
	v39 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k1)
	v40 = base.Simd_g_i8x16_shuffle2(v5, v38, base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k3))
	v41 = base.Simd_g_i32x4_extend_low_i16x8_u(v40)
	v42 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k4)
	v44 = base.Simd_g_i32x4_extend_high_i16x8_u(v40)
	v46 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k5)
	v49 = base.Simd_g_v128_load64_zero(m, l0, v37)
	v51 = base.Simd_g_i8x16_shuffle2(v5, v49, base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k3))
	v53 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k6)
	v58 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v51), v53), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v51), v53), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8))
	v60 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k9)
	v62 = int32(6)
	v64 = int32(8)
	v65 = base.Simd_g_v128_load64_zero(m, l1, v64)
	v67 = base.Simd_g_i8x16_shuffle2(v5, v65, base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k3))
	v68 = base.Simd_g_i32x4_extend_low_i16x8_u(v67)
	v70 = base.Simd_g_i32x4_extend_high_i16x8_u(v67)
	v75 = base.Simd_g_v128_load64_zero(m, l0, v64)
	v77 = base.Simd_g_i8x16_shuffle2(v5, v75, base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k3))
	v83 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v77), v53), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v77), v53), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8))
	v88 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v41, v42), base.Simd_g_i32x4_mul(v44, v42), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8)), v58), v60), v62), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v68, v42), base.Simd_g_i32x4_mul(v70, v42), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8)), v83), v60), v62))
	v89 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k10)
	v91 = int32(16)
	v92 = base.Simd_g_v128_load64_zero(m, l1, v91)
	v94 = base.Simd_g_i8x16_shuffle2(v5, v92, base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k3))
	v95 = base.Simd_g_i32x4_extend_low_i16x8_u(v94)
	v97 = base.Simd_g_i32x4_extend_high_i16x8_u(v94)
	v102 = base.Simd_g_v128_load64_zero(m, l0, v91)
	v104 = base.Simd_g_i8x16_shuffle2(v5, v102, base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k3))
	v110 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v104), v53), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v104), v53), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8))
	v115 = int32(24)
	v116 = base.Simd_g_v128_load64_zero(m, l1, v115)
	v118 = base.Simd_g_i8x16_shuffle2(v5, v116, base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k3))
	v119 = base.Simd_g_i32x4_extend_low_i16x8_u(v118)
	v121 = base.Simd_g_i32x4_extend_high_i16x8_u(v118)
	v126 = base.Simd_g_v128_load64_zero(m, l0, v115)
	v128 = base.Simd_g_i8x16_shuffle2(v5, v126, base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k3))
	v134 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v128), v53), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v128), v53), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8))
	v139 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v95, v42), base.Simd_g_i32x4_mul(v97, v42), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8)), v110), v60), v62), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v119, v42), base.Simd_g_i32x4_mul(v121, v42), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8)), v134), v60), v62))
	v141 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v88, v89), base.Simd_g_v128_and(v139, v89))
	v145 = base.Simd_g_v128_load64_zero(m, l2, v37)
	v147 = base.Simd_g_i8x16_shuffle2(v5, v145, base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k3))
	v148 = base.Simd_g_i32x4_extend_low_i16x8_u(v147)
	v149 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k11)
	v151 = base.Simd_g_i32x4_extend_high_i16x8_u(v147)
	v156 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k12)
	v161 = base.Simd_g_v128_load64_zero(m, l2, v64)
	v163 = base.Simd_g_i8x16_shuffle2(v5, v161, base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k3))
	v164 = base.Simd_g_i32x4_extend_low_i16x8_u(v163)
	v166 = base.Simd_g_i32x4_extend_high_i16x8_u(v163)
	v174 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v58, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v148, v149), base.Simd_g_i32x4_mul(v151, v149), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8))), v156), v62), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v83, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v164, v149), base.Simd_g_i32x4_mul(v166, v149), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8))), v156), v62))
	v178 = base.Simd_g_v128_load64_zero(m, l2, v91)
	v180 = base.Simd_g_i8x16_shuffle2(v5, v178, base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k3))
	v181 = base.Simd_g_i32x4_extend_low_i16x8_u(v180)
	v183 = base.Simd_g_i32x4_extend_high_i16x8_u(v180)
	v192 = base.Simd_g_v128_load64_zero(m, l2, v115)
	v194 = base.Simd_g_i8x16_shuffle2(v5, v192, base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k3))
	v195 = base.Simd_g_i32x4_extend_low_i16x8_u(v194)
	v197 = base.Simd_g_i32x4_extend_high_i16x8_u(v194)
	v205 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v110, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v181, v149), base.Simd_g_i32x4_mul(v183, v149), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8))), v156), v62), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v134, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v195, v149), base.Simd_g_i32x4_mul(v197, v149), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8))), v156), v62))
	v208 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v174, v64), base.Simd_g_i16x8_shr_u(v205, v64))
	v211 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v141, v64), base.Simd_g_i16x8_shr_u(v208, v64))
	v213 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k13)
	v218 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k14)
	v225 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k15)
	v242 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v58, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v41, v213), base.Simd_g_i32x4_mul(v44, v213), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v148, v218), base.Simd_g_i32x4_mul(v151, v218), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8)))), v225), v62), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v83, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v68, v213), base.Simd_g_i32x4_mul(v70, v213), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v164, v218), base.Simd_g_i32x4_mul(v166, v218), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8)))), v225), v62))
	v271 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v110, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v95, v213), base.Simd_g_i32x4_mul(v97, v213), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v181, v218), base.Simd_g_i32x4_mul(v183, v218), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8)))), v225), v62), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v134, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v119, v213), base.Simd_g_i32x4_mul(v121, v213), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v195, v218), base.Simd_g_i32x4_mul(v197, v218), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k7), base.Simd_g_const(&F_VP8YuvToRgb32_SSE2__k8)))), v225), v62))
	v274 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v242, v64), base.Simd_g_i16x8_shr_u(v271, v64))
	v281 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v88, v64), base.Simd_g_i16x8_shr_u(v139, v64))
	v284 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v274, v64), base.Simd_g_i16x8_shr_u(v281, v64))
	v286 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v211, v89), base.Simd_g_v128_and(v284, v89))
	v291 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v174, v89), base.Simd_g_v128_and(v205, v89))
	v295 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v242, v89), base.Simd_g_v128_and(v271, v89))
	v297 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v291, v89), base.Simd_g_v128_and(v295, v89))
	v302 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v141, v89), base.Simd_g_v128_and(v208, v89))
	v305 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v297, v64), base.Simd_g_i16x8_shr_u(v302, v64))
	v308 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v286, v64), base.Simd_g_i16x8_shr_u(v305, v64))
	v313 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v274, v89), base.Simd_g_v128_and(v281, v89))
	v320 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v291, v64), base.Simd_g_i16x8_shr_u(v295, v64))
	v323 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v313, v64), base.Simd_g_i16x8_shr_u(v320, v64))
	v330 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v211, v64), base.Simd_g_i16x8_shr_u(v284, v64))
	v333 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v323, v64), base.Simd_g_i16x8_shr_u(v330, v64))
	v336 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v308, v64), base.Simd_g_i16x8_shr_u(v333, v64))
	base.Simd_g_v128_store(m, l3, int32(80), v336)
	v341 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v323, v89), base.Simd_g_v128_and(v330, v89))
	v346 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v297, v89), base.Simd_g_v128_and(v302, v89))
	v351 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v313, v89), base.Simd_g_v128_and(v320, v89))
	v354 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v346, v64), base.Simd_g_i16x8_shr_u(v351, v64))
	v357 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v341, v64), base.Simd_g_i16x8_shr_u(v354, v64))
	base.Simd_g_v128_store(m, l3, int32(64), v357)
	v362 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v346, v89), base.Simd_g_v128_and(v351, v89))
	v367 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v286, v89), base.Simd_g_v128_and(v305, v89))
	v370 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v362, v64), base.Simd_g_i16x8_shr_u(v367, v64))
	base.Simd_g_v128_store(m, l3, int32(48), v370)
	v375 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v308, v89), base.Simd_g_v128_and(v333, v89))
	base.Simd_g_v128_store(m, l3, int32(32), v375)
	v380 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v341, v89), base.Simd_g_v128_and(v354, v89))
	base.Simd_g_v128_store(m, l3, v91, v380)
	v385 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v362, v89), base.Simd_g_v128_and(v367, v89))
	base.Simd_g_v128_store(m, l3, v37, v385)
	return
}

var F_VP8YuvToRgb32_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_VP8YuvToRgb32_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_VP8YuvToRgb32_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_VP8YuvToRgb32_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_VP8YuvToRgb32_SSE2__k4 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_VP8YuvToRgb32_SSE2__k5 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_VP8YuvToRgb32_SSE2__k6 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_VP8YuvToRgb32_SSE2__k7 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_VP8YuvToRgb32_SSE2__k8 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_VP8YuvToRgb32_SSE2__k9 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_VP8YuvToRgb32_SSE2__k10 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_VP8YuvToRgb32_SSE2__k11 = [2]uint64{0x662500006625, 0x662500006625}
var F_VP8YuvToRgb32_SSE2__k12 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_VP8YuvToRgb32_SSE2__k13 = [2]uint64{0x191300001913, 0x191300001913}
var F_VP8YuvToRgb32_SSE2__k14 = [2]uint64{0x340800003408, 0x340800003408}
var F_VP8YuvToRgb32_SSE2__k15 = [2]uint64{0x2204220422042204, 0x2204220422042204}

func F_VP8YuvToRgb32_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	var v5 base.V128
	_ = v5
	var v32 int32
	_ = v32
	var v33 base.V128
	_ = v33
	var v34 base.V128
	_ = v34
	var v35 base.V128
	_ = v35
	var v37 base.V128
	_ = v37
	var v41 base.V128
	_ = v41
	var v42 base.V128
	_ = v42
	var v44 base.V128
	_ = v44
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v48 base.V128
	_ = v48
	var v50 base.V128
	_ = v50
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
	var v67 base.V128
	_ = v67
	var v69 int32
	_ = v69
	var v71 int32
	_ = v71
	var v72 base.V128
	_ = v72
	var v74 base.V128
	_ = v74
	var v80 base.V128
	_ = v80
	var v82 base.V128
	_ = v82
	var v84 base.V128
	_ = v84
	var v85 base.V128
	_ = v85
	var v87 base.V128
	_ = v87
	var v92 base.V128
	_ = v92
	var v94 base.V128
	_ = v94
	var v95 base.V128
	_ = v95
	var v97 base.V128
	_ = v97
	var v106 base.V128
	_ = v106
	var v107 base.V128
	_ = v107
	var v109 base.V128
	_ = v109
	var v115 base.V128
	_ = v115
	var v127 base.V128
	_ = v127
	var v128 base.V128
	_ = v128
	var v131 base.V128
	_ = v131
	var v137 base.V128
	_ = v137
	var v149 base.V128
	_ = v149
	var v150 base.V128
	_ = v150
	var v152 base.V128
	_ = v152
	var v155 base.V128
	_ = v155
	var v157 base.V128
	_ = v157
	var v160 base.V128
	_ = v160
	var v162 base.V128
	_ = v162
	var v165 base.V128
	_ = v165
	var v167 base.V128
	_ = v167
	var v170 base.V128
	_ = v170
	var v172 base.V128
	_ = v172
	var v175 int32
	_ = v175
	var v176 base.V128
	_ = v176
	var v178 base.V128
	_ = v178
	var v184 base.V128
	_ = v184
	var v186 base.V128
	_ = v186
	var v188 base.V128
	_ = v188
	var v189 base.V128
	_ = v189
	var v191 base.V128
	_ = v191
	var v196 base.V128
	_ = v196
	var v198 base.V128
	_ = v198
	var v199 base.V128
	_ = v199
	var v201 base.V128
	_ = v201
	var v210 int32
	_ = v210
	var v211 base.V128
	_ = v211
	var v213 base.V128
	_ = v213
	var v219 base.V128
	_ = v219
	var v221 base.V128
	_ = v221
	var v223 base.V128
	_ = v223
	var v224 base.V128
	_ = v224
	var v226 base.V128
	_ = v226
	var v231 base.V128
	_ = v231
	var v233 base.V128
	_ = v233
	var v234 base.V128
	_ = v234
	var v236 base.V128
	_ = v236
	var v245 base.V128
	_ = v245
	var v263 base.V128
	_ = v263
	var v282 base.V128
	_ = v282
	var v284 base.V128
	_ = v284
	var v291 base.V128
	_ = v291
	var v298 base.V128
	_ = v298
	v5 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k0)
	v32 = int32(16)
	v33 = base.Simd_g_v128_load64_zero(m, l0, v32)
	v34 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k1)
	v35 = base.Simd_g_i8x16_shuffle2(v5, v33, base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k3))
	v37 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k4)
	v41 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k5)
	v42 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v35), v37), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v35), v37), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7))
	v44 = base.Simd_g_v128_load64_zero(m, l1, v32)
	v46 = base.Simd_g_i8x16_shuffle2(v5, v44, base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k3))
	v47 = base.Simd_g_i32x4_extend_low_i16x8_u(v46)
	v48 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k8)
	v50 = base.Simd_g_i32x4_extend_high_i16x8_u(v46)
	v55 = base.Simd_g_v128_load64_zero(m, l2, v32)
	v57 = base.Simd_g_i8x16_shuffle2(v5, v55, base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k3))
	v58 = base.Simd_g_i32x4_extend_low_i16x8_u(v57)
	v59 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k9)
	v61 = base.Simd_g_i32x4_extend_high_i16x8_u(v57)
	v67 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k10)
	v69 = int32(6)
	v71 = int32(24)
	v72 = base.Simd_g_v128_load64_zero(m, l0, v71)
	v74 = base.Simd_g_i8x16_shuffle2(v5, v72, base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k3))
	v80 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v74), v37), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v74), v37), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7))
	v82 = base.Simd_g_v128_load64_zero(m, l1, v71)
	v84 = base.Simd_g_i8x16_shuffle2(v5, v82, base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k3))
	v85 = base.Simd_g_i32x4_extend_low_i16x8_u(v84)
	v87 = base.Simd_g_i32x4_extend_high_i16x8_u(v84)
	v92 = base.Simd_g_v128_load64_zero(m, l2, v71)
	v94 = base.Simd_g_i8x16_shuffle2(v5, v92, base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k3))
	v95 = base.Simd_g_i32x4_extend_low_i16x8_u(v94)
	v97 = base.Simd_g_i32x4_extend_high_i16x8_u(v94)
	v106 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v42, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v47, v48), base.Simd_g_i32x4_mul(v50, v48), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v58, v59), base.Simd_g_i32x4_mul(v61, v59), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7)))), v67), v69), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v80, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v85, v48), base.Simd_g_i32x4_mul(v87, v48), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v95, v59), base.Simd_g_i32x4_mul(v97, v59), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7)))), v67), v69))
	v107 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k11)
	v109 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k12)
	v115 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k13)
	v127 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v42, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v58, v109), base.Simd_g_i32x4_mul(v61, v109), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7))), v115), v69), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v80, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v95, v109), base.Simd_g_i32x4_mul(v97, v109), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7))), v115), v69))
	v128 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k14)
	v131 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k15)
	v137 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k16)
	v149 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v47, v131), base.Simd_g_i32x4_mul(v50, v131), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7)), v42), v137), v69), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v85, v131), base.Simd_g_i32x4_mul(v87, v131), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7)), v80), v137), v69))
	v150 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k17)
	v152 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v106, v107), base.Simd_g_i8x16_swizzle(v127, v128)), base.Simd_g_i8x16_swizzle(v149, v150))
	base.Simd_g_v128_store(m, l3, int32(80), v152)
	v155 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k18)
	v157 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k19)
	v160 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k20)
	v162 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v106, v155), base.Simd_g_i8x16_swizzle(v127, v157)), base.Simd_g_i8x16_swizzle(v149, v160))
	base.Simd_g_v128_store(m, l3, int32(64), v162)
	v165 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k21)
	v167 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k22)
	v170 = base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k23)
	v172 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v106, v165), base.Simd_g_i8x16_swizzle(v127, v167)), base.Simd_g_i8x16_swizzle(v149, v170))
	base.Simd_g_v128_store(m, l3, int32(48), v172)
	v175 = int32(0)
	v176 = base.Simd_g_v128_load64_zero(m, l0, v175)
	v178 = base.Simd_g_i8x16_shuffle2(v5, v176, base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k3))
	v184 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v178), v37), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v178), v37), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7))
	v186 = base.Simd_g_v128_load64_zero(m, l1, v175)
	v188 = base.Simd_g_i8x16_shuffle2(v5, v186, base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k3))
	v189 = base.Simd_g_i32x4_extend_low_i16x8_u(v188)
	v191 = base.Simd_g_i32x4_extend_high_i16x8_u(v188)
	v196 = base.Simd_g_v128_load64_zero(m, l2, v175)
	v198 = base.Simd_g_i8x16_shuffle2(v5, v196, base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k3))
	v199 = base.Simd_g_i32x4_extend_low_i16x8_u(v198)
	v201 = base.Simd_g_i32x4_extend_high_i16x8_u(v198)
	v210 = int32(8)
	v211 = base.Simd_g_v128_load64_zero(m, l0, v210)
	v213 = base.Simd_g_i8x16_shuffle2(v5, v211, base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k3))
	v219 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v213), v37), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v213), v37), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7))
	v221 = base.Simd_g_v128_load64_zero(m, l1, v210)
	v223 = base.Simd_g_i8x16_shuffle2(v5, v221, base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k3))
	v224 = base.Simd_g_i32x4_extend_low_i16x8_u(v223)
	v226 = base.Simd_g_i32x4_extend_high_i16x8_u(v223)
	v231 = base.Simd_g_v128_load64_zero(m, l2, v210)
	v233 = base.Simd_g_i8x16_shuffle2(v5, v231, base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k2), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k3))
	v234 = base.Simd_g_i32x4_extend_low_i16x8_u(v233)
	v236 = base.Simd_g_i32x4_extend_high_i16x8_u(v233)
	v245 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v184, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v189, v48), base.Simd_g_i32x4_mul(v191, v48), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v199, v59), base.Simd_g_i32x4_mul(v201, v59), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7)))), v67), v69), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v219, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v224, v48), base.Simd_g_i32x4_mul(v226, v48), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v234, v59), base.Simd_g_i32x4_mul(v236, v59), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7)))), v67), v69))
	v263 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v184, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v199, v109), base.Simd_g_i32x4_mul(v201, v109), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7))), v115), v69), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v219, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v234, v109), base.Simd_g_i32x4_mul(v236, v109), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7))), v115), v69))
	v282 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v189, v131), base.Simd_g_i32x4_mul(v191, v131), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7)), v184), v137), v69), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v224, v131), base.Simd_g_i32x4_mul(v226, v131), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k6), base.Simd_g_const(&F_VP8YuvToRgb32_SSE41__k7)), v219), v137), v69))
	v284 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v245, v107), base.Simd_g_i8x16_swizzle(v263, v128)), base.Simd_g_i8x16_swizzle(v282, v150))
	base.Simd_g_v128_store(m, l3, int32(32), v284)
	v291 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v245, v155), base.Simd_g_i8x16_swizzle(v263, v157)), base.Simd_g_i8x16_swizzle(v282, v160))
	base.Simd_g_v128_store(m, l3, v32, v291)
	v298 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v245, v165), base.Simd_g_i8x16_swizzle(v263, v167)), base.Simd_g_i8x16_swizzle(v282, v170))
	base.Simd_g_v128_store(m, l3, v175, v298)
	return
}

var F_VP8YuvToRgb32_SSE41__k0 = [2]uint64{0x0, 0x0}
var F_VP8YuvToRgb32_SSE41__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_VP8YuvToRgb32_SSE41__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_VP8YuvToRgb32_SSE41__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_VP8YuvToRgb32_SSE41__k4 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_VP8YuvToRgb32_SSE41__k5 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_VP8YuvToRgb32_SSE41__k6 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_VP8YuvToRgb32_SSE41__k7 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_VP8YuvToRgb32_SSE41__k8 = [2]uint64{0x191300001913, 0x191300001913}
var F_VP8YuvToRgb32_SSE41__k9 = [2]uint64{0x340800003408, 0x340800003408}
var F_VP8YuvToRgb32_SSE41__k10 = [2]uint64{0x2204220422042204, 0x2204220422042204}
var F_VP8YuvToRgb32_SSE41__k11 = [2]uint64{0x8f8f0c8f8f0b8f8f, 0x8f0f8f8f0e8f8f0d}
var F_VP8YuvToRgb32_SSE41__k12 = [2]uint64{0x662500006625, 0x662500006625}
var F_VP8YuvToRgb32_SSE41__k13 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_VP8YuvToRgb32_SSE41__k14 = [2]uint64{0xd8f8f0c8f8f0b8f, 0x8f8f0f8f8f0e8f8f}
var F_VP8YuvToRgb32_SSE41__k15 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_VP8YuvToRgb32_SSE41__k16 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_VP8YuvToRgb32_SSE41__k17 = [2]uint64{0x8f0c8f8f0b8f8f0a, 0xf8f8f0e8f8f0d8f}
var F_VP8YuvToRgb32_SSE41__k18 = [2]uint64{0x8f078f8f068f8f05, 0xa8f8f098f8f088f}
var F_VP8YuvToRgb32_SSE41__k19 = [2]uint64{0x8f8f078f8f068f8f, 0x8f0a8f8f098f8f08}
var F_VP8YuvToRgb32_SSE41__k20 = [2]uint64{0x78f8f068f8f058f, 0x8f8f098f8f088f8f}
var F_VP8YuvToRgb32_SSE41__k21 = [2]uint64{0x28f8f018f8f008f, 0x8f8f048f8f038f8f}
var F_VP8YuvToRgb32_SSE41__k22 = [2]uint64{0x8f028f8f018f8f00, 0x58f8f048f8f038f}
var F_VP8YuvToRgb32_SSE41__k23 = [2]uint64{0x8f8f018f8f008f8f, 0x8f048f8f038f8f02}

func F_VP8YuvToRgb56532_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v28 base.V128
	_ = v28
	var v30 int32
	_ = v30
	var v31 base.V128
	_ = v31
	var v32 base.V128
	_ = v32
	var v33 base.V128
	_ = v33
	var v35 base.V128
	_ = v35
	var v39 base.V128
	_ = v39
	var v40 base.V128
	_ = v40
	var v43 base.V128
	_ = v43
	var v45 base.V128
	_ = v45
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v49 base.V128
	_ = v49
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
	var v69 int32
	_ = v69
	var v70 base.V128
	_ = v70
	var v71 base.V128
	_ = v71
	var v76 base.V128
	_ = v76
	var v85 base.V128
	_ = v85
	var v90 base.V128
	_ = v90
	var v99 base.V128
	_ = v99
	var v101 int32
	_ = v101
	var v111 base.V128
	_ = v111
	v18 = l3
	v19 = int32(0)
	for {
		v28 = base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k0)
		v30 = int32(0)
		v31 = base.Simd_g_v128_load64_zero(m, l0+v19, v30)
		v32 = base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k1)
		v33 = base.Simd_g_i8x16_shuffle2(v28, v31, base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k3))
		v35 = base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k4)
		v39 = base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k5)
		v40 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v33), v35), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v33), v35), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k6), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k7))
		v43 = base.Simd_g_v128_load64_zero(m, l1+v19, v30)
		v45 = base.Simd_g_i8x16_shuffle2(v28, v43, base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k3))
		v46 = base.Simd_g_i32x4_extend_low_i16x8_u(v45)
		v47 = base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k8)
		v49 = base.Simd_g_i32x4_extend_high_i16x8_u(v45)
		v55 = base.Simd_g_v128_load64_zero(m, l2+v19, v30)
		v57 = base.Simd_g_i8x16_shuffle2(v28, v55, base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k3))
		v58 = base.Simd_g_i32x4_extend_low_i16x8_u(v57)
		v59 = base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k9)
		v61 = base.Simd_g_i32x4_extend_high_i16x8_u(v57)
		v69 = int32(6)
		v70 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v40, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v46, v47), base.Simd_g_i32x4_mul(v49, v47), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k6), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v58, v59), base.Simd_g_i32x4_mul(v61, v59), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k6), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k7)))), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k10)), v69)
		v71 = base.Simd_g_i8x16_narrow_i16x8_u(v70, v70)
		v76 = base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k11)
		v85 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v40, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v58, v76), base.Simd_g_i32x4_mul(v61, v76), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k6), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k7))), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k12)), v69)
		v90 = base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k13)
		v99 = base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v46, v90), base.Simd_g_i32x4_mul(v49, v90), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k6), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k7)), v40), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k14)), v69)
		v101 = int32(3)
		v111 = base.Simd_g_i8x16_shuffle2(base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(v71, int32(5)), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k15)), base.Simd_g_v128_and(base.Simd_g_i8x16_narrow_i16x8_u(v85, v85), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k16))), base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_narrow_i16x8_u(v99, v99), v101), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k17)), base.Simd_g_v128_and(base.Simd_g_i16x8_shl(v71, v101), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k18))), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgb56532_SSE2__k3))
		base.Simd_g_v128_store(m, v18, v30, v111)
		if base.Ui32(v19) < base.Ui32(int32(24)) {
			v18 = v18 + int32(16)
			v19 = v19 + int32(8)
			continue
		} else {
			break
		}
		break
	}
	return
}

var F_VP8YuvToRgb56532_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_VP8YuvToRgb56532_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_VP8YuvToRgb56532_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_VP8YuvToRgb56532_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_VP8YuvToRgb56532_SSE2__k4 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_VP8YuvToRgb56532_SSE2__k5 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_VP8YuvToRgb56532_SSE2__k6 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_VP8YuvToRgb56532_SSE2__k7 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_VP8YuvToRgb56532_SSE2__k8 = [2]uint64{0x191300001913, 0x191300001913}
var F_VP8YuvToRgb56532_SSE2__k9 = [2]uint64{0x340800003408, 0x340800003408}
var F_VP8YuvToRgb56532_SSE2__k10 = [2]uint64{0x2204220422042204, 0x2204220422042204}
var F_VP8YuvToRgb56532_SSE2__k11 = [2]uint64{0x662500006625, 0x662500006625}
var F_VP8YuvToRgb56532_SSE2__k12 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_VP8YuvToRgb56532_SSE2__k13 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_VP8YuvToRgb56532_SSE2__k14 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_VP8YuvToRgb56532_SSE2__k15 = [2]uint64{0x707070707070707, 0x707070707070707}
var F_VP8YuvToRgb56532_SSE2__k16 = [2]uint64{0xf8f8f8f8f8f8f8f8, 0xf8f8f8f8f8f8f8f8}
var F_VP8YuvToRgb56532_SSE2__k17 = [2]uint64{0x1f1f1f1f1f1f1f1f, 0x1f1f1f1f1f1f1f1f}
var F_VP8YuvToRgb56532_SSE2__k18 = [2]uint64{0xe0e0e0e0e0e0e0e0, 0xe0e0e0e0e0e0e0e0}

func F_VP8YuvToRgba32_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v28 base.V128
	_ = v28
	var v30 int32
	_ = v30
	var v31 base.V128
	_ = v31
	var v32 base.V128
	_ = v32
	var v33 base.V128
	_ = v33
	var v35 base.V128
	_ = v35
	var v39 base.V128
	_ = v39
	var v40 base.V128
	_ = v40
	var v43 base.V128
	_ = v43
	var v45 base.V128
	_ = v45
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v49 base.V128
	_ = v49
	var v56 int32
	_ = v56
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v63 base.V128
	_ = v63
	var v64 base.V128
	_ = v64
	var v66 base.V128
	_ = v66
	var v75 base.V128
	_ = v75
	var v76 base.V128
	_ = v76
	var v81 base.V128
	_ = v81
	var v93 base.V128
	_ = v93
	var v95 base.V128
	_ = v95
	var v97 base.V128
	_ = v97
	var v99 base.V128
	_ = v99
	var v103 base.V128
	_ = v103
	v18 = l3
	v19 = int32(0)
	for {
		v28 = base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k0)
		v30 = int32(0)
		v31 = base.Simd_g_v128_load64_zero(m, l0+v19, v30)
		v32 = base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k1)
		v33 = base.Simd_g_i8x16_shuffle2(v28, v31, base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k3))
		v35 = base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k4)
		v39 = base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k5)
		v40 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v33), v35), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v33), v35), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k7))
		v43 = base.Simd_g_v128_load64_zero(m, l2+v19, v30)
		v45 = base.Simd_g_i8x16_shuffle2(v28, v43, base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k3))
		v46 = base.Simd_g_i32x4_extend_low_i16x8_u(v45)
		v47 = base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k8)
		v49 = base.Simd_g_i32x4_extend_high_i16x8_u(v45)
		v56 = int32(6)
		v60 = base.Simd_g_v128_load64_zero(m, l1+v19, v30)
		v62 = base.Simd_g_i8x16_shuffle2(v28, v60, base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k3))
		v63 = base.Simd_g_i32x4_extend_low_i16x8_u(v62)
		v64 = base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k9)
		v66 = base.Simd_g_i32x4_extend_high_i16x8_u(v62)
		v75 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v40, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v46, v47), base.Simd_g_i32x4_mul(v49, v47), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k7))), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k10)), v56), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v63, v64), base.Simd_g_i32x4_mul(v66, v64), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k7)), v40), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k11)), v56))
		v76 = base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k12)
		v81 = base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k13)
		v93 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v40, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v63, v76), base.Simd_g_i32x4_mul(v66, v76), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v46, v81), base.Simd_g_i32x4_mul(v49, v81), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k6), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k7)))), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k14)), v56), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k15))
		v95 = base.Simd_g_i8x16_shuffle2(v75, v93, base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k3))
		v97 = base.Simd_g_i8x16_shuffle2(v75, v93, base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k16), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k17))
		v99 = base.Simd_g_i8x16_shuffle2(v95, v97, base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k18), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k19))
		base.Simd_g_v128_store(m, v18, int32(16), v99)
		v103 = base.Simd_g_i8x16_shuffle2(v95, v97, base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k20), base.Simd_g_const(&F_VP8YuvToRgba32_SSE2__k21))
		base.Simd_g_v128_store(m, v18, v30, v103)
		if base.Ui32(v19) < base.Ui32(int32(24)) {
			v18 = v18 + int32(32)
			v19 = v19 + int32(8)
			continue
		} else {
			break
		}
		break
	}
	return
}

var F_VP8YuvToRgba32_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_VP8YuvToRgba32_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_VP8YuvToRgba32_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_VP8YuvToRgba32_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_VP8YuvToRgba32_SSE2__k4 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_VP8YuvToRgba32_SSE2__k5 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_VP8YuvToRgba32_SSE2__k6 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_VP8YuvToRgba32_SSE2__k7 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_VP8YuvToRgba32_SSE2__k8 = [2]uint64{0x662500006625, 0x662500006625}
var F_VP8YuvToRgba32_SSE2__k9 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_VP8YuvToRgba32_SSE2__k10 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_VP8YuvToRgba32_SSE2__k11 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_VP8YuvToRgba32_SSE2__k12 = [2]uint64{0x191300001913, 0x191300001913}
var F_VP8YuvToRgba32_SSE2__k13 = [2]uint64{0x340800003408, 0x340800003408}
var F_VP8YuvToRgba32_SSE2__k14 = [2]uint64{0x2204220422042204, 0x2204220422042204}
var F_VP8YuvToRgba32_SSE2__k15 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_VP8YuvToRgba32_SSE2__k16 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_VP8YuvToRgba32_SSE2__k17 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_VP8YuvToRgba32_SSE2__k18 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_VP8YuvToRgba32_SSE2__k19 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_VP8YuvToRgba32_SSE2__k20 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_VP8YuvToRgba32_SSE2__k21 = [2]uint64{0x302808001008080, 0x706808005048080}

func F_VP8YuvToRgba444432_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v26 base.V128
	_ = v26
	var v28 int32
	_ = v28
	var v29 base.V128
	_ = v29
	var v30 base.V128
	_ = v30
	var v31 base.V128
	_ = v31
	var v33 base.V128
	_ = v33
	var v37 base.V128
	_ = v37
	var v38 base.V128
	_ = v38
	var v41 base.V128
	_ = v41
	var v43 base.V128
	_ = v43
	var v44 base.V128
	_ = v44
	var v45 base.V128
	_ = v45
	var v47 base.V128
	_ = v47
	var v54 int32
	_ = v54
	var v58 base.V128
	_ = v58
	var v60 base.V128
	_ = v60
	var v61 base.V128
	_ = v61
	var v62 base.V128
	_ = v62
	var v64 base.V128
	_ = v64
	var v68 base.V128
	_ = v68
	var v79 base.V128
	_ = v79
	var v80 base.V128
	_ = v80
	var v91 base.V128
	_ = v91
	var v102 base.V128
	_ = v102
	v17 = l3
	v18 = int32(0)
	for {
		v26 = base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k0)
		v28 = int32(0)
		v29 = base.Simd_g_v128_load64_zero(m, l0+v18, v28)
		v30 = base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k1)
		v31 = base.Simd_g_i8x16_shuffle2(v26, v29, base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k3))
		v33 = base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k4)
		v37 = base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k5)
		v38 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v31), v33), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v31), v33), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k6), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k7))
		v41 = base.Simd_g_v128_load64_zero(m, l2+v18, v28)
		v43 = base.Simd_g_i8x16_shuffle2(v26, v41, base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k3))
		v44 = base.Simd_g_i32x4_extend_low_i16x8_u(v43)
		v45 = base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k8)
		v47 = base.Simd_g_i32x4_extend_high_i16x8_u(v43)
		v54 = int32(6)
		v58 = base.Simd_g_v128_load64_zero(m, l1+v18, v28)
		v60 = base.Simd_g_i8x16_shuffle2(v26, v58, base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k3))
		v61 = base.Simd_g_i32x4_extend_low_i16x8_u(v60)
		v62 = base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k9)
		v64 = base.Simd_g_i32x4_extend_high_i16x8_u(v60)
		v68 = base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k10)
		v79 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v38, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v44, v45), base.Simd_g_i32x4_mul(v47, v45), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k6), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k7))), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k11)), v54), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v38, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v61, v62), base.Simd_g_i32x4_mul(v64, v62), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k6), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v44, v68), base.Simd_g_i32x4_mul(v47, v68), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k6), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k7)))), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k12)), v54))
		v80 = base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k13)
		v91 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v61, v80), base.Simd_g_i32x4_mul(v64, v80), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k6), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k7)), v38), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k14)), v54), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k15))
		v102 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_shuffle2(v79, v91, base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k16), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k17)), int32(4)), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k18)), base.Simd_g_v128_and(base.Simd_g_i8x16_shuffle2(v79, v91, base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k2), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k3)), base.Simd_g_const(&F_VP8YuvToRgba444432_SSE2__k19)))
		base.Simd_g_v128_store(m, v17, v28, v102)
		if base.Ui32(v18) < base.Ui32(int32(24)) {
			v17 = v17 + int32(16)
			v18 = v18 + int32(8)
			continue
		} else {
			break
		}
		break
	}
	return
}

var F_VP8YuvToRgba444432_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_VP8YuvToRgba444432_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_VP8YuvToRgba444432_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_VP8YuvToRgba444432_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_VP8YuvToRgba444432_SSE2__k4 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_VP8YuvToRgba444432_SSE2__k5 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_VP8YuvToRgba444432_SSE2__k6 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_VP8YuvToRgba444432_SSE2__k7 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_VP8YuvToRgba444432_SSE2__k8 = [2]uint64{0x662500006625, 0x662500006625}
var F_VP8YuvToRgba444432_SSE2__k9 = [2]uint64{0x191300001913, 0x191300001913}
var F_VP8YuvToRgba444432_SSE2__k10 = [2]uint64{0x340800003408, 0x340800003408}
var F_VP8YuvToRgba444432_SSE2__k11 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_VP8YuvToRgba444432_SSE2__k12 = [2]uint64{0x2204220422042204, 0x2204220422042204}
var F_VP8YuvToRgba444432_SSE2__k13 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_VP8YuvToRgba444432_SSE2__k14 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_VP8YuvToRgba444432_SSE2__k15 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_VP8YuvToRgba444432_SSE2__k16 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_VP8YuvToRgba444432_SSE2__k17 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_VP8YuvToRgba444432_SSE2__k18 = [2]uint64{0xf0f0f0f0f0f0f0f, 0xf0f0f0f0f0f0f0f}
var F_VP8YuvToRgba444432_SSE2__k19 = [2]uint64{0xf0f0f0f0f0f0f0f0, 0xf0f0f0f0f0f0f0f0}

func F_VR4_SSE2(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 base.V128
	_ = v2
	var v3 int32
	_ = v3
	var v12 int32
	_ = v12
	var v14 base.V128
	_ = v14
	var v16 base.V128
	_ = v16
	var v17 base.V128
	_ = v17
	var v18 base.V128
	_ = v18
	var v25 int32
	_ = v25
	var v30 int32
	_ = v30
	var v33 base.V128
	_ = v33
	var v39 base.V128
	_ = v39
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v53 int32
	_ = v53
	var v56 int32
	_ = v56
	var v61 int32
	_ = v61
	var v69 int32
	_ = v69
	v2 = base.Simd_g_const(&F_VR4_SSE2__k0)
	v3 = int32(0)
	v12 = l0 + int32(-33)
	v14 = base.Simd_g_v128_load64_zero(m, v12, v3)
	v16 = base.Simd_g_i8x16_shuffle2(v14, v2, base.Simd_g_const(&F_VR4_SSE2__k1), base.Simd_g_const(&F_VR4_SSE2__k2))
	v17 = base.Simd_g_i8x16_avgr_u(v14, v16)
	v18 = base.Simd_g_const(&F_VR4_SSE2__k3)
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(64), base.Simd_g_i8x16_shuffle2(v2, v17, base.Simd_g_const(&F_VR4_SSE2__k4), base.Simd_g_const(&F_VR4_SSE2__k5)))
	v25 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v12))))
	v30 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-1)))))
	v33 = base.Simd_g_i16x8_replace_lane_l0(base.Simd_g_i8x16_shuffle2(v2, v14, base.Simd_g_const(&F_VR4_SSE2__k4), base.Simd_g_const(&F_VR4_SSE2__k5)), v25<<(uint(int32(8))%32)|v30)
	v39 = base.Simd_g_i8x16_avgr_u(v14, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_avgr_u(v33, v16), base.Simd_g_v128_and(base.Simd_g_v128_xor(v16, v33), base.Simd_g_const(&F_VR4_SSE2__k6))))
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(96), base.Simd_g_i8x16_shuffle2(v2, v39, base.Simd_g_const(&F_VR4_SSE2__k4), base.Simd_g_const(&F_VR4_SSE2__k5)))
	base.Simd_g_v128_store32_lane_l0(m, l0, v3, v17)
	v48 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+31)))
	v49 = int32(1)
	v53 = int32(2)
	v56 = int32(base.Ui32(v25+(v48+v30<<(uint(v49)%32))+v53) >> (uint(v53) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+64)) = uint8(v56)
	base.Simd_g_v128_store32_lane_l0(m, l0, int32(32), v39)
	v61 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+63)))
	v69 = int32(base.Ui32(v30+v61+v48<<(uint(v49)%32)+v53) >> (uint(v53) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+96)) = uint8(v69)
	return
}

var F_VR4_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_VR4_SSE2__k1 = [2]uint64{0x807060504030201, 0x800f0e0d0c0b0a09}
var F_VR4_SSE2__k2 = [2]uint64{0x8080808080808080, 0x80808080808080}
var F_VR4_SSE2__k3 = [2]uint64{0x161514131211100f, 0x1e1d1c1b1a191817}
var F_VR4_SSE2__k4 = [2]uint64{0x808080808080800f, 0x8080808080808080}
var F_VR4_SSE2__k5 = [2]uint64{0x605040302010080, 0xe0d0c0b0a090807}
var F_VR4_SSE2__k6 = [2]uint64{0x101010101010101, 0x101010101010101}

func F_VectorMismatch_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v17 base.V128
	_ = v17
	var v19 base.V128
	_ = v19
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v29 int32
	_ = v29
	var v31 int32
	_ = v31
	var v32 base.V128
	_ = v32
	var v35 base.V128
	_ = v35
	var v37 int32
	_ = v37
	var v40 int32
	_ = v40
	var v44 int32
	_ = v44
	var v50 int32
	_ = v50
	var v53 int32
	_ = v53
	var v54 base.V128
	_ = v54
	var v56 base.V128
	_ = v56
	var v66 int32
	_ = v66
	var v67 base.V128
	_ = v67
	var v69 base.V128
	_ = v69
	var v74 int32
	_ = v74
	var v78 int32
	_ = v78
	var v83 int32
	_ = v83
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v89 int32
	_ = v89
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v95 int32
	_ = v95
	var v100 int32
	_ = v100
	var v105 int32
	_ = v105
	if l2 < int32(12) {
		goto L2
	} else {
		goto L3
	}
L1:
	;
	if l2 <= v78 {
		v105 = v78
		goto L18
	} else {
		goto L19
	}
L2:
	;
	v50 = int32(0)
	if l2 < int32(4) {
		v78 = v50
		goto L1
	} else {
		goto L12
	}
L3:
	;
	v13 = int32(0)
	v14 = l1
	v15 = l0
	goto L4
L4:
	;
	v16 = int32(0)
	v17 = base.Simd_g_v128_load(m, v15, v16)
	v19 = base.Simd_g_v128_load(m, v14, v16)
	if base.Simd_g_i8x16_bitmask(base.Simd_g_i32x4_eq(v17, v19)) != int32(_a_F_VectorMismatch_SSE2_0) {
		v78 = v13
		goto L1
	} else {
		goto L6
	}
L6:
	;
	v25 = v13 + int32(8)
	v27 = v13 + int32(4)
	v29 = v27 << (uint(int32(2)) % 32)
	v31 = int32(0)
	v32 = base.Simd_g_v128_load(m, l0+v29, v31)
	v35 = base.Simd_g_v128_load(m, l1+v29, v31)
	v37 = base.Simd_g_i8x16_bitmask(base.Simd_g_i32x4_eq(v32, v35))
	if v37 == int32(_a_F_VectorMismatch_SSE2_0) {
		goto L7
	} else {
		goto L8
	}
L7:
	;
	v40 = v25
	goto L9
L8:
	;
	v40 = v27
	goto L9
L9:
	;
	if v37 != int32(_a_F_VectorMismatch_SSE2_0) {
		v78 = v40
		goto L1
	} else {
		goto L10
	}
L10:
	;
	v44 = v25 << (uint(int32(2)) % 32)
	if v40+int32(12) < l2 {
		v13 = v40
		v14 = l1 + v44
		v15 = l0 + v44
		goto L4
	} else {
		goto L11
	}
L11:
	;
	v78 = v40
	goto L1
L12:
	;
	v53 = int32(0)
	v54 = base.Simd_g_v128_load(m, l0, v53)
	v56 = base.Simd_g_v128_load(m, l1, v53)
	if base.Simd_g_i8x16_bitmask(base.Simd_g_i32x4_eq(v54, v56)) != int32(_a_F_VectorMismatch_SSE2_0) {
		v78 = v50
		goto L1
	} else {
		goto L13
	}
L13:
	;
	if base.Ui32(l2) < base.Ui32(int32(8)) {
		v78 = int32(4)
		goto L1
	} else {
		goto L14
	}
L14:
	;
	v66 = int32(16)
	v67 = base.Simd_g_v128_load(m, l0, v66)
	v69 = base.Simd_g_v128_load(m, l1, v66)
	if base.Simd_g_i8x16_bitmask(base.Simd_g_i32x4_eq(v67, v69)) == int32(_a_F_VectorMismatch_SSE2_0) {
		goto L15
	} else {
		goto L16
	}
L15:
	;
	v74 = int32(8)
	goto L17
L16:
	;
	v74 = int32(4)
	goto L17
L17:
	;
	v78 = v74
	goto L1
L18:
	;
	return v105
L19:
	;
	v83 = v78 << (uint(int32(2)) % 32)
	v86 = l0 + v83
	v87 = l1 + v83
	v89 = v78
	goto L20
L20:
	;
	v92 = *(*int32)(unsafe.Add(mBase, uint32(v86)))
	v93 = *(*int32)(unsafe.Add(mBase, uint32(v87)))
	if v92 != v93 {
		v105 = v89
		goto L18
	} else {
		goto L22
	}
L21:
	;
	v105 = l2
	goto L18
L22:
	;
	v95 = int32(4)
	v100 = v89 + int32(1)
	if l2 != v100 {
		v86 = v86 + v95
		v87 = v87 + v95
		v89 = v100
		goto L20
	} else {
		goto L23
	}
L23:
	;
	goto L21
}
func F_VerticalFilter_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v38 int32
	_ = v38
	var v39 base.V128
	_ = v39
	var v41 base.V128
	_ = v41
	var v42 base.V128
	_ = v42
	var v45 int32
	_ = v45
	var v50 int32
	_ = v50
	var v57 int32
	_ = v57
	var v76 int32
	_ = v76
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v84 int32
	_ = v84
	var v85 int32
	_ = v85
	var v90 int32
	_ = v90
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v122 int32
	_ = v122
	var v128 int32
	_ = v128
	var v151 int32
	_ = v151
	var v154 int32
	_ = v154
	var v156 int32
	_ = v156
	var v160 int32
	_ = v160
	var v164 int32
	_ = v164
	var v167 int32
	_ = v167
	var v184 int32
	_ = v184
	var v188 int32
	_ = v188
	var v189 base.V128
	_ = v189
	var v192 base.V128
	_ = v192
	var v193 base.V128
	_ = v193
	var v197 int32
	_ = v197
	var v208 int32
	_ = v208
	var v220 int32
	_ = v220
	var v222 int32
	_ = v222
	var v223 int32
	_ = v223
	var v227 int32
	_ = v227
	var v228 int32
	_ = v228
	var v229 int32
	_ = v229
	var v231 int32
	_ = v231
	var v236 int32
	_ = v236
	var v251 int32
	_ = v251
	v12 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	*(*uint8)(unsafe.Add(mBase, uint32(l4))) = uint8(v12)
	if l1 < int32(2) {
	} else {
		if base.Ui32(l1) < base.Ui32(int32(17)) {
			v57 = int32(0)
			v76 = v57 + l4 + int32(1)
			v78 = v57 ^ int32(-1) + l1
			v79 = l0 + v57
			for {
				v81 = int32(1)
				v82 = v79 + v81
				v83 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v82))))
				v84 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v79))))
				v85 = v83 - v84
				*(*uint8)(unsafe.Add(mBase, uint32(v76))) = uint8(v85)
				v90 = v78 + int32(-1)
				if v90 != 0 {
					v76 = v76 + v81
					v78 = v90
					v79 = v82
					continue
				} else {
					break
				}
				break
			}
		} else {
			v22 = l1 + int32(-1)
			v24 = v22 & int32(-16)
			v31 = l4 + int32(1)
			v33 = v24
			v34 = l0
			for {
				v38 = int32(0)
				v39 = base.Simd_g_v128_load_rng(m, v34+int32(1), v38, int32(-1), int32(17))
				v41 = base.Simd_g_v128_load_nc(m, v34, v38)
				v42 = base.Simd_g_i8x16_sub(v39, v41)
				base.Simd_g_v128_store(m, v31, v38, v42)
				v45 = int32(16)
				v50 = v33 + int32(-16)
				if v50 != 0 {
					v31 = v31 + v45
					v33 = v50
					v34 = v34 + v45
					continue
				} else {
					break
				}
				break
			}
			if v22 == v24 {
			} else {
				v57 = v24
				v76 = v57 + l4 + int32(1)
				v78 = v57 ^ int32(-1) + l1
				v79 = l0 + v57
				for {
					v81 = int32(1)
					v82 = v79 + v81
					v83 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v82))))
					v84 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v79))))
					v85 = v83 - v84
					*(*uint8)(unsafe.Add(mBase, uint32(v76))) = uint8(v85)
					v90 = v78 + int32(-1)
					if v90 != 0 {
						v76 = v76 + v81
						v78 = v90
						v79 = v82
						continue
					} else {
						break
					}
					break
				}
			}
		}
	}
	if l2 < int32(2) {
	} else {
		if int32(0) < l1 {
			v156 = l1 & int32(2147483632)
			v160 = l0
			v164 = l4
			v167 = int32(1)
			for {
				if base.Ui32(l1) < base.Ui32(int32(16)) {
					v208 = int32(0)
					v220 = v160 + v208
					v222 = l1 - v208
					v223 = l3 + v208
					for {
						v227 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v160+v223))))
						v228 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v220))))
						v229 = v227 - v228
						*(*uint8)(unsafe.Add(mBase, uint32(v164+v223))) = uint8(v229)
						v231 = int32(1)
						v236 = v222 + int32(-1)
						if v236 != 0 {
							v220 = v220 + v231
							v222 = v236
							v223 = v223 + v231
							continue
						} else {
							break
						}
						break
					}
				} else {
					v184 = int32(0)
					for {
						v188 = int32(0)
						v189 = base.Simd_g_v128_load(m, v160+l3+v184, v188)
						v192 = base.Simd_g_v128_load(m, v160+v184, v188)
						v193 = base.Simd_g_i8x16_sub(v189, v192)
						base.Simd_g_v128_store(m, v164+l3+v184, v188, v193)
						v197 = v184 + int32(16)
						if v156 != v197 {
							v184 = v197
							continue
						} else {
							break
						}
						break
					}
					if v156 == l1 {
					} else {
						v208 = v156
						v220 = v160 + v208
						v222 = l1 - v208
						v223 = l3 + v208
						for {
							v227 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v160+v223))))
							v228 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v220))))
							v229 = v227 - v228
							*(*uint8)(unsafe.Add(mBase, uint32(v164+v223))) = uint8(v229)
							v231 = int32(1)
							v236 = v222 + int32(-1)
							if v236 != 0 {
								v220 = v220 + v231
								v222 = v236
								v223 = v223 + v231
								continue
							} else {
								break
							}
							break
						}
					}
				}
				v251 = v167 + int32(1)
				if v251 != l2 {
					v160 = v160 + l3
					v164 = v164 + l3
					v167 = v251
					continue
				} else {
					break
				}
				break
			}
		} else {
			v107 = l2 + int32(-1)
			v108 = int32(7)
			v109 = v107 & v108
			if base.Ui32(l2+int32(-2)) < base.Ui32(v108) {
			} else {
				v122 = v107 & int32(-8)
				for {
					v128 = v122 + int32(-8)
					if v128 != 0 {
						v122 = v128
						continue
					} else {
						break
					}
					break
				}
			}
			if v109 == int32(0) {
			} else {
				v151 = v109
				for {
					v154 = v151 + int32(-1)
					if v154 != 0 {
						v151 = v154
						continue
					} else {
						break
					}
					break
				}
			}
		}
	}
	return
}
func F_VerticalFilter_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v26 int32
	_ = v26
	var v29 int32
	_ = v29
	var v31 int32
	_ = v31
	var v36 int32
	_ = v36
	var v54 int32
	_ = v54
	var v56 int32
	_ = v56
	var v72 int32
	_ = v72
	var v73 int32
	_ = v73
	var v75 int32
	_ = v75
	var v78 int32
	_ = v78
	var v79 base.V128
	_ = v79
	var v83 base.V128
	_ = v83
	var v84 base.V128
	_ = v84
	var v87 int32
	_ = v87
	var v92 base.V128
	_ = v92
	var v94 base.V128
	_ = v94
	var v95 base.V128
	_ = v95
	var v98 int32
	_ = v98
	var v103 base.V128
	_ = v103
	var v107 base.V128
	_ = v107
	var v108 base.V128
	_ = v108
	var v111 int32
	_ = v111
	var v116 base.V128
	_ = v116
	var v120 base.V128
	_ = v120
	var v121 base.V128
	_ = v121
	var v125 int32
	_ = v125
	var v127 int32
	_ = v127
	var v135 int32
	_ = v135
	var v155 int32
	_ = v155
	var v157 int32
	_ = v157
	var v160 int32
	_ = v160
	var v161 int32
	_ = v161
	var v162 base.V128
	_ = v162
	var v164 base.V128
	_ = v164
	var v165 base.V128
	_ = v165
	var v168 int32
	_ = v168
	var v169 base.V128
	_ = v169
	var v173 base.V128
	_ = v173
	var v174 base.V128
	_ = v174
	var v186 int32
	_ = v186
	var v207 int32
	_ = v207
	var v215 int32
	_ = v215
	var v223 int32
	_ = v223
	var v226 int32
	_ = v226
	var v227 int32
	_ = v227
	var v244 int32
	_ = v244
	var v245 base.V128
	_ = v245
	var v247 base.V128
	_ = v247
	var v248 base.V128
	_ = v248
	var v251 int32
	_ = v251
	var v256 int32
	_ = v256
	var v265 int32
	_ = v265
	var v296 int32
	_ = v296
	var v297 int32
	_ = v297
	var v300 int32
	_ = v300
	var v315 int32
	_ = v315
	var v316 int32
	_ = v316
	var v317 int32
	_ = v317
	var v318 int32
	_ = v318
	var v319 int32
	_ = v319
	var v324 int32
	_ = v324
	var v353 int32
	_ = v353
	var v360 int32
	_ = v360
	var v361 int32
	_ = v361
	var v367 int32
	_ = v367
	var v369 int32
	_ = v369
	var v371 int32
	_ = v371
	var v376 int32
	_ = v376
	var v380 int32
	_ = v380
	var v384 int32
	_ = v384
	var v396 int32
	_ = v396
	var v397 int32
	_ = v397
	var v398 int32
	_ = v398
	var v399 int32
	_ = v399
	var v400 int32
	_ = v400
	var v401 int32
	_ = v401
	var v414 int32
	_ = v414
	var v415 int32
	_ = v415
	var v417 int32
	_ = v417
	var v418 int32
	_ = v418
	var v433 int32
	_ = v433
	var v434 int32
	_ = v434
	var v435 int32
	_ = v435
	var v436 base.V128
	_ = v436
	var v437 int32
	_ = v437
	var v439 base.V128
	_ = v439
	var v440 base.V128
	_ = v440
	var v443 int32
	_ = v443
	var v448 base.V128
	_ = v448
	var v452 base.V128
	_ = v452
	var v453 base.V128
	_ = v453
	var v456 int32
	_ = v456
	var v463 int32
	_ = v463
	var v474 int32
	_ = v474
	var v476 int32
	_ = v476
	var v477 int32
	_ = v477
	var v493 int32
	_ = v493
	var v494 int32
	_ = v494
	var v495 base.V128
	_ = v495
	var v497 base.V128
	_ = v497
	var v498 base.V128
	_ = v498
	var v501 int32
	_ = v501
	var v508 int32
	_ = v508
	var v517 int32
	_ = v517
	var v545 int32
	_ = v545
	var v546 int32
	_ = v546
	var v549 int32
	_ = v549
	var v566 int32
	_ = v566
	var v567 int32
	_ = v567
	var v568 int32
	_ = v568
	var v570 int32
	_ = v570
	var v575 int32
	_ = v575
	var v607 int32
	_ = v607
	var v612 int32
	_ = v612
	var v613 int32
	_ = v613
	var v614 int32
	_ = v614
	var v628 int32
	_ = v628
	var v647 int32
	_ = v647
	var v681 int32
	_ = v681
	var v701 int32
	_ = v701
	var v703 int32
	_ = v703
	var v707 int32
	_ = v707
	var v711 int32
	_ = v711
	var v712 int32
	_ = v712
	var v742 int32
	_ = v742
	var v763 int32
	_ = v763
	var v764 base.V128
	_ = v764
	var v767 base.V128
	_ = v767
	var v768 base.V128
	_ = v768
	var v772 int32
	_ = v772
	var v785 int32
	_ = v785
	var v809 int32
	_ = v809
	var v810 int32
	_ = v810
	var v813 int32
	_ = v813
	var v830 int32
	_ = v830
	var v831 int32
	_ = v831
	var v832 int32
	_ = v832
	var v834 int32
	_ = v834
	var v839 int32
	_ = v839
	var v868 int32
	_ = v868
	v26 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	*(*uint8)(unsafe.Add(mBase, uint32(l4))) = uint8(v26)
	v29 = l1 + int32(-1)
	v31 = v29 & int32(-32)
	if int32(1) <= v31 {
		v36 = v31 + int32(-1)
		if v36 != int32(31) {
			v54 = int32(0)
			v56 = (int32(base.Ui32(v36)>>(uint(int32(5))%32)) + int32(1)) & int32(268435454)
			for {
				v72 = l4 + v54
				v73 = int32(17)
				v75 = l0 + v54
				v78 = int32(0)
				v79 = base.Simd_g_v128_load_rng(m, v75+v73, v78, int32(-1), int32(17))
				v83 = base.Simd_g_v128_load_nc(m, v75+int32(16), v78)
				v84 = base.Simd_g_i8x16_sub(v79, v83)
				base.Simd_g_v128_store(m, v72+v73, v78, v84)
				v87 = int32(1)
				v92 = base.Simd_g_v128_load_rng(m, v75+v87, v78, int32(-1), int32(17))
				v94 = base.Simd_g_v128_load_nc(m, v75, v78)
				v95 = base.Simd_g_i8x16_sub(v92, v94)
				base.Simd_g_v128_store(m, v72+v87, v78, v95)
				v98 = int32(49)
				v103 = base.Simd_g_v128_load_rng(m, v75+v98, v78, int32(-1), int32(17))
				v107 = base.Simd_g_v128_load_nc(m, v75+int32(48), v78)
				v108 = base.Simd_g_i8x16_sub(v103, v107)
				base.Simd_g_v128_store(m, v72+v98, v78, v108)
				v111 = int32(33)
				v116 = base.Simd_g_v128_load_rng(m, v75+v111, v78, int32(-1), int32(17))
				v120 = base.Simd_g_v128_load_nc(m, v75+int32(32), v78)
				v121 = base.Simd_g_i8x16_sub(v116, v120)
				base.Simd_g_v128_store(m, v72+v111, v78, v121)
				v125 = v54 + int32(64)
				v127 = v56 + int32(-2)
				if v127 != 0 {
					v54 = v125
					v56 = v127
					continue
				} else {
					break
				}
				break
			}
			v135 = v125
		} else {
			v135 = int32(0)
		}
		if v36&int32(32) != 0 {
			v186 = v135
		} else {
			v155 = int32(1)
			v157 = l4 + v155 + v135
			v160 = l0 + v155 + v135
			v161 = int32(16)
			v162 = base.Simd_g_v128_load_rng(m, v160, v161, int32(15), int32(17))
			v164 = base.Simd_g_v128_load_nc(m, v160, int32(15))
			v165 = base.Simd_g_i8x16_sub(v162, v164)
			base.Simd_g_v128_store(m, v157, v161, v165)
			v168 = int32(0)
			v169 = base.Simd_g_v128_load(m, v160, v168)
			v173 = base.Simd_g_v128_load(m, v160+int32(-1), v168)
			v174 = base.Simd_g_i8x16_sub(v169, v173)
			base.Simd_g_v128_store(m, v157, v168, v174)
			v186 = v135 + int32(32)
		}
	} else {
		v186 = int32(0)
	}
	if v29 <= v186 {
	} else {
		v207 = v186 ^ int32(-1) + l1
		if base.Ui32(v207) <= base.Ui32(int32(15)) {
			v265 = v186
			v296 = l0 + v265
			v297 = v265 ^ int32(-1) + l1
			v300 = v265 + l4 + int32(1)
			for {
				v315 = int32(1)
				v316 = v296 + v315
				v317 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v316))))
				v318 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296))))
				v319 = v317 - v318
				*(*uint8)(unsafe.Add(mBase, uint32(v300))) = uint8(v319)
				v324 = v297 + int32(-1)
				if v324 != 0 {
					v296 = v316
					v297 = v324
					v300 = v300 + v315
					continue
				} else {
					break
				}
				break
			}
		} else {
			v215 = v207 & int32(-16)
			v223 = l0 + v186
			v226 = v215
			v227 = v186 + l4 + int32(1)
			for {
				v244 = int32(0)
				v245 = base.Simd_g_v128_load_rng(m, v223+int32(1), v244, int32(-1), int32(17))
				v247 = base.Simd_g_v128_load_nc(m, v223, v244)
				v248 = base.Simd_g_i8x16_sub(v245, v247)
				base.Simd_g_v128_store(m, v227, v244, v248)
				v251 = int32(16)
				v256 = v226 + int32(-16)
				if v256 != 0 {
					v223 = v223 + v251
					v226 = v256
					v227 = v227 + v251
					continue
				} else {
					break
				}
				break
			}
			if v207 == v215 {
			} else {
				v265 = v186 + v215
				v296 = l0 + v265
				v297 = v265 ^ int32(-1) + l1
				v300 = v265 + l4 + int32(1)
				for {
					v315 = int32(1)
					v316 = v296 + v315
					v317 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v316))))
					v318 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296))))
					v319 = v317 - v318
					*(*uint8)(unsafe.Add(mBase, uint32(v300))) = uint8(v319)
					v324 = v297 + int32(-1)
					if v324 != 0 {
						v296 = v316
						v297 = v324
						v300 = v300 + v315
						continue
					} else {
						break
					}
					break
				}
			}
		}
	}
	if l2 < int32(2) {
	} else {
		v353 = l1 & int32(-32)
		if v353 < int32(1) {
			if int32(0) < l1 {
				v703 = l1 & int32(2147483632)
				v707 = l0
				v711 = l4
				v712 = int32(1)
				for {
					if base.Ui32(l1) < base.Ui32(int32(16)) {
						v785 = int32(0)
						v809 = l3 + v785
						v810 = v707 + v785
						v813 = l1 - v785
						for {
							v830 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v707+v809))))
							v831 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v810))))
							v832 = v830 - v831
							*(*uint8)(unsafe.Add(mBase, uint32(v711+v809))) = uint8(v832)
							v834 = int32(1)
							v839 = v813 + int32(-1)
							if v839 != 0 {
								v809 = v809 + v834
								v810 = v810 + v834
								v813 = v839
								continue
							} else {
								break
							}
							break
						}
					} else {
						v742 = int32(0)
						for {
							v763 = int32(0)
							v764 = base.Simd_g_v128_load(m, v707+l3+v742, v763)
							v767 = base.Simd_g_v128_load(m, v707+v742, v763)
							v768 = base.Simd_g_i8x16_sub(v764, v767)
							base.Simd_g_v128_store(m, v711+l3+v742, v763, v768)
							v772 = v742 + int32(16)
							if v703 != v772 {
								v742 = v772
								continue
							} else {
								break
							}
							break
						}
						if v703 == l1 {
						} else {
							v785 = v703
							v809 = l3 + v785
							v810 = v707 + v785
							v813 = l1 - v785
							for {
								v830 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v707+v809))))
								v831 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v810))))
								v832 = v830 - v831
								*(*uint8)(unsafe.Add(mBase, uint32(v711+v809))) = uint8(v832)
								v834 = int32(1)
								v839 = v813 + int32(-1)
								if v839 != 0 {
									v809 = v809 + v834
									v810 = v810 + v834
									v813 = v839
									continue
								} else {
									break
								}
								break
							}
						}
					}
					v868 = v712 + int32(1)
					if v868 != l2 {
						v707 = v707 + l3
						v711 = v711 + l3
						v712 = v868
						continue
					} else {
						break
					}
					break
				}
			} else {
				v612 = l2 + int32(-1)
				v613 = int32(7)
				v614 = v612 & v613
				if base.Ui32(l2+int32(-2)) < base.Ui32(v613) {
				} else {
					v628 = v612 & int32(-8)
					for {
						v647 = v628 + int32(-8)
						if v647 != 0 {
							v628 = v647
							continue
						} else {
							break
						}
						break
					}
				}
				if v614 == int32(0) {
				} else {
					v681 = v614
					for {
						v701 = v681 + int32(-1)
						if v701 != 0 {
							v681 = v701
							continue
						} else {
							break
						}
						break
					}
				}
			}
		} else {
			v360 = int32(-32)
			v361 = (v353 + int32(-1)) & v360
			v367 = l1 - v361 + v360
			v369 = l1 & int32(15)
			v371 = int32(32)
			v376 = l3 + v371
			v380 = l0
			v384 = l4
			v396 = l0 + v371
			v397 = l4 + l3
			v398 = l0 + l3
			v399 = l0 + v376
			v400 = l4 + v376
			v401 = int32(1)
			for {
				v414 = int32(0)
				v415 = v399
				v417 = v400
				v418 = v396
				for {
					v433 = v397 + v414
					v434 = v398 + v414
					v435 = int32(0)
					v436 = base.Simd_g_v128_load(m, v434, v435)
					v437 = v380 + v414
					v439 = base.Simd_g_v128_load(m, v437, v435)
					v440 = base.Simd_g_i8x16_sub(v436, v439)
					base.Simd_g_v128_store(m, v433, v435, v440)
					v443 = int32(16)
					v448 = base.Simd_g_v128_load(m, v434+v443, v435)
					v452 = base.Simd_g_v128_load(m, v437+v443, v435)
					v453 = base.Simd_g_i8x16_sub(v448, v452)
					base.Simd_g_v128_store(m, v433+v443, v435, v453)
					v456 = int32(32)
					v463 = v414 + v456
					if v463 < v353 {
						v414 = v463
						v415 = v415 + v456
						v417 = v417 + v456
						v418 = v418 + v456
						continue
					} else {
						break
					}
					break
				}
				if l1 <= v463 {
				} else {
					if base.Ui32(v367) <= base.Ui32(int32(15)) {
						v517 = v463
						v545 = l1 - v517
						v546 = l3 + v517
						v549 = v380 + v517
						for {
							v566 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v380+v546))))
							v567 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v549))))
							v568 = v566 - v567
							*(*uint8)(unsafe.Add(mBase, uint32(v384+v546))) = uint8(v568)
							v570 = int32(1)
							v575 = v545 + int32(-1)
							if v575 != 0 {
								v545 = v575
								v546 = v546 + v570
								v549 = v549 + v570
								continue
							} else {
								break
							}
							break
						}
					} else {
						v474 = v415
						v476 = v367 - v369
						v477 = v418
						v493 = v417
						for {
							v494 = int32(0)
							v495 = base.Simd_g_v128_load(m, v474, v494)
							v497 = base.Simd_g_v128_load(m, v477, v494)
							v498 = base.Simd_g_i8x16_sub(v495, v497)
							base.Simd_g_v128_store(m, v493, v494, v498)
							v501 = int32(16)
							v508 = v476 + int32(-16)
							if v508 != 0 {
								v474 = v474 + v501
								v476 = v508
								v477 = v477 + v501
								v493 = v493 + v501
								continue
							} else {
								break
							}
							break
						}
						if v369 == int32(0) {
						} else {
							v517 = l1&int32(-16) - v361 + v360 + v463
							v545 = l1 - v517
							v546 = l3 + v517
							v549 = v380 + v517
							for {
								v566 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v380+v546))))
								v567 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v549))))
								v568 = v566 - v567
								*(*uint8)(unsafe.Add(mBase, uint32(v384+v546))) = uint8(v568)
								v570 = int32(1)
								v575 = v545 + int32(-1)
								if v575 != 0 {
									v545 = v575
									v546 = v546 + v570
									v549 = v549 + v570
									continue
								} else {
									break
								}
								break
							}
						}
					}
				}
				v607 = v401 + int32(1)
				if v607 != l2 {
					v380 = v380 + l3
					v384 = v384 + l3
					v396 = v396 + l3
					v397 = v397 + l3
					v398 = v398 + l3
					v399 = v399 + l3
					v400 = v400 + l3
					v401 = v607
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
func F_VerticalUnfilter_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v24 int32
	_ = v24
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v35 base.V128
	_ = v35
	var v37 base.V128
	_ = v37
	var v38 base.V128
	_ = v38
	var v41 int32
	_ = v41
	var v48 int32
	_ = v48
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v68 int32
	_ = v68
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v83 int32
	_ = v83
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v90 int32
	_ = v90
	var v94 int32
	_ = v94
	var v95 int32
	_ = v95
	var v98 int32
	_ = v98
	var v100 int32
	_ = v100
	var v106 int32
	_ = v106
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v119 int32
	_ = v119
	var v121 int32
	_ = v121
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v125 int32
	_ = v125
	var v130 int32
	_ = v130
	var v135 int32
	_ = v135
	var v146 int32
	_ = v146
	var v148 int32
	_ = v148
	var v149 int32
	_ = v149
	var v153 int32
	_ = v153
	var v157 int32
	_ = v157
	var v158 int32
	_ = v158
	var v159 int32
	_ = v159
	var v160 int32
	_ = v160
	var v166 int32
	_ = v166
	var v167 int32
	_ = v167
	var v168 int32
	_ = v168
	var v169 int32
	_ = v169
	var v170 int32
	_ = v170
	var v171 int32
	_ = v171
	var v173 int32
	_ = v173
	var v177 int32
	_ = v177
	var v180 int32
	_ = v180
	var v181 int32
	_ = v181
	var v183 int32
	_ = v183
	var v190 int32
	_ = v190
	if l0 == int32(0) {
		if l3 < int32(1) {
		} else {
			v53 = l3 & int32(3)
			v54 = int32(0)
			if base.Ui32(l3) < base.Ui32(int32(4)) {
				v100 = v54
				v106 = v54
			} else {
				v60 = int32(0)
				v62 = v60
				v68 = v60
				for {
					v71 = l2 + v62
					v72 = l1 + v62
					v73 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v72))))
					v74 = v73 + v68
					*(*uint8)(unsafe.Add(mBase, uint32(v71))) = uint8(v74)
					v76 = int32(1)
					v80 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v72+v76))))
					v81 = v80 + v74
					*(*uint8)(unsafe.Add(mBase, uint32(v71+v76))) = uint8(v81)
					v83 = int32(2)
					v87 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v72+v83))))
					v88 = v87 + v81
					*(*uint8)(unsafe.Add(mBase, uint32(v71+v83))) = uint8(v88)
					v90 = int32(3)
					v94 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v72+v90))))
					v95 = v94 + v88
					*(*uint8)(unsafe.Add(mBase, uint32(v71+v90))) = uint8(v95)
					v98 = v62 + int32(4)
					if l3&int32(2147483644) != v98 {
						v62 = v98
						v68 = v95
						continue
					} else {
						break
					}
					break
				}
				v100 = v98
				v106 = v95
			}
			if v53 == int32(0) {
			} else {
				v114 = l1 + v100
				v115 = l2 + v100
				v119 = v106
				v121 = v53
				for {
					v122 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v114))))
					v123 = v122 + v119
					*(*uint8)(unsafe.Add(mBase, uint32(v115))) = uint8(v123)
					v125 = int32(1)
					v130 = v121 + int32(-1)
					if v130 != 0 {
						v114 = v114 + v125
						v115 = v115 + v125
						v119 = v123
						v121 = v130
						continue
					} else {
						break
					}
					break
				}
			}
		}
	} else {
		if l3 < int32(1) {
		} else {
			v14 = int32(0)
			if base.Ui32(l3) <= base.Ui32(int32(15)) {
				v135 = v14
				if l3&int32(1) == int32(0) {
					v153 = v135
				} else {
					v146 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v135))))
					v148 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v135))))
					v149 = v146 + v148
					*(*uint8)(unsafe.Add(mBase, uint32(l2+v135))) = uint8(v149)
					v153 = v135 | int32(1)
				}
				if v135 == l3+int32(-1) {
				} else {
					v157 = l0
					v158 = l1
					v159 = l2
					v160 = l3
					for {
						v166 = v159 + v153
						v167 = v158 + v153
						v168 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v167))))
						v169 = v157 + v153
						v170 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v169))))
						v171 = v168 + v170
						*(*uint8)(unsafe.Add(mBase, uint32(v166))) = uint8(v171)
						v173 = int32(1)
						v177 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v167+v173))))
						v180 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v169+v173))))
						v181 = v177 + v180
						*(*uint8)(unsafe.Add(mBase, uint32(v166+v173))) = uint8(v181)
						v183 = int32(2)
						v190 = v160 + int32(-2)
						if v153 != v190 {
							v157 = v157 + v183
							v158 = v158 + v183
							v159 = v159 + v183
							v160 = v190
							continue
						} else {
							break
						}
						break
					}
				}
			} else {
				if base.Ui32(l2-l0) < base.Ui32(int32(16)) {
					v135 = v14
					if l3&int32(1) == int32(0) {
						v153 = v135
					} else {
						v146 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v135))))
						v148 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v135))))
						v149 = v146 + v148
						*(*uint8)(unsafe.Add(mBase, uint32(l2+v135))) = uint8(v149)
						v153 = v135 | int32(1)
					}
					if v135 == l3+int32(-1) {
					} else {
						v157 = l0
						v158 = l1
						v159 = l2
						v160 = l3
						for {
							v166 = v159 + v153
							v167 = v158 + v153
							v168 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v167))))
							v169 = v157 + v153
							v170 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v169))))
							v171 = v168 + v170
							*(*uint8)(unsafe.Add(mBase, uint32(v166))) = uint8(v171)
							v173 = int32(1)
							v177 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v167+v173))))
							v180 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v169+v173))))
							v181 = v177 + v180
							*(*uint8)(unsafe.Add(mBase, uint32(v166+v173))) = uint8(v181)
							v183 = int32(2)
							v190 = v160 + int32(-2)
							if v153 != v190 {
								v157 = v157 + v183
								v158 = v158 + v183
								v159 = v159 + v183
								v160 = v190
								continue
							} else {
								break
							}
							break
						}
					}
				} else {
					if base.Ui32(l2-l1) < base.Ui32(int32(16)) {
						v135 = v14
						if l3&int32(1) == int32(0) {
							v153 = v135
						} else {
							v146 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v135))))
							v148 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v135))))
							v149 = v146 + v148
							*(*uint8)(unsafe.Add(mBase, uint32(l2+v135))) = uint8(v149)
							v153 = v135 | int32(1)
						}
						if v135 == l3+int32(-1) {
						} else {
							v157 = l0
							v158 = l1
							v159 = l2
							v160 = l3
							for {
								v166 = v159 + v153
								v167 = v158 + v153
								v168 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v167))))
								v169 = v157 + v153
								v170 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v169))))
								v171 = v168 + v170
								*(*uint8)(unsafe.Add(mBase, uint32(v166))) = uint8(v171)
								v173 = int32(1)
								v177 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v167+v173))))
								v180 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v169+v173))))
								v181 = v177 + v180
								*(*uint8)(unsafe.Add(mBase, uint32(v166+v173))) = uint8(v181)
								v183 = int32(2)
								v190 = v160 + int32(-2)
								if v153 != v190 {
									v157 = v157 + v183
									v158 = v158 + v183
									v159 = v159 + v183
									v160 = v190
									continue
								} else {
									break
								}
								break
							}
						}
					} else {
						v24 = l3 & int32(2147483632)
						v30 = v24
						v31 = l2
						v32 = l1
						v33 = l0
						for {
							v34 = int32(0)
							v35 = base.Simd_g_v128_load(m, v32, v34)
							v37 = base.Simd_g_v128_load(m, v33, v34)
							v38 = base.Simd_g_i8x16_add(v35, v37)
							base.Simd_g_v128_store(m, v31, v34, v38)
							v41 = int32(16)
							v48 = v30 + int32(-16)
							if v48 != 0 {
								v30 = v48
								v31 = v31 + v41
								v32 = v32 + v41
								v33 = v33 + v41
								continue
							} else {
								break
							}
							break
						}
						if v24 != l3 {
							v135 = v24
							if l3&int32(1) == int32(0) {
								v153 = v135
							} else {
								v146 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v135))))
								v148 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v135))))
								v149 = v146 + v148
								*(*uint8)(unsafe.Add(mBase, uint32(l2+v135))) = uint8(v149)
								v153 = v135 | int32(1)
							}
							if v135 == l3+int32(-1) {
							} else {
								v157 = l0
								v158 = l1
								v159 = l2
								v160 = l3
								for {
									v166 = v159 + v153
									v167 = v158 + v153
									v168 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v167))))
									v169 = v157 + v153
									v170 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v169))))
									v171 = v168 + v170
									*(*uint8)(unsafe.Add(mBase, uint32(v166))) = uint8(v171)
									v173 = int32(1)
									v177 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v167+v173))))
									v180 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v169+v173))))
									v181 = v177 + v180
									*(*uint8)(unsafe.Add(mBase, uint32(v166+v173))) = uint8(v181)
									v183 = int32(2)
									v190 = v160 + int32(-2)
									if v153 != v190 {
										v157 = v157 + v183
										v158 = v158 + v183
										v159 = v159 + v183
										v160 = v190
										continue
									} else {
										break
									}
									break
								}
							}
						} else {
						}
					}
				}
			}
		}
	}
	return
}
func F_VerticalUnfilter_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v21 int32
	_ = v21
	var v26 int32
	_ = v26
	var v33 int32
	_ = v33
	var v36 base.V128
	_ = v36
	var v43 base.V128
	_ = v43
	var v45 int32
	_ = v45
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v50 base.V128
	_ = v50
	var v53 base.V128
	_ = v53
	var v56 base.V128
	_ = v56
	var v74 int32
	_ = v74
	var v86 int32
	_ = v86
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v94 int32
	_ = v94
	var v96 int32
	_ = v96
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v109 int32
	_ = v109
	var v111 int32
	_ = v111
	var v112 int32
	_ = v112
	var v116 int32
	_ = v116
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v146 int32
	_ = v146
	var v150 int32
	_ = v150
	var v151 int32
	_ = v151
	var v153 int32
	_ = v153
	var v157 int32
	_ = v157
	var v158 int32
	_ = v158
	var v160 int32
	_ = v160
	var v164 int32
	_ = v164
	var v165 int32
	_ = v165
	var v167 int32
	_ = v167
	var v172 int32
	_ = v172
	var v174 int32
	_ = v174
	var v176 int32
	_ = v176
	var v183 int32
	_ = v183
	var v192 int32
	_ = v192
	var v193 int32
	_ = v193
	var v195 int32
	_ = v195
	var v196 base.V128
	_ = v196
	var v197 int32
	_ = v197
	var v201 base.V128
	_ = v201
	var v202 int32
	_ = v202
	var v204 base.V128
	_ = v204
	var v206 base.V128
	_ = v206
	var v207 base.V128
	_ = v207
	var v212 base.V128
	_ = v212
	var v216 int32
	_ = v216
	var v222 int32
	_ = v222
	var v232 int32
	_ = v232
	var v245 int32
	_ = v245
	var v252 int32
	_ = v252
	var v253 int32
	_ = v253
	var v256 int32
	_ = v256
	var v257 int32
	_ = v257
	var v260 int32
	_ = v260
	var v261 base.V128
	_ = v261
	var v263 base.V128
	_ = v263
	var v264 base.V128
	_ = v264
	var v267 int32
	_ = v267
	var v274 int32
	_ = v274
	var v280 int32
	_ = v280
	var v296 int32
	_ = v296
	var v298 int32
	_ = v298
	var v299 int32
	_ = v299
	var v303 int32
	_ = v303
	var v307 int32
	_ = v307
	var v308 int32
	_ = v308
	var v309 int32
	_ = v309
	var v310 int32
	_ = v310
	var v320 int32
	_ = v320
	var v321 int32
	_ = v321
	var v322 int32
	_ = v322
	var v323 int32
	_ = v323
	var v324 int32
	_ = v324
	var v325 int32
	_ = v325
	var v327 int32
	_ = v327
	var v331 int32
	_ = v331
	var v334 int32
	_ = v334
	var v335 int32
	_ = v335
	var v337 int32
	_ = v337
	var v344 int32
	_ = v344
	if l0 != 0 {
		v174 = int32(0)
		v176 = l3 & int32(-32)
		if v176 < int32(1) {
			v222 = v174
		} else {
			v183 = v174
			for {
				v192 = l0 + v183
				v193 = int32(16)
				v195 = int32(0)
				v196 = base.Simd_g_v128_load_rng(m, v192+v193, v195, int32(-16), int32(32))
				v197 = l1 + v183
				v201 = base.Simd_g_v128_load_rng(m, v197+v193, v195, int32(-16), int32(32))
				v202 = l2 + v183
				v204 = base.Simd_g_v128_load_nc(m, v192, v195)
				v206 = base.Simd_g_v128_load_nc(m, v197, v195)
				v207 = base.Simd_g_i8x16_add(v204, v206)
				base.Simd_g_v128_store(m, v202, v195, v207)
				v212 = base.Simd_g_i8x16_add(v196, v201)
				base.Simd_g_v128_store(m, v202+v193, v195, v212)
				v216 = v183 + int32(32)
				if v216 < v176 {
					v183 = v216
					continue
				} else {
					break
				}
				break
			}
			v222 = v216
		}
		if l3 <= v222 {
		} else {
			v232 = l3 - v222
			if base.Ui32(v232) <= base.Ui32(int32(15)) {
				v280 = v222
				if (l3-v280)&int32(1) == int32(0) {
					v303 = v280
				} else {
					v296 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v280))))
					v298 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v280))))
					v299 = v296 + v298
					*(*uint8)(unsafe.Add(mBase, uint32(l2+v280))) = uint8(v299)
					v303 = v280 + int32(1)
				}
				if v280 == l3+int32(-1) {
				} else {
					v307 = l0
					v308 = l1
					v309 = l2
					v310 = l3
					for {
						v320 = v309 + v303
						v321 = v307 + v303
						v322 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321))))
						v323 = v308 + v303
						v324 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v323))))
						v325 = v322 + v324
						*(*uint8)(unsafe.Add(mBase, uint32(v320))) = uint8(v325)
						v327 = int32(1)
						v331 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321+v327))))
						v334 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v323+v327))))
						v335 = v331 + v334
						*(*uint8)(unsafe.Add(mBase, uint32(v320+v327))) = uint8(v335)
						v337 = int32(2)
						v344 = v310 + int32(-2)
						if v303 != v344 {
							v307 = v307 + v337
							v308 = v308 + v337
							v309 = v309 + v337
							v310 = v344
							continue
						} else {
							break
						}
						break
					}
				}
			} else {
				if base.Ui32(l2-l1) < base.Ui32(int32(16)) {
					v280 = v222
					if (l3-v280)&int32(1) == int32(0) {
						v303 = v280
					} else {
						v296 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v280))))
						v298 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v280))))
						v299 = v296 + v298
						*(*uint8)(unsafe.Add(mBase, uint32(l2+v280))) = uint8(v299)
						v303 = v280 + int32(1)
					}
					if v280 == l3+int32(-1) {
					} else {
						v307 = l0
						v308 = l1
						v309 = l2
						v310 = l3
						for {
							v320 = v309 + v303
							v321 = v307 + v303
							v322 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321))))
							v323 = v308 + v303
							v324 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v323))))
							v325 = v322 + v324
							*(*uint8)(unsafe.Add(mBase, uint32(v320))) = uint8(v325)
							v327 = int32(1)
							v331 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321+v327))))
							v334 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v323+v327))))
							v335 = v331 + v334
							*(*uint8)(unsafe.Add(mBase, uint32(v320+v327))) = uint8(v335)
							v337 = int32(2)
							v344 = v310 + int32(-2)
							if v303 != v344 {
								v307 = v307 + v337
								v308 = v308 + v337
								v309 = v309 + v337
								v310 = v344
								continue
							} else {
								break
							}
							break
						}
					}
				} else {
					if base.Ui32(l2-l0) < base.Ui32(int32(16)) {
						v280 = v222
						if (l3-v280)&int32(1) == int32(0) {
							v303 = v280
						} else {
							v296 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v280))))
							v298 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v280))))
							v299 = v296 + v298
							*(*uint8)(unsafe.Add(mBase, uint32(l2+v280))) = uint8(v299)
							v303 = v280 + int32(1)
						}
						if v280 == l3+int32(-1) {
						} else {
							v307 = l0
							v308 = l1
							v309 = l2
							v310 = l3
							for {
								v320 = v309 + v303
								v321 = v307 + v303
								v322 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321))))
								v323 = v308 + v303
								v324 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v323))))
								v325 = v322 + v324
								*(*uint8)(unsafe.Add(mBase, uint32(v320))) = uint8(v325)
								v327 = int32(1)
								v331 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321+v327))))
								v334 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v323+v327))))
								v335 = v331 + v334
								*(*uint8)(unsafe.Add(mBase, uint32(v320+v327))) = uint8(v335)
								v337 = int32(2)
								v344 = v310 + int32(-2)
								if v303 != v344 {
									v307 = v307 + v337
									v308 = v308 + v337
									v309 = v309 + v337
									v310 = v344
									continue
								} else {
									break
								}
								break
							}
						}
					} else {
						v245 = v232 & int32(-16)
						v252 = l0 + v222
						v253 = l2 + v222
						v256 = l1 + v222
						v257 = v245
						for {
							v260 = int32(0)
							v261 = base.Simd_g_v128_load(m, v252, v260)
							v263 = base.Simd_g_v128_load(m, v256, v260)
							v264 = base.Simd_g_i8x16_add(v261, v263)
							base.Simd_g_v128_store(m, v253, v260, v264)
							v267 = int32(16)
							v274 = v257 + int32(-16)
							if v274 != 0 {
								v252 = v252 + v267
								v253 = v253 + v267
								v256 = v256 + v267
								v257 = v274
								continue
							} else {
								break
							}
							break
						}
						if v232 == v245 {
						} else {
							v280 = v222 + v245
							if (l3-v280)&int32(1) == int32(0) {
								v303 = v280
							} else {
								v296 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v280))))
								v298 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v280))))
								v299 = v296 + v298
								*(*uint8)(unsafe.Add(mBase, uint32(l2+v280))) = uint8(v299)
								v303 = v280 + int32(1)
							}
							if v280 == l3+int32(-1) {
							} else {
								v307 = l0
								v308 = l1
								v309 = l2
								v310 = l3
								for {
									v320 = v309 + v303
									v321 = v307 + v303
									v322 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321))))
									v323 = v308 + v303
									v324 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v323))))
									v325 = v322 + v324
									*(*uint8)(unsafe.Add(mBase, uint32(v320))) = uint8(v325)
									v327 = int32(1)
									v331 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321+v327))))
									v334 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v323+v327))))
									v335 = v331 + v334
									*(*uint8)(unsafe.Add(mBase, uint32(v320+v327))) = uint8(v335)
									v337 = int32(2)
									v344 = v310 + int32(-2)
									if v303 != v344 {
										v307 = v307 + v337
										v308 = v308 + v337
										v309 = v309 + v337
										v310 = v344
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
		}
	} else {
		v14 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
		*(*uint8)(unsafe.Add(mBase, uint32(l2))) = uint8(v14)
		if l3 < int32(2) {
		} else {
			if base.Ui32(l3) < base.Ui32(int32(9)) {
				v74 = int32(1)
			} else {
				v21 = int32(1)
				v26 = int32(0)
				v33 = v26
				v36 = base.Simd_g_i32x4_replace_lane_l0(base.Simd_g_const(&F_VerticalUnfilter_SSE2__k0), v14)
				for {
					v43 = base.Simd_g_const(&F_VerticalUnfilter_SSE2__k0)
					v45 = int32(0)
					v46 = base.Simd_g_v128_load64_zero(m, l1+v21+v33, v45)
					v47 = base.Simd_g_i8x16_add(v46, v36)
					v50 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(v43, v47, base.Simd_g_const(&F_VerticalUnfilter_SSE2__k1), base.Simd_g_const(&F_VerticalUnfilter_SSE2__k2)), v47)
					v53 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(v43, v50, base.Simd_g_const(&F_VerticalUnfilter_SSE2__k3), base.Simd_g_const(&F_VerticalUnfilter_SSE2__k4)), v50)
					v56 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(v43, v53, base.Simd_g_const(&F_VerticalUnfilter_SSE2__k5), base.Simd_g_const(&F_VerticalUnfilter_SSE2__k6)), v53)
					base.Simd_g_v128_store64_lane_l0(m, l2+v21+v33, v45, v56)
					if v33+int32(17) <= l3 {
						v33 = v33 + int32(8)
						v36 = base.Simd_g_i64x2_shr_u(v56, int32(56))
						continue
					} else {
						break
					}
					break
				}
				v74 = v33 + int32(9)
			}
			if l3 <= v74 {
			} else {
				v86 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v74+int32(-1)))))
				v89 = (l3 - v74) & int32(3)
				if v89 != 0 {
					v90 = v86
					v94 = v74
					v96 = v89
					for {
						v105 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v94))))
						v106 = v90 + v105
						*(*uint8)(unsafe.Add(mBase, uint32(l2+v94))) = uint8(v106)
						v109 = v94 + int32(1)
						v111 = v96 + int32(-1)
						if v111 != 0 {
							v90 = v106
							v94 = v109
							v96 = v111
							continue
						} else {
							break
						}
						break
					}
					v112 = v106
					v116 = v109
				} else {
					v112 = v86
					v116 = v74
				}
				if base.Ui32(int32(-4)) < base.Ui32(v74-l3) {
				} else {
					v128 = v112
					v129 = l1
					v130 = l2
					v131 = l3
					for {
						v141 = v130 + v116
						v142 = v129 + v116
						v143 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v142))))
						v144 = v128 + v143
						*(*uint8)(unsafe.Add(mBase, uint32(v141))) = uint8(v144)
						v146 = int32(1)
						v150 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v142+v146))))
						v151 = v144 + v150
						*(*uint8)(unsafe.Add(mBase, uint32(v141+v146))) = uint8(v151)
						v153 = int32(2)
						v157 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v142+v153))))
						v158 = v151 + v157
						*(*uint8)(unsafe.Add(mBase, uint32(v141+v153))) = uint8(v158)
						v160 = int32(3)
						v164 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v142+v160))))
						v165 = v158 + v164
						*(*uint8)(unsafe.Add(mBase, uint32(v141+v160))) = uint8(v165)
						v167 = int32(4)
						v172 = v131 + int32(-4)
						if v116 != v172 {
							v128 = v165
							v129 = v129 + v167
							v130 = v130 + v167
							v131 = v172
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

var F_VerticalUnfilter_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_VerticalUnfilter_SSE2__k1 = [2]uint64{0x808080808080800f, 0x8080808080808080}
var F_VerticalUnfilter_SSE2__k2 = [2]uint64{0x605040302010080, 0xe0d0c0b0a090807}
var F_VerticalUnfilter_SSE2__k3 = [2]uint64{0x8080808080800f0e, 0x8080808080808080}
var F_VerticalUnfilter_SSE2__k4 = [2]uint64{0x504030201008080, 0xd0c0b0a09080706}
var F_VerticalUnfilter_SSE2__k5 = [2]uint64{0x808080800f0e0d0c, 0x8080808080808080}
var F_VerticalUnfilter_SSE2__k6 = [2]uint64{0x302010080808080, 0xb0a090807060504}
