//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package base

import (
	"math"
	"math/bits"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

type Module struct {
	Memory      []byte
	MaxMem      uint64
	M           unsafe.Pointer
	T0          []any
	G0          int32
	G1          int32
	G2          int32
	G3          int32
	G4          int32
	G5          int32
	G6          int32
	G7          int32
	G8          int32
	G9          int32
	G10         int32
	G11         int32
	G12         int32
	G13         int32
	G14         int32
	G15         int32
	G16         int32
	G17         int32
	G18         int32
	G19         int32
	G20         int32
	G21         int32
	G22         int32
	G23         int32
	G24         int32
	G25         int32
	G26         int32
	G27         int32
	G28         int32
	G29         int32
	G30         int32
	G31         int32
	G32         int32
	G33         int32
	G34         int32
	G35         int32
	G36         int32
	G37         int32
	G38         int32
	G39         int32
	G40         int32
	G41         int32
	G42         int32
	G43         int32
	G44         int32
	G45         int32
	G46         int32
	G47         int32
	G48         int32
	G49         int32
	G50         int32
	G51         int32
	G52         int32
	G53         int32
	G54         int32
	G55         int32
	G56         int32
	G57         int32
	G58         int32
	G59         int32
	G60         int32
	G61         int32
	G62         int32
	G63         int32
	G64         int32
	G65         int32
	G66         int32
	G67         int32
	G68         int32
	G69         int32
	G70         int32
	G71         int32
	G72         int32
	G73         int32
	G74         int32
	G75         int32
	G76         int32
	G77         int32
	G78         int32
	G79         int32
	G80         int32
	G81         int32
	G82         int32
	G83         int32
	G84         int32
	G85         int32
	G86         int32
	G87         int32
	G88         int32
	G89         int32
	G90         int32
	G91         int32
	G92         int32
	G93         int32
	G94         int32
	G95         int32
	G96         int32
	G97         int32
	G98         int32
	G99         int32
	G100        int32
	G101        int32
	G102        int32
	G103        int32
	G104        int32
	G105        int32
	G106        int32
	G107        int32
	G108        int32
	G109        int32
	G110        int32
	G111        int32
	G112        int32
	G113        int32
	G114        int32
	G115        int32
	G116        int32
	G117        int32
	G118        int32
	G119        int32
	G120        int32
	G121        int32
	G122        int32
	G123        int32
	G124        int32
	G125        int32
	G126        int32
	G127        int32
	G128        int32
	G129        int32
	G130        int32
	MemMu       *sync.Mutex
	MemSize     *atomic.Uint64
	DataEnd     uint32
	MemShared   bool
	Threads     *ThreadPool
	ThreadStart func(*Module, int32, int32)
}

func I32(x int32) int32 { return x }

func I64(x int64) int64 { return x }

// ui32 / ui64 reinterpret a signed integer as its unsigned bit
// equivalent at runtime. Used for the operands of wasm unsigned
// comparisons (i32.lt_u etc.) — emitting `uint32(int32(-N))` directly
// fails Go's compile-time constant rule because the negative typed
// constant isn't representable in uint32; routing through these
// function-call boundaries forces runtime conversion.
func Ui32(x int32) uint32 { return uint32(x) }

func Ui64(x int64) uint64 { return uint64(x) }

// b2i32 materialises a wasm comparison result — an i32 that is 0 or 1 — from
// the Go bool the comparison expression evaluates to.
//
// It exists as a named helper rather than an inline `func() int32 { ... }()`
// because the gcasm backend requires every direct call left in the compiled
// output to be either a package-local FnN or something the Go inliner removed.
// A func literal is normally inlined at its call site, but the inliner gives up
// once the ENCLOSING function grows past its budget — and a single wasm function
// can translate to tens of thousands of lines of Go, as an interpreter's
// bytecode dispatch loop does. The literal is then outlined into a real closure
// symbol (FnN.funcA.funcB), which reaches the assembler as a direct call gcasm
// cannot marshal. A named helper this small is always inlined, and if it ever
// were not, it would fail loudly at its own symbol rather than as a nested
// closure.
func B2i32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func F32(x float32) float32 { runtime.KeepAlive(&x); return x }

func F64(x float64) float64 { runtime.KeepAlive(&x); return x }

//go:noinline
func Wasm_trap_div_zero() { panic("wasm: integer divide by zero") }

//go:noinline
func Wasm_trap_int_overflow() { panic("wasm: integer overflow") }

//go:noinline
func Wasm_trap_invalid_conv() { panic("wasm: invalid conversion to integer") }

//go:noinline
func Wasm_trap_unreachable() { panic("wasm: unreachable") }

//go:noinline
func Wasm_trap_memfill_oob() { panic("wasm: memory.fill out of bounds") }

//go:noinline
func Wasm_trap_memcopy_oob() { panic("wasm: memory.copy out of bounds") }

func I32_div_s(x, y int32) int32 {
	if y == -1 && x == math.MinInt32 {
		Wasm_trap_int_overflow()
	}
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x / y
}

func I64_div_s(x, y int64) int64 {
	if y == -1 && x == math.MinInt64 {
		Wasm_trap_int_overflow()
	}
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x / y
}

func I32_div_u(x, y uint32) uint32 {
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x / y
}

func I64_div_u(x, y uint64) uint64 {
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x / y
}

func I32_rem_u(x, y uint32) uint32 {
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x % y
}

func I64_rem_u(x, y uint64) uint64 {
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x % y
}

func I32_rotl(x, y int32) int32 { return int32(bits.RotateLeft32(uint32(x), int(y&31))) }

func I64_rotl(x, y int64) int64 { return int64(bits.RotateLeft64(uint64(x), int(y&63))) }

func F32_abs(x float32) float32 {
	return math.Float32frombits(math.Float32bits(x) &^ (1 << 31))
}

func F64_abs(x float64) float64 {
	return math.Float64frombits(math.Float64bits(x) &^ (1 << 63))
}

func F32_neg(x float32) float32 {
	return math.Float32frombits(math.Float32bits(x) ^ (1 << 31))
}

func F64_neg(x float64) float64 {
	return math.Float64frombits(math.Float64bits(x) ^ (1 << 63))
}

func I32_trunc_f32_s(x float32) int32 {
	if x != x {
		Wasm_trap_invalid_conv()
	}

	if !(x > -2147483904.0 && x < 2147483648.0) {
		Wasm_trap_int_overflow()
	}
	return int32(x)
}

func I32_trunc_f32_u(x float32) int32 {
	if x != x {
		Wasm_trap_invalid_conv()
	}
	if !(x > -1.0 && x < 4294967296.0) {
		Wasm_trap_int_overflow()
	}
	return int32(uint32(x))
}

func I32_trunc_f64_s(x float64) int32 {
	if x != x {
		Wasm_trap_invalid_conv()
	}

	if !(x > -2147483649.0 && x < 2147483648.0) {
		Wasm_trap_int_overflow()
	}
	return int32(x)
}

func I32_trunc_f64_u(x float64) int32 {
	if x != x {
		Wasm_trap_invalid_conv()
	}
	if !(x > -1.0 && x < 4294967296.0) {
		Wasm_trap_int_overflow()
	}
	return int32(uint32(x))
}

func I64_trunc_f64_u(x float64) int64 {
	if x != x {
		Wasm_trap_invalid_conv()
	}
	if !(x > -1.0 && x < 18446744073709551616.0) {
		Wasm_trap_int_overflow()
	}
	return int64(uint64(x))
}

// memorySize returns the current size of m.memory in wasm pages (each
// page is 64 KiB).
func MemorySize(m *Module) int32 {
	return int32(m.MemSize.Load() >> 16)
}

// wasmMemHardCap is the implementation limit on linear-memory size:
// 65534 pages, two short of wasm32's architectural 65536. Growth past
// it fails with -1 like any resource limit (the JS API allows an
// engine to refuse any grow). Keeping memSize strictly below 2^32
// minus a 128 KiB margin is what makes the coalesced SIMD bounds check
// (simd_v128_load_rng) exact: a group whose unwrapped address range
// reaches past memSize can then never be a group whose members all
// individually landed in bounds via u32 wraparound.
//
// A function rather than a const because the helper extractor carries
// only function declarations into the output (it must stay in sync
// with codegen's wasmMemHardCapBytes).
func WasmMemHardCap() uint64 { return (1 << 32) - (1 << 17) }

// memoryGrow grows m.memory by n wasm pages (64 KiB each). Returns the
// previous page count, or -1 if the new size would exceed maxMem or
// wasmMemHardCap. n may be 0, which simply returns the current size.
//
// len(m.memory) must always equal the exact wasm memory size (memory.size
// and every bounds check depend on it), but the backing array is grown
// GEOMETRICALLY: a sequence of small memory.grow calls — which a C++ heap
// does constantly during start-up — would otherwise reallocate and recopy
// the whole linear memory on every page, i.e. O(n^2) total copying. Spare
// capacity makes the common grow a zero-copy reslice and amortizes the
// reallocations to O(n).
func MemoryGrow(m *Module, n int32) int32 {

	m.MemMu.Lock()
	defer m.MemMu.Unlock()
	cur := m.MemSize.Load()
	prev := int32(cur >> 16)
	if n == 0 {
		return prev
	}
	if n < 0 {
		return -1
	}
	want := cur + uint64(n)*65536
	if m.MaxMem != 0 && want > m.MaxMem {
		return -1
	}
	if want > WasmMemHardCap() {
		return -1
	}
	if m.MemShared {

		if want > uint64(len(m.Memory)) {
			return -1
		}
		m.MemSize.Store(want)
		return prev
	}
	if want <= uint64(cap(m.Memory)) {

		m.Memory = m.Memory[:want]
		m.MemSize.Store(want)
		return prev
	}

	newCap := uint64(cap(m.Memory)) * 2
	if newCap < want {
		newCap = want
	}
	if m.MaxMem != 0 && newCap > m.MaxMem {
		newCap = m.MaxMem
	}
	if newCap > WasmMemHardCap() {
		newCap = WasmMemHardCap()
	}
	grown := make([]byte, want, newCap)
	copy(grown, m.Memory)
	m.Memory = grown
	m.MemSize.Store(want)

	m.M = unsafe.Pointer(unsafe.SliceData(m.Memory))
	return prev
}

// accessMemory runs f with the module's current linear memory while
// holding the same lock memoryGrow takes to mutate the memory slice
// header or relocate its backing array. It is the ONE safe way to
// touch linear memory from OUTSIDE the module's execution goroutine —
// e.g. a watchdog goroutine raising CPython's eval-breaker bit while
// an evaluation is running. For the duration of f the memory can
// neither be resliced nor relocated, so f's writes land in the array
// the guest observes; a grow that raced in just before blocks until f
// returns and then copies f's writes forward with the rest of the
// contents. Determinism notes for callers:
//
//   - f MUST NOT call back into the module or into memoryGrow — that
//     would self-deadlock.
//   - f should be short: a running guest blocks inside memory.grow
//     until f returns (ordinary guest loads/stores do not block).
//   - Bytes the guest reads or writes concurrently with f (that is
//     the point of an eval-breaker-style flag) are exchanged with
//     plain single-word accesses; keep such shared words
//     word-aligned and word-sized.
func AccessMemory(m *Module, f func(mem []byte)) {
	m.MemMu.Lock()
	defer m.MemMu.Unlock()
	f(m.Memory)
}

func I32_div_u_s(x, y int32) int32 { return int32(I32_div_u(uint32(x), uint32(y))) }
func I32_rem_u_s(x, y int32) int32 { return int32(I32_rem_u(uint32(x), uint32(y))) }
func I64_div_u_s(x, y int64) int64 { return int64(I64_div_u(uint64(x), uint64(y))) }
func I64_rem_u_s(x, y int64) int64 { return int64(I64_rem_u(uint64(x), uint64(y))) }

// The explicit same-type conversions are NOT redundant: they are
// rounding points. Once these helpers inline, gc is free to fuse a
// multiply feeding an add into a single FMA — legal Go, but wasm
// requires every operation individually rounded, and a fused result
// diverges from every wasm runtime (bitwise, and observably in greedy
// sampling). A float conversion forces the intermediate rounding and
// forbids the fusion (spec: Conversions, "rounds to the precision of
// the target type"; the same rule math.FMA documents).
func F32_add(x, y float32) float32 { return float32(x + y) }
func F32_sub(x, y float32) float32 { return float32(x - y) }
func F32_mul(x, y float32) float32 { return float32(x * y) }
func F32_div(x, y float32) float32 { return float32(x / y) }
func F64_add(x, y float64) float64 { return float64(x + y) }
func F64_sub(x, y float64) float64 { return float64(x - y) }
func F64_mul(x, y float64) float64 { return float64(x * y) }
func F64_div(x, y float64) float64 { return float64(x / y) }

func I32_clz(x int32) int32    { return int32(bits.LeadingZeros32(uint32(x))) }
func I32_ctz(x int32) int32    { return int32(bits.TrailingZeros32(uint32(x))) }
func I32_popcnt(x int32) int32 { return int32(bits.OnesCount32(uint32(x))) }

func F32_ceil(x float32) float32 { return float32(math.Ceil(float64(x))) }

func F32_floor(x float32) float32 { return float32(math.Floor(float64(x))) }

func F32_sqrt(x float32) float32 { return float32(math.Sqrt(float64(x))) }

func F32_eq(x, y float32) int32 {
	if x == y {
		return 1
	}
	return 0
}
func F32_ne(x, y float32) int32 {
	if x != y {
		return 1
	}
	return 0
}
func F32_lt(x, y float32) int32 {
	if x < y {
		return 1
	}
	return 0
}
func F32_gt(x, y float32) int32 {
	if x > y {
		return 1
	}
	return 0
}
func F32_le(x, y float32) int32 {
	if x <= y {
		return 1
	}
	return 0
}
func F32_ge(x, y float32) int32 {
	if x >= y {
		return 1
	}
	return 0
}

func F64_eq(x, y float64) int32 {
	if x == y {
		return 1
	}
	return 0
}
func F64_ne(x, y float64) int32 {
	if x != y {
		return 1
	}
	return 0
}
func F64_lt(x, y float64) int32 {
	if x < y {
		return 1
	}
	return 0
}
func F64_gt(x, y float64) int32 {
	if x > y {
		return 1
	}
	return 0
}
func F64_le(x, y float64) int32 {
	if x <= y {
		return 1
	}
	return 0
}
func F64_ge(x, y float64) int32 {
	if x >= y {
		return 1
	}
	return 0
}

func I32_wrap_i64(x int64) int32       { return int32(x) }
func I64_extend_i32_s(x int32) int64   { return int64(x) }
func I64_extend_i32_u(x int32) int64   { return int64(uint32(x)) }
func F32_demote_f64(x float64) float32 { return float32(x) }
func F64_promote_f32(x float32) float64 {

	if math.IsNaN(float64(x)) {

		return float64(x)
	}
	return float64(x)
}

func F32_convert_i32_s(x int32) float32 { return float32(x) }
func F32_convert_i32_u(x int32) float32 { return float32(uint32(x)) }

func F64_convert_i32_s(x int32) float64 { return float64(x) }
func F64_convert_i32_u(x int32) float64 { return float64(uint32(x)) }

func F64_convert_i64_u(x int64) float64 { return float64(uint64(x)) }

func I32_reinterpret_f32(x float32) int32 { return int32(math.Float32bits(x)) }
func I64_reinterpret_f64(x float64) int64 { return int64(math.Float64bits(x)) }
func F32_reinterpret_i32(x int32) float32 { return math.Float32frombits(uint32(x)) }
func F64_reinterpret_i64(x int64) float64 { return math.Float64frombits(uint64(x)) }

func I32_extend8_s(x int32) int32  { return int32(int8(x)) }
func I32_extend16_s(x int32) int32 { return int32(int16(x)) }

func I64_extend32_s(x int64) int64 { return int64(int32(x)) }

func MemoryFill(m *Module, dst int32, val int32, n int32) {
	if n == 0 {
		return
	}
	end := uint64(uint32(dst)) + uint64(uint32(n))
	if end > MemBound(m) {
		Wasm_trap_memfill_oob()
	}
	b := m.Memory[uint32(dst):uint32(end)]
	v := byte(val)

	if v == 0 {
		for k := range b {
			b[k] = 0
		}
		return
	}
	b[0] = v
	for filled := 1; filled < len(b); filled *= 2 {
		copy(b[filled:], b[:filled])
	}
}

func MemoryCopy(m *Module, dst int32, src int32, n int32) {
	if n == 0 {
		return
	}
	srcEnd := uint64(uint32(src)) + uint64(uint32(n))
	dstEnd := uint64(uint32(dst)) + uint64(uint32(n))
	if size := MemBound(m); srcEnd > size || dstEnd > size {
		Wasm_trap_memcopy_oob()
	}
	copy(m.Memory[uint32(dst):uint32(dstEnd)], m.Memory[uint32(src):uint32(srcEnd)])
}

// memBound is the highest address a bounds-checked bulk or atomic access
// may reach. A shared memory's slice spans the whole declared maximum from
// the start, so only memSize says how much of it the guest may touch (and
// reading it atomically keeps growth race-free without a lock). Otherwise
// the slice is the truth: an embedder that maps the whole growable range
// up front and aliases shared segments above the guest-visible size into
// it keeps every byte of the slice addressable, exactly as the unchecked
// load/store paths do.
func MemBound(m *Module) uint64 {
	if m.MemShared {
		return m.MemSize.Load()
	}
	return uint64(len(m.Memory))
}

// forceContendedAtomics makes every atomic helper of m take its LOCKed
// path from now on, as if a wasi thread had been spawned. For an embedder
// that shares part of m's linear memory with OTHER instances (several
// single-threaded modules aliasing one segment, the way processes share
// System V memory): no thread of m ever exists, yet the atomics in the
// shared range race with those instances' goroutines. Host-facing API;
// nothing in the generated code calls it.
func ForceContendedAtomics(m *Module) {
	if m.Threads == nil {
		m.Threads = &ThreadPool{}
	}
	if m.Threads.nextTID.Load() == 0 {
		m.Threads.nextTID.Store(1)
	}
}

//go:noinline
func Wasm_trap_simd_oob() { panic("wasm: v128 memory access out of bounds") }

// gcasmMemProbe anchors the Module field offsets the gcasm memory-op
// splices hardcode. The splices read m.M and m.memSize straight off the
// receiver in generated assembly, and the offsets of those fields
// depend on the module (the import-interface fields between them vary).
// Rather than re-deriving Go's struct layout, gcasm extracts the two
// offsets from THIS function's captured assembly — two loads off R0/AX,
// M first — so they always come from the same compile that produced the
// code being spliced. Never called at run time.
//
//go:noinline
func GcasmMemProbe(m *Module) (unsafe.Pointer, *atomic.Uint64) {
	return m.M, m.MemSize
}

func SimdEA(m *Module, addr int32, offset int32, size uint64) uint64 {

	ea := uint64(uint32(addr)) + uint64(uint32(offset))
	if ea+size > m.MemSize.Load() {
		Wasm_trap_simd_oob()
	}
	return ea
}

// simd_v128_load_rng is the range-checked first load of a coalesced
// group (pass.CoalesceSimdBounds): one trap decision covers every
// access in [addr+rlo, addr+rlo+span), then it loads at its own
// addr+offset. The group's other loads use the unchecked _nc form
// below — which is only sound BECAUSE this ran first; nothing else may
// emit either form.
//
// rlo is SIGNED: the group minimum may lie below this load's own
// address when the group's loads appear out of address order, and when
// a member's u32 address arithmetic wrapped, addr+rlo can go negative
// — a negative start means some member sits just below 2^32 unwrapped,
// which the per-load checks would have trapped (memSize can never
// reach 2^32: memoryGrow stops at wasmMemHardCap), so trapping on
// start < 0 reproduces the original semantics exactly.
//
//go:noinline
func Simd_v128_load_rng(m *Module, addr int32, offset int32, rlo int32, span int32) [2]uint64 {
	start := int64(uint64(uint32(addr))) + int64(rlo)
	if start < 0 || uint64(start)+uint64(uint32(span)) > m.MemSize.Load() {
		Wasm_trap_simd_oob()
	}
	ea := uint64(uint32(addr)) + uint64(uint32(offset))
	p := unsafe.Add(m.M, uintptr(ea))
	return [2]uint64{*(*uint64)(p), *(*uint64)(unsafe.Add(p, 8))}
}

// simd_v128_load_nc is simd_v128_load minus the bounds check; emitted
// only by the bounds-coalescing pass, always behind a covering
// simd_v128_load_rng.
//
//go:noinline
func Simd_v128_load_nc(m *Module, addr int32, offset int32) [2]uint64 {
	ea := uint64(uint32(addr)) + uint64(uint32(offset))
	p := unsafe.Add(m.M, uintptr(ea))
	return [2]uint64{*(*uint64)(p), *(*uint64)(unsafe.Add(p, 8))}
}

//go:noinline
func Simd_v128_load(m *Module, addr int32, offset int32) [2]uint64 {
	ea := SimdEA(m, addr, offset, 16)
	p := unsafe.Add(m.M, uintptr(ea))
	return [2]uint64{*(*uint64)(p), *(*uint64)(unsafe.Add(p, 8))}
}

//go:noinline
func Simd_scalar_i32_shl(v int32, s int32) int32 { return v << (uint(s) % 32) }

//go:noinline
func Simd_scalar_i32_add(a int32, b int32) int32 { return a + b }

//go:noinline
func Simd_v128_f16x4_cvt_store(m *Module, addr int32, offset int32, v [2]uint64) int32 {

	ea := SimdEA(m, addr, offset, 8)
	var out uint64
	for i := 0; i < 4; i++ {
		w := uint32(v[i/2] >> (32 * uint(i) % 64))
		shl1w := w + w
		sign := w & 0x80000000
		var h uint32
		if shl1w > 0xFF000000 {
			h = (sign >> 16) | 0x7E00
		} else {
			bias := shl1w & 0xFF000000
			if bias < 0x71000000 {
				bias = 0x71000000
			}
			f := math.Float32frombits(w&0x7FFFFFFF) * 0x1p+112 * 0x1p-110
			f += math.Float32frombits((bias >> 1) + 0x07800000)
			fbits := math.Float32bits(f)
			h = (sign >> 16) | (fbits>>13)&0x7C00 + fbits&0xFFF
		}
		out |= uint64(uint16(h)) << (16 * uint(i))
	}
	*(*uint64)(unsafe.Add(m.M, uintptr(ea))) = out
	return 0
}

//go:noinline
func Simd_v128_store(m *Module, addr int32, offset int32, v [2]uint64) int32 {
	ea := SimdEA(m, addr, offset, 16)
	p := unsafe.Add(m.M, uintptr(ea))
	*(*uint64)(p) = v[0]
	*(*uint64)(unsafe.Add(p, 8)) = v[1]
	return 0
}

//go:noinline
func Simd_v128_load8x8_u(m *Module, addr int32, offset int32) [2]uint64 {
	ea := SimdEA(m, addr, offset, 8)
	p := unsafe.Add(m.M, uintptr(ea))
	var out [2]uint64
	for i := 0; i < 8; i++ {
		x := *(*uint8)(unsafe.Add(p, 1*i))
		out[i*16/64] |= uint64(uint16(x)) << (16 * uint(i) % 64)
	}
	return out
}

// f16BitsToF32Bits is the IEEE binary16 -> binary32 conversion,
// bit-exact including subnormals, infinities and NaN payloads.
func F16BitsToF32Bits(h uint16) uint32 {
	sign := uint32(h>>15) << 31
	exp := uint32(h>>10) & 0x1F
	man := uint32(h) & 0x3FF
	switch exp {
	case 0:
		if man == 0 {
			return sign
		}
		e := uint32(113)
		for man&0x400 == 0 {
			man <<= 1
			e--
		}
		return sign | e<<23 | (man&0x3FF)<<13
	case 0x1F:
		return sign | 0xFF<<23 | man<<13
	}
	return sign | (exp+112)<<23 | man<<13
}

// simd_f16x4_cvt converts four f16 values (widened to the low 16 bits
// of each i32x4 lane, the v128_load16x4_u result shape) to f32 lanes.
// Emitted only after the transpiler verified the module's conversion
// table is the IEEE map, so this computed conversion is bit-identical
// to the table reads it replaces (and to XTN+FCVTL on arm64).
//
//go:noinline
func Simd_f16x4_cvt(v [2]uint64) [2]uint64 {
	var out [2]uint64
	for i := 0; i < 4; i++ {
		bits := uint16(v[i/2] >> (32 * uint(i) % 64))
		out[i/2] |= uint64(F16BitsToF32Bits(bits)) << (32 * uint(i) % 64)
	}
	return out
}

//go:noinline
func Simd_v128_load8_splat(m *Module, addr int32, offset int32) [2]uint64 {
	ea := SimdEA(m, addr, offset, 1)
	x := *(*uint8)(unsafe.Add(m.M, uintptr(ea)))
	var out [2]uint64
	for i := 0; i < 16; i++ {
		out[i*8/64] |= uint64(x) << (8 * uint(i) % 64)
	}
	return out
}

//go:noinline
func Simd_v128_load16_splat(m *Module, addr int32, offset int32) [2]uint64 {
	ea := SimdEA(m, addr, offset, 2)
	x := *(*uint16)(unsafe.Add(m.M, uintptr(ea)))
	var out [2]uint64
	for i := 0; i < 8; i++ {
		out[i*16/64] |= uint64(x) << (16 * uint(i) % 64)
	}
	return out
}

//go:noinline
func Simd_v128_load32_splat(m *Module, addr int32, offset int32) [2]uint64 {
	ea := SimdEA(m, addr, offset, 4)
	x := *(*uint32)(unsafe.Add(m.M, uintptr(ea)))
	var out [2]uint64
	for i := 0; i < 4; i++ {
		out[i*32/64] |= uint64(x) << (32 * uint(i) % 64)
	}
	return out
}

//go:noinline
func Simd_v128_load64_splat(m *Module, addr int32, offset int32) [2]uint64 {
	ea := SimdEA(m, addr, offset, 8)
	x := *(*uint64)(unsafe.Add(m.M, uintptr(ea)))
	var out [2]uint64
	for i := 0; i < 2; i++ {
		out[i*64/64] |= uint64(x) << (64 * uint(i) % 64)
	}
	return out
}

//go:noinline
func Simd_v128_load32_zero(m *Module, addr int32, offset int32) [2]uint64 {
	ea := SimdEA(m, addr, offset, 4)
	return [2]uint64{uint64(*(*uint32)(unsafe.Add(m.M, uintptr(ea)))), 0}
}

//go:noinline
func Simd_v128_load64_zero(m *Module, addr int32, offset int32) [2]uint64 {
	ea := SimdEA(m, addr, offset, 8)
	return [2]uint64{*(*uint64)(unsafe.Add(m.M, uintptr(ea))), 0}
}

//go:noinline
func Simd_v128_load8_lane(m *Module, addr int32, offset int32, lane int32, v [2]uint64) [2]uint64 {
	ea := SimdEA(m, addr, offset, 1)
	x := *(*uint8)(unsafe.Add(m.M, uintptr(ea)))
	sh := 8 * uint(lane) % 64
	i := int(lane) * 8 / 64
	v[i] = v[i]&^(uint64(uint8(^uint8(0)))<<sh) | uint64(x)<<sh
	return v
}

//go:noinline
func Simd_v128_load16_lane(m *Module, addr int32, offset int32, lane int32, v [2]uint64) [2]uint64 {
	ea := SimdEA(m, addr, offset, 2)
	x := *(*uint16)(unsafe.Add(m.M, uintptr(ea)))
	sh := 16 * uint(lane) % 64
	i := int(lane) * 16 / 64
	v[i] = v[i]&^(uint64(uint16(^uint16(0)))<<sh) | uint64(x)<<sh
	return v
}

//go:noinline
func Simd_v128_load32_lane(m *Module, addr int32, offset int32, lane int32, v [2]uint64) [2]uint64 {
	ea := SimdEA(m, addr, offset, 4)
	x := *(*uint32)(unsafe.Add(m.M, uintptr(ea)))
	sh := 32 * uint(lane) % 64
	i := int(lane) * 32 / 64
	v[i] = v[i]&^(uint64(uint32(^uint32(0)))<<sh) | uint64(x)<<sh
	return v
}

//go:noinline
func Simd_v128_store8_lane(m *Module, addr int32, offset int32, lane int32, v [2]uint64) int32 {
	ea := SimdEA(m, addr, offset, 1)
	x := uint8(v[int(lane)*8/64] >> (8 * uint(lane) % 64))
	*(*uint8)(unsafe.Add(m.M, uintptr(ea))) = x
	return 0
}

//go:noinline
func Simd_v128_store16_lane(m *Module, addr int32, offset int32, lane int32, v [2]uint64) int32 {
	ea := SimdEA(m, addr, offset, 2)
	x := uint16(v[int(lane)*16/64] >> (16 * uint(lane) % 64))
	*(*uint16)(unsafe.Add(m.M, uintptr(ea))) = x
	return 0
}

//go:noinline
func Simd_v128_store32_lane(m *Module, addr int32, offset int32, lane int32, v [2]uint64) int32 {
	ea := SimdEA(m, addr, offset, 4)
	x := uint32(v[int(lane)*32/64] >> (32 * uint(lane) % 64))
	*(*uint32)(unsafe.Add(m.M, uintptr(ea))) = x
	return 0
}

//go:noinline
func Simd_v128_store64_lane(m *Module, addr int32, offset int32, lane int32, v [2]uint64) int32 {
	ea := SimdEA(m, addr, offset, 8)
	x := uint64(v[int(lane)*64/64] >> (64 * uint(lane) % 64))
	*(*uint64)(unsafe.Add(m.M, uintptr(ea))) = x
	return 0
}

//go:noinline
func Simd_p_v128_load_rng(m *Module, addr int32, offset int32, rlo int32, span int32) (uint64, uint64) {
	r := Simd_v128_load_rng(m, addr, offset, rlo, span)
	return r[0], r[1]
}

//go:noinline
func Simd_p_v128_load_nc(m *Module, addr int32, offset int32) (uint64, uint64) {
	r := Simd_v128_load_nc(m, addr, offset)
	return r[0], r[1]
}

//go:noinline
func Simd_p_v128_load(m *Module, addr int32, offset int32) (uint64, uint64) {
	r := Simd_v128_load(m, addr, offset)
	return r[0], r[1]
}

//go:noinline
func Simd_p_v128_load8x8_u(m *Module, addr int32, offset int32) (uint64, uint64) {
	r := Simd_v128_load8x8_u(m, addr, offset)
	return r[0], r[1]
}

//go:noinline
func Simd_p_v128_load8_splat(m *Module, addr int32, offset int32) (uint64, uint64) {
	r := Simd_v128_load8_splat(m, addr, offset)
	return r[0], r[1]
}

//go:noinline
func Simd_p_v128_load16_splat(m *Module, addr int32, offset int32) (uint64, uint64) {
	r := Simd_v128_load16_splat(m, addr, offset)
	return r[0], r[1]
}

//go:noinline
func Simd_p_v128_load32_splat(m *Module, addr int32, offset int32) (uint64, uint64) {
	r := Simd_v128_load32_splat(m, addr, offset)
	return r[0], r[1]
}

//go:noinline
func Simd_p_v128_load64_splat(m *Module, addr int32, offset int32) (uint64, uint64) {
	r := Simd_v128_load64_splat(m, addr, offset)
	return r[0], r[1]
}

//go:noinline
func Simd_p_v128_load32_zero(m *Module, addr int32, offset int32) (uint64, uint64) {
	r := Simd_v128_load32_zero(m, addr, offset)
	return r[0], r[1]
}

//go:noinline
func Simd_p_v128_load64_zero(m *Module, addr int32, offset int32) (uint64, uint64) {
	r := Simd_v128_load64_zero(m, addr, offset)
	return r[0], r[1]
}

//go:noinline
func Simd_p_v128_load8_lane(m *Module, addr int32, offset int32, lane int32, v0, v1 uint64) (uint64, uint64) {
	r := Simd_v128_load8_lane(m, addr, offset, lane, [2]uint64{v0, v1})
	return r[0], r[1]
}

//go:noinline
func Simd_p_v128_load16_lane(m *Module, addr int32, offset int32, lane int32, v0, v1 uint64) (uint64, uint64) {
	r := Simd_v128_load16_lane(m, addr, offset, lane, [2]uint64{v0, v1})
	return r[0], r[1]
}

//go:noinline
func Simd_p_v128_store(m *Module, addr int32, offset int32, v0, v1 uint64) int32 {
	return Simd_v128_store(m, addr, offset, [2]uint64{v0, v1})
}

//go:noinline
func Simd_p_v128_store8_lane(m *Module, addr int32, offset int32, lane int32, v0, v1 uint64) int32 {
	return Simd_v128_store8_lane(m, addr, offset, lane, [2]uint64{v0, v1})
}

//go:noinline
func Simd_p_v128_store16_lane(m *Module, addr int32, offset int32, lane int32, v0, v1 uint64) int32 {
	return Simd_v128_store16_lane(m, addr, offset, lane, [2]uint64{v0, v1})
}

//go:noinline
func Simd_p_v128_store32_lane(m *Module, addr int32, offset int32, lane int32, v0, v1 uint64) int32 {
	return Simd_v128_store32_lane(m, addr, offset, lane, [2]uint64{v0, v1})
}

//go:noinline
func Simd_p_v128_store64_lane(m *Module, addr int32, offset int32, lane int32, v0, v1 uint64) int32 {
	return Simd_v128_store64_lane(m, addr, offset, lane, [2]uint64{v0, v1})
}

// simdEA64 is simdEA for 64-bit addresses: u64 effective address with
// an overflow-safe range check.
func SimdEA64(m *Module, addr int64, offset int64, size uint64) uint64 {
	ea := uint64(addr) + uint64(offset)
	end := ea + size
	if ea < uint64(addr) || end < ea || end > m.MemSize.Load() {
		Wasm_trap_simd_oob()
	}
	return ea
}

//go:noinline
func Simd_m64_v128_f16x4_cvt_store(m *Module, addr int64, offset int64, v [2]uint64) int32 {

	ea := SimdEA64(m, addr, offset, 8)
	var out uint64
	for i := 0; i < 4; i++ {
		w := uint32(v[i/2] >> (32 * uint(i) % 64))
		shl1w := w + w
		sign := w & 0x80000000
		var h uint32
		if shl1w > 0xFF000000 {
			h = (sign >> 16) | 0x7E00
		} else {
			bias := shl1w & 0xFF000000
			if bias < 0x71000000 {
				bias = 0x71000000
			}
			f := math.Float32frombits(w&0x7FFFFFFF) * 0x1p+112 * 0x1p-110
			f += math.Float32frombits((bias >> 1) + 0x07800000)
			fbits := math.Float32bits(f)
			h = (sign >> 16) | (fbits>>13)&0x7C00 + fbits&0xFFF
		}
		out |= uint64(uint16(h)) << (16 * uint(i))
	}
	*(*uint64)(unsafe.Add(m.M, uintptr(ea))) = out
	return 0
}

//go:noinline
func Simd_p_fx0(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n2 := Simd_i16x8_mul(n0, n1)
	n3 := Simd_i16x8_add(n2, [2]uint64{p4, p4h})
	return n3[0], n3[1]
}

//go:noinline
func Simd_p_fx1(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_extend_low_i16x8_u([2]uint64{p0, p0h})
	n1 := Simd_i32x4_mul(n0, [2]uint64{p1, p1h})
	n2 := Simd_i32x4_extend_high_i16x8_u([2]uint64{p0, p0h})
	n3 := Simd_i32x4_mul(n2, [2]uint64{p1, p1h})
	n4 := Simd_i8x16_shuffle(n1, n3, [2]uint64{1084816697938281218, 2242259463347507986})
	n5 := Simd_i8x16_narrow_i16x8_u(n4, [2]uint64{p2, p2h})
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx2(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1369958511735279616, 1659319203087586308})
	n1 := Simd_v128_or(n0, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p1, p1h}, [2]uint64{361421592565516038, 1084818905618843912})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p1, p1h}, [2]uint64{506097522914230528, 940142975270129422})
	n4 := Simd_i16x8_mul(n3, n0)
	n5 := Simd_i16x8_add(n4, [2]uint64{p3, p3h})
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx3(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i32x4_shr_u(n0, 24)
	n2 := Simd_i32x4_mul(n1, [2]uint64{p0, p0h})
	return n0[0], n0[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx4(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_v128_and([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i32x4_mul([2]uint64{p0, p0h}, n0)
	n2 := Simd_i32x4_add(n1, [2]uint64{p3, p3h})
	n3 := Simd_i32x4_shr_u(n2, 24)
	n4 := Simd_v128_and([2]uint64{p1, p1h}, [2]uint64{p4, p4h})
	n5 := Simd_v128_or(n3, n4)
	n6 := Simd_i32x4_shr_u([2]uint64{p1, p1h}, 8)
	n7 := Simd_v128_and(n6, [2]uint64{p2, p2h})
	n8 := Simd_i32x4_mul([2]uint64{p0, p0h}, n7)
	n9 := Simd_i32x4_add(n8, [2]uint64{p3, p3h})
	n10 := Simd_i32x4_shr_u(n9, 16)
	n11 := Simd_v128_and(n10, [2]uint64{p5, p5h})
	n12 := Simd_v128_or(n5, n11)
	n13 := Simd_i32x4_shr_u([2]uint64{p1, p1h}, 16)
	n14 := Simd_v128_and(n13, [2]uint64{p2, p2h})
	n15 := Simd_i32x4_mul([2]uint64{p0, p0h}, n14)
	n16 := Simd_i32x4_add(n15, [2]uint64{p3, p3h})
	n17 := Simd_i32x4_shr_u(n16, 8)
	n18 := Simd_v128_and(n17, [2]uint64{p6, p6h})
	n19 := Simd_v128_or(n12, n18)
	return n19[0], n19[1]
}

//go:noinline
func Simd_p_fx5(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i32x4_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i32x4_lt_u(n1, [2]uint64{p1, p1h})
	return n0[0], n0[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx6(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p0, p0h}, [2]uint64{1369958511735279616, 1659319203087586308})
	n2 := Simd_v128_or(n1, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n4 := Simd_i8x16_shuffle(n3, [2]uint64{p0, p0h}, [2]uint64{p3, p3h})
	n5 := Simd_i16x8_mul(n4, n1)
	n6 := Simd_i32x4_extend_low_i16x8_u(n5)
	n7 := Simd_i32x4_mul(n6, [2]uint64{p4, p4h})
	n8 := Simd_i8x16_shuffle(n0, [2]uint64{p0, p0h}, [2]uint64{1948679894439893000, 2238040585792199692})
	n9 := Simd_v128_or(n8, [2]uint64{p1, p1h})
	n10 := Simd_i8x16_shuffle(n9, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n11 := Simd_i8x16_shuffle(n10, [2]uint64{p0, p0h}, [2]uint64{p3, p3h})
	n12 := Simd_i16x8_mul(n11, n8)
	n13 := Simd_i32x4_extend_low_i16x8_u(n12)
	n14 := Simd_i32x4_mul(n13, [2]uint64{p4, p4h})
	n15 := Simd_i32x4_extend_high_i16x8_u(n5)
	n16 := Simd_i32x4_mul(n15, [2]uint64{p4, p4h})
	n17 := Simd_i8x16_shuffle(n7, n16, [2]uint64{p5, p5h})
	n18 := Simd_i16x8_shr_u(n17, 7)
	n19 := Simd_i32x4_extend_high_i16x8_u(n12)
	n20 := Simd_i32x4_mul(n19, [2]uint64{p4, p4h})
	n21 := Simd_i8x16_shuffle(n14, n20, [2]uint64{p5, p5h})
	n22 := Simd_i16x8_shr_u(n21, 7)
	n23 := Simd_i8x16_narrow_i16x8_u(n18, n22)
	_ = Simd_v128_store(m, s0, 0, n23)
	return
}

//go:noinline
func Simd_p_fx7(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p0, p0h}, [2]uint64{1948679894439893000, 2238040585792199692})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	return n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx8(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1369958511735279616, 1659319203087586308})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	return n0[0], n0[1], n1[0], n1[1]
}

//go:noinline
func Simd_p_fx9(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1369958511735279616, 1659319203087586308})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p1, p1h}, [2]uint64{1952885526417115400, 2242246217769422092})
	return n0[0], n0[1], n1[0], n1[1]
}

//go:noinline
func Simd_p_fx10(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_splat(s0)
	n1 := Simd_i32x4_add(n0, [2]uint64{p0, p0h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx11(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_extend_low_i8x16_u([2]uint64{p1, p1h})
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	n2 := Simd_v128_and([2]uint64{p0, p0h}, n1)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx12(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1084818905618843912, 216736831629295872})
	n1 := Simd_v128_and([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx13(m *Module, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{216736831696667908, 216736831629295872})
	n1 := Simd_v128_and([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx14(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load32_zero(m, n0, 0)
	n2 := Simd_i16x8_extend_low_i8x16_u(n1)
	n3 := Simd_i32x4_extend_low_i16x8_u(n2)
	n4 := Simd_i32x4_shl(n3, 8)
	_ = Simd_v128_store(m, s2, 0, n4)
	return
}

//go:noinline
func Simd_p_fx15(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n1, [2]uint64{1369958511735279616, 1659319203087586308})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n4 := Simd_i8x16_shuffle(n2, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n5 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n1, [2]uint64{1948679894439893000, 2238040585792199692})
	n6 := Simd_i8x16_shuffle(n5, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n7 := Simd_i8x16_shuffle(n5, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	_ = Simd_v128_store(m, s2, 0, n3)
	_ = Simd_v128_store(m, s2+16, 0, n4)
	_ = Simd_v128_store(m, s2+32, 0, n6)
	_ = Simd_v128_store(m, s2+48, 0, n7)
	return
}

//go:noinline
func Simd_p_fx16(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load32_zero(m, s0, 0)
	n1 := Simd_i16x8_extend_low_i8x16_u(n0)
	n2 := Simd_i32x4_extend_low_i16x8_u(n1)
	n3 := Simd_i32x4_shl(n2, 8)
	_ = Simd_v128_store(m, s1, 0, n3)
	return
}

//go:noinline
func Simd_p_fx17(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 64)
	n1 := Simd_v128_and(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_load_nc(m, s0, 16)
	n3 := Simd_v128_and(n2, [2]uint64{p0, p0h})
	n4 := Simd_i16x8_narrow_i32x4_s(n1, n3)
	n5 := Simd_i8x16_narrow_i16x8_u(n4, n4)
	n6 := Simd_v128_and([2]uint64{p1, p1h}, n5)
	n7 := Simd_v128_load_nc(m, s0, 32)
	n8 := Simd_v128_and(n7, [2]uint64{p0, p0h})
	n9 := Simd_v128_load_nc(m, s0, 48)
	n10 := Simd_v128_and(n9, [2]uint64{p0, p0h})
	n11 := Simd_i16x8_narrow_i32x4_s(n8, n10)
	n12 := Simd_i8x16_narrow_i16x8_u(n11, n11)
	n13 := Simd_i8x16_shuffle(n5, n12, [2]uint64{506097522914230528, 1663540288323457296})
	n14 := Simd_v128_and(n6, n12)
	n15 := Simd_scalar_i32_add(s1, s2)
	_ = Simd_v128_store(m, n15, 0, n13)
	return n14[0], n14[1]
}

//go:noinline
func Simd_p_fx18(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 32)
	n1 := Simd_v128_and(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_load_nc(m, s0, 16)
	n3 := Simd_v128_and(n2, [2]uint64{p0, p0h})
	n4 := Simd_i16x8_narrow_i32x4_s(n1, n3)
	n5 := Simd_i8x16_narrow_i16x8_u(n4, n4)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx19(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_splat(s0)
	n1 := Simd_v128_or(n0, [2]uint64{p0, p0h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx20(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 64)
	n1 := Simd_i32x4_shr_u(n0, 8)
	n2 := Simd_v128_and(n1, [2]uint64{p0, p0h})
	n3 := Simd_v128_load_nc(m, s0, 16)
	n4 := Simd_i32x4_shr_u(n3, 8)
	n5 := Simd_v128_and(n4, [2]uint64{p0, p0h})
	n6 := Simd_i16x8_narrow_i32x4_s(n2, n5)
	n7 := Simd_v128_load_nc(m, s0, 32)
	n8 := Simd_i32x4_shr_u(n7, 8)
	n9 := Simd_v128_and(n8, [2]uint64{p0, p0h})
	n10 := Simd_v128_load_nc(m, s0, 48)
	n11 := Simd_i32x4_shr_u(n10, 8)
	n12 := Simd_v128_and(n11, [2]uint64{p0, p0h})
	n13 := Simd_i16x8_narrow_i32x4_s(n9, n12)
	n14 := Simd_i8x16_narrow_i16x8_u(n6, n13)
	n15 := Simd_scalar_i32_add(s1, s2)
	_ = Simd_v128_store(m, n15, 0, n14)
	return
}

//go:noinline
func Simd_p_fx21(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 32)
	n1 := Simd_i32x4_shr_u(n0, 8)
	n2 := Simd_v128_and(n1, [2]uint64{p0, p0h})
	n3 := Simd_v128_load_nc(m, s0, 16)
	n4 := Simd_i32x4_shr_u(n3, 8)
	n5 := Simd_v128_and(n4, [2]uint64{p0, p0h})
	n6 := Simd_i16x8_narrow_i32x4_s(n2, n5)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx22(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_shr_u([2]uint64{p0, p0h}, 8)
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p1, p1h}, [2]uint64{201851904, 0})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx23(m *Module, s0 int32, s1 int32) (uint64, uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx24(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_and([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_i16x8_narrow_i32x4_s(n0, n1)
	n3 := Simd_v128_and([2]uint64{p3, p3h}, [2]uint64{p1, p1h})
	n4 := Simd_v128_and([2]uint64{p4, p4h}, [2]uint64{p1, p1h})
	n5 := Simd_i16x8_narrow_i32x4_s(n3, n4)
	n6 := Simd_i8x16_narrow_i16x8_u(n2, n5)
	n7 := Simd_i8x16_eq(n6, [2]uint64{p5, p5h})
	return n7[0], n7[1]
}

//go:noinline
func Simd_p_fx25(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 32)
	n1 := Simd_v128_and(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_load_nc(m, s0, 16)
	n3 := Simd_v128_and(n2, [2]uint64{p0, p0h})
	n4 := Simd_i16x8_narrow_i32x4_s(n1, n3)
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx26(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_narrow_i16x8_u([2]uint64{p0, p0h}, [2]uint64{p0, p0h})
	n1 := Simd_i8x16_eq(n0, [2]uint64{p1, p1h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx27(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i32x4_gt_u(n0, [2]uint64{p1, p1h})
	n2 := Simd_v128_bitselect(n0, [2]uint64{p0, p0h}, n1)
	_ = Simd_v128_store(m, s0, 0, n2)
	n4 := Simd_v128_load(m, s0+16, 0)
	n5 := Simd_i32x4_gt_u(n4, [2]uint64{p1, p1h})
	n6 := Simd_v128_bitselect(n4, [2]uint64{p0, p0h}, n5)
	_ = Simd_v128_store(m, s0+16, 0, n6)
	n8 := Simd_v128_load(m, s0+32, 0)
	n9 := Simd_i32x4_gt_u(n8, [2]uint64{p1, p1h})
	n10 := Simd_v128_bitselect(n8, [2]uint64{p0, p0h}, n9)
	_ = Simd_v128_store(m, s0+32, 0, n10)
	n12 := Simd_v128_load(m, s0+48, 0)
	n13 := Simd_i32x4_gt_u(n12, [2]uint64{p1, p1h})
	n14 := Simd_v128_bitselect(n12, [2]uint64{p0, p0h}, n13)
	_ = Simd_v128_store(m, s0+48, 0, n14)
	return
}

//go:noinline
func Simd_p_fx28(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_v128_load(m, s0, 16)
	n1 := Simd_i32x4_gt_u(n0, [2]uint64{p1, p1h})
	n2 := Simd_v128_bitselect(n0, [2]uint64{p0, p0h}, n1)
	_ = Simd_v128_store(m, s0, 16, n2)
	n4 := Simd_v128_load(m, s0, 0)
	n5 := Simd_i32x4_gt_u(n4, [2]uint64{p1, p1h})
	n6 := Simd_v128_bitselect(n4, [2]uint64{p0, p0h}, n5)
	_ = Simd_v128_store(m, s0, 0, n6)
	return
}

//go:noinline
func Simd_p_fx29(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i32x4_lt_u(n0, [2]uint64{p0, p0h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx30(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 16, 0, 64)
	n1 := Simd_i8x16_swizzle(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_load_nc(m, s0, 0)
	n3 := Simd_i8x16_swizzle(n2, [2]uint64{p1, p1h})
	n4 := Simd_v128_or(n1, n3)
	n5 := Simd_v128_load_nc(m, s0, 32)
	n6 := Simd_i8x16_swizzle(n5, [2]uint64{p2, p2h})
	n7 := Simd_v128_or(n4, n6)
	n8 := Simd_v128_load_nc(m, s0, 48)
	n9 := Simd_i8x16_swizzle(n8, [2]uint64{p3, p3h})
	n10 := Simd_v128_or(n7, n9)
	n11 := Simd_v128_and([2]uint64{p4, p4h}, n10)
	n12 := Simd_scalar_i32_add(s1, s2)
	_ = Simd_v128_store(m, n12, 0, n10)
	return n11[0], n11[1]
}

//go:noinline
func Simd_p_fx31(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{216736831696667908, 216736831629295872})
	n1 := Simd_v128_and([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx32(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_shl([2]uint64{p0, p0h}, 1)
	n1 := Simd_v128_or(n0, [2]uint64{p1, p1h})
	return n0[0], n0[1], n1[0], n1[1]
}

//go:noinline
func Simd_p_fx33(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_extend_low_i8x16_u(n0)
	n2 := Simd_i32x4_extend_low_i16x8_u(n1)
	n3 := Simd_i32x4_mul(n2, [2]uint64{p2, p2h})
	return n0[0], n0[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx34(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_or(n0, [2]uint64{p2, p2h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx35(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i16x8_extend_low_i8x16_u(n0)
	n2 := Simd_i32x4_extend_low_i16x8_u(n1)
	n3 := Simd_i32x4_mul(n2, [2]uint64{p3, p3h})
	return n3[0], n3[1]
}

//go:noinline
func Simd_p_fx36(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_extend_low_i8x16_u([2]uint64{p1, p1h})
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	n2 := Simd_i32x4_mul([2]uint64{p0, p0h}, n1)
	n3 := Simd_i32x4_shr_u(n2, 16)
	n4 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n5 := Simd_i16x8_extend_low_i8x16_u(n4)
	n6 := Simd_i32x4_extend_low_i16x8_u(n5)
	n7 := Simd_i32x4_mul([2]uint64{p2, p2h}, n6)
	n8 := Simd_i32x4_shr_u(n7, 16)
	n9 := Simd_i16x8_narrow_i32x4_u(n3, n8)
	return n9[0], n9[1]
}

//go:noinline
func Simd_p_fx37(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n1 := Simd_i16x8_extend_low_i8x16_u(n0)
	n2 := Simd_i32x4_extend_low_i16x8_u(n1)
	n3 := Simd_i32x4_mul([2]uint64{p0, p0h}, n2)
	n4 := Simd_i32x4_shr_u(n3, 16)
	n5 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p5, p5h})
	n6 := Simd_i16x8_extend_low_i8x16_u(n5)
	n7 := Simd_i32x4_extend_low_i16x8_u(n6)
	n8 := Simd_i32x4_mul([2]uint64{p4, p4h}, n7)
	n9 := Simd_i32x4_shr_u(n8, 16)
	n10 := Simd_i16x8_narrow_i32x4_u(n4, n9)
	return n10[0], n10[1]
}

//go:noinline
func Simd_p_fx38(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_extend_low_i8x16_u([2]uint64{p1, p1h})
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	n2 := Simd_i32x4_mul([2]uint64{p0, p0h}, n1)
	n3 := Simd_i32x4_shr_u(n2, 20)
	n4 := Simd_v128_and(n3, [2]uint64{p2, p2h})
	n5 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p1, p1h}, [2]uint64{p4, p4h})
	n6 := Simd_i16x8_extend_low_i8x16_u(n5)
	n7 := Simd_i32x4_extend_low_i16x8_u(n6)
	n8 := Simd_i32x4_mul([2]uint64{p3, p3h}, n7)
	n9 := Simd_i32x4_shr_u(n8, 20)
	n10 := Simd_v128_and(n9, [2]uint64{p2, p2h})
	n11 := Simd_i16x8_narrow_i32x4_u(n4, n10)
	return n11[0], n11[1]
}

//go:noinline
func Simd_p_fx39(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i16x8_extend_low_i8x16_u(n0)
	n2 := Simd_i32x4_extend_low_i16x8_u(n1)
	n3 := Simd_i32x4_mul([2]uint64{p0, p0h}, n2)
	n4 := Simd_i32x4_shr_u(n3, 20)
	n5 := Simd_v128_and(n4, [2]uint64{p3, p3h})
	n6 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p1, p1h}, [2]uint64{p5, p5h})
	n7 := Simd_i16x8_extend_low_i8x16_u(n6)
	n8 := Simd_i32x4_extend_low_i16x8_u(n7)
	n9 := Simd_i32x4_mul([2]uint64{p4, p4h}, n8)
	n10 := Simd_i32x4_shr_u(n9, 20)
	n11 := Simd_v128_and(n10, [2]uint64{p3, p3h})
	n12 := Simd_i16x8_narrow_i32x4_u(n5, n11)
	return n12[0], n12[1]
}

//go:noinline
func Simd_p_fx40(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_v128_load32_zero(m, s0, 0)
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p0, p0h}, [2]uint64{4294967296, 12884901890})
	n2 := Simd_i32x4_shl(n1, 16)
	n3 := Simd_v128_load32_zero(m, s1, 0)
	n4 := Simd_i16x8_extend_low_i8x16_u(n3)
	n5 := Simd_i32x4_extend_low_i16x8_u(n4)
	n6 := Simd_i32x4_shl(n5, 8)
	n7 := Simd_v128_or(n2, n6)
	n8 := Simd_v128_load32_zero(m, s2, 0)
	n9 := Simd_i16x8_extend_low_i8x16_u(n8)
	n10 := Simd_i32x4_extend_low_i16x8_u(n9)
	n11 := Simd_v128_or(n7, n10)
	n12 := Simd_v128_or(n11, [2]uint64{p1, p1h})
	_ = Simd_v128_store(m, s3, 0, n12)
	return
}

//go:noinline
func Simd_p_fx41(m *Module, s0 int32, s1 int32) (uint64, uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load32_zero(m, n0, 0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx42(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1084818905618843912, 0})
	n1 := Simd_v128_and([2]uint64{p0, p0h}, n0)
	n2 := Simd_i8x16_shuffle(n1, n1, [2]uint64{117835012, 0})
	n3 := Simd_v128_and(n1, n2)
	n4 := Simd_i8x16_shuffle(n3, n3, [2]uint64{770, 0})
	n5 := Simd_v128_and(n3, n4)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx43(m *Module, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{1, 0})
	n1 := Simd_v128_and([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx44(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_shr_u([2]uint64{p0, p0h}, 8)
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p1, p1h}, [2]uint64{201851904, 0})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx45(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	n2 := Simd_i32x4_mul(n1, [2]uint64{p3, p3h})
	n3 := Simd_i32x4_extend_high_i16x8_u(n0)
	n4 := Simd_i32x4_mul(n3, [2]uint64{p3, p3h})
	n5 := Simd_i8x16_shuffle(n2, n4, [2]uint64{p4, p4h})
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx46(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	n2 := Simd_i32x4_extend_high_i16x8_u(n0)
	return n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx47(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_mul([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p3, p3h})
	n3 := Simd_i32x4_mul([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n4 := Simd_i32x4_mul([2]uint64{p6, p6h}, [2]uint64{p5, p5h})
	n5 := Simd_i8x16_shuffle(n3, n4, [2]uint64{p3, p3h})
	n6 := Simd_i16x8_add(n2, n5)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx48(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i32x4_mul([2]uint64{p3, p3h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p4, p4h})
	n3 := Simd_i16x8_add([2]uint64{p0, p0h}, n2)
	n4 := Simd_i16x8_add(n3, [2]uint64{p5, p5h})
	n5 := Simd_i16x8_shr_s(n4, 6)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx49(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_mul([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p3, p3h})
	n3 := Simd_i16x8_add_sat_u(n2, [2]uint64{p4, p4h})
	n4 := Simd_i16x8_sub_sat_u(n3, [2]uint64{p5, p5h})
	n5 := Simd_i16x8_shr_u(n4, 6)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx50(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_i8x16_swizzle([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_swizzle([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_v128_or(n0, n1)
	n3 := Simd_i8x16_swizzle([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n4 := Simd_v128_or(n2, n3)
	_ = Simd_v128_store(m, s0, 80, n4)
	return
}

//go:noinline
func Simd_p_fx51(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_i8x16_swizzle([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_swizzle([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_v128_or(n0, n1)
	n3 := Simd_i8x16_swizzle([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n4 := Simd_v128_or(n2, n3)
	_ = Simd_v128_store(m, s0, 64, n4)
	return
}

//go:noinline
func Simd_p_fx52(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_i8x16_swizzle([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_swizzle([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_v128_or(n0, n1)
	n3 := Simd_i8x16_swizzle([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n4 := Simd_v128_or(n2, n3)
	_ = Simd_v128_store(m, s0, 48, n4)
	return
}

//go:noinline
func Simd_p_fx53(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_i8x16_swizzle([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_swizzle([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_v128_or(n0, n1)
	n3 := Simd_i8x16_swizzle([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n4 := Simd_v128_or(n2, n3)
	_ = Simd_v128_store(m, s0, 32, n4)
	return
}

//go:noinline
func Simd_p_fx54(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_i8x16_swizzle([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_swizzle([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_v128_or(n0, n1)
	n3 := Simd_i8x16_swizzle([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n4 := Simd_v128_or(n2, n3)
	_ = Simd_v128_store(m, s0, 16, n4)
	return
}

//go:noinline
func Simd_p_fx55(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_i8x16_swizzle([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_swizzle([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_v128_or(n0, n1)
	n3 := Simd_i8x16_swizzle([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n4 := Simd_v128_or(n2, n3)
	_ = Simd_v128_store(m, s0, 0, n4)
	return
}

//go:noinline
func Simd_p_fx56(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 128)
	n1 := Simd_i8x16_swizzle(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_load_nc(m, s0+16, 0)
	n3 := Simd_i8x16_swizzle(n2, [2]uint64{p0, p0h})
	n4 := Simd_v128_load_nc(m, s0+32, 0)
	n5 := Simd_i8x16_swizzle(n4, [2]uint64{p0, p0h})
	n6 := Simd_v128_load_nc(m, s0+48, 0)
	n7 := Simd_i8x16_swizzle(n6, [2]uint64{p0, p0h})
	return n1[0], n1[1], n3[0], n3[1], n5[0], n5[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx57(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p4, p4h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p5, p5h})
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx58(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i32x4_dot_i16x8_s(n0, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p4, p4h})
	n3 := Simd_i32x4_dot_i16x8_s(n2, [2]uint64{p3, p3h})
	n4 := Simd_i16x8_narrow_i32x4_s(n1, n3)
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx59(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p4, p4h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p5, p5h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx60(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i32x4_dot_i16x8_s(n0, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p4, p4h})
	n3 := Simd_i32x4_dot_i16x8_s(n2, [2]uint64{p3, p3h})
	n4 := Simd_i16x8_narrow_i32x4_s(n1, n3)
	n5 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, n4, [2]uint64{p6, p6h})
	return n4[0], n4[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx61(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n2 := Simd_i32x4_dot_i16x8_s(n1, [2]uint64{p5, p5h})
	n3 := Simd_i8x16_shuffle(n0, [2]uint64{p3, p3h}, [2]uint64{p6, p6h})
	n4 := Simd_i32x4_dot_i16x8_s(n3, [2]uint64{p5, p5h})
	n5 := Simd_i16x8_narrow_i32x4_s(n2, n4)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx62(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p0, p0h}, [2]uint64{p4, p4h})
	n2 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p4, p4h})
	n3 := Simd_v128_load_nc(m, s0+64, 0)
	n4 := Simd_i8x16_swizzle(n3, [2]uint64{p5, p5h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx63(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_nc(m, s0+80, 0)
	n1 := Simd_i8x16_swizzle(n0, [2]uint64{p0, p0h})
	n2 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, n1, [2]uint64{p2, p2h})
	n3 := Simd_v128_load_nc(m, s0+96, 0)
	n4 := Simd_i8x16_swizzle(n3, [2]uint64{p0, p0h})
	n5 := Simd_v128_load_nc(m, s0+112, 0)
	n6 := Simd_i8x16_swizzle(n5, [2]uint64{p0, p0h})
	n7 := Simd_i8x16_shuffle(n4, n6, [2]uint64{p2, p2h})
	n8 := Simd_i8x16_shuffle(n2, n7, [2]uint64{p3, p3h})
	return n1[0], n1[1], n4[0], n4[1], n6[0], n6[1], n8[0], n8[1]
}

//go:noinline
func Simd_p_fx64(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_dot_i16x8_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_dot_i16x8_s([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i32x4_add(n0, n1)
	n3 := Simd_i32x4_add(n2, [2]uint64{p4, p4h})
	n4 := Simd_i32x4_shr_s(n3, 18)
	n5 := Simd_i32x4_dot_i16x8_s([2]uint64{p5, p5h}, [2]uint64{p1, p1h})
	n6 := Simd_i32x4_dot_i16x8_s([2]uint64{p6, p6h}, [2]uint64{p3, p3h})
	n7 := Simd_i32x4_add(n5, n6)
	n8 := Simd_i32x4_add(n7, [2]uint64{p4, p4h})
	n9 := Simd_i32x4_shr_s(n8, 18)
	n10 := Simd_i16x8_narrow_i32x4_s(n4, n9)
	return n10[0], n10[1]
}

//go:noinline
func Simd_p_fx65(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i8x16_avgr_u([2]uint64{p0, p0h}, n0)
	n2 := Simd_v128_load(m, s1, 0)
	n3 := Simd_i8x16_avgr_u([2]uint64{p1, p1h}, n2)
	return n1[0], n1[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx66(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, s5 int32, s6 int32, s7 int32, p0, p0h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load32_splat(m, s0, 0)
	n1 := Simd_v128_load32_lane(m, s1, 0, 1, n0)
	n2 := Simd_v128_load32_lane(m, s2, 0, 2, n1)
	n3 := Simd_v128_load32_lane(m, s3, 0, 3, n2)
	n4 := Simd_i32x4_shr_u(n3, 15)
	n5 := Simd_v128_and(n4, [2]uint64{p0, p0h})
	n6 := Simd_i32x4_shr_u(n3, 7)
	n7 := Simd_v128_and(n6, [2]uint64{p0, p0h})
	n8 := Simd_i32x4_shl(n3, 1)
	n9 := Simd_v128_and(n8, [2]uint64{p0, p0h})
	n10 := Simd_v128_load32_splat(m, s4, 0)
	n11 := Simd_v128_load32_lane(m, s5, 0, 1, n10)
	n12 := Simd_v128_load32_lane(m, s6, 0, 2, n11)
	n13 := Simd_v128_load32_lane(m, s7, 0, 3, n12)
	n14 := Simd_i32x4_shr_u(n13, 15)
	n15 := Simd_v128_and(n14, [2]uint64{p0, p0h})
	n16 := Simd_i32x4_add(n5, n15)
	n17 := Simd_i32x4_shr_u(n13, 7)
	n18 := Simd_v128_and(n17, [2]uint64{p0, p0h})
	n19 := Simd_i32x4_add(n7, n18)
	n20 := Simd_i32x4_shl(n13, 1)
	n21 := Simd_v128_and(n20, [2]uint64{p0, p0h})
	n22 := Simd_i32x4_add(n9, n21)
	return n16[0], n16[1], n19[0], n19[1], n22[0], n22[1]
}

//go:noinline
func Simd_p_fx67(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_mul([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i32x4_add(n0, n1)
	n3 := Simd_i32x4_mul([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n4 := Simd_i32x4_add(n2, n3)
	n5 := Simd_i32x4_add(n4, [2]uint64{p6, p6h})
	n6 := Simd_i32x4_shr_u(n5, 18)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx68(m *Module, s0 int32) (uint64, uint64) {
	n0 := Simd_v128_load32_zero(m, s0, 0)
	n1 := Simd_i16x8_extend_low_i8x16_u(n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx69(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_v128_or([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n2 := Simd_i16x8_shr_u(n1, 1)
	n3 := Simd_i16x8_sub(n0, n2)
	n4 := Simd_i8x16_shuffle(n3, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx70(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_mul([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i32x4_add(n0, n1)
	n3 := Simd_i32x4_mul([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n4 := Simd_i32x4_add(n2, n3)
	n5 := Simd_i32x4_add(n4, [2]uint64{p6, p6h})
	n6 := Simd_i32x4_shr_u(n5, 18)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx71(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_v128_or([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n2 := Simd_i16x8_shr_u(n1, 1)
	n3 := Simd_i16x8_sub(n0, n2)
	n4 := Simd_i8x16_shuffle(n3, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx72(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_mul([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i32x4_add(n0, n1)
	n3 := Simd_i32x4_mul([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n4 := Simd_i32x4_add(n2, n3)
	n5 := Simd_i32x4_add(n4, [2]uint64{p6, p6h})
	n6 := Simd_i32x4_shr_u(n5, 18)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx73(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_mul([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i32x4_add(n0, n1)
	n3 := Simd_i32x4_mul([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n4 := Simd_i32x4_add(n2, n3)
	n5 := Simd_i32x4_add(n4, [2]uint64{p6, p6h})
	n6 := Simd_i32x4_shr_u(n5, 18)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx74(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 64)
	n1 := Simd_i8x16_swizzle(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_load_nc(m, s0+16, 0)
	n3 := Simd_i8x16_swizzle(n2, [2]uint64{p0, p0h})
	n4 := Simd_i8x16_shuffle(n1, n3, [2]uint64{p1, p1h})
	n5 := Simd_i8x16_shuffle(n1, n3, [2]uint64{p5, p5h})
	n6 := Simd_v128_load_nc(m, s0+32, 0)
	n7 := Simd_i8x16_swizzle(n6, [2]uint64{p0, p0h})
	n8 := Simd_v128_load_nc(m, s0+48, 0)
	n9 := Simd_i8x16_swizzle(n8, [2]uint64{p0, p0h})
	n10 := Simd_i8x16_shuffle(n7, n9, [2]uint64{p1, p1h})
	n11 := Simd_i8x16_shuffle(n4, n10, [2]uint64{p2, p2h})
	n12 := Simd_i8x16_shuffle(n11, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n13 := Simd_i8x16_shuffle(n7, n9, [2]uint64{p5, p5h})
	return n11[0], n11[1], n12[0], n12[1], n5[0], n5[1], n13[0], n13[1]
}

//go:noinline
func Simd_p_fx75(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1084818905618843912, 2242261671028070680})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p4, p4h})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx76(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i32x4_dot_i16x8_s(n0, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p4, p4h}, [2]uint64{p2, p2h})
	n3 := Simd_i32x4_dot_i16x8_s(n2, [2]uint64{p5, p5h})
	n4 := Simd_i32x4_add(n1, n3)
	n5 := Simd_i32x4_add(n4, [2]uint64{p6, p6h})
	n6 := Simd_i32x4_shr_s(n5, 16)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx77(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_mul(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_shr_u([2]uint64{p0, p0h}, 16)
	n3 := Simd_v128_and(n2, [2]uint64{p1, p1h})
	n4 := Simd_i32x4_mul(n3, [2]uint64{p3, p3h})
	n5 := Simd_i32x4_add(n1, n4)
	n6 := Simd_i32x4_shr_u([2]uint64{p0, p0h}, 8)
	n7 := Simd_v128_and(n6, [2]uint64{p1, p1h})
	n8 := Simd_i32x4_mul(n7, [2]uint64{p4, p4h})
	n9 := Simd_i32x4_add(n5, n8)
	n10 := Simd_i32x4_add(n9, [2]uint64{p5, p5h})
	n11 := Simd_i32x4_shr_u(n10, 16)
	n12 := Simd_i8x16_shuffle(n11, [2]uint64{p0, p0h}, [2]uint64{201851904, 0})
	return n12[0], n12[1]
}

//go:noinline
func Simd_p_fx78(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 64, 48, 48)
	n1 := Simd_i8x16_swizzle(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_load_nc(m, s0, 48)
	n3 := Simd_i8x16_swizzle(n2, [2]uint64{p1, p1h})
	n4 := Simd_v128_or(n1, n3)
	n5 := Simd_v128_load_nc(m, s0, 80)
	n6 := Simd_i8x16_swizzle(n5, [2]uint64{p2, p2h})
	n7 := Simd_v128_or(n4, n6)
	return n0[0], n0[1], n2[0], n2[1], n5[0], n5[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx79(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_swizzle([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_swizzle([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_v128_or(n0, n1)
	n3 := Simd_i8x16_swizzle([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n4 := Simd_v128_or(n2, n3)
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx80(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	_ = Simd_v128_store(m, s0+16, 0, [2]uint64{p0, p0h})
	n1 := Simd_v128_load_rng(m, s1, 16, 0, 48)
	n2 := Simd_i8x16_swizzle(n1, [2]uint64{p1, p1h})
	n3 := Simd_v128_load_nc(m, s1, 0)
	n4 := Simd_i8x16_swizzle(n3, [2]uint64{p2, p2h})
	n5 := Simd_v128_or(n2, n4)
	n6 := Simd_v128_load_nc(m, s1, 32)
	n7 := Simd_i8x16_swizzle(n6, [2]uint64{p3, p3h})
	n8 := Simd_v128_or(n5, n7)
	return n1[0], n1[1], n3[0], n3[1], n6[0], n6[1], n8[0], n8[1]
}

//go:noinline
func Simd_p_fx81(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_extend_low_i8x16_u([2]uint64{p0, p0h})
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	n2 := Simd_i32x4_mul(n1, [2]uint64{p1, p1h})
	n3 := Simd_i16x8_extend_low_i8x16_u([2]uint64{p2, p2h})
	n4 := Simd_i32x4_extend_low_i16x8_u(n3)
	n5 := Simd_i32x4_mul(n4, [2]uint64{p3, p3h})
	n6 := Simd_i32x4_add(n2, n5)
	n7 := Simd_i16x8_extend_low_i8x16_u([2]uint64{p4, p4h})
	n8 := Simd_i32x4_extend_low_i16x8_u(n7)
	n9 := Simd_i32x4_mul(n8, [2]uint64{p5, p5h})
	n10 := Simd_i32x4_add(n6, n9)
	n11 := Simd_i32x4_add(n10, [2]uint64{p6, p6h})
	n12 := Simd_i32x4_shr_u(n11, s0)
	return n12[0], n12[1]
}

//go:noinline
func Simd_p_fx82(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i32x4_dot_i16x8_s(n0, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle([2]uint64{p4, p4h}, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n3 := Simd_i32x4_dot_i16x8_s(n2, [2]uint64{p5, p5h})
	n4 := Simd_i32x4_add(n1, n3)
	n5 := Simd_i32x4_add(n4, [2]uint64{p6, p6h})
	n6 := Simd_i32x4_shr_s(n5, 16)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx83(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 128)
	n1 := Simd_i8x16_swizzle(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_load_nc(m, s0, 16)
	n3 := Simd_i8x16_swizzle(n2, [2]uint64{p1, p1h})
	n4 := Simd_v128_or(n3, n1)
	n5 := Simd_i8x16_shuffle(n1, n3, [2]uint64{p2, p2h})
	n6 := Simd_v128_load_nc(m, s0, 32)
	n7 := Simd_i8x16_swizzle(n6, [2]uint64{p0, p0h})
	n8 := Simd_v128_load_nc(m, s0, 48)
	n9 := Simd_i8x16_swizzle(n8, [2]uint64{p1, p1h})
	n10 := Simd_v128_or(n9, n7)
	n11 := Simd_i8x16_shuffle(n7, n9, [2]uint64{p2, p2h})
	n12 := Simd_i8x16_shuffle(n4, n10, [2]uint64{p4, p4h})
	n13 := Simd_i8x16_shuffle(n5, n11, [2]uint64{p3, p3h})
	n14 := Simd_i8x16_shuffle(n5, n11, [2]uint64{p4, p4h})
	n15 := Simd_i8x16_shuffle(n14, n12, [2]uint64{p5, p5h})
	n16 := Simd_i8x16_shuffle(n14, n12, [2]uint64{p6, p6h})
	n17 := Simd_i8x16_shuffle(n13, n14, [2]uint64{p5, p5h})
	n18 := Simd_i8x16_shuffle(n13, n14, [2]uint64{p6, p6h})
	return n17[0], n17[1], n15[0], n15[1], n18[0], n18[1], n16[0], n16[1]
}

//go:noinline
func Simd_p_fx84(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_nc(m, s0, 64)
	n1 := Simd_i8x16_swizzle(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_load_nc(m, s0, 80)
	n3 := Simd_i8x16_swizzle(n2, [2]uint64{p1, p1h})
	n4 := Simd_v128_or(n3, n1)
	n5 := Simd_i8x16_shuffle(n1, n3, [2]uint64{p2, p2h})
	n6 := Simd_v128_load_nc(m, s0, 96)
	n7 := Simd_i8x16_swizzle(n6, [2]uint64{p0, p0h})
	n8 := Simd_v128_load_nc(m, s0, 112)
	n9 := Simd_i8x16_swizzle(n8, [2]uint64{p1, p1h})
	n10 := Simd_v128_or(n9, n7)
	n11 := Simd_i8x16_shuffle(n7, n9, [2]uint64{p2, p2h})
	n12 := Simd_i8x16_shuffle(n4, n10, [2]uint64{p4, p4h})
	n13 := Simd_i8x16_shuffle(n5, n11, [2]uint64{p3, p3h})
	n14 := Simd_i8x16_shuffle(n5, n11, [2]uint64{p4, p4h})
	n15 := Simd_i8x16_shuffle(n14, n12, [2]uint64{p5, p5h})
	n16 := Simd_i8x16_shuffle(n14, n12, [2]uint64{p6, p6h})
	n17 := Simd_i8x16_shuffle(n13, n14, [2]uint64{p5, p5h})
	n18 := Simd_i8x16_shuffle(n13, n14, [2]uint64{p6, p6h})
	return n17[0], n17[1], n15[0], n15[1], n18[0], n18[1], n16[0], n16[1]
}

//go:noinline
func Simd_p_fx85(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_mul([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i32x4_add(n0, n1)
	n3 := Simd_i32x4_mul([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n4 := Simd_i32x4_add(n2, n3)
	n5 := Simd_i32x4_add(n4, [2]uint64{p6, p6h})
	n6 := Simd_i32x4_shr_s(n5, 18)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx86(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_extend_low_i8x16_u([2]uint64{p0, p0h})
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	n2 := Simd_i32x4_mul(n1, [2]uint64{p1, p1h})
	n3 := Simd_i16x8_extend_low_i8x16_u([2]uint64{p2, p2h})
	n4 := Simd_i32x4_extend_low_i16x8_u(n3)
	n5 := Simd_i32x4_mul(n4, [2]uint64{p3, p3h})
	n6 := Simd_i32x4_add(n2, n5)
	n7 := Simd_i16x8_extend_low_i8x16_u([2]uint64{p4, p4h})
	n8 := Simd_i32x4_extend_low_i16x8_u(n7)
	n9 := Simd_i32x4_mul(n8, [2]uint64{p5, p5h})
	n10 := Simd_i32x4_add(n6, n9)
	n11 := Simd_i32x4_add(n10, [2]uint64{p6, p6h})
	n12 := Simd_i32x4_shr_u(n11, 16)
	return n12[0], n12[1]
}

//go:noinline
func Simd_p_fx87(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i32x4_mul([2]uint64{p3, p3h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p4, p4h})
	n3 := Simd_i16x8_add([2]uint64{p0, p0h}, n2)
	n4 := Simd_i16x8_add(n3, [2]uint64{p5, p5h})
	n5 := Simd_i16x8_shr_s(n4, 6)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx88(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_mul([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p3, p3h})
	n3 := Simd_i16x8_add_sat_u(n2, [2]uint64{p4, p4h})
	n4 := Simd_i16x8_sub_sat_u(n3, [2]uint64{p5, p5h})
	n5 := Simd_i16x8_shr_u(n4, 6)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx89(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1948679894439893000, 2238040585792199692})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{1952885526417115400, 2242246217769422092})
	n3 := Simd_i8x16_shuffle(n0, n1, [2]uint64{1374164143712502016, 1663524835064808708})
	_ = Simd_v128_store(m, s0, 16, n2)
	_ = Simd_v128_store(m, s0, 0, n3)
	return
}

//go:noinline
func Simd_p_fx90(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_mul([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p3, p3h})
	n3 := Simd_i16x8_add_sat_u(n2, [2]uint64{p4, p4h})
	n4 := Simd_i16x8_sub_sat_u(n3, [2]uint64{p5, p5h})
	n5 := Simd_i16x8_shr_u(n4, 6)
	n6 := Simd_i8x16_narrow_i16x8_u(n5, [2]uint64{p6, p6h})
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx91(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1948679894439893000, 2238040585792199692})
	n1 := Simd_i16x8_shr_u(n0, 4)
	n2 := Simd_v128_and(n1, [2]uint64{p2, p2h})
	n3 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p3, p3h})
	n4 := Simd_v128_and(n3, [2]uint64{p4, p4h})
	n5 := Simd_v128_or(n2, n4)
	_ = Simd_v128_store(m, s0, 0, n5)
	return
}

//go:noinline
func Simd_p_fx92(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_shr_u([2]uint64{p0, p0h}, 5)
	n1 := Simd_v128_and(n0, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_narrow_i16x8_u([2]uint64{p2, p2h}, [2]uint64{p2, p2h})
	n3 := Simd_v128_and(n2, [2]uint64{p3, p3h})
	n4 := Simd_v128_or(n1, n3)
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx93(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_narrow_i16x8_u([2]uint64{p0, p0h}, [2]uint64{p0, p0h})
	n1 := Simd_i16x8_shr_u(n0, 3)
	n2 := Simd_v128_and(n1, [2]uint64{p1, p1h})
	n3 := Simd_i16x8_shl([2]uint64{p2, p2h}, 3)
	n4 := Simd_v128_and(n3, [2]uint64{p3, p3h})
	n5 := Simd_v128_or(n2, n4)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx94(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_and([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_narrow_i16x8_u(n0, n1)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx95(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_shr_u([2]uint64{p0, p0h}, 8)
	n1 := Simd_i16x8_shr_u([2]uint64{p1, p1h}, 8)
	n2 := Simd_i8x16_narrow_i16x8_u(n0, n1)
	n3 := Simd_i16x8_shr_u(n2, 8)
	n4 := Simd_i16x8_shr_u([2]uint64{p2, p2h}, 8)
	n5 := Simd_i8x16_narrow_i16x8_u(n4, n3)
	return n2[0], n2[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx96(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_shr_u([2]uint64{p0, p0h}, 8)
	n1 := Simd_i16x8_shr_u([2]uint64{p1, p1h}, 8)
	n2 := Simd_i8x16_narrow_i16x8_u(n0, n1)
	n3 := Simd_i16x8_shr_u(n2, 8)
	n4 := Simd_i16x8_shr_u([2]uint64{p2, p2h}, 8)
	n5 := Simd_i16x8_shr_u([2]uint64{p3, p3h}, 8)
	n6 := Simd_i8x16_narrow_i16x8_u(n4, n5)
	n7 := Simd_i16x8_shr_u(n6, 8)
	n8 := Simd_i8x16_narrow_i16x8_u(n3, n7)
	n9 := Simd_v128_and(n8, [2]uint64{p5, p5h})
	n10 := Simd_v128_and([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n11 := Simd_i8x16_narrow_i16x8_u(n10, n9)
	return n2[0], n2[1], n6[0], n6[1], n8[0], n8[1], n11[0], n11[1]
}

//go:noinline
func Simd_p_fx97(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_and([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_narrow_i16x8_u(n0, n1)
	n3 := Simd_v128_and(n2, [2]uint64{p1, p1h})
	n4 := Simd_v128_and([2]uint64{p3, p3h}, [2]uint64{p1, p1h})
	n5 := Simd_v128_and([2]uint64{p4, p4h}, [2]uint64{p1, p1h})
	n6 := Simd_i8x16_narrow_i16x8_u(n4, n5)
	n7 := Simd_v128_and(n6, [2]uint64{p1, p1h})
	n8 := Simd_i8x16_narrow_i16x8_u(n3, n7)
	n9 := Simd_v128_and([2]uint64{p5, p5h}, [2]uint64{p1, p1h})
	n10 := Simd_v128_and([2]uint64{p6, p6h}, [2]uint64{p1, p1h})
	n11 := Simd_i8x16_narrow_i16x8_u(n9, n10)
	return n2[0], n2[1], n6[0], n6[1], n8[0], n8[1], n11[0], n11[1]
}

//go:noinline
func Simd_p_fx98(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_shr_u([2]uint64{p0, p0h}, 8)
	n1 := Simd_i16x8_shr_u([2]uint64{p1, p1h}, 8)
	n2 := Simd_i8x16_narrow_i16x8_u(n0, n1)
	n3 := Simd_i16x8_shr_u(n2, 8)
	n4 := Simd_i16x8_shr_u([2]uint64{p2, p2h}, 8)
	n5 := Simd_i8x16_narrow_i16x8_u(n4, n3)
	n6 := Simd_v128_and([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n7 := Simd_v128_and([2]uint64{p5, p5h}, [2]uint64{p4, p4h})
	n8 := Simd_i8x16_narrow_i16x8_u(n6, n7)
	return n2[0], n2[1], n5[0], n5[1], n8[0], n8[1]
}

//go:noinline
func Simd_p_fx99(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_shr_u([2]uint64{p0, p0h}, 8)
	n1 := Simd_i16x8_shr_u([2]uint64{p1, p1h}, 8)
	n2 := Simd_i8x16_narrow_i16x8_u(n0, n1)
	n3 := Simd_i16x8_shr_u(n2, 8)
	n4 := Simd_i16x8_shr_u([2]uint64{p2, p2h}, 8)
	n5 := Simd_i8x16_narrow_i16x8_u(n4, n3)
	n6 := Simd_i16x8_shr_u(n5, 8)
	n7 := Simd_v128_and(n5, [2]uint64{p6, p6h})
	n8 := Simd_i16x8_shr_u([2]uint64{p3, p3h}, 8)
	n9 := Simd_i16x8_shr_u([2]uint64{p4, p4h}, 8)
	n10 := Simd_i8x16_narrow_i16x8_u(n8, n9)
	n11 := Simd_i16x8_shr_u(n10, 8)
	n12 := Simd_i8x16_narrow_i16x8_u(n6, n11)
	n13 := Simd_i16x8_shr_u(n12, 8)
	n14 := Simd_v128_and(n10, [2]uint64{p6, p6h})
	n15 := Simd_i8x16_narrow_i16x8_u(n7, n14)
	n16 := Simd_i16x8_shr_u([2]uint64{p5, p5h}, 8)
	n17 := Simd_i8x16_narrow_i16x8_u(n16, n13)
	_ = Simd_v128_store(m, s0, 80, n17)
	return n2[0], n2[1], n12[0], n12[1], n15[0], n15[1]
}

//go:noinline
func Simd_p_fx100(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_and([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_narrow_i16x8_u(n0, n1)
	n3 := Simd_i16x8_shr_u(n2, 8)
	n4 := Simd_v128_and(n2, [2]uint64{p1, p1h})
	n5 := Simd_v128_and([2]uint64{p3, p3h}, [2]uint64{p1, p1h})
	n6 := Simd_v128_and([2]uint64{p4, p4h}, [2]uint64{p1, p1h})
	n7 := Simd_i8x16_narrow_i16x8_u(n5, n6)
	n8 := Simd_i16x8_shr_u(n7, 8)
	n9 := Simd_i8x16_narrow_i16x8_u(n3, n8)
	n10 := Simd_i16x8_shr_u(n9, 8)
	n11 := Simd_v128_and(n7, [2]uint64{p1, p1h})
	n12 := Simd_i8x16_narrow_i16x8_u(n4, n11)
	n13 := Simd_i16x8_shr_u([2]uint64{p5, p5h}, 8)
	n14 := Simd_i8x16_narrow_i16x8_u(n13, n10)
	_ = Simd_v128_store(m, s0, 64, n14)
	return n9[0], n9[1], n12[0], n12[1]
}

//go:noinline
func Simd_p_fx101(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_and([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_narrow_i16x8_u(n0, n1)
	n3 := Simd_i16x8_shr_u(n2, 8)
	n4 := Simd_i16x8_shr_u([2]uint64{p3, p3h}, 8)
	n5 := Simd_i8x16_narrow_i16x8_u(n4, n3)
	n6 := Simd_v128_and([2]uint64{p4, p4h}, [2]uint64{p1, p1h})
	n7 := Simd_v128_and([2]uint64{p5, p5h}, [2]uint64{p1, p1h})
	n8 := Simd_i8x16_narrow_i16x8_u(n6, n7)
	_ = Simd_v128_store(m, s0, 48, n5)
	_ = Simd_v128_store(m, s0, 32, n8)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx102(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) {
	n0 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_and([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_narrow_i16x8_u(n0, n1)
	n3 := Simd_v128_and([2]uint64{p3, p3h}, [2]uint64{p1, p1h})
	n4 := Simd_v128_and([2]uint64{p4, p4h}, [2]uint64{p1, p1h})
	n5 := Simd_i8x16_narrow_i16x8_u(n3, n4)
	_ = Simd_v128_store(m, s0, 16, n2)
	_ = Simd_v128_store(m, s0, 0, n5)
	return
}

//go:noinline
func Simd_p_fx103(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 128)
	n1 := Simd_v128_load_nc(m, s0+16, 0)
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p0, p0h})
	n3 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p1, p1h})
	n4 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p0, p0h})
	n5 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p1, p1h})
	n6 := Simd_i8x16_shuffle(n4, n5, [2]uint64{p1, p1h})
	n7 := Simd_i8x16_shuffle(n4, n5, [2]uint64{p0, p0h})
	n8 := Simd_v128_load_nc(m, s0+32, 0)
	n9 := Simd_v128_load_nc(m, s0+48, 0)
	n10 := Simd_i8x16_shuffle(n8, n9, [2]uint64{p0, p0h})
	n11 := Simd_i8x16_shuffle(n8, n9, [2]uint64{p1, p1h})
	n12 := Simd_i8x16_shuffle(n10, n11, [2]uint64{p0, p0h})
	n13 := Simd_i8x16_shuffle(n10, n11, [2]uint64{p1, p1h})
	n14 := Simd_i8x16_shuffle(n12, n13, [2]uint64{p1, p1h})
	n15 := Simd_i8x16_shuffle(n6, n14, [2]uint64{p2, p2h})
	n16 := Simd_i8x16_shuffle(n15, [2]uint64{p3, p3h}, [2]uint64{p0, p0h})
	n17 := Simd_i32x4_dot_i16x8_s(n16, [2]uint64{p4, p4h})
	n18 := Simd_i8x16_shuffle(n15, [2]uint64{p3, p3h}, [2]uint64{p1, p1h})
	n19 := Simd_i32x4_dot_i16x8_s(n18, [2]uint64{p4, p4h})
	n20 := Simd_i16x8_narrow_i32x4_s(n17, n19)
	n21 := Simd_i8x16_shuffle(n12, n13, [2]uint64{p0, p0h})
	n22 := Simd_i8x16_shuffle(n7, n21, [2]uint64{p5, p5h})
	n23 := Simd_i8x16_shuffle(n22, [2]uint64{p3, p3h}, [2]uint64{p0, p0h})
	n24 := Simd_i32x4_dot_i16x8_s(n23, [2]uint64{p4, p4h})
	n25 := Simd_i8x16_shuffle(n22, [2]uint64{p3, p3h}, [2]uint64{p1, p1h})
	n26 := Simd_i32x4_dot_i16x8_s(n25, [2]uint64{p4, p4h})
	n27 := Simd_i16x8_narrow_i32x4_s(n24, n26)
	n28 := Simd_i8x16_shuffle(n20, n27, [2]uint64{p6, p6h})
	n29 := Simd_i8x16_shuffle(n7, n21, [2]uint64{p2, p2h})
	n30 := Simd_i8x16_shuffle(n29, [2]uint64{p3, p3h}, [2]uint64{p0, p0h})
	n31 := Simd_i32x4_dot_i16x8_s(n30, [2]uint64{p4, p4h})
	n32 := Simd_i8x16_shuffle(n29, [2]uint64{p3, p3h}, [2]uint64{p1, p1h})
	n33 := Simd_i32x4_dot_i16x8_s(n32, [2]uint64{p4, p4h})
	n34 := Simd_i16x8_narrow_i32x4_s(n31, n33)
	return n20[0], n20[1], n27[0], n27[1], n28[0], n28[1], n34[0], n34[1]
}

//go:noinline
func Simd_p_fx104(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p3, p3h}, [2]uint64{p2, p2h})
	n2 := Simd_v128_load_nc(m, s0+64, 0)
	n3 := Simd_v128_load_nc(m, s0+80, 0)
	n4 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p4, p4h})
	n5 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p5, p5h})
	n6 := Simd_i8x16_shuffle(n4, n5, [2]uint64{p4, p4h})
	n7 := Simd_i8x16_shuffle(n4, n5, [2]uint64{p5, p5h})
	return n0[0], n0[1], n1[0], n1[1], n6[0], n6[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx105(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, [2]uint64{p3, p3h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, [2]uint64{p3, p3h}, [2]uint64{p0, p0h})
	n2 := Simd_v128_load_nc(m, s0+96, 0)
	n3 := Simd_v128_load_nc(m, s0+112, 0)
	n4 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p0, p0h})
	n5 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p1, p1h})
	n6 := Simd_i8x16_shuffle(n4, n5, [2]uint64{p0, p0h})
	n7 := Simd_i8x16_shuffle(n4, n5, [2]uint64{p1, p1h})
	n8 := Simd_i8x16_shuffle(n6, n7, [2]uint64{p1, p1h})
	n9 := Simd_i8x16_shuffle(n0, n8, [2]uint64{p4, p4h})
	n10 := Simd_i8x16_shuffle(n9, [2]uint64{p5, p5h}, [2]uint64{p0, p0h})
	n11 := Simd_i32x4_dot_i16x8_s(n10, [2]uint64{p6, p6h})
	n12 := Simd_i8x16_shuffle(n9, [2]uint64{p5, p5h}, [2]uint64{p1, p1h})
	n13 := Simd_i32x4_dot_i16x8_s(n12, [2]uint64{p6, p6h})
	n14 := Simd_i16x8_narrow_i32x4_s(n11, n13)
	n15 := Simd_i8x16_shuffle(n6, n7, [2]uint64{p0, p0h})
	return n14[0], n14[1], n1[0], n1[1], n15[0], n15[1]
}

//go:noinline
func Simd_p_fx106(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i8x16_avgr_u([2]uint64{p1, p1h}, n0)
	n2 := Simd_v128_load(m, s1, 0)
	n3 := Simd_i8x16_avgr_u([2]uint64{p0, p0h}, n2)
	return n3[0], n3[1], n1[0], n1[1]
}

//go:noinline
func Simd_p_fx107(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 64)
	n1 := Simd_v128_load_nc(m, s0+16, 0)
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p0, p0h})
	n3 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p1, p1h})
	n4 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p0, p0h})
	n5 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p1, p1h})
	n6 := Simd_i8x16_shuffle(n4, n5, [2]uint64{p1, p1h})
	n7 := Simd_i8x16_shuffle(n4, n5, [2]uint64{p0, p0h})
	n8 := Simd_v128_load_nc(m, s0+32, 0)
	n9 := Simd_v128_load_nc(m, s0+48, 0)
	n10 := Simd_i8x16_shuffle(n8, n9, [2]uint64{p0, p0h})
	n11 := Simd_i8x16_shuffle(n8, n9, [2]uint64{p1, p1h})
	n12 := Simd_i8x16_shuffle(n10, n11, [2]uint64{p0, p0h})
	n13 := Simd_i8x16_shuffle(n10, n11, [2]uint64{p1, p1h})
	n14 := Simd_i8x16_shuffle(n12, n13, [2]uint64{p1, p1h})
	n15 := Simd_i8x16_shuffle(n6, n14, [2]uint64{p2, p2h})
	n16 := Simd_i8x16_shuffle(n15, [2]uint64{p3, p3h}, [2]uint64{p0, p0h})
	n17 := Simd_i8x16_shuffle(n12, n13, [2]uint64{p0, p0h})
	return n15[0], n15[1], n16[0], n16[1], n7[0], n7[1], n17[0], n17[1]
}

//go:noinline
func Simd_p_fx108(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 96)
	n1 := Simd_v128_load_nc(m, s0, 48)
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p0, p0h})
	n3 := Simd_v128_load_nc(m, s0, 16)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx109(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p4, p4h}, [2]uint64{p5, p5h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_load_nc(m, s0, 64)
	n2 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n1, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, n2, [2]uint64{p3, p3h})
	return n1[0], n1[1], n2[0], n2[1], n3[0], n3[1], n0[0], n0[1]
}

//go:noinline
func Simd_p_fx110(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_nc(m, s0, 32)
	n1 := Simd_v128_load_nc(m, s0, 80)
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p0, p0h})
	n3 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, n2, [2]uint64{p2, p2h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx111(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p2, p2h})
	n3 := Simd_i8x16_shuffle([2]uint64{p6, p6h}, n2, [2]uint64{p5, p5h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx112(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n2 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, [2]uint64{p6, p6h}, [2]uint64{p2, p2h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx113(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p5, p5h})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p6, p6h}, [2]uint64{p2, p2h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx114(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p4, p4h}, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, [2]uint64{p6, p6h}, [2]uint64{p3, p3h})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p4, p4h}, [2]uint64{p3, p3h})
	_ = Simd_v128_store(m, s0+16, 0, [2]uint64{p0, p0h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx115(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p3, p3h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle(n0, [2]uint64{p3, p3h}, [2]uint64{p5, p5h})
	n3 := Simd_i8x16_shuffle([2]uint64{p4, p4h}, [2]uint64{p3, p3h}, [2]uint64{p5, p5h})
	n4 := Simd_i8x16_shuffle([2]uint64{p6, p6h}, [2]uint64{p3, p3h}, [2]uint64{p5, p5h})
	return n1[0], n1[1], n3[0], n3[1], n4[0], n4[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx116(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p5, p5h})
	n3 := Simd_i8x16_shuffle([2]uint64{p6, p6h}, n2, [2]uint64{p2, p2h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx117(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n2 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, [2]uint64{p6, p6h}, [2]uint64{p4, p4h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx118(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p2, p2h})
	n3 := Simd_i8x16_shuffle([2]uint64{p6, p6h}, n2, [2]uint64{p2, p2h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx119(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 128)
	n1 := Simd_v128_load_nc(m, s0, 16)
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p0, p0h})
	n3 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p1, p1h})
	n4 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p0, p0h})
	n5 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p1, p1h})
	n6 := Simd_v128_load_nc(m, s0, 32)
	n7 := Simd_v128_load_nc(m, s0, 48)
	n8 := Simd_i8x16_shuffle(n6, n7, [2]uint64{p0, p0h})
	n9 := Simd_i8x16_shuffle(n6, n7, [2]uint64{p1, p1h})
	n10 := Simd_i8x16_shuffle(n8, n9, [2]uint64{p0, p0h})
	n11 := Simd_i8x16_shuffle(n8, n9, [2]uint64{p1, p1h})
	n12 := Simd_i8x16_shuffle(n5, n11, [2]uint64{p2, p2h})
	n13 := Simd_i8x16_shuffle(n4, n10, [2]uint64{p2, p2h})
	n14 := Simd_i8x16_shuffle(n4, n10, [2]uint64{p3, p3h})
	n15 := Simd_i8x16_shuffle(n14, n12, [2]uint64{p0, p0h})
	n16 := Simd_i8x16_shuffle(n14, n12, [2]uint64{p1, p1h})
	n17 := Simd_i8x16_shuffle(n13, n14, [2]uint64{p0, p0h})
	n18 := Simd_i8x16_shuffle(n13, n14, [2]uint64{p1, p1h})
	return n17[0], n17[1], n15[0], n15[1], n18[0], n18[1], n16[0], n16[1]
}

//go:noinline
func Simd_p_fx120(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_nc(m, s0, 64)
	n1 := Simd_v128_load_nc(m, s0, 80)
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p0, p0h})
	n3 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p1, p1h})
	n4 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p0, p0h})
	n5 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p1, p1h})
	n6 := Simd_v128_load_nc(m, s0, 96)
	n7 := Simd_v128_load_nc(m, s0, 112)
	n8 := Simd_i8x16_shuffle(n6, n7, [2]uint64{p0, p0h})
	n9 := Simd_i8x16_shuffle(n6, n7, [2]uint64{p1, p1h})
	n10 := Simd_i8x16_shuffle(n8, n9, [2]uint64{p0, p0h})
	n11 := Simd_i8x16_shuffle(n8, n9, [2]uint64{p1, p1h})
	n12 := Simd_i8x16_shuffle(n5, n11, [2]uint64{p2, p2h})
	n13 := Simd_i8x16_shuffle(n4, n10, [2]uint64{p2, p2h})
	n14 := Simd_i8x16_shuffle(n4, n10, [2]uint64{p3, p3h})
	n15 := Simd_i8x16_shuffle(n14, n12, [2]uint64{p0, p0h})
	n16 := Simd_i8x16_shuffle(n14, n12, [2]uint64{p1, p1h})
	n17 := Simd_i8x16_shuffle(n13, n14, [2]uint64{p0, p0h})
	n18 := Simd_i8x16_shuffle(n13, n14, [2]uint64{p1, p1h})
	return n17[0], n17[1], n15[0], n15[1], n18[0], n18[1], n16[0], n16[1]
}

//go:noinline
func Simd_p_fx121(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, s5 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 17)
	n1 := Simd_scalar_i32_add(s1, s2)
	n2 := Simd_v128_load_rng(m, n1, 0, 0, 17)
	n3 := Simd_v128_load_nc(m, s0, 1)
	n4 := Simd_v128_xor(n3, n2)
	n5 := Simd_i8x16_avgr_u(n2, n3)
	n6 := Simd_scalar_i32_add(s1, s2)
	n7 := Simd_v128_load_nc(m, n6, 1)
	n8 := Simd_i8x16_avgr_u(n0, n7)
	n9 := Simd_v128_xor(n0, n7)
	n10 := Simd_v128_xor(n5, n8)
	n11 := Simd_v128_and(n10, n4)
	n12 := Simd_v128_or(n4, n9)
	n13 := Simd_v128_or(n12, n10)
	n14 := Simd_v128_and(n13, [2]uint64{p0, p0h})
	n15 := Simd_v128_and(n10, n9)
	n16 := Simd_i8x16_avgr_u(n5, n8)
	n17 := Simd_i8x16_sub(n16, n14)
	n18 := Simd_i8x16_avgr_u(n17, n5)
	n19 := Simd_v128_xor(n17, n5)
	n20 := Simd_v128_or(n19, n11)
	n21 := Simd_v128_and(n20, [2]uint64{p0, p0h})
	n22 := Simd_i8x16_sub(n18, n21)
	n23 := Simd_i8x16_avgr_u(n0, n22)
	n24 := Simd_i8x16_avgr_u(n7, n22)
	n25 := Simd_i8x16_avgr_u(n17, n8)
	n26 := Simd_v128_xor(n17, n8)
	n27 := Simd_v128_or(n26, n15)
	n28 := Simd_v128_and(n27, [2]uint64{p0, p0h})
	n29 := Simd_i8x16_sub(n25, n28)
	n30 := Simd_i8x16_avgr_u(n3, n29)
	n31 := Simd_i8x16_avgr_u(n2, n29)
	n32 := Simd_i8x16_shuffle(n31, n24, [2]uint64{p1, p1h})
	n33 := Simd_i8x16_shuffle(n31, n24, [2]uint64{p2, p2h})
	n34 := Simd_i8x16_shuffle(n23, n30, [2]uint64{p1, p1h})
	n35 := Simd_i8x16_shuffle(n23, n30, [2]uint64{p2, p2h})
	_ = Simd_v128_store(m, s3, 80, n34)
	_ = Simd_v128_store(m, s3, 64, n35)
	_ = Simd_v128_store(m, s3, 16, n32)
	_ = Simd_v128_store(m, s3, 0, n33)
	n40 := Simd_scalar_i32_add(s4, s2)
	n41 := Simd_v128_load_rng(m, n40, 0, 0, 17)
	n42 := Simd_scalar_i32_add(s5, s2)
	n43 := Simd_v128_load_rng(m, n42, 0, 0, 17)
	n44 := Simd_scalar_i32_add(s4, s2)
	n45 := Simd_v128_load_nc(m, n44, 1)
	n46 := Simd_v128_xor(n45, n43)
	n47 := Simd_i8x16_avgr_u(n43, n45)
	n48 := Simd_scalar_i32_add(s5, s2)
	n49 := Simd_v128_load_nc(m, n48, 1)
	n50 := Simd_i8x16_avgr_u(n41, n49)
	n51 := Simd_v128_xor(n41, n49)
	n52 := Simd_v128_xor(n47, n50)
	n53 := Simd_v128_and(n52, n46)
	n54 := Simd_v128_or(n46, n51)
	n55 := Simd_v128_or(n54, n52)
	n56 := Simd_v128_and(n55, [2]uint64{p0, p0h})
	n57 := Simd_v128_and(n52, n51)
	n58 := Simd_i8x16_avgr_u(n47, n50)
	n59 := Simd_i8x16_sub(n58, n56)
	n60 := Simd_i8x16_avgr_u(n59, n47)
	n61 := Simd_v128_xor(n59, n47)
	n62 := Simd_v128_or(n61, n53)
	n63 := Simd_v128_and(n62, [2]uint64{p0, p0h})
	n64 := Simd_i8x16_sub(n60, n63)
	n65 := Simd_i8x16_avgr_u(n41, n64)
	n66 := Simd_i8x16_avgr_u(n49, n64)
	n67 := Simd_i8x16_avgr_u(n59, n50)
	n68 := Simd_v128_xor(n59, n50)
	n69 := Simd_v128_or(n68, n57)
	n70 := Simd_v128_and(n69, [2]uint64{p0, p0h})
	n71 := Simd_i8x16_sub(n67, n70)
	n72 := Simd_i8x16_avgr_u(n45, n71)
	n73 := Simd_i8x16_avgr_u(n43, n71)
	n74 := Simd_i8x16_shuffle(n73, n66, [2]uint64{p1, p1h})
	n75 := Simd_i8x16_shuffle(n73, n66, [2]uint64{p2, p2h})
	n76 := Simd_i8x16_shuffle(n65, n72, [2]uint64{p1, p1h})
	n77 := Simd_i8x16_shuffle(n65, n72, [2]uint64{p2, p2h})
	_ = Simd_v128_store(m, s3, 112, n76)
	_ = Simd_v128_store(m, s3, 96, n77)
	_ = Simd_v128_store(m, s3, 48, n74)
	_ = Simd_v128_store(m, s3, 32, n75)
	return
}

//go:noinline
func Simd_p_fx122(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 49)
	n1 := Simd_v128_load_nc(m, s0, 32)
	n2 := Simd_v128_load_nc(m, s0, 1)
	n3 := Simd_v128_xor(n2, n1)
	n4 := Simd_i8x16_avgr_u(n1, n2)
	n5 := Simd_v128_load_nc(m, s0, 33)
	n6 := Simd_i8x16_avgr_u(n0, n5)
	n7 := Simd_v128_xor(n0, n5)
	n8 := Simd_v128_xor(n4, n6)
	n9 := Simd_v128_and(n8, n3)
	n10 := Simd_v128_or(n3, n7)
	n11 := Simd_v128_or(n10, n8)
	n12 := Simd_v128_and(n11, [2]uint64{p0, p0h})
	n13 := Simd_v128_and(n8, n7)
	n14 := Simd_i8x16_avgr_u(n4, n6)
	n15 := Simd_i8x16_sub(n14, n12)
	n16 := Simd_i8x16_avgr_u(n15, n4)
	n17 := Simd_v128_xor(n15, n4)
	n18 := Simd_v128_or(n17, n9)
	n19 := Simd_v128_and(n18, [2]uint64{p0, p0h})
	n20 := Simd_i8x16_sub(n16, n19)
	n21 := Simd_i8x16_avgr_u(n0, n20)
	n22 := Simd_i8x16_avgr_u(n5, n20)
	n23 := Simd_i8x16_avgr_u(n15, n6)
	n24 := Simd_v128_xor(n15, n6)
	n25 := Simd_v128_or(n24, n13)
	n26 := Simd_v128_and(n25, [2]uint64{p0, p0h})
	n27 := Simd_i8x16_sub(n23, n26)
	n28 := Simd_i8x16_avgr_u(n2, n27)
	n29 := Simd_i8x16_avgr_u(n1, n27)
	n30 := Simd_i8x16_shuffle(n29, n22, [2]uint64{p1, p1h})
	n31 := Simd_i8x16_shuffle(n29, n22, [2]uint64{p2, p2h})
	n32 := Simd_i8x16_shuffle(n21, n28, [2]uint64{p1, p1h})
	n33 := Simd_i8x16_shuffle(n21, n28, [2]uint64{p2, p2h})
	_ = Simd_v128_store(m, s1, 80, n32)
	_ = Simd_v128_store(m, s1, 64, n33)
	_ = Simd_v128_store(m, s1, 16, n30)
	_ = Simd_v128_store(m, s1, 0, n31)
	return
}

//go:noinline
func Simd_p_fx123(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 49)
	n1 := Simd_v128_load_nc(m, s0, 32)
	n2 := Simd_v128_load_nc(m, s0, 1)
	n3 := Simd_v128_xor(n2, n1)
	n4 := Simd_i8x16_avgr_u(n1, n2)
	n5 := Simd_v128_load_nc(m, s0, 33)
	n6 := Simd_i8x16_avgr_u(n0, n5)
	n7 := Simd_v128_xor(n0, n5)
	n8 := Simd_v128_xor(n4, n6)
	n9 := Simd_v128_and(n8, n3)
	n10 := Simd_v128_or(n3, n7)
	n11 := Simd_v128_or(n10, n8)
	n12 := Simd_v128_and(n11, [2]uint64{p0, p0h})
	n13 := Simd_v128_and(n8, n7)
	n14 := Simd_i8x16_avgr_u(n4, n6)
	n15 := Simd_i8x16_sub(n14, n12)
	n16 := Simd_i8x16_avgr_u(n15, n4)
	n17 := Simd_v128_xor(n15, n4)
	n18 := Simd_v128_or(n17, n9)
	n19 := Simd_v128_and(n18, [2]uint64{p0, p0h})
	n20 := Simd_i8x16_sub(n16, n19)
	n21 := Simd_i8x16_avgr_u(n0, n20)
	n22 := Simd_i8x16_avgr_u(n5, n20)
	n23 := Simd_i8x16_avgr_u(n15, n6)
	n24 := Simd_v128_xor(n15, n6)
	n25 := Simd_v128_or(n24, n13)
	n26 := Simd_v128_and(n25, [2]uint64{p0, p0h})
	n27 := Simd_i8x16_sub(n23, n26)
	n28 := Simd_i8x16_avgr_u(n2, n27)
	n29 := Simd_i8x16_avgr_u(n1, n27)
	n30 := Simd_i8x16_shuffle(n29, n22, [2]uint64{p1, p1h})
	n31 := Simd_i8x16_shuffle(n29, n22, [2]uint64{p2, p2h})
	n32 := Simd_i8x16_shuffle(n21, n28, [2]uint64{p1, p1h})
	n33 := Simd_i8x16_shuffle(n21, n28, [2]uint64{p2, p2h})
	_ = Simd_v128_store(m, s1, 112, n32)
	_ = Simd_v128_store(m, s1, 96, n33)
	_ = Simd_v128_store(m, s1, 48, n30)
	_ = Simd_v128_store(m, s1, 32, n31)
	return
}

//go:noinline
func Simd_p_fx124(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	_ = Simd_v128_store(m, s1+32, 0, n0)
	n2 := Simd_v128_load(m, s2, 0)
	_ = Simd_v128_store(m, s1+16, 0, n2)
	n4 := Simd_v128_load(m, s3, 0)
	_ = Simd_v128_store(m, s1, 0, n4)
	n6 := Simd_v128_load(m, s4, 0)
	_ = Simd_v128_store(m, s1+48, 0, n6)
	return
}

//go:noinline
func Simd_p_fx125(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_v128_load32_zero(m, s0, 0)
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n2 := Simd_v128_load32_zero(m, s1, 0)
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n4 := Simd_i16x8_add(n1, n3)
	n5 := Simd_i16x8_shr_u(n4, 1)
	n6 := Simd_v128_load32_zero(m, s0+-4, 0)
	n7 := Simd_i8x16_shuffle(n6, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n8 := Simd_i16x8_sub(n5, n7)
	n9 := Simd_i16x8_lt_s(n5, n7)
	n10 := Simd_i16x8_sub(n8, n9)
	n11 := Simd_i16x8_shr_s(n10, 1)
	n12 := Simd_i16x8_add(n11, n5)
	return n12[0], n12[1]
}

//go:noinline
func Simd_p_fx126(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_v128_load32_zero(m, s0, 0)
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n2 := Simd_v128_load32_zero(m, s1, 0)
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n4 := Simd_i16x8_add(n1, n3)
	n5 := Simd_v128_load32_zero(m, s0+-4, 0)
	n6 := Simd_i8x16_shuffle(n5, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n7 := Simd_i16x8_sub(n4, n6)
	return n7[0], n7[1]
}

//go:noinline
func Simd_p_fx127(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_v128_load32_zero(m, s0+-4, 0)
	n1 := Simd_i8x16_sub_sat_u(n0, [2]uint64{p0, p0h})
	n2 := Simd_i8x16_sub_sat_u([2]uint64{p0, p0h}, n0)
	n3 := Simd_v128_or(n1, n2)
	n4 := Simd_i8x16_shuffle(n3, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n5 := Simd_i8x16_sub_sat_u(n0, [2]uint64{p3, p3h})
	n6 := Simd_i8x16_sub_sat_u([2]uint64{p3, p3h}, n0)
	n7 := Simd_v128_or(n5, n6)
	n8 := Simd_i8x16_shuffle(n7, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n9 := Simd_i16x8_sub(n4, n8)
	return n9[0], n9[1]
}

//go:noinline
func Simd_p_fx128(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_v128_load32_zero(m, s0, 4)
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n2 := Simd_v128_load32_zero(m, s0, 0)
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n4 := Simd_i16x8_add(n1, n3)
	n5 := Simd_i16x8_shr_u(n4, 1)
	n6 := Simd_v128_load32_zero(m, s0+-4, 0)
	n7 := Simd_i8x16_shuffle(n6, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n8 := Simd_v128_load32_zero(m, s1, 0)
	n9 := Simd_i8x16_shuffle(n8, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n10 := Simd_i16x8_add(n7, n9)
	n11 := Simd_i16x8_shr_u(n10, 1)
	n12 := Simd_v128_and(n5, n11)
	n13 := Simd_v128_xor(n5, n11)
	n14 := Simd_i16x8_shr_u(n13, 1)
	n15 := Simd_i16x8_add(n12, n14)
	return n15[0], n15[1]
}

//go:noinline
func Simd_p_fx129(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_avgr_u([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_xor([2]uint64{p1, p1h}, [2]uint64{p0, p0h})
	n2 := Simd_v128_and(n1, [2]uint64{p2, p2h})
	n3 := Simd_i8x16_sub(n0, n2)
	return n3[0], n3[1]
}

//go:noinline
func Simd_p_fx130(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_v128_load32_zero(m, s0, 4)
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n2 := Simd_v128_load32_zero(m, s1, 0)
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n4 := Simd_i16x8_add(n1, n3)
	n5 := Simd_i16x8_shr_u(n4, 1)
	n6 := Simd_v128_load32_zero(m, s0, 0)
	n7 := Simd_i8x16_shuffle(n6, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n8 := Simd_i16x8_add(n5, n7)
	n9 := Simd_i16x8_shr_u(n8, 1)
	return n9[0], n9[1]
}

//go:noinline
func Simd_p_fx131(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32) {
	n0 := Simd_v128_load(m, s0+4, 0)
	n1 := Simd_scalar_i32_add(s1, s2)
	n2 := Simd_v128_load(m, n1, 0)
	n3 := Simd_i8x16_add(n0, n2)
	_ = Simd_v128_store(m, s3, 0, n3)
	n5 := Simd_v128_load(m, s0+20, 0)
	n6 := Simd_v128_load(m, s4+16, 0)
	n7 := Simd_i8x16_add(n5, n6)
	_ = Simd_v128_store(m, s3+16, 0, n7)
	return
}

//go:noinline
func Simd_p_fx132(m *Module, s0 int32, s1 int32, s2 int32, s3 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_scalar_i32_add(s1, s2)
	n2 := Simd_v128_load(m, n1, 0)
	n3 := Simd_i8x16_add(n0, n2)
	n4 := Simd_scalar_i32_add(s3, s2)
	_ = Simd_v128_store(m, n4, 0, n3)
	return
}

//go:noinline
func Simd_p_fx133(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1, 0)
	n2 := Simd_i8x16_add(n0, n1)
	_ = Simd_v128_store(m, s2, 0, n2)
	n4 := Simd_v128_load(m, s0+16, 0)
	n5 := Simd_v128_load(m, s1+16, 0)
	n6 := Simd_i8x16_add(n4, n5)
	_ = Simd_v128_store(m, s2+16, 0, n6)
	return
}

//go:noinline
func Simd_p_fx134(m *Module, s0 int32, s1 int32, s2 int32, s3 int32) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_scalar_i32_add(s2, s1)
	n3 := Simd_v128_load(m, n2, 0)
	n4 := Simd_i8x16_add(n1, n3)
	n5 := Simd_scalar_i32_add(s3, s1)
	_ = Simd_v128_store(m, n5, 0, n4)
	return
}

//go:noinline
func Simd_p_fx135(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n0, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_add(n1, n0)
	n3 := Simd_i8x16_add(n2, [2]uint64{p2, p2h})
	n4 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n2, [2]uint64{p3, p3h})
	n5 := Simd_i8x16_add(n3, n4)
	n6 := Simd_i8x16_shuffle(n5, [2]uint64{p0, p0h}, [2]uint64{p4, p4h})
	_ = Simd_v128_store(m, s1, 0, n5)
	n8 := Simd_v128_load(m, s0+16, 0)
	n9 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n8, [2]uint64{p1, p1h})
	n10 := Simd_i8x16_add(n9, n8)
	n11 := Simd_i8x16_add(n10, n6)
	n12 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n10, [2]uint64{p3, p3h})
	n13 := Simd_i8x16_add(n11, n12)
	n14 := Simd_i8x16_shuffle(n13, [2]uint64{p0, p0h}, [2]uint64{p4, p4h})
	_ = Simd_v128_store(m, s1+16, 0, n13)
	return n14[0], n14[1]
}

//go:noinline
func Simd_p_fx136(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n1, [2]uint64{1374179596971150604, 1952900979675763988})
	n3 := Simd_i8x16_add(n2, n1)
	n4 := Simd_i8x16_add(n3, [2]uint64{p1, p1h})
	n5 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n3, [2]uint64{1084818905618843912, 1663540288323457296})
	n6 := Simd_i8x16_add(n4, n5)
	n7 := Simd_scalar_i32_add(s2, s1)
	_ = Simd_v128_store(m, n7, 0, n6)
	return
}

//go:noinline
func Simd_p_fx137(m *Module, s0 int32, s1 int32, p0, p0h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i8x16_add(n0, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s1, 0, n1)
	n3 := Simd_v128_load(m, s0+16, 0)
	n4 := Simd_i8x16_add(n3, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s1+16, 0, n4)
	n6 := Simd_v128_load(m, s0+32, 0)
	n7 := Simd_i8x16_add(n6, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s1+32, 0, n7)
	n9 := Simd_v128_load(m, s0+48, 0)
	n10 := Simd_i8x16_add(n9, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s1+48, 0, n10)
	return
}

//go:noinline
func Simd_p_fx138(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_i8x16_add(n1, [2]uint64{p0, p0h})
	n3 := Simd_scalar_i32_add(s2, s1)
	_ = Simd_v128_store(m, n3, 0, n2)
	return
}

//go:noinline
func Simd_p_fx139(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_load32_zero(m, s0, 0)
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n3 := Simd_i16x8_add(n2, n0)
	n4 := Simd_i16x8_shr_u(n3, 1)
	n5 := Simd_v128_load32_zero(m, s1, 0)
	n6 := Simd_i8x16_shuffle(n5, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n7 := Simd_i16x8_add(n4, n6)
	n8 := Simd_i16x8_shr_u(n7, 1)
	return n8[0], n8[1]
}

//go:noinline
func Simd_p_fx140(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32) {
	n0 := Simd_v128_load(m, s0+-4, 0)
	n1 := Simd_scalar_i32_add(s1, s2)
	n2 := Simd_v128_load(m, n1, 0)
	n3 := Simd_i8x16_add(n0, n2)
	_ = Simd_v128_store(m, s3, 0, n3)
	n5 := Simd_v128_load(m, s0+12, 0)
	n6 := Simd_v128_load(m, s4+16, 0)
	n7 := Simd_i8x16_add(n5, n6)
	_ = Simd_v128_store(m, s3+16, 0, n7)
	return
}

//go:noinline
func Simd_p_fx141(m *Module, s0 int32, s1 int32, s2 int32, s3 int32) {
	n0 := Simd_scalar_i32_add(s0, -4)
	n1 := Simd_scalar_i32_add(n0, s1)
	n2 := Simd_v128_load(m, n1, 0)
	n3 := Simd_scalar_i32_add(s2, s1)
	n4 := Simd_v128_load(m, n3, 0)
	n5 := Simd_i8x16_add(n2, n4)
	n6 := Simd_scalar_i32_add(s3, s1)
	_ = Simd_v128_store(m, n6, 0, n5)
	return
}

//go:noinline
func Simd_p_fx142(m *Module, s0 int32) (uint64, uint64) {
	n0 := Simd_v128_load32_zero(m, s0+-4, 0)
	return n0[0], n0[1]
}

//go:noinline
func Simd_p_fx143(m *Module, s0 int32, s1 int32, p0, p0h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load_rng(m, s1+4, 0, -4, 20)
	n2 := Simd_v128_load_nc(m, s1, 0)
	n3 := Simd_i8x16_shuffle(n2, n1, [2]uint64{p0, p0h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx144(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_sub_sat_u(n0, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_sub_sat_u([2]uint64{p3, p3h}, n0)
	n3 := Simd_v128_or(n2, n1)
	n4 := Simd_i16x8_shr_u(n3, 8)
	n5 := Simd_v128_and(n3, [2]uint64{p4, p4h})
	n6 := Simd_i16x8_add(n4, n5)
	n7 := Simd_i32x4_shl(n6, 16)
	n8 := Simd_i16x8_add(n6, n7)
	n9 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n10 := Simd_i8x16_sub_sat_u(n9, [2]uint64{p3, p3h})
	n11 := Simd_i8x16_sub_sat_u([2]uint64{p3, p3h}, n9)
	n12 := Simd_v128_or(n11, n10)
	n13 := Simd_i16x8_shr_u(n12, 8)
	n14 := Simd_v128_and(n12, [2]uint64{p4, p4h})
	n15 := Simd_i16x8_add(n13, n14)
	n16 := Simd_i32x4_shl(n15, 16)
	n17 := Simd_i16x8_add(n15, n16)
	n18 := Simd_i64x2_shl(n17, 32)
	n19 := Simd_i16x8_add(n17, n18)
	n20 := Simd_i64x2_shr_u(n19, 48)
	n21 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, [2]uint64{p1, p1h}, [2]uint64{p6, p6h})
	n22 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p1, p1h}, [2]uint64{p6, p6h})
	n23 := Simd_i8x16_sub_sat_u(n22, n21)
	n24 := Simd_i8x16_sub_sat_u(n21, n22)
	n25 := Simd_v128_or(n24, n23)
	n26 := Simd_i16x8_shr_u(n25, 8)
	n27 := Simd_v128_and(n25, [2]uint64{p4, p4h})
	n28 := Simd_i16x8_add(n26, n27)
	n29 := Simd_i32x4_shl(n28, 16)
	n30 := Simd_i16x8_add(n28, n29)
	n31 := Simd_i64x2_shl(n30, 32)
	n32 := Simd_i16x8_add(n30, n31)
	n33 := Simd_i64x2_shr_u(n32, 48)
	n34 := Simd_i16x8_narrow_i32x4_s(n20, n33)
	return n8[0], n8[1], n34[0], n34[1]
}

//go:noinline
func Simd_p_fx145(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i64x2_shl([2]uint64{p3, p3h}, 32)
	n1 := Simd_i16x8_add([2]uint64{p3, p3h}, n0)
	n2 := Simd_i64x2_shr_u(n1, 48)
	n3 := Simd_i32x4_gt_s(n2, [2]uint64{p4, p4h})
	n4 := Simd_v128_bitselect([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, n3)
	n5 := Simd_i8x16_add([2]uint64{p0, p0h}, n4)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx146(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle([2]uint64{p4, p4h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n3 := Simd_i8x16_shuffle(n2, n1, [2]uint64{p5, p5h})
	n4 := Simd_i8x16_shuffle([2]uint64{p6, p6h}, n1, [2]uint64{p5, p5h})
	n5 := Simd_i8x16_sub_sat_u(n4, n3)
	n6 := Simd_i8x16_sub_sat_u(n3, n4)
	n7 := Simd_v128_or(n6, n5)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx147(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_shr_u([2]uint64{p0, p0h}, 8)
	n1 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n2 := Simd_i16x8_add(n0, n1)
	n3 := Simd_i32x4_shl(n2, 16)
	n4 := Simd_i16x8_add(n2, n3)
	n5 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	return n4[0], n4[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx148(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle(n1, n0, [2]uint64{p4, p4h})
	n3 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, n0, [2]uint64{p4, p4h})
	n4 := Simd_i8x16_sub_sat_u(n3, n2)
	n5 := Simd_i8x16_sub_sat_u(n2, n3)
	n6 := Simd_v128_or(n5, n4)
	n7 := Simd_i16x8_shr_u(n6, 8)
	n8 := Simd_v128_and(n6, [2]uint64{p6, p6h})
	n9 := Simd_i16x8_add(n7, n8)
	n10 := Simd_i32x4_shl(n9, 16)
	n11 := Simd_i16x8_add(n9, n10)
	return n0[0], n0[1], n11[0], n11[1]
}

//go:noinline
func Simd_p_fx149(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i64x2_shl([2]uint64{p5, p5h}, 32)
	n2 := Simd_i16x8_add([2]uint64{p5, p5h}, n1)
	n3 := Simd_i64x2_shr_u(n2, 48)
	n4 := Simd_i8x16_shuffle([2]uint64{p6, p6h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n5 := Simd_i32x4_gt_s(n3, n4)
	n6 := Simd_v128_bitselect([2]uint64{p3, p3h}, [2]uint64{p4, p4h}, n5)
	n7 := Simd_i8x16_add(n0, n6)
	return n7[0], n7[1]
}

//go:noinline
func Simd_p_fx150(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1374179596769034496, 216736831629295872})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p2, p2h}, [2]uint64{506097522914230528, 216736831898784016})
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p3, p3h}, [2]uint64{506097522914230528, 1374179596903778568})
	_ = Simd_v128_store(m, s0, 0, n2)
	return
}

//go:noinline
func Simd_p_fx151(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0+4, 0, -4, 24)
	n1 := Simd_v128_load_nc(m, s0+8, 0)
	n2 := Simd_v128_xor(n1, n0)
	n3 := Simd_v128_and(n2, [2]uint64{p0, p0h})
	n4 := Simd_i8x16_avgr_u(n0, n1)
	n5 := Simd_i8x16_sub(n4, n3)
	n6 := Simd_i8x16_shuffle(n5, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n7 := Simd_i8x16_shuffle(n6, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n8 := Simd_i8x16_shuffle(n7, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n9 := Simd_v128_load_nc(m, s0, 0)
	n10 := Simd_v128_xor(n9, [2]uint64{p1, p1h})
	n11 := Simd_v128_and(n10, [2]uint64{p0, p0h})
	n12 := Simd_i8x16_avgr_u([2]uint64{p1, p1h}, n9)
	n13 := Simd_i8x16_sub(n12, n11)
	n14 := Simd_i8x16_shuffle(n9, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n15 := Simd_i8x16_shuffle(n14, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n16 := Simd_i8x16_shuffle(n15, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n17 := Simd_i8x16_avgr_u(n5, n13)
	n18 := Simd_v128_xor(n5, n13)
	n19 := Simd_v128_and(n18, [2]uint64{p0, p0h})
	n20 := Simd_i8x16_sub(n17, n19)
	n21 := Simd_v128_load(m, s1, 0)
	n22 := Simd_i8x16_shuffle(n21, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n23 := Simd_i8x16_shuffle(n22, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n24 := Simd_i8x16_shuffle(n23, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n25 := Simd_i8x16_add(n20, n21)
	n26 := Simd_i8x16_avgr_u(n25, n14)
	n27 := Simd_v128_xor(n25, n14)
	n28 := Simd_v128_and(n27, [2]uint64{p0, p0h})
	n29 := Simd_i8x16_sub(n26, n28)
	n30 := Simd_v128_xor(n29, n6)
	n31 := Simd_v128_and(n30, [2]uint64{p0, p0h})
	n32 := Simd_i8x16_avgr_u(n6, n29)
	n33 := Simd_i8x16_sub(n32, n31)
	n34 := Simd_i8x16_add(n33, n22)
	n35 := Simd_i8x16_avgr_u(n34, n15)
	n36 := Simd_v128_xor(n34, n15)
	n37 := Simd_v128_and(n36, [2]uint64{p0, p0h})
	n38 := Simd_i8x16_sub(n35, n37)
	n39 := Simd_v128_xor(n38, n7)
	n40 := Simd_v128_and(n39, [2]uint64{p0, p0h})
	n41 := Simd_i8x16_avgr_u(n7, n38)
	n42 := Simd_i8x16_sub(n41, n40)
	n43 := Simd_i8x16_add(n42, n23)
	n44 := Simd_i8x16_avgr_u(n43, n16)
	n45 := Simd_v128_xor(n43, n16)
	n46 := Simd_v128_and(n45, [2]uint64{p0, p0h})
	n47 := Simd_i8x16_sub(n44, n46)
	n48 := Simd_v128_xor(n47, n8)
	n49 := Simd_v128_and(n48, [2]uint64{p0, p0h})
	n50 := Simd_i8x16_avgr_u(n8, n47)
	n51 := Simd_i8x16_add(n50, n24)
	n52 := Simd_i8x16_sub(n51, n49)
	n53 := Simd_i8x16_shuffle(n25, n34, [2]uint64{1374179596769034496, 216736831629295872})
	n54 := Simd_i8x16_shuffle(n53, n43, [2]uint64{506097522914230528, 216736831898784016})
	n55 := Simd_i8x16_shuffle(n54, n52, [2]uint64{506097522914230528, 1374179596903778568})
	_ = Simd_v128_store(m, s2, 0, n55)
	return n52[0], n52[1]
}

//go:noinline
func Simd_p_fx152(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 20)
	n1 := Simd_v128_load_nc(m, s0+4, 0)
	n2 := Simd_i8x16_avgr_u(n0, n1)
	n3 := Simd_v128_xor(n0, n1)
	n4 := Simd_v128_and(n3, [2]uint64{p0, p0h})
	n5 := Simd_v128_load(m, s1, 0)
	n6 := Simd_i8x16_add(n2, n5)
	n7 := Simd_i8x16_sub(n6, n4)
	_ = Simd_v128_store(m, s2, 0, n7)
	return
}

//go:noinline
func Simd_p_fx153(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, p0, p0h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s0+-4, 0)
	n2 := Simd_i8x16_avgr_u(n0, n1)
	n3 := Simd_v128_xor(n0, n1)
	n4 := Simd_v128_and(n3, [2]uint64{p0, p0h})
	n5 := Simd_scalar_i32_add(s1, s2)
	n6 := Simd_v128_load(m, n5, 0)
	n7 := Simd_i8x16_add(n2, n6)
	n8 := Simd_i8x16_sub(n7, n4)
	_ = Simd_v128_store(m, s3, 0, n8)
	n10 := Simd_v128_load_rng(m, s0+16, 0, -4, 20)
	n11 := Simd_v128_load_nc(m, s0+12, 0)
	n12 := Simd_i8x16_avgr_u(n10, n11)
	n13 := Simd_v128_xor(n10, n11)
	n14 := Simd_v128_and(n13, [2]uint64{p0, p0h})
	n15 := Simd_v128_load(m, s4+16, 0)
	n16 := Simd_i8x16_add(n12, n15)
	n17 := Simd_i8x16_sub(n16, n14)
	_ = Simd_v128_store(m, s3+16, 0, n17)
	return
}

//go:noinline
func Simd_p_fx154(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, p0, p0h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s0+-4, 0)
	n2 := Simd_i8x16_avgr_u(n0, n1)
	n3 := Simd_v128_xor(n0, n1)
	n4 := Simd_v128_and(n3, [2]uint64{p0, p0h})
	n5 := Simd_scalar_i32_add(s1, s2)
	n6 := Simd_v128_load(m, n5, 0)
	n7 := Simd_i8x16_add(n2, n6)
	n8 := Simd_i8x16_sub(n7, n4)
	n9 := Simd_scalar_i32_add(s3, s2)
	_ = Simd_v128_store(m, n9, 0, n8)
	return
}

//go:noinline
func Simd_p_fx155(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_load32_zero(m, s0, 0)
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n3 := Simd_i16x8_add(n2, n0)
	n4 := Simd_i16x8_shr_u(n3, 1)
	n5 := Simd_v128_load32_zero(m, s1, 0)
	n6 := Simd_i8x16_shuffle(n5, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n7 := Simd_i16x8_sub(n4, n6)
	n8 := Simd_i16x8_lt_s(n4, n6)
	n9 := Simd_i16x8_sub(n7, n8)
	n10 := Simd_i16x8_shr_s(n9, 1)
	n11 := Simd_i16x8_add(n10, n4)
	return n11[0], n11[1]
}

//go:noinline
func Simd_p_fx156(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_v128_load32_zero(m, s0+-4, 0)
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p0, p0h}, [2]uint64{1369958511735279616, 1659319203087586308})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx157(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0+4, 0, -4, 20)
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_shuffle(n0, [2]uint64{p0, p0h}, [2]uint64{p5, p5h})
	n3 := Simd_v128_load_nc(m, s0, 0)
	n4 := Simd_i8x16_shuffle(n3, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n5 := Simd_i16x8_sub(n1, n4)
	n6 := Simd_i16x8_add(n5, [2]uint64{p2, p2h})
	n7 := Simd_i8x16_narrow_i16x8_u(n6, n6)
	n8 := Simd_i8x16_shuffle(n5, [2]uint64{p0, p0h}, [2]uint64{p3, p3h})
	n9 := Simd_i8x16_shuffle(n3, [2]uint64{p0, p0h}, [2]uint64{p5, p5h})
	n10 := Simd_i16x8_sub(n2, n9)
	n11 := Simd_i8x16_shuffle(n10, [2]uint64{p0, p0h}, [2]uint64{p3, p3h})
	n12 := Simd_v128_load(m, s1, 0)
	n13 := Simd_i8x16_shuffle(n12, [2]uint64{p0, p0h}, [2]uint64{p4, p4h})
	n14 := Simd_i8x16_shuffle(n13, [2]uint64{p0, p0h}, [2]uint64{p4, p4h})
	n15 := Simd_i8x16_shuffle(n14, [2]uint64{p0, p0h}, [2]uint64{p4, p4h})
	n16 := Simd_i8x16_add(n7, n12)
	n17 := Simd_i8x16_shuffle(n16, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n18 := Simd_i16x8_add(n17, n8)
	n19 := Simd_i8x16_narrow_i16x8_u(n18, n18)
	n20 := Simd_i8x16_add(n19, n13)
	n21 := Simd_i8x16_shuffle(n20, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n22 := Simd_i16x8_add(n10, n21)
	n23 := Simd_i8x16_narrow_i16x8_u(n22, n22)
	n24 := Simd_i8x16_add(n23, n14)
	n25 := Simd_i8x16_shuffle(n24, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n26 := Simd_i16x8_add(n25, n11)
	n27 := Simd_i8x16_narrow_i16x8_u(n26, n26)
	n28 := Simd_i8x16_add(n27, n15)
	n29 := Simd_i8x16_shuffle(n16, n20, [2]uint64{1374179596769034496, 216736831629295872})
	n30 := Simd_i8x16_shuffle(n29, n24, [2]uint64{506097522914230528, 216736831898784016})
	n31 := Simd_i8x16_shuffle(n30, n28, [2]uint64{506097522914230528, 1374179596903778568})
	_ = Simd_v128_store(m, s2, 0, n31)
	return n28[0], n28[1]
}

//go:noinline
func Simd_p_fx158(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i16x8_shr_u(n0, 8)
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n4 := Simd_i8x16_add(n3, n0)
	_ = Simd_v128_store(m, s1, 0, n4)
	n6 := Simd_v128_load(m, s0+16, 0)
	n7 := Simd_i16x8_shr_u(n6, 8)
	n8 := Simd_i8x16_shuffle(n7, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n9 := Simd_i8x16_shuffle(n8, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n10 := Simd_i8x16_add(n9, n6)
	_ = Simd_v128_store(m, s1+16, 0, n10)
	return
}

//go:noinline
func Simd_p_fx159(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_i16x8_shr_u(n1, 8)
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p0, p0h}, [2]uint64{361419384851267840, 1084818905618843912})
	n4 := Simd_i8x16_shuffle(n3, [2]uint64{p0, p0h}, [2]uint64{506097522914230528, 940140767555881224})
	n5 := Simd_i8x16_add(n4, n1)
	n6 := Simd_scalar_i32_add(s2, s1)
	_ = Simd_v128_store(m, n6, 0, n5)
	return
}

//go:noinline
func Simd_p_fx160(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i32x4_shr_u(n0, 8)
	n2 := Simd_v128_and(n1, [2]uint64{p0, p0h})
	n3 := Simd_i32x4_shl(n1, 16)
	n4 := Simd_v128_and(n0, [2]uint64{p1, p1h})
	n5 := Simd_i32x4_add(n2, n4)
	n6 := Simd_i32x4_add(n5, n3)
	n7 := Simd_v128_and(n6, [2]uint64{p1, p1h})
	n8 := Simd_v128_and(n0, [2]uint64{p2, p2h})
	n9 := Simd_v128_or(n7, n8)
	_ = Simd_v128_store(m, s1, 0, n9)
	return
}

//go:noinline
func Simd_p_fx161(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) {
	n0 := Simd_i32x4_extend_low_i16x8_s([2]uint64{p2, p2h})
	n1 := Simd_i32x4_extend_high_i16x8_s([2]uint64{p2, p2h})
	n2 := Simd_i32x4_extend_low_i16x8_s([2]uint64{p4, p4h})
	n3 := Simd_i32x4_extend_high_i16x8_s([2]uint64{p4, p4h})
	n4 := Simd_v128_load(m, s0, 0)
	n5 := Simd_v128_and(n4, [2]uint64{p0, p0h})
	n6 := Simd_i8x16_shuffle(n5, [2]uint64{p1, p1h}, [2]uint64{361419384851267840, 1084818905618843912})
	n7 := Simd_i8x16_shuffle(n6, [2]uint64{p1, p1h}, [2]uint64{506097522914230528, 940140767555881224})
	n8 := Simd_i32x4_extend_low_i16x8_s(n7)
	n9 := Simd_i32x4_mul(n8, n0)
	n10 := Simd_i32x4_extend_high_i16x8_s(n7)
	n11 := Simd_i32x4_mul(n10, n1)
	n12 := Simd_i8x16_shuffle(n9, n11, [2]uint64{p3, p3h})
	n13 := Simd_i8x16_add(n12, n4)
	n14 := Simd_i16x8_shl(n13, 8)
	n15 := Simd_i32x4_extend_low_i16x8_s(n14)
	n16 := Simd_i32x4_mul(n15, n2)
	n17 := Simd_i32x4_extend_high_i16x8_s(n14)
	n18 := Simd_i32x4_mul(n17, n3)
	n19 := Simd_i8x16_shuffle(n16, n18, [2]uint64{p3, p3h})
	n20 := Simd_i32x4_shr_u(n19, 8)
	n21 := Simd_i8x16_add(n20, n14)
	n22 := Simd_i16x8_shr_u(n21, 8)
	n23 := Simd_v128_or(n5, n22)
	_ = Simd_v128_store(m, s1, 0, n23)
	return
}

//go:noinline
func Simd_p_fx162(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_i32x4_splat(s1)
	n1 := Simd_i32x4_splat(s2)
	n2 := Simd_i32x4_splat(s3)
	n3 := Simd_v128_load(m, s0, 0)
	n4 := Simd_i32x4_shl(n3, 16)
	n5 := Simd_i32x4_shr_s(n4, 24)
	n6 := Simd_i32x4_mul(n5, n0)
	n7 := Simd_i32x4_shr_s(n6, 5)
	n8 := Simd_i32x4_mul(n5, n1)
	n9 := Simd_i32x4_shr_u(n8, 5)
	n10 := Simd_i32x4_add(n9, n3)
	n11 := Simd_i32x4_shr_u(n3, 16)
	n12 := Simd_i32x4_add(n7, n11)
	n13 := Simd_i32x4_shl(n12, 16)
	n14 := Simd_v128_and(n13, [2]uint64{p0, p0h})
	n15 := Simd_v128_and(n3, [2]uint64{p1, p1h})
	n16 := Simd_v128_or(n14, n15)
	n17 := Simd_i32x4_shl(n12, 24)
	n18 := Simd_i32x4_shr_s(n17, 24)
	n19 := Simd_i32x4_mul(n18, n2)
	n20 := Simd_i32x4_shr_u(n19, 5)
	n21 := Simd_i32x4_add(n10, n20)
	n22 := Simd_v128_and(n21, [2]uint64{p2, p2h})
	n23 := Simd_v128_or(n16, n22)
	_ = Simd_v128_store(m, s4, 0, n23)
	return
}

//go:noinline
func Simd_p_fx163(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 128)
	n1 := Simd_v128_load_nc(m, s0, 16)
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p0, p0h})
	n3 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p1, p1h})
	n4 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p0, p0h})
	n5 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p1, p1h})
	n6 := Simd_i8x16_shuffle(n4, n5, [2]uint64{p0, p0h})
	n7 := Simd_v128_load_nc(m, s0, 32)
	return n4[0], n4[1], n5[0], n5[1], n6[0], n6[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx164(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_nc(m, s0, 48)
	n1 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n0, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n0, [2]uint64{p2, p2h})
	n3 := Simd_i8x16_shuffle(n1, n2, [2]uint64{p1, p1h})
	n4 := Simd_i8x16_shuffle(n1, n2, [2]uint64{p2, p2h})
	n5 := Simd_i8x16_shuffle(n3, n4, [2]uint64{p1, p1h})
	n6 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n5, [2]uint64{p4, p4h})
	return n3[0], n3[1], n4[0], n4[1], n5[0], n5[1], n6[0], n6[1]
}

//go:noinline
func Simd_p_fx165(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_nc(m, s0, 64)
	n1 := Simd_v128_load_nc(m, s0, 80)
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p0, p0h})
	n3 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p1, p1h})
	n4 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p0, p0h})
	n5 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p1, p1h})
	n6 := Simd_i8x16_shuffle(n4, n5, [2]uint64{p0, p0h})
	n7 := Simd_v128_load_nc(m, s0, 96)
	return n4[0], n4[1], n5[0], n5[1], n6[0], n6[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx166(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_nc(m, s0, 112)
	n1 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n0, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n0, [2]uint64{p2, p2h})
	n3 := Simd_i8x16_shuffle(n1, n2, [2]uint64{p1, p1h})
	n4 := Simd_i8x16_shuffle(n1, n2, [2]uint64{p2, p2h})
	n5 := Simd_i8x16_shuffle(n3, n4, [2]uint64{p1, p1h})
	n6 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n5, [2]uint64{p4, p4h})
	return n3[0], n3[1], n4[0], n4[1], n5[0], n5[1], n6[0], n6[1]
}

//go:noinline
func Simd_p_fx167(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p4, p4h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p5, p5h})
	n3 := Simd_i16x8_shr_u(n2, 8)
	n4 := Simd_i16x8_shr_u([2]uint64{p6, p6h}, 8)
	n5 := Simd_i8x16_narrow_i16x8_u(n4, n3)
	return n2[0], n2[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx168(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_shr_u([2]uint64{p0, p0h}, 8)
	n1 := Simd_i16x8_shr_u([2]uint64{p1, p1h}, 8)
	n2 := Simd_i8x16_narrow_i16x8_u(n0, n1)
	n3 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n4 := Simd_i16x8_shr_u(n3, 8)
	n5 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, [2]uint64{p6, p6h}, [2]uint64{p4, p4h})
	n6 := Simd_i16x8_shr_u(n5, 8)
	n7 := Simd_i8x16_narrow_i16x8_u(n4, n6)
	return n2[0], n2[1], n3[0], n3[1], n5[0], n5[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx169(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_shr_u([2]uint64{p0, p0h}, 8)
	n1 := Simd_i16x8_shr_u([2]uint64{p1, p1h}, 8)
	n2 := Simd_i8x16_narrow_i16x8_u(n0, n1)
	n3 := Simd_i16x8_shr_u(n2, 8)
	n4 := Simd_i16x8_shr_u([2]uint64{p2, p2h}, 8)
	n5 := Simd_i8x16_narrow_i16x8_u(n4, n3)
	n6 := Simd_v128_and(n5, [2]uint64{p4, p4h})
	n7 := Simd_v128_and([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n8 := Simd_i8x16_narrow_i16x8_u(n7, n6)
	n9 := Simd_v128_and([2]uint64{p5, p5h}, [2]uint64{p4, p4h})
	n10 := Simd_v128_and([2]uint64{p6, p6h}, [2]uint64{p4, p4h})
	n11 := Simd_i8x16_narrow_i16x8_u(n9, n10)
	return n2[0], n2[1], n5[0], n5[1], n8[0], n8[1], n11[0], n11[1]
}

//go:noinline
func Simd_p_fx170(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_and([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_narrow_i16x8_u(n0, n1)
	n3 := Simd_v128_and(n2, [2]uint64{p1, p1h})
	n4 := Simd_v128_and([2]uint64{p3, p3h}, [2]uint64{p1, p1h})
	n5 := Simd_i8x16_narrow_i16x8_u(n4, n3)
	n6 := Simd_i16x8_shr_u(n5, 8)
	n7 := Simd_v128_and([2]uint64{p4, p4h}, [2]uint64{p1, p1h})
	n8 := Simd_v128_and([2]uint64{p5, p5h}, [2]uint64{p1, p1h})
	n9 := Simd_i8x16_narrow_i16x8_u(n7, n8)
	n10 := Simd_i16x8_shr_u(n9, 8)
	n11 := Simd_i8x16_narrow_i16x8_u(n6, n10)
	return n2[0], n2[1], n5[0], n5[1], n9[0], n9[1], n11[0], n11[1]
}

//go:noinline
func Simd_p_fx171(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_shr_u([2]uint64{p0, p0h}, 8)
	n1 := Simd_i16x8_shr_u([2]uint64{p1, p1h}, 8)
	n2 := Simd_i8x16_narrow_i16x8_u(n0, n1)
	n3 := Simd_v128_and([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n4 := Simd_v128_and([2]uint64{p4, p4h}, [2]uint64{p3, p3h})
	n5 := Simd_i8x16_narrow_i16x8_u(n3, n4)
	n6 := Simd_i16x8_shr_u(n5, 8)
	n7 := Simd_i16x8_shr_u([2]uint64{p5, p5h}, 8)
	n8 := Simd_i16x8_shr_u([2]uint64{p6, p6h}, 8)
	n9 := Simd_i8x16_narrow_i16x8_u(n7, n8)
	n10 := Simd_i16x8_shr_u(n9, 8)
	n11 := Simd_i8x16_narrow_i16x8_u(n6, n10)
	return n2[0], n2[1], n5[0], n5[1], n9[0], n9[1], n11[0], n11[1]
}

//go:noinline
func Simd_p_fx172(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_shr_u([2]uint64{p0, p0h}, 8)
	n1 := Simd_i16x8_shr_u([2]uint64{p1, p1h}, 8)
	n2 := Simd_i8x16_narrow_i16x8_u(n0, n1)
	n3 := Simd_i16x8_shr_u(n2, 8)
	n4 := Simd_v128_and(n2, [2]uint64{p4, p4h})
	n5 := Simd_i16x8_shr_u([2]uint64{p2, p2h}, 8)
	n6 := Simd_i8x16_narrow_i16x8_u(n5, n3)
	n7 := Simd_i16x8_shr_u(n6, 8)
	n8 := Simd_i16x8_shr_u([2]uint64{p3, p3h}, 8)
	n9 := Simd_i8x16_narrow_i16x8_u(n8, n7)
	n10 := Simd_v128_and([2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	n11 := Simd_i8x16_narrow_i16x8_u(n10, n4)
	n12 := Simd_v128_and([2]uint64{p5, p5h}, [2]uint64{p4, p4h})
	n13 := Simd_v128_and([2]uint64{p6, p6h}, [2]uint64{p4, p4h})
	n14 := Simd_i8x16_narrow_i16x8_u(n12, n13)
	_ = Simd_v128_store(m, s0, 80, n9)
	return n6[0], n6[1], n11[0], n11[1], n14[0], n14[1]
}

//go:noinline
func Simd_p_fx173(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_and([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_narrow_i16x8_u(n0, n1)
	n3 := Simd_i16x8_shr_u(n2, 8)
	n4 := Simd_v128_and(n2, [2]uint64{p1, p1h})
	n5 := Simd_i16x8_shr_u([2]uint64{p3, p3h}, 8)
	n6 := Simd_i8x16_narrow_i16x8_u(n5, n3)
	n7 := Simd_i16x8_shr_u(n6, 8)
	n8 := Simd_i16x8_shr_u([2]uint64{p4, p4h}, 8)
	n9 := Simd_i8x16_narrow_i16x8_u(n8, n7)
	n10 := Simd_v128_and([2]uint64{p3, p3h}, [2]uint64{p1, p1h})
	n11 := Simd_i8x16_narrow_i16x8_u(n10, n4)
	n12 := Simd_i16x8_shr_u(n11, 8)
	n13 := Simd_v128_and([2]uint64{p5, p5h}, [2]uint64{p1, p1h})
	n14 := Simd_v128_and([2]uint64{p6, p6h}, [2]uint64{p1, p1h})
	n15 := Simd_i8x16_narrow_i16x8_u(n13, n14)
	n16 := Simd_i16x8_shr_u(n15, 8)
	n17 := Simd_i8x16_narrow_i16x8_u(n12, n16)
	_ = Simd_v128_store(m, s0, 64, n9)
	_ = Simd_v128_store(m, s0, 48, n17)
	return n6[0], n6[1], n11[0], n11[1], n15[0], n15[1]
}

//go:noinline
func Simd_p_fx174(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) {
	n0 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_and([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_narrow_i16x8_u(n0, n1)
	n3 := Simd_v128_and([2]uint64{p3, p3h}, [2]uint64{p1, p1h})
	n4 := Simd_v128_and([2]uint64{p4, p4h}, [2]uint64{p1, p1h})
	n5 := Simd_i8x16_narrow_i16x8_u(n3, n4)
	n6 := Simd_v128_and([2]uint64{p5, p5h}, [2]uint64{p1, p1h})
	n7 := Simd_v128_and([2]uint64{p6, p6h}, [2]uint64{p1, p1h})
	n8 := Simd_i8x16_narrow_i16x8_u(n6, n7)
	_ = Simd_v128_store(m, s0, 32, n2)
	_ = Simd_v128_store(m, s0, 16, n5)
	_ = Simd_v128_store(m, s0, 0, n8)
	return
}

//go:noinline
func Simd_p_fx175(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) {
	n0 := Simd_v128_load(m, s0, 16)
	n1 := Simd_v128_and(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_and(n0, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n4 := Simd_i8x16_shuffle(n3, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	n5 := Simd_v128_or(n1, n4)
	_ = Simd_v128_store(m, s1, 16, n5)
	n7 := Simd_v128_load(m, s0, 0)
	n8 := Simd_v128_and(n7, [2]uint64{p0, p0h})
	n9 := Simd_v128_and(n7, [2]uint64{p1, p1h})
	n10 := Simd_i8x16_shuffle(n9, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n11 := Simd_i8x16_shuffle(n10, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	n12 := Simd_v128_or(n8, n11)
	_ = Simd_v128_store(m, s1, 0, n12)
	return
}

//go:noinline
func Simd_p_fx176(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 32)
	n1 := Simd_v128_load_nc(m, s0, 16)
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p0, p0h})
	n3 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p1, p1h})
	n4 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p0, p0h})
	n5 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p1, p1h})
	n6 := Simd_i8x16_shuffle(n4, n5, [2]uint64{p0, p0h})
	n7 := Simd_i8x16_shuffle(n4, n5, [2]uint64{p1, p1h})
	n8 := Simd_i8x16_shuffle(n7, n6, [2]uint64{506097522914230528, 1663540288323457296})
	n9 := Simd_v128_and(n8, [2]uint64{p3, p3h})
	n10 := Simd_i8x16_shuffle(n6, n7, [2]uint64{1084818905618843912, 2242261671028070680})
	n11 := Simd_i16x8_shr_u(n10, 4)
	n12 := Simd_v128_and(n11, [2]uint64{p2, p2h})
	n13 := Simd_v128_or(n12, n9)
	n14 := Simd_i8x16_shuffle(n13, [2]uint64{p4, p4h}, [2]uint64{1084818905618843912, 1663540288323457296})
	n15 := Simd_i8x16_shuffle(n13, n14, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s1, 0, n15)
	return
}

//go:noinline
func Simd_p_fx177(m *Module, s0 int32, s1 int32, p0, p0h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, s1)
	n1 := Simd_v128_and(n0, [2]uint64{p0, p0h})
	n2 := Simd_i32x4_shr_u(n0, 28)
	n3 := Simd_v128_or(n1, n2)
	return n0[0], n0[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx178(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_shr_u([2]uint64{p0, p0h}, s0)
	n1 := Simd_v128_and(n0, [2]uint64{p1, p1h})
	n2 := Simd_i32x4_shr_u([2]uint64{p0, p0h}, 12)
	n3 := Simd_v128_and(n2, [2]uint64{p2, p2h})
	n4 := Simd_v128_or(n1, n3)
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx179(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 32)
	n1 := Simd_v128_load_nc(m, s0, 16)
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p0, p0h})
	n3 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p1, p1h})
	n4 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p0, p0h})
	n5 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p1, p1h})
	n6 := Simd_i8x16_shuffle(n4, n5, [2]uint64{p0, p0h})
	n7 := Simd_i8x16_shuffle(n4, n5, [2]uint64{p1, p1h})
	n8 := Simd_i8x16_shuffle(n7, n6, [2]uint64{506097522914230528, 1663540288323457296})
	n9 := Simd_v128_and(n8, [2]uint64{p2, p2h})
	n10 := Simd_i8x16_shuffle(n9, [2]uint64{p5, p5h}, [2]uint64{1084818905618843912, 1663540288323457296})
	n11 := Simd_i16x8_shr_u(n10, 3)
	n12 := Simd_i8x16_shuffle(n6, n7, [2]uint64{1084818905618843912, 2242261671028070680})
	n13 := Simd_i16x8_shr_u(n12, 5)
	n14 := Simd_v128_and(n13, [2]uint64{p3, p3h})
	n15 := Simd_v128_or(n14, n9)
	n16 := Simd_i16x8_shl(n12, 3)
	n17 := Simd_v128_and(n16, [2]uint64{p4, p4h})
	n18 := Simd_v128_or(n17, n11)
	n19 := Simd_i8x16_shuffle(n15, n18, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s1, 0, n19)
	return
}

//go:noinline
func Simd_p_fx180(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, s1)
	n1 := Simd_i32x4_shr_u(n0, s2)
	n2 := Simd_v128_and(n1, [2]uint64{p0, p0h})
	n3 := Simd_i32x4_shr_u(n0, 13)
	n4 := Simd_v128_and(n3, [2]uint64{p1, p1h})
	n5 := Simd_v128_or(n2, n4)
	return n0[0], n0[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx181(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_shr_u([2]uint64{p0, p0h}, s0)
	n1 := Simd_v128_and(n0, [2]uint64{p1, p1h})
	n2 := Simd_i32x4_shr_u([2]uint64{p0, p0h}, s1)
	n3 := Simd_v128_and(n2, [2]uint64{p2, p2h})
	n4 := Simd_v128_or(n1, n3)
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx182(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i64x2_shr_u(n0, 8)
	n2 := Simd_v128_and(n1, [2]uint64{p0, p0h})
	n3 := Simd_v128_and(n0, [2]uint64{p1, p1h})
	n4 := Simd_v128_or(n2, n3)
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx183(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 16)
	n1 := Simd_i64x2_shr_u(n0, 8)
	n2 := Simd_v128_and(n1, [2]uint64{p0, p0h})
	n3 := Simd_v128_and(n0, [2]uint64{p1, p1h})
	n4 := Simd_v128_or(n2, n3)
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx184(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1084818905618843912, 506097522914230528})
	n1 := Simd_v128_and(n0, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{216736831831411980, 795458214266537220})
	n3 := Simd_v128_and(n2, [2]uint64{p2, p2h})
	n4 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{795458214266537220, 216736831831411980})
	n5 := Simd_v128_and(n4, [2]uint64{p2, p2h})
	n6 := Simd_v128_load_rng(m, s0, 32, 32, 32)
	n7 := Simd_i8x16_swizzle(n6, n1)
	n8 := Simd_v128_load_nc(m, s0, 48)
	n9 := Simd_i8x16_swizzle(n8, n3)
	n10 := Simd_i8x16_shuffle(n7, n9, [2]uint64{1663540288053969152, 2242261671028070680})
	_ = Simd_v128_store(m, s1, 32, n10)
	n12 := Simd_v128_load(m, s0, 16)
	n13 := Simd_i8x16_swizzle(n12, n5)
	n14 := Simd_i8x16_shuffle(n13, n7, [2]uint64{506097522914230528, 2242261671028070680})
	_ = Simd_v128_store(m, s1, 16, n14)
	n16 := Simd_v128_load(m, s0, 0)
	n17 := Simd_i8x16_swizzle(n16, [2]uint64{p3, p3h})
	n18 := Simd_i8x16_shuffle(n17, n13, [2]uint64{506097522914230528, 2242261670758582536})
	_ = Simd_v128_store(m, s1, 0, n18)
	return
}

//go:noinline
func Simd_p_fx185(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_i32x4_extend_low_i16x8_s([2]uint64{p1, p1h})
	n1 := Simd_i32x4_extend_high_i16x8_s([2]uint64{p1, p1h})
	n2 := Simd_i32x4_extend_low_i16x8_s([2]uint64{p4, p4h})
	n3 := Simd_i32x4_extend_high_i16x8_s([2]uint64{p4, p4h})
	n4 := Simd_v128_load(m, s0, 0)
	n5 := Simd_i8x16_swizzle(n4, [2]uint64{p0, p0h})
	n6 := Simd_i32x4_extend_low_i16x8_s(n5)
	n7 := Simd_i32x4_mul(n6, n0)
	n8 := Simd_i32x4_extend_high_i16x8_s(n5)
	n9 := Simd_i32x4_mul(n8, n1)
	n10 := Simd_i8x16_shuffle(n7, n9, [2]uint64{p2, p2h})
	n11 := Simd_i8x16_add(n4, n10)
	n12 := Simd_i8x16_swizzle(n11, [2]uint64{p3, p3h})
	n13 := Simd_i32x4_extend_low_i16x8_s(n12)
	n14 := Simd_i32x4_mul(n13, n2)
	n15 := Simd_i32x4_extend_high_i16x8_s(n12)
	n16 := Simd_i32x4_mul(n15, n3)
	n17 := Simd_i8x16_shuffle(n14, n16, [2]uint64{p2, p2h})
	n18 := Simd_i8x16_add(n17, n11)
	n19 := Simd_v128_bitselect(n4, n18, [2]uint64{p5, p5h})
	_ = Simd_v128_store(m, s1, 0, n19)
	return
}

//go:noinline
func Simd_p_fx186(m *Module, s0 int32, s1 int32, p0, p0h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i32x4_add(n0, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s1, 0, n1)
	return
}

//go:noinline
func Simd_p_fx187(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_v128_load_rng(m, s0+4, 0, -4, 20)
	n1 := Simd_v128_load_nc(m, s0, 0)
	n2 := Simd_v128_xor(n0, n1)
	n3 := Simd_i32x4_shr_u(n2, 1)
	n4 := Simd_v128_and(n3, [2]uint64{p0, p0h})
	n5 := Simd_v128_and(n0, n1)
	n6 := Simd_i32x4_add(n4, n5)
	n7 := Simd_v128_and(n6, [2]uint64{p1, p1h})
	n8 := Simd_v128_and(n6, [2]uint64{p2, p2h})
	n9 := Simd_v128_load(m, s1, 0)
	n10 := Simd_v128_and(n9, [2]uint64{p1, p1h})
	n11 := Simd_i32x4_add(n7, n10)
	n12 := Simd_v128_and(n11, [2]uint64{p1, p1h})
	n13 := Simd_v128_and(n9, [2]uint64{p2, p2h})
	n14 := Simd_i32x4_add(n8, n13)
	n15 := Simd_v128_and(n14, [2]uint64{p2, p2h})
	n16 := Simd_v128_or(n12, n15)
	_ = Simd_v128_store(m, s2, 0, n16)
	return
}

//go:noinline
func Simd_p_fx188(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_and(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_and(n0, [2]uint64{p1, p1h})
	n3 := Simd_v128_load(m, s1, 0)
	n4 := Simd_v128_and(n3, [2]uint64{p0, p0h})
	n5 := Simd_i32x4_add(n1, n4)
	n6 := Simd_v128_and(n5, [2]uint64{p0, p0h})
	n7 := Simd_v128_and(n3, [2]uint64{p1, p1h})
	n8 := Simd_i32x4_add(n2, n7)
	n9 := Simd_v128_and(n8, [2]uint64{p1, p1h})
	n10 := Simd_v128_or(n6, n9)
	_ = Simd_v128_store(m, s2, 0, n10)
	return
}

//go:noinline
func Simd_p_fx189(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_f64x2_convert_low_i32x4_u([2]uint64{p0, p0h})
	n1 := Simd_f64x2_mul(n0, [2]uint64{p1, p1h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx190(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_f64x2_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_f64x2_add(n0, [2]uint64{p2, p2h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx191(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1084818905618843912, 216736831629295872})
	n1 := Simd_f64x2_convert_low_i32x4_u(n0)
	n2 := Simd_f64x2_mul(n1, [2]uint64{p2, p2h})
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx192(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{1084818905618843912, 216736831629295872})
	n1 := Simd_f64x2_convert_low_i32x4_u(n0)
	n2 := Simd_f64x2_mul(n1, [2]uint64{p1, p1h})
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx193(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_extend_low_i8x16_u(n0)
	n2 := Simd_i32x4_extend_low_i16x8_u(n1)
	n3 := Simd_i32x4_mul(n2, [2]uint64{p2, p2h})
	n4 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n5 := Simd_i16x8_extend_low_i8x16_u(n4)
	n6 := Simd_i32x4_extend_low_i16x8_u(n5)
	n7 := Simd_i32x4_mul(n6, [2]uint64{p4, p4h})
	n8 := Simd_i32x4_add(n3, n7)
	n9 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n10 := Simd_i16x8_extend_low_i8x16_u(n9)
	n11 := Simd_i32x4_extend_low_i16x8_u(n10)
	n12 := Simd_i32x4_mul(n11, [2]uint64{p6, p6h})
	n13 := Simd_i32x4_add(n8, n12)
	return n13[0], n13[1]
}

//go:noinline
func Simd_p_fx194(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load32_zero(m, s0, s1)
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p0, p0h}, [2]uint64{4294967296, 12884901890})
	n2 := Simd_i32x4_shl(n1, 24)
	n3 := Simd_v128_load(m, s2, s1)
	n4 := Simd_v128_and(n3, [2]uint64{p1, p1h})
	n5 := Simd_v128_or(n2, n4)
	_ = Simd_v128_store(m, s2, s1, n5)
	return n0[0], n0[1], n3[0], n3[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx195(m *Module, s0 int32, s1 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_i8x16_eq(n1, [2]uint64{p0, p0h})
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx196(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i16x8_mul(n0, [2]uint64{p3, p3h})
	return n0[0], n0[1], n1[0], n1[1]
}

//go:noinline
func Simd_p_fx197(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_dot_i16x8_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_dot_i16x8_s([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i32x4_add(n0, n1)
	n3 := Simd_i32x4_dot_i16x8_s([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n4 := Simd_i32x4_add(n2, n3)
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx198(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i32x4_add([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx199(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_add(n0, [2]uint64{p2, p2h})
	n2 := Simd_i16x8_add(n1, [2]uint64{p3, p3h})
	n3 := Simd_i16x8_add(n2, [2]uint64{p4, p4h})
	n4 := Simd_i16x8_add(n3, [2]uint64{p5, p5h})
	n5 := Simd_i16x8_add(n4, [2]uint64{p6, p6h})
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx200(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i16x8_add([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx201(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_sub_sat_u([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_sub_sat_u([2]uint64{p1, p1h}, [2]uint64{p0, p0h})
	n2 := Simd_v128_or(n0, n1)
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n4 := Simd_i32x4_dot_i16x8_s(n3, n3)
	n5 := Simd_i32x4_add(n4, [2]uint64{p5, p5h})
	n6 := Simd_i8x16_shuffle(n2, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	n7 := Simd_i32x4_dot_i16x8_s(n6, n6)
	n8 := Simd_i32x4_add(n5, n7)
	n9 := Simd_v128_load_rng(m, s0+16, 0, 0, 32)
	n10 := Simd_v128_load_rng(m, s1+16, 0, 0, 32)
	n11 := Simd_i8x16_sub_sat_u(n10, n9)
	n12 := Simd_i8x16_sub_sat_u(n9, n10)
	n13 := Simd_v128_or(n12, n11)
	n14 := Simd_i8x16_shuffle(n13, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n15 := Simd_i32x4_dot_i16x8_s(n14, n14)
	n16 := Simd_i8x16_shuffle(n13, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	n17 := Simd_i32x4_dot_i16x8_s(n16, n16)
	n18 := Simd_i32x4_add(n8, n15)
	n19 := Simd_i32x4_add(n18, n17)
	n20 := Simd_v128_load_nc(m, s1+32, 0)
	n21 := Simd_v128_load_nc(m, s0+32, 0)
	return n19[0], n19[1], n20[0], n20[1], n21[0], n21[1]
}

//go:noinline
func Simd_p_fx202(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_sub_sat_u([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_sub_sat_u([2]uint64{p1, p1h}, [2]uint64{p0, p0h})
	n2 := Simd_v128_or(n0, n1)
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p2, p2h}, [2]uint64{1948679894439893000, 2238040585792199692})
	n4 := Simd_i32x4_dot_i16x8_s(n3, n3)
	n5 := Simd_i32x4_add(n4, [2]uint64{p3, p3h})
	n6 := Simd_i8x16_shuffle(n2, [2]uint64{p2, p2h}, [2]uint64{1369958511735279616, 1659319203087586308})
	n7 := Simd_i32x4_dot_i16x8_s(n6, n6)
	n8 := Simd_i32x4_add(n5, n7)
	return n8[0], n8[1]
}

//go:noinline
func Simd_p_fx203(m *Module, s0 int32, s1 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_v128_load32_zero(m, s0, 0)
	n1 := Simd_i16x8_extend_low_i8x16_u(n0)
	n2 := Simd_i32x4_extend_low_i16x8_u(n1)
	n3 := Simd_v128_load32_zero(m, s1, 0)
	n4 := Simd_i16x8_extend_low_i8x16_u(n3)
	n5 := Simd_i32x4_extend_low_i16x8_u(n4)
	n6 := Simd_i32x4_sub(n2, n5)
	n7 := Simd_i32x4_mul(n6, n6)
	n8 := Simd_i32x4_add(n7, [2]uint64{p0, p0h})
	return n8[0], n8[1]
}

//go:noinline
func Simd_p_fx204(m *Module, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{1084818905618843912, 216736831629295872})
	n1 := Simd_i32x4_add([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx205(m *Module, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{216736831696667908, 216736831629295872})
	n1 := Simd_i32x4_add([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx206(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, p0, p0h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_splat(s1)
	n1 := Simd_v128_load(m, s0, 0)
	n2 := Simd_i32x4_mul(n0, n1)
	n3 := Simd_scalar_i32_add(s2, s3)
	n4 := Simd_v128_load32_zero(m, n3, 0)
	n5 := Simd_i16x8_extend_low_i8x16_u(n4)
	n6 := Simd_i32x4_extend_low_i16x8_u(n5)
	n7 := Simd_i32x4_mul(n2, n6)
	n8 := Simd_i32x4_add(n7, [2]uint64{p0, p0h})
	return n2[0], n2[1], n6[0], n6[1], n7[0], n7[1], n8[0], n8[1]
}

//go:noinline
func Simd_p_fx207(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n1 := Simd_i32x4_add(n0, [2]uint64{p4, p4h})
	n2 := Simd_scalar_i32_add(s0, s1)
	n3 := Simd_v128_load32_zero(m, n2, 0)
	n4 := Simd_i16x8_extend_low_i8x16_u(n3)
	n5 := Simd_i32x4_extend_low_i16x8_u(n4)
	n6 := Simd_i32x4_mul([2]uint64{p0, p0h}, n5)
	n7 := Simd_i32x4_add(n6, [2]uint64{p1, p1h})
	return n5[0], n5[1], n6[0], n6[1], n7[0], n7[1], n1[0], n1[1]
}

//go:noinline
func Simd_p_fx208(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_add(n0, [2]uint64{p2, p2h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx209(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_add([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx210(m *Module, s0 int32, s1 int32, s2 int32, s3 int32) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_v128_load(m, s2, 0)
	n3 := Simd_v128_load(m, s3, 0)
	n4 := Simd_i8x16_sub_sat_u(n3, n2)
	n5 := Simd_i8x16_sub_sat_u(n2, n3)
	n6 := Simd_v128_or(n5, n4)
	return n1[0], n1[1], n2[0], n2[1], n3[0], n3[1], n6[0], n6[1]
}

//go:noinline
func Simd_p_fx211(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, p0, p0h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_i8x16_sub_sat_u(n1, [2]uint64{p0, p0h})
	n3 := Simd_i8x16_sub_sat_u([2]uint64{p0, p0h}, n1)
	n4 := Simd_v128_or(n2, n3)
	n5 := Simd_scalar_i32_add(s2, s3)
	n6 := Simd_v128_load(m, n5, 0)
	n7 := Simd_v128_load(m, s0, 0)
	return n1[0], n1[1], n6[0], n6[1], n4[0], n4[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx212(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_sub_sat_u([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_sub_sat_u([2]uint64{p1, p1h}, [2]uint64{p0, p0h})
	n2 := Simd_v128_or(n0, n1)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx213(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_add_sat_u([2]uint64{p0, p0h}, [2]uint64{p0, p0h})
	n1 := Simd_i8x16_sub_sat_u([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_sub_sat_u([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n3 := Simd_v128_or(n1, n2)
	n4 := Simd_i16x8_shr_u(n3, 1)
	n5 := Simd_v128_and(n4, [2]uint64{p3, p3h})
	n6 := Simd_i8x16_add_sat_u(n0, n5)
	n7 := Simd_i8x16_splat(s0)
	n8 := Simd_i8x16_sub_sat_u(n6, n7)
	return n8[0], n8[1]
}

//go:noinline
func Simd_p_fx214(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_sub_sat_u([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_sub_sat_u([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_v128_or(n0, n1)
	n3 := Simd_i8x16_max_u([2]uint64{p0, p0h}, n2)
	n4 := Simd_i8x16_sub_sat_u([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n5 := Simd_i8x16_sub_sat_u([2]uint64{p3, p3h}, [2]uint64{p2, p2h})
	n6 := Simd_v128_or(n4, n5)
	n7 := Simd_i8x16_max_u(n3, n6)
	n8 := Simd_i8x16_max_u(n7, [2]uint64{p4, p4h})
	n9 := Simd_i8x16_sub_sat_u([2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n10 := Simd_i8x16_sub_sat_u([2]uint64{p6, p6h}, [2]uint64{p5, p5h})
	n11 := Simd_v128_or(n9, n10)
	n12 := Simd_i8x16_max_u(n8, n11)
	return n12[0], n12[1]
}

//go:noinline
func Simd_p_fx215(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_xor([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_sub_sat_s(n0, n1)
	n3 := Simd_v128_xor([2]uint64{p3, p3h}, [2]uint64{p1, p1h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx216(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_sub_sat_s([2]uint64{p3, p3h}, n0)
	n2 := Simd_i8x16_add_sat_s([2]uint64{p2, p2h}, n1)
	n3 := Simd_i8x16_add_sat_s([2]uint64{p2, p2h}, n2)
	n4 := Simd_i8x16_add_sat_s([2]uint64{p2, p2h}, n3)
	n5 := Simd_i8x16_max_u([2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n6 := Simd_i8x16_splat(s0)
	n7 := Simd_i8x16_sub_sat_u(n5, n6)
	n8 := Simd_i8x16_eq([2]uint64{p4, p4h}, n7)
	return n0[0], n0[1], n4[0], n4[1], n8[0], n8[1]
}

//go:noinline
func Simd_p_fx217(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_and(n0, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n1, [2]uint64{p4, p4h})
	n3 := Simd_i32x4_extend_low_i16x8_s(n2)
	n4 := Simd_i32x4_mul(n3, [2]uint64{p5, p5h})
	n5 := Simd_i32x4_extend_high_i16x8_s(n2)
	n6 := Simd_i32x4_mul(n5, [2]uint64{p5, p5h})
	n7 := Simd_i8x16_shuffle(n4, n6, [2]uint64{p6, p6h})
	return n1[0], n1[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx218(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_shr_s(n0, 7)
	n2 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n3 := Simd_i32x4_extend_low_i16x8_s(n2)
	n4 := Simd_i32x4_mul(n3, [2]uint64{p5, p5h})
	n5 := Simd_i32x4_extend_high_i16x8_s(n2)
	n6 := Simd_i32x4_mul(n5, [2]uint64{p5, p5h})
	n7 := Simd_i8x16_shuffle(n4, n6, [2]uint64{p6, p6h})
	n8 := Simd_i16x8_add(n7, [2]uint64{p1, p1h})
	n9 := Simd_i16x8_shr_s(n8, 7)
	n10 := Simd_i8x16_narrow_i16x8_s(n1, n9)
	return n0[0], n0[1], n7[0], n7[1], n8[0], n8[1], n10[0], n10[1]
}

//go:noinline
func Simd_p_fx219(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_add_sat_s(n0, [2]uint64{p2, p2h})
	n2 := Simd_v128_xor(n1, [2]uint64{p1, p1h})
	n3 := Simd_i16x8_add([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n4 := Simd_i16x8_shr_s(n3, 7)
	n5 := Simd_i16x8_add([2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n6 := Simd_i16x8_shr_s(n5, 7)
	n7 := Simd_i8x16_narrow_i16x8_s(n4, n6)
	_ = Simd_v128_store(m, s0, 0, n2)
	return n3[0], n3[1], n5[0], n5[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx220(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_add_sat_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_xor(n0, [2]uint64{p2, p2h})
	n2 := Simd_v128_andnot([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n3 := Simd_v128_and(n2, [2]uint64{p5, p5h})
	n4 := Simd_i8x16_add_sat_s(n3, [2]uint64{p6, p6h})
	_ = Simd_v128_store(m, s0, 0, n1)
	return n3[0], n3[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx221(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_shr_s(n0, 7)
	n2 := Simd_i16x8_add([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n3 := Simd_i16x8_shr_s(n2, 7)
	n4 := Simd_i8x16_narrow_i16x8_s(n1, n3)
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx222(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n1 := Simd_i16x8_shr_s(n0, 11)
	n2 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	n3 := Simd_i16x8_shr_s(n2, 11)
	n4 := Simd_i8x16_narrow_i16x8_s(n1, n3)
	n5 := Simd_i8x16_add_sat_s([2]uint64{p0, p0h}, n4)
	n6 := Simd_i8x16_add_sat_s(n5, [2]uint64{p5, p5h})
	n7 := Simd_v128_xor(n6, [2]uint64{p6, p6h})
	_ = Simd_v128_store(m, s0, 0, n7)
	return
}

//go:noinline
func Simd_p_fx223(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n1 := Simd_i16x8_shr_s(n0, 11)
	n2 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	n3 := Simd_i16x8_shr_s(n2, 11)
	n4 := Simd_i8x16_narrow_i16x8_s(n1, n3)
	n5 := Simd_i8x16_sub_sat_s([2]uint64{p0, p0h}, n4)
	n6 := Simd_i8x16_sub_sat_s(n5, [2]uint64{p5, p5h})
	n7 := Simd_v128_xor(n6, [2]uint64{p6, p6h})
	_ = Simd_v128_store(m, s0, 0, n7)
	return
}

//go:noinline
func Simd_p_fx224(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) {
	n0 := Simd_i8x16_sub_sat_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_xor(n0, [2]uint64{p2, p2h})
	n2 := Simd_v128_xor([2]uint64{p3, p3h}, [2]uint64{p2, p2h})
	n3 := Simd_i8x16_sub_sat_s(n2, [2]uint64{p4, p4h})
	n4 := Simd_v128_xor(n3, [2]uint64{p2, p2h})
	_ = Simd_v128_store(m, s0, 0, n1)
	_ = Simd_v128_store(m, s1, 0, n4)
	return
}

//go:noinline
func Simd_p_fx225(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_extend_low_i16x8_s([2]uint64{p0, p0h})
	n1 := Simd_i32x4_mul(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_mul(n0, [2]uint64{p4, p4h})
	n3 := Simd_i32x4_extend_high_i16x8_s([2]uint64{p0, p0h})
	n4 := Simd_i32x4_mul(n3, [2]uint64{p2, p2h})
	n5 := Simd_i8x16_shuffle(n1, n4, [2]uint64{p3, p3h})
	n6 := Simd_i32x4_mul(n3, [2]uint64{p4, p4h})
	n7 := Simd_i8x16_shuffle(n2, n6, [2]uint64{p3, p3h})
	n8 := Simd_i32x4_extend_low_i16x8_s([2]uint64{p1, p1h})
	n9 := Simd_i32x4_mul(n8, [2]uint64{p4, p4h})
	n10 := Simd_i32x4_mul(n8, [2]uint64{p2, p2h})
	n11 := Simd_i32x4_extend_high_i16x8_s([2]uint64{p1, p1h})
	n12 := Simd_i32x4_mul(n11, [2]uint64{p4, p4h})
	n13 := Simd_i8x16_shuffle(n9, n12, [2]uint64{p3, p3h})
	n14 := Simd_i32x4_mul(n11, [2]uint64{p2, p2h})
	n15 := Simd_i8x16_shuffle(n10, n14, [2]uint64{p3, p3h})
	n16 := Simd_i16x8_sub(n15, n7)
	n17 := Simd_i16x8_add([2]uint64{p1, p1h}, [2]uint64{p0, p0h})
	n18 := Simd_i16x8_add(n17, n5)
	n19 := Simd_i16x8_add(n18, n13)
	n20 := Simd_i16x8_add([2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n21 := Simd_i16x8_add(n19, n20)
	n22 := Simd_i16x8_sub([2]uint64{p1, p1h}, [2]uint64{p0, p0h})
	n23 := Simd_i16x8_add(n16, n22)
	return n19[0], n19[1], n20[0], n20[1], n21[0], n21[1], n23[0], n23[1]
}

//go:noinline
func Simd_p_fx226(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_sub(n0, [2]uint64{p2, p2h})
	n2 := Simd_i16x8_add([2]uint64{p2, p2h}, n0)
	n3 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n2, [2]uint64{p4, p4h})
	n4 := Simd_i16x8_sub([2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	return n2[0], n2[1], n3[0], n3[1], n1[0], n1[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx227(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p4, p4h})
	return n0[0], n0[1], n1[0], n1[1]
}

//go:noinline
func Simd_p_fx228(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p4, p4h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p5, p5h})
	n3 := Simd_i32x4_extend_low_i16x8_s(n2)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx229(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_extend_high_i16x8_s([2]uint64{p0, p0h})
	n1 := Simd_i32x4_mul(n0, [2]uint64{p3, p3h})
	n2 := Simd_i32x4_extend_low_i16x8_s([2]uint64{p1, p1h})
	n3 := Simd_i32x4_mul(n2, [2]uint64{p5, p5h})
	n4 := Simd_i32x4_extend_high_i16x8_s([2]uint64{p1, p1h})
	n5 := Simd_i32x4_mul(n4, [2]uint64{p5, p5h})
	n6 := Simd_i8x16_shuffle(n3, n5, [2]uint64{p4, p4h})
	n7 := Simd_i16x8_add([2]uint64{p1, p1h}, [2]uint64{p0, p0h})
	n8 := Simd_i32x4_mul([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n9 := Simd_i8x16_shuffle(n8, n1, [2]uint64{p4, p4h})
	n10 := Simd_i16x8_add(n7, n9)
	n11 := Simd_i16x8_add(n10, n6)
	return n0[0], n0[1], n2[0], n2[1], n4[0], n4[1], n11[0], n11[1]
}

//go:noinline
func Simd_p_fx230(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i16x8_add(n0, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle([2]uint64{p4, p4h}, [2]uint64{p5, p5h}, [2]uint64{p2, p2h})
	n3 := Simd_i16x8_add(n1, n2)
	n4 := Simd_i16x8_add([2]uint64{p6, p6h}, n3)
	n5 := Simd_i16x8_shr_s(n4, 3)
	return n1[0], n1[1], n2[0], n2[1], n3[0], n3[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx231(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_mul([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p3, p3h})
	n3 := Simd_i32x4_mul([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n4 := Simd_i32x4_mul([2]uint64{p6, p6h}, [2]uint64{p5, p5h})
	n5 := Simd_i8x16_shuffle(n3, n4, [2]uint64{p3, p3h})
	n6 := Simd_i16x8_sub(n2, n5)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx232(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_sub(n0, [2]uint64{p2, p2h})
	n2 := Simd_i16x8_shr_s(n1, 3)
	n3 := Simd_i16x8_add([2]uint64{p2, p2h}, n0)
	n4 := Simd_i16x8_shr_s(n3, 3)
	n5 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n4, [2]uint64{p4, p4h})
	n6 := Simd_i16x8_sub([2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n7 := Simd_i16x8_shr_s(n6, 3)
	return n4[0], n4[1], n5[0], n5[1], n2[0], n2[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx233(m *Module, s0 int32) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load32_zero(m, s0, 96)
	n1 := Simd_v128_load32_zero(m, s0, 64)
	n2 := Simd_v128_load32_zero(m, s0, 32)
	n3 := Simd_v128_load32_zero(m, s0, 0)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx234(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n2 := Simd_i16x8_add(n0, n1)
	n3 := Simd_i8x16_narrow_i16x8_u(n2, n2)
	return n3[0], n3[1]
}

//go:noinline
func Simd_p_fx235(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, s5 int32, s6 int32, s7 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load32_splat(m, s0, 0)
	n1 := Simd_v128_load32_lane(m, s1, 0, 1, n0)
	n2 := Simd_scalar_i32_add(s0, s2)
	n3 := Simd_v128_load32_lane(m, n2, 0, 2, n1)
	n4 := Simd_scalar_i32_add(s0, s3)
	n5 := Simd_v128_load32_lane(m, n4, 0, 3, n3)
	n6 := Simd_scalar_i32_add(s0, s4)
	n7 := Simd_v128_load32_splat(m, n6, 0)
	n8 := Simd_v128_load32_lane(m, s5, 0, 1, n7)
	n9 := Simd_v128_load32_lane(m, s6, 0, 2, n8)
	n10 := Simd_v128_load32_lane(m, s7, 0, 3, n9)
	n11 := Simd_i8x16_shuffle(n5, n10, [2]uint64{p0, p0h})
	n12 := Simd_i8x16_shuffle(n5, n10, [2]uint64{p1, p1h})
	n13 := Simd_i8x16_shuffle(n11, n12, [2]uint64{p2, p2h})
	return n11[0], n11[1], n12[0], n12[1], n13[0], n13[1]
}

//go:noinline
func Simd_p_fx236(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p4, p4h})
	n2 := Simd_scalar_i32_shl(s1, 3)
	n3 := Simd_scalar_i32_add(s0, n2)
	n4 := Simd_v128_load32_splat(m, n3, 0)
	n5 := Simd_scalar_i32_shl(s1, 3)
	n6 := Simd_scalar_i32_add(s0, n5)
	n7 := Simd_scalar_i32_add(n6, s2)
	n8 := Simd_v128_load32_lane(m, n7, 0, 1, n4)
	n9 := Simd_scalar_i32_shl(s1, 3)
	n10 := Simd_scalar_i32_add(s0, n9)
	n11 := Simd_scalar_i32_add(n10, s3)
	n12 := Simd_v128_load32_lane(m, n11, 0, 2, n8)
	n13 := Simd_scalar_i32_shl(s1, 3)
	n14 := Simd_scalar_i32_add(s0, n13)
	n15 := Simd_scalar_i32_add(n14, s4)
	n16 := Simd_v128_load32_lane(m, n15, 0, 3, n12)
	n17 := Simd_scalar_i32_shl(s1, 3)
	n18 := Simd_scalar_i32_add(s0, n17)
	n19 := Simd_scalar_i32_add(n18, s1)
	n20 := Simd_v128_load32_splat(m, n19, 0)
	return n0[0], n0[1], n1[0], n1[1], n16[0], n16[1], n20[0], n20[1]
}

//go:noinline
func Simd_p_fx237(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load32_lane(m, n0, 0, 1, [2]uint64{p0, p0h})
	n2 := Simd_scalar_i32_add(s0, s2)
	n3 := Simd_v128_load32_lane(m, n2, 0, 2, n1)
	n4 := Simd_scalar_i32_add(s0, s3)
	n5 := Simd_v128_load32_lane(m, n4, 0, 3, n3)
	n6 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, n5, [2]uint64{p2, p2h})
	n7 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, n5, [2]uint64{p3, p3h})
	n8 := Simd_i8x16_shuffle(n6, n7, [2]uint64{p4, p4h})
	return n6[0], n6[1], n7[0], n7[1], n8[0], n8[1]
}

//go:noinline
func Simd_p_fx238(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p4, p4h})
	n2 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, n1, [2]uint64{p6, p6h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx239(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p4, p4h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p5, p5h})
	n3 := Simd_v128_load32_splat(m, s0, 0)
	n4 := Simd_scalar_i32_add(s0, s1)
	n5 := Simd_v128_load32_lane(m, n4, 0, 1, n3)
	n6 := Simd_scalar_i32_add(s0, s2)
	n7 := Simd_v128_load32_lane(m, n6, 0, 2, n5)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx240(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, s5 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load32_lane(m, n0, 0, 3, [2]uint64{p0, p0h})
	n2 := Simd_scalar_i32_add(s0, s2)
	n3 := Simd_v128_load32_splat(m, n2, 0)
	n4 := Simd_scalar_i32_add(s0, s3)
	n5 := Simd_v128_load32_lane(m, n4, 0, 1, n3)
	n6 := Simd_scalar_i32_add(s0, s4)
	n7 := Simd_v128_load32_lane(m, n6, 0, 2, n5)
	n8 := Simd_scalar_i32_add(s0, s5)
	n9 := Simd_v128_load32_lane(m, n8, 0, 3, n7)
	n10 := Simd_i8x16_shuffle(n1, n9, [2]uint64{p1, p1h})
	n11 := Simd_i8x16_shuffle(n1, n9, [2]uint64{p2, p2h})
	n12 := Simd_i8x16_shuffle(n10, n11, [2]uint64{p3, p3h})
	return n10[0], n10[1], n11[0], n11[1], n12[0], n12[1]
}

//go:noinline
func Simd_p_fx241(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p4, p4h})
	n2 := Simd_scalar_i32_add(s0, s1)
	n3 := Simd_v128_load32_splat(m, n2, 0)
	n4 := Simd_scalar_i32_add(s0, s1)
	n5 := Simd_scalar_i32_add(n4, s2)
	n6 := Simd_v128_load32_lane(m, n5, 0, 1, n3)
	n7 := Simd_scalar_i32_add(s0, s1)
	n8 := Simd_scalar_i32_add(n7, s3)
	n9 := Simd_v128_load32_lane(m, n8, 0, 2, n6)
	n10 := Simd_scalar_i32_add(s0, s1)
	n11 := Simd_scalar_i32_add(n10, s4)
	n12 := Simd_v128_load32_lane(m, n11, 0, 3, n9)
	return n0[0], n0[1], n1[0], n1[1], n12[0], n12[1]
}

//go:noinline
func Simd_p_fx242(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load32_splat(m, s0, 0)
	n1 := Simd_scalar_i32_add(s1, s2)
	n2 := Simd_v128_load32_lane(m, n1, 0, 1, n0)
	n3 := Simd_scalar_i32_add(s1, s3)
	n4 := Simd_v128_load32_lane(m, n3, 0, 2, n2)
	n5 := Simd_scalar_i32_add(s1, s4)
	n6 := Simd_v128_load32_lane(m, n5, 0, 3, n4)
	n7 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n6, [2]uint64{p1, p1h})
	n8 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n6, [2]uint64{p2, p2h})
	n9 := Simd_i8x16_shuffle(n7, n8, [2]uint64{p3, p3h})
	n10 := Simd_i8x16_shuffle(n7, n8, [2]uint64{p4, p4h})
	return n9[0], n9[1], n10[0], n10[1]
}

//go:noinline
func Simd_p_fx243(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p4, p4h})
	n2 := Simd_i8x16_sub_sat_u(n1, [2]uint64{p5, p5h})
	n3 := Simd_i8x16_sub_sat_u([2]uint64{p5, p5h}, n1)
	n4 := Simd_v128_or(n3, n2)
	return n0[0], n0[1], n1[0], n1[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx244(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_sub_sat_u(n0, [2]uint64{p6, p6h})
	n2 := Simd_i8x16_sub_sat_u([2]uint64{p6, p6h}, n0)
	n3 := Simd_v128_or(n1, n2)
	n4 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	return n0[0], n0[1], n4[0], n4[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx245(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p4, p4h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p5, p5h})
	n3 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p6, p6h})
	return n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx246(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_add_sat_s(n0, [2]uint64{p2, p2h})
	n2 := Simd_v128_xor(n1, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n2, [2]uint64{p4, p4h})
	n4 := Simd_i16x8_add([2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	return n2[0], n2[1], n3[0], n3[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx247(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_shr_s(n0, 7)
	n2 := Simd_i16x8_shr_s([2]uint64{p2, p2h}, 7)
	n3 := Simd_i8x16_narrow_i16x8_s(n2, n1)
	n4 := Simd_i8x16_add_sat_s([2]uint64{p3, p3h}, n3)
	n5 := Simd_v128_xor(n4, [2]uint64{p4, p4h})
	return n0[0], n0[1], n3[0], n3[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx248(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_andnot([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_and(n0, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_add_sat_s(n1, [2]uint64{p3, p3h})
	return n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx249(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n1 := Simd_i16x8_shr_s(n0, 11)
	n2 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	n3 := Simd_i16x8_shr_s(n2, 11)
	n4 := Simd_i8x16_narrow_i16x8_s(n1, n3)
	n5 := Simd_i8x16_add_sat_s([2]uint64{p0, p0h}, n4)
	n6 := Simd_i8x16_add_sat_s(n5, [2]uint64{p5, p5h})
	n7 := Simd_v128_xor(n6, [2]uint64{p6, p6h})
	return n7[0], n7[1]
}

//go:noinline
func Simd_p_fx250(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n1 := Simd_i16x8_shr_s(n0, 11)
	n2 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	n3 := Simd_i16x8_shr_s(n2, 11)
	n4 := Simd_i8x16_narrow_i16x8_s(n1, n3)
	n5 := Simd_i8x16_sub_sat_s([2]uint64{p0, p0h}, n4)
	n6 := Simd_i8x16_sub_sat_s(n5, [2]uint64{p5, p5h})
	n7 := Simd_v128_xor(n6, [2]uint64{p6, p6h})
	return n7[0], n7[1]
}

//go:noinline
func Simd_p_fx251(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_sub_sat_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_xor(n0, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n1, [2]uint64{p4, p4h})
	n3 := Simd_v128_xor([2]uint64{p5, p5h}, [2]uint64{p2, p2h})
	n4 := Simd_i8x16_sub_sat_s(n3, [2]uint64{p6, p6h})
	n5 := Simd_v128_xor(n4, [2]uint64{p2, p2h})
	return n1[0], n1[1], n2[0], n2[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx252(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_sub_sat_u(n0, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_sub_sat_u([2]uint64{p3, p3h}, n0)
	n3 := Simd_v128_or(n2, n1)
	return n0[0], n0[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx253(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_add_sat_s(n0, [2]uint64{p2, p2h})
	n2 := Simd_v128_xor(n1, [2]uint64{p1, p1h})
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx254(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_shr_s(n0, 7)
	n2 := Simd_i16x8_add([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n3 := Simd_i16x8_shr_s(n2, 7)
	n4 := Simd_i8x16_narrow_i16x8_s(n1, n3)
	n5 := Simd_i8x16_add_sat_s([2]uint64{p4, p4h}, n4)
	n6 := Simd_v128_xor(n5, [2]uint64{p5, p5h})
	return n0[0], n0[1], n2[0], n2[1], n4[0], n4[1], n6[0], n6[1]
}

//go:noinline
func Simd_p_fx255(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_sub_sat_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_xor(n0, [2]uint64{p2, p2h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx256(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_sub_sat_s(n0, [2]uint64{p2, p2h})
	n2 := Simd_v128_xor(n1, [2]uint64{p1, p1h})
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx257(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p4, p4h})
	n2 := Simd_v128_load32_splat(m, s0+-4, 0)
	n3 := Simd_scalar_i32_add(s0, -4)
	n4 := Simd_scalar_i32_add(n3, s1)
	n5 := Simd_v128_load32_lane(m, n4, 0, 1, n2)
	n6 := Simd_scalar_i32_add(s0, -4)
	n7 := Simd_scalar_i32_add(n6, s2)
	n8 := Simd_v128_load32_lane(m, n7, 0, 2, n5)
	n9 := Simd_scalar_i32_add(s0, -4)
	n10 := Simd_scalar_i32_add(n9, s3)
	n11 := Simd_v128_load32_lane(m, n10, 0, 3, n8)
	n12 := Simd_scalar_i32_add(s0, -4)
	n13 := Simd_scalar_i32_add(n12, s4)
	n14 := Simd_v128_load32_splat(m, n13, 0)
	return n0[0], n0[1], n1[0], n1[1], n11[0], n11[1], n14[0], n14[1]
}

//go:noinline
func Simd_p_fx258(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p4, p4h})
	n2 := Simd_v128_load32_splat(m, s0, 0)
	n3 := Simd_scalar_i32_add(s0, s1)
	n4 := Simd_v128_load32_lane(m, n3, 0, 1, n2)
	n5 := Simd_scalar_i32_add(s0, s2)
	n6 := Simd_v128_load32_lane(m, n5, 0, 2, n4)
	n7 := Simd_scalar_i32_add(s0, s3)
	n8 := Simd_v128_load32_lane(m, n7, 0, 3, n6)
	n9 := Simd_scalar_i32_add(s0, s4)
	n10 := Simd_v128_load32_splat(m, n9, 0)
	return n0[0], n0[1], n1[0], n1[1], n8[0], n8[1], n10[0], n10[1]
}

//go:noinline
func Simd_p_fx259(m *Module, s0 int32, s1 int32, s2 int32) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_v128_load(m, s2, 0)
	n3 := Simd_scalar_i32_shl(s1, 1)
	n4 := Simd_scalar_i32_add(s0, n3)
	n5 := Simd_v128_load(m, n4, 0)
	n6 := Simd_v128_load(m, s0, 0)
	return n1[0], n1[1], n2[0], n2[1], n5[0], n5[1], n6[0], n6[1]
}

//go:noinline
func Simd_p_fx260(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_sub_sat_u([2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_sub_sat_u([2]uint64{p2, p2h}, [2]uint64{p0, p0h})
	n3 := Simd_v128_or(n1, n2)
	n4 := Simd_v128_load(m, s0, 0)
	n5 := Simd_v128_xor(n4, [2]uint64{p1, p1h})
	return n0[0], n0[1], n4[0], n4[1], n5[0], n5[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx261(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_splat(s1)
	n1 := Simd_v128_load(m, s0, 0)
	n2 := Simd_i8x16_sub_sat_u(n1, [2]uint64{p0, p0h})
	n3 := Simd_i8x16_sub_sat_u([2]uint64{p0, p0h}, n1)
	n4 := Simd_v128_or(n2, n3)
	n5 := Simd_i8x16_max_u([2]uint64{p2, p2h}, n4)
	n6 := Simd_i8x16_sub_sat_u(n5, n0)
	n7 := Simd_i8x16_eq([2]uint64{p1, p1h}, n6)
	n8 := Simd_v128_xor(n1, [2]uint64{p3, p3h})
	return n1[0], n1[1], n4[0], n4[1], n7[0], n7[1], n8[0], n8[1]
}

//go:noinline
func Simd_p_fx262(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_sub_sat_s([2]uint64{p2, p2h}, n0)
	n2 := Simd_v128_load(m, s0, 0)
	n3 := Simd_v128_load(m, s1, 0)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx263(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_sub_sat_s([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_v128_bitselect([2]uint64{p0, p0h}, n0, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_add_sat_s(n1, [2]uint64{p4, p4h})
	n3 := Simd_i8x16_add_sat_s(n2, [2]uint64{p4, p4h})
	n4 := Simd_i8x16_add_sat_s(n3, [2]uint64{p4, p4h})
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx264(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_sub_sat_u([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_sub_sat_u([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_v128_or(n0, n1)
	n3 := Simd_i8x16_max_u([2]uint64{p0, p0h}, n2)
	n4 := Simd_i8x16_sub_sat_u([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n5 := Simd_i8x16_sub_sat_u([2]uint64{p3, p3h}, [2]uint64{p2, p2h})
	n6 := Simd_v128_or(n4, n5)
	n7 := Simd_i8x16_max_u(n3, n6)
	n8 := Simd_i8x16_sub_sat_u([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n9 := Simd_i8x16_sub_sat_u([2]uint64{p5, p5h}, [2]uint64{p4, p4h})
	n10 := Simd_v128_or(n8, n9)
	n11 := Simd_i8x16_max_u(n7, n10)
	n12 := Simd_i8x16_max_u(n11, [2]uint64{p6, p6h})
	return n12[0], n12[1]
}

//go:noinline
func Simd_p_fx265(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_add_sat_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, n0, [2]uint64{p3, p3h})
	n2 := Simd_i16x8_shr_s(n1, 11)
	n3 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, n0, [2]uint64{p4, p4h})
	n4 := Simd_i16x8_shr_s(n3, 11)
	n5 := Simd_i8x16_narrow_i16x8_s(n2, n4)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx266(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_avgr_u(n0, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_add(n1, [2]uint64{p3, p3h})
	n3 := Simd_v128_bitselect(n2, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	n4 := Simd_i8x16_add_sat_s([2]uint64{p5, p5h}, n3)
	n5 := Simd_v128_xor(n4, [2]uint64{p1, p1h})
	n6 := Simd_scalar_i32_add(s0, s1)
	_ = Simd_v128_store(m, n6, 0, n5)
	return n3[0], n3[1]
}

//go:noinline
func Simd_p_fx267(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_add_sat_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p4, p4h})
	n2 := Simd_i16x8_shr_s(n1, 11)
	n3 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p5, p5h})
	n4 := Simd_i16x8_shr_s(n3, 11)
	n5 := Simd_i8x16_narrow_i16x8_s(n2, n4)
	n6 := Simd_i8x16_add_sat_s([2]uint64{p2, p2h}, n5)
	n7 := Simd_v128_xor(n6, [2]uint64{p6, p6h})
	return n7[0], n7[1]
}

//go:noinline
func Simd_p_fx268(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_sub_sat_s([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_v128_xor(n0, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_sub_sat_s([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n3 := Simd_v128_xor(n2, [2]uint64{p3, p3h})
	n4 := Simd_scalar_i32_add(s0, s1)
	_ = Simd_v128_store(m, n4, 0, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s2, 0, n1)
	return n1[0], n1[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx269(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, s5 int32, s6 int32, s7 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load32_splat(m, s0, 0)
	n1 := Simd_scalar_i32_add(s0, s1)
	n2 := Simd_v128_load32_lane(m, n1, 0, 1, n0)
	n3 := Simd_scalar_i32_add(s0, s2)
	n4 := Simd_v128_load32_lane(m, n3, 0, 2, n2)
	n5 := Simd_scalar_i32_add(s0, s3)
	n6 := Simd_v128_load32_lane(m, n5, 0, 3, n4)
	n7 := Simd_scalar_i32_add(s0, s4)
	n8 := Simd_v128_load32_splat(m, n7, 0)
	n9 := Simd_v128_load32_lane(m, s5, 0, 1, n8)
	n10 := Simd_v128_load32_lane(m, s6, 0, 2, n9)
	n11 := Simd_v128_load32_lane(m, s7, 0, 3, n10)
	n12 := Simd_i8x16_shuffle(n6, n11, [2]uint64{p0, p0h})
	n13 := Simd_i8x16_shuffle(n6, n11, [2]uint64{p1, p1h})
	n14 := Simd_i8x16_shuffle(n12, n13, [2]uint64{p2, p2h})
	return n12[0], n12[1], n13[0], n13[1], n14[0], n14[1]
}

//go:noinline
func Simd_p_fx270(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, s5 int32, s6 int32, s7 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_load32_splat(m, s0, 4)
	n2 := Simd_v128_load32_lane(m, s1+4, 0, 1, n1)
	n3 := Simd_v128_load32_lane(m, s2+4, 0, 2, n2)
	n4 := Simd_v128_load32_lane(m, s3+4, 0, 3, n3)
	n5 := Simd_v128_load32_splat(m, s4+4, 0)
	n6 := Simd_v128_load32_lane(m, s5+4, 0, 1, n5)
	n7 := Simd_v128_load32_lane(m, s6+4, 0, 2, n6)
	n8 := Simd_v128_load32_lane(m, s7+4, 0, 3, n7)
	n9 := Simd_i8x16_shuffle(n4, n8, [2]uint64{p2, p2h})
	return n0[0], n0[1], n4[0], n4[1], n8[0], n8[1], n9[0], n9[1]
}

//go:noinline
func Simd_p_fx271(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p4, p4h})
	n2 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p5, p5h})
	n3 := Simd_i8x16_shuffle(n1, n2, [2]uint64{p6, p6h})
	n4 := Simd_v128_load32_splat(m, s0+4, 0)
	return n1[0], n1[1], n2[0], n2[1], n3[0], n3[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx272(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, s5 int32, s6 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load32_lane(m, s0+4, 0, 1, [2]uint64{p0, p0h})
	n1 := Simd_v128_load32_lane(m, s1+4, 0, 2, n0)
	n2 := Simd_v128_load32_lane(m, s2+4, 0, 3, n1)
	n3 := Simd_v128_load32_splat(m, s3+4, 0)
	n4 := Simd_v128_load32_lane(m, s4+4, 0, 1, n3)
	n5 := Simd_v128_load32_lane(m, s5+4, 0, 2, n4)
	n6 := Simd_v128_load32_lane(m, s6+4, 0, 3, n5)
	n7 := Simd_i8x16_shuffle(n2, n6, [2]uint64{p1, p1h})
	n8 := Simd_i8x16_shuffle(n2, n6, [2]uint64{p2, p2h})
	n9 := Simd_i8x16_shuffle(n7, n8, [2]uint64{p3, p3h})
	return n7[0], n7[1], n8[0], n8[1], n9[0], n9[1]
}

//go:noinline
func Simd_p_fx273(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_sub_sat_u([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_sub_sat_u([2]uint64{p3, p3h}, [2]uint64{p2, p2h})
	n3 := Simd_v128_or(n1, n2)
	n4 := Simd_i8x16_shuffle([2]uint64{p4, p4h}, [2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n5 := Simd_i8x16_sub_sat_u(n4, [2]uint64{p0, p0h})
	n6 := Simd_i8x16_sub_sat_u([2]uint64{p0, p0h}, n4)
	n7 := Simd_v128_or(n5, n6)
	return n0[0], n0[1], n3[0], n3[1], n4[0], n4[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx274(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_max_u([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_splat(s0)
	n2 := Simd_i8x16_sub_sat_u(n0, n1)
	n3 := Simd_i8x16_eq([2]uint64{p0, p0h}, n2)
	n4 := Simd_v128_xor([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n5 := Simd_v128_xor([2]uint64{p5, p5h}, [2]uint64{p4, p4h})
	n6 := Simd_i8x16_sub_sat_s(n4, n5)
	return n3[0], n3[1], n4[0], n4[1], n5[0], n5[1], n6[0], n6[1]
}

//go:noinline
func Simd_p_fx275(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_avgr_u(n0, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_add(n1, [2]uint64{p3, p3h})
	n3 := Simd_v128_bitselect(n2, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	n4 := Simd_i8x16_add_sat_s([2]uint64{p5, p5h}, n3)
	n5 := Simd_v128_xor(n4, [2]uint64{p1, p1h})
	return n3[0], n3[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx276(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_sub_sat_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_xor(n0, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_sub_sat_s([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n3 := Simd_v128_xor(n2, [2]uint64{p2, p2h})
	n4 := Simd_i8x16_shuffle(n1, n3, [2]uint64{p5, p5h})
	return n1[0], n1[1], n3[0], n3[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx277(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_v128_xor(n0, [2]uint64{p3, p3h})
	return n0[0], n0[1], n1[0], n1[1]
}

//go:noinline
func Simd_p_fx278(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_sub_sat_u(n0, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_sub_sat_u([2]uint64{p3, p3h}, n0)
	n3 := Simd_v128_or(n2, n1)
	n4 := Simd_i8x16_max_u([2]uint64{p5, p5h}, n3)
	n5 := Simd_v128_xor(n0, [2]uint64{p6, p6h})
	n6 := Simd_i8x16_splat(s0)
	n7 := Simd_i8x16_sub_sat_u(n4, n6)
	n8 := Simd_i8x16_eq([2]uint64{p4, p4h}, n7)
	return n0[0], n0[1], n3[0], n3[1], n8[0], n8[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx279(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_sub_sat_s([2]uint64{p2, p2h}, n0)
	return n0[0], n0[1], n1[0], n1[1]
}

//go:noinline
func Simd_p_fx280(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, s5 int32, s6 int32, s7 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_load32_splat(m, s0, 4)
	n2 := Simd_scalar_i32_add(s0, 4)
	n3 := Simd_scalar_i32_add(n2, s1)
	n4 := Simd_v128_load32_lane(m, n3, 0, 1, n1)
	n5 := Simd_scalar_i32_add(s0, 4)
	n6 := Simd_scalar_i32_add(n5, s2)
	n7 := Simd_v128_load32_lane(m, n6, 0, 2, n4)
	n8 := Simd_scalar_i32_add(s0, 4)
	n9 := Simd_scalar_i32_add(n8, s3)
	n10 := Simd_v128_load32_lane(m, n9, 0, 3, n7)
	n11 := Simd_scalar_i32_add(s0, 4)
	n12 := Simd_scalar_i32_add(n11, s4)
	n13 := Simd_v128_load32_splat(m, n12, 0)
	n14 := Simd_scalar_i32_add(s0, 4)
	n15 := Simd_scalar_i32_add(n14, s5)
	n16 := Simd_v128_load32_lane(m, n15, 0, 1, n13)
	n17 := Simd_scalar_i32_add(s0, 4)
	n18 := Simd_scalar_i32_add(n17, s6)
	n19 := Simd_v128_load32_lane(m, n18, 0, 2, n16)
	n20 := Simd_scalar_i32_add(s0, 4)
	n21 := Simd_scalar_i32_add(n20, s7)
	n22 := Simd_v128_load32_lane(m, n21, 0, 3, n19)
	n23 := Simd_i8x16_shuffle(n10, n22, [2]uint64{p2, p2h})
	return n0[0], n0[1], n10[0], n10[1], n22[0], n22[1], n23[0], n23[1]
}

//go:noinline
func Simd_p_fx281(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p4, p4h})
	n2 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p5, p5h})
	n3 := Simd_i8x16_shuffle(n1, n2, [2]uint64{p6, p6h})
	n4 := Simd_v128_load32_splat(m, s0, 4)
	return n1[0], n1[1], n2[0], n2[1], n3[0], n3[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx282(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, s5 int32, s6 int32, s7 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load32_lane(m, n0, 0, 1, [2]uint64{p0, p0h})
	n2 := Simd_scalar_i32_add(s0, s2)
	n3 := Simd_v128_load32_lane(m, n2, 0, 2, n1)
	n4 := Simd_scalar_i32_add(s0, s3)
	n5 := Simd_v128_load32_lane(m, n4, 0, 3, n3)
	n6 := Simd_scalar_i32_add(s0, s4)
	n7 := Simd_v128_load32_splat(m, n6, 0)
	n8 := Simd_scalar_i32_add(s0, s5)
	n9 := Simd_v128_load32_lane(m, n8, 0, 1, n7)
	n10 := Simd_scalar_i32_add(s0, s6)
	n11 := Simd_v128_load32_lane(m, n10, 0, 2, n9)
	n12 := Simd_scalar_i32_add(s0, s7)
	n13 := Simd_v128_load32_lane(m, n12, 0, 3, n11)
	n14 := Simd_i8x16_shuffle(n5, n13, [2]uint64{p1, p1h})
	n15 := Simd_i8x16_shuffle(n5, n13, [2]uint64{p2, p2h})
	return n14[0], n14[1], n15[0], n15[1]
}

//go:noinline
func Simd_p_fx283(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p4, p4h})
	n3 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, n2, [2]uint64{p6, p6h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx284(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n2 := Simd_i8x16_sub_sat_u(n1, [2]uint64{p5, p5h})
	n3 := Simd_i8x16_sub_sat_u([2]uint64{p5, p5h}, n1)
	n4 := Simd_v128_or(n3, n2)
	return n0[0], n0[1], n1[0], n1[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx285(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_sub_sat_s([2]uint64{p2, p2h}, n0)
	n2 := Simd_i8x16_sub_sat_u([2]uint64{p0, p0h}, [2]uint64{p3, p3h})
	n3 := Simd_i8x16_sub_sat_u([2]uint64{p3, p3h}, [2]uint64{p0, p0h})
	n4 := Simd_v128_or(n2, n3)
	n5 := Simd_i8x16_shuffle([2]uint64{p4, p4h}, [2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	return n0[0], n0[1], n1[0], n1[1], n4[0], n4[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx286(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p4, p4h})
	n2 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p5, p5h})
	return n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx287(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_splat(s4)
	n1 := Simd_v128_load(m, s0, 0)
	n2 := Simd_v128_xor(n1, [2]uint64{p0, p0h})
	n3 := Simd_v128_load(m, s1, 0)
	n4 := Simd_v128_xor(n3, [2]uint64{p0, p0h})
	n5 := Simd_i8x16_sub_sat_s(n4, n2)
	n6 := Simd_i8x16_sub_sat_u(n1, n3)
	n7 := Simd_i8x16_sub_sat_u(n3, n1)
	n8 := Simd_v128_or(n6, n7)
	n9 := Simd_i8x16_add_sat_u(n8, n8)
	n10 := Simd_v128_load(m, s2, 0)
	n11 := Simd_v128_xor(n10, [2]uint64{p0, p0h})
	n12 := Simd_scalar_i32_add(s1, s3)
	n13 := Simd_v128_load(m, n12, 0)
	n14 := Simd_v128_xor(n13, [2]uint64{p0, p0h})
	n15 := Simd_i8x16_sub_sat_s(n11, n14)
	n16 := Simd_i8x16_add_sat_s(n5, n15)
	n17 := Simd_i8x16_add_sat_s(n5, n16)
	n18 := Simd_i8x16_add_sat_s(n5, n17)
	n19 := Simd_i8x16_sub_sat_u(n10, n13)
	n20 := Simd_i8x16_sub_sat_u(n13, n10)
	n21 := Simd_v128_or(n19, n20)
	n22 := Simd_i16x8_shr_u(n21, 1)
	n23 := Simd_v128_and(n22, [2]uint64{p2, p2h})
	n24 := Simd_i8x16_add_sat_u(n9, n23)
	n25 := Simd_i8x16_sub_sat_u(n24, n0)
	n26 := Simd_i8x16_eq([2]uint64{p1, p1h}, n25)
	n27 := Simd_v128_bitselect(n18, [2]uint64{p1, p1h}, n26)
	n28 := Simd_i8x16_add_sat_s(n27, [2]uint64{p3, p3h})
	return n2[0], n2[1], n4[0], n4[1], n27[0], n27[1], n28[0], n28[1]
}

//go:noinline
func Simd_p_fx288(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n1 := Simd_i16x8_shr_s(n0, 11)
	n2 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	n3 := Simd_i16x8_shr_s(n2, 11)
	n4 := Simd_i8x16_narrow_i16x8_s(n1, n3)
	n5 := Simd_i8x16_add_sat_s([2]uint64{p0, p0h}, n4)
	n6 := Simd_v128_xor(n5, [2]uint64{p5, p5h})
	_ = Simd_v128_store(m, s0, 0, n6)
	return
}

//go:noinline
func Simd_p_fx289(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) {
	n0 := Simd_i8x16_add_sat_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p4, p4h})
	n2 := Simd_i16x8_shr_s(n1, 11)
	n3 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p5, p5h})
	n4 := Simd_i16x8_shr_s(n3, 11)
	n5 := Simd_i8x16_narrow_i16x8_s(n2, n4)
	n6 := Simd_i8x16_sub_sat_s([2]uint64{p2, p2h}, n5)
	n7 := Simd_v128_xor(n6, [2]uint64{p6, p6h})
	_ = Simd_v128_store(m, s0, 0, n7)
	return
}

//go:noinline
func Simd_p_fx290(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_v128_xor(n0, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle([2]uint64{p4, p4h}, [2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx291(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p4, p4h})
	n2 := Simd_v128_xor(n1, [2]uint64{p5, p5h})
	n3 := Simd_i8x16_sub_sat_s(n2, [2]uint64{p6, p6h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx292(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_v128_xor([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_v128_xor([2]uint64{p3, p3h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_sub_sat_s(n0, n1)
	n3 := Simd_i8x16_add_sat_s([2]uint64{p0, p0h}, n2)
	n4 := Simd_i8x16_add_sat_s([2]uint64{p0, p0h}, n3)
	n5 := Simd_i8x16_add_sat_s([2]uint64{p0, p0h}, n4)
	n6 := Simd_i8x16_add_sat_u([2]uint64{p5, p5h}, [2]uint64{p5, p5h})
	n7 := Simd_i8x16_sub_sat_u([2]uint64{p1, p1h}, [2]uint64{p3, p3h})
	n8 := Simd_i8x16_sub_sat_u([2]uint64{p3, p3h}, [2]uint64{p1, p1h})
	n9 := Simd_v128_or(n7, n8)
	n10 := Simd_i16x8_shr_u(n9, 1)
	n11 := Simd_v128_and(n10, [2]uint64{p6, p6h})
	n12 := Simd_i8x16_add_sat_u(n6, n11)
	n13 := Simd_i8x16_splat(s0)
	n14 := Simd_i8x16_sub_sat_u(n12, n13)
	n15 := Simd_i8x16_eq([2]uint64{p4, p4h}, n14)
	n16 := Simd_v128_bitselect(n5, [2]uint64{p4, p4h}, n15)
	return n16[0], n16[1]
}

//go:noinline
func Simd_p_fx293(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_add_sat_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p4, p4h})
	n2 := Simd_i16x8_shr_s(n1, 11)
	n3 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p5, p5h})
	n4 := Simd_i16x8_shr_s(n3, 11)
	n5 := Simd_i8x16_narrow_i16x8_s(n2, n4)
	n6 := Simd_i8x16_sub_sat_s([2]uint64{p2, p2h}, n5)
	n7 := Simd_v128_xor(n6, [2]uint64{p6, p6h})
	return n7[0], n7[1]
}

//go:noinline
func Simd_p_fx294(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_splat(s4)
	n1 := Simd_v128_load(m, s0, 0)
	n2 := Simd_v128_xor(n1, [2]uint64{p0, p0h})
	n3 := Simd_v128_load(m, s1, 0)
	n4 := Simd_v128_xor(n3, [2]uint64{p0, p0h})
	n5 := Simd_i8x16_sub_sat_s(n4, n2)
	n6 := Simd_i8x16_sub_sat_u(n1, n3)
	n7 := Simd_i8x16_sub_sat_u(n3, n1)
	n8 := Simd_v128_or(n6, n7)
	n9 := Simd_i8x16_add_sat_u(n8, n8)
	n10 := Simd_v128_load(m, s2, 0)
	n11 := Simd_v128_xor(n10, [2]uint64{p0, p0h})
	n12 := Simd_scalar_i32_add(s1, s3)
	n13 := Simd_v128_load(m, n12, 0)
	n14 := Simd_v128_xor(n13, [2]uint64{p0, p0h})
	n15 := Simd_i8x16_sub_sat_s(n11, n14)
	n16 := Simd_i8x16_add_sat_s(n5, n15)
	n17 := Simd_i8x16_add_sat_s(n5, n16)
	n18 := Simd_i8x16_add_sat_s(n5, n17)
	n19 := Simd_i8x16_sub_sat_u(n10, n13)
	n20 := Simd_i8x16_sub_sat_u(n13, n10)
	n21 := Simd_v128_or(n19, n20)
	n22 := Simd_i16x8_shr_u(n21, 1)
	n23 := Simd_v128_and(n22, [2]uint64{p2, p2h})
	n24 := Simd_i8x16_add_sat_u(n9, n23)
	n25 := Simd_i8x16_sub_sat_u(n24, n0)
	n26 := Simd_i8x16_eq([2]uint64{p1, p1h}, n25)
	n27 := Simd_v128_bitselect(n18, [2]uint64{p1, p1h}, n26)
	return n2[0], n2[1], n4[0], n4[1], n0[0], n0[1], n27[0], n27[1]
}

//go:noinline
func Simd_p_fx295(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) {
	n0 := Simd_i8x16_add_sat_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p4, p4h})
	n2 := Simd_i16x8_shr_s(n1, 11)
	n3 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p5, p5h})
	n4 := Simd_i16x8_shr_s(n3, 11)
	n5 := Simd_i8x16_narrow_i16x8_s(n2, n4)
	n6 := Simd_i8x16_add_sat_s([2]uint64{p2, p2h}, n5)
	n7 := Simd_v128_xor(n6, [2]uint64{p6, p6h})
	_ = Simd_v128_store(m, s0, 0, n7)
	return
}

//go:noinline
func Simd_p_fx296(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_xor(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_load(m, s1, 0)
	n3 := Simd_v128_xor(n2, [2]uint64{p0, p0h})
	n4 := Simd_i8x16_sub_sat_s(n3, n1)
	n5 := Simd_i8x16_sub_sat_u(n0, n2)
	n6 := Simd_i8x16_sub_sat_u(n2, n0)
	n7 := Simd_v128_or(n5, n6)
	n8 := Simd_i8x16_add_sat_u(n7, n7)
	n9 := Simd_v128_load(m, s2, 0)
	n10 := Simd_v128_xor(n9, [2]uint64{p0, p0h})
	n11 := Simd_scalar_i32_add(s1, s3)
	n12 := Simd_v128_load(m, n11, 0)
	n13 := Simd_v128_xor(n12, [2]uint64{p0, p0h})
	n14 := Simd_i8x16_sub_sat_s(n10, n13)
	n15 := Simd_i8x16_add_sat_s(n4, n14)
	n16 := Simd_i8x16_add_sat_s(n4, n15)
	n17 := Simd_i8x16_add_sat_s(n4, n16)
	n18 := Simd_i8x16_sub_sat_u(n9, n12)
	n19 := Simd_i8x16_sub_sat_u(n12, n9)
	n20 := Simd_v128_or(n18, n19)
	n21 := Simd_i16x8_shr_u(n20, 1)
	n22 := Simd_v128_and(n21, [2]uint64{p2, p2h})
	n23 := Simd_i8x16_add_sat_u(n8, n22)
	n24 := Simd_i8x16_sub_sat_u(n23, [2]uint64{p3, p3h})
	n25 := Simd_i8x16_eq([2]uint64{p1, p1h}, n24)
	n26 := Simd_v128_bitselect(n17, [2]uint64{p1, p1h}, n25)
	n27 := Simd_i8x16_add_sat_s(n26, [2]uint64{p4, p4h})
	return n1[0], n1[1], n3[0], n3[1], n26[0], n26[1], n27[0], n27[1]
}

//go:noinline
func Simd_p_fx297(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_avgr_u([2]uint64{p0, p0h}, n0)
	return n0[0], n0[1], n1[0], n1[1]
}

//go:noinline
func Simd_p_fx298(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{650777868590383874, 1229499251294997258})
	n1 := Simd_i8x16_avgr_u(n0, [2]uint64{p2, p2h})
	n2 := Simd_v128_xor(n1, [2]uint64{p3, p3h})
	n3 := Simd_i8x16_avgr_u([2]uint64{p3, p3h}, n1)
	n4 := Simd_v128_xor([2]uint64{p2, p2h}, n0)
	n5 := Simd_v128_xor([2]uint64{p2, p2h}, [2]uint64{p0, p0h})
	n6 := Simd_v128_or(n5, n4)
	n7 := Simd_v128_and(n2, n6)
	n8 := Simd_v128_and(n7, [2]uint64{p4, p4h})
	n9 := Simd_i8x16_sub_sat_u(n3, n8)
	return n9[0], n9[1]
}

//go:noinline
func Simd_p_fx299(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_avgr_u([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_xor([2]uint64{p1, p1h}, [2]uint64{p0, p0h})
	n2 := Simd_v128_and(n1, [2]uint64{p2, p2h})
	n3 := Simd_i8x16_sub_sat_u(n0, n2)
	n4 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n5 := Simd_i8x16_avgr_u(n3, n4)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx300(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{578437695752307201, 1157159078456920585})
	n1 := Simd_i8x16_avgr_u([2]uint64{p0, p0h}, n0)
	return n0[0], n0[1], n1[0], n1[1]
}

//go:noinline
func Simd_p_fx301(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_avgr_u([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_v128_xor([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_v128_and(n1, [2]uint64{p3, p3h})
	n3 := Simd_i8x16_sub_sat_u(n0, n2)
	n4 := Simd_i8x16_avgr_u([2]uint64{p0, p0h}, n3)
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx302(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1374179596971150604, 1952900979675763988})
	n1 := Simd_v128_or(n0, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p0, p0h}, [2]uint64{p3, p3h})
	n3 := Simd_i8x16_avgr_u(n2, n1)
	n4 := Simd_v128_xor(n2, n1)
	n5 := Simd_v128_and(n4, [2]uint64{p4, p4h})
	n6 := Simd_i8x16_sub_sat_u(n3, n5)
	n7 := Simd_i8x16_shuffle(n1, [2]uint64{p0, p0h}, [2]uint64{p5, p5h})
	n8 := Simd_i8x16_avgr_u(n6, n7)
	return n8[0], n8[1]
}

//go:noinline
func Simd_p_fx303(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{650777868590383874, 1229499251294997258})
	n1 := Simd_v128_xor(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_and(n1, [2]uint64{p2, p2h})
	n3 := Simd_i8x16_avgr_u([2]uint64{p0, p0h}, n0)
	n4 := Simd_i8x16_sub_sat_u(n3, n2)
	n5 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{578437695752307201, 1157159078456920585})
	n6 := Simd_i8x16_avgr_u(n4, n5)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx304(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_v128_load32_zero(m, s0+-32, 0)
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p0, p0h}, [2]uint64{1369958511735279616, 1659319203087586308})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx305(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i8x16_narrow_i16x8_u(n1, [2]uint64{p1, p1h})
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx306(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0+-32, 0)
	n1 := Simd_i8x16_sub_sat_u(n0, [2]uint64{p0, p0h})
	n2 := Simd_i8x16_sub_sat_u([2]uint64{p0, p0h}, n0)
	n3 := Simd_v128_or(n2, n1)
	n4 := Simd_i16x8_shr_u(n3, 8)
	n5 := Simd_v128_and(n3, [2]uint64{p1, p1h})
	n6 := Simd_i16x8_add(n4, n5)
	n7 := Simd_i32x4_shl(n6, 16)
	n8 := Simd_i16x8_add(n6, n7)
	n9 := Simd_i64x2_shl(n8, 32)
	n10 := Simd_i16x8_add(n8, n9)
	n11 := Simd_i64x2_shr_u(n10, 48)
	return n11[0], n11[1]
}

//go:noinline
func Simd_p_fx307(m *Module, s0 int32, s1 int32) {
	n0 := Simd_i8x16_splat(s0)
	_ = Simd_v128_store(m, s1, 480, n0)
	_ = Simd_v128_store(m, s1, 448, n0)
	_ = Simd_v128_store(m, s1, 416, n0)
	_ = Simd_v128_store(m, s1, 384, n0)
	_ = Simd_v128_store(m, s1, 352, n0)
	_ = Simd_v128_store(m, s1, 320, n0)
	_ = Simd_v128_store(m, s1, 288, n0)
	_ = Simd_v128_store(m, s1, 256, n0)
	_ = Simd_v128_store(m, s1, 224, n0)
	_ = Simd_v128_store(m, s1, 192, n0)
	_ = Simd_v128_store(m, s1, 160, n0)
	_ = Simd_v128_store(m, s1, 128, n0)
	_ = Simd_v128_store(m, s1, 96, n0)
	_ = Simd_v128_store(m, s1, 64, n0)
	_ = Simd_v128_store(m, s1, 32, n0)
	_ = Simd_v128_store(m, s1, 0, n0)
	return
}

//go:noinline
func Simd_p_fx308(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{216736831764039944, 216736831629295872})
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx309(m *Module, s0 int32) {
	n0 := Simd_v128_load(m, s0+-32, 0)
	_ = Simd_v128_store(m, s0, 480, n0)
	_ = Simd_v128_store(m, s0, 448, n0)
	_ = Simd_v128_store(m, s0, 416, n0)
	_ = Simd_v128_store(m, s0, 384, n0)
	_ = Simd_v128_store(m, s0, 352, n0)
	_ = Simd_v128_store(m, s0, 320, n0)
	_ = Simd_v128_store(m, s0, 288, n0)
	_ = Simd_v128_store(m, s0, 256, n0)
	_ = Simd_v128_store(m, s0, 224, n0)
	_ = Simd_v128_store(m, s0, 192, n0)
	_ = Simd_v128_store(m, s0, 160, n0)
	_ = Simd_v128_store(m, s0, 128, n0)
	_ = Simd_v128_store(m, s0, 96, n0)
	_ = Simd_v128_store(m, s0, 64, n0)
	_ = Simd_v128_store(m, s0, 32, n0)
	_ = Simd_v128_store(m, s0, 0, n0)
	return
}

//go:noinline
func Simd_p_fx310(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1369958511735279616, 1659319203087586308})
	n2 := Simd_i16x8_add(n0, n1)
	n3 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1948679894439893000, 2238040585792199692})
	n4 := Simd_i16x8_add(n0, n3)
	n5 := Simd_i8x16_narrow_i16x8_u(n2, n4)
	_ = Simd_v128_store(m, s1, 0, n5)
	return
}

//go:noinline
func Simd_p_fx311(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_sub_sat_u([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_sub_sat_u([2]uint64{p1, p1h}, [2]uint64{p0, p0h})
	n2 := Simd_v128_or(n0, n1)
	n3 := Simd_i16x8_shr_u(n2, 8)
	n4 := Simd_v128_and(n2, [2]uint64{p2, p2h})
	n5 := Simd_i16x8_add(n3, n4)
	n6 := Simd_i32x4_shl(n5, 16)
	n7 := Simd_i16x8_add(n5, n6)
	return n7[0], n7[1]
}

//go:noinline
func Simd_p_fx312(m *Module, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_i64x2_shl([2]uint64{p0, p0h}, 32)
	n1 := Simd_i16x8_add([2]uint64{p0, p0h}, n0)
	n2 := Simd_i64x2_shr_u(n1, 48)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx313(m *Module, s0 int32, p0, p0h uint64) {
	n0 := Simd_v128_load32_zero(m, s0, 28)
	n1 := Simd_i8x16_swizzle(n0, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 32, n1)
	n3 := Simd_v128_load32_zero(m, s0, 60)
	n4 := Simd_i8x16_swizzle(n3, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 64, n4)
	n6 := Simd_v128_load32_zero(m, s0, 92)
	n7 := Simd_i8x16_swizzle(n6, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 96, n7)
	n9 := Simd_v128_load32_zero(m, s0, 124)
	n10 := Simd_i8x16_swizzle(n9, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 128, n10)
	n12 := Simd_v128_load32_zero(m, s0, 156)
	n13 := Simd_i8x16_swizzle(n12, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 160, n13)
	n15 := Simd_v128_load32_zero(m, s0, 188)
	n16 := Simd_i8x16_swizzle(n15, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 192, n16)
	n18 := Simd_v128_load32_zero(m, s0, 220)
	n19 := Simd_i8x16_swizzle(n18, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 224, n19)
	n21 := Simd_v128_load32_zero(m, s0, 252)
	n22 := Simd_i8x16_swizzle(n21, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 256, n22)
	n24 := Simd_v128_load32_zero(m, s0+-4, 0)
	n25 := Simd_i8x16_swizzle(n24, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 0, n25)
	n27 := Simd_v128_load32_zero(m, s0, 284)
	n28 := Simd_i8x16_swizzle(n27, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 288, n28)
	n30 := Simd_v128_load32_zero(m, s0, 316)
	n31 := Simd_i8x16_swizzle(n30, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 320, n31)
	n33 := Simd_v128_load32_zero(m, s0, 348)
	n34 := Simd_i8x16_swizzle(n33, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 352, n34)
	n36 := Simd_v128_load32_zero(m, s0, 380)
	n37 := Simd_i8x16_swizzle(n36, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 384, n37)
	n39 := Simd_v128_load32_zero(m, s0, 412)
	n40 := Simd_i8x16_swizzle(n39, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 416, n40)
	n42 := Simd_v128_load32_zero(m, s0, 444)
	n43 := Simd_i8x16_swizzle(n42, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 448, n43)
	n45 := Simd_v128_load32_zero(m, s0, 476)
	n46 := Simd_i8x16_swizzle(n45, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 480, n46)
	return
}

//go:noinline
func Simd_p_fx314(m *Module, s0 int32) {
	n0 := Simd_v128_load(m, s0+-32, 0)
	_ = Simd_v128_store(m, s0, 0, n0)
	_ = Simd_v128_store(m, s0, 32, n0)
	_ = Simd_v128_store(m, s0, 64, n0)
	_ = Simd_v128_store(m, s0, 96, n0)
	_ = Simd_v128_store(m, s0, 128, n0)
	_ = Simd_v128_store(m, s0, 256, n0)
	_ = Simd_v128_store(m, s0, 224, n0)
	_ = Simd_v128_store(m, s0, 192, n0)
	_ = Simd_v128_store(m, s0, 160, n0)
	_ = Simd_v128_store(m, s0, 384, n0)
	_ = Simd_v128_store(m, s0, 352, n0)
	_ = Simd_v128_store(m, s0, 320, n0)
	_ = Simd_v128_store(m, s0, 288, n0)
	_ = Simd_v128_store(m, s0, 480, n0)
	_ = Simd_v128_store(m, s0, 448, n0)
	_ = Simd_v128_store(m, s0, 416, n0)
	return
}

//go:noinline
func Simd_p_fx315(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i8x16_sub_sat_u(n0, [2]uint64{p0, p0h})
	n2 := Simd_i8x16_sub_sat_u([2]uint64{p0, p0h}, n0)
	n3 := Simd_v128_or(n2, n1)
	n4 := Simd_i16x8_shr_u(n3, 8)
	n5 := Simd_v128_and(n3, [2]uint64{p1, p1h})
	n6 := Simd_i16x8_add(n4, n5)
	n7 := Simd_i32x4_shl(n6, 16)
	n8 := Simd_i16x8_add(n6, n7)
	n9 := Simd_i64x2_shl(n8, 32)
	n10 := Simd_i16x8_add(n8, n9)
	n11 := Simd_i64x2_shr_u(n10, 48)
	return n11[0], n11[1]
}

//go:noinline
func Simd_p_fx316(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{216736831764039944, 216736831629295872})
	n1 := Simd_i32x4_add(n0, [2]uint64{p0, p0h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx317(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i8x16_sub_sat_u(n0, [2]uint64{p0, p0h})
	n2 := Simd_i8x16_sub_sat_u([2]uint64{p0, p0h}, n0)
	n3 := Simd_v128_or(n2, n1)
	n4 := Simd_i16x8_shr_u(n3, 8)
	n5 := Simd_v128_and(n3, [2]uint64{p1, p1h})
	n6 := Simd_i16x8_add(n4, n5)
	n7 := Simd_i32x4_shl(n6, 16)
	n8 := Simd_i16x8_add(n6, n7)
	n9 := Simd_i64x2_shl(n8, 32)
	n10 := Simd_i16x8_add(n8, n9)
	n11 := Simd_i64x2_shr_u(n10, 48)
	n12 := Simd_v128_load(m, s1, 0)
	n13 := Simd_i8x16_sub_sat_u(n12, [2]uint64{p0, p0h})
	n14 := Simd_i8x16_sub_sat_u([2]uint64{p0, p0h}, n12)
	n15 := Simd_v128_or(n14, n13)
	n16 := Simd_i16x8_shr_u(n15, 8)
	n17 := Simd_v128_and(n15, [2]uint64{p1, p1h})
	n18 := Simd_i16x8_add(n16, n17)
	n19 := Simd_i32x4_shl(n18, 16)
	n20 := Simd_i16x8_add(n18, n19)
	n21 := Simd_i64x2_shl(n20, 32)
	n22 := Simd_i16x8_add(n20, n21)
	n23 := Simd_i64x2_shr_u(n22, 48)
	return n11[0], n11[1], n23[0], n23[1]
}

//go:noinline
func Simd_p_fx318(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i32x4_add(n0, [2]uint64{p0, p0h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx319(m *Module, s0 int32, s1 int32) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	_ = Simd_v128_store(m, s1, 768, n0)
	_ = Simd_v128_store(m, s1, 736, n0)
	_ = Simd_v128_store(m, s1, 704, n0)
	_ = Simd_v128_store(m, s1, 672, n0)
	_ = Simd_v128_store(m, s1, 640, n0)
	_ = Simd_v128_store(m, s1, 608, n0)
	_ = Simd_v128_store(m, s1, 576, n0)
	_ = Simd_v128_store(m, s1, 544, n0)
	_ = Simd_v128_store(m, s1, 512, n0)
	return n0[0], n0[1]
}

//go:noinline
func Simd_p_fx320(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	_ = Simd_v128_store(m, s1, 480, n0)
	_ = Simd_v128_store(m, s1, 448, n0)
	_ = Simd_v128_store(m, s1, 416, n0)
	_ = Simd_v128_store(m, s1, 384, n0)
	_ = Simd_v128_store(m, s1, 352, n0)
	_ = Simd_v128_store(m, s1, 320, n0)
	_ = Simd_v128_store(m, s1, 288, n0)
	_ = Simd_v128_store(m, s1, 256, n0)
	_ = Simd_v128_store(m, s1, 224, n0)
	_ = Simd_v128_store(m, s1, 192, n0)
	_ = Simd_v128_store(m, s1, 160, n0)
	_ = Simd_v128_store(m, s1, 128, n0)
	_ = Simd_v128_store(m, s1, 96, n0)
	_ = Simd_v128_store(m, s1, 64, n0)
	_ = Simd_v128_store(m, s1, 32, n0)
	_ = Simd_v128_store(m, s1, 0, n0)
	return
}

//go:noinline
func Simd_p_fx321(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_v128_load(m, s1, 0)
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p0, p0h}, [2]uint64{1369958511735279616, 1659319203087586308})
	n3 := Simd_i8x16_shuffle(n1, [2]uint64{p0, p0h}, [2]uint64{1948679894439893000, 2238040585792199692})
	n4 := Simd_i16x8_add(n0, n2)
	n5 := Simd_i16x8_add(n0, n3)
	n6 := Simd_i8x16_narrow_i16x8_u(n4, n5)
	_ = Simd_v128_store(m, s2, 0, n6)
	return n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx322(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_add(n0, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_narrow_i16x8_u(n1, n2)
	_ = Simd_v128_store(m, s1, 32, n3)
	return
}

//go:noinline
func Simd_p_fx323(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_add(n0, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_narrow_i16x8_u(n1, n2)
	_ = Simd_v128_store(m, s1, 64, n3)
	return
}

//go:noinline
func Simd_p_fx324(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_add(n0, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_narrow_i16x8_u(n1, n2)
	_ = Simd_v128_store(m, s1, 96, n3)
	return
}

//go:noinline
func Simd_p_fx325(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_add(n0, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_narrow_i16x8_u(n1, n2)
	_ = Simd_v128_store(m, s1, 128, n3)
	return
}

//go:noinline
func Simd_p_fx326(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_add(n0, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_narrow_i16x8_u(n1, n2)
	_ = Simd_v128_store(m, s1, 160, n3)
	return
}

//go:noinline
func Simd_p_fx327(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_add(n0, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_narrow_i16x8_u(n1, n2)
	_ = Simd_v128_store(m, s1, 192, n3)
	return
}

//go:noinline
func Simd_p_fx328(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_add(n0, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_narrow_i16x8_u(n1, n2)
	_ = Simd_v128_store(m, s1, 224, n3)
	return
}

//go:noinline
func Simd_p_fx329(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_add(n0, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_narrow_i16x8_u(n1, n2)
	_ = Simd_v128_store(m, s1, 256, n3)
	return
}

//go:noinline
func Simd_p_fx330(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_add(n0, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_narrow_i16x8_u(n1, n2)
	_ = Simd_v128_store(m, s1, 288, n3)
	return
}

//go:noinline
func Simd_p_fx331(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_add(n0, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_narrow_i16x8_u(n1, n2)
	_ = Simd_v128_store(m, s1, 320, n3)
	return
}

//go:noinline
func Simd_p_fx332(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_add(n0, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_narrow_i16x8_u(n1, n2)
	_ = Simd_v128_store(m, s1, 352, n3)
	return
}

//go:noinline
func Simd_p_fx333(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_add(n0, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_narrow_i16x8_u(n1, n2)
	_ = Simd_v128_store(m, s1, 384, n3)
	return
}

//go:noinline
func Simd_p_fx334(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_add(n0, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_narrow_i16x8_u(n1, n2)
	_ = Simd_v128_store(m, s1, 416, n3)
	return
}

//go:noinline
func Simd_p_fx335(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_add(n0, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_narrow_i16x8_u(n1, n2)
	_ = Simd_v128_store(m, s1, 448, n3)
	return
}

//go:noinline
func Simd_p_fx336(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i16x8_splat(s0)
	n1 := Simd_i16x8_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_add(n0, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_narrow_i16x8_u(n1, n2)
	_ = Simd_v128_store(m, s1, 480, n3)
	return
}

//go:noinline
func Simd_p_fx337(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n2 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, [2]uint64{p6, p6h}, [2]uint64{p2, p2h})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n4 := Simd_i16x8_sub(n1, n3)
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx338(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p3, p3h})
	n2 := Simd_i16x8_add(n0, n1)
	n3 := Simd_i32x4_dot_i16x8_s(n2, [2]uint64{p4, p4h})
	n4 := Simd_i32x4_dot_i16x8_s(n2, [2]uint64{p5, p5h})
	n5 := Simd_i16x8_narrow_i32x4_s(n3, n4)
	n6 := Simd_i16x8_sub(n0, n1)
	return n5[0], n5[1], n6[0], n6[1]
}

//go:noinline
func Simd_p_fx339(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_dot_i16x8_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_add(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_shr_s(n1, 9)
	n3 := Simd_i32x4_dot_i16x8_s([2]uint64{p0, p0h}, [2]uint64{p3, p3h})
	n4 := Simd_i32x4_add(n3, [2]uint64{p4, p4h})
	n5 := Simd_i32x4_shr_s(n4, 9)
	n6 := Simd_i16x8_narrow_i32x4_s(n2, n5)
	n7 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, n6, [2]uint64{p6, p6h})
	n8 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, n6, [2]uint64{1952885526417115400, 2242246217769422092})
	return n7[0], n7[1], n8[0], n8[1]
}

//go:noinline
func Simd_p_fx340(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1952900979473647880, 2242261670825954572})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p2, p2h}, [2]uint64{1084818905618843912, 506097522914230528})
	n2 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1374179596769034496, 1663540288121341188})
	n3 := Simd_i16x8_add(n1, n2)
	n4 := Simd_i16x8_add(n3, [2]uint64{p3, p3h})
	n5 := Simd_i8x16_shuffle(n3, n3, [2]uint64{p4, p4h})
	return n1[0], n1[1], n2[0], n2[1], n4[0], n4[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx341(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_shuffle(n0, n0, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle(n1, n0, [2]uint64{p3, p3h})
	n3 := Simd_i32x4_dot_i16x8_s(n2, [2]uint64{p4, p4h})
	n4 := Simd_i32x4_add(n3, [2]uint64{p5, p5h})
	n5 := Simd_i32x4_shr_s(n4, 16)
	return n2[0], n2[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx342(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_shr_s(n0, 4)
	n2 := Simd_i16x8_narrow_i32x4_s([2]uint64{p2, p2h}, [2]uint64{p2, p2h})
	n3 := Simd_i16x8_eq([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n4 := Simd_i16x8_add(n2, n3)
	n5 := Simd_i8x16_shuffle(n1, n4, [2]uint64{p5, p5h})
	n6 := Simd_i16x8_abs(n5)
	n7 := Simd_i16x8_shr_s(n6, 3)
	n8 := Simd_i16x8_min_s(n7, [2]uint64{p6, p6h})
	return n8[0], n8[1]
}

//go:noinline
func Simd_p_fx343(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_dot_i16x8_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_add(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_shr_s(n1, 16)
	n3 := Simd_i16x8_narrow_i32x4_s(n2, n2)
	n4 := Simd_i16x8_sub([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n5 := Simd_i16x8_shr_s(n4, 4)
	n6 := Simd_i8x16_shuffle(n5, n3, [2]uint64{p5, p5h})
	n7 := Simd_i16x8_abs(n6)
	n8 := Simd_i16x8_shr_s(n7, 3)
	n9 := Simd_i16x8_min_s(n8, [2]uint64{p6, p6h})
	return n9[0], n9[1]
}

//go:noinline
func Simd_p_fx344(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{506097522914230528, 1663540288323457296})
	n1 := Simd_i8x16_sub_sat_u(n0, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_sub_sat_u([2]uint64{p2, p2h}, n0)
	n3 := Simd_v128_or(n2, n1)
	n4 := Simd_i16x8_shr_u(n3, 8)
	n5 := Simd_v128_and(n3, [2]uint64{p3, p3h})
	n6 := Simd_i16x8_add(n4, n5)
	n7 := Simd_i32x4_shl(n6, 16)
	n8 := Simd_i16x8_add(n6, n7)
	n9 := Simd_i64x2_shl(n8, 32)
	n10 := Simd_i16x8_add(n8, n9)
	n11 := Simd_i64x2_shr_u(n10, 48)
	return n11[0], n11[1]
}

//go:noinline
func Simd_p_fx345(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_v128_load32_zero(m, s0, 0)
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p0, p0h}, [2]uint64{1369958511735279616, 1659319203087586308})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx346(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_avgr_u(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_xor(n0, [2]uint64{p0, p0h})
	n3 := Simd_v128_and(n2, [2]uint64{p3, p3h})
	n4 := Simd_i8x16_sub_sat_u(n1, n3)
	n5 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p4, p4h})
	n6 := Simd_i8x16_avgr_u(n4, n5)
	n7 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, [2]uint64{p1, p1h}, [2]uint64{p4, p4h})
	n8 := Simd_i8x16_avgr_u([2]uint64{p5, p5h}, n7)
	n9 := Simd_i8x16_shuffle([2]uint64{p6, p6h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	return n6[0], n6[1], n7[0], n7[1], n8[0], n8[1], n9[0], n9[1]
}

//go:noinline
func Simd_p_fx347(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_avgr_u([2]uint64{p0, p0h}, [2]uint64{p3, p3h})
	n2 := Simd_v128_xor([2]uint64{p3, p3h}, [2]uint64{p0, p0h})
	n3 := Simd_v128_and(n2, [2]uint64{p4, p4h})
	n4 := Simd_i8x16_sub_sat_u(n1, n3)
	n5 := Simd_i8x16_avgr_u(n4, n0)
	n6 := Simd_i8x16_avgr_u([2]uint64{p0, p0h}, n0)
	return n0[0], n0[1], n5[0], n5[1], n6[0], n6[1]
}

//go:noinline
func Simd_p_fx348(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p5, p5h})
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx349(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_avgr_u([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_xor([2]uint64{p1, p1h}, [2]uint64{p0, p0h})
	n2 := Simd_v128_and(n1, [2]uint64{p2, p2h})
	n3 := Simd_i8x16_sub_sat_u(n0, n2)
	n4 := Simd_i8x16_avgr_u(n3, [2]uint64{p3, p3h})
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx350(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_avgr_u([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_xor(n0, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_avgr_u([2]uint64{p2, p2h}, n0)
	n3 := Simd_v128_xor([2]uint64{p1, p1h}, [2]uint64{p3, p3h})
	n4 := Simd_v128_xor([2]uint64{p1, p1h}, [2]uint64{p0, p0h})
	n5 := Simd_v128_or(n3, n4)
	n6 := Simd_v128_and(n1, n5)
	n7 := Simd_v128_and(n6, [2]uint64{p4, p4h})
	n8 := Simd_i8x16_sub_sat_u(n2, n7)
	return n8[0], n8[1]
}

//go:noinline
func Simd_p_fx351(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n1 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n0, [2]uint64{p4, p4h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx352(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_v128_xor(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_and(n1, [2]uint64{p3, p3h})
	n3 := Simd_i8x16_avgr_u([2]uint64{p0, p0h}, n0)
	n4 := Simd_i8x16_sub_sat_u(n3, n2)
	n5 := Simd_i8x16_avgr_u(n4, [2]uint64{p4, p4h})
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx353(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx354(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx355(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 96, 16, 208)
	n1 := Simd_v128_load_nc(m, s0, 208)
	n2 := Simd_v128_load(m, s1, 16)
	n3 := Simd_i16x8_abs(n2)
	n4 := Simd_i16x8_add(n1, n3)
	n5 := Simd_i16x8_shr_s(n2, 15)
	n6 := Simd_v128_load_nc(m, s0, 48)
	n7 := Simd_i16x8_mul(n4, n6)
	n8 := Simd_i32x4_extmul_low_i16x8_u(n4, n6)
	n9 := Simd_i32x4_extmul_high_i16x8_u(n4, n6)
	n10 := Simd_i8x16_shuffle(n8, n9, [2]uint64{p0, p0h})
	n11 := Simd_i8x16_shuffle(n7, n10, [2]uint64{p1, p1h})
	n12 := Simd_i32x4_add(n0, n11)
	n13 := Simd_i32x4_shr_s(n12, 17)
	n14 := Simd_i8x16_shuffle(n7, n10, [2]uint64{p2, p2h})
	n15 := Simd_v128_load_nc(m, s0, 112)
	n16 := Simd_i32x4_add(n15, n14)
	n17 := Simd_i32x4_shr_s(n16, 17)
	n18 := Simd_i16x8_narrow_i32x4_s(n13, n17)
	n19 := Simd_i16x8_min_s(n18, [2]uint64{p3, p3h})
	n20 := Simd_v128_xor(n19, n5)
	n21 := Simd_i16x8_sub(n20, n5)
	n22 := Simd_i8x16_shuffle(n21, [2]uint64{p4, p4h}, [2]uint64{506093107721536258, 1084818905618843912})
	n23 := Simd_i8x16_shuffle(n22, [2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n24 := Simd_i8x16_shuffle(n23, [2]uint64{p4, p4h}, [2]uint64{216743454502551808, 1084818905618843912})
	n25 := Simd_v128_load_nc(m, s0, 16)
	n26 := Simd_i16x8_mul(n21, n25)
	_ = Simd_v128_store(m, s1, 16, n26)
	n28 := Simd_v128_load_rng(m, s0, 64, 0, 208)
	n29 := Simd_v128_load_nc(m, s0, 192)
	n30 := Simd_v128_load(m, s1, 0)
	n31 := Simd_i16x8_abs(n30)
	n32 := Simd_i16x8_add(n29, n31)
	n33 := Simd_i16x8_shr_s(n30, 15)
	n34 := Simd_v128_load_nc(m, s0, 32)
	n35 := Simd_i16x8_mul(n32, n34)
	n36 := Simd_i32x4_extmul_low_i16x8_u(n32, n34)
	n37 := Simd_i32x4_extmul_high_i16x8_u(n32, n34)
	n38 := Simd_i8x16_shuffle(n36, n37, [2]uint64{p0, p0h})
	n39 := Simd_i8x16_shuffle(n35, n38, [2]uint64{p1, p1h})
	n40 := Simd_i32x4_add(n28, n39)
	n41 := Simd_i32x4_shr_s(n40, 17)
	n42 := Simd_i8x16_shuffle(n35, n38, [2]uint64{p2, p2h})
	n43 := Simd_v128_load_nc(m, s0, 80)
	n44 := Simd_i32x4_add(n43, n42)
	n45 := Simd_i32x4_shr_s(n44, 17)
	n46 := Simd_i16x8_narrow_i32x4_s(n41, n45)
	n47 := Simd_i16x8_min_s(n46, [2]uint64{p3, p3h})
	n48 := Simd_v128_xor(n47, n33)
	n49 := Simd_i16x8_sub(n48, n33)
	n50 := Simd_i8x16_shuffle(n49, [2]uint64{p4, p4h}, [2]uint64{506097522914230528, 940138560043747592})
	n51 := Simd_i8x16_shuffle(n50, [2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n52 := Simd_i8x16_shuffle(n51, [2]uint64{p4, p4h}, [2]uint64{506097522914230528, 1084816697971969292})
	n53 := Simd_v128_load_nc(m, s0, 0)
	n54 := Simd_i16x8_mul(n49, n53)
	_ = Simd_v128_store(m, s1, 0, n54)
	_ = Simd_v128_store(m, s2, 0, n52)
	_ = Simd_v128_store(m, s2, 16, n24)
	return n52[0], n52[1], n24[0], n24[1]
}

//go:noinline
func Simd_p_fx356(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_narrow_i16x8_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_eq(n0, [2]uint64{p2, p2h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx357(m *Module, s0 int32, s1 int32) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 96, 16, 208)
	n1 := Simd_v128_load_nc(m, s0, 208)
	n2 := Simd_v128_load(m, s1, 16)
	n3 := Simd_i16x8_abs(n2)
	n4 := Simd_i16x8_add(n1, n3)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx358(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_shr_s([2]uint64{p2, p2h}, 15)
	n1 := Simd_v128_load_nc(m, s0, 48)
	n2 := Simd_i16x8_mul([2]uint64{p0, p0h}, n1)
	n3 := Simd_i32x4_extmul_low_i16x8_u([2]uint64{p0, p0h}, n1)
	n4 := Simd_i32x4_extmul_high_i16x8_u([2]uint64{p0, p0h}, n1)
	n5 := Simd_i8x16_shuffle(n3, n4, [2]uint64{p1, p1h})
	n6 := Simd_i8x16_shuffle(n2, n5, [2]uint64{p4, p4h})
	n7 := Simd_i32x4_add([2]uint64{p3, p3h}, n6)
	n8 := Simd_i32x4_shr_s(n7, 17)
	n9 := Simd_i8x16_shuffle(n2, n5, [2]uint64{p5, p5h})
	n10 := Simd_v128_load_nc(m, s0, 112)
	n11 := Simd_i32x4_add(n10, n9)
	n12 := Simd_i32x4_shr_s(n11, 17)
	n13 := Simd_i16x8_narrow_i32x4_s(n8, n12)
	n14 := Simd_i16x8_min_s(n13, [2]uint64{p6, p6h})
	n15 := Simd_v128_xor(n14, n0)
	n16 := Simd_i16x8_sub(n15, n0)
	n17 := Simd_v128_load_nc(m, s0, 16)
	return n1[0], n1[1], n10[0], n10[1], n16[0], n16[1], n17[0], n17[1]
}

//go:noinline
func Simd_p_fx359(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	_ = Simd_v128_store(m, s0, 16, n0)
	n2 := Simd_v128_load_rng(m, s1, 64, 0, 208)
	n3 := Simd_v128_load_nc(m, s1, 192)
	n4 := Simd_v128_load(m, s0, 0)
	n5 := Simd_i16x8_abs(n4)
	n6 := Simd_i16x8_add(n3, n5)
	return n2[0], n2[1], n3[0], n3[1], n4[0], n4[1], n6[0], n6[1]
}

//go:noinline
func Simd_p_fx360(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_shr_s([2]uint64{p2, p2h}, 15)
	n1 := Simd_v128_load_nc(m, s0, 32)
	n2 := Simd_i16x8_mul([2]uint64{p0, p0h}, n1)
	n3 := Simd_i32x4_extmul_low_i16x8_u([2]uint64{p0, p0h}, n1)
	n4 := Simd_i32x4_extmul_high_i16x8_u([2]uint64{p0, p0h}, n1)
	n5 := Simd_i8x16_shuffle(n3, n4, [2]uint64{p1, p1h})
	n6 := Simd_i8x16_shuffle(n2, n5, [2]uint64{p4, p4h})
	n7 := Simd_i32x4_add([2]uint64{p3, p3h}, n6)
	n8 := Simd_i32x4_shr_s(n7, 17)
	n9 := Simd_i8x16_shuffle(n2, n5, [2]uint64{p5, p5h})
	n10 := Simd_v128_load_nc(m, s0, 80)
	n11 := Simd_i32x4_add(n10, n9)
	n12 := Simd_i32x4_shr_s(n11, 17)
	n13 := Simd_i16x8_narrow_i32x4_s(n8, n12)
	n14 := Simd_i16x8_min_s(n13, [2]uint64{p6, p6h})
	n15 := Simd_v128_xor(n14, n0)
	n16 := Simd_i16x8_sub(n15, n0)
	n17 := Simd_v128_load_nc(m, s0, 0)
	return n1[0], n1[1], n10[0], n10[1], n16[0], n16[1], n17[0], n17[1]
}

//go:noinline
func Simd_p_fx361(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p2, p2h}, [2]uint64{p5, p5h})
	_ = Simd_v128_store(m, s0, 0, n0)
	_ = Simd_v128_store(m, s1, 0, n3)
	return n3[0], n3[1]
}

//go:noinline
func Simd_p_fx362(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p1, p1h}, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p1, p1h}, [2]uint64{p4, p4h})
	_ = Simd_v128_store(m, s0, 16, n2)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx363(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 32)
	n1 := Simd_i16x8_abs(n0)
	n2 := Simd_i16x8_add([2]uint64{p0, p0h}, n1)
	n3 := Simd_i32x4_extmul_low_i16x8_u(n2, [2]uint64{p1, p1h})
	n4 := Simd_i16x8_shr_s(n0, 15)
	n5 := Simd_i16x8_mul([2]uint64{p1, p1h}, n2)
	n6 := Simd_i32x4_extmul_high_i16x8_u(n2, [2]uint64{p1, p1h})
	n7 := Simd_i8x16_shuffle(n3, n6, [2]uint64{p2, p2h})
	return n5[0], n5[1], n7[0], n7[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx364(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n1 := Simd_i32x4_add([2]uint64{p0, p0h}, n0)
	n2 := Simd_i32x4_shr_s(n1, 17)
	n3 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p5, p5h})
	n4 := Simd_i32x4_add([2]uint64{p4, p4h}, n3)
	n5 := Simd_i32x4_shr_s(n4, 17)
	n6 := Simd_i16x8_narrow_i32x4_s(n2, n5)
	n7 := Simd_i16x8_min_s(n6, [2]uint64{p6, p6h})
	return n7[0], n7[1]
}

//go:noinline
func Simd_p_fx365(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	_ = Simd_v128_store(m, s0, 32, n0)
	n2 := Simd_v128_load(m, s0, 48)
	n3 := Simd_i16x8_abs(n2)
	n4 := Simd_i16x8_add([2]uint64{p2, p2h}, n3)
	n5 := Simd_i32x4_extmul_low_i16x8_u(n4, [2]uint64{p3, p3h})
	n6 := Simd_i16x8_shr_s(n2, 15)
	n7 := Simd_i16x8_mul([2]uint64{p3, p3h}, n4)
	n8 := Simd_i32x4_extmul_high_i16x8_u(n4, [2]uint64{p3, p3h})
	n9 := Simd_i8x16_shuffle(n5, n8, [2]uint64{p4, p4h})
	return n7[0], n7[1], n9[0], n9[1], n6[0], n6[1]
}

//go:noinline
func Simd_p_fx366(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p2, p2h}, [2]uint64{p5, p5h})
	_ = Simd_v128_store(m, s0, 48, n0)
	_ = Simd_v128_store(m, s1, 48, n3)
	return n3[0], n3[1]
}

//go:noinline
func Simd_p_fx367(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p1, p1h}, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p1, p1h}, [2]uint64{p4, p4h})
	_ = Simd_v128_store(m, s0, 32, n2)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx368(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 96, 16, 112)
	n1 := Simd_v128_load(m, s1, 16)
	n2 := Simd_i16x8_abs(n1)
	n3 := Simd_i16x8_shr_s(n1, 15)
	n4 := Simd_v128_load_nc(m, s0, 48)
	n5 := Simd_i16x8_mul(n2, n4)
	n6 := Simd_i32x4_extmul_low_i16x8_u(n2, n4)
	n7 := Simd_i32x4_extmul_high_i16x8_u(n2, n4)
	n8 := Simd_i8x16_shuffle(n6, n7, [2]uint64{p0, p0h})
	n9 := Simd_i8x16_shuffle(n5, n8, [2]uint64{p1, p1h})
	n10 := Simd_i32x4_add(n0, n9)
	n11 := Simd_i32x4_shr_s(n10, 17)
	n12 := Simd_i8x16_shuffle(n5, n8, [2]uint64{p2, p2h})
	n13 := Simd_v128_load_nc(m, s0, 112)
	n14 := Simd_i32x4_add(n13, n12)
	n15 := Simd_i32x4_shr_s(n14, 17)
	n16 := Simd_i16x8_narrow_i32x4_s(n11, n15)
	n17 := Simd_i16x8_min_s(n16, [2]uint64{p3, p3h})
	n18 := Simd_v128_xor(n17, n3)
	n19 := Simd_i16x8_sub(n18, n3)
	n20 := Simd_i8x16_shuffle(n19, [2]uint64{p4, p4h}, [2]uint64{506093107721536258, 1084818905618843912})
	n21 := Simd_i8x16_shuffle(n20, [2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n22 := Simd_i8x16_shuffle(n21, [2]uint64{p4, p4h}, [2]uint64{216743454502551808, 1084818905618843912})
	n23 := Simd_v128_load_nc(m, s0, 16)
	n24 := Simd_i16x8_mul(n19, n23)
	_ = Simd_v128_store(m, s1, 16, n24)
	n26 := Simd_v128_load_rng(m, s0, 64, 0, 96)
	n27 := Simd_v128_load(m, s1, 0)
	n28 := Simd_i16x8_abs(n27)
	n29 := Simd_i16x8_shr_s(n27, 15)
	n30 := Simd_v128_load_nc(m, s0, 32)
	n31 := Simd_i16x8_mul(n28, n30)
	n32 := Simd_i32x4_extmul_low_i16x8_u(n28, n30)
	n33 := Simd_i32x4_extmul_high_i16x8_u(n28, n30)
	n34 := Simd_i8x16_shuffle(n32, n33, [2]uint64{p0, p0h})
	n35 := Simd_i8x16_shuffle(n31, n34, [2]uint64{p1, p1h})
	n36 := Simd_i32x4_add(n26, n35)
	n37 := Simd_i32x4_shr_s(n36, 17)
	n38 := Simd_i8x16_shuffle(n31, n34, [2]uint64{p2, p2h})
	n39 := Simd_v128_load_nc(m, s0, 80)
	n40 := Simd_i32x4_add(n39, n38)
	n41 := Simd_i32x4_shr_s(n40, 17)
	n42 := Simd_i16x8_narrow_i32x4_s(n37, n41)
	n43 := Simd_i16x8_min_s(n42, [2]uint64{p3, p3h})
	n44 := Simd_v128_xor(n43, n29)
	n45 := Simd_i16x8_sub(n44, n29)
	n46 := Simd_i8x16_shuffle(n45, [2]uint64{p4, p4h}, [2]uint64{506097522914230528, 940138560043747592})
	n47 := Simd_i8x16_shuffle(n46, [2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n48 := Simd_i8x16_shuffle(n47, [2]uint64{p4, p4h}, [2]uint64{506097522914230528, 1084816697971969292})
	n49 := Simd_v128_load_nc(m, s0, 0)
	n50 := Simd_i16x8_mul(n45, n49)
	_ = Simd_v128_store(m, s1, 0, n50)
	_ = Simd_v128_store(m, s2, 0, n48)
	_ = Simd_v128_store(m, s2, 16, n22)
	return n48[0], n48[1], n22[0], n22[1]
}

//go:noinline
func Simd_p_fx369(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0+16, 0, -16, 32)
	n1 := Simd_i8x16_shuffle(n0, n0, [2]uint64{p0, p0h})
	n2 := Simd_i32x4_extend_low_i16x8_s(n1)
	n3 := Simd_i32x4_mul(n2, [2]uint64{p2, p2h})
	n4 := Simd_i32x4_extend_high_i16x8_s(n1)
	n5 := Simd_i32x4_mul(n4, [2]uint64{p1, p1h})
	n6 := Simd_i8x16_shuffle(n3, n5, [2]uint64{p3, p3h})
	n7 := Simd_v128_load_nc(m, s0, 0)
	n8 := Simd_i8x16_shuffle(n7, n7, [2]uint64{p0, p0h})
	n9 := Simd_i32x4_extend_low_i16x8_s(n8)
	n10 := Simd_i32x4_mul(n9, [2]uint64{p1, p1h})
	n11 := Simd_i32x4_extend_high_i16x8_s(n8)
	n12 := Simd_i32x4_mul(n11, [2]uint64{p2, p2h})
	n13 := Simd_i8x16_shuffle(n10, n12, [2]uint64{p3, p3h})
	n14 := Simd_i16x8_add(n0, n7)
	n15 := Simd_i16x8_add(n14, n13)
	n16 := Simd_i16x8_add(n15, n6)
	n17 := Simd_i8x16_shuffle(n16, n16, [2]uint64{p0, p0h})
	n18 := Simd_i16x8_sub(n7, n0)
	n19 := Simd_i8x16_shuffle(n18, n18, [2]uint64{p0, p0h})
	n20 := Simd_i16x8_add(n13, n19)
	n21 := Simd_i16x8_sub(n20, n6)
	n22 := Simd_i8x16_shuffle(n17, n21, [2]uint64{p4, p4h})
	n23 := Simd_i8x16_shuffle(n14, n18, [2]uint64{p4, p4h})
	n24 := Simd_i16x8_sub(n23, n22)
	n25 := Simd_i8x16_shuffle(n24, [2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n26 := Simd_i16x8_add(n22, n23)
	return n26[0], n26[1], n25[0], n25[1]
}

//go:noinline
func Simd_p_fx370(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p2, p2h})
	n3 := Simd_i16x8_add(n2, [2]uint64{p4, p4h})
	n4 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p3, p3h})
	n5 := Simd_i16x8_add(n3, n4)
	return n2[0], n2[1], n3[0], n3[1], n4[0], n4[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx371(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_extend_low_i16x8_s(n0)
	n2 := Simd_i32x4_mul(n1, [2]uint64{p2, p2h})
	n3 := Simd_i32x4_extend_high_i16x8_s(n0)
	n4 := Simd_i32x4_mul(n3, [2]uint64{p3, p3h})
	n5 := Simd_i8x16_shuffle(n2, n4, [2]uint64{p4, p4h})
	n6 := Simd_i16x8_add([2]uint64{p6, p6h}, n5)
	n7 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, [2]uint64{p5, p5h}, [2]uint64{p1, p1h})
	n8 := Simd_i32x4_extend_low_i16x8_s(n7)
	n9 := Simd_i32x4_mul(n8, [2]uint64{p3, p3h})
	n10 := Simd_i32x4_extend_high_i16x8_s(n7)
	n11 := Simd_i32x4_mul(n10, [2]uint64{p2, p2h})
	n12 := Simd_i8x16_shuffle(n9, n11, [2]uint64{p4, p4h})
	n13 := Simd_i16x8_add(n6, n12)
	return n5[0], n5[1], n12[0], n12[1], n13[0], n13[1]
}

//go:noinline
func Simd_p_fx372(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_shuffle(n0, n0, [2]uint64{p3, p3h})
	n2 := Simd_i16x8_sub(n1, [2]uint64{p4, p4h})
	n3 := Simd_i16x8_add(n2, [2]uint64{p5, p5h})
	n4 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n5 := Simd_i8x16_shuffle(n4, n3, [2]uint64{p6, p6h})
	return n0[0], n0[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx373(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i16x8_sub(n0, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n3 := Simd_i16x8_shr_s(n2, 3)
	n4 := Simd_i16x8_add([2]uint64{p3, p3h}, n0)
	n5 := Simd_i16x8_shr_s(n4, 3)
	n6 := Simd_i8x16_shuffle(n5, n3, [2]uint64{p6, p6h})
	return n5[0], n5[1], n3[0], n3[1], n6[0], n6[1]
}

//go:noinline
func Simd_p_fx374(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p6, p6h})
	n2 := Simd_v128_load32_zero(m, s0, 0)
	n3 := Simd_v128_load32_zero(m, s0, 32)
	n4 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p3, p3h})
	n5 := Simd_i8x16_shuffle(n4, [2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n6 := Simd_i16x8_add(n0, n5)
	n7 := Simd_v128_load32_zero(m, s0, 64)
	n8 := Simd_v128_load32_zero(m, s0, 96)
	n9 := Simd_i8x16_shuffle(n7, n8, [2]uint64{p3, p3h})
	n10 := Simd_i8x16_shuffle(n9, [2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n11 := Simd_i16x8_add(n1, n10)
	n12 := Simd_i8x16_narrow_i16x8_u(n6, n11)
	return n12[0], n12[1]
}

//go:noinline
func Simd_p_fx375(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 16, 0, 64)
	n1 := Simd_v128_load_nc(m, s0, 48)
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p0, p0h})
	n3 := Simd_v128_load_nc(m, s0, 0)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx376(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_nc(m, s0, 32)
	n1 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n0, [2]uint64{p1, p1h})
	n2 := Simd_i32x4_extend_low_i16x8_s(n1)
	n3 := Simd_i32x4_extend_high_i16x8_s(n1)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx377(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_extend_low_i16x8_s([2]uint64{p0, p0h})
	n1 := Simd_i32x4_mul(n0, [2]uint64{p6, p6h})
	n2 := Simd_i32x4_extend_high_i16x8_s([2]uint64{p0, p0h})
	n3 := Simd_i32x4_mul(n2, [2]uint64{p6, p6h})
	n4 := Simd_i8x16_shuffle(n1, n3, [2]uint64{p5, p5h})
	n5 := Simd_i16x8_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n6 := Simd_i32x4_mul([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n7 := Simd_i32x4_mul([2]uint64{p4, p4h}, [2]uint64{p3, p3h})
	n8 := Simd_i8x16_shuffle(n6, n7, [2]uint64{p5, p5h})
	n9 := Simd_i16x8_add(n5, n8)
	n10 := Simd_i16x8_add(n9, n4)
	return n0[0], n0[1], n2[0], n2[1], n10[0], n10[1]
}

//go:noinline
func Simd_p_fx378(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p4, p4h}, [2]uint64{p2, p2h})
	n2 := Simd_i16x8_add(n0, n1)
	n3 := Simd_i16x8_add([2]uint64{p5, p5h}, n2)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx379(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n2 := Simd_i16x8_add(n0, n1)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx380(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_shr_s(n0, 4)
	n2 := Simd_i16x8_narrow_i32x4_s([2]uint64{p2, p2h}, [2]uint64{p2, p2h})
	n3 := Simd_i8x16_shuffle(n1, n2, [2]uint64{p3, p3h})
	n4 := Simd_i32x4_dot_i16x8_s([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n5 := Simd_i32x4_add(n4, [2]uint64{p6, p6h})
	n6 := Simd_i32x4_shr_s(n5, 16)
	_ = Simd_v128_store(m, s0, 16, n3)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx381(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_i16x8_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_shr_s(n0, 4)
	n2 := Simd_i16x8_narrow_i32x4_s([2]uint64{p2, p2h}, [2]uint64{p2, p2h})
	n3 := Simd_i16x8_eq([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n4 := Simd_i16x8_add(n2, n3)
	n5 := Simd_i8x16_shuffle(n1, n4, [2]uint64{p5, p5h})
	_ = Simd_v128_store(m, s0, 0, n5)
	return
}

//go:noinline
func Simd_p_fx382(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n2 := Simd_i16x8_sub(n0, n1)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx383(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n2 := Simd_i16x8_sub(n0, n1)
	n3 := Simd_i8x16_shuffle([2]uint64{p4, p4h}, n2, [2]uint64{p5, p5h})
	n4 := Simd_i8x16_shuffle(n3, [2]uint64{p1, p1h}, [2]uint64{p6, p6h})
	return n2[0], n2[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx384(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_dot_i16x8_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_add(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_shr_s(n1, 9)
	n3 := Simd_i32x4_dot_i16x8_s([2]uint64{p0, p0h}, [2]uint64{p3, p3h})
	n4 := Simd_i32x4_add(n3, [2]uint64{p4, p4h})
	n5 := Simd_i32x4_shr_s(n4, 9)
	n6 := Simd_i16x8_narrow_i32x4_s(n2, n5)
	n7 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, n6, [2]uint64{p6, p6h})
	return n6[0], n6[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx385(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, n0, [2]uint64{p4, p4h})
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	return n0[0], n0[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx386(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i16x8_sub(n0, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle(n1, n1, [2]uint64{p5, p5h})
	n3 := Simd_i8x16_shuffle(n2, n1, [2]uint64{p6, p6h})
	n4 := Simd_i16x8_add([2]uint64{p3, p3h}, n0)
	n5 := Simd_i16x8_add(n4, [2]uint64{p4, p4h})
	n6 := Simd_i8x16_shuffle(n4, n4, [2]uint64{p5, p5h})
	return n0[0], n0[1], n5[0], n5[1], n6[0], n6[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx387(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_i32x4_dot_i16x8_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_add(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_shr_s(n1, 16)
	n3 := Simd_i16x8_narrow_i32x4_s(n2, n2)
	n4 := Simd_i16x8_sub([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n5 := Simd_i16x8_shr_s(n4, 4)
	n6 := Simd_i8x16_shuffle(n5, n3, [2]uint64{p5, p5h})
	_ = Simd_v128_store(m, s0, 48, n6)
	return
}

//go:noinline
func Simd_p_fx388(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n2 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, n1, [2]uint64{p6, p6h})
	return n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx389(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i16x8_add([2]uint64{p3, p3h}, n0)
	n2 := Simd_i32x4_dot_i16x8_s(n1, [2]uint64{p4, p4h})
	n3 := Simd_i32x4_dot_i16x8_s(n1, [2]uint64{p5, p5h})
	n4 := Simd_i16x8_narrow_i32x4_s(n2, n3)
	n5 := Simd_i16x8_sub([2]uint64{p3, p3h}, n0)
	return n4[0], n4[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx390(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_i32x4_dot_i16x8_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_add(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_shr_s(n1, 16)
	n3 := Simd_i16x8_narrow_i32x4_s(n2, n2)
	n4 := Simd_i16x8_sub([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n5 := Simd_i16x8_shr_s(n4, 4)
	n6 := Simd_i8x16_shuffle(n5, n3, [2]uint64{p5, p5h})
	_ = Simd_v128_store(m, s0, 16, n6)
	return
}

//go:noinline
func Simd_p_fx391(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_dot_i16x8_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_add(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_shr_s(n1, 16)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx392(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_i16x8_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_shr_s(n0, 4)
	n2 := Simd_i16x8_narrow_i32x4_s([2]uint64{p2, p2h}, [2]uint64{p2, p2h})
	n3 := Simd_i16x8_eq([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n4 := Simd_i16x8_add(n2, n3)
	n5 := Simd_i8x16_shuffle(n1, n4, [2]uint64{p5, p5h})
	_ = Simd_v128_store(m, s0, 32, n5)
	return
}

//go:noinline
func Simd_p_fx393(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i16x8_add_sat_s([2]uint64{p3, p3h}, n0)
	n2 := Simd_i16x8_sub_sat_s([2]uint64{p3, p3h}, n0)
	n3 := Simd_i8x16_shuffle(n2, n1, [2]uint64{p4, p4h})
	n4 := Simd_i8x16_shuffle(n1, n2, [2]uint64{p4, p4h})
	n5 := Simd_i8x16_shuffle(n4, n3, [2]uint64{p5, p5h})
	n6 := Simd_i32x4_dot_i16x8_s(n5, [2]uint64{p6, p6h})
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx394(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_sub([2]uint64{p1, p1h}, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_narrow_i32x4_s(n0, n1)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx395(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) {
	n0 := Simd_i32x4_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_sub([2]uint64{p1, p1h}, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_narrow_i32x4_s(n0, n1)
	n3 := Simd_i16x8_sub(n2, [2]uint64{p2, p2h})
	n4 := Simd_i8x16_shuffle(n3, n3, [2]uint64{1084818905618843912, 2242261671028070680})
	n5 := Simd_i8x16_shuffle(n4, n3, [2]uint64{p3, p3h})
	n6 := Simd_i16x8_shr_s(n5, 1)
	n7 := Simd_i16x8_add([2]uint64{p2, p2h}, n2)
	n8 := Simd_i16x8_shr_s(n7, 1)
	_ = Simd_v128_store(m, s0, 0, n8)
	_ = Simd_v128_store(m, s0, 16, n6)
	return
}

//go:noinline
func Simd_p_fx396(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 448, 0, 496)
	n1 := Simd_v128_load_rng(m, s1, 448, 0, 496)
	n2 := Simd_i8x16_sub_sat_u(n1, n0)
	n3 := Simd_i8x16_sub_sat_u(n0, n1)
	n4 := Simd_v128_or(n3, n2)
	n5 := Simd_i8x16_shuffle(n4, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n6 := Simd_i32x4_dot_i16x8_s(n5, n5)
	n7 := Simd_i8x16_shuffle(n4, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n8 := Simd_i32x4_dot_i16x8_s(n7, n7)
	n9 := Simd_v128_load_nc(m, s0, 384)
	n10 := Simd_v128_load_nc(m, s1, 384)
	n11 := Simd_i8x16_sub_sat_u(n10, n9)
	n12 := Simd_i8x16_sub_sat_u(n9, n10)
	n13 := Simd_v128_or(n12, n11)
	n14 := Simd_i8x16_shuffle(n13, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n15 := Simd_i32x4_dot_i16x8_s(n14, n14)
	n16 := Simd_i8x16_shuffle(n13, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n17 := Simd_i32x4_dot_i16x8_s(n16, n16)
	n18 := Simd_i32x4_add(n6, n15)
	n19 := Simd_v128_load_nc(m, s0, 320)
	n20 := Simd_v128_load_nc(m, s1, 320)
	n21 := Simd_i8x16_sub_sat_u(n20, n19)
	n22 := Simd_i8x16_sub_sat_u(n19, n20)
	n23 := Simd_v128_or(n22, n21)
	n24 := Simd_i8x16_shuffle(n23, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n25 := Simd_i32x4_dot_i16x8_s(n24, n24)
	n26 := Simd_i8x16_shuffle(n23, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n27 := Simd_i32x4_dot_i16x8_s(n26, n26)
	n28 := Simd_i32x4_add(n18, n25)
	n29 := Simd_v128_load_nc(m, s0, 256)
	n30 := Simd_v128_load_nc(m, s1, 256)
	n31 := Simd_i8x16_sub_sat_u(n30, n29)
	n32 := Simd_i8x16_sub_sat_u(n29, n30)
	n33 := Simd_v128_or(n32, n31)
	n34 := Simd_i8x16_shuffle(n33, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n35 := Simd_i32x4_dot_i16x8_s(n34, n34)
	n36 := Simd_i8x16_shuffle(n33, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n37 := Simd_i32x4_dot_i16x8_s(n36, n36)
	n38 := Simd_i32x4_add(n28, n35)
	n39 := Simd_v128_load_nc(m, s0, 192)
	n40 := Simd_v128_load_nc(m, s1, 192)
	n41 := Simd_i8x16_sub_sat_u(n40, n39)
	n42 := Simd_i8x16_sub_sat_u(n39, n40)
	n43 := Simd_v128_or(n42, n41)
	n44 := Simd_i8x16_shuffle(n43, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n45 := Simd_i32x4_dot_i16x8_s(n44, n44)
	n46 := Simd_i8x16_shuffle(n43, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n47 := Simd_i32x4_dot_i16x8_s(n46, n46)
	n48 := Simd_i32x4_add(n38, n45)
	n49 := Simd_v128_load_nc(m, s0, 128)
	n50 := Simd_v128_load_nc(m, s1, 128)
	n51 := Simd_i8x16_sub_sat_u(n50, n49)
	n52 := Simd_i8x16_sub_sat_u(n49, n50)
	n53 := Simd_v128_or(n52, n51)
	n54 := Simd_i8x16_shuffle(n53, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n55 := Simd_i32x4_dot_i16x8_s(n54, n54)
	n56 := Simd_i8x16_shuffle(n53, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n57 := Simd_i32x4_dot_i16x8_s(n56, n56)
	n58 := Simd_i32x4_add(n48, n55)
	n59 := Simd_v128_load_nc(m, s0, 64)
	n60 := Simd_v128_load_nc(m, s1, 64)
	n61 := Simd_i8x16_sub_sat_u(n60, n59)
	n62 := Simd_i8x16_sub_sat_u(n59, n60)
	n63 := Simd_v128_or(n62, n61)
	n64 := Simd_i8x16_shuffle(n63, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n65 := Simd_i32x4_dot_i16x8_s(n64, n64)
	n66 := Simd_i8x16_shuffle(n63, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n67 := Simd_i32x4_dot_i16x8_s(n66, n66)
	n68 := Simd_i32x4_add(n58, n65)
	n69 := Simd_v128_load_nc(m, s0, 0)
	n70 := Simd_v128_load_nc(m, s1, 0)
	n71 := Simd_i8x16_sub_sat_u(n70, n69)
	n72 := Simd_i8x16_sub_sat_u(n69, n70)
	n73 := Simd_v128_or(n72, n71)
	n74 := Simd_i8x16_shuffle(n73, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n75 := Simd_i32x4_dot_i16x8_s(n74, n74)
	n76 := Simd_i8x16_shuffle(n73, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n77 := Simd_i32x4_dot_i16x8_s(n76, n76)
	n78 := Simd_i32x4_add(n68, n75)
	n79 := Simd_i32x4_add(n78, n77)
	n80 := Simd_v128_load_nc(m, s0, 32)
	n81 := Simd_v128_load_nc(m, s1, 32)
	n82 := Simd_i8x16_sub_sat_u(n81, n80)
	n83 := Simd_i8x16_sub_sat_u(n80, n81)
	n84 := Simd_v128_or(n83, n82)
	n85 := Simd_i8x16_shuffle(n84, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n86 := Simd_i32x4_dot_i16x8_s(n85, n85)
	n87 := Simd_i8x16_shuffle(n84, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n88 := Simd_i32x4_dot_i16x8_s(n87, n87)
	n89 := Simd_i32x4_add(n79, n86)
	n90 := Simd_i32x4_add(n89, n88)
	n91 := Simd_i32x4_add(n90, n67)
	n92 := Simd_v128_load_nc(m, s0, 96)
	n93 := Simd_v128_load_nc(m, s1, 96)
	n94 := Simd_i8x16_sub_sat_u(n93, n92)
	n95 := Simd_i8x16_sub_sat_u(n92, n93)
	n96 := Simd_v128_or(n95, n94)
	n97 := Simd_i8x16_shuffle(n96, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n98 := Simd_i32x4_dot_i16x8_s(n97, n97)
	n99 := Simd_i8x16_shuffle(n96, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n100 := Simd_i32x4_dot_i16x8_s(n99, n99)
	n101 := Simd_i32x4_add(n91, n98)
	n102 := Simd_i32x4_add(n101, n100)
	n103 := Simd_i32x4_add(n102, n57)
	n104 := Simd_v128_load_nc(m, s0, 160)
	n105 := Simd_v128_load_nc(m, s1, 160)
	n106 := Simd_i8x16_sub_sat_u(n105, n104)
	n107 := Simd_i8x16_sub_sat_u(n104, n105)
	n108 := Simd_v128_or(n107, n106)
	n109 := Simd_i8x16_shuffle(n108, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n110 := Simd_i32x4_dot_i16x8_s(n109, n109)
	n111 := Simd_i8x16_shuffle(n108, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n112 := Simd_i32x4_dot_i16x8_s(n111, n111)
	n113 := Simd_i32x4_add(n103, n110)
	n114 := Simd_i32x4_add(n113, n112)
	n115 := Simd_i32x4_add(n114, n47)
	n116 := Simd_v128_load_nc(m, s0, 224)
	n117 := Simd_v128_load_nc(m, s1, 224)
	n118 := Simd_i8x16_sub_sat_u(n117, n116)
	n119 := Simd_i8x16_sub_sat_u(n116, n117)
	n120 := Simd_v128_or(n119, n118)
	n121 := Simd_i8x16_shuffle(n120, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n122 := Simd_i32x4_dot_i16x8_s(n121, n121)
	n123 := Simd_i8x16_shuffle(n120, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n124 := Simd_i32x4_dot_i16x8_s(n123, n123)
	n125 := Simd_i32x4_add(n115, n122)
	n126 := Simd_i32x4_add(n125, n124)
	n127 := Simd_i32x4_add(n126, n37)
	n128 := Simd_v128_load_nc(m, s0, 288)
	n129 := Simd_v128_load_nc(m, s1, 288)
	n130 := Simd_i8x16_sub_sat_u(n129, n128)
	n131 := Simd_i8x16_sub_sat_u(n128, n129)
	n132 := Simd_v128_or(n131, n130)
	n133 := Simd_i8x16_shuffle(n132, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n134 := Simd_i32x4_dot_i16x8_s(n133, n133)
	n135 := Simd_i8x16_shuffle(n132, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n136 := Simd_i32x4_dot_i16x8_s(n135, n135)
	n137 := Simd_i32x4_add(n127, n134)
	n138 := Simd_i32x4_add(n137, n136)
	n139 := Simd_i32x4_add(n138, n27)
	n140 := Simd_v128_load_nc(m, s0, 352)
	n141 := Simd_v128_load_nc(m, s1, 352)
	n142 := Simd_i8x16_sub_sat_u(n141, n140)
	n143 := Simd_i8x16_sub_sat_u(n140, n141)
	n144 := Simd_v128_or(n143, n142)
	n145 := Simd_i8x16_shuffle(n144, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n146 := Simd_i32x4_dot_i16x8_s(n145, n145)
	n147 := Simd_i8x16_shuffle(n144, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n148 := Simd_i32x4_dot_i16x8_s(n147, n147)
	n149 := Simd_i32x4_add(n139, n146)
	n150 := Simd_i32x4_add(n149, n148)
	n151 := Simd_i32x4_add(n150, n17)
	n152 := Simd_v128_load_nc(m, s0, 416)
	n153 := Simd_v128_load_nc(m, s1, 416)
	n154 := Simd_i8x16_sub_sat_u(n153, n152)
	n155 := Simd_i8x16_sub_sat_u(n152, n153)
	n156 := Simd_v128_or(n155, n154)
	n157 := Simd_i8x16_shuffle(n156, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n158 := Simd_i32x4_dot_i16x8_s(n157, n157)
	n159 := Simd_i8x16_shuffle(n156, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n160 := Simd_i32x4_dot_i16x8_s(n159, n159)
	n161 := Simd_i32x4_add(n151, n158)
	n162 := Simd_i32x4_add(n161, n160)
	n163 := Simd_i32x4_add(n162, n8)
	n164 := Simd_v128_load_nc(m, s0, 480)
	n165 := Simd_v128_load_nc(m, s1, 480)
	n166 := Simd_i8x16_sub_sat_u(n165, n164)
	n167 := Simd_i8x16_sub_sat_u(n164, n165)
	n168 := Simd_v128_or(n167, n166)
	n169 := Simd_i8x16_shuffle(n168, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n170 := Simd_i32x4_dot_i16x8_s(n169, n169)
	n171 := Simd_i8x16_shuffle(n168, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n172 := Simd_i32x4_dot_i16x8_s(n171, n171)
	n173 := Simd_i32x4_add(n163, n170)
	n174 := Simd_i32x4_add(n173, n172)
	return n174[0], n174[1]
}

//go:noinline
func Simd_p_fx397(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 192, 0, 240)
	n1 := Simd_v128_load_rng(m, s1, 192, 0, 240)
	n2 := Simd_i8x16_sub_sat_u(n1, n0)
	n3 := Simd_i8x16_sub_sat_u(n0, n1)
	n4 := Simd_v128_or(n3, n2)
	n5 := Simd_i8x16_shuffle(n4, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n6 := Simd_i32x4_dot_i16x8_s(n5, n5)
	n7 := Simd_i8x16_shuffle(n4, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n8 := Simd_i32x4_dot_i16x8_s(n7, n7)
	n9 := Simd_v128_load_nc(m, s0, 128)
	n10 := Simd_v128_load_nc(m, s1, 128)
	n11 := Simd_i8x16_sub_sat_u(n10, n9)
	n12 := Simd_i8x16_sub_sat_u(n9, n10)
	n13 := Simd_v128_or(n12, n11)
	n14 := Simd_i8x16_shuffle(n13, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n15 := Simd_i32x4_dot_i16x8_s(n14, n14)
	n16 := Simd_i8x16_shuffle(n13, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n17 := Simd_i32x4_dot_i16x8_s(n16, n16)
	n18 := Simd_i32x4_add(n6, n15)
	n19 := Simd_v128_load_nc(m, s0, 64)
	n20 := Simd_v128_load_nc(m, s1, 64)
	n21 := Simd_i8x16_sub_sat_u(n20, n19)
	n22 := Simd_i8x16_sub_sat_u(n19, n20)
	n23 := Simd_v128_or(n22, n21)
	n24 := Simd_i8x16_shuffle(n23, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n25 := Simd_i32x4_dot_i16x8_s(n24, n24)
	n26 := Simd_i8x16_shuffle(n23, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n27 := Simd_i32x4_dot_i16x8_s(n26, n26)
	n28 := Simd_i32x4_add(n18, n25)
	n29 := Simd_v128_load_nc(m, s0, 0)
	n30 := Simd_v128_load_nc(m, s1, 0)
	n31 := Simd_i8x16_sub_sat_u(n30, n29)
	n32 := Simd_i8x16_sub_sat_u(n29, n30)
	n33 := Simd_v128_or(n32, n31)
	n34 := Simd_i8x16_shuffle(n33, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n35 := Simd_i32x4_dot_i16x8_s(n34, n34)
	n36 := Simd_i8x16_shuffle(n33, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n37 := Simd_i32x4_dot_i16x8_s(n36, n36)
	n38 := Simd_i32x4_add(n28, n35)
	n39 := Simd_i32x4_add(n38, n37)
	n40 := Simd_v128_load_nc(m, s0, 32)
	n41 := Simd_v128_load_nc(m, s1, 32)
	n42 := Simd_i8x16_sub_sat_u(n41, n40)
	n43 := Simd_i8x16_sub_sat_u(n40, n41)
	n44 := Simd_v128_or(n43, n42)
	n45 := Simd_i8x16_shuffle(n44, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n46 := Simd_i32x4_dot_i16x8_s(n45, n45)
	n47 := Simd_i8x16_shuffle(n44, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n48 := Simd_i32x4_dot_i16x8_s(n47, n47)
	n49 := Simd_i32x4_add(n39, n46)
	n50 := Simd_i32x4_add(n49, n48)
	n51 := Simd_i32x4_add(n50, n27)
	n52 := Simd_v128_load_nc(m, s0, 96)
	n53 := Simd_v128_load_nc(m, s1, 96)
	n54 := Simd_i8x16_sub_sat_u(n53, n52)
	n55 := Simd_i8x16_sub_sat_u(n52, n53)
	n56 := Simd_v128_or(n55, n54)
	n57 := Simd_i8x16_shuffle(n56, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n58 := Simd_i32x4_dot_i16x8_s(n57, n57)
	n59 := Simd_i8x16_shuffle(n56, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n60 := Simd_i32x4_dot_i16x8_s(n59, n59)
	n61 := Simd_i32x4_add(n51, n58)
	n62 := Simd_i32x4_add(n61, n60)
	n63 := Simd_i32x4_add(n62, n17)
	n64 := Simd_v128_load_nc(m, s0, 160)
	n65 := Simd_v128_load_nc(m, s1, 160)
	n66 := Simd_i8x16_sub_sat_u(n65, n64)
	n67 := Simd_i8x16_sub_sat_u(n64, n65)
	n68 := Simd_v128_or(n67, n66)
	n69 := Simd_i8x16_shuffle(n68, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n70 := Simd_i32x4_dot_i16x8_s(n69, n69)
	n71 := Simd_i8x16_shuffle(n68, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n72 := Simd_i32x4_dot_i16x8_s(n71, n71)
	n73 := Simd_i32x4_add(n63, n70)
	n74 := Simd_i32x4_add(n73, n72)
	n75 := Simd_i32x4_add(n74, n8)
	n76 := Simd_v128_load_nc(m, s0, 224)
	n77 := Simd_v128_load_nc(m, s1, 224)
	n78 := Simd_i8x16_sub_sat_u(n77, n76)
	n79 := Simd_i8x16_sub_sat_u(n76, n77)
	n80 := Simd_v128_or(n79, n78)
	n81 := Simd_i8x16_shuffle(n80, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n82 := Simd_i32x4_dot_i16x8_s(n81, n81)
	n83 := Simd_i8x16_shuffle(n80, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n84 := Simd_i32x4_dot_i16x8_s(n83, n83)
	n85 := Simd_i32x4_add(n75, n82)
	n86 := Simd_i32x4_add(n85, n84)
	return n86[0], n86[1]
}

//go:noinline
func Simd_p_fx398(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n2 := Simd_i16x8_sub_sat_s(n0, n1)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx399(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_dot_i16x8_s([2]uint64{p0, p0h}, [2]uint64{p0, p0h})
	n1 := Simd_i32x4_dot_i16x8_s([2]uint64{p1, p1h}, [2]uint64{p1, p1h})
	n2 := Simd_i32x4_dot_i16x8_s([2]uint64{p2, p2h}, [2]uint64{p2, p2h})
	n3 := Simd_i32x4_dot_i16x8_s([2]uint64{p3, p3h}, [2]uint64{p3, p3h})
	n4 := Simd_i32x4_dot_i16x8_s([2]uint64{p4, p4h}, [2]uint64{p4, p4h})
	n5 := Simd_i32x4_add(n3, n4)
	n6 := Simd_i32x4_add(n2, n5)
	n7 := Simd_i32x4_dot_i16x8_s([2]uint64{p5, p5h}, [2]uint64{p5, p5h})
	n8 := Simd_i32x4_add(n6, n7)
	n9 := Simd_i32x4_add(n1, n8)
	n10 := Simd_i32x4_dot_i16x8_s([2]uint64{p6, p6h}, [2]uint64{p6, p6h})
	n11 := Simd_i32x4_add(n9, n10)
	n12 := Simd_i32x4_add(n0, n11)
	return n12[0], n12[1]
}

//go:noinline
func Simd_p_fx400(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n2 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, [2]uint64{p6, p6h}, [2]uint64{p2, p2h})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n4 := Simd_i16x8_sub_sat_s(n1, n3)
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx401(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_dot_i16x8_s([2]uint64{p0, p0h}, [2]uint64{p0, p0h})
	n1 := Simd_i32x4_dot_i16x8_s([2]uint64{p1, p1h}, [2]uint64{p1, p1h})
	n2 := Simd_i32x4_add(n0, n1)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx402(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n2 := Simd_i16x8_add([2]uint64{p5, p5h}, n1)
	return n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx403(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n2 := Simd_i16x8_add([2]uint64{p5, p5h}, n1)
	n3 := Simd_i16x8_add([2]uint64{p6, p6h}, n2)
	return n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx404(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_sub([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i16x8_sub(n1, n0)
	n3 := Simd_i16x8_add(n0, n1)
	n4 := Simd_i8x16_shuffle([2]uint64{p4, p4h}, n3, [2]uint64{p5, p5h})
	return n3[0], n3[1], n4[0], n4[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx405(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, n0, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_shuffle([2]uint64{p4, p4h}, n1, [2]uint64{p5, p5h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx406(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p4, p4h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{p5, p5h})
	n3 := Simd_i16x8_add([2]uint64{p6, p6h}, n2)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx407(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_sub([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i16x8_add(n0, n1)
	n3 := Simd_i16x8_sub(n0, n1)
	n4 := Simd_i16x8_sub([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n5 := Simd_v128_load_rng(m, s0, 0, 0, 32)
	return n2[0], n2[1], n5[0], n5[1], n3[0], n3[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx408(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i16x8_abs(n0)
	n2 := Simd_i32x4_dot_i16x8_s(n1, [2]uint64{p3, p3h})
	n3 := Simd_i8x16_shuffle([2]uint64{p4, p4h}, [2]uint64{p5, p5h}, [2]uint64{p2, p2h})
	n4 := Simd_i16x8_abs(n3)
	n5 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p6, p6h})
	n6 := Simd_i16x8_abs(n5)
	n7 := Simd_i32x4_dot_i16x8_s(n6, [2]uint64{p3, p3h})
	n8 := Simd_i8x16_shuffle([2]uint64{p4, p4h}, [2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n9 := Simd_i16x8_abs(n8)
	n10 := Simd_v128_load_nc(m, s0+16, 0)
	n11 := Simd_i32x4_dot_i16x8_s(n4, n10)
	n12 := Simd_i32x4_add(n2, n11)
	n13 := Simd_i32x4_dot_i16x8_s(n9, n10)
	n14 := Simd_i32x4_add(n7, n13)
	n15 := Simd_i32x4_sub(n12, n14)
	return n15[0], n15[1]
}

//go:noinline
func Simd_p_fx409(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_sub([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i16x8_add(n0, n1)
	n3 := Simd_i16x8_sub(n0, n1)
	n4 := Simd_i16x8_sub([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	return n2[0], n2[1], n3[0], n3[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx410(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i16x8_abs(n0)
	n2 := Simd_i32x4_dot_i16x8_s(n1, [2]uint64{p3, p3h})
	n3 := Simd_i8x16_shuffle([2]uint64{p4, p4h}, [2]uint64{p5, p5h}, [2]uint64{p2, p2h})
	n4 := Simd_i16x8_abs(n3)
	n5 := Simd_i32x4_dot_i16x8_s(n4, [2]uint64{p6, p6h})
	n6 := Simd_i32x4_add(n2, n5)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx411(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 32, 0, 112)
	n1 := Simd_v128_and(n0, [2]uint64{p0, p0h})
	n2 := Simd_i16x8_shr_u(n0, 8)
	n3 := Simd_v128_load_nc(m, s0, 0)
	n4 := Simd_v128_and(n3, [2]uint64{p0, p0h})
	n5 := Simd_i32x4_add(n1, n4)
	n6 := Simd_i16x8_shr_u(n3, 8)
	n7 := Simd_i32x4_add(n5, n6)
	n8 := Simd_v128_load_nc(m, s0, 64)
	n9 := Simd_v128_and(n8, [2]uint64{p0, p0h})
	n10 := Simd_i32x4_add(n7, n9)
	n11 := Simd_i32x4_add(n10, n2)
	n12 := Simd_i16x8_shr_u(n8, 8)
	n13 := Simd_v128_load_nc(m, s0, 96)
	n14 := Simd_v128_and(n13, [2]uint64{p0, p0h})
	n15 := Simd_i32x4_add(n11, n14)
	n16 := Simd_i32x4_add(n15, n12)
	n17 := Simd_i16x8_shr_u(n13, 8)
	n18 := Simd_i32x4_add(n16, n17)
	return n18[0], n18[1]
}

//go:noinline
func Simd_p_fx412(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 96, 16, 208)
	n1 := Simd_v128_load_nc(m, s0, 208)
	n2 := Simd_v128_load(m, s1, 16)
	n3 := Simd_i16x8_abs(n2)
	n4 := Simd_i16x8_add(n1, n3)
	n5 := Simd_i16x8_shr_s(n2, 15)
	n6 := Simd_i16x8_eq(n2, [2]uint64{p1, p1h})
	n7 := Simd_v128_load_nc(m, s0, 48)
	n8 := Simd_i16x8_mul(n4, n7)
	n9 := Simd_i32x4_extmul_low_i16x8_u(n4, n7)
	n10 := Simd_i32x4_extmul_high_i16x8_u(n4, n7)
	n11 := Simd_i8x16_shuffle(n9, n10, [2]uint64{p0, p0h})
	n12 := Simd_i8x16_shuffle(n8, n11, [2]uint64{p2, p2h})
	n13 := Simd_i32x4_add(n0, n12)
	n14 := Simd_i32x4_shr_s(n13, 17)
	n15 := Simd_i8x16_shuffle(n8, n11, [2]uint64{p3, p3h})
	n16 := Simd_v128_load_nc(m, s0, 112)
	n17 := Simd_i32x4_add(n16, n15)
	n18 := Simd_i32x4_shr_s(n17, 17)
	n19 := Simd_i16x8_narrow_i32x4_s(n14, n18)
	n20 := Simd_i16x8_min_s(n19, [2]uint64{p4, p4h})
	n21 := Simd_v128_bitselect([2]uint64{p1, p1h}, n20, n6)
	n22 := Simd_i16x8_add(n21, n5)
	n23 := Simd_v128_xor(n22, n5)
	n24 := Simd_v128_load_nc(m, s0, 16)
	n25 := Simd_i16x8_mul(n23, n24)
	_ = Simd_v128_store(m, s1, 16, n25)
	n27 := Simd_v128_load_rng(m, s0, 64, 0, 208)
	n28 := Simd_v128_load_nc(m, s0, 192)
	n29 := Simd_v128_load(m, s1, 0)
	n30 := Simd_i16x8_abs(n29)
	n31 := Simd_i16x8_add(n28, n30)
	n32 := Simd_i16x8_shr_s(n29, 15)
	n33 := Simd_i16x8_eq(n29, [2]uint64{p1, p1h})
	n34 := Simd_v128_load_nc(m, s0, 32)
	n35 := Simd_i16x8_mul(n31, n34)
	n36 := Simd_i32x4_extmul_low_i16x8_u(n31, n34)
	n37 := Simd_i32x4_extmul_high_i16x8_u(n31, n34)
	n38 := Simd_i8x16_shuffle(n36, n37, [2]uint64{p0, p0h})
	n39 := Simd_i8x16_shuffle(n35, n38, [2]uint64{p2, p2h})
	n40 := Simd_i32x4_add(n27, n39)
	n41 := Simd_i32x4_shr_s(n40, 17)
	n42 := Simd_i8x16_shuffle(n35, n38, [2]uint64{p3, p3h})
	n43 := Simd_v128_load_nc(m, s0, 80)
	n44 := Simd_i32x4_add(n43, n42)
	n45 := Simd_i32x4_shr_s(n44, 17)
	n46 := Simd_i16x8_narrow_i32x4_s(n41, n45)
	n47 := Simd_i16x8_min_s(n46, [2]uint64{p4, p4h})
	n48 := Simd_v128_bitselect([2]uint64{p1, p1h}, n47, n33)
	n49 := Simd_i16x8_add(n48, n32)
	n50 := Simd_v128_xor(n49, n32)
	n51 := Simd_v128_load_nc(m, s0, 0)
	n52 := Simd_i16x8_mul(n50, n51)
	_ = Simd_v128_store(m, s1, 0, n52)
	return n23[0], n23[1], n50[0], n50[1]
}

//go:noinline
func Simd_p_fx413(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_swizzle([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_swizzle([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_v128_or(n0, n1)
	n3 := Simd_i8x16_swizzle([2]uint64{p0, p0h}, [2]uint64{p4, p4h})
	n4 := Simd_i8x16_swizzle([2]uint64{p2, p2h}, [2]uint64{p5, p5h})
	n5 := Simd_v128_or(n3, n4)
	_ = Simd_v128_store(m, s0, 0, n2)
	_ = Simd_v128_store(m, s0, 16, n5)
	return n2[0], n2[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx414(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i16x8_abs(n0)
	n2 := Simd_i16x8_shr_s(n1, 3)
	n3 := Simd_i16x8_min_s(n2, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 0, n3)
	return n3[0], n3[1]
}

//go:noinline
func Simd_p_fx415(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 16)
	n1 := Simd_i16x8_abs(n0)
	n2 := Simd_i16x8_shr_s(n1, 3)
	n3 := Simd_i16x8_min_s(n2, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s0, 16, n3)
	return n3[0], n3[1]
}

//go:noinline
func Simd_p_fx416(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_nc(m, s0, 48)
	n1 := Simd_i16x8_mul([2]uint64{p0, p0h}, n0)
	n2 := Simd_i32x4_extmul_low_i16x8_u([2]uint64{p0, p0h}, n0)
	n3 := Simd_i32x4_extmul_high_i16x8_u([2]uint64{p0, p0h}, n0)
	n4 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p1, p1h})
	n5 := Simd_v128_load_nc(m, s0, 112)
	return n0[0], n0[1], n1[0], n1[1], n4[0], n4[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx417(m *Module, s0 int32, s1 int32, p0, p0h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_nc(m, s0, 16)
	n1 := Simd_i16x8_mul([2]uint64{p0, p0h}, n0)
	_ = Simd_v128_store(m, s1, 16, n1)
	n3 := Simd_v128_load_rng(m, s0, 64, 0, 208)
	n4 := Simd_v128_load_nc(m, s0, 192)
	n5 := Simd_v128_load(m, s1, 0)
	return n0[0], n0[1], n3[0], n3[1], n4[0], n4[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx418(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_abs([2]uint64{p1, p1h})
	n1 := Simd_i16x8_add([2]uint64{p0, p0h}, n0)
	n2 := Simd_v128_load_nc(m, s0, 32)
	n3 := Simd_i16x8_mul(n1, n2)
	n4 := Simd_i32x4_extmul_low_i16x8_u(n1, n2)
	n5 := Simd_i32x4_extmul_high_i16x8_u(n1, n2)
	n6 := Simd_i8x16_shuffle(n4, n5, [2]uint64{p2, p2h})
	n7 := Simd_v128_load_nc(m, s0, 80)
	return n2[0], n2[1], n3[0], n3[1], n6[0], n6[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx419(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_swizzle([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_swizzle([2]uint64{p0, p0h}, [2]uint64{p3, p3h})
	n2 := Simd_v128_or(n0, n1)
	n3 := Simd_i8x16_swizzle([2]uint64{p1, p1h}, [2]uint64{p4, p4h})
	n4 := Simd_i8x16_swizzle([2]uint64{p0, p0h}, [2]uint64{p5, p5h})
	n5 := Simd_v128_or(n3, n4)
	n6 := Simd_v128_load_nc(m, s0, 0)
	n7 := Simd_i16x8_mul([2]uint64{p0, p0h}, n6)
	_ = Simd_v128_store(m, s1, 0, n7)
	_ = Simd_v128_store(m, s2, 0, n2)
	_ = Simd_v128_store(m, s2, 16, n5)
	n11 := Simd_v128_load(m, s1, 32)
	return n6[0], n6[1], n2[0], n2[1], n5[0], n5[1], n11[0], n11[1]
}

//go:noinline
func Simd_p_fx420(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_abs([2]uint64{p1, p1h})
	n1 := Simd_i16x8_add([2]uint64{p0, p0h}, n0)
	n2 := Simd_i32x4_extmul_low_i16x8_u(n1, [2]uint64{p2, p2h})
	n3 := Simd_i16x8_mul([2]uint64{p2, p2h}, n1)
	n4 := Simd_i32x4_extmul_high_i16x8_u(n1, [2]uint64{p2, p2h})
	n5 := Simd_i8x16_shuffle(n2, n4, [2]uint64{p3, p3h})
	n6 := Simd_i16x8_shr_s([2]uint64{p1, p1h}, 15)
	return n3[0], n3[1], n5[0], n5[1], n6[0], n6[1]
}

//go:noinline
func Simd_p_fx421(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	_ = Simd_v128_store(m, s0, 32, n0)
	n2 := Simd_v128_load(m, s0, 48)
	n3 := Simd_i16x8_abs(n2)
	n4 := Simd_i16x8_add([2]uint64{p2, p2h}, n3)
	n5 := Simd_i32x4_extmul_low_i16x8_u(n4, [2]uint64{p3, p3h})
	n6 := Simd_i16x8_shr_s(n2, 15)
	n7 := Simd_i16x8_mul([2]uint64{p3, p3h}, n4)
	n8 := Simd_i32x4_extmul_high_i16x8_u(n4, [2]uint64{p3, p3h})
	n9 := Simd_i8x16_shuffle(n5, n8, [2]uint64{p4, p4h})
	return n2[0], n2[1], n7[0], n7[1], n9[0], n9[1], n6[0], n6[1]
}

//go:noinline
func Simd_p_fx422(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_swizzle([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_swizzle([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n3 := Simd_v128_or(n1, n2)
	_ = Simd_v128_store(m, s0, 48, n0)
	_ = Simd_v128_store(m, s1, 48, n3)
	return n3[0], n3[1]
}

//go:noinline
func Simd_p_fx423(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_swizzle([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_swizzle([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_v128_or(n0, n1)
	_ = Simd_v128_store(m, s0, 32, n2)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx424(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 96, 16, 112)
	n1 := Simd_v128_load_nc(m, s0, 48)
	n2 := Simd_v128_load(m, s1, 16)
	n3 := Simd_i16x8_abs(n2)
	n4 := Simd_i16x8_mul(n1, n3)
	n5 := Simd_i32x4_extmul_low_i16x8_u(n1, n3)
	n6 := Simd_i32x4_extmul_high_i16x8_u(n1, n3)
	n7 := Simd_i8x16_shuffle(n5, n6, [2]uint64{p0, p0h})
	n8 := Simd_i16x8_shr_s(n2, 15)
	n9 := Simd_i16x8_eq(n2, [2]uint64{p1, p1h})
	n10 := Simd_i8x16_shuffle(n4, n7, [2]uint64{p2, p2h})
	n11 := Simd_i32x4_add(n0, n10)
	n12 := Simd_i32x4_shr_s(n11, 17)
	n13 := Simd_i8x16_shuffle(n4, n7, [2]uint64{p3, p3h})
	n14 := Simd_v128_load_nc(m, s0, 112)
	n15 := Simd_i32x4_add(n14, n13)
	n16 := Simd_i32x4_shr_s(n15, 17)
	n17 := Simd_i16x8_narrow_i32x4_s(n12, n16)
	n18 := Simd_i16x8_min_s(n17, [2]uint64{p4, p4h})
	n19 := Simd_v128_bitselect([2]uint64{p1, p1h}, n18, n9)
	n20 := Simd_i16x8_add(n19, n8)
	n21 := Simd_v128_xor(n20, n8)
	n22 := Simd_v128_load_nc(m, s0, 16)
	n23 := Simd_i16x8_mul(n21, n22)
	_ = Simd_v128_store(m, s1, 16, n23)
	n25 := Simd_v128_load_rng(m, s0, 64, 0, 96)
	n26 := Simd_v128_load_nc(m, s0, 32)
	n27 := Simd_v128_load(m, s1, 0)
	n28 := Simd_i16x8_abs(n27)
	n29 := Simd_i16x8_mul(n26, n28)
	n30 := Simd_i32x4_extmul_low_i16x8_u(n26, n28)
	n31 := Simd_i32x4_extmul_high_i16x8_u(n26, n28)
	n32 := Simd_i8x16_shuffle(n30, n31, [2]uint64{p0, p0h})
	n33 := Simd_i16x8_shr_s(n27, 15)
	n34 := Simd_i16x8_eq(n27, [2]uint64{p1, p1h})
	n35 := Simd_i8x16_shuffle(n29, n32, [2]uint64{p2, p2h})
	n36 := Simd_i32x4_add(n25, n35)
	n37 := Simd_i32x4_shr_s(n36, 17)
	n38 := Simd_i8x16_shuffle(n29, n32, [2]uint64{p3, p3h})
	n39 := Simd_v128_load_nc(m, s0, 80)
	n40 := Simd_i32x4_add(n39, n38)
	n41 := Simd_i32x4_shr_s(n40, 17)
	n42 := Simd_i16x8_narrow_i32x4_s(n37, n41)
	n43 := Simd_i16x8_min_s(n42, [2]uint64{p4, p4h})
	n44 := Simd_v128_bitselect([2]uint64{p1, p1h}, n43, n34)
	n45 := Simd_i16x8_add(n44, n33)
	n46 := Simd_v128_xor(n45, n33)
	n47 := Simd_v128_load_nc(m, s0, 0)
	n48 := Simd_i16x8_mul(n46, n47)
	_ = Simd_v128_store(m, s1, 0, n48)
	return n21[0], n21[1], n46[0], n46[1]
}

//go:noinline
func Simd_p_fx425(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i16x8_extend_low_i8x16_u(n0)
	n2 := Simd_v128_load_rng(m, s0, 32, 0, 80)
	n3 := Simd_v128_load_rng(m, s1, 32, 0, 80)
	n4 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p2, p2h})
	n5 := Simd_i16x8_extend_low_i8x16_u(n4)
	n6 := Simd_i16x8_add(n1, n5)
	n7 := Simd_v128_load_nc(m, s0, 64)
	n8 := Simd_v128_load_nc(m, s1, 64)
	n9 := Simd_i8x16_shuffle(n7, n8, [2]uint64{p2, p2h})
	n10 := Simd_i16x8_extend_low_i8x16_u(n9)
	return n1[0], n1[1], n5[0], n5[1], n6[0], n6[1], n10[0], n10[1]
}

//go:noinline
func Simd_p_fx426(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_sub([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n1 := Simd_v128_load_nc(m, s0, 0)
	n2 := Simd_v128_load_nc(m, s1, 0)
	n3 := Simd_i8x16_shuffle(n1, n2, [2]uint64{p0, p0h})
	n4 := Simd_i16x8_extend_low_i8x16_u(n3)
	n5 := Simd_i16x8_sub(n4, [2]uint64{p1, p1h})
	n6 := Simd_i16x8_add([2]uint64{p1, p1h}, n4)
	n7 := Simd_i16x8_add([2]uint64{p2, p2h}, n6)
	return n6[0], n6[1], n7[0], n7[1], n0[0], n0[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx427(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, n0, [2]uint64{p3, p3h})
	n2 := Simd_i16x8_sub([2]uint64{p1, p1h}, [2]uint64{p0, p0h})
	n3 := Simd_i16x8_sub([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx428(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i16x8_extend_low_i8x16_u(n0)
	n2 := Simd_v128_load_rng(m, s0+36, 0, -32, 80)
	n3 := Simd_v128_load_rng(m, s1+36, 0, -32, 80)
	n4 := Simd_i8x16_shuffle(n2, n3, [2]uint64{p2, p2h})
	n5 := Simd_i16x8_extend_low_i8x16_u(n4)
	n6 := Simd_i16x8_add(n1, n5)
	n7 := Simd_v128_load_nc(m, s0+68, 0)
	n8 := Simd_v128_load_nc(m, s1+68, 0)
	n9 := Simd_i8x16_shuffle(n7, n8, [2]uint64{p2, p2h})
	n10 := Simd_i16x8_extend_low_i8x16_u(n9)
	return n1[0], n1[1], n5[0], n5[1], n6[0], n6[1], n10[0], n10[1]
}

//go:noinline
func Simd_p_fx429(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_sub([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n1 := Simd_v128_load_nc(m, s0+4, 0)
	n2 := Simd_v128_load_nc(m, s1+4, 0)
	n3 := Simd_i8x16_shuffle(n1, n2, [2]uint64{p0, p0h})
	n4 := Simd_i16x8_extend_low_i8x16_u(n3)
	n5 := Simd_i16x8_sub(n4, [2]uint64{p1, p1h})
	n6 := Simd_i16x8_add([2]uint64{p1, p1h}, n4)
	n7 := Simd_i16x8_add([2]uint64{p2, p2h}, n6)
	return n6[0], n6[1], n7[0], n7[1], n0[0], n0[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx430(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_i32x4_max_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_min_s(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_max_s([2]uint64{p3, p3h}, [2]uint64{p1, p1h})
	n3 := Simd_i32x4_min_s(n2, [2]uint64{p2, p2h})
	n4 := Simd_i16x8_narrow_i32x4_u(n1, n3)
	n5 := Simd_i32x4_max_s([2]uint64{p4, p4h}, [2]uint64{p1, p1h})
	n6 := Simd_i32x4_min_s(n5, [2]uint64{p2, p2h})
	n7 := Simd_i32x4_max_s([2]uint64{p5, p5h}, [2]uint64{p1, p1h})
	n8 := Simd_i32x4_min_s(n7, [2]uint64{p2, p2h})
	n9 := Simd_i16x8_narrow_i32x4_u(n6, n8)
	n10 := Simd_i8x16_narrow_i16x8_u(n4, n9)
	_ = Simd_v128_store(m, s0, 0, n10)
	return
}

//go:noinline
func Simd_p_fx431(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) {
	n0 := Simd_i32x4_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_max_s(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_min_s(n1, [2]uint64{p3, p3h})
	n3 := Simd_i32x4_add([2]uint64{p4, p4h}, [2]uint64{p1, p1h})
	n4 := Simd_i32x4_max_s(n3, [2]uint64{p2, p2h})
	n5 := Simd_i32x4_min_s(n4, [2]uint64{p3, p3h})
	n6 := Simd_i16x8_narrow_i32x4_u(n2, n5)
	n7 := Simd_i32x4_add([2]uint64{p5, p5h}, [2]uint64{p1, p1h})
	n8 := Simd_i32x4_max_s(n7, [2]uint64{p2, p2h})
	n9 := Simd_i32x4_min_s(n8, [2]uint64{p3, p3h})
	n10 := Simd_i32x4_add([2]uint64{p6, p6h}, [2]uint64{p1, p1h})
	n11 := Simd_i32x4_max_s(n10, [2]uint64{p2, p2h})
	n12 := Simd_i32x4_min_s(n11, [2]uint64{p3, p3h})
	n13 := Simd_i16x8_narrow_i32x4_u(n9, n12)
	n14 := Simd_i8x16_narrow_i16x8_u(n6, n13)
	_ = Simd_v128_store(m, s0, 0, n14)
	return
}

//go:noinline
func Simd_p_fx432(m *Module, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_extend_low_i8x16_u([2]uint64{p0, p0h})
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx433(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_extend_low_i8x16_u([2]uint64{p0, p0h})
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	n2 := Simd_i32x4_add([2]uint64{p1, p1h}, n1)
	return n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx434(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_extend_low_i8x16_u([2]uint64{p0, p0h})
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	n2 := Simd_i32x4_add([2]uint64{p1, p1h}, n1)
	n3 := Simd_i32x4_add([2]uint64{p2, p2h}, n2)
	return n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx435(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_extend_low_i8x16_u([2]uint64{p0, p0h})
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	n2 := Simd_i32x4_add([2]uint64{p1, p1h}, n1)
	n3 := Simd_i32x4_add([2]uint64{p2, p2h}, n2)
	n4 := Simd_i32x4_add([2]uint64{p3, p3h}, n3)
	return n1[0], n1[1], n2[0], n2[1], n3[0], n3[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx436(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_extend_low_i8x16_u([2]uint64{p0, p0h})
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	n2 := Simd_i32x4_sub([2]uint64{p1, p1h}, n1)
	return n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx437(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_extend_low_i8x16_u([2]uint64{p0, p0h})
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	n2 := Simd_i32x4_sub([2]uint64{p1, p1h}, n1)
	n3 := Simd_i32x4_add([2]uint64{p2, p2h}, n2)
	return n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx438(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i16x8_extend_low_i8x16_u([2]uint64{p0, p0h})
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	n2 := Simd_i32x4_sub([2]uint64{p1, p1h}, n1)
	n3 := Simd_i32x4_add([2]uint64{p2, p2h}, n2)
	n4 := Simd_i32x4_add([2]uint64{p3, p3h}, n3)
	return n1[0], n1[1], n2[0], n2[1], n3[0], n3[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx439(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_add([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i32x4_add(n0, n1)
	n3 := Simd_i32x4_add([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx440(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_add([2]uint64{p2, p2h}, n0)
	n2 := Simd_i32x4_add([2]uint64{p3, p3h}, n1)
	n3 := Simd_i32x4_add([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx441(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_add([2]uint64{p2, p2h}, n0)
	n2 := Simd_i32x4_add([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n3 := Simd_i32x4_add([2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx442(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_sub(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_add([2]uint64{p2, p2h}, n0)
	n3 := Simd_i32x4_sub([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n4 := Simd_i32x4_sub([2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	return n2[0], n2[1], n1[0], n1[1], n3[0], n3[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx443(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_sub([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i32x4_sub([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	n3 := Simd_i32x4_add(n1, n2)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx444(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_sub([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i32x4_sub(n1, n0)
	n3 := Simd_i32x4_add(n0, n1)
	n4 := Simd_i32x4_sub([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	return n3[0], n3[1], n4[0], n4[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx445(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_sub([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i32x4_add(n0, n1)
	n3 := Simd_i32x4_sub([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx446(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_sub(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_add([2]uint64{p2, p2h}, n0)
	n3 := Simd_i32x4_sub([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n4 := Simd_i32x4_sub([2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	return n2[0], n2[1], n3[0], n3[1], n1[0], n1[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx447(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_sub([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i32x4_add(n0, n1)
	n3 := Simd_i32x4_add([2]uint64{p4, p4h}, n2)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx448(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_add([2]uint64{p2, p2h}, n0)
	n2 := Simd_i32x4_sub(n1, [2]uint64{p3, p3h})
	n3 := Simd_i32x4_add([2]uint64{p3, p3h}, n1)
	n4 := Simd_i32x4_sub([2]uint64{p4, p4h}, [2]uint64{p5, p5h})
	return n0[0], n0[1], n3[0], n3[1], n4[0], n4[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx449(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_sub(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_add([2]uint64{p2, p2h}, n0)
	n3 := Simd_i32x4_sub([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n4 := Simd_i32x4_sub([2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n5 := Simd_i32x4_sub(n4, n3)
	n6 := Simd_i32x4_add(n3, n4)
	return n2[0], n2[1], n6[0], n6[1], n1[0], n1[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx450(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_mul(n0, [2]uint64{p2, p2h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx451(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_abs(n0)
	n2 := Simd_i32x4_mul(n1, [2]uint64{p2, p2h})
	n3 := Simd_i32x4_add([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n4 := Simd_i32x4_mul(n3, [2]uint64{p5, p5h})
	n5 := Simd_i32x4_add(n2, n4)
	n6 := Simd_i32x4_sub([2]uint64{p4, p4h}, [2]uint64{p3, p3h})
	n7 := Simd_i32x4_abs(n6)
	n8 := Simd_i32x4_mul(n7, [2]uint64{p6, p6h})
	n9 := Simd_i32x4_add(n5, n8)
	return n9[0], n9[1]
}

//go:noinline
func Simd_p_fx452(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_abs(n0)
	n2 := Simd_i32x4_mul(n1, [2]uint64{p2, p2h})
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx453(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_abs(n0)
	n2 := Simd_i32x4_mul(n1, [2]uint64{p2, p2h})
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx454(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1084818905618843912, 216736831629295872})
	n1 := Simd_i32x4_add([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx455(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{216736831696667908, 216736831629295872})
	n1 := Simd_i32x4_add([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx456(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_i16x8_extend_low_i8x16_u(n1)
	n3 := Simd_i32x4_extend_low_i16x8_u(n2)
	n4 := Simd_i8x16_shuffle(n1, n1, [2]uint64{p0, p0h})
	n5 := Simd_i16x8_extend_low_i8x16_u(n4)
	n6 := Simd_i32x4_extend_low_i16x8_u(n5)
	n7 := Simd_i8x16_shuffle(n1, n1, [2]uint64{p1, p1h})
	n8 := Simd_i16x8_extend_low_i8x16_u(n7)
	n9 := Simd_i32x4_extend_low_i16x8_u(n8)
	n10 := Simd_i8x16_shuffle(n1, n1, [2]uint64{p2, p2h})
	n11 := Simd_i16x8_extend_low_i8x16_u(n10)
	n12 := Simd_i32x4_extend_low_i16x8_u(n11)
	n13 := Simd_scalar_i32_add(s2, s1)
	n14 := Simd_v128_load(m, n13, 0)
	n15 := Simd_i16x8_extend_low_i8x16_u(n14)
	n16 := Simd_i32x4_extend_low_i16x8_u(n15)
	n17 := Simd_i32x4_sub(n3, n16)
	n18 := Simd_i32x4_mul(n17, n17)
	n19 := Simd_i8x16_shuffle(n14, n1, [2]uint64{p0, p0h})
	n20 := Simd_i16x8_extend_low_i8x16_u(n19)
	n21 := Simd_i32x4_extend_low_i16x8_u(n20)
	n22 := Simd_i32x4_sub(n6, n21)
	n23 := Simd_i32x4_mul(n22, n22)
	n24 := Simd_i32x4_add(n18, n23)
	n25 := Simd_i8x16_shuffle(n14, n1, [2]uint64{p1, p1h})
	n26 := Simd_i16x8_extend_low_i8x16_u(n25)
	n27 := Simd_i32x4_extend_low_i16x8_u(n26)
	n28 := Simd_i32x4_sub(n9, n27)
	n29 := Simd_i32x4_mul(n28, n28)
	n30 := Simd_i8x16_shuffle(n14, n1, [2]uint64{p2, p2h})
	n31 := Simd_i16x8_extend_low_i8x16_u(n30)
	n32 := Simd_i32x4_extend_low_i16x8_u(n31)
	n33 := Simd_i32x4_sub(n12, n32)
	n34 := Simd_i32x4_mul(n33, n33)
	n35 := Simd_i32x4_add(n24, n29)
	n36 := Simd_i32x4_add(n35, n34)
	n37 := Simd_i8x16_shuffle(n36, n36, [2]uint64{1084818905618843912, 216736831629295872})
	n38 := Simd_i32x4_add(n36, n37)
	return n38[0], n38[1]
}

//go:noinline
func Simd_p_fx457(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, s5 int32, s6 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	_ = Simd_v128_store(m, s1, 0, n0)
	_ = Simd_v128_store(m, s2, 0, n0)
	_ = Simd_v128_store(m, s3, 0, n0)
	_ = Simd_v128_store(m, s4, 0, n0)
	_ = Simd_v128_store(m, s5, 0, n0)
	_ = Simd_v128_store(m, s6, 0, n0)
	return
}

//go:noinline
func Simd_p_fx458(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	_ = Simd_v128_store(m, s1, 0, n0)
	n2 := Simd_v128_load(m, s0, 32)
	_ = Simd_v128_store(m, s1, 32, n2)
	n4 := Simd_v128_load(m, s0, 64)
	_ = Simd_v128_store(m, s1, 64, n4)
	n6 := Simd_v128_load(m, s0, 96)
	_ = Simd_v128_store(m, s1, 96, n6)
	n8 := Simd_v128_load(m, s0, 128)
	_ = Simd_v128_store(m, s1, 128, n8)
	n10 := Simd_v128_load(m, s0, 160)
	_ = Simd_v128_store(m, s1, 160, n10)
	n12 := Simd_v128_load(m, s0, 192)
	_ = Simd_v128_store(m, s1, 192, n12)
	n14 := Simd_v128_load(m, s0, 224)
	_ = Simd_v128_store(m, s1, 224, n14)
	return
}

//go:noinline
func Simd_p_fx459(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i16x8_abs(n0)
	_ = Simd_v128_store(m, s1, 0, n1)
	n3 := Simd_v128_load(m, s0, 16)
	n4 := Simd_i16x8_abs(n3)
	n5 := Simd_i8x16_narrow_i16x8_s(n1, n4)
	n6 := Simd_i8x16_min_u(n5, [2]uint64{p0, p0h})
	n7 := Simd_i8x16_min_u(n5, [2]uint64{p1, p1h})
	_ = Simd_v128_store(m, s1, 16, n4)
	_ = Simd_v128_store(m, s1, 32, n6)
	_ = Simd_v128_store(m, s1, 48, n7)
	return
}

//go:noinline
func Simd_p_fx460(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_narrow_i16x8_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_eq(n0, [2]uint64{p2, p2h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx461(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_f64x2_mul([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_f64x2_add([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx462(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_f64x2_sub([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_f64x2_mul(n0, [2]uint64{p2, p2h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx463(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_min_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_max_s(n0, [2]uint64{p2, p2h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx464(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	_ = Simd_v128_store(m, s2, 0, n1)
	return
}

//go:noinline
func Simd_p_fx465(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{72058693566333184, 216736831595610368})
	_ = Simd_v128_store(m, s0, 12, n0)
	return
}

//go:noinline
func Simd_p_fx466(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load(m, s0, 64)
	_ = Simd_v128_store(m, s1, 0, n0)
	n2 := Simd_v128_load(m, s0, 80)
	_ = Simd_v128_store(m, s1, 16, n2)
	return
}

//go:noinline
func Simd_p_fx467(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	_ = Simd_v128_store(m, s1, 0, n0)
	n2 := Simd_v128_load(m, s0+16, 0)
	_ = Simd_v128_store(m, s1+16, 0, n2)
	return
}

//go:noinline
func Simd_p_fx468(m *Module, s0 int32, s1 int32, s2 int32) (uint64, uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, s2)
	_ = Simd_v128_store(m, s0, s2, n1)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx469(m *Module, s0 int32, s1 int32) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0+-32, s1)
	_ = Simd_v128_store(m, s0, s1, n0)
	return n0[0], n0[1]
}

//go:noinline
func Simd_p_fx470(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	_ = Simd_v128_store(m, s1, 0, n0)
	return
}

//go:noinline
func Simd_p_fx471(m *Module, s0 int32, s1 int32, s2 int32) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, s1)
	_ = Simd_v128_store(m, s2, s1, n0)
	return n0[0], n0[1]
}

//go:noinline
func Simd_p_fx472(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load(m, s0, 480)
	_ = Simd_v128_store(m, s1, 0, n0)
	return
}

//go:noinline
func Simd_p_fx473(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load(m, s0, 240)
	_ = Simd_v128_store(m, s1, 0, n0)
	return
}

//go:noinline
func Simd_p_fx474(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_splat(s0)
	n1 := Simd_i32x4_mul([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx475(m *Module, s0 int32, s1 int32, p0, p0h uint64) {
	n0 := Simd_i32x4_splat(s0)
	n1 := Simd_i32x4_add([2]uint64{p0, p0h}, n0)
	_ = Simd_v128_store(m, s1, 0, n1)
	return
}

//go:noinline
func Simd_p_fx476(m *Module, s0 int32) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_splat(s0)
	return n0[0], n0[1], n0[0], n0[1]
}

//go:noinline
func Simd_p_fx477(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i32x4_max_s([2]uint64{p0, p0h}, n0)
	n2 := Simd_i32x4_min_s([2]uint64{p1, p1h}, n0)
	return n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx478(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i32x4_max_s([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx479(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_min_s([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx480(m *Module, s0 int32, s1 int32, p0, p0h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_splat(s1)
	n1 := Simd_v128_load(m, s0, 0)
	n2 := Simd_i32x4_sub(n1, n0)
	n3 := Simd_i32x4_mul(n2, [2]uint64{p0, p0h})
	return n1[0], n1[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx481(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_min_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_max_s(n0, [2]uint64{p2, p2h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx482(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_splat(s0)
	n1 := Simd_i32x4_sub([2]uint64{p0, p0h}, n0)
	n2 := Simd_i32x4_mul(n1, [2]uint64{p1, p1h})
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx483(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_min_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_max_s(n0, [2]uint64{p2, p2h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx484(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n0, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n1, [2]uint64{p2, p2h})
	n3 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n2, [2]uint64{p2, p2h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx485(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n0, [2]uint64{p2, p2h})
	return n0[0], n0[1], n1[0], n1[1]
}

//go:noinline
func Simd_p_fx486(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_narrow_i16x8_u(n0, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_add(n1, [2]uint64{p3, p3h})
	n3 := Simd_v128_and(n2, [2]uint64{p4, p4h})
	return n3[0], n3[1]
}

//go:noinline
func Simd_p_fx487(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p1, p1h}, [2]uint64{p4, p4h})
	n2 := Simd_i16x8_add([2]uint64{p0, p0h}, n1)
	n3 := Simd_i8x16_narrow_i16x8_u(n2, [2]uint64{p5, p5h})
	n4 := Simd_i8x16_add(n3, [2]uint64{p6, p6h})
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx488(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_v128_or([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_or(n0, [2]uint64{p2, p2h})
	n2 := Simd_v128_or(n1, [2]uint64{p3, p3h})
	n3 := Simd_v128_or(n2, [2]uint64{p4, p4h})
	n4 := Simd_v128_or(n3, [2]uint64{p5, p5h})
	n5 := Simd_v128_or(n4, [2]uint64{p6, p6h})
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx489(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, n0, [2]uint64{1591200115485380623, 2169921498189994007})
	n2 := Simd_i8x16_add(n1, n0)
	n3 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, n2, [2]uint64{1518859942647303950, 2097581325351917334})
	n4 := Simd_i8x16_add(n3, n2)
	n5 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, n4, [2]uint64{1374179596971150604, 1952900979675763988})
	n6 := Simd_i8x16_add(n5, n4)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx490(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32) {
	n0 := Simd_v128_load_rng(m, s0+16, 0, -16, 32)
	n1 := Simd_v128_load_rng(m, s1+16, 0, -16, 32)
	n2 := Simd_i8x16_add(n0, n1)
	n3 := Simd_v128_load_nc(m, s0, 0)
	n4 := Simd_v128_load_nc(m, s1, 0)
	n5 := Simd_i8x16_add(n3, n4)
	n6 := Simd_scalar_i32_add(s2, s3)
	_ = Simd_v128_store(m, n6, 0, n5)
	_ = Simd_v128_store(m, s4+16, 0, n2)
	return
}

//go:noinline
func Simd_p_fx491(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1, 0)
	n2 := Simd_i8x16_add(n0, n1)
	_ = Simd_v128_store(m, s2, 0, n2)
	return
}

//go:noinline
func Simd_p_fx492(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load_rng(m, s0+17, 0, -1, 17)
	n1 := Simd_v128_load_nc(m, s0+16, 0)
	n2 := Simd_i8x16_sub(n0, n1)
	_ = Simd_v128_store(m, s1+17, 0, n2)
	n4 := Simd_v128_load_rng(m, s0+1, 0, -1, 17)
	n5 := Simd_v128_load_nc(m, s0, 0)
	n6 := Simd_i8x16_sub(n4, n5)
	_ = Simd_v128_store(m, s1+1, 0, n6)
	n8 := Simd_v128_load_rng(m, s0+49, 0, -1, 17)
	n9 := Simd_v128_load_nc(m, s0+48, 0)
	n10 := Simd_i8x16_sub(n8, n9)
	_ = Simd_v128_store(m, s1+49, 0, n10)
	n12 := Simd_v128_load_rng(m, s0+33, 0, -1, 17)
	n13 := Simd_v128_load_nc(m, s0+32, 0)
	n14 := Simd_i8x16_sub(n12, n13)
	_ = Simd_v128_store(m, s1+33, 0, n14)
	return
}

//go:noinline
func Simd_p_fx493(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load_rng(m, s0, 16, 15, 17)
	n1 := Simd_v128_load_nc(m, s0, 15)
	n2 := Simd_i8x16_sub(n0, n1)
	_ = Simd_v128_store(m, s1, 16, n2)
	n4 := Simd_v128_load(m, s0, 0)
	n5 := Simd_v128_load(m, s0+-1, 0)
	n6 := Simd_i8x16_sub(n4, n5)
	_ = Simd_v128_store(m, s1, 0, n6)
	return
}

//go:noinline
func Simd_p_fx494(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load_rng(m, s0+1, 0, -1, 17)
	n1 := Simd_v128_load_nc(m, s0, 0)
	n2 := Simd_i8x16_sub(n0, n1)
	_ = Simd_v128_store(m, s1, 0, n2)
	return
}

//go:noinline
func Simd_p_fx495(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) {
	n0 := Simd_v128_load_rng(m, s0+1, 0, -1, 17)
	n1 := Simd_v128_load_rng(m, s1+1, 0, -1, 17)
	n2 := Simd_i16x8_extend_low_i8x16_u(n1)
	n3 := Simd_i32x4_extend_low_i16x8_u(n2)
	n4 := Simd_i8x16_shuffle(n1, n1, [2]uint64{p2, p2h})
	n5 := Simd_i16x8_extend_low_i8x16_u(n4)
	n6 := Simd_i32x4_extend_low_i16x8_u(n5)
	n7 := Simd_i8x16_shuffle(n1, n1, [2]uint64{p3, p3h})
	n8 := Simd_i16x8_extend_low_i8x16_u(n7)
	n9 := Simd_i32x4_extend_low_i16x8_u(n8)
	n10 := Simd_i8x16_shuffle(n1, n1, [2]uint64{p4, p4h})
	n11 := Simd_i16x8_extend_low_i8x16_u(n10)
	n12 := Simd_i32x4_extend_low_i16x8_u(n11)
	n13 := Simd_v128_load_nc(m, s0, 0)
	n14 := Simd_i16x8_extend_low_i8x16_u(n13)
	n15 := Simd_i32x4_extend_low_i16x8_u(n14)
	n16 := Simd_i32x4_add(n3, n15)
	n17 := Simd_i8x16_shuffle(n13, n1, [2]uint64{p2, p2h})
	n18 := Simd_i16x8_extend_low_i8x16_u(n17)
	n19 := Simd_i32x4_extend_low_i16x8_u(n18)
	n20 := Simd_i32x4_add(n6, n19)
	n21 := Simd_i8x16_shuffle(n13, n1, [2]uint64{p3, p3h})
	n22 := Simd_i16x8_extend_low_i8x16_u(n21)
	n23 := Simd_i32x4_extend_low_i16x8_u(n22)
	n24 := Simd_i32x4_add(n9, n23)
	n25 := Simd_i8x16_shuffle(n13, n1, [2]uint64{p4, p4h})
	n26 := Simd_i16x8_extend_low_i8x16_u(n25)
	n27 := Simd_i32x4_extend_low_i16x8_u(n26)
	n28 := Simd_i32x4_add(n12, n27)
	n29 := Simd_v128_load_nc(m, s1, 0)
	n30 := Simd_i16x8_extend_low_i8x16_u(n29)
	n31 := Simd_i32x4_extend_low_i16x8_u(n30)
	n32 := Simd_i32x4_sub(n16, n31)
	n33 := Simd_i32x4_max_s(n32, [2]uint64{p0, p0h})
	n34 := Simd_i32x4_min_s(n33, [2]uint64{p1, p1h})
	n35 := Simd_i8x16_shuffle(n29, n1, [2]uint64{p2, p2h})
	n36 := Simd_i16x8_extend_low_i8x16_u(n35)
	n37 := Simd_i32x4_extend_low_i16x8_u(n36)
	n38 := Simd_i32x4_sub(n20, n37)
	n39 := Simd_i32x4_max_s(n38, [2]uint64{p0, p0h})
	n40 := Simd_i32x4_min_s(n39, [2]uint64{p1, p1h})
	n41 := Simd_i16x8_narrow_i32x4_u(n34, n40)
	n42 := Simd_i8x16_shuffle(n29, n1, [2]uint64{p3, p3h})
	n43 := Simd_i16x8_extend_low_i8x16_u(n42)
	n44 := Simd_i32x4_extend_low_i16x8_u(n43)
	n45 := Simd_i32x4_sub(n24, n44)
	n46 := Simd_i32x4_max_s(n45, [2]uint64{p0, p0h})
	n47 := Simd_i32x4_min_s(n46, [2]uint64{p1, p1h})
	n48 := Simd_i8x16_shuffle(n29, n1, [2]uint64{p4, p4h})
	n49 := Simd_i16x8_extend_low_i8x16_u(n48)
	n50 := Simd_i32x4_extend_low_i16x8_u(n49)
	n51 := Simd_i32x4_sub(n28, n50)
	n52 := Simd_i32x4_max_s(n51, [2]uint64{p0, p0h})
	n53 := Simd_i32x4_min_s(n52, [2]uint64{p1, p1h})
	n54 := Simd_i16x8_narrow_i32x4_u(n47, n53)
	n55 := Simd_i8x16_narrow_i16x8_u(n41, n54)
	n56 := Simd_i8x16_sub(n0, n55)
	n57 := Simd_scalar_i32_add(s2, s3)
	_ = Simd_v128_store(m, n57, 0, n56)
	return
}

//go:noinline
func Simd_p_fx496(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n1 := Simd_i8x16_shuffle([2]uint64{p4, p4h}, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i16x8_add(n0, n1)
	n3 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n4 := Simd_i16x8_sub(n2, n3)
	n5 := Simd_i8x16_narrow_i16x8_u(n4, [2]uint64{p2, p2h})
	n6 := Simd_i8x16_sub([2]uint64{p0, p0h}, n5)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx497(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) {
	n0 := Simd_v128_load_rng(m, s0+1, 0, -1, 17)
	n1 := Simd_v128_load_rng(m, s1+1, 0, -1, 17)
	n2 := Simd_i16x8_extend_low_i8x16_u(n1)
	n3 := Simd_i32x4_extend_low_i16x8_u(n2)
	n4 := Simd_i8x16_shuffle(n1, n1, [2]uint64{p2, p2h})
	n5 := Simd_i16x8_extend_low_i8x16_u(n4)
	n6 := Simd_i32x4_extend_low_i16x8_u(n5)
	n7 := Simd_i8x16_shuffle(n1, n1, [2]uint64{p3, p3h})
	n8 := Simd_i16x8_extend_low_i8x16_u(n7)
	n9 := Simd_i32x4_extend_low_i16x8_u(n8)
	n10 := Simd_i8x16_shuffle(n1, n1, [2]uint64{p4, p4h})
	n11 := Simd_i16x8_extend_low_i8x16_u(n10)
	n12 := Simd_i32x4_extend_low_i16x8_u(n11)
	n13 := Simd_v128_load_nc(m, s0, 0)
	n14 := Simd_i16x8_extend_low_i8x16_u(n13)
	n15 := Simd_i32x4_extend_low_i16x8_u(n14)
	n16 := Simd_i32x4_add(n3, n15)
	n17 := Simd_i8x16_shuffle(n13, n1, [2]uint64{p2, p2h})
	n18 := Simd_i16x8_extend_low_i8x16_u(n17)
	n19 := Simd_i32x4_extend_low_i16x8_u(n18)
	n20 := Simd_i32x4_add(n6, n19)
	n21 := Simd_i8x16_shuffle(n13, n1, [2]uint64{p3, p3h})
	n22 := Simd_i16x8_extend_low_i8x16_u(n21)
	n23 := Simd_i32x4_extend_low_i16x8_u(n22)
	n24 := Simd_i32x4_add(n9, n23)
	n25 := Simd_i8x16_shuffle(n13, n1, [2]uint64{p4, p4h})
	n26 := Simd_i16x8_extend_low_i8x16_u(n25)
	n27 := Simd_i32x4_extend_low_i16x8_u(n26)
	n28 := Simd_i32x4_add(n12, n27)
	n29 := Simd_v128_load_nc(m, s1, 0)
	n30 := Simd_i16x8_extend_low_i8x16_u(n29)
	n31 := Simd_i32x4_extend_low_i16x8_u(n30)
	n32 := Simd_i32x4_sub(n16, n31)
	n33 := Simd_i32x4_max_s(n32, [2]uint64{p0, p0h})
	n34 := Simd_i32x4_min_s(n33, [2]uint64{p1, p1h})
	n35 := Simd_i8x16_shuffle(n29, n1, [2]uint64{p2, p2h})
	n36 := Simd_i16x8_extend_low_i8x16_u(n35)
	n37 := Simd_i32x4_extend_low_i16x8_u(n36)
	n38 := Simd_i32x4_sub(n20, n37)
	n39 := Simd_i32x4_max_s(n38, [2]uint64{p0, p0h})
	n40 := Simd_i32x4_min_s(n39, [2]uint64{p1, p1h})
	n41 := Simd_i16x8_narrow_i32x4_u(n34, n40)
	n42 := Simd_i8x16_shuffle(n29, n1, [2]uint64{p3, p3h})
	n43 := Simd_i16x8_extend_low_i8x16_u(n42)
	n44 := Simd_i32x4_extend_low_i16x8_u(n43)
	n45 := Simd_i32x4_sub(n24, n44)
	n46 := Simd_i32x4_max_s(n45, [2]uint64{p0, p0h})
	n47 := Simd_i32x4_min_s(n46, [2]uint64{p1, p1h})
	n48 := Simd_i8x16_shuffle(n29, n1, [2]uint64{p4, p4h})
	n49 := Simd_i16x8_extend_low_i8x16_u(n48)
	n50 := Simd_i32x4_extend_low_i16x8_u(n49)
	n51 := Simd_i32x4_sub(n28, n50)
	n52 := Simd_i32x4_max_s(n51, [2]uint64{p0, p0h})
	n53 := Simd_i32x4_min_s(n52, [2]uint64{p1, p1h})
	n54 := Simd_i16x8_narrow_i32x4_u(n47, n53)
	n55 := Simd_i8x16_narrow_i16x8_u(n41, n54)
	n56 := Simd_i8x16_sub(n0, n55)
	_ = Simd_v128_store(m, s2, 0, n56)
	return
}

//go:noinline
func Simd_p_fx498(m *Module, s0 int32, s1 int32, s2 int32, s3 int32) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_scalar_i32_add(n0, s2)
	n2 := Simd_v128_load(m, n1, 0)
	n3 := Simd_scalar_i32_add(s0, s2)
	n4 := Simd_v128_load(m, n3, 0)
	n5 := Simd_i8x16_sub(n2, n4)
	n6 := Simd_scalar_i32_add(s3, s1)
	n7 := Simd_scalar_i32_add(n6, s2)
	_ = Simd_v128_store(m, n7, 0, n5)
	return
}

//go:noinline
func Simd_p_fx499(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_scalar_i32_add(s1, s2)
	n2 := Simd_v128_load(m, n1, 0)
	n3 := Simd_i8x16_sub(n0, n2)
	_ = Simd_v128_store(m, s3, 0, n3)
	n5 := Simd_v128_load(m, s0+16, 0)
	n6 := Simd_v128_load(m, s4+16, 0)
	n7 := Simd_i8x16_sub(n5, n6)
	_ = Simd_v128_store(m, s3+16, 0, n7)
	return
}

//go:noinline
func Simd_p_fx500(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1, 0)
	n2 := Simd_i8x16_sub(n0, n1)
	_ = Simd_v128_store(m, s2, 0, n2)
	return
}

//go:noinline
func Simd_p_fx501(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_i8x16_shuffle(n1, n1, [2]uint64{100992003, 0})
	n3 := Simd_i16x8_extend_low_i8x16_u(n2)
	n4 := Simd_i32x4_extend_low_i16x8_u(n3)
	n5 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n1, [2]uint64{1591200115485380623, 2169921498189994007})
	n6 := Simd_i16x8_extend_low_i8x16_u(n5)
	n7 := Simd_i32x4_extend_low_i16x8_u(n6)
	n8 := Simd_i8x16_shuffle(n1, n1, [2]uint64{168364039, 0})
	n9 := Simd_i16x8_extend_low_i8x16_u(n8)
	n10 := Simd_i32x4_extend_low_i16x8_u(n9)
	n11 := Simd_i8x16_shuffle(n1, n1, [2]uint64{235736075, 0})
	n12 := Simd_i16x8_extend_low_i8x16_u(n11)
	n13 := Simd_i32x4_extend_low_i16x8_u(n12)
	n14 := Simd_v128_load_rng(m, s0+1, 0, -1, 17)
	n15 := Simd_i16x8_extend_low_i8x16_u(n14)
	n16 := Simd_i32x4_extend_low_i16x8_u(n15)
	n17 := Simd_i32x4_add(n16, n7)
	n18 := Simd_i8x16_shuffle(n14, n1, [2]uint64{p3, p3h})
	n19 := Simd_i16x8_extend_low_i8x16_u(n18)
	n20 := Simd_i32x4_extend_low_i16x8_u(n19)
	n21 := Simd_i32x4_add(n20, n4)
	n22 := Simd_i8x16_shuffle(n14, n1, [2]uint64{p4, p4h})
	n23 := Simd_i16x8_extend_low_i8x16_u(n22)
	n24 := Simd_i32x4_extend_low_i16x8_u(n23)
	n25 := Simd_i32x4_add(n24, n10)
	n26 := Simd_i8x16_shuffle(n14, n1, [2]uint64{p5, p5h})
	n27 := Simd_i16x8_extend_low_i8x16_u(n26)
	n28 := Simd_i32x4_extend_low_i16x8_u(n27)
	n29 := Simd_i32x4_add(n28, n13)
	n30 := Simd_v128_load_nc(m, s0, 0)
	n31 := Simd_i16x8_extend_low_i8x16_u(n30)
	n32 := Simd_i32x4_extend_low_i16x8_u(n31)
	n33 := Simd_i32x4_sub(n17, n32)
	n34 := Simd_i32x4_max_s(n33, [2]uint64{p1, p1h})
	n35 := Simd_i32x4_min_s(n34, [2]uint64{p2, p2h})
	n36 := Simd_i8x16_shuffle(n30, n1, [2]uint64{p3, p3h})
	n37 := Simd_i16x8_extend_low_i8x16_u(n36)
	n38 := Simd_i32x4_extend_low_i16x8_u(n37)
	n39 := Simd_i32x4_sub(n21, n38)
	n40 := Simd_i32x4_max_s(n39, [2]uint64{p1, p1h})
	n41 := Simd_i32x4_min_s(n40, [2]uint64{p2, p2h})
	n42 := Simd_i16x8_narrow_i32x4_u(n35, n41)
	n43 := Simd_i8x16_shuffle(n30, n1, [2]uint64{p4, p4h})
	n44 := Simd_i16x8_extend_low_i8x16_u(n43)
	n45 := Simd_i32x4_extend_low_i16x8_u(n44)
	n46 := Simd_i32x4_sub(n25, n45)
	n47 := Simd_i32x4_max_s(n46, [2]uint64{p1, p1h})
	n48 := Simd_i32x4_min_s(n47, [2]uint64{p2, p2h})
	n49 := Simd_i8x16_shuffle(n30, n1, [2]uint64{p5, p5h})
	n50 := Simd_i16x8_extend_low_i8x16_u(n49)
	n51 := Simd_i32x4_extend_low_i16x8_u(n50)
	n52 := Simd_i32x4_sub(n29, n51)
	n53 := Simd_i32x4_max_s(n52, [2]uint64{p1, p1h})
	n54 := Simd_i32x4_min_s(n53, [2]uint64{p2, p2h})
	n55 := Simd_i16x8_narrow_i32x4_u(n48, n54)
	n56 := Simd_i8x16_narrow_i16x8_u(n42, n55)
	n57 := Simd_i8x16_sub(n1, n56)
	n58 := Simd_scalar_i32_add(s2, s1)
	_ = Simd_v128_store(m, n58, 0, n57)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx502(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_i32x4_extend_low_i16x8_s([2]uint64{p2, p2h})
	n1 := Simd_i32x4_extend_high_i16x8_s([2]uint64{p2, p2h})
	n2 := Simd_i32x4_extend_low_i16x8_s([2]uint64{p4, p4h})
	n3 := Simd_i32x4_extend_high_i16x8_s([2]uint64{p4, p4h})
	n4 := Simd_v128_load(m, s0, 0)
	n5 := Simd_v128_and(n4, [2]uint64{p0, p0h})
	n6 := Simd_i8x16_shuffle(n5, [2]uint64{p1, p1h}, [2]uint64{361419384851267840, 1084818905618843912})
	n7 := Simd_i8x16_shuffle(n6, [2]uint64{p1, p1h}, [2]uint64{506097522914230528, 940140767555881224})
	n8 := Simd_i32x4_extend_low_i16x8_s(n7)
	n9 := Simd_i32x4_mul(n8, n0)
	n10 := Simd_i32x4_extend_high_i16x8_s(n7)
	n11 := Simd_i32x4_mul(n10, n1)
	n12 := Simd_i8x16_shuffle(n9, n11, [2]uint64{p3, p3h})
	n13 := Simd_i16x8_shl(n4, 8)
	n14 := Simd_i32x4_extend_low_i16x8_s(n13)
	n15 := Simd_i32x4_mul(n14, n2)
	n16 := Simd_i32x4_extend_high_i16x8_s(n13)
	n17 := Simd_i32x4_mul(n16, n3)
	n18 := Simd_i8x16_shuffle(n15, n17, [2]uint64{p3, p3h})
	n19 := Simd_i32x4_shr_u(n18, 16)
	n20 := Simd_i8x16_add(n12, n19)
	n21 := Simd_v128_and(n20, [2]uint64{p5, p5h})
	n22 := Simd_i8x16_sub(n4, n21)
	_ = Simd_v128_store(m, s0, 0, n22)
	return
}

//go:noinline
func Simd_p_fx503(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_i32x4_splat(s1)
	n1 := Simd_i32x4_splat(s2)
	n2 := Simd_i32x4_splat(s3)
	n3 := Simd_v128_load(m, s0, 0)
	n4 := Simd_i32x4_shr_u(n3, 16)
	n5 := Simd_i32x4_shl(n4, 24)
	n6 := Simd_i32x4_shr_s(n5, 24)
	n7 := Simd_i32x4_mul(n6, n2)
	n8 := Simd_i32x4_shr_u(n7, 5)
	n9 := Simd_i32x4_shl(n3, 16)
	n10 := Simd_i32x4_shr_s(n9, 24)
	n11 := Simd_i32x4_mul(n10, n0)
	n12 := Simd_i32x4_shr_u(n11, 5)
	n13 := Simd_i32x4_sub(n4, n12)
	n14 := Simd_i32x4_shl(n13, 16)
	n15 := Simd_v128_and(n14, [2]uint64{p0, p0h})
	n16 := Simd_i32x4_mul(n10, n1)
	n17 := Simd_i32x4_shr_u(n16, 5)
	n18 := Simd_i32x4_add(n17, n8)
	n19 := Simd_i32x4_sub(n3, n18)
	n20 := Simd_v128_and(n19, [2]uint64{p2, p2h})
	n21 := Simd_v128_and(n3, [2]uint64{p1, p1h})
	n22 := Simd_v128_or(n15, n21)
	n23 := Simd_v128_or(n22, n20)
	_ = Simd_v128_store(m, s0, 0, n23)
	return
}

//go:noinline
func Simd_p_fx504(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i16x8_shr_u(n0, 8)
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n4 := Simd_i8x16_sub(n0, n3)
	_ = Simd_v128_store(m, s0, 0, n4)
	n6 := Simd_v128_load(m, s0+16, 0)
	n7 := Simd_i16x8_shr_u(n6, 8)
	n8 := Simd_i8x16_shuffle(n7, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n9 := Simd_i8x16_shuffle(n8, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n10 := Simd_i8x16_sub(n6, n9)
	_ = Simd_v128_store(m, s0+16, 0, n10)
	return
}

//go:noinline
func Simd_p_fx505(m *Module, s0 int32, p0, p0h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i16x8_shr_u(n0, 8)
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p0, p0h}, [2]uint64{361419384851267840, 1084818905618843912})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p0, p0h}, [2]uint64{506097522914230528, 940140767555881224})
	n4 := Simd_i8x16_sub(n0, n3)
	_ = Simd_v128_store(m, s0, 0, n4)
	return
}

//go:noinline
func Simd_p_fx506(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i32x4_shr_u(n0, 8)
	n2 := Simd_i32x4_shl(n1, 16)
	n3 := Simd_i32x4_sub(n0, n1)
	n4 := Simd_v128_and(n3, [2]uint64{p0, p0h})
	n5 := Simd_i32x4_sub(n0, n2)
	n6 := Simd_v128_and(n5, [2]uint64{p2, p2h})
	n7 := Simd_v128_and(n0, [2]uint64{p1, p1h})
	n8 := Simd_v128_or(n4, n7)
	n9 := Simd_v128_or(n8, n6)
	_ = Simd_v128_store(m, s0, 0, n9)
	return
}

//go:noinline
func Simd_p_fx507(m *Module, s0 int32, s1 int32) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_splat(s0)
	n1 := Simd_i32x4_extend_high_i16x8_s(n0)
	n2 := Simd_i32x4_extend_low_i16x8_s(n0)
	n3 := Simd_i32x4_splat(s1)
	n4 := Simd_i32x4_extend_high_i16x8_s(n3)
	n5 := Simd_i32x4_extend_low_i16x8_s(n3)
	return n1[0], n1[1], n2[0], n2[1], n4[0], n4[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx508(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 32)
	n1 := Simd_v128_and(n0, [2]uint64{p0, p0h})
	n2 := Simd_i32x4_extend_low_i16x8_s(n1)
	n3 := Simd_i32x4_mul(n2, [2]uint64{p1, p1h})
	n4 := Simd_i32x4_extend_high_i16x8_s(n1)
	n5 := Simd_i32x4_mul(n4, [2]uint64{p2, p2h})
	n6 := Simd_i8x16_shuffle(n3, n5, [2]uint64{p3, p3h})
	n7 := Simd_i16x8_shl(n0, 8)
	n8 := Simd_i32x4_extend_low_i16x8_s(n7)
	n9 := Simd_i32x4_mul(n8, [2]uint64{p4, p4h})
	n10 := Simd_i32x4_extend_high_i16x8_s(n7)
	n11 := Simd_i32x4_mul(n10, [2]uint64{p5, p5h})
	n12 := Simd_i8x16_shuffle(n9, n11, [2]uint64{p3, p3h})
	n13 := Simd_i32x4_shr_u(n12, 16)
	n14 := Simd_i8x16_add(n6, n13)
	n15 := Simd_i8x16_sub(n0, n14)
	n16 := Simd_v128_and(n15, [2]uint64{p6, p6h})
	n17 := Simd_v128_load_nc(m, s0+16, 0)
	n18 := Simd_v128_and(n17, [2]uint64{p0, p0h})
	n19 := Simd_i32x4_extend_low_i16x8_s(n18)
	n20 := Simd_i32x4_mul(n19, [2]uint64{p1, p1h})
	n21 := Simd_i32x4_extend_high_i16x8_s(n18)
	n22 := Simd_i32x4_mul(n21, [2]uint64{p2, p2h})
	n23 := Simd_i8x16_shuffle(n20, n22, [2]uint64{p3, p3h})
	n24 := Simd_i16x8_shl(n17, 8)
	n25 := Simd_i32x4_extend_low_i16x8_s(n24)
	n26 := Simd_i32x4_mul(n25, [2]uint64{p4, p4h})
	n27 := Simd_i32x4_extend_high_i16x8_s(n24)
	n28 := Simd_i32x4_mul(n27, [2]uint64{p5, p5h})
	n29 := Simd_i8x16_shuffle(n26, n28, [2]uint64{p3, p3h})
	n30 := Simd_i32x4_shr_u(n29, 16)
	n31 := Simd_i8x16_add(n23, n30)
	n32 := Simd_i8x16_sub(n17, n31)
	n33 := Simd_v128_and(n32, [2]uint64{p6, p6h})
	n34 := Simd_i16x8_narrow_i32x4_s(n16, n33)
	return n34[0], n34[1]
}

//go:noinline
func Simd_p_fx509(m *Module, s0 int32) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_splat(s0)
	n1 := Simd_i32x4_extend_high_i16x8_s(n0)
	n2 := Simd_i32x4_extend_low_i16x8_s(n0)
	return n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx510(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 32)
	n1 := Simd_v128_and(n0, [2]uint64{p0, p0h})
	n2 := Simd_i32x4_extend_low_i16x8_s(n1)
	n3 := Simd_i32x4_mul(n2, [2]uint64{p1, p1h})
	n4 := Simd_i32x4_shr_u(n0, 16)
	n5 := Simd_i32x4_extend_high_i16x8_s(n1)
	n6 := Simd_i32x4_mul(n5, [2]uint64{p2, p2h})
	n7 := Simd_i8x16_shuffle(n3, n6, [2]uint64{p3, p3h})
	n8 := Simd_i8x16_sub(n4, n7)
	n9 := Simd_v128_and(n8, [2]uint64{p4, p4h})
	n10 := Simd_v128_load_nc(m, s0+16, 0)
	n11 := Simd_v128_and(n10, [2]uint64{p0, p0h})
	n12 := Simd_i32x4_extend_low_i16x8_s(n11)
	n13 := Simd_i32x4_mul(n12, [2]uint64{p1, p1h})
	n14 := Simd_i32x4_shr_u(n10, 16)
	n15 := Simd_i32x4_extend_high_i16x8_s(n11)
	n16 := Simd_i32x4_mul(n15, [2]uint64{p2, p2h})
	n17 := Simd_i8x16_shuffle(n13, n16, [2]uint64{p3, p3h})
	n18 := Simd_i8x16_sub(n14, n17)
	n19 := Simd_v128_and(n18, [2]uint64{p4, p4h})
	n20 := Simd_i16x8_narrow_i32x4_s(n9, n19)
	return n20[0], n20[1]
}

//go:noinline
func Simd_p_fx511(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_scalar_i32_add(s1, s2)
	n2 := Simd_v128_load(m, n1, 0)
	n3 := Simd_i32x4_add(n0, n2)
	_ = Simd_v128_store(m, s3, 0, n3)
	n5 := Simd_v128_load(m, s0+48, 0)
	n6 := Simd_v128_load(m, s4+48, 0)
	n7 := Simd_i32x4_add(n5, n6)
	_ = Simd_v128_store(m, s3+48, 0, n7)
	n9 := Simd_v128_load(m, s0+32, 0)
	n10 := Simd_v128_load(m, s4+32, 0)
	n11 := Simd_i32x4_add(n9, n10)
	_ = Simd_v128_store(m, s3+32, 0, n11)
	n13 := Simd_v128_load(m, s0+16, 0)
	n14 := Simd_v128_load(m, s4+16, 0)
	n15 := Simd_i32x4_add(n13, n14)
	_ = Simd_v128_store(m, s3+16, 0, n15)
	return
}

//go:noinline
func Simd_p_fx512(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, s5 int32, s6 int32) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_scalar_i32_add(s2, s1)
	n3 := Simd_v128_load(m, n2, 0)
	n4 := Simd_i32x4_add(n1, n3)
	n5 := Simd_scalar_i32_add(s3, s1)
	_ = Simd_v128_store(m, n5, 0, n4)
	n7 := Simd_v128_load(m, s4, 0)
	n8 := Simd_v128_load(m, s5, 0)
	n9 := Simd_i32x4_add(n7, n8)
	_ = Simd_v128_store(m, s6, 0, n9)
	return
}

//go:noinline
func Simd_p_fx513(m *Module, s0 int32, s1 int32, s2 int32, s3 int32) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_scalar_i32_add(s2, s1)
	n3 := Simd_v128_load(m, n2, 0)
	n4 := Simd_i32x4_add(n1, n3)
	n5 := Simd_scalar_i32_add(s3, s1)
	_ = Simd_v128_store(m, n5, 0, n4)
	return
}

//go:noinline
func Simd_p_fx514(m *Module, s0 int32, s1 int32, s2 int32, s3 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_scalar_i32_add(s1, s2)
	n2 := Simd_v128_load(m, n1, 0)
	n3 := Simd_i32x4_add(n0, n2)
	_ = Simd_v128_store(m, s0, 0, n3)
	n5 := Simd_v128_load(m, s0+48, 0)
	n6 := Simd_v128_load(m, s3+48, 0)
	n7 := Simd_i32x4_add(n5, n6)
	_ = Simd_v128_store(m, s0+48, 0, n7)
	n9 := Simd_v128_load(m, s0+32, 0)
	n10 := Simd_v128_load(m, s3+32, 0)
	n11 := Simd_i32x4_add(n9, n10)
	_ = Simd_v128_store(m, s0+32, 0, n11)
	n13 := Simd_v128_load(m, s0+16, 0)
	n14 := Simd_v128_load(m, s3+16, 0)
	n15 := Simd_i32x4_add(n13, n14)
	_ = Simd_v128_store(m, s0+16, 0, n15)
	return
}

//go:noinline
func Simd_p_fx515(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_scalar_i32_shl(s2, 2)
	n2 := Simd_scalar_i32_add(s1, n1)
	n3 := Simd_v128_load(m, n2, 0)
	n4 := Simd_scalar_i32_shl(s2, 2)
	n5 := Simd_scalar_i32_add(s3, n4)
	n6 := Simd_v128_load(m, n5, 0)
	n7 := Simd_i32x4_add(n3, n6)
	n8 := Simd_scalar_i32_shl(s2, 2)
	n9 := Simd_scalar_i32_add(s1, n8)
	_ = Simd_v128_store(m, n9, 0, n7)
	n11 := Simd_v128_load(m, s4, 0)
	n12 := Simd_i32x4_add(n0, n11)
	_ = Simd_v128_store(m, s0, 0, n12)
	return
}

//go:noinline
func Simd_p_fx516(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_scalar_i32_add(s1, s2)
	n2 := Simd_v128_load(m, n1, 0)
	n3 := Simd_i32x4_add(n0, n2)
	_ = Simd_v128_store(m, s0, 0, n3)
	return
}

//go:noinline
func Simd_p_fx517(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_narrow_i32x4_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i16x8_narrow_i32x4_s([2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_i8x16_narrow_i16x8_s(n0, n1)
	n3 := Simd_i8x16_gt_s(n2, [2]uint64{p4, p4h})
	return n3[0], n3[1]
}

//go:noinline
func Simd_p_fx518(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n0, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n3 := Simd_i8x16_shuffle(n1, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	n4 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n0, [2]uint64{p5, p5h})
	n5 := Simd_i8x16_shuffle(n4, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n6 := Simd_i8x16_shuffle(n4, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	_ = Simd_v128_store(m, s1, 112, n2)
	_ = Simd_v128_store(m, s1, 96, n3)
	_ = Simd_v128_store(m, s1, 80, n5)
	_ = Simd_v128_store(m, s1+64, 0, n6)
	n11 := Simd_v128_load(m, s0+-16, 0)
	n12 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n11, [2]uint64{p1, p1h})
	n13 := Simd_i8x16_shuffle(n12, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n14 := Simd_i8x16_shuffle(n12, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	n15 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n11, [2]uint64{p5, p5h})
	n16 := Simd_i8x16_shuffle(n15, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n17 := Simd_i8x16_shuffle(n15, [2]uint64{p2, p2h}, [2]uint64{p4, p4h})
	_ = Simd_v128_store(m, s1, 48, n13)
	_ = Simd_v128_store(m, s1, 32, n14)
	_ = Simd_v128_store(m, s1, 16, n16)
	_ = Simd_v128_store(m, s1, 0, n17)
	return
}

//go:noinline
func Simd_p_fx519(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n1, [2]uint64{1948679894439893000, 2238040585792199692})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n4 := Simd_i8x16_shuffle(n2, [2]uint64{p1, p1h}, [2]uint64{p3, p3h})
	n5 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n1, [2]uint64{1369958511735279616, 1659319203087586308})
	n6 := Simd_i8x16_shuffle(n5, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n7 := Simd_i8x16_shuffle(n5, [2]uint64{p1, p1h}, [2]uint64{p3, p3h})
	_ = Simd_v128_store(m, s2, 48, n3)
	_ = Simd_v128_store(m, s2, 32, n4)
	_ = Simd_v128_store(m, s2, 16, n6)
	_ = Simd_v128_store(m, s2, 0, n7)
	return
}

//go:noinline
func Simd_p_fx520(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i16x8_mul(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_and(n1, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_shuffle(n2, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n4 := Simd_i8x16_shuffle(n2, [2]uint64{p1, p1h}, [2]uint64{p3, p3h})
	_ = Simd_v128_store(m, s1, 48, n3)
	_ = Simd_v128_store(m, s1+32, 0, n4)
	n7 := Simd_v128_load(m, s0+-16, 0)
	n8 := Simd_i16x8_mul(n7, [2]uint64{p0, p0h})
	n9 := Simd_v128_and(n8, [2]uint64{p1, p1h})
	n10 := Simd_i8x16_shuffle(n9, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n11 := Simd_i8x16_shuffle(n9, [2]uint64{p1, p1h}, [2]uint64{p3, p3h})
	_ = Simd_v128_store(m, s1, 16, n10)
	_ = Simd_v128_store(m, s1, 0, n11)
	return
}

//go:noinline
func Simd_p_fx521(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_i16x8_mul(n1, [2]uint64{p0, p0h})
	n3 := Simd_v128_and(n2, [2]uint64{p1, p1h})
	n4 := Simd_i8x16_shuffle(n3, [2]uint64{p1, p1h}, [2]uint64{1952885526417115400, 2242246217769422092})
	n5 := Simd_i8x16_shuffle(n3, [2]uint64{p1, p1h}, [2]uint64{1374164143712502016, 1663524835064808708})
	_ = Simd_v128_store(m, s2, 16, n4)
	_ = Simd_v128_store(m, s2, 0, n5)
	return
}

//go:noinline
func Simd_p_fx522(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i16x8_mul(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_and(n1, [2]uint64{p1, p1h})
	n3 := Simd_i32x4_shr_u(n2, 12)
	n4 := Simd_v128_or(n2, n3)
	n5 := Simd_v128_or(n4, [2]uint64{p2, p2h})
	_ = Simd_v128_store(m, s1, 0, n5)
	n7 := Simd_v128_load(m, s0+16, 0)
	n8 := Simd_i16x8_mul(n7, [2]uint64{p0, p0h})
	n9 := Simd_v128_and(n8, [2]uint64{p1, p1h})
	n10 := Simd_i32x4_shr_u(n9, 12)
	n11 := Simd_v128_or(n9, n10)
	n12 := Simd_v128_or(n11, [2]uint64{p2, p2h})
	_ = Simd_v128_store(m, s1+16, 0, n12)
	return
}

//go:noinline
func Simd_p_fx523(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_i16x8_mul(n1, [2]uint64{p0, p0h})
	n3 := Simd_v128_and(n2, [2]uint64{p1, p1h})
	n4 := Simd_i32x4_shr_u(n3, 12)
	n5 := Simd_v128_or(n3, n4)
	n6 := Simd_v128_or(n5, [2]uint64{p2, p2h})
	_ = Simd_v128_store(m, s2, 0, n6)
	return
}

//go:noinline
func Simd_p_fx524(m *Module, s0 int32, s1 int32, p0, p0h uint64) {
	n0 := Simd_v128_load32_zero(m, s0, 0)
	n1 := Simd_i16x8_extend_low_i8x16_u(n0)
	n2 := Simd_i32x4_extend_low_i16x8_u(n1)
	n3 := Simd_i32x4_shl(n2, 8)
	n4 := Simd_v128_or(n3, [2]uint64{p0, p0h})
	_ = Simd_v128_store(m, s1, 0, n4)
	return
}

//go:noinline
func Simd_p_fx525(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1+4, 0)
	n2 := Simd_i8x16_sub(n0, n1)
	_ = Simd_v128_store(m, s2, 0, n2)
	n4 := Simd_v128_load(m, s0+16, 0)
	n5 := Simd_v128_load(m, s1+20, 0)
	n6 := Simd_i8x16_sub(n4, n5)
	_ = Simd_v128_store(m, s2+16, 0, n6)
	return
}

//go:noinline
func Simd_p_fx526(m *Module, s0 int32, s1 int32, s2 int32, s3 int32) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_v128_load(m, s2, 0)
	n3 := Simd_i8x16_sub(n1, n2)
	n4 := Simd_scalar_i32_add(s3, s1)
	_ = Simd_v128_store(m, n4, 0, n3)
	return
}

//go:noinline
func Simd_p_fx527(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1, 0)
	n2 := Simd_i8x16_sub(n0, n1)
	_ = Simd_v128_store(m, s2, 0, n2)
	n4 := Simd_v128_load(m, s0+16, 0)
	n5 := Simd_v128_load(m, s1+16, 0)
	n6 := Simd_i8x16_sub(n4, n5)
	_ = Simd_v128_store(m, s2+16, 0, n6)
	return
}

//go:noinline
func Simd_p_fx528(m *Module, s0 int32, s1 int32, s2 int32, s3 int32) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_scalar_i32_add(s2, s1)
	n3 := Simd_v128_load(m, n2, 0)
	n4 := Simd_i8x16_sub(n1, n3)
	n5 := Simd_scalar_i32_add(s3, s1)
	_ = Simd_v128_store(m, n5, 0, n4)
	return
}

//go:noinline
func Simd_p_fx529(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s0+-4, 0)
	n2 := Simd_i8x16_sub(n0, n1)
	_ = Simd_v128_store(m, s1, 0, n2)
	n4 := Simd_v128_load_rng(m, s0+16, 0, -4, 20)
	n5 := Simd_v128_load_nc(m, s0+12, 0)
	n6 := Simd_i8x16_sub(n4, n5)
	_ = Simd_v128_store(m, s1+16, 0, n6)
	return
}

//go:noinline
func Simd_p_fx530(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s0+-4, 0)
	n2 := Simd_i8x16_sub(n0, n1)
	n3 := Simd_scalar_i32_add(s1, s2)
	_ = Simd_v128_store(m, n3, 0, n2)
	return
}

//go:noinline
func Simd_p_fx531(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, p0, p0h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s0+-4, 0)
	n2 := Simd_scalar_i32_add(s1, s2)
	n3 := Simd_v128_load(m, n2, 0)
	n4 := Simd_v128_xor(n3, n1)
	n5 := Simd_v128_and(n4, [2]uint64{p0, p0h})
	n6 := Simd_i8x16_avgr_u(n1, n3)
	n7 := Simd_i8x16_sub(n0, n6)
	n8 := Simd_i8x16_add(n7, n5)
	_ = Simd_v128_store(m, s3, 0, n8)
	n10 := Simd_v128_load_rng(m, s0+16, 0, -4, 20)
	n11 := Simd_v128_load_nc(m, s0+12, 0)
	n12 := Simd_v128_load(m, s4+16, 0)
	n13 := Simd_v128_xor(n12, n11)
	n14 := Simd_v128_and(n13, [2]uint64{p0, p0h})
	n15 := Simd_i8x16_avgr_u(n11, n12)
	n16 := Simd_i8x16_sub(n10, n15)
	n17 := Simd_i8x16_add(n16, n14)
	_ = Simd_v128_store(m, s3+16, 0, n17)
	return
}

//go:noinline
func Simd_p_fx532(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, p0, p0h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s0+-4, 0)
	n2 := Simd_scalar_i32_add(s1, s2)
	n3 := Simd_v128_load(m, n2, 0)
	n4 := Simd_v128_xor(n3, n1)
	n5 := Simd_v128_and(n4, [2]uint64{p0, p0h})
	n6 := Simd_i8x16_avgr_u(n1, n3)
	n7 := Simd_i8x16_sub(n0, n6)
	n8 := Simd_i8x16_add(n7, n5)
	n9 := Simd_scalar_i32_add(s3, s2)
	_ = Simd_v128_store(m, n9, 0, n8)
	return
}

//go:noinline
func Simd_p_fx533(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64) {
	n0 := Simd_v128_load_rng(m, s0+4, 0, -4, 20)
	n1 := Simd_v128_load_nc(m, s0, 0)
	n2 := Simd_v128_load(m, s1, 0)
	n3 := Simd_v128_xor(n2, n1)
	n4 := Simd_v128_and(n3, [2]uint64{p0, p0h})
	n5 := Simd_i8x16_avgr_u(n1, n2)
	n6 := Simd_i8x16_sub(n0, n5)
	n7 := Simd_i8x16_add(n6, n4)
	_ = Simd_v128_store(m, s2, 0, n7)
	return
}

//go:noinline
func Simd_p_fx534(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64) {
	n0 := Simd_v128_load_rng(m, s0+4, 0, -4, 20)
	n1 := Simd_v128_load_nc(m, s0, 0)
	n2 := Simd_v128_load_rng(m, s1+4, 0, -4, 20)
	n3 := Simd_v128_xor(n2, n1)
	n4 := Simd_v128_and(n3, [2]uint64{p0, p0h})
	n5 := Simd_i8x16_avgr_u(n1, n2)
	n6 := Simd_i8x16_sub(n5, n4)
	n7 := Simd_v128_load_nc(m, s1, 0)
	n8 := Simd_i8x16_avgr_u(n6, n7)
	n9 := Simd_i8x16_sub(n0, n8)
	n10 := Simd_v128_xor(n6, n7)
	n11 := Simd_v128_and(n10, [2]uint64{p0, p0h})
	n12 := Simd_i8x16_add(n9, n11)
	_ = Simd_v128_store(m, s2, 0, n12)
	return
}

//go:noinline
func Simd_p_fx535(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1+-4, 0)
	n2 := Simd_i8x16_sub(n0, n1)
	_ = Simd_v128_store(m, s2, 0, n2)
	n4 := Simd_v128_load(m, s0+16, 0)
	n5 := Simd_v128_load(m, s1+12, 0)
	n6 := Simd_i8x16_sub(n4, n5)
	_ = Simd_v128_store(m, s2+16, 0, n6)
	return
}

//go:noinline
func Simd_p_fx536(m *Module, s0 int32, s1 int32, s2 int32, s3 int32) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_scalar_i32_add(s2, -4)
	n3 := Simd_scalar_i32_add(n2, s1)
	n4 := Simd_v128_load(m, n3, 0)
	n5 := Simd_i8x16_sub(n1, n4)
	n6 := Simd_scalar_i32_add(s3, s1)
	_ = Simd_v128_store(m, n6, 0, n5)
	return
}

//go:noinline
func Simd_p_fx537(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_v128_load_rng(m, s0+4, 0, -4, 20)
	n1 := Simd_v128_load_nc(m, s0, 0)
	n2 := Simd_i8x16_shuffle(n1, n1, [2]uint64{p0, p0h})
	n3 := Simd_i8x16_shuffle(n1, n1, [2]uint64{p2, p2h})
	n4 := Simd_v128_load_rng(m, s1+4, 0, -4, 20)
	n5 := Simd_i8x16_shuffle(n4, n4, [2]uint64{p0, p0h})
	n6 := Simd_i8x16_shuffle(n4, n4, [2]uint64{p2, p2h})
	n7 := Simd_v128_load_nc(m, s1, 0)
	n8 := Simd_i8x16_shuffle(n7, n1, [2]uint64{p0, p0h})
	n9 := Simd_i8x16_sub_sat_u(n8, n2)
	n10 := Simd_i8x16_sub_sat_u(n2, n8)
	n11 := Simd_v128_or(n9, n10)
	n12 := Simd_i16x8_shr_u(n11, 8)
	n13 := Simd_v128_and(n11, [2]uint64{p1, p1h})
	n14 := Simd_i16x8_add(n12, n13)
	n15 := Simd_i32x4_shl(n14, 16)
	n16 := Simd_i16x8_add(n14, n15)
	n17 := Simd_i64x2_shl(n16, 32)
	n18 := Simd_i16x8_add(n16, n17)
	n19 := Simd_i64x2_shr_u(n18, 48)
	n20 := Simd_i8x16_shuffle(n7, n1, [2]uint64{p2, p2h})
	n21 := Simd_i8x16_sub_sat_u(n20, n3)
	n22 := Simd_i8x16_sub_sat_u(n3, n20)
	n23 := Simd_v128_or(n21, n22)
	n24 := Simd_i16x8_shr_u(n23, 8)
	n25 := Simd_v128_and(n23, [2]uint64{p1, p1h})
	n26 := Simd_i16x8_add(n24, n25)
	n27 := Simd_i32x4_shl(n26, 16)
	n28 := Simd_i16x8_add(n26, n27)
	n29 := Simd_i64x2_shl(n28, 32)
	n30 := Simd_i16x8_add(n28, n29)
	n31 := Simd_i64x2_shr_u(n30, 48)
	n32 := Simd_i16x8_narrow_i32x4_s(n19, n31)
	n33 := Simd_i8x16_shuffle(n7, n4, [2]uint64{p0, p0h})
	n34 := Simd_i8x16_sub_sat_u(n33, n5)
	n35 := Simd_i8x16_sub_sat_u(n5, n33)
	n36 := Simd_v128_or(n34, n35)
	n37 := Simd_i16x8_shr_u(n36, 8)
	n38 := Simd_v128_and(n36, [2]uint64{p1, p1h})
	n39 := Simd_i16x8_add(n37, n38)
	n40 := Simd_i32x4_shl(n39, 16)
	n41 := Simd_i16x8_add(n39, n40)
	n42 := Simd_i64x2_shl(n41, 32)
	n43 := Simd_i8x16_shuffle(n7, n4, [2]uint64{p2, p2h})
	n44 := Simd_i8x16_sub_sat_u(n43, n6)
	n45 := Simd_i8x16_sub_sat_u(n6, n43)
	n46 := Simd_v128_or(n44, n45)
	n47 := Simd_i16x8_shr_u(n46, 8)
	n48 := Simd_v128_and(n46, [2]uint64{p1, p1h})
	n49 := Simd_i16x8_add(n47, n48)
	n50 := Simd_i32x4_shl(n49, 16)
	n51 := Simd_i16x8_add(n49, n50)
	n52 := Simd_i64x2_shl(n51, 32)
	n53 := Simd_i16x8_add(n41, n42)
	n54 := Simd_i64x2_shr_u(n53, 48)
	n55 := Simd_i16x8_add(n51, n52)
	n56 := Simd_i64x2_shr_u(n55, 48)
	n57 := Simd_i16x8_narrow_i32x4_s(n54, n56)
	n58 := Simd_i32x4_gt_s(n32, n57)
	n59 := Simd_v128_bitselect(n1, n4, n58)
	n60 := Simd_i8x16_sub(n0, n59)
	_ = Simd_v128_store(m, s2, 0, n60)
	return
}

//go:noinline
func Simd_p_fx538(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64) {
	n0 := Simd_v128_load_rng(m, s0+4, 0, -4, 20)
	n1 := Simd_v128_load_rng(m, s1+4, 0, -4, 24)
	n2 := Simd_v128_load_nc(m, s1+8, 0)
	n3 := Simd_v128_xor(n2, n1)
	n4 := Simd_v128_and(n3, [2]uint64{p0, p0h})
	n5 := Simd_i8x16_avgr_u(n1, n2)
	n6 := Simd_i8x16_sub(n5, n4)
	n7 := Simd_v128_load_nc(m, s0, 0)
	n8 := Simd_v128_load_nc(m, s1, 0)
	n9 := Simd_v128_xor(n8, n7)
	n10 := Simd_v128_and(n9, [2]uint64{p0, p0h})
	n11 := Simd_i8x16_avgr_u(n7, n8)
	n12 := Simd_i8x16_sub(n11, n10)
	n13 := Simd_i8x16_avgr_u(n6, n12)
	n14 := Simd_i8x16_sub(n0, n13)
	n15 := Simd_v128_xor(n6, n12)
	n16 := Simd_v128_and(n15, [2]uint64{p0, p0h})
	n17 := Simd_i8x16_add(n14, n16)
	_ = Simd_v128_store(m, s2, 0, n17)
	return
}

//go:noinline
func Simd_p_fx539(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load_rng(m, s1, 0, 0, 20)
	n2 := Simd_v128_load_nc(m, s1+4, 0)
	n3 := Simd_v128_xor(n2, n1)
	n4 := Simd_v128_and(n3, [2]uint64{p0, p0h})
	n5 := Simd_i8x16_avgr_u(n1, n2)
	n6 := Simd_i8x16_sub(n0, n5)
	n7 := Simd_i8x16_add(n6, n4)
	_ = Simd_v128_store(m, s2, 0, n7)
	return
}

//go:noinline
func Simd_p_fx540(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1+-4, 0)
	n2 := Simd_v128_load(m, s1, 0)
	n3 := Simd_v128_xor(n2, n1)
	n4 := Simd_v128_and(n3, [2]uint64{p0, p0h})
	n5 := Simd_i8x16_avgr_u(n1, n2)
	n6 := Simd_i8x16_sub(n0, n5)
	n7 := Simd_i8x16_add(n6, n4)
	_ = Simd_v128_store(m, s2, 0, n7)
	n9 := Simd_v128_load(m, s0+16, 0)
	n10 := Simd_v128_load_rng(m, s1+12, 0, 0, 20)
	n11 := Simd_v128_load_nc(m, s1+16, 0)
	n12 := Simd_v128_xor(n11, n10)
	n13 := Simd_v128_and(n12, [2]uint64{p0, p0h})
	n14 := Simd_i8x16_avgr_u(n10, n11)
	n15 := Simd_i8x16_sub(n9, n14)
	n16 := Simd_i8x16_add(n15, n13)
	_ = Simd_v128_store(m, s2+16, 0, n16)
	return
}

//go:noinline
func Simd_p_fx541(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, p0, p0h uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_v128_load(m, s2+-4, 0)
	n3 := Simd_v128_load(m, s2, 0)
	n4 := Simd_v128_xor(n3, n2)
	n5 := Simd_v128_and(n4, [2]uint64{p0, p0h})
	n6 := Simd_i8x16_avgr_u(n2, n3)
	n7 := Simd_i8x16_sub(n1, n6)
	n8 := Simd_i8x16_add(n7, n5)
	n9 := Simd_scalar_i32_add(s3, s1)
	_ = Simd_v128_store(m, n9, 0, n8)
	return
}

//go:noinline
func Simd_p_fx542(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_v128_load_rng(m, s0+4, 0, -4, 20)
	n1 := Simd_v128_load_rng(m, s1+4, 0, -4, 20)
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_shuffle(n1, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n4 := Simd_v128_load_nc(m, s0, 0)
	n5 := Simd_i8x16_shuffle(n4, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n6 := Simd_i16x8_add(n2, n5)
	n7 := Simd_i16x8_shr_u(n6, 1)
	n8 := Simd_i8x16_shuffle(n4, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n9 := Simd_i16x8_add(n3, n8)
	n10 := Simd_i16x8_shr_u(n9, 1)
	n11 := Simd_v128_load_nc(m, s1, 0)
	n12 := Simd_i8x16_shuffle(n11, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n13 := Simd_i8x16_shuffle(n11, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n14 := Simd_i16x8_sub(n7, n12)
	n15 := Simd_i16x8_lt_s(n7, n12)
	n16 := Simd_i16x8_sub(n14, n15)
	n17 := Simd_i16x8_shr_s(n16, 1)
	n18 := Simd_i16x8_add(n17, n7)
	n19 := Simd_i16x8_sub(n10, n13)
	n20 := Simd_i16x8_lt_s(n10, n13)
	n21 := Simd_i16x8_sub(n19, n20)
	n22 := Simd_i16x8_shr_s(n21, 1)
	n23 := Simd_i16x8_add(n22, n10)
	n24 := Simd_i8x16_narrow_i16x8_u(n18, n23)
	n25 := Simd_i8x16_sub(n0, n24)
	_ = Simd_v128_store(m, s2, 0, n25)
	return
}

//go:noinline
func Simd_p_fx543(m *Module, s0 int32, s1 int32, s2 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_v128_load_rng(m, s0+4, 0, -4, 20)
	n1 := Simd_v128_load_rng(m, s1+4, 0, -4, 20)
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n3 := Simd_i8x16_shuffle(n1, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n4 := Simd_v128_load_nc(m, s0, 0)
	n5 := Simd_i8x16_shuffle(n4, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n6 := Simd_i16x8_add(n2, n5)
	n7 := Simd_i8x16_shuffle(n4, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n8 := Simd_i16x8_add(n3, n7)
	n9 := Simd_v128_load_nc(m, s1, 0)
	n10 := Simd_i8x16_shuffle(n9, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n11 := Simd_i16x8_sub(n6, n10)
	n12 := Simd_i8x16_shuffle(n9, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n13 := Simd_i16x8_sub(n8, n12)
	n14 := Simd_i8x16_narrow_i16x8_u(n11, n13)
	n15 := Simd_i8x16_sub(n0, n14)
	_ = Simd_v128_store(m, s2, 0, n15)
	return
}

//go:noinline
func Simd_p_fx544(m *Module, s0 int32, p0, p0h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i8x16_swizzle(n0, [2]uint64{p0, p0h})
	n2 := Simd_i8x16_sub(n0, n1)
	_ = Simd_v128_store(m, s0, 0, n2)
	n4 := Simd_v128_load(m, s0+16, 0)
	n5 := Simd_i8x16_swizzle(n4, [2]uint64{p0, p0h})
	n6 := Simd_i8x16_sub(n4, n5)
	_ = Simd_v128_store(m, s0+16, 0, n6)
	n8 := Simd_v128_load(m, s0+32, 0)
	n9 := Simd_i8x16_swizzle(n8, [2]uint64{p0, p0h})
	n10 := Simd_i8x16_sub(n8, n9)
	_ = Simd_v128_store(m, s0+32, 0, n10)
	n12 := Simd_v128_load(m, s0+48, 0)
	n13 := Simd_i8x16_swizzle(n12, [2]uint64{p0, p0h})
	n14 := Simd_i8x16_sub(n12, n13)
	_ = Simd_v128_store(m, s0+48, 0, n14)
	return
}

//go:noinline
func Simd_p_fx545(m *Module, s0 int32, p0, p0h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i8x16_swizzle(n0, [2]uint64{p0, p0h})
	n2 := Simd_i8x16_sub(n0, n1)
	_ = Simd_v128_store(m, s0, 0, n2)
	return
}

//go:noinline
func Simd_p_fx546(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 32)
	n1 := Simd_v128_load_nc(m, s0+16, 0)
	n2 := Simd_i8x16_shuffle(n0, n1, [2]uint64{795458214199165184, 1952900979608391952})
	n3 := Simd_i8x16_shuffle(n0, n1, [2]uint64{1084818905551471876, 2242261670960698644})
	n4 := Simd_i32x4_add(n2, n3)
	n5 := Simd_i32x4_mul(n4, [2]uint64{p0, p0h})
	n6 := Simd_i32x4_add(n5, [2]uint64{p1, p1h})
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx547(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1084818905618843912, 1663540288323457296})
	n1 := Simd_i32x4_add([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx548(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{795458214266537220, 1374179596971150604})
	n1 := Simd_i32x4_add([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx549(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i8x16_swizzle(n0, [2]uint64{p0, p0h})
	n2 := Simd_i32x4_extend_low_i16x8_s(n1)
	n3 := Simd_i32x4_mul(n2, [2]uint64{p1, p1h})
	n4 := Simd_i32x4_extend_high_i16x8_s(n1)
	n5 := Simd_i32x4_mul(n4, [2]uint64{p2, p2h})
	n6 := Simd_i8x16_shuffle(n3, n5, [2]uint64{1084816697938281218, 2242259463347507986})
	n7 := Simd_i16x8_sub(n0, n6)
	n8 := Simd_i32x4_shr_u(n7, 16)
	n9 := Simd_i16x8_add(n7, n8)
	return n9[0], n9[1]
}

//go:noinline
func Simd_p_fx550(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_and(n0, [2]uint64{p0, p0h})
	n2 := Simd_i32x4_dot_i16x8_s(n1, [2]uint64{p1, p1h})
	n3 := Simd_i16x8_sub(n0, n2)
	return n3[0], n3[1]
}

//go:noinline
func Simd_p_fx551(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_add([2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n1 := Simd_v128_load32_splat(m, s0, 12)
	n2 := Simd_v128_load32_lane(m, s1+12, 0, 1, n1)
	n3 := Simd_v128_load32_lane(m, s2+12, 0, 2, n2)
	n4 := Simd_v128_load32_lane(m, s3+12, 0, 3, n3)
	n5 := Simd_v128_load32_splat(m, s0, 8)
	n6 := Simd_v128_load32_lane(m, s1+8, 0, 1, n5)
	n7 := Simd_v128_load32_lane(m, s2+8, 0, 2, n6)
	n8 := Simd_v128_load32_lane(m, s3+8, 0, 3, n7)
	n9 := Simd_i32x4_add(n4, n8)
	n10 := Simd_i32x4_mul(n9, [2]uint64{p0, p0h})
	n11 := Simd_i32x4_add(n10, [2]uint64{p1, p1h})
	return n11[0], n11[1], n0[0], n0[1]
}

//go:noinline
func Simd_p_fx552(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1, 0)
	n2 := Simd_i32x4_add(n0, n1)
	_ = Simd_v128_store(m, s2, 0, n2)
	return
}

//go:noinline
func Simd_p_fx553(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1, 0)
	n2 := Simd_i32x4_add(n0, n1)
	_ = Simd_v128_store(m, s0, 0, n2)
	return
}

//go:noinline
func Simd_p_fx554(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) {
	n0 := Simd_v128_load_rng(m, s0+4, 0, -4, 20)
	n1 := Simd_v128_or(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_or(n0, [2]uint64{p2, p2h})
	n3 := Simd_v128_load_nc(m, s0, 0)
	n4 := Simd_v128_and(n3, [2]uint64{p1, p1h})
	n5 := Simd_i32x4_sub(n1, n4)
	n6 := Simd_v128_and(n5, [2]uint64{p1, p1h})
	n7 := Simd_v128_and(n3, [2]uint64{p3, p3h})
	n8 := Simd_i32x4_sub(n2, n7)
	n9 := Simd_v128_and(n8, [2]uint64{p3, p3h})
	n10 := Simd_v128_or(n6, n9)
	_ = Simd_v128_store(m, s1, 0, n10)
	return
}

//go:noinline
func Simd_p_fx555(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load(m, s0, 256)
	_ = Simd_v128_store(m, s1, 3264, n0)
	n2 := Simd_v128_load(m, s0, 272)
	_ = Simd_v128_store(m, s1, 3280, n2)
	return
}

//go:noinline
func Simd_p_fx556(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load(m, s0+304, 0)
	_ = Simd_v128_store(m, s1+48, 0, n0)
	n2 := Simd_v128_load(m, s0+288, 0)
	_ = Simd_v128_store(m, s1+32, 0, n2)
	return
}

//go:noinline
func Simd_p_fx557(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1+-64, 0)
	_ = Simd_v128_store(m, s0, 0, n1)
	n3 := Simd_v128_load(m, s0+48, 0)
	n4 := Simd_v128_load(m, s1+-16, 0)
	_ = Simd_v128_store(m, s0+48, 0, n4)
	n6 := Simd_v128_load(m, s0+32, 0)
	n7 := Simd_v128_load(m, s1+-32, 0)
	_ = Simd_v128_store(m, s0+32, 0, n7)
	n9 := Simd_v128_load(m, s0+16, 0)
	n10 := Simd_v128_load(m, s1+-48, 0)
	_ = Simd_v128_store(m, s0+16, 0, n10)
	_ = Simd_v128_store(m, s2+16, 0, n9)
	_ = Simd_v128_store(m, s2+32, 0, n6)
	_ = Simd_v128_store(m, s2+48, 0, n3)
	_ = Simd_v128_store(m, s2, 0, n0)
	return
}

//go:noinline
func Simd_p_fx558(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	_ = Simd_v128_store(m, s1, 0, n0)
	n2 := Simd_v128_load(m, s0+48, 0)
	_ = Simd_v128_store(m, s1+48, 0, n2)
	n4 := Simd_v128_load(m, s0+32, 0)
	_ = Simd_v128_store(m, s1+32, 0, n4)
	n6 := Simd_v128_load(m, s0+16, 0)
	_ = Simd_v128_store(m, s1+16, 0, n6)
	return
}

//go:noinline
func Simd_p_fx559(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	_ = Simd_v128_store(m, s1, 0, n0)
	n2 := Simd_v128_load(m, s0+48, 0)
	_ = Simd_v128_store(m, s1+48, 0, n2)
	n4 := Simd_v128_load(m, s0+32, 0)
	_ = Simd_v128_store(m, s1+32, 0, n4)
	n6 := Simd_v128_load(m, s0+16, 0)
	_ = Simd_v128_store(m, s2, 0, n6)
	return
}

//go:noinline
func Simd_p_fx560(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1, 0)
	_ = Simd_v128_store(m, s0, 0, n1)
	n3 := Simd_v128_load(m, s0+48, 0)
	n4 := Simd_v128_load(m, s1+48, 0)
	_ = Simd_v128_store(m, s0+48, 0, n4)
	n6 := Simd_v128_load(m, s0+32, 0)
	n7 := Simd_v128_load(m, s1+32, 0)
	_ = Simd_v128_store(m, s0+32, 0, n7)
	n9 := Simd_v128_load(m, s0+16, 0)
	n10 := Simd_v128_load(m, s1+16, 0)
	_ = Simd_v128_store(m, s0+16, 0, n10)
	_ = Simd_v128_store(m, s2+16, 0, n9)
	_ = Simd_v128_store(m, s2+32, 0, n6)
	_ = Simd_v128_store(m, s2+48, 0, n3)
	_ = Simd_v128_store(m, s2, 0, n0)
	return
}

//go:noinline
func Simd_p_fx561(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i32x4_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i32x4_ne(n0, [2]uint64{p2, p2h})
	n3 := Simd_i32x4_sub([2]uint64{p1, p1h}, n2)
	return n1[0], n1[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx562(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0+52, 0, -48, 64)
	n1 := Simd_i32x4_ne(n0, [2]uint64{p2, p2h})
	n2 := Simd_v128_load_nc(m, s0+36, 0)
	n3 := Simd_i32x4_ne(n2, [2]uint64{p2, p2h})
	n4 := Simd_i32x4_add(n0, n2)
	n5 := Simd_v128_load_nc(m, s0+20, 0)
	n6 := Simd_i32x4_ne(n5, [2]uint64{p2, p2h})
	n7 := Simd_i32x4_add(n4, n5)
	n8 := Simd_v128_load_nc(m, s0+4, 0)
	n9 := Simd_i32x4_ne(n8, [2]uint64{p2, p2h})
	n10 := Simd_i32x4_add(n7, n8)
	n11 := Simd_i32x4_add(n10, [2]uint64{p0, p0h})
	n12 := Simd_i32x4_sub([2]uint64{p1, p1h}, n9)
	n13 := Simd_i32x4_sub(n12, n6)
	n14 := Simd_i32x4_sub(n13, n3)
	n15 := Simd_i32x4_sub(n14, n1)
	return n11[0], n11[1], n15[0], n15[1]
}

//go:noinline
func Simd_p_fx563(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0+1076, 0, -48, 64)
	n1 := Simd_i32x4_ne(n0, [2]uint64{p2, p2h})
	n2 := Simd_v128_load_nc(m, s0+1060, 0)
	n3 := Simd_i32x4_ne(n2, [2]uint64{p2, p2h})
	n4 := Simd_i32x4_add(n0, n2)
	n5 := Simd_v128_load_nc(m, s0+1044, 0)
	n6 := Simd_i32x4_ne(n5, [2]uint64{p2, p2h})
	n7 := Simd_i32x4_add(n4, n5)
	n8 := Simd_v128_load_nc(m, s0+1028, 0)
	n9 := Simd_i32x4_ne(n8, [2]uint64{p2, p2h})
	n10 := Simd_i32x4_add(n7, n8)
	n11 := Simd_i32x4_add(n10, [2]uint64{p0, p0h})
	n12 := Simd_i32x4_sub([2]uint64{p1, p1h}, n9)
	n13 := Simd_i32x4_sub(n12, n6)
	n14 := Simd_i32x4_sub(n13, n3)
	n15 := Simd_i32x4_sub(n14, n1)
	return n11[0], n11[1], n15[0], n15[1]
}

//go:noinline
func Simd_p_fx564(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0+2100, 0, -48, 64)
	n1 := Simd_i32x4_ne(n0, [2]uint64{p2, p2h})
	n2 := Simd_v128_load_nc(m, s0+2084, 0)
	n3 := Simd_i32x4_ne(n2, [2]uint64{p2, p2h})
	n4 := Simd_i32x4_add(n0, n2)
	n5 := Simd_v128_load_nc(m, s0+2068, 0)
	n6 := Simd_i32x4_ne(n5, [2]uint64{p2, p2h})
	n7 := Simd_i32x4_add(n4, n5)
	n8 := Simd_v128_load_nc(m, s0+2052, 0)
	n9 := Simd_i32x4_ne(n8, [2]uint64{p2, p2h})
	n10 := Simd_i32x4_add(n7, n8)
	n11 := Simd_i32x4_add(n10, [2]uint64{p0, p0h})
	n12 := Simd_i32x4_sub([2]uint64{p1, p1h}, n9)
	n13 := Simd_i32x4_sub(n12, n6)
	n14 := Simd_i32x4_sub(n13, n3)
	n15 := Simd_i32x4_sub(n14, n1)
	return n11[0], n11[1], n15[0], n15[1]
}

//go:noinline
func Simd_p_fx565(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_ne([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_load_rng(m, s0, 3096, 3096, 112)
	n2 := Simd_i32x4_ne(n1, [2]uint64{p1, p1h})
	n3 := Simd_v128_load_nc(m, s0, 3112)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx566(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_ne([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_load_nc(m, s0, 3128)
	n2 := Simd_i32x4_ne(n1, [2]uint64{p1, p1h})
	n3 := Simd_v128_load_nc(m, s0, 3192)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx567(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_ne([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_load_nc(m, s0, 3176)
	n2 := Simd_i32x4_ne(n1, [2]uint64{p1, p1h})
	n3 := Simd_v128_load_nc(m, s0, 3160)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx568(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_ne([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_load_nc(m, s0, 3144)
	n2 := Simd_i32x4_ne(n1, [2]uint64{p1, p1h})
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx569(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_splat(s0)
	n1 := Simd_i32x4_add(n0, [2]uint64{p0, p0h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx570(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_add(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_add([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	_ = Simd_v128_store(m, s0, 0, n0)
	return n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx571(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_shr_u([2]uint64{p0, p0h}, 8)
	n1 := Simd_v128_and(n0, [2]uint64{p1, p1h})
	n2 := Simd_i32x4_shl(n0, 16)
	n3 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n4 := Simd_i32x4_add(n1, n3)
	n5 := Simd_i32x4_add(n4, n2)
	n6 := Simd_v128_and(n5, [2]uint64{p2, p2h})
	n7 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p3, p3h})
	n8 := Simd_v128_or(n6, n7)
	return n8[0], n8[1]
}

//go:noinline
func Simd_p_fx572(m *Module, s0 int32, s1 int32) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1+188, 0)
	_ = Simd_v128_store(m, s0, 0, n1)
	return n0[0], n0[1]
}

//go:noinline
func Simd_p_fx573(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	_ = Simd_v128_store(m, s1, 32, n0)
	return
}

//go:noinline
func Simd_p_fx574(m *Module, s0 int32, s1 int32) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1+8, 0)
	_ = Simd_v128_store(m, s0, 0, n1)
	return n0[0], n0[1]
}

//go:noinline
func Simd_p_fx575(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	_ = Simd_v128_store(m, s1, 2112, n0)
	return
}

//go:noinline
func Simd_p_fx576(m *Module, s0 int32, s1 int32) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1+1096, 0)
	_ = Simd_v128_store(m, s0, 0, n1)
	return n0[0], n0[1]
}

//go:noinline
func Simd_p_fx577(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_v128_load32_splat(m, s0, 0)
	n1 := Simd_v128_load32_lane(m, s0+12, 0, 1, n0)
	n2 := Simd_v128_load32_lane(m, s0+24, 0, 2, n1)
	n3 := Simd_v128_load32_lane(m, s0+36, 0, 3, n2)
	n4 := Simd_i32x4_max_s([2]uint64{p0, p0h}, n3)
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx578(m *Module, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{1084818905618843912, 216736831629295872})
	n1 := Simd_i32x4_max_s([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx579(m *Module, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{216736831696667908, 216736831629295872})
	n1 := Simd_i32x4_max_s([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx580(m *Module, s0 int32, s1 int32) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1+16, 0)
	_ = Simd_v128_store(m, s0, 0, n1)
	return n0[0], n0[1]
}

//go:noinline
func Simd_p_fx581(m *Module, s0 int32) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i16x8_extend_low_i8x16_u(n0)
	n2 := Simd_i32x4_extend_low_i16x8_u(n1)
	return n0[0], n0[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx582(m *Module, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{117835012, 0})
	n1 := Simd_i16x8_extend_low_i8x16_u(n0)
	n2 := Simd_i32x4_extend_low_i16x8_u(n1)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx583(m *Module, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{185207048, 0})
	n1 := Simd_i16x8_extend_low_i8x16_u(n0)
	n2 := Simd_i32x4_extend_low_i16x8_u(n1)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx584(m *Module, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{252579084, 0})
	n1 := Simd_i16x8_extend_low_i8x16_u(n0)
	n2 := Simd_i32x4_extend_low_i16x8_u(n1)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx585(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_ne([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n2 := Simd_v128_and(n1, [2]uint64{p4, p4h})
	n3 := Simd_i32x4_add([2]uint64{p0, p0h}, n2)
	n4 := Simd_i8x16_ne([2]uint64{p5, p5h}, [2]uint64{p2, p2h})
	n5 := Simd_i8x16_shuffle(n4, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n6 := Simd_v128_and(n5, [2]uint64{p4, p4h})
	n7 := Simd_i32x4_add(n3, n6)
	n8 := Simd_i8x16_ne([2]uint64{p6, p6h}, [2]uint64{p2, p2h})
	n9 := Simd_i8x16_shuffle(n8, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n10 := Simd_v128_and(n9, [2]uint64{p4, p4h})
	n11 := Simd_i32x4_add(n7, n10)
	n12 := Simd_v128_load32_zero(m, s0+12, 0)
	n13 := Simd_i8x16_ne(n12, [2]uint64{p2, p2h})
	n14 := Simd_i8x16_shuffle(n13, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n15 := Simd_v128_and(n14, [2]uint64{p4, p4h})
	n16 := Simd_i32x4_add(n11, n15)
	return n16[0], n16[1]
}

//go:noinline
func Simd_p_fx586(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load(m, s0, 200)
	_ = Simd_v128_store(m, s1, 148, n0)
	n2 := Simd_v128_load(m, s0, 216)
	_ = Simd_v128_store(m, s1, 164, n2)
	return
}

//go:noinline
func Simd_p_fx587(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_v128_load32_splat(m, s0+-96, 0)
	n1 := Simd_v128_load32_lane(m, s0+-64, 0, 1, n0)
	n2 := Simd_v128_load32_lane(m, s0+-32, 0, 2, n1)
	n3 := Simd_v128_load32_lane(m, s0, 0, 3, n2)
	n4 := Simd_i32x4_add(n3, [2]uint64{p0, p0h})
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx588(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	_ = Simd_v128_store(m, s2, 16, n1)
	return
}

//go:noinline
func Simd_p_fx589(m *Module, s0 int32, s1 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	_ = Simd_v128_store(m, s1, 92, n0)
	return
}

//go:noinline
func Simd_p_fx590(m *Module, s0 int32) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 200)
	n1 := Simd_v128_load32_lane(m, s0+136, 0, 1, n0)
	n2 := Simd_v128_load(m, s0, 196)
	n3 := Simd_v128_load32_lane(m, s0+132, 0, 1, n2)
	n4 := Simd_v128_load(m, s0, 204)
	n5 := Simd_v128_load32_lane(m, s0+140, 0, 1, n4)
	n6 := Simd_v128_load(m, s0, 208)
	n7 := Simd_v128_load32_lane(m, s0+144, 0, 1, n6)
	return n1[0], n1[1], n3[0], n3[1], n5[0], n5[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx591(m *Module, s0 int32) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 212)
	n1 := Simd_v128_load32_lane(m, s0+148, 0, 1, n0)
	n2 := Simd_v128_load(m, s0, 216)
	n3 := Simd_v128_load32_lane(m, s0+152, 0, 1, n2)
	n4 := Simd_v128_load(m, s0, 220)
	n5 := Simd_v128_load32_lane(m, s0+156, 0, 1, n4)
	n6 := Simd_v128_load(m, s0, 224)
	n7 := Simd_v128_load32_lane(m, s0+160, 0, 1, n6)
	return n1[0], n1[1], n3[0], n3[1], n5[0], n5[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx592(m *Module, s0 int32) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 228)
	n1 := Simd_v128_load32_lane(m, s0+164, 0, 1, n0)
	n2 := Simd_v128_load(m, s0, 232)
	n3 := Simd_v128_load32_lane(m, s0+168, 0, 1, n2)
	n4 := Simd_v128_load(m, s0, 236)
	n5 := Simd_v128_load32_lane(m, s0+172, 0, 1, n4)
	n6 := Simd_v128_load(m, s0, 240)
	n7 := Simd_v128_load32_lane(m, s0+176, 0, 1, n6)
	return n1[0], n1[1], n3[0], n3[1], n5[0], n5[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx593(m *Module, s0 int32) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load32_splat(m, s0, 244)
	n1 := Simd_v128_load32_lane(m, s0+180, 0, 1, n0)
	n2 := Simd_v128_load32_splat(m, s0, 248)
	n3 := Simd_v128_load32_lane(m, s0+184, 0, 1, n2)
	n4 := Simd_v128_load32_splat(m, s0, 252)
	n5 := Simd_v128_load32_lane(m, s0+188, 0, 1, n4)
	return n1[0], n1[1], n3[0], n3[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx594(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_gt_s([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx595(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_gt_s([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx596(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_gt_s([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx597(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_gt_s([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx598(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_gt_s([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx599(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_gt_s([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx600(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_gt_s([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx601(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_gt_s([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx602(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_gt_s([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx603(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_gt_s([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx604(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_gt_s([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx605(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_gt_s([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx606(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_gt_s([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, n0)
	n2 := Simd_i32x4_gt_s([2]uint64{p4, p4h}, [2]uint64{p1, p1h})
	n3 := Simd_v128_bitselect([2]uint64{p3, p3h}, [2]uint64{p1, p1h}, n2)
	n4 := Simd_i32x4_gt_s([2]uint64{p5, p5h}, [2]uint64{p1, p1h})
	n5 := Simd_i32x4_sub(n3, n4)
	n6 := Simd_i32x4_add(n1, n5)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx607(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_or(n0, [2]uint64{p1, p1h})
	n2 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, n0, [2]uint64{1374179596971150604, 1952900979675763988})
	n3 := Simd_v128_and(n2, [2]uint64{p2, p2h})
	n4 := Simd_i32x4_sub(n1, n3)
	n5 := Simd_v128_and(n4, [2]uint64{p3, p3h})
	n6 := Simd_v128_and(n2, [2]uint64{p1, p1h})
	n7 := Simd_i32x4_sub(n0, n6)
	n8 := Simd_i32x4_shr_u(n7, 8)
	n9 := Simd_v128_and(n8, [2]uint64{p3, p3h})
	n10 := Simd_i32x4_shr_u(n4, 16)
	n11 := Simd_v128_and(n10, [2]uint64{p3, p3h})
	return n0[0], n0[1], n5[0], n5[1], n9[0], n9[1], n11[0], n11[1]
}

//go:noinline
func Simd_p_fx608(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_lt_u([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n2 := Simd_v128_bitselect([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, n1)
	n3 := Simd_i32x4_eq([2]uint64{p3, p3h}, [2]uint64{p0, p0h})
	n4 := Simd_i8x16_shuffle(n3, [2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n5 := Simd_v128_bitselect([2]uint64{p0, p0h}, n2, n4)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx609(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_lt_u([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n2 := Simd_v128_bitselect([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, n1)
	n3 := Simd_i32x4_eq([2]uint64{p3, p3h}, [2]uint64{p0, p0h})
	n4 := Simd_i8x16_shuffle(n3, [2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n5 := Simd_v128_bitselect([2]uint64{p0, p0h}, n2, n4)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx610(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_lt_u([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n2 := Simd_v128_bitselect([2]uint64{p1, p1h}, [2]uint64{p2, p2h}, n1)
	n3 := Simd_i32x4_eq([2]uint64{p3, p3h}, [2]uint64{p0, p0h})
	n4 := Simd_i8x16_shuffle(n3, [2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n5 := Simd_v128_bitselect([2]uint64{p0, p0h}, n2, n4)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx611(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{770, 0})
	n1 := Simd_v128_or([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx612(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1, 0})
	n1 := Simd_v128_or([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx613(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i32x4_add(n0, [2]uint64{p0, p0h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx614(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i32x4_ne(n0, [2]uint64{p1, p1h})
	n2 := Simd_i32x4_sub([2]uint64{p0, p0h}, n1)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx615(m *Module, s0 int32, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_v128_load32_zero(m, s0, 0)
	n1 := Simd_i16x8_extend_low_i8x16_u(n0)
	n2 := Simd_i32x4_extend_low_i16x8_u(n1)
	n3 := Simd_i32x4_max_s([2]uint64{p0, p0h}, n2)
	return n3[0], n3[1]
}

//go:noinline
func Simd_p_fx616(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) {
	n0 := Simd_f64x2_convert_low_i32x4_u([2]uint64{p1, p1h})
	n1 := Simd_f64x2_mul([2]uint64{p0, p0h}, n0)
	n2 := Simd_f64x2_div(n1, [2]uint64{p2, p2h})
	n3 := Simd_f64x2_add(n2, [2]uint64{p3, p3h})
	_ = Simd_v128_store(m, s0, 0, n3)
	return
}

//go:noinline
func Simd_p_fx617(m *Module, s0 int32) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i8x16_shuffle(n0, n0, [2]uint64{252579084, 0})
	n2 := Simd_i16x8_extend_low_i8x16_u(n1)
	n3 := Simd_i32x4_extend_low_i16x8_u(n2)
	n4 := Simd_i8x16_shuffle(n0, n0, [2]uint64{185207048, 0})
	n5 := Simd_i16x8_extend_low_i8x16_u(n4)
	n6 := Simd_i32x4_extend_low_i16x8_u(n5)
	n7 := Simd_i8x16_shuffle(n0, n0, [2]uint64{117835012, 0})
	n8 := Simd_i16x8_extend_low_i8x16_u(n7)
	n9 := Simd_i32x4_extend_low_i16x8_u(n8)
	n10 := Simd_i16x8_extend_low_i8x16_u(n0)
	n11 := Simd_i32x4_extend_low_i16x8_u(n10)
	return n3[0], n3[1], n6[0], n6[1], n9[0], n9[1], n11[0], n11[1]
}

//go:noinline
func Simd_p_fx618(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1, 0)
	n2 := Simd_i16x8_sub(n0, n1)
	n3 := Simd_v128_load(m, s2, 0)
	n4 := Simd_i16x8_add(n2, n3)
	_ = Simd_v128_store(m, s2, 0, n4)
	n6 := Simd_v128_load(m, s0+16, 0)
	n7 := Simd_v128_load(m, s1+16, 0)
	n8 := Simd_i16x8_sub(n6, n7)
	n9 := Simd_v128_load(m, s2+16, 0)
	n10 := Simd_i16x8_add(n8, n9)
	_ = Simd_v128_store(m, s2+16, 0, n10)
	return
}

//go:noinline
func Simd_p_fx619(m *Module, s0 int32, s1 int32, s2 int32, s3 int32) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_scalar_i32_add(s2, s1)
	n3 := Simd_v128_load(m, n2, 0)
	n4 := Simd_i16x8_sub(n1, n3)
	n5 := Simd_v128_load(m, s3, 0)
	n6 := Simd_i16x8_add(n4, n5)
	_ = Simd_v128_store(m, s3, 0, n6)
	return
}

//go:noinline
func Simd_p_fx620(m *Module, s0 int32, s1 int32, s2 int32) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_v128_load(m, s1, 0)
	n2 := Simd_i16x8_sub(n0, n1)
	n3 := Simd_v128_load(m, s2, 0)
	n4 := Simd_i16x8_add(n2, n3)
	_ = Simd_v128_store(m, s2, 0, n4)
	return
}

//go:noinline
func Simd_p_fx621(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_splat(s3)
	n1 := Simd_v128_load(m, s0, 0)
	n2 := Simd_v128_load(m, s1, 0)
	n3 := Simd_v128_load(m, s2, 0)
	n4 := Simd_i16x8_sub(n2, n3)
	n5 := Simd_i16x8_shr_s(n4, 15)
	n6 := Simd_v128_or(n5, [2]uint64{p1, p1h})
	n7 := Simd_i16x8_add(n1, n4)
	n8 := Simd_i16x8_min_s(n7, n0)
	n9 := Simd_i16x8_max_s(n8, [2]uint64{p0, p0h})
	n10 := Simd_i32x4_dot_i16x8_s(n4, n6)
	n11 := Simd_i32x4_add(n10, [2]uint64{p2, p2h})
	_ = Simd_v128_store(m, s0, 0, n9)
	return n11[0], n11[1]
}

//go:noinline
func Simd_p_fx622(m *Module, s0 int32, s1 int32, s2 int32) (uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load32_zero(m, s0, 0)
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	n2 := Simd_v128_load32_zero(m, s1, 0)
	n3 := Simd_i32x4_extend_low_i16x8_u(n2)
	n4 := Simd_i32x4_sub(n1, n3)
	n5 := Simd_v128_load32_zero(m, s2, 0)
	n6 := Simd_i32x4_extend_low_i16x8_u(n5)
	n7 := Simd_i32x4_add(n4, n6)
	return n4[0], n4[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx623(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_splat(s0)
	n1 := Simd_i32x4_min_s([2]uint64{p1, p1h}, n0)
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n3 := Simd_i32x4_lt_s([2]uint64{p1, p1h}, [2]uint64{p0, p0h})
	n4 := Simd_i8x16_shuffle(n3, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	n5 := Simd_v128_bitselect([2]uint64{p0, p0h}, n2, n4)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx624(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_abs([2]uint64{p1, p1h})
	n1 := Simd_i64x2_extend_low_i32x4_u(n0)
	n2 := Simd_i64x2_add([2]uint64{p0, p0h}, n1)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx625(m *Module, p0, p0h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{1084818905618843912, 506097522914230528})
	n1 := Simd_i64x2_add([2]uint64{p0, p0h}, n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx626(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_shr_s(n0, 16)
	n2 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, [2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n3 := Simd_i32x4_shr_s(n2, 16)
	n4 := Simd_i32x4_add(n3, n1)
	return n1[0], n1[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx627(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_shr_s(n0, 16)
	n2 := Simd_i8x16_shuffle([2]uint64{p2, p2h}, [2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n3 := Simd_i32x4_shr_s(n2, 16)
	n4 := Simd_i32x4_add(n3, n1)
	n5 := Simd_i32x4_shl(n4, 1)
	n6 := Simd_i32x4_add([2]uint64{p3, p3h}, n4)
	n7 := Simd_i32x4_add(n6, [2]uint64{p4, p4h})
	n8 := Simd_i32x4_add(n5, n7)
	n9 := Simd_i32x4_shr_s(n8, 3)
	n10 := Simd_v128_and(n9, [2]uint64{p5, p5h})
	n11 := Simd_v128_xor(n9, [2]uint64{p5, p5h})
	n12 := Simd_i32x4_shr_s(n11, 1)
	n13 := Simd_i32x4_add(n10, n12)
	n14 := Simd_i32x4_shl([2]uint64{p3, p3h}, 1)
	n15 := Simd_i32x4_add(n7, n14)
	n16 := Simd_i32x4_shr_s(n15, 3)
	n17 := Simd_v128_and(n16, n1)
	n18 := Simd_v128_xor(n16, n1)
	n19 := Simd_i32x4_shr_s(n18, 1)
	n20 := Simd_i32x4_add(n17, n19)
	n21 := Simd_v128_load(m, s0, 0)
	return n20[0], n20[1], n13[0], n13[1], n21[0], n21[1]
}

//go:noinline
func Simd_p_fx628(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1374179596769034496, 1663540288121341188})
	n1 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1952900979473647880, 2242261670825954572})
	n2 := Simd_i16x8_narrow_i32x4_s(n0, n1)
	n3 := Simd_i16x8_add(n2, [2]uint64{p2, p2h})
	n4 := Simd_i16x8_min_s(n3, [2]uint64{p3, p3h})
	n5 := Simd_i16x8_max_s(n4, [2]uint64{p4, p4h})
	_ = Simd_v128_store(m, s0, 0, n5)
	return
}

//go:noinline
func Simd_p_fx629(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_splat(s0)
	n1 := Simd_i32x4_add(n0, [2]uint64{p0, p0h})
	n2 := Simd_i32x4_add(n0, [2]uint64{p1, p1h})
	n3 := Simd_i32x4_splat(s1)
	return n3[0], n3[1], n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx630(m *Module, s0 int32, s1 int32) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i32x4_extend_low_i16x8_s(n0)
	n2 := Simd_v128_load(m, s1, 0)
	n3 := Simd_i32x4_extend_low_i16x8_s(n2)
	n4 := Simd_i32x4_add(n1, n3)
	return n0[0], n0[1], n2[0], n2[1], n3[0], n3[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx631(m *Module, s0 int32, s1 int32, s2 int32) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i32x4_extend_low_i16x8_s(n0)
	n2 := Simd_scalar_i32_add(s1, s2)
	n3 := Simd_v128_load(m, n2, 0)
	n4 := Simd_i32x4_extend_low_i16x8_s(n3)
	n5 := Simd_i32x4_add(n4, n1)
	return n0[0], n0[1], n1[0], n1[1], n3[0], n3[1], n5[0], n5[1]
}

//go:noinline
func Simd_p_fx632(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_add(n0, [2]uint64{p2, p2h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx633(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_shl([2]uint64{p0, p0h}, 1)
	n1 := Simd_i32x4_shl([2]uint64{p1, p1h}, 3)
	n2 := Simd_i32x4_add(n0, n1)
	n3 := Simd_i32x4_add(n2, [2]uint64{p2, p2h})
	n4 := Simd_i32x4_shr_s(n3, 4)
	n5 := Simd_i32x4_extend_low_i16x8_u([2]uint64{p3, p3h})
	n6 := Simd_i32x4_add(n4, n5)
	n7 := Simd_i8x16_shuffle([2]uint64{p4, p4h}, [2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n8 := Simd_i32x4_extend_low_i16x8_s(n7)
	n9 := Simd_i8x16_shuffle([2]uint64{p5, p5h}, [2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n10 := Simd_i32x4_extend_low_i16x8_s(n9)
	n11 := Simd_i32x4_add(n10, n8)
	return n6[0], n6[1], n8[0], n8[1], n11[0], n11[1]
}

//go:noinline
func Simd_p_fx634(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i32x4_extend_low_i16x8_s(n0)
	n2 := Simd_i8x16_shuffle([2]uint64{p3, p3h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n3 := Simd_i32x4_extend_low_i16x8_s(n2)
	n4 := Simd_i32x4_add(n3, n1)
	n5 := Simd_i32x4_add(n4, [2]uint64{p1, p1h})
	n6 := Simd_i32x4_add(n5, [2]uint64{p4, p4h})
	n7 := Simd_i32x4_shl([2]uint64{p5, p5h}, 1)
	return n1[0], n1[1], n4[0], n4[1], n6[0], n6[1], n7[0], n7[1]
}

//go:noinline
func Simd_p_fx635(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_shl([2]uint64{p0, p0h}, 1)
	n1 := Simd_i32x4_shl([2]uint64{p1, p1h}, 3)
	n2 := Simd_i32x4_add(n0, n1)
	n3 := Simd_i32x4_add(n2, [2]uint64{p2, p2h})
	n4 := Simd_i32x4_shr_s(n3, 4)
	n5 := Simd_i32x4_extend_low_i16x8_u([2]uint64{p3, p3h})
	n6 := Simd_i32x4_add(n4, n5)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx636(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_min_s([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_v128_and(n0, [2]uint64{p3, p3h})
	n2 := Simd_i32x4_min_s([2]uint64{p4, p4h}, [2]uint64{p2, p2h})
	n3 := Simd_v128_and(n2, [2]uint64{p3, p3h})
	n4 := Simd_i16x8_narrow_i32x4_u(n1, n3)
	n5 := Simd_i32x4_lt_s([2]uint64{p1, p1h}, [2]uint64{p0, p0h})
	n6 := Simd_i32x4_lt_s([2]uint64{p4, p4h}, [2]uint64{p0, p0h})
	n7 := Simd_i8x16_shuffle(n5, n6, [2]uint64{p5, p5h})
	n8 := Simd_v128_bitselect([2]uint64{p0, p0h}, n4, n7)
	return n8[0], n8[1]
}

//go:noinline
func Simd_p_fx637(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_shl([2]uint64{p0, p0h}, 3)
	n1 := Simd_i32x4_shl([2]uint64{p1, p1h}, 1)
	n2 := Simd_i32x4_add(n0, n1)
	n3 := Simd_i32x4_add(n2, [2]uint64{p2, p2h})
	n4 := Simd_i32x4_shr_s(n3, 4)
	n5 := Simd_i32x4_extend_low_i16x8_u([2]uint64{p3, p3h})
	n6 := Simd_i32x4_add(n4, n5)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx638(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 32)
	n1 := Simd_v128_load_nc(m, s0+16, 0)
	n2 := Simd_v128_load_rng(m, s1, 0, 0, 18)
	n3 := Simd_v128_load_rng(m, s2+2, 0, -2, 18)
	n4 := Simd_i16x8_add(n2, n3)
	n5 := Simd_i16x8_shl(n4, 1)
	n6 := Simd_v128_load_nc(m, s1+2, 0)
	n7 := Simd_v128_load_nc(m, s2, 0)
	n8 := Simd_i16x8_add(n6, n7)
	n9 := Simd_i16x8_shl(n8, 1)
	n10 := Simd_i16x8_add(n4, n8)
	n11 := Simd_i16x8_add(n10, [2]uint64{p0, p0h})
	n12 := Simd_i16x8_add(n11, n5)
	n13 := Simd_i16x8_shr_s(n12, 3)
	n14 := Simd_i16x8_add(n13, n7)
	n15 := Simd_i16x8_shr_s(n14, 1)
	n16 := Simd_i16x8_add(n9, n11)
	n17 := Simd_i16x8_shr_s(n16, 3)
	n18 := Simd_i16x8_add(n17, n3)
	n19 := Simd_i16x8_shr_s(n18, 1)
	n20 := Simd_i8x16_shuffle(n15, n19, [2]uint64{1952885526417115400, 2242246217769422092})
	n21 := Simd_i16x8_add(n1, n20)
	n22 := Simd_i16x8_min_s(n21, [2]uint64{p1, p1h})
	n23 := Simd_i16x8_max_s(n22, [2]uint64{p2, p2h})
	n24 := Simd_i8x16_shuffle(n15, n19, [2]uint64{1374164143712502016, 1663524835064808708})
	n25 := Simd_i16x8_add(n0, n24)
	n26 := Simd_i16x8_min_s(n25, [2]uint64{p1, p1h})
	n27 := Simd_i16x8_max_s(n26, [2]uint64{p2, p2h})
	_ = Simd_v128_store(m, s3+16, 0, n23)
	_ = Simd_v128_store(m, s3, 0, n27)
	return
}

//go:noinline
func Simd_p_fx639(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_splat(s0)
	n1 := Simd_v128_or(n0, [2]uint64{p0, p0h})
	n2 := Simd_v128_or(n0, [2]uint64{p1, p1h})
	n3 := Simd_i32x4_splat(s1)
	return n3[0], n3[1], n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx640(m *Module, s0 int32, s1 int32) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_rng(m, s0, 0, 0, 18)
	n1 := Simd_i32x4_extend_low_i16x8_s(n0)
	n2 := Simd_v128_load_rng(m, s1+2, 0, -2, 18)
	n3 := Simd_i32x4_extend_low_i16x8_s(n2)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx641(m *Module, s0 int32, s1 int32) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load_nc(m, s0, 0)
	n1 := Simd_i32x4_extend_low_i16x8_s(n0)
	n2 := Simd_v128_load_nc(m, s1+2, 0)
	n3 := Simd_i32x4_extend_low_i16x8_s(n2)
	return n0[0], n0[1], n1[0], n1[1], n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx642(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_add(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_add([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n3 := Simd_i32x4_mul(n2, [2]uint64{p5, p5h})
	n4 := Simd_i32x4_add(n1, n3)
	n5 := Simd_i32x4_add(n4, [2]uint64{p6, p6h})
	n6 := Simd_i32x4_shr_s(n5, 4)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx643(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_extend_low_i16x8_s(n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx644(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i32x4_extend_low_i16x8_s(n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx645(m *Module, s0 int32, p0, p0h uint64) {
	n0 := Simd_i16x8_shl([2]uint64{p0, p0h}, 2)
	_ = Simd_v128_store(m, s0, 0, n0)
	return
}

//go:noinline
func Simd_p_fx646(m *Module, s0 int32, s1 int32, p0, p0h uint64) {
	n0 := Simd_i16x8_shl([2]uint64{p0, p0h}, 2)
	n1 := Simd_scalar_i32_add(s0, s1)
	_ = Simd_v128_store(m, n1, 0, n0)
	return
}

//go:noinline
func Simd_p_fx647(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_lt_u([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_and(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_lt_u([2]uint64{p3, p3h}, [2]uint64{p1, p1h})
	n3 := Simd_v128_and(n2, [2]uint64{p2, p2h})
	n4 := Simd_i16x8_narrow_i32x4_u(n1, n3)
	return n4[0], n4[1]
}

//go:noinline
func Simd_p_fx648(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, s5 int32, s6 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n2 := Simd_i32x4_shl(n1, s1)
	n3 := Simd_v128_and(n2, [2]uint64{p0, p0h})
	n4 := Simd_i8x16_shuffle(n0, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n5 := Simd_i32x4_shl(n4, s1)
	n6 := Simd_v128_and(n5, [2]uint64{p0, p0h})
	n7 := Simd_i16x8_narrow_i32x4_u(n3, n6)
	_ = Simd_v128_store(m, s2, 0, n7)
	n9 := Simd_v128_load(m, s3, 0)
	n10 := Simd_i8x16_shuffle(n9, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n11 := Simd_i32x4_shl(n10, s1)
	n12 := Simd_v128_and(n11, [2]uint64{p0, p0h})
	n13 := Simd_i8x16_shuffle(n9, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n14 := Simd_i32x4_shl(n13, s1)
	n15 := Simd_v128_and(n14, [2]uint64{p0, p0h})
	n16 := Simd_i16x8_narrow_i32x4_u(n12, n15)
	_ = Simd_v128_store(m, s4, 0, n16)
	n18 := Simd_v128_load(m, s5, 0)
	n19 := Simd_i8x16_shuffle(n18, [2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n20 := Simd_i32x4_shl(n19, s1)
	n21 := Simd_v128_and(n20, [2]uint64{p0, p0h})
	n22 := Simd_i8x16_shuffle(n18, [2]uint64{p0, p0h}, [2]uint64{p2, p2h})
	n23 := Simd_i32x4_shl(n22, s1)
	n24 := Simd_v128_and(n23, [2]uint64{p0, p0h})
	n25 := Simd_i16x8_narrow_i32x4_u(n21, n24)
	n26 := Simd_scalar_i32_add(s2, s6)
	_ = Simd_v128_store(m, n26, 0, n25)
	return
}

//go:noinline
func Simd_p_fx649(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, s5 int32, s6 int32, p0, p0h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	n2 := Simd_i32x4_shr_u(n1, s1)
	n3 := Simd_i8x16_shuffle(n0, n0, [2]uint64{p0, p0h})
	n4 := Simd_i32x4_extend_low_i16x8_u(n3)
	n5 := Simd_i32x4_shr_u(n4, s1)
	n6 := Simd_i16x8_narrow_i32x4_u(n2, n5)
	_ = Simd_v128_store(m, s2, 0, n6)
	n8 := Simd_v128_load(m, s3, 0)
	n9 := Simd_i32x4_extend_low_i16x8_u(n8)
	n10 := Simd_i32x4_shr_u(n9, s1)
	n11 := Simd_i8x16_shuffle(n8, n8, [2]uint64{p0, p0h})
	n12 := Simd_i32x4_extend_low_i16x8_u(n11)
	n13 := Simd_i32x4_shr_u(n12, s1)
	n14 := Simd_i16x8_narrow_i32x4_u(n10, n13)
	_ = Simd_v128_store(m, s4, 0, n14)
	n16 := Simd_v128_load(m, s5, 0)
	n17 := Simd_i32x4_extend_low_i16x8_u(n16)
	n18 := Simd_i32x4_shr_u(n17, s1)
	n19 := Simd_i8x16_shuffle(n16, n16, [2]uint64{p0, p0h})
	n20 := Simd_i32x4_extend_low_i16x8_u(n19)
	n21 := Simd_i32x4_shr_u(n20, s1)
	n22 := Simd_i16x8_narrow_i32x4_u(n18, n21)
	n23 := Simd_scalar_i32_add(s2, s6)
	_ = Simd_v128_store(m, n23, 0, n22)
	return
}

//go:noinline
func Simd_p_fx650(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_i32x4_extend_low_i16x8_u(n1)
	n3 := Simd_i32x4_mul(n2, [2]uint64{p0, p0h})
	n4 := Simd_i8x16_shuffle(n1, n1, [2]uint64{p4, p4h})
	n5 := Simd_i32x4_extend_low_i16x8_u(n4)
	n6 := Simd_i32x4_mul(n5, [2]uint64{p0, p0h})
	n7 := Simd_scalar_i32_add(s2, s1)
	n8 := Simd_v128_load(m, n7, 0)
	n9 := Simd_i32x4_extend_low_i16x8_u(n8)
	n10 := Simd_i32x4_mul(n9, [2]uint64{p1, p1h})
	n11 := Simd_i32x4_add(n3, n10)
	n12 := Simd_i8x16_shuffle(n8, n1, [2]uint64{p4, p4h})
	n13 := Simd_i32x4_extend_low_i16x8_u(n12)
	n14 := Simd_i32x4_mul(n13, [2]uint64{p1, p1h})
	n15 := Simd_i32x4_add(n6, n14)
	n16 := Simd_scalar_i32_add(s3, s1)
	n17 := Simd_v128_load(m, n16, 0)
	n18 := Simd_i32x4_extend_low_i16x8_u(n17)
	n19 := Simd_i32x4_mul(n18, [2]uint64{p2, p2h})
	n20 := Simd_i32x4_add(n11, n19)
	n21 := Simd_i32x4_add(n20, [2]uint64{p3, p3h})
	n22 := Simd_i32x4_shr_u(n21, 16)
	n23 := Simd_i8x16_shuffle(n17, n1, [2]uint64{p4, p4h})
	n24 := Simd_i32x4_extend_low_i16x8_u(n23)
	n25 := Simd_i32x4_mul(n24, [2]uint64{p2, p2h})
	n26 := Simd_i32x4_add(n15, n25)
	n27 := Simd_i32x4_add(n26, [2]uint64{p3, p3h})
	n28 := Simd_i32x4_shr_u(n27, 16)
	n29 := Simd_i16x8_narrow_i32x4_u(n22, n28)
	n30 := Simd_scalar_i32_add(s4, s1)
	_ = Simd_v128_store(m, n30, 0, n29)
	return
}

//go:noinline
func Simd_p_fx651(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	n2 := Simd_i32x4_mul(n1, [2]uint64{p0, p0h})
	n3 := Simd_i8x16_shuffle(n0, n0, [2]uint64{p4, p4h})
	n4 := Simd_i32x4_extend_low_i16x8_u(n3)
	n5 := Simd_i32x4_mul(n4, [2]uint64{p0, p0h})
	n6 := Simd_v128_load(m, s1, 0)
	n7 := Simd_i32x4_extend_low_i16x8_u(n6)
	n8 := Simd_i32x4_mul(n7, [2]uint64{p1, p1h})
	n9 := Simd_i32x4_add(n2, n8)
	n10 := Simd_i8x16_shuffle(n6, n0, [2]uint64{p4, p4h})
	n11 := Simd_i32x4_extend_low_i16x8_u(n10)
	n12 := Simd_i32x4_mul(n11, [2]uint64{p1, p1h})
	n13 := Simd_i32x4_add(n5, n12)
	n14 := Simd_v128_load(m, s2, 0)
	n15 := Simd_i32x4_extend_low_i16x8_u(n14)
	n16 := Simd_i32x4_mul(n15, [2]uint64{p2, p2h})
	n17 := Simd_i32x4_add(n9, n16)
	n18 := Simd_i32x4_add(n17, [2]uint64{p3, p3h})
	n19 := Simd_i32x4_shr_u(n18, 16)
	n20 := Simd_i8x16_shuffle(n14, n0, [2]uint64{p4, p4h})
	n21 := Simd_i32x4_extend_low_i16x8_u(n20)
	n22 := Simd_i32x4_mul(n21, [2]uint64{p2, p2h})
	n23 := Simd_i32x4_add(n13, n22)
	n24 := Simd_i32x4_add(n23, [2]uint64{p3, p3h})
	n25 := Simd_i32x4_shr_u(n24, 16)
	n26 := Simd_i16x8_narrow_i32x4_u(n19, n25)
	_ = Simd_v128_store(m, s3, 0, n26)
	return
}

//go:noinline
func Simd_p_fx652(m *Module, s0 int32) (uint64, uint64) {
	n0 := Simd_v128_load32_zero(m, s0, 0)
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx653(m *Module, s0 int32, s1 int32, s2 int32, s3 int32, s4 int32, s5 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_splat(s0)
	n1 := Simd_i32x4_extend_low_i16x8_s([2]uint64{p0, p0h})
	n2 := Simd_i32x4_add(n1, [2]uint64{p1, p1h})
	n3 := Simd_i32x4_mul(n0, n2)
	n4 := Simd_i32x4_splat(s1)
	n5 := Simd_i32x4_add(n3, n4)
	n6 := Simd_i32x4_splat(s2)
	n7 := Simd_i32x4_extend_low_i16x8_s([2]uint64{p2, p2h})
	n8 := Simd_i32x4_add(n7, [2]uint64{p1, p1h})
	n9 := Simd_i32x4_mul(n6, n8)
	n10 := Simd_i32x4_add(n5, n9)
	n11 := Simd_i32x4_splat(s3)
	n12 := Simd_i32x4_extend_low_i16x8_s([2]uint64{p3, p3h})
	n13 := Simd_i32x4_add(n12, [2]uint64{p1, p1h})
	n14 := Simd_i32x4_mul(n11, n13)
	n15 := Simd_i32x4_add(n10, n14)
	n16 := Simd_i32x4_splat(s4)
	n17 := Simd_i32x4_add(n15, n16)
	n18 := Simd_i32x4_shr_s(n17, s5)
	n19 := Simd_i8x16_shuffle(n18, n18, [2]uint64{72058693633704192, 72058693566333184})
	return n18[0], n18[1], n19[0], n19[1]
}

//go:noinline
func Simd_p_fx654(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p0, p0h}, [2]uint64{1024, 0})
	n1 := Simd_i16x8_gt_s([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n2 := Simd_i8x16_shuffle(n1, [2]uint64{p1, p1h}, [2]uint64{p3, p3h})
	n3 := Simd_i16x8_lt_u([2]uint64{p1, p1h}, [2]uint64{p4, p4h})
	n4 := Simd_i8x16_shuffle(n3, [2]uint64{p1, p1h}, [2]uint64{p3, p3h})
	n5 := Simd_v128_bitselect(n0, n2, n4)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx655(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64) (uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_i32x4_add([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_scalar_i32_add(s0, s1)
	n2 := Simd_v128_load(m, n1, 0)
	n3 := Simd_i32x4_extend_low_i16x8_u(n2)
	return n2[0], n2[1], n3[0], n3[1], n0[0], n0[1]
}

//go:noinline
func Simd_p_fx656(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_extend_low_i16x8_s([2]uint64{p1, p1h})
	n1 := Simd_i32x4_add(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_mul(n1, [2]uint64{p3, p3h})
	n3 := Simd_i32x4_add([2]uint64{p0, p0h}, n2)
	n4 := Simd_i32x4_extend_low_i16x8_s([2]uint64{p4, p4h})
	n5 := Simd_i32x4_add(n4, [2]uint64{p2, p2h})
	n6 := Simd_i32x4_mul(n5, [2]uint64{p5, p5h})
	n7 := Simd_i32x4_add(n3, n6)
	return n7[0], n7[1]
}

//go:noinline
func Simd_p_fx657(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_extend_low_i16x8_s([2]uint64{p0, p0h})
	n1 := Simd_i32x4_add(n0, [2]uint64{p1, p1h})
	n2 := Simd_i32x4_mul(n1, [2]uint64{p2, p2h})
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx658(m *Module, p0, p0h uint64, p1, p1h uint64) (uint64, uint64) {
	n0 := Simd_i8x16_shuffle([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, [2]uint64{1084818905618843912, 72058693566333184})
	n1 := Simd_i32x4_extend_low_i16x8_u(n0)
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx659(m *Module, s0 int32, s1 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64) {
	n0 := Simd_i32x4_shl([2]uint64{p1, p1h}, 16)
	n1 := Simd_i32x4_shr_s(n0, 16)
	n2 := Simd_i32x4_gt_s(n1, [2]uint64{p0, p0h})
	n3 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, n2)
	n4 := Simd_v128_and(n3, [2]uint64{p2, p2h})
	n5 := Simd_i32x4_shl([2]uint64{p3, p3h}, 16)
	n6 := Simd_i32x4_shr_s(n5, 16)
	n7 := Simd_i32x4_gt_s(n6, [2]uint64{p0, p0h})
	n8 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p3, p3h}, n7)
	n9 := Simd_v128_and(n8, [2]uint64{p2, p2h})
	n10 := Simd_i16x8_narrow_i32x4_u(n4, n9)
	n11 := Simd_v128_and([2]uint64{p1, p1h}, [2]uint64{p5, p5h})
	n12 := Simd_i32x4_eq(n11, [2]uint64{p4, p4h})
	n13 := Simd_v128_and([2]uint64{p3, p3h}, [2]uint64{p5, p5h})
	n14 := Simd_i32x4_eq(n13, [2]uint64{p4, p4h})
	n15 := Simd_i8x16_shuffle(n12, n14, [2]uint64{940136352262127872, 2097579117671354640})
	n16 := Simd_v128_bitselect(n10, [2]uint64{p4, p4h}, n15)
	n17 := Simd_scalar_i32_add(s0, s1)
	_ = Simd_v128_store(m, n17, 0, n16)
	return
}

//go:noinline
func Simd_p_fx660(m *Module, s0 int32, s1 int32) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_v128_load(m, s0, 0)
	n1 := Simd_i32x4_extend_low_i16x8_s(n0)
	n2 := Simd_scalar_i32_add(s0, s1)
	n3 := Simd_v128_load(m, n2, 0)
	n4 := Simd_i32x4_extend_low_i16x8_s(n3)
	return n0[0], n0[1], n1[0], n1[1], n3[0], n3[1], n4[0], n4[1]
}

//go:noinline
func Simd_p_fx661(m *Module, s0 int32, s1 int32) (uint64, uint64, uint64, uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_i32x4_extend_low_i16x8_s(n1)
	return n1[0], n1[1], n2[0], n2[1]
}

//go:noinline
func Simd_p_fx662(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i32x4_add(n0, [2]uint64{p2, p2h})
	n2 := Simd_i32x4_mul([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n3 := Simd_i32x4_add(n1, n2)
	n4 := Simd_i32x4_mul([2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n5 := Simd_i32x4_add(n3, n4)
	return n5[0], n5[1]
}

//go:noinline
func Simd_p_fx663(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64) (uint64, uint64) {
	n0 := Simd_v128_and([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_v128_and([2]uint64{p2, p2h}, [2]uint64{p1, p1h})
	n2 := Simd_i16x8_narrow_i32x4_u(n0, n1)
	return n2[0], n2[1]
}

//go:noinline
func Simd_p_fx664(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_gt_s([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx665(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64) (uint64, uint64) {
	n0 := Simd_i16x8_lt_u([2]uint64{p0, p0h}, [2]uint64{p1, p1h})
	n1 := Simd_i8x16_shuffle(n0, [2]uint64{p2, p2h}, [2]uint64{p3, p3h})
	return n1[0], n1[1]
}

//go:noinline
func Simd_p_fx666(m *Module, s0 int32, s1 int32, s2 int32) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_v128_load(m, n0, 0)
	n2 := Simd_i32x4_extend_low_i16x8_s(n1)
	n3 := Simd_scalar_i32_add(s0, s2)
	n4 := Simd_scalar_i32_add(n3, s1)
	n5 := Simd_v128_load(m, n4, 0)
	n6 := Simd_i32x4_extend_low_i16x8_s(n5)
	return n1[0], n1[1], n2[0], n2[1], n5[0], n5[1], n6[0], n6[1]
}

//go:noinline
func Simd_p_fx667(m *Module, s0 int32, s1 int32, s2 int32) (uint64, uint64, uint64, uint64) {
	n0 := Simd_scalar_i32_add(s0, s1)
	n1 := Simd_scalar_i32_add(n0, s2)
	n2 := Simd_v128_load(m, n1, 0)
	n3 := Simd_i32x4_extend_low_i16x8_s(n2)
	return n2[0], n2[1], n3[0], n3[1]
}

//go:noinline
func Simd_p_fx668(m *Module, s0 int32, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_mul([2]uint64{p1, p1h}, [2]uint64{p2, p2h})
	n1 := Simd_i32x4_add([2]uint64{p0, p0h}, n0)
	n2 := Simd_i32x4_mul([2]uint64{p3, p3h}, [2]uint64{p4, p4h})
	n3 := Simd_i32x4_add(n1, n2)
	n4 := Simd_i32x4_mul([2]uint64{p5, p5h}, [2]uint64{p6, p6h})
	n5 := Simd_i32x4_add(n3, n4)
	n6 := Simd_i32x4_shr_s(n5, s0)
	return n6[0], n6[1]
}

//go:noinline
func Simd_p_fx669(m *Module, p0, p0h uint64, p1, p1h uint64, p2, p2h uint64, p3, p3h uint64, p4, p4h uint64, p5, p5h uint64, p6, p6h uint64) (uint64, uint64) {
	n0 := Simd_i32x4_shl([2]uint64{p1, p1h}, 16)
	n1 := Simd_i32x4_shr_s(n0, 16)
	n2 := Simd_i32x4_gt_s(n1, [2]uint64{p0, p0h})
	n3 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p1, p1h}, n2)
	n4 := Simd_v128_and(n3, [2]uint64{p2, p2h})
	n5 := Simd_i32x4_shl([2]uint64{p3, p3h}, 16)
	n6 := Simd_i32x4_shr_s(n5, 16)
	n7 := Simd_i32x4_gt_s(n6, [2]uint64{p0, p0h})
	n8 := Simd_v128_bitselect([2]uint64{p0, p0h}, [2]uint64{p3, p3h}, n7)
	n9 := Simd_v128_and(n8, [2]uint64{p2, p2h})
	n10 := Simd_i16x8_narrow_i32x4_u(n4, n9)
	n11 := Simd_v128_and([2]uint64{p1, p1h}, [2]uint64{p5, p5h})
	n12 := Simd_i32x4_eq(n11, [2]uint64{p4, p4h})
	n13 := Simd_v128_and([2]uint64{p3, p3h}, [2]uint64{p5, p5h})
	n14 := Simd_i32x4_eq(n13, [2]uint64{p4, p4h})
	n15 := Simd_i8x16_shuffle(n12, n14, [2]uint64{p6, p6h})
	n16 := Simd_v128_bitselect(n10, [2]uint64{p4, p4h}, n15)
	return n16[0], n16[1]
}

var spinAgents int32
var spinOversubscribed uint32

type ThreadPool struct {
	nextTID atomic.Int32
	wg      sync.WaitGroup

	parkMu sync.Mutex
	parked map[uint64][]chan struct{}
}

// wake releases up to count waiters on ea and reports how many it woke.
func (p *ThreadPool) wake(ea uint64, count int32) int32 {
	p.parkMu.Lock()
	defer p.parkMu.Unlock()
	waiters := p.parked[ea]
	n := int32(len(waiters))
	if count >= 0 && count < n {
		n = count
	}
	for _, ch := range waiters[:n] {
		close(ch)
	}
	if int(n) == len(waiters) {
		delete(p.parked, ea)
	} else {
		p.parked[ea] = waiters[n:]
	}
	return n
}

// SaveGlobals returns the module's mutable globals, in a form that can be handed back
// to RestoreGlobals. It is how a snapshot of an instance captures the state that does not
// live in linear memory.
func SaveGlobals(m *Module) []uint64 {
	g := make([]uint64, 1)
	g[0] = uint64(uint32(m.G0))
	return g
}

// RestoreGlobals puts a snapshot's globals back. A snapshot from a different module (or a
// different build of the same one) has a different global count; rather than
// index out of bounds, take what fits and leave the rest at their declared
// initializers.
func RestoreGlobals(m *Module, g []uint64) {
	if len(g) != 1 {
		return
	}
	m.G0 = int32(uint32(g[0]))
}
