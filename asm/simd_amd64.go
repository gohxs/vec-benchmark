package asm

// func Mulf32x4(vec1 []float32, vec2 []float32) []float32
/*
TEXT ·Mulf32x4(SB),4,$24-48
	MOVQ         $0, ret0+32(FP)
  MOVQ         $0, ret0+40(FP)
block0:
  MOVUPS       vec1+0(FP), X14    # Move thing to X14
  MOVUPS       vec2+0(FP), X15    # Move thing to X15
  MULPS        X14, X15           # Multiply X14, X15
  MOVUPS       X14, ret0+32(FP)   # move to ret?
  RET
*/

//func Mulf32x4([]float32, []float32) []float32

//Mulf32x4 Multiplies 4 float32 elements from each a,b into out
func Mulf32x4(a, b []float32, out []float32)

// Thing to call makeslice
/*LEAQ    type·float32(SB), AX
MOVQ    AX, (SP)
MOVQ    R8, 8(SP)
MOVQ    R8, 16(SP)
PCDATA  $0, $0
CALL    runtime·makeslice(SB)
MOVQ    32(SP), R13
//MOVQ    24(SP), CX
//MOVQ    40(SP), DX
*/

//VecMul multiplies each element into out
//from assembler
func VecMul(a, b []float32, out []float32)

// Fully assembler
//func MulVec(a, b []float32) []float32
