Researching vector multiplication optimizations
==================================================

My studies [research.md](research.md)

Running
-----------

```bash
  go get github.com/gohxs/vec-benchmark
  go test github.com/gohxs/vec-benchmark -bench .
```

The optimizations are based on this function which multiply each element of two
vectors into a third vector

all vectors needs to be same size, out vector can only be bigger or equal in size

```go
func Mul(vec1, vec2, out []float32) {
  for i := 0; i < len(vec1); i++ {
    out[i] = vec1[i] * vec2[i]
  }
}
```

Sample result for

```
NWorkers = 4
vecSize  = 10000 * NWorkers * 8 // Vector size is aligned to fit in 256bit and NWorkers
```

```
goos: linux
goarch: amd64
pkg: github.com/gohxs/vec-benchmark
BenchmarkVecSingle/Mul-4                   10000            156385 ns/op
BenchmarkVecSingle/MulFunc-4               10000            163037 ns/op
BenchmarkVecSingle/asm.Mulf32x4sse-4       20000             62565 ns/op
BenchmarkVecSingle/asm.Mulf32x8avx-4       50000             32298 ns/op
BenchmarkVecSingle/cgo.Mulf32x4sse-4       20000             63110 ns/op
BenchmarkVecSingle/cgo.Mulf32x8avx-4       30000             59067 ns/op
BenchmarkVecRoutines/Mul-4                 10000            100538 ns/op
BenchmarkVecRoutines/MulFunc-4             20000             97674 ns/op
BenchmarkVecRoutines/asm.Mulf32x4sse-4     30000             57022 ns/op
BenchmarkVecRoutines/asm.Mulf32x8avx-4     50000             29526 ns/op
BenchmarkVecRoutines/cgo.Mulf32x4sse-4     30000             56668 ns/op
BenchmarkVecRoutines/cgo.Mulf32x8avx-4     30000             56053 ns/op
BenchmarkVecWorker/Mul-4                   20000             97252 ns/op
BenchmarkVecWorker/MulFunc-4               20000             97171 ns/op
BenchmarkVecWorker/asm.Mulf32x4sse-4       30000             58726 ns/op
BenchmarkVecWorker/asm.Mulf32x8avx-4       50000             30492 ns/op
BenchmarkVecWorker/cgo.Mulf32x4sse-4       30000             57325 ns/op
BenchmarkVecWorker/cgo.Mulf32x8avx-4       30000             55132 ns/op
PASS
ok      github.com/gohxs/vec-benchmark  38.849s
```


