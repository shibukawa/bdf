//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	"unsafe"
)

func F_SSE16x16_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v29 int32
	_ = v29
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v67 int32
	_ = v67
	var v69 int32
	_ = v69
	var v72 int32
	_ = v72
	var v73 int32
	_ = v73
	var v75 int32
	_ = v75
	var v77 int32
	_ = v77
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v83 int32
	_ = v83
	var v85 int32
	_ = v85
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v93 int32
	_ = v93
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v99 int32
	_ = v99
	var v101 int32
	_ = v101
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v107 int32
	_ = v107
	var v109 int32
	_ = v109
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v115 int32
	_ = v115
	var v117 int32
	_ = v117
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v123 int32
	_ = v123
	var v125 int32
	_ = v125
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v131 int32
	_ = v131
	var v133 int32
	_ = v133
	var v136 int32
	_ = v136
	var v137 int32
	_ = v137
	var v139 int32
	_ = v139
	var v140 int32
	_ = v140
	var v141 int32
	_ = v141
	var v158 int32
	_ = v158
	var v160 int32
	_ = v160
	v3 = int32(0)
	v12 = v3
	v13 = v3
	for {
		v17 = l0 + v12
		v18 = int32(15)
		v20 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v18))))
		v21 = l1 + v12
		v24 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v18))))
		v25 = v20 - v24
		v27 = int32(14)
		v29 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v27))))
		v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v27))))
		v33 = v29 - v32
		v35 = int32(13)
		v37 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v35))))
		v40 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v35))))
		v41 = v37 - v40
		v43 = int32(12)
		v45 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v43))))
		v48 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v43))))
		v49 = v45 - v48
		v51 = int32(11)
		v53 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v51))))
		v56 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v51))))
		v57 = v53 - v56
		v59 = int32(10)
		v61 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v59))))
		v64 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v59))))
		v65 = v61 - v64
		v67 = int32(9)
		v69 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v67))))
		v72 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v67))))
		v73 = v69 - v72
		v75 = int32(8)
		v77 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v75))))
		v80 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v75))))
		v81 = v77 - v80
		v83 = int32(7)
		v85 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v83))))
		v88 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v83))))
		v89 = v85 - v88
		v91 = int32(6)
		v93 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v91))))
		v96 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v91))))
		v97 = v93 - v96
		v99 = int32(5)
		v101 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v99))))
		v104 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v99))))
		v105 = v101 - v104
		v107 = int32(4)
		v109 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v107))))
		v112 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v107))))
		v113 = v109 - v112
		v115 = int32(3)
		v117 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v115))))
		v120 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v115))))
		v121 = v117 - v120
		v123 = int32(2)
		v125 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v123))))
		v128 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v123))))
		v129 = v125 - v128
		v131 = int32(1)
		v133 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v131))))
		v136 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v131))))
		v137 = v133 - v136
		v139 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17))))
		v140 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21))))
		v141 = v139 - v140
		v158 = v25*v25 + (v33*v33 + (v41*v41 + (v49*v49 + (v57*v57 + (v65*v65 + (v73*v73 + (v81*v81 + (v89*v89 + (v97*v97 + (v105*v105 + (v113*v113 + (v121*v121 + (v129*v129 + (v137*v137 + (v141*v141 + v13)))))))))))))))
		v160 = v12 + int32(32)
		if v160 != int32(512) {
			v12 = v160
			v13 = v158
			continue
		} else {
			break
		}
		break
	}
	return v158
}
func F_SSE16x8_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v29 int32
	_ = v29
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v67 int32
	_ = v67
	var v69 int32
	_ = v69
	var v72 int32
	_ = v72
	var v73 int32
	_ = v73
	var v75 int32
	_ = v75
	var v77 int32
	_ = v77
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v83 int32
	_ = v83
	var v85 int32
	_ = v85
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v93 int32
	_ = v93
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v99 int32
	_ = v99
	var v101 int32
	_ = v101
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v107 int32
	_ = v107
	var v109 int32
	_ = v109
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v115 int32
	_ = v115
	var v117 int32
	_ = v117
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v123 int32
	_ = v123
	var v125 int32
	_ = v125
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v131 int32
	_ = v131
	var v133 int32
	_ = v133
	var v136 int32
	_ = v136
	var v137 int32
	_ = v137
	var v139 int32
	_ = v139
	var v140 int32
	_ = v140
	var v141 int32
	_ = v141
	var v158 int32
	_ = v158
	var v160 int32
	_ = v160
	v3 = int32(0)
	v12 = v3
	v13 = v3
	for {
		v17 = l0 + v12
		v18 = int32(15)
		v20 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v18))))
		v21 = l1 + v12
		v24 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v18))))
		v25 = v20 - v24
		v27 = int32(14)
		v29 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v27))))
		v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v27))))
		v33 = v29 - v32
		v35 = int32(13)
		v37 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v35))))
		v40 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v35))))
		v41 = v37 - v40
		v43 = int32(12)
		v45 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v43))))
		v48 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v43))))
		v49 = v45 - v48
		v51 = int32(11)
		v53 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v51))))
		v56 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v51))))
		v57 = v53 - v56
		v59 = int32(10)
		v61 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v59))))
		v64 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v59))))
		v65 = v61 - v64
		v67 = int32(9)
		v69 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v67))))
		v72 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v67))))
		v73 = v69 - v72
		v75 = int32(8)
		v77 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v75))))
		v80 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v75))))
		v81 = v77 - v80
		v83 = int32(7)
		v85 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v83))))
		v88 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v83))))
		v89 = v85 - v88
		v91 = int32(6)
		v93 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v91))))
		v96 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v91))))
		v97 = v93 - v96
		v99 = int32(5)
		v101 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v99))))
		v104 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v99))))
		v105 = v101 - v104
		v107 = int32(4)
		v109 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v107))))
		v112 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v107))))
		v113 = v109 - v112
		v115 = int32(3)
		v117 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v115))))
		v120 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v115))))
		v121 = v117 - v120
		v123 = int32(2)
		v125 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v123))))
		v128 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v123))))
		v129 = v125 - v128
		v131 = int32(1)
		v133 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17+v131))))
		v136 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v131))))
		v137 = v133 - v136
		v139 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17))))
		v140 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21))))
		v141 = v139 - v140
		v158 = v25*v25 + (v33*v33 + (v41*v41 + (v49*v49 + (v57*v57 + (v65*v65 + (v73*v73 + (v81*v81 + (v89*v89 + (v97*v97 + (v105*v105 + (v113*v113 + (v121*v121 + (v129*v129 + (v137*v137 + (v141*v141 + v13)))))))))))))))
		v160 = v12 + int32(32)
		if v160 != int32(256) {
			v12 = v160
			v13 = v158
			continue
		} else {
			break
		}
		break
	}
	return v158
}
func F_SSE4x4_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v5 int32
	_ = v5
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v42 int32
	_ = v42
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
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	v4 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+99)))
	v5 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+99)))
	v6 = v4 - v5
	v8 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+98)))
	v9 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+98)))
	v10 = v8 - v9
	v12 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+97)))
	v13 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+97)))
	v14 = v12 - v13
	v16 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+96)))
	v17 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+96)))
	v18 = v16 - v17
	v20 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+67)))
	v21 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+67)))
	v22 = v20 - v21
	v24 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+66)))
	v25 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+66)))
	v26 = v24 - v25
	v28 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+65)))
	v29 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+65)))
	v30 = v28 - v29
	v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+64)))
	v33 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+64)))
	v34 = v32 - v33
	v36 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+35)))
	v37 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+35)))
	v38 = v36 - v37
	v40 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+34)))
	v41 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+34)))
	v42 = v40 - v41
	v44 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+33)))
	v45 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+33)))
	v46 = v44 - v45
	v48 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+32)))
	v49 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+32)))
	v50 = v48 - v49
	v52 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)))
	v53 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3)))
	v54 = v52 - v53
	v56 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)))
	v57 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+2)))
	v58 = v56 - v57
	v60 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)))
	v61 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+1)))
	v62 = v60 - v61
	v64 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v65 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
	v66 = v64 - v65
	return v6*v6 + (v10*v10 + (v14*v14 + (v18*v18 + (v22*v22 + (v26*v26 + (v30*v30 + (v34*v34 + (v38*v38 + (v42*v42 + (v46*v46 + (v50*v50 + (v54*v54 + (v58*v58 + (v62*v62 + v66*v66))))))))))))))
}
func F_SSE8x8_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v5 int32
	_ = v5
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v42 int32
	_ = v42
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
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v68 int32
	_ = v68
	var v69 int32
	_ = v69
	var v70 int32
	_ = v70
	var v72 int32
	_ = v72
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v77 int32
	_ = v77
	var v78 int32
	_ = v78
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v84 int32
	_ = v84
	var v85 int32
	_ = v85
	var v86 int32
	_ = v86
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v94 int32
	_ = v94
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v100 int32
	_ = v100
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v116 int32
	_ = v116
	var v117 int32
	_ = v117
	var v118 int32
	_ = v118
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v122 int32
	_ = v122
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v126 int32
	_ = v126
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v132 int32
	_ = v132
	var v133 int32
	_ = v133
	var v134 int32
	_ = v134
	var v136 int32
	_ = v136
	var v137 int32
	_ = v137
	var v138 int32
	_ = v138
	var v140 int32
	_ = v140
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v144 int32
	_ = v144
	var v145 int32
	_ = v145
	var v146 int32
	_ = v146
	var v148 int32
	_ = v148
	var v149 int32
	_ = v149
	var v150 int32
	_ = v150
	var v152 int32
	_ = v152
	var v153 int32
	_ = v153
	var v154 int32
	_ = v154
	var v156 int32
	_ = v156
	var v157 int32
	_ = v157
	var v158 int32
	_ = v158
	var v160 int32
	_ = v160
	var v161 int32
	_ = v161
	var v162 int32
	_ = v162
	var v164 int32
	_ = v164
	var v165 int32
	_ = v165
	var v166 int32
	_ = v166
	var v168 int32
	_ = v168
	var v169 int32
	_ = v169
	var v170 int32
	_ = v170
	var v172 int32
	_ = v172
	var v173 int32
	_ = v173
	var v174 int32
	_ = v174
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v178 int32
	_ = v178
	var v180 int32
	_ = v180
	var v181 int32
	_ = v181
	var v182 int32
	_ = v182
	var v184 int32
	_ = v184
	var v185 int32
	_ = v185
	var v186 int32
	_ = v186
	var v188 int32
	_ = v188
	var v189 int32
	_ = v189
	var v190 int32
	_ = v190
	var v192 int32
	_ = v192
	var v193 int32
	_ = v193
	var v194 int32
	_ = v194
	var v196 int32
	_ = v196
	var v197 int32
	_ = v197
	var v198 int32
	_ = v198
	var v200 int32
	_ = v200
	var v201 int32
	_ = v201
	var v202 int32
	_ = v202
	var v204 int32
	_ = v204
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
	var v212 int32
	_ = v212
	var v213 int32
	_ = v213
	var v214 int32
	_ = v214
	var v216 int32
	_ = v216
	var v217 int32
	_ = v217
	var v218 int32
	_ = v218
	var v220 int32
	_ = v220
	var v221 int32
	_ = v221
	var v222 int32
	_ = v222
	var v224 int32
	_ = v224
	var v225 int32
	_ = v225
	var v226 int32
	_ = v226
	var v228 int32
	_ = v228
	var v229 int32
	_ = v229
	var v230 int32
	_ = v230
	var v232 int32
	_ = v232
	var v233 int32
	_ = v233
	var v234 int32
	_ = v234
	var v236 int32
	_ = v236
	var v237 int32
	_ = v237
	var v238 int32
	_ = v238
	var v240 int32
	_ = v240
	var v241 int32
	_ = v241
	var v242 int32
	_ = v242
	var v244 int32
	_ = v244
	var v245 int32
	_ = v245
	var v246 int32
	_ = v246
	var v248 int32
	_ = v248
	var v249 int32
	_ = v249
	var v250 int32
	_ = v250
	var v252 int32
	_ = v252
	var v253 int32
	_ = v253
	var v254 int32
	_ = v254
	var v256 int32
	_ = v256
	var v257 int32
	_ = v257
	var v258 int32
	_ = v258
	v4 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+231)))
	v5 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+231)))
	v6 = v4 - v5
	v8 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+230)))
	v9 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+230)))
	v10 = v8 - v9
	v12 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+229)))
	v13 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+229)))
	v14 = v12 - v13
	v16 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+228)))
	v17 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+228)))
	v18 = v16 - v17
	v20 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+227)))
	v21 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+227)))
	v22 = v20 - v21
	v24 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+226)))
	v25 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+226)))
	v26 = v24 - v25
	v28 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+225)))
	v29 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+225)))
	v30 = v28 - v29
	v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+224)))
	v33 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+224)))
	v34 = v32 - v33
	v36 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+199)))
	v37 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+199)))
	v38 = v36 - v37
	v40 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+198)))
	v41 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+198)))
	v42 = v40 - v41
	v44 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+197)))
	v45 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+197)))
	v46 = v44 - v45
	v48 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+196)))
	v49 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+196)))
	v50 = v48 - v49
	v52 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+195)))
	v53 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+195)))
	v54 = v52 - v53
	v56 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+194)))
	v57 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+194)))
	v58 = v56 - v57
	v60 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+193)))
	v61 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+193)))
	v62 = v60 - v61
	v64 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+192)))
	v65 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+192)))
	v66 = v64 - v65
	v68 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+167)))
	v69 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+167)))
	v70 = v68 - v69
	v72 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+166)))
	v73 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+166)))
	v74 = v72 - v73
	v76 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+165)))
	v77 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+165)))
	v78 = v76 - v77
	v80 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+164)))
	v81 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+164)))
	v82 = v80 - v81
	v84 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+163)))
	v85 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+163)))
	v86 = v84 - v85
	v88 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+162)))
	v89 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+162)))
	v90 = v88 - v89
	v92 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+161)))
	v93 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+161)))
	v94 = v92 - v93
	v96 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+160)))
	v97 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+160)))
	v98 = v96 - v97
	v100 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+135)))
	v101 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+135)))
	v102 = v100 - v101
	v104 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+134)))
	v105 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+134)))
	v106 = v104 - v105
	v108 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+133)))
	v109 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+133)))
	v110 = v108 - v109
	v112 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+132)))
	v113 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+132)))
	v114 = v112 - v113
	v116 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+131)))
	v117 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+131)))
	v118 = v116 - v117
	v120 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+130)))
	v121 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+130)))
	v122 = v120 - v121
	v124 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+129)))
	v125 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+129)))
	v126 = v124 - v125
	v128 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+128)))
	v129 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+128)))
	v130 = v128 - v129
	v132 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+103)))
	v133 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+103)))
	v134 = v132 - v133
	v136 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+102)))
	v137 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+102)))
	v138 = v136 - v137
	v140 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+101)))
	v141 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+101)))
	v142 = v140 - v141
	v144 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+100)))
	v145 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+100)))
	v146 = v144 - v145
	v148 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+99)))
	v149 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+99)))
	v150 = v148 - v149
	v152 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+98)))
	v153 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+98)))
	v154 = v152 - v153
	v156 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+97)))
	v157 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+97)))
	v158 = v156 - v157
	v160 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+96)))
	v161 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+96)))
	v162 = v160 - v161
	v164 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+71)))
	v165 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+71)))
	v166 = v164 - v165
	v168 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+70)))
	v169 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+70)))
	v170 = v168 - v169
	v172 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+69)))
	v173 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+69)))
	v174 = v172 - v173
	v176 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+68)))
	v177 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+68)))
	v178 = v176 - v177
	v180 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+67)))
	v181 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+67)))
	v182 = v180 - v181
	v184 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+66)))
	v185 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+66)))
	v186 = v184 - v185
	v188 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+65)))
	v189 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+65)))
	v190 = v188 - v189
	v192 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+64)))
	v193 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+64)))
	v194 = v192 - v193
	v196 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+39)))
	v197 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+39)))
	v198 = v196 - v197
	v200 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+38)))
	v201 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+38)))
	v202 = v200 - v201
	v204 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+37)))
	v205 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+37)))
	v206 = v204 - v205
	v208 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+36)))
	v209 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+36)))
	v210 = v208 - v209
	v212 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+35)))
	v213 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+35)))
	v214 = v212 - v213
	v216 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+34)))
	v217 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+34)))
	v218 = v216 - v217
	v220 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+33)))
	v221 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+33)))
	v222 = v220 - v221
	v224 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+32)))
	v225 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+32)))
	v226 = v224 - v225
	v228 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+7)))
	v229 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+7)))
	v230 = v228 - v229
	v232 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+6)))
	v233 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+6)))
	v234 = v232 - v233
	v236 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+5)))
	v237 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+5)))
	v238 = v236 - v237
	v240 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+4)))
	v241 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+4)))
	v242 = v240 - v241
	v244 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)))
	v245 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3)))
	v246 = v244 - v245
	v248 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)))
	v249 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+2)))
	v250 = v248 - v249
	v252 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)))
	v253 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+1)))
	v254 = v252 - v253
	v256 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v257 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
	v258 = v256 - v257
	return v6*v6 + (v10*v10 + (v14*v14 + (v18*v18 + (v22*v22 + (v26*v26 + (v30*v30 + (v34*v34 + (v38*v38 + (v42*v42 + (v46*v46 + (v50*v50 + (v54*v54 + (v58*v58 + (v62*v62 + (v66*v66 + (v70*v70 + (v74*v74 + (v78*v78 + (v82*v82 + (v86*v86 + (v90*v90 + (v94*v94 + (v98*v98 + (v102*v102 + (v106*v106 + (v110*v110 + (v114*v114 + (v118*v118 + (v122*v122 + (v126*v126 + (v130*v130 + (v134*v134 + (v138*v138 + (v142*v142 + (v146*v146 + (v150*v150 + (v154*v154 + (v158*v158 + (v162*v162 + (v166*v166 + (v170*v170 + (v174*v174 + (v178*v178 + (v182*v182 + (v186*v186 + (v190*v190 + (v194*v194 + (v198*v198 + (v202*v202 + (v206*v206 + (v210*v210 + (v214*v214 + (v218*v218 + (v222*v222 + (v226*v226 + (v230*v230 + (v234*v234 + (v238*v238 + (v242*v242 + (v246*v246 + (v250*v250 + (v254*v254 + v258*v258))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))
}
func F_SSIMGetClipped_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32) float64 {
	mBase := m.M
	_ = mBase
	var v9 int64
	_ = v9
	var v10 int32
	_ = v10
	var v36 int32
	_ = v36
	var v39 int32
	_ = v39
	var v41 int32
	_ = v41
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v47 int32
	_ = v47
	var v49 int64
	_ = v49
	var v53 int64
	_ = v53
	var v57 int32
	_ = v57
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v66 int32
	_ = v66
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v73 int32
	_ = v73
	var v86 int32
	_ = v86
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
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
	var v118 int32
	_ = v118
	var v128 int32
	_ = v128
	var v134 int32
	_ = v134
	var v135 int32
	_ = v135
	var v139 int32
	_ = v139
	var v142 int32
	_ = v142
	var v144 int32
	_ = v144
	var v157 int32
	_ = v157
	var v158 int32
	_ = v158
	var v159 int32
	_ = v159
	var v160 int32
	_ = v160
	var v161 int32
	_ = v161
	var v168 int32
	_ = v168
	var v169 int32
	_ = v169
	var v171 int32
	_ = v171
	var v172 int32
	_ = v172
	var v173 int32
	_ = v173
	var v175 int32
	_ = v175
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v179 int32
	_ = v179
	var v181 int32
	_ = v181
	var v183 int32
	_ = v183
	var v188 int32
	_ = v188
	var v190 int32
	_ = v190
	var v209 int64
	_ = v209
	var v210 int32
	_ = v210
	var v214 int64
	_ = v214
	var v215 int64
	_ = v215
	var v216 int64
	_ = v216
	var v217 int64
	_ = v217
	var v237 int64
	_ = v237
	var v238 int32
	_ = v238
	var v243 int64
	_ = v243
	var v245 int64
	_ = v245
	var v246 int64
	_ = v246
	var v247 int64
	_ = v247
	var v250 int64
	_ = v250
	var v251 int64
	_ = v251
	var v255 int64
	_ = v255
	var v257 int64
	_ = v257
	var v263 int64
	_ = v263
	var v279 float64
	_ = v279
	v9 = int64(0)
	v10 = int32(0)
	v36 = int32(3)
	if v36 < l5 {
		v39 = l5
	} else {
		v39 = v36
	}
	v41 = v39 + int32(-3)
	v43 = l5 + int32(3)
	v45 = l7 + int32(-1)
	if v43 < v45 {
		v47 = v43
	} else {
		v47 = v45
	}
	if v41 <= v47 {
		v53 = int64(0)
		v57 = int32(3)
		if v57 < l4 {
			v60 = l4
		} else {
			v60 = v57
		}
		v62 = v60 + int32(-3)
		v64 = l4 + int32(3)
		v66 = l6 + int32(-1)
		if v64 < v66 {
			v68 = v64
		} else {
			v68 = v66
		}
		if v68 < v62 {
			v209 = v9
			v210 = v10
			v214 = v53
			v215 = v53
			v216 = v53
			v217 = v53
		} else {
			v70 = m.G1
			v73 = int32(2)
			v86 = int32(0)
			v101 = v86
			v102 = v41
			v112 = l2 + v41*l3
			v113 = l0 + v41*l1
			v114 = v86
			v115 = v86
			v116 = v86
			v117 = v86
			v118 = v86
			for {
				v128 = m.G1
				v134 = *(*int32)(unsafe.Add(mBase, uint32((v102-l5)<<(uint(int32(2))%32)+(v128+int32(_a_F_SSIMGetClipped_C_0))+int32(12))))
				v135 = v68 - v60 + int32(4)
				v139 = v62
				v142 = v70 + int32(_a_F_SSIMGetClipped_C_0) + (v60<<(uint(v73)%32) - l4<<(uint(v73)%32))
				v144 = v101
				v157 = v114
				v158 = v115
				v159 = v116
				v160 = v117
				v161 = v118
				for {
					v168 = *(*int32)(unsafe.Add(mBase, uint32(v142)))
					v169 = v134 * v168
					v171 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v112+v139))))
					v172 = v169 * v171
					v173 = v172 + v160
					v175 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v113+v139))))
					v176 = v169 * v175
					v177 = v176 + v161
					v179 = v172*v171 + v157
					v181 = v176*v171 + v158
					v183 = v176*v175 + v159
					v188 = v169 + v144
					v190 = v135 + int32(-1)
					if v190 != 0 {
						v135 = v190
						v139 = v139 + int32(1)
						v142 = v142 + int32(4)
						v144 = v188
						v157 = v179
						v158 = v181
						v159 = v183
						v160 = v173
						v161 = v177
						continue
					} else {
						break
					}
					break
				}
				if v102 < v47 {
					v101 = v188
					v102 = v102 + int32(1)
					v112 = v112 + l3
					v113 = v113 + l1
					v114 = v179
					v115 = v181
					v116 = v183
					v117 = v173
					v118 = v177
					continue
				} else {
					break
				}
				break
			}
			v209 = base.I64_extend_i32_u(v177)
			v210 = v188
			v214 = base.I64_extend_i32_u(v173)
			v215 = base.I64_extend_i32_u(v183)
			v216 = base.I64_extend_i32_u(v181)
			v217 = base.I64_extend_i32_u(v179)
		}
	} else {
		v49 = int64(0)
		v209 = v9
		v210 = v10
		v214 = v49
		v215 = v49
		v216 = v49
		v217 = v49
	}
	v237 = v214*v214 + v209*v209
	v238 = v210 * v210
	if base.Ui64(v237) < base.Ui64(base.I64_extend_i32_u(v238<<(uint(int32(6))%32))) {
		v279 = float64(1)
	} else {
		v243 = base.I64_extend_i32_u(v210)
		v245 = v214 * v209
		v246 = v216*v243 - v245
		v247 = int64(0)
		if v247 < v246 {
			v250 = v246
		} else {
			v250 = v247
		}
		v251 = int64(1)
		v255 = base.I64_extend_i32_u(v238 * int32(60))
		v257 = int64(8)
		v263 = base.I64_extend_i32_u(v238 * int32(20))
		v279 = base.F64_div(base.F64_convert_i64_u(int64(base.Ui64(v250<<(uint(v251)%64)+v255)>>(uint(v257)%64))*(v245<<(uint(v251)%64)+v263)), base.F64_convert_i64_u(int64(base.Ui64(v255-v237+(v217+v215)*v243)>>(uint(v257)%64))*(v237+v263)))
	}
	return v279
}
func F_SSIMGet_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) float64 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
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
	var v91 int32
	_ = v91
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v103 int32
	_ = v103
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v107 int32
	_ = v107
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v111 int32
	_ = v111
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
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v126 int32
	_ = v126
	var v127 int32
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
	var v134 int32
	_ = v134
	var v135 int32
	_ = v135
	var v136 int32
	_ = v136
	var v137 int32
	_ = v137
	var v138 int32
	_ = v138
	var v145 int32
	_ = v145
	var v159 int32
	_ = v159
	var v173 int32
	_ = v173
	var v187 int32
	_ = v187
	var v191 int32
	_ = v191
	var v193 int64
	_ = v193
	var v195 int64
	_ = v195
	var v197 int64
	_ = v197
	var v203 int64
	_ = v203
	var v204 int64
	_ = v204
	var v205 int64
	_ = v205
	var v208 int64
	_ = v208
	var v209 int64
	_ = v209
	var v211 int64
	_ = v211
	var v213 int64
	_ = v213
	var v217 int64
	_ = v217
	var v236 float64
	_ = v236
	v5 = int32(0)
	v49 = l0
	v51 = l2
	v53 = v5
	v54 = int32(-28)
	v55 = v5
	v56 = v5
	v57 = v5
	v58 = v5
	for {
		v91 = m.G1
		v97 = *(*int32)(unsafe.Add(mBase, uint32(v91+int32(_a_F_SSIMGet_C_0)+v54+int32(28))))
		v98 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v51)+6)))
		v99 = v97 * v98
		v101 = v97 << (uint(int32(1)) % 32)
		v102 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v51)+5)))
		v103 = v101 * v102
		v105 = v97 * int32(3)
		v106 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v51)+4)))
		v107 = v105 * v106
		v109 = v97 << (uint(int32(2)) % 32)
		v110 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v51)+3)))
		v111 = v109 * v110
		v112 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v51)+2)))
		v113 = v105 * v112
		v114 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v51)+1)))
		v115 = v101 * v114
		v116 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v51))))
		v117 = v97 * v116
		v124 = v99 + (v103 + (v107 + (v111 + (v113 + (v115 + (v117 + v57))))))
		v125 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v49)+6)))
		v126 = v97 * v125
		v127 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v49)+5)))
		v128 = v101 * v127
		v129 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v49)+4)))
		v130 = v105 * v129
		v131 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v49)+3)))
		v132 = v109 * v131
		v133 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v49)+2)))
		v134 = v105 * v133
		v135 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v49)+1)))
		v136 = v101 * v135
		v137 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v49))))
		v138 = v97 * v137
		v145 = v126 + (v128 + (v130 + (v132 + (v134 + (v136 + (v138 + v58))))))
		v159 = v99*v98 + (v103*v102 + (v107*v106 + (v111*v110 + (v113*v112 + (v115*v114 + (v117*v116 + v53))))))
		v173 = v126*v98 + (v128*v102 + (v130*v106 + (v132*v110 + (v134*v112 + (v136*v114 + (v138*v116 + v55))))))
		v187 = v126*v125 + (v128*v127 + (v130*v129 + (v132*v131 + (v134*v133 + (v136*v135 + (v138*v137 + v56))))))
		v191 = v54 + int32(4)
		if v191 != 0 {
			v49 = v49 + l1
			v51 = v51 + l3
			v53 = v159
			v54 = v191
			v55 = v173
			v56 = v187
			v57 = v124
			v58 = v145
			continue
		} else {
			break
		}
		break
	}
	v193 = base.I64_extend_i32_u(v124)
	v195 = base.I64_extend_i32_u(v145)
	v197 = v193*v193 + v195*v195
	if base.Ui64(v197) < base.Ui64(int64(4194304)) {
		v236 = float64(1)
	} else {
		v203 = v193 * v195
		v204 = base.I64_extend_i32_u(v173)<<(uint(int64(8))%64) - v203
		v205 = int64(0)
		if v205 < v204 {
			v208 = v204
		} else {
			v208 = v205
		}
		v209 = int64(1)
		v211 = int64(3932160)
		v213 = int64(8)
		v217 = int64(1310720)
		v236 = base.F64_div(base.F64_convert_i64_u(int64(base.Ui64(v208<<(uint(v209)%64)+v211)>>(uint(v213)%64))*(v203<<(uint(v209)%64)+v217)), base.F64_convert_i64_u(int64(base.Ui64((base.I64_extend_i32_u(v159)+base.I64_extend_i32_u(v187))<<(uint(v213)%64)-v197+v211)>>(uint(v213)%64))*(v197+v217)))
	}
	return v236
}
func F_SearchColorNoIdx(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v14 int32
	_ = v14
	var v16 int32
	_ = v16
	var v21 int32
	_ = v21
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v33 int32
	_ = v33
	v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if v9 == l1 {
		v33 = int32(0)
	} else {
		v14 = l2
		v16 = int32(0)
		for {
			v21 = (v14 + v16) >> (uint(int32(1)) % 32)
			v25 = *(*int32)(unsafe.Add(mBase, uint32(l0+v21<<(uint(int32(2))%32))))
			v26 = base.B2i32(base.Ui32(v25) < base.Ui32(l1))
			if base.Ui32(v25) < base.Ui32(l1) {
				v27 = v14
			} else {
				v27 = v21
			}
			if base.Ui32(v25) < base.Ui32(l1) {
				v28 = v21
			} else {
				v28 = v16
			}
			if v25 != l1 {
				v14 = v27
				v16 = v28
				continue
			} else {
				break
			}
			break
		}
		v33 = v21
	}
	return v33
}
func F_SetBitDepths(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v9 int32
	_ = v9
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v16 int32
	_ = v16
	var v18 int32
	_ = v18
	var v20 int32
	_ = v20
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v32 int32
	_ = v32
	var v34 int32
	_ = v34
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v41 int32
	_ = v41
	var v44 int32
	_ = v44
	var v46 int32
	_ = v46
	var v49 int32
	_ = v49
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v56 int32
	_ = v56
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	if v6 < int32(0) {
		v56 = l0
		v59 = l3
	} else {
		v9 = l0
		v12 = l3
		v13 = v6
		for {
			v16 = l1 + v13<<(uint(int32(4))%32)
			v18 = v12 + int32(1)
			v20 = *(*int32)(unsafe.Add(mBase, uint32(v16)+8))
			if v20 < int32(0) {
				v41 = v16
				v44 = v18
			} else {
				v23 = v16
				v26 = v18
				v27 = v20
				for {
					v28 = int32(4)
					v32 = v26 + int32(1)
					F_SetBitDepths(m, l1+v27<<(uint(v28)%32), l1, l2, v32)
					mBase = m.M
					v34 = *(*int32)(unsafe.Add(mBase, uint32(v23)+12))
					v37 = l1 + v34<<(uint(v28)%32)
					v38 = *(*int32)(unsafe.Add(mBase, uint32(v37)+8))
					if int32(-1) < v38 {
						v23 = v37
						v26 = v32
						v27 = v38
						continue
					} else {
						break
					}
					break
				}
				v41 = v37
				v44 = v32
			}
			v46 = *(*int32)(unsafe.Add(mBase, uint32(v41)+4))
			*(*uint8)(unsafe.Add(mBase, uint32(l2+v46))) = uint8(v44)
			v49 = *(*int32)(unsafe.Add(mBase, uint32(v9)+12))
			v52 = l1 + v49<<(uint(int32(4))%32)
			v53 = *(*int32)(unsafe.Add(mBase, uint32(v52)+8))
			if int32(-1) < v53 {
				v9 = v52
				v12 = v18
				v13 = v53
				continue
			} else {
				break
			}
			break
		}
		v56 = v52
		v59 = v18
	}
	v61 = *(*int32)(unsafe.Add(mBase, uint32(v56)+4))
	*(*uint8)(unsafe.Add(mBase, uint32(l2+v61))) = uint8(v59)
	return
}
func F_SetLoopParams(m *base.Module, l0 int32, l1 float32) {
	mBase := m.M
	_ = mBase
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v18 float32
	_ = v18
	var v21 float32
	_ = v21
	var v24 float32
	_ = v24
	var v26 int64
	_ = v26
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v47 int32
	_ = v47
	var v50 int32
	_ = v50
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v68 int32
	_ = v68
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v81 int32
	_ = v81
	var v89 int32
	_ = v89
	var v100 int32
	_ = v100
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v126 int64
	_ = v126
	var v130 int64
	_ = v130
	var v132 int32
	_ = v132
	var v135 int32
	_ = v135
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
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v147 int32
	_ = v147
	var v151 int32
	_ = v151
	var v152 int32
	_ = v152
	var v157 int32
	_ = v157
	var v161 int32
	_ = v161
	var v162 int32
	_ = v162
	var v168 int32
	_ = v168
	var v172 int32
	_ = v172
	var v173 int32
	_ = v173
	var v175 int32
	_ = v175
	var v179 int32
	_ = v179
	var v183 int32
	_ = v183
	var v186 int32
	_ = v186
	var v188 int32
	_ = v188
	var v189 int32
	_ = v189
	var v193 int32
	_ = v193
	var v198 int32
	_ = v198
	var v201 int32
	_ = v201
	var v207 int32
	_ = v207
	var v208 int32
	_ = v208
	var v209 int32
	_ = v209
	var v211 int32
	_ = v211
	var v216 int32
	_ = v216
	var v217 int32
	_ = v217
	var v218 int32
	_ = v218
	var v221 int32
	_ = v221
	var v222 int32
	_ = v222
	var v223 int32
	_ = v223
	var v224 int32
	_ = v224
	var v229 int32
	_ = v229
	var v234 int32
	_ = v234
	var v238 int32
	_ = v238
	var v239 int32
	_ = v239
	var v243 int32
	_ = v243
	var v244 int32
	_ = v244
	var v246 int32
	_ = v246
	var v248 int32
	_ = v248
	var v251 int32
	_ = v251
	var v257 int32
	_ = v257
	var v265 int32
	_ = v265
	var v274 int32
	_ = v274
	var v282 int32
	_ = v282
	var v293 int32
	_ = v293
	var v300 int32
	_ = v300
	var v352 int32
	_ = v352
	var v367 int32
	_ = v367
	var v368 int32
	_ = v368
	var v369 int32
	_ = v369
	var v371 int32
	_ = v371
	var v374 int32
	_ = v374
	var v377 int32
	_ = v377
	var v380 int32
	_ = v380
	var v383 int32
	_ = v383
	var v386 int32
	_ = v386
	var v389 int32
	_ = v389
	var v392 int32
	_ = v392
	var v395 int32
	_ = v395
	var v398 int32
	_ = v398
	var v401 int32
	_ = v401
	var v404 int32
	_ = v404
	var v407 int32
	_ = v407
	var v410 int32
	_ = v410
	var v413 int32
	_ = v413
	var v418 int32
	_ = v418
	var v419 int32
	_ = v419
	var v420 int32
	_ = v420
	var v441 int32
	_ = v441
	var v460 int32
	_ = v460
	var v492 int32
	_ = v492
	var v493 int32
	_ = v493
	var v494 int32
	_ = v494
	var v495 int32
	_ = v495
	var v508 int32
	_ = v508
	var v509 int32
	_ = v509
	var v510 int32
	_ = v510
	var v513 int32
	_ = v513
	var v514 int32
	_ = v514
	var v518 int32
	_ = v518
	var v525 int32
	_ = v525
	var v558 int32
	_ = v558
	var v569 int32
	_ = v569
	var v574 int32
	_ = v574
	var v575 int32
	_ = v575
	var v577 int32
	_ = v577
	var v586 int32
	_ = v586
	var v613 int32
	_ = v613
	var v615 int32
	_ = v615
	var v616 int32
	_ = v616
	var v624 int32
	_ = v624
	var v625 int32
	_ = v625
	var v627 int32
	_ = v627
	var v636 int32
	_ = v636
	var v638 int32
	_ = v638
	var v640 int32
	_ = v640
	var v684 int32
	_ = v684
	var v689 int32
	_ = v689
	var v692 int32
	_ = v692
	var v695 int32
	_ = v695
	var v698 int32
	_ = v698
	var v699 int32
	_ = v699
	var v700 int32
	_ = v700
	var v704 int32
	_ = v704
	var v705 int32
	_ = v705
	var v706 int32
	_ = v706
	var v711 int32
	_ = v711
	var v712 int32
	_ = v712
	var v719 int32
	_ = v719
	var v755 int32
	_ = v755
	var v766 int32
	_ = v766
	var v771 int32
	_ = v771
	var v772 int32
	_ = v772
	var v774 int32
	_ = v774
	var v783 int32
	_ = v783
	var v810 int32
	_ = v810
	var v812 int32
	_ = v812
	var v813 int32
	_ = v813
	var v821 int32
	_ = v821
	var v822 int32
	_ = v822
	var v824 int32
	_ = v824
	var v833 int32
	_ = v833
	var v835 int32
	_ = v835
	var v837 int32
	_ = v837
	var v881 int32
	_ = v881
	var v886 int32
	_ = v886
	var v889 int32
	_ = v889
	var v892 int32
	_ = v892
	var v895 int32
	_ = v895
	var v896 int32
	_ = v896
	var v897 int32
	_ = v897
	var v901 int32
	_ = v901
	var v902 int32
	_ = v902
	var v903 int32
	_ = v903
	var v908 int32
	_ = v908
	var v909 int32
	_ = v909
	var v916 int32
	_ = v916
	var v952 int32
	_ = v952
	var v963 int32
	_ = v963
	var v968 int32
	_ = v968
	var v969 int32
	_ = v969
	var v971 int32
	_ = v971
	var v980 int32
	_ = v980
	var v1007 int32
	_ = v1007
	var v1009 int32
	_ = v1009
	var v1010 int32
	_ = v1010
	var v1018 int32
	_ = v1018
	var v1019 int32
	_ = v1019
	var v1021 int32
	_ = v1021
	var v1030 int32
	_ = v1030
	var v1032 int32
	_ = v1032
	var v1034 int32
	_ = v1034
	var v1078 int32
	_ = v1078
	var v1083 int32
	_ = v1083
	var v1086 int32
	_ = v1086
	var v1089 int32
	_ = v1089
	var v1092 int32
	_ = v1092
	var v1099 int32
	_ = v1099
	var v1104 int32
	_ = v1104
	var v1105 int32
	_ = v1105
	var v1107 int32
	_ = v1107
	var v1109 int32
	_ = v1109
	var v1111 int32
	_ = v1111
	var v1113 int32
	_ = v1113
	var v1115 int32
	_ = v1115
	var v1117 int32
	_ = v1117
	var v1119 int32
	_ = v1119
	var v1121 int32
	_ = v1121
	var v1123 int32
	_ = v1123
	var v1125 int32
	_ = v1125
	var v1127 int32
	_ = v1127
	var v1129 int32
	_ = v1129
	var v1131 int32
	_ = v1131
	var v1133 int32
	_ = v1133
	var v1135 int32
	_ = v1135
	var v1137 int32
	_ = v1137
	var v1140 int32
	_ = v1140
	var v1233 int32
	_ = v1233
	var v1240 int32
	_ = v1240
	var v1286 int64
	_ = v1286
	v13 = m.G0
	v15 = v13 - int32(16)
	m.G0 = v15
	v18 = float32(100)
	if base.F32_gt(l1, v18) != 0 {
		v21 = v18
	} else {
		v21 = l1
	}
	if base.F32_lt(l1, float32(0)) != 0 {
		v24 = float32(0)
	} else {
		v24 = v21
	}
	F_VP8SetSegmentParams(m, l0, v24)
	mBase = m.M
	v26 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v15)+8)) = v26
	*(*int64)(unsafe.Add(mBase, uint32(v15))) = v26
	v30 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v31 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v32 = v30 * v31
	if v32 < int32(1) {
	} else {
		v35 = int32(1)
		v37 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_SetLoopParams[0])))
		if v32 == v35 {
			v89 = int32(0)
		} else {
			v47 = v37
			v50 = int32(0)
			for {
				v56 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v47))))
				v57 = int32(3)
				v59 = int32(12)
				v61 = v15 | int32(base.Ui32(v56)>>(uint(v57)%32))&v59
				v62 = *(*int32)(unsafe.Add(mBase, uint32(v61)))
				v63 = int32(1)
				*(*int32)(unsafe.Add(mBase, uint32(v61))) = v62 + v63
				v68 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v47+int32(4)))))
				v73 = v15 | int32(base.Ui32(v68)>>(uint(v57)%32))&v59
				v74 = *(*int32)(unsafe.Add(mBase, uint32(v73)))
				*(*int32)(unsafe.Add(mBase, uint32(v73))) = v74 + v63
				v81 = v50 + int32(2)
				if v32&int32(2147483646) != v81 {
					v47 = v47 + int32(8)
					v50 = v81
					continue
				} else {
					break
				}
				break
			}
			v89 = v81
		}
		if v32&v35 == int32(0) {
		} else {
			v100 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v37+v89<<(uint(int32(2))%32)))))
			v105 = v15 | int32(base.Ui32(v100)>>(uint(int32(3))%32))&int32(12)
			v106 = *(*int32)(unsafe.Add(mBase, uint32(v105)))
			*(*int32)(unsafe.Add(mBase, uint32(v105))) = v106 + int32(1)
		}
	}
	v122 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v123 = *(*int32)(unsafe.Add(mBase, uint32(v122)+88))
	if v123 == int32(0) {
	} else {
		v126 = *(*int64)(unsafe.Add(mBase, uint32(v15)))
		*(*int64)(unsafe.Add(mBase, uint32(v123)+92)) = v126
		v130 = *(*int64)(unsafe.Add(mBase, uint32(v15)+8))
		*(*int64)(unsafe.Add(mBase, uint32(v123+int32(100)))) = v130
	}
	v132 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	if v132 < int32(2) {
		v224 = int32(0)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = v224
		v300 = v224
	} else {
		v135 = int32(255)
		v137 = *(*int32)(unsafe.Add(mBase, uint32(v15)+12))
		v138 = *(*int32)(unsafe.Add(mBase, uint32(v15)+8))
		v139 = v137 + v138
		v140 = *(*int32)(unsafe.Add(mBase, uint32(v15)+4))
		v141 = *(*int32)(unsafe.Add(mBase, uint32(v15)))
		v142 = v140 + v141
		v143 = v139 + v142
		if v143 == int32(0) {
			v152 = v135
		} else {
			v147 = base.I32_div_s(v143, int32(2))
			v151 = base.I32_div_s(v147+v142*int32(255), v143)
			v152 = v151
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l0)+3416)) = uint8(v152)
		if v142 == int32(0) {
			v162 = v135
		} else {
			v157 = base.I32_div_s(v142, int32(2))
			v161 = base.I32_div_s(v157+v141*int32(255), v142)
			v162 = v161
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l0)+3417)) = uint8(v162)
		if v139 == int32(0) {
			v173 = int32(255)
		} else {
			v168 = base.I32_div_s(v139, int32(2))
			v172 = base.I32_div_s(v168+v138*int32(255), v139)
			v173 = v172
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l0)+3418)) = uint8(v173)
		v175 = int32(255)
		if v152&v175 != v175 {
			*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
			v234 = v152
			v238 = v173
			v239 = v162
		} else {
			v179 = int32(255)
			if v162&v179 != v179 {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
				v234 = v152
				v238 = v173
				v239 = v162
			} else {
				v183 = int32(255)
				v186 = base.B2i32(v173&v183 != v183)
				*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = v186
				if v173&v183 != v183 {
					v229 = int32(255)
					v234 = v229
					v238 = v173
					v239 = v229
				} else {
					v188 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
					v189 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
					if v188*v189 < int32(1) {
						v229 = int32(255)
						v234 = v229
						v238 = v173
						v239 = v229
					} else {
						v193 = int32(0)
						v198 = v193
						v201 = v193
						for {
							v207 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_SetLoopParams[0])))
							v208 = v207 + v198
							v209 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v208))))
							v211 = v209 & int32(159)
							*(*uint8)(unsafe.Add(mBase, uint32(v208))) = uint8(v211)
							v216 = v201 + int32(1)
							v217 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
							v218 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
							if v216 < v217*v218 {
								v198 = v198 + int32(4)
								v201 = v216
								continue
							} else {
								break
							}
							break
						}
						v221 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3418)))
						v222 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3417)))
						v223 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3416)))
						v234 = v223
						v238 = v221
						v239 = v222
					}
				}
			}
		}
		v243 = m.G24
		v244 = int32(-1)
		v246 = int32(255)
		v248 = int32(1)
		v251 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v243+(v239^v244)&v246<<(uint(v248)%32)))))
		v257 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v243+v234&v246<<(uint(v248)%32)))))
		v265 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v243+v239&v246<<(uint(v248)%32)))))
		v274 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v243+v238&v246<<(uint(v248)%32)))))
		v282 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v243+(v234^v244)&v246<<(uint(v248)%32)))))
		v293 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v243+(v238^v244)&v246<<(uint(v248)%32)))))
		v300 = (v251+v257)*v140 + (v265+v257)*v141 + (v274+v282)*v138 + (v293+v282)*v137
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = v300
	v352 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_SetLoopParams[1])))
	if v352 == int32(0) {
	} else {
		v367 = m.G23
		v368 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367)+15)))
		v369 = int32(408)
		v371 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367)+14)))
		v374 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367)+13)))
		v377 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367)+12)))
		v380 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367)+11)))
		v383 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367)+10)))
		v386 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367)+9)))
		v389 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367)+8)))
		v392 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367)+7)))
		v395 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367)+6)))
		v398 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367)+5)))
		v401 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367)+4)))
		v404 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367)+3)))
		v407 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367)+2)))
		v410 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367)+1)))
		v413 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367))))
		v418 = l0 + int32(3444)
		v419 = l0 + int32(3433)
		v420 = l0 + int32(3422)
		v441 = int32(0)
		for {
			v460 = l0 + int32(_a_F_SetLoopParams_0) + v441*int32(3264)
			v492 = v418
			v493 = v419
			v494 = v420
			v495 = int32(0)
			for {
				v508 = l0 + int32(3420) + v441*int32(264) + v495*int32(33)
				v509 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v508)+1)))
				v510 = int32(1)
				v513 = v460 + v495*int32(408)
				v514 = m.G24
				v518 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v514+v509<<(uint(v510)%32)))))
				*(*uint16)(unsafe.Add(mBase, uint32(v513))) = uint16(v518)
				v525 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v514+(v509^int32(255))<<(uint(v510)%32)))))
				v558 = v510
				for {
					v569 = m.G1
					v574 = v558<<(uint(int32(2))%32) + (v569 + int32(_a_F_SetLoopParams_1)) + int32(-4)
					v575 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v574))))
					if v575 != 0 {
						v577 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v574)+2)))
						v586 = v575
						v613 = v577
						v615 = int32(0)
						v616 = v494
						for {
							if v586&int32(1) == int32(0) {
								v638 = v615
							} else {
								v624 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v616))))
								v625 = m.G24
								v627 = int32(1)
								v636 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v625+(v624^(int32(0)-v613&v627)&int32(255))<<(uint(v627)%32)))))
								v638 = v615 + v636
							}
							v640 = int32(1)
							if base.Ui32(v640) < base.Ui32(v586) {
								v586 = int32(base.Ui32(v586) >> (uint(v640) % 32))
								v613 = int32(base.Ui32(v613) >> (uint(v640) % 32))
								v615 = v638
								v616 = v616 + v640
								continue
							} else {
								break
							}
							break
						}
						v684 = v638
					} else {
						v684 = int32(0)
					}
					v689 = int32(1)
					v692 = v525 + v684
					*(*uint16)(unsafe.Add(mBase, uint32(v513+v558<<(uint(v689)%32)))) = uint16(v692)
					v695 = v558 + v689
					if v695 != int32(68) {
						v558 = v695
						continue
					} else {
						break
					}
					break
				}
				v698 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v508)+12)))
				v699 = int32(1)
				v700 = m.G24
				v704 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v700+v698<<(uint(v699)%32)))))
				v705 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v508)+11)))
				v706 = int32(255)
				v711 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v700+(v705^v706)<<(uint(v699)%32)))))
				v712 = v704 + v711
				*(*uint16)(unsafe.Add(mBase, uint32(v513)+136)) = uint16(v712)
				v719 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v700+(v698^v706)<<(uint(v699)%32)))))
				v755 = v699
				for {
					v766 = m.G1
					v771 = v755<<(uint(int32(2))%32) + (v766 + int32(_a_F_SetLoopParams_1)) + int32(-4)
					v772 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v771))))
					if v772 != 0 {
						v774 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v771)+2)))
						v783 = v772
						v810 = v774
						v812 = int32(0)
						v813 = v493
						for {
							if v783&int32(1) == int32(0) {
								v835 = v812
							} else {
								v821 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v813))))
								v822 = m.G24
								v824 = int32(1)
								v833 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v822+(v821^(int32(0)-v810&v824)&int32(255))<<(uint(v824)%32)))))
								v835 = v812 + v833
							}
							v837 = int32(1)
							if base.Ui32(v837) < base.Ui32(v783) {
								v783 = int32(base.Ui32(v783) >> (uint(v837) % 32))
								v810 = int32(base.Ui32(v810) >> (uint(v837) % 32))
								v812 = v835
								v813 = v813 + v837
								continue
							} else {
								break
							}
							break
						}
						v881 = v835
					} else {
						v881 = int32(0)
					}
					v886 = int32(1)
					v889 = v711 + v719 + v881
					*(*uint16)(unsafe.Add(mBase, uint32(v513+int32(136)+v755<<(uint(v886)%32)))) = uint16(v889)
					v892 = v755 + v886
					if v892 != int32(68) {
						v755 = v892
						continue
					} else {
						break
					}
					break
				}
				v895 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v508)+23)))
				v896 = int32(1)
				v897 = m.G24
				v901 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v897+v895<<(uint(v896)%32)))))
				v902 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v508)+22)))
				v903 = int32(255)
				v908 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v897+(v902^v903)<<(uint(v896)%32)))))
				v909 = v901 + v908
				*(*uint16)(unsafe.Add(mBase, uint32(v513)+272)) = uint16(v909)
				v916 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v897+(v895^v903)<<(uint(v896)%32)))))
				v952 = v896
				for {
					v963 = m.G1
					v968 = v952<<(uint(int32(2))%32) + (v963 + int32(_a_F_SetLoopParams_1)) + int32(-4)
					v969 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v968))))
					if v969 != 0 {
						v971 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v968)+2)))
						v980 = v969
						v1007 = v971
						v1009 = int32(0)
						v1010 = v492
						for {
							if v980&int32(1) == int32(0) {
								v1032 = v1009
							} else {
								v1018 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1010))))
								v1019 = m.G24
								v1021 = int32(1)
								v1030 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1019+(v1018^(int32(0)-v1007&v1021)&int32(255))<<(uint(v1021)%32)))))
								v1032 = v1009 + v1030
							}
							v1034 = int32(1)
							if base.Ui32(v1034) < base.Ui32(v980) {
								v980 = int32(base.Ui32(v980) >> (uint(v1034) % 32))
								v1007 = int32(base.Ui32(v1007) >> (uint(v1034) % 32))
								v1009 = v1032
								v1010 = v1010 + v1034
								continue
							} else {
								break
							}
							break
						}
						v1078 = v1032
					} else {
						v1078 = int32(0)
					}
					v1083 = int32(1)
					v1086 = v908 + v916 + v1078
					*(*uint16)(unsafe.Add(mBase, uint32(v513+int32(272)+v952<<(uint(v1083)%32)))) = uint16(v1086)
					v1089 = v952 + v1083
					if v1089 != int32(68) {
						v952 = v1089
						continue
					} else {
						break
					}
					break
				}
				v1092 = int32(33)
				v1099 = v495 + int32(1)
				if v1099 != int32(8) {
					v492 = v492 + v1092
					v493 = v493 + v1092
					v494 = v494 + v1092
					v495 = v1099
					continue
				} else {
					break
				}
				break
			}
			v1104 = l0 + int32(_a_F_SetLoopParams_2) + v441*int32(192)
			v1105 = v460 + v368*v369
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+180)) = v1105
			v1107 = v460 + v371*v369
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+168)) = v1107
			v1109 = v460 + v374*v369
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+156)) = v1109
			v1111 = v460 + v377*v369
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+144)) = v1111
			v1113 = v460 + v380*v369
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+132)) = v1113
			v1115 = v460 + v383*v369
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+120)) = v1115
			v1117 = v460 + v386*v369
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+108)) = v1117
			v1119 = v460 + v389*v369
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+96)) = v1119
			v1121 = v460 + v392*v369
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+84)) = v1121
			v1123 = v460 + v395*v369
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+72)) = v1123
			v1125 = v460 + v398*v369
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+60)) = v1125
			v1127 = v460 + v401*v369
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+48)) = v1127
			v1129 = v460 + v404*v369
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+36)) = v1129
			v1131 = v460 + v407*v369
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+24)) = v1131
			v1133 = v460 + v410*v369
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+12)) = v1133
			v1135 = v460 + v413*v369
			*(*int32)(unsafe.Add(mBase, uint32(v1104))) = v1135
			v1137 = int32(272)
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+188)) = v1105 + v1137
			v1140 = int32(136)
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+184)) = v1105 + v1140
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+176)) = v1107 + v1137
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+172)) = v1107 + v1140
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+164)) = v1109 + v1137
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+160)) = v1109 + v1140
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+152)) = v1111 + v1137
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+148)) = v1111 + v1140
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+140)) = v1113 + v1137
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+136)) = v1113 + v1140
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+128)) = v1115 + v1137
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+124)) = v1115 + v1140
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+116)) = v1117 + v1137
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+112)) = v1117 + v1140
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+104)) = v1119 + v1137
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+100)) = v1119 + v1140
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+92)) = v1121 + v1137
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+88)) = v1121 + v1140
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+80)) = v1123 + v1137
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+76)) = v1123 + v1140
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+68)) = v1125 + v1137
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+64)) = v1125 + v1140
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+56)) = v1127 + v1137
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+52)) = v1127 + v1140
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+44)) = v1129 + v1137
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+40)) = v1129 + v1140
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+32)) = v1131 + v1137
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+28)) = v1131 + v1140
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+20)) = v1133 + v1137
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+16)) = v1133 + v1140
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+8)) = v1135 + v1137
			*(*int32)(unsafe.Add(mBase, uint32(v1104)+4)) = v1135 + v1140
			v1233 = int32(264)
			v1240 = v441 + int32(1)
			if v1240 != int32(4) {
				v418 = v418 + v1233
				v419 = v419 + v1233
				v420 = v420 + v1233
				v441 = v1240
				continue
			} else {
				break
			}
			break
		}
		*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_SetLoopParams[1]))) = int32(0)
	}
	v1286 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_SetLoopParams[2]))) = v1286
	*(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_SetLoopParams[3]))) = v1286
	*(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_SetLoopParams[4]))) = v1286
	*(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_SetLoopParams[5]))) = v1286
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_SetLoopParams[6]))) = int32(0)
	m.G0 = v15 + int32(16)
	return
}
func F_SetResidualCoeffs_C(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v10 int32
	_ = v10
	var v14 int32
	_ = v14
	var v18 int32
	_ = v18
	var v22 int32
	_ = v22
	var v26 int32
	_ = v26
	var v30 int32
	_ = v30
	var v34 int32
	_ = v34
	var v38 int32
	_ = v38
	var v42 int32
	_ = v42
	var v46 int32
	_ = v46
	var v50 int32
	_ = v50
	var v54 int32
	_ = v54
	var v58 int32
	_ = v58
	var v62 int32
	_ = v62
	var v66 int32
	_ = v66
	var v70 int32
	_ = v70
	*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = int32(-1)
	v6 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+30)))
	if v6 == int32(0) {
		v10 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+28)))
		if v10 == int32(0) {
			v14 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+26)))
			if v14 == int32(0) {
				v18 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+24)))
				if v18 == int32(0) {
					v22 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+22)))
					if v22 == int32(0) {
						v26 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+20)))
						if v26 == int32(0) {
							v30 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+18)))
							if v30 == int32(0) {
								v34 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+16)))
								if v34 == int32(0) {
									v38 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+14)))
									if v38 == int32(0) {
										v42 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+12)))
										if v42 == int32(0) {
											v46 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+10)))
											if v46 == int32(0) {
												v50 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+8)))
												if v50 == int32(0) {
													v54 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+6)))
													if v54 == int32(0) {
														v58 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+4)))
														if v58 == int32(0) {
															v62 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+2)))
															if v62 == int32(0) {
																v66 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0))))
																if v66 == int32(0) {
																} else {
																	v70 = int32(0)
																	*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v70
																}
															} else {
																v70 = int32(1)
																*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v70
															}
														} else {
															v70 = int32(2)
															*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v70
														}
													} else {
														v70 = int32(3)
														*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v70
													}
												} else {
													v70 = int32(4)
													*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v70
												}
											} else {
												v70 = int32(5)
												*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v70
											}
										} else {
											v70 = int32(6)
											*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v70
										}
									} else {
										v70 = int32(7)
										*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v70
									}
								} else {
									v70 = int32(8)
									*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v70
								}
							} else {
								v70 = int32(9)
								*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v70
							}
						} else {
							v70 = int32(10)
							*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v70
						}
					} else {
						v70 = int32(11)
						*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v70
					}
				} else {
					v70 = int32(12)
					*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v70
				}
			} else {
				v70 = int32(13)
				*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v70
			}
		} else {
			v70 = int32(14)
			*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v70
		}
	} else {
		v70 = int32(15)
		*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v70
	}
	*(*int32)(unsafe.Add(mBase, uint32(l1)+8)) = l0
	return
}
func F_ShannonEntropy_C(m *base.Module, l0 int32, l1 int32) int64 {
	mBase := m.M
	_ = mBase
	var v3 int64
	_ = v3
	var v4 int32
	_ = v4
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int64
	_ = v13
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v22 int32
	_ = v22
	var v28 int64
	_ = v28
	var v29 int32
	_ = v29
	var v32 int32
	_ = v32
	var v33 int64
	_ = v33
	var v34 int64
	_ = v34
	var v37 int64
	_ = v37
	var v38 int32
	_ = v38
	var v43 int32
	_ = v43
	var v48 int64
	_ = v48
	var v49 int32
	_ = v49
	var v52 int32
	_ = v52
	var v58 int64
	_ = v58
	var v61 int32
	_ = v61
	var v64 int32
	_ = v64
	var v65 int64
	_ = v65
	v3 = int64(0)
	v4 = int32(0)
	if l1 < int32(1) {
		v48 = v3
		v49 = v4
		v52 = m.G1
		v58 = *(*int64)(unsafe.Add(mBase, uint32(v52+int32(_a_F_ShannonEntropy_C_0)+v49<<(uint(int32(3))%32))))
		return v58 - v48
	} else {
		v11 = l0
		v12 = l1
		v13 = v3
		v14 = v4
		for {
			v17 = *(*int32)(unsafe.Add(mBase, uint32(v11)))
			if v17 == int32(0) {
				v37 = v13
				v38 = v14
			} else {
				if base.Ui32(int32(255)) < base.Ui32(v17) {
					v29 = m.G1
					v32 = *(*int32)(unsafe.Add(mBase, uint32(v29)+uint32(_c_F_ShannonEntropy_C[0])))
					v33 = m.T0[v32].(func(*base.Module, int32) int64)(m, v17)
					mBase = m.M
					v34 = v33
				} else {
					v22 = m.G1
					v28 = *(*int64)(unsafe.Add(mBase, uint32(v22+int32(_a_F_ShannonEntropy_C_0)+v17<<(uint(int32(3))%32))))
					v34 = v28
				}
				v37 = v34 + v13
				v38 = v17 + v14
			}
			v43 = v12 + int32(-1)
			if v43 != 0 {
				v11 = v11 + int32(4)
				v12 = v43
				v13 = v37
				v14 = v38
				continue
			} else {
				break
			}
			break
		}
		if base.Ui32(int32(255)) < base.Ui32(v38) {
			v61 = m.G1
			v64 = *(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_ShannonEntropy_C[0])))
			v65 = m.T0[v64].(func(*base.Module, int32) int64)(m, v38)
			mBase = m.M
			return v65 - v37
		} else {
			v48 = v37
			v49 = v38
			v52 = m.G1
			v58 = *(*int64)(unsafe.Add(mBase, uint32(v52+int32(_a_F_ShannonEntropy_C_0)+v49<<(uint(int32(3))%32))))
			return v58 - v48
		}
	}
}
func F_SharpYuvFilterRow_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v30 int32
	_ = v30
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
	var v48 int32
	_ = v48
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v69 int32
	_ = v69
	var v70 int32
	_ = v70
	var v72 int32
	_ = v72
	var v75 int32
	_ = v75
	var v77 int32
	_ = v77
	var v81 int32
	_ = v81
	var v86 int32
	_ = v86
	if l2 < int32(1) {
	} else {
		v14 = int32(-1)
		v17 = v14<<(uint(l5)%32) ^ v14
		v18 = l0
		v19 = l1
		v20 = l2
		v21 = l3
		v22 = l4
		for {
			v30 = int32(*(*int16)(unsafe.Add(mBase, uint32(v18))))
			v33 = int32(*(*int16)(unsafe.Add(mBase, uint32(v19)+2)))
			v35 = int32(*(*int16)(unsafe.Add(mBase, uint32(v19))))
			v36 = int32(*(*int16)(unsafe.Add(mBase, uint32(v18)+2)))
			v45 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v21))))
			v46 = (v30*int32(9)+v33+(v35+v36)*int32(3)+int32(8))>>(uint(int32(4))%32) + v45
			if v46 < v17 {
				v48 = v46
			} else {
				v48 = v17
			}
			if v46 < int32(0) {
				v51 = int32(0)
			} else {
				v51 = v48
			}
			*(*uint16)(unsafe.Add(mBase, uint32(v22))) = uint16(v51)
			v53 = int32(2)
			v69 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v21+v53))))
			v70 = (v35+v36*int32(9)+(v33+v30)*int32(3)+int32(8))>>(uint(int32(4))%32) + v69
			if v70 < v17 {
				v72 = v70
			} else {
				v72 = v17
			}
			if v70 < int32(0) {
				v75 = int32(0)
			} else {
				v75 = v72
			}
			*(*uint16)(unsafe.Add(mBase, uint32(v22+v53))) = uint16(v75)
			v77 = int32(4)
			v81 = int32(2)
			v86 = v20 + int32(-1)
			if v86 != 0 {
				v18 = v18 + v81
				v19 = v19 + v81
				v20 = v86
				v21 = v21 + v77
				v22 = v22 + v77
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_SharpYuvUpdateRGB_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v13 int32
	_ = v13
	var v24 int32
	_ = v24
	var v26 int32
	_ = v26
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v41 int32
	_ = v41
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
	var v63 int32
	_ = v63
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	if l3 < int32(1) {
	} else {
		v13 = int32(1)
		if l3 == v13 {
			v63 = int32(0)
		} else {
			v24 = l2
			v26 = int32(0)
			v28 = l1
			v29 = l0
			for {
				v31 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v29))))
				v32 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v28))))
				v34 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v24))))
				v35 = v31 - v32 + v34
				*(*uint16)(unsafe.Add(mBase, uint32(v24))) = uint16(v35)
				v37 = int32(2)
				v38 = v24 + v37
				v41 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v29+v37))))
				v44 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v28+v37))))
				v46 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v38))))
				v47 = v41 - v44 + v46
				*(*uint16)(unsafe.Add(mBase, uint32(v38))) = uint16(v47)
				v49 = int32(4)
				v56 = v26 + v37
				if l3&int32(2147483646) != v56 {
					v24 = v24 + v49
					v26 = v56
					v28 = v28 + v49
					v29 = v29 + v49
					continue
				} else {
					break
				}
				break
			}
			v63 = v56
		}
		if l3&v13 == int32(0) {
		} else {
			v71 = v63 << (uint(int32(1)) % 32)
			v72 = l2 + v71
			v74 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0+v71))))
			v76 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l1+v71))))
			v78 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v72))))
			v79 = v74 - v76 + v78
			*(*uint16)(unsafe.Add(mBase, uint32(v72))) = uint16(v79)
		}
	}
	return
}
func F_SharpYuvUpdateY_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) int64 {
	mBase := m.M
	_ = mBase
	var v13 int32
	_ = v13
	var v16 int32
	_ = v16
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v24 int64
	_ = v24
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v36 int32
	_ = v36
	var v39 int32
	_ = v39
	var v43 int64
	_ = v43
	var v44 int32
	_ = v44
	var v51 int32
	_ = v51
	if int32(1) <= l3 {
		v13 = int32(-1)
		v16 = v13<<(uint(l4)%32) ^ v13
		v18 = l0
		v19 = l1
		v20 = l2
		v21 = l3
		v24 = int64(0)
		for {
			v27 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v18))))
			v28 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v19))))
			v29 = v27 - v28
			v30 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v20))))
			v31 = v29 + v30
			if v31 < v16 {
				v33 = v31
			} else {
				v33 = v16
			}
			if v31 < int32(0) {
				v36 = int32(0)
			} else {
				v36 = v33
			}
			*(*uint16)(unsafe.Add(mBase, uint32(v20))) = uint16(v36)
			v39 = v29 >> (uint(int32(31)) % 32)
			v43 = v24 + base.I64_extend_i32_u(v29^v39-v39)
			v44 = int32(2)
			v51 = v21 + int32(-1)
			if v51 != 0 {
				v18 = v18 + v44
				v19 = v19 + v44
				v20 = v20 + v44
				v21 = v51
				v24 = v43
				continue
			} else {
				break
			}
			break
		}
		return v43
	} else {
		return int64(0)
	}
}
func F_SimpleHFilter16_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v17 int32
	_ = v17
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v37 int32
	_ = v37
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
	var v52 int32
	_ = v52
	var v55 int32
	_ = v55
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	var v65 int32
	_ = v65
	var v68 int32
	_ = v68
	var v69 int32
	_ = v69
	var v75 int32
	_ = v75
	var v82 int32
	_ = v82
	var v84 int32
	_ = v84
	var v88 int32
	_ = v88
	var v94 int32
	_ = v94
	v17 = int32(1)
	v21 = m.G12
	v22 = *(*int32)(unsafe.Add(mBase, uint32(v21)))
	v23 = m.G13
	v24 = *(*int32)(unsafe.Add(mBase, uint32(v23)))
	v25 = m.G14
	v26 = *(*int32)(unsafe.Add(mBase, uint32(v25)))
	v27 = m.G15
	v28 = *(*int32)(unsafe.Add(mBase, uint32(v27)))
	v30 = l0 + int32(-2)
	v37 = int32(16)
	for {
		v45 = v30 + int32(1)
		v46 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v45))))
		v47 = int32(2)
		v48 = v30 + v47
		v49 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48))))
		v52 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28+(v46-v49)))))
		v55 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v30))))
		v58 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v30+int32(3)))))
		v59 = v55 - v58
		v61 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28+v59))))
		if l2<<(uint(v17)%32)|v17 < v52<<(uint(v47)%32)+v61 {
		} else {
			v65 = int32(3)
			v68 = int32(*(*int8)(unsafe.Add(mBase, uint32(v26+v59))))
			v69 = (v49-v46)*v65 + v68
			v75 = int32(*(*int8)(unsafe.Add(mBase, uint32(v24+(v69+int32(4))>>(uint(v65)%32)))))
			v82 = int32(*(*int8)(unsafe.Add(mBase, uint32(v24+(v69+v65)>>(uint(v65)%32)))))
			v84 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v22+v46+v82))))
			*(*uint8)(unsafe.Add(mBase, uint32(v45))) = uint8(v84)
			v88 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v22+(v49-v75)))))
			*(*uint8)(unsafe.Add(mBase, uint32(v48))) = uint8(v88)
		}
		v94 = v37 + int32(-1)
		if v94 != 0 {
			v30 = v30 + l1
			v37 = v94
			continue
		} else {
			break
		}
		break
	}
	return
}
func F_SimpleHFilter16i_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v18 int32
	_ = v18
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v34 int32
	_ = v34
	var v39 int32
	_ = v39
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v54 int32
	_ = v54
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v67 int32
	_ = v67
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v77 int32
	_ = v77
	var v84 int32
	_ = v84
	var v86 int32
	_ = v86
	var v90 int32
	_ = v90
	var v96 int32
	_ = v96
	var v103 int32
	_ = v103
	var v108 int32
	_ = v108
	var v116 int32
	_ = v116
	var v117 int32
	_ = v117
	var v119 int32
	_ = v119
	var v120 int32
	_ = v120
	var v123 int32
	_ = v123
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v132 int32
	_ = v132
	var v136 int32
	_ = v136
	var v139 int32
	_ = v139
	var v140 int32
	_ = v140
	var v146 int32
	_ = v146
	var v153 int32
	_ = v153
	var v155 int32
	_ = v155
	var v159 int32
	_ = v159
	var v165 int32
	_ = v165
	var v172 int32
	_ = v172
	var v177 int32
	_ = v177
	var v185 int32
	_ = v185
	var v186 int32
	_ = v186
	var v188 int32
	_ = v188
	var v189 int32
	_ = v189
	var v192 int32
	_ = v192
	var v197 int32
	_ = v197
	var v198 int32
	_ = v198
	var v199 int32
	_ = v199
	var v201 int32
	_ = v201
	var v205 int32
	_ = v205
	var v208 int32
	_ = v208
	var v209 int32
	_ = v209
	var v215 int32
	_ = v215
	var v222 int32
	_ = v222
	var v224 int32
	_ = v224
	var v228 int32
	_ = v228
	var v234 int32
	_ = v234
	v18 = int32(1)
	v21 = l2<<(uint(v18)%32) | v18
	v22 = m.G12
	v23 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
	v24 = m.G13
	v25 = *(*int32)(unsafe.Add(mBase, uint32(v24)))
	v26 = m.G14
	v27 = *(*int32)(unsafe.Add(mBase, uint32(v26)))
	v28 = m.G15
	v29 = *(*int32)(unsafe.Add(mBase, uint32(v28)))
	v34 = l0 + int32(5)
	v39 = int32(16)
	for {
		v47 = v34 + int32(-2)
		v48 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v47))))
		v50 = v34 + int32(-1)
		v51 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v50))))
		v54 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v29+(v48-v51)))))
		v59 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v34+int32(-3)))))
		v60 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v34))))
		v61 = v59 - v60
		v63 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v29+v61))))
		if v21 < v54<<(uint(int32(2))%32)+v63 {
		} else {
			v67 = int32(3)
			v70 = int32(*(*int8)(unsafe.Add(mBase, uint32(v27+v61))))
			v71 = (v51-v48)*v67 + v70
			v77 = int32(*(*int8)(unsafe.Add(mBase, uint32(v25+(v71+int32(4))>>(uint(v67)%32)))))
			v84 = int32(*(*int8)(unsafe.Add(mBase, uint32(v25+(v71+v67)>>(uint(v67)%32)))))
			v86 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v23+v48+v84))))
			*(*uint8)(unsafe.Add(mBase, uint32(v47))) = uint8(v86)
			v90 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v23+(v51-v77)))))
			*(*uint8)(unsafe.Add(mBase, uint32(v50))) = uint8(v90)
		}
		v96 = v39 + int32(-1)
		if v96 != 0 {
			v34 = v34 + l1
			v39 = v96
			continue
		} else {
			break
		}
		break
	}
	v103 = l0 + int32(9)
	v108 = int32(16)
	for {
		v116 = v103 + int32(-2)
		v117 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v116))))
		v119 = v103 + int32(-1)
		v120 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v119))))
		v123 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v29+(v117-v120)))))
		v128 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v103+int32(-3)))))
		v129 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v103))))
		v130 = v128 - v129
		v132 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v29+v130))))
		if v21 < v123<<(uint(int32(2))%32)+v132 {
		} else {
			v136 = int32(3)
			v139 = int32(*(*int8)(unsafe.Add(mBase, uint32(v27+v130))))
			v140 = (v120-v117)*v136 + v139
			v146 = int32(*(*int8)(unsafe.Add(mBase, uint32(v25+(v140+int32(4))>>(uint(v136)%32)))))
			v153 = int32(*(*int8)(unsafe.Add(mBase, uint32(v25+(v140+v136)>>(uint(v136)%32)))))
			v155 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v23+v117+v153))))
			*(*uint8)(unsafe.Add(mBase, uint32(v116))) = uint8(v155)
			v159 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v23+(v120-v146)))))
			*(*uint8)(unsafe.Add(mBase, uint32(v119))) = uint8(v159)
		}
		v165 = v108 + int32(-1)
		if v165 != 0 {
			v103 = v103 + l1
			v108 = v165
			continue
		} else {
			break
		}
		break
	}
	v172 = l0 + int32(13)
	v177 = int32(16)
	for {
		v185 = v172 + int32(-2)
		v186 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v185))))
		v188 = v172 + int32(-1)
		v189 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v188))))
		v192 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v29+(v186-v189)))))
		v197 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v172+int32(-3)))))
		v198 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v172))))
		v199 = v197 - v198
		v201 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v29+v199))))
		if v21 < v192<<(uint(int32(2))%32)+v201 {
		} else {
			v205 = int32(3)
			v208 = int32(*(*int8)(unsafe.Add(mBase, uint32(v27+v199))))
			v209 = (v189-v186)*v205 + v208
			v215 = int32(*(*int8)(unsafe.Add(mBase, uint32(v25+(v209+int32(4))>>(uint(v205)%32)))))
			v222 = int32(*(*int8)(unsafe.Add(mBase, uint32(v25+(v209+v205)>>(uint(v205)%32)))))
			v224 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v23+v186+v222))))
			*(*uint8)(unsafe.Add(mBase, uint32(v185))) = uint8(v224)
			v228 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v23+(v189-v215)))))
			*(*uint8)(unsafe.Add(mBase, uint32(v188))) = uint8(v228)
		}
		v234 = v177 + int32(-1)
		if v234 != 0 {
			v172 = v172 + l1
			v177 = v234
			continue
		} else {
			break
		}
		break
	}
	return
}
func F_SimpleVFilter16_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v19 int32
	_ = v19
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v36 int32
	_ = v36
	var v51 int32
	_ = v51
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v57 int32
	_ = v57
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
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
	var v80 int32
	_ = v80
	var v87 int32
	_ = v87
	var v89 int32
	_ = v89
	var v93 int32
	_ = v93
	var v98 int32
	_ = v98
	v19 = int32(1)
	v26 = m.G12
	v27 = *(*int32)(unsafe.Add(mBase, uint32(v26)))
	v28 = m.G13
	v29 = *(*int32)(unsafe.Add(mBase, uint32(v28)))
	v30 = m.G14
	v31 = *(*int32)(unsafe.Add(mBase, uint32(v30)))
	v32 = m.G15
	v33 = *(*int32)(unsafe.Add(mBase, uint32(v32)))
	v36 = int32(0)
	for {
		v51 = l0 - l1 + v36
		v52 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v51))))
		v53 = l0 + v36
		v54 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v53))))
		v57 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v33+(v52-v54)))))
		v61 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0-l1<<(uint(v19)%32)+v36))))
		v63 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+l1+v36))))
		v64 = v61 - v63
		v66 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v33+v64))))
		if l2<<(uint(v19)%32)|v19 < v57<<(uint(int32(2))%32)+v66 {
		} else {
			v70 = int32(3)
			v73 = int32(*(*int8)(unsafe.Add(mBase, uint32(v31+v64))))
			v74 = (v54-v52)*v70 + v73
			v80 = int32(*(*int8)(unsafe.Add(mBase, uint32(v29+(v74+int32(4))>>(uint(v70)%32)))))
			v87 = int32(*(*int8)(unsafe.Add(mBase, uint32(v29+(v74+v70)>>(uint(v70)%32)))))
			v89 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v27+v52+v87))))
			*(*uint8)(unsafe.Add(mBase, uint32(v51))) = uint8(v89)
			v93 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v27+(v54-v80)))))
			*(*uint8)(unsafe.Add(mBase, uint32(v53))) = uint8(v93)
		}
		v98 = v36 + int32(1)
		if v98 != int32(16) {
			v36 = v98
			continue
		} else {
			break
		}
		break
	}
	return
}
func F_SimpleVFilter16i_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v20 int32
	_ = v20
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
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
	var v57 int32
	_ = v57
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v67 int32
	_ = v67
	var v70 int32
	_ = v70
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v77 int32
	_ = v77
	var v79 int32
	_ = v79
	var v83 int32
	_ = v83
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v93 int32
	_ = v93
	var v100 int32
	_ = v100
	var v102 int32
	_ = v102
	var v106 int32
	_ = v106
	var v111 int32
	_ = v111
	var v114 int32
	_ = v114
	var v137 int32
	_ = v137
	var v144 int32
	_ = v144
	var v145 int32
	_ = v145
	var v146 int32
	_ = v146
	var v147 int32
	_ = v147
	var v150 int32
	_ = v150
	var v154 int32
	_ = v154
	var v156 int32
	_ = v156
	var v157 int32
	_ = v157
	var v159 int32
	_ = v159
	var v163 int32
	_ = v163
	var v166 int32
	_ = v166
	var v167 int32
	_ = v167
	var v173 int32
	_ = v173
	var v180 int32
	_ = v180
	var v182 int32
	_ = v182
	var v186 int32
	_ = v186
	var v191 int32
	_ = v191
	var v217 int32
	_ = v217
	var v224 int32
	_ = v224
	var v225 int32
	_ = v225
	var v226 int32
	_ = v226
	var v227 int32
	_ = v227
	var v230 int32
	_ = v230
	var v234 int32
	_ = v234
	var v236 int32
	_ = v236
	var v237 int32
	_ = v237
	var v239 int32
	_ = v239
	var v243 int32
	_ = v243
	var v246 int32
	_ = v246
	var v247 int32
	_ = v247
	var v253 int32
	_ = v253
	var v260 int32
	_ = v260
	var v262 int32
	_ = v262
	var v266 int32
	_ = v266
	var v271 int32
	_ = v271
	v20 = int32(1)
	v30 = l1 << (uint(int32(2)) % 32)
	v31 = l0 + v30
	v35 = l2<<(uint(v20)%32) | v20
	v36 = m.G12
	v37 = *(*int32)(unsafe.Add(mBase, uint32(v36)))
	v38 = m.G13
	v39 = *(*int32)(unsafe.Add(mBase, uint32(v38)))
	v40 = m.G14
	v41 = *(*int32)(unsafe.Add(mBase, uint32(v40)))
	v42 = m.G15
	v43 = *(*int32)(unsafe.Add(mBase, uint32(v42)))
	v57 = int32(0)
	for {
		v64 = l0 + l1*int32(3) + v57
		v65 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v64))))
		v66 = v31 + v57
		v67 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v66))))
		v70 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v43+(v65-v67)))))
		v74 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+l1<<(uint(v20)%32)+v57))))
		v76 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+l1*int32(5)+v57))))
		v77 = v74 - v76
		v79 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v43+v77))))
		if v35 < v70<<(uint(int32(2))%32)+v79 {
		} else {
			v83 = int32(3)
			v86 = int32(*(*int8)(unsafe.Add(mBase, uint32(v41+v77))))
			v87 = (v67-v65)*v83 + v86
			v93 = int32(*(*int8)(unsafe.Add(mBase, uint32(v39+(v87+int32(4))>>(uint(v83)%32)))))
			v100 = int32(*(*int8)(unsafe.Add(mBase, uint32(v39+(v87+v83)>>(uint(v83)%32)))))
			v102 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v37+v65+v100))))
			*(*uint8)(unsafe.Add(mBase, uint32(v64))) = uint8(v102)
			v106 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v37+(v67-v93)))))
			*(*uint8)(unsafe.Add(mBase, uint32(v66))) = uint8(v106)
		}
		v111 = v57 + int32(1)
		if v111 != int32(16) {
			v57 = v111
			continue
		} else {
			break
		}
		break
	}
	v114 = v31 + v30
	v137 = int32(0)
	for {
		v144 = l0 + l1*int32(7) + v137
		v145 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v144))))
		v146 = v114 + v137
		v147 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v146))))
		v150 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v43+(v145-v147)))))
		v154 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+l1*int32(6)+v137))))
		v156 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+l1*int32(9)+v137))))
		v157 = v154 - v156
		v159 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v43+v157))))
		if v35 < v150<<(uint(int32(2))%32)+v159 {
		} else {
			v163 = int32(3)
			v166 = int32(*(*int8)(unsafe.Add(mBase, uint32(v41+v157))))
			v167 = (v147-v145)*v163 + v166
			v173 = int32(*(*int8)(unsafe.Add(mBase, uint32(v39+(v167+int32(4))>>(uint(v163)%32)))))
			v180 = int32(*(*int8)(unsafe.Add(mBase, uint32(v39+(v167+v163)>>(uint(v163)%32)))))
			v182 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v37+v145+v180))))
			*(*uint8)(unsafe.Add(mBase, uint32(v144))) = uint8(v182)
			v186 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v37+(v147-v173)))))
			*(*uint8)(unsafe.Add(mBase, uint32(v146))) = uint8(v186)
		}
		v191 = v137 + int32(1)
		if v191 != int32(16) {
			v137 = v191
			continue
		} else {
			break
		}
		break
	}
	v217 = int32(0)
	for {
		v224 = l0 + l1*int32(11) + v217
		v225 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v224))))
		v226 = v114 + v30 + v217
		v227 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v226))))
		v230 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v43+(v225-v227)))))
		v234 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+l1*int32(10)+v217))))
		v236 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+l1*int32(13)+v217))))
		v237 = v234 - v236
		v239 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v43+v237))))
		if v35 < v230<<(uint(int32(2))%32)+v239 {
		} else {
			v243 = int32(3)
			v246 = int32(*(*int8)(unsafe.Add(mBase, uint32(v41+v237))))
			v247 = (v227-v225)*v243 + v246
			v253 = int32(*(*int8)(unsafe.Add(mBase, uint32(v39+(v247+int32(4))>>(uint(v243)%32)))))
			v260 = int32(*(*int8)(unsafe.Add(mBase, uint32(v39+(v247+v243)>>(uint(v243)%32)))))
			v262 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v37+v225+v260))))
			*(*uint8)(unsafe.Add(mBase, uint32(v224))) = uint8(v262)
			v266 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v37+(v227-v253)))))
			*(*uint8)(unsafe.Add(mBase, uint32(v226))) = uint8(v266)
		}
		v271 = v217 + int32(1)
		if v271 != int32(16) {
			v217 = v271
			continue
		} else {
			break
		}
		break
	}
	return
}
func F_StoreHuffmanCode(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	var v18 int32
	_ = v18
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v36 int32
	_ = v36
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v66 int32
	_ = v66
	var v67 int32
	_ = v67
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v84 int32
	_ = v84
	var v85 int32
	_ = v85
	var v93 int32
	_ = v93
	var v94 int32
	_ = v94
	var v95 int32
	_ = v95
	var v96 int32
	_ = v96
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v103 int32
	_ = v103
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v115 int32
	_ = v115
	var v116 int32
	_ = v116
	var v119 int64
	_ = v119
	var v122 int32
	_ = v122
	var v124 int64
	_ = v124
	var v127 int32
	_ = v127
	var v134 int32
	_ = v134
	var v136 int32
	_ = v136
	var v140 int32
	_ = v140
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v152 int32
	_ = v152
	var v153 int32
	_ = v153
	var v154 int32
	_ = v154
	var v156 int32
	_ = v156
	var v161 int32
	_ = v161
	var v162 int32
	_ = v162
	var v169 int32
	_ = v169
	var v172 int32
	_ = v172
	var v174 int32
	_ = v174
	var v178 int32
	_ = v178
	var v179 int32
	_ = v179
	var v180 int32
	_ = v180
	var v181 int32
	_ = v181
	var v220 int32
	_ = v220
	var v223 int32
	_ = v223
	var v226 int32
	_ = v226
	var v239 int32
	_ = v239
	var v240 int32
	_ = v240
	var v244 int32
	_ = v244
	var v245 int32
	_ = v245
	var v253 int32
	_ = v253
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
	var v263 int32
	_ = v263
	var v264 int32
	_ = v264
	var v265 int32
	_ = v265
	var v266 int32
	_ = v266
	var v275 int32
	_ = v275
	var v276 int32
	_ = v276
	var v279 int64
	_ = v279
	var v282 int32
	_ = v282
	var v284 int64
	_ = v284
	var v287 int32
	_ = v287
	var v294 int32
	_ = v294
	var v296 int32
	_ = v296
	var v300 int32
	_ = v300
	var v301 int32
	_ = v301
	var v302 int32
	_ = v302
	var v303 int32
	_ = v303
	var v312 int32
	_ = v312
	var v313 int32
	_ = v313
	var v314 int32
	_ = v314
	var v316 int32
	_ = v316
	var v321 int32
	_ = v321
	var v322 int32
	_ = v322
	var v329 int32
	_ = v329
	var v332 int32
	_ = v332
	var v334 int32
	_ = v334
	var v338 int32
	_ = v338
	var v339 int32
	_ = v339
	var v340 int32
	_ = v340
	var v341 int32
	_ = v341
	var v379 int32
	_ = v379
	var v380 int32
	_ = v380
	var v392 int32
	_ = v392
	var v393 int32
	_ = v393
	var v397 int32
	_ = v397
	var v398 int32
	_ = v398
	var v406 int32
	_ = v406
	var v407 int32
	_ = v407
	var v408 int32
	_ = v408
	var v409 int32
	_ = v409
	var v411 int32
	_ = v411
	var v412 int32
	_ = v412
	var v416 int32
	_ = v416
	var v417 int32
	_ = v417
	var v418 int32
	_ = v418
	var v419 int32
	_ = v419
	var v428 int32
	_ = v428
	var v429 int32
	_ = v429
	var v432 int64
	_ = v432
	var v435 int32
	_ = v435
	var v437 int64
	_ = v437
	var v440 int32
	_ = v440
	var v447 int32
	_ = v447
	var v449 int32
	_ = v449
	var v453 int32
	_ = v453
	var v454 int32
	_ = v454
	var v455 int32
	_ = v455
	var v456 int32
	_ = v456
	var v465 int32
	_ = v465
	var v466 int32
	_ = v466
	var v467 int32
	_ = v467
	var v469 int32
	_ = v469
	var v474 int32
	_ = v474
	var v475 int32
	_ = v475
	var v482 int32
	_ = v482
	var v485 int32
	_ = v485
	var v487 int32
	_ = v487
	var v491 int32
	_ = v491
	var v492 int32
	_ = v492
	var v493 int32
	_ = v493
	var v494 int32
	_ = v494
	var v531 int32
	_ = v531
	var v532 int32
	_ = v532
	var v545 int32
	_ = v545
	var v546 int32
	_ = v546
	var v550 int32
	_ = v550
	var v551 int32
	_ = v551
	var v559 int32
	_ = v559
	var v560 int32
	_ = v560
	var v561 int32
	_ = v561
	var v562 int32
	_ = v562
	var v564 int32
	_ = v564
	var v565 int32
	_ = v565
	var v569 int32
	_ = v569
	var v570 int32
	_ = v570
	var v571 int32
	_ = v571
	var v572 int32
	_ = v572
	var v581 int32
	_ = v581
	var v582 int32
	_ = v582
	var v585 int64
	_ = v585
	var v588 int32
	_ = v588
	var v590 int64
	_ = v590
	var v593 int32
	_ = v593
	var v600 int32
	_ = v600
	var v602 int32
	_ = v602
	var v606 int32
	_ = v606
	var v607 int32
	_ = v607
	var v608 int32
	_ = v608
	var v609 int32
	_ = v609
	var v618 int32
	_ = v618
	var v619 int32
	_ = v619
	var v620 int32
	_ = v620
	var v622 int32
	_ = v622
	var v627 int32
	_ = v627
	var v628 int32
	_ = v628
	var v635 int32
	_ = v635
	var v638 int32
	_ = v638
	var v640 int32
	_ = v640
	var v644 int32
	_ = v644
	var v645 int32
	_ = v645
	var v646 int32
	_ = v646
	var v647 int32
	_ = v647
	var v686 int32
	_ = v686
	var v698 int32
	_ = v698
	var v699 int32
	_ = v699
	var v703 int32
	_ = v703
	var v704 int32
	_ = v704
	var v712 int32
	_ = v712
	var v713 int32
	_ = v713
	var v714 int32
	_ = v714
	var v715 int32
	_ = v715
	var v717 int32
	_ = v717
	var v718 int32
	_ = v718
	var v722 int32
	_ = v722
	var v723 int32
	_ = v723
	var v724 int32
	_ = v724
	var v725 int32
	_ = v725
	var v734 int32
	_ = v734
	var v735 int32
	_ = v735
	var v738 int64
	_ = v738
	var v741 int32
	_ = v741
	var v743 int64
	_ = v743
	var v746 int32
	_ = v746
	var v753 int32
	_ = v753
	var v755 int32
	_ = v755
	var v759 int32
	_ = v759
	var v760 int32
	_ = v760
	var v761 int32
	_ = v761
	var v762 int32
	_ = v762
	var v771 int32
	_ = v771
	var v772 int32
	_ = v772
	var v773 int32
	_ = v773
	var v775 int32
	_ = v775
	var v780 int32
	_ = v780
	var v781 int32
	_ = v781
	var v788 int32
	_ = v788
	var v791 int32
	_ = v791
	var v793 int32
	_ = v793
	var v797 int32
	_ = v797
	var v798 int32
	_ = v798
	var v799 int32
	_ = v799
	var v800 int32
	_ = v800
	var v839 int32
	_ = v839
	var v851 int32
	_ = v851
	var v852 int32
	_ = v852
	var v856 int32
	_ = v856
	var v857 int32
	_ = v857
	var v865 int32
	_ = v865
	var v866 int32
	_ = v866
	var v867 int32
	_ = v867
	var v868 int32
	_ = v868
	var v870 int32
	_ = v870
	var v871 int32
	_ = v871
	var v875 int32
	_ = v875
	var v876 int32
	_ = v876
	var v877 int32
	_ = v877
	var v878 int32
	_ = v878
	var v887 int32
	_ = v887
	var v888 int32
	_ = v888
	var v891 int64
	_ = v891
	var v894 int32
	_ = v894
	var v896 int64
	_ = v896
	var v899 int32
	_ = v899
	var v906 int32
	_ = v906
	var v908 int32
	_ = v908
	var v912 int32
	_ = v912
	var v913 int32
	_ = v913
	var v914 int32
	_ = v914
	var v915 int32
	_ = v915
	var v924 int32
	_ = v924
	var v925 int32
	_ = v925
	var v926 int32
	_ = v926
	var v928 int32
	_ = v928
	var v933 int32
	_ = v933
	var v934 int32
	_ = v934
	var v941 int32
	_ = v941
	var v944 int32
	_ = v944
	var v946 int32
	_ = v946
	var v950 int32
	_ = v950
	var v951 int32
	_ = v951
	var v952 int32
	_ = v952
	var v953 int32
	_ = v953
	var v992 int32
	_ = v992
	var v995 int64
	_ = v995
	var v1015 int32
	_ = v1015
	var v1025 int32
	_ = v1025
	var v1037 int32
	_ = v1037
	var v1038 int32
	_ = v1038
	var v1042 int32
	_ = v1042
	var v1043 int32
	_ = v1043
	var v1051 int32
	_ = v1051
	var v1052 int32
	_ = v1052
	var v1053 int32
	_ = v1053
	var v1054 int32
	_ = v1054
	var v1056 int32
	_ = v1056
	var v1057 int32
	_ = v1057
	var v1061 int32
	_ = v1061
	var v1062 int32
	_ = v1062
	var v1063 int32
	_ = v1063
	var v1064 int32
	_ = v1064
	var v1073 int32
	_ = v1073
	var v1074 int32
	_ = v1074
	var v1077 int64
	_ = v1077
	var v1080 int32
	_ = v1080
	var v1082 int64
	_ = v1082
	var v1085 int32
	_ = v1085
	var v1092 int32
	_ = v1092
	var v1094 int32
	_ = v1094
	var v1098 int32
	_ = v1098
	var v1099 int32
	_ = v1099
	var v1100 int32
	_ = v1100
	var v1101 int32
	_ = v1101
	var v1110 int32
	_ = v1110
	var v1111 int32
	_ = v1111
	var v1112 int32
	_ = v1112
	var v1114 int32
	_ = v1114
	var v1119 int32
	_ = v1119
	var v1120 int32
	_ = v1120
	var v1127 int32
	_ = v1127
	var v1130 int32
	_ = v1130
	var v1132 int32
	_ = v1132
	var v1136 int32
	_ = v1136
	var v1137 int32
	_ = v1137
	var v1138 int32
	_ = v1138
	var v1139 int32
	_ = v1139
	var v1176 int32
	_ = v1176
	var v1303 int64
	_ = v1303
	var v1310 int32
	_ = v1310
	var v1323 int32
	_ = v1323
	var v1325 int32
	_ = v1325
	var v1328 int32
	_ = v1328
	var v1329 int32
	_ = v1329
	var v1330 int32
	_ = v1330
	var v1332 int32
	_ = v1332
	var v1333 int32
	_ = v1333
	var v1334 int32
	_ = v1334
	var v1341 int32
	_ = v1341
	var v1344 int32
	_ = v1344
	var v1345 int32
	_ = v1345
	var v1351 int32
	_ = v1351
	var v1353 int32
	_ = v1353
	var v1356 int32
	_ = v1356
	var v1357 int32
	_ = v1357
	var v1365 int32
	_ = v1365
	var v1368 int32
	_ = v1368
	var v1369 int32
	_ = v1369
	var v1376 int32
	_ = v1376
	var v1387 int32
	_ = v1387
	var v1399 int32
	_ = v1399
	var v1401 int32
	_ = v1401
	var v1407 int32
	_ = v1407
	var v1408 int32
	_ = v1408
	var v1410 int32
	_ = v1410
	var v1411 int32
	_ = v1411
	var v1418 int32
	_ = v1418
	var v1438 int32
	_ = v1438
	var v1439 int32
	_ = v1439
	var v1443 int32
	_ = v1443
	var v1447 int32
	_ = v1447
	var v1451 int32
	_ = v1451
	var v1455 int32
	_ = v1455
	var v1459 int32
	_ = v1459
	var v1463 int32
	_ = v1463
	var v1467 int32
	_ = v1467
	var v1471 int32
	_ = v1471
	var v1475 int32
	_ = v1475
	var v1479 int32
	_ = v1479
	var v1483 int32
	_ = v1483
	var v1487 int32
	_ = v1487
	var v1493 int32
	_ = v1493
	var v1494 int32
	_ = v1494
	var v1495 int32
	_ = v1495
	var v1497 int32
	_ = v1497
	var v1498 int32
	_ = v1498
	var v1510 int32
	_ = v1510
	var v1511 int32
	_ = v1511
	var v1515 int32
	_ = v1515
	var v1516 int32
	_ = v1516
	var v1524 int32
	_ = v1524
	var v1525 int32
	_ = v1525
	var v1526 int32
	_ = v1526
	var v1527 int32
	_ = v1527
	var v1529 int32
	_ = v1529
	var v1530 int32
	_ = v1530
	var v1534 int32
	_ = v1534
	var v1535 int32
	_ = v1535
	var v1536 int32
	_ = v1536
	var v1537 int32
	_ = v1537
	var v1546 int32
	_ = v1546
	var v1547 int32
	_ = v1547
	var v1550 int64
	_ = v1550
	var v1553 int32
	_ = v1553
	var v1555 int64
	_ = v1555
	var v1558 int32
	_ = v1558
	var v1565 int32
	_ = v1565
	var v1567 int32
	_ = v1567
	var v1571 int32
	_ = v1571
	var v1572 int32
	_ = v1572
	var v1573 int32
	_ = v1573
	var v1574 int32
	_ = v1574
	var v1583 int32
	_ = v1583
	var v1584 int32
	_ = v1584
	var v1585 int32
	_ = v1585
	var v1587 int32
	_ = v1587
	var v1592 int32
	_ = v1592
	var v1593 int32
	_ = v1593
	var v1600 int32
	_ = v1600
	var v1603 int32
	_ = v1603
	var v1605 int32
	_ = v1605
	var v1609 int32
	_ = v1609
	var v1610 int32
	_ = v1610
	var v1611 int32
	_ = v1611
	var v1612 int32
	_ = v1612
	var v1657 int32
	_ = v1657
	var v1663 int32
	_ = v1663
	var v1667 int32
	_ = v1667
	var v1669 int32
	_ = v1669
	var v1670 int32
	_ = v1670
	var v1682 int32
	_ = v1682
	var v1683 int32
	_ = v1683
	var v1687 int32
	_ = v1687
	var v1688 int32
	_ = v1688
	var v1696 int32
	_ = v1696
	var v1697 int32
	_ = v1697
	var v1698 int32
	_ = v1698
	var v1699 int32
	_ = v1699
	var v1701 int32
	_ = v1701
	var v1702 int32
	_ = v1702
	var v1706 int32
	_ = v1706
	var v1707 int32
	_ = v1707
	var v1708 int32
	_ = v1708
	var v1709 int32
	_ = v1709
	var v1718 int32
	_ = v1718
	var v1719 int32
	_ = v1719
	var v1722 int64
	_ = v1722
	var v1725 int32
	_ = v1725
	var v1727 int64
	_ = v1727
	var v1730 int32
	_ = v1730
	var v1737 int32
	_ = v1737
	var v1739 int32
	_ = v1739
	var v1743 int32
	_ = v1743
	var v1744 int32
	_ = v1744
	var v1745 int32
	_ = v1745
	var v1746 int32
	_ = v1746
	var v1755 int32
	_ = v1755
	var v1756 int32
	_ = v1756
	var v1757 int32
	_ = v1757
	var v1759 int32
	_ = v1759
	var v1764 int32
	_ = v1764
	var v1765 int32
	_ = v1765
	var v1772 int32
	_ = v1772
	var v1775 int32
	_ = v1775
	var v1777 int32
	_ = v1777
	var v1781 int32
	_ = v1781
	var v1782 int32
	_ = v1782
	var v1783 int32
	_ = v1783
	var v1784 int32
	_ = v1784
	var v1822 int32
	_ = v1822
	var v1824 int32
	_ = v1824
	var v1827 int32
	_ = v1827
	var v1834 int32
	_ = v1834
	var v1835 int32
	_ = v1835
	var v1836 int32
	_ = v1836
	var v1840 int32
	_ = v1840
	var v1846 int32
	_ = v1846
	var v1851 int32
	_ = v1851
	var v1852 int32
	_ = v1852
	var v1859 int32
	_ = v1859
	var v1861 int32
	_ = v1861
	var v1865 int32
	_ = v1865
	var v1867 int32
	_ = v1867
	var v1869 int32
	_ = v1869
	var v1876 int32
	_ = v1876
	var v1877 int32
	_ = v1877
	var v1904 int32
	_ = v1904
	var v1905 int32
	_ = v1905
	var v1908 int32
	_ = v1908
	var v1910 int32
	_ = v1910
	var v1922 int32
	_ = v1922
	var v1923 int32
	_ = v1923
	var v1930 int32
	_ = v1930
	var v1940 int32
	_ = v1940
	var v1942 int32
	_ = v1942
	var v1952 int32
	_ = v1952
	var v1954 int32
	_ = v1954
	var v1956 int32
	_ = v1956
	var v1958 int32
	_ = v1958
	var v1971 int32
	_ = v1971
	var v1972 int32
	_ = v1972
	var v1976 int32
	_ = v1976
	var v1977 int32
	_ = v1977
	var v1985 int32
	_ = v1985
	var v1986 int32
	_ = v1986
	var v1987 int32
	_ = v1987
	var v1988 int32
	_ = v1988
	var v1990 int32
	_ = v1990
	var v1991 int32
	_ = v1991
	var v1995 int32
	_ = v1995
	var v1996 int32
	_ = v1996
	var v1997 int32
	_ = v1997
	var v1998 int32
	_ = v1998
	var v2007 int32
	_ = v2007
	var v2008 int32
	_ = v2008
	var v2011 int64
	_ = v2011
	var v2014 int32
	_ = v2014
	var v2016 int64
	_ = v2016
	var v2019 int32
	_ = v2019
	var v2026 int32
	_ = v2026
	var v2028 int32
	_ = v2028
	var v2032 int32
	_ = v2032
	var v2033 int32
	_ = v2033
	var v2034 int32
	_ = v2034
	var v2035 int32
	_ = v2035
	var v2044 int32
	_ = v2044
	var v2045 int32
	_ = v2045
	var v2046 int32
	_ = v2046
	var v2048 int32
	_ = v2048
	var v2053 int32
	_ = v2053
	var v2054 int32
	_ = v2054
	var v2061 int32
	_ = v2061
	var v2064 int32
	_ = v2064
	var v2066 int32
	_ = v2066
	var v2070 int32
	_ = v2070
	var v2071 int32
	_ = v2071
	var v2072 int32
	_ = v2072
	var v2073 int32
	_ = v2073
	var v2112 int32
	_ = v2112
	var v2115 int32
	_ = v2115
	var v2116 int32
	_ = v2116
	var v2128 int32
	_ = v2128
	var v2129 int32
	_ = v2129
	var v2133 int32
	_ = v2133
	var v2134 int32
	_ = v2134
	var v2142 int32
	_ = v2142
	var v2143 int32
	_ = v2143
	var v2144 int32
	_ = v2144
	var v2145 int32
	_ = v2145
	var v2147 int32
	_ = v2147
	var v2148 int32
	_ = v2148
	var v2152 int32
	_ = v2152
	var v2153 int32
	_ = v2153
	var v2154 int32
	_ = v2154
	var v2155 int32
	_ = v2155
	var v2164 int32
	_ = v2164
	var v2165 int32
	_ = v2165
	var v2168 int64
	_ = v2168
	var v2171 int32
	_ = v2171
	var v2173 int64
	_ = v2173
	var v2176 int32
	_ = v2176
	var v2183 int32
	_ = v2183
	var v2185 int32
	_ = v2185
	var v2189 int32
	_ = v2189
	var v2190 int32
	_ = v2190
	var v2191 int32
	_ = v2191
	var v2192 int32
	_ = v2192
	var v2201 int32
	_ = v2201
	var v2202 int32
	_ = v2202
	var v2203 int32
	_ = v2203
	var v2205 int32
	_ = v2205
	var v2210 int32
	_ = v2210
	var v2211 int32
	_ = v2211
	var v2218 int32
	_ = v2218
	var v2221 int32
	_ = v2221
	var v2223 int32
	_ = v2223
	var v2227 int32
	_ = v2227
	var v2228 int32
	_ = v2228
	var v2229 int32
	_ = v2229
	var v2230 int32
	_ = v2230
	var v2268 int32
	_ = v2268
	var v2271 int32
	_ = v2271
	var v2273 int32
	_ = v2273
	var v2274 int32
	_ = v2274
	var v2286 int32
	_ = v2286
	var v2287 int32
	_ = v2287
	var v2291 int32
	_ = v2291
	var v2292 int32
	_ = v2292
	var v2300 int32
	_ = v2300
	var v2301 int32
	_ = v2301
	var v2302 int32
	_ = v2302
	var v2303 int32
	_ = v2303
	var v2305 int32
	_ = v2305
	var v2306 int32
	_ = v2306
	var v2310 int32
	_ = v2310
	var v2311 int32
	_ = v2311
	var v2312 int32
	_ = v2312
	var v2313 int32
	_ = v2313
	var v2322 int32
	_ = v2322
	var v2323 int32
	_ = v2323
	var v2326 int64
	_ = v2326
	var v2329 int32
	_ = v2329
	var v2331 int64
	_ = v2331
	var v2334 int32
	_ = v2334
	var v2341 int32
	_ = v2341
	var v2343 int32
	_ = v2343
	var v2347 int32
	_ = v2347
	var v2348 int32
	_ = v2348
	var v2349 int32
	_ = v2349
	var v2350 int32
	_ = v2350
	var v2359 int32
	_ = v2359
	var v2360 int32
	_ = v2360
	var v2361 int32
	_ = v2361
	var v2363 int32
	_ = v2363
	var v2368 int32
	_ = v2368
	var v2369 int32
	_ = v2369
	var v2376 int32
	_ = v2376
	var v2379 int32
	_ = v2379
	var v2381 int32
	_ = v2381
	var v2385 int32
	_ = v2385
	var v2386 int32
	_ = v2386
	var v2387 int32
	_ = v2387
	var v2388 int32
	_ = v2388
	var v2428 int32
	_ = v2428
	var v2440 int32
	_ = v2440
	var v2441 int32
	_ = v2441
	var v2445 int32
	_ = v2445
	var v2446 int32
	_ = v2446
	var v2454 int32
	_ = v2454
	var v2455 int32
	_ = v2455
	var v2456 int32
	_ = v2456
	var v2457 int32
	_ = v2457
	var v2459 int32
	_ = v2459
	var v2460 int32
	_ = v2460
	var v2464 int32
	_ = v2464
	var v2465 int32
	_ = v2465
	var v2466 int32
	_ = v2466
	var v2467 int32
	_ = v2467
	var v2476 int32
	_ = v2476
	var v2477 int32
	_ = v2477
	var v2480 int64
	_ = v2480
	var v2483 int32
	_ = v2483
	var v2485 int64
	_ = v2485
	var v2488 int32
	_ = v2488
	var v2495 int32
	_ = v2495
	var v2497 int32
	_ = v2497
	var v2501 int32
	_ = v2501
	var v2502 int32
	_ = v2502
	var v2503 int32
	_ = v2503
	var v2504 int32
	_ = v2504
	var v2513 int32
	_ = v2513
	var v2514 int32
	_ = v2514
	var v2515 int32
	_ = v2515
	var v2517 int32
	_ = v2517
	var v2522 int32
	_ = v2522
	var v2523 int32
	_ = v2523
	var v2530 int32
	_ = v2530
	var v2533 int32
	_ = v2533
	var v2535 int32
	_ = v2535
	var v2539 int32
	_ = v2539
	var v2540 int32
	_ = v2540
	var v2541 int32
	_ = v2541
	var v2542 int32
	_ = v2542
	var v2582 int32
	_ = v2582
	var v2585 int32
	_ = v2585
	var v2591 int32
	_ = v2591
	var v2594 int32
	_ = v2594
	var v2596 int32
	_ = v2596
	var v2597 int32
	_ = v2597
	var v2598 int32
	_ = v2598
	var v2602 int32
	_ = v2602
	var v2603 int32
	_ = v2603
	var v2605 int32
	_ = v2605
	var v2617 int32
	_ = v2617
	var v2618 int32
	_ = v2618
	var v2622 int32
	_ = v2622
	var v2623 int32
	_ = v2623
	var v2631 int32
	_ = v2631
	var v2632 int32
	_ = v2632
	var v2633 int32
	_ = v2633
	var v2634 int32
	_ = v2634
	var v2636 int32
	_ = v2636
	var v2637 int32
	_ = v2637
	var v2641 int32
	_ = v2641
	var v2642 int32
	_ = v2642
	var v2643 int32
	_ = v2643
	var v2644 int32
	_ = v2644
	var v2653 int32
	_ = v2653
	var v2654 int32
	_ = v2654
	var v2657 int64
	_ = v2657
	var v2660 int32
	_ = v2660
	var v2662 int64
	_ = v2662
	var v2665 int32
	_ = v2665
	var v2672 int32
	_ = v2672
	var v2674 int32
	_ = v2674
	var v2678 int32
	_ = v2678
	var v2679 int32
	_ = v2679
	var v2680 int32
	_ = v2680
	var v2681 int32
	_ = v2681
	var v2690 int32
	_ = v2690
	var v2691 int32
	_ = v2691
	var v2692 int32
	_ = v2692
	var v2694 int32
	_ = v2694
	var v2699 int32
	_ = v2699
	var v2700 int32
	_ = v2700
	var v2707 int32
	_ = v2707
	var v2710 int32
	_ = v2710
	var v2712 int32
	_ = v2712
	var v2716 int32
	_ = v2716
	var v2717 int32
	_ = v2717
	var v2718 int32
	_ = v2718
	var v2719 int32
	_ = v2719
	var v2759 int32
	_ = v2759
	var v2762 int32
	_ = v2762
	var v2768 int32
	_ = v2768
	var v2780 int32
	_ = v2780
	var v2781 int32
	_ = v2781
	var v2785 int32
	_ = v2785
	var v2786 int32
	_ = v2786
	var v2794 int32
	_ = v2794
	var v2795 int32
	_ = v2795
	var v2796 int32
	_ = v2796
	var v2797 int32
	_ = v2797
	var v2799 int32
	_ = v2799
	var v2800 int32
	_ = v2800
	var v2804 int32
	_ = v2804
	var v2805 int32
	_ = v2805
	var v2806 int32
	_ = v2806
	var v2807 int32
	_ = v2807
	var v2816 int32
	_ = v2816
	var v2817 int32
	_ = v2817
	var v2820 int64
	_ = v2820
	var v2823 int32
	_ = v2823
	var v2825 int64
	_ = v2825
	var v2828 int32
	_ = v2828
	var v2835 int32
	_ = v2835
	var v2837 int32
	_ = v2837
	var v2841 int32
	_ = v2841
	var v2842 int32
	_ = v2842
	var v2843 int32
	_ = v2843
	var v2844 int32
	_ = v2844
	var v2853 int32
	_ = v2853
	var v2854 int32
	_ = v2854
	var v2855 int32
	_ = v2855
	var v2857 int32
	_ = v2857
	var v2862 int32
	_ = v2862
	var v2863 int32
	_ = v2863
	var v2870 int32
	_ = v2870
	var v2873 int32
	_ = v2873
	var v2875 int32
	_ = v2875
	var v2879 int32
	_ = v2879
	var v2880 int32
	_ = v2880
	var v2881 int32
	_ = v2881
	var v2882 int32
	_ = v2882
	var v2922 int32
	_ = v2922
	v12 = m.G0
	v14 = v12 - int32(224)
	m.G0 = v14
	*(*int64)(unsafe.Add(mBase, uint32(v14)+8)) = int64(0)
	v18 = *(*int32)(unsafe.Add(mBase, uint32(l3)))
	if v18 < int32(1) {
		goto L3
	} else {
		goto L4
	}
L1:
	;
	m.G0 = v14 + int32(224)
	return
L2:
	;
	if int32(2) < v49 {
		goto L41
	} else {
		goto L42
	}
L3:
	;
	v66 = int32(1)
	v67 = int32(4)
	goto L17
L4:
	;
	v21 = *(*int32)(unsafe.Add(mBase, uint32(l3)+4))
	v22 = int32(0)
	v31 = v22
	v32 = v22
	goto L5
L5:
	;
	v36 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+v31))))
	if v36 == int32(0) {
		v49 = v32
		goto L7
	} else {
		goto L8
	}
L6:
	;
	if v49 != 0 {
		goto L2
	} else {
		goto L14
	}
L7:
	;
	v51 = v31 + int32(1)
	if v18 <= v51 {
		goto L11
	} else {
		goto L12
	}
L8:
	;
	if int32(1) < v32 {
		goto L9
	} else {
		goto L10
	}
L9:
	;
	v49 = v32 + int32(1)
	goto L7
L10:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v14+int32(8)+v32<<(uint(int32(2))%32)))) = v31
	goto L9
L11:
	;
	goto L6
L12:
	;
	if v49 < int32(3) {
		v31 = v51
		v32 = v49
		goto L5
	} else {
		goto L13
	}
L13:
	;
	goto L11
L14:
	;
	goto L3
L15:
	;
	goto L1
L16:
	;
	goto L15
L17:
	;
	v79 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v80 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v80+v67 < int32(32) {
		goto L20
	} else {
		goto L21
	}
L18:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v181 + v179
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v178<<(uint(v181)%32) | v180
	goto L16
L19:
	;
	v98 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v99 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v103 = v95
	v104 = v96
	v105 = v99
	v106 = v98
	goto L23
L20:
	;
	if v80 < int32(16) {
		v178 = v66
		v179 = v67
		v180 = v79
		v181 = v80
		goto L18
	} else {
		goto L22
	}
L21:
	;
	v84 = int32(32)
	v85 = v84 - v80
	v93 = int32(base.Ui32(v66) >> (uint(v85) % 32))
	v94 = v67 - v85
	v95 = v66<<(uint(v80)%32) | v79
	v96 = v84
	goto L19
L22:
	;
	v93 = v66
	v94 = v67
	v95 = v79
	v96 = v80
	goto L19
L23:
	;
	if base.Ui32(v105+int32(2)) <= base.Ui32(v106) {
		v161 = v105
		v162 = v106
		goto L25
	} else {
		goto L26
	}
L24:
	;
	v178 = v93
	v179 = v94
	v180 = v174
	v181 = v172
	goto L18
L25:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v161))) = uint16(v103)
	v169 = v161 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v169
	v172 = v104 + int32(-16)
	v174 = int32(base.Ui32(v103) >> (uint(int32(16)) % 32))
	if int32(31) < v104 {
		v103 = v174
		v104 = v172
		v105 = v169
		v106 = v162
		goto L23
	} else {
		goto L40
	}
L26:
	;
	v115 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v116 = v106 - v115
	v119 = base.I64_extend_i32_s(v116) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v119) {
		v143 = v115
		goto L28
	} else {
		goto L29
	}
L27:
	;
	if v105 == v115 {
		goto L38
	} else {
		goto L39
	}
L28:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v143
	goto L15
L29:
	;
	v122 = v105 - v115
	v124 = v119 + base.I64_extend_i32_u(v122)
	if base.Ui64(int64(4294967295)) < base.Ui64(v124) {
		v143 = v115
		goto L28
	} else {
		goto L30
	}
L30:
	;
	v127 = base.I32_wrap_i64(v124)
	if v106 == v115 {
		goto L31
	} else {
		goto L32
	}
L31:
	;
	v134 = int32(base.Ui32(v116*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v127) < base.Ui32(v134) {
		goto L34
	} else {
		goto L35
	}
L32:
	;
	if base.Ui32(v127) <= base.Ui32(v116) {
		v161 = v105
		v162 = v106
		goto L25
	} else {
		goto L33
	}
L33:
	;
	goto L31
L34:
	;
	v136 = v134
	goto L36
L35:
	;
	v136 = v127
	goto L36
L36:
	;
	v140 = v136&int32(-1024) + int32(1024)
	v141 = F_WebPSafeMalloc(m, int64(1), v140)
	mBase = m.M
	if v141 != 0 {
		goto L27
	} else {
		goto L37
	}
L37:
	;
	v142 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v143 = v142
	goto L28
L38:
	;
	v154 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v154)
	mBase = m.M
	v156 = v141 + v140
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v156
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v141
	v161 = v141 + v122
	v162 = v156
	goto L25
L39:
	;
	v152 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v153 = F_memcpy(m, v141, v152, v122)
	mBase = m.M
	goto L38
L40:
	;
	goto L24
L41:
	;
	v992 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v14)+207)) = v992
	v995 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v14)+200)) = v995
	*(*int64)(unsafe.Add(mBase, uint32(v14)+192)) = v995
	*(*int64)(unsafe.Add(mBase, uint32(v14+int32(174)))) = v995
	*(*int64)(unsafe.Add(mBase, uint32(v14+int32(168)))) = v995
	*(*int64)(unsafe.Add(mBase, uint32(v14+int32(160)))) = v995
	*(*int64)(unsafe.Add(mBase, uint32(v14)+152)) = v995
	*(*int64)(unsafe.Add(mBase, uint32(v14)+144)) = v995
	v1015 = int32(19)
	*(*int32)(unsafe.Add(mBase, uint32(v14)+132)) = v1015
	*(*int32)(unsafe.Add(mBase, uint32(v14)+140)) = v14 + int32(144)
	*(*int32)(unsafe.Add(mBase, uint32(v14)+136)) = v14 + int32(192)
	v1025 = int32(1)
	goto L181
L42:
	;
	v220 = *(*int32)(unsafe.Add(mBase, uint32(v14)+8))
	if int32(255) < v220 {
		goto L41
	} else {
		goto L43
	}
L43:
	;
	v223 = *(*int32)(unsafe.Add(mBase, uint32(v14)+12))
	if int32(255) < v223 {
		goto L41
	} else {
		goto L44
	}
L44:
	;
	v226 = int32(1)
	goto L47
L45:
	;
	v379 = v49 + int32(-1)
	v380 = int32(1)
	goto L73
L46:
	;
	goto L45
L47:
	;
	v239 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v240 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v240+v226 < int32(32) {
		goto L50
	} else {
		goto L51
	}
L48:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v341 + v339
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v338<<(uint(v341)%32) | v340
	goto L46
L49:
	;
	v258 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v259 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v263 = v255
	v264 = v256
	v265 = v259
	v266 = v258
	goto L53
L50:
	;
	if v240 < int32(16) {
		v338 = v226
		v339 = v226
		v340 = v239
		v341 = v240
		goto L48
	} else {
		goto L52
	}
L51:
	;
	v244 = int32(32)
	v245 = v244 - v240
	v253 = int32(base.Ui32(v226) >> (uint(v245) % 32))
	v254 = v226 - v245
	v255 = v226<<(uint(v240)%32) | v239
	v256 = v244
	goto L49
L52:
	;
	v253 = v226
	v254 = v226
	v255 = v239
	v256 = v240
	goto L49
L53:
	;
	if base.Ui32(v265+int32(2)) <= base.Ui32(v266) {
		v321 = v265
		v322 = v266
		goto L55
	} else {
		goto L56
	}
L54:
	;
	v338 = v253
	v339 = v254
	v340 = v334
	v341 = v332
	goto L48
L55:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v321))) = uint16(v263)
	v329 = v321 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v329
	v332 = v264 + int32(-16)
	v334 = int32(base.Ui32(v263) >> (uint(int32(16)) % 32))
	if int32(31) < v264 {
		v263 = v334
		v264 = v332
		v265 = v329
		v266 = v322
		goto L53
	} else {
		goto L70
	}
L56:
	;
	v275 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v276 = v266 - v275
	v279 = base.I64_extend_i32_s(v276) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v279) {
		v303 = v275
		goto L58
	} else {
		goto L59
	}
L57:
	;
	if v265 == v275 {
		goto L68
	} else {
		goto L69
	}
L58:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v303
	goto L45
L59:
	;
	v282 = v265 - v275
	v284 = v279 + base.I64_extend_i32_u(v282)
	if base.Ui64(int64(4294967295)) < base.Ui64(v284) {
		v303 = v275
		goto L58
	} else {
		goto L60
	}
L60:
	;
	v287 = base.I32_wrap_i64(v284)
	if v266 == v275 {
		goto L61
	} else {
		goto L62
	}
L61:
	;
	v294 = int32(base.Ui32(v276*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v287) < base.Ui32(v294) {
		goto L64
	} else {
		goto L65
	}
L62:
	;
	if base.Ui32(v287) <= base.Ui32(v276) {
		v321 = v265
		v322 = v266
		goto L55
	} else {
		goto L63
	}
L63:
	;
	goto L61
L64:
	;
	v296 = v294
	goto L66
L65:
	;
	v296 = v287
	goto L66
L66:
	;
	v300 = v296&int32(-1024) + int32(1024)
	v301 = F_WebPSafeMalloc(m, int64(1), v300)
	mBase = m.M
	if v301 != 0 {
		goto L57
	} else {
		goto L67
	}
L67:
	;
	v302 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v303 = v302
	goto L58
L68:
	;
	v314 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v314)
	mBase = m.M
	v316 = v301 + v300
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v316
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v301
	v321 = v301 + v282
	v322 = v316
	goto L55
L69:
	;
	v312 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v313 = F_memcpy(m, v301, v312, v282)
	mBase = m.M
	goto L68
L70:
	;
	goto L54
L71:
	;
	v531 = int32(1)
	v532 = base.B2i32(v531 < v220)
	goto L99
L72:
	;
	goto L71
L73:
	;
	v392 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v393 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v393+v380 < int32(32) {
		goto L76
	} else {
		goto L77
	}
L74:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v494 + v492
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v491<<(uint(v494)%32) | v493
	goto L72
L75:
	;
	v411 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v412 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v416 = v408
	v417 = v409
	v418 = v412
	v419 = v411
	goto L79
L76:
	;
	if v393 < int32(16) {
		v491 = v379
		v492 = v380
		v493 = v392
		v494 = v393
		goto L74
	} else {
		goto L78
	}
L77:
	;
	v397 = int32(32)
	v398 = v397 - v393
	v406 = int32(base.Ui32(v379) >> (uint(v398) % 32))
	v407 = v380 - v398
	v408 = v379<<(uint(v393)%32) | v392
	v409 = v397
	goto L75
L78:
	;
	v406 = v379
	v407 = v380
	v408 = v392
	v409 = v393
	goto L75
L79:
	;
	if base.Ui32(v418+int32(2)) <= base.Ui32(v419) {
		v474 = v418
		v475 = v419
		goto L81
	} else {
		goto L82
	}
L80:
	;
	v491 = v406
	v492 = v407
	v493 = v487
	v494 = v485
	goto L74
L81:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v474))) = uint16(v416)
	v482 = v474 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v482
	v485 = v417 + int32(-16)
	v487 = int32(base.Ui32(v416) >> (uint(int32(16)) % 32))
	if int32(31) < v417 {
		v416 = v487
		v417 = v485
		v418 = v482
		v419 = v475
		goto L79
	} else {
		goto L96
	}
L82:
	;
	v428 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v429 = v419 - v428
	v432 = base.I64_extend_i32_s(v429) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v432) {
		v456 = v428
		goto L84
	} else {
		goto L85
	}
L83:
	;
	if v418 == v428 {
		goto L94
	} else {
		goto L95
	}
L84:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v456
	goto L71
L85:
	;
	v435 = v418 - v428
	v437 = v432 + base.I64_extend_i32_u(v435)
	if base.Ui64(int64(4294967295)) < base.Ui64(v437) {
		v456 = v428
		goto L84
	} else {
		goto L86
	}
L86:
	;
	v440 = base.I32_wrap_i64(v437)
	if v419 == v428 {
		goto L87
	} else {
		goto L88
	}
L87:
	;
	v447 = int32(base.Ui32(v429*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v440) < base.Ui32(v447) {
		goto L90
	} else {
		goto L91
	}
L88:
	;
	if base.Ui32(v440) <= base.Ui32(v429) {
		v474 = v418
		v475 = v419
		goto L81
	} else {
		goto L89
	}
L89:
	;
	goto L87
L90:
	;
	v449 = v447
	goto L92
L91:
	;
	v449 = v440
	goto L92
L92:
	;
	v453 = v449&int32(-1024) + int32(1024)
	v454 = F_WebPSafeMalloc(m, int64(1), v453)
	mBase = m.M
	if v454 != 0 {
		goto L83
	} else {
		goto L93
	}
L93:
	;
	v455 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v456 = v455
	goto L84
L94:
	;
	v467 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v467)
	mBase = m.M
	v469 = v454 + v453
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v469
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v454
	v474 = v454 + v435
	v475 = v469
	goto L81
L95:
	;
	v465 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v466 = F_memcpy(m, v454, v465, v435)
	mBase = m.M
	goto L94
L96:
	;
	goto L80
L97:
	;
	if v531 < v220 {
		goto L123
	} else {
		goto L124
	}
L98:
	;
	goto L97
L99:
	;
	v545 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v546 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v546+v531 < int32(32) {
		goto L102
	} else {
		goto L103
	}
L100:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v647 + v645
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v644<<(uint(v647)%32) | v646
	goto L98
L101:
	;
	v564 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v565 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v569 = v561
	v570 = v562
	v571 = v565
	v572 = v564
	goto L105
L102:
	;
	if v546 < int32(16) {
		v644 = v532
		v645 = v531
		v646 = v545
		v647 = v546
		goto L100
	} else {
		goto L104
	}
L103:
	;
	v550 = int32(32)
	v551 = v550 - v546
	v559 = int32(base.Ui32(v532) >> (uint(v551) % 32))
	v560 = v531 - v551
	v561 = v532<<(uint(v546)%32) | v545
	v562 = v550
	goto L101
L104:
	;
	v559 = v532
	v560 = v531
	v561 = v545
	v562 = v546
	goto L101
L105:
	;
	if base.Ui32(v571+int32(2)) <= base.Ui32(v572) {
		v627 = v571
		v628 = v572
		goto L107
	} else {
		goto L108
	}
L106:
	;
	v644 = v559
	v645 = v560
	v646 = v640
	v647 = v638
	goto L100
L107:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v627))) = uint16(v569)
	v635 = v627 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v635
	v638 = v570 + int32(-16)
	v640 = int32(base.Ui32(v569) >> (uint(int32(16)) % 32))
	if int32(31) < v570 {
		v569 = v640
		v570 = v638
		v571 = v635
		v572 = v628
		goto L105
	} else {
		goto L122
	}
L108:
	;
	v581 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v582 = v572 - v581
	v585 = base.I64_extend_i32_s(v582) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v585) {
		v609 = v581
		goto L110
	} else {
		goto L111
	}
L109:
	;
	if v571 == v581 {
		goto L120
	} else {
		goto L121
	}
L110:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v609
	goto L97
L111:
	;
	v588 = v571 - v581
	v590 = v585 + base.I64_extend_i32_u(v588)
	if base.Ui64(int64(4294967295)) < base.Ui64(v590) {
		v609 = v581
		goto L110
	} else {
		goto L112
	}
L112:
	;
	v593 = base.I32_wrap_i64(v590)
	if v572 == v581 {
		goto L113
	} else {
		goto L114
	}
L113:
	;
	v600 = int32(base.Ui32(v582*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v593) < base.Ui32(v600) {
		goto L116
	} else {
		goto L117
	}
L114:
	;
	if base.Ui32(v593) <= base.Ui32(v582) {
		v627 = v571
		v628 = v572
		goto L107
	} else {
		goto L115
	}
L115:
	;
	goto L113
L116:
	;
	v602 = v600
	goto L118
L117:
	;
	v602 = v593
	goto L118
L118:
	;
	v606 = v602&int32(-1024) + int32(1024)
	v607 = F_WebPSafeMalloc(m, int64(1), v606)
	mBase = m.M
	if v607 != 0 {
		goto L109
	} else {
		goto L119
	}
L119:
	;
	v608 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v609 = v608
	goto L110
L120:
	;
	v620 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v620)
	mBase = m.M
	v622 = v607 + v606
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v622
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v607
	v627 = v607 + v588
	v628 = v622
	goto L107
L121:
	;
	v618 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v619 = F_memcpy(m, v607, v618, v588)
	mBase = m.M
	goto L120
L122:
	;
	goto L106
L123:
	;
	v686 = int32(8)
	goto L125
L124:
	;
	v686 = int32(1)
	goto L125
L125:
	;
	if v686 < int32(1) {
		goto L127
	} else {
		goto L128
	}
L126:
	;
	if v49 != int32(2) {
		goto L1
	} else {
		goto L152
	}
L127:
	;
	goto L126
L128:
	;
	v698 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v699 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v699+v686 < int32(32) {
		goto L131
	} else {
		goto L132
	}
L129:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v800 + v798
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v797<<(uint(v800)%32) | v799
	goto L127
L130:
	;
	v717 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v718 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v722 = v714
	v723 = v715
	v724 = v718
	v725 = v717
	goto L134
L131:
	;
	if v699 < int32(16) {
		v797 = v220
		v798 = v686
		v799 = v698
		v800 = v699
		goto L129
	} else {
		goto L133
	}
L132:
	;
	v703 = int32(32)
	v704 = v703 - v699
	v712 = int32(base.Ui32(v220) >> (uint(v704) % 32))
	v713 = v686 - v704
	v714 = v220<<(uint(v699)%32) | v698
	v715 = v703
	goto L130
L133:
	;
	v712 = v220
	v713 = v686
	v714 = v698
	v715 = v699
	goto L130
L134:
	;
	if base.Ui32(v724+int32(2)) <= base.Ui32(v725) {
		v780 = v724
		v781 = v725
		goto L136
	} else {
		goto L137
	}
L135:
	;
	v797 = v712
	v798 = v713
	v799 = v793
	v800 = v791
	goto L129
L136:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v780))) = uint16(v722)
	v788 = v780 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v788
	v791 = v723 + int32(-16)
	v793 = int32(base.Ui32(v722) >> (uint(int32(16)) % 32))
	if int32(31) < v723 {
		v722 = v793
		v723 = v791
		v724 = v788
		v725 = v781
		goto L134
	} else {
		goto L151
	}
L137:
	;
	v734 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v735 = v725 - v734
	v738 = base.I64_extend_i32_s(v735) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v738) {
		v762 = v734
		goto L139
	} else {
		goto L140
	}
L138:
	;
	if v724 == v734 {
		goto L149
	} else {
		goto L150
	}
L139:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v762
	goto L126
L140:
	;
	v741 = v724 - v734
	v743 = v738 + base.I64_extend_i32_u(v741)
	if base.Ui64(int64(4294967295)) < base.Ui64(v743) {
		v762 = v734
		goto L139
	} else {
		goto L141
	}
L141:
	;
	v746 = base.I32_wrap_i64(v743)
	if v725 == v734 {
		goto L142
	} else {
		goto L143
	}
L142:
	;
	v753 = int32(base.Ui32(v735*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v746) < base.Ui32(v753) {
		goto L145
	} else {
		goto L146
	}
L143:
	;
	if base.Ui32(v746) <= base.Ui32(v735) {
		v780 = v724
		v781 = v725
		goto L136
	} else {
		goto L144
	}
L144:
	;
	goto L142
L145:
	;
	v755 = v753
	goto L147
L146:
	;
	v755 = v746
	goto L147
L147:
	;
	v759 = v755&int32(-1024) + int32(1024)
	v760 = F_WebPSafeMalloc(m, int64(1), v759)
	mBase = m.M
	if v760 != 0 {
		goto L138
	} else {
		goto L148
	}
L148:
	;
	v761 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v762 = v761
	goto L139
L149:
	;
	v773 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v773)
	mBase = m.M
	v775 = v760 + v759
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v775
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v760
	v780 = v760 + v741
	v781 = v775
	goto L136
L150:
	;
	v771 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v772 = F_memcpy(m, v760, v771, v741)
	mBase = m.M
	goto L149
L151:
	;
	goto L135
L152:
	;
	v839 = int32(8)
	goto L155
L153:
	;
	goto L1
L154:
	;
	goto L153
L155:
	;
	v851 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v852 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v852+v839 < int32(32) {
		goto L158
	} else {
		goto L159
	}
L156:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v953 + v951
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v950<<(uint(v953)%32) | v952
	goto L154
L157:
	;
	v870 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v871 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v875 = v867
	v876 = v868
	v877 = v871
	v878 = v870
	goto L161
L158:
	;
	if v852 < int32(16) {
		v950 = v223
		v951 = v839
		v952 = v851
		v953 = v852
		goto L156
	} else {
		goto L160
	}
L159:
	;
	v856 = int32(32)
	v857 = v856 - v852
	v865 = int32(base.Ui32(v223) >> (uint(v857) % 32))
	v866 = v839 - v857
	v867 = v223<<(uint(v852)%32) | v851
	v868 = v856
	goto L157
L160:
	;
	v865 = v223
	v866 = v839
	v867 = v851
	v868 = v852
	goto L157
L161:
	;
	if base.Ui32(v877+int32(2)) <= base.Ui32(v878) {
		v933 = v877
		v934 = v878
		goto L163
	} else {
		goto L164
	}
L162:
	;
	v950 = v865
	v951 = v866
	v952 = v946
	v953 = v944
	goto L156
L163:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v933))) = uint16(v875)
	v941 = v933 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v941
	v944 = v876 + int32(-16)
	v946 = int32(base.Ui32(v875) >> (uint(int32(16)) % 32))
	if int32(31) < v876 {
		v875 = v946
		v876 = v944
		v877 = v941
		v878 = v934
		goto L161
	} else {
		goto L178
	}
L164:
	;
	v887 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v888 = v878 - v887
	v891 = base.I64_extend_i32_s(v888) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v891) {
		v915 = v887
		goto L166
	} else {
		goto L167
	}
L165:
	;
	if v877 == v887 {
		goto L176
	} else {
		goto L177
	}
L166:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v915
	goto L153
L167:
	;
	v894 = v877 - v887
	v896 = v891 + base.I64_extend_i32_u(v894)
	if base.Ui64(int64(4294967295)) < base.Ui64(v896) {
		v915 = v887
		goto L166
	} else {
		goto L168
	}
L168:
	;
	v899 = base.I32_wrap_i64(v896)
	if v878 == v887 {
		goto L169
	} else {
		goto L170
	}
L169:
	;
	v906 = int32(base.Ui32(v888*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v899) < base.Ui32(v906) {
		goto L172
	} else {
		goto L173
	}
L170:
	;
	if base.Ui32(v899) <= base.Ui32(v888) {
		v933 = v877
		v934 = v878
		goto L163
	} else {
		goto L171
	}
L171:
	;
	goto L169
L172:
	;
	v908 = v906
	goto L174
L173:
	;
	v908 = v899
	goto L174
L174:
	;
	v912 = v908&int32(-1024) + int32(1024)
	v913 = F_WebPSafeMalloc(m, int64(1), v912)
	mBase = m.M
	if v913 != 0 {
		goto L165
	} else {
		goto L175
	}
L175:
	;
	v914 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v915 = v914
	goto L166
L176:
	;
	v926 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v926)
	mBase = m.M
	v928 = v913 + v912
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v928
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v913
	v933 = v913 + v894
	v934 = v928
	goto L163
L177:
	;
	v924 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v925 = F_memcpy(m, v913, v924, v894)
	mBase = m.M
	goto L176
L178:
	;
	goto L162
L179:
	;
	v1176 = F_VP8LCreateCompressedHuffmanTree(m, l3, l2, v18)
	mBase = m.M
	goto L207
L180:
	;
	goto L179
L181:
	;
	v1037 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v1038 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v1038+v1025 < int32(32) {
		goto L184
	} else {
		goto L185
	}
L182:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1139 + v1137
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1136<<(uint(v1139)%32) | v1138
	goto L180
L183:
	;
	v1056 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v1057 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v1061 = v1053
	v1062 = v1054
	v1063 = v1057
	v1064 = v1056
	goto L187
L184:
	;
	if v1038 < int32(16) {
		v1136 = v992
		v1137 = v1025
		v1138 = v1037
		v1139 = v1038
		goto L182
	} else {
		goto L186
	}
L185:
	;
	v1042 = int32(32)
	v1043 = v1042 - v1038
	v1051 = int32(base.Ui32(v992) >> (uint(v1043) % 32))
	v1052 = v1025 - v1043
	v1053 = v992<<(uint(v1038)%32) | v1037
	v1054 = v1042
	goto L183
L186:
	;
	v1051 = v992
	v1052 = v1025
	v1053 = v1037
	v1054 = v1038
	goto L183
L187:
	;
	if base.Ui32(v1063+int32(2)) <= base.Ui32(v1064) {
		v1119 = v1063
		v1120 = v1064
		goto L189
	} else {
		goto L190
	}
L188:
	;
	v1136 = v1051
	v1137 = v1052
	v1138 = v1132
	v1139 = v1130
	goto L182
L189:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1119))) = uint16(v1061)
	v1127 = v1119 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1127
	v1130 = v1062 + int32(-16)
	v1132 = int32(base.Ui32(v1061) >> (uint(int32(16)) % 32))
	if int32(31) < v1062 {
		v1061 = v1132
		v1062 = v1130
		v1063 = v1127
		v1064 = v1120
		goto L187
	} else {
		goto L204
	}
L190:
	;
	v1073 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1074 = v1064 - v1073
	v1077 = base.I64_extend_i32_s(v1074) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1077) {
		v1101 = v1073
		goto L192
	} else {
		goto L193
	}
L191:
	;
	if v1063 == v1073 {
		goto L202
	} else {
		goto L203
	}
L192:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1101
	goto L179
L193:
	;
	v1080 = v1063 - v1073
	v1082 = v1077 + base.I64_extend_i32_u(v1080)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1082) {
		v1101 = v1073
		goto L192
	} else {
		goto L194
	}
L194:
	;
	v1085 = base.I32_wrap_i64(v1082)
	if v1064 == v1073 {
		goto L195
	} else {
		goto L196
	}
L195:
	;
	v1092 = int32(base.Ui32(v1074*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v1085) < base.Ui32(v1092) {
		goto L198
	} else {
		goto L199
	}
L196:
	;
	if base.Ui32(v1085) <= base.Ui32(v1074) {
		v1119 = v1063
		v1120 = v1064
		goto L189
	} else {
		goto L197
	}
L197:
	;
	goto L195
L198:
	;
	v1094 = v1092
	goto L200
L199:
	;
	v1094 = v1085
	goto L200
L200:
	;
	v1098 = v1094&int32(-1024) + int32(1024)
	v1099 = F_WebPSafeMalloc(m, int64(1), v1098)
	mBase = m.M
	if v1099 != 0 {
		goto L191
	} else {
		goto L201
	}
L201:
	;
	v1100 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1101 = v1100
	goto L192
L202:
	;
	v1112 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v1112)
	mBase = m.M
	v1114 = v1099 + v1098
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v1114
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v1099
	v1119 = v1099 + v1080
	v1120 = v1114
	goto L189
L203:
	;
	v1110 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1111 = F_memcpy(m, v1099, v1110, v1080)
	mBase = m.M
	goto L202
L204:
	;
	goto L188
L205:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v14)+31)) = int32(0)
	v1303 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v14)+24)) = v1303
	*(*int64)(unsafe.Add(mBase, uint32(v14)+16)) = v1303
	if v1176 < int32(1) {
		goto L219
	} else {
		goto L220
	}
L207:
	;
	base.MemoryFill(m, v14+int32(48), int32(0), int32(76))
	goto L205
L219:
	;
	F_VP8LCreateHuffmanTree(m, v14+int32(48), int32(7), v14+int32(16), l1, v14+int32(132))
	mBase = m.M
	v1438 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+207)))
	if v1438 != 0 {
		v1495 = v1015
		goto L230
	} else {
		goto L231
	}
L220:
	;
	v1310 = v1176 & int32(3)
	if base.Ui32(v1176) < base.Ui32(int32(4)) {
		v1387 = v992
		goto L221
	} else {
		goto L222
	}
L221:
	;
	if v1310 == int32(0) {
		goto L219
	} else {
		goto L226
	}
L222:
	;
	v1323 = l2
	v1325 = int32(0)
	goto L223
L223:
	;
	v1328 = v14 + int32(48)
	v1329 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1323))))
	v1330 = int32(2)
	v1332 = v1328 + v1329<<(uint(v1330)%32)
	v1333 = *(*int32)(unsafe.Add(mBase, uint32(v1332)))
	v1334 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v1332))) = v1333 + v1334
	v1341 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1323+v1330))))
	v1344 = v1328 + v1341<<(uint(v1330)%32)
	v1345 = *(*int32)(unsafe.Add(mBase, uint32(v1344)))
	*(*int32)(unsafe.Add(mBase, uint32(v1344))) = v1345 + v1334
	v1351 = int32(4)
	v1353 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1323+v1351))))
	v1356 = v1328 + v1353<<(uint(v1330)%32)
	v1357 = *(*int32)(unsafe.Add(mBase, uint32(v1356)))
	*(*int32)(unsafe.Add(mBase, uint32(v1356))) = v1357 + v1334
	v1365 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1323+int32(6)))))
	v1368 = v1328 + v1365<<(uint(v1330)%32)
	v1369 = *(*int32)(unsafe.Add(mBase, uint32(v1368)))
	*(*int32)(unsafe.Add(mBase, uint32(v1368))) = v1369 + v1334
	v1376 = v1325 + v1351
	if v1176&int32(2147483644) != v1376 {
		v1323 = v1323 + int32(8)
		v1325 = v1376
		goto L223
	} else {
		goto L225
	}
L224:
	;
	v1387 = v1376
	goto L221
L225:
	;
	goto L224
L226:
	;
	v1399 = v1310
	v1401 = l2 + v1387<<(uint(int32(1))%32)
	goto L227
L227:
	;
	v1407 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1401))))
	v1408 = int32(2)
	v1410 = v14 + int32(48) + v1407<<(uint(v1408)%32)
	v1411 = *(*int32)(unsafe.Add(mBase, uint32(v1410)))
	*(*int32)(unsafe.Add(mBase, uint32(v1410))) = v1411 + int32(1)
	v1418 = v1399 + int32(-1)
	if v1418 != 0 {
		v1399 = v1418
		v1401 = v1401 + v1408
		goto L227
	} else {
		goto L229
	}
L228:
	;
	goto L219
L229:
	;
	goto L228
L230:
	;
	v1497 = v1495 + int32(-4)
	v1498 = int32(4)
	goto L263
L231:
	;
	v1439 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+206)))
	if v1439 == int32(0) {
		goto L232
	} else {
		goto L233
	}
L232:
	;
	v1443 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+205)))
	if v1443 == int32(0) {
		goto L234
	} else {
		goto L235
	}
L233:
	;
	v1495 = int32(18)
	goto L230
L234:
	;
	v1447 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+204)))
	if v1447 == int32(0) {
		goto L236
	} else {
		goto L237
	}
L235:
	;
	v1495 = int32(17)
	goto L230
L236:
	;
	v1451 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+203)))
	if v1451 == int32(0) {
		goto L238
	} else {
		goto L239
	}
L237:
	;
	v1495 = int32(16)
	goto L230
L238:
	;
	v1455 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+202)))
	if v1455 == int32(0) {
		goto L240
	} else {
		goto L241
	}
L239:
	;
	v1495 = int32(15)
	goto L230
L240:
	;
	v1459 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+201)))
	if v1459 == int32(0) {
		goto L242
	} else {
		goto L243
	}
L241:
	;
	v1495 = int32(14)
	goto L230
L242:
	;
	v1463 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+200)))
	if v1463 == int32(0) {
		goto L244
	} else {
		goto L245
	}
L243:
	;
	v1495 = int32(13)
	goto L230
L244:
	;
	v1467 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+199)))
	if v1467 == int32(0) {
		goto L246
	} else {
		goto L247
	}
L245:
	;
	v1495 = int32(12)
	goto L230
L246:
	;
	v1471 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+198)))
	if v1471 == int32(0) {
		goto L248
	} else {
		goto L249
	}
L247:
	;
	v1495 = int32(11)
	goto L230
L248:
	;
	v1475 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+208)))
	if v1475 == int32(0) {
		goto L250
	} else {
		goto L251
	}
L249:
	;
	v1495 = int32(10)
	goto L230
L250:
	;
	v1479 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+197)))
	if v1479 == int32(0) {
		goto L252
	} else {
		goto L253
	}
L251:
	;
	v1495 = int32(9)
	goto L230
L252:
	;
	v1483 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+196)))
	if v1483 == int32(0) {
		goto L254
	} else {
		goto L255
	}
L253:
	;
	v1495 = int32(8)
	goto L230
L254:
	;
	v1487 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+195)))
	if v1487 == int32(0) {
		goto L256
	} else {
		goto L257
	}
L255:
	;
	v1495 = int32(7)
	goto L230
L256:
	;
	v1493 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+194)))
	if v1493 != 0 {
		goto L258
	} else {
		goto L259
	}
L257:
	;
	v1495 = int32(6)
	goto L230
L258:
	;
	v1494 = int32(5)
	goto L260
L259:
	;
	v1494 = int32(4)
	goto L260
L260:
	;
	v1495 = v1494
	goto L230
L261:
	;
	v1657 = int32(0)
	goto L287
L262:
	;
	goto L261
L263:
	;
	v1510 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v1511 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v1511+v1498 < int32(32) {
		goto L266
	} else {
		goto L267
	}
L264:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1612 + v1610
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1609<<(uint(v1612)%32) | v1611
	goto L262
L265:
	;
	v1529 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v1530 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v1534 = v1526
	v1535 = v1527
	v1536 = v1530
	v1537 = v1529
	goto L269
L266:
	;
	if v1511 < int32(16) {
		v1609 = v1497
		v1610 = v1498
		v1611 = v1510
		v1612 = v1511
		goto L264
	} else {
		goto L268
	}
L267:
	;
	v1515 = int32(32)
	v1516 = v1515 - v1511
	v1524 = int32(base.Ui32(v1497) >> (uint(v1516) % 32))
	v1525 = v1498 - v1516
	v1526 = v1497<<(uint(v1511)%32) | v1510
	v1527 = v1515
	goto L265
L268:
	;
	v1524 = v1497
	v1525 = v1498
	v1526 = v1510
	v1527 = v1511
	goto L265
L269:
	;
	if base.Ui32(v1536+int32(2)) <= base.Ui32(v1537) {
		v1592 = v1536
		v1593 = v1537
		goto L271
	} else {
		goto L272
	}
L270:
	;
	v1609 = v1524
	v1610 = v1525
	v1611 = v1605
	v1612 = v1603
	goto L264
L271:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1592))) = uint16(v1534)
	v1600 = v1592 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1600
	v1603 = v1535 + int32(-16)
	v1605 = int32(base.Ui32(v1534) >> (uint(int32(16)) % 32))
	if int32(31) < v1535 {
		v1534 = v1605
		v1535 = v1603
		v1536 = v1600
		v1537 = v1593
		goto L269
	} else {
		goto L286
	}
L272:
	;
	v1546 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1547 = v1537 - v1546
	v1550 = base.I64_extend_i32_s(v1547) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1550) {
		v1574 = v1546
		goto L274
	} else {
		goto L275
	}
L273:
	;
	if v1536 == v1546 {
		goto L284
	} else {
		goto L285
	}
L274:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1574
	goto L261
L275:
	;
	v1553 = v1536 - v1546
	v1555 = v1550 + base.I64_extend_i32_u(v1553)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1555) {
		v1574 = v1546
		goto L274
	} else {
		goto L276
	}
L276:
	;
	v1558 = base.I32_wrap_i64(v1555)
	if v1537 == v1546 {
		goto L277
	} else {
		goto L278
	}
L277:
	;
	v1565 = int32(base.Ui32(v1547*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v1558) < base.Ui32(v1565) {
		goto L280
	} else {
		goto L281
	}
L278:
	;
	if base.Ui32(v1558) <= base.Ui32(v1547) {
		v1592 = v1536
		v1593 = v1537
		goto L271
	} else {
		goto L279
	}
L279:
	;
	goto L277
L280:
	;
	v1567 = v1565
	goto L282
L281:
	;
	v1567 = v1558
	goto L282
L282:
	;
	v1571 = v1567&int32(-1024) + int32(1024)
	v1572 = F_WebPSafeMalloc(m, int64(1), v1571)
	mBase = m.M
	if v1572 != 0 {
		goto L273
	} else {
		goto L283
	}
L283:
	;
	v1573 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1574 = v1573
	goto L274
L284:
	;
	v1585 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v1585)
	mBase = m.M
	v1587 = v1572 + v1571
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v1587
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v1572
	v1592 = v1572 + v1553
	v1593 = v1587
	goto L271
L285:
	;
	v1583 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1584 = F_memcpy(m, v1572, v1583, v1553)
	mBase = m.M
	goto L284
L286:
	;
	goto L270
L287:
	;
	v1663 = m.G1
	v1667 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1663+int32(_a_F_StoreHuffmanCode_0)+v1657))))
	v1669 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14+int32(192)+v1667))))
	v1670 = int32(3)
	goto L291
L288:
	;
	v1824 = *(*int32)(unsafe.Add(mBase, uint32(v14)+132))
	if v1824 < int32(1) {
		goto L316
	} else {
		goto L317
	}
L289:
	;
	v1822 = v1657 + int32(1)
	if v1495 != v1822 {
		v1657 = v1822
		goto L287
	} else {
		goto L315
	}
L290:
	;
	goto L289
L291:
	;
	v1682 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v1683 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v1683+v1670 < int32(32) {
		goto L294
	} else {
		goto L295
	}
L292:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1784 + v1782
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1781<<(uint(v1784)%32) | v1783
	goto L290
L293:
	;
	v1701 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v1702 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v1706 = v1698
	v1707 = v1699
	v1708 = v1702
	v1709 = v1701
	goto L297
L294:
	;
	if v1683 < int32(16) {
		v1781 = v1669
		v1782 = v1670
		v1783 = v1682
		v1784 = v1683
		goto L292
	} else {
		goto L296
	}
L295:
	;
	v1687 = int32(32)
	v1688 = v1687 - v1683
	v1696 = int32(base.Ui32(v1669) >> (uint(v1688) % 32))
	v1697 = v1670 - v1688
	v1698 = v1669<<(uint(v1683)%32) | v1682
	v1699 = v1687
	goto L293
L296:
	;
	v1696 = v1669
	v1697 = v1670
	v1698 = v1682
	v1699 = v1683
	goto L293
L297:
	;
	if base.Ui32(v1708+int32(2)) <= base.Ui32(v1709) {
		v1764 = v1708
		v1765 = v1709
		goto L299
	} else {
		goto L300
	}
L298:
	;
	v1781 = v1696
	v1782 = v1697
	v1783 = v1777
	v1784 = v1775
	goto L292
L299:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1764))) = uint16(v1706)
	v1772 = v1764 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1772
	v1775 = v1707 + int32(-16)
	v1777 = int32(base.Ui32(v1706) >> (uint(int32(16)) % 32))
	if int32(31) < v1707 {
		v1706 = v1777
		v1707 = v1775
		v1708 = v1772
		v1709 = v1765
		goto L297
	} else {
		goto L314
	}
L300:
	;
	v1718 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1719 = v1709 - v1718
	v1722 = base.I64_extend_i32_s(v1719) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1722) {
		v1746 = v1718
		goto L302
	} else {
		goto L303
	}
L301:
	;
	if v1708 == v1718 {
		goto L312
	} else {
		goto L313
	}
L302:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1746
	goto L289
L303:
	;
	v1725 = v1708 - v1718
	v1727 = v1722 + base.I64_extend_i32_u(v1725)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1727) {
		v1746 = v1718
		goto L302
	} else {
		goto L304
	}
L304:
	;
	v1730 = base.I32_wrap_i64(v1727)
	if v1709 == v1718 {
		goto L305
	} else {
		goto L306
	}
L305:
	;
	v1737 = int32(base.Ui32(v1719*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v1730) < base.Ui32(v1737) {
		goto L308
	} else {
		goto L309
	}
L306:
	;
	if base.Ui32(v1730) <= base.Ui32(v1719) {
		v1764 = v1708
		v1765 = v1709
		goto L299
	} else {
		goto L307
	}
L307:
	;
	goto L305
L308:
	;
	v1739 = v1737
	goto L310
L309:
	;
	v1739 = v1730
	goto L310
L310:
	;
	v1743 = v1739&int32(-1024) + int32(1024)
	v1744 = F_WebPSafeMalloc(m, int64(1), v1743)
	mBase = m.M
	if v1744 != 0 {
		goto L301
	} else {
		goto L311
	}
L311:
	;
	v1745 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1746 = v1745
	goto L302
L312:
	;
	v1757 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v1757)
	mBase = m.M
	v1759 = v1744 + v1743
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v1759
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v1744
	v1764 = v1744 + v1725
	v1765 = v1759
	goto L299
L313:
	;
	v1755 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1756 = F_memcpy(m, v1744, v1755, v1725)
	mBase = m.M
	goto L312
L314:
	;
	goto L298
L315:
	;
	goto L288
L316:
	;
	if int32(1) <= v1176 {
		goto L328
	} else {
		goto L329
	}
L317:
	;
	v1827 = *(*int32)(unsafe.Add(mBase, uint32(v14)+136))
	v1834 = v1824
	v1835 = int32(0)
	v1836 = v1827
	goto L318
L318:
	;
	v1840 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1836))))
	if v1840 == int32(0) {
		v1846 = v1835
		goto L320
	} else {
		goto L321
	}
L319:
	;
	v1852 = int32(0)
	v1859 = v1852
	v1861 = v1852
	goto L324
L320:
	;
	v1851 = v1834 + int32(-1)
	if v1851 != 0 {
		v1834 = v1851
		v1835 = v1846
		v1836 = v1836 + int32(1)
		goto L318
	} else {
		goto L323
	}
L321:
	;
	if int32(0) < v1835 {
		goto L316
	} else {
		goto L322
	}
L322:
	;
	v1846 = int32(1)
	goto L320
L323:
	;
	goto L319
L324:
	;
	v1865 = *(*int32)(unsafe.Add(mBase, uint32(v14)+136))
	v1867 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v1865+v1859))) = uint8(v1867)
	v1869 = *(*int32)(unsafe.Add(mBase, uint32(v14)+140))
	*(*uint16)(unsafe.Add(mBase, uint32(v1869+v1861))) = uint16(v1867)
	v1876 = v1859 + int32(1)
	v1877 = *(*int32)(unsafe.Add(mBase, uint32(v14)+132))
	if v1876 < v1877 {
		v1859 = v1876
		v1861 = v1861 + int32(2)
		goto L324
	} else {
		goto L326
	}
L325:
	;
	goto L316
L326:
	;
	goto L325
L327:
	;
	v1956 = int32(1)
	v1958 = base.B2i32(v1956 < v1954) & v1952
	goto L341
L328:
	;
	v1904 = int32(0)
	v1905 = v1176<<(uint(int32(1))%32) + l2 + int32(-2)
	v1908 = v1176
	goto L331
L329:
	;
	v1952 = int32(0)
	v1954 = v1176
	goto L327
L330:
	;
	v1952 = base.B2i32(int32(12) < v1940)
	v1954 = v1942
	goto L327
L331:
	;
	v1910 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1905))))
	if base.Ui32(int32(18)) < base.Ui32(v1910) {
		v1940 = v1904
		v1942 = v1908
		goto L330
	} else {
		goto L333
	}
L332:
	;
	v1940 = v1930
	v1942 = int32(0)
	goto L330
L333:
	;
	if int32(1)<<(uint(v1910)%32)&int32(393217) == int32(0) {
		v1940 = v1904
		v1942 = v1908
		goto L330
	} else {
		goto L334
	}
L334:
	;
	v1922 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14+int32(192)+v1910))))
	v1923 = v1904 + v1922
	switch v1910 + int32(-17) {
	case 0:
		goto L337
	case 1:
		goto L336
	default:
		v1930 = v1923
		goto L335
	}
L335:
	;
	if base.Ui32(int32(1)) < base.Ui32(v1908+int32(0)) {
		v1904 = v1930
		v1905 = v1905 + int32(-2)
		v1908 = v1908 + int32(-1)
		goto L331
	} else {
		goto L338
	}
L336:
	;
	v1930 = v1923 + int32(7)
	goto L335
L337:
	;
	v1930 = v1923 + int32(3)
	goto L335
L338:
	;
	goto L332
L339:
	;
	if v1958 != int32(1) {
		goto L366
	} else {
		goto L367
	}
L340:
	;
	goto L339
L341:
	;
	v1971 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v1972 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v1972+v1956 < int32(32) {
		goto L344
	} else {
		goto L345
	}
L342:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v2073 + v2071
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v2070<<(uint(v2073)%32) | v2072
	goto L340
L343:
	;
	v1990 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v1991 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v1995 = v1987
	v1996 = v1988
	v1997 = v1991
	v1998 = v1990
	goto L347
L344:
	;
	if v1972 < int32(16) {
		v2070 = v1958
		v2071 = v1956
		v2072 = v1971
		v2073 = v1972
		goto L342
	} else {
		goto L346
	}
L345:
	;
	v1976 = int32(32)
	v1977 = v1976 - v1972
	v1985 = int32(base.Ui32(v1958) >> (uint(v1977) % 32))
	v1986 = v1956 - v1977
	v1987 = v1958<<(uint(v1972)%32) | v1971
	v1988 = v1976
	goto L343
L346:
	;
	v1985 = v1958
	v1986 = v1956
	v1987 = v1971
	v1988 = v1972
	goto L343
L347:
	;
	if base.Ui32(v1997+int32(2)) <= base.Ui32(v1998) {
		v2053 = v1997
		v2054 = v1998
		goto L349
	} else {
		goto L350
	}
L348:
	;
	v2070 = v1985
	v2071 = v1986
	v2072 = v2066
	v2073 = v2064
	goto L342
L349:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2053))) = uint16(v1995)
	v2061 = v2053 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2061
	v2064 = v1996 + int32(-16)
	v2066 = int32(base.Ui32(v1995) >> (uint(int32(16)) % 32))
	if int32(31) < v1996 {
		v1995 = v2066
		v1996 = v2064
		v1997 = v2061
		v1998 = v2054
		goto L347
	} else {
		goto L364
	}
L350:
	;
	v2007 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2008 = v1998 - v2007
	v2011 = base.I64_extend_i32_s(v2008) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2011) {
		v2035 = v2007
		goto L352
	} else {
		goto L353
	}
L351:
	;
	if v1997 == v2007 {
		goto L362
	} else {
		goto L363
	}
L352:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2035
	goto L339
L353:
	;
	v2014 = v1997 - v2007
	v2016 = v2011 + base.I64_extend_i32_u(v2014)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2016) {
		v2035 = v2007
		goto L352
	} else {
		goto L354
	}
L354:
	;
	v2019 = base.I32_wrap_i64(v2016)
	if v1998 == v2007 {
		goto L355
	} else {
		goto L356
	}
L355:
	;
	v2026 = int32(base.Ui32(v2008*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v2019) < base.Ui32(v2026) {
		goto L358
	} else {
		goto L359
	}
L356:
	;
	if base.Ui32(v2019) <= base.Ui32(v2008) {
		v2053 = v1997
		v2054 = v1998
		goto L349
	} else {
		goto L357
	}
L357:
	;
	goto L355
L358:
	;
	v2028 = v2026
	goto L360
L359:
	;
	v2028 = v2019
	goto L360
L360:
	;
	v2032 = v2028&int32(-1024) + int32(1024)
	v2033 = F_WebPSafeMalloc(m, int64(1), v2032)
	mBase = m.M
	if v2033 != 0 {
		goto L351
	} else {
		goto L361
	}
L361:
	;
	v2034 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2035 = v2034
	goto L352
L362:
	;
	v2046 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v2046)
	mBase = m.M
	v2048 = v2033 + v2032
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v2048
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v2033
	v2053 = v2033 + v2014
	v2054 = v2048
	goto L349
L363:
	;
	v2044 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2045 = F_memcpy(m, v2033, v2044, v2014)
	mBase = m.M
	goto L362
L364:
	;
	goto L348
L365:
	;
	v2585 = l2
	v2591 = v2582
	goto L449
L366:
	;
	if v1176 < int32(1) {
		goto L1
	} else {
		goto L448
	}
L367:
	;
	v2112 = int32(2)
	if v1954 != v2112 {
		goto L368
	} else {
		goto L369
	}
L368:
	;
	v2268 = v1954 + int32(-2)
	v2271 = base.I32_clz(v2268) ^ int32(31)
	v2273 = int32(base.Ui32(v2271) >> (uint(int32(1)) % 32))
	v2274 = int32(3)
	goto L398
L369:
	;
	v2115 = int32(0)
	v2116 = int32(5)
	goto L372
L370:
	;
	v2582 = v2112
	goto L365
L371:
	;
	goto L370
L372:
	;
	v2128 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v2129 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v2129+v2116 < int32(32) {
		goto L375
	} else {
		goto L376
	}
L373:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v2230 + v2228
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v2227<<(uint(v2230)%32) | v2229
	goto L371
L374:
	;
	v2147 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v2148 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v2152 = v2144
	v2153 = v2145
	v2154 = v2148
	v2155 = v2147
	goto L378
L375:
	;
	if v2129 < int32(16) {
		v2227 = v2115
		v2228 = v2116
		v2229 = v2128
		v2230 = v2129
		goto L373
	} else {
		goto L377
	}
L376:
	;
	v2133 = int32(32)
	v2134 = v2133 - v2129
	v2142 = int32(base.Ui32(v2115) >> (uint(v2134) % 32))
	v2143 = v2116 - v2134
	v2144 = v2115<<(uint(v2129)%32) | v2128
	v2145 = v2133
	goto L374
L377:
	;
	v2142 = v2115
	v2143 = v2116
	v2144 = v2128
	v2145 = v2129
	goto L374
L378:
	;
	if base.Ui32(v2154+int32(2)) <= base.Ui32(v2155) {
		v2210 = v2154
		v2211 = v2155
		goto L380
	} else {
		goto L381
	}
L379:
	;
	v2227 = v2142
	v2228 = v2143
	v2229 = v2223
	v2230 = v2221
	goto L373
L380:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2210))) = uint16(v2152)
	v2218 = v2210 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2218
	v2221 = v2153 + int32(-16)
	v2223 = int32(base.Ui32(v2152) >> (uint(int32(16)) % 32))
	if int32(31) < v2153 {
		v2152 = v2223
		v2153 = v2221
		v2154 = v2218
		v2155 = v2211
		goto L378
	} else {
		goto L395
	}
L381:
	;
	v2164 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2165 = v2155 - v2164
	v2168 = base.I64_extend_i32_s(v2165) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2168) {
		v2192 = v2164
		goto L383
	} else {
		goto L384
	}
L382:
	;
	if v2154 == v2164 {
		goto L393
	} else {
		goto L394
	}
L383:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2192
	goto L370
L384:
	;
	v2171 = v2154 - v2164
	v2173 = v2168 + base.I64_extend_i32_u(v2171)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2173) {
		v2192 = v2164
		goto L383
	} else {
		goto L385
	}
L385:
	;
	v2176 = base.I32_wrap_i64(v2173)
	if v2155 == v2164 {
		goto L386
	} else {
		goto L387
	}
L386:
	;
	v2183 = int32(base.Ui32(v2165*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v2176) < base.Ui32(v2183) {
		goto L389
	} else {
		goto L390
	}
L387:
	;
	if base.Ui32(v2176) <= base.Ui32(v2165) {
		v2210 = v2154
		v2211 = v2155
		goto L380
	} else {
		goto L388
	}
L388:
	;
	goto L386
L389:
	;
	v2185 = v2183
	goto L391
L390:
	;
	v2185 = v2176
	goto L391
L391:
	;
	v2189 = v2185&int32(-1024) + int32(1024)
	v2190 = F_WebPSafeMalloc(m, int64(1), v2189)
	mBase = m.M
	if v2190 != 0 {
		goto L382
	} else {
		goto L392
	}
L392:
	;
	v2191 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2192 = v2191
	goto L383
L393:
	;
	v2203 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v2203)
	mBase = m.M
	v2205 = v2190 + v2189
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v2205
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v2190
	v2210 = v2190 + v2171
	v2211 = v2205
	goto L380
L394:
	;
	v2201 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2202 = F_memcpy(m, v2190, v2201, v2171)
	mBase = m.M
	goto L393
L395:
	;
	goto L379
L396:
	;
	v2428 = v2271&int32(62) + int32(2)
	if v2428 < int32(1) {
		goto L423
	} else {
		goto L424
	}
L397:
	;
	goto L396
L398:
	;
	v2286 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v2287 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v2287+v2274 < int32(32) {
		goto L401
	} else {
		goto L402
	}
L399:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v2388 + v2386
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v2385<<(uint(v2388)%32) | v2387
	goto L397
L400:
	;
	v2305 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v2306 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v2310 = v2302
	v2311 = v2303
	v2312 = v2306
	v2313 = v2305
	goto L404
L401:
	;
	if v2287 < int32(16) {
		v2385 = v2273
		v2386 = v2274
		v2387 = v2286
		v2388 = v2287
		goto L399
	} else {
		goto L403
	}
L402:
	;
	v2291 = int32(32)
	v2292 = v2291 - v2287
	v2300 = int32(base.Ui32(v2273) >> (uint(v2292) % 32))
	v2301 = v2274 - v2292
	v2302 = v2273<<(uint(v2287)%32) | v2286
	v2303 = v2291
	goto L400
L403:
	;
	v2300 = v2273
	v2301 = v2274
	v2302 = v2286
	v2303 = v2287
	goto L400
L404:
	;
	if base.Ui32(v2312+int32(2)) <= base.Ui32(v2313) {
		v2368 = v2312
		v2369 = v2313
		goto L406
	} else {
		goto L407
	}
L405:
	;
	v2385 = v2300
	v2386 = v2301
	v2387 = v2381
	v2388 = v2379
	goto L399
L406:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2368))) = uint16(v2310)
	v2376 = v2368 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2376
	v2379 = v2311 + int32(-16)
	v2381 = int32(base.Ui32(v2310) >> (uint(int32(16)) % 32))
	if int32(31) < v2311 {
		v2310 = v2381
		v2311 = v2379
		v2312 = v2376
		v2313 = v2369
		goto L404
	} else {
		goto L421
	}
L407:
	;
	v2322 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2323 = v2313 - v2322
	v2326 = base.I64_extend_i32_s(v2323) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2326) {
		v2350 = v2322
		goto L409
	} else {
		goto L410
	}
L408:
	;
	if v2312 == v2322 {
		goto L419
	} else {
		goto L420
	}
L409:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2350
	goto L396
L410:
	;
	v2329 = v2312 - v2322
	v2331 = v2326 + base.I64_extend_i32_u(v2329)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2331) {
		v2350 = v2322
		goto L409
	} else {
		goto L411
	}
L411:
	;
	v2334 = base.I32_wrap_i64(v2331)
	if v2313 == v2322 {
		goto L412
	} else {
		goto L413
	}
L412:
	;
	v2341 = int32(base.Ui32(v2323*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v2334) < base.Ui32(v2341) {
		goto L415
	} else {
		goto L416
	}
L413:
	;
	if base.Ui32(v2334) <= base.Ui32(v2323) {
		v2368 = v2312
		v2369 = v2313
		goto L406
	} else {
		goto L414
	}
L414:
	;
	goto L412
L415:
	;
	v2343 = v2341
	goto L417
L416:
	;
	v2343 = v2334
	goto L417
L417:
	;
	v2347 = v2343&int32(-1024) + int32(1024)
	v2348 = F_WebPSafeMalloc(m, int64(1), v2347)
	mBase = m.M
	if v2348 != 0 {
		goto L408
	} else {
		goto L418
	}
L418:
	;
	v2349 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2350 = v2349
	goto L409
L419:
	;
	v2361 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v2361)
	mBase = m.M
	v2363 = v2348 + v2347
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v2363
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v2348
	v2368 = v2348 + v2329
	v2369 = v2363
	goto L406
L420:
	;
	v2359 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2360 = F_memcpy(m, v2348, v2359, v2329)
	mBase = m.M
	goto L419
L421:
	;
	goto L405
L422:
	;
	v2582 = v1954
	goto L365
L423:
	;
	goto L422
L424:
	;
	v2440 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v2441 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v2441+v2428 < int32(32) {
		goto L427
	} else {
		goto L428
	}
L425:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v2542 + v2540
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v2539<<(uint(v2542)%32) | v2541
	goto L423
L426:
	;
	v2459 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v2460 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v2464 = v2456
	v2465 = v2457
	v2466 = v2460
	v2467 = v2459
	goto L430
L427:
	;
	if v2441 < int32(16) {
		v2539 = v2268
		v2540 = v2428
		v2541 = v2440
		v2542 = v2441
		goto L425
	} else {
		goto L429
	}
L428:
	;
	v2445 = int32(32)
	v2446 = v2445 - v2441
	v2454 = int32(base.Ui32(v2268) >> (uint(v2446) % 32))
	v2455 = v2428 - v2446
	v2456 = v2268<<(uint(v2441)%32) | v2440
	v2457 = v2445
	goto L426
L429:
	;
	v2454 = v2268
	v2455 = v2428
	v2456 = v2440
	v2457 = v2441
	goto L426
L430:
	;
	if base.Ui32(v2466+int32(2)) <= base.Ui32(v2467) {
		v2522 = v2466
		v2523 = v2467
		goto L432
	} else {
		goto L433
	}
L431:
	;
	v2539 = v2454
	v2540 = v2455
	v2541 = v2535
	v2542 = v2533
	goto L425
L432:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2522))) = uint16(v2464)
	v2530 = v2522 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2530
	v2533 = v2465 + int32(-16)
	v2535 = int32(base.Ui32(v2464) >> (uint(int32(16)) % 32))
	if int32(31) < v2465 {
		v2464 = v2535
		v2465 = v2533
		v2466 = v2530
		v2467 = v2523
		goto L430
	} else {
		goto L447
	}
L433:
	;
	v2476 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2477 = v2467 - v2476
	v2480 = base.I64_extend_i32_s(v2477) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2480) {
		v2504 = v2476
		goto L435
	} else {
		goto L436
	}
L434:
	;
	if v2466 == v2476 {
		goto L445
	} else {
		goto L446
	}
L435:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2504
	goto L422
L436:
	;
	v2483 = v2466 - v2476
	v2485 = v2480 + base.I64_extend_i32_u(v2483)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2485) {
		v2504 = v2476
		goto L435
	} else {
		goto L437
	}
L437:
	;
	v2488 = base.I32_wrap_i64(v2485)
	if v2467 == v2476 {
		goto L438
	} else {
		goto L439
	}
L438:
	;
	v2495 = int32(base.Ui32(v2477*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v2488) < base.Ui32(v2495) {
		goto L441
	} else {
		goto L442
	}
L439:
	;
	if base.Ui32(v2488) <= base.Ui32(v2477) {
		v2522 = v2466
		v2523 = v2467
		goto L432
	} else {
		goto L440
	}
L440:
	;
	goto L438
L441:
	;
	v2497 = v2495
	goto L443
L442:
	;
	v2497 = v2488
	goto L443
L443:
	;
	v2501 = v2497&int32(-1024) + int32(1024)
	v2502 = F_WebPSafeMalloc(m, int64(1), v2501)
	mBase = m.M
	if v2502 != 0 {
		goto L434
	} else {
		goto L444
	}
L444:
	;
	v2503 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2504 = v2503
	goto L435
L445:
	;
	v2515 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v2515)
	mBase = m.M
	v2517 = v2502 + v2501
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v2517
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v2502
	v2522 = v2502 + v2483
	v2523 = v2517
	goto L432
L446:
	;
	v2513 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2514 = F_memcpy(m, v2502, v2513, v2483)
	mBase = m.M
	goto L445
L447:
	;
	goto L431
L448:
	;
	v2582 = v1176
	goto L365
L449:
	;
	v2594 = int32(1)
	v2596 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2585+v2594))))
	v2597 = *(*int32)(unsafe.Add(mBase, uint32(v14)+140))
	v2598 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2585))))
	v2602 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2597+v2598<<(uint(v2594)%32)))))
	v2603 = *(*int32)(unsafe.Add(mBase, uint32(v14)+136))
	v2605 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2603+v2598))))
	if v2605 < v2594 {
		goto L452
	} else {
		goto L453
	}
L450:
	;
	goto L1
L451:
	;
	v2759 = (v2598 + int32(-16)) & int32(255)
	if base.Ui32(int32(2)) < base.Ui32(v2759) {
		goto L477
	} else {
		goto L478
	}
L452:
	;
	goto L451
L453:
	;
	v2617 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v2618 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v2618+v2605 < int32(32) {
		goto L456
	} else {
		goto L457
	}
L454:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v2719 + v2717
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v2716<<(uint(v2719)%32) | v2718
	goto L452
L455:
	;
	v2636 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v2637 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v2641 = v2633
	v2642 = v2634
	v2643 = v2637
	v2644 = v2636
	goto L459
L456:
	;
	if v2618 < int32(16) {
		v2716 = v2602
		v2717 = v2605
		v2718 = v2617
		v2719 = v2618
		goto L454
	} else {
		goto L458
	}
L457:
	;
	v2622 = int32(32)
	v2623 = v2622 - v2618
	v2631 = int32(base.Ui32(v2602) >> (uint(v2623) % 32))
	v2632 = v2605 - v2623
	v2633 = v2602<<(uint(v2618)%32) | v2617
	v2634 = v2622
	goto L455
L458:
	;
	v2631 = v2602
	v2632 = v2605
	v2633 = v2617
	v2634 = v2618
	goto L455
L459:
	;
	if base.Ui32(v2643+int32(2)) <= base.Ui32(v2644) {
		v2699 = v2643
		v2700 = v2644
		goto L461
	} else {
		goto L462
	}
L460:
	;
	v2716 = v2631
	v2717 = v2632
	v2718 = v2712
	v2719 = v2710
	goto L454
L461:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2699))) = uint16(v2641)
	v2707 = v2699 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2707
	v2710 = v2642 + int32(-16)
	v2712 = int32(base.Ui32(v2641) >> (uint(int32(16)) % 32))
	if int32(31) < v2642 {
		v2641 = v2712
		v2642 = v2710
		v2643 = v2707
		v2644 = v2700
		goto L459
	} else {
		goto L476
	}
L462:
	;
	v2653 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2654 = v2644 - v2653
	v2657 = base.I64_extend_i32_s(v2654) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2657) {
		v2681 = v2653
		goto L464
	} else {
		goto L465
	}
L463:
	;
	if v2643 == v2653 {
		goto L474
	} else {
		goto L475
	}
L464:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2681
	goto L451
L465:
	;
	v2660 = v2643 - v2653
	v2662 = v2657 + base.I64_extend_i32_u(v2660)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2662) {
		v2681 = v2653
		goto L464
	} else {
		goto L466
	}
L466:
	;
	v2665 = base.I32_wrap_i64(v2662)
	if v2644 == v2653 {
		goto L467
	} else {
		goto L468
	}
L467:
	;
	v2672 = int32(base.Ui32(v2654*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v2665) < base.Ui32(v2672) {
		goto L470
	} else {
		goto L471
	}
L468:
	;
	if base.Ui32(v2665) <= base.Ui32(v2654) {
		v2699 = v2643
		v2700 = v2644
		goto L461
	} else {
		goto L469
	}
L469:
	;
	goto L467
L470:
	;
	v2674 = v2672
	goto L472
L471:
	;
	v2674 = v2665
	goto L472
L472:
	;
	v2678 = v2674&int32(-1024) + int32(1024)
	v2679 = F_WebPSafeMalloc(m, int64(1), v2678)
	mBase = m.M
	if v2679 != 0 {
		goto L463
	} else {
		goto L473
	}
L473:
	;
	v2680 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2681 = v2680
	goto L464
L474:
	;
	v2692 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v2692)
	mBase = m.M
	v2694 = v2679 + v2678
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v2694
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v2679
	v2699 = v2679 + v2660
	v2700 = v2694
	goto L461
L475:
	;
	v2690 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2691 = F_memcpy(m, v2679, v2690, v2660)
	mBase = m.M
	goto L474
L476:
	;
	goto L460
L477:
	;
	v2922 = v2591 + int32(-1)
	if v2922 != 0 {
		v2585 = v2585 + int32(2)
		v2591 = v2922
		goto L449
	} else {
		goto L505
	}
L478:
	;
	v2762 = m.G1
	v2768 = *(*int32)(unsafe.Add(mBase, uint32(v2762+int32(_a_F_StoreHuffmanCode_1)+v2759<<(uint(int32(2))%32))))
	if v2768 < int32(1) {
		goto L480
	} else {
		goto L481
	}
L479:
	;
	goto L477
L480:
	;
	goto L479
L481:
	;
	v2780 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v2781 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v2781+v2768 < int32(32) {
		goto L484
	} else {
		goto L485
	}
L482:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v2882 + v2880
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v2879<<(uint(v2882)%32) | v2881
	goto L480
L483:
	;
	v2799 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v2800 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v2804 = v2796
	v2805 = v2797
	v2806 = v2800
	v2807 = v2799
	goto L487
L484:
	;
	if v2781 < int32(16) {
		v2879 = v2596
		v2880 = v2768
		v2881 = v2780
		v2882 = v2781
		goto L482
	} else {
		goto L486
	}
L485:
	;
	v2785 = int32(32)
	v2786 = v2785 - v2781
	v2794 = int32(base.Ui32(v2596) >> (uint(v2786) % 32))
	v2795 = v2768 - v2786
	v2796 = v2596<<(uint(v2781)%32) | v2780
	v2797 = v2785
	goto L483
L486:
	;
	v2794 = v2596
	v2795 = v2768
	v2796 = v2780
	v2797 = v2781
	goto L483
L487:
	;
	if base.Ui32(v2806+int32(2)) <= base.Ui32(v2807) {
		v2862 = v2806
		v2863 = v2807
		goto L489
	} else {
		goto L490
	}
L488:
	;
	v2879 = v2794
	v2880 = v2795
	v2881 = v2875
	v2882 = v2873
	goto L482
L489:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2862))) = uint16(v2804)
	v2870 = v2862 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2870
	v2873 = v2805 + int32(-16)
	v2875 = int32(base.Ui32(v2804) >> (uint(int32(16)) % 32))
	if int32(31) < v2805 {
		v2804 = v2875
		v2805 = v2873
		v2806 = v2870
		v2807 = v2863
		goto L487
	} else {
		goto L504
	}
L490:
	;
	v2816 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2817 = v2807 - v2816
	v2820 = base.I64_extend_i32_s(v2817) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2820) {
		v2844 = v2816
		goto L492
	} else {
		goto L493
	}
L491:
	;
	if v2806 == v2816 {
		goto L502
	} else {
		goto L503
	}
L492:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2844
	goto L479
L493:
	;
	v2823 = v2806 - v2816
	v2825 = v2820 + base.I64_extend_i32_u(v2823)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2825) {
		v2844 = v2816
		goto L492
	} else {
		goto L494
	}
L494:
	;
	v2828 = base.I32_wrap_i64(v2825)
	if v2807 == v2816 {
		goto L495
	} else {
		goto L496
	}
L495:
	;
	v2835 = int32(base.Ui32(v2817*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v2828) < base.Ui32(v2835) {
		goto L498
	} else {
		goto L499
	}
L496:
	;
	if base.Ui32(v2828) <= base.Ui32(v2817) {
		v2862 = v2806
		v2863 = v2807
		goto L489
	} else {
		goto L497
	}
L497:
	;
	goto L495
L498:
	;
	v2837 = v2835
	goto L500
L499:
	;
	v2837 = v2828
	goto L500
L500:
	;
	v2841 = v2837&int32(-1024) + int32(1024)
	v2842 = F_WebPSafeMalloc(m, int64(1), v2841)
	mBase = m.M
	if v2842 != 0 {
		goto L491
	} else {
		goto L501
	}
L501:
	;
	v2843 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2844 = v2843
	goto L492
L502:
	;
	v2855 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v2855)
	mBase = m.M
	v2857 = v2842 + v2841
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v2857
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v2842
	v2862 = v2842 + v2823
	v2863 = v2857
	goto L489
L503:
	;
	v2853 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2854 = F_memcpy(m, v2842, v2853, v2823)
	mBase = m.M
	goto L502
L504:
	;
	goto L488
L505:
	;
	goto L450
}
func F_StoreImageToBitMask(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32) int32 {
	mBase := m.M
	_ = mBase
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
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
	var v35 int32
	_ = v35
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v44 int32
	_ = v44
	var v47 int32
	_ = v47
	var v51 int32
	_ = v51
	var v52 int32
	_ = v52
	var v59 int32
	_ = v59
	var v67 int32
	_ = v67
	var v73 int32
	_ = v73
	var v78 int32
	_ = v78
	var v80 int32
	_ = v80
	var v83 int32
	_ = v83
	var v84 int32
	_ = v84
	var v85 int32
	_ = v85
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v97 int32
	_ = v97
	var v104 int32
	_ = v104
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
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
	var v117 int32
	_ = v117
	var v121 int32
	_ = v121
	var v122 int32
	_ = v122
	var v125 int32
	_ = v125
	var v128 int32
	_ = v128
	var v130 int32
	_ = v130
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v147 int32
	_ = v147
	var v148 int32
	_ = v148
	var v156 int32
	_ = v156
	var v157 int32
	_ = v157
	var v158 int32
	_ = v158
	var v159 int32
	_ = v159
	var v161 int32
	_ = v161
	var v162 int32
	_ = v162
	var v166 int32
	_ = v166
	var v167 int32
	_ = v167
	var v168 int32
	_ = v168
	var v169 int32
	_ = v169
	var v178 int32
	_ = v178
	var v179 int32
	_ = v179
	var v182 int64
	_ = v182
	var v185 int32
	_ = v185
	var v187 int64
	_ = v187
	var v190 int32
	_ = v190
	var v197 int32
	_ = v197
	var v199 int32
	_ = v199
	var v203 int32
	_ = v203
	var v204 int32
	_ = v204
	var v205 int32
	_ = v205
	var v206 int32
	_ = v206
	var v215 int32
	_ = v215
	var v216 int32
	_ = v216
	var v217 int32
	_ = v217
	var v219 int32
	_ = v219
	var v224 int32
	_ = v224
	var v225 int32
	_ = v225
	var v232 int32
	_ = v232
	var v235 int32
	_ = v235
	var v237 int32
	_ = v237
	var v241 int32
	_ = v241
	var v242 int32
	_ = v242
	var v243 int32
	_ = v243
	var v244 int32
	_ = v244
	var v283 int32
	_ = v283
	var v286 int32
	_ = v286
	var v287 int32
	_ = v287
	var v290 int32
	_ = v290
	var v293 int32
	_ = v293
	var v295 int32
	_ = v295
	var v307 int32
	_ = v307
	var v308 int32
	_ = v308
	var v312 int32
	_ = v312
	var v313 int32
	_ = v313
	var v321 int32
	_ = v321
	var v322 int32
	_ = v322
	var v323 int32
	_ = v323
	var v324 int32
	_ = v324
	var v326 int32
	_ = v326
	var v327 int32
	_ = v327
	var v331 int32
	_ = v331
	var v332 int32
	_ = v332
	var v333 int32
	_ = v333
	var v334 int32
	_ = v334
	var v343 int32
	_ = v343
	var v344 int32
	_ = v344
	var v347 int64
	_ = v347
	var v350 int32
	_ = v350
	var v352 int64
	_ = v352
	var v355 int32
	_ = v355
	var v362 int32
	_ = v362
	var v364 int32
	_ = v364
	var v368 int32
	_ = v368
	var v369 int32
	_ = v369
	var v370 int32
	_ = v370
	var v371 int32
	_ = v371
	var v380 int32
	_ = v380
	var v381 int32
	_ = v381
	var v382 int32
	_ = v382
	var v384 int32
	_ = v384
	var v389 int32
	_ = v389
	var v390 int32
	_ = v390
	var v397 int32
	_ = v397
	var v400 int32
	_ = v400
	var v402 int32
	_ = v402
	var v406 int32
	_ = v406
	var v407 int32
	_ = v407
	var v408 int32
	_ = v408
	var v409 int32
	_ = v409
	var v448 int32
	_ = v448
	var v449 int32
	_ = v449
	var v450 int32
	_ = v450
	var v453 int32
	_ = v453
	var v456 int32
	_ = v456
	var v458 int32
	_ = v458
	var v470 int32
	_ = v470
	var v471 int32
	_ = v471
	var v475 int32
	_ = v475
	var v476 int32
	_ = v476
	var v484 int32
	_ = v484
	var v485 int32
	_ = v485
	var v486 int32
	_ = v486
	var v487 int32
	_ = v487
	var v489 int32
	_ = v489
	var v490 int32
	_ = v490
	var v494 int32
	_ = v494
	var v495 int32
	_ = v495
	var v496 int32
	_ = v496
	var v497 int32
	_ = v497
	var v506 int32
	_ = v506
	var v507 int32
	_ = v507
	var v510 int64
	_ = v510
	var v513 int32
	_ = v513
	var v515 int64
	_ = v515
	var v518 int32
	_ = v518
	var v525 int32
	_ = v525
	var v527 int32
	_ = v527
	var v531 int32
	_ = v531
	var v532 int32
	_ = v532
	var v533 int32
	_ = v533
	var v534 int32
	_ = v534
	var v543 int32
	_ = v543
	var v544 int32
	_ = v544
	var v545 int32
	_ = v545
	var v547 int32
	_ = v547
	var v552 int32
	_ = v552
	var v553 int32
	_ = v553
	var v560 int32
	_ = v560
	var v563 int32
	_ = v563
	var v565 int32
	_ = v565
	var v569 int32
	_ = v569
	var v570 int32
	_ = v570
	var v571 int32
	_ = v571
	var v572 int32
	_ = v572
	var v611 int32
	_ = v611
	var v612 int32
	_ = v612
	var v614 int32
	_ = v614
	var v615 int32
	_ = v615
	var v618 int32
	_ = v618
	var v621 int32
	_ = v621
	var v623 int32
	_ = v623
	var v635 int32
	_ = v635
	var v636 int32
	_ = v636
	var v640 int32
	_ = v640
	var v641 int32
	_ = v641
	var v649 int32
	_ = v649
	var v650 int32
	_ = v650
	var v651 int32
	_ = v651
	var v652 int32
	_ = v652
	var v654 int32
	_ = v654
	var v655 int32
	_ = v655
	var v659 int32
	_ = v659
	var v660 int32
	_ = v660
	var v661 int32
	_ = v661
	var v662 int32
	_ = v662
	var v671 int32
	_ = v671
	var v672 int32
	_ = v672
	var v675 int64
	_ = v675
	var v678 int32
	_ = v678
	var v680 int64
	_ = v680
	var v683 int32
	_ = v683
	var v690 int32
	_ = v690
	var v692 int32
	_ = v692
	var v696 int32
	_ = v696
	var v697 int32
	_ = v697
	var v698 int32
	_ = v698
	var v699 int32
	_ = v699
	var v708 int32
	_ = v708
	var v709 int32
	_ = v709
	var v710 int32
	_ = v710
	var v712 int32
	_ = v712
	var v717 int32
	_ = v717
	var v718 int32
	_ = v718
	var v725 int32
	_ = v725
	var v728 int32
	_ = v728
	var v730 int32
	_ = v730
	var v734 int32
	_ = v734
	var v735 int32
	_ = v735
	var v736 int32
	_ = v736
	var v737 int32
	_ = v737
	var v776 int32
	_ = v776
	var v778 int32
	_ = v778
	var v779 int32
	_ = v779
	var v782 int32
	_ = v782
	var v785 int32
	_ = v785
	var v787 int32
	_ = v787
	var v799 int32
	_ = v799
	var v800 int32
	_ = v800
	var v804 int32
	_ = v804
	var v805 int32
	_ = v805
	var v813 int32
	_ = v813
	var v814 int32
	_ = v814
	var v815 int32
	_ = v815
	var v816 int32
	_ = v816
	var v818 int32
	_ = v818
	var v819 int32
	_ = v819
	var v823 int32
	_ = v823
	var v824 int32
	_ = v824
	var v825 int32
	_ = v825
	var v826 int32
	_ = v826
	var v835 int32
	_ = v835
	var v836 int32
	_ = v836
	var v839 int64
	_ = v839
	var v842 int32
	_ = v842
	var v844 int64
	_ = v844
	var v847 int32
	_ = v847
	var v854 int32
	_ = v854
	var v856 int32
	_ = v856
	var v860 int32
	_ = v860
	var v861 int32
	_ = v861
	var v862 int32
	_ = v862
	var v863 int32
	_ = v863
	var v872 int32
	_ = v872
	var v873 int32
	_ = v873
	var v874 int32
	_ = v874
	var v876 int32
	_ = v876
	var v881 int32
	_ = v881
	var v882 int32
	_ = v882
	var v889 int32
	_ = v889
	var v892 int32
	_ = v892
	var v894 int32
	_ = v894
	var v898 int32
	_ = v898
	var v899 int32
	_ = v899
	var v900 int32
	_ = v900
	var v901 int32
	_ = v901
	var v938 int32
	_ = v938
	var v941 int32
	_ = v941
	var v943 int32
	_ = v943
	var v944 int32
	_ = v944
	var v947 int32
	_ = v947
	var v948 int32
	_ = v948
	var v949 int32
	_ = v949
	var v950 int32
	_ = v950
	var v951 int32
	_ = v951
	var v954 int32
	_ = v954
	var v956 int32
	_ = v956
	var v958 int32
	_ = v958
	var v968 int32
	_ = v968
	var v969 int32
	_ = v969
	var v970 int32
	_ = v970
	var v973 int32
	_ = v973
	var v975 int32
	_ = v975
	var v977 int32
	_ = v977
	var v981 int32
	_ = v981
	var v982 int32
	_ = v982
	var v985 int32
	_ = v985
	var v986 int32
	_ = v986
	var v987 int32
	_ = v987
	var v999 int32
	_ = v999
	var v1000 int32
	_ = v1000
	var v1004 int32
	_ = v1004
	var v1005 int32
	_ = v1005
	var v1013 int32
	_ = v1013
	var v1014 int32
	_ = v1014
	var v1015 int32
	_ = v1015
	var v1016 int32
	_ = v1016
	var v1018 int32
	_ = v1018
	var v1019 int32
	_ = v1019
	var v1023 int32
	_ = v1023
	var v1024 int32
	_ = v1024
	var v1025 int32
	_ = v1025
	var v1026 int32
	_ = v1026
	var v1035 int32
	_ = v1035
	var v1036 int32
	_ = v1036
	var v1039 int64
	_ = v1039
	var v1042 int32
	_ = v1042
	var v1044 int64
	_ = v1044
	var v1047 int32
	_ = v1047
	var v1054 int32
	_ = v1054
	var v1056 int32
	_ = v1056
	var v1060 int32
	_ = v1060
	var v1061 int32
	_ = v1061
	var v1062 int32
	_ = v1062
	var v1063 int32
	_ = v1063
	var v1072 int32
	_ = v1072
	var v1073 int32
	_ = v1073
	var v1074 int32
	_ = v1074
	var v1076 int32
	_ = v1076
	var v1081 int32
	_ = v1081
	var v1082 int32
	_ = v1082
	var v1089 int32
	_ = v1089
	var v1092 int32
	_ = v1092
	var v1094 int32
	_ = v1094
	var v1098 int32
	_ = v1098
	var v1099 int32
	_ = v1099
	var v1100 int32
	_ = v1100
	var v1101 int32
	_ = v1101
	var v1140 int32
	_ = v1140
	var v1142 int32
	_ = v1142
	var v1143 int32
	_ = v1143
	var v1146 int32
	_ = v1146
	var v1147 int32
	_ = v1147
	var v1148 int32
	_ = v1148
	var v1149 int32
	_ = v1149
	var v1150 int32
	_ = v1150
	var v1153 int32
	_ = v1153
	var v1155 int32
	_ = v1155
	var v1157 int32
	_ = v1157
	var v1167 int32
	_ = v1167
	var v1168 int32
	_ = v1168
	var v1169 int32
	_ = v1169
	var v1172 int32
	_ = v1172
	var v1173 int32
	_ = v1173
	var v1176 int32
	_ = v1176
	var v1179 int32
	_ = v1179
	var v1181 int32
	_ = v1181
	var v1193 int32
	_ = v1193
	var v1194 int32
	_ = v1194
	var v1198 int32
	_ = v1198
	var v1199 int32
	_ = v1199
	var v1207 int32
	_ = v1207
	var v1208 int32
	_ = v1208
	var v1209 int32
	_ = v1209
	var v1210 int32
	_ = v1210
	var v1212 int32
	_ = v1212
	var v1213 int32
	_ = v1213
	var v1217 int32
	_ = v1217
	var v1218 int32
	_ = v1218
	var v1219 int32
	_ = v1219
	var v1220 int32
	_ = v1220
	var v1229 int32
	_ = v1229
	var v1230 int32
	_ = v1230
	var v1233 int64
	_ = v1233
	var v1236 int32
	_ = v1236
	var v1238 int64
	_ = v1238
	var v1241 int32
	_ = v1241
	var v1248 int32
	_ = v1248
	var v1250 int32
	_ = v1250
	var v1254 int32
	_ = v1254
	var v1255 int32
	_ = v1255
	var v1256 int32
	_ = v1256
	var v1257 int32
	_ = v1257
	var v1266 int32
	_ = v1266
	var v1267 int32
	_ = v1267
	var v1268 int32
	_ = v1268
	var v1270 int32
	_ = v1270
	var v1275 int32
	_ = v1275
	var v1276 int32
	_ = v1276
	var v1283 int32
	_ = v1283
	var v1286 int32
	_ = v1286
	var v1288 int32
	_ = v1288
	var v1292 int32
	_ = v1292
	var v1293 int32
	_ = v1293
	var v1294 int32
	_ = v1294
	var v1295 int32
	_ = v1295
	var v1343 int32
	_ = v1343
	var v1344 int32
	_ = v1344
	var v1348 int32
	_ = v1348
	var v1349 int32
	_ = v1349
	var v1357 int32
	_ = v1357
	var v1358 int32
	_ = v1358
	var v1359 int32
	_ = v1359
	var v1360 int32
	_ = v1360
	var v1362 int32
	_ = v1362
	var v1363 int32
	_ = v1363
	var v1367 int32
	_ = v1367
	var v1368 int32
	_ = v1368
	var v1369 int32
	_ = v1369
	var v1370 int32
	_ = v1370
	var v1379 int32
	_ = v1379
	var v1380 int32
	_ = v1380
	var v1383 int64
	_ = v1383
	var v1386 int32
	_ = v1386
	var v1388 int64
	_ = v1388
	var v1391 int32
	_ = v1391
	var v1398 int32
	_ = v1398
	var v1400 int32
	_ = v1400
	var v1404 int32
	_ = v1404
	var v1405 int32
	_ = v1405
	var v1406 int32
	_ = v1406
	var v1407 int32
	_ = v1407
	var v1416 int32
	_ = v1416
	var v1417 int32
	_ = v1417
	var v1418 int32
	_ = v1418
	var v1420 int32
	_ = v1420
	var v1425 int32
	_ = v1425
	var v1426 int32
	_ = v1426
	var v1433 int32
	_ = v1433
	var v1436 int32
	_ = v1436
	var v1438 int32
	_ = v1438
	var v1442 int32
	_ = v1442
	var v1443 int32
	_ = v1443
	var v1444 int32
	_ = v1444
	var v1445 int32
	_ = v1445
	var v1489 int32
	_ = v1489
	var v1490 int32
	_ = v1490
	var v1495 int32
	_ = v1495
	var v1500 int32
	_ = v1500
	var v1514 int32
	_ = v1514
	var v1515 int32
	_ = v1515
	var v1520 int32
	_ = v1520
	var v1525 int32
	_ = v1525
	var v1538 int32
	_ = v1538
	var v1540 int32
	_ = v1540
	var v1542 int32
	_ = v1542
	var v1545 int32
	_ = v1545
	var v1549 int32
	_ = v1549
	var v1550 int32
	_ = v1550
	var v1551 int32
	_ = v1551
	var v1553 int32
	_ = v1553
	var v1554 int32
	_ = v1554
	var v1558 int32
	_ = v1558
	var v1559 int32
	_ = v1559
	var v1563 int32
	_ = v1563
	var v1564 int32
	_ = v1564
	var v1586 int32
	_ = v1586
	var v1590 int32
	_ = v1590
	var v1593 int32
	_ = v1593
	v22 = m.G0
	v24 = v22 - int32(16)
	m.G0 = v24
	v26 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
	v28 = v24 + int32(4)
	v30 = *(*int32)(unsafe.Add(mBase, uint32(l3)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v28)+4)) = v30
	if v30 != 0 {
		goto L3
	} else {
		goto L4
	}
L1:
	;
	v44 = *(*int32)(unsafe.Add(mBase, uint32(v24)+4))
	if v44 == int32(0) {
		goto L5
	} else {
		goto L6
	}
L2:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v28)+8)) = v39
	*(*int32)(unsafe.Add(mBase, uint32(v28))) = v40
	goto L1
L3:
	;
	v34 = *(*int32)(unsafe.Add(mBase, uint32(v30)+4))
	v35 = *(*int32)(unsafe.Add(mBase, uint32(v30)+8))
	v39 = v34 + v35<<(uint(int32(3))%32)
	v40 = v34
	goto L2
L4:
	;
	v32 = int32(0)
	v39 = v32
	v40 = v32
	goto L2
L5:
	;
	v1586 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	if v1586 == int32(0) {
		v1593 = int32(1)
		goto L249
	} else {
		goto L250
	}
L6:
	;
	v47 = int32(0)
	if l2 != 0 {
		goto L7
	} else {
		goto L8
	}
L7:
	;
	v51 = int32(-1) << (uint(l2) % 32)
	goto L9
L8:
	;
	v51 = v47
	goto L9
L9:
	;
	v52 = int32(1)
	if l2 != 0 {
		goto L10
	} else {
		goto L11
	}
L10:
	;
	v59 = int32(base.Ui32(l1+v52<<(uint(l2)%32)+int32(-1)) >> (uint(l2) % 32))
	goto L12
L11:
	;
	v59 = v52
	goto L12
L12:
	;
	v67 = int32(0)
	v73 = v47
	v78 = v67
	v80 = v44
	v83 = l5 + int32(base.Ui32(v26)>>(uint(int32(8))%32))&int32(_a_F_StoreImageToBitMask_0)*int32(60)
	v84 = v67
	v85 = v67
	goto L13
L13:
	;
	v91 = v78 & v51
	v92 = v73 & v51
	if v84 != v92 {
		goto L16
	} else {
		goto L17
	}
L14:
	;
	goto L5
L15:
	;
	v112 = v80 + int32(4)
	v113 = *(*int32)(unsafe.Add(mBase, uint32(v112)))
	v114 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v80))))
	switch v114 {
	case 0:
		goto L22
	case 1:
		goto L21
	default:
		goto L20
	}
L16:
	;
	v97 = int32(2)
	v104 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l4+v78>>(uint(l2)%32)*v59<<(uint(v97)%32)+v73>>(uint(l2)%32)<<(uint(v97)%32))+1)))
	v108 = l5 + v104*int32(60)
	v109 = v92
	v110 = v91
	goto L15
L17:
	;
	if v85 == v91 {
		v108 = v83
		v109 = v84
		v110 = v85
		goto L15
	} else {
		goto L18
	}
L18:
	;
	goto L16
L19:
	;
	v1489 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v80+int32(2)))))
	v1490 = v73 + v1489
	if v1490 < l1 {
		v1520 = v1490
		v1525 = v78
		goto L237
	} else {
		goto L238
	}
L20:
	;
	v938 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v80)+2)))
	if base.Ui32(int32(511)) < base.Ui32(v938) {
		goto L154
	} else {
		goto L155
	}
L21:
	;
	v776 = *(*int32)(unsafe.Add(mBase, uint32(v108+int32(8))))
	v778 = v113 + int32(280)
	v779 = int32(1)
	v782 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v776+v778<<(uint(v779)%32)))))
	v785 = *(*int32)(unsafe.Add(mBase, uint32(v108+int32(4))))
	v787 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v785+v778))))
	if v787 < v779 {
		goto L128
	} else {
		goto L129
	}
L22:
	;
	v115 = int32(8)
	v117 = *(*int32)(unsafe.Add(mBase, uint32(v108+v115)))
	v121 = int32(base.Ui32(v113)>>(uint(v115)%32)) & int32(255)
	v122 = int32(1)
	v125 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v117+v121<<(uint(v122)%32)))))
	v128 = *(*int32)(unsafe.Add(mBase, uint32(v108+int32(4))))
	v130 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v128+v121))))
	if v130 < v122 {
		goto L24
	} else {
		goto L25
	}
L23:
	;
	v283 = *(*int32)(unsafe.Add(mBase, uint32(v108+int32(20))))
	v286 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v80+int32(6)))))
	v287 = int32(1)
	v290 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v283+v286<<(uint(v287)%32)))))
	v293 = *(*int32)(unsafe.Add(mBase, uint32(v108+int32(16))))
	v295 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v293+v286))))
	if v295 < v287 {
		goto L50
	} else {
		goto L51
	}
L24:
	;
	goto L23
L25:
	;
	v142 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v143 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v143+v130 < int32(32) {
		goto L28
	} else {
		goto L29
	}
L26:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v244 + v242
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v241<<(uint(v244)%32) | v243
	goto L24
L27:
	;
	v161 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v162 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v166 = v158
	v167 = v159
	v168 = v162
	v169 = v161
	goto L31
L28:
	;
	if v143 < int32(16) {
		v241 = v125
		v242 = v130
		v243 = v142
		v244 = v143
		goto L26
	} else {
		goto L30
	}
L29:
	;
	v147 = int32(32)
	v148 = v147 - v143
	v156 = int32(base.Ui32(v125) >> (uint(v148) % 32))
	v157 = v130 - v148
	v158 = v125<<(uint(v143)%32) | v142
	v159 = v147
	goto L27
L30:
	;
	v156 = v125
	v157 = v130
	v158 = v142
	v159 = v143
	goto L27
L31:
	;
	if base.Ui32(v168+int32(2)) <= base.Ui32(v169) {
		v224 = v168
		v225 = v169
		goto L33
	} else {
		goto L34
	}
L32:
	;
	v241 = v156
	v242 = v157
	v243 = v237
	v244 = v235
	goto L26
L33:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v224))) = uint16(v166)
	v232 = v224 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v232
	v235 = v167 + int32(-16)
	v237 = int32(base.Ui32(v166) >> (uint(int32(16)) % 32))
	if int32(31) < v167 {
		v166 = v237
		v167 = v235
		v168 = v232
		v169 = v225
		goto L31
	} else {
		goto L48
	}
L34:
	;
	v178 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v179 = v169 - v178
	v182 = base.I64_extend_i32_s(v179) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v182) {
		v206 = v178
		goto L36
	} else {
		goto L37
	}
L35:
	;
	if v168 == v178 {
		goto L46
	} else {
		goto L47
	}
L36:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v206
	goto L23
L37:
	;
	v185 = v168 - v178
	v187 = v182 + base.I64_extend_i32_u(v185)
	if base.Ui64(int64(4294967295)) < base.Ui64(v187) {
		v206 = v178
		goto L36
	} else {
		goto L38
	}
L38:
	;
	v190 = base.I32_wrap_i64(v187)
	if v169 == v178 {
		goto L39
	} else {
		goto L40
	}
L39:
	;
	v197 = int32(base.Ui32(v179*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v190) < base.Ui32(v197) {
		goto L42
	} else {
		goto L43
	}
L40:
	;
	if base.Ui32(v190) <= base.Ui32(v179) {
		v224 = v168
		v225 = v169
		goto L33
	} else {
		goto L41
	}
L41:
	;
	goto L39
L42:
	;
	v199 = v197
	goto L44
L43:
	;
	v199 = v190
	goto L44
L44:
	;
	v203 = v199&int32(-1024) + int32(1024)
	v204 = F_WebPSafeMalloc(m, int64(1), v203)
	mBase = m.M
	if v204 != 0 {
		goto L35
	} else {
		goto L45
	}
L45:
	;
	v205 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v206 = v205
	goto L36
L46:
	;
	v217 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v217)
	mBase = m.M
	v219 = v204 + v203
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v219
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v204
	v224 = v204 + v185
	v225 = v219
	goto L33
L47:
	;
	v215 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v216 = F_memcpy(m, v204, v215, v185)
	mBase = m.M
	goto L46
L48:
	;
	goto L32
L49:
	;
	v448 = *(*int32)(unsafe.Add(mBase, uint32(v108+int32(32))))
	v449 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v112))))
	v450 = int32(1)
	v453 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v448+v449<<(uint(v450)%32)))))
	v456 = *(*int32)(unsafe.Add(mBase, uint32(v108+int32(28))))
	v458 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v456+v449))))
	if v458 < v450 {
		goto L76
	} else {
		goto L77
	}
L50:
	;
	goto L49
L51:
	;
	v307 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v308 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v308+v295 < int32(32) {
		goto L54
	} else {
		goto L55
	}
L52:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v409 + v407
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v406<<(uint(v409)%32) | v408
	goto L50
L53:
	;
	v326 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v327 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v331 = v323
	v332 = v324
	v333 = v327
	v334 = v326
	goto L57
L54:
	;
	if v308 < int32(16) {
		v406 = v290
		v407 = v295
		v408 = v307
		v409 = v308
		goto L52
	} else {
		goto L56
	}
L55:
	;
	v312 = int32(32)
	v313 = v312 - v308
	v321 = int32(base.Ui32(v290) >> (uint(v313) % 32))
	v322 = v295 - v313
	v323 = v290<<(uint(v308)%32) | v307
	v324 = v312
	goto L53
L56:
	;
	v321 = v290
	v322 = v295
	v323 = v307
	v324 = v308
	goto L53
L57:
	;
	if base.Ui32(v333+int32(2)) <= base.Ui32(v334) {
		v389 = v333
		v390 = v334
		goto L59
	} else {
		goto L60
	}
L58:
	;
	v406 = v321
	v407 = v322
	v408 = v402
	v409 = v400
	goto L52
L59:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v389))) = uint16(v331)
	v397 = v389 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v397
	v400 = v332 + int32(-16)
	v402 = int32(base.Ui32(v331) >> (uint(int32(16)) % 32))
	if int32(31) < v332 {
		v331 = v402
		v332 = v400
		v333 = v397
		v334 = v390
		goto L57
	} else {
		goto L74
	}
L60:
	;
	v343 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v344 = v334 - v343
	v347 = base.I64_extend_i32_s(v344) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v347) {
		v371 = v343
		goto L62
	} else {
		goto L63
	}
L61:
	;
	if v333 == v343 {
		goto L72
	} else {
		goto L73
	}
L62:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v371
	goto L49
L63:
	;
	v350 = v333 - v343
	v352 = v347 + base.I64_extend_i32_u(v350)
	if base.Ui64(int64(4294967295)) < base.Ui64(v352) {
		v371 = v343
		goto L62
	} else {
		goto L64
	}
L64:
	;
	v355 = base.I32_wrap_i64(v352)
	if v334 == v343 {
		goto L65
	} else {
		goto L66
	}
L65:
	;
	v362 = int32(base.Ui32(v344*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v355) < base.Ui32(v362) {
		goto L68
	} else {
		goto L69
	}
L66:
	;
	if base.Ui32(v355) <= base.Ui32(v344) {
		v389 = v333
		v390 = v334
		goto L59
	} else {
		goto L67
	}
L67:
	;
	goto L65
L68:
	;
	v364 = v362
	goto L70
L69:
	;
	v364 = v355
	goto L70
L70:
	;
	v368 = v364&int32(-1024) + int32(1024)
	v369 = F_WebPSafeMalloc(m, int64(1), v368)
	mBase = m.M
	if v369 != 0 {
		goto L61
	} else {
		goto L71
	}
L71:
	;
	v370 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v371 = v370
	goto L62
L72:
	;
	v382 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v382)
	mBase = m.M
	v384 = v369 + v368
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v384
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v369
	v389 = v369 + v350
	v390 = v384
	goto L59
L73:
	;
	v380 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v381 = F_memcpy(m, v369, v380, v350)
	mBase = m.M
	goto L72
L74:
	;
	goto L58
L75:
	;
	v611 = *(*int32)(unsafe.Add(mBase, uint32(v108+int32(44))))
	v612 = *(*int32)(unsafe.Add(mBase, uint32(v112)))
	v614 = int32(base.Ui32(v612) >> (uint(int32(24)) % 32))
	v615 = int32(1)
	v618 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v611+v614<<(uint(v615)%32)))))
	v621 = *(*int32)(unsafe.Add(mBase, uint32(v108+int32(40))))
	v623 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v621+v614))))
	if v623 < v615 {
		goto L102
	} else {
		goto L103
	}
L76:
	;
	goto L75
L77:
	;
	v470 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v471 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v471+v458 < int32(32) {
		goto L80
	} else {
		goto L81
	}
L78:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v572 + v570
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v569<<(uint(v572)%32) | v571
	goto L76
L79:
	;
	v489 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v490 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v494 = v486
	v495 = v487
	v496 = v490
	v497 = v489
	goto L83
L80:
	;
	if v471 < int32(16) {
		v569 = v453
		v570 = v458
		v571 = v470
		v572 = v471
		goto L78
	} else {
		goto L82
	}
L81:
	;
	v475 = int32(32)
	v476 = v475 - v471
	v484 = int32(base.Ui32(v453) >> (uint(v476) % 32))
	v485 = v458 - v476
	v486 = v453<<(uint(v471)%32) | v470
	v487 = v475
	goto L79
L82:
	;
	v484 = v453
	v485 = v458
	v486 = v470
	v487 = v471
	goto L79
L83:
	;
	if base.Ui32(v496+int32(2)) <= base.Ui32(v497) {
		v552 = v496
		v553 = v497
		goto L85
	} else {
		goto L86
	}
L84:
	;
	v569 = v484
	v570 = v485
	v571 = v565
	v572 = v563
	goto L78
L85:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v552))) = uint16(v494)
	v560 = v552 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v560
	v563 = v495 + int32(-16)
	v565 = int32(base.Ui32(v494) >> (uint(int32(16)) % 32))
	if int32(31) < v495 {
		v494 = v565
		v495 = v563
		v496 = v560
		v497 = v553
		goto L83
	} else {
		goto L100
	}
L86:
	;
	v506 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v507 = v497 - v506
	v510 = base.I64_extend_i32_s(v507) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v510) {
		v534 = v506
		goto L88
	} else {
		goto L89
	}
L87:
	;
	if v496 == v506 {
		goto L98
	} else {
		goto L99
	}
L88:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v534
	goto L75
L89:
	;
	v513 = v496 - v506
	v515 = v510 + base.I64_extend_i32_u(v513)
	if base.Ui64(int64(4294967295)) < base.Ui64(v515) {
		v534 = v506
		goto L88
	} else {
		goto L90
	}
L90:
	;
	v518 = base.I32_wrap_i64(v515)
	if v497 == v506 {
		goto L91
	} else {
		goto L92
	}
L91:
	;
	v525 = int32(base.Ui32(v507*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v518) < base.Ui32(v525) {
		goto L94
	} else {
		goto L95
	}
L92:
	;
	if base.Ui32(v518) <= base.Ui32(v507) {
		v552 = v496
		v553 = v497
		goto L85
	} else {
		goto L93
	}
L93:
	;
	goto L91
L94:
	;
	v527 = v525
	goto L96
L95:
	;
	v527 = v518
	goto L96
L96:
	;
	v531 = v527&int32(-1024) + int32(1024)
	v532 = F_WebPSafeMalloc(m, int64(1), v531)
	mBase = m.M
	if v532 != 0 {
		goto L87
	} else {
		goto L97
	}
L97:
	;
	v533 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v534 = v533
	goto L88
L98:
	;
	v545 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v545)
	mBase = m.M
	v547 = v532 + v531
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v547
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v532
	v552 = v532 + v513
	v553 = v547
	goto L85
L99:
	;
	v543 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v544 = F_memcpy(m, v532, v543, v513)
	mBase = m.M
	goto L98
L100:
	;
	goto L84
L101:
	;
	goto L19
L102:
	;
	goto L101
L103:
	;
	v635 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v636 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v636+v623 < int32(32) {
		goto L106
	} else {
		goto L107
	}
L104:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v737 + v735
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v734<<(uint(v737)%32) | v736
	goto L102
L105:
	;
	v654 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v655 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v659 = v651
	v660 = v652
	v661 = v655
	v662 = v654
	goto L109
L106:
	;
	if v636 < int32(16) {
		v734 = v618
		v735 = v623
		v736 = v635
		v737 = v636
		goto L104
	} else {
		goto L108
	}
L107:
	;
	v640 = int32(32)
	v641 = v640 - v636
	v649 = int32(base.Ui32(v618) >> (uint(v641) % 32))
	v650 = v623 - v641
	v651 = v618<<(uint(v636)%32) | v635
	v652 = v640
	goto L105
L108:
	;
	v649 = v618
	v650 = v623
	v651 = v635
	v652 = v636
	goto L105
L109:
	;
	if base.Ui32(v661+int32(2)) <= base.Ui32(v662) {
		v717 = v661
		v718 = v662
		goto L111
	} else {
		goto L112
	}
L110:
	;
	v734 = v649
	v735 = v650
	v736 = v730
	v737 = v728
	goto L104
L111:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v717))) = uint16(v659)
	v725 = v717 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v725
	v728 = v660 + int32(-16)
	v730 = int32(base.Ui32(v659) >> (uint(int32(16)) % 32))
	if int32(31) < v660 {
		v659 = v730
		v660 = v728
		v661 = v725
		v662 = v718
		goto L109
	} else {
		goto L126
	}
L112:
	;
	v671 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v672 = v662 - v671
	v675 = base.I64_extend_i32_s(v672) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v675) {
		v699 = v671
		goto L114
	} else {
		goto L115
	}
L113:
	;
	if v661 == v671 {
		goto L124
	} else {
		goto L125
	}
L114:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v699
	goto L101
L115:
	;
	v678 = v661 - v671
	v680 = v675 + base.I64_extend_i32_u(v678)
	if base.Ui64(int64(4294967295)) < base.Ui64(v680) {
		v699 = v671
		goto L114
	} else {
		goto L116
	}
L116:
	;
	v683 = base.I32_wrap_i64(v680)
	if v662 == v671 {
		goto L117
	} else {
		goto L118
	}
L117:
	;
	v690 = int32(base.Ui32(v672*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v683) < base.Ui32(v690) {
		goto L120
	} else {
		goto L121
	}
L118:
	;
	if base.Ui32(v683) <= base.Ui32(v672) {
		v717 = v661
		v718 = v662
		goto L111
	} else {
		goto L119
	}
L119:
	;
	goto L117
L120:
	;
	v692 = v690
	goto L122
L121:
	;
	v692 = v683
	goto L122
L122:
	;
	v696 = v692&int32(-1024) + int32(1024)
	v697 = F_WebPSafeMalloc(m, int64(1), v696)
	mBase = m.M
	if v697 != 0 {
		goto L113
	} else {
		goto L123
	}
L123:
	;
	v698 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v699 = v698
	goto L114
L124:
	;
	v710 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v710)
	mBase = m.M
	v712 = v697 + v696
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v712
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v697
	v717 = v697 + v678
	v718 = v712
	goto L111
L125:
	;
	v708 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v709 = F_memcpy(m, v697, v708, v678)
	mBase = m.M
	goto L124
L126:
	;
	goto L110
L127:
	;
	goto L19
L128:
	;
	goto L127
L129:
	;
	v799 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v800 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v800+v787 < int32(32) {
		goto L132
	} else {
		goto L133
	}
L130:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v901 + v899
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v898<<(uint(v901)%32) | v900
	goto L128
L131:
	;
	v818 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v819 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v823 = v815
	v824 = v816
	v825 = v819
	v826 = v818
	goto L135
L132:
	;
	if v800 < int32(16) {
		v898 = v782
		v899 = v787
		v900 = v799
		v901 = v800
		goto L130
	} else {
		goto L134
	}
L133:
	;
	v804 = int32(32)
	v805 = v804 - v800
	v813 = int32(base.Ui32(v782) >> (uint(v805) % 32))
	v814 = v787 - v805
	v815 = v782<<(uint(v800)%32) | v799
	v816 = v804
	goto L131
L134:
	;
	v813 = v782
	v814 = v787
	v815 = v799
	v816 = v800
	goto L131
L135:
	;
	if base.Ui32(v825+int32(2)) <= base.Ui32(v826) {
		v881 = v825
		v882 = v826
		goto L137
	} else {
		goto L138
	}
L136:
	;
	v898 = v813
	v899 = v814
	v900 = v894
	v901 = v892
	goto L130
L137:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v881))) = uint16(v823)
	v889 = v881 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v889
	v892 = v824 + int32(-16)
	v894 = int32(base.Ui32(v823) >> (uint(int32(16)) % 32))
	if int32(31) < v824 {
		v823 = v894
		v824 = v892
		v825 = v889
		v826 = v882
		goto L135
	} else {
		goto L152
	}
L138:
	;
	v835 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v836 = v826 - v835
	v839 = base.I64_extend_i32_s(v836) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v839) {
		v863 = v835
		goto L140
	} else {
		goto L141
	}
L139:
	;
	if v825 == v835 {
		goto L150
	} else {
		goto L151
	}
L140:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v863
	goto L127
L141:
	;
	v842 = v825 - v835
	v844 = v839 + base.I64_extend_i32_u(v842)
	if base.Ui64(int64(4294967295)) < base.Ui64(v844) {
		v863 = v835
		goto L140
	} else {
		goto L142
	}
L142:
	;
	v847 = base.I32_wrap_i64(v844)
	if v826 == v835 {
		goto L143
	} else {
		goto L144
	}
L143:
	;
	v854 = int32(base.Ui32(v836*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v847) < base.Ui32(v854) {
		goto L146
	} else {
		goto L147
	}
L144:
	;
	if base.Ui32(v847) <= base.Ui32(v836) {
		v881 = v825
		v882 = v826
		goto L137
	} else {
		goto L145
	}
L145:
	;
	goto L143
L146:
	;
	v856 = v854
	goto L148
L147:
	;
	v856 = v847
	goto L148
L148:
	;
	v860 = v856&int32(-1024) + int32(1024)
	v861 = F_WebPSafeMalloc(m, int64(1), v860)
	mBase = m.M
	if v861 != 0 {
		goto L139
	} else {
		goto L149
	}
L149:
	;
	v862 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v863 = v862
	goto L140
L150:
	;
	v874 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v874)
	mBase = m.M
	v876 = v861 + v860
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v876
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v861
	v881 = v861 + v842
	v882 = v876
	goto L137
L151:
	;
	v872 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v873 = F_memcpy(m, v861, v872, v842)
	mBase = m.M
	goto L150
L152:
	;
	goto L136
L153:
	;
	v973 = *(*int32)(unsafe.Add(mBase, uint32(v108+int32(4))))
	v975 = v970 + int32(256)
	v977 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v973+v975))))
	v981 = *(*int32)(unsafe.Add(mBase, uint32(v108+int32(8))))
	v982 = int32(1)
	v985 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v981+v975<<(uint(v982)%32)))))
	v986 = v969<<(uint(v977)%32) | v985
	v987 = v968 + v977
	if v987 < v982 {
		goto L157
	} else {
		goto L158
	}
L154:
	;
	v950 = int32(-1)
	v951 = v938 + v950
	v954 = base.I32_clz(v951) ^ int32(31)
	v956 = v954 + v950
	v958 = int32(1)
	v968 = v956
	v969 = v951 & (v950<<(uint(v956)%32) ^ v950)
	v970 = int32(base.Ui32(v951)>>(uint(v956)%32))&v958 | v954<<(uint(v958)%32)
	goto L153
L155:
	;
	v941 = m.G73
	v943 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v941+v938))))
	v944 = m.G54
	v947 = v944 + v938<<(uint(int32(1))%32)
	v948 = int32(*(*int8)(unsafe.Add(mBase, uint32(v947)+1)))
	v949 = int32(*(*int8)(unsafe.Add(mBase, uint32(v947))))
	v968 = v948
	v969 = v943
	v970 = v949
	goto L153
L156:
	;
	if int32(511) < v113 {
		goto L183
	} else {
		goto L184
	}
L157:
	;
	goto L156
L158:
	;
	v999 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v1000 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v1000+v987 < int32(32) {
		goto L161
	} else {
		goto L162
	}
L159:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1101 + v1099
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1098<<(uint(v1101)%32) | v1100
	goto L157
L160:
	;
	v1018 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v1019 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v1023 = v1015
	v1024 = v1016
	v1025 = v1019
	v1026 = v1018
	goto L164
L161:
	;
	if v1000 < int32(16) {
		v1098 = v986
		v1099 = v987
		v1100 = v999
		v1101 = v1000
		goto L159
	} else {
		goto L163
	}
L162:
	;
	v1004 = int32(32)
	v1005 = v1004 - v1000
	v1013 = int32(base.Ui32(v986) >> (uint(v1005) % 32))
	v1014 = v987 - v1005
	v1015 = v986<<(uint(v1000)%32) | v999
	v1016 = v1004
	goto L160
L163:
	;
	v1013 = v986
	v1014 = v987
	v1015 = v999
	v1016 = v1000
	goto L160
L164:
	;
	if base.Ui32(v1025+int32(2)) <= base.Ui32(v1026) {
		v1081 = v1025
		v1082 = v1026
		goto L166
	} else {
		goto L167
	}
L165:
	;
	v1098 = v1013
	v1099 = v1014
	v1100 = v1094
	v1101 = v1092
	goto L159
L166:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1081))) = uint16(v1023)
	v1089 = v1081 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1089
	v1092 = v1024 + int32(-16)
	v1094 = int32(base.Ui32(v1023) >> (uint(int32(16)) % 32))
	if int32(31) < v1024 {
		v1023 = v1094
		v1024 = v1092
		v1025 = v1089
		v1026 = v1082
		goto L164
	} else {
		goto L181
	}
L167:
	;
	v1035 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1036 = v1026 - v1035
	v1039 = base.I64_extend_i32_s(v1036) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1039) {
		v1063 = v1035
		goto L169
	} else {
		goto L170
	}
L168:
	;
	if v1025 == v1035 {
		goto L179
	} else {
		goto L180
	}
L169:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1063
	goto L156
L170:
	;
	v1042 = v1025 - v1035
	v1044 = v1039 + base.I64_extend_i32_u(v1042)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1044) {
		v1063 = v1035
		goto L169
	} else {
		goto L171
	}
L171:
	;
	v1047 = base.I32_wrap_i64(v1044)
	if v1026 == v1035 {
		goto L172
	} else {
		goto L173
	}
L172:
	;
	v1054 = int32(base.Ui32(v1036*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v1047) < base.Ui32(v1054) {
		goto L175
	} else {
		goto L176
	}
L173:
	;
	if base.Ui32(v1047) <= base.Ui32(v1036) {
		v1081 = v1025
		v1082 = v1026
		goto L166
	} else {
		goto L174
	}
L174:
	;
	goto L172
L175:
	;
	v1056 = v1054
	goto L177
L176:
	;
	v1056 = v1047
	goto L177
L177:
	;
	v1060 = v1056&int32(-1024) + int32(1024)
	v1061 = F_WebPSafeMalloc(m, int64(1), v1060)
	mBase = m.M
	if v1061 != 0 {
		goto L168
	} else {
		goto L178
	}
L178:
	;
	v1062 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1063 = v1062
	goto L169
L179:
	;
	v1074 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v1074)
	mBase = m.M
	v1076 = v1061 + v1060
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v1076
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v1061
	v1081 = v1061 + v1042
	v1082 = v1076
	goto L166
L180:
	;
	v1072 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1073 = F_memcpy(m, v1061, v1072, v1042)
	mBase = m.M
	goto L179
L181:
	;
	goto L165
L182:
	;
	v1172 = *(*int32)(unsafe.Add(mBase, uint32(v108+int32(56))))
	v1173 = int32(1)
	v1176 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1172+v1167<<(uint(v1173)%32)))))
	v1179 = *(*int32)(unsafe.Add(mBase, uint32(v108+int32(52))))
	v1181 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1179+v1167))))
	if v1181 < v1173 {
		goto L186
	} else {
		goto L187
	}
L183:
	;
	v1149 = int32(-1)
	v1150 = v113 + v1149
	v1153 = base.I32_clz(v1150) ^ int32(31)
	v1155 = v1153 + v1149
	v1157 = int32(1)
	v1167 = int32(base.Ui32(v1150)>>(uint(v1155)%32))&v1157 | v1153<<(uint(v1157)%32)
	v1168 = v1155
	v1169 = v1150 & (v1149<<(uint(v1155)%32) ^ v1149)
	goto L182
L184:
	;
	v1140 = m.G73
	v1142 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1140+v113))))
	v1143 = m.G54
	v1146 = v1143 + v113<<(uint(int32(1))%32)
	v1147 = int32(*(*int8)(unsafe.Add(mBase, uint32(v1146)+1)))
	v1148 = int32(*(*int8)(unsafe.Add(mBase, uint32(v1146))))
	v1167 = v1148
	v1168 = v1147
	v1169 = v1142
	goto L182
L185:
	;
	if v1168 < int32(1) {
		goto L212
	} else {
		goto L213
	}
L186:
	;
	goto L185
L187:
	;
	v1193 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v1194 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v1194+v1181 < int32(32) {
		goto L190
	} else {
		goto L191
	}
L188:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1295 + v1293
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1292<<(uint(v1295)%32) | v1294
	goto L186
L189:
	;
	v1212 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v1213 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v1217 = v1209
	v1218 = v1210
	v1219 = v1213
	v1220 = v1212
	goto L193
L190:
	;
	if v1194 < int32(16) {
		v1292 = v1176
		v1293 = v1181
		v1294 = v1193
		v1295 = v1194
		goto L188
	} else {
		goto L192
	}
L191:
	;
	v1198 = int32(32)
	v1199 = v1198 - v1194
	v1207 = int32(base.Ui32(v1176) >> (uint(v1199) % 32))
	v1208 = v1181 - v1199
	v1209 = v1176<<(uint(v1194)%32) | v1193
	v1210 = v1198
	goto L189
L192:
	;
	v1207 = v1176
	v1208 = v1181
	v1209 = v1193
	v1210 = v1194
	goto L189
L193:
	;
	if base.Ui32(v1219+int32(2)) <= base.Ui32(v1220) {
		v1275 = v1219
		v1276 = v1220
		goto L195
	} else {
		goto L196
	}
L194:
	;
	v1292 = v1207
	v1293 = v1208
	v1294 = v1288
	v1295 = v1286
	goto L188
L195:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1275))) = uint16(v1217)
	v1283 = v1275 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1283
	v1286 = v1218 + int32(-16)
	v1288 = int32(base.Ui32(v1217) >> (uint(int32(16)) % 32))
	if int32(31) < v1218 {
		v1217 = v1288
		v1218 = v1286
		v1219 = v1283
		v1220 = v1276
		goto L193
	} else {
		goto L210
	}
L196:
	;
	v1229 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1230 = v1220 - v1229
	v1233 = base.I64_extend_i32_s(v1230) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1233) {
		v1257 = v1229
		goto L198
	} else {
		goto L199
	}
L197:
	;
	if v1219 == v1229 {
		goto L208
	} else {
		goto L209
	}
L198:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1257
	goto L185
L199:
	;
	v1236 = v1219 - v1229
	v1238 = v1233 + base.I64_extend_i32_u(v1236)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1238) {
		v1257 = v1229
		goto L198
	} else {
		goto L200
	}
L200:
	;
	v1241 = base.I32_wrap_i64(v1238)
	if v1220 == v1229 {
		goto L201
	} else {
		goto L202
	}
L201:
	;
	v1248 = int32(base.Ui32(v1230*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v1241) < base.Ui32(v1248) {
		goto L204
	} else {
		goto L205
	}
L202:
	;
	if base.Ui32(v1241) <= base.Ui32(v1230) {
		v1275 = v1219
		v1276 = v1220
		goto L195
	} else {
		goto L203
	}
L203:
	;
	goto L201
L204:
	;
	v1250 = v1248
	goto L206
L205:
	;
	v1250 = v1241
	goto L206
L206:
	;
	v1254 = v1250&int32(-1024) + int32(1024)
	v1255 = F_WebPSafeMalloc(m, int64(1), v1254)
	mBase = m.M
	if v1255 != 0 {
		goto L197
	} else {
		goto L207
	}
L207:
	;
	v1256 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1257 = v1256
	goto L198
L208:
	;
	v1268 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v1268)
	mBase = m.M
	v1270 = v1255 + v1254
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v1270
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v1255
	v1275 = v1255 + v1236
	v1276 = v1270
	goto L195
L209:
	;
	v1266 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1267 = F_memcpy(m, v1255, v1266, v1236)
	mBase = m.M
	goto L208
L210:
	;
	goto L194
L211:
	;
	goto L19
L212:
	;
	goto L211
L213:
	;
	v1343 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v1344 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v1344+v1168 < int32(32) {
		goto L216
	} else {
		goto L217
	}
L214:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1445 + v1443
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1442<<(uint(v1445)%32) | v1444
	goto L212
L215:
	;
	v1362 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v1363 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v1367 = v1359
	v1368 = v1360
	v1369 = v1363
	v1370 = v1362
	goto L219
L216:
	;
	if v1344 < int32(16) {
		v1442 = v1169
		v1443 = v1168
		v1444 = v1343
		v1445 = v1344
		goto L214
	} else {
		goto L218
	}
L217:
	;
	v1348 = int32(32)
	v1349 = v1348 - v1344
	v1357 = int32(base.Ui32(v1169) >> (uint(v1349) % 32))
	v1358 = v1168 - v1349
	v1359 = v1169<<(uint(v1344)%32) | v1343
	v1360 = v1348
	goto L215
L218:
	;
	v1357 = v1169
	v1358 = v1168
	v1359 = v1343
	v1360 = v1344
	goto L215
L219:
	;
	if base.Ui32(v1369+int32(2)) <= base.Ui32(v1370) {
		v1425 = v1369
		v1426 = v1370
		goto L221
	} else {
		goto L222
	}
L220:
	;
	v1442 = v1357
	v1443 = v1358
	v1444 = v1438
	v1445 = v1436
	goto L214
L221:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1425))) = uint16(v1367)
	v1433 = v1425 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1433
	v1436 = v1368 + int32(-16)
	v1438 = int32(base.Ui32(v1367) >> (uint(int32(16)) % 32))
	if int32(31) < v1368 {
		v1367 = v1438
		v1368 = v1436
		v1369 = v1433
		v1370 = v1426
		goto L219
	} else {
		goto L236
	}
L222:
	;
	v1379 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1380 = v1370 - v1379
	v1383 = base.I64_extend_i32_s(v1380) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1383) {
		v1407 = v1379
		goto L224
	} else {
		goto L225
	}
L223:
	;
	if v1369 == v1379 {
		goto L234
	} else {
		goto L235
	}
L224:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1407
	goto L211
L225:
	;
	v1386 = v1369 - v1379
	v1388 = v1383 + base.I64_extend_i32_u(v1386)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1388) {
		v1407 = v1379
		goto L224
	} else {
		goto L226
	}
L226:
	;
	v1391 = base.I32_wrap_i64(v1388)
	if v1370 == v1379 {
		goto L227
	} else {
		goto L228
	}
L227:
	;
	v1398 = int32(base.Ui32(v1380*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v1391) < base.Ui32(v1398) {
		goto L230
	} else {
		goto L231
	}
L228:
	;
	if base.Ui32(v1391) <= base.Ui32(v1380) {
		v1425 = v1369
		v1426 = v1370
		goto L221
	} else {
		goto L229
	}
L229:
	;
	goto L227
L230:
	;
	v1400 = v1398
	goto L232
L231:
	;
	v1400 = v1391
	goto L232
L232:
	;
	v1404 = v1400&int32(-1024) + int32(1024)
	v1405 = F_WebPSafeMalloc(m, int64(1), v1404)
	mBase = m.M
	if v1405 != 0 {
		goto L223
	} else {
		goto L233
	}
L233:
	;
	v1406 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1407 = v1406
	goto L224
L234:
	;
	v1418 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v1418)
	mBase = m.M
	v1420 = v1405 + v1404
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v1420
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v1405
	v1425 = v1405 + v1386
	v1426 = v1420
	goto L221
L235:
	;
	v1416 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1417 = F_memcpy(m, v1405, v1416, v1386)
	mBase = m.M
	goto L234
L236:
	;
	goto L220
L237:
	;
	v1538 = *(*int32)(unsafe.Add(mBase, uint32(v24)+4))
	v1540 = v1538 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v24)+4)) = v1540
	v1542 = *(*int32)(unsafe.Add(mBase, uint32(v24)+12))
	if v1540 != v1542 {
		v1564 = v1540
		goto L242
	} else {
		goto L243
	}
L238:
	;
	v1495 = v1490
	v1500 = v78
	goto L239
L239:
	;
	v1514 = v1500 + int32(1)
	v1515 = v1495 - l1
	if l1 <= v1515 {
		v1495 = v1515
		v1500 = v1514
		goto L239
	} else {
		goto L241
	}
L240:
	;
	v1520 = v1515
	v1525 = v1514
	goto L237
L241:
	;
	goto L240
L242:
	;
	if v1564 != 0 {
		v73 = v1520
		v78 = v1525
		v80 = v1564
		v83 = v108
		v84 = v109
		v85 = v110
		goto L13
	} else {
		goto L248
	}
L243:
	;
	v1545 = v24 + int32(4)
	v1549 = *(*int32)(unsafe.Add(mBase, uint32(v1545)+4))
	v1550 = *(*int32)(unsafe.Add(mBase, uint32(v1549)))
	if v1550 != 0 {
		goto L246
	} else {
		goto L247
	}
L244:
	;
	v1563 = *(*int32)(unsafe.Add(mBase, uint32(v24)+4))
	v1564 = v1563
	goto L242
L245:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1545)+8)) = v1559
	*(*int32)(unsafe.Add(mBase, uint32(v1545))) = v1558
	*(*int32)(unsafe.Add(mBase, uint32(v1545)+4)) = v1550
	goto L244
L246:
	;
	v1553 = *(*int32)(unsafe.Add(mBase, uint32(v1550)+4))
	v1554 = *(*int32)(unsafe.Add(mBase, uint32(v1550)+8))
	v1558 = v1553
	v1559 = v1553 + v1554<<(uint(int32(3))%32)
	goto L245
L247:
	;
	v1551 = int32(0)
	v1558 = v1551
	v1559 = v1551
	goto L245
L248:
	;
	goto L14
L249:
	;
	m.G0 = v24 + int32(16)
	return v1593
L250:
	;
	v1590 = *(*int32)(unsafe.Add(mBase, uint32(l6)+92))
	if v1590 != 0 {
		goto L252
	} else {
		goto L253
	}
L251:
	;
	v1593 = int32(0)
	goto L249
L252:
	;
	goto L251
L253:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l6)+92)) = int32(1)
	goto L252
}
func F_StoreSideInfo(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v20 int64
	_ = v20
	var v24 int32
	_ = v24
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int64
	_ = v31
	var v35 int32
	_ = v35
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v41 int64
	_ = v41
	var v45 int64
	_ = v45
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v61 int32
	_ = v61
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v77 int32
	_ = v77
	var v80 int32
	_ = v80
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v87 int32
	_ = v87
	var v90 int32
	_ = v90
	var v93 int32
	_ = v93
	var v98 int32
	_ = v98
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v115 int32
	_ = v115
	var v116 int32
	_ = v116
	var v117 int32
	_ = v117
	var v122 int64
	_ = v122
	var v123 int64
	_ = v123
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v133 int32
	_ = v133
	var v134 int32
	_ = v134
	var v136 int32
	_ = v136
	v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
	v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v11 = *(*int32)(unsafe.Add(mBase, uint32(v10)+4))
	v12 = *(*int32)(unsafe.Add(mBase, uint32(v11)+88))
	if v12 == int32(0) {
	} else {
		v15 = m.G29
		v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
		v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
		v18 = *(*int32)(unsafe.Add(mBase, uint32(v15)))
		v19 = m.T0[v18].(func(*base.Module, int32, int32) int32)(m, v16, v17)
		mBase = m.M
		v20 = *(*int64)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_StoreSideInfo[0])))
		*(*int64)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_StoreSideInfo[0]))) = v20 + base.I64_extend_i32_s(v19)
		v24 = int32(16)
		v28 = m.G82
		v29 = *(*int32)(unsafe.Add(mBase, uint32(v28)))
		v30 = m.T0[v29].(func(*base.Module, int32, int32) int32)(m, v16+v24, v17+v24)
		mBase = m.M
		v31 = *(*int64)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_StoreSideInfo[1])))
		*(*int64)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_StoreSideInfo[1]))) = v31 + base.I64_extend_i32_s(v30)
		v35 = int32(24)
		v39 = *(*int32)(unsafe.Add(mBase, uint32(v28)))
		v40 = m.T0[v39].(func(*base.Module, int32, int32) int32)(m, v16+v35, v17+v35)
		mBase = m.M
		v41 = *(*int64)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_StoreSideInfo[2])))
		*(*int64)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_StoreSideInfo[2]))) = v41 + base.I64_extend_i32_s(v40)
		v45 = *(*int64)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_StoreSideInfo[3])))
		*(*int64)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_StoreSideInfo[3]))) = v45 + int64(256)
		v49 = *(*int32)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_StoreSideInfo[4])))
		v50 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v9))))
		v51 = int32(3)
		*(*int32)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_StoreSideInfo[4]))) = v49 + base.B2i32(v50&v51 == int32(0))
		v57 = *(*int32)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_StoreSideInfo[5])))
		v58 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v9))))
		v61 = int32(1)
		*(*int32)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_StoreSideInfo[5]))) = v57 + base.B2i32(v58&v51 == v61)
		v65 = *(*int32)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_StoreSideInfo[6])))
		v66 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v9))))
		*(*int32)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_StoreSideInfo[6]))) = v65 + int32(base.Ui32(v66)>>(uint(int32(4))%32))&v61
	}
	v77 = *(*int32)(unsafe.Add(mBase, uint32(v11)+84))
	if v77 == int32(0) {
	} else {
		v80 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		v82 = *(*int32)(unsafe.Add(mBase, uint32(v10)+40))
		v83 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
		v87 = *(*int32)(unsafe.Add(mBase, uint32(v11)+80))
		switch v87 + int32(-1) {
		case 0:
			v90 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v9))))
			v136 = v90 & int32(3)
		case 1:
			v93 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v9))))
			v136 = int32(base.Ui32(v93)>>(uint(int32(5))%32)) & int32(3)
		case 2:
			v98 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v9))))
			v108 = *(*int32)(unsafe.Add(mBase, uint32(v10+int32(base.Ui32(v98)>>(uint(int32(5))%32))&int32(3)*int32(744)+int32(1088))))
			v136 = v108
		case 3:
			v109 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v9))))
			if v109&int32(3) == int32(1) {
				v115 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
				v116 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v115))))
				v136 = v116
			} else {
				v136 = int32(255)
			}
		case 4:
			v117 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v9))))
			v136 = int32(base.Ui32(v117)>>(uint(int32(2))%32)) & int32(3)
		case 5:
			v122 = *(*int64)(unsafe.Add(mBase, uint32(l0)+264))
			v123 = *(*int64)(unsafe.Add(mBase, uint32(l0)+272))
			v129 = base.I32_wrap_i64(int64(base.Ui64(v122+v123+int64(7)) >> (uint(int64(3)) % 64)))
			v130 = int32(255)
			if v129 < v130 {
				v133 = v129
			} else {
				v133 = v130
			}
			v136 = v133
		case 6:
			v134 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v9)+1)))
			v136 = v134
		default:
			v136 = int32(0)
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v77+v80+v82*v83))) = uint8(v136)
	}
	return
}
func F_StoreStats(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v16 int32
	_ = v16
	var v18 int32
	_ = v18
	var v20 int32
	_ = v20
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
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
	var v36 int32
	_ = v36
	var v38 int32
	_ = v38
	var v40 int32
	_ = v40
	var v42 int32
	_ = v42
	var v44 int32
	_ = v44
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v50 int32
	_ = v50
	var v52 int32
	_ = v52
	var v54 int32
	_ = v54
	var v56 int64
	_ = v56
	var v57 float32
	_ = v57
	var v59 int64
	_ = v59
	var v61 int32
	_ = v61
	var v69 float64
	_ = v69
	var v80 int64
	_ = v80
	var v95 int32
	_ = v95
	var v97 int64
	_ = v97
	var v106 int64
	_ = v106
	var v111 int64
	_ = v111
	var v112 int32
	_ = v112
	var v114 int32
	_ = v114
	var v116 int32
	_ = v116
	var v120 float64
	_ = v120
	var v122 float64
	_ = v122
	var v135 float64
	_ = v135
	var v138 float64
	_ = v138
	var v143 float64
	_ = v143
	var v144 float64
	_ = v144
	var v145 float64
	_ = v145
	var v146 float64
	_ = v146
	var v151 float64
	_ = v151
	var v152 float64
	_ = v152
	var v153 float64
	_ = v153
	var v178 float64
	_ = v178
	var v190 float64
	_ = v190
	var v212 float64
	_ = v212
	var v216 float32
	_ = v216
	var v219 int64
	_ = v219
	var v220 int64
	_ = v220
	var v222 int32
	_ = v222
	var v229 float64
	_ = v229
	var v240 int64
	_ = v240
	var v255 int32
	_ = v255
	var v257 int64
	_ = v257
	var v266 int64
	_ = v266
	var v271 int64
	_ = v271
	var v272 int32
	_ = v272
	var v274 int32
	_ = v274
	var v276 int32
	_ = v276
	var v280 float64
	_ = v280
	var v282 float64
	_ = v282
	var v295 float64
	_ = v295
	var v298 float64
	_ = v298
	var v303 float64
	_ = v303
	var v304 float64
	_ = v304
	var v305 float64
	_ = v305
	var v306 float64
	_ = v306
	var v311 float64
	_ = v311
	var v312 float64
	_ = v312
	var v313 float64
	_ = v313
	var v338 float64
	_ = v338
	var v350 float64
	_ = v350
	var v372 float64
	_ = v372
	var v376 float32
	_ = v376
	var v378 int64
	_ = v378
	var v379 float32
	_ = v379
	var v388 float64
	_ = v388
	var v399 int64
	_ = v399
	var v414 int32
	_ = v414
	var v416 int64
	_ = v416
	var v425 int64
	_ = v425
	var v430 int64
	_ = v430
	var v431 int32
	_ = v431
	var v433 int32
	_ = v433
	var v435 int32
	_ = v435
	var v439 float64
	_ = v439
	var v441 float64
	_ = v441
	var v454 float64
	_ = v454
	var v457 float64
	_ = v457
	var v462 float64
	_ = v462
	var v463 float64
	_ = v463
	var v464 float64
	_ = v464
	var v465 float64
	_ = v465
	var v470 float64
	_ = v470
	var v471 float64
	_ = v471
	var v472 float64
	_ = v472
	var v497 float64
	_ = v497
	var v509 float64
	_ = v509
	var v531 float64
	_ = v531
	var v535 float32
	_ = v535
	var v538 int64
	_ = v538
	var v542 int64
	_ = v542
	var v551 float64
	_ = v551
	var v562 int64
	_ = v562
	var v577 int32
	_ = v577
	var v579 int64
	_ = v579
	var v588 int64
	_ = v588
	var v593 int64
	_ = v593
	var v594 int32
	_ = v594
	var v596 int32
	_ = v596
	var v598 int32
	_ = v598
	var v602 float64
	_ = v602
	var v604 float64
	_ = v604
	var v617 float64
	_ = v617
	var v620 float64
	_ = v620
	var v625 float64
	_ = v625
	var v626 float64
	_ = v626
	var v627 float64
	_ = v627
	var v628 float64
	_ = v628
	var v633 float64
	_ = v633
	var v634 float64
	_ = v634
	var v635 float64
	_ = v635
	var v660 float64
	_ = v660
	var v672 float64
	_ = v672
	var v694 float64
	_ = v694
	var v699 float32
	_ = v699
	var v701 float32
	_ = v701
	var v702 int64
	_ = v702
	var v709 float64
	_ = v709
	var v720 int64
	_ = v720
	var v735 int32
	_ = v735
	var v737 int64
	_ = v737
	var v746 int64
	_ = v746
	var v751 int64
	_ = v751
	var v752 int32
	_ = v752
	var v754 int32
	_ = v754
	var v756 int32
	_ = v756
	var v760 float64
	_ = v760
	var v762 float64
	_ = v762
	var v775 float64
	_ = v775
	var v778 float64
	_ = v778
	var v783 float64
	_ = v783
	var v784 float64
	_ = v784
	var v785 float64
	_ = v785
	var v786 float64
	_ = v786
	var v791 float64
	_ = v791
	var v792 float64
	_ = v792
	var v793 float64
	_ = v793
	var v818 float64
	_ = v818
	var v830 float64
	_ = v830
	var v852 float64
	_ = v852
	var v857 float32
	_ = v857
	var v858 int32
	_ = v858
	var v860 int32
	_ = v860
	var v862 int32
	_ = v862
	var v865 int32
	_ = v865
	v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v13 = *(*int32)(unsafe.Add(mBase, uint32(v12)+88))
	if v13 == int32(0) {
	} else {
		v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1092))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+124)) = v16
		v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1088))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+108)) = v18
		v20 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[0])))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+44)) = v20
		v22 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[1])))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+60)) = v22
		v24 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[2])))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+76)) = v24
		v26 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1836))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+128)) = v26
		v28 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1832))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+112)) = v28
		v30 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[3])))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+48)) = v30
		v32 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[4])))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+64)) = v32
		v34 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[5])))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+80)) = v34
		v36 = *(*int32)(unsafe.Add(mBase, uint32(l0)+2580))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+132)) = v36
		v38 = *(*int32)(unsafe.Add(mBase, uint32(l0)+2576))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+116)) = v38
		v40 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[6])))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+52)) = v40
		v42 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[7])))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+68)) = v42
		v44 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[8])))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+84)) = v44
		v46 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3324))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+136)) = v46
		v48 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3320))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+120)) = v48
		v50 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[9])))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+56)) = v50
		v52 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[10])))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+72)) = v52
		v54 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[11])))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+88)) = v54
		v56 = *(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[12])))
		v57 = float32(99)
		v59 = *(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[13])))
		v61 = base.B2i32(v59 == int64(0))
		if v59 == int64(0) {
			v216 = v57
		} else {
			if v56 == int64(0) {
				v216 = float32(99)
			} else {
				v69 = base.F64_div(base.F64_mul(base.F64_convert_i64_u(v59), float64(65025)), base.F64_convert_i64_u(v56))
				v80 = base.I64_reinterpret_f64(v69)
				if int64(4503599627370495) < v80 {
					if base.Ui64(int64(9218868437227405311)) < base.Ui64(v80) {
						v190 = v69
						v212 = v190
					} else {
						v95 = int32(-1023)
						v97 = int64(base.Ui64(v80) >> (uint(int64(32)) % 64))
						if v97 == int64(1072693248) {
							if base.I32_wrap_i64(v80) != 0 {
								v111 = v80
								v112 = v95
								v114 = int32(1072693248)
								v116 = v114 + int32(614242)
								v120 = base.F64_convert_i32_s(v112 + int32(base.Ui32(v116)>>(uint(int32(20))%32)))
								v122 = base.F64_mul(v120, float64(0.30102999566361177))
								v135 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v116&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v111&int64(4294967295)), float64(-1))
								v138 = base.F64_mul(v135, base.F64_mul(v135, float64(0.5)))
								v143 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v135, v138)) & int64(-4294967296))
								v144 = float64(0.4342944818781689)
								v145 = base.F64_mul(v143, v144)
								v146 = base.F64_add(v122, v145)
								v151 = base.F64_div(v135, base.F64_add(v135, float64(2)))
								v152 = base.F64_mul(v151, v151)
								v153 = base.F64_mul(v152, v152)
								v178 = base.F64_add(base.F64_mul(v151, base.F64_add(v138, base.F64_add(base.F64_mul(v153, base.F64_add(base.F64_mul(v153, base.F64_add(base.F64_mul(v153, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v152, base.F64_add(base.F64_mul(v153, base.F64_add(base.F64_mul(v153, base.F64_add(base.F64_mul(v153, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v135, v143), v138))
								v190 = base.F64_add(v146, base.F64_add(base.F64_add(v145, base.F64_sub(v122, v146)), base.F64_add(base.F64_mul(v178, v144), base.F64_add(base.F64_mul(v120, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v178, v143), float64(2.5082946711645275e-11))))))
								v212 = v190
							} else {
								v212 = float64(0)
							}
						} else {
							v111 = v80
							v112 = v95
							v114 = base.I32_wrap_i64(v97)
							v116 = v114 + int32(614242)
							v120 = base.F64_convert_i32_s(v112 + int32(base.Ui32(v116)>>(uint(int32(20))%32)))
							v122 = base.F64_mul(v120, float64(0.30102999566361177))
							v135 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v116&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v111&int64(4294967295)), float64(-1))
							v138 = base.F64_mul(v135, base.F64_mul(v135, float64(0.5)))
							v143 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v135, v138)) & int64(-4294967296))
							v144 = float64(0.4342944818781689)
							v145 = base.F64_mul(v143, v144)
							v146 = base.F64_add(v122, v145)
							v151 = base.F64_div(v135, base.F64_add(v135, float64(2)))
							v152 = base.F64_mul(v151, v151)
							v153 = base.F64_mul(v152, v152)
							v178 = base.F64_add(base.F64_mul(v151, base.F64_add(v138, base.F64_add(base.F64_mul(v153, base.F64_add(base.F64_mul(v153, base.F64_add(base.F64_mul(v153, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v152, base.F64_add(base.F64_mul(v153, base.F64_add(base.F64_mul(v153, base.F64_add(base.F64_mul(v153, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v135, v143), v138))
							v190 = base.F64_add(v146, base.F64_add(base.F64_add(v145, base.F64_sub(v122, v146)), base.F64_add(base.F64_mul(v178, v144), base.F64_add(base.F64_mul(v120, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v178, v143), float64(2.5082946711645275e-11))))))
							v212 = v190
						}
					}
				} else {
					if base.F64_ne(v69, float64(0)) != 0 {
						if int64(-1) < v80 {
							v106 = base.I64_reinterpret_f64(base.F64_mul(v69, float64(1.8014398509481984e+16)))
							v111 = v106
							v112 = int32(-1077)
							v114 = base.I32_wrap_i64(int64(base.Ui64(v106) >> (uint(int64(32)) % 64)))
							v116 = v114 + int32(614242)
							v120 = base.F64_convert_i32_s(v112 + int32(base.Ui32(v116)>>(uint(int32(20))%32)))
							v122 = base.F64_mul(v120, float64(0.30102999566361177))
							v135 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v116&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v111&int64(4294967295)), float64(-1))
							v138 = base.F64_mul(v135, base.F64_mul(v135, float64(0.5)))
							v143 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v135, v138)) & int64(-4294967296))
							v144 = float64(0.4342944818781689)
							v145 = base.F64_mul(v143, v144)
							v146 = base.F64_add(v122, v145)
							v151 = base.F64_div(v135, base.F64_add(v135, float64(2)))
							v152 = base.F64_mul(v151, v151)
							v153 = base.F64_mul(v152, v152)
							v178 = base.F64_add(base.F64_mul(v151, base.F64_add(v138, base.F64_add(base.F64_mul(v153, base.F64_add(base.F64_mul(v153, base.F64_add(base.F64_mul(v153, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v152, base.F64_add(base.F64_mul(v153, base.F64_add(base.F64_mul(v153, base.F64_add(base.F64_mul(v153, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v135, v143), v138))
							v190 = base.F64_add(v146, base.F64_add(base.F64_add(v145, base.F64_sub(v122, v146)), base.F64_add(base.F64_mul(v178, v144), base.F64_add(base.F64_mul(v120, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v178, v143), float64(2.5082946711645275e-11))))))
							v212 = v190
						} else {
							v212 = base.F64_div(base.F64_sub(v69, v69), float64(0))
						}
					} else {
						v212 = base.F64_div(float64(-1), base.F64_mul(v69, v69))
					}
				}
				v216 = base.F32_demote_f64(base.F64_mul(v212, float64(10)))
			}
		}
		*(*float32)(unsafe.Add(mBase, uint32(v13)+4)) = v216
		v219 = int64(base.Ui64(v59) >> (uint(int64(2)) % 64))
		v220 = *(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[14])))
		v222 = base.B2i32(base.Ui64(v59) < base.Ui64(int64(4)))
		if base.Ui64(v59) < base.Ui64(int64(4)) {
			v376 = v57
		} else {
			if v220 == int64(0) {
				v376 = v57
			} else {
				v229 = base.F64_div(base.F64_mul(base.F64_convert_i64_u(v219), float64(65025)), base.F64_convert_i64_u(v220))
				v240 = base.I64_reinterpret_f64(v229)
				if int64(4503599627370495) < v240 {
					if base.Ui64(int64(9218868437227405311)) < base.Ui64(v240) {
						v350 = v229
						v372 = v350
					} else {
						v255 = int32(-1023)
						v257 = int64(base.Ui64(v240) >> (uint(int64(32)) % 64))
						if v257 == int64(1072693248) {
							if base.I32_wrap_i64(v240) != 0 {
								v271 = v240
								v272 = v255
								v274 = int32(1072693248)
								v276 = v274 + int32(614242)
								v280 = base.F64_convert_i32_s(v272 + int32(base.Ui32(v276)>>(uint(int32(20))%32)))
								v282 = base.F64_mul(v280, float64(0.30102999566361177))
								v295 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v276&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v271&int64(4294967295)), float64(-1))
								v298 = base.F64_mul(v295, base.F64_mul(v295, float64(0.5)))
								v303 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v295, v298)) & int64(-4294967296))
								v304 = float64(0.4342944818781689)
								v305 = base.F64_mul(v303, v304)
								v306 = base.F64_add(v282, v305)
								v311 = base.F64_div(v295, base.F64_add(v295, float64(2)))
								v312 = base.F64_mul(v311, v311)
								v313 = base.F64_mul(v312, v312)
								v338 = base.F64_add(base.F64_mul(v311, base.F64_add(v298, base.F64_add(base.F64_mul(v313, base.F64_add(base.F64_mul(v313, base.F64_add(base.F64_mul(v313, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v312, base.F64_add(base.F64_mul(v313, base.F64_add(base.F64_mul(v313, base.F64_add(base.F64_mul(v313, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v295, v303), v298))
								v350 = base.F64_add(v306, base.F64_add(base.F64_add(v305, base.F64_sub(v282, v306)), base.F64_add(base.F64_mul(v338, v304), base.F64_add(base.F64_mul(v280, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v338, v303), float64(2.5082946711645275e-11))))))
								v372 = v350
							} else {
								v372 = float64(0)
							}
						} else {
							v271 = v240
							v272 = v255
							v274 = base.I32_wrap_i64(v257)
							v276 = v274 + int32(614242)
							v280 = base.F64_convert_i32_s(v272 + int32(base.Ui32(v276)>>(uint(int32(20))%32)))
							v282 = base.F64_mul(v280, float64(0.30102999566361177))
							v295 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v276&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v271&int64(4294967295)), float64(-1))
							v298 = base.F64_mul(v295, base.F64_mul(v295, float64(0.5)))
							v303 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v295, v298)) & int64(-4294967296))
							v304 = float64(0.4342944818781689)
							v305 = base.F64_mul(v303, v304)
							v306 = base.F64_add(v282, v305)
							v311 = base.F64_div(v295, base.F64_add(v295, float64(2)))
							v312 = base.F64_mul(v311, v311)
							v313 = base.F64_mul(v312, v312)
							v338 = base.F64_add(base.F64_mul(v311, base.F64_add(v298, base.F64_add(base.F64_mul(v313, base.F64_add(base.F64_mul(v313, base.F64_add(base.F64_mul(v313, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v312, base.F64_add(base.F64_mul(v313, base.F64_add(base.F64_mul(v313, base.F64_add(base.F64_mul(v313, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v295, v303), v298))
							v350 = base.F64_add(v306, base.F64_add(base.F64_add(v305, base.F64_sub(v282, v306)), base.F64_add(base.F64_mul(v338, v304), base.F64_add(base.F64_mul(v280, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v338, v303), float64(2.5082946711645275e-11))))))
							v372 = v350
						}
					}
				} else {
					if base.F64_ne(v229, float64(0)) != 0 {
						if int64(-1) < v240 {
							v266 = base.I64_reinterpret_f64(base.F64_mul(v229, float64(1.8014398509481984e+16)))
							v271 = v266
							v272 = int32(-1077)
							v274 = base.I32_wrap_i64(int64(base.Ui64(v266) >> (uint(int64(32)) % 64)))
							v276 = v274 + int32(614242)
							v280 = base.F64_convert_i32_s(v272 + int32(base.Ui32(v276)>>(uint(int32(20))%32)))
							v282 = base.F64_mul(v280, float64(0.30102999566361177))
							v295 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v276&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v271&int64(4294967295)), float64(-1))
							v298 = base.F64_mul(v295, base.F64_mul(v295, float64(0.5)))
							v303 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v295, v298)) & int64(-4294967296))
							v304 = float64(0.4342944818781689)
							v305 = base.F64_mul(v303, v304)
							v306 = base.F64_add(v282, v305)
							v311 = base.F64_div(v295, base.F64_add(v295, float64(2)))
							v312 = base.F64_mul(v311, v311)
							v313 = base.F64_mul(v312, v312)
							v338 = base.F64_add(base.F64_mul(v311, base.F64_add(v298, base.F64_add(base.F64_mul(v313, base.F64_add(base.F64_mul(v313, base.F64_add(base.F64_mul(v313, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v312, base.F64_add(base.F64_mul(v313, base.F64_add(base.F64_mul(v313, base.F64_add(base.F64_mul(v313, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v295, v303), v298))
							v350 = base.F64_add(v306, base.F64_add(base.F64_add(v305, base.F64_sub(v282, v306)), base.F64_add(base.F64_mul(v338, v304), base.F64_add(base.F64_mul(v280, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v338, v303), float64(2.5082946711645275e-11))))))
							v372 = v350
						} else {
							v372 = base.F64_div(base.F64_sub(v229, v229), float64(0))
						}
					} else {
						v372 = base.F64_div(float64(-1), base.F64_mul(v229, v229))
					}
				}
				v376 = base.F32_demote_f64(base.F64_mul(v372, float64(10)))
			}
		}
		*(*float32)(unsafe.Add(mBase, uint32(v13)+8)) = v376
		v378 = *(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[15])))
		v379 = float32(99)
		if base.Ui64(v59) < base.Ui64(int64(4)) {
			v535 = v379
		} else {
			if v378 == int64(0) {
				v535 = float32(99)
			} else {
				v388 = base.F64_div(base.F64_mul(base.F64_convert_i64_u(v219), float64(65025)), base.F64_convert_i64_u(v378))
				v399 = base.I64_reinterpret_f64(v388)
				if int64(4503599627370495) < v399 {
					if base.Ui64(int64(9218868437227405311)) < base.Ui64(v399) {
						v509 = v388
						v531 = v509
					} else {
						v414 = int32(-1023)
						v416 = int64(base.Ui64(v399) >> (uint(int64(32)) % 64))
						if v416 == int64(1072693248) {
							if base.I32_wrap_i64(v399) != 0 {
								v430 = v399
								v431 = v414
								v433 = int32(1072693248)
								v435 = v433 + int32(614242)
								v439 = base.F64_convert_i32_s(v431 + int32(base.Ui32(v435)>>(uint(int32(20))%32)))
								v441 = base.F64_mul(v439, float64(0.30102999566361177))
								v454 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v435&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v430&int64(4294967295)), float64(-1))
								v457 = base.F64_mul(v454, base.F64_mul(v454, float64(0.5)))
								v462 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v454, v457)) & int64(-4294967296))
								v463 = float64(0.4342944818781689)
								v464 = base.F64_mul(v462, v463)
								v465 = base.F64_add(v441, v464)
								v470 = base.F64_div(v454, base.F64_add(v454, float64(2)))
								v471 = base.F64_mul(v470, v470)
								v472 = base.F64_mul(v471, v471)
								v497 = base.F64_add(base.F64_mul(v470, base.F64_add(v457, base.F64_add(base.F64_mul(v472, base.F64_add(base.F64_mul(v472, base.F64_add(base.F64_mul(v472, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v471, base.F64_add(base.F64_mul(v472, base.F64_add(base.F64_mul(v472, base.F64_add(base.F64_mul(v472, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v454, v462), v457))
								v509 = base.F64_add(v465, base.F64_add(base.F64_add(v464, base.F64_sub(v441, v465)), base.F64_add(base.F64_mul(v497, v463), base.F64_add(base.F64_mul(v439, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v497, v462), float64(2.5082946711645275e-11))))))
								v531 = v509
							} else {
								v531 = float64(0)
							}
						} else {
							v430 = v399
							v431 = v414
							v433 = base.I32_wrap_i64(v416)
							v435 = v433 + int32(614242)
							v439 = base.F64_convert_i32_s(v431 + int32(base.Ui32(v435)>>(uint(int32(20))%32)))
							v441 = base.F64_mul(v439, float64(0.30102999566361177))
							v454 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v435&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v430&int64(4294967295)), float64(-1))
							v457 = base.F64_mul(v454, base.F64_mul(v454, float64(0.5)))
							v462 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v454, v457)) & int64(-4294967296))
							v463 = float64(0.4342944818781689)
							v464 = base.F64_mul(v462, v463)
							v465 = base.F64_add(v441, v464)
							v470 = base.F64_div(v454, base.F64_add(v454, float64(2)))
							v471 = base.F64_mul(v470, v470)
							v472 = base.F64_mul(v471, v471)
							v497 = base.F64_add(base.F64_mul(v470, base.F64_add(v457, base.F64_add(base.F64_mul(v472, base.F64_add(base.F64_mul(v472, base.F64_add(base.F64_mul(v472, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v471, base.F64_add(base.F64_mul(v472, base.F64_add(base.F64_mul(v472, base.F64_add(base.F64_mul(v472, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v454, v462), v457))
							v509 = base.F64_add(v465, base.F64_add(base.F64_add(v464, base.F64_sub(v441, v465)), base.F64_add(base.F64_mul(v497, v463), base.F64_add(base.F64_mul(v439, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v497, v462), float64(2.5082946711645275e-11))))))
							v531 = v509
						}
					}
				} else {
					if base.F64_ne(v388, float64(0)) != 0 {
						if int64(-1) < v399 {
							v425 = base.I64_reinterpret_f64(base.F64_mul(v388, float64(1.8014398509481984e+16)))
							v430 = v425
							v431 = int32(-1077)
							v433 = base.I32_wrap_i64(int64(base.Ui64(v425) >> (uint(int64(32)) % 64)))
							v435 = v433 + int32(614242)
							v439 = base.F64_convert_i32_s(v431 + int32(base.Ui32(v435)>>(uint(int32(20))%32)))
							v441 = base.F64_mul(v439, float64(0.30102999566361177))
							v454 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v435&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v430&int64(4294967295)), float64(-1))
							v457 = base.F64_mul(v454, base.F64_mul(v454, float64(0.5)))
							v462 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v454, v457)) & int64(-4294967296))
							v463 = float64(0.4342944818781689)
							v464 = base.F64_mul(v462, v463)
							v465 = base.F64_add(v441, v464)
							v470 = base.F64_div(v454, base.F64_add(v454, float64(2)))
							v471 = base.F64_mul(v470, v470)
							v472 = base.F64_mul(v471, v471)
							v497 = base.F64_add(base.F64_mul(v470, base.F64_add(v457, base.F64_add(base.F64_mul(v472, base.F64_add(base.F64_mul(v472, base.F64_add(base.F64_mul(v472, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v471, base.F64_add(base.F64_mul(v472, base.F64_add(base.F64_mul(v472, base.F64_add(base.F64_mul(v472, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v454, v462), v457))
							v509 = base.F64_add(v465, base.F64_add(base.F64_add(v464, base.F64_sub(v441, v465)), base.F64_add(base.F64_mul(v497, v463), base.F64_add(base.F64_mul(v439, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v497, v462), float64(2.5082946711645275e-11))))))
							v531 = v509
						} else {
							v531 = base.F64_div(base.F64_sub(v388, v388), float64(0))
						}
					} else {
						v531 = base.F64_div(float64(-1), base.F64_mul(v388, v388))
					}
				}
				v535 = base.F32_demote_f64(base.F64_mul(v531, float64(10)))
			}
		}
		*(*float32)(unsafe.Add(mBase, uint32(v13)+12)) = v535
		v538 = v59 * int64(3)
		if base.Ui64(v538) < base.Ui64(int64(2)) {
			v699 = v379
		} else {
			v542 = v220 + v56 + v378
			if v542 == int64(0) {
				v699 = v379
			} else {
				v551 = base.F64_div(base.F64_mul(base.F64_convert_i64_u(int64(base.Ui64(v538)>>(uint(int64(1))%64))), float64(65025)), base.F64_convert_i64_u(v542))
				v562 = base.I64_reinterpret_f64(v551)
				if int64(4503599627370495) < v562 {
					if base.Ui64(int64(9218868437227405311)) < base.Ui64(v562) {
						v672 = v551
						v694 = v672
					} else {
						v577 = int32(-1023)
						v579 = int64(base.Ui64(v562) >> (uint(int64(32)) % 64))
						if v579 == int64(1072693248) {
							if base.I32_wrap_i64(v562) != 0 {
								v593 = v562
								v594 = v577
								v596 = int32(1072693248)
								v598 = v596 + int32(614242)
								v602 = base.F64_convert_i32_s(v594 + int32(base.Ui32(v598)>>(uint(int32(20))%32)))
								v604 = base.F64_mul(v602, float64(0.30102999566361177))
								v617 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v598&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v593&int64(4294967295)), float64(-1))
								v620 = base.F64_mul(v617, base.F64_mul(v617, float64(0.5)))
								v625 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v617, v620)) & int64(-4294967296))
								v626 = float64(0.4342944818781689)
								v627 = base.F64_mul(v625, v626)
								v628 = base.F64_add(v604, v627)
								v633 = base.F64_div(v617, base.F64_add(v617, float64(2)))
								v634 = base.F64_mul(v633, v633)
								v635 = base.F64_mul(v634, v634)
								v660 = base.F64_add(base.F64_mul(v633, base.F64_add(v620, base.F64_add(base.F64_mul(v635, base.F64_add(base.F64_mul(v635, base.F64_add(base.F64_mul(v635, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v634, base.F64_add(base.F64_mul(v635, base.F64_add(base.F64_mul(v635, base.F64_add(base.F64_mul(v635, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v617, v625), v620))
								v672 = base.F64_add(v628, base.F64_add(base.F64_add(v627, base.F64_sub(v604, v628)), base.F64_add(base.F64_mul(v660, v626), base.F64_add(base.F64_mul(v602, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v660, v625), float64(2.5082946711645275e-11))))))
								v694 = v672
							} else {
								v694 = float64(0)
							}
						} else {
							v593 = v562
							v594 = v577
							v596 = base.I32_wrap_i64(v579)
							v598 = v596 + int32(614242)
							v602 = base.F64_convert_i32_s(v594 + int32(base.Ui32(v598)>>(uint(int32(20))%32)))
							v604 = base.F64_mul(v602, float64(0.30102999566361177))
							v617 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v598&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v593&int64(4294967295)), float64(-1))
							v620 = base.F64_mul(v617, base.F64_mul(v617, float64(0.5)))
							v625 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v617, v620)) & int64(-4294967296))
							v626 = float64(0.4342944818781689)
							v627 = base.F64_mul(v625, v626)
							v628 = base.F64_add(v604, v627)
							v633 = base.F64_div(v617, base.F64_add(v617, float64(2)))
							v634 = base.F64_mul(v633, v633)
							v635 = base.F64_mul(v634, v634)
							v660 = base.F64_add(base.F64_mul(v633, base.F64_add(v620, base.F64_add(base.F64_mul(v635, base.F64_add(base.F64_mul(v635, base.F64_add(base.F64_mul(v635, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v634, base.F64_add(base.F64_mul(v635, base.F64_add(base.F64_mul(v635, base.F64_add(base.F64_mul(v635, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v617, v625), v620))
							v672 = base.F64_add(v628, base.F64_add(base.F64_add(v627, base.F64_sub(v604, v628)), base.F64_add(base.F64_mul(v660, v626), base.F64_add(base.F64_mul(v602, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v660, v625), float64(2.5082946711645275e-11))))))
							v694 = v672
						}
					}
				} else {
					if base.F64_ne(v551, float64(0)) != 0 {
						if int64(-1) < v562 {
							v588 = base.I64_reinterpret_f64(base.F64_mul(v551, float64(1.8014398509481984e+16)))
							v593 = v588
							v594 = int32(-1077)
							v596 = base.I32_wrap_i64(int64(base.Ui64(v588) >> (uint(int64(32)) % 64)))
							v598 = v596 + int32(614242)
							v602 = base.F64_convert_i32_s(v594 + int32(base.Ui32(v598)>>(uint(int32(20))%32)))
							v604 = base.F64_mul(v602, float64(0.30102999566361177))
							v617 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v598&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v593&int64(4294967295)), float64(-1))
							v620 = base.F64_mul(v617, base.F64_mul(v617, float64(0.5)))
							v625 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v617, v620)) & int64(-4294967296))
							v626 = float64(0.4342944818781689)
							v627 = base.F64_mul(v625, v626)
							v628 = base.F64_add(v604, v627)
							v633 = base.F64_div(v617, base.F64_add(v617, float64(2)))
							v634 = base.F64_mul(v633, v633)
							v635 = base.F64_mul(v634, v634)
							v660 = base.F64_add(base.F64_mul(v633, base.F64_add(v620, base.F64_add(base.F64_mul(v635, base.F64_add(base.F64_mul(v635, base.F64_add(base.F64_mul(v635, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v634, base.F64_add(base.F64_mul(v635, base.F64_add(base.F64_mul(v635, base.F64_add(base.F64_mul(v635, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v617, v625), v620))
							v672 = base.F64_add(v628, base.F64_add(base.F64_add(v627, base.F64_sub(v604, v628)), base.F64_add(base.F64_mul(v660, v626), base.F64_add(base.F64_mul(v602, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v660, v625), float64(2.5082946711645275e-11))))))
							v694 = v672
						} else {
							v694 = base.F64_div(base.F64_sub(v551, v551), float64(0))
						}
					} else {
						v694 = base.F64_div(float64(-1), base.F64_mul(v551, v551))
					}
				}
				v699 = base.F32_demote_f64(base.F64_mul(v694, float64(10)))
			}
		}
		*(*float32)(unsafe.Add(mBase, uint32(v13)+16)) = v699
		v701 = float32(99)
		if v59 == int64(0) {
			v857 = v701
		} else {
			v702 = *(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[16])))
			if v702 == int64(0) {
				v857 = v701
			} else {
				v709 = base.F64_div(base.F64_mul(base.F64_convert_i64_u(v59), float64(65025)), base.F64_convert_i64_u(v702))
				v720 = base.I64_reinterpret_f64(v709)
				if int64(4503599627370495) < v720 {
					if base.Ui64(int64(9218868437227405311)) < base.Ui64(v720) {
						v830 = v709
						v852 = v830
					} else {
						v735 = int32(-1023)
						v737 = int64(base.Ui64(v720) >> (uint(int64(32)) % 64))
						if v737 == int64(1072693248) {
							if base.I32_wrap_i64(v720) != 0 {
								v751 = v720
								v752 = v735
								v754 = int32(1072693248)
								v756 = v754 + int32(614242)
								v760 = base.F64_convert_i32_s(v752 + int32(base.Ui32(v756)>>(uint(int32(20))%32)))
								v762 = base.F64_mul(v760, float64(0.30102999566361177))
								v775 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v756&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v751&int64(4294967295)), float64(-1))
								v778 = base.F64_mul(v775, base.F64_mul(v775, float64(0.5)))
								v783 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v775, v778)) & int64(-4294967296))
								v784 = float64(0.4342944818781689)
								v785 = base.F64_mul(v783, v784)
								v786 = base.F64_add(v762, v785)
								v791 = base.F64_div(v775, base.F64_add(v775, float64(2)))
								v792 = base.F64_mul(v791, v791)
								v793 = base.F64_mul(v792, v792)
								v818 = base.F64_add(base.F64_mul(v791, base.F64_add(v778, base.F64_add(base.F64_mul(v793, base.F64_add(base.F64_mul(v793, base.F64_add(base.F64_mul(v793, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v792, base.F64_add(base.F64_mul(v793, base.F64_add(base.F64_mul(v793, base.F64_add(base.F64_mul(v793, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v775, v783), v778))
								v830 = base.F64_add(v786, base.F64_add(base.F64_add(v785, base.F64_sub(v762, v786)), base.F64_add(base.F64_mul(v818, v784), base.F64_add(base.F64_mul(v760, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v818, v783), float64(2.5082946711645275e-11))))))
								v852 = v830
							} else {
								v852 = float64(0)
							}
						} else {
							v751 = v720
							v752 = v735
							v754 = base.I32_wrap_i64(v737)
							v756 = v754 + int32(614242)
							v760 = base.F64_convert_i32_s(v752 + int32(base.Ui32(v756)>>(uint(int32(20))%32)))
							v762 = base.F64_mul(v760, float64(0.30102999566361177))
							v775 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v756&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v751&int64(4294967295)), float64(-1))
							v778 = base.F64_mul(v775, base.F64_mul(v775, float64(0.5)))
							v783 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v775, v778)) & int64(-4294967296))
							v784 = float64(0.4342944818781689)
							v785 = base.F64_mul(v783, v784)
							v786 = base.F64_add(v762, v785)
							v791 = base.F64_div(v775, base.F64_add(v775, float64(2)))
							v792 = base.F64_mul(v791, v791)
							v793 = base.F64_mul(v792, v792)
							v818 = base.F64_add(base.F64_mul(v791, base.F64_add(v778, base.F64_add(base.F64_mul(v793, base.F64_add(base.F64_mul(v793, base.F64_add(base.F64_mul(v793, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v792, base.F64_add(base.F64_mul(v793, base.F64_add(base.F64_mul(v793, base.F64_add(base.F64_mul(v793, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v775, v783), v778))
							v830 = base.F64_add(v786, base.F64_add(base.F64_add(v785, base.F64_sub(v762, v786)), base.F64_add(base.F64_mul(v818, v784), base.F64_add(base.F64_mul(v760, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v818, v783), float64(2.5082946711645275e-11))))))
							v852 = v830
						}
					}
				} else {
					if base.F64_ne(v709, float64(0)) != 0 {
						if int64(-1) < v720 {
							v746 = base.I64_reinterpret_f64(base.F64_mul(v709, float64(1.8014398509481984e+16)))
							v751 = v746
							v752 = int32(-1077)
							v754 = base.I32_wrap_i64(int64(base.Ui64(v746) >> (uint(int64(32)) % 64)))
							v756 = v754 + int32(614242)
							v760 = base.F64_convert_i32_s(v752 + int32(base.Ui32(v756)>>(uint(int32(20))%32)))
							v762 = base.F64_mul(v760, float64(0.30102999566361177))
							v775 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v756&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v751&int64(4294967295)), float64(-1))
							v778 = base.F64_mul(v775, base.F64_mul(v775, float64(0.5)))
							v783 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v775, v778)) & int64(-4294967296))
							v784 = float64(0.4342944818781689)
							v785 = base.F64_mul(v783, v784)
							v786 = base.F64_add(v762, v785)
							v791 = base.F64_div(v775, base.F64_add(v775, float64(2)))
							v792 = base.F64_mul(v791, v791)
							v793 = base.F64_mul(v792, v792)
							v818 = base.F64_add(base.F64_mul(v791, base.F64_add(v778, base.F64_add(base.F64_mul(v793, base.F64_add(base.F64_mul(v793, base.F64_add(base.F64_mul(v793, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v792, base.F64_add(base.F64_mul(v793, base.F64_add(base.F64_mul(v793, base.F64_add(base.F64_mul(v793, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v775, v783), v778))
							v830 = base.F64_add(v786, base.F64_add(base.F64_add(v785, base.F64_sub(v762, v786)), base.F64_add(base.F64_mul(v818, v784), base.F64_add(base.F64_mul(v760, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v818, v783), float64(2.5082946711645275e-11))))))
							v852 = v830
						} else {
							v852 = base.F64_div(base.F64_sub(v709, v709), float64(0))
						}
					} else {
						v852 = base.F64_div(float64(-1), base.F64_mul(v709, v709))
					}
				}
				v857 = base.F32_demote_f64(base.F64_mul(v852, float64(10)))
			}
		}
		v858 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[17])))
		*(*int32)(unsafe.Add(mBase, uint32(v13))) = v858
		v860 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[18])))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+24)) = v860
		v862 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[19])))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+28)) = v862
		*(*float32)(unsafe.Add(mBase, uint32(v13)+20)) = v857
		v865 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_StoreStats[20])))
		*(*int32)(unsafe.Add(mBase, uint32(v13)+32)) = v865
	}
	return
}
func F_Sync(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	v2 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	return base.B2i32(v2 == int32(0))
}
func F_sbrk(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	if l0 != 0 {
		if l0&int32(_a_F_sbrk_0) != 0 {
			base.Wasm_trap_unreachable()
			for {
			}
		} else {
			if l0 <= int32(-1) {
				base.Wasm_trap_unreachable()
				for {
				}
			} else {
				v12 = base.MemoryGrow(m, int32(base.Ui32(l0)>>(uint(int32(16))%32)))
				mBase = m.M
				if v12 != int32(-1) {
					return v12 << (uint(int32(16)) % 32)
				} else {
					*(*int32)(unsafe.Add(mBase, _c_F_sbrk[0])) = int32(48)
					return int32(-1)
				}
			}
		}
	} else {
		return base.MemorySize(m) << (uint(int32(16)) % 32)
	}
}
func F_sift(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v17 int32
	_ = v17
	var v31 int32
	_ = v31
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v41 int32
	_ = v41
	var v43 int32
	_ = v43
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v52 int32
	_ = v52
	var v55 int32
	_ = v55
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v63 int32
	_ = v63
	var v66 int32
	_ = v66
	var v71 int32
	_ = v71
	var v78 int32
	_ = v78
	var v89 int32
	_ = v89
	var v102 int32
	_ = v102
	var v103 int32
	_ = v103
	var v104 int32
	_ = v104
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v113 int32
	_ = v113
	var v117 int32
	_ = v117
	var v121 int32
	_ = v121
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v126 int32
	_ = v126
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v132 int32
	_ = v132
	var v136 int32
	_ = v136
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v145 int32
	_ = v145
	v15 = m.G0
	v17 = v15 - int32(496)
	m.G0 = v17
	*(*int32)(unsafe.Add(mBase, uint32(v17))) = l0
	if l4 <= int32(1) {
		goto L1
	} else {
		goto L2
	}
L1:
	;
	m.G0 = v17 + int32(496)
	return
L2:
	;
	v31 = l4
	v35 = v17 | int32(4)
	v36 = int32(1)
	v37 = l0
	goto L4
L3:
	;
	if base.Ui32(v71) < base.Ui32(int32(2)) {
		goto L1
	} else {
		goto L16
	}
L4:
	;
	v41 = v37 + (int32(0) - l1)
	v43 = v31 + int32(-2)
	v47 = *(*int32)(unsafe.Add(mBase, uint32(l5+v43<<(uint(int32(2))%32))))
	v48 = v41 - v47
	v49 = m.T0[l2].(func(*base.Module, int32, int32, int32) int32)(m, l0, v48, l3)
	mBase = m.M
	if v49 < int32(0) {
		goto L6
	} else {
		goto L7
	}
L5:
	;
	v71 = v63
	goto L3
L6:
	;
	v55 = m.T0[l2].(func(*base.Module, int32, int32, int32) int32)(m, v48, v41, l3)
	mBase = m.M
	v57 = base.B2i32(int32(-1) < v55)
	if int32(-1) < v55 {
		goto L9
	} else {
		goto L10
	}
L7:
	;
	v52 = m.T0[l2].(func(*base.Module, int32, int32, int32) int32)(m, l0, v41, l3)
	mBase = m.M
	if int32(-1) < v52 {
		v71 = v36
		goto L3
	} else {
		goto L8
	}
L8:
	;
	goto L6
L9:
	;
	v58 = v48
	goto L11
L10:
	;
	v58 = v41
	goto L11
L11:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v58
	v63 = v36 + int32(1)
	if int32(-1) < v55 {
		goto L12
	} else {
		goto L13
	}
L12:
	;
	v66 = v31 + int32(-1)
	goto L14
L13:
	;
	v66 = v43
	goto L14
L14:
	;
	if int32(1) < v66 {
		v31 = v66
		v35 = v35 + int32(4)
		v36 = v63
		v37 = v58
		goto L4
	} else {
		goto L15
	}
L15:
	;
	goto L5
L16:
	;
	v78 = v17 + v71<<(uint(int32(2))%32)
	*(*int32)(unsafe.Add(mBase, uint32(v78))) = v17 + int32(240)
	if l1 == int32(0) {
		goto L1
	} else {
		goto L17
	}
L17:
	;
	v89 = l1
	goto L18
L18:
	;
	v102 = *(*int32)(unsafe.Add(mBase, uint32(v78)))
	v103 = *(*int32)(unsafe.Add(mBase, uint32(v17)))
	v104 = int32(256)
	if base.Ui32(v89) < base.Ui32(v104) {
		goto L20
	} else {
		goto L21
	}
L19:
	;
	goto L1
L20:
	;
	v107 = v89
	goto L22
L21:
	;
	v107 = v104
	goto L22
L22:
	;
	v108 = F_memcpy(m, v102, v103, v107)
	mBase = m.M
	v113 = v17
	v117 = v103
	v121 = v71 & int32(-2)
	goto L23
L23:
	;
	v124 = v113 + int32(4)
	v125 = *(*int32)(unsafe.Add(mBase, uint32(v124)))
	v126 = F_memcpy(m, v117, v125, v107)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v113))) = v126 + v107
	v130 = v113 + int32(8)
	v131 = *(*int32)(unsafe.Add(mBase, uint32(v130)))
	v132 = F_memcpy(m, v125, v131, v107)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v124))) = v132 + v107
	v136 = v121 + int32(-2)
	if v136 != 0 {
		v113 = v130
		v117 = v131
		v121 = v136
		goto L23
	} else {
		goto L25
	}
L24:
	;
	if v71&int32(1) == int32(0) {
		goto L26
	} else {
		goto L27
	}
L25:
	;
	goto L24
L26:
	;
	v145 = v89 - v107
	if v145 != 0 {
		v89 = v145
		goto L18
	} else {
		goto L28
	}
L27:
	;
	v141 = *(*int32)(unsafe.Add(mBase, uint32(v113+int32(12))))
	v142 = F_memcpy(m, v131, v141, v107)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v130))) = v142 + v107
	goto L26
L28:
	;
	goto L19
}
func F_specialcase(m *base.Module, l0 float64, l1 int64, l2 int64) float64 {
	var v14 float64
	_ = v14
	var v21 int64
	_ = v21
	var v22 float64
	_ = v22
	var v23 float64
	_ = v23
	var v24 float64
	_ = v24
	var v37 float64
	_ = v37
	var v38 float64
	_ = v38
	var v45 float64
	_ = v45
	var v48 float64
	_ = v48
	var v49 float64
	_ = v49
	if l2&int64(2147483648) != int64(0) {
		v21 = l1 + int64(4602678819172646912)
		v22 = base.F64_reinterpret_i64(v21)
		v23 = base.F64_mul(v22, l0)
		v24 = base.F64_add(v23, v22)
		if base.F64_lt(base.F64_abs(v24), float64(1)) == int32(0) {
			v49 = v24
		} else {
			if base.F64_lt(v24, float64(0)) != 0 {
				v37 = float64(-1)
			} else {
				v37 = float64(1)
			}
			v38 = base.F64_add(v24, v37)
			v45 = base.F64_sub(base.F64_add(v38, base.F64_add(base.F64_add(v23, base.F64_sub(v22, v24)), base.F64_add(v24, base.F64_sub(v37, v38)))), v37)
			if base.F64_eq(v45, float64(0)) != 0 {
				v48 = base.F64_reinterpret_i64(v21 & int64(-9223372036854775807-1))
			} else {
				v48 = v45
			}
			v49 = v48
		}
		return base.F64_mul(v49, float64(2.2250738585072014e-308))
	} else {
		v14 = base.F64_reinterpret_i64(l1 + int64(-4544132024016830464))
		return base.F64_mul(base.F64_add(base.F64_mul(v14, l0), v14), float64(5.486124068793689e+303))
	}
}
