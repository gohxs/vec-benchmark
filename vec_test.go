package vec_test

import (
	"reflect"
	"sync"
	"testing"
	"unsafe"

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
		{"MulASMx4g", vec.MulASMx4g},
		{"MulASMx4", vec.MulASMx4},
		{"MulCGOx4", vec.MulCGOx4},
		{"MulCGOx4g", vec.MulCGOx4g},
		{"MulCGOx8", vec.MulCGOx8},
		//{"vecCGo512", vec.VecCGo512},
	}

	NWorkers = 2
	vecSize  = 32 * NWorkers * 8 // 8 floats to do 256bit operation4 million

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

// Adding Big vectors
//
func TestSlice(t *testing.T) {
	for i := 0; i < 1000; i++ {
		r := make([]float32, len(vec1))
		addr := []uintptr{}
		//diff := []int{}
		for i := range r {
			u := uintptr(unsafe.Pointer(&r[i]))
			addr = append(addr, u)
			if i != 0 {
				if u-addr[i-1] != 4 {
					t.Fatal("We have a subject here!")
				}
				//diff = append(diff, int(u-addr[i-1]))
			}
		}
	}
}

func TestVectorSingle(t *testing.T) {
	for _, f := range testFuncs {
		t.Run(f.name, func(t *testing.T) {
			out := make([]float32, len(vec1))
			f.fn(vec1, vec2, out)
			if !reflect.DeepEqual(sample, out) {
				t.Fatal("Value mismatch")
			}
			t.Log(out)
		})
	}
}
func TestVectorRoutines(t *testing.T) {
	for _, f := range testFuncs {
		out := make([]float32, len(vec1))
		t.Run(f.name, func(t *testing.T) {
			goVector(vec1, vec2, out, f.fn)
			if !reflect.DeepEqual(sample, out) {
				t.Fatal("Value mismatch")
			}
			t.Log(out)
		})
	}
}
func TestVectorWorker(t *testing.T) {
	for _, f := range testFuncs {
		out := make([]float32, len(vec1))
		t.Run(f.name, func(t *testing.T) {
			goWorkerVector(vec1, vec2, out, f.fn)
			if !reflect.DeepEqual(sample, out) {
				t.Fatal("Value mismatch")
			}
			t.Log(out)
		})
	}
}

func BenchmarkVectorSingle(b *testing.B) {
	for _, f := range testFuncs {
		b.Run(f.name, func(b *testing.B) {
			out := make([]float32, len(vec1))
			for n := b.N; n >= 0; n-- { // is this safe?
				f.fn(vec1, vec2, out)
			}
		})
	}
}
func BenchmarkVectorRoutines(b *testing.B) {
	for _, f := range testFuncs {
		b.Run(f.name, func(b *testing.B) {
			out := make([]float32, len(vec1))
			for n := b.N; n >= 0; n-- { // is this safe?
				goVector(vec1, vec2, out, f.fn)
			}
		})
	}
}
func BenchmarkVectorWorker(b *testing.B) {
	for _, f := range testFuncs {
		b.Run(f.name, func(b *testing.B) {
			out := make([]float32, len(vec1))
			for n := b.N; n >= 0; n-- { // is this safe?
				goWorkerVector(vec1, vec2, out, f.fn)
			}
		})
	}
}

type vecFunc func(a, b, c []float32)

func goVector(vec1, vec2, out []float32, fn vecFunc) {
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

// Warmed go rountines
// Data chunk
// worker data
type workerData struct {
	vec1, vec2, out []float32
	fn              vecFunc
}
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

func goWorkerVector(vec1, vec2, out []float32, fn vecFunc) {
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
