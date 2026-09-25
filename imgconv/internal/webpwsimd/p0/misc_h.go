//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_HD4_C(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v19 int32
	_ = v19
	var v24 int32
	_ = v24
	var v29 int32
	_ = v29
	var v32 int32
	_ = v32
	var v37 int32
	_ = v37
	var v43 int32
	_ = v43
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v54 int32
	_ = v54
	var v57 int32
	_ = v57
	var v61 int32
	_ = v61
	var v67 int32
	_ = v67
	var v72 int32
	_ = v72
	var v81 int32
	_ = v81
	var v84 int32
	_ = v84
	var v90 int32
	_ = v90
	var v98 int32
	_ = v98
	var v106 int32
	_ = v106
	v11 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+31)))
	v12 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+63)))
	v14 = int32(1)
	v17 = int32(base.Ui32(v11+v12+v14) >> (uint(v14) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+98)) = uint8(v17)
	v19 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+95)))
	v24 = int32(base.Ui32(v12+v19+v14) >> (uint(v14) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+96)) = uint8(v24)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+64)) = uint8(v17)
	v29 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-1)))))
	v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-33)))))
	v37 = int32(base.Ui32(v29+v32+v14) >> (uint(v14) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+34)) = uint8(v37)
	v43 = int32(base.Ui32(v29+v11+v14) >> (uint(v14) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+66)) = uint8(v43)
	*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v37)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+32)) = uint8(v43)
	v47 = int32(2)
	v48 = v29 + v47
	v54 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-32)))))
	v57 = int32(base.Ui32(v48+v32<<(uint(v14)%32)+v54) >> (uint(v47) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+35)) = uint8(v57)
	v61 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-31)))))
	v67 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-30)))))
	v72 = int32(base.Ui32(v54+v61<<(uint(v14)%32)+v67+v47) >> (uint(v47) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)) = uint8(v72)
	v81 = int32(base.Ui32(v61+(v32+v54<<(uint(v14)%32))+v47) >> (uint(v47) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v81)
	v84 = v11 + v47
	v90 = int32(base.Ui32(v32+(v84+v29<<(uint(v14)%32))) >> (uint(v47) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+67)) = uint8(v90)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v57)
	v98 = int32(base.Ui32(v48+v12+v11<<(uint(v14)%32)) >> (uint(v47) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+99)) = uint8(v98)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+33)) = uint8(v90)
	v106 = int32(base.Ui32(v84+v19+v12<<(uint(v14)%32)) >> (uint(v47) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+97)) = uint8(v106)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+65)) = uint8(v98)
	return
}
func F_HE16_C(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v3 int64
	_ = v3
	var v4 int64
	_ = v4
	var v5 int64
	_ = v5
	var v10 int64
	_ = v10
	var v12 int64
	_ = v12
	var v17 int64
	_ = v17
	var v19 int64
	_ = v19
	var v24 int64
	_ = v24
	var v26 int64
	_ = v26
	var v33 int64
	_ = v33
	var v35 int64
	_ = v35
	var v40 int64
	_ = v40
	var v42 int64
	_ = v42
	var v49 int64
	_ = v49
	var v51 int64
	_ = v51
	var v56 int64
	_ = v56
	var v58 int64
	_ = v58
	var v63 int64
	_ = v63
	var v65 int64
	_ = v65
	var v70 int64
	_ = v70
	var v72 int64
	_ = v72
	var v77 int64
	_ = v77
	var v79 int64
	_ = v79
	var v84 int64
	_ = v84
	var v86 int64
	_ = v86
	var v91 int64
	_ = v91
	var v93 int64
	_ = v93
	var v98 int64
	_ = v98
	var v100 int64
	_ = v100
	var v105 int64
	_ = v105
	var v107 int64
	_ = v107
	var v112 int64
	_ = v112
	var v114 int64
	_ = v114
	v3 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+31)))
	v4 = int64(72340172838076673)
	v5 = v3 * v4
	*(*int64)(unsafe.Add(mBase, uint32(l0)+32)) = v5
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(40)))) = v5
	v10 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+63)))
	v12 = v10 * v4
	*(*int64)(unsafe.Add(mBase, uint32(l0)+64)) = v12
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(72)))) = v12
	v17 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+95)))
	v19 = v17 * v4
	*(*int64)(unsafe.Add(mBase, uint32(l0)+96)) = v19
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(104)))) = v19
	v24 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+127)))
	v26 = v24 * v4
	*(*int64)(unsafe.Add(mBase, uint32(l0)+128)) = v26
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(136)))) = v26
	v33 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+159)))
	v35 = v33 * v4
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(168)))) = v35
	*(*int64)(unsafe.Add(mBase, uint32(l0)+160)) = v35
	v40 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-1)))))
	v42 = v40 * v4
	*(*int64)(unsafe.Add(mBase, uint32(l0))) = v42
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(8)))) = v42
	v49 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+191)))
	v51 = v49 * v4
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(200)))) = v51
	*(*int64)(unsafe.Add(mBase, uint32(l0)+192)) = v51
	v56 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+223)))
	v58 = v56 * v4
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(232)))) = v58
	*(*int64)(unsafe.Add(mBase, uint32(l0)+224)) = v58
	v63 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+255)))
	v65 = v63 * v4
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(264)))) = v65
	*(*int64)(unsafe.Add(mBase, uint32(l0)+256)) = v65
	v70 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+287)))
	v72 = v70 * v4
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(296)))) = v72
	*(*int64)(unsafe.Add(mBase, uint32(l0)+288)) = v72
	v77 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+319)))
	v79 = v77 * v4
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(328)))) = v79
	*(*int64)(unsafe.Add(mBase, uint32(l0)+320)) = v79
	v84 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+351)))
	v86 = v84 * v4
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(360)))) = v86
	*(*int64)(unsafe.Add(mBase, uint32(l0)+352)) = v86
	v91 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+383)))
	v93 = v91 * v4
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(392)))) = v93
	*(*int64)(unsafe.Add(mBase, uint32(l0)+384)) = v93
	v98 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+415)))
	v100 = v98 * v4
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(424)))) = v100
	*(*int64)(unsafe.Add(mBase, uint32(l0)+416)) = v100
	v105 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+447)))
	v107 = v105 * v4
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(456)))) = v107
	*(*int64)(unsafe.Add(mBase, uint32(l0)+448)) = v107
	v112 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+479)))
	v114 = v112 * v4
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(488)))) = v114
	*(*int64)(unsafe.Add(mBase, uint32(l0)+480)) = v114
	return
}
func F_HE4_C(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v17 int32
	_ = v17
	var v20 int32
	_ = v20
	var v22 int32
	_ = v22
	var v34 int32
	_ = v34
	var v46 int32
	_ = v46
	v7 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+63)))
	v8 = int32(2)
	v9 = v7 + v8
	v10 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+95)))
	v12 = int32(1)
	v17 = int32(16843009)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+96)) = int32(base.Ui32(v9+v10+v10<<(uint(v12)%32))>>(uint(v8)%32)) * v17
	v20 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+31)))
	v22 = v20 + v8
	*(*int32)(unsafe.Add(mBase, uint32(l0)+64)) = int32(base.Ui32(v10+(v22+v7<<(uint(v12)%32)))>>(uint(v8)%32)) * v17
	v34 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-1)))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = int32(base.Ui32(v9+v34+v20<<(uint(v12)%32))>>(uint(v8)%32)) * v17
	v46 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-33)))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = int32(base.Ui32(v22+v46+v34<<(uint(v12)%32))>>(uint(v8)%32)) * v17
	return
}
func F_HE8uv_C(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int64
	_ = v2
	var v3 int64
	_ = v3
	var v6 int64
	_ = v6
	var v10 int64
	_ = v10
	var v14 int64
	_ = v14
	var v18 int64
	_ = v18
	var v22 int64
	_ = v22
	var v26 int64
	_ = v26
	var v32 int64
	_ = v32
	v2 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+31)))
	v3 = int64(72340172838076673)
	*(*int64)(unsafe.Add(mBase, uint32(l0)+32)) = v2 * v3
	v6 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+63)))
	*(*int64)(unsafe.Add(mBase, uint32(l0)+64)) = v6 * v3
	v10 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+95)))
	*(*int64)(unsafe.Add(mBase, uint32(l0)+96)) = v10 * v3
	v14 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+127)))
	*(*int64)(unsafe.Add(mBase, uint32(l0)+128)) = v14 * v3
	v18 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+159)))
	*(*int64)(unsafe.Add(mBase, uint32(l0)+160)) = v18 * v3
	v22 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+191)))
	*(*int64)(unsafe.Add(mBase, uint32(l0)+192)) = v22 * v3
	v26 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0)+223)))
	*(*int64)(unsafe.Add(mBase, uint32(l0)+224)) = v26 * v3
	v32 = int64(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-1)))))
	*(*int64)(unsafe.Add(mBase, uint32(l0))) = v32 * v3
	return
}
func F_HFilter16_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v35 int32
	_ = v35
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
	v35 = int32(1)
	v55 = m.G54
	v56 = *(*int32)(unsafe.Add(mBase, uint32(v55)))
	v57 = m.G53
	v58 = *(*int32)(unsafe.Add(mBase, uint32(v57)))
	v59 = m.G55
	v60 = *(*int32)(unsafe.Add(mBase, uint32(v59)))
	v61 = m.G56
	v62 = *(*int32)(unsafe.Add(mBase, uint32(v61)))
	v63 = l0
	v67 = int32(17)
	for {
		v95 = v63 + int32(-1)
		v96 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v95))))
		v97 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63))))
		v100 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v96-v97)))))
		v103 = v63 + int32(-2)
		v104 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v103))))
		v105 = v63 + int32(1)
		v106 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v105))))
		v107 = v104 - v106
		v109 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+v107))))
		if l2<<(uint(v35)%32)|v35 < v100<<(uint(int32(2))%32)+v109 {
		} else {
			v113 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+int32(-4)))))
			v114 = v63 + int32(-3)
			v115 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v114))))
			v118 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v113-v115)))))
			if l3 < v118 {
			} else {
				v122 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v115-v104)))))
				if base.Ui32(l3) < base.Ui32(v122) {
				} else {
					v126 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v104-v96)))))
					if base.Ui32(l3) < base.Ui32(v126) {
					} else {
						v129 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+int32(3)))))
						v130 = v63 + int32(2)
						v131 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v130))))
						v134 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v129-v131)))))
						if base.Ui32(l3) < base.Ui32(v134) {
						} else {
							v138 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v131-v106)))))
							if base.Ui32(l3) < base.Ui32(v138) {
							} else {
								v142 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v106-v97)))))
								if base.Ui32(l3) < base.Ui32(v142) {
								} else {
									v146 = (v97 - v96) * int32(3)
									v148 = int32(*(*int8)(unsafe.Add(mBase, uint32(v60+v107))))
									if l4 < v126 {
										v151 = v146 + v148
										v154 = int32(3)
										v157 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+(v151+int32(4))>>(uint(v154)%32)))))
										v164 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+(v151+v154)>>(uint(v154)%32)))))
										v166 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v96+v164))))
										*(*uint8)(unsafe.Add(mBase, uint32(v95))) = uint8(v166)
										v211 = v97 - v157
										v216 = v63
									} else {
										if base.Ui32(v142) <= base.Ui32(l4) {
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
			v63 = v63 + l1
			v67 = v232
			continue
		} else {
			break
		}
		break
	}
	return
}
func F_HFilter16i_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v33 int32
	_ = v33
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
	var v65 int32
	_ = v65
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v94 int32
	_ = v94
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v101 int32
	_ = v101
	var v103 int32
	_ = v103
	var v107 int32
	_ = v107
	var v109 int32
	_ = v109
	var v112 int32
	_ = v112
	var v116 int32
	_ = v116
	var v120 int32
	_ = v120
	var v123 int32
	_ = v123
	var v125 int32
	_ = v125
	var v128 int32
	_ = v128
	var v132 int32
	_ = v132
	var v136 int32
	_ = v136
	var v140 int32
	_ = v140
	var v144 int32
	_ = v144
	var v145 int32
	_ = v145
	var v148 int32
	_ = v148
	var v151 int32
	_ = v151
	var v158 int32
	_ = v158
	var v160 int32
	_ = v160
	var v163 int32
	_ = v163
	var v168 int32
	_ = v168
	var v174 int32
	_ = v174
	var v175 int32
	_ = v175
	var v178 int32
	_ = v178
	var v181 int32
	_ = v181
	var v185 int32
	_ = v185
	var v189 int32
	_ = v189
	var v192 int32
	_ = v192
	var v193 int32
	_ = v193
	var v199 int32
	_ = v199
	var v210 int32
	_ = v210
	var v268 int32
	_ = v268
	var v288 int32
	_ = v288
	var v289 int32
	_ = v289
	var v290 int32
	_ = v290
	var v291 int32
	_ = v291
	var v292 int32
	_ = v292
	var v293 int32
	_ = v293
	var v294 int32
	_ = v294
	var v295 int32
	_ = v295
	var v296 int32
	_ = v296
	var v300 int32
	_ = v300
	var v324 int32
	_ = v324
	var v325 int32
	_ = v325
	var v326 int32
	_ = v326
	var v329 int32
	_ = v329
	var v332 int32
	_ = v332
	var v333 int32
	_ = v333
	var v334 int32
	_ = v334
	var v335 int32
	_ = v335
	var v336 int32
	_ = v336
	var v338 int32
	_ = v338
	var v342 int32
	_ = v342
	var v344 int32
	_ = v344
	var v347 int32
	_ = v347
	var v351 int32
	_ = v351
	var v355 int32
	_ = v355
	var v358 int32
	_ = v358
	var v360 int32
	_ = v360
	var v363 int32
	_ = v363
	var v367 int32
	_ = v367
	var v371 int32
	_ = v371
	var v375 int32
	_ = v375
	var v379 int32
	_ = v379
	var v380 int32
	_ = v380
	var v383 int32
	_ = v383
	var v386 int32
	_ = v386
	var v393 int32
	_ = v393
	var v395 int32
	_ = v395
	var v398 int32
	_ = v398
	var v403 int32
	_ = v403
	var v409 int32
	_ = v409
	var v410 int32
	_ = v410
	var v413 int32
	_ = v413
	var v416 int32
	_ = v416
	var v420 int32
	_ = v420
	var v424 int32
	_ = v424
	var v427 int32
	_ = v427
	var v428 int32
	_ = v428
	var v434 int32
	_ = v434
	var v445 int32
	_ = v445
	var v503 int32
	_ = v503
	var v523 int32
	_ = v523
	var v524 int32
	_ = v524
	var v525 int32
	_ = v525
	var v526 int32
	_ = v526
	var v527 int32
	_ = v527
	var v528 int32
	_ = v528
	var v529 int32
	_ = v529
	var v530 int32
	_ = v530
	var v531 int32
	_ = v531
	var v535 int32
	_ = v535
	var v559 int32
	_ = v559
	var v560 int32
	_ = v560
	var v561 int32
	_ = v561
	var v564 int32
	_ = v564
	var v567 int32
	_ = v567
	var v568 int32
	_ = v568
	var v569 int32
	_ = v569
	var v570 int32
	_ = v570
	var v571 int32
	_ = v571
	var v573 int32
	_ = v573
	var v577 int32
	_ = v577
	var v579 int32
	_ = v579
	var v582 int32
	_ = v582
	var v586 int32
	_ = v586
	var v590 int32
	_ = v590
	var v593 int32
	_ = v593
	var v595 int32
	_ = v595
	var v598 int32
	_ = v598
	var v602 int32
	_ = v602
	var v606 int32
	_ = v606
	var v610 int32
	_ = v610
	var v614 int32
	_ = v614
	var v615 int32
	_ = v615
	var v618 int32
	_ = v618
	var v621 int32
	_ = v621
	var v628 int32
	_ = v628
	var v630 int32
	_ = v630
	var v633 int32
	_ = v633
	var v638 int32
	_ = v638
	var v644 int32
	_ = v644
	var v645 int32
	_ = v645
	var v648 int32
	_ = v648
	var v651 int32
	_ = v651
	var v655 int32
	_ = v655
	var v659 int32
	_ = v659
	var v662 int32
	_ = v662
	var v663 int32
	_ = v663
	var v669 int32
	_ = v669
	var v680 int32
	_ = v680
	v33 = int32(1)
	v53 = m.G55
	v54 = *(*int32)(unsafe.Add(mBase, uint32(v53)))
	v55 = m.G53
	v56 = *(*int32)(unsafe.Add(mBase, uint32(v55)))
	v57 = m.G54
	v58 = *(*int32)(unsafe.Add(mBase, uint32(v57)))
	v59 = m.G56
	v60 = *(*int32)(unsafe.Add(mBase, uint32(v59)))
	v61 = l0 + int32(4)
	v65 = int32(17)
	for {
		v89 = v61 + int32(-1)
		v90 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v89))))
		v91 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v61))))
		v94 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v60+(v90-v91)))))
		v97 = v61 + int32(-2)
		v98 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v97))))
		v99 = v61 + int32(1)
		v100 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v99))))
		v101 = v98 - v100
		v103 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v60+v101))))
		if l2<<(uint(v33)%32)|v33 < v94<<(uint(int32(2))%32)+v103 {
		} else {
			v107 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v61+int32(-4)))))
			v109 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v61+int32(-3)))))
			v112 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v60+(v107-v109)))))
			if l3 < v112 {
			} else {
				v116 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v60+(v109-v98)))))
				if base.Ui32(l3) < base.Ui32(v116) {
				} else {
					v120 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v60+(v98-v90)))))
					if base.Ui32(l3) < base.Ui32(v120) {
					} else {
						v123 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v61+int32(3)))))
						v125 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v61+int32(2)))))
						v128 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v60+(v123-v125)))))
						if base.Ui32(l3) < base.Ui32(v128) {
						} else {
							v132 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v60+(v125-v100)))))
							if base.Ui32(l3) < base.Ui32(v132) {
							} else {
								v136 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v60+(v100-v91)))))
								if base.Ui32(l3) < base.Ui32(v136) {
								} else {
									v140 = (v91 - v90) * int32(3)
									if l4 < v120 {
										v144 = int32(*(*int8)(unsafe.Add(mBase, uint32(v54+v101))))
										v145 = v140 + v144
										v148 = int32(3)
										v151 = int32(*(*int8)(unsafe.Add(mBase, uint32(v58+(v145+int32(4))>>(uint(v148)%32)))))
										v158 = int32(*(*int8)(unsafe.Add(mBase, uint32(v58+(v145+v148)>>(uint(v148)%32)))))
										v160 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v56+v90+v158))))
										*(*uint8)(unsafe.Add(mBase, uint32(v89))) = uint8(v160)
										v192 = v91 - v151
										v193 = v61
									} else {
										if base.Ui32(v136) <= base.Ui32(l4) {
											v163 = int32(3)
											v168 = int32(*(*int8)(unsafe.Add(mBase, uint32(v58+(v140+v163)>>(uint(v163)%32)))))
											v174 = int32(*(*int8)(unsafe.Add(mBase, uint32(v58+(v140+int32(4))>>(uint(v163)%32)))))
											v175 = int32(1)
											v178 = (v174 + v175) >> (uint(v175) % 32)
											v181 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v56+v178+v98))))
											*(*uint8)(unsafe.Add(mBase, uint32(v97))) = uint8(v181)
											v185 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v168+(v56+v90)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v89))) = uint8(v185)
											v189 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v56+(v91-v174)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v61))) = uint8(v189)
											v192 = v100 - v178
											v193 = v99
										} else {
											v144 = int32(*(*int8)(unsafe.Add(mBase, uint32(v54+v101))))
											v145 = v140 + v144
											v148 = int32(3)
											v151 = int32(*(*int8)(unsafe.Add(mBase, uint32(v58+(v145+int32(4))>>(uint(v148)%32)))))
											v158 = int32(*(*int8)(unsafe.Add(mBase, uint32(v58+(v145+v148)>>(uint(v148)%32)))))
											v160 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v56+v90+v158))))
											*(*uint8)(unsafe.Add(mBase, uint32(v89))) = uint8(v160)
											v192 = v91 - v151
											v193 = v61
										}
									}
									v199 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v56+v192))))
									*(*uint8)(unsafe.Add(mBase, uint32(v193))) = uint8(v199)
								}
							}
						}
					}
				}
			}
		}
		v210 = v65 + int32(-1)
		if base.Ui32(int32(1)) < base.Ui32(v210) {
			v61 = v61 + l1
			v65 = v210
			continue
		} else {
			break
		}
		break
	}
	v268 = int32(1)
	v288 = m.G55
	v289 = *(*int32)(unsafe.Add(mBase, uint32(v288)))
	v290 = m.G53
	v291 = *(*int32)(unsafe.Add(mBase, uint32(v290)))
	v292 = m.G54
	v293 = *(*int32)(unsafe.Add(mBase, uint32(v292)))
	v294 = m.G56
	v295 = *(*int32)(unsafe.Add(mBase, uint32(v294)))
	v296 = l0 + int32(8)
	v300 = int32(17)
	for {
		v324 = v296 + int32(-1)
		v325 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324))))
		v326 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296))))
		v329 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v295+(v325-v326)))))
		v332 = v296 + int32(-2)
		v333 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v332))))
		v334 = v296 + int32(1)
		v335 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v334))))
		v336 = v333 - v335
		v338 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v295+v336))))
		if l2<<(uint(v268)%32)|v268 < v329<<(uint(int32(2))%32)+v338 {
		} else {
			v342 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+int32(-4)))))
			v344 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+int32(-3)))))
			v347 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v295+(v342-v344)))))
			if l3 < v347 {
			} else {
				v351 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v295+(v344-v333)))))
				if base.Ui32(l3) < base.Ui32(v351) {
				} else {
					v355 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v295+(v333-v325)))))
					if base.Ui32(l3) < base.Ui32(v355) {
					} else {
						v358 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+int32(3)))))
						v360 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+int32(2)))))
						v363 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v295+(v358-v360)))))
						if base.Ui32(l3) < base.Ui32(v363) {
						} else {
							v367 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v295+(v360-v335)))))
							if base.Ui32(l3) < base.Ui32(v367) {
							} else {
								v371 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v295+(v335-v326)))))
								if base.Ui32(l3) < base.Ui32(v371) {
								} else {
									v375 = (v326 - v325) * int32(3)
									if l4 < v355 {
										v379 = int32(*(*int8)(unsafe.Add(mBase, uint32(v289+v336))))
										v380 = v375 + v379
										v383 = int32(3)
										v386 = int32(*(*int8)(unsafe.Add(mBase, uint32(v293+(v380+int32(4))>>(uint(v383)%32)))))
										v393 = int32(*(*int8)(unsafe.Add(mBase, uint32(v293+(v380+v383)>>(uint(v383)%32)))))
										v395 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v291+v325+v393))))
										*(*uint8)(unsafe.Add(mBase, uint32(v324))) = uint8(v395)
										v427 = v326 - v386
										v428 = v296
									} else {
										if base.Ui32(v371) <= base.Ui32(l4) {
											v398 = int32(3)
											v403 = int32(*(*int8)(unsafe.Add(mBase, uint32(v293+(v375+v398)>>(uint(v398)%32)))))
											v409 = int32(*(*int8)(unsafe.Add(mBase, uint32(v293+(v375+int32(4))>>(uint(v398)%32)))))
											v410 = int32(1)
											v413 = (v409 + v410) >> (uint(v410) % 32)
											v416 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v291+v413+v333))))
											*(*uint8)(unsafe.Add(mBase, uint32(v332))) = uint8(v416)
											v420 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v403+(v291+v325)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v324))) = uint8(v420)
											v424 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v291+(v326-v409)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v296))) = uint8(v424)
											v427 = v335 - v413
											v428 = v334
										} else {
											v379 = int32(*(*int8)(unsafe.Add(mBase, uint32(v289+v336))))
											v380 = v375 + v379
											v383 = int32(3)
											v386 = int32(*(*int8)(unsafe.Add(mBase, uint32(v293+(v380+int32(4))>>(uint(v383)%32)))))
											v393 = int32(*(*int8)(unsafe.Add(mBase, uint32(v293+(v380+v383)>>(uint(v383)%32)))))
											v395 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v291+v325+v393))))
											*(*uint8)(unsafe.Add(mBase, uint32(v324))) = uint8(v395)
											v427 = v326 - v386
											v428 = v296
										}
									}
									v434 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v291+v427))))
									*(*uint8)(unsafe.Add(mBase, uint32(v428))) = uint8(v434)
								}
							}
						}
					}
				}
			}
		}
		v445 = v300 + int32(-1)
		if base.Ui32(int32(1)) < base.Ui32(v445) {
			v296 = v296 + l1
			v300 = v445
			continue
		} else {
			break
		}
		break
	}
	v503 = int32(1)
	v523 = m.G55
	v524 = *(*int32)(unsafe.Add(mBase, uint32(v523)))
	v525 = m.G53
	v526 = *(*int32)(unsafe.Add(mBase, uint32(v525)))
	v527 = m.G54
	v528 = *(*int32)(unsafe.Add(mBase, uint32(v527)))
	v529 = m.G56
	v530 = *(*int32)(unsafe.Add(mBase, uint32(v529)))
	v531 = l0 + int32(12)
	v535 = int32(17)
	for {
		v559 = v531 + int32(-1)
		v560 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v559))))
		v561 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v531))))
		v564 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v530+(v560-v561)))))
		v567 = v531 + int32(-2)
		v568 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v567))))
		v569 = v531 + int32(1)
		v570 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v569))))
		v571 = v568 - v570
		v573 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v530+v571))))
		if l2<<(uint(v503)%32)|v503 < v564<<(uint(int32(2))%32)+v573 {
		} else {
			v577 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v531+int32(-4)))))
			v579 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v531+int32(-3)))))
			v582 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v530+(v577-v579)))))
			if l3 < v582 {
			} else {
				v586 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v530+(v579-v568)))))
				if base.Ui32(l3) < base.Ui32(v586) {
				} else {
					v590 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v530+(v568-v560)))))
					if base.Ui32(l3) < base.Ui32(v590) {
					} else {
						v593 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v531+int32(3)))))
						v595 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v531+int32(2)))))
						v598 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v530+(v593-v595)))))
						if base.Ui32(l3) < base.Ui32(v598) {
						} else {
							v602 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v530+(v595-v570)))))
							if base.Ui32(l3) < base.Ui32(v602) {
							} else {
								v606 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v530+(v570-v561)))))
								if base.Ui32(l3) < base.Ui32(v606) {
								} else {
									v610 = (v561 - v560) * int32(3)
									if l4 < v590 {
										v614 = int32(*(*int8)(unsafe.Add(mBase, uint32(v524+v571))))
										v615 = v610 + v614
										v618 = int32(3)
										v621 = int32(*(*int8)(unsafe.Add(mBase, uint32(v528+(v615+int32(4))>>(uint(v618)%32)))))
										v628 = int32(*(*int8)(unsafe.Add(mBase, uint32(v528+(v615+v618)>>(uint(v618)%32)))))
										v630 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v526+v560+v628))))
										*(*uint8)(unsafe.Add(mBase, uint32(v559))) = uint8(v630)
										v662 = v561 - v621
										v663 = v531
									} else {
										if base.Ui32(v606) <= base.Ui32(l4) {
											v633 = int32(3)
											v638 = int32(*(*int8)(unsafe.Add(mBase, uint32(v528+(v610+v633)>>(uint(v633)%32)))))
											v644 = int32(*(*int8)(unsafe.Add(mBase, uint32(v528+(v610+int32(4))>>(uint(v633)%32)))))
											v645 = int32(1)
											v648 = (v644 + v645) >> (uint(v645) % 32)
											v651 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v526+v648+v568))))
											*(*uint8)(unsafe.Add(mBase, uint32(v567))) = uint8(v651)
											v655 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v638+(v526+v560)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v559))) = uint8(v655)
											v659 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v526+(v561-v644)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v531))) = uint8(v659)
											v662 = v570 - v648
											v663 = v569
										} else {
											v614 = int32(*(*int8)(unsafe.Add(mBase, uint32(v524+v571))))
											v615 = v610 + v614
											v618 = int32(3)
											v621 = int32(*(*int8)(unsafe.Add(mBase, uint32(v528+(v615+int32(4))>>(uint(v618)%32)))))
											v628 = int32(*(*int8)(unsafe.Add(mBase, uint32(v528+(v615+v618)>>(uint(v618)%32)))))
											v630 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v526+v560+v628))))
											*(*uint8)(unsafe.Add(mBase, uint32(v559))) = uint8(v630)
											v662 = v561 - v621
											v663 = v531
										}
									}
									v669 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v526+v662))))
									*(*uint8)(unsafe.Add(mBase, uint32(v663))) = uint8(v669)
								}
							}
						}
					}
				}
			}
		}
		v680 = v535 + int32(-1)
		if base.Ui32(int32(1)) < base.Ui32(v680) {
			v531 = v531 + l1
			v535 = v680
			continue
		} else {
			break
		}
		break
	}
	return
}
func F_HFilter8_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v36 int32
	_ = v36
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
	var v64 int32
	_ = v64
	var v68 int32
	_ = v68
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v101 int32
	_ = v101
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v110 int32
	_ = v110
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v116 int32
	_ = v116
	var v119 int32
	_ = v119
	var v123 int32
	_ = v123
	var v127 int32
	_ = v127
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v132 int32
	_ = v132
	var v135 int32
	_ = v135
	var v139 int32
	_ = v139
	var v143 int32
	_ = v143
	var v147 int32
	_ = v147
	var v149 int32
	_ = v149
	var v152 int32
	_ = v152
	var v155 int32
	_ = v155
	var v158 int32
	_ = v158
	var v165 int32
	_ = v165
	var v167 int32
	_ = v167
	var v172 int32
	_ = v172
	var v175 int32
	_ = v175
	var v177 int32
	_ = v177
	var v178 int32
	_ = v178
	var v181 int32
	_ = v181
	var v188 int32
	_ = v188
	var v191 int32
	_ = v191
	var v198 int32
	_ = v198
	var v201 int32
	_ = v201
	var v205 int32
	_ = v205
	var v209 int32
	_ = v209
	var v212 int32
	_ = v212
	var v217 int32
	_ = v217
	var v219 int32
	_ = v219
	var v233 int32
	_ = v233
	var v297 int32
	_ = v297
	var v317 int32
	_ = v317
	var v318 int32
	_ = v318
	var v319 int32
	_ = v319
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
	var v329 int32
	_ = v329
	var v357 int32
	_ = v357
	var v358 int32
	_ = v358
	var v359 int32
	_ = v359
	var v362 int32
	_ = v362
	var v365 int32
	_ = v365
	var v366 int32
	_ = v366
	var v367 int32
	_ = v367
	var v368 int32
	_ = v368
	var v369 int32
	_ = v369
	var v371 int32
	_ = v371
	var v375 int32
	_ = v375
	var v376 int32
	_ = v376
	var v377 int32
	_ = v377
	var v380 int32
	_ = v380
	var v384 int32
	_ = v384
	var v388 int32
	_ = v388
	var v391 int32
	_ = v391
	var v392 int32
	_ = v392
	var v393 int32
	_ = v393
	var v396 int32
	_ = v396
	var v400 int32
	_ = v400
	var v404 int32
	_ = v404
	var v408 int32
	_ = v408
	var v410 int32
	_ = v410
	var v413 int32
	_ = v413
	var v416 int32
	_ = v416
	var v419 int32
	_ = v419
	var v426 int32
	_ = v426
	var v428 int32
	_ = v428
	var v433 int32
	_ = v433
	var v436 int32
	_ = v436
	var v438 int32
	_ = v438
	var v439 int32
	_ = v439
	var v442 int32
	_ = v442
	var v449 int32
	_ = v449
	var v452 int32
	_ = v452
	var v459 int32
	_ = v459
	var v462 int32
	_ = v462
	var v466 int32
	_ = v466
	var v470 int32
	_ = v470
	var v473 int32
	_ = v473
	var v478 int32
	_ = v478
	var v480 int32
	_ = v480
	var v494 int32
	_ = v494
	v36 = int32(1)
	v56 = m.G54
	v57 = *(*int32)(unsafe.Add(mBase, uint32(v56)))
	v58 = m.G53
	v59 = *(*int32)(unsafe.Add(mBase, uint32(v58)))
	v60 = m.G55
	v61 = *(*int32)(unsafe.Add(mBase, uint32(v60)))
	v62 = m.G56
	v63 = *(*int32)(unsafe.Add(mBase, uint32(v62)))
	v64 = l0
	v68 = int32(9)
	for {
		v96 = v64 + int32(-1)
		v97 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v96))))
		v98 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v64))))
		v101 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v97-v98)))))
		v104 = v64 + int32(-2)
		v105 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v104))))
		v106 = v64 + int32(1)
		v107 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v106))))
		v108 = v105 - v107
		v110 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+v108))))
		if l3<<(uint(v36)%32)|v36 < v101<<(uint(int32(2))%32)+v110 {
		} else {
			v114 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v64+int32(-4)))))
			v115 = v64 + int32(-3)
			v116 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v115))))
			v119 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v114-v116)))))
			if l4 < v119 {
			} else {
				v123 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v116-v105)))))
				if base.Ui32(l4) < base.Ui32(v123) {
				} else {
					v127 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v105-v97)))))
					if base.Ui32(l4) < base.Ui32(v127) {
					} else {
						v130 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v64+int32(3)))))
						v131 = v64 + int32(2)
						v132 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v131))))
						v135 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v130-v132)))))
						if base.Ui32(l4) < base.Ui32(v135) {
						} else {
							v139 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v132-v107)))))
							if base.Ui32(l4) < base.Ui32(v139) {
							} else {
								v143 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v107-v98)))))
								if base.Ui32(l4) < base.Ui32(v143) {
								} else {
									v147 = (v98 - v97) * int32(3)
									v149 = int32(*(*int8)(unsafe.Add(mBase, uint32(v61+v108))))
									if l5 < v127 {
										v152 = v147 + v149
										v155 = int32(3)
										v158 = int32(*(*int8)(unsafe.Add(mBase, uint32(v57+(v152+int32(4))>>(uint(v155)%32)))))
										v165 = int32(*(*int8)(unsafe.Add(mBase, uint32(v57+(v152+v155)>>(uint(v155)%32)))))
										v167 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v97+v165))))
										*(*uint8)(unsafe.Add(mBase, uint32(v96))) = uint8(v167)
										v212 = v98 - v158
										v217 = v64
									} else {
										if base.Ui32(v143) <= base.Ui32(l5) {
											v172 = int32(*(*int8)(unsafe.Add(mBase, uint32(v61+v147+v149))))
											v175 = int32(63)
											v177 = int32(7)
											v178 = (v172*int32(9) + v175) >> (uint(v177) % 32)
											v181 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v178+v116))))
											*(*uint8)(unsafe.Add(mBase, uint32(v115))) = uint8(v181)
											v188 = (v172*int32(18) + v175) >> (uint(v177) % 32)
											v191 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v188+v105))))
											*(*uint8)(unsafe.Add(mBase, uint32(v104))) = uint8(v191)
											v198 = (v172*int32(27) + v175) >> (uint(v177) % 32)
											v201 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v198+v97))))
											*(*uint8)(unsafe.Add(mBase, uint32(v96))) = uint8(v201)
											v205 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+(v98-v198)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v64))) = uint8(v205)
											v209 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+(v107-v188)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v106))) = uint8(v209)
											v212 = v132 - v178
											v217 = v131
										} else {
											v152 = v147 + v149
											v155 = int32(3)
											v158 = int32(*(*int8)(unsafe.Add(mBase, uint32(v57+(v152+int32(4))>>(uint(v155)%32)))))
											v165 = int32(*(*int8)(unsafe.Add(mBase, uint32(v57+(v152+v155)>>(uint(v155)%32)))))
											v167 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v97+v165))))
											*(*uint8)(unsafe.Add(mBase, uint32(v96))) = uint8(v167)
											v212 = v98 - v158
											v217 = v64
										}
									}
									v219 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v212))))
									*(*uint8)(unsafe.Add(mBase, uint32(v217))) = uint8(v219)
								}
							}
						}
					}
				}
			}
		}
		v233 = v68 + int32(-1)
		if base.Ui32(int32(1)) < base.Ui32(v233) {
			v64 = v64 + l2
			v68 = v233
			continue
		} else {
			break
		}
		break
	}
	v297 = int32(1)
	v317 = m.G54
	v318 = *(*int32)(unsafe.Add(mBase, uint32(v317)))
	v319 = m.G53
	v320 = *(*int32)(unsafe.Add(mBase, uint32(v319)))
	v321 = m.G55
	v322 = *(*int32)(unsafe.Add(mBase, uint32(v321)))
	v323 = m.G56
	v324 = *(*int32)(unsafe.Add(mBase, uint32(v323)))
	v325 = l1
	v329 = int32(9)
	for {
		v357 = v325 + int32(-1)
		v358 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v357))))
		v359 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v325))))
		v362 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324+(v358-v359)))))
		v365 = v325 + int32(-2)
		v366 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v365))))
		v367 = v325 + int32(1)
		v368 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367))))
		v369 = v366 - v368
		v371 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324+v369))))
		if l3<<(uint(v297)%32)|v297 < v362<<(uint(int32(2))%32)+v371 {
		} else {
			v375 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v325+int32(-4)))))
			v376 = v325 + int32(-3)
			v377 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v376))))
			v380 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324+(v375-v377)))))
			if l4 < v380 {
			} else {
				v384 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324+(v377-v366)))))
				if base.Ui32(l4) < base.Ui32(v384) {
				} else {
					v388 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324+(v366-v358)))))
					if base.Ui32(l4) < base.Ui32(v388) {
					} else {
						v391 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v325+int32(3)))))
						v392 = v325 + int32(2)
						v393 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v392))))
						v396 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324+(v391-v393)))))
						if base.Ui32(l4) < base.Ui32(v396) {
						} else {
							v400 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324+(v393-v368)))))
							if base.Ui32(l4) < base.Ui32(v400) {
							} else {
								v404 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324+(v368-v359)))))
								if base.Ui32(l4) < base.Ui32(v404) {
								} else {
									v408 = (v359 - v358) * int32(3)
									v410 = int32(*(*int8)(unsafe.Add(mBase, uint32(v322+v369))))
									if l5 < v388 {
										v413 = v408 + v410
										v416 = int32(3)
										v419 = int32(*(*int8)(unsafe.Add(mBase, uint32(v318+(v413+int32(4))>>(uint(v416)%32)))))
										v426 = int32(*(*int8)(unsafe.Add(mBase, uint32(v318+(v413+v416)>>(uint(v416)%32)))))
										v428 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+v358+v426))))
										*(*uint8)(unsafe.Add(mBase, uint32(v357))) = uint8(v428)
										v473 = v359 - v419
										v478 = v325
									} else {
										if base.Ui32(v404) <= base.Ui32(l5) {
											v433 = int32(*(*int8)(unsafe.Add(mBase, uint32(v322+v408+v410))))
											v436 = int32(63)
											v438 = int32(7)
											v439 = (v433*int32(9) + v436) >> (uint(v438) % 32)
											v442 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+v439+v377))))
											*(*uint8)(unsafe.Add(mBase, uint32(v376))) = uint8(v442)
											v449 = (v433*int32(18) + v436) >> (uint(v438) % 32)
											v452 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+v449+v366))))
											*(*uint8)(unsafe.Add(mBase, uint32(v365))) = uint8(v452)
											v459 = (v433*int32(27) + v436) >> (uint(v438) % 32)
											v462 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+v459+v358))))
											*(*uint8)(unsafe.Add(mBase, uint32(v357))) = uint8(v462)
											v466 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+(v359-v459)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v325))) = uint8(v466)
											v470 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+(v368-v449)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v367))) = uint8(v470)
											v473 = v393 - v439
											v478 = v392
										} else {
											v413 = v408 + v410
											v416 = int32(3)
											v419 = int32(*(*int8)(unsafe.Add(mBase, uint32(v318+(v413+int32(4))>>(uint(v416)%32)))))
											v426 = int32(*(*int8)(unsafe.Add(mBase, uint32(v318+(v413+v416)>>(uint(v416)%32)))))
											v428 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+v358+v426))))
											*(*uint8)(unsafe.Add(mBase, uint32(v357))) = uint8(v428)
											v473 = v359 - v419
											v478 = v325
										}
									}
									v480 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+v473))))
									*(*uint8)(unsafe.Add(mBase, uint32(v478))) = uint8(v480)
								}
							}
						}
					}
				}
			}
		}
		v494 = v329 + int32(-1)
		if base.Ui32(int32(1)) < base.Ui32(v494) {
			v325 = v325 + l2
			v329 = v494
			continue
		} else {
			break
		}
		break
	}
	return
}
func F_HFilter8i_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v34 int32
	_ = v34
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
	var v66 int32
	_ = v66
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v95 int32
	_ = v95
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v104 int32
	_ = v104
	var v108 int32
	_ = v108
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v117 int32
	_ = v117
	var v121 int32
	_ = v121
	var v124 int32
	_ = v124
	var v126 int32
	_ = v126
	var v129 int32
	_ = v129
	var v133 int32
	_ = v133
	var v137 int32
	_ = v137
	var v141 int32
	_ = v141
	var v145 int32
	_ = v145
	var v146 int32
	_ = v146
	var v149 int32
	_ = v149
	var v152 int32
	_ = v152
	var v159 int32
	_ = v159
	var v161 int32
	_ = v161
	var v164 int32
	_ = v164
	var v169 int32
	_ = v169
	var v175 int32
	_ = v175
	var v176 int32
	_ = v176
	var v179 int32
	_ = v179
	var v182 int32
	_ = v182
	var v186 int32
	_ = v186
	var v190 int32
	_ = v190
	var v193 int32
	_ = v193
	var v194 int32
	_ = v194
	var v200 int32
	_ = v200
	var v211 int32
	_ = v211
	var v269 int32
	_ = v269
	var v289 int32
	_ = v289
	var v290 int32
	_ = v290
	var v291 int32
	_ = v291
	var v292 int32
	_ = v292
	var v293 int32
	_ = v293
	var v294 int32
	_ = v294
	var v295 int32
	_ = v295
	var v296 int32
	_ = v296
	var v297 int32
	_ = v297
	var v301 int32
	_ = v301
	var v325 int32
	_ = v325
	var v326 int32
	_ = v326
	var v327 int32
	_ = v327
	var v330 int32
	_ = v330
	var v333 int32
	_ = v333
	var v334 int32
	_ = v334
	var v335 int32
	_ = v335
	var v336 int32
	_ = v336
	var v337 int32
	_ = v337
	var v339 int32
	_ = v339
	var v343 int32
	_ = v343
	var v345 int32
	_ = v345
	var v348 int32
	_ = v348
	var v352 int32
	_ = v352
	var v356 int32
	_ = v356
	var v359 int32
	_ = v359
	var v361 int32
	_ = v361
	var v364 int32
	_ = v364
	var v368 int32
	_ = v368
	var v372 int32
	_ = v372
	var v376 int32
	_ = v376
	var v380 int32
	_ = v380
	var v381 int32
	_ = v381
	var v384 int32
	_ = v384
	var v387 int32
	_ = v387
	var v394 int32
	_ = v394
	var v396 int32
	_ = v396
	var v399 int32
	_ = v399
	var v404 int32
	_ = v404
	var v410 int32
	_ = v410
	var v411 int32
	_ = v411
	var v414 int32
	_ = v414
	var v417 int32
	_ = v417
	var v421 int32
	_ = v421
	var v425 int32
	_ = v425
	var v428 int32
	_ = v428
	var v429 int32
	_ = v429
	var v435 int32
	_ = v435
	var v446 int32
	_ = v446
	v34 = int32(1)
	v54 = m.G55
	v55 = *(*int32)(unsafe.Add(mBase, uint32(v54)))
	v56 = m.G53
	v57 = *(*int32)(unsafe.Add(mBase, uint32(v56)))
	v58 = m.G54
	v59 = *(*int32)(unsafe.Add(mBase, uint32(v58)))
	v60 = m.G56
	v61 = *(*int32)(unsafe.Add(mBase, uint32(v60)))
	v62 = l0 + int32(4)
	v66 = int32(9)
	for {
		v90 = v62 + int32(-1)
		v91 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v90))))
		v92 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62))))
		v95 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v61+(v91-v92)))))
		v98 = v62 + int32(-2)
		v99 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v98))))
		v100 = v62 + int32(1)
		v101 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v100))))
		v102 = v99 - v101
		v104 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v61+v102))))
		if l3<<(uint(v34)%32)|v34 < v95<<(uint(int32(2))%32)+v104 {
		} else {
			v108 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+int32(-4)))))
			v110 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+int32(-3)))))
			v113 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v61+(v108-v110)))))
			if l4 < v113 {
			} else {
				v117 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v61+(v110-v99)))))
				if base.Ui32(l4) < base.Ui32(v117) {
				} else {
					v121 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v61+(v99-v91)))))
					if base.Ui32(l4) < base.Ui32(v121) {
					} else {
						v124 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+int32(3)))))
						v126 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+int32(2)))))
						v129 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v61+(v124-v126)))))
						if base.Ui32(l4) < base.Ui32(v129) {
						} else {
							v133 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v61+(v126-v101)))))
							if base.Ui32(l4) < base.Ui32(v133) {
							} else {
								v137 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v61+(v101-v92)))))
								if base.Ui32(l4) < base.Ui32(v137) {
								} else {
									v141 = (v92 - v91) * int32(3)
									if l5 < v121 {
										v145 = int32(*(*int8)(unsafe.Add(mBase, uint32(v55+v102))))
										v146 = v141 + v145
										v149 = int32(3)
										v152 = int32(*(*int8)(unsafe.Add(mBase, uint32(v59+(v146+int32(4))>>(uint(v149)%32)))))
										v159 = int32(*(*int8)(unsafe.Add(mBase, uint32(v59+(v146+v149)>>(uint(v149)%32)))))
										v161 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57+v91+v159))))
										*(*uint8)(unsafe.Add(mBase, uint32(v90))) = uint8(v161)
										v193 = v92 - v152
										v194 = v62
									} else {
										if base.Ui32(v137) <= base.Ui32(l5) {
											v164 = int32(3)
											v169 = int32(*(*int8)(unsafe.Add(mBase, uint32(v59+(v141+v164)>>(uint(v164)%32)))))
											v175 = int32(*(*int8)(unsafe.Add(mBase, uint32(v59+(v141+int32(4))>>(uint(v164)%32)))))
											v176 = int32(1)
											v179 = (v175 + v176) >> (uint(v176) % 32)
											v182 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57+v179+v99))))
											*(*uint8)(unsafe.Add(mBase, uint32(v98))) = uint8(v182)
											v186 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v169+(v57+v91)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v90))) = uint8(v186)
											v190 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57+(v92-v175)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v62))) = uint8(v190)
											v193 = v101 - v179
											v194 = v100
										} else {
											v145 = int32(*(*int8)(unsafe.Add(mBase, uint32(v55+v102))))
											v146 = v141 + v145
											v149 = int32(3)
											v152 = int32(*(*int8)(unsafe.Add(mBase, uint32(v59+(v146+int32(4))>>(uint(v149)%32)))))
											v159 = int32(*(*int8)(unsafe.Add(mBase, uint32(v59+(v146+v149)>>(uint(v149)%32)))))
											v161 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57+v91+v159))))
											*(*uint8)(unsafe.Add(mBase, uint32(v90))) = uint8(v161)
											v193 = v92 - v152
											v194 = v62
										}
									}
									v200 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57+v193))))
									*(*uint8)(unsafe.Add(mBase, uint32(v194))) = uint8(v200)
								}
							}
						}
					}
				}
			}
		}
		v211 = v66 + int32(-1)
		if base.Ui32(int32(1)) < base.Ui32(v211) {
			v62 = v62 + l2
			v66 = v211
			continue
		} else {
			break
		}
		break
	}
	v269 = int32(1)
	v289 = m.G55
	v290 = *(*int32)(unsafe.Add(mBase, uint32(v289)))
	v291 = m.G53
	v292 = *(*int32)(unsafe.Add(mBase, uint32(v291)))
	v293 = m.G54
	v294 = *(*int32)(unsafe.Add(mBase, uint32(v293)))
	v295 = m.G56
	v296 = *(*int32)(unsafe.Add(mBase, uint32(v295)))
	v297 = l1 + int32(4)
	v301 = int32(9)
	for {
		v325 = v297 + int32(-1)
		v326 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v325))))
		v327 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297))))
		v330 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+(v326-v327)))))
		v333 = v297 + int32(-2)
		v334 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v333))))
		v335 = v297 + int32(1)
		v336 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v335))))
		v337 = v334 - v336
		v339 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+v337))))
		if l3<<(uint(v269)%32)|v269 < v330<<(uint(int32(2))%32)+v339 {
		} else {
			v343 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297+int32(-4)))))
			v345 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297+int32(-3)))))
			v348 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+(v343-v345)))))
			if l4 < v348 {
			} else {
				v352 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+(v345-v334)))))
				if base.Ui32(l4) < base.Ui32(v352) {
				} else {
					v356 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+(v334-v326)))))
					if base.Ui32(l4) < base.Ui32(v356) {
					} else {
						v359 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297+int32(3)))))
						v361 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297+int32(2)))))
						v364 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+(v359-v361)))))
						if base.Ui32(l4) < base.Ui32(v364) {
						} else {
							v368 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+(v361-v336)))))
							if base.Ui32(l4) < base.Ui32(v368) {
							} else {
								v372 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+(v336-v327)))))
								if base.Ui32(l4) < base.Ui32(v372) {
								} else {
									v376 = (v327 - v326) * int32(3)
									if l5 < v356 {
										v380 = int32(*(*int8)(unsafe.Add(mBase, uint32(v290+v337))))
										v381 = v376 + v380
										v384 = int32(3)
										v387 = int32(*(*int8)(unsafe.Add(mBase, uint32(v294+(v381+int32(4))>>(uint(v384)%32)))))
										v394 = int32(*(*int8)(unsafe.Add(mBase, uint32(v294+(v381+v384)>>(uint(v384)%32)))))
										v396 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v292+v326+v394))))
										*(*uint8)(unsafe.Add(mBase, uint32(v325))) = uint8(v396)
										v428 = v327 - v387
										v429 = v297
									} else {
										if base.Ui32(v372) <= base.Ui32(l5) {
											v399 = int32(3)
											v404 = int32(*(*int8)(unsafe.Add(mBase, uint32(v294+(v376+v399)>>(uint(v399)%32)))))
											v410 = int32(*(*int8)(unsafe.Add(mBase, uint32(v294+(v376+int32(4))>>(uint(v399)%32)))))
											v411 = int32(1)
											v414 = (v410 + v411) >> (uint(v411) % 32)
											v417 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v292+v414+v334))))
											*(*uint8)(unsafe.Add(mBase, uint32(v333))) = uint8(v417)
											v421 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v404+(v292+v326)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v325))) = uint8(v421)
											v425 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v292+(v327-v410)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v297))) = uint8(v425)
											v428 = v336 - v414
											v429 = v335
										} else {
											v380 = int32(*(*int8)(unsafe.Add(mBase, uint32(v290+v337))))
											v381 = v376 + v380
											v384 = int32(3)
											v387 = int32(*(*int8)(unsafe.Add(mBase, uint32(v294+(v381+int32(4))>>(uint(v384)%32)))))
											v394 = int32(*(*int8)(unsafe.Add(mBase, uint32(v294+(v381+v384)>>(uint(v384)%32)))))
											v396 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v292+v326+v394))))
											*(*uint8)(unsafe.Add(mBase, uint32(v325))) = uint8(v396)
											v428 = v327 - v387
											v429 = v297
										}
									}
									v435 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v292+v428))))
									*(*uint8)(unsafe.Add(mBase, uint32(v429))) = uint8(v435)
								}
							}
						}
					}
				}
			}
		}
		v446 = v301 + int32(-1)
		if base.Ui32(int32(1)) < base.Ui32(v446) {
			v297 = v297 + l2
			v301 = v446
			continue
		} else {
			break
		}
		break
	}
	return
}
func F_HU4_C(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v7 int32
	_ = v7
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v18 int32
	_ = v18
	var v24 int32
	_ = v24
	var v30 int32
	_ = v30
	var v35 int32
	_ = v35
	var v41 int32
	_ = v41
	var v44 int32
	_ = v44
	var v47 int32
	_ = v47
	var v53 int32
	_ = v53
	var v60 int32
	_ = v60
	v7 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+95)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+99)) = uint8(v7)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+98)) = uint8(v7)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+97)) = uint8(v7)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+96)) = uint8(v7)
	v12 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+31)))
	v13 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+63)))
	v15 = int32(1)
	v18 = int32(base.Ui32(v12+v13+v15) >> (uint(v15) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+32)) = uint8(v18)
	v24 = int32(base.Ui32(v13+v7+v15) >> (uint(v15) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+64)) = uint8(v24)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v18)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+34)) = uint8(v24)
	v30 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-1)))))
	v35 = int32(base.Ui32(v12+v30+v15) >> (uint(v15) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v35)
	v41 = int32(2)
	v44 = int32(base.Ui32(v7+v12+v13<<(uint(v15)%32)+v41) >> (uint(v41) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+33)) = uint8(v44)
	v47 = v13 + v41
	v53 = int32(base.Ui32(v30+v47+v12<<(uint(v15)%32)) >> (uint(v41) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v53)
	v60 = int32(base.Ui32(v47+v7+v7<<(uint(v15)%32)) >> (uint(v41) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+65)) = uint8(v60)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)) = uint8(v44)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+35)) = uint8(v60)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+67)) = uint8(v7)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+66)) = uint8(v7)
	return
}
func F_HasAlpha32b_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v18 int32
	_ = v18
	var v24 int32
	_ = v24
	v3 = int32(0)
	if l1 < int32(1) {
		v24 = v3
		goto L1
	} else {
		goto L2
	}
L1:
	;
	return v24
L2:
	;
	v9 = l0
	v10 = l1 + int32(1)
	goto L4
L3:
	;
	v24 = int32(1)
	goto L1
L4:
	;
	v12 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v9))))
	if v12 != int32(255) {
		goto L3
	} else {
		goto L6
	}
L6:
	;
	v18 = v10 + int32(-1)
	if v18 < int32(2) {
		v24 = v3
		goto L1
	} else {
		goto L7
	}
L7:
	;
	v9 = v9 + int32(4)
	v10 = v18
	goto L4
}
func F_HasAlpha8b_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v5 int32
	_ = v5
	var v13 int32
	_ = v13
	v4 = l0
	v5 = l1
	goto L1
L1:
	;
	if int32(1) <= v5 {
		goto L3
	} else {
		goto L4
	}
L2:
	;
	return int32(1)
L3:
	;
	v13 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4))))
	if v13 == int32(255) {
		v4 = v4 + int32(1)
		v5 = v5 + int32(-1)
		goto L1
	} else {
		goto L5
	}
L4:
	;
	return int32(0)
L5:
	;
	goto L2
}
func F_HistogramAdd(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
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
	var v35 int32
	_ = v35
	var v39 int32
	_ = v39
	var v41 int32
	_ = v41
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v51 int32
	_ = v51
	var v52 int32
	_ = v52
	var v56 int32
	_ = v56
	var v59 int32
	_ = v59
	var v62 int32
	_ = v62
	var v66 int32
	_ = v66
	var v67 int32
	_ = v67
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v74 int32
	_ = v74
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v86 int32
	_ = v86
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v94 int32
	_ = v94
	var v95 int32
	_ = v95
	var v98 int32
	_ = v98
	var v102 int32
	_ = v102
	var v103 int32
	_ = v103
	var v106 int32
	_ = v106
	var v108 int32
	_ = v108
	var v110 int32
	_ = v110
	var v115 int32
	_ = v115
	var v116 int32
	_ = v116
	var v117 int32
	_ = v117
	var v118 int32
	_ = v118
	var v119 int32
	_ = v119
	var v126 int32
	_ = v126
	var v127 int32
	_ = v127
	var v131 int32
	_ = v131
	var v133 int32
	_ = v133
	var v138 int32
	_ = v138
	var v139 int32
	_ = v139
	var v150 int32
	_ = v150
	var v172 int32
	_ = v172
	var v175 int32
	_ = v175
	var v176 int32
	_ = v176
	var v184 int32
	_ = v184
	var v185 int32
	_ = v185
	var v220 int32
	_ = v220
	var v221 int32
	_ = v221
	var v226 int64
	_ = v226
	var v229 int32
	_ = v229
	var v230 int32
	_ = v230
	var v242 int32
	_ = v242
	var v262 int32
	_ = v262
	var v263 int32
	_ = v263
	var v264 int32
	_ = v264
	var v272 int32
	_ = v272
	var v273 int32
	_ = v273
	var v276 int32
	_ = v276
	var v282 int32
	_ = v282
	var v406 int32
	_ = v406
	var v407 int32
	_ = v407
	var v408 int32
	_ = v408
	var v416 int32
	_ = v416
	var v417 int32
	_ = v417
	var v420 int32
	_ = v420
	var v426 int32
	_ = v426
	var v550 int32
	_ = v550
	var v551 int32
	_ = v551
	var v552 int32
	_ = v552
	var v560 int32
	_ = v560
	var v561 int32
	_ = v561
	var v564 int32
	_ = v564
	var v570 int32
	_ = v570
	var v694 int32
	_ = v694
	var v695 int32
	_ = v695
	var v696 int32
	_ = v696
	var v704 int32
	_ = v704
	var v705 int32
	_ = v705
	var v708 int32
	_ = v708
	var v714 int32
	_ = v714
	var v845 int32
	_ = v845
	var v847 int32
	_ = v847
	var v849 int32
	_ = v849
	var v851 int32
	_ = v851
	var v853 int32
	_ = v853
	var v854 int32
	_ = v854
	var v857 int32
	_ = v857
	var v859 int32
	_ = v859
	var v861 int32
	_ = v861
	var v863 int32
	_ = v863
	var v865 int32
	_ = v865
	var v866 int32
	_ = v866
	var v869 int32
	_ = v869
	var v871 int32
	_ = v871
	var v873 int32
	_ = v873
	var v875 int32
	_ = v875
	var v877 int32
	_ = v877
	var v879 int32
	_ = v879
	var v880 int32
	_ = v880
	var v883 int32
	_ = v883
	var v885 int32
	_ = v885
	var v887 int32
	_ = v887
	var v889 int32
	_ = v889
	var v891 int32
	_ = v891
	var v892 int32
	_ = v892
	var v895 int32
	_ = v895
	var v897 int32
	_ = v897
	var v899 int32
	_ = v899
	var v901 int32
	_ = v901
	var v904 int32
	_ = v904
	var v905 int32
	_ = v905
	var v908 int32
	_ = v908
	v17 = int32(4)
	v18 = l2 + v17
	v20 = l0 + v17
	v21 = int32(1028)
	v22 = l2 + v21
	v24 = l0 + v21
	v25 = int32(2052)
	v26 = l2 + v25
	v28 = l0 + v25
	v29 = int32(3076)
	v30 = l2 + v29
	v32 = l0 + v29
	v33 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if l1 != l2 {
		v108 = *(*int32)(unsafe.Add(mBase, uint32(l2)+3236))
		v110 = int32(280)
		if int32(0) < v108 {
			v115 = int32(1)<<(uint(v108)%32) + v110
		} else {
			v115 = v110
		}
		v116 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3304)))
		v117 = *(*int32)(unsafe.Add(mBase, uint32(l2)))
		v118 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
		v119 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3304)))
		if v119 == int32(0) {
			v133 = v115 << (uint(int32(2)) % 32)
			if v116&int32(255) == int32(0) {
				v139 = int32(0)
				if base.Ui32(v133) < base.Ui32(int32(33)) {
					if v133 == int32(0) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v117))) = uint8(v139)
						v150 = v117 + v133
						*(*uint8)(unsafe.Add(mBase, uint32(v150+int32(-1)))) = uint8(v139)
						if base.Ui32(v133) < base.Ui32(int32(3)) {
						} else {
							*(*uint8)(unsafe.Add(mBase, uint32(v117)+2)) = uint8(v139)
							*(*uint8)(unsafe.Add(mBase, uint32(v117)+1)) = uint8(v139)
							*(*uint8)(unsafe.Add(mBase, uint32(v150+int32(-3)))) = uint8(v139)
							*(*uint8)(unsafe.Add(mBase, uint32(v150+int32(-2)))) = uint8(v139)
							if base.Ui32(v133) < base.Ui32(int32(7)) {
							} else {
								*(*uint8)(unsafe.Add(mBase, uint32(v117)+3)) = uint8(v139)
								*(*uint8)(unsafe.Add(mBase, uint32(v150+int32(-4)))) = uint8(v139)
								if base.Ui32(v133) < base.Ui32(int32(9)) {
								} else {
									v172 = int32(0)
									v175 = (v172 - v117) & int32(3)
									v176 = v117 + v175
									*(*int32)(unsafe.Add(mBase, uint32(v176))) = v172
									v184 = (v133 - v175) & int32(60)
									v185 = v176 + v184
									*(*int32)(unsafe.Add(mBase, uint32(v185+int32(-4)))) = v172
									if base.Ui32(v184) < base.Ui32(int32(9)) {
									} else {
										*(*int32)(unsafe.Add(mBase, uint32(v176)+8)) = v172
										*(*int32)(unsafe.Add(mBase, uint32(v176)+4)) = v172
										*(*int32)(unsafe.Add(mBase, uint32(v185+int32(-8)))) = v172
										*(*int32)(unsafe.Add(mBase, uint32(v185+int32(-12)))) = v172
										if base.Ui32(v184) < base.Ui32(int32(25)) {
										} else {
											*(*int32)(unsafe.Add(mBase, uint32(v176)+24)) = v172
											*(*int32)(unsafe.Add(mBase, uint32(v176)+20)) = v172
											*(*int32)(unsafe.Add(mBase, uint32(v176)+16)) = v172
											*(*int32)(unsafe.Add(mBase, uint32(v176)+12)) = v172
											*(*int32)(unsafe.Add(mBase, uint32(v185+int32(-16)))) = v172
											*(*int32)(unsafe.Add(mBase, uint32(v185+int32(-20)))) = v172
											*(*int32)(unsafe.Add(mBase, uint32(v185+int32(-24)))) = v172
											*(*int32)(unsafe.Add(mBase, uint32(v185+int32(-28)))) = v172
											v220 = v176&int32(4) | int32(24)
											v221 = v184 - v220
											if base.Ui32(v221) < base.Ui32(int32(32)) {
											} else {
												v226 = base.I64_extend_i32_u(v172) * int64(4294967297)
												v229 = v221
												v230 = v176 + v220
												for {
													*(*int64)(unsafe.Add(mBase, uint32(v230)+24)) = v226
													*(*int64)(unsafe.Add(mBase, uint32(v230)+16)) = v226
													*(*int64)(unsafe.Add(mBase, uint32(v230)+8)) = v226
													*(*int64)(unsafe.Add(mBase, uint32(v230))) = v226
													v242 = v229 + int32(-32)
													if base.Ui32(int32(31)) < base.Ui32(v242) {
														v229 = v242
														v230 = v230 + int32(32)
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
				} else {
					base.MemoryFill(m, v117, v139, v133)
				}
			} else {
				v138 = F_memcpy(m, v117, v118, v133)
				mBase = m.M
			}
		} else {
			if v116&int32(255) == int32(0) {
				v131 = F_memcpy(m, v117, v33, v115<<(uint(int32(2))%32))
				mBase = m.M
			} else {
				v126 = m.G97
				v127 = *(*int32)(unsafe.Add(mBase, uint32(v126)))
				m.T0[v127].(func(*base.Module, int32, int32, int32, int32))(m, v33, v118, v117, v115)
				mBase = m.M
			}
		}
		v262 = l1 + int32(4)
		v263 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3305)))
		v264 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3305)))
		if v264 == int32(0) {
			if v263&int32(255) == int32(0) {
				base.MemoryFill(m, v18, int32(0), int32(1024))
			} else {
				v282 = F_memcpy(m, v18, v262, int32(1024))
				mBase = m.M
			}
		} else {
			if v263&int32(255) == int32(0) {
				v276 = F_memcpy(m, v18, v20, int32(1024))
				mBase = m.M
			} else {
				v272 = m.G97
				v273 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
				m.T0[v273].(func(*base.Module, int32, int32, int32, int32))(m, v20, v262, v18, int32(256))
				mBase = m.M
			}
		}
		v406 = l1 + int32(1028)
		v407 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3306)))
		v408 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3306)))
		if v408 == int32(0) {
			if v407&int32(255) == int32(0) {
				base.MemoryFill(m, v22, int32(0), int32(1024))
			} else {
				v426 = F_memcpy(m, v22, v406, int32(1024))
				mBase = m.M
			}
		} else {
			if v407&int32(255) == int32(0) {
				v420 = F_memcpy(m, v22, v24, int32(1024))
				mBase = m.M
			} else {
				v416 = m.G97
				v417 = *(*int32)(unsafe.Add(mBase, uint32(v416)))
				m.T0[v417].(func(*base.Module, int32, int32, int32, int32))(m, v24, v406, v22, int32(256))
				mBase = m.M
			}
		}
		v550 = l1 + int32(2052)
		v551 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3307)))
		v552 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3307)))
		if v552 == int32(0) {
			if v551&int32(255) == int32(0) {
				base.MemoryFill(m, v26, int32(0), int32(1024))
			} else {
				v570 = F_memcpy(m, v26, v550, int32(1024))
				mBase = m.M
			}
		} else {
			if v551&int32(255) == int32(0) {
				v564 = F_memcpy(m, v26, v28, int32(1024))
				mBase = m.M
			} else {
				v560 = m.G97
				v561 = *(*int32)(unsafe.Add(mBase, uint32(v560)))
				m.T0[v561].(func(*base.Module, int32, int32, int32, int32))(m, v28, v550, v26, int32(256))
				mBase = m.M
			}
		}
		v694 = l1 + int32(3076)
		v695 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3308)))
		v696 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3308)))
		if v696 == int32(0) {
			if v695&int32(255) == int32(0) {
				base.MemoryFill(m, v30, int32(0), int32(160))
			} else {
				v714 = F_memcpy(m, v30, v694, int32(160))
				mBase = m.M
			}
		} else {
			if v695&int32(255) == int32(0) {
				v708 = F_memcpy(m, v30, v32, int32(160))
				mBase = m.M
			} else {
				v704 = m.G97
				v705 = *(*int32)(unsafe.Add(mBase, uint32(v704)))
				m.T0[v705].(func(*base.Module, int32, int32, int32, int32))(m, v32, v694, v30, int32(40))
				mBase = m.M
			}
		}
	} else {
		v35 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3304)))
		if v35 == int32(0) {
		} else {
			v39 = *(*int32)(unsafe.Add(mBase, uint32(l2)+3236))
			v41 = int32(280)
			if int32(0) < v39 {
				v46 = int32(1)<<(uint(v39)%32) + v41
			} else {
				v46 = v41
			}
			v47 = *(*int32)(unsafe.Add(mBase, uint32(l2)))
			v48 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2)+3304)))
			if v48 == int32(0) {
				v56 = F_memcpy(m, v47, v33, v46<<(uint(int32(2))%32))
				mBase = m.M
			} else {
				v51 = m.G98
				v52 = *(*int32)(unsafe.Add(mBase, uint32(v51)))
				m.T0[v52].(func(*base.Module, int32, int32, int32))(m, v33, v47, v46)
				mBase = m.M
			}
		}
		v59 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3305)))
		if v59 == int32(0) {
		} else {
			v62 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2)+3305)))
			if v62 == int32(0) {
				v70 = F_memcpy(m, v18, v20, int32(1024))
				mBase = m.M
			} else {
				v66 = m.G98
				v67 = *(*int32)(unsafe.Add(mBase, uint32(v66)))
				m.T0[v67].(func(*base.Module, int32, int32, int32))(m, v20, v18, int32(256))
				mBase = m.M
			}
		}
		v71 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3306)))
		if v71 == int32(0) {
		} else {
			v74 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2)+3306)))
			if v74 == int32(0) {
				v82 = F_memcpy(m, v22, v24, int32(1024))
				mBase = m.M
			} else {
				v78 = m.G98
				v79 = *(*int32)(unsafe.Add(mBase, uint32(v78)))
				m.T0[v79].(func(*base.Module, int32, int32, int32))(m, v24, v22, int32(256))
				mBase = m.M
			}
		}
		v83 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3307)))
		if v83 == int32(0) {
		} else {
			v86 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2)+3307)))
			if v86 == int32(0) {
				v94 = F_memcpy(m, v26, v28, int32(1024))
				mBase = m.M
			} else {
				v90 = m.G98
				v91 = *(*int32)(unsafe.Add(mBase, uint32(v90)))
				m.T0[v91].(func(*base.Module, int32, int32, int32))(m, v28, v26, int32(256))
				mBase = m.M
			}
		}
		v95 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3308)))
		if v95 == int32(0) {
		} else {
			v98 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2)+3308)))
			if v98 == int32(0) {
				v106 = F_memcpy(m, v30, v32, int32(160))
				mBase = m.M
			} else {
				v102 = m.G98
				v103 = *(*int32)(unsafe.Add(mBase, uint32(v102)))
				m.T0[v103].(func(*base.Module, int32, int32, int32))(m, v32, v30, int32(40))
				mBase = m.M
			}
		}
	}
	v845 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+3240)))
	v847 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l1)+3240)))
	if v845 == v847 {
		v849 = v845
	} else {
		v849 = int32(-1)
	}
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+3240)) = uint16(v849)
	v851 = int32(1)
	v853 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3304)))
	if v853 != 0 {
		v857 = v851
	} else {
		v854 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3304)))
		v857 = base.B2i32(v854 != int32(0))
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l2)+3304)) = uint8(v857)
	v859 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+3242)))
	v861 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l1)+3242)))
	if v859 == v861 {
		v863 = v859
	} else {
		v863 = int32(-1)
	}
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+3242)) = uint16(v863)
	v865 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3305)))
	if v865 != 0 {
		v869 = v851
	} else {
		v866 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3305)))
		v869 = base.B2i32(v866 != int32(0))
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l2)+3305)) = uint8(v869)
	v871 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+3244)))
	v873 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l1)+3244)))
	if v871 == v873 {
		v875 = v871
	} else {
		v875 = int32(-1)
	}
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+3244)) = uint16(v875)
	v877 = int32(1)
	v879 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3306)))
	if v879 != 0 {
		v883 = v877
	} else {
		v880 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3306)))
		v883 = base.B2i32(v880 != int32(0))
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l2)+3306)) = uint8(v883)
	v885 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+3246)))
	v887 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l1)+3246)))
	if v885 == v887 {
		v889 = v885
	} else {
		v889 = int32(-1)
	}
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+3246)) = uint16(v889)
	v891 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3307)))
	if v891 != 0 {
		v895 = v877
	} else {
		v892 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3307)))
		v895 = base.B2i32(v892 != int32(0))
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l2)+3307)) = uint8(v895)
	v897 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+3248)))
	v899 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l1)+3248)))
	if v897 == v899 {
		v901 = v897
	} else {
		v901 = int32(-1)
	}
	*(*uint16)(unsafe.Add(mBase, uint32(l2)+3248)) = uint16(v901)
	v904 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3308)))
	if v904 != 0 {
		v908 = int32(1)
	} else {
		v905 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3308)))
		v908 = base.B2i32(v905 != int32(0))
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l2)+3308)) = uint8(v908)
	return
}
func F_HistogramAddSinglePixOrCopy(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
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
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v44 int32
	_ = v44
	var v50 int32
	_ = v50
	var v53 int32
	_ = v53
	var v61 int32
	_ = v61
	var v64 int32
	_ = v64
	var v68 int32
	_ = v68
	var v69 int32
	_ = v69
	var v70 int32
	_ = v70
	var v73 int32
	_ = v73
	var v77 int32
	_ = v77
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v95 int32
	_ = v95
	var v98 int32
	_ = v98
	var v102 int32
	_ = v102
	var v103 int32
	_ = v103
	var v104 int32
	_ = v104
	var v107 int32
	_ = v107
	var v111 int32
	_ = v111
	var v116 int32
	_ = v116
	var v119 int32
	_ = v119
	var v123 int32
	_ = v123
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v128 int32
	_ = v128
	var v132 int32
	_ = v132
	var v137 int32
	_ = v137
	var v143 int32
	_ = v143
	var v146 int32
	_ = v146
	v6 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
	switch v6 {
	case 0:
		v9 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+int32(7)))))
		v10 = int32(2)
		v14 = l0 + v9<<(uint(v10)%32) + int32(2052)
		v15 = *(*int32)(unsafe.Add(mBase, uint32(v14)))
		v16 = int32(1)
		*(*int32)(unsafe.Add(mBase, uint32(v14))) = v15 + v16
		v21 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+int32(6)))))
		v25 = int32(4)
		v26 = l0 + v21<<(uint(v10)%32) + v25
		v27 = *(*int32)(unsafe.Add(mBase, uint32(v26)))
		*(*int32)(unsafe.Add(mBase, uint32(v26))) = v27 + v16
		v31 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		v34 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+int32(5)))))
		v37 = v31 + v34<<(uint(v10)%32)
		v38 = *(*int32)(unsafe.Add(mBase, uint32(v37)))
		*(*int32)(unsafe.Add(mBase, uint32(v37))) = v38 + v16
		v44 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v25))))
		v143 = l0 + v44<<(uint(v10)%32) + int32(1028)
	case 1:
		v50 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		v53 = *(*int32)(unsafe.Add(mBase, uint32(l1+int32(4))))
		v143 = v50 + v53<<(uint(int32(2))%32) + int32(1120)
	default:
		v61 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l1+int32(2)))))
		if base.Ui32(int32(511)) < base.Ui32(v61) {
			v69 = int32(-1)
			v70 = v61 + v69
			v73 = base.I32_clz(v70) ^ int32(31)
			v77 = int32(1)
			v82 = int32(base.Ui32(v70)>>(uint(v73+v69)%32))&v77 | v73<<(uint(v77)%32)
		} else {
			v64 = m.G113
			v68 = int32(*(*int8)(unsafe.Add(mBase, uint32(v64+v61<<(uint(int32(1))%32)))))
			v82 = v68
		}
		v83 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		v88 = v83 + v82<<(uint(int32(2))%32) + int32(1024)
		v89 = *(*int32)(unsafe.Add(mBase, uint32(v88)))
		*(*int32)(unsafe.Add(mBase, uint32(v88))) = v89 + int32(1)
		v95 = *(*int32)(unsafe.Add(mBase, uint32(l1+int32(4))))
		if l2 != 0 {
			v116 = m.T0[l2].(func(*base.Module, int32, int32) int32)(m, l3, v95)
			mBase = m.M
			if int32(511) < v116 {
				v124 = int32(-1)
				v125 = v116 + v124
				v128 = base.I32_clz(v125) ^ int32(31)
				v132 = int32(1)
				v137 = int32(base.Ui32(v125)>>(uint(v128+v124)%32))&v132 | v128<<(uint(v132)%32)
			} else {
				v119 = m.G113
				v123 = int32(*(*int8)(unsafe.Add(mBase, uint32(v119+v116<<(uint(int32(1))%32)))))
				v137 = v123
			}
		} else {
			if int32(511) < v95 {
				v103 = int32(-1)
				v104 = v95 + v103
				v107 = base.I32_clz(v104) ^ int32(31)
				v111 = int32(1)
				v137 = int32(base.Ui32(v104)>>(uint(v107+v103)%32))&v111 | v107<<(uint(v111)%32)
			} else {
				v98 = m.G113
				v102 = int32(*(*int8)(unsafe.Add(mBase, uint32(v98+v95<<(uint(int32(1))%32)))))
				v137 = v102
			}
		}
		v143 = l0 + v137<<(uint(int32(2))%32) + int32(3076)
	}
	v146 = *(*int32)(unsafe.Add(mBase, uint32(v143)))
	*(*int32)(unsafe.Add(mBase, uint32(v143))) = v146 + int32(1)
	return
}
func F_HorizontalUnfilter_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v15 int32
	_ = v15
	var v22 int32
	_ = v22
	var v27 int32
	_ = v27
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v42 int32
	_ = v42
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v49 int32
	_ = v49
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v57 int32
	_ = v57
	var v59 int32
	_ = v59
	var v64 int32
	_ = v64
	var v71 int32
	_ = v71
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v76 int32
	_ = v76
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v82 int32
	_ = v82
	var v87 int32
	_ = v87
	if l0 != 0 {
		v10 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
		v11 = v10
	} else {
		v11 = int32(0)
	}
	if l3 < int32(1) {
	} else {
		v15 = l3 & int32(3)
		if base.Ui32(l3) < base.Ui32(int32(4)) {
			v59 = v11
			v64 = int32(0)
		} else {
			v22 = v11
			v27 = int32(0)
			for {
				v30 = l2 + v27
				v31 = l1 + v27
				v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v31))))
				v33 = v32 + v22
				*(*uint8)(unsafe.Add(mBase, uint32(v30))) = uint8(v33)
				v35 = int32(1)
				v39 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v31+v35))))
				v40 = v39 + v33
				*(*uint8)(unsafe.Add(mBase, uint32(v30+v35))) = uint8(v40)
				v42 = int32(2)
				v46 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v31+v42))))
				v47 = v46 + v40
				*(*uint8)(unsafe.Add(mBase, uint32(v30+v42))) = uint8(v47)
				v49 = int32(3)
				v53 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v31+v49))))
				v54 = v53 + v47
				*(*uint8)(unsafe.Add(mBase, uint32(v30+v49))) = uint8(v54)
				v57 = v27 + int32(4)
				if l3&int32(2147483644) != v57 {
					v22 = v54
					v27 = v57
					continue
				} else {
					break
				}
				break
			}
			v59 = v54
			v64 = v57
		}
		if v15 == int32(0) {
		} else {
			v71 = v59
			v74 = l1 + v64
			v75 = v15
			v76 = l2 + v64
			for {
				v79 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v74))))
				v80 = v79 + v71
				*(*uint8)(unsafe.Add(mBase, uint32(v76))) = uint8(v80)
				v82 = int32(1)
				v87 = v75 + int32(-1)
				if v87 != 0 {
					v71 = v80
					v74 = v74 + v82
					v75 = v87
					v76 = v76 + v82
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
