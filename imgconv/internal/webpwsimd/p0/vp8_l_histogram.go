//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_VP8LFreeHistogram(m *base.Module, l0 int32) {
	F_free(m, l0)
	return
}
func F_VP8LFreeHistogramSet(m *base.Module, l0 int32) {
	F_free(m, l0)
	return
}
func F_VP8LHistogramEstimateBits(m *base.Module, l0 int32) int64 {
	mBase := m.M
	_ = mBase
	var v8 int32
	_ = v8
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v20 int64
	_ = v20
	var v23 int32
	_ = v23
	var v26 int64
	_ = v26
	var v32 int64
	_ = v32
	var v38 int64
	_ = v38
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v44 int64
	_ = v44
	var v45 int32
	_ = v45
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	v8 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3236))
	v12 = int32(280)
	if int32(0) < v10 {
		v17 = int32(1)<<(uint(v10)%32) + v12
	} else {
		v17 = v12
	}
	v18 = int32(0)
	v20 = F_PopulationCost(m, v8, v17, v18, v18)
	mBase = m.M
	v23 = int32(256)
	v26 = F_PopulationCost(m, l0+int32(4), v23, v18, v18)
	mBase = m.M
	v32 = F_PopulationCost(m, l0+int32(1028), v23, v18, v18)
	mBase = m.M
	v38 = F_PopulationCost(m, l0+int32(2052), v23, v18, v18)
	mBase = m.M
	v40 = l0 + int32(3076)
	v41 = int32(40)
	v44 = F_PopulationCost(m, v40, v41, v18, v18)
	mBase = m.M
	v45 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v49 = m.G112
	v50 = *(*int32)(unsafe.Add(mBase, uint32(v49)))
	v51 = m.T0[v50].(func(*base.Module, int32, int32) int32)(m, v45+int32(1024), int32(24))
	mBase = m.M
	v53 = *(*int32)(unsafe.Add(mBase, uint32(v49)))
	v54 = m.T0[v53].(func(*base.Module, int32, int32) int32)(m, v40, v41)
	mBase = m.M
	return base.I64_extend_i32_u(v51+v54)<<(uint(int64(23))%64) + (v44 + (v38 + (v32 + (v26 + v20))))
}
func F_VP8LHistogramSetClear(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v27 int32
	_ = v27
	var v38 int32
	_ = v38
	var v60 int32
	_ = v60
	var v63 int32
	_ = v63
	var v64 int32
	_ = v64
	var v72 int32
	_ = v72
	var v73 int32
	_ = v73
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v114 int64
	_ = v114
	var v117 int32
	_ = v117
	var v118 int32
	_ = v118
	var v130 int32
	_ = v130
	var v150 int32
	_ = v150
	var v155 int32
	_ = v155
	var v159 int32
	_ = v159
	var v162 int32
	_ = v162
	var v165 int32
	_ = v165
	var v170 int32
	_ = v170
	var v171 int32
	_ = v171
	var v177 int32
	_ = v177
	var v179 int32
	_ = v179
	var v181 int32
	_ = v181
	var v183 int32
	_ = v183
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
	var v193 int32
	_ = v193
	var v200 int32
	_ = v200
	var v202 int32
	_ = v202
	var v204 int32
	_ = v204
	var v210 int32
	_ = v210
	var v212 int32
	_ = v212
	var v215 int32
	_ = v215
	var v221 int32
	_ = v221
	var v223 int32
	_ = v223
	var v227 int32
	_ = v227
	var v232 int32
	_ = v232
	var v234 int32
	_ = v234
	var v236 int32
	_ = v236
	var v241 int32
	_ = v241
	var v249 int32
	_ = v249
	var v250 int32
	_ = v250
	var v259 int32
	_ = v259
	var v262 int32
	_ = v262
	var v265 int32
	_ = v265
	var v268 int32
	_ = v268
	var v269 int32
	_ = v269
	var v277 int32
	_ = v277
	var v283 int32
	_ = v283
	var v293 int32
	_ = v293
	var v302 int32
	_ = v302
	var v303 int32
	_ = v303
	var v308 int32
	_ = v308
	v2 = int32(0)
	v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v15 = *(*int32)(unsafe.Add(mBase, uint32(v14)))
	v16 = *(*int32)(unsafe.Add(mBase, uint32(v15)+3236))
	v17 = int32(4) << (uint(v16) % 32)
	v18 = int32(_a_F_VP8LHistogramSetClear_0)
	v22 = base.B2i32(v2 < v16)
	if v2 < v16 {
		v23 = v17 + v18
	} else {
		v23 = v18
	}
	v24 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v27 = v23*v24 + int32(12)
	if base.Ui32(v27) < base.Ui32(int32(33)) {
		if v27 == int32(0) {
		} else {
			*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v2)
			v38 = l0 + v27
			*(*uint8)(unsafe.Add(mBase, uint32(v38+int32(-1)))) = uint8(v2)
			if base.Ui32(v27) < base.Ui32(int32(3)) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v2)
				*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v2)
				*(*uint8)(unsafe.Add(mBase, uint32(v38+int32(-3)))) = uint8(v2)
				*(*uint8)(unsafe.Add(mBase, uint32(v38+int32(-2)))) = uint8(v2)
				if base.Ui32(v27) < base.Ui32(int32(7)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)) = uint8(v2)
					*(*uint8)(unsafe.Add(mBase, uint32(v38+int32(-4)))) = uint8(v2)
					if base.Ui32(v27) < base.Ui32(int32(9)) {
					} else {
						v60 = int32(0)
						v63 = (v60 - l0) & int32(3)
						v64 = l0 + v63
						*(*int32)(unsafe.Add(mBase, uint32(v64))) = v60
						v72 = (v27 - v63) & int32(60)
						v73 = v64 + v72
						*(*int32)(unsafe.Add(mBase, uint32(v73+int32(-4)))) = v60
						if base.Ui32(v72) < base.Ui32(int32(9)) {
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(v64)+8)) = v60
							*(*int32)(unsafe.Add(mBase, uint32(v64)+4)) = v60
							*(*int32)(unsafe.Add(mBase, uint32(v73+int32(-8)))) = v60
							*(*int32)(unsafe.Add(mBase, uint32(v73+int32(-12)))) = v60
							if base.Ui32(v72) < base.Ui32(int32(25)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v64)+24)) = v60
								*(*int32)(unsafe.Add(mBase, uint32(v64)+20)) = v60
								*(*int32)(unsafe.Add(mBase, uint32(v64)+16)) = v60
								*(*int32)(unsafe.Add(mBase, uint32(v64)+12)) = v60
								*(*int32)(unsafe.Add(mBase, uint32(v73+int32(-16)))) = v60
								*(*int32)(unsafe.Add(mBase, uint32(v73+int32(-20)))) = v60
								*(*int32)(unsafe.Add(mBase, uint32(v73+int32(-24)))) = v60
								*(*int32)(unsafe.Add(mBase, uint32(v73+int32(-28)))) = v60
								v108 = v64&int32(4) | int32(24)
								v109 = v72 - v108
								if base.Ui32(v109) < base.Ui32(int32(32)) {
								} else {
									v114 = base.I64_extend_i32_u(v60) * int64(4294967297)
									v117 = v109
									v118 = v64 + v108
									for {
										*(*int64)(unsafe.Add(mBase, uint32(v118)+24)) = v114
										*(*int64)(unsafe.Add(mBase, uint32(v118)+16)) = v114
										*(*int64)(unsafe.Add(mBase, uint32(v118)+8)) = v114
										*(*int64)(unsafe.Add(mBase, uint32(v118))) = v114
										v130 = v117 + int32(-32)
										if base.Ui32(int32(31)) < base.Ui32(v130) {
											v117 = v130
											v118 = v118 + int32(32)
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
		base.MemoryFill(m, l0, v2, v27)
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v24
	v150 = l0 + int32(12)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v150
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v24
	if v24 < int32(1) {
	} else {
		v155 = int32(1)
		v159 = v150 + v24<<(uint(int32(2))%32)
		if v24 == v155 {
			v215 = v2
			v221 = v150
			v223 = v159
		} else {
			v162 = int32(_a_F_VP8LHistogramSetClear_1)
			if v2 < v16 {
				v165 = v17 + v162
			} else {
				v165 = v162
			}
			v170 = int32(4)
			v171 = int32(0)
			v177 = v150
			v179 = v159
			for {
				v181 = int32(-4)
				v183 = int32(31)
				v185 = int32(-32)
				v186 = (v179 + v183) & v185
				*(*int32)(unsafe.Add(mBase, uint32(v177+v170+v181))) = v186
				v188 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				v189 = v188 + v170
				v192 = *(*int32)(unsafe.Add(mBase, uint32(v189+v181)))
				v193 = int32(3312)
				*(*int32)(unsafe.Add(mBase, uint32(v192))) = v186 + v193
				v200 = (v186 + v165 + v183) & v185
				*(*int32)(unsafe.Add(mBase, uint32(v189))) = v200
				v202 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				v204 = *(*int32)(unsafe.Add(mBase, uint32(v202+v170)))
				*(*int32)(unsafe.Add(mBase, uint32(v204))) = v200 + v193
				v210 = v200 + v165
				v212 = v171 + int32(2)
				if v24&int32(2147483646) != v212 {
					v170 = v170 + int32(8)
					v171 = v212
					v177 = v202
					v179 = v210
					continue
				} else {
					break
				}
				break
			}
			v215 = v212
			v221 = v202
			v223 = v210
		}
		if v24&v155 == int32(0) {
			v241 = v221
		} else {
			v227 = v215 << (uint(int32(2)) % 32)
			v232 = (v223 + int32(31)) & int32(-32)
			*(*int32)(unsafe.Add(mBase, uint32(v221+v227))) = v232
			v234 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
			v236 = *(*int32)(unsafe.Add(mBase, uint32(v234+v227)))
			*(*int32)(unsafe.Add(mBase, uint32(v236))) = v232 + int32(3312)
			v241 = v234
		}
		if v24 < int32(1) {
		} else {
			if base.Ui32(v24) < base.Ui32(int32(4)) {
				v283 = int32(0)
				v293 = v241 + v283<<(uint(int32(2))%32)
				v302 = v24 - v283
				for {
					v303 = *(*int32)(unsafe.Add(mBase, uint32(v293)))
					*(*int32)(unsafe.Add(mBase, uint32(v303)+3236)) = v16
					v308 = v302 + int32(-1)
					if v308 != 0 {
						v293 = v293 + int32(4)
						v302 = v308
						continue
					} else {
						break
					}
					break
				}
			} else {
				v249 = v24 & int32(2147483644)
				v250 = v241
				v259 = v249
				for {
					v262 = *(*int32)(unsafe.Add(mBase, uint32(v250+int32(12))))
					v265 = *(*int32)(unsafe.Add(mBase, uint32(v250+int32(8))))
					v268 = *(*int32)(unsafe.Add(mBase, uint32(v250+int32(4))))
					v269 = *(*int32)(unsafe.Add(mBase, uint32(v250)))
					*(*int32)(unsafe.Add(mBase, uint32(v269)+3236)) = v16
					*(*int32)(unsafe.Add(mBase, uint32(v268)+3236)) = v16
					*(*int32)(unsafe.Add(mBase, uint32(v265)+3236)) = v16
					*(*int32)(unsafe.Add(mBase, uint32(v262)+3236)) = v16
					v277 = v259 + int32(-4)
					if v277 != 0 {
						v250 = v250 + int32(16)
						v259 = v277
						continue
					} else {
						break
					}
					break
				}
				if v24 == v249 {
				} else {
					v283 = v249
					v293 = v241 + v283<<(uint(int32(2))%32)
					v302 = v24 - v283
					for {
						v303 = *(*int32)(unsafe.Add(mBase, uint32(v293)))
						*(*int32)(unsafe.Add(mBase, uint32(v303)+3236)) = v16
						v308 = v302 + int32(-1)
						if v308 != 0 {
							v293 = v293 + int32(4)
							v302 = v308
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
func F_VP8LHistogramStoreRefs(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v11 int32
	_ = v11
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v29 int32
	_ = v29
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v42 int32
	_ = v42
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
	var v55 int32
	_ = v55
	var v56 int32
	_ = v56
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	v6 = m.G0
	v8 = v6 - int32(16)
	m.G0 = v8
	v11 = v8 + int32(4)
	v13 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v11)+4)) = v13
	if v13 != 0 {
		v17 = *(*int32)(unsafe.Add(mBase, uint32(v13)+4))
		v18 = *(*int32)(unsafe.Add(mBase, uint32(v13)+8))
		v22 = v17 + v18<<(uint(int32(3))%32)
		v23 = v17
	} else {
		v15 = int32(0)
		v22 = v15
		v23 = v15
	}
	*(*int32)(unsafe.Add(mBase, uint32(v11)+8)) = v22
	*(*int32)(unsafe.Add(mBase, uint32(v11))) = v23
	v26 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
	if v26 == int32(0) {
	} else {
		v29 = v26
		for {
			F_HistogramAddSinglePixOrCopy(m, l3, v29, l1, l2)
			mBase = m.M
			v35 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
			v37 = v35 + int32(8)
			*(*int32)(unsafe.Add(mBase, uint32(v8)+4)) = v37
			v39 = *(*int32)(unsafe.Add(mBase, uint32(v8)+12))
			if v37 != v39 {
				v61 = v37
			} else {
				v42 = v8 + int32(4)
				v46 = *(*int32)(unsafe.Add(mBase, uint32(v42)+4))
				v47 = *(*int32)(unsafe.Add(mBase, uint32(v46)))
				if v47 != 0 {
					v50 = *(*int32)(unsafe.Add(mBase, uint32(v47)+4))
					v51 = *(*int32)(unsafe.Add(mBase, uint32(v47)+8))
					v55 = v50
					v56 = v50 + v51<<(uint(int32(3))%32)
				} else {
					v48 = int32(0)
					v55 = v48
					v56 = v48
				}
				*(*int32)(unsafe.Add(mBase, uint32(v42)+8)) = v56
				*(*int32)(unsafe.Add(mBase, uint32(v42))) = v55
				*(*int32)(unsafe.Add(mBase, uint32(v42)+4)) = v47
				v60 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
				v61 = v60
			}
			if v61 != 0 {
				v29 = v61
				continue
			} else {
				break
			}
			break
		}
	}
	m.G0 = v8 + int32(16)
	return
}
