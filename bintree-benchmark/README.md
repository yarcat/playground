Various binary tree comparisons.

Again, we aren't trying to create something super efficient, but rather trying
to understand differencies between recursive and non-recursive operations.

## 2021/03/19 results

Another day, another numbers (-;

This time I've included a benchmark for a "generic" tree implementation (the one
that stores values as `interface{}`, and uses a `less` function to compare
elements). It is expected that this implementation's gonna be slower, and it is
slower (at least in my benchmarks).

```
$ go test -bench .
goos: linux
goarch: amd64
pkg: example.com
BenchmarkIntTreeIterativeInsert/100-8                     167382              6279 ns/op
BenchmarkIntTreeIterativeInsert/1000-8                     14323            108405 ns/op
BenchmarkIntTreeIterativeInsert/10000-8                      847           1348536 ns/op
BenchmarkIntTreeIterativeInsert/100000-8                      33          33769556 ns/op
BenchmarkIntTreeIterativeInsertParentPtr/100-8            202290              6053 ns/op
BenchmarkIntTreeIterativeInsertParentPtr/1000-8            10000            105248 ns/op
BenchmarkIntTreeIterativeInsertParentPtr/10000-8             906           1423558 ns/op
BenchmarkIntTreeIterativeInsertParentPtr/100000-8             36          32489886 ns/op
BenchmarkIntTreeRecursiveInsert/100-8                     219153              6764 ns/op
BenchmarkIntTreeRecursiveInsert/1000-8                     10094            129623 ns/op
BenchmarkIntTreeRecursiveInsert/10000-8                      508           2084201 ns/op
BenchmarkIntTreeRecursiveInsert/100000-8                      26          45146209 ns/op
BenchmarkIntTreeRecursiveInsertParentPtr/100-8            173766              6349 ns/op
BenchmarkIntTreeRecursiveInsertParentPtr/1000-8             9370            130910 ns/op
BenchmarkIntTreeRecursiveInsertParentPtr/10000-8             452           2653407 ns/op
BenchmarkIntTreeRecursiveInsertParentPtr/100000-8             24          41995683 ns/op
BenchmarkTreeIterativeInsertParentPtr/100-8               129399              9442 ns/op
BenchmarkTreeIterativeInsertParentPtr/1000-8                5833            177970 ns/op
BenchmarkTreeIterativeInsertParentPtr/10000-8                391           3334678 ns/op
BenchmarkTreeIterativeInsertParentPtr/100000-8                19          59854362 ns/op
PASS
ok      example.com     35.621s
```

## 2021/03/18 results

```
 go test -bench . -benchmem
goos: linux
goarch: amd64
pkg: github.com/yarcat/playground/bintree-benchmark
BenchmarkIterativeInsert/100-8                            235442              4642 ns/op            2408 B/op        101 allocs/op
BenchmarkIterativeInsert/1000-8                            14610             79274 ns/op           24008 B/op       1001 allocs/op
BenchmarkIterativeInsert/10000-8                             962           1192958 ns/op          240008 B/op      10001 allocs/op
BenchmarkIterativeInsert/100000-8                             51          21224437 ns/op         2400095 B/op     100001 allocs/op
BenchmarkIterativeInsertParentPtr/100-8                   208026              5532 ns/op            2408 B/op        101 allocs/op
BenchmarkIterativeInsertParentPtr/1000-8                   13893             81823 ns/op           24008 B/op       1001 allocs/op
BenchmarkIterativeInsertParentPtr/10000-8                    966           1163609 ns/op          240008 B/op      10001 allocs/op
BenchmarkIterativeInsertParentPtr/100000-8                    49          21375318 ns/op         2400019 B/op     100001 allocs/op
BenchmarkRecursiveInsert/100-8                            249668              4759 ns/op            2408 B/op        101 allocs/op
BenchmarkRecursiveInsert/1000-8                            10000            103145 ns/op           24008 B/op       1001 allocs/op
BenchmarkRecursiveInsert/10000-8                             724           1587742 ns/op          240009 B/op      10001 allocs/op
BenchmarkRecursiveInsert/100000-8                             40          27473843 ns/op         2400069 B/op     100001 allocs/op
BenchmarkRecursiveInsertParentPtr/100-8                   222666              5055 ns/op            2408 B/op        101 allocs/op
BenchmarkRecursiveInsertParentPtr/1000-8                   10000            110604 ns/op           24008 B/op       1001 allocs/op
BenchmarkRecursiveInsertParentPtr/10000-8                    661           1699921 ns/op          240010 B/op      10001 allocs/op
BenchmarkRecursiveInsertParentPtr/100000-8                    37          28851233 ns/op         2400046 B/op     100001 allocs/op
PASS
ok      github.com/yarcat/playground/bintree-benchmark  21.518s
```
