package cgo

// #cgo CFLAGS: -mavx -O2
//
// void VecMul256(int sz, float const *vec1,float const *vec2, float *out);
// void VecMul128(int sz, float const *vec1,float const *vec2, float *out);
// void Mul128(int sz, float const *vec1, float const *vec2, float *out);
import "C"
import "unsafe"

// VecMulf32x8 AVX multiply using Intel intrinsics
func VecMulf32x8(vec1, vec2, out []float32) {
	C.VecMul256(
		C.int(len(vec1)),
		(*C.float)(unsafe.Pointer(&vec1[0])),
		(*C.float)(unsafe.Pointer(&vec2[0])),
		(*C.float)(unsafe.Pointer(&out[0])),
	)
}

// VecMulf32x4 SSE intel extension
func VecMulf32x4(vec1, vec2, out []float32) {
	C.VecMul128(
		C.int(len(vec1)),
		(*C.float)(unsafe.Pointer(&vec1[0])),
		(*C.float)(unsafe.Pointer(&vec2[0])),
		(*C.float)(unsafe.Pointer(&out[0])),
	)
}
