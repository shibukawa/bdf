//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"math"
	"unsafe"
)

func F_LD4_C(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v14 int32
	_ = v14
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v29 int32
	_ = v29
	var v35 int32
	_ = v35
	var v42 int32
	_ = v42
	var v45 int32
	_ = v45
	var v54 int32
	_ = v54
	var v57 int32
	_ = v57
	var v68 int32
	_ = v68
	var v73 int32
	_ = v73
	var v80 int32
	_ = v80
	var v85 int32
	_ = v85
	var v96 int32
	_ = v96
	v9 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-29)))))
	v10 = int32(2)
	v11 = v9 + v10
	v14 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-31)))))
	v18 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-30)))))
	v19 = int32(1)
	v23 = int32(base.Ui32(v11+v14+v18<<(uint(v19)%32)) >> (uint(v10) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+32)) = uint8(v23)
	v26 = v18 + v10
	v29 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-32)))))
	v35 = int32(base.Ui32(v26+v29+v14<<(uint(v19)%32)) >> (uint(v10) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v35)
	v42 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-28)))))
	v45 = int32(base.Ui32(v26+v9<<(uint(v19)%32)+v42) >> (uint(v10) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+64)) = uint8(v45)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v23)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+33)) = uint8(v45)
	v54 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-27)))))
	v57 = int32(base.Ui32(v11+v42<<(uint(v19)%32)+v54) >> (uint(v10) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+96)) = uint8(v57)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v45)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+65)) = uint8(v57)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+34)) = uint8(v57)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)) = uint8(v57)
	v68 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-26)))))
	v73 = int32(base.Ui32(v42+v54<<(uint(v19)%32)+v68+v10) >> (uint(v10) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+97)) = uint8(v73)
	v80 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-25)))))
	v85 = int32(base.Ui32(v54+v68<<(uint(v19)%32)+v80+v10) >> (uint(v10) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+98)) = uint8(v85)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+35)) = uint8(v73)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+66)) = uint8(v73)
	v96 = int32(base.Ui32(v68+v80+v80<<(uint(v19)%32)+v10) >> (uint(v10) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+99)) = uint8(v96)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+67)) = uint8(v85)
	return
}
func F_Launch(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v6 int32
	_ = v6
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	if v3 == int32(0) {
	} else {
		v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
		v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
		v8 = m.T0[v3].(func(*base.Module, int32, int32) int32)(m, v6, v7)
		mBase = m.M
		v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
		*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v9 | base.B2i32(v8 == int32(0))
	}
	return
}
func F_log(m *base.Module, l0 float64) float64 {
	mBase := m.M
	_ = mBase
	var v11 int64
	_ = v11
	var v17 float64
	_ = v17
	var v19 float64
	_ = v19
	var v21 float64
	_ = v21
	var v23 int32
	_ = v23
	var v24 float64
	_ = v24
	var v25 float64
	_ = v25
	var v26 float64
	_ = v26
	var v27 float64
	_ = v27
	var v28 float64
	_ = v28
	var v30 float64
	_ = v30
	var v33 float64
	_ = v33
	var v36 float64
	_ = v36
	var v39 float64
	_ = v39
	var v45 float64
	_ = v45
	var v48 float64
	_ = v48
	var v51 float64
	_ = v51
	var v57 float64
	_ = v57
	var v60 float64
	_ = v60
	var v63 float64
	_ = v63
	var v80 int32
	_ = v80
	var v96 int32
	_ = v96
	var v100 float64
	_ = v100
	var v108 int64
	_ = v108
	var v110 int64
	_ = v110
	var v114 float64
	_ = v114
	var v115 int32
	_ = v115
	var v116 float64
	_ = v116
	var v124 int32
	_ = v124
	var v127 float64
	_ = v127
	var v128 float64
	_ = v128
	var v131 float64
	_ = v131
	var v138 float64
	_ = v138
	var v142 float64
	_ = v142
	var v144 float64
	_ = v144
	var v145 float64
	_ = v145
	var v146 float64
	_ = v146
	var v149 float64
	_ = v149
	var v152 float64
	_ = v152
	var v156 float64
	_ = v156
	var v159 float64
	_ = v159
	var v164 float64
	_ = v164
	var v167 float64
	_ = v167
	var v175 float64
	_ = v175
	v11 = base.I64_reinterpret_f64(l0)
	if base.Ui64(int64(854320534781951)) < base.Ui64(v11+int64(-4606619468846596096)) {
		v80 = base.I32_wrap_i64(int64(base.Ui64(v11) >> (uint(int64(48)) % 64)))
		if base.Ui32(int32(-32737)) < base.Ui32(v80+int32(-32752)) {
			v108 = v11
			v110 = v108 + int64(-4604367669032910848)
			v114 = base.F64_convert_i32_s(base.I32_wrap_i64(v110 >> (uint(int64(52)) % 64)))
			v115 = int32(0)
			v116 = *(*float64)(unsafe.Add(mBase, _c_F_log[0]))
			v124 = base.I32_wrap_i64(int64(base.Ui64(v110)>>(uint(int64(45))%64))) & int32(127) << (uint(int32(4)) % 32)
			v127 = *(*float64)(unsafe.Add(mBase, uint32(v124)+uint32(_c_F_log[1])))
			v128 = base.F64_add(base.F64_mul(v114, v116), v127)
			v131 = *(*float64)(unsafe.Add(mBase, uint32(v124)+uint32(_c_F_log[2])))
			v138 = *(*float64)(unsafe.Add(mBase, uint32(v124)+uint32(_c_F_log[3])))
			v142 = *(*float64)(unsafe.Add(mBase, uint32(v124)+uint32(_c_F_log[4])))
			v144 = base.F64_mul(v131, base.F64_sub(base.F64_sub(base.F64_reinterpret_i64(v108-v110&int64(-4503599627370496)), v138), v142))
			v145 = base.F64_add(v128, v144)
			v146 = base.F64_mul(v144, v144)
			v149 = *(*float64)(unsafe.Add(mBase, _c_F_log[5]))
			v152 = *(*float64)(unsafe.Add(mBase, _c_F_log[6]))
			v156 = *(*float64)(unsafe.Add(mBase, _c_F_log[7]))
			v159 = *(*float64)(unsafe.Add(mBase, _c_F_log[8]))
			v164 = *(*float64)(unsafe.Add(mBase, _c_F_log[9]))
			v167 = *(*float64)(unsafe.Add(mBase, _c_F_log[10]))
			v175 = base.F64_add(v145, base.F64_add(base.F64_mul(base.F64_mul(v144, v146), base.F64_add(base.F64_mul(v146, base.F64_add(base.F64_mul(v144, v149), v152)), base.F64_add(base.F64_mul(v144, v156), v159))), base.F64_add(base.F64_mul(v146, v164), base.F64_add(base.F64_mul(v114, v167), base.F64_add(v144, base.F64_sub(v128, v145))))))
			return v175
		} else {
			if base.F64_ne(l0, float64(0)) != 0 {
				if v11 == int64(9218868437227405312) {
					v175 = l0
					return v175
				} else {
					if v11 < int64(0) {
						v100 = base.F64_sub(l0, l0)
						return base.F64_div(v100, v100)
					} else {
						v96 = int32(_a_F_log_0)
						if v80&v96 != v96 {
							v108 = base.I64_reinterpret_f64(base.F64_mul(l0, float64(4.503599627370496e+15))) + int64(-234187180623265792)
							v110 = v108 + int64(-4604367669032910848)
							v114 = base.F64_convert_i32_s(base.I32_wrap_i64(v110 >> (uint(int64(52)) % 64)))
							v115 = int32(0)
							v116 = *(*float64)(unsafe.Add(mBase, _c_F_log[0]))
							v124 = base.I32_wrap_i64(int64(base.Ui64(v110)>>(uint(int64(45))%64))) & int32(127) << (uint(int32(4)) % 32)
							v127 = *(*float64)(unsafe.Add(mBase, uint32(v124)+uint32(_c_F_log[1])))
							v128 = base.F64_add(base.F64_mul(v114, v116), v127)
							v131 = *(*float64)(unsafe.Add(mBase, uint32(v124)+uint32(_c_F_log[2])))
							v138 = *(*float64)(unsafe.Add(mBase, uint32(v124)+uint32(_c_F_log[3])))
							v142 = *(*float64)(unsafe.Add(mBase, uint32(v124)+uint32(_c_F_log[4])))
							v144 = base.F64_mul(v131, base.F64_sub(base.F64_sub(base.F64_reinterpret_i64(v108-v110&int64(-4503599627370496)), v138), v142))
							v145 = base.F64_add(v128, v144)
							v146 = base.F64_mul(v144, v144)
							v149 = *(*float64)(unsafe.Add(mBase, _c_F_log[5]))
							v152 = *(*float64)(unsafe.Add(mBase, _c_F_log[6]))
							v156 = *(*float64)(unsafe.Add(mBase, _c_F_log[7]))
							v159 = *(*float64)(unsafe.Add(mBase, _c_F_log[8]))
							v164 = *(*float64)(unsafe.Add(mBase, _c_F_log[9]))
							v167 = *(*float64)(unsafe.Add(mBase, _c_F_log[10]))
							v175 = base.F64_add(v145, base.F64_add(base.F64_mul(base.F64_mul(v144, v146), base.F64_add(base.F64_mul(v146, base.F64_add(base.F64_mul(v144, v149), v152)), base.F64_add(base.F64_mul(v144, v156), v159))), base.F64_add(base.F64_mul(v146, v164), base.F64_add(base.F64_mul(v114, v167), base.F64_add(v144, base.F64_sub(v128, v145))))))
							return v175
						} else {
							v100 = base.F64_sub(l0, l0)
							return base.F64_div(v100, v100)
						}
					}
				}
			} else {
				return math.Float64frombits(uint64(0xfff0000000000000))
			}
		}
	} else {
		v17 = base.F64_add(l0, float64(-1))
		v19 = base.F64_mul(v17, float64(1.34217728e+08))
		v21 = base.F64_sub(base.F64_add(v17, v19), v19)
		v23 = int32(0)
		v24 = *(*float64)(unsafe.Add(mBase, _c_F_log[11]))
		v25 = base.F64_mul(base.F64_mul(v21, v21), v24)
		v26 = base.F64_add(v17, v25)
		v27 = base.F64_mul(v17, v17)
		v28 = base.F64_mul(v17, v27)
		v30 = *(*float64)(unsafe.Add(mBase, _c_F_log[12]))
		v33 = *(*float64)(unsafe.Add(mBase, _c_F_log[13]))
		v36 = *(*float64)(unsafe.Add(mBase, _c_F_log[14]))
		v39 = *(*float64)(unsafe.Add(mBase, _c_F_log[15]))
		v45 = *(*float64)(unsafe.Add(mBase, _c_F_log[16]))
		v48 = *(*float64)(unsafe.Add(mBase, _c_F_log[17]))
		v51 = *(*float64)(unsafe.Add(mBase, _c_F_log[18]))
		v57 = *(*float64)(unsafe.Add(mBase, _c_F_log[19]))
		v60 = *(*float64)(unsafe.Add(mBase, _c_F_log[20]))
		v63 = *(*float64)(unsafe.Add(mBase, _c_F_log[21]))
		return base.F64_add(v26, base.F64_add(base.F64_mul(v28, base.F64_add(base.F64_mul(v28, base.F64_add(base.F64_mul(v28, base.F64_add(base.F64_mul(v28, v30), base.F64_add(base.F64_mul(v27, v33), base.F64_add(base.F64_mul(v17, v36), v39)))), base.F64_add(base.F64_mul(v27, v45), base.F64_add(base.F64_mul(v17, v48), v51)))), base.F64_add(base.F64_mul(v27, v57), base.F64_add(base.F64_mul(v17, v60), v63)))), base.F64_add(base.F64_mul(base.F64_mul(base.F64_sub(v17, v21), v24), base.F64_add(v17, v21)), base.F64_add(v25, base.F64_sub(v17, v26)))))
	}
}
func F_log10(m *base.Module, l0 float64) float64 {
	var v12 int64
	_ = v12
	var v29 int32
	_ = v29
	var v31 int64
	_ = v31
	var v41 int64
	_ = v41
	var v46 int64
	_ = v46
	var v47 int32
	_ = v47
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v55 float64
	_ = v55
	var v57 float64
	_ = v57
	var v70 float64
	_ = v70
	var v73 float64
	_ = v73
	var v78 float64
	_ = v78
	var v79 float64
	_ = v79
	var v80 float64
	_ = v80
	var v81 float64
	_ = v81
	var v86 float64
	_ = v86
	var v87 float64
	_ = v87
	var v88 float64
	_ = v88
	var v113 float64
	_ = v113
	var v125 float64
	_ = v125
	v12 = base.I64_reinterpret_f64(l0)
	if int64(4503599627370495) < v12 {
		if base.Ui64(int64(9218868437227405311)) < base.Ui64(v12) {
			v125 = l0
			return v125
		} else {
			v29 = int32(-1023)
			v31 = int64(base.Ui64(v12) >> (uint(int64(32)) % 64))
			if v31 == int64(1072693248) {
				if base.I32_wrap_i64(v12) != 0 {
					v46 = v12
					v47 = v29
					v49 = int32(1072693248)
					v51 = v49 + int32(614242)
					v55 = base.F64_convert_i32_s(v47 + int32(base.Ui32(v51)>>(uint(int32(20))%32)))
					v57 = base.F64_mul(v55, float64(0.30102999566361177))
					v70 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v51&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v46&int64(4294967295)), float64(-1))
					v73 = base.F64_mul(v70, base.F64_mul(v70, float64(0.5)))
					v78 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v70, v73)) & int64(-4294967296))
					v79 = float64(0.4342944818781689)
					v80 = base.F64_mul(v78, v79)
					v81 = base.F64_add(v57, v80)
					v86 = base.F64_div(v70, base.F64_add(v70, float64(2)))
					v87 = base.F64_mul(v86, v86)
					v88 = base.F64_mul(v87, v87)
					v113 = base.F64_add(base.F64_mul(v86, base.F64_add(v73, base.F64_add(base.F64_mul(v88, base.F64_add(base.F64_mul(v88, base.F64_add(base.F64_mul(v88, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v87, base.F64_add(base.F64_mul(v88, base.F64_add(base.F64_mul(v88, base.F64_add(base.F64_mul(v88, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v70, v78), v73))
					v125 = base.F64_add(v81, base.F64_add(base.F64_add(v80, base.F64_sub(v57, v81)), base.F64_add(base.F64_mul(v113, v79), base.F64_add(base.F64_mul(v55, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v113, v78), float64(2.5082946711645275e-11))))))
					return v125
				} else {
					return float64(0)
				}
			} else {
				v46 = v12
				v47 = v29
				v49 = base.I32_wrap_i64(v31)
				v51 = v49 + int32(614242)
				v55 = base.F64_convert_i32_s(v47 + int32(base.Ui32(v51)>>(uint(int32(20))%32)))
				v57 = base.F64_mul(v55, float64(0.30102999566361177))
				v70 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v51&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v46&int64(4294967295)), float64(-1))
				v73 = base.F64_mul(v70, base.F64_mul(v70, float64(0.5)))
				v78 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v70, v73)) & int64(-4294967296))
				v79 = float64(0.4342944818781689)
				v80 = base.F64_mul(v78, v79)
				v81 = base.F64_add(v57, v80)
				v86 = base.F64_div(v70, base.F64_add(v70, float64(2)))
				v87 = base.F64_mul(v86, v86)
				v88 = base.F64_mul(v87, v87)
				v113 = base.F64_add(base.F64_mul(v86, base.F64_add(v73, base.F64_add(base.F64_mul(v88, base.F64_add(base.F64_mul(v88, base.F64_add(base.F64_mul(v88, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v87, base.F64_add(base.F64_mul(v88, base.F64_add(base.F64_mul(v88, base.F64_add(base.F64_mul(v88, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v70, v78), v73))
				v125 = base.F64_add(v81, base.F64_add(base.F64_add(v80, base.F64_sub(v57, v81)), base.F64_add(base.F64_mul(v113, v79), base.F64_add(base.F64_mul(v55, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v113, v78), float64(2.5082946711645275e-11))))))
				return v125
			}
		}
	} else {
		if base.F64_ne(l0, float64(0)) != 0 {
			if int64(-1) < v12 {
				v41 = base.I64_reinterpret_f64(base.F64_mul(l0, float64(1.8014398509481984e+16)))
				v46 = v41
				v47 = int32(-1077)
				v49 = base.I32_wrap_i64(int64(base.Ui64(v41) >> (uint(int64(32)) % 64)))
				v51 = v49 + int32(614242)
				v55 = base.F64_convert_i32_s(v47 + int32(base.Ui32(v51)>>(uint(int32(20))%32)))
				v57 = base.F64_mul(v55, float64(0.30102999566361177))
				v70 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v51&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v46&int64(4294967295)), float64(-1))
				v73 = base.F64_mul(v70, base.F64_mul(v70, float64(0.5)))
				v78 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v70, v73)) & int64(-4294967296))
				v79 = float64(0.4342944818781689)
				v80 = base.F64_mul(v78, v79)
				v81 = base.F64_add(v57, v80)
				v86 = base.F64_div(v70, base.F64_add(v70, float64(2)))
				v87 = base.F64_mul(v86, v86)
				v88 = base.F64_mul(v87, v87)
				v113 = base.F64_add(base.F64_mul(v86, base.F64_add(v73, base.F64_add(base.F64_mul(v88, base.F64_add(base.F64_mul(v88, base.F64_add(base.F64_mul(v88, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v87, base.F64_add(base.F64_mul(v88, base.F64_add(base.F64_mul(v88, base.F64_add(base.F64_mul(v88, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v70, v78), v73))
				v125 = base.F64_add(v81, base.F64_add(base.F64_add(v80, base.F64_sub(v57, v81)), base.F64_add(base.F64_mul(v113, v79), base.F64_add(base.F64_mul(v55, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v113, v78), float64(2.5082946711645275e-11))))))
				return v125
			} else {
				return base.F64_div(base.F64_sub(l0, l0), float64(0))
			}
		} else {
			return base.F64_div(float64(-1), base.F64_mul(l0, l0))
		}
	}
}
func F_logf(m *base.Module, l0 float32) float32 {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v12 int32
	_ = v12
	var v24 float32
	_ = v24
	var v32 int32
	_ = v32
	var v34 int32
	_ = v34
	var v35 float64
	_ = v35
	var v37 int32
	_ = v37
	var v46 int32
	_ = v46
	var v49 float64
	_ = v49
	var v52 float64
	_ = v52
	var v53 float64
	_ = v53
	var v56 float64
	_ = v56
	var v59 float64
	_ = v59
	var v67 float64
	_ = v67
	var v71 float64
	_ = v71
	var v76 float32
	_ = v76
	v6 = base.I32_reinterpret_f32(l0)
	if base.Ui32(int32(-2130706433)) < base.Ui32(v6+int32(-2139095040)) {
		v32 = v6
		v34 = int32(0)
		v35 = *(*float64)(unsafe.Add(mBase, _c_F_logf[0]))
		v37 = v32 + int32(-1060306944)
		v46 = int32(base.Ui32(v37)>>(uint(int32(15))%32)) & int32(240)
		v49 = *(*float64)(unsafe.Add(mBase, uint32(v46)+uint32(_c_F_logf[1])))
		v52 = base.F64_add(base.F64_mul(base.F64_promote_f32(base.F32_reinterpret_i32(v32-v37&int32(-8388608))), v49), float64(-1))
		v53 = base.F64_mul(v52, v52)
		v56 = *(*float64)(unsafe.Add(mBase, _c_F_logf[2]))
		v59 = *(*float64)(unsafe.Add(mBase, _c_F_logf[3]))
		v67 = *(*float64)(unsafe.Add(mBase, _c_F_logf[4]))
		v71 = *(*float64)(unsafe.Add(mBase, uint32(v46)+uint32(_c_F_logf[5])))
		v76 = base.F32_demote_f64(base.F64_add(base.F64_mul(base.F64_add(base.F64_mul(v35, v53), base.F64_add(base.F64_mul(v56, v52), v59)), v53), base.F64_add(base.F64_add(base.F64_mul(base.F64_convert_i32_s(v37>>(uint(int32(23))%32)), v67), v71), v52)))
		return v76
	} else {
		v12 = v6 << (uint(int32(1)) % 32)
		if v12 != 0 {
			if v6 == int32(2139095040) {
				v76 = l0
				return v76
			} else {
				if v6 < int32(0) {
					v24 = base.F32_sub(l0, l0)
					return base.F32_div(v24, v24)
				} else {
					if base.Ui32(v12) < base.Ui32(int32(-16777216)) {
						v32 = base.I32_reinterpret_f32(base.F32_mul(l0, float32(8.388608e+06))) + int32(-192937984)
						v34 = int32(0)
						v35 = *(*float64)(unsafe.Add(mBase, _c_F_logf[0]))
						v37 = v32 + int32(-1060306944)
						v46 = int32(base.Ui32(v37)>>(uint(int32(15))%32)) & int32(240)
						v49 = *(*float64)(unsafe.Add(mBase, uint32(v46)+uint32(_c_F_logf[1])))
						v52 = base.F64_add(base.F64_mul(base.F64_promote_f32(base.F32_reinterpret_i32(v32-v37&int32(-8388608))), v49), float64(-1))
						v53 = base.F64_mul(v52, v52)
						v56 = *(*float64)(unsafe.Add(mBase, _c_F_logf[2]))
						v59 = *(*float64)(unsafe.Add(mBase, _c_F_logf[3]))
						v67 = *(*float64)(unsafe.Add(mBase, _c_F_logf[4]))
						v71 = *(*float64)(unsafe.Add(mBase, uint32(v46)+uint32(_c_F_logf[5])))
						v76 = base.F32_demote_f64(base.F64_add(base.F64_mul(base.F64_add(base.F64_mul(v35, v53), base.F64_add(base.F64_mul(v56, v52), v59)), v53), base.F64_add(base.F64_add(base.F64_mul(base.F64_convert_i32_s(v37>>(uint(int32(23))%32)), v67), v71), v52)))
						return v76
					} else {
						v24 = base.F32_sub(l0, l0)
						return base.F32_div(v24, v24)
					}
				}
			}
		} else {
			return math.Float32frombits(uint32(0xff800000))
		}
	}
}
