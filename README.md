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
big   vec = 100000 * NWorkers * 8
small vec = 10 * NWorkers * 8
```

```
goos: linux
goarch: amd64
pkg: github.com/gohxs/vec-benchmark
BenchmarkVecSmall/Single/_________VecMulgo-4            10000000               168 ns/op               0 B/op          0 allocs/op
BenchmarkVecSmall/Single/_____VecMulFuncgo-4            10000000               172 ns/op               0 B/op          0 allocs/op
BenchmarkVecSmall/Single/asm.VecMulf32x4sse-4           50000000                29.1 ns/op             0 B/op          0 allocs/op
BenchmarkVecSmall/Single/asm.VecMulf32x8avx-4           100000000               19.1 ns/op             0 B/op          0 allocs/op
BenchmarkVecSmall/Single/cgo.VecMulf32x4sse-4           20000000                89.0 ns/op             0 B/op          0 allocs/op
BenchmarkVecSmall/Single/cgo.VecMulf32x8avx-4           20000000                86.8 ns/op             0 B/op          0 allocs/op
BenchmarkVecSmall/Routine/_________VecMulgo-4            1000000              1003 ns/op              16 B/op          1 allocs/op
BenchmarkVecSmall/Routine/_____VecMulFuncgo-4            2000000               970 ns/op              16 B/op          1 allocs/op
BenchmarkVecSmall/Routine/asm.VecMulf32x4sse-4           2000000               761 ns/op              16 B/op          1 allocs/op
BenchmarkVecSmall/Routine/asm.VecMulf32x8avx-4           2000000               774 ns/op              16 B/op          1 allocs/op
BenchmarkVecSmall/Routine/cgo.VecMulf32x4sse-4           2000000               948 ns/op              16 B/op          1 allocs/op
BenchmarkVecSmall/Routine/cgo.VecMulf32x8avx-4           2000000               949 ns/op              16 B/op          1 allocs/op
BenchmarkVecSmall/Worker/_________VecMulgo-4             1000000              1178 ns/op              16 B/op          1 allocs/op
BenchmarkVecSmall/Worker/_____VecMulFuncgo-4             1000000              1191 ns/op              16 B/op          1 allocs/op
BenchmarkVecSmall/Worker/asm.VecMulf32x4sse-4            1000000              1012 ns/op              16 B/op          1 allocs/op
BenchmarkVecSmall/Worker/asm.VecMulf32x8avx-4            1000000              1093 ns/op              16 B/op          1 allocs/op
BenchmarkVecSmall/Worker/cgo.VecMulf32x4sse-4            1000000              1258 ns/op              16 B/op          1 allocs/op
BenchmarkVecSmall/Worker/cgo.VecMulf32x8avx-4            1000000              1289 ns/op              16 B/op          1 allocs/op
BenchmarkVecBig/Single/_________VecMulgo-4                  1000           1802189 ns/op               0 B/op          0 allocs/op
BenchmarkVecBig/Single/_____VecMulFuncgo-4                  1000           1845062 ns/op               0 B/op          0 allocs/op
BenchmarkVecBig/Single/asm.VecMulf32x4sse-4                 1000           1440795 ns/op               0 B/op          0 allocs/op
BenchmarkVecBig/Single/asm.VecMulf32x8avx-4                 2000            601437 ns/op               0 B/op          0 allocs/op
BenchmarkVecBig/Single/cgo.VecMulf32x4sse-4                 1000           1457430 ns/op               0 B/op          0 allocs/op
BenchmarkVecBig/Single/cgo.VecMulf32x8avx-4                 1000           1469456 ns/op               0 B/op          0 allocs/op
BenchmarkVecBig/Routine/_________VecMulgo-4                 1000           1413166 ns/op              16 B/op          1 allocs/op
BenchmarkVecBig/Routine/_____VecMulFuncgo-4                 1000           1412326 ns/op              16 B/op          1 allocs/op
BenchmarkVecBig/Routine/asm.VecMulf32x4sse-4                1000           1375778 ns/op              16 B/op          1 allocs/op
BenchmarkVecBig/Routine/asm.VecMulf32x8avx-4                3000            573044 ns/op              16 B/op          1 allocs/op
BenchmarkVecBig/Routine/cgo.VecMulf32x4sse-4                1000           1399443 ns/op              16 B/op          1 allocs/op
BenchmarkVecBig/Routine/cgo.VecMulf32x8avx-4                1000           1430045 ns/op              16 B/op          1 allocs/op
BenchmarkVecBig/Worker/_________VecMulgo-4                  1000           1413125 ns/op              16 B/op          1 allocs/op
BenchmarkVecBig/Worker/_____VecMulFuncgo-4                  1000           1407591 ns/op              16 B/op          1 allocs/op
BenchmarkVecBig/Worker/asm.VecMulf32x4sse-4                 1000           1384380 ns/op              16 B/op          1 allocs/op
BenchmarkVecBig/Worker/asm.VecMulf32x8avx-4                 3000            573784 ns/op              16 B/op          1 allocs/op
BenchmarkVecBig/Worker/cgo.VecMulf32x4sse-4                 1000           1408114 ns/op              27 B/op          1 allocs/op
BenchmarkVecBig/Worker/cgo.VecMulf32x8avx-4                 1000           1482808 ns/op              21 B/op          1 allocs/op
PASS
ok      github.com/gohxs/vec-benchmark  61.542s
```


