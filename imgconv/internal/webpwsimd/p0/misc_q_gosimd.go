//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_Quantize2Blocks_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	var v4 base.V128
	_ = v4
	var v21 base.V128
	_ = v21
	var v23 base.V128
	_ = v23
	var v24 int32
	_ = v24
	var v25 base.V128
	_ = v25
	var v27 base.V128
	_ = v27
	var v28 int32
	_ = v28
	var v29 base.V128
	_ = v29
	var v30 base.V128
	_ = v30
	var v33 base.V128
	_ = v33
	var v34 base.V128
	_ = v34
	var v35 base.V128
	_ = v35
	var v38 int32
	_ = v38
	var v41 base.V128
	_ = v41
	var v42 base.V128
	_ = v42
	var v48 base.V128
	_ = v48
	var v50 int32
	_ = v50
	var v51 base.V128
	_ = v51
	var v53 base.V128
	_ = v53
	var v55 base.V128
	_ = v55
	var v56 base.V128
	_ = v56
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v63 int32
	_ = v63
	var v64 base.V128
	_ = v64
	var v66 base.V128
	_ = v66
	var v67 int32
	_ = v67
	var v68 base.V128
	_ = v68
	var v69 base.V128
	_ = v69
	var v73 base.V128
	_ = v73
	var v80 base.V128
	_ = v80
	var v89 base.V128
	_ = v89
	var v91 base.V128
	_ = v91
	var v93 base.V128
	_ = v93
	var v94 base.V128
	_ = v94
	var v98 base.V128
	_ = v98
	var v100 base.V128
	_ = v100
	var v102 base.V128
	_ = v102
	var v103 base.V128
	_ = v103
	var v106 base.V128
	_ = v106
	var v110 base.V128
	_ = v110
	var v111 base.V128
	_ = v111
	var v115 int32
	_ = v115
	var v118 int32
	_ = v118
	var v121 base.V128
	_ = v121
	var v123 base.V128
	_ = v123
	var v124 base.V128
	_ = v124
	var v128 base.V128
	_ = v128
	var v142 base.V128
	_ = v142
	var v144 base.V128
	_ = v144
	var v145 base.V128
	_ = v145
	var v149 base.V128
	_ = v149
	var v151 base.V128
	_ = v151
	var v152 base.V128
	_ = v152
	var v156 base.V128
	_ = v156
	var v170 base.V128
	_ = v170
	var v172 base.V128
	_ = v172
	var v173 base.V128
	_ = v173
	var v181 base.V128
	_ = v181
	var v189 base.V128
	_ = v189
	var v201 int32
	_ = v201
	v4 = base.Simd_g_const(&F_Quantize2Blocks_SSE2__k0)
	v21 = base.Simd_g_v128_load_rng(m, l2, int32(96), int32(16), int32(208))
	v23 = base.Simd_g_v128_load_nc(m, l2, int32(208))
	v24 = int32(16)
	v25 = base.Simd_g_v128_load(m, l0, v24)
	v27 = base.Simd_g_i16x8_add(v23, base.Simd_g_i16x8_abs(v25))
	v28 = int32(48)
	v29 = base.Simd_g_v128_load_nc(m, l2, v28)
	v30 = base.Simd_g_i16x8_mul(v27, v29)
	v33 = base.Simd_g_const(&F_Quantize2Blocks_SSE2__k1)
	v34 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_extmul_low_i16x8_u(v27, v29), base.Simd_g_i32x4_extmul_high_i16x8_u(v27, v29), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k2), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k3))
	v35 = base.Simd_g_const(&F_Quantize2Blocks_SSE2__k4)
	v38 = int32(17)
	v41 = base.Simd_g_v128_load_nc(m, l2, int32(112))
	v42 = base.Simd_g_const(&F_Quantize2Blocks_SSE2__k5)
	v48 = base.Simd_g_const(&F_Quantize2Blocks_SSE2__k6)
	v50 = int32(15)
	v51 = base.Simd_g_i16x8_shr_s(v25, v50)
	v53 = base.Simd_g_i16x8_sub(base.Simd_g_v128_xor(base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v21, base.Simd_g_i8x16_shuffle2(v30, v34, base.Simd_g_const(&F_Quantize2Blocks_SSE2__k7), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k8))), v38), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v41, base.Simd_g_i8x16_shuffle2(v30, v34, base.Simd_g_const(&F_Quantize2Blocks_SSE2__k9), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k10))), v38)), v48), v51), v51)
	v55 = base.Simd_g_v128_load_nc(m, l2, v24)
	v56 = base.Simd_g_i16x8_mul(v53, v55)
	base.Simd_g_v128_store(m, l0, v24, v56)
	v60 = base.Simd_g_v128_load_rng(m, l2, int32(64), int32(0), int32(208))
	v62 = base.Simd_g_v128_load_nc(m, l2, int32(192))
	v63 = int32(0)
	v64 = base.Simd_g_v128_load(m, l0, v63)
	v66 = base.Simd_g_i16x8_add(v62, base.Simd_g_i16x8_abs(v64))
	v67 = int32(32)
	v68 = base.Simd_g_v128_load_nc(m, l2, v67)
	v69 = base.Simd_g_i16x8_mul(v66, v68)
	v73 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_extmul_low_i16x8_u(v66, v68), base.Simd_g_i32x4_extmul_high_i16x8_u(v66, v68), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k2), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k3))
	v80 = base.Simd_g_v128_load_nc(m, l2, int32(80))
	v89 = base.Simd_g_i16x8_shr_s(v64, v50)
	v91 = base.Simd_g_i16x8_sub(base.Simd_g_v128_xor(base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v60, base.Simd_g_i8x16_shuffle2(v69, v73, base.Simd_g_const(&F_Quantize2Blocks_SSE2__k7), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k8))), v38), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v80, base.Simd_g_i8x16_shuffle2(v69, v73, base.Simd_g_const(&F_Quantize2Blocks_SSE2__k9), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k10))), v38)), v48), v89), v89)
	v93 = base.Simd_g_v128_load_nc(m, l2, v63)
	v94 = base.Simd_g_i16x8_mul(v91, v93)
	base.Simd_g_v128_store(m, l0, v63, v94)
	v98 = base.Simd_g_const(&F_Quantize2Blocks_SSE2__k11)
	v100 = base.Simd_g_const(&F_Quantize2Blocks_SSE2__k12)
	v102 = base.Simd_g_const(&F_Quantize2Blocks_SSE2__k13)
	v103 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(v91, base.Simd_g_const(&F_Quantize2Blocks_SSE2__k11)), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k12)), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k13))
	base.Simd_g_v128_store(m, l1, v63, v103)
	v106 = base.Simd_g_const(&F_Quantize2Blocks_SSE2__k14)
	v110 = base.Simd_g_const(&F_Quantize2Blocks_SSE2__k15)
	v111 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(v53, base.Simd_g_const(&F_Quantize2Blocks_SSE2__k14)), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k12)), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k15))
	base.Simd_g_v128_store(m, l1, v24, v111)
	v115 = int32(4)
	base.Simd_g_v128_store16_lane_l4(m, l1, int32(6), v111)
	v118 = int32(3)
	base.Simd_g_v128_store16_lane_l3(m, l1, int32(24), v103)
	v121 = base.Simd_g_v128_load(m, l0, v67)
	v123 = base.Simd_g_i16x8_add(v62, base.Simd_g_i16x8_abs(v121))
	v124 = base.Simd_g_i16x8_mul(v68, v123)
	v128 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_extmul_low_i16x8_u(v123, v68), base.Simd_g_i32x4_extmul_high_i16x8_u(v123, v68), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k2), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k3))
	v142 = base.Simd_g_i16x8_shr_s(v121, v50)
	v144 = base.Simd_g_i16x8_sub(base.Simd_g_v128_xor(base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v60, base.Simd_g_i8x16_shuffle2(v124, v128, base.Simd_g_const(&F_Quantize2Blocks_SSE2__k7), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k8))), v38), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v80, base.Simd_g_i8x16_shuffle2(v124, v128, base.Simd_g_const(&F_Quantize2Blocks_SSE2__k9), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k10))), v38)), v48), v142), v142)
	v145 = base.Simd_g_i16x8_mul(v93, v144)
	base.Simd_g_v128_store(m, l0, v67, v145)
	v149 = base.Simd_g_v128_load(m, l0, v28)
	v151 = base.Simd_g_i16x8_add(v23, base.Simd_g_i16x8_abs(v149))
	v152 = base.Simd_g_i16x8_mul(v29, v151)
	v156 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_extmul_low_i16x8_u(v151, v29), base.Simd_g_i32x4_extmul_high_i16x8_u(v151, v29), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k2), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k3))
	v170 = base.Simd_g_i16x8_shr_s(v149, v50)
	v172 = base.Simd_g_i16x8_sub(base.Simd_g_v128_xor(base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v21, base.Simd_g_i8x16_shuffle2(v152, v156, base.Simd_g_const(&F_Quantize2Blocks_SSE2__k7), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k8))), v38), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v41, base.Simd_g_i8x16_shuffle2(v152, v156, base.Simd_g_const(&F_Quantize2Blocks_SSE2__k9), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k10))), v38)), v48), v170), v170)
	v173 = base.Simd_g_i16x8_mul(v55, v172)
	base.Simd_g_v128_store(m, l0, v28, v173)
	v181 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(v172, base.Simd_g_const(&F_Quantize2Blocks_SSE2__k14)), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k12)), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k15))
	base.Simd_g_v128_store(m, l1, v28, v181)
	v189 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(v144, base.Simd_g_const(&F_Quantize2Blocks_SSE2__k11)), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k12)), base.Simd_g_const(&F_Quantize2Blocks_SSE2__k13))
	base.Simd_g_v128_store(m, l1, v67, v189)
	base.Simd_g_v128_store16_lane_l3(m, l1, int32(56), v189)
	base.Simd_g_v128_store16_lane_l4(m, l1, int32(38), v181)
	v201 = int32(_a_F_Quantize2Blocks_SSE2_0)
	return base.B2i32(base.Simd_g_i8x16_bitmask(base.Simd_g_i8x16_eq(base.Simd_g_i8x16_narrow_i16x8_s(v189, v181), v4)) != v201)<<(uint(int32(1))%32) | base.B2i32(base.Simd_g_i8x16_bitmask(base.Simd_g_i8x16_eq(base.Simd_g_i8x16_narrow_i16x8_s(v103, v111), v4)) != v201)
}

var F_Quantize2Blocks_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_Quantize2Blocks_SSE2__k1 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_Quantize2Blocks_SSE2__k2 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_Quantize2Blocks_SSE2__k3 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_Quantize2Blocks_SSE2__k4 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_Quantize2Blocks_SSE2__k5 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_Quantize2Blocks_SSE2__k6 = [2]uint64{0x7ff07ff07ff07ff, 0x7ff07ff07ff07ff}
var F_Quantize2Blocks_SSE2__k7 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_Quantize2Blocks_SSE2__k8 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_Quantize2Blocks_SSE2__k9 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_Quantize2Blocks_SSE2__k10 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_Quantize2Blocks_SSE2__k11 = [2]uint64{0x706050403020100, 0xd0c0b0a0f0e0908}
var F_Quantize2Blocks_SSE2__k12 = [2]uint64{0xb0a090803020100, 0xf0e0d0c07060504}
var F_Quantize2Blocks_SSE2__k13 = [2]uint64{0x706050403020100, 0xf0e0b0a09080d0c}
var F_Quantize2Blocks_SSE2__k14 = [2]uint64{0x706010005040302, 0xf0e0d0c0b0a0908}
var F_Quantize2Blocks_SSE2__k15 = [2]uint64{0x302070605040100, 0xf0e0d0c0b0a0908}

func F_Quantize2Blocks_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	var v4 base.V128
	_ = v4
	var v26 base.V128
	_ = v26
	var v28 base.V128
	_ = v28
	var v29 int32
	_ = v29
	var v30 base.V128
	_ = v30
	var v32 base.V128
	_ = v32
	var v33 int32
	_ = v33
	var v34 base.V128
	_ = v34
	var v35 base.V128
	_ = v35
	var v38 base.V128
	_ = v38
	var v39 base.V128
	_ = v39
	var v40 base.V128
	_ = v40
	var v43 int32
	_ = v43
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v53 base.V128
	_ = v53
	var v57 int32
	_ = v57
	var v58 base.V128
	_ = v58
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v63 base.V128
	_ = v63
	var v67 base.V128
	_ = v67
	var v69 base.V128
	_ = v69
	var v70 int32
	_ = v70
	var v71 base.V128
	_ = v71
	var v73 base.V128
	_ = v73
	var v74 int32
	_ = v74
	var v75 base.V128
	_ = v75
	var v76 base.V128
	_ = v76
	var v80 base.V128
	_ = v80
	var v87 base.V128
	_ = v87
	var v98 base.V128
	_ = v98
	var v100 base.V128
	_ = v100
	var v102 base.V128
	_ = v102
	var v103 base.V128
	_ = v103
	var v106 base.V128
	_ = v106
	var v108 base.V128
	_ = v108
	var v110 base.V128
	_ = v110
	var v113 base.V128
	_ = v113
	var v115 base.V128
	_ = v115
	var v117 base.V128
	_ = v117
	var v121 base.V128
	_ = v121
	var v123 base.V128
	_ = v123
	var v124 base.V128
	_ = v124
	var v128 base.V128
	_ = v128
	var v144 base.V128
	_ = v144
	var v146 base.V128
	_ = v146
	var v147 base.V128
	_ = v147
	var v151 base.V128
	_ = v151
	var v153 base.V128
	_ = v153
	var v154 base.V128
	_ = v154
	var v158 base.V128
	_ = v158
	var v174 base.V128
	_ = v174
	var v176 base.V128
	_ = v176
	var v177 base.V128
	_ = v177
	var v182 base.V128
	_ = v182
	var v187 base.V128
	_ = v187
	var v194 int32
	_ = v194
	v4 = base.Simd_g_const(&F_Quantize2Blocks_SSE41__k0)
	v26 = base.Simd_g_v128_load_rng(m, l2, int32(96), int32(16), int32(208))
	v28 = base.Simd_g_v128_load_nc(m, l2, int32(208))
	v29 = int32(16)
	v30 = base.Simd_g_v128_load(m, l0, v29)
	v32 = base.Simd_g_i16x8_add(v28, base.Simd_g_i16x8_abs(v30))
	v33 = int32(48)
	v34 = base.Simd_g_v128_load_nc(m, l2, v33)
	v35 = base.Simd_g_i16x8_mul(v32, v34)
	v38 = base.Simd_g_const(&F_Quantize2Blocks_SSE41__k1)
	v39 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_extmul_low_i16x8_u(v32, v34), base.Simd_g_i32x4_extmul_high_i16x8_u(v32, v34), base.Simd_g_const(&F_Quantize2Blocks_SSE41__k2), base.Simd_g_const(&F_Quantize2Blocks_SSE41__k3))
	v40 = base.Simd_g_const(&F_Quantize2Blocks_SSE41__k4)
	v43 = int32(17)
	v46 = base.Simd_g_v128_load_nc(m, l2, int32(112))
	v47 = base.Simd_g_const(&F_Quantize2Blocks_SSE41__k5)
	v53 = base.Simd_g_const(&F_Quantize2Blocks_SSE41__k6)
	v57 = int32(15)
	v58 = base.Simd_g_i16x8_shr_s(v30, v57)
	v60 = base.Simd_g_v128_xor(base.Simd_g_i16x8_add(base.Simd_g_v128_bitselect(v4, base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v26, base.Simd_g_i8x16_shuffle2(v35, v39, base.Simd_g_const(&F_Quantize2Blocks_SSE41__k7), base.Simd_g_const(&F_Quantize2Blocks_SSE41__k8))), v43), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v46, base.Simd_g_i8x16_shuffle2(v35, v39, base.Simd_g_const(&F_Quantize2Blocks_SSE41__k9), base.Simd_g_const(&F_Quantize2Blocks_SSE41__k10))), v43)), v53), base.Simd_g_i16x8_eq(v30, v4)), v58), v58)
	v62 = base.Simd_g_v128_load_nc(m, l2, v29)
	v63 = base.Simd_g_i16x8_mul(v60, v62)
	base.Simd_g_v128_store(m, l0, v29, v63)
	v67 = base.Simd_g_v128_load_rng(m, l2, int32(64), int32(0), int32(208))
	v69 = base.Simd_g_v128_load_nc(m, l2, int32(192))
	v70 = int32(0)
	v71 = base.Simd_g_v128_load(m, l0, v70)
	v73 = base.Simd_g_i16x8_add(v69, base.Simd_g_i16x8_abs(v71))
	v74 = int32(32)
	v75 = base.Simd_g_v128_load_nc(m, l2, v74)
	v76 = base.Simd_g_i16x8_mul(v73, v75)
	v80 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_extmul_low_i16x8_u(v73, v75), base.Simd_g_i32x4_extmul_high_i16x8_u(v73, v75), base.Simd_g_const(&F_Quantize2Blocks_SSE41__k2), base.Simd_g_const(&F_Quantize2Blocks_SSE41__k3))
	v87 = base.Simd_g_v128_load_nc(m, l2, int32(80))
	v98 = base.Simd_g_i16x8_shr_s(v71, v57)
	v100 = base.Simd_g_v128_xor(base.Simd_g_i16x8_add(base.Simd_g_v128_bitselect(v4, base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v67, base.Simd_g_i8x16_shuffle2(v76, v80, base.Simd_g_const(&F_Quantize2Blocks_SSE41__k7), base.Simd_g_const(&F_Quantize2Blocks_SSE41__k8))), v43), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v87, base.Simd_g_i8x16_shuffle2(v76, v80, base.Simd_g_const(&F_Quantize2Blocks_SSE41__k9), base.Simd_g_const(&F_Quantize2Blocks_SSE41__k10))), v43)), v53), base.Simd_g_i16x8_eq(v71, v4)), v98), v98)
	v102 = base.Simd_g_v128_load_nc(m, l2, v70)
	v103 = base.Simd_g_i16x8_mul(v100, v102)
	base.Simd_g_v128_store(m, l0, v70, v103)
	v106 = base.Simd_g_const(&F_Quantize2Blocks_SSE41__k11)
	v108 = base.Simd_g_const(&F_Quantize2Blocks_SSE41__k12)
	v110 = base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v60, v106), base.Simd_g_i8x16_swizzle(v100, v108))
	base.Simd_g_v128_store(m, l1, v70, v110)
	v113 = base.Simd_g_const(&F_Quantize2Blocks_SSE41__k13)
	v115 = base.Simd_g_const(&F_Quantize2Blocks_SSE41__k14)
	v117 = base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v60, v113), base.Simd_g_i8x16_swizzle(v100, v115))
	base.Simd_g_v128_store(m, l1, v29, v117)
	v121 = base.Simd_g_v128_load(m, l0, v74)
	v123 = base.Simd_g_i16x8_add(v69, base.Simd_g_i16x8_abs(v121))
	v124 = base.Simd_g_i16x8_mul(v75, v123)
	v128 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_extmul_low_i16x8_u(v123, v75), base.Simd_g_i32x4_extmul_high_i16x8_u(v123, v75), base.Simd_g_const(&F_Quantize2Blocks_SSE41__k2), base.Simd_g_const(&F_Quantize2Blocks_SSE41__k3))
	v144 = base.Simd_g_i16x8_shr_s(v121, v57)
	v146 = base.Simd_g_v128_xor(base.Simd_g_i16x8_add(base.Simd_g_v128_bitselect(v4, base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v67, base.Simd_g_i8x16_shuffle2(v124, v128, base.Simd_g_const(&F_Quantize2Blocks_SSE41__k7), base.Simd_g_const(&F_Quantize2Blocks_SSE41__k8))), v43), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v87, base.Simd_g_i8x16_shuffle2(v124, v128, base.Simd_g_const(&F_Quantize2Blocks_SSE41__k9), base.Simd_g_const(&F_Quantize2Blocks_SSE41__k10))), v43)), v53), base.Simd_g_i16x8_eq(v121, v4)), v144), v144)
	v147 = base.Simd_g_i16x8_mul(v102, v146)
	base.Simd_g_v128_store(m, l0, v74, v147)
	v151 = base.Simd_g_v128_load(m, l0, v33)
	v153 = base.Simd_g_i16x8_add(v28, base.Simd_g_i16x8_abs(v151))
	v154 = base.Simd_g_i16x8_mul(v34, v153)
	v158 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_extmul_low_i16x8_u(v153, v34), base.Simd_g_i32x4_extmul_high_i16x8_u(v153, v34), base.Simd_g_const(&F_Quantize2Blocks_SSE41__k2), base.Simd_g_const(&F_Quantize2Blocks_SSE41__k3))
	v174 = base.Simd_g_i16x8_shr_s(v151, v57)
	v176 = base.Simd_g_v128_xor(base.Simd_g_i16x8_add(base.Simd_g_v128_bitselect(v4, base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v26, base.Simd_g_i8x16_shuffle2(v154, v158, base.Simd_g_const(&F_Quantize2Blocks_SSE41__k7), base.Simd_g_const(&F_Quantize2Blocks_SSE41__k8))), v43), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v46, base.Simd_g_i8x16_shuffle2(v154, v158, base.Simd_g_const(&F_Quantize2Blocks_SSE41__k9), base.Simd_g_const(&F_Quantize2Blocks_SSE41__k10))), v43)), v53), base.Simd_g_i16x8_eq(v151, v4)), v174), v174)
	v177 = base.Simd_g_i16x8_mul(v62, v176)
	base.Simd_g_v128_store(m, l0, v33, v177)
	v182 = base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v176, v113), base.Simd_g_i8x16_swizzle(v146, v115))
	base.Simd_g_v128_store(m, l1, v33, v182)
	v187 = base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v176, v106), base.Simd_g_i8x16_swizzle(v146, v108))
	base.Simd_g_v128_store(m, l1, v74, v187)
	v194 = int32(_a_F_Quantize2Blocks_SSE41_0)
	return base.B2i32(base.Simd_g_i8x16_bitmask(base.Simd_g_i8x16_eq(base.Simd_g_i8x16_narrow_i16x8_s(v187, v182), v4)) != v194)<<(uint(int32(1))%32) | base.B2i32(base.Simd_g_i8x16_bitmask(base.Simd_g_i8x16_eq(base.Simd_g_i8x16_narrow_i16x8_s(v110, v117), v4)) != v194)
}

var F_Quantize2Blocks_SSE41__k0 = [2]uint64{0x0, 0x0}
var F_Quantize2Blocks_SSE41__k1 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_Quantize2Blocks_SSE41__k2 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_Quantize2Blocks_SSE41__k3 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_Quantize2Blocks_SSE41__k4 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_Quantize2Blocks_SSE41__k5 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_Quantize2Blocks_SSE41__k6 = [2]uint64{0x7ff07ff07ff07ff, 0x7ff07ff07ff07ff}
var F_Quantize2Blocks_SSE41__k7 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_Quantize2Blocks_SSE41__k8 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_Quantize2Blocks_SSE41__k9 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_Quantize2Blocks_SSE41__k10 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_Quantize2Blocks_SSE41__k11 = [2]uint64{0x1008f8e8f8e8f8e, 0x8f8e8f8e8f8e8f8e}
var F_Quantize2Blocks_SSE41__k12 = [2]uint64{0x8f8e090803020100, 0xd0c070605040b0a}
var F_Quantize2Blocks_SSE41__k13 = [2]uint64{0x5040b0a09080302, 0xf0e0d0c07068f8e}
var F_Quantize2Blocks_SSE41__k14 = [2]uint64{0x8f8e8f8e8f8e8f8e, 0x8f8e8f8e8f8e0f0e}

func F_QuantizeBlockWHT_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	var v4 base.V128
	_ = v4
	var v11 base.V128
	_ = v11
	var v12 int32
	_ = v12
	var v13 base.V128
	_ = v13
	var v14 base.V128
	_ = v14
	var v16 base.V128
	_ = v16
	var v17 base.V128
	_ = v17
	var v20 base.V128
	_ = v20
	var v21 base.V128
	_ = v21
	var v22 base.V128
	_ = v22
	var v25 int32
	_ = v25
	var v28 base.V128
	_ = v28
	var v29 base.V128
	_ = v29
	var v35 base.V128
	_ = v35
	var v37 int32
	_ = v37
	var v38 base.V128
	_ = v38
	var v40 base.V128
	_ = v40
	var v42 base.V128
	_ = v42
	var v43 base.V128
	_ = v43
	var v47 base.V128
	_ = v47
	var v48 int32
	_ = v48
	var v49 base.V128
	_ = v49
	var v50 base.V128
	_ = v50
	var v52 base.V128
	_ = v52
	var v53 base.V128
	_ = v53
	var v57 base.V128
	_ = v57
	var v64 base.V128
	_ = v64
	var v73 base.V128
	_ = v73
	var v75 base.V128
	_ = v75
	var v77 base.V128
	_ = v77
	var v78 base.V128
	_ = v78
	var v84 base.V128
	_ = v84
	var v87 base.V128
	_ = v87
	var v95 base.V128
	_ = v95
	var v99 int32
	_ = v99
	var v102 int32
	_ = v102
	v4 = base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k0)
	v11 = base.Simd_g_v128_load_rng(m, l2, int32(96), int32(16), int32(112))
	v12 = int32(16)
	v13 = base.Simd_g_v128_load(m, l0, v12)
	v14 = base.Simd_g_i16x8_abs(v13)
	v16 = base.Simd_g_v128_load_nc(m, l2, int32(48))
	v17 = base.Simd_g_i16x8_mul(v14, v16)
	v20 = base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k1)
	v21 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_extmul_low_i16x8_u(v14, v16), base.Simd_g_i32x4_extmul_high_i16x8_u(v14, v16), base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k2), base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k3))
	v22 = base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k4)
	v25 = int32(17)
	v28 = base.Simd_g_v128_load_nc(m, l2, int32(112))
	v29 = base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k5)
	v35 = base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k6)
	v37 = int32(15)
	v38 = base.Simd_g_i16x8_shr_s(v13, v37)
	v40 = base.Simd_g_i16x8_sub(base.Simd_g_v128_xor(base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v11, base.Simd_g_i8x16_shuffle2(v17, v21, base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k7), base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k8))), v25), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v28, base.Simd_g_i8x16_shuffle2(v17, v21, base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k9), base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k10))), v25)), v35), v38), v38)
	v42 = base.Simd_g_v128_load_nc(m, l2, v12)
	v43 = base.Simd_g_i16x8_mul(v40, v42)
	base.Simd_g_v128_store(m, l0, v12, v43)
	v47 = base.Simd_g_v128_load_rng(m, l2, int32(64), int32(0), int32(96))
	v48 = int32(0)
	v49 = base.Simd_g_v128_load(m, l0, v48)
	v50 = base.Simd_g_i16x8_abs(v49)
	v52 = base.Simd_g_v128_load_nc(m, l2, int32(32))
	v53 = base.Simd_g_i16x8_mul(v50, v52)
	v57 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_extmul_low_i16x8_u(v50, v52), base.Simd_g_i32x4_extmul_high_i16x8_u(v50, v52), base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k2), base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k3))
	v64 = base.Simd_g_v128_load_nc(m, l2, int32(80))
	v73 = base.Simd_g_i16x8_shr_s(v49, v37)
	v75 = base.Simd_g_i16x8_sub(base.Simd_g_v128_xor(base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v47, base.Simd_g_i8x16_shuffle2(v53, v57, base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k7), base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k8))), v25), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v64, base.Simd_g_i8x16_shuffle2(v53, v57, base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k9), base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k10))), v25)), v35), v73), v73)
	v77 = base.Simd_g_v128_load_nc(m, l2, v48)
	v78 = base.Simd_g_i16x8_mul(v75, v77)
	base.Simd_g_v128_store(m, l0, v48, v78)
	v84 = base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k11)
	v87 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(v75, base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k12)), base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k11)), base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k13))
	base.Simd_g_v128_store(m, l1, v48, v87)
	v95 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(v40, base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k14)), base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k11)), base.Simd_g_const(&F_QuantizeBlockWHT_SSE2__k15))
	base.Simd_g_v128_store(m, l1, v12, v95)
	v99 = int32(4)
	base.Simd_g_v128_store16_lane_l4(m, l1, int32(6), v95)
	v102 = int32(3)
	base.Simd_g_v128_store16_lane_l3(m, l1, int32(24), v87)
	return base.B2i32(base.Simd_g_i8x16_bitmask(base.Simd_g_i8x16_eq(base.Simd_g_i8x16_narrow_i16x8_s(v87, v95), v4)) != int32(_a_F_QuantizeBlockWHT_SSE2_0))
}

var F_QuantizeBlockWHT_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_QuantizeBlockWHT_SSE2__k1 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_QuantizeBlockWHT_SSE2__k2 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_QuantizeBlockWHT_SSE2__k3 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_QuantizeBlockWHT_SSE2__k4 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_QuantizeBlockWHT_SSE2__k5 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_QuantizeBlockWHT_SSE2__k6 = [2]uint64{0x7ff07ff07ff07ff, 0x7ff07ff07ff07ff}
var F_QuantizeBlockWHT_SSE2__k7 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_QuantizeBlockWHT_SSE2__k8 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_QuantizeBlockWHT_SSE2__k9 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_QuantizeBlockWHT_SSE2__k10 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_QuantizeBlockWHT_SSE2__k11 = [2]uint64{0xb0a090803020100, 0xf0e0d0c07060504}
var F_QuantizeBlockWHT_SSE2__k12 = [2]uint64{0x706050403020100, 0xd0c0b0a0f0e0908}
var F_QuantizeBlockWHT_SSE2__k13 = [2]uint64{0x706050403020100, 0xf0e0b0a09080d0c}
var F_QuantizeBlockWHT_SSE2__k14 = [2]uint64{0x706010005040302, 0xf0e0d0c0b0a0908}
var F_QuantizeBlockWHT_SSE2__k15 = [2]uint64{0x302070605040100, 0xf0e0d0c0b0a0908}

func F_QuantizeBlockWHT_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	var v4 base.V128
	_ = v4
	var v13 base.V128
	_ = v13
	var v15 base.V128
	_ = v15
	var v16 int32
	_ = v16
	var v17 base.V128
	_ = v17
	var v18 base.V128
	_ = v18
	var v19 base.V128
	_ = v19
	var v22 base.V128
	_ = v22
	var v23 base.V128
	_ = v23
	var v24 base.V128
	_ = v24
	var v27 int32
	_ = v27
	var v30 base.V128
	_ = v30
	var v31 base.V128
	_ = v31
	var v37 base.V128
	_ = v37
	var v41 int32
	_ = v41
	var v42 base.V128
	_ = v42
	var v44 base.V128
	_ = v44
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v51 base.V128
	_ = v51
	var v53 base.V128
	_ = v53
	var v54 int32
	_ = v54
	var v55 base.V128
	_ = v55
	var v56 base.V128
	_ = v56
	var v57 base.V128
	_ = v57
	var v61 base.V128
	_ = v61
	var v68 base.V128
	_ = v68
	var v79 base.V128
	_ = v79
	var v81 base.V128
	_ = v81
	var v83 base.V128
	_ = v83
	var v84 base.V128
	_ = v84
	var v91 base.V128
	_ = v91
	var v98 base.V128
	_ = v98
	v4 = base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k0)
	v13 = base.Simd_g_v128_load_rng(m, l2, int32(96), int32(16), int32(112))
	v15 = base.Simd_g_v128_load_nc(m, l2, int32(48))
	v16 = int32(16)
	v17 = base.Simd_g_v128_load(m, l0, v16)
	v18 = base.Simd_g_i16x8_abs(v17)
	v19 = base.Simd_g_i16x8_mul(v15, v18)
	v22 = base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k1)
	v23 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_extmul_low_i16x8_u(v15, v18), base.Simd_g_i32x4_extmul_high_i16x8_u(v15, v18), base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k2), base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k3))
	v24 = base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k4)
	v27 = int32(17)
	v30 = base.Simd_g_v128_load_nc(m, l2, int32(112))
	v31 = base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k5)
	v37 = base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k6)
	v41 = int32(15)
	v42 = base.Simd_g_i16x8_shr_s(v17, v41)
	v44 = base.Simd_g_v128_xor(base.Simd_g_i16x8_add(base.Simd_g_v128_bitselect(v4, base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v13, base.Simd_g_i8x16_shuffle2(v19, v23, base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k7), base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k8))), v27), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v30, base.Simd_g_i8x16_shuffle2(v19, v23, base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k9), base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k10))), v27)), v37), base.Simd_g_i16x8_eq(v17, v4)), v42), v42)
	v46 = base.Simd_g_v128_load_nc(m, l2, v16)
	v47 = base.Simd_g_i16x8_mul(v44, v46)
	base.Simd_g_v128_store(m, l0, v16, v47)
	v51 = base.Simd_g_v128_load_rng(m, l2, int32(64), int32(0), int32(96))
	v53 = base.Simd_g_v128_load_nc(m, l2, int32(32))
	v54 = int32(0)
	v55 = base.Simd_g_v128_load(m, l0, v54)
	v56 = base.Simd_g_i16x8_abs(v55)
	v57 = base.Simd_g_i16x8_mul(v53, v56)
	v61 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_extmul_low_i16x8_u(v53, v56), base.Simd_g_i32x4_extmul_high_i16x8_u(v53, v56), base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k2), base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k3))
	v68 = base.Simd_g_v128_load_nc(m, l2, int32(80))
	v79 = base.Simd_g_i16x8_shr_s(v55, v41)
	v81 = base.Simd_g_v128_xor(base.Simd_g_i16x8_add(base.Simd_g_v128_bitselect(v4, base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v51, base.Simd_g_i8x16_shuffle2(v57, v61, base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k7), base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k8))), v27), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v68, base.Simd_g_i8x16_shuffle2(v57, v61, base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k9), base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k10))), v27)), v37), base.Simd_g_i16x8_eq(v55, v4)), v79), v79)
	v83 = base.Simd_g_v128_load_nc(m, l2, v54)
	v84 = base.Simd_g_i16x8_mul(v81, v83)
	base.Simd_g_v128_store(m, l0, v54, v84)
	v91 = base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v44, base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k11)), base.Simd_g_i8x16_swizzle(v81, base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k12)))
	base.Simd_g_v128_store(m, l1, v54, v91)
	v98 = base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v44, base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k13)), base.Simd_g_i8x16_swizzle(v81, base.Simd_g_const(&F_QuantizeBlockWHT_SSE41__k14)))
	base.Simd_g_v128_store(m, l1, v16, v98)
	return base.B2i32(base.Simd_g_i8x16_bitmask(base.Simd_g_i8x16_eq(base.Simd_g_i8x16_narrow_i16x8_s(v91, v98), v4)) != int32(_a_F_QuantizeBlockWHT_SSE41_0))
}

var F_QuantizeBlockWHT_SSE41__k0 = [2]uint64{0x0, 0x0}
var F_QuantizeBlockWHT_SSE41__k1 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_QuantizeBlockWHT_SSE41__k2 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_QuantizeBlockWHT_SSE41__k3 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_QuantizeBlockWHT_SSE41__k4 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_QuantizeBlockWHT_SSE41__k5 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_QuantizeBlockWHT_SSE41__k6 = [2]uint64{0x7ff07ff07ff07ff, 0x7ff07ff07ff07ff}
var F_QuantizeBlockWHT_SSE41__k7 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_QuantizeBlockWHT_SSE41__k8 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_QuantizeBlockWHT_SSE41__k9 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_QuantizeBlockWHT_SSE41__k10 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_QuantizeBlockWHT_SSE41__k11 = [2]uint64{0x1008f8e8f8e8f8e, 0x8f8e8f8e8f8e8f8e}
var F_QuantizeBlockWHT_SSE41__k12 = [2]uint64{0x8f8e090803020100, 0xd0c070605040b0a}
var F_QuantizeBlockWHT_SSE41__k13 = [2]uint64{0x5040b0a09080302, 0xf0e0d0c07068f8e}
var F_QuantizeBlockWHT_SSE41__k14 = [2]uint64{0x8f8e8f8e8f8e8f8e, 0x8f8e8f8e8f8e0f0e}

func F_QuantizeBlock_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	var v4 base.V128
	_ = v4
	var v11 base.V128
	_ = v11
	var v13 base.V128
	_ = v13
	var v14 int32
	_ = v14
	var v15 base.V128
	_ = v15
	var v17 base.V128
	_ = v17
	var v19 base.V128
	_ = v19
	var v20 base.V128
	_ = v20
	var v23 base.V128
	_ = v23
	var v24 base.V128
	_ = v24
	var v25 base.V128
	_ = v25
	var v28 int32
	_ = v28
	var v31 base.V128
	_ = v31
	var v32 base.V128
	_ = v32
	var v38 base.V128
	_ = v38
	var v40 int32
	_ = v40
	var v41 base.V128
	_ = v41
	var v43 base.V128
	_ = v43
	var v45 base.V128
	_ = v45
	var v46 base.V128
	_ = v46
	var v50 base.V128
	_ = v50
	var v52 base.V128
	_ = v52
	var v53 int32
	_ = v53
	var v54 base.V128
	_ = v54
	var v56 base.V128
	_ = v56
	var v58 base.V128
	_ = v58
	var v59 base.V128
	_ = v59
	var v63 base.V128
	_ = v63
	var v70 base.V128
	_ = v70
	var v79 base.V128
	_ = v79
	var v81 base.V128
	_ = v81
	var v83 base.V128
	_ = v83
	var v84 base.V128
	_ = v84
	var v90 base.V128
	_ = v90
	var v93 base.V128
	_ = v93
	var v101 base.V128
	_ = v101
	var v105 int32
	_ = v105
	var v108 int32
	_ = v108
	v4 = base.Simd_g_const(&F_QuantizeBlock_SSE2__k0)
	v11 = base.Simd_g_v128_load_rng(m, l2, int32(96), int32(16), int32(208))
	v13 = base.Simd_g_v128_load_nc(m, l2, int32(208))
	v14 = int32(16)
	v15 = base.Simd_g_v128_load(m, l0, v14)
	v17 = base.Simd_g_i16x8_add(v13, base.Simd_g_i16x8_abs(v15))
	v19 = base.Simd_g_v128_load_nc(m, l2, int32(48))
	v20 = base.Simd_g_i16x8_mul(v17, v19)
	v23 = base.Simd_g_const(&F_QuantizeBlock_SSE2__k1)
	v24 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_extmul_low_i16x8_u(v17, v19), base.Simd_g_i32x4_extmul_high_i16x8_u(v17, v19), base.Simd_g_const(&F_QuantizeBlock_SSE2__k2), base.Simd_g_const(&F_QuantizeBlock_SSE2__k3))
	v25 = base.Simd_g_const(&F_QuantizeBlock_SSE2__k4)
	v28 = int32(17)
	v31 = base.Simd_g_v128_load_nc(m, l2, int32(112))
	v32 = base.Simd_g_const(&F_QuantizeBlock_SSE2__k5)
	v38 = base.Simd_g_const(&F_QuantizeBlock_SSE2__k6)
	v40 = int32(15)
	v41 = base.Simd_g_i16x8_shr_s(v15, v40)
	v43 = base.Simd_g_i16x8_sub(base.Simd_g_v128_xor(base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v11, base.Simd_g_i8x16_shuffle2(v20, v24, base.Simd_g_const(&F_QuantizeBlock_SSE2__k7), base.Simd_g_const(&F_QuantizeBlock_SSE2__k8))), v28), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v31, base.Simd_g_i8x16_shuffle2(v20, v24, base.Simd_g_const(&F_QuantizeBlock_SSE2__k9), base.Simd_g_const(&F_QuantizeBlock_SSE2__k10))), v28)), v38), v41), v41)
	v45 = base.Simd_g_v128_load_nc(m, l2, v14)
	v46 = base.Simd_g_i16x8_mul(v43, v45)
	base.Simd_g_v128_store(m, l0, v14, v46)
	v50 = base.Simd_g_v128_load_rng(m, l2, int32(64), int32(0), int32(208))
	v52 = base.Simd_g_v128_load_nc(m, l2, int32(192))
	v53 = int32(0)
	v54 = base.Simd_g_v128_load(m, l0, v53)
	v56 = base.Simd_g_i16x8_add(v52, base.Simd_g_i16x8_abs(v54))
	v58 = base.Simd_g_v128_load_nc(m, l2, int32(32))
	v59 = base.Simd_g_i16x8_mul(v56, v58)
	v63 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_extmul_low_i16x8_u(v56, v58), base.Simd_g_i32x4_extmul_high_i16x8_u(v56, v58), base.Simd_g_const(&F_QuantizeBlock_SSE2__k2), base.Simd_g_const(&F_QuantizeBlock_SSE2__k3))
	v70 = base.Simd_g_v128_load_nc(m, l2, int32(80))
	v79 = base.Simd_g_i16x8_shr_s(v54, v40)
	v81 = base.Simd_g_i16x8_sub(base.Simd_g_v128_xor(base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v50, base.Simd_g_i8x16_shuffle2(v59, v63, base.Simd_g_const(&F_QuantizeBlock_SSE2__k7), base.Simd_g_const(&F_QuantizeBlock_SSE2__k8))), v28), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v70, base.Simd_g_i8x16_shuffle2(v59, v63, base.Simd_g_const(&F_QuantizeBlock_SSE2__k9), base.Simd_g_const(&F_QuantizeBlock_SSE2__k10))), v28)), v38), v79), v79)
	v83 = base.Simd_g_v128_load_nc(m, l2, v53)
	v84 = base.Simd_g_i16x8_mul(v81, v83)
	base.Simd_g_v128_store(m, l0, v53, v84)
	v90 = base.Simd_g_const(&F_QuantizeBlock_SSE2__k11)
	v93 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(v81, base.Simd_g_const(&F_QuantizeBlock_SSE2__k12)), base.Simd_g_const(&F_QuantizeBlock_SSE2__k11)), base.Simd_g_const(&F_QuantizeBlock_SSE2__k13))
	base.Simd_g_v128_store(m, l1, v53, v93)
	v101 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(v43, base.Simd_g_const(&F_QuantizeBlock_SSE2__k14)), base.Simd_g_const(&F_QuantizeBlock_SSE2__k11)), base.Simd_g_const(&F_QuantizeBlock_SSE2__k15))
	base.Simd_g_v128_store(m, l1, v14, v101)
	v105 = int32(4)
	base.Simd_g_v128_store16_lane_l4(m, l1, int32(6), v101)
	v108 = int32(3)
	base.Simd_g_v128_store16_lane_l3(m, l1, int32(24), v93)
	return base.B2i32(base.Simd_g_i8x16_bitmask(base.Simd_g_i8x16_eq(base.Simd_g_i8x16_narrow_i16x8_s(v93, v101), v4)) != int32(_a_F_QuantizeBlock_SSE2_0))
}

var F_QuantizeBlock_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_QuantizeBlock_SSE2__k1 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_QuantizeBlock_SSE2__k2 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_QuantizeBlock_SSE2__k3 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_QuantizeBlock_SSE2__k4 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_QuantizeBlock_SSE2__k5 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_QuantizeBlock_SSE2__k6 = [2]uint64{0x7ff07ff07ff07ff, 0x7ff07ff07ff07ff}
var F_QuantizeBlock_SSE2__k7 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_QuantizeBlock_SSE2__k8 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_QuantizeBlock_SSE2__k9 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_QuantizeBlock_SSE2__k10 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_QuantizeBlock_SSE2__k11 = [2]uint64{0xb0a090803020100, 0xf0e0d0c07060504}
var F_QuantizeBlock_SSE2__k12 = [2]uint64{0x706050403020100, 0xd0c0b0a0f0e0908}
var F_QuantizeBlock_SSE2__k13 = [2]uint64{0x706050403020100, 0xf0e0b0a09080d0c}
var F_QuantizeBlock_SSE2__k14 = [2]uint64{0x706010005040302, 0xf0e0d0c0b0a0908}
var F_QuantizeBlock_SSE2__k15 = [2]uint64{0x302070605040100, 0xf0e0d0c0b0a0908}

func F_QuantizeBlock_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	var v4 base.V128
	_ = v4
	var v13 base.V128
	_ = v13
	var v15 base.V128
	_ = v15
	var v16 int32
	_ = v16
	var v17 base.V128
	_ = v17
	var v19 base.V128
	_ = v19
	var v21 base.V128
	_ = v21
	var v22 base.V128
	_ = v22
	var v25 base.V128
	_ = v25
	var v26 base.V128
	_ = v26
	var v27 base.V128
	_ = v27
	var v30 int32
	_ = v30
	var v33 base.V128
	_ = v33
	var v34 base.V128
	_ = v34
	var v40 base.V128
	_ = v40
	var v44 int32
	_ = v44
	var v45 base.V128
	_ = v45
	var v47 base.V128
	_ = v47
	var v49 base.V128
	_ = v49
	var v50 base.V128
	_ = v50
	var v54 base.V128
	_ = v54
	var v56 base.V128
	_ = v56
	var v57 int32
	_ = v57
	var v58 base.V128
	_ = v58
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v63 base.V128
	_ = v63
	var v67 base.V128
	_ = v67
	var v74 base.V128
	_ = v74
	var v85 base.V128
	_ = v85
	var v87 base.V128
	_ = v87
	var v89 base.V128
	_ = v89
	var v90 base.V128
	_ = v90
	var v97 base.V128
	_ = v97
	var v104 base.V128
	_ = v104
	v4 = base.Simd_g_const(&F_QuantizeBlock_SSE41__k0)
	v13 = base.Simd_g_v128_load_rng(m, l2, int32(96), int32(16), int32(208))
	v15 = base.Simd_g_v128_load_nc(m, l2, int32(208))
	v16 = int32(16)
	v17 = base.Simd_g_v128_load(m, l0, v16)
	v19 = base.Simd_g_i16x8_add(v15, base.Simd_g_i16x8_abs(v17))
	v21 = base.Simd_g_v128_load_nc(m, l2, int32(48))
	v22 = base.Simd_g_i16x8_mul(v19, v21)
	v25 = base.Simd_g_const(&F_QuantizeBlock_SSE41__k1)
	v26 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_extmul_low_i16x8_u(v19, v21), base.Simd_g_i32x4_extmul_high_i16x8_u(v19, v21), base.Simd_g_const(&F_QuantizeBlock_SSE41__k2), base.Simd_g_const(&F_QuantizeBlock_SSE41__k3))
	v27 = base.Simd_g_const(&F_QuantizeBlock_SSE41__k4)
	v30 = int32(17)
	v33 = base.Simd_g_v128_load_nc(m, l2, int32(112))
	v34 = base.Simd_g_const(&F_QuantizeBlock_SSE41__k5)
	v40 = base.Simd_g_const(&F_QuantizeBlock_SSE41__k6)
	v44 = int32(15)
	v45 = base.Simd_g_i16x8_shr_s(v17, v44)
	v47 = base.Simd_g_v128_xor(base.Simd_g_i16x8_add(base.Simd_g_v128_bitselect(v4, base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v13, base.Simd_g_i8x16_shuffle2(v22, v26, base.Simd_g_const(&F_QuantizeBlock_SSE41__k7), base.Simd_g_const(&F_QuantizeBlock_SSE41__k8))), v30), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v33, base.Simd_g_i8x16_shuffle2(v22, v26, base.Simd_g_const(&F_QuantizeBlock_SSE41__k9), base.Simd_g_const(&F_QuantizeBlock_SSE41__k10))), v30)), v40), base.Simd_g_i16x8_eq(v17, v4)), v45), v45)
	v49 = base.Simd_g_v128_load_nc(m, l2, v16)
	v50 = base.Simd_g_i16x8_mul(v47, v49)
	base.Simd_g_v128_store(m, l0, v16, v50)
	v54 = base.Simd_g_v128_load_rng(m, l2, int32(64), int32(0), int32(208))
	v56 = base.Simd_g_v128_load_nc(m, l2, int32(192))
	v57 = int32(0)
	v58 = base.Simd_g_v128_load(m, l0, v57)
	v60 = base.Simd_g_i16x8_add(v56, base.Simd_g_i16x8_abs(v58))
	v62 = base.Simd_g_v128_load_nc(m, l2, int32(32))
	v63 = base.Simd_g_i16x8_mul(v60, v62)
	v67 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_extmul_low_i16x8_u(v60, v62), base.Simd_g_i32x4_extmul_high_i16x8_u(v60, v62), base.Simd_g_const(&F_QuantizeBlock_SSE41__k2), base.Simd_g_const(&F_QuantizeBlock_SSE41__k3))
	v74 = base.Simd_g_v128_load_nc(m, l2, int32(80))
	v85 = base.Simd_g_i16x8_shr_s(v58, v44)
	v87 = base.Simd_g_v128_xor(base.Simd_g_i16x8_add(base.Simd_g_v128_bitselect(v4, base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v54, base.Simd_g_i8x16_shuffle2(v63, v67, base.Simd_g_const(&F_QuantizeBlock_SSE41__k7), base.Simd_g_const(&F_QuantizeBlock_SSE41__k8))), v30), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v74, base.Simd_g_i8x16_shuffle2(v63, v67, base.Simd_g_const(&F_QuantizeBlock_SSE41__k9), base.Simd_g_const(&F_QuantizeBlock_SSE41__k10))), v30)), v40), base.Simd_g_i16x8_eq(v58, v4)), v85), v85)
	v89 = base.Simd_g_v128_load_nc(m, l2, v57)
	v90 = base.Simd_g_i16x8_mul(v87, v89)
	base.Simd_g_v128_store(m, l0, v57, v90)
	v97 = base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v47, base.Simd_g_const(&F_QuantizeBlock_SSE41__k11)), base.Simd_g_i8x16_swizzle(v87, base.Simd_g_const(&F_QuantizeBlock_SSE41__k12)))
	base.Simd_g_v128_store(m, l1, v57, v97)
	v104 = base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v47, base.Simd_g_const(&F_QuantizeBlock_SSE41__k13)), base.Simd_g_i8x16_swizzle(v87, base.Simd_g_const(&F_QuantizeBlock_SSE41__k14)))
	base.Simd_g_v128_store(m, l1, v16, v104)
	return base.B2i32(base.Simd_g_i8x16_bitmask(base.Simd_g_i8x16_eq(base.Simd_g_i8x16_narrow_i16x8_s(v97, v104), v4)) != int32(_a_F_QuantizeBlock_SSE41_0))
}

var F_QuantizeBlock_SSE41__k0 = [2]uint64{0x0, 0x0}
var F_QuantizeBlock_SSE41__k1 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_QuantizeBlock_SSE41__k2 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_QuantizeBlock_SSE41__k3 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_QuantizeBlock_SSE41__k4 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_QuantizeBlock_SSE41__k5 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_QuantizeBlock_SSE41__k6 = [2]uint64{0x7ff07ff07ff07ff, 0x7ff07ff07ff07ff}
var F_QuantizeBlock_SSE41__k7 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_QuantizeBlock_SSE41__k8 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_QuantizeBlock_SSE41__k9 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_QuantizeBlock_SSE41__k10 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_QuantizeBlock_SSE41__k11 = [2]uint64{0x1008f8e8f8e8f8e, 0x8f8e8f8e8f8e8f8e}
var F_QuantizeBlock_SSE41__k12 = [2]uint64{0x8f8e090803020100, 0xd0c070605040b0a}
var F_QuantizeBlock_SSE41__k13 = [2]uint64{0x5040b0a09080302, 0xf0e0d0c07068f8e}
var F_QuantizeBlock_SSE41__k14 = [2]uint64{0x8f8e8f8e8f8e8f8e, 0x8f8e8f8e8f8e0f0e}

func F_QuantizeLevels(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) int32 {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v34 int32
	_ = v34
	var v36 int32
	_ = v36
	var v421 int32
	_ = v421
	var v423 int32
	_ = v423
	var v425 int32
	_ = v425
	var v429 int32
	_ = v429
	var v435 int32
	_ = v435
	var v441 int32
	_ = v441
	var v447 int32
	_ = v447
	var v448 int32
	_ = v448
	var v449 int32
	_ = v449
	var v473 int32
	_ = v473
	var v474 int32
	_ = v474
	var v475 int32
	_ = v475
	var v476 int32
	_ = v476
	var v478 int32
	_ = v478
	var v479 int32
	_ = v479
	var v480 int32
	_ = v480
	var v487 int32
	_ = v487
	var v490 int32
	_ = v490
	var v491 int32
	_ = v491
	var v496 int32
	_ = v496
	var v498 int32
	_ = v498
	var v500 int32
	_ = v500
	var v502 int32
	_ = v502
	var v503 int32
	_ = v503
	var v508 int32
	_ = v508
	var v510 int32
	_ = v510
	var v514 int32
	_ = v514
	var v520 int32
	_ = v520
	var v521 int32
	_ = v521
	var v522 int32
	_ = v522
	var v550 int32
	_ = v550
	var v553 int32
	_ = v553
	var v554 int32
	_ = v554
	var v559 int32
	_ = v559
	var v561 int32
	_ = v561
	var v573 int32
	_ = v573
	var v574 int32
	_ = v574
	var v575 int32
	_ = v575
	var v601 int32
	_ = v601
	var v607 int32
	_ = v607
	var v608 float64
	_ = v608
	var v610 float64
	_ = v610
	var v611 float64
	_ = v611
	var v621 int32
	_ = v621
	var v623 int32
	_ = v623
	var v624 int32
	_ = v624
	var v646 base.V128
	_ = v646
	var v658 base.V128
	_ = v658
	var v666 int32
	_ = v666
	var v674 int32
	_ = v674
	var v707 int32
	_ = v707
	var v712 int32
	_ = v712
	var v747 int32
	_ = v747
	var v764 int32
	_ = v764
	var v785 int32
	_ = v785
	var v786 int32
	_ = v786
	var v787 int32
	_ = v787
	var v788 int32
	_ = v788
	var v789 int32
	_ = v789
	var v797 int32
	_ = v797
	var v803 int32
	_ = v803
	var v836 float64
	_ = v836
	var v848 int32
	_ = v848
	var v849 int32
	_ = v849
	var v1095 int32
	_ = v1095
	var v1099 int32
	_ = v1099
	var v1103 int32
	_ = v1103
	var v1131 int32
	_ = v1131
	var v1141 int32
	_ = v1141
	var v1146 int32
	_ = v1146
	var v1176 float64
	_ = v1176
	var v1178 int32
	_ = v1178
	var v1179 float64
	_ = v1179
	var v1186 int32
	_ = v1186
	var v1191 int32
	_ = v1191
	var v1193 int32
	_ = v1193
	var v1199 int32
	_ = v1199
	var v1200 int32
	_ = v1200
	var v1201 float64
	_ = v1201
	var v1206 int32
	_ = v1206
	var v1207 float64
	_ = v1207
	var v1256 int32
	_ = v1256
	var v1259 int32
	_ = v1259
	var v1264 int32
	_ = v1264
	var v1291 int32
	_ = v1291
	var v1294 float64
	_ = v1294
	var v1302 int32
	_ = v1302
	var v1309 float64
	_ = v1309
	var v1314 float64
	_ = v1314
	var v1322 int32
	_ = v1322
	var v1329 float64
	_ = v1329
	var v1335 int32
	_ = v1335
	var v1345 int32
	_ = v1345
	var v1375 int32
	_ = v1375
	var v1377 float64
	_ = v1377
	var v1388 float64
	_ = v1388
	var v1428 int32
	_ = v1428
	var v1432 float64
	_ = v1432
	var v1433 float64
	_ = v1433
	var v1434 int32
	_ = v1434
	var v1440 int32
	_ = v1440
	var v1441 float64
	_ = v1441
	var v1446 int32
	_ = v1446
	var v1452 int32
	_ = v1452
	var v1453 int32
	_ = v1453
	var v1457 int32
	_ = v1457
	var v1469 float64
	_ = v1469
	var v1485 int32
	_ = v1485
	var v1488 int32
	_ = v1488
	var v1489 int32
	_ = v1489
	var v1491 int32
	_ = v1491
	var v1492 int32
	_ = v1492
	var v1495 float64
	_ = v1495
	var v1496 float64
	_ = v1496
	var v1499 int32
	_ = v1499
	var v1506 int32
	_ = v1506
	var v1510 float64
	_ = v1510
	var v1511 float64
	_ = v1511
	var v1512 int32
	_ = v1512
	var v1517 float64
	_ = v1517
	var v1520 int32
	_ = v1520
	var v1543 float64
	_ = v1543
	var v1561 int32
	_ = v1561
	var v1567 int32
	_ = v1567
	var v1568 int32
	_ = v1568
	var v1572 int32
	_ = v1572
	var v1583 int32
	_ = v1583
	var v1584 int32
	_ = v1584
	var v1588 int32
	_ = v1588
	var v1616 int32
	_ = v1616
	var v1617 int32
	_ = v1617
	var v1618 int32
	_ = v1618
	var v1621 int32
	_ = v1621
	var v1622 base.V128
	_ = v1622
	var v1627 int32
	_ = v1627
	var v1631 float64
	_ = v1631
	var v1632 int32
	_ = v1632
	var v1635 base.V128
	_ = v1635
	var v1637 float64
	_ = v1637
	var v1645 int32
	_ = v1645
	var v1647 int32
	_ = v1647
	var v1648 int32
	_ = v1648
	var v1649 float64
	_ = v1649
	var v1657 int32
	_ = v1657
	var v1659 int32
	_ = v1659
	var v1663 int32
	_ = v1663
	var v1671 int32
	_ = v1671
	var v1688 int32
	_ = v1688
	var v1718 int32
	_ = v1718
	var v1719 int32
	_ = v1719
	var v1723 int32
	_ = v1723
	var v1752 int32
	_ = v1752
	var v1756 float64
	_ = v1756
	var v1758 float64
	_ = v1758
	var v1766 int32
	_ = v1766
	var v1768 int32
	_ = v1768
	var v1775 int32
	_ = v1775
	var v1815 int32
	_ = v1815
	var v1817 int32
	_ = v1817
	var v1822 int32
	_ = v1822
	var v1850 int32
	_ = v1850
	var v1851 int32
	_ = v1851
	var v1852 base.V128
	_ = v1852
	var v1856 base.V128
	_ = v1856
	var v1857 int32
	_ = v1857
	var v1862 int32
	_ = v1862
	var v1867 int32
	_ = v1867
	var v1880 base.V128
	_ = v1880
	var v1904 base.V128
	_ = v1904
	var v1926 base.V128
	_ = v1926
	var v1946 base.V128
	_ = v1946
	var v1949 base.V128
	_ = v1949
	var v1952 base.V128
	_ = v1952
	var v1955 base.V128
	_ = v1955
	var v1958 base.V128
	_ = v1958
	var v1961 base.V128
	_ = v1961
	var v1964 base.V128
	_ = v1964
	var v1967 base.V128
	_ = v1967
	var v1970 base.V128
	_ = v1970
	var v1973 base.V128
	_ = v1973
	var v1976 base.V128
	_ = v1976
	var v1979 base.V128
	_ = v1979
	var v1982 base.V128
	_ = v1982
	var v1985 base.V128
	_ = v1985
	var v1988 base.V128
	_ = v1988
	var v1991 base.V128
	_ = v1991
	var v1997 int32
	_ = v1997
	var v2001 int32
	_ = v2001
	var v2035 int32
	_ = v2035
	var v2040 int32
	_ = v2040
	var v2069 int32
	_ = v2069
	var v2071 int32
	_ = v2071
	var v2076 int32
	_ = v2076
	var v2117 int64
	_ = v2117
	var v2135 int64
	_ = v2135
	var v2152 int32
	_ = v2152
	var v2162 int32
	_ = v2162
	v6 = int32(0)
	v34 = m.G0
	v36 = v34 - int32(_a_F_QuantizeLevels_0)
	m.G0 = v36
	goto L3
L1:
	;
	goto L17
L3:
	;
	base.MemoryFill(m, v36+int32(_a_F_QuantizeLevels_1), v6, int32(1024))
	goto L1
L15:
	;
	goto L31
L17:
	;
	base.MemoryFill(m, v36+int32(_a_F_QuantizeLevels_2), int32(0), int32(1024))
	goto L15
L29:
	;
	if l0 == int32(0) {
		v2162 = v6
		goto L43
	} else {
		goto L44
	}
L31:
	;
	base.MemoryFill(m, v36+int32(_a_F_QuantizeLevels_3), int32(0), int32(2048))
	goto L29
L43:
	;
	m.G0 = v36 + int32(_a_F_QuantizeLevels_0)
	return v2162
L44:
	;
	if l1 < int32(1) {
		v2162 = v6
		goto L43
	} else {
		goto L45
	}
L45:
	;
	if l2 < int32(1) {
		v2162 = v6
		goto L43
	} else {
		goto L46
	}
L46:
	;
	if base.Ui32(l3+int32(-257)) < base.Ui32(int32(-255)) {
		v2162 = v6
		goto L43
	} else {
		goto L47
	}
L47:
	;
	v421 = l2 * l1
	if v421 != 0 {
		goto L49
	} else {
		goto L50
	}
L48:
	;
	if v574 <= l3 {
		v2135 = int64(0)
		goto L76
	} else {
		goto L77
	}
L49:
	;
	v425 = int32(1)
	if v421 != v425 {
		goto L52
	} else {
		goto L53
	}
L50:
	;
	v423 = int32(0)
	v573 = int32(255)
	v574 = v423
	v575 = v423
	goto L48
L51:
	;
	if v421&v425 == int32(0) {
		v573 = v520
		v574 = v521
		v575 = v522
		goto L48
	} else {
		goto L69
	}
L52:
	;
	v435 = int32(0)
	v441 = v435
	v447 = int32(255)
	v448 = v435
	v449 = v435
	goto L54
L53:
	;
	v429 = int32(0)
	v514 = v429
	v520 = int32(255)
	v521 = v429
	v522 = v429
	goto L51
L54:
	;
	v473 = v36 + int32(_a_F_QuantizeLevels_1)
	v474 = l0 + v441
	v475 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v474))))
	v476 = int32(2)
	v478 = v473 + v475<<(uint(v476)%32)
	v479 = *(*int32)(unsafe.Add(mBase, uint32(v478)))
	v480 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v478))) = v479 + v480
	v487 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v474+v480))))
	v490 = v473 + v487<<(uint(v476)%32)
	v491 = *(*int32)(unsafe.Add(mBase, uint32(v490)))
	*(*int32)(unsafe.Add(mBase, uint32(v490))) = v491 + v480
	if v475 < v449 {
		goto L56
	} else {
		goto L57
	}
L55:
	;
	v514 = v510
	v520 = v502
	v521 = v508
	v522 = v498
	goto L51
L56:
	;
	v496 = v449
	goto L58
L57:
	;
	v496 = v475
	goto L58
L58:
	;
	if v487 < v496 {
		goto L59
	} else {
		goto L60
	}
L59:
	;
	v498 = v496
	goto L61
L60:
	;
	v498 = v487
	goto L61
L61:
	;
	if v447 < v475 {
		goto L62
	} else {
		goto L63
	}
L62:
	;
	v500 = v447
	goto L64
L63:
	;
	v500 = v475
	goto L64
L64:
	;
	if v500 < v487 {
		goto L65
	} else {
		goto L66
	}
L65:
	;
	v502 = v500
	goto L67
L66:
	;
	v502 = v487
	goto L67
L67:
	;
	v503 = int32(0)
	v508 = v448 + base.B2i32(v479 == v503) + base.B2i32(v491 == v503)
	v510 = v441 + int32(2)
	if v421&int32(-2) != v510 {
		v441 = v510
		v447 = v502
		v448 = v508
		v449 = v498
		goto L54
	} else {
		goto L68
	}
L68:
	;
	goto L55
L69:
	;
	v550 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v514))))
	v553 = v36 + int32(_a_F_QuantizeLevels_1) + v550<<(uint(int32(2))%32)
	v554 = *(*int32)(unsafe.Add(mBase, uint32(v553)))
	*(*int32)(unsafe.Add(mBase, uint32(v553))) = v554 + int32(1)
	if v550 < v522 {
		goto L70
	} else {
		goto L71
	}
L70:
	;
	v559 = v522
	goto L72
L71:
	;
	v559 = v550
	goto L72
L72:
	;
	if v520 < v550 {
		goto L73
	} else {
		goto L74
	}
L73:
	;
	v561 = v520
	goto L75
L74:
	;
	v561 = v550
	goto L75
L75:
	;
	v573 = v561
	v574 = v521 + base.B2i32(v554 == int32(0))
	v575 = v559
	goto L48
L76:
	;
	v2152 = int32(1)
	if l4 == int32(0) {
		v2162 = v2152
		goto L43
	} else {
		goto L194
	}
L77:
	;
	v601 = int32(0)
	if v601 < l3 {
		goto L79
	} else {
		goto L80
	}
L78:
	;
	v785 = v36 + int32(_a_F_QuantizeLevels_2)
	v786 = int32(2)
	v787 = v573 << (uint(v786) % 32)
	v788 = v785 + v787
	v789 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v788))) = v789
	*(*int32)(unsafe.Add(mBase, uint32(v785+v575<<(uint(v786)%32)))) = v764
	v797 = int32(1)
	v803 = v764 + int32(-1)
	v836 = float64(1e+38)
	v848 = v789
	goto L90
L79:
	;
	v607 = l3 + int32(-1)
	v608 = base.F64_convert_i32_u(v607)
	v610 = base.F64_convert_i32_s(v575 - v573)
	v611 = base.F64_convert_i32_u(v573)
	if l3 == int32(1) {
		v674 = v601
		goto L81
	} else {
		goto L82
	}
L80:
	;
	v764 = l3 + int32(-1)
	goto L78
L81:
	;
	v707 = v36 + int32(_a_F_QuantizeLevels_3) + v674<<(uint(int32(3))%32)
	v712 = v674
	goto L87
L82:
	;
	v621 = l3 & int32(510)
	v623 = v36 + int32(_a_F_QuantizeLevels_3)
	v624 = v621
	v646 = base.Simd_g_const(&F_QuantizeLevels__k0)
	goto L83
L83:
	;
	v658 = base.Simd_g_f64x2_add(base.Simd_g_f64x2_div(base.Simd_g_f64x2_mul(base.Simd_g_f64x2_splat(v610), base.Simd_g_f64x2_convert_low_i32x4_u(v646)), base.Simd_g_f64x2_splat(v608)), base.Simd_g_f64x2_splat(v611))
	base.Simd_g_v128_store(m, v623, int32(0), v658)
	v666 = v624 + int32(-2)
	if v666 != 0 {
		v623 = v623 + int32(16)
		v624 = v666
		v646 = base.Simd_g_i32x4_add(v646, base.Simd_g_const(&F_QuantizeLevels__k1))
		goto L83
	} else {
		goto L85
	}
L84:
	;
	if v621 == l3 {
		v764 = v607
		goto L78
	} else {
		goto L86
	}
L85:
	;
	goto L84
L86:
	;
	v674 = v621
	goto L81
L87:
	;
	*(*float64)(unsafe.Add(mBase, uint32(v707))) = base.F64_add(base.F64_div(base.F64_mul(v610, base.F64_convert_i32_u(v712)), v608), v611)
	v747 = v712 + int32(1)
	if l3 != v747 {
		v707 = v707 + int32(8)
		v712 = v747
		goto L87
	} else {
		goto L89
	}
L88:
	;
	v764 = v607
	goto L78
L89:
	;
	goto L88
L90:
	;
	v849 = int32(2048)
	goto L94
L91:
	;
	if v575 < v573 {
		goto L160
	} else {
		goto L161
	}
L92:
	;
	goto L108
L94:
	;
	base.MemoryFill(m, v36+v849, int32(0), v849)
	goto L92
L106:
	;
	v1095 = base.B2i32(v575 < v573)
	if v575 < v573 {
		goto L120
	} else {
		goto L121
	}
L108:
	;
	base.MemoryFill(m, v36, int32(0), int32(2048))
	goto L106
L120:
	;
	if l3 < int32(3) {
		goto L136
	} else {
		goto L137
	}
L121:
	;
	v1099 = v573
	v1103 = int32(0)
	goto L122
L122:
	;
	if v764 < v1103 {
		goto L124
	} else {
		goto L125
	}
L123:
	;
	goto L120
L124:
	;
	v1131 = v1103
	goto L126
L125:
	;
	v1131 = v764
	goto L126
L126:
	;
	v1141 = v36 + int32(_a_F_QuantizeLevels_3) + v1103<<(uint(int32(3))%32)
	v1146 = v1103
	goto L128
L127:
	;
	v1191 = v1099 << (uint(int32(2)) % 32)
	v1193 = *(*int32)(unsafe.Add(mBase, uint32(v36+int32(_a_F_QuantizeLevels_1)+v1191)))
	if v1193 < int32(1) {
		goto L133
	} else {
		goto L134
	}
L128:
	;
	if v1131 != v1146 {
		goto L130
	} else {
		goto L131
	}
L129:
	;
	v1186 = v1146 + int32(0)
	goto L127
L130:
	;
	v1176 = *(*float64)(unsafe.Add(mBase, uint32(v1141)))
	v1178 = v1141 + int32(8)
	v1179 = *(*float64)(unsafe.Add(mBase, uint32(v1178)))
	if base.F64_lt(base.F64_add(v1176, v1179), base.F64_convert_i32_s(v1099<<(uint(int32(1))%32))) != 0 {
		v1141 = v1178
		v1146 = v1146 + int32(1)
		goto L128
	} else {
		goto L132
	}
L131:
	;
	v1186 = v1131
	goto L127
L132:
	;
	goto L129
L133:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v36+int32(_a_F_QuantizeLevels_2)+v1191))) = v1186
	if v1099 != v575 {
		v1099 = v1099 + int32(1)
		v1103 = v1186
		goto L122
	} else {
		goto L135
	}
L134:
	;
	v1199 = v1186 << (uint(int32(3)) % 32)
	v1200 = v36 + int32(2048) + v1199
	v1201 = *(*float64)(unsafe.Add(mBase, uint32(v1200)))
	*(*float64)(unsafe.Add(mBase, uint32(v1200))) = base.F64_add(v1201, base.F64_convert_i32_s(v1193*v1099))
	v1206 = v36 + v1199
	v1207 = *(*float64)(unsafe.Add(mBase, uint32(v1206)))
	*(*float64)(unsafe.Add(mBase, uint32(v1206))) = base.F64_add(v1207, base.F64_convert_i32_u(v1193))
	goto L133
L135:
	;
	goto L123
L136:
	;
	if v575 < v573 {
		v1543 = float64(0)
		goto L149
	} else {
		goto L150
	}
L137:
	;
	if v764 == int32(2) {
		v1345 = int32(1)
		goto L138
	} else {
		goto L139
	}
L138:
	;
	if v803&v797 == int32(0) {
		goto L136
	} else {
		goto L147
	}
L139:
	;
	v1256 = int32(0)
	v1259 = v1256
	v1264 = v1256
	goto L140
L140:
	;
	v1291 = v36 + v1264
	v1294 = *(*float64)(unsafe.Add(mBase, uint32(v1291+int32(8))))
	if base.F64_gt(v1294, float64(0)) == int32(0) {
		goto L142
	} else {
		goto L143
	}
L141:
	;
	v1345 = v1259 + int32(3)
	goto L138
L142:
	;
	v1314 = *(*float64)(unsafe.Add(mBase, uint32(v1291+int32(16))))
	if base.F64_gt(v1314, float64(0)) == int32(0) {
		goto L144
	} else {
		goto L145
	}
L143:
	;
	v1302 = int32(8)
	v1309 = *(*float64)(unsafe.Add(mBase, uint32(v36+int32(2048)+v1264+v1302)))
	*(*float64)(unsafe.Add(mBase, uint32(v36+int32(_a_F_QuantizeLevels_3)+v1264+v1302))) = base.F64_div(v1309, v1294)
	goto L142
L144:
	;
	v1335 = v1259 + int32(2)
	if v803&int32(-2) != v1335 {
		v1259 = v1335
		v1264 = v1264 + int32(16)
		goto L140
	} else {
		goto L146
	}
L145:
	;
	v1322 = int32(16)
	v1329 = *(*float64)(unsafe.Add(mBase, uint32(v36+int32(2048)+v1264+v1322)))
	*(*float64)(unsafe.Add(mBase, uint32(v36+int32(_a_F_QuantizeLevels_3)+v1264+v1322))) = base.F64_div(v1329, v1314)
	goto L144
L146:
	;
	goto L141
L147:
	;
	v1375 = v1345 << (uint(int32(3)) % 32)
	v1377 = *(*float64)(unsafe.Add(mBase, uint32(v36+v1375)))
	if base.F64_gt(v1377, float64(0)) == int32(0) {
		goto L136
	} else {
		goto L148
	}
L148:
	;
	v1388 = *(*float64)(unsafe.Add(mBase, uint32(v36+int32(2048)+v1375)))
	*(*float64)(unsafe.Add(mBase, uint32(v36+int32(_a_F_QuantizeLevels_3)+v1375))) = base.F64_div(v1388, v1377)
	goto L136
L149:
	;
	if base.F64_lt(base.F64_sub(v836, v1543), base.F64_mul(base.F64_convert_i32_u(v421), float64(0.0001))) != 0 {
		goto L157
	} else {
		goto L158
	}
L150:
	;
	if (v575+v573)&v797 != 0 {
		v1440 = v573
		v1441 = float64(0)
		goto L151
	} else {
		goto L152
	}
L151:
	;
	if v575 == v573 {
		v1543 = v1441
		goto L149
	} else {
		goto L153
	}
L152:
	;
	v1428 = *(*int32)(unsafe.Add(mBase, uint32(v788)))
	v1432 = *(*float64)(unsafe.Add(mBase, uint32(v36+int32(_a_F_QuantizeLevels_3)+v1428<<(uint(int32(3))%32))))
	v1433 = base.F64_sub(base.F64_convert_i32_s(v573), v1432)
	v1434 = *(*int32)(unsafe.Add(mBase, uint32(v36+int32(_a_F_QuantizeLevels_1)+v787)))
	v1440 = v573 + v797
	v1441 = base.F64_add(base.F64_mul(base.F64_mul(v1433, base.F64_convert_i32_s(v1434)), v1433), float64(0))
	goto L151
L153:
	;
	v1446 = v1440 << (uint(int32(2)) % 32)
	v1452 = v36 + int32(_a_F_QuantizeLevels_2) + v1446
	v1453 = v36 + int32(_a_F_QuantizeLevels_1) + v1446
	v1457 = v1440
	v1469 = v1441
	goto L154
L154:
	;
	v1485 = v1457 + int32(1)
	v1488 = v36 + int32(_a_F_QuantizeLevels_3)
	v1489 = int32(4)
	v1491 = *(*int32)(unsafe.Add(mBase, uint32(v1452+v1489)))
	v1492 = int32(3)
	v1495 = *(*float64)(unsafe.Add(mBase, uint32(v1488+v1491<<(uint(v1492)%32))))
	v1496 = base.F64_sub(base.F64_convert_i32_s(v1485), v1495)
	v1499 = *(*int32)(unsafe.Add(mBase, uint32(v1453+v1489)))
	v1506 = *(*int32)(unsafe.Add(mBase, uint32(v1452)))
	v1510 = *(*float64)(unsafe.Add(mBase, uint32(v1488+v1506<<(uint(v1492)%32))))
	v1511 = base.F64_sub(base.F64_convert_i32_s(v1457), v1510)
	v1512 = *(*int32)(unsafe.Add(mBase, uint32(v1453)))
	v1517 = base.F64_add(base.F64_mul(base.F64_mul(v1496, base.F64_convert_i32_s(v1499)), v1496), base.F64_add(base.F64_mul(base.F64_mul(v1511, base.F64_convert_i32_s(v1512)), v1511), v1469))
	v1520 = int32(8)
	if v1485 != v575 {
		v1452 = v1452 + v1520
		v1453 = v1453 + v1520
		v1457 = v1457 + int32(2)
		v1469 = v1517
		goto L154
	} else {
		goto L156
	}
L155:
	;
	v1543 = v1517
	goto L149
L156:
	;
	goto L155
L157:
	;
	goto L91
L158:
	;
	v1561 = v848 + int32(1)
	if v1561 != int32(6) {
		v836 = v1543
		v848 = v1561
		goto L90
	} else {
		goto L159
	}
L159:
	;
	goto L157
L160:
	;
	if v421 == int32(0) {
		goto L181
	} else {
		goto L182
	}
L161:
	;
	v1567 = int32(1)
	v1568 = v575 - v573 + v1567
	if base.Ui32(v1567) < base.Ui32(v1568) {
		goto L163
	} else {
		goto L164
	}
L162:
	;
	v1718 = v36 + int32(2048) + v1688
	v1719 = v575 - v1688 + int32(1)
	v1723 = v36 + int32(_a_F_QuantizeLevels_2) + v1688<<(uint(int32(2))%32)
	goto L175
L163:
	;
	v1572 = v1568 & int32(-2)
	v1583 = v36 + int32(2048) + v573
	v1584 = v1572
	v1588 = v36 + int32(_a_F_QuantizeLevels_2) + v573<<(uint(int32(2))%32)
	goto L165
L164:
	;
	v1688 = v573
	goto L162
L165:
	;
	v1616 = v36 + int32(_a_F_QuantizeLevels_3)
	v1617 = *(*int32)(unsafe.Add(mBase, uint32(v1588)))
	v1618 = int32(3)
	v1621 = int32(0)
	v1622 = base.Simd_g_v128_load64_splat(m, v1616+v1617<<(uint(v1618)%32), v1621)
	v1627 = *(*int32)(unsafe.Add(mBase, uint32(v1588+int32(4))))
	v1631 = *(*float64)(unsafe.Add(mBase, uint32(v1616+v1627<<(uint(v1618)%32))))
	v1632 = int32(1)
	v1635 = base.Simd_g_f64x2_add(base.Simd_g_f64x2_replace_lane_l1(v1622, v1631), base.Simd_g_const(&F_QuantizeLevels__k2))
	v1637 = base.Simd_g_f64x2_extract_lane_l1(v1635)
	if base.F64_lt(v1637, float64(4.294967296e+09))&base.F64_ge(v1637, float64(0)) == v1621 {
		goto L168
	} else {
		goto L169
	}
L166:
	;
	if v1568 == v1572 {
		goto L160
	} else {
		goto L174
	}
L167:
	;
	v1648 = int32(0)
	v1649 = base.Simd_g_f64x2_extract_lane_l0(v1635)
	if base.F64_lt(v1649, float64(4.294967296e+09))&base.F64_ge(v1649, float64(0)) == v1648 {
		goto L171
	} else {
		goto L172
	}
L168:
	;
	v1647 = int32(0)
	goto L167
L169:
	;
	v1645 = base.I32_trunc_f64_u(v1637)
	v1647 = v1645
	goto L167
L170:
	;
	v1663 = int32(0)
	base.Simd_g_v128_store16_lane_l0(m, v1583, v1663, base.Simd_g_i8x16_replace_lane_l1(base.Simd_g_i8x16_splat(v1659), v1647))
	v1671 = v1584 + int32(-2)
	if v1671 != 0 {
		v1583 = v1583 + int32(2)
		v1584 = v1671
		v1588 = v1588 + int32(8)
		goto L165
	} else {
		goto L173
	}
L171:
	;
	v1659 = int32(0)
	goto L170
L172:
	;
	v1657 = base.I32_trunc_f64_u(v1649)
	v1659 = v1657
	goto L170
L173:
	;
	goto L166
L174:
	;
	v1688 = v573 + v1572
	goto L162
L175:
	;
	v1752 = *(*int32)(unsafe.Add(mBase, uint32(v1723)))
	v1756 = *(*float64)(unsafe.Add(mBase, uint32(v36+int32(_a_F_QuantizeLevels_3)+v1752<<(uint(int32(3))%32))))
	v1758 = base.F64_add(v1756, float64(0.5))
	if base.F64_lt(v1758, float64(4.294967296e+09))&base.F64_ge(v1758, float64(0)) == int32(0) {
		goto L178
	} else {
		goto L179
	}
L176:
	;
	goto L160
L177:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1718))) = uint8(v1768)
	v1775 = v1719 + int32(-1)
	if v1775 != 0 {
		v1718 = v1718 + int32(1)
		v1719 = v1775
		v1723 = v1723 + int32(4)
		goto L175
	} else {
		goto L180
	}
L178:
	;
	v1768 = int32(0)
	goto L177
L179:
	;
	v1766 = base.I32_trunc_f64_u(v1758)
	v1768 = v1766
	goto L177
L180:
	;
	goto L176
L181:
	;
	if base.F64_lt(v1543, float64(1.8446744073709552e+19))&base.F64_ge(v1543, float64(0)) == int32(0) {
		goto L192
	} else {
		goto L193
	}
L182:
	;
	if base.Ui32(v421) <= base.Ui32(int32(15)) {
		v2001 = int32(0)
		goto L183
	} else {
		goto L184
	}
L183:
	;
	v2035 = v421 - v2001
	v2040 = l0 + v2001
	goto L189
L184:
	;
	v1815 = v421 & int32(-16)
	v1817 = v1815
	v1822 = l0
	goto L185
L185:
	;
	v1850 = v36 + int32(2048)
	v1851 = int32(0)
	v1852 = base.Simd_g_v128_load(m, v1822, v1851)
	v1856 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_swizzle_c(v1852, base.Simd_g_const(&F_QuantizeLevels__k3))))
	v1857 = int32(3)
	v1862 = int32(2)
	v1867 = int32(1)
	v1880 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_swizzle_c(v1852, base.Simd_g_const(&F_QuantizeLevels__k4))))
	v1904 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_swizzle_c(v1852, base.Simd_g_const(&F_QuantizeLevels__k5))))
	v1926 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v1852))
	v1946 = base.Simd_g_v128_load8_splat(m, v1850+base.Simd_g_i32x4_extract_lane_l0(v1926), v1851)
	v1949 = base.Simd_g_v128_load8_lane_l1(m, v1850+base.Simd_g_i32x4_extract_lane_l1(v1926), v1851, v1946)
	v1952 = base.Simd_g_v128_load8_lane_l2(m, v1850+base.Simd_g_i32x4_extract_lane_l2(v1926), v1851, v1949)
	v1955 = base.Simd_g_v128_load8_lane_l3(m, v1850+base.Simd_g_i32x4_extract_lane_l3(v1926), v1851, v1952)
	v1958 = base.Simd_g_v128_load8_lane_l4(m, v1850+base.Simd_g_i32x4_extract_lane_l0(v1904), v1851, v1955)
	v1961 = base.Simd_g_v128_load8_lane_l5(m, v1850+base.Simd_g_i32x4_extract_lane_l1(v1904), v1851, v1958)
	v1964 = base.Simd_g_v128_load8_lane_l6(m, v1850+base.Simd_g_i32x4_extract_lane_l2(v1904), v1851, v1961)
	v1967 = base.Simd_g_v128_load8_lane_l7(m, v1850+base.Simd_g_i32x4_extract_lane_l3(v1904), v1851, v1964)
	v1970 = base.Simd_g_v128_load8_lane_l8(m, v1850+base.Simd_g_i32x4_extract_lane_l0(v1880), v1851, v1967)
	v1973 = base.Simd_g_v128_load8_lane_l9(m, v1850+base.Simd_g_i32x4_extract_lane_l1(v1880), v1851, v1970)
	v1976 = base.Simd_g_v128_load8_lane_l10(m, v1850+base.Simd_g_i32x4_extract_lane_l2(v1880), v1851, v1973)
	v1979 = base.Simd_g_v128_load8_lane_l11(m, v1850+base.Simd_g_i32x4_extract_lane_l3(v1880), v1851, v1976)
	v1982 = base.Simd_g_v128_load8_lane_l12(m, v1850+base.Simd_g_i32x4_extract_lane_l0(v1856), v1851, v1979)
	v1985 = base.Simd_g_v128_load8_lane_l13(m, v1850+base.Simd_g_i32x4_extract_lane_l1(v1856), v1851, v1982)
	v1988 = base.Simd_g_v128_load8_lane_l14(m, v1850+base.Simd_g_i32x4_extract_lane_l2(v1856), v1851, v1985)
	v1991 = base.Simd_g_v128_load8_lane_l15(m, v1850+base.Simd_g_i32x4_extract_lane_l3(v1856), v1851, v1988)
	base.Simd_g_v128_store(m, v1822, v1851, v1991)
	v1997 = v1817 + int32(-16)
	if v1997 != 0 {
		v1817 = v1997
		v1822 = v1822 + int32(16)
		goto L185
	} else {
		goto L187
	}
L186:
	;
	if v421 == v1815 {
		goto L181
	} else {
		goto L188
	}
L187:
	;
	goto L186
L188:
	;
	v2001 = v1815
	goto L183
L189:
	;
	v2069 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2040))))
	v2071 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v36+int32(2048)+v2069))))
	*(*uint8)(unsafe.Add(mBase, uint32(v2040))) = uint8(v2071)
	v2076 = v2035 + int32(-1)
	if v2076 != 0 {
		v2035 = v2076
		v2040 = v2040 + int32(1)
		goto L189
	} else {
		goto L191
	}
L190:
	;
	goto L181
L191:
	;
	goto L190
L192:
	;
	v2135 = int64(0)
	goto L76
L193:
	;
	v2117 = base.I64_trunc_f64_u(v1543)
	v2135 = v2117
	goto L76
L194:
	;
	*(*int64)(unsafe.Add(mBase, uint32(l4))) = v2135
	v2162 = v2152
	goto L43
}

var F_QuantizeLevels__k0 = [2]uint64{0x100000000, 0x0}
var F_QuantizeLevels__k1 = [2]uint64{0x200000002, 0x200000002}
var F_QuantizeLevels__k2 = [2]uint64{0x3fe0000000000000, 0x3fe0000000000000}
var F_QuantizeLevels__k3 = [2]uint64{0xf0e0d0c, 0x0}
var F_QuantizeLevels__k4 = [2]uint64{0xb0a0908, 0x0}
var F_QuantizeLevels__k5 = [2]uint64{0x7060504, 0x0}
