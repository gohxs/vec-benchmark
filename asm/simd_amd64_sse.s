// +build !noasm,!appengine

#include "textflag.h"

#define NFLOATS $16

// func VecMulf32x4(a, b, out []float32)
TEXT ·VecMulf32x4(SB), $0
	MOVQ    a_base+0(FP), SI  // SI = &a
	MOVQ    b_base+24(FP), DX  // DX = &b
	MOVQ    out_base+48(FP), DI // DI = &out
	MOVQ    out_len+56(FP), CX  // CX = len(out)

	// Smaller size for CX
	CMPQ    a_len+8(FP), CX   // CX = max( len(out), len(a), len(b) )
	CMOVQLE a_len+8(FP), CX  
	CMPQ    b_len+32(FP), CX
	CMOVQLE b_len+32(FP), CX

	CMPQ CX,$0
	JE done

alignment: // Align destination at least?
	MOVQ    DI, BX		// Check on DI
	ANDQ    $15, BX
	JZ      aligned

	MOVSS (SI), X0    // X0 = s[i]
	MULSS (DX), X0    // X0 *= t[i]
	MOVSS  X0, (DI)   // dst[i] = X0

	ADDQ $4, SI
	ADDQ $4, DX
	ADDQ $4, DI
	DECQ  CX                // --CX
	JZ    done              // if CX == 0 { return }
	JMP alignment           // until aligned


aligned:
	SUBQ NFLOATS, CX                  // take 16 floats 4sse * 4unroll
	JL remainder                      // if less than 0

loop:													      // Loop unrolled 4x   do {
	// # MEM TO REG ptr increment 
	MOVUPS   (SI), X0
	MOVUPS 16(SI), X1
	MOVUPS 32(SI), X2
	MOVUPS 48(SI), X3

	MULPS    (DX), X0        // X0 /= y[i:i+1]
	MULPS  16(DX), X1
	MULPS  32(DX), X2
	MULPS  48(DX), X3

	MOVUPS X0,   (DI)        // dst[i:i+1] = X0
	MOVUPS X1, 16(DI)
	MOVUPS X2, 32(DI)
	MOVUPS X3, 48(DI)

	// this is faster than do a single add 
	// a single reg and offseting in MOV ptrs
	ADDQ NFLOATS*4, SI
	ADDQ NFLOATS*4, DI
	ADDQ NFLOATS*4, DX
	
	SUBQ NFLOATS, CX // Take 16 floats
	JGE loop

remainder:                // Reset loop registers
	ADDQ NFLOATS,CX                   // Add back since its negative
	JE done

remainderloop:                     // do { // Last couple of things
	MOVSS (SI), X0
	MULSS (DX), X0
	MOVSS X0, (DI)

	ADDQ $4, SI
	ADDQ $4, DI
	ADDQ $4, DX

	LOOP  remainderloop              // } while --CX > 0

done:
	RET
