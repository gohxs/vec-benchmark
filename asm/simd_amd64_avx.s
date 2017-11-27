#include "textflag.h"

// Based on Chewxy vec64f
// func VecMulf32x8(a, b, out []float64)
TEXT ·VecMulf32x8(SB), $0-72

	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DX
	MOVQ out_base+48(FP),DI   // Destination

	MOVQ a_len+8(FP), AX      // len(a) into AX
	MOVQ b_len+32(FP), BX     // len(b) into BX
	MOVQ out_len+56(FP), CX   // len(out) into DX

	CMPQ AX, BX   // Check if a,b are same lenght
	JNE  panic  
	CMPQ AX, CX
	JG   panic    // if output is smaller than inputs 
	
	
	SUBQ $32, AX   // n floats per loop
	JL   remainder

loop:
	// a[0]

	BYTE $0xC5; BYTE $0xFC; BYTE $0x10; BYTE $0x06           //  vmovups ymm0,yword [rsi]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x10; BYTE $0x0A           //  vmovups ymm1,yword [rdx]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x59; BYTE $0xC1           //  vmulps ymm0,ymm0,ymm1
	BYTE $0xC5; BYTE $0xFC; BYTE $0x11; BYTE $0x07           //  vmovups yword [rdi],ymm0

	
	BYTE $0xC5; BYTE $0xFC; BYTE $0x10; BYTE$0x46; BYTE $32  //  vmovups ymm0,yword [rsi+0x20]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x10; BYTE$0x4A; BYTE $32  //  vmovups ymm1,yword [rdx+0x20]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x59; BYTE$0xC1;           //  vmulps ymm0,ymm0,ymm1
	BYTE $0xC5; BYTE $0xFC; BYTE $0x11; BYTE$0x47; BYTE $32  //  vmovups yword [rdi+0x20],ymm0

	BYTE $0xC5; BYTE $0xFC; BYTE $0x10; BYTE$0x46; BYTE $64  //  vmovups ymm0,yword [rsi+0x40]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x10; BYTE$0x4A; BYTE $64  //  vmovups ymm1,yword [rdx+0x40]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x59; BYTE$0xC1;           //  vmulps ymm0,ymm0,ymm1
	BYTE $0xC5; BYTE $0xFC; BYTE $0x11; BYTE$0x47; BYTE $64  //  vmovups yword [rdi+0x40],ymm0


	BYTE $0xC5; BYTE $0xFC; BYTE $0x10; BYTE$0x46; BYTE $96  //  vmovups ymm0,yword [rsi+0x60]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x10; BYTE$0x4A; BYTE $96  //  vmovups ymm1,yword [rdx+0x60]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x59; BYTE$0xC1;           //  vmulps ymm0,ymm0,ymm1
	BYTE $0xC5; BYTE $0xFC; BYTE $0x11; BYTE$0x47; BYTE $96  //  vmovups yword [rdi+0x60],ymm0
	
	


	// This is faster on goroutines why?
	// $128 for 4 steps
	// If these are $64 goroutine is half of the time
	// Maybe due to alignment?

	ADDQ $128, SI         // increment sizeof(float32)4 * n floats
	ADDQ $128, DI 
	ADDQ $128, DX

	SUBQ $32, AX           // Count down n floats
	JGE  loop             // Repeat

remainder:
	ADDQ $32, AX
	JE   done

remainderloop:        // 1 by 1
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
