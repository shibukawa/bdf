//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	"unsafe"
)

func F_FTransform2_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v8 int32
	_ = v8
	var v10 int32
	_ = v10
	var v16 int32
	_ = v16
	v5 = m.G1
	v8 = *(*int32)(unsafe.Add(mBase, uint32(v5)+uint32(_c_F_FTransform2_C[0])))
	m.T0[v8].(func(*base.Module, int32, int32, int32))(m, l0, l1, l2)
	mBase = m.M
	v10 = int32(4)
	v16 = *(*int32)(unsafe.Add(mBase, uint32(v5)+uint32(_c_F_FTransform2_C[0])))
	m.T0[v16].(func(*base.Module, int32, int32, int32))(m, l0+v10, l1+v10, l2+int32(32))
	mBase = m.M
	return
}
func F_FTransformWHT_C(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	var v44 int32
	_ = v44
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v55 int32
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
	var v63 int32
	_ = v63
	var v64 int32
	_ = v64
	var v66 int32
	_ = v66
	var v67 int32
	_ = v67
	var v70 int32
	_ = v70
	var v74 int32
	_ = v74
	var v78 int32
	_ = v78
	var v80 int32
	_ = v80
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
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v96 int32
	_ = v96
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v102 int32
	_ = v102
	var v106 int32
	_ = v106
	var v110 int32
	_ = v110
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v116 int32
	_ = v116
	var v117 int32
	_ = v117
	var v120 int32
	_ = v120
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v126 int32
	_ = v126
	var v130 int32
	_ = v130
	var v134 int32
	_ = v134
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
	var v144 int32
	_ = v144
	var v146 int32
	_ = v146
	var v147 int32
	_ = v147
	var v150 int32
	_ = v150
	var v154 int32
	_ = v154
	var v158 int32
	_ = v158
	v32 = int32(*(*int16)(unsafe.Add(mBase, uint32(l0)+352)))
	v33 = int32(*(*int16)(unsafe.Add(mBase, uint32(l0)+288)))
	v34 = v32 + v33
	v35 = int32(*(*int16)(unsafe.Add(mBase, uint32(l0)+320)))
	v36 = int32(*(*int16)(unsafe.Add(mBase, uint32(l0)+256)))
	v37 = v35 + v36
	v38 = v34 + v37
	v39 = int32(*(*int16)(unsafe.Add(mBase, uint32(l0)+96)))
	v40 = int32(*(*int16)(unsafe.Add(mBase, uint32(l0)+32)))
	v41 = v39 + v40
	v42 = int32(*(*int16)(unsafe.Add(mBase, uint32(l0)+64)))
	v43 = int32(*(*int16)(unsafe.Add(mBase, uint32(l0))))
	v44 = v42 + v43
	v45 = v41 + v44
	v46 = v38 + v45
	v47 = int32(*(*int16)(unsafe.Add(mBase, uint32(l0)+480)))
	v48 = int32(*(*int16)(unsafe.Add(mBase, uint32(l0)+416)))
	v49 = v47 + v48
	v50 = int32(*(*int16)(unsafe.Add(mBase, uint32(l0)+448)))
	v51 = int32(*(*int16)(unsafe.Add(mBase, uint32(l0)+384)))
	v52 = v50 + v51
	v53 = v49 + v52
	v54 = int32(*(*int16)(unsafe.Add(mBase, uint32(l0)+224)))
	v55 = int32(*(*int16)(unsafe.Add(mBase, uint32(l0)+160)))
	v56 = v54 + v55
	v57 = int32(*(*int16)(unsafe.Add(mBase, uint32(l0)+192)))
	v58 = int32(*(*int16)(unsafe.Add(mBase, uint32(l0)+128)))
	v59 = v57 + v58
	v60 = v56 + v59
	v61 = v53 + v60
	v63 = int32(1)
	v64 = int32(base.Ui32(v46-v61) >> (uint(v63) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l1)+24)) = uint16(v64)
	v66 = v45 - v38
	v67 = v60 - v53
	v70 = int32(base.Ui32(v66-v67) >> (uint(v63) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l1)+16)) = uint16(v70)
	v74 = int32(base.Ui32(v67+v66) >> (uint(v63) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l1)+8)) = uint16(v74)
	v78 = int32(base.Ui32(v61+v46) >> (uint(v63) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l1))) = uint16(v78)
	v80 = v33 - v32
	v81 = v36 - v35
	v82 = v80 + v81
	v83 = v40 - v39
	v84 = v43 - v42
	v85 = v83 + v84
	v86 = v82 + v85
	v87 = v48 - v47
	v88 = v51 - v50
	v89 = v87 + v88
	v90 = v55 - v54
	v91 = v58 - v57
	v92 = v90 + v91
	v93 = v89 + v92
	v96 = int32(base.Ui32(v86-v93) >> (uint(v63) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l1)+26)) = uint16(v96)
	v98 = v85 - v82
	v99 = v92 - v89
	v102 = int32(base.Ui32(v98-v99) >> (uint(v63) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l1)+18)) = uint16(v102)
	v106 = int32(base.Ui32(v99+v98) >> (uint(v63) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l1)+10)) = uint16(v106)
	v110 = int32(base.Ui32(v93+v86) >> (uint(v63) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l1)+2)) = uint16(v110)
	v112 = v81 - v80
	v113 = v84 - v83
	v114 = v112 + v113
	v115 = v88 - v87
	v116 = v91 - v90
	v117 = v115 + v116
	v120 = int32(base.Ui32(v114-v117) >> (uint(v63) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l1)+28)) = uint16(v120)
	v122 = v113 - v112
	v123 = v116 - v115
	v126 = int32(base.Ui32(v122-v123) >> (uint(v63) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l1)+20)) = uint16(v126)
	v130 = int32(base.Ui32(v123+v122) >> (uint(v63) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l1)+12)) = uint16(v130)
	v134 = int32(base.Ui32(v117+v114) >> (uint(v63) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l1)+4)) = uint16(v134)
	v136 = v37 - v34
	v137 = v44 - v41
	v138 = v136 + v137
	v139 = v52 - v49
	v140 = v59 - v56
	v141 = v139 + v140
	v144 = int32(base.Ui32(v138-v141) >> (uint(v63) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l1)+30)) = uint16(v144)
	v146 = v137 - v136
	v147 = v140 - v139
	v150 = int32(base.Ui32(v146-v147) >> (uint(v63) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l1)+22)) = uint16(v150)
	v154 = int32(base.Ui32(v147+v146) >> (uint(v63) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l1)+14)) = uint16(v154)
	v158 = int32(base.Ui32(v141+v138) >> (uint(v63) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l1)+6)) = uint16(v158)
	return
}
func F_FTransform_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	var v44 int32
	_ = v44
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v55 int32
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
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v66 int32
	_ = v66
	var v67 int32
	_ = v67
	var v68 int32
	_ = v68
	var v69 int32
	_ = v69
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v76 int32
	_ = v76
	var v77 int32
	_ = v77
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v84 int32
	_ = v84
	var v85 int32
	_ = v85
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v94 int32
	_ = v94
	var v95 int32
	_ = v95
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v109 int32
	_ = v109
	var v113 int32
	_ = v113
	var v116 int32
	_ = v116
	var v117 int32
	_ = v117
	var v120 int32
	_ = v120
	var v123 int32
	_ = v123
	var v124 int32
	_ = v124
	var v127 int32
	_ = v127
	var v131 int32
	_ = v131
	var v133 int32
	_ = v133
	var v134 int32
	_ = v134
	var v136 int32
	_ = v136
	var v137 int32
	_ = v137
	var v140 int32
	_ = v140
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v145 int32
	_ = v145
	var v148 int32
	_ = v148
	var v155 int32
	_ = v155
	var v157 int32
	_ = v157
	var v162 int32
	_ = v162
	var v167 int32
	_ = v167
	var v169 int32
	_ = v169
	var v172 int32
	_ = v172
	var v176 int32
	_ = v176
	var v178 int32
	_ = v178
	var v179 int32
	_ = v179
	var v180 int32
	_ = v180
	var v183 int32
	_ = v183
	var v190 int32
	_ = v190
	var v192 int32
	_ = v192
	var v193 int32
	_ = v193
	var v194 int32
	_ = v194
	var v197 int32
	_ = v197
	var v204 int32
	_ = v204
	var v205 int32
	_ = v205
	var v208 int32
	_ = v208
	var v215 int32
	_ = v215
	var v216 int32
	_ = v216
	var v219 int32
	_ = v219
	var v223 int32
	_ = v223
	var v235 int32
	_ = v235
	var v242 int32
	_ = v242
	var v245 int32
	_ = v245
	var v254 int32
	_ = v254
	var v257 int32
	_ = v257
	var v266 int32
	_ = v266
	var v275 int32
	_ = v275
	var v276 int32
	_ = v276
	var v279 int32
	_ = v279
	var v283 int32
	_ = v283
	var v285 int32
	_ = v285
	var v288 int32
	_ = v288
	var v295 int32
	_ = v295
	var v297 int32
	_ = v297
	var v300 int32
	_ = v300
	var v307 int32
	_ = v307
	var v319 int32
	_ = v319
	var v331 int32
	_ = v331
	v33 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+99)))
	v34 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+99)))
	v35 = v33 - v34
	v36 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+96)))
	v37 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+96)))
	v38 = v36 - v37
	v39 = v35 + v38
	v40 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+98)))
	v41 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+98)))
	v42 = v40 - v41
	v43 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+97)))
	v44 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+97)))
	v45 = v43 - v44
	v46 = v42 + v45
	v48 = int32(3)
	v49 = (v39 + v46) << (uint(v48) % 32)
	v50 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)))
	v51 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3)))
	v52 = v50 - v51
	v53 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v54 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
	v55 = v53 - v54
	v56 = v52 + v55
	v57 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)))
	v58 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+2)))
	v59 = v57 - v58
	v60 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)))
	v61 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+1)))
	v62 = v60 - v61
	v63 = v59 + v62
	v66 = (v56 + v63) << (uint(v48) % 32)
	v67 = v49 + v66
	v68 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+67)))
	v69 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+67)))
	v70 = v68 - v69
	v71 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+64)))
	v72 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+64)))
	v73 = v71 - v72
	v74 = v70 + v73
	v75 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+66)))
	v76 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+66)))
	v77 = v75 - v76
	v78 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+65)))
	v79 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+65)))
	v80 = v78 - v79
	v81 = v77 + v80
	v84 = (v74 + v81) << (uint(v48) % 32)
	v85 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+35)))
	v86 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+35)))
	v87 = v85 - v86
	v88 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+32)))
	v89 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+32)))
	v90 = v88 - v89
	v91 = v87 + v90
	v92 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+34)))
	v93 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+34)))
	v94 = v92 - v93
	v95 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+33)))
	v96 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+33)))
	v97 = v95 - v96
	v98 = v94 + v97
	v101 = (v91 + v98) << (uint(v48) % 32)
	v102 = v84 + v101
	v104 = int32(4)
	v105 = int32(base.Ui32(v67-v102) >> (uint(v104) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+16)) = uint16(v105)
	v109 = int32(base.Ui32(v67+v102) >> (uint(v104) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l2))) = uint16(v109)
	v113 = (v39 - v46) << (uint(v48) % 32)
	v116 = (v56 - v63) << (uint(v48) % 32)
	v117 = v113 + v116
	v120 = (v74 - v81) << (uint(v48) % 32)
	v123 = (v91 - v98) << (uint(v48) % 32)
	v124 = v120 + v123
	v127 = int32(base.Ui32(v117-v124) >> (uint(v104) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+20)) = uint16(v127)
	v131 = int32(base.Ui32(v117+v124) >> (uint(v104) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+4)) = uint16(v131)
	v133 = v66 - v49
	v134 = int32(2217)
	v136 = v101 - v84
	v137 = int32(-5352)
	v140 = int32(_a_F_FTransform_C_0)
	v142 = int32(16)
	v143 = int32(base.Ui32(v133*v134+v136*v137+v140) >> (uint(v142) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+24)) = uint16(v143)
	v145 = v116 - v113
	v148 = v123 - v120
	v155 = int32(base.Ui32(v145*v134+v148*v137+v140) >> (uint(v142) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+28)) = uint16(v155)
	v157 = int32(_a_F_FTransform_C_1)
	v162 = int32(_a_F_FTransform_C_2)
	v167 = int32(base.Ui32(v133*v157+v136*v134+v162)>>(uint(v142)%32)) + base.B2i32(v66 != v49)
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+8)) = uint16(v167)
	v169 = v45 - v42
	v172 = v38 - v35
	v176 = int32(1812)
	v178 = int32(9)
	v179 = (v169*v134 + v172*v157 + v176) >> (uint(v178) % 32)
	v180 = v62 - v59
	v183 = v55 - v52
	v190 = (v180*v134 + v183*v157 + v176) >> (uint(v178) % 32)
	v192 = int32(7)
	v193 = v179 + v190 + v192
	v194 = v80 - v77
	v197 = v73 - v70
	v204 = (v194*v134 + v197*v157 + v176) >> (uint(v178) % 32)
	v205 = v97 - v94
	v208 = v90 - v87
	v215 = (v205*v134 + v208*v157 + v176) >> (uint(v178) % 32)
	v216 = v204 + v215
	v219 = int32(base.Ui32(v193-v216) >> (uint(v104) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+18)) = uint16(v219)
	v223 = int32(base.Ui32(v193+v216) >> (uint(v104) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+2)) = uint16(v223)
	v235 = int32(base.Ui32(v145*v157+v148*v134+v162)>>(uint(v142)%32)) + base.B2i32(v116 != v113)
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+12)) = uint16(v235)
	v242 = int32(937)
	v245 = (v169*v137 + v172*v134 + v242) >> (uint(v178) % 32)
	v254 = (v180*v137 + v183*v134 + v242) >> (uint(v178) % 32)
	v257 = v245 + v254 + v192
	v266 = (v194*v137 + v197*v134 + v242) >> (uint(v178) % 32)
	v275 = (v205*v137 + v208*v134 + v242) >> (uint(v178) % 32)
	v276 = v266 + v275
	v279 = int32(base.Ui32(v257-v276) >> (uint(v104) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+22)) = uint16(v279)
	v283 = int32(base.Ui32(v257+v276) >> (uint(v104) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+6)) = uint16(v283)
	v285 = v190 - v179
	v288 = v215 - v204
	v295 = int32(base.Ui32(v285*v134+v288*v137+v140) >> (uint(v142) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+26)) = uint16(v295)
	v297 = v254 - v245
	v300 = v275 - v266
	v307 = int32(base.Ui32(v297*v134+v300*v137+v140) >> (uint(v142) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+30)) = uint16(v307)
	v319 = int32(base.Ui32(v285*v157+v288*v134+v162)>>(uint(v142)%32)) + base.B2i32(v190 != v179)
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+10)) = uint16(v319)
	v331 = int32(base.Ui32(v297*v157+v300*v134+v162)>>(uint(v142)%32)) + base.B2i32(v254 != v245)
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+14)) = uint16(v331)
	return
}
func F_FastLog2Slow_C(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v8 int32
	_ = v8
	var v15 int32
	_ = v15
	var v20 int32
	_ = v20
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v34 int64
	_ = v34
	var v38 int64
	_ = v38
	var v43 float64
	_ = v43
	var v47 float64
	_ = v47
	var v55 int32
	_ = v55
	var v59 int32
	_ = v59
	if base.Ui32(int32(_a_F_FastLog2Slow_C_0)) < base.Ui32(l0) {
		v43 = F_log(m, base.F64_convert_i32_u(l0))
		mBase = m.M
		v47 = base.F64_add(base.F64_mul(v43, float64(1.2102203161561485e+07)), float64(0.5))
		if base.F64_lt(v47, float64(4.294967296e+09))&base.F64_ge(v47, float64(0)) == int32(0) {
			v59 = int32(0)
			return v59
		} else {
			v55 = base.I32_trunc_f64_u(v47)
			return v55
		}
	} else {
		v8 = m.G1
		v15 = base.I32_clz(l0) ^ int32(31) + int32(-7)
		v20 = *(*int32)(unsafe.Add(mBase, uint32(v8+int32(_a_F_FastLog2Slow_C_1)+int32(base.Ui32(l0)>>(uint(v15)%32))<<(uint(int32(2))%32))))
		v23 = v20 + v15<<(uint(int32(23))%32)
		if base.Ui32(l0) < base.Ui32(int32(_a_F_FastLog2Slow_C_2)) {
			v59 = v23
			return v59
		} else {
			v26 = int32(-1)
			v34 = base.I64_extend_i32_u(l0)
			v38 = base.I64_div_u_s(base.I64_extend_i32_u((v26<<(uint(v15)%32)^v26)&l0)*int64(12102203)+int64(base.Ui64(v34)>>(uint(int64(1))%64)), v34)
			return v23 + base.I32_wrap_i64(v38)
		}
	}
}
func F_FastSLog2Slow_C(m *base.Module, l0 int32) int64 {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v14 int32
	_ = v14
	var v21 int64
	_ = v21
	var v25 int32
	_ = v25
	var v35 float64
	_ = v35
	var v38 float64
	_ = v38
	var v41 float64
	_ = v41
	var v49 int64
	_ = v49
	if base.Ui32(int32(_a_F_FastSLog2Slow_C_0)) < base.Ui32(l0) {
		v35 = base.F64_convert_i32_u(l0)
		v38 = F_log(m, v35)
		mBase = m.M
		v41 = base.F64_add(base.F64_mul(base.F64_mul(v35, float64(1.2102203161561485e+07)), v38), float64(0.5))
		if base.F64_lt(v41, float64(1.8446744073709552e+19))&base.F64_ge(v41, float64(0)) == int32(0) {
			return int64(0)
		} else {
			v49 = base.I64_trunc_f64_u(v41)
			return v49
		}
	} else {
		v10 = base.I32_clz(l0) ^ int32(31) + int32(-7)
		v14 = m.G1
		v21 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v14+int32(_a_F_FastSLog2Slow_C_1)+int32(base.Ui32(l0)>>(uint(v10)%32))<<(uint(int32(2))%32)))))
		v25 = int32(-1)
		return (base.I64_extend_i32_s(v10)<<(uint(int64(23))%64)+v21)*base.I64_extend_i32_u(l0) + base.I64_extend_i32_u((v25<<(uint(v10)%32)^v25)&l0)*int64(12102203)
	}
}
func F_FilterLoop24_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32) {
	mBase := m.M
	_ = mBase
	var v31 int32
	_ = v31
	var v39 int32
	_ = v39
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v55 int32
	_ = v55
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v63 int32
	_ = v63
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v92 int32
	_ = v92
	var v95 int32
	_ = v95
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v101 int32
	_ = v101
	var v105 int32
	_ = v105
	var v107 int32
	_ = v107
	var v110 int32
	_ = v110
	var v114 int32
	_ = v114
	var v118 int32
	_ = v118
	var v121 int32
	_ = v121
	var v123 int32
	_ = v123
	var v126 int32
	_ = v126
	var v130 int32
	_ = v130
	var v134 int32
	_ = v134
	var v138 int32
	_ = v138
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v146 int32
	_ = v146
	var v149 int32
	_ = v149
	var v156 int32
	_ = v156
	var v158 int32
	_ = v158
	var v161 int32
	_ = v161
	var v166 int32
	_ = v166
	var v172 int32
	_ = v172
	var v173 int32
	_ = v173
	var v176 int32
	_ = v176
	var v179 int32
	_ = v179
	var v183 int32
	_ = v183
	var v187 int32
	_ = v187
	var v190 int32
	_ = v190
	var v191 int32
	_ = v191
	var v197 int32
	_ = v197
	var v208 int32
	_ = v208
	if l3 < int32(1) {
	} else {
		v31 = int32(1)
		v39 = int32(0)
		v49 = l1 << (uint(v31) % 32)
		v51 = m.G14
		v52 = *(*int32)(unsafe.Add(mBase, uint32(v51)))
		v53 = m.G12
		v54 = *(*int32)(unsafe.Add(mBase, uint32(v53)))
		v55 = m.G13
		v56 = *(*int32)(unsafe.Add(mBase, uint32(v55)))
		v57 = m.G15
		v58 = *(*int32)(unsafe.Add(mBase, uint32(v57)))
		v59 = l0
		v63 = l3 + v31
		for {
			v87 = v59 + (v39 - l1)
			v88 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v87))))
			v89 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59))))
			v92 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+(v88-v89)))))
			v95 = v59 + (v39 - v49)
			v96 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v95))))
			v97 = v59 + l1
			v98 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v97))))
			v99 = v96 - v98
			v101 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v99))))
			if l4<<(uint(v31)%32)|v31 < v92<<(uint(int32(2))%32)+v101 {
			} else {
				v105 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+(v39-l1<<(uint(int32(2))%32))))))
				v107 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+l1*int32(-3)))))
				v110 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+(v105-v107)))))
				if l5 < v110 {
				} else {
					v114 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+(v107-v96)))))
					if base.Ui32(l5) < base.Ui32(v114) {
					} else {
						v118 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+(v96-v88)))))
						if base.Ui32(l5) < base.Ui32(v118) {
						} else {
							v121 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+l1*int32(3)))))
							v123 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v49))))
							v126 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+(v121-v123)))))
							if base.Ui32(l5) < base.Ui32(v126) {
							} else {
								v130 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+(v123-v98)))))
								if base.Ui32(l5) < base.Ui32(v130) {
								} else {
									v134 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+(v98-v89)))))
									if base.Ui32(l5) < base.Ui32(v134) {
									} else {
										v138 = (v89 - v88) * int32(3)
										if l6 < v118 {
											v142 = int32(*(*int8)(unsafe.Add(mBase, uint32(v52+v99))))
											v143 = v138 + v142
											v146 = int32(3)
											v149 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+(v143+int32(4))>>(uint(v146)%32)))))
											v156 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+(v143+v146)>>(uint(v146)%32)))))
											v158 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v54+v88+v156))))
											*(*uint8)(unsafe.Add(mBase, uint32(v87))) = uint8(v158)
											v190 = v89 - v149
											v191 = v59
										} else {
											if base.Ui32(v134) <= base.Ui32(l6) {
												v161 = int32(3)
												v166 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+(v138+v161)>>(uint(v161)%32)))))
												v172 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+(v138+int32(4))>>(uint(v161)%32)))))
												v173 = int32(1)
												v176 = (v172 + v173) >> (uint(v173) % 32)
												v179 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v54+v176+v96))))
												*(*uint8)(unsafe.Add(mBase, uint32(v95))) = uint8(v179)
												v183 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v166+(v54+v88)))))
												*(*uint8)(unsafe.Add(mBase, uint32(v87))) = uint8(v183)
												v187 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v54+(v89-v172)))))
												*(*uint8)(unsafe.Add(mBase, uint32(v59))) = uint8(v187)
												v190 = v98 - v176
												v191 = v97
											} else {
												v142 = int32(*(*int8)(unsafe.Add(mBase, uint32(v52+v99))))
												v143 = v138 + v142
												v146 = int32(3)
												v149 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+(v143+int32(4))>>(uint(v146)%32)))))
												v156 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+(v143+v146)>>(uint(v146)%32)))))
												v158 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v54+v88+v156))))
												*(*uint8)(unsafe.Add(mBase, uint32(v87))) = uint8(v158)
												v190 = v89 - v149
												v191 = v59
											}
										}
										v197 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v54+v190))))
										*(*uint8)(unsafe.Add(mBase, uint32(v191))) = uint8(v197)
									}
								}
							}
						}
					}
				}
			}
			v208 = v63 + int32(-1)
			if base.Ui32(int32(1)) < base.Ui32(v208) {
				v59 = v59 + l2
				v63 = v208
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_FilterLoop26_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32) {
	mBase := m.M
	_ = mBase
	var v35 int32
	_ = v35
	var v43 int32
	_ = v43
	var v53 int32
	_ = v53
	var v55 int32
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
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v67 int32
	_ = v67
	var v95 int32
	_ = v95
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v100 int32
	_ = v100
	var v103 int32
	_ = v103
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v107 int32
	_ = v107
	var v109 int32
	_ = v109
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v118 int32
	_ = v118
	var v122 int32
	_ = v122
	var v126 int32
	_ = v126
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v134 int32
	_ = v134
	var v138 int32
	_ = v138
	var v142 int32
	_ = v142
	var v146 int32
	_ = v146
	var v148 int32
	_ = v148
	var v151 int32
	_ = v151
	var v154 int32
	_ = v154
	var v157 int32
	_ = v157
	var v164 int32
	_ = v164
	var v166 int32
	_ = v166
	var v171 int32
	_ = v171
	var v174 int32
	_ = v174
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v180 int32
	_ = v180
	var v187 int32
	_ = v187
	var v190 int32
	_ = v190
	var v197 int32
	_ = v197
	var v200 int32
	_ = v200
	var v204 int32
	_ = v204
	var v208 int32
	_ = v208
	var v211 int32
	_ = v211
	var v216 int32
	_ = v216
	var v218 int32
	_ = v218
	var v232 int32
	_ = v232
	if l3 < int32(1) {
	} else {
		v35 = int32(1)
		v43 = int32(0)
		v53 = l1 << (uint(v35) % 32)
		v55 = m.G13
		v56 = *(*int32)(unsafe.Add(mBase, uint32(v55)))
		v57 = m.G12
		v58 = *(*int32)(unsafe.Add(mBase, uint32(v57)))
		v59 = m.G14
		v60 = *(*int32)(unsafe.Add(mBase, uint32(v59)))
		v61 = m.G15
		v62 = *(*int32)(unsafe.Add(mBase, uint32(v61)))
		v63 = l0
		v67 = l3 + v35
		for {
			v95 = v63 + (v43 - l1)
			v96 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v95))))
			v97 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63))))
			v100 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v96-v97)))))
			v103 = v63 + (v43 - v53)
			v104 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v103))))
			v105 = v63 + l1
			v106 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v105))))
			v107 = v104 - v106
			v109 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+v107))))
			if l4<<(uint(v35)%32)|v35 < v100<<(uint(int32(2))%32)+v109 {
			} else {
				v113 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v43-l1<<(uint(int32(2))%32))))))
				v114 = v63 + l1*int32(-3)
				v115 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v114))))
				v118 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v113-v115)))))
				if l5 < v118 {
				} else {
					v122 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v115-v104)))))
					if base.Ui32(l5) < base.Ui32(v122) {
					} else {
						v126 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v104-v96)))))
						if base.Ui32(l5) < base.Ui32(v126) {
						} else {
							v129 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+l1*int32(3)))))
							v130 = v63 + v53
							v131 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v130))))
							v134 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v129-v131)))))
							if base.Ui32(l5) < base.Ui32(v134) {
							} else {
								v138 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v131-v106)))))
								if base.Ui32(l5) < base.Ui32(v138) {
								} else {
									v142 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v106-v97)))))
									if base.Ui32(l5) < base.Ui32(v142) {
									} else {
										v146 = (v97 - v96) * int32(3)
										v148 = int32(*(*int8)(unsafe.Add(mBase, uint32(v60+v107))))
										if l6 < v126 {
											v151 = v146 + v148
											v154 = int32(3)
											v157 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+(v151+int32(4))>>(uint(v154)%32)))))
											v164 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+(v151+v154)>>(uint(v154)%32)))))
											v166 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v96+v164))))
											*(*uint8)(unsafe.Add(mBase, uint32(v95))) = uint8(v166)
											v211 = v97 - v157
											v216 = v63
										} else {
											if base.Ui32(v142) <= base.Ui32(l6) {
												v171 = int32(*(*int8)(unsafe.Add(mBase, uint32(v60+v146+v148))))
												v174 = int32(63)
												v176 = int32(7)
												v177 = (v171*int32(9) + v174) >> (uint(v176) % 32)
												v180 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v177+v115))))
												*(*uint8)(unsafe.Add(mBase, uint32(v114))) = uint8(v180)
												v187 = (v171*int32(18) + v174) >> (uint(v176) % 32)
												v190 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v187+v104))))
												*(*uint8)(unsafe.Add(mBase, uint32(v103))) = uint8(v190)
												v197 = (v171*int32(27) + v174) >> (uint(v176) % 32)
												v200 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v197+v96))))
												*(*uint8)(unsafe.Add(mBase, uint32(v95))) = uint8(v200)
												v204 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+(v97-v197)))))
												*(*uint8)(unsafe.Add(mBase, uint32(v63))) = uint8(v204)
												v208 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+(v106-v187)))))
												*(*uint8)(unsafe.Add(mBase, uint32(v105))) = uint8(v208)
												v211 = v131 - v177
												v216 = v130
											} else {
												v151 = v146 + v148
												v154 = int32(3)
												v157 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+(v151+int32(4))>>(uint(v154)%32)))))
												v164 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+(v151+v154)>>(uint(v154)%32)))))
												v166 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v96+v164))))
												*(*uint8)(unsafe.Add(mBase, uint32(v95))) = uint8(v166)
												v211 = v97 - v157
												v216 = v63
											}
										}
										v218 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v211))))
										*(*uint8)(unsafe.Add(mBase, uint32(v216))) = uint8(v218)
									}
								}
							}
						}
					}
				}
			}
			v232 = v67 + int32(-1)
			if base.Ui32(int32(1)) < base.Ui32(v232) {
				v63 = v63 + l2
				v67 = v232
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_FinalizeTokenProbas(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v39 int32
	_ = v39
	var v44 int32
	_ = v44
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v83 int32
	_ = v83
	var v84 int32
	_ = v84
	var v85 int32
	_ = v85
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v119 int32
	_ = v119
	var v120 int32
	_ = v120
	var v130 int32
	_ = v130
	var v132 int32
	_ = v132
	var v133 int32
	_ = v133
	var v135 int32
	_ = v135
	var v136 int32
	_ = v136
	var v138 int32
	_ = v138
	var v140 int32
	_ = v140
	var v142 int32
	_ = v142
	var v145 int32
	_ = v145
	var v147 int32
	_ = v147
	var v148 int32
	_ = v148
	var v151 int32
	_ = v151
	var v152 int32
	_ = v152
	var v155 int32
	_ = v155
	var v157 int32
	_ = v157
	var v162 int32
	_ = v162
	var v168 int32
	_ = v168
	var v169 int32
	_ = v169
	var v171 int32
	_ = v171
	var v177 int32
	_ = v177
	var v182 int32
	_ = v182
	var v190 int32
	_ = v190
	var v193 int32
	_ = v193
	var v195 int32
	_ = v195
	var v202 int32
	_ = v202
	var v203 int32
	_ = v203
	var v217 int32
	_ = v217
	var v218 int32
	_ = v218
	var v222 int32
	_ = v222
	var v234 int32
	_ = v234
	var v235 int32
	_ = v235
	var v244 int32
	_ = v244
	var v250 int32
	_ = v250
	var v254 int32
	_ = v254
	var v256 int32
	_ = v256
	var v257 int32
	_ = v257
	var v258 int32
	_ = v258
	var v260 int32
	_ = v260
	var v262 int32
	_ = v262
	var v263 int32
	_ = v263
	var v267 int32
	_ = v267
	var v269 int32
	_ = v269
	var v271 int32
	_ = v271
	var v274 int32
	_ = v274
	var v276 int32
	_ = v276
	var v277 int32
	_ = v277
	var v280 int32
	_ = v280
	var v281 int32
	_ = v281
	var v284 int32
	_ = v284
	var v286 int32
	_ = v286
	var v291 int32
	_ = v291
	var v297 int32
	_ = v297
	var v300 int32
	_ = v300
	var v306 int32
	_ = v306
	var v311 int32
	_ = v311
	var v319 int32
	_ = v319
	var v323 int32
	_ = v323
	var v324 int32
	_ = v324
	var v331 int32
	_ = v331
	var v332 int32
	_ = v332
	var v341 int32
	_ = v341
	var v342 int32
	_ = v342
	var v346 int32
	_ = v346
	var v358 int32
	_ = v358
	var v359 int32
	_ = v359
	var v368 int32
	_ = v368
	var v374 int32
	_ = v374
	var v378 int32
	_ = v378
	var v380 int32
	_ = v380
	var v381 int32
	_ = v381
	var v382 int32
	_ = v382
	var v384 int32
	_ = v384
	var v386 int32
	_ = v386
	var v387 int32
	_ = v387
	var v391 int32
	_ = v391
	var v393 int32
	_ = v393
	var v395 int32
	_ = v395
	var v398 int32
	_ = v398
	var v400 int32
	_ = v400
	var v401 int32
	_ = v401
	var v404 int32
	_ = v404
	var v405 int32
	_ = v405
	var v408 int32
	_ = v408
	var v410 int32
	_ = v410
	var v415 int32
	_ = v415
	var v421 int32
	_ = v421
	var v424 int32
	_ = v424
	var v430 int32
	_ = v430
	var v435 int32
	_ = v435
	var v443 int32
	_ = v443
	var v447 int32
	_ = v447
	var v448 int32
	_ = v448
	var v455 int32
	_ = v455
	var v456 int32
	_ = v456
	var v465 int32
	_ = v465
	var v466 int32
	_ = v466
	var v470 int32
	_ = v470
	var v473 int32
	_ = v473
	var v475 int32
	_ = v475
	var v486 int32
	_ = v486
	var v489 int32
	_ = v489
	var v491 int32
	_ = v491
	var v502 int32
	_ = v502
	v2 = int32(0)
	v39 = l0 + int32(1060)
	v44 = l0 + int32(1148)
	v45 = l0 + int32(26)
	v46 = l0 + int32(1104)
	v47 = l0 + int32(15)
	v48 = v2
	v50 = v39
	v51 = v2
	v52 = v2
	v53 = v2
	for {
		v80 = v51
		v81 = v52
		v83 = v44
		v84 = v45
		v85 = v46
		v86 = v47
		v87 = v48
		v88 = v50
		v89 = int32(0)
		for {
			v109 = v80
			v110 = v81
			v119 = int32(0)
			v120 = v87
			for {
				v130 = *(*int32)(unsafe.Add(mBase, uint32(v88+v119)))
				v132 = int32(base.Ui32(v130) >> (uint(int32(16)) % 32))
				v133 = m.G80
				v135 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v133+v120))))
				v136 = m.G81
				v138 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v136+v120))))
				v140 = v130 & int32(_a_F_FinalizeTokenProbas_0)
				if v140 != 0 {
					v142 = int32(255)
					v145 = base.I32_div_u_s(v140*v142, v132)
					v147 = v142 - v145
				} else {
					v147 = int32(255)
				}
				v148 = m.G24
				v151 = v132 - v140
				v152 = int32(1)
				v155 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v148+v135<<(uint(v152)%32)))))
				v157 = int32(255)
				v162 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v148+(v135^v157)<<(uint(v152)%32)))))
				v168 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v148+v138<<(uint(v152)%32)))))
				v169 = v151*v155 + v140*v162 + v168
				v171 = v147 & v157
				v177 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v148+(v171^v157)<<(uint(v152)%32)))))
				v182 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v148+v171<<(uint(v152)%32)))))
				v190 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v148+(v138^v157)<<(uint(v152)%32)))))
				v193 = v140*v177 + v151*v182 + v190 + int32(2048)
				if v193 < v169 {
					v195 = int32(-1)
				} else {
					v195 = int32(0)
				}
				v202 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v148+(v138^v195)&int32(255)<<(uint(int32(1))%32)))))
				v203 = v110 + v202
				if v169 <= v193 {
					*(*uint8)(unsafe.Add(mBase, uint32(v39+v120+int32(-1056)))) = uint8(v135)
					v217 = v109
					v218 = v203
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v39+v120+int32(-1056)))) = uint8(v147)
					v217 = v109 | base.B2i32(v147 != v135)
					v218 = v203 + int32(2048)
				}
				v222 = v119 + int32(4)
				if v222 != int32(44) {
					v109 = v217
					v110 = v218
					v119 = v222
					v120 = v120 + int32(1)
					continue
				} else {
					break
				}
				break
			}
			v234 = v217
			v235 = v218
			v244 = v85
			v250 = int32(0)
			for {
				v254 = *(*int32)(unsafe.Add(mBase, uint32(v244)))
				v256 = int32(base.Ui32(v254) >> (uint(int32(16)) % 32))
				v257 = m.G80
				v258 = v87 + v250
				v260 = int32(11)
				v262 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v257+v258+v260))))
				v263 = m.G81
				v267 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v263+v258+v260))))
				v269 = v254 & int32(_a_F_FinalizeTokenProbas_0)
				if v269 != 0 {
					v271 = int32(255)
					v274 = base.I32_div_u_s(v269*v271, v256)
					v276 = v271 - v274
				} else {
					v276 = int32(255)
				}
				v277 = m.G24
				v280 = v256 - v269
				v281 = int32(1)
				v284 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v277+v262<<(uint(v281)%32)))))
				v286 = int32(255)
				v291 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v277+(v262^v286)<<(uint(v281)%32)))))
				v297 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v277+v267<<(uint(v281)%32)))))
				v300 = v276 & v286
				v306 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v277+(v300^v286)<<(uint(v281)%32)))))
				v311 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v277+v300<<(uint(v281)%32)))))
				v319 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v277+(v267^v286)<<(uint(v281)%32)))))
				v323 = base.B2i32(v269*v306+v280*v311+v319+int32(2048) < v280*v284+v269*v291+v297)
				if v269*v306+v280*v311+v319+int32(2048) < v280*v284+v269*v291+v297 {
					v324 = int32(-1)
				} else {
					v324 = int32(0)
				}
				v331 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v277+(v267^v324)&int32(255)<<(uint(int32(1))%32)))))
				v332 = v235 + v331
				if v269*v306+v280*v311+v319+int32(2048) < v280*v284+v269*v291+v297 {
					*(*uint8)(unsafe.Add(mBase, uint32(v86+v250))) = uint8(v276)
					v341 = v234 | base.B2i32(v276 != v262)
					v342 = v332 + int32(2048)
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v86+v250))) = uint8(v262)
					v341 = v234
					v342 = v332
				}
				v346 = v250 + int32(1)
				if v346 != int32(11) {
					v234 = v341
					v235 = v342
					v244 = v244 + int32(4)
					v250 = v346
					continue
				} else {
					break
				}
				break
			}
			v358 = v341
			v359 = v342
			v368 = v83
			v374 = int32(0)
			for {
				v378 = *(*int32)(unsafe.Add(mBase, uint32(v368)))
				v380 = int32(base.Ui32(v378) >> (uint(int32(16)) % 32))
				v381 = m.G80
				v382 = v87 + v374
				v384 = int32(22)
				v386 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v381+v382+v384))))
				v387 = m.G81
				v391 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v387+v382+v384))))
				v393 = v378 & int32(_a_F_FinalizeTokenProbas_0)
				if v393 != 0 {
					v395 = int32(255)
					v398 = base.I32_div_u_s(v393*v395, v380)
					v400 = v395 - v398
				} else {
					v400 = int32(255)
				}
				v401 = m.G24
				v404 = v380 - v393
				v405 = int32(1)
				v408 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v401+v386<<(uint(v405)%32)))))
				v410 = int32(255)
				v415 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v401+(v386^v410)<<(uint(v405)%32)))))
				v421 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v401+v391<<(uint(v405)%32)))))
				v424 = v400 & v410
				v430 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v401+(v424^v410)<<(uint(v405)%32)))))
				v435 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v401+v424<<(uint(v405)%32)))))
				v443 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v401+(v391^v410)<<(uint(v405)%32)))))
				v447 = base.B2i32(v393*v430+v404*v435+v443+int32(2048) < v404*v408+v393*v415+v421)
				if v393*v430+v404*v435+v443+int32(2048) < v404*v408+v393*v415+v421 {
					v448 = int32(-1)
				} else {
					v448 = int32(0)
				}
				v455 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v401+(v391^v448)&int32(255)<<(uint(int32(1))%32)))))
				v456 = v359 + v455
				if v393*v430+v404*v435+v443+int32(2048) < v404*v408+v393*v415+v421 {
					*(*uint8)(unsafe.Add(mBase, uint32(v84+v374))) = uint8(v400)
					v465 = v358 | base.B2i32(v400 != v386)
					v466 = v456 + int32(2048)
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v84+v374))) = uint8(v386)
					v465 = v358
					v466 = v456
				}
				v470 = v374 + int32(1)
				if v470 != int32(11) {
					v358 = v465
					v359 = v466
					v368 = v368 + int32(4)
					v374 = v470
					continue
				} else {
					break
				}
				break
			}
			v473 = int32(132)
			v475 = int32(33)
			v486 = v89 + int32(1)
			if v486 != int32(8) {
				v80 = v465
				v81 = v466
				v83 = v83 + v473
				v84 = v84 + v475
				v85 = v85 + v473
				v86 = v86 + v475
				v87 = v87 + v475
				v88 = v88 + v473
				v89 = v486
				continue
			} else {
				break
			}
			break
		}
		v489 = int32(1056)
		v491 = int32(264)
		v502 = v53 + int32(1)
		if v502 != int32(4) {
			v44 = v44 + v489
			v45 = v45 + v491
			v46 = v46 + v489
			v47 = v47 + v491
			v48 = v48 + v491
			v50 = v50 + v489
			v51 = v465
			v52 = v466
			v53 = v502
			continue
		} else {
			break
		}
		break
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_FinalizeTokenProbas[0]))) = v465
	return v466
}
func F_Flush(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v8 int32
	_ = v8
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v19 int32
	_ = v19
	var v23 int32
	_ = v23
	var v27 int32
	_ = v27
	var v29 int64
	_ = v29
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v37 int64
	_ = v37
	var v39 int32
	_ = v39
	var v41 int32
	_ = v41
	var v42 int32
	_ = v42
	var v45 int32
	_ = v45
	var v52 int64
	_ = v52
	var v53 int32
	_ = v53
	var v54 int64
	_ = v54
	var v64 int32
	_ = v64
	var v66 int32
	_ = v66
	var v70 int32
	_ = v70
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v83 int32
	_ = v83
	var v88 int32
	_ = v88
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v94 int32
	_ = v94
	var v97 int32
	_ = v97
	var v102 int32
	_ = v102
	var v104 int32
	_ = v104
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v117 int32
	_ = v117
	var v118 int32
	_ = v118
	var v122 int32
	_ = v122
	var v128 int32
	_ = v128
	var v134 int32
	_ = v134
	v8 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v8 + int32(-8)
	v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v14 = v8 + int32(8)
	v15 = v12 >> (uint(v14) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v12 - v15<<(uint(v14)%32)
	v19 = int32(255)
	if v15&v19 == v19 {
		v134 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
		*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v134 + int32(1)
		return
	} else {
		v23 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
		v27 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
		v29 = base.I64_extend_i32_u(v23+int32(1)) + base.I64_extend_i32_u(v27)
		if base.Ui64(v29) < base.Ui64(int64(4294967296)) {
			v34 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
			v35 = base.I32_wrap_i64(v29)
			if base.Ui32(v35) <= base.Ui32(v34) {
				v83 = v15 & int32(256)
				if v83 == int32(0) {
				} else {
					if v27 == int32(0) {
					} else {
						v88 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
						v91 = v88 + v27 + int32(-1)
						v92 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v91))))
						v94 = v92 + int32(1)
						*(*uint8)(unsafe.Add(mBase, uint32(v91))) = uint8(v94)
					}
				}
				v97 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				if v97 < int32(1) {
					v122 = v27
				} else {
					if v83 != 0 {
						v102 = int32(0)
					} else {
						v102 = int32(-1)
					}
					v104 = v27
					for {
						v110 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
						*(*uint8)(unsafe.Add(mBase, uint32(v110+v104))) = uint8(v102)
						v113 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
						*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v113 + int32(-1)
						v117 = int32(1)
						v118 = v104 + v117
						if v117 < v113 {
							v104 = v118
							continue
						} else {
							break
						}
						break
					}
					v122 = v118
				}
				v128 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
				*(*uint8)(unsafe.Add(mBase, uint32(v128+v122))) = uint8(v15)
				*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v122 + int32(1)
				return
			} else {
				v37 = int64(1)
				v39 = v34 << (uint(int32(1)) % 32)
				if base.Ui32(v35) < base.Ui32(v39) {
					v41 = v39
				} else {
					v41 = v35
				}
				v42 = int32(1024)
				if base.Ui32(v42) < base.Ui32(v41) {
					v45 = v41
				} else {
					v45 = v42
				}
				v52 = base.I64_div_u_s(int64(2147418112), v37)
				v53 = int32(0)
				v54 = base.I64_extend_i32_u(v45)
				if base.Ui64(int64(4294967295)) < base.Ui64(v54*v37) {
					v66 = v53
				} else {
					if base.Ui64(v52) < base.Ui64(v54) {
						v66 = v53
					} else {
						v64 = F_malloc(m, base.I32_wrap_i64(v37)*v45)
						mBase = m.M
						v66 = v64
					}
				}
				if v66 != 0 {
					v70 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
					if v70 == int32(0) {
					} else {
						v73 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
						v74 = F_memcpy(m, v66, v73, v70)
						mBase = m.M
					}
					v75 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
					F_free(m, v75)
					mBase = m.M
					*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v45
					*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v66
					v83 = v15 & int32(256)
					if v83 == int32(0) {
					} else {
						if v27 == int32(0) {
						} else {
							v88 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
							v91 = v88 + v27 + int32(-1)
							v92 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v91))))
							v94 = v92 + int32(1)
							*(*uint8)(unsafe.Add(mBase, uint32(v91))) = uint8(v94)
						}
					}
					v97 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					if v97 < int32(1) {
						v122 = v27
					} else {
						if v83 != 0 {
							v102 = int32(0)
						} else {
							v102 = int32(-1)
						}
						v104 = v27
						for {
							v110 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
							*(*uint8)(unsafe.Add(mBase, uint32(v110+v104))) = uint8(v102)
							v113 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v113 + int32(-1)
							v117 = int32(1)
							v118 = v104 + v117
							if v117 < v113 {
								v104 = v118
								continue
							} else {
								break
							}
							break
						}
						v122 = v118
					}
					v128 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
					*(*uint8)(unsafe.Add(mBase, uint32(v128+v122))) = uint8(v15)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v122 + int32(1)
					return
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
					return
				}
			}
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
			return
		}
	}
}
func F_free(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v15 int32
	_ = v15
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v32 int32
	_ = v32
	var v34 int32
	_ = v34
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v41 int32
	_ = v41
	var v48 int32
	_ = v48
	var v50 int32
	_ = v50
	var v53 int32
	_ = v53
	var v58 int32
	_ = v58
	var v63 int32
	_ = v63
	var v64 int32
	_ = v64
	var v69 int32
	_ = v69
	var __phi69 int32
	_ = __phi69
	var v70 int32
	_ = v70
	var __phi70 int32
	_ = __phi70
	var v75 int32
	_ = v75
	var v78 int32
	_ = v78
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v100 int32
	_ = v100
	var v108 int32
	_ = v108
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v116 int32
	_ = v116
	var v118 int32
	_ = v118
	var v125 int32
	_ = v125
	var v127 int32
	_ = v127
	var v133 int32
	_ = v133
	var v138 int32
	_ = v138
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v152 int32
	_ = v152
	var v160 int32
	_ = v160
	var v162 int32
	_ = v162
	var v166 int32
	_ = v166
	var v167 int32
	_ = v167
	var v173 int32
	_ = v173
	var v175 int32
	_ = v175
	var v182 int32
	_ = v182
	var v184 int32
	_ = v184
	var v188 int32
	_ = v188
	var v189 int32
	_ = v189
	var v198 int32
	_ = v198
	var v199 int32
	_ = v199
	var v202 int32
	_ = v202
	var v204 int32
	_ = v204
	var v206 int32
	_ = v206
	var v215 int32
	_ = v215
	var v217 int32
	_ = v217
	var v220 int32
	_ = v220
	var v225 int32
	_ = v225
	var v230 int32
	_ = v230
	var v231 int32
	_ = v231
	var v236 int32
	_ = v236
	var __phi236 int32
	_ = __phi236
	var v237 int32
	_ = v237
	var __phi237 int32
	_ = __phi237
	var v242 int32
	_ = v242
	var v245 int32
	_ = v245
	var v259 int32
	_ = v259
	var v267 int32
	_ = v267
	var v269 int32
	_ = v269
	var v272 int32
	_ = v272
	var v275 int32
	_ = v275
	var v277 int32
	_ = v277
	var v284 int32
	_ = v284
	var v286 int32
	_ = v286
	var v292 int32
	_ = v292
	var v297 int32
	_ = v297
	var v316 int32
	_ = v316
	var v320 int32
	_ = v320
	var v331 int32
	_ = v331
	var v333 int32
	_ = v333
	var v335 int32
	_ = v335
	var v339 int32
	_ = v339
	var v344 int32
	_ = v344
	var v345 int32
	_ = v345
	var v356 int32
	_ = v356
	var v359 int32
	_ = v359
	var v366 int32
	_ = v366
	var v371 int32
	_ = v371
	var v373 int32
	_ = v373
	var v375 int32
	_ = v375
	var v377 int32
	_ = v377
	var v391 int32
	_ = v391
	var v393 int32
	_ = v393
	var v396 int32
	_ = v396
	var v399 int32
	_ = v399
	var v402 int32
	_ = v402
	var v414 int32
	_ = v414
	var v415 int32
	_ = v415
	var v418 int32
	_ = v418
	var v420 int32
	_ = v420
	var v421 int32
	_ = v421
	var v423 int32
	_ = v423
	var v426 int32
	_ = v426
	var v428 int32
	_ = v428
	var v433 int32
	_ = v433
	var v435 int32
	_ = v435
	var v436 int32
	_ = v436
	var v437 int32
	_ = v437
	var v438 int32
	_ = v438
	var v440 int32
	_ = v440
	var v447 int32
	_ = v447
	var v449 int32
	_ = v449
	var v450 int32
	_ = v450
	var v451 int32
	_ = v451
	var v453 int32
	_ = v453
	if l0 == int32(0) {
		goto L2
	} else {
		goto L3
	}
L1:
	;
	return
L2:
	;
	goto L1
L3:
	;
	v11 = int32(-8)
	v12 = l0 + v11
	v15 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(-4))))
	v17 = v15 & v11
	v18 = v12 + v17
	if v15&int32(1) != 0 {
		v143 = v17
		v144 = v12
		goto L4
	} else {
		goto L5
	}
L4:
	;
	if base.Ui32(v18) <= base.Ui32(v144) {
		goto L2
	} else {
		goto L39
	}
L5:
	;
	if v15&int32(2) == int32(0) {
		goto L2
	} else {
		goto L6
	}
L6:
	;
	v25 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
	v26 = v12 - v25
	v28 = *(*int32)(unsafe.Add(mBase, _c_F_free[0]))
	if base.Ui32(v26) < base.Ui32(v28) {
		goto L2
	} else {
		goto L7
	}
L7:
	;
	v30 = v25 + v17
	v32 = *(*int32)(unsafe.Add(mBase, _c_F_free[1]))
	if v26 == v32 {
		goto L11
	} else {
		goto L12
	}
L8:
	;
	if v48 == int32(0) {
		v143 = v30
		v144 = v26
		goto L4
	} else {
		goto L27
	}
L9:
	;
	v100 = int32(0)
	goto L8
L10:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34)+8)) = v37
	*(*int32)(unsafe.Add(mBase, uint32(v37)+12)) = v34
	v143 = v30
	v144 = v26
	goto L4
L11:
	;
	v81 = *(*int32)(unsafe.Add(mBase, uint32(v18)+4))
	v82 = int32(3)
	if v81&v82 != v82 {
		v143 = v30
		v144 = v26
		goto L4
	} else {
		goto L26
	}
L12:
	;
	v34 = *(*int32)(unsafe.Add(mBase, uint32(v26)+12))
	if base.Ui32(int32(255)) < base.Ui32(v25) {
		goto L13
	} else {
		goto L14
	}
L13:
	;
	v48 = *(*int32)(unsafe.Add(mBase, uint32(v26)+24))
	if v34 == v26 {
		goto L16
	} else {
		goto L17
	}
L14:
	;
	v37 = *(*int32)(unsafe.Add(mBase, uint32(v26)+8))
	if v34 != v37 {
		goto L10
	} else {
		goto L15
	}
L15:
	;
	v39 = int32(0)
	v41 = *(*int32)(unsafe.Add(mBase, _c_F_free[2]))
	*(*int32)(unsafe.Add(mBase, _c_F_free[2])) = v41 & base.I32_rotl(int32(-2), int32(base.Ui32(v25)>>(uint(int32(3))%32)))
	v143 = v30
	v144 = v26
	goto L4
L16:
	;
	v53 = *(*int32)(unsafe.Add(mBase, uint32(v26)+20))
	if v53 == int32(0) {
		goto L19
	} else {
		goto L20
	}
L17:
	;
	v50 = *(*int32)(unsafe.Add(mBase, uint32(v26)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v50)+12)) = v34
	*(*int32)(unsafe.Add(mBase, uint32(v34)+8)) = v50
	v100 = v34
	goto L8
L18:
	;
	__phi69 = v63
	__phi70 = v64
	v69 = __phi69
	v70 = __phi70
	goto L22
L19:
	;
	v58 = *(*int32)(unsafe.Add(mBase, uint32(v26)+16))
	if v58 == int32(0) {
		goto L9
	} else {
		goto L21
	}
L20:
	;
	v63 = v53
	v64 = v26 + int32(20)
	goto L18
L21:
	;
	v63 = v58
	v64 = v26 + int32(16)
	goto L18
L22:
	;
	v75 = *(*int32)(unsafe.Add(mBase, uint32(v69)+20))
	if v75 != 0 {
		__phi69 = v75
		__phi70 = v69 + int32(20)
		v69 = __phi69
		v70 = __phi70
		goto L22
	} else {
		goto L24
	}
L23:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v70))) = int32(0)
	v100 = v69
	goto L8
L24:
	;
	v78 = *(*int32)(unsafe.Add(mBase, uint32(v69)+16))
	if v78 != 0 {
		__phi69 = v78
		__phi70 = v69 + int32(16)
		v69 = __phi69
		v70 = __phi70
		goto L22
	} else {
		goto L25
	}
L25:
	;
	goto L23
L26:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v18)+4)) = v81 & int32(-2)
	*(*int32)(unsafe.Add(mBase, _c_F_free[3])) = v30
	*(*int32)(unsafe.Add(mBase, uint32(v18))) = v30
	*(*int32)(unsafe.Add(mBase, uint32(v26)+4)) = v30 | int32(1)
	goto L1
L27:
	;
	v108 = *(*int32)(unsafe.Add(mBase, uint32(v26)+28))
	v110 = v108 << (uint(int32(2)) % 32)
	v113 = *(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_free[4])))
	if v26 != v113 {
		goto L29
	} else {
		goto L30
	}
L28:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v100)+24)) = v48
	v133 = *(*int32)(unsafe.Add(mBase, uint32(v26)+16))
	if v133 == int32(0) {
		goto L36
	} else {
		goto L37
	}
L29:
	;
	v125 = *(*int32)(unsafe.Add(mBase, uint32(v48)+16))
	if v125 == v26 {
		goto L32
	} else {
		goto L33
	}
L30:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_free[4]))) = v100
	if v100 != 0 {
		goto L28
	} else {
		goto L31
	}
L31:
	;
	v116 = int32(0)
	v118 = *(*int32)(unsafe.Add(mBase, _c_F_free[5]))
	*(*int32)(unsafe.Add(mBase, _c_F_free[5])) = v118 & base.I32_rotl(int32(-2), v108)
	v143 = v30
	v144 = v26
	goto L4
L32:
	;
	v127 = int32(16)
	goto L34
L33:
	;
	v127 = int32(20)
	goto L34
L34:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v48+v127))) = v100
	if v100 == int32(0) {
		v143 = v30
		v144 = v26
		goto L4
	} else {
		goto L35
	}
L35:
	;
	goto L28
L36:
	;
	v138 = *(*int32)(unsafe.Add(mBase, uint32(v26)+20))
	if v138 == int32(0) {
		v143 = v30
		v144 = v26
		goto L4
	} else {
		goto L38
	}
L37:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v100)+16)) = v133
	*(*int32)(unsafe.Add(mBase, uint32(v133)+24)) = v100
	goto L36
L38:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v100)+20)) = v138
	*(*int32)(unsafe.Add(mBase, uint32(v138)+24)) = v100
	v143 = v30
	v144 = v26
	goto L4
L39:
	;
	v152 = *(*int32)(unsafe.Add(mBase, uint32(v18)+4))
	if v152&int32(1) == int32(0) {
		goto L2
	} else {
		goto L40
	}
L40:
	;
	if v152&int32(2) != 0 {
		goto L45
	} else {
		goto L46
	}
L41:
	;
	if base.Ui32(int32(255)) < base.Ui32(v320) {
		goto L79
	} else {
		goto L80
	}
L42:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v144+v198))) = v198
	*(*int32)(unsafe.Add(mBase, uint32(v144)+4)) = v198 | int32(1)
	v316 = *(*int32)(unsafe.Add(mBase, _c_F_free[1]))
	if v144 != v316 {
		v320 = v198
		goto L41
	} else {
		goto L78
	}
L43:
	;
	if v215 == int32(0) {
		goto L42
	} else {
		goto L66
	}
L44:
	;
	v259 = int32(0)
	goto L43
L45:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v18)+4)) = v152 & int32(-2)
	*(*int32)(unsafe.Add(mBase, uint32(v144+v143))) = v143
	*(*int32)(unsafe.Add(mBase, uint32(v144)+4)) = v143 | int32(1)
	v320 = v143
	goto L41
L46:
	;
	v160 = *(*int32)(unsafe.Add(mBase, _c_F_free[6]))
	if v18 != v160 {
		goto L47
	} else {
		goto L48
	}
L47:
	;
	v182 = *(*int32)(unsafe.Add(mBase, _c_F_free[1]))
	if v18 != v182 {
		goto L50
	} else {
		goto L51
	}
L48:
	;
	v162 = int32(0)
	*(*int32)(unsafe.Add(mBase, _c_F_free[6])) = v144
	v166 = *(*int32)(unsafe.Add(mBase, _c_F_free[7]))
	v167 = v166 + v143
	*(*int32)(unsafe.Add(mBase, _c_F_free[7])) = v167
	*(*int32)(unsafe.Add(mBase, uint32(v144)+4)) = v167 | int32(1)
	v173 = *(*int32)(unsafe.Add(mBase, _c_F_free[1]))
	if v144 != v173 {
		goto L2
	} else {
		goto L49
	}
L49:
	;
	v175 = int32(0)
	*(*int32)(unsafe.Add(mBase, _c_F_free[3])) = v175
	*(*int32)(unsafe.Add(mBase, _c_F_free[1])) = v175
	goto L1
L50:
	;
	v198 = v152&int32(-8) + v143
	v199 = *(*int32)(unsafe.Add(mBase, uint32(v18)+12))
	if base.Ui32(int32(255)) < base.Ui32(v152) {
		goto L52
	} else {
		goto L53
	}
L51:
	;
	v184 = int32(0)
	*(*int32)(unsafe.Add(mBase, _c_F_free[1])) = v144
	v188 = *(*int32)(unsafe.Add(mBase, _c_F_free[3]))
	v189 = v188 + v143
	*(*int32)(unsafe.Add(mBase, _c_F_free[3])) = v189
	*(*int32)(unsafe.Add(mBase, uint32(v144)+4)) = v189 | int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v144+v189))) = v189
	goto L1
L52:
	;
	v215 = *(*int32)(unsafe.Add(mBase, uint32(v18)+24))
	if v199 == v18 {
		goto L56
	} else {
		goto L57
	}
L53:
	;
	v202 = *(*int32)(unsafe.Add(mBase, uint32(v18)+8))
	if v199 != v202 {
		goto L54
	} else {
		goto L55
	}
L54:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v199)+8)) = v202
	*(*int32)(unsafe.Add(mBase, uint32(v202)+12)) = v199
	goto L42
L55:
	;
	v204 = int32(0)
	v206 = *(*int32)(unsafe.Add(mBase, _c_F_free[2]))
	*(*int32)(unsafe.Add(mBase, _c_F_free[2])) = v206 & base.I32_rotl(int32(-2), int32(base.Ui32(v152)>>(uint(int32(3))%32)))
	goto L42
L56:
	;
	v220 = *(*int32)(unsafe.Add(mBase, uint32(v18)+20))
	if v220 == int32(0) {
		goto L59
	} else {
		goto L60
	}
L57:
	;
	v217 = *(*int32)(unsafe.Add(mBase, uint32(v18)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v217)+12)) = v199
	*(*int32)(unsafe.Add(mBase, uint32(v199)+8)) = v217
	v259 = v199
	goto L43
L58:
	;
	__phi236 = v230
	__phi237 = v231
	v236 = __phi236
	v237 = __phi237
	goto L62
L59:
	;
	v225 = *(*int32)(unsafe.Add(mBase, uint32(v18)+16))
	if v225 == int32(0) {
		goto L44
	} else {
		goto L61
	}
L60:
	;
	v230 = v220
	v231 = v18 + int32(20)
	goto L58
L61:
	;
	v230 = v225
	v231 = v18 + int32(16)
	goto L58
L62:
	;
	v242 = *(*int32)(unsafe.Add(mBase, uint32(v236)+20))
	if v242 != 0 {
		__phi236 = v242
		__phi237 = v236 + int32(20)
		v236 = __phi236
		v237 = __phi237
		goto L62
	} else {
		goto L64
	}
L63:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v237))) = int32(0)
	v259 = v236
	goto L43
L64:
	;
	v245 = *(*int32)(unsafe.Add(mBase, uint32(v236)+16))
	if v245 != 0 {
		__phi236 = v245
		__phi237 = v236 + int32(16)
		v236 = __phi236
		v237 = __phi237
		goto L62
	} else {
		goto L65
	}
L65:
	;
	goto L63
L66:
	;
	v267 = *(*int32)(unsafe.Add(mBase, uint32(v18)+28))
	v269 = v267 << (uint(int32(2)) % 32)
	v272 = *(*int32)(unsafe.Add(mBase, uint32(v269)+uint32(_c_F_free[4])))
	if v18 != v272 {
		goto L68
	} else {
		goto L69
	}
L67:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v259)+24)) = v215
	v292 = *(*int32)(unsafe.Add(mBase, uint32(v18)+16))
	if v292 == int32(0) {
		goto L75
	} else {
		goto L76
	}
L68:
	;
	v284 = *(*int32)(unsafe.Add(mBase, uint32(v215)+16))
	if v284 == v18 {
		goto L71
	} else {
		goto L72
	}
L69:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v269)+uint32(_c_F_free[4]))) = v259
	if v259 != 0 {
		goto L67
	} else {
		goto L70
	}
L70:
	;
	v275 = int32(0)
	v277 = *(*int32)(unsafe.Add(mBase, _c_F_free[5]))
	*(*int32)(unsafe.Add(mBase, _c_F_free[5])) = v277 & base.I32_rotl(int32(-2), v267)
	goto L42
L71:
	;
	v286 = int32(16)
	goto L73
L72:
	;
	v286 = int32(20)
	goto L73
L73:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v215+v286))) = v259
	if v259 == int32(0) {
		goto L42
	} else {
		goto L74
	}
L74:
	;
	goto L67
L75:
	;
	v297 = *(*int32)(unsafe.Add(mBase, uint32(v18)+20))
	if v297 == int32(0) {
		goto L42
	} else {
		goto L77
	}
L76:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v259)+16)) = v292
	*(*int32)(unsafe.Add(mBase, uint32(v292)+24)) = v259
	goto L75
L77:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v259)+20)) = v297
	*(*int32)(unsafe.Add(mBase, uint32(v297)+24)) = v259
	goto L42
L78:
	;
	*(*int32)(unsafe.Add(mBase, _c_F_free[3])) = v198
	goto L1
L79:
	;
	if base.Ui32(int32(16777215)) < base.Ui32(v320) {
		v366 = int32(31)
		goto L84
	} else {
		goto L85
	}
L80:
	;
	v331 = v320 & int32(-8)
	v333 = v331 + int32(_a_F_free_0)
	v335 = *(*int32)(unsafe.Add(mBase, _c_F_free[2]))
	v339 = int32(1) << (uint(int32(base.Ui32(v320)>>(uint(int32(3))%32))) % 32)
	if v335&v339 != 0 {
		goto L82
	} else {
		goto L83
	}
L81:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v345)+12)) = v144
	*(*int32)(unsafe.Add(mBase, uint32(v331)+uint32(_c_F_free[8]))) = v144
	*(*int32)(unsafe.Add(mBase, uint32(v144)+12)) = v333
	*(*int32)(unsafe.Add(mBase, uint32(v144)+8)) = v345
	goto L1
L82:
	;
	v344 = *(*int32)(unsafe.Add(mBase, uint32(v331)+uint32(_c_F_free[8])))
	v345 = v344
	goto L81
L83:
	;
	*(*int32)(unsafe.Add(mBase, _c_F_free[2])) = v335 | v339
	v345 = v333
	goto L81
L84:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v144)+28)) = v366
	*(*int64)(unsafe.Add(mBase, uint32(v144)+16)) = int64(0)
	v371 = v366 << (uint(int32(2)) % 32)
	v373 = v371 + int32(_a_F_free_1)
	v375 = *(*int32)(unsafe.Add(mBase, _c_F_free[5]))
	v377 = int32(1) << (uint(v366) % 32)
	if v375&v377 != 0 {
		goto L89
	} else {
		goto L90
	}
L85:
	;
	v356 = base.I32_clz(int32(base.Ui32(v320) >> (uint(int32(8)) % 32)))
	v359 = int32(1)
	v366 = int32(base.Ui32(v320)>>(uint(int32(38)-v356)%32))&v359 - v356<<(uint(v359)%32) + int32(62)
	goto L84
L86:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v436))) = v144
	*(*int32)(unsafe.Add(mBase, uint32(v144+v435))) = v438
	*(*int32)(unsafe.Add(mBase, uint32(v144)+12)) = v437
	*(*int32)(unsafe.Add(mBase, uint32(v144+v433))) = v440
	v447 = int32(0)
	v449 = *(*int32)(unsafe.Add(mBase, _c_F_free[9]))
	v450 = int32(-1)
	v451 = v449 + v450
	if v451 != 0 {
		goto L98
	} else {
		goto L99
	}
L87:
	;
	v426 = *(*int32)(unsafe.Add(mBase, uint32(v399)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v426)+12)) = v144
	v428 = int32(8)
	v433 = int32(24)
	v435 = v428
	v436 = v399 + v428
	v437 = v399
	v438 = v426
	v440 = int32(0)
	goto L86
L88:
	;
	v433 = v418
	v435 = v420
	v436 = v421
	v437 = v144
	v438 = v423
	v440 = v144
	goto L86
L89:
	;
	if v366 == int32(31) {
		goto L91
	} else {
		goto L92
	}
L90:
	;
	*(*int32)(unsafe.Add(mBase, _c_F_free[5])) = v375 | v377
	v418 = int32(8)
	v420 = int32(24)
	v421 = v373
	v423 = v373
	goto L88
L91:
	;
	v391 = int32(0)
	goto L93
L92:
	;
	v391 = int32(25) - int32(base.Ui32(v366)>>(uint(int32(1))%32))
	goto L93
L93:
	;
	v393 = *(*int32)(unsafe.Add(mBase, uint32(v371)+uint32(_c_F_free[4])))
	v396 = v320 << (uint(v391) % 32)
	v399 = v393
	goto L94
L94:
	;
	v402 = *(*int32)(unsafe.Add(mBase, uint32(v399)+4))
	if v402&int32(-8) == v320 {
		goto L87
	} else {
		goto L96
	}
L95:
	;
	v418 = int32(8)
	v420 = int32(24)
	v421 = v414
	v423 = v399
	goto L88
L96:
	;
	v414 = v399 + int32(base.Ui32(v396)>>(uint(int32(29))%32))&int32(4) + int32(16)
	v415 = *(*int32)(unsafe.Add(mBase, uint32(v414)))
	if v415 != 0 {
		v396 = v396 << (uint(int32(1)) % 32)
		v399 = v415
		goto L94
	} else {
		goto L97
	}
L97:
	;
	goto L95
L98:
	;
	v453 = v451
	goto L100
L99:
	;
	v453 = v450
	goto L100
L100:
	;
	*(*int32)(unsafe.Add(mBase, _c_F_free[9])) = v453
	goto L2
}
