package vec

// VecMulFunc prototype for multiplication funcs
type MulFunc func(a, b, c []float32)

// Mul multiplies each element of vec1 to vec2 placing result in out
func Mul(vec1, vec2, out []float32) {
	for i := 0; i < len(vec1); i++ {
		out[i] = vec1[i] * vec2[i]
	}
}

// For Benchmark testing purposes
func baseEMulFunc(v1, v2 float32) float32 {
	return v1 * v2
}

// MulFunc add each element of array using a function
func MulEFunc(vec1, vec2, out []float32) {
	for i := 0; i < len(vec1); i++ {
		out[i] = baseEMulFunc(vec1[i], vec2[i])
	}
}
