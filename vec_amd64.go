package vec

import (
	"github.com/gohxs/vec-benchmark/asm"
	"github.com/gohxs/vec-benchmark/cgo"
)

// Funcs to test

// MulASMSSEx4gi SSE goassembler multiply 4 at a time elements from vec1 and vec2 to out
// interpolated by go
// arrays must be 4 floats aligned and same size
func MulASMSSEx4gi(vec1, vec2, out []float32) {
	max := len(vec1) - 4
	var e int
	for i := 0; i <= max; i += 4 {
		e = i + 4
		asm.Mulf32x4(vec1[i:e], vec2[i:e], out[i:])
	}
}

// MulASMSSEx4 goassembler vector multiplication entirely in assembler (broken)
// interpolated by asm
// using sse instructions
var MulASMSSEx4 = asm.VecMul

//func MulASMSSEx4(vec1, vec2, out []float32) {
//	asm.VecMul(vec1, vec2, out)
//}

//MulASMChewxy goassembler https://github.com/chewxy/vecf64/blob/master/asm_vecMul_sse.s vecf64 implementation
var MulASMChewxy = asm.VecMulChewxy

// MulCGOXVAx8 vector multiplication using intrinsict AVX from cgo
// interpolated by cgo
var MulCGOXVAx8 = cgo.VecMulx8

// MulCGOSSEx4 vector multiplication using intrinsict SSE from cgo
// interpolated by cgo
var MulCGOSSEx4 = cgo.VecMulx4

// MulCGOSSEx4gi multiply 4 elements only
// interpolated by go
func MulCGOSSEx4gi(vec1, vec2, out []float32) {
	max := len(vec1) - 4
	var e int
	for i := 0; i <= max; i += 4 {
		e = i + 4
		cgo.Mulx4(vec1[i:e], vec2[i:e], out[i:])
	}
}
