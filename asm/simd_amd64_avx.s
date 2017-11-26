#include "textflag.h"

// Based on Chewxy vec64f
// func VecMulf32x8(a, b, out []float64)
TEXT ·VecMulf32x8(SB), $0-72

	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DX
	MOVQ out_base+48(FP),DI   // Destination

	MOVQ a_len+8(FP), AX    // len(a) into AX
	MOVQ b_len+32(FP), BX   // len(b) into BX
	MOVQ out_len+56(FP), CX // len(out) into DX

	// check if they are the same length
	CMPQ AX, BX
	JNE  panic  // jump to panic if not the same length. TODO: return bloody errors
	CMPQ AX, CX
	JG   panic  // jump to panic if not the same length. TODO: return bloody errors
	

	// check if there are at least 8 elements
	SUBQ $16, AX
	JL   remainder

loop:
	// a[0]

	WORD $0xC5FC2806          //vmovaps ymm0,yword [rsi]
	WORD $0xC5FC280A          //vmovaps ymm1,yword [rdx]
	WORD $0xC5FC59C1          //vmulps ymm0,ymm0,ymm1
  WORD $0xC5FC2907          //vmovaps yword [rdi],ymm0


	WORD $0xC5FC2846; BYTE $0x20    //  vmovaps ymm0,yword [rsi+0x20]
	WORD $0xC5FC284A; BYTE $0x20    //  vmovaps ymm1,yword [rdx+0x20]
	WORD $0xC5FC59C1                //  vmulps ymm0,ymm0,ymm1
	WORD $0xC5FC2947; BYTE $0x20    //  vmovaps yword [rdi+0x20],ymm0




	ADDQ $32, SI         // increment 8 iterations 4 * 16
	ADDQ $32, DI
	ADDQ $32, DX

	SUBQ $16, AX         // Count down 1*8 floats
	JGE  loop            // Repeat

remainder:
	ADDQ $16, AX
	JE   done

remainderloop:

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
