Note: Of course it is stupid to generate long linked lists, but it was done to prove the point.

```
go test -bench . insert_test.go 
goos: linux
goarch: amd64
BenchmarkInsertHeadSlice10-2            16823118                77.7 ns/op
BenchmarkInsertHeadSlice100-2            5883536               173 ns/op
BenchmarkInsertHeadSlice1000-2           1000000              1062 ns/op
BenchmarkInsertHeadSlice10000-2           123944              9333 ns/op
BenchmarkInsertHeadSlice100000-2           12022             88350 ns/op
BenchmarkInsertHeadSlice1000000-2           1440            773285 ns/op
BenchmarkInsertHeadList10-2              9093858               132 ns/op
BenchmarkInsertHeadList100-2            11352920               153 ns/op
BenchmarkInsertHeadList1000-2           10429632               152 ns/op
BenchmarkInsertHeadList10000-2          11601094               144 ns/op
BenchmarkInsertHeadList100000-2         11034750               161 ns/op
BenchmarkInsertHeadList1000000-2        12002636               144 ns/op
PASS
ok      command-line-arguments  29.799s
```
