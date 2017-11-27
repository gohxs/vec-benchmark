
#include "textflag.h"

// Based on Chewxy vec64f
// func VecMul(a, b, out []float64)
TEXT ·VecMulf32x4(SB), $0-72

	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DX
	MOVQ out_base+48(FP),DI   // Destination

	MOVQ a_len+8(FP), AX    // len(a) into AX
	MOVQ b_len+32(FP), BX   // len(b) into BX
	MOVQ out_len+56(FP), CX // len(out) into DX


	CMPQ AX, BX   // Check if a,b are same lenght
	JNE  panic  
	CMPQ AX, CX
	JG   panic    // if output is smaller than inputs 
	
	SUBQ $16, AX         // n floats per loop
	JL   remainder

loop:
	// a[0]
	MOVUPS (SI), X0
	MOVUPS (DX), X1
	MULPS  X0, X1
	MOVUPS X1, (DI) 
	
	MOVUPS 16(SI), X0    // Next 16 bytes (each float32 is 4bytes) * 4 floats 16) - 4 float32
	MOVUPS 16(DX), X1
	MULPS  X0, X1
	MOVUPS X1, 16(DI)

	MOVUPS 32(SI), X4
	MOVUPS 32(DX), X5
	MULPS  X4, X5
	MOVUPS X5, 32(DI)

	MOVUPS 48(SI), X6
	MOVUPS 48(DX), X7
	MULPS  X6, X7
	MOVUPS X7, 48(DI)

	ADDQ $64, SI         // increment sizeof(float32)4 * n
	ADDQ $64, DI
	ADDQ $64, DX

	SUBQ $16, AX         // Count down n floats
	JGE  loop            // Repeat

remainder:
	ADDQ $16, AX         // Re add n elems
	JE   done            // if is 0 go to end

remainderloop:         // 1 by 1
RET //temp
	MOVSS (SI), X0
	MOVSS (DX), X1
	MULSS X0, X1
	MOVSS X1, (DI)

	// update pointer to the top of the data
	ADDQ $4, SI
	ADDQ $4, DI
	ADDQ $4, DX

	DECQ AX
	JNE  remainderloop

done:
	RET

panic:
	CALL runtime·panicindex(SB)
	RET
