```
$ go test -bench . -benchmem
goos: linux
goarch: amd64
pkg: github.com/yarcat/playground/shift-benchmark
BenchmarkShiftSmall-2            1950262               661.8 ns/op             0 B/op          0 allocs/op
BenchmarkShiftLarge-2                162           7371765 ns/op               0 B/op          0 allocs/op
BenchmarkShift2Small-2           1499940               717.2 ns/op             0 B/op          0 allocs/op
BenchmarkShift2Large-2               147           8719832 ns/op               0 B/op          0 allocs/op
BenchmarkShiftRevSmall-2         2475946               460.3 ns/op             0 B/op          0 allocs/op
BenchmarkShiftRevLarge-2             192           5629785 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/yarcat/playground/shift-benchmark    12.485s
```
