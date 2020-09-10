```
go test -bench .
goos: linux
goarch: amd64
BenchmarkSliceCopyWithAppendEmpty-12            428990575                2.42 ns/op
BenchmarkSliceCopyWithAppendSmall-12              202466              6092 ns/op
BenchmarkSliceCopyWithAppendMedium-12               1164            876474 ns/op
BenchmarkSliceCopyWithAppendLarge-12                 106          10908358 ns/op
BenchmarkSliceCopyWithPreallocEmpty-12          491223627                2.28 ns/op
BenchmarkSliceCopyWithPreallocSmall-12            416534              3217 ns/op
BenchmarkSliceCopyWithPreallocMedium-12             4124            246446 ns/op
BenchmarkSliceCopyWithPreallocLarge-12               499           2233837 ns/op
PASS
ok      _/usr/local/google/home/yarcat/src/slice-copy-benchmark 10.940s
```
