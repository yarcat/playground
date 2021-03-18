Various binary tree comparisons.

Again, we aren't trying to create something super efficient, but rather trying
to understand differencies between recursive and non-recursive operations.

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
