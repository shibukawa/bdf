//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_WebPCopyPlane(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	var v12 int32
	_ = v12
	var v15 int32
	_ = v15
	var v17 int32
	_ = v17
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v35 int32
	_ = v35
	var v40 int32
	_ = v40
	var v42 int32
	_ = v42
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v61 int32
	_ = v61
	if l5 < int32(1) {
	} else {
		v12 = l5 & int32(3)
		if v12 != 0 {
			v15 = l0
			v17 = l2
			v21 = v12
			for {
				v23 = F_memcpy(m, v17, v15, l4)
				v24 = v23 + l3
				v25 = v15 + l1
				v27 = v21 + int32(-1)
				if v27 != 0 {
					v15 = v25
					v17 = v24
					v21 = v27
					continue
				} else {
					break
				}
				break
			}
			v28 = v25
			v30 = v24
			v35 = l5 & int32(-4)
		} else {
			v28 = l0
			v30 = l2
			v35 = l5
		}
		if base.Ui32(l5) < base.Ui32(int32(4)) {
		} else {
			v40 = v28
			v42 = v30
			v46 = v35 + int32(-1)
			for {
				v48 = F_memcpy(m, v42, v40, l4)
				v50 = v40 + l1
				v51 = F_memcpy(m, v48+l3, v50, l4)
				v53 = v50 + l1
				v54 = F_memcpy(m, v51+l3, v53, l4)
				v56 = v53 + l1
				v57 = F_memcpy(m, v54+l3, v56, l4)
				v61 = v46 + int32(-4)
				if base.Ui32(v61) < base.Ui32(int32(-2)) {
					v40 = v56 + l1
					v42 = v57 + l3
					v46 = v61
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
func F_WebPEncode(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v14 int32
	_ = v14
	var v21 float32
	_ = v21
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v30 float32
	_ = v30
	var v33 int32
	_ = v33
	var v36 int32
	_ = v36
	var v41 int32
	_ = v41
	var v44 int32
	_ = v44
	var v47 int32
	_ = v47
	var v50 int32
	_ = v50
	var v53 int32
	_ = v53
	var v56 int32
	_ = v56
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v65 int32
	_ = v65
	var v69 int32
	_ = v69
	var v72 int32
	_ = v72
	var v75 int32
	_ = v75
	var v78 int32
	_ = v78
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v85 int32
	_ = v85
	var v86 int32
	_ = v86
	var v89 int32
	_ = v89
	var v92 int32
	_ = v92
	var v95 int32
	_ = v95
	var v98 int32
	_ = v98
	var v101 int32
	_ = v101
	var v104 int32
	_ = v104
	var v107 int32
	_ = v107
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v116 int32
	_ = v116
	var v121 int32
	_ = v121
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v134 int32
	_ = v134
	var v139 int32
	_ = v139
	var v140 int32
	_ = v140
	var v142 int32
	_ = v142
	var v144 int32
	_ = v144
	var v147 int32
	_ = v147
	var v150 int32
	_ = v150
	var v153 int32
	_ = v153
	var v156 int32
	_ = v156
	var v161 int32
	_ = v161
	var v286 int32
	_ = v286
	var v287 int32
	_ = v287
	var v288 int32
	_ = v288
	var v291 int32
	_ = v291
	var v294 int32
	_ = v294
	var v295 int32
	_ = v295
	var v296 int32
	_ = v296
	var v304 int32
	_ = v304
	var v306 int32
	_ = v306
	var v309 int32
	_ = v309
	var v311 int32
	_ = v311
	var v316 int32
	_ = v316
	var v321 int32
	_ = v321
	var v323 int32
	_ = v323
	var v331 float32
	_ = v331
	var v333 float32
	_ = v333
	var v334 float32
	_ = v334
	var v340 float32
	_ = v340
	var v345 int32
	_ = v345
	var v347 int32
	_ = v347
	var v354 int32
	_ = v354
	var v356 int32
	_ = v356
	var v363 int32
	_ = v363
	var v367 int32
	_ = v367
	var v369 int32
	_ = v369
	var v374 int32
	_ = v374
	var v376 int32
	_ = v376
	var v379 int32
	_ = v379
	var v382 int32
	_ = v382
	var v383 int32
	_ = v383
	var v386 int32
	_ = v386
	var v389 int32
	_ = v389
	var v390 int32
	_ = v390
	var v393 int32
	_ = v393
	var v396 int32
	_ = v396
	var v404 int32
	_ = v404
	var v411 int32
	_ = v411
	var v412 int32
	_ = v412
	var v417 int32
	_ = v417
	var v418 int32
	_ = v418
	var v424 int32
	_ = v424
	var v432 int32
	_ = v432
	var v440 int32
	_ = v440
	var v441 int32
	_ = v441
	var v442 int32
	_ = v442
	var v446 int32
	_ = v446
	var v449 int32
	_ = v449
	var v450 int32
	_ = v450
	v3 = int32(0)
	if l1 == v3 {
		v450 = v3
		goto L1
	} else {
		goto L2
	}
L1:
	;
	return v450
L2:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(0)
	if l0 != 0 {
		goto L3
	} else {
		goto L4
	}
L3:
	;
	v14 = int32(0)
	if l0 == v14 {
		v116 = v14
		goto L7
	} else {
		goto L8
	}
L4:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(3)
	return int32(0)
L5:
	;
	if l1 != 0 {
		goto L42
	} else {
		goto L43
	}
L6:
	;
	if v116 != 0 {
		goto L5
	} else {
		goto L38
	}
L7:
	;
	goto L6
L8:
	;
	v21 = *(*float32)(unsafe.Add(mBase, uint32(l0)+4))
	if base.F32_lt(v21, float32(0)) != 0 {
		v116 = v14
		goto L7
	} else {
		goto L9
	}
L9:
	;
	if base.F32_gt(v21, float32(100)) != 0 {
		v116 = v14
		goto L7
	} else {
		goto L10
	}
L10:
	;
	v26 = int32(0)
	v27 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	if v27 < v26 {
		v116 = v26
		goto L7
	} else {
		goto L11
	}
L11:
	;
	v30 = *(*float32)(unsafe.Add(mBase, uint32(l0)+20))
	if base.F32_lt(v30, float32(0)) != 0 {
		v116 = v26
		goto L7
	} else {
		goto L12
	}
L12:
	;
	v33 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	if base.Ui32(int32(6)) < base.Ui32(v33) {
		v116 = v26
		goto L7
	} else {
		goto L13
	}
L13:
	;
	v36 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	if base.Ui32(v36+int32(-5)) < base.Ui32(int32(-4)) {
		v116 = v26
		goto L7
	} else {
		goto L14
	}
L14:
	;
	v41 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
	if base.Ui32(int32(100)) < base.Ui32(v41) {
		v116 = v26
		goto L7
	} else {
		goto L15
	}
L15:
	;
	v44 = *(*int32)(unsafe.Add(mBase, uint32(l0)+32))
	if base.Ui32(int32(100)) < base.Ui32(v44) {
		v116 = v26
		goto L7
	} else {
		goto L16
	}
L16:
	;
	v47 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
	if base.Ui32(int32(7)) < base.Ui32(v47) {
		v116 = v26
		goto L7
	} else {
		goto L17
	}
L17:
	;
	v50 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	if base.Ui32(int32(1)) < base.Ui32(v50) {
		v116 = v26
		goto L7
	} else {
		goto L18
	}
L18:
	;
	v53 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	if base.Ui32(int32(1)) < base.Ui32(v53) {
		v116 = v26
		goto L7
	} else {
		goto L19
	}
L19:
	;
	v56 = *(*int32)(unsafe.Add(mBase, uint32(l0)+60))
	if base.Ui32(v56+int32(-11)) < base.Ui32(int32(-10)) {
		v116 = v26
		goto L7
	} else {
		goto L20
	}
L20:
	;
	v61 = int32(0)
	v62 = *(*int32)(unsafe.Add(mBase, uint32(l0)+108))
	if v62 < v61 {
		v116 = v61
		goto L7
	} else {
		goto L21
	}
L21:
	;
	v65 = *(*int32)(unsafe.Add(mBase, uint32(l0)+112))
	if int32(100) < v65 {
		v116 = v61
		goto L7
	} else {
		goto L22
	}
L22:
	;
	if v65 < v62 {
		v116 = v61
		goto L7
	} else {
		goto L23
	}
L23:
	;
	v69 = *(*int32)(unsafe.Add(mBase, uint32(l0)+64))
	if base.Ui32(int32(1)) < base.Ui32(v69) {
		v116 = v61
		goto L7
	} else {
		goto L24
	}
L24:
	;
	v72 = *(*int32)(unsafe.Add(mBase, uint32(l0)+68))
	if base.Ui32(int32(7)) < base.Ui32(v72) {
		v116 = v61
		goto L7
	} else {
		goto L25
	}
L25:
	;
	v75 = *(*int32)(unsafe.Add(mBase, uint32(l0)+72))
	if base.Ui32(int32(3)) < base.Ui32(v75) {
		v116 = v61
		goto L7
	} else {
		goto L26
	}
L26:
	;
	v78 = *(*int32)(unsafe.Add(mBase, uint32(l0)+76))
	if base.Ui32(int32(100)) < base.Ui32(v78) {
		v116 = v61
		goto L7
	} else {
		goto L27
	}
L27:
	;
	v81 = int32(0)
	v82 = *(*int32)(unsafe.Add(mBase, uint32(l0)+48))
	if v82 < v81 {
		v116 = v81
		goto L7
	} else {
		goto L28
	}
L28:
	;
	v85 = int32(0)
	v86 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if v86 < v85 {
		v116 = v85
		goto L7
	} else {
		goto L29
	}
L29:
	;
	v89 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
	if base.Ui32(int32(100)) < base.Ui32(v89) {
		v116 = v85
		goto L7
	} else {
		goto L30
	}
L30:
	;
	v92 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if base.Ui32(int32(1)) < base.Ui32(v92) {
		v116 = v85
		goto L7
	} else {
		goto L31
	}
L31:
	;
	v95 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
	if base.Ui32(int32(100)) < base.Ui32(v95) {
		v116 = v85
		goto L7
	} else {
		goto L32
	}
L32:
	;
	v98 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	if base.Ui32(int32(3)) < base.Ui32(v98) {
		v116 = v85
		goto L7
	} else {
		goto L33
	}
L33:
	;
	v101 = *(*int32)(unsafe.Add(mBase, uint32(l0)+80))
	if base.Ui32(int32(1)) < base.Ui32(v101) {
		v116 = v85
		goto L7
	} else {
		goto L34
	}
L34:
	;
	v104 = *(*int32)(unsafe.Add(mBase, uint32(l0)+84))
	if base.Ui32(int32(1)) < base.Ui32(v104) {
		v116 = v85
		goto L7
	} else {
		goto L35
	}
L35:
	;
	v107 = *(*int32)(unsafe.Add(mBase, uint32(l0)+88))
	if base.Ui32(int32(1)) < base.Ui32(v107) {
		v116 = v85
		goto L7
	} else {
		goto L36
	}
L36:
	;
	v110 = *(*int32)(unsafe.Add(mBase, uint32(l0)+96))
	if base.Ui32(int32(1)) < base.Ui32(v110) {
		v116 = v85
		goto L7
	} else {
		goto L37
	}
L37:
	;
	v113 = *(*int32)(unsafe.Add(mBase, uint32(l0)+104))
	v116 = base.B2i32(base.Ui32(v113) < base.Ui32(int32(2)))
	goto L7
L38:
	;
	v121 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v121 != 0 {
		v450 = int32(0)
		goto L1
	} else {
		goto L39
	}
L39:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(4)
	return int32(0)
L40:
	;
	v150 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	if int32(_a_F_WebPEncode_0) < v150 {
		goto L50
	} else {
		goto L51
	}
L41:
	;
	if v147 != 0 {
		goto L40
	} else {
		goto L48
	}
L42:
	;
	v130 = int32(5)
	v131 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	if v131 < int32(1) {
		v140 = v130
		goto L45
	} else {
		goto L46
	}
L43:
	;
	v147 = int32(0)
	goto L41
L44:
	;
	v147 = v144
	goto L41
L45:
	;
	v142 = F_WebPEncodingSetError(m, l1, v140)
	mBase = m.M
	v144 = v142
	goto L44
L46:
	;
	v134 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	if v134 < int32(1) {
		v140 = v130
		goto L45
	} else {
		goto L47
	}
L47:
	;
	v139 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
	switch v139 {
	case 0, 4:
		v144 = int32(1)
		goto L44
	default:
		v140 = int32(4)
		goto L45
	}
L48:
	;
	return int32(0)
L49:
	;
	v161 = *(*int32)(unsafe.Add(mBase, uint32(l1)+88))
	if v161 == int32(0) {
		goto L54
	} else {
		goto L55
	}
L50:
	;
	v156 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v156 != 0 {
		v450 = int32(0)
		goto L1
	} else {
		goto L53
	}
L51:
	;
	v153 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	if v153 < int32(_a_F_WebPEncode_1) {
		goto L49
	} else {
		goto L52
	}
L52:
	;
	goto L50
L53:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(5)
	return int32(0)
L54:
	;
	v286 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if v286 != 0 {
		goto L70
	} else {
		goto L71
	}
L55:
	;
	goto L58
L56:
	;
	goto L54
L58:
	;
	base.MemoryFill(m, v161, int32(0), int32(188))
	goto L56
L70:
	;
	v440 = *(*int32)(unsafe.Add(mBase, uint32(l1)+52))
	if v440 != 0 {
		goto L123
	} else {
		goto L124
	}
L71:
	;
	v287 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	if v287 != 0 {
		goto L73
	} else {
		goto L74
	}
L72:
	;
	v374 = *(*int32)(unsafe.Add(mBase, uint32(l0)+96))
	if v374 != 0 {
		goto L98
	} else {
		goto L99
	}
L73:
	;
	v295 = *(*int32)(unsafe.Add(mBase, uint32(l0)+104))
	if v295 != 0 {
		goto L79
	} else {
		goto L80
	}
L74:
	;
	v288 = *(*int32)(unsafe.Add(mBase, uint32(l1)+16))
	if v288 == int32(0) {
		goto L73
	} else {
		goto L75
	}
L75:
	;
	v291 = *(*int32)(unsafe.Add(mBase, uint32(l1)+20))
	if v291 == int32(0) {
		goto L73
	} else {
		goto L76
	}
L76:
	;
	v294 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	if v294 != 0 {
		goto L72
	} else {
		goto L77
	}
L77:
	;
	goto L73
L78:
	;
	if v296&int32(2) == int32(0) {
		v340 = float32(0)
		goto L88
	} else {
		goto L89
	}
L79:
	;
	if l1 != 0 {
		goto L83
	} else {
		goto L84
	}
L80:
	;
	v296 = *(*int32)(unsafe.Add(mBase, uint32(l0)+68))
	if v296&int32(4) == int32(0) {
		goto L78
	} else {
		goto L81
	}
L81:
	;
	goto L79
L82:
	;
	if v323 != 0 {
		goto L72
	} else {
		goto L87
	}
L83:
	;
	v304 = *(*int32)(unsafe.Add(mBase, uint32(l1)+52))
	if v304 != 0 {
		goto L85
	} else {
		goto L86
	}
L84:
	;
	v323 = int32(0)
	goto L82
L85:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = int32(0)
	v309 = int32(2)
	v311 = int32(1)
	v316 = *(*int32)(unsafe.Add(mBase, uint32(l1)+56))
	v321 = F_ImportYUVAFromRGBA(m, v304+v309, v304+v311, v304, v304+int32(3), int32(4), v316<<(uint(v309)%32), float32(0), v311, l1)
	mBase = m.M
	v323 = v321
	goto L82
L86:
	;
	v306 = F_WebPEncodingSetError(m, l1, int32(3))
	mBase = m.M
	v323 = v306
	goto L82
L87:
	;
	return int32(0)
L88:
	;
	if l1 != 0 {
		goto L91
	} else {
		goto L92
	}
L89:
	;
	v331 = *(*float32)(unsafe.Add(mBase, uint32(l0)+4))
	v333 = base.F32_div(v331, float32(100))
	v334 = base.F32_mul(v333, v333)
	v340 = base.F32_add(base.F32_mul(base.F32_mul(v334, float32(-0.5)), v334), float32(1))
	goto L88
L90:
	;
	if v369 == int32(0) {
		v450 = int32(0)
		goto L1
	} else {
		goto L97
	}
L91:
	;
	v345 = *(*int32)(unsafe.Add(mBase, uint32(l1)+52))
	if v345 != 0 {
		goto L93
	} else {
		goto L94
	}
L92:
	;
	v369 = int32(0)
	goto L90
L93:
	;
	goto L95
L94:
	;
	v347 = F_WebPEncodingSetError(m, l1, int32(3))
	mBase = m.M
	v369 = v347
	goto L90
L95:
	;
	v354 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v354
	v356 = int32(2)
	v363 = *(*int32)(unsafe.Add(mBase, uint32(l1)+56))
	v367 = F_ImportYUVAFromRGBA(m, v345+v356, v345+int32(1), v345, v345+int32(3), int32(4), v363<<(uint(v356)%32), v340, v354, l1)
	mBase = m.M
	v369 = v367
	goto L90
L97:
	;
	goto L72
L98:
	;
	v376 = F_InitVP8Encoder(m, l0, l1)
	mBase = m.M
	if v376 != 0 {
		goto L100
	} else {
		goto L101
	}
L99:
	;
	F_WebPCleanupTransparentArea(m, l1)
	mBase = m.M
	goto L98
L100:
	;
	v379 = F_VP8EncAnalyze(m, v376)
	mBase = m.M
	if v379 == int32(0) {
		goto L105
	} else {
		goto L106
	}
L101:
	;
	return int32(0)
L102:
	;
	v432 = F_VP8EncDeleteAlpha(m, v376)
	mBase = m.M
	F_VP8TBufferClear(m, v376+int32(344))
	mBase = m.M
	F_WebPSafeFree(m, v376)
	mBase = m.M
	goto L122
L103:
	;
	F_VP8BitWriterWipeOut(m, v376+int32(56))
	mBase = m.M
	v404 = *(*int32)(unsafe.Add(mBase, uint32(v376)+52))
	if v404 < int32(1) {
		goto L116
	} else {
		goto L117
	}
L104:
	;
	v396 = F_VP8EncWrite(m, v376)
	mBase = m.M
	F_StoreStats(m, v376)
	mBase = m.M
	if v396 != 0 {
		goto L102
	} else {
		goto L114
	}
L105:
	;
	F_StoreStats(m, v376)
	mBase = m.M
	goto L103
L106:
	;
	v382 = F_VP8EncStartAlpha(m, v376)
	mBase = m.M
	v383 = *(*int32)(unsafe.Add(mBase, uint32(v376)+uint32(_c_F_WebPEncode[0])))
	if v383 != 0 {
		goto L108
	} else {
		goto L109
	}
L107:
	;
	if v390 == int32(0) {
		goto L105
	} else {
		goto L112
	}
L108:
	;
	if v382 == int32(0) {
		goto L105
	} else {
		goto L111
	}
L109:
	;
	if v382 == int32(0) {
		goto L105
	} else {
		goto L110
	}
L110:
	;
	v386 = F_VP8EncLoop(m, v376)
	mBase = m.M
	v390 = v386
	goto L107
L111:
	;
	v389 = F_VP8EncTokenLoop(m, v376)
	mBase = m.M
	v390 = v389
	goto L107
L112:
	;
	v393 = F_VP8EncFinishAlpha(m, v376)
	mBase = m.M
	if v393 != 0 {
		goto L104
	} else {
		goto L113
	}
L113:
	;
	goto L105
L114:
	;
	goto L103
L115:
	;
	v424 = F_VP8EncDeleteAlpha(m, v376)
	mBase = m.M
	F_VP8TBufferClear(m, v376+int32(344))
	mBase = m.M
	F_WebPSafeFree(m, v376)
	mBase = m.M
	goto L121
L116:
	;
	goto L115
L117:
	;
	v411 = v376 + int32(88)
	v412 = int32(0)
	goto L118
L118:
	;
	F_VP8BitWriterWipeOut(m, v411)
	mBase = m.M
	v417 = v412 + int32(1)
	v418 = *(*int32)(unsafe.Add(mBase, uint32(v376)+52))
	if v417 < v418 {
		v411 = v411 + int32(32)
		v412 = v417
		goto L118
	} else {
		goto L120
	}
L119:
	;
	goto L116
L120:
	;
	goto L119
L121:
	;
	return int32(0)
L122:
	;
	return v432 & int32(1)
L123:
	;
	v446 = *(*int32)(unsafe.Add(mBase, uint32(l0)+96))
	if v446 != 0 {
		goto L126
	} else {
		goto L127
	}
L124:
	;
	v441 = int32(0)
	v442 = F_WebPPictureYUVAToARGB(m, l1)
	mBase = m.M
	if v442 == v441 {
		v450 = v441
		goto L1
	} else {
		goto L125
	}
L125:
	;
	goto L123
L126:
	;
	v449 = F_VP8LEncodeImage(m, l0, l1)
	mBase = m.M
	v450 = v449
	goto L1
L127:
	;
	F_WebPReplaceTransparentPixels(m, l1, int32(0))
	mBase = m.M
	goto L126
}
func F_WebPEncodingSetError(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
	if v3 != 0 {
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = l1
	}
	return int32(0)
}
func F_WebPGetLinePairConverter(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v7 int32
	_ = v7
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v19 int32
	_ = v19
	var v22 int32
	_ = v22
	var v27 int32
	_ = v27
	var v30 int32
	_ = v30
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v56 int32
	_ = v56
	var v59 int32
	_ = v59
	var v64 int32
	_ = v64
	var v77 int32
	_ = v77
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v96 int32
	_ = v96
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v106 int32
	_ = v106
	var v111 int32
	_ = v111
	var v113 int32
	_ = v113
	v7 = m.G1
	v10 = *(*int32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_WebPGetLinePairConverter[0])))
	v11 = m.G13
	v12 = *(*int32)(unsafe.Add(mBase, uint32(v11)))
	if v10 == v12 {
	} else {
		v14 = m.G2
		v15 = m.G1
		v19 = v14 + int32(50)
		*(*int32)(unsafe.Add(mBase, uint32(v15)+uint32(_c_F_WebPGetLinePairConverter[1]))) = v19
		v22 = v14 + int32(51)
		*(*int32)(unsafe.Add(mBase, uint32(v15)+uint32(_c_F_WebPGetLinePairConverter[2]))) = v22
		*(*int32)(unsafe.Add(mBase, uint32(v15)+uint32(_c_F_WebPGetLinePairConverter[3]))) = v19
		*(*int32)(unsafe.Add(mBase, uint32(v15)+uint32(_c_F_WebPGetLinePairConverter[4]))) = v22
		v27 = v14 + int32(52)
		*(*int32)(unsafe.Add(mBase, uint32(v15)+uint32(_c_F_WebPGetLinePairConverter[5]))) = v27
		v30 = v14 + int32(53)
		*(*int32)(unsafe.Add(mBase, uint32(v15)+uint32(_c_F_WebPGetLinePairConverter[6]))) = v30
		*(*int32)(unsafe.Add(mBase, uint32(v15)+uint32(_c_F_WebPGetLinePairConverter[7]))) = v14 + int32(54)
		*(*int32)(unsafe.Add(mBase, uint32(v15)+uint32(_c_F_WebPGetLinePairConverter[8]))) = v27
		*(*int32)(unsafe.Add(mBase, uint32(v15)+uint32(_c_F_WebPGetLinePairConverter[9]))) = v30
		*(*int32)(unsafe.Add(mBase, uint32(v15)+uint32(_c_F_WebPGetLinePairConverter[10]))) = v14 + int32(55)
		*(*int32)(unsafe.Add(mBase, uint32(v15)+uint32(_c_F_WebPGetLinePairConverter[11]))) = v14 + int32(56)
		if v12 == int32(0) {
		} else {
			v45 = int32(0)
			v46 = m.T0[v12].(func(*base.Module, int32) int32)(m, v45)
			mBase = m.M
			if v46 == v45 {
			} else {
				v53 = m.G2
				v54 = m.G19
				v56 = v53 + int32(41)
				*(*int32)(unsafe.Add(mBase, uint32(v54)+32)) = v56
				v59 = v53 + int32(42)
				*(*int32)(unsafe.Add(mBase, uint32(v54)+28)) = v59
				*(*int32)(unsafe.Add(mBase, uint32(v54)+12)) = v56
				*(*int32)(unsafe.Add(mBase, uint32(v54)+4)) = v59
				v64 = v53 + int32(43)
				*(*int32)(unsafe.Add(mBase, uint32(v54)+36)) = v64
				*(*int32)(unsafe.Add(mBase, uint32(v54)+16)) = v64
				*(*int32)(unsafe.Add(mBase, uint32(v54)+8)) = v53 + int32(44)
				*(*int32)(unsafe.Add(mBase, uint32(v54))) = v53 + int32(45)
				*(*int32)(unsafe.Add(mBase, uint32(v54)+24)) = v53 + int32(46)
				v77 = v53 + int32(47)
				*(*int32)(unsafe.Add(mBase, uint32(v54)+40)) = v77
				*(*int32)(unsafe.Add(mBase, uint32(v54)+20)) = v77
			}
			v81 = m.G13
			v82 = *(*int32)(unsafe.Add(mBase, uint32(v81)))
			v83 = m.T0[v82].(func(*base.Module, int32) int32)(m, int32(3))
			mBase = m.M
			if v83 == int32(0) {
			} else {
				v88 = m.G2
				v89 = m.G19
				*(*int32)(unsafe.Add(mBase, uint32(v89)+8)) = v88 + int32(48)
				*(*int32)(unsafe.Add(mBase, uint32(v89))) = v88 + int32(49)
			}
		}
		v96 = m.G1
		v99 = m.G13
		v100 = *(*int32)(unsafe.Add(mBase, uint32(v99)))
		*(*int32)(unsafe.Add(mBase, uint32(v96)+uint32(_c_F_WebPGetLinePairConverter[0]))) = v100
	}
	v106 = m.G1
	if l0 != 0 {
		v111 = int32(12)
	} else {
		v111 = int32(16)
	}
	v113 = *(*int32)(unsafe.Add(mBase, uint32(v106+int32(_a_F_WebPGetLinePairConverter_0)+v111)))
	return v113
}
func F_WebPGetWorkerInterface(m *base.Module) int32 {
	var v1 int32
	_ = v1
	v1 = m.G1
	return v1 + int32(_a_F_WebPGetWorkerInterface_0)
}
func F_WebPInitAlphaProcessing(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v75 int32
	_ = v75
	var v76 int32
	_ = v76
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v85 int32
	_ = v85
	var v89 int32
	_ = v89
	var v93 int32
	_ = v93
	var v97 int32
	_ = v97
	var v101 int32
	_ = v101
	var v105 int32
	_ = v105
	var v109 int32
	_ = v109
	var v113 int32
	_ = v113
	var v117 int32
	_ = v117
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v124 int32
	_ = v124
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v133 int32
	_ = v133
	var v136 int32
	_ = v136
	var v137 int32
	_ = v137
	v4 = m.G1
	v7 = *(*int32)(unsafe.Add(mBase, uint32(v4)+uint32(_c_F_WebPInitAlphaProcessing[0])))
	v8 = m.G13
	v9 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
	if v7 == v9 {
	} else {
		v11 = m.G2
		v12 = m.G1
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[1]))) = v11 + int32(14)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[2]))) = v11 + int32(15)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[3]))) = v11 + int32(16)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[4]))) = v11 + int32(17)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[5]))) = v11 + int32(18)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[6]))) = v11 + int32(19)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[7]))) = v11 + int32(20)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[8]))) = v11 + int32(21)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[9]))) = v11 + int32(22)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[10]))) = v11 + int32(23)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[11]))) = v11 + int32(24)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[12]))) = v11 + int32(25)
		if v9 == int32(0) {
		} else {
			v75 = int32(0)
			v76 = m.T0[v9].(func(*base.Module, int32) int32)(m, v75)
			mBase = m.M
			if v76 == v75 {
			} else {
				v80 = m.G2
				v81 = m.G3
				*(*int32)(unsafe.Add(mBase, uint32(v81))) = v80 + int32(3)
				v85 = m.G4
				*(*int32)(unsafe.Add(mBase, uint32(v85))) = v80 + int32(4)
				v89 = m.G5
				*(*int32)(unsafe.Add(mBase, uint32(v89))) = v80 + int32(5)
				v93 = m.G6
				*(*int32)(unsafe.Add(mBase, uint32(v93))) = v80 + int32(6)
				v97 = m.G7
				*(*int32)(unsafe.Add(mBase, uint32(v97))) = v80 + int32(7)
				v101 = m.G8
				*(*int32)(unsafe.Add(mBase, uint32(v101))) = v80 + int32(8)
				v105 = m.G9
				*(*int32)(unsafe.Add(mBase, uint32(v105))) = v80 + int32(9)
				v109 = m.G10
				*(*int32)(unsafe.Add(mBase, uint32(v109))) = v80 + int32(10)
				v113 = m.G11
				*(*int32)(unsafe.Add(mBase, uint32(v113))) = v80 + int32(11)
				v117 = m.G12
				*(*int32)(unsafe.Add(mBase, uint32(v117))) = v80 + int32(12)
				v122 = m.G13
				v123 = *(*int32)(unsafe.Add(mBase, uint32(v122)))
				v124 = m.T0[v123].(func(*base.Module, int32) int32)(m, int32(3))
				mBase = m.M
				if v124 == int32(0) {
				} else {
					v128 = m.G2
					v129 = m.G8
					*(*int32)(unsafe.Add(mBase, uint32(v129))) = v128 + int32(13)
				}
			}
		}
		v133 = m.G1
		v136 = m.G13
		v137 = *(*int32)(unsafe.Add(mBase, uint32(v136)))
		*(*int32)(unsafe.Add(mBase, uint32(v133)+uint32(_c_F_WebPInitAlphaProcessing[0]))) = v137
	}
	return
}
func F_WebPInitAlphaProcessingSSE2(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v3 int32
	_ = v3
	var v7 int32
	_ = v7
	var v11 int32
	_ = v11
	var v15 int32
	_ = v15
	var v19 int32
	_ = v19
	var v23 int32
	_ = v23
	var v27 int32
	_ = v27
	var v31 int32
	_ = v31
	var v35 int32
	_ = v35
	var v39 int32
	_ = v39
	v2 = m.G2
	v3 = m.G3
	*(*int32)(unsafe.Add(mBase, uint32(v3))) = v2 + int32(3)
	v7 = m.G4
	*(*int32)(unsafe.Add(mBase, uint32(v7))) = v2 + int32(4)
	v11 = m.G5
	*(*int32)(unsafe.Add(mBase, uint32(v11))) = v2 + int32(5)
	v15 = m.G6
	*(*int32)(unsafe.Add(mBase, uint32(v15))) = v2 + int32(6)
	v19 = m.G7
	*(*int32)(unsafe.Add(mBase, uint32(v19))) = v2 + int32(7)
	v23 = m.G8
	*(*int32)(unsafe.Add(mBase, uint32(v23))) = v2 + int32(8)
	v27 = m.G9
	*(*int32)(unsafe.Add(mBase, uint32(v27))) = v2 + int32(9)
	v31 = m.G10
	*(*int32)(unsafe.Add(mBase, uint32(v31))) = v2 + int32(10)
	v35 = m.G11
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2 + int32(11)
	v39 = m.G12
	*(*int32)(unsafe.Add(mBase, uint32(v39))) = v2 + int32(12)
	return
}
func F_WebPInitAlphaProcessingSSE41(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v3 int32
	_ = v3
	v2 = m.G2
	v3 = m.G8
	*(*int32)(unsafe.Add(mBase, uint32(v3))) = v2 + int32(13)
	return
}
func F_WebPInitConvertARGBToYUV(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v45 int32
	_ = v45
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
	var v67 int32
	_ = v67
	var v68 int32
	_ = v68
	var v69 int32
	_ = v69
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v78 int32
	_ = v78
	var v82 int32
	_ = v82
	var v86 int32
	_ = v86
	var v90 int32
	_ = v90
	var v94 int32
	_ = v94
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	v4 = m.G1
	v7 = *(*int32)(unsafe.Add(mBase, uint32(v4)+uint32(_c_F_WebPInitConvertARGBToYUV[0])))
	v8 = m.G13
	v9 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
	if v7 == v9 {
	} else {
		v11 = m.G2
		v12 = m.G1
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitConvertARGBToYUV[1]))) = v11 + int32(31)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitConvertARGBToYUV[2]))) = v11 + int32(32)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitConvertARGBToYUV[3]))) = v11 + int32(33)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitConvertARGBToYUV[4]))) = v11 + int32(34)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitConvertARGBToYUV[5]))) = v11 + int32(35)
		if v9 == int32(0) {
		} else {
			v40 = int32(0)
			v41 = m.T0[v9].(func(*base.Module, int32) int32)(m, v40)
			mBase = m.M
			if v41 == v40 {
			} else {
				v45 = m.G2
				v46 = m.G14
				*(*int32)(unsafe.Add(mBase, uint32(v46))) = v45 + int32(36)
				v50 = m.G15
				*(*int32)(unsafe.Add(mBase, uint32(v50))) = v45 + int32(37)
				v54 = m.G16
				*(*int32)(unsafe.Add(mBase, uint32(v54))) = v45 + int32(38)
				v58 = m.G17
				*(*int32)(unsafe.Add(mBase, uint32(v58))) = v45 + int32(39)
				v62 = m.G18
				*(*int32)(unsafe.Add(mBase, uint32(v62))) = v45 + int32(40)
			}
			v67 = m.G13
			v68 = *(*int32)(unsafe.Add(mBase, uint32(v67)))
			v69 = m.T0[v68].(func(*base.Module, int32) int32)(m, int32(3))
			mBase = m.M
			if v69 == int32(0) {
			} else {
				v73 = m.G2
				v74 = m.G14
				*(*int32)(unsafe.Add(mBase, uint32(v74))) = v73 + int32(26)
				v78 = m.G15
				*(*int32)(unsafe.Add(mBase, uint32(v78))) = v73 + int32(27)
				v82 = m.G16
				*(*int32)(unsafe.Add(mBase, uint32(v82))) = v73 + int32(28)
				v86 = m.G17
				*(*int32)(unsafe.Add(mBase, uint32(v86))) = v73 + int32(29)
				v90 = m.G18
				*(*int32)(unsafe.Add(mBase, uint32(v90))) = v73 + int32(30)
			}
		}
		v94 = m.G1
		v97 = m.G13
		v98 = *(*int32)(unsafe.Add(mBase, uint32(v97)))
		*(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_WebPInitConvertARGBToYUV[0]))) = v98
	}
	return
}
func F_WebPInitConvertARGBToYUVSSE2(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v3 int32
	_ = v3
	var v7 int32
	_ = v7
	var v11 int32
	_ = v11
	var v15 int32
	_ = v15
	var v19 int32
	_ = v19
	v2 = m.G2
	v3 = m.G14
	*(*int32)(unsafe.Add(mBase, uint32(v3))) = v2 + int32(36)
	v7 = m.G15
	*(*int32)(unsafe.Add(mBase, uint32(v7))) = v2 + int32(37)
	v11 = m.G16
	*(*int32)(unsafe.Add(mBase, uint32(v11))) = v2 + int32(38)
	v15 = m.G17
	*(*int32)(unsafe.Add(mBase, uint32(v15))) = v2 + int32(39)
	v19 = m.G18
	*(*int32)(unsafe.Add(mBase, uint32(v19))) = v2 + int32(40)
	return
}
func F_WebPInitConvertARGBToYUVSSE41(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v3 int32
	_ = v3
	var v7 int32
	_ = v7
	var v11 int32
	_ = v11
	var v15 int32
	_ = v15
	var v19 int32
	_ = v19
	v2 = m.G2
	v3 = m.G14
	*(*int32)(unsafe.Add(mBase, uint32(v3))) = v2 + int32(26)
	v7 = m.G15
	*(*int32)(unsafe.Add(mBase, uint32(v7))) = v2 + int32(27)
	v11 = m.G16
	*(*int32)(unsafe.Add(mBase, uint32(v11))) = v2 + int32(28)
	v15 = m.G17
	*(*int32)(unsafe.Add(mBase, uint32(v15))) = v2 + int32(29)
	v19 = m.G18
	*(*int32)(unsafe.Add(mBase, uint32(v19))) = v2 + int32(30)
	return
}
func F_WebPInitUpsamplersSSE2(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v11 int32
	_ = v11
	var v16 int32
	_ = v16
	var v29 int32
	_ = v29
	v5 = m.G2
	v6 = m.G19
	v8 = v5 + int32(41)
	*(*int32)(unsafe.Add(mBase, uint32(v6)+32)) = v8
	v11 = v5 + int32(42)
	*(*int32)(unsafe.Add(mBase, uint32(v6)+28)) = v11
	*(*int32)(unsafe.Add(mBase, uint32(v6)+12)) = v8
	*(*int32)(unsafe.Add(mBase, uint32(v6)+4)) = v11
	v16 = v5 + int32(43)
	*(*int32)(unsafe.Add(mBase, uint32(v6)+36)) = v16
	*(*int32)(unsafe.Add(mBase, uint32(v6)+16)) = v16
	*(*int32)(unsafe.Add(mBase, uint32(v6)+8)) = v5 + int32(44)
	*(*int32)(unsafe.Add(mBase, uint32(v6))) = v5 + int32(45)
	*(*int32)(unsafe.Add(mBase, uint32(v6)+24)) = v5 + int32(46)
	v29 = v5 + int32(47)
	*(*int32)(unsafe.Add(mBase, uint32(v6)+40)) = v29
	*(*int32)(unsafe.Add(mBase, uint32(v6)+20)) = v29
	return
}
func F_WebPInitUpsamplersSSE41(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v4 int32
	_ = v4
	v3 = m.G2
	v4 = m.G19
	*(*int32)(unsafe.Add(mBase, uint32(v4)+8)) = v3 + int32(48)
	*(*int32)(unsafe.Add(mBase, uint32(v4))) = v3 + int32(49)
	return
}
func F_WebPMemoryWrite(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v8 int32
	_ = v8
	var v11 int64
	_ = v11
	var v13 int64
	_ = v13
	var v14 int64
	_ = v14
	var v17 int64
	_ = v17
	var v19 int64
	_ = v19
	var v20 int64
	_ = v20
	var v23 int64
	_ = v23
	var v24 int32
	_ = v24
	var v31 int64
	_ = v31
	var v32 int32
	_ = v32
	var v33 int64
	_ = v33
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v49 int32
	_ = v49
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v67 int32
	_ = v67
	var v68 int32
	_ = v68
	v8 = *(*int32)(unsafe.Add(mBase, uint32(l2)+76))
	if v8 == int32(0) {
		return int32(1)
	} else {
		v11 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v8)+4)))
		v13 = v11 + base.I64_extend_i32_u(l1)
		v14 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v8)+8)))
		if base.Ui64(v13) <= base.Ui64(v14) {
			if l1 == int32(0) {
			} else {
				v64 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
				v65 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
				v67 = F_memcpy(m, v64+v65, l0, l1)
				mBase = m.M
				v68 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
				*(*int32)(unsafe.Add(mBase, uint32(v8)+4)) = v68 + l1
			}
			return int32(1)
		} else {
			v17 = v14 << (uint(int64(1)) % 64)
			if base.Ui64(v13) < base.Ui64(v17) {
				v19 = v17
			} else {
				v19 = v13
			}
			v20 = int64(8192)
			if base.Ui64(v20) < base.Ui64(v19) {
				v23 = v19
			} else {
				v23 = v20
			}
			v24 = int32(1)
			if v23 == int64(0) {
				v43 = F_malloc(m, base.I32_wrap_i64(v23)*v24)
				mBase = m.M
				v45 = v43
			} else {
				v31 = base.I64_div_u_s(int64(2147418112), v23)
				v32 = int32(0)
				v33 = base.I64_extend_i32_u(v24)
				if base.Ui64(int64(4294967295)) < base.Ui64(v33*v23) {
					v45 = v32
				} else {
					if base.Ui64(v31) < base.Ui64(v33) {
						v45 = v32
					} else {
						v43 = F_malloc(m, base.I32_wrap_i64(v23)*v24)
						mBase = m.M
						v45 = v43
					}
				}
			}
			if v45 != 0 {
				v49 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
				if v49 == int32(0) {
				} else {
					v52 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
					v53 = F_memcpy(m, v45, v52, v49)
					mBase = m.M
				}
				v54 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
				F_free(m, v54)
				mBase = m.M
				*(*uint32)(unsafe.Add(mBase, uint32(v8)+8)) = uint32(v23)
				*(*int32)(unsafe.Add(mBase, uint32(v8))) = v45
				if l1 == int32(0) {
				} else {
					v64 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
					v65 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
					v67 = F_memcpy(m, v64+v65, l0, l1)
					mBase = m.M
					v68 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
					*(*int32)(unsafe.Add(mBase, uint32(v8)+4)) = v68 + l1
				}
				return int32(1)
			} else {
				return int32(0)
			}
		}
	}
}
func F_WebPMemoryWriterClear(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	if l0 == int32(0) {
	} else {
		v4 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		F_free(m, v4)
		mBase = m.M
		*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = int32(0)
		*(*int64)(unsafe.Add(mBase, uint32(l0))) = int64(0)
	}
	return
}
func F_WebPMemoryWriterInit(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = int32(0)
	*(*int64)(unsafe.Add(mBase, uint32(l0))) = int64(0)
	return
}
func F_WebPPictureARGBToYUVADithered(m *base.Module, l0 int32, l1 int32, l2 float32) int32 {
	mBase := m.M
	_ = mBase
	var v7 int32
	_ = v7
	var v9 int32
	_ = v9
	var v18 int32
	_ = v18
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	var v31 int32
	_ = v31
	var v35 int32
	_ = v35
	if l0 != 0 {
		v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
		if v7 != 0 {
			if l1&int32(3) == int32(0) {
				v22 = int32(0)
				*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v22
				v24 = int32(2)
				v31 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
				v35 = F_ImportYUVAFromRGBA(m, v7+v24, v7+int32(1), v7, v7+int32(3), int32(4), v31<<(uint(v24)%32), l2, v22, l0)
				mBase = m.M
				return v35
			} else {
				v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
				if v18 != 0 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(4)
				}
				return int32(0)
			}
		} else {
			v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
			if v9 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(3)
			}
			return int32(0)
		}
	} else {
		return int32(0)
	}
}
func F_WebPPictureAllocARGB(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v6 int32
	_ = v6
	var v7 int32
	_ = v7
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v21 int32
	_ = v21
	var v23 int64
	_ = v23
	var v31 int64
	_ = v31
	var v32 int32
	_ = v32
	var v39 int64
	_ = v39
	var v40 int32
	_ = v40
	var v41 int64
	_ = v41
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v56 int32
	_ = v56
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v6 = int32(5)
	v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	if v7 < int32(1) {
		v14 = v6
		v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
		if v15 != 0 {
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v14
		}
		return int32(0)
	} else {
		if v5 < int32(1) {
			v14 = v6
			v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
			if v15 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v14
			}
			return int32(0)
		} else {
			v13 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
			switch v13 {
			case 0, 4:
				v21 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
				F_free(m, v21)
				mBase = m.M
				v23 = int64(0)
				*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = v23
				*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
				v31 = base.I64_extend_i32_s(v5)*base.I64_extend_i32_s(v7) + int64(31)
				v32 = int32(4)
				if v31 == v23 {
					v51 = F_malloc(m, base.I32_wrap_i64(v31)*v32)
					mBase = m.M
					v53 = v51
				} else {
					v39 = base.I64_div_u_s(int64(2147418112), v31)
					v40 = int32(0)
					v41 = base.I64_extend_i32_u(v32)
					if base.Ui64(int64(4294967295)) < base.Ui64(v41*v31) {
						v53 = v40
					} else {
						if base.Ui64(v39) < base.Ui64(v41) {
							v53 = v40
						} else {
							v51 = F_malloc(m, base.I32_wrap_i64(v31)*v32)
							mBase = m.M
							v53 = v51
						}
					}
				}
				if v53 != 0 {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v7
					*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v53
					*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v53 + int32(31)) & int32(-32)
					return int32(1)
				} else {
					v56 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
					if v56 != 0 {
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(1)
					}
					return int32(0)
				}
			default:
				v14 = int32(4)
				v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
				if v15 != 0 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v14
				}
				return int32(0)
			}
		}
	}
}
func F_WebPPictureHasTransparency(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v8 int32
	_ = v8
	var v11 int32
	_ = v11
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v26 int32
	_ = v26
	var v29 int32
	_ = v29
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v40 int32
	_ = v40
	var v43 int32
	_ = v43
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v54 int32
	_ = v54
	var v57 int32
	_ = v57
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v68 int32
	_ = v68
	if l0 == int32(0) {
		goto L1
	} else {
		goto L2
	}
L1:
	;
	return int32(0)
L2:
	;
	v8 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if v8 == int32(0) {
		goto L3
	} else {
		goto L4
	}
L3:
	;
	v43 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
	if v43 == int32(0) {
		goto L1
	} else {
		goto L12
	}
L4:
	;
	v11 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if v11 == int32(0) {
		goto L1
	} else {
		goto L5
	}
L5:
	;
	v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
	v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPInitAlphaProcessing(m)
	mBase = m.M
	if v15 < int32(1) {
		goto L1
	} else {
		goto L6
	}
L6:
	;
	v26 = v15 + int32(1)
	v29 = v11 + int32(3)
	goto L7
L7:
	;
	v31 = m.G11
	v32 = *(*int32)(unsafe.Add(mBase, uint32(v31)))
	v33 = m.T0[v32].(func(*base.Module, int32, int32) int32)(m, v29, v16)
	mBase = m.M
	if v33 == int32(0) {
		goto L9
	} else {
		goto L10
	}
L9:
	;
	v40 = v26 + int32(-1)
	if int32(1) < v40 {
		v26 = v40
		v29 = v29 + v14<<(uint(int32(2))%32)
		goto L7
	} else {
		goto L11
	}
L10:
	;
	return int32(1)
L11:
	;
	goto L1
L12:
	;
	v46 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v47 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v48 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	F_WebPInitAlphaProcessing(m)
	mBase = m.M
	if v48 < int32(1) {
		goto L1
	} else {
		goto L13
	}
L13:
	;
	v54 = v48 + int32(1)
	v57 = v43
	goto L14
L14:
	;
	v59 = m.G10
	v60 = *(*int32)(unsafe.Add(mBase, uint32(v59)))
	v61 = m.T0[v60].(func(*base.Module, int32, int32) int32)(m, v57, v47)
	mBase = m.M
	if v61 == int32(0) {
		goto L16
	} else {
		goto L17
	}
L15:
	;
	goto L1
L16:
	;
	v68 = v54 + int32(-1)
	if int32(1) < v68 {
		v54 = v68
		v57 = v57 + v46
		goto L14
	} else {
		goto L18
	}
L17:
	;
	return int32(1)
L18:
	;
	goto L15
}
func F_WebPPictureInitInternal(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v134 int32
	_ = v134
	var v139 int32
	_ = v139
	var v143 int32
	_ = v143
	if l1&int32(-256) != int32(512) {
		v143 = int32(0)
	} else {
		v9 = int32(1)
		if l0 == int32(0) {
			v143 = v9
		} else {
			base.MemoryFill(m, l0, int32(0), int32(172))
			v134 = m.G2
			*(*int32)(unsafe.Add(mBase, uint32(l0)+72)) = v134 + int32(1)
			v139 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
			if v139 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(0)
			}
			v143 = v9
		}
	}
	return v143
}
func F_WebPPictureSharpARGBToYUVA(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v7 int32
	_ = v7
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v20 int32
	_ = v20
	var v25 int32
	_ = v25
	if l0 != 0 {
		v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
		if v5 != 0 {
			*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = int32(0)
			v13 = int32(2)
			v15 = int32(1)
			v20 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
			v25 = F_ImportYUVAFromRGBA(m, v5+v13, v5+v15, v5, v5+int32(3), int32(4), v20<<(uint(v13)%32), float32(0), v15, l0)
			mBase = m.M
			return v25
		} else {
			v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
			if v7 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(3)
			}
			return int32(0)
		}
	} else {
		return int32(0)
	}
}
func F_WebPReplaceTransparentPixels(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v18 int32
	_ = v18
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v27 int32
	_ = v27
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v43 int32
	_ = v43
	var v44 int32
	_ = v44
	var v45 int32
	_ = v45
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
	var v54 int32
	_ = v54
	var v59 int32
	_ = v59
	if l0 == int32(0) {
	} else {
		v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		if v9 == int32(0) {
		} else {
			v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
			v13 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
			F_WebPInitAlphaProcessing(m)
			mBase = m.M
			if v12 < int32(1) {
			} else {
				v18 = l1 & int32(16777215)
				if v12&int32(1) != 0 {
					v21 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					v22 = m.G12
					v23 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
					m.T0[v23].(func(*base.Module, int32, int32, int32))(m, v13, v21, v18)
					mBase = m.M
					v27 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
					v31 = v13 + v27<<(uint(int32(2))%32)
					v32 = v12 + int32(-1)
				} else {
					v31 = v13
					v32 = v12
				}
				if v12 == int32(1) {
				} else {
					v39 = v32 + int32(-1)
					v40 = v31
					for {
						v43 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
						v44 = m.G12
						v45 = *(*int32)(unsafe.Add(mBase, uint32(v44)))
						m.T0[v45].(func(*base.Module, int32, int32, int32))(m, v40, v43, v18)
						mBase = m.M
						v47 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
						v48 = int32(2)
						v50 = v40 + v47<<(uint(v48)%32)
						v51 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
						v52 = *(*int32)(unsafe.Add(mBase, uint32(v44)))
						m.T0[v52].(func(*base.Module, int32, int32, int32))(m, v50, v51, v18)
						mBase = m.M
						v54 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
						v59 = v39 + int32(-2)
						if base.Ui32(v59) < base.Ui32(v39) {
							v39 = v59
							v40 = v50 + v54<<(uint(v48)%32)
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
func F_WebPReportProgress(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v8 int32
	_ = v8
	var v11 int32
	_ = v11
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v20 int32
	_ = v20
	v5 = int32(1)
	if l2 == int32(0) {
		v20 = v5
	} else {
		v8 = *(*int32)(unsafe.Add(mBase, uint32(l2)))
		if v8 == l1 {
			v20 = v5
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l2))) = l1
			v11 = *(*int32)(unsafe.Add(mBase, uint32(l0)+96))
			if v11 == int32(0) {
				v20 = v5
			} else {
				v14 = m.T0[v11].(func(*base.Module, int32, int32) int32)(m, l1, l0)
				mBase = m.M
				if v14 != 0 {
					v20 = v5
				} else {
					v15 = int32(0)
					v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
					if v16 != 0 {
						v20 = v15
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(10)
						v20 = v15
					}
				}
			}
		}
	}
	return v20
}
func F_WebPSafeCalloc(m *base.Module, l0 int64, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v9 int64
	_ = v9
	var v10 int32
	_ = v10
	var v11 int64
	_ = v11
	var v19 int32
	_ = v19
	var v25 int64
	_ = v25
	var v26 int32
	_ = v26
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v44 int32
	_ = v44
	var v50 int32
	_ = v50
	var v52 int32
	_ = v52
	if l0 == int64(0) {
		v19 = base.I32_wrap_i64(l0)
		if v19 != 0 {
			v25 = base.I64_extend_i32_u(v19) * base.I64_extend_i32_u(l1)
			v26 = base.I32_wrap_i64(v25)
			if base.Ui32(l1|v19) < base.Ui32(int32(_a_F_WebPSafeCalloc_0)) {
				v37 = v26
			} else {
				if base.I32_wrap_i64(int64(base.Ui64(v25)>>(uint(int64(32))%64))) != int32(0) {
					v36 = int32(-1)
				} else {
					v36 = v26
				}
				v37 = v36
			}
		} else {
			v37 = int32(0)
		}
		v39 = F_dlmalloc(m, v37)
		mBase = m.M
		if v39 == int32(0) {
		} else {
			v44 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v39+int32(-4)))))
			if v44&int32(3) == int32(0) {
			} else {
				v50 = F_memset(m, v39, int32(0), v37)
				mBase = m.M
			}
		}
		v52 = v39
	} else {
		v9 = base.I64_div_u_s(int64(2147418112), l0)
		v10 = int32(0)
		v11 = base.I64_extend_i32_u(l1)
		if base.Ui64(int64(4294967295)) < base.Ui64(v11*l0) {
			v52 = v10
		} else {
			if base.Ui64(v9) < base.Ui64(v11) {
				v52 = v10
			} else {
				v19 = base.I32_wrap_i64(l0)
				if v19 != 0 {
					v25 = base.I64_extend_i32_u(v19) * base.I64_extend_i32_u(l1)
					v26 = base.I32_wrap_i64(v25)
					if base.Ui32(l1|v19) < base.Ui32(int32(_a_F_WebPSafeCalloc_0)) {
						v37 = v26
					} else {
						if base.I32_wrap_i64(int64(base.Ui64(v25)>>(uint(int64(32))%64))) != int32(0) {
							v36 = int32(-1)
						} else {
							v36 = v26
						}
						v37 = v36
					}
				} else {
					v37 = int32(0)
				}
				v39 = F_dlmalloc(m, v37)
				mBase = m.M
				if v39 == int32(0) {
				} else {
					v44 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v39+int32(-4)))))
					if v44&int32(3) == int32(0) {
					} else {
						v50 = F_memset(m, v39, int32(0), v37)
						mBase = m.M
					}
				}
				v52 = v39
			}
		}
	}
	return v52
}
func F_WebPSafeFree(m *base.Module, l0 int32) {
	F_dlfree(m, l0)
	return
}
func F_WebPSafeMalloc(m *base.Module, l0 int64, l1 int32) int32 {
	var v9 int64
	_ = v9
	var v10 int32
	_ = v10
	var v11 int64
	_ = v11
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	if l0 == int64(0) {
		v21 = F_dlmalloc(m, base.I32_wrap_i64(l0)*l1)
		v23 = v21
	} else {
		v9 = base.I64_div_u_s(int64(2147418112), l0)
		v10 = int32(0)
		v11 = base.I64_extend_i32_u(l1)
		if base.Ui64(int64(4294967295)) < base.Ui64(v11*l0) {
			v23 = v10
		} else {
			if base.Ui64(v9) < base.Ui64(v11) {
				v23 = v10
			} else {
				v21 = F_dlmalloc(m, base.I32_wrap_i64(l0)*l1)
				v23 = v21
			}
		}
	}
	return v23
}
func F_WebPValidateConfig(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v9 float32
	_ = v9
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v18 float32
	_ = v18
	var v21 int32
	_ = v21
	var v24 int32
	_ = v24
	var v29 int32
	_ = v29
	var v32 int32
	_ = v32
	var v35 int32
	_ = v35
	var v38 int32
	_ = v38
	var v41 int32
	_ = v41
	var v44 int32
	_ = v44
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v53 int32
	_ = v53
	var v57 int32
	_ = v57
	var v60 int32
	_ = v60
	var v63 int32
	_ = v63
	var v66 int32
	_ = v66
	var v69 int32
	_ = v69
	var v70 int32
	_ = v70
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v77 int32
	_ = v77
	var v80 int32
	_ = v80
	var v83 int32
	_ = v83
	var v86 int32
	_ = v86
	var v89 int32
	_ = v89
	var v92 int32
	_ = v92
	var v95 int32
	_ = v95
	var v98 int32
	_ = v98
	var v101 int32
	_ = v101
	var v104 int32
	_ = v104
	v2 = int32(0)
	if l0 == v2 {
		v104 = v2
	} else {
		v9 = *(*float32)(unsafe.Add(mBase, uint32(l0)+4))
		if base.F32_lt(v9, float32(0)) != 0 {
			v104 = v2
		} else {
			if base.F32_gt(v9, float32(100)) != 0 {
				v104 = v2
			} else {
				v14 = int32(0)
				v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
				if v15 < v14 {
					v104 = v14
				} else {
					v18 = *(*float32)(unsafe.Add(mBase, uint32(l0)+20))
					if base.F32_lt(v18, float32(0)) != 0 {
						v104 = v14
					} else {
						v21 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
						if base.Ui32(int32(6)) < base.Ui32(v21) {
							v104 = v14
						} else {
							v24 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
							if base.Ui32(v24+int32(-5)) < base.Ui32(int32(-4)) {
								v104 = v14
							} else {
								v29 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
								if base.Ui32(int32(100)) < base.Ui32(v29) {
									v104 = v14
								} else {
									v32 = *(*int32)(unsafe.Add(mBase, uint32(l0)+32))
									if base.Ui32(int32(100)) < base.Ui32(v32) {
										v104 = v14
									} else {
										v35 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
										if base.Ui32(int32(7)) < base.Ui32(v35) {
											v104 = v14
										} else {
											v38 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
											if base.Ui32(int32(1)) < base.Ui32(v38) {
												v104 = v14
											} else {
												v41 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
												if base.Ui32(int32(1)) < base.Ui32(v41) {
													v104 = v14
												} else {
													v44 = *(*int32)(unsafe.Add(mBase, uint32(l0)+60))
													if base.Ui32(v44+int32(-11)) < base.Ui32(int32(-10)) {
														v104 = v14
													} else {
														v49 = int32(0)
														v50 = *(*int32)(unsafe.Add(mBase, uint32(l0)+108))
														if v50 < v49 {
															v104 = v49
														} else {
															v53 = *(*int32)(unsafe.Add(mBase, uint32(l0)+112))
															if int32(100) < v53 {
																v104 = v49
															} else {
																if v53 < v50 {
																	v104 = v49
																} else {
																	v57 = *(*int32)(unsafe.Add(mBase, uint32(l0)+64))
																	if base.Ui32(int32(1)) < base.Ui32(v57) {
																		v104 = v49
																	} else {
																		v60 = *(*int32)(unsafe.Add(mBase, uint32(l0)+68))
																		if base.Ui32(int32(7)) < base.Ui32(v60) {
																			v104 = v49
																		} else {
																			v63 = *(*int32)(unsafe.Add(mBase, uint32(l0)+72))
																			if base.Ui32(int32(3)) < base.Ui32(v63) {
																				v104 = v49
																			} else {
																				v66 = *(*int32)(unsafe.Add(mBase, uint32(l0)+76))
																				if base.Ui32(int32(100)) < base.Ui32(v66) {
																					v104 = v49
																				} else {
																					v69 = int32(0)
																					v70 = *(*int32)(unsafe.Add(mBase, uint32(l0)+48))
																					if v70 < v69 {
																						v104 = v69
																					} else {
																						v73 = int32(0)
																						v74 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
																						if v74 < v73 {
																							v104 = v73
																						} else {
																							v77 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
																							if base.Ui32(int32(100)) < base.Ui32(v77) {
																								v104 = v73
																							} else {
																								v80 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
																								if base.Ui32(int32(1)) < base.Ui32(v80) {
																									v104 = v73
																								} else {
																									v83 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
																									if base.Ui32(int32(100)) < base.Ui32(v83) {
																										v104 = v73
																									} else {
																										v86 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
																										if base.Ui32(int32(3)) < base.Ui32(v86) {
																											v104 = v73
																										} else {
																											v89 = *(*int32)(unsafe.Add(mBase, uint32(l0)+80))
																											if base.Ui32(int32(1)) < base.Ui32(v89) {
																												v104 = v73
																											} else {
																												v92 = *(*int32)(unsafe.Add(mBase, uint32(l0)+84))
																												if base.Ui32(int32(1)) < base.Ui32(v92) {
																													v104 = v73
																												} else {
																													v95 = *(*int32)(unsafe.Add(mBase, uint32(l0)+88))
																													if base.Ui32(int32(1)) < base.Ui32(v95) {
																														v104 = v73
																													} else {
																														v98 = *(*int32)(unsafe.Add(mBase, uint32(l0)+96))
																														if base.Ui32(int32(1)) < base.Ui32(v98) {
																															v104 = v73
																														} else {
																															v101 = *(*int32)(unsafe.Add(mBase, uint32(l0)+104))
																															v104 = base.B2i32(base.Ui32(v101) < base.Ui32(int32(2)))
																														}
																													}
																												}
																											}
																										}
																									}
																								}
																							}
																						}
																					}
																				}
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return v104
}
func F_WebPValidatePicture(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v7 int32
	_ = v7
	var v10 int32
	_ = v10
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v18 int32
	_ = v18
	var v22 int32
	_ = v22
	if l0 != 0 {
		v6 = int32(5)
		v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
		if v7 < int32(1) {
			v16 = v6
			v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
			if v18 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v16
			}
			v22 = int32(0)
		} else {
			v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
			if v10 < int32(1) {
				v16 = v6
				v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
				if v18 != 0 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v16
				}
				v22 = int32(0)
			} else {
				v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
				switch v15 {
				case 0, 4:
					v22 = int32(1)
				default:
					v16 = int32(4)
					v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
					if v18 != 0 {
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v16
					}
					v22 = int32(0)
				}
			}
		}
		return v22
	} else {
		return int32(0)
	}
}
