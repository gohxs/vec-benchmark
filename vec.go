package vec

// Mul multiplies each element of vec1 to vec2 placing result in out
func Mul(vec1, vec2, out []float32) {
	for i := 0; i < len(vec1); i++ {
		out[i] = vec1[i] * vec2[i]
	}
}

// For Benchmark testing purposes
func baseMulFunc(v1, v2 float32) float32 {
	return v1 * v2
}

// MulFunc add each element of array using a function
func MulFunc(vec1, vec2, out []float32) {
	for i := 0; i < len(vec1); i++ {
		out[i] = baseMulFunc(vec1[i], vec2[i])
	}
}
