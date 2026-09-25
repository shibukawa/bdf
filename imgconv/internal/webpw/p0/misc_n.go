//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	"unsafe"
)

func F_NearLossless(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32) {
	mBase := m.M
	_ = mBase
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v38 int32
	_ = v38
	var v42 int32
	_ = v42
	var v45 int32
	_ = v45
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v55 int32
	_ = v55
	var v58 int32
	_ = v58
	var v63 int32
	_ = v63
	var v66 int32
	_ = v66
	var __phi66 int32
	_ = __phi66
	var v69 int32
	_ = v69
	var __phi69 int32
	_ = __phi69
	var v70 int32
	_ = v70
	var __phi70 int32
	_ = __phi70
	var v72 int32
	_ = v72
	var __phi72 int32
	_ = __phi72
	var v73 int32
	_ = v73
	var __phi73 int32
	_ = __phi73
	var v75 int32
	_ = v75
	var __phi75 int32
	_ = __phi75
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v104 int32
	_ = v104
	var v106 int32
	_ = v106
	var v111 int32
	_ = v111
	var v133 int32
	_ = v133
	var v140 int32
	_ = v140
	var v143 int32
	_ = v143
	var v146 int32
	_ = v146
	var v147 int32
	_ = v147
	var v149 int32
	_ = v149
	var v150 int32
	_ = v150
	var v153 int32
	_ = v153
	var v160 int32
	_ = v160
	var v163 int32
	_ = v163
	var v165 int32
	_ = v165
	var v166 int32
	_ = v166
	var v171 int32
	_ = v171
	var v174 int32
	_ = v174
	var v175 int32
	_ = v175
	var v178 int32
	_ = v178
	var v183 int32
	_ = v183
	var v186 int32
	_ = v186
	var v193 int32
	_ = v193
	var v200 int32
	_ = v200
	var v205 int32
	_ = v205
	var v209 int32
	_ = v209
	var v212 int32
	_ = v212
	var v219 int32
	_ = v219
	var v226 int32
	_ = v226
	var v231 int32
	_ = v231
	var v235 int32
	_ = v235
	var v238 int32
	_ = v238
	var v245 int32
	_ = v245
	var v252 int32
	_ = v252
	var v257 int32
	_ = v257
	var v265 int32
	_ = v265
	var v266 int32
	_ = v266
	var v271 int32
	_ = v271
	var v277 int32
	_ = v277
	var v278 int32
	_ = v278
	var v283 int32
	_ = v283
	var v287 int32
	_ = v287
	var v290 int32
	_ = v290
	var v292 int32
	_ = v292
	var v293 int32
	_ = v293
	var v298 int32
	_ = v298
	var v304 int32
	_ = v304
	var v311 int32
	_ = v311
	var v317 int32
	_ = v317
	var v319 int32
	_ = v319
	var v326 int32
	_ = v326
	var v332 int32
	_ = v332
	var v333 int32
	_ = v333
	var v368 int32
	_ = v368
	var v370 int32
	_ = v370
	var v373 int32
	_ = v373
	var v374 int32
	_ = v374
	var v376 int32
	_ = v376
	var v378 int32
	_ = v378
	var v380 int32
	_ = v380
	var v388 int32
	_ = v388
	var __phi388 int32
	_ = __phi388
	var v391 int32
	_ = v391
	var __phi391 int32
	_ = __phi391
	var v392 int32
	_ = v392
	var __phi392 int32
	_ = __phi392
	var v394 int32
	_ = v394
	var __phi394 int32
	_ = __phi394
	var v395 int32
	_ = v395
	var __phi395 int32
	_ = __phi395
	var v396 int32
	_ = v396
	var __phi396 int32
	_ = __phi396
	var v415 int32
	_ = v415
	var v416 int32
	_ = v416
	var v419 int32
	_ = v419
	var v421 int32
	_ = v421
	var v422 int32
	_ = v422
	var v426 int32
	_ = v426
	var v429 int32
	_ = v429
	v30 = int32(2)
	v31 = l0 << (uint(v30) % 32)
	v32 = l5 + v31
	v33 = F_memcpy(m, v32, l2, v31)
	mBase = m.M
	v34 = v33 + v31
	v38 = F_memcpy(m, v34, l2+l3<<(uint(v30)%32), v31)
	mBase = m.M
	if l1 < int32(1) {
	} else {
		v42 = l1 + int32(-1)
		if l0 < int32(3) {
			v370 = F_memcpy(m, l6, l2, v31)
			mBase = m.M
			if l1 == int32(1) {
			} else {
				v373 = int32(2)
				v374 = l3 << (uint(v373) % 32)
				v376 = l0 << (uint(v373) % 32)
				v378 = int32(-4)
				v380 = int32(3)
				__phi388 = l2
				__phi391 = l5
				__phi392 = l6
				__phi394 = v32
				__phi395 = v34
				__phi396 = v42
				v388 = __phi388
				v391 = __phi391
				v392 = __phi392
				v394 = __phi394
				v395 = __phi395
				v396 = __phi396
				for {
					v415 = v392 + v376
					v416 = v388 + v374
					if v396 != int32(1) {
						v421 = F_memcpy(m, v391, v388+l3<<(uint(v380)%32), v31)
						mBase = m.M
						v422 = *(*int32)(unsafe.Add(mBase, uint32(v416)))
						*(*int32)(unsafe.Add(mBase, uint32(v415))) = v422
						v426 = *(*int32)(unsafe.Add(mBase, uint32(v388+(v374+v376+v378))))
						*(*int32)(unsafe.Add(mBase, uint32(v392+(l0<<(uint(v380)%32)+v378)))) = v426
					} else {
						v419 = F_memcpy(m, v415, v416, v31)
						mBase = m.M
					}
					v429 = v396 + int32(-1)
					if v429 != 0 {
						__phi388 = v416
						__phi391 = v394
						__phi392 = v415
						__phi394 = v395
						__phi395 = v391
						__phi396 = v429
						v388 = __phi388
						v391 = __phi391
						v392 = __phi392
						v394 = __phi394
						v395 = __phi395
						v396 = __phi396
						continue
					} else {
						break
					}
					break
				}
			}
		} else {
			v45 = int32(0)
			v47 = int32(1)
			v48 = v47 << (uint(l4) % 32)
			v49 = v45 - v48
			v50 = int32(-1)
			v51 = v50 << (uint(l4) % 32)
			v55 = int32(base.Ui32(v51^v50) >> (uint(v47) % 32))
			v58 = int32(2)
			v63 = (l0 + v50) << (uint(v58) % 32)
			__phi66 = l2
			__phi69 = l5
			__phi70 = l6
			__phi72 = v32
			__phi73 = v34
			__phi75 = v45
			v66 = __phi66
			v69 = __phi69
			v70 = __phi70
			v72 = __phi72
			v73 = __phi73
			v75 = __phi75
			for {
				if v75 == int32(0) {
					v333 = F_memcpy(m, v70, v66, v31)
					mBase = m.M
				} else {
					if v75 == v42 {
						v333 = F_memcpy(m, v70, v66, v31)
						mBase = m.M
					} else {
						v99 = F_memcpy(m, v73, v66+l3<<(uint(int32(2))%32), v31)
						mBase = m.M
						v100 = *(*int32)(unsafe.Add(mBase, uint32(v66)))
						*(*int32)(unsafe.Add(mBase, uint32(v70))) = v100
						v104 = *(*int32)(unsafe.Add(mBase, uint32(v66+v63)))
						*(*int32)(unsafe.Add(mBase, uint32(v70+v63))) = v104
						v106 = int32(4)
						v111 = int32(0)
						v133 = l0 + int32(-2)
						for {
							v140 = v72 + v111
							v143 = *(*int32)(unsafe.Add(mBase, uint32(v140+int32(4))))
							v146 = int32(255)
							v147 = int32(base.Ui32(v143)>>(uint(int32(8))%32)) & v146
							v149 = v143 & v146
							v150 = *(*int32)(unsafe.Add(mBase, uint32(v140)))
							v153 = v149 - v150&v146
							if v48 <= v153 {
								v265 = int32(24)
								v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
								v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
								if base.Ui32(int32(255)) < base.Ui32(v271) {
									v277 = int32(-16777216)
								} else {
									v277 = v271 & v51 << (uint(v265) % 32)
								}
								v278 = int32(255)
								v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
								if base.Ui32(v278) < base.Ui32(v283) {
									v287 = v278
								} else {
									v287 = v283 & v51
								}
								v290 = int32(16)
								v292 = int32(255)
								v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
								v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
								if base.Ui32(v292) < base.Ui32(v298) {
									v304 = int32(16711680)
								} else {
									v304 = v298 & v51 << (uint(v290) % 32)
								}
								v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
								if base.Ui32(int32(255)) < base.Ui32(v311) {
									v317 = int32(_a_F_NearLossless_0)
								} else {
									v317 = v311 & v51 << (uint(int32(8)) % 32)
								}
								v319 = v277 | v287 | v304 | v317
							} else {
								if v153 <= v49 {
									v265 = int32(24)
									v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
									v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
									if base.Ui32(int32(255)) < base.Ui32(v271) {
										v277 = int32(-16777216)
									} else {
										v277 = v271 & v51 << (uint(v265) % 32)
									}
									v278 = int32(255)
									v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
									if base.Ui32(v278) < base.Ui32(v283) {
										v287 = v278
									} else {
										v287 = v283 & v51
									}
									v290 = int32(16)
									v292 = int32(255)
									v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
									v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
									if base.Ui32(v292) < base.Ui32(v298) {
										v304 = int32(16711680)
									} else {
										v304 = v298 & v51 << (uint(v290) % 32)
									}
									v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
									if base.Ui32(int32(255)) < base.Ui32(v311) {
										v317 = int32(_a_F_NearLossless_0)
									} else {
										v317 = v311 & v51 << (uint(int32(8)) % 32)
									}
									v319 = v277 | v287 | v304 | v317
								} else {
									v160 = v147 - int32(base.Ui32(v150)>>(uint(int32(8))%32))&int32(255)
									if v48 <= v160 {
										v265 = int32(24)
										v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
										v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
										if base.Ui32(int32(255)) < base.Ui32(v271) {
											v277 = int32(-16777216)
										} else {
											v277 = v271 & v51 << (uint(v265) % 32)
										}
										v278 = int32(255)
										v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
										if base.Ui32(v278) < base.Ui32(v283) {
											v287 = v278
										} else {
											v287 = v283 & v51
										}
										v290 = int32(16)
										v292 = int32(255)
										v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
										v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
										if base.Ui32(v292) < base.Ui32(v298) {
											v304 = int32(16711680)
										} else {
											v304 = v298 & v51 << (uint(v290) % 32)
										}
										v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
										if base.Ui32(int32(255)) < base.Ui32(v311) {
											v317 = int32(_a_F_NearLossless_0)
										} else {
											v317 = v311 & v51 << (uint(int32(8)) % 32)
										}
										v319 = v277 | v287 | v304 | v317
									} else {
										if v160 <= v49 {
											v265 = int32(24)
											v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
											v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
											if base.Ui32(int32(255)) < base.Ui32(v271) {
												v277 = int32(-16777216)
											} else {
												v277 = v271 & v51 << (uint(v265) % 32)
											}
											v278 = int32(255)
											v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
											if base.Ui32(v278) < base.Ui32(v283) {
												v287 = v278
											} else {
												v287 = v283 & v51
											}
											v290 = int32(16)
											v292 = int32(255)
											v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
											v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
											if base.Ui32(v292) < base.Ui32(v298) {
												v304 = int32(16711680)
											} else {
												v304 = v298 & v51 << (uint(v290) % 32)
											}
											v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
											if base.Ui32(int32(255)) < base.Ui32(v311) {
												v317 = int32(_a_F_NearLossless_0)
											} else {
												v317 = v311 & v51 << (uint(int32(8)) % 32)
											}
											v319 = v277 | v287 | v304 | v317
										} else {
											v163 = int32(16)
											v165 = int32(255)
											v166 = int32(base.Ui32(v143)>>(uint(v163)%32)) & v165
											v171 = v166 - int32(base.Ui32(v150)>>(uint(v163)%32))&v165
											if v48 <= v171 {
												v265 = int32(24)
												v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
												v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
												if base.Ui32(int32(255)) < base.Ui32(v271) {
													v277 = int32(-16777216)
												} else {
													v277 = v271 & v51 << (uint(v265) % 32)
												}
												v278 = int32(255)
												v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
												if base.Ui32(v278) < base.Ui32(v283) {
													v287 = v278
												} else {
													v287 = v283 & v51
												}
												v290 = int32(16)
												v292 = int32(255)
												v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
												v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
												if base.Ui32(v292) < base.Ui32(v298) {
													v304 = int32(16711680)
												} else {
													v304 = v298 & v51 << (uint(v290) % 32)
												}
												v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
												if base.Ui32(int32(255)) < base.Ui32(v311) {
													v317 = int32(_a_F_NearLossless_0)
												} else {
													v317 = v311 & v51 << (uint(int32(8)) % 32)
												}
												v319 = v277 | v287 | v304 | v317
											} else {
												if v171 <= v49 {
													v265 = int32(24)
													v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
													v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
													if base.Ui32(int32(255)) < base.Ui32(v271) {
														v277 = int32(-16777216)
													} else {
														v277 = v271 & v51 << (uint(v265) % 32)
													}
													v278 = int32(255)
													v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
													if base.Ui32(v278) < base.Ui32(v283) {
														v287 = v278
													} else {
														v287 = v283 & v51
													}
													v290 = int32(16)
													v292 = int32(255)
													v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
													v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
													if base.Ui32(v292) < base.Ui32(v298) {
														v304 = int32(16711680)
													} else {
														v304 = v298 & v51 << (uint(v290) % 32)
													}
													v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
													if base.Ui32(int32(255)) < base.Ui32(v311) {
														v317 = int32(_a_F_NearLossless_0)
													} else {
														v317 = v311 & v51 << (uint(int32(8)) % 32)
													}
													v319 = v277 | v287 | v304 | v317
												} else {
													v174 = int32(24)
													v175 = int32(base.Ui32(v143) >> (uint(v174) % 32))
													v178 = v175 - int32(base.Ui32(v150)>>(uint(v174)%32))
													if v48 <= v178 {
														v265 = int32(24)
														v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
														v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
														if base.Ui32(int32(255)) < base.Ui32(v271) {
															v277 = int32(-16777216)
														} else {
															v277 = v271 & v51 << (uint(v265) % 32)
														}
														v278 = int32(255)
														v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
														if base.Ui32(v278) < base.Ui32(v283) {
															v287 = v278
														} else {
															v287 = v283 & v51
														}
														v290 = int32(16)
														v292 = int32(255)
														v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
														v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
														if base.Ui32(v292) < base.Ui32(v298) {
															v304 = int32(16711680)
														} else {
															v304 = v298 & v51 << (uint(v290) % 32)
														}
														v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
														if base.Ui32(int32(255)) < base.Ui32(v311) {
															v317 = int32(_a_F_NearLossless_0)
														} else {
															v317 = v311 & v51 << (uint(int32(8)) % 32)
														}
														v319 = v277 | v287 | v304 | v317
													} else {
														if v178 <= v49 {
															v265 = int32(24)
															v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
															v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
															if base.Ui32(int32(255)) < base.Ui32(v271) {
																v277 = int32(-16777216)
															} else {
																v277 = v271 & v51 << (uint(v265) % 32)
															}
															v278 = int32(255)
															v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
															if base.Ui32(v278) < base.Ui32(v283) {
																v287 = v278
															} else {
																v287 = v283 & v51
															}
															v290 = int32(16)
															v292 = int32(255)
															v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
															v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
															if base.Ui32(v292) < base.Ui32(v298) {
																v304 = int32(16711680)
															} else {
																v304 = v298 & v51 << (uint(v290) % 32)
															}
															v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
															if base.Ui32(int32(255)) < base.Ui32(v311) {
																v317 = int32(_a_F_NearLossless_0)
															} else {
																v317 = v311 & v51 << (uint(int32(8)) % 32)
															}
															v319 = v277 | v287 | v304 | v317
														} else {
															v183 = *(*int32)(unsafe.Add(mBase, uint32(v140+int32(8))))
															v186 = v149 - v183&int32(255)
															if v48 <= v186 {
																v265 = int32(24)
																v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																if base.Ui32(int32(255)) < base.Ui32(v271) {
																	v277 = int32(-16777216)
																} else {
																	v277 = v271 & v51 << (uint(v265) % 32)
																}
																v278 = int32(255)
																v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																if base.Ui32(v278) < base.Ui32(v283) {
																	v287 = v278
																} else {
																	v287 = v283 & v51
																}
																v290 = int32(16)
																v292 = int32(255)
																v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																if base.Ui32(v292) < base.Ui32(v298) {
																	v304 = int32(16711680)
																} else {
																	v304 = v298 & v51 << (uint(v290) % 32)
																}
																v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																if base.Ui32(int32(255)) < base.Ui32(v311) {
																	v317 = int32(_a_F_NearLossless_0)
																} else {
																	v317 = v311 & v51 << (uint(int32(8)) % 32)
																}
																v319 = v277 | v287 | v304 | v317
															} else {
																if v186 <= v49 {
																	v265 = int32(24)
																	v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																	v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																	if base.Ui32(int32(255)) < base.Ui32(v271) {
																		v277 = int32(-16777216)
																	} else {
																		v277 = v271 & v51 << (uint(v265) % 32)
																	}
																	v278 = int32(255)
																	v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																	if base.Ui32(v278) < base.Ui32(v283) {
																		v287 = v278
																	} else {
																		v287 = v283 & v51
																	}
																	v290 = int32(16)
																	v292 = int32(255)
																	v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																	v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																	if base.Ui32(v292) < base.Ui32(v298) {
																		v304 = int32(16711680)
																	} else {
																		v304 = v298 & v51 << (uint(v290) % 32)
																	}
																	v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																	if base.Ui32(int32(255)) < base.Ui32(v311) {
																		v317 = int32(_a_F_NearLossless_0)
																	} else {
																		v317 = v311 & v51 << (uint(int32(8)) % 32)
																	}
																	v319 = v277 | v287 | v304 | v317
																} else {
																	v193 = v147 - int32(base.Ui32(v183)>>(uint(int32(8))%32))&int32(255)
																	if v48 <= v193 {
																		v265 = int32(24)
																		v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																		v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																		if base.Ui32(int32(255)) < base.Ui32(v271) {
																			v277 = int32(-16777216)
																		} else {
																			v277 = v271 & v51 << (uint(v265) % 32)
																		}
																		v278 = int32(255)
																		v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																		if base.Ui32(v278) < base.Ui32(v283) {
																			v287 = v278
																		} else {
																			v287 = v283 & v51
																		}
																		v290 = int32(16)
																		v292 = int32(255)
																		v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																		v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																		if base.Ui32(v292) < base.Ui32(v298) {
																			v304 = int32(16711680)
																		} else {
																			v304 = v298 & v51 << (uint(v290) % 32)
																		}
																		v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																		if base.Ui32(int32(255)) < base.Ui32(v311) {
																			v317 = int32(_a_F_NearLossless_0)
																		} else {
																			v317 = v311 & v51 << (uint(int32(8)) % 32)
																		}
																		v319 = v277 | v287 | v304 | v317
																	} else {
																		if v193 <= v49 {
																			v265 = int32(24)
																			v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																			v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																			if base.Ui32(int32(255)) < base.Ui32(v271) {
																				v277 = int32(-16777216)
																			} else {
																				v277 = v271 & v51 << (uint(v265) % 32)
																			}
																			v278 = int32(255)
																			v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																			if base.Ui32(v278) < base.Ui32(v283) {
																				v287 = v278
																			} else {
																				v287 = v283 & v51
																			}
																			v290 = int32(16)
																			v292 = int32(255)
																			v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																			v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																			if base.Ui32(v292) < base.Ui32(v298) {
																				v304 = int32(16711680)
																			} else {
																				v304 = v298 & v51 << (uint(v290) % 32)
																			}
																			v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																			if base.Ui32(int32(255)) < base.Ui32(v311) {
																				v317 = int32(_a_F_NearLossless_0)
																			} else {
																				v317 = v311 & v51 << (uint(int32(8)) % 32)
																			}
																			v319 = v277 | v287 | v304 | v317
																		} else {
																			v200 = v166 - int32(base.Ui32(v183)>>(uint(int32(16))%32))&int32(255)
																			if v48 <= v200 {
																				v265 = int32(24)
																				v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																				v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																				if base.Ui32(int32(255)) < base.Ui32(v271) {
																					v277 = int32(-16777216)
																				} else {
																					v277 = v271 & v51 << (uint(v265) % 32)
																				}
																				v278 = int32(255)
																				v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																				if base.Ui32(v278) < base.Ui32(v283) {
																					v287 = v278
																				} else {
																					v287 = v283 & v51
																				}
																				v290 = int32(16)
																				v292 = int32(255)
																				v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																				v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																				if base.Ui32(v292) < base.Ui32(v298) {
																					v304 = int32(16711680)
																				} else {
																					v304 = v298 & v51 << (uint(v290) % 32)
																				}
																				v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																				if base.Ui32(int32(255)) < base.Ui32(v311) {
																					v317 = int32(_a_F_NearLossless_0)
																				} else {
																					v317 = v311 & v51 << (uint(int32(8)) % 32)
																				}
																				v319 = v277 | v287 | v304 | v317
																			} else {
																				if v200 <= v49 {
																					v265 = int32(24)
																					v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																					v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																					if base.Ui32(int32(255)) < base.Ui32(v271) {
																						v277 = int32(-16777216)
																					} else {
																						v277 = v271 & v51 << (uint(v265) % 32)
																					}
																					v278 = int32(255)
																					v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																					if base.Ui32(v278) < base.Ui32(v283) {
																						v287 = v278
																					} else {
																						v287 = v283 & v51
																					}
																					v290 = int32(16)
																					v292 = int32(255)
																					v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																					v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																					if base.Ui32(v292) < base.Ui32(v298) {
																						v304 = int32(16711680)
																					} else {
																						v304 = v298 & v51 << (uint(v290) % 32)
																					}
																					v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																					if base.Ui32(int32(255)) < base.Ui32(v311) {
																						v317 = int32(_a_F_NearLossless_0)
																					} else {
																						v317 = v311 & v51 << (uint(int32(8)) % 32)
																					}
																					v319 = v277 | v287 | v304 | v317
																				} else {
																					v205 = v175 - int32(base.Ui32(v183)>>(uint(int32(24))%32))
																					if v48 <= v205 {
																						v265 = int32(24)
																						v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																						v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																						if base.Ui32(int32(255)) < base.Ui32(v271) {
																							v277 = int32(-16777216)
																						} else {
																							v277 = v271 & v51 << (uint(v265) % 32)
																						}
																						v278 = int32(255)
																						v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																						if base.Ui32(v278) < base.Ui32(v283) {
																							v287 = v278
																						} else {
																							v287 = v283 & v51
																						}
																						v290 = int32(16)
																						v292 = int32(255)
																						v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																						v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																						if base.Ui32(v292) < base.Ui32(v298) {
																							v304 = int32(16711680)
																						} else {
																							v304 = v298 & v51 << (uint(v290) % 32)
																						}
																						v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																						if base.Ui32(int32(255)) < base.Ui32(v311) {
																							v317 = int32(_a_F_NearLossless_0)
																						} else {
																							v317 = v311 & v51 << (uint(int32(8)) % 32)
																						}
																						v319 = v277 | v287 | v304 | v317
																					} else {
																						if v205 <= v49 {
																							v265 = int32(24)
																							v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																							v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																							if base.Ui32(int32(255)) < base.Ui32(v271) {
																								v277 = int32(-16777216)
																							} else {
																								v277 = v271 & v51 << (uint(v265) % 32)
																							}
																							v278 = int32(255)
																							v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																							if base.Ui32(v278) < base.Ui32(v283) {
																								v287 = v278
																							} else {
																								v287 = v283 & v51
																							}
																							v290 = int32(16)
																							v292 = int32(255)
																							v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																							v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																							if base.Ui32(v292) < base.Ui32(v298) {
																								v304 = int32(16711680)
																							} else {
																								v304 = v298 & v51 << (uint(v290) % 32)
																							}
																							v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																							if base.Ui32(int32(255)) < base.Ui32(v311) {
																								v317 = int32(_a_F_NearLossless_0)
																							} else {
																								v317 = v311 & v51 << (uint(int32(8)) % 32)
																							}
																							v319 = v277 | v287 | v304 | v317
																						} else {
																							v209 = *(*int32)(unsafe.Add(mBase, uint32(v69+v106+v111)))
																							v212 = v149 - v209&int32(255)
																							if v48 <= v212 {
																								v265 = int32(24)
																								v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																								v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																								if base.Ui32(int32(255)) < base.Ui32(v271) {
																									v277 = int32(-16777216)
																								} else {
																									v277 = v271 & v51 << (uint(v265) % 32)
																								}
																								v278 = int32(255)
																								v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																								if base.Ui32(v278) < base.Ui32(v283) {
																									v287 = v278
																								} else {
																									v287 = v283 & v51
																								}
																								v290 = int32(16)
																								v292 = int32(255)
																								v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																								v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																								if base.Ui32(v292) < base.Ui32(v298) {
																									v304 = int32(16711680)
																								} else {
																									v304 = v298 & v51 << (uint(v290) % 32)
																								}
																								v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																								if base.Ui32(int32(255)) < base.Ui32(v311) {
																									v317 = int32(_a_F_NearLossless_0)
																								} else {
																									v317 = v311 & v51 << (uint(int32(8)) % 32)
																								}
																								v319 = v277 | v287 | v304 | v317
																							} else {
																								if v212 <= v49 {
																									v265 = int32(24)
																									v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																									v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																									if base.Ui32(int32(255)) < base.Ui32(v271) {
																										v277 = int32(-16777216)
																									} else {
																										v277 = v271 & v51 << (uint(v265) % 32)
																									}
																									v278 = int32(255)
																									v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																									if base.Ui32(v278) < base.Ui32(v283) {
																										v287 = v278
																									} else {
																										v287 = v283 & v51
																									}
																									v290 = int32(16)
																									v292 = int32(255)
																									v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																									v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																									if base.Ui32(v292) < base.Ui32(v298) {
																										v304 = int32(16711680)
																									} else {
																										v304 = v298 & v51 << (uint(v290) % 32)
																									}
																									v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																									if base.Ui32(int32(255)) < base.Ui32(v311) {
																										v317 = int32(_a_F_NearLossless_0)
																									} else {
																										v317 = v311 & v51 << (uint(int32(8)) % 32)
																									}
																									v319 = v277 | v287 | v304 | v317
																								} else {
																									v219 = v147 - int32(base.Ui32(v209)>>(uint(int32(8))%32))&int32(255)
																									if v48 <= v219 {
																										v265 = int32(24)
																										v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																										v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																										if base.Ui32(int32(255)) < base.Ui32(v271) {
																											v277 = int32(-16777216)
																										} else {
																											v277 = v271 & v51 << (uint(v265) % 32)
																										}
																										v278 = int32(255)
																										v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																										if base.Ui32(v278) < base.Ui32(v283) {
																											v287 = v278
																										} else {
																											v287 = v283 & v51
																										}
																										v290 = int32(16)
																										v292 = int32(255)
																										v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																										v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																										if base.Ui32(v292) < base.Ui32(v298) {
																											v304 = int32(16711680)
																										} else {
																											v304 = v298 & v51 << (uint(v290) % 32)
																										}
																										v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																										if base.Ui32(int32(255)) < base.Ui32(v311) {
																											v317 = int32(_a_F_NearLossless_0)
																										} else {
																											v317 = v311 & v51 << (uint(int32(8)) % 32)
																										}
																										v319 = v277 | v287 | v304 | v317
																									} else {
																										if v219 <= v49 {
																											v265 = int32(24)
																											v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																											v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																											if base.Ui32(int32(255)) < base.Ui32(v271) {
																												v277 = int32(-16777216)
																											} else {
																												v277 = v271 & v51 << (uint(v265) % 32)
																											}
																											v278 = int32(255)
																											v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																											if base.Ui32(v278) < base.Ui32(v283) {
																												v287 = v278
																											} else {
																												v287 = v283 & v51
																											}
																											v290 = int32(16)
																											v292 = int32(255)
																											v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																											v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																											if base.Ui32(v292) < base.Ui32(v298) {
																												v304 = int32(16711680)
																											} else {
																												v304 = v298 & v51 << (uint(v290) % 32)
																											}
																											v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																											if base.Ui32(int32(255)) < base.Ui32(v311) {
																												v317 = int32(_a_F_NearLossless_0)
																											} else {
																												v317 = v311 & v51 << (uint(int32(8)) % 32)
																											}
																											v319 = v277 | v287 | v304 | v317
																										} else {
																											v226 = v166 - int32(base.Ui32(v209)>>(uint(int32(16))%32))&int32(255)
																											if v48 <= v226 {
																												v265 = int32(24)
																												v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																												v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																												if base.Ui32(int32(255)) < base.Ui32(v271) {
																													v277 = int32(-16777216)
																												} else {
																													v277 = v271 & v51 << (uint(v265) % 32)
																												}
																												v278 = int32(255)
																												v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																												if base.Ui32(v278) < base.Ui32(v283) {
																													v287 = v278
																												} else {
																													v287 = v283 & v51
																												}
																												v290 = int32(16)
																												v292 = int32(255)
																												v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																												v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																												if base.Ui32(v292) < base.Ui32(v298) {
																													v304 = int32(16711680)
																												} else {
																													v304 = v298 & v51 << (uint(v290) % 32)
																												}
																												v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																												if base.Ui32(int32(255)) < base.Ui32(v311) {
																													v317 = int32(_a_F_NearLossless_0)
																												} else {
																													v317 = v311 & v51 << (uint(int32(8)) % 32)
																												}
																												v319 = v277 | v287 | v304 | v317
																											} else {
																												if v226 <= v49 {
																													v265 = int32(24)
																													v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																													v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																													if base.Ui32(int32(255)) < base.Ui32(v271) {
																														v277 = int32(-16777216)
																													} else {
																														v277 = v271 & v51 << (uint(v265) % 32)
																													}
																													v278 = int32(255)
																													v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																													if base.Ui32(v278) < base.Ui32(v283) {
																														v287 = v278
																													} else {
																														v287 = v283 & v51
																													}
																													v290 = int32(16)
																													v292 = int32(255)
																													v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																													v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																													if base.Ui32(v292) < base.Ui32(v298) {
																														v304 = int32(16711680)
																													} else {
																														v304 = v298 & v51 << (uint(v290) % 32)
																													}
																													v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																													if base.Ui32(int32(255)) < base.Ui32(v311) {
																														v317 = int32(_a_F_NearLossless_0)
																													} else {
																														v317 = v311 & v51 << (uint(int32(8)) % 32)
																													}
																													v319 = v277 | v287 | v304 | v317
																												} else {
																													v231 = v175 - int32(base.Ui32(v209)>>(uint(int32(24))%32))
																													if v48 <= v231 {
																														v265 = int32(24)
																														v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																														v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																														if base.Ui32(int32(255)) < base.Ui32(v271) {
																															v277 = int32(-16777216)
																														} else {
																															v277 = v271 & v51 << (uint(v265) % 32)
																														}
																														v278 = int32(255)
																														v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																														if base.Ui32(v278) < base.Ui32(v283) {
																															v287 = v278
																														} else {
																															v287 = v283 & v51
																														}
																														v290 = int32(16)
																														v292 = int32(255)
																														v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																														v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																														if base.Ui32(v292) < base.Ui32(v298) {
																															v304 = int32(16711680)
																														} else {
																															v304 = v298 & v51 << (uint(v290) % 32)
																														}
																														v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																														if base.Ui32(int32(255)) < base.Ui32(v311) {
																															v317 = int32(_a_F_NearLossless_0)
																														} else {
																															v317 = v311 & v51 << (uint(int32(8)) % 32)
																														}
																														v319 = v277 | v287 | v304 | v317
																													} else {
																														if v231 <= v49 {
																															v265 = int32(24)
																															v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																															v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																															if base.Ui32(int32(255)) < base.Ui32(v271) {
																																v277 = int32(-16777216)
																															} else {
																																v277 = v271 & v51 << (uint(v265) % 32)
																															}
																															v278 = int32(255)
																															v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																															if base.Ui32(v278) < base.Ui32(v283) {
																																v287 = v278
																															} else {
																																v287 = v283 & v51
																															}
																															v290 = int32(16)
																															v292 = int32(255)
																															v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																															v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																															if base.Ui32(v292) < base.Ui32(v298) {
																																v304 = int32(16711680)
																															} else {
																																v304 = v298 & v51 << (uint(v290) % 32)
																															}
																															v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																															if base.Ui32(int32(255)) < base.Ui32(v311) {
																																v317 = int32(_a_F_NearLossless_0)
																															} else {
																																v317 = v311 & v51 << (uint(int32(8)) % 32)
																															}
																															v319 = v277 | v287 | v304 | v317
																														} else {
																															v235 = *(*int32)(unsafe.Add(mBase, uint32(v99+v106+v111)))
																															v238 = v149 - v235&int32(255)
																															if v48 <= v238 {
																																v265 = int32(24)
																																v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																																v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																																if base.Ui32(int32(255)) < base.Ui32(v271) {
																																	v277 = int32(-16777216)
																																} else {
																																	v277 = v271 & v51 << (uint(v265) % 32)
																																}
																																v278 = int32(255)
																																v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																																if base.Ui32(v278) < base.Ui32(v283) {
																																	v287 = v278
																																} else {
																																	v287 = v283 & v51
																																}
																																v290 = int32(16)
																																v292 = int32(255)
																																v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																																v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																																if base.Ui32(v292) < base.Ui32(v298) {
																																	v304 = int32(16711680)
																																} else {
																																	v304 = v298 & v51 << (uint(v290) % 32)
																																}
																																v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																																if base.Ui32(int32(255)) < base.Ui32(v311) {
																																	v317 = int32(_a_F_NearLossless_0)
																																} else {
																																	v317 = v311 & v51 << (uint(int32(8)) % 32)
																																}
																																v319 = v277 | v287 | v304 | v317
																															} else {
																																if v238 <= v49 {
																																	v265 = int32(24)
																																	v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																																	v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																																	if base.Ui32(int32(255)) < base.Ui32(v271) {
																																		v277 = int32(-16777216)
																																	} else {
																																		v277 = v271 & v51 << (uint(v265) % 32)
																																	}
																																	v278 = int32(255)
																																	v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																																	if base.Ui32(v278) < base.Ui32(v283) {
																																		v287 = v278
																																	} else {
																																		v287 = v283 & v51
																																	}
																																	v290 = int32(16)
																																	v292 = int32(255)
																																	v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																																	v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																																	if base.Ui32(v292) < base.Ui32(v298) {
																																		v304 = int32(16711680)
																																	} else {
																																		v304 = v298 & v51 << (uint(v290) % 32)
																																	}
																																	v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																																	if base.Ui32(int32(255)) < base.Ui32(v311) {
																																		v317 = int32(_a_F_NearLossless_0)
																																	} else {
																																		v317 = v311 & v51 << (uint(int32(8)) % 32)
																																	}
																																	v319 = v277 | v287 | v304 | v317
																																} else {
																																	v245 = v147 - int32(base.Ui32(v235)>>(uint(int32(8))%32))&int32(255)
																																	if v48 <= v245 {
																																		v265 = int32(24)
																																		v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																																		v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																																		if base.Ui32(int32(255)) < base.Ui32(v271) {
																																			v277 = int32(-16777216)
																																		} else {
																																			v277 = v271 & v51 << (uint(v265) % 32)
																																		}
																																		v278 = int32(255)
																																		v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																																		if base.Ui32(v278) < base.Ui32(v283) {
																																			v287 = v278
																																		} else {
																																			v287 = v283 & v51
																																		}
																																		v290 = int32(16)
																																		v292 = int32(255)
																																		v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																																		v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																																		if base.Ui32(v292) < base.Ui32(v298) {
																																			v304 = int32(16711680)
																																		} else {
																																			v304 = v298 & v51 << (uint(v290) % 32)
																																		}
																																		v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																																		if base.Ui32(int32(255)) < base.Ui32(v311) {
																																			v317 = int32(_a_F_NearLossless_0)
																																		} else {
																																			v317 = v311 & v51 << (uint(int32(8)) % 32)
																																		}
																																		v319 = v277 | v287 | v304 | v317
																																	} else {
																																		if v245 <= v49 {
																																			v265 = int32(24)
																																			v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																																			v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																																			if base.Ui32(int32(255)) < base.Ui32(v271) {
																																				v277 = int32(-16777216)
																																			} else {
																																				v277 = v271 & v51 << (uint(v265) % 32)
																																			}
																																			v278 = int32(255)
																																			v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																																			if base.Ui32(v278) < base.Ui32(v283) {
																																				v287 = v278
																																			} else {
																																				v287 = v283 & v51
																																			}
																																			v290 = int32(16)
																																			v292 = int32(255)
																																			v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																																			v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																																			if base.Ui32(v292) < base.Ui32(v298) {
																																				v304 = int32(16711680)
																																			} else {
																																				v304 = v298 & v51 << (uint(v290) % 32)
																																			}
																																			v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																																			if base.Ui32(int32(255)) < base.Ui32(v311) {
																																				v317 = int32(_a_F_NearLossless_0)
																																			} else {
																																				v317 = v311 & v51 << (uint(int32(8)) % 32)
																																			}
																																			v319 = v277 | v287 | v304 | v317
																																		} else {
																																			v252 = v166 - int32(base.Ui32(v235)>>(uint(int32(16))%32))&int32(255)
																																			if v48 <= v252 {
																																				v265 = int32(24)
																																				v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																																				v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																																				if base.Ui32(int32(255)) < base.Ui32(v271) {
																																					v277 = int32(-16777216)
																																				} else {
																																					v277 = v271 & v51 << (uint(v265) % 32)
																																				}
																																				v278 = int32(255)
																																				v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																																				if base.Ui32(v278) < base.Ui32(v283) {
																																					v287 = v278
																																				} else {
																																					v287 = v283 & v51
																																				}
																																				v290 = int32(16)
																																				v292 = int32(255)
																																				v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																																				v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																																				if base.Ui32(v292) < base.Ui32(v298) {
																																					v304 = int32(16711680)
																																				} else {
																																					v304 = v298 & v51 << (uint(v290) % 32)
																																				}
																																				v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																																				if base.Ui32(int32(255)) < base.Ui32(v311) {
																																					v317 = int32(_a_F_NearLossless_0)
																																				} else {
																																					v317 = v311 & v51 << (uint(int32(8)) % 32)
																																				}
																																				v319 = v277 | v287 | v304 | v317
																																			} else {
																																				if v252 <= v49 {
																																					v265 = int32(24)
																																					v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																																					v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																																					if base.Ui32(int32(255)) < base.Ui32(v271) {
																																						v277 = int32(-16777216)
																																					} else {
																																						v277 = v271 & v51 << (uint(v265) % 32)
																																					}
																																					v278 = int32(255)
																																					v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																																					if base.Ui32(v278) < base.Ui32(v283) {
																																						v287 = v278
																																					} else {
																																						v287 = v283 & v51
																																					}
																																					v290 = int32(16)
																																					v292 = int32(255)
																																					v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																																					v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																																					if base.Ui32(v292) < base.Ui32(v298) {
																																						v304 = int32(16711680)
																																					} else {
																																						v304 = v298 & v51 << (uint(v290) % 32)
																																					}
																																					v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																																					if base.Ui32(int32(255)) < base.Ui32(v311) {
																																						v317 = int32(_a_F_NearLossless_0)
																																					} else {
																																						v317 = v311 & v51 << (uint(int32(8)) % 32)
																																					}
																																					v319 = v277 | v287 | v304 | v317
																																				} else {
																																					v257 = v175 - int32(base.Ui32(v235)>>(uint(int32(24))%32))
																																					if v48 <= v257 {
																																						v265 = int32(24)
																																						v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																																						v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																																						if base.Ui32(int32(255)) < base.Ui32(v271) {
																																							v277 = int32(-16777216)
																																						} else {
																																							v277 = v271 & v51 << (uint(v265) % 32)
																																						}
																																						v278 = int32(255)
																																						v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																																						if base.Ui32(v278) < base.Ui32(v283) {
																																							v287 = v278
																																						} else {
																																							v287 = v283 & v51
																																						}
																																						v290 = int32(16)
																																						v292 = int32(255)
																																						v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																																						v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																																						if base.Ui32(v292) < base.Ui32(v298) {
																																							v304 = int32(16711680)
																																						} else {
																																							v304 = v298 & v51 << (uint(v290) % 32)
																																						}
																																						v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																																						if base.Ui32(int32(255)) < base.Ui32(v311) {
																																							v317 = int32(_a_F_NearLossless_0)
																																						} else {
																																							v317 = v311 & v51 << (uint(int32(8)) % 32)
																																						}
																																						v319 = v277 | v287 | v304 | v317
																																					} else {
																																						if v49 < v257 {
																																							v319 = v143
																																						} else {
																																							v265 = int32(24)
																																							v266 = int32(base.Ui32(v143) >> (uint(v265) % 32))
																																							v271 = v266 + v55 + int32(base.Ui32(v266)>>(uint(l4)%32))&int32(1)
																																							if base.Ui32(int32(255)) < base.Ui32(v271) {
																																								v277 = int32(-16777216)
																																							} else {
																																								v277 = v271 & v51 << (uint(v265) % 32)
																																							}
																																							v278 = int32(255)
																																							v283 = v149 + v55 + int32(base.Ui32(v149)>>(uint(l4)%32))&int32(1)
																																							if base.Ui32(v278) < base.Ui32(v283) {
																																								v287 = v278
																																							} else {
																																								v287 = v283 & v51
																																							}
																																							v290 = int32(16)
																																							v292 = int32(255)
																																							v293 = int32(base.Ui32(v143)>>(uint(v290)%32)) & v292
																																							v298 = v293 + v55 + int32(base.Ui32(v293)>>(uint(l4)%32))&int32(1)
																																							if base.Ui32(v292) < base.Ui32(v298) {
																																								v304 = int32(16711680)
																																							} else {
																																								v304 = v298 & v51 << (uint(v290) % 32)
																																							}
																																							v311 = v147 + v55 + int32(base.Ui32(v147)>>(uint(l4)%32))&int32(1)
																																							if base.Ui32(int32(255)) < base.Ui32(v311) {
																																								v317 = int32(_a_F_NearLossless_0)
																																							} else {
																																								v317 = v311 & v51 << (uint(int32(8)) % 32)
																																							}
																																							v319 = v277 | v287 | v304 | v317
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
								}
							}
							v326 = int32(4)
							*(*int32)(unsafe.Add(mBase, uint32(v70+v111+v326))) = v319
							v332 = v133 + int32(-1)
							if v332 != 0 {
								v111 = v111 + v326
								v133 = v332
								continue
							} else {
								break
							}
							break
						}
					}
				}
				v368 = v75 + int32(1)
				if v368 != l1 {
					__phi66 = v66 + l3<<(uint(int32(2))%32)
					__phi69 = v72
					__phi70 = v70 + l0<<(uint(v58)%32)
					__phi72 = v73
					__phi73 = v69
					__phi75 = v368
					v66 = __phi66
					v69 = __phi69
					v70 = __phi70
					v72 = __phi72
					v73 = __phi73
					v75 = __phi75
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
func F_NoneUnfilter_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	var v6 int32
	_ = v6
	if l2 == l1 {
	} else {
		v6 = F_memcpy(m, l2, l1, l3)
	}
	return
}
