Researching vector multiplication optimizations
==================================================

My studies [research.md](research.md)

Running
-----------

```bash
  go get github.com/gohxs/vec-benchmark
  go test github.com/gohxs/vec-benchmark -bench .
```

Sample data:

```go
var (
  NWorkers = 2                 // Workers for multiple go routines
  vecSize  = 32 * NWorkers * 8 // Aligned to NWorkers and maximum 8 floats (512 float32)

  vec1   = make([]float32, vecSize)
  vec2   = make([]float32, vecSize)
  sample = make([]float32, vecSize)
)
```

Sample test

```go
out := make([]float32,vecSize)
vec.Mul(vec1,vec2,out)

```

Sample result

```
goos: linux
goarch: amd64
pkg: github.com/gohxs/vec-benchmark
BenchmarkVecSingle/Mul-4                   5000000               258 ns/op
BenchmarkVecSingle/MulFunc-4               5000000               267 ns/op
BenchmarkVecSingle/asm.MulSSEx4gi-4        3000000               467 ns/op
BenchmarkVecSingle/cgo.MulSSEx4-4         20000000               108 ns/op
BenchmarkVecSingle/cgo.MulSSEx4gi-4         200000              7081 ns/op
BenchmarkVecSingle/cgo.MulXVAx8-4         20000000                83 ns/op
BenchmarkVecRoutines/Mul-4                 2000000               725 ns/op
BenchmarkVecRoutines/MulFunc-4             2000000               728 ns/op
BenchmarkVecRoutines/asm.MulSSEx4gi-4      2000000               980 ns/op
BenchmarkVecRoutines/cgo.MulSSEx4-4        2000000               613 ns/op
BenchmarkVecRoutines/cgo.MulSSEx4gi-4       200000             10079 ns/op
BenchmarkVecRoutines/cgo.MulXVAx8-4        3000000               590 ns/op
BenchmarkVecWorker/Mul-4                   2000000               799 ns/op
BenchmarkVecWorker/MulFunc-4               2000000               813 ns/op
BenchmarkVecWorker/asm.MulSSEx4gi-4        1000000              1030 ns/op
BenchmarkVecWorker/cgo.MulSSEx4-4          2000000               685 ns/op
BenchmarkVecWorker/cgo.MulSSEx4gi-4         200000              9866 ns/op
BenchmarkVecWorker/cgo.MulXVAx8-4          2000000               644 ns/op
PASS
ok      github.com/gohxs/vec-benchmark  36.255s
```


