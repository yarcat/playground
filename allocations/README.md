# Alloc-free

## Description

An attempt to see what leads to alloc-free code in Golang and what not.

Set: The most straight-forward alloc-free implementation where nothing goes outside of the client.

     Pros: The major plus is that it's very easy to implement.

     Cons: It's very hard to handle different data types. Using those would require interfaces or
           function pointers... And it makes it alloc S: See `Exec`.

Exec: An attempt to use the "options pattern" to handle different data types. Unfortunately those
      dynamic functions allow and don't inline nicely. I may work w/o `ExecArgFunc` factories though,
      since Go does nice job inlining simple functions.

     Cons: Freaking allocs S:

Exec2: Expose mutator as a pointer.

       Cons: Mutator is alloc'ed... Functions look weird (but I like it).

Exec3: Expose mutator as a value, return new state when done.

       Pros: Zero allocs!

       Cons: Funcstions look weird (but I still like it).

Exec4: Pre-allocate the mutator.

       Pros: Zero allocs!

       Cons: Tiny bit higher memory consumption (seems to be 1 byte + pointer size).

## Current Results

```
go test -bench=. -benchmem -memprofile=mem.out .
goos: linux
goarch: amd64
pkg: github.com/yarcat/playground/allocations
BenchmarkSet/Set-8      33414662                31.71 ns/op            0 B/op          0 allocs/op
BenchmarkExec/Exec-8             9790566               121.7 ns/op            56 B/op          2 allocs/op
BenchmarkExec2/Exec2-8          17658217                65.39 ns/op           24 B/op          1 allocs/op
BenchmarkExec3/Exec3-8          29295375                35.57 ns/op            0 B/op          0 allocs/op
BenchmarkExec4/Exec4-8          33040213                34.83 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/yarcat/playground/allocations        5.935s
```
