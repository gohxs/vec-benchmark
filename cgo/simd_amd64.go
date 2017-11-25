package cgo

// #cgo CFLAGS: -mavx -O2
//
// void VecMul256(int sz, float const *vec1,float const *vec2, float *out);
// void VecMul128(int sz, float const *vec1,float const *vec2, float *out);
// void Mul128(int sz, float const *vec1, float const *vec2, float *out);
import "C"
import "unsafe"

// VecMulx8 AVX multiply using Intel intrinsics
func VecMulx8(vec1, vec2, out []float32) {
	C.VecMul256(
		C.int(len(vec1)),
		(*C.float)(unsafe.Pointer(&vec1[0])),
		(*C.float)(unsafe.Pointer(&vec2[0])),
		(*C.float)(unsafe.Pointer(&out[0])),
	)
}

// VecMulx4 SSE intel extension
func VecMulx4(vec1, vec2, out []float32) {
	C.VecMul128(
		C.int(len(vec1)),
		(*C.float)(unsafe.Pointer(&vec1[0])),
		(*C.float)(unsafe.Pointer(&vec2[0])),
		(*C.float)(unsafe.Pointer(&out[0])),
	)
}

// VecMulx4 SSE intel extension
func Mulx4(vec1, vec2, out []float32) {
	C.Mul128(
		C.int(len(vec1)),
		(*C.float)(unsafe.Pointer(&vec1[0])),
		(*C.float)(unsafe.Pointer(&vec2[0])),
		(*C.float)(unsafe.Pointer(&out[0])),
	)
}
