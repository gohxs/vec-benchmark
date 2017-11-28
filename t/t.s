[bits 64]

;nasm t.s ; ndisasm -b 64 t
vmovups ymm0, yword[rsi]
vmovups ymm1, yword[rsi+0x20]
vmovups ymm2, yword[rsi+0x40]
vmovups ymm3, yword[rsi+0x60]

vmulps ymm0, yword[rdx]
vmulps ymm1, yword[rdx+0x20]
vmulps ymm2, yword[rdx+0x40]
vmulps ymm3, yword[rdx+0x60]

vmovups yword[rdi], ymm0
vmovups yword[rdi+0x20], ymm1
vmovups yword[rdi+0x40], ymm2
vmovups yword[rdi+0x60], ymm3





;aligned
vmovaps ymm0, yword[rsi]
vmovaps ymm1, yword[rdx]
vmulps ymm0,ymm1
vmovaps yword [rdi],ymm0

vmovaps ymm0, yword[rsi+0x20]
vmovaps ymm1, yword[rdx+0x20]
vmulps ymm0,ymm1
vmovaps yword [rdi+0x20],ymm0

vmovups ymm0, yword[rsi]
vmovups ymm1, yword[rdx]
vmulps ymm0,ymm1
vmovups yword [rdi],ymm0

vmovups ymm0, yword[rsi+0x20]
vmovups ymm1, yword[rdx+0x20]
vmulps ymm0,ymm1
vmovups yword [rdi+0x20],ymm0
