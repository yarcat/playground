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
BenchmarkSet/Set-8      35552845                32.49 ns/op            0 B/op          0 allocs/op
BenchmarkExec/Factory-8                  9541806               123.0 ns/op            56 B/op          2 allocs/op
BenchmarkExec/Lambda-8                  27838359                43.37 ns/op            0 B/op          0 allocs/op
BenchmarkExec/Functors-8                27430675                41.27 ns/op            0 B/op          0 allocs/op
BenchmarkExecI/ExecI-8                  26037552                40.00 ns/op            0 B/op          0 allocs/op
BenchmarkExec2/Exec2-8                  16896396                62.85 ns/op           24 B/op          1 allocs/op
BenchmarkExec3/Exec3-8                  30386654                35.93 ns/op            0 B/op          0 allocs/op
BenchmarkExec4/Exec4-8                  31253431                35.20 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/yarcat/playground/allocations        10.444s
```
