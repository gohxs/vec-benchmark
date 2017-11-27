[bits 64]


vmovups ymm0, yword[rsi]
vmovups ymm1, yword[rdx]
vmulps ymm0,ymm1
vmovups yword [rdi],ymm0

vmovups ymm0, yword[rsi+0x20]
vmovups ymm1, yword[rdx+0x20]
vmulps ymm0,ymm1
vmovups yword [rdi+0x20],ymm0
