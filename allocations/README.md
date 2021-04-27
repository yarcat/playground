# Alloc-free

An attempt to see what leads to alloc-free code in Golang and what not.

## Description

### Set

The most straight-forward alloc-free implementation where nothing goes outside of the client.

Pros: The major plus is that it's very easy to implement.

Cons: It's very hard to handle different data types. Using those would require interfaces or
      function pointers... And it makes it alloc S: See `Exec`.

### Exec/Factor

An attempt to use the "options pattern" to handle different data types. Unfortunately those
dynamic functions allow and don't inline nicely.

Cons: Freaking allocs S:


### Exec/Lambda

Using "options pattern" w/o factories.

Pros: No allocs!

Cons: Too many lambdas.

### Exec3

Expose mutator as a pointer.

Cons: Mutator is alloc'ed... Functions look weird (but I like it).

### Exec3

Expose mutator as a value, return new state when done.

Pros: Zero allocs!

Cons: Funcstions look weird (but I still like it).

### Exec4

Pre-allocate the mutator.

Pros: Zero allocs!

Cons: Tiny bit higher memory consumption (seems to be 1 byte + pointer size).

## Current Results

```
go test -bench=. -benchmem .
goos: linux
goarch: amd64
pkg: github.com/yarcat/playground/allocations
BenchmarkSet/Set-8      37535660                31.68 ns/op            0 B/op          0 allocs/op
BenchmarkExec/Factory-8                  9246202               125.2 ns/op            56 B/op          2 allocs/op
BenchmarkExec/Lambda-8                  29503152                40.07 ns/op            0 B/op          0 allocs/op
BenchmarkExec2/Exec2-8                  18385735                64.96 ns/op           24 B/op          1 allocs/op
BenchmarkExec3/Exec3-8                  31633879                35.63 ns/op            0 B/op          0 allocs/op
BenchmarkExec4/Exec4-8                  33905007                35.81 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/yarcat/playground/allocations        11.213s
```
