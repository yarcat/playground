## For 1_000 items

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

## For 10_000 items

```text
$ go test -bench . -benchmem
goos: windows
goarch: amd64
pkg: tree-benchmark
cpu: Intel(R) Core(TM) i7-7700K CPU @ 4.20GHz
BenchmarkNew-8              1017           1139137 ns/op          240018 B/op      10001 allocs/op
BenchmarkArena-8            1426            826852 ns/op          245824 B/op          4 allocs/op
PASS
ok      tree-benchmark  2.703s
```

## For 100_000 items

```text
$ go test -bench . -benchmem
goos: windows
goarch: amd64
pkg: tree-benchmark
cpu: Intel(R) Core(TM) i7-7700K CPU @ 4.20GHz
BenchmarkNew-8                64          20377164 ns/op         2400086 B/op     100001 allocs/op
BenchmarkArena-8              70          16032143 ns/op         2400323 B/op          4 allocs/op
PASS
ok      tree-benchmark  2.631s
```

## For 1M items

```text
$ go test -bench . -benchmem
goos: windows
goarch: amd64
pkg: tree-benchmark
cpu: Intel(R) Core(TM) i7-7700K CPU @ 4.20GHz
BenchmarkNew-8                 3         466900133 ns/op        24000048 B/op    1000001 allocs/op
BenchmarkArena-8               3         413100700 ns/op        24002624 B/op          4 allocs/op
PASS
ok      tree-benchmark  5.519s
```
