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
BenchmarkSortLinkedList/256-8              31489             31960 ns/op            4112 B/op        257 allocs/op
BenchmarkSortLinkedList/512-8               9152            119512 ns/op            8208 B/op        513 allocs/op
BenchmarkSortLinkedList/1024-8              2636            419320 ns/op           16403 B/op       1025 allocs/op
BenchmarkSortLinkedList/2048-8               680           1755330 ns/op           32808 B/op       2049 allocs/op
BenchmarkSortSkipList/256-8                11283             97768 ns/op           39071 B/op       1025 allocs/op
BenchmarkSortSkipList/512-8                 4910            227451 ns/op           78060 B/op       2049 allocs/op
BenchmarkSortSkipList/1024-8                2592            428158 ns/op          156051 B/op       4097 allocs/op
BenchmarkSortSkipList/2048-8                1130            984413 ns/op          312043 B/op       8193 allocs/op
BenchmarkSortSkipListInt/256-8             17046             74918 ns/op           32926 B/op        769 allocs/op
BenchmarkSortSkipListInt/512-8              7482            163165 ns/op           65779 B/op       1537 allocs/op
BenchmarkSortSkipListInt/1024-8             3532            342018 ns/op          131466 B/op       3073 allocs/op
BenchmarkSortSkipListInt/2048-8             1876            650391 ns/op          262872 B/op       6145 allocs/op
PASS
ok      example.org     20.998s
```
