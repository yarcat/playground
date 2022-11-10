No attempt to write perfect sort functions. Just comparing the asymptotics of
quick and merge sorts for *linked lists*. Please note that quick uses element
[0] as a pivot point. This means that we expect the random case to perform
approx equally, and sorted/reversed cases to suck.

Also benchmarks generate lots of garbage, so I'm not sure how reliable it gets
after some time. Most probably not reliable at all. But at least we can see
our expectations clearly.

```
$ go test -bench . -benchmem 
goos: windows
goarch: amd64
pkg: sort-benchmark
cpu: Intel(R) Core(TM) i7-7700K CPU @ 4.20GHz
BenchmarkSort/old_qsort/ordered-8                      8         132012675 ns/op           20000 B/op       1250 allocs/op
BenchmarkSort/new_qsort/ordered-8                      6         171538533 ns/op           26666 B/op       1666 allocs/op
BenchmarkSort/old_qsort/random-8                    1100           1093667 ns/op             145 B/op          9 allocs/op
BenchmarkSort/new_qsort/random-8                    1153           1026220 ns/op             138 B/op          8 allocs/op
BenchmarkSort/old_qsort/reversed-8                     8         133900512 ns/op           20000 B/op       1250 allocs/op
BenchmarkSort/new_qsort/reversed-8                     7         151973714 ns/op           22857 B/op       1428 allocs/op
BenchmarkSort/old_msort/ordered-8                   1118           1038488 ns/op             143 B/op          8 allocs/op
BenchmarkSort/new_msort/ordered-8                   2212            485322 ns/op              72 B/op          4 allocs/op
BenchmarkSort/old_msort/random-8                     762           1510275 ns/op             209 B/op         13 allocs/op
BenchmarkSort/new_msort/random-8                    1023           1090734 ns/op             156 B/op          9 allocs/op
BenchmarkSort/old_msort/reversed-8                  1117           1040278 ns/op             143 B/op          8 allocs/op
BenchmarkSort/new_msort/reversed-8                  2398            471350 ns/op              66 B/op          4 allocs/op
PASS
ok      sort-benchmark  19.511s
```
