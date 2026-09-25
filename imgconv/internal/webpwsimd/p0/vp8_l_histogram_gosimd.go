//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_VP8LAllocateHistogram(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v11 int32
	_ = v11
	var v12 int64
	_ = v12
	var v13 int32
	_ = v13
	var v20 int64
	_ = v20
	var v21 int32
	_ = v21
	var v22 int64
	_ = v22
	var v32 int32
	_ = v32
	var v34 int32
	_ = v34
	var v43 base.V128
	_ = v43
	var v51 int32
	_ = v51
	var v55 int32
	_ = v55
	var v59 int32
	_ = v59
	v6 = int32(_a_F_VP8LAllocateHistogram_0)
	if int32(0) < l0 {
		v11 = int32(4)<<(uint(l0)%32) + v6
	} else {
		v11 = v6
	}
	v12 = base.I64_extend_i32_s(v11)
	v13 = int32(1)
	if v12 == int64(0) {
		v32 = F_malloc(m, base.I32_wrap_i64(v12)*v13)
		mBase = m.M
		v34 = v32
	} else {
		v20 = base.I64_div_u_s(int64(2147418112), v12)
		v21 = int32(0)
		v22 = base.I64_extend_i32_u(v13)
		if base.Ui64(int64(4294967295)) < base.Ui64(v22*v12) {
			v34 = v21
		} else {
			if base.Ui64(v20) < base.Ui64(v22) {
				v34 = v21
			} else {
				v32 = F_malloc(m, base.I32_wrap_i64(v12)*v13)
				mBase = m.M
				v34 = v32
			}
		}
	}
	if v34 == int32(0) {
	} else {
		*(*int64)(unsafe.Add(mBase, uint32(v34)+3240)) = int64(-1)
		*(*int32)(unsafe.Add(mBase, uint32(v34)+3236)) = l0
		*(*int32)(unsafe.Add(mBase, uint32(v34)+3304)) = int32(16843009)
		v43 = base.Simd_g_const(&F_VP8LAllocateHistogram__k0)
		base.Simd_g_v128_store(m, v34, int32(3256), v43)
		*(*int32)(unsafe.Add(mBase, uint32(v34))) = v34 + int32(3312)
		v51 = int32(-1)
		*(*uint16)(unsafe.Add(mBase, uint32(v34+int32(3248)))) = uint16(v51)
		v55 = int32(1)
		*(*uint8)(unsafe.Add(mBase, uint32(v34+int32(3308)))) = uint8(v55)
		v59 = int32(0)
		base.Simd_g_v128_store(m, v34+int32(3272), v59, v43)
		base.Simd_g_v128_store(m, v34+int32(3288), v59, v43)
	}
	return v34
}

var F_VP8LAllocateHistogram__k0 = [2]uint64{0x0, 0x0}

func F_VP8LAllocateHistogramSet(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v24 int64
	_ = v24
	var v25 int32
	_ = v25
	var v32 int64
	_ = v32
	var v33 int32
	_ = v33
	var v34 int64
	_ = v34
	var v44 int32
	_ = v44
	var v46 int32
	_ = v46
	var v53 int32
	_ = v53
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v65 int32
	_ = v65
	var v68 int32
	_ = v68
	var v75 int32
	_ = v75
	var v76 int32
	_ = v76
	var v78 int32
	_ = v78
	var v80 int32
	_ = v80
	var v86 int32
	_ = v86
	var v88 int32
	_ = v88
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v93 int32
	_ = v93
	var v94 int32
	_ = v94
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v105 int32
	_ = v105
	var v107 int32
	_ = v107
	var v109 int32
	_ = v109
	var v115 int32
	_ = v115
	var v117 int32
	_ = v117
	var v122 int32
	_ = v122
	var v124 int32
	_ = v124
	var v126 int32
	_ = v126
	var v134 int32
	_ = v134
	var v139 int32
	_ = v139
	var v141 int32
	_ = v141
	var v143 int32
	_ = v143
	var v163 int32
	_ = v163
	var v165 int32
	_ = v165
	var v170 int32
	_ = v170
	var v174 int32
	_ = v174
	var v175 int64
	_ = v175
	var v177 int32
	_ = v177
	var v179 base.V128
	_ = v179
	var v180 int32
	_ = v180
	var v182 int32
	_ = v182
	var v184 int32
	_ = v184
	var v186 int32
	_ = v186
	var v188 int32
	_ = v188
	var v190 int32
	_ = v190
	var v192 int32
	_ = v192
	var v194 int32
	_ = v194
	var v199 int32
	_ = v199
	var v201 int32
	_ = v201
	var v228 int32
	_ = v228
	var v235 int32
	_ = v235
	var v244 int32
	_ = v244
	var v248 int32
	_ = v248
	var v254 base.V128
	_ = v254
	var v259 int32
	_ = v259
	var v263 int32
	_ = v263
	var v267 int32
	_ = v267
	v14 = int32(4) << (uint(l1) % 32)
	v15 = int32(_a_F_VP8LAllocateHistogramSet_0)
	v19 = base.B2i32(int32(0) < l1)
	if int32(0) < l1 {
		v20 = v14 + v15
	} else {
		v20 = v15
	}
	v24 = base.I64_extend_i32_u(v20*l0 + int32(12))
	v25 = int32(1)
	if v24 == int64(0) {
		v44 = F_malloc(m, base.I32_wrap_i64(v24)*v25)
		mBase = m.M
		v46 = v44
	} else {
		v32 = base.I64_div_u_s(int64(2147418112), v24)
		v33 = int32(0)
		v34 = base.I64_extend_i32_u(v25)
		if base.Ui64(int64(4294967295)) < base.Ui64(v34*v24) {
			v46 = v33
		} else {
			if base.Ui64(v32) < base.Ui64(v34) {
				v46 = v33
			} else {
				v44 = F_malloc(m, base.I32_wrap_i64(v24)*v25)
				mBase = m.M
				v46 = v44
			}
		}
	}
	if v46 == int32(0) {
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(v46)+4)) = l0
		*(*int32)(unsafe.Add(mBase, uint32(v46))) = l0
		v53 = v46 + int32(12)
		*(*int32)(unsafe.Add(mBase, uint32(v46)+8)) = v53
		if l0 < int32(1) {
		} else {
			v61 = v53 + l0<<(uint(int32(2))%32)
			v63 = l0 + int32(-1)
			if v63 != 0 {
				v65 = int32(_a_F_VP8LAllocateHistogramSet_1)
				if int32(0) < l1 {
					v68 = v14 + v65
				} else {
					v68 = v65
				}
				v75 = int32(4)
				v76 = int32(0)
				v78 = v53
				v80 = v61
				for {
					v86 = int32(-4)
					v88 = int32(31)
					v90 = int32(-32)
					v91 = (v80 + v88) & v90
					*(*int32)(unsafe.Add(mBase, uint32(v78+v75+v86))) = v91
					v93 = *(*int32)(unsafe.Add(mBase, uint32(v46)+8))
					v94 = v93 + v75
					v97 = *(*int32)(unsafe.Add(mBase, uint32(v94+v86)))
					v98 = int32(3312)
					*(*int32)(unsafe.Add(mBase, uint32(v97))) = v91 + v98
					v105 = (v91 + v68 + v88) & v90
					*(*int32)(unsafe.Add(mBase, uint32(v94))) = v105
					v107 = *(*int32)(unsafe.Add(mBase, uint32(v46)+8))
					v109 = *(*int32)(unsafe.Add(mBase, uint32(v107+v75)))
					*(*int32)(unsafe.Add(mBase, uint32(v109))) = v105 + v98
					v115 = v105 + v68
					v117 = v76 + int32(2)
					if l0&int32(2147483646) != v117 {
						v75 = v75 + int32(8)
						v76 = v117
						v78 = v107
						v80 = v115
						continue
					} else {
						break
					}
					break
				}
				v122 = v117
				v124 = v107
				v126 = v115
			} else {
				v122 = int32(0)
				v124 = v53
				v126 = v61
			}
			if l0&int32(1) == int32(0) {
			} else {
				v134 = v122 << (uint(int32(2)) % 32)
				v139 = (v126 + int32(31)) & int32(-32)
				*(*int32)(unsafe.Add(mBase, uint32(v124+v134))) = v139
				v141 = *(*int32)(unsafe.Add(mBase, uint32(v46)+8))
				v143 = *(*int32)(unsafe.Add(mBase, uint32(v141+v134)))
				*(*int32)(unsafe.Add(mBase, uint32(v143))) = v139 + int32(3312)
			}
			if l0 < int32(1) {
			} else {
				if v63 != 0 {
					v163 = int32(0)
					v165 = int32(4)
					for {
						v170 = *(*int32)(unsafe.Add(mBase, uint32(v46)+8))
						v174 = *(*int32)(unsafe.Add(mBase, uint32(v170+v165+int32(-4))))
						v175 = int64(-1)
						*(*int64)(unsafe.Add(mBase, uint32(v174)+3240)) = v175
						v177 = int32(16843009)
						*(*int32)(unsafe.Add(mBase, uint32(v174)+3304)) = v177
						v179 = base.Simd_g_const(&F_VP8LAllocateHistogramSet__k0)
						v180 = int32(3256)
						base.Simd_g_v128_store(m, v174, v180, v179)
						v182 = int32(3248)
						v184 = int32(-1)
						*(*uint16)(unsafe.Add(mBase, uint32(v174+v182))) = uint16(v184)
						v186 = int32(3308)
						v188 = int32(1)
						*(*uint8)(unsafe.Add(mBase, uint32(v174+v186))) = uint8(v188)
						v190 = int32(3272)
						v192 = int32(0)
						base.Simd_g_v128_store(m, v174+v190, v192, v179)
						v194 = int32(3288)
						base.Simd_g_v128_store(m, v174+v194, v192, v179)
						*(*int32)(unsafe.Add(mBase, uint32(v174)+3236)) = l1
						v199 = *(*int32)(unsafe.Add(mBase, uint32(v46)+8))
						v201 = *(*int32)(unsafe.Add(mBase, uint32(v199+v165)))
						*(*int64)(unsafe.Add(mBase, uint32(v201)+3240)) = v175
						*(*int32)(unsafe.Add(mBase, uint32(v201)+3304)) = v177
						base.Simd_g_v128_store(m, v201, v180, v179)
						*(*int32)(unsafe.Add(mBase, uint32(v201)+3236)) = l1
						*(*uint16)(unsafe.Add(mBase, uint32(v201+v182))) = uint16(v184)
						*(*uint8)(unsafe.Add(mBase, uint32(v201+v186))) = uint8(v188)
						base.Simd_g_v128_store(m, v201+v190, v192, v179)
						base.Simd_g_v128_store(m, v201+v194, v192, v179)
						v228 = v163 + int32(2)
						if l0&int32(2147483646) != v228 {
							v163 = v228
							v165 = v165 + int32(8)
							continue
						} else {
							break
						}
						break
					}
					v235 = v228
				} else {
					v235 = int32(0)
				}
				if l0&int32(1) == int32(0) {
				} else {
					v244 = *(*int32)(unsafe.Add(mBase, uint32(v46)+8))
					v248 = *(*int32)(unsafe.Add(mBase, uint32(v244+v235<<(uint(int32(2))%32))))
					*(*int64)(unsafe.Add(mBase, uint32(v248)+3240)) = int64(-1)
					*(*int32)(unsafe.Add(mBase, uint32(v248)+3236)) = l1
					*(*int32)(unsafe.Add(mBase, uint32(v248)+3304)) = int32(16843009)
					v254 = base.Simd_g_const(&F_VP8LAllocateHistogramSet__k0)
					base.Simd_g_v128_store(m, v248, int32(3256), v254)
					v259 = int32(-1)
					*(*uint16)(unsafe.Add(mBase, uint32(v248+int32(3248)))) = uint16(v259)
					v263 = int32(1)
					*(*uint8)(unsafe.Add(mBase, uint32(v248+int32(3308)))) = uint8(v263)
					v267 = int32(0)
					base.Simd_g_v128_store(m, v248+int32(3272), v267, v254)
					base.Simd_g_v128_store(m, v248+int32(3288), v267, v254)
				}
			}
		}
	}
	return v46
}

var F_VP8LAllocateHistogramSet__k0 = [2]uint64{0x0, 0x0}

func F_VP8LHistogramCreate(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v7 int32
	_ = v7
	var v9 int32
	_ = v9
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v19 int32
	_ = v19
	var v24 int32
	_ = v24
	var v35 int32
	_ = v35
	var v57 int32
	_ = v57
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v69 int32
	_ = v69
	var v70 int32
	_ = v70
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v111 int64
	_ = v111
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v127 int32
	_ = v127
	var v151 int32
	_ = v151
	var v157 int32
	_ = v157
	var v159 base.V128
	_ = v159
	var v164 int32
	_ = v164
	var v171 int32
	_ = v171
	var v173 int32
	_ = v173
	var v175 int32
	_ = v175
	var v177 int32
	_ = v177
	var v178 int32
	_ = v178
	var v182 int32
	_ = v182
	var v183 int32
	_ = v183
	var v186 int32
	_ = v186
	var v191 int32
	_ = v191
	var v195 int32
	_ = v195
	var v198 int32
	_ = v198
	var v200 int32
	_ = v200
	var v202 int32
	_ = v202
	var v205 int32
	_ = v205
	var v209 int32
	_ = v209
	var v210 int32
	_ = v210
	var v211 int32
	_ = v211
	var v213 int32
	_ = v213
	var v214 int32
	_ = v214
	var v218 int32
	_ = v218
	var v219 int32
	_ = v219
	var v223 int32
	_ = v223
	var v224 int32
	_ = v224
	v7 = m.G0
	v9 = v7 - int32(16)
	m.G0 = v9
	if int32(-1) < l2 {
		v14 = l2
	} else {
		v13 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3236))
		v14 = v13
	}
	v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v16 = int32(0)
	v19 = int32(_a_F_VP8LHistogramCreate_0)
	if v16 < v14 {
		v24 = int32(4)<<(uint(v14)%32) + v19
	} else {
		v24 = v19
	}
	if base.Ui32(v24) < base.Ui32(int32(33)) {
		if v24 == int32(0) {
		} else {
			*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v16)
			v35 = l0 + v24
			*(*uint8)(unsafe.Add(mBase, uint32(v35+int32(-1)))) = uint8(v16)
			if base.Ui32(v24) < base.Ui32(int32(3)) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v16)
				*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v16)
				*(*uint8)(unsafe.Add(mBase, uint32(v35+int32(-3)))) = uint8(v16)
				*(*uint8)(unsafe.Add(mBase, uint32(v35+int32(-2)))) = uint8(v16)
				if base.Ui32(v24) < base.Ui32(int32(7)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)) = uint8(v16)
					*(*uint8)(unsafe.Add(mBase, uint32(v35+int32(-4)))) = uint8(v16)
					if base.Ui32(v24) < base.Ui32(int32(9)) {
					} else {
						v57 = int32(0)
						v60 = (v57 - l0) & int32(3)
						v61 = l0 + v60
						*(*int32)(unsafe.Add(mBase, uint32(v61))) = v57
						v69 = (v24 - v60) & int32(60)
						v70 = v61 + v69
						*(*int32)(unsafe.Add(mBase, uint32(v70+int32(-4)))) = v57
						if base.Ui32(v69) < base.Ui32(int32(9)) {
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(v61)+8)) = v57
							*(*int32)(unsafe.Add(mBase, uint32(v61)+4)) = v57
							*(*int32)(unsafe.Add(mBase, uint32(v70+int32(-8)))) = v57
							*(*int32)(unsafe.Add(mBase, uint32(v70+int32(-12)))) = v57
							if base.Ui32(v69) < base.Ui32(int32(25)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v61)+24)) = v57
								*(*int32)(unsafe.Add(mBase, uint32(v61)+20)) = v57
								*(*int32)(unsafe.Add(mBase, uint32(v61)+16)) = v57
								*(*int32)(unsafe.Add(mBase, uint32(v61)+12)) = v57
								*(*int32)(unsafe.Add(mBase, uint32(v70+int32(-16)))) = v57
								*(*int32)(unsafe.Add(mBase, uint32(v70+int32(-20)))) = v57
								*(*int32)(unsafe.Add(mBase, uint32(v70+int32(-24)))) = v57
								*(*int32)(unsafe.Add(mBase, uint32(v70+int32(-28)))) = v57
								v105 = v61&int32(4) | int32(24)
								v106 = v69 - v105
								if base.Ui32(v106) < base.Ui32(int32(32)) {
								} else {
									v111 = base.I64_extend_i32_u(v57) * int64(4294967297)
									v114 = v106
									v115 = v61 + v105
									for {
										*(*int64)(unsafe.Add(mBase, uint32(v115)+24)) = v111
										*(*int64)(unsafe.Add(mBase, uint32(v115)+16)) = v111
										*(*int64)(unsafe.Add(mBase, uint32(v115)+8)) = v111
										*(*int64)(unsafe.Add(mBase, uint32(v115))) = v111
										v127 = v114 + int32(-32)
										if base.Ui32(int32(31)) < base.Ui32(v127) {
											v114 = v127
											v115 = v115 + int32(32)
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
		base.MemoryFill(m, l0, v16, v24)
	}
	*(*int64)(unsafe.Add(mBase, uint32(l0)+3240)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v15
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3236)) = v14
	v151 = int32(-1)
	*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(3248)))) = uint16(v151)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3304)) = int32(16843009)
	v157 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(3308)))) = uint8(v157)
	v159 = base.Simd_g_const(&F_VP8LHistogramCreate__k0)
	base.Simd_g_v128_store(m, l0, int32(3256), v159)
	v164 = int32(0)
	base.Simd_g_v128_store(m, l0+int32(3272), v164, v159)
	base.Simd_g_v128_store(m, l0+int32(3288), v164, v159)
	v171 = v9 + int32(4)
	v173 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v171)+4)) = v173
	if v173 != 0 {
		v177 = *(*int32)(unsafe.Add(mBase, uint32(v173)+4))
		v178 = *(*int32)(unsafe.Add(mBase, uint32(v173)+8))
		v182 = v177 + v178<<(uint(int32(3))%32)
		v183 = v177
	} else {
		v175 = int32(0)
		v182 = v175
		v183 = v175
	}
	*(*int32)(unsafe.Add(mBase, uint32(v171)+8)) = v182
	*(*int32)(unsafe.Add(mBase, uint32(v171))) = v183
	v186 = *(*int32)(unsafe.Add(mBase, uint32(v9)+4))
	if v186 == int32(0) {
	} else {
		v191 = v186
		for {
			v195 = int32(0)
			F_HistogramAddSinglePixOrCopy(m, l0, v191, v195, v195)
			mBase = m.M
			v198 = *(*int32)(unsafe.Add(mBase, uint32(v9)+4))
			v200 = v198 + int32(8)
			*(*int32)(unsafe.Add(mBase, uint32(v9)+4)) = v200
			v202 = *(*int32)(unsafe.Add(mBase, uint32(v9)+12))
			if v200 != v202 {
				v224 = v200
			} else {
				v205 = v9 + int32(4)
				v209 = *(*int32)(unsafe.Add(mBase, uint32(v205)+4))
				v210 = *(*int32)(unsafe.Add(mBase, uint32(v209)))
				if v210 != 0 {
					v213 = *(*int32)(unsafe.Add(mBase, uint32(v210)+4))
					v214 = *(*int32)(unsafe.Add(mBase, uint32(v210)+8))
					v218 = v213
					v219 = v213 + v214<<(uint(int32(3))%32)
				} else {
					v211 = int32(0)
					v218 = v211
					v219 = v211
				}
				*(*int32)(unsafe.Add(mBase, uint32(v205)+8)) = v219
				*(*int32)(unsafe.Add(mBase, uint32(v205))) = v218
				*(*int32)(unsafe.Add(mBase, uint32(v205)+4)) = v210
				v223 = *(*int32)(unsafe.Add(mBase, uint32(v9)+4))
				v224 = v223
			}
			if v224 != 0 {
				v191 = v224
				continue
			} else {
				break
			}
			break
		}
	}
	m.G0 = v9 + int32(16)
	return
}

var F_VP8LHistogramCreate__k0 = [2]uint64{0x0, 0x0}

func F_VP8LHistogramInit(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v13 int32
	_ = v13
	var v18 int32
	_ = v18
	var v29 int32
	_ = v29
	var v51 int32
	_ = v51
	var v54 int32
	_ = v54
	var v55 int32
	_ = v55
	var v63 int32
	_ = v63
	var v64 int32
	_ = v64
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v105 int64
	_ = v105
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v121 int32
	_ = v121
	var v147 base.V128
	_ = v147
	var v152 int32
	_ = v152
	var v156 int32
	_ = v156
	var v160 int32
	_ = v160
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3236)) = l1
	if l2 == int32(0) {
	} else {
		v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		v10 = int32(0)
		v13 = int32(_a_F_VP8LHistogramInit_0)
		if v10 < l1 {
			v18 = int32(4)<<(uint(l1)%32) + v13
		} else {
			v18 = v13
		}
		if base.Ui32(v18) < base.Ui32(int32(33)) {
			if v18 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v10)
				v29 = l0 + v18
				*(*uint8)(unsafe.Add(mBase, uint32(v29+int32(-1)))) = uint8(v10)
				if base.Ui32(v18) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v10)
					*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v10)
					*(*uint8)(unsafe.Add(mBase, uint32(v29+int32(-3)))) = uint8(v10)
					*(*uint8)(unsafe.Add(mBase, uint32(v29+int32(-2)))) = uint8(v10)
					if base.Ui32(v18) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)) = uint8(v10)
						*(*uint8)(unsafe.Add(mBase, uint32(v29+int32(-4)))) = uint8(v10)
						if base.Ui32(v18) < base.Ui32(int32(9)) {
						} else {
							v51 = int32(0)
							v54 = (v51 - l0) & int32(3)
							v55 = l0 + v54
							*(*int32)(unsafe.Add(mBase, uint32(v55))) = v51
							v63 = (v18 - v54) & int32(60)
							v64 = v55 + v63
							*(*int32)(unsafe.Add(mBase, uint32(v64+int32(-4)))) = v51
							if base.Ui32(v63) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v55)+8)) = v51
								*(*int32)(unsafe.Add(mBase, uint32(v55)+4)) = v51
								*(*int32)(unsafe.Add(mBase, uint32(v64+int32(-8)))) = v51
								*(*int32)(unsafe.Add(mBase, uint32(v64+int32(-12)))) = v51
								if base.Ui32(v63) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v55)+24)) = v51
									*(*int32)(unsafe.Add(mBase, uint32(v55)+20)) = v51
									*(*int32)(unsafe.Add(mBase, uint32(v55)+16)) = v51
									*(*int32)(unsafe.Add(mBase, uint32(v55)+12)) = v51
									*(*int32)(unsafe.Add(mBase, uint32(v64+int32(-16)))) = v51
									*(*int32)(unsafe.Add(mBase, uint32(v64+int32(-20)))) = v51
									*(*int32)(unsafe.Add(mBase, uint32(v64+int32(-24)))) = v51
									*(*int32)(unsafe.Add(mBase, uint32(v64+int32(-28)))) = v51
									v99 = v55&int32(4) | int32(24)
									v100 = v63 - v99
									if base.Ui32(v100) < base.Ui32(int32(32)) {
									} else {
										v105 = base.I64_extend_i32_u(v51) * int64(4294967297)
										v108 = v100
										v109 = v55 + v99
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v109)+24)) = v105
											*(*int64)(unsafe.Add(mBase, uint32(v109)+16)) = v105
											*(*int64)(unsafe.Add(mBase, uint32(v109)+8)) = v105
											*(*int64)(unsafe.Add(mBase, uint32(v109))) = v105
											v121 = v108 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v121) {
												v108 = v121
												v109 = v109 + int32(32)
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
			base.MemoryFill(m, l0, v10, v18)
		}
		*(*int32)(unsafe.Add(mBase, uint32(l0))) = v9
		*(*int32)(unsafe.Add(mBase, uint32(l0)+3236)) = l1
	}
	*(*int64)(unsafe.Add(mBase, uint32(l0)+3240)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3304)) = int32(16843009)
	v147 = base.Simd_g_const(&F_VP8LHistogramInit__k0)
	base.Simd_g_v128_store(m, l0, int32(3256), v147)
	v152 = int32(-1)
	*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(3248)))) = uint16(v152)
	v156 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(3308)))) = uint8(v156)
	v160 = int32(0)
	base.Simd_g_v128_store(m, l0+int32(3272), v160, v147)
	base.Simd_g_v128_store(m, l0+int32(3288), v160, v147)
	return
}

var F_VP8LHistogramInit__k0 = [2]uint64{0x0, 0x0}
