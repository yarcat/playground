package diff

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

func init() { rand.Seed(time.Now().UnixNano()) }

var (
	rand10      = makeSortedRandSlice(10)
	rand100     = makeSortedRandSlice(100)
	rand1000    = makeSortedRandSlice(1000)
	rand10000   = makeSortedRandSlice(10000)
	rand100000  = makeSortedRandSlice(100000)
	rand1000000 = makeSortedRandSlice(1000000)
	rand10000000 = makeSortedRandSlice(10000000)
)

func makeSortedRandSlice(n int) []int {
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = rand.Int()
	}
	sort.Ints(s)
	return s
}

func BenchmarkDiffBisect10(b *testing.B)      { benchmarkDiffBisect(b, rand10, rand10) }
func BenchmarkDiffBisect100(b *testing.B)     { benchmarkDiffBisect(b, rand100, rand100) }
func BenchmarkDiffBisect1000(b *testing.B)    { benchmarkDiffBisect(b, rand1000, rand1000) }
func BenchmarkDiffBisect10000(b *testing.B)   { benchmarkDiffBisect(b, rand10000, rand10000) }
func BenchmarkDiffBisect100000(b *testing.B)  { benchmarkDiffBisect(b, rand100000, rand100000) }
func BenchmarkDiffBisect1000000(b *testing.B) { benchmarkDiffBisect(b, rand1000000, rand1000000) }
func BenchmarkDiffBisect10000000(b *testing.B) { benchmarkDiffBisect(b, rand10000000, rand10000000) }

func BenchmarkDiffMap10(b *testing.B)      { benchmarkDiffMap(b, rand10, rand10) }
func BenchmarkDiffMap100(b *testing.B)     { benchmarkDiffMap(b, rand100, rand100) }
func BenchmarkDiffMap1000(b *testing.B)    { benchmarkDiffMap(b, rand1000, rand1000) }
func BenchmarkDiffMap10000(b *testing.B)   { benchmarkDiffMap(b, rand10000, rand10000) }
func BenchmarkDiffMap100000(b *testing.B)  { benchmarkDiffMap(b, rand100000, rand100000) }
func BenchmarkDiffMap1000000(b *testing.B) { benchmarkDiffMap(b, rand1000000, rand1000000) }
func BenchmarkDiffMap10000000(b *testing.B) { benchmarkDiffMap(b, rand10000000, rand10000000) }

func benchmarkDiffBisect(b *testing.B, sa, sb []int) {
	for n := 0; n < b.N; n++ {
		diffBisect(sa, sb)
	}
}
func benchmarkDiffMap(b *testing.B, sa, sb []int) {
	for n := 0; n < b.N; n++ {
		diffMap(sa, sb)
	}
}

func find(b []int, x int) int {
	min, max := 0, len(b)
	for min != max {
		i := (min + max) / 2
		if v := b[i]; v == x {
			return i
		} else if v < x {
			min = i + 1
		} else {
			max = i
		}
	}
	return -1
}

func diffBisect(a, b []int) (diff []int) {
	if len(a) > len(b) {
		a, b = b, a
	}
	for _, x := range a {
		if find(b, x) < 0 {
			diff = append(diff, x)
		}
	}
	return
}

func diffMap(a, b []int) (diff []int) {
	if len(a) > len(b) {
		a, b = b, a
	}
	m := make(map[int]struct{}, len(a))
	for _, x := range a {
		m[x] = struct{}{}
	}
	for _, x := range b {
		if _, ok := m[x]; !ok {
			diff = append(diff, x)
		}
	}
	return
}
