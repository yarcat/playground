Benchmark various implementations of filter+map sequences with limited
concurrency.

## ConcurrentFilterMap1

Using buffered output. Since output size is now known in advance we buffer it
using `len(in)`.

`maxConcurrency` go-routines is spawned to do work.

## ConcurrentFilterMap2

Using two extra go-routines to avoid buffering.

`maxConcurrency` go-routines is spawned to do work.

## ConcurrentFilterMap3

The most idiomatic implementation that resonates with "Don't communicate by
sharing memory, share memory by communicating."

Concurrency is limited using semaphore (a `maxConcurrency` -buffered channel).
Go-routine is spawned for every task.

See also:
- https://go-proverbs.github.io/
- https://golang.org/doc/effective_go.html#channels
- https://github.com/duffn/gophercon2018#rethinking-classical-concurrency-patterns

## Benchmark Results

The benchmarks were generated on my old chrome book (it is quite slow).

```
goos: linux
goarch: amd64
BenchmarkConcurrentFilterMap1-4             6348            180627 ns/op
BenchmarkConcurrentFilterMap2-4             1936            573523 ns/op
BenchmarkConcurrentFilterMap3-4             1567            669313 ns/op
```
