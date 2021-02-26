No attempt to write perfect sort functions. Just comparing the asymptotics of
quick and merge sorts. Please note that quick uses element [0] as a pivot point.
This means that we expect the random case to perform approx equally, and
sorted/reversed cases to suck.

```
$ go test -bench .
goos: linux
goarch: amd64
pkg: github.com/yarcat/playground/sort-benchmark
BenchmarkQuickOrdered-2              176           5978517 ns/op
BenchmarkQuickRandom-2              3156            371981 ns/op
BenchmarkQuickRev-2                  186           6142657 ns/op
BenchmarkMergeOrdered-2             3580            319631 ns/op
BenchmarkMergeRandom-2              2956            368074 ns/op
BenchmarkMergeRev-2                 3082            334254 ns/op
PASS
ok      github.com/yarcat/playground/sort-benchmark     8.886s
```
