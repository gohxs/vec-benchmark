package vec

import (
	"github.com/gohxs/vec-benchmark/asm"
	"github.com/gohxs/vec-benchmark/cgo"
)

// Funcs to test

// MulASMSSEx4 goassembler vector multiplication entirely in assembler (broken)
// interpolated by asm
// using sse instructions

//MulASMf32x4sse goassembler https://github.com/chewxy/vecf64/blob/master/asm_vecMul_sse.s vecf64 implementation
var MulASMf32x4sse = asm.VecMulf32x4

// MulCGOf32x8xva vector multiplication using intrinsict AVX from cgo
// interpolated by cgo
var MulCGOf32x8xva = cgo.VecMulf32x8

// MulCGOf32x4sse vector multiplication using intrinsict SSE from cgo
// interpolated by cgo
var MulCGOf32x4sse = cgo.VecMulf32x4
