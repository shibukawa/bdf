//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_VP8IteratorBytesToNz(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v3 int32
	_ = v3
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
	v2 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l0)+96))
	v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
	v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+100))
	v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)+104))
	v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+108))
	v22 = *(*int32)(unsafe.Add(mBase, uint32(l0)+112))
	v26 = *(*int32)(unsafe.Add(mBase, uint32(l0)+116))
	v30 = *(*int32)(unsafe.Add(mBase, uint32(l0)+120))
	v34 = *(*int32)(unsafe.Add(mBase, uint32(l0)+124))
	v38 = *(*int32)(unsafe.Add(mBase, uint32(l0)+128))
	v42 = *(*int32)(unsafe.Add(mBase, uint32(l0)+132))
	v46 = *(*int32)(unsafe.Add(mBase, uint32(l0)+136))
	v50 = *(*int32)(unsafe.Add(mBase, uint32(l0)+144))
	v54 = *(*int32)(unsafe.Add(mBase, uint32(l0)+152))
	*(*int32)(unsafe.Add(mBase, uint32(v2))) = v3<<(uint(int32(13))%32) | v6<<(uint(int32(12))%32) | v10<<(uint(int32(14))%32) | v14<<(uint(int32(15))%32) | v18<<(uint(int32(18))%32) | v22<<(uint(int32(19))%32) | v26<<(uint(int32(22))%32) | v30<<(uint(int32(23))%32) | v34<<(uint(int32(24))%32) | v38<<(uint(int32(3))%32) | v42<<(uint(int32(7))%32) | v46<<(uint(int32(11))%32) | v50<<(uint(int32(17))%32) | v54<<(uint(int32(21))%32)
	return
}
func F_VP8IteratorExport(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
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
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v54 int32
	_ = v54
	var v55 int32
	_ = v55
	var v60 int32
	_ = v60
	var v68 int32
	_ = v68
	var v73 int32
	_ = v73
	var v76 int32
	_ = v76
	var v83 int32
	_ = v83
	var v91 int32
	_ = v91
	var v94 int32
	_ = v94
	var v98 int32
	_ = v98
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v103 int32
	_ = v103
	var v104 int32
	_ = v104
	var v107 int32
	_ = v107
	var v109 int32
	_ = v109
	var v112 int32
	_ = v112
	var v119 int32
	_ = v119
	var v124 int32
	_ = v124
	var v126 int32
	_ = v126
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v131 int32
	_ = v131
	var v132 int32
	_ = v132
	var v136 int32
	_ = v136
	var v144 int32
	_ = v144
	var v150 int32
	_ = v150
	var v157 int32
	_ = v157
	var v162 int32
	_ = v162
	var v164 int32
	_ = v164
	var v174 int32
	_ = v174
	var v176 int32
	_ = v176
	var v178 int32
	_ = v178
	var v181 int32
	_ = v181
	var v203 int32
	_ = v203
	var v204 int32
	_ = v204
	var v210 int32
	_ = v210
	var v216 int32
	_ = v216
	var v218 int32
	_ = v218
	var v223 int32
	_ = v223
	v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v16 = *(*int32)(unsafe.Add(mBase, uint32(v15)))
	v17 = *(*int32)(unsafe.Add(mBase, uint32(v16)+64))
	if v17 == int32(0) {
	} else {
		v20 = *(*int32)(unsafe.Add(mBase, uint32(v15)+4))
		v21 = *(*int32)(unsafe.Add(mBase, uint32(v20)+12))
		v22 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
		v25 = v21 - v22<<(uint(int32(4))%32)
		v26 = int32(16)
		if v25 < v26 {
			v29 = v25
		} else {
			v29 = v26
		}
		v30 = *(*int32)(unsafe.Add(mBase, uint32(v20)+8))
		v31 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		v34 = v30 - v31<<(uint(int32(4))%32)
		v35 = int32(16)
		if v34 < v35 {
			v38 = v34
		} else {
			v38 = v35
		}
		v39 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
		v40 = *(*int32)(unsafe.Add(mBase, uint32(v20)+24))
		v41 = *(*int32)(unsafe.Add(mBase, uint32(v20)+20))
		v42 = *(*int32)(unsafe.Add(mBase, uint32(v20)+32))
		if v25 < int32(1) {
			v83 = v42
		} else {
			v47 = *(*int32)(unsafe.Add(mBase, uint32(v20)+16))
			v48 = *(*int32)(unsafe.Add(mBase, uint32(v20)+28))
			v54 = v39
			v55 = v29 + int32(1)
			v60 = v47 + (v48*v22+v31)<<(uint(int32(4))%32)
			for {
				v68 = F_memcpy(m, v60, v54, v38)
				mBase = m.M
				v73 = v55 + int32(-1)
				if base.Ui32(int32(1)) < base.Ui32(v73) {
					v54 = v54 + int32(32)
					v55 = v73
					v60 = v68 + v48
					continue
				} else {
					break
				}
				break
			}
			v76 = *(*int32)(unsafe.Add(mBase, uint32(v20)+32))
			v83 = v76
		}
		v91 = int32(1)
		v94 = (v29 + v91) >> (uint(v91) % 32)
		if v94 < v91 {
		} else {
			v98 = v39 + int32(16)
			v101 = int32(3)
			v102 = (v42*v22 + v31) << (uint(v101) % 32)
			v103 = v41 + v102
			v104 = int32(1)
			v107 = (v38 + v104) >> (uint(v104) % 32)
			v109 = v94 & v101
			if v109 != 0 {
				v112 = v98
				v119 = v109
				v124 = v103
				for {
					v126 = F_memcpy(m, v124, v112, v107)
					mBase = m.M
					v128 = v112 + int32(32)
					v129 = v126 + v83
					v131 = v119 + int32(-1)
					if v131 != 0 {
						v112 = v128
						v119 = v131
						v124 = v129
						continue
					} else {
						break
					}
					break
				}
				v132 = v128
				v136 = v94 & int32(-4)
				v144 = v129
			} else {
				v132 = v98
				v136 = v94
				v144 = v103
			}
			if base.Ui32(v94) < base.Ui32(int32(4)) {
			} else {
				v150 = v132
				v157 = v136 + int32(-1)
				v162 = v144
				for {
					v164 = F_memcpy(m, v162, v150, v107)
					mBase = m.M
					v174 = F_memcpy(m, v164+v83, v150+int32(32), v107)
					mBase = m.M
					v176 = F_memcpy(m, v174+v83, v150+int32(64), v107)
					mBase = m.M
					v178 = F_memcpy(m, v176+v83, v150+int32(96), v107)
					mBase = m.M
					v181 = v157 + int32(-4)
					if base.Ui32(v181) < base.Ui32(int32(-2)) {
						v150 = v150 + int32(128)
						v157 = v181
						v162 = v178 + v83
						continue
					} else {
						break
					}
					break
				}
			}
			v203 = *(*int32)(unsafe.Add(mBase, uint32(v20)+32))
			v204 = v39 + int32(24)
			v210 = v40 + v102
			v216 = v94 + int32(1)
			for {
				v218 = F_memcpy(m, v210, v204, v107)
				mBase = m.M
				v223 = v216 + int32(-1)
				if base.Ui32(int32(1)) < base.Ui32(v223) {
					v204 = v204 + int32(32)
					v210 = v218 + v203
					v216 = v223
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
func F_VP8IteratorInit(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v8 int32
	_ = v8
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v17 int32
	_ = v17
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v27 int32
	_ = v27
	var v44 int32
	_ = v44
	var v46 int32
	_ = v46
	var v48 int64
	_ = v48
	var v50 int32
	_ = v50
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v57 int32
	_ = v57
	var v62 int32
	_ = v62
	var v63 int64
	_ = v63
	var v69 int32
	_ = v69
	var v72 int32
	_ = v72
	var v77 int32
	_ = v77
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v84 int32
	_ = v84
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v92 int32
	_ = v92
	var v103 int32
	_ = v103
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v133 int32
	_ = v133
	var v137 int32
	_ = v137
	var v138 int32
	_ = v138
	var v173 int32
	_ = v173
	var v174 int32
	_ = v174
	var v179 int64
	_ = v179
	var v182 int32
	_ = v182
	var v183 int32
	_ = v183
	var v195 int32
	_ = v195
	var v213 int32
	_ = v213
	var v214 int32
	_ = v214
	var v215 int32
	_ = v215
	var v217 int32
	_ = v217
	var v228 int32
	_ = v228
	var v250 int32
	_ = v250
	var v253 int32
	_ = v253
	var v254 int32
	_ = v254
	var v262 int32
	_ = v262
	var v263 int32
	_ = v263
	var v298 int32
	_ = v298
	var v299 int32
	_ = v299
	var v304 int64
	_ = v304
	var v307 int32
	_ = v307
	var v308 int32
	_ = v308
	var v320 int32
	_ = v320
	var v338 int32
	_ = v338
	var v341 int32
	_ = v341
	var v342 int32
	_ = v342
	var v344 int32
	_ = v344
	var v355 int32
	_ = v355
	var v377 int32
	_ = v377
	var v380 int32
	_ = v380
	var v381 int32
	_ = v381
	var v389 int32
	_ = v389
	var v390 int32
	_ = v390
	var v425 int32
	_ = v425
	var v426 int32
	_ = v426
	var v431 int64
	_ = v431
	var v434 int32
	_ = v434
	var v435 int32
	_ = v435
	var v447 int32
	_ = v447
	v3 = int32(0)
	*(*int64)(unsafe.Add(mBase, uint32(l1))) = int64(0)
	*(*int32)(unsafe.Add(mBase, uint32(l1)+24)) = l0
	v8 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8IteratorInit[0])))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+280)) = v8
	v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+296)) = v10
	v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8IteratorInit[1])))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+304)) = v12
	*(*int32)(unsafe.Add(mBase, uint32(l1)+32)) = l0 + int32(88)
	v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8IteratorInit[2])))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+36)) = v17
	v21 = int32(-32)
	v22 = (l1 + int32(447)) & v21
	*(*int32)(unsafe.Add(mBase, uint32(l1)+8)) = v22
	v27 = (l1 + int32(360)) & v21
	*(*int32)(unsafe.Add(mBase, uint32(l1)+308)) = v27
	*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = v22 + int32(1536)
	*(*int32)(unsafe.Add(mBase, uint32(l1)+16)) = v22 + int32(1024)
	*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v22 + int32(512)
	*(*int32)(unsafe.Add(mBase, uint32(l1)+316)) = v27 + int32(48)
	*(*int32)(unsafe.Add(mBase, uint32(l1)+312)) = v27 + int32(32)
	v44 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8IteratorInit[3])))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+40)) = v44
	v46 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8IteratorInit[4])))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+28)) = v46
	v48 = *(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8IteratorInit[5])))
	*(*int64)(unsafe.Add(mBase, uint32(l1)+320)) = v48
	v50 = int32(127)
	*(*uint8)(unsafe.Add(mBase, uint32(v27)+47)) = uint8(v50)
	v52 = *(*int32)(unsafe.Add(mBase, uint32(l1)+312))
	v53 = int32(-1)
	*(*uint8)(unsafe.Add(mBase, uint32(v52+v53))) = uint8(v50)
	v57 = *(*int32)(unsafe.Add(mBase, uint32(l1)+308))
	*(*uint8)(unsafe.Add(mBase, uint32(v57+v53))) = uint8(v50)
	v62 = *(*int32)(unsafe.Add(mBase, uint32(l1)+308))
	v63 = int64(-9114861777597660799)
	*(*int64)(unsafe.Add(mBase, uint32(v62))) = v63
	*(*int64)(unsafe.Add(mBase, uint32(v62+int32(8)))) = v63
	v69 = *(*int32)(unsafe.Add(mBase, uint32(l1)+312))
	*(*int64)(unsafe.Add(mBase, uint32(v69))) = v63
	v72 = *(*int32)(unsafe.Add(mBase, uint32(l1)+316))
	*(*int64)(unsafe.Add(mBase, uint32(v72))) = v63
	*(*int32)(unsafe.Add(mBase, uint32(l1)+160)) = v3
	v77 = *(*int32)(unsafe.Add(mBase, uint32(l1)+304))
	if v77 == v3 {
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l1)+300)) = int32(0)
	}
	v82 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v83 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v84 = v82 * v83
	*(*int32)(unsafe.Add(mBase, uint32(l1)+292)) = v84
	*(*int32)(unsafe.Add(mBase, uint32(l1)+288)) = v84
	v87 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v88 = *(*int32)(unsafe.Add(mBase, uint32(v87)+uint32(_c_F_VP8IteratorInit[5])))
	v89 = int32(127)
	v90 = *(*int32)(unsafe.Add(mBase, uint32(v87)+40))
	v92 = v90 << (uint(int32(5)) % 32)
	if base.Ui32(v92) < base.Ui32(int32(33)) {
		if v92 == int32(0) {
		} else {
			*(*uint8)(unsafe.Add(mBase, uint32(v88))) = uint8(v89)
			v103 = v88 + v92
			*(*uint8)(unsafe.Add(mBase, uint32(v103+int32(-1)))) = uint8(v89)
			if base.Ui32(v92) < base.Ui32(int32(3)) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v88)+2)) = uint8(v89)
				*(*uint8)(unsafe.Add(mBase, uint32(v88)+1)) = uint8(v89)
				*(*uint8)(unsafe.Add(mBase, uint32(v103+int32(-3)))) = uint8(v89)
				*(*uint8)(unsafe.Add(mBase, uint32(v103+int32(-2)))) = uint8(v89)
				if base.Ui32(v92) < base.Ui32(int32(7)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v88)+3)) = uint8(v89)
					*(*uint8)(unsafe.Add(mBase, uint32(v103+int32(-4)))) = uint8(v89)
					if base.Ui32(v92) < base.Ui32(int32(9)) {
					} else {
						v128 = (int32(0) - v88) & int32(3)
						v129 = v88 + v128
						v133 = int32(2139062143)
						*(*int32)(unsafe.Add(mBase, uint32(v129))) = v133
						v137 = (v92 - v128) & int32(60)
						v138 = v129 + v137
						*(*int32)(unsafe.Add(mBase, uint32(v138+int32(-4)))) = v133
						if base.Ui32(v137) < base.Ui32(int32(9)) {
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(v129)+8)) = v133
							*(*int32)(unsafe.Add(mBase, uint32(v129)+4)) = v133
							*(*int32)(unsafe.Add(mBase, uint32(v138+int32(-8)))) = v133
							*(*int32)(unsafe.Add(mBase, uint32(v138+int32(-12)))) = v133
							if base.Ui32(v137) < base.Ui32(int32(25)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v129)+24)) = v133
								*(*int32)(unsafe.Add(mBase, uint32(v129)+20)) = v133
								*(*int32)(unsafe.Add(mBase, uint32(v129)+16)) = v133
								*(*int32)(unsafe.Add(mBase, uint32(v129)+12)) = v133
								*(*int32)(unsafe.Add(mBase, uint32(v138+int32(-16)))) = v133
								*(*int32)(unsafe.Add(mBase, uint32(v138+int32(-20)))) = v133
								*(*int32)(unsafe.Add(mBase, uint32(v138+int32(-24)))) = v133
								*(*int32)(unsafe.Add(mBase, uint32(v138+int32(-28)))) = v133
								v173 = v129&int32(4) | int32(24)
								v174 = v137 - v173
								if base.Ui32(v174) < base.Ui32(int32(32)) {
								} else {
									v179 = base.I64_extend_i32_u(v133) * int64(4294967297)
									v182 = v174
									v183 = v129 + v173
									for {
										*(*int64)(unsafe.Add(mBase, uint32(v183)+24)) = v179
										*(*int64)(unsafe.Add(mBase, uint32(v183)+16)) = v179
										*(*int64)(unsafe.Add(mBase, uint32(v183)+8)) = v179
										*(*int64)(unsafe.Add(mBase, uint32(v183))) = v179
										v195 = v182 + int32(-32)
										if base.Ui32(int32(31)) < base.Ui32(v195) {
											v182 = v195
											v183 = v183 + int32(32)
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
		base.MemoryFill(m, v88, v89, v92)
	}
	v213 = *(*int32)(unsafe.Add(mBase, uint32(v87)+uint32(_c_F_VP8IteratorInit[3])))
	v214 = int32(0)
	v215 = *(*int32)(unsafe.Add(mBase, uint32(v87)+40))
	v217 = v215 << (uint(int32(2)) % 32)
	if base.Ui32(v217) < base.Ui32(int32(33)) {
		if v217 == int32(0) {
		} else {
			*(*uint8)(unsafe.Add(mBase, uint32(v213))) = uint8(v214)
			v228 = v213 + v217
			*(*uint8)(unsafe.Add(mBase, uint32(v228+int32(-1)))) = uint8(v214)
			if base.Ui32(v217) < base.Ui32(int32(3)) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v213)+2)) = uint8(v214)
				*(*uint8)(unsafe.Add(mBase, uint32(v213)+1)) = uint8(v214)
				*(*uint8)(unsafe.Add(mBase, uint32(v228+int32(-3)))) = uint8(v214)
				*(*uint8)(unsafe.Add(mBase, uint32(v228+int32(-2)))) = uint8(v214)
				if base.Ui32(v217) < base.Ui32(int32(7)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v213)+3)) = uint8(v214)
					*(*uint8)(unsafe.Add(mBase, uint32(v228+int32(-4)))) = uint8(v214)
					if base.Ui32(v217) < base.Ui32(int32(9)) {
					} else {
						v250 = int32(0)
						v253 = (v250 - v213) & int32(3)
						v254 = v213 + v253
						*(*int32)(unsafe.Add(mBase, uint32(v254))) = v250
						v262 = (v217 - v253) & int32(60)
						v263 = v254 + v262
						*(*int32)(unsafe.Add(mBase, uint32(v263+int32(-4)))) = v250
						if base.Ui32(v262) < base.Ui32(int32(9)) {
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(v254)+8)) = v250
							*(*int32)(unsafe.Add(mBase, uint32(v254)+4)) = v250
							*(*int32)(unsafe.Add(mBase, uint32(v263+int32(-8)))) = v250
							*(*int32)(unsafe.Add(mBase, uint32(v263+int32(-12)))) = v250
							if base.Ui32(v262) < base.Ui32(int32(25)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v254)+24)) = v250
								*(*int32)(unsafe.Add(mBase, uint32(v254)+20)) = v250
								*(*int32)(unsafe.Add(mBase, uint32(v254)+16)) = v250
								*(*int32)(unsafe.Add(mBase, uint32(v254)+12)) = v250
								*(*int32)(unsafe.Add(mBase, uint32(v263+int32(-16)))) = v250
								*(*int32)(unsafe.Add(mBase, uint32(v263+int32(-20)))) = v250
								*(*int32)(unsafe.Add(mBase, uint32(v263+int32(-24)))) = v250
								*(*int32)(unsafe.Add(mBase, uint32(v263+int32(-28)))) = v250
								v298 = v254&int32(4) | int32(24)
								v299 = v262 - v298
								if base.Ui32(v299) < base.Ui32(int32(32)) {
								} else {
									v304 = base.I64_extend_i32_u(v250) * int64(4294967297)
									v307 = v299
									v308 = v254 + v298
									for {
										*(*int64)(unsafe.Add(mBase, uint32(v308)+24)) = v304
										*(*int64)(unsafe.Add(mBase, uint32(v308)+16)) = v304
										*(*int64)(unsafe.Add(mBase, uint32(v308)+8)) = v304
										*(*int64)(unsafe.Add(mBase, uint32(v308))) = v304
										v320 = v307 + int32(-32)
										if base.Ui32(int32(31)) < base.Ui32(v320) {
											v307 = v320
											v308 = v308 + int32(32)
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
		base.MemoryFill(m, v213, v214, v217)
	}
	v338 = *(*int32)(unsafe.Add(mBase, uint32(v87)+uint32(_c_F_VP8IteratorInit[1])))
	if v338 == int32(0) {
	} else {
		v341 = int32(0)
		v342 = *(*int32)(unsafe.Add(mBase, uint32(v87)+40))
		v344 = v342 << (uint(int32(2)) % 32)
		if base.Ui32(v344) < base.Ui32(int32(33)) {
			if v344 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v338))) = uint8(v341)
				v355 = v338 + v344
				*(*uint8)(unsafe.Add(mBase, uint32(v355+int32(-1)))) = uint8(v341)
				if base.Ui32(v344) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v338)+2)) = uint8(v341)
					*(*uint8)(unsafe.Add(mBase, uint32(v338)+1)) = uint8(v341)
					*(*uint8)(unsafe.Add(mBase, uint32(v355+int32(-3)))) = uint8(v341)
					*(*uint8)(unsafe.Add(mBase, uint32(v355+int32(-2)))) = uint8(v341)
					if base.Ui32(v344) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v338)+3)) = uint8(v341)
						*(*uint8)(unsafe.Add(mBase, uint32(v355+int32(-4)))) = uint8(v341)
						if base.Ui32(v344) < base.Ui32(int32(9)) {
						} else {
							v377 = int32(0)
							v380 = (v377 - v338) & int32(3)
							v381 = v338 + v380
							*(*int32)(unsafe.Add(mBase, uint32(v381))) = v377
							v389 = (v344 - v380) & int32(60)
							v390 = v381 + v389
							*(*int32)(unsafe.Add(mBase, uint32(v390+int32(-4)))) = v377
							if base.Ui32(v389) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v381)+8)) = v377
								*(*int32)(unsafe.Add(mBase, uint32(v381)+4)) = v377
								*(*int32)(unsafe.Add(mBase, uint32(v390+int32(-8)))) = v377
								*(*int32)(unsafe.Add(mBase, uint32(v390+int32(-12)))) = v377
								if base.Ui32(v389) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v381)+24)) = v377
									*(*int32)(unsafe.Add(mBase, uint32(v381)+20)) = v377
									*(*int32)(unsafe.Add(mBase, uint32(v381)+16)) = v377
									*(*int32)(unsafe.Add(mBase, uint32(v381)+12)) = v377
									*(*int32)(unsafe.Add(mBase, uint32(v390+int32(-16)))) = v377
									*(*int32)(unsafe.Add(mBase, uint32(v390+int32(-20)))) = v377
									*(*int32)(unsafe.Add(mBase, uint32(v390+int32(-24)))) = v377
									*(*int32)(unsafe.Add(mBase, uint32(v390+int32(-28)))) = v377
									v425 = v381&int32(4) | int32(24)
									v426 = v389 - v425
									if base.Ui32(v426) < base.Ui32(int32(32)) {
									} else {
										v431 = base.I64_extend_i32_u(v377) * int64(4294967297)
										v434 = v426
										v435 = v381 + v425
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v435)+24)) = v431
											*(*int64)(unsafe.Add(mBase, uint32(v435)+16)) = v431
											*(*int64)(unsafe.Add(mBase, uint32(v435)+8)) = v431
											*(*int64)(unsafe.Add(mBase, uint32(v435))) = v431
											v447 = v434 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v447) {
												v434 = v447
												v435 = v435 + int32(32)
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
			base.MemoryFill(m, v338, v341, v344)
		}
	}
	base.MemoryFill(m, l1+int32(168), int32(0), int32(96))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+284)) = int32(0)
	return
}
func F_VP8IteratorIsDone(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	v2 = *(*int32)(unsafe.Add(mBase, uint32(l0)+288))
	return base.B2i32(v2 < int32(1))
}
func F_VP8IteratorNext(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v15 int32
	_ = v15
	var v17 int64
	_ = v17
	var v19 int32
	_ = v19
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v26 int32
	_ = v26
	var v30 int32
	_ = v30
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v46 int32
	_ = v46
	var v53 int32
	_ = v53
	var v55 int32
	_ = v55
	var v56 int32
	_ = v56
	var v59 int32
	_ = v59
	var v63 int32
	_ = v63
	var v64 int64
	_ = v64
	var v70 int32
	_ = v70
	var v73 int32
	_ = v73
	var v76 int32
	_ = v76
	var v78 int32
	_ = v78
	var v83 int32
	_ = v83
	var v84 int32
	_ = v84
	var v87 int32
	_ = v87
	var v91 int32
	_ = v91
	var v95 int32
	_ = v95
	var v96 int32
	_ = v96
	var v99 int32
	_ = v99
	var v106 int32
	_ = v106
	v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v8 = v6 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v8
	v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v11 = *(*int32)(unsafe.Add(mBase, uint32(v10)+40))
	if v8 != v11 {
		v83 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
		v84 = int32(4)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = v83 + v84
		v87 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
		*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = v87 + v84
		v91 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
		*(*int32)(unsafe.Add(mBase, uint32(l0)+40)) = v91 + v84
		v95 = *(*int32)(unsafe.Add(mBase, uint32(l0)+320))
		v96 = int32(16)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+320)) = v95 + v96
		v99 = *(*int32)(unsafe.Add(mBase, uint32(l0)+324))
		*(*int32)(unsafe.Add(mBase, uint32(l0)+324)) = v99 + v96
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l0))) = int32(0)
		v15 = *(*int32)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_VP8IteratorNext[0])))
		*(*int32)(unsafe.Add(mBase, uint32(l0)+40)) = v15
		v17 = *(*int64)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_VP8IteratorNext[1])))
		*(*int64)(unsafe.Add(mBase, uint32(l0)+320)) = v17
		v19 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
		v21 = v19 + int32(1)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v21
		v23 = *(*int32)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_VP8IteratorNext[2])))
		v24 = *(*int32)(unsafe.Add(mBase, uint32(v10)+48))
		v26 = int32(2)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = v23 + v21*v24<<(uint(v26)%32)
		v30 = *(*int32)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_VP8IteratorNext[3])))
		*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = v30 + v21*v8<<(uint(v26)%32)
		v36 = *(*int32)(unsafe.Add(mBase, uint32(v10)+52))
		v37 = int32(-1)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = v10 + (v36+v37)&v21<<(uint(int32(5))%32) + int32(88)
		v46 = *(*int32)(unsafe.Add(mBase, uint32(l0)+316))
		if v37 < v19 {
			v53 = int32(-127)
		} else {
			v53 = int32(127)
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v46+v37))) = uint8(v53)
		v55 = *(*int32)(unsafe.Add(mBase, uint32(l0)+312))
		v56 = int32(-1)
		*(*uint8)(unsafe.Add(mBase, uint32(v55+v56))) = uint8(v53)
		v59 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		*(*uint8)(unsafe.Add(mBase, uint32(v59+v56))) = uint8(v53)
		v63 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v64 = int64(-9114861777597660799)
		*(*int64)(unsafe.Add(mBase, uint32(v63))) = v64
		*(*int64)(unsafe.Add(mBase, uint32(v63+int32(8)))) = v64
		v70 = *(*int32)(unsafe.Add(mBase, uint32(l0)+312))
		*(*int64)(unsafe.Add(mBase, uint32(v70))) = v64
		v73 = *(*int32)(unsafe.Add(mBase, uint32(l0)+316))
		*(*int64)(unsafe.Add(mBase, uint32(v73))) = v64
		v76 = int32(0)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v76
		v78 = *(*int32)(unsafe.Add(mBase, uint32(l0)+304))
		if v78 == v76 {
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l0)+300)) = int32(0)
		}
	}
	v106 = *(*int32)(unsafe.Add(mBase, uint32(l0)+288))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+288)) = v106 + int32(-1)
	return base.B2i32(int32(1) < v106)
}
func F_VP8IteratorNzToBytes(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v11 int32
	_ = v11
	var v14 int32
	_ = v14
	var v24 int32
	_ = v24
	var v34 int32
	_ = v34
	v4 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v7 = *(*int32)(unsafe.Add(mBase, uint32(v4+int32(-4))))
	v8 = *(*int32)(unsafe.Add(mBase, uint32(v4)))
	v11 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+124)) = int32(base.Ui32(v8)>>(uint(int32(24))%32)) & v11
	v14 = int32(23)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+120)) = int32(base.Ui32(v8)>>(uint(v14)%32)) & v11
	*(*int32)(unsafe.Add(mBase, uint32(l0)+116)) = int32(base.Ui32(v8)>>(uint(int32(22))%32)) & v11
	v24 = int32(19)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+112)) = int32(base.Ui32(v8)>>(uint(v24)%32)) & v11
	*(*int32)(unsafe.Add(mBase, uint32(l0)+108)) = int32(base.Ui32(v8)>>(uint(int32(18))%32)) & v11
	v34 = int32(15)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+104)) = int32(base.Ui32(v8)>>(uint(v34)%32)) & v11
	*(*int32)(unsafe.Add(mBase, uint32(l0)+100)) = int32(base.Ui32(v8)>>(uint(int32(14))%32)) & v11
	*(*int32)(unsafe.Add(mBase, uint32(l0)+96)) = int32(base.Ui32(v8)>>(uint(int32(13))%32)) & v11
	*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(base.Ui32(v8)>>(uint(int32(12))%32)) & v11
	*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = int32(base.Ui32(v7)>>(uint(v14)%32)) & v11
	*(*int32)(unsafe.Add(mBase, uint32(l0)+152)) = int32(base.Ui32(v7)>>(uint(int32(21))%32)) & v11
	*(*int32)(unsafe.Add(mBase, uint32(l0)+148)) = int32(base.Ui32(v7)>>(uint(v24)%32)) & v11
	*(*int32)(unsafe.Add(mBase, uint32(l0)+144)) = int32(base.Ui32(v7)>>(uint(int32(17))%32)) & v11
	*(*int32)(unsafe.Add(mBase, uint32(l0)+140)) = int32(base.Ui32(v7)>>(uint(v34)%32)) & v11
	*(*int32)(unsafe.Add(mBase, uint32(l0)+136)) = int32(base.Ui32(v7)>>(uint(int32(11))%32)) & v11
	*(*int32)(unsafe.Add(mBase, uint32(l0)+132)) = int32(base.Ui32(v7)>>(uint(int32(7))%32)) & v11
	*(*int32)(unsafe.Add(mBase, uint32(l0)+128)) = int32(base.Ui32(v7)>>(uint(int32(3))%32)) & v11
	return
}
func F_VP8IteratorProgress(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	var v27 int32
	_ = v27
	var v29 int32
	_ = v29
	v6 = int32(1)
	if l1 == int32(0) {
		v29 = v6
	} else {
		v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
		v10 = *(*int32)(unsafe.Add(mBase, uint32(v9)+4))
		v11 = *(*int32)(unsafe.Add(mBase, uint32(v10)+96))
		if v11 == int32(0) {
			v29 = v6
		} else {
			v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)+292))
			if int32(0) < v14 {
				v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+288))
				v21 = base.I32_div_s((v14-v18)*l1, v14)
				v22 = *(*int32)(unsafe.Add(mBase, uint32(l0)+296))
				v24 = v21 + v22
			} else {
				v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)+296))
				v24 = v17
			}
			v27 = F_WebPReportProgress(m, v10, v24, v9+int32(368))
			mBase = m.M
			v29 = v27
		}
	}
	return v29
}
func F_VP8IteratorRotateI4(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v6 int32
	_ = v6
	var v9 int32
	_ = v9
	var v13 int32
	_ = v13
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
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v34 int32
	_ = v34
	var v36 int32
	_ = v36
	var v38 int32
	_ = v38
	var v42 int32
	_ = v42
	var v44 int32
	_ = v44
	var v46 int32
	_ = v46
	var v51 int32
	_ = v51
	var v55 int32
	_ = v55
	var v61 int32
	_ = v61
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+88))
	v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)+84))
	v9 = m.G88
	v13 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v9+v5<<(uint(int32(1))%32)))))
	v14 = l1 + v13
	v15 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+96)))
	*(*uint8)(unsafe.Add(mBase, uint32(v6+int32(-4)))) = uint8(v15)
	v19 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+97)))
	*(*uint8)(unsafe.Add(mBase, uint32(v6+int32(-3)))) = uint8(v19)
	v23 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+98)))
	*(*uint8)(unsafe.Add(mBase, uint32(v6+int32(-2)))) = uint8(v23)
	v27 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+99)))
	*(*uint8)(unsafe.Add(mBase, uint32(v6+int32(-1)))) = uint8(v27)
	v29 = *(*int32)(unsafe.Add(mBase, uint32(l0)+88))
	v30 = int32(3)
	if v29&v30 == v30 {
		v42 = *(*int32)(unsafe.Add(mBase, uint32(v6+int32(4))))
		*(*int32)(unsafe.Add(mBase, uint32(v6))) = v42
	} else {
		v34 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+67)))
		*(*uint8)(unsafe.Add(mBase, uint32(v6))) = uint8(v34)
		v36 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+35)))
		*(*uint8)(unsafe.Add(mBase, uint32(v6)+1)) = uint8(v36)
		v38 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v14)+3)))
		*(*uint8)(unsafe.Add(mBase, uint32(v6)+2)) = uint8(v38)
	}
	v44 = *(*int32)(unsafe.Add(mBase, uint32(l0)+88))
	v46 = v44 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+88)) = v46
	if v46 == int32(16) {
		v61 = int32(0)
	} else {
		v51 = m.G1
		v55 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v51+int32(_a_F_VP8IteratorRotateI4_0)+v46))))
		*(*int32)(unsafe.Add(mBase, uint32(l0)+84)) = l0 + v55 + int32(44)
		v61 = int32(1)
	}
	return v61
}
func F_VP8IteratorSetCountDown(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	*(*int32)(unsafe.Add(mBase, uint32(l0)+288)) = l1
	*(*int32)(unsafe.Add(mBase, uint32(l0)+292)) = l1
	return
}
func F_VP8IteratorSetRow(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v10 int64
	_ = v10
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v36 int32
	_ = v36
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v49 int32
	_ = v49
	var v53 int32
	_ = v53
	var v54 int64
	_ = v54
	var v60 int32
	_ = v60
	var v63 int32
	_ = v63
	var v66 int32
	_ = v66
	var v68 int32
	_ = v68
	v3 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = l1
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v3
	v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v8 = *(*int32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_VP8IteratorSetRow[0])))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+40)) = v8
	v10 = *(*int64)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_VP8IteratorSetRow[1])))
	*(*int64)(unsafe.Add(mBase, uint32(l0)+320)) = v10
	v12 = *(*int32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_VP8IteratorSetRow[2])))
	v13 = *(*int32)(unsafe.Add(mBase, uint32(v7)+48))
	v15 = int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = v12 + l1*v13<<(uint(v15)%32)
	v19 = *(*int32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_VP8IteratorSetRow[3])))
	v20 = *(*int32)(unsafe.Add(mBase, uint32(v7)+40))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = v19 + v20*l1<<(uint(v15)%32)
	v26 = *(*int32)(unsafe.Add(mBase, uint32(v7)+52))
	v27 = int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = v7 + (v26+v27)&l1<<(uint(int32(5))%32) + int32(88)
	v36 = *(*int32)(unsafe.Add(mBase, uint32(l0)+316))
	if v3 < l1 {
		v43 = int32(-127)
	} else {
		v43 = int32(127)
	}
	*(*uint8)(unsafe.Add(mBase, uint32(v36+v27))) = uint8(v43)
	v45 = *(*int32)(unsafe.Add(mBase, uint32(l0)+312))
	v46 = int32(-1)
	*(*uint8)(unsafe.Add(mBase, uint32(v45+v46))) = uint8(v43)
	v49 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
	*(*uint8)(unsafe.Add(mBase, uint32(v49+v46))) = uint8(v43)
	v53 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
	v54 = int64(-9114861777597660799)
	*(*int64)(unsafe.Add(mBase, uint32(v53))) = v54
	*(*int64)(unsafe.Add(mBase, uint32(v53+int32(8)))) = v54
	v60 = *(*int32)(unsafe.Add(mBase, uint32(l0)+312))
	*(*int64)(unsafe.Add(mBase, uint32(v60))) = v54
	v63 = *(*int32)(unsafe.Add(mBase, uint32(l0)+316))
	*(*int64)(unsafe.Add(mBase, uint32(v63))) = v54
	v66 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v66
	v68 = *(*int32)(unsafe.Add(mBase, uint32(l0)+304))
	if v68 == v66 {
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l0)+300)) = int32(0)
	}
	return
}
func F_VP8IteratorStartI4(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v17 int32
	_ = v17
	var v19 int32
	_ = v19
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v29 int32
	_ = v29
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
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
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v55 int32
	_ = v55
	var v57 int32
	_ = v57
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v65 int32
	_ = v65
	var v67 int32
	_ = v67
	var v69 int32
	_ = v69
	var v71 int32
	_ = v71
	var v73 int32
	_ = v73
	var v75 int32
	_ = v75
	var v77 int32
	_ = v77
	var v79 int32
	_ = v79
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v93 int32
	_ = v93
	var v95 int32
	_ = v95
	var v96 int32
	_ = v96
	var v98 int32
	_ = v98
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v105 int32
	_ = v105
	var v108 int32
	_ = v108
	var v118 int32
	_ = v118
	var v128 int32
	_ = v128
	*(*int32)(unsafe.Add(mBase, uint32(l0)+88)) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+84)) = l0 + int32(61)
	v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
	v11 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10)+15)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+44)) = uint8(v11)
	v13 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10)+14)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+45)) = uint8(v13)
	v15 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10)+13)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+46)) = uint8(v15)
	v17 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10)+12)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+47)) = uint8(v17)
	v19 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10)+11)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+48)) = uint8(v19)
	v21 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10)+10)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+49)) = uint8(v21)
	v23 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10)+9)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+50)) = uint8(v23)
	v25 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10)+8)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+51)) = uint8(v25)
	v27 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10)+7)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+52)) = uint8(v27)
	v29 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10)+6)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+53)) = uint8(v29)
	v31 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10)+5)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+54)) = uint8(v31)
	v33 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10)+4)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+55)) = uint8(v33)
	v35 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10)+3)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+56)) = uint8(v35)
	v37 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10)+2)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+57)) = uint8(v37)
	v39 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10)+1)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+58)) = uint8(v39)
	v41 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10))))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+59)) = uint8(v41)
	v43 = int32(-1)
	v45 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10+v43))))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+60)) = uint8(v45)
	v47 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v48 = *(*int32)(unsafe.Add(mBase, uint32(l0)+320))
	v49 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48))))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+61)) = uint8(v49)
	v51 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+1)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+62)) = uint8(v51)
	v53 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+2)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+63)) = uint8(v53)
	v55 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+3)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+64)) = uint8(v55)
	v57 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+4)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+65)) = uint8(v57)
	v59 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+5)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+66)) = uint8(v59)
	v61 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+6)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+67)) = uint8(v61)
	v63 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+7)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+68)) = uint8(v63)
	v65 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+8)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+69)) = uint8(v65)
	v67 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+9)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+70)) = uint8(v67)
	v69 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+10)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+71)) = uint8(v69)
	v71 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+11)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+72)) = uint8(v71)
	v73 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+12)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+73)) = uint8(v73)
	v75 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+13)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+74)) = uint8(v75)
	v77 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+14)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+75)) = uint8(v77)
	v79 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+15)))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+76)) = uint8(v79)
	v81 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v82 = *(*int32)(unsafe.Add(mBase, uint32(v47)+40))
	if v81 < v82+v43 {
		v89 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+16)))
		*(*uint8)(unsafe.Add(mBase, uint32(l0)+77)) = uint8(v89)
		v91 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+17)))
		*(*uint8)(unsafe.Add(mBase, uint32(l0)+78)) = uint8(v91)
		v93 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+18)))
		*(*uint8)(unsafe.Add(mBase, uint32(l0)+79)) = uint8(v93)
		v95 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48)+19)))
		v96 = v95
	} else {
		*(*uint8)(unsafe.Add(mBase, uint32(l0)+79)) = uint8(v79)
		*(*uint8)(unsafe.Add(mBase, uint32(l0)+78)) = uint8(v79)
		*(*uint8)(unsafe.Add(mBase, uint32(l0)+77)) = uint8(v79)
		v96 = v79
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+80)) = uint8(v96)
	v98 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v101 = *(*int32)(unsafe.Add(mBase, uint32(v98+int32(-4))))
	v102 = *(*int32)(unsafe.Add(mBase, uint32(v98)))
	v105 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+124)) = int32(base.Ui32(v102)>>(uint(int32(24))%32)) & v105
	v108 = int32(23)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+120)) = int32(base.Ui32(v102)>>(uint(v108)%32)) & v105
	*(*int32)(unsafe.Add(mBase, uint32(l0)+116)) = int32(base.Ui32(v102)>>(uint(int32(22))%32)) & v105
	v118 = int32(19)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+112)) = int32(base.Ui32(v102)>>(uint(v118)%32)) & v105
	*(*int32)(unsafe.Add(mBase, uint32(l0)+108)) = int32(base.Ui32(v102)>>(uint(int32(18))%32)) & v105
	v128 = int32(15)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+104)) = int32(base.Ui32(v102)>>(uint(v128)%32)) & v105
	*(*int32)(unsafe.Add(mBase, uint32(l0)+100)) = int32(base.Ui32(v102)>>(uint(int32(14))%32)) & v105
	*(*int32)(unsafe.Add(mBase, uint32(l0)+96)) = int32(base.Ui32(v102)>>(uint(int32(13))%32)) & v105
	*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(base.Ui32(v102)>>(uint(int32(12))%32)) & v105
	*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = int32(base.Ui32(v101)>>(uint(v108)%32)) & v105
	*(*int32)(unsafe.Add(mBase, uint32(l0)+152)) = int32(base.Ui32(v101)>>(uint(int32(21))%32)) & v105
	*(*int32)(unsafe.Add(mBase, uint32(l0)+148)) = int32(base.Ui32(v101)>>(uint(v118)%32)) & v105
	*(*int32)(unsafe.Add(mBase, uint32(l0)+144)) = int32(base.Ui32(v101)>>(uint(int32(17))%32)) & v105
	*(*int32)(unsafe.Add(mBase, uint32(l0)+140)) = int32(base.Ui32(v101)>>(uint(v128)%32)) & v105
	*(*int32)(unsafe.Add(mBase, uint32(l0)+136)) = int32(base.Ui32(v101)>>(uint(int32(11))%32)) & v105
	*(*int32)(unsafe.Add(mBase, uint32(l0)+132)) = int32(base.Ui32(v101)>>(uint(int32(7))%32)) & v105
	*(*int32)(unsafe.Add(mBase, uint32(l0)+128)) = int32(base.Ui32(v101)>>(uint(int32(3))%32)) & v105
	return
}
