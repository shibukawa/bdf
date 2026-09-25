//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
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
func F_QuantizeLevels(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) int32 {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v418 int32
	_ = v418
	var v420 int32
	_ = v420
	var v422 int32
	_ = v422
	var v426 int32
	_ = v426
	var v432 int32
	_ = v432
	var v438 int32
	_ = v438
	var v444 int32
	_ = v444
	var v445 int32
	_ = v445
	var v446 int32
	_ = v446
	var v467 int32
	_ = v467
	var v468 int32
	_ = v468
	var v469 int32
	_ = v469
	var v470 int32
	_ = v470
	var v472 int32
	_ = v472
	var v473 int32
	_ = v473
	var v474 int32
	_ = v474
	var v481 int32
	_ = v481
	var v484 int32
	_ = v484
	var v485 int32
	_ = v485
	var v490 int32
	_ = v490
	var v492 int32
	_ = v492
	var v494 int32
	_ = v494
	var v496 int32
	_ = v496
	var v497 int32
	_ = v497
	var v502 int32
	_ = v502
	var v504 int32
	_ = v504
	var v508 int32
	_ = v508
	var v514 int32
	_ = v514
	var v515 int32
	_ = v515
	var v516 int32
	_ = v516
	var v541 int32
	_ = v541
	var v544 int32
	_ = v544
	var v545 int32
	_ = v545
	var v550 int32
	_ = v550
	var v552 int32
	_ = v552
	var v564 int32
	_ = v564
	var v565 int32
	_ = v565
	var v566 int32
	_ = v566
	var v589 int32
	_ = v589
	var v594 int32
	_ = v594
	var v597 int32
	_ = v597
	var v598 float64
	_ = v598
	var v600 float64
	_ = v600
	var v601 float64
	_ = v601
	var v611 int32
	_ = v611
	var v616 int32
	_ = v616
	var v631 float64
	_ = v631
	var v658 int32
	_ = v658
	var v661 int32
	_ = v661
	var v717 int32
	_ = v717
	var v735 int32
	_ = v735
	var v736 int32
	_ = v736
	var v737 int32
	_ = v737
	var v738 int32
	_ = v738
	var v739 int32
	_ = v739
	var v747 int32
	_ = v747
	var v750 int32
	_ = v750
	var v757 int32
	_ = v757
	var v787 float64
	_ = v787
	var v797 int32
	_ = v797
	var v798 int32
	_ = v798
	var v1044 int32
	_ = v1044
	var v1048 int32
	_ = v1048
	var v1052 int32
	_ = v1052
	var v1077 int32
	_ = v1077
	var v1087 int32
	_ = v1087
	var v1092 int32
	_ = v1092
	var v1119 float64
	_ = v1119
	var v1121 int32
	_ = v1121
	var v1122 float64
	_ = v1122
	var v1129 int32
	_ = v1129
	var v1134 int32
	_ = v1134
	var v1136 int32
	_ = v1136
	var v1142 int32
	_ = v1142
	var v1143 int32
	_ = v1143
	var v1144 float64
	_ = v1144
	var v1149 int32
	_ = v1149
	var v1150 float64
	_ = v1150
	var v1196 int32
	_ = v1196
	var v1199 int32
	_ = v1199
	var v1204 int32
	_ = v1204
	var v1228 int32
	_ = v1228
	var v1231 float64
	_ = v1231
	var v1239 int32
	_ = v1239
	var v1246 float64
	_ = v1246
	var v1251 float64
	_ = v1251
	var v1259 int32
	_ = v1259
	var v1266 float64
	_ = v1266
	var v1272 int32
	_ = v1272
	var v1282 int32
	_ = v1282
	var v1309 int32
	_ = v1309
	var v1311 float64
	_ = v1311
	var v1322 float64
	_ = v1322
	var v1361 int32
	_ = v1361
	var v1365 float64
	_ = v1365
	var v1366 float64
	_ = v1366
	var v1367 int32
	_ = v1367
	var v1373 int32
	_ = v1373
	var v1374 float64
	_ = v1374
	var v1379 int32
	_ = v1379
	var v1385 int32
	_ = v1385
	var v1386 int32
	_ = v1386
	var v1390 int32
	_ = v1390
	var v1405 float64
	_ = v1405
	var v1415 int32
	_ = v1415
	var v1418 int32
	_ = v1418
	var v1419 int32
	_ = v1419
	var v1421 int32
	_ = v1421
	var v1422 int32
	_ = v1422
	var v1425 float64
	_ = v1425
	var v1426 float64
	_ = v1426
	var v1429 int32
	_ = v1429
	var v1436 int32
	_ = v1436
	var v1440 float64
	_ = v1440
	var v1441 float64
	_ = v1441
	var v1442 int32
	_ = v1442
	var v1447 float64
	_ = v1447
	var v1450 int32
	_ = v1450
	var v1476 float64
	_ = v1476
	var v1488 int32
	_ = v1488
	var v1504 int32
	_ = v1504
	var v1508 float64
	_ = v1508
	var v1510 float64
	_ = v1510
	var v1518 int32
	_ = v1518
	var v1520 int32
	_ = v1520
	var v1527 int32
	_ = v1527
	var v1543 int32
	_ = v1543
	var v1544 int32
	_ = v1544
	var v1548 int32
	_ = v1548
	var v1574 int32
	_ = v1574
	var v1578 float64
	_ = v1578
	var v1580 float64
	_ = v1580
	var v1588 int32
	_ = v1588
	var v1590 int32
	_ = v1590
	var v1596 int32
	_ = v1596
	var v1600 float64
	_ = v1600
	var v1602 float64
	_ = v1602
	var v1610 int32
	_ = v1610
	var v1612 int32
	_ = v1612
	var v1621 int32
	_ = v1621
	var v1655 int32
	_ = v1655
	var v1663 int32
	_ = v1663
	var v1692 int32
	_ = v1692
	var v1694 int32
	_ = v1694
	var v1695 int32
	_ = v1695
	var v1697 int32
	_ = v1697
	var v1700 int32
	_ = v1700
	var v1703 int32
	_ = v1703
	var v1705 int32
	_ = v1705
	var v1708 int32
	_ = v1708
	var v1711 int32
	_ = v1711
	var v1713 int32
	_ = v1713
	var v1716 int32
	_ = v1716
	var v1719 int32
	_ = v1719
	var v1721 int32
	_ = v1721
	var v1724 int32
	_ = v1724
	var v1727 int32
	_ = v1727
	var v1761 int32
	_ = v1761
	var v1765 int32
	_ = v1765
	var v1791 int32
	_ = v1791
	var v1793 int32
	_ = v1793
	var v1798 int32
	_ = v1798
	var v1836 int64
	_ = v1836
	var v1854 int64
	_ = v1854
	var v1868 int32
	_ = v1868
	var v1878 int32
	_ = v1878
	v6 = int32(0)
	v31 = m.G0
	v33 = v31 - int32(_a_F_QuantizeLevels_0)
	m.G0 = v33
	goto L3
L1:
	;
	goto L17
L3:
	;
	base.MemoryFill(m, v33+int32(_a_F_QuantizeLevels_1), v6, int32(1024))
	goto L1
L15:
	;
	goto L31
L17:
	;
	base.MemoryFill(m, v33+int32(_a_F_QuantizeLevels_2), int32(0), int32(1024))
	goto L15
L29:
	;
	if l0 == int32(0) {
		v1878 = v6
		goto L43
	} else {
		goto L44
	}
L31:
	;
	base.MemoryFill(m, v33+int32(_a_F_QuantizeLevels_3), int32(0), int32(2048))
	goto L29
L43:
	;
	m.G0 = v33 + int32(_a_F_QuantizeLevels_0)
	return v1878
L44:
	;
	if l1 < int32(1) {
		v1878 = v6
		goto L43
	} else {
		goto L45
	}
L45:
	;
	if l2 < int32(1) {
		v1878 = v6
		goto L43
	} else {
		goto L46
	}
L46:
	;
	if base.Ui32(l3+int32(-257)) < base.Ui32(int32(-255)) {
		v1878 = v6
		goto L43
	} else {
		goto L47
	}
L47:
	;
	v418 = l2 * l1
	if v418 != 0 {
		goto L49
	} else {
		goto L50
	}
L48:
	;
	if v565 <= l3 {
		v1854 = int64(0)
		goto L76
	} else {
		goto L77
	}
L49:
	;
	v422 = int32(1)
	if v418 != v422 {
		goto L52
	} else {
		goto L53
	}
L50:
	;
	v420 = int32(0)
	v564 = int32(255)
	v565 = v420
	v566 = v420
	goto L48
L51:
	;
	if v418&v422 == int32(0) {
		v564 = v514
		v565 = v515
		v566 = v516
		goto L48
	} else {
		goto L69
	}
L52:
	;
	v432 = int32(0)
	v438 = v432
	v444 = int32(255)
	v445 = v432
	v446 = v432
	goto L54
L53:
	;
	v426 = int32(0)
	v508 = v426
	v514 = int32(255)
	v515 = v426
	v516 = v426
	goto L51
L54:
	;
	v467 = v33 + int32(_a_F_QuantizeLevels_1)
	v468 = l0 + v438
	v469 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v468))))
	v470 = int32(2)
	v472 = v467 + v469<<(uint(v470)%32)
	v473 = *(*int32)(unsafe.Add(mBase, uint32(v472)))
	v474 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v472))) = v473 + v474
	v481 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v468+v474))))
	v484 = v467 + v481<<(uint(v470)%32)
	v485 = *(*int32)(unsafe.Add(mBase, uint32(v484)))
	*(*int32)(unsafe.Add(mBase, uint32(v484))) = v485 + v474
	if v469 < v446 {
		goto L56
	} else {
		goto L57
	}
L55:
	;
	v508 = v504
	v514 = v496
	v515 = v502
	v516 = v492
	goto L51
L56:
	;
	v490 = v446
	goto L58
L57:
	;
	v490 = v469
	goto L58
L58:
	;
	if v481 < v490 {
		goto L59
	} else {
		goto L60
	}
L59:
	;
	v492 = v490
	goto L61
L60:
	;
	v492 = v481
	goto L61
L61:
	;
	if v444 < v469 {
		goto L62
	} else {
		goto L63
	}
L62:
	;
	v494 = v444
	goto L64
L63:
	;
	v494 = v469
	goto L64
L64:
	;
	if v494 < v481 {
		goto L65
	} else {
		goto L66
	}
L65:
	;
	v496 = v494
	goto L67
L66:
	;
	v496 = v481
	goto L67
L67:
	;
	v497 = int32(0)
	v502 = v445 + base.B2i32(v473 == v497) + base.B2i32(v485 == v497)
	v504 = v438 + int32(2)
	if v418&int32(-2) != v504 {
		v438 = v504
		v444 = v496
		v445 = v502
		v446 = v492
		goto L54
	} else {
		goto L68
	}
L68:
	;
	goto L55
L69:
	;
	v541 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v508))))
	v544 = v33 + int32(_a_F_QuantizeLevels_1) + v541<<(uint(int32(2))%32)
	v545 = *(*int32)(unsafe.Add(mBase, uint32(v544)))
	*(*int32)(unsafe.Add(mBase, uint32(v544))) = v545 + int32(1)
	if v541 < v516 {
		goto L70
	} else {
		goto L71
	}
L70:
	;
	v550 = v516
	goto L72
L71:
	;
	v550 = v541
	goto L72
L72:
	;
	if v514 < v541 {
		goto L73
	} else {
		goto L74
	}
L73:
	;
	v552 = v514
	goto L75
L74:
	;
	v552 = v541
	goto L75
L75:
	;
	v564 = v552
	v565 = v515 + base.B2i32(v545 == int32(0))
	v566 = v550
	goto L48
L76:
	;
	v1868 = int32(1)
	if l4 == int32(0) {
		v1878 = v1868
		goto L43
	} else {
		goto L187
	}
L77:
	;
	v589 = int32(0)
	if v589 < l3 {
		goto L79
	} else {
		goto L80
	}
L78:
	;
	v735 = v33 + int32(_a_F_QuantizeLevels_2)
	v736 = int32(2)
	v737 = v564 << (uint(v736) % 32)
	v738 = v735 + v737
	v739 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v738))) = v739
	*(*int32)(unsafe.Add(mBase, uint32(v735+v566<<(uint(v736)%32)))) = v717
	v747 = int32(1)
	v750 = v717 + int32(-1)
	v757 = v566 - v564 + v747
	v787 = float64(1e+38)
	v797 = v739
	goto L87
L79:
	;
	v594 = int32(1)
	v597 = l3 + int32(-1)
	v598 = base.F64_convert_i32_u(v597)
	v600 = base.F64_convert_i32_s(v566 - v564)
	v601 = base.F64_convert_i32_u(v564)
	if l3 == v594 {
		v661 = v589
		goto L81
	} else {
		goto L82
	}
L80:
	;
	v717 = l3 + int32(-1)
	goto L78
L81:
	;
	if l3&v594 == int32(0) {
		v717 = v597
		goto L78
	} else {
		goto L86
	}
L82:
	;
	v611 = int32(0)
	v616 = v33 + int32(_a_F_QuantizeLevels_3)
	v631 = float64(0)
	goto L83
L83:
	;
	*(*float64)(unsafe.Add(mBase, uint32(v616))) = base.F64_add(base.F64_div(base.F64_mul(v600, v631), v598), v601)
	*(*float64)(unsafe.Add(mBase, uint32(v616+int32(8)))) = base.F64_add(base.F64_div(base.F64_mul(v600, base.F64_convert_i32_u(v611+int32(1))), v598), v601)
	v658 = v611 + int32(2)
	if l3&int32(510) != v658 {
		v611 = v658
		v616 = v616 + int32(16)
		v631 = base.F64_add(v631, float64(2))
		goto L83
	} else {
		goto L85
	}
L84:
	;
	v661 = v658
	goto L81
L85:
	;
	goto L84
L86:
	;
	*(*float64)(unsafe.Add(mBase, uint32(v33+int32(_a_F_QuantizeLevels_3)+v661<<(uint(int32(3))%32)))) = base.F64_add(base.F64_div(base.F64_mul(v600, base.F64_convert_i32_u(v661)), v598), v601)
	v717 = v597
	goto L78
L87:
	;
	v798 = int32(2048)
	goto L91
L88:
	;
	if v566 < v564 {
		goto L157
	} else {
		goto L158
	}
L89:
	;
	goto L105
L91:
	;
	base.MemoryFill(m, v33+v798, int32(0), v798)
	goto L89
L103:
	;
	v1044 = base.B2i32(v566 < v564)
	if v566 < v564 {
		goto L117
	} else {
		goto L118
	}
L105:
	;
	base.MemoryFill(m, v33, int32(0), int32(2048))
	goto L103
L117:
	;
	if l3 < int32(3) {
		goto L133
	} else {
		goto L134
	}
L118:
	;
	v1048 = v564
	v1052 = int32(0)
	goto L119
L119:
	;
	if v717 < v1052 {
		goto L121
	} else {
		goto L122
	}
L120:
	;
	goto L117
L121:
	;
	v1077 = v1052
	goto L123
L122:
	;
	v1077 = v717
	goto L123
L123:
	;
	v1087 = v33 + int32(_a_F_QuantizeLevels_3) + v1052<<(uint(int32(3))%32)
	v1092 = v1052
	goto L125
L124:
	;
	v1134 = v1048 << (uint(int32(2)) % 32)
	v1136 = *(*int32)(unsafe.Add(mBase, uint32(v33+int32(_a_F_QuantizeLevels_1)+v1134)))
	if v1136 < int32(1) {
		goto L130
	} else {
		goto L131
	}
L125:
	;
	if v1077 != v1092 {
		goto L127
	} else {
		goto L128
	}
L126:
	;
	v1129 = v1092 + int32(0)
	goto L124
L127:
	;
	v1119 = *(*float64)(unsafe.Add(mBase, uint32(v1087)))
	v1121 = v1087 + int32(8)
	v1122 = *(*float64)(unsafe.Add(mBase, uint32(v1121)))
	if base.F64_lt(base.F64_add(v1119, v1122), base.F64_convert_i32_s(v1048<<(uint(int32(1))%32))) != 0 {
		v1087 = v1121
		v1092 = v1092 + int32(1)
		goto L125
	} else {
		goto L129
	}
L128:
	;
	v1129 = v1077
	goto L124
L129:
	;
	goto L126
L130:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v33+int32(_a_F_QuantizeLevels_2)+v1134))) = v1129
	if v1048 != v566 {
		v1048 = v1048 + int32(1)
		v1052 = v1129
		goto L119
	} else {
		goto L132
	}
L131:
	;
	v1142 = v1129 << (uint(int32(3)) % 32)
	v1143 = v33 + int32(2048) + v1142
	v1144 = *(*float64)(unsafe.Add(mBase, uint32(v1143)))
	*(*float64)(unsafe.Add(mBase, uint32(v1143))) = base.F64_add(v1144, base.F64_convert_i32_s(v1136*v1048))
	v1149 = v33 + v1142
	v1150 = *(*float64)(unsafe.Add(mBase, uint32(v1149)))
	*(*float64)(unsafe.Add(mBase, uint32(v1149))) = base.F64_add(v1150, base.F64_convert_i32_u(v1136))
	goto L130
L132:
	;
	goto L120
L133:
	;
	if v566 < v564 {
		v1476 = float64(0)
		goto L146
	} else {
		goto L147
	}
L134:
	;
	if v717 == int32(2) {
		v1282 = int32(1)
		goto L135
	} else {
		goto L136
	}
L135:
	;
	if v750&v747 == int32(0) {
		goto L133
	} else {
		goto L144
	}
L136:
	;
	v1196 = int32(0)
	v1199 = v1196
	v1204 = v1196
	goto L137
L137:
	;
	v1228 = v33 + v1204
	v1231 = *(*float64)(unsafe.Add(mBase, uint32(v1228+int32(8))))
	if base.F64_gt(v1231, float64(0)) == int32(0) {
		goto L139
	} else {
		goto L140
	}
L138:
	;
	v1282 = v1199 + int32(3)
	goto L135
L139:
	;
	v1251 = *(*float64)(unsafe.Add(mBase, uint32(v1228+int32(16))))
	if base.F64_gt(v1251, float64(0)) == int32(0) {
		goto L141
	} else {
		goto L142
	}
L140:
	;
	v1239 = int32(8)
	v1246 = *(*float64)(unsafe.Add(mBase, uint32(v33+int32(2048)+v1204+v1239)))
	*(*float64)(unsafe.Add(mBase, uint32(v33+int32(_a_F_QuantizeLevels_3)+v1204+v1239))) = base.F64_div(v1246, v1231)
	goto L139
L141:
	;
	v1272 = v1199 + int32(2)
	if v750&int32(-2) != v1272 {
		v1199 = v1272
		v1204 = v1204 + int32(16)
		goto L137
	} else {
		goto L143
	}
L142:
	;
	v1259 = int32(16)
	v1266 = *(*float64)(unsafe.Add(mBase, uint32(v33+int32(2048)+v1204+v1259)))
	*(*float64)(unsafe.Add(mBase, uint32(v33+int32(_a_F_QuantizeLevels_3)+v1204+v1259))) = base.F64_div(v1266, v1251)
	goto L141
L143:
	;
	goto L138
L144:
	;
	v1309 = v1282 << (uint(int32(3)) % 32)
	v1311 = *(*float64)(unsafe.Add(mBase, uint32(v33+v1309)))
	if base.F64_gt(v1311, float64(0)) == int32(0) {
		goto L133
	} else {
		goto L145
	}
L145:
	;
	v1322 = *(*float64)(unsafe.Add(mBase, uint32(v33+int32(2048)+v1309)))
	*(*float64)(unsafe.Add(mBase, uint32(v33+int32(_a_F_QuantizeLevels_3)+v1309))) = base.F64_div(v1322, v1311)
	goto L133
L146:
	;
	if base.F64_lt(base.F64_sub(v787, v1476), base.F64_mul(base.F64_convert_i32_u(v418), float64(0.0001))) != 0 {
		goto L154
	} else {
		goto L155
	}
L147:
	;
	if v757&v747 == int32(0) {
		v1373 = v564
		v1374 = float64(0)
		goto L148
	} else {
		goto L149
	}
L148:
	;
	if v566 == v564 {
		v1476 = v1374
		goto L146
	} else {
		goto L150
	}
L149:
	;
	v1361 = *(*int32)(unsafe.Add(mBase, uint32(v738)))
	v1365 = *(*float64)(unsafe.Add(mBase, uint32(v33+int32(_a_F_QuantizeLevels_3)+v1361<<(uint(int32(3))%32))))
	v1366 = base.F64_sub(base.F64_convert_i32_s(v564), v1365)
	v1367 = *(*int32)(unsafe.Add(mBase, uint32(v33+int32(_a_F_QuantizeLevels_1)+v737)))
	v1373 = v564 + v747
	v1374 = base.F64_add(base.F64_mul(base.F64_mul(v1366, base.F64_convert_i32_s(v1367)), v1366), float64(0))
	goto L148
L150:
	;
	v1379 = v1373 << (uint(int32(2)) % 32)
	v1385 = v33 + int32(_a_F_QuantizeLevels_2) + v1379
	v1386 = v33 + int32(_a_F_QuantizeLevels_1) + v1379
	v1390 = v1373
	v1405 = v1374
	goto L151
L151:
	;
	v1415 = v1390 + int32(1)
	v1418 = v33 + int32(_a_F_QuantizeLevels_3)
	v1419 = int32(4)
	v1421 = *(*int32)(unsafe.Add(mBase, uint32(v1385+v1419)))
	v1422 = int32(3)
	v1425 = *(*float64)(unsafe.Add(mBase, uint32(v1418+v1421<<(uint(v1422)%32))))
	v1426 = base.F64_sub(base.F64_convert_i32_s(v1415), v1425)
	v1429 = *(*int32)(unsafe.Add(mBase, uint32(v1386+v1419)))
	v1436 = *(*int32)(unsafe.Add(mBase, uint32(v1385)))
	v1440 = *(*float64)(unsafe.Add(mBase, uint32(v1418+v1436<<(uint(v1422)%32))))
	v1441 = base.F64_sub(base.F64_convert_i32_s(v1390), v1440)
	v1442 = *(*int32)(unsafe.Add(mBase, uint32(v1386)))
	v1447 = base.F64_add(base.F64_mul(base.F64_mul(v1426, base.F64_convert_i32_s(v1429)), v1426), base.F64_add(base.F64_mul(base.F64_mul(v1441, base.F64_convert_i32_s(v1442)), v1441), v1405))
	v1450 = int32(8)
	if v1415 != v566 {
		v1385 = v1385 + v1450
		v1386 = v1386 + v1450
		v1390 = v1390 + int32(2)
		v1405 = v1447
		goto L151
	} else {
		goto L153
	}
L152:
	;
	v1476 = v1447
	goto L146
L153:
	;
	goto L152
L154:
	;
	goto L88
L155:
	;
	v1488 = v797 + int32(1)
	if v1488 != int32(6) {
		v787 = v1476
		v797 = v1488
		goto L87
	} else {
		goto L156
	}
L156:
	;
	goto L154
L157:
	;
	if v418 == int32(0) {
		goto L174
	} else {
		goto L175
	}
L158:
	;
	if v757&int32(1) == int32(0) {
		v1527 = v564
		goto L159
	} else {
		goto L160
	}
L159:
	;
	if v566 == v564 {
		goto L157
	} else {
		goto L164
	}
L160:
	;
	v1504 = *(*int32)(unsafe.Add(mBase, uint32(v33+int32(_a_F_QuantizeLevels_2)+v564<<(uint(int32(2))%32))))
	v1508 = *(*float64)(unsafe.Add(mBase, uint32(v33+int32(_a_F_QuantizeLevels_3)+v1504<<(uint(int32(3))%32))))
	v1510 = base.F64_add(v1508, float64(0.5))
	if base.F64_lt(v1510, float64(4.294967296e+09))&base.F64_ge(v1510, float64(0)) == int32(0) {
		goto L162
	} else {
		goto L163
	}
L161:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v33+int32(2048)+v564))) = uint8(v1520)
	v1527 = v564 + int32(1)
	goto L159
L162:
	;
	v1520 = int32(0)
	goto L161
L163:
	;
	v1518 = base.I32_trunc_f64_u(v1510)
	v1520 = v1518
	goto L161
L164:
	;
	v1543 = v33 + int32(2048) + v1527
	v1544 = v566 - v1527 + int32(1)
	v1548 = v33 + int32(_a_F_QuantizeLevels_2) + v1527<<(uint(int32(2))%32)
	goto L165
L165:
	;
	v1574 = *(*int32)(unsafe.Add(mBase, uint32(v1548)))
	v1578 = *(*float64)(unsafe.Add(mBase, uint32(v33+int32(_a_F_QuantizeLevels_3)+v1574<<(uint(int32(3))%32))))
	v1580 = base.F64_add(v1578, float64(0.5))
	if base.F64_lt(v1580, float64(4.294967296e+09))&base.F64_ge(v1580, float64(0)) == int32(0) {
		goto L168
	} else {
		goto L169
	}
L166:
	;
	goto L157
L167:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1543))) = uint8(v1590)
	v1596 = *(*int32)(unsafe.Add(mBase, uint32(v1548+int32(4))))
	v1600 = *(*float64)(unsafe.Add(mBase, uint32(v33+int32(_a_F_QuantizeLevels_3)+v1596<<(uint(int32(3))%32))))
	v1602 = base.F64_add(v1600, float64(0.5))
	if base.F64_lt(v1602, float64(4.294967296e+09))&base.F64_ge(v1602, float64(0)) == int32(0) {
		goto L171
	} else {
		goto L172
	}
L168:
	;
	v1590 = int32(0)
	goto L167
L169:
	;
	v1588 = base.I32_trunc_f64_u(v1580)
	v1590 = v1588
	goto L167
L170:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1543+int32(1)))) = uint8(v1612)
	v1621 = v1544 + int32(-2)
	if v1621 != 0 {
		v1543 = v1543 + int32(2)
		v1544 = v1621
		v1548 = v1548 + int32(8)
		goto L165
	} else {
		goto L173
	}
L171:
	;
	v1612 = int32(0)
	goto L170
L172:
	;
	v1610 = base.I32_trunc_f64_u(v1602)
	v1612 = v1610
	goto L170
L173:
	;
	goto L166
L174:
	;
	if base.F64_lt(v1476, float64(1.8446744073709552e+19))&base.F64_ge(v1476, float64(0)) == int32(0) {
		goto L185
	} else {
		goto L186
	}
L175:
	;
	v1655 = v418 & int32(3)
	if base.Ui32(v418) < base.Ui32(int32(4)) {
		v1727 = int32(0)
		goto L176
	} else {
		goto L177
	}
L176:
	;
	if v1655 == int32(0) {
		goto L174
	} else {
		goto L181
	}
L177:
	;
	v1663 = int32(0)
	goto L178
L178:
	;
	v1692 = l0 + v1663
	v1694 = v33 + int32(2048)
	v1695 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1692))))
	v1697 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1694+v1695))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1692))) = uint8(v1697)
	v1700 = v1692 + int32(1)
	v1703 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1700))))
	v1705 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1694+v1703))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1700))) = uint8(v1705)
	v1708 = v1692 + int32(2)
	v1711 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1708))))
	v1713 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1694+v1711))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1708))) = uint8(v1713)
	v1716 = v1692 + int32(3)
	v1719 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1716))))
	v1721 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1694+v1719))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1716))) = uint8(v1721)
	v1724 = v1663 + int32(4)
	if v418&int32(-4) != v1724 {
		v1663 = v1724
		goto L178
	} else {
		goto L180
	}
L179:
	;
	v1727 = v1724
	goto L176
L180:
	;
	goto L179
L181:
	;
	v1761 = v1655
	v1765 = l0 + v1727
	goto L182
L182:
	;
	v1791 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1765))))
	v1793 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v33+int32(2048)+v1791))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1765))) = uint8(v1793)
	v1798 = v1761 + int32(-1)
	if v1798 != 0 {
		v1761 = v1798
		v1765 = v1765 + int32(1)
		goto L182
	} else {
		goto L184
	}
L183:
	;
	goto L174
L184:
	;
	goto L183
L185:
	;
	v1854 = int64(0)
	goto L76
L186:
	;
	v1836 = base.I64_trunc_f64_u(v1476)
	v1854 = v1836
	goto L76
L187:
	;
	*(*int64)(unsafe.Add(mBase, uint32(l4))) = v1854
	v1878 = v1868
	goto L43
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
	F___qsort_r(m, l0, l1, l2, int32(182), l3)
	return
}
