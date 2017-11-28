// +build !noasm,!appengine

#include "textflag.h"

// func VecMulf32x8(a, b, out []float64)
TEXT ·VecMulf32x8(SB), $0-72
	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DX
	MOVQ out_base+48(FP),DI   // Destination
	MOVQ out_len+56(FP), CX

	// Smaller size for CX
	CMPQ    a_len+8(FP), CX   // CX = max( len(out), len(a), len(b) )
	CMOVQLE a_len+8(FP), CX  
	CMPQ    b_len+32(FP), CX
	CMOVQLE b_len+32(FP), CX

	MOVQ    DX, BX
	ANDQ    $15, BX            // BX = &y & OxF
	JZ      no_align           // if BX == 0 { goto div_no_trim }

	// An alignment could happen here?
	// Align on 16-bit boundary test
	MOVSS (SI), X0    // X0 = s[i]
	MULSS (DX), X0    // X0 *= t[i]
	MOVSS  X0, (DI)   // dst[i] = X0

	ADDQ $4, SI
	ADDQ $4, DX
	ADDQ $4, DI
	DECQ  CX                // --CX

no_align:
	SUBQ $32, CX   // n floats per loop
	JL   remainder
loop:
	// a[0]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x28; BYTE $0x06;            	//C5FC1006          vmovups ymm0,yword [rsi]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x28; BYTE $0x4E;BYTE $0x20; 	//C5FC104E20        vmovups ymm1,yword [rsi+0x20]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x28; BYTE $0x56;BYTE $0x40; 	//C5FC105640        vmovups ymm2,yword [rsi+0x40]
	BYTE $0xC5; BYTE $0xFC; BYTE $0x28; BYTE $0x5E;BYTE $0x60; 	//C5FC105E60        vmovups ymm3,yword [rsi+0x60]

	BYTE $0xC5; BYTE $0xFC; BYTE $0x59; BYTE $0x02;            	//C5FC5902          vmulps ymm0,ymm0,yword [rdx]
	BYTE $0xC5; BYTE $0xF4; BYTE $0x59; BYTE $0x4A;BYTE $0x20; 	//C5F4594A20        vmulps ymm1,ymm1,yword [rdx+0x20]
	BYTE $0xC5; BYTE $0xEC; BYTE $0x59; BYTE $0x52;BYTE $0x40; 	//C5EC595240        vmulps ymm2,ymm2,yword [rdx+0x40]
	BYTE $0xC5; BYTE $0xE4; BYTE $0x59; BYTE $0x5A;BYTE $0x60; 	//C5E4595A60        vmulps ymm3,ymm3,yword [rdx+0x60]

	BYTE $0xC5; BYTE $0xFC; BYTE $0x29; BYTE $0x07;            	//C5FC1107          vmovups yword [rdi],ymm0
	BYTE $0xC5; BYTE $0xFC; BYTE $0x29; BYTE $0x4F;BYTE $0x20; 	//C5FC114F20        vmovups yword [rdi+0x20],ymm1
	BYTE $0xC5; BYTE $0xFC; BYTE $0x29; BYTE $0x57;BYTE $0x40; 	//C5FC115740        vmovups yword [rdi+0x40],ymm2
	BYTE $0xC5; BYTE $0xFC; BYTE $0x29; BYTE $0x5F;BYTE $0x60; 	//C5FC115F60        vmovups yword [rdi+0x60],ymm3

	// this is faster than do a single add 
	// a single reg and offseting in MOV ptrs
	ADDQ $4*32, SI         // increment sizeof(float32)4 * n floats
	ADDQ $4*32, DI 
	ADDQ $4*32, DX

	SUBQ $32, CX          // Count down n floats
	JGE  loop             // Repeat

remainder:
	ADDQ $32, CX
	JE   done

remainderloop:        // 1 by 1
	MOVSS (SI), X0
	MULSS (DX), X0
	MOVSS X0, (DI)

	ADDQ $4, SI
	ADDQ $4, DI
	ADDQ $4, DX

	LOOP remainderloop

done:
	RET
