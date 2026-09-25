//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_MapARGB_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v16 int32
	_ = v16
	var v19 int32
	_ = v19
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v33 int32
	_ = v33
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v41 int32
	_ = v41
	var v47 int32
	_ = v47
	var v51 int32
	_ = v51
	var v57 int32
	_ = v57
	var v61 int32
	_ = v61
	var v67 int32
	_ = v67
	var v71 int32
	_ = v71
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v81 int32
	_ = v81
	var v90 int32
	_ = v90
	var v92 int32
	_ = v92
	var v95 int32
	_ = v95
	var v99 int32
	_ = v99
	var v103 int32
	_ = v103
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v108 int32
	_ = v108
	var v110 int32
	_ = v110
	var v111 int32
	_ = v111
	var v113 int32
	_ = v113
	var v121 int32
	_ = v121
	if l4 <= l3 {
	} else {
		if l5 < int32(1) {
		} else {
			v16 = l5 & int32(3)
			v19 = l0
			v21 = l2
			v22 = l3
			for {
				if base.Ui32(l5) < base.Ui32(int32(4)) {
					v79 = v19
					v81 = v21
				} else {
					v28 = v19
					v30 = v21
					v33 = l5 & int32(2147483644)
					for {
						v37 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28)+1)))
						v38 = int32(2)
						v41 = *(*int32)(unsafe.Add(mBase, uint32(l1+v37<<(uint(v38)%32))))
						*(*int32)(unsafe.Add(mBase, uint32(v30))) = v41
						v47 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28+int32(5)))))
						v51 = *(*int32)(unsafe.Add(mBase, uint32(l1+v47<<(uint(v38)%32))))
						*(*int32)(unsafe.Add(mBase, uint32(v30+int32(4)))) = v51
						v57 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28+int32(9)))))
						v61 = *(*int32)(unsafe.Add(mBase, uint32(l1+v57<<(uint(v38)%32))))
						*(*int32)(unsafe.Add(mBase, uint32(v30+int32(8)))) = v61
						v67 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28+int32(13)))))
						v71 = *(*int32)(unsafe.Add(mBase, uint32(l1+v67<<(uint(v38)%32))))
						*(*int32)(unsafe.Add(mBase, uint32(v30+int32(12)))) = v71
						v73 = int32(16)
						v74 = v30 + v73
						v76 = v28 + v73
						v78 = v33 + int32(-4)
						if v78 != 0 {
							v28 = v76
							v30 = v74
							v33 = v78
							continue
						} else {
							break
						}
						break
					}
					v79 = v76
					v81 = v74
				}
				if v16 == int32(0) {
					v111 = v79
					v113 = v81
				} else {
					v90 = v79
					v92 = v81
					v95 = v16
					for {
						v99 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v90)+1)))
						v103 = *(*int32)(unsafe.Add(mBase, uint32(l1+v99<<(uint(int32(2))%32))))
						*(*int32)(unsafe.Add(mBase, uint32(v92))) = v103
						v105 = int32(4)
						v106 = v92 + v105
						v108 = v90 + v105
						v110 = v95 + int32(-1)
						if v110 != 0 {
							v90 = v108
							v92 = v106
							v95 = v110
							continue
						} else {
							break
						}
						break
					}
					v111 = v108
					v113 = v106
				}
				v121 = v22 + int32(1)
				if v121 != l4 {
					v19 = v111
					v21 = v113
					v22 = v121
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
func F_MapAlpha_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v16 int32
	_ = v16
	var v19 int32
	_ = v19
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v33 int32
	_ = v33
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v41 int32
	_ = v41
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v49 int32
	_ = v49
	var v53 int32
	_ = v53
	var v55 int32
	_ = v55
	var v61 int32
	_ = v61
	var v65 int32
	_ = v65
	var v67 int32
	_ = v67
	var v69 int32
	_ = v69
	var v73 int32
	_ = v73
	var v77 int32
	_ = v77
	var v79 int32
	_ = v79
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v84 int32
	_ = v84
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v89 int32
	_ = v89
	var v98 int32
	_ = v98
	var v100 int32
	_ = v100
	var v103 int32
	_ = v103
	var v107 int32
	_ = v107
	var v111 int32
	_ = v111
	var v113 int32
	_ = v113
	var v115 int32
	_ = v115
	var v116 int32
	_ = v116
	var v118 int32
	_ = v118
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v123 int32
	_ = v123
	var v131 int32
	_ = v131
	if l4 <= l3 {
	} else {
		if l5 < int32(1) {
		} else {
			v16 = l5 & int32(3)
			v19 = l0
			v21 = l2
			v22 = l3
			for {
				if base.Ui32(l5) < base.Ui32(int32(4)) {
					v87 = v19
					v89 = v21
				} else {
					v28 = v19
					v30 = v21
					v33 = l5 & int32(2147483644)
					for {
						v37 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28))))
						v38 = int32(2)
						v41 = *(*int32)(unsafe.Add(mBase, uint32(l1+v37<<(uint(v38)%32))))
						v42 = int32(8)
						v43 = int32(base.Ui32(v41) >> (uint(v42) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v30))) = uint8(v43)
						v45 = int32(1)
						v49 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28+v45))))
						v53 = *(*int32)(unsafe.Add(mBase, uint32(l1+v49<<(uint(v38)%32))))
						v55 = int32(base.Ui32(v53) >> (uint(v42) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v30+v45))) = uint8(v55)
						v61 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28+v38))))
						v65 = *(*int32)(unsafe.Add(mBase, uint32(l1+v61<<(uint(v38)%32))))
						v67 = int32(base.Ui32(v65) >> (uint(v42) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v30+v38))) = uint8(v67)
						v69 = int32(3)
						v73 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28+v69))))
						v77 = *(*int32)(unsafe.Add(mBase, uint32(l1+v73<<(uint(v38)%32))))
						v79 = int32(base.Ui32(v77) >> (uint(v42) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v30+v69))) = uint8(v79)
						v81 = int32(4)
						v82 = v30 + v81
						v84 = v28 + v81
						v86 = v33 + int32(-4)
						if v86 != 0 {
							v28 = v84
							v30 = v82
							v33 = v86
							continue
						} else {
							break
						}
						break
					}
					v87 = v84
					v89 = v82
				}
				if v16 == int32(0) {
					v121 = v87
					v123 = v89
				} else {
					v98 = v87
					v100 = v89
					v103 = v16
					for {
						v107 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v98))))
						v111 = *(*int32)(unsafe.Add(mBase, uint32(l1+v107<<(uint(int32(2))%32))))
						v113 = int32(base.Ui32(v111) >> (uint(int32(8)) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v100))) = uint8(v113)
						v115 = int32(1)
						v116 = v100 + v115
						v118 = v98 + v115
						v120 = v103 + int32(-1)
						if v120 != 0 {
							v98 = v118
							v100 = v116
							v103 = v120
							continue
						} else {
							break
						}
						break
					}
					v121 = v118
					v123 = v116
				}
				v131 = v22 + int32(1)
				if v131 != l4 {
					v19 = v121
					v21 = v123
					v22 = v131
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
func F_Mean16x4_C(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v4 int32
	_ = v4
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
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
	var v35 int32
	_ = v35
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
	var v56 int32
	_ = v56
	var v58 int32
	_ = v58
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v67 int32
	_ = v67
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v72 int32
	_ = v72
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v78 int32
	_ = v78
	var v80 int32
	_ = v80
	var v82 int32
	_ = v82
	var v84 int32
	_ = v84
	var v86 int32
	_ = v86
	var v88 int32
	_ = v88
	var v90 int32
	_ = v90
	var v92 int32
	_ = v92
	var v94 int32
	_ = v94
	var v96 int32
	_ = v96
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v102 int32
	_ = v102
	var v104 int32
	_ = v104
	var v106 int32
	_ = v106
	var v108 int32
	_ = v108
	var v110 int32
	_ = v110
	var v112 int32
	_ = v112
	var v114 int32
	_ = v114
	var v116 int32
	_ = v116
	var v118 int32
	_ = v118
	var v120 int32
	_ = v120
	var v122 int32
	_ = v122
	var v124 int32
	_ = v124
	var v126 int32
	_ = v126
	var v128 int32
	_ = v128
	v3 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v4 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)))
	v6 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)))
	v8 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)))
	v10 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+32)))
	v12 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+33)))
	v14 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+34)))
	v16 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+35)))
	v18 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+64)))
	v20 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+65)))
	v22 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+66)))
	v24 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+67)))
	v26 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+96)))
	v28 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+97)))
	v30 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+98)))
	v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+99)))
	*(*int32)(unsafe.Add(mBase, uint32(l1))) = v3 + v4 + v6 + v8 + v10 + v12 + v14 + v16 + v18 + v20 + v22 + v24 + v26 + v28 + v30 + v32
	v35 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+4)))
	v36 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+5)))
	v38 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+6)))
	v40 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+7)))
	v42 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+36)))
	v44 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+37)))
	v46 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+38)))
	v48 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+39)))
	v50 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+68)))
	v52 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+69)))
	v54 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+70)))
	v56 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+71)))
	v58 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+100)))
	v60 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+101)))
	v62 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+102)))
	v64 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+103)))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v35 + v36 + v38 + v40 + v42 + v44 + v46 + v48 + v50 + v52 + v54 + v56 + v58 + v60 + v62 + v64
	v67 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+8)))
	v68 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+9)))
	v70 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+10)))
	v72 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+11)))
	v74 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+40)))
	v76 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+41)))
	v78 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+42)))
	v80 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+43)))
	v82 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+72)))
	v84 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+73)))
	v86 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+74)))
	v88 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+75)))
	v90 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+104)))
	v92 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+105)))
	v94 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+106)))
	v96 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+107)))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+8)) = v67 + v68 + v70 + v72 + v74 + v76 + v78 + v80 + v82 + v84 + v86 + v88 + v90 + v92 + v94 + v96
	v99 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+12)))
	v100 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+13)))
	v102 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+14)))
	v104 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+15)))
	v106 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+44)))
	v108 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+45)))
	v110 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+46)))
	v112 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+47)))
	v114 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+76)))
	v116 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+77)))
	v118 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+78)))
	v120 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+79)))
	v122 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+108)))
	v124 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+109)))
	v126 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+110)))
	v128 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+111)))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v99 + v100 + v102 + v104 + v106 + v108 + v110 + v112 + v114 + v116 + v118 + v120 + v122 + v124 + v126 + v128
	return
}
func F_malloc(m *base.Module, l0 int32) int32 {
	var v2 int32
	_ = v2
	v2 = F_dlmalloc(m, l0)
	return v2
}
func F_memcmp(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v19 int32
	_ = v19
	var v24 int32
	_ = v24
	var v29 int32
	_ = v29
	v4 = int32(0)
	if l2 == v4 {
		v29 = v4
		goto L1
	} else {
		goto L2
	}
L1:
	;
	return v29
L2:
	;
	v10 = l0
	v11 = l1
	v12 = l2
	goto L4
L3:
	;
	v29 = v16 - v17
	goto L1
L4:
	;
	v16 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10))))
	v17 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v11))))
	if v16 != v17 {
		goto L3
	} else {
		goto L6
	}
L6:
	;
	v19 = int32(1)
	v24 = v12 + int32(-1)
	if v24 != 0 {
		v10 = v10 + v19
		v11 = v11 + v19
		v12 = v24
		goto L4
	} else {
		goto L7
	}
L7:
	;
	v29 = v4
	goto L1
}
func F_memcpy(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v16 int32
	_ = v16
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	var v30 int32
	_ = v30
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v44 int32
	_ = v44
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v58 int32
	_ = v58
	var v62 int32
	_ = v62
	var v68 int32
	_ = v68
	var v69 int32
	_ = v69
	var v70 int32
	_ = v70
	var v72 int32
	_ = v72
	var v76 int32
	_ = v76
	var v79 int64
	_ = v79
	var v81 int64
	_ = v81
	var v83 int32
	_ = v83
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v94 int32
	_ = v94
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v99 int64
	_ = v99
	var v101 int64
	_ = v101
	var v103 int64
	_ = v103
	var v105 int64
	_ = v105
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v110 int32
	_ = v110
	var v112 int32
	_ = v112
	var v117 int32
	_ = v117
	var v119 int32
	_ = v119
	var v120 int32
	_ = v120
	var v124 int64
	_ = v124
	var v126 int32
	_ = v126
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v136 int32
	_ = v136
	var v138 int32
	_ = v138
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v148 int32
	_ = v148
	var v150 int32
	_ = v150
	var v154 int32
	_ = v154
	var v155 int32
	_ = v155
	var v160 int32
	_ = v160
	var v165 int32
	_ = v165
	var v170 int32
	_ = v170
	var v174 int64
	_ = v174
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v183 int32
	_ = v183
	var v187 int32
	_ = v187
	var v190 int32
	_ = v190
	var v194 int64
	_ = v194
	var v196 int32
	_ = v196
	var v203 int32
	_ = v203
	var v207 int32
	_ = v207
	var v210 int32
	_ = v210
	var v214 int32
	_ = v214
	var v216 int32
	_ = v216
	var v218 int64
	_ = v218
	var v220 int32
	_ = v220
	var v222 int32
	_ = v222
	var v224 int32
	_ = v224
	var v228 int32
	_ = v228
	var v229 int32
	_ = v229
	var v233 int32
	_ = v233
	var v235 int32
	_ = v235
	var v236 int32
	_ = v236
	var v240 int64
	_ = v240
	var v242 int32
	_ = v242
	var v249 int32
	_ = v249
	var v253 int32
	_ = v253
	var v256 int32
	_ = v256
	var v258 int32
	_ = v258
	var v259 int32
	_ = v259
	var v260 int32
	_ = v260
	var v261 int32
	_ = v261
	var v262 int32
	_ = v262
	var v265 int32
	_ = v265
	var v266 int32
	_ = v266
	var v267 int32
	_ = v267
	var v270 int64
	_ = v270
	var v272 int32
	_ = v272
	var v276 int32
	_ = v276
	var v277 int32
	_ = v277
	var v278 int32
	_ = v278
	var v285 int32
	_ = v285
	var v287 int32
	_ = v287
	var v291 int32
	_ = v291
	var v292 int32
	_ = v292
	var v297 int32
	_ = v297
	var v299 int32
	_ = v299
	var v303 int32
	_ = v303
	var v304 int32
	_ = v304
	var v309 int32
	_ = v309
	if base.Ui32(int32(32)) < base.Ui32(l2) {
		base.MemoryCopy(m, l0, l1, l2)
		return l0
	} else {
		if l1&int32(3) == int32(0) {
			v68 = l2
			v69 = l0
			v70 = l1
		} else {
			if l2 == int32(0) {
				v68 = l2
				v69 = l0
				v70 = l1
			} else {
				v16 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
				*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v16)
				v19 = l2 + int32(-1)
				v20 = int32(1)
				v21 = l0 + v20
				v23 = l1 + v20
				if v23&int32(3) == int32(0) {
					v68 = v19
					v69 = v21
					v70 = v23
				} else {
					if v19 == int32(0) {
						v68 = v19
						v69 = v21
						v70 = v23
					} else {
						v30 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+1)))
						*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v30)
						v33 = l2 + int32(-2)
						v34 = int32(2)
						v35 = l0 + v34
						v37 = l1 + v34
						if v37&int32(3) == int32(0) {
							v68 = v33
							v69 = v35
							v70 = v37
						} else {
							if v33 == int32(0) {
								v68 = v33
								v69 = v35
								v70 = v37
							} else {
								v44 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+2)))
								*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v44)
								v47 = l2 + int32(-3)
								v48 = int32(3)
								v49 = l0 + v48
								v51 = l1 + v48
								if v51&v48 == int32(0) {
									v68 = v47
									v69 = v49
									v70 = v51
								} else {
									if v47 == int32(0) {
										v68 = v47
										v69 = v49
										v70 = v51
									} else {
										v58 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3)))
										*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)) = uint8(v58)
										v62 = int32(4)
										v68 = l2 + int32(-4)
										v69 = l0 + v62
										v70 = l1 + v62
									}
								}
							}
						}
					}
				}
			}
		}
		v72 = v69 & int32(3)
		if v72 != 0 {
			if base.Ui32(v68) < base.Ui32(int32(32)) {
				if base.Ui32(int32(16)) <= base.Ui32(v68) {
					v214 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v70))))
					*(*uint8)(unsafe.Add(mBase, uint32(v69))) = uint8(v214)
					v216 = *(*int32)(unsafe.Add(mBase, uint32(v70)+1))
					*(*int32)(unsafe.Add(mBase, uint32(v69)+1)) = v216
					v218 = *(*int64)(unsafe.Add(mBase, uint32(v70)+5))
					*(*int64)(unsafe.Add(mBase, uint32(v69)+5)) = v218
					v220 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v70)+13)))
					*(*uint16)(unsafe.Add(mBase, uint32(v69)+13)) = uint16(v220)
					v222 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v70)+15)))
					*(*uint8)(unsafe.Add(mBase, uint32(v69)+15)) = uint8(v222)
					v224 = int32(16)
					v228 = v70 + v224
					v229 = v69 + v224
				} else {
					v228 = v70
					v229 = v69
				}
				if v68&int32(8) != 0 {
					v265 = v228
					v266 = v229
					v267 = v68
					v270 = *(*int64)(unsafe.Add(mBase, uint32(v265)))
					*(*int64)(unsafe.Add(mBase, uint32(v266))) = v270
					v272 = int32(8)
					v276 = v265 + v272
					v277 = v266 + v272
					v278 = v267
				} else {
					v276 = v228
					v277 = v229
					v278 = v68
				}
			} else {
				v165 = *(*int32)(unsafe.Add(mBase, uint32(v70)))
				*(*uint8)(unsafe.Add(mBase, uint32(v69))) = uint8(v165)
				switch v72 + int32(-1) {
				default:
					v233 = int32(base.Ui32(v165) >> (uint(int32(16)) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v69)+2)) = uint8(v233)
					v235 = int32(8)
					v236 = int32(base.Ui32(v165) >> (uint(v235) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v69)+1)) = uint8(v236)
					v240 = *(*int64)(unsafe.Add(mBase, uint32(v70+int32(7))))
					*(*int64)(unsafe.Add(mBase, uint32(v69)+7)) = v240
					v242 = *(*int32)(unsafe.Add(mBase, uint32(v70)+4))
					*(*int32)(unsafe.Add(mBase, uint32(v69)+3)) = v242<<(uint(v235)%32) | int32(base.Ui32(v165)>>(uint(int32(24))%32))
					v249 = int32(19)
					v253 = int32(15)
					v256 = *(*int32)(unsafe.Add(mBase, uint32(v70+v253)))
					v258 = v70 + v249
					v259 = v69 + v249
					v260 = int32(13)
					v261 = v256
					v262 = v253
				case 1:
					v170 = int32(base.Ui32(v165) >> (uint(int32(8)) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v69)+1)) = uint8(v170)
					v174 = *(*int64)(unsafe.Add(mBase, uint32(v70+int32(6))))
					*(*int64)(unsafe.Add(mBase, uint32(v69)+6)) = v174
					v176 = *(*int32)(unsafe.Add(mBase, uint32(v70)+4))
					v177 = int32(16)
					*(*int32)(unsafe.Add(mBase, uint32(v69)+2)) = v176<<(uint(v177)%32) | int32(base.Ui32(v165)>>(uint(v177)%32))
					v183 = int32(18)
					v187 = int32(14)
					v190 = *(*int32)(unsafe.Add(mBase, uint32(v70+v187)))
					v258 = v70 + v183
					v259 = v69 + v183
					v260 = v187
					v261 = v190
					v262 = v187
				case 2:
					v194 = *(*int64)(unsafe.Add(mBase, uint32(v70+int32(5))))
					*(*int64)(unsafe.Add(mBase, uint32(v69)+5)) = v194
					v196 = *(*int32)(unsafe.Add(mBase, uint32(v70)+4))
					*(*int32)(unsafe.Add(mBase, uint32(v69)+1)) = v196<<(uint(int32(24))%32) | int32(base.Ui32(v165)>>(uint(int32(8))%32))
					v203 = int32(17)
					v207 = int32(13)
					v210 = *(*int32)(unsafe.Add(mBase, uint32(v70+v207)))
					v258 = v70 + v203
					v259 = v69 + v203
					v260 = int32(15)
					v261 = v210
					v262 = v207
				}
				*(*int32)(unsafe.Add(mBase, uint32(v69+v262))) = v261
				v265 = v258
				v266 = v259
				v267 = v260
				v270 = *(*int64)(unsafe.Add(mBase, uint32(v265)))
				*(*int64)(unsafe.Add(mBase, uint32(v266))) = v270
				v272 = int32(8)
				v276 = v265 + v272
				v277 = v266 + v272
				v278 = v267
			}
			if v278&int32(4) == int32(0) {
				v291 = v276
				v292 = v277
			} else {
				v285 = *(*int32)(unsafe.Add(mBase, uint32(v276)))
				*(*int32)(unsafe.Add(mBase, uint32(v277))) = v285
				v287 = int32(4)
				v291 = v276 + v287
				v292 = v277 + v287
			}
			if v278&int32(2) == int32(0) {
				v303 = v291
				v304 = v292
			} else {
				v297 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v291))))
				*(*uint16)(unsafe.Add(mBase, uint32(v292))) = uint16(v297)
				v299 = int32(2)
				v303 = v291 + v299
				v304 = v292 + v299
			}
			if v278&int32(1) == int32(0) {
			} else {
				v309 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v303))))
				*(*uint8)(unsafe.Add(mBase, uint32(v304))) = uint8(v309)
			}
			return l0
		} else {
			if base.Ui32(int32(16)) <= base.Ui32(v68) {
				v76 = v68 + int32(-16)
				if v76&int32(16) != 0 {
					v87 = v68
					v88 = v69
					v89 = v70
				} else {
					v79 = *(*int64)(unsafe.Add(mBase, uint32(v70)))
					*(*int64)(unsafe.Add(mBase, uint32(v69))) = v79
					v81 = *(*int64)(unsafe.Add(mBase, uint32(v70)+8))
					*(*int64)(unsafe.Add(mBase, uint32(v69)+8)) = v81
					v83 = int32(16)
					v87 = v76
					v88 = v69 + v83
					v89 = v70 + v83
				}
				if base.Ui32(v76) < base.Ui32(int32(16)) {
					v117 = v76
					v119 = v88
					v120 = v89
				} else {
					v94 = v87
					v96 = v88
					v97 = v89
					for {
						v99 = *(*int64)(unsafe.Add(mBase, uint32(v97)))
						*(*int64)(unsafe.Add(mBase, uint32(v96))) = v99
						v101 = *(*int64)(unsafe.Add(mBase, uint32(v97)+8))
						*(*int64)(unsafe.Add(mBase, uint32(v96)+8)) = v101
						v103 = *(*int64)(unsafe.Add(mBase, uint32(v97)+16))
						*(*int64)(unsafe.Add(mBase, uint32(v96)+16)) = v103
						v105 = *(*int64)(unsafe.Add(mBase, uint32(v97)+24))
						*(*int64)(unsafe.Add(mBase, uint32(v96)+24)) = v105
						v107 = int32(32)
						v108 = v96 + v107
						v110 = v97 + v107
						v112 = v94 + int32(-32)
						if base.Ui32(int32(15)) < base.Ui32(v112) {
							v94 = v112
							v96 = v108
							v97 = v110
							continue
						} else {
							break
						}
						break
					}
					v117 = v112
					v119 = v108
					v120 = v110
				}
			} else {
				v117 = v68
				v119 = v69
				v120 = v70
			}
			if base.Ui32(v117) < base.Ui32(int32(8)) {
				v130 = v119
				v131 = v120
			} else {
				v124 = *(*int64)(unsafe.Add(mBase, uint32(v120)))
				*(*int64)(unsafe.Add(mBase, uint32(v119))) = v124
				v126 = int32(8)
				v130 = v119 + v126
				v131 = v120 + v126
			}
			if v117&int32(4) == int32(0) {
				v142 = v130
				v143 = v131
			} else {
				v136 = *(*int32)(unsafe.Add(mBase, uint32(v131)))
				*(*int32)(unsafe.Add(mBase, uint32(v130))) = v136
				v138 = int32(4)
				v142 = v130 + v138
				v143 = v131 + v138
			}
			if v117&int32(2) == int32(0) {
				v154 = v142
				v155 = v143
			} else {
				v148 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v143))))
				*(*uint16)(unsafe.Add(mBase, uint32(v142))) = uint16(v148)
				v150 = int32(2)
				v154 = v142 + v150
				v155 = v143 + v150
			}
			if v117&int32(1) == int32(0) {
				return l0
			} else {
				v160 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v155))))
				*(*uint8)(unsafe.Add(mBase, uint32(v154))) = uint8(v160)
				return l0
			}
		}
	}
}
func F_memmove(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v22 int32
	_ = v22
	var v30 int32
	_ = v30
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v42 int32
	_ = v42
	var v45 int32
	_ = v45
	var v47 int32
	_ = v47
	var v54 int32
	_ = v54
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v66 int32
	_ = v66
	var v68 int32
	_ = v68
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v84 int32
	_ = v84
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v94 int32
	_ = v94
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v104 int32
	_ = v104
	var v111 int32
	_ = v111
	var v114 int32
	_ = v114
	var v116 int32
	_ = v116
	var v122 int32
	_ = v122
	var v128 int32
	_ = v128
	var v131 int32
	_ = v131
	var v137 int32
	_ = v137
	var v138 int32
	_ = v138
	var v144 int32
	_ = v144
	var v147 int32
	_ = v147
	var v149 int32
	_ = v149
	var v152 int32
	_ = v152
	var v159 int32
	_ = v159
	var v165 int32
	_ = v165
	var v170 int32
	_ = v170
	var v171 int32
	_ = v171
	var v173 int32
	_ = v173
	var v176 int32
	_ = v176
	var v178 int32
	_ = v178
	var v182 int32
	_ = v182
	var v184 int32
	_ = v184
	var v188 int32
	_ = v188
	var v190 int32
	_ = v190
	var v193 int32
	_ = v193
	var v198 int32
	_ = v198
	var v206 int32
	_ = v206
	var v209 int32
	_ = v209
	var v216 int32
	_ = v216
	var v217 int32
	_ = v217
	var v222 int32
	_ = v222
	var v224 int32
	_ = v224
	var v225 int32
	_ = v225
	var v227 int32
	_ = v227
	var v231 int32
	_ = v231
	var v237 int32
	_ = v237
	var v244 int32
	_ = v244
	var v248 int32
	_ = v248
	var v249 int32
	_ = v249
	var v251 int32
	_ = v251
	var v254 int32
	_ = v254
	var v256 int32
	_ = v256
	var v260 int32
	_ = v260
	var v262 int32
	_ = v262
	var v266 int32
	_ = v266
	var v268 int32
	_ = v268
	var v271 int32
	_ = v271
	var v272 int32
	_ = v272
	var v273 int32
	_ = v273
	var v274 int32
	_ = v274
	var v278 int32
	_ = v278
	var v284 int32
	_ = v284
	var v291 int32
	_ = v291
	var v292 int32
	_ = v292
	var v293 int32
	_ = v293
	var v297 int32
	_ = v297
	var v299 int32
	_ = v299
	var v300 int32
	_ = v300
	var v302 int32
	_ = v302
	var v304 int32
	_ = v304
	var v306 int32
	_ = v306
	var v308 int32
	_ = v308
	var v310 int32
	_ = v310
	var v315 int32
	_ = v315
	var v317 int32
	_ = v317
	var v319 int32
	_ = v319
	var v321 int32
	_ = v321
	var v323 int32
	_ = v323
	var v325 int32
	_ = v325
	var v327 int32
	_ = v327
	var v329 int32
	_ = v329
	var v331 int32
	_ = v331
	var v333 int32
	_ = v333
	var v335 int32
	_ = v335
	var v337 int32
	_ = v337
	var v338 int32
	_ = v338
	var v340 int32
	_ = v340
	var v342 int32
	_ = v342
	var v346 int32
	_ = v346
	var v348 int32
	_ = v348
	var v350 int32
	_ = v350
	var v355 int32
	_ = v355
	var v359 int32
	_ = v359
	var v360 int32
	_ = v360
	var v361 int32
	_ = v361
	var v365 int32
	_ = v365
	var v367 int32
	_ = v367
	var v368 int32
	_ = v368
	var v370 int32
	_ = v370
	var v372 int32
	_ = v372
	var v374 int32
	_ = v374
	var v376 int32
	_ = v376
	var v377 int32
	_ = v377
	var v383 int32
	_ = v383
	var v385 int32
	_ = v385
	var v386 int32
	_ = v386
	var v389 int32
	_ = v389
	var v391 int32
	_ = v391
	var v393 int32
	_ = v393
	var v395 int32
	_ = v395
	var v397 int32
	_ = v397
	var v399 int32
	_ = v399
	var v401 int32
	_ = v401
	var v403 int32
	_ = v403
	var v405 int32
	_ = v405
	var v410 int32
	_ = v410
	if base.Ui32(int32(33)) <= base.Ui32(l2) {
		base.MemoryCopy(m, l0, l1, l2)
	} else {
		if l0 == l1 {
		} else {
			v11 = l0 + l2
			if base.Ui32(int32(0)-l2<<(uint(int32(1))%32)) < base.Ui32(l1-v11) {
				v22 = (l1 ^ l0) & int32(3)
				if base.Ui32(l1) <= base.Ui32(l0) {
					if v22 != 0 {
						v198 = l2
						if v198 == int32(0) {
						} else {
							v206 = v198 & int32(3)
							if v206 == int32(0) {
								v231 = v198
							} else {
								v209 = int32(-1)
								v216 = v198
								v217 = v206
								for {
									v222 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v209+v216))))
									*(*uint8)(unsafe.Add(mBase, uint32(l0+v209+v216))) = uint8(v222)
									v224 = int32(-1)
									v225 = v216 + v224
									v227 = v217 + v224
									if v227 != 0 {
										v216 = v225
										v217 = v227
										continue
									} else {
										break
									}
									break
								}
								v231 = v225
							}
							if base.Ui32(v198) < base.Ui32(int32(4)) {
							} else {
								v237 = int32(-4)
								v244 = v231
								for {
									v248 = l0 + v237 + v244
									v249 = int32(3)
									v251 = l1 + v237 + v244
									v254 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v249))))
									*(*uint8)(unsafe.Add(mBase, uint32(v248+v249))) = uint8(v254)
									v256 = int32(2)
									v260 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v256))))
									*(*uint8)(unsafe.Add(mBase, uint32(v248+v256))) = uint8(v260)
									v262 = int32(1)
									v266 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v262))))
									*(*uint8)(unsafe.Add(mBase, uint32(v248+v262))) = uint8(v266)
									v268 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251))))
									*(*uint8)(unsafe.Add(mBase, uint32(v248))) = uint8(v268)
									v271 = v244 + int32(-4)
									if v271 != 0 {
										v244 = v271
										continue
									} else {
										break
									}
									break
								}
							}
						}
					} else {
						if v11&int32(3) == int32(0) {
							v116 = l2
							if base.Ui32(v116) < base.Ui32(int32(4)) {
								v198 = v116
							} else {
								v122 = v116 + int32(-4)
								v128 = (int32(base.Ui32(v122)>>(uint(int32(2))%32)) + int32(1)) & int32(3)
								if v128 == int32(0) {
									v152 = v116
								} else {
									v131 = int32(-4)
									v137 = v116
									v138 = v128
									for {
										v144 = *(*int32)(unsafe.Add(mBase, uint32(l1+v131+v137)))
										*(*int32)(unsafe.Add(mBase, uint32(l0+v131+v137))) = v144
										v147 = v137 + int32(-4)
										v149 = v138 + int32(-1)
										if v149 != 0 {
											v137 = v147
											v138 = v149
											continue
										} else {
											break
										}
										break
									}
									v152 = v147
								}
								if base.Ui32(v122) < base.Ui32(int32(12)) {
									v198 = v152
								} else {
									v159 = int32(-16)
									v165 = v152
									for {
										v170 = l0 + v159 + v165
										v171 = int32(12)
										v173 = l1 + v159 + v165
										v176 = *(*int32)(unsafe.Add(mBase, uint32(v173+v171)))
										*(*int32)(unsafe.Add(mBase, uint32(v170+v171))) = v176
										v178 = int32(8)
										v182 = *(*int32)(unsafe.Add(mBase, uint32(v173+v178)))
										*(*int32)(unsafe.Add(mBase, uint32(v170+v178))) = v182
										v184 = int32(4)
										v188 = *(*int32)(unsafe.Add(mBase, uint32(v173+v184)))
										*(*int32)(unsafe.Add(mBase, uint32(v170+v184))) = v188
										v190 = *(*int32)(unsafe.Add(mBase, uint32(v173)))
										*(*int32)(unsafe.Add(mBase, uint32(v170))) = v190
										v193 = v165 + int32(-16)
										if base.Ui32(int32(3)) < base.Ui32(v193) {
											v165 = v193
											continue
										} else {
											break
										}
										break
									}
									v198 = v193
								}
							}
							if v198 == int32(0) {
							} else {
								v206 = v198 & int32(3)
								if v206 == int32(0) {
									v231 = v198
								} else {
									v209 = int32(-1)
									v216 = v198
									v217 = v206
									for {
										v222 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v209+v216))))
										*(*uint8)(unsafe.Add(mBase, uint32(l0+v209+v216))) = uint8(v222)
										v224 = int32(-1)
										v225 = v216 + v224
										v227 = v217 + v224
										if v227 != 0 {
											v216 = v225
											v217 = v227
											continue
										} else {
											break
										}
										break
									}
									v231 = v225
								}
								if base.Ui32(v198) < base.Ui32(int32(4)) {
								} else {
									v237 = int32(-4)
									v244 = v231
									for {
										v248 = l0 + v237 + v244
										v249 = int32(3)
										v251 = l1 + v237 + v244
										v254 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v249))))
										*(*uint8)(unsafe.Add(mBase, uint32(v248+v249))) = uint8(v254)
										v256 = int32(2)
										v260 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v256))))
										*(*uint8)(unsafe.Add(mBase, uint32(v248+v256))) = uint8(v260)
										v262 = int32(1)
										v266 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v262))))
										*(*uint8)(unsafe.Add(mBase, uint32(v248+v262))) = uint8(v266)
										v268 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251))))
										*(*uint8)(unsafe.Add(mBase, uint32(v248))) = uint8(v268)
										v271 = v244 + int32(-4)
										if v271 != 0 {
											v244 = v271
											continue
										} else {
											break
										}
										break
									}
								}
							}
						} else {
							if l2 == int32(0) {
							} else {
								v81 = l2 + int32(-1)
								v82 = l0 + v81
								v84 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v81))))
								*(*uint8)(unsafe.Add(mBase, uint32(v82))) = uint8(v84)
								if v82&int32(3) != 0 {
									if v81 == int32(0) {
									} else {
										v91 = l2 + int32(-2)
										v92 = l0 + v91
										v94 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v91))))
										*(*uint8)(unsafe.Add(mBase, uint32(v92))) = uint8(v94)
										if v92&int32(3) != 0 {
											if v91 == int32(0) {
											} else {
												v101 = l2 + int32(-3)
												v102 = l0 + v101
												v104 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v101))))
												*(*uint8)(unsafe.Add(mBase, uint32(v102))) = uint8(v104)
												if v102&int32(3) != 0 {
													if v101 == int32(0) {
													} else {
														v111 = l2 + int32(-4)
														v114 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v111))))
														*(*uint8)(unsafe.Add(mBase, uint32(l0+v111))) = uint8(v114)
														v116 = v111
														if base.Ui32(v116) < base.Ui32(int32(4)) {
															v198 = v116
														} else {
															v122 = v116 + int32(-4)
															v128 = (int32(base.Ui32(v122)>>(uint(int32(2))%32)) + int32(1)) & int32(3)
															if v128 == int32(0) {
																v152 = v116
															} else {
																v131 = int32(-4)
																v137 = v116
																v138 = v128
																for {
																	v144 = *(*int32)(unsafe.Add(mBase, uint32(l1+v131+v137)))
																	*(*int32)(unsafe.Add(mBase, uint32(l0+v131+v137))) = v144
																	v147 = v137 + int32(-4)
																	v149 = v138 + int32(-1)
																	if v149 != 0 {
																		v137 = v147
																		v138 = v149
																		continue
																	} else {
																		break
																	}
																	break
																}
																v152 = v147
															}
															if base.Ui32(v122) < base.Ui32(int32(12)) {
																v198 = v152
															} else {
																v159 = int32(-16)
																v165 = v152
																for {
																	v170 = l0 + v159 + v165
																	v171 = int32(12)
																	v173 = l1 + v159 + v165
																	v176 = *(*int32)(unsafe.Add(mBase, uint32(v173+v171)))
																	*(*int32)(unsafe.Add(mBase, uint32(v170+v171))) = v176
																	v178 = int32(8)
																	v182 = *(*int32)(unsafe.Add(mBase, uint32(v173+v178)))
																	*(*int32)(unsafe.Add(mBase, uint32(v170+v178))) = v182
																	v184 = int32(4)
																	v188 = *(*int32)(unsafe.Add(mBase, uint32(v173+v184)))
																	*(*int32)(unsafe.Add(mBase, uint32(v170+v184))) = v188
																	v190 = *(*int32)(unsafe.Add(mBase, uint32(v173)))
																	*(*int32)(unsafe.Add(mBase, uint32(v170))) = v190
																	v193 = v165 + int32(-16)
																	if base.Ui32(int32(3)) < base.Ui32(v193) {
																		v165 = v193
																		continue
																	} else {
																		break
																	}
																	break
																}
																v198 = v193
															}
														}
														if v198 == int32(0) {
														} else {
															v206 = v198 & int32(3)
															if v206 == int32(0) {
																v231 = v198
															} else {
																v209 = int32(-1)
																v216 = v198
																v217 = v206
																for {
																	v222 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v209+v216))))
																	*(*uint8)(unsafe.Add(mBase, uint32(l0+v209+v216))) = uint8(v222)
																	v224 = int32(-1)
																	v225 = v216 + v224
																	v227 = v217 + v224
																	if v227 != 0 {
																		v216 = v225
																		v217 = v227
																		continue
																	} else {
																		break
																	}
																	break
																}
																v231 = v225
															}
															if base.Ui32(v198) < base.Ui32(int32(4)) {
															} else {
																v237 = int32(-4)
																v244 = v231
																for {
																	v248 = l0 + v237 + v244
																	v249 = int32(3)
																	v251 = l1 + v237 + v244
																	v254 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v249))))
																	*(*uint8)(unsafe.Add(mBase, uint32(v248+v249))) = uint8(v254)
																	v256 = int32(2)
																	v260 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v256))))
																	*(*uint8)(unsafe.Add(mBase, uint32(v248+v256))) = uint8(v260)
																	v262 = int32(1)
																	v266 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v262))))
																	*(*uint8)(unsafe.Add(mBase, uint32(v248+v262))) = uint8(v266)
																	v268 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251))))
																	*(*uint8)(unsafe.Add(mBase, uint32(v248))) = uint8(v268)
																	v271 = v244 + int32(-4)
																	if v271 != 0 {
																		v244 = v271
																		continue
																	} else {
																		break
																	}
																	break
																}
															}
														}
													}
												} else {
													v116 = v101
													if base.Ui32(v116) < base.Ui32(int32(4)) {
														v198 = v116
													} else {
														v122 = v116 + int32(-4)
														v128 = (int32(base.Ui32(v122)>>(uint(int32(2))%32)) + int32(1)) & int32(3)
														if v128 == int32(0) {
															v152 = v116
														} else {
															v131 = int32(-4)
															v137 = v116
															v138 = v128
															for {
																v144 = *(*int32)(unsafe.Add(mBase, uint32(l1+v131+v137)))
																*(*int32)(unsafe.Add(mBase, uint32(l0+v131+v137))) = v144
																v147 = v137 + int32(-4)
																v149 = v138 + int32(-1)
																if v149 != 0 {
																	v137 = v147
																	v138 = v149
																	continue
																} else {
																	break
																}
																break
															}
															v152 = v147
														}
														if base.Ui32(v122) < base.Ui32(int32(12)) {
															v198 = v152
														} else {
															v159 = int32(-16)
															v165 = v152
															for {
																v170 = l0 + v159 + v165
																v171 = int32(12)
																v173 = l1 + v159 + v165
																v176 = *(*int32)(unsafe.Add(mBase, uint32(v173+v171)))
																*(*int32)(unsafe.Add(mBase, uint32(v170+v171))) = v176
																v178 = int32(8)
																v182 = *(*int32)(unsafe.Add(mBase, uint32(v173+v178)))
																*(*int32)(unsafe.Add(mBase, uint32(v170+v178))) = v182
																v184 = int32(4)
																v188 = *(*int32)(unsafe.Add(mBase, uint32(v173+v184)))
																*(*int32)(unsafe.Add(mBase, uint32(v170+v184))) = v188
																v190 = *(*int32)(unsafe.Add(mBase, uint32(v173)))
																*(*int32)(unsafe.Add(mBase, uint32(v170))) = v190
																v193 = v165 + int32(-16)
																if base.Ui32(int32(3)) < base.Ui32(v193) {
																	v165 = v193
																	continue
																} else {
																	break
																}
																break
															}
															v198 = v193
														}
													}
													if v198 == int32(0) {
													} else {
														v206 = v198 & int32(3)
														if v206 == int32(0) {
															v231 = v198
														} else {
															v209 = int32(-1)
															v216 = v198
															v217 = v206
															for {
																v222 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v209+v216))))
																*(*uint8)(unsafe.Add(mBase, uint32(l0+v209+v216))) = uint8(v222)
																v224 = int32(-1)
																v225 = v216 + v224
																v227 = v217 + v224
																if v227 != 0 {
																	v216 = v225
																	v217 = v227
																	continue
																} else {
																	break
																}
																break
															}
															v231 = v225
														}
														if base.Ui32(v198) < base.Ui32(int32(4)) {
														} else {
															v237 = int32(-4)
															v244 = v231
															for {
																v248 = l0 + v237 + v244
																v249 = int32(3)
																v251 = l1 + v237 + v244
																v254 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v249))))
																*(*uint8)(unsafe.Add(mBase, uint32(v248+v249))) = uint8(v254)
																v256 = int32(2)
																v260 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v256))))
																*(*uint8)(unsafe.Add(mBase, uint32(v248+v256))) = uint8(v260)
																v262 = int32(1)
																v266 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v262))))
																*(*uint8)(unsafe.Add(mBase, uint32(v248+v262))) = uint8(v266)
																v268 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251))))
																*(*uint8)(unsafe.Add(mBase, uint32(v248))) = uint8(v268)
																v271 = v244 + int32(-4)
																if v271 != 0 {
																	v244 = v271
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
										} else {
											v116 = v91
											if base.Ui32(v116) < base.Ui32(int32(4)) {
												v198 = v116
											} else {
												v122 = v116 + int32(-4)
												v128 = (int32(base.Ui32(v122)>>(uint(int32(2))%32)) + int32(1)) & int32(3)
												if v128 == int32(0) {
													v152 = v116
												} else {
													v131 = int32(-4)
													v137 = v116
													v138 = v128
													for {
														v144 = *(*int32)(unsafe.Add(mBase, uint32(l1+v131+v137)))
														*(*int32)(unsafe.Add(mBase, uint32(l0+v131+v137))) = v144
														v147 = v137 + int32(-4)
														v149 = v138 + int32(-1)
														if v149 != 0 {
															v137 = v147
															v138 = v149
															continue
														} else {
															break
														}
														break
													}
													v152 = v147
												}
												if base.Ui32(v122) < base.Ui32(int32(12)) {
													v198 = v152
												} else {
													v159 = int32(-16)
													v165 = v152
													for {
														v170 = l0 + v159 + v165
														v171 = int32(12)
														v173 = l1 + v159 + v165
														v176 = *(*int32)(unsafe.Add(mBase, uint32(v173+v171)))
														*(*int32)(unsafe.Add(mBase, uint32(v170+v171))) = v176
														v178 = int32(8)
														v182 = *(*int32)(unsafe.Add(mBase, uint32(v173+v178)))
														*(*int32)(unsafe.Add(mBase, uint32(v170+v178))) = v182
														v184 = int32(4)
														v188 = *(*int32)(unsafe.Add(mBase, uint32(v173+v184)))
														*(*int32)(unsafe.Add(mBase, uint32(v170+v184))) = v188
														v190 = *(*int32)(unsafe.Add(mBase, uint32(v173)))
														*(*int32)(unsafe.Add(mBase, uint32(v170))) = v190
														v193 = v165 + int32(-16)
														if base.Ui32(int32(3)) < base.Ui32(v193) {
															v165 = v193
															continue
														} else {
															break
														}
														break
													}
													v198 = v193
												}
											}
											if v198 == int32(0) {
											} else {
												v206 = v198 & int32(3)
												if v206 == int32(0) {
													v231 = v198
												} else {
													v209 = int32(-1)
													v216 = v198
													v217 = v206
													for {
														v222 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v209+v216))))
														*(*uint8)(unsafe.Add(mBase, uint32(l0+v209+v216))) = uint8(v222)
														v224 = int32(-1)
														v225 = v216 + v224
														v227 = v217 + v224
														if v227 != 0 {
															v216 = v225
															v217 = v227
															continue
														} else {
															break
														}
														break
													}
													v231 = v225
												}
												if base.Ui32(v198) < base.Ui32(int32(4)) {
												} else {
													v237 = int32(-4)
													v244 = v231
													for {
														v248 = l0 + v237 + v244
														v249 = int32(3)
														v251 = l1 + v237 + v244
														v254 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v249))))
														*(*uint8)(unsafe.Add(mBase, uint32(v248+v249))) = uint8(v254)
														v256 = int32(2)
														v260 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v256))))
														*(*uint8)(unsafe.Add(mBase, uint32(v248+v256))) = uint8(v260)
														v262 = int32(1)
														v266 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v262))))
														*(*uint8)(unsafe.Add(mBase, uint32(v248+v262))) = uint8(v266)
														v268 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251))))
														*(*uint8)(unsafe.Add(mBase, uint32(v248))) = uint8(v268)
														v271 = v244 + int32(-4)
														if v271 != 0 {
															v244 = v271
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
								} else {
									v116 = v81
									if base.Ui32(v116) < base.Ui32(int32(4)) {
										v198 = v116
									} else {
										v122 = v116 + int32(-4)
										v128 = (int32(base.Ui32(v122)>>(uint(int32(2))%32)) + int32(1)) & int32(3)
										if v128 == int32(0) {
											v152 = v116
										} else {
											v131 = int32(-4)
											v137 = v116
											v138 = v128
											for {
												v144 = *(*int32)(unsafe.Add(mBase, uint32(l1+v131+v137)))
												*(*int32)(unsafe.Add(mBase, uint32(l0+v131+v137))) = v144
												v147 = v137 + int32(-4)
												v149 = v138 + int32(-1)
												if v149 != 0 {
													v137 = v147
													v138 = v149
													continue
												} else {
													break
												}
												break
											}
											v152 = v147
										}
										if base.Ui32(v122) < base.Ui32(int32(12)) {
											v198 = v152
										} else {
											v159 = int32(-16)
											v165 = v152
											for {
												v170 = l0 + v159 + v165
												v171 = int32(12)
												v173 = l1 + v159 + v165
												v176 = *(*int32)(unsafe.Add(mBase, uint32(v173+v171)))
												*(*int32)(unsafe.Add(mBase, uint32(v170+v171))) = v176
												v178 = int32(8)
												v182 = *(*int32)(unsafe.Add(mBase, uint32(v173+v178)))
												*(*int32)(unsafe.Add(mBase, uint32(v170+v178))) = v182
												v184 = int32(4)
												v188 = *(*int32)(unsafe.Add(mBase, uint32(v173+v184)))
												*(*int32)(unsafe.Add(mBase, uint32(v170+v184))) = v188
												v190 = *(*int32)(unsafe.Add(mBase, uint32(v173)))
												*(*int32)(unsafe.Add(mBase, uint32(v170))) = v190
												v193 = v165 + int32(-16)
												if base.Ui32(int32(3)) < base.Ui32(v193) {
													v165 = v193
													continue
												} else {
													break
												}
												break
											}
											v198 = v193
										}
									}
									if v198 == int32(0) {
									} else {
										v206 = v198 & int32(3)
										if v206 == int32(0) {
											v231 = v198
										} else {
											v209 = int32(-1)
											v216 = v198
											v217 = v206
											for {
												v222 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v209+v216))))
												*(*uint8)(unsafe.Add(mBase, uint32(l0+v209+v216))) = uint8(v222)
												v224 = int32(-1)
												v225 = v216 + v224
												v227 = v217 + v224
												if v227 != 0 {
													v216 = v225
													v217 = v227
													continue
												} else {
													break
												}
												break
											}
											v231 = v225
										}
										if base.Ui32(v198) < base.Ui32(int32(4)) {
										} else {
											v237 = int32(-4)
											v244 = v231
											for {
												v248 = l0 + v237 + v244
												v249 = int32(3)
												v251 = l1 + v237 + v244
												v254 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v249))))
												*(*uint8)(unsafe.Add(mBase, uint32(v248+v249))) = uint8(v254)
												v256 = int32(2)
												v260 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v256))))
												*(*uint8)(unsafe.Add(mBase, uint32(v248+v256))) = uint8(v260)
												v262 = int32(1)
												v266 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251+v262))))
												*(*uint8)(unsafe.Add(mBase, uint32(v248+v262))) = uint8(v266)
												v268 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v251))))
												*(*uint8)(unsafe.Add(mBase, uint32(v248))) = uint8(v268)
												v271 = v244 + int32(-4)
												if v271 != 0 {
													v244 = v271
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
					if v22 == int32(0) {
						if l0&int32(3) != 0 {
							if l2 == int32(0) {
							} else {
								v30 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
								*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v30)
								v33 = l2 + int32(-1)
								v35 = l0 + int32(1)
								if v35&int32(3) != 0 {
									if v33 == int32(0) {
									} else {
										v42 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+1)))
										*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v42)
										v45 = l2 + int32(-2)
										v47 = l0 + int32(2)
										if v47&int32(3) != 0 {
											if v45 == int32(0) {
											} else {
												v54 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+2)))
												*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v54)
												v57 = l2 + int32(-3)
												v58 = int32(3)
												v59 = l0 + v58
												if v59&v58 != 0 {
													if v57 == int32(0) {
													} else {
														v66 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3)))
														*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)) = uint8(v66)
														v68 = int32(4)
														v272 = l1 + v68
														v273 = l0 + v68
														v274 = l2 + int32(-4)
														if base.Ui32(v274) < base.Ui32(int32(4)) {
															v346 = v272
															v348 = v273
															v350 = v274
														} else {
															v278 = v274 + int32(-4)
															v284 = (int32(base.Ui32(v278)>>(uint(int32(2))%32)) + int32(1)) & int32(7)
															if v284 == int32(0) {
																v306 = v272
																v308 = v273
																v310 = v274
															} else {
																v291 = v272
																v292 = v284
																v293 = v273
																for {
																	v297 = *(*int32)(unsafe.Add(mBase, uint32(v291)))
																	*(*int32)(unsafe.Add(mBase, uint32(v293))) = v297
																	v299 = int32(4)
																	v300 = v291 + v299
																	v302 = v293 + v299
																	v304 = v292 + int32(-1)
																	if v304 != 0 {
																		v291 = v300
																		v292 = v304
																		v293 = v302
																		continue
																	} else {
																		break
																	}
																	break
																}
																v306 = v300
																v308 = v302
																v310 = v274 - v284<<(uint(int32(2))%32)
															}
															if base.Ui32(v278) < base.Ui32(int32(28)) {
																v346 = v306
																v348 = v308
																v350 = v310
															} else {
																v315 = v306
																v317 = v308
																v319 = v310
																for {
																	v321 = *(*int32)(unsafe.Add(mBase, uint32(v315)))
																	*(*int32)(unsafe.Add(mBase, uint32(v317))) = v321
																	v323 = *(*int32)(unsafe.Add(mBase, uint32(v315)+4))
																	*(*int32)(unsafe.Add(mBase, uint32(v317)+4)) = v323
																	v325 = *(*int32)(unsafe.Add(mBase, uint32(v315)+8))
																	*(*int32)(unsafe.Add(mBase, uint32(v317)+8)) = v325
																	v327 = *(*int32)(unsafe.Add(mBase, uint32(v315)+12))
																	*(*int32)(unsafe.Add(mBase, uint32(v317)+12)) = v327
																	v329 = *(*int32)(unsafe.Add(mBase, uint32(v315)+16))
																	*(*int32)(unsafe.Add(mBase, uint32(v317)+16)) = v329
																	v331 = *(*int32)(unsafe.Add(mBase, uint32(v315)+20))
																	*(*int32)(unsafe.Add(mBase, uint32(v317)+20)) = v331
																	v333 = *(*int32)(unsafe.Add(mBase, uint32(v315)+24))
																	*(*int32)(unsafe.Add(mBase, uint32(v317)+24)) = v333
																	v335 = *(*int32)(unsafe.Add(mBase, uint32(v315)+28))
																	*(*int32)(unsafe.Add(mBase, uint32(v317)+28)) = v335
																	v337 = int32(32)
																	v338 = v315 + v337
																	v340 = v317 + v337
																	v342 = v319 + int32(-32)
																	if base.Ui32(int32(3)) < base.Ui32(v342) {
																		v315 = v338
																		v317 = v340
																		v319 = v342
																		continue
																	} else {
																		break
																	}
																	break
																}
																v346 = v338
																v348 = v340
																v350 = v342
															}
														}
														if v350 == int32(0) {
														} else {
															v355 = v350 & int32(7)
															if v355 != 0 {
																v359 = v346
																v360 = v355
																v361 = v348
																for {
																	v365 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v359))))
																	*(*uint8)(unsafe.Add(mBase, uint32(v361))) = uint8(v365)
																	v367 = int32(1)
																	v368 = v361 + v367
																	v370 = v359 + v367
																	v372 = v360 + int32(-1)
																	if v372 != 0 {
																		v359 = v370
																		v360 = v372
																		v361 = v368
																		continue
																	} else {
																		break
																	}
																	break
																}
																v374 = v370
																v376 = v368
																v377 = v350 & int32(-8)
															} else {
																v374 = v346
																v376 = v348
																v377 = v350
															}
															if base.Ui32(v350) < base.Ui32(int32(8)) {
															} else {
																v383 = v374
																v385 = v376
																v386 = v377
																for {
																	v389 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383))))
																	*(*uint8)(unsafe.Add(mBase, uint32(v385))) = uint8(v389)
																	v391 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+1)))
																	*(*uint8)(unsafe.Add(mBase, uint32(v385)+1)) = uint8(v391)
																	v393 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+2)))
																	*(*uint8)(unsafe.Add(mBase, uint32(v385)+2)) = uint8(v393)
																	v395 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+3)))
																	*(*uint8)(unsafe.Add(mBase, uint32(v385)+3)) = uint8(v395)
																	v397 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+4)))
																	*(*uint8)(unsafe.Add(mBase, uint32(v385)+4)) = uint8(v397)
																	v399 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+5)))
																	*(*uint8)(unsafe.Add(mBase, uint32(v385)+5)) = uint8(v399)
																	v401 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+6)))
																	*(*uint8)(unsafe.Add(mBase, uint32(v385)+6)) = uint8(v401)
																	v403 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+7)))
																	*(*uint8)(unsafe.Add(mBase, uint32(v385)+7)) = uint8(v403)
																	v405 = int32(8)
																	v410 = v386 + int32(-8)
																	if v410 != 0 {
																		v383 = v383 + v405
																		v385 = v385 + v405
																		v386 = v410
																		continue
																	} else {
																		break
																	}
																	break
																}
															}
														}
													}
												} else {
													v272 = l1 + int32(3)
													v273 = v59
													v274 = v57
													if base.Ui32(v274) < base.Ui32(int32(4)) {
														v346 = v272
														v348 = v273
														v350 = v274
													} else {
														v278 = v274 + int32(-4)
														v284 = (int32(base.Ui32(v278)>>(uint(int32(2))%32)) + int32(1)) & int32(7)
														if v284 == int32(0) {
															v306 = v272
															v308 = v273
															v310 = v274
														} else {
															v291 = v272
															v292 = v284
															v293 = v273
															for {
																v297 = *(*int32)(unsafe.Add(mBase, uint32(v291)))
																*(*int32)(unsafe.Add(mBase, uint32(v293))) = v297
																v299 = int32(4)
																v300 = v291 + v299
																v302 = v293 + v299
																v304 = v292 + int32(-1)
																if v304 != 0 {
																	v291 = v300
																	v292 = v304
																	v293 = v302
																	continue
																} else {
																	break
																}
																break
															}
															v306 = v300
															v308 = v302
															v310 = v274 - v284<<(uint(int32(2))%32)
														}
														if base.Ui32(v278) < base.Ui32(int32(28)) {
															v346 = v306
															v348 = v308
															v350 = v310
														} else {
															v315 = v306
															v317 = v308
															v319 = v310
															for {
																v321 = *(*int32)(unsafe.Add(mBase, uint32(v315)))
																*(*int32)(unsafe.Add(mBase, uint32(v317))) = v321
																v323 = *(*int32)(unsafe.Add(mBase, uint32(v315)+4))
																*(*int32)(unsafe.Add(mBase, uint32(v317)+4)) = v323
																v325 = *(*int32)(unsafe.Add(mBase, uint32(v315)+8))
																*(*int32)(unsafe.Add(mBase, uint32(v317)+8)) = v325
																v327 = *(*int32)(unsafe.Add(mBase, uint32(v315)+12))
																*(*int32)(unsafe.Add(mBase, uint32(v317)+12)) = v327
																v329 = *(*int32)(unsafe.Add(mBase, uint32(v315)+16))
																*(*int32)(unsafe.Add(mBase, uint32(v317)+16)) = v329
																v331 = *(*int32)(unsafe.Add(mBase, uint32(v315)+20))
																*(*int32)(unsafe.Add(mBase, uint32(v317)+20)) = v331
																v333 = *(*int32)(unsafe.Add(mBase, uint32(v315)+24))
																*(*int32)(unsafe.Add(mBase, uint32(v317)+24)) = v333
																v335 = *(*int32)(unsafe.Add(mBase, uint32(v315)+28))
																*(*int32)(unsafe.Add(mBase, uint32(v317)+28)) = v335
																v337 = int32(32)
																v338 = v315 + v337
																v340 = v317 + v337
																v342 = v319 + int32(-32)
																if base.Ui32(int32(3)) < base.Ui32(v342) {
																	v315 = v338
																	v317 = v340
																	v319 = v342
																	continue
																} else {
																	break
																}
																break
															}
															v346 = v338
															v348 = v340
															v350 = v342
														}
													}
													if v350 == int32(0) {
													} else {
														v355 = v350 & int32(7)
														if v355 != 0 {
															v359 = v346
															v360 = v355
															v361 = v348
															for {
																v365 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v359))))
																*(*uint8)(unsafe.Add(mBase, uint32(v361))) = uint8(v365)
																v367 = int32(1)
																v368 = v361 + v367
																v370 = v359 + v367
																v372 = v360 + int32(-1)
																if v372 != 0 {
																	v359 = v370
																	v360 = v372
																	v361 = v368
																	continue
																} else {
																	break
																}
																break
															}
															v374 = v370
															v376 = v368
															v377 = v350 & int32(-8)
														} else {
															v374 = v346
															v376 = v348
															v377 = v350
														}
														if base.Ui32(v350) < base.Ui32(int32(8)) {
														} else {
															v383 = v374
															v385 = v376
															v386 = v377
															for {
																v389 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383))))
																*(*uint8)(unsafe.Add(mBase, uint32(v385))) = uint8(v389)
																v391 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+1)))
																*(*uint8)(unsafe.Add(mBase, uint32(v385)+1)) = uint8(v391)
																v393 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+2)))
																*(*uint8)(unsafe.Add(mBase, uint32(v385)+2)) = uint8(v393)
																v395 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+3)))
																*(*uint8)(unsafe.Add(mBase, uint32(v385)+3)) = uint8(v395)
																v397 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+4)))
																*(*uint8)(unsafe.Add(mBase, uint32(v385)+4)) = uint8(v397)
																v399 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+5)))
																*(*uint8)(unsafe.Add(mBase, uint32(v385)+5)) = uint8(v399)
																v401 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+6)))
																*(*uint8)(unsafe.Add(mBase, uint32(v385)+6)) = uint8(v401)
																v403 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+7)))
																*(*uint8)(unsafe.Add(mBase, uint32(v385)+7)) = uint8(v403)
																v405 = int32(8)
																v410 = v386 + int32(-8)
																if v410 != 0 {
																	v383 = v383 + v405
																	v385 = v385 + v405
																	v386 = v410
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
										} else {
											v272 = l1 + int32(2)
											v273 = v47
											v274 = v45
											if base.Ui32(v274) < base.Ui32(int32(4)) {
												v346 = v272
												v348 = v273
												v350 = v274
											} else {
												v278 = v274 + int32(-4)
												v284 = (int32(base.Ui32(v278)>>(uint(int32(2))%32)) + int32(1)) & int32(7)
												if v284 == int32(0) {
													v306 = v272
													v308 = v273
													v310 = v274
												} else {
													v291 = v272
													v292 = v284
													v293 = v273
													for {
														v297 = *(*int32)(unsafe.Add(mBase, uint32(v291)))
														*(*int32)(unsafe.Add(mBase, uint32(v293))) = v297
														v299 = int32(4)
														v300 = v291 + v299
														v302 = v293 + v299
														v304 = v292 + int32(-1)
														if v304 != 0 {
															v291 = v300
															v292 = v304
															v293 = v302
															continue
														} else {
															break
														}
														break
													}
													v306 = v300
													v308 = v302
													v310 = v274 - v284<<(uint(int32(2))%32)
												}
												if base.Ui32(v278) < base.Ui32(int32(28)) {
													v346 = v306
													v348 = v308
													v350 = v310
												} else {
													v315 = v306
													v317 = v308
													v319 = v310
													for {
														v321 = *(*int32)(unsafe.Add(mBase, uint32(v315)))
														*(*int32)(unsafe.Add(mBase, uint32(v317))) = v321
														v323 = *(*int32)(unsafe.Add(mBase, uint32(v315)+4))
														*(*int32)(unsafe.Add(mBase, uint32(v317)+4)) = v323
														v325 = *(*int32)(unsafe.Add(mBase, uint32(v315)+8))
														*(*int32)(unsafe.Add(mBase, uint32(v317)+8)) = v325
														v327 = *(*int32)(unsafe.Add(mBase, uint32(v315)+12))
														*(*int32)(unsafe.Add(mBase, uint32(v317)+12)) = v327
														v329 = *(*int32)(unsafe.Add(mBase, uint32(v315)+16))
														*(*int32)(unsafe.Add(mBase, uint32(v317)+16)) = v329
														v331 = *(*int32)(unsafe.Add(mBase, uint32(v315)+20))
														*(*int32)(unsafe.Add(mBase, uint32(v317)+20)) = v331
														v333 = *(*int32)(unsafe.Add(mBase, uint32(v315)+24))
														*(*int32)(unsafe.Add(mBase, uint32(v317)+24)) = v333
														v335 = *(*int32)(unsafe.Add(mBase, uint32(v315)+28))
														*(*int32)(unsafe.Add(mBase, uint32(v317)+28)) = v335
														v337 = int32(32)
														v338 = v315 + v337
														v340 = v317 + v337
														v342 = v319 + int32(-32)
														if base.Ui32(int32(3)) < base.Ui32(v342) {
															v315 = v338
															v317 = v340
															v319 = v342
															continue
														} else {
															break
														}
														break
													}
													v346 = v338
													v348 = v340
													v350 = v342
												}
											}
											if v350 == int32(0) {
											} else {
												v355 = v350 & int32(7)
												if v355 != 0 {
													v359 = v346
													v360 = v355
													v361 = v348
													for {
														v365 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v359))))
														*(*uint8)(unsafe.Add(mBase, uint32(v361))) = uint8(v365)
														v367 = int32(1)
														v368 = v361 + v367
														v370 = v359 + v367
														v372 = v360 + int32(-1)
														if v372 != 0 {
															v359 = v370
															v360 = v372
															v361 = v368
															continue
														} else {
															break
														}
														break
													}
													v374 = v370
													v376 = v368
													v377 = v350 & int32(-8)
												} else {
													v374 = v346
													v376 = v348
													v377 = v350
												}
												if base.Ui32(v350) < base.Ui32(int32(8)) {
												} else {
													v383 = v374
													v385 = v376
													v386 = v377
													for {
														v389 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383))))
														*(*uint8)(unsafe.Add(mBase, uint32(v385))) = uint8(v389)
														v391 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+1)))
														*(*uint8)(unsafe.Add(mBase, uint32(v385)+1)) = uint8(v391)
														v393 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+2)))
														*(*uint8)(unsafe.Add(mBase, uint32(v385)+2)) = uint8(v393)
														v395 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+3)))
														*(*uint8)(unsafe.Add(mBase, uint32(v385)+3)) = uint8(v395)
														v397 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+4)))
														*(*uint8)(unsafe.Add(mBase, uint32(v385)+4)) = uint8(v397)
														v399 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+5)))
														*(*uint8)(unsafe.Add(mBase, uint32(v385)+5)) = uint8(v399)
														v401 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+6)))
														*(*uint8)(unsafe.Add(mBase, uint32(v385)+6)) = uint8(v401)
														v403 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+7)))
														*(*uint8)(unsafe.Add(mBase, uint32(v385)+7)) = uint8(v403)
														v405 = int32(8)
														v410 = v386 + int32(-8)
														if v410 != 0 {
															v383 = v383 + v405
															v385 = v385 + v405
															v386 = v410
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
								} else {
									v272 = l1 + int32(1)
									v273 = v35
									v274 = v33
									if base.Ui32(v274) < base.Ui32(int32(4)) {
										v346 = v272
										v348 = v273
										v350 = v274
									} else {
										v278 = v274 + int32(-4)
										v284 = (int32(base.Ui32(v278)>>(uint(int32(2))%32)) + int32(1)) & int32(7)
										if v284 == int32(0) {
											v306 = v272
											v308 = v273
											v310 = v274
										} else {
											v291 = v272
											v292 = v284
											v293 = v273
											for {
												v297 = *(*int32)(unsafe.Add(mBase, uint32(v291)))
												*(*int32)(unsafe.Add(mBase, uint32(v293))) = v297
												v299 = int32(4)
												v300 = v291 + v299
												v302 = v293 + v299
												v304 = v292 + int32(-1)
												if v304 != 0 {
													v291 = v300
													v292 = v304
													v293 = v302
													continue
												} else {
													break
												}
												break
											}
											v306 = v300
											v308 = v302
											v310 = v274 - v284<<(uint(int32(2))%32)
										}
										if base.Ui32(v278) < base.Ui32(int32(28)) {
											v346 = v306
											v348 = v308
											v350 = v310
										} else {
											v315 = v306
											v317 = v308
											v319 = v310
											for {
												v321 = *(*int32)(unsafe.Add(mBase, uint32(v315)))
												*(*int32)(unsafe.Add(mBase, uint32(v317))) = v321
												v323 = *(*int32)(unsafe.Add(mBase, uint32(v315)+4))
												*(*int32)(unsafe.Add(mBase, uint32(v317)+4)) = v323
												v325 = *(*int32)(unsafe.Add(mBase, uint32(v315)+8))
												*(*int32)(unsafe.Add(mBase, uint32(v317)+8)) = v325
												v327 = *(*int32)(unsafe.Add(mBase, uint32(v315)+12))
												*(*int32)(unsafe.Add(mBase, uint32(v317)+12)) = v327
												v329 = *(*int32)(unsafe.Add(mBase, uint32(v315)+16))
												*(*int32)(unsafe.Add(mBase, uint32(v317)+16)) = v329
												v331 = *(*int32)(unsafe.Add(mBase, uint32(v315)+20))
												*(*int32)(unsafe.Add(mBase, uint32(v317)+20)) = v331
												v333 = *(*int32)(unsafe.Add(mBase, uint32(v315)+24))
												*(*int32)(unsafe.Add(mBase, uint32(v317)+24)) = v333
												v335 = *(*int32)(unsafe.Add(mBase, uint32(v315)+28))
												*(*int32)(unsafe.Add(mBase, uint32(v317)+28)) = v335
												v337 = int32(32)
												v338 = v315 + v337
												v340 = v317 + v337
												v342 = v319 + int32(-32)
												if base.Ui32(int32(3)) < base.Ui32(v342) {
													v315 = v338
													v317 = v340
													v319 = v342
													continue
												} else {
													break
												}
												break
											}
											v346 = v338
											v348 = v340
											v350 = v342
										}
									}
									if v350 == int32(0) {
									} else {
										v355 = v350 & int32(7)
										if v355 != 0 {
											v359 = v346
											v360 = v355
											v361 = v348
											for {
												v365 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v359))))
												*(*uint8)(unsafe.Add(mBase, uint32(v361))) = uint8(v365)
												v367 = int32(1)
												v368 = v361 + v367
												v370 = v359 + v367
												v372 = v360 + int32(-1)
												if v372 != 0 {
													v359 = v370
													v360 = v372
													v361 = v368
													continue
												} else {
													break
												}
												break
											}
											v374 = v370
											v376 = v368
											v377 = v350 & int32(-8)
										} else {
											v374 = v346
											v376 = v348
											v377 = v350
										}
										if base.Ui32(v350) < base.Ui32(int32(8)) {
										} else {
											v383 = v374
											v385 = v376
											v386 = v377
											for {
												v389 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383))))
												*(*uint8)(unsafe.Add(mBase, uint32(v385))) = uint8(v389)
												v391 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+1)))
												*(*uint8)(unsafe.Add(mBase, uint32(v385)+1)) = uint8(v391)
												v393 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+2)))
												*(*uint8)(unsafe.Add(mBase, uint32(v385)+2)) = uint8(v393)
												v395 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+3)))
												*(*uint8)(unsafe.Add(mBase, uint32(v385)+3)) = uint8(v395)
												v397 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+4)))
												*(*uint8)(unsafe.Add(mBase, uint32(v385)+4)) = uint8(v397)
												v399 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+5)))
												*(*uint8)(unsafe.Add(mBase, uint32(v385)+5)) = uint8(v399)
												v401 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+6)))
												*(*uint8)(unsafe.Add(mBase, uint32(v385)+6)) = uint8(v401)
												v403 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+7)))
												*(*uint8)(unsafe.Add(mBase, uint32(v385)+7)) = uint8(v403)
												v405 = int32(8)
												v410 = v386 + int32(-8)
												if v410 != 0 {
													v383 = v383 + v405
													v385 = v385 + v405
													v386 = v410
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
						} else {
							v272 = l1
							v273 = l0
							v274 = l2
							if base.Ui32(v274) < base.Ui32(int32(4)) {
								v346 = v272
								v348 = v273
								v350 = v274
							} else {
								v278 = v274 + int32(-4)
								v284 = (int32(base.Ui32(v278)>>(uint(int32(2))%32)) + int32(1)) & int32(7)
								if v284 == int32(0) {
									v306 = v272
									v308 = v273
									v310 = v274
								} else {
									v291 = v272
									v292 = v284
									v293 = v273
									for {
										v297 = *(*int32)(unsafe.Add(mBase, uint32(v291)))
										*(*int32)(unsafe.Add(mBase, uint32(v293))) = v297
										v299 = int32(4)
										v300 = v291 + v299
										v302 = v293 + v299
										v304 = v292 + int32(-1)
										if v304 != 0 {
											v291 = v300
											v292 = v304
											v293 = v302
											continue
										} else {
											break
										}
										break
									}
									v306 = v300
									v308 = v302
									v310 = v274 - v284<<(uint(int32(2))%32)
								}
								if base.Ui32(v278) < base.Ui32(int32(28)) {
									v346 = v306
									v348 = v308
									v350 = v310
								} else {
									v315 = v306
									v317 = v308
									v319 = v310
									for {
										v321 = *(*int32)(unsafe.Add(mBase, uint32(v315)))
										*(*int32)(unsafe.Add(mBase, uint32(v317))) = v321
										v323 = *(*int32)(unsafe.Add(mBase, uint32(v315)+4))
										*(*int32)(unsafe.Add(mBase, uint32(v317)+4)) = v323
										v325 = *(*int32)(unsafe.Add(mBase, uint32(v315)+8))
										*(*int32)(unsafe.Add(mBase, uint32(v317)+8)) = v325
										v327 = *(*int32)(unsafe.Add(mBase, uint32(v315)+12))
										*(*int32)(unsafe.Add(mBase, uint32(v317)+12)) = v327
										v329 = *(*int32)(unsafe.Add(mBase, uint32(v315)+16))
										*(*int32)(unsafe.Add(mBase, uint32(v317)+16)) = v329
										v331 = *(*int32)(unsafe.Add(mBase, uint32(v315)+20))
										*(*int32)(unsafe.Add(mBase, uint32(v317)+20)) = v331
										v333 = *(*int32)(unsafe.Add(mBase, uint32(v315)+24))
										*(*int32)(unsafe.Add(mBase, uint32(v317)+24)) = v333
										v335 = *(*int32)(unsafe.Add(mBase, uint32(v315)+28))
										*(*int32)(unsafe.Add(mBase, uint32(v317)+28)) = v335
										v337 = int32(32)
										v338 = v315 + v337
										v340 = v317 + v337
										v342 = v319 + int32(-32)
										if base.Ui32(int32(3)) < base.Ui32(v342) {
											v315 = v338
											v317 = v340
											v319 = v342
											continue
										} else {
											break
										}
										break
									}
									v346 = v338
									v348 = v340
									v350 = v342
								}
							}
							if v350 == int32(0) {
							} else {
								v355 = v350 & int32(7)
								if v355 != 0 {
									v359 = v346
									v360 = v355
									v361 = v348
									for {
										v365 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v359))))
										*(*uint8)(unsafe.Add(mBase, uint32(v361))) = uint8(v365)
										v367 = int32(1)
										v368 = v361 + v367
										v370 = v359 + v367
										v372 = v360 + int32(-1)
										if v372 != 0 {
											v359 = v370
											v360 = v372
											v361 = v368
											continue
										} else {
											break
										}
										break
									}
									v374 = v370
									v376 = v368
									v377 = v350 & int32(-8)
								} else {
									v374 = v346
									v376 = v348
									v377 = v350
								}
								if base.Ui32(v350) < base.Ui32(int32(8)) {
								} else {
									v383 = v374
									v385 = v376
									v386 = v377
									for {
										v389 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383))))
										*(*uint8)(unsafe.Add(mBase, uint32(v385))) = uint8(v389)
										v391 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+1)))
										*(*uint8)(unsafe.Add(mBase, uint32(v385)+1)) = uint8(v391)
										v393 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+2)))
										*(*uint8)(unsafe.Add(mBase, uint32(v385)+2)) = uint8(v393)
										v395 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+3)))
										*(*uint8)(unsafe.Add(mBase, uint32(v385)+3)) = uint8(v395)
										v397 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+4)))
										*(*uint8)(unsafe.Add(mBase, uint32(v385)+4)) = uint8(v397)
										v399 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+5)))
										*(*uint8)(unsafe.Add(mBase, uint32(v385)+5)) = uint8(v399)
										v401 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+6)))
										*(*uint8)(unsafe.Add(mBase, uint32(v385)+6)) = uint8(v401)
										v403 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+7)))
										*(*uint8)(unsafe.Add(mBase, uint32(v385)+7)) = uint8(v403)
										v405 = int32(8)
										v410 = v386 + int32(-8)
										if v410 != 0 {
											v383 = v383 + v405
											v385 = v385 + v405
											v386 = v410
											continue
										} else {
											break
										}
										break
									}
								}
							}
						}
					} else {
						v346 = l1
						v348 = l0
						v350 = l2
						if v350 == int32(0) {
						} else {
							v355 = v350 & int32(7)
							if v355 != 0 {
								v359 = v346
								v360 = v355
								v361 = v348
								for {
									v365 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v359))))
									*(*uint8)(unsafe.Add(mBase, uint32(v361))) = uint8(v365)
									v367 = int32(1)
									v368 = v361 + v367
									v370 = v359 + v367
									v372 = v360 + int32(-1)
									if v372 != 0 {
										v359 = v370
										v360 = v372
										v361 = v368
										continue
									} else {
										break
									}
									break
								}
								v374 = v370
								v376 = v368
								v377 = v350 & int32(-8)
							} else {
								v374 = v346
								v376 = v348
								v377 = v350
							}
							if base.Ui32(v350) < base.Ui32(int32(8)) {
							} else {
								v383 = v374
								v385 = v376
								v386 = v377
								for {
									v389 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383))))
									*(*uint8)(unsafe.Add(mBase, uint32(v385))) = uint8(v389)
									v391 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+1)))
									*(*uint8)(unsafe.Add(mBase, uint32(v385)+1)) = uint8(v391)
									v393 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+2)))
									*(*uint8)(unsafe.Add(mBase, uint32(v385)+2)) = uint8(v393)
									v395 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+3)))
									*(*uint8)(unsafe.Add(mBase, uint32(v385)+3)) = uint8(v395)
									v397 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+4)))
									*(*uint8)(unsafe.Add(mBase, uint32(v385)+4)) = uint8(v397)
									v399 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+5)))
									*(*uint8)(unsafe.Add(mBase, uint32(v385)+5)) = uint8(v399)
									v401 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+6)))
									*(*uint8)(unsafe.Add(mBase, uint32(v385)+6)) = uint8(v401)
									v403 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v383)+7)))
									*(*uint8)(unsafe.Add(mBase, uint32(v385)+7)) = uint8(v403)
									v405 = int32(8)
									v410 = v386 + int32(-8)
									if v410 != 0 {
										v383 = v383 + v405
										v385 = v385 + v405
										v386 = v410
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
			} else {
				base.MemoryCopy(m, l0, l1, l2)
			}
		}
	}
	return l0
}
func F_memset(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v15 int32
	_ = v15
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v45 int32
	_ = v45
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v85 int32
	_ = v85
	var v86 int32
	_ = v86
	var v91 int64
	_ = v91
	var v94 int32
	_ = v94
	var v95 int32
	_ = v95
	var v107 int32
	_ = v107
	v2 = l1
	if base.Ui32(l2) < base.Ui32(int32(33)) {
		if l2 == int32(0) {
		} else {
			*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v2)
			v15 = l0 + l2
			*(*uint8)(unsafe.Add(mBase, uint32(v15+int32(-1)))) = uint8(v2)
			if base.Ui32(l2) < base.Ui32(int32(3)) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v2)
				*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v2)
				*(*uint8)(unsafe.Add(mBase, uint32(v15+int32(-3)))) = uint8(v2)
				*(*uint8)(unsafe.Add(mBase, uint32(v15+int32(-2)))) = uint8(v2)
				if base.Ui32(l2) < base.Ui32(int32(7)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)) = uint8(v2)
					*(*uint8)(unsafe.Add(mBase, uint32(v15+int32(-4)))) = uint8(v2)
					if base.Ui32(l2) < base.Ui32(int32(9)) {
					} else {
						v40 = (int32(0) - l0) & int32(3)
						v41 = l0 + v40
						v45 = v2 & int32(255) * int32(16843009)
						*(*int32)(unsafe.Add(mBase, uint32(v41))) = v45
						v49 = (l2 - v40) & int32(60)
						v50 = v41 + v49
						*(*int32)(unsafe.Add(mBase, uint32(v50+int32(-4)))) = v45
						if base.Ui32(v49) < base.Ui32(int32(9)) {
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(v41)+8)) = v45
							*(*int32)(unsafe.Add(mBase, uint32(v41)+4)) = v45
							*(*int32)(unsafe.Add(mBase, uint32(v50+int32(-8)))) = v45
							*(*int32)(unsafe.Add(mBase, uint32(v50+int32(-12)))) = v45
							if base.Ui32(v49) < base.Ui32(int32(25)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v41)+24)) = v45
								*(*int32)(unsafe.Add(mBase, uint32(v41)+20)) = v45
								*(*int32)(unsafe.Add(mBase, uint32(v41)+16)) = v45
								*(*int32)(unsafe.Add(mBase, uint32(v41)+12)) = v45
								*(*int32)(unsafe.Add(mBase, uint32(v50+int32(-16)))) = v45
								*(*int32)(unsafe.Add(mBase, uint32(v50+int32(-20)))) = v45
								*(*int32)(unsafe.Add(mBase, uint32(v50+int32(-24)))) = v45
								*(*int32)(unsafe.Add(mBase, uint32(v50+int32(-28)))) = v45
								v85 = v41&int32(4) | int32(24)
								v86 = v49 - v85
								if base.Ui32(v86) < base.Ui32(int32(32)) {
								} else {
									v91 = base.I64_extend_i32_u(v45) * int64(4294967297)
									v94 = v86
									v95 = v41 + v85
									for {
										*(*int64)(unsafe.Add(mBase, uint32(v95)+24)) = v91
										*(*int64)(unsafe.Add(mBase, uint32(v95)+16)) = v91
										*(*int64)(unsafe.Add(mBase, uint32(v95)+8)) = v91
										*(*int64)(unsafe.Add(mBase, uint32(v95))) = v91
										v107 = v94 + int32(-32)
										if base.Ui32(int32(31)) < base.Ui32(v107) {
											v94 = v107
											v95 = v95 + int32(32)
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
		}
		return l0
	} else {
		base.MemoryFill(m, l0, v2, l2)
		return l0
	}
}
