package main

import "testing"

var input = generateInput(1000)

func generateInput(n int) (in []T1) {
	for i := 0; i < n; i++ {
		in = append(in, i)
	}
	return
}

func benchmarkConcurrentFilterMap(b *testing.B, filterMap func([]T1, func(T1) (T2, bool), int) []T2) {
	for n := 0; n < b.N; n++ {
		filterMap(input, f, 10)
	}
}

func BenchmarkConcurrentFilterMap1(b *testing.B) {
	benchmarkConcurrentFilterMap(b, ConcurrentFilterMap1)
}
func BenchmarkConcurrentFilterMap2(b *testing.B) {
	benchmarkConcurrentFilterMap(b, ConcurrentFilterMap2)
}
func BenchmarkConcurrentFilterMap3(b *testing.B) {
	benchmarkConcurrentFilterMap(b, ConcurrentFilterMap3)
}
