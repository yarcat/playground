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

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

func Sort[T constraints.Ordered](s []T) {
	x := make([]T, len(s))
	from, to := &s, &x
	merge := func(seg, out int) int {
		a, b := out, out+seg
		A, B := min(seg, len(s)-a), min(seg, len(s)-b)
		for ; A > 0 && B > 0; out++ {
			if (*from)[a] <= (*from)[b] {
				(*to)[out] = (*from)[a]
				a, A = a+1, A-1
			} else {
				(*to)[out] = (*from)[b]
				b, B = b+1, B-1
			}
		}
		if A > 0 {
			out += copy((*to)[out:], (*from)[a:a+A])
		} else {
			out += copy((*to)[out:], (*from)[b:b+B])
		}
		return out
	}
	for seg := 1; seg < len(s); seg *= 2 {
		for out := 0; out < len(s); {
			out = merge(seg, out)
		}
		from, to = to, from
	}
	if from == &x {
		copy(s, *from)
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
