```bash
$ go test -bench .
goos: linux
goarch: amd64
pkg: sorts
cpu: Intel(R) Celeron(R) CPU N3350 @ 1.10GHz
BenchmarkSort-2                      624           1763481 ns/op
BenchmarkSortBubble-2                  5         239264213 ns/op
```
