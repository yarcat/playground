No attempt to write perfect sort functions. Just comparing the asymptotics of
quick and merge sorts for *linked lists*. Please note that quick uses element
[0] as a pivot point. This means that we expect the random case to perform
approx equally, and sorted/reversed cases to suck.

Also benchmarks generate lots of garbage, so I'm not sure how reliable it gets
after some time. Most probably not reliable at all. But at least we can see
our expectations clearly.

```
$ go test -bench . -benchmem
goos: linux
goarch: amd64
pkg: github.com/yarcat/playground/sort-benchmark
BenchmarkQuickOrdered-2              188           6168298 ns/op           16000 B/op       1000 allocs/op
BenchmarkQuickRandom-2              3397            326665 ns/op           16000 B/op       1000 allocs/op
BenchmarkQuickRev-2                  194           5985385 ns/op           16000 B/op       1000 allocs/op
BenchmarkMergeOrdered-2             3400            326412 ns/op           16000 B/op       1000 allocs/op
BenchmarkMergeRandom-2              3241            364131 ns/op           16000 B/op       1000 allocs/op
BenchmarkMergeRev-2                 3565            332913 ns/op           16000 B/op       1000 allocs/op
```
