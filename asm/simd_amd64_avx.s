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
	SUBQ $32, AX // 32 floats per loop
	JL   remainder

loop:
	// a[0]

	BYTE $0xC5; BYTE $0xFC; BYTE $0x28; BYTE $0x06          //vmovaps ymm0,yword [rsi]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x28; BYTE $0x0A          //vmovaps ymm1,yword [rdx]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x59; BYTE $0xC1          //vmulps ymm0,ymm0,ymm1
  BYTE $0xC5; BYTE $0xFC; BYTE $0x29; BYTE $0x07          //vmovaps yword [rdi],ymm0

	BYTE $0xC5; BYTE $0xFC; BYTE $0x28; BYTE$0x46; BYTE $32  //  vmovaps ymm0,yword [rsi+0x20]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x28; BYTE$0x4A; BYTE $32  //  vmovaps ymm1,yword [rdx+0x20]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x59; BYTE$0xC1;             //  vmulps ymm0,ymm0,ymm1
	BYTE $0xC5; BYTE $0xFC; BYTE $0x29; BYTE$0x47; BYTE $32  //  vmovaps yword [rdi+0x20],ymm0
	
	BYTE $0xC5; BYTE $0xFC; BYTE $0x28; BYTE$0x46; BYTE $64  //  vmovaps ymm0,yword [rsi+0x40]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x28; BYTE$0x4A; BYTE $64  //  vmovaps ymm1,yword [rdx+0x40]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x59; BYTE$0xC1;             //  vmulps ymm0,ymm0,ymm1
	BYTE $0xC5; BYTE $0xFC; BYTE $0x29; BYTE$0x47; BYTE $64  //  vmovaps yword [rdi+0x40],ymm0


	BYTE $0xC5; BYTE $0xFC; BYTE $0x28; BYTE$0x46; BYTE $96  //  vmovaps ymm0,yword [rsi+0x60]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x28; BYTE$0x4A; BYTE $96  //  vmovaps ymm1,yword [rdx+0x60]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x59; BYTE$0xC1;             //  vmulps ymm0,ymm0,ymm1
	BYTE $0xC5; BYTE $0xFC; BYTE $0x29; BYTE$0x47; BYTE $96  //  vmovaps yword [rdi+0x60],ymm0
	





	ADDQ $64, SI         // increment 8 iterations 4 * 16
	ADDQ $64, DI
	ADDQ $64, DX

	SUBQ $32, AX         // Count down 2*8 floats
	JGE  loop            // Repeat

remainder:
	ADDQ $32, AX
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
