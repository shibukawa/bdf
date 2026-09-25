//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_Quantize2Blocks_C(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	v5 = m.G1
	v8 = *(*int32)(unsafe.Add(mBase, uint32(v5)+uint32(_c_F_Quantize2Blocks_C[0])))
	v9 = m.T0[v8].(func(*base.Module, int32, int32, int32) int32)(m, l0, l1, l2)
	mBase = m.M
	v10 = int32(32)
	v14 = *(*int32)(unsafe.Add(mBase, uint32(v5)+uint32(_c_F_Quantize2Blocks_C[0])))
	v15 = m.T0[v14].(func(*base.Module, int32, int32, int32) int32)(m, l0+v10, l1+v10, l2)
	mBase = m.M
	return v9 | v15<<(uint(int32(1))%32)
}
func F_QuantizeBlock_C(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v26 int32
	_ = v26
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v39 int32
	_ = v39
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
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
	var v56 int32
	_ = v56
	var v58 int32
	_ = v58
	var v61 int32
	_ = v61
	var v64 int32
	_ = v64
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
	var v87 int32
	_ = v87
	var v92 int32
	_ = v92
	v26 = l1
	v32 = int32(-1)
	v33 = int32(0)
	for {
		v39 = m.G1
		v43 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v39+int32(_a_F_QuantizeBlock_C_0)+v33))))
		v45 = v43 << (uint(int32(1)) % 32)
		v46 = l0 + v45
		v47 = int32(*(*int16)(unsafe.Add(mBase, uint32(v46))))
		v49 = v47 >> (uint(int32(31)) % 32)
		v53 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2+int32(192)+v45))))
		v54 = v47 ^ v49 - v49 + v53
		v56 = v43 << (uint(int32(2)) % 32)
		v58 = *(*int32)(unsafe.Add(mBase, uint32(l2+int32(128)+v56)))
		if base.Ui32(v54) <= base.Ui32(v58) {
			v83 = int32(0)
			*(*uint16)(unsafe.Add(mBase, uint32(v26))) = uint16(v83)
			*(*uint16)(unsafe.Add(mBase, uint32(v46))) = uint16(v83)
			v87 = v32
		} else {
			v61 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2+v45))))
			v64 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l2+int32(32)+v45))))
			v67 = *(*int32)(unsafe.Add(mBase, uint32(l2+int32(64)+v56)))
			v70 = int32(base.Ui32(v54*v64+v67) >> (uint(int32(17)) % 32))
			v71 = int32(2047)
			if base.Ui32(v70) < base.Ui32(v71) {
				v74 = v70
			} else {
				v74 = v71
			}
			if v47 < int32(0) {
				v78 = int32(0) - v74
			} else {
				v78 = v74
			}
			v79 = v61 * v78
			*(*uint16)(unsafe.Add(mBase, uint32(v46))) = uint16(v79)
			*(*uint16)(unsafe.Add(mBase, uint32(v26))) = uint16(v78)
			if v78 != 0 {
				v82 = v33
			} else {
				v82 = v32
			}
			v87 = v82
		}
		v92 = v33 + int32(1)
		if v92 != int32(16) {
			v26 = v26 + int32(2)
			v32 = v87
			v33 = v92
			continue
		} else {
			break
		}
		break
	}
	return int32(base.Ui32(v87^int32(-1)) >> (uint(int32(31)) % 32))
}
func F___qsort_r(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v17 int32
	_ = v17
	var v21 int32
	_ = v21
	var v27 int32
	_ = v27
	var v40 int32
	_ = v40
	var __phi40 int32
	_ = __phi40
	var v41 int32
	_ = v41
	var __phi41 int32
	_ = __phi41
	var v42 int32
	_ = v42
	var __phi42 int32
	_ = __phi42
	var v47 int32
	_ = v47
	var v53 int32
	_ = v53
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
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
	var v78 int32
	_ = v78
	var v85 int32
	_ = v85
	var v86 int32
	_ = v86
	var v98 int32
	_ = v98
	var v102 int32
	_ = v102
	var v115 int32
	_ = v115
	var v119 int32
	_ = v119
	var v125 int32
	_ = v125
	var v126 int32
	_ = v126
	var v129 int32
	_ = v129
	var v132 int32
	_ = v132
	var v136 int32
	_ = v136
	var v140 int32
	_ = v140
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v146 int32
	_ = v146
	var v149 int32
	_ = v149
	var v151 int32
	_ = v151
	var v159 int32
	_ = v159
	var v171 int32
	_ = v171
	var v172 int32
	_ = v172
	var v186 int32
	_ = v186
	var v188 int32
	_ = v188
	var v189 int32
	_ = v189
	var v190 int32
	_ = v190
	var v198 int32
	_ = v198
	var v199 int32
	_ = v199
	var v204 int32
	_ = v204
	var v206 int32
	_ = v206
	var v209 int32
	_ = v209
	var v212 int32
	_ = v212
	var v214 int32
	_ = v214
	var v217 int32
	_ = v217
	var v230 int32
	_ = v230
	var v232 int32
	_ = v232
	var v236 int32
	_ = v236
	var v239 int32
	_ = v239
	var v244 int32
	_ = v244
	var v247 int32
	_ = v247
	var v263 int32
	_ = v263
	var v264 int32
	_ = v264
	var v266 int32
	_ = v266
	var v270 int32
	_ = v270
	var v271 int32
	_ = v271
	var v272 int32
	_ = v272
	var v273 int32
	_ = v273
	var v274 int32
	_ = v274
	var v280 int32
	_ = v280
	var v285 int32
	_ = v285
	var v286 int32
	_ = v286
	var v287 int32
	_ = v287
	var v289 int32
	_ = v289
	v15 = m.G0
	v17 = v15 - int32(208)
	m.G0 = v17
	*(*int64)(unsafe.Add(mBase, uint32(v17)+8)) = int64(1)
	v21 = l2 * l1
	if v21 == int32(0) {
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(v17)+16)) = l2
		*(*int32)(unsafe.Add(mBase, uint32(v17)+20)) = l2
		v27 = int32(0) - l2
		__phi40 = v17 + int32(16) | int32(8)
		__phi41 = l2
		__phi42 = l2
		v40 = __phi40
		v41 = __phi41
		v42 = __phi42
		for {
			v47 = v42 + l2 + v41
			*(*int32)(unsafe.Add(mBase, uint32(v40))) = v47
			if base.Ui32(v47) < base.Ui32(v21) {
				__phi40 = v40 + int32(4)
				__phi41 = v47
				__phi42 = v41
				v40 = __phi40
				v41 = __phi41
				v42 = __phi42
				continue
			} else {
				break
			}
			break
		}
		v53 = l0 + v21 + v27
		if base.Ui32(l0) < base.Ui32(v53) {
			v59 = *(*int32)(unsafe.Add(mBase, uint32(v17)+12))
			v60 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
			v62 = int32(1)
			v64 = l0
			v65 = v60
			v70 = v59
			v72 = v62
			v73 = v62
			v74 = l2 * (l1 + int32(-1))
			v76 = int32(0)
			for {
				v78 = int32(3)
				if v73&v78 != v78 {
					v98 = v72 + int32(-1)
					v102 = *(*int32)(unsafe.Add(mBase, uint32(v17+int32(16)+v98<<(uint(int32(2))%32))))
					if base.Ui32(v102) < base.Ui32(v74) {
						F_sift(m, v64, l2, l3, l4, v72, v17+int32(16))
						mBase = m.M
					} else {
						F_trinkle(m, v64, l2, l3, l4, v17+int32(8), v72, int32(0), v17+int32(16))
						mBase = m.M
					}
					if v72 != int32(1) {
						v125 = base.B2i32(base.Ui32(int32(31)) < base.Ui32(v98))
						if base.Ui32(int32(31)) < base.Ui32(v98) {
							v126 = v65
						} else {
							v126 = v70
						}
						if base.Ui32(int32(31)) < base.Ui32(v98) {
							v129 = v72 + int32(-33)
						} else {
							v129 = v98
						}
						if base.Ui32(int32(31)) < base.Ui32(v98) {
							v132 = int32(0)
						} else {
							v132 = v73
						}
						v136 = v126<<(uint(v129)%32) | int32(base.Ui32(v132)>>(uint(int32(32)-v129)%32))
						*(*int32)(unsafe.Add(mBase, uint32(v17)+12)) = v136
						v140 = v136
						v141 = int32(1)
						v142 = v132 << (uint(v129) % 32)
					} else {
						v115 = int32(1)
						v119 = v76<<(uint(v115)%32) | int32(base.Ui32(v73)>>(uint(int32(31))%32))
						*(*int32)(unsafe.Add(mBase, uint32(v17)+12)) = v119
						v140 = v119
						v141 = int32(0)
						v142 = v73 << (uint(v115) % 32)
					}
				} else {
					F_sift(m, v64, l2, l3, l4, v72, v17+int32(16))
					mBase = m.M
					v85 = int32(2)
					v86 = int32(base.Ui32(v76) >> (uint(v85) % 32))
					*(*int32)(unsafe.Add(mBase, uint32(v17)+12)) = v86
					v140 = v86
					v141 = v72 + v85
					v142 = v76<<(uint(int32(30))%32) | int32(base.Ui32(v73)>>(uint(v85)%32))
				}
				v146 = v142 | int32(1)
				*(*int32)(unsafe.Add(mBase, uint32(v17)+8)) = v146
				v149 = v64 + l2
				if base.Ui32(v149) < base.Ui32(v53) {
					v64 = v149
					v65 = v146
					v70 = v140
					v72 = v141
					v73 = v146
					v74 = v74 - l2
					v76 = v140
					continue
				} else {
					break
				}
				break
			}
			v151 = v149
			v159 = v141
		} else {
			v151 = l0
			v159 = int32(1)
		}
		F_trinkle(m, v151, l2, l3, l4, v17+int32(8), v159, int32(0), v17+int32(16))
		mBase = m.M
		v171 = *(*int32)(unsafe.Add(mBase, uint32(v17)+12))
		v172 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
		if v159 != int32(1) {
			v186 = v171
			v188 = v159
			v189 = v172
			v190 = v151 + v27
			for {
				if int32(1) < v188 {
					v209 = int32(2)
					v212 = int32(base.Ui32(v189) >> (uint(int32(30)) % 32))
					v214 = int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(v17)+12)) = int32(base.Ui32(v186<<(uint(v209)%32)|v212) >> (uint(v214) % 32))
					v217 = int32(31)
					*(*int32)(unsafe.Add(mBase, uint32(v17)+8)) = v212<<(uint(v217)%32) | v189<<(uint(v214)%32)&int32(2147483646) ^ int32(3)
					v230 = v17 + int32(16)
					v232 = v188 + int32(-2)
					v236 = *(*int32)(unsafe.Add(mBase, uint32(v230+v232<<(uint(v209)%32))))
					v239 = v17 + int32(8)
					F_trinkle(m, v190-v236, l2, l3, l4, v239, v188+int32(-1), v214, v230)
					mBase = m.M
					v244 = *(*int32)(unsafe.Add(mBase, uint32(v17)+12))
					v247 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
					*(*int32)(unsafe.Add(mBase, uint32(v17)+12)) = v244<<(uint(v214)%32) | int32(base.Ui32(v247)>>(uint(v217)%32))
					*(*int32)(unsafe.Add(mBase, uint32(v17)+8)) = v247<<(uint(v214)%32) | v214
					F_trinkle(m, v190, l2, l3, l4, v239, v232, v214, v230)
					mBase = m.M
					v263 = *(*int32)(unsafe.Add(mBase, uint32(v17)+12))
					v264 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
					v285 = v263
					v286 = v232
					v287 = v264
				} else {
					v198 = base.I32_ctz(v189 + int32(-1))
					if v198 != 0 {
						v206 = int32(32)
						if base.Ui32(v206) <= base.Ui32(v198) {
							v266 = v206
							v270 = v266 + int32(-32)
							v271 = int32(0)
							v272 = v186
							v273 = v266
						} else {
							v270 = v198
							v271 = v186
							v272 = v189
							v273 = v198
						}
					} else {
						v199 = base.I32_ctz(v186)
						if v199 == int32(0) {
							v204 = int32(0)
							v270 = v204
							v271 = v186
							v272 = v189
							v273 = v204
						} else {
							v266 = v199 + int32(32)
							v270 = v266 + int32(-32)
							v271 = int32(0)
							v272 = v186
							v273 = v266
						}
					}
					v274 = int32(base.Ui32(v271) >> (uint(v270) % 32))
					*(*int32)(unsafe.Add(mBase, uint32(v17)+12)) = v274
					v280 = v271<<(uint(int32(32)-v270)%32) | int32(base.Ui32(v272)>>(uint(v270)%32))
					*(*int32)(unsafe.Add(mBase, uint32(v17)+8)) = v280
					v285 = v274
					v286 = v273 + v188
					v287 = v280
				}
				v289 = v190 + v27
				if v286 != int32(1) {
					v186 = v285
					v188 = v286
					v189 = v287
					v190 = v289
					continue
				} else {
				}
				if v287 != int32(1) {
					v186 = v285
					v188 = v286
					v189 = v287
					v190 = v289
					continue
				} else {
				}
				if v285 != 0 {
					v186 = v285
					v188 = v286
					v189 = v287
					v190 = v289
					continue
				} else {
					break
				}
				break
			}
		} else {
			if v172 != int32(1) {
				v186 = v171
				v188 = v159
				v189 = v172
				v190 = v151 + v27
				for {
					if int32(1) < v188 {
						v209 = int32(2)
						v212 = int32(base.Ui32(v189) >> (uint(int32(30)) % 32))
						v214 = int32(1)
						*(*int32)(unsafe.Add(mBase, uint32(v17)+12)) = int32(base.Ui32(v186<<(uint(v209)%32)|v212) >> (uint(v214) % 32))
						v217 = int32(31)
						*(*int32)(unsafe.Add(mBase, uint32(v17)+8)) = v212<<(uint(v217)%32) | v189<<(uint(v214)%32)&int32(2147483646) ^ int32(3)
						v230 = v17 + int32(16)
						v232 = v188 + int32(-2)
						v236 = *(*int32)(unsafe.Add(mBase, uint32(v230+v232<<(uint(v209)%32))))
						v239 = v17 + int32(8)
						F_trinkle(m, v190-v236, l2, l3, l4, v239, v188+int32(-1), v214, v230)
						mBase = m.M
						v244 = *(*int32)(unsafe.Add(mBase, uint32(v17)+12))
						v247 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
						*(*int32)(unsafe.Add(mBase, uint32(v17)+12)) = v244<<(uint(v214)%32) | int32(base.Ui32(v247)>>(uint(v217)%32))
						*(*int32)(unsafe.Add(mBase, uint32(v17)+8)) = v247<<(uint(v214)%32) | v214
						F_trinkle(m, v190, l2, l3, l4, v239, v232, v214, v230)
						mBase = m.M
						v263 = *(*int32)(unsafe.Add(mBase, uint32(v17)+12))
						v264 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
						v285 = v263
						v286 = v232
						v287 = v264
					} else {
						v198 = base.I32_ctz(v189 + int32(-1))
						if v198 != 0 {
							v206 = int32(32)
							if base.Ui32(v206) <= base.Ui32(v198) {
								v266 = v206
								v270 = v266 + int32(-32)
								v271 = int32(0)
								v272 = v186
								v273 = v266
							} else {
								v270 = v198
								v271 = v186
								v272 = v189
								v273 = v198
							}
						} else {
							v199 = base.I32_ctz(v186)
							if v199 == int32(0) {
								v204 = int32(0)
								v270 = v204
								v271 = v186
								v272 = v189
								v273 = v204
							} else {
								v266 = v199 + int32(32)
								v270 = v266 + int32(-32)
								v271 = int32(0)
								v272 = v186
								v273 = v266
							}
						}
						v274 = int32(base.Ui32(v271) >> (uint(v270) % 32))
						*(*int32)(unsafe.Add(mBase, uint32(v17)+12)) = v274
						v280 = v271<<(uint(int32(32)-v270)%32) | int32(base.Ui32(v272)>>(uint(v270)%32))
						*(*int32)(unsafe.Add(mBase, uint32(v17)+8)) = v280
						v285 = v274
						v286 = v273 + v188
						v287 = v280
					}
					v289 = v190 + v27
					if v286 != int32(1) {
						v186 = v285
						v188 = v286
						v189 = v287
						v190 = v289
						continue
					} else {
					}
					if v287 != int32(1) {
						v186 = v285
						v188 = v286
						v189 = v287
						v190 = v289
						continue
					} else {
					}
					if v285 != 0 {
						v186 = v285
						v188 = v286
						v189 = v287
						v190 = v289
						continue
					} else {
						break
					}
					break
				}
			} else {
				if v171 == int32(0) {
				} else {
					v186 = v171
					v188 = v159
					v189 = v172
					v190 = v151 + v27
					for {
						if int32(1) < v188 {
							v209 = int32(2)
							v212 = int32(base.Ui32(v189) >> (uint(int32(30)) % 32))
							v214 = int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(v17)+12)) = int32(base.Ui32(v186<<(uint(v209)%32)|v212) >> (uint(v214) % 32))
							v217 = int32(31)
							*(*int32)(unsafe.Add(mBase, uint32(v17)+8)) = v212<<(uint(v217)%32) | v189<<(uint(v214)%32)&int32(2147483646) ^ int32(3)
							v230 = v17 + int32(16)
							v232 = v188 + int32(-2)
							v236 = *(*int32)(unsafe.Add(mBase, uint32(v230+v232<<(uint(v209)%32))))
							v239 = v17 + int32(8)
							F_trinkle(m, v190-v236, l2, l3, l4, v239, v188+int32(-1), v214, v230)
							mBase = m.M
							v244 = *(*int32)(unsafe.Add(mBase, uint32(v17)+12))
							v247 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
							*(*int32)(unsafe.Add(mBase, uint32(v17)+12)) = v244<<(uint(v214)%32) | int32(base.Ui32(v247)>>(uint(v217)%32))
							*(*int32)(unsafe.Add(mBase, uint32(v17)+8)) = v247<<(uint(v214)%32) | v214
							F_trinkle(m, v190, l2, l3, l4, v239, v232, v214, v230)
							mBase = m.M
							v263 = *(*int32)(unsafe.Add(mBase, uint32(v17)+12))
							v264 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
							v285 = v263
							v286 = v232
							v287 = v264
						} else {
							v198 = base.I32_ctz(v189 + int32(-1))
							if v198 != 0 {
								v206 = int32(32)
								if base.Ui32(v206) <= base.Ui32(v198) {
									v266 = v206
									v270 = v266 + int32(-32)
									v271 = int32(0)
									v272 = v186
									v273 = v266
								} else {
									v270 = v198
									v271 = v186
									v272 = v189
									v273 = v198
								}
							} else {
								v199 = base.I32_ctz(v186)
								if v199 == int32(0) {
									v204 = int32(0)
									v270 = v204
									v271 = v186
									v272 = v189
									v273 = v204
								} else {
									v266 = v199 + int32(32)
									v270 = v266 + int32(-32)
									v271 = int32(0)
									v272 = v186
									v273 = v266
								}
							}
							v274 = int32(base.Ui32(v271) >> (uint(v270) % 32))
							*(*int32)(unsafe.Add(mBase, uint32(v17)+12)) = v274
							v280 = v271<<(uint(int32(32)-v270)%32) | int32(base.Ui32(v272)>>(uint(v270)%32))
							*(*int32)(unsafe.Add(mBase, uint32(v17)+8)) = v280
							v285 = v274
							v286 = v273 + v188
							v287 = v280
						}
						v289 = v190 + v27
						if v286 != int32(1) {
							v186 = v285
							v188 = v286
							v189 = v287
							v190 = v289
							continue
						} else {
						}
						if v287 != int32(1) {
							v186 = v285
							v188 = v286
							v189 = v287
							v190 = v289
							continue
						} else {
						}
						if v285 != 0 {
							v186 = v285
							v188 = v286
							v189 = v287
							v190 = v289
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
	m.G0 = v17 + int32(208)
	return
}
func F_qsort(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	F___qsort_r(m, l0, l1, l2, int32(344), l3)
	return
}
