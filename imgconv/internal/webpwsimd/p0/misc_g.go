//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_GetBestColorTransformForTile(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32, l10 int32, l11 int32) {
	mBase := m.M
	_ = mBase
	var v13 int32
	_ = v13
	var v32 int32
	_ = v32
	var v34 int32
	_ = v34
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v42 int32
	_ = v42
	var v45 int32
	_ = v45
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
	var v55 int32
	_ = v55
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v182 int32
	_ = v182
	var v183 int32
	_ = v183
	var v185 int32
	_ = v185
	var v186 int32
	_ = v186
	var v187 int64
	_ = v187
	var v188 int64
	_ = v188
	var v198 int32
	_ = v198
	var v203 int32
	_ = v203
	var v215 int64
	_ = v215
	var v216 int64
	_ = v216
	var v228 int32
	_ = v228
	var v230 int32
	_ = v230
	var v233 int64
	_ = v233
	var v236 int64
	_ = v236
	var v239 int64
	_ = v239
	var v240 int64
	_ = v240
	var v247 int64
	_ = v247
	var v252 int64
	_ = v252
	var v254 int32
	_ = v254
	var v261 int32
	_ = v261
	var v268 int64
	_ = v268
	var v271 int64
	_ = v271
	var v272 int64
	_ = v272
	var v276 int32
	_ = v276
	var v277 int64
	_ = v277
	var v281 int32
	_ = v281
	var v282 int64
	_ = v282
	var v289 int32
	_ = v289
	var v303 int32
	_ = v303
	var v310 int64
	_ = v310
	var v318 int32
	_ = v318
	var v323 int32
	_ = v323
	var v339 int32
	_ = v339
	var v340 int32
	_ = v340
	var v347 int64
	_ = v347
	var v359 int32
	_ = v359
	var v482 int32
	_ = v482
	var v483 int32
	_ = v483
	var v485 int32
	_ = v485
	var v486 int32
	_ = v486
	var v487 int64
	_ = v487
	var v488 int64
	_ = v488
	var v496 int32
	_ = v496
	var v501 int32
	_ = v501
	var v513 int64
	_ = v513
	var v514 int64
	_ = v514
	var v526 int32
	_ = v526
	var v528 int32
	_ = v528
	var v531 int64
	_ = v531
	var v534 int64
	_ = v534
	var v537 int64
	_ = v537
	var v538 int64
	_ = v538
	var v545 int64
	_ = v545
	var v550 int64
	_ = v550
	var v552 int32
	_ = v552
	var v559 int64
	_ = v559
	var v562 int64
	_ = v562
	var v563 int64
	_ = v563
	var v567 int32
	_ = v567
	var v569 int64
	_ = v569
	var v573 int64
	_ = v573
	var v576 int64
	_ = v576
	var v577 int32
	_ = v577
	var v578 int64
	_ = v578
	var v579 int32
	_ = v579
	var v580 int32
	_ = v580
	var v597 int32
	_ = v597
	var v617 int32
	_ = v617
	var v618 int32
	_ = v618
	var v619 int32
	_ = v619
	var v620 int32
	_ = v620
	var v621 int32
	_ = v621
	var v745 int32
	_ = v745
	var v746 int32
	_ = v746
	var v748 int32
	_ = v748
	var v749 int32
	_ = v749
	var v750 int64
	_ = v750
	var v751 int64
	_ = v751
	var v761 int32
	_ = v761
	var v766 int32
	_ = v766
	var v778 int64
	_ = v778
	var v779 int64
	_ = v779
	var v791 int32
	_ = v791
	var v793 int32
	_ = v793
	var v796 int64
	_ = v796
	var v799 int64
	_ = v799
	var v802 int64
	_ = v802
	var v803 int64
	_ = v803
	var v810 int64
	_ = v810
	var v815 int64
	_ = v815
	var v817 int32
	_ = v817
	var v824 int64
	_ = v824
	var v827 int64
	_ = v827
	var v828 int64
	_ = v828
	var v832 int32
	_ = v832
	var v833 int64
	_ = v833
	var v837 int32
	_ = v837
	var v838 int64
	_ = v838
	var v842 int32
	_ = v842
	var v843 int64
	_ = v843
	var v847 int32
	_ = v847
	var v848 int64
	_ = v848
	var v855 int32
	_ = v855
	var v862 int32
	_ = v862
	var v871 int32
	_ = v871
	var v878 int32
	_ = v878
	var v883 int32
	_ = v883
	var v884 int64
	_ = v884
	var v899 int32
	_ = v899
	var v903 int32
	_ = v903
	var v906 int32
	_ = v906
	var v918 int32
	_ = v918
	var v923 int32
	_ = v923
	var v924 int64
	_ = v924
	var v1058 int32
	_ = v1058
	var v1063 int32
	_ = v1063
	var v1064 int32
	_ = v1064
	var v1066 int32
	_ = v1066
	var v1067 int32
	_ = v1067
	var v1069 int32
	_ = v1069
	var v1070 int32
	_ = v1070
	var v1071 int32
	_ = v1071
	var v1073 int32
	_ = v1073
	var v1074 int32
	_ = v1074
	var v1075 int64
	_ = v1075
	var v1076 int64
	_ = v1076
	var v1084 int32
	_ = v1084
	var v1089 int32
	_ = v1089
	var v1101 int64
	_ = v1101
	var v1102 int64
	_ = v1102
	var v1114 int32
	_ = v1114
	var v1116 int32
	_ = v1116
	var v1119 int64
	_ = v1119
	var v1122 int64
	_ = v1122
	var v1125 int64
	_ = v1125
	var v1126 int64
	_ = v1126
	var v1133 int64
	_ = v1133
	var v1138 int64
	_ = v1138
	var v1140 int32
	_ = v1140
	var v1147 int64
	_ = v1147
	var v1150 int64
	_ = v1150
	var v1151 int64
	_ = v1151
	var v1155 int32
	_ = v1155
	var v1157 int64
	_ = v1157
	var v1161 int64
	_ = v1161
	var v1165 int32
	_ = v1165
	var v1167 int64
	_ = v1167
	var v1171 int64
	_ = v1171
	var v1174 int64
	_ = v1174
	var v1177 int64
	_ = v1177
	var v1178 int32
	_ = v1178
	var v1179 int64
	_ = v1179
	var v1180 int32
	_ = v1180
	var v1181 int32
	_ = v1181
	var v1195 int32
	_ = v1195
	var v1197 int32
	_ = v1197
	var v1200 int32
	_ = v1200
	var v1201 int32
	_ = v1201
	v13 = int32(0)
	v32 = m.G0
	v34 = v32 - int32(1024)
	m.G0 = v34
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+1)) = uint16(v13)
	v38 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v39 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5))))
	v40 = l2 << (uint(l3) % 32)
	v42 = int32(2)
	v45 = l1 << (uint(l3) % 32)
	v48 = l11 + v40*l7<<(uint(v42)%32) + v45<<(uint(v42)%32)
	v50 = int32(1) << (uint(l3) % 32)
	v51 = v45 + v50
	if v51 < l7 {
		goto L1
	} else {
		goto L2
	}
L1:
	;
	v53 = v51
	goto L3
L2:
	;
	v53 = l7
	goto L3
L3:
	;
	v54 = v53 - v45
	v55 = v40 + v50
	if v55 < l8 {
		goto L4
	} else {
		goto L5
	}
L4:
	;
	v57 = v55
	goto L6
L5:
	;
	v57 = l8
	goto L6
L6:
	;
	v58 = v57 - v40
	v59 = int32(0)
	goto L9
L7:
	;
	v182 = m.G96
	v183 = *(*int32)(unsafe.Add(mBase, uint32(v182)))
	m.T0[v183].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v48, l7, v54, v58, v59, v34)
	mBase = m.M
	v185 = m.G99
	v186 = *(*int32)(unsafe.Add(mBase, uint32(v185)))
	v187 = m.T0[v186].(func(*base.Module, int32, int32) int64)(m, v34, l9)
	mBase = m.M
	v188 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v34))))
	v198 = v34 + int32(1020)
	v203 = int32(4)
	v215 = v188 * int64(25165824)
	v216 = int64(2013265920)
	goto L21
L9:
	;
	base.MemoryFill(m, v34, v59, int32(1024))
	goto L7
L21:
	;
	v228 = *(*int32)(unsafe.Add(mBase, uint32(v198)))
	v230 = *(*int32)(unsafe.Add(mBase, uint32(v34+v203)))
	v233 = v216 * base.I64_extend_i32_u(v228+v230)
	if v233 < int64(0) {
		goto L23
	} else {
		goto L24
	}
L22:
	;
	v261 = l6 * int32(7) >> (uint(int32(8)) % 32)
	if v261 < int32(-3) {
		v597 = int32(0)
		goto L30
	} else {
		goto L31
	}
L23:
	;
	v236 = int64(-50)
	goto L25
L24:
	;
	v236 = int64(50)
	goto L25
L25:
	;
	v239 = base.I64_div_s(v236+v233, int64(100))
	v240 = v239 + v215
	if v216 < int64(0) {
		goto L26
	} else {
		goto L27
	}
L26:
	;
	v247 = int64(-5)
	goto L28
L27:
	;
	v247 = int64(5)
	goto L28
L28:
	;
	v252 = base.I64_div_s(v247+v216*int64(6), int64(10))
	v254 = v203 + int32(4)
	if v254 != int32(64) {
		v198 = v198 + int32(-4)
		v203 = v254
		v215 = v240
		v216 = v252
		goto L21
	} else {
		goto L29
	}
L29:
	;
	goto L22
L30:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v597)
	v617 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4)+1)))
	v618 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4)+2)))
	v619 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5)+1)))
	v620 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5)+2)))
	v621 = int32(0)
	goto L90
L31:
	;
	if v240 < int64(0) {
		goto L32
	} else {
		goto L33
	}
L32:
	;
	v268 = int64(-5)
	goto L34
L33:
	;
	v268 = int64(5)
	goto L34
L34:
	;
	v271 = base.I64_div_s(v268+v240, int64(-10))
	v272 = v271 + v187
	v276 = v38 & int32(255)
	if v276 != 0 {
		goto L35
	} else {
		goto L36
	}
L35:
	;
	v277 = v272
	goto L37
L36:
	;
	v277 = v272 + int64(-25165824)
	goto L37
L37:
	;
	v281 = v39 & int32(255)
	if v281 != 0 {
		goto L38
	} else {
		goto L39
	}
L38:
	;
	v282 = v277
	goto L40
L39:
	;
	v282 = v277 + int64(-25165824)
	goto L40
L40:
	;
	v289 = int32(0)
	v303 = v289
	v310 = v282 + int64(-25165824)
	v318 = v289
	goto L41
L41:
	;
	v323 = int32(base.Ui32(int32(32)) >> (uint(v318) % 32))
	v339 = int32(0) - v323
	v340 = v303
	v347 = v310
	goto L43
L42:
	;
	v597 = v579
	goto L30
L43:
	;
	v359 = v339 + v340
	goto L47
L44:
	;
	if v318 != v261+int32(3) {
		v303 = v579
		v310 = v578
		v318 = v318 + int32(1)
		goto L41
	} else {
		goto L87
	}
L45:
	;
	v482 = m.G96
	v483 = *(*int32)(unsafe.Add(mBase, uint32(v482)))
	m.T0[v483].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v48, l7, v54, v58, v359, v34)
	mBase = m.M
	v485 = m.G99
	v486 = *(*int32)(unsafe.Add(mBase, uint32(v485)))
	v487 = m.T0[v486].(func(*base.Module, int32, int32) int64)(m, v34, l9)
	mBase = m.M
	v488 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v34))))
	v496 = int32(4)
	v501 = v34 + int32(1020)
	v513 = v488 * int64(25165824)
	v514 = int64(2013265920)
	goto L59
L47:
	;
	base.MemoryFill(m, v34, int32(0), int32(1024))
	goto L45
L59:
	;
	v526 = *(*int32)(unsafe.Add(mBase, uint32(v501)))
	v528 = *(*int32)(unsafe.Add(mBase, uint32(v34+v496)))
	v531 = v514 * base.I64_extend_i32_u(v526+v528)
	if v531 < int64(0) {
		goto L61
	} else {
		goto L62
	}
L60:
	;
	if v538 < int64(0) {
		goto L68
	} else {
		goto L69
	}
L61:
	;
	v534 = int64(-50)
	goto L63
L62:
	;
	v534 = int64(50)
	goto L63
L63:
	;
	v537 = base.I64_div_s(v534+v531, int64(100))
	v538 = v537 + v513
	if v514 < int64(0) {
		goto L64
	} else {
		goto L65
	}
L64:
	;
	v545 = int64(-5)
	goto L66
L65:
	;
	v545 = int64(5)
	goto L66
L66:
	;
	v550 = base.I64_div_s(v545+v514*int64(6), int64(10))
	v552 = v496 + int32(4)
	if v552 != int32(64) {
		v496 = v552
		v501 = v501 + int32(-4)
		v513 = v538
		v514 = v550
		goto L59
	} else {
		goto L67
	}
L67:
	;
	goto L60
L68:
	;
	v559 = int64(-5)
	goto L70
L69:
	;
	v559 = int64(5)
	goto L70
L70:
	;
	v562 = base.I64_div_s(v559+v538, int64(-10))
	v563 = v562 + v487
	v567 = v359 & int32(255)
	if v276 == v567 {
		goto L71
	} else {
		goto L72
	}
L71:
	;
	v569 = v563 + int64(-25165824)
	goto L73
L72:
	;
	v569 = v563
	goto L73
L73:
	;
	if v281 == v567 {
		goto L74
	} else {
		goto L75
	}
L74:
	;
	v573 = v569 + int64(-25165824)
	goto L76
L75:
	;
	v573 = v569
	goto L76
L76:
	;
	if v359 != 0 {
		goto L77
	} else {
		goto L78
	}
L77:
	;
	v576 = v573
	goto L79
L78:
	;
	v576 = v573 + int64(-25165824)
	goto L79
L79:
	;
	v577 = base.B2i32(v576 < v347)
	if v576 < v347 {
		goto L80
	} else {
		goto L81
	}
L80:
	;
	v578 = v576
	goto L82
L81:
	;
	v578 = v347
	goto L82
L82:
	;
	if v576 < v347 {
		goto L83
	} else {
		goto L84
	}
L83:
	;
	v579 = v359
	goto L85
L84:
	;
	v579 = v340
	goto L85
L85:
	;
	v580 = v339 + v323<<(uint(int32(1))%32)
	if v580 <= v323 {
		v339 = v580
		v340 = v579
		v347 = v578
		goto L43
	} else {
		goto L86
	}
L86:
	;
	goto L44
L87:
	;
	goto L42
L88:
	;
	v745 = m.G95
	v746 = *(*int32)(unsafe.Add(mBase, uint32(v745)))
	m.T0[v746].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32))(m, v48, l7, v54, v58, v621, v621, v34)
	mBase = m.M
	v748 = m.G99
	v749 = *(*int32)(unsafe.Add(mBase, uint32(v748)))
	v750 = m.T0[v749].(func(*base.Module, int32, int32) int64)(m, v34, l10)
	mBase = m.M
	v751 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v34))))
	v761 = v34 + int32(1020)
	v766 = int32(4)
	v778 = v751 * int64(25165824)
	v779 = int64(2013265920)
	goto L102
L90:
	;
	base.MemoryFill(m, v34, v621, int32(1024))
	goto L88
L102:
	;
	v791 = *(*int32)(unsafe.Add(mBase, uint32(v761)))
	v793 = *(*int32)(unsafe.Add(mBase, uint32(v34+v766)))
	v796 = v779 * base.I64_extend_i32_u(v791+v793)
	if v796 < int64(0) {
		goto L104
	} else {
		goto L105
	}
L103:
	;
	if v803 < int64(0) {
		goto L111
	} else {
		goto L112
	}
L104:
	;
	v799 = int64(-50)
	goto L106
L105:
	;
	v799 = int64(50)
	goto L106
L106:
	;
	v802 = base.I64_div_s(v799+v796, int64(100))
	v803 = v802 + v778
	if v779 < int64(0) {
		goto L107
	} else {
		goto L108
	}
L107:
	;
	v810 = int64(-5)
	goto L109
L108:
	;
	v810 = int64(5)
	goto L109
L109:
	;
	v815 = base.I64_div_s(v810+v779*int64(6), int64(10))
	v817 = v766 + int32(4)
	if v817 != int32(64) {
		v761 = v761 + int32(-4)
		v766 = v817
		v778 = v803
		v779 = v815
		goto L102
	} else {
		goto L110
	}
L110:
	;
	goto L103
L111:
	;
	v824 = int64(-5)
	goto L113
L112:
	;
	v824 = int64(5)
	goto L113
L113:
	;
	v827 = base.I64_div_s(v824+v803, int64(-10))
	v828 = v827 + v750
	v832 = v617 & int32(255)
	if v832 != 0 {
		goto L114
	} else {
		goto L115
	}
L114:
	;
	v833 = v828
	goto L116
L115:
	;
	v833 = v828 + int64(-25165824)
	goto L116
L116:
	;
	v837 = v619 & int32(255)
	if v837 != 0 {
		goto L117
	} else {
		goto L118
	}
L117:
	;
	v838 = v833
	goto L119
L118:
	;
	v838 = v833 + int64(-25165824)
	goto L119
L119:
	;
	v842 = v618 & int32(255)
	if v842 != 0 {
		goto L120
	} else {
		goto L121
	}
L120:
	;
	v843 = v838
	goto L122
L121:
	;
	v843 = v838 + int64(-25165824)
	goto L122
L122:
	;
	v847 = v620 & int32(255)
	if v847 != 0 {
		goto L123
	} else {
		goto L124
	}
L123:
	;
	v848 = v843
	goto L125
L124:
	;
	v848 = v843 + int64(-25165824)
	goto L125
L125:
	;
	if base.Ui32(int32(50)) < base.Ui32(l6) {
		goto L126
	} else {
		goto L127
	}
L126:
	;
	v855 = int32(7)
	goto L128
L127:
	;
	v855 = int32(4)
	goto L128
L128:
	;
	v862 = int32(0)
	v871 = v862
	v878 = v862
	v883 = v862
	v884 = v848 + int64(-50331648)
	goto L130
L129:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v1200)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v1201)
	m.G0 = v34 + int32(1024)
	return
L130:
	;
	v899 = m.G1
	v903 = int32(*(*int8)(unsafe.Add(mBase, uint32(v899+int32(_a_F_GetBestColorTransformForTile_0)+v871))))
	v906 = int32(0)
	v918 = v878
	v923 = v883
	v924 = v884
	goto L132
L131:
	;
	v1197 = int32(0)
	v1200 = v1197
	v1201 = v1197
	goto L129
L132:
	;
	goto L136
L133:
	;
	if base.Ui32(int32(2)) < base.Ui32(v871+int32(-4)) {
		goto L189
	} else {
		goto L190
	}
L134:
	;
	v1058 = m.G1
	v1063 = v1058 + int32(_a_F_GetBestColorTransformForTile_1) + v906<<(uint(int32(1))%32)
	v1064 = int32(*(*int8)(unsafe.Add(mBase, uint32(v1063))))
	v1066 = v1064*v903 + v923
	v1067 = int32(*(*int8)(unsafe.Add(mBase, uint32(v1063)+1)))
	v1069 = v1067*v903 + v918
	v1070 = m.G95
	v1071 = *(*int32)(unsafe.Add(mBase, uint32(v1070)))
	m.T0[v1071].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32))(m, v48, l7, v54, v58, v1066, v1069, v34)
	mBase = m.M
	v1073 = m.G99
	v1074 = *(*int32)(unsafe.Add(mBase, uint32(v1073)))
	v1075 = m.T0[v1074].(func(*base.Module, int32, int32) int64)(m, v34, l10)
	mBase = m.M
	v1076 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v34))))
	v1084 = int32(4)
	v1089 = v34 + int32(1020)
	v1101 = v1076 * int64(25165824)
	v1102 = int64(2013265920)
	goto L148
L136:
	;
	base.MemoryFill(m, v34, int32(0), int32(1024))
	goto L134
L148:
	;
	v1114 = *(*int32)(unsafe.Add(mBase, uint32(v1089)))
	v1116 = *(*int32)(unsafe.Add(mBase, uint32(v34+v1084)))
	v1119 = v1102 * base.I64_extend_i32_u(v1114+v1116)
	if v1119 < int64(0) {
		goto L150
	} else {
		goto L151
	}
L149:
	;
	if v1126 < int64(0) {
		goto L157
	} else {
		goto L158
	}
L150:
	;
	v1122 = int64(-50)
	goto L152
L151:
	;
	v1122 = int64(50)
	goto L152
L152:
	;
	v1125 = base.I64_div_s(v1122+v1119, int64(100))
	v1126 = v1125 + v1101
	if v1102 < int64(0) {
		goto L153
	} else {
		goto L154
	}
L153:
	;
	v1133 = int64(-5)
	goto L155
L154:
	;
	v1133 = int64(5)
	goto L155
L155:
	;
	v1138 = base.I64_div_s(v1133+v1102*int64(6), int64(10))
	v1140 = v1084 + int32(4)
	if v1140 != int32(64) {
		v1084 = v1140
		v1089 = v1089 + int32(-4)
		v1101 = v1126
		v1102 = v1138
		goto L148
	} else {
		goto L156
	}
L156:
	;
	goto L149
L157:
	;
	v1147 = int64(-5)
	goto L159
L158:
	;
	v1147 = int64(5)
	goto L159
L159:
	;
	v1150 = base.I64_div_s(v1147+v1126, int64(-10))
	v1151 = v1150 + v1075
	v1155 = v1066 & int32(255)
	if v832 == v1155 {
		goto L160
	} else {
		goto L161
	}
L160:
	;
	v1157 = v1151 + int64(-25165824)
	goto L162
L161:
	;
	v1157 = v1151
	goto L162
L162:
	;
	if v837 == v1155 {
		goto L163
	} else {
		goto L164
	}
L163:
	;
	v1161 = v1157 + int64(-25165824)
	goto L165
L164:
	;
	v1161 = v1157
	goto L165
L165:
	;
	v1165 = v1069 & int32(255)
	if v842 == v1165 {
		goto L166
	} else {
		goto L167
	}
L166:
	;
	v1167 = v1161 + int64(-25165824)
	goto L168
L167:
	;
	v1167 = v1161
	goto L168
L168:
	;
	if v847 == v1165 {
		goto L169
	} else {
		goto L170
	}
L169:
	;
	v1171 = v1167 + int64(-25165824)
	goto L171
L170:
	;
	v1171 = v1167
	goto L171
L171:
	;
	if v1066 != 0 {
		goto L172
	} else {
		goto L173
	}
L172:
	;
	v1174 = v1171
	goto L174
L173:
	;
	v1174 = v1171 + int64(-25165824)
	goto L174
L174:
	;
	if v1069 != 0 {
		goto L175
	} else {
		goto L176
	}
L175:
	;
	v1177 = v1174
	goto L177
L176:
	;
	v1177 = v1174 + int64(-25165824)
	goto L177
L177:
	;
	v1178 = base.B2i32(v1177 < v924)
	if v1177 < v924 {
		goto L178
	} else {
		goto L179
	}
L178:
	;
	v1179 = v1177
	goto L180
L179:
	;
	v1179 = v924
	goto L180
L180:
	;
	if v1177 < v924 {
		goto L181
	} else {
		goto L182
	}
L181:
	;
	v1180 = v1066
	goto L183
L182:
	;
	v1180 = v923
	goto L183
L183:
	;
	if v1177 < v924 {
		goto L184
	} else {
		goto L185
	}
L184:
	;
	v1181 = v1069
	goto L186
L185:
	;
	v1181 = v918
	goto L186
L186:
	;
	if (base.B2i32(int32(24) < l6)|base.B2i32(v871 != int32(4)))&base.B2i32(base.Ui32(v906) < base.Ui32(int32(7))) != 0 {
		v906 = v906 + int32(1)
		v918 = v1181
		v923 = v1180
		v924 = v1179
		goto L132
	} else {
		goto L187
	}
L187:
	;
	goto L133
L188:
	;
	goto L131
L189:
	;
	if l6 < int32(25) {
		v1200 = v1181
		v1201 = v1180
		goto L129
	} else {
		goto L192
	}
L190:
	;
	if v1180|v1181 == int32(0) {
		goto L188
	} else {
		goto L191
	}
L191:
	;
	goto L189
L192:
	;
	v1195 = v871 + int32(1)
	if base.Ui32(v855) <= base.Ui32(v1195) {
		v1200 = v1181
		v1201 = v1180
		goto L129
	} else {
		goto L193
	}
L193:
	;
	v871 = v1195
	v878 = v1181
	v883 = v1180
	v884 = v1179
	goto L130
}
func F_GetColorPalette(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v264 int32
	_ = v264
	var v267 int32
	_ = v267
	var v270 int32
	_ = v270
	var v271 int32
	_ = v271
	var v274 int32
	_ = v274
	var v278 int32
	_ = v278
	var v279 int32
	_ = v279
	var v283 int32
	_ = v283
	var v284 int32
	_ = v284
	var v291 int32
	_ = v291
	var v292 int32
	_ = v292
	var v298 int32
	_ = v298
	var v304 int32
	_ = v304
	var v311 int32
	_ = v311
	var v312 int32
	_ = v312
	var v313 int32
	_ = v313
	var v318 int32
	_ = v318
	var v331 int32
	_ = v331
	var v338 int32
	_ = v338
	var v339 int32
	_ = v339
	var v340 int32
	_ = v340
	var v343 int32
	_ = v343
	var v352 int32
	_ = v352
	var v353 int32
	_ = v353
	var v366 int32
	_ = v366
	var v377 int32
	_ = v377
	var v378 int32
	_ = v378
	var v388 int32
	_ = v388
	var v390 int32
	_ = v390
	var v395 int32
	_ = v395
	var v400 int32
	_ = v400
	var v411 int32
	_ = v411
	var v415 int32
	_ = v415
	var v416 int32
	_ = v416
	var v423 int32
	_ = v423
	var v427 int32
	_ = v427
	var v428 int32
	_ = v428
	var v434 int32
	_ = v434
	var v438 int32
	_ = v438
	var v441 int32
	_ = v441
	var v449 int32
	_ = v449
	var v453 int32
	_ = v453
	var v457 int32
	_ = v457
	var v461 int32
	_ = v461
	var v470 int32
	_ = v470
	v3 = int32(0)
	v13 = m.G0
	v15 = v13 - int32(_a_F_GetColorPalette_0)
	m.G0 = v15
	goto L3
L1:
	;
	goto L17
L3:
	;
	base.MemoryFill(m, v15+int32(_a_F_GetColorPalette_1), v3, int32(1024))
	goto L1
L15:
	;
	v264 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	if v264 < int32(1) {
		v400 = v3
		goto L31
	} else {
		goto L32
	}
L17:
	;
	base.MemoryFill(m, v15, int32(0), int32(_a_F_GetColorPalette_1))
	goto L15
L29:
	;
	m.G0 = v15 + int32(_a_F_GetColorPalette_0)
	return v470
L30:
	;
	v470 = int32(257)
	goto L29
L31:
	;
	if l1 == int32(0) {
		v470 = v400
		goto L29
	} else {
		goto L50
	}
L32:
	;
	v267 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	if v267 < int32(1) {
		v400 = v3
		goto L31
	} else {
		goto L33
	}
L33:
	;
	v270 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	v271 = *(*int32)(unsafe.Add(mBase, uint32(v270)))
	v274 = int32(0)
	v278 = v271 ^ int32(-1)
	v279 = v274
	v283 = v270
	v284 = v274
	goto L34
L34:
	;
	v291 = v278
	v292 = v279
	v298 = int32(0)
	goto L36
L35:
	;
	v400 = v378
	goto L31
L36:
	;
	v304 = *(*int32)(unsafe.Add(mBase, uint32(v283+v298<<(uint(int32(2))%32))))
	if v304 == v291 {
		v377 = v291
		v378 = v292
		goto L38
	} else {
		goto L39
	}
L37:
	;
	v390 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
	v395 = v284 + int32(1)
	if v395 != v264 {
		v278 = v377
		v279 = v378
		v283 = v283 + v390<<(uint(int32(2))%32)
		v284 = v395
		goto L34
	} else {
		goto L49
	}
L38:
	;
	v388 = v298 + int32(1)
	if v388 != v267 {
		v291 = v377
		v292 = v378
		v298 = v388
		goto L36
	} else {
		goto L48
	}
L39:
	;
	v311 = int32(base.Ui32(v304*int32(506832829)) >> (uint(int32(22)) % 32))
	v312 = v15 + int32(_a_F_GetColorPalette_1) + v311
	v313 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v312))))
	if v313 == int32(0) {
		v343 = v311
		v352 = v312
		goto L41
	} else {
		goto L42
	}
L40:
	;
	v377 = v304
	v378 = v366
	goto L38
L41:
	;
	v353 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v352))) = uint8(v353)
	*(*int32)(unsafe.Add(mBase, uint32(v15+v343<<(uint(int32(2))%32)))) = v304
	if int32(255) < v292 {
		goto L30
	} else {
		goto L47
	}
L42:
	;
	v318 = v311
	goto L43
L43:
	;
	v331 = *(*int32)(unsafe.Add(mBase, uint32(v15+v318<<(uint(int32(2))%32))))
	if v331 == v304 {
		v366 = v292
		goto L40
	} else {
		goto L45
	}
L44:
	;
	v343 = v338
	v352 = v339
	goto L41
L45:
	;
	v338 = (v318 + int32(1)) & int32(1023)
	v339 = v15 + int32(_a_F_GetColorPalette_1) + v338
	v340 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v339))))
	if v340 != 0 {
		v318 = v338
		goto L43
	} else {
		goto L46
	}
L46:
	;
	goto L44
L47:
	;
	v366 = v292 + int32(1)
	goto L40
L48:
	;
	goto L37
L49:
	;
	goto L35
L50:
	;
	v411 = int32(0)
	v415 = v15
	v416 = v411
	v423 = v411
	goto L51
L51:
	;
	v427 = v15 + int32(_a_F_GetColorPalette_1) + v423
	v428 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v427))))
	if v428 == int32(0) {
		v438 = v416
		goto L53
	} else {
		goto L54
	}
L52:
	;
	v461 = m.G2
	F___qsort_r(m, l1, v453, int32(4), int32(344), v461+int32(328))
	mBase = m.M
	goto L58
L53:
	;
	v441 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v427+int32(1)))))
	if v441 == int32(0) {
		v453 = v438
		goto L55
	} else {
		goto L56
	}
L54:
	;
	v434 = *(*int32)(unsafe.Add(mBase, uint32(v415)))
	*(*int32)(unsafe.Add(mBase, uint32(l1+v416<<(uint(int32(2))%32)))) = v434
	v438 = v416 + int32(1)
	goto L53
L55:
	;
	v457 = v423 + int32(2)
	if v457 != int32(1024) {
		v415 = v415 + int32(8)
		v416 = v453
		v423 = v457
		goto L51
	} else {
		goto L57
	}
L56:
	;
	v449 = *(*int32)(unsafe.Add(mBase, uint32(v415+int32(4))))
	*(*int32)(unsafe.Add(mBase, uint32(l1+v438<<(uint(int32(2))%32)))) = v449
	v453 = v438 + int32(1)
	goto L55
L57:
	;
	goto L52
L58:
	;
	v470 = v453
	goto L29
}
func F_GetCombinedHistogramEntropy(m *base.Module, l0 int32, l1 int32, l2 int64, l3 int32, l4 int32) int32 {
	mBase := m.M
	_ = mBase
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v33 int32
	_ = v33
	var v37 int32
	_ = v37
	var v41 int32
	_ = v41
	var v45 int32
	_ = v45
	var v49 int32
	_ = v49
	var v72 int32
	_ = v72
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v82 int32
	_ = v82
	var v86 int32
	_ = v86
	var v88 int32
	_ = v88
	var v90 int32
	_ = v90
	var v98 int32
	_ = v98
	var v106 int64
	_ = v106
	var v108 int64
	_ = v108
	var v109 int32
	_ = v109
	var v111 int32
	_ = v111
	var v113 int32
	_ = v113
	var v118 int32
	_ = v118
	var v119 int32
	_ = v119
	var v120 int32
	_ = v120
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v124 int32
	_ = v124
	var v127 int32
	_ = v127
	var v128 int32
	_ = v128
	var v130 int32
	_ = v130
	var v142 int64
	_ = v142
	var v145 int64
	_ = v145
	var v146 int64
	_ = v146
	var v149 int64
	_ = v149
	var v152 int64
	_ = v152
	var v154 int64
	_ = v154
	var v155 int64
	_ = v155
	var v158 int32
	_ = v158
	var v161 int32
	_ = v161
	var v170 int64
	_ = v170
	var v173 int64
	_ = v173
	var v176 int64
	_ = v176
	var v178 int64
	_ = v178
	var v179 int64
	_ = v179
	var v181 int32
	_ = v181
	var v184 int32
	_ = v184
	var v188 int32
	_ = v188
	var v192 int32
	_ = v192
	var v196 int32
	_ = v196
	var v200 int32
	_ = v200
	var v212 int64
	_ = v212
	var v219 int64
	_ = v219
	var v220 int64
	_ = v220
	var v227 int32
	_ = v227
	var v229 int32
	_ = v229
	var v276 int32
	_ = v276
	v25 = m.G0
	v27 = v25 - int32(48)
	m.G0 = v27
	if l2 < int64(1) {
		goto L2
	} else {
		goto L3
	}
L1:
	;
	m.G0 = v27 + int32(48)
	return v276
L2:
	;
	v276 = int32(0)
	goto L1
L3:
	;
	*(*int64)(unsafe.Add(mBase, uint32(l3))) = int64(0)
	v33 = int32(4)
	v37 = int32(1028)
	v41 = int32(2052)
	v45 = int32(3076)
	v49 = int32(3304)
	v72 = int32(3264)
	v73 = int32(3240)
	v74 = int32(0)
	goto L4
L4:
	;
	v82 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0+v73))))
	if v82 == int32(_a_F_GetCombinedHistogramEntropy_0) {
		v88 = int32(1)
		goto L6
	} else {
		goto L7
	}
L6:
	;
	v90 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v49+v74))))
	if v88 == int32(0) {
		goto L10
	} else {
		goto L11
	}
L7:
	;
	v86 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l1+v73))))
	v88 = base.B2i32(v82 != v86)
	goto L6
L8:
	;
	*(*int64)(unsafe.Add(mBase, uint32(l4+v72+int32(-3264)))) = v212
	v219 = *(*int64)(unsafe.Add(mBase, uint32(l3)))
	v220 = v219 + v212
	*(*int64)(unsafe.Add(mBase, uint32(l3))) = v220
	if base.Ui64(l2) <= base.Ui64(v220) {
		goto L2
	} else {
		goto L40
	}
L9:
	;
	v109 = int32(256)
	switch v74 {
	default:
		goto L20
	case 1:
		v122 = l0 + v33
		v123 = l1 + v33
		v124 = v109
		goto L16
	case 2:
		goto L19
	case 3:
		goto L18
	case 4:
		goto L17
	}
L10:
	;
	if v90&int32(255) == int32(0) {
		goto L14
	} else {
		goto L15
	}
L11:
	;
	if v90&int32(255) == int32(0) {
		goto L10
	} else {
		goto L12
	}
L12:
	;
	v98 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v49+v74))))
	if v98&int32(255) != 0 {
		goto L9
	} else {
		goto L13
	}
L13:
	;
	goto L10
L14:
	;
	v108 = *(*int64)(unsafe.Add(mBase, uint32(l1+v72)))
	v212 = v108
	goto L8
L15:
	;
	v106 = *(*int64)(unsafe.Add(mBase, uint32(l0+v72)))
	v212 = v106
	goto L8
L16:
	;
	v127 = m.G115
	v128 = *(*int32)(unsafe.Add(mBase, uint32(v127)))
	m.T0[v128].(func(*base.Module, int32, int32, int32, int32, int32))(m, v122, v123, v124, v27, v27+int32(24))
	mBase = m.M
	v130 = *(*int32)(unsafe.Add(mBase, uint32(v27)+12))
	if v130 <= int32(4) {
		goto L26
	} else {
		goto L27
	}
L17:
	;
	v122 = l0 + v45
	v123 = l1 + v45
	v124 = int32(40)
	goto L16
L18:
	;
	v122 = l0 + v41
	v123 = l1 + v41
	v124 = v109
	goto L16
L19:
	;
	v122 = l0 + v37
	v123 = l1 + v37
	v124 = v109
	goto L16
L20:
	;
	v111 = *(*int32)(unsafe.Add(mBase, uint32(l1)+3236))
	v113 = int32(280)
	if int32(0) < v111 {
		goto L21
	} else {
		goto L22
	}
L21:
	;
	v118 = int32(1)<<(uint(v111)%32) + v113
	goto L23
L22:
	;
	v118 = v113
	goto L23
L23:
	;
	v119 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	v120 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v122 = v120
	v123 = v119
	v124 = v118
	goto L16
L24:
	;
	v181 = *(*int32)(unsafe.Add(mBase, uint32(v27)+36))
	v184 = *(*int32)(unsafe.Add(mBase, uint32(v27)+24))
	v188 = *(*int32)(unsafe.Add(mBase, uint32(v27)+28))
	v192 = *(*int32)(unsafe.Add(mBase, uint32(v27)+44))
	v196 = *(*int32)(unsafe.Add(mBase, uint32(v27)+32))
	v200 = *(*int32)(unsafe.Add(mBase, uint32(v27)+40))
	v212 = v179 + base.I64_extend_i32_u(v181*int32(240)+v184*int32(1600)+v188*int32(2640)+v192*int32(720)+v196*int32(1840)+v200*int32(3360))<<(uint(int64(13))%64) + int64(401814323)
	goto L8
L25:
	;
	v155 = *(*int64)(unsafe.Add(mBase, uint32(v27)))
	v158 = *(*int32)(unsafe.Add(mBase, uint32(v27)+8))
	v161 = *(*int32)(unsafe.Add(mBase, uint32(v27)+16))
	v170 = v154*base.I64_extend_i32_u(v158<<(uint(int32(1))%32)-v161)<<(uint(int64(23))%64) + v155*(int64(1000)-v154)
	if v170 < int64(0) {
		goto L34
	} else {
		goto L35
	}
L26:
	;
	if v130 < int32(2) {
		v179 = int64(0)
		goto L24
	} else {
		goto L28
	}
L27:
	;
	v154 = int64(627)
	goto L25
L28:
	;
	switch v130 + int32(-2) {
	case 0:
		goto L30
	case 1:
		v154 = int64(950)
		goto L25
	default:
		goto L29
	}
L29:
	;
	v154 = int64(700)
	goto L25
L30:
	;
	v142 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v27)+8)))
	v145 = *(*int64)(unsafe.Add(mBase, uint32(v27)))
	v146 = v142*int64(830472192) + v145
	if v146 < int64(0) {
		goto L31
	} else {
		goto L32
	}
L31:
	;
	v149 = int64(-50)
	goto L33
L32:
	;
	v149 = int64(50)
	goto L33
L33:
	;
	v152 = base.I64_div_s(v149+v146, int64(100))
	v179 = v152
	goto L24
L34:
	;
	v173 = int64(-500)
	goto L36
L35:
	;
	v173 = int64(500)
	goto L36
L36:
	;
	v176 = base.I64_div_s(v173+v170, int64(1000))
	if base.Ui64(v176) < base.Ui64(v155) {
		goto L37
	} else {
		goto L38
	}
L37:
	;
	v178 = v155
	goto L39
L38:
	;
	v178 = v176
	goto L39
L39:
	;
	v179 = v178
	goto L24
L40:
	;
	v227 = int32(1)
	v229 = v74 + v227
	if v229 == int32(5) {
		v276 = v227
		goto L1
	} else {
		goto L41
	}
L41:
	;
	v72 = v72 + int32(8)
	v73 = v73 + int32(2)
	v74 = v229
	goto L4
}
func F_GetHuffBitLengthsAndCodes(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v19 int32
	_ = v19
	var v25 int32
	_ = v25
	var v26 int64
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v41 int32
	_ = v41
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v55 int32
	_ = v55
	var v60 int32
	_ = v60
	var v65 int64
	_ = v65
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v79 int64
	_ = v79
	var v80 int32
	_ = v80
	var v81 int64
	_ = v81
	var v90 int32
	_ = v90
	var v92 int32
	_ = v92
	var v96 int32
	_ = v96
	var v116 int32
	_ = v116
	var v122 int32
	_ = v122
	var v128 int32
	_ = v128
	var v134 int32
	_ = v134
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v145 int32
	_ = v145
	var v146 int32
	_ = v146
	var v151 int32
	_ = v151
	var v154 int32
	_ = v154
	var v158 int32
	_ = v158
	var v166 int32
	_ = v166
	var v167 int32
	_ = v167
	var v168 int32
	_ = v168
	var v172 int32
	_ = v172
	var v173 int32
	_ = v173
	var v182 int32
	_ = v182
	var v183 int32
	_ = v183
	var v189 int32
	_ = v189
	var v192 int32
	_ = v192
	var v195 int32
	_ = v195
	var v197 int32
	_ = v197
	var v198 int32
	_ = v198
	var v201 int32
	_ = v201
	var v205 int32
	_ = v205
	var v214 int32
	_ = v214
	var v215 int32
	_ = v215
	var v219 int32
	_ = v219
	var v220 int32
	_ = v220
	var v225 int32
	_ = v225
	var v228 int32
	_ = v228
	var v230 int32
	_ = v230
	var v239 int32
	_ = v239
	var v240 int32
	_ = v240
	var v245 int32
	_ = v245
	var v246 int64
	_ = v246
	var v253 int64
	_ = v253
	var v254 int32
	_ = v254
	var v255 int64
	_ = v255
	var v265 int32
	_ = v265
	var v267 int32
	_ = v267
	var v271 int64
	_ = v271
	var v272 int32
	_ = v272
	var v279 int64
	_ = v279
	var v280 int32
	_ = v280
	var v281 int64
	_ = v281
	var v291 int32
	_ = v291
	var v293 int32
	_ = v293
	var v299 int32
	_ = v299
	var v300 int32
	_ = v300
	var v310 int32
	_ = v310
	var v316 int32
	_ = v316
	var v317 int32
	_ = v317
	var v318 int32
	_ = v318
	var v320 int32
	_ = v320
	var v321 int32
	_ = v321
	var v322 int32
	_ = v322
	var v324 int32
	_ = v324
	var v353 int32
	_ = v353
	var v356 int32
	_ = v356
	var v364 int32
	_ = v364
	var v370 int32
	_ = v370
	var v372 int32
	_ = v372
	var v375 int32
	_ = v375
	var v376 int32
	_ = v376
	var v378 int32
	_ = v378
	var v387 int32
	_ = v387
	var v389 int32
	_ = v389
	var v400 int32
	_ = v400
	var v422 int32
	_ = v422
	var v425 int32
	_ = v425
	var v426 int32
	_ = v426
	var v434 int32
	_ = v434
	var v435 int32
	_ = v435
	var v470 int32
	_ = v470
	var v471 int32
	_ = v471
	var v476 int64
	_ = v476
	var v479 int32
	_ = v479
	var v480 int32
	_ = v480
	var v492 int32
	_ = v492
	v15 = int32(1)
	v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if v16 < v15 {
		v96 = int32(0)
		v116 = F_calloc(m, base.I32_wrap_i64(int64(0)), int32(3))
		mBase = m.M
		if v116 != 0 {
			v239 = v96
			v240 = v116
			v245 = int32(0)
			v246 = int64(1)
			v253 = base.I64_div_u_s(int64(2147418112), v246)
			v254 = int32(0)
			v255 = base.I64_extend_i32_u(v239)
			if base.Ui64(int64(4294967295)) < base.Ui64(v255*v246) {
				v267 = v254
			} else {
				if base.Ui64(v253) < base.Ui64(v255) {
					v267 = v254
				} else {
					v265 = F_malloc(m, base.I32_wrap_i64(v246)*v239)
					mBase = m.M
					v267 = v265
				}
			}
			v271 = base.I64_extend_i32_u(v239) * int64(3)
			v272 = int32(16)
			if v271 == int64(0) {
				v291 = F_malloc(m, base.I32_wrap_i64(v271)*v272)
				mBase = m.M
				v293 = v291
			} else {
				v279 = base.I64_div_u_s(int64(2147418112), v271)
				v280 = int32(0)
				v281 = base.I64_extend_i32_u(v272)
				if base.Ui64(int64(4294967295)) < base.Ui64(v281*v271) {
					v293 = v280
				} else {
					if base.Ui64(v279) < base.Ui64(v281) {
						v293 = v280
					} else {
						v291 = F_malloc(m, base.I32_wrap_i64(v271)*v272)
						mBase = m.M
						v293 = v291
					}
				}
			}
			if v267 == int32(0) {
				v356 = v15
				v364 = v245
			} else {
				if v293 == int32(0) {
					v356 = v15
					v364 = v245
				} else {
					v299 = int32(1)
					v300 = int32(0)
					if v16 < v299 {
						v356 = v300
						v364 = v299
					} else {
						v310 = l1
						v316 = v16
						v317 = int32(0)
						for {
							v318 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
							v320 = *(*int32)(unsafe.Add(mBase, uint32(v318+v317)))
							v321 = *(*int32)(unsafe.Add(mBase, uint32(v320)))
							v322 = int32(15)
							F_VP8LCreateHuffmanTree(m, v321, v322, v267, v293, v310)
							mBase = m.M
							v324 = int32(4)
							F_VP8LCreateHuffmanTree(m, v320+v324, v322, v267, v293, v310+int32(12))
							mBase = m.M
							F_VP8LCreateHuffmanTree(m, v320+int32(1028), v322, v267, v293, v310+int32(24))
							mBase = m.M
							F_VP8LCreateHuffmanTree(m, v320+int32(2052), v322, v267, v293, v310+int32(36))
							mBase = m.M
							F_VP8LCreateHuffmanTree(m, v320+int32(3076), v322, v267, v293, v310+int32(48))
							mBase = m.M
							v353 = v316 + int32(-1)
							if v353 != 0 {
								v310 = v310 + int32(60)
								v316 = v353
								v317 = v317 + v324
								continue
							} else {
								break
							}
							break
						}
						v356 = v300
						v364 = v299
					}
				}
			}
			v370 = v356
			v372 = v267
			v375 = v293
			v376 = v240
			v378 = v364
		} else {
			v122 = v15
			v128 = v96
			v134 = int32(0)
			v370 = v122
			v372 = v134
			v375 = v134
			v376 = v128
			v378 = v134
		}
	} else {
		v19 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
		v25 = v19
		v26 = int64(0)
		v27 = l1
		v28 = v16
		for {
			*(*int32)(unsafe.Add(mBase, uint32(v27+int32(48)))) = int32(40)
			v41 = int32(256)
			*(*int32)(unsafe.Add(mBase, uint32(v27+int32(36)))) = v41
			*(*int32)(unsafe.Add(mBase, uint32(v27+int32(24)))) = v41
			*(*int32)(unsafe.Add(mBase, uint32(v27+int32(12)))) = v41
			v52 = *(*int32)(unsafe.Add(mBase, uint32(v25)))
			v53 = *(*int32)(unsafe.Add(mBase, uint32(v52)+3236))
			v55 = int32(280)
			if int32(0) < v53 {
				v60 = int32(1)<<(uint(v53)%32) + v55
			} else {
				v60 = v55
			}
			*(*int32)(unsafe.Add(mBase, uint32(v27))) = v60
			v65 = v26 + base.I64_extend_i32_s(v60) + int64(808)
			v71 = v28 + int32(-1)
			if v71 != 0 {
				v25 = v25 + int32(4)
				v26 = v65
				v27 = v27 + int32(60)
				v28 = v71
				continue
			} else {
				break
			}
			break
		}
		v72 = int32(3)
		if v65 == int64(0) {
			v90 = F_calloc(m, base.I32_wrap_i64(v65), v72)
			mBase = m.M
			v92 = v90
		} else {
			v79 = base.I64_div_u_s(int64(2147418112), v65)
			v80 = int32(0)
			v81 = base.I64_extend_i32_u(v72)
			if base.Ui64(int64(4294967295)) < base.Ui64(v81*v65) {
				v92 = v80
			} else {
				if base.Ui64(v79) < base.Ui64(v81) {
					v92 = v80
				} else {
					v90 = F_calloc(m, base.I32_wrap_i64(v65), v72)
					mBase = m.M
					v92 = v90
				}
			}
		}
		if v92 != 0 {
			if int32(1) <= v16 {
				v141 = v16 * int32(5)
				v142 = int32(1)
				if v142 < v141 {
					v145 = v141
				} else {
					v145 = v142
				}
				v146 = int32(1)
				v151 = v92 + base.I32_wrap_i64(v65)<<(uint(v146)%32)
				if int32(2) <= v141 {
					v158 = int32(0)
					v166 = l1
					v167 = v151
					v168 = v158
					v172 = v158
					v173 = v92
					for {
						*(*int32)(unsafe.Add(mBase, uint32(v166+int32(8)))) = v173
						*(*int32)(unsafe.Add(mBase, uint32(v166+int32(4)))) = v167
						v182 = *(*int32)(unsafe.Add(mBase, uint32(v166)))
						v183 = v167 + v182
						*(*int32)(unsafe.Add(mBase, uint32(v166+int32(16)))) = v183
						v189 = v173 + v182<<(uint(int32(1))%32)
						*(*int32)(unsafe.Add(mBase, uint32(v166+int32(20)))) = v189
						if v182 < v168 {
							v192 = v168
						} else {
							v192 = v182
						}
						v195 = *(*int32)(unsafe.Add(mBase, uint32(v166+int32(12))))
						if v195 < v192 {
							v197 = v192
						} else {
							v197 = v195
						}
						v198 = v183 + v195
						v201 = v189 + v195<<(uint(int32(1))%32)
						v205 = v172 + int32(2)
						if v145&int32(2147483646) != v205 {
							v166 = v166 + int32(24)
							v167 = v198
							v168 = v197
							v172 = v205
							v173 = v201
							continue
						} else {
							break
						}
						break
					}
					v214 = v198
					v215 = v197
					v219 = v205
					v220 = v201
				} else {
					v154 = int32(0)
					v214 = v151
					v215 = v154
					v219 = v154
					v220 = v92
				}
				if v145&v146 == int32(0) {
					v239 = v215
					v240 = v92
				} else {
					v225 = l1 + v219*int32(12)
					*(*int32)(unsafe.Add(mBase, uint32(v225)+8)) = v220
					*(*int32)(unsafe.Add(mBase, uint32(v225)+4)) = v214
					v228 = *(*int32)(unsafe.Add(mBase, uint32(v225)))
					if v228 < v215 {
						v230 = v215
					} else {
						v230 = v228
					}
					v239 = v230
					v240 = v92
				}
			} else {
				v239 = int32(0)
				v240 = v92
			}
			v245 = int32(0)
			v246 = int64(1)
			v253 = base.I64_div_u_s(int64(2147418112), v246)
			v254 = int32(0)
			v255 = base.I64_extend_i32_u(v239)
			if base.Ui64(int64(4294967295)) < base.Ui64(v255*v246) {
				v267 = v254
			} else {
				if base.Ui64(v253) < base.Ui64(v255) {
					v267 = v254
				} else {
					v265 = F_malloc(m, base.I32_wrap_i64(v246)*v239)
					mBase = m.M
					v267 = v265
				}
			}
			v271 = base.I64_extend_i32_u(v239) * int64(3)
			v272 = int32(16)
			if v271 == int64(0) {
				v291 = F_malloc(m, base.I32_wrap_i64(v271)*v272)
				mBase = m.M
				v293 = v291
			} else {
				v279 = base.I64_div_u_s(int64(2147418112), v271)
				v280 = int32(0)
				v281 = base.I64_extend_i32_u(v272)
				if base.Ui64(int64(4294967295)) < base.Ui64(v281*v271) {
					v293 = v280
				} else {
					if base.Ui64(v279) < base.Ui64(v281) {
						v293 = v280
					} else {
						v291 = F_malloc(m, base.I32_wrap_i64(v271)*v272)
						mBase = m.M
						v293 = v291
					}
				}
			}
			if v267 == int32(0) {
				v356 = v15
				v364 = v245
			} else {
				if v293 == int32(0) {
					v356 = v15
					v364 = v245
				} else {
					v299 = int32(1)
					v300 = int32(0)
					if v16 < v299 {
						v356 = v300
						v364 = v299
					} else {
						v310 = l1
						v316 = v16
						v317 = int32(0)
						for {
							v318 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
							v320 = *(*int32)(unsafe.Add(mBase, uint32(v318+v317)))
							v321 = *(*int32)(unsafe.Add(mBase, uint32(v320)))
							v322 = int32(15)
							F_VP8LCreateHuffmanTree(m, v321, v322, v267, v293, v310)
							mBase = m.M
							v324 = int32(4)
							F_VP8LCreateHuffmanTree(m, v320+v324, v322, v267, v293, v310+int32(12))
							mBase = m.M
							F_VP8LCreateHuffmanTree(m, v320+int32(1028), v322, v267, v293, v310+int32(24))
							mBase = m.M
							F_VP8LCreateHuffmanTree(m, v320+int32(2052), v322, v267, v293, v310+int32(36))
							mBase = m.M
							F_VP8LCreateHuffmanTree(m, v320+int32(3076), v322, v267, v293, v310+int32(48))
							mBase = m.M
							v353 = v316 + int32(-1)
							if v353 != 0 {
								v310 = v310 + int32(60)
								v316 = v353
								v317 = v317 + v324
								continue
							} else {
								break
							}
							break
						}
						v356 = v300
						v364 = v299
					}
				}
			}
			v370 = v356
			v372 = v267
			v375 = v293
			v376 = v240
			v378 = v364
		} else {
			v122 = int32(1)
			v128 = int32(0)
			v134 = int32(0)
			v370 = v122
			v372 = v134
			v375 = v134
			v376 = v128
			v378 = v134
		}
	}
	F_free(m, v375)
	mBase = m.M
	F_free(m, v372)
	mBase = m.M
	if v370 == int32(0) {
	} else {
		F_free(m, v376)
		mBase = m.M
		v387 = int32(0)
		v389 = v16 * int32(60)
		if base.Ui32(v389) < base.Ui32(int32(33)) {
			if v389 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(l1))) = uint8(v387)
				v400 = l1 + v389
				*(*uint8)(unsafe.Add(mBase, uint32(v400+int32(-1)))) = uint8(v387)
				if base.Ui32(v389) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(l1)+2)) = uint8(v387)
					*(*uint8)(unsafe.Add(mBase, uint32(l1)+1)) = uint8(v387)
					*(*uint8)(unsafe.Add(mBase, uint32(v400+int32(-3)))) = uint8(v387)
					*(*uint8)(unsafe.Add(mBase, uint32(v400+int32(-2)))) = uint8(v387)
					if base.Ui32(v389) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(l1)+3)) = uint8(v387)
						*(*uint8)(unsafe.Add(mBase, uint32(v400+int32(-4)))) = uint8(v387)
						if base.Ui32(v389) < base.Ui32(int32(9)) {
						} else {
							v422 = int32(0)
							v425 = (v422 - l1) & int32(3)
							v426 = l1 + v425
							*(*int32)(unsafe.Add(mBase, uint32(v426))) = v422
							v434 = (v389 - v425) & int32(60)
							v435 = v426 + v434
							*(*int32)(unsafe.Add(mBase, uint32(v435+int32(-4)))) = v422
							if base.Ui32(v434) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v426)+8)) = v422
								*(*int32)(unsafe.Add(mBase, uint32(v426)+4)) = v422
								*(*int32)(unsafe.Add(mBase, uint32(v435+int32(-8)))) = v422
								*(*int32)(unsafe.Add(mBase, uint32(v435+int32(-12)))) = v422
								if base.Ui32(v434) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v426)+24)) = v422
									*(*int32)(unsafe.Add(mBase, uint32(v426)+20)) = v422
									*(*int32)(unsafe.Add(mBase, uint32(v426)+16)) = v422
									*(*int32)(unsafe.Add(mBase, uint32(v426)+12)) = v422
									*(*int32)(unsafe.Add(mBase, uint32(v435+int32(-16)))) = v422
									*(*int32)(unsafe.Add(mBase, uint32(v435+int32(-20)))) = v422
									*(*int32)(unsafe.Add(mBase, uint32(v435+int32(-24)))) = v422
									*(*int32)(unsafe.Add(mBase, uint32(v435+int32(-28)))) = v422
									v470 = v426&int32(4) | int32(24)
									v471 = v434 - v470
									if base.Ui32(v471) < base.Ui32(int32(32)) {
									} else {
										v476 = base.I64_extend_i32_u(v422) * int64(4294967297)
										v479 = v471
										v480 = v426 + v470
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v480)+24)) = v476
											*(*int64)(unsafe.Add(mBase, uint32(v480)+16)) = v476
											*(*int64)(unsafe.Add(mBase, uint32(v480)+8)) = v476
											*(*int64)(unsafe.Add(mBase, uint32(v480))) = v476
											v492 = v479 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v492) {
												v479 = v492
												v480 = v480 + int32(32)
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
			base.MemoryFill(m, l1, v387, v389)
		}
	}
	return v378
}
func F_GetMBSSIM(m *base.Module, l0 int32, l1 int32) float64 {
	mBase := m.M
	_ = mBase
	var v12 float64
	_ = v12
	var v13 int32
	_ = v13
	var v17 int32
	_ = v17
	var v20 int32
	_ = v20
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v24 float64
	_ = v24
	var v31 int32
	_ = v31
	var v32 float64
	_ = v32
	var v39 int32
	_ = v39
	var v40 float64
	_ = v40
	var v47 int32
	_ = v47
	var v48 float64
	_ = v48
	var v55 int32
	_ = v55
	var v56 float64
	_ = v56
	var v63 int32
	_ = v63
	var v64 float64
	_ = v64
	var v71 int32
	_ = v71
	var v72 float64
	_ = v72
	var v79 int32
	_ = v79
	var v80 float64
	_ = v80
	var v87 int32
	_ = v87
	var v88 float64
	_ = v88
	var v95 int32
	_ = v95
	var v96 float64
	_ = v96
	var v97 float64
	_ = v97
	var v99 int32
	_ = v99
	var v102 int32
	_ = v102
	var v103 int32
	_ = v103
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v107 int32
	_ = v107
	var v109 int32
	_ = v109
	var v113 float64
	_ = v113
	var v114 int32
	_ = v114
	var v118 int32
	_ = v118
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v123 int32
	_ = v123
	var v124 int32
	_ = v124
	var v125 float64
	_ = v125
	var v132 int32
	_ = v132
	var v133 float64
	_ = v133
	var v137 int32
	_ = v137
	var v140 int32
	_ = v140
	var v141 float64
	_ = v141
	var v148 int32
	_ = v148
	var v149 float64
	_ = v149
	var v153 int32
	_ = v153
	var v156 int32
	_ = v156
	var v157 float64
	_ = v157
	var v164 int32
	_ = v164
	var v165 float64
	_ = v165
	var v169 int32
	_ = v169
	var v172 int32
	_ = v172
	var v173 float64
	_ = v173
	var v180 int32
	_ = v180
	var v181 float64
	_ = v181
	var v185 int32
	_ = v185
	var v188 int32
	_ = v188
	var v189 float64
	_ = v189
	var v196 int32
	_ = v196
	var v197 float64
	_ = v197
	var v201 int32
	_ = v201
	var v204 int32
	_ = v204
	var v205 float64
	_ = v205
	var v212 int32
	_ = v212
	var v213 float64
	_ = v213
	var v214 float64
	_ = v214
	var v216 int32
	_ = v216
	v12 = float64(0)
	v13 = int32(3)
	for {
		v17 = int32(32)
		v20 = int32(16)
		v22 = m.G57
		v23 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
		v24 = m.T0[v23].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, l0, v17, l1, v17, int32(3), v13, v20, v20)
		mBase = m.M
		v31 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
		v32 = m.T0[v31].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, l0, v17, l1, v17, int32(4), v13, v20, v20)
		mBase = m.M
		v39 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
		v40 = m.T0[v39].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, l0, v17, l1, v17, int32(5), v13, v20, v20)
		mBase = m.M
		v47 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
		v48 = m.T0[v47].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, l0, v17, l1, v17, int32(6), v13, v20, v20)
		mBase = m.M
		v55 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
		v56 = m.T0[v55].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, l0, v17, l1, v17, int32(7), v13, v20, v20)
		mBase = m.M
		v63 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
		v64 = m.T0[v63].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, l0, v17, l1, v17, int32(8), v13, v20, v20)
		mBase = m.M
		v71 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
		v72 = m.T0[v71].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, l0, v17, l1, v17, int32(9), v13, v20, v20)
		mBase = m.M
		v79 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
		v80 = m.T0[v79].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, l0, v17, l1, v17, int32(10), v13, v20, v20)
		mBase = m.M
		v87 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
		v88 = m.T0[v87].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, l0, v17, l1, v17, int32(11), v13, v20, v20)
		mBase = m.M
		v95 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
		v96 = m.T0[v95].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, l0, v17, l1, v17, int32(12), v13, v20, v20)
		mBase = m.M
		v97 = base.F64_add(base.F64_add(base.F64_add(base.F64_add(base.F64_add(base.F64_add(base.F64_add(base.F64_add(base.F64_add(base.F64_add(v12, v24), v32), v40), v48), v56), v64), v72), v80), v88), v96)
		v99 = v13 + int32(1)
		if v99 != int32(13) {
			v12 = v97
			v13 = v99
			continue
		} else {
			break
		}
		break
	}
	v102 = int32(24)
	v103 = l1 + v102
	v105 = l0 + v102
	v106 = int32(16)
	v107 = l1 + v106
	v109 = l0 + v106
	v113 = v97
	v114 = int32(1)
	for {
		v118 = int32(32)
		v120 = int32(1)
		v121 = int32(8)
		v123 = m.G57
		v124 = *(*int32)(unsafe.Add(mBase, uint32(v123)))
		v125 = m.T0[v124].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, v109, v118, v107, v118, v114, v120, v121, v121)
		mBase = m.M
		v132 = *(*int32)(unsafe.Add(mBase, uint32(v123)))
		v133 = m.T0[v132].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, v105, v118, v103, v118, v114, v120, v121, v121)
		mBase = m.M
		v137 = int32(2)
		v140 = *(*int32)(unsafe.Add(mBase, uint32(v123)))
		v141 = m.T0[v140].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, v109, v118, v107, v118, v114, v137, v121, v121)
		mBase = m.M
		v148 = *(*int32)(unsafe.Add(mBase, uint32(v123)))
		v149 = m.T0[v148].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, v105, v118, v103, v118, v114, v137, v121, v121)
		mBase = m.M
		v153 = int32(3)
		v156 = *(*int32)(unsafe.Add(mBase, uint32(v123)))
		v157 = m.T0[v156].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, v109, v118, v107, v118, v114, v153, v121, v121)
		mBase = m.M
		v164 = *(*int32)(unsafe.Add(mBase, uint32(v123)))
		v165 = m.T0[v164].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, v105, v118, v103, v118, v114, v153, v121, v121)
		mBase = m.M
		v169 = int32(4)
		v172 = *(*int32)(unsafe.Add(mBase, uint32(v123)))
		v173 = m.T0[v172].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, v109, v118, v107, v118, v114, v169, v121, v121)
		mBase = m.M
		v180 = *(*int32)(unsafe.Add(mBase, uint32(v123)))
		v181 = m.T0[v180].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, v105, v118, v103, v118, v114, v169, v121, v121)
		mBase = m.M
		v185 = int32(5)
		v188 = *(*int32)(unsafe.Add(mBase, uint32(v123)))
		v189 = m.T0[v188].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, v109, v118, v107, v118, v114, v185, v121, v121)
		mBase = m.M
		v196 = *(*int32)(unsafe.Add(mBase, uint32(v123)))
		v197 = m.T0[v196].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, v105, v118, v103, v118, v114, v185, v121, v121)
		mBase = m.M
		v201 = int32(6)
		v204 = *(*int32)(unsafe.Add(mBase, uint32(v123)))
		v205 = m.T0[v204].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, v109, v118, v107, v118, v114, v201, v121, v121)
		mBase = m.M
		v212 = *(*int32)(unsafe.Add(mBase, uint32(v123)))
		v213 = m.T0[v212].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32) float64)(m, v105, v118, v103, v118, v114, v201, v121, v121)
		mBase = m.M
		v214 = base.F64_add(base.F64_add(base.F64_add(base.F64_add(base.F64_add(base.F64_add(base.F64_add(base.F64_add(base.F64_add(base.F64_add(base.F64_add(base.F64_add(v113, v125), v133), v141), v149), v157), v165), v173), v181), v189), v197), v205), v213)
		v216 = v114 + v120
		if v216 != int32(7) {
			v113 = v214
			v114 = v216
			continue
		} else {
			break
		}
		break
	}
	return v214
}
func F_GetResidual(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32, l10 int32, l11 int32, l12 int32) {
	mBase := m.M
	_ = mBase
	var v33 int32
	_ = v33
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v55 int32
	_ = v55
	var v56 int32
	_ = v56
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v62 int32
	_ = v62
	var v66 int32
	_ = v66
	var v68 int32
	_ = v68
	var v81 int32
	_ = v81
	var v85 int32
	_ = v85
	var v87 int32
	_ = v87
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v111 int32
	_ = v111
	var v112 int32
	_ = v112
	var v117 int32
	_ = v117
	var v124 int32
	_ = v124
	var v137 int32
	_ = v137
	var v142 int32
	_ = v142
	var v144 int32
	_ = v144
	var v147 int32
	_ = v147
	var v150 int32
	_ = v150
	var v163 int32
	_ = v163
	var v165 int32
	_ = v165
	var v169 int32
	_ = v169
	var v193 int32
	_ = v193
	var v194 int32
	_ = v194
	var v195 int32
	_ = v195
	var v196 int32
	_ = v196
	var v197 int32
	_ = v197
	var v202 int32
	_ = v202
	var v209 int32
	_ = v209
	var v219 int32
	_ = v219
	var v220 int32
	_ = v220
	var v222 int32
	_ = v222
	var v224 int32
	_ = v224
	var v229 int32
	_ = v229
	var v243 int32
	_ = v243
	var v268 int32
	_ = v268
	var v270 int32
	_ = v270
	var v273 int32
	_ = v273
	var v275 int32
	_ = v275
	var v280 int32
	_ = v280
	var v283 int32
	_ = v283
	var v290 int32
	_ = v290
	var v296 int32
	_ = v296
	var v297 int32
	_ = v297
	var v299 int32
	_ = v299
	var v304 int32
	_ = v304
	var v305 int32
	_ = v305
	var v311 int32
	_ = v311
	var v312 int32
	_ = v312
	var v325 int32
	_ = v325
	var v348 int32
	_ = v348
	var v350 int32
	_ = v350
	var v351 int32
	_ = v351
	var v353 int32
	_ = v353
	var v354 int32
	_ = v354
	var v361 int32
	_ = v361
	var v362 int32
	_ = v362
	var v364 int32
	_ = v364
	var v365 int32
	_ = v365
	var v366 int32
	_ = v366
	var v371 int32
	_ = v371
	var v374 int32
	_ = v374
	var v378 int32
	_ = v378
	var v380 int32
	_ = v380
	var v382 int32
	_ = v382
	var v384 int32
	_ = v384
	var v386 int32
	_ = v386
	var v387 int32
	_ = v387
	var v391 int32
	_ = v391
	var v392 int32
	_ = v392
	var v394 int32
	_ = v394
	var v396 int32
	_ = v396
	var v397 int32
	_ = v397
	var v398 int32
	_ = v398
	var v399 int32
	_ = v399
	var v401 int32
	_ = v401
	var v408 int32
	_ = v408
	var v411 int32
	_ = v411
	var v415 int32
	_ = v415
	var v417 int32
	_ = v417
	var v419 int32
	_ = v419
	var v421 int32
	_ = v421
	var v422 int32
	_ = v422
	var v425 int32
	_ = v425
	var v428 int32
	_ = v428
	var v429 int32
	_ = v429
	var v431 int32
	_ = v431
	var v433 int32
	_ = v433
	var v434 int32
	_ = v434
	var v435 int32
	_ = v435
	var v436 int32
	_ = v436
	var v438 int32
	_ = v438
	var v445 int32
	_ = v445
	var v447 int32
	_ = v447
	var v451 int32
	_ = v451
	var v456 int32
	_ = v456
	var v458 int32
	_ = v458
	var v460 int32
	_ = v460
	var v462 int32
	_ = v462
	var v463 int32
	_ = v463
	var v464 int32
	_ = v464
	var v466 int32
	_ = v466
	var v467 int32
	_ = v467
	var v468 int32
	_ = v468
	var v469 int32
	_ = v469
	var v471 int32
	_ = v471
	var v478 int32
	_ = v478
	var v483 int32
	_ = v483
	var v485 int32
	_ = v485
	var v487 int32
	_ = v487
	var v489 int32
	_ = v489
	var v490 int32
	_ = v490
	var v510 int32
	_ = v510
	var v518 int32
	_ = v518
	var v520 int32
	_ = v520
	var v540 int32
	_ = v540
	var v545 int32
	_ = v545
	var v550 int32
	_ = v550
	var v553 int32
	_ = v553
	var v560 int32
	_ = v560
	var v562 int32
	_ = v562
	var v591 int32
	_ = v591
	var v592 int32
	_ = v592
	var v594 int32
	_ = v594
	var v601 int32
	_ = v601
	var v611 int32
	_ = v611
	var v613 int32
	_ = v613
	var v617 int32
	_ = v617
	var v641 int32
	_ = v641
	var v642 int32
	_ = v642
	var v643 int32
	_ = v643
	var v644 int32
	_ = v644
	var v645 int32
	_ = v645
	var v650 int32
	_ = v650
	var v657 int32
	_ = v657
	var v667 int32
	_ = v667
	var v668 int32
	_ = v668
	var v670 int32
	_ = v670
	var v672 int32
	_ = v672
	var v677 int32
	_ = v677
	if l10 == int32(0) {
		v68 = m.G20
		if l7 <= l6 {
		} else {
			if l8 != 0 {
				v144 = int32(2)
				v147 = *(*int32)(unsafe.Add(mBase, uint32(v68+l5<<(uint(v144)%32))))
				v150 = l2 + l0<<(uint(v144)%32)
				if l5 == int32(0) {
					v163 = l6
					v165 = l6 << (uint(int32(2)) % 32)
					v169 = l12
					for {
						if v163 == int32(0) {
							v194 = *(*int32)(unsafe.Add(mBase, uint32(l2)))
							v195 = v194
						} else {
							v193 = m.T0[v147].(func(*base.Module, int32, int32) int32)(m, l3+v165+int32(-4), l2+v165)
							mBase = m.M
							v195 = v193
						}
						v196 = l3 + v165
						v197 = *(*int32)(unsafe.Add(mBase, uint32(v196)))
						v202 = v197 | int32(16711680) - v195&int32(-16711936)
						if base.Ui32(v197) <= base.Ui32(int32(16777215)) {
							*(*int32)(unsafe.Add(mBase, uint32(v196))) = v195 & int32(16777215)
							v219 = v202 & int32(-16777216)
							if v163 != 0 {
								v222 = v219
							} else {
								v220 = *(*int32)(unsafe.Add(mBase, uint32(l3)))
								*(*int32)(unsafe.Add(mBase, uint32(v150))) = v220
								v222 = v219
							}
						} else {
							v209 = int32(16711935)
							v222 = v202&int32(-16711936) | (v197|int32(_a_F_GetResidual_0)-v195&v209)&v209
						}
						*(*int32)(unsafe.Add(mBase, uint32(v169))) = v222
						v224 = int32(4)
						v229 = v163 + int32(1)
						if l7 != v229 {
							v163 = v229
							v165 = v165 + v224
							v169 = v169 + v224
							continue
						} else {
							break
						}
						break
					}
				} else {
					if l9 != int32(1) {
						if l1+int32(-1) == l8 {
							v611 = l6
							v613 = l6 << (uint(int32(2)) % 32)
							v617 = l12
							for {
								if v611 == int32(0) {
									v642 = *(*int32)(unsafe.Add(mBase, uint32(l2)))
									v643 = v642
								} else {
									v641 = m.T0[v147].(func(*base.Module, int32, int32) int32)(m, l3+v613+int32(-4), l2+v613)
									mBase = m.M
									v643 = v641
								}
								v644 = l3 + v613
								v645 = *(*int32)(unsafe.Add(mBase, uint32(v644)))
								v650 = v645 | int32(16711680) - v643&int32(-16711936)
								if base.Ui32(v645) <= base.Ui32(int32(16777215)) {
									*(*int32)(unsafe.Add(mBase, uint32(v644))) = v643 & int32(16777215)
									v667 = v650 & int32(-16777216)
									if v611 != 0 {
										v670 = v667
									} else {
										v668 = *(*int32)(unsafe.Add(mBase, uint32(l3)))
										*(*int32)(unsafe.Add(mBase, uint32(v150))) = v668
										v670 = v667
									}
								} else {
									v657 = int32(16711935)
									v670 = v650&int32(-16711936) | (v645|int32(_a_F_GetResidual_0)-v643&v657)&v657
								}
								*(*int32)(unsafe.Add(mBase, uint32(v617))) = v670
								v672 = int32(4)
								v677 = v611 + int32(1)
								if l7 != v677 {
									v611 = v677
									v613 = v613 + v672
									v617 = v617 + v672
									continue
								} else {
									break
								}
								break
							}
						} else {
							v243 = l6
							for {
								if v243 != 0 {
									v270 = v243 << (uint(int32(2)) % 32)
									v273 = m.T0[v147].(func(*base.Module, int32, int32) int32)(m, l3+int32(-4)+v270, l2+v270)
									mBase = m.M
									if v243 != l0+int32(-1) {
										v296 = l3 + v270
										v297 = *(*int32)(unsafe.Add(mBase, uint32(v296)))
										v299 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4+v243))))
										if base.Ui32(int32(3)) <= base.Ui32(v299) {
											v325 = l9
											for {
												v348 = int32(base.Ui32(v325) >> (uint(int32(1)) % 32))
												if v299 <= v325 {
													v325 = v348
													continue
												} else {
													break
												}
												break
											}
											v350 = int32(24)
											v351 = int32(base.Ui32(v297) >> (uint(v350) % 32))
											v353 = int32(base.Ui32(v273) >> (uint(v350) % 32))
											v354 = v351 - v353
											if base.Ui32(v297) < base.Ui32(int32(16777216)) {
												v386 = v354
												v387 = int32(0) - v325
											} else {
												if v351 != int32(255) {
													v361 = int32(255)
													v362 = v354 & v361
													v364 = int32(0) - v325
													v365 = v362 & v364
													v366 = v365 + v348
													v371 = v353 ^ v361
													v374 = v365 + v325
													if base.B2i32(base.Ui32(v351^v361) < base.Ui32(v371))-v362+v374 <= v362-v365 {
														if v371 < v374 {
															v382 = v366
														} else {
															v382 = v374
														}
														if base.Ui32(v362) <= base.Ui32(v371) {
															v384 = v382
														} else {
															v384 = v374
														}
														v386 = v384
														v387 = v364
													} else {
														if base.Ui32(v371) < base.Ui32(v365) {
															v378 = v365
														} else {
															v378 = v366
														}
														if base.Ui32(v362) <= base.Ui32(v371) {
															v380 = v365
														} else {
															v380 = v378
														}
														v386 = v380
														v387 = v364
													}
												} else {
													v386 = v354
													v387 = int32(0) - v325
												}
											}
											v391 = int32(8)
											v392 = int32(base.Ui32(v297) >> (uint(v391) % 32))
											v394 = int32(base.Ui32(v273) >> (uint(v391) % 32))
											v396 = int32(255)
											v397 = (v392 - v394) & v396
											v398 = v387 & v397
											v399 = v398 + v348
											v401 = int32(-1)
											v408 = (v394 ^ v401) & v396
											v411 = v398 + v325
											if base.B2i32(base.Ui32((v392^v401)&v396) < base.Ui32(v408))-v397+v411 <= v397-v398 {
												if v408 < v411 {
													v419 = v399
												} else {
													v419 = v411
												}
												if base.Ui32(v397) <= base.Ui32(v408) {
													v421 = v419
												} else {
													v421 = v411
												}
												v422 = v421
											} else {
												if base.Ui32(v408) < base.Ui32(v398) {
													v415 = v398
												} else {
													v415 = v399
												}
												if base.Ui32(v397) <= base.Ui32(v408) {
													v417 = v398
												} else {
													v417 = v415
												}
												v422 = v417
											}
											v425 = v422 + v394
											if l11 != 0 {
												v428 = v425 - v392
											} else {
												v428 = int32(0)
											}
											v429 = int32(base.Ui32(v297)>>(uint(int32(16))%32)) - v428
											v431 = int32(base.Ui32(v273) >> (uint(int32(16)) % 32))
											v433 = int32(255)
											v434 = (v429 - v431) & v433
											v435 = v434 & v387
											v436 = v435 + v348
											v438 = v435 + v325
											if l11 != 0 {
												v445 = (v425 ^ int32(-1)) & v433
											} else {
												v445 = v433
											}
											v447 = int32(255)
											v451 = (v445 - v431) & v447
											if v438-v434+base.B2i32(base.Ui32((v445-v429)&v447) < base.Ui32(v451)) <= v434-v435 {
												if v451 < v438 {
													v460 = v436
												} else {
													v460 = v438
												}
												if base.Ui32(v434) <= base.Ui32(v451) {
													v462 = v460
												} else {
													v462 = v438
												}
												v463 = v462
											} else {
												if base.Ui32(v451) < base.Ui32(v435) {
													v456 = v435
												} else {
													v456 = v436
												}
												if base.Ui32(v434) <= base.Ui32(v451) {
													v458 = v435
												} else {
													v458 = v456
												}
												v463 = v458
											}
											v464 = v297 - v428
											v466 = int32(255)
											v467 = (v464 - v273) & v466
											v468 = v467 & v387
											v469 = v468 + v348
											v471 = v468 + v325
											v478 = (v445 - v273) & v466
											if v471-v467+base.B2i32(base.Ui32((v445-v464)&v466) < base.Ui32(v478)) <= v467-v468 {
												if v478 < v471 {
													v487 = v469
												} else {
													v487 = v471
												}
												if base.Ui32(v467) <= base.Ui32(v478) {
													v489 = v487
												} else {
													v489 = v471
												}
												v490 = v489
											} else {
												if base.Ui32(v478) < base.Ui32(v468) {
													v483 = v468
												} else {
													v483 = v469
												}
												if base.Ui32(v467) <= base.Ui32(v478) {
													v485 = v468
												} else {
													v485 = v483
												}
												v490 = v485
											}
											v510 = v273 & int32(16711935)
											v518 = v490&int32(255) | (v463<<(uint(int32(16))%32)&int32(16711680) | v386<<(uint(int32(24))%32) | v422<<(uint(int32(8))%32)&int32(_a_F_GetResidual_0))
											v520 = v273 & int32(-16711936)
										} else {
											v304 = int32(-16711936)
											v305 = v273 & v304
											v311 = int32(16711935)
											v312 = v273 & v311
											v510 = v312
											v518 = (v297|int32(16711680)-v305)&v304 | (v297|int32(_a_F_GetResidual_0)-v312)&v311
											v520 = v305
										}
										v540 = int32(-16711936)
										v545 = int32(16711935)
										v550 = (v518&v540+v520)&v540 | (v518&v545+v510)&v545
										*(*int32)(unsafe.Add(mBase, uint32(v296))) = v550
										v553 = v273
										v560 = v518
										v562 = v550
									} else {
										v275 = v273
										v280 = *(*int32)(unsafe.Add(mBase, uint32(l3+v243<<(uint(int32(2))%32))))
										v283 = int32(-16711936)
										v290 = int32(16711935)
										v553 = v275
										v560 = (v280|int32(16711680)-v275&v283)&v283 | (v280|int32(_a_F_GetResidual_0)-v275&v290)&v290
										v562 = v280
									}
								} else {
									v268 = *(*int32)(unsafe.Add(mBase, uint32(l2)))
									v275 = v268
									v280 = *(*int32)(unsafe.Add(mBase, uint32(l3+v243<<(uint(int32(2))%32))))
									v283 = int32(-16711936)
									v290 = int32(16711935)
									v553 = v275
									v560 = (v280|int32(16711680)-v275&v283)&v283 | (v280|int32(_a_F_GetResidual_0)-v275&v290)&v290
									v562 = v280
								}
								if base.Ui32(int32(16777215)) < base.Ui32(v562) {
									v594 = v560
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(l3+v243<<(uint(int32(2))%32)))) = v553 & int32(16777215)
									v591 = v560 & int32(-16777216)
									if v243 != 0 {
										v594 = v591
									} else {
										v592 = *(*int32)(unsafe.Add(mBase, uint32(l3)))
										*(*int32)(unsafe.Add(mBase, uint32(v150))) = v592
										v594 = v591
									}
								}
								*(*int32)(unsafe.Add(mBase, uint32(l12+(v243-l6)<<(uint(int32(2))%32)))) = v594
								v601 = v243 + int32(1)
								if v601 != l7 {
									v243 = v601
									continue
								} else {
									break
								}
								break
							}
						}
					} else {
						v163 = l6
						v165 = l6 << (uint(int32(2)) % 32)
						v169 = l12
						for {
							if v163 == int32(0) {
								v194 = *(*int32)(unsafe.Add(mBase, uint32(l2)))
								v195 = v194
							} else {
								v193 = m.T0[v147].(func(*base.Module, int32, int32) int32)(m, l3+v165+int32(-4), l2+v165)
								mBase = m.M
								v195 = v193
							}
							v196 = l3 + v165
							v197 = *(*int32)(unsafe.Add(mBase, uint32(v196)))
							v202 = v197 | int32(16711680) - v195&int32(-16711936)
							if base.Ui32(v197) <= base.Ui32(int32(16777215)) {
								*(*int32)(unsafe.Add(mBase, uint32(v196))) = v195 & int32(16777215)
								v219 = v202 & int32(-16777216)
								if v163 != 0 {
									v222 = v219
								} else {
									v220 = *(*int32)(unsafe.Add(mBase, uint32(l3)))
									*(*int32)(unsafe.Add(mBase, uint32(v150))) = v220
									v222 = v219
								}
							} else {
								v209 = int32(16711935)
								v222 = v202&int32(-16711936) | (v197|int32(_a_F_GetResidual_0)-v195&v209)&v209
							}
							*(*int32)(unsafe.Add(mBase, uint32(v169))) = v222
							v224 = int32(4)
							v229 = v163 + int32(1)
							if l7 != v229 {
								v163 = v229
								v165 = v165 + v224
								v169 = v169 + v224
								continue
							} else {
								break
							}
							break
						}
					}
				}
			} else {
				v81 = l6
				v85 = l6<<(uint(int32(2))%32) + l3 + int32(-4)
				v87 = l12
				for {
					if v81 == int32(0) {
						v109 = int32(-16777216)
					} else {
						v108 = *(*int32)(unsafe.Add(mBase, uint32(v85)))
						v109 = v108
					}
					v111 = v85 + int32(4)
					v112 = *(*int32)(unsafe.Add(mBase, uint32(v111)))
					v117 = v112 | int32(16711680) - v109&int32(-16711936)
					if base.Ui32(v112) <= base.Ui32(int32(16777215)) {
						*(*int32)(unsafe.Add(mBase, uint32(v85+int32(4)))) = v109 & int32(16777215)
						v137 = v117 & int32(-16777216)
					} else {
						v124 = int32(16711935)
						v137 = v117&int32(-16711936) | (v112|int32(_a_F_GetResidual_0)-v109&v124)&v124
					}
					*(*int32)(unsafe.Add(mBase, uint32(v87))) = v137
					v142 = v81 + int32(1)
					if l7 != v142 {
						v81 = v142
						v85 = v111
						v87 = v87 + int32(4)
						continue
					} else {
						break
					}
					break
				}
			}
		}
		return
	} else {
		v33 = l7 - l6
		if l6 != 0 {
			v48 = l6
			v49 = v33
			v50 = l12
		} else {
			if l8 != 0 {
				v40 = m.G102
				v41 = *(*int32)(unsafe.Add(mBase, uint32(v40)+8))
				m.T0[v41].(func(*base.Module, int32, int32, int32, int32))(m, l3, l2, int32(1), l12)
				mBase = m.M
			} else {
				v36 = m.G102
				v37 = *(*int32)(unsafe.Add(mBase, uint32(v36)))
				m.T0[v37].(func(*base.Module, int32, int32, int32, int32))(m, l3, int32(0), int32(1), l12)
				mBase = m.M
			}
			v48 = int32(1)
			v49 = v33 + int32(-1)
			v50 = l12 + int32(4)
		}
		if l8 != 0 {
			v58 = int32(2)
			v59 = v48 << (uint(v58) % 32)
			v62 = m.G102
			v66 = *(*int32)(unsafe.Add(mBase, uint32(v62+l5<<(uint(v58)%32))))
			m.T0[v66].(func(*base.Module, int32, int32, int32, int32))(m, l3+v59, l2+v59, v49, v50)
			mBase = m.M
			return
		} else {
			v55 = m.G102
			v56 = *(*int32)(unsafe.Add(mBase, uint32(v55)+4))
			m.T0[v56].(func(*base.Module, int32, int32, int32, int32))(m, l3+v48<<(uint(int32(2))%32), int32(0), v49, v50)
			mBase = m.M
			return
		}
	}
}
func F_GetResidualCost_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v27 int32
	_ = v27
	var v29 int32
	_ = v29
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v42 int32
	_ = v42
	var v48 int32
	_ = v48
	var v53 int32
	_ = v53
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v63 int32
	_ = v63
	var v66 int32
	_ = v66
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v85 int32
	_ = v85
	var v87 int32
	_ = v87
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v95 int32
	_ = v95
	var v101 int32
	_ = v101
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v109 int32
	_ = v109
	var v115 int32
	_ = v115
	var v118 int32
	_ = v118
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v126 int32
	_ = v126
	var v127 int32
	_ = v127
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v137 int32
	_ = v137
	var v139 int32
	_ = v139
	var v142 int32
	_ = v142
	var v148 int32
	_ = v148
	var v150 int32
	_ = v150
	var v158 int32
	_ = v158
	var v160 int32
	_ = v160
	var v164 int32
	_ = v164
	var v166 int32
	_ = v166
	var v169 int32
	_ = v169
	var v175 int32
	_ = v175
	v11 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	v12 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	v19 = *(*int32)(unsafe.Add(mBase, uint32(v11+v12*int32(12)+l0<<(uint(int32(2))%32))))
	v20 = *(*int32)(unsafe.Add(mBase, uint32(l1)+16))
	v27 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v20+v12*int32(33)+l0*int32(11)))))
	if l0 != 0 {
		v38 = int32(0)
	} else {
		v29 = m.G1
		v37 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v29+int32(_a_F_GetResidualCost_C_0)+(v27^int32(255))<<(uint(int32(1))%32)))))
		v38 = v37
	}
	v39 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
	if v39 < int32(0) {
		v169 = m.G1
		v175 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v169+int32(_a_F_GetResidualCost_C_0)+v27<<(uint(int32(1))%32)))))
		return v175
	} else {
		v42 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
		if v39 <= v12 {
			v105 = v12
			v106 = v19
			v109 = v38
		} else {
			v48 = int32(12)
			v53 = v42 + v12<<(uint(int32(1))%32)
			v56 = v12*v48 + v11 + v48
			v57 = v19
			v59 = v39 - v12
			v60 = v38
			for {
				v63 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v53))))
				v66 = base.I32_extend16_s(v63) >> (uint(int32(15)) % 32)
				v70 = (v63 ^ v66 - v66) & int32(_a_F_GetResidualCost_C_1)
				v71 = int32(67)
				if base.Ui32(v70) < base.Ui32(v71) {
					v74 = v70
				} else {
					v74 = v71
				}
				v75 = int32(1)
				v78 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v57+v74<<(uint(v75)%32)))))
				v79 = m.G1
				v85 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v79+int32(_a_F_GetResidualCost_C_2)+v70<<(uint(v75)%32)))))
				v87 = v78 + (v60 + v85)
				if base.Ui32(v75) < base.Ui32(v70) {
					v91 = int32(2)
				} else {
					v91 = v70
				}
				v92 = int32(2)
				v95 = *(*int32)(unsafe.Add(mBase, uint32(v56+v91<<(uint(v92)%32))))
				v101 = v59 + int32(-1)
				if v101 != 0 {
					v53 = v53 + v92
					v56 = v56 + int32(12)
					v57 = v95
					v59 = v101
					v60 = v87
					continue
				} else {
					break
				}
				break
			}
			v105 = v39
			v106 = v95
			v109 = v87
		}
		v115 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v42+v105<<(uint(int32(1))%32)))))
		v118 = base.I32_extend16_s(v115) >> (uint(int32(15)) % 32)
		v122 = (v115 ^ v118 - v118) & int32(_a_F_GetResidualCost_C_1)
		v123 = int32(67)
		if base.Ui32(v122) < base.Ui32(v123) {
			v126 = v122
		} else {
			v126 = v123
		}
		v127 = int32(1)
		v130 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v106+v126<<(uint(v127)%32)))))
		v131 = m.G1
		v137 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v131+int32(_a_F_GetResidualCost_C_2)+v122<<(uint(v127)%32)))))
		v139 = v130 + (v109 + v137)
		if int32(14) < v105 {
			v166 = v139
		} else {
			v142 = m.G1
			v148 = int32(1)
			v150 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v105+(v142+int32(_a_F_GetResidualCost_C_3))+v148))))
			if v122 == v148 {
				v158 = int32(11)
			} else {
				v158 = int32(22)
			}
			v160 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v20+v150*int32(33)+v158))))
			v164 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v142+int32(_a_F_GetResidualCost_C_0)+v160<<(uint(int32(1))%32)))))
			v166 = v139 + v164
		}
		return v166
	}
}
func F_GradientUnfilter_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v19 int32
	_ = v19
	var v21 int32
	_ = v21
	var v26 int32
	_ = v26
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v34 int32
	_ = v34
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v41 int32
	_ = v41
	var v45 int32
	_ = v45
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
	var v63 int32
	_ = v63
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v74 int32
	_ = v74
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v81 int32
	_ = v81
	var v86 int32
	_ = v86
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v95 int32
	_ = v95
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v116 int32
	_ = v116
	var v123 int32
	_ = v123
	if l0 != 0 {
		if l3 < int32(1) {
		} else {
			v89 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
			v90 = l0
			v91 = l1
			v92 = l2
			v93 = l3
			v95 = v89
			v97 = v89
			for {
				v98 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v91))))
				v99 = int32(255)
				v104 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v90))))
				v105 = v95&v99 - v97&v99 + v104
				v106 = int32(0)
				if v106 < v105 {
					v109 = v105
				} else {
					v109 = v106
				}
				v110 = int32(255)
				if v109 < v110 {
					v113 = v109
				} else {
					v113 = v110
				}
				v114 = v98 + v113
				*(*uint8)(unsafe.Add(mBase, uint32(v92))) = uint8(v114)
				v116 = int32(1)
				v123 = v93 + int32(-1)
				if v123 != 0 {
					v90 = v90 + v116
					v91 = v91 + v116
					v92 = v92 + v116
					v93 = v123
					v95 = v114
					v97 = v104
					continue
				} else {
					break
				}
				break
			}
		}
	} else {
		if l3 < int32(1) {
		} else {
			v12 = l3 & int32(3)
			v13 = int32(0)
			if base.Ui32(l3) < base.Ui32(int32(4)) {
				v58 = v13
				v63 = v13
			} else {
				v19 = int32(0)
				v21 = v19
				v26 = v19
				for {
					v29 = l2 + v26
					v30 = l1 + v26
					v31 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v30))))
					v32 = v31 + v21
					*(*uint8)(unsafe.Add(mBase, uint32(v29))) = uint8(v32)
					v34 = int32(1)
					v38 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v30+v34))))
					v39 = v38 + v32
					*(*uint8)(unsafe.Add(mBase, uint32(v29+v34))) = uint8(v39)
					v41 = int32(2)
					v45 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v30+v41))))
					v46 = v45 + v39
					*(*uint8)(unsafe.Add(mBase, uint32(v29+v41))) = uint8(v46)
					v48 = int32(3)
					v52 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v30+v48))))
					v53 = v52 + v46
					*(*uint8)(unsafe.Add(mBase, uint32(v29+v48))) = uint8(v53)
					v56 = v26 + int32(4)
					if l3&int32(2147483644) != v56 {
						v21 = v53
						v26 = v56
						continue
					} else {
						break
					}
					break
				}
				v58 = v53
				v63 = v56
			}
			if v12 == int32(0) {
			} else {
				v70 = v58
				v71 = l1 + v63
				v72 = l2 + v63
				v74 = v12
				for {
					v78 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v71))))
					v79 = v78 + v70
					*(*uint8)(unsafe.Add(mBase, uint32(v72))) = uint8(v79)
					v81 = int32(1)
					v86 = v74 + int32(-1)
					if v86 != 0 {
						v70 = v79
						v71 = v71 + v81
						v72 = v72 + v81
						v74 = v86
						continue
					} else {
						break
					}
					break
				}
			}
		}
	}
	return
}
