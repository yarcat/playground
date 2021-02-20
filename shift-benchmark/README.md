```
$ go test -bench .
goos: linux
goarch: amd64
pkg: github.com/yarcat/playground/shift-benchmark
BenchmarkShiftSmall-4            1777986               657.1 ns/op
BenchmarkShiftLarge-4                180           6531875 ns/op
BenchmarkShiftRevSmall-4         1805539               665.9 ns/op
BenchmarkShiftRevLarge-4             180           6562823 ns/op
PASS
ok      github.com/yarcat/playground/shift-benchmark    14.305s
```