```bash
$ go test -bench .
goos: windows
goarch: amd64
pkg: sorts
cpu: Intel(R) Core(TM) i7-7700K CPU @ 4.20GHz
BenchmarkIntSort/builtin-8                  1036            979482 ns/op
BenchmarkIntSort/merge-8                    1688            675672 ns/op
BenchmarkIntSort/bubble-8                     13          83916777 ns/op
BenchmarkIntSort/insert-8                    106          10326882 ns/op
BenchmarkIntSort/insert2-8                    91          11037781 ns/op
BenchmarkIntSort/insert_bisect-8             438           2728966 ns/op
BenchmarkIntSort/quick-8                    1748            580117 ns/op
BenchmarkUintSort/radix_count-8             3075            348839 ns/op
PASS
```
