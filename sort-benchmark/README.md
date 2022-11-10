No attempt to write perfect sort functions. Just comparing the asymptotics of
quick and merge sorts for *linked lists*. Please note that quick uses element
[0] as a pivot point. This means that we expect the random case to perform
approx equally, and sorted/reversed cases to suck.

Also benchmarks generate lots of garbage, so I'm not sure how reliable it gets
after some time. Most probably not reliable at all. But at least we can see
our expectations clearly.

```
go test -bench . -benchmem
goos: linux
goarch: amd64
pkg: sort-benchmark
cpu: Intel(R) Celeron(R) CPU N3350 @ 1.10GHz
BenchmarkSort/ord/old_qsort-2                  2         531386944 ns/op               4 B/op          0 allocs/op
BenchmarkSort/ord/new_qsort-2                  2         530056282 ns/op               0 B/op          0 allocs/op
BenchmarkSort/ord/old_msort-2                345           3335615 ns/op               0 B/op          0 allocs/op
BenchmarkSort/ord/new_msort-2                705           1678797 ns/op               0 B/op          0 allocs/op
BenchmarkSort/rng/old_qsort-2                409           3139043 ns/op               0 B/op          0 allocs/op
BenchmarkSort/rng/new_qsort-2                464           2621409 ns/op               0 B/op          0 allocs/op
BenchmarkSort/rng/old_msort-2                297           4167340 ns/op               0 B/op          0 allocs/op
BenchmarkSort/rng/new_msort-2                405           2916307 ns/op               0 B/op          0 allocs/op
BenchmarkSort/rev/old_qsort-2                  3         483659589 ns/op               0 B/op          0 allocs/op
BenchmarkSort/rev/new_qsort-2                  3         494323130 ns/op               0 B/op          0 allocs/op
BenchmarkSort/rev/old_msort-2                352           3366407 ns/op               0 B/op          0 allocs/op
BenchmarkSort/rev/new_msort-2                848           1402031 ns/op               0 B/op          0 allocs/op
PASS
ok      sort-benchmark  26.145s
```
