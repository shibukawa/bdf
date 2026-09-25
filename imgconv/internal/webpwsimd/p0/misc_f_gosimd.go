//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"

func F_FTransform2_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	var v4 base.V128
	_ = v4
	var v24 int32
	_ = v24
	var v25 base.V128
	_ = v25
	var v27 base.V128
	_ = v27
	var v30 base.V128
	_ = v30
	var v33 base.V128
	_ = v33
	var v34 int32
	_ = v34
	var v35 base.V128
	_ = v35
	var v39 base.V128
	_ = v39
	var v42 base.V128
	_ = v42
	var v43 base.V128
	_ = v43
	var v45 base.V128
	_ = v45
	var v46 base.V128
	_ = v46
	var v47 int32
	_ = v47
	var v48 base.V128
	_ = v48
	var v52 base.V128
	_ = v52
	var v55 base.V128
	_ = v55
	var v56 int32
	_ = v56
	var v57 base.V128
	_ = v57
	var v61 base.V128
	_ = v61
	var v64 base.V128
	_ = v64
	var v68 base.V128
	_ = v68
	var v69 base.V128
	_ = v69
	var v70 base.V128
	_ = v70
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
	var v78 base.V128
	_ = v78
	var v79 base.V128
	_ = v79
	var v80 base.V128
	_ = v80
	var v82 base.V128
	_ = v82
	var v84 int32
	_ = v84
	var v86 base.V128
	_ = v86
	var v88 base.V128
	_ = v88
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
	var v99 base.V128
	_ = v99
	var v100 base.V128
	_ = v100
	var v101 base.V128
	_ = v101
	var v102 base.V128
	_ = v102
	var v103 base.V128
	_ = v103
	var v104 base.V128
	_ = v104
	var v105 base.V128
	_ = v105
	var v107 base.V128
	_ = v107
	var v109 int32
	_ = v109
	var v111 base.V128
	_ = v111
	var v115 base.V128
	_ = v115
	var v116 base.V128
	_ = v116
	var v118 base.V128
	_ = v118
	var v120 int32
	_ = v120
	var v121 base.V128
	_ = v121
	var v124 base.V128
	_ = v124
	var v130 base.V128
	_ = v130
	var v134 base.V128
	_ = v134
	var v136 base.V128
	_ = v136
	var v138 base.V128
	_ = v138
	var v139 base.V128
	_ = v139
	var v142 base.V128
	_ = v142
	var v143 base.V128
	_ = v143
	var v152 base.V128
	_ = v152
	var v154 base.V128
	_ = v154
	var v156 base.V128
	_ = v156
	var v160 base.V128
	_ = v160
	var v162 base.V128
	_ = v162
	var v163 base.V128
	_ = v163
	var v164 base.V128
	_ = v164
	var v166 base.V128
	_ = v166
	var v170 base.V128
	_ = v170
	var v174 base.V128
	_ = v174
	var v178 base.V128
	_ = v178
	var v181 base.V128
	_ = v181
	var v187 base.V128
	_ = v187
	var v189 base.V128
	_ = v189
	var v192 base.V128
	_ = v192
	var v197 base.V128
	_ = v197
	var v206 base.V128
	_ = v206
	var v211 base.V128
	_ = v211
	v4 = base.Simd_g_const(&F_FTransform2_SSE2__k0)
	v24 = int32(0)
	v25 = base.Simd_g_v128_load64_zero(m, l0, v24)
	v27 = base.Simd_g_const(&F_FTransform2_SSE2__k1)
	v30 = base.Simd_g_v128_load64_zero(m, l1, v24)
	v33 = base.Simd_g_i16x8_sub(base.Simd_g_i8x16_shuffle2(v25, v4, base.Simd_g_const(&F_FTransform2_SSE2__k2), base.Simd_g_const(&F_FTransform2_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v30, v4, base.Simd_g_const(&F_FTransform2_SSE2__k2), base.Simd_g_const(&F_FTransform2_SSE2__k3)))
	v34 = int32(32)
	v35 = base.Simd_g_v128_load64_zero(m, l0, v34)
	v39 = base.Simd_g_v128_load64_zero(m, l1, v34)
	v42 = base.Simd_g_i16x8_sub(base.Simd_g_i8x16_shuffle2(v35, v4, base.Simd_g_const(&F_FTransform2_SSE2__k2), base.Simd_g_const(&F_FTransform2_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v39, v4, base.Simd_g_const(&F_FTransform2_SSE2__k2), base.Simd_g_const(&F_FTransform2_SSE2__k3)))
	v43 = base.Simd_g_const(&F_FTransform2_SSE2__k4)
	v45 = base.Simd_g_const(&F_FTransform2_SSE2__k5)
	v46 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_shuffle2(v33, v42, base.Simd_g_const(&F_FTransform2_SSE2__k6), base.Simd_g_const(&F_FTransform2_SSE2__k7)), base.Simd_g_const(&F_FTransform2_SSE2__k5))
	v47 = int32(64)
	v48 = base.Simd_g_v128_load64_zero(m, l0, v47)
	v52 = base.Simd_g_v128_load64_zero(m, l1, v47)
	v55 = base.Simd_g_i16x8_sub(base.Simd_g_i8x16_shuffle2(v48, v4, base.Simd_g_const(&F_FTransform2_SSE2__k2), base.Simd_g_const(&F_FTransform2_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v52, v4, base.Simd_g_const(&F_FTransform2_SSE2__k2), base.Simd_g_const(&F_FTransform2_SSE2__k3)))
	v56 = int32(96)
	v57 = base.Simd_g_v128_load64_zero(m, l0, v56)
	v61 = base.Simd_g_v128_load64_zero(m, l1, v56)
	v64 = base.Simd_g_i16x8_sub(base.Simd_g_i8x16_shuffle2(v57, v4, base.Simd_g_const(&F_FTransform2_SSE2__k2), base.Simd_g_const(&F_FTransform2_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v61, v4, base.Simd_g_const(&F_FTransform2_SSE2__k2), base.Simd_g_const(&F_FTransform2_SSE2__k3)))
	v68 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_shuffle2(v55, v64, base.Simd_g_const(&F_FTransform2_SSE2__k6), base.Simd_g_const(&F_FTransform2_SSE2__k7)), base.Simd_g_const(&F_FTransform2_SSE2__k5))
	v69 = base.Simd_g_const(&F_FTransform2_SSE2__k8)
	v70 = base.Simd_g_i8x16_shuffle2(v46, v68, base.Simd_g_const(&F_FTransform2_SSE2__k9), base.Simd_g_const(&F_FTransform2_SSE2__k10))
	v71 = base.Simd_g_const(&F_FTransform2_SSE2__k11)
	v72 = base.Simd_g_i8x16_shuffle2(v46, v68, base.Simd_g_const(&F_FTransform2_SSE2__k12), base.Simd_g_const(&F_FTransform2_SSE2__k13))
	v73 = base.Simd_g_i16x8_add(v70, v72)
	v74 = base.Simd_g_const(&F_FTransform2_SSE2__k14)
	v76 = base.Simd_g_const(&F_FTransform2_SSE2__k15)
	v78 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_dot_i16x8_s(v73, v74), base.Simd_g_i32x4_dot_i16x8_s(v73, v76))
	v79 = base.Simd_g_i16x8_sub(v70, v72)
	v80 = base.Simd_g_const(&F_FTransform2_SSE2__k16)
	v82 = base.Simd_g_const(&F_FTransform2_SSE2__k17)
	v84 = int32(9)
	v86 = base.Simd_g_const(&F_FTransform2_SSE2__k18)
	v88 = base.Simd_g_const(&F_FTransform2_SSE2__k19)
	v92 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v79, v80), v82), v84), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v79, v86), v88), v84))
	v93 = base.Simd_g_const(&F_FTransform2_SSE2__k20)
	v94 = base.Simd_g_i8x16_shuffle2(v78, v92, base.Simd_g_const(&F_FTransform2_SSE2__k21), base.Simd_g_const(&F_FTransform2_SSE2__k22))
	v95 = base.Simd_g_const(&F_FTransform2_SSE2__k23)
	v96 = base.Simd_g_i8x16_shuffle2(v78, v92, base.Simd_g_const(&F_FTransform2_SSE2__k24), base.Simd_g_const(&F_FTransform2_SSE2__k25))
	v99 = base.Simd_g_const(&F_FTransform2_SSE2__k26)
	v100 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_shuffle2(v94, v96, base.Simd_g_const(&F_FTransform2_SSE2__k6), base.Simd_g_const(&F_FTransform2_SSE2__k7)), base.Simd_g_const(&F_FTransform2_SSE2__k26))
	v101 = base.Simd_g_const(&F_FTransform2_SSE2__k27)
	v102 = base.Simd_g_i8x16_shuffle2(v94, v96, base.Simd_g_const(&F_FTransform2_SSE2__k28), base.Simd_g_const(&F_FTransform2_SSE2__k29))
	v103 = base.Simd_g_i16x8_add(v100, v102)
	v104 = base.Simd_g_const(&F_FTransform2_SSE2__k30)
	v105 = base.Simd_g_i16x8_add(v103, v104)
	v107 = base.Simd_g_i8x16_shuffle2(v103, v103, base.Simd_g_const(&F_FTransform2_SSE2__k12), base.Simd_g_const(&F_FTransform2_SSE2__k13))
	v109 = int32(4)
	v111 = base.Simd_g_i16x8_sub(v102, v100)
	v115 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v111, v111, base.Simd_g_const(&F_FTransform2_SSE2__k12), base.Simd_g_const(&F_FTransform2_SSE2__k13)), v111, base.Simd_g_const(&F_FTransform2_SSE2__k21), base.Simd_g_const(&F_FTransform2_SSE2__k22))
	v116 = base.Simd_g_const(&F_FTransform2_SSE2__k31)
	v118 = base.Simd_g_const(&F_FTransform2_SSE2__k32)
	v120 = int32(16)
	v121 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v115, v116), v118), v120)
	v124 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_sub(v105, v107), v109), base.Simd_g_i16x8_narrow_i32x4_s(v121, v121), base.Simd_g_const(&F_FTransform2_SSE2__k9), base.Simd_g_const(&F_FTransform2_SSE2__k10))
	base.Simd_g_v128_store(m, l2, int32(48), v124)
	v130 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_shuffle2(v33, v42, base.Simd_g_const(&F_FTransform2_SSE2__k28), base.Simd_g_const(&F_FTransform2_SSE2__k29)), base.Simd_g_const(&F_FTransform2_SSE2__k5))
	v134 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_shuffle2(v55, v64, base.Simd_g_const(&F_FTransform2_SSE2__k28), base.Simd_g_const(&F_FTransform2_SSE2__k29)), base.Simd_g_const(&F_FTransform2_SSE2__k5))
	v136 = base.Simd_g_i8x16_shuffle2(v130, v134, base.Simd_g_const(&F_FTransform2_SSE2__k9), base.Simd_g_const(&F_FTransform2_SSE2__k10))
	v138 = base.Simd_g_i8x16_shuffle2(v130, v134, base.Simd_g_const(&F_FTransform2_SSE2__k12), base.Simd_g_const(&F_FTransform2_SSE2__k13))
	v139 = base.Simd_g_i16x8_add(v136, v138)
	v142 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_dot_i16x8_s(v139, v74), base.Simd_g_i32x4_dot_i16x8_s(v139, v76))
	v143 = base.Simd_g_i16x8_sub(v136, v138)
	v152 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v143, v80), v82), v84), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v143, v86), v88), v84))
	v154 = base.Simd_g_i8x16_shuffle2(v142, v152, base.Simd_g_const(&F_FTransform2_SSE2__k21), base.Simd_g_const(&F_FTransform2_SSE2__k22))
	v156 = base.Simd_g_i8x16_shuffle2(v142, v152, base.Simd_g_const(&F_FTransform2_SSE2__k24), base.Simd_g_const(&F_FTransform2_SSE2__k25))
	v160 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_shuffle2(v154, v156, base.Simd_g_const(&F_FTransform2_SSE2__k6), base.Simd_g_const(&F_FTransform2_SSE2__k7)), base.Simd_g_const(&F_FTransform2_SSE2__k26))
	v162 = base.Simd_g_i8x16_shuffle2(v154, v156, base.Simd_g_const(&F_FTransform2_SSE2__k28), base.Simd_g_const(&F_FTransform2_SSE2__k29))
	v163 = base.Simd_g_i16x8_add(v160, v162)
	v164 = base.Simd_g_i16x8_add(v163, v104)
	v166 = base.Simd_g_i8x16_shuffle2(v163, v163, base.Simd_g_const(&F_FTransform2_SSE2__k12), base.Simd_g_const(&F_FTransform2_SSE2__k13))
	v170 = base.Simd_g_i16x8_sub(v162, v160)
	v174 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v170, v170, base.Simd_g_const(&F_FTransform2_SSE2__k12), base.Simd_g_const(&F_FTransform2_SSE2__k13)), v170, base.Simd_g_const(&F_FTransform2_SSE2__k21), base.Simd_g_const(&F_FTransform2_SSE2__k22))
	v178 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v174, v116), v118), v120)
	v181 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_sub(v164, v166), v109), base.Simd_g_i16x8_narrow_i32x4_s(v178, v178), base.Simd_g_const(&F_FTransform2_SSE2__k9), base.Simd_g_const(&F_FTransform2_SSE2__k10))
	base.Simd_g_v128_store(m, l2, v120, v181)
	v187 = base.Simd_g_const(&F_FTransform2_SSE2__k33)
	v189 = base.Simd_g_const(&F_FTransform2_SSE2__k34)
	v192 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v115, v187), v189), v120)
	v197 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(v105, v107), v109), base.Simd_g_i16x8_add(base.Simd_g_i16x8_narrow_i32x4_s(v192, v192), base.Simd_g_i16x8_eq(v102, v100)), base.Simd_g_const(&F_FTransform2_SSE2__k9), base.Simd_g_const(&F_FTransform2_SSE2__k10))
	base.Simd_g_v128_store(m, l2, v34, v197)
	v206 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v174, v187), v189), v120)
	v211 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(v164, v166), v109), base.Simd_g_i16x8_add(base.Simd_g_i16x8_narrow_i32x4_s(v206, v206), base.Simd_g_i16x8_eq(v162, v160)), base.Simd_g_const(&F_FTransform2_SSE2__k9), base.Simd_g_const(&F_FTransform2_SSE2__k10))
	base.Simd_g_v128_store(m, l2, v24, v211)
	return
}

var F_FTransform2_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_FTransform2_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_FTransform2_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_FTransform2_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_FTransform2_SSE2__k4 = [2]uint64{0x1b1a19180b0a0908, 0x1f1e1d1c0f0e0d0c}
var F_FTransform2_SSE2__k5 = [2]uint64{0x706050403020100, 0xd0c0f0e09080b0a}
var F_FTransform2_SSE2__k6 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_FTransform2_SSE2__k7 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_FTransform2_SSE2__k8 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_FTransform2_SSE2__k9 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_FTransform2_SSE2__k10 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_FTransform2_SSE2__k11 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_FTransform2_SSE2__k12 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_FTransform2_SSE2__k13 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_FTransform2_SSE2__k14 = [2]uint64{0x8000800080008, 0x8000800080008}
var F_FTransform2_SSE2__k15 = [2]uint64{0xfff80008fff80008, 0xfff80008fff80008}
var F_FTransform2_SSE2__k16 = [2]uint64{0x8a914e808a914e8, 0x8a914e808a914e8}
var F_FTransform2_SSE2__k17 = [2]uint64{0x71400000714, 0x71400000714}
var F_FTransform2_SSE2__k18 = [2]uint64{0xeb1808a9eb1808a9, 0xeb1808a9eb1808a9}
var F_FTransform2_SSE2__k19 = [2]uint64{0x3a9000003a9, 0x3a9000003a9}
var F_FTransform2_SSE2__k20 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_FTransform2_SSE2__k21 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_FTransform2_SSE2__k22 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_FTransform2_SSE2__k23 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_FTransform2_SSE2__k24 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_FTransform2_SSE2__k25 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_FTransform2_SSE2__k26 = [2]uint64{0xf0e0d0c0b0a0908, 0x706050403020100}
var F_FTransform2_SSE2__k27 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_FTransform2_SSE2__k28 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_FTransform2_SSE2__k29 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_FTransform2_SSE2__k30 = [2]uint64{0x7000700070007, 0x7000700070007}
var F_FTransform2_SSE2__k31 = [2]uint64{0x8a9eb1808a9eb18, 0x8a9eb1808a9eb18}
var F_FTransform2_SSE2__k32 = [2]uint64{0xc7380000c738, 0xc7380000c738}
var F_FTransform2_SSE2__k33 = [2]uint64{0x14e808a914e808a9, 0x14e808a914e808a9}
var F_FTransform2_SSE2__k34 = [2]uint64{0x12ee000012ee0, 0x12ee000012ee0}

func F_FTransformWHT_SSE2(m *base.Module, l0 int32, l1 int32) {
	var v10 base.V128
	_ = v10
	var v12 base.V128
	_ = v12
	var v13 base.V128
	_ = v13
	var v14 base.V128
	_ = v14
	var v16 base.V128
	_ = v16
	var v18 base.V128
	_ = v18
	var v20 base.V128
	_ = v20
	var v21 base.V128
	_ = v21
	var v22 base.V128
	_ = v22
	var v23 base.V128
	_ = v23
	var v27 base.V128
	_ = v27
	var v29 base.V128
	_ = v29
	var v30 base.V128
	_ = v30
	var v32 base.V128
	_ = v32
	var v34 base.V128
	_ = v34
	var v36 base.V128
	_ = v36
	var v38 base.V128
	_ = v38
	var v40 base.V128
	_ = v40
	var v42 base.V128
	_ = v42
	var v43 base.V128
	_ = v43
	var v44 base.V128
	_ = v44
	var v51 base.V128
	_ = v51
	var v54 base.V128
	_ = v54
	var v56 base.V128
	_ = v56
	var v58 base.V128
	_ = v58
	var v60 base.V128
	_ = v60
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
	var v75 base.V128
	_ = v75
	var v76 int32
	_ = v76
	var v77 base.V128
	_ = v77
	var v79 base.V128
	_ = v79
	var v81 base.V128
	_ = v81
	var v83 base.V128
	_ = v83
	var v85 base.V128
	_ = v85
	var v87 base.V128
	_ = v87
	var v88 base.V128
	_ = v88
	var v89 base.V128
	_ = v89
	var v96 base.V128
	_ = v96
	var v99 base.V128
	_ = v99
	var v101 int32
	_ = v101
	var v102 base.V128
	_ = v102
	var v105 base.V128
	_ = v105
	var v111 base.V128
	_ = v111
	v10 = base.Simd_g_v128_load64_zero(m, l0, int32(384))
	v12 = base.Simd_g_v128_load64_zero(m, l0, int32(416))
	v13 = base.Simd_g_const(&F_FTransformWHT_SSE2__k0)
	v14 = base.Simd_g_i8x16_shuffle2(v10, v12, base.Simd_g_const(&F_FTransformWHT_SSE2__k1), base.Simd_g_const(&F_FTransformWHT_SSE2__k2))
	v16 = base.Simd_g_v128_load64_zero(m, l0, int32(448))
	v18 = base.Simd_g_v128_load64_zero(m, l0, int32(480))
	v20 = base.Simd_g_i8x16_shuffle2(v16, v18, base.Simd_g_const(&F_FTransformWHT_SSE2__k1), base.Simd_g_const(&F_FTransformWHT_SSE2__k2))
	v21 = base.Simd_g_i16x8_add_sat_s(v14, v20)
	v22 = base.Simd_g_i16x8_sub_sat_s(v14, v20)
	v23 = base.Simd_g_const(&F_FTransformWHT_SSE2__k3)
	v27 = base.Simd_g_const(&F_FTransformWHT_SSE2__k4)
	v29 = base.Simd_g_const(&F_FTransformWHT_SSE2__k5)
	v30 = base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v21, v22, base.Simd_g_const(&F_FTransformWHT_SSE2__k6), base.Simd_g_const(&F_FTransformWHT_SSE2__k7)), base.Simd_g_i8x16_shuffle2(v22, v21, base.Simd_g_const(&F_FTransformWHT_SSE2__k6), base.Simd_g_const(&F_FTransformWHT_SSE2__k7)), base.Simd_g_const(&F_FTransformWHT_SSE2__k8), base.Simd_g_const(&F_FTransformWHT_SSE2__k9)), v29)
	v32 = base.Simd_g_v128_load64_zero(m, l0, int32(128))
	v34 = base.Simd_g_v128_load64_zero(m, l0, int32(160))
	v36 = base.Simd_g_i8x16_shuffle2(v32, v34, base.Simd_g_const(&F_FTransformWHT_SSE2__k1), base.Simd_g_const(&F_FTransformWHT_SSE2__k2))
	v38 = base.Simd_g_v128_load64_zero(m, l0, int32(192))
	v40 = base.Simd_g_v128_load64_zero(m, l0, int32(224))
	v42 = base.Simd_g_i8x16_shuffle2(v38, v40, base.Simd_g_const(&F_FTransformWHT_SSE2__k1), base.Simd_g_const(&F_FTransformWHT_SSE2__k2))
	v43 = base.Simd_g_i16x8_add_sat_s(v36, v42)
	v44 = base.Simd_g_i16x8_sub_sat_s(v36, v42)
	v51 = base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v43, v44, base.Simd_g_const(&F_FTransformWHT_SSE2__k6), base.Simd_g_const(&F_FTransformWHT_SSE2__k7)), base.Simd_g_i8x16_shuffle2(v44, v43, base.Simd_g_const(&F_FTransformWHT_SSE2__k6), base.Simd_g_const(&F_FTransformWHT_SSE2__k7)), base.Simd_g_const(&F_FTransformWHT_SSE2__k8), base.Simd_g_const(&F_FTransformWHT_SSE2__k9)), v29)
	v54 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_add(v30, v51), base.Simd_g_i32x4_sub(v51, v30))
	v56 = base.Simd_g_v128_load64_zero(m, l0, int32(256))
	v58 = base.Simd_g_v128_load64_zero(m, l0, int32(288))
	v60 = base.Simd_g_i8x16_shuffle2(v56, v58, base.Simd_g_const(&F_FTransformWHT_SSE2__k1), base.Simd_g_const(&F_FTransformWHT_SSE2__k2))
	v62 = base.Simd_g_v128_load64_zero(m, l0, int32(320))
	v64 = base.Simd_g_v128_load64_zero(m, l0, int32(352))
	v66 = base.Simd_g_i8x16_shuffle2(v62, v64, base.Simd_g_const(&F_FTransformWHT_SSE2__k1), base.Simd_g_const(&F_FTransformWHT_SSE2__k2))
	v67 = base.Simd_g_i16x8_add_sat_s(v60, v66)
	v68 = base.Simd_g_i16x8_sub_sat_s(v60, v66)
	v75 = base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v67, v68, base.Simd_g_const(&F_FTransformWHT_SSE2__k6), base.Simd_g_const(&F_FTransformWHT_SSE2__k7)), base.Simd_g_i8x16_shuffle2(v68, v67, base.Simd_g_const(&F_FTransformWHT_SSE2__k6), base.Simd_g_const(&F_FTransformWHT_SSE2__k7)), base.Simd_g_const(&F_FTransformWHT_SSE2__k8), base.Simd_g_const(&F_FTransformWHT_SSE2__k9)), v29)
	v76 = int32(0)
	v77 = base.Simd_g_v128_load64_zero(m, l0, v76)
	v79 = base.Simd_g_v128_load64_zero(m, l0, int32(32))
	v81 = base.Simd_g_i8x16_shuffle2(v77, v79, base.Simd_g_const(&F_FTransformWHT_SSE2__k1), base.Simd_g_const(&F_FTransformWHT_SSE2__k2))
	v83 = base.Simd_g_v128_load64_zero(m, l0, int32(64))
	v85 = base.Simd_g_v128_load64_zero(m, l0, int32(96))
	v87 = base.Simd_g_i8x16_shuffle2(v83, v85, base.Simd_g_const(&F_FTransformWHT_SSE2__k1), base.Simd_g_const(&F_FTransformWHT_SSE2__k2))
	v88 = base.Simd_g_i16x8_add_sat_s(v81, v87)
	v89 = base.Simd_g_i16x8_sub_sat_s(v81, v87)
	v96 = base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v88, v89, base.Simd_g_const(&F_FTransformWHT_SSE2__k6), base.Simd_g_const(&F_FTransformWHT_SSE2__k7)), base.Simd_g_i8x16_shuffle2(v89, v88, base.Simd_g_const(&F_FTransformWHT_SSE2__k6), base.Simd_g_const(&F_FTransformWHT_SSE2__k7)), base.Simd_g_const(&F_FTransformWHT_SSE2__k8), base.Simd_g_const(&F_FTransformWHT_SSE2__k9)), v29)
	v99 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_add(v75, v96), base.Simd_g_i32x4_sub(v96, v75))
	v101 = int32(1)
	v102 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(v54, v99), v101)
	base.Simd_g_v128_store(m, l1, v76, v102)
	v105 = base.Simd_g_i16x8_sub(v99, v54)
	v111 = base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v105, v105, base.Simd_g_const(&F_FTransformWHT_SSE2__k10), base.Simd_g_const(&F_FTransformWHT_SSE2__k11)), v105, base.Simd_g_const(&F_FTransformWHT_SSE2__k8), base.Simd_g_const(&F_FTransformWHT_SSE2__k9)), v101)
	base.Simd_g_v128_store(m, l1, int32(16), v111)
	return
}

var F_FTransformWHT_SSE2__k0 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_FTransformWHT_SSE2__k1 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_FTransformWHT_SSE2__k2 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_FTransformWHT_SSE2__k3 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_FTransformWHT_SSE2__k4 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_FTransformWHT_SSE2__k5 = [2]uint64{0x1000100010001, 0xffff0001ffff0001}
var F_FTransformWHT_SSE2__k6 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_FTransformWHT_SSE2__k7 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_FTransformWHT_SSE2__k8 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_FTransformWHT_SSE2__k9 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_FTransformWHT_SSE2__k10 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_FTransformWHT_SSE2__k11 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}

func F_FTransform_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	var v4 base.V128
	_ = v4
	var v10 int32
	_ = v10
	var v11 base.V128
	_ = v11
	var v12 int32
	_ = v12
	var v13 base.V128
	_ = v13
	var v14 base.V128
	_ = v14
	var v17 base.V128
	_ = v17
	var v20 base.V128
	_ = v20
	var v22 base.V128
	_ = v22
	var v28 base.V128
	_ = v28
	var v29 base.V128
	_ = v29
	var v30 int32
	_ = v30
	var v31 base.V128
	_ = v31
	var v32 int32
	_ = v32
	var v33 base.V128
	_ = v33
	var v39 base.V128
	_ = v39
	var v41 base.V128
	_ = v41
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
	var v58 base.V128
	_ = v58
	var v59 base.V128
	_ = v59
	var v64 int32
	_ = v64
	var v72 base.V128
	_ = v72
	var v74 base.V128
	_ = v74
	var v76 base.V128
	_ = v76
	var v80 base.V128
	_ = v80
	var v82 base.V128
	_ = v82
	var v83 base.V128
	_ = v83
	var v85 base.V128
	_ = v85
	var v87 base.V128
	_ = v87
	var v89 int32
	_ = v89
	var v91 base.V128
	_ = v91
	var v95 base.V128
	_ = v95
	var v100 int32
	_ = v100
	var v101 base.V128
	_ = v101
	var v104 base.V128
	_ = v104
	var v115 base.V128
	_ = v115
	var v120 base.V128
	_ = v120
	v4 = base.Simd_g_const(&F_FTransform_SSE2__k0)
	v10 = int32(0)
	v11 = base.Simd_g_v128_load64_zero(m, l0, v10)
	v12 = int32(32)
	v13 = base.Simd_g_v128_load64_zero(m, l0, v12)
	v14 = base.Simd_g_const(&F_FTransform_SSE2__k1)
	v17 = base.Simd_g_const(&F_FTransform_SSE2__k2)
	v20 = base.Simd_g_v128_load64_zero(m, l1, v10)
	v22 = base.Simd_g_v128_load64_zero(m, l1, v12)
	v28 = base.Simd_g_const(&F_FTransform_SSE2__k3)
	v29 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_sub(base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v11, v13, base.Simd_g_const(&F_FTransform_SSE2__k4), base.Simd_g_const(&F_FTransform_SSE2__k5)), v4, base.Simd_g_const(&F_FTransform_SSE2__k6), base.Simd_g_const(&F_FTransform_SSE2__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v20, v22, base.Simd_g_const(&F_FTransform_SSE2__k4), base.Simd_g_const(&F_FTransform_SSE2__k5)), v4, base.Simd_g_const(&F_FTransform_SSE2__k6), base.Simd_g_const(&F_FTransform_SSE2__k7))), base.Simd_g_const(&F_FTransform_SSE2__k3))
	v30 = int32(64)
	v31 = base.Simd_g_v128_load64_zero(m, l0, v30)
	v32 = int32(96)
	v33 = base.Simd_g_v128_load64_zero(m, l0, v32)
	v39 = base.Simd_g_v128_load64_zero(m, l1, v30)
	v41 = base.Simd_g_v128_load64_zero(m, l1, v32)
	v48 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_sub(base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v31, v33, base.Simd_g_const(&F_FTransform_SSE2__k4), base.Simd_g_const(&F_FTransform_SSE2__k5)), v4, base.Simd_g_const(&F_FTransform_SSE2__k6), base.Simd_g_const(&F_FTransform_SSE2__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v39, v41, base.Simd_g_const(&F_FTransform_SSE2__k4), base.Simd_g_const(&F_FTransform_SSE2__k5)), v4, base.Simd_g_const(&F_FTransform_SSE2__k6), base.Simd_g_const(&F_FTransform_SSE2__k7))), base.Simd_g_const(&F_FTransform_SSE2__k3))
	v49 = base.Simd_g_const(&F_FTransform_SSE2__k8)
	v50 = base.Simd_g_i8x16_shuffle2(v29, v48, base.Simd_g_const(&F_FTransform_SSE2__k9), base.Simd_g_const(&F_FTransform_SSE2__k10))
	v51 = base.Simd_g_const(&F_FTransform_SSE2__k11)
	v52 = base.Simd_g_i8x16_shuffle2(v29, v48, base.Simd_g_const(&F_FTransform_SSE2__k12), base.Simd_g_const(&F_FTransform_SSE2__k13))
	v53 = base.Simd_g_i16x8_add(v50, v52)
	v58 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_dot_i16x8_s(v53, base.Simd_g_const(&F_FTransform_SSE2__k14)), base.Simd_g_i32x4_dot_i16x8_s(v53, base.Simd_g_const(&F_FTransform_SSE2__k15)))
	v59 = base.Simd_g_i16x8_sub(v50, v52)
	v64 = int32(9)
	v72 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v59, base.Simd_g_const(&F_FTransform_SSE2__k16)), base.Simd_g_const(&F_FTransform_SSE2__k17)), v64), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v59, base.Simd_g_const(&F_FTransform_SSE2__k18)), base.Simd_g_const(&F_FTransform_SSE2__k19)), v64))
	v74 = base.Simd_g_i8x16_shuffle2(v58, v72, base.Simd_g_const(&F_FTransform_SSE2__k4), base.Simd_g_const(&F_FTransform_SSE2__k5))
	v76 = base.Simd_g_i8x16_shuffle2(v58, v72, base.Simd_g_const(&F_FTransform_SSE2__k20), base.Simd_g_const(&F_FTransform_SSE2__k21))
	v80 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_shuffle2(v74, v76, base.Simd_g_const(&F_FTransform_SSE2__k22), base.Simd_g_const(&F_FTransform_SSE2__k23)), base.Simd_g_const(&F_FTransform_SSE2__k24))
	v82 = base.Simd_g_i8x16_shuffle2(v74, v76, base.Simd_g_const(&F_FTransform_SSE2__k25), base.Simd_g_const(&F_FTransform_SSE2__k26))
	v83 = base.Simd_g_i16x8_add(v80, v82)
	v85 = base.Simd_g_i16x8_add(v83, base.Simd_g_const(&F_FTransform_SSE2__k27))
	v87 = base.Simd_g_i8x16_shuffle2(v83, v83, base.Simd_g_const(&F_FTransform_SSE2__k12), base.Simd_g_const(&F_FTransform_SSE2__k13))
	v89 = int32(4)
	v91 = base.Simd_g_i16x8_sub(v82, v80)
	v95 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v91, v91, base.Simd_g_const(&F_FTransform_SSE2__k12), base.Simd_g_const(&F_FTransform_SSE2__k13)), v91, base.Simd_g_const(&F_FTransform_SSE2__k4), base.Simd_g_const(&F_FTransform_SSE2__k5))
	v100 = int32(16)
	v101 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v95, base.Simd_g_const(&F_FTransform_SSE2__k28)), base.Simd_g_const(&F_FTransform_SSE2__k29)), v100)
	v104 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_sub(v85, v87), v89), base.Simd_g_i16x8_narrow_i32x4_s(v101, v101), base.Simd_g_const(&F_FTransform_SSE2__k9), base.Simd_g_const(&F_FTransform_SSE2__k10))
	base.Simd_g_v128_store(m, l2, v100, v104)
	v115 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v95, base.Simd_g_const(&F_FTransform_SSE2__k30)), base.Simd_g_const(&F_FTransform_SSE2__k31)), v100)
	v120 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(v85, v87), v89), base.Simd_g_i16x8_add(base.Simd_g_i16x8_narrow_i32x4_s(v115, v115), base.Simd_g_i16x8_eq(v82, v80)), base.Simd_g_const(&F_FTransform_SSE2__k9), base.Simd_g_const(&F_FTransform_SSE2__k10))
	base.Simd_g_v128_store(m, l2, v10, v120)
	return
}

var F_FTransform_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_FTransform_SSE2__k1 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_FTransform_SSE2__k2 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_FTransform_SSE2__k3 = [2]uint64{0x706050403020100, 0xd0c0f0e09080b0a}
var F_FTransform_SSE2__k4 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_FTransform_SSE2__k5 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_FTransform_SSE2__k6 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_FTransform_SSE2__k7 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_FTransform_SSE2__k8 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_FTransform_SSE2__k9 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_FTransform_SSE2__k10 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_FTransform_SSE2__k11 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_FTransform_SSE2__k12 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_FTransform_SSE2__k13 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_FTransform_SSE2__k14 = [2]uint64{0x8000800080008, 0x8000800080008}
var F_FTransform_SSE2__k15 = [2]uint64{0xfff80008fff80008, 0xfff80008fff80008}
var F_FTransform_SSE2__k16 = [2]uint64{0x8a914e808a914e8, 0x8a914e808a914e8}
var F_FTransform_SSE2__k17 = [2]uint64{0x71400000714, 0x71400000714}
var F_FTransform_SSE2__k18 = [2]uint64{0xeb1808a9eb1808a9, 0xeb1808a9eb1808a9}
var F_FTransform_SSE2__k19 = [2]uint64{0x3a9000003a9, 0x3a9000003a9}
var F_FTransform_SSE2__k20 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_FTransform_SSE2__k21 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_FTransform_SSE2__k22 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_FTransform_SSE2__k23 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_FTransform_SSE2__k24 = [2]uint64{0xf0e0d0c0b0a0908, 0x706050403020100}
var F_FTransform_SSE2__k25 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_FTransform_SSE2__k26 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_FTransform_SSE2__k27 = [2]uint64{0x7000700070007, 0x7000700070007}
var F_FTransform_SSE2__k28 = [2]uint64{0x8a9eb1808a9eb18, 0x8a9eb1808a9eb18}
var F_FTransform_SSE2__k29 = [2]uint64{0xc7380000c738, 0xc7380000c738}
var F_FTransform_SSE2__k30 = [2]uint64{0x14e808a914e808a9, 0x14e808a914e808a9}
var F_FTransform_SSE2__k31 = [2]uint64{0x12ee000012ee0, 0x12ee000012ee0}
