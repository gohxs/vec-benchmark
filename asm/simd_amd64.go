package asm

//Mulf32x4 Multiplies 4 float32 elements from each a,b into out
func Mulf32x4(a, b []float32, out []float32)

//VecMul multiplies each element into out
func VecMul(a, b []float32, out []float32)

//VecMulChewxy multiply two float32 based on vecf64 package from chewxy with slight changes for 32bits
func VecMulChewxy(a, b, out []float32)
