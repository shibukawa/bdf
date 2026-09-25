//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	"unsafe"
)

func F_WebPConvertARGBToUV_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v41 int32
	_ = v41
	var v43 int32
	_ = v43
	var v48 int32
	_ = v48
	var v51 int32
	_ = v51
	var v59 int32
	_ = v59
	var v63 int32
	_ = v63
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v75 int32
	_ = v75
	var v77 int32
	_ = v77
	var v78 int32
	_ = v78
	var v92 int32
	_ = v92
	var v97 int32
	_ = v97
	var v106 int32
	_ = v106
	var v107 int32
	_ = v107
	var v113 int32
	_ = v113
	var v116 int32
	_ = v116
	var v117 int32
	_ = v117
	var v119 int32
	_ = v119
	var v121 int32
	_ = v121
	var v126 int32
	_ = v126
	var v129 int32
	_ = v129
	var v137 int32
	_ = v137
	var v141 int32
	_ = v141
	var v149 int32
	_ = v149
	var v150 int32
	_ = v150
	var v153 int32
	_ = v153
	var v155 int32
	_ = v155
	var v157 int32
	_ = v157
	var v162 int32
	_ = v162
	var v164 int32
	_ = v164
	var v177 int32
	_ = v177
	var v182 int32
	_ = v182
	var v187 int32
	_ = v187
	var v194 int32
	_ = v194
	var v209 int32
	_ = v209
	var v212 int32
	_ = v212
	var v213 int32
	_ = v213
	var v214 int32
	_ = v214
	var v219 int32
	_ = v219
	var v226 int32
	_ = v226
	var v230 int32
	_ = v230
	var v232 int32
	_ = v232
	var v233 int32
	_ = v233
	var v245 int32
	_ = v245
	var v252 int32
	_ = v252
	var v253 int32
	_ = v253
	var v255 int32
	_ = v255
	var v258 int32
	_ = v258
	var v260 int32
	_ = v260
	var v261 int32
	_ = v261
	var v266 int32
	_ = v266
	v14 = int32(1)
	v15 = l3 >> (uint(v14) % 32)
	if v14 <= v15 {
		if l4 == int32(0) {
			v106 = int32(0)
			v107 = l0
			for {
				v113 = l1 + v106
				v116 = *(*int32)(unsafe.Add(mBase, uint32(v107+int32(4))))
				v117 = int32(15)
				v119 = int32(510)
				v121 = *(*int32)(unsafe.Add(mBase, uint32(v107)))
				v126 = int32(base.Ui32(v116)>>(uint(v117)%32))&v119 + int32(base.Ui32(v121)>>(uint(v117)%32))&v119
				v129 = int32(7)
				v137 = int32(base.Ui32(v116)>>(uint(v129)%32))&v119 + int32(base.Ui32(v121)>>(uint(v129)%32))&v119
				v141 = int32(1)
				v149 = v116<<(uint(v141)%32)&v119 + v121<<(uint(v141)%32)&v119
				v150 = int32(_a_F_WebPConvertARGBToUV_C_0)
				v153 = int32(33685504)
				v155 = int32(18)
				v157 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v113))))
				v162 = int32(base.Ui32(int32(base.Ui32(v126*int32(-9719)+v137*int32(-19081)+v149*v150+v153)>>(uint(v155)%32))+v157+v141) >> (uint(v141) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v113))) = uint8(v162)
				v164 = l2 + v106
				v177 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v164))))
				v182 = int32(base.Ui32(int32(base.Ui32(v126*v150+v137*int32(-24116)+v149*int32(-4684)+v153)>>(uint(v155)%32))+v177+v141) >> (uint(v141) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v164))) = uint8(v182)
				v187 = v106 + v141
				if v15 != v187 {
					v106 = v187
					v107 = v107 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			v194 = v15
		} else {
			v28 = int32(0)
			v29 = l0
			for {
				v38 = *(*int32)(unsafe.Add(mBase, uint32(v29+int32(4))))
				v39 = int32(15)
				v41 = int32(510)
				v43 = *(*int32)(unsafe.Add(mBase, uint32(v29)))
				v48 = int32(base.Ui32(v38)>>(uint(v39)%32))&v41 + int32(base.Ui32(v43)>>(uint(v39)%32))&v41
				v51 = int32(7)
				v59 = int32(base.Ui32(v38)>>(uint(v51)%32))&v41 + int32(base.Ui32(v43)>>(uint(v51)%32))&v41
				v63 = int32(1)
				v71 = v38<<(uint(v63)%32)&v41 + v43<<(uint(v63)%32)&v41
				v72 = int32(_a_F_WebPConvertARGBToUV_C_0)
				v75 = int32(33685504)
				v77 = int32(18)
				v78 = int32(base.Ui32(v48*int32(67099145)+v59*int32(67089783)+v71*v72+v75) >> (uint(v77) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(l1+v28))) = uint8(v78)
				v92 = int32(base.Ui32(v48*v72+v59*int32(67084748)+v71*int32(67104180)+v75) >> (uint(v77) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(l2+v28))) = uint8(v92)
				v97 = v28 + v63
				if v15 != v97 {
					v28 = v97
					v29 = v29 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			v194 = v15
		}
	} else {
		v194 = int32(0)
	}
	if l3&int32(1) == int32(0) {
		return
	} else {
		v209 = *(*int32)(unsafe.Add(mBase, uint32(l0+v194<<(uint(int32(3))%32))))
		v212 = int32(1020)
		v213 = int32(base.Ui32(v209)>>(uint(int32(14))%32)) & v212
		v214 = int32(_a_F_WebPConvertARGBToUV_C_0)
		v219 = int32(base.Ui32(v209)>>(uint(int32(6))%32)) & v212
		v226 = v209 << (uint(int32(2)) % 32) & v212
		v230 = int32(33685504)
		v232 = int32(18)
		v233 = int32(base.Ui32(v213*v214+v219*int32(-24116)+v226*int32(-4684)+v230) >> (uint(v232) % 32))
		v245 = int32(base.Ui32(v213*int32(-9719)+v219*int32(-19081)+v226*v214+v230) >> (uint(v232) % 32))
		if l4 == int32(0) {
			v252 = l1 + v194
			v253 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v252))))
			v255 = int32(1)
			v258 = int32(base.Ui32(v245+v253+v255) >> (uint(v255) % 32))
			*(*uint8)(unsafe.Add(mBase, uint32(v252))) = uint8(v258)
			v260 = l2 + v194
			v261 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v260))))
			v266 = int32(base.Ui32(v233+v261+v255) >> (uint(v255) % 32))
			*(*uint8)(unsafe.Add(mBase, uint32(v260))) = uint8(v266)
			return
		} else {
			*(*uint8)(unsafe.Add(mBase, uint32(l2+v194))) = uint8(v233)
			*(*uint8)(unsafe.Add(mBase, uint32(l1+v194))) = uint8(v245)
			return
		}
	}
}
func F_WebPConvertRGBA32ToUV_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v19 int32
	_ = v19
	var v22 int32
	_ = v22
	var v26 int32
	_ = v26
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v41 int32
	_ = v41
	var v54 int32
	_ = v54
	var v55 int32
	_ = v55
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v71 int32
	_ = v71
	if l3 < int32(1) {
	} else {
		v11 = l0
		v12 = l1
		v13 = l2
		v14 = l3
		for {
			v19 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v11))))
			v22 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v11)+2)))
			v26 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v11)+4)))
			v33 = (v19*int32(-9719) + v22*int32(-19081) + v26*int32(_a_F_WebPConvertRGBA32ToUV_C_0) + int32(33685504)) >> (uint(int32(18)) % 32)
			v34 = int32(0)
			if v34 < v33 {
				v37 = v33
			} else {
				v37 = v34
			}
			v38 = int32(255)
			if v37 < v38 {
				v41 = v37
			} else {
				v41 = v38
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v12))) = uint8(v41)
			v54 = (v19*int32(_a_F_WebPConvertRGBA32ToUV_C_0) + v22*int32(-24116) + v26*int32(-4684) + int32(33685504)) >> (uint(int32(18)) % 32)
			v55 = int32(0)
			if v55 < v54 {
				v58 = v54
			} else {
				v58 = v55
			}
			v59 = int32(255)
			if v58 < v59 {
				v62 = v58
			} else {
				v62 = v59
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v13))) = uint8(v62)
			v64 = int32(1)
			v71 = v14 + int32(-1)
			if v71 != 0 {
				v11 = v11 + int32(8)
				v12 = v12 + v64
				v13 = v13 + v64
				v14 = v71
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_WebPMultARGBRow_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v13 int32
	_ = v13
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v33 int32
	_ = v33
	var v40 int32
	_ = v40
	var v57 int32
	_ = v57
	var v63 int32
	_ = v63
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v68 int32
	_ = v68
	var v74 int32
	_ = v74
	var v77 int32
	_ = v77
	var v78 int32
	_ = v78
	var v81 int32
	_ = v81
	var v88 int32
	_ = v88
	var v95 int32
	_ = v95
	var v112 int32
	_ = v112
	var v118 int32
	_ = v118
	if l1 < int32(1) {
	} else {
		if l2 == int32(0) {
			v64 = l0
			v65 = l1
			for {
				v68 = *(*int32)(unsafe.Add(mBase, uint32(v64)))
				if base.Ui32(int32(-16777217)) < base.Ui32(v68) {
				} else {
					if base.Ui32(v68) < base.Ui32(int32(16777216)) {
						v112 = int32(0)
					} else {
						v74 = int32(24)
						v77 = int32(base.Ui32(v68)>>(uint(v74)%32)) * int32(_a_F_WebPMultARGBRow_C_0)
						v78 = int32(255)
						v81 = int32(8388608)
						v88 = int32(8)
						v95 = int32(16)
						v112 = int32(base.Ui32(v77*(v68&v78)+v81)>>(uint(v74)%32)) | v68&int32(-16777216) | int32(base.Ui32(v77*(int32(base.Ui32(v68)>>(uint(v88)%32))&v78)+v81)>>(uint(v95)%32))&int32(_a_F_WebPMultARGBRow_C_1) | int32(base.Ui32(v77*(int32(base.Ui32(v68)>>(uint(v95)%32))&v78)+v81)>>(uint(v88)%32))&int32(16711680)
					}
					*(*int32)(unsafe.Add(mBase, uint32(v64))) = v112
				}
				v118 = v65 + int32(-1)
				if v118 != 0 {
					v64 = v64 + int32(4)
					v65 = v118
					continue
				} else {
					break
				}
				break
			}
		} else {
			v9 = l0
			v10 = l1
			for {
				v13 = *(*int32)(unsafe.Add(mBase, uint32(v9)))
				if base.Ui32(int32(-16777217)) < base.Ui32(v13) {
				} else {
					if base.Ui32(v13) < base.Ui32(int32(16777216)) {
						v57 = int32(0)
					} else {
						v19 = int32(-16777216)
						v20 = int32(24)
						v22 = base.I32_div_u_s(v19, int32(base.Ui32(v13)>>(uint(v20)%32)))
						v23 = int32(255)
						v26 = int32(8388608)
						v33 = int32(8)
						v40 = int32(16)
						v57 = int32(base.Ui32(v22*(v13&v23)+v26)>>(uint(v20)%32)) | v13&v19 | int32(base.Ui32(v22*(int32(base.Ui32(v13)>>(uint(v33)%32))&v23)+v26)>>(uint(v40)%32))&int32(_a_F_WebPMultARGBRow_C_1) | int32(base.Ui32(v22*(int32(base.Ui32(v13)>>(uint(v40)%32))&v23)+v26)>>(uint(v33)%32))&int32(16711680)
					}
					*(*int32)(unsafe.Add(mBase, uint32(v9))) = v57
				}
				v63 = v10 + int32(-1)
				if v63 != 0 {
					v9 = v9 + int32(4)
					v10 = v63
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
func F_WebPMultRow_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v13 int32
	_ = v13
	var v16 int32
	_ = v16
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v40 int32
	_ = v40
	var v42 int32
	_ = v42
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v56 int32
	_ = v56
	if l2 < int32(1) {
	} else {
		if l3 == int32(0) {
			v33 = l0
			v34 = l1
			v35 = l2
			for {
				v37 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v34))))
				if v37 == int32(255) {
				} else {
					if v37 != 0 {
						v42 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v33))))
						v49 = int32(base.Ui32(v37*v42*int32(_a_F_WebPMultRow_C_0)+int32(8388608)) >> (uint(int32(24)) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v33))) = uint8(v49)
					} else {
						v40 = int32(0)
						*(*uint8)(unsafe.Add(mBase, uint32(v33))) = uint8(v40)
					}
				}
				v51 = int32(1)
				v56 = v35 + int32(-1)
				if v56 != 0 {
					v33 = v33 + v51
					v34 = v34 + v51
					v35 = v56
					continue
				} else {
					break
				}
				break
			}
		} else {
			v9 = l0
			v10 = l1
			v11 = l2
			for {
				v13 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10))))
				if v13 == int32(255) {
				} else {
					if v13 != 0 {
						v19 = base.I32_div_u_s(int32(-16777216), v13)
						v20 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v9))))
						v25 = int32(base.Ui32(v19*v20+int32(8388608)) >> (uint(int32(24)) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v9))) = uint8(v25)
					} else {
						v16 = int32(0)
						*(*uint8)(unsafe.Add(mBase, uint32(v9))) = uint8(v16)
					}
				}
				v27 = int32(1)
				v32 = v11 + int32(-1)
				if v32 != 0 {
					v9 = v9 + v27
					v10 = v10 + v27
					v11 = v32
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
func F_WriteImage(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v11 int32
	_ = v11
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v31 int64
	_ = v31
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v56 int32
	_ = v56
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v67 int32
	_ = v67
	var v68 int32
	_ = v68
	var v72 int32
	_ = v72
	var v75 int32
	_ = v75
	var v82 int32
	_ = v82
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v108 int32
	_ = v108
	var v118 int32
	_ = v118
	var v119 int32
	_ = v119
	var v121 int32
	_ = v121
	var v125 int32
	_ = v125
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v134 int64
	_ = v134
	var v138 int64
	_ = v138
	var v140 int32
	_ = v140
	var v142 int64
	_ = v142
	var v149 int32
	_ = v149
	var v150 int32
	_ = v150
	var v151 int32
	_ = v151
	var v153 int32
	_ = v153
	var v154 int32
	_ = v154
	var v156 int32
	_ = v156
	var v158 int32
	_ = v158
	var v159 int32
	_ = v159
	var v161 int32
	_ = v161
	var v162 int32
	_ = v162
	var v165 int32
	_ = v165
	var v169 int32
	_ = v169
	var v172 int32
	_ = v172
	var v175 int32
	_ = v175
	var v178 int32
	_ = v178
	var v179 int32
	_ = v179
	var v182 int32
	_ = v182
	var v183 int32
	_ = v183
	var v185 int32
	_ = v185
	var v190 int32
	_ = v190
	var v193 int32
	_ = v193
	var v194 int32
	_ = v194
	var v196 int32
	_ = v196
	var v203 int32
	_ = v203
	v9 = m.G0
	v11 = v9 - int32(32)
	m.G0 = v11
	v21 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	v22 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v23 = v21 - v22
	v25 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
	v31 = base.I64_extend_i32_u(v23) + base.I64_extend_i32_u((v25+int32(7))>>(uint(int32(3))%32))
	if base.Ui64(v31) < base.Ui64(int64(4294967296)) {
		v37 = base.I32_wrap_i64(v31)
		v38 = *(*int32)(unsafe.Add(mBase, uint32(l1)+16))
		v39 = v38 - v22
		if v38 == v22 {
			v46 = int32(base.Ui32(v39*int32(3)) >> (uint(int32(1)) % 32))
			if base.Ui32(v37) < base.Ui32(v46) {
				v48 = v46
			} else {
				v48 = v37
			}
			v52 = v48&int32(-1024) + int32(1024)
			v53 = F_WebPSafeMalloc(m, int64(1), v52)
			mBase = m.M
			if v53 != 0 {
				if v21 == v22 {
				} else {
					v58 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
					v59 = F_memcpy(m, v53, v58, v23)
					mBase = m.M
				}
				v60 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
				F_WebPSafeFree(m, v60)
				mBase = m.M
				*(*int32)(unsafe.Add(mBase, uint32(l1)+16)) = v53 + v52
				*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v53 + v23
				*(*int32)(unsafe.Add(mBase, uint32(l1)+8)) = v53
				v67 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
				v68 = v67
				if v68 < int32(1) {
				} else {
					v72 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
					v75 = v72
					for {
						v82 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
						*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v82 + int32(1)
						*(*uint8)(unsafe.Add(mBase, uint32(v82))) = uint8(v75)
						v87 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
						v88 = int32(8)
						v89 = int32(base.Ui32(v87) >> (uint(v88) % 32))
						*(*int32)(unsafe.Add(mBase, uint32(l1))) = v89
						v91 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
						*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v91 + int32(-8)
						if v88 < v91 {
							v75 = v89
							continue
						} else {
							break
						}
						break
					}
				}
				*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = int32(0)
				v108 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
				v118 = v108
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = int32(1)
				v56 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
				v118 = v56
			}
		} else {
			if base.Ui32(v37) <= base.Ui32(v39) {
				v68 = v25
				if v68 < int32(1) {
				} else {
					v72 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
					v75 = v72
					for {
						v82 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
						*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v82 + int32(1)
						*(*uint8)(unsafe.Add(mBase, uint32(v82))) = uint8(v75)
						v87 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
						v88 = int32(8)
						v89 = int32(base.Ui32(v87) >> (uint(v88) % 32))
						*(*int32)(unsafe.Add(mBase, uint32(l1))) = v89
						v91 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
						*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v91 + int32(-8)
						if v88 < v91 {
							v75 = v89
							continue
						} else {
							break
						}
						break
					}
				}
				*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = int32(0)
				v108 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
				v118 = v108
			} else {
				v46 = int32(base.Ui32(v39*int32(3)) >> (uint(int32(1)) % 32))
				if base.Ui32(v37) < base.Ui32(v46) {
					v48 = v46
				} else {
					v48 = v37
				}
				v52 = v48&int32(-1024) + int32(1024)
				v53 = F_WebPSafeMalloc(m, int64(1), v52)
				mBase = m.M
				if v53 != 0 {
					if v21 == v22 {
					} else {
						v58 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
						v59 = F_memcpy(m, v53, v58, v23)
						mBase = m.M
					}
					v60 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
					F_WebPSafeFree(m, v60)
					mBase = m.M
					*(*int32)(unsafe.Add(mBase, uint32(l1)+16)) = v53 + v52
					*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v53 + v23
					*(*int32)(unsafe.Add(mBase, uint32(l1)+8)) = v53
					v67 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
					v68 = v67
					if v68 < int32(1) {
					} else {
						v72 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
						v75 = v72
						for {
							v82 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
							*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v82 + int32(1)
							*(*uint8)(unsafe.Add(mBase, uint32(v82))) = uint8(v75)
							v87 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
							v88 = int32(8)
							v89 = int32(base.Ui32(v87) >> (uint(v88) % 32))
							*(*int32)(unsafe.Add(mBase, uint32(l1))) = v89
							v91 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
							*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v91 + int32(-8)
							if v88 < v91 {
								v75 = v89
								continue
							} else {
								break
							}
							break
						}
					}
					*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = int32(0)
					v108 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
					v118 = v108
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = int32(1)
					v56 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
					v118 = v56
				}
			}
		}
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = int32(1)
		v36 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
		v118 = v36
	}
	v119 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l2))) = v119
	v121 = *(*int32)(unsafe.Add(mBase, uint32(l1)+20))
	if v121 == v119 {
		v128 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
		v129 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
		v130 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
		v131 = m.G1
		v134 = *(*int64)(unsafe.Add(mBase, uint32(v131)+uint32(_c_F_WriteImage[0])))
		*(*int64)(unsafe.Add(mBase, uint32(v11))) = v134
		v138 = *(*int64)(unsafe.Add(mBase, uint32(v131)+uint32(_c_F_WriteImage[1])))
		*(*int64)(unsafe.Add(mBase, uint32(v11)+13)) = v138
		v140 = int32(8)
		v142 = *(*int64)(unsafe.Add(mBase, uint32(v131)+uint32(_c_F_WriteImage[2])))
		*(*int64)(unsafe.Add(mBase, uint32(v11)+8)) = v142
		v149 = (v130+int32(7))>>(uint(int32(3))%32) + (v129 - v128)
		v150 = int32(1)
		v151 = v149 + v150
		v153 = v151 & v150
		v154 = v149 + v153
		v156 = v154 + int32(13)
		*(*uint8)(unsafe.Add(mBase, uint32(v11)+4)) = uint8(v156)
		v158 = int32(24)
		v159 = int32(base.Ui32(v156) >> (uint(v158) % 32))
		*(*uint8)(unsafe.Add(mBase, uint32(v11)+7)) = uint8(v159)
		v161 = int32(16)
		v162 = int32(base.Ui32(v156) >> (uint(v161) % 32))
		*(*uint8)(unsafe.Add(mBase, uint32(v11)+6)) = uint8(v162)
		v165 = int32(base.Ui32(v156) >> (uint(v140) % 32))
		*(*uint8)(unsafe.Add(mBase, uint32(v11)+5)) = uint8(v165)
		*(*uint8)(unsafe.Add(mBase, uint32(v11)+16)) = uint8(v151)
		v169 = int32(base.Ui32(v151) >> (uint(v158) % 32))
		*(*uint8)(unsafe.Add(mBase, uint32(v11)+19)) = uint8(v169)
		v172 = int32(base.Ui32(v151) >> (uint(v161) % 32))
		*(*uint8)(unsafe.Add(mBase, uint32(v11)+18)) = uint8(v172)
		v175 = int32(base.Ui32(v151) >> (uint(v140) % 32))
		*(*uint8)(unsafe.Add(mBase, uint32(v11)+17)) = uint8(v175)
		v178 = *(*int32)(unsafe.Add(mBase, uint32(l0)+72))
		v179 = m.T0[v178].(func(*base.Module, int32, int32, int32) int32)(m, v11, int32(21), l0)
		mBase = m.M
		if v179 == int32(0) {
			v185 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
			if v185 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(8)
			}
			v203 = int32(0)
		} else {
			v182 = *(*int32)(unsafe.Add(mBase, uint32(l0)+72))
			v183 = m.T0[v182].(func(*base.Module, int32, int32, int32) int32)(m, v118, v149, l0)
			mBase = m.M
			if v183 != 0 {
				if v153 == int32(0) {
					*(*int32)(unsafe.Add(mBase, uint32(l2))) = v154 + int32(21)
					v203 = int32(1)
				} else {
					v190 = int32(0)
					*(*uint8)(unsafe.Add(mBase, uint32(v11))) = uint8(v190)
					v193 = *(*int32)(unsafe.Add(mBase, uint32(l0)+72))
					v194 = m.T0[v193].(func(*base.Module, int32, int32, int32) int32)(m, v11, int32(1), l0)
					mBase = m.M
					if v194 != 0 {
						*(*int32)(unsafe.Add(mBase, uint32(l2))) = v154 + int32(21)
						v203 = int32(1)
					} else {
						v196 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
						if v196 != 0 {
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(8)
						}
						v203 = int32(0)
					}
				}
			} else {
				v185 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
				if v185 != 0 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(8)
				}
				v203 = int32(0)
			}
		}
	} else {
		v125 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
		if v125 != 0 {
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(1)
		}
		v203 = int32(0)
	}
	m.G0 = v11 + int32(32)
	return v203
}
func F___wasm_call_ctors(m *base.Module) {
	return
}
func F_wrapper_cmp(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	var v4 int32
	_ = v4
	v4 = m.T0[l2].(func(*base.Module, int32, int32) int32)(m, l0, l1)
	return v4
}
