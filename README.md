Researching vector multiplication optimizations
==================================================

My studies [research.md](research.md)

Running
-----------

```bash
  go get github.com/gohxs/vec-benchmark
  go test github.com/gohxs/vec-benchmark -bench .
```

Sample result for

```
NWorkers = 4
vecSize  = 10000 * NWorkers * 8
```

```
goos: linux
goarch: amd64
pkg: github.com/gohxs/vec-benchmark
BenchmarkVecSingle/Mul-4                   10000            157489 ns/op
BenchmarkVecSingle/MulFunc-4               10000            161500 ns/op
BenchmarkVecSingle/asm.Mulf32x4sse-4       20000             61391 ns/op
BenchmarkVecSingle/cgo.Mulf32x4sse-4       20000             62902 ns/op
BenchmarkVecSingle/cgo.Mulf32x8xva-4       20000             60549 ns/op
BenchmarkVecRoutines/Mul-4                 20000            102947 ns/op
BenchmarkVecRoutines/MulFunc-4             10000            104255 ns/op
BenchmarkVecRoutines/asm.Mulf32x4sse-4     30000             56851 ns/op
BenchmarkVecRoutines/cgo.Mulf32x4sse-4     30000             56759 ns/op
BenchmarkVecRoutines/cgo.Mulf32x8xva-4     30000             54143 ns/op
BenchmarkVecWorker/Mul-4                   10000            101114 ns/op
BenchmarkVecWorker/MulFunc-4               20000             96828 ns/op
BenchmarkVecWorker/asm.Mulf32x4sse-4       30000             57762 ns/op
BenchmarkVecWorker/cgo.Mulf32x4sse-4       30000             58025 ns/op
BenchmarkVecWorker/cgo.Mulf32x8xva-4       30000             53736 ns/op
PASS
ok      github.com/gohxs/vec-benchmark  32.354s
```


