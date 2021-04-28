# Alloc-free

An attempt to see what leads to alloc-free code in Golang and what not.

## Surprises

TBD

## Current Results

```
go test -bench=. -benchmem .
goos: linux
goarch: amd64
pkg: github.com/yarcat/playground/allocations
BenchmarkSet-8                                  31356510                31.91 ns/op            0 B/op          0 allocs/op
BenchmarkExecFuncBufio/Factory-8                 9852060               116.30 ns/op           56 B/op          2 allocs/op
BenchmarkExecFuncBufio/Lambda-8                 31124596                37.93 ns/op            0 B/op          0 allocs/op
BenchmarkExecFuncBufio/Functors/Ptr-8           29809599                39.97 ns/op            0 B/op          0 allocs/op
BenchmarkExecFuncBufio/Functors/Val-8           26675809                41.56 ns/op            0 B/op          0 allocs/op
BenchmarkExecFuncBufio/Functors/Hlp-8           13596487                86.17 ns/op            0 B/op          0 allocs/op
BenchmarkExecI-8                                28307331                39.82 ns/op            0 B/op          0 allocs/op
BenchmarkExecFuncSenderPtr-8                    17440234                63.30 ns/op           24 B/op          1 allocs/op
BenchmarkExecFuncSenderVal-8                    32717740                34.39 ns/op            0 B/op          0 allocs/op
BenchmarkExecFuncSenderPrealloc-8               32952367                33.92 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/yarcat/playground/allocations        13.690s
```
