[bits 64]

;nasm t.s ; ndisasm -b 64 t
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
