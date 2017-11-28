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

27-11-2017

```
goos: linux
goarch: amd64
pkg: github.com/gohxs/vec-benchmark
BenchmarkVecSmall/Single/_________VecMulgo-4      10000000        166 ns/op
BenchmarkVecSmall/Single/_____VecMulFuncgo-4      10000000        191 ns/op
BenchmarkVecSmall/Single/asm.VecMulf32x4sse-4     50000000         29.1 ns/op
BenchmarkVecSmall/Single/asm.VecMulf32x8avx-4     100000000        15.3 ns/op
BenchmarkVecSmall/Single/cgo.VecMulf32x4sse-4     20000000         86.6 ns/op
BenchmarkVecSmall/Single/cgo.VecMulf32x8avx-4     20000000         85.5 ns/op
BenchmarkVecSmall/Routine/_________VecMulgo-4      2000000        914 ns/op
BenchmarkVecSmall/Routine/_____VecMulFuncgo-4      2000000        929 ns/op
BenchmarkVecSmall/Routine/asm.VecMulf32x4sse-4     2000000        734 ns/op
BenchmarkVecSmall/Routine/asm.VecMulf32x8avx-4     2000000        808 ns/op
BenchmarkVecSmall/Routine/cgo.VecMulf32x4sse-4     2000000        939 ns/op
BenchmarkVecSmall/Routine/cgo.VecMulf32x8avx-4     2000000        957 ns/op
BenchmarkVecSmall/Worker/_________VecMulgo-4       1000000       1190 ns/op
BenchmarkVecSmall/Worker/_____VecMulFuncgo-4       1000000       1225 ns/op
BenchmarkVecSmall/Worker/asm.VecMulf32x4sse-4      1000000       1005 ns/op
BenchmarkVecSmall/Worker/asm.VecMulf32x8avx-4      1000000       1099 ns/op
BenchmarkVecSmall/Worker/cgo.VecMulf32x4sse-4      1000000       1300 ns/op
BenchmarkVecSmall/Worker/cgo.VecMulf32x8avx-4      1000000       1273 ns/op
BenchmarkVecBig/Single/_________VecMulgo-4           10000     155749 ns/op
BenchmarkVecBig/Single/_____VecMulFuncgo-4           10000     181747 ns/op
BenchmarkVecBig/Single/asm.VecMulf32x4sse-4          20000      63763 ns/op
BenchmarkVecBig/Single/asm.VecMulf32x8avx-4          20000      58711 ns/op
BenchmarkVecBig/Single/cgo.VecMulf32x4sse-4          20000      62572 ns/op
BenchmarkVecBig/Single/cgo.VecMulf32x8avx-4          20000      62099 ns/op
BenchmarkVecBig/Routine/_________VecMulgo-4          20000      93637 ns/op
BenchmarkVecBig/Routine/_____VecMulFuncgo-4          10000     102044 ns/op
BenchmarkVecBig/Routine/asm.VecMulf32x4sse-4         30000      60345 ns/op
BenchmarkVecBig/Routine/asm.VecMulf32x8avx-4         20000      53928 ns/op
BenchmarkVecBig/Routine/cgo.VecMulf32x4sse-4         30000      60901 ns/op
BenchmarkVecBig/Routine/cgo.VecMulf32x8avx-4         30000      58339 ns/op
BenchmarkVecBig/Worker/_________VecMulgo-4           20000     103466 ns/op
BenchmarkVecBig/Worker/_____VecMulFuncgo-4           10000     110428 ns/op
BenchmarkVecBig/Worker/asm.VecMulf32x4sse-4          20000      57864 ns/op
BenchmarkVecBig/Worker/asm.VecMulf32x8avx-4          30000      56748 ns/op
BenchmarkVecBig/Worker/cgo.VecMulf32x4sse-4          30000      61035 ns/op
BenchmarkVecBig/Worker/cgo.VecMulf32x8avx-4          30000      61112 ns/op
PASS
ok      github.com/gohxs/vec-benchmark  70.360s
```

26-11-2017:

```
goos: linux
goarch: amd64
pkg: github.com/gohxs/vec-benchmark
BenchmarkVecSmall/Single/_________VecMulgo-4     10000000       167 ns/op
BenchmarkVecSmall/Single/_____VecMulFuncgo-4     10000000       171 ns/op
BenchmarkVecSmall/Single/asm.VecMulf32x4sse-4    50000000        32.4 ns/op
BenchmarkVecSmall/Single/asm.VecMulf32x8avx-4    100000000       15.9 ns/op
BenchmarkVecSmall/Single/cgo.VecMulf32x4sse-4    20000000        87.0 ns/op
BenchmarkVecSmall/Single/cgo.VecMulf32x8avx-4    20000000        85.4 ns/op
BenchmarkVecSmall/Routine/_________VecMulgo-4     2000000       933 ns/op
BenchmarkVecSmall/Routine/_____VecMulFuncgo-4     2000000       924 ns/op
BenchmarkVecSmall/Routine/asm.VecMulf32x4sse-4    2000000       726 ns/op
BenchmarkVecSmall/Routine/asm.VecMulf32x8avx-4    2000000       797 ns/op
BenchmarkVecSmall/Routine/cgo.VecMulf32x4sse-4    2000000       930 ns/op
BenchmarkVecSmall/Routine/cgo.VecMulf32x8avx-4    2000000      1022 ns/op
BenchmarkVecSmall/Worker/_________VecMulgo-4      1000000      1182 ns/op
BenchmarkVecSmall/Worker/_____VecMulFuncgo-4      1000000      1196 ns/op
BenchmarkVecSmall/Worker/asm.VecMulf32x4sse-4     1000000      1022 ns/op
BenchmarkVecSmall/Worker/asm.VecMulf32x8avx-4     1000000      1090 ns/op
BenchmarkVecSmall/Worker/cgo.VecMulf32x4sse-4     1000000      1268 ns/op
BenchmarkVecSmall/Worker/cgo.VecMulf32x8avx-4     1000000      1295 ns/op
BenchmarkVecBig/Single/_________VecMulgo-4          10000    156223 ns/op
BenchmarkVecBig/Single/_____VecMulFuncgo-4          10000    161283 ns/op
BenchmarkVecBig/Single/asm.VecMulf32x4sse-4         20000     64221 ns/op
BenchmarkVecBig/Single/asm.VecMulf32x8avx-4         20000     58880 ns/op
BenchmarkVecBig/Single/cgo.VecMulf32x4sse-4         20000     63158 ns/op
BenchmarkVecBig/Single/cgo.VecMulf32x8avx-4         20000     62082 ns/op
BenchmarkVecBig/Routine/_________VecMulgo-4         20000     95356 ns/op
BenchmarkVecBig/Routine/_____VecMulFuncgo-4         20000     98905 ns/op
BenchmarkVecBig/Routine/asm.VecMulf32x4sse-4        30000     55805 ns/op
BenchmarkVecBig/Routine/asm.VecMulf32x8avx-4        30000     54938 ns/op
BenchmarkVecBig/Routine/cgo.VecMulf32x4sse-4        30000     56425 ns/op
BenchmarkVecBig/Routine/cgo.VecMulf32x8avx-4        30000     56495 ns/op
BenchmarkVecBig/Worker/_________VecMulgo-4          20000     98456 ns/op
BenchmarkVecBig/Worker/_____VecMulFuncgo-4          20000    104448 ns/op
BenchmarkVecBig/Worker/asm.VecMulf32x4sse-4         30000     56922 ns/op
BenchmarkVecBig/Worker/asm.VecMulf32x8avx-4         30000     54801 ns/op
BenchmarkVecBig/Worker/cgo.VecMulf32x4sse-4         30000     56743 ns/op
BenchmarkVecBig/Worker/cgo.VecMulf32x8avx-4         30000     57845 ns/op
PASS
ok      github.com/gohxs/vec-benchmark  74.490s
```


