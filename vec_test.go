package vec_test

import (
	"reflect"
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
	vec1, vec2, out = createVecs(100000 * NWorkers * 8)
	worker          = vec.NewWorkerPool(NWorkers)
)

func createVecs(vecSize int) ([]float32, []float32, []float32) {
	vec1 := make([]float32, vecSize)
	vec2 := make([]float32, vecSize)
	out := make([]float32, vecSize)
	for i := 0; i < len(vec1); i++ {
		vec1[i] = float32(i)
		vec2[i] = 2
	}
	return vec1, vec2, out
}

func TestCalc(t *testing.T) {
	vecSize := 1 * 10
	vec1 := make([]float32, vecSize)
	vec2 := make([]float32, vecSize)
	out := make([]float32, vecSize*2)
	sample := make([]float32, vecSize)
	for i := 0; i < len(vec1); i++ {
		vec1[i] = float32(i)
		vec2[i] = 2
		sample[i] = vec1[i] * vec2[i]
	}

	for _, f := range testFuncs {
		t.Run(f.name, func(t *testing.T) {
			// 3x per type?
			f.fn(vec1, vec2, out)
			t.Logf("Single Vec1:   %2v", vec1)
			t.Logf("Single Out:    %2v", out)
			t.Logf("Single Sample: %2v", sample)
			for i := range sample {
				if sample[i] != out[i] {
					t.Fatal("Value mismatch")
				}
			}

			vec.GoVecMul(NWorkers, vec1, vec2, out, f.fn)
			t.Logf("goVec  Vec1:   %2v", vec1)
			t.Logf("goVec  Out:    %2v", out)
			t.Logf("goVec  Sample: %2v", sample)
			for i := range sample {
				if sample[i] != out[i] {
					t.Fatal("Value mismatch")
				}
			}
			worker.VecMul(vec1, vec2, out, f.fn)
			t.Logf("Worker Vec1:   %2v", vec1)
			t.Logf("Worker Out:    %2v", out)
			t.Logf("Worker Sample: %2v", sample)
			for i := range sample {
				if sample[i] != out[i] {
					t.Fatal("Value mismatch")
				}
			}
		})
	}

}

func TestWorker(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			vecSize := 4 * 8
			vec1, vec2, out := createVecs(vecSize)
			worker.VecMul(vec1, vec2, out, vec.Mul)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestVec(t *testing.T) {
	sample := make([]float32, len(vec1))
	vec.Mul(vec1, vec2, sample) // Safe implementation

	for _, f := range testFuncs {
		t.Run("Single"+f.name, func(t *testing.T) {
			f.fn(vec1, vec2, out)
			if !reflect.DeepEqual(sample, out) {
				t.Fatal("Value mismatch")
			}
		})
		t.Run("Routine/"+f.name, func(t *testing.T) {
			vec.GoVecMul(NWorkers, vec1, vec2, out, f.fn)
			if !reflect.DeepEqual(sample, out) {
				t.Fatal("Value mismatch")
			}
		})
		t.Run("Worker"+f.name, func(t *testing.T) {
			worker.VecMul(vec1, vec2, out, f.fn)
			if !reflect.DeepEqual(sample, out) {
				t.Fatal("Value mismatch")
			}
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
