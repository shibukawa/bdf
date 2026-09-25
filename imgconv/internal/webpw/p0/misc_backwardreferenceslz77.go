//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	"unsafe"
)

func F_BackwardReferencesLz77(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) int32 {
	mBase := m.M
	_ = mBase
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v28 int32
	_ = v28
	var v37 int32
	_ = v37
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
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
	var v69 int32
	_ = v69
	var v70 int32
	_ = v70
	var v72 int32
	_ = v72
	var v82 int32
	_ = v82
	var v93 int32
	_ = v93
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
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
	var v109 int32
	_ = v109
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v122 int32
	_ = v122
	var v125 int32
	_ = v125
	var v128 int32
	_ = v128
	var v135 int32
	_ = v135
	var v147 int64
	_ = v147
	var v150 int32
	_ = v150
	var v151 int32
	_ = v151
	var v154 int32
	_ = v154
	var v155 int64
	_ = v155
	var v156 int32
	_ = v156
	var v160 int32
	_ = v160
	var v167 int64
	_ = v167
	var v168 int32
	_ = v168
	var v169 int64
	_ = v169
	var v179 int32
	_ = v179
	var v181 int32
	_ = v181
	var v183 int32
	_ = v183
	var v184 int32
	_ = v184
	var v191 int32
	_ = v191
	var v193 int32
	_ = v193
	var v194 int32
	_ = v194
	var v196 int32
	_ = v196
	var v203 int32
	_ = v203
	var v204 int32
	_ = v204
	var v205 int32
	_ = v205
	var v209 int32
	_ = v209
	var v218 int32
	_ = v218
	var v219 int32
	_ = v219
	var v230 int32
	_ = v230
	var v239 int32
	_ = v239
	var v240 int32
	_ = v240
	var v243 int32
	_ = v243
	var v244 int64
	_ = v244
	var v245 int32
	_ = v245
	var v249 int32
	_ = v249
	var v256 int64
	_ = v256
	var v257 int32
	_ = v257
	var v258 int64
	_ = v258
	var v268 int32
	_ = v268
	var v270 int32
	_ = v270
	var v272 int32
	_ = v272
	var v279 int32
	_ = v279
	var v281 int32
	_ = v281
	var v282 int32
	_ = v282
	var v284 int32
	_ = v284
	var v291 int32
	_ = v291
	var v292 int32
	_ = v292
	var v296 int32
	_ = v296
	var v319 int32
	_ = v319
	var v323 int32
	_ = v323
	var v330 int32
	_ = v330
	var v350 int32
	_ = v350
	v19 = l1 * l0
	v20 = *(*int32)(unsafe.Add(mBase, uint32(l4)+12))
	if v20 == int32(0) {
		goto L1
	} else {
		goto L2
	}
L1:
	;
	v25 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l4)+20)) = v25
	v28 = *(*int32)(unsafe.Add(mBase, uint32(l4)+8))
	*(*int32)(unsafe.Add(mBase, uint32(l4)+16)) = v28
	*(*int32)(unsafe.Add(mBase, uint32(l4)+12)) = l4 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(l4)+8)) = v25
	if v19 < int32(1) {
		goto L3
	} else {
		goto L4
	}
L2:
	;
	v23 = *(*int32)(unsafe.Add(mBase, uint32(l4)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v20))) = v23
	goto L1
L3:
	;
	v350 = *(*int32)(unsafe.Add(mBase, uint32(l4)+4))
	return base.B2i32(v350 == int32(0))
L4:
	;
	v37 = int32(-1)
	v46 = v25
	v47 = v37
	goto L5
L5:
	;
	v58 = *(*int32)(unsafe.Add(mBase, uint32(l3)))
	v60 = v46 << (uint(int32(2)) % 32)
	v62 = *(*int32)(unsafe.Add(mBase, uint32(v58+v60)))
	v64 = v62 & int32(4095)
	if base.Ui32(int32(3)) < base.Ui32(v64) {
		goto L11
	} else {
		goto L12
	}
L6:
	;
	goto L3
L7:
	;
	v330 = v323 + v46
	if v330 < v19 {
		v46 = v330
		v47 = v319
		goto L5
	} else {
		goto L62
	}
L8:
	;
	if v219 == int32(0) {
		goto L48
	} else {
		goto L49
	}
L9:
	;
	v218 = *(*int32)(unsafe.Add(mBase, uint32(l4)+20))
	v219 = v218
	v230 = v64
	goto L8
L10:
	;
	v147 = int64(*(*uint32)(unsafe.Add(mBase, uint32(l2+v60))))
	if v128 == int32(0) {
		goto L33
	} else {
		goto L34
	}
L11:
	;
	if v47 < v46 {
		goto L13
	} else {
		goto L14
	}
L12:
	;
	v67 = *(*int32)(unsafe.Add(mBase, uint32(l4)+20))
	v128 = v67
	v135 = v47
	goto L10
L13:
	;
	v69 = v46
	goto L15
L14:
	;
	v69 = v47
	goto L15
L15:
	;
	v70 = v64 + v46
	if v70 < v19 {
		goto L16
	} else {
		goto L17
	}
L16:
	;
	v72 = v70
	goto L18
L17:
	;
	v72 = v19 + v37
	goto L18
L18:
	;
	if v72 <= v69 {
		goto L9
	} else {
		goto L19
	}
L19:
	;
	v82 = v58 + v69<<(uint(int32(2))%32) + int32(4)
	v93 = v64
	v96 = int32(0)
	v97 = v69
	goto L22
L20:
	;
	v125 = *(*int32)(unsafe.Add(mBase, uint32(l4)+20))
	if v122 != int32(1) {
		v219 = v125
		v230 = v122
		goto L8
	} else {
		goto L31
	}
L21:
	;
	v122 = v97 - v46 + int32(1)
	goto L20
L22:
	;
	v100 = int32(1)
	v101 = v97 + v100
	v102 = *(*int32)(unsafe.Add(mBase, uint32(v82)))
	v104 = v102 & int32(4095)
	if base.Ui32(int32(3)) < base.Ui32(v104) {
		goto L25
	} else {
		goto L26
	}
L24:
	;
	if v72 != v101 {
		v82 = v82 + int32(4)
		v93 = v113
		v96 = v114
		v97 = v101
		goto L22
	} else {
		goto L30
	}
L25:
	;
	v108 = v104
	goto L27
L26:
	;
	v108 = v100
	goto L27
L27:
	;
	v109 = v101 + v108
	if v109 <= v96 {
		v113 = v93
		v114 = v96
		goto L24
	} else {
		goto L28
	}
L28:
	;
	if v19 <= v109 {
		goto L21
	} else {
		goto L29
	}
L29:
	;
	v113 = int32(1) - v46 + v97
	v114 = v109
	goto L24
L30:
	;
	v122 = v113
	goto L20
L31:
	;
	v128 = v125
	v135 = v69
	goto L10
L32:
	;
	v205 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v203)+8)) = v204 + v205
	v209 = *(*int32)(unsafe.Add(mBase, uint32(v203)+4))
	*(*int64)(unsafe.Add(mBase, uint32(v209+v204<<(uint(int32(3))%32)))) = v147<<(uint(int64(32))%64) | int64(65536)
	v319 = v135
	v323 = v205
	goto L7
L33:
	;
	v154 = *(*int32)(unsafe.Add(mBase, uint32(l4)+16))
	if v154 != 0 {
		goto L37
	} else {
		goto L38
	}
L34:
	;
	v150 = *(*int32)(unsafe.Add(mBase, uint32(v128)+8))
	v151 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
	if v150 != v151 {
		v203 = v128
		v204 = v150
		goto L32
	} else {
		goto L35
	}
L35:
	;
	goto L33
L36:
	;
	v194 = *(*int32)(unsafe.Add(mBase, uint32(l4)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v194))) = v193
	v196 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v193)+8)) = v196
	*(*int32)(unsafe.Add(mBase, uint32(l4)+20)) = v193
	*(*int32)(unsafe.Add(mBase, uint32(l4)+12)) = v193
	*(*int32)(unsafe.Add(mBase, uint32(v193))) = v196
	v203 = v193
	v204 = v196
	goto L32
L37:
	;
	v191 = *(*int32)(unsafe.Add(mBase, uint32(v154)))
	*(*int32)(unsafe.Add(mBase, uint32(l4)+16)) = v191
	v193 = v154
	goto L36
L38:
	;
	v155 = int64(1)
	v156 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
	v160 = v156<<(uint(int32(3))%32) + int32(12)
	goto L43
L39:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v181)+4)) = v181 + int32(12)
	v193 = v181
	goto L36
L40:
	;
	if v181 != 0 {
		goto L39
	} else {
		goto L46
	}
L41:
	;
	goto L40
L42:
	;
	v179 = F_malloc(m, base.I32_wrap_i64(v155)*v160)
	mBase = m.M
	v181 = v179
	goto L41
L43:
	;
	v167 = base.I64_div_u_s(int64(2147418112), v155)
	v168 = int32(0)
	v169 = base.I64_extend_i32_u(v160)
	if base.Ui64(int64(4294967295)) < base.Ui64(v169*v155) {
		v181 = v168
		goto L41
	} else {
		goto L44
	}
L44:
	;
	if base.Ui64(v167) < base.Ui64(v169) {
		v181 = v168
		goto L41
	} else {
		goto L45
	}
L45:
	;
	goto L42
L46:
	;
	v183 = int32(1)
	v184 = *(*int32)(unsafe.Add(mBase, uint32(l4)+4))
	*(*int32)(unsafe.Add(mBase, uint32(l4)+4)) = v184 | v183
	v319 = v135
	v323 = v183
	goto L7
L47:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v291)+8)) = v292 + int32(1)
	v296 = *(*int32)(unsafe.Add(mBase, uint32(v291)+4))
	*(*int64)(unsafe.Add(mBase, uint32(v296+v292<<(uint(int32(3))%32)))) = base.I64_extend_i32_u(int32(base.Ui32(v62)>>(uint(int32(12))%32)))<<(uint(int64(32))%64) | base.I64_extend_i32_u(v230<<(uint(int32(16))%32)) | int64(2)
	v319 = v69
	v323 = v230
	goto L7
L48:
	;
	v243 = *(*int32)(unsafe.Add(mBase, uint32(l4)+16))
	if v243 != 0 {
		goto L52
	} else {
		goto L53
	}
L49:
	;
	v239 = *(*int32)(unsafe.Add(mBase, uint32(v219)+8))
	v240 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
	if v239 != v240 {
		v291 = v219
		v292 = v239
		goto L47
	} else {
		goto L50
	}
L50:
	;
	goto L48
L51:
	;
	v282 = *(*int32)(unsafe.Add(mBase, uint32(l4)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v282))) = v281
	v284 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v281)+8)) = v284
	*(*int32)(unsafe.Add(mBase, uint32(l4)+20)) = v281
	*(*int32)(unsafe.Add(mBase, uint32(l4)+12)) = v281
	*(*int32)(unsafe.Add(mBase, uint32(v281))) = v284
	v291 = v281
	v292 = v284
	goto L47
L52:
	;
	v279 = *(*int32)(unsafe.Add(mBase, uint32(v243)))
	*(*int32)(unsafe.Add(mBase, uint32(l4)+16)) = v279
	v281 = v243
	goto L51
L53:
	;
	v244 = int64(1)
	v245 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
	v249 = v245<<(uint(int32(3))%32) + int32(12)
	goto L58
L54:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v270)+4)) = v270 + int32(12)
	v281 = v270
	goto L51
L55:
	;
	if v270 != 0 {
		goto L54
	} else {
		goto L61
	}
L56:
	;
	goto L55
L57:
	;
	v268 = F_malloc(m, base.I32_wrap_i64(v244)*v249)
	mBase = m.M
	v270 = v268
	goto L56
L58:
	;
	v256 = base.I64_div_u_s(int64(2147418112), v244)
	v257 = int32(0)
	v258 = base.I64_extend_i32_u(v249)
	if base.Ui64(int64(4294967295)) < base.Ui64(v258*v244) {
		v270 = v257
		goto L56
	} else {
		goto L59
	}
L59:
	;
	if base.Ui64(v256) < base.Ui64(v258) {
		v270 = v257
		goto L56
	} else {
		goto L60
	}
L60:
	;
	goto L57
L61:
	;
	v272 = *(*int32)(unsafe.Add(mBase, uint32(l4)+4))
	*(*int32)(unsafe.Add(mBase, uint32(l4)+4)) = v272 | int32(1)
	v319 = v69
	v323 = v230
	goto L7
L62:
	;
	goto L6
}
