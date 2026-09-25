//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_WebPConfigInitInternal(m *base.Module, l0 int32, l1 int32, l2 float32, l3 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v13 int64
	_ = v13
	var v18 base.V128
	_ = v18
	var v21 base.V128
	_ = v21
	var v24 base.V128
	_ = v24
	var v27 int32
	_ = v27
	var v35 base.V128
	_ = v35
	var v62 int32
	_ = v62
	var v69 float32
	_ = v69
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v78 float32
	_ = v78
	var v81 int32
	_ = v81
	var v84 int32
	_ = v84
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
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v117 int32
	_ = v117
	var v120 int32
	_ = v120
	var v123 int32
	_ = v123
	var v126 int32
	_ = v126
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v133 int32
	_ = v133
	var v134 int32
	_ = v134
	var v137 int32
	_ = v137
	var v140 int32
	_ = v140
	var v143 int32
	_ = v143
	var v146 int32
	_ = v146
	var v149 int32
	_ = v149
	var v152 int32
	_ = v152
	var v155 int32
	_ = v155
	var v158 int32
	_ = v158
	var v161 int32
	_ = v161
	var v164 int32
	_ = v164
	var v168 int32
	_ = v168
	v5 = int32(0)
	if l0 == v5 {
		v168 = v5
	} else {
		if l3&int32(-256) != int32(512) {
			v168 = v5
		} else {
			v13 = int64(0)
			*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = v13
			*(*float32)(unsafe.Add(mBase, uint32(l0)+4)) = l2
			*(*int32)(unsafe.Add(mBase, uint32(l0)+112)) = int32(100)
			v18 = base.Simd_g_const(&F_WebPConfigInitInternal__k0)
			base.Simd_g_v128_store(m, l0, int32(24), v18)
			v21 = base.Simd_g_const(&F_WebPConfigInitInternal__k1)
			base.Simd_g_v128_store(m, l0, int32(56), v21)
			v24 = base.Simd_g_const(&F_WebPConfigInitInternal__k2)
			base.Simd_g_v128_store(m, l0, int32(40), v24)
			v27 = int32(0)
			*(*int32)(unsafe.Add(mBase, uint32(l0)+96)) = v27
			*(*int32)(unsafe.Add(mBase, uint32(l0))) = v27
			*(*int64)(unsafe.Add(mBase, uint32(l0)+104)) = v13
			*(*int64)(unsafe.Add(mBase, uint32(l0)+88)) = int64(429496729600)
			v35 = base.Simd_g_const(&F_WebPConfigInitInternal__k3)
			base.Simd_g_v128_store(m, l0, int32(72), v35)
			*(*int64)(unsafe.Add(mBase, uint32(l0)+8)) = int64(4)
			switch l1 + int32(-1) {
			case 0:
				*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = int32(4)
				*(*int64)(unsafe.Add(mBase, uint32(l0)+28)) = int64(150323855440)
			case 1:
				*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = int32(3)
				*(*int64)(unsafe.Add(mBase, uint32(l0)+28)) = int64(128849018960)
				*(*int32)(unsafe.Add(mBase, uint32(l0)+68)) = int32(2)
			case 2:
				*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = int32(6)
				*(*int64)(unsafe.Add(mBase, uint32(l0)+28)) = int64(42949672985)
			case 3:
				*(*int64)(unsafe.Add(mBase, uint32(l0)+28)) = int64(0)
			case 4:
				*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = int32(0)
				*(*int64)(unsafe.Add(mBase, uint32(l0)+24)) = int64(2)
			default:
			}
			v62 = int32(0)
			if l0 == v62 {
				v164 = v62
			} else {
				v69 = *(*float32)(unsafe.Add(mBase, uint32(l0)+4))
				if base.F32_lt(v69, float32(0)) != 0 {
					v164 = v62
				} else {
					if base.F32_gt(v69, float32(100)) != 0 {
						v164 = v62
					} else {
						v74 = int32(0)
						v75 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
						if v75 < v74 {
							v164 = v74
						} else {
							v78 = *(*float32)(unsafe.Add(mBase, uint32(l0)+20))
							if base.F32_lt(v78, float32(0)) != 0 {
								v164 = v74
							} else {
								v81 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
								if base.Ui32(int32(6)) < base.Ui32(v81) {
									v164 = v74
								} else {
									v84 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
									if base.Ui32(v84+int32(-5)) < base.Ui32(int32(-4)) {
										v164 = v74
									} else {
										v89 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
										if base.Ui32(int32(100)) < base.Ui32(v89) {
											v164 = v74
										} else {
											v92 = *(*int32)(unsafe.Add(mBase, uint32(l0)+32))
											if base.Ui32(int32(100)) < base.Ui32(v92) {
												v164 = v74
											} else {
												v95 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
												if base.Ui32(int32(7)) < base.Ui32(v95) {
													v164 = v74
												} else {
													v98 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
													if base.Ui32(int32(1)) < base.Ui32(v98) {
														v164 = v74
													} else {
														v101 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
														if base.Ui32(int32(1)) < base.Ui32(v101) {
															v164 = v74
														} else {
															v104 = *(*int32)(unsafe.Add(mBase, uint32(l0)+60))
															if base.Ui32(v104+int32(-11)) < base.Ui32(int32(-10)) {
																v164 = v74
															} else {
																v109 = int32(0)
																v110 = *(*int32)(unsafe.Add(mBase, uint32(l0)+108))
																if v110 < v109 {
																	v164 = v109
																} else {
																	v113 = *(*int32)(unsafe.Add(mBase, uint32(l0)+112))
																	if int32(100) < v113 {
																		v164 = v109
																	} else {
																		if v113 < v110 {
																			v164 = v109
																		} else {
																			v117 = *(*int32)(unsafe.Add(mBase, uint32(l0)+64))
																			if base.Ui32(int32(1)) < base.Ui32(v117) {
																				v164 = v109
																			} else {
																				v120 = *(*int32)(unsafe.Add(mBase, uint32(l0)+68))
																				if base.Ui32(int32(7)) < base.Ui32(v120) {
																					v164 = v109
																				} else {
																					v123 = *(*int32)(unsafe.Add(mBase, uint32(l0)+72))
																					if base.Ui32(int32(3)) < base.Ui32(v123) {
																						v164 = v109
																					} else {
																						v126 = *(*int32)(unsafe.Add(mBase, uint32(l0)+76))
																						if base.Ui32(int32(100)) < base.Ui32(v126) {
																							v164 = v109
																						} else {
																							v129 = int32(0)
																							v130 = *(*int32)(unsafe.Add(mBase, uint32(l0)+48))
																							if v130 < v129 {
																								v164 = v129
																							} else {
																								v133 = int32(0)
																								v134 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
																								if v134 < v133 {
																									v164 = v133
																								} else {
																									v137 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
																									if base.Ui32(int32(100)) < base.Ui32(v137) {
																										v164 = v133
																									} else {
																										v140 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
																										if base.Ui32(int32(1)) < base.Ui32(v140) {
																											v164 = v133
																										} else {
																											v143 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
																											if base.Ui32(int32(100)) < base.Ui32(v143) {
																												v164 = v133
																											} else {
																												v146 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
																												if base.Ui32(int32(3)) < base.Ui32(v146) {
																													v164 = v133
																												} else {
																													v149 = *(*int32)(unsafe.Add(mBase, uint32(l0)+80))
																													if base.Ui32(int32(1)) < base.Ui32(v149) {
																														v164 = v133
																													} else {
																														v152 = *(*int32)(unsafe.Add(mBase, uint32(l0)+84))
																														if base.Ui32(int32(1)) < base.Ui32(v152) {
																															v164 = v133
																														} else {
																															v155 = *(*int32)(unsafe.Add(mBase, uint32(l0)+88))
																															if base.Ui32(int32(1)) < base.Ui32(v155) {
																																v164 = v133
																															} else {
																																v158 = *(*int32)(unsafe.Add(mBase, uint32(l0)+96))
																																if base.Ui32(int32(1)) < base.Ui32(v158) {
																																	v164 = v133
																																} else {
																																	v161 = *(*int32)(unsafe.Add(mBase, uint32(l0)+104))
																																	v164 = base.B2i32(base.Ui32(v161) < base.Ui32(int32(2)))
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
			v168 = v164
		}
	}
	return v168
}

var F_WebPConfigInitInternal__k0 = [2]uint64{0x3200000004, 0x3c}
var F_WebPConfigInitInternal__k1 = [2]uint64{0x100000064, 0x0}
var F_WebPConfigInitInternal__k2 = [2]uint64{0x1, 0x100000001}
var F_WebPConfigInitInternal__k3 = [2]uint64{0x0, 0x0}

func F_WebPEstimateBestFilter(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) int32 {
	mBase := m.M
	_ = mBase
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v174 int32
	_ = v174
	var v197 int32
	_ = v197
	var v220 int32
	_ = v220
	var v230 int32
	_ = v230
	var v231 int32
	_ = v231
	var v233 int32
	_ = v233
	var v250 int32
	_ = v250
	var v251 int32
	_ = v251
	var v256 int32
	_ = v256
	var v287 int32
	_ = v287
	var v293 int32
	_ = v293
	var v296 int32
	_ = v296
	var v336 int32
	_ = v336
	var v337 int32
	_ = v337
	var v339 int32
	_ = v339
	var v340 int32
	_ = v340
	var v341 int32
	_ = v341
	var v342 int32
	_ = v342
	var v345 int32
	_ = v345
	var v352 int32
	_ = v352
	var v353 int32
	_ = v353
	var v355 int32
	_ = v355
	var v360 int32
	_ = v360
	var v365 int32
	_ = v365
	var v368 int32
	_ = v368
	var v369 int32
	_ = v369
	var v371 int32
	_ = v371
	var v382 int32
	_ = v382
	var v383 int32
	_ = v383
	var v384 int32
	_ = v384
	var v387 int32
	_ = v387
	var v388 int32
	_ = v388
	var v391 int32
	_ = v391
	var v392 int32
	_ = v392
	var v394 int32
	_ = v394
	var v397 int32
	_ = v397
	var v419 int32
	_ = v419
	var v422 int32
	_ = v422
	var v423 int32
	_ = v423
	var v426 int32
	_ = v426
	var v428 int32
	_ = v428
	var v429 int32
	_ = v429
	var v432 int32
	_ = v432
	var v434 int32
	_ = v434
	var v435 int32
	_ = v435
	var v438 int32
	_ = v438
	var v440 int32
	_ = v440
	var v441 int32
	_ = v441
	var v444 int32
	_ = v444
	var v446 int32
	_ = v446
	var v447 int32
	_ = v447
	var v450 int32
	_ = v450
	var v452 int32
	_ = v452
	var v453 int32
	_ = v453
	var v456 int32
	_ = v456
	var v458 int32
	_ = v458
	var v459 int32
	_ = v459
	var v462 int32
	_ = v462
	var v464 int32
	_ = v464
	var v465 int32
	_ = v465
	var v468 int32
	_ = v468
	var v470 int32
	_ = v470
	var v471 int32
	_ = v471
	var v474 int32
	_ = v474
	var v476 int32
	_ = v476
	var v477 int32
	_ = v477
	var v480 int32
	_ = v480
	var v482 int32
	_ = v482
	var v483 int32
	_ = v483
	var v486 int32
	_ = v486
	var v488 int32
	_ = v488
	var v489 int32
	_ = v489
	var v492 int32
	_ = v492
	var v494 int32
	_ = v494
	var v495 int32
	_ = v495
	var v498 int32
	_ = v498
	var v500 int32
	_ = v500
	var v501 int32
	_ = v501
	var v504 int32
	_ = v504
	var v505 int32
	_ = v505
	var v506 int32
	_ = v506
	var v508 int32
	_ = v508
	var v511 int32
	_ = v511
	var v514 int32
	_ = v514
	var v516 int32
	_ = v516
	var v519 int32
	_ = v519
	var v521 int32
	_ = v521
	var v526 int32
	_ = v526
	var v531 int32
	_ = v531
	var v540 int32
	_ = v540
	var v542 int32
	_ = v542
	var v543 int32
	_ = v543
	var v544 int32
	_ = v544
	var v545 int32
	_ = v545
	var v546 int32
	_ = v546
	var v547 int32
	_ = v547
	var v548 int32
	_ = v548
	var v549 int32
	_ = v549
	var v550 int32
	_ = v550
	var v551 int32
	_ = v551
	var v552 int32
	_ = v552
	var v553 int32
	_ = v553
	var v554 int32
	_ = v554
	var v555 int32
	_ = v555
	var v556 int32
	_ = v556
	var v557 int32
	_ = v557
	var v558 int32
	_ = v558
	var v559 int32
	_ = v559
	var v560 int32
	_ = v560
	var v561 int32
	_ = v561
	var v583 int32
	_ = v583
	var v584 int32
	_ = v584
	var v585 int32
	_ = v585
	var v586 int32
	_ = v586
	var v587 int32
	_ = v587
	var v588 int32
	_ = v588
	var v589 int32
	_ = v589
	var v590 int32
	_ = v590
	var v591 int32
	_ = v591
	var v595 base.V128
	_ = v595
	var v596 int32
	_ = v596
	var v597 int32
	_ = v597
	var v598 base.V128
	_ = v598
	var v602 base.V128
	_ = v602
	var v605 base.V128
	_ = v605
	var v609 base.V128
	_ = v609
	var v612 base.V128
	_ = v612
	var v616 base.V128
	_ = v616
	var v619 base.V128
	_ = v619
	var v623 base.V128
	_ = v623
	var v626 base.V128
	_ = v626
	var v630 base.V128
	_ = v630
	var v633 base.V128
	_ = v633
	var v637 base.V128
	_ = v637
	var v640 base.V128
	_ = v640
	var v644 base.V128
	_ = v644
	var v647 base.V128
	_ = v647
	var v651 base.V128
	_ = v651
	var v654 base.V128
	_ = v654
	var v658 base.V128
	_ = v658
	var v661 base.V128
	_ = v661
	var v665 base.V128
	_ = v665
	var v668 base.V128
	_ = v668
	var v672 base.V128
	_ = v672
	var v675 base.V128
	_ = v675
	var v679 base.V128
	_ = v679
	var v682 base.V128
	_ = v682
	var v686 base.V128
	_ = v686
	var v689 base.V128
	_ = v689
	var v693 base.V128
	_ = v693
	var v696 base.V128
	_ = v696
	var v706 int32
	_ = v706
	var v708 int32
	_ = v708
	var v711 int32
	_ = v711
	var v713 int32
	_ = v713
	var v716 int32
	_ = v716
	var v718 int32
	_ = v718
	var v721 int32
	_ = v721
	var v723 int32
	_ = v723
	var v726 int32
	_ = v726
	var v728 int32
	_ = v728
	var v731 int32
	_ = v731
	var v733 int32
	_ = v733
	var v736 int32
	_ = v736
	var v737 int32
	_ = v737
	var v745 int32
	_ = v745
	var v759 int32
	_ = v759
	var v773 int32
	_ = v773
	var v774 int32
	_ = v774
	var v776 base.V128
	_ = v776
	var v832 base.V128
	_ = v832
	var v834 int32
	_ = v834
	var v835 int32
	_ = v835
	var v837 int32
	_ = v837
	var v839 int32
	_ = v839
	var v841 base.V128
	_ = v841
	var v847 int32
	_ = v847
	var v852 int32
	_ = v852
	v48 = m.G0
	v49 = int32(256)
	v50 = v48 - v49
	m.G0 = v50
	base.MemoryFill(m, v50, int32(0), v49)
	v174 = int32(0)
	if l2 < int32(4) {
		v540 = v174
		v542 = v174
		v543 = v174
		v544 = v174
		v545 = v174
		v546 = v174
		v547 = v174
		v548 = v174
		v549 = v174
		v550 = v174
		v551 = v174
		v552 = v174
		v553 = v174
		v554 = v174
		v555 = v174
		v556 = v174
		v557 = v174
		v558 = v174
		v559 = v174
		v560 = v174
		v561 = v174
	} else {
		v197 = int32(0)
		if l1 < int32(4) {
			v540 = v197
			v542 = v197
			v543 = v197
			v544 = v197
			v545 = v197
			v546 = v197
			v547 = v197
			v548 = v197
			v549 = v197
			v550 = v197
			v551 = v197
			v552 = v197
			v553 = v197
			v554 = v197
			v555 = v197
			v556 = v197
			v557 = v197
			v558 = v197
			v559 = v197
			v560 = v197
			v561 = v197
		} else {
			v220 = int32(-1)
			v230 = int32(1)
			v231 = l3 << (uint(v230) % 32)
			v233 = v231 | v230
			v250 = l0 + v233
			v251 = l0 + (v233 - l1)
			v256 = int32(2)
			for {
				v287 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v256*l3))))
				v293 = int32(0)
				v296 = v287
				for {
					v336 = v250 + v293
					v337 = int32(1)
					v339 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v336+v337))))
					v340 = v339 - v296
					v341 = int32(31)
					v342 = v340 >> (uint(v341) % 32)
					v345 = int32(2)
					*(*int32)(unsafe.Add(mBase, uint32(v50+int32(base.Ui32(v340^v342-v342)>>(uint(v345)%32))&int32(1073741820)))) = v337
					v352 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v336))))
					v353 = v339 - v352
					v355 = v353 >> (uint(v341) % 32)
					v360 = int32(252)
					*(*int32)(unsafe.Add(mBase, uint32(v50+int32(64)+int32(base.Ui32(v353^v355-v355)>>(uint(v345)%32))&v360))) = v337
					v365 = v251 + v293
					v368 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v365+v337))))
					v369 = v339 - v368
					v371 = v369 >> (uint(v341) % 32)
					*(*int32)(unsafe.Add(mBase, uint32(v50+int32(128)+int32(base.Ui32(v369^v371-v371)>>(uint(v345)%32))&v360))) = v337
					v382 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v365))))
					v383 = v368 + v352 - v382
					v384 = int32(0)
					if v384 < v383 {
						v387 = v383
					} else {
						v387 = v384
					}
					v388 = int32(255)
					if v387 < v388 {
						v391 = v387
					} else {
						v391 = v388
					}
					v392 = v339 - v391
					v394 = v392 >> (uint(int32(31)) % 32)
					v397 = int32(2)
					*(*int32)(unsafe.Add(mBase, uint32(v50+int32(192)+int32(base.Ui32(v392^v394-v394)>>(uint(v397)%32))&int32(252)))) = int32(1)
					if v293+int32(4) < l1+v220 {
						v293 = v293 + v397
						v296 = int32(base.Ui32(v339+v296*int32(3)+v397) >> (uint(v397) % 32))
						continue
					} else {
						break
					}
					break
				}
				v419 = v256 + int32(2)
				if v419 < l2+v220 {
					v250 = v250 + v231
					v251 = v251 + v231
					v256 = v419
					continue
				} else {
					break
				}
				break
			}
			v422 = int32(0)
			v423 = *(*int32)(unsafe.Add(mBase, uint32(v50)+88))
			if v422 < v423 {
				v426 = int32(6)
			} else {
				v426 = v422
			}
			v428 = int32(0)
			v429 = *(*int32)(unsafe.Add(mBase, uint32(v50)+84))
			if v428 < v429 {
				v432 = int32(5)
			} else {
				v432 = v428
			}
			v434 = int32(0)
			v435 = *(*int32)(unsafe.Add(mBase, uint32(v50)+76))
			if v434 < v435 {
				v438 = int32(3)
			} else {
				v438 = v434
			}
			v440 = int32(0)
			v441 = *(*int32)(unsafe.Add(mBase, uint32(v50)+60))
			if v440 < v441 {
				v444 = int32(15)
			} else {
				v444 = v440
			}
			v446 = int32(0)
			v447 = *(*int32)(unsafe.Add(mBase, uint32(v50)+56))
			if v446 < v447 {
				v450 = int32(14)
			} else {
				v450 = v446
			}
			v452 = int32(0)
			v453 = *(*int32)(unsafe.Add(mBase, uint32(v50)+52))
			if v452 < v453 {
				v456 = int32(13)
			} else {
				v456 = v452
			}
			v458 = int32(0)
			v459 = *(*int32)(unsafe.Add(mBase, uint32(v50)+48))
			if v458 < v459 {
				v462 = int32(12)
			} else {
				v462 = v458
			}
			v464 = int32(0)
			v465 = *(*int32)(unsafe.Add(mBase, uint32(v50)+44))
			if v464 < v465 {
				v468 = int32(11)
			} else {
				v468 = v464
			}
			v470 = int32(0)
			v471 = *(*int32)(unsafe.Add(mBase, uint32(v50)+40))
			if v470 < v471 {
				v474 = int32(10)
			} else {
				v474 = v470
			}
			v476 = int32(0)
			v477 = *(*int32)(unsafe.Add(mBase, uint32(v50)+36))
			if v476 < v477 {
				v480 = int32(9)
			} else {
				v480 = v476
			}
			v482 = int32(0)
			v483 = *(*int32)(unsafe.Add(mBase, uint32(v50)+28))
			if v482 < v483 {
				v486 = int32(7)
			} else {
				v486 = v482
			}
			v488 = int32(0)
			v489 = *(*int32)(unsafe.Add(mBase, uint32(v50)+24))
			if v488 < v489 {
				v492 = int32(6)
			} else {
				v492 = v488
			}
			v494 = int32(0)
			v495 = *(*int32)(unsafe.Add(mBase, uint32(v50)+20))
			if v494 < v495 {
				v498 = int32(5)
			} else {
				v498 = v494
			}
			v500 = int32(0)
			v501 = *(*int32)(unsafe.Add(mBase, uint32(v50)+12))
			if v500 < v501 {
				v504 = int32(3)
			} else {
				v504 = v500
			}
			v505 = *(*int32)(unsafe.Add(mBase, uint32(v50)+68))
			v506 = int32(0)
			v508 = *(*int32)(unsafe.Add(mBase, uint32(v50)+4))
			v511 = *(*int32)(unsafe.Add(mBase, uint32(v50)+80))
			v514 = int32(2)
			v516 = *(*int32)(unsafe.Add(mBase, uint32(v50)+72))
			v519 = int32(1)
			v521 = *(*int32)(unsafe.Add(mBase, uint32(v50)+32))
			v526 = *(*int32)(unsafe.Add(mBase, uint32(v50)+16))
			v531 = *(*int32)(unsafe.Add(mBase, uint32(v50)+8))
			v540 = v426
			v542 = v432
			v543 = base.B2i32(v506 < v511) << (uint(v514) % 32)
			v544 = v438
			v545 = base.B2i32(v506 < v516) << (uint(v519) % 32)
			v546 = base.B2i32(v506 < v505)
			v547 = v444
			v548 = v450
			v549 = v456
			v550 = v462
			v551 = v468
			v552 = v474
			v553 = v480
			v554 = base.B2i32(v506 < v521) << (uint(int32(3)) % 32)
			v555 = v486
			v556 = v492
			v557 = v498
			v558 = base.B2i32(v506 < v526) << (uint(v514) % 32)
			v559 = v504
			v560 = base.B2i32(v506 < v531) << (uint(v519) % 32)
			v561 = base.B2i32(v506 < v508)
		}
	}
	v583 = *(*int32)(unsafe.Add(mBase, uint32(v50)+92))
	v584 = *(*int32)(unsafe.Add(mBase, uint32(v50)+96))
	v585 = *(*int32)(unsafe.Add(mBase, uint32(v50)+100))
	v586 = *(*int32)(unsafe.Add(mBase, uint32(v50)+104))
	v587 = *(*int32)(unsafe.Add(mBase, uint32(v50)+108))
	v588 = *(*int32)(unsafe.Add(mBase, uint32(v50)+112))
	v589 = *(*int32)(unsafe.Add(mBase, uint32(v50)+116))
	v590 = *(*int32)(unsafe.Add(mBase, uint32(v50)+120))
	v591 = *(*int32)(unsafe.Add(mBase, uint32(v50)+124))
	v595 = base.Simd_g_v128_load(m, v50, int32(200))
	v596 = int32(0)
	v597 = int32(1)
	v598 = base.Simd_g_v128_load32_lane_l1(m, v50+int32(136), v596, v595)
	v602 = base.Simd_g_v128_load(m, v50, int32(196))
	v605 = base.Simd_g_v128_load32_lane_l1(m, v50+int32(132), v596, v602)
	v609 = base.Simd_g_v128_load(m, v50, int32(204))
	v612 = base.Simd_g_v128_load32_lane_l1(m, v50+int32(140), v596, v609)
	v616 = base.Simd_g_v128_load(m, v50, int32(208))
	v619 = base.Simd_g_v128_load32_lane_l1(m, v50+int32(144), v596, v616)
	v623 = base.Simd_g_v128_load(m, v50, int32(212))
	v626 = base.Simd_g_v128_load32_lane_l1(m, v50+int32(148), v596, v623)
	v630 = base.Simd_g_v128_load(m, v50, int32(216))
	v633 = base.Simd_g_v128_load32_lane_l1(m, v50+int32(152), v596, v630)
	v637 = base.Simd_g_v128_load(m, v50, int32(220))
	v640 = base.Simd_g_v128_load32_lane_l1(m, v50+int32(156), v596, v637)
	v644 = base.Simd_g_v128_load(m, v50, int32(224))
	v647 = base.Simd_g_v128_load32_lane_l1(m, v50+int32(160), v596, v644)
	v651 = base.Simd_g_v128_load(m, v50, int32(228))
	v654 = base.Simd_g_v128_load32_lane_l1(m, v50+int32(164), v596, v651)
	v658 = base.Simd_g_v128_load(m, v50, int32(232))
	v661 = base.Simd_g_v128_load32_lane_l1(m, v50+int32(168), v596, v658)
	v665 = base.Simd_g_v128_load(m, v50, int32(236))
	v668 = base.Simd_g_v128_load32_lane_l1(m, v50+int32(172), v596, v665)
	v672 = base.Simd_g_v128_load(m, v50, int32(240))
	v675 = base.Simd_g_v128_load32_lane_l1(m, v50+int32(176), v596, v672)
	v679 = base.Simd_g_v128_load32_splat(m, v50, int32(244))
	v682 = base.Simd_g_v128_load32_lane_l1(m, v50+int32(180), v596, v679)
	v686 = base.Simd_g_v128_load32_splat(m, v50, int32(248))
	v689 = base.Simd_g_v128_load32_lane_l1(m, v50+int32(184), v596, v686)
	v693 = base.Simd_g_v128_load32_splat(m, v50, int32(252))
	v696 = base.Simd_g_v128_load32_lane_l1(m, v50+int32(188), v596, v693)
	m.G0 = v50 + int32(256)
	if v596 < v591 {
		v706 = int32(15)
	} else {
		v706 = v596
	}
	v708 = int32(0)
	if v708 < v590 {
		v711 = int32(14)
	} else {
		v711 = v708
	}
	v713 = int32(0)
	if v713 < v589 {
		v716 = int32(13)
	} else {
		v716 = v713
	}
	v718 = int32(0)
	if v718 < v588 {
		v721 = int32(12)
	} else {
		v721 = v718
	}
	v723 = int32(0)
	if v723 < v587 {
		v726 = int32(11)
	} else {
		v726 = v723
	}
	v728 = int32(0)
	if v728 < v586 {
		v731 = int32(10)
	} else {
		v731 = v728
	}
	v733 = int32(0)
	if v733 < v585 {
		v736 = int32(9)
	} else {
		v736 = v733
	}
	v737 = int32(0)
	if v737 < v583 {
		v745 = int32(7)
	} else {
		v745 = v737
	}
	v759 = v706 + (v711 + (v716 + (v721 + (v726 + (v731 + (v736 + (base.B2i32(v737 < v584)<<(uint(int32(3))%32) + (v745 + (v540 + (v542 + (v543 + (v544 + (v545 | v546)))))))))))))
	v773 = v547 + (v548 + (v549 + (v550 + (v551 + (v552 + (v553 + (v554 + (v555 + (v556 + (v557 + (v558 + (v559 + (v560 | v561)))))))))))))
	v774 = base.B2i32(base.Ui32(v759) < base.Ui32(v773))
	v776 = base.Simd_g_const(&F_WebPEstimateBestFilter__k0)
	v832 = base.Simd_g_i32x4_add(base.Simd_g_v128_bitselect(base.Simd_g_const(&F_WebPEstimateBestFilter__k1), v776, base.Simd_g_i32x4_gt_s(v696, v776)), base.Simd_g_i32x4_add(base.Simd_g_v128_bitselect(base.Simd_g_const(&F_WebPEstimateBestFilter__k2), v776, base.Simd_g_i32x4_gt_s(v689, v776)), base.Simd_g_i32x4_add(base.Simd_g_v128_bitselect(base.Simd_g_const(&F_WebPEstimateBestFilter__k3), v776, base.Simd_g_i32x4_gt_s(v682, v776)), base.Simd_g_i32x4_add(base.Simd_g_v128_bitselect(base.Simd_g_const(&F_WebPEstimateBestFilter__k4), v776, base.Simd_g_i32x4_gt_s(v675, v776)), base.Simd_g_i32x4_add(base.Simd_g_v128_bitselect(base.Simd_g_const(&F_WebPEstimateBestFilter__k5), v776, base.Simd_g_i32x4_gt_s(v668, v776)), base.Simd_g_i32x4_add(base.Simd_g_v128_bitselect(base.Simd_g_const(&F_WebPEstimateBestFilter__k6), v776, base.Simd_g_i32x4_gt_s(v661, v776)), base.Simd_g_i32x4_add(base.Simd_g_v128_bitselect(base.Simd_g_const(&F_WebPEstimateBestFilter__k7), v776, base.Simd_g_i32x4_gt_s(v654, v776)), base.Simd_g_i32x4_add(base.Simd_g_v128_bitselect(base.Simd_g_const(&F_WebPEstimateBestFilter__k8), v776, base.Simd_g_i32x4_gt_s(v647, v776)), base.Simd_g_i32x4_add(base.Simd_g_v128_bitselect(base.Simd_g_const(&F_WebPEstimateBestFilter__k9), v776, base.Simd_g_i32x4_gt_s(v640, v776)), base.Simd_g_i32x4_add(base.Simd_g_v128_bitselect(base.Simd_g_const(&F_WebPEstimateBestFilter__k10), v776, base.Simd_g_i32x4_gt_s(v633, v776)), base.Simd_g_i32x4_add(base.Simd_g_v128_bitselect(base.Simd_g_const(&F_WebPEstimateBestFilter__k11), v776, base.Simd_g_i32x4_gt_s(v626, v776)), base.Simd_g_i32x4_add(base.Simd_g_v128_bitselect(base.Simd_g_const(&F_WebPEstimateBestFilter__k12), v776, base.Simd_g_i32x4_gt_s(v619, v776)), base.Simd_g_i32x4_add(base.Simd_g_v128_bitselect(base.Simd_g_const(&F_WebPEstimateBestFilter__k13), v776, base.Simd_g_i32x4_gt_s(v612, v776)), base.Simd_g_i32x4_sub(base.Simd_g_v128_bitselect(base.Simd_g_const(&F_WebPEstimateBestFilter__k14), v776, base.Simd_g_i32x4_gt_s(v598, v776)), base.Simd_g_i32x4_gt_s(v605, v776)))))))))))))))
	v834 = base.Simd_g_i32x4_extract_lane_l1(v832)
	if base.Ui32(v759) < base.Ui32(v773) {
		v835 = v759
	} else {
		v835 = v773
	}
	if base.Ui32(v834) < base.Ui32(v835) {
		v837 = v834
	} else {
		v837 = v835
	}
	v839 = int32(1)
	v841 = base.Simd_g_i32x4_lt_u(v832, base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(v837), v835))
	if base.Simd_g_i32x4_extract_lane_l2(base.Simd_g_i64x2_extend_low_i32x4_s(v841))&v839 != 0 {
		v847 = int32(2)
	} else {
		v847 = v774
	}
	if base.Simd_g_i32x4_extract_lane_l0(v841)&int32(1) != 0 {
		v852 = int32(3)
	} else {
		v852 = v847
	}
	return v852
}

var F_WebPEstimateBestFilter__k0 = [2]uint64{0x0, 0x0}
var F_WebPEstimateBestFilter__k1 = [2]uint64{0xf0000000f, 0xf0000000f}
var F_WebPEstimateBestFilter__k2 = [2]uint64{0xe0000000e, 0xe0000000e}
var F_WebPEstimateBestFilter__k3 = [2]uint64{0xd0000000d, 0xd0000000d}
var F_WebPEstimateBestFilter__k4 = [2]uint64{0xc0000000c, 0xc0000000c}
var F_WebPEstimateBestFilter__k5 = [2]uint64{0xb0000000b, 0xb0000000b}
var F_WebPEstimateBestFilter__k6 = [2]uint64{0xa0000000a, 0xa0000000a}
var F_WebPEstimateBestFilter__k7 = [2]uint64{0x900000009, 0x900000009}
var F_WebPEstimateBestFilter__k8 = [2]uint64{0x800000008, 0x800000008}
var F_WebPEstimateBestFilter__k9 = [2]uint64{0x700000007, 0x700000007}
var F_WebPEstimateBestFilter__k10 = [2]uint64{0x600000006, 0x600000006}
var F_WebPEstimateBestFilter__k11 = [2]uint64{0x500000005, 0x500000005}
var F_WebPEstimateBestFilter__k12 = [2]uint64{0x400000004, 0x400000004}
var F_WebPEstimateBestFilter__k13 = [2]uint64{0x300000003, 0x300000003}
var F_WebPEstimateBestFilter__k14 = [2]uint64{0x200000002, 0x200000002}

func F_WebPPictureAlloc(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v6 int32
	_ = v6
	var v8 int64
	_ = v8
	var v16 base.V128
	_ = v16
	var v17 int32
	_ = v17
	var v23 int32
	_ = v23
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v38 int32
	_ = v38
	var v41 int32
	_ = v41
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v50 int32
	_ = v50
	var v54 int32
	_ = v54
	var v58 base.V128
	_ = v58
	var v65 int64
	_ = v65
	var v66 int64
	_ = v66
	var v69 int64
	_ = v69
	var v70 int32
	_ = v70
	var v74 int64
	_ = v74
	var v75 int64
	_ = v75
	var v76 int64
	_ = v76
	var v85 int32
	_ = v85
	var v90 int32
	_ = v90
	var v92 int64
	_ = v92
	var v93 int64
	_ = v93
	var v95 int64
	_ = v95
	var v98 int64
	_ = v98
	var v103 int32
	_ = v103
	var v105 int32
	_ = v105
	var v112 int32
	_ = v112
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v131 int32
	_ = v131
	var v136 int32
	_ = v136
	var v137 int32
	_ = v137
	var v138 int32
	_ = v138
	var v144 int32
	_ = v144
	var v145 int32
	_ = v145
	var v146 int32
	_ = v146
	var v149 int32
	_ = v149
	var v161 int32
	_ = v161
	var v163 int32
	_ = v163
	var v174 int32
	_ = v174
	if l0 != 0 {
		v4 = *(*int32)(unsafe.Add(mBase, uint32(l0)+156))
		F_free(m, v4)
		mBase = m.M
		v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
		F_free(m, v6)
		mBase = m.M
		v8 = int64(0)
		*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = v8
		*(*int64)(unsafe.Add(mBase, uint32(l0)+156)) = v8
		*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = v8
		v16 = base.Simd_g_const(&F_WebPPictureAlloc__k0)
		v17 = int32(0)
		base.Simd_g_v128_store(m, l0+int32(24), v17, v16)
		*(*int32)(unsafe.Add(mBase, uint32(l0+int32(40)))) = v17
		v23 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		if v23 != 0 {
			v136 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
			v137 = int32(5)
			v138 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
			if v138 < int32(1) {
				v145 = v137
				v146 = F_WebPEncodingSetError(m, l0, v145)
				mBase = m.M
				if v146 != 0 {
					v149 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
					F_WebPSafeFree(m, v149)
					mBase = m.M
					*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
					v161 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v136)*base.I64_extend_i32_s(v138)+int64(31), int32(4))
					mBase = m.M
					if v161 != 0 {
						*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v138
						*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v161
						*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v161 + int32(31)) & int32(-32)
						v174 = int32(1)
					} else {
						v163 = F_WebPEncodingSetError(m, l0, int32(1))
						mBase = m.M
						v174 = v163
					}
				} else {
					v174 = int32(0)
				}
			} else {
				if v136 < int32(1) {
					v145 = v137
					v146 = F_WebPEncodingSetError(m, l0, v145)
					mBase = m.M
					if v146 != 0 {
						v149 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
						F_WebPSafeFree(m, v149)
						mBase = m.M
						*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
						*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
						v161 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v136)*base.I64_extend_i32_s(v138)+int64(31), int32(4))
						mBase = m.M
						if v161 != 0 {
							*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v138
							*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v161
							*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v161 + int32(31)) & int32(-32)
							v174 = int32(1)
						} else {
							v163 = F_WebPEncodingSetError(m, l0, int32(1))
							mBase = m.M
							v174 = v163
						}
					} else {
						v174 = int32(0)
					}
				} else {
					v144 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					switch v144 {
					case 0, 4:
						v149 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
						F_WebPSafeFree(m, v149)
						mBase = m.M
						*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
						*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
						v161 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v136)*base.I64_extend_i32_s(v138)+int64(31), int32(4))
						mBase = m.M
						if v161 != 0 {
							*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v138
							*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v161
							*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v161 + int32(31)) & int32(-32)
							v174 = int32(1)
						} else {
							v163 = F_WebPEncodingSetError(m, l0, int32(1))
							mBase = m.M
							v174 = v163
						}
					default:
						v145 = int32(4)
						v146 = F_WebPEncodingSetError(m, l0, v145)
						mBase = m.M
						if v146 != 0 {
							v149 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
							F_WebPSafeFree(m, v149)
							mBase = m.M
							*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
							v161 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v136)*base.I64_extend_i32_s(v138)+int64(31), int32(4))
							mBase = m.M
							if v161 != 0 {
								*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v138
								*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v161
								*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v161 + int32(31)) & int32(-32)
								v174 = int32(1)
							} else {
								v163 = F_WebPEncodingSetError(m, l0, int32(1))
								mBase = m.M
								v174 = v163
							}
						} else {
							v174 = int32(0)
						}
					}
				}
			}
			return v174
		} else {
			v34 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
			v35 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
			v36 = int32(1)
			v38 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
			v41 = base.B2i32(v35 < v36) | base.B2i32(v38 < v36)
			if v41 == int32(0) {
				switch v34 {
				case 0, 4:
					v50 = *(*int32)(unsafe.Add(mBase, uint32(l0)+156))
					F_WebPSafeFree(m, v50)
					mBase = m.M
					*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = int64(0)
					v54 = int32(0)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v54
					v58 = base.Simd_g_const(&F_WebPPictureAlloc__k0)
					base.Simd_g_v128_store(m, l0+int32(24), v54, v58)
					*(*int32)(unsafe.Add(mBase, uint32(l0+int32(40)))) = v54
					v65 = base.I64_extend_i32_s(v35)
					v66 = int64(1)
					v69 = int64(base.Ui64(v65+v66) >> (uint(v66) % 64))
					v70 = base.I32_wrap_i64(v69)
					if v41|base.B2i32(v70 < int32(1)) != 0 {
						v85 = F_WebPEncodingSetError(m, l0, int32(5))
						mBase = m.M
						v131 = v85
					} else {
						v74 = base.I64_extend_i32_s(v38)
						v75 = int64(1)
						v76 = v74 + v75
						if int32(0) < base.I32_wrap_i64(int64(base.Ui64(v76)>>(uint(v75)%64))) {
							v90 = v34 << (uint(int32(29)) % 32) >> (uint(int32(31)) % 32) & v35
							v92 = base.I64_extend_i32_s(v90) * v74
							v93 = v74 * v65
							v95 = int64(1)
							v98 = v76 >> (uint(v95) % 64) * base.I64_extend32_s(v69)
							v103 = F_WebPSafeMalloc(m, v92+v93+v98<<(uint(v95)%64), int32(1))
							mBase = m.M
							if v103 != 0 {
								*(*int32)(unsafe.Add(mBase, uint32(l0)+40)) = v90
								*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = v70
								*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = v35
								*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v103
								*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v103
								v112 = v103 + base.I32_wrap_i64(v93)
								*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v112
								v114 = base.I32_wrap_i64(v98)
								v115 = v112 + v114
								*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v115
								if v92 == int64(0) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = v115 + v114
								}
								v131 = int32(1)
							} else {
								v105 = F_WebPEncodingSetError(m, l0, int32(1))
								mBase = m.M
								v131 = v105
							}
						} else {
							v85 = F_WebPEncodingSetError(m, l0, int32(5))
							mBase = m.M
							v131 = v85
						}
					}
				default:
					v46 = int32(4)
					v47 = F_WebPEncodingSetError(m, l0, v46)
					mBase = m.M
					if v47 != 0 {
						v50 = *(*int32)(unsafe.Add(mBase, uint32(l0)+156))
						F_WebPSafeFree(m, v50)
						mBase = m.M
						*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = int64(0)
						v54 = int32(0)
						*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v54
						v58 = base.Simd_g_const(&F_WebPPictureAlloc__k0)
						base.Simd_g_v128_store(m, l0+int32(24), v54, v58)
						*(*int32)(unsafe.Add(mBase, uint32(l0+int32(40)))) = v54
						v65 = base.I64_extend_i32_s(v35)
						v66 = int64(1)
						v69 = int64(base.Ui64(v65+v66) >> (uint(v66) % 64))
						v70 = base.I32_wrap_i64(v69)
						if v41|base.B2i32(v70 < int32(1)) != 0 {
							v85 = F_WebPEncodingSetError(m, l0, int32(5))
							mBase = m.M
							v131 = v85
						} else {
							v74 = base.I64_extend_i32_s(v38)
							v75 = int64(1)
							v76 = v74 + v75
							if int32(0) < base.I32_wrap_i64(int64(base.Ui64(v76)>>(uint(v75)%64))) {
								v90 = v34 << (uint(int32(29)) % 32) >> (uint(int32(31)) % 32) & v35
								v92 = base.I64_extend_i32_s(v90) * v74
								v93 = v74 * v65
								v95 = int64(1)
								v98 = v76 >> (uint(v95) % 64) * base.I64_extend32_s(v69)
								v103 = F_WebPSafeMalloc(m, v92+v93+v98<<(uint(v95)%64), int32(1))
								mBase = m.M
								if v103 != 0 {
									*(*int32)(unsafe.Add(mBase, uint32(l0)+40)) = v90
									*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = v70
									*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = v35
									*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v103
									*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v103
									v112 = v103 + base.I32_wrap_i64(v93)
									*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v112
									v114 = base.I32_wrap_i64(v98)
									v115 = v112 + v114
									*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v115
									if v92 == int64(0) {
									} else {
										*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = v115 + v114
									}
									v131 = int32(1)
								} else {
									v105 = F_WebPEncodingSetError(m, l0, int32(1))
									mBase = m.M
									v131 = v105
								}
							} else {
								v85 = F_WebPEncodingSetError(m, l0, int32(5))
								mBase = m.M
								v131 = v85
							}
						}
					} else {
						v131 = int32(0)
					}
				}
			} else {
				v46 = int32(5)
				v47 = F_WebPEncodingSetError(m, l0, v46)
				mBase = m.M
				if v47 != 0 {
					v50 = *(*int32)(unsafe.Add(mBase, uint32(l0)+156))
					F_WebPSafeFree(m, v50)
					mBase = m.M
					*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = int64(0)
					v54 = int32(0)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v54
					v58 = base.Simd_g_const(&F_WebPPictureAlloc__k0)
					base.Simd_g_v128_store(m, l0+int32(24), v54, v58)
					*(*int32)(unsafe.Add(mBase, uint32(l0+int32(40)))) = v54
					v65 = base.I64_extend_i32_s(v35)
					v66 = int64(1)
					v69 = int64(base.Ui64(v65+v66) >> (uint(v66) % 64))
					v70 = base.I32_wrap_i64(v69)
					if v41|base.B2i32(v70 < int32(1)) != 0 {
						v85 = F_WebPEncodingSetError(m, l0, int32(5))
						mBase = m.M
						v131 = v85
					} else {
						v74 = base.I64_extend_i32_s(v38)
						v75 = int64(1)
						v76 = v74 + v75
						if int32(0) < base.I32_wrap_i64(int64(base.Ui64(v76)>>(uint(v75)%64))) {
							v90 = v34 << (uint(int32(29)) % 32) >> (uint(int32(31)) % 32) & v35
							v92 = base.I64_extend_i32_s(v90) * v74
							v93 = v74 * v65
							v95 = int64(1)
							v98 = v76 >> (uint(v95) % 64) * base.I64_extend32_s(v69)
							v103 = F_WebPSafeMalloc(m, v92+v93+v98<<(uint(v95)%64), int32(1))
							mBase = m.M
							if v103 != 0 {
								*(*int32)(unsafe.Add(mBase, uint32(l0)+40)) = v90
								*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = v70
								*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = v35
								*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v103
								*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v103
								v112 = v103 + base.I32_wrap_i64(v93)
								*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v112
								v114 = base.I32_wrap_i64(v98)
								v115 = v112 + v114
								*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v115
								if v92 == int64(0) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = v115 + v114
								}
								v131 = int32(1)
							} else {
								v105 = F_WebPEncodingSetError(m, l0, int32(1))
								mBase = m.M
								v131 = v105
							}
						} else {
							v85 = F_WebPEncodingSetError(m, l0, int32(5))
							mBase = m.M
							v131 = v85
						}
					}
				} else {
					v131 = int32(0)
				}
			}
			return v131
		}
	} else {
		return int32(1)
	}
}

var F_WebPPictureAlloc__k0 = [2]uint64{0x0, 0x0}

func F_WebPPictureAllocYUVA(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v16 int32
	_ = v16
	var v19 int32
	_ = v19
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v31 int32
	_ = v31
	var v35 int32
	_ = v35
	var v39 base.V128
	_ = v39
	var v46 int64
	_ = v46
	var v47 int64
	_ = v47
	var v50 int64
	_ = v50
	var v51 int32
	_ = v51
	var v55 int64
	_ = v55
	var v56 int64
	_ = v56
	var v57 int64
	_ = v57
	var v66 int32
	_ = v66
	var v74 int32
	_ = v74
	var v76 int64
	_ = v76
	var v77 int64
	_ = v77
	var v79 int64
	_ = v79
	var v82 int64
	_ = v82
	var v85 int64
	_ = v85
	var v86 int32
	_ = v86
	var v93 int64
	_ = v93
	var v94 int32
	_ = v94
	var v95 int64
	_ = v95
	var v105 int32
	_ = v105
	var v107 int32
	_ = v107
	var v110 int32
	_ = v110
	var v120 int32
	_ = v120
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v13 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v14 = int32(1)
	v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v19 = base.B2i32(v13 < v14) | base.B2i32(v16 < v14)
	if v19 == int32(0) {
		switch v12 {
		case 0, 4:
			v31 = *(*int32)(unsafe.Add(mBase, uint32(l0)+156))
			F_free(m, v31)
			mBase = m.M
			*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = int64(0)
			v35 = int32(0)
			*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v35
			v39 = base.Simd_g_const(&F_WebPPictureAllocYUVA__k0)
			base.Simd_g_v128_store(m, l0+int32(24), v35, v39)
			*(*int32)(unsafe.Add(mBase, uint32(l0+int32(40)))) = v35
			v46 = base.I64_extend_i32_s(v13)
			v47 = int64(1)
			v50 = int64(base.Ui64(v46+v47) >> (uint(v47) % 64))
			v51 = base.I32_wrap_i64(v50)
			if v19|base.B2i32(v51 < int32(1)) != 0 {
				v66 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
				if v66 != 0 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(5)
				}
				return int32(0)
			} else {
				v55 = base.I64_extend_i32_s(v16)
				v56 = int64(1)
				v57 = v55 + v56
				if int32(0) < base.I32_wrap_i64(int64(base.Ui64(v57)>>(uint(v56)%64))) {
					v74 = v12 << (uint(int32(29)) % 32) >> (uint(int32(31)) % 32) & v13
					v76 = base.I64_extend_i32_s(v74) * v55
					v77 = v55 * v46
					v79 = int64(1)
					v82 = v57 >> (uint(v79) % 64) * base.I64_extend32_s(v50)
					v85 = v76 + v77 + v82<<(uint(v79)%64)
					v86 = int32(1)
					if v85 == int64(0) {
						v105 = F_malloc(m, base.I32_wrap_i64(v85)*v86)
						mBase = m.M
						v107 = v105
					} else {
						v93 = base.I64_div_u_s(int64(2147418112), v85)
						v94 = int32(0)
						v95 = base.I64_extend_i32_u(v86)
						if base.Ui64(int64(4294967295)) < base.Ui64(v95*v85) {
							v107 = v94
						} else {
							if base.Ui64(v93) < base.Ui64(v95) {
								v107 = v94
							} else {
								v105 = F_malloc(m, base.I32_wrap_i64(v85)*v86)
								mBase = m.M
								v107 = v105
							}
						}
					}
					if v107 != 0 {
						*(*int32)(unsafe.Add(mBase, uint32(l0)+40)) = v74
						*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = v51
						*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = v13
						*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v107
						*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v107
						v120 = v107 + base.I32_wrap_i64(v77)
						*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v120
						v122 = base.I32_wrap_i64(v82)
						v123 = v120 + v122
						*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v123
						if v76 == int64(0) {
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = v123 + v122
						}
						return int32(1)
					} else {
						v110 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
						if v110 != 0 {
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(1)
						}
						return int32(0)
					}
				} else {
					v66 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
					if v66 != 0 {
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(5)
					}
					return int32(0)
				}
			}
		default:
			v24 = int32(4)
			v25 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
			if v25 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v24
			}
			return int32(0)
		}
	} else {
		v24 = int32(5)
		v25 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
		if v25 != 0 {
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v24
		}
		return int32(0)
	}
}

var F_WebPPictureAllocYUVA__k0 = [2]uint64{0x0, 0x0}

func F_WebPPictureFree(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v6 int32
	_ = v6
	var v8 int64
	_ = v8
	var v16 base.V128
	_ = v16
	var v17 int32
	_ = v17
	if l0 == int32(0) {
	} else {
		v4 = *(*int32)(unsafe.Add(mBase, uint32(l0)+156))
		F_free(m, v4)
		mBase = m.M
		v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
		F_free(m, v6)
		mBase = m.M
		v8 = int64(0)
		*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = v8
		*(*int64)(unsafe.Add(mBase, uint32(l0)+156)) = v8
		*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = v8
		v16 = base.Simd_g_const(&F_WebPPictureFree__k0)
		v17 = int32(0)
		base.Simd_g_v128_store(m, l0+int32(24), v17, v16)
		*(*int32)(unsafe.Add(mBase, uint32(l0+int32(40)))) = v17
	}
	return
}

var F_WebPPictureFree__k0 = [2]uint64{0x0, 0x0}

func F_WebPPictureImportRGBA(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v16 int32
	_ = v16
	var v19 int32
	_ = v19
	var v23 int32
	_ = v23
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v41 int64
	_ = v41
	var v49 base.V128
	_ = v49
	var v50 int32
	_ = v50
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v62 int32
	_ = v62
	var v67 int32
	_ = v67
	var v70 int32
	_ = v70
	var v76 int32
	_ = v76
	var v79 int32
	_ = v79
	var v81 int32
	_ = v81
	var v84 int32
	_ = v84
	var v85 int32
	_ = v85
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v94 int32
	_ = v94
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v100 int32
	_ = v100
	var v102 int32
	_ = v102
	var v105 int32
	_ = v105
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v118 int32
	_ = v118
	v4 = int32(0)
	if l0 == v4 {
		v118 = v4
		return v118
	} else {
		if l1 == int32(0) {
			v118 = v4
			return v118
		} else {
			v16 = l2 >> (uint(int32(31)) % 32)
			v19 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
			if l2^v16-v16 < v19<<(uint(int32(2))%32) {
				v118 = v4
				return v118
			} else {
				v23 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				if v23 != 0 {
					v35 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					if l0 != 0 {
						v37 = *(*int32)(unsafe.Add(mBase, uint32(l0)+156))
						F_WebPSafeFree(m, v37)
						mBase = m.M
						v39 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
						F_WebPSafeFree(m, v39)
						mBase = m.M
						v41 = int64(0)
						*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = v41
						*(*int64)(unsafe.Add(mBase, uint32(l0)+156)) = v41
						*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = v41
						v49 = base.Simd_g_const(&F_WebPPictureImportRGBA__k0)
						v50 = int32(0)
						base.Simd_g_v128_store(m, l0+int32(24), v50, v49)
						*(*int32)(unsafe.Add(mBase, uint32(l0+int32(40)))) = v50
						v56 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						if v56 != 0 {
							v58 = F_WebPPictureAllocARGB(m, l0)
							mBase = m.M
							v59 = v58
						} else {
							v57 = F_WebPPictureAllocYUVA(m, l0)
							mBase = m.M
							v59 = v57
						}
					} else {
						v59 = int32(1)
					}
					if v59 == int32(0) {
						v118 = v4
					} else {
						v62 = int32(1)
						F_VP8LDspInit(m)
						mBase = m.M
						F_WebPInitAlphaProcessing(m)
						mBase = m.M
						if v35 < v62 {
							v118 = v62
						} else {
							v67 = int32(1)
							v70 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
							if v35 == v67 {
								v102 = l1
								v105 = v70
							} else {
								v76 = l1
								v79 = v70
								v81 = v35 & int32(-2)
								for {
									v84 = m.G25
									v85 = *(*int32)(unsafe.Add(mBase, uint32(v84)))
									m.T0[v85].(func(*base.Module, int32, int32, int32))(m, v76, v19, v79)
									mBase = m.M
									v87 = v76 + l2
									v88 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
									v89 = int32(2)
									v91 = v79 + v88<<(uint(v89)%32)
									v92 = *(*int32)(unsafe.Add(mBase, uint32(v84)))
									m.T0[v92].(func(*base.Module, int32, int32, int32))(m, v87, v19, v91)
									mBase = m.M
									v94 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
									v97 = v91 + v94<<(uint(v89)%32)
									v98 = v87 + l2
									v100 = v81 + int32(-2)
									if v100 != 0 {
										v76 = v98
										v79 = v97
										v81 = v100
										continue
									} else {
										break
									}
									break
								}
								v102 = v98
								v105 = v97
							}
							if v35&v67 == int32(0) {
								v118 = v67
							} else {
								v112 = m.G25
								v113 = *(*int32)(unsafe.Add(mBase, uint32(v112)))
								m.T0[v113].(func(*base.Module, int32, int32, int32))(m, v102, v19, v105)
								mBase = m.M
								v118 = v67
							}
						}
					}
					return v118
				} else {
					v33 = F_ImportYUVAFromRGBA(m, l1, l1+int32(1), l1+int32(2), l1+int32(3), int32(4), l2, float32(0), int32(0), l0)
					mBase = m.M
					return v33
				}
			}
		}
	}
}

var F_WebPPictureImportRGBA__k0 = [2]uint64{0x0, 0x0}

func F_WebPPictureResetBuffers(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int64
	_ = v2
	var v10 base.V128
	_ = v10
	var v11 int32
	_ = v11
	v2 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = v2
	*(*int64)(unsafe.Add(mBase, uint32(l0)+156)) = v2
	*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = v2
	v10 = base.Simd_g_const(&F_WebPPictureResetBuffers__k0)
	v11 = int32(0)
	base.Simd_g_v128_store(m, l0+int32(24), v11, v10)
	*(*int32)(unsafe.Add(mBase, uint32(l0+int32(40)))) = v11
	return
}

var F_WebPPictureResetBuffers__k0 = [2]uint64{0x0, 0x0}

func F_WebPPictureView(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) int32 {
	mBase := m.M
	_ = mBase
	var v7 int32
	_ = v7
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v21 int32
	_ = v21
	var v30 int32
	_ = v30
	var v33 int32
	_ = v33
	var v37 int32
	_ = v37
	var v38 int64
	_ = v38
	var v46 base.V128
	_ = v46
	var v47 int32
	_ = v47
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v57 int32
	_ = v57
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	var v66 int32
	_ = v66
	var v67 int32
	_ = v67
	var v70 int32
	_ = v70
	var v73 int32
	_ = v73
	var v76 int32
	_ = v76
	var v80 int32
	_ = v80
	var v83 int32
	_ = v83
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v92 int32
	_ = v92
	var v100 int32
	_ = v100
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v105 int32
	_ = v105
	var v116 int32
	_ = v116
	v7 = int32(0)
	if l0 == v7 {
		v116 = v7
	} else {
		if l5 == int32(0) {
			v116 = v7
		} else {
			v14 = int32(0)
			v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
			if v17 != 0 {
				v18 = l1
			} else {
				v18 = l1 & int32(-2)
			}
			if v17 != 0 {
				v21 = l2
			} else {
				v21 = l2 & int32(-2)
			}
			if v18|v21 < int32(0) {
				v116 = v14
			} else {
				if l3 < int32(1) {
					v116 = v14
				} else {
					if l4 < int32(1) {
						v116 = v14
					} else {
						v30 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
						if v30 < v18+l3 {
							v116 = v14
						} else {
							v33 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							if v33 < v21+l4 {
								v116 = v14
							} else {
								if l0 == l5 {
									v54 = v17
								} else {
									v37 = F_memcpy(m, l5, l0, int32(172))
									mBase = m.M
									v38 = int64(0)
									*(*int64)(unsafe.Add(mBase, uint32(v37)+52)) = v38
									*(*int64)(unsafe.Add(mBase, uint32(v37)+156)) = v38
									*(*int64)(unsafe.Add(mBase, uint32(v37)+16)) = v38
									v46 = base.Simd_g_const(&F_WebPPictureView__k0)
									v47 = int32(0)
									base.Simd_g_v128_store(m, v37+int32(24), v47, v46)
									*(*int32)(unsafe.Add(mBase, uint32(v37+int32(40)))) = v47
									v53 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
									v54 = v53
								}
								*(*int32)(unsafe.Add(mBase, uint32(l5)+12)) = l4
								*(*int32)(unsafe.Add(mBase, uint32(l5)+8)) = l3
								if v54 != 0 {
									v89 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
									v90 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
									v92 = int32(2)
									v100 = v90
									v101 = int32(52)
									v102 = int32(56)
									v105 = v89 + v90*v21<<(uint(v92)%32) + v18<<(uint(v92)%32)
									*(*int32)(unsafe.Add(mBase, uint32(l5+v102))) = v100
									*(*int32)(unsafe.Add(mBase, uint32(l5+v101))) = v105
									v116 = int32(1)
								} else {
									v57 = *(*int32)(unsafe.Add(mBase, uint32(l0)+32))
									*(*int32)(unsafe.Add(mBase, uint32(l5)+32)) = v57
									v59 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
									*(*int32)(unsafe.Add(mBase, uint32(l5)+28)) = v59
									v61 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									*(*int32)(unsafe.Add(mBase, uint32(l5)+16)) = v61 + v59*v21 + v18
									v66 = int32(1)
									v67 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
									v70 = v57 * int32(base.Ui32(v21)>>(uint(v66)%32))
									v73 = int32(base.Ui32(v18) >> (uint(v66) % 32))
									*(*int32)(unsafe.Add(mBase, uint32(l5)+20)) = v67 + v70 + v73
									v76 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
									*(*int32)(unsafe.Add(mBase, uint32(l5)+24)) = v76 + v70 + v73
									v80 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
									if v80 == int32(0) {
										v116 = v66
									} else {
										v83 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
										v100 = v83
										v101 = int32(36)
										v102 = int32(40)
										v105 = v80 + v83*v21 + v18
										*(*int32)(unsafe.Add(mBase, uint32(l5+v102))) = v100
										*(*int32)(unsafe.Add(mBase, uint32(l5+v101))) = v105
										v116 = int32(1)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return v116
}

var F_WebPPictureView__k0 = [2]uint64{0x0, 0x0}

func F_WebPPictureYUVAToARGB(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v28 int32
	_ = v28
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v39 int32
	_ = v39
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v61 int32
	_ = v61
	var v73 int32
	_ = v73
	var v75 int32
	_ = v75
	var v86 int32
	_ = v86
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v94 int32
	_ = v94
	var v95 int32
	_ = v95
	var v96 int32
	_ = v96
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
	var v106 int32
	_ = v106
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v116 int32
	_ = v116
	var v118 int32
	_ = v118
	var v119 int32
	_ = v119
	var v123 int32
	_ = v123
	var v134 int32
	_ = v134
	var v135 int32
	_ = v135
	var v136 int32
	_ = v136
	var v138 int32
	_ = v138
	var v140 int32
	_ = v140
	var v143 int32
	_ = v143
	var v145 int32
	_ = v145
	var v149 int32
	_ = v149
	var v150 int32
	_ = v150
	var v151 int32
	_ = v151
	var v154 int32
	_ = v154
	var v172 int32
	_ = v172
	var v175 int32
	_ = v175
	var v176 int32
	_ = v176
	var v183 int32
	_ = v183
	var v189 int32
	_ = v189
	var v190 int32
	_ = v190
	var v192 int32
	_ = v192
	var v196 int32
	_ = v196
	var v206 int32
	_ = v206
	var v210 int32
	_ = v210
	var v220 int32
	_ = v220
	var v221 int32
	_ = v221
	var v222 int32
	_ = v222
	var v223 int32
	_ = v223
	var v226 int32
	_ = v226
	var v227 int32
	_ = v227
	var v238 int32
	_ = v238
	var v239 int32
	_ = v239
	var v240 int32
	_ = v240
	var v257 int32
	_ = v257
	var v258 base.V128
	_ = v258
	var v264 base.V128
	_ = v264
	var v267 base.V128
	_ = v267
	var v275 int32
	_ = v275
	var v279 int32
	_ = v279
	var v298 int32
	_ = v298
	var v307 int32
	_ = v307
	var v309 int32
	_ = v309
	var v320 int32
	_ = v320
	var v321 int32
	_ = v321
	var v322 int32
	_ = v322
	var v341 int32
	_ = v341
	var v347 int32
	_ = v347
	var v354 int32
	_ = v354
	var v377 int32
	_ = v377
	var v379 int32
	_ = v379
	var v381 int32
	_ = v381
	var v383 int32
	_ = v383
	var v387 int32
	_ = v387
	v2 = int32(0)
	if l0 == v2 {
		v387 = v2
	} else {
		v24 = int32(3)
		v25 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
		if v25 == int32(0) {
			v381 = v24
			v383 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
			if v383 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v381
			}
			v387 = int32(0)
		} else {
			v28 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
			if v28 == int32(0) {
				v381 = v24
				v383 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
				if v383 != 0 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v381
				}
				v387 = int32(0)
			} else {
				v31 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
				if v31 == int32(0) {
					v381 = v24
					v383 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
					if v383 != 0 {
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v381
					}
					v387 = int32(0)
				} else {
					v34 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					if v34&int32(4) == int32(0) {
						if v34&int32(3) != 0 {
							v381 = int32(4)
							v383 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
							if v383 != 0 {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v381
							}
							v387 = int32(0)
						} else {
							v48 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							v49 = int32(5)
							v50 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
							if v50 < int32(1) {
								v57 = v49
								v58 = F_WebPEncodingSetError(m, l0, v57)
								mBase = m.M
								if v58 != 0 {
									v61 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
									F_WebPSafeFree(m, v61)
									mBase = m.M
									*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
									*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
									v73 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v48)*base.I64_extend_i32_s(v50)+int64(31), int32(4))
									mBase = m.M
									if v73 != 0 {
										*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v50
										*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v73
										*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v73 + int32(31)) & int32(-32)
										v86 = int32(1)
									} else {
										v75 = F_WebPEncodingSetError(m, l0, int32(1))
										mBase = m.M
										v86 = v75
									}
								} else {
									v86 = int32(0)
								}
							} else {
								if v48 < int32(1) {
									v57 = v49
									v58 = F_WebPEncodingSetError(m, l0, v57)
									mBase = m.M
									if v58 != 0 {
										v61 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
										F_WebPSafeFree(m, v61)
										mBase = m.M
										*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
										*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
										v73 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v48)*base.I64_extend_i32_s(v50)+int64(31), int32(4))
										mBase = m.M
										if v73 != 0 {
											*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v50
											*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v73
											*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v73 + int32(31)) & int32(-32)
											v86 = int32(1)
										} else {
											v75 = F_WebPEncodingSetError(m, l0, int32(1))
											mBase = m.M
											v86 = v75
										}
									} else {
										v86 = int32(0)
									}
								} else {
									v56 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
									switch v56 {
									case 0, 4:
										v61 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
										F_WebPSafeFree(m, v61)
										mBase = m.M
										*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
										*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
										v73 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v48)*base.I64_extend_i32_s(v50)+int64(31), int32(4))
										mBase = m.M
										if v73 != 0 {
											*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v50
											*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v73
											*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v73 + int32(31)) & int32(-32)
											v86 = int32(1)
										} else {
											v75 = F_WebPEncodingSetError(m, l0, int32(1))
											mBase = m.M
											v86 = v75
										}
									default:
										v57 = int32(4)
										v58 = F_WebPEncodingSetError(m, l0, v57)
										mBase = m.M
										if v58 != 0 {
											v61 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
											F_WebPSafeFree(m, v61)
											mBase = m.M
											*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
											*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
											v73 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v48)*base.I64_extend_i32_s(v50)+int64(31), int32(4))
											mBase = m.M
											if v73 != 0 {
												*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v50
												*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v73
												*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v73 + int32(31)) & int32(-32)
												v86 = int32(1)
											} else {
												v75 = F_WebPEncodingSetError(m, l0, int32(1))
												mBase = m.M
												v86 = v75
											}
										} else {
											v86 = int32(0)
										}
									}
								}
							}
							if v86 == int32(0) {
								v387 = v2
							} else {
								v89 = int32(1)
								*(*int32)(unsafe.Add(mBase, uint32(l0))) = v89
								v91 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
								v92 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
								v93 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
								v94 = int32(0)
								v95 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
								v96 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
								v97 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
								v99 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
								v101 = F_WebPGetLinePairConverter(m, v89)
								mBase = m.M
								m.T0[v101].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32, int32))(m, v93, v94, v95, v96, v95, v96, v97, v94, v99)
								mBase = m.M
								v104 = v92 << (uint(int32(2)) % 32)
								v105 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
								v106 = v93 + v105
								if int32(3) <= v91 {
									v114 = v106
									v115 = v95
									v116 = v96
									v118 = int32(2)
									v119 = v97
									v123 = v105
									for {
										v134 = *(*int32)(unsafe.Add(mBase, uint32(l0)+32))
										v135 = v115 + v134
										v136 = v116 + v134
										v138 = v119 + v92<<(uint(int32(3))%32)
										m.T0[v101].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32, int32))(m, v114, v114+v123, v115, v116, v135, v136, v119+v104, v138, v99)
										mBase = m.M
										v140 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
										v143 = v114 + v140<<(uint(int32(1))%32)
										v145 = v118 + int32(2)
										if v145 < v91 {
											v114 = v143
											v115 = v135
											v116 = v136
											v118 = v145
											v119 = v138
											v123 = v140
											continue
										} else {
											break
										}
										break
									}
									v149 = v143
									v150 = v135
									v151 = v136
									v154 = v138 + v104
								} else {
									v149 = v106
									v150 = v95
									v151 = v96
									v154 = v97 + v104
								}
								if v91 < int32(2) {
								} else {
									if v91&int32(1) != 0 {
									} else {
										v172 = int32(0)
										m.T0[v101].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32, int32))(m, v149, v172, v150, v151, v150, v151, v154, v172, v99)
										mBase = m.M
									}
								}
								v175 = int32(1)
								v176 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+4)))
								if v176&int32(4) == int32(0) {
									v387 = v175
								} else {
									if v91 < int32(1) {
										v387 = v175
									} else {
										v183 = int32(1)
										if v99 < v183 {
											v387 = v183
										} else {
											v189 = v99 & int32(2147483644)
											v190 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
											v192 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
											v196 = int32(0)
											v206 = v196
											v210 = v196
											for {
												v220 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
												v221 = v220 * v206
												v222 = v190 + v221
												v223 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
												v226 = v223 * v206 << (uint(int32(2)) % 32)
												v227 = v192 + v226
												if base.Ui32(v99) < base.Ui32(int32(4)) {
													v279 = int32(0)
													v298 = v279 | int32(1)
													if v99&int32(1) == int32(0) {
														v309 = v279
													} else {
														v307 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v222+v279))))
														*(*uint8)(unsafe.Add(mBase, uint32(v227+v279<<(uint(int32(2))%32)+int32(3)))) = uint8(v307)
														v309 = v298
													}
													if v99 == v298 {
													} else {
														v320 = v192 + (v223*v210 + v309<<(uint(int32(2))%32))
														v321 = v190 + (v309 + v221)
														v322 = v99 - v309
														for {
															v341 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321))))
															*(*uint8)(unsafe.Add(mBase, uint32(v320+int32(3)))) = uint8(v341)
															v347 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321+int32(1)))))
															*(*uint8)(unsafe.Add(mBase, uint32(v320+int32(7)))) = uint8(v347)
															v354 = v322 + int32(-2)
															if v354 != 0 {
																v320 = v320 + int32(8)
																v321 = v321 + int32(2)
																v322 = v354
																continue
															} else {
																break
															}
															break
														}
													}
												} else {
													if base.Ui32(v190+v99+v221) <= base.Ui32(v227) {
														v238 = v192 + v223*v210
														v239 = v222
														v240 = v189
														for {
															v257 = int32(0)
															v258 = base.Simd_g_v128_load32_zero(m, v239, v257)
															v264 = base.Simd_g_v128_load(m, v238, v257)
															v267 = base.Simd_g_v128_or(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v258, base.Simd_g_const(&F_WebPPictureYUVAToARGB__k0)), int32(24)), base.Simd_g_v128_and(v264, base.Simd_g_const(&F_WebPPictureYUVAToARGB__k1)))
															base.Simd_g_v128_store(m, v238, v257, v267)
															v275 = v240 + int32(-4)
															if v275 != 0 {
																v238 = v238 + int32(16)
																v239 = v239 + int32(4)
																v240 = v275
																continue
															} else {
																break
															}
															break
														}
														if v99 == v189 {
														} else {
															v279 = v189
															v298 = v279 | int32(1)
															if v99&int32(1) == int32(0) {
																v309 = v279
															} else {
																v307 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v222+v279))))
																*(*uint8)(unsafe.Add(mBase, uint32(v227+v279<<(uint(int32(2))%32)+int32(3)))) = uint8(v307)
																v309 = v298
															}
															if v99 == v298 {
															} else {
																v320 = v192 + (v223*v210 + v309<<(uint(int32(2))%32))
																v321 = v190 + (v309 + v221)
																v322 = v99 - v309
																for {
																	v341 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321))))
																	*(*uint8)(unsafe.Add(mBase, uint32(v320+int32(3)))) = uint8(v341)
																	v347 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321+int32(1)))))
																	*(*uint8)(unsafe.Add(mBase, uint32(v320+int32(7)))) = uint8(v347)
																	v354 = v322 + int32(-2)
																	if v354 != 0 {
																		v320 = v320 + int32(8)
																		v321 = v321 + int32(2)
																		v322 = v354
																		continue
																	} else {
																		break
																	}
																	break
																}
															}
														}
													} else {
														if base.Ui32(v222) < base.Ui32(v192+v99<<(uint(int32(2))%32)+v226) {
															v279 = int32(0)
															v298 = v279 | int32(1)
															if v99&int32(1) == int32(0) {
																v309 = v279
															} else {
																v307 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v222+v279))))
																*(*uint8)(unsafe.Add(mBase, uint32(v227+v279<<(uint(int32(2))%32)+int32(3)))) = uint8(v307)
																v309 = v298
															}
															if v99 == v298 {
															} else {
																v320 = v192 + (v223*v210 + v309<<(uint(int32(2))%32))
																v321 = v190 + (v309 + v221)
																v322 = v99 - v309
																for {
																	v341 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321))))
																	*(*uint8)(unsafe.Add(mBase, uint32(v320+int32(3)))) = uint8(v341)
																	v347 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321+int32(1)))))
																	*(*uint8)(unsafe.Add(mBase, uint32(v320+int32(7)))) = uint8(v347)
																	v354 = v322 + int32(-2)
																	if v354 != 0 {
																		v320 = v320 + int32(8)
																		v321 = v321 + int32(2)
																		v322 = v354
																		continue
																	} else {
																		break
																	}
																	break
																}
															}
														} else {
															v238 = v192 + v223*v210
															v239 = v222
															v240 = v189
															for {
																v257 = int32(0)
																v258 = base.Simd_g_v128_load32_zero(m, v239, v257)
																v264 = base.Simd_g_v128_load(m, v238, v257)
																v267 = base.Simd_g_v128_or(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v258, base.Simd_g_const(&F_WebPPictureYUVAToARGB__k0)), int32(24)), base.Simd_g_v128_and(v264, base.Simd_g_const(&F_WebPPictureYUVAToARGB__k1)))
																base.Simd_g_v128_store(m, v238, v257, v267)
																v275 = v240 + int32(-4)
																if v275 != 0 {
																	v238 = v238 + int32(16)
																	v239 = v239 + int32(4)
																	v240 = v275
																	continue
																} else {
																	break
																}
																break
															}
															if v99 == v189 {
															} else {
																v279 = v189
																v298 = v279 | int32(1)
																if v99&int32(1) == int32(0) {
																	v309 = v279
																} else {
																	v307 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v222+v279))))
																	*(*uint8)(unsafe.Add(mBase, uint32(v227+v279<<(uint(int32(2))%32)+int32(3)))) = uint8(v307)
																	v309 = v298
																}
																if v99 == v298 {
																} else {
																	v320 = v192 + (v223*v210 + v309<<(uint(int32(2))%32))
																	v321 = v190 + (v309 + v221)
																	v322 = v99 - v309
																	for {
																		v341 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321))))
																		*(*uint8)(unsafe.Add(mBase, uint32(v320+int32(3)))) = uint8(v341)
																		v347 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321+int32(1)))))
																		*(*uint8)(unsafe.Add(mBase, uint32(v320+int32(7)))) = uint8(v347)
																		v354 = v322 + int32(-2)
																		if v354 != 0 {
																			v320 = v320 + int32(8)
																			v321 = v321 + int32(2)
																			v322 = v354
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
												v377 = int32(1)
												v379 = v206 + v377
												if v379 != v91 {
													v206 = v379
													v210 = v210 + int32(4)
													continue
												} else {
													break
												}
												break
											}
											v387 = v377
										}
									}
								}
							}
						}
					} else {
						v39 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
						if v39 == int32(0) {
							v381 = v24
							v383 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
							if v383 != 0 {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v381
							}
							v387 = int32(0)
						} else {
							if v34&int32(3) != 0 {
								v381 = int32(4)
								v383 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
								if v383 != 0 {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v381
								}
								v387 = int32(0)
							} else {
								v48 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
								v49 = int32(5)
								v50 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
								if v50 < int32(1) {
									v57 = v49
									v58 = F_WebPEncodingSetError(m, l0, v57)
									mBase = m.M
									if v58 != 0 {
										v61 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
										F_WebPSafeFree(m, v61)
										mBase = m.M
										*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
										*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
										v73 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v48)*base.I64_extend_i32_s(v50)+int64(31), int32(4))
										mBase = m.M
										if v73 != 0 {
											*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v50
											*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v73
											*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v73 + int32(31)) & int32(-32)
											v86 = int32(1)
										} else {
											v75 = F_WebPEncodingSetError(m, l0, int32(1))
											mBase = m.M
											v86 = v75
										}
									} else {
										v86 = int32(0)
									}
								} else {
									if v48 < int32(1) {
										v57 = v49
										v58 = F_WebPEncodingSetError(m, l0, v57)
										mBase = m.M
										if v58 != 0 {
											v61 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
											F_WebPSafeFree(m, v61)
											mBase = m.M
											*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
											*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
											v73 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v48)*base.I64_extend_i32_s(v50)+int64(31), int32(4))
											mBase = m.M
											if v73 != 0 {
												*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v50
												*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v73
												*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v73 + int32(31)) & int32(-32)
												v86 = int32(1)
											} else {
												v75 = F_WebPEncodingSetError(m, l0, int32(1))
												mBase = m.M
												v86 = v75
											}
										} else {
											v86 = int32(0)
										}
									} else {
										v56 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
										switch v56 {
										case 0, 4:
											v61 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
											F_WebPSafeFree(m, v61)
											mBase = m.M
											*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
											*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
											v73 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v48)*base.I64_extend_i32_s(v50)+int64(31), int32(4))
											mBase = m.M
											if v73 != 0 {
												*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v50
												*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v73
												*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v73 + int32(31)) & int32(-32)
												v86 = int32(1)
											} else {
												v75 = F_WebPEncodingSetError(m, l0, int32(1))
												mBase = m.M
												v86 = v75
											}
										default:
											v57 = int32(4)
											v58 = F_WebPEncodingSetError(m, l0, v57)
											mBase = m.M
											if v58 != 0 {
												v61 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
												F_WebPSafeFree(m, v61)
												mBase = m.M
												*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
												*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
												v73 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v48)*base.I64_extend_i32_s(v50)+int64(31), int32(4))
												mBase = m.M
												if v73 != 0 {
													*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v50
													*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v73
													*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v73 + int32(31)) & int32(-32)
													v86 = int32(1)
												} else {
													v75 = F_WebPEncodingSetError(m, l0, int32(1))
													mBase = m.M
													v86 = v75
												}
											} else {
												v86 = int32(0)
											}
										}
									}
								}
								if v86 == int32(0) {
									v387 = v2
								} else {
									v89 = int32(1)
									*(*int32)(unsafe.Add(mBase, uint32(l0))) = v89
									v91 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
									v92 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
									v93 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									v94 = int32(0)
									v95 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
									v96 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
									v97 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
									v99 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
									v101 = F_WebPGetLinePairConverter(m, v89)
									mBase = m.M
									m.T0[v101].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32, int32))(m, v93, v94, v95, v96, v95, v96, v97, v94, v99)
									mBase = m.M
									v104 = v92 << (uint(int32(2)) % 32)
									v105 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
									v106 = v93 + v105
									if int32(3) <= v91 {
										v114 = v106
										v115 = v95
										v116 = v96
										v118 = int32(2)
										v119 = v97
										v123 = v105
										for {
											v134 = *(*int32)(unsafe.Add(mBase, uint32(l0)+32))
											v135 = v115 + v134
											v136 = v116 + v134
											v138 = v119 + v92<<(uint(int32(3))%32)
											m.T0[v101].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32, int32))(m, v114, v114+v123, v115, v116, v135, v136, v119+v104, v138, v99)
											mBase = m.M
											v140 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
											v143 = v114 + v140<<(uint(int32(1))%32)
											v145 = v118 + int32(2)
											if v145 < v91 {
												v114 = v143
												v115 = v135
												v116 = v136
												v118 = v145
												v119 = v138
												v123 = v140
												continue
											} else {
												break
											}
											break
										}
										v149 = v143
										v150 = v135
										v151 = v136
										v154 = v138 + v104
									} else {
										v149 = v106
										v150 = v95
										v151 = v96
										v154 = v97 + v104
									}
									if v91 < int32(2) {
									} else {
										if v91&int32(1) != 0 {
										} else {
											v172 = int32(0)
											m.T0[v101].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32, int32))(m, v149, v172, v150, v151, v150, v151, v154, v172, v99)
											mBase = m.M
										}
									}
									v175 = int32(1)
									v176 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+4)))
									if v176&int32(4) == int32(0) {
										v387 = v175
									} else {
										if v91 < int32(1) {
											v387 = v175
										} else {
											v183 = int32(1)
											if v99 < v183 {
												v387 = v183
											} else {
												v189 = v99 & int32(2147483644)
												v190 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
												v192 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
												v196 = int32(0)
												v206 = v196
												v210 = v196
												for {
													v220 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
													v221 = v220 * v206
													v222 = v190 + v221
													v223 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
													v226 = v223 * v206 << (uint(int32(2)) % 32)
													v227 = v192 + v226
													if base.Ui32(v99) < base.Ui32(int32(4)) {
														v279 = int32(0)
														v298 = v279 | int32(1)
														if v99&int32(1) == int32(0) {
															v309 = v279
														} else {
															v307 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v222+v279))))
															*(*uint8)(unsafe.Add(mBase, uint32(v227+v279<<(uint(int32(2))%32)+int32(3)))) = uint8(v307)
															v309 = v298
														}
														if v99 == v298 {
														} else {
															v320 = v192 + (v223*v210 + v309<<(uint(int32(2))%32))
															v321 = v190 + (v309 + v221)
															v322 = v99 - v309
															for {
																v341 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321))))
																*(*uint8)(unsafe.Add(mBase, uint32(v320+int32(3)))) = uint8(v341)
																v347 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321+int32(1)))))
																*(*uint8)(unsafe.Add(mBase, uint32(v320+int32(7)))) = uint8(v347)
																v354 = v322 + int32(-2)
																if v354 != 0 {
																	v320 = v320 + int32(8)
																	v321 = v321 + int32(2)
																	v322 = v354
																	continue
																} else {
																	break
																}
																break
															}
														}
													} else {
														if base.Ui32(v190+v99+v221) <= base.Ui32(v227) {
															v238 = v192 + v223*v210
															v239 = v222
															v240 = v189
															for {
																v257 = int32(0)
																v258 = base.Simd_g_v128_load32_zero(m, v239, v257)
																v264 = base.Simd_g_v128_load(m, v238, v257)
																v267 = base.Simd_g_v128_or(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v258, base.Simd_g_const(&F_WebPPictureYUVAToARGB__k0)), int32(24)), base.Simd_g_v128_and(v264, base.Simd_g_const(&F_WebPPictureYUVAToARGB__k1)))
																base.Simd_g_v128_store(m, v238, v257, v267)
																v275 = v240 + int32(-4)
																if v275 != 0 {
																	v238 = v238 + int32(16)
																	v239 = v239 + int32(4)
																	v240 = v275
																	continue
																} else {
																	break
																}
																break
															}
															if v99 == v189 {
															} else {
																v279 = v189
																v298 = v279 | int32(1)
																if v99&int32(1) == int32(0) {
																	v309 = v279
																} else {
																	v307 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v222+v279))))
																	*(*uint8)(unsafe.Add(mBase, uint32(v227+v279<<(uint(int32(2))%32)+int32(3)))) = uint8(v307)
																	v309 = v298
																}
																if v99 == v298 {
																} else {
																	v320 = v192 + (v223*v210 + v309<<(uint(int32(2))%32))
																	v321 = v190 + (v309 + v221)
																	v322 = v99 - v309
																	for {
																		v341 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321))))
																		*(*uint8)(unsafe.Add(mBase, uint32(v320+int32(3)))) = uint8(v341)
																		v347 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321+int32(1)))))
																		*(*uint8)(unsafe.Add(mBase, uint32(v320+int32(7)))) = uint8(v347)
																		v354 = v322 + int32(-2)
																		if v354 != 0 {
																			v320 = v320 + int32(8)
																			v321 = v321 + int32(2)
																			v322 = v354
																			continue
																		} else {
																			break
																		}
																		break
																	}
																}
															}
														} else {
															if base.Ui32(v222) < base.Ui32(v192+v99<<(uint(int32(2))%32)+v226) {
																v279 = int32(0)
																v298 = v279 | int32(1)
																if v99&int32(1) == int32(0) {
																	v309 = v279
																} else {
																	v307 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v222+v279))))
																	*(*uint8)(unsafe.Add(mBase, uint32(v227+v279<<(uint(int32(2))%32)+int32(3)))) = uint8(v307)
																	v309 = v298
																}
																if v99 == v298 {
																} else {
																	v320 = v192 + (v223*v210 + v309<<(uint(int32(2))%32))
																	v321 = v190 + (v309 + v221)
																	v322 = v99 - v309
																	for {
																		v341 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321))))
																		*(*uint8)(unsafe.Add(mBase, uint32(v320+int32(3)))) = uint8(v341)
																		v347 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321+int32(1)))))
																		*(*uint8)(unsafe.Add(mBase, uint32(v320+int32(7)))) = uint8(v347)
																		v354 = v322 + int32(-2)
																		if v354 != 0 {
																			v320 = v320 + int32(8)
																			v321 = v321 + int32(2)
																			v322 = v354
																			continue
																		} else {
																			break
																		}
																		break
																	}
																}
															} else {
																v238 = v192 + v223*v210
																v239 = v222
																v240 = v189
																for {
																	v257 = int32(0)
																	v258 = base.Simd_g_v128_load32_zero(m, v239, v257)
																	v264 = base.Simd_g_v128_load(m, v238, v257)
																	v267 = base.Simd_g_v128_or(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v258, base.Simd_g_const(&F_WebPPictureYUVAToARGB__k0)), int32(24)), base.Simd_g_v128_and(v264, base.Simd_g_const(&F_WebPPictureYUVAToARGB__k1)))
																	base.Simd_g_v128_store(m, v238, v257, v267)
																	v275 = v240 + int32(-4)
																	if v275 != 0 {
																		v238 = v238 + int32(16)
																		v239 = v239 + int32(4)
																		v240 = v275
																		continue
																	} else {
																		break
																	}
																	break
																}
																if v99 == v189 {
																} else {
																	v279 = v189
																	v298 = v279 | int32(1)
																	if v99&int32(1) == int32(0) {
																		v309 = v279
																	} else {
																		v307 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v222+v279))))
																		*(*uint8)(unsafe.Add(mBase, uint32(v227+v279<<(uint(int32(2))%32)+int32(3)))) = uint8(v307)
																		v309 = v298
																	}
																	if v99 == v298 {
																	} else {
																		v320 = v192 + (v223*v210 + v309<<(uint(int32(2))%32))
																		v321 = v190 + (v309 + v221)
																		v322 = v99 - v309
																		for {
																			v341 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321))))
																			*(*uint8)(unsafe.Add(mBase, uint32(v320+int32(3)))) = uint8(v341)
																			v347 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v321+int32(1)))))
																			*(*uint8)(unsafe.Add(mBase, uint32(v320+int32(7)))) = uint8(v347)
																			v354 = v322 + int32(-2)
																			if v354 != 0 {
																				v320 = v320 + int32(8)
																				v321 = v321 + int32(2)
																				v322 = v354
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
													v377 = int32(1)
													v379 = v206 + v377
													if v379 != v91 {
														v206 = v379
														v210 = v210 + int32(4)
														continue
													} else {
														break
													}
													break
												}
												v387 = v377
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
	return v387
}

var F_WebPPictureYUVAToARGB__k0 = [2]uint64{0x100000000, 0x300000002}
var F_WebPPictureYUVAToARGB__k1 = [2]uint64{0xffffff00ffffff, 0xffffff00ffffff}
