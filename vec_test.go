package vec_test

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"testing"

	vec "github.com/gohxs/vec-benchmark"
	"github.com/gohxs/vec-benchmark/asm"
	"github.com/gohxs/vec-benchmark/cgo"
	"github.com/gohxs/vec-benchmark/cl"
)

var (
	// TestFuncs
	testFuncs = []struct {
		name string
		fn   vec.MulFunc
	}{
		{"   vec.VecMulf32x1", vec.Mul},
		{" vec.VecMulf32x1Fn", vec.MulEFunc},
		{"asm.VecMulf32x4sse", asm.VecMulf32x4},
		{"asm.VecMulf32x8avx", asm.VecMulf32x8},
		{"cgo.VecMulf32x4sse", cgo.VecMulf32x4},
		{"cgo.VecMulf32x8avx", cgo.VecMulf32x8},
		{"      cl.VecMulf32", cl.VecMulf32},
	}
	NWorkers = runtime.NumCPU() // Workers for multiple go routines
	vecSize  int
	//vec1, vec2, out []float32
	//sample          = make([]float32, len(vec1))
	worker = vec.NewWorkerPool(NWorkers)
)

func init() {
	log.Println("Workers:", runtime.GOMAXPROCS(NWorkers))
	vecSize = 10 * NWorkers * 8 // Default 10 floats
	vecSizeEnv := os.Getenv("VECSIZE")
	if vecSizeEnv != "" {
		newVecSize, err := strconv.Atoi(vecSizeEnv)
		if err != nil {
			log.Println("Environ var VECSIZE has wrong entry, using default.")
		} else {
			vecSize = newVecSize
		}
	}
	log.Println("Running with vecSize:", vecSize)
}

func TestAlignment(t *testing.T) {
	vec1, vec2, out := createVecs(32, 32*2)

	// simulate
	//
	for _, f := range testFuncs {
		t.Run("Single/"+f.name, func(t *testing.T) {
			// Pre align?
			for i := 0; i < 8; i++ {
				t.Log("Starting at:", i, &out[i])
				f.fn(vec1[i:], vec2[i:], out[i+2:])
				checkVec(16, out[i:])
				f.fn(vec1[i:], vec2[i:], out[i+1:])
				checkVec(16, out[i:])
				f.fn(vec1[i+1:], vec2[i:], out[i:])
				checkVec(16, out[i:])
				f.fn(vec1[i:], vec2[i:], out[i:])
				checkVec(16, out[i:])
			}
			//}
		})
	}
}

func TestCalc(t *testing.T) {
	// Visual test
	//vecSize := 10 * 8 * NWorkers
	sample := make([]float32, vecSize)
	for i := range sample {
		sample[i] = float32(i+1) * 2
	}
	for _, f := range testFuncs {
		t.Run("Single/"+f.name, func(t *testing.T) {
			vec1, vec2, out := createVecs(vecSize, vecSize*2)
			f.fn(vec1, vec2, out)
			t.Logf("Vector:\nVec1:  %2v\nOut:   %2v\nSample:%2v\n", vec1, out, sample)
			if err := checkVec(vecSize, out); err != nil {
				t.Fatal(err)
			}
		})
		t.Run("Routine/"+f.name, func(t *testing.T) {
			vec1, vec2, out := createVecs(vecSize, vecSize*2)
			vec.GoVecMul(NWorkers, vec1, vec2, out, f.fn)
			t.Logf("Vector:\nVec1:  %2v\nOut:   %2v\nSample:%2v\n", vec1, out, sample)
			if err := checkVec(vecSize, out); err != nil {
				t.Fatal(err)
			}
		})
		t.Run("Worker/"+f.name, func(t *testing.T) {
			vec1, vec2, out := createVecs(vecSize, vecSize*2)
			worker.VecMul(vec1, vec2, out, f.fn)
			t.Logf("Vector:\nVec1:  %2v\nOut:   %2v\nSample:%2v\n", vec1, out, sample)
			if err := checkVec(vecSize, out); err != nil {
				t.Fatal(err)
			}
		})
	}

}

func TestVec(t *testing.T) {
	for _, f := range testFuncs {
		t.Run("Single"+f.name, func(t *testing.T) {
			vec1, vec2, out := createVecs(vecSize, vecSize)
			f.fn(vec1, vec2, out)
			if err := checkVec(vecSize, out); err != nil {
				log.Fatal(err)
			}
		})
		t.Run("Routine/"+f.name, func(t *testing.T) {
			vec1, vec2, out := createVecs(vecSize, vecSize)
			vec.GoVecMul(NWorkers, vec1, vec2, out, f.fn)
			if err := checkVec(vecSize, out); err != nil {
				log.Fatal(err)
			}
		})
		t.Run("Worker"+f.name, func(t *testing.T) {
			vec1, vec2, out := createVecs(vecSize, vecSize)
			worker.VecMul(vec1, vec2, out, f.fn)
			if err := checkVec(vecSize, out); err != nil {
				log.Fatal(err)
			}
		})
	}
}

func BenchmarkVec(b *testing.B) {
	vec1, vec2, out := createVecs(vecSize, vecSize)
	benchHelper(b, vec1, vec2, out)
}

///////////////////////////////////////////////////////////////////////////////
//                             HELPER FUNCS                                  //
///////////////////////////////////////////////////////////////////////////////
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
func checkVec(vecSize int, out []float32) error {
	for i := 0; i < vecSize; i++ {
		// Same as vec1 (which is indexes) * vec2 (which is constant 2)
		if out[i] != float32(i+1)*2 {
			return fmt.Errorf("Value mismatch at %d expected: %f got %f", i, float32(i)*2, out[i])
		}
	}
	// Check remaining output for possible overflow
	for i := vecSize; i < len(out); i++ {
		if out[i] != 0 {
			return fmt.Errorf("Value mismatch at %d expected: %f got %f", i, float32(i)*2, out[i])
			//t.Fatalf("Possible overflow at %d, expected:0 got: %f", i, out[i])
		}
	}
	return nil // no error
}

func benchHelper(b *testing.B, vec1, vec2, out []float32) {
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
