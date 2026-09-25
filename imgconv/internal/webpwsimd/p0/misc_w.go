//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_WebPMultRow_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v13 int32
	_ = v13
	var v16 int32
	_ = v16
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v40 int32
	_ = v40
	var v42 int32
	_ = v42
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v56 int32
	_ = v56
	if l2 < int32(1) {
	} else {
		if l3 == int32(0) {
			v33 = l0
			v34 = l1
			v35 = l2
			for {
				v37 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v34))))
				if v37 == int32(255) {
				} else {
					if v37 != 0 {
						v42 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v33))))
						v49 = int32(base.Ui32(v37*v42*int32(_a_F_WebPMultRow_C_0)+int32(8388608)) >> (uint(int32(24)) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v33))) = uint8(v49)
					} else {
						v40 = int32(0)
						*(*uint8)(unsafe.Add(mBase, uint32(v33))) = uint8(v40)
					}
				}
				v51 = int32(1)
				v56 = v35 + int32(-1)
				if v56 != 0 {
					v33 = v33 + v51
					v34 = v34 + v51
					v35 = v56
					continue
				} else {
					break
				}
				break
			}
		} else {
			v9 = l0
			v10 = l1
			v11 = l2
			for {
				v13 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v10))))
				if v13 == int32(255) {
				} else {
					if v13 != 0 {
						v19 = base.I32_div_u_s(int32(-16777216), v13)
						v20 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v9))))
						v25 = int32(base.Ui32(v19*v20+int32(8388608)) >> (uint(int32(24)) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v9))) = uint8(v25)
					} else {
						v16 = int32(0)
						*(*uint8)(unsafe.Add(mBase, uint32(v9))) = uint8(v16)
					}
				}
				v27 = int32(1)
				v32 = v11 + int32(-1)
				if v32 != 0 {
					v9 = v9 + v27
					v10 = v10 + v27
					v11 = v32
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
func F___wasm_call_ctors(m *base.Module) {
	return
}
func F_wasmCPUInfo_1(m *base.Module, l0 int32) int32 {
	return base.B2i32(base.Ui32(l0) < base.Ui32(int32(4)))
}
func F_wasmCPUInfo_2(m *base.Module, l0 int32) int32 {
	return base.B2i32(base.Ui32(l0) < base.Ui32(int32(4)))
}
func F_wrapper_cmp(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	var v4 int32
	_ = v4
	v4 = m.T0[l2].(func(*base.Module, int32, int32) int32)(m, l0, l1)
	return v4
}
