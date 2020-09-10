```
BenchmarkSliceCopyWithAppendEmpty-12                            429394074                2.43 ns/op
BenchmarkSliceCopyWithAppendSmall-12                              233440              5928 ns/op
BenchmarkSliceCopyWithAppendMedium-12                               1222            872521 ns/op
BenchmarkSliceCopyWithAppendLarge-12                                 104          11173250 ns/op
BenchmarkSliceCopyWithPreallocEmpty-12                          493602181                2.16 ns/op
BenchmarkSliceCopyWithPreallocSmall-12                            487420              3147 ns/op
BenchmarkSliceCopyWithPreallocMedium-12                             4681            246923 ns/op
BenchmarkSliceCopyWithPreallocLarge-12                               484           2335213 ns/op
BenchmarkSliceCopyWithIndexAndPreallocEmpty-12                  483791526                2.13 ns/op
BenchmarkSliceCopyWithIndexAndPreallocSmall-12                    366908              3156 ns/op
BenchmarkSliceCopyWithIndexAndPreallocMedium-12                     4519            238571 ns/op
BenchmarkSliceCopyWithIndexAndPreallocLarge-12                       492           2259025 ns/op

BenchmarkSliceCopyWithTransformAppendEmpty-12                   425280856                2.46 ns/op
BenchmarkSliceCopyWithTransformAppendSmall-12                     198354              6289 ns/op
BenchmarkSliceCopyWithTransformAppendMedium-12                      1251            858746 ns/op
BenchmarkSliceCopyWithTransformAppendLarge-12                        106          11412955 ns/op
BenchmarkSliceCopyWithTransformPreallocEmpty-12                 464395524                2.27 ns/op
BenchmarkSliceCopyWithTransformPreallocSmall-12                   181555              6556 ns/op
BenchmarkSliceCopyWithTransformPreallocMedium-12                    2072            735380 ns/op
BenchmarkSliceCopyWithTransformPreallocLarge-12                      204           5931064 ns/op
BenchmarkSliceCopyWithTransformIndexAndPreallocEmpty-12         489999499                2.18 ns/op
BenchmarkSliceCopyWithTransformIndexAndPreallocSmall-12           344100              3120 ns/op
BenchmarkSliceCopyWithTransformIndexAndPreallocMedium-12            4641            244839 ns/op
BenchmarkSliceCopyWithTransformIndexAndPreallocLarge-12              490           2252869 ns/op
```
