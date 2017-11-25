
#include "textflag.h"

// func Mulf32x4(a, b []float32, out float32) float32
TEXT ·Mulf32x4(SB),4,$0-72
	MOVQ a+0(FP), R15 
	MOVQ b+24(FP), R14 
	MOVQ out+48(FP),R13  // Grab Slice?

	MOVOU (R15), X15
	MOVOU (R14), X14
	MULPS   X15, X14
	MOVOU X14, (R13)
RET


// Broken
// func MulVec( a, b []float32, out *[]float32)
TEXT ·VecMul(SB), NOSPLIT, $0-72

	MOVQ a+0(FP), R15						// Grab slice data pointer of A
	MOVQ a_len+8(FP), AX				// Vector size
	MOVQ b+24(FP), R14					// Grab slice data pointer of B
	MOVQ out+48(FP), R13        // Grab output slice data pointer


	MOVQ $16, DX								// Size index
	IMULQ DX, R8								// bytes to quads
	SUBQ $16, R8
loop:

	MOVOU (R15), X15      // Move &A to SSE registers
	MOVOU (R14), X14      // Move &B to SSE Registers
	MULPS  X15 , X14      // Vec multiply
	MOVOU  X14 ,(R13)     // Copy To R13 address


	ADDQ DX, R15         // Increment A pointer
	ADDQ DX, R14         // Increment B pointer
	ADDQ DX, R13         // Increment Output pointer

	SUBQ $1,AX           // Next
	JGE loop
done:
	MOVQ R13, ret+48(SP)
	RET


