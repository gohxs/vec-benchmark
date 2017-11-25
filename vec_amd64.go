package vec

import (
	"github.com/gohxs/vec-benchmark/asm"
	"github.com/gohxs/vec-benchmark/cgo"
)

// Funcs to test

// MulASMx4g SSE goassembler multiply 4 at a time elements from vec1 and vec2 to out
// interpolated by go
// arrays must be 4 floats aligned and same size
func MulASMx4g(vec1, vec2, out []float32) {
	max := len(vec1) - 4
	var e int
	for i := 0; i <= max; i += 4 {
		e = i + 4
		asm.Mulf32x4(vec1[i:e], vec2[i:e], out[i:])
	}
}

// MulASMx4 goassembler vector multiplication entirely in assembler (broken)
// interpolated by asm
// using sse instructions
func MulASMx4(vec1, vec2, out []float32) {
	asm.VecMul(vec1, vec2, out)
}

// MulCGOx8 vector multiplication using intrinsict AVX from cgo
// interpolated by cgo
func MulCGOx8(vec1, vec2, out []float32) {
	cgo.VecMulx8(vec1, vec2, out)
}

// MulCGOx4 vector multiplication using intrinsict SSE from cgo
// interpolated by cgo
func MulCGOx4(vec1, vec2, out []float32) {
	cgo.VecMulx4(vec1, vec2, out)
}

// MulCGOx4g multiply 4 elements only
// interpolated by go
func MulCGOx4g(vec1, vec2, out []float32) {
	max := len(vec1) - 4
	var e int
	for i := 0; i <= max; i += 4 {
		e = i + 4
		cgo.Mulx4(vec1[i:e], vec2[i:e], out[i:])
	}
}
