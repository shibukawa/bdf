//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_UpsampleArgbLinePair_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v162 int32
	_ = v162
	var v164 int32
	_ = v164
	var v165 int32
	_ = v165
	var v166 int32
	_ = v166
	var v171 int32
	_ = v171
	var v174 int32
	_ = v174
	var v177 int32
	_ = v177
	var v179 int32
	_ = v179
	var v183 int32
	_ = v183
	var v184 int32
	_ = v184
	var v186 int32
	_ = v186
	var v193 int32
	_ = v193
	var v196 int32
	_ = v196
	var v198 int32
	_ = v198
	var v199 int32
	_ = v199
	var v201 int32
	_ = v201
	var v204 int32
	_ = v204
	var v207 int32
	_ = v207
	var v212 int32
	_ = v212
	var v214 int32
	_ = v214
	var v221 int32
	_ = v221
	var v224 int32
	_ = v224
	var v228 int32
	_ = v228
	var v235 int32
	_ = v235
	var v237 int32
	_ = v237
	var v244 int32
	_ = v244
	var v247 int32
	_ = v247
	var v251 int32
	_ = v251
	var v253 int32
	_ = v253
	var v256 int32
	_ = v256
	var v257 int32
	_ = v257
	var v260 int32
	_ = v260
	var v265 int32
	_ = v265
	var v267 int32
	_ = v267
	var v274 int32
	_ = v274
	var v277 int32
	_ = v277
	var v281 int32
	_ = v281
	var v286 int32
	_ = v286
	var v288 int32
	_ = v288
	var v295 int32
	_ = v295
	var v298 int32
	_ = v298
	var v302 int32
	_ = v302
	var v309 int32
	_ = v309
	var v311 int32
	_ = v311
	var v318 int32
	_ = v318
	var v321 int32
	_ = v321
	var v329 int32
	_ = v329
	var v331 int32
	_ = v331
	var v334 int32
	_ = v334
	var v338 int32
	_ = v338
	var v346 int32
	_ = v346
	var v358 int32
	_ = v358
	var v359 int32
	_ = v359
	var v363 int32
	_ = v363
	var v364 int32
	_ = v364
	var v380 int32
	_ = v380
	var v381 int32
	_ = v381
	var v382 base.V128
	_ = v382
	var v383 int32
	_ = v383
	var v385 base.V128
	_ = v385
	var v386 int32
	_ = v386
	var v387 base.V128
	_ = v387
	var v388 base.V128
	_ = v388
	var v390 base.V128
	_ = v390
	var v391 base.V128
	_ = v391
	var v393 base.V128
	_ = v393
	var v394 base.V128
	_ = v394
	var v396 base.V128
	_ = v396
	var v398 base.V128
	_ = v398
	var v400 base.V128
	_ = v400
	var v406 base.V128
	_ = v406
	var v407 base.V128
	_ = v407
	var v413 base.V128
	_ = v413
	var v414 base.V128
	_ = v414
	var v415 base.V128
	_ = v415
	var v416 base.V128
	_ = v416
	var v419 base.V128
	_ = v419
	var v420 base.V128
	_ = v420
	var v423 base.V128
	_ = v423
	var v424 base.V128
	_ = v424
	var v426 base.V128
	_ = v426
	var v430 base.V128
	_ = v430
	var v433 int32
	_ = v433
	var v435 base.V128
	_ = v435
	var v436 int32
	_ = v436
	var v438 base.V128
	_ = v438
	var v440 base.V128
	_ = v440
	var v441 base.V128
	_ = v441
	var v443 base.V128
	_ = v443
	var v444 base.V128
	_ = v444
	var v446 base.V128
	_ = v446
	var v447 base.V128
	_ = v447
	var v449 base.V128
	_ = v449
	var v452 base.V128
	_ = v452
	var v458 base.V128
	_ = v458
	var v459 base.V128
	_ = v459
	var v465 base.V128
	_ = v465
	var v466 base.V128
	_ = v466
	var v468 base.V128
	_ = v468
	var v472 base.V128
	_ = v472
	var v475 base.V128
	_ = v475
	var v476 base.V128
	_ = v476
	var v478 base.V128
	_ = v478
	var v482 base.V128
	_ = v482
	var v499 int32
	_ = v499
	var v500 int32
	_ = v500
	var v510 base.V128
	_ = v510
	var v512 int32
	_ = v512
	var v513 base.V128
	_ = v513
	var v514 base.V128
	_ = v514
	var v515 base.V128
	_ = v515
	var v517 base.V128
	_ = v517
	var v521 base.V128
	_ = v521
	var v522 base.V128
	_ = v522
	var v525 base.V128
	_ = v525
	var v527 base.V128
	_ = v527
	var v528 base.V128
	_ = v528
	var v529 base.V128
	_ = v529
	var v531 base.V128
	_ = v531
	var v537 base.V128
	_ = v537
	var v539 base.V128
	_ = v539
	var v540 base.V128
	_ = v540
	var v541 base.V128
	_ = v541
	var v543 base.V128
	_ = v543
	var v551 int32
	_ = v551
	var v553 base.V128
	_ = v553
	var v554 base.V128
	_ = v554
	var v564 base.V128
	_ = v564
	var v574 base.V128
	_ = v574
	var v576 base.V128
	_ = v576
	var v578 base.V128
	_ = v578
	var v580 base.V128
	_ = v580
	var v584 base.V128
	_ = v584
	var v609 int32
	_ = v609
	var v610 int32
	_ = v610
	var v620 base.V128
	_ = v620
	var v622 int32
	_ = v622
	var v623 base.V128
	_ = v623
	var v624 base.V128
	_ = v624
	var v625 base.V128
	_ = v625
	var v627 base.V128
	_ = v627
	var v631 base.V128
	_ = v631
	var v632 base.V128
	_ = v632
	var v635 base.V128
	_ = v635
	var v637 base.V128
	_ = v637
	var v638 base.V128
	_ = v638
	var v639 base.V128
	_ = v639
	var v641 base.V128
	_ = v641
	var v647 base.V128
	_ = v647
	var v649 base.V128
	_ = v649
	var v650 base.V128
	_ = v650
	var v651 base.V128
	_ = v651
	var v653 base.V128
	_ = v653
	var v661 int32
	_ = v661
	var v663 base.V128
	_ = v663
	var v664 base.V128
	_ = v664
	var v674 base.V128
	_ = v674
	var v684 base.V128
	_ = v684
	var v686 base.V128
	_ = v686
	var v688 base.V128
	_ = v688
	var v690 base.V128
	_ = v690
	var v694 base.V128
	_ = v694
	var v703 int32
	_ = v703
	var v708 int32
	_ = v708
	var v726 int32
	_ = v726
	var v727 int32
	_ = v727
	var v750 int32
	_ = v750
	var v753 int32
	_ = v753
	var v759 int32
	_ = v759
	var v760 int32
	_ = v760
	var v762 int32
	_ = v762
	var v764 int32
	_ = v764
	var v765 int32
	_ = v765
	var v769 int32
	_ = v769
	var v770 int32
	_ = v770
	var v771 int32
	_ = v771
	var v773 int32
	_ = v773
	var v784 int32
	_ = v784
	var v809 int32
	_ = v809
	var v810 int32
	_ = v810
	var v814 int32
	_ = v814
	var v818 int32
	_ = v818
	var v819 int32
	_ = v819
	var v854 int32
	_ = v854
	var v855 int32
	_ = v855
	var v860 int64
	_ = v860
	var v863 int32
	_ = v863
	var v864 int32
	_ = v864
	var v876 int32
	_ = v876
	var v894 int32
	_ = v894
	var v895 int32
	_ = v895
	var v896 int32
	_ = v896
	var v907 int32
	_ = v907
	var v932 int32
	_ = v932
	var v933 int32
	_ = v933
	var v937 int32
	_ = v937
	var v941 int32
	_ = v941
	var v942 int32
	_ = v942
	var v977 int32
	_ = v977
	var v978 int32
	_ = v978
	var v983 int64
	_ = v983
	var v986 int32
	_ = v986
	var v987 int32
	_ = v987
	var v999 int32
	_ = v999
	var v1017 int32
	_ = v1017
	var v1018 base.V128
	_ = v1018
	var v1019 int32
	_ = v1019
	var v1020 base.V128
	_ = v1020
	var v1022 base.V128
	_ = v1022
	var v1023 base.V128
	_ = v1023
	var v1024 int32
	_ = v1024
	var v1025 base.V128
	_ = v1025
	var v1026 base.V128
	_ = v1026
	var v1028 base.V128
	_ = v1028
	var v1029 base.V128
	_ = v1029
	var v1031 base.V128
	_ = v1031
	var v1033 base.V128
	_ = v1033
	var v1035 base.V128
	_ = v1035
	var v1041 base.V128
	_ = v1041
	var v1042 base.V128
	_ = v1042
	var v1048 base.V128
	_ = v1048
	var v1049 base.V128
	_ = v1049
	var v1050 base.V128
	_ = v1050
	var v1051 base.V128
	_ = v1051
	var v1054 base.V128
	_ = v1054
	var v1055 base.V128
	_ = v1055
	var v1058 base.V128
	_ = v1058
	var v1059 base.V128
	_ = v1059
	var v1061 base.V128
	_ = v1061
	var v1065 base.V128
	_ = v1065
	var v1071 int32
	_ = v1071
	var v1073 int32
	_ = v1073
	var v1074 int32
	_ = v1074
	var v1085 int32
	_ = v1085
	var v1110 int32
	_ = v1110
	var v1111 int32
	_ = v1111
	var v1115 int32
	_ = v1115
	var v1119 int32
	_ = v1119
	var v1120 int32
	_ = v1120
	var v1155 int32
	_ = v1155
	var v1156 int32
	_ = v1156
	var v1161 int64
	_ = v1161
	var v1164 int32
	_ = v1164
	var v1165 int32
	_ = v1165
	var v1177 int32
	_ = v1177
	var v1195 int32
	_ = v1195
	var v1206 int32
	_ = v1206
	var v1231 int32
	_ = v1231
	var v1232 int32
	_ = v1232
	var v1236 int32
	_ = v1236
	var v1240 int32
	_ = v1240
	var v1241 int32
	_ = v1241
	var v1276 int32
	_ = v1276
	var v1277 int32
	_ = v1277
	var v1282 int64
	_ = v1282
	var v1285 int32
	_ = v1285
	var v1286 int32
	_ = v1286
	var v1298 int32
	_ = v1298
	var v1317 base.V128
	_ = v1317
	var v1318 int32
	_ = v1318
	var v1319 base.V128
	_ = v1319
	var v1321 base.V128
	_ = v1321
	var v1322 base.V128
	_ = v1322
	var v1324 base.V128
	_ = v1324
	var v1325 base.V128
	_ = v1325
	var v1327 base.V128
	_ = v1327
	var v1328 base.V128
	_ = v1328
	var v1330 base.V128
	_ = v1330
	var v1333 base.V128
	_ = v1333
	var v1339 base.V128
	_ = v1339
	var v1340 base.V128
	_ = v1340
	var v1346 base.V128
	_ = v1346
	var v1347 base.V128
	_ = v1347
	var v1348 base.V128
	_ = v1348
	var v1349 base.V128
	_ = v1349
	var v1352 base.V128
	_ = v1352
	var v1353 base.V128
	_ = v1353
	var v1356 base.V128
	_ = v1356
	var v1357 base.V128
	_ = v1357
	var v1359 base.V128
	_ = v1359
	var v1363 base.V128
	_ = v1363
	var v1369 int32
	_ = v1369
	var v1370 int32
	_ = v1370
	var v1372 int32
	_ = v1372
	var v1386 int32
	_ = v1386
	var v1387 int32
	_ = v1387
	var v1397 base.V128
	_ = v1397
	var v1399 int32
	_ = v1399
	var v1400 base.V128
	_ = v1400
	var v1401 base.V128
	_ = v1401
	var v1402 base.V128
	_ = v1402
	var v1404 base.V128
	_ = v1404
	var v1408 base.V128
	_ = v1408
	var v1409 base.V128
	_ = v1409
	var v1412 base.V128
	_ = v1412
	var v1414 base.V128
	_ = v1414
	var v1415 base.V128
	_ = v1415
	var v1416 base.V128
	_ = v1416
	var v1418 base.V128
	_ = v1418
	var v1424 base.V128
	_ = v1424
	var v1426 base.V128
	_ = v1426
	var v1427 base.V128
	_ = v1427
	var v1428 base.V128
	_ = v1428
	var v1430 base.V128
	_ = v1430
	var v1438 int32
	_ = v1438
	var v1440 base.V128
	_ = v1440
	var v1441 base.V128
	_ = v1441
	var v1451 base.V128
	_ = v1451
	var v1461 base.V128
	_ = v1461
	var v1463 base.V128
	_ = v1463
	var v1465 base.V128
	_ = v1465
	var v1467 base.V128
	_ = v1467
	var v1471 base.V128
	_ = v1471
	var v1480 int32
	_ = v1480
	var v1485 int32
	_ = v1485
	var v1489 int32
	_ = v1489
	var v1503 int32
	_ = v1503
	var v1504 int32
	_ = v1504
	var v1514 base.V128
	_ = v1514
	var v1516 int32
	_ = v1516
	var v1517 base.V128
	_ = v1517
	var v1518 base.V128
	_ = v1518
	var v1519 base.V128
	_ = v1519
	var v1521 base.V128
	_ = v1521
	var v1525 base.V128
	_ = v1525
	var v1526 base.V128
	_ = v1526
	var v1529 base.V128
	_ = v1529
	var v1531 base.V128
	_ = v1531
	var v1532 base.V128
	_ = v1532
	var v1533 base.V128
	_ = v1533
	var v1535 base.V128
	_ = v1535
	var v1541 base.V128
	_ = v1541
	var v1543 base.V128
	_ = v1543
	var v1544 base.V128
	_ = v1544
	var v1545 base.V128
	_ = v1545
	var v1547 base.V128
	_ = v1547
	var v1555 int32
	_ = v1555
	var v1557 base.V128
	_ = v1557
	var v1558 base.V128
	_ = v1558
	var v1568 base.V128
	_ = v1568
	var v1578 base.V128
	_ = v1578
	var v1580 base.V128
	_ = v1580
	var v1582 base.V128
	_ = v1582
	var v1584 base.V128
	_ = v1584
	var v1588 base.V128
	_ = v1588
	var v1602 int32
	_ = v1602
	var v1616 int32
	_ = v1616
	var v1617 int32
	_ = v1617
	var v1627 base.V128
	_ = v1627
	var v1629 int32
	_ = v1629
	var v1630 base.V128
	_ = v1630
	var v1631 base.V128
	_ = v1631
	var v1632 base.V128
	_ = v1632
	var v1634 base.V128
	_ = v1634
	var v1638 base.V128
	_ = v1638
	var v1639 base.V128
	_ = v1639
	var v1642 base.V128
	_ = v1642
	var v1644 base.V128
	_ = v1644
	var v1645 base.V128
	_ = v1645
	var v1646 base.V128
	_ = v1646
	var v1648 base.V128
	_ = v1648
	var v1654 base.V128
	_ = v1654
	var v1656 base.V128
	_ = v1656
	var v1657 base.V128
	_ = v1657
	var v1658 base.V128
	_ = v1658
	var v1660 base.V128
	_ = v1660
	var v1668 int32
	_ = v1668
	var v1670 base.V128
	_ = v1670
	var v1671 base.V128
	_ = v1671
	var v1681 base.V128
	_ = v1681
	var v1691 base.V128
	_ = v1691
	var v1693 base.V128
	_ = v1693
	var v1695 base.V128
	_ = v1695
	var v1697 base.V128
	_ = v1697
	var v1701 base.V128
	_ = v1701
	var v1710 int32
	_ = v1710
	var v1711 int32
	_ = v1711
	var v1714 int32
	_ = v1714
	var v1715 int32
	_ = v1715
	var v1717 int32
	_ = v1717
	v10 = int32(0)
	v33 = m.G0
	v35 = v33 - int32(528)
	m.G0 = v35
	base.MemoryFill(m, v35+int32(64), v10, int32(463))
	v162 = int32(255)
	*(*uint8)(unsafe.Add(mBase, uint32(l6))) = uint8(v162)
	v164 = int32(1)
	v165 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v166 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2))))
	v171 = int32(base.Ui32(v165+v166)>>(uint(v164)%32)) + v164
	v174 = int32(base.Ui32(v171+v166) >> (uint(v164) % 32))
	v177 = int32(8)
	v179 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v183 = int32(base.Ui32(v179*int32(_a_F_UpsampleArgbLinePair_SSE2_0)) >> (uint(v177) % 32))
	v184 = int32(base.Ui32(v174*int32(_a_F_UpsampleArgbLinePair_SSE2_1))>>(uint(v177)%32)) + v183
	v186 = v184 + int32(-17685)
	if base.Ui32(v184) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_SSE2_2)) {
		v193 = int32(0)
	} else {
		v193 = v162
	}
	if base.Ui32(v186) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_SSE2_3)) {
		v196 = int32(base.Ui32(v186) >> (uint(int32(6)) % 32))
	} else {
		v196 = v193
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+3)) = uint8(v196)
	v198 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5))))
	v199 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3))))
	v201 = int32(1)
	v204 = int32(base.Ui32(v198+v199)>>(uint(v201)%32)) + v201
	v207 = int32(base.Ui32(v204+v199) >> (uint(v201) % 32))
	v212 = int32(base.Ui32(v207*int32(_a_F_UpsampleArgbLinePair_SSE2_4))>>(uint(int32(8))%32)) + v183
	v214 = v212 + int32(-14234)
	if base.Ui32(v212) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_SSE2_5)) {
		v221 = int32(0)
	} else {
		v221 = int32(255)
	}
	if base.Ui32(v214) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_SSE2_3)) {
		v224 = int32(base.Ui32(v214) >> (uint(int32(6)) % 32))
	} else {
		v224 = v221
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+1)) = uint8(v224)
	v228 = int32(8)
	v235 = v183 - (int32(base.Ui32(v174*int32(_a_F_UpsampleArgbLinePair_SSE2_6))>>(uint(v228)%32)) + int32(base.Ui32(v207*int32(_a_F_UpsampleArgbLinePair_SSE2_7))>>(uint(v228)%32)))
	v237 = v235 + int32(_a_F_UpsampleArgbLinePair_SSE2_8)
	if v235 < int32(-8708) {
		v244 = int32(0)
	} else {
		v244 = int32(255)
	}
	if base.Ui32(v237) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_SSE2_3)) {
		v247 = int32(base.Ui32(v237) >> (uint(int32(6)) % 32))
	} else {
		v247 = v244
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+2)) = uint8(v247)
	if l1 == int32(0) {
	} else {
		v251 = int32(255)
		*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v251)
		v253 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
		v256 = int32(8)
		v257 = int32(base.Ui32(v253*int32(_a_F_UpsampleArgbLinePair_SSE2_0)) >> (uint(v256) % 32))
		v260 = int32(base.Ui32(v171+v165) >> (uint(int32(1)) % 32))
		v265 = v257 + int32(base.Ui32(v260*int32(_a_F_UpsampleArgbLinePair_SSE2_1))>>(uint(v256)%32))
		v267 = v265 + int32(-17685)
		if base.Ui32(v265) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_SSE2_2)) {
			v274 = int32(0)
		} else {
			v274 = v251
		}
		if base.Ui32(v267) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_SSE2_3)) {
			v277 = int32(base.Ui32(v267) >> (uint(int32(6)) % 32))
		} else {
			v277 = v274
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+3)) = uint8(v277)
		v281 = int32(base.Ui32(v204+v198) >> (uint(int32(1)) % 32))
		v286 = v257 + int32(base.Ui32(v281*int32(_a_F_UpsampleArgbLinePair_SSE2_4))>>(uint(int32(8))%32))
		v288 = v286 + int32(-14234)
		if base.Ui32(v286) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_SSE2_5)) {
			v295 = int32(0)
		} else {
			v295 = int32(255)
		}
		if base.Ui32(v288) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_SSE2_3)) {
			v298 = int32(base.Ui32(v288) >> (uint(int32(6)) % 32))
		} else {
			v298 = v295
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+1)) = uint8(v298)
		v302 = int32(8)
		v309 = v257 - (int32(base.Ui32(v281*int32(_a_F_UpsampleArgbLinePair_SSE2_7))>>(uint(v302)%32)) + int32(base.Ui32(v260*int32(_a_F_UpsampleArgbLinePair_SSE2_6))>>(uint(v302)%32)))
		v311 = v309 + int32(_a_F_UpsampleArgbLinePair_SSE2_8)
		if v309 < int32(-8708) {
			v318 = int32(0)
		} else {
			v318 = int32(255)
		}
		if base.Ui32(v311) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_SSE2_3)) {
			v321 = int32(base.Ui32(v311) >> (uint(int32(6)) % 32))
		} else {
			v321 = v318
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+2)) = uint8(v321)
	}
	v329 = v35 + int32(96)
	v331 = v35 + int32(64)
	if l8 < int32(34) {
		v726 = v10
		v727 = v164
	} else {
		v334 = int32(4)
		v338 = int32(1)
		v346 = int32(0)
		v358 = v346
		v359 = l7 + v334
		v363 = v346
		v364 = l6 + v334
		for {
			v380 = l4 + v358
			v381 = int32(0)
			v382 = base.Simd_g_v128_load_rng(m, v380, v381, int32(0), int32(17))
			v383 = l2 + v358
			v385 = base.Simd_g_v128_load_rng(m, v383, v381, int32(0), int32(17))
			v386 = int32(1)
			v387 = base.Simd_g_v128_load_nc(m, v380, v386)
			v388 = base.Simd_g_i8x16_avgr_u(v385, v387)
			v390 = base.Simd_g_v128_load_nc(m, v383, v386)
			v391 = base.Simd_g_i8x16_avgr_u(v382, v390)
			v393 = base.Simd_g_v128_xor(v387, v385)
			v394 = base.Simd_g_v128_xor(v382, v390)
			v396 = base.Simd_g_v128_xor(v388, v391)
			v398 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k0)
			v400 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v388, v391), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v393, v394), v396), v398))
			v406 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v400, v388), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v400, v388), base.Simd_g_v128_and(v396, v393)), v398))
			v407 = base.Simd_g_i8x16_avgr_u(v382, v406)
			v413 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v400, v391), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v400, v391), base.Simd_g_v128_and(v396, v394)), v398))
			v414 = base.Simd_g_i8x16_avgr_u(v387, v413)
			v415 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k1)
			v416 = base.Simd_g_i8x16_shuffle2(v407, v414, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v331, int32(80), v416)
			v419 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k4)
			v420 = base.Simd_g_i8x16_shuffle2(v407, v414, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v331, int32(64), v420)
			v423 = base.Simd_g_i8x16_avgr_u(v385, v413)
			v424 = base.Simd_g_i8x16_avgr_u(v390, v406)
			v426 = base.Simd_g_i8x16_shuffle2(v423, v424, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v331, int32(16), v426)
			v430 = base.Simd_g_i8x16_shuffle2(v423, v424, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v331, v381, v430)
			v433 = l5 + v358
			v435 = base.Simd_g_v128_load_rng(m, v433, v381, int32(0), int32(17))
			v436 = l3 + v358
			v438 = base.Simd_g_v128_load_rng(m, v436, v381, int32(0), int32(17))
			v440 = base.Simd_g_v128_load_nc(m, v433, v386)
			v441 = base.Simd_g_i8x16_avgr_u(v438, v440)
			v443 = base.Simd_g_v128_load_nc(m, v436, v386)
			v444 = base.Simd_g_i8x16_avgr_u(v435, v443)
			v446 = base.Simd_g_v128_xor(v440, v438)
			v447 = base.Simd_g_v128_xor(v435, v443)
			v449 = base.Simd_g_v128_xor(v441, v444)
			v452 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v441, v444), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v446, v447), v449), v398))
			v458 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v452, v441), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v452, v441), base.Simd_g_v128_and(v449, v446)), v398))
			v459 = base.Simd_g_i8x16_avgr_u(v435, v458)
			v465 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v452, v444), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v452, v444), base.Simd_g_v128_and(v449, v447)), v398))
			v466 = base.Simd_g_i8x16_avgr_u(v440, v465)
			v468 = base.Simd_g_i8x16_shuffle2(v459, v466, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v331, int32(112), v468)
			v472 = base.Simd_g_i8x16_shuffle2(v459, v466, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v331, int32(96), v472)
			v475 = base.Simd_g_i8x16_avgr_u(v438, v465)
			v476 = base.Simd_g_i8x16_avgr_u(v443, v458)
			v478 = base.Simd_g_i8x16_shuffle2(v475, v476, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v331, int32(48), v478)
			v482 = base.Simd_g_i8x16_shuffle2(v475, v476, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v331, int32(32), v482)
			v499 = v364
			v500 = v381
			for {
				v510 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k7)
				v512 = int32(0)
				v513 = base.Simd_g_v128_load64_zero(m, l0+v338+v363+v500, v512)
				v514 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k4)
				v515 = base.Simd_g_i8x16_shuffle2(v510, v513, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
				v517 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k8)
				v521 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k9)
				v522 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v515), v517), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v515), v517), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11))
				v525 = base.Simd_g_v128_load64_zero(m, v331+v500, v512)
				v527 = base.Simd_g_i8x16_shuffle2(v510, v525, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
				v528 = base.Simd_g_i32x4_extend_low_i16x8_u(v527)
				v529 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k12)
				v531 = base.Simd_g_i32x4_extend_high_i16x8_u(v527)
				v537 = base.Simd_g_v128_load64_zero(m, v329+v500, v512)
				v539 = base.Simd_g_i8x16_shuffle2(v510, v537, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
				v540 = base.Simd_g_i32x4_extend_low_i16x8_u(v539)
				v541 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k13)
				v543 = base.Simd_g_i32x4_extend_high_i16x8_u(v539)
				v551 = int32(6)
				v553 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k14), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v522, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v528, v529), base.Simd_g_i32x4_mul(v531, v529), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v540, v541), base.Simd_g_i32x4_mul(v543, v541), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k15)), v551))
				v554 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k16)
				v564 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k17)
				v574 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v522, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v540, v554), base.Simd_g_i32x4_mul(v543, v554), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k18)), v551), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v528, v564), base.Simd_g_i32x4_mul(v531, v564), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11)), v522), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k19)), v551))
				v576 = base.Simd_g_i8x16_shuffle2(v553, v574, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
				v578 = base.Simd_g_i8x16_shuffle2(v553, v574, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k3))
				v580 = base.Simd_g_i8x16_shuffle2(v576, v578, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k20), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k21))
				base.Simd_g_v128_store(m, v499, int32(16), v580)
				v584 = base.Simd_g_i8x16_shuffle2(v576, v578, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k22), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k23))
				base.Simd_g_v128_store(m, v499, v512, v584)
				if base.Ui32(v500) < base.Ui32(int32(24)) {
					v499 = v499 + int32(32)
					v500 = v500 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			if l1 == int32(0) {
			} else {
				v609 = v359
				v610 = int32(0)
				for {
					v620 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k7)
					v622 = int32(0)
					v623 = base.Simd_g_v128_load64_zero(m, l1+v338+v363+v610, v622)
					v624 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k4)
					v625 = base.Simd_g_i8x16_shuffle2(v620, v623, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
					v627 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k8)
					v631 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k9)
					v632 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v625), v627), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v625), v627), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11))
					v635 = base.Simd_g_v128_load64_zero(m, v35+int32(128)+v610, v622)
					v637 = base.Simd_g_i8x16_shuffle2(v620, v635, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
					v638 = base.Simd_g_i32x4_extend_low_i16x8_u(v637)
					v639 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k12)
					v641 = base.Simd_g_i32x4_extend_high_i16x8_u(v637)
					v647 = base.Simd_g_v128_load64_zero(m, v35+int32(160)+v610, v622)
					v649 = base.Simd_g_i8x16_shuffle2(v620, v647, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
					v650 = base.Simd_g_i32x4_extend_low_i16x8_u(v649)
					v651 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k13)
					v653 = base.Simd_g_i32x4_extend_high_i16x8_u(v649)
					v661 = int32(6)
					v663 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k14), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v632, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v638, v639), base.Simd_g_i32x4_mul(v641, v639), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v650, v651), base.Simd_g_i32x4_mul(v653, v651), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k15)), v661))
					v664 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k16)
					v674 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k17)
					v684 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v632, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v650, v664), base.Simd_g_i32x4_mul(v653, v664), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k18)), v661), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v638, v674), base.Simd_g_i32x4_mul(v641, v674), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11)), v632), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k19)), v661))
					v686 = base.Simd_g_i8x16_shuffle2(v663, v684, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
					v688 = base.Simd_g_i8x16_shuffle2(v663, v684, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k3))
					v690 = base.Simd_g_i8x16_shuffle2(v686, v688, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k20), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k21))
					base.Simd_g_v128_store(m, v609, int32(16), v690)
					v694 = base.Simd_g_i8x16_shuffle2(v686, v688, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k22), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k23))
					base.Simd_g_v128_store(m, v609, v622, v694)
					if base.Ui32(v610) < base.Ui32(int32(24)) {
						v609 = v609 + int32(32)
						v610 = v610 + int32(8)
						continue
					} else {
						break
					}
					break
				}
			}
			v703 = int32(128)
			v708 = v358 + int32(16)
			if v363+int32(66) <= l8 {
				v358 = v708
				v359 = v359 + v703
				v363 = v363 + int32(32)
				v364 = v364 + v703
				continue
			} else {
				break
			}
			break
		}
		v726 = v708
		v727 = v363 + int32(33)
	}
	if l8 < int32(2) {
	} else {
		v750 = int32(32)
		v753 = int32(1)
		v759 = int32(base.Ui32(l8+v753)>>(uint(v753)%32)) - int32(base.Ui32(v727)>>(uint(v753)%32))
		v760 = F_memcpy(m, v35+v750, l2+v726, v759)
		mBase = m.M
		v762 = F_memcpy(m, v35, l4+v726, v759)
		mBase = m.M
		v764 = v762 + v750
		v765 = v764 + v759
		v769 = v759 + int32(-1)
		v770 = v764 + v769
		v771 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v770))))
		v773 = int32(17) - v759
		if base.Ui32(v773) < base.Ui32(int32(33)) {
			if v773 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v765))) = uint8(v771)
				v784 = v765 + v773
				*(*uint8)(unsafe.Add(mBase, uint32(v784+int32(-1)))) = uint8(v771)
				if base.Ui32(v773) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v765)+2)) = uint8(v771)
					*(*uint8)(unsafe.Add(mBase, uint32(v765)+1)) = uint8(v771)
					*(*uint8)(unsafe.Add(mBase, uint32(v784+int32(-3)))) = uint8(v771)
					*(*uint8)(unsafe.Add(mBase, uint32(v784+int32(-2)))) = uint8(v771)
					if base.Ui32(v773) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v765)+3)) = uint8(v771)
						*(*uint8)(unsafe.Add(mBase, uint32(v784+int32(-4)))) = uint8(v771)
						if base.Ui32(v773) < base.Ui32(int32(9)) {
						} else {
							v809 = (int32(0) - v765) & int32(3)
							v810 = v765 + v809
							v814 = v771 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v810))) = v814
							v818 = (v773 - v809) & int32(60)
							v819 = v810 + v818
							*(*int32)(unsafe.Add(mBase, uint32(v819+int32(-4)))) = v814
							if base.Ui32(v818) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v810)+8)) = v814
								*(*int32)(unsafe.Add(mBase, uint32(v810)+4)) = v814
								*(*int32)(unsafe.Add(mBase, uint32(v819+int32(-8)))) = v814
								*(*int32)(unsafe.Add(mBase, uint32(v819+int32(-12)))) = v814
								if base.Ui32(v818) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v810)+24)) = v814
									*(*int32)(unsafe.Add(mBase, uint32(v810)+20)) = v814
									*(*int32)(unsafe.Add(mBase, uint32(v810)+16)) = v814
									*(*int32)(unsafe.Add(mBase, uint32(v810)+12)) = v814
									*(*int32)(unsafe.Add(mBase, uint32(v819+int32(-16)))) = v814
									*(*int32)(unsafe.Add(mBase, uint32(v819+int32(-20)))) = v814
									*(*int32)(unsafe.Add(mBase, uint32(v819+int32(-24)))) = v814
									*(*int32)(unsafe.Add(mBase, uint32(v819+int32(-28)))) = v814
									v854 = v810&int32(4) | int32(24)
									v855 = v818 - v854
									if base.Ui32(v855) < base.Ui32(int32(32)) {
									} else {
										v860 = base.I64_extend_i32_u(v814) * int64(4294967297)
										v863 = v855
										v864 = v810 + v854
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v864)+24)) = v860
											*(*int64)(unsafe.Add(mBase, uint32(v864)+16)) = v860
											*(*int64)(unsafe.Add(mBase, uint32(v864)+8)) = v860
											*(*int64)(unsafe.Add(mBase, uint32(v864))) = v860
											v876 = v863 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v876) {
												v863 = v876
												v864 = v864 + int32(32)
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
			base.MemoryFill(m, v765, v771, v773)
		}
		v894 = v762 + v759
		v895 = v762 + v769
		v896 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v895))))
		if base.Ui32(v773) < base.Ui32(int32(33)) {
			if v773 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v894))) = uint8(v896)
				v907 = v894 + v773
				*(*uint8)(unsafe.Add(mBase, uint32(v907+int32(-1)))) = uint8(v896)
				if base.Ui32(v773) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v894)+2)) = uint8(v896)
					*(*uint8)(unsafe.Add(mBase, uint32(v894)+1)) = uint8(v896)
					*(*uint8)(unsafe.Add(mBase, uint32(v907+int32(-3)))) = uint8(v896)
					*(*uint8)(unsafe.Add(mBase, uint32(v907+int32(-2)))) = uint8(v896)
					if base.Ui32(v773) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v894)+3)) = uint8(v896)
						*(*uint8)(unsafe.Add(mBase, uint32(v907+int32(-4)))) = uint8(v896)
						if base.Ui32(v773) < base.Ui32(int32(9)) {
						} else {
							v932 = (int32(0) - v894) & int32(3)
							v933 = v894 + v932
							v937 = v896 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v933))) = v937
							v941 = (v773 - v932) & int32(60)
							v942 = v933 + v941
							*(*int32)(unsafe.Add(mBase, uint32(v942+int32(-4)))) = v937
							if base.Ui32(v941) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v933)+8)) = v937
								*(*int32)(unsafe.Add(mBase, uint32(v933)+4)) = v937
								*(*int32)(unsafe.Add(mBase, uint32(v942+int32(-8)))) = v937
								*(*int32)(unsafe.Add(mBase, uint32(v942+int32(-12)))) = v937
								if base.Ui32(v941) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v933)+24)) = v937
									*(*int32)(unsafe.Add(mBase, uint32(v933)+20)) = v937
									*(*int32)(unsafe.Add(mBase, uint32(v933)+16)) = v937
									*(*int32)(unsafe.Add(mBase, uint32(v933)+12)) = v937
									*(*int32)(unsafe.Add(mBase, uint32(v942+int32(-16)))) = v937
									*(*int32)(unsafe.Add(mBase, uint32(v942+int32(-20)))) = v937
									*(*int32)(unsafe.Add(mBase, uint32(v942+int32(-24)))) = v937
									*(*int32)(unsafe.Add(mBase, uint32(v942+int32(-28)))) = v937
									v977 = v933&int32(4) | int32(24)
									v978 = v941 - v977
									if base.Ui32(v978) < base.Ui32(int32(32)) {
									} else {
										v983 = base.I64_extend_i32_u(v937) * int64(4294967297)
										v986 = v978
										v987 = v933 + v977
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v987)+24)) = v983
											*(*int64)(unsafe.Add(mBase, uint32(v987)+16)) = v983
											*(*int64)(unsafe.Add(mBase, uint32(v987)+8)) = v983
											*(*int64)(unsafe.Add(mBase, uint32(v987))) = v983
											v999 = v986 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v999) {
												v986 = v999
												v987 = v987 + int32(32)
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
			base.MemoryFill(m, v894, v896, v773)
		}
		v1017 = int32(0)
		v1018 = base.Simd_g_v128_load_rng(m, v762, v1017, int32(0), int32(49))
		v1019 = int32(32)
		v1020 = base.Simd_g_v128_load_nc(m, v762, v1019)
		v1022 = base.Simd_g_v128_load_nc(m, v762, int32(1))
		v1023 = base.Simd_g_i8x16_avgr_u(v1020, v1022)
		v1024 = int32(33)
		v1025 = base.Simd_g_v128_load_nc(m, v762, v1024)
		v1026 = base.Simd_g_i8x16_avgr_u(v1018, v1025)
		v1028 = base.Simd_g_v128_xor(v1022, v1020)
		v1029 = base.Simd_g_v128_xor(v1018, v1025)
		v1031 = base.Simd_g_v128_xor(v1023, v1026)
		v1033 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k0)
		v1035 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1023, v1026), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v1028, v1029), v1031), v1033))
		v1041 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1035, v1023), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1035, v1023), base.Simd_g_v128_and(v1031, v1028)), v1033))
		v1042 = base.Simd_g_i8x16_avgr_u(v1018, v1041)
		v1048 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1035, v1026), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1035, v1026), base.Simd_g_v128_and(v1031, v1029)), v1033))
		v1049 = base.Simd_g_i8x16_avgr_u(v1022, v1048)
		v1050 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k1)
		v1051 = base.Simd_g_i8x16_shuffle2(v1042, v1049, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v331, int32(80), v1051)
		v1054 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k4)
		v1055 = base.Simd_g_i8x16_shuffle2(v1042, v1049, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v331, int32(64), v1055)
		v1058 = base.Simd_g_i8x16_avgr_u(v1020, v1048)
		v1059 = base.Simd_g_i8x16_avgr_u(v1025, v1041)
		v1061 = base.Simd_g_i8x16_shuffle2(v1058, v1059, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v331, int32(16), v1061)
		v1065 = base.Simd_g_i8x16_shuffle2(v1058, v1059, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v331, v1017, v1065)
		v1071 = F_memcpy(m, v762+v1019, l3+v726, v759)
		mBase = m.M
		v1073 = F_memcpy(m, v762, l5+v726, v759)
		mBase = m.M
		v1074 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v770))))
		if base.Ui32(v773) < base.Ui32(v1024) {
			if v773 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v765))) = uint8(v1074)
				v1085 = v765 + v773
				*(*uint8)(unsafe.Add(mBase, uint32(v1085+int32(-1)))) = uint8(v1074)
				if base.Ui32(v773) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v765)+2)) = uint8(v1074)
					*(*uint8)(unsafe.Add(mBase, uint32(v765)+1)) = uint8(v1074)
					*(*uint8)(unsafe.Add(mBase, uint32(v1085+int32(-3)))) = uint8(v1074)
					*(*uint8)(unsafe.Add(mBase, uint32(v1085+int32(-2)))) = uint8(v1074)
					if base.Ui32(v773) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v765)+3)) = uint8(v1074)
						*(*uint8)(unsafe.Add(mBase, uint32(v1085+int32(-4)))) = uint8(v1074)
						if base.Ui32(v773) < base.Ui32(int32(9)) {
						} else {
							v1110 = (int32(0) - v765) & int32(3)
							v1111 = v765 + v1110
							v1115 = v1074 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1111))) = v1115
							v1119 = (v773 - v1110) & int32(60)
							v1120 = v1111 + v1119
							*(*int32)(unsafe.Add(mBase, uint32(v1120+int32(-4)))) = v1115
							if base.Ui32(v1119) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1111)+8)) = v1115
								*(*int32)(unsafe.Add(mBase, uint32(v1111)+4)) = v1115
								*(*int32)(unsafe.Add(mBase, uint32(v1120+int32(-8)))) = v1115
								*(*int32)(unsafe.Add(mBase, uint32(v1120+int32(-12)))) = v1115
								if base.Ui32(v1119) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1111)+24)) = v1115
									*(*int32)(unsafe.Add(mBase, uint32(v1111)+20)) = v1115
									*(*int32)(unsafe.Add(mBase, uint32(v1111)+16)) = v1115
									*(*int32)(unsafe.Add(mBase, uint32(v1111)+12)) = v1115
									*(*int32)(unsafe.Add(mBase, uint32(v1120+int32(-16)))) = v1115
									*(*int32)(unsafe.Add(mBase, uint32(v1120+int32(-20)))) = v1115
									*(*int32)(unsafe.Add(mBase, uint32(v1120+int32(-24)))) = v1115
									*(*int32)(unsafe.Add(mBase, uint32(v1120+int32(-28)))) = v1115
									v1155 = v1111&int32(4) | int32(24)
									v1156 = v1119 - v1155
									if base.Ui32(v1156) < base.Ui32(int32(32)) {
									} else {
										v1161 = base.I64_extend_i32_u(v1115) * int64(4294967297)
										v1164 = v1156
										v1165 = v1111 + v1155
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1165)+24)) = v1161
											*(*int64)(unsafe.Add(mBase, uint32(v1165)+16)) = v1161
											*(*int64)(unsafe.Add(mBase, uint32(v1165)+8)) = v1161
											*(*int64)(unsafe.Add(mBase, uint32(v1165))) = v1161
											v1177 = v1164 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1177) {
												v1164 = v1177
												v1165 = v1165 + int32(32)
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
			base.MemoryFill(m, v765, v1074, v773)
		}
		v1195 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v895))))
		if base.Ui32(v773) < base.Ui32(int32(33)) {
			if v773 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v894))) = uint8(v1195)
				v1206 = v894 + v773
				*(*uint8)(unsafe.Add(mBase, uint32(v1206+int32(-1)))) = uint8(v1195)
				if base.Ui32(v773) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v894)+2)) = uint8(v1195)
					*(*uint8)(unsafe.Add(mBase, uint32(v894)+1)) = uint8(v1195)
					*(*uint8)(unsafe.Add(mBase, uint32(v1206+int32(-3)))) = uint8(v1195)
					*(*uint8)(unsafe.Add(mBase, uint32(v1206+int32(-2)))) = uint8(v1195)
					if base.Ui32(v773) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v894)+3)) = uint8(v1195)
						*(*uint8)(unsafe.Add(mBase, uint32(v1206+int32(-4)))) = uint8(v1195)
						if base.Ui32(v773) < base.Ui32(int32(9)) {
						} else {
							v1231 = (int32(0) - v894) & int32(3)
							v1232 = v894 + v1231
							v1236 = v1195 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1232))) = v1236
							v1240 = (v773 - v1231) & int32(60)
							v1241 = v1232 + v1240
							*(*int32)(unsafe.Add(mBase, uint32(v1241+int32(-4)))) = v1236
							if base.Ui32(v1240) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1232)+8)) = v1236
								*(*int32)(unsafe.Add(mBase, uint32(v1232)+4)) = v1236
								*(*int32)(unsafe.Add(mBase, uint32(v1241+int32(-8)))) = v1236
								*(*int32)(unsafe.Add(mBase, uint32(v1241+int32(-12)))) = v1236
								if base.Ui32(v1240) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1232)+24)) = v1236
									*(*int32)(unsafe.Add(mBase, uint32(v1232)+20)) = v1236
									*(*int32)(unsafe.Add(mBase, uint32(v1232)+16)) = v1236
									*(*int32)(unsafe.Add(mBase, uint32(v1232)+12)) = v1236
									*(*int32)(unsafe.Add(mBase, uint32(v1241+int32(-16)))) = v1236
									*(*int32)(unsafe.Add(mBase, uint32(v1241+int32(-20)))) = v1236
									*(*int32)(unsafe.Add(mBase, uint32(v1241+int32(-24)))) = v1236
									*(*int32)(unsafe.Add(mBase, uint32(v1241+int32(-28)))) = v1236
									v1276 = v1232&int32(4) | int32(24)
									v1277 = v1240 - v1276
									if base.Ui32(v1277) < base.Ui32(int32(32)) {
									} else {
										v1282 = base.I64_extend_i32_u(v1236) * int64(4294967297)
										v1285 = v1277
										v1286 = v1232 + v1276
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1286)+24)) = v1282
											*(*int64)(unsafe.Add(mBase, uint32(v1286)+16)) = v1282
											*(*int64)(unsafe.Add(mBase, uint32(v1286)+8)) = v1282
											*(*int64)(unsafe.Add(mBase, uint32(v1286))) = v1282
											v1298 = v1285 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1298) {
												v1285 = v1298
												v1286 = v1286 + int32(32)
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
			base.MemoryFill(m, v894, v1195, v773)
		}
		v1317 = base.Simd_g_v128_load_rng(m, v1073, int32(0), int32(0), int32(49))
		v1318 = int32(32)
		v1319 = base.Simd_g_v128_load_nc(m, v1073, v1318)
		v1321 = base.Simd_g_v128_load_nc(m, v1073, int32(1))
		v1322 = base.Simd_g_i8x16_avgr_u(v1319, v1321)
		v1324 = base.Simd_g_v128_load_nc(m, v1073, int32(33))
		v1325 = base.Simd_g_i8x16_avgr_u(v1317, v1324)
		v1327 = base.Simd_g_v128_xor(v1321, v1319)
		v1328 = base.Simd_g_v128_xor(v1317, v1324)
		v1330 = base.Simd_g_v128_xor(v1322, v1325)
		v1333 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1322, v1325), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v1327, v1328), v1330), v1033))
		v1339 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1333, v1322), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1333, v1322), base.Simd_g_v128_and(v1330, v1327)), v1033))
		v1340 = base.Simd_g_i8x16_avgr_u(v1317, v1339)
		v1346 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1333, v1325), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1333, v1325), base.Simd_g_v128_and(v1330, v1328)), v1033))
		v1347 = base.Simd_g_i8x16_avgr_u(v1321, v1346)
		v1348 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k1)
		v1349 = base.Simd_g_i8x16_shuffle2(v1340, v1347, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v331, int32(112), v1349)
		v1352 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k4)
		v1353 = base.Simd_g_i8x16_shuffle2(v1340, v1347, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v331, int32(96), v1353)
		v1356 = base.Simd_g_i8x16_avgr_u(v1319, v1346)
		v1357 = base.Simd_g_i8x16_avgr_u(v1324, v1339)
		v1359 = base.Simd_g_i8x16_shuffle2(v1356, v1357, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v331, int32(48), v1359)
		v1363 = base.Simd_g_i8x16_shuffle2(v1356, v1357, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v331, v1318, v1363)
		v1369 = l8 - v727
		v1370 = F_memcpy(m, v35+int32(448), l0+v727, v1369)
		mBase = m.M
		v1372 = v35 + int32(192)
		if l1 != 0 {
			v1489 = F_memcpy(m, v35+int32(480), l1+v727, v1369)
			mBase = m.M
			v1503 = v1372
			v1504 = int32(0)
			for {
				v1514 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k7)
				v1516 = int32(0)
				v1517 = base.Simd_g_v128_load64_zero(m, v1370+v1504, v1516)
				v1518 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k4)
				v1519 = base.Simd_g_i8x16_shuffle2(v1514, v1517, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
				v1521 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k8)
				v1525 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k9)
				v1526 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1519), v1521), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1519), v1521), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11))
				v1529 = base.Simd_g_v128_load64_zero(m, v331+v1504, v1516)
				v1531 = base.Simd_g_i8x16_shuffle2(v1514, v1529, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
				v1532 = base.Simd_g_i32x4_extend_low_i16x8_u(v1531)
				v1533 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k12)
				v1535 = base.Simd_g_i32x4_extend_high_i16x8_u(v1531)
				v1541 = base.Simd_g_v128_load64_zero(m, v329+v1504, v1516)
				v1543 = base.Simd_g_i8x16_shuffle2(v1514, v1541, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
				v1544 = base.Simd_g_i32x4_extend_low_i16x8_u(v1543)
				v1545 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k13)
				v1547 = base.Simd_g_i32x4_extend_high_i16x8_u(v1543)
				v1555 = int32(6)
				v1557 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k14), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1526, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1532, v1533), base.Simd_g_i32x4_mul(v1535, v1533), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1544, v1545), base.Simd_g_i32x4_mul(v1547, v1545), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k15)), v1555))
				v1558 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k16)
				v1568 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k17)
				v1578 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1526, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1544, v1558), base.Simd_g_i32x4_mul(v1547, v1558), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k18)), v1555), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1532, v1568), base.Simd_g_i32x4_mul(v1535, v1568), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11)), v1526), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k19)), v1555))
				v1580 = base.Simd_g_i8x16_shuffle2(v1557, v1578, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
				v1582 = base.Simd_g_i8x16_shuffle2(v1557, v1578, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k3))
				v1584 = base.Simd_g_i8x16_shuffle2(v1580, v1582, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k20), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k21))
				base.Simd_g_v128_store(m, v1503, int32(16), v1584)
				v1588 = base.Simd_g_i8x16_shuffle2(v1580, v1582, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k22), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k23))
				base.Simd_g_v128_store(m, v1503, v1516, v1588)
				if base.Ui32(v1504) < base.Ui32(int32(24)) {
					v1503 = v1503 + int32(32)
					v1504 = v1504 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			v1602 = v35 + int32(320)
			v1616 = v1602
			v1617 = int32(0)
			for {
				v1627 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k7)
				v1629 = int32(0)
				v1630 = base.Simd_g_v128_load64_zero(m, v1489+v1617, v1629)
				v1631 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k4)
				v1632 = base.Simd_g_i8x16_shuffle2(v1627, v1630, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
				v1634 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k8)
				v1638 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k9)
				v1639 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1632), v1634), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1632), v1634), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11))
				v1642 = base.Simd_g_v128_load64_zero(m, v35+int32(128)+v1617, v1629)
				v1644 = base.Simd_g_i8x16_shuffle2(v1627, v1642, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
				v1645 = base.Simd_g_i32x4_extend_low_i16x8_u(v1644)
				v1646 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k12)
				v1648 = base.Simd_g_i32x4_extend_high_i16x8_u(v1644)
				v1654 = base.Simd_g_v128_load64_zero(m, v35+int32(160)+v1617, v1629)
				v1656 = base.Simd_g_i8x16_shuffle2(v1627, v1654, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
				v1657 = base.Simd_g_i32x4_extend_low_i16x8_u(v1656)
				v1658 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k13)
				v1660 = base.Simd_g_i32x4_extend_high_i16x8_u(v1656)
				v1668 = int32(6)
				v1670 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k14), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1639, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1645, v1646), base.Simd_g_i32x4_mul(v1648, v1646), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1657, v1658), base.Simd_g_i32x4_mul(v1660, v1658), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k15)), v1668))
				v1671 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k16)
				v1681 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k17)
				v1691 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1639, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1657, v1671), base.Simd_g_i32x4_mul(v1660, v1671), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k18)), v1668), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1645, v1681), base.Simd_g_i32x4_mul(v1648, v1681), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11)), v1639), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k19)), v1668))
				v1693 = base.Simd_g_i8x16_shuffle2(v1670, v1691, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
				v1695 = base.Simd_g_i8x16_shuffle2(v1670, v1691, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k3))
				v1697 = base.Simd_g_i8x16_shuffle2(v1693, v1695, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k20), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k21))
				base.Simd_g_v128_store(m, v1616, int32(16), v1697)
				v1701 = base.Simd_g_i8x16_shuffle2(v1693, v1695, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k22), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k23))
				base.Simd_g_v128_store(m, v1616, v1629, v1701)
				if base.Ui32(v1617) < base.Ui32(int32(24)) {
					v1616 = v1616 + int32(32)
					v1617 = v1617 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			v1710 = int32(2)
			v1711 = v727 << (uint(v1710) % 32)
			v1714 = v1369 << (uint(v1710) % 32)
			v1715 = F_memcpy(m, l6+v1711, v1372, v1714)
			mBase = m.M
			v1717 = F_memcpy(m, l7+v1711, v1602, v1714)
			mBase = m.M
		} else {
			v1386 = v1372
			v1387 = int32(0)
			for {
				v1397 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k7)
				v1399 = int32(0)
				v1400 = base.Simd_g_v128_load64_zero(m, v1370+v1387, v1399)
				v1401 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k4)
				v1402 = base.Simd_g_i8x16_shuffle2(v1397, v1400, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
				v1404 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k8)
				v1408 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k9)
				v1409 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1402), v1404), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1402), v1404), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11))
				v1412 = base.Simd_g_v128_load64_zero(m, v331+v1387, v1399)
				v1414 = base.Simd_g_i8x16_shuffle2(v1397, v1412, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
				v1415 = base.Simd_g_i32x4_extend_low_i16x8_u(v1414)
				v1416 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k12)
				v1418 = base.Simd_g_i32x4_extend_high_i16x8_u(v1414)
				v1424 = base.Simd_g_v128_load64_zero(m, v329+v1387, v1399)
				v1426 = base.Simd_g_i8x16_shuffle2(v1397, v1424, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
				v1427 = base.Simd_g_i32x4_extend_low_i16x8_u(v1426)
				v1428 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k13)
				v1430 = base.Simd_g_i32x4_extend_high_i16x8_u(v1426)
				v1438 = int32(6)
				v1440 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k14), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1409, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1415, v1416), base.Simd_g_i32x4_mul(v1418, v1416), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1427, v1428), base.Simd_g_i32x4_mul(v1430, v1428), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11)))), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k15)), v1438))
				v1441 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k16)
				v1451 = base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k17)
				v1461 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1409, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1427, v1441), base.Simd_g_i32x4_mul(v1430, v1441), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11))), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k18)), v1438), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1415, v1451), base.Simd_g_i32x4_mul(v1418, v1451), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k11)), v1409), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k19)), v1438))
				v1463 = base.Simd_g_i8x16_shuffle2(v1440, v1461, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k6))
				v1465 = base.Simd_g_i8x16_shuffle2(v1440, v1461, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k3))
				v1467 = base.Simd_g_i8x16_shuffle2(v1463, v1465, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k20), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k21))
				base.Simd_g_v128_store(m, v1386, int32(16), v1467)
				v1471 = base.Simd_g_i8x16_shuffle2(v1463, v1465, base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k22), base.Simd_g_const(&F_UpsampleArgbLinePair_SSE2__k23))
				base.Simd_g_v128_store(m, v1386, v1399, v1471)
				if base.Ui32(v1387) < base.Ui32(int32(24)) {
					v1386 = v1386 + int32(32)
					v1387 = v1387 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			v1480 = int32(2)
			v1485 = F_memcpy(m, l6+v727<<(uint(v1480)%32), v1372, v1369<<(uint(v1480)%32))
			mBase = m.M
		}
	}
	m.G0 = v35 + int32(528)
	return
}

var F_UpsampleArgbLinePair_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}
var F_UpsampleArgbLinePair_SSE2__k1 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_UpsampleArgbLinePair_SSE2__k2 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_UpsampleArgbLinePair_SSE2__k3 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_UpsampleArgbLinePair_SSE2__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_UpsampleArgbLinePair_SSE2__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_UpsampleArgbLinePair_SSE2__k6 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_UpsampleArgbLinePair_SSE2__k7 = [2]uint64{0x0, 0x0}
var F_UpsampleArgbLinePair_SSE2__k8 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_UpsampleArgbLinePair_SSE2__k9 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_UpsampleArgbLinePair_SSE2__k10 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_UpsampleArgbLinePair_SSE2__k11 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_UpsampleArgbLinePair_SSE2__k12 = [2]uint64{0x191300001913, 0x191300001913}
var F_UpsampleArgbLinePair_SSE2__k13 = [2]uint64{0x340800003408, 0x340800003408}
var F_UpsampleArgbLinePair_SSE2__k14 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_UpsampleArgbLinePair_SSE2__k15 = [2]uint64{0x2204220422042204, 0x2204220422042204}
var F_UpsampleArgbLinePair_SSE2__k16 = [2]uint64{0x662500006625, 0x662500006625}
var F_UpsampleArgbLinePair_SSE2__k17 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_UpsampleArgbLinePair_SSE2__k18 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_UpsampleArgbLinePair_SSE2__k19 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_UpsampleArgbLinePair_SSE2__k20 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_UpsampleArgbLinePair_SSE2__k21 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_UpsampleArgbLinePair_SSE2__k22 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_UpsampleArgbLinePair_SSE2__k23 = [2]uint64{0x302808001008080, 0x706808005048080}
