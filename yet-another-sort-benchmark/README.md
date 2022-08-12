```bash
$ go test -bench .
goos: linux
goarch: amd64
pkg: sorts
cpu: Intel(R) Celeron(R) CPU N3350 @ 1.10GHz
BenchmarkSort-2                      691           1522877 ns/op
BenchmarkSortBubble-2                  5         238475533 ns/op
BenchmarkSortInsert-2                 20          57414062 ns/op
BenchmarkSortInsertBisect-2           48          21471761 ns/op
BenchmarkSortQuick-2                1045           1134086 ns/op
PASS
ok      sorts   9.525s
```
