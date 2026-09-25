//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_GetCombinedEntropyUnrefined_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v6 int64
	_ = v6
	var v19 base.V128
	_ = v19
	var v25 int32
	_ = v25
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v39 int32
	_ = v39
	var v43 int32
	_ = v43
	var v51 int32
	_ = v51
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v61 int64
	_ = v61
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v67 int32
	_ = v67
	var v69 int64
	_ = v69
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v76 int32
	_ = v76
	var v78 int32
	_ = v78
	var v83 int32
	_ = v83
	var v86 int32
	_ = v86
	var v90 int32
	_ = v90
	var v96 int64
	_ = v96
	var v97 int32
	_ = v97
	var v100 int32
	_ = v100
	var v101 int64
	_ = v101
	var v102 int64
	_ = v102
	var v106 int64
	_ = v106
	var v110 int64
	_ = v110
	var v111 int32
	_ = v111
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v116 int32
	_ = v116
	var v117 int32
	_ = v117
	var v119 int32
	_ = v119
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v122 int32
	_ = v122
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v135 int64
	_ = v135
	var v136 int32
	_ = v136
	var v137 int32
	_ = v137
	var v138 int32
	_ = v138
	var v139 int32
	_ = v139
	var v140 int32
	_ = v140
	var v141 int32
	_ = v141
	var v147 int32
	_ = v147
	var v152 int32
	_ = v152
	var v156 int32
	_ = v156
	var v159 int64
	_ = v159
	var v160 int32
	_ = v160
	var v161 int32
	_ = v161
	var v162 int32
	_ = v162
	var v163 int32
	_ = v163
	var v164 int32
	_ = v164
	var v172 int32
	_ = v172
	var v179 int32
	_ = v179
	var v183 int32
	_ = v183
	var v189 int64
	_ = v189
	var v190 int32
	_ = v190
	var v193 int32
	_ = v193
	var v194 int64
	_ = v194
	var v195 int64
	_ = v195
	var v198 int64
	_ = v198
	var v201 int64
	_ = v201
	var v202 int32
	_ = v202
	var v205 int32
	_ = v205
	var v206 int32
	_ = v206
	var v208 int32
	_ = v208
	var v209 int32
	_ = v209
	var v210 int32
	_ = v210
	var v211 int32
	_ = v211
	var v221 int32
	_ = v221
	var v222 int32
	_ = v222
	var v227 int32
	_ = v227
	var v233 int64
	_ = v233
	var v234 int32
	_ = v234
	var v237 int32
	_ = v237
	var v238 int64
	_ = v238
	var v239 int64
	_ = v239
	v6 = int64(0)
	v19 = base.Simd_g_const(&F_GetCombinedEntropyUnrefined_C__k0)
	base.Simd_g_v128_store(m, l4, int32(0), v19)
	*(*int64)(unsafe.Add(mBase, uint32(l3)+8)) = v6
	v25 = int32(16)
	*(*int64)(unsafe.Add(mBase, uint32(l4+v25))) = v6
	*(*int64)(unsafe.Add(mBase, uint32(l3+v25))) = int64(281470681743360)
	v33 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	v34 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v35 = v33 + v34
	if int32(2) <= l2 {
		v43 = int32(4)
		v51 = int32(0)
		v56 = int32(1)
		v57 = l1 + v43
		v61 = int64(0)
		v62 = v35
		v63 = v51
		v64 = v51
		v65 = v51
		v66 = v51
		v67 = l0 + v43
		v69 = int64(1)
		for {
			v74 = *(*int32)(unsafe.Add(mBase, uint32(v57)))
			v75 = *(*int32)(unsafe.Add(mBase, uint32(v67)))
			v76 = v74 + v75
			if v76 != v62 {
				v78 = v56 - v66
				if v62 == int32(0) {
					v110 = v61
					v111 = v63
					v112 = v64
					v113 = v65
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l3)+20)) = v66
					v83 = v64 + v56 - v66
					*(*int32)(unsafe.Add(mBase, uint32(l3)+12)) = v83
					v86 = v65 + v78*v62
					*(*int32)(unsafe.Add(mBase, uint32(l3)+8)) = v86
					if base.Ui32(int32(255)) < base.Ui32(v62) {
						v97 = m.G1
						v100 = *(*int32)(unsafe.Add(mBase, uint32(v97)+uint32(_c_F_GetCombinedEntropyUnrefined_C[0])))
						v101 = m.T0[v100].(func(*base.Module, int32) int64)(m, v62)
						mBase = m.M
						v102 = v101
					} else {
						v90 = m.G1
						v96 = *(*int64)(unsafe.Add(mBase, uint32(v90+int32(_a_F_GetCombinedEntropyUnrefined_C_0)+v62<<(uint(int32(3))%32))))
						v102 = v96
					}
					v106 = v61 + v102*(v69-base.I64_extend_i32_s(v66))
					*(*int64)(unsafe.Add(mBase, uint32(l3))) = v106
					if base.Ui32(v62) <= base.Ui32(v63) {
						v110 = v106
						v111 = v63
						v112 = v83
						v113 = v86
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(l3)+16)) = v62
						v110 = v106
						v111 = v62
						v112 = v83
						v113 = v86
					}
				}
				v116 = base.B2i32(v62 != int32(0))
				v117 = int32(2)
				v119 = l4 + v116<<(uint(v117)%32)
				v120 = *(*int32)(unsafe.Add(mBase, uint32(v119)))
				v121 = int32(3)
				v122 = base.B2i32(v121 < v78)
				*(*int32)(unsafe.Add(mBase, uint32(v119))) = v120 + v122
				v130 = l4 + int32(8) + v116<<(uint(v121)%32) + v122<<(uint(v117)%32)
				v131 = *(*int32)(unsafe.Add(mBase, uint32(v130)))
				*(*int32)(unsafe.Add(mBase, uint32(v130))) = v56 + (v131 - v66)
				v135 = v110
				v136 = v76
				v137 = v111
				v138 = v112
				v139 = v113
				v140 = v56
				v141 = v76
			} else {
				v135 = v61
				v136 = v62
				v137 = v63
				v138 = v64
				v139 = v65
				v140 = v66
				v141 = v62
			}
			v147 = int32(4)
			v152 = v56 + int32(1)
			if l2 != v152 {
				v56 = v152
				v57 = v57 + v147
				v61 = v135
				v62 = v136
				v63 = v137
				v64 = v138
				v65 = v139
				v66 = v140
				v67 = v67 + v147
				v69 = v69 + int64(1)
				continue
			} else {
				break
			}
			break
		}
		v156 = l2
		v159 = v135
		v160 = v141
		v161 = v137
		v162 = v138
		v163 = v139
		v164 = v140
	} else {
		v39 = int32(0)
		v156 = int32(1)
		v159 = v6
		v160 = v35
		v161 = v39
		v162 = v39
		v163 = v39
		v164 = v39
	}
	v172 = v156 - v164
	if v160 == int32(0) {
		v201 = v159
		v202 = v163
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l3)+20)) = v164
		*(*int32)(unsafe.Add(mBase, uint32(l3)+12)) = v162 + v172
		v179 = v163 + v172*v160
		*(*int32)(unsafe.Add(mBase, uint32(l3)+8)) = v179
		if base.Ui32(int32(255)) < base.Ui32(v160) {
			v190 = m.G1
			v193 = *(*int32)(unsafe.Add(mBase, uint32(v190)+uint32(_c_F_GetCombinedEntropyUnrefined_C[0])))
			v194 = m.T0[v193].(func(*base.Module, int32) int64)(m, v160)
			mBase = m.M
			v195 = v194
		} else {
			v183 = m.G1
			v189 = *(*int64)(unsafe.Add(mBase, uint32(v183+int32(_a_F_GetCombinedEntropyUnrefined_C_0)+v160<<(uint(int32(3))%32))))
			v195 = v189
		}
		v198 = v159 + v195*base.I64_extend_i32_s(v172)
		if base.Ui32(v160) <= base.Ui32(v161) {
			v201 = v198
			v202 = v179
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l3)+16)) = v160
			v201 = v198
			v202 = v179
		}
	}
	v205 = base.B2i32(v160 != int32(0))
	v206 = int32(2)
	v208 = l4 + v205<<(uint(v206)%32)
	v209 = *(*int32)(unsafe.Add(mBase, uint32(v208)))
	v210 = int32(3)
	v211 = base.B2i32(v210 < v172)
	*(*int32)(unsafe.Add(mBase, uint32(v208))) = v209 + v211
	v221 = l4 + v205<<(uint(v210)%32) + v211<<(uint(v206)%32) + int32(8)
	v222 = *(*int32)(unsafe.Add(mBase, uint32(v221)))
	*(*int32)(unsafe.Add(mBase, uint32(v221))) = v222 + v172
	if base.Ui32(int32(255)) < base.Ui32(v202) {
		v234 = m.G1
		v237 = *(*int32)(unsafe.Add(mBase, uint32(v234)+uint32(_c_F_GetCombinedEntropyUnrefined_C[0])))
		v238 = m.T0[v237].(func(*base.Module, int32) int64)(m, v202)
		mBase = m.M
		v239 = v238
	} else {
		v227 = m.G1
		v233 = *(*int64)(unsafe.Add(mBase, uint32(v227+int32(_a_F_GetCombinedEntropyUnrefined_C_0)+v202<<(uint(int32(3))%32))))
		v239 = v233
	}
	*(*int64)(unsafe.Add(mBase, uint32(l3))) = v239 - v201
	return
}

var F_GetCombinedEntropyUnrefined_C__k0 = [2]uint64{0x0, 0x0}

func F_GetEntropyUnrefined_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v5 int64
	_ = v5
	var v18 base.V128
	_ = v18
	var v24 int32
	_ = v24
	var v32 int32
	_ = v32
	var v36 int32
	_ = v36
	var v46 int32
	_ = v46
	var v51 int32
	_ = v51
	var v55 int64
	_ = v55
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v63 int64
	_ = v63
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v75 int32
	_ = v75
	var v78 int32
	_ = v78
	var v82 int32
	_ = v82
	var v88 int64
	_ = v88
	var v89 int32
	_ = v89
	var v92 int32
	_ = v92
	var v93 int64
	_ = v93
	var v94 int64
	_ = v94
	var v98 int64
	_ = v98
	var v102 int64
	_ = v102
	var v103 int32
	_ = v103
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v111 int32
	_ = v111
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v127 int64
	_ = v127
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v132 int32
	_ = v132
	var v133 int32
	_ = v133
	var v142 int32
	_ = v142
	var v145 int32
	_ = v145
	var v148 int64
	_ = v148
	var v149 int32
	_ = v149
	var v150 int32
	_ = v150
	var v151 int32
	_ = v151
	var v152 int32
	_ = v152
	var v153 int32
	_ = v153
	var v161 int32
	_ = v161
	var v168 int32
	_ = v168
	var v172 int32
	_ = v172
	var v178 int64
	_ = v178
	var v179 int32
	_ = v179
	var v182 int32
	_ = v182
	var v183 int64
	_ = v183
	var v184 int64
	_ = v184
	var v187 int64
	_ = v187
	var v190 int64
	_ = v190
	var v191 int32
	_ = v191
	var v194 int32
	_ = v194
	var v195 int32
	_ = v195
	var v197 int32
	_ = v197
	var v198 int32
	_ = v198
	var v199 int32
	_ = v199
	var v200 int32
	_ = v200
	var v210 int32
	_ = v210
	var v211 int32
	_ = v211
	var v216 int32
	_ = v216
	var v222 int64
	_ = v222
	var v223 int32
	_ = v223
	var v226 int32
	_ = v226
	var v227 int64
	_ = v227
	var v228 int64
	_ = v228
	v5 = int64(0)
	v18 = base.Simd_g_const(&F_GetEntropyUnrefined_C__k0)
	base.Simd_g_v128_store(m, l3, int32(0), v18)
	*(*int64)(unsafe.Add(mBase, uint32(l2)+8)) = v5
	v24 = int32(16)
	*(*int64)(unsafe.Add(mBase, uint32(l3+v24))) = v5
	*(*int64)(unsafe.Add(mBase, uint32(l2+v24))) = int64(281470681743360)
	v32 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if int32(2) <= l1 {
		v46 = int32(0)
		v51 = int32(1)
		v55 = int64(0)
		v56 = v32
		v57 = v46
		v58 = v46
		v59 = v46
		v60 = v46
		v61 = l0 + int32(4)
		v63 = int64(1)
		for {
			v68 = *(*int32)(unsafe.Add(mBase, uint32(v61)))
			if v68 != v56 {
				v70 = v51 - v60
				if v56 == int32(0) {
					v102 = v55
					v103 = v57
					v104 = v58
					v105 = v59
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = v60
					v75 = v58 + v51 - v60
					*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v75
					v78 = v59 + v70*v56
					*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v78
					if base.Ui32(int32(255)) < base.Ui32(v56) {
						v89 = m.G1
						v92 = *(*int32)(unsafe.Add(mBase, uint32(v89)+uint32(_c_F_GetEntropyUnrefined_C[0])))
						v93 = m.T0[v92].(func(*base.Module, int32) int64)(m, v56)
						mBase = m.M
						v94 = v93
					} else {
						v82 = m.G1
						v88 = *(*int64)(unsafe.Add(mBase, uint32(v82+int32(_a_F_GetEntropyUnrefined_C_0)+v56<<(uint(int32(3))%32))))
						v94 = v88
					}
					v98 = v55 + v94*(v63-base.I64_extend_i32_s(v60))
					*(*int64)(unsafe.Add(mBase, uint32(l2))) = v98
					if base.Ui32(v56) <= base.Ui32(v57) {
						v102 = v98
						v103 = v57
						v104 = v75
						v105 = v78
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(l2)+16)) = v56
						v102 = v98
						v103 = v56
						v104 = v75
						v105 = v78
					}
				}
				v108 = base.B2i32(v56 != int32(0))
				v109 = int32(2)
				v111 = l3 + v108<<(uint(v109)%32)
				v112 = *(*int32)(unsafe.Add(mBase, uint32(v111)))
				v113 = int32(3)
				v114 = base.B2i32(v113 < v70)
				*(*int32)(unsafe.Add(mBase, uint32(v111))) = v112 + v114
				v122 = l3 + int32(8) + v108<<(uint(v113)%32) + v114<<(uint(v109)%32)
				v123 = *(*int32)(unsafe.Add(mBase, uint32(v122)))
				*(*int32)(unsafe.Add(mBase, uint32(v122))) = v51 + (v123 - v60)
				v127 = v102
				v128 = v68
				v129 = v103
				v130 = v104
				v131 = v105
				v132 = v51
				v133 = v68
			} else {
				v127 = v55
				v128 = v56
				v129 = v57
				v130 = v58
				v131 = v59
				v132 = v60
				v133 = v56
			}
			v142 = v51 + int32(1)
			if l1 != v142 {
				v51 = v142
				v55 = v127
				v56 = v128
				v57 = v129
				v58 = v130
				v59 = v131
				v60 = v132
				v61 = v61 + int32(4)
				v63 = v63 + int64(1)
				continue
			} else {
				break
			}
			break
		}
		v145 = l1
		v148 = v127
		v149 = v133
		v150 = v129
		v151 = v130
		v152 = v131
		v153 = v132
	} else {
		v36 = int32(0)
		v145 = int32(1)
		v148 = v5
		v149 = v32
		v150 = v36
		v151 = v36
		v152 = v36
		v153 = v36
	}
	v161 = v145 - v153
	if v149 == int32(0) {
		v190 = v148
		v191 = v152
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = v153
		*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v151 + v161
		v168 = v152 + v161*v149
		*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v168
		if base.Ui32(int32(255)) < base.Ui32(v149) {
			v179 = m.G1
			v182 = *(*int32)(unsafe.Add(mBase, uint32(v179)+uint32(_c_F_GetEntropyUnrefined_C[0])))
			v183 = m.T0[v182].(func(*base.Module, int32) int64)(m, v149)
			mBase = m.M
			v184 = v183
		} else {
			v172 = m.G1
			v178 = *(*int64)(unsafe.Add(mBase, uint32(v172+int32(_a_F_GetEntropyUnrefined_C_0)+v149<<(uint(int32(3))%32))))
			v184 = v178
		}
		v187 = v148 + v184*base.I64_extend_i32_s(v161)
		if base.Ui32(v149) <= base.Ui32(v150) {
			v190 = v187
			v191 = v168
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l2)+16)) = v149
			v190 = v187
			v191 = v168
		}
	}
	v194 = base.B2i32(v149 != int32(0))
	v195 = int32(2)
	v197 = l3 + v194<<(uint(v195)%32)
	v198 = *(*int32)(unsafe.Add(mBase, uint32(v197)))
	v199 = int32(3)
	v200 = base.B2i32(v199 < v161)
	*(*int32)(unsafe.Add(mBase, uint32(v197))) = v198 + v200
	v210 = l3 + v194<<(uint(v199)%32) + v200<<(uint(v195)%32) + int32(8)
	v211 = *(*int32)(unsafe.Add(mBase, uint32(v210)))
	*(*int32)(unsafe.Add(mBase, uint32(v210))) = v211 + v161
	if base.Ui32(int32(255)) < base.Ui32(v191) {
		v223 = m.G1
		v226 = *(*int32)(unsafe.Add(mBase, uint32(v223)+uint32(_c_F_GetEntropyUnrefined_C[0])))
		v227 = m.T0[v226].(func(*base.Module, int32) int64)(m, v191)
		mBase = m.M
		v228 = v227
	} else {
		v216 = m.G1
		v222 = *(*int64)(unsafe.Add(mBase, uint32(v216+int32(_a_F_GetEntropyUnrefined_C_0)+v191<<(uint(int32(3))%32))))
		v228 = v222
	}
	*(*int64)(unsafe.Add(mBase, uint32(l2))) = v228 - v190
	return
}

var F_GetEntropyUnrefined_C__k0 = [2]uint64{0x0, 0x0}

func F_GetResidualCost_SSE2(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v41 int32
	_ = v41
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	var v46 int32
	_ = v46
	var v50 int32
	_ = v50
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v54 base.V128
	_ = v54
	var v55 base.V128
	_ = v55
	var v58 int32
	_ = v58
	var v59 base.V128
	_ = v59
	var v60 base.V128
	_ = v60
	var v63 base.V128
	_ = v63
	var v65 base.V128
	_ = v65
	var v69 base.V128
	_ = v69
	var v74 int32
	_ = v74
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v95 int32
	_ = v95
	var v96 int32
	_ = v96
	var v100 int32
	_ = v100
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v103 int32
	_ = v103
	var v106 int32
	_ = v106
	var v108 int32
	_ = v108
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v116 int32
	_ = v116
	var v118 int32
	_ = v118
	var v122 int32
	_ = v122
	var v128 int32
	_ = v128
	var v133 int32
	_ = v133
	var v134 int32
	_ = v134
	var v137 int32
	_ = v137
	var v142 int32
	_ = v142
	var v145 int32
	_ = v145
	var v146 int32
	_ = v146
	var v150 int32
	_ = v150
	var v155 int32
	_ = v155
	var v159 int32
	_ = v159
	var v160 int32
	_ = v160
	var v166 int32
	_ = v166
	var v167 int32
	_ = v167
	var v169 int32
	_ = v169
	var v171 int32
	_ = v171
	var v178 int32
	_ = v178
	var v179 int32
	_ = v179
	var v183 int32
	_ = v183
	var v186 int32
	_ = v186
	v14 = m.G0
	v16 = v14 - int32(64)
	v17 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v18 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	v25 = *(*int32)(unsafe.Add(mBase, uint32(v17+v18*int32(12)+l0<<(uint(int32(2))%32))))
	v26 = *(*int32)(unsafe.Add(mBase, uint32(l1)+16))
	v33 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v26+v18*int32(33)+l0*int32(11)))))
	if l0 != 0 {
		v42 = int32(0)
	} else {
		v35 = m.G79
		v41 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v35+(v33^int32(255))<<(uint(int32(1))%32)))))
		v42 = v41
	}
	v43 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
	if int32(-1) < v43 {
		v52 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
		v53 = int32(0)
		v54 = base.Simd_g_v128_load(m, v52, v53)
		v55 = base.Simd_g_i16x8_abs(v54)
		base.Simd_g_v128_store(m, v16, v53, v55)
		v58 = int32(16)
		v59 = base.Simd_g_v128_load(m, v52, v58)
		v60 = base.Simd_g_i16x8_abs(v59)
		base.Simd_g_v128_store(m, v16, v58, v60)
		v63 = base.Simd_g_i8x16_narrow_i16x8_s(v55, v60)
		v65 = base.Simd_g_i8x16_min_u(v63, base.Simd_g_const(&F_GetResidualCost_SSE2__k0))
		base.Simd_g_v128_store(m, v16, int32(32), v65)
		v69 = base.Simd_g_i8x16_min_u(v63, base.Simd_g_const(&F_GetResidualCost_SSE2__k1))
		base.Simd_g_v128_store(m, v16, int32(48), v69)
		if v43 <= v18 {
			v133 = v18
			v134 = v25
			v137 = v42
		} else {
			v74 = int32(12)
			v88 = v16 + v18<<(uint(int32(1))%32)
			v89 = v18*v74 + v17 + v74
			v92 = v16 + int32(48) + v18
			v93 = v25
			v95 = v16 + int32(32) + v18
			v96 = v42
			v100 = v43 - v18
			for {
				v101 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v88))))
				v102 = m.G80
				v103 = int32(1)
				v106 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v102+v101<<(uint(v103)%32)))))
				v108 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v92))))
				v112 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v93+v108<<(uint(v103)%32)))))
				v113 = v96 + v106 + v112
				v116 = int32(2)
				v118 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v95))))
				v122 = *(*int32)(unsafe.Add(mBase, uint32(v89+v118<<(uint(v116)%32))))
				v128 = v100 + int32(-1)
				if v128 != 0 {
					v88 = v88 + v116
					v89 = v89 + int32(12)
					v92 = v92 + v103
					v93 = v122
					v95 = v95 + v103
					v96 = v113
					v100 = v128
					continue
				} else {
					break
				}
				break
			}
			v133 = v43
			v134 = v122
			v137 = v113
		}
		v142 = int32(1)
		v145 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v16+v133<<(uint(v142)%32)))))
		v146 = m.G80
		v150 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v146+v145<<(uint(v142)%32)))))
		v155 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v16+int32(48)+v133))))
		v159 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v134+v155<<(uint(v142)%32)))))
		v160 = v137 + v150 + v159
		if int32(14) < v133 {
			v186 = v160
		} else {
			v166 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v16+int32(32)+v133))))
			v167 = m.G81
			v169 = int32(1)
			v171 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v133+v167+v169))))
			v178 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v26+v171*int32(33)+v166*int32(11)))))
			v179 = m.G79
			v183 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v179+v178<<(uint(v169)%32)))))
			v186 = v160 + v183
		}
		return v186
	} else {
		v46 = m.G79
		v50 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v46+v33<<(uint(int32(1))%32)))))
		return v50
	}
}

var F_GetResidualCost_SSE2__k0 = [2]uint64{0x202020202020202, 0x202020202020202}
var F_GetResidualCost_SSE2__k1 = [2]uint64{0x4343434343434343, 0x4343434343434343}

func F_GradientFilter_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v27 int32
	_ = v27
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v47 int32
	_ = v47
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v68 int32
	_ = v68
	var v69 base.V128
	_ = v69
	var v71 base.V128
	_ = v71
	var v72 base.V128
	_ = v72
	var v75 int32
	_ = v75
	var v80 int32
	_ = v80
	var v88 int32
	_ = v88
	var v122 int32
	_ = v122
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v145 int32
	_ = v145
	var v150 int32
	_ = v150
	var v179 int32
	_ = v179
	var v181 int32
	_ = v181
	var v182 int32
	_ = v182
	var v184 int32
	_ = v184
	var v186 int32
	_ = v186
	var v192 int32
	_ = v192
	var v196 int32
	_ = v196
	var v197 int32
	_ = v197
	var v203 int32
	_ = v203
	var v209 int32
	_ = v209
	var v218 int32
	_ = v218
	var v219 int32
	_ = v219
	var v220 int32
	_ = v220
	var v221 int32
	_ = v221
	var v230 int32
	_ = v230
	var v232 int32
	_ = v232
	var v235 int32
	_ = v235
	var v246 base.V128
	_ = v246
	var v253 int32
	_ = v253
	var v254 base.V128
	_ = v254
	var v258 base.V128
	_ = v258
	var v267 base.V128
	_ = v267
	var v271 base.V128
	_ = v271
	var v273 base.V128
	_ = v273
	var v275 base.V128
	_ = v275
	var v292 base.V128
	_ = v292
	var v308 base.V128
	_ = v308
	var v326 base.V128
	_ = v326
	var v329 int32
	_ = v329
	var v334 int32
	_ = v334
	var v343 int32
	_ = v343
	var v345 int32
	_ = v345
	var v372 int32
	_ = v372
	var v374 int32
	_ = v374
	var v376 int32
	_ = v376
	var v377 int32
	_ = v377
	var v397 int32
	_ = v397
	var v399 int32
	_ = v399
	var v400 int32
	_ = v400
	var v402 int32
	_ = v402
	var v403 int32
	_ = v403
	var v404 int32
	_ = v404
	var v407 int32
	_ = v407
	var v408 int32
	_ = v408
	var v411 int32
	_ = v411
	var v412 int32
	_ = v412
	var v417 int32
	_ = v417
	var v446 int32
	_ = v446
	var v451 int32
	_ = v451
	var v452 int32
	_ = v452
	var v453 int32
	_ = v453
	var v463 int32
	_ = v463
	var v466 int32
	_ = v466
	var v469 int32
	_ = v469
	var v480 int32
	_ = v480
	var v484 int32
	_ = v484
	var v485 int32
	_ = v485
	var v503 int32
	_ = v503
	var v504 int32
	_ = v504
	var v508 int32
	_ = v508
	var v509 int32
	_ = v509
	var v513 int32
	_ = v513
	var v514 int32
	_ = v514
	var v518 int32
	_ = v518
	var v519 int32
	_ = v519
	var v521 int32
	_ = v521
	var v523 int32
	_ = v523
	var v526 int32
	_ = v526
	var v530 int32
	_ = v530
	var v531 int32
	_ = v531
	var v559 int32
	_ = v559
	var v561 int32
	_ = v561
	var v564 int32
	_ = v564
	var v582 int32
	_ = v582
	var v583 int32
	_ = v583
	var v587 int32
	_ = v587
	v27 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	*(*uint8)(unsafe.Add(mBase, uint32(l4))) = uint8(v27)
	if l1 < int32(2) {
		if l2 < int32(2) {
		} else {
			v451 = l2 + int32(-1)
			v452 = int32(3)
			v453 = v451 & v452
			if base.Ui32(l2+int32(-2)) < base.Ui32(v452) {
				v526 = l0
				v530 = l4
				v531 = v27
			} else {
				v463 = l3 << (uint(int32(1)) % 32)
				v466 = l3 * int32(3)
				v469 = l3 << (uint(int32(2)) % 32)
				v480 = v27
				v484 = v451 & int32(-4)
				v485 = int32(0)
				for {
					v503 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+l3+v485))))
					v504 = v503 - v480
					*(*uint8)(unsafe.Add(mBase, uint32(l4+l3+v485))) = uint8(v504)
					v508 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v463+v485))))
					v509 = v508 - v503
					*(*uint8)(unsafe.Add(mBase, uint32(l4+v463+v485))) = uint8(v509)
					v513 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v466+v485))))
					v514 = v513 - v508
					*(*uint8)(unsafe.Add(mBase, uint32(l4+v466+v485))) = uint8(v514)
					v518 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v469+v485))))
					v519 = v518 - v513
					*(*uint8)(unsafe.Add(mBase, uint32(l4+v469+v485))) = uint8(v519)
					v521 = v485 + v469
					v523 = v484 + int32(-4)
					if v523 != 0 {
						v480 = v518
						v484 = v523
						v485 = v521
						continue
					} else {
						break
					}
					break
				}
				v526 = l0 + v521
				v530 = l4 + v521
				v531 = v518
			}
			if v453 == int32(0) {
			} else {
				v559 = v531
				v561 = v453
				v564 = l3
				for {
					v582 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v526+v564))))
					v583 = v582 - v559
					*(*uint8)(unsafe.Add(mBase, uint32(v530+v564))) = uint8(v583)
					v587 = v561 + int32(-1)
					if v587 != 0 {
						v559 = v582
						v561 = v587
						v564 = v564 + l3
						continue
					} else {
						break
					}
					break
				}
			}
		}
	} else {
		if base.Ui32(l1) < base.Ui32(int32(17)) {
			v88 = int32(0)
			v122 = v88 + l4 + int32(1)
			v124 = v88 ^ int32(-1) + l1
			v125 = l0 + v88
			for {
				v141 = int32(1)
				v142 = v125 + v141
				v143 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v142))))
				v144 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v125))))
				v145 = v143 - v144
				*(*uint8)(unsafe.Add(mBase, uint32(v122))) = uint8(v145)
				v150 = v124 + int32(-1)
				if v150 != 0 {
					v122 = v122 + v141
					v124 = v150
					v125 = v142
					continue
				} else {
					break
				}
				break
			}
		} else {
			v37 = l1 + int32(-1)
			v39 = v37 & int32(-16)
			v47 = l4 + int32(1)
			v49 = v39
			v50 = l0
			for {
				v68 = int32(0)
				v69 = base.Simd_g_v128_load_rng(m, v50+int32(1), v68, int32(-1), int32(17))
				v71 = base.Simd_g_v128_load_nc(m, v50, v68)
				v72 = base.Simd_g_i8x16_sub(v69, v71)
				base.Simd_g_v128_store(m, v47, v68, v72)
				v75 = int32(16)
				v80 = v49 + int32(-16)
				if v80 != 0 {
					v47 = v47 + v75
					v49 = v80
					v50 = v50 + v75
					continue
				} else {
					break
				}
				break
			}
			if v37 == v39 {
			} else {
				v88 = v39
				v122 = v88 + l4 + int32(1)
				v124 = v88 ^ int32(-1) + l1
				v125 = l0 + v88
				for {
					v141 = int32(1)
					v142 = v125 + v141
					v143 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v142))))
					v144 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v125))))
					v145 = v143 - v144
					*(*uint8)(unsafe.Add(mBase, uint32(v122))) = uint8(v145)
					v150 = v124 + int32(-1)
					if v150 != 0 {
						v122 = v122 + v141
						v124 = v150
						v125 = v142
						continue
					} else {
						break
					}
					break
				}
			}
		}
		if l2 < int32(2) {
		} else {
			v179 = int32(-1)
			v181 = int32(1)
			v182 = l3 + v181
			v184 = l1 + v179
			v186 = v184 & int32(-16)
			v192 = l0
			v196 = l4
			v197 = v27
			v203 = l0 + v179
			v209 = v181
			for {
				v218 = v196 + l3
				v219 = v192 + l3
				v220 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v219))))
				v221 = v220 - v197
				*(*uint8)(unsafe.Add(mBase, uint32(v218))) = uint8(v221)
				if base.Ui32(l1) < base.Ui32(int32(17)) {
					v343 = v220
					v345 = int32(1)
					v372 = v343
					v374 = v203 + v345
					v376 = l1 - v345
					v377 = l3 + v345
					for {
						v397 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v192+v377))))
						v399 = v374 + int32(1)
						v400 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v399))))
						v402 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v374))))
						v403 = v372&int32(255) + v400 - v402
						v404 = int32(0)
						if v404 < v403 {
							v407 = v403
						} else {
							v407 = v404
						}
						v408 = int32(255)
						if v407 < v408 {
							v411 = v407
						} else {
							v411 = v408
						}
						v412 = v397 - v411
						*(*uint8)(unsafe.Add(mBase, uint32(v196+v377))) = uint8(v412)
						v417 = v376 + int32(-1)
						if v417 != 0 {
							v372 = v397
							v374 = v399
							v376 = v417
							v377 = v377 + int32(1)
							continue
						} else {
							break
						}
						break
					}
				} else {
					v230 = v186
					v232 = v196
					v235 = v192
					v246 = base.Simd_g_i8x16_splat(v220)
					for {
						v253 = int32(0)
						v254 = base.Simd_g_v128_load(m, v235+v182, v253)
						v258 = base.Simd_g_v128_load_rng(m, v235+int32(1), v253, int32(-1), int32(17))
						v267 = base.Simd_g_v128_load_nc(m, v235, v253)
						v271 = base.Simd_g_const(&F_GradientFilter_C__k0)
						v273 = base.Simd_g_const(&F_GradientFilter_C__k1)
						v275 = base.Simd_g_const(&F_GradientFilter_C__k2)
						v292 = base.Simd_g_const(&F_GradientFilter_C__k3)
						v308 = base.Simd_g_const(&F_GradientFilter_C__k4)
						v326 = base.Simd_g_i8x16_sub(v254, base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v258)), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v246, v254, base.Simd_g_const(&F_GradientFilter_C__k5), base.Simd_g_const(&F_GradientFilter_C__k6))))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v267))), v271), v273), base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v258, v254, base.Simd_g_const(&F_GradientFilter_C__k2), base.Simd_g_const(&F_GradientFilter_C__k7)))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v254, v254, base.Simd_g_const(&F_GradientFilter_C__k8), base.Simd_g_const(&F_GradientFilter_C__k7))))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v267, v254, base.Simd_g_const(&F_GradientFilter_C__k2), base.Simd_g_const(&F_GradientFilter_C__k7))))), v271), v273)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v258, v254, base.Simd_g_const(&F_GradientFilter_C__k3), base.Simd_g_const(&F_GradientFilter_C__k7)))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v254, v254, base.Simd_g_const(&F_GradientFilter_C__k9), base.Simd_g_const(&F_GradientFilter_C__k7))))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v267, v254, base.Simd_g_const(&F_GradientFilter_C__k3), base.Simd_g_const(&F_GradientFilter_C__k7))))), v271), v273), base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v258, v254, base.Simd_g_const(&F_GradientFilter_C__k4), base.Simd_g_const(&F_GradientFilter_C__k7)))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v254, v254, base.Simd_g_const(&F_GradientFilter_C__k10), base.Simd_g_const(&F_GradientFilter_C__k7))))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v267, v254, base.Simd_g_const(&F_GradientFilter_C__k4), base.Simd_g_const(&F_GradientFilter_C__k7))))), v271), v273))))
						base.Simd_g_v128_store(m, v232+v182, v253, v326)
						v329 = int32(16)
						v334 = v230 + int32(-16)
						if v334 != 0 {
							v230 = v334
							v232 = v232 + v329
							v235 = v235 + v329
							v246 = v254
							continue
						} else {
							break
						}
						break
					}
					if v184 == v186 {
					} else {
						v343 = base.Simd_g_i8x16_extract_lane_u_l15(v254)
						v345 = v186 | v181
						v372 = v343
						v374 = v203 + v345
						v376 = l1 - v345
						v377 = l3 + v345
						for {
							v397 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v192+v377))))
							v399 = v374 + int32(1)
							v400 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v399))))
							v402 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v374))))
							v403 = v372&int32(255) + v400 - v402
							v404 = int32(0)
							if v404 < v403 {
								v407 = v403
							} else {
								v407 = v404
							}
							v408 = int32(255)
							if v407 < v408 {
								v411 = v407
							} else {
								v411 = v408
							}
							v412 = v397 - v411
							*(*uint8)(unsafe.Add(mBase, uint32(v196+v377))) = uint8(v412)
							v417 = v376 + int32(-1)
							if v417 != 0 {
								v372 = v397
								v374 = v399
								v376 = v417
								v377 = v377 + int32(1)
								continue
							} else {
								break
							}
							break
						}
					}
				}
				v446 = v209 + int32(1)
				if v446 != l2 {
					v192 = v219
					v196 = v218
					v197 = v220
					v203 = v203 + l3
					v209 = v446
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

var F_GradientFilter_C__k0 = [2]uint64{0x0, 0x0}
var F_GradientFilter_C__k1 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_GradientFilter_C__k2 = [2]uint64{0x7060504, 0x0}
var F_GradientFilter_C__k3 = [2]uint64{0xb0a0908, 0x0}
var F_GradientFilter_C__k4 = [2]uint64{0xf0e0d0c, 0x0}
var F_GradientFilter_C__k5 = [2]uint64{0x808080808080800f, 0x8080808080808080}
var F_GradientFilter_C__k6 = [2]uint64{0x605040302010080, 0xe0d0c0b0a090807}
var F_GradientFilter_C__k7 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_GradientFilter_C__k8 = [2]uint64{0x6050403, 0x0}
var F_GradientFilter_C__k9 = [2]uint64{0xa090807, 0x0}
var F_GradientFilter_C__k10 = [2]uint64{0xe0d0c0b, 0x0}

func F_GradientFilter_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v30 int32
	_ = v30
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v40 int32
	_ = v40
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v83 int32
	_ = v83
	var v86 int32
	_ = v86
	var v87 base.V128
	_ = v87
	var v91 base.V128
	_ = v91
	var v92 base.V128
	_ = v92
	var v95 int32
	_ = v95
	var v100 base.V128
	_ = v100
	var v102 base.V128
	_ = v102
	var v103 base.V128
	_ = v103
	var v106 int32
	_ = v106
	var v111 base.V128
	_ = v111
	var v115 base.V128
	_ = v115
	var v116 base.V128
	_ = v116
	var v119 int32
	_ = v119
	var v124 base.V128
	_ = v124
	var v128 base.V128
	_ = v128
	var v129 base.V128
	_ = v129
	var v133 int32
	_ = v133
	var v135 int32
	_ = v135
	var v144 int32
	_ = v144
	var v167 int32
	_ = v167
	var v169 int32
	_ = v169
	var v172 int32
	_ = v172
	var v173 int32
	_ = v173
	var v174 base.V128
	_ = v174
	var v176 base.V128
	_ = v176
	var v177 base.V128
	_ = v177
	var v180 int32
	_ = v180
	var v181 base.V128
	_ = v181
	var v185 base.V128
	_ = v185
	var v186 base.V128
	_ = v186
	var v199 int32
	_ = v199
	var v223 int32
	_ = v223
	var v231 int32
	_ = v231
	var v240 int32
	_ = v240
	var v243 int32
	_ = v243
	var v244 int32
	_ = v244
	var v264 int32
	_ = v264
	var v265 base.V128
	_ = v265
	var v267 base.V128
	_ = v267
	var v268 base.V128
	_ = v268
	var v271 int32
	_ = v271
	var v276 int32
	_ = v276
	var v286 int32
	_ = v286
	var v321 int32
	_ = v321
	var v322 int32
	_ = v322
	var v325 int32
	_ = v325
	var v343 int32
	_ = v343
	var v344 int32
	_ = v344
	var v345 int32
	_ = v345
	var v346 int32
	_ = v346
	var v347 int32
	_ = v347
	var v352 int32
	_ = v352
	var v385 int32
	_ = v385
	var v390 int32
	_ = v390
	var v392 int32
	_ = v392
	var v395 int32
	_ = v395
	var v403 int32
	_ = v403
	var v405 int32
	_ = v405
	var v409 int32
	_ = v409
	var v413 int32
	_ = v413
	var v414 int32
	_ = v414
	var v423 int32
	_ = v423
	var v424 int32
	_ = v424
	var v425 int32
	_ = v425
	var v426 int32
	_ = v426
	var v427 int32
	_ = v427
	var v431 int32
	_ = v431
	var v438 int32
	_ = v438
	var v439 int32
	_ = v439
	var v440 int32
	_ = v440
	var v441 int32
	_ = v441
	var v451 int32
	_ = v451
	var v452 int32
	_ = v452
	var v454 int32
	_ = v454
	var v455 int32
	_ = v455
	var v474 int32
	_ = v474
	var v475 int32
	_ = v475
	var v477 int32
	_ = v477
	var v478 base.V128
	_ = v478
	var v479 int32
	_ = v479
	var v483 base.V128
	_ = v483
	var v484 base.V128
	_ = v484
	var v485 base.V128
	_ = v485
	var v488 base.V128
	_ = v488
	var v493 base.V128
	_ = v493
	var v503 int32
	_ = v503
	var v510 int32
	_ = v510
	var v517 int32
	_ = v517
	var v524 int32
	_ = v524
	var v525 int32
	_ = v525
	var v528 int32
	_ = v528
	var v545 int32
	_ = v545
	var v547 int32
	_ = v547
	var v548 base.V128
	_ = v548
	var v552 base.V128
	_ = v552
	var v556 base.V128
	_ = v556
	var v561 base.V128
	_ = v561
	var v565 base.V128
	_ = v565
	var v567 base.V128
	_ = v567
	var v569 base.V128
	_ = v569
	var v586 base.V128
	_ = v586
	var v602 base.V128
	_ = v602
	var v620 base.V128
	_ = v620
	var v623 int32
	_ = v623
	var v630 int32
	_ = v630
	var v639 int32
	_ = v639
	var v668 int32
	_ = v668
	var v691 int32
	_ = v691
	var v692 int32
	_ = v692
	var v694 int32
	_ = v694
	var v695 int32
	_ = v695
	var v698 int32
	_ = v698
	var v699 int32
	_ = v699
	var v701 int32
	_ = v701
	var v702 int32
	_ = v702
	var v703 int32
	_ = v703
	var v706 int32
	_ = v706
	var v707 int32
	_ = v707
	var v710 int32
	_ = v710
	var v711 int32
	_ = v711
	var v714 int32
	_ = v714
	var v751 int32
	_ = v751
	var v756 int32
	_ = v756
	var v757 int32
	_ = v757
	var v758 int32
	_ = v758
	var v768 int32
	_ = v768
	var v771 int32
	_ = v771
	var v774 int32
	_ = v774
	var v785 int32
	_ = v785
	var v787 int32
	_ = v787
	var v791 int32
	_ = v791
	var v811 int32
	_ = v811
	var v812 int32
	_ = v812
	var v816 int32
	_ = v816
	var v817 int32
	_ = v817
	var v821 int32
	_ = v821
	var v822 int32
	_ = v822
	var v826 int32
	_ = v826
	var v827 int32
	_ = v827
	var v829 int32
	_ = v829
	var v831 int32
	_ = v831
	var v834 int32
	_ = v834
	var v838 int32
	_ = v838
	var v839 int32
	_ = v839
	var v870 int32
	_ = v870
	var v872 int32
	_ = v872
	var v873 int32
	_ = v873
	var v896 int32
	_ = v896
	var v897 int32
	_ = v897
	var v901 int32
	_ = v901
	var v904 int32
	_ = v904
	var v906 int32
	_ = v906
	var v911 int32
	_ = v911
	var v912 int32
	_ = v912
	var v915 int32
	_ = v915
	var v916 int32
	_ = v916
	var v920 int32
	_ = v920
	var v921 int32
	_ = v921
	var v940 int32
	_ = v940
	var v941 int32
	_ = v941
	var v942 int32
	_ = v942
	var v943 int32
	_ = v943
	var v954 int32
	_ = v954
	var v977 int32
	_ = v977
	var v978 int32
	_ = v978
	var v980 int32
	_ = v980
	var v981 base.V128
	_ = v981
	var v982 int32
	_ = v982
	var v986 base.V128
	_ = v986
	var v990 base.V128
	_ = v990
	var v995 base.V128
	_ = v995
	var v999 base.V128
	_ = v999
	var v1001 base.V128
	_ = v1001
	var v1003 base.V128
	_ = v1003
	var v1020 base.V128
	_ = v1020
	var v1036 base.V128
	_ = v1036
	var v1054 base.V128
	_ = v1054
	var v1058 int32
	_ = v1058
	var v1068 int32
	_ = v1068
	var v1097 int32
	_ = v1097
	var v1120 int32
	_ = v1120
	var v1121 int32
	_ = v1121
	var v1123 int32
	_ = v1123
	var v1124 int32
	_ = v1124
	var v1127 int32
	_ = v1127
	var v1128 int32
	_ = v1128
	var v1130 int32
	_ = v1130
	var v1131 int32
	_ = v1131
	var v1132 int32
	_ = v1132
	var v1135 int32
	_ = v1135
	var v1136 int32
	_ = v1136
	var v1139 int32
	_ = v1139
	var v1140 int32
	_ = v1140
	var v1143 int32
	_ = v1143
	var v1177 int32
	_ = v1177
	v30 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	*(*uint8)(unsafe.Add(mBase, uint32(l4))) = uint8(v30)
	v33 = l1 + int32(-1)
	v35 = v33 & int32(-32)
	if int32(1) <= v35 {
		v40 = v35 + int32(-1)
		if v40 != int32(31) {
			v59 = int32(0)
			v61 = (int32(base.Ui32(v40)>>(uint(int32(5))%32)) + int32(1)) & int32(268435454)
			for {
				v80 = l4 + v59
				v81 = int32(17)
				v83 = l0 + v59
				v86 = int32(0)
				v87 = base.Simd_g_v128_load_rng(m, v83+v81, v86, int32(-1), int32(17))
				v91 = base.Simd_g_v128_load_nc(m, v83+int32(16), v86)
				v92 = base.Simd_g_i8x16_sub(v87, v91)
				base.Simd_g_v128_store(m, v80+v81, v86, v92)
				v95 = int32(1)
				v100 = base.Simd_g_v128_load_rng(m, v83+v95, v86, int32(-1), int32(17))
				v102 = base.Simd_g_v128_load_nc(m, v83, v86)
				v103 = base.Simd_g_i8x16_sub(v100, v102)
				base.Simd_g_v128_store(m, v80+v95, v86, v103)
				v106 = int32(49)
				v111 = base.Simd_g_v128_load_rng(m, v83+v106, v86, int32(-1), int32(17))
				v115 = base.Simd_g_v128_load_nc(m, v83+int32(48), v86)
				v116 = base.Simd_g_i8x16_sub(v111, v115)
				base.Simd_g_v128_store(m, v80+v106, v86, v116)
				v119 = int32(33)
				v124 = base.Simd_g_v128_load_rng(m, v83+v119, v86, int32(-1), int32(17))
				v128 = base.Simd_g_v128_load_nc(m, v83+int32(32), v86)
				v129 = base.Simd_g_i8x16_sub(v124, v128)
				base.Simd_g_v128_store(m, v80+v119, v86, v129)
				v133 = v59 + int32(64)
				v135 = v61 + int32(-2)
				if v135 != 0 {
					v59 = v133
					v61 = v135
					continue
				} else {
					break
				}
				break
			}
			v144 = v133
		} else {
			v144 = int32(0)
		}
		if v40&int32(32) != 0 {
			v199 = v144
		} else {
			v167 = int32(1)
			v169 = l4 + v167 + v144
			v172 = l0 + v167 + v144
			v173 = int32(16)
			v174 = base.Simd_g_v128_load_rng(m, v172, v173, int32(15), int32(17))
			v176 = base.Simd_g_v128_load_nc(m, v172, int32(15))
			v177 = base.Simd_g_i8x16_sub(v174, v176)
			base.Simd_g_v128_store(m, v169, v173, v177)
			v180 = int32(0)
			v181 = base.Simd_g_v128_load(m, v172, v180)
			v185 = base.Simd_g_v128_load(m, v172+int32(-1), v180)
			v186 = base.Simd_g_i8x16_sub(v181, v185)
			base.Simd_g_v128_store(m, v169, v180, v186)
			v199 = v144 + int32(32)
		}
	} else {
		v199 = int32(0)
	}
	if v33 <= v199 {
	} else {
		v223 = v199 ^ int32(-1) + l1
		if base.Ui32(v223) <= base.Ui32(int32(15)) {
			v286 = v199
			v321 = l0 + v286
			v322 = v286 ^ int32(-1) + l1
			v325 = v286 + l4 + int32(1)
			for {
				v343 = int32(1)
				v344 = v321 + v343
				v345 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v344))))
				v346 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321))))
				v347 = v345 - v346
				*(*uint8)(unsafe.Add(mBase, uint32(v325))) = uint8(v347)
				v352 = v322 + int32(-1)
				if v352 != 0 {
					v321 = v344
					v322 = v352
					v325 = v325 + v343
					continue
				} else {
					break
				}
				break
			}
		} else {
			v231 = v223 & int32(-16)
			v240 = l0 + v199
			v243 = v231
			v244 = v199 + l4 + int32(1)
			for {
				v264 = int32(0)
				v265 = base.Simd_g_v128_load_rng(m, v240+int32(1), v264, int32(-1), int32(17))
				v267 = base.Simd_g_v128_load_nc(m, v240, v264)
				v268 = base.Simd_g_i8x16_sub(v265, v267)
				base.Simd_g_v128_store(m, v244, v264, v268)
				v271 = int32(16)
				v276 = v243 + int32(-16)
				if v276 != 0 {
					v240 = v240 + v271
					v243 = v276
					v244 = v244 + v271
					continue
				} else {
					break
				}
				break
			}
			if v223 == v231 {
			} else {
				v286 = v199 + v231
				v321 = l0 + v286
				v322 = v286 ^ int32(-1) + l1
				v325 = v286 + l4 + int32(1)
				for {
					v343 = int32(1)
					v344 = v321 + v343
					v345 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v344))))
					v346 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321))))
					v347 = v345 - v346
					*(*uint8)(unsafe.Add(mBase, uint32(v325))) = uint8(v347)
					v352 = v322 + int32(-1)
					if v352 != 0 {
						v321 = v344
						v322 = v352
						v325 = v325 + v343
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
		v385 = v33 & int32(-8)
		if v385 < int32(1) {
			if int32(1) < l1 {
				v904 = v33 & int32(-16)
				v906 = int32(1)
				v911 = l0
				v912 = v906
				v915 = l4
				v916 = v30
				v920 = l3 + l4 + v906
				v921 = l0 + l3
				for {
					v940 = v915 + l3
					v941 = v911 + l3
					v942 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v941))))
					v943 = v942 - v916
					*(*uint8)(unsafe.Add(mBase, uint32(v940))) = uint8(v943)
					if base.Ui32(l1) < base.Ui32(int32(17)) {
						v1068 = int32(0)
						v1097 = v1068
						for {
							v1120 = v921 + v1097
							v1121 = int32(1)
							v1123 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1120+v1121))))
							v1124 = v911 + v1097
							v1127 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1124+v1121))))
							v1128 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1120))))
							v1130 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1124))))
							v1131 = v1127 + v1128 - v1130
							v1132 = int32(0)
							if v1132 < v1131 {
								v1135 = v1131
							} else {
								v1135 = v1132
							}
							v1136 = int32(255)
							if v1135 < v1136 {
								v1139 = v1135
							} else {
								v1139 = v1136
							}
							v1140 = v1123 - v1139
							*(*uint8)(unsafe.Add(mBase, uint32(v920+v1097))) = uint8(v1140)
							v1143 = v1097 + int32(1)
							if v33 != v1143 {
								v1097 = v1143
								continue
							} else {
								break
							}
							break
						}
					} else {
						v954 = int32(0)
						for {
							v977 = v921 + v954
							v978 = int32(1)
							v980 = int32(0)
							v981 = base.Simd_g_v128_load_rng(m, v977+v978, v980, int32(-1), int32(17))
							v982 = v911 + v954
							v986 = base.Simd_g_v128_load_rng(m, v982+v978, v980, int32(-1), int32(17))
							v990 = base.Simd_g_v128_load_nc(m, v977, v980)
							v995 = base.Simd_g_v128_load_nc(m, v982, v980)
							v999 = base.Simd_g_const(&F_GradientFilter_SSE2__k0)
							v1001 = base.Simd_g_const(&F_GradientFilter_SSE2__k1)
							v1003 = base.Simd_g_const(&F_GradientFilter_SSE2__k2)
							v1020 = base.Simd_g_const(&F_GradientFilter_SSE2__k3)
							v1036 = base.Simd_g_const(&F_GradientFilter_SSE2__k4)
							v1054 = base.Simd_g_i8x16_sub(v981, base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v986)), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v990))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v995))), v999), v1001), base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v986, v986, base.Simd_g_const(&F_GradientFilter_SSE2__k2), base.Simd_g_const(&F_GradientFilter_SSE2__k5)))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v990, v986, base.Simd_g_const(&F_GradientFilter_SSE2__k2), base.Simd_g_const(&F_GradientFilter_SSE2__k5))))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v995, v986, base.Simd_g_const(&F_GradientFilter_SSE2__k2), base.Simd_g_const(&F_GradientFilter_SSE2__k5))))), v999), v1001)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v986, v986, base.Simd_g_const(&F_GradientFilter_SSE2__k3), base.Simd_g_const(&F_GradientFilter_SSE2__k5)))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v990, v986, base.Simd_g_const(&F_GradientFilter_SSE2__k3), base.Simd_g_const(&F_GradientFilter_SSE2__k5))))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v995, v986, base.Simd_g_const(&F_GradientFilter_SSE2__k3), base.Simd_g_const(&F_GradientFilter_SSE2__k5))))), v999), v1001), base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v986, v986, base.Simd_g_const(&F_GradientFilter_SSE2__k4), base.Simd_g_const(&F_GradientFilter_SSE2__k5)))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v990, v986, base.Simd_g_const(&F_GradientFilter_SSE2__k4), base.Simd_g_const(&F_GradientFilter_SSE2__k5))))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v995, v986, base.Simd_g_const(&F_GradientFilter_SSE2__k4), base.Simd_g_const(&F_GradientFilter_SSE2__k5))))), v999), v1001))))
							base.Simd_g_v128_store(m, v920+v954, v980, v1054)
							v1058 = v954 + int32(16)
							if v904 != v1058 {
								v954 = v1058
								continue
							} else {
								break
							}
							break
						}
						if v33 == v904 {
						} else {
							v1068 = v904
							v1097 = v1068
							for {
								v1120 = v921 + v1097
								v1121 = int32(1)
								v1123 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1120+v1121))))
								v1124 = v911 + v1097
								v1127 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1124+v1121))))
								v1128 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1120))))
								v1130 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1124))))
								v1131 = v1127 + v1128 - v1130
								v1132 = int32(0)
								if v1132 < v1131 {
									v1135 = v1131
								} else {
									v1135 = v1132
								}
								v1136 = int32(255)
								if v1135 < v1136 {
									v1139 = v1135
								} else {
									v1139 = v1136
								}
								v1140 = v1123 - v1139
								*(*uint8)(unsafe.Add(mBase, uint32(v920+v1097))) = uint8(v1140)
								v1143 = v1097 + int32(1)
								if v33 != v1143 {
									v1097 = v1143
									continue
								} else {
									break
								}
								break
							}
						}
					}
					v1177 = v912 + int32(1)
					if v1177 != l2 {
						v911 = v941
						v912 = v1177
						v915 = v940
						v916 = v942
						v920 = v920 + l3
						v921 = v921 + l3
						continue
					} else {
						break
					}
					break
				}
			} else {
				v756 = l2 + int32(-1)
				v757 = int32(3)
				v758 = v756 & v757
				if base.Ui32(l2+int32(-2)) < base.Ui32(v757) {
					v834 = l0
					v838 = l4
					v839 = v30
				} else {
					v768 = l3 << (uint(int32(1)) % 32)
					v771 = l3 * int32(3)
					v774 = l3 << (uint(int32(2)) % 32)
					v785 = v30
					v787 = int32(0)
					v791 = v756 & int32(-4)
					for {
						v811 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+l3+v787))))
						v812 = v811 - v785
						*(*uint8)(unsafe.Add(mBase, uint32(l4+l3+v787))) = uint8(v812)
						v816 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v768+v787))))
						v817 = v816 - v811
						*(*uint8)(unsafe.Add(mBase, uint32(l4+v768+v787))) = uint8(v817)
						v821 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v771+v787))))
						v822 = v821 - v816
						*(*uint8)(unsafe.Add(mBase, uint32(l4+v771+v787))) = uint8(v822)
						v826 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v774+v787))))
						v827 = v826 - v821
						*(*uint8)(unsafe.Add(mBase, uint32(l4+v774+v787))) = uint8(v827)
						v829 = v787 + v774
						v831 = v791 + int32(-4)
						if v831 != 0 {
							v785 = v826
							v787 = v829
							v791 = v831
							continue
						} else {
							break
						}
						break
					}
					v834 = l0 + v829
					v838 = l4 + v829
					v839 = v826
				}
				if v758 == int32(0) {
				} else {
					v870 = v839
					v872 = l3
					v873 = v758
					for {
						v896 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v834+v872))))
						v897 = v896 - v870
						*(*uint8)(unsafe.Add(mBase, uint32(v838+v872))) = uint8(v897)
						v901 = v873 + int32(-1)
						if v901 != 0 {
							v870 = v896
							v872 = v872 + l3
							v873 = v901
							continue
						} else {
							break
						}
						break
					}
				}
			}
		} else {
			v390 = l0 + int32(8)
			v392 = l3 + l4
			v395 = int32(1)
			v403 = l1 - (v385+int32(-1))&int32(-8) + int32(-9)
			v405 = v403 & int32(-16)
			v409 = l0
			v413 = l4
			v414 = v30
			v423 = l0 + l3
			v424 = v390
			v425 = v390 + l3
			v426 = v392 + int32(9)
			v427 = v392 + v395
			v431 = v395
			for {
				v438 = v413 + l3
				v439 = v409 + l3
				v440 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v439))))
				v441 = v440 - v414
				*(*uint8)(unsafe.Add(mBase, uint32(v438))) = uint8(v441)
				v451 = int32(0)
				v452 = v425
				v454 = v426
				v455 = v424
				for {
					v474 = v423 + v451
					v475 = int32(1)
					v477 = int32(0)
					v478 = base.Simd_g_v128_load64_splat(m, v474+v475, v477)
					v479 = v409 + v451
					v483 = base.Simd_g_v128_load64_zero(m, v479+v475, v477)
					v484 = base.Simd_g_const(&F_GradientFilter_SSE2__k0)
					v485 = base.Simd_g_const(&F_GradientFilter_SSE2__k6)
					v488 = base.Simd_g_v128_load64_zero(m, v474, v477)
					v493 = base.Simd_g_v128_load64_zero(m, v479, v477)
					base.Simd_g_v128_store64_lane_l0(m, v427+v451, v477, base.Simd_g_i8x16_sub(v478, base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_sub(base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v483, v484, base.Simd_g_const(&F_GradientFilter_SSE2__k7), base.Simd_g_const(&F_GradientFilter_SSE2__k8)), base.Simd_g_i8x16_shuffle2(v488, v484, base.Simd_g_const(&F_GradientFilter_SSE2__k7), base.Simd_g_const(&F_GradientFilter_SSE2__k8))), base.Simd_g_i8x16_shuffle2(v493, v484, base.Simd_g_const(&F_GradientFilter_SSE2__k7), base.Simd_g_const(&F_GradientFilter_SSE2__k8))), v484)))
					v503 = int32(8)
					v510 = v451 + v503
					if v510 < v385 {
						v451 = v510
						v452 = v452 + v503
						v454 = v454 + v503
						v455 = v455 + v503
						continue
					} else {
						break
					}
					break
				}
				if v33 <= v510 {
				} else {
					if base.B2i32(base.Ui32(int32(15)) < base.Ui32(v403)) == int32(0) {
						v639 = v510
						v668 = v639
						for {
							v691 = v423 + v668
							v692 = int32(1)
							v694 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v691+v692))))
							v695 = v409 + v668
							v698 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v695+v692))))
							v699 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v691))))
							v701 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v695))))
							v702 = v698 + v699 - v701
							v703 = int32(0)
							if v703 < v702 {
								v706 = v702
							} else {
								v706 = v703
							}
							v707 = int32(255)
							if v706 < v707 {
								v710 = v706
							} else {
								v710 = v707
							}
							v711 = v694 - v710
							*(*uint8)(unsafe.Add(mBase, uint32(v427+v668))) = uint8(v711)
							v714 = v668 + int32(1)
							if v33 != v714 {
								v668 = v714
								continue
							} else {
								break
							}
							break
						}
					} else {
						v517 = v454
						v524 = v405
						v525 = v452
						v528 = v455
						for {
							v545 = int32(1)
							v547 = int32(0)
							v548 = base.Simd_g_v128_load_rng(m, v525+v545, v547, int32(-1), int32(17))
							v552 = base.Simd_g_v128_load_rng(m, v528+v545, v547, int32(-1), int32(17))
							v556 = base.Simd_g_v128_load_nc(m, v525, v547)
							v561 = base.Simd_g_v128_load_nc(m, v528, v547)
							v565 = base.Simd_g_const(&F_GradientFilter_SSE2__k0)
							v567 = base.Simd_g_const(&F_GradientFilter_SSE2__k1)
							v569 = base.Simd_g_const(&F_GradientFilter_SSE2__k2)
							v586 = base.Simd_g_const(&F_GradientFilter_SSE2__k3)
							v602 = base.Simd_g_const(&F_GradientFilter_SSE2__k4)
							v620 = base.Simd_g_i8x16_sub(v548, base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v552)), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v556))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v561))), v565), v567), base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v552, v552, base.Simd_g_const(&F_GradientFilter_SSE2__k2), base.Simd_g_const(&F_GradientFilter_SSE2__k5)))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v556, v552, base.Simd_g_const(&F_GradientFilter_SSE2__k2), base.Simd_g_const(&F_GradientFilter_SSE2__k5))))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v561, v552, base.Simd_g_const(&F_GradientFilter_SSE2__k2), base.Simd_g_const(&F_GradientFilter_SSE2__k5))))), v565), v567)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v552, v552, base.Simd_g_const(&F_GradientFilter_SSE2__k3), base.Simd_g_const(&F_GradientFilter_SSE2__k5)))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v556, v552, base.Simd_g_const(&F_GradientFilter_SSE2__k3), base.Simd_g_const(&F_GradientFilter_SSE2__k5))))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v561, v552, base.Simd_g_const(&F_GradientFilter_SSE2__k3), base.Simd_g_const(&F_GradientFilter_SSE2__k5))))), v565), v567), base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_sub(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v552, v552, base.Simd_g_const(&F_GradientFilter_SSE2__k4), base.Simd_g_const(&F_GradientFilter_SSE2__k5)))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v556, v552, base.Simd_g_const(&F_GradientFilter_SSE2__k4), base.Simd_g_const(&F_GradientFilter_SSE2__k5))))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v561, v552, base.Simd_g_const(&F_GradientFilter_SSE2__k4), base.Simd_g_const(&F_GradientFilter_SSE2__k5))))), v565), v567))))
							base.Simd_g_v128_store(m, v517, v547, v620)
							v623 = int32(16)
							v630 = v524 + int32(-16)
							if v630 != 0 {
								v517 = v517 + v623
								v524 = v630
								v525 = v525 + v623
								v528 = v528 + v623
								continue
							} else {
								break
							}
							break
						}
						if v403 == v405 {
						} else {
							v639 = v405 + v510
							v668 = v639
							for {
								v691 = v423 + v668
								v692 = int32(1)
								v694 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v691+v692))))
								v695 = v409 + v668
								v698 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v695+v692))))
								v699 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v691))))
								v701 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v695))))
								v702 = v698 + v699 - v701
								v703 = int32(0)
								if v703 < v702 {
									v706 = v702
								} else {
									v706 = v703
								}
								v707 = int32(255)
								if v706 < v707 {
									v710 = v706
								} else {
									v710 = v707
								}
								v711 = v694 - v710
								*(*uint8)(unsafe.Add(mBase, uint32(v427+v668))) = uint8(v711)
								v714 = v668 + int32(1)
								if v33 != v714 {
									v668 = v714
									continue
								} else {
									break
								}
								break
							}
						}
					}
				}
				v751 = v431 + int32(1)
				if v751 != l2 {
					v409 = v439
					v413 = v438
					v414 = v440
					v423 = v423 + l3
					v424 = v424 + l3
					v425 = v425 + l3
					v426 = v426 + l3
					v427 = v427 + l3
					v431 = v751
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

var F_GradientFilter_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_GradientFilter_SSE2__k1 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_GradientFilter_SSE2__k2 = [2]uint64{0x7060504, 0x0}
var F_GradientFilter_SSE2__k3 = [2]uint64{0xb0a0908, 0x0}
var F_GradientFilter_SSE2__k4 = [2]uint64{0xf0e0d0c, 0x0}
var F_GradientFilter_SSE2__k5 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_GradientFilter_SSE2__k6 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_GradientFilter_SSE2__k7 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_GradientFilter_SSE2__k8 = [2]uint64{0x380028001800080, 0x780068005800480}

func F_GradientUnfilter_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v24 int32
	_ = v24
	var v31 int32
	_ = v31
	var v36 int32
	_ = v36
	var v43 int32
	_ = v43
	var v45 base.V128
	_ = v45
	var v63 base.V128
	_ = v63
	var v65 int32
	_ = v65
	var v66 base.V128
	_ = v66
	var v67 base.V128
	_ = v67
	var v70 base.V128
	_ = v70
	var v73 base.V128
	_ = v73
	var v76 base.V128
	_ = v76
	var v94 int32
	_ = v94
	var v116 int32
	_ = v116
	var v119 int32
	_ = v119
	var v120 int32
	_ = v120
	var v124 int32
	_ = v124
	var v128 int32
	_ = v128
	var v145 int32
	_ = v145
	var v146 int32
	_ = v146
	var v149 int32
	_ = v149
	var v151 int32
	_ = v151
	var v156 int32
	_ = v156
	var v160 int32
	_ = v160
	var v179 int32
	_ = v179
	var v180 int32
	_ = v180
	var v181 int32
	_ = v181
	var v186 int32
	_ = v186
	var v201 int32
	_ = v201
	var v202 int32
	_ = v202
	var v203 int32
	_ = v203
	var v204 int32
	_ = v204
	var v206 int32
	_ = v206
	var v210 int32
	_ = v210
	var v211 int32
	_ = v211
	var v213 int32
	_ = v213
	var v217 int32
	_ = v217
	var v218 int32
	_ = v218
	var v220 int32
	_ = v220
	var v224 int32
	_ = v224
	var v225 int32
	_ = v225
	var v227 int32
	_ = v227
	var v232 int32
	_ = v232
	var v234 int32
	_ = v234
	var v235 int32
	_ = v235
	var v240 int32
	_ = v240
	var v242 int32
	_ = v242
	var v244 int32
	_ = v244
	var v248 base.V128
	_ = v248
	var v251 int32
	_ = v251
	var v255 base.V128
	_ = v255
	var v256 base.V128
	_ = v256
	var v258 base.V128
	_ = v258
	var v260 base.V128
	_ = v260
	var v262 base.V128
	_ = v262
	var v264 base.V128
	_ = v264
	var v266 base.V128
	_ = v266
	var v274 int32
	_ = v274
	var v282 base.V128
	_ = v282
	var v294 int32
	_ = v294
	var v297 int32
	_ = v297
	var v298 base.V128
	_ = v298
	var v299 base.V128
	_ = v299
	var v302 base.V128
	_ = v302
	var v305 base.V128
	_ = v305
	var v307 base.V128
	_ = v307
	var v311 base.V128
	_ = v311
	var v314 base.V128
	_ = v314
	var v315 base.V128
	_ = v315
	var v322 base.V128
	_ = v322
	var v330 base.V128
	_ = v330
	var v339 base.V128
	_ = v339
	var v348 base.V128
	_ = v348
	var v357 base.V128
	_ = v357
	var v366 base.V128
	_ = v366
	var v375 base.V128
	_ = v375
	var v384 int32
	_ = v384
	var v390 int32
	_ = v390
	var v412 int32
	_ = v412
	var v419 int32
	_ = v419
	var v420 int32
	_ = v420
	var v421 int32
	_ = v421
	var v423 int32
	_ = v423
	var v424 int32
	_ = v424
	var v428 int32
	_ = v428
	var v443 int32
	_ = v443
	var v445 int32
	_ = v445
	var v446 int32
	_ = v446
	var v450 int32
	_ = v450
	var v451 int32
	_ = v451
	var v452 int32
	_ = v452
	var v455 int32
	_ = v455
	var v456 int32
	_ = v456
	var v459 int32
	_ = v459
	var v460 int32
	_ = v460
	var v462 int32
	_ = v462
	var v467 int32
	_ = v467
	v24 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
	if l0 != 0 {
		v234 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
		v235 = v234 + v24
		*(*uint8)(unsafe.Add(mBase, uint32(l2))) = uint8(v235)
		if l3 < int32(2) {
		} else {
			v240 = l3 + int32(-1)
			v242 = v240 & int32(2147483640)
			if v242 != 0 {
				v244 = int32(1)
				v248 = base.Simd_g_const(&F_GradientUnfilter_SSE2__k0)
				v251 = int32(0)
				v255 = base.Simd_g_const(&F_GradientUnfilter_SSE2__k1)
				v256 = base.Simd_g_i8x16_shuffle2(v248, base.Simd_g_const(&F_GradientUnfilter_SSE2__k2), base.Simd_g_const(&F_GradientUnfilter_SSE2__k3), base.Simd_g_const(&F_GradientUnfilter_SSE2__k4))
				v258 = base.Simd_g_i8x16_shuffle2(v248, v256, base.Simd_g_const(&F_GradientUnfilter_SSE2__k3), base.Simd_g_const(&F_GradientUnfilter_SSE2__k4))
				v260 = base.Simd_g_i8x16_shuffle2(v248, v258, base.Simd_g_const(&F_GradientUnfilter_SSE2__k3), base.Simd_g_const(&F_GradientUnfilter_SSE2__k4))
				v262 = base.Simd_g_i8x16_shuffle2(v248, v260, base.Simd_g_const(&F_GradientUnfilter_SSE2__k3), base.Simd_g_const(&F_GradientUnfilter_SSE2__k4))
				v264 = base.Simd_g_i8x16_shuffle2(v248, v262, base.Simd_g_const(&F_GradientUnfilter_SSE2__k3), base.Simd_g_const(&F_GradientUnfilter_SSE2__k4))
				v266 = base.Simd_g_i8x16_shuffle2(v248, v264, base.Simd_g_const(&F_GradientUnfilter_SSE2__k3), base.Simd_g_const(&F_GradientUnfilter_SSE2__k4))
				v274 = v251
				v282 = base.Simd_g_i32x4_replace_lane_l0(v248, v235&int32(255))
				for {
					v294 = l0 + v274
					v297 = int32(0)
					v298 = base.Simd_g_v128_load64_zero(m, v294+int32(1), v297)
					v299 = base.Simd_g_const(&F_GradientUnfilter_SSE2__k5)
					v302 = base.Simd_g_v128_load64_zero(m, v294, v297)
					v305 = base.Simd_g_i16x8_sub(base.Simd_g_i8x16_shuffle2(v298, v248, base.Simd_g_const(&F_GradientUnfilter_SSE2__k6), base.Simd_g_const(&F_GradientUnfilter_SSE2__k7)), base.Simd_g_i8x16_shuffle2(v302, v248, base.Simd_g_const(&F_GradientUnfilter_SSE2__k6), base.Simd_g_const(&F_GradientUnfilter_SSE2__k7)))
					v307 = base.Simd_g_const(&F_GradientUnfilter_SSE2__k0)
					v311 = base.Simd_g_v128_load64_zero(m, l1+v244+v274, v297)
					v314 = base.Simd_g_v128_and(base.Simd_g_i8x16_add(base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v305, v282), v307), v311), base.Simd_g_const(&F_GradientUnfilter_SSE2__k2))
					v315 = base.Simd_g_const(&F_GradientUnfilter_SSE2__k1)
					v322 = base.Simd_g_v128_and(base.Simd_g_i8x16_add(base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v305, base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v248, v314, base.Simd_g_const(&F_GradientUnfilter_SSE2__k3), base.Simd_g_const(&F_GradientUnfilter_SSE2__k4)), v248, base.Simd_g_const(&F_GradientUnfilter_SSE2__k6), base.Simd_g_const(&F_GradientUnfilter_SSE2__k7))), v307), v311), v256)
					v330 = base.Simd_g_v128_and(base.Simd_g_i8x16_add(base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v305, base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v248, v322, base.Simd_g_const(&F_GradientUnfilter_SSE2__k3), base.Simd_g_const(&F_GradientUnfilter_SSE2__k4)), v248, base.Simd_g_const(&F_GradientUnfilter_SSE2__k6), base.Simd_g_const(&F_GradientUnfilter_SSE2__k7))), v307), v311), v258)
					v339 = base.Simd_g_v128_and(base.Simd_g_i8x16_add(base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v305, base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v248, v330, base.Simd_g_const(&F_GradientUnfilter_SSE2__k3), base.Simd_g_const(&F_GradientUnfilter_SSE2__k4)), v248, base.Simd_g_const(&F_GradientUnfilter_SSE2__k6), base.Simd_g_const(&F_GradientUnfilter_SSE2__k7))), v307), v311), v260)
					v348 = base.Simd_g_v128_and(base.Simd_g_i8x16_add(base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v305, base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v248, v339, base.Simd_g_const(&F_GradientUnfilter_SSE2__k3), base.Simd_g_const(&F_GradientUnfilter_SSE2__k4)), v248, base.Simd_g_const(&F_GradientUnfilter_SSE2__k6), base.Simd_g_const(&F_GradientUnfilter_SSE2__k7))), v307), v311), v262)
					v357 = base.Simd_g_v128_and(base.Simd_g_i8x16_add(base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v305, base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v248, v348, base.Simd_g_const(&F_GradientUnfilter_SSE2__k3), base.Simd_g_const(&F_GradientUnfilter_SSE2__k4)), v248, base.Simd_g_const(&F_GradientUnfilter_SSE2__k6), base.Simd_g_const(&F_GradientUnfilter_SSE2__k7))), v307), v311), v264)
					v366 = base.Simd_g_v128_and(base.Simd_g_i8x16_add(base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v305, base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v248, v357, base.Simd_g_const(&F_GradientUnfilter_SSE2__k3), base.Simd_g_const(&F_GradientUnfilter_SSE2__k4)), v248, base.Simd_g_const(&F_GradientUnfilter_SSE2__k6), base.Simd_g_const(&F_GradientUnfilter_SSE2__k7))), v307), v311), v266)
					v375 = base.Simd_g_v128_and(base.Simd_g_i8x16_add(base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(v305, base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v248, v366, base.Simd_g_const(&F_GradientUnfilter_SSE2__k3), base.Simd_g_const(&F_GradientUnfilter_SSE2__k4)), v248, base.Simd_g_const(&F_GradientUnfilter_SSE2__k6), base.Simd_g_const(&F_GradientUnfilter_SSE2__k7))), v307), v311), base.Simd_g_i8x16_shuffle2(v248, v266, base.Simd_g_const(&F_GradientUnfilter_SSE2__k3), base.Simd_g_const(&F_GradientUnfilter_SSE2__k4)))
					base.Simd_g_v128_store64_lane_l0(m, l2+v244+v274, v297, base.Simd_g_v128_or(v314, base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_or(v330, v322), v339), v348), v357), v366), v375)))
					v384 = v274 + int32(8)
					if base.Ui32(v384) < base.Ui32(v242) {
						v274 = v384
						v282 = base.Simd_g_i8x16_shuffle2(v375, v248, base.Simd_g_const(&F_GradientUnfilter_SSE2__k8), base.Simd_g_const(&F_GradientUnfilter_SSE2__k9))
						continue
					} else {
						break
					}
					break
				}
				v390 = v384
			} else {
				v390 = int32(0)
			}
			if v240 <= v390 {
			} else {
				v412 = v390 + int32(1)
				v419 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v390))))
				v420 = l1 + v412
				v421 = l2 + v412
				v423 = v390 ^ int32(-1) + l3
				v424 = v419
				v428 = l0 + v390
				for {
					v443 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v420))))
					v445 = v428 + int32(1)
					v446 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v445))))
					v450 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v428))))
					v451 = v446 + v424&int32(255) - v450
					v452 = int32(0)
					if v452 < v451 {
						v455 = v451
					} else {
						v455 = v452
					}
					v456 = int32(255)
					if v455 < v456 {
						v459 = v455
					} else {
						v459 = v456
					}
					v460 = v443 + v459
					*(*uint8)(unsafe.Add(mBase, uint32(v421))) = uint8(v460)
					v462 = int32(1)
					v467 = v423 + int32(-1)
					if v467 != 0 {
						v420 = v420 + v462
						v421 = v421 + v462
						v423 = v467
						v424 = v460
						v428 = v445
						continue
					} else {
						break
					}
					break
				}
			}
		}
	} else {
		*(*uint8)(unsafe.Add(mBase, uint32(l2))) = uint8(v24)
		if l3 < int32(2) {
		} else {
			if base.Ui32(l3) < base.Ui32(int32(9)) {
				v94 = int32(1)
			} else {
				v31 = int32(1)
				v36 = int32(0)
				v43 = v36
				v45 = base.Simd_g_i32x4_replace_lane_l0(base.Simd_g_const(&F_GradientUnfilter_SSE2__k0), v24)
				for {
					v63 = base.Simd_g_const(&F_GradientUnfilter_SSE2__k0)
					v65 = int32(0)
					v66 = base.Simd_g_v128_load64_zero(m, l1+v31+v43, v65)
					v67 = base.Simd_g_i8x16_add(v66, v45)
					v70 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(v63, v67, base.Simd_g_const(&F_GradientUnfilter_SSE2__k3), base.Simd_g_const(&F_GradientUnfilter_SSE2__k4)), v67)
					v73 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(v63, v70, base.Simd_g_const(&F_GradientUnfilter_SSE2__k10), base.Simd_g_const(&F_GradientUnfilter_SSE2__k11)), v70)
					v76 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(v63, v73, base.Simd_g_const(&F_GradientUnfilter_SSE2__k12), base.Simd_g_const(&F_GradientUnfilter_SSE2__k13)), v73)
					base.Simd_g_v128_store64_lane_l0(m, l2+v31+v43, v65, v76)
					if v43+int32(17) <= l3 {
						v43 = v43 + int32(8)
						v45 = base.Simd_g_i64x2_shr_u(v76, int32(56))
						continue
					} else {
						break
					}
					break
				}
				v94 = v43 + int32(9)
			}
			if l3 <= v94 {
			} else {
				v116 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v94+int32(-1)))))
				v119 = (l3 - v94) & int32(3)
				if v119 != 0 {
					v120 = v119
					v124 = v94
					v128 = v116
					for {
						v145 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v124))))
						v146 = v128 + v145
						*(*uint8)(unsafe.Add(mBase, uint32(l2+v124))) = uint8(v146)
						v149 = v124 + int32(1)
						v151 = v120 + int32(-1)
						if v151 != 0 {
							v120 = v151
							v124 = v149
							v128 = v146
							continue
						} else {
							break
						}
						break
					}
					v156 = v149
					v160 = v146
				} else {
					v156 = v94
					v160 = v116
				}
				if base.Ui32(int32(-4)) < base.Ui32(v94-l3) {
				} else {
					v179 = l1
					v180 = l2
					v181 = l3
					v186 = v160
					for {
						v201 = v180 + v156
						v202 = v179 + v156
						v203 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v202))))
						v204 = v186 + v203
						*(*uint8)(unsafe.Add(mBase, uint32(v201))) = uint8(v204)
						v206 = int32(1)
						v210 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v202+v206))))
						v211 = v204 + v210
						*(*uint8)(unsafe.Add(mBase, uint32(v201+v206))) = uint8(v211)
						v213 = int32(2)
						v217 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v202+v213))))
						v218 = v211 + v217
						*(*uint8)(unsafe.Add(mBase, uint32(v201+v213))) = uint8(v218)
						v220 = int32(3)
						v224 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v202+v220))))
						v225 = v218 + v224
						*(*uint8)(unsafe.Add(mBase, uint32(v201+v220))) = uint8(v225)
						v227 = int32(4)
						v232 = v181 + int32(-4)
						if v156 != v232 {
							v179 = v179 + v227
							v180 = v180 + v227
							v181 = v232
							v186 = v225
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

var F_GradientUnfilter_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_GradientUnfilter_SSE2__k1 = [2]uint64{0x161514131211100f, 0x1e1d1c1b1a191817}
var F_GradientUnfilter_SSE2__k2 = [2]uint64{0xff, 0x0}
var F_GradientUnfilter_SSE2__k3 = [2]uint64{0x808080808080800f, 0x8080808080808080}
var F_GradientUnfilter_SSE2__k4 = [2]uint64{0x605040302010080, 0xe0d0c0b0a090807}
var F_GradientUnfilter_SSE2__k5 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_GradientUnfilter_SSE2__k6 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_GradientUnfilter_SSE2__k7 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_GradientUnfilter_SSE2__k8 = [2]uint64{0xe0d0c0b0a090807, 0x808080808080800f}
var F_GradientUnfilter_SSE2__k9 = [2]uint64{0x8080808080808080, 0x605040302010080}
var F_GradientUnfilter_SSE2__k10 = [2]uint64{0x8080808080800f0e, 0x8080808080808080}
var F_GradientUnfilter_SSE2__k11 = [2]uint64{0x504030201008080, 0xd0c0b0a09080706}
var F_GradientUnfilter_SSE2__k12 = [2]uint64{0x808080800f0e0d0c, 0x8080808080808080}
var F_GradientUnfilter_SSE2__k13 = [2]uint64{0x302010080808080, 0xb0a090807060504}
