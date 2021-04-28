# Alloc-free

An attempt to see what leads to alloc-free code in Golang and what not.

## Surprises

I've managed to make the `ExecInterface` version to be allocation-free. However,
I had to use pointers to a string and a byte-slice for this. If `StringWriter`
(or `BytesWriter`) is simply based on a string, then there is 16 bytes per string
and 24 per slice allocation for every option. Which makes sense, since interfaces
cannot contain values, but this is a good reminder to think more than twice while
using interfaces.

## Current Results

```
go test -bench=. -benchmem .
goos: linux
goarch: amd64
pkg: github.com/yarcat/playground/allocations
BenchmarkSet/Set-8                      38127736                31.16 ns/op            0 B/op          0 allocs/op
BenchmarkExec/Factory-8                  9092862               114.8 ns/op            56 B/op          2 allocs/op
BenchmarkExec/Lambda-8                  27203308                40.27 ns/op            0 B/op          0 allocs/op
BenchmarkExec/Functors/Ptr-8            26227058                40.09 ns/op            0 B/op          0 allocs/op
BenchmarkExec/Functors/Val-8            26946501                40.84 ns/op            0 B/op          0 allocs/op
BenchmarkExec/Functors/Hlp-8            13166934                84.49 ns/op            0 B/op          0 allocs/op
BenchmarkExecI/ExecI-8                  27920743                39.74 ns/op            0 B/op          0 allocs/op
BenchmarkExec2/Exec2-8                  18925665                63.70 ns/op           24 B/op          1 allocs/op
BenchmarkExec3/Exec3-8                  33926996                34.66 ns/op            0 B/op          0 allocs/op
BenchmarkExec4/Exec4-8                  31617224                34.48 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/yarcat/playground/allocations        14.660s
```
