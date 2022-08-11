package sort_test

import (
	"math/rand"
	"testing"

	"golang.org/x/exp/constraints"
)

var s []int

func init() {
	// rand.Seed(time.Now().UnixMicro())
	s = make([]int, 10_000)
	for i := range s {
		s[i] = rand.Intn(10_000)
	}
}

func benchmarkSort(b *testing.B, sort func([]int)) {
	s2 := make([]int, len(s))
	for n := 0; n < b.N; n++ {
		copy(s2, s)
		sort(s2)
	}
}
func BenchmarkSort(b *testing.B)       { benchmarkSort(b, Sort[int]) }
func BenchmarkSortBubble(b *testing.B) { benchmarkSort(b, SortBubble[int]) }

func Sort[T constraints.Ordered](s []T) {
	x := make([]T, len(s))
	slices, from, to := [2][]T{s, x}, 0, 1
	for seg := 1; seg < len(s); seg *= 2 {
		out := 0
		for out < len(s) {
			a, b := out, out+seg
			A, B := seg, seg
			if out+A >= len(s) {
				A, B = len(s)-out, 0
			} else if out+A+B > len(s) {
				B = len(s) - A - out
			}
			for A > 0 && B > 0 {
				if slices[from][a] <= slices[from][b] {
					slices[to][out] = slices[from][a]
					a++
					A--
				} else {
					slices[to][out] = slices[from][b]
					b++
					B--
				}
				out++
			}
			if A > 0 {
				copy(slices[to][out:], slices[from][a:a+A])
				out += A
			} else {
				copy(slices[to][out:], slices[from][b:b+B])
				out += B
			}
		}
		from, to = to, from
	}
	if from == 1 {
		copy(s, slices[1])
	}
}

func SortBubble[T constraints.Ordered](s []T) {
	for i := len(s); i > 0; i-- {
		swapped := false
		for j := 1; j < i; j++ {
			if s[j-1] > s[j] {
				s[j-1], s[j] = s[j], s[j-1]
				swapped = true
			}
		}
		if !swapped {
			return
		}
	}
}
