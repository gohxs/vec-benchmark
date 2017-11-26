package cgo

// #cgo CFLAGS: -mavx -O2
//
// void VecMulf32x8(int sz, float const *vec1,float const *vec2, float *out);
// void VecMulf32x4(int sz, float const *vec1,float const *vec2, float *out);
import "C"
import "unsafe"

// VecMulf32x8 AVX multiply using Intel intrinsics
func VecMulf32x8(vec1, vec2, out []float32) {
	C.VecMulf32x8(
		C.int(len(vec1)),
		(*C.float)(unsafe.Pointer(&vec1[0])),
		(*C.float)(unsafe.Pointer(&vec2[0])),
		(*C.float)(unsafe.Pointer(&out[0])),
	)
}

// VecMulf32x4 SSE intel extension
func VecMulf32x4(vec1, vec2, out []float32) {
	C.VecMulf32x4(
		C.int(len(vec1)),
		(*C.float)(unsafe.Pointer(&vec1[0])),
		(*C.float)(unsafe.Pointer(&vec2[0])),
		(*C.float)(unsafe.Pointer(&out[0])),
	)
}
