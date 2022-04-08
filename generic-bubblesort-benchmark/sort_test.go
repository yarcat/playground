package main

import (
	"math/rand"
	"sort"
	"testing"
)

var input []int

func init() {
	rand.Seed(17)

	input = make([]int, 1000)
	for i := range input {
		input[i] = rand.Int()
	}
}

func benchmarkSort(b *testing.B, sort func([]int)) {
	data := make([]int, 1000)
	for n := 0; n < b.N; n++ {
		copy(data, input)
		sort(data)
	}
}

func BenchmarkSortGen(b *testing.B) {
	benchmarkSort(b, func(input []int) { SortGen(sort.IntSlice(input)) })
}

func BenchmarkSortInt(b *testing.B) {
	benchmarkSort(b, func(input []int) { SortInt(sort.IntSlice(input)) })
}

