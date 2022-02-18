```text
$ go test -bench . -benchmem
goos: windows
goarch: amd64
pkg: tree-benchmark
cpu: Intel(R) Core(TM) i7-7700K CPU @ 4.20GHz
BenchmarkNew-8             14665             80681 ns/op           24016 B/op       1001 allocs/op
BenchmarkArena-8           24033             50119 ns/op           24640 B/op          4 allocs/op
PASS
ok      tree-benchmark  3.873s
```
