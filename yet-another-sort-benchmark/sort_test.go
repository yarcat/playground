package sort_test

import (
	"math/rand"
	"reflect"
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

func testSort(t *testing.T, sort func([]int)) {
	for _, tc := range []struct {
		name        string
		input, want []int
	}{
		{"nil", nil, nil},
		{"empty", []int{}, []int{}},
		{"single", []int{123}, []int{123}},
		{"two_sorted", []int{1, 2}, []int{1, 2}},
		{"two_reversed", []int{2, 1}, []int{1, 2}},
		{"three_random", []int{2, 3, 1}, []int{1, 2, 3}},
		{"three_sorted", []int{1, 2, 3}, []int{1, 2, 3}},
		{"three_reversed", []int{3, 2, 1}, []int{1, 2, 3}},
		{"four_reversed", []int{4, 3, 2, 1}, []int{1, 2, 3, 4}},
		{"ten_random", []int{4, 3, 7, 1, 0, 5, 2, 6, 9, 8}, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"ten_sorted", []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"ten_reversed", []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got []int
			if tc.input != nil {
				got = make([]int, len(tc.input))
				copy(got, tc.input)
			}
			sort(got)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got = %v, want = %v", got, tc.want)
			}
		})
	}
}

func TestSort(t *testing.T)       { testSort(t, Sort[int]) }
func TestSortBubble(t *testing.T) { testSort(t, SortBubble[int]) }

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

func Sort[T constraints.Ordered](s []T) {
	from, to, fromS := s, make([]T, len(s)), true
	merge := func(seg, out int) int {
		a, b := out, out+seg
		A, B := min(seg, len(s)-a), min(seg, len(s)-b)
		for ; A > 0 && B > 0; out++ {
			if from[a] <= from[b] {
				to[out] = from[a]
				a, A = a+1, A-1
			} else {
				to[out] = from[b]
				b, B = b+1, B-1
			}
		}
		if A > 0 {
			out += copy(to[out:], from[a:a+A])
		} else {
			out += copy(to[out:], from[b:b+B])
		}
		return out
	}
	for seg := 1; seg < len(s); seg *= 2 {
		for out := 0; out < len(s); {
			out = merge(seg, out)
		}
		from, to, fromS = to, from, !fromS
	}
	if !fromS {
		copy(s, from)
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
