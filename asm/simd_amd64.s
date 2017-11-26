
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

// func VecMulChewxy(a, b, out []float64)
TEXT ·VecMulChewxy(SB), NOSPLIT, $0
	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI // use destination index register for this
	MOVQ out_base+48(FP),DX

	MOVQ a_len+8(FP), AX  // len(a) into AX
	//MOVQ b_len+32(FP), BX // len(b) into BX

	// check if they are the same length
	//CMPQ AX, BX
	//JNE  panic  // jump to panic if not the same length. TODO: return bloody errors

	// check if there are at least 8 elements
	SUBQ $8, AX
	JL   remainder

loop:
	// a[0]
	MOVAPS (SI), X0
	MOVAPS (DI), X1
	MULPS  X0, X1
	MOVAPS X1, (DX)

	MOVAPS 16(SI), X2
	MOVAPS 16(DI), X3
	MULPS  X2, X3
	MOVAPS X3, 16(DX)

	MOVAPS 32(SI), X4
	MOVAPS 32(DI), X5
	MULPS  X4, X5
	MOVAPS X5, 32(DX)

	MOVAPS 48(SI), X6
	MOVAPS 48(DI), X7
	MULPS  X6, X7
	MOVAPS X7, 48(DX)

	// update pointers. 4 registers, 2 elements at once, each element is 8 bytes
	ADDQ $64, SI
	ADDQ $64, DI
	ADDQ $64, DX

	// len(a) is now 4*2 elements less
	SUBQ $8, AX
	JGE  loop

remainder:
	ADDQ $8, AX
	JE   done

remainderloop:
	MOVSS (SI), X0
	MOVSS (DI), X1
	MULSD X0, X1
	MOVSD X1, (SI)

	// update pointer to the top of the data
	ADDQ $8, SI
	ADDQ $8, DI

	DECQ AX
	JNE  remainderloop

done:
	RET

panic:
	CALL runtime·panicindex(SB)
	RET
