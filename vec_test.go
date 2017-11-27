package vec_test

import (
	"log"
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

	NWorkers = runtime.NumCPU()      // Workers for multiple go routines
	vecSize  = 100000 * NWorkers * 8 // 8 floats to do 256bit operation

	vec1   = make([]float32, vecSize)
	vec2   = make([]float32, vecSize)
	out    = make([]float32, vecSize)
	sample = make([]float32, vecSize)

	worker = vec.NewWorkerPool(NWorkers)
)

// Move this to other place
func init() {
	for i := 0; i < len(vec1); i++ {
		vec1[i] = float32(i)
		vec2[i] = 2
	}
	vec.Mul(vec1, vec2, sample)
}

func TestCalc(t *testing.T) {
	vecSize = 4 * 8
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

			goVec(vec1, vec2, out, f.fn)
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
	log.Println("Testing multiple workers")
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			vecSize = 4 * 8 * NWorkers
			vec1 := make([]float32, vecSize)
			vec2 := make([]float32, vecSize)
			out := make([]float32, vecSize)
			worker.VecMul(vec1, vec2, out, vec.Mul)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestVecSingle(t *testing.T) {
	for _, f := range testFuncs {
		t.Run(f.name, func(t *testing.T) {
			f.fn(vec1, vec2, out)
			if !reflect.DeepEqual(sample, out) {
				t.Fatal("Value mismatch")
			}
		})
	}
}
func TestVecRoutines(t *testing.T) {
	for _, f := range testFuncs {
		t.Run(f.name, func(t *testing.T) {
			goVec(vec1, vec2, out, f.fn)
			if !reflect.DeepEqual(sample, out) {
				t.Fatal("Value mismatch")
			}
		})
	}
}
func TestVecWorker(t *testing.T) {
	for _, f := range testFuncs {
		t.Run(f.name, func(t *testing.T) {
			worker.VecMul(vec1, vec2, out, f.fn)
			if !reflect.DeepEqual(sample, out) {
				t.Fatal("Value mismatch")
			}
		})
	}
}

// Benchmarks
func BenchmarkVecSingle(b *testing.B) {
	for _, f := range testFuncs {
		b.Run(f.name, func(b *testing.B) {
			for n := b.N; n >= 0; n-- { // is this safe?
				f.fn(vec1, vec2, out)
			}
		})
	}
}
func BenchmarkVecRoutines(b *testing.B) {
	for _, f := range testFuncs {
		b.Run(f.name, func(b *testing.B) {
			for n := b.N; n >= 0; n-- { // is this safe?
				goVec(vec1, vec2, out, f.fn)
			}
		})
	}
}
func BenchmarkVecWorker(b *testing.B) {
	for _, f := range testFuncs {
		b.Run(f.name, func(b *testing.B) {
			for n := b.N; n >= 0; n-- { // is this safe?
				worker.VecMul(vec1, vec2, out, f.fn)
			}
		})
	}
}

// routine Helper
func goVec(vec1, vec2, out []float32, fn vec.MulFunc) {
	wg := sync.WaitGroup{}
	wg.Add(NWorkers)

	for i := 0; i < NWorkers; i++ { // Divide workload between cores?
		sz := len(vec1) / NWorkers
		go func(offs int) {
			fn(
				vec1[offs:offs+sz],
				vec2[offs:offs+sz],
				out[offs:offs+sz],
			)
			wg.Done()
		}(i * sz)
	}
	wg.Wait()
}
