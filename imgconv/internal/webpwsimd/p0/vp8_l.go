//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_VP8LBackwardRefsClear(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v7 int32
	_ = v7
	var v9 int32
	_ = v9
	var v11 int32
	_ = v11
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	v4 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	if v4 == int32(0) {
	} else {
		v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
		*(*int32)(unsafe.Add(mBase, uint32(v4))) = v7
	}
	v9 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v9
	v11 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v11
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = l0 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v9
	if v11 == v9 {
	} else {
		v22 = v11
		for {
			v23 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
			F_free(m, v22)
			mBase = m.M
			*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v23
			if v23 != 0 {
				v22 = v23
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_VP8LBackwardRefsCursorAdd(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v12 int32
	_ = v12
	var v13 int64
	_ = v13
	var v14 int32
	_ = v14
	var v18 int32
	_ = v18
	var v25 int64
	_ = v25
	var v26 int32
	_ = v26
	var v27 int64
	_ = v27
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v41 int32
	_ = v41
	var v48 int32
	_ = v48
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v65 int32
	_ = v65
	var v69 int64
	_ = v69
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	if v5 == int32(0) {
		v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
		if v12 != 0 {
			v48 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
			*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v48
			v50 = v12
			v51 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
			*(*int32)(unsafe.Add(mBase, uint32(v51))) = v50
			v53 = int32(0)
			*(*int32)(unsafe.Add(mBase, uint32(v50)+8)) = v53
			*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v50
			*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v50
			*(*int32)(unsafe.Add(mBase, uint32(v50))) = v53
			v60 = v50
			v61 = v53
			*(*int32)(unsafe.Add(mBase, uint32(v60)+8)) = v61 + int32(1)
			v65 = *(*int32)(unsafe.Add(mBase, uint32(v60)+4))
			v69 = *(*int64)(unsafe.Add(mBase, uint32(l1)))
			*(*int64)(unsafe.Add(mBase, uint32(v65+v61<<(uint(int32(3))%32)))) = v69
			return
		} else {
			v13 = int64(1)
			v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
			v18 = v14<<(uint(int32(3))%32) + int32(12)
			v25 = base.I64_div_u_s(int64(2147418112), v13)
			v26 = int32(0)
			v27 = base.I64_extend_i32_u(v18)
			if base.Ui64(int64(4294967295)) < base.Ui64(v27*v13) {
				v39 = v26
			} else {
				if base.Ui64(v25) < base.Ui64(v27) {
					v39 = v26
				} else {
					v37 = F_malloc(m, base.I32_wrap_i64(v13)*v18)
					mBase = m.M
					v39 = v37
				}
			}
			if v39 != 0 {
				*(*int32)(unsafe.Add(mBase, uint32(v39)+4)) = v39 + int32(12)
				v50 = v39
				v51 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
				*(*int32)(unsafe.Add(mBase, uint32(v51))) = v50
				v53 = int32(0)
				*(*int32)(unsafe.Add(mBase, uint32(v50)+8)) = v53
				*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v50
				*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v50
				*(*int32)(unsafe.Add(mBase, uint32(v50))) = v53
				v60 = v50
				v61 = v53
				*(*int32)(unsafe.Add(mBase, uint32(v60)+8)) = v61 + int32(1)
				v65 = *(*int32)(unsafe.Add(mBase, uint32(v60)+4))
				v69 = *(*int64)(unsafe.Add(mBase, uint32(l1)))
				*(*int64)(unsafe.Add(mBase, uint32(v65+v61<<(uint(int32(3))%32)))) = v69
				return
			} else {
				v41 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
				*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v41 | int32(1)
				return
			}
		}
	} else {
		v8 = *(*int32)(unsafe.Add(mBase, uint32(v5)+8))
		v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		if v8 != v9 {
			v60 = v5
			v61 = v8
			*(*int32)(unsafe.Add(mBase, uint32(v60)+8)) = v61 + int32(1)
			v65 = *(*int32)(unsafe.Add(mBase, uint32(v60)+4))
			v69 = *(*int64)(unsafe.Add(mBase, uint32(l1)))
			*(*int64)(unsafe.Add(mBase, uint32(v65+v61<<(uint(int32(3))%32)))) = v69
			return
		} else {
			v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
			if v12 != 0 {
				v48 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
				*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v48
				v50 = v12
				v51 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
				*(*int32)(unsafe.Add(mBase, uint32(v51))) = v50
				v53 = int32(0)
				*(*int32)(unsafe.Add(mBase, uint32(v50)+8)) = v53
				*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v50
				*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v50
				*(*int32)(unsafe.Add(mBase, uint32(v50))) = v53
				v60 = v50
				v61 = v53
				*(*int32)(unsafe.Add(mBase, uint32(v60)+8)) = v61 + int32(1)
				v65 = *(*int32)(unsafe.Add(mBase, uint32(v60)+4))
				v69 = *(*int64)(unsafe.Add(mBase, uint32(l1)))
				*(*int64)(unsafe.Add(mBase, uint32(v65+v61<<(uint(int32(3))%32)))) = v69
				return
			} else {
				v13 = int64(1)
				v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				v18 = v14<<(uint(int32(3))%32) + int32(12)
				v25 = base.I64_div_u_s(int64(2147418112), v13)
				v26 = int32(0)
				v27 = base.I64_extend_i32_u(v18)
				if base.Ui64(int64(4294967295)) < base.Ui64(v27*v13) {
					v39 = v26
				} else {
					if base.Ui64(v25) < base.Ui64(v27) {
						v39 = v26
					} else {
						v37 = F_malloc(m, base.I32_wrap_i64(v13)*v18)
						mBase = m.M
						v39 = v37
					}
				}
				if v39 != 0 {
					*(*int32)(unsafe.Add(mBase, uint32(v39)+4)) = v39 + int32(12)
					v50 = v39
					v51 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					*(*int32)(unsafe.Add(mBase, uint32(v51))) = v50
					v53 = int32(0)
					*(*int32)(unsafe.Add(mBase, uint32(v50)+8)) = v53
					*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v50
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v50
					*(*int32)(unsafe.Add(mBase, uint32(v50))) = v53
					v60 = v50
					v61 = v53
					*(*int32)(unsafe.Add(mBase, uint32(v60)+8)) = v61 + int32(1)
					v65 = *(*int32)(unsafe.Add(mBase, uint32(v60)+4))
					v69 = *(*int64)(unsafe.Add(mBase, uint32(l1)))
					*(*int64)(unsafe.Add(mBase, uint32(v65+v61<<(uint(int32(3))%32)))) = v69
					return
				} else {
					v41 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v41 | int32(1)
					return
				}
			}
		}
	}
}
func F_VP8LBitWriterClone(m *base.Module, l0 int32, l1 int32) int32 {
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
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v20 int64
	_ = v20
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v32 int64
	_ = v32
	var v36 int32
	_ = v36
	var v38 int32
	_ = v38
	var v42 int32
	_ = v42
	var v49 int64
	_ = v49
	var v50 int32
	_ = v50
	var v51 int64
	_ = v51
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v83 int64
	_ = v83
	var v85 int32
	_ = v85
	var v87 int32
	_ = v87
	v12 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	v13 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v14 = v12 - v13
	v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v18 = v16 - v17
	v20 = base.I64_extend_i32_u(v14) + base.I64_extend_i32_u(v18)
	if base.Ui64(v20) < base.Ui64(int64(4294967296)) {
		v27 = base.I32_wrap_i64(v20)
		v28 = *(*int32)(unsafe.Add(mBase, uint32(l1)+16))
		v29 = v28 - v13
		if v28 == v13 {
			v32 = int64(1)
			v36 = int32(base.Ui32(v29*int32(3)) >> (uint(int32(1)) % 32))
			if base.Ui32(v27) < base.Ui32(v36) {
				v38 = v36
			} else {
				v38 = v27
			}
			v42 = v38&int32(-1024) + int32(1024)
			v49 = base.I64_div_u_s(int64(2147418112), v32)
			v50 = int32(0)
			v51 = base.I64_extend_i32_u(v42)
			if base.Ui64(int64(4294967295)) < base.Ui64(v51*v32) {
				v63 = v50
			} else {
				if base.Ui64(v49) < base.Ui64(v51) {
					v63 = v50
				} else {
					v61 = F_malloc(m, base.I32_wrap_i64(v32)*v42)
					mBase = m.M
					v63 = v61
				}
			}
			if v63 != 0 {
				if v12 == v13 {
				} else {
					v70 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
					v71 = F_memcpy(m, v63, v70, v14)
					mBase = m.M
				}
				v72 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
				F_free(m, v72)
				mBase = m.M
				*(*int32)(unsafe.Add(mBase, uint32(l1)+8)) = v63
				*(*int32)(unsafe.Add(mBase, uint32(l1)+16)) = v63 + v42
				*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v63 + v14
				v79 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				v80 = v79
				v81 = v63
				v82 = F_memcpy(m, v81, v80, v18)
				mBase = m.M
				v83 = *(*int64)(unsafe.Add(mBase, uint32(l0)))
				*(*int64)(unsafe.Add(mBase, uint32(l1))) = v83
				v85 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
				*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = v85
				v87 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
				*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v87 + v18
				return int32(1)
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = int32(1)
				return int32(0)
			}
		} else {
			if base.Ui32(v29) < base.Ui32(v27) {
				v32 = int64(1)
				v36 = int32(base.Ui32(v29*int32(3)) >> (uint(int32(1)) % 32))
				if base.Ui32(v27) < base.Ui32(v36) {
					v38 = v36
				} else {
					v38 = v27
				}
				v42 = v38&int32(-1024) + int32(1024)
				v49 = base.I64_div_u_s(int64(2147418112), v32)
				v50 = int32(0)
				v51 = base.I64_extend_i32_u(v42)
				if base.Ui64(int64(4294967295)) < base.Ui64(v51*v32) {
					v63 = v50
				} else {
					if base.Ui64(v49) < base.Ui64(v51) {
						v63 = v50
					} else {
						v61 = F_malloc(m, base.I32_wrap_i64(v32)*v42)
						mBase = m.M
						v63 = v61
					}
				}
				if v63 != 0 {
					if v12 == v13 {
					} else {
						v70 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
						v71 = F_memcpy(m, v63, v70, v14)
						mBase = m.M
					}
					v72 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
					F_free(m, v72)
					mBase = m.M
					*(*int32)(unsafe.Add(mBase, uint32(l1)+8)) = v63
					*(*int32)(unsafe.Add(mBase, uint32(l1)+16)) = v63 + v42
					*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v63 + v14
					v79 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					v80 = v79
					v81 = v63
					v82 = F_memcpy(m, v81, v80, v18)
					mBase = m.M
					v83 = *(*int64)(unsafe.Add(mBase, uint32(l0)))
					*(*int64)(unsafe.Add(mBase, uint32(l1))) = v83
					v85 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
					*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = v85
					v87 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
					*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v87 + v18
					return int32(1)
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = int32(1)
					return int32(0)
				}
			} else {
				v80 = v17
				v81 = v13
				v82 = F_memcpy(m, v81, v80, v18)
				mBase = m.M
				v83 = *(*int64)(unsafe.Add(mBase, uint32(l0)))
				*(*int64)(unsafe.Add(mBase, uint32(l1))) = v83
				v85 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
				*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = v85
				v87 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
				*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v87 + v18
				return int32(1)
			}
		}
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = int32(1)
		return int32(0)
	}
}
func F_VP8LBitWriterFinish(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	var v20 int64
	_ = v20
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v32 int64
	_ = v32
	var v36 int32
	_ = v36
	var v38 int32
	_ = v38
	var v42 int32
	_ = v42
	var v49 int64
	_ = v49
	var v50 int32
	_ = v50
	var v51 int64
	_ = v51
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v67 int32
	_ = v67
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v84 int32
	_ = v84
	var v87 int32
	_ = v87
	var v94 int32
	_ = v94
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v101 int32
	_ = v101
	var v103 int32
	_ = v103
	var v120 int32
	_ = v120
	v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v11 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v12 = v10 - v11
	v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v20 = base.I64_extend_i32_u(v12) + base.I64_extend_i32_u((v14+int32(7))>>(uint(int32(3))%32))
	if base.Ui64(v20) < base.Ui64(int64(4294967296)) {
		v27 = base.I32_wrap_i64(v20)
		v28 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
		v29 = v28 - v11
		if v28 == v11 {
			v32 = int64(1)
			v36 = int32(base.Ui32(v29*int32(3)) >> (uint(int32(1)) % 32))
			if base.Ui32(v27) < base.Ui32(v36) {
				v38 = v36
			} else {
				v38 = v27
			}
			v42 = v38&int32(-1024) + int32(1024)
			v49 = base.I64_div_u_s(int64(2147418112), v32)
			v50 = int32(0)
			v51 = base.I64_extend_i32_u(v42)
			if base.Ui64(int64(4294967295)) < base.Ui64(v51*v32) {
				v63 = v50
			} else {
				if base.Ui64(v49) < base.Ui64(v51) {
					v63 = v50
				} else {
					v61 = F_malloc(m, base.I32_wrap_i64(v32)*v42)
					mBase = m.M
					v63 = v61
				}
			}
			if v63 != 0 {
				if v10 == v11 {
				} else {
					v70 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					v71 = F_memcpy(m, v63, v70, v12)
					mBase = m.M
				}
				v72 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				F_free(m, v72)
				mBase = m.M
				*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v63 + v42
				*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v63 + v12
				*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v63
				v79 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
				v80 = v79
				if v80 < int32(1) {
				} else {
					v84 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
					v87 = v84
					for {
						v94 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
						*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v94 + int32(1)
						*(*uint8)(unsafe.Add(mBase, uint32(v94))) = uint8(v87)
						v99 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v100 = int32(8)
						v101 = int32(base.Ui32(v99) >> (uint(v100) % 32))
						*(*int32)(unsafe.Add(mBase, uint32(l0))) = v101
						v103 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
						*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v103 + int32(-8)
						if v100 < v103 {
							v87 = v101
							continue
						} else {
							break
						}
						break
					}
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = int32(0)
				v120 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				return v120
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
				v67 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				return v67
			}
		} else {
			if base.Ui32(v27) <= base.Ui32(v29) {
				v80 = v14
				if v80 < int32(1) {
				} else {
					v84 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
					v87 = v84
					for {
						v94 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
						*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v94 + int32(1)
						*(*uint8)(unsafe.Add(mBase, uint32(v94))) = uint8(v87)
						v99 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v100 = int32(8)
						v101 = int32(base.Ui32(v99) >> (uint(v100) % 32))
						*(*int32)(unsafe.Add(mBase, uint32(l0))) = v101
						v103 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
						*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v103 + int32(-8)
						if v100 < v103 {
							v87 = v101
							continue
						} else {
							break
						}
						break
					}
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = int32(0)
				v120 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				return v120
			} else {
				v32 = int64(1)
				v36 = int32(base.Ui32(v29*int32(3)) >> (uint(int32(1)) % 32))
				if base.Ui32(v27) < base.Ui32(v36) {
					v38 = v36
				} else {
					v38 = v27
				}
				v42 = v38&int32(-1024) + int32(1024)
				v49 = base.I64_div_u_s(int64(2147418112), v32)
				v50 = int32(0)
				v51 = base.I64_extend_i32_u(v42)
				if base.Ui64(int64(4294967295)) < base.Ui64(v51*v32) {
					v63 = v50
				} else {
					if base.Ui64(v49) < base.Ui64(v51) {
						v63 = v50
					} else {
						v61 = F_malloc(m, base.I32_wrap_i64(v32)*v42)
						mBase = m.M
						v63 = v61
					}
				}
				if v63 != 0 {
					if v10 == v11 {
					} else {
						v70 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
						v71 = F_memcpy(m, v63, v70, v12)
						mBase = m.M
					}
					v72 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					F_free(m, v72)
					mBase = m.M
					*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v63 + v42
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v63 + v12
					*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v63
					v79 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v80 = v79
					if v80 < int32(1) {
					} else {
						v84 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v87 = v84
						for {
							v94 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v94 + int32(1)
							*(*uint8)(unsafe.Add(mBase, uint32(v94))) = uint8(v87)
							v99 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
							v100 = int32(8)
							v101 = int32(base.Ui32(v99) >> (uint(v100) % 32))
							*(*int32)(unsafe.Add(mBase, uint32(l0))) = v101
							v103 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v103 + int32(-8)
							if v100 < v103 {
								v87 = v101
								continue
							} else {
								break
							}
							break
						}
					}
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = int32(0)
					v120 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					return v120
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
					v67 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					return v67
				}
			}
		}
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
		v25 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
		return v25
	}
}
func F_VP8LBitWriterReset(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v3 int64
	_ = v3
	var v5 int32
	_ = v5
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	v3 = *(*int64)(unsafe.Add(mBase, uint32(l0)))
	*(*int64)(unsafe.Add(mBase, uint32(l1))) = v3
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = v5
	v7 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v8 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v7 + (v8 - v9)
	return
}
func F_VP8LBitsEntropy(m *base.Module, l0 int32, l1 int32) int64 {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v13 int32
	_ = v13
	var v25 int64
	_ = v25
	var v28 int64
	_ = v28
	var v29 int64
	_ = v29
	var v32 int64
	_ = v32
	var v35 int64
	_ = v35
	var v37 int64
	_ = v37
	var v38 int64
	_ = v38
	var v41 int32
	_ = v41
	var v44 int32
	_ = v44
	var v53 int64
	_ = v53
	var v56 int64
	_ = v56
	var v59 int64
	_ = v59
	var v61 int64
	_ = v61
	var v62 int64
	_ = v62
	v6 = m.G0
	v8 = v6 - int32(32)
	m.G0 = v8
	F_VP8LBitsEntropyUnrefined(m, l0, l1, v8+int32(8))
	mBase = m.M
	v13 = *(*int32)(unsafe.Add(mBase, uint32(v8)+20))
	if v13 <= int32(4) {
		if v13 < int32(2) {
			v62 = int64(0)
		} else {
			switch v13 + int32(-2) {
			case 0:
				v25 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v8)+16)))
				v28 = *(*int64)(unsafe.Add(mBase, uint32(v8)+8))
				v29 = v25*int64(830472192) + v28
				if v29 < int64(0) {
					v32 = int64(-50)
				} else {
					v32 = int64(50)
				}
				v35 = base.I64_div_s(v32+v29, int64(100))
				v62 = v35
			case 1:
				v37 = int64(950)
				v38 = *(*int64)(unsafe.Add(mBase, uint32(v8)+8))
				v41 = *(*int32)(unsafe.Add(mBase, uint32(v8)+16))
				v44 = *(*int32)(unsafe.Add(mBase, uint32(v8)+24))
				v53 = v37*base.I64_extend_i32_u(v41<<(uint(int32(1))%32)-v44)<<(uint(int64(23))%64) + v38*(int64(1000)-v37)
				if v53 < int64(0) {
					v56 = int64(-500)
				} else {
					v56 = int64(500)
				}
				v59 = base.I64_div_s(v56+v53, int64(1000))
				if base.Ui64(v59) < base.Ui64(v38) {
					v61 = v38
				} else {
					v61 = v59
				}
				v62 = v61
			default:
				v37 = int64(700)
				v38 = *(*int64)(unsafe.Add(mBase, uint32(v8)+8))
				v41 = *(*int32)(unsafe.Add(mBase, uint32(v8)+16))
				v44 = *(*int32)(unsafe.Add(mBase, uint32(v8)+24))
				v53 = v37*base.I64_extend_i32_u(v41<<(uint(int32(1))%32)-v44)<<(uint(int64(23))%64) + v38*(int64(1000)-v37)
				if v53 < int64(0) {
					v56 = int64(-500)
				} else {
					v56 = int64(500)
				}
				v59 = base.I64_div_s(v56+v53, int64(1000))
				if base.Ui64(v59) < base.Ui64(v38) {
					v61 = v38
				} else {
					v61 = v59
				}
				v62 = v61
			}
		}
	} else {
		v37 = int64(627)
		v38 = *(*int64)(unsafe.Add(mBase, uint32(v8)+8))
		v41 = *(*int32)(unsafe.Add(mBase, uint32(v8)+16))
		v44 = *(*int32)(unsafe.Add(mBase, uint32(v8)+24))
		v53 = v37*base.I64_extend_i32_u(v41<<(uint(int32(1))%32)-v44)<<(uint(int64(23))%64) + v38*(int64(1000)-v37)
		if v53 < int64(0) {
			v56 = int64(-500)
		} else {
			v56 = int64(500)
		}
		v59 = base.I64_div_s(v56+v53, int64(1000))
		if base.Ui64(v59) < base.Ui64(v38) {
			v61 = v38
		} else {
			v61 = v59
		}
		v62 = v61
	}
	m.G0 = v8 + int32(32)
	return v62
}
func F_VP8LBitsEntropyUnrefined(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v4 int64
	_ = v4
	var v22 int32
	_ = v22
	var v26 int32
	_ = v26
	var v29 int64
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v36 int32
	_ = v36
	var v40 int32
	_ = v40
	var v43 int32
	_ = v43
	var v47 int32
	_ = v47
	var v53 int64
	_ = v53
	var v54 int32
	_ = v54
	var v57 int32
	_ = v57
	var v58 int64
	_ = v58
	var v59 int64
	_ = v59
	var v60 int64
	_ = v60
	var v63 int64
	_ = v63
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v71 int32
	_ = v71
	var v78 int64
	_ = v78
	var v79 int32
	_ = v79
	var v85 int32
	_ = v85
	var v91 int64
	_ = v91
	var v92 int32
	_ = v92
	var v95 int32
	_ = v95
	var v96 int64
	_ = v96
	var v100 int64
	_ = v100
	var v106 int64
	_ = v106
	v4 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(l2)+8)) = v4
	*(*int64)(unsafe.Add(mBase, uint32(l2+int32(16)))) = int64(281470681743360)
	if int32(1) <= l1 {
		v22 = int32(0)
		v26 = l0
		v29 = int64(0)
		v30 = v22
		v31 = v22
		v32 = v22
		v33 = v22
		for {
			v36 = *(*int32)(unsafe.Add(mBase, uint32(v26)))
			if v36 == int32(0) {
				v63 = v29
				v64 = v30
				v65 = v31
				v66 = v32
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = v33
				v40 = v30 + v36
				*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v40
				v43 = v32 + int32(1)
				*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v43
				if base.Ui32(int32(255)) < base.Ui32(v36) {
					v54 = m.G1
					v57 = *(*int32)(unsafe.Add(mBase, uint32(v54)+uint32(_c_F_VP8LBitsEntropyUnrefined[0])))
					v58 = m.T0[v57].(func(*base.Module, int32) int64)(m, v36)
					mBase = m.M
					v59 = v58
				} else {
					v47 = m.G1
					v53 = *(*int64)(unsafe.Add(mBase, uint32(v47+int32(_a_F_VP8LBitsEntropyUnrefined_0)+v36<<(uint(int32(3))%32))))
					v59 = v53
				}
				v60 = v29 + v59
				if base.Ui32(v36) <= base.Ui32(v31) {
					v63 = v60
					v64 = v40
					v65 = v31
					v66 = v43
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l2)+16)) = v36
					v63 = v60
					v64 = v40
					v65 = v36
					v66 = v43
				}
			}
			v71 = v33 + int32(1)
			if l1 != v71 {
				v26 = v26 + int32(4)
				v29 = v63
				v30 = v64
				v31 = v65
				v32 = v66
				v33 = v71
				continue
			} else {
				break
			}
			break
		}
		if base.Ui32(int32(255)) < base.Ui32(v64) {
			v92 = m.G1
			v95 = *(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_VP8LBitsEntropyUnrefined[0])))
			v96 = m.T0[v95].(func(*base.Module, int32) int64)(m, v64)
			mBase = m.M
			v100 = v63
			v106 = v96
		} else {
			v78 = v63
			v79 = v64
			v85 = m.G1
			v91 = *(*int64)(unsafe.Add(mBase, uint32(v85+int32(_a_F_VP8LBitsEntropyUnrefined_0)+v79<<(uint(int32(3))%32))))
			v100 = v78
			v106 = v91
		}
	} else {
		v78 = v4
		v79 = int32(0)
		v85 = m.G1
		v91 = *(*int64)(unsafe.Add(mBase, uint32(v85+int32(_a_F_VP8LBitsEntropyUnrefined_0)+v79<<(uint(int32(3))%32))))
		v100 = v78
		v106 = v91
	}
	*(*int64)(unsafe.Add(mBase, uint32(l2))) = v106 - v100
	return
}
func F_VP8LClearBackwardRefs(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v10 int32
	_ = v10
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	if v3 == int32(0) {
	} else {
		v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
		*(*int32)(unsafe.Add(mBase, uint32(v3))) = v6
	}
	v8 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v8
	v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v10
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = l0 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v8
	return
}
func F_VP8LColorCacheClear(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	if l0 == int32(0) {
	} else {
		v4 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		F_free(m, v4)
		mBase = m.M
		*(*int32)(unsafe.Add(mBase, uint32(l0))) = int32(0)
	}
	return
}
func F_VP8LColorCacheInit(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v6 int64
	_ = v6
	var v7 int32
	_ = v7
	var v14 int64
	_ = v14
	var v15 int32
	_ = v15
	var v16 int64
	_ = v16
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	v6 = base.I64_extend_i32_s(int32(1) << (uint(l1) % 32))
	v7 = int32(4)
	if v6 == int64(0) {
		v25 = F_calloc(m, base.I32_wrap_i64(v6), v7)
		mBase = m.M
		v27 = v25
	} else {
		v14 = base.I64_div_u_s(int64(2147418112), v6)
		v15 = int32(0)
		v16 = base.I64_extend_i32_u(v7)
		if base.Ui64(int64(4294967295)) < base.Ui64(v16*v6) {
			v27 = v15
		} else {
			if base.Ui64(v14) < base.Ui64(v16) {
				v27 = v15
			} else {
				v25 = F_calloc(m, base.I32_wrap_i64(v6), v7)
				mBase = m.M
				v27 = v25
			}
		}
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v27
	if v27 != 0 {
		*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = l1
		*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = int32(32) - l1
		return int32(1)
	} else {
		return int32(0)
	}
}
func F_VP8LColorSpaceTransform(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32) int32 {
	mBase := m.M
	_ = mBase
	var v45 int32
	_ = v45
	var v47 int32
	_ = v47
	var v49 int32
	_ = v49
	var v298 int32
	_ = v298
	var v306 int32
	_ = v306
	var v308 int32
	_ = v308
	var v310 int32
	_ = v310
	var v312 int32
	_ = v312
	var v316 int32
	_ = v316
	var v320 int32
	_ = v320
	var v323 int32
	_ = v323
	var v325 int32
	_ = v325
	var v348 int32
	_ = v348
	var v371 int32
	_ = v371
	var v372 int32
	_ = v372
	var v374 int32
	_ = v374
	var v375 int32
	_ = v375
	var v377 int32
	_ = v377
	var v381 int32
	_ = v381
	var v387 int32
	_ = v387
	var v429 int32
	_ = v429
	var v440 int32
	_ = v440
	var v447 int32
	_ = v447
	var v450 int32
	_ = v450
	var v453 int32
	_ = v453
	var v456 int32
	_ = v456
	var v457 int32
	_ = v457
	var v460 int32
	_ = v460
	var v465 int32
	_ = v465
	var v466 int32
	_ = v466
	var v476 int32
	_ = v476
	var v478 int32
	_ = v478
	var v480 int32
	_ = v480
	var v490 int32
	_ = v490
	var v497 int32
	_ = v497
	var v499 int32
	_ = v499
	var v506 int32
	_ = v506
	var v510 int32
	_ = v510
	var v517 int32
	_ = v517
	var v519 int32
	_ = v519
	var v522 int32
	_ = v522
	var v537 int32
	_ = v537
	var v558 int32
	_ = v558
	var v569 int32
	_ = v569
	var v570 int32
	_ = v570
	var v572 int32
	_ = v572
	var v574 int32
	_ = v574
	var v589 int32
	_ = v589
	var v610 int32
	_ = v610
	var v635 int32
	_ = v635
	var v657 int32
	_ = v657
	var v666 int32
	_ = v666
	var v667 int32
	_ = v667
	var v668 int32
	_ = v668
	var v672 int32
	_ = v672
	var v673 int32
	_ = v673
	var v677 int32
	_ = v677
	var v678 int32
	_ = v678
	var v682 int32
	_ = v682
	var v683 int32
	_ = v683
	var v687 int32
	_ = v687
	var v735 int32
	_ = v735
	var v737 int32
	_ = v737
	var v739 int32
	_ = v739
	var v785 int32
	_ = v785
	var v786 int32
	_ = v786
	var v787 int32
	_ = v787
	var v790 int32
	_ = v790
	var v791 int32
	_ = v791
	var v808 int32
	_ = v808
	var v811 int32
	_ = v811
	var v830 int32
	_ = v830
	var v836 int32
	_ = v836
	var v840 int32
	_ = v840
	var v843 int32
	_ = v843
	var v847 int32
	_ = v847
	var v850 int32
	_ = v850
	var v851 int32
	_ = v851
	var v853 int32
	_ = v853
	var v855 int32
	_ = v855
	var v858 int32
	_ = v858
	var v862 int32
	_ = v862
	var v870 int32
	_ = v870
	var v871 int32
	_ = v871
	var v872 int32
	_ = v872
	var v881 int32
	_ = v881
	var v882 int32
	_ = v882
	var v888 int32
	_ = v888
	var v895 int32
	_ = v895
	var v943 int32
	_ = v943
	var v990 int32
	_ = v990
	var v993 int32
	_ = v993
	var v995 int32
	_ = v995
	var v998 int32
	_ = v998
	var v1014 int32
	_ = v1014
	var v1017 int32
	_ = v1017
	var v1044 int32
	_ = v1044
	var v1046 int32
	_ = v1046
	var v1051 int32
	_ = v1051
	var v1113 int32
	_ = v1113
	var v1149 int32
	_ = v1149
	v45 = m.G0
	v47 = v45 - int32(2080)
	m.G0 = v47
	v49 = *(*int32)(unsafe.Add(mBase, uint32(l8)))
	goto L3
L1:
	;
	goto L17
L3:
	;
	base.MemoryFill(m, v47+int32(1040), int32(0), int32(1024))
	goto L1
L15:
	;
	v298 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v47)+14)) = uint8(v298)
	*(*uint16)(unsafe.Add(mBase, uint32(v47)+12)) = uint16(v298)
	*(*uint8)(unsafe.Add(mBase, uint32(v47)+11)) = uint8(v298)
	*(*uint16)(unsafe.Add(mBase, uint32(v47)+9)) = uint16(v298)
	v306 = int32(1)
	v308 = v306 << (uint(l2) % 32)
	v310 = v308 + int32(-1)
	v312 = int32(base.Ui32(v310+l1) >> (uint(l2) % 32))
	if v312 < v306 {
		goto L30
	} else {
		goto L31
	}
L17:
	;
	base.MemoryFill(m, v47+int32(16), int32(0), int32(1024))
	goto L15
L29:
	;
	m.G0 = v47 + int32(2080)
	return v1113
L30:
	;
	F_VP8LOptimizeSampling(m, l5, l0, l1, l2, int32(9), l9)
	mBase = m.M
	v1113 = v306
	goto L29
L31:
	;
	v316 = int32(base.Ui32(v310+l0) >> (uint(l2) % 32))
	if int32(0) < v316 {
		goto L33
	} else {
		goto L34
	}
L32:
	;
	v1014 = int32(0)
	v1017 = v312
	goto L91
L33:
	;
	v320 = int32(2)
	v323 = l0 << (uint(v320) % 32)
	v325 = l4 + int32(-8)
	v348 = int32(0)
	goto L35
L34:
	;
	goto L32
L35:
	;
	v371 = v348 << (uint(l2) % 32)
	v372 = v371 + v308
	if v372 < l1 {
		goto L37
	} else {
		goto L38
	}
L37:
	;
	v374 = v372
	goto L39
L38:
	;
	v374 = l1
	goto L39
L39:
	;
	v375 = l1 - v371
	if v308 < v375 {
		goto L40
	} else {
		goto L41
	}
L40:
	;
	v377 = v308
	goto L42
L41:
	;
	v377 = v375
	goto L42
L42:
	;
	v381 = v377 & int32(3)
	v387 = v371 * l0
	v429 = int32(0)
	goto L43
L43:
	;
	v440 = v429 + v348*v316
	if v348 == int32(0) {
		goto L46
	} else {
		goto L47
	}
L44:
	;
	v993 = base.I32_div_s(v348*l7, v312)
	v995 = F_WebPReportProgress(m, l6, v993+v49, l8)
	mBase = m.M
	if v995 != 0 {
		goto L87
	} else {
		goto L88
	}
L45:
	;
	v457 = v429 << (uint(l2) % 32)
	v460 = int32(2)
	v1149 = int32(8)
	v465 = v47 + int32(14)
	v466 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v465))))
	*(*uint8)(unsafe.Add(mBase, uint32(v47+v1149))) = uint8(v466)
	v476 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v47+int32(11)))))
	*(*uint8)(unsafe.Add(mBase, uint32(v47+int32(4)))) = uint8(v476)
	v478 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v47)+12)))
	*(*uint16)(unsafe.Add(mBase, uint32(v47)+6)) = uint16(v478)
	v480 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v47)+9)))
	*(*uint16)(unsafe.Add(mBase, uint32(v47)+2)) = uint16(v480)
	v490 = int32(16)
	F_GetBestColorTransformForTile(m, v47+int32(2072), v429, v456, l2, v47+int32(6), v47+v460, l3, l0, l1, v47+int32(1040), v47+v490, l4)
	mBase = m.M
	v497 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v47+int32(2074)))))
	*(*uint8)(unsafe.Add(mBase, uint32(v465))) = uint8(v497)
	v499 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v47)+2072)))
	*(*uint16)(unsafe.Add(mBase, uint32(v47)+12)) = uint16(v499)
	v506 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v47)+13)))
	v510 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v47)+12)))
	*(*int32)(unsafe.Add(mBase, uint32(l5+v440<<(uint(v460)%32)))) = v497<<(uint(v490)%32) | v506<<(uint(v1149)%32) | v510 | int32(-16777216)
	*(*uint8)(unsafe.Add(mBase, uint32(v47)+2074)) = uint8(v497)
	*(*uint16)(unsafe.Add(mBase, uint32(v47)+2072)) = uint16(v499)
	if v377 < int32(1) {
		goto L48
	} else {
		goto L49
	}
L46:
	;
	v456 = int32(0)
	goto L45
L47:
	;
	v447 = *(*int32)(unsafe.Add(mBase, uint32(l5+(v440-v316)<<(uint(int32(2))%32))))
	*(*uint8)(unsafe.Add(mBase, uint32(v47)+9)) = uint8(v447)
	v450 = int32(base.Ui32(v447) >> (uint(int32(16)) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v47)+11)) = uint8(v450)
	v453 = int32(base.Ui32(v447) >> (uint(int32(8)) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v47)+10)) = uint8(v453)
	v456 = v348
	goto L45
L48:
	;
	if v374 <= v371 {
		goto L63
	} else {
		goto L64
	}
L49:
	;
	v517 = l0 - v457
	if v308 < v517 {
		goto L50
	} else {
		goto L51
	}
L50:
	;
	v519 = v308
	goto L52
L51:
	;
	v519 = v517
	goto L52
L52:
	;
	v522 = l4 + v387<<(uint(int32(2))%32) + v457<<(uint(int32(2))%32)
	if v381 != 0 {
		goto L54
	} else {
		goto L55
	}
L53:
	;
	if base.Ui32(v377) < base.Ui32(int32(4)) {
		goto L48
	} else {
		goto L59
	}
L54:
	;
	v537 = v522
	v558 = v381
	goto L56
L55:
	;
	v589 = v522
	v610 = v377
	goto L53
L56:
	;
	v569 = m.G93
	v570 = *(*int32)(unsafe.Add(mBase, uint32(v569)))
	m.T0[v570].(func(*base.Module, int32, int32, int32))(m, v47+int32(2072), v537, v519)
	mBase = m.M
	v572 = v537 + v323
	v574 = v558 + int32(-1)
	if v574 != 0 {
		v537 = v572
		v558 = v574
		goto L56
	} else {
		goto L58
	}
L57:
	;
	v589 = v572
	v610 = v377 & int32(-4)
	goto L53
L58:
	;
	goto L57
L59:
	;
	v635 = v589
	v657 = v610 + int32(-1)
	goto L60
L60:
	;
	v666 = v47 + int32(2072)
	v667 = m.G93
	v668 = *(*int32)(unsafe.Add(mBase, uint32(v667)))
	m.T0[v668].(func(*base.Module, int32, int32, int32))(m, v666, v635, v519)
	mBase = m.M
	v672 = v635 + v323
	v673 = *(*int32)(unsafe.Add(mBase, uint32(v667)))
	m.T0[v673].(func(*base.Module, int32, int32, int32))(m, v666, v672, v519)
	mBase = m.M
	v677 = v672 + v323
	v678 = *(*int32)(unsafe.Add(mBase, uint32(v667)))
	m.T0[v678].(func(*base.Module, int32, int32, int32))(m, v666, v677, v519)
	mBase = m.M
	v682 = v677 + v323
	v683 = *(*int32)(unsafe.Add(mBase, uint32(v667)))
	m.T0[v683].(func(*base.Module, int32, int32, int32))(m, v666, v682, v519)
	mBase = m.M
	v687 = v657 + int32(-4)
	if base.Ui32(v687) < base.Ui32(int32(-2)) {
		v635 = v682 + v323
		v657 = v687
		goto L60
	} else {
		goto L62
	}
L61:
	;
	goto L48
L62:
	;
	goto L61
L63:
	;
	v990 = v429 + int32(1)
	if v990 != v316 {
		v429 = v990
		goto L43
	} else {
		goto L86
	}
L64:
	;
	v735 = v457 + v308
	if v735 < l0 {
		goto L65
	} else {
		goto L66
	}
L65:
	;
	v737 = v735
	goto L67
L66:
	;
	v737 = l0
	goto L67
L67:
	;
	v739 = int32(2)
	v785 = v325 + v323*(v371+int32(-1)) + v457<<(uint(v739)%32)
	v786 = v325 + (v387+v457)<<(uint(v739)%32)
	v787 = v371
	goto L68
L68:
	;
	v790 = v787 * l0
	v791 = v790 + v457
	if v790+v737 <= v791 {
		goto L70
	} else {
		goto L71
	}
L69:
	;
	goto L63
L70:
	;
	v943 = v787 + int32(1)
	if v943 != v374 {
		v785 = v785 + v323
		v786 = v786 + v323
		v787 = v943
		goto L68
	} else {
		goto L85
	}
L71:
	;
	v808 = v786
	v811 = v791
	v830 = v785
	v836 = v737 - v457
	goto L72
L72:
	;
	v840 = *(*int32)(unsafe.Add(mBase, uint32(v808+int32(8))))
	if v811 < int32(2) {
		goto L75
	} else {
		goto L76
	}
L73:
	;
	goto L70
L74:
	;
	v888 = int32(4)
	v895 = v836 + int32(-1)
	if v895 != 0 {
		v808 = v808 + v888
		v811 = v811 + int32(1)
		v830 = v830 + v888
		v836 = v895
		goto L72
	} else {
		goto L84
	}
L75:
	;
	if v811 < l0+v320 {
		goto L79
	} else {
		goto L80
	}
L76:
	;
	v843 = *(*int32)(unsafe.Add(mBase, uint32(v808)))
	if v840 != v843 {
		goto L75
	} else {
		goto L77
	}
L77:
	;
	v847 = *(*int32)(unsafe.Add(mBase, uint32(v808+int32(4))))
	if v840 == v847 {
		goto L74
	} else {
		goto L78
	}
L78:
	;
	goto L75
L79:
	;
	v870 = v47 + int32(1040) + int32(base.Ui32(v840)>>(uint(int32(14))%32))&int32(1020)
	v871 = *(*int32)(unsafe.Add(mBase, uint32(v870)))
	v872 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v870))) = v871 + v872
	v881 = v47 + int32(16) + v840&int32(255)<<(uint(int32(2))%32)
	v882 = *(*int32)(unsafe.Add(mBase, uint32(v881)))
	*(*int32)(unsafe.Add(mBase, uint32(v881))) = v882 + v872
	goto L74
L80:
	;
	v850 = *(*int32)(unsafe.Add(mBase, uint32(v808)))
	v851 = *(*int32)(unsafe.Add(mBase, uint32(v830)))
	if v850 != v851 {
		goto L79
	} else {
		goto L81
	}
L81:
	;
	v853 = int32(4)
	v855 = *(*int32)(unsafe.Add(mBase, uint32(v808+v853)))
	v858 = *(*int32)(unsafe.Add(mBase, uint32(v830+v853)))
	if v855 != v858 {
		goto L79
	} else {
		goto L82
	}
L82:
	;
	v862 = *(*int32)(unsafe.Add(mBase, uint32(v830+int32(8))))
	if v840 == v862 {
		goto L74
	} else {
		goto L83
	}
L83:
	;
	goto L79
L84:
	;
	goto L73
L85:
	;
	goto L69
L86:
	;
	goto L44
L87:
	;
	v998 = v348 + int32(1)
	if v998 != v312 {
		v348 = v998
		goto L35
	} else {
		goto L89
	}
L88:
	;
	v1113 = int32(0)
	goto L29
L89:
	;
	goto L30
L90:
	;
	v1113 = int32(0)
	goto L29
L91:
	;
	v1044 = base.I32_div_s(v1014, v312)
	v1046 = F_WebPReportProgress(m, l6, v1044+v49, l8)
	mBase = m.M
	if v1046 == int32(0) {
		goto L90
	} else {
		goto L93
	}
L93:
	;
	v1051 = v1017 + int32(-1)
	if v1051 == int32(0) {
		goto L30
	} else {
		goto L94
	}
L94:
	;
	v1014 = v1014 + l7
	v1017 = v1051
	goto L91
}
func F_VP8LCreateCompressedHuffmanTree(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v40 int32
	_ = v40
	var v44 int32
	_ = v44
	var v53 int32
	_ = v53
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v65 int32
	_ = v65
	var v68 int32
	_ = v68
	var v74 int32
	_ = v74
	var v77 int32
	_ = v77
	var v86 int32
	_ = v86
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v95 int32
	_ = v95
	var v106 int32
	_ = v106
	var v128 int32
	_ = v128
	var v131 int32
	_ = v131
	var v132 int32
	_ = v132
	var v140 int32
	_ = v140
	var v141 int32
	_ = v141
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v182 int64
	_ = v182
	var v185 int32
	_ = v185
	var v186 int32
	_ = v186
	var v198 int32
	_ = v198
	var v219 int32
	_ = v219
	var v224 int32
	_ = v224
	var v228 int32
	_ = v228
	var v235 int32
	_ = v235
	var v241 int32
	_ = v241
	var v244 int32
	_ = v244
	var v255 int32
	_ = v255
	var v262 int32
	_ = v262
	var v263 int32
	_ = v263
	var v269 int32
	_ = v269
	var v275 int32
	_ = v275
	var v280 int32
	_ = v280
	var v286 int32
	_ = v286
	var v290 int32
	_ = v290
	var v293 int32
	_ = v293
	var v295 int32
	_ = v295
	var v299 int32
	_ = v299
	var v306 int32
	_ = v306
	var v311 int32
	_ = v311
	var v325 int32
	_ = v325
	var v335 int32
	_ = v335
	var v336 int32
	_ = v336
	var v361 int32
	_ = v361
	var v363 int32
	_ = v363
	var v367 int32
	_ = v367
	var v381 int32
	_ = v381
	var v386 int32
	_ = v386
	var v392 int32
	_ = v392
	var v396 int32
	_ = v396
	var v398 int32
	_ = v398
	var v402 int32
	_ = v402
	var v404 int32
	_ = v404
	var v417 int32
	_ = v417
	v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if v15 < int32(1) {
		v417 = l1
		goto L1
	} else {
		goto L2
	}
L1:
	;
	return (v417 - l1) >> (uint(int32(1)) % 32)
L2:
	;
	v23 = l1
	v25 = int32(8)
	v26 = int32(0)
	goto L3
L3:
	;
	v35 = v26 + int32(1)
	if v35 < v15 {
		goto L5
	} else {
		goto L6
	}
L4:
	;
	v417 = v402
	goto L1
L5:
	;
	v37 = v15
	goto L7
L6:
	;
	v37 = v35
	goto L7
L7:
	;
	v40 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v44 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v40+v26))))
	v53 = v26
	goto L9
L8:
	;
	v68 = v65 - v26
	if v44 != 0 {
		goto L15
	} else {
		goto L16
	}
L9:
	;
	if v37+int32(-1) != v53 {
		goto L11
	} else {
		goto L12
	}
L10:
	;
	v65 = v62
	goto L8
L11:
	;
	v62 = v53 + int32(1)
	v63 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v40+int32(1)+v53))))
	if v63 == v44 {
		v53 = v62
		goto L9
	} else {
		goto L13
	}
L12:
	;
	v65 = v37
	goto L8
L13:
	;
	goto L10
L14:
	;
	if v65 < v15 {
		v23 = v402
		v25 = v404
		v26 = v65
		goto L3
	} else {
		goto L62
	}
L15:
	;
	if v25 == v44 {
		v262 = v23
		v263 = v68
		goto L41
	} else {
		goto L42
	}
L16:
	;
	if v68 < int32(1) {
		v402 = v23
		v404 = v25
		goto L14
	} else {
		goto L17
	}
L17:
	;
	v74 = base.I32_rem_u_s(v68+int32(-1), int32(138))
	v77 = v68 << (uint(int32(1)) % 32)
	v86 = v68
	v88 = v77
	v89 = v77
	v91 = int32(0)
	goto L18
L18:
	;
	v92 = v23 + v91
	if base.Ui32(int32(2)) < base.Ui32(v86) {
		goto L20
	} else {
		goto L21
	}
L19:
	;
	v402 = v23 + v244
	v404 = v25
	goto L14
L20:
	;
	if base.Ui32(int32(10)) < base.Ui32(v86) {
		goto L36
	} else {
		goto L37
	}
L21:
	;
	v95 = int32(0)
	if base.Ui32(v88) < base.Ui32(int32(33)) {
		goto L23
	} else {
		goto L24
	}
L22:
	;
	v402 = v23 + v89
	v404 = v25
	goto L14
L23:
	;
	if v88 == int32(0) {
		goto L25
	} else {
		goto L26
	}
L24:
	;
	base.MemoryFill(m, v92, v95, v88)
	goto L22
L25:
	;
	goto L22
L26:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v92))) = uint8(v95)
	v106 = v92 + v88
	*(*uint8)(unsafe.Add(mBase, uint32(v106+int32(-1)))) = uint8(v95)
	if base.Ui32(v88) < base.Ui32(int32(3)) {
		goto L25
	} else {
		goto L27
	}
L27:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v92)+2)) = uint8(v95)
	*(*uint8)(unsafe.Add(mBase, uint32(v92)+1)) = uint8(v95)
	*(*uint8)(unsafe.Add(mBase, uint32(v106+int32(-3)))) = uint8(v95)
	*(*uint8)(unsafe.Add(mBase, uint32(v106+int32(-2)))) = uint8(v95)
	if base.Ui32(v88) < base.Ui32(int32(7)) {
		goto L25
	} else {
		goto L28
	}
L28:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v92)+3)) = uint8(v95)
	*(*uint8)(unsafe.Add(mBase, uint32(v106+int32(-4)))) = uint8(v95)
	if base.Ui32(v88) < base.Ui32(int32(9)) {
		goto L25
	} else {
		goto L29
	}
L29:
	;
	v128 = int32(0)
	v131 = (v128 - v92) & int32(3)
	v132 = v92 + v131
	*(*int32)(unsafe.Add(mBase, uint32(v132))) = v128
	v140 = (v88 - v131) & int32(60)
	v141 = v132 + v140
	*(*int32)(unsafe.Add(mBase, uint32(v141+int32(-4)))) = v128
	if base.Ui32(v140) < base.Ui32(int32(9)) {
		goto L25
	} else {
		goto L30
	}
L30:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v132)+8)) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v132)+4)) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v141+int32(-8)))) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v141+int32(-12)))) = v128
	if base.Ui32(v140) < base.Ui32(int32(25)) {
		goto L25
	} else {
		goto L31
	}
L31:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v132)+24)) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v132)+20)) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v132)+16)) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v132)+12)) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v141+int32(-16)))) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v141+int32(-20)))) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v141+int32(-24)))) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v141+int32(-28)))) = v128
	v176 = v132&int32(4) | int32(24)
	v177 = v140 - v176
	if base.Ui32(v177) < base.Ui32(int32(32)) {
		goto L25
	} else {
		goto L32
	}
L32:
	;
	v182 = base.I64_extend_i32_u(v128) * int64(4294967297)
	v185 = v177
	v186 = v132 + v176
	goto L33
L33:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v186)+24)) = v182
	*(*int64)(unsafe.Add(mBase, uint32(v186)+16)) = v182
	*(*int64)(unsafe.Add(mBase, uint32(v186)+8)) = v182
	*(*int64)(unsafe.Add(mBase, uint32(v186))) = v182
	v198 = v185 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v198) {
		v185 = v198
		v186 = v186 + int32(32)
		goto L33
	} else {
		goto L35
	}
L34:
	;
	goto L25
L35:
	;
	goto L34
L36:
	;
	v228 = int32(18)
	*(*uint8)(unsafe.Add(mBase, uint32(v92))) = uint8(v228)
	if base.Ui32(int32(138)) < base.Ui32(v86) {
		goto L38
	} else {
		goto L39
	}
L37:
	;
	v219 = int32(17)
	*(*uint8)(unsafe.Add(mBase, uint32(v92))) = uint8(v219)
	v224 = v86 + int32(-3)
	*(*uint8)(unsafe.Add(mBase, uint32(v92+int32(1)))) = uint8(v224)
	v402 = v92 + int32(2)
	v404 = v25
	goto L14
L38:
	;
	v241 = int32(127)
	*(*uint8)(unsafe.Add(mBase, uint32(v92+int32(1)))) = uint8(v241)
	v244 = v91 + int32(2)
	if int32(138) < v86 {
		v86 = v86 + int32(-138)
		v88 = v88 + int32(-276)
		v89 = v89 + int32(-274)
		v91 = v244
		goto L18
	} else {
		goto L40
	}
L39:
	;
	v235 = v74 + int32(-10)
	*(*uint8)(unsafe.Add(mBase, uint32(v92+int32(1)))) = uint8(v235)
	v402 = v92 + int32(2)
	v404 = v25
	goto L14
L40:
	;
	goto L19
L41:
	;
	if int32(1) <= v263 {
		goto L43
	} else {
		goto L44
	}
L42:
	;
	v255 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v23)+1)) = uint8(v255)
	*(*uint8)(unsafe.Add(mBase, uint32(v23))) = uint8(v44)
	v262 = v23 + int32(2)
	v263 = v68 + int32(-1)
	goto L41
L43:
	;
	v269 = base.I32_rem_u_s(v263+int32(-1), int32(6))
	if base.Ui32(v263) < base.Ui32(int32(3)) {
		v306 = v262
		v311 = v263
		goto L45
	} else {
		goto L46
	}
L44:
	;
	v402 = v262
	v404 = v44
	goto L14
L45:
	;
	if base.Ui32(v311+int32(-1)) < base.Ui32(int32(7)) {
		v367 = v306
		goto L52
	} else {
		goto L53
	}
L46:
	;
	v275 = v262
	v280 = v263
	goto L48
L47:
	;
	v299 = v269 + int32(-2)
	*(*uint8)(unsafe.Add(mBase, uint32(v275)+1)) = uint8(v299)
	v402 = v275 + int32(2)
	v404 = v44
	goto L14
L48:
	;
	v286 = int32(16)
	*(*uint8)(unsafe.Add(mBase, uint32(v275))) = uint8(v286)
	if base.Ui32(v280) <= base.Ui32(int32(6)) {
		goto L47
	} else {
		goto L50
	}
L50:
	;
	v290 = int32(3)
	*(*uint8)(unsafe.Add(mBase, uint32(v275)+1)) = uint8(v290)
	v293 = v275 + int32(2)
	v295 = v280 + int32(-6)
	if base.Ui32(v295) < base.Ui32(v290) {
		v306 = v293
		v311 = v295
		goto L45
	} else {
		goto L51
	}
L51:
	;
	v275 = v293
	v280 = v295
	goto L48
L52:
	;
	if v311 != 0 {
		goto L57
	} else {
		goto L58
	}
L53:
	;
	v325 = v306
	v335 = int32(0)
	goto L54
L54:
	;
	v336 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+15)) = uint8(v336)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+14)) = uint8(v44)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+13)) = uint8(v336)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+12)) = uint8(v44)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+11)) = uint8(v336)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+10)) = uint8(v44)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+9)) = uint8(v336)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+8)) = uint8(v44)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+7)) = uint8(v336)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+6)) = uint8(v44)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+5)) = uint8(v336)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+4)) = uint8(v44)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+3)) = uint8(v336)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+2)) = uint8(v44)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+1)) = uint8(v336)
	*(*uint8)(unsafe.Add(mBase, uint32(v325))) = uint8(v44)
	v361 = v325 + int32(16)
	v363 = v335 + int32(-8)
	if v363 != 0 {
		v325 = v361
		v335 = v363
		goto L54
	} else {
		goto L56
	}
L55:
	;
	v367 = v361
	goto L52
L56:
	;
	goto L55
L57:
	;
	v381 = v367
	v386 = v311
	goto L59
L58:
	;
	v402 = v367
	v404 = v44
	goto L14
L59:
	;
	v392 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v381)+1)) = uint8(v392)
	*(*uint8)(unsafe.Add(mBase, uint32(v381))) = uint8(v44)
	v396 = v381 + int32(2)
	v398 = v386 + int32(-1)
	if v398 != 0 {
		v381 = v396
		v386 = v398
		goto L59
	} else {
		goto L61
	}
L60:
	;
	v402 = v396
	v404 = v44
	goto L14
L61:
	;
	goto L60
L62:
	;
	goto L4
}
func F_VP8LDistanceToPlaneCode(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v7 int32
	_ = v7
	var v12 int32
	_ = v12
	var v21 int32
	_ = v21
	var v34 int32
	_ = v34
	var v40 int32
	_ = v40
	v5 = base.I32_div_s(l1, l0)
	v7 = l1 - v5*l0
	if int32(7) < v5 {
		if int32(6) < v5 {
			return l1 + int32(120)
		} else {
			if v7 <= l0+int32(-8) {
				return l1 + int32(120)
			} else {
				v34 = m.G1
				v40 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v5<<(uint(int32(4))%32)-v7+(v34+int32(_a_F_VP8LDistanceToPlaneCode_0))+int32(24)))))
				return v40 + int32(1)
			}
		}
	} else {
		if int32(8) < v7 {
			if int32(6) < v5 {
				return l1 + int32(120)
			} else {
				if v7 <= l0+int32(-8) {
					return l1 + int32(120)
				} else {
					v34 = m.G1
					v40 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v5<<(uint(int32(4))%32)-v7+(v34+int32(_a_F_VP8LDistanceToPlaneCode_0))+int32(24)))))
					return v40 + int32(1)
				}
			}
		} else {
			v12 = m.G1
			v21 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v12+int32(_a_F_VP8LDistanceToPlaneCode_0)+(v5<<(uint(int32(4))%32)|int32(8)-v7)))))
			return v21 + int32(1)
		}
	}
}
func F_VP8LDspInitSSE41(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v5 int32
	_ = v5
	var v7 int32
	_ = v7
	var v9 int32
	_ = v9
	var v11 int32
	_ = v11
	var v13 int32
	_ = v13
	var v17 int32
	_ = v17
	var v19 int32
	_ = v19
	v4 = m.G2
	v5 = m.G24
	v7 = v4 + int32(87)
	*(*int32)(unsafe.Add(mBase, uint32(v5))) = v7
	v9 = m.G23
	v11 = v4 + int32(88)
	*(*int32)(unsafe.Add(mBase, uint32(v9))) = v11
	v13 = m.G28
	*(*int32)(unsafe.Add(mBase, uint32(v13))) = v4 + int32(89)
	v17 = m.G31
	*(*int32)(unsafe.Add(mBase, uint32(v17))) = v11
	v19 = m.G32
	*(*int32)(unsafe.Add(mBase, uint32(v19))) = v7
	return
}
func F_VP8LEncDspInitSSE41(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v5 int32
	_ = v5
	var v7 int32
	_ = v7
	var v9 int32
	_ = v9
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v17 int32
	_ = v17
	var v19 int32
	_ = v19
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	v4 = m.G2
	v5 = m.G94
	v7 = v4 + int32(292)
	*(*int32)(unsafe.Add(mBase, uint32(v5))) = v7
	v9 = m.G112
	*(*int32)(unsafe.Add(mBase, uint32(v9))) = v4 + int32(293)
	v13 = m.G95
	v15 = v4 + int32(294)
	*(*int32)(unsafe.Add(mBase, uint32(v13))) = v15
	v17 = m.G96
	v19 = v4 + int32(295)
	*(*int32)(unsafe.Add(mBase, uint32(v17))) = v19
	v21 = m.G103
	*(*int32)(unsafe.Add(mBase, uint32(v21))) = v7
	v23 = m.G105
	*(*int32)(unsafe.Add(mBase, uint32(v23))) = v15
	v25 = m.G106
	*(*int32)(unsafe.Add(mBase, uint32(v25))) = v19
	return
}
func F_VP8LHashChainClear(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	v2 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	F_free(m, v2)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(l0))) = int64(0)
	return
}
func F_VP8LHashChainFill(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) int32 {
	mBase := m.M
	_ = mBase
	var v33 int32
	_ = v33
	var v45 int32
	_ = v45
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v56 int32
	_ = v56
	var v62 int32
	_ = v62
	var v63 int64
	_ = v63
	var v64 int32
	_ = v64
	var v71 int64
	_ = v71
	var v72 int32
	_ = v72
	var v73 int64
	_ = v73
	var v83 int32
	_ = v83
	var v85 int32
	_ = v85
	var v88 int32
	_ = v88
	var v94 int32
	_ = v94
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v100 int32
	_ = v100
	var v226 int32
	_ = v226
	var v228 int32
	_ = v228
	var v230 int32
	_ = v230
	var v231 int32
	_ = v231
	var v235 int32
	_ = v235
	var v255 int32
	_ = v255
	var v267 int32
	_ = v267
	var v268 int32
	_ = v268
	var v269 int32
	_ = v269
	var v272 int32
	_ = v272
	var v274 int32
	_ = v274
	var v278 int32
	_ = v278
	var v282 int32
	_ = v282
	var v295 int32
	_ = v295
	var v303 int32
	_ = v303
	var v323 int32
	_ = v323
	var v328 int32
	_ = v328
	var v334 int32
	_ = v334
	var v337 int32
	_ = v337
	var v338 int32
	_ = v338
	var v340 int32
	_ = v340
	var v342 int32
	_ = v342
	var v353 int32
	_ = v353
	var v378 int32
	_ = v378
	var v379 int32
	_ = v379
	var v383 int32
	_ = v383
	var v387 int32
	_ = v387
	var v388 int32
	_ = v388
	var v423 int32
	_ = v423
	var v424 int32
	_ = v424
	var v429 int64
	_ = v429
	var v432 int32
	_ = v432
	var v433 int32
	_ = v433
	var v445 int32
	_ = v445
	var v469 int32
	_ = v469
	var v477 int32
	_ = v477
	var v478 int32
	_ = v478
	var v483 int32
	_ = v483
	var v485 int32
	_ = v485
	var v488 int32
	_ = v488
	var v491 int32
	_ = v491
	var v500 int32
	_ = v500
	var v503 int32
	_ = v503
	var v504 int32
	_ = v504
	var v507 int32
	_ = v507
	var v509 int32
	_ = v509
	var v510 int32
	_ = v510
	var v513 int32
	_ = v513
	var v522 int32
	_ = v522
	var v531 int32
	_ = v531
	var v561 int32
	_ = v561
	var v563 int32
	_ = v563
	var v565 int32
	_ = v565
	var v566 int32
	_ = v566
	var v567 int32
	_ = v567
	var v570 int32
	_ = v570
	var v572 int32
	_ = v572
	var v573 int32
	_ = v573
	var v575 int32
	_ = v575
	var v579 int32
	_ = v579
	var v583 int32
	_ = v583
	var v585 int32
	_ = v585
	var v586 int32
	_ = v586
	var v587 int32
	_ = v587
	var v589 int32
	_ = v589
	var v590 int32
	_ = v590
	var v591 int32
	_ = v591
	var v592 int32
	_ = v592
	var v594 int32
	_ = v594
	var v595 int32
	_ = v595
	var v597 int32
	_ = v597
	var v598 int32
	_ = v598
	var v599 int32
	_ = v599
	var v601 int32
	_ = v601
	var v603 int32
	_ = v603
	var v605 int32
	_ = v605
	var v607 int32
	_ = v607
	var v609 int32
	_ = v609
	var v611 int32
	_ = v611
	var v612 int32
	_ = v612
	var v613 int32
	_ = v613
	var v614 int32
	_ = v614
	var v617 int32
	_ = v617
	var v618 int32
	_ = v618
	var v621 int32
	_ = v621
	var v623 int32
	_ = v623
	var v626 int32
	_ = v626
	var v627 int32
	_ = v627
	var v628 int32
	_ = v628
	var v631 int32
	_ = v631
	var v634 int32
	_ = v634
	var v637 int32
	_ = v637
	var v641 int32
	_ = v641
	var v643 int32
	_ = v643
	var v646 int32
	_ = v646
	var v654 int32
	_ = v654
	var v663 int32
	_ = v663
	var v664 int32
	_ = v664
	var v675 int32
	_ = v675
	var v678 int32
	_ = v678
	var v679 int32
	_ = v679
	var v680 int32
	_ = v680
	var v684 int32
	_ = v684
	var v686 int32
	_ = v686
	var v687 int32
	_ = v687
	var v688 int32
	_ = v688
	var v690 int32
	_ = v690
	var v695 int32
	_ = v695
	var v696 int32
	_ = v696
	var v697 int32
	_ = v697
	var v698 int32
	_ = v698
	var v701 int32
	_ = v701
	var v704 int32
	_ = v704
	var v725 int32
	_ = v725
	var v735 int32
	_ = v735
	var v736 int32
	_ = v736
	var v738 int32
	_ = v738
	var v742 int32
	_ = v742
	var v758 int32
	_ = v758
	var v761 int32
	_ = v761
	var v769 int32
	_ = v769
	var v773 int32
	_ = v773
	var v780 int32
	_ = v780
	var v790 int32
	_ = v790
	var v791 int32
	_ = v791
	var v799 int32
	_ = v799
	var v800 int32
	_ = v800
	var v803 int32
	_ = v803
	var v804 int32
	_ = v804
	var v809 int32
	_ = v809
	var v815 int32
	_ = v815
	var v818 int32
	_ = v818
	var v826 int32
	_ = v826
	var v830 int32
	_ = v830
	var v846 int32
	_ = v846
	var v849 int32
	_ = v849
	var v853 int32
	_ = v853
	var v857 int32
	_ = v857
	var v862 int32
	_ = v862
	var v868 int32
	_ = v868
	var v898 int32
	_ = v898
	var v900 int32
	_ = v900
	var v904 int32
	_ = v904
	var v907 int32
	_ = v907
	var v931 int32
	_ = v931
	var v939 int32
	_ = v939
	var v952 int32
	_ = v952
	var v953 int32
	_ = v953
	var v960 int32
	_ = v960
	var v962 int32
	_ = v962
	var v963 int32
	_ = v963
	var v973 int32
	_ = v973
	var v976 int32
	_ = v976
	var v984 int32
	_ = v984
	var v995 int32
	_ = v995
	var v1004 int32
	_ = v1004
	var v1006 int32
	_ = v1006
	var v1008 int32
	_ = v1008
	var v1009 int32
	_ = v1009
	var v1020 int32
	_ = v1020
	var v1021 int32
	_ = v1021
	var v1027 int32
	_ = v1027
	var v1033 int32
	_ = v1033
	var v1035 int32
	_ = v1035
	var v1055 int32
	_ = v1055
	var v1067 int32
	_ = v1067
	var v1069 int32
	_ = v1069
	v33 = l4 * l3
	if int32(75) < l1 {
		v47 = int32(1048456)
		goto L1
	} else {
		goto L2
	}
L1:
	;
	v48 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if int32(2) < v33 {
		goto L8
	} else {
		goto L9
	}
L2:
	;
	if l1 < int32(51) {
		goto L3
	} else {
		goto L4
	}
L3:
	;
	if int32(25) < l1 {
		goto L5
	} else {
		goto L6
	}
L4:
	;
	v47 = l3 << (uint(int32(8)) % 32)
	goto L1
L5:
	;
	v45 = int32(6)
	goto L7
L6:
	;
	v45 = int32(4)
	goto L7
L7:
	;
	v47 = l3 << (uint(v45) % 32)
	goto L1
L8:
	;
	v62 = *(*int32)(unsafe.Add(mBase, uint32(l8)))
	v63 = int64(262144)
	v64 = int32(4)
	goto L14
L9:
	;
	v56 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v48+v33<<(uint(int32(2))%32)+int32(-4)))) = v56
	*(*int32)(unsafe.Add(mBase, uint32(v48))) = v56
	return int32(1)
L10:
	;
	v94 = int32(base.Ui32(l1*l1) >> (uint(int32(7)) % 32))
	v96 = v94 + int32(8)
	v97 = int32(1048456)
	if v47 < v97 {
		goto L21
	} else {
		goto L22
	}
L11:
	;
	if v85 != 0 {
		goto L10
	} else {
		goto L17
	}
L12:
	;
	goto L11
L13:
	;
	v83 = F_malloc(m, base.I32_wrap_i64(v63)*v64)
	mBase = m.M
	v85 = v83
	goto L12
L14:
	;
	v71 = base.I64_div_u_s(int64(2147418112), v63)
	v72 = int32(0)
	v73 = base.I64_extend_i32_u(v64)
	if base.Ui64(int64(4294967295)) < base.Ui64(v73*v63) {
		v85 = v72
		goto L12
	} else {
		goto L15
	}
L15:
	;
	if base.Ui64(v71) < base.Ui64(v73) {
		v85 = v72
		goto L12
	} else {
		goto L16
	}
L16:
	;
	goto L13
L17:
	;
	v88 = *(*int32)(unsafe.Add(mBase, uint32(l6)+92))
	if v88 != 0 {
		goto L19
	} else {
		goto L20
	}
L18:
	;
	return int32(0)
L19:
	;
	goto L18
L20:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l6)+92)) = int32(1)
	goto L19
L21:
	;
	v100 = v47
	goto L23
L22:
	;
	v100 = v97
	goto L23
L23:
	;
	goto L26
L24:
	;
	v226 = v33 + int32(-2)
	v228 = base.I32_div_s(l7, int32(2))
	v230 = *(*int32)(unsafe.Add(mBase, uint32(l2)))
	v231 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	v235 = int32(0)
	v255 = base.B2i32(v230 == v231)
	goto L38
L26:
	;
	base.MemoryFill(m, v85, int32(255), int32(1048576))
	goto L24
L38:
	;
	if v226 <= v235 {
		goto L42
	} else {
		goto L43
	}
L39:
	;
	F_free(m, v85)
	mBase = m.M
	goto L150
L40:
	;
	v1067 = base.I32_div_s(v1035*v228, v226)
	v1069 = F_WebPReportProgress(m, l6, v1067+v62, l8)
	mBase = m.M
	if v1069 != 0 {
		v235 = v1035
		v255 = v1055
		goto L38
	} else {
		goto L149
	}
L41:
	;
	v939 = v282 * int32(1540483478)
	if v931&int32(1) != 0 {
		goto L143
	} else {
		goto L144
	}
L42:
	;
	v483 = v235 << (uint(int32(2)) % 32)
	v485 = l2 + v483
	v488 = *(*int32)(unsafe.Add(mBase, uint32(v485+int32(4))))
	v491 = *(*int32)(unsafe.Add(mBase, uint32(v485)))
	v500 = *(*int32)(unsafe.Add(mBase, uint32(v85+int32(base.Ui32(v488*int32(-962287725)+v491*int32(1540483478))>>(uint(int32(12))%32))&int32(1048572))))
	*(*int32)(unsafe.Add(mBase, uint32(v48+v483))) = v500
	F_free(m, v85)
	mBase = m.M
	goto L71
L43:
	;
	v267 = int32(2)
	v268 = v235 << (uint(v267) % 32)
	v269 = l2 + v268
	v272 = *(*int32)(unsafe.Add(mBase, uint32(v269+int32(8))))
	v274 = v235 + int32(1)
	v278 = *(*int32)(unsafe.Add(mBase, uint32(l2+v274<<(uint(v267)%32))))
	if v255 == int32(0) {
		goto L44
	} else {
		goto L45
	}
L44:
	;
	v469 = *(*int32)(unsafe.Add(mBase, uint32(v269)))
	v477 = v85 + int32(base.Ui32(v278*int32(-962287725)+v469*int32(1540483478))>>(uint(int32(12))%32))&int32(1048572)
	v478 = *(*int32)(unsafe.Add(mBase, uint32(v477)))
	*(*int32)(unsafe.Add(mBase, uint32(v48+v268))) = v478
	*(*int32)(unsafe.Add(mBase, uint32(v477))) = v235
	v1035 = v274
	v1055 = base.B2i32(v278 == v272)
	goto L40
L45:
	;
	if v278 != v272 {
		goto L44
	} else {
		goto L46
	}
L46:
	;
	v282 = *(*int32)(unsafe.Add(mBase, uint32(v269)))
	if v33 <= v235+int32(3) {
		v907 = v235
		v931 = int32(1)
		goto L41
	} else {
		goto L47
	}
L47:
	;
	v295 = l2 + int32(12) + v268
	v303 = int32(3)
	goto L50
L48:
	;
	if base.Ui32(v334) < base.Ui32(int32(_a_F_VP8LHashChainFill_0)) {
		goto L54
	} else {
		goto L55
	}
L49:
	;
	v334 = v303 + int32(-2)
	goto L48
L50:
	;
	v323 = *(*int32)(unsafe.Add(mBase, uint32(v295)))
	if v323 != v282 {
		goto L49
	} else {
		goto L52
	}
L52:
	;
	v328 = v303 + int32(1)
	if v33-v235 != v328 {
		v295 = v295 + int32(4)
		v303 = v328
		goto L50
	} else {
		goto L53
	}
L53:
	;
	v334 = v226 - v235
	goto L48
L54:
	;
	if v334 != 0 {
		v907 = v235
		v931 = v334
		goto L41
	} else {
		goto L70
	}
L55:
	;
	v337 = v48 + v268
	v338 = int32(255)
	v340 = v334 + int32(-4095)
	v342 = v340 << (uint(int32(2)) % 32)
	if base.Ui32(v342) < base.Ui32(int32(33)) {
		goto L57
	} else {
		goto L58
	}
L56:
	;
	v907 = v340 + v235
	v931 = int32(4095)
	goto L41
L57:
	;
	if v342 == int32(0) {
		goto L59
	} else {
		goto L60
	}
L58:
	;
	base.MemoryFill(m, v337, v338, v342)
	goto L56
L59:
	;
	goto L56
L60:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v337))) = uint8(v338)
	v353 = v337 + v342
	*(*uint8)(unsafe.Add(mBase, uint32(v353+int32(-1)))) = uint8(v338)
	if base.Ui32(v342) < base.Ui32(int32(3)) {
		goto L59
	} else {
		goto L61
	}
L61:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v337)+2)) = uint8(v338)
	*(*uint8)(unsafe.Add(mBase, uint32(v337)+1)) = uint8(v338)
	*(*uint8)(unsafe.Add(mBase, uint32(v353+int32(-3)))) = uint8(v338)
	*(*uint8)(unsafe.Add(mBase, uint32(v353+int32(-2)))) = uint8(v338)
	if base.Ui32(v342) < base.Ui32(int32(7)) {
		goto L59
	} else {
		goto L62
	}
L62:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v337)+3)) = uint8(v338)
	*(*uint8)(unsafe.Add(mBase, uint32(v353+int32(-4)))) = uint8(v338)
	if base.Ui32(v342) < base.Ui32(int32(9)) {
		goto L59
	} else {
		goto L63
	}
L63:
	;
	v378 = (int32(0) - v337) & int32(3)
	v379 = v337 + v378
	v383 = int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v379))) = v383
	v387 = (v342 - v378) & int32(60)
	v388 = v379 + v387
	*(*int32)(unsafe.Add(mBase, uint32(v388+int32(-4)))) = v383
	if base.Ui32(v387) < base.Ui32(int32(9)) {
		goto L59
	} else {
		goto L64
	}
L64:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v379)+8)) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v379)+4)) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v388+int32(-8)))) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v388+int32(-12)))) = v383
	if base.Ui32(v387) < base.Ui32(int32(25)) {
		goto L59
	} else {
		goto L65
	}
L65:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v379)+24)) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v379)+20)) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v379)+16)) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v379)+12)) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v388+int32(-16)))) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v388+int32(-20)))) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v388+int32(-24)))) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v388+int32(-28)))) = v383
	v423 = v379&int32(4) | int32(24)
	v424 = v387 - v423
	if base.Ui32(v424) < base.Ui32(int32(32)) {
		goto L59
	} else {
		goto L66
	}
L66:
	;
	v429 = base.I64_extend_i32_u(v383) * int64(4294967297)
	v432 = v424
	v433 = v379 + v423
	goto L67
L67:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v433)+24)) = v429
	*(*int64)(unsafe.Add(mBase, uint32(v433)+16)) = v429
	*(*int64)(unsafe.Add(mBase, uint32(v433)+8)) = v429
	*(*int64)(unsafe.Add(mBase, uint32(v433))) = v429
	v445 = v432 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v445) {
		v432 = v445
		v433 = v433 + int32(32)
		goto L67
	} else {
		goto L69
	}
L68:
	;
	goto L59
L69:
	;
	goto L68
L70:
	;
	v1035 = v235
	v1055 = int32(0)
	goto L40
L71:
	;
	v503 = v62 + v228
	v504 = F_WebPReportProgress(m, l6, v503, l8)
	mBase = m.M
	if v504 != 0 {
		goto L72
	} else {
		goto L73
	}
L72:
	;
	v507 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v509 = v33 + int32(-1)
	v510 = int32(2)
	v513 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v507+v509<<(uint(v510)%32)))) = v513
	*(*int32)(unsafe.Add(mBase, uint32(v507))) = v513
	v522 = v94 + int32(7)
	v531 = v226
	goto L75
L73:
	;
	return int32(0)
L74:
	;
	v904 = F_WebPReportProgress(m, l6, v62+l7, l8)
	mBase = m.M
	return v904
L75:
	;
	if v531 == int32(0) {
		goto L74
	} else {
		goto L77
	}
L76:
	;
	return int32(0)
L77:
	;
	v561 = int32(0)
	v563 = v531 - v100
	if base.Ui32(v531) < base.Ui32(v563) {
		goto L78
	} else {
		goto L79
	}
L78:
	;
	v565 = v561
	goto L80
L79:
	;
	v565 = v563
	goto L80
L80:
	;
	v566 = v509 - v531
	v567 = int32(4095)
	if v566 < v567 {
		goto L81
	} else {
		goto L82
	}
L81:
	;
	v570 = v566
	goto L83
L82:
	;
	v570 = v567
	goto L83
L83:
	;
	v572 = v531 << (uint(int32(2)) % 32)
	v573 = l2 + v572
	v575 = *(*int32)(unsafe.Add(mBase, uint32(v48+v572)))
	if l5 == int32(0) {
		goto L85
	} else {
		goto L86
	}
L84:
	;
	if v627 < v565 {
		v704 = v626
		v725 = v631
		goto L107
	} else {
		goto L108
	}
L85:
	;
	v579 = int32(0)
	if base.Ui32(v531) < base.Ui32(l3) {
		v598 = v579
		v599 = v579
		v601 = v96
		goto L87
	} else {
		goto L88
	}
L86:
	;
	v626 = v561
	v627 = v575
	v628 = v96
	v631 = int32(0)
	goto L84
L87:
	;
	v603 = v573 + int32(-4)
	v605 = v598 << (uint(int32(2)) % 32)
	v607 = *(*int32)(unsafe.Add(mBase, uint32(v603+v605)))
	v609 = *(*int32)(unsafe.Add(mBase, uint32(v573+v605)))
	if v607 != v609 {
		v614 = v579
		goto L96
	} else {
		goto L97
	}
L88:
	;
	v583 = int32(0)
	v585 = v573 + (v513-l3)<<(uint(v510)%32)
	v586 = *(*int32)(unsafe.Add(mBase, uint32(v585)))
	v587 = *(*int32)(unsafe.Add(mBase, uint32(v573)))
	if v586 != v587 {
		v598 = v583
		v599 = v583
		v601 = v522
		goto L87
	} else {
		goto L89
	}
L89:
	;
	v589 = m.G100
	v590 = *(*int32)(unsafe.Add(mBase, uint32(v589)))
	v591 = m.T0[v590].(func(*base.Module, int32, int32, int32) int32)(m, v585, v573, v570)
	mBase = m.M
	v592 = int32(0)
	v594 = base.B2i32(v592 < v591)
	if v592 < v591 {
		goto L90
	} else {
		goto L91
	}
L90:
	;
	v595 = v591
	goto L92
L91:
	;
	v595 = v592
	goto L92
L92:
	;
	if v592 < v591 {
		goto L93
	} else {
		goto L94
	}
L93:
	;
	v597 = l3
	goto L95
L94:
	;
	v597 = int32(0)
	goto L95
L95:
	;
	v598 = v595
	v599 = v597
	v601 = v522
	goto L87
L96:
	;
	v617 = base.B2i32(v598 < v614)
	if v598 < v614 {
		goto L98
	} else {
		goto L99
	}
L97:
	;
	v611 = m.G100
	v612 = *(*int32)(unsafe.Add(mBase, uint32(v611)))
	v613 = m.T0[v612].(func(*base.Module, int32, int32, int32) int32)(m, v603, v573, v570)
	mBase = m.M
	v614 = v613
	goto L96
L98:
	;
	v618 = v614
	goto L100
L99:
	;
	v618 = v598
	goto L100
L100:
	;
	if v618 == int32(4095) {
		goto L101
	} else {
		goto L102
	}
L101:
	;
	v621 = v565 + int32(-1)
	goto L103
L102:
	;
	v621 = v575
	goto L103
L103:
	;
	if v598 < v614 {
		goto L104
	} else {
		goto L105
	}
L104:
	;
	v623 = int32(1)
	goto L106
L105:
	;
	v623 = v599
	goto L106
L106:
	;
	v626 = v618
	v627 = v621
	v628 = v601 + int32(-1)
	v631 = v623
	goto L84
L107:
	;
	v735 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v736 = v735 + v572
	v738 = v725 << (uint(int32(12)) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v736))) = v704 | v738
	v742 = v531 + int32(-1)
	if v725 == int32(0) {
		v868 = v742
		goto L121
	} else {
		goto L122
	}
L108:
	;
	v634 = int32(256)
	if v566 < v634 {
		goto L109
	} else {
		goto L110
	}
L109:
	;
	v637 = v566
	goto L111
L110:
	;
	v637 = v634
	goto L111
L111:
	;
	v641 = *(*int32)(unsafe.Add(mBase, uint32(v573+v626<<(uint(int32(2))%32))))
	v643 = v626
	v646 = v627
	v654 = v628
	v663 = v641
	v664 = v631
	goto L112
L112:
	;
	v675 = v654 + int32(-1)
	if v675 == int32(0) {
		v704 = v643
		v725 = v664
		goto L107
	} else {
		goto L114
	}
L113:
	;
	v704 = v696
	v725 = v698
	goto L107
L114:
	;
	v678 = int32(2)
	v679 = v646 << (uint(v678) % 32)
	v680 = l2 + v679
	v684 = *(*int32)(unsafe.Add(mBase, uint32(v680+v643<<(uint(v678)%32))))
	if v684 != v663 {
		v696 = v643
		v697 = v663
		v698 = v664
		goto L115
	} else {
		goto L116
	}
L115:
	;
	v701 = *(*int32)(unsafe.Add(mBase, uint32(v48+v679)))
	if v565 <= v701 {
		v643 = v696
		v646 = v701
		v654 = v675
		v663 = v697
		v664 = v698
		goto L112
	} else {
		goto L120
	}
L116:
	;
	v686 = m.G100
	v687 = *(*int32)(unsafe.Add(mBase, uint32(v686)))
	v688 = m.T0[v687].(func(*base.Module, int32, int32, int32) int32)(m, v680, v573, v570)
	mBase = m.M
	if v688 <= v643 {
		v696 = v643
		v697 = v663
		v698 = v664
		goto L115
	} else {
		goto L117
	}
L117:
	;
	v690 = v531 - v646
	if v688 < v637 {
		goto L118
	} else {
		goto L119
	}
L118:
	;
	v695 = *(*int32)(unsafe.Add(mBase, uint32(v573+v688<<(uint(int32(2))%32))))
	v696 = v688
	v697 = v695
	v698 = v690
	goto L115
L119:
	;
	v704 = v688
	v725 = v690
	goto L107
L120:
	;
	goto L113
L121:
	;
	v898 = base.I32_div_u_s((v226-v868)*(l7-v228), v226)
	v900 = F_WebPReportProgress(m, l6, v898+v503, l8)
	mBase = m.M
	if v900 != 0 {
		v531 = v868
		goto L75
	} else {
		goto L141
	}
L122:
	;
	if v742 == int32(0) {
		v868 = v742
		goto L121
	} else {
		goto L123
	}
L123:
	;
	if base.Ui32(v742) < base.Ui32(v725) {
		v868 = v742
		goto L121
	} else {
		goto L124
	}
L124:
	;
	if v725 == int32(1) {
		goto L125
	} else {
		goto L126
	}
L125:
	;
	v815 = v704
	v818 = v742
	v826 = l2 + int32(-8) + v572
	v830 = v736 + int32(-4)
	goto L137
L126:
	;
	v758 = v704
	v761 = v742
	v769 = l2 + int32(-4) + v572
	v773 = v531
	v780 = v736 + int32(-4)
	goto L127
L127:
	;
	v790 = *(*int32)(unsafe.Add(mBase, uint32(v769+(int32(0)-v725<<(uint(int32(2))%32)))))
	v791 = *(*int32)(unsafe.Add(mBase, uint32(v769)))
	if v790 != v791 {
		v868 = v761
		goto L121
	} else {
		goto L129
	}
L129:
	;
	if v758 != int32(4095) {
		goto L130
	} else {
		goto L131
	}
L130:
	;
	v799 = base.B2i32(v758 < int32(4095))
	v800 = v758 + v799
	*(*int32)(unsafe.Add(mBase, uint32(v780))) = v800 | v738
	if v758 < int32(4095) {
		goto L133
	} else {
		goto L134
	}
L131:
	;
	if base.Ui32(v761+int32(4095)) < base.Ui32(v773) {
		v868 = v761
		goto L121
	} else {
		goto L132
	}
L132:
	;
	goto L130
L133:
	;
	v803 = v761
	goto L135
L134:
	;
	v803 = v773
	goto L135
L135:
	;
	v804 = int32(-4)
	v809 = v761 + int32(-1)
	if base.Ui32(v725) <= base.Ui32(v809) {
		v758 = v800
		v761 = v809
		v769 = v769 + v804
		v773 = v803
		v780 = v780 + v804
		goto L127
	} else {
		goto L136
	}
L136:
	;
	v868 = v809
	goto L121
L137:
	;
	v846 = *(*int32)(unsafe.Add(mBase, uint32(v826)))
	v849 = *(*int32)(unsafe.Add(mBase, uint32(v826+int32(4))))
	if v846 != v849 {
		v868 = v818
		goto L121
	} else {
		goto L139
	}
L138:
	;
	v868 = int32(0)
	goto L121
L139:
	;
	v853 = v815 + base.B2i32(v815 < int32(4095))
	*(*int32)(unsafe.Add(mBase, uint32(v830))) = v853 | int32(_a_F_VP8LHashChainFill_0)
	v857 = int32(-4)
	v862 = v818 + int32(-1)
	if v862 != 0 {
		v815 = v853
		v818 = v862
		v826 = v826 + v857
		v830 = v830 + v857
		goto L137
	} else {
		goto L140
	}
L140:
	;
	goto L138
L141:
	;
	goto L76
L142:
	;
	v963 = int32(0)
	if v931 == int32(1) {
		v1035 = v960
		v1055 = v963
		goto L40
	} else {
		goto L145
	}
L143:
	;
	v952 = v85 + int32(base.Ui32(v931*int32(-962287725)+v939)>>(uint(int32(12))%32))&int32(1048572)
	v953 = *(*int32)(unsafe.Add(mBase, uint32(v952)))
	*(*int32)(unsafe.Add(mBase, uint32(v48+v907<<(uint(int32(2))%32)))) = v953
	*(*int32)(unsafe.Add(mBase, uint32(v952))) = v907
	v960 = v907 + int32(1)
	v962 = v931 + int32(-1)
	goto L142
L144:
	;
	v960 = v907
	v962 = v931
	goto L142
L145:
	;
	v973 = v960
	v976 = v939 - v962*int32(962287725)
	v984 = v48 + v960<<(uint(int32(2))%32)
	v995 = v962
	goto L146
L146:
	;
	v1004 = int32(12)
	v1006 = int32(1048572)
	v1008 = v85 + int32(base.Ui32(v976)>>(uint(v1004)%32))&v1006
	v1009 = *(*int32)(unsafe.Add(mBase, uint32(v1008)))
	*(*int32)(unsafe.Add(mBase, uint32(v984))) = v1009
	*(*int32)(unsafe.Add(mBase, uint32(v1008))) = v973
	v1020 = v85 + int32(base.Ui32(v976+int32(962287725))>>(uint(v1004)%32))&v1006
	v1021 = *(*int32)(unsafe.Add(mBase, uint32(v1020)))
	*(*int32)(unsafe.Add(mBase, uint32(v984+int32(4)))) = v1021
	*(*int32)(unsafe.Add(mBase, uint32(v1020))) = v973 + int32(1)
	v1027 = v973 + int32(2)
	v1033 = v995 + int32(-2)
	if v1033 != 0 {
		v973 = v1027
		v976 = v976 + int32(1924575450)
		v984 = v984 + int32(8)
		v995 = v1033
		goto L146
	} else {
		goto L148
	}
L147:
	;
	v1035 = v1027
	v1055 = v963
	goto L40
L148:
	;
	goto L147
L149:
	;
	goto L39
L150:
	;
	return int32(0)
}
func F_VP8LHashChainInit(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int64
	_ = v4
	var v5 int32
	_ = v5
	var v12 int64
	_ = v12
	var v13 int32
	_ = v13
	var v14 int64
	_ = v14
	var v24 int32
	_ = v24
	var v26 int32
	_ = v26
	v4 = base.I64_extend_i32_s(l1)
	v5 = int32(4)
	if v4 == int64(0) {
		v24 = F_malloc(m, base.I32_wrap_i64(v4)*v5)
		mBase = m.M
		v26 = v24
	} else {
		v12 = base.I64_div_u_s(int64(2147418112), v4)
		v13 = int32(0)
		v14 = base.I64_extend_i32_u(v5)
		if base.Ui64(int64(4294967295)) < base.Ui64(v14*v4) {
			v26 = v13
		} else {
			if base.Ui64(v12) < base.Ui64(v14) {
				v26 = v13
			} else {
				v24 = F_malloc(m, base.I32_wrap_i64(v4)*v5)
				mBase = m.M
				v26 = v24
			}
		}
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v26
	if v26 != 0 {
		*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = l1
		return int32(1)
	} else {
		return int32(0)
	}
}
func F_VP8LOptimizeSampling(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v49 int32
	_ = v49
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v65 int32
	_ = v65
	var v69 int32
	_ = v69
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v96 int32
	_ = v96
	var v102 int32
	_ = v102
	var v103 int32
	_ = v103
	var v104 int32
	_ = v104
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v111 int32
	_ = v111
	var v116 int32
	_ = v116
	var v121 int32
	_ = v121
	var v125 int32
	_ = v125
	var v163 int32
	_ = v163
	var v180 int32
	_ = v180
	var v198 int32
	_ = v198
	var v211 int32
	_ = v211
	var v236 int32
	_ = v236
	var v237 int32
	_ = v237
	var v254 int32
	_ = v254
	var v257 int32
	_ = v257
	var v266 int32
	_ = v266
	var v268 int32
	_ = v268
	var v273 int32
	_ = v273
	var v274 int32
	_ = v274
	var v296 int32
	_ = v296
	var v298 int32
	_ = v298
	var v301 int32
	_ = v301
	var v304 int32
	_ = v304
	var v312 int32
	_ = v312
	var v349 int32
	_ = v349
	var v361 int32
	_ = v361
	var v362 int32
	_ = v362
	var v364 int32
	_ = v364
	var v368 int32
	_ = v368
	var v371 int32
	_ = v371
	var v377 int32
	_ = v377
	var v378 int32
	_ = v378
	var v385 int32
	_ = v385
	var v395 int32
	_ = v395
	var v396 int32
	_ = v396
	var v409 int32
	_ = v409
	var v422 int32
	_ = v422
	var v430 int32
	_ = v430
	var v432 int32
	_ = v432
	var v435 int32
	_ = v435
	var v437 int32
	_ = v437
	var v445 int32
	_ = v445
	var v455 int32
	_ = v455
	var v465 int32
	_ = v465
	var v470 int32
	_ = v470
	var v487 int32
	_ = v487
	var v503 int32
	_ = v503
	var v507 int32
	_ = v507
	var v516 int32
	_ = v516
	var v528 int32
	_ = v528
	var v535 int32
	_ = v535
	var v562 int32
	_ = v562
	*(*int32)(unsafe.Add(mBase, uint32(l5))) = l3
	if l4 <= l3 {
		goto L1
	} else {
		goto L2
	}
L1:
	;
	return
L2:
	;
	v27 = int32(1) << (uint(l3) % 32)
	v28 = int32(-1)
	v29 = l2 + v28
	v31 = int32(base.Ui32(v27+v29) >> (uint(l3) % 32))
	v33 = l1 + v28
	v35 = int32(base.Ui32(v27+v33) >> (uint(l3) % 32))
	v37 = v35 << (uint(int32(2)) % 32)
	v49 = l3
	goto L4
L3:
	;
	if v163 == l3 {
		goto L1
	} else {
		goto L21
	}
L4:
	;
	v61 = int32(1)
	v62 = v49 + v61
	v65 = v61 << (uint(v49-l3) % 32)
	if v31 <= v65 {
		goto L6
	} else {
		goto L7
	}
L5:
	;
	v163 = l4
	goto L3
L6:
	;
	if v62 != l4 {
		v49 = v62
		goto L4
	} else {
		goto L20
	}
L7:
	;
	v69 = int32(1) << (uint(v62-l3) % 32)
	v73 = v65
	v74 = l0
	goto L8
L8:
	;
	v96 = int32(0)
	if v37 == v96 {
		v121 = v96
		goto L11
	} else {
		goto L12
	}
L9:
	;
	goto L6
L10:
	;
	if v121 != 0 {
		v163 = v49
		goto L3
	} else {
		goto L18
	}
L11:
	;
	goto L10
L12:
	;
	v102 = v74
	v103 = v74 + v37*v65
	v104 = v37
	goto L14
L13:
	;
	v121 = v108 - v109
	goto L11
L14:
	;
	v108 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v102))))
	v109 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v103))))
	if v108 != v109 {
		goto L13
	} else {
		goto L16
	}
L16:
	;
	v111 = int32(1)
	v116 = v104 + int32(-1)
	if v116 != 0 {
		v102 = v102 + v111
		v103 = v103 + v111
		v104 = v116
		goto L14
	} else {
		goto L17
	}
L17:
	;
	v121 = v96
	goto L11
L18:
	;
	v125 = v73 + v69
	if v125 < v31 {
		v73 = v125
		v74 = v74 + v37*v69
		goto L8
	} else {
		goto L19
	}
L19:
	;
	goto L9
L20:
	;
	goto L5
L21:
	;
	if v31 < int32(1) {
		v349 = v163
		goto L22
	} else {
		goto L23
	}
L22:
	;
	if v349 == l3 {
		goto L1
	} else {
		goto L49
	}
L23:
	;
	if v35 < int32(1) {
		v349 = v163
		goto L22
	} else {
		goto L24
	}
L24:
	;
	if v163 < l3 {
		goto L25
	} else {
		goto L26
	}
L25:
	;
	v180 = v163
	goto L27
L26:
	;
	v180 = l3
	goto L27
L27:
	;
	if v163 <= l3 {
		goto L28
	} else {
		goto L29
	}
L28:
	;
	v349 = v180
	goto L22
L29:
	;
	v198 = v163
	goto L30
L30:
	;
	v211 = int32(1) << (uint(v198-l3) % 32)
	v236 = int32(0)
	v237 = l0 + int32(4)
	goto L32
L32:
	;
	v254 = v237
	v257 = int32(0)
	goto L34
L34:
	;
	v266 = v257 + v211
	if v266 < v35 {
		goto L36
	} else {
		goto L37
	}
L35:
	;
	v312 = v236 + int32(1)
	if v312 < v31 {
		v236 = v312
		v237 = v237 + v35<<(uint(int32(2))%32)
		goto L32
	} else {
		goto L48
	}
L36:
	;
	v268 = v266
	goto L38
L37:
	;
	v268 = v35
	goto L38
L38:
	;
	v273 = v257
	v274 = v254
	goto L40
L39:
	;
	if v296 < v268 {
		goto L45
	} else {
		goto L46
	}
L40:
	;
	v296 = v273 + int32(1)
	if v268 <= v296 {
		goto L39
	} else {
		goto L42
	}
L41:
	;
	v304 = v198 + int32(-1)
	if l3 < v304 {
		v198 = v304
		goto L30
	} else {
		goto L44
	}
L42:
	;
	v298 = *(*int32)(unsafe.Add(mBase, uint32(v274)))
	v301 = *(*int32)(unsafe.Add(mBase, uint32(l0+v236*v35<<(uint(int32(2))%32)+v257<<(uint(int32(2))%32))))
	if v298 == v301 {
		v273 = v296
		v274 = v274 + int32(4)
		goto L40
	} else {
		goto L43
	}
L43:
	;
	goto L41
L44:
	;
	goto L28
L45:
	;
	goto L35
L46:
	;
	if v266 < v35 {
		v254 = v254 + v211<<(uint(int32(2))%32)
		v257 = v266
		goto L34
	} else {
		goto L47
	}
L47:
	;
	goto L45
L48:
	;
	v349 = v198
	goto L22
L49:
	;
	v361 = int32(1)
	v362 = v361 << (uint(v349) % 32)
	v364 = int32(base.Ui32(v362+v29) >> (uint(v349) % 32))
	if v364 < v361 {
		goto L50
	} else {
		goto L51
	}
L50:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l5))) = v349
	goto L1
L51:
	;
	v368 = int32(base.Ui32(v362+v33) >> (uint(v349) % 32))
	if v368 < int32(1) {
		goto L50
	} else {
		goto L52
	}
L52:
	;
	v371 = v349 - l3
	v377 = v368 & int32(3)
	v378 = int32(0)
	v385 = v378
	v395 = v378
	v396 = l0
	goto L53
L53:
	;
	if base.Ui32(v368) < base.Ui32(int32(4)) {
		v487 = int32(0)
		goto L55
	} else {
		goto L56
	}
L54:
	;
	goto L50
L55:
	;
	if v377 == int32(0) {
		goto L60
	} else {
		goto L61
	}
L56:
	;
	v409 = v396
	v422 = int32(0)
	goto L57
L57:
	;
	v430 = v395 + v422
	v432 = int32(2)
	v435 = *(*int32)(unsafe.Add(mBase, uint32(l0+v430<<(uint(v371)%32)<<(uint(v432)%32))))
	*(*int32)(unsafe.Add(mBase, uint32(v409))) = v435
	v437 = int32(4)
	v445 = *(*int32)(unsafe.Add(mBase, uint32(l0+(v430+int32(1))<<(uint(v371)%32)<<(uint(v432)%32))))
	*(*int32)(unsafe.Add(mBase, uint32(v409+v437))) = v445
	v455 = *(*int32)(unsafe.Add(mBase, uint32(l0+(v430+v432)<<(uint(v371)%32)<<(uint(v432)%32))))
	*(*int32)(unsafe.Add(mBase, uint32(v409+int32(8)))) = v455
	v465 = *(*int32)(unsafe.Add(mBase, uint32(l0+(v430+int32(3))<<(uint(v371)%32)<<(uint(v432)%32))))
	*(*int32)(unsafe.Add(mBase, uint32(v409+int32(12)))) = v465
	v470 = v422 + v437
	if v368&int32(2147483644) != v470 {
		v409 = v409 + int32(16)
		v422 = v470
		goto L57
	} else {
		goto L59
	}
L58:
	;
	v487 = v470
	goto L55
L59:
	;
	goto L58
L60:
	;
	v562 = v385 + int32(1)
	if v562 != v364 {
		v385 = v562
		v395 = v395 + v35
		v396 = v396 + v368<<(uint(int32(2))%32)
		goto L53
	} else {
		goto L65
	}
L61:
	;
	v503 = v487 + v395
	v507 = v396 + v487<<(uint(int32(2))%32)
	v516 = v377
	goto L62
L62:
	;
	v528 = *(*int32)(unsafe.Add(mBase, uint32(l0+v503<<(uint(v371)%32)<<(uint(int32(2))%32))))
	*(*int32)(unsafe.Add(mBase, uint32(v507))) = v528
	v535 = v516 + int32(-1)
	if v535 != 0 {
		v503 = v503 + int32(1)
		v507 = v507 + int32(4)
		v516 = v535
		goto L62
	} else {
		goto L64
	}
L63:
	;
	goto L60
L64:
	;
	goto L63
L65:
	;
	goto L54
}
func F_VP8LPutBitsInternal(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
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
	var v35 int32
	_ = v35
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v42 int32
	_ = v42
	var v51 int32
	_ = v51
	var v52 int32
	_ = v52
	var v55 int64
	_ = v55
	var v58 int32
	_ = v58
	var v60 int64
	_ = v60
	var v63 int32
	_ = v63
	var v66 int64
	_ = v66
	var v70 int32
	_ = v70
	var v72 int32
	_ = v72
	var v76 int32
	_ = v76
	var v83 int64
	_ = v83
	var v84 int32
	_ = v84
	var v85 int64
	_ = v85
	var v95 int32
	_ = v95
	var v97 int32
	_ = v97
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v111 int32
	_ = v111
	var v113 int32
	_ = v113
	var v118 int32
	_ = v118
	var v119 int32
	_ = v119
	var v126 int32
	_ = v126
	var v129 int32
	_ = v129
	var v131 int32
	_ = v131
	var v135 int32
	_ = v135
	var v136 int32
	_ = v136
	var v137 int32
	_ = v137
	var v138 int32
	_ = v138
	if l2 < int32(1) {
		goto L1
	} else {
		goto L2
	}
L1:
	;
	return
L2:
	;
	v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v16+l2 < int32(32) {
		goto L5
	} else {
		goto L6
	}
L3:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v138 + v136
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v135<<(uint(v138)%32) | v137
	goto L1
L4:
	;
	v34 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v35 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v39 = v31
	v40 = v32
	v41 = v35
	v42 = v34
	goto L8
L5:
	;
	if v16 < int32(16) {
		v135 = l1
		v136 = l2
		v137 = v15
		v138 = v16
		goto L3
	} else {
		goto L7
	}
L6:
	;
	v20 = int32(32)
	v21 = v20 - v16
	v29 = int32(base.Ui32(l1) >> (uint(v21) % 32))
	v30 = l2 - v21
	v31 = l1<<(uint(v16)%32) | v15
	v32 = v20
	goto L4
L7:
	;
	v29 = l1
	v30 = l2
	v31 = v15
	v32 = v16
	goto L4
L8:
	;
	if base.Ui32(v41+int32(2)) <= base.Ui32(v42) {
		v118 = v41
		v119 = v42
		goto L10
	} else {
		goto L11
	}
L9:
	;
	v135 = v29
	v136 = v30
	v137 = v131
	v138 = v129
	goto L3
L10:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v118))) = uint16(v39)
	v126 = v118 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v126
	v129 = v40 + int32(-16)
	v131 = int32(base.Ui32(v39) >> (uint(int32(16)) % 32))
	if int32(31) < v40 {
		v39 = v131
		v40 = v129
		v41 = v126
		v42 = v119
		goto L8
	} else {
		goto L32
	}
L11:
	;
	v51 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v52 = v42 - v51
	v55 = base.I64_extend_i32_s(v52) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v55) {
		v100 = v51
		goto L13
	} else {
		goto L14
	}
L12:
	;
	if v41 == v51 {
		goto L29
	} else {
		goto L30
	}
L13:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v100
	return
L14:
	;
	v58 = v41 - v51
	v60 = v55 + base.I64_extend_i32_u(v58)
	if base.Ui64(int64(4294967295)) < base.Ui64(v60) {
		v100 = v51
		goto L13
	} else {
		goto L15
	}
L15:
	;
	v63 = base.I32_wrap_i64(v60)
	if v42 == v51 {
		goto L16
	} else {
		goto L17
	}
L16:
	;
	v66 = int64(1)
	v70 = int32(base.Ui32(v52*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v63) < base.Ui32(v70) {
		goto L19
	} else {
		goto L20
	}
L17:
	;
	if base.Ui32(v63) <= base.Ui32(v52) {
		v118 = v41
		v119 = v42
		goto L10
	} else {
		goto L18
	}
L18:
	;
	goto L16
L19:
	;
	v72 = v70
	goto L21
L20:
	;
	v72 = v63
	goto L21
L21:
	;
	v76 = v72&int32(-1024) + int32(1024)
	goto L25
L22:
	;
	if v97 != 0 {
		goto L12
	} else {
		goto L28
	}
L23:
	;
	goto L22
L24:
	;
	v95 = F_malloc(m, base.I32_wrap_i64(v66)*v76)
	mBase = m.M
	v97 = v95
	goto L23
L25:
	;
	v83 = base.I64_div_u_s(int64(2147418112), v66)
	v84 = int32(0)
	v85 = base.I64_extend_i32_u(v76)
	if base.Ui64(int64(4294967295)) < base.Ui64(v85*v66) {
		v97 = v84
		goto L23
	} else {
		goto L26
	}
L26:
	;
	if base.Ui64(v83) < base.Ui64(v85) {
		v97 = v84
		goto L23
	} else {
		goto L27
	}
L27:
	;
	goto L24
L28:
	;
	v99 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v100 = v99
	goto L13
L29:
	;
	v111 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_free(m, v111)
	mBase = m.M
	goto L31
L30:
	;
	v109 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v110 = F_memcpy(m, v97, v109, v58)
	mBase = m.M
	goto L29
L31:
	;
	v113 = v97 + v76
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v113
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v97
	v118 = v97 + v58
	v119 = v113
	goto L10
L32:
	;
	goto L9
}
func F_VP8LRefsCursorInit(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	v4 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v4
	if v4 != 0 {
		v8 = *(*int32)(unsafe.Add(mBase, uint32(v4)+4))
		v9 = *(*int32)(unsafe.Add(mBase, uint32(v4)+8))
		v13 = v8 + v9<<(uint(int32(3))%32)
		v14 = v8
	} else {
		v6 = int32(0)
		v13 = v6
		v14 = v6
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v13
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v14
	return
}
func F_VP8LRefsCursorNextBlock(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v6 int32
	_ = v6
	var v7 int32
	_ = v7
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v6 = *(*int32)(unsafe.Add(mBase, uint32(v5)))
	if v6 != 0 {
		v9 = *(*int32)(unsafe.Add(mBase, uint32(v6)+4))
		v10 = *(*int32)(unsafe.Add(mBase, uint32(v6)+8))
		v14 = v9
		v15 = v9 + v10<<(uint(int32(3))%32)
	} else {
		v7 = int32(0)
		v14 = v7
		v15 = v7
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v15
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v14
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v6
	return
}
