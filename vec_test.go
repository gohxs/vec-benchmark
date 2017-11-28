package vec_test

import (
	"runtime"
	"sync"
	"testing"

	vec "github.com/gohxs/vec-benchmark"
	"github.com/gohxs/vec-benchmark/asm"
	"github.com/gohxs/vec-benchmark/cgo"
)

var (
	// TestFuncs
	testFuncs = []struct {
		name string
		fn   vec.MulFunc
	}{
		{"         VecMulgo", vec.Mul},
		{"     VecMulFuncgo", vec.MulEFunc},
		{"asm.VecMulf32x4sse", asm.VecMulf32x4},
		{"asm.VecMulf32x8avx", asm.VecMulf32x8},
		{"cgo.VecMulf32x4sse", cgo.VecMulf32x4},
		{"cgo.VecMulf32x8avx", cgo.VecMulf32x8},
	}
	NWorkers        = runtime.NumCPU() // Workers for multiple go routines
	vec1, vec2, out = createVecs(10000*NWorkers*8, 10000*NWorkers*8)
	//sample          = make([]float32, len(vec1))
	worker = vec.NewWorkerPool(NWorkers)
)

func init() {
	runtime.GOMAXPROCS(NWorkers)
}

func createVecs(vecLen, outLen int) ([]float32, []float32, []float32) {
	vec1 := make([]float32, vecLen)
	vec2 := make([]float32, vecLen)
	out := make([]float32, outLen)
	for i := 0; i < vecLen && i < outLen; i++ {
		vec1[i] = float32(i + 1)
		vec2[i] = 2
	}
	return vec1, vec2, out
}

func checkVec(t *testing.T, vecSize int, out []float32) {
	for i := range out {
		if i < vecSize {

			if out[i] != float32(i+1)*2 {
				t.Logf("Worker Out:    %2v", out)
				t.Fatalf("Value mismatch at %d expected: %f got %f", i, float32(i)*2, out[i])
			}
		} else if out[i] != 0 {
			t.Logf("Worker Out:    %2v", out)
			t.Fatalf("Possible overflow at %d, expected:0 got: %f", i, out[i])
		}
	}
}

func TestCalc(t *testing.T) {
	vecSize := 3
	sample := make([]float32, vecSize)
	for i := range sample {
		sample[i] = float32(i+1) * 2
	}

	for _, f := range testFuncs {
		t.Run("Single/"+f.name, func(t *testing.T) {
			vec1, vec2, out := createVecs(vecSize, vecSize*2)
			// 3x per type?
			f.fn(vec1, vec2, out)
			t.Logf("Single Vec1:   %2v", vec1)
			t.Logf("Single Sample: %2v", sample)
			checkVec(t, vecSize, out)
		})
		t.Run("Routine/"+f.name, func(t *testing.T) {
			vec1, vec2, out := createVecs(vecSize, vecSize*2)
			vec.GoVecMul(NWorkers, vec1, vec2, out, f.fn)
			t.Logf("goVec  Vec1:   %2v", vec1)
			t.Logf("goVec  Sample: %2v", sample)
			checkVec(t, vecSize, out)
		})
		t.Run("Worker/"+f.name, func(t *testing.T) {
			vec1, vec2, out := createVecs(vecSize, vecSize*2)
			worker.VecMul(vec1, vec2, out, f.fn)
			t.Logf("Worker Vec1:   %2v", vec1)
			t.Logf("Worker Sample: %2v", sample)
			checkVec(t, vecSize, out)
		})
	}

}

func TestWorker(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			vecSize := 4 * 8
			vec1, vec2, out := createVecs(vecSize, vecSize)
			worker.VecMul(vec1, vec2, out, vec.Mul)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestVec(t *testing.T) {
	vecSize := 10 * NWorkers * 8
	sample := make([]float32, vecSize)
	for i := range sample {
		sample[i] = float32(i+1) * 2
	}
	for _, f := range testFuncs {
		t.Run("Single"+f.name, func(t *testing.T) {
			vec1, vec2, out := createVecs(vecSize, vecSize)
			f.fn(vec1, vec2, out)
			checkVec(t, vecSize, out)

		})
		t.Run("Routine/"+f.name, func(t *testing.T) {
			vec1, vec2, out := createVecs(vecSize, vecSize)
			vec.GoVecMul(NWorkers, vec1, vec2, out, f.fn)
			checkVec(t, vecSize, out)
		})
		t.Run("Worker"+f.name, func(t *testing.T) {
			vec1, vec2, out := createVecs(vecSize, vecSize)
			worker.VecMul(vec1, vec2, out, f.fn)
			checkVec(t, vecSize, out)
		})
	}
}

func BenchmarkVecSmall(b *testing.B) {
	// Create new local vectors from global ones
	vec1 := vec1[:10*NWorkers*8]
	vec2 := vec2[:10*NWorkers*8]
	out := out[:10*NWorkers*8]
	for _, f := range testFuncs {
		b.Run("Single/"+f.name, func(b *testing.B) {
			for n := b.N; n >= 0; n-- { // is this safe?
				f.fn(vec1, vec2, out)
			}
		})
	}
	for _, f := range testFuncs {
		b.Run("Routine/"+f.name, func(b *testing.B) {
			for n := b.N; n >= 0; n-- { // is this safe?
				vec.GoVecMul(NWorkers, vec1, vec2, out, f.fn)
			}
		})
	}
	for _, f := range testFuncs {
		b.Run("Worker/"+f.name, func(b *testing.B) {
			for n := b.N; n >= 0; n-- { // is this safe?
				worker.VecMul(vec1, vec2, out, f.fn)
			}
		})

	}
}

// Benchmarks
func BenchmarkVecBig(b *testing.B) {
	for _, f := range testFuncs {
		b.Run("Single/"+f.name, func(b *testing.B) {
			for n := b.N; n >= 0; n-- { // is this safe?
				f.fn(vec1, vec2, out)
			}
		})
	}
	for _, f := range testFuncs {
		b.Run("Routine/"+f.name, func(b *testing.B) {
			for n := b.N; n >= 0; n-- { // is this safe?
				vec.GoVecMul(NWorkers, vec1, vec2, out, f.fn)
			}
		})
	}
	for _, f := range testFuncs {
		b.Run("Worker/"+f.name, func(b *testing.B) {
			for n := b.N; n >= 0; n-- { // is this safe?
				worker.VecMul(vec1, vec2, out, f.fn)
			}
		})

	}
}
