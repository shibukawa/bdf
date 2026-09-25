//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	"unsafe"
)

func F_RD4_C(m *base.Module, l0 int32) {
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
	var v11 int32
	_ = v11
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	var v29 int32
	_ = v29
	var v37 int32
	_ = v37
	var v40 int32
	_ = v40
	var v46 int32
	_ = v46
	var v52 int32
	_ = v52
	var v59 int32
	_ = v59
	var v62 int32
	_ = v62
	var v65 int32
	_ = v65
	var v73 int32
	_ = v73
	var v84 int32
	_ = v84
	var v94 int32
	_ = v94
	v7 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+31)))
	v8 = int32(2)
	v9 = v7 + v8
	v10 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+63)))
	v11 = int32(1)
	v14 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+95)))
	v17 = int32(base.Ui32(v9+v10<<(uint(v11)%32)+v14) >> (uint(v8) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+96)) = uint8(v17)
	v21 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-1)))))
	v23 = v21 + v8
	v29 = int32(base.Ui32(v10+(v23+v7<<(uint(v11)%32))) >> (uint(v8) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+97)) = uint8(v29)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+64)) = uint8(v29)
	v37 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-33)))))
	v40 = int32(base.Ui32(v9+v21<<(uint(v11)%32)+v37) >> (uint(v8) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+98)) = uint8(v40)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+65)) = uint8(v40)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+32)) = uint8(v40)
	v46 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-32)))))
	v52 = int32(base.Ui32(v23+v46+v37<<(uint(v11)%32)) >> (uint(v8) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+99)) = uint8(v52)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+66)) = uint8(v52)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+33)) = uint8(v52)
	*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v52)
	v59 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-29)))))
	v62 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-30)))))
	v65 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-31)))))
	v73 = int32(base.Ui32(v37+v65+v46<<(uint(v11)%32)+v8) >> (uint(v8) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+67)) = uint8(v73)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+34)) = uint8(v73)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v73)
	v84 = int32(base.Ui32(v46+v62+v65<<(uint(v11)%32)+v8) >> (uint(v8) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+35)) = uint8(v84)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v84)
	v94 = int32(base.Ui32(v65+v59+v62<<(uint(v11)%32)+v8) >> (uint(v8) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)) = uint8(v94)
	return
}
func F_ReconstructUV(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) int32 {
	mBase := m.M
	_ = mBase
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
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
	var v33 int32
	_ = v33
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v42 int32
	_ = v42
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v55 int32
	_ = v55
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v63 int32
	_ = v63
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v77 int32
	_ = v77
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v85 int32
	_ = v85
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v95 int32
	_ = v95
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
	var v111 int32
	_ = v111
	var v112 int32
	_ = v112
	var v116 int32
	_ = v116
	var v118 int32
	_ = v118
	var v119 int32
	_ = v119
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v128 int32
	_ = v128
	var v132 int32
	_ = v132
	var v134 int32
	_ = v134
	var v137 int32
	_ = v137
	var v139 int32
	_ = v139
	var v141 int32
	_ = v141
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v146 int32
	_ = v146
	var v150 int32
	_ = v150
	var v151 int32
	_ = v151
	var v155 int32
	_ = v155
	var v157 int32
	_ = v157
	var v158 int32
	_ = v158
	var v161 int32
	_ = v161
	var v166 int32
	_ = v166
	var v168 int32
	_ = v168
	var v169 int32
	_ = v169
	var v174 int32
	_ = v174
	var v177 int32
	_ = v177
	var v179 int32
	_ = v179
	var v181 int32
	_ = v181
	var v184 int32
	_ = v184
	var v185 int32
	_ = v185
	var v187 int32
	_ = v187
	var v191 int32
	_ = v191
	var v192 int32
	_ = v192
	var v196 int32
	_ = v196
	var v198 int32
	_ = v198
	var v199 int32
	_ = v199
	var v202 int32
	_ = v202
	var v203 int32
	_ = v203
	var v207 int32
	_ = v207
	var v209 int32
	_ = v209
	var v216 int32
	_ = v216
	var v218 int32
	_ = v218
	var v220 int32
	_ = v220
	var v222 int32
	_ = v222
	var v223 int32
	_ = v223
	var v225 int32
	_ = v225
	var v229 int32
	_ = v229
	var v230 int32
	_ = v230
	var v234 int32
	_ = v234
	var v236 int32
	_ = v236
	var v237 int32
	_ = v237
	var v241 int32
	_ = v241
	var v246 int32
	_ = v246
	var v248 int32
	_ = v248
	var v251 int32
	_ = v251
	var v253 int32
	_ = v253
	var v256 int32
	_ = v256
	var v257 int32
	_ = v257
	var v259 int32
	_ = v259
	var v265 int32
	_ = v265
	var v267 int32
	_ = v267
	var v269 int32
	_ = v269
	var v271 int32
	_ = v271
	var v274 int32
	_ = v274
	var v275 int32
	_ = v275
	var v277 int32
	_ = v277
	var v281 int32
	_ = v281
	var v282 int32
	_ = v282
	var v286 int32
	_ = v286
	var v288 int32
	_ = v288
	var v290 int32
	_ = v290
	var v292 int32
	_ = v292
	var v293 int32
	_ = v293
	var v298 int32
	_ = v298
	var v302 int32
	_ = v302
	var v304 int32
	_ = v304
	var v307 int32
	_ = v307
	var v309 int32
	_ = v309
	var v311 int32
	_ = v311
	var v313 int32
	_ = v313
	var v314 int32
	_ = v314
	var v316 int32
	_ = v316
	var v320 int32
	_ = v320
	var v321 int32
	_ = v321
	var v325 int32
	_ = v325
	var v327 int32
	_ = v327
	var v328 int32
	_ = v328
	var v331 int32
	_ = v331
	var v336 int32
	_ = v336
	var v338 int32
	_ = v338
	var v339 int32
	_ = v339
	var v344 int32
	_ = v344
	var v347 int32
	_ = v347
	var v349 int32
	_ = v349
	var v351 int32
	_ = v351
	var v354 int32
	_ = v354
	var v355 int32
	_ = v355
	var v357 int32
	_ = v357
	var v361 int32
	_ = v361
	var v362 int32
	_ = v362
	var v366 int32
	_ = v366
	var v368 int32
	_ = v368
	var v369 int32
	_ = v369
	var v372 int32
	_ = v372
	var v373 int32
	_ = v373
	var v377 int32
	_ = v377
	var v379 int32
	_ = v379
	var v386 int32
	_ = v386
	var v388 int32
	_ = v388
	var v390 int32
	_ = v390
	var v392 int32
	_ = v392
	var v393 int32
	_ = v393
	var v395 int32
	_ = v395
	var v399 int32
	_ = v399
	var v400 int32
	_ = v400
	var v404 int32
	_ = v404
	var v406 int32
	_ = v406
	var v408 int32
	_ = v408
	var v411 int32
	_ = v411
	var v415 int32
	_ = v415
	var v417 int32
	_ = v417
	var v433 int32
	_ = v433
	var v434 int32
	_ = v434
	var v435 int32
	_ = v435
	var v436 int32
	_ = v436
	var v439 int32
	_ = v439
	var v440 int32
	_ = v440
	var v443 int32
	_ = v443
	var v444 int32
	_ = v444
	var v445 int32
	_ = v445
	var v446 int32
	_ = v446
	var v449 int32
	_ = v449
	var v450 int32
	_ = v450
	var v451 int32
	_ = v451
	var v456 int32
	_ = v456
	var v461 int32
	_ = v461
	var v466 int32
	_ = v466
	v22 = m.G0
	v24 = v22 - int32(256)
	m.G0 = v24
	v26 = m.G1
	v27 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v28 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
	v29 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28))))
	v30 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v33 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	v39 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v26+int32(_a_F_ReconstructUV_0)+l3<<(uint(int32(1))%32)))))
	v40 = v33 + v39
	v41 = m.G43
	v42 = *(*int32)(unsafe.Add(mBase, uint32(v41)))
	m.T0[v42].(func(*base.Module, int32, int32, int32))(m, v30+int32(16), v40, v24)
	mBase = m.M
	v46 = int32(128)
	v47 = v40 + v46
	v49 = v24 + int32(64)
	v50 = *(*int32)(unsafe.Add(mBase, uint32(v41)))
	m.T0[v50].(func(*base.Module, int32, int32, int32))(m, v30+int32(144), v47, v49)
	mBase = m.M
	v55 = v40 + int32(8)
	v57 = v24 + v46
	v58 = *(*int32)(unsafe.Add(mBase, uint32(v41)))
	m.T0[v58].(func(*base.Module, int32, int32, int32))(m, v30+int32(24), v55, v57)
	mBase = m.M
	v63 = v40 + int32(136)
	v65 = v24 + int32(192)
	v66 = *(*int32)(unsafe.Add(mBase, uint32(v41)))
	m.T0[v66].(func(*base.Module, int32, int32, int32))(m, v30+int32(152), v63, v65)
	mBase = m.M
	v74 = v27 + int32(base.Ui32(v29)>>(uint(int32(5))%32))&int32(3)*int32(744)
	v76 = v74 + int32(408)
	v77 = *(*int32)(unsafe.Add(mBase, uint32(l0)+304))
	if v77 == int32(0) {
	} else {
		v81 = *(*int32)(unsafe.Add(mBase, uint32(v76)+576))
		v82 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+300)))
		v83 = int32(3)
		v85 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		v88 = v77 + v85<<(uint(int32(2))%32)
		v89 = int32(*(*int8)(unsafe.Add(mBase, uint32(v88))))
		v95 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v24))))
		v97 = base.I32_extend16_s(int32(base.Ui32(v82<<(uint(v83)%32)+v89*int32(7))>>(uint(v83)%32)) + v95)
		v99 = v97 >> (uint(int32(31)) % 32)
		v101 = v97 ^ v99 - v99
		if v81 < v101 {
			v104 = int32(0)
			v105 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76)+480)))
			v107 = *(*int32)(unsafe.Add(mBase, uint32(v76)+512))
			v111 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76)+448)))
			v112 = int32(base.Ui32(v105*v101+v107)>>(uint(int32(17))%32)) * v111
			if v97 < v104 {
				v116 = v104 - v112
			} else {
				v116 = v112
			}
			v118 = v101 - v112
			v119 = v116
		} else {
			v118 = v101
			v119 = int32(0)
		}
		*(*uint16)(unsafe.Add(mBase, uint32(v24))) = uint16(v119)
		v122 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v24)+32)))
		v123 = int32(*(*int8)(unsafe.Add(mBase, uint32(v88)+1)))
		v128 = int32(0)
		if v97 < v128 {
			v132 = v128 - v118
		} else {
			v132 = v118
		}
		v134 = v132 >> (uint(int32(1)) % 32)
		v137 = base.I32_extend16_s(v122 + (int32(base.Ui32(v123*int32(7))>>(uint(int32(3))%32)) + v134))
		v139 = v137 >> (uint(int32(31)) % 32)
		v141 = v137 ^ v139 - v139
		if v141 <= v81 {
			v157 = v141
			v158 = int32(0)
		} else {
			v143 = int32(0)
			v144 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76)+480)))
			v146 = *(*int32)(unsafe.Add(mBase, uint32(v76)+512))
			v150 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76)+448)))
			v151 = int32(base.Ui32(v144*v141+v146)>>(uint(int32(17))%32)) * v150
			if v137 < v143 {
				v155 = v143 - v151
			} else {
				v155 = v151
			}
			v157 = v141 - v151
			v158 = v155
		}
		*(*uint16)(unsafe.Add(mBase, uint32(v24)+32)) = uint16(v158)
		v161 = int32(0)
		if v137 < v161 {
			v166 = v161 - v157
		} else {
			v166 = v157
		}
		v168 = v166 >> (uint(int32(1)) % 32)
		v169 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v24)+64)))
		v174 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+301)))
		v177 = base.I32_extend16_s(v169 + (int32(base.Ui32(v134*int32(7))>>(uint(int32(3))%32)) + v174))
		v179 = v177 >> (uint(int32(31)) % 32)
		v181 = v177 ^ v179 - v179
		if v81 < v181 {
			v184 = int32(0)
			v185 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76)+480)))
			v187 = *(*int32)(unsafe.Add(mBase, uint32(v76)+512))
			v191 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76)+448)))
			v192 = int32(base.Ui32(v185*v181+v187)>>(uint(int32(17))%32)) * v191
			if v177 < v184 {
				v196 = v184 - v192
			} else {
				v196 = v192
			}
			v198 = v181 - v192
			v199 = v196
		} else {
			v198 = v181
			v199 = int32(0)
		}
		*(*uint16)(unsafe.Add(mBase, uint32(v24)+64)) = uint16(v199)
		v202 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v24)+96)))
		v203 = int32(0)
		if v177 < v203 {
			v207 = v203 - v198
		} else {
			v207 = v198
		}
		v209 = v207 >> (uint(int32(1)) % 32)
		v216 = base.I32_extend16_s(v202 + (v209 + int32(base.Ui32(v168*int32(7))>>(uint(int32(3))%32))))
		v218 = v216 >> (uint(int32(31)) % 32)
		v220 = v216 ^ v218 - v218
		if v220 <= v81 {
			v236 = v220
			v237 = v161
		} else {
			v222 = int32(0)
			v223 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76)+480)))
			v225 = *(*int32)(unsafe.Add(mBase, uint32(v76)+512))
			v229 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76)+448)))
			v230 = int32(base.Ui32(v223*v220+v225)>>(uint(int32(17))%32)) * v229
			if v216 < v222 {
				v234 = v222 - v230
			} else {
				v234 = v230
			}
			v236 = v220 - v230
			v237 = v234
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l1)+869)) = uint8(v209)
		*(*uint8)(unsafe.Add(mBase, uint32(l1)+868)) = uint8(v168)
		v241 = int32(0)
		if v216 < v241 {
			v246 = v241 - v236
		} else {
			v246 = v236
		}
		v248 = int32(base.Ui32(v246) >> (uint(int32(1)) % 32))
		*(*uint8)(unsafe.Add(mBase, uint32(l1)+870)) = uint8(v248)
		*(*uint16)(unsafe.Add(mBase, uint32(v24)+96)) = uint16(v237)
		v251 = int32(2)
		v253 = v77 + v85<<(uint(v251)%32)
		v256 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+302)))
		v257 = int32(3)
		v259 = int32(*(*int8)(unsafe.Add(mBase, uint32(v253)+2)))
		v265 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v24)+128)))
		v267 = base.I32_extend16_s(int32(base.Ui32(v256<<(uint(v257)%32)+v259*int32(7))>>(uint(v257)%32)) + v265)
		v269 = v267 >> (uint(int32(31)) % 32)
		v271 = v267 ^ v269 - v269
		if v81 < v271 {
			v274 = int32(0)
			v275 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76)+480)))
			v277 = *(*int32)(unsafe.Add(mBase, uint32(v76)+512))
			v281 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76)+448)))
			v282 = int32(base.Ui32(v275*v271+v277)>>(uint(int32(17))%32)) * v281
			if v267 < v274 {
				v286 = v274 - v282
			} else {
				v286 = v282
			}
			v288 = v271 - v282
			v290 = v286
		} else {
			v288 = v271
			v290 = int32(0)
		}
		*(*uint16)(unsafe.Add(mBase, uint32(v24)+128)) = uint16(v290)
		v292 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v24)+160)))
		v293 = int32(*(*int8)(unsafe.Add(mBase, uint32(v253+v251)+1)))
		v298 = int32(0)
		if v267 < v298 {
			v302 = v298 - v288
		} else {
			v302 = v288
		}
		v304 = v302 >> (uint(int32(1)) % 32)
		v307 = base.I32_extend16_s(v292 + (int32(base.Ui32(v293*int32(7))>>(uint(int32(3))%32)) + v304))
		v309 = v307 >> (uint(int32(31)) % 32)
		v311 = v307 ^ v309 - v309
		if v311 <= v81 {
			v327 = v311
			v328 = v241
		} else {
			v313 = int32(0)
			v314 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76)+480)))
			v316 = *(*int32)(unsafe.Add(mBase, uint32(v76)+512))
			v320 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76)+448)))
			v321 = int32(base.Ui32(v314*v311+v316)>>(uint(int32(17))%32)) * v320
			if v307 < v313 {
				v325 = v313 - v321
			} else {
				v325 = v321
			}
			v327 = v311 - v321
			v328 = v325
		}
		*(*uint16)(unsafe.Add(mBase, uint32(v24)+160)) = uint16(v328)
		v331 = int32(0)
		if v307 < v331 {
			v336 = v331 - v327
		} else {
			v336 = v327
		}
		v338 = v336 >> (uint(int32(1)) % 32)
		v339 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v24)+192)))
		v344 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+303)))
		v347 = base.I32_extend16_s(v339 + (int32(base.Ui32(v304*int32(7))>>(uint(int32(3))%32)) + v344))
		v349 = v347 >> (uint(int32(31)) % 32)
		v351 = v347 ^ v349 - v349
		if v81 < v351 {
			v354 = int32(0)
			v355 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76)+480)))
			v357 = *(*int32)(unsafe.Add(mBase, uint32(v76)+512))
			v361 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76)+448)))
			v362 = int32(base.Ui32(v355*v351+v357)>>(uint(int32(17))%32)) * v361
			if v347 < v354 {
				v366 = v354 - v362
			} else {
				v366 = v362
			}
			v368 = v351 - v362
			v369 = v366
		} else {
			v368 = v351
			v369 = int32(0)
		}
		*(*uint16)(unsafe.Add(mBase, uint32(v24)+192)) = uint16(v369)
		v372 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v24)+224)))
		v373 = int32(0)
		if v347 < v373 {
			v377 = v373 - v368
		} else {
			v377 = v368
		}
		v379 = v377 >> (uint(int32(1)) % 32)
		v386 = base.I32_extend16_s(v372 + (v379 + int32(base.Ui32(v338*int32(7))>>(uint(int32(3))%32))))
		v388 = v386 >> (uint(int32(31)) % 32)
		v390 = v386 ^ v388 - v388
		if v390 <= v81 {
			v406 = v390
			v408 = v331
		} else {
			v392 = int32(0)
			v393 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76)+480)))
			v395 = *(*int32)(unsafe.Add(mBase, uint32(v76)+512))
			v399 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76)+448)))
			v400 = int32(base.Ui32(v393*v390+v395)>>(uint(int32(17))%32)) * v399
			if v386 < v392 {
				v404 = v392 - v400
			} else {
				v404 = v400
			}
			v406 = v390 - v400
			v408 = v404
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l1)+872)) = uint8(v379)
		*(*uint8)(unsafe.Add(mBase, uint32(l1)+871)) = uint8(v338)
		v411 = int32(0)
		if v386 < v411 {
			v415 = v411 - v406
		} else {
			v415 = v406
		}
		v417 = int32(base.Ui32(v415) >> (uint(int32(1)) % 32))
		*(*uint8)(unsafe.Add(mBase, uint32(l1)+873)) = uint8(v417)
		*(*uint16)(unsafe.Add(mBase, uint32(v24)+224)) = uint16(v408)
	}
	v433 = v74 + int32(856)
	v434 = m.G46
	v435 = *(*int32)(unsafe.Add(mBase, uint32(v434)))
	v436 = m.T0[v435].(func(*base.Module, int32, int32, int32) int32)(m, v24, l1+int32(584), v433)
	mBase = m.M
	v439 = *(*int32)(unsafe.Add(mBase, uint32(v434)))
	v440 = m.T0[v439].(func(*base.Module, int32, int32, int32) int32)(m, v49, l1+int32(648), v433)
	mBase = m.M
	v443 = *(*int32)(unsafe.Add(mBase, uint32(v434)))
	v444 = m.T0[v443].(func(*base.Module, int32, int32, int32) int32)(m, v57, l1+int32(712), v433)
	mBase = m.M
	v445 = *(*int32)(unsafe.Add(mBase, uint32(v434)))
	v446 = m.G36
	v449 = m.T0[v445].(func(*base.Module, int32, int32, int32) int32)(m, v65, l1+int32(776), v433)
	mBase = m.M
	v450 = int32(1)
	v451 = *(*int32)(unsafe.Add(mBase, uint32(v446)))
	m.T0[v451].(func(*base.Module, int32, int32, int32, int32))(m, v40, v24, l2, v450)
	mBase = m.M
	v456 = *(*int32)(unsafe.Add(mBase, uint32(v446)))
	m.T0[v456].(func(*base.Module, int32, int32, int32, int32))(m, v47, v49, l2+int32(128), v450)
	mBase = m.M
	v461 = *(*int32)(unsafe.Add(mBase, uint32(v446)))
	m.T0[v461].(func(*base.Module, int32, int32, int32, int32))(m, v55, v57, l2+int32(8), v450)
	mBase = m.M
	v466 = *(*int32)(unsafe.Add(mBase, uint32(v446)))
	m.T0[v466].(func(*base.Module, int32, int32, int32, int32))(m, v63, v65, l2+int32(136), v450)
	mBase = m.M
	m.G0 = v24 + int32(256)
	return v449<<(uint(int32(22))%32) | (v444<<(uint(int32(20))%32) | (v440<<(uint(int32(18))%32) | v436<<(uint(int32(16))%32)))
}
func F_Reset(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(0)
	v4 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v4 != 0 {
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = int32(1)
	}
	return int32(1)
}
