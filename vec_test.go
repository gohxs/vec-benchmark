package vec_test

import (
	"reflect"
	"sync"
	"testing"

	vec "github.com/gohxs/vec-benchmark"
)

var (
	// TestFuncs
	testFuncs = []struct {
		name string
		fn   vecFunc
	}{
		{"Mul", vec.Mul},
		{"MulFunc", vec.MulFunc},
		{"asm.MulSSEx4gi", vec.MulASMSSEx4gi},
		{"asm.MulChewxy", vec.MulASMChewxy},
		//{"MulASMx4", vec.MulASMx4},
		{"cgo.MulSSEx4", vec.MulCGOSSEx4},
		{"cgo.MulSSEx4gi", vec.MulCGOSSEx4gi},
		{"cgo.MulXVAx8", vec.MulCGOXVAx8},
	}

	NWorkers = 2                 // Workers for multiple go routines
	vecSize  = 20 * NWorkers * 8 // 8 floats to do 256bit operation

	vec1   = make([]float32, vecSize)
	vec2   = make([]float32, vecSize)
	sample = make([]float32, vecSize)
)

func init() {
	workersLaunch()
	for i := 0; i < len(vec1); i++ {
		vec1[i] = float32(i)
		vec2[i] = 2
	}
	vec.Mul(vec1, vec2, sample)
}

func TestVecSingle(t *testing.T) {
	out := make([]float32, vecSize)
	for _, f := range testFuncs {
		t.Run(f.name, func(t *testing.T) {
			f.fn(vec1, vec2, out)
			t.Log(sample)
			t.Log(out)
			if !reflect.DeepEqual(sample, out) {
				t.Fatal("Value mismatch")
			}
		})
	}
}
func TestVecRoutines(t *testing.T) {
	out := make([]float32, vecSize)
	for _, f := range testFuncs {
		t.Run(f.name, func(t *testing.T) {
			goVec(vec1, vec2, out, f.fn)
			t.Log(sample)
			t.Log(out)
			if !reflect.DeepEqual(sample, out) {
				t.Fatal("Value mismatch")
			}
		})
	}
}
func TestVecWorker(t *testing.T) {
	out := make([]float32, vecSize)
	for _, f := range testFuncs {
		t.Run(f.name, func(t *testing.T) {
			goWorkerVec(vec1, vec2, out, f.fn)
			t.Log(sample)
			t.Log(out)
			if !reflect.DeepEqual(sample, out) {
				t.Fatal("Value mismatch")
			}
		})
	}
}

// Benchmarks
func BenchmarkVecSingle(b *testing.B) {
	out := make([]float32, vecSize)
	for _, f := range testFuncs {
		b.Run(f.name, func(b *testing.B) {
			for n := b.N; n >= 0; n-- { // is this safe?
				f.fn(vec1, vec2, out)
			}
		})
	}
}
func BenchmarkVecRoutines(b *testing.B) {
	out := make([]float32, vecSize)
	for _, f := range testFuncs {
		b.Run(f.name, func(b *testing.B) {
			for n := b.N; n >= 0; n-- { // is this safe?
				goVec(vec1, vec2, out, f.fn)
			}
		})
	}
}
func BenchmarkVecWorker(b *testing.B) {
	out := make([]float32, vecSize)
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
