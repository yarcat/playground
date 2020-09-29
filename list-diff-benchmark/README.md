```
$ go test -bench . diff_test.go
goos: linux
goarch: amd64
BenchmarkDiffBisect10-2                 11214938               106 ns/op
BenchmarkDiffBisect100-2                  755307              1633 ns/op
BenchmarkDiffBisect1000-2                  27706             42183 ns/op
BenchmarkDiffBisect10000-2                  1875            536959 ns/op
BenchmarkDiffBisect100000-2                  189           6399612 ns/op
BenchmarkDiffBisect1000000-2                  15          69505401 ns/op
BenchmarkDiffBisect10000000-2                  2         752527940 ns/op
BenchmarkDiffMap10-2                     2077942               579 ns/op
BenchmarkDiffMap100-2                     140996              7655 ns/op
BenchmarkDiffMap1000-2                     16776             69910 ns/op
BenchmarkDiffMap10000-2                     1422            736660 ns/op
BenchmarkDiffMap100000-2                     120           9673775 ns/op
BenchmarkDiffMap1000000-2                      5         223671942 ns/op
BenchmarkDiffMap10000000-2                     1        2648435463 ns/op
PASS
ok      command-line-arguments  28.421s
```
