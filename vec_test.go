package vec_test

import (
	"reflect"
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
		fn   vecFunc
	}{
		{"         VecMul", vec.Mul},
		{"     VecMulFunc", vec.MulFunc},
		{"asm.VecMulf32x4", asm.VecMulf32x4},
		{"asm.VecMulf32x8", asm.VecMulf32x8},
		{"cgo.VecMulf32x4", cgo.VecMulf32x4},
		{"cgo.VecMulf32x8", cgo.VecMulf32x8},
	}

	NWorkers = 4                    // Workers for multiple go routines
	vecSize  = 10000 * NWorkers * 8 // 8 floats to do 256bit operation

	vec1   = make([]float32, vecSize)
	vec2   = make([]float32, vecSize)
	out    = make([]float32, vecSize)
	sample = make([]float32, vecSize)
)

// Move this to other place
func init() {
	workersLaunch()
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
			goWorkerVec(vec1, vec2, out, f.fn)
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
			goWorkerVec(vec1, vec2, out, f.fn)
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
				goWorkerVec(vec1, vec2, out, f.fn)
			}
		})
	}
}

type vecFunc func(a, b, c []float32)

// routine Helper
func goVec(vec1, vec2, out []float32, fn vecFunc) {
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

// Workers go routines
// worker data
type workerData struct {
	vec1, vec2, out []float32
	fn              vecFunc
}

//worker
type workerChan struct {
	in   chan workerData
	done chan int
}

var workers []workerChan

func workerStart(in chan workerData, done chan int) {
	for {
		d := <-in
		d.fn(d.vec1, d.vec2, d.out)
		done <- 1
	}
}

func workersLaunch() {
	//gin := make(chan _bgData, 8) // Single in
	for i := 0; i < NWorkers; i++ {
		newWorker := workerChan{
			in:   make(chan workerData, 1),
			done: make(chan int, 1),
		}
		workers = append(workers, newWorker)

		go func(worker workerChan) {
			for {
				d := <-worker.in
				d.fn(d.vec1, d.vec2, d.out)
				worker.done <- 1
			}
		}(newWorker)
	}
}

func goWorkerVec(vec1, vec2, out []float32, fn vecFunc) {
	// send to bgRoutines
	for i, ch := range workers { // Divide workload between cores?
		sz := len(vec1) / NWorkers
		s := i * sz
		e := s + sz

		ch.in <- workerData{
			vec1[s:e],
			vec2[s:e],
			out[s:e],
			fn,
		}
	}
	// Wait for outs from all the workers
	for _, ch := range workers {
		<-ch.done
	}
}
