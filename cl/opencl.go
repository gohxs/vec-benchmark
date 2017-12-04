package cl

// #cgo LDFLAGS: -lOpenCL
//
// int platformIdCount();
// void VecMulf32(int sz, float const *vec1, float const *vec2, float *out);
import "C"
import "unsafe"

func PlatformIDCount() int {
	return int(C.platformIdCount())
}

// VecMulf32x8 AVX multiply using Intel intrinsics
func VecMulf32(vec1, vec2, out []float32) {
	if len(vec1) == 0 {
		return
	}
	C.VecMulf32(
		C.int(len(vec1)),
		(*C.float)(unsafe.Pointer(&vec1[0])),
		(*C.float)(unsafe.Pointer(&vec2[0])),
		(*C.float)(unsafe.Pointer(&out[0])),
	)
}
