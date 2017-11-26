


#### Anatomy of simple vector add func

```go
func VecGo(vec1, vec2, out []float32) {
  for i := 0; i < len(vec1); i++ {
    out[i] = vec1[i] * vec2[i]
  }
}
```

`go compile base.go -S` 

Outputs:

```asm
0x0000 00000 (base.go:3)        TEXT    "".VecGo(SB), NOSPLIT, $8-72
0x0000 00000 (base.go:3)        SUBQ    $8, SP
0x0004 00004 (base.go:3)        MOVQ    BP, (SP)
0x0008 00008 (base.go:3)        LEAQ    (SP), BP
0x000c 00012 (base.go:3)        FUNCDATA        $0, gclocals·1c3c8a9d47ed40f27c10312f31f2a755(SB)
0x000c 00012 (base.go:3)        FUNCDATA        $1, gclocals·69c1753bd5f81501d95132d08af04464(SB)
0x000c 00012 (base.go:3)        MOVQ    "".vec1+24(SP), AX
0x0011 00017 (base.go:3)        MOVQ    "".out+64(SP), CX
0x0016 00022 (base.go:3)        MOVQ    "".out+72(SP), DX
0x001b 00027 (base.go:3)        MOVQ    "".vec2+48(SP), BX
0x0020 00032 (base.go:3)        MOVQ    "".vec2+40(SP), SI
0x0025 00037 (base.go:3)        MOVQ    "".vec1+16(SP), DI
0x002a 00042 (base.go:3)        MOVL    $0, R8
0x002d 00045 (base.go:4)        JMP     56
0x002f 00047 (base.go:5)        MOVSS   X0, (CX)(R8*4)
0x0035 00053 (base.go:4)        INCQ    R8
0x0038 00056 (base.go:4)        CMPQ    R8, AX
0x003b 00059 (base.go:4)        JGE     89
0x003d 00061 (base.go:5)        MOVSS   (DI)(R8*4), X0
0x0043 00067 (base.go:5)        CMPQ    R8, BX
0x0046 00070 (base.go:5)        JCC     98
0x0048 00072 (base.go:5)        MOVSS   (SI)(R8*4), X1
0x004e 00078 (base.go:5)        MULSS   X1, X0
0x0052 00082 (base.go:5)        CMPQ    R8, DX
0x0055 00085 (base.go:5)        JCS     47
0x0057 00087 (base.go:5)        JMP     98
0x0059 00089 (base.go:7)        MOVQ    (SP), BP
0x005d 00093 (base.go:7)        ADDQ    $8, SP
0x0061 00097 (base.go:7)        RET
0x0062 00098 (base.go:5)        PCDATA  $0, $1
0x0062 00098 (base.go:5)        CALL    runtime.panicindex(SB)
0x0067 00103 (base.go:5)        UNDEF
```


