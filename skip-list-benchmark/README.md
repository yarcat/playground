Comparing insertion sorts using regular linked and skip lists. There are two
implementations of the skip lists used here -- a "generic" one (that stores
interface{}) and fixed to ints.

It's expected that linked list is faster at first since it does less allocations,
but after certain amount of elements skip lists should become faster. On my
machine the border was around 1'000 elements.

## Disclaimer

This is a fun benchmark. I'm not trying to implement an efficient skip list here.

## Results

```
$ go test -bench . -benchmem
goos: linux
goarch: amd64
pkg: example.org
BenchmarkSortLinkedList/256-8              31333             33521 ns/op            4112 B/op        257 allocs/op
BenchmarkSortLinkedList/512-8               9612            112572 ns/op            8208 B/op        513 allocs/op
BenchmarkSortLinkedList/1024-8              2893            406987 ns/op           16402 B/op       1025 allocs/op
BenchmarkSortLinkedList/2048-8               634           1631921 ns/op           32809 B/op       2049 allocs/op
BenchmarkSortSkipList/256-8                15334             83614 ns/op           18591 B/op        769 allocs/op
BenchmarkSortSkipList/512-8                 6874            155489 ns/op           37103 B/op       1537 allocs/op
BenchmarkSortSkipList/1024-8                3490            359360 ns/op           74135 B/op       3073 allocs/op
BenchmarkSortSkipList/2048-8                1254            905418 ns/op          148193 B/op       6145 allocs/op
BenchmarkSortSkipListInt/256-8             24240             47080 ns/op           28752 B/op        257 allocs/op
BenchmarkSortSkipListInt/512-8             12099            105872 ns/op           57424 B/op        513 allocs/op
BenchmarkSortSkipListInt/1024-8             4898            235058 ns/op          114770 B/op       1025 allocs/op
BenchmarkSortSkipListInt/2048-8             2482            476615 ns/op          229463 B/op       2049 allocs/op
PASS
ok      example.org     21.301s
```
